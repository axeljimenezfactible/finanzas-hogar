"""Oráculo S2 · devengado y efectivo (S2-interfaz.md §1.3, §3.2–§3.7; RN-02, RN-21, K19, K20, K39).

Método propio (S2-03, §0.4): arma un libro de costos y movimientos y recorre el calendario día por día;
los meses se devengan al pasar por el día 1 (o por el inicio de cada ventana). Sumas exactas en enteros;
HALF_UP solo donde lo dice §1.3.
"""
from collections import Counter, defaultdict
from fractions import Fraction

from fin_comun import (ErrorMotor, como_datos, dias, etiqueta, hu, lim_mes, mapa, mes_de, mes_mas, mes_n, partidas_estables,
                       por_categoria, ultimo)


def porciones(D, total, partes, n):
    """Devengo de un costo en n meses (RN-02, §1.3): HALF_UP(C/n) y el último ajusta; externas igual; el titular, el resto."""
    base = hu(Fraction(total, n))
    otros = [(p, v) for p, v in partes.items() if p != D.holder]
    res = []
    for k in range(n):
        ult = k == n - 1
        t = total - (n - 1) * base if ult else base
        by = {}
        for p, v in otros:
            b = hu(Fraction(v, n))
            by[p] = v - (n - 1) * b if ult else b
        by[D.holder] = t - sum(by.values())
        res.append((t, by))
    return res


def libro(D, gastos, compras):
    """Costos devengables y movimientos de efectivo de gastos y compras no canceladas."""
    vivas = {i: c for i, c in compras.items() if c['status'] != 'cancelled'}
    por_id = {g['id']: g for g in gastos if g.get('id') is not None}
    cat = {}
    for i, c in vivas.items():
        if c['expense_id'] is not None:
            g = por_id.get(c['expense_id'])
            if g is None:
                raise ErrorMotor('invalid_argument', f'compra {i}: expense_id sin gasto')
            cat[i] = g['cat']
        else:
            cat[i] = c['cat']
    orden = sorted((c['date'], i) for i, c in vivas.items() if cat[i] == D.groceries)
    seq = {i: k + 1 for k, (_, i) in enumerate(orden)}
    costos, caja = [], []
    for g in gastos:
        papel = g['cat'] == D.paper
        if g['method'] == 'credit':
            c = compras[g['cpid']]
            if c['status'] == 'cancelled':
                continue
            costos.append(dict(date=g['date'], n=g['n'], total=c['total'], parts=c['share_tot'], cat=g['cat'],
                               seq=seq.get(c['id']), paper_cash=False, pid=c['id']))
        else:
            costos.append(dict(date=g['date'], n=g['n'], total=g['amount'], parts=g['shares'], cat=g['cat'],
                               seq=None, paper_cash=papel, pid=None))
            caja.append(dict(date=g['date'], total=g['amount'], by={g['paid_by']: g['amount']}, cat=g['cat'],
                             pago=None, paper_cash=papel, gasto=g))
    for i, c in vivas.items():
        if c['expense_id'] is None:
            costos.append(dict(date=c['date'], n=c['n'], total=c['total'], parts=c['share_tot'], cat=c['cat'],
                               seq=seq.get(i), paper_cash=False, pid=i))
        for p in c['payments']:
            caja.append(dict(date=p['due'], total=p['amount'], by=p['shares'], cat=cat[i], paper_cash=False, gasto=None,
                             pago=dict(purchase_id=i, purchase_seq=seq.get(i), seq=p['seq'], due_date=p['due'],
                                       amount_cents=p['amount'], shares=p['shares'])))
    return dict(costos=costos, caja=caja, seq=seq, cat=cat, vivas=vivas)


def sin_cubrir(tramos, d):
    """d cae antes de la primera ventana o en un hueco entre ventanas: no está en ninguna y no es posterior a la última.

    FIRME desde la puerta 2 (§2.12 regla 2, Q2-25 = A: desde el primer evento; regla 3, Q2-26 = H1: sin eventos en un
    hueco): un costo o un pago así es invalid_argument. Lo posterior a la última ventana no entra en ningún tramo.
    Con la regla 4 (contiguas, PROVISIONAL) ya no queda hueco, pero la regla 3 se comprueba aparte porque es FIRME."""
    if any(a <= d <= b for _, a, b in tramos):
        return False
    return d <= max(b for _, _, b in tramos)


