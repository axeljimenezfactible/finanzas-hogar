"""Regenera docs/contratos/diccionario-datos.md desde los comentarios @ de TODAS las migraciones
(backend/db/migrations/*.sql, en orden) y la introspección de una base en memoria migrada con todas.
Uso: python3 -I backend/db/tools/gen_dict.py [salida.md]
Sin argumento reescribe docs/contratos/diccionario-datos.md; con argumento escribe ahí (para comparar
con cmp). El texto fijo (antes de la primera sección de tablas y desde «## Datos de referencia») se
toma siempre del diccionario del repositorio.
Anotaciones: las de 00001 (@table, @title, @entity, @desc, @used, @idx, comentario por columna
«descripción | regla» y @rule) y, desde 00002, @alter <tabla> (columnas de ALTER TABLE ... ADD COLUMN),
@used_new <texto> (quién usa las columnas nuevas) y @after <tabla> (lugar de una tabla nueva).
Origen: Anexo A de docs/estado/wf0/S1-modelo-datos.md (WF-0); adaptado a varias migraciones en P1."""
import glob, os, re, sqlite3, sys
R = os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))))
DIC = os.path.join(R, 'docs', 'contratos', 'diccionario-datos.md')
out_path = sys.argv[1] if len(sys.argv) > 1 else DIC
MIGS = sorted(glob.glob(os.path.join(R, 'backend', 'db', 'migrations', '*.sql')))
UPS = [open(p, encoding='utf-8').read().split('-- +goose Down', 1)[0] for p in MIGS]
db = sqlite3.connect(':memory:')
for u in UPS:
    db.executescript(u)
mig = '\n'.join(UPS)

meta, order, sections = {}, [], []
cur = in_table = pending_idx = None
for line in mig.splitlines():
    if line.startswith('-- +goose Down'):
        break
    hm = re.match(r'^-- ([A-F](?: \(cont\.\))?\. .+)$', line)
    if hm:
        sections.append((hm.group(1), len(order)))
        continue
    m = re.match(r'^-- @(table|title|entity|desc|used|idx|alter|used_new|after) (.*)$', line)
    if m:
        k, v = m.groups()
        if k == 'table':
            cur = v; meta[cur] = dict(cols={}, rules=[], idx={}); order.append(cur)
        elif k == 'alter':
            cur = v
        elif k == 'after':
            order.remove(cur); p = order.index(v); order.insert(p + 1, cur)
            sections = [(n, i + 1 if i > p else i) for n, i in sections]
        elif k == 'idx':
            pending_idx = v
        else:
            meta[cur][k] = v
        continue
    m = re.match(r'^ALTER TABLE (\w+) ADD COLUMN ([a-z_][a-z0-9_]*) (?:INTEGER|TEXT|REAL|BLOB)\b', line)
    if m:
        comment = line.split(' -- ', 1)[1].strip() if ' -- ' in line else ''
        desc, _, rule = comment.partition(' | ')
        meta[m.group(1)]['cols'][m.group(2)] = (desc.strip(), rule.strip())
        continue
    m = re.match(r'^CREATE TABLE (\w+) \(', line)
    if m:
        in_table = m.group(1); continue
    if in_table:
        if line.startswith(') STRICT;'):
            in_table = None; continue
        s = line.strip()
        comment = line.split(' -- ', 1)[1].strip() if ' -- ' in line else None
        mc = re.match(r'^([a-z_][a-z0-9_]*) (INTEGER|TEXT|REAL|BLOB)\b', s)
        if mc and not s.startswith(('CHECK', 'UNIQUE')):
            desc, _, rule = (comment or '').partition(' | ')
            meta[in_table]['cols'][mc.group(1)] = (desc.strip(), rule.strip())
        elif comment and '@rule' in comment:
            meta[in_table]['rules'].append(comment.split('@rule', 1)[1].strip())
        continue
    m = re.match(r'^CREATE (?:UNIQUE )?INDEX (\w+) ON (\w+)', line)
    if m:
        meta[m.group(2)]['idx'][m.group(1)] = pending_idx; pending_idx = None

