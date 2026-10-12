"""Oráculo S5 · cálculo de las funciones consultables (S5-interfaz.md §4 a §10, §13, §14.1, §15.2).

Todo exacto con Fraction (precio por unidad base por multiplicación cruzada); un solo redondeo HALF_UP al publicar.
Los puntos que el contrato deja abiertos están en docs/estado/wf1/oraculo-S5.md (Supuestos).
"""
from fractions import Fraction

from pre_comun import ErrorMotor, frac, half_up, techo
from pre_datos import COMPITEN

PER = {'g': 'kg', 'ml': 'l', 'piece': 'piece'}
MULT = {'g': 1000000, 'ml': 1000000, 'piece': 1000}


# ---------- §4 unidad base ----------
def unidad_de(d, pres):
    return d.products[d.pres[pres]['product']]['unit']


def exacto_base(d, pres, cents):
    """Centavos por 1 kg, 1 l o 1 pieza, sin redondear (Fraction)."""
    return Fraction(cents * MULT[unidad_de(d, pres)], d.pres[pres]['size'])


def base_unit_price(d, pres, cents):
    if pres not in d.pres:
        raise ErrorMotor('unknown_presentation', pres)
    ex = exacto_base(d, pres, cents)
    return {'base_unit_price_cents': half_up(ex), 'per': PER[unidad_de(d, pres)], 'exact': frac(ex)}


# ---------- §5 frescura ----------
def fresco(d, obs, hoy, fuente):
    edad = (hoy - obs).days
    if edad < 0:
        raise ErrorMotor('future_date', 'observed_date')
    if edad < d.dv:
        w = Fraction(1)
    elif edad <= d.d0:
        w = Fraction(d.d0 - edad, d.d0 - d.dv + 1)
    else:
        w = Fraction(0)
    return {'age': edad, 'stale': edad >= d.dv, 'w': w, 'competes': fuente in COMPITEN and w > 0}


def freshness(d, obs, fuente):
    f = fresco(d, obs, d.today, fuente)
    w = f['w']
    disp = half_up(w * 10000)
    return {'age_days': f['age'], 'is_stale': f['stale'], 'weight': frac(w), 'weight_display': f'{disp // 10000}.{disp % 10000:04d}',
            'competes': f['competes']}


# ---------- §6 precio vigente ----------
def _de(d, pres, store, hoy):
    return [p for p in d.prices if p['pres'] == pres and p['store'] == store and p['date'] <= hoy]


def vigente(d, pres, store, hoy):
    c = [p for p in _de(d, pres, store, hoy) if not p['promo'] and p['source'] in COMPITEN]
    return max(c, key=lambda p: (p['date'], p['created'], p['n'])) if c else None


def current_price(d, pres, store):
    if pres not in d.pres:
        raise ErrorMotor('unknown_presentation', pres)
    if store not in d.stores:
        raise ErrorMotor('unknown_store', store)
    v = vigente(d, pres, store, d.today)
    refs = sorted([p for p in _de(d, pres, store, d.today) if not p['promo'] and p['source'] not in COMPITEN],
                  key=lambda p: (p['date'], p['created'], p['n']), reverse=True)
    pro = [p for p in _de(d, pres, store, d.today) if p['promo'] and p['source'] in COMPITEN and p['date'] == d.today]
    pr = max(pro, key=lambda p: (p['created'], p['n'])) if pro else None
    f = lambda p: {'price_cents': p['cents'], 'observed_date': p['date'].isoformat()}
    return {'price': None if v is None else {**f(v), 'source': v['source']},
            'references': [{'source': p['source'], **f(p)} for p in refs],
            'promotion': None if pr is None else f(pr)}


# ---------- §7.4, §14.1 cotización ----------
def compite_tienda(d, st, product):
    s = d.stores[st]
    if s['kind'] != 'market':
        return True
    g = d.products[product]['group']
    return g is not None and g in s['groups']