def desordenadas(tramos):
    """Etiqueta de la primera ventana que no empieza después de que terminaron todas las anteriores; None si van bien.

    FIRME desde la puerta 2 (§2.12 regla 1, Q2-27 = O1): fuera de orden, traslapadas (aunque sea un día) o una dentro de
    otra se rechazan. Se compara cada inicio con el fin más tardío visto hasta ahí."""
    ultimo_fin = None
    for etiqueta_v, inicio, fin in tramos:
        if ultimo_fin is not None and inicio <= ultimo_fin:
            return etiqueta_v
        ultimo_fin = fin if ultimo_fin is None else max(ultimo_fin, fin)
    return None


def con_hueco(tramos):
    """Etiqueta de la primera ventana que no empieza justo el día siguiente al fin de la anterior; None si no hay hueco.

    PROVISIONAL (§2.12 regla 4; docs/estado/bloqueos/S2-ventanas.md, pregunta 4, caso 1; Q2-01 = A): toda lista con
    hueco es invalid_argument aunque el hueco no tenga eventos, porque un devengo que lo cruza mandaría la parte del
    hueco a la ventana siguiente (REVISION-S2-r3 M1). Se mide en días: la distancia entre fin e inicio debe ser 1."""
    for (_, _, fin), (etiqueta_v, inicio, _) in zip(tramos, tramos[1:]):
        if (inicio - fin).days != 1:
            return etiqueta_v
    return None


def resumen(D, tramos, modo, lb):
    """Resumen por tramo (§3.2). tramos: [(label, inicio, fin)]; modo 'mes' o 'ventana'. Devuelve [(salida, extra)]."""
    if modo == 'ventana':
        desordenada = desordenadas(tramos)
        if desordenada is not None:
            raise ErrorMotor('invalid_argument', f'windows: la ventana {desordenada!r} no empieza después de todas las anteriores; '
                                                 'van ascendentes y disjuntas (§2.12 regla 1)')
        hueco = con_hueco(tramos)
        if hueco is not None:
            raise ErrorMotor('invalid_argument', f'windows: la ventana {hueco!r} no empieza el día siguiente al fin de la '
                                                 'anterior; van contiguas (§2.12 regla 4, PROVISIONAL, S2-ventanas)')
        for tipo, lista in (('costo', lb['costos']), ('pago', lb['caja'])):
            for e in lista:
                if sin_cubrir(tramos, e['date']):
                    raise ErrorMotor('invalid_argument', f"windows: {tipo} del {e['date']} antes de la primera ventana o en un "
                                                         'hueco entre ventanas (§2.12 reglas 2 y 3)')
    asignado = defaultdict(list)
    for c in lb['costos']:
        if modo == 'mes':
            k0 = mes_n(mes_de(c['date']))
        else:
            # La parte j sigue a la ventana y no al mes, aunque la ventana no dure un mes: PROVISIONAL (S2-ventanas,
            # pregunta 4, caso 2; §3.2), regla en uso sin decidir.
            k0 = next((i for i, (_, a, b) in enumerate(tramos) if a <= c['date'] <= b), None)
            if k0 is None:
                continue  # posterior a la última ventana (sin_cubrir ya rechazó lo anterior y los huecos)
        for j, (t, by) in enumerate(porciones(D, c['total'], c['parts'], c['n'])):
            asignado[k0 + j].append((t, by, c))
    caja_dia = defaultdict(list)
    for e in lb['caja']:
        caja_dia[e['date']].append(e)
    out = []
    for i, (label, a, b) in enumerate(tramos):
        at, ap, ac, apap, cubre = 0, Counter(), Counter(), 0, None
        ct, cpag, cgas, cp, cc, cpap, pagos = 0, 0, 0, Counter(), Counter(), 0, []
        for d in dias(a, b):
            claves = [mes_n(mes_de(d))] if (modo == 'mes' and d.day == 1) else ([i] if (modo == 'ventana' and d == a) else [])
            for kk in claves:
                for t, by, c in asignado.get(kk, []):
                    at += t
                    for p, v in by.items():
                        ap[p] += v
                    ac[c['cat']] += t
                    if c['paper_cash']:
                        apap += t
                    if c['seq'] is not None and (cubre is None or c['seq'] > cubre):
                        cubre = c['seq']
            for e in caja_dia.get(d, []):
                ct += e['total']
                for p, v in e['by'].items():
                    cp[p] += v
                cc[e['cat']] += e['total']
                if e['pago'] is not None:
                    cpag += e['total']
                    pagos.append(e['pago'])
                else:
                    cgas += e['total']
                if e['paper_cash']:
                    cpap += e['total']
        lim = mes_n(mes_de(b)) if modo == 'mes' else i
        acum_dev = sum(t for k, lst in asignado.items() if k <= lim for t, _, _ in lst)
        acum_caja = sum(e['total'] for e in lb['caja'] if e['date'] <= b)
        pagos.sort(key=lambda p: (p['due_date'], p['purchase_id'], p['seq']))
        sal = dict(label=label, start=a, end=b, days=(b - a).days + 1,
                   accrued=dict(total_cents=at, by_person=mapa(D, ap), by_category=por_categoria(ac)),
                   cash=dict(total_cents=ct, payments_cents=cpag, cash_expenses_cents=cgas, by_person=mapa(D, cp),
                             by_category=por_categoria(cc), payment_count=len(pagos),
                             payments=[dict(purchase_id=p['purchase_id'], purchase_seq=p['purchase_seq'], seq=p['seq'],
                                            due_date=p['due_date'], amount_cents=p['amount_cents'], shares=mapa(D, p['shares']))
                                       for p in pagos]),
                   cumulative_diff_cents=acum_dev - acum_caja, covering_groceries_seq=cubre)
        out.append((sal, dict(paper_accrued=apap, paper_cash=cpap)))
    return out


