#!/usr/bin/env python3
"""Generador determinista de las fixtures de S3 que no son del comparador (S3-interfaz.md §10.0 y §10).

Escribe oraculo/fixtures/s3/_datos/kueski-semilla.json y un archivo por caso de la tabla de §10 salvo los
s3-comparador-* (autor: optimizador S3). Reglas de autoría (§10, «De dónde sale cada esperado»):
- Canon: las hojas de oraculo/model.json convertidas con FORMATO.md §3 leyendo solo model.json (y data.json
  para R-05 y R-06). Nunca se llama a cred_lib ni al motor Go para un «esperado».
- Sin canon: cifras de la interfaz, de S3.md o de las decisiones (citadas en «fuente»), o cuentas a mano
  hechas aquí con decimal.Decimal (método distinto del oráculo, que usa fractions.Fraction) y explicadas en
  «descripcion».
Uso: python3 -I oraculo/cred_genera_fixtures.py [--out oraculo/fixtures/s3]. Solo biblioteca estándar.
"""
import argparse
import json
import os
import sys
from decimal import ROUND_HALF_UP, Decimal, getcontext

BASE = os.path.dirname(os.path.abspath(__file__))
RAIZ = os.path.dirname(BASE)
CONTEOS = {'n', 'n_pagos', 'despensa', 'despensa_que_cubre', 'paquete'}
B_TARJETA = 'docs/estado/bloqueos/S3-formula-tarjeta.md'
B_PRESTAMO = 'docs/estado/bloqueos/S3-convencion-prestamo.md'
B_CARGO = 'docs/estado/bloqueos/S3-origen-cargo-normal.md'
B_SERIE = 'docs/estado/bloqueos/S3-serie-sin-intereses.md'
B_IVA = 'docs/estado/bloqueos/U4-iva-credito.md'
getcontext().prec = 80


def leer(ruta):
    with open(ruta, encoding='utf-8') as fh:
        return json.load(fh, parse_float=Decimal)


M = leer(os.path.join(BASE, 'model.json'))
DJ = leer(os.path.join(RAIZ, 'datos', 'semilla', 'data.json'))


# ------------------------------------------------------------------ conversión (FORMATO.md §3)
def conv(v, clave):
    if isinstance(v, (bool, str)) or v is None:
        return v
    if clave in CONTEOS:
        return int(v)
    x = Decimal(v) * 100
    if x != x.to_integral_value():
        raise SystemExit(f'dinero con fracción de centavo: {clave}={v}')
    return int(x)


def arbol(x, clave=None, quitar=()):
    if isinstance(x, dict):
        return {k: arbol(v, k) for k, v in x.items() if k not in quitar}
    if isinstance(x, list):
        return [arbol(v, clave) for v in x]
    return conv(x, clave)


def C(v):
    return conv(v, 'dinero')


def hu(x):
    """HALF_UP a entero de un Decimal (cuentas a mano)."""
    return int(Decimal(x).quantize(Decimal(1), rounding=ROUND_HALF_UP))


# ------------------------------------------------------------------ piezas de entrada (§5, §5.7)
def credito(cid, name, kind, holder, limit, freq, fpmd=None, reading='B', confirmed=False, anchor=None,
            sd=None, dd=None, annual=None, base='', vat=None, pct=None, fixed=None, term_rates=(), promos=()):
    return {'id': cid, 'name': name, 'kind': kind, 'holder': holder, 'limit_cents': limit, 'frequency': freq,
            'first_payment_min_days': fpmd, 'reading': reading, 'reading_confirmed': confirmed,
            'reading_a_anchor': anchor, 'statement_day': sd, 'payment_due_day': dd, 'annual_rate': annual,
            'rate_base': base, 'vat_on_interest_bp': vat, 'min_payment_pct_bp': pct,
            'min_payment_fixed_cents': fixed, 'term_rates': list(term_rates), 'promotions': list(promos)}


def tasa(num, den, dec):
    return {'num': num, 'den': den, 'decimales': dec}


def tr(n, num, den, dec):
    return {'term_count': n, 'rate': tasa(num, den, dec), 'source': 'manual', 'quote_date': None}


def promo(pid, cid, merchant, n, minimo, vf, vt, rate=(0, 1, 0)):
    return {'id': pid, 'credit_id': cid, 'merchant': merchant, 'store_id': None, 'term_count': n,
            'rate': tasa(*rate), 'min_purchase_cents': minimo, 'valid_from': vf, 'valid_to': vt, 'active': True}


def parte(person, ext, P, joint=None, paper=None):
    return {'person': person, 'external': ext, 'principal_cents': P, 'joint_cents': joint, 'paper_in_cents': paper}


def compra(date, P, n, shares, quoted=None, rd=None, pid=None, merchant='', fpd=None):
    return {'date': date, 'principal_cents': P, 'term_count': n, 'quoted_interest_cents': quoted,
            'rate_decimals': rd, 'promotion_id': pid, 'store_id': None, 'merchant': merchant,
            'first_payment_date': fpd, 'shares': shares}


def estado(credit, purchases=(), payments=(), statements=()):
    return {'credit': credit, 'purchases': list(purchases), 'payments': list(payments),
            'statements': list(statements)}


def pst(ref, inp, status='active'):
    return {'ref': ref, 'status': status, 'input': inp}


def corte(sdate, due, bal, daily=(), source='bank'):
    return {'statement_date': sdate, 'due_date': due, 'balance_cents': bal, 'daily_balances_cents': list(daily),
            'source': source}


def consulta(qid, fn, args):
    return {'id': qid, 'funcion': fn, 'argumentos': args}


def src(tipo, ruta):
    return {'tipo': tipo, 'ruta': ruta}


def prov(bloqueo, que):
    return {'bloqueo': bloqueo, 'que': que}


SALIDA = {}


def fixture(caso, desc, fuente, provisional, datos, consultas, esperado):
    SALIDA[caso + '.json'] = {'formato': 'fixture/1', 'motor': 's3', 'caso': caso, 'descripcion': desc,
                              'fuente': fuente, 'provisional': provisional,
                              'entradas': {'datos': datos, 'consultas': consultas}, 'esperado': esperado}


