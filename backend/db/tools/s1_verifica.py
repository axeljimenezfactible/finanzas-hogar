# s1_verifica.py · Verificación reejecutable de S1 (esquema 1.2, diccionario y reglas) con ida y vuelta del oráculo.
# Uso: python3 -I backend/db/tools/s1_verifica.py   (crea datos/tmp/s1-verifica.db, la borra al final; sale con 1 si algo falla)
# Aplica TODAS las migraciones de backend/db/migrations en orden. goose: variable GOOSE o goose en el PATH.
# Origen: Anexo B de docs/estado/wf0/S1-modelo-datos.md (WF-0). Versión 1.2 (agente S1-esquema-1.2, flujo P1):
# Down de la 00002 contra el esquema 1.1 byte por byte, goose up / down-to 0 / up, regresión Kueski sobre 1.2,
# K12 en pantalla, PROVISIONAL solo con bloqueo abierto y pruebas de las columnas y la tabla nuevas.
# Revisión S1-r1 (agente correccion-S1-r1): S1-03 marca las cifras de dos decimales que no están en model.json ni en
# data.json dentro de los ejemplos sin rótulo, acepta las cuentas citadas («ruta / n = valor») y cuenta aparte los
# ilustrativos; NEG12 prueba que shopping_list empiece en lunes (K8).
# Revisión S1-r2 (agente correccion-S1-r2, hallazgo R2-01): S1-03 también marca las fechas (ISO, dd-mmm-aaaa y dd-mmm
# con el año deducido de la frase) que no están en model.json ni en data.json dentro de frases sin rótulo, y las fechas
# imposibles aunque estén rotuladas; una autoprueba con las frases que señaló la revisión comprueba que la marca funciona.
# Versión 1.3 (agente S1-esquema-1.3, flujo P1c): aplica también la 00003 (K36 y K44). DOWN13 prueba que el Down de la 00003
# deja sqlite_master, app_setting y los datos de 1.2 idénticos; DDL comprueba que el DDL de 00001 y 00002 (sin comentarios)
# es el de git HEAD; K44 prueba columnas y claves nuevas con sus valores iniciales, las negativas (ajuste −1, parámetro que no
# es entero), las cifras de K36, K39 y K40 en las reglas y que ninguna marca PROVISIONAL cite un bloqueo ya resuelto.
# Corrección S1-r1 de P1c (H-01 y H-02): K44 incluye las claves alerts.* y exige que params_json ya no las lleve; AVISOS13
# prueba el traspaso params_json <-> app_setting con valores que no son los iniciales; S1-13 exige cada clave @param en la
# fila de su sección de «Ajustes: 14 secciones»; DOWN13 exige los avisos de 1.2 byte por byte (params_json y updated_at).
# S1-03 acepta cuentas citadas encadenadas («ruta × 7 / 31 = valor», un solo redondeo HALF_UP al final, K39).
# Versión 1.4 (agente S1-puerta2, flujo P2c): aplica también la 00004 (puerta 2: Q2-11 B/K47, K50, Q2-66 B, Q2-67). K47 prueba las
# claves nuevas con su valor inicial, sus rangos y horas (negativas y positivas) y credit.vat_included; AVISOS14 el traspaso
# params_json <-> app_setting de weekday y day_of_month con valores que no son los iniciales y la migración de vat_included;
# PROG14 sustituye a PROG13 (los cuatro parámetros de avisos ya no son PROVISIONAL); DOWN14 prueba que el Down de la 00004 deja
# el esquema 1.3, sus ajustes, sus avisos y sus créditos idénticos. S1-16 y K44 tratan como historia las marcas PROVISIONAL de
# una migración anterior que citan un bloqueo que una migración posterior declara con «-- @resuelve» (00001 a 00003 no se editan).
# Versión 1.5 (agente s1-c1-v15, WF-3c, K54): aplica también la 00005 (tablero de control). K44 cuenta 30 anotaciones @param (25 + 5 de
# K54); K54 prueba las cinco claves nuevas con su valor inicial y sus rangos (negativas y positivas); PROG14 exige «sin PROVISIONAL»
# solo en las filas de las claves de la 00004 (las de K54 tienen valores PROVISIONALES con bloqueo abierto, C1-tablero y U3-tablero-*).
import datetime, functools, glob, json, os, re, shutil, sqlite3, subprocess, sys
from fractions import Fraction
R = os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__)))))
MIGS = sorted(glob.glob(R + '/backend/db/migrations/*.sql'))
DIC = open(R + '/docs/contratos/diccionario-datos.md', encoding='utf-8').read()
REG = open(R + '/docs/contratos/reglas-negocio.md', encoding='utf-8').read()
m = json.load(open(R + '/oraculo/model.json', encoding='utf-8'))
dj = json.load(open(R + '/datos/semilla/data.json', encoding='utf-8'))
TMP = R + '/datos/tmp'; DB = TMP + '/s1-verifica.db'
os.makedirs(TMP, exist_ok=True)
if os.path.exists(DB): os.remove(DB)
TXT = [open(q, encoding='utf-8').read() for q in MIGS]
UPS = [q.split('-- +goose Down', 1)[0] for q in TXT]; DOWNS = [q.split('-- +goose Down', 1)[1] for q in TXT]
mig = '\n'.join(TXT); up = '\n'.join(UPS)
def esquema(con): return con.execute("SELECT type, name, tbl_name, sql FROM sqlite_master WHERE name NOT LIKE 'sqlite_%' ORDER BY type, name").fetchall()
def ajustes(con): return con.execute('SELECT setting_key, user_account_id, value_json, description FROM app_setting ORDER BY setting_key').fetchall()
db = sqlite3.connect(DB); db.execute('PRAGMA foreign_keys=ON'); x = db.execute
FALLAS = []
def out(cid, ok, txt):
    print(f"{cid} {'OK  ' if ok else 'FALLA'} {txt}")
    if not ok: FALLAS.append(cid)
c = lambda v: int(Fraction(str(v)) * 100)  # pesos con 2 decimales -> centavos exactos
def one(q, a=()): return x(q, a).fetchone()[0]
def lid(): return one('SELECT last_insert_rowid()')
db.executescript(UPS[0]); E11 = esquema(db); A11 = ajustes(db)
db.executescript(UPS[1]); E12 = esquema(db); A12 = ajustes(db)
cols12 = {t: [r[1] for r in x(f'PRAGMA table_info({t})')] for (t,) in x("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").fetchall()}
N12 = x(f"SELECT {', '.join(cols12['notification'])} FROM notification ORDER BY rowid").fetchall()  # avisos de 1.2: el Down de la 00003 los devuelve byte por byte
db.executescript(UPS[2]); E13 = esquema(db); A13 = ajustes(db)
cols13t = {t: [r[1] for r in x(f'PRAGMA table_info({t})')] for (t,) in x("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").fetchall()}
N13 = x(f"SELECT {', '.join(cols13t['notification'])} FROM notification ORDER BY rowid").fetchall()  # avisos de 1.3: el Down de la 00004 los devuelve byte por byte
for u in UPS[3:]: db.executescript(u)
PATB = r'\b((?:S1|C1|U1|U2|U2a|S6|criterios)-[a-z0-9-]+[a-z0-9])\b'
RESUELVE = {}
for i_, t_ in enumerate(TXT):
    for b_ in re.findall(r'^-- @resuelve ([A-Za-z0-9-]+)', t_, re.M): RESUELVE.setdefault(b_, i_)
BLQ_FS = set(os.listdir(R + '/docs/estado/bloqueos'))
def historia(i_, l, b):  # marca PROVISIONAL de una migración ya cerrada que cita solo bloqueos que resuelve una migración posterior
    nm = {z for z in (set(re.findall(PATB, l)) or set(re.findall(PATB, b))) if z + '.md' in BLQ_FS}
    return bool(nm) and all(RESUELVE.get(z, -1) > i_ for z in nm)
EF = esquema(db)
tabs = [r[0] for r in x("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name")]
cols = {t: [r[1] for r in x(f'PRAGMA table_info({t})')] for t in tabs}
def has(tc):
    t, _, col = tc.partition('.')
    return t in cols and (not col or col in cols[t])
# ---------- S1-02
ni = one("SELECT count(*) FROM sqlite_master WHERE type='index'"); nt = one("SELECT count(*) FROM sqlite_master WHERE type='trigger'")
ns = one("SELECT count(*) FROM pragma_table_list WHERE schema='main' AND strict=1 AND name NOT LIKE 'sqlite_%'")
ic = one('PRAGMA integrity_check'); fk = x('PRAGMA foreign_key_check').fetchall()
out('S1-02', ic == 'ok' and not fk and ns == len(tabs), f"Up de {', '.join(os.path.basename(q) for q in MIGS)} aplicado con sqlite3 {sqlite3.sqlite_version}: tablas={len(tabs)} strict={ns} índices={ni} disparadores={nt} · integrity_check={ic} · foreign_key_check={len(fk)} filas")
# ---------- S1-01
ENT = ['persona', 'cuenta_usuario', 'sesion_dispositivo', 'periodo', 'ingreso', 'categoria', 'gasto', 'regla_reparto', 'meta_ahorro', 'movimiento_ahorro', 'credito', 'tasa_por_plazo', 'compra_credito', 'pago_programado', 'promocion', 'tienda', 'producto', 'presentacion', 'precio', 'receta', 'ingrediente_receta', 'tiempo_comida', 'plan_comida', 'inventario_item', 'perfil_salud', 'peso_historico', 'nutricion_ingrediente', 'notificacion', 'destinatario', 'job', 'respaldo']
mapa = re.findall(r'^\| ([^|]+?) \| `([a-z_]+)` \| [^|]+ \|$', DIC.split('## Mapa entidad del plan')[1].split('Tablas complementarias')[0], re.M)
ent2tab = {}
for e, t in mapa:
    for name in re.split(r'[ ,]+', e.split(' (')[0]):
        ent2tab.setdefault(name, t)
falt = [e for e in ENT if e not in ent2tab or ent2tab[e] not in tabs]
sin_fila = [t for t in tabs if f'### `{t}` ·' not in DIC]
out('S1-01', not falt and not sin_fila, f"{len(ENT) - len(falt)} de {len(ENT)} entidades con tabla: " + ', '.join(f'{e}→{ent2tab.get(e)}' for e in ENT) + f" · faltantes={falt} · tablas sin fila en el diccionario={sin_fila}")
# ---------- S1-04
usado = re.findall(r'^### `([a-z_]+)` ·.*?\n- \*\*Usado por:\*\* *(.*)$', DIC, re.M | re.S)
vacias = [t for t, u in re.findall(r'^### `([a-z_]+)` ·[^\n]*\n(?:[^\n]*\n){2}- \*\*Usado por:\*\* *([^\n]*)', DIC, re.M) if not u.strip()]
n_usado = len(re.findall(r'^- \*\*Usado por:\*\* \S', DIC, re.M))
out('S1-04', n_usado == len(tabs) and not vacias, f"tablas con «Usado por» lleno: {n_usado} de {len(tabs)} · vacías={vacias}")
# ---------- S1-05
reales = x("SELECT m.name, p.name, p.type FROM sqlite_master m, pragma_table_info(m.name) p WHERE m.type='table' AND upper(p.type) IN ('REAL','FLOAT','DOUBLE')").fetchall()
cents = {}
for t in tabs:
    for r in x(f'PRAGMA table_info({t})'):
        if r[1].endswith('_cents'): cents.setdefault(r[2], []).append(f'{t}.{r[1]}')
