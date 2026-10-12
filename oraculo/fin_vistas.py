"""Oráculo S2 · vistas con las claves literales de model.json (S2-interfaz.md §3.14)."""
from fin_comun import MESES, ErrorMotor, como_datos, lim_mes
import fin_devengo as dv


def _forma(D):
    if len(D.externas) != 1:
        raise ErrorMotor('oracle_view_shape', 'se espera un titular y una sola persona externa')
    return D.holder, D.externas[0]


def vista_periodo(datos, m):
    D = como_datos(datos)
    h, e = _forma(D)
    s, x = dv.periodos_int(D, meses=[m])[0]
    a, b = lim_mes(m)
    ac, ca = s['accrued'], s['cash']
    return {'despensa_que_cubre': s['covering_groceries_seq'],
            'devengado': {'despensa_total': ac['by_category'].get(D.groceries, 0), 'papel_contado': x['paper_accrued'],
                          'total': ac['total_cents'], h: ac['by_person'][h], e: ac['by_person'][e]},
            'efectivo': {'pagos': ca['payments_cents'], 'papel': x['paper_cash'], 'total': ca['total_cents'],
                         h: ca['by_person'][h], e: ca['by_person'][e], 'n_pagos': ca['payment_count'],
                         'detalle_pagos': [{'despensa': p['purchase_seq'], 'fecha': p['due_date'], 'monto': p['amount_cents'],
                                            h: p['shares'][h], e: p['shares'][e]} for p in ca['payments']]},
            'diferencia_acumulada': s['cumulative_diff_cents'], 'etiqueta': s['label'], 'inicio': a, 'fin': b}


def vista_horizonte(datos, inicio, fin, meses):
    D = como_datos(datos)
    _forma(D)
    r = dv.horizonte(D, inicio, fin, meses)
    out = {}
    for p in r['payments_by_purchase']:
        if p['purchase_seq'] is not None:
            k = f"pagos_despensa_{p['purchase_seq']}" + ('' if p['complete'] else '_en_horizonte')
            out[k] = p['in_horizon_cents']
    lv = r['live_debt']
    out.update({'papel_contado': r['paper_cash_cents'], 'total_efectivo': r['cash_total_cents'],
                'suma_periodos_efectivo': r['sum_of_periods_cash_cents'], 'total_devengado': r['accrued_total_cents'],
                f'deuda_viva_{fin.day}_{MESES[fin.month - 1]}': lv['live_debt_cents'],
                'deuda_viva_detalle': {(f"despensa_{p['purchase_seq']}_pendiente" if p['started'] else f"despensa_{p['purchase_seq']}"): p['pending_cents']
                                       for p in lv['by_purchase'] if p['purchase_seq'] is not None},
                'promedio_mensual_efectivo': r['monthly_avg_cash_cents'], 'promedio_mensual_devengado': r['monthly_avg_accrued_cents'],
                'diferencia_final': r['final_diff_cents']})
    return out


def vista_valor_estable(datos, items=None):
    D = como_datos(datos)
    h, e = _forma(D)
    r = dv.valor_estable_de(D, items)
    comp = {c['label']: c for c in r['components']}
    if set(comp) != {'despensa', 'papel_contado'}:
        raise ErrorMotor('oracle_view_shape', 'partidas despensa y papel_contado')
    m, q = r['month'], r['fortnight']
    return {'mes': {'total': m['total_cents'], h: m['by_person'][h], e: m['by_person'][e]},
            'quincena': {'total': q['total_cents'], h: q['by_person'][h], e: q['by_person'][e],
                         f'suma_{h}_{e}': q['sum_persons_cents'], 'diferencia_centavo': q['cent_diff_cents']},
            'componentes': {'mitad_despensa': comp['despensa']['monthly_cents'], 'papel_contado_mes': comp['papel_contado']['monthly_cents'],
                            f'{e}_mitad_despensa': comp['despensa']['by_person'][e], f'{h}_mitad_despensa': comp['despensa']['by_person'][h]}}


def vista_papel(datos):
    D = como_datos(datos)
    _forma(D)
    return [{'paquete': p['package'], 'compra': p['date'], 'quien': p['person_name'],
             'forma': f"dentro de despensa {p['purchase_seq']}" if p['form'] == 'in_groceries' else 'de contado',
             'monto': p['amount_cents']} for p in dv.papel(D)]


def vista_sensibilidad(datos, items):
    D = como_datos(datos)
    h, e = _forma(D)
    r = dv.valor_estable_de(D, items)
    return {'costo_mensual_total': r['month']['total_cents'], f'{h}_mensual': r['month']['by_person'][h],
            f'{e}_mensual': r['month']['by_person'][e]}


def vista_interes_por_monto(datos, purchase_id):
    D = como_datos(datos)
    _forma(D)
    c = D.purchases.get(purchase_id)
    if c is None:
        raise ErrorMotor('invalid_argument', 'purchase_id')
    lb = dv.libro(D, D.expenses, D.purchases)
    dev = sum(t for x in lb['costos'] if x['pid'] == purchase_id for t, _ in dv.porciones(D, x['total'], x['parts'], x['n']))
    caja = sum(e['total'] for e in lb['caja'] if e['pago'] is not None and e['pago']['purchase_id'] == purchase_id)
    if dev != caja:
        raise ErrorMotor('sum_mismatch', f'compra {purchase_id}: devengado {dev} != efectivo {caja}')
    uno = next((p for p in c['payments'] if p['seq'] == 1), None)
    return {'monto': c['principal'], 'intereses': dev - c['principal'], 'total': dev,
            'pago_quincenal': uno['amount'] if uno else None}
