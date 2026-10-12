#!/usr/bin/env python3
"""Modelo v3 de despensa (MXN). Decimal, ROUND_HALF_UP a centavos. Karen fija, credito sin intereses."""
import json, calendar, datetime as dt, os
from decimal import Decimal as D, ROUND_HALF_UP, ROUND_DOWN

def q(x): return D(x).quantize(D('0.01'), rounding=ROUND_HALF_UP)
def f(x): return float(q(x))
def iso(d): return d.isoformat()
def S(it): return sum(it, D(0))

R = D('800.50') / D('3094.50')
PAPEL = D('199.00'); LIMITE = D('10000.00')
START, END = dt.date(2026,10,10), dt.date(2027,4,10)
PERIODS = [(dt.date(2026,10,10),dt.date(2026,11,9)),(dt.date(2026,11,10),dt.date(2026,12,9)),
           (dt.date(2026,12,10),dt.date(2027,1,9)),(dt.date(2027,1,10),dt.date(2027,2,9)),
           (dt.date(2027,2,10),dt.date(2027,3,9)),(dt.date(2027,3,10),dt.date(2027,4,9))]
MESES = ['ene','feb','mar','abr','may','jun','jul','ago','sep','oct','nov','dic']
def lab(d): return f"{d.day}-{MESES[d.month-1]}-{d.year}"

def sched(buy):
    c = set()
    for k in range(8):
        yy = buy.year + (buy.month-1+k)//12; mm = (buy.month-1+k) % 12 + 1
        last = calendar.monthrange(yy, mm)[1]
        for day in (15, 30):
            d = dt.date(yy, mm, min(day, last))
            if (d-buy).days >= 15: c.add(d)
    return sorted(c)[:4]

# ---- constantes Karen fija / conjunto
K_CAP, K_INT, K_TOT, K_PAGO = D('1191.50'), D('308.22'), D('1499.72'), D('374.93')
assert K_PAGO*4 == K_TOT and K_CAP + K_INT == K_TOT
CONJ = D('851.00'); K_PROPIO = K_CAP - CONJ   # 340.50

def mk(n, buy, cap, inter, total, montos, ax_cap, ax_int, papel_ax, papel_ka, ax_pagos=None, ka_pagos=None):
    dates = sched(buy)
    pagos = []
    for i, (d, m) in enumerate(zip(dates, montos)):
        ka = ka_pagos[i]; ax = ax_pagos[i] if ax_pagos else m - ka
        cr = [D('773.63')]*3 + [D('773.61')] if n == 1 else [cap/4]*4
        pagos.append(dict(fecha=d, monto=m, axel=ax, karen=ka, capital_repuesto=cr[i]))
    ka_cap = cap - ax_cap
    def pers(capx, intx, papel, pagos_p):
        conj = CONJ; propio = capx - conj
        return dict(capital=capx, conjunto=conj, propio=propio, papel_dentro=papel, libre=capx-papel,
                    intereses=intx, total=capx+intx, pago_por_quincena=pagos_p)
    return dict(n=n, fecha=buy, capital=cap, intereses=inter, total=total, pagos=pagos,
        axel=pers(ax_cap, ax_int, papel_ax, pagos[0]['axel']), karen=pers(ka_cap, inter-ax_int, papel_ka, pagos[0]['karen']))

desp = {}
desp[1] = mk(1, START, D('3094.50'), D('800.50'), D('3895.00'), [D('973.75')]*4,
             D('1903.00'), D('492.28'), D(0), PAPEL, ax_pagos=[D('598.82')]*4, ka_pagos=[K_PAGO]*4)