out('S1-05', not reales and list(cents) == ['INTEGER'], f"columnas REAL/FLOAT/DOUBLE={len(reales)} · columnas _cents por tipo: " + ', '.join(f'{k}={len(v)}' for k, v in cents.items()) + " · convención: enteros de centavos")
# ---------- S1-06
sin_coment = "\n".join(l.split(" -- ")[0] for l in up.splitlines() if not l.lstrip().startswith("--"))
pant = re.findall(r"\d{1,2}-(?:ene|feb|mar|abr|may|jun|jul|ago|sep|oct|nov|dic)-20\d\d", sin_coment)
crit = len(re.findall(r"oct-2026|ene-20", mig))
out('S1-06', not pant, f"fechas en formato de pantalla en datos o código SQL (sin comentarios): {len(pant)} · grep literal del criterio «oct-2026|ene-20» en la migración: {crit} · zona America/Mexico_City declarada en el diccionario: {'America/Mexico_City' in DIC}")
# ---------- datos de la ida y vuelta (S1-12) ----------
H = '$argon2id$v=19$m=65536,t=3,p=4$czFwcnVlYmE$cHJ1ZWJh'  # cadena de prueba con formato PHC; no es credencial
x("INSERT INTO person(code,display_name,participation) VALUES ('axel','Axel','holder')")
x("INSERT INTO person(code,display_name,participation) VALUES ('ximena','Ximena','co_holder')")
x("INSERT INTO person(code,display_name,participation,pays_share_directly) VALUES ('karen','Karen','external',1)")
P = dict(x('SELECT code,id FROM person'))
for u in ('axel', 'ximena'): x('INSERT INTO user_account(person_id,username,password_hash) VALUES (?,?,?)', (P[u], u, H))
x("INSERT INTO cost_share_rule(name,remainder_person_id,is_template,is_default) VALUES ('Hogar',?,1,1)", (P['axel'],)); r_hogar = lid()
x("INSERT INTO cost_share_rule(name,remainder_person_id,is_template,is_default) VALUES ('Conjunto 50/50',?,1,0)", (P['axel'],)); r_conj = lid()
x('INSERT INTO cost_share_rule_part(rule_id,person_id,pct_bp) VALUES (?,?,5000)', (r_conj, P['karen']))
x("INSERT INTO category(code,name,kind,default_cost_share_rule_id) VALUES ('groceries','Despensa','variable',?)", (r_conj,)); cat = lid()
x("INSERT INTO category(code,name,kind) VALUES ('household_goods','Papel','variable')"); cat2 = lid()
x("INSERT INTO credit(name,kind,lender,holder_person_id,credit_limit_cents,payment_frequency,first_payment_min_days,reading_a_anchor_date) VALUES ('Kueski Pay','revolving_line','Kueski',?,?,'fortnightly',15,'2026-10-10')", (P['axel'], c(m['parametros']['limite_credito']))); cr = lid()
x("INSERT INTO credit_term_rate(credit_id,term_count,rate_num,rate_den,source) VALUES (?,2,53450,309450,'manual'),(?,4,80050,309450,'manual')", (cr, cr))
for d in m['despensas']:
    x("INSERT INTO expense(kind,description,category_id,amount_cents,expense_date,accrual_months,payment_method,cost_share_rule_id) VALUES ('variable',?,?,?,?,2,'credit',?)", (f"Despensa {d['n']}", cat, c(d['capital']), d['fecha'], r_conj)); eid = lid()
    for w in ('axel', 'karen'): x('INSERT INTO expense_share(expense_id,person_id,amount_cents) VALUES (?,?,?)', (eid, P[w], c(d[w]['capital'])))
    x("INSERT INTO credit_purchase(credit_id,expense_id,kind,purchase_date,principal_cents,interest_cents,total_cents,term_count,payment_frequency,first_payment_date,interest_source,interest_includes_vat) VALUES (?,?,'installment_plan',?,?,?,?,4,'fortnightly',?,'quoted',1)", (cr, eid, d['fecha'], c(d['capital']), c(d['intereses']), c(d['total']), d['pagos'][0]['fecha'])); pid = lid()
    for w in ('axel', 'karen'): x('INSERT INTO credit_purchase_share(credit_purchase_id,person_id,principal_cents,interest_cents,total_cents) VALUES (?,?,?,?,?)', (pid, P[w], c(d[w]['capital']), c(d[w]['intereses']), c(d[w]['total'])))
    for i, p in enumerate(d['pagos'], 1):
        x('INSERT INTO scheduled_payment(credit_purchase_id,seq,due_date,amount_cents,principal_cents,interest_cents) VALUES (?,?,?,?,?,?)', (pid, i, p['fecha'], c(p['monto']), c(p['capital_repuesto']), c(p['monto']) - c(p['capital_repuesto']))); sid = lid()
        for w in ('axel', 'karen'): x('INSERT INTO scheduled_payment_share(scheduled_payment_id,person_id,amount_cents) VALUES (?,?,?)', (sid, P[w], c(p[w])))
for e in m['credito']['eventos']:
    if e['tipo'] == 'papel':
        x("INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method,paid_by_person_id,person_id) VALUES ('variable',?,?,?,'cash',?,?)", (cat2, c(e['monto']), e['fecha'], P['karen'], P['karen'])); eid = lid()
        x('INSERT INTO expense_share(expense_id,person_id,amount_cents) VALUES (?,?,?)', (eid, P['karen'], c(e['monto'])))
db.commit()
n = 0; dif = []
def eq(a, b, what):
    global n; n += 1
    if a != b: dif.append((what, a, b))
rows = x('SELECT id,purchase_date,principal_cents,interest_cents,total_cents FROM credit_purchase ORDER BY purchase_date').fetchall()
eq(len(rows), 4, 'número de despensas')
npag = 0
for d, (pid, f, pr, it, tt) in zip(m['despensas'], rows):
    eq((f, pr, it, tt), (d['fecha'], c(d['capital']), c(d['intereses']), c(d['total'])), f"despensa {d['n']}")
    for w in ('axel', 'karen'):
        eq(x('SELECT principal_cents,interest_cents,total_cents FROM credit_purchase_share WHERE credit_purchase_id=? AND person_id=?', (pid, P[w])).fetchone(), (c(d[w]['capital']), c(d[w]['intereses']), c(d[w]['total'])), f"despensa {d['n']} parte {w}")
        eq(one('SELECT amount_cents FROM expense_share s JOIN credit_purchase cp ON cp.expense_id=s.expense_id WHERE cp.id=? AND s.person_id=?', (pid, P[w])), c(d[w]['capital']), f"gasto {d['n']} parte {w}")
    for p, (sid, due, amt, prin) in zip(d['pagos'], x('SELECT id,due_date,amount_cents,principal_cents FROM scheduled_payment WHERE credit_purchase_id=? ORDER BY seq', (pid,)).fetchall()):
        npag += 1
        eq((due, amt, prin), (p['fecha'], c(p['monto']), c(p['capital_repuesto'])), f'pago {due}')
        for w in ('axel', 'karen'): eq(one('SELECT amount_cents FROM scheduled_payment_share WHERE scheduled_payment_id=? AND person_id=?', (sid, P[w])), c(p[w]), f'parte pago {due} {w}')
eq(npag, 16, 'número de pagos')
s_d1 = one("SELECT sum(sp.amount_cents) FROM scheduled_payment sp JOIN credit_purchase cp ON cp.id=sp.credit_purchase_id WHERE cp.purchase_date='2026-10-10'")
a_d1 = one("SELECT sum(s.amount_cents) FROM scheduled_payment_share s JOIN scheduled_payment sp ON sp.id=s.scheduled_payment_id JOIN credit_purchase cp ON cp.id=sp.credit_purchase_id WHERE cp.purchase_date='2026-10-10' AND s.person_id=?", (P['axel'],))
eq(s_d1, 389500, 'suma pagos d1'); eq(a_d1, 239528, 'parte Axel d1')
lim, anc, lect = x('SELECT credit_limit_cents, reading_a_anchor_date, availability_reading FROM credit WHERE id=?', (cr,)).fetchone()
fin = m['parametros']['horizonte'][1]
viva = one('SELECT sum(principal_cents) FROM credit_purchase WHERE purchase_date <= ?', (anc,))
A, B = lim, lim - viva
ev = [(anc, 0, 'inicio', 0)]
ev += [(f, 1, 'compra', p) for f, p in x('SELECT purchase_date, principal_cents FROM credit_purchase WHERE purchase_date > ? AND purchase_date <= ?', (anc, fin))]
ev += [(f, 2, 'pago', p) for f, p in x('SELECT due_date, principal_cents FROM scheduled_payment WHERE due_date > ? AND due_date <= ?', (anc, fin))]
ev += [(f, 3, 'papel', 0) for (f,) in x("SELECT expense_date FROM expense WHERE payment_method='cash' AND expense_date > ? AND expense_date <= ?", (anc, fin))]
ev.sort(); serie = []; sA = []; sB = []
for f, _, t, p in ev:
    if t == 'compra': A -= p; B -= p
    elif t == 'pago': A += p; B += p
    serie.append((f, t, A, B))
    if t in ('compra', 'pago'): sA.append(A); sB.append(B)
eo = m['credito']['eventos']; eq(len(serie), len(eo), 'número de eventos')
for (f, t, a, b), e in zip(serie, eo): eq((f, t, a, b), (e['fecha'], e['tipo'], c(e['disponible_A']), c(e['disponible_B'])), f'evento {f}')
cs = m['credito']['comparacion_series']
for v1, v2, k in ((sA, cs['esperado_A'], 'A'), (sB, cs['esperado_B'], 'B')):
    eq(len(v1), 14, f'largo serie {k}')
    for i, (a_, b_) in enumerate(zip(v1, v2)): eq(a_, c(b_), f'serie {k}[{i}]')
eq(min(sA + [lim]), c(m['credito']['minimo_A']), 'mínimo A'); eq(min(sB + [lim - viva]), c(m['credito']['minimo_B']), 'mínimo B')
acum = Fraction(0); per_ok = 0
for p in m['periodos']:
    ym = p['inicio'][:7]
    def mes(q, a=()): return one(q, a) or 0
    efs = {w: mes("SELECT sum(s.amount_cents) FROM scheduled_payment_share s JOIN scheduled_payment sp ON sp.id=s.scheduled_payment_id WHERE strftime('%Y-%m',sp.due_date)=? AND s.person_id=?", (ym, P[w])) + mes("SELECT sum(s.amount_cents) FROM expense_share s JOIN expense e ON e.id=s.expense_id WHERE e.payment_method='cash' AND strftime('%Y-%m',e.expense_date)=? AND s.person_id=?", (ym, P[w])) for w in ('axel', 'karen')}
    dvs = {w: sum((Fraction(v, 2) for (v,) in x("SELECT s.total_cents FROM credit_purchase_share s JOIN credit_purchase cp ON cp.id=s.credit_purchase_id WHERE s.person_id=? AND (strftime('%Y-%m',cp.purchase_date)=? OR strftime('%Y-%m',date(cp.purchase_date,'start of month','+1 month'))=?)", (P[w], ym, ym))), Fraction(0)) + mes("SELECT sum(s.amount_cents) FROM expense_share s JOIN expense e ON e.id=s.expense_id WHERE e.payment_method='cash' AND strftime('%Y-%m',e.expense_date)=? AND s.person_id=?", (ym, P[w])) for w in ('axel', 'karen')}
    ef, dv = sum(efs.values()), sum(dvs.values())
    eq(ef, c(p['efectivo']['total']), f"efectivo {p['periodo']}"); eq(dv, c(p['devengado']['total']), f"devengado {p['periodo']}")
    per_ok += (ef == c(p['efectivo']['total'])) + (dv == c(p['devengado']['total']))
    for w in ('axel', 'karen'):
        eq(efs[w], c(p['efectivo'][w]), f"efectivo {p['periodo']} {w}"); eq(dvs[w], c(p['devengado'][w]), f"devengado {p['periodo']} {w}")
    acum += dv - ef; eq(acum, c(p['diferencia_acumulada']), f"diferencia acumulada {p['periodo']}")
fechas_1_9 = [f for (f,) in x("SELECT purchase_date FROM credit_purchase UNION ALL SELECT due_date FROM scheduled_payment UNION ALL SELECT expense_date FROM expense WHERE payment_method='cash'") if int(f[8:]) <= 9]
out('S1-12', not dif, f"ida y vuelta de oraculo/model.json: 4 despensas, {npag} pagos, {len(serie)} eventos, series A/B de {len(sA)}/{len(sB)} valores, mínimos {min(sA)}/{min(sB)}, P1–P6 por persona · valores comparados={n} · diferencias={len(dif)} · suma pagos d1={s_d1} · parte Axel d1={a_d1}" + (f' · {dif[:5]}' if dif else ''))
out('KUESKI', s_d1 == 389500 and a_d1 == 239528 and len(MIGS) >= 2 and not dif, f"regresión Kueski sobre el esquema 1.4 ({', '.join(os.path.basename(q) for q in MIGS)}): suma de pagos de la despensa 1 = {s_d1} (esperado 389500) · parte de Axel = {a_d1} (esperado 239528) · diferencias contra model.json = {len(dif)}")
# ---------- K12: la tasa se muestra con los decimales capturados y se calcula con la fracción
x('UPDATE credit_term_rate SET rate_decimals = CASE term_count WHEN 2 THEN 2 ELSE 4 END WHERE credit_id = ?', (cr,))
def pct_txt(num, den, d):
    n_ = int(Fraction(num * 100, den) * 10 ** d + Fraction(1, 2))  # HALF_UP para positivos
    s_ = str(n_).rjust(d + 1, '0')
    return s_[:-d] + '.' + s_[-d:] if d else s_
mostradas = {tt: pct_txt(nm, dn, dd) for tt, nm, dn, dd in x('SELECT term_count, rate_num, rate_den, rate_decimals FROM credit_term_rate WHERE credit_id = ? ORDER BY term_count', (cr,))}
i4000 = int(Fraction(400000 * 80050, 309450) + Fraction(1, 2))
out('K12', mostradas == {2: '17.27', 4: '25.8685'} and i4000 == c(m['interes_por_monto'][5]['intereses']), f"tasas mostradas con rate_decimals: 2 quincenas → {mostradas.get(2)} %, 4 quincenas → {mostradas.get(4)} % · intereses de 4,000.00 con la fracción 80050/309450 = {i4000} centavos (model.json: interes_por_monto[5].intereses → {c(m['interes_por_monto'][5]['intereses'])})")
db.commit()
# ---------- S1-09
mes_ids = {}
for y_, mo in ((2026, 10), (2026, 11), (2026, 12), (2027, 1), (2027, 2), (2027, 3)):
    ini = f'{y_:04d}-{mo:02d}-01'; finm = one("SELECT date(?, '+1 month', '-1 day')", (ini,))
    x("INSERT INTO period(kind,year,month,start_date,end_date) VALUES ('month',?,?,?,?)", (y_, mo, ini, finm)); mid = lid(); mes_ids[(y_, mo)] = mid
    x("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES ('fortnight',?,?,?,1,?,?,?)", (mid, y_, mo, ini, f'{y_:04d}-{mo:02d}-15', f'{y_:04d}-{mo:02d}-15'))
    x("INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES ('fortnight',?,?,?,2,?,?,?)", (mid, y_, mo, f'{y_:04d}-{mo:02d}-16', finm, min(f'{y_:04d}-{mo:02d}-30', finm)))
