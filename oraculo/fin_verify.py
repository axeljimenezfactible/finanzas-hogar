#!/usr/bin/env python3
"""Segunda vía del oráculo S2 en float (patrón de verify3.py): recalcula las 159 hojas de S2 desde
_datos/semilla-oraculo.json agrupando por mes con aritmética float en pesos y redondeo a 2 decimales, y las compara
con la vía exacta (oraculo/salidas/s2.json, Fraction/enteros) y con model.json. Solo lee; no escribe nada.
Uso: python3 -I oraculo/fin_verify.py   (antes: python3 -I oraculo/fin_run_fixtures.py)
"""
import datetime as dt
import json
import re
import sys
from decimal import ROUND_HALF_UP, Decimal
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[1]
RAMAS = ('periodos', 'valor_estable', 'horizonte', 'papel')
CONTEOS = {'n', 'n_pagos', 'despensa', 'despensa_que_cubre', 'paquete'}


def q(x):
    return float(Decimal(repr(x)).quantize(Decimal('0.01'), ROUND_HALF_UP))


def pasos(r):
    return [(n if n else int(i)) for n, i in re.findall(r'([^.\[\]]+)|\[(\d+)\]', r)]


def nodo(x, r):
    for k in pasos(r):
        x = x[k]
    return x


def numericas(x, r):
    if isinstance(x, dict):
        for k, v in x.items():
            yield from numericas(v, f'{r}.{k}')
    elif isinstance(x, list):
        for i, v in enumerate(x):
            yield from numericas(v, f'{r}[{i}]')
    elif isinstance(x, (int, float)) and not isinstance(x, bool):
        yield r, x


def calcula(s):
    f = dt.date.fromisoformat
    gastos = {g['id']: g for g in s['expenses']}
    compras = []
    for c in s['credit_purchases']:
        g = gastos[c['expense_id']]
        tot = c['total_cents'] / 100
        kar = sum(x['total_cents'] for x in c['shares'] if x['person'] == 'karen') / 100
        compras.append(dict(id=c['id'], fecha=f(g['expense_date']), tot=tot, kar=kar, emb=c['embedded'],
                            pagos=[dict(f=f(p['due_date']), m=p['amount_cents'] / 100,
                                        a=sum(x['amount_cents'] for x in p['shares'] if x['person'] == 'axel') / 100,
                                        k=sum(x['amount_cents'] for x in p['shares'] if x['person'] == 'karen') / 100) for p in c['payments']]))
    compras.sort(key=lambda c: (c['fecha'], c['id']))
    for n, c in enumerate(compras):
        c['seq'] = n + 1
    papel = [dict(f=f(g['expense_date']), m=g['amount_cents'] / 100, quien=g['paid_by']) for g in s['expenses']
             if g['payment_method'] != 'credit' and g['category'] == s['params']['paper_category']]
    meses = [(2026, 10), (2026, 11), (2026, 12), (2027, 1), (2027, 2), (2027, 3)]
    per, ca, ce = [], 0.0, 0.0
    for (y, m) in meses:
        dev_d = dev_k = 0.0
        cubre = None
        for c in compras:
            idx = (y * 12 + m) - (c['fecha'].year * 12 + c['fecha'].month)
            if idx in (0, 1):
                h, hk = q(c['tot'] / 2), q(c['kar'] / 2)
                if idx == 1:
                    h, hk = c['tot'] - h, c['kar'] - hk
                dev_d += h
                dev_k += hk
                cubre = max(cubre or 0, c['seq'])
        pap = sum(p['m'] for p in papel if (p['f'].year, p['f'].month) == (y, m))
        pk = sum(p['m'] for p in papel if (p['f'].year, p['f'].month) == (y, m) and p['quien'] == 'karen')
        pagos = sorted([(p['f'], c['id'], c['seq'], p) for c in compras for p in c['pagos'] if (p['f'].year, p['f'].month) == (y, m)], key=lambda t: t[:2])
        ps = sum(p['m'] for *_, p in pagos)
        dev = dict(despensa_total=q(dev_d), papel_contado=q(pap), total=q(dev_d + pap), axel=q(dev_d - dev_k + pap - pk), karen=q(dev_k + pk))
        ef = dict(pagos=q(ps), papel=q(pap), total=q(ps + pap), axel=q(sum(p['a'] for *_, p in pagos) + pap - pk),
                  karen=q(sum(p['k'] for *_, p in pagos) + pk), n_pagos=len(pagos),
                  detalle_pagos=[dict(despensa=sq, monto=q(p['m']), axel=q(p['a']), karen=q(p['k'])) for _, _, sq, p in pagos])
        ca += dev['total']
        ce += ef['total']
        per.append(dict(despensa_que_cubre=cubre, devengado=dev, efectivo=ef, diferencia_acumulada=q(ca - ce)))
    comp = {}
    for it in s['stable_items']:
        t = q(it['total_cents'] / 100 / it['cycle_months'])
        k = q(sum(x['total_cents'] for x in it['shares'] if x['person'] == 'karen') / 100 / it['cycle_months'])
        comp[it['label']] = (t, k, q(t - k))
    mt = q(sum(v[0] for v in comp.values()))
    mk = q(sum(v[1] for v in comp.values()))
    ma = q(sum(v[2] for v in comp.values()))
    qa, qk = q(ma / 2), q(mk / 2)
    ve = dict(mes=dict(total=mt, axel=ma, karen=mk),
              quincena=dict(total=q(mt / 2), axel=qa, karen=qk, suma_axel_karen=q(qa + qk), diferencia_centavo=q(q(mt / 2) - qa - qk)),
              componentes=dict(mitad_despensa=comp['despensa'][0], papel_contado_mes=comp['papel_contado'][0],
                               karen_mitad_despensa=comp['despensa'][1], axel_mitad_despensa=comp['despensa'][2]))
    ini, fin = dt.date(2026, 10, 10), dt.date(2027, 4, 10)
    hor = {}
    for c in compras:
        dentro = [p['m'] for p in c['pagos'] if ini <= p['f'] <= fin]
        if dentro:
            hor[f"pagos_despensa_{c['seq']}" + ('' if len(dentro) == len(c['pagos']) else '_en_horizonte')] = q(sum(dentro))
    te = q(sum(p['m'] for c in compras for p in c['pagos'] if ini <= p['f'] <= fin) + sum(p['m'] for p in papel if ini <= p['f'] <= fin))
    sp = q(sum(x['efectivo']['total'] for x in per))
    td = q(sum(x['devengado']['total'] for x in per))
    det = {}
    for c in compras:
        pend = [p['m'] for p in c['pagos'] if p['f'] > fin]
        if c['fecha'] <= fin and pend:
            det[f"despensa_{c['seq']}" + ('_pendiente' if len(pend) < len(c['pagos']) else '')] = q(sum(pend))
    hor.update(papel_contado=q(sum(p['m'] for p in papel if ini <= p['f'] <= fin)), total_efectivo=te, suma_periodos_efectivo=sp,
               total_devengado=td, deuda_viva_10_abr=q(sum(det.values())), deuda_viva_detalle=det,
               promedio_mensual_efectivo=q(te / 6), promedio_mensual_devengado=q(td / 6), diferencia_final=q(td - te))
    paq = [(c['fecha'], 0, c['id'], e['amount_cents'] / 100) for c in compras for e in c['emb'] if e['category'] == s['params']['paper_category']]
    paq += [(p['f'], 1, 0, p['m']) for p in papel]
    pap = [dict(paquete=i + 1, monto=q(x[3])) for i, x in enumerate(sorted(paq))]
    return dict(periodos=per, valor_estable=ve, horizonte=hor, papel=pap)


