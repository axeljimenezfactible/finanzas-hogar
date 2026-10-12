"""Oráculo S2 · metas y fondos (S2-interfaz.md §3.13; RN-18, S2-12, M-07, M-16).

FIRMES desde la puerta 2: months_left (Q2-16 = A: del mes de hoy al mes de target_date, los dos incluidos, mínimo 0),
`state` (Q2-18 = A: achieved, overdue o in_progress) y el resumen solo con metas activas (Q2-15 = A).
PROVISIONAL (docs/estado/bloqueos/S2-metas-vencida.md): una meta vencida sin cumplir (months_left 0) pone todo el
restante en el mes de hoy (plan de un renglón y necesaria = restante), la regla en uso.
"""
from fractions import Fraction

from fin_comun import como_datos, hu, mes_de, mes_mas, mes_n, mes_txt


def metas(datos):
    D = como_datos(datos)
    T = D.today
    mT = mes_de(T)
    filas, ahorro, aportado, planeado = [], 0, 0, 0
    for g in sorted((x for x in D.goals if x['active']), key=lambda x: x['id']):
        movs = [x for x in D.movements if x['goal'] == g['id'] and x['date'] <= T]
        saved = sum(x['amount'] if x['kind'] == 'deposit' else -x['amount'] for x in movs)
        aportado += sum(x['amount'] for x in movs if x['kind'] == 'deposit' and mes_de(x['date']) == mT)
        tgt, fd = g['target'], g['target_date']
        decimas = hu(Fraction(saved * 1000, tgt)) if tgt else None
        pct = hu(Fraction(saved * 100, tgt)) if tgt else None
        faltan = max(0, mes_n(mes_de(fd)) - mes_n(mT) + 1) if fd else None
        plan, necesaria = [], None
        if tgt is not None and fd is not None:
            if saved >= tgt:
                necesaria = 0
            elif faltan == 0:  # vencida sin cumplir: PROVISIONAL (S2-metas-vencida)
                plan = [dict(month=mes_txt(mT), cents=tgt - saved)]
                necesaria = tgt - saved
            else:
                resto = tgt - saved
                base = hu(Fraction(resto, faltan))
                plan = [dict(month=mes_txt(mes_mas(mT, j)), cents=base) for j in range(faltan - 1)]
                plan.append(dict(month=mes_txt(mes_mas(mT, faltan - 1)), cents=resto - (faltan - 1) * base))
                necesaria = base
        estado = 'achieved' if (tgt is not None and saved >= tgt) else ('overdue' if (fd is not None and fd < T) else 'in_progress')
        filas.append(dict(goal_id=g['id'], kind=g['kind'], saved_cents=saved, target_cents=tgt, progress_tenths_pct=decimas,
                          progress_pct=pct, months_left=faltan, needed_monthly_cents=necesaria, plan=plan, state=estado))
        ahorro += saved
        planeado += g['planned'] or 0
    return dict(month_summary=dict(saved_cents=ahorro, contributed_this_month_cents=aportado, planned_this_month_cents=planeado),
                goals=filas)