def candidatos(d, product, hoy, store=None, pres=None):
    tiendas = [store] if store is not None else list(d.stores)
    out = []
    for st in tiendas:
        if not compite_tienda(d, st, product):
            continue
        for pc in d.products[product]['pres']:
            if pres is not None and pc != pres:
                continue
            p = vigente(d, pc, st, hoy)
            if p is None:
                continue
            f = fresco(d, p['date'], hoy, p['source'])
            if not f['competes']:
                continue
            size = d.pres[pc]['size']
            clave = (Fraction(p['cents'], size) / f['w'], p['cents'], d.stores[st]['idx'], d.pres[pc]['idx'])
            out.append((clave, st, pc, p, f))
    return out


def cotiza(d, product, hoy, store=None, pres=None):
    c = candidatos(d, product, hoy, store, pres)
    if not c:
        return None
    clave, st, pc, p, f = min(c, key=lambda x: x[0])
    return {'product': product, 'store': st, 'presentation': pc, 'price_cents': p['cents'], 'size_milli': d.pres[pc]['size'],
            'observed_date': p['date'].isoformat(), 'is_stale': f['stale'], 'recommended': store is None and pres is None,
            '_age': f['age'], '_w': f['w']}


def quote_for(d, product, store, pres):
    if product not in d.products:
        raise ErrorMotor('invalid_argument', 'product')
    if store is not None and store not in d.stores:
        raise ErrorMotor('unknown_store', store)
    if pres is not None and (pres not in d.pres or d.pres[pres]['product'] != product):
        raise ErrorMotor('unknown_presentation', pres)
    q = cotiza(d, product, d.today, store, pres)
    return {'quote': None if q is None else {k: v for k, v in q.items() if not k.startswith('_')}}


# ---------- §7.2 compra contra carrito ----------
def _agrupa(cart):
    g = {}
    for ln in cart['lines']:
        e = g.setdefault(ln['pres'], {'pres': ln['pres'], 'qty': 0, 'total': 0, 'unit': ln['unit'], 'exacto': 0})
        e['qty'] += ln['qty']
        e['exacto'] += ln['qty'] * ln['unit']  # importe exacto en milésimas de centavo, sin redondeo
        e['total'] += half_up(Fraction(ln['qty'] * ln['unit'], 1000))
    return list(g.values())


def _piezas(qty):
    return qty // 1000 if qty % 1000 == 0 else 1


def _unit(e):
    return Fraction(e['total'] * 1000, e['qty'])


def _linea(d, st, pres, qty, unit):
    tot = half_up(Fraction(qty * unit.numerator, 1000 * unit.denominator))
    return {'store': st, 'presentation': pres, 'quantity_milli': qty, 'total_cents': tot,
            'base_unit_price_cents': half_up(exacto_base(d, pres, unit))}


def _ganador(vals):
    """vals: lista de (store, Fraction). Gana el menor valor; null si el menor se repite."""
    m = min(v for _, v in vals)
    g = [s for s, v in vals if v == m]
    return g[0] if len(g) == 1 else None


