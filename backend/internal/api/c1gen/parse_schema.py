# Generador C1 · etapa 1: esquema de S1 -> schema.json (intermedio en C1GEN_WORK).
# Aplica TODAS las migraciones de backend/db/migrations/ en orden (sección Up de goose) en una base SQLite temporal
# dentro de C1GEN_WORK (bajo datos/tmp/) y lee de ella sqlite_master, PRAGMA table_info y PRAGMA foreign_key_list.
# Los metadatos de tabla (@table, @title, @entity, @desc, @used) y los comentarios de columna se leen de los archivos SQL:
# cada columna se busca en todas las migraciones (CREATE TABLE y ALTER TABLE ... ADD COLUMN); manda la más reciente.
# Ninguna ruta fija a una migración concreta (C1 v1.1).
# C1 v1.2 (puerta 1b, K44): también lee el inventario de parámetros operativos (anotaciones @param y @param_sin_valor de
# cualquier migración, cabecera incluida) y lo escribe en params.json; cada valor inicial se compara con la siembra real de
# app_setting en la base temporal (si no coincide, se detiene). Las notas @col_note se suman a la descripción de su columna.
import re, json, sqlite3, os, glob
from c1paths import ROOT, WORK
MIGDIR = os.path.join(ROOT, 'backend/db/migrations')
MIGS = sorted(glob.glob(os.path.join(MIGDIR, '[0-9]*.sql')))
if not MIGS:
    raise SystemExit('parse_schema: no hay migraciones en ' + MIGDIR)
NEW_PREFIX = re.compile(r'^Nueva en \d+\.\d+ \([^)]*\)\.\s*')  # «Nueva en 1.2 (K3).» es procedencia del diccionario, no descripción


def up_section(txt):
    return txt.split('-- +goose Down')[0].split('-- +goose Up', 1)[1]


def tag_value(block, key):
    m = re.search(r'-- @%s (.*)' % key, block)
    return m.group(1).strip() if m else ''


def enum_of(pattern, rest):
    e = re.search(pattern, rest)
    return [v.strip().strip("'") for v in e.group(1).split(',')] if e else None


meta, colcomments, enums = {}, {}, {}
dbpath = os.path.join(WORK, 'schema.db')
if os.path.exists(dbpath):
    os.remove(dbpath)
db = sqlite3.connect(dbpath)
for path in MIGS:
    up = up_section(open(path, encoding='utf-8').read())
    db.executescript(up.replace('-- +goose StatementBegin', '').replace('-- +goose StatementEnd', ''))
    for b in re.split(r'\n(?=-- @table )', up):
        m = re.match(r'-- @table (\w+)', b)
        if m:
            meta[m.group(1)] = {k: tag_value(b, k) for k in ('title', 'entity', 'desc', 'used')}
    for t, body in re.findall(r'CREATE TABLE (\w+) \((.*?)\n\) STRICT;', up, re.S):
        for line in body.split('\n'):
            mm = re.match(r'(\w+) (INTEGER|TEXT)\b(.*)', line.strip())
            if not mm:
                continue
            name, _typ, rest = mm.groups()
            comment = ''
            if '--' in rest:
                rest, comment = rest.split('--', 1)
            colcomments[(t, name)] = comment.strip()
            enums[(t, name)] = enum_of(r"CHECK \(%s IN \(([^)]*)\)" % name, rest)
    for t, name, rest, comment in re.findall(r'^ALTER TABLE (\w+) ADD COLUMN (\w+) (?:INTEGER|TEXT)\b(.*?);[ \t]*(?:--[ \t]*(.*))?$', up, re.M):
        colcomments[(t, name)] = NEW_PREFIX.sub('', comment.strip())
        enums[(t, name)] = enum_of(r"\b%s IN \(([^)]*)\)" % name, rest)
# C1 v1.2 (K44): inventario de parámetros operativos y notas de columna de 1.3.
params, params_sin_valor = [], []
for path in MIGS:
    head = open(path, encoding='utf-8').read().split('-- +goose Down')[0]
    for line in re.findall(r'^-- @param (.*)$', head, re.M):
        k, typ, unit, ini, dec, sec, mig = [x.strip() for x in line.split('|')]
        params.append(dict(key=k, type=typ, unit=unit, initial_json=ini, decision=dec, section=sec, migration=mig))
    for line in re.findall(r'^-- @param_sin_valor (.*)$', head, re.M):
        k, why, who = [x.strip().replace('|', '/') for x in line.split('|', 2)]
        params_sin_valor.append(dict(name=k, purpose=why, reason=who))
    for tc, note in re.findall(r'^-- @col_note (\w+\.\w+) \| (.*)$', head, re.M):
        t_, c_ = tc.split('.'); base, sep, rest = colcomments.get((t_, c_), '').partition('|')
        colcomments[(t_, c_)] = (base.strip() + ' ' + note.strip()).strip() + ((' |' + rest) if sep else '')
keys = [p['key'] for p in params]
if len(keys) != len(set(keys)):
    raise SystemExit('parse_schema: clave @param repetida')
for p in params:
    row = db.execute('SELECT value_json, description FROM app_setting WHERE setting_key = ? AND user_account_id IS NULL', (p['key'],)).fetchone()
    if row is None or json.loads(row[0]) != json.loads(p['initial_json']):
        raise SystemExit(f"parse_schema: {p['key']}: la siembra ({row and row[0]}) no coincide con @param ({p['initial_json']})")
    p['initial'] = json.loads(row[0]); p['seed_description'] = (row[1] or '').strip()
json.dump(dict(params=params, params_sin_valor=params_sin_valor), open(os.path.join(WORK, 'params.json'), 'w', encoding='utf-8'), ensure_ascii=False, indent=1)
db.commit()
names = [r[0] for r in db.execute("SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND name <> 'goose_db_version' ORDER BY rowid")]
tables = {}
for t in names:
    cols = [dict(name=name, type=typ, notnull=bool(notnull), default=dflt, pk=bool(pk), comment=colcomments.get((t, name), ''), enum=enums.get((t, name)))
            for _cid, name, typ, notnull, dflt, pk in db.execute(f'PRAGMA table_info({t})')]
    fks = [dict(col=r[3], ref=r[2], refcol=r[4], on_delete=r[6]) for r in db.execute(f'PRAGMA foreign_key_list({t})')]
    tables[t] = dict(meta=meta.get(t, dict(title='', entity='', desc='', used='')), cols=cols, fks=fks)
db.close()
os.remove(dbpath)
json.dump(tables, open(os.path.join(WORK, 'schema.json'), 'w', encoding='utf-8'), ensure_ascii=False, indent=1)
print(len(tables), 'tablas y', len(params), 'parámetros (K44) de', len(MIGS), 'migraciones:', ', '.join(os.path.basename(p) for p in MIGS))
