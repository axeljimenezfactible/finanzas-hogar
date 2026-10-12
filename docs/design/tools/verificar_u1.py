"""Verificación de U1 (docs/design/U1-arquitectura.md), versión 1.2.

Uso: python3 -I docs/design/tools/verificar_u1.py [raiz_repo]

Lee U1, `oraculo/model.json`, `datos/semilla/data.json`, `docs/contratos/C1-datos.md`,
`docs/contratos/openapi.yaml` y las migraciones `backend/db/migrations/00001_esquema.sql` y
`00002_puerta1.sql` (estos tres los carga el script por sí mismo: nadie necesita abrirlos; de las
migraciones solo busca las columnas de K18).
Acepta como excepciones técnicas las operaciones de la tabla de C1 §11.2 y como necesidades sin API
las que U1 marca «Local» en la columna Actualización. Las N-xxx nuevas de la versión vigente (1.2) que aún
no tienen operación se reportan como PENDIENTE C1 (esperado hasta C1 v1.2), no como falla.
Origen: anexo `u1_verify.py` de docs/estado/wf0/U1-arquitectura.md, ampliado en la puerta 1 y en la
corrección r1 (hallazgos M1 a M3 de docs/estado/p1/REVISION-U1-r1.md: N-090 con la regla de N-318,
F-7 con la rama «Faltan restricciones» y la regla del domingo escrita solo en §0.2 y §7) y en la
corrección r2 (hallazgo A1 de docs/estado/p1/REVISION-U1-r2.md: K18 en P-26.05, N-283, N-289 y P-18
con «Falta capturar», las columnas de S1 1.2 que lo sostienen y los campos que C1 v1.1 debe exponer,
reportados como PENDIENTE C1 mientras openapi.yaml no los tenga).
En la versión 1.2 (puerta 1b, flujo P1c) suma K38 (domingo FIRME), K43 (N-023, N-064, N-226, N-297, N-298 y
N-300 FIRMES), K44 (tabla «Parámetros operativos» de §2.3 y una necesidad por grupo), K36 (ajuste por
objetivo de la casa y por persona), K35 (Pista 3 en §7) y «Descartar chat detectado» (N-324).
Corrección r1 del flujo P1c (H-01 de docs/estado/p1c/REVISION-U1-r1.md): el aviso de pago próximo usa
`days_before` con el valor inicial que siembra 00001 ([1, 0], decisión F) en N-298, N-300 y la tabla de §2.3,
y `advance_days` es solo la anticipación del aviso de precios viejos, como define `00003_puerta1b.sql`
(el script lee 00001 y 00003 por sí mismo).
Corrección r2 del flujo P1c (H-r2-01 de docs/estado/p1c/REVISION-U1-r2.md): desde la corrección S1-r1, 00003 mueve
los días antes del pago próximo y los umbrales del presupuesto a claves de app_setting (`alerts.payment_due_days_before`
[1, 0], `alerts.budget_warning_pct` 80 y `alerts.budget_limit_pct` 100); el script lee sus líneas `@param` y exige
esas claves en N-298, N-300, N-323 y la tabla de §2.3. También comprueba que la sección «Para C1» que lee C1 v1.2
(la de docs/estado/p1c/U1-1.2.md si ya está corregida; si no, «Para C1 (corregida r2)» de
docs/estado/p1c/correccion-U1-r2.md) trae las N nuevas y las viñetas de N-297, N-298, N-300 y N-323 coherentes con
U1 1.2 y 00003, y la línea PENDIENTE apunta a esa sección.
"""
import re, json, sys, os, datetime as dt

AQUI = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.abspath(sys.argv[1]) if len(sys.argv) > 1 else os.path.abspath(os.path.join(AQUI, '..', '..', '..'))
doc = open(os.path.join(ROOT, 'docs/design/U1-arquitectura.md'), encoding='utf-8').read()
L = doc.split('\n')
fallas = []

def ok(c, msg):
    print(('OK   ' if c else 'FALLA') + ' ' + msg)
    if not c:
        fallas.append(msg)

def info(msg):
    print('INFO  ' + msg)

def section(text, start, end):
    a = text.index(start); b = text.index(end, a + len(start)); return text[a:b]

def sec(start, end):
    return section(doc, start, end)

def rows(text):
    return [[c.strip() for c in l.strip().strip('|').split('|')] for l in text.split('\n') if l.startswith('| ') and not l.startswith('|---')]

def lineas(pat, desde=0):
    return [i + 1 for i, l in enumerate(L) if i >= desde and re.search(pat, l)]

METODO = r'\b(GET|POST|PUT|PATCH|DELETE)\s+`?(/[^\s`|)]*)'
def norm(op):
    m, p = op.split(' ', 1)
    return m + ' ' + (p[len('/api/v1'):] or '/' if p.startswith('/api/v1') else p)

# Versión y lista de cambios
i0 = doc.index('## 0. Convenciones'); cab = doc[:i0]; cuerpo0 = doc[:i0].count('\n')
ok(re.search(r'\| Versión \| \*\*1\.2\*\*', cab) is not None and 0 <= cab.find('## Cambios de la versión 1.2') < cab.find('## Cambios de la versión 1.1'),
   f'Versión 1.2 con su lista de cambios al inicio, antes de la de 1.1 (línea {lineas(r"^## Cambios de la versión 1\.2")})')