STD = {
    'id': 'Identificador interno (alias de rowid).',
    'created_at': 'Alta del registro, instante UTC ISO 8601 con «Z».',
    'updated_at': 'Última modificación, instante UTC ISO 8601. La mantiene el disparador `trg_<tabla>_updated_at`.',
    'is_active': '0 = archivado: deja de ofrecerse en capturas nuevas y se conserva para la historia.',
    'notes': 'Notas libres.',
}
ON_DEL = {'RESTRICT': 'RESTRICT (no se puede borrar la madre)', 'CASCADE': 'CASCADE (se borra con la madre)',
          'SET NULL': 'SET NULL (queda sin liga)', 'NO ACTION': 'NO ACTION'}
def esc(x):
    return (x or '').replace('|', '/').replace('\n', ' ')

fk_all = {t: db.execute(f'PRAGMA foreign_key_list({t})').fetchall() for t in order}
inbound = {t: [] for t in order}
for t, fks in fk_all.items():
    for fk in fks:
        inbound.setdefault(fk[2], []).append(f'`{t}.{fk[3]}`')

def table_md(t):
    m = meta[t]; o = []
    o.append(f"### `{t}` · {m['title']}\n")
    o.append(f"- **Entidad del plan (§3.4.1):** {m['entity']}")
    o.append(f"- **Descripción:** {m['desc']}")
    o.append(f"- **Usado por:** {m['used']}" + ('' if m.get('used_new') else '\n'))
    if m.get('used_new'):
        o.append(f"- **Usado por las columnas nuevas de 1.2:** {m['used_new']}\n")
    cols = db.execute(f'PRAGMA table_xinfo({t})').fetchall()
    fks = {fk[3]: fk for fk in fk_all[t]}
    uniq1 = set()
    for il in db.execute(f'PRAGMA index_list({t})').fetchall():
        if il[2] and il[3] == 'u' and not il[4]:
            cs = db.execute(f'PRAGMA index_info({il[1]})').fetchall()
            if len(cs) == 1:
                uniq1.add(cs[0][2])
    o.append('| nombre_codigo | tipo | nulo | descripción | regla |')
    o.append('|---|---|---|---|---|')
    for cid, name, typ, notnull, dflt, pk, hidden in cols:
        desc, rule = m['cols'].get(name, ('', ''))
        desc = desc or STD.get(name, '')
        parts = []
        if pk: parts.append('PK')
        if name in fks:
            fk = fks[name]; parts.append(f"FK → `{fk[2]}.{fk[4]}` al borrar: {fk[6]}")
        if name in uniq1: parts.append('único')
        if dflt is not None:
            parts.append('def: ahora (UTC)' if 'strftime' in dflt else f'def: {dflt}')
        if rule: parts.append(rule)
        if name in ('created_at', 'updated_at') and not rule: parts.append('ISO 8601 UTC')
        o.append(f"| `{name}` | {typ} | {'no' if (notnull or pk) else 'sí'} | {esc(desc)} | {esc('; '.join(parts))} |")
    if m['rules']:
        o.append('\n**Reglas de tabla (CHECK / UNIQUE):**')
        o += [f'- {esc(r)}' for r in m['rules']]
    o.append('\n**Relaciones:**')
    if fk_all[t]:
        o += [f"- `{fk[3]}` → `{fk[2]}.{fk[4]}` · al borrar la madre: {ON_DEL.get(fk[6], fk[6])}" for fk in fk_all[t]]
    else:
        o.append('- Sin llaves foráneas salientes.')
    if inbound.get(t):
        o.append(f"- La referencian: {', '.join(sorted(set(inbound[t])))}")
    idx_rows = db.execute("SELECT name, sql FROM sqlite_master WHERE type='index' AND tbl_name=? ORDER BY name", (t,)).fetchall()
    if idx_rows:
        o.append('\n**Índices:**\n')
        o.append('| índice | definición | propósito |')
        o.append('|---|---|---|')
        for name, sql in idx_rows:
            if sql is None:
                cs = [c[2] for c in db.execute(f'PRAGMA index_info({name})').fetchall()]
                o.append(f"| `{name}` | UNIQUE ({', '.join(cs)}) automático | Lo crea SQLite por la restricción UNIQUE o la columna única. |")
            else:
                d = sql.split(' ON ', 1)[1]
                o.append(f"| `{name}` | {'UNIQUE ' if sql.startswith('CREATE UNIQUE') else ''}`{esc(d)}` | {esc(m['idx'].get(name) or '')} |")
    trg = db.execute("SELECT name, sql FROM sqlite_master WHERE type='trigger' AND tbl_name=? ORDER BY name", (t,)).fetchall()
    if trg:
        o.append('\n**Disparadores:**')
        for name, sql in trg:
            mm = re.search(r"RAISE\(ABORT, '([^']+)'\)", sql)
            o.append(f"- `{name}`: " + (f"rechaza con «{mm.group(1)}»." if mm else 'pone updated_at en el instante actual si el UPDATE no lo trae.'))
    return '\n'.join(o) + '\n'

