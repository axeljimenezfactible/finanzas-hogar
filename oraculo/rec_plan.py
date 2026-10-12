"""Oráculo S4 · plan semanal (§5) y lista de compras neteada (§6, §6.5).

PROVISIONAL (bloqueos S4-reglas-lista y S4-asignacion-mercado): sobrantes previstos, reposición por ritmo, asignación de mercados,
orden y total de groups[]. Cada decisión que el contrato no fija está en el reporte (supuestos S-O*).
"""
import datetime as dt
from fractions import Fraction

from rec_calculo import cotiza_rige, detalle, existencia, receta, salud
from rec_comun import ErrorMotor, half_up, lunes_a_semana, semana_a_lunes, techo

SLOTS = ('breakfast', 'lunch', 'dinner', 'snack')
OBLIGATORIOS = ('breakfast', 'lunch', 'dinner')


# ---- §5.2 semana por omisión ----
def default_week(d, hoy):
    lunes = hoy - dt.timedelta(days=hoy.weekday())
    if hoy.weekday() == 6:  # domingo: abre la semana siguiente (K38)
        lunes = hoy + dt.timedelta(days=1)
    return {'week': lunes_a_semana(lunes), 'monday': lunes}


# ---- validación del plan (§5.1) ----
def valida_plan(d):
    for e in d.plan:
        p = d.person.get(e['person'])
        if p is None:
            raise ErrorMotor('invalid_argument', f"meal_plan[{e['idx']}].person")
        if not p['plan']:
            raise ErrorMotor('person_without_plan', e['person'])
        if e['out']:
            if e['recipe'] is not None:
                raise ErrorMotor('eating_out_conflict', f"meal_plan[{e['idx']}]")
            if e['serv'] is not None:
                raise ErrorMotor('invalid_argument', f"meal_plan[{e['idx']}].servings_milli")
        else:
            receta(d, e['recipe'])
            s = e['serv']
            if not isinstance(s, int) or isinstance(s, bool):
                raise ErrorMotor('invalid_argument', f"meal_plan[{e['idx']}].servings_milli")
            if s <= 0:
                raise ErrorMotor('invalid_servings', f"meal_plan[{e['idx']}].servings_milli")


def entradas_semana(d, lunes):
    """Entradas de la semana, ordenadas por (fecha, tiempo, persona según `persons`, orden de entrada)."""
    orden_p = {p['code']: i for i, p in enumerate(d.persons)}
    es = [e for e in d.plan if lunes <= e['date'] <= lunes + dt.timedelta(days=6)]
    return sorted(es, key=lambda e: (e['date'], SLOTS.index(e['slot']), orden_p[e['person']], e['idx']))


def unplanned(d, lunes):
    valida_plan(d)
    con_plan = [p['code'] for p in d.persons if p['plan']]
    llenas = {(e['date'], e['slot'], e['person']) for e in d.plan}
    n = 0
    for k in range(7):
        dia = lunes + dt.timedelta(days=k)
        for s in OBLIGATORIOS:
            if any((dia, s, p) not in llenas for p in con_plan):
                n += 1
    return n


def unplanned_slots(d, semana):
    return {'unplanned_slots': unplanned(d, semana_a_lunes(semana, 'week'))}


# ---- §6.4 sobrantes registrados ----
def cubre_sobrantes(d, entradas):
    """Cada sobrante (inventario con receta) cubre primero las entradas de su receta, en orden de fecha y tiempo, si no caducó.
    Devuelve ({idx de entrada: porciones cubiertas}, {código de sobrante: porciones restantes})."""
    resto = {it['code']: it['q'] for it in d.inv if it['recipe'] is not None}
    cubierto = {}
    for e in entradas:
        if e['out']:
            continue
        falta = e['serv']
        for it in d.inv:
            if it['recipe'] != e['recipe'] or falta == 0 or resto[it['code']] == 0:
                continue
            if it['expires'] is not None and it['expires'] < e['date']:
                continue
            t = min(resto[it['code']], falta)
            resto[it['code']] -= t
            falta -= t
            cubierto[e['idx']] = cubierto.get(e['idx'], 0) + t
    return cubierto, resto


# ---- §5.3 plan semanal ----
def meta_rige(d, persona):
    return salud(d, persona)['daily_kcal_target']