# ------------------------------------------------------------------ semilla (§10, «Datos compartidos»)
D = M['despensas']
PAR = M['parametros']
DQ = DJ['financiamiento']['dos_quincenas']
HOY = D[0]['fecha']                                   # 2026-10-10 (hoja de texto de model.json)
CAP1, INT1 = C(D[0]['capital']), C(D[0]['intereses'])  # 309450, 80050
INT2Q = C(DQ['intereses'])                            # 53450 (data.json)
N4 = len(D[0]['pagos'])                               # 4
TASA4 = tasa(INT1, CAP1, 4)                           # K12: «25.8685 %» se capturó con 4 decimales
SEMILLA_REF = {'$ref': '_datos/kueski-semilla.json'}


def partes_de(k):
    a, kk = D[k]['axel'], D[k]['karen']
    return [parte('axel', False, C(a['capital']), C(a['conjunto']), C(a['papel_dentro'])),
            parte('karen', True, C(kk['capital']), C(kk['conjunto']), C(kk['papel_dentro']))]


def compra_d(k, **kw):
    base = dict(quoted=INT1, rd=4) if k == 0 else {}
    base.update(kw)
    return compra(D[k]['fecha'], C(D[k]['capital']), len(D[k]['pagos']), partes_de(k), **base)


def kueski(promos=()):
    return credito(1, 'Kueski', 'revolving_line', 'axel', C(PAR['limite_credito']), 'fortnightly', 15,
                   anchor=HOY, term_rates=[tr(2, INT2Q, CAP1, 2), tr(4, INT1, CAP1, 4)], promos=promos)


PAPEL = C(PAR['papel'])
SEMILLA = {'today': HOY,
           'credits': [estado(kueski(), [pst('d%d' % (k + 1), compra_d(k)) for k in range(len(D))])],
           'cash': [{'date': '2027-01-10', 'amount_cents': PAPEL, 'person': 'karen', 'description': 'Papel de contado'},
                    {'date': '2027-03-10', 'amount_cents': PAPEL, 'person': 'karen', 'description': 'Papel de contado'}]}


# ------------------------------------------------------------------ casos canon (§10, tabla)
def casos_canon():
    for k in range(len(D)):
        n = k + 1
        desp = arbol(D[k], quitar=('fecha',))
        disp = arbol(M['disposicion_por_despensa'][k], quitar=('fecha',))
        cons = [consulta('despensas[%d]' % k, 'vista_despensa', {'credit_id': 1, 'n': n}),
                consulta('disposicion_por_despensa[%d]' % k, 'vista_disposicion', {'credit_id': 1, 'n': n})]
        esp = {'despensas[%d]' % k: desp, 'disposicion_por_despensa[%d]' % k: disp}
        fuente = [src('model.json', 'despensas[%d]' % k), src('model.json', 'disposicion_por_despensa[%d]' % k)]
        crit = ['S3-01', 'S3-04', 'S3-19']
        desc = 'Despensa %d del canon con la semilla: capital, intereses, pagos, partes y disposición' % n
        if k == 0:
            r = {'num': INT1, 'den': CAP1}
            dosq = compra_d(0, quoted=None, rd=None)
            dosq['term_count'] = 2
            dec = -Decimal(str(DQ['tasa_pct'])).as_tuple().exponent
            cons = [consulta('parametros', 'vista_parametros', {'credit_id': 1})] + cons + [
                consulta('data.json:financiamiento.dos_quincenas', 'vista_dos_quincenas',
                         {'credit_id': 1, 'purchase': dosq}),
                consulta('dos-quincenas', 'cotizar_compra', {'credit_id': 1, 'purchase': dosq})]
            esp = {'parametros': {'tasa': r, 'tasa_pct': dict(r), 'karen_fija': arbol(PAR['karen_fija'])}, **esp,
                   'data.json:financiamiento.dos_quincenas': {
                       'total': C(DQ['total']), 'intereses': INT2Q,
                       'tasa_pct': tasa(INT2Q, C(DQ['total']) - INT2Q, dec),
                       'axel': arbol(DQ['axel']), 'karen': arbol(DQ['karen'])},
                   'dos-quincenas': {'interest_cents': INT2Q, 'total_cents': C(DQ['total']),
                                     'source': 'term_table', 'rate_display': '17.27 %'}}
            fuente = [src('model.json', 'parametros')] + fuente + [src('data.json', 'financiamiento.dos_quincenas')]
            crit = ['S3-01', 'S3-04', 'S3-06', 'S3-19']
            desc += ('; tasa por compra 800.50/3,094.50 (K12, 4 decimales); R-05 y R-06: la misma compra a 2 '
                     'quincenas sin cotización toma la tasa de la tabla 53450/309450 (2 decimales) y da 534.50 '
                     'de intereses y 3,629.00 de total')
        fixture('s3-kueski-despensa-%d' % n, desc, fuente + [src('criterio', c) for c in crit], None,
                SEMILLA_REF, cons, esp)
    # series A y B
    cred = arbol(M['credito'], quitar=('nota',))
    for e in cred['eventos']:
        e.pop('descripcion', None)
    fixture('s3-series-a-b',
            'Línea de tiempo del canon del 10-oct-2026 al 10-abr-2027: eventos, series A y B de 14 valores, '
            'mínimos 8,094.50 y 5,000.00; con la lectura omitida rige B y sale «pendiente confirmar»',
            [src('model.json', 'credito'), src('model.json', 'parametros'), src('criterio', 'S3-03'),
             src('criterio', 'S3-05')], None, SEMILLA_REF,
            [consulta('credito', 'vista_credito', {'credit_id': 1, 'today': HOY, 'from': HOY, 'to': '2027-04-10'}),
             consulta('parametros', 'vista_parametros', {'credit_id': 1}),
             consulta('lectura-omitida', 'linea_tiempo',
                      {'credit_id': 1, 'today': HOY, 'from': HOY, 'to': '2027-04-10', 'reading': ''})],
            {'credito': cred, 'parametros': {'papel': PAPEL, 'limite_credito': C(PAR['limite_credito'])},
             'lectura-omitida': {'reading': 'B', 'pending_confirmation': True,
                                 'series_b_cents': cred['comparacion_series']['esperado_B']}})
    # sensibilidad e interés por monto (receta de §10, corrige r4-P1)
    rate = tasa(INT1, CAP1, 4)
    sens_args = {'amounts_cents': [C(r['N']) for r in M['sensibilidad']], 'rate': rate, 'term_count': N4,
                 'external': parte('karen', True, C(PAR['karen_fija']['capital'])),
                 'holder_current_principal_cents': C(D[0]['axel']['capital']), 'accrual_months': 2,
                 'cash_monthly_cents': hu(Decimal(PAPEL) / 2)}
    fixture('s3-sensibilidad', 'Tabla de sensibilidad del canon (5 montos de 3,000.00 a 5,000.00) con el insumo '
            'de sensitivityFromModel: tasa de 4 q, n 4, Karen 1,191.50, Axel actual 1,903.00, 2 meses y 99.50',
            [src('model.json', 'sensibilidad'), src('criterio', 'S3-01'), src('criterio', 'S3-02')], None,
            {'today': HOY, 'credits': [], 'cash': []},
            [consulta('sensibilidad', 'vista_sensibilidad', sens_args)],
            {'sensibilidad': arbol(M['sensibilidad'])})
    fixture('s3-interes-por-monto', 'Interés por monto del canon (7 montos de 1.00 a 5,000.00) con la tasa de 4 q y '
            'n 4; pago de referencia desde el total sin redondear',
            [src('model.json', 'interes_por_monto'), src('criterio', 'S3-01'), src('criterio', 'S3-02')], None,
            {'today': HOY, 'credits': [], 'cash': []},
            [consulta('interes_por_monto', 'vista_interes_por_monto',
                      {'amounts_cents': [C(r['monto']) for r in M['interes_por_monto']], 'rate': rate,
                       'term_count': N4})],
            {'interes_por_monto': arbol(M['interes_por_monto'])})