mnu = re.search(r'las nuevas son (N-\d{3}) a (N-\d{3})', cab)

# U1-01 barra lateral
tree = sec('```\nFinanzas Hogar', '\n```')
top = []
for l in tree.split('\n'):
    if l.startswith('├─ ') or l.startswith('└─ '):
        top.append(re.split(r'\s{2,}|→', l[3:])[0].replace('▸', '').strip())
secc = [t for t in top if t not in ('Buscar', 'Registrar gasto')]
esperadas = ['Inicio', 'Gastos', 'Presupuesto', 'Comidas', 'Compras y tiendas', 'Créditos', 'Salud', 'Reportes', 'Ajustes']
ok(secc == esperadas, f'U1-01 barra lateral: {len(secc)} entradas de primer nivel en orden: {secc}')

# U1-02 y U1-08 inventario y Ajustes
inv = rows(sec('### 2.1 Qué es cada pantalla', '### 2.2'))[1:]
pids = [r[0] for r in inv]
ok(len(pids) == len(set(pids)), f'U1-02 ids de pantalla únicos: {len(pids)} ({pids[0]}..{pids[-1]})')
ok(all(len(r) >= 3 and r[1] and r[2] for r in inv), 'U1-02 cada pantalla con nombre y sección')
ajs = rows(sec('### 2.3 Ajustes: las 14 secciones', '### 2.4').split('\nEstados de Ajustes')[0])[1:]  # solo la primera tabla
nombres14 = ['Cuenta y roles', 'Presupuesto', 'Ingresos', 'Categorías', 'Créditos y condiciones (incluye promociones)', "Tiendas y precios (incluye mercados y membresía Sam's)", 'Recetas', 'Reportes y frecuencias', 'Notificaciones', 'Apariencia', 'Datos y respaldos', 'Personas y reparto', 'Salud y metas de ahorro', 'Inventario']
ok([r[1] for r in ajs] == nombres14 and [r[0] for r in ajs] == [f'P-26.{i:02d}' for i in range(1, 15)], f'U1-08 Ajustes: {len(ajs)} secciones con nombre y orden exactos')
todas_p = set(pids) | {r[0] for r in ajs}

# Necesidades N-xxx
nec = sec('## 3. Necesidades de datos', '### 3.29')
defs = re.findall(r'^\| (N-\d{3}) \| ([LE]) \|', nec, re.M)
nids = [d[0] for d in defs]; nL = sum(1 for d in defs if d[1] == 'L'); nE = len(defs) - nL
nset = set(nids)
ok(len(nids) == len(set(nids)), f'N-xxx definidas y únicas: {len(nids)} (L={nL}, E={nE}); última: {max(nids)}')
refs = set(re.findall(r'N-\d{3}', doc))
ok(refs <= nset, f'Referencias a N-xxx inexistentes: {sorted(refs - nset)}')
nuevas = {f'N-{i:03d}' for i in range(int(mnu.group(1)[2:]), int(mnu.group(2)[2:]) + 1)} if mnu else set()
ok(bool(nuevas) and nuevas <= nset and min(nuevas) > max(nset - nuevas), f'N-xxx nuevas de la versión vigente a continuación de la última, sin reusar: {sorted(nuevas)}')
locales = {r[0] for r in rows(nec) if re.fullmatch(r'N-\d{3}', r[0]) and r[-1] == 'Local'}
info(f'N-xxx Local (columna Actualización = Local): {sorted(locales)}')
m329 = sec('### 3.29', '### 3.30')
mc = re.search(r'\*\*(\d+)\*\* necesidades \(N-001 a (N-\d{3}); (\d+) de lectura y (\d+) de escritura\)', m329)
mk = re.search(r'Necesidades: \*\*(\d+)\*\* \(antes \d+\): (\d+) de lectura y (\d+) de escritura', cab)
ok(bool(mc and mk) and (int(mc.group(1)), mc.group(2), int(mc.group(3)), int(mc.group(4))) == (len(nids), max(nids), nL, nE) and (int(mk.group(1)), int(mk.group(2)), int(mk.group(3))) == (len(nids), nL, nE),
   f'Conteos declarados (§3.29 y lista de cambios) = definidos: {len(nids)} = {nL} L + {nE} E')

# Matriz 3.29
mat = rows(m329)[1:]
def expand(cell):
    out = set()
    for a, b in re.findall(r'(N-\d{3}) a (N-\d{3})', cell):
        out |= {n for n in nset if a <= n <= b}
    return out | set(re.findall(r'N-\d{3}', cell))
en_matriz = set(); pant_matriz = set()
for r in mat:
    pant_matriz |= set(re.findall(r'P-\d{2}(?:\.\d{2})?', r[0]))
    en_matriz |= expand(r[1]) | expand(r[2])
ok(nset <= en_matriz, f'Toda N-xxx aparece en la matriz 3.29; faltan: {sorted(nset - en_matriz)}')
falt_p = sorted((todas_p - {'P-26'}) - pant_matriz)
ok(not falt_p, f'Toda pantalla y sección de Ajustes está en la matriz; faltan: {falt_p}')

