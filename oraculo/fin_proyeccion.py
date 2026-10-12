"""Oráculo S2 · proyección de 12 meses (S2-interfaz.md §3.12; S2-09, RN-28, K21, M-11).

Supuesto del oráculo (no cambia cifras canon): las ocurrencias planificadas se generan solo dentro de los 12 meses
proyectados. PROVISIONAL: cadencia encendida (S2-proyeccion-interruptor) y liga (S2-flujo-quincenal).
"""
from collections import Counter

from fin_comun import ErrorMotor, como_datos, lim_mes, mes_de, mes_mas, mes_txt, reparto
import fin_devengo as dv
from fin_flujo import BLOQUEO_LIGA, capturas, ligada, ocurrencias, partidas_vivas

BLOQUEO_CADENCIA = 'S2-proyeccion-interruptor'


def proyeccion(datos, desde=None):
    D = como_datos(datos)
    m0 = desde or mes_de(D.today)
    meses = [mes_mas(m0, j) for j in range(12)]
    a12, b12 = lim_mes(meses[0])[0], lim_mes(meses[-1])[1]
    cad = D.repeat
    if cad:
        m7 = lim_mes(meses[6])[0]
        for x in list(D.cad_purchases.values()):
            if x['date'] < m7:
                raise ErrorMotor('cadence_before_month_7', f"cadence.purchases id {x['id']}")
        for g in D.cad_expenses:
            if g['date'] < m7:
                raise ErrorMotor('cadence_before_month_7', f"cadence.expenses id {g['id']}")
    cap = capturas(D)
    marcados, ing_occ, plan_m, sinteticos = set(), Counter(), Counter(), []
    for it in partidas_vivas(D):
        for d, monto in ocurrencias(it, a12, b12):
            if ligada(it, d, cap):
                marcados.add(mes_de(d))
                continue
            if it['kind'] == 'income':
                ing_occ[mes_de(d)] += monto
            else:
                plan_m[mes_de(d)] += monto
                sinteticos.append(dict(id=None, kind='fixed' if it['kind'] == 'fixed_expense' else 'upcoming', cat=it['cat'],
                                       amount=monto, date=d, n=it['n'], method='cash', paid_by=D.holder,
                                       shares=reparto(it['rule'], monto), cpid=None, pid=it['id']))
    gastos = list(D.expenses) + (list(D.cad_expenses) if cad else []) + sinteticos
    compras = dict(D.purchases)
    if cad:
        compras.update(D.cad_purchases)
    lb = dv.libro(D, gastos, compras)
    res = dv.resumen(D, dv.tramos_mes(meses), 'mes', lb)
    out = []
    for j, (s, _) in enumerate(res):
        m = meses[j]
        a, b = lim_mes(m)
        ing = sum(x['amount'] for x in D.incomes if a <= x['date'] <= b) + ing_occ[m]
        prov = sorted(([BLOQUEO_LIGA] if m in marcados else []) + ([BLOQUEO_CADENCIA] if cad and j >= 6 else [])) or None
        out.append(dict(index=j + 1, month=mes_txt(m), accrued=s['accrued'], cash=s['cash'], income_cents=ing,
                        planned_expense_cents=plan_m[m], net_cash_cents=ing - s['cash']['by_person'][D.holder],
                        cumulative_diff_cents=s['cumulative_diff_cents'], provisional=prov))
    return dict(repeat_cadence_enabled=cad, months=out)
