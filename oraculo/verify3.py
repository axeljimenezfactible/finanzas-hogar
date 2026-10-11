import json, datetime as dt
from decimal import Decimal as D, ROUND_HALF_UP
def q(x): return float(D(str(x)).quantize(D('0.01'), ROUND_HALF_UP))
d = dt.date
r = 800.50/3094.50
# despensas: (compra, pagos total, pagos karen, pagos axel, capital, fechas)
K = 374.93
pay = []  # (date, desp, total, karen, axel, capital)
def add(desp, fechas, tot, cap):
    cp = [q(cap/4)]*3
    cp = [cap/4]*4
    for f, t in zip(fechas, tot):
        pay.append((f, desp, t, K if desp > 1 else None, None, None))
D1 = [d(2026,10,30), d(2026,11,15), d(2026,11,30), d(2026,12,15)]
D2 = [d(2026,12,30), d(2027,1,15), d(2027,1,30), d(2027,2,15)]
D3 = [d(2027,2,28), d(2027,3,15), d(2027,3,30), d(2027,4,15)]
D4 = [d(2027,4,30), d(2027,5,15), d(2027,5,30), d(2027,6,15)]
P = []
def ins(desp, fechas, tot, kar, ax, cap):
    for i, f in enumerate(fechas):
        P.append(dict(f=f, desp=desp, tot=tot[i], k=kar[i], a=ax[i], cap=cap[i]))
ins(1, D1, [973.75]*4, [374.93]*4, [598.82]*4, [773.63]*3+[773.61])
for n, F in ((2, D2), (3, D3), (4, D4)):
    ins(n, F, [1258.68]*3+[1258.70], [K]*4, [883.75]*3+[883.77], [1000.0]*4)
periods = [(d(2026,10,10), d(2026,11,9)), (d(2026,11,10), d(2026,12,9)), (d(2026,12,10), d(2027,1,9)),
           (d(2027,1,10), d(2027,2,9)), (d(2027,2,10), d(2027,3,9)), (d(2027,3,10), d(2027,4,9))]
papel_cont = [0, 0, 0, 199.0, 0, 199.0]
# devengado
dev_desp = {1: (1197.64, 749.86, [0, 1]), 2: (1767.51, 749.86, [2, 3]), 3: (1767.51, 749.86, [4, 5])}
dev_a = [0.0]*6; dev_k = [0.0]*6
for n, (a, k, ps) in dev_desp.items():
    for p in ps: dev_a[p] += a; dev_k[p] += k
for i in range(6): dev_k[i] += papel_cont[i]
dev_t = [q(a+k) for a, k in zip(dev_a, dev_k)]
dev_a = [q(x) for x in dev_a]; dev_k = [q(x) for x in dev_k]
# efectivo
ef_a, ef_k, ef_t = [], [], []
for i, (s, e) in enumerate(periods):
    ps = [p for p in P if s <= p['f'] <= e]
    a = q(sum(p['a'] for p in ps)); k = q(sum(p['k'] for p in ps) + papel_cont[i])
    ef_a.append(a); ef_k.append(k); ef_t.append(q(sum(p['tot'] for p in ps) + papel_cont[i]))
def cum(x):
    o, s = [], 0
    for v in x: s += v; o.append(q(s))
    return o
dif_t = [q(a-b) for a, b in zip(cum(dev_t), cum(ef_t))]
dif_a = [q(a-b) for a, b in zip(cum(dev_a), cum(ef_a))]
dif_k = [q(a-b) for a, b in zip(cum(dev_k), cum(ef_k))]
# credito
events = []
for p in P:
    if p['f'] <= d(2027,4,10): events.append((p['f'], 1, 'pago', p['cap']))
for f, c in ((d(2026,10,10), 3094.50), (d(2026,12,10), 4000.0), (d(2027,2,10), 4000.0), (d(2027,4,10), 4000.0)):
    events.append((f, 0, 'compra', c))
events.sort(key=lambda e: (e[0], e[1]))
A, B = [], []
# credito A
a = 10000.0 - 3094.50  # tras compra 10-oct (no se lista) -> hoy 10,000 libres: la serie parte con A=10000 tras compra?
# Lectura A: hoy hay 10000 libres (ya descontando la compra actual)
a = 10000.0; b = 10000.0 - 3094.50
for f, o, t, v in events:
    if t == 'pago':
        a += v; b += v
    else:
        if f == d(2026,10,10): continue  # compra actual ya descontada
        a -= v; b -= v
    A.append((str(f), q(a))); B.append((str(f), q(b)))
# estable
est_mes = dict(total=q(2517.37+99.50), karen=q(749.86+99.50), axel=1767.51)
est_q = dict(total=q(est_mes['total']/2), karen=q(est_mes['karen']/2), axel=q(est_mes['axel']/2))
# sensibilidad
sens = []
for N in (3000, 3500, 4000, 4500, 5000):
    tot = N*(1+r)
    sens.append(dict(N=N, total=q(tot), axel_total=q(tot-1499.72), mensual=q(tot/2+99.50), axel_mensual=q((tot-1499.72)/2), karen_mensual=849.36))
# intereses por monto
inte = []
for N in (1, 100, 1000, 2000, 3094.50, 4000, 5000):
    i = N*r; inte.append(dict(N=N, int=q(i), total=q(N+i), quinc=q((N+i)/4)))
# horizonte
hor = d(2027,4,10)
pagos_h = q(sum(p['tot'] for p in P if p['f'] <= hor))
efectivo_h = q(pagos_h + 398.00)
deuda = q(sum(p['tot'] for p in P if p['f'] > hor and p['desp'] in (3, 4)))
out = dict(dev_t=dev_t, dev_a=dev_a, dev_k=dev_k, ef_t=ef_t, ef_a=ef_a, ef_k=ef_k, dif_t=dif_t, dif_a=dif_a, dif_k=dif_k,
           A=A, B=B, est_mes=est_mes, est_q=est_q, sens=sens, interes=inte, pagos_h=pagos_h, efectivo_h=efectivo_h,
           suma_periodos=q(sum(ef_t)), deuda=deuda,
           pagos=[dict(f=str(p['f']), desp=p['desp'], tot=p['tot'], k=p['k'], a=p['a']) for p in P])
json.dump(out, open('/home/deepsec/reporte-compra-latex/model/verify3.json', 'w'), indent=1)
if __name__ == '__main__':
    for k in ('dev_t','ef_t','ef_a','ef_k','dif_t','pagos_h','efectivo_h','suma_periodos','deuda','est_mes','est_q'):
        print(k, out[k])
    print('A', [x[1] for x in A]); print('B', [x[1] for x in B])
    print(sens); print(inte)
