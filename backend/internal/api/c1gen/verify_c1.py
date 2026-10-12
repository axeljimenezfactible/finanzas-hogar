# Verificación de C1 (criterios C1-01 a C1-15, C1-23, consistencia de la matriz §10 y paridad de rutas S6-03).
# Uso, desde cualquier carpeta (lee el repositorio que contiene este archivo):
#   uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py [--db BASE] [--root RAIZ]
#   --db BASE   base SQLite ya migrada (p. ej. goose -dir backend/db/migrations sqlite3 BASE up) para C1-03, abierta solo lectura;
#               sin --db se aplican en memoria, en orden, las secciones Up de backend/db/migrations/*.sql.
#   --root RAIZ otra copia del repositorio (por defecto, cuatro niveles arriba de esta carpeta).
# Código de salida: 0 si todo OK; 1 si hay alguna FALLA. Origen: anexo A del reporte S6a-contrato-c1 (corrección r1).
import argparse, json, re, os, sqlite3, random, sys, yaml
_ap = argparse.ArgumentParser(description='Verificación de C1 (C1-01 a C1-15, C1-23, matriz §10 y S6-03)')
_ap.add_argument('--root', default=os.path.join(os.path.dirname(os.path.abspath(__file__)), os.pardir, os.pardir, os.pardir, os.pardir))
_ap.add_argument('--db', default=None)
ARGS = _ap.parse_args()
ROOT = os.path.abspath(ARGS.root)
print('Repositorio:', ROOT)
OA = os.path.join(ROOT, 'docs/contratos/openapi.yaml'); C1 = os.path.join(ROOT, 'docs/contratos/C1-datos.md')
d = yaml.safe_load(open(OA, encoding='utf-8')); c1 = open(C1, encoding='utf-8').read()
fails = []
def ok(cond, msg):
    print(('OK   ' if cond else 'FALLA ') + msg)
    if not cond: fails.append(msg)
OPS = [(m.upper(), p, o) for p, ops in d['paths'].items() for m, o in ops.items()]
def deref(x):
    while isinstance(x, dict) and '$ref' in x:
        cur = d
        for part in x['$ref'][2:].split('/'): cur = cur[part]
        x = cur
    return x
# C1-01
print('C1-01 openapi =', d['openapi'], '· rutas =', len(d['paths']), '· operaciones =', len(OPS))
ok(str(d['openapi']).startswith('3.1.'), 'C1-01 versión 3.1.x')
# C1-02 (texto vigente del criterio: una N-xxx sin operación que no es Local, bloqueada o no, no cuenta como cubierta)
u1 = open(os.path.join(ROOT, 'docs/design/U1-arquitectura.md'), encoding='utf-8').read()
NSET = set(re.findall(r'^\| (N-\d{3}) \|', u1, re.M))
LOCAL_U1 = set(re.findall(r'^\| (N-\d{3}) \| [LE] \| .*\| Local \|$', u1, re.M))  # U1 declara «Local» en la última columna
_i, _j = u1.index('### 3.29'), u1.index('### 3.30')
_nums = sorted(int(n[2:]) for n in NSET)
def _expand(cell):
    out = []
    for a, b in re.findall(r'N-(\d{3})(?: a N-(\d{3}))?', cell):
        out += [f'N-{n:03d}' for n in _nums if int(a) <= n <= int(b)] if b else [f'N-{a}']
    return out
pant = {}  # pantalla de U1 §3.29 -> necesidades
for l in u1[_i:_j].split('\n'):
    m_ = re.match(r'^\| (Global \(G-01 a G-07\)|P-[\d.]+) \| (.*?) \| (.*?) \| (.*?) \|$', l)
    if m_:
        scr = 'Global' if m_.group(1).startswith('Global') else m_.group(1)
        pant.setdefault(scr, set()).update(_expand(m_.group(2)) + _expand(m_.group(3)))
pares_u1 = {(s, n) for s, ns in pant.items() for n in ns}
s10 = c1[c1.index('## 10. Matriz'):c1.index('## 11. Cobertura')]
filas = re.findall(r'^\| (P-\d{2}(?:\.\d{2})?|Global) \| (N-\d{3}) \| ([LE]) \| (.*?) \| (.*?) \| (.*?) \|$', s10, re.M)
pares_c1 = {(f[0], f[1]) for f in filas}
m10 = {}
for f in filas: m10.setdefault(f[1], set()).update(f'{a} {b}' for a, b in re.findall(r'\b(GET|POST|PUT|PATCH|DELETE) `([^`]+)`', f[4]))
xnm = {}
for m, p, o in OPS:
    for n in o.get('x-necesidades', []): xnm.setdefault(n, set()).add(f'{m} {p}')
opset = {f'{m} {p}' for m, p, _ in OPS}
usados = set().union(*m10.values()) if m10 else set()
xn = set(xnm)
sin_op = sorted(n for n in NSET if not m10.get(n))
sin_op_ni_local = sorted(set(sin_op) - LOCAL_U1)
bloq = sorted({f[1] for f in filas if f[5].startswith('BLOQUEO')})
# inventario de pantallas de U1 (§2.1 P-xx y §2.3 P-26.xx); P-26 se cubre con sus secciones P-26.xx (criterio C1-02)
P_INV = re.findall(r'^\| (P-\d{2}) \|', u1[u1.index('### 2.1'):u1.index('### 2.2')], re.M) + re.findall(r'^\| (P-26\.\d{2}) \|', u1[u1.index('### 2.3'):u1.index('### 2.4')], re.M)
def _needs_of(s):
    if s == 'P-26' and s not in pant: return set().union(*(ns for k, ns in pant.items() if k.startswith('P-26.')))
    return pant.get(s, set())