# Flujos (U1-04, U1-05, U1-15)
fl = sec('## 4. Flujos clave', '### Resumen de toques')
ok('K10' in fl.split('### F-1')[0], 'U1-04 regla de conteo de toques escrita con la decisión K10 (opción a)')
bloques = re.split(r'\n(?=### F-\d|\*\*F-\d)', fl)
ids_g = set(re.findall(r'^\| (G-\d{2}) \|', doc, re.M))
titulos = {}
for b in bloques:
    m = re.match(r'(?:### |\*\*)(F-\d)', b)
    if not m: continue
    fid = m.group(1); titulos[fid] = b.split('\n')[0]
    pasos = [r for r in rows(b) if r and r[0].isdigit()]
    sinp = [r[0] for r in pasos if not re.search(r'[PG]-\d{2}', r[1])]
    malos = [x for r in pasos for x in re.findall(r'[PG]-\d{2}(?:\.\d{2})?', r[1]) if x not in todas_p | ids_g]
    ok(pasos and not sinp and not malos, f'{fid}: {len(pasos)} pasos numerados, todos con pantalla existente (sin pantalla: {sinp}, inexistentes: {malos})')
    if fid == 'F-1':
        ok(pasos[-1][-1].strip('*') == '3', f'U1-04 F-1 registrar gasto: toques acumulados al guardar = {pasos[-1][-1]}')
    if fid in ('F-1', 'F-6'):
        pant = 'P-03' if fid == 'F-1' else 'P-22'
        v = [x for x in b.split('\n\n') if x.startswith('**Variante «Sin conexión»')]
        ok(bool(v) and pant in v[0] and 'N-317' in v[0] and '3 toques' in v[0], f'U1-15 {fid}: variante «Sin conexión» con {pant}, cola N-317 y 3 toques')
oblig = {'a': ('F-1', 'Registrar gasto'), 'b': (None, 'comparador'), 'c': (None, 'semana'), 'd': (None, 'icket')}
falt_o = [k for k, (fid, pat) in oblig.items() if not any((fid is None or f == fid) and pat in t for f, t in titulos.items())]
ok(len(titulos) >= 6 and not falt_o, f'U1-05 flujos con tabla: {len(titulos)} (≥ 6) y los 4 obligatorios (a) a (d); faltan: {falt_o}')
ok(all('idempotencia' in next(r for r in rows(nec) if r[0] == n)[-2] for n in ('N-054', 'N-240')) and 'N-317' in locales,
   'U1-15 N-054 y N-240 reciben clave de idempotencia y fecha de captura; N-317 (cola) es Local')
info(f'U1-15 líneas con «Sin conexión»: {lineas("Sin conexión", cuerpo0)}')

# U1-16 lista vigente y U1-11 objetivo táctil
ok(doc.count('**lista de compras vigente**') == 1 and len(lineas(r'\*\*lista de compras vigente\*\*', cuerpo0)) == 1, f'U1-16 «lista de compras vigente» definida una sola vez (línea {lineas(r"[*][*]lista de compras vigente[*][*]", cuerpo0)})')
lista = ['N-073', 'N-127', 'N-150', 'N-151', 'N-152', 'N-153', 'N-160', 'N-165']
sin02 = [n for n in lista if '§0.2' not in ' '.join(next(r for r in rows(nec) if r[0] == n))]
ok(not sin02, f'U1-16 necesidades de la lista con semana explícita que remiten a §0.2; faltan: {sin02}')
info(f'U1-16 líneas con «vigente»: {lineas("vigente", cuerpo0)}')
ok('44 × 44 px' in sec('### 6.1', '### 6.2'), 'U1-11 objetivo táctil ≥ 44 × 44 px en §6.1 (24 × 24 px solo como piso)')

# Corrección r1 · M3 (U1-14, U1-16): la regla PROVISIONAL del domingo solo se escribe en §0.2 y §7; el resto remite a §0.2
def rango(ini, fin):
    a = doc.index(ini) + 1; b = doc.find(fin, a)
    return doc[:a].count('\n') + 1, (doc[:b].count('\n') + 1 if b >= 0 else len(L))
r02, r7 = rango('\n### 0.2 ', '\n### 0.3 '), rango('\n## 7. ', '\n## ')
dom = [n for n in lineas(r'(?i)domingos?[\s,;:()]+(?:la\s+)?(?:semana\s+)?siguiente') if not (r02[0] <= n <= r02[1] or r7[0] <= n <= r7[1])]
ok(not dom, f'U1-16 la regla del domingo («domingos, la siguiente») solo aparece en §0.2 (líneas {r02[0]}-{r02[1]}) y §7 (desde {r7[0]}); fuera: {dom}')
pc = next((l for l in L if l.startswith('| plan_comida |') and 'periodo (semana)' in l), '')
f2 = next((r for r in rows(sec('### F-2', '### F-3')) if r and r[0] == '1'), [])
ok('§0.2' in pc and '§0.2' in ' '.join(f2), 'U1-16 la vista de plan_comida (§2.5) y el paso 1 de F-2 remiten a la semana por defecto de §0.2')