def compare_carts(d, base, others):
    codes = [base] + list(others)
    for c in codes:
        if c not in d.carts:
            raise ErrorMotor('invalid_argument', f'cart {c}')
    if len(set(codes)) != len(codes):
        raise ErrorMotor('invalid_argument', 'carts repetidos')
    gr = {c: _agrupa(d.carts[c]) for c in codes}
    por = {c: {e['pres']: e for e in gr[c]} for c in codes}
    tienda = {c: d.carts[c]['store'] for c in codes}
    rows_i, rows_s, rows_o = [], [], []
    rest = {c: {e['pres']: e['qty'] for e in gr[c]} for c in codes}  # cantidad no emparejada por presentación
    orden_pres = []
    for c in codes:
        for e in gr[c]:
            if e['pres'] not in orden_pres:
                orden_pres.append(e['pres'])
    # idénticos: la presentación está en todas las canastas; cantidad emparejada = la menor
    for pc in orden_pres:
        if all(pc in por[c] for c in codes):
            q = min(por[c][pc]['qty'] for c in codes)
            lines, vals = [], []
            for c in codes:
                u = _unit(por[c][pc])
                lines.append(_linea(d, tienda[c], pc, q, u))
                vals.append((tienda[c], u))
                rest[c][pc] -= q
            rows_i.append({'kind': 'identical', 'product': d.pres[pc]['product'], 'lines': lines, 'winner_store': _ganador(vals),
                           '_q': q, '_pc': pc})
    # sustitutos: el producto aparece con presentaciones sin emparejar en dos o más canastas
    productos = []
    for pc in orden_pres:
        pr = d.pres[pc]['product']
        if pr not in productos:
            productos.append(pr)
    for pr in productos:
        en = {c: [pc for pc in por[c] if d.pres[pc]['product'] == pr and rest[c][pc] > 0 and not all(pc in por[x] for x in codes)] for c in codes}
        con = [c for c in codes if en[c]]
        if len(con) >= 2:
            lines, vals = [], []
            # ganador (PROVISIONAL, bloqueo S5-carritos-detalle, opción A): precio por unidad base de lo que esa tienda
            # cobró por el producto = importe exacto total ÷ contenido total de sus presentaciones sin emparejar
            for c in con:
                costo = contenido = 0
                for pc in en[c]:
                    e = por[c][pc]
                    lines.append(_linea(d, tienda[c], pc, rest[c][pc], _unit(e)))
                    assert rest[c][pc] == e['qty']  # una presentación sin emparejar nunca se consumió en parte
                    costo += e['exacto']
                    contenido += e['qty'] * d.pres[pc]['size']
                    rest[c][pc] = 0
                vals.append((tienda[c], Fraction(costo, contenido)))
            rows_s.append({'kind': 'substitute', 'product': pr, 'lines': lines, 'winner_store': _ganador(vals), '_pcs': [x['presentation'] for x in lines]})
    # solo en una tienda: lo que queda
    for c in codes:
        for e in gr[c]:
            q = rest[c][e['pres']]
            if q > 0:
                rows_o.append({'kind': 'only_one', 'product': d.pres[e['pres']]['product'],
                               'lines': [_linea(d, tienda[c], e['pres'], q, _unit(e))], 'winner_store': None, '_pc': e['pres'], '_c': c})
    filas = rows_i + rows_s + rows_o
    ult = [{'store': tienda[c], 'total_cents': 0, 'identical_count': 0, 'identical_total_cents': 0, 'substitute_count': 0, 'piece_count': 0,
            'missing_count': 0, 'stale_count': 0, 'difference_vs_cheapest_cents': 0, 'assumes_membership': False, '_i': i}
           for i, c in enumerate(codes)]
    idx = {tienda[c]: i for i, c in enumerate(codes)}
    for r in filas:
        for ln in r['lines']:
            s = ult[idx[ln['store']]]
            s['total_cents'] += ln['total_cents']
            s['piece_count'] += _piezas(ln['quantity_milli'])
            if r['kind'] == 'identical':
                s['identical_total_cents'] += ln['total_cents']
    for r in rows_i:
        for ln in r['lines']:
            ult[idx[ln['store']]]['identical_count'] += 1
    for r in rows_s:
        for st in {ln['store'] for ln in r['lines']}:
            ult[idx[st]]['substitute_count'] += 1
    for i, c in enumerate(codes):
        s = ult[i]
        edad = (d.today - d.carts[c]['date']).days
        s['stale_count'] = len(gr[c]) if edad >= d.dv else 0  # renglones agrupados (bloqueo S5-carritos-detalle, opción A)
        # faltantes (bloqueo S5-carritos-detalle, opción A): productos de la unión de las canastas sin ninguna presentación aquí
        s['missing_count'] = sum(1 for pr in productos if not any(d.pres[e['pres']]['product'] == pr for e in gr[c]))
        st = d.stores[tienda[c]]
        s['assumes_membership'] = st['requires'] and not st['has']
    minimo = min(s['total_cents'] for s in ult)
    for s in ult:
        s['difference_vs_cheapest_cents'] = s['total_cents'] - minimo
    ult.sort(key=lambda s: (s['total_cents'], s['_i']))
    # winners[] en el orden de stores[] (convención 1 de oraculo-S5.md), también con count 0
    ganan = [{'store': s['store'], 'count': sum(1 for r in rows_i if r['winner_store'] == s['store'])} for s in ult]
    limpia = lambda r: {k: v for k, v in r.items() if not k.startswith('_')}
    return {'stores': [limpia(s) for s in ult], 'identical_pieces': sum(_piezas(r['_q']) for r in rows_i), 'winners': ganan,
            'ties': sum(1 for r in rows_i if r['winner_store'] is None), 'rows': [limpia(r) for r in filas]}


