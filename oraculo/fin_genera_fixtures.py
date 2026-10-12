#!/usr/bin/env python3
"""Genera oraculo/fixtures/s2/ (27 fixtures y 3 archivos de _datos/) de forma determinista (S2-interfaz.md §5, §5.0–§5.3).

Lee oraculo/model.json y datos/semilla/data.json (FORMATO §3). Consulta con id ruta: esperado = valor convertido de
model.json. Consulta de nombre libre: cifras fijadas por la interfaz o los criterios (citadas en `fuente`) y el resto
calculado por el oráculo. Si el oráculo difiere de una cifra fijada, se escribe la fijada y se reporta HALLAZGO (protocolo §6).
Uso: python3 -I oraculo/fin_genera_fixtures.py
"""
import datetime as dt
import json
import re
import sys
from decimal import Decimal
from fractions import Fraction
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(Path(__file__).resolve().parent))
import fin_run_fixtures as rf  # noqa: E402
from fin_comun import hu  # noqa: E402
import model as modelo_v3  # noqa: E402  (sched: calendario de pagos del oráculo; importarlo ya no escribe, S2-04)

CARPETA = RAIZ / 'oraculo' / 'fixtures' / 's2'
with open(RAIZ / 'oraculo' / 'model.json', encoding='utf-8') as fh:
    M = json.load(fh, parse_float=Decimal)
with open(RAIZ / 'datos' / 'semilla' / 'data.json', encoding='utf-8') as fh:
    DJ = json.load(fh, parse_float=Decimal)
CONTEOS = {'n', 'n_pagos', 'despensa', 'despensa_que_cubre', 'paquete'}
TEXTO = {'fecha', 'compra', 'quien', 'forma'}
MESES6 = ['2026-10', '2026-11', '2026-12', '2027-01', '2027-02', '2027-03']
P = ['axel', 'karen']
HALLAZGOS = []
SEMILLA = {'$ref': '_datos/semilla-oraculo.json'}


def c(x):
    v = Decimal(x) * 100
    if v != v.to_integral_value():
        raise SystemExit(f'no es centavo exacto: {x}')
    return int(v)


def pasos(ruta):
    return [(n if n else int(i)) for n, i in re.findall(r'([^.\[\]]+)|\[(\d+)\]', ruta)]


def nodo(raiz, ruta):
    x = raiz
    for k in pasos(ruta):
        x = x[k]
    return x


def convierte(x):
    if isinstance(x, list):
        return [convierte(v) for v in x]
    out = {}
    for k, v in x.items():
        if isinstance(v, (dict, list)):
            out[k] = convierte(v)
        elif isinstance(v, str):
            if k in TEXTO:
                out[k] = v
        elif isinstance(v, bool) or v is None:
            continue
        else:
            out[k] = int(v) if k in CONTEOS else c(v)
    return out