SCR = ['Global'] + P_INV
p_local = sorted(s for s in SCR if _needs_of(s) and _needs_of(s) <= LOCAL_U1)
p_sin_op = sorted(s for s in SCR if not any(m10.get(n) for n in _needs_of(s)) and s not in p_local)
p_fuera = sorted(set(pant) - set(SCR))
sin_n = sorted(f'{m} {p}' for m, p, o in OPS if not o.get('x-necesidades'))
inexist = sorted(usados - opset)
d1 = sorted((n, k) for n in m10 for k in m10[n] if k not in xnm.get(n, set()))
d2 = sorted((n, k) for n in xnm for k in xnm[n] if k not in m10.get(n, set()))
print(f'C1-02 N-xxx en U1: {len(NSET)} · con operación en la matriz §10: {len(NSET) - len(sin_op)} · Local según U1: {sorted(LOCAL_U1)} · N-xxx sin operación ni Local: {len(sin_op_ni_local)} {sin_op_ni_local} (con estado BLOQUEO en §10: {bloq})')
print(f'C1-02 pantallas de U1 (§2.1 y §2.3, más Global): {len(SCR)} · en la matriz U1 §3.29: {len(pant)} · P-26 por sus secciones P-26.xx · Local: {p_local} · P-xx sin operación ni Local: {len(p_sin_op)} {p_sin_op} · filas de §3.29 que no son pantallas del inventario: {p_fuera}')
print(f'C1-02 operaciones: {len(opset)} · operaciones sin N-xxx: {len(sin_n)} {sin_n} · N-xxx citadas que no existen en U1: {sorted(xn - NSET)} · rutas de la matriz que no existen en openapi.yaml: {inexist}')
print(f'C1-02 matriz §10: {len(filas)} filas · pares pantalla-necesidad de U1 §3.29: {len(pares_u1)} · solo en U1: {sorted(pares_u1 - pares_c1)[:5]} · solo en C1: {sorted(pares_c1 - pares_u1)[:5]} · pares necesidad-operación en §10 sin x-necesidades: {len(d1)} {d1[:5]} · en x-necesidades sin §10: {len(d2)} {d2[:5]}')
ok(not sin_op_ni_local, 'C1-02 0 N-xxx sin operación ni Local (una N-xxx bloqueada no cuenta como cubierta)')
ok(not p_sin_op and not p_fuera, 'C1-02 0 pantallas P-xx sin operación ni Local')
ok(sin_n == ['GET /openapi.yaml'], 'C1-02 única operación sin N-xxx = GET /openapi.yaml (excepción «estáticos», C1 §11.2)')
ok(not (xn - NSET) and not inexist, 'C1-02 0 N-xxx citadas que no existen en U1 y 0 rutas de la matriz que no existen en openapi.yaml')
ok(pares_u1 == pares_c1 and len(filas) == len(pares_c1), 'C1-02 matriz §10 = pares pantalla-necesidad de U1 §3.29 (una fila por par)')
ok(not d1 and not d2, 'C1-02 matriz §10 = x-necesidades de openapi.yaml')
# C1-03
if ARGS.db:
    db = sqlite3.connect(f'file:{os.path.abspath(ARGS.db)}?mode=ro', uri=True)
    print(f'C1-03 base: {os.path.abspath(ARGS.db)} (migrada por fuera, abierta solo lectura; se excluye goose_db_version)')
else:
    import glob
    db = sqlite3.connect(':memory:')
    migs = sorted(glob.glob(os.path.join(ROOT, 'backend/db/migrations/[0-9]*.sql')))
    for mp in migs:
        sql = open(mp, encoding='utf-8').read()
        db.executescript(sql.split('-- +goose Down')[0].split('-- +goose Up', 1)[1].replace('-- +goose StatementBegin', '').replace('-- +goose StatementEnd', ''))
    print('C1-03 base: secciones Up de', ', '.join(os.path.basename(m) for m in migs), 'aplicadas en orden en memoria')
tabs = {r[0] for r in db.execute("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name <> 'goose_db_version'")}
used = {t for _, _, o in OPS for t in o.get('x-tablas', [])}
sinx = [o['operationId'] for _, _, o in OPS if 'x-tablas' not in o]
vac = [o['operationId'] for _, _, o in OPS if o.get('x-tablas') == []]
print(f'C1-03 tablas en sqlite_master: {len(tabs)} · citadas por x-tablas: {len(used)} · inexistentes: {sorted(used - tabs)} · tablas sin operación: {sorted(tabs - used)} · operaciones con x-tablas vacía (archivo o proceso, sin base): {vac}')
ok(not (used - tabs) and not sinx and tabs <= used, 'C1-03 0 tablas inexistentes, toda operación declara x-tablas y toda tabla tiene operación')
ent_props = bad_props = 0; badlist = []
cols = {t: {r[1] for r in db.execute(f'PRAGMA table_info({t})')} for t in tabs}
for n, s in d['components']['schemas'].items():
    t = s.get('x-tabla')
    if not t: continue
    for k, v in s.get('properties', {}).items():
        ent_props += 1
        if k not in cols[t] and not v.get('x-calculado') and not v.get('x-entrada') and not v.get('x-columna') and not (isinstance(v, dict) and v.get('type') == 'array' and v.get('items', {}).get('$ref')) and not v.get('$ref', '').endswith(('ExpenseCreate',)):
            bad_props += 1; badlist.append(f'{n}.{k}')
print(f'C1-03 propiedades de esquemas de entidad: {ent_props} · sin columna ni marca de calculado/entrada: {bad_props} {badlist[:20]}')
# Corrección r2: la cuenta de propiedades ahora sí falla, y cada x-columna debe nombrar una columna que existe.
xcol_ok = 0; xcol_bad = []
for n, s in d['components']['schemas'].items():
    for k, v in (s.get('properties') or {}).items():
        xc = v.get('x-columna') if isinstance(v, dict) else None
        if not xc: continue
        t2, _, c2 = xc.rpartition('.'); t2 = t2 or s.get('x-tabla')
        if not t2: continue  # sin tabla declarada (p. ej. MealPlanEntryCreate): la columna la nombra su descripción
        if c2 in cols.get(t2, set()): xcol_ok += 1
        else: xcol_bad.append(f'{n}.{k} -> {xc}')
print(f'C1-03 x-columna comprobadas contra la base: {xcol_ok} · que no existen: {xcol_bad}')
ok(bad_props == 0 and not xcol_bad, 'C1-03 toda propiedad de esquema de entidad es columna de su tabla o campo calculado o de entrada, y toda x-columna existe')
# C1-04
nerr = nref = 0; malos = []
for m, p, o in OPS:
    for code, r in o['responses'].items():
        if code == 'default' or code[0] in '45':
            nerr += 1
            if '$ref' in r and r['$ref'].startswith('#/components/responses/') and deref(r)['content']['application/json']['schema'] == {'$ref': '#/components/schemas/Error'}: nref += 1
            else: malos.append(f'{m} {p} {code}')
otros = [n for n in d['components']['schemas'] if 'error' in n.lower() and n not in ('Error', 'FieldError')]
print(f'C1-04 respuestas de error: {nerr} · con referencia al esquema único Error: {nref} · excepciones: {malos[:5]} · esquemas de error alternativos: {otros}')
ok(nerr == nref and not otros, 'C1-04 100 % de respuestas 4xx/5xx/default con components/responses → Error')
# C1-05
listas = []; dif = []
for m, p, o in OPS:
    sch = o['responses'].get('200', {}).get('content', {}).get('application/json', {}).get('schema', {})
    if m == 'GET' and sch.get('$ref', '').endswith('Page'):
        listas.append(o['operationId'])
        names = {x['$ref'].split('/')[-1] for x in o.get('parameters', []) if '$ref' in x}
        pg = deref(sch)
        if not {'Cursor', 'Limit', 'Sort'} <= names or pg['allOf'][0] != {'$ref': '#/components/schemas/Page'}: dif.append(o['operationId'])
    elif m == 'GET' and (sch.get('type') == 'array'): dif.append(o['operationId'] + ' (arreglo sin envoltura)')
