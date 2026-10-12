"""Oráculo S2 · presupuesto (S2-interfaz.md §3.8–§3.10; RN-22, K2, M-01, M-09, M-14, S2-11)."""
import datetime as dt
from collections import Counter
from fractions import Fraction

from fin_comun import ErrorMotor, como_datos, hu, lim_mes, mes_de, mes_mas, mes_txt, renglones, ultimo
import fin_devengo as dv


def _renglon(budgets, cat, m):
    a, b = lim_mes(m)
    for r in budgets:
        if r['category'] == cat and r['valid_from'] <= a and (r['valid_to'] is None or b <= r['valid_to']):
            return r
    return None


def _banderas(D, actual, plan):
    """used_pct = HALF_UP(actual·100/plan); banderas con comparación exacta en enteros (S2-11)."""
    if plan is None or plan == 0:  # plan 0: sin porcentaje ni banderas (FIRME, Q2-20 = A)
        return None, False, False, 'no_plan'
    used = hu(Fraction(actual * 100, plan))
    w = actual * 100 >= D.warn * plan
    lim = actual * 100 >= D.limit * plan
    return used, w, lim, ('exceeded' if lim else 'warning' if w else 'normal')


def _hogar_por_categoria(D, m):
    """Devengado de la parte del hogar (titular) por categoría en el mes m, día por día desde el día 1."""
    lb = dv.libro(D, D.expenses, D.purchases)
    acc = Counter()
    for c in lb['costos']:
        for j, (_, by) in enumerate(dv.porciones(D, c['total'], c['parts'], c['n'])):
            if mes_mas(mes_de(c['date']), j) == m:
                acc[c['cat']] += by.get(D.holder, 0)
    return acc


def presupuesto_mes(datos, m):
    D = como_datos(datos)
    a, b = lim_mes(m)
    hog = _hogar_por_categoria(D, m)
    filas, tp, ta = [], 0, 0
    for code in sorted((k for k, v in D.categories.items() if v['active']), key=lambda s: s.encode()):
        r = _renglon(D.budgets, code, m)
        act = hog.get(code, 0)
        if r is None:
            filas.append(dict(category=code, plan_cents=None, actual_cents=act, remaining_cents=None, used_pct=None,
                              warning=False, limit=False, state='no_plan'))
        else:
            used, w, lim, st = _banderas(D, act, r['amount'])
            filas.append(dict(category=code, plan_cents=r['amount'], actual_cents=act, remaining_cents=r['amount'] - act,
                              used_pct=used, warning=w, limit=lim, state=st))
            tp += r['amount']
        ta += act
    used, w, lim, st = _banderas(D, ta, tp)
    tot = _renglon(D.budgets, None, m)
    per = dv.periodos(D, meses=[m])[0]
    T = D.today
    quedan = 0 if T > b else (ultimo(*m) if T < a else (b - T).days + 1)
    return dict(month=mes_txt(m), days_in_month=ultimo(*m), days_left=quedan, categories=filas,
                total=dict(plan_cents=tp, actual_cents=ta, remaining_cents=tp - ta, used_pct=used, warning=w, limit=lim, state=st,
                           declared_total_row_cents=tot['amount'] if tot else None),
                accrued_household_cents=per['accrued']['by_person'][D.holder],
                cash_household_cents=per['cash']['by_person'][D.holder])


def aplicar_cambio_presupuesto(renglones_raw, cambio):
    """§3.9 con las reglas de updateBudget (K2). `renglones_raw`: BudgetRow de entradas; `cambio`: dict validado."""
    rows = [dict(r) for r in (renglones_raw if renglones_raw and isinstance(renglones_raw[0].get('valid_from'), dt.date) else renglones(renglones_raw))]
    m, monto, alcance = cambio['month'], cambio['amount_cents'], cambio.get('scope') or 'onward'
    r = next((x for x in rows if x['id'] == cambio['row_id']), None)
    if r is None:
        raise ErrorMotor('invalid_argument', 'row_id')
    a, b = lim_mes(m)
    if not (r['valid_from'] <= a and (r['valid_to'] is None or b <= r['valid_to'])):
        raise ErrorMotor('out_of_range', 'month fuera de la vigencia')
    sig = max(x['id'] for x in rows) + 1
    borrados = []
    fin_original = r['valid_to']
    if alcance == 'onward':
        borrados = sorted(x['id'] for x in rows if x['category'] == r['category'] and x['valid_from'] > r['valid_from'])
        rows = [x for x in rows if x['id'] not in borrados]
        if r['valid_from'] == a:
            r['amount'], r['valid_to'] = monto, None
            efectivo = r['id']
        else:
            r['valid_to'] = a - dt.timedelta(days=1)
            rows.append(dict(id=sig, category=r['category'], amount=monto, valid_from=a, valid_to=None))
            efectivo = sig
    else:
        anterior = r['amount']
        if r['valid_from'] == a:
            r['amount'], r['valid_to'] = monto, b
            efectivo = r['id']
        else:
            r['valid_to'] = a - dt.timedelta(days=1)
            rows.append(dict(id=sig, category=r['category'], amount=monto, valid_from=a, valid_to=b))
            efectivo = sig
            sig += 1
        if fin_original is None or fin_original > b:
            rows.append(dict(id=sig, category=r['category'], amount=anterior, valid_from=b + dt.timedelta(days=1), valid_to=fin_original))
    rows.sort(key=lambda x: (x['category'] is not None, (x['category'] or '').encode(), x['valid_from']))
    return dict(rows=[dict(id=x['id'], category=x['category'], amount_cents=x['amount'], valid_from=x['valid_from'], valid_to=x['valid_to']) for x in rows],
                effective_row_id=efectivo, deleted_row_ids=borrados)


def plan_contra_real(datos, desde, n, categoria=None):
    D = como_datos(datos)
    if not 1 <= n <= 12:
        raise ErrorMotor('invalid_argument', 'count')
    if categoria is not None and not (categoria in D.categories and D.categories[categoria]['active']):
        raise ErrorMotor('invalid_argument', 'category')
    out = []
    for j in range(n):
        m = mes_mas(desde, j)
        p = presupuesto_mes(D, m)
        fila = p['total'] if categoria is None else next(x for x in p['categories'] if x['category'] == categoria)
        out.append(dict(month=mes_txt(m), plan_cents=fila['plan_cents'], actual_cents=fila['actual_cents']))
    return out