def escribe(nombre, contenido):
    p = CARPETA / nombre
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(json.dumps(contenido, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')


def fx(caso, descripcion, fuente, consultas, fijos=None, modelo=None, datos=None, provisional=None, recorte=None):
    datos = datos or SEMILLA
    entrada = rf.resuelve_datos(datos, CARPETA)
    esperado = {}
    for qid, fn, args in consultas:
        if modelo and qid in modelo:
            esperado[qid] = modelo[qid]
            continue
        calc = rf.calcula(entrada, {'id': qid, 'funcion': fn, 'argumentos': args})
        if recorte and qid in recorte:
            calc = recorte[qid](calc)
        for ruta, v in (fijos or {}).get(qid, {}).items():
            ks = pasos(ruta)
            try:
                padre = nodo(calc, ruta.rsplit('.', 1)[0]) if '.' in ruta else calc
                actual = nodo(calc, ruta)
            except (KeyError, IndexError, TypeError):
                HALLAZGOS.append(f'{caso} {qid}.{ruta}: fijado={v!r} oraculo=__falta__')
                continue
            if not rf.igual(actual, v) if not isinstance(v, (list, dict)) else actual != v:
                HALLAZGOS.append(f'{caso} {qid}.{ruta}: fijado={v!r} oraculo={actual!r}')
            padre[ks[-1]] = v
        esperado[qid] = calc
    escribe(f'{caso}.json', {'formato': 'fixture/1', 'motor': 's2', 'caso': caso, 'descripcion': descripcion, 'fuente': fuente,
                             'provisional': provisional,
                             'entradas': {'datos': datos, 'consultas': [{'id': q, 'funcion': f, 'argumentos': a} for q, f, a in consultas]},
                             'esperado': esperado})


def F(tipo, ruta):
    return {'tipo': tipo, 'ruta': ruta}


# ---- _datos/semilla-oraculo.json (§5.3) ------------------------------------
D, PA, PAR = M['despensas'], M['papel'], M['parametros']
PAPEL = c(PAR['papel'])
PERSONAS = [{'code': 'axel', 'display_name': 'Axel', 'participation': 'holder', 'pays_share_directly': False},
            {'code': 'karen', 'display_name': 'Karen', 'participation': 'external', 'pays_share_directly': True}]
CATS = [{'code': 'despensa', 'name': 'Despensa', 'kind': 'variable', 'is_active': True},
        {'code': 'papel', 'name': 'Papel', 'kind': 'variable', 'is_active': True}]
PARAMS = {'today': '2026-10-11', 'holder': 'axel', 'groceries_category': 'despensa', 'paper_category': 'papel', 'repeat_cadence_enabled': False}


def gasto_despensa(k, d, ident, cpid, proyectado=False):
    g = {'id': ident, 'kind': 'variable', 'category': 'despensa', 'amount_cents': c(d['capital']), 'expense_date': d['fecha'],
         'accrual_months': 2, 'payment_method': 'credit', 'paid_by': None,
         'shares': [{'person': p, 'amount_cents': c(d[p]['capital'])} for p in P], 'credit_purchase_id': cpid, 'planned_item_id': None}
    if proyectado:
        g['projected'] = True
    return g


def compra(d, ident, gid, fechas=None, proyectado=False):
    pagos = []
    for i, x in enumerate(d['pagos']):
        pagos.append({'seq': i + 1, 'due_date': fechas[i] if fechas else x['fecha'], 'amount_cents': c(x['monto']),
                      'principal_cents': c(x['capital_repuesto']), 'interest_cents': c(x['monto']) - c(x['capital_repuesto']),
                      'shares': [{'person': p, 'amount_cents': c(x[p])} for p in P], 'paid': []})
    return {'id': ident, 'credit_id': 1, 'expense_id': gid, 'purchase_date': d['fecha'] if not fechas else d['_fecha'],
            'principal_cents': c(d['capital']), 'interest_cents': c(d['intereses']), 'total_cents': c(d['total']), 'status': 'active',
            'shares': [{'person': p, 'principal_cents': c(d[p]['capital']), 'interest_cents': c(d[p]['intereses']), 'total_cents': c(d[p]['total'])} for p in P],
            'embedded': [{'category': 'papel', 'person': p, 'amount_cents': c(d[p]['papel_dentro'])} for p in P if d[p]['papel_dentro'] > 0],
            'payments': pagos, 'projected': proyectado}


def papel_contado(ident, i):
    q = PA[i]['quien'].lower()
    return {'id': ident, 'kind': 'variable', 'category': 'papel', 'amount_cents': PAPEL, 'expense_date': PA[i]['compra'], 'accrual_months': 1,
            'payment_method': 'cash', 'paid_by': q, 'shares': [{'person': q, 'amount_cents': PAPEL}], 'credit_purchase_id': None, 'planned_item_id': None}


GASTOS = [gasto_despensa(k, D[k], k + 1, k + 1) for k in range(4)] + [papel_contado(5, 2), papel_contado(6, 4)]
semilla = {'params': PARAMS, 'persons': PERSONAS, 'categories': CATS, 'expenses': GASTOS,
           'credit_purchases': [compra(D[k], k + 1, k + 1) for k in range(4)],
           'incomes': [], 'planned_items': [], 'budgets': [], 'savings_goals': [], 'savings_movements': [],
           'cadence': {'purchases': [], 'expenses': []},
           'stable_items': [{'label': 'despensa', 'cycle_months': 2, 'total_cents': c(D[1]['total']),
                             'shares': [{'person': 'karen', 'total_cents': c(D[1]['karen']['total'])}, {'person': 'axel', 'total_cents': c(D[1]['axel']['total'])}]},
                            {'label': 'papel_contado', 'cycle_months': 2, 'total_cents': PAPEL, 'shares': [{'person': 'karen', 'total_cents': PAPEL}]}]}
escribe('_datos/semilla-oraculo.json', semilla)

# ---- _datos/insumo-s3.json (§5.1) -------------------------------------------
filas = {}
def fila(r, llaves):
    filas[r] = {}
    for ll in llaves:
        if ll.startswith('data.json:'):
            v = nodo(DJ, ll[len('data.json:'):])
        else:
            v = nodo(M, ll)
        k = ll.rsplit('.', 1)[-1]
        if k == 'tasa_pct' and ll.startswith('data.json:'):
            filas[r][ll] = {'num': c(DJ['financiamiento']['dos_quincenas']['intereses']), 'den': c(D[0]['capital']), 'decimales': 2}
        elif k == 'tasa_pct':
            filas[r][ll] = {'num': c(D[0]['intereses']), 'den': c(D[0]['capital'])}
        elif isinstance(v, list):
            filas[r][ll] = [c(x) for x in v]
        elif isinstance(v, str):
            filas[r][ll] = v
        else:
            filas[r][ll] = int(v) if k in CONTEOS else c(v)
fila('R-01', ['despensas[0].capital']); fila('R-02', ['despensas[0].intereses']); fila('R-03', ['despensas[0].total'])
fila('R-04', ['parametros.tasa_pct']); fila('R-05', ['data.json:financiamiento.dos_quincenas.tasa_pct'])
fila('R-06', ['data.json:financiamiento.dos_quincenas.total', 'data.json:financiamiento.dos_quincenas.intereses'])
for r, (p, k) in zip(['R-07', 'R-08', 'R-09', 'R-10', 'R-11', 'R-12'], [(p, k) for p in P for k in ('total', 'pago_por_quincena', 'intereses')]):
    fila(r, [f'despensas[0].{p}.{k}'])
for i in range(4):
    fila(f'R-{13 + i}', [f'despensas[0].pagos[{i}].{k}' for k in ('fecha', 'monto', 'axel', 'karen', 'capital_repuesto')])
for k, r0 in ((1, 17), (2, 22), (3, 27)):
    fila(f'R-{r0}', [f'despensas[{k}].{x}' for x in ('capital', 'intereses', 'total')])
    for i in range(4):
        fila(f'R-{r0 + 1 + i}', [f'despensas[{k}].pagos[{i}].{x}' for x in ('fecha', 'monto', 'axel', 'karen')])
for r, k in zip(['R-32', 'R-33', 'R-34', 'R-35', 'R-36'], ['total', 'pago', 'capital', 'intereses', 'conjunto']):
    fila(r, [f'parametros.karen_fija.{k}'])
fila('R-61', ['credito.minimo_A']); fila('R-62', ['credito.minimo_B']); fila('R-63', ['parametros.limite_credito'])
fila('R-64', ['credito.eventos[0].disponible_A', 'credito.eventos[0].disponible_B'])
fila('R-65', ['credito.comparacion_series.esperado_A']); fila('R-66', ['credito.comparacion_series.esperado_B'])
fila('R-67', [f'credito.eventos[{i}].{x}' for i in range(17) for x in ('fecha', 'tipo', 'monto', 'disponible_A', 'disponible_B')])
for k in range(5):
    fila(f'R-{74 + k}', [f'sensibilidad[{k}].{x}' for x in ('intereses', 'total', 'pago_quincenal', 'axel_total')])
for k in range(7):
    fila(f'R-{79 + k}', [f'interes_por_monto[{k}].{x}' for x in ('intereses', 'total', 'pago_quincenal')])
assert len(filas) == 55 and len(M['credito']['eventos']) == 17
tp = float(Fraction(c(D[0]['intereses']) * 100, c(D[0]['capital'])))
if tp != float(M['parametros']['tasa_pct']):
    HALLAZGOS.append(f"R-04: binary64(100·80050/309450)={tp!r} model.json={M['parametros']['tasa_pct']}")
r5 = hu(Fraction(c(DJ['financiamiento']['dos_quincenas']['intereses']) * 10000, c(D[0]['capital'])))
if Decimal(r5) / 100 != DJ['financiamiento']['dos_quincenas']['tasa_pct']:
    HALLAZGOS.append(f'R-05: HALF_UP={r5} data.json={DJ["financiamiento"]["dos_quincenas"]["tasa_pct"]}')
escribe('_datos/insumo-s3.json', {'formato': 'insumo-s3/1', 'descripcion': 'Filas S3 de la regresión tal como S3 las entrega; S2 no las calcula', 'filas': filas})

# ---- _datos/interes-por-monto.json (§5.2) ----------------------------------
fechas_d1 = [d.isoformat() for d in modelo_v3.sched(dt.date(2026, 10, 10))]
assert fechas_d1 == [x['fecha'] for x in D[0]['pagos']]
ipm = []
for k, x in enumerate(M['interes_por_monto']):
    cap, tot, pq = c(x['monto']), c(x['total']), c(x['pago_quincenal'])
    montos = [pq] * 3 + [tot - 3 * pq]
    base = hu(Fraction(cap, 4))
    caps = [base] * 3 + [cap - 3 * base]
    ipm.append({'id': k + 1, 'credit_id': 1, 'expense_id': None, 'category': 'despensa', 'purchase_date': '2026-10-10',
                'principal_cents': cap, 'interest_cents': c(x['intereses']), 'total_cents': tot, 'accrual_months': 2, 'status': 'active',
                'shares': [{'person': 'axel', 'principal_cents': cap, 'interest_cents': c(x['intereses']), 'total_cents': tot}], 'embedded': [],
                'payments': [{'seq': i + 1, 'due_date': fechas_d1[i], 'amount_cents': montos[i], 'principal_cents': caps[i],
                              'interest_cents': montos[i] - caps[i], 'shares': [{'person': 'axel', 'amount_cents': montos[i]}], 'paid': []} for i in range(4)],
                'projected': False})
escribe('_datos/interes-por-monto.json', {**semilla, 'expenses': [], 'credit_purchases': ipm, 'stable_items': []})

# ---- casos canon (ids = rutas de model.json) --------------------------------
for k in range(6):
    fx(f's2-periodos-p{k + 1}', f'Mes calendario {MESES6[k]} con la semilla; reproduce P{k + 1} (S2-09, nota (e) de S2.md §6)',
       [F('model.json', f'periodos[{k}]'), F('criterio', 'S2-01'), F('criterio', 'S2-09')],
       [(f'periodos[{k}]', 'vista_periodo', {'month': MESES6[k]})], modelo={f'periodos[{k}]': convierte(M['periodos'][k])})
fx('s2-valor-estable', 'Valor estable de la semilla (R-37…R-44)', [F('model.json', 'valor_estable'), F('criterio', 'S2-01')],
   [('valor_estable', 'vista_valor_estable', {})], modelo={'valor_estable': convierte(M['valor_estable'])})
fx('s2-horizonte', 'Horizonte 10-oct-2026 a 10-abr-2027 con los meses oct-2026…mar-2027 (R-51…R-60)',
   [F('model.json', 'horizonte'), F('criterio', 'S2-01')],
   [('horizonte', 'vista_horizonte', {'start': '2026-10-10', 'end': '2027-04-10', 'months': MESES6})], modelo={'horizonte': convierte(M['horizonte'])})
fx('s2-papel', 'Seis paquetes de papel: dentro de despensa y de contado (R-68…R-73)', [F('model.json', 'papel'), F('criterio', 'S2-01')],
   [('papel', 'vista_papel', {})], modelo={'papel': convierte(M['papel'])})
S = M['sensibilidad']
fx('s2-sensibilidad', 'Cruce con S3: valor estable con la despensa de N (5 filas de sensibilidad, 15 hojas)',
   [F('model.json', f'sensibilidad[{k}]') for k in range(5)] + [F('criterio', 'S2-03')],
   [(f'sensibilidad[{k}]', 'vista_sensibilidad', {'items': [
       {'label': 'despensa', 'cycle_months': 2, 'total_cents': c(S[k]['total']),
        'shares': [{'person': 'karen', 'total_cents': c(PAR['karen_fija']['total'])}, {'person': 'axel', 'total_cents': c(S[k]['axel_total'])}]},
       {'label': 'papel_contado', 'cycle_months': 2, 'total_cents': PAPEL, 'shares': [{'person': 'karen', 'total_cents': PAPEL}]}]}) for k in range(5)],
   modelo={f'sensibilidad[{k}]': {x: c(S[k][x]) for x in ('costo_mensual_total', 'axel_mensual', 'karen_mensual')} for k in range(5)})
fx('s2-interes-por-monto', 'Cruce con S3: 7 compras de interes_por_monto confirmadas por devengado y efectivo (28 hojas)',
   [F('model.json', f'interes_por_monto[{k}]') for k in range(7)] + [F('criterio', 'S2-03')],
   [(f'interes_por_monto[{k}]', 'vista_interes_por_monto', {'purchase_id': k + 1}) for k in range(7)],
   modelo={f'interes_por_monto[{k}]': convierte(M['interes_por_monto'][k]) for k in range(7)}, datos={'$ref': '_datos/interes-por-monto.json'})

# ---- casos de nombre libre ---------------------------------------------------
CIERRES = [('2026-10-31', 292125, 97375, 194750), ('2026-11-30', 97375, 97375, 0), ('2026-12-31', 377606, 125869, 251737),
           ('2027-01-31', 125870, 125870, 0), ('2027-02-28', 377606, 125869, 251737), ('2027-03-31', 125870, 125870, 0),
           ('2027-04-10', 629344, 125870, 503474)]
fx('s2-deuda-viva-identidad', 'Identidad de K19 en los 7 cierres (S2-interfaz §3.4); resto calculado por el oráculo',
   [F('criterio', 'S2-07'), F('decision', 'K19'), F('regla', 'RN-21'), F('model.json', 'horizonte.deuda_viva_10_abr')],
   [(f'cierre-{d}', 'deuda_viva', {'cut': d}) for d, *_ in CIERRES],
   fijos={f'cierre-{d}': {'live_debt_cents': v, 'cumulative_diff_cents': df, 'unaccrued_balance_cents': s, 'identity_holds': True,
                          'identity_without_balance_holds': s == 0} for d, v, df, s in CIERRES})
pg = lambda k, i, p: c(D[k]['pagos'][i][p])  # noqa: E731
fx('s2-fin-de-mes-feb-2027', 'Febrero de 2027: quincena 2 del 16 al 28, paga el 28 (R-23); resto calculado por el oráculo',
   [F('criterio', 'S2-06'), F('decision', 'K8'), F('model.json', 'despensas[2].pagos[0]'), F('model.json', 'periodos[4]')],
   [('feb-20', 'calendario', {'date': '2027-02-20'}), ('feb-2027', 'periodos', {'months': ['2027-02']}),
    ('flujo-feb-2027', 'flujo_quincenal', {'from': '2027-02-q1', 'count': 2})],
   fijos={'feb-20': {'fortnight': '2027-02-q2', 'fortnight_start': '2027-02-16', 'fortnight_end': '2027-02-28', 'pay_date': '2027-02-28', 'days_in_month': 28},
          'feb-2027': {'[0].days': 28, '[0].accrued.total_cents': c(M['periodos'][4]['devengado']['total']), '[0].cash.total_cents': c(M['periodos'][4]['efectivo']['total'])},
          'flujo-feb-2027': {'[0].credit_payments.household_cents': pg(1, 3, 'axel'), '[1].end': '2027-02-28', '[1].pay_date': '2027-02-28',
                             '[1].credit_payments.household_cents': pg(2, 0, 'axel')}})
fx('s2-mes-30-dias', 'Noviembre de 2026 (30 días): quincena 2 del 16 al 30, paga el 30; resto calculado por el oráculo',
   [F('criterio', 'S2-06'), F('model.json', 'despensas[0].pagos[1]'), F('model.json', 'despensas[0].pagos[2]'), F('model.json', 'periodos[1]')],
   [('nov-30', 'calendario', {'date': '2026-11-30'}), ('nov-2026', 'periodos', {'months': ['2026-11']}),
    ('flujo-nov-2026', 'flujo_quincenal', {'from': '2026-11-q1', 'count': 2})],
   fijos={'nov-30': {'fortnight': '2026-11-q2', 'fortnight_start': '2026-11-16', 'fortnight_end': '2026-11-30', 'pay_date': '2026-11-30', 'days_in_month': 30},
          'nov-2026': {'[0].days': 30, '[0].cash.total_cents': c(M['periodos'][1]['efectivo']['total']), '[0].cash.payment_count': 2},
          'flujo-nov-2026': {'[0].credit_payments.household_cents': pg(0, 1, 'axel'), '[1].credit_payments.household_cents': pg(0, 2, 'axel'), '[1].end': '2026-11-30'}})
VARIOS = {'code': 'varios', 'name': 'Varios', 'kind': 'variable', 'is_active': True}


def gasto_fx(ident, cat, monto, fecha, kind='variable'):
    return {'id': ident, 'kind': kind, 'category': cat, 'amount_cents': monto, 'expense_date': fecha, 'accrual_months': 1,
            'payment_method': 'cash', 'paid_by': None, 'shares': [{'person': 'axel', 'amount_cents': monto}],
            'credit_purchase_id': None, 'planned_item_id': None}


fx('s2-mes-31-dias', 'Octubre de 2026 (31 días) con un gasto de fixture del 31-oct (50000, de contado, del titular) y hoy = 31-oct; resto calculado por el oráculo',
   [F('criterio', 'S2-06'), F('manual', 'gasto del 2026-10-31 calculado a mano'), F('model.json', 'despensas[0].pagos[0]')],
   [('oct-31', 'calendario', {'date': '2026-10-31'}), ('flujo-oct-2026', 'flujo_quincenal', {'from': '2026-10-q1', 'count': 2}),
    ('oct-2026', 'periodos', {'months': ['2026-10']})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {'params': {'today': '2026-10-31'}, 'categories': CATS + [VARIOS],
                                                           'expenses': GASTOS + [gasto_fx(7, 'varios', 50000, '2026-10-31')]}},
   fijos={'oct-31': {'fortnight': '2026-10-q2', 'fortnight_start': '2026-10-16', 'fortnight_end': '2026-10-31', 'pay_date': '2026-10-30', 'days_in_month': 31},
          'flujo-oct-2026': {'[0].variable_spent_cents': 0, '[0].spent_cents': 0, '[1].variable_spent_cents': 50000, '[1].spent_cents': 50000,
                             '[1].credit_payments.household_pending_cents': pg(0, 0, 'axel')},
          'oct-2026': {'[0].cash.cash_expenses_cents': 50000, '[0].cash.total_cents': 147375}})
fx('s2-semana', 'Semanas de lunes a domingo (K8, ISO 8601): 14-oct-2026, 10-oct-2026 y 1-ene-2027 (semana 53 de 2026); resto calculado por el oráculo',
   [F('criterio', 'S2-06'), F('regla', 'RN-06'), F('manual', 'semana ISO del 1-ene-2027')],
   [('mie-14-oct', 'calendario', {'date': '2026-10-14'}), ('sab-10-oct', 'calendario', {'date': '2026-10-10'}), ('vie-1-ene', 'calendario', {'date': '2027-01-01'})],
   fijos={'mie-14-oct': {'week': '2026-W42', 'week_start': '2026-10-12', 'week_end': '2026-10-18'},
          'sab-10-oct': {'week': '2026-W41', 'week_start': '2026-10-05', 'week_end': '2026-10-11'},
          'vie-1-ene': {'week': '2026-W53', 'week_start': '2026-12-28', 'week_end': '2027-01-03', 'fortnight': '2027-01-q1'}})
fx('s2-zona-medianoche', 'Hoy en America/Mexico_City (UTC−6 todo el año) cerca de la medianoche; un instante con desfase +00:00 o −06:00 '
   'en lugar de «Z» es invalid_date (§1.1, REVISION-S2-r3 B1); resto calculado por el oráculo',
   [F('criterio', 'S2-14'), F('manual', 'instante de verano sin horario de verano'), F('manual', 'S2-interfaz §1.1: instante RFC 3339 en UTC con «Z»')],
   [('antes', 'calendario', {'now': '2026-10-31T05:59:59Z'}), ('despues', 'calendario', {'now': '2026-10-31T06:00:00Z'}),
    ('verano', 'calendario', {'now': '2027-07-01T05:30:00Z'}), ('dos-formas', 'calendario', {'date': '2026-10-31', 'now': '2026-10-31T06:00:00Z'}),
    ('desfase-cero', 'calendario', {'now': '2026-10-31T05:59:59+00:00'}), ('desfase-cdmx', 'calendario', {'now': '2026-10-31T00:30:00-06:00'})],
   fijos={'antes': {'today': '2026-10-30', 'fortnight': '2026-10-q2'}, 'despues': {'today': '2026-10-31'}, 'verano': {'today': '2027-06-30'},
          'dos-formas': {'error': 'invalid_argument'}, 'desfase-cero': {'error': 'invalid_date'}, 'desfase-cdmx': {'error': 'invalid_date'}})
fx('s2-promedios', 'Promedios diario, semanal (K39) y mensual (R-56, R-57) con los meses calendario',
   [F('criterio', 'S2-08'), F('decision', 'K20'), F('decision', 'K39'), F('regla', 'RN-27'), F('model.json', 'horizonte.promedio_mensual_efectivo'), F('model.json', 'horizonte.promedio_mensual_devengado')],
   [('oct-2026-efectivo', 'promedios', {'month': '2026-10', 'basis': 'cash'}), ('nov-2026-efectivo', 'promedios', {'month': '2026-11', 'basis': 'cash'}),
    ('feb-2027-efectivo', 'promedios', {'month': '2027-02', 'basis': 'cash'}), ('oct-2026-devengado-despensa', 'promedios', {'month': '2026-10', 'basis': 'accrued_groceries'}),
    ('oct-2026-total-explicito', 'promedios', {'total_cents': 97375, 'start': '2026-10-01', 'end': '2026-10-31'}),
    ('mensual-efectivo', 'promedio_mensual', {'total_cents': c(M['horizonte']['total_efectivo']), 'month_count': 6}),
    ('mensual-devengado', 'promedio_mensual', {'total_cents': c(M['horizonte']['total_devengado']), 'month_count': 6})],
   fijos={'oct-2026-efectivo': {'days': 31, 'daily_cents': 3141, 'weekly_cents': 21988}, 'nov-2026-efectivo': {'days': 30, 'daily_cents': 6492, 'weekly_cents': 45442},
          'feb-2027-efectivo': {'days': 28, 'daily_cents': 8991, 'weekly_cents': 62935}, 'oct-2026-devengado-despensa': {'days': 31, 'daily_cents': 6282, 'weekly_cents': 43976},
          'oct-2026-total-explicito': {'days': 31, 'daily_cents': 3141, 'weekly_cents': 21988},
          'mensual-efectivo': {'monthly_cents': c(M['horizonte']['promedio_mensual_efectivo'])}, 'mensual-devengado': {'monthly_cents': c(M['horizonte']['promedio_mensual_devengado'])}})
VENT = [{'label': p['etiqueta'], 'start': p['inicio'], 'end': p['fin']} for p in M['periodos']]
fx('s2-ventanas-dia-5', 'Gasto de fixture del 5-nov-2026 (30000): con las ventanas del oráculo cae en P1 y con meses en noviembre; esperado del oráculo con las ventanas declaradas (S2-09)',
   [F('criterio', 'S2-09'), F('manual', 'gasto del 2026-11-05 calculado a mano')],
   [('ventanas', 'periodos', {'windows': VENT}), ('meses', 'periodos', {'months': ['2026-10', '2026-11']})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {'categories': CATS + [VARIOS], 'expenses': GASTOS + [gasto_fx(7, 'varios', 30000, '2026-11-05')]}},
   fijos={'ventanas': {'[0].accrued.total_cents': 224750, '[0].cash.total_cents': 127375, '[1].accrued.total_cents': 194750},
          'meses': {'[0].cash.total_cents': 97375, '[1].accrued.total_cents': 224750, '[1].cash.total_cents': 224750}})
fij = {f'months[{j}].{k}': c(M['periodos'][j][a][b]) for j in range(6) for k, a, b in
       (('accrued.total_cents', 'devengado', 'total'), ('cash.total_cents', 'efectivo', 'total'), ('cumulative_diff_cents', 'diferencia_acumulada', None))
       if b} | {f'months[{j}].cumulative_diff_cents': c(M['periodos'][j]['diferencia_acumulada']) for j in range(6)}
fij |= {f'months[{j}].provisional': None for j in range(12)} | {'repeat_cadence_enabled': False}
fij |= {'months[6].cash.payment_count': 2, 'months[7].cash.payment_count': 2, 'months[8].cash.payment_count': 1,
        'months[9].cash.payment_count': 0, 'months[10].cash.payment_count': 0, 'months[11].cash.payment_count': 0}
fx('s2-proyeccion-7-12-defecto', 'Proyección de 12 meses con el interruptor apagado y sin partidas: meses 1–6 = P1–P6; 7–12 solo pagos programados; resto calculado por el oráculo',
   [F('criterio', 'S2-09'), F('decision', 'K21'), F('regla', 'RN-28')] + [F('model.json', f'periodos[{j}]') for j in range(6)] + [F('model.json', f'despensas[{k}].pagos') for k in (2, 3)],
   [('proyeccion', 'proyeccion', {})], fijos={'proyeccion': fij})
# cadencia de prueba (opción A del bloqueo S2-proyeccion-interruptor): despensas 5 y 6 y papel de contado de Karen en meses alternos
cad_p, cad_g = [], []
for n, (fc, gid) in enumerate(((dt.date(2027, 6, 10), 7), (dt.date(2027, 8, 10), 8))):
    d = dict(D[1])
    d['fecha'] = fc.isoformat()
    cad_g.append(gasto_despensa(None, d, gid, n + 5, proyectado=True))
    d['_fecha'] = fc.isoformat()
    cad_p.append(compra(d, n + 5, gid, fechas=[x.isoformat() for x in modelo_v3.sched(fc)], proyectado=True))
for ident, f in ((9, '2027-05-10'), (10, '2027-07-10'), (11, '2027-09-10')):
    g = papel_contado(ident, 2)
    g['expense_date'] = f
    g['projected'] = True
    cad_g.append(g)
fx('s2-proyeccion-7-12-cadencia', 'Proyección con el interruptor encendido y una cadencia de prueba (valores de la opción A del bloqueo; FIRME desde la '
   'puerta 2, Q2-21 a Q2-24: provisional null en los 12 meses); todo calculado por el oráculo',
   [F('criterio', 'S2-09'), F('decision', 'K21'), F('model.json', 'despensas[1]'), F('model.json', 'papel[2]')],
   [('proyeccion', 'proyeccion', {'from': '2026-10'})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {'params': {'repeat_cadence_enabled': True}, 'cadence': {'purchases': cad_p, 'expenses': cad_g}}},
   fijos={'proyeccion': {f'months[{j}].provisional': None for j in range(12)} | {'repeat_cadence_enabled': True}})
HOGAR = {'parts': [], 'rest_person': 'axel'}
CONJ = {'parts': [{'person': 'karen', 'pct_bp': 5000}], 'rest_person': 'axel'}
SUELDO = {'id': 1, 'kind': 'income', 'description': 'Sueldo (dato de fixture)', 'amount_cents': 1200000, 'second_half_amount_cents': 1250000,
          'frequency': 'fortnightly', 'interval_count': 1, 'start_date': '2026-10-15', 'end_date': None, 'category': None, 'payment_method': None,
          'accrual_months': None, 'person': 'axel', 'share_rule': None, 'credit_id': None, 'is_active': True}
RENTA = {'id': 2, 'kind': 'fixed_expense', 'description': 'Renta (dato de fixture)', 'amount_cents': 800000, 'second_half_amount_cents': None,
         'frequency': 'monthly', 'interval_count': 1, 'start_date': '2026-10-20', 'end_date': None, 'category': 'vivienda', 'payment_method': 'transfer',
         'accrual_months': 1, 'person': None, 'share_rule': HOGAR, 'credit_id': None, 'is_active': True}
VENIDERO = {'id': 3, 'kind': 'upcoming_expense', 'description': 'Gasto venidero conjunto (dato de fixture)', 'amount_cents': 300001, 'second_half_amount_cents': None,
            'frequency': 'once', 'interval_count': 1, 'start_date': '2027-01-25', 'end_date': None, 'category': 'varios', 'payment_method': 'cash',
            'accrual_months': 2, 'person': None, 'share_rule': CONJ, 'credit_id': None, 'is_active': True}
VIVIENDA = {'code': 'vivienda', 'name': 'Vivienda', 'kind': 'fixed', 'is_active': True}
fx('s2-proyeccion-csv', 'Proyección apagada con partidas de fixture (sueldo quincenal, renta mensual «Hogar», venidero «Conjunto 50/50» a 2 meses) sin capturas ligadas; firme; resto calculado por el oráculo',
   [F('criterio', 'S2-09'), F('decision', 'K21'), F('regla', 'RN-03'), F('regla', 'RN-02'), F('manual', 'ingresos, partidas y devengo de ene-feb 2027 calculados a mano')],
   [('proyeccion', 'proyeccion', {'from': '2026-10'})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {'categories': CATS + [VIVIENDA, VARIOS], 'planned_items': [SUELDO, RENTA, VENIDERO]}},
   fijos={'proyeccion': {f'months[{j}].provisional': None for j in range(12)} |
          {'months[0].income_cents': 2450000, 'months[4].income_cents': 2450000, 'months[0].planned_expense_cents': 800000,
           'months[3].planned_expense_cents': 1100001, 'months[3].accrued.total_cents': 1221638, 'months[3].accrued.by_person.karen': 169887,
           'months[4].accrued.total_cents': 1201737}})
fx('s2-liga-provisional', 'Sueldo quincenal con una captura del 14-oct ligada a la ocurrencia del 15-oct (FIRME desde la puerta 2, Q2-14 = A: '
   'provisional null; el id de la fixture se conserva); resto calculado por el oráculo',
   [F('criterio', 'S2-09'), F('criterio', 'S2-10'), F('manual', 'captura del 2026-10-14 de fixture')],
   [('flujo', 'flujo_quincenal', {'from': '2026-10-q1', 'count': 2}), ('proyeccion', 'proyeccion', {'from': '2026-10'})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {'planned_items': [SUELDO], 'incomes': [
       {'id': 1, 'person': 'axel', 'kind': 'fixed', 'amount_cents': 1195000, 'income_date': '2026-10-14', 'planned_item_id': 1}]}},
   fijos={'flujo': {'[0].provisional': None, '[1].provisional': None, '[0].income.expected_cents': 0, '[0].income.fixed_cents': 1195000,
                    '[1].income.expected_cents': 1250000},
          'proyeccion': {f'months[{j}].provisional': None for j in range(12)} | {'months[0].income_cents': 2445000}})
fl = {}
for j in range(16):
    k, i = divmod(j, 4)
    fl[f'[{j}].credit_payments.household_cents'] = pg(k, i, 'axel')
    fl[f'[{j}].credit_payments.by_person.karen'] = pg(k, i, 'karen')
    fl[f'[{j}].credit_payments.external_informative_cents'] = pg(k, i, 'karen')
    fl[f'[{j}].pay_date'] = D[k]['pagos'][i]['fecha']
fx('s2-flujo-quincenal', '16 quincenas desde 2026-10-q2 (un pago de despensa por quincena) y valor estable para el contraste (S2-10); resto calculado por el oráculo',
   [F('criterio', 'S2-10')] + [F('model.json', f'despensas[{k}].pagos') for k in range(4)] + [F('model.json', 'valor_estable')],
   [('flujo', 'flujo_quincenal', {'from': '2026-10-q2', 'count': 16}), ('valor-estable', 'valor_estable', {})],
   fijos={'flujo': fl, 'valor-estable': {'month.by_person.axel': c(M['valor_estable']['mes']['axel']), 'month.by_person.karen': c(M['valor_estable']['mes']['karen']),
                                         'fortnight.by_person.axel': c(M['valor_estable']['quincena']['axel'])}})
PRES = {'id': 1, 'category': 'despensa', 'amount_cents': 200000, 'valid_from': '2026-10-01', 'valid_to': None}
fx('s2-presupuesto-scope', 'Cambio de presupuesto con alcance (K2), plan contra real (M-14) y el ejemplo de RN-22; resto calculado por el oráculo',
   [F('decision', 'K2'), F('regla', 'RN-22'), F('criterio', 'S2-11'), F('model.json', 'periodos[0].devengado.axel'), F('model.json', 'periodos[2].devengado.axel')],
   [('onward-dic', 'aplicar_cambio_presupuesto', {'row_id': 1, 'month': '2026-12', 'amount_cents': 250000, 'scope': 'onward'}),
    ('solo-dic', 'aplicar_cambio_presupuesto', {'row_id': 1, 'month': '2026-12', 'amount_cents': 250000, 'scope': 'this_month'}),
    ('sin-alcance-dic', 'aplicar_cambio_presupuesto', {'row_id': 1, 'month': '2026-12', 'amount_cents': 250000}),
    ('solo-oct-inicio', 'aplicar_cambio_presupuesto', {'row_id': 1, 'month': '2026-10', 'amount_cents': 180000, 'scope': 'this_month'}),
    ('fuera-de-vigencia', 'aplicar_cambio_presupuesto', {'row_id': 1, 'month': '2026-09', 'amount_cents': 180000}),
    ('dic-2026', 'presupuesto_mes', {'month': '2026-12'}), ('oct-2026', 'presupuesto_mes', {'month': '2026-10'}),
    ('plan-real', 'plan_contra_real', {'from': '2026-10', 'count': 3, 'category': 'despensa'})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {'budgets': [PRES]}},
   fijos={'onward-dic': {'rows': [{'id': 1, 'category': 'despensa', 'amount_cents': 200000, 'valid_from': '2026-10-01', 'valid_to': '2026-11-30'},
                                  {'id': 2, 'category': 'despensa', 'amount_cents': 250000, 'valid_from': '2026-12-01', 'valid_to': None}], 'effective_row_id': 2},
          'solo-dic': {'rows': [{'id': 1, 'category': 'despensa', 'amount_cents': 200000, 'valid_from': '2026-10-01', 'valid_to': '2026-11-30'},
                                {'id': 2, 'category': 'despensa', 'amount_cents': 250000, 'valid_from': '2026-12-01', 'valid_to': '2026-12-31'},
                                {'id': 3, 'category': 'despensa', 'amount_cents': 200000, 'valid_from': '2027-01-01', 'valid_to': None}], 'effective_row_id': 2},
          'sin-alcance-dic': {'effective_row_id': 2, 'deleted_row_ids': []},
          'fuera-de-vigencia': {'error': 'out_of_range'},
          'dic-2026': {'categories[0].plan_cents': 200000, 'categories[0].actual_cents': c(M['periodos'][2]['devengado']['axel']), 'categories[0].used_pct': 88,
                       'categories[0].warning': True, 'categories[0].limit': False, 'categories[0].state': 'warning'},
          'oct-2026': {'accrued_household_cents': c(M['periodos'][0]['devengado']['axel']), 'cash_household_cents': c(M['periodos'][0]['efectivo']['axel'])},
          'plan-real': {'[0].actual_cents': 119764, '[1].actual_cents': 119764, '[2].actual_cents': 176751, '[2].plan_cents': 200000}})
UMB = [('alerta_079', 79000, 79, False, False, 'normal'), ('alerta_079_999', 79999, 80, False, False, 'normal'), ('alerta_080', 80000, 80, True, False, 'warning'),
       ('alerta_099', 99000, 99, True, False, 'warning'), ('alerta_100', 100000, 100, True, True, 'exceeded'), ('alerta_101', 101000, 101, True, True, 'exceeded')]
fa = {'days_left': 21}
for i, (cat, _, u, w, lim, st) in enumerate(UMB):
    fa |= {f'categories[{i}].category': cat, f'categories[{i}].used_pct': u, f'categories[{i}].warning': w, f'categories[{i}].limit': lim, f'categories[{i}].state': st}
fx('s2-alertas-80-100', 'Banderas 80/100 con frontera exacta: 79 %, 79.999 %, 80 %, 99 %, 100 % y 101 % sobre planes de 100000 (S2-11); resto calculado por el oráculo',
   [F('criterio', 'S2-11'), F('regla', 'RN-22'), F('decision', 'K44'), F('manual', 'gastos de fixture del 2026-10-05')],
   [('oct-2026', 'presupuesto_mes', {'month': '2026-10'})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {
       'categories': [{'code': u[0], 'name': u[0], 'kind': 'variable', 'is_active': True} for u in UMB] + CATS,
       'expenses': GASTOS + [gasto_fx(7 + i, u[0], u[1], '2026-10-05') for i, u in enumerate(UMB)],
       'budgets': [{'id': i + 1, 'category': u[0], 'amount_cents': 100000, 'valid_from': '2026-10-01', 'valid_to': None} for i, u in enumerate(UMB)]}},
   fijos={'oct-2026': fa})
METAS = [{'id': 1, 'name': 'Meses exactos', 'kind': 'goal', 'target_cents': 600000, 'target_date': '2027-03-31', 'planned_monthly_cents': 100000, 'start_date': '2026-10-01', 'is_active': True},
         {'id': 2, 'name': 'Con residuo (RN-18)', 'kind': 'goal', 'target_cents': c(M['horizonte']['deuda_viva_10_abr']), 'target_date': '2027-04-10', 'planned_monthly_cents': 44953, 'start_date': '2026-10-01', 'is_active': True},
         {'id': 3, 'name': 'Cumplida', 'kind': 'goal', 'target_cents': 100000, 'target_date': '2026-12-31', 'planned_monthly_cents': None, 'start_date': '2026-09-01', 'is_active': True},
         {'id': 4, 'name': 'Fondo', 'kind': 'fund', 'target_cents': None, 'target_date': None, 'planned_monthly_cents': 20000, 'start_date': '2026-10-01', 'is_active': True},
         {'id': 5, 'name': 'Archivada', 'kind': 'goal', 'target_cents': 50000, 'target_date': '2026-12-31', 'planned_monthly_cents': 10000, 'start_date': '2026-09-01', 'is_active': False}]
MOVS = [{'goal_id': 2, 'movement_date': '2026-10-05', 'kind': 'deposit', 'amount_cents': 314672},
        {'goal_id': 3, 'movement_date': '2026-09-15', 'kind': 'deposit', 'amount_cents': 120000},
        {'goal_id': 3, 'movement_date': '2026-10-02', 'kind': 'withdrawal', 'amount_cents': 10000},
        {'goal_id': 4, 'movement_date': '2026-10-10', 'kind': 'deposit', 'amount_cents': 50000},
        {'goal_id': 4, 'movement_date': '2026-10-20', 'kind': 'deposit', 'amount_cents': 7000},
        {'goal_id': 5, 'movement_date': '2026-10-01', 'kind': 'deposit', 'amount_cents': 5000}]


fx('s2-metas', 'Metas: meses exactos, residuo de centavos (RN-18: 50.0 %), meta cumplida, fondo sin monto y una archivada; FIRME desde la puerta 2 '
   '(Q2-15 a Q2-18), state incluido; resto calculado por el oráculo',
   [F('criterio', 'S2-12'), F('regla', 'RN-18'), F('model.json', 'horizonte.deuda_viva_10_abr')],
   [('metas', 'metas', {})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {'savings_goals': METAS, 'savings_movements': MOVS}},
   fijos={'metas': {'goals[0].state': 'in_progress', 'goals[1].state': 'in_progress', 'goals[2].state': 'achieved', 'goals[3].state': 'in_progress','goals[1].saved_cents': 314672, 'goals[1].progress_tenths_pct': 500, 'goals[1].progress_pct': 50,
                    'goals[2].plan': [], 'goals[2].needed_monthly_cents': 0, 'goals[3].target_cents': None, 'goals[3].progress_pct': None,
                    'goals[3].needed_monthly_cents': None, 'goals[3].plan': [], 'goals[3].saved_cents': 50000}})

# Meta vencida sin cumplir (P2c, puerta 2): state y months_left FIRMES (Q2-16 = A, Q2-18 = A); needed_monthly_cents y plan
# de la vencida PROVISIONALES (bloqueo S2-metas-vencida: todo el restante en el mes de hoy). Hoy = 2026-10-11.
VENCIDAS = [{'id': 1, 'name': 'Vencida este mes', 'kind': 'goal', 'target_cents': 500, 'target_date': '2026-10-05', 'planned_monthly_cents': None, 'start_date': '2026-01-01', 'is_active': True},
            {'id': 2, 'name': 'Vencida el mes anterior', 'kind': 'goal', 'target_cents': 500, 'target_date': '2026-09-30', 'planned_monthly_cents': None, 'start_date': '2026-01-01', 'is_active': True},
            {'id': 3, 'name': 'Fondo cumplido sin fecha', 'kind': 'fund', 'target_cents': 500, 'target_date': None, 'planned_monthly_cents': None, 'start_date': '2026-01-01', 'is_active': True},
            {'id': 4, 'name': 'Fondo con fecha pasada', 'kind': 'fund', 'target_cents': None, 'target_date': '2026-09-30', 'planned_monthly_cents': None, 'start_date': '2026-01-01', 'is_active': True},
            {'id': 5, 'name': 'Fondo sin monto ni fecha', 'kind': 'fund', 'target_cents': None, 'target_date': None, 'planned_monthly_cents': None, 'start_date': '2026-01-01', 'is_active': True}]
MOVS_V = [{'goal_id': 1, 'movement_date': '2026-10-01', 'kind': 'deposit', 'amount_cents': 100},
          {'goal_id': 2, 'movement_date': '2026-10-01', 'kind': 'deposit', 'amount_cents': 100},
          {'goal_id': 3, 'movement_date': '2026-10-01', 'kind': 'deposit', 'amount_cents': 600}]
fv = {}
for i, (st, ml, nd, pl) in enumerate([('overdue', 1, 400, [{'month': '2026-10', 'cents': 400}]), ('overdue', 0, 400, [{'month': '2026-10', 'cents': 400}]),
                                       ('achieved', None, None, []), ('overdue', 0, None, []), ('in_progress', None, None, [])]):
    fv |= {f'goals[{i}].state': st, f'goals[{i}].months_left': ml, f'goals[{i}].needed_monthly_cents': nd, f'goals[{i}].plan': pl}
fx('s2-metas-vencida', 'Metas vencidas y fondos: target_date anterior a hoy en el mes de hoy (months_left 1) y en el mes anterior (months_left 0), '
   'fondo con monto cumplido sin fecha (achieved), fondo sin monto con fecha pasada (overdue) y fondo sin monto ni fecha (in_progress); '
   'la vencida sin cumplir pone todo el restante en el mes de hoy (PROVISIONAL, S2-metas-vencida); cifras a mano',
   [F('criterio', 'S2-12'), F('regla', 'RN-18'), F('manual', 'S2-interfaz §3.13: months_left, state y meta vencida; restante 500 − 100 = 400')],
   [('metas', 'metas', {})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {'savings_goals': VENCIDAS, 'savings_movements': MOVS_V}},
   provisional={'bloqueo': 'docs/estado/bloqueos/S2-metas-vencida.md',
                'que': 'needed_monthly_cents y plan de una meta vencida sin cumplir (todo el restante en el mes de hoy, regla en uso)'},
   fijos={'metas': fv})

# ---- Errores de entrada (correccion-S2-r1: hallazgos M1 y M2 de REVISION-S2-r1) ----
ctl = {}
for j in range(2):
    ctl |= {f'[{j}].accrued.total_cents': c(M['periodos'][j]['devengado']['total']), f'[{j}].cash.total_cents': c(M['periodos'][j]['efectivo']['total']),
            f'[{j}].cumulative_diff_cents': c(M['periodos'][j]['diferencia_acumulada'])}
# correccion-S2-r2 (hallazgo M1 de REVISION-S2-r2): ventanas fuera de orden, traslapadas o contenidas.
ERR_ORDEN = {'desordenadas-p2-p1': [VENT[1], VENT[0]],
             'desordenadas-p3-p1-p2': [VENT[2], VENT[0], VENT[1]],
             'traslape-p1-x': [VENT[0], {'label': 'X', 'start': '2026-10-25', 'end': '2026-12-09'}],
             'traslape-un-dia': [VENT[0], {'label': 'Y', 'start': '2026-11-09', 'end': '2026-12-09'}],
             'contenida-a-b': [{'label': 'A', 'start': '2026-10-10', 'end': '2027-01-09'},
                               {'label': 'B', 'start': '2026-11-10', 'end': '2026-12-09'}]}
fx('s2-ventanas-sin-cubrir', 'Ventanas que no cubren un costo o un pago: una sola ventana posterior al primer evento (P2), un hueco con solo pagos (P1 y P3) '
   'y un hueco con solo un costo (el 10-dic-2026) dan invalid_argument; también ventanas que no van ascendentes y disjuntas: [P2, P1], [P3, P1, P2], '
   'P1 + X (2026-10-25 a 2026-12-09), P1 + Y (empieza el 2026-11-09, el último día de P1) y A (2026-10-10 a 2027-01-09) con B (2026-11-10 a 2026-12-09) '
   'dentro (§2.12 reglas 1 a 3, FIRMES); puerta 2: un hueco sin eventos (2026-11-16 a 2026-11-29) que cruza el devengo de la despensa 1 es '
   'invalid_argument (regla 4, PROVISIONAL) y [d = 2026-10-10, r = 2026-10-11 a 2026-11-09] da 292125 al 2026-11-09 (§3.2 caso 2, PROVISIONAL); '
   'control: P1 y P2 contiguas = R-45 y R-46',
   [F('criterio', 'S2-03'), F('criterio', 'S2-16'), F('manual', 'S2-interfaz §2.12 reglas 1 a 4 y §3.2 (docs/estado/bloqueos/S2-ventanas.md, preguntas 1 a 4)'),
    F('manual', 'REVISION-S2-r3 M1: [d, r] = 389500 devengado − 97375 pagado = 292125'), F('model.json', 'periodos[0]'), F('model.json', 'periodos[1]')],
   [('control-p1-p2', 'periodos', {'windows': VENT[:2]}), ('solo-p2', 'periodos', {'windows': [VENT[1]]}),
    ('hueco-con-pagos', 'periodos', {'windows': [VENT[0], VENT[2]]}),
    ('hueco-con-costo', 'periodos', {'windows': [{'label': 'A', 'start': '2026-10-10', 'end': '2026-12-09'},
                                                 {'label': 'B', 'start': '2026-12-11', 'end': '2027-01-09'}]})]
   + [(k, 'periodos', {'windows': v}) for k, v in ERR_ORDEN.items()]
   + [('hueco-sin-eventos', 'periodos', {'windows': [{'label': 'A1', 'start': '2026-10-10', 'end': '2026-11-15'},
                                                    {'label': 'C1', 'start': '2026-11-30', 'end': '2027-01-09'}]}),
      ('d-r-caso-2', 'periodos', {'windows': [{'label': 'd', 'start': '2026-10-10', 'end': '2026-10-10'},
                                             {'label': 'r', 'start': '2026-10-11', 'end': '2026-11-09'}]})],
   provisional={'bloqueo': 'docs/estado/bloqueos/S2-ventanas.md',
                'que': 'pregunta 4: contigüidad de las ventanas (caso 1, regla 4 de §2.12) y devengo con ventanas que no duran un mes (caso 2, §3.2)'},
   fijos={'control-p1-p2': ctl, 'solo-p2': {'error': 'invalid_argument'}, 'hueco-con-pagos': {'error': 'invalid_argument'},
          'hueco-con-costo': {'error': 'invalid_argument'}, 'hueco-sin-eventos': {'error': 'invalid_argument'},
          'd-r-caso-2': {'[1].cumulative_diff_cents': 292125}} | {k: {'error': 'invalid_argument'} for k in ERR_ORDEN})
# Puerta 2 (REVISION-S2-r3 M1, consulta A + C): la semilla sin despensas a crédito más un gasto de contado de 100000 el 2026-10-10 con
# accrual_months 2. Va en su propia fixture porque el hueco de A + C debe quedar sin eventos y la semilla tiene pagos en él.
VA, VB, VC = ({'label': l, 'start': a, 'end': b} for l, a, b in (('A', '2026-10-10', '2026-11-09'), ('B', '2026-11-10', '2026-12-09'),
                                                                   ('C', '2026-12-10', '2027-01-09')))
fx('s2-ventanas-a-mas-c', 'Ventanas A y C con un hueco sin eventos (B) que cruza un devengo de 2 meses: invalid_argument (§2.12 regla 4, PROVISIONAL); '
   'control A, B y C contiguas: A 50000/100000/−50000, B devenga 50000 y C 0 con diferencia 0; en meses [oct, dic] diciembre devenga 0. '
   'Datos de REVISION-S2-r3 M1: la semilla sin despensas a crédito más un gasto de contado de 100000 el 2026-10-10 con accrual_months 2; cifras a mano',
   [F('criterio', 'S2-03'), F('criterio', 'S2-16'), F('manual', 'S2-interfaz §2.12 regla 4 y REVISION-S2-r3 M1: 100000 / 2 = 50000 por mes')],
   [('a-mas-c', 'periodos', {'windows': [VA, VC]}), ('a-b-c', 'periodos', {'windows': [VA, VB, VC]}),
    ('meses-oct-dic', 'periodos', {'months': ['2026-10', '2026-12']})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {
       'categories': CATS, 'credit_purchases': [],
       'expenses': [g for g in GASTOS if g['category'] != 'despensa'] + [gasto_fx(7, 'despensa', 100000, '2026-10-10') | {'accrual_months': 2}]}},
   provisional={'bloqueo': 'docs/estado/bloqueos/S2-ventanas.md', 'que': 'pregunta 4, caso 1: contigüidad de las ventanas (regla 4 de §2.12)'},
   fijos={'a-mas-c': {'error': 'invalid_argument'},
          'a-b-c': {'[0].accrued.total_cents': 50000, '[0].cash.total_cents': 100000, '[0].cumulative_diff_cents': -50000,
                    '[1].accrued.total_cents': 50000, '[2].accrued.total_cents': 0, '[2].cumulative_diff_cents': 0},
          'meses-oct-dic': {'[1].accrued.total_cents': 0}})
fx('s2-rango-devengo-0', 'Gasto de fixture del 5-dic-2026 con accrual_months 0 explícito: fuera del rango 1–24 (S2-interfaz §2.4); ausente o null valdría 1',
   [F('criterio', 'S2-03'), F('manual', 'S2-interfaz §2.4: accrual_months entero 1–24')],
   [('dic-2026', 'periodos', {'months': ['2026-12']})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {'categories': CATS + [VARIOS],
                                                           'expenses': GASTOS + [gasto_fx(7, 'varios', 30000, '2026-12-05') | {'accrual_months': 0}]}},
   fijos={'dic-2026': {'error': 'invalid_argument'}})
fx('s2-rango-intervalo-0', 'Partida mensual de fixture (renta) con interval_count 0 explícito: fuera del rango 1–24 (S2-interfaz §2.7); ausente o null valdría 1',
   [F('criterio', 'S2-03'), F('manual', 'S2-interfaz §2.7: interval_count entero 1–24')],
   [('proyeccion', 'proyeccion', {'from': '2026-10'})],
   datos={'$ref': '_datos/semilla-oraculo.json', 'parche': {'categories': CATS + [VIVIENDA], 'planned_items': [RENTA | {'interval_count': 0}]}},
   fijos={'proyeccion': {'error': 'invalid_argument'}})

casos = sorted(p.stem for p in CARPETA.glob('*.json') if not p.name.startswith('_'))
print(f'{len(casos)} fixtures y 3 archivos de _datos escritos en {CARPETA.relative_to(RAIZ)}')
for h in HALLAZGOS:
    print('HALLAZGO', h)