filas = x("SELECT kind, year, month, half, start_date, end_date, pay_date FROM period ORDER BY start_date, kind DESC").fetchall()
feb2 = one("SELECT pay_date FROM period WHERE kind='fortnight' AND year=2027 AND month=2 AND half=2")
try:
    x("INSERT INTO period(kind,year,month,start_date,end_date) VALUES ('month',2026,10,'2026-10-10','2026-11-09')"); ventana = 'ACEPTADA'
except sqlite3.IntegrityError: ventana = 'rechazada'
out('S1-09', len(filas) == 18 and feb2 == '2027-02-28' and ventana == 'rechazada' and per_ok == 12 and not fechas_1_9,
    f"{len(filas)} filas recuperadas (6 meses oct-2026…mar-2027 y 12 quincenas): " + '; '.join(f"{k[0][0]}{'' if k[3] is None else k[3]} {k[4]}…{k[5]}" + (f' paga {k[6]}' if k[6] else '') for k in filas) + f" · quincena 2 de feb-2027 paga {feb2} · ventana 2026-10-10…2026-11-09 como mes: {ventana} · equivalencia con P1–P6 (efectivo y devengado totales): {per_ok} de 12 · eventos de la semilla entre los días 1 y 9: {len(fechas_1_9)}")
# ---------- S1-10
kinds = re.search(r"kind TEXT NOT NULL CHECK \(kind IN \(([^)]*)\)\), -- revolving_line", up)
tasas = x('SELECT term_count, rate_num, rate_den FROM credit_term_rate ORDER BY term_count').fetchall()
pct = {t: (Fraction(nm, dn) * 100) for t, nm, dn in tasas}
r2 = round(float(pct[2]), 2); r4 = round(float(pct[4]), 4)
promo_cols = [k for k in ('merchant', 'term_count', 'min_purchase_cents', 'start_date', 'end_date') if k in cols['promotion']]
out('S1-10', lect == 'B' and r2 == 17.27 and r4 == 25.8685 and has('credit.availability_confirmed'), f"credit.kind = {kinds.group(1) if kinds else '?'} · tasas semilla como fracción: 2q 53450/309450 = {r2} %, 4q 80050/309450 = {r4} % · lectura por defecto al insertar sin valor = {lect} · availability_confirmed={has('credit.availability_confirmed')} · promotion: {', '.join(cols['promotion'][1:8])}")
# ---------- S1-11
lista = [('tiempos 4', 'meal_slot.code'), ('fuera sin receta', 'meal_plan_entry.is_eating_out'), ('porciones base 2 y escalado', 'recipe.base_servings'), ('perfil por persona', 'health_profile.height_mm'), ('historial de peso', 'weight_entry.weight_g'), ('restricciones vacías', 'health_profile.dietary_restrictions'), ('kcal con fuente', 'nutrition_fact.source'), ('inventario híbrido y sobrantes', 'inventory_item.tracking_mode'), ('cadena/sucursal/mercado con día y qué se compra', 'store.visit_weekday'), ('membresía Sam’s', 'store.membership_annual_fee_cents'), ('precio con tienda, fecha y fuente', 'price.observed_date'), ('unidad base y macros previstos', 'nutrition_fact.protein_mg_per_100g')]
mal = [tc for _, tc in lista if not has(tc)]
out('S1-11', not mal and one("SELECT count(*) FROM meal_slot") == 4, f"{len(lista) - len(mal)} de 12 puntos con tabla.columna: " + '; '.join(f'{a} → {b}' for a, b in lista))
# ---------- S1-13
sec = DIC.split('## Ajustes: 14 secciones')[1].split('\n## ')[0]
filas13 = re.findall(r'^\| (\d+) \| ([^|]+) \| ([^|]+) \| ([^|]+) \|$', sec, re.M)
NOMBRES = ['Cuenta y roles', 'Presupuesto', 'Ingresos', 'Categorías', 'Créditos y condiciones', 'Tiendas y precios', 'Recetas', 'Reportes y frecuencias', 'Notificaciones', 'Apariencia', 'Datos y respaldos', 'Personas y reparto', 'Salud y metas de ahorro', 'Inventario']
okn = [f[1].startswith(nm) for f, nm in zip(filas13, NOMBRES)]
claves = {r[0] for r in x('SELECT setting_key FROM app_setting')} | set(re.findall(r'^\| `([a-z_.]+)` \|', DIC.split('## Claves de app_setting previstas sin valor')[1].split('\n## ')[0], re.M))
refs = set(re.findall(r'`([a-z_]+\.[a-z_]+)`', sec)) - claves; malos = sorted(r for r in refs if r.split('.')[0] in cols and not has(r))
PSEC = [(l.split(' | ')[0].strip(), l.split(' | ')[5].strip()) for l in re.findall(r'^-- @param (.+)$', mig, re.M)]
mal13 = [k_ for k_, s_ in PSEC if not any(f[1].startswith(s_) and f'`{k_}`' in f[3] for f in filas13)]
tablas13 = sorted({r.split('.')[0] for r in refs if r.split('.')[0] in cols})
out('S1-13', len(filas13) == 14 and all(okn) and all(f[3].strip() for f in filas13) and not malos and PSEC and not mal13, f"diccionario «Ajustes: 14 secciones → campos → tabla.columna»: {len(filas13)} de 14 secciones en orden, 0 destinos vacíos={all(f[3].strip() for f in filas13)}, referencias tabla.columna={len(refs)} (inexistentes: {malos}), tablas citadas={len(tablas13)}; Apariencia local por dispositivo (N-301) · claves @param en la fila de su sección: {len(PSEC) - len(mal13)} de {len(PSEC)} (faltan: {mal13})")
# ---------- S1-14
sens = sorted(set(re.findall(r'\b\w*(?:password|token|secret|apikey)\w*\b', mig, re.I)))
tg = re.findall(r'[0-9]{8,}:AA', mig)
out('S1-14', not tg and 'password_hash' in sens, f"coincidencias de password/token/secret/apikey en la migración (solo nombres de columna o texto): {sens} · patrón de token de Telegram: {len(tg)} · password_hash CHECK argon2id: {'$argon2id$' in mig}")
# ---------- S1-15
fkm = {t: len(x(f'PRAGMA foreign_key_list({t})').fetchall()) for t in ('expense', 'credit_purchase', 'scheduled_payment_share', 'price', 'meal_plan_entry', 'shopping_list_item')}
idx = {r[0]: r[1] for r in x("SELECT name, sql FROM sqlite_master WHERE type='index' AND sql IS NOT NULL")}
def ix_de(t, col):
    return [il[1] for il in x(f'PRAGMA index_list({t})') if x(f'PRAGMA index_info({il[1]})').fetchone()[2] == col]
filtros = {'persona': ('expense', 'person_id'), 'familia (partes por persona)': ('expense_share', 'person_id'), 'categoría': ('expense', 'category_id'), 'tienda': ('expense', 'store_id'), 'periodo (fechas)': ('expense', 'expense_date'), 'crédito': ('credit_purchase', 'credit_id')}
fi = {k: ix_de(*v) for k, v in filtros.items()}
out('S1-15', all(fi.values()) and all(fkm.values()), f"foreign_key_list: {fkm} · índices de filtros C1: " + ', '.join(f'{k}={v[0] if v else None}' for k, v in fi.items()) + " · PRAGMA foreign_keys documentado en la cabecera de la migración y en C-13")
# ---------- S1-16
def bloques(t, pat):
    return re.split(pat, t)
prov = []
for i_, t_ in enumerate(TXT):
    for b in re.split(r'\n(?=-- @table )', t_):
        prov += [('migración', l, b) for l in b.splitlines() if 'PROVISIONAL' in l and not historia(i_, l, b)]
for f, t, pat in (('diccionario', DIC, r'\n(?=### `|## )'), ('reglas', REG, r'\n(?=## |- \*\*C-)')):
    for b in bloques(t, pat):
        prov += [(f, l, b) for l in b.splitlines() if 'PROVISIONAL' in l]
blq = set(os.listdir(R + '/docs/estado/bloqueos'))
sin_ref = []; refs16 = set()
for f, l, b in prov:
    mrn = re.match(r'\| (RN-\d\d) \|', l)
    if mrn: b = REG.split('## ' + mrn.group(1))[1].split('\n## ')[0]
    nombres = set(re.findall(r'\b((?:S1|C1|U1|U2|S6|criterios)-[a-z0-9-]+[a-z0-9])\b', l)) or set(re.findall(r'\b((?:S1|C1|U1|U2|S6|criterios)-[a-z0-9-]+[a-z0-9])\b', b))
    nombres = {nm for nm in nombres if nm + '.md' in blq and 'Respuesta de Axel' not in open(R + '/docs/estado/bloqueos/' + nm + '.md', encoding='utf-8').read()}
    refs16 |= nombres
    if not nombres and not re.search(r'PROVISIONAL\*\* marca|PROVISIONAL\*\* marca|marca lo que depende|PROVISIONAL\b.*bloqueo abierto', l): sin_ref.append((f, l[:90]))
out('S1-16', not sin_ref, f"líneas con PROVISIONAL: {len(prov)} (migración {sum(1 for f, _, _ in prov if f == 'migración')}, diccionario {sum(1 for f, _, _ in prov if f == 'diccionario')}, reglas {sum(1 for f, _, _ in prov if f == 'reglas')}) · bloqueos abiertos citados (sin «Respuesta de Axel»): {sorted(refs16)} · líneas sin bloqueo en su bloque={len(sin_ref)} {sin_ref}")
# ---------- S1-17
ident = [t for t in tabs] + [f'{t}.{cl}' for t in tabs for cl in cols[t]]
malid = [i for i in ident if not re.fullmatch(r'[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)?', i)]
glos = re.findall(r'^\| [^|]+ \| (.+) \|$', DIC.split('## Glosario de nombre_codigo')[1].split('\n## ')[0], re.M)[1:]
out('S1-17', not malid and len(tabs) == len(set(tabs)), f"identificadores (tablas y columnas): {len(ident)}, fuera de [a-z0-9_]: {len(malid)} · tablas únicas: {len(set(tabs))} · filas del glosario: {len(glos) - 1}, repetidas: {len(glos) - len(set(glos))}")
# ---------- S1-03
secs = re.split(r'\n(?=## RN-\d\d)', REG)
def path_get(root, p):
    objs = [root]
    for tok in re.findall(r'[^.\[\]]+|\[\d+(?:\.\.\d+)?\]', p):
        if tok.startswith('['):
            a, _, b = tok[1:-1].partition('..')
            objs = [o[i] for o in objs for i in range(int(a), int(b or a) + 1)]
        else:
            objs = [o[tok] for o in objs]
    return objs if '..' in p else objs[0]
def fmt(v): return f'{v:,.2f}'
FIG2 = re.compile(r'(?<![\d,.])\d{1,3}(?:,\d{3})*\.\d{2}(?![\d])')  # cifra con dos decimales (pesos, pesos de frescura)
ORA = set()  # cifras del oráculo: hojas numéricas de model.json y data.json y cifras escritas dentro de sus textos
def hojas(o):
    if isinstance(o, dict):
        for w in o.values(): hojas(w)
    elif isinstance(o, list):
        for w in o: hojas(w)
    elif isinstance(o, bool): pass
    elif isinstance(o, (int, float)): ORA.add(Fraction(str(o)))
    elif isinstance(o, str): ORA.update(Fraction(f.replace(',', '')) for f in FIG2.findall(o))
hojas(m); hojas(dj)
# Fechas (revisión S1-r2, R2-01): ISO (2026-10-10), dd-mmm-aaaa (10-oct-2026) y dd-mmm sin año (15-oct). El año de una
# fecha sin año sale de la fecha completa o del mes-año (oct-2026) más cercano a su derecha en la misma frase (un año
# menos si su mes es mayor); si no hay, del más cercano a su izquierda o del último del ejemplo (un año más si su mes es
# menor). Si el ejemplo no da ningún año, la fecha queda «????» y cuenta como fuera del oráculo.
MESES = {'ene': 1, 'feb': 2, 'mar': 3, 'abr': 4, 'may': 5, 'jun': 6, 'jul': 7, 'ago': 8, 'sep': 9, 'oct': 10, 'nov': 11, 'dic': 12}
FISO = re.compile(r'(?<![\d-])(\d{4})-(\d{2})-(\d{2})(?!\d)')
FDMA = re.compile(r'(?<![\w-])(\d{1,2})-(' + '|'.join(MESES) + r')-(\d{4})(?!\d)')
FDM = re.compile(r'(?<![\w-])(\d{1,2})-(' + '|'.join(MESES) + r')(?![\w-])')
FMA = re.compile(r'(?<![\w-])(' + '|'.join(MESES) + r')-(\d{4})(?!\d)')  # mes-año: solo sirve para deducir el año
def iso(y, mo, d):  # fecha ISO; una fecha imposible (31-nov) queda marcada con «!»
    try: return datetime.date(int(y), int(mo), int(d)).isoformat()
    except ValueError: return f'{int(y)}-{int(mo):02d}-{int(d):02d}!'
