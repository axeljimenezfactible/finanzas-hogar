# Verificación de C1 (criterios C1-01 a C1-15, consistencia de la matriz §10 y paridad de rutas S6-03).
# Uso, desde cualquier carpeta (lee el repositorio que contiene este archivo):
#   uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py [--db BASE] [--root RAIZ]
#   --db BASE   base SQLite ya migrada (p. ej. goose -dir backend/db/migrations sqlite3 BASE up) para C1-03, abierta solo lectura;
#               sin --db se aplican en memoria, en orden, las secciones Up de backend/db/migrations/*.sql.
#   --root RAIZ otra copia del repositorio (por defecto, cuatro niveles arriba de esta carpeta).
# Código de salida: 0 si todo OK; 1 si hay alguna FALLA. Origen: anexo A del reporte S6a-contrato-c1 (corrección r1).
import argparse, json, re, os, sqlite3, random, sys, yaml
_ap = argparse.ArgumentParser(description='Verificación de C1 (C1-01 a C1-15, matriz §10 y S6-03)')
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
        if re.search(r'password|hash|token|secret', k): hits.append(path + '.' + k)
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
ok(len(ids) == len(set(ids)) and not badtag and [t['name'] for t in d['tags']] == TAGS9 and 'C1 v1.1' in c1 and d['info']['version'] == '1.1.0', 'C1-13 ids únicas, 9 etiquetas de la barra lateral, C1 v1.1')
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
# S6-03 paridad de rutas
rg = open(os.path.join(ROOT, 'backend/internal/api/routes.go'), encoding='utf-8').read()
rr = set(re.findall(r'Method: "(\w+)", Pattern: "([^"]+)", OperationID: "(\w+)"', rg))
oa = {(m, p, o['operationId']) for m, p, o in OPS}
print(f'S6-03 rutas en routes.go: {len(rr)} · en openapi.yaml: {len(oa)} · sobrantes: {sorted(rr - oa)[:3]} · faltantes: {sorted(oa - rr)[:3]}')
ok(rr == oa, 'S6-03 conjunto de rutas del router idéntico al contrato')
print('\nRESULTADO:', 'TODO OK' if not fails else f'{len(fails)} falla(s)')
sys.exit(1 if fails else 0)