def week_plan(d, semana):
    lunes = semana_a_lunes(semana, 'week')
    valida_plan(d)
    es = entradas_semana(d, lunes)
    celdas = []
    for e in es:
        c = {'date': e['date'], 'slot': e['slot'], 'person': e['person'], 'recipe': e['recipe'], 'servings_milli': e['serv'],
             'eating_out': e['out'], 'cost_cents': None, 'kcal': None, 'price_store': None, 'price_observed_date': None}
        if not e['out']:
            x = detalle(d, e['recipe'], e['serv'])
            c['cost_cents'] = None if x['cost'] is None else half_up(x['cost'])
            c['kcal'] = None if x['kcal_milli'] is None else half_up(x['kcal_milli'] / 1000)
            if x['quotes'] and x['cost'] is not None:  # supuesto S-O3: con algún ingrediente sin precio, la celda no declara tienda ni fecha
                tiendas = {q['store'] for q in x['quotes']}
                c['price_store'] = next(iter(tiendas)) if len(tiendas) == 1 else None
                c['price_observed_date'] = min(q['observed'] for q in x['quotes'])
        celdas.append(c)
    dias, total_semana = [], 0
    for k in range(7):
        dia = lunes + dt.timedelta(days=k)
        del_dia = [c for c in celdas if c['date'] == dia and not c['eating_out']]
        costo = None if any(c['cost_cents'] is None for c in del_dia) else sum(c['cost_cents'] for c in del_dia)
        pers = []
        for p in d.persons:
            if not p['plan']:
                continue
            mias = [c for c in del_dia if c['person'] == p['code']]
            kcal = None if any(c['kcal'] is None for c in mias) else sum(c['kcal'] for c in mias)
            meta = meta_rige(d, p['code'])
            pers.append({'person': p['code'], 'kcal': kcal, 'daily_kcal_target': meta,
                         'difference_kcal': None if kcal is None or meta is None else kcal - meta})
        dias.append({'date': dia, 'cost_cents': costo, 'persons': pers})
        total_semana = None if costo is None or total_semana is None else total_semana + costo
    _, resto = cubre_sobrantes(d, es)
    sobran = [{'recipe': it['recipe'], 'servings_milli': resto[it['code']]} for it in d.inv
              if it['recipe'] is not None and resto[it['code']] > 0]
    return {'week': lunes_a_semana(lunes), 'cells': celdas, 'days': dias, 'week_cost_cents': total_semana,
            'unplanned_slots': unplanned(d, lunes), 'eating_out_slots': sum(1 for e in es if e['out']), 'leftovers': sobran}


# ---- §6 lista de compras ----
def _vertical(d, producto):
    g = d.products[producto]['group']
    if g is None:
        raise ErrorMotor('product_without_channel', producto)
    return d.groups[g]['channel']


def _recomendada(d, producto):
    return next((q for q in d.quotes.get(producto, []) if q['rec']), None)


def _tienda(d, producto, vertical, eleccion, cot):
    """§6.5: elegida; chain -> tienda de la cotización que rige; market -> mercado que surte el grupo con visita más cercana
    (empate: el primero en `stores`). PROVISIONAL (bloqueo S4-asignacion-mercado, opción A)."""
    if eleccion is not None and eleccion['store'] is not None:
        return eleccion['store']
    if vertical == 'chain':
        return None if cot is None else cot['store']
    grupo = d.products[producto]['group']
    cand = [(s['visit'], i) for i, s in enumerate(d.stores) if s['kind'] == 'market' and grupo in s['groups']]
    return d.stores[min(cand)[1]]['code'] if cand else None


def _articulo(d, p, kind, needed, reasons, recs, pres_ingrediente):
    """Un artículo de la lista (§6.2 pasos 5 y 6, §6.5, §6.6). `needed` = None en `manual`."""
    ch = d.choices.get((p, kind))
    cot = cotiza_rige(d, p, ch)
    fijada = None if ch is None else ch['q']
    comprar = fijada if fijada is not None else (needed or 0)  # cantidad a comprar (§6.6)
    if ch is not None and ch['pres'] is not None:
        pres = ch['pres']
    elif cot is not None:
        pres = cot['pres']
    else:
        pres = pres_ingrediente
    size = None
    if pres is not None:
        size = cot['size'] if cot is not None and cot['pres'] == pres else d.pres[pres]['size']
    paquetes = None if size is None else techo(Fraction(comprar, size))
    rec = _recomendada(d, p)
    eligio_tienda = ch is not None and ch['store'] is not None
    eligio_pres = ch is not None and ch['pres'] is not None
    a = {'product': p, 'kind': kind, 'needed_milli': needed, 'quantity_milli': fijada, 'presentation': pres,
         'packages_to_buy': paquetes,
         'leftover_milli': None if size is None else max(0, paquetes * size - (needed or 0)),
         'reasons': reasons, 'recipes': recs, 'store': None, 'is_bought': ch['bought'] if ch else False,
         'recommended_store': (rec['store'] if rec else None) if eligio_tienda else (cot['store'] if cot else None),
         'recommended_presentation': (rec['pres'] if rec else None) if eligio_pres else pres,
         'price_store': cot['store'] if cot else None, 'price_observed_date': cot['observed'] if cot else None,
         'is_stale_price': cot['stale'] if cot else None,
         'estimated_cents': None if cot is None or paquetes is None else paquetes * cot['price'], '_cot': cot}
    return a