def completas(t):  # [(posición, fecha ISO)] de las fechas con año de un texto
    return sorted([(a.start(), iso(*a.groups())) for a in FISO.finditer(t)] + [(a.start(), iso(a.group(3), MESES[a.group(2)], a.group(1))) for a in FDMA.finditer(t)])
ORA_F = set()  # fechas del oráculo: las que aparecen en valores, llaves o textos de model.json y data.json
def hojas_f(o):
    if isinstance(o, dict):
        for k, w in o.items(): ORA_F.update(f for _, f in completas(str(k))); hojas_f(w)
    elif isinstance(o, list):
        for w in o: hojas_f(w)
    elif isinstance(o, str): ORA_F.update(f for _, f in completas(o))
hojas_f(m); hojas_f(dj)
def fechas(q, ult):  # fechas de una frase en ISO (las sin año, con el año deducido) y el último (año, mes) visto
    comp = completas(q); res = [f for _, f in comp]
    ctx = sorted([(p, int(f[:4]), int(f[5:7])) for p, f in comp] + [(a.start(), int(a.group(2)), MESES[a.group(1)]) for a in FMA.finditer(q)])
    for a in FDM.finditer(q):
        d, mo = int(a.group(1)), MESES[a.group(2)]
        der = [(y, mm) for p, y, mm in ctx if p > a.start()]; izq = [(y, mm) for p, y, mm in ctx if p < a.start()] or ([ult] if ult else [])
        if der: y = der[0][0] - (mo > der[0][1])
        elif izq: y = izq[-1][0] + (mo < izq[-1][1])
        else: res.append(f'????-{mo:02d}-{d:02d}'); continue
        res.append(iso(y, mo, d))
    return res, (ctx[-1][1:] if ctx else ult)
def fuera_f(fs): return [f for f in fs if f not in ORA_F]  # sin año, imposibles o ausentes del oráculo
OPS = {'/': lambda a, b: a / b, '×': lambda a, b: a * b, '*': lambda a, b: a * b}
def half_up(q):  # Fraction -> Fraction a centavos con ROUND_HALF_UP
    k = int(abs(q) * 100 + Fraction(1, 2))
    return Fraction(k if q >= 0 else -k, 100)
def frac(t): return Fraction(t.strip().replace(',', ''))
tabla3 = []; malas3 = []; fuera3 = []; ilus_ent = []; ilus_frag = []; fuera3f = []; imposibles3 = []
n_ley = 0; n_cuentas = 0; n_cif_ilus = 0; n_fec = 0; n_fec_ilus = 0
for s in secs:
    mm = re.match(r'## (RN-\d\d)', s)
    if not mm: continue
    rn = mm.group(1)
    citas = re.findall(r'`(model|data)\.json: ([^`=]+?)(?: = ([^`]+))?`', s)
    okc = 0; tot = 0; der = set()
    for src, pth, val in citas:
        pth = pth.strip()
        if pth == 'ruta': continue
        tot += 1
        op = re.fullmatch(r'(.+?)((?:\s*[/×*]\s*\d+)+)', pth)  # cuenta citada: «ruta / n = valor» o encadenada «ruta × 7 / 31 = valor», a centavos con HALF_UP al final
        if op: pth = op.group(1); pasos = re.findall(r'([/×*])\s*(\d+)', op.group(2))
        try:
            v = path_get(m if src == 'model' else dj, pth)
        except (KeyError, IndexError, TypeError):
            malas3.append((rn, pth, 'ruta inexistente')); continue
        if op:
            good = bool(val) and isinstance(v, (int, float)) and re.fullmatch(r'-?[\d,]+(\.\d+)?', val.strip()) is not None and half_up(functools.reduce(lambda a, st: OPS[st[0]](a, Fraction(st[1])), pasos, Fraction(str(v)))) == frac(val)
            if good: der.add(frac(val)); n_cuentas += 1
        elif val:
            good = str(v) == val.strip() or (isinstance(v, (int, float)) and re.fullmatch(r'-?[\d,]+(\.\d+)?', val.strip()) and Fraction(str(v)) == frac(val))
        else:
            vs = v if isinstance(v, list) and '..' in pth else [v]
            good = all((not isinstance(w, (int, float))) or fmt(w) in s or str(w) in s or f'{w * 100:g}' in s for w in vs)
        okc += bool(good)
        if not good: malas3.append((rn, pth, v))
    i = s.find('**Ejemplo')
    ejt = s[i:].split('**Esquema.**')[0] if i >= 0 else ''
    entero = ejt.startswith('**Ejemplo ilustrativo')  # todo el ejemplo es ilustrativo
    frases = [q for q in re.split(r'(?<=[.;])\s+', ejt) if q.strip()]  # el rótulo de una frase solo cubre esa frase
    rot = [q for q in frases if 'ilustrativ' in q.lower()]
    n_ley += ejt.lower().count('ilustrativo (sin dato en el oráculo)')
    sin_rot = []; sin_rot_f = []; ult = None
    for q in frases:
        nuevas = [f for f in FIG2.findall(q) if frac(f) not in ORA and frac(f) not in der]
        fs, ult = fechas(q, ult); n_fec += len(fs); nuevas_f = fuera_f(fs)
        imposibles3 += [(rn, f) for f in fs if f.endswith('!')]  # una fecha imposible falla aunque esté rotulada
        if entero or 'ilustrativ' in q.lower(): n_cif_ilus += len(nuevas); n_fec_ilus += len(nuevas_f)
        else: sin_rot += nuevas; sin_rot_f += nuevas_f
    if sin_rot: fuera3.append((rn, sin_rot))
    if sin_rot_f: fuera3f.append((rn, sin_rot_f))
    if entero: ilus_ent.append(rn)
    elif rot: ilus_frag.append((rn, len(rot)))
    marca = ' (ilustrativo)' if entero else (f' ({len(rot)} frases ilustrativas)' if rot else '')
    tabla3.append(f"{rn}: ejemplo={'sí' if i >= 0 else 'NO'}{marca}, citas {okc}/{tot}" + (f', cuentas citadas {len(der)}' if der else '') + (f', CIFRAS FUERA DEL ORÁCULO SIN RÓTULO {sin_rot}' if sin_rot else '') + (f', FECHAS FUERA DEL ORÁCULO SIN RÓTULO {sin_rot_f}' if sin_rot_f else ''))
# Autoprueba de las fechas: las frases de RN-06 y RN-11 que señaló la revisión S1-r2 (sin rótulo) deben dar exactamente
# sus fechas fuera del oráculo; la última prueba el año deducido al cruzar de diciembre a enero y una fecha imposible.
AUTO_F = [('El sábado 10-oct-2026 cae en la semana del lunes 5-oct-2026 al domingo 11-oct-2026, y el 31-oct-2026 en la quincena 2 de octubre, que paga el 30-oct-2026.', ['2026-10-05', '2026-10-11', '2026-10-31']),
          ('moverla al 16-dic-2026 la cumple.', ['2026-12-16']),
          ('Compra del 10-oct-2026: el 15-oct queda a 5 días (no cuenta) y el 30-oct a 20 días, así que los pagos son 30-oct, 15-nov, 30-nov y 15-dic-2026.', ['2026-10-15']),
          ('Con «deuda en cero antes del 1-may-2027», la despensa 4 se descarta porque paga hasta el 15-jun-2027.', ['2027-05-01']),
          ('Pagos del 30-dic y del 15-ene-2027; el 31-nov-2026 no existe y el 2026-10-12 no está.', ['2026-10-12', '2026-11-31!'])]
auto_f = [(q[:30], sorted(fuera_f(fechas(q, None)[0])), sorted(e)) for q, e in AUTO_F]
auto_ok = all(a == b for _, a, b in auto_f)
sin_ej = [t for t in tabla3 if 'ejemplo=NO' in t]
frag_txt = ', '.join(f'{a} ({b})' for a, b in ilus_frag)
out('S1-03', not sin_ej and not malas3 and not fuera3 and not fuera3f and not imposibles3 and auto_ok,
    f"{len(tabla3)} reglas · " + ' · '.join(tabla3)
    + f" · reglas sin ejemplo={len(sin_ej)} · citas que no coinciden con su ruta={malas3}"
    + f" · ejemplos con cifras de dos decimales que no están en model.json ni en data.json y sin rótulo={len(fuera3)}{(' ' + str(fuera3)) if fuera3 else ''}"
    + f" · ILUSTRATIVOS, aparte: {len(ilus_ent)} ejemplos completos [{', '.join(ilus_ent)}] y {sum(b for _, b in ilus_frag)} frases rotuladas en {len(ilus_frag)} ejemplos con dato [{frag_txt}]; cifras ilustrativas que no están en el oráculo={n_cif_ilus}; rótulos con la leyenda literal «ilustrativo (sin dato en el oráculo)»={n_ley}"
    + f" · cuentas citadas que coinciden={n_cuentas}"
    + f" · FECHAS (ISO, dd-mmm-aaaa y dd-mmm con el año deducido): {n_fec} en los ejemplos, {len(ORA_F)} distintas en el oráculo; ejemplos con fechas que no están en model.json ni en data.json y sin rótulo={len(fuera3f)}{(' ' + str(fuera3f)) if fuera3f else ''}; fechas imposibles={imposibles3}; fechas ilustrativas que no están en el oráculo={n_fec_ilus}"
    + f" · autoprueba de fechas con las frases de la revisión S1-r2: {'pasa' if auto_ok else 'FALLA ' + str([t for t in auto_f if t[1] != t[2]])} ({len(AUTO_F)} casos)")