# ---------- §7.3 canasta por necesidad ----------
def costo_en_tienda(d, st, product, qty, hoy):
    """Presentación de `st` que cubre `qty` con el menor costo de paquetes completos (empate: menor precio por unidad base, luego menor id)."""
    if not compite_tienda(d, st, product):
        return None
    mejor = None
    for pc in d.products[product]['pres']:
        p = vigente(d, pc, st, hoy)
        if p is None:
            continue
        f = fresco(d, p['date'], hoy, p['source'])
        if not f['competes']:
            continue
        size = d.pres[pc]['size']
        paq = techo(Fraction(qty, size))
        clave = (paq * p['cents'], Fraction(p['cents'], size), d.pres[pc]['idx'])
        if mejor is None or clave < mejor[0]:
            mejor = (clave, pc, paq, p, f)
    if mejor is None:
        return None
    clave, pc, paq, p, f = mejor
    return {'pres': pc, 'packages': paq, 'cost': clave[0], 'price': p, 'fresh': f}


def compare_basket(d):
    filas = []
    for it in d.shopping:
        q = cotiza(d, it['product'], d.today)
        filas.append((it, q))
    est = {st: {'store': st, 'total_cents': 0, 'identical_count': 0, 'identical_total_cents': 0, 'substitute_count': 0, 'piece_count': 0,
                'missing_count': 0, 'stale_count': 0, 'difference_vs_cheapest_cents': 0,
                'assumes_membership': d.stores[st]['requires'] and not d.stores[st]['has'], '_rec': 0, '_i': d.stores[st]['idx']}
           for st in d.stores}
    items = []
    for it, q in filas:
        prices = []
        for st in d.stores:
            c = costo_en_tienda(d, st, it['product'], it['qty'], d.today)
            s = est[st]
            if c is None:
                s['missing_count'] += 1
                continue
            sub = it['rec'] is not None and c['pres'] != it['rec']
            s['total_cents'] += c['cost']
            s['piece_count'] += c['packages']
            if sub:
                s['substitute_count'] += 1
            else:
                s['identical_count'] += 1
                s['identical_total_cents'] += c['cost']
            if c['fresh']['stale']:
                s['stale_count'] += 1
            prices.append({'store': st, 'presentation': c['pres'], 'price_cents': c['price']['cents'],
                           'base_unit_price_cents': half_up(exacto_base(d, c['pres'], c['price']['cents'])),
                           'observed_date': c['price']['date'].isoformat(), 'is_substitute': sub})
        if q is not None:
            est[q['store']]['_rec'] += 1
            estimado = techo(Fraction(it['qty'], q['size_milli'])) * q['price_cents']
        items.append({'product': it['product'], 'quantity_milli': it['qty'],
                      'recommended_store': None if q is None else q['store'],
                      'recommended_presentation': None if q is None else q['presentation'],
                      'estimated_price_cents': None if q is None else estimado,
                      'is_stale_price': False if q is None else q['is_stale'], 'prices': prices})
    lista_est = sorted(est.values(), key=lambda s: (s['total_cents'], -s['_rec'], s['stale_count'], s['_i']))
    completas = [s for s in lista_est if s['missing_count'] == 0]
    barata = completas[0] if completas else None
    ref = barata['total_cents'] if barata else (min(s['total_cents'] for s in lista_est) if lista_est else 0)
    for s in lista_est:
        s['difference_vs_cheapest_cents'] = s['total_cents'] - ref
    return {'stores': [{k: v for k, v in s.items() if not k.startswith('_')} for s in lista_est], 'items': items,
            'cheapest_store': None if barata is None else barata['store']}