# Corrección r1 · M1 (U1-12, K21): N-090 y N-091 aplican en los meses 7 a 12 la regla de N-318
fila = lambda n: ' '.join(next(r for r in rows(nec) if r[0] == n))
ok(all('N-318' in fila(n) and 'K21' in fila(n) and 'Meses 7 a 12' in fila(n) for n in ('N-090', 'N-091'))
   and all('solo compromisos contratados' in fila(n) for n in ('N-090', 'N-091', 'N-318')),
   'U1-12 K21: N-090 (por quincena) y N-091 (por mes) usan en los meses 7 a 12 la regla de N-318 (solo compromisos contratados, gastos fijos e ingresos del CSV)')

# Corrección r1 · M2 (U1-05, U1-12, K16): F-7 no da por hecha la recomendación ni restricciones por defecto
f7 = sec('\n**F-7', '\n**F-8')
ok('Faltan restricciones' in f7 and 'K16' in f7 and 'pesos por defecto' not in f7 and re.search(r'^\| 11 \|[^\n]*F-7', cab, re.M) is not None,
   'U1-05 F-7 con la rama «Faltan restricciones» (K16), sin restricciones por defecto, y citado en el cambio 11')

# Corrección r2 · A1 (U1-09, U1-12, K18): datos del estado de cuenta de tarjeta bancaria y crédito de tienda
K18_DATOS = ('tasa anual', 'base diaria o mensual', 'incluye IVA', 'porcentaje de pago mínimo', 'mínimo fijo')
K18_COLS = ('rate_base', 'min_payment_pct_bp', 'min_payment_fixed_cents', 'annual_rate_decimals')
p2605 = ' '.join(next((r for r in ajs if r[0] == 'P-26.05'), []))
ok(all(c in p2605 for c in K18_DATOS) and all(k in p2605 for k in ('K18', 'K12', 'sin valor por defecto', 'Falta capturar')),
   'U1-09 K18: P-26.05 (§2.3) contiene los cinco datos del estado de cuenta, sin valor por defecto y con «Falta capturar»')
ok(all(c in fila('N-283') for c in K18_DATOS) and all(k in fila('N-283') for k in ('K18', 'K12', 'sin valor por defecto', 'Falta capturar')),
   'U1-09 K18: N-283 lee los cinco datos (la tasa anual como se capturó, K12), sin valor por defecto y con «Falta capturar»')
ok(all(k in fila('N-289') for k in ('K18', 'N-283', 'sin valor por defecto', 'Falta capturar')), 'U1-09 K18: N-289 edita los cinco datos de N-283 sin valor por defecto')
p18 = ' '.join(next((r for r in rows(sec('### 2.2 ', '### 2.3 ')) if r and r[0] == 'P-18'), []))
sin18 = [n for n in ('N-201', 'N-205') if not all(k in fila(n) for k in ('K18', 'pago mínimo', 'Falta capturar', 'P-26.05', 'valores supuestos'))]
ok(not sin18 and all(k in p18 for k in ('K18', 'Falta capturar', 'P-26.05')) and all(k in fila('N-200') for k in ('K18', 'N-283', 'Falta capturar')),
   f'U1-09 K18: P-18 (§2.2, N-200, N-201 y N-205) avisa «Falta capturar», abre P-26.05 y no calcula con valores supuestos; sin aviso: {sin18}')
ok('**Falta capturar** (K18)' in sec('### 2.4 ', '### 2.5 '), 'U1-09 K18: «Falta capturar» es un estado común (§2.4) para los intereses y el pago mínimo de esa línea')
c18 = [r for r in rows(cab) if len(r) >= 4 and 'K18' in r[2]]
ok(bool(c18) and all(x in c18[0][3] for x in ('P-26.05', 'P-18', 'N-200', 'N-201', 'N-205', 'N-283', 'N-289')),
   f'U1-09 K18 en la lista de cambios (fila {c18[0][0] if c18 else "—"}) con P-26.05, P-18, N-200, N-201, N-205, N-283 y N-289')
migs = [os.path.join(ROOT, 'backend/db/migrations', f) for f in ('00001_esquema.sql', '00002_puerta1.sql')]
if all(os.path.exists(f) for f in migs):
    sql = '\n'.join(open(f, encoding='utf-8').read() for f in migs)
    falt_cols = [c for c in K18_COLS if not re.search(rf'^ALTER TABLE credit ADD COLUMN {c}\b', sql, re.M)]
    ok(not falt_cols and re.search(r'^\s+vat_on_interest_bp INTEGER', sql, re.M) is not None,
       f'U1-12 K18: S1 1.2 tiene en credit las columnas que leen N-200 y N-283 y edita N-289 ({", ".join(K18_COLS)} y vat_on_interest_bp); faltan: {falt_cols}')
else:
    print('PENDIENTE U1-12 K18: no existen las migraciones 00001 y 00002')

# PROVISIONAL restantes: cada una cita un bloqueo existente; si el bloqueo está resuelto, la línea cita K31
bdir = os.path.join(ROOT, 'docs/estado/bloqueos'); malas_p = []
for n in lineas('PROVISIONAL', cuerpo0):
    bl = [b for b in re.findall(r'`([\w.-]+)`', L[n - 1]) if os.path.exists(os.path.join(bdir, b + '.md'))]
    estado = ['resuelto' if 'RESUELTO' in open(os.path.join(bdir, b + '.md'), encoding='utf-8').read() else 'abierto' for b in bl]
    bien = bool(bl) and all(e == 'abierto' or 'K31' in L[n - 1] for e in estado)
    info(f'PROVISIONAL línea {n}: bloqueo {bl} ({", ".join(estado) or "sin bloqueo"}){"" if bien else " ← sin justificar"}')
    if not bien: malas_p.append(n)