def shopping_list(d, semana):
    lunes = semana_a_lunes(semana, 'week')
    domingo = lunes + dt.timedelta(days=6)
    valida_plan(d)
    es = entradas_semana(d, lunes)
    cubierto, _ = cubre_sobrantes(d, es)
    # N(p): necesidad de recetas con las porciones que los sobrantes no cubren (§6.2 paso 1, §6.4)
    need, recetas, pres_ing = {}, {}, {}
    for e in es:
        if e['out']:
            continue
        r = receta(d, e['recipe'])
        serv = e['serv'] - cubierto.get(e['idx'], 0)
        for i in r['ings']:
            if i['optional']:
                continue
            p = i['product']
            pres_ing.setdefault(p, i['presentation'])
            if serv > 0:
                need[p] = need.get(p, Fraction(0)) + Fraction(i['q'] * serv, r['base'] * 1000)
                recetas.setdefault(p, set()).add(e['recipe'])
    # E(p) al lunes y R(p) por ritmo (opción A de S4-reglas-lista, pregunta 2: PROVISIONAL)
    exist, repo = {}, {}
    for it in d.inv:
        p = it['product']
        if p is None:
            continue
        exist[p] = exist.get(p, 0) + existencia(it, lunes)
        if it['mode'] == 'consumption_rate' and it['level'] is not None and existencia(it, domingo) <= it['level']:
            repo[p] = repo.get(p, 0) + half_up(Fraction(it['rq'] * 7, it['rp']))
    items = []
    for p in sorted(set(need) | set(repo) | {p for (p, k) in d.choices if k == 'calculated'}):
        neto = need.get(p, Fraction(0)) + repo.get(p, 0) - exist.get(p, 0)
        necesidad = techo(neto) if neto > 0 else 0
        if necesidad == 0 and (p, 'calculated') not in d.choices:
            continue  # §6.2 paso 7
        items.append(_articulo(d, p, 'calculated', necesidad,
                               [r for r, ok in (('recipe', p in recetas), ('replenishment', p in repo)) if ok],
                               sorted(recetas.get(p, [])), pres_ing.get(p)))
    for p in sorted(p for (p, k) in d.choices if k == 'manual'):
        items.append(_articulo(d, p, 'manual', None, ['manual'], [], None))
    # tienda y grupos (§6.5)
    pares = {}
    for a in items:
        vert = _vertical(d, a['product'])
        a['store'] = _tienda(d, a['product'], vert, d.choices.get((a['product'], a['kind'])), a['_cot'])
        pares.setdefault((vert, a['store']), []).append(a)
    orden_t = {s['code']: i for i, s in enumerate(d.stores)}

    def visita(t):
        return d.store[t]['visit'] if t is not None and t in d.store and d.store[t]['kind'] == 'market' else None

    def clave(par):
        v, t = par
        vi = visita(t)
        return (0 if v == 'market' else 1, vi is None, vi or dt.date.min, t is None, orden_t.get(t, len(orden_t)), t or '')
    grupos = []
    for par in sorted(pares, key=clave):
        v, t = par
        total = sum(a['estimated_cents'] for a in pares[par] if a['estimated_cents'] is not None)
        grupos.append({'vertical': v, 'store': t, 'visit_date': visita(t), 'estimated_total_cents': total})
    sin_precio = sorted({a['product'] for a in items if a['_cot'] is None})
    for a in items:
        del a['_cot']
    return {'week': lunes_a_semana(lunes), 'items': items, 'groups': grupos,
            'estimated_total_cents': sum(g['estimated_total_cents'] for g in grupos), 'missing_prices': sin_precio}