# ---------- negativas (S1-07, S1-08, S1-05, S1-06 y tablas nuevas)
x("INSERT INTO product(name,base_unit) VALUES ('Arroz','g')"); prod = lid()
x("INSERT INTO recipe(name,origin) VALUES ('Arroz rojo','claude_seed')"); rec = lid()
x("INSERT INTO inventory_item(product_id,tracking_mode,quantity_milli) VALUES (?,'stock',1000)", (prod,)); inv = lid()
x("INSERT INTO shopping_list(start_date,end_date) VALUES ('2026-10-12','2026-10-18')"); sl = lid()
x("INSERT INTO shopping_list_item(shopping_list_id,origin,product_id,needed_milli) VALUES (?,'calculated',?,500000)", (sl, prod))
x("INSERT INTO recipe_meal_slot(recipe_id,meal_slot_id) VALUES (?,2)", (rec,))
KAREN, X_, CAT, CR = P['karen'], P['ximena'], cat, cr
NEG = [
 ('karen_con_cuenta', f"INSERT INTO user_account(person_id,username,password_hash) VALUES ({KAREN},'karen','{H}')"),
 ('tercera_cuenta_otra_co_titular', "INSERT INTO person(code,display_name,participation) VALUES ('otra','Otra','co_holder')"),
 ('ximena_pasa_a_externa_con_cuenta', f"UPDATE person SET participation='external' WHERE id={X_}"),
 ('ingreso_de_ximena', f"INSERT INTO income(person_id,kind,description,amount_cents,income_date) VALUES ({X_},'extra','x',100,'2026-10-15')"),
 ('hash_no_argon2id', f"UPDATE user_account SET password_hash='bcrypt$abc' WHERE person_id={X_}"),
 ('segundo_titular', "INSERT INTO person(code,display_name,participation) VALUES ('otro','Otro','holder')"),
 ('compra_total_distinto', f"INSERT INTO credit_purchase(credit_id,kind,purchase_date,principal_cents,interest_cents,total_cents,term_count,payment_frequency,first_payment_date,interest_source) VALUES ({CR},'installment_plan','2026-10-10',100,10,111,4,'fortnightly','2026-10-30','quoted')"),
 ('pago_distinto', "INSERT INTO scheduled_payment(credit_purchase_id,seq,due_date,amount_cents,principal_cents,interest_cents) VALUES (1,5,'2026-12-30',100,50,49)"),
 ('pct_mayor_100', f"INSERT INTO cost_share_rule_part(rule_id,person_id,pct_bp) VALUES ({r_hogar},{KAREN},10001)"),
 ('parte_pct_y_fijo', f"INSERT INTO cost_share_rule_part(rule_id,person_id,pct_bp,fixed_cents) VALUES ({r_hogar},{KAREN},5000,100)"),
 ('fecha_30_feb', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',{CAT},100,'2027-02-30','cash')"),
 ('fecha_pantalla', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',{CAT},100,'10-oct-2026','cash')"),
 ('quincena_paga_14', "INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES ('fortnight',(SELECT id FROM period WHERE kind='month' AND year=2027 AND month=2),2027,2,1,'2027-02-01','2027-02-15','2027-02-14')"),
 ('quincena2_feb_paga_30', "INSERT INTO period(kind,parent_id,year,month,half,start_date,end_date,pay_date) VALUES ('fortnight',(SELECT id FROM period WHERE kind='month' AND year=2027 AND month=2),2027,2,2,'2027-02-16','2027-02-28','2027-02-30')"),
 ('real_en_cents', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',{CAT},12.5,'2026-10-10','cash')"),
 ('fk_inexistente', "INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',999,100,'2026-10-10','cash')"),
 ('plan_comida_de_karen', f"INSERT INTO meal_plan_entry(plan_date,meal_slot_id,person_id,is_eating_out) VALUES ('2026-10-12',2,{KAREN},1)"),
 ('fuera_con_porciones', f"INSERT INTO meal_plan_entry(plan_date,meal_slot_id,person_id,servings_milli,is_eating_out) VALUES ('2026-10-12',2,{P['axel']},1000,0)"),
 ('lectura_A_sin_ancla', f"INSERT INTO credit(name,kind,holder_person_id,credit_limit_cents,payment_frequency,availability_reading) VALUES ('Otra','revolving_line',{P['axel']},500000,'fortnightly','A')"),
 ('tarjeta_sin_corte', f"INSERT INTO credit(name,kind,holder_person_id,credit_limit_cents,payment_frequency) VALUES ('Tarjeta X','bank_card',{P['axel']},500000,'monthly')"),
 ('mercado_sin_dia', "INSERT INTO store(kind,name) VALUES ('market','Tianguis sin día')"),
 ('credito_con_pagador_de_contado', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method,paid_by_person_id) VALUES ('variable',{CAT},100,'2026-10-10','credit',{KAREN})"),
 ('borrar_karen_con_historia', f"DELETE FROM person WHERE id={KAREN}"),
 ('bitacora_inmutable', "UPDATE audit_log SET table_name='x'"),
 ('descripcion_en_blanco', f"INSERT INTO expense(kind,description,category_id,amount_cents,expense_date,payment_method) VALUES ('variable','  ',{CAT},100,'2026-10-10','cash')"),
 ('categoria_tipo_invalido', "INSERT INTO category(name,kind) VALUES ('X','upcoming')"),
 ('presupuesto_no_inicia_dia_1', f"INSERT INTO budget(category_id,amount_cents,valid_from_date) VALUES ({CAT},100,'2026-10-10')"),
 ('presupuesto_fin_no_es_fin_de_mes', f"INSERT INTO budget(category_id,amount_cents,valid_from_date,valid_to_date) VALUES ({CAT},100,'2026-10-01','2026-11-15')"),
 ('presupuesto_fin_antes_de_inicio', f"INSERT INTO budget(category_id,amount_cents,valid_from_date,valid_to_date) VALUES ({CAT},100,'2026-11-01','2026-10-31')"),
 ('presupuesto_mismo_inicio_dos_veces', f"INSERT INTO budget(category_id,amount_cents,valid_from_date) VALUES ({CAT},100,'2026-12-01'); INSERT INTO budget(category_id,amount_cents,valid_from_date) VALUES ({CAT},200,'2026-12-01')"),
 ('segunda_quincena_en_mensual', f"INSERT INTO planned_item(kind,description,amount_cents,second_half_amount_cents,frequency,start_date,category_id) VALUES ('fixed_expense','Renta',100,200,'monthly','2026-10-01',{CAT})"),
 ('receta_tiempo_repetido', f"INSERT INTO recipe_meal_slot(recipe_id,meal_slot_id) VALUES ({rec},2)"),
 ('ajuste_sin_motivo', f"INSERT INTO inventory_movement(inventory_item_id,movement_date,kind,quantity_delta_milli) VALUES ({inv},'2026-10-12','adjustment',-100)"),
 ('compra_con_motivo', f"INSERT INTO inventory_movement(inventory_item_id,movement_date,kind,reason,quantity_delta_milli) VALUES ({inv},'2026-10-12','purchase','count',100)"),
 ('motivo_invalido', f"INSERT INTO inventory_movement(inventory_item_id,movement_date,kind,reason,quantity_delta_milli) VALUES ({inv},'2026-10-12','adjustment','robo',-100)"),
 ('lista_no_dura_7_dias', "INSERT INTO shopping_list(start_date,end_date) VALUES ('2026-10-19','2026-10-24')"),
 ('lista_semana_repetida', "INSERT INTO shopping_list(start_date,end_date) VALUES ('2026-10-12','2026-10-18')"),
 ('articulo_producto_y_texto', f"INSERT INTO shopping_list_item(shopping_list_id,origin,product_id,free_text,free_text_unit,quantity_milli) VALUES ({sl},'manual',{prod},'Tortillas','g',1000)"),
 ('texto_libre_sin_unidad', f"INSERT INTO shopping_list_item(shopping_list_id,origin,free_text,quantity_milli) VALUES ({sl},'manual','Tortillas',1000)"),
 ('manual_sin_cantidad', f"INSERT INTO shopping_list_item(shopping_list_id,origin,free_text,free_text_unit) VALUES ({sl},'manual','Tortillas','g')"),
 ('calculado_sin_necesidad', f"INSERT INTO shopping_list_item(shopping_list_id,origin,product_id) VALUES ({sl},'calculated',{prod})"),
 ('calculado_repetido', f"INSERT INTO shopping_list_item(shopping_list_id,origin,product_id,needed_milli) VALUES ({sl},'calculated',{prod},1000)"),
 ('hora_de_compra_sin_comprado', f"INSERT INTO shopping_list_item(shopping_list_id,origin,free_text,free_text_unit,quantity_milli,bought_at) VALUES ({sl},'manual','Pan','piece',2000,'2026-10-12T18:00:00.000Z')"),
 ('clave_idempotencia_corta', "INSERT INTO idempotency_key(user_account_id,idempotency_key,http_method,request_path,request_sha256) VALUES (1,'corta','POST','/api/v1/expenses','" + '0' * 64 + "')"),
 ('clave_idempotencia_con_espacio', "INSERT INTO idempotency_key(user_account_id,idempotency_key,http_method,request_path,request_sha256) VALUES (1,'clave con espacio 123','POST','/api/v1/expenses','" + '0' * 64 + "')"),
 ('clave_idempotencia_repetida', "INSERT INTO idempotency_key(user_account_id,idempotency_key,http_method,request_path,request_sha256) VALUES (1,'7f3c2a10-aaaa-4bbb-8ccc-0123456789ab','POST','/api/v1/expenses','" + '1' * 64 + "'); INSERT INTO idempotency_key(user_account_id,idempotency_key,http_method,request_path,request_sha256) VALUES (1,'7f3c2a10-aaaa-4bbb-8ccc-0123456789ab','POST','/api/v1/expenses','" + '2' * 64 + "')"),
]
x("INSERT INTO audit_log(action,table_name,row_id,after_json) VALUES ('insert','person',1,'{\"id\":1}')"); db.commit()
aceptadas = []
for nombre, sql in NEG:
    db.execute('SAVEPOINT neg')
    try:
        db.executescript('SAVEPOINT n2; ' + sql + '; RELEASE n2;') if False else [db.execute(q) for q in sql.split('; ')]
        aceptadas.append(nombre)
    except (sqlite3.IntegrityError, sqlite3.DatabaseError):
        pass
    db.execute('ROLLBACK TO neg'); db.execute('RELEASE neg')
POS = [
 ('gasto_sin_descripcion', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',{CAT},100,'2026-10-12','cash')"),
 ('ingreso_15_y_30_distintos', f"INSERT INTO planned_item(kind,description,amount_cents,second_half_amount_cents,frequency,start_date,person_id) VALUES ('income','Sueldo',100000,120000,'fortnightly','2026-10-15',{P['axel']})"),
 ('presupuesto_un_mes_y_vigente', f"INSERT INTO budget(category_id,amount_cents,valid_from_date,valid_to_date) VALUES ({CAT},200000,'2026-10-01','2026-10-31'); INSERT INTO budget(category_id,amount_cents,valid_from_date) VALUES ({CAT},210000,'2026-11-01')"),
 ('receta_en_dos_tiempos', f"INSERT INTO recipe_meal_slot(recipe_id,meal_slot_id) VALUES ({rec},3)"),
 ('ajuste_por_merma', f"INSERT INTO inventory_movement(inventory_item_id,movement_date,kind,reason,quantity_delta_milli) VALUES ({inv},'2026-10-12','adjustment','waste',-250)"),
 ('articulo_a_mano_comprado', f"INSERT INTO shopping_list_item(shopping_list_id,origin,free_text,free_text_unit,quantity_milli,store_id,is_bought,bought_at) VALUES ({sl},'manual','Tortillas','g',1000000,NULL,1,'2026-10-12T18:00:00.000Z')"),
 ('clave_idempotencia_valida', "INSERT INTO idempotency_key(user_account_id,idempotency_key,http_method,request_path,request_sha256,response_status,response_body) VALUES (1,'7f3c2a10-aaaa-4bbb-8ccc-0123456789ab','POST','/api/v1/expenses','" + '3' * 64 + "',201,'{\"id\":1}')"),
 ('regla_60_mas_50_documentada_invariante_8', f"INSERT INTO cost_share_rule(name,remainder_person_id) VALUES ('Prueba 60/50',{P['axel']}); INSERT INTO cost_share_rule_part(rule_id,person_id,pct_bp) VALUES (last_insert_rowid(),{KAREN},6000)"),
]
rechazadas_pos = []
for nombre, sql in POS:
    db.execute('SAVEPOINT pos')
    try:
        [db.execute(q) for q in sql.split('; ')]
    except sqlite3.DatabaseError as e:
        rechazadas_pos.append((nombre, str(e)))
    db.execute('ROLLBACK TO pos'); db.execute('RELEASE pos')
neg_ok = not aceptadas and not rechazadas_pos
out('S1-07', 'karen_con_cuenta' not in aceptadas and 'tercera_cuenta_otra_co_titular' not in aceptadas and 'ximena_pasa_a_externa_con_cuenta' not in aceptadas and 'ingreso_de_ximena' not in aceptadas and 'role' not in ' '.join(cols['user_account']),
    f"cuenta para Karen: rechazada · tercera cuenta (otra co-titular): rechazada · Ximena con cuenta pasa a externa: rechazada · ingreso de Ximena: rechazado · columnas de user_account sin rol: {cols['user_account']} · tablas de finanzas de la familia de Karen: {[t for t in tabs if 'karen' in t or 'family' in t]}")
ax_d = c(dj['cuentas']['conjunto']['mitad']) + sum(c(l.get('total', l.get('precio', 0))) for l in []) 
out('S1-08', 'regla_60_mas_50_documentada_invariante_8' not in [p for p, _ in rechazadas_pos], f"regla con partes 60 % + 50 %: la base la acepta y la rechaza la API (invariante 8 del diccionario, prueba de S6) · despensa 1 sin pérdida: partes del gasto Axel {c(m['despensas'][0]['axel']['capital'])} + Karen {c(m['despensas'][0]['karen']['capital'])} = {c(m['despensas'][0]['capital'])} (ver S1-12) · conjunto {c(dj['cuentas']['conjunto']['total'])} → mitad {ax_d} cada quien (data.json: cuentas.conjunto.mitad)")
out('NEG', neg_ok, f"pruebas negativas: {len(NEG)}, rechazadas {len(NEG) - len(aceptadas)}, aceptadas por error {aceptadas} · positivas de las tablas nuevas: {len(POS)}, aceptadas {len(POS) - len(rechazadas_pos)}, rechazadas por error {rechazadas_pos}")
# ---------- columnas y tabla nuevas de 1.2 (K3, K4, K9, K12, K18, K23): pruebas negativas y positivas
AX = P['axel']
x('INSERT INTO health_profile(person_id) VALUES (?)', (AX,))
x("INSERT INTO weight_entry(person_id,measured_date,weight_g,idempotency_key,captured_at) VALUES (?,'2026-10-10',72500,'cola-peso-0001-abcdef','2026-10-10T19:20:00.000Z')", (AX,))
x("INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method,idempotency_key,captured_at) VALUES ('variable',?,100,'2026-10-10','cash','cola-gasto-0001-abcdef','2026-10-10T19:20:00.000Z')", (CAT,))
x("INSERT INTO attachment(kind,file_path,mime_type,size_bytes,sha256,ocr_status) VALUES ('ticket_photo','tickets/2026-10-10.jpg','image/jpeg',1000,?,'done')", ('a' * 64,)); ATT = lid()
x("INSERT INTO telegram_chat_candidate(chat_id,telegram_name) VALUES ('123456789','Ximena')"); TCC = lid()
db.commit()
T18 = "INSERT INTO credit(name,kind,holder_person_id,credit_limit_cents,payment_frequency,statement_day,payment_due_day,"
NEG12 = [
 ('lista_empieza_martes', "INSERT INTO shopping_list(start_date,end_date) VALUES ('2026-10-13','2026-10-19')"),
 ('lista_cambia_a_miercoles', f"UPDATE shopping_list SET start_date='2026-10-14', end_date='2026-10-20' WHERE id={sl}"),
 ('gasto_clave_repetida', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method,idempotency_key) VALUES ('variable',{CAT},100,'2026-10-11','cash','cola-gasto-0001-abcdef')"),
 ('peso_clave_repetida', f"INSERT INTO weight_entry(person_id,measured_date,weight_g,idempotency_key) VALUES ({AX},'2026-10-11',72400,'cola-peso-0001-abcdef')"),
 ('gasto_clave_corta', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method,idempotency_key) VALUES ('variable',{CAT},100,'2026-10-11','cash','corta')"),
 ('peso_clave_con_espacio', f"INSERT INTO weight_entry(person_id,measured_date,weight_g,idempotency_key) VALUES ({AX},'2026-10-11',72400,'clave con espacio 123')"),
 ('gasto_captura_sin_z', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method,captured_at) VALUES ('variable',{CAT},100,'2026-10-11','cash','2026-10-11 10:00:00')"),
 ('peso_captura_formato_pantalla', f"INSERT INTO weight_entry(person_id,measured_date,weight_g,captured_at) VALUES ({AX},'2026-10-12',72400,'11-oct-2026')"),
 ('objetivo_fuera_de_enumerado', f"UPDATE health_profile SET goal='bulk' WHERE person_id={AX}"),
 ('base_de_tasa_invalida', T18 + f"rate_base) VALUES ('Tarjeta A','bank_card',{AX},500000,'monthly',5,25,'weekly')"),
 ('base_de_tasa_en_linea_revolvente', f"UPDATE credit SET rate_base='daily' WHERE id={CR}"),
 ('pago_minimo_mayor_a_100', T18 + f"min_payment_pct_bp) VALUES ('Tarjeta B','bank_card',{AX},500000,'monthly',5,25,10001)"),
 ('minimo_fijo_negativo', T18 + f"min_payment_fixed_cents) VALUES ('Tarjeta C','bank_card',{AX},500000,'monthly',5,25,-1)"),
 ('decimales_de_tasa_anual_sin_tasa', T18 + f"annual_rate_decimals) VALUES ('Tarjeta D','bank_card',{AX},500000,'monthly',5,25,1)"),
 ('decimales_negativos', f"UPDATE credit_term_rate SET rate_decimals=-1 WHERE credit_id={CR}"),
 ('decimales_de_compra_sin_tasa', "UPDATE credit_purchase SET rate_decimals=2 WHERE rate_num IS NULL"),
 ('borrador_que_no_es_objeto', f"UPDATE attachment SET ocr_draft_json=json_array(1,2), ocr_draft_saved_at='2026-10-10T19:20:00.000Z' WHERE id={ATT}"),
 ('borrador_json_invalido', f"UPDATE attachment SET ocr_draft_json='x', ocr_draft_saved_at='2026-10-10T19:20:00.000Z' WHERE id={ATT}"),
 ('borrador_sin_fecha', f"UPDATE attachment SET ocr_draft_json=json_object() WHERE id={ATT}"),
 ('fecha_de_borrador_sin_borrador', f"UPDATE attachment SET ocr_draft_saved_at='2026-10-10T19:20:00.000Z' WHERE id={ATT}"),
 ('chat_confirmado_sin_persona', f"UPDATE telegram_chat_candidate SET confirmed_at='2026-10-11T01:00:00.000Z' WHERE id={TCC}"),
 ('chat_con_persona_sin_confirmar', f"UPDATE telegram_chat_candidate SET person_id={X_} WHERE id={TCC}"),
 ('chat_repetido', "INSERT INTO telegram_chat_candidate(chat_id,telegram_name) VALUES ('123456789','Otra')"),
 ('chat_invalido', "INSERT INTO telegram_chat_candidate(chat_id,telegram_name) VALUES ('abc','X')"),
 ('chat_sin_nombre', "INSERT INTO telegram_chat_candidate(chat_id,telegram_name) VALUES ('987654321','  ')"),
 ('ajuste_de_1_2_repetido', "INSERT INTO app_setting(setting_key,value_json) VALUES ('projection.repeat_cadence_enabled','true')"),
]
# 1.4 (correccion-S1-r1, M-01): con la 00004 la tarjeta con porcentaje de IVA lleva vat_included = 0 (Q2-67 opción 1); el CHECK rechaza (NULL, 1600)
VI18 = (',vat_included', ',0') if 'vat_included' in [r_[1] for r_ in db.execute('PRAGMA table_info(credit)')] else ('', '')
POS12 = [
 ('lista_empieza_lunes', "INSERT INTO shopping_list(start_date,end_date) VALUES ('2026-10-26','2026-11-01')"),
 ('lista_cambia_a_otro_lunes', f"UPDATE shopping_list SET start_date='2026-11-02', end_date='2026-11-08' WHERE id={sl}"),
 ('gasto_desde_la_cola', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method,idempotency_key,captured_at) VALUES ('variable',{CAT},250,'2026-10-11','cash','cola-gasto-0002-abcdef','2026-10-11T19:20:00Z')"),
 ('dos_gastos_sin_clave', f"INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',{CAT},100,'2026-10-11','cash'); INSERT INTO expense(kind,category_id,amount_cents,expense_date,payment_method) VALUES ('variable',{CAT},100,'2026-10-11','cash')"),
 ('peso_desde_la_cola', f"INSERT INTO weight_entry(person_id,measured_date,weight_g,idempotency_key,captured_at) VALUES ({AX},'2026-10-12',72300,'cola-peso-0002-abcdef','2026-10-12T14:00:00.000Z')"),
 ('borrador_guardado_y_vaciado_al_confirmar', f"UPDATE attachment SET ocr_draft_json=json_object('lineas', json_array()), ocr_draft_saved_at='2026-10-11T02:00:00.000Z' WHERE id={ATT}; UPDATE attachment SET ocr_draft_json=NULL, ocr_draft_saved_at=NULL, ocr_status='confirmed' WHERE id={ATT}"),
 ('chat_confirmado_con_su_destinatario', f"UPDATE telegram_chat_candidate SET confirmed_at='2026-10-11T01:00:00.000Z', person_id={X_} WHERE id={TCC}; INSERT INTO notification_recipient(person_id,chat_id) VALUES ({X_},'123456789')"),
 ('tarjeta_con_datos_k18', T18 + f"annual_rate_num,annual_rate_den,annual_rate_decimals,rate_base,vat_on_interest_bp,min_payment_pct_bp,min_payment_fixed_cents{VI18[0]}) VALUES ('Tarjeta ilustrativa','bank_card',{AX},500000,'monthly',5,25,455,1000,1,'daily',1600,500,20000{VI18[1]})"),
 ('credito_de_tienda_base_mensual', f"INSERT INTO credit(name,kind,holder_person_id,payment_frequency,rate_base) VALUES ('Crédito de tienda ilustrativo','store_credit',{AX},'monthly','monthly')"),
 ('objetivo_bajar', f"UPDATE health_profile SET goal='lose' WHERE person_id={AX}"),
 ('compra_con_tasa_y_decimales', "UPDATE credit_purchase SET rate_num=80050, rate_den=309450, rate_basis='term_total', rate_decimals=4 WHERE id=1"),
 ('promocion_0_por_ciento', "INSERT INTO promotion(merchant_name,term_count,valid_from,valid_to,rate_decimals) VALUES ('Chedraui',4,'2026-10-30','2026-10-30',0)"),
]
def corre(lista, negativas):
    malas = []
    for nombre, sql in lista:
        db.execute('SAVEPOINT p12')
        try:
            [db.execute(q) for q in sql.split('; ')]
            if negativas: malas.append(nombre)
        except sqlite3.DatabaseError as e:
            if not negativas: malas.append((nombre, str(e)))
        db.execute('ROLLBACK TO p12'); db.execute('RELEASE p12')
    return malas
acept12 = corre(NEG12, True); rech12 = corre(POS12, False)
sembradas = dict(x("SELECT setting_key, value_json FROM app_setting WHERE setting_key IN ('projection.repeat_cadence_enabled', 'session.max_age_days', 'idempotency.retention_hours')").fetchall())
k16 = [k for (k,) in x("SELECT setting_key FROM app_setting WHERE setting_key IN ('optimizer.min_available_buffer_cents', 'optimizer.max_live_purchases', 'optimizer.zero_debt_by_date')")]  # 1.4: optimizer.date_window_days sí se siembra (Q2-75)
out('NEG12', not acept12 and not rech12 and sembradas == {'projection.repeat_cadence_enabled': 'false', 'session.max_age_days': '90', 'idempotency.retention_hours': '24'} and not k16,
    f"pruebas negativas de 1.2: {len(NEG12)}, rechazadas {len(NEG12) - len(acept12)}, aceptadas por error {acept12} · positivas: {len(POS12)}, aceptadas {len(POS12) - len(rech12)}, rechazadas por error {rech12} · ajustes sembrados por la 00002: {sembradas} · claves optimizer.* de restricciones sembradas (K16, deben ser 0): {k16}")
# ---------- K44 y K36 (esquema 1.3): columnas, claves sembradas, negativas, cifras de las reglas y marcas PROVISIONAL
SEM13 = {'health.goal_adjust_lose_kcal': '0', 'health.goal_adjust_gain_kcal': '0', 'price.zero_weight_after_days': '90', 'auth.password_min_length': '12', 'attachment.max_size_mb': '10', 'price.profeco_import_max_gb': '1', 'comparator.timeout_seconds': '5', 'alerts.budget_warning_pct': '80', 'alerts.budget_limit_pct': '100', 'alerts.payment_due_days_before': '[1, 0]'}
PREV = {'auth.max_failed_attempts': '5', 'auth.lockout_window_minutes': '15', 'price.stale_after_days': '30', 'session.max_age_days': '90', 'idempotency.retention_hours': '24'}
casa = dict(x('SELECT setting_key, value_json FROM app_setting WHERE user_account_id IS NULL').fetchall())
anot = {l.split(' | ')[0].strip(): l.split(' | ')[3].strip() for l in re.findall(r'^-- @param (.+)$', mig, re.M)}
cols13 = [cc for cc in ('goal_adjust_lose_kcal', 'goal_adjust_gain_kcal') if has('health_profile.' + cc)]
x('UPDATE health_profile SET goal_adjust_lose_kcal = 300 WHERE person_id = ?', (AX,)); db.commit()  # dato que el Down de la 00003 debe dejar intacto en las columnas de 1.2
NEG13 = [
 ('ajuste_bajar_-1', f"UPDATE health_profile SET goal_adjust_lose_kcal = -1 WHERE person_id = {AX}"),
 ('ajuste_subir_-1', f"UPDATE health_profile SET goal_adjust_gain_kcal = -1 WHERE person_id = {AX}"),
 ('ajuste_bajar_texto', f"UPDATE health_profile SET goal_adjust_lose_kcal = '300 kcal' WHERE person_id = {AX}"),
 ('casa_ajuste_-1', "UPDATE app_setting SET value_json = '-1' WHERE setting_key = 'health.goal_adjust_lose_kcal' AND user_account_id IS NULL"),
 ('adjunto_como_texto', "UPDATE app_setting SET value_json = '\"10\"' WHERE setting_key = 'attachment.max_size_mb' AND user_account_id IS NULL"),
 ('comparador_decimal', "UPDATE app_setting SET value_json = '5.5' WHERE setting_key = 'comparator.timeout_seconds' AND user_account_id IS NULL"),
 ('intentos_booleano', "UPDATE app_setting SET value_json = 'true' WHERE setting_key = 'auth.max_failed_attempts' AND user_account_id IS NULL"),
 ('clave_repetida', "INSERT INTO app_setting(setting_key, value_json) VALUES ('price.zero_weight_after_days', '120')"),
 ('alta_por_cuenta_negativa', "INSERT INTO app_setting(setting_key, user_account_id, value_json) VALUES ('session.max_age_days', 1, '-90')"),
 ('umbral_como_texto', "UPDATE app_setting SET value_json = '\"80\"' WHERE setting_key = 'alerts.budget_warning_pct' AND user_account_id IS NULL"),
 ('tope_negativo', "UPDATE app_setting SET value_json = '-100' WHERE setting_key = 'alerts.budget_limit_pct' AND user_account_id IS NULL"),
 ('dias_antes_no_lista', "UPDATE app_setting SET value_json = '1' WHERE setting_key = 'alerts.payment_due_days_before' AND user_account_id IS NULL"),
 ('dias_antes_negativo', "UPDATE app_setting SET value_json = '[1, -1]' WHERE setting_key = 'alerts.payment_due_days_before' AND user_account_id IS NULL"),
 ('dias_antes_decimal', "UPDATE app_setting SET value_json = '[1.5, 0]' WHERE setting_key = 'alerts.payment_due_days_before' AND user_account_id IS NULL"),
 ('dias_antes_texto_por_cuenta', "INSERT INTO app_setting(setting_key, user_account_id, value_json) VALUES ('alerts.payment_due_days_before', 1, '[\"1\"]')"),
]
POS13 = [
 ('ajuste_bajar_500', f"UPDATE health_profile SET goal_adjust_lose_kcal = 500, goal_adjust_gain_kcal = 0 WHERE person_id = {AX}"),
 ('ajuste_null_valor_casa', f"UPDATE health_profile SET goal_adjust_lose_kcal = NULL, goal_adjust_gain_kcal = NULL WHERE person_id = {AX}"),
 ('cambiar_intentos_a_3', "UPDATE app_setting SET value_json = '3' WHERE setting_key = 'auth.max_failed_attempts' AND user_account_id IS NULL"),
 ('restablecer_peso_cero', "UPDATE app_setting SET value_json = '120' WHERE setting_key = 'price.zero_weight_after_days'; UPDATE app_setting SET value_json = '90' WHERE setting_key = 'price.zero_weight_after_days'"),
 ('otro_ajuste_no_entero', "UPDATE app_setting SET value_json = 'true' WHERE setting_key = 'projection.repeat_cadence_enabled'"),
 ('cambiar_aviso_a_70', "UPDATE app_setting SET value_json = '70' WHERE setting_key = 'alerts.budget_warning_pct' AND user_account_id IS NULL"),
 ('dias_antes_3_1_0', "UPDATE app_setting SET value_json = '[3, 1, 0]' WHERE setting_key = 'alerts.payment_due_days_before' AND user_account_id IS NULL"),
]
acept13 = corre(NEG13, True); rech13 = corre(POS13, False)
fr = lambda q, d: int(q * 10**d + Fraction(1, 2)) / Fraction(10**d)  # HALF_UP a d decimales (valores positivos)
msj = (10 * Fraction('72.5') + Fraction('6.25') * 175 - 5 * 35 + 5) * Fraction('1.55')
k36 = (fr(msj, 0), fr(msj - 500, 0)); k40 = [fr(Fraction(90 - d, 61), 4) for d in (30, 60, 90)]
k39 = (half_up(Fraction(str(m['periodos'][0]['efectivo']['total'])) * 7 / 31), half_up(Fraction(str(m['periodos'][4]['efectivo']['total'])) * 7 / 28))
txt_ok = k36 == (2556, 2056) and k40 == [Fraction('0.9836'), Fraction('0.4918'), 0] and k39 == (Fraction('219.88'), Fraction('629.35')) and all(v in REG for v in ('2,556', '2,056', '0.9836', '0.4918', '219.88', '629.35', '(90 − edad) / 61'))
resueltos = sorted(f[:-3] for f in os.listdir(R + '/docs/estado/bloqueos') if f.endswith('.md') and 'Respuesta de Axel' in open(R + '/docs/estado/bloqueos/' + f, encoding='utf-8').read())
prov_res = [(f, nm) for f, t in (('diccionario', DIC), ('reglas', REG)) for l in t.splitlines() if 'PROVISIONAL' in l for nm in re.findall(PATB, l) if nm in resueltos]
prov_res += [('migraciones', nm) for i_, t_ in enumerate(TXT) for l in t_.splitlines() if 'PROVISIONAL' in l for nm in re.findall(PATB, l) if nm in resueltos and not RESUELVE.get(nm, -1) > i_]
k5 = ['S1-ajuste-objetivo-kcal', 'S1-fuentes-precio', 'criterios-promedios-redondeo', 'criterios-frescura-extremos', 'C1-valores-k31']
notif13 = dict(x("SELECT code, params_json FROM notification WHERE code IN ('payment_due', 'budget_threshold')").fetchall())
notif_ok = notif13 == {'payment_due': '{}', 'budget_threshold': '{"scopes": ["category", "total"]}'} and x('SELECT updated_at FROM notification ORDER BY rowid').fetchall() == [(r[-1],) for r in N12]
ok44 = (len(cols13) == 2 and all(casa.get(k_) == v for k_, v in {**SEM13, **PREV}.items()) and all(anot.get(k_) == casa.get(k_) for k_ in anot) and len(anot) == 30
        and not acept13 and not rech13 and txt_ok and not prov_res and all(k_ in resueltos for k_ in k5) and notif_ok)
out('K44', ok44, f"columnas nuevas de health_profile={cols13} · claves sembradas por la 00003 con su valor inicial={all(casa.get(k_) == v for k_, v in SEM13.items())} ({len(SEM13)}) · las 5 de 00001/00002 intactas={all(casa.get(k_) == v for k_, v in PREV.items())} · anotaciones @param={len(anot)}, iguales a lo sembrado={all(anot.get(k_) == casa.get(k_) for k_ in anot)}"
    f" · negativas: {len(NEG13)}, rechazadas {len(NEG13) - len(acept13)}, aceptadas por error {acept13} · positivas: {len(POS13)}, aceptadas {len(POS13) - len(rech13)}, rechazadas por error {rech13}"
    f" · cifras: RN-17 meta {k36[0]} y con 500 kcal para bajar {k36[1]}; RN-12 pesos {[str(float(v)) for v in k40]}; RN-27 semanal {[str(float(v)) for v in k39]}; presentes en las reglas={txt_ok}"
    f" · marcas PROVISIONAL que citan un bloqueo resuelto={prov_res} · los 5 bloqueos de esta puerta resueltos={all(k_ in resueltos for k_ in k5)} · params_json de payment_due y budget_threshold sin days_before ni thresholds_pct, con el updated_at de 1.2={notif_ok} ({notif13})")
# ---------- AVISOS13 (corrección S1-r1 de P1c, H-01): traspaso params_json <-> app_setting con valores que no son los iniciales
av = sqlite3.connect(':memory:'); av.executescript(UPS[0]); av.executescript(UPS[1])
av.execute("UPDATE notification SET params_json = '{\"thresholds_pct\": [75, 95], \"scopes\": [\"total\"]}' WHERE code = 'budget_threshold'")
av.execute("UPDATE notification SET params_json = '{\"days_before\": [3, 1, 0]}' WHERE code = 'payment_due'"); av.commit()
qn = "SELECT code, params_json, updated_at FROM notification WHERE code IN ('payment_due', 'budget_threshold') ORDER BY code"
qa = "SELECT setting_key, value_json FROM app_setting WHERE setting_key LIKE 'alerts.%' AND user_account_id IS NULL"
n_a = av.execute(qn).fetchall(); av.executescript(UPS[2]); a_up = dict(av.execute(qa).fetchall()); n_up = av.execute(qn).fetchall()
av.execute("UPDATE app_setting SET value_json = '70' WHERE setting_key = 'alerts.budget_warning_pct'")
av.execute("UPDATE app_setting SET value_json = '[2,0]' WHERE setting_key = 'alerts.payment_due_days_before'"); av.commit()
av.executescript(DOWNS[2]); n_dn = av.execute(qn).fetchall(); av.close()
ok_av = (a_up == {'alerts.budget_warning_pct': '75', 'alerts.budget_limit_pct': '95', 'alerts.payment_due_days_before': '[3, 1, 0]'}
         and [json.loads(r[1]) for r in n_up] == [{'scopes': ['total']}, {}]
         and [json.loads(r[1]) for r in n_dn] == [{'thresholds_pct': [70, 95], 'scopes': ['total']}, {'days_before': [2, 0]}]
         and [r[2] for r in n_a] == [r[2] for r in n_up] == [r[2] for r in n_dn])
out('AVISOS13', ok_av, f"Up con avisos cambiados en 1.2 → app_setting {a_up}, params_json {[r[1] for r in n_up]} · Down tras poner 70 % y [2,0] en 1.3 → {[r[1] for r in n_dn]} · updated_at conservado={[r[2] for r in n_a] == [r[2] for r in n_up] == [r[2] for r in n_dn]}")
# ---------- PROG14 (puerta 2, Q2-11 B, K47): weekday y day_of_month ya no están en notification.params_json; los cuatro
# parámetros de avisos son claves de app_setting sin marca PROVISIONAL en la migración nueva, el diccionario ni RN-29
prog = dict(x("SELECT code, params_json FROM notification WHERE code IN ('weekly_summary', 'meal_plan_missing', 'monthly_report')").fetchall())
SEM14 = {'reports.daily_send_time': '"21:30"', 'reports.weekly_weekday': '7', 'reports.weekly_send_time': '"19:00"', 'reports.monthly_day_of_month': '1', 'reports.monthly_send_time': '"08:00"', 'alerts.send_time': '"09:00"', 'alerts.stale_prices_advance_days': '1', 'alerts.meal_plan_missing_weekday': '7', 'optimizer.date_window_days': '7', 'ocr.timeout_seconds': '30'}
cerr = {l.split(' | ')[0].strip() for l in re.findall(r'^-- @param_sin_valor_cerrado (.+)$', mig, re.M)}
nv_l = [l.split(' | ')[0].strip() for l in re.findall(r'^-- @param_sin_valor (.+)$', mig, re.M)]
nv_vivos = [n_ for n_ in nv_l if n_ not in cerr and any(k_ in n_ for k_ in ('send_time', 'advance_days', 'weekday', 'day_of_month', 'date_window_days'))]
s44 = DIC.split('## Parámetros operativos de app_setting (K44)')[1].split('\n## ')[0]
def _sin_prov14(sec_): return not any('PROVISIONAL' in l_ for l_ in sec_.splitlines() if any(f'`{k_}`' in l_ for k_ in SEM14))  # 1.5: solo las filas de la 00004
dic_p = all(f'| `{k_}` |' in s44 for k_ in SEM14) and _sin_prov14(s44)
rn29 = REG.split('\n## RN-29 ')[1].split('\n## ')[0]
reg_p = all(f'`{k_}`' in rn29 for k_ in SEM14) and _sin_prov14(rn29) and 'Q2-66' in rn29
blq = open(R + '/docs/estado/bloqueos/S1-parametros-avisos.md', encoding='utf-8').read()
ok_pg = (prog == {'weekly_summary': '{}', 'meal_plan_missing': '{}', 'monthly_report': '{}'} and all(casa.get(k_) == v_ for k_, v_ in SEM14.items())
         and not nv_vivos and len(cerr) == 5 and dic_p and reg_p and 'Respuesta de Axel' in blq and RESUELVE.get('S1-parametros-avisos') == 3)
out('PROG14', ok_pg, f"params_json en 1.4 {prog} · claves de la 00004 con su valor inicial={all(casa.get(k_) == v_ for k_, v_ in SEM14.items())} ({len(SEM14)}) · @param_sin_valor cerrados={len(cerr)}, vivos de avisos o ventana={nv_vivos} · diccionario K44 con las 10 y sin PROVISIONAL={dic_p} · RN-29 con las 10, sin PROVISIONAL y con Q2-66={reg_p} · bloqueo con respuesta y @resuelve en la 00004={'Respuesta de Axel' in blq and RESUELVE.get('S1-parametros-avisos') == 3}")
# ---------- K47 (puerta 2): rangos de una sola clave, horas «HH:MM» y credit.vat_included (negativas y positivas)
def aj(k_, v_): return f"UPDATE app_setting SET value_json = '{v_}' WHERE setting_key = '{k_}' AND user_account_id IS NULL"
NEG14 = [('semana_8', aj('reports.weekly_weekday', '8')), ('plan_semana_0', aj('alerts.meal_plan_missing_weekday', '0')), ('dia_32', aj('reports.monthly_day_of_month', '32')),
 ('dia_0', aj('reports.monthly_day_of_month', '0')), ('dia_texto', aj('reports.monthly_day_of_month', '"1"')), ('hora_24', aj('alerts.send_time', '"24:00"')), ('hora_un_digito', aj('reports.daily_send_time', '"9:00"')),
 ('hora_numero', aj('reports.weekly_send_time', '930')), ('hora_minutos_60', aj('reports.monthly_send_time', '"08:60"')), ('comparador_31', aj('comparator.timeout_seconds', '31')),
 ('comparador_0', aj('comparator.timeout_seconds', '0')), ('viejo_0', aj('price.stale_after_days', '0')), ('ventana_negativa', aj('optimizer.date_window_days', '-1')),
 ('ocr_texto', aj('ocr.timeout_seconds', '"30"')), ('anticipacion_decimal', aj('alerts.stale_prices_advance_days', '1.5')),
 ('iva_incluido_con_porcentaje', T18 + f"annual_rate_num,annual_rate_den,annual_rate_decimals,rate_base,vat_on_interest_bp,min_payment_pct_bp,min_payment_fixed_cents,vat_included) VALUES ('Tarjeta K47','bank_card',{AX},500000,'monthly',5,25,455,1000,1,'daily',1600,500,20000,1)"),
 ('iva_sin_dato_con_porcentaje', T18 + f"annual_rate_num,annual_rate_den,annual_rate_decimals,rate_base,vat_on_interest_bp,min_payment_pct_bp,min_payment_fixed_cents) VALUES ('Tarjeta K47','bank_card',{AX},500000,'monthly',5,25,455,1000,1,'daily',1600,500,20000)"),
 ('iva_sin_dato_con_porcentaje_al_editar', "UPDATE credit SET vat_included = NULL, vat_on_interest_bp = 1600 WHERE id = (SELECT min(id) FROM credit)"),
 ('iva_2', "UPDATE credit SET vat_included = 2 WHERE id = (SELECT min(id) FROM credit)")]
POS14 = [('dia_31', aj('reports.monthly_day_of_month', '31')), ('hora_medianoche', aj('alerts.send_time', '"00:00"')), ('hora_2359', aj('reports.daily_send_time', '"23:59"')),
 ('comparador_30_y_1', aj('comparator.timeout_seconds', '30') + '; ' + aj('comparator.timeout_seconds', '1')), ('ventana_0', aj('optimizer.date_window_days', '0')), ('viejo_1', aj('price.stale_after_days', '1')),
 ('semana_lunes', aj('reports.weekly_weekday', '1')), ('anticipacion_mismo_dia', aj('alerts.stale_prices_advance_days', '0')),
 ('tarjeta_iva_aparte', T18 + f"annual_rate_num,annual_rate_den,annual_rate_decimals,rate_base,vat_on_interest_bp,min_payment_pct_bp,min_payment_fixed_cents,vat_included) VALUES ('Tarjeta K47','bank_card',{AX},500000,'monthly',5,25,455,1000,1,'daily',1600,500,20000,0)"),
 ('tarjeta_iva_por_capturar', T18 + f"annual_rate_num,annual_rate_den,annual_rate_decimals,rate_base,min_payment_pct_bp,min_payment_fixed_cents) VALUES ('Tarjeta K47','bank_card',{AX},500000,'monthly',5,25,455,1000,1,'daily',500,20000)"),
 ('linea_incluye_iva', "UPDATE credit SET vat_included = 1, vat_on_interest_bp = NULL WHERE id = (SELECT min(id) FROM credit)")]
acept14 = corre(NEG14, True); rech14 = corre(POS14, False)
nul = one("SELECT count(*) FROM credit WHERE vat_included IS NULL")
out('K47', not acept14 and not rech14 and has('credit.vat_included'), f"credit.vat_included={has('credit.vat_included')} (créditos capturados después de migrar, NULL = falta capturar: {nul}) · negativas: {len(NEG14)}, rechazadas {len(NEG14) - len(acept14)}, aceptadas por error {acept14} · positivas: {len(POS14)}, aceptadas {len(POS14) - len(rech14)}, rechazadas por error {rech14}")
# ---------- K54 (1.5, 00005): cinco claves del tablero con su valor inicial y su rango (PROVISIONAL, bloqueo C1-tablero)
SEM15 = {'debt.capacity_warning_pct': '30', 'debt.capacity_limit_pct': '40', 'dashboard.payment_reminder_count': '5', 'budget.signal_months': '3', 'budget.signal_margin_pct': '10'}
NEG15 = [('umbral_0', aj('debt.capacity_warning_pct', '0')), ('tope_101', aj('debt.capacity_limit_pct', '101')), ('umbral_texto', aj('debt.capacity_warning_pct', '"30"')),
 ('tope_decimal', aj('debt.capacity_limit_pct', '40.5')), ('pagos_0', aj('dashboard.payment_reminder_count', '0')), ('pagos_21', aj('dashboard.payment_reminder_count', '21')),
 ('meses_1', aj('budget.signal_months', '1')), ('meses_13', aj('budget.signal_months', '13')), ('margen_0', aj('budget.signal_margin_pct', '0')), ('margen_51', aj('budget.signal_margin_pct', '51')),
 ('alta_fuera_de_rango', "INSERT INTO app_setting (setting_key, value_json, user_account_id) VALUES ('budget.signal_months', '20', NULL)")]
POS15 = [('umbral_1', aj('debt.capacity_warning_pct', '1')), ('tope_100', aj('debt.capacity_limit_pct', '100')), ('pagos_1_y_20', aj('dashboard.payment_reminder_count', '1') + '; ' + aj('dashboard.payment_reminder_count', '20')),
 ('meses_2_y_12', aj('budget.signal_months', '2') + '; ' + aj('budget.signal_months', '12')), ('margen_1_y_50', aj('budget.signal_margin_pct', '1') + '; ' + aj('budget.signal_margin_pct', '50'))]
acept15 = corre(NEG15, True); rech15 = corre(POS15, False)
sem15_ok = all(casa.get(k_) == v_ for k_, v_ in SEM15.items())
out('K54', sem15_ok and not acept15 and not rech15 and all(f'`{k_}`' in DIC and f'`{k_}`' in REG for k_ in SEM15),
    f"claves de la 00005 con su valor inicial={sem15_ok} ({len(SEM15)}) · negativas: {len(NEG15)}, rechazadas {len(NEG15) - len(acept15)}, aceptadas por error {acept15} · positivas: {len(POS15)}, aceptadas {len(POS15) - len(rech15)}, rechazadas por error {rech15} · en diccionario y reglas={all(f'`{k_}`' in DIC and f'`{k_}`' in REG for k_ in SEM15)}")
# ---------- AVISOS14: traspaso params_json <-> app_setting con valores que no son los iniciales y migración de vat_included
av = sqlite3.connect(':memory:'); [av.executescript(u) for u in UPS[:3]]
av.execute("UPDATE notification SET params_json = '{\"weekday\": 5}' WHERE code = 'weekly_summary'")
av.execute("UPDATE notification SET params_json = '{\"day_of_month\": 31, \"extra\": true}' WHERE code = 'monthly_report'")
av.execute("UPDATE notification SET params_json = '{\"weekday\": 9}' WHERE code = 'meal_plan_missing'")
av.execute("INSERT INTO person(code, display_name, participation) VALUES ('axel', 'Axel', 'holder')")
av.execute("INSERT INTO credit(name, kind, holder_person_id, payment_frequency, vat_on_interest_bp) VALUES ('Préstamo con IVA aparte', 'term_loan', 1, 'monthly', 1600), ('Préstamo con IVA incluido', 'term_loan', 1, 'monthly', NULL)"); av.commit()
qn = "SELECT code, params_json, updated_at FROM notification WHERE code IN ('weekly_summary', 'monthly_report', 'meal_plan_missing') ORDER BY id"
qc = "SELECT name, vat_included, updated_at FROM credit ORDER BY id"
n_a = av.execute(qn).fetchall(); c_a = av.execute("SELECT updated_at FROM credit ORDER BY id").fetchall()
av.executescript(UPS[3]); a_up = dict(av.execute("SELECT setting_key, value_json FROM app_setting WHERE setting_key IN ('reports.weekly_weekday', 'reports.monthly_day_of_month', 'alerts.meal_plan_missing_weekday')").fetchall())
n_up = av.execute(qn).fetchall(); c_up = av.execute(qc).fetchall()
av.execute("UPDATE app_setting SET value_json = '3' WHERE setting_key = 'reports.weekly_weekday'"); av.commit()
av.executescript(DOWNS[3]); n_dn = av.execute(qn).fetchall(); av.close()
ok_a14 = (a_up == {'reports.weekly_weekday': '5', 'reports.monthly_day_of_month': '31', 'alerts.meal_plan_missing_weekday': '7'}
          and [json.loads(r[1]) for r in n_up] == [{}, {'extra': True}, {}] and [r[1] for r in c_up] == [0, 1]
          and [json.loads(r[1]) for r in n_dn] == [{'weekday': 3}, {'day_of_month': 31, 'extra': True}, {'weekday': 7}]
          and [r[2] for r in n_a] == [r[2] for r in n_up] == [r[2] for r in n_dn] and [r[0] for r in c_a] == [r[2] for r in c_up])
out('AVISOS14', ok_a14, f"Up con avisos cambiados en 1.3 (semana 5, día 31 con otro dato, plan 9 inválido) → app_setting {a_up}, params_json {[r[1] for r in n_up]} · vat_included de los créditos de 1.3 (con 1600 y sin porcentaje) = {[r[1] for r in c_up]} · Down tras poner semana 3 → {[r[1] for r in n_dn]} · updated_at conservado (avisos y créditos)={[r[2] for r in n_a] == [r[2] for r in n_up] == [r[2] for r in n_dn] and [r[0] for r in c_a] == [r[2] for r in c_up]}")
# ---------- DDL de 00001 y 00002 sin cambios contra git HEAD (sin comentarios)
def sin_coment(t): return '\n'.join(ln for ln in (re.sub(r'\s*--.*$', '', l).rstrip() for l in t.splitlines()) if ln)
ddl = []
for q in MIGS[:3]:  # 1.4: también la 00003, ya en git HEAD (00001 a 00003 no se editan)
    r_ = subprocess.run(['git', '-C', R, 'show', 'HEAD:' + os.path.relpath(q, R)], capture_output=True, text=True)
    ddl.append((os.path.basename(q), r_.returncode == 0 and sin_coment(r_.stdout) == sin_coment(open(q, encoding='utf-8').read())))
out('DDL', all(ok_ for _, ok_ in ddl) and len(ddl) == 3, f"DDL sin comentarios igual al de git HEAD: {ddl}")
# ---------- Down de la 00003 (vuelve el esquema 1.2 con sus datos), Down de la 00002 (vuelve el 1.1), Down de la 00001 y Up otra vez
db.close(); db = sqlite3.connect(DB)
def datos12(): return {t: db.execute(f'SELECT {", ".join(cs)} FROM {t} ORDER BY rowid').fetchall() for t, cs in cols12.items() if t != 'app_setting'}
D_antes = datos12()
def datos13(t_): return db.execute(f'SELECT {", ".join(cols13t[t_])} FROM {t_} ORDER BY rowid').fetchall()
C_antes = datos13('credit')
for d_ in reversed(DOWNS[3:]): db.executescript(d_)
E_d4 = esquema(db); A_d4 = ajustes(db); N_d4 = datos13('notification'); C_d4 = datos13('credit')
out('DOWN14', E_d4 == E13 and A_d4 == A13 and N_d4 == N13 and C_d4 == C_antes and len(C_antes) > 0, f"Down de la 00004: sqlite_master igual al esquema 1.3 = {E_d4 == E13} ({len(E13)} objetos, diferencias={len(set(E13) ^ set(E_d4))}) · app_setting igual a 1.3 = {A_d4 == A13} ({len(A13)} claves) · avisos de 1.3 byte por byte (params_json y updated_at) = {N_d4 == N13} · créditos con sus columnas de 1.3 iguales = {C_d4 == C_antes} ({len(C_antes)} filas)")
db.executescript(DOWNS[2])
E_d3 = esquema(db); A_d3 = ajustes(db); D_d3 = datos12(); dif3 = set(E12) ^ set(E_d3)
n_filas = sum(len(v) for v in D_d3.values())
out('DOWN13', E_d3 == E12 and A_d3 == A12 and D_d3 == {**D_antes, 'notification': N12} and n_filas > 0, f"Down de la 00003: sqlite_master igual al esquema 1.2 = {E_d3 == E12} ({len(E12)} objetos, diferencias={len(dif3)}) · app_setting igual a 1.2 = {A_d3 == A12} ({len(A12)} claves) · datos de las columnas de 1.2 en las otras {len(D_d3)} tablas iguales (notification con el params_json y el updated_at de 1.2) = {D_d3 == {**D_antes, 'notification': N12}} ({n_filas} filas, incluido el perfil con ajuste de kcal)")
db.executescript(DOWNS[1])
E_down = esquema(db); A_down = ajustes(db); difs = set(E11) ^ set(E_down)
db.executescript(DOWNS[0])
resto = db.execute("SELECT count(*) FROM sqlite_master WHERE name NOT LIKE 'sqlite_%'").fetchone()[0]
for u in UPS: db.executescript(u)
t2 = db.execute("SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").fetchone()[0]
ic2 = db.execute('PRAGMA integrity_check').fetchone()[0]; E_again = esquema(db); db.close()
out('DOWN', E_down == E11 and A_down == A11 and resto == 0 and t2 == len(tabs) and ic2 == 'ok' and E_again == EF, f"Down de la 00002: sqlite_master igual al esquema 1.1 = {E_down == E11} ({len(E11)} objetos, diferencias={len(difs)}) · app_setting igual a 1.1 = {A_down == A11} · Down de la 00001 deja {resto} objetos · Up otra vez: {t2} tablas, integrity_check={ic2}, esquema igual al final ({len(MIGS)} migraciones) = {E_again == EF}")
gs = os.environ.get('GOOSE') or shutil.which('goose')
if gs:
    gdb = TMP + '/s1-verifica-goose.db'
    if os.path.exists(gdb): os.remove(gdb)
    res = []; obj = []
    for a_ in (['up'], ['down-to', '0'], ['up']):
        res.append(subprocess.run([gs, '-dir', R + '/backend/db/migrations', 'sqlite3', gdb] + a_, capture_output=True, text=True).returncode)
        g = sqlite3.connect(gdb); obj.append(g.execute("SELECT count(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name <> 'goose_db_version'").fetchone()[0]); g.close()
    g = sqlite3.connect(gdb)
    ver = g.execute('SELECT max(version_id) FROM goose_db_version WHERE is_applied').fetchone()[0]
    cols_g = {tt: [r[1] for r in g.execute(f'PRAGMA table_info({tt})')] for tt in tabs}
    nom_g = sorted(g.execute("SELECT type, name FROM sqlite_master WHERE name NOT LIKE 'sqlite_%' AND name <> 'goose_db_version'").fetchall()); g.close(); os.remove(gdb)
    nom_p = sorted((r[0], r[1]) for r in EF)
    out('GOOSE', res == [0, 0, 0] and obj == [len(tabs), 0, len(tabs)] and ver == len(MIGS) and cols_g == cols and nom_g == nom_p, f"goose ({os.path.basename(gs)}) up / down-to 0 / up → códigos {res}, tablas tras cada paso={obj}, versión final={ver}, mismas columnas y objetos que con sqlite3={cols_g == cols and nom_g == nom_p}")
else:
    print('GOOSE -    goose no está en el PATH ni en $GOOSE: se omitió (S1-02 ya probó el SQL con sqlite3)')
os.remove(DB)
if not os.listdir(TMP): os.rmdir(TMP)
print('RESUMEN', 'TODO OK' if not FALLAS else f'FALLAS: {FALLAS}')
sys.exit(1 if FALLAS else 0)