TOT_N = q(D(4000)*(1+R)); INT_N = TOT_N - D(4000)
assert TOT_N == D('5034.74') and INT_N == D('1034.74')
base = (TOT_N/4).quantize(D('0.01'), rounding=ROUND_DOWN)
MONTOS = [base]*3 + [TOT_N - base*3]
assert MONTOS == [D('1258.68')]*3 + [D('1258.70')]
AX_CAP = D(4000) - K_CAP; AX_INT = INT_N - K_INT
assert AX_CAP == D('2808.50') and AX_INT == D('726.52') and AX_CAP + AX_INT == D('3535.02')
for n, buy in ((2, dt.date(2026,12,10)), (3, dt.date(2027,2,10)), (4, dt.date(2027,4,10))):
    desp[n] = mk(n, buy, D(4000), INT_N, TOT_N, MONTOS, AX_CAP, AX_INT, PAPEL, D(0), ka_pagos=[K_PAGO]*4)
assert [p['axel'] for p in desp[2]['pagos']] == [D('883.75')]*3 + [D('883.77')]
assert [iso(p['fecha']) for p in desp[2]['pagos']] == ['2026-12-30','2027-01-15','2027-01-30','2027-02-15']
assert [iso(p['fecha']) for p in desp[3]['pagos']] == ['2027-02-28','2027-03-15','2027-03-30','2027-04-15']
assert [iso(p['fecha']) for p in desp[4]['pagos']] == ['2027-04-30','2027-05-15','2027-05-30','2027-06-15']
assert [iso(p['fecha']) for p in desp[1]['pagos']] == ['2026-10-30','2026-11-15','2026-11-30','2026-12-15']

for k, d in desp.items():
    assert S(p['monto'] for p in d['pagos']) == d['total'], k
    assert S(p['capital_repuesto'] for p in d['pagos']) == d['capital'], k
    assert d['capital'] + d['intereses'] == d['total']
    for p in d['pagos']: assert p['axel'] + p['karen'] == p['monto']
    for w in ('axel','karen'):
        assert S(p[w] for p in d['pagos']) == d[w]['total'], (k, w)
    assert d['axel']['total'] + d['karen']['total'] == d['total']
    assert d['axel']['capital'] + d['karen']['capital'] == d['capital']
    assert d['axel']['intereses'] + d['karen']['intereses'] == d['intereses']
    assert d['karen']['capital'] == K_CAP and d['karen']['intereses'] == K_INT and d['karen']['total'] == K_TOT
assert desp[1]['axel']['total'] == D('2395.28') and desp[1]['axel']['pago_por_quincena'] == D('598.82')

PAPEL_CONTADO = [(dt.date(2027,1,10),'Karen'),(dt.date(2027,3,10),'Karen')]