print(f'C1-05 colecciones GET: {len(listas)} · con cursor+limit+sort y envoltura Page: {len(listas) - len(dif)} · diferencias: {dif}')
ok(not dif, f'C1-05 {len(listas)} de {len(listas)} con la misma paginación y envoltura')
# C1-06
F6 = {'person_id': 'PersonId', 'family': 'Family', 'category_id': 'CategoryId', 'store_id': 'StoreId', 'period': 'Period', 'credit_id': 'CreditId'}
SIN = {'person', 'persona', 'familia', 'family_id', 'category', 'categoria', 'store', 'tienda', 'credit', 'credito', 'periodo', 'month', 'week', 'from', 'to', 'date_from', 'date_to'}
inl = []; sinon = []; usos = {k: 0 for k in F6}
for m, p, o in OPS:
    for x in o.get('parameters', []):
        if '$ref' in x:
            nm = x['$ref'].split('/')[-1]
            for k, v in F6.items():
                if v == nm: usos[k] += 1
            continue
        if x['in'] == 'query' and x['name'] in F6: inl.append(f"{o['operationId']}.{x['name']}")
        if x['in'] == 'query' and x['name'] in SIN: sinon.append(f"{o['operationId']}.{x['name']}")
defs = {k: d['components']['parameters'][v]['name'] == k for k, v in F6.items()}
print('C1-06 usos por filtro:', usos, '· definidos inline:', inl, '· sinónimos:', sinon)
ok(all(defs.values()) and not inl and not sinon, f'C1-06 {sum(defs.values())} de 6 filtros con definición única por $ref, 0 sinónimos')
# H-r2-01 (revisión C1-r2): cada nombre de parámetro tiene un solo componente; ningún otro componente usa el nombre de un filtro;
# los componentes con nombre de periodo que no son el filtro (llaves) se documentan en §6 de C1-datos.md.
PC = d['components']['parameters']
cnt = {}
for k, v in PC.items(): cnt.setdefault((v['name'], v['in']), []).append(k)
dup = sorted(f'{n} ({i}): {ks}' for (n, i), ks in cnt.items() if len(ks) > 1)
otros = sorted(f"{k} → {v['name']}" for k, v in PC.items() if v['name'] in F6 and F6[v['name']] != k)
llaves = sorted((k, v['name']) for k, v in PC.items() if v['name'] in SIN)
s6 = c1[c1.index('## 6. '):c1.index('## 7. ')]
nodoc = [f'{k} → {n}' for k, n in llaves if f'`{n}`' not in s6 or k not in s6]
print('C1-06 nombres con más de un componente:', dup, '· otros componentes con nombre de filtro:', otros, '· llaves que no son filtro:', [f'{k} → {n}' for k, n in llaves], '· sin documentar en §6:', nodoc)
ok(not dup and not otros and not nodoc, f'C1-06 {len(PC)} componentes de parámetro, cada nombre con un solo componente; {len(llaves)} llaves documentadas en §6')
# H-r3-01 (revisión C1-r3): toda operación que lee o escribe la lista de compras (x-tablas con shopping_list o shopping_list_item)
# recibe la llave `week` por $ref a ListWeek, obligatoria, y no recibe el filtro `period`.
lw = PC.get('ListWeek', {})
sl_ops = [(o['operationId'], {x['$ref'].split('/')[-1] for x in o.get('parameters', []) if '$ref' in x}) for m, p, o in OPS if {'shopping_list', 'shopping_list_item'} & set(o.get('x-tablas', []))]
sin_week = sorted(oid for oid, rs in sl_ops if 'ListWeek' not in rs or 'Period' in rs)
print(f"C1-06 operaciones que leen o escriben shopping_list: {len(sl_ops)} · sin `week` o con `period`: {sin_week} · ListWeek = week obligatoria: {lw.get('name') == 'week' and lw.get('required') is True}")
ok(bool(sl_ops) and not sin_week and lw.get('name') == 'week' and lw.get('required') is True, f'C1-06 las {len(sl_ops)} operaciones que leen shopping_list exigen la llave `week` (ListWeek)')
# C1-07
fmts = set(); exs = []
def walk(o, path=''):
    if isinstance(o, dict):
        if 'format' in o and isinstance(o['format'], str): fmts.add(o['format'])
        for key in ('examples', 'example', 'default'):
            if key in o: exs.append((path + '/' + key, o[key]))
        for k, v in o.items(): walk(v, path + '/' + str(k))
    elif isinstance(o, list):
        for i, v in enumerate(o): walk(v, f'{path}[{i}]')
walk(d)
pantalla = re.compile(r'\b\d{1,2}-(ene|feb|mar|abr|may|jun|jul|ago|sep|oct|nov|dic)-\d{4}\b')
fechas = [(p, v) for p, v in exs for x in (v if isinstance(v, list) else [v]) if isinstance(x, str) and (re.search(r'\d{4}-\d{2}-\d{2}', x) or pantalla.search(x))]
mal = [(p, v) for p, v in fechas for x in (v if isinstance(v, list) else [v]) if isinstance(x, str) and (pantalla.search(x) or not re.match(r'^[0-9]{4}-[0-9]{2}-[0-9]{2}', x))]
print(f'C1-07 formatos usados: {sorted(fmts)} · ejemplos con fecha: {len(fechas)} · fuera de ISO: {mal}')
ok(fmts <= {'date', 'date-time'} and not mal, 'C1-07 solo date y date-time; 0 ejemplos de fecha fuera de ISO')
ok('el «30» es el **último día**' in c1 and '2027-02-28' in c1, 'C1-07 C1-datos.md declara las quincenas 15 y 30 y que en febrero el «30» es el 28')
# C1-08
money = nomoney = 0; badm = []
MW = re.compile(r'(amount|total|principal|interest|price|balance|available|payment|fee|_cents)')
def isref(v):
    return v.get('$ref') == '#/components/schemas/MoneyCents' or any(isinstance(a, dict) and a.get('$ref') == '#/components/schemas/MoneyCents' for a in v.get('anyOf', []))
def walkp(o, path=''):
    global money, nomoney
    if isinstance(o, dict):
        for k, v in o.get('properties', {}).items():
            if isinstance(v, dict):
                if k.endswith('_cents'):
                    money += 1
                    if not isref(v): badm.append(path + '.' + k)
                elif MW.search(k) and (v.get('type') == 'number' or isref(v)):
                    nomoney += 1; badm.append(path + '.' + k + ' (sin sufijo _cents)')
        for k, v in o.items(): walkp(v, path + '/' + str(k))
    elif isinstance(o, list):
        for i, v in enumerate(o): walkp(v, path)
walkp(d['components']['schemas'], 'schemas'); walkp(d['paths'], 'paths')
nnum = open(OA, encoding='utf-8').read().count('type: "number"') + open(OA, encoding='utf-8').read().count('type: number')
print(f'C1-08 propiedades de dinero (_cents): {money} · sin referencia a MoneyCents: {badm} · type: number en el archivo: {nnum}')
ok(not badm and nnum == 0, 'C1-08 100 % del dinero por $ref a MoneyCents; 0 number')
# C1-12
pub = sorted(f"{m} {p}" for m, p, o in OPS if o.get('security') == [])
seen = set(); hits = []
def walks(s, path):
    s = deref(s) if isinstance(s, dict) and '$ref' in s else s
    if not isinstance(s, dict) or id(s) in seen: return
    seen.add(id(s))
    for k, v in s.get('properties', {}).items():
        # C1 v1.2 (K44): una clave de app_setting expuesta como parámetro (x-app-setting = su nombre, p. ej. auth.password_min_length)
        # es configuración (un entero), no un secreto.
        if re.search(r'password|hash|token|secret', k) and not (isinstance(v, dict) and v.get('x-app-setting') == k): hits.append(path + '.' + k)
        walks(v, path + '.' + k)
    for key in ('items', 'allOf', 'anyOf', 'oneOf'):
        v = s.get(key)
        for x in (v if isinstance(v, list) else [v] if v else []): walks(x, path)