# ------------------------------------------------------------------ casos sin canon o con cruces (§10, tabla)
KIND = {'inicio': 'start', 'compra': 'purchase', 'papel': 'cash', 'pago': 'payment'}


def desde_otra_fecha():
    ev = M['credito']['eventos']
    hasta = '2027-04-10'
    cons, esp = [], {}
    for qid, f in (('desde-2026-12-16', '2026-12-16'), ('desde-2027-02-10', '2027-02-10')):
        prev = [e for e in ev if e['fecha'] <= f][-1]
        ini = {'date': f, 'kind': 'start', 'amount_cents': 0, 'principal_restored_cents': 0,
               'interest_not_restored_cents': 0, 'available_a_cents': C(prev['disponible_A']),
               'available_b_cents': C(prev['disponible_B'])}
        sig = [{'date': e['fecha'], 'kind': KIND[e['tipo']], 'amount_cents': C(e['monto']),
                'principal_restored_cents': C(e['capital_repuesto']),
                'interest_not_restored_cents': C(e['interes_no_repuesto']),
                'available_a_cents': C(e['disponible_A']), 'available_b_cents': C(e['disponible_B'])}
               for e in ev if f < e['fecha'] <= hasta]
        todos = [ini] + sig
        sa = [e['available_a_cents'] for e in sig if e['kind'] in ('purchase', 'payment')]
        sb = [e['available_b_cents'] for e in sig if e['kind'] in ('purchase', 'payment')]
        cons.append(consulta(qid, 'linea_tiempo', {'credit_id': 1, 'today': HOY, 'from': f, 'to': hasta}))
        esp[qid] = {'events': todos, 'series_a_cents': sa, 'series_b_cents': sb,
                    'min_a_cents': min(e['available_a_cents'] for e in todos),
                    'min_b_cents': min(e['available_b_cents'] for e in todos), 'reading': 'B'}
    fixture('s3-series-desde-otra-fecha',
            'Línea de tiempo con From posterior al ancla (§5.2): el inicio vale el A y el B del canon tras el último '
            'evento con fecha <= From (9,094.50 / 6,000.00 el 16-dic-2026; 8,094.50 / 5,000.00 el 10-feb-2027) y '
            'los eventos siguientes son los de credito.eventos; valores copiados de model.json con id libres',
            [src('model.json', 'credito.eventos'), src('criterio', 'S3-03'), src('manual', 'descripcion')], None,
            SEMILLA_REF, cons, esp)


def reposicion_separada():
    ev = lambda d, k, ref, seq, amt, kr, ki, a, b: {  # noqa: E731
        'date': d, 'kind': k, 'purchase_ref': ref, 'seq': seq, 'amount_cents': amt, 'principal_restored_cents': kr,
        'interest_not_restored_cents': ki, 'available_a_cents': a, 'available_b_cents': b}
    datos = {'today': HOY, 'credits': [estado(kueski(), [pst('d1', compra_d(0))], [
        {'purchase_ref': 'd1', 'seq': 1, 'date': '2026-10-29', 'payer': 'karen', 'amount_cents': 37493},
        {'purchase_ref': 'd1', 'seq': 1, 'date': '2026-10-31', 'payer': 'axel', 'amount_cents': 59882}])],
        'cash': []}
    fixture('s3-reposicion-separada',
            'Despensa 1 con la parte de Karen pagada aparte el 29-oct y la de Axel el 31-oct: el capital del pago 1 '
            '(773.63) se repone el 31-oct y no antes (K15); el pago 2, sin abonos, se proyecta completo el 15-nov (§5.2)',
            [src('decision', 'K15'), src('criterio', 'S3-05'), src('manual', 'descripcion')], None, datos,
            [consulta('linea', 'linea_tiempo', {'credit_id': 1, 'today': HOY, 'from': HOY, 'to': '2026-11-15'})],
            {'linea': {'events': [ev(HOY, 'start', '', 0, 0, 0, 0, 1000000, 690550),
                                  ev('2026-10-29', 'payment', 'd1', 1, 37493, 0, 0, 1000000, 690550),
                                  ev('2026-10-31', 'payment', 'd1', 1, 59882, 77363, 20012, 1077363, 767913),
                                  ev('2026-11-15', 'payment', 'd1', 2, 97375, 77363, 20012, 1154726, 845276)],
                       'series_a_cents': [1000000, 1077363, 1154726], 'series_b_cents': [690550, 767913, 845276],
                       'min_a_cents': 1000000, 'min_b_cents': 690550, 'reading': 'B', 'pending_confirmation': True,
                       'a_anchor': HOY, 'a_anchor_defaulted': False}})