body = []
sec_at = {i: name for name, i in sections}
for i, t in enumerate(order):
    if i in sec_at:
        body.append(f'## {sec_at[i]}\n')
    body.append(table_md(t))

core = [t for t in order if not meta[t]['entity'].startswith('(complementaria)')]
comp = [t for t in order if meta[t]['entity'].startswith('(complementaria)')]
map_rows = '\n'.join(f"| {meta[t]['entity']} | `{t}` | {meta[t]['title']} |" for t in core)
comp_rows = '\n'.join(f"| `{t}` | {meta[t]['title']} | {meta[t]['entity'].replace('(complementaria) ', '')} |" for t in comp)
n_tab = len(order); n_idx = db.execute("SELECT count(*) FROM sqlite_master WHERE type='index'").fetchone()[0]
n_trg = db.execute("SELECT count(*) FROM sqlite_master WHERE type='trigger'").fetchone()[0]
n_upd = db.execute("SELECT count(*) FROM sqlite_master WHERE type='trigger' AND name LIKE 'trg\\_%\\_updated\\_at' ESCAPE '\\'").fetchone()[0]
# Plantilla: el propio diccionario (texto fijo antes de la primera sección de tablas y desde «## Datos de referencia»).
old = open(DIC, encoding='utf-8').read()
first = '\n## ' + sections[0][0] + '\n'
intro = old[:old.index(first) + 1]
tail = old[old.index('## Datos de referencia que crea la migración'):]
intro = re.sub(r'^Resumen: .*$', f'Resumen: **{n_tab} tablas STRICT**, {n_idx} índices (incluye los automáticos de UNIQUE) y {n_trg} disparadores ({n_trg - n_upd} de reglas de negocio y {n_upd} de `updated_at`).', intro, count=1, flags=re.M)
def put_rows(text, header, rows):
    k = text.index(header) + len(header)
    end = text.index('\n\n', k)
    return text[:k] + rows + text[end:]
intro = put_rows(intro, '| Entidad del plan (§3.4.1) | nombre_codigo | Título |\n|---|---|---|\n', map_rows)
intro = put_rows(intro, '| nombre_codigo | Título | Por qué existe |\n|---|---|---|\n', comp_rows)
open(out_path, 'w', encoding='utf-8').write(intro + '\n'.join(body) + '\n' + tail)
print('tablas', n_tab, 'núcleo', len(core), 'complementarias', len(comp), 'índices', n_idx, 'disparadores', n_trg)
missing = [(t, c) for t in order for c in [x[1] for x in db.execute(f'PRAGMA table_xinfo({t})')] if not (meta[t]['cols'].get(c, ('',))[0] or c in STD)]
print('columnas sin descripción:', missing)
