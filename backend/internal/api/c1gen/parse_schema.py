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
rangos, cerrados, colbase = [], [], {}
for path in MIGS:
    head = open(path, encoding='utf-8').read().split('-- +goose Down')[0]
    for line in re.findall(r'^-- @param (.*)$', head, re.M):
        k, typ, unit, ini, dec, sec, mig = [x.strip() for x in line.split('|')]
        params.append(dict(key=k, type=typ, unit=unit, initial_json=ini, decision=dec, section=sec, migration=mig))
    for line in re.findall(r'^-- @param_sin_valor (.*)$', head, re.M):
        k, why, who = [x.strip().replace('|', '/') for x in line.split('|', 2)]
        params_sin_valor.append(dict(name=k, purpose=why, reason=who))
    for line in re.findall(r'^-- @param_rango (.*)$', head, re.M):
        k, rng, dec = [x.strip() for x in line.split('|')]
        rangos.append(dict(key=k, note=rng, decision=dec))
    for line in re.findall(r'^-- @param_sin_valor_cerrado (.*)$', head, re.M):
        k, where = [x.strip() for x in line.split('|', 1)]
        cerrados.append(dict(name=k, closed=where.replace('|', '/'), migration=os.path.basename(path)[:5]))
    for tc, note in re.findall(r'^-- @col_note (\w+\.\w+) \| (.*)$', head, re.M):
        # C1 v1.3 (puerta 2): como gen_dict.py de S1 1.4, la última @col_note sustituye a la anterior (no se acumulan).
        t_, c_ = tc.split('.'); colbase.setdefault((t_, c_), colcomments.get((t_, c_), '')); base, sep, rest = colbase[(t_, c_)].partition('|')
        colcomments[(t_, c_)] = (base.strip() + ' ' + note.strip()).strip() + ((' |' + rest) if sep else '')
keys = [p['key'] for p in params]
if len(keys) != len(set(keys)):
    raise SystemExit('parse_schema: clave @param repetida')
for p in params:
    row = db.execute('SELECT value_json, description FROM app_setting WHERE setting_key = ? AND user_account_id IS NULL', (p['key'],)).fetchone()
    if row is None or json.loads(row[0]) != json.loads(p['initial_json']):
        raise SystemExit(f"parse_schema: {p['key']}: la siembra ({row and row[0]}) no coincide con @param ({p['initial_json']})")
    p['initial'] = json.loads(row[0]); p['seed_description'] = (row[1] or '').strip()
# C1 v1.3 (puerta 2): @param_rango se suma a su clave; @param_sin_valor_cerrado saca de la lista un @param_sin_valor anterior.
# El rango de cada clave se lee de la base, no se copia: cada valor se prueba con un UPDATE dentro de un SAVEPOINT que se
# deshace, así minimum, maximum y el patrón de la hora salen de los disparadores de S1 (00003 y 00004).
HHMM = '^([01][0-9]|2[0-3]):[0-5][0-9]$'
pk = {p['key']: p for p in params}
for r in rangos:
    if r['key'] not in pk: raise SystemExit(f"parse_schema: @param_rango de una clave sin @param: {r['key']}")
    pk[r['key']].setdefault('range_notes', []).append(dict(note=r['note'], decision=r['decision']))
_sv = {s['name'] for s in params_sin_valor}
for c in cerrados:
    if c['name'] not in _sv: raise SystemExit(f"parse_schema: @param_sin_valor_cerrado sin @param_sin_valor: {c['name']}")
params_sin_valor_cerrados = [dict(s, closed=c['closed'], closed_by=c['migration']) for c in cerrados for s in params_sin_valor if s['name'] == c['name']]
params_sin_valor = [s for s in params_sin_valor if s['name'] not in {c['name'] for c in cerrados}]
db.isolation_level = None
def _acepta(key, val):
    db.execute('SAVEPOINT rango')
    try:
        db.execute('UPDATE app_setting SET value_json = ? WHERE setting_key = ? AND user_account_id IS NULL', (json.dumps(val), key)); return True
    except sqlite3.DatabaseError:
        return False
    finally:
        db.execute('ROLLBACK TO rango'); db.execute('RELEASE rango')
TOPE = 10 ** 12
for p in params:
    k = p['key']
    if not _acepta(k, p['initial']): raise SystemExit(f'parse_schema: {k}: la base rechaza su propio valor inicial')
    if p['type'] == 'entero':
        if _acepta(k, -1) or _acepta(k, '1') or _acepta(k, 1.5): raise SystemExit(f'parse_schema: {k}: la base acepta un valor que no es entero >= 0')
        lo = 0 if _acepta(k, 0) else 1 if _acepta(k, 1) else None
        if lo is None: raise SystemExit(f'parse_schema: {k}: no acepta 0 ni 1; rango no previsto por el generador')
        p['minimum'] = lo
        if not _acepta(k, TOPE):
            a, b = p['initial'], TOPE
            while b - a > 1:
                m = (a + b) // 2
                if _acepta(k, m): a = m
                else: b = m
            p['maximum'] = a
    elif p['type'] == 'hora HH:MM':
        _rx = re.compile(HHMM)
        _malas = [s for s in [f'{h:02d}:{m:02d}' for h in range(30) for m in (0, 5, 59, 60, 99)] + ['9:00', '12:5', '1230', '08:00 ', ''] if bool(_rx.fullmatch(s)) != _acepta(k, s)]
        if _malas or _acepta(k, 900): raise SystemExit(f'parse_schema: {k}: el patrón {HHMM} no coincide con la base en {_malas[:5]}')
        p['pattern'] = HHMM
    elif p['type'] == 'lista de enteros':
        if _acepta(k, [-1]): raise SystemExit(f'parse_schema: {k}: la base acepta un elemento negativo')
        p['minimum'] = 0
    else:
        raise SystemExit(f"parse_schema: {k}: tipo de parámetro desconocido: {p['type']}")
json.dump(dict(params=params, params_sin_valor=params_sin_valor, params_sin_valor_cerrados=params_sin_valor_cerrados), open(os.path.join(WORK, 'params.json'), 'w', encoding='utf-8'), ensure_ascii=False, indent=1)
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