def deuda_viva():
    H = M['horizonte']
    det = H['deuda_viva_detalle']
    cons = [consulta('deuda-viva', 'deuda_viva', {'as_of': '2027-04-10'})]
    esp = {'deuda-viva': {'as_of': '2027-04-10', 'total_cents': C(H['deuda_viva_10_abr']),
                          'by_purchase': [{'credit_id': 1, 'purchase_ref': 'd3',
                                           'pending_cents': C(det['despensa_3_pendiente']), 'payment_count': 1},
                                          {'credit_id': 1, 'purchase_ref': 'd4', 'pending_cents': C(det['despensa_4']),
                                           'payment_count': 4}],
                          'by_credit': [{'credit_id': 1, 'purchase_ref': '', 'pending_cents': C(H['deuda_viva_10_abr']),
                                         'payment_count': 5}]}}
    for k, p in enumerate(M['periodos'][:6]):
        qid = 'efectivo-p%d' % (k + 1)
        cons.append(consulta(qid, 'efectivo_ventana', {'from': p['inicio'], 'to': p['fin']}))
        esp[qid] = {'total_cents': C(p['efectivo']['total'])}
    cons.append(consulta('efectivo-horizonte', 'efectivo_ventana', {'from': HOY, 'to': '2027-04-10'}))
    esp['efectivo-horizonte'] = {'cash_cents': 2 * PAPEL, 'total_cents': C(H['total_efectivo'])}
    cons.append(consulta('vivas', 'compras_vivas', {'credit_id': 1, 'at': '2027-04-10'}))
    esp['vivas'] = {'refs': ['d3', 'd4']}
    fixture('s3-deuda-viva-10-abr-2027',
            'Cruce con S2: deuda viva al 10-abr-2027 (6,293.44 = 1,258.70 de d3, 1 pago, + 5,034.74 de d4, 4 pagos), '
            'efectivo de las ventanas del 10 al 9 (P1 a P6), efectivo del horizonte con 398.00 de contado y compras '
            'vivas d3 y d4; valores de las ramas de S2 de model.json en centavos',
            [src('criterio', 'S3-01'), src('regla', 'RN-21'), src('model.json', 'horizonte'),
             src('model.json', 'periodos')], None, SEMILLA_REF, cons, esp)


def feb_2027():
    sh = partes_de(2)
    f27 = [p['fecha'] for p in D[2]['pagos']]
    f28 = ['2028-02-29', '2028-03-15', '2028-03-30', '2028-04-15']
    fixture('s3-feb-2027',
            'Calendario quincenal de RN-06 con febrero: compra del 10-feb-2027 paga el 28-feb (el «30» es el último '
            'día) y la del 10-feb-2028 (bisiesto) el 29-feb; a mano: primera quincena a >= 15 días y las siguientes',
            [src('criterio', 'S3-04'), src('model.json', 'despensas[2].pagos'), src('manual', 'descripcion')], None,
            SEMILLA_REF,
            [consulta('feb-2027', 'cotizar_compra', {'credit_id': 1, 'purchase': compra('2027-02-10', 400000, 4, sh)}),
             consulta('feb-2028', 'cotizar_compra', {'credit_id': 1, 'purchase': compra('2028-02-10', 400000, 4, sh)})],
            {'feb-2027': {'payments': [{'due_date': f} for f in f27]},
             'feb-2028': {'payments': [{'due_date': f} for f in f28]}})


def tasa_plazo():
    sh = partes_de(1)
    cons = [consulta('plazo-%d' % n, 'cotizar_compra',
                     {'credit_id': 1, 'purchase': compra('2026-12-10', 400000, n, sh)}) for n in (6, 8, 10, 12)]
    esp = {'plazo-%d' % n: {'error': 'no_rate'} for n in (6, 8, 10, 12)}
    cons += [consulta('plazo-4', 'cotizar_compra', {'credit_id': 1, 'purchase': compra('2026-12-10', 400000, 4, sh)}),
             consulta('d1-cotizada', 'cotizar_compra', {'credit_id': 1, 'purchase': compra_d(0)}),
             consulta('d1-sin-decimales', 'cotizar_compra', {'credit_id': 1, 'purchase': compra_d(0, rd=None)}),
             consulta('d1-decimales-negativos', 'cotizar_compra', {'credit_id': 1, 'purchase': compra_d(0, rd=-1)})]
    esp.update({'plazo-4': {'interest_cents': C(D[1]['intereses']), 'total_cents': C(D[1]['total']),
                            'source': 'term_table', 'rate_display': '25.8685 %'},
                'd1-cotizada': {'interest_cents': INT1, 'rate': tasa(INT1, CAP1, 4), 'rate_display': '25.8685 %',
                                'source': 'quoted'},
                'd1-sin-decimales': {'interest_cents': INT1, 'rate': tasa(INT1, CAP1, -1), 'rate_display': ''},
                'd1-decimales-negativos': {'error': 'invalid'}})
    fixture('s3-tasa-plazo',
            'Tabla por plazo: 6, 8, 10 y 12 quincenas sin tasa capturada dan «sin tasa» (nunca un número); 4,000.00 a '
            '4 q da 1,034.74 (R-17); d1 cotizada con 4 decimales se muestra «25.8685 %» y sin rate_decimals la tasa '
            'lleva decimales -1 y texto vacío (§5.1, K12); rate_decimals negativo es entrada inválida',
            [src('criterio', 'S3-06'), src('decision', 'K12'), src('model.json', 'despensas[1]')], None,
            SEMILLA_REF, cons, esp)