for m, p, o in OPS:
    for code, r in o['responses'].items():
        for ct in deref(r).get('content', {}).values(): walks(ct.get('schema', {}), f'{m} {p} {code}')
sch = d['components']['securitySchemes']['cookieAuth']
cuentas = [o['operationId'] for m, p, o in OPS if re.search(r'(create|post).*(user|account)', o['operationId'], re.I) or (m == 'POST' and p.rstrip('/').endswith(('/users', '/persons')))]
print(f'C1-12 operaciones públicas: {pub} · campos sensibles en respuestas: {hits} · cookie: {sch["in"]} {sch["name"]} · operaciones de alta de cuentas o personas: {cuentas}')
ok(pub == ['GET /health', 'POST /api/v1/auth/login'] and not hits and not cuentas and 'HttpOnly' in sch['description'] and 'SameSite' in sch['description'], 'C1-12 solo login y salud públicos; 0 campos sensibles; sin alta de cuentas')
ok('5 intentos fallidos en 15 minutos' in d['paths']['/api/v1/auth/login']['post']['description'] and '429' in d['paths']['/api/v1/auth/login']['post']['responses'], 'C1-12 límite 5 en 15 min descrito en login con 429')
# C1-13
ids = [o['operationId'] for _, _, o in OPS]
TAGS9 = ['Inicio', 'Gastos', 'Presupuesto', 'Comidas', 'Compras y tiendas', 'Créditos', 'Salud', 'Reportes', 'Ajustes']
badtag = [o['operationId'] for _, _, o in OPS if not o.get('tags') or not set(o['tags']) <= set(TAGS9)]
print(f'C1-13 operationId duplicadas: {len(ids) - len(set(ids))} · sin etiqueta válida: {badtag} · etiquetas: {[t["name"] for t in d["tags"]]} · info.version: {d["info"]["version"]}')
ok(len(ids) == len(set(ids)) and not badtag and [t['name'] for t in d['tags']] == TAGS9 and 'C1 v1.3' in c1 and d['info']['version'] == '1.3.0', 'C1-13 ids únicas, 9 etiquetas de la barra lateral, C1 v1.3')
# K44 (C1 v1.2, puerta 1b): cada parámetro @param de las migraciones tiene una sola lectura (GET con parameters) y un solo cambio con
# «restablecer» (PATCH con parameters y reset_parameters); la tabla de §11.7 da su valor inicial y la base sembró ese valor.
PRM = []
MD = os.path.join(ROOT, 'backend/db/migrations')
for f in sorted(os.listdir(MD)):
    if re.match(r'\d+.*\.sql$', f):
        PRM += [[x.strip() for x in l.split('|')] for l in re.findall(r'^-- @param (.*)$', open(os.path.join(MD, f), encoding='utf-8').read().split('-- +goose Down')[0], re.M)]
SCH = d['components']['schemas']
def _pk(ref_, campo):
    s = SCH.get((ref_ or '').split('/')[-1], {}).get('properties', {}).get(campo, {})
    return set(s.get('properties', {})) if campo == 'parameters' else set(s.get('items', {}).get('enum', []))
lee, cambia = {}, {}
for m, p, o in OPS:
    if not isinstance(o, dict): continue
    rs = o.get('responses', {}).get('200', {}).get('content', {}).get('application/json', {}).get('schema', {}).get('$ref')
    rq = o.get('requestBody', {}).get('content', {}).get('application/json', {}).get('schema', {}).get('$ref')
    if m == 'GET':
        for k in _pk(rs, 'parameters'): lee.setdefault(k, []).append(o['operationId'])
    if m == 'PATCH':
        for k in _pk(rq, 'parameters') & _pk(rq, 'reset_parameters'): cambia.setdefault(k, []).append(o['operationId'])
s117 = c1.split('### 11.7', 1)[1].split('\n## ', 1)[0] if '### 11.7' in c1 else ''
fila = {r[0]: next((l for l in s117.split('\n') if l.startswith(f'| `{r[0]}` |')), '') for r in PRM}
pmal = [r[0] for r in PRM if len(lee.get(r[0], [])) != 1 or len(cambia.get(r[0], [])) != 1 or f'| `{r[3]}` |' not in fila[r[0]]]
semb = {}
for r in PRM:
    v = db.execute('SELECT value_json FROM app_setting WHERE setting_key = ? AND user_account_id IS NULL', (r[0],)).fetchone()
    semb[r[0]] = v is not None and json.loads(v[0]) == json.loads(r[3])
sobra = sorted(set(lee) - {r[0] for r in PRM})
print(f'K44 parámetros @param: {len(PRM)} · con una lectura y un cambio con restablecer, y en §11.7 con su valor inicial: {len(PRM) - len(pmal)} · fallan: {pmal} · sembrados con ese valor: {sum(semb.values())} · expuestos sin @param: {sobra}')
ok(PRM and not pmal and all(semb.values()) and not sobra, 'K44 todo parámetro operativo se lee y se cambia o restablece por clave, con su valor inicial sembrado')
# C1-14
try:
    from jsonschema import Draft202012Validator
    from referencing import Registry, Resource
    reg = Registry().with_resource('urn:c1', Resource.from_contents(d, default_specification=__import__('referencing.jsonschema', fromlist=['DRAFT202012']).DRAFT202012))
    nex = bad = 0; badex = []
    def walke(o, ptr):
        global nex, bad
        if isinstance(o, dict):
            if isinstance(o.get('examples'), list) and ('type' in o or '$ref' in o or 'properties' in o or 'pattern' in o):
                for ex in o['examples']:
                    nex += 1
                    errs = list(Draft202012Validator({'$ref': 'urn:c1#' + ptr}, registry=reg).iter_errors(ex))
                    if errs: bad += 1; badex.append((ptr, ex, errs[0].message))
            for k, v in o.items(): walke(v, ptr + '/' + str(k).replace('~', '~0').replace('/', '~1'))
        elif isinstance(o, list):
            for i, v in enumerate(o): walke(v, f'{ptr}/{i}')
    walke(d['components'], '/components')
    print(f'C1-14 ejemplos validados con jsonschema: {nex} · inválidos: {bad} {badex[:3]}')
    ok(bad == 0, 'C1-14 cada example valida contra su esquema')
except ImportError as e:
    ok(False, f'C1-14 jsonschema no disponible: {e}')
EN = re.compile(r'\b(the|and|with|of|for|is|are|this|that|from)\b')
descs = [o for o in [] ]
allds = []
def walkd(o):
    if isinstance(o, dict):
        if isinstance(o.get('description'), str): allds.append(o['description'])
        if isinstance(o.get('summary'), str): allds.append(o['summary'])
        for v in o.values(): walkd(v)
    elif isinstance(o, list):
        for v in o: walkd(v)