ok(not malas_p, f'Toda marca PROVISIONAL cita un bloqueo abierto o la regla K31; sin justificar: {malas_p}')

# Puerta 1b (U1 1.2): K38, K43, K44, K36, K35 y «Descartar chat detectado»
def fila_de(n):
    return next((r for r in rows(nec) if r and r[0] == n), [])
txt = lambda n: ' '.join(fila_de(n))
l02 = next((l for l in L if '**lista de compras vigente**' in l), '')
s7 = sec('\n## 7. ', '\nObservación para S1')
dom7 = next((r for r in rows(s7) if r and r[0].startswith('Semana por defecto los domingos')), [])
ok('K38' in l02 and '**FIRME**' in l02 and 'PROVISIONAL' not in l02 and bool(dom7) and 'K38' in dom7[1] and 'FIRME' in dom7[1] and 'PROVISIONAL' not in ' '.join(dom7),
   'K38 la regla del domingo es FIRME en §0.2 y en §7, sin marca PROVISIONAL')
ok(all('§0.2' in txt(n) and 'PROVISIONAL' not in txt(n) for n in ('N-120', 'N-150', 'N-160')), 'K38 N-120, N-150 y N-160 remiten a §0.2 sin marca PROVISIONAL')
k43 = {'N-023': ('FIRME', 'K43'), 'N-064': ('10 MB', 'K43', 'N-325'), 'N-226': ('5 s', 'K43', 'N-321'), 'N-297': ('send_time', 'K43'), 'N-298': ('advance_days', 'send_time', 'K43'), 'N-300': ('advance_days', 'send_time', 'K43')}
mal43 = [n for n, ks in k43.items() if not all(k in txt(n) for k in ks) or 'PROVISIONAL' in txt(n)]
ok(not mal43, f'K43 N-023, N-064, N-226, N-297, N-298 y N-300 FIRMES con sus valores (10 MB, 5 s, send_time, advance_days); mal: {mal43}')
t23 = sec('### 2.3 Ajustes: las 14 secciones', '### 2.4')
reglas = t23.split('**Parámetros operativos (K44).**')[1] if '**Parámetros operativos (K44).**' in t23 else ''
par = [r for r in rows(reglas) if len(r) >= 5 and r[0].startswith('P-26.')]
ESP44 = [('P-26.01', '5', 'N-320'), ('P-26.01', '15 min', 'N-320'), ('P-26.01', '90 días', 'N-320'), ('P-26.01', '12 caracteres', 'N-320'), ('P-26.05', '5 s', 'N-321'), ('P-26.06', '1 GB', 'N-322'),
         ('P-26.06', '30 días', 'N-322'), ('P-26.06', '90 días', 'N-322'), ('P-26.09', '80 % y 100 %', 'N-323'),
         ('P-26.11', '10 MB', 'N-325'), ('P-26.11', '24 h', 'N-325'), ('P-26.13', '0 kcal y 0 kcal', 'N-326')]
falt44 = [e for e in ESP44 if not any(r[0].startswith(e[0]) and r[2] == e[1] and e[2] in r[4] for r in par)]
sec_ok = all(r[0][:7] in {a[0] for a in ajs} and next(a[1] for a in ajs if a[0] == r[0][:7]).startswith(r[0][8:]) for r in par)
ok(bool(par) and not falt44 and sec_ok and all(k in reglas for k in ('Valor inicial', 'Sin valor inventado', 'Restablecer al valor inicial', 'Sin valor inicial', 'sin sección nueva', 'puerta 2')),
   f'K44 tabla «Parámetros operativos» (§2.3): {len(par)} filas, cada parámetro en una de las 14 secciones con su valor inicial, decisión y necesidad; reglas comunes escritas; faltan: {falt44}')
# Corrección r1 y r2 del flujo P1c (H-01 y H-r2-01): días antes del pago próximo y umbrales del presupuesto como claves
# alerts.* de app_setting (00003 tras la corrección S1-r1; origen en 00001, decisión F); advance_days solo de precios viejos
def _mig(nombre):
    p = os.path.join(ROOT, 'backend/db/migrations', nombre)
    return open(p, encoding='utf-8').read() if os.path.exists(p) else ''
_m1, _m3 = _mig('00001_esquema.sql'), _mig('00003_puerta1b.sql')
_pd = re.search(r"\('payment_due',[^)]*'(\{[^']*\})'\)", _m1)
db_ini = json.loads(_pd.group(1)).get('days_before') if _pd else None
def _param3(clave):
    l = next((l for l in _m3.split('\n') if re.match(rf'^-- @param {re.escape(clave)} \|', l)), '')
    c = [x.strip() for x in l.split('|')]
    return c[3] if len(c) > 3 else None