def promocion_0():
    pr = promo(1, 1, 'Chedraui', 4, 200000, '2026-10-28', '2026-10-30')
    datos = {'today': HOY, 'credits': [estado(kueski([pr]))], 'cash': []}
    sh = partes_de(0)
    sh_m = [parte('axel', False, 199999 - 119150), parte('karen', True, 119150)]

    def cp(qid, d, P, shares, pid):
        return consulta(qid, 'cotizar_compra', {'credit_id': 1, 'purchase': compra(d, P, 4, shares, pid=pid,
                                                                                  merchant='Chedraui')})

    def ap(qid, d, P):
        return consulta(qid, 'promociones_aplicables', {'credit_id': 1, 'date': d, 'principal_cents': P,
                                                        'term_count': 4, 'store_id': None, 'merchant': 'Chedraui'})
    cero = {'interest_cents': 0, 'total_cents': CAP1, 'source': 'promotion', 'promotion_id': 1,
            'payments': [{'amount_cents': a} for a in (77362, 77362, 77362, 77364)],
            'provisional': ['S3-serie-sin-intereses']}
    fixture('s3-promocion-0',
            'Promoción 0 % de fixture (Chedraui, 4 q, mínimo 2,000.00, del 28 al 30-oct-2026): 0.00 de intereses el '
            '28 y el 30 con montos 773.62 x3 + 773.64 (RN-07, truncado); el 27 y con 1,999.99 no aplica: pedirla da '
            'promotion_not_applicable y sin ella se cobra la tasa de 4 q (800.50; a mano HALF_UP(199,999 x 80050 / '
            '309450) = 517.37); el capital por pago no se pide (PROVISIONAL S3-serie-sin-intereses)',
            [src('criterio', 'S3-07'), src('regla', 'RN-10'), src('manual', 'descripcion')], None, datos,
            [cp('promo-28', '2026-10-28', CAP1, sh, 1), cp('promo-30', '2026-10-30', CAP1, sh, 1),
             cp('promo-27', '2026-10-27', CAP1, sh, 1), cp('normal-27', '2026-10-27', CAP1, sh, None),
             cp('promo-minimo-menos-1', '2026-10-28', 199999, sh_m, 1),
             cp('normal-minimo-menos-1', '2026-10-28', 199999, sh_m, None),
             ap('aplica-28', '2026-10-28', CAP1), ap('aplica-30', '2026-10-30', CAP1),
             ap('aplica-27', '2026-10-27', CAP1), ap('aplica-minimo-menos-1', '2026-10-28', 199999),
             ap('aplica-minimo-exacto', '2026-10-28', 200000)],
            {'promo-28': cero, 'promo-30': dict(cero), 'promo-27': {'error': 'promotion_not_applicable'},
             'normal-27': {'interest_cents': INT1, 'source': 'term_table', 'promotion_id': None},
             'promo-minimo-menos-1': {'error': 'promotion_not_applicable'},
             'normal-minimo-menos-1': {'interest_cents': 51737, 'source': 'term_table'},
             'aplica-28': {'promotion_ids': [1]}, 'aplica-30': {'promotion_ids': [1]},
             'aplica-27': {'promotion_ids': []}, 'aplica-minimo-menos-1': {'promotion_ids': []},
             'aplica-minimo-exacto': {'promotion_ids': [1]}})


# ------------------------------------------------------------------ tarjeta, préstamo, tienda y calendario
def tarjeta_k18(cid=2, promos=()):
    return credito(cid, 'Tarjeta con datos', 'bank_card', 'axel', 2000000, 'monthly', sd=5, dd=25,
                   annual=tasa(455, 1000, 1), base='daily', vat=1600, pct=500, fixed=20000, promos=promos)


def tarjeta_sin(cid=3):
    return credito(cid, 'Tarjeta sin datos', 'bank_card', 'axel', 2000000, 'monthly', sd=5, dd=25)


MSI = promo(21, 2, 'Liverpool', 3, 90000, '2026-10-01', '2026-12-31')
AX = lambda P: [parte('axel', False, P)]  # noqa: E731
P_CARGO_L = ['S3-formula-tarjeta', 'S3-origen-cargo-normal']
P_MSI_L = ['S3-formula-tarjeta', 'S3-serie-sin-intereses']
P_FORM_L = ['S3-formula-tarjeta', 'U4-iva-credito']