walkd(d)
eng = [x for x in allds if EN.search(re.sub(r'`[^`]*`', '', x))]
random.seed(20261010); muestra = random.sample(allds, 20)
print(f'C1-14 descripciones y resúmenes: {len(allds)} · con palabras frecuentes del inglés fuera de código: {len(eng)} {eng[:3]}')
print('C1-14 muestra de 20 descripciones para el revisor:'); [print('   -', x[:110]) for x in muestra]
ok(not eng, 'C1-14 descripciones en español')
# C1-14 (P2c, corrección C1 r1, H-1): nombre de propiedad = nombre_codigo. O() usa el «_» final como marca de opcional y lo quita; si una
# propiedad publicada termina en «_» (p. ej. por _ins, que no lo quita), el nombre no es el de C1-datos. Se revisa openapi.yaml y cada extracto.
# Además, cada nombre entre comillas invertidas de «C1 v1.3 (puerta 2)» de C1-datos.md debe existir en openapi.yaml con su tipo de nombre.
_PROPN, _ENUMV, _XTAB, _XSET = set(), set(), set(), set()
def _walkn(o, props):
    if isinstance(o, dict):
        if isinstance(o.get('properties'), dict): props.update(k for k in o['properties'] if isinstance(k, str))
        if isinstance(o.get('enum'), list): _ENUMV.update(v for v in o['enum'] if isinstance(v, str))
        if isinstance(o.get('x-tabla'), str): _XTAB.add(o['x-tabla'])
        if isinstance(o.get('x-tablas'), list): _XTAB.update(o['x-tablas'])
        if isinstance(o.get('x-app-setting'), str): _XSET.add(o['x-app-setting'])
        for v in o.values(): _walkn(v, props)
    elif isinstance(o, list):
        for v in o: _walkn(v, props)
_walkn(d, _PROPN)
_guion = sorted(n for n in _PROPN if n.endswith('_'))
_EXD14 = os.path.join(ROOT, 'docs/contratos/extractos'); _guion_ex = []
for _f in sorted(os.listdir(_EXD14)) if os.path.isdir(_EXD14) else []:
    if _f.endswith('.yaml'):
        _pe = set(); _walkn(yaml.safe_load(open(os.path.join(_EXD14, _f), encoding='utf-8')), _pe)
        _guion_ex += [f'{_f}:{n}' for n in sorted(_pe) if n.endswith('_')]
print(f'C1-14 propiedades distintas: {len(_PROPN)} · que terminan en «_»: openapi {_guion[:5]} · extractos {_guion_ex[:5]}')
ok(not _guion and not _guion_ex, 'C1-14 0 propiedades que terminan en «_» (openapi.yaml y extractos)')
_SCH = d['components']['schemas']; _OPIDS = {o['operationId'] for _, _, o in OPS}
def _prop14(s, n):
    s = deref(s)
    for alt in [s] + [deref(x) for k in ('allOf', 'anyOf', 'oneOf') for x in s.get(k, [])]:
        if isinstance(alt, dict) and n in alt.get('properties', {}): return alt['properties'][n]
    return None
def _ruta14(s, partes):
    for p in partes:
        arr = p.endswith('[]'); s = _prop14(s, p[:-2] if arr else p)
        if s is None: return False
        if arr:
            s = deref(s); its = [s.get('items')] + [deref(x).get('items') for x in s.get('anyOf', []) if isinstance(deref(x), dict)]
            s = next((x for x in its if x is not None), None)
            if s is None: return False
    return True
_KW14 = {'minimum', 'maximum', 'pattern', 'true', 'false'}   # palabras de JSON Schema o de JSON, no campos
_RET14 = {'send_time', 'advance_days'}   # la sección dice que params_json ya NO los declara (lo comprueba C1-23)
_PJ14 = json.dumps(_prop14(_SCH['Notification'], 'params_json') or {}, ensure_ascii=False)   # claves propias de params_json: se nombran en su descripción
_sec14 = c1.split('### C1 v1.3 (puerta 2)', 1)[1].split('\n### ', 1)[0] if '### C1 v1.3 (puerta 2)' in c1 else ''
_tok14 = sorted(set(re.findall(r'`([^`\s]+)`', _sec14))); _falt14, _otros14 = [], []
for t in _tok14:
    if re.fullmatch(r'[A-Z]+', t): _otros14.append(t); continue   # PATCH, UPDATE: verbos HTTP o SQL
    if t.startswith('/api/'): good = t in d['paths']
    elif re.fullmatch(r'[A-Z][A-Za-z0-9]*', t): good = t in _SCH
    elif re.fullmatch(r'[a-z]+[A-Z][A-Za-z0-9]*', t): good = t in _OPIDS
    elif re.fullmatch(r'[A-Z][A-Za-z0-9]*(\.[a-z_][a-z0-9_]*(\[\])?)+', t):
        _s, *_ps = t.split('.'); good = _s in _SCH and _ruta14(_SCH[_s], _ps)
    elif '[]' in t and re.fullmatch(r'[a-z_][a-z0-9_]*(\[\])?(\.[a-z_][a-z0-9_]*(\[\])?)*', t): good = any(_ruta14(s, t.split('.')) for s in _SCH.values())
    elif re.fullmatch(r'[a-z]+\.[a-z0-9_]+', t) and not re.search(r'\.(py|sql|md|yaml|go)$', t): good = t in _XSET or t in _PROPN
    elif re.fullmatch(r'[a-z][a-z0-9_]*', t): good = t in _PROPN or t in _ENUMV or t in _XTAB or t in _KW14 or t in _RET14 or f'`{t}`' in _PJ14
    else: _otros14.append(t); continue
    if not good: _falt14.append(t)
print(f'C1-14 nombres en «C1 v1.3 (puerta 2)» de C1-datos.md: {len(_tok14)} · comprobados: {len(_tok14) - len(_otros14)} · que no existen en openapi.yaml: {_falt14} · no son nombres de campo (archivos, anotaciones, formatos): {_otros14}')
ok(bool(_sec14) and not _falt14, 'C1-14 cada esquema, operación, propiedad (con su ruta), clave de app_setting y valor que nombra «C1 v1.3 (puerta 2)» de C1-datos.md existe en openapi.yaml')
# C1-09 (P2c, corrección C1 r1, H-2): M-11 nombra los acumulados de Q2-65 A, y las cifras de su ejemplo (oct-2026, ventana desde nov-2026
# y mar-2027) se recalculan aquí desde oraculo/model.json (Σ periodos[i] con Decimal), no se copian a mano.
from decimal import Decimal as _D9
_mj9 = json.load(open(os.path.join(ROOT, 'oraculo/model.json'), encoding='utf-8'), parse_float=_D9)
_ca9 = _cc9 = _D9(0); _ac9 = []
for _p9 in _mj9['periodos']:
    _ca9 += _D9(str(_p9['devengado']['total'])); _cc9 += _D9(str(_p9['efectivo']['total']))
    _ac9.append((_ca9, _cc9, _ca9 - _cc9, _D9(str(_p9['diferencia_acumulada']))))