def tramos_mes(meses):
    res = []
    for m in meses:
        a, b = lim_mes(m)
        res.append((etiqueta(a, b), a, b))
    return res


def periodos_int(D, meses=None, ventanas=None, lb=None):
    lb = lb or libro(D, D.expenses, D.purchases)
    if meses is not None:
        return resumen(D, tramos_mes(meses), 'mes', lb)
    return resumen(D, list(ventanas), 'ventana', lb)


def periodos(datos, meses=None, ventanas=None):
    """§3.2: un resumen por mes (lista de (año, mes)) o por ventana (lista de (label, inicio, fin))."""
    D = como_datos(datos)
    if (meses is None) == (ventanas is None):
        raise ErrorMotor('invalid_argument', 'exactamente uno de months/windows')
    return [s for s, _ in periodos_int(D, meses, ventanas)]


def deuda_viva(datos, corte):
    """§3.4: deuda viva (RN-21) e identidad de K19 con devengado de meses cerrados."""
    D = como_datos(datos)
    lb = libro(D, D.expenses, D.purchases)
    viva, detalle = 0, []
    for i, c in sorted(lb['vivas'].items(), key=lambda kv: (kv[1]['date'], kv[0])):
        if c['date'] > corte:
            continue
        pend = sum(p['amount'] for p in c['payments'] if p['due'] > corte)
        viva += pend
        if pend > 0:
            detalle.append(dict(purchase_id=i, purchase_seq=lb['seq'].get(i), pending_cents=pend,
                                started=any(p['due'] <= corte for p in c['payments'])))
    dev = 0
    for c in lb['costos']:
        for j, (t, _) in enumerate(porciones(D, c['total'], c['parts'], c['n'])):
            if lim_mes(mes_mas(mes_de(c['date']), j))[1] <= corte:
                dev += t
    caja = sum(e['total'] for e in lb['caja'] if e['date'] <= corte)
    saldo = sum(c['total'] for c in lb['costos'] if c['date'] <= corte) - dev
    dif = dev - caja
    return dict(cut=corte, live_debt_cents=viva, by_purchase=detalle, accrued_to_date_cents=dev, cash_to_date_cents=caja,
                cumulative_diff_cents=dif, unaccrued_balance_cents=saldo, identity_holds=viva == dif + saldo,
                identity_without_balance_holds=viva == dif)


