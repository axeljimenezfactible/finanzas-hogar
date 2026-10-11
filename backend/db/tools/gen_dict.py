"""Regenera docs/contratos/diccionario-datos.md desde los comentarios @ de TODAS las migraciones
(backend/db/migrations/*.sql, en orden) y la introspección de una base en memoria migrada con todas.
Uso: python3 -I backend/db/tools/gen_dict.py [salida.md]
Sin argumento reescribe docs/contratos/diccionario-datos.md; con argumento escribe ahí (para comparar
con cmp). El texto fijo (antes de la primera sección de tablas y desde «## Datos de referencia») se
toma siempre del diccionario del repositorio.
Anotaciones: las de 00001 (@table, @title, @entity, @desc, @used, @idx, comentario por columna
«descripción | regla» y @rule) y, desde 00002, @alter <tabla> (columnas de ALTER TABLE ... ADD COLUMN),
@used_new <texto> (quién usa las columnas nuevas) y @after <tabla> (lugar de una tabla nueva).
Desde 00003 (esquema 1.3): @param <clave> | tipo | unidad | valor inicial | decisión | sección | migración y
@param_sin_valor <parámetro> | para qué | por qué no se siembra (inventario de parámetros de K44); genera la
sección «Parámetros operativos de app_setting (K44)», comprueba cada valor inicial contra lo sembrado y pone
la versión del esquema (la de «esquema a → b» en la cabecera de la última migración) en el encabezado.
También desde 00003: @col_note <tabla>.<columna> | nota (se añade a la descripción de la columna) y, en la tabla
«Ajustes: 14 secciones», cada fila nombra las claves @param de su sección con su valor inicial (idempotente).
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
VERS = [(re.search(r'esquema \S+ → (\d+\.\d+)', u[:600]) or [None, '1.0'])[1] for u in UPS]
VER = VERS[-1]
PARAMS = [[c.strip() for c in l.split(' | ')] for l in re.findall(r'^-- @param (.+)$', mig, re.M)]
PARAMS_NV = [[c.strip() for c in l.split(' | ')] for l in re.findall(r'^-- @param_sin_valor (.+)$', mig, re.M)]

meta, order, sections = {}, [], []
cur = in_table = pending_idx = None
for line, ver in ((ln, VERS[k]) for k, u in enumerate(UPS) for ln in u.splitlines()):
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
        elif k == 'used_new':
            meta[cur].setdefault('used_new', []).append((ver, v))
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
    for k, (ver, txt) in enumerate(m.get('used_new', [])):
        o.append(f"- **Usado por las columnas nuevas de {ver}:** {txt}" + ('\n' if k == len(m['used_new']) - 1 else ''))
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

# Notas de 1.3 sobre columnas de migraciones anteriores (-- @col_note tabla.columna | nota).
COLN = [[c.strip() for c in l.split(' | ', 1)] for l in re.findall(r'^-- @col_note (.+)$', mig, re.M)]
for tc, nota in COLN:
    t_, c_ = tc.split('.', 1); d_, r_ = meta[t_]['cols'][c_]
    meta[t_]['cols'][c_] = (d_ + '. ' + nota + '.', r_)
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
# Esquema 1.3: encabezado, cambios, semilla de la 00003, parámetros de K44 y fila 13 de Ajustes (todo idempotente).
nombres = [os.path.basename(p) for p in MIGS]
mig_txt = ', '.join(f'`backend/db/migrations/{n}`' for n in nombres[:-1]) + f' y `backend/db/migrations/{nombres[-1]}`'
intro = re.sub(r'^\*\*Versión:\*\* .*$', lambda _: f'**Versión:** {VER} · migraciones {mig_txt} · 11-oct-2026 (respuestas de la puerta 1 y de la puerta 1b, `docs/plan/decisiones.md` §K)', intro, count=1, flags=re.M)
intro = re.sub(r'^(\*\*Autor:\*\* [^\n]*?S1-esquema-1\.2 \(flujo P1\))(?:; versión 1\.3: [^\n]*)?$', lambda mm: mm.group(1) + '; versión 1.3: agente S1-esquema-1.3 (flujo P1c)', intro, count=1, flags=re.M)
intro = intro.replace('`-- @used_new` y `-- @after`)', '`-- @used_new` y `-- @after`; desde 1.3 `-- @param` y `-- @param_sin_valor`)')
C13 = ('**Cambios en 1.3** (migración `00003_puerta1b.sql`, respuestas de Axel de la puerta 1b, `docs/plan/decisiones.md` K36 a K44): '
       'columnas nuevas `health_profile.goal_adjust_lose_kcal` y `.goal_adjust_gain_kcal` (K36, ajuste de la meta de calorías por persona; NULL = valor de la casa); '
       'claves nuevas de `app_setting` sembradas con su valor decidido (K44): `health.goal_adjust_lose_kcal` = 0 y `health.goal_adjust_gain_kcal` = 0 (K36), '
       '`price.zero_weight_after_days` = 90 (K22, K40), `auth.password_min_length` = 12 y `price.profeco_import_max_gb` = 1 (K6), `attachment.max_size_mb` = 10 y `comparator.timeout_seconds` = 5 (K43), y `alerts.budget_warning_pct` = 80, `alerts.budget_limit_pct` = 100 y `alerts.payment_due_days_before` = [1, 0] (decisión F), que salen de `notification.params_json` (`payment_due` deja `days_before` y `budget_threshold` deja `thresholds_pct`; el Down los repone); '
       'disparadores `trg_app_setting_param_int_ins` y `trg_app_setting_param_int_upd` (los parámetros enteros solo aceptan un entero JSON mayor o igual a 0) y `trg_app_setting_param_int_list_ins` y `trg_app_setting_param_int_list_upd` (los días antes del pago próximo, una lista JSON de enteros mayores o iguales a 0); '
       'sección nueva «Parámetros operativos de app_setting (K44)». Ninguna tabla nueva; el DDL de 00001 y 00002 no cambia. '
       'Corrección de la revisión S1-r1 de P1c (`docs/estado/p1c/REVISION-S1-r1.md`, agente correccion-S1-r1): H-01, los umbrales del presupuesto y los días antes del pago próximo son claves de `app_setting` (K44); '
       'H-02, cada fila de «Ajustes: 14 secciones» nombra los parámetros de su sección con su valor inicial. '
       'Corrección de la revisión S1-r2 de P1c (`docs/estado/p1c/REVISION-S1-r2.md`, agente correccion-S1-r2): M-01, el inventario lista `weekday` (7) y `day_of_month` (1) de la decisión F, que siguen en `notification.params_json`; dónde viven es PROVISIONAL, bloqueo `S1-parametros-avisos`. Ningún valor cambia.')
intro = re.sub(r'^\*\*Cambios en 1\.3\*\*.*\n', '', intro, flags=re.M)
k = intro.index('**Correcciones de la revisión S1-r1**')
intro = intro[:k] + C13 + '\n' + intro[k:]
seeded = dict(db.execute('SELECT setting_key, value_json FROM app_setting WHERE user_account_id IS NULL').fetchall())
desc_s = dict(db.execute('SELECT setting_key, description FROM app_setting WHERE user_account_id IS NULL').fetchall())
mal_p = [p[0] for p in PARAMS if seeded.get(p[0]) != p[3]]
L = ['## Parámetros operativos de app_setting (K44)',
     f'Inventario de K44 («todo parámetro operativo existe y es configurable»), generado desde las anotaciones `-- @param` y `-- @param_sin_valor` de `{nombres[-1]}` y comprobado contra lo sembrado en la base migrada ({len(PARAMS) - len(mal_p)} de {len(PARAMS)} valores iniciales coinciden). Todos son del hogar (`user_account_id` NULL), enteros JSON mayores o iguales a 0 (disparadores `trg_app_setting_param_int_ins` y `_upd`), salvo `alerts.payment_due_days_before`, una lista JSON de enteros mayores o iguales a 0 (`trg_app_setting_param_int_list_ins` y `_upd`), y sin mínimos ni máximos en el esquema: esos los propone S1 y se aprueban en la puerta 2. La sección de Ajustes es la PROPUESTA del orquestador en K44. Regla de negocio: RN-29.',
     '', '| clave | tipo | unidad | valor inicial | decisión | sección de Ajustes | migración | descripción sembrada |', '|---|---|---|---|---|---|---|---|']
L += [f'| `{p[0]}` | {p[1]} | {p[2]} | {p[3]} | {p[4]} | {p[5]} | {p[6]} | {esc(desc_s.get(p[0]))} |' for p in PARAMS]
L += ['', '**Existen sin valor sembrado o viven en otra tabla** (no se inventa un valor):', '', '| parámetro | para qué | por qué no se siembra |', '|---|---|---|']
L += [f'| {esc(p[0])} | {esc(p[1])} | {esc(p[2])} |' for p in PARAMS_NV]
PSEC = '\n'.join(L) + '\n\n'
if '## Parámetros operativos de app_setting (K44)' in tail:
    a = tail.index('## Parámetros operativos de app_setting (K44)'); b = tail.index('## Claves de app_setting previstas sin valor')
    tail = tail[:a] + tail[b:]
b = tail.index('## Claves de app_setting previstas sin valor')
tail = tail[:b] + PSEC + tail[b:]
S03 = '- `app_setting` sembrado por la migración 00003 (puerta 1b, K44): ' + ', '.join(f'`{p[0]}` = {p[3]} ({p[4]})' for p in PARAMS if p[6] == '00003') + '.\n'
tail = re.sub(r'^- `app_setting` sembrado por la migración 00003 .*\n', '', tail, flags=re.M)
k = tail.index('\n', tail.index('- `app_setting` sembrado por la migración 00002')) + 1
NOT13 = ('- `notification` en la migración 00003 (K44): ' + ' y '.join(f'`{c}` queda con `{pj}`' for c, pj in db.execute("SELECT code, params_json FROM notification WHERE code IN ('payment_due', 'budget_threshold') ORDER BY code"))
         + '; los días antes y los umbrales viven en `alerts.payment_due_days_before`, `alerts.budget_warning_pct` y `alerts.budget_limit_pct` de `app_setting` (el Down de la 00003 los repone).\n')
tail = re.sub(r'^- `notification` en la migración 00003 .*\n', '', tail, flags=re.M)
tail = tail[:k] + S03 + NOT13 + tail[k:]
tail = tail.replace('`health_profile.height_mm`, `.birth_date`, `.sex`, `.activity_level`, `.goal`, `.daily_kcal_target`',
                    '`health_profile.height_mm`, `.birth_date`, `.sex`, `.activity_level`, `.goal`, `.goal_adjust_lose_kcal`, `.goal_adjust_gain_kcal` (K36; NULL = valor de la casa en `app_setting` `health.goal_adjust_lose_kcal` y `health.goal_adjust_gain_kcal`), `.daily_kcal_target`')
# Fila 9 (1.3, K44) y, en cada fila de «Ajustes: 14 secciones», las claves @param de su sección con su valor inicial
# (corrección S1-r1 de P1c, H-02). Idempotente: el sufijo generado se quita antes de volver a ponerlo.
tail = tail.replace('`notification.params_json` (`days_before`, `thresholds_pct`, `scopes`, `weekday`)',
                    '`notification.params_json` (`scopes`, `weekday`; desde 1.3 los umbrales y los días antes son claves de `app_setting`, K44)')
SUF = ' · parámetros de K44 en `app_setting` (valor inicial, RN-29): '
a13 = tail.index('## Ajustes: 14 secciones'); b13 = tail.find('\n## ', a13 + 1); b13 = len(tail) if b13 < 0 else b13
filas, usados = [], set()
for ln in tail[a13:b13].split('\n'):
    fm = re.match(r'^\| (\d+) \| ([^|]+) \| ([^|]+) \| ([^|]+) \|$', ln)
    if fm:
        dest = fm.group(4).split(SUF)[0]
        ps = [p for p in PARAMS if fm.group(2).startswith(p[5])]
        usados |= {p[0] for p in ps}
        if ps: dest += SUF + ', '.join(f'`{p[0]}` ({p[3]})' for p in ps)
        ln = f'| {fm.group(1)} | {fm.group(2)} | {fm.group(3)} | {dest} |'
    filas.append(ln)
tail = tail[:a13] + '\n'.join(filas) + tail[b13:]
sin_fila = [p[0] for p in PARAMS if p[0] not in usados]
open(out_path, 'w', encoding='utf-8').write(intro + '\n'.join(body) + '\n' + tail)
print('parámetros sin fila en «Ajustes: 14 secciones»:', sin_fila)
if sin_fila: sys.exit(1)
print('versión', VER, '· parámetros', len(PARAMS), '· sin valor', len(PARAMS_NV), '· valores iniciales que no coinciden con lo sembrado:', mal_p)
if mal_p: sys.exit(1)
print('tablas', n_tab, 'núcleo', len(core), 'complementarias', len(comp), 'índices', n_idx, 'disparadores', n_trg)
missing = [(t, c) for t in order for c in [x[1] for x in db.execute(f'PRAGMA table_xinfo({t})')] if not (meta[t]['cols'].get(c, ('',))[0] or c in STD)]
print('columnas sin descripción:', missing)