# ---------- §8 membresía ----------
def membership(d, store, ahorro_dado):
    if store not in d.stores:
        raise ErrorMotor('unknown_store', store)
    s = d.stores[store]
    if ahorro_dado is not None:
        ahorro = ahorro_dado
    else:
        ahorro = 0
        for it in d.shopping:
            cs = costo_en_tienda(d, store, it['product'], it['qty'], d.today)
            otros = [costo_en_tienda(d, o, it['product'], it['qty'], d.today) for o, so in d.stores.items() if o != store and not so['requires']]
            otros = [c['cost'] for c in otros if c is not None]
            if cs is not None and otros:
                ahorro += max(0, min(otros) - cs['cost'])
    compras = sum(1 for g in d.grocery if d.stores[g['store']]['kind'] != 'market' and 0 <= (d.today - g['date']).days < 365)
    cuota = s['fee']
    equilibrio = None if cuota is None or ahorro == 0 else techo(Fraction(cuota, ahorro))
    anual = ahorro * compras
    if cuota is None:
        veredicto = 'missing_fee'
    elif ahorro == 0:
        veredicto = 'not_worth_it'
    else:
        veredicto = 'worth_it' if anual >= cuota else 'not_worth_it'
    return {'store': store, 'assumes_membership': s['requires'] and not s['has'], 'basket_savings_cents': ahorro,
            'annual_fee_cents': cuota, 'break_even_purchases': equilibrio, 'purchases_per_year': compras,
            'annual_savings_cents': anual, 'verdict': veredicto}


# ---------- §9 mercados ----------
def next_visit(d, store, desde):
    if store not in d.stores or d.stores[store]['visit_weekday'] is None:
        raise ErrorMotor('unknown_store' if store not in d.stores else 'invalid_argument', store)
    wd = d.stores[store]['visit_weekday']
    import datetime as dt
    n = (wd - desde.isoweekday()) % 7
    return {'date': (desde + dt.timedelta(days=n)).isoformat()}


# ---------- §10 precios viejos ----------
def stale_prices_for_plan(d, products, asignadas, fecha_despensa):
    stale, sin = [], []
    for pr in products:
        if pr not in d.products:
            raise ErrorMotor('invalid_argument', pr)
        st = asignadas.get(pr)
        if st is not None and st not in d.stores:
            raise ErrorMotor('unknown_store', st)
        q = cotiza(d, pr, fecha_despensa, st)
        if q is None:
            sin.append(pr)
        elif q['is_stale']:
            stale.append({'product': pr, 'store': q['store'], 'presentation': q['presentation'], 'observed_date': q['observed_date'],
                          'age_days': q['_age'], 'weight': frac(q['_w'])})
    return {'stale': stale, 'without_price': sin}


# ---------- §13 lectura asistida ----------
def assisted_can_enable(d, store):
    if store not in d.stores:
        raise ErrorMotor('unknown_store', store)
    ok = d.assisted_on and d.stores[store]['assisted']
    t = [x for x in d.terms if x['site'] == store]
    ok = ok and any(all(isinstance(x[k], str) and x[k] != '' for k in ('verdict', 'url', 'read_date')) for x in t)
    if not ok:
        raise ErrorMotor('no_terms_verdict', store)
    return {'can_enable': True}