def horizonte(datos, inicio, fin, meses):
    """§3.3 (R-51…R-60)."""
    D = como_datos(datos)
    lb = libro(D, D.expenses, D.purchases)
    pagos = []
    for i, c in sorted(lb['vivas'].items(), key=lambda kv: (kv[1]['date'], kv[0])):
        dentro = [p for p in c['payments'] if inicio <= p['due'] <= fin]
        if dentro:
            pagos.append(dict(purchase_id=i, purchase_seq=lb['seq'].get(i), in_horizon_cents=sum(p['amount'] for p in dentro),
                              complete=len(dentro) == len(c['payments'])))
    papel = sum(e['total'] for e in lb['caja'] if e['paper_cash'] and inicio <= e['date'] <= fin)
    caja = sum(e['total'] for e in lb['caja'] if inicio <= e['date'] <= fin)
    per = periodos(D, meses=meses)
    suma_caja = sum(p['cash']['total_cents'] for p in per)
    dev = sum(p['accrued']['total_cents'] for p in per)
    return dict(payments_by_purchase=pagos, paper_cash_cents=papel, cash_total_cents=caja, sum_of_periods_cash_cents=suma_caja,
                accrued_total_cents=dev, live_debt=deuda_viva(D, fin),
                monthly_avg_cash_cents=hu(Fraction(caja, len(meses))), monthly_avg_accrued_cents=hu(Fraction(dev, len(meses))),
                final_diff_cents=dev - caja)


def promedios(total, inicio, fin):
    """§3.5 (K20, K39): diario = HALF_UP(T/d); semanal = HALF_UP(7·T/d), un solo redondeo."""
    if fin < inicio:
        raise ErrorMotor('invalid_argument', 'start > end')
    d = (fin - inicio).days + 1
    return dict(days=d, daily_cents=hu(Fraction(total, d)), weekly_cents=hu(Fraction(7 * total, d)))


def promedios_mes(datos, m, base):
    D = como_datos(datos)
    p = periodos(D, meses=[m])[0]
    if base == 'cash':
        t = p['cash']['total_cents']
    elif base == 'accrued':
        t = p['accrued']['total_cents']
    elif base == 'accrued_groceries':
        t = p['accrued']['by_category'].get(D.groceries, 0)
    else:
        raise ErrorMotor('invalid_argument', 'basis')
    a, b = lim_mes(m)
    return promedios(t, a, b)


def promedio_mensual(total, meses):
    if meses < 1:
        raise ErrorMotor('invalid_argument', 'month_count')
    return dict(monthly_cents=hu(Fraction(total, meses)))


def valor_estable(partidas, personas):
    """§3.6 (R-37…R-44): cada partida aporta HALF_UP(total/ciclo); quincena = HALF_UP(x/2) por renglón."""
    D = personas
    comps, mt, mp = [], 0, Counter()
    for it in partidas:
        t = hu(Fraction(it['total'], it['cycle']))
        by = {p: hu(Fraction(v, it['cycle'])) for p, v in it['shares'].items() if p != D.holder}
        by[D.holder] = t - sum(by.values())
        comps.append(dict(label=it['label'], monthly_cents=t, by_person=mapa(D, by)))
        mt += t
        for p, v in by.items():
            mp[p] += v
    qp = {p: hu(Fraction(mp.get(p, 0), 2)) for p in D.claves}
    qt = hu(Fraction(mt, 2))
    return dict(month=dict(total_cents=mt, by_person=mapa(D, mp)),
                fortnight=dict(total_cents=qt, by_person=mapa(D, qp), sum_persons_cents=sum(qp.values()),
                               cent_diff_cents=qt - sum(qp.values())),
                components=comps)


def valor_estable_de(datos, items=None):
    D = como_datos(datos)
    partidas = D.stable_items if items is None else partidas_estables(items, D, 'argumentos.items')
    return valor_estable(partidas, D)


def papel(datos):
    """§3.7 (R-68…R-73): paquetes dentro de compras y de contado, por fecha."""
    D = como_datos(datos)
    lb = libro(D, D.expenses, D.purchases)
    items = []
    for i, c in lb['vivas'].items():
        for pos, e in enumerate(c['embedded']):
            if e['cat'] == D.paper:
                items.append(((c['date'], 0, i, pos), dict(date=c['date'], person=e['person'], form='in_groceries',
                                                           purchase_seq=lb['seq'].get(i), amount_cents=e['amount'])))
    for g in D.expenses:
        if g['method'] != 'credit' and g['cat'] == D.paper:
            items.append(((g['date'], 1, g['id'], 0), dict(date=g['date'], person=g['paid_by'], form='cash',
                                                           purchase_seq=None, amount_cents=g['amount'])))
    items.sort(key=lambda x: x[0])
    return [dict(package=k + 1, date=x['date'], person=x['person'], person_name=D.nombre.get(x['person'], x['person']),
                 form=x['form'], purchase_seq=x['purchase_seq'], amount_cents=x['amount_cents']) for k, (_, x) in enumerate(items)]


def dias_del_mes(m):
    return ultimo(*m)