def tarjetas():
    datos = {'today': HOY, 'credits': [estado(tarjeta_k18(promos=[MSI])), estado(tarjeta_sin())], 'cash': []}
    st = corte('2026-11-05', '2026-11-25', 150000)
    cons, esp = [], {}
    for cid, nom in ((2, 'con-datos'), (3, 'sin-datos')):
        cons.append(consulta('cargo-' + nom, 'cotizar_compra', {'credit_id': cid, 'purchase': compra('2026-10-12', 150000, 1, AX(150000))}))
        esp['cargo-' + nom] = {'interest_cents': 0, 'total_cents': 150000, 'payments': [{'amount_cents': 150000}],
                               'provisional': P_CARGO_L}
        for qid, d in (('pago-total-dia-antes-', '2026-11-24'), ('pago-total-mismo-dia-', '2026-11-25')):
            cons.append(consulta(qid + nom, 'interes_tarjeta', {'credit_id': cid, 'statement': st,
                                                                 'paid_cents': 150000, 'paid_date': d}))
            esp[qid + nom] = {'interest_cents': 0, 'provisional': []}
    cons.append(consulta('msi-3', 'cotizar_compra', {'credit_id': 2, 'purchase': compra('2026-10-20', 100000, 3, AX(100000),
                                                                                     pid=21, merchant='Liverpool')}))
    esp['msi-3'] = {'interest_cents': 0, 'total_cents': 100000,
                    'payments': [{'amount_cents': a} for a in (33333, 33333, 33334)], 'provisional': P_MSI_L}
    fixture('s3-tarjeta-sin-intereses',
            'Tarjetas de fixture (límite 20,000.00, corte 5, pago 25), una con datos K18 y otra sin ellos: el cargo '
            'normal pagado completo un día antes o el mismo día de la fecha de pago da 0.00 con o sin datos (S3-08); '
            'MSI a 3 de 1,000.00: 333.33 x2 + 333.34 (truncado, el último ajusta); solo hojas FIRMES',
            [src('criterio', 'S3-08'), src('decision', 'K18'), src('manual', 'descripcion')], None, datos, cons, esp)
    # mínimo e interés con la fórmula PROVISIONAL (§5.3, opción b)
    diario = corte('2026-11-05', '2026-11-25', 150000, [150000] * 30)
    datos2 = {'today': HOY, 'credits': [estado(tarjeta_k18()), estado(tarjeta_sin())], 'cash': []}
    sigma = Decimal(150000 * 30)
    base = hu(sigma * Decimal('0.455') / 360)          # 4,500,000 x 0.455 / 360 = 5,687.5 -> 5,688
    iva = hu(Decimal(base) * Decimal('0.16'))          # 910.08 -> 910
    fixture('s3-tarjeta-minimo',
            'Tarjeta con datos de fixture (45.5 % anual, base diaria, IVA 16 %, mínimo 5 % o 200.00 fijo): saldo '
            'diario 1,500.00 por 30 días pagado un día después: base HALF_UP(45,000.00 x 0.455 / 360) = 56.88 e IVA '
            'HALF_UP(56.88 x 0.16) = 9.10, interés 65.98; mínimo máx(75.00, 200.00) = 200.00; con saldo 10,000.00 el '
            '5 % = 500.00; con saldo 150.00 el mínimo no pasa del saldo; sin datos, «falta capturar»',
            [src('manual', 'descripcion'), src('decision', 'K18'), src('criterio', 'S3-08')],
            prov(B_TARJETA, 'fórmula del interés del periodo y del mínimo (opción b); otros bloqueos: '
                 + B_IVA + ' (datos de «falta capturar» y lectura del IVA)'), datos2,
            [consulta('interes-dia-despues', 'interes_tarjeta', {'credit_id': 2, 'statement': diario,
                                                                 'paid_cents': 150000, 'paid_date': '2026-11-26'}),
             consulta('minimo', 'minimo_tarjeta', {'credit_id': 2, 'statement': diario}),
             consulta('minimo-saldo-alto', 'minimo_tarjeta', {'credit_id': 2, 'statement': corte('2026-11-05', '2026-11-25', 1000000)}),
             consulta('minimo-saldo-bajo', 'minimo_tarjeta', {'credit_id': 2, 'statement': corte('2026-11-05', '2026-11-25', 15000)}),
             consulta('interes-sin-datos', 'interes_tarjeta', {'credit_id': 3, 'statement': diario,
                                                               'paid_cents': 150000, 'paid_date': '2026-11-26'}),
             consulta('minimo-sin-datos', 'minimo_tarjeta', {'credit_id': 3, 'statement': diario})],
            {'interes-dia-despues': {'interest_cents': base + iva, 'provisional': P_FORM_L},
             'minimo': {'minimum_cents': 20000, 'provisional': P_FORM_L},
             'minimo-saldo-alto': {'minimum_cents': 50000, 'provisional': P_FORM_L},
             'minimo-saldo-bajo': {'minimum_cents': 15000, 'provisional': P_FORM_L},
             'interes-sin-datos': {'error': 'missing_capture'}, 'minimo-sin-datos': {'error': 'missing_capture'}})


def amort_mano(P, r, n):
    """Tabla francesa a mano con Decimal (80 dígitos): pago fijo, interés del saldo y último que ajusta."""
    exacto = Decimal(P) * r / (1 - 1 / (1 + r) ** n)
    if abs(exacto - exacto.to_integral_value(rounding='ROUND_FLOOR') - Decimal('0.5')) < Decimal('1e-40'):
        raise SystemExit('pago fijo demasiado cerca de medio centavo')
    pago, bal, filas = hu(exacto), P, []
    for k in range(n):
        i = hu(Decimal(bal) * r)
        a = pago - i if k < n - 1 else bal
        filas.append({'amount_cents': a + i, 'principal_cents': a, 'interest_cents': i})
        bal -= a
    return pago, filas