# ---- periodos
periodos = []; cd = ce = D(0)
for i, (a, b) in enumerate(PERIODS):
    d = desp[i//2+1]
    h = q(d['total']/2); hk = q(d['karen']['total']/2); ha = h - hk
    if d['n'] == 1:
        assert (ha, hk) == (D('1197.64'), D('749.86'))
    pc = [x for x in PAPEL_CONTADO if a <= x[0] <= b]
    pap = PAPEL*len(pc); pk = pap  # todo contado es de Karen
    dev = dict(despensa_total=h, papel_contado=pap, total=h+pap, axel=ha, karen=hk+pk)
    pays = [(x['n'], p) for x in desp.values() for p in x['pagos'] if a <= p['fecha'] <= b]
    ps = S(p['monto'] for _, p in pays)
    ef = dict(pagos=ps, papel=pap, total=ps+pap, axel=S(p['axel'] for _, p in pays),
              karen=S(p['karen'] for _, p in pays)+pk, n_pagos=len(pays),
              detalle_pagos=[dict(despensa=n, fecha=iso(p['fecha']), monto=f(p['monto']), axel=f(p['axel']), karen=f(p['karen'])) for n, p in pays])
    assert dev['axel']+dev['karen'] == dev['total'] and ef['axel']+ef['karen'] == ef['total']
    cd += dev['total']; ce += ef['total']
    periodos.append(dict(periodo=f"P{i+1}", etiqueta=f"{lab(a)} a {lab(b)}", inicio=a, fin=b, despensa=d['n'], dev=dev, ef=ef, dif=cd-ce))
exp_dev = ['1947.50','1947.50','2517.37','2716.37','2517.37','2716.37']
exp_ef = ['973.75','1947.50','2232.43','2716.36','2517.38','2716.36']
assert [p['dev']['total'] for p in periodos] == [D(x) for x in exp_dev]
assert [p['ef']['total'] for p in periodos] == [D(x) for x in exp_ef]
assert (periodos[2]['ef']['axel'], periodos[2]['ef']['karen']) == (D('1482.57'), D('749.86'))
assert (periodos[3]['ef']['axel'], periodos[3]['ef']['karen']) == (D('1767.50'), D('948.86'))
assert (periodos[4]['ef']['axel'], periodos[4]['ef']['karen']) == (D('1767.52'), D('749.86'))
assert (periodos[5]['ef']['axel'], periodos[5]['ef']['karen']) == (D('1767.50'), D('948.86'))
assert periodos[1]['ef']['karen'] == D('749.86')
sum_dev = S(p['dev']['total'] for p in periodos); sum_ef = S(p['ef']['total'] for p in periodos)
pagos_hor = S(x['monto'] for d in desp.values() for x in d['pagos'] if x['fecha'] <= END)
assert sum_ef == pagos_hor + PAPEL*2 == S(p['ef']['pagos'] for p in periodos) + S(p['ef']['papel'] for p in periodos)
assert pagos_hor == D('3895.00') + D('5034.74') + D('3776.04')
assert sum_ef == D('13103.78')

# ---- credito (sin intereses)
ev = [(START,'inicio','Disponible hoy',D(0),D(0),D(0))]
for d in desp.values():
    if d['n'] > 1: ev.append((d['fecha'],'compra',f"Despensa {d['n']} (a credito, solo capital)",d['capital'],-d['capital'],D(0)))
    for p in d['pagos']: ev.append((p['fecha'],'pago',f"Pago despensa {d['n']} (repone capital)",p['monto'],p['capital_repuesto'],p['monto']-p['capital_repuesto']))
for dd, o in PAPEL_CONTADO: ev.append((dd,'papel',f"Papel de contado ({o}); no mueve el credito",PAPEL,D(0),D(0)))
order = {'inicio':0,'compra':1,'papel':2,'pago':3}
ev = sorted([e for e in ev if e[0] <= END], key=lambda e: (e[0], order[e[1]]))
A = LIMITE; B = LIMITE - D('3094.50'); cred = []
for d, k, l, m, delta, intr in ev:
    A += delta; B = min(LIMITE, B + delta)
    assert A >= 0 and B >= 0
    cred.append(dict(fecha=iso(d), tipo=k, descripcion=l, monto=f(m), capital_repuesto=f(delta) if k=='pago' else 0.0,
                     interes_no_repuesto=f(intr), disponible_A=f(A), disponible_B=f(B)))
expA = ['10773.63','11547.26','12320.89','8320.89','9094.50','10094.50','11094.50','12094.50','8094.50','9094.50','10094.50','11094.50','12094.50','8094.50']
expB = ['7679.13','8452.76','9226.39','5226.39','6000.00','7000.00','8000.00','9000.00','5000.00','6000.00','7000.00','8000.00','9000.00','5000.00']
ser = [c for c in cred if c['tipo'] != 'inicio' and c['tipo'] != 'papel']
assert [c['disponible_A'] for c in ser] == [float(x) for x in expA], [c['disponible_A'] for c in ser]
assert [c['disponible_B'] for c in ser] == [float(x) for x in expB], [c['disponible_B'] for c in ser]
cred_cmp = dict(esperado_A=[float(x) for x in expA], esperado_B=[float(x) for x in expB],
                calculado_A=[c['disponible_A'] for c in ser], calculado_B=[c['disponible_B'] for c in ser], coinciden=True)

# ---- valor estable
MITAD = q(TOT_N/2); PAP_M = PAPEL/2   # 99.50
mes = MITAD + PAP_M; ka_m = q(desp[2]['karen']['total']/2) + PAP_M; ax_m = MITAD - q(desp[2]['karen']['total']/2)
assert (q(mes), ka_m, ax_m) == (D('2616.87'), D('849.36'), D('1767.51')) and ka_m + ax_m == q(mes)
qm, qa, qk = q(mes/2), q(ax_m/2), q(ka_m/2)
ve = dict(mes=dict(total=f(mes), axel=f(ax_m), karen=f(ka_m)),
          quincena=dict(total=f(qm), axel=f(qa), karen=f(qk), suma_axel_karen=f(qa+qk), diferencia_centavo=f(qm-qa-qk),
                        nota="redondeo a centavo; diferencia de 0.01 si suma != total"),
          componentes=dict(mitad_despensa=f(MITAD), papel_contado_mes=f(PAP_M), karen_mitad_despensa=f(ka_m-PAP_M),
                           axel_mitad_despensa=f(ax_m), formula="5034.74/2 + 199/2"))
assert qm == D('1308.44') and qa == D('883.76') and qk == D('424.68')

# ---- sensibilidad V11
sens = []
for N in (3000,3500,4000,4500,5000):
    t = D(N)*(1+R); ax_t = t - K_TOT; ax_cap = D(N) - K_CAP
    sens.append(dict(N=N, intereses=f(t-N), total=f(t), pago_quincenal=f(t/4),
        axel_total=f(ax_t), axel_capital=f(ax_cap), axel_menor_que_actual=bool(ax_cap < D('1903.00')),
        costo_mensual_total=f(t/2+PAP_M), axel_mensual=f(ax_t/2), karen_mensual=f(ka_m)))
assert [s for s in sens if s['N']==4000][0]['costo_mensual_total'] == 2616.87

# ---- papel
papel = [
 dict(paquete=1, cubre='noviembre 2026', compra='2026-10-10', quien='Karen', forma='dentro de despensa 1', nota='ya incluido en sus 1,191.50; no hay compra el 10-nov'),
 dict(paquete=2, cubre='diciembre 2026', compra='2026-12-10', quien='Axel', forma='dentro de despensa 2', nota=''),
 dict(paquete=3, cubre='enero 2027', compra='2027-01-10', quien='Karen', forma='de contado', nota='credito sigue activo; fuera del cap de 374.93'),
 dict(paquete=4, cubre='febrero 2027', compra='2027-02-10', quien='Axel', forma='dentro de despensa 3', nota=''),
 dict(paquete=5, cubre='marzo 2027', compra='2027-03-10', quien='Karen', forma='de contado', nota='fuera del cap de 374.93'),
 dict(paquete=6, cubre='abril 2027', compra='2027-04-10', quien='Axel', forma='dentro de despensa 4', nota=''),
]
for p in papel: p['monto'] = f(PAPEL)

# ---- disposicion V12
disp = []
for n in (1,2,3,4):
    d = desp[n]
    row = dict(despensa=n, fecha=iso(d['fecha']))
    for w in ('axel','karen'):
        x = d[w]; row[w] = {k: f(v) for k, v in x.items()}
    disp.append(row)

# ---- horizonte V13
pend = S(p['monto'] for d in desp.values() for p in d['pagos'] if p['fecha'] > END)
pend_d3 = S(p['monto'] for p in desp[3]['pagos'] if p['fecha'] > END)
assert pend_d3 == D('1258.70') and pend == D('6293.44') and pend == pend_d3 + desp[4]['total']
horizonte = dict(inicio=iso(START), fin=iso(END), pagos_despensa_1=f(desp[1]['total']), pagos_despensa_2=f(desp[2]['total']),
    pagos_despensa_3_en_horizonte=f(S(p['monto'] for p in desp[3]['pagos'] if p['fecha'] <= END)),
    papel_contado=f(PAPEL*2), total_efectivo=f(sum_ef), suma_periodos_efectivo=f(sum_ef), total_devengado=f(sum_dev),
    deuda_viva_10_abr=f(pend), deuda_viva_detalle=dict(despensa_3_pendiente=f(pend_d3), despensa_4=f(desp[4]['total'])),
    promedio_mensual_efectivo=f(sum_ef/6), promedio_mensual_devengado=f(sum_dev/6), diferencia_final=f(sum_dev-sum_ef))
assert horizonte['pagos_despensa_3_en_horizonte'] == 3776.04
ipm = [dict(monto=N, intereses=f(D(str(N))*R), total=f(D(str(N))*(1+R)), pago_quincenal=f(D(str(N))*(1+R)/4)) for N in (1,100,1000,2000,3094.5,4000,5000)]

def pj(x): return {k: f(v) for k, v in x.items()}
def desp_json(d):
    return dict(n=d['n'], fecha=iso(d['fecha']), capital=f(d['capital']), intereses=f(d['intereses']), total=f(d['total']),
        pagos=[dict(fecha=iso(p['fecha']), monto=f(p['monto']), axel=f(p['axel']), karen=f(p['karen']), capital_repuesto=f(p['capital_repuesto'])) for p in d['pagos']],
        axel=pj(d['axel']), karen=pj(d['karen']))
def per_json(p):
    ef = {k: (f(v) if k not in ('n_pagos','detalle_pagos') else v) for k, v in p['ef'].items()}
    return dict(periodo=p['periodo'], etiqueta=p['etiqueta'], inicio=iso(p['inicio']), fin=iso(p['fin']), despensa_que_cubre=p['despensa'],
                devengado=pj(p['dev']), efectivo=ef, diferencia_acumulada=f(p['dif']))
model = dict(version='v3',
    parametros=dict(tasa=float(R), tasa_pct=float(R*100), papel=f(PAPEL), limite_credito=f(LIMITE), horizonte=[iso(START), iso(END)],
                    karen_fija=dict(capital=f(K_CAP), intereses=f(K_INT), total=f(K_TOT), pago=f(K_PAGO), conjunto=f(CONJ))),
    periodos=[per_json(p) for p in periodos],
    despensas=[desp_json(desp[i]) for i in (1,2,3,4)],
    credito=dict(eventos=cred, comparacion_series=cred_cmp, minimo_A=min(c['disponible_A'] for c in cred), minimo_B=min(c['disponible_B'] for c in cred),
                 nota="Solo capital usa/repone la linea; intereses nunca regresan; papel de contado no mueve el credito. A: 10,000 libres hoy y los pagos suman encima; B: 10,000 es tope total (inicial 6,905.50)"),
    valor_estable=ve, sensibilidad=sens, papel=papel, disposicion_por_despensa=disp, horizonte=horizonte, interes_por_monto=ipm)
if __name__ == '__main__':
    # S2-04: importar este módulo no escribe nada; ejecutado, regenera en datos/tmp/model.json (ruta relativa
    # a la raíz del repositorio) y lo compara con oraculo/model.json, que es canon y nunca se reescribe.
    _raiz = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    out = os.path.join(_raiz, 'datos', 'tmp', 'model.json')
    os.makedirs(os.path.dirname(out), exist_ok=True)
    with open(out, 'w', encoding='utf-8') as _fh:
        json.dump(model, _fh, ensure_ascii=False, indent=1)
    with open(out, 'rb') as _a, open(os.path.join(_raiz, 'oraculo', 'model.json'), 'rb') as _b:
        print('datos/tmp/model.json igual a oraculo/model.json:', _a.read() == _b.read())
    print('OK asserts')
    for p in model['periodos']: print(p['periodo'], p['devengado'], p['efectivo']['total'], p['diferencia_acumulada'])
    print(ve['quincena']); print(horizonte)