_m11 = [l for l in c1.splitlines() if l.startswith('| M-11 |')]
_esp9 = [f'{x:,.2f}' for i in (0, 1, len(_ac9) - 1) for x in _ac9[i][:3]]
_m11c = [c for c in ('cumulative_accrued_cents', 'cumulative_cash_cents', 'cumulative_difference_cents') if len(_m11) != 1 or f'`{c}`' not in _m11[0]]
_m11x = [x for x in _esp9 if len(_m11) != 1 or x not in _m11[0]]
_dif9 = [i for i, a in enumerate(_ac9) if a[2] != a[3]]
print(f'C1-09 M-11: filas {len(_m11)} · cifras del oráculo (P1, P2 y P{len(_ac9)}: devengado acumulado, efectivo acumulado, diferencia) {_esp9} · faltan en la fila: {_m11x} · campos que no nombra: {_m11c} · periodos donde Σ devengado − Σ efectivo ≠ diferencia_acumulada: {_dif9}')
ok(len(_m11) == 1 and not _m11c and not _m11x and not _dif9, 'C1-09 M-11 nombra cumulative_accrued_cents, cumulative_cash_cents y cumulative_difference_cents, y su ejemplo coincide con oraculo/model.json')
# C1-15
B = os.path.join(ROOT, 'docs/estado/bloqueos'); bl = {f[:-3] for f in os.listdir(B) if f.endswith('.md')}
# C1 v1.1: una marca solo es válida si cita un bloqueo ABIERTO (sin «Respuesta de Axel» ni «Estado: RESUELTO»)
abiertos = {b for b in bl if not re.search(r'Respuesta de Axel|Estado: RESUELTO', open(os.path.join(B, b + '.md'), encoding='utf-8').read())}
nprov = 0; huerf = []
def walkv(o, path):
    global nprov
    if isinstance(o, dict):
        de = o.get('description', '') if isinstance(o.get('description'), str) else ''
        if o.get('x-provisional') or 'PROVISIONAL' in de:
            nprov += 1
            if not any(b in de for b in abiertos): huerf.append(path)
        for k, v in o.items(): walkv(v, path + '/' + str(k))
    elif isinstance(o, list):
        for i, v in enumerate(o): walkv(v, f'{path}[{i}]')
walkv(d, '')
lin = [l for l in c1.split('\n') if 'PROVISIONAL' in l]
c1h = [l[:80] for l in lin if not any(b in l for b in abiertos) and not l.startswith('#') and not l.startswith('| C1 v') and not re.match(r'^\| (Tema|#) \|', l) and 'Estado:' not in l and 'Todo lo que no sale' not in l and 'C1 tomó en cada caso' not in l]  # títulos y encabezados de tabla no afirman nada
refs = sorted({r for l in c1.split('\n') for r in re.findall(r'`((?:C1|S1|S6|U1|criterios)-[a-z0-9-]+)`', l)})
print(f'C1-15 marcas PROVISIONAL en openapi.yaml: {nprov} · bloqueos abiertos citados: {sorted(b for b in abiertos if b in str(d))} · sin bloqueo abierto: {huerf[:5]} · líneas PROVISIONAL en C1-datos.md: {len(lin)} · sin bloqueo: {c1h[:5]} · bloqueos citados en C1: {len(refs)} · citados que no existen: {[r for r in refs if r not in bl]}')
ok(not huerf and not c1h and all(r in bl for r in refs), 'C1-15 toda marca PROVISIONAL remite a un bloqueo abierto (sin «Respuesta de Axel»)')
# H-r3-02 (revisión C1-r3): mientras el bloqueo U1-lista-vigente-domingo siga abierto, toda descripción que aplique la regla del domingo
# («los domingos»; «lunes a domingo» de K8 y «7 = domingo» no son la regla) lleva x-provisional en su objeto o en su operación y cita el
# bloqueo; en C1-datos.md, la línea cita el bloqueo.
DOM = 'U1-lista-vigente-domingo'; dsin = []
def walkdom(o, path, prov=False):
    if isinstance(o, dict):
        prov = prov or bool(o.get('x-provisional'))
        de = o['description'] if isinstance(o.get('description'), str) else ''
        if re.search(r'(?i)\bdomingos\b', de) and not (prov and 'PROVISIONAL' in de and DOM in de): dsin.append(path)
        for k, v in o.items(): walkdom(v, path + '/' + str(k), prov)
    elif isinstance(o, list):
        for i, v in enumerate(o): walkdom(v, f'{path}[{i}]', prov)
if DOM in abiertos: walkdom(d, '')
cdom = [l[:80] for l in c1.split('\n') if DOM in abiertos and re.search(r'(?i)\bdomingos\b', l) and DOM not in l]
print(f'C1-15 regla del domingo (FIRME por K38; bloqueo de origen abierto: {DOM in abiertos}) · en openapi sin marca: {dsin[:5]} · en C1-datos.md sin bloqueo: {cdom[:5]}')
ok(not dsin and not cdom, 'C1-15 la regla del domingo («los domingos») es FIRME por K38 y no deja marcas pendientes')
# EXT extractos (protocolo-agentes §4, K32): cada operación en un solo extracto, ≤ 80 KB y referencias resueltas
EXD = os.path.join(ROOT, 'docs/contratos/extractos')
exf = sorted(f for f in os.listdir(EXD) if f.endswith('.yaml')) if os.path.isdir(EXD) else []
exdocs = {f: yaml.safe_load(open(os.path.join(EXD, f), encoding='utf-8')) for f in exf}
exbig = [f for f in exf if os.path.getsize(os.path.join(EXD, f)) > 80 * 1024]
excnt = {}
for e in exdocs.values():
    for it in (e.get('paths') or {}).values():
        for o_ in it.values():
            excnt[o_['operationId']] = excnt.get(o_['operationId'], 0) + 1
allids = [o['operationId'] for _, _, o in OPS]
exrep = sorted(k for k, n in excnt.items() if n != 1); exmiss = sorted(set(allids) - set(excnt)); exextra = sorted(set(excnt) - set(allids))
def exrefs(o, acc):
    if isinstance(o, dict):
        if isinstance(o.get('$ref'), str): acc.append(o['$ref'])
        for x in o.values(): exrefs(x, acc)
    elif isinstance(o, list):
        for x in o: exrefs(x, acc)
    return acc
exbad = []
for f, e in exdocs.items():
    for r in exrefs(e, []):
        fn, _, frag = r.partition('#'); node = exdocs.get(fn) if fn else e
        for part in frag.strip('/').split('/'):
            node = node.get(part) if isinstance(node, dict) else None
        if node is None: exbad.append(f'{f}: {r}')
