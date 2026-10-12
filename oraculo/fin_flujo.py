"""Oráculo S2 · flujo por quincena y ocurrencias planificadas (S2-interfaz.md §2.7, §3.11; S2-10, M-02, M-10).

FIRME desde la puerta 2 (Q2-14 = A; S2-flujo-quincenal pregunta 3): la liga ocurrencia ↔ captura es una captura con el
mismo planned_item_id en la misma quincena (weekly, fortnightly) o en el mismo mes (monthly, yearly, once). Por eso
`provisional` queda en None en cada quincena; el campo se conserva para marcas futuras.
"""
import datetime as dt
from collections import Counter

from fin_comun import (ErrorMotor, como_datos, lim_mes, lim_quincena, mapa, mes_de, mes_mas, q_mas, q_txt, quincena_de,
                       ultimo)



def ocurrencias(it, a, b):
    """[(fecha, monto)] de una partida en [a, b], nunca antes de start_date ni después de end_date (§2.7)."""
    s, fin, k, f = it['start'], it['end'], it['interval'], it['freq']
    tope = min(b, fin) if fin else b
    res = []

    def mete(d, monto):
        if a <= d <= tope and d >= s:
            res.append((d, monto))
    if f == 'once':
        mete(s, it['amount'])
    elif f == 'weekly':
        d = s
        while d <= tope:
            mete(d, it['amount'])
            d += dt.timedelta(days=7 * k)
    elif f in ('monthly', 'yearly'):
        j = 0
        while True:
            m = mes_mas(mes_de(s), j * k * (12 if f == 'yearly' else 1))
            d = dt.date(m[0], m[1], min(s.day, ultimo(*m)))
            if d > tope:
                break
            mete(d, it['amount'])
            j += 1
    else:  # fortnightly: días de pago 15 y 30 (o último), cada k quincenas desde la de start_date
        j = 0
        q0 = quincena_de(s)
        while True:
            q = q_mas(q0, j * k)
            pay = lim_quincena(q)[2]
            if pay > tope:
                break
            mete(pay, it['second'] if (q[2] == 2 and it['second'] is not None) else it['amount'])
            j += 1
    return res


def capturas(D):
    cap = {}
    for x in list(D.expenses) + list(D.incomes):
        if x.get('pid') is not None:
            cap.setdefault(x['pid'], []).append(x['date'])
    return cap


def ligada(it, d, cap):
    """Liga ocurrencia ↔ captura (FIRME, Q2-14 = A)."""
    if it['freq'] in ('weekly', 'fortnightly'):
        a, b, _ = lim_quincena(quincena_de(d))
    else:
        a, b = lim_mes(mes_de(d))
    return any(a <= x <= b for x in cap.get(it['id'], []))


def partidas_vivas(D):
    return [it for it in D.planned if it['active'] and it['credit_id'] is None]


def flujo_quincenal(datos, desde, n):
    D = como_datos(datos)
    if not 1 <= n <= 24:
        raise ErrorMotor('invalid_argument', 'count')
    T, cap, filas, acum = D.today, capturas(D), [], 0
    vivas = {i: c for i, c in D.purchases.items() if c['status'] != 'cancelled'}
    for j in range(n):
        q = q_mas(desde, j)
        a, b, pay = lim_quincena(q)
        fijo = sum(x['amount'] for x in D.incomes if x['kind'] == 'fixed' and a <= x['date'] <= b)
        extra = sum(x['amount'] for x in D.incomes if x['kind'] == 'extra' and a <= x['date'] <= b)
        esperado = plan_fu = pend_occ = 0
        for it in partidas_vivas(D):
            for d, monto in ocurrencias(it, a, b):
                lig = ligada(it, d, cap)
                if it['kind'] == 'income':
                    esperado += 0 if lig else monto
                else:
                    plan_fu += monto
                    pend_occ += 0 if lig else monto
        by, inf, pagado = Counter(), 0, 0
        for c in vivas.values():
            for p in c['payments']:
                if a <= p['due'] <= b:
                    for per, v in p['shares'].items():
                        by[per] += v
                        if per in D.directo:
                            inf += v
                    h = p['shares'].get(D.holder, 0)
                    abonos = sum(x['amount'] for x in p['paid'] if x['person'] == D.holder and x['date'] <= T)
                    pagado += min(h, abonos)
        hogar = by.get(D.holder, 0)

        def caja_hogar(kinds, cond):
            return sum(g['amount'] for g in D.expenses if g['method'] != 'credit' and g['kind'] in kinds and a <= g['date'] <= b
                       and cond(g['date']) and g['paid_by'] == D.holder)
        fu_pag = caja_hogar(('fixed', 'upcoming'), lambda d: d <= T)
        fu_pend = pend_occ + caja_hogar(('fixed', 'upcoming'), lambda d: d > T)
        var = caja_hogar(('variable',), lambda d: d <= T)
        total_ing = fijo + extra + esperado
        gastado = var + fu_pag + pagado
        comprometido = (hogar - pagado) + fu_pend
        saldo = total_ing - gastado - comprometido
        acum += saldo
        filas.append(dict(fortnight=q_txt(q), start=a, end=b, pay_date=pay,
                          income=dict(fixed_cents=fijo, extra_cents=extra, expected_cents=esperado, income_is_expected=esperado > 0,
                                      total_cents=total_ing),
                          credit_payments=dict(by_person=mapa(D, by), household_cents=hogar, external_informative_cents=inf,
                                               household_paid_cents=pagado, household_pending_cents=hogar - pagado),
                          fixed_upcoming=dict(planned_cents=plan_fu, paid_cents=fu_pag, pending_cents=fu_pend),
                          # None mientras sigan PROVISIONALES: S2-flujo-quincenal pregunta 1 (Q2-12 sin respuesta) y S2-ahorro-quincenal.
                          variable_spent_cents=var, variable_estimated_cents=None, savings_planned_cents=None,
                          spent_cents=gastado, committed_cents=comprometido, balance_cents=saldo, cumulative_balance_cents=acum,
                          provisional=None))
    return filas