ALERTS = ('alerts.payment_due_days_before', 'alerts.budget_warning_pct', 'alerts.budget_limit_pct')
p3 = {k: _param3(k) for k in ALERTS}
ini_txt = p3['alerts.payment_due_days_before'] or '¿?'
adv3 = next((l for l in _m3.split('\n') if '@param_sin_valor' in l and 'advance_days' in l), '')
fp = next((r for r in par if 'pago próximo' in r[1]), [])
fu = next((r for r in par if 'Umbrales del aviso de presupuesto' in r[1]), [])
fa = next((r for r in par if 'advance_days' in r[1]), [])
mal_h01 = []
if db_ini != [1, 0]: mal_h01.append(f'00001 payment_due days_before = {db_ini}')
if p3 != dict(zip(ALERTS, ('[1, 0]', '80', '100'))): mal_h01.append(f'00003 @param alerts.* = {p3}')
if 'precios viejos' not in adv3: mal_h01.append('00003 no define advance_days como anticipación del aviso de precios viejos')
mal_h01 += [n for n in ('N-298', 'N-300') if not all(k in txt(n) for k in ('alerts.payment_due_days_before', ini_txt, 'decisión F', 'N-323', 'advance_days', 'precios viejos'))]
if not all(k in txt('N-323') for k in ALERTS + ('80 %', '100 %', ini_txt, 'decisión F', 'app_setting')): mal_h01.append('N-323')
mal_h01 += [r[0] for r in rows(nec) if r and '`days_before`' in ' '.join(r)] + [r[1] for r in par if '`days_before`' in ' '.join(r)]
if not (fp and 'alerts.payment_due_days_before' in fp[1] and ini_txt in fp[2] and 'N-323' in fp[4]): mal_h01.append('fila «pago próximo» de §2.3')
if not (fu and 'alerts.budget_warning_pct' in fu[1] and 'alerts.budget_limit_pct' in fu[1]): mal_h01.append('fila «umbrales» de §2.3')
if not (fa and 'precios viejos' in fa[1] and 'demás avisos' not in fa[1] and 'pago' not in fa[1]): mal_h01.append('fila `advance_days` de §2.3')
ok(not mal_h01, f'H-01/H-r2-01 (P1c) pago próximo con `alerts.payment_due_days_before` {ini_txt} y umbrales `alerts.budget_warning_pct` {p3["alerts.budget_warning_pct"]} y `alerts.budget_limit_pct` {p3["alerts.budget_limit_pct"]} (00003 tras S1-r1; origen 00001, decisión F) en N-298, N-300, N-323 y la tabla de §2.3, sin `days_before` suelto; `advance_days` solo del aviso de precios viejos; mal: {mal_h01}')
# H-r2-01 (P1c r2): la sección «Para C1» que lee C1 v1.2 trae las N nuevas y es coherente con U1 1.2 y 00003
FUENTES_C1 = [('docs/estado/p1c/U1-1.2.md', 'Para C1'), ('docs/estado/p1c/correccion-U1-r2.md', 'Para C1 (corregida r2)')]
def _para_c1(ruta, titulo):
    p = os.path.join(ROOT, ruta)
    t = open(p, encoding='utf-8').read() if os.path.exists(p) else ''
    return t.split(f'\n## {titulo}\n', 1)[1].split('\n## ', 1)[0] if f'\n## {titulo}\n' in t else ''
def _vineta(t, n):
    return next((l for l in t.split('\n') if l.startswith(f'- **{n}**')), '')
def _c1_ok(t):
    v, u = _vineta(t, 'N-297'), _vineta(t, 'N-323')
    return (bool(t) and all(n in t for n in nuevas) and 'hoy no tienen valor sembrado' not in v
            and all(k in v for k in ('alerts.payment_due_days_before', ini_txt, 'send_time', 'advance_days', 'S1-parametros-avisos'))
            and all(k in u for k in ALERTS))
fuente_c1 = next(((r, ti) for r, ti in FUENTES_C1 if _c1_ok(_para_c1(r, ti))), None)
ok(fuente_c1 is not None, f'H-r2-01 «Para C1» coherente con U1 1.2 y 00003 ({len(nuevas)} N nuevas; N-297, N-298 y N-300 con `alerts.payment_due_days_before` {ini_txt} para «Restablecer»; N-323 con las tres claves alerts.*); fuente: ' + (f'«{fuente_c1[1]}» de {fuente_c1[0]}' if fuente_c1 else f'ninguna de {FUENTES_C1}'))
mrow = {r[0]: expand(r[1]) | expand(r[2]) for r in mat}
grupos = {'N-320': '01', 'N-321': '05', 'N-322': '06', 'N-323': '09', 'N-325': '11', 'N-326': '13'}
mal_g = [n for n, sc in grupos.items() if not (len(fila_de(n)) >= 6 and fila_de(n)[1] == 'E' and fila_de(n)[2] == sc and 'valor inicial' in txt(n).lower()
         and 'restablece' in txt(n).lower() and 'K44' in txt(n) + ' '.join(r[3] for r in par if r[4] == n) and n in mrow.get(f'P-26.{sc}', set()))]
