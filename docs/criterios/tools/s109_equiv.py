#!/usr/bin/env python3
"""Equivalencia P1-P6 del oraculo = meses calendario oct-2026..mar-2027 (S1-09, texto nuevo).
Lee solo oraculo/model.json (no escribe). Dinero en centavos enteros."""
import json, datetime as dt
from decimal import Decimal as D, ROUND_HALF_UP
M = json.load(open('/home/deepsec/finanzas-hogar/oraculo/model.json'))
c = lambda x: int((D(str(x)) * 100).quantize(D(1), rounding=ROUND_HALF_UP))
F = dt.date.fromisoformat
# eventos fechados de la semilla
pagos = [(F(p['fecha']), c(p['monto'])) for d in M['despensas'] for p in d['pagos']]
papel_contado = [(F(x['compra']), c(x['monto'])) for x in M['papel'] if 'contado' in x['forma']]
compras = [F(d['fecha']) for d in M['despensas']] if 'fecha' in M['despensas'][0] else []
eventos = [e[0] for e in pagos + papel_contado] + compras + [F(x['compra']) for x in M['papel']]
print('eventos fechados:', len(eventos), '| entre los dias 1 y 9 de un mes:', sum(1 for e in eventos if e.day <= 9))
meses = [(2026, 10), (2026, 11), (2026, 12), (2027, 1), (2027, 2), (2027, 3)]
dif = 0; n = 0
for i, (y, m) in enumerate(meses):
    P = M['periodos'][i]
    ef = sum(v for d, v in pagos if (d.year, d.month) == (y, m)) + sum(v for d, v in papel_contado if (d.year, d.month) == (y, m))
    papel_m = sum(v for d, v in papel_contado if (d.year, d.month) == (y, m))
    desp = M['despensas'][i // 2]              # cada despensa cubre 2 meses (bimestral)
    dev = c(D(str(desp['total'])) / 2) + papel_m
    for nombre, mio, ora in (('efectivo', ef, c(P['efectivo']['total'])), ('devengado', dev, c(P['devengado']['total']))):
        n += 1
        ok = mio == ora
        dif += (not ok)
        print(f"  {P['periodo']} <-> {y}-{m:02d} {nombre:9s} calendario={mio/100:>9,.2f} oraculo={ora/100:>9,.2f} {'OK' if ok else 'DIFIERE'}")
print(f'cifras comparadas: {n} | diferencias: {dif}')