print(f'EXT extractos: {len(exf)} archivos (incluye componentes-comunes.yaml) · operaciones en extractos: {sum(excnt.values())} · en openapi.yaml: {len(allids)} · repetidas: {exrep[:3]} · faltantes: {exmiss[:3]} · sobrantes: {exextra[:3]} · > 80 KB: {exbig} · referencias sin resolver: {len(exbad)} {exbad[:3]} · INDICE.md: {os.path.isfile(os.path.join(EXD, "INDICE.md"))}')
ok(bool(exf) and 'componentes-comunes.yaml' in exf and not exrep and not exmiss and not exextra and not exbig and not exbad and os.path.isfile(os.path.join(EXD, 'INDICE.md')), 'EXT extractos: cada operación en un solo extracto, ≤ 80 KB y referencias resueltas')
# C1-04 (P2c, corrección C1 r2, M-1): qué responde la API al CHECK completo de credit.vat_included (00004). Las combinaciones se
# prueban contra la base (una copia en memoria, cada una dentro de un SAVEPOINT que se deshace), no se copian: alta de un
# store_credit con vat_on_interest_bp = 1600 y vat_included true, null, false o ausente; las mismas sin porcentaje como control; y
# Merge Patch sobre un crédito con (false, 1600). Cada combinación que la base rechaza tiene que estar en la regla 422 de CreditCreate,
# CreditUpdate, createCredit y updateCredit (openapi.yaml y su extracto), y la que acepta, como la única que se admite. En la copia
# se apaga la llave a person: aquí solo se prueba el CHECK de credit.
_dbv = sqlite3.connect(':memory:'); db.backup(_dbv); _dbv.isolation_level = None; _dbv.execute('PRAGMA foreign_keys = OFF')
_JV = {1: 'true', None: 'null', 0: 'false'}
_NIVA = 1600
_INSV = "INSERT INTO credit (name, kind, holder_person_id, payment_frequency, vat_on_interest_bp{c}) VALUES ('c1-04 iva', 'store_credit', 1, 'monthly', ?{v})"
def _civa(ins, args, upd=None):
    _dbv.execute('SAVEPOINT iva')
    try:
        _dbv.execute(ins, args)
        if upd: _dbv.execute(upd)
        return 'acepta'
    except sqlite3.IntegrityError as e:
        return 'rechaza' if 'CHECK constraint failed' in str(e) else f'otro: {e}'
    except sqlite3.DatabaseError as e:
        return f'otro: {e}'
    finally:
        _dbv.execute('ROLLBACK TO iva'); _dbv.execute('RELEASE iva')
_alta4 = {_JV[v]: _civa(_INSV.format(c=', vat_included', v=', ?'), (_NIVA, v)) for v in (1, None, 0)}
_alta4['ausente'] = _civa(_INSV.format(c='', v=''), (_NIVA,))
_ctl4 = {_JV[v]: _civa(_INSV.format(c=', vat_included', v=', ?'), (None, v)) for v in (1, None, 0)}
_PAT4 = {'{"vat_included": true}': 'vat_included = 1', '{"vat_included": null}': 'vat_included = NULL',
         '{"vat_included": true, "vat_on_interest_bp": null}': 'vat_included = 1, vat_on_interest_bp = NULL'}
_pat4 = {k: _civa(_INSV.format(c=', vat_included', v=', ?'), (_NIVA, 0), f"UPDATE credit SET {s} WHERE name = 'c1-04 iva'") for k, s in _PAT4.items()}
_ORD4 = ('true', 'null', 'false', 'ausente')
def _lista4(xs): return (', '.join(xs[:-1]) + ' o ' + xs[-1]) if len(xs) > 1 else ''.join(xs)
_rech4 = [k for k in _ORD4 if _alta4[k] == 'rechaza']; _adm4 = [k for k in _ORD4 if _alta4[k] == 'acepta']
_cl4 = [f'vat_on_interest_bp con valor y vat_included = {_lista4(_rech4)}', '422 validation_failed', 'vat_on_interest_bp', 'no cambia el otro campo por su cuenta']
_cls4 = _cl4 + [f'vat_on_interest_bp con valor solo se admite con vat_included = {_lista4(_adm4)}', 'no guarda nada', 'field = vat_on_interest_bp']
_cls4 += ([f'"vat_on_interest_bp": {_NIVA}}} sin vat_included'] if _alta4['ausente'] == 'rechaza' else []) + list(_PAT4)
_txt4 = [('openapi.yaml', 'CreditCreate', d['components']['schemas']['CreditCreate'].get('description', ''), _cls4),
         ('openapi.yaml', 'CreditUpdate', d['components']['schemas']['CreditUpdate'].get('description', ''), _cls4)]
_txt4 += [('openapi.yaml', o['operationId'], o.get('description', ''), _cl4) for _, _, o in OPS if o['operationId'] in ('createCredit', 'updateCredit')]
_sin4 = [o['operationId'] for _, _, o in OPS if o['operationId'] in ('createCredit', 'updateCredit') and '422' not in o.get('responses', {})]
for _f4, _e4 in exdocs.items():
    for _n4 in ('CreditCreate', 'CreditUpdate'):
        if _n4 in ((_e4.get('components') or {}).get('schemas') or {}): _txt4.append((_f4, _n4, _e4['components']['schemas'][_n4].get('description', ''), _cls4))
    for _it4 in (_e4.get('paths') or {}).values():
        for _o4 in _it4.values():
            if _o4['operationId'] in ('createCredit', 'updateCredit'): _txt4.append((_f4, _o4['operationId'], _o4.get('description', ''), _cl4))
_falt4 = [f'{a}:{b}: {c!r}' for a, b, t, cs in _txt4 for c in cs if c not in t]
_otro4 = [r for r in list(_alta4.values()) + list(_ctl4.values()) + list(_pat4.values()) if r not in ('acepta', 'rechaza')]
_pok4 = all((_pat4[k] == 'rechaza') == (k != '{"vat_included": true, "vat_on_interest_bp": null}') for k in _PAT4)
print(f'C1-04 IVA (CHECK de credit.vat_included, base): alta con {_NIVA} {_alta4} · sin porcentaje {_ctl4} · Merge Patch sobre (false, {_NIVA}) {_pat4} · textos revisados: {len(_txt4)} · operaciones sin respuesta 422: {_sin4} · lo que falta en los textos: {_falt4[:6]}{" …" if len(_falt4) > 6 else ""} ({len(_falt4)})')
ok(len(_txt4) == 8 and not _sin4 and not _otro4 and _rech4 and _adm4 and all(v == 'acepta' for v in _ctl4.values()) and _pok4 and not _falt4,
   'C1-04 cada combinación de vat_included y vat_on_interest_bp que el CHECK de la 00004 rechaza (alta y Merge Patch) está en la regla 422 de CreditCreate, CreditUpdate, createCredit y updateCredit')