ok(not mal_g, f'K44 una necesidad de escritura por grupo ({", ".join(grupos)}) en su sección de §3.28, con valor inicial y restablecer, y en su fila de la matriz; mal: {mal_g}')
f327, f230 = txt('N-327').lower(), txt('N-230').lower()
ok('0 kcal y 0 kcal' in txt('N-326') and 'K36' in txt('N-326') and all(k in f327 for k in ('el de la casa', 'restablecer', 'kcal', 'k36'))
   and 'N-327' in mrow.get('P-21', set()) and 'N-327' in mrow.get('P-26.13', set()) and 'el de la casa' in f230 and 'k36' in f230 and 'N-327' in sec('### 3.23 ', '### 3.24 '),
   'K36 ajuste por objetivo: casa 0 y 0 kcal en P-26.13 (N-326); por persona en P-21 (N-327, §3.23) con «el de la casa» y «Restablecer»; N-230 lo muestra')
ok('K45' in s7 and 'docs/design/entrada/' in s7, 'K45 §7: la Pista 3 la ejecuta el orquestador con Artifacts; docs/design/entrada/ queda como brief de referencia (sustituye a K35)')
ok(fila_de('N-324')[1:3] == ['E', '09'] and 'Descartar' in txt('N-324') and 'N-324' in mrow.get('P-26.09', set()) and 'N-324' in txt('N-299'),
   'Descartar chat detectado: N-324 (E, P-26.09), citado en N-299 y en la matriz')

# Menú / y atajos
cmd = rows(sec('### 5.2 Lista de comandos', '### 5.3'))[1:]
ok(not [r[0] for r in cmd if not r[3] or r[3] == '—'], f'U1-07 comandos «/»: {len(cmd)} filas, todas con destino')
claves = [r[0] for r in rows(sec('### 5.3 Atajos', '### 5.4'))[1:]]
dup = sorted({c for c in claves if claves.count(c) > 1})
ok(not dup, f'U1-07 atajos: {len(claves)} filas, duplicados: {dup}')

# Cifras contra el oráculo y la semilla
def nums(o, acc):
    if isinstance(o, dict): [nums(v, acc) for v in o.values()]
    elif isinstance(o, list): [nums(v, acc) for v in o]
    elif isinstance(o, (int, float)) and not isinstance(o, bool): acc.add(round(float(o), 2))
    return acc
m = json.load(open(os.path.join(ROOT, 'oraculo/model.json'), encoding='utf-8'))
d = json.load(open(os.path.join(ROOT, 'datos/semilla/data.json'), encoding='utf-8'))
canon = nums(m, set()) | nums(d, set()); t = d['totales']
canon |= {round(t['walmart'] - t['bodega'], 2), round(t['walmart_identicos'] - t['bodega_identicos'], 2)}
cifras = re.findall(r'\$(\d{1,3}(?:,\d{3})*\.\d{2})', doc)
fuera = sorted({c for c in cifras if round(float(c.replace(',', '')), 2) not in canon | {99999.99, 0.0}})
ok(not fuera, f'Cifras en $ ({len(cifras)} apariciones, {len(set(cifras))} distintas) presentes en model.json/data.json o ejemplos de formato; fuera: {fuera}')

# Fechas
mes = dict(ene=1, feb=2, mar=3, abr=4, may=5, jun=6, jul=7, ago=8, sep=9, oct=10, nov=11, dic=12)
dsem = ['Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb', 'Dom']; malas = []
for wd, dd, mm, yy in re.findall(r'\b(Lun|Mar|Mié|Jue|Vie|Sáb|Dom) (\d{1,2}) (ene|feb|mar|abr|may|jun|jul|ago|sep|oct|nov|dic) (\d{4})', doc):
    if dsem[dt.date(int(yy), mes[mm], int(dd)).weekday()] != wd: malas.append(f'{wd} {dd} {mm} {yy}')
for dd, mm, yy in re.findall(r'\b(\d{1,2})-(ene|feb|mar|abr|may|jun|jul|ago|sep|oct|nov|dic)-(\d{4})', doc):
    try: dt.date(int(yy), mes[mm], int(dd))
    except ValueError: malas.append(f'{dd}-{mm}-{yy}')
ok(not malas, f'Fechas con día de la semana y fechas dd-mmm-aaaa válidas; erróneas: {malas}')

# Idioma y Karen
visibles = [r[1] for r in inv] + [r[1] for r in ajs] + secc + [r[2] for r in cmd]
ingles = re.compile(r'\b(Dashboard|Settings|Home|Save|Board|Gallery|Table|Calendar|Search|Login|Logout|Report|Budget|Credit|Shopping|Inventory|Health|Expense|Income|Recipe|Store)\b', re.I)
hits = sorted({v for v in visibles if ingles.search(v)})
ok(not hits, f'U1-13 nombres visibles sin palabras en inglés ({len(visibles)} revisados); coincidencias: {hits}')
ok(not any('Karen' in r[1] and 'Estado de cuenta' not in r[1] for r in inv), 'U1-10 ninguna pantalla propia de Karen salvo «Estado de cuenta de Karen»')