def caso_de(r):
    k = pasos(r)
    if k[0] == 'periodos':
        return f's2-periodos-p{k[1] + 1}', f'periodos[{k[1]}]'
    return {'valor_estable': 's2-valor-estable', 'horizonte': 's2-horizonte', 'papel': 's2-papel'}[k[0]], k[0]


def main():
    with open(RAIZ / 'oraculo' / 'fixtures' / 's2' / '_datos' / 'semilla-oraculo.json', encoding='utf-8') as fh:
        s = json.load(fh)
    with open(RAIZ / 'oraculo' / 'model.json', encoding='utf-8') as fh:
        modelo = json.load(fh)
    with open(RAIZ / 'oraculo' / 'salidas' / 's2.json', encoding='utf-8') as fh:
        sal = json.load(fh)
    fv = calcula(s)
    hojas = [(r, v) for rama in RAMAS for r, v in numericas(modelo[rama], rama)]
    ok_ex, ok_mod, difs = 0, 0, []
    for r, mv in hojas:
        conteo = pasos(r)[-1] in CONTEOS
        try:
            v = nodo(fv, r)
        except (KeyError, IndexError):
            v = None
        caso, i = caso_de(r)
        ex = nodo(sal['casos'][caso][i], r[len(i):]) if r != i else sal['casos'][caso][i]
        exv = ex if conteo else (ex / 100 if isinstance(ex, int) else None)
        a = v is not None and exv is not None and (v == exv if conteo else abs(round(v, 2) - exv) < 0.005)
        b = v is not None and (v == mv if conteo else abs(round(v, 2) - mv) < 0.005)
        ok_ex += a
        ok_mod += b
        if not (a and b):
            difs.append(f'DIF {r} float={v} exacta={exv} model.json={mv}')
    print(f'fin_verify: {ok_ex} de {len(hojas)} hojas de S2 coinciden a 2 decimales (float contra la vía exacta)')
    print(f'fin_verify: {ok_mod} de {len(hojas)} hojas de S2 coinciden a 2 decimales (float contra model.json)')
    for d in difs:
        print(d)
    return 0 if not difs else 1


if __name__ == '__main__':
    sys.exit(main())