# S6-03 paridad de rutas
rg = open(os.path.join(ROOT, 'backend/internal/api/routes.go'), encoding='utf-8').read()
rr = set(re.findall(r'Method: "(\w+)", Pattern: "([^"]+)", OperationID: "(\w+)"', rg))
oa = {(m, p, o['operationId']) for m, p, o in OPS}
print(f'S6-03 rutas en routes.go: {len(rr)} · en openapi.yaml: {len(oa)} · sobrantes: {sorted(rr - oa)[:3]} · faltantes: {sorted(oa - rr)[:3]}')
ok(rr == oa, 'S6-03 conjunto de rutas del router idéntico al contrato')
# C1-23 (C1 v1.2 · corrección r1, hallazgo A-01 de docs/estado/p1c/REVISION-C1-r1.md; K43): send_time (hora local HH:MM de 24 horas)
# y advance_days (entero >= 0) son propiedades de params_json en Notification y NotificationUpdate, sin valor por omisión; weekday y
# day_of_month entran por additionalProperties: true y la marca de lugar x-provisional cita S1-parametros-avisos. Se comprueba en
# openapi.yaml, en el extracto que trae esos esquemas y en C1-datos.md, junto con el resto de C1-23 (10 MB y 413, 503, N-xxx, M-05).
# C1 v1.3 (puerta 2, Q2-11 B, K47, K50): send_time y advance_days dejan params_json (y weekday y day_of_month); la hora, el día y la
# anticipación son claves de app_setting con su rango en el cambio de su sección. Se comprueba: params_json sin esas claves ni marca de
# lugar en openapi.yaml y en el extracto que trae NotificationUpdate; las horas con el patrón HH:MM; los rangos aprobados (S1 1.4,
# docs/estado/p2c/S1-puerta2.md) sin valor por omisión, y la columna «Rango (cambio)» de C1-datos.md §11.7.
PAT_K43 = '^([01][0-9]|2[0-3]):[0-5][0-9]$'
RANGOS_P2 = {'reports.daily_send_time': 'p', 'reports.weekly_send_time': 'p', 'reports.monthly_send_time': 'p', 'alerts.send_time': 'p',
             'reports.weekly_weekday': (1, 7), 'alerts.meal_plan_missing_weekday': (1, 7), 'reports.monthly_day_of_month': (1, 31),
             'comparator.timeout_seconds': (1, 30), 'price.stale_after_days': (1, None), 'alerts.stale_prices_advance_days': (0, None),
             'optimizer.date_window_days': (0, None), 'ocr.timeout_seconds': (0, None)}
def _k43(doc, n):
    pj = (((doc.get('components') or {}).get('schemas') or {}).get(n) or {}).get('properties', {}).get('params_json', {}); m = []
    if set(pj.get('properties', {})) & {'send_time', 'advance_days', 'weekday', 'day_of_month'}: m.append('claves')
    if pj.get('x-provisional') or 'S1-parametros-avisos' in pj.get('description', ''): m.append('lugar')
    if pj.get('type') != 'object': m.append('tipo')
    return m
k43bad = {n: v for n in ('Notification', 'NotificationUpdate') for v in [_k43(d, n)] if v}
k43ex = sorted(f for f, e in exdocs.items() if 'NotificationUpdate' in (((e or {}).get('components') or {}).get('schemas') or {}))
k43exbad = {f: v for f in k43ex for v in [_k43(exdocs[f], 'Notification') + _k43(exdocs[f], 'NotificationUpdate')] if v}
_wq = {}
for _m, _p, _o in OPS:
    if _m == 'PATCH' and isinstance(_o, dict):
        _rq = _o.get('requestBody', {}).get('content', {}).get('application/json', {}).get('schema', {}).get('$ref', '') or ''
        _wq.update(SCH.get(_rq.split('/')[-1], {}).get('properties', {}).get('parameters', {}).get('properties', {}))
def _rok(k, r):
    v = _wq.get(k, {})
    if r == 'p': return v.get('type') == 'string' and v.get('pattern') == PAT_K43 and 'default' not in v
    return v.get('type') == 'integer' and v.get('minimum') == r[0] and v.get('maximum') == r[1] and 'default' not in v
k47rng = sorted(k for k, r in RANGOS_P2.items() if not _rok(k, r))
def _rtxt(r): return '`pattern: ' + PAT_K43.replace('|', '\\|') + '`' if r == 'p' else f'`minimum: {r[0]}`' + (f', `maximum: {r[1]}`' if r[1] is not None else '')
k43c1 = [k for k, r in RANGOS_P2.items() if not re.search(r'^\| `%s` \|.*\| %s \|' % (re.escape(k), re.escape(_rtxt(r))), c1, re.M)]
_rx = re.compile(PAT_K43)
k43rx = all(_rx.fullmatch(v) for v in ('00:00', '07:30', '19:05', '23:59')) and not any(_rx.fullmatch(v) for v in ('24:00', '7:30', '12:60', '12:5', '1230', '08:00 '))
_byid = {o['operationId']: o for _, _, o in OPS}
_up, _cmp = _byid.get('uploadExpenseAttachment', {}), _byid.get('compareScenarios', {})
k43ops = '413' in _up.get('responses', {}) and '10 MB' in json.dumps(_up, ensure_ascii=False) and '503' in _cmp.get('responses', {}) and 'unavailable' in json.dumps(deref(_cmp['responses']['503']) if '503' in _cmp.get('responses', {}) else {}, ensure_ascii=False) + json.dumps(_cmp, ensure_ascii=False)
_nrow = {n: [l for l in c1.splitlines() if l.startswith('| P-') and f'| {n} |' in l] for n in ('N-023', 'N-064', 'N-226', 'N-297', 'N-300')}
k43n = sorted(n for n, ls in _nrow.items() if not ls or any(not re.search(r'FIRME \(K43[,)]', l) or 'C1-valores-k31' in l for l in ls))
k43m05 = bool(re.search(r'^\| M-05 \|.*\b21\b', c1, re.M)); k43grep = len(re.findall(r'PROVISIONAL.*C1-valores-k31', c1))
print(f'C1-23 params_json sin send_time, advance_days, weekday ni day_of_month y sin marca de lugar: mal en openapi: {k43bad} · extractos con NotificationUpdate: {k43ex} mal: {k43exbad} · claves de puerta 2 con rango distinto del aprobado: {k47rng} · patrón acepta 00:00–23:59 y rechaza 24:00, 7:30, 12:60: {k43rx} · filas de §11.7 sin su rango: {k43c1}')
print(f'C1-23 10 MB y 413 en uploadExpenseAttachment, 503 unavailable en compareScenarios: {k43ops} · filas §10 sin FIRME (K43) o con C1-valores-k31: {k43n} · M-05 con 21: {k43m05} · grep PROVISIONAL.*C1-valores-k31: {k43grep}')
ok(not k43bad and k43ex and not k43exbad and not k47rng and k43rx and not k43c1, 'C1-23 hora, día y anticipación como claves de app_setting con su rango aprobado (patrón HH:MM, 1 a 7, 1 a 31, 1 a 30, 1 o más) y fuera de params_json (openapi.yaml, extracto y C1-datos.md), sin valor por omisión')
ok(k43ops and not k43n and k43m05 and k43grep == 0, 'C1-23 adjunto 10 MB con 413, comparador 503 unavailable, N-023/N-064/N-226/N-297/N-300 FIRMES (K43) y M-05 con máximo 21')
print('\nRESULTADO:', 'TODO OK' if not fails else f'{len(fails)} falla(s)')
sys.exit(1 if fails else 0)