# U1-03 cruce con C1: openapi.yaml (x-necesidades por operación) + excepciones de C1 §11.2
oa = os.path.join(ROOT, 'docs/contratos/openapi.yaml'); c1 = os.path.join(ROOT, 'docs/contratos/C1-datos.md')
if os.path.exists(oa) and os.path.exists(c1):
    ver = None; ops = {}; cur = op = None
    for l in open(oa, encoding='utf-8'):
        l = l.rstrip('\n')
        if ver is None and re.match(r'^  version:', l): ver = l.split(':', 1)[1].strip().strip('"\'')
        if re.match(r'^\S', l): cur = op = None; continue
        mp = re.match(r'^  (/\S*):\s*$', l)
        if mp: cur, op = mp.group(1), None; continue
        if re.match(r'^  \S', l): cur = op = None; continue
        mm_ = re.match(r'^    (get|post|put|patch|delete):\s*$', l)
        if mm_ and cur: op = f'{mm_.group(1).upper()} {cur}'; ops[op] = {'id': None, 'n': set()}; continue
        if re.match(r'^    \S', l): op = None; continue
        if op:
            mi = re.match(r'^\s+operationId:\s*"?([\w-]+)"?', l)
            if mi and ops[op]['id'] is None: ops[op]['id'] = mi.group(1)
            ops[op]['n'].update(re.findall(r'^\s+- "?(N-\d{3})"?\s*$', l))
    c1t = open(c1, encoding='utf-8').read()
    exc = {norm(f'{a} {b}') for r in rows(section(c1t, '### 11.2', '### 11.3'))[1:] for a, b in re.findall(METODO, r[0])}
    n_ops = {}
    for o, v in ops.items():
        for n in v['n']: n_ops.setdefault(n, set()).add(o)
    info(f'openapi.yaml versión {ver}: {len(ops)} operaciones, {sum(1 for v in ops.values() if v["n"])} con x-necesidades; excepciones de C1 §11.2: {sorted(exc)}')
    sin_op = sorted(n for n in nset - locales if n not in n_ops)
    pend = sorted(set(sin_op) & nuevas); reales = sorted(set(sin_op) - nuevas)
    ok(not reales, f'U1-03 N-xxx sin operación ni Local (sin contar las nuevas de la versión vigente): {reales}')
    if pend: print(f'PENDIENTE C1 v1.2: {len(pend)} N-xxx nuevas de U1 1.2 sin operación todavía (esperado; van en ' + (f'«{fuente_c1[1]}» de {fuente_c1[0]}' if fuente_c1 else '«Para C1» (sin fuente coherente: ver H-r2-01)') + f'): {pend}')
    oal = open(oa, encoding='utf-8').read().split('\n')
    def bloque(nombre):
        if f'    {nombre}:' not in oal: return ''
        out = []
        for l in oal[oal.index(f'    {nombre}:') + 1:]:
            if l.strip() and not l.startswith('     '): break
            out.append(l)
        return '\n'.join(out)
    k18_falt = {s: [c for c in K18_COLS if not re.search(rf'^\s+{c}:\s*$', bloque(s), re.M)] for s in ('Credit', 'CreditCreate', 'CreditUpdate')}
    k18_falt = {s: v for s, v in k18_falt.items() if v}
    if k18_falt: print(f'PENDIENTE C1 v1.1: campos de K18 que N-200, N-283 y N-289 necesitan y openapi.yaml {ver} aún no expone (esperado; van en «Para C1» de docs/estado/p1/correccion-U1-r2.md): {k18_falt}')
    else: ok(True, 'U1-03 K18: Credit, CreditCreate y CreditUpdate de openapi.yaml exponen ' + ', '.join(K18_COLS))
    ok(not sorted(n for n in locales if n in n_ops), f'U1-03 ninguna N-xxx Local tiene operación: {sorted(n for n in locales if n in n_ops)}')
    sin_n = sorted(f'{o} ({ops[o]["id"]})' for o in ops if not ops[o]['n'] and norm(o) not in exc)
    ok(not sin_n, f'U1-03 operaciones de openapi.yaml sin N-xxx (fuera de las excepciones de C1 §11.2): {sin_n}')
    exc_con_n = sorted(o for o in ops if norm(o) in exc and ops[o]['n'])
    ok(not (exc - {norm(o) for o in ops}) and not exc_con_n, f'U1-03 excepciones de C1 §11.2 existen en openapi.yaml y no tienen N-xxx (con N-xxx: {exc_con_n})')
    ok(not sorted(set(n_ops) - nset), f'U1-03 N-xxx citadas por openapi.yaml que no existen en U1: {sorted(set(n_ops) - nset)}')
    pant_sin = [r[0] for r in mat if not any(n in n_ops for n in expand(r[1]) | expand(r[2])) and not ((expand(r[1]) | expand(r[2])) <= locales)]
    ok(not pant_sin, f'U1-03 filas de la matriz sin ninguna N-xxx con operación ni Local: {pant_sin}')
    c1map = {}
    for l in c1t.split('\n'):
        for n in re.findall(r'N-\d{3}', l): c1map.setdefault(n, set()).update(norm(f'{a} {b}') for a, b in re.findall(METODO, l))
    difs = sorted(n for n in n_ops if c1map.get(n) and not ({norm(o) for o in n_ops[n]} & c1map[n]))
    info(f'Cruce openapi.yaml ↔ C1-datos.md: N-xxx cuyas operaciones no coinciden en ninguna: {difs}')
else:
    print('PENDIENTE U1-03: docs/contratos/openapi.yaml o C1-datos.md no existen')
print('Menciones de «Karen» en el documento:', doc.count('Karen'))
print('\nRESULTADO:', 'TODO OK' if not fallas else f'{len(fallas)} falla(s)')
sys.exit(1 if fallas else 0)