def pesos(c):
    return '%d.%02d' % (c // 100, c % 100)


def mensual(y, m, d, n):
    import calendar
    out = []
    for k in range(n):
        yy, mm = y + (m - 1 + k) // 12, (m - 1 + k) % 12 + 1
        out.append('%04d-%02d-%02d' % (yy, mm, min(d, calendar.monthrange(yy, mm)[1])))
    return out


def esperado_prestamo(P, r, n, fechas, provl):
    pago, filas = amort_mano(P, r, n)
    for f, row in zip(fechas, filas):
        row['due_date'] = f
    tot = sum(x['amount_cents'] for x in filas)
    return {'principal_cents': P, 'interest_cents': tot - P, 'total_cents': tot, 'reference_payment_cents': pago,
            'source': 'amortization', 'payments': [{'due_date': x['due_date'], 'amount_cents': x['amount_cents'],
                                                    'principal_cents': x['principal_cents'],
                                                    'interest_cents': x['interest_cents']} for x in filas],
            'shares': [], 'provisional': provl}


def prestamo():
    r_iva, r_sin = Decimal('0.24') / 12 * Decimal('1.16'), Decimal('0.24') / 12
    a = esperado_prestamo(1000000, r_iva, 12, mensual(2026, 11, 15, 12), ['S3-convencion-prestamo'])
    b = esperado_prestamo(100000, r_iva, 7, mensual(2027, 1, 31, 7), ['S3-convencion-prestamo'])
    c = esperado_prestamo(100000, r_sin, 7, mensual(2027, 1, 31, 7), ['S3-convencion-prestamo', 'U4-iva-credito'])
    args = lambda P, n, vat, f: {'principal_cents': P, 'annual_rate': tasa(24, 100, 0), 'vat_bp': vat,  # noqa: E731
                                 'term_count': n, 'frequency': 'monthly', 'first_date': f}
    fixture('s3-prestamo-frances',
            'Tabla francesa a mano (Decimal de 80 dígitos, método distinto del oráculo): r = 24 %% / 12 x 1.16 = '
            '0.0232; pago = HALF_UP(P r / (1 - (1 + r)^-n)); interés = HALF_UP(saldo x r); el último pago = saldo + '
            'su interés. 10,000.00 a 12 meses desde el 15-nov-2026; 1,000.00 a 7 meses desde el 31-ene-2027 (fechas '
            'recortadas a fin de mes; el último ajusta: pago fijo %s y último %s); la misma sin IVA (r = 0.02) '
            'lleva además U4-iva-credito. Propiedades de S3-09: suma de capital = capital, saldo final 0.00'
            % (pesos(b['reference_payment_cents']), pesos(b['payments'][-1]['amount_cents'])),
            [src('manual', 'descripcion'), src('criterio', 'S3-09')],
            prov(B_PRESTAMO, 'tasa del periodo = anual / 12 x (1 + IVA) (opción a); otros bloqueos: ' + B_IVA
                 + ' (lectura de vat_bp nulo en la variante sin IVA)'),
            {'today': HOY, 'credits': [], 'cash': []},
            [consulta('prestamo-12', 'amortizar', args(1000000, 12, 1600, '2026-11-15')),
             consulta('prestamo-7-ajuste', 'amortizar', args(100000, 7, 1600, '2027-01-31')),
             consulta('prestamo-7-sin-iva', 'amortizar', args(100000, 7, None, '2027-01-31'))],
            {'prestamo-12': a, 'prestamo-7-ajuste': b, 'prestamo-7-sin-iva': c})
    return b


def tiendas():
    t4 = credito(4, 'Tienda sin interés', 'store_credit', 'axel', None, 'monthly', term_rates=[tr(3, 0, 1, 0)])
    t5 = credito(5, 'Tienda sin datos', 'store_credit', 'axel', None, 'monthly')
    datos = {'today': HOY, 'credits': [estado(t4), estado(t5)], 'cash': []}
    p = compra('2026-10-20', 100000, 3, AX(100000), fpd='2026-11-20')
    fixture('s3-tienda-abonos',
            'Crédito de tienda con 3 abonos sin interés (tasa 0 en term_rates, §5.8): 1,000.00 da 333.33, 333.33 y '
            '333.34 (truncado, el último ajusta), 0.00 de intereses, saldo final 0.00 y fechas mensuales desde la '
            'capturada; el otro crédito sin term_rates ni datos K18 da «falta capturar»; el capital por pago no se '
            'pide (PROVISIONAL S3-serie-sin-intereses)',
            [src('criterio', 'S3-10'), src('decision', 'K18'), src('manual', 'descripcion')], None, datos,
            [consulta('abonos-3', 'cotizar_compra', {'credit_id': 4, 'purchase': p}),
             consulta('sin-datos', 'cotizar_compra', {'credit_id': 5, 'purchase': p})],
            {'abonos-3': {'interest_cents': 0, 'total_cents': 100000, 'source': 'term_table',
                          'payments': [{'due_date': f, 'amount_cents': a} for f, a in
                                       zip(mensual(2026, 11, 20, 3), (33333, 33333, 33334))],
                          'provisional': ['S3-serie-sin-intereses']},
             'sin-datos': {'error': 'missing_capture'}})
    t6 = credito(6, 'Tienda con interés', 'store_credit', 'axel', None, 'monthly', annual=tasa(60, 100, 0),
                 base='monthly', vat=1600, pct=1000, fixed=10000)
    st = corte('2026-11-05', '2026-11-25', 80000, source='computed')
    base = hu(Decimal(80000) * Decimal('0.60') / 12)   # 4,000
    iva = hu(Decimal(base) * Decimal('0.16'))          # 640
    fixture('s3-tienda-abonos-interes',
            'Crédito de tienda con datos K18 de fixture (60 % anual, base mensual, IVA 16 %, mínimo 10 % o 100.00) y '
            'sin term_rates: corte de 800.00 pagado un día después: HALF_UP(800.00 x 0.60 / 12) = 40.00 más IVA 6.40 '
            '= 46.40; mínimo máx(80.00, 100.00) = 100.00; la compra en 3 abonos con interés no tiene regla: unsupported',
            [src('manual', 'descripcion'), src('criterio', 'S3-10')],
            prov(B_TARJETA, 'pregunta 3: fórmula del periodo aplicada a la tienda y abonos con interés sin regla; '
                 'otros bloqueos: ' + B_IVA + ' (datos de «falta capturar» y lectura del IVA)'),
            {'today': HOY, 'credits': [estado(t6)], 'cash': []},
            [consulta('interes', 'interes_tarjeta', {'credit_id': 6, 'statement': st, 'paid_cents': 80000,
                                                     'paid_date': '2026-11-26'}),
             consulta('minimo', 'minimo_tarjeta', {'credit_id': 6, 'statement': st}),
             consulta('compra-3', 'cotizar_compra', {'credit_id': 6, 'purchase': p})],
            {'interes': {'interest_cents': base + iva, 'provisional': P_FORM_L},
             'minimo': {'minimum_cents': 10000, 'provisional': P_FORM_L},
             'compra-3': {'error': 'unsupported'}})


def dia_de_pago(f):
    import calendar
    y, m, d = int(f[:4]), int(f[5:7]), int(f[8:])
    return '%04d-%02d-%02d' % (y, m, 15 if d <= 15 else min(30, calendar.monthrange(y, m)[1]))


def unificado():
    tarjeta = tarjeta_k18(promos=[MSI])
    prest = credito(3, 'Préstamo', 'term_loan', 'axel', None, 'monthly', annual=tasa(24, 100, 0), vat=1600)
    tienda = credito(4, 'Tienda sin interés', 'store_credit', 'axel', None, 'monthly', term_rates=[tr(3, 0, 1, 0)])
    base_cred = [estado(kueski(), [pst('d1', compra_d(0))]),
                 None,
                 estado(prest, [pst('l1', compra('2026-10-15', 100000, 3, AX(100000), fpd='2026-11-15'))]),
                 estado(tienda, [pst('s1', compra('2026-10-20', 60000, 3, AX(60000), fpd='2026-11-20'))])]
    tar_compras = [pst('c1', compra('2026-10-12', 150000, 1, AX(150000))),
                   pst('c2', compra('2026-10-20', 90000, 3, AX(90000), pid=21, merchant='Liverpool'))]
    _, filas_l = amort_mano(100000, Decimal('0.24') / 12 * Decimal('1.16'), 3)
    kq = [(p['fecha'], C(p['monto']), C(p['capital_repuesto']), C(p['axel']), C(p['karen'])) for p in D[0]['pagos']]
    ent = []   # (fecha, crédito, tipo, source, ref, seq, monto, capital, partes, provisional)
    for s, (f, m, k, a, kk) in enumerate(kq, 1):
        ent.append((f, 1, 'revolving_line', 'scheduled', 'd1', s, m, k, [('axel', a), ('karen', kk)], []))
    ent.append(('2026-11-25', 2, 'bank_card', 'scheduled', 'c1', 1, 150000, 150000, [('axel', 150000)], P_CARGO_L))
    for s, f in enumerate(('2026-11-25', '2026-12-25', '2027-01-25'), 1):
        ent.append((f, 2, 'bank_card', 'scheduled', 'c2', s, 30000, 30000, [('axel', 30000)], P_MSI_L))
    for s, (f, row) in enumerate(zip(mensual(2026, 11, 15, 3), filas_l), 1):
        ent.append((f, 3, 'term_loan', 'scheduled', 'l1', s, row['amount_cents'], row['principal_cents'],
                    [('axel', row['amount_cents'])], ['S3-convencion-prestamo']))
    for s, f in enumerate(mensual(2026, 11, 20, 3), 1):
        ent.append((f, 4, 'store_credit', 'scheduled', 's1', s, 20000, 20000, [('axel', 20000)],
                    ['S3-serie-sin-intereses']))

    def esperado(entradas):
        entradas = sorted(entradas, key=lambda e: (e[0], e[1], e[3] != 'scheduled', e[4], e[5]))
        out, por = [], {}
        for f, cid, kind, source, ref, seq, m, k, partes, pl in entradas:
            out.append({'date': f, 'credit_id': cid, 'credit_kind': kind, 'source': source, 'purchase_ref': ref,
                        'seq': seq, 'statement_date': '2026-11-05' if source == 'statement' else None,
                        'amount_cents': m, 'principal_cents': k, 'interest_cents': m - k,
                        'shares': [{'person': p, 'amount_cents': x} for p, x in partes], 'status': 'pending',
                        'provisional': pl})
            q = por.setdefault(dia_de_pago(f), {})
            for p, x in partes:
                q[p] = q.get(p, 0) + x
        return {'entries': out,
                'by_fortnight': [{'pay_day': d, 'total_cents': sum(por[d].values()),
                                  'by_person': [{'person': p, 'amount_cents': por[d][p]} for p in sorted(por[d])]}
                                 for d in sorted(por)]}

    def datos(cortes):
        cr = list(base_cred)
        cr[1] = estado(tarjeta, tar_compras, [], cortes)
        return {'today': HOY, 'credits': cr, 'cash': []}
    q = [consulta('calendario', 'calendario_unificado', {'from': HOY, 'to': '2027-01-31'})]
    otros = (B_PRESTAMO + ' (tasa del periodo del préstamo), %s' + B_SERIE
             + ' (serie de capital de MSI y de la tienda sin interés)')
    fixture('s3-unificado-sin-corte',
            'Calendario unificado del 10-oct-2026 al 31-ene-2027 con los 4 tipos: Kueski d1 (canon), tarjeta (cargo '
            'normal de 1,500.00 del 12-oct y MSI 3 de 900.00 del 20-oct: cortes los 5 y pagos los 25), préstamo de '
            '1,000.00 a 3 meses al 0.0232 (tabla francesa a mano) y tienda sin interés de 600.00 en 3 abonos; orden '
            'por fecha, crédito y pago; suma por quincena = suma de entradas; sin duplicados',
            [src('criterio', 'S3-11'), src('manual', 'descripcion'), src('model.json', 'despensas[0].pagos')],
            prov(B_TARJETA, 'fechas de los pagos de la tarjeta (pregunta 4); otros bloqueos: '
                 + otros % (B_CARGO + ' (source del cargo normal), ')), datos([]), q, {'calendario': esperado(ent)})
    sust = [e for e in ent if not (e[1] == 2 and e[0] == '2026-11-25')]
    sust.append(('2026-11-25', 2, 'bank_card', 'statement', '', 0, 182500, 180000, [('axel', 182500)],
                 ['S3-formula-tarjeta']))
    fixture('s3-unificado',
            'Lo mismo con un corte de banco del 5-nov-2026 (pago el 25-nov) de 1,825.00: sustituye al cargo normal '
            '(1,500.00) y a la primera mensualidad MSI (300.00), con capital 1,800.00 e interés 25.00; suma del '
            'calendario = programados no sustituidos + saldo del corte',
            [src('criterio', 'S3-11'), src('manual', 'descripcion')],
            prov(B_TARJETA, 'fechas de los pagos de la tarjeta (pregunta 4) y el corte que sustituye pagos; otros '
                 'bloqueos: ' + otros % ''),
            datos([corte('2026-11-05', '2026-11-25', 182500)]), q, {'calendario': esperado(sust)})
    fixture('s3-unificado-duplicado', 'Dos cortes de la tarjeta con la misma fecha de corte: entrada inválida',
            [src('criterio', 'S3-11')], None,
            datos([corte('2026-11-05', '2026-11-25', 182500), corte('2026-11-05', '2026-11-26', 182500)]), q,
            {'calendario': {'error': 'invalid'}})


def main(argv=None):
    ap = argparse.ArgumentParser()
    ap.add_argument('--out', default=os.path.join(BASE, 'fixtures', 's3'))
    a = ap.parse_args(argv)
    casos_canon()
    desde_otra_fecha()
    reposicion_separada()
    deuda_viva()
    feb_2027()
    tasa_plazo()
    promocion_0()
    tarjetas()
    prestamo()
    tiendas()
    unificado()
    os.makedirs(os.path.join(a.out, '_datos'), exist_ok=True)
    escritos = {'_datos/kueski-semilla.json': SEMILLA, **SALIDA}
    for nombre, obj in escritos.items():
        with open(os.path.join(a.out, nombre), 'w', encoding='utf-8') as fh:
            json.dump(obj, fh, ensure_ascii=False, indent=2)
            fh.write('\n')
    print('cred_genera_fixtures: %d casos + semilla en %s' % (len(SALIDA), a.out))
    return 0


if __name__ == '__main__':
    sys.exit(main())
