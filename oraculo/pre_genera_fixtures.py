#!/usr/bin/env python3
"""Genera las fixtures de oraculo/fixtures/s5/ (S5-interfaz.md §15) y sus archivos de _datos/.

Los `esperado` salen de cálculos a mano escritos aquí (literales y aritmética de enteros/Fraction con la fórmula del contrato)
o de data.json leído directo; NO se llama al oráculo (pre_calculo.py), para que el esperado sea una segunda derivación.
Uso: python3 -I oraculo/pre_genera_fixtures.py   (desde cualquier carpeta; escribe en oraculo/fixtures/s5/)
"""
import datetime as dt
import json
from decimal import Decimal
from fractions import Fraction
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[1]
OUT = RAIZ / 'oraculo' / 'fixtures' / 's5'
DATOS = OUT / '_datos'
CAT = 'datos/semilla/productos.json'
BLOQ = 'docs/estado/bloqueos/'


def hu(x):
    x = Fraction(x)
    n, d = x.numerator, x.denominator
    q = (2 * abs(n) + d) // (2 * d)
    return q if n >= 0 else -q


def cts(v):
    c = Decimal(v) * 100
    assert c == c.to_integral_value()
    return int(c)


def fr(x):
    x = Fraction(x)
    return {'num': x.numerator, 'den': x.denominator}


def escribe(nombre, fx):
    OUT.mkdir(parents=True, exist_ok=True)
    (OUT / f'{nombre}.json').write_text(json.dumps(fx, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')


def fixture(caso, desc, fuente, datos, consultas, esperado, prov=None):
    escribe(caso, {'formato': 'fixture/1', 'motor': 's5', 'caso': caso, 'descripcion': desc,
                   'fuente': [{'tipo': t, 'ruta': r} for t, r in fuente], 'provisional': prov,
                   'entradas': {'datos': datos, 'consultas': consultas}, 'esperado': esperado})


def prov(bloqueo, que):
    return {'bloqueo': BLOQ + bloqueo + '.md', 'que': que}


def tienda(code, kind='chain', name=None, **kw):
    t = {'code': code, 'kind': kind, 'parent': kw.get('parent'), 'name': name or code, 'profeco_name': kw.get('profeco_name'),
         'municipality': kw.get('municipality'), 'visit_weekday': kw.get('visit_weekday'),
         'requires_membership': kw.get('requires', False), 'has_membership': kw.get('has', False),
         'membership_fee_cents': kw.get('fee'), 'assisted_reading': kw.get('assisted', False), 'groups': kw.get('groups', [])}
    return t


def precio(pres, store, cents, fecha, source='manual', promo=False, created='2026-10-10T12:00:00Z', ref=None):
    return {'presentation': pres, 'store': store, 'price_cents': cents, 'observed_date': fecha, 'source': source,
            'is_promotional': promo, 'source_ref': ref, 'created_at': created}


def catalogo(prods):
    """prods: lista de (producto, unidad, [(presentación, milésimas)]). Un grupo `abarrotes` (cadena) y otro `verduras` (mercado)."""
    return {'groups': [{'code': 'abarrotes', 'name': 'Abarrotes', 'purchase_channel': 'chain'},
                       {'code': 'verduras', 'name': 'Frutas y verduras', 'purchase_channel': 'market'}],
            'products': [{'code': p, 'name': p.capitalize(), 'base_unit': u, 'group': g, 'grams_per_unit_milli': None} for p, u, g, _ in prods],
            'presentations': [{'code': c, 'product': p, 'name': c, 'size_milli': s} for p, _, _, ps in prods for c, s in ps]}


def lista_item(prod, nec=None, cant=None, rec=None):
    return {'product': prod, 'needed_milli': nec, 'quantity_milli': cant, 'recommended_presentation': rec}


def consulta(id_, funcion, **args):
    return {'id': id_, 'funcion': funcion, 'argumentos': args}


# ---------------------------------------------------------------- data.json
REP = json.load(open(RAIZ / 'datos' / 'semilla' / 'data.json', encoding='utf-8'), parse_float=Decimal)
PRODS = json.load(open(RAIZ / 'datos' / 'semilla' / 'productos.json', encoding='utf-8'))
NOMBRE_PROD = {x['id']: p['id'] for p in PRODS['productos'] for x in p['presentaciones']}
NOMBRE = {x['nombre']: (p['id'], x['id'], x['contenido_milli'], p['unidad_base']) for p in PRODS['productos'] for x in p['presentaciones']}
HOY = '2026-10-12'
TIENDAS_CANON = [tienda('bodega', name='Bodega Aurrera (sucursal Vallejo)'), tienda('walmart', name='Walmart (sucursal Acueducto de Guadalupe)')]


def canon():
    t = REP['totales']
    idn = REP['items_identicos']
    sus = REP['sustitutos']
    tb, tw = cts(t['bodega']), cts(t['walmart'])
    ib, iw = cts(t['bodega_identicos']), cts(t['walmart_identicos'])
    assert sum(cts(i['bodega_total']) for i in idn) == ib and sum(cts(i['walmart_total']) for i in idn) == iw
    assert ib + sum(cts(s['bodega_total']) for s in sus) + cts(REP['atun_sexta_pieza_bodega']) == tb
    assert iw + sum(cts(s['walmart_total']) for s in sus) == tw
    piezas_b = t['articulos_bodega']
    assert sum(i['cantidad'] for i in idn) == t['articulos_identicos']
    gana_b = sum(1 for i in idn if i['bodega_unit'] < i['walmart_unit'])
    gana_w = sum(1 for i in idn if i['bodega_unit'] > i['walmart_unit'])
    empates = sum(1 for i in idn if i['bodega_unit'] == i['walmart_unit'])
    assert (gana_b, empates, gana_w) == (REP['resultados_por_producto']['bodega_gana'], REP['resultados_por_producto']['empates'],
                                         REP['resultados_por_producto']['walmart_gana'])
    filas = []
    for i in idn:
        g = 'bodega' if i['bodega_unit'] < i['walmart_unit'] else ('walmart' if i['bodega_unit'] > i['walmart_unit'] else None)
        filas.append({'kind': 'identical', 'product': NOMBRE[i['producto']][0], 'winner_store': g})
    for s in sus:
        ub = Fraction(cts(s['bodega_total']) * 1000000, NOMBRE[s['bodega']][2] * s['cant_b'])
        uw = Fraction(cts(s['walmart_total']) * 1000000, NOMBRE[s['walmart']][2] * s['cant_w'])
        filas.append({'kind': 'substitute', 'product': NOMBRE[s['bodega']][0], 'winner_store': 'bodega' if ub < uw else 'walmart'})
    filas.append({'kind': 'only_one', 'product': NOMBRE['Atún Dolores aleta amarilla en agua con aceite 140 g (5 pzas)'][0], 'winner_store': None})
    esp_tot = {'stores': [
        {'store': 'bodega', 'total_cents': tb, 'identical_count': len(idn), 'identical_total_cents': ib, 'substitute_count': len(sus),
         'piece_count': piezas_b, 'missing_count': 0, 'stale_count': 0, 'difference_vs_cheapest_cents': 0},
        {'store': 'walmart', 'total_cents': tw, 'identical_count': len(idn), 'identical_total_cents': iw, 'substitute_count': len(sus),
         'piece_count': t['articulos_walmart'], 'missing_count': 0, 'stale_count': 0, 'difference_vs_cheapest_cents': tw - tb}],
        'identical_pieces': t['articulos_identicos'],
        'winners': [{'store': 'bodega', 'count': gana_b}, {'store': 'walmart', 'count': gana_w}], 'ties': empates}
    assert tw - tb == 13850 and iw - ib == 10050
    datos = {'today': HOY, 'catalog_file': CAT, 'stores': TIENDAS_CANON, 'report_file': 'datos/semilla/data.json'}
    fixture('s5-comparativa-canon',
            'Compra registrada de Bodega Aurrera contra el carrito de Walmart del 10-oct-2026 (S5-interfaz §3, §7.2). Todo sale de data.json '
            '(report_file): K-01 3,094.50; K-02 3,233.00; K-03 +138.50 = +100.50 (idénticos 2,874.00 - 2,773.50) + 59.00 (sustitutos 359.00 - 300.00) '
            '- 21.00 (la lata que solo compró Bodega); K-04/K-05 2,773.50 y 2,874.00; K-06 +100.50; K-07 61/60/55 piezas; K-08 33/3 renglones; '
            'ganadores 22/8/3. Los renglones se piden aparte (canon-filas) porque su orden es una convención del oráculo.',
            [('data.json', 'totales'), ('data.json', 'resultados_por_producto'), ('data.json', 'items_identicos'), ('data.json', 'sustitutos')],
            datos, [consulta('canon-totales', 'compare_carts', base='bodega', others=['walmart']),
                    consulta('canon-filas', 'compare_carts', base='bodega', others=['walmart'])],
            {'canon-totales': esp_tot, 'canon-filas': {'rows': filas}})
    return tb, tw, ib, iw, idn, sus


def unidad():
    sus = REP['sustitutos']
    consultas, esperado = [], {}
    for s in sus:
        for lado, cant, tot in (('bodega', s['cant_b'], s['bodega_total']), ('walmart', s['cant_w'], s['walmart_total'])):
            nom = s[lado]
            prod, pres, size, uni = NOMBRE[nom]
            unit = cts(tot) // cant
            assert unit * cant == cts(tot)
            ex = Fraction(unit * (1000 if uni == 'pieza' else 1000000), size)
            id_ = f'{prod}-{lado}'
            consultas.append(consulta(id_, 'base_unit_price', presentation=pres, price_cents=unit))
            esperado[id_] = {'base_unit_price_cents': hu(ex), 'per': {'g': 'kg', 'ml': 'l', 'pieza': 'piece'}[uni], 'exact': fr(ex)}
    # regresión del contrato §4: 2933/2800, 12692/16406, 3375/2833
    assert [esperado[k]['base_unit_price_cents'] for k in esperado] == [2933, 2800, 12692, 16406, 3375, 2833]
    # redondeo HALF_UP en la mitad exacta (g y pieza) y en pieza
    extra = [('redondeo-g', 'jamon-pavo.sabori-400g', 6199, 1000000, 'kg'), ('redondeo-pieza', 'caldo-pollo-cubos.knorr-24-cubos', 12900, 1000, 'piece'),
             ('exacto-g', 'queso-manchego.nochebuena-200g', 6100, 1000000, 'kg')]
    for id_, pres, cents, mult, per in extra:
        size = {'jamon-pavo.sabori-400g': 400000, 'caldo-pollo-cubos.knorr-24-cubos': 24000, 'queso-manchego.nochebuena-200g': 200000}[pres]
        ex = Fraction(cents * mult, size)
        consultas.append(consulta(id_, 'base_unit_price', presentation=pres, price_cents=cents))
        esperado[id_] = {'base_unit_price_cents': hu(ex), 'per': per, 'exact': fr(ex)}
    assert esperado['redondeo-g']['base_unit_price_cents'] == 15498 and esperado['redondeo-pieza']['base_unit_price_cents'] == 538
    # el más barato por renglón: carritos solo con los 3 sustitutos (suman 300.00 y 359.00)
    lin = lambda lado, c: [{'presentation': NOMBRE[s[lado]][1], 'quantity_milli': s[c] * 1000, 'unit_price_cents': cts(s[lado + '_total']) // s[c]}
                           for s in sus]
    carts = [{'code': 'bodega', 'store': 'bodega', 'kind': 'expense', 'observed_date': '2026-10-10', 'source': 'manual', 'lines': lin('bodega', 'cant_b')},
             {'code': 'walmart', 'store': 'walmart', 'kind': 'cart', 'observed_date': '2026-10-10', 'source': 'manual', 'lines': lin('walmart', 'cant_w')}]
    consultas.append(consulta('sustitutos-carritos', 'compare_carts', base='bodega', others=['walmart']))
    esperado['sustitutos-carritos'] = {
        'stores': [{'store': 'bodega', 'total_cents': 30000, 'identical_count': 0, 'substitute_count': 3, 'piece_count': 5, 'difference_vs_cheapest_cents': 0},
                   {'store': 'walmart', 'total_cents': 35900, 'identical_count': 0, 'substitute_count': 3, 'piece_count': 5, 'difference_vs_cheapest_cents': 5900}],
        'ties': 0, 'rows': [{'kind': 'substitute', 'product': 'arroz-grano-largo', 'winner_store': 'walmart'},
                            {'kind': 'substitute', 'product': 'sustituto-crema-cafe', 'winner_store': 'bodega'},
                            {'kind': 'substitute', 'product': 'leche', 'winner_store': 'walmart'}]}
    fixture('s5-unidad-sustitutos',
            'Precio por unidad base de los 3 sustitutos de data.json (arroz 2933/2800, crema 12692/16406, leche 3375/2833, sin redondeo intermedio), '
            'redondeo HALF_UP en mitades exactas (6199 por 400 g = 15497.5 -> 15498; 12900 por 24 cubos = 537.5 por pieza -> 538) y el más barato por renglón '
            'con carritos de solo los 3 sustitutos (Bodega 300.00, Walmart 359.00).',
            [('data.json', 'sustitutos'), ('regla', 'RN-20')],
            {'today': HOY, 'catalog_file': CAT, 'stores': TIENDAS_CANON, 'carts': carts}, consultas, esperado)


def frescura():
    hoy = dt.date(2026, 12, 31)
    f = lambda n: (hoy - dt.timedelta(days=n)).isoformat()

    def peso(n, dv=30, d0=90):
        if n < dv:
            return Fraction(1)
        return Fraction(d0 - n, d0 - dv + 1) if n <= d0 else Fraction(0)

    def disp(w):
        x = hu(w * 10000)
        return f'{x // 10000}.{x % 10000:04d}'

    def esp(n, dv=30, d0=90, src_comp=True):
        w = peso(n, dv, d0)
        return {'age_days': n, 'is_stale': n >= dv, 'weight': fr(w), 'weight_display': disp(w), 'competes': src_comp and w > 0}
    consultas, esperado = [], {}
    for n in [0] + list(range(29, 92)):
        id_ = f'edad-{n:03d}'
        consultas.append(consulta(id_, 'freshness', observed_date=f(n), source='manual'))
        esperado[id_] = esp(n)
    assert esperado['edad-030']['weight_display'] == '0.9836' and esperado['edad-060']['weight_display'] == '0.4918'
    assert esperado['edad-090']['weight'] == {'num': 0, 'den': 1} and esperado['edad-091']['competes'] is False
    for src in ('manual', 'ticket_manual', 'ticket_ocr', 'site_assisted', 'profeco'):
        id_ = f'fuente-{src.replace("_", "-")}'
        consultas.append(consulta(id_, 'freshness', observed_date=f(5), source=src))
        esperado[id_] = esp(5, src_comp=src in ('manual', 'ticket_manual', 'ticket_ocr'))
    consultas.append(consulta('futura', 'freshness', observed_date=(hoy + dt.timedelta(days=1)).isoformat(), source='manual'))
    esperado['futura'] = {'error': 'future_date'}
    cat = catalogo([('arroz', 'g', 'abarrotes', [('arroz-1kg', 1000000)])])
    fixture('s5-frescura-fronteras',
            'Frescura (RN-12, K22, K40): hoy 2026-12-31 y precios con edad de 0 y de 29 a 91 días. Peso 1 hasta el día 29; (90 - días) / 61 de 30 a 90 '
            '(0.9836, 0.4918, 0); día 90 pesa 0 y no compite; día 91 no compite. La diferencia entre días consecutivos es 1/61. Las 5 fuentes a los 5 días: '
            'compiten manual, ticket_manual y ticket_ocr; no compiten site_assisted ni profeco. Una fecha futura es future_date.',
            [('decision', 'K40'), ('decision', 'K22'), ('regla', 'RN-12'), ('criterio', 'S5-04')],
            {'today': hoy.isoformat(), 'catalog': cat}, consultas, esperado)
    # parámetros configurables (K44): viejo a los 10 días, peso cero a los 20
    consultas, esperado = [], {}
    for n in (9, 10, 15, 20, 21):
        id_ = f'edad-{n:03d}'
        consultas.append(consulta(id_, 'freshness', observed_date=f(n), source='ticket_ocr'))
        esperado[id_] = esp(n, 10, 20)
    assert esperado['edad-015']['weight'] == {'num': 5, 'den': 11}
    fixture('s5-frescura-parametros',
            'K44: con price.stale_after_days = 10 y price.zero_weight_after_days = 20 el peso es (20 - días) / 11 entre los días 10 y 20 (15 días = 5/11).',
            [('decision', 'K44')],
            {'today': hoy.isoformat(), 'settings': {'stale_after_days': 10, 'zero_weight_after_days': 20}, 'catalog': cat}, consultas, esperado)


# ---------------------------------------------------------------- fuentes y promociones
def fuentes():
    pres = 'arroz-grano-largo.sos-uruguayo-1kg'
    cod = ['tienda-manual', 'tienda-ticket', 'tienda-ocr', 'tienda-sitio', 'tienda-profeco']
    pr = [precio(pres, 'tienda-manual', 2900, '2026-10-10', 'manual'), precio(pres, 'tienda-ticket', 2900, '2026-10-10', 'ticket_manual'),
          precio(pres, 'tienda-ocr', 2900, '2026-10-10', 'ticket_ocr', ref='attachment 7'),
          precio(pres, 'tienda-sitio', 1000, '2026-10-10', 'site_assisted'), precio(pres, 'tienda-profeco', 900, '2026-10-10', 'profeco', ref='QQP 2026 a.csv:12')]
    datos = {'today': HOY, 'catalog_file': CAT, 'stores': [tienda(c) for c in cod], 'prices': pr}
    q = [consulta(f'precio-{c}', 'current_price', presentation=pres, store=c) for c in cod]
    q += [consulta('cotiza-arroz', 'quote_for', product='arroz-grano-largo'),
          consulta('cotiza-arroz-sitio', 'quote_for', product='arroz-grano-largo', store='tienda-sitio'),
          consulta('cotiza-arroz-profeco', 'quote_for', product='arroz-grano-largo', store='tienda-profeco')]
    vig = lambda src: {'price': {'price_cents': 2900, 'observed_date': '2026-10-10', 'source': src}, 'references': [], 'promotion': None}
    e = {'precio-tienda-manual': vig('manual'), 'precio-tienda-ticket': vig('ticket_manual'), 'precio-tienda-ocr': vig('ticket_ocr'),
         'precio-tienda-sitio': {'price': None, 'references': [{'source': 'site_assisted', 'price_cents': 1000, 'observed_date': '2026-10-10'}], 'promotion': None},
         'precio-tienda-profeco': {'price': None, 'references': [{'source': 'profeco', 'price_cents': 900, 'observed_date': '2026-10-10'}], 'promotion': None},
         'cotiza-arroz': {'quote': {'product': 'arroz-grano-largo', 'store': 'tienda-manual', 'presentation': pres, 'price_cents': 2900,
                                    'size_milli': 1000000, 'observed_date': '2026-10-10', 'is_stale': False, 'recommended': True}},
         'cotiza-arroz-sitio': {'quote': None}, 'cotiza-arroz-profeco': {'quote': None}}
    fixture('s5-fuentes-compiten',
            'Cinco fuentes, un precio cada una del 2026-10-10 (K22, K37, S5-04, S5-10): manual, ticket_manual y ticket_ocr valen 2900 y compiten con el mismo peso; '
            'site_assisted (1000) y profeco (900) son más baratos pero solo referencia: no aparecen como precio vigente ni ganan. Empate de las tres que compiten: gana la primera de `stores`.',
            [('decision', 'K37'), ('decision', 'K22'), ('criterio', 'S5-04'), ('criterio', 'S5-10')], datos, q, e)
    # promociones (vigencia: PROVISIONAL, bloqueo S5-reglas-comparativa pregunta 3, opción A: solo en su observed_date)
    pr = [precio(pres, 'tienda-manual', 2900, '2026-10-10', 'manual'), precio(pres, 'tienda-manual', 2000, '2026-10-05', 'manual', promo=True),
          precio(pres, 'tienda-ticket', 2900, '2026-10-10', 'ticket_manual'), precio(pres, 'tienda-ticket', 2500, HOY, 'ticket_manual', promo=True)]
    datos = {'today': HOY, 'catalog_file': CAT, 'stores': [tienda('tienda-manual'), tienda('tienda-ticket')], 'prices': pr}
    q = [consulta('promo-vencida', 'current_price', presentation=pres, store='tienda-manual'),
         consulta('promo-vigente', 'current_price', presentation=pres, store='tienda-ticket'),
         consulta('cotiza-arroz', 'quote_for', product='arroz-grano-largo')]
    e = {'promo-vencida': {'price': {'price_cents': 2900, 'observed_date': '2026-10-10', 'source': 'manual'}, 'references': [], 'promotion': None},
         'promo-vigente': {'price': {'price_cents': 2900, 'observed_date': '2026-10-10', 'source': 'ticket_manual'}, 'references': [],
                           'promotion': {'price_cents': 2500, 'observed_date': HOY}},
         'cotiza-arroz': {'quote': {'product': 'arroz-grano-largo', 'store': 'tienda-manual', 'presentation': pres, 'price_cents': 2900,
                                    'size_milli': 1000000, 'observed_date': '2026-10-10', 'is_stale': False, 'recommended': True}}}
    fixture('s5-promocion-vigencia',
            'Promociones (RN-12, S5-10): la del 2026-10-05 (2000) está vencida y no se muestra ni cuenta; la de hoy (2500) se muestra junto al precio regular, '
            'y ninguna entra en la cotización (gana el regular 2900, primera tienda).',
            [('regla', 'RN-12'), ('criterio', 'S5-10')], datos, q, e,
            prov('S5-reglas-comparativa', 'pregunta 3: una promoción solo es vigente en su observed_date (opción A)'))


# ---------------------------------------------------------------- canasta del 10-oct y mejor tienda
def canasta():
    idn, sus = REP['items_identicos'], REP['sustitutos']
    pr, lst = [], []
    for i in idn:
        prod, pres, size, _ = NOMBRE[i['producto']]
        for t, k in (('bodega', 'bodega_unit'), ('walmart', 'walmart_unit')):
            pr.append(precio(pres, t, cts(i[k]), '2026-10-10'))
        cant = i['cantidad'] + (1 if i['producto'].startswith('Atún') else 0)
        lst.append(lista_item(prod, nec=cant * size, rec=pres))
    cuota = {'arroz-grano-largo': (2250000, 'arroz-grano-largo.sos-uruguayo-1kg'), 'sustituto-crema-cafe': (640000, 'sustituto-crema-cafe.coffee-mate-pack-780g'),
             'leche': (4000000, 'leche.santa-clara-entera-6x1l')}
    for s in sus:
        for t, k, c in (('bodega', 'bodega', 'cant_b'), ('walmart', 'walmart', 'cant_w')):
            pr.append(precio(NOMBRE[s[k]][1], t, cts(s[k + '_total']) // s[c], '2026-10-10'))
        p = NOMBRE[s['bodega']][0]
        lst.append(lista_item(p, nec=cuota[p][0], rec=cuota[p][1]))
    tb = cts(REP['totales']['bodega'])
    tw = cts(REP['totales']['walmart']) + cts(REP['items_identicos'][[i['producto'].startswith('Atún') for i in idn].index(True)]['walmart_unit'])
    datos = {'today': HOY, 'catalog_file': CAT, 'stores': TIENDAS_CANON, 'prices': pr, 'shopping_list': lst}

    def cot(prod, tienda_, pres, cents, size):
        return {'quote': {'product': prod, 'store': tienda_, 'presentation': pres, 'price_cents': cents, 'size_milli': size,
                          'observed_date': '2026-10-10', 'is_stale': False, 'recommended': True}}
    q = [consulta('canasta', 'compare_basket'),
         consulta('mejor-arroz', 'quote_for', product='arroz-grano-largo'), consulta('mejor-crema', 'quote_for', product='sustituto-crema-cafe'),
         consulta('mejor-leche', 'quote_for', product='leche'), consulta('mejor-aceite-soya', 'quote_for', product='aceite-soya'),
         consulta('mejor-frijoles', 'quote_for', product='frijoles-refritos')]
    aceite = [i for i in idn if i['producto'].startswith('Aceite de soya')][0]
    frij = [i for i in idn if i['producto'].startswith('Frijoles bayos')][0]
    assert aceite['walmart_unit'] < aceite['bodega_unit'] and frij['bodega_unit'] < frij['walmart_unit']
    e = {'canasta': {'stores': [
            {'store': 'bodega', 'total_cents': tb, 'identical_count': 34, 'substitute_count': 2, 'missing_count': 0, 'stale_count': 0, 'difference_vs_cheapest_cents': 0},
            {'store': 'walmart', 'total_cents': tw, 'identical_count': 35, 'substitute_count': 1, 'missing_count': 0, 'stale_count': 0,
             'difference_vs_cheapest_cents': tw - tb}], 'cheapest_store': 'bodega'},
         'mejor-arroz': cot('arroz-grano-largo', 'walmart', 'arroz-grano-largo.sos-uruguayo-1kg', 2800, 1000000),
         'mejor-crema': cot('sustituto-crema-cafe', 'bodega', 'sustituto-crema-cafe.coffee-mate-pack-780g', 9900, 780000),
         'mejor-leche': cot('leche', 'walmart', 'leche.santa-clara-entera-6x1l', 17000, 6000000),
         'mejor-aceite-soya': cot('aceite-soya', 'walmart', 'aceite-soya.nutrioli-946ml', cts(aceite['walmart_unit']), 946000),
         'mejor-frijoles': cot('frijoles-refritos', 'bodega', 'frijoles-refritos.la-costena-220g', cts(frij['bodega_unit']), 220000)}
    fixture('s5-mejor-tienda-canasta',
            'Canasta por necesidad con los precios del 10-oct-2026 de data.json (§7.3, §7.4). Necesidades: las piezas de los 33 renglones idénticos (el atún, 6 latas) y, '
            'de los sustitutos, arroz 2250 g, crema 640 g y leche 4 l; cada tienda cubre con el menor costo de paquetes completos. Bodega suma 3,094.50 (igual a K-01: '
            'compra las mismas piezas); Walmart, 3,233.00 + 21.00 de la sexta lata = 3,254.00. Bodega es la más barata. Por ingrediente (unidad base): arroz Walmart 1 kg, '
            'crema Bodega 780 g, leche Walmart 6 l, aceite de soya Walmart, frijoles Bodega. Con las presentaciones recomendadas de la lista: Bodega 34 idénticos y 2 sustitutos '
            '(arroz 750 g, leche 4 l); Walmart 35 y 1 (crema 640 g).',
            [('data.json', 'totales'), ('criterio', 'S5-03')], datos, q, e)


def empates():
    cat = catalogo([('a', 'g', 'abarrotes', [('a-1kg', 1000000)]), ('b', 'g', 'abarrotes', [('b-1kg', 1000000)]), ('c', 'g', 'abarrotes', [('c-1kg', 1000000)])])
    pv = prov('S5-reglas-comparativa', 'pregunta 2: desempate por renglones recomendados, luego precios viejos, luego orden de `stores` (opción A)')
    nec = lambda p: lista_item(p, nec=1000000, rec=p + '-1kg')
    # a) mismo total; gana la que tiene más renglones recomendados aunque vaya segunda en `stores`
    pr = [precio('a-1kg', 't1', 1100, '2026-10-10'), precio('a-1kg', 't2', 1000, '2026-10-10'), precio('b-1kg', 't1', 1900, '2026-10-10'), precio('b-1kg', 't2', 2100, '2026-10-10'),
          precio('c-1kg', 't1', 1100, '2026-10-10'), precio('c-1kg', 't2', 1000, '2026-10-10')]
    st = lambda *cs: [tienda(c) for c in cs]
    e = {'canasta': {'stores': [{'store': 't2', 'total_cents': 4100, 'difference_vs_cheapest_cents': 0}, {'store': 't1', 'total_cents': 4100, 'difference_vs_cheapest_cents': 0}],
                     'items': [{'recommended_store': 't2'}, {'recommended_store': 't1'}, {'recommended_store': 't2'}], 'cheapest_store': 't2'}}
    fixture('s5-mejor-tienda-empate-recomendados',
            'Empate de total en la canasta (4100 las dos): t2 es la recomendada en 2 renglones (a y c) y t1 en 1 (b), así que t2 va primero aunque t1 vaya antes en `stores`.',
            [('regla', 'RN-12'), ('criterio', 'S5-03')],
            {'today': HOY, 'catalog': cat, 'stores': st('t1', 't2'), 'prices': pr, 'shopping_list': [nec('a'), nec('b'), nec('c')]},
            [consulta('canasta', 'compare_basket')], e, pv)
    # b) mismo total y mismos renglones recomendados (1 y 1): gana la de menos precios viejos
    pr = [precio('a-1kg', 't1', 1000, '2026-10-10'), precio('a-1kg', 't2', 1100, '2026-09-02'), precio('b-1kg', 't1', 2100, '2026-10-10'), precio('b-1kg', 't2', 2000, '2026-10-10')]
    e = {'canasta': {'stores': [{'store': 't1', 'total_cents': 3100, 'stale_count': 0}, {'store': 't2', 'total_cents': 3100, 'stale_count': 1}],
                     'items': [{'recommended_store': 't1'}, {'recommended_store': 't2'}], 'cheapest_store': 't1'}}
    fixture('s5-mejor-tienda-empate-viejos',
            'Empate de total (3100) y de renglones recomendados (1 y 1): el precio de t2 para a tiene 40 días (viejo, no es la recomendada: su precio ajustado es mayor), '
            'así que t1 (0 viejos) va antes que t2 (1 viejo) aunque t2 vaya primero en `stores`.',
            [('regla', 'RN-12'), ('criterio', 'S5-03')],
            {'today': HOY, 'catalog': cat, 'stores': st('t2', 't1'), 'prices': pr, 'shopping_list': [nec('a'), nec('b')]},
            [consulta('canasta', 'compare_basket')], e, pv)
    # c) todo empata: orden de `stores`
    pr = [precio('a-1kg', 'zeta', 1000, '2026-10-10'), precio('a-1kg', 'alfa', 1100, '2026-10-10'), precio('b-1kg', 'zeta', 2100, '2026-10-10'), precio('b-1kg', 'alfa', 2000, '2026-10-10')]
    e = {'canasta': {'stores': [{'store': 'zeta', 'total_cents': 3100}, {'store': 'alfa', 'total_cents': 3100}],
                     'items': [{'recommended_store': 'zeta'}, {'recommended_store': 'alfa'}], 'cheapest_store': 'zeta'}}
    fixture('s5-mejor-tienda-empate-orden',
            'Empate total (3100), de renglones recomendados (1 y 1) y de precios viejos (0): queda el orden de `stores` (zeta antes que alfa, no el alfabético).',
            [('regla', 'RN-12'), ('criterio', 'S5-03')],
            {'today': HOY, 'catalog': cat, 'stores': st('zeta', 'alfa'), 'prices': pr, 'shopping_list': [nec('a'), nec('b')]},
            [consulta('canasta', 'compare_basket')], e, pv)
    # d) canasta incompleta: la tienda sin precio de b se lista pero no puede ser la más barata
    pr = [precio('a-1kg', 't1', 1000, '2026-10-10'), precio('a-1kg', 't2', 1100, '2026-10-10'), precio('a-1kg', 't3', 500, '2026-10-10'),
          precio('b-1kg', 't1', 2100, '2026-10-10'), precio('b-1kg', 't2', 2000, '2026-10-10')]
    e = {'canasta': {'stores': [{'store': 't3', 'total_cents': 500, 'missing_count': 1}, {'store': 't2', 'total_cents': 3100, 'missing_count': 0},
                                {'store': 't1', 'total_cents': 3100, 'missing_count': 0}], 'cheapest_store': 't2'}}
    fixture('s5-mejor-tienda-incompleta',
            't3 solo tiene precio de a (500): se lista con missing_count 1 y total 500 (ordenada por total ascendente), pero la más barata es t2: gana el empate de total con t1 '
            'porque es la recomendada de b (t3 lo es de a, t1 de ninguno).',
            [('regla', 'RN-12'), ('criterio', 'S5-03')],
            {'today': HOY, 'catalog': cat, 'stores': st('t1', 't2', 't3'), 'prices': pr, 'shopping_list': [nec('a'), nec('b')]},
            [consulta('canasta', 'compare_basket')], e, pv)


# ---------------------------------------------------------------- cotización recomendada y regla única (§14.1, §7.3; ILUSTRATIVOS)
def cotizacion():
    cat = catalogo([('arroz', 'g', 'abarrotes', [('arroz-1kg', 1000000), ('arroz-750', 750000), ('arroz-500', 500000)]),
                    ('frijol', 'g', 'abarrotes', [('frijol-1kg', 1000000), ('frijol-500', 500000)])])
    pr = [precio('arroz-1kg', 'walmart', 2800, '2026-10-10'), precio('arroz-750', 'walmart', 2200, '2026-10-10'),
          precio('arroz-750', 'bodega', 2200, '2026-10-10'), precio('arroz-500', 'bodega', 1500, '2026-10-10'),
          precio('frijol-1kg', 'walmart', 3000, '2026-10-10'), precio('frijol-500', 'bodega', 1500, '2026-10-10')]
    stores = [tienda('walmart'), tienda('bodega')]
    qte = lambda prod, st, pres, c, size, rec: {'quote': {'product': prod, 'store': st, 'presentation': pres, 'price_cents': c, 'size_milli': size,
                                                           'observed_date': '2026-10-10', 'is_stale': False, 'recommended': rec}}
    q = [consulta('arroz-libre', 'quote_for', product='arroz'), consulta('arroz-bodega', 'quote_for', product='arroz', store='bodega'),
         consulta('arroz-750', 'quote_for', product='arroz', presentation='arroz-750'), consulta('frijol-libre', 'quote_for', product='frijol'),
         consulta('canasta', 'compare_basket')]
    e = {'arroz-libre': qte('arroz', 'walmart', 'arroz-1kg', 2800, 1000000, True),
         'arroz-bodega': qte('arroz', 'bodega', 'arroz-750', 2200, 750000, False),
         'arroz-750': qte('arroz', 'walmart', 'arroz-750', 2200, 750000, False),
         'frijol-libre': qte('frijol', 'bodega', 'frijol-500', 1500, 500000, True),
         'canasta': {'stores': [{'store': 'bodega', 'total_cents': 1500, 'difference_vs_cheapest_cents': 0, 'substitute_count': 1, 'identical_count': 0},
                                {'store': 'walmart', 'total_cents': 2200, 'difference_vs_cheapest_cents': 700, 'substitute_count': 1, 'identical_count': 0}],
                     'items': [{'product': 'arroz', 'quantity_milli': 500000, 'recommended_store': 'walmart', 'recommended_presentation': 'arroz-1kg',
                                'estimated_price_cents': 2800, 'is_stale_price': False,
                                'prices': [{'store': 'walmart', 'presentation': 'arroz-750', 'price_cents': 2200, 'base_unit_price_cents': 2933,
                                            'observed_date': '2026-10-10', 'is_substitute': True},
                                           {'store': 'bodega', 'presentation': 'arroz-500', 'price_cents': 1500, 'base_unit_price_cents': 3000,
                                            'observed_date': '2026-10-10', 'is_substitute': True}]}],
                     'cheapest_store': 'bodega'}}
    fixture('s5-cotizacion-recomendada',
            'Ejemplo de §14.1 (hoy 2026-10-12; precios del 10-oct, peso 1; stores = walmart, bodega). Arroz: Walmart 1 kg 2800 y 750 g 2200; Bodega 750 g 2200 y 500 g 1500. '
            'Sin elección gana Walmart 1 kg (2800 por kg contra 2933 y 3000); con store bodega, 750 g (2933 contra 3000); con presentación 750 g, empatan en unidad base '
            'y en precio y gana Walmart por el orden de `stores`. Frijol: Walmart 1 kg 3000 y Bodega 500 g 1500 empatan por kg (3000); gana Bodega por menor precio de paquete. '
            'Canasta (necesidad 500 g): recomendada Walmart 1 kg (2800); cada tienda cubre con su menor costo, Walmart 750 g a 2200 y Bodega 500 g a 1500, los dos sustitutos; '
            'Bodega la más barata, Walmart +700.',
            [('criterio', 'S5-03'), ('manual', '§14.1 y §7.3 de S5-interfaz.md')],
            {'today': HOY, 'catalog': cat, 'stores': stores, 'prices': pr, 'shopping_list': [lista_item('arroz', nec=500000, rec='arroz-1kg')]}, q, e,
            prov('S5-reglas-comparativa', 'pregunta 2: desempate por precio del paquete, orden de `stores` y orden del catálogo (opción A)'))
    # regla única (contraejemplo de M-r4-1)
    cat = catalogo([('arroz', 'g', 'abarrotes', [('arroz-1kg', 1000000), ('arroz-500', 500000)]), ('aceite-soya', 'ml', 'abarrotes', [('aceite-946', 946000)])])
    pr = [precio('arroz-1kg', 'walmart', 2800, '2026-10-10'), precio('arroz-500', 'walmart', 1450, '2026-10-10'), precio('arroz-500', 'bodega', 1460, '2026-10-10'),
          precio('aceite-946', 'walmart', 3990, '2026-10-10'), precio('aceite-946', 'bodega', 4000, '2026-10-10')]
    lista = [lista_item('arroz', nec=500000, rec='arroz-1kg'), lista_item('aceite-soya', cant=1892000, rec='aceite-946')]
    e = {'canasta': {'stores': [{'store': 'walmart', 'total_cents': 9430, 'difference_vs_cheapest_cents': 0, 'identical_count': 1, 'substitute_count': 1},
                                {'store': 'bodega', 'total_cents': 9460, 'difference_vs_cheapest_cents': 30, 'identical_count': 1, 'substitute_count': 1}],
                     'items': [{'product': 'arroz', 'quantity_milli': 500000, 'recommended_store': 'walmart', 'recommended_presentation': 'arroz-1kg', 'estimated_price_cents': 2800,
                                'prices': [{'store': 'walmart', 'presentation': 'arroz-500', 'price_cents': 1450, 'base_unit_price_cents': 2900, 'is_substitute': True},
                                           {'store': 'bodega', 'presentation': 'arroz-500', 'price_cents': 1460, 'base_unit_price_cents': 2920, 'is_substitute': True}]},
                               {'product': 'aceite-soya', 'quantity_milli': 1892000, 'recommended_store': 'walmart', 'recommended_presentation': 'aceite-946', 'estimated_price_cents': 7980,
                                'prices': [{'store': 'walmart', 'presentation': 'aceite-946', 'price_cents': 3990, 'base_unit_price_cents': 4218, 'is_substitute': False},
                                           {'store': 'bodega', 'presentation': 'aceite-946', 'price_cents': 4000, 'base_unit_price_cents': 4228, 'is_substitute': False}]}],
                     'cheapest_store': 'walmart'}}
    fixture('s5-canasta-regla-unica',
            'Contraejemplo de la revisión r4 (§7.3): el total de cada tienda usa una sola regla. Arroz (necesidad 500 g): Walmart 1 kg 2800 y 500 g 1450, Bodega 500 g 1460; '
            'recomendado Walmart 1 kg (2800), pero Walmart cubre con 500 g a 1450. Aceite de soya `manual` (cantidad fijada 1892 ml, rec. 946 ml): Walmart 3990 y Bodega 4000 por botella, '
            '2 paquetes: 7980 y 8000. Totales Walmart 9430 y Bodega 9460 (+30); la más barata es Walmart. Si la tienda recomendada se cotizara con la presentación recomendada, '
            'Walmart saldría en 10780 y la canasta señalaría a Bodega.',
            [('criterio', 'S5-03'), ('manual', '§7.3 de S5-interfaz.md y REVISION-contrato-S4S5-r4.md')],
            {'today': HOY, 'catalog': cat, 'stores': [tienda('walmart'), tienda('bodega')], 'prices': pr, 'shopping_list': lista},
            [consulta('canasta', 'compare_basket')], e,
            prov('S5-reglas-comparativa', 'pregunta 2: desempate (aquí no interviene; la recomendación es única por unidad base)'))


# ---------------------------------------------------------------- membresía
def membresia():
    cat = catalogo([('a', 'g', 'abarrotes', [('a-1kg', 1000000)]), ('b', 'g', 'abarrotes', [('b-1kg', 1000000)]), ('c', 'g', 'abarrotes', [('c-1kg', 1000000)])])
    stores = [tienda('sams', requires=True, fee=100000), tienda('sams-b', requires=True, fee=80000), tienda('sams-sin-costo', requires=True), tienda('chedraui'),
              tienda('tianguis', kind='market', visit_weekday=2, groups=['verduras'])]
    hoy = dt.date.fromisoformat(HOY)
    compras = [{'date': (hoy - dt.timedelta(days=n)).isoformat(), 'store': 'chedraui'} for n in (3, 40, 80, 120, 200, 300)]
    compras += [{'date': (hoy - dt.timedelta(days=400)).isoformat(), 'store': 'chedraui'}, {'date': (hoy - dt.timedelta(days=10)).isoformat(), 'store': 'tianguis'}]
    consultas = [consulta(i, 'membership', store=s, savings_per_purchase_cents=a) for i, s, a in
                 (('exacto', 'sams', 12500), ('hacia-arriba', 'sams', 12000), ('rn13', 'sams', 13850), ('rn13-barata', 'sams-b', 13850),
                  ('sin-costo', 'sams-sin-costo', 12500), ('sin-ahorro', 'sams', 0))]
    ent = lambda fee, a: {'store': None}
    e = {'exacto': {'basket_savings_cents': 12500, 'annual_fee_cents': 100000, 'break_even_purchases': 8, 'assumes_membership': True},
         'sin-costo': {'annual_fee_cents': None, 'break_even_purchases': None, 'verdict': 'missing_fee'},
         'sin-ahorro': {'basket_savings_cents': 0, 'break_even_purchases': None, 'verdict': 'not_worth_it'}}
    fixture('s5-membresia-equilibrio',
            'Punto de equilibrio de Sam\'s (S5-09, RN-13; ILUSTRATIVO): costo 1,000.00 y ahorro por compra 125.00 dan 8 compras al año (división exacta); costo ausente -> '
            '`missing_fee` ("sin costo capturado", sin número inventado); ahorro 0 -> equilibrio ausente y `not_worth_it`. Sam\'s compite como si tuvieran membresía (assumes_membership).',
            [('criterio', 'S5-09'), ('regla', 'RN-13'), ('decision', 'K7')],
            {'today': HOY, 'catalog': cat, 'stores': stores, 'grocery_purchases': compras}, consultas, e)
    e = {'hacia-arriba': {'break_even_purchases': 9, 'purchases_per_year': 6, 'annual_savings_cents': 72000, 'verdict': 'not_worth_it'},
         'rn13': {'break_even_purchases': 8, 'purchases_per_year': 6, 'annual_savings_cents': 83100, 'verdict': 'not_worth_it'},
         'rn13-barata': {'annual_fee_cents': 80000, 'break_even_purchases': 6, 'annual_savings_cents': 83100, 'verdict': 'worth_it'},
         'exacto': {'annual_savings_cents': 75000, 'verdict': 'not_worth_it'}}
    fixture('s5-membresia-anual',
            'Anualización del punto de equilibrio (S5-09, ILUSTRATIVO): compras al año = compras de despensa en cadenas en los últimos 365 días (6 en cadena; la de 400 días y la del tianguis no cuentan). '
            '1000.00 / 120.00 = 8.33 -> 9 (hacia arriba: con 8 no se cubre). RN-13: 138.50 x 6 = 831.00 (no cubre 1,000.00 pero sí 800.00 de `sams-b`, que equilibra en 6).',
            [('criterio', 'S5-09'), ('regla', 'RN-13')],
            {'today': HOY, 'catalog': cat, 'stores': stores, 'grocery_purchases': compras}, consultas, e,
            prov('S5-reglas-comparativa', 'preguntas 4a y 4b: compras al año = despensas en cadenas de los últimos 365 días; equilibrio hacia arriba'))
    # ahorro calculado con la canasta: A 200, B 0, C 200 -> 400; equilibrio 100000 / 400 = 250 (exacto)
    pr = [precio('a-1kg', 'sams', 1000, '2026-10-10'), precio('a-1kg', 'chedraui', 1200, '2026-10-10'), precio('a-1kg', 'sams-b', 700, '2026-10-10'),
          precio('b-1kg', 'sams', 2500, '2026-10-10'), precio('b-1kg', 'chedraui', 2300, '2026-10-10'),
          precio('c-1kg', 'sams', 800, '2026-10-10'), precio('c-1kg', 'chedraui', 1000, '2026-10-10')]
    lista = [lista_item('a', nec=1000000, rec='a-1kg'), lista_item('b', nec=1000000, rec='b-1kg'), lista_item('c', nec=2000000, rec='c-1kg')]
    fixture('s5-membresia-canasta',
            'Ahorro de la canasta de Sam\'s (RN-13): solo cuentan los renglones donde Sam\'s es más barato que la mejor tienda sin membresía, con las cantidades de la lista. '
            'a: 1000 contra 1200 -> 200; b: 2500 contra 2300 -> 0; c (2 kg): 1600 contra 2000 -> 400; ahorro 600. '
            '`sams-b` (700 en a) también exige membresía y no cuenta como comparador. Un artículo sin precio en Sam\'s no suma.',
            [('regla', 'RN-13')],
            {'today': HOY, 'catalog': cat, 'stores': stores, 'prices': pr, 'shopping_list': lista},
            [consulta('sams', 'membership', store='sams')],
            {'sams': {'basket_savings_cents': 600, 'annual_fee_cents': 100000, 'assumes_membership': True}})


# ---------------------------------------------------------------- mercados, precios viejos, lectura asistida, ingesta
def mercado():
    pr = [precio('uva-verde.pieza-1kg', 'tianguis', 6000, '2026-10-10'), precio('uva-verde.pieza-1kg', 'bodega', 8700, '2026-10-10'),
          precio('avena.3-minutos-400g', 'tianguis', 100, '2026-10-10'), precio('avena.3-minutos-400g', 'bodega', 3600, '2026-10-10')]
    stores = [tienda('bodega'), tienda('tianguis', kind='market', name='Tianguis de martes', visit_weekday=2, groups=['frutas-verduras'])]
    cotiza = lambda prod, st, pres, c, size: {'quote': {'product': prod, 'store': st, 'presentation': pres, 'price_cents': c, 'size_milli': size,
                                                        'observed_date': '2026-10-10', 'is_stale': False, 'recommended': True}}
    q = [consulta('visita-sabado', 'next_visit', store='tianguis', **{'from': '2026-10-10'}),
         consulta('visita-lunes', 'next_visit', store='tianguis', **{'from': '2026-10-12'}),
         consulta('visita-martes', 'next_visit', store='tianguis', **{'from': '2026-10-13'}),
         consulta('visita-miercoles', 'next_visit', store='tianguis', **{'from': '2026-10-14'}),
         consulta('uva', 'quote_for', product='uva-verde'), consulta('avena', 'quote_for', product='avena'),
         consulta('avena-en-tianguis', 'quote_for', product='avena', store='tianguis')]
    e = {'visita-sabado': {'date': '2026-10-13'}, 'visita-lunes': {'date': '2026-10-13'}, 'visita-martes': {'date': '2026-10-13'},
         'visita-miercoles': {'date': '2026-10-20'},
         'uva': cotiza('uva-verde', 'tianguis', 'uva-verde.pieza-1kg', 6000, 1000000),
         'avena': cotiza('avena', 'bodega', 'avena.3-minutos-400g', 3600, 400000), 'avena-en-tianguis': {'quote': None}}
    fixture('s5-mercado-visita',
            'Mercados (RN-16, S5-11): tianguis de martes (2). Desde el sábado 10-oct-2026 la próxima visita es el martes 13; si hoy es día de visita, hoy cuenta (supuesto S-10); '
            'desde el miércoles 14, el martes 20. El tianguis declara solo "frutas-verduras": gana en uva (6000 contra 8700) pero no compite en avena, aunque tenga un precio de 100.',
            [('regla', 'RN-16'), ('criterio', 'S5-11')], {'today': HOY, 'catalog_file': CAT, 'stores': stores, 'prices': pr}, q, e)


def viejos():
    cat = catalogo([(p, 'g', 'abarrotes', [(p + '-1kg', 1000000)]) for p in ('p1', 'p2', 'p3', 'p4', 'p5', 'p6')])
    g = dt.date(2026, 11, 9)
    fch = lambda n: (g - dt.timedelta(days=n)).isoformat()
    pr = [precio('p1-1kg', 't1', 1000, fch(30)), precio('p2-1kg', 't1', 1000, fch(28)), precio('p3-1kg', 't2', 1000, fch(89)), precio('p4-1kg', 't1', 1000, fch(90)),
          precio('p6-1kg', 't1', 1000, fch(28)), precio('p6-1kg', 't2', 1000, fch(70))]
    q = [consulta('despensa', 'stale_prices_for_plan', products=['p1', 'p2', 'p3', 'p4', 'p5', 'p6'], assigned_stores={'p6': 't2'}, grocery_date=g.isoformat()),
         consulta('despensa-sin-asignar', 'stale_prices_for_plan', products=['p6'], assigned_stores={}, grocery_date=g.isoformat())]
    def reg(p, st, n, w):
        return {'product': p, 'store': st, 'presentation': p + '-1kg', 'observed_date': fch(n), 'age_days': n, 'weight': fr(w)}
    e = {'despensa': {'stale': [reg('p1', 't1', 30, Fraction(60, 61)), reg('p3', 't2', 89, Fraction(1, 61)), reg('p6', 't2', 70, Fraction(20, 61))],
                      'without_price': ['p4', 'p5']},
         'despensa-sin-asignar': {'stale': [], 'without_price': []}}
    fixture('s5-precios-viejos',
            'Insumo de la alerta de precios viejos (S5-12; ILUSTRATIVO): despensa el 2026-11-09 (se evalúa en esa fecha, no en hoy). p1 a 30 días (viejo, 60/61); p2 a 28 (fresco); '
            'p3 a 89 (viejo, 1/61); p4 a 90 (peso 0: no compite, sin precio); p5 sin precio; p6 con tienda asignada t2 a 70 días (viejo, 20/61) aunque t1 tenga uno fresco; '
            'sin tienda asignada p6 se cotiza en t1 (28 días) y no aparece.',
            [('criterio', 'S5-12')], {'today': HOY, 'catalog': cat, 'stores': [tienda('t1'), tienda('t2')], 'prices': pr}, q, e)


def asistida():
    cat = catalogo([('a', 'g', 'abarrotes', [('a-1kg', 1000000)])])
    stores = [tienda('walmart', assisted=True), tienda('bodega_aurrera', assisted=True), tienda('sams_club', assisted=False)]
    terms = [{'site': 'walmart', 'verdict': 'permitido con baja frecuencia', 'url': 'https://ejemplo.invalid/terminos-walmart', 'read_date': '2026-10-11'},
             {'site': 'bodega_aurrera', 'verdict': 'prohibido', 'url': '', 'read_date': '2026-10-11'},
             {'site': 'sams_club', 'verdict': 'permitido', 'url': 'https://ejemplo.invalid/terminos-sams', 'read_date': '2026-10-11'}]
    q = [consulta(s, 'assisted_can_enable', store=s) for s in ('walmart', 'bodega_aurrera', 'sams_club')]
    e = {'walmart': {'can_enable': True}, 'bodega_aurrera': {'error': 'no_terms_verdict'}, 'sams_club': {'error': 'no_terms_verdict'}}
    fixture('s5-lectura-asistida',
            'Interruptor activado (S5-07): walmart con veredicto completo (URL y fecha) y la bandera de la tienda -> se puede activar; bodega_aurrera con la URL vacía (campo vacío = faltante) '
            'y sams_club con la bandera de la tienda apagada -> no_terms_verdict.',
            [('criterio', 'S5-07')], {'today': HOY, 'settings': {'assisted_reading_enabled': True}, 'catalog': cat, 'stores': stores, 'terms': terms}, q, e)
    e = {'walmart': {'error': 'no_terms_verdict'}, 'bodega_aurrera': {'error': 'no_terms_verdict'}, 'sams_club': {'error': 'no_terms_verdict'}}
    fixture('s5-lectura-asistida-apagada',
            'Por omisión price.assisted_reading_enabled es false: aun con veredicto completo y bandera de tienda, ninguna cadena se puede activar (S5-07, K44).',
            [('criterio', 'S5-07'), ('decision', 'K44')], {'today': HOY, 'catalog': cat, 'stores': stores, 'terms': terms}, q, e)
    fixture('s5-lectura-asistida-sin-veredicto',
            'Interruptor activado y bandera de tienda encendida, pero sin ningún veredicto en `terms`: no_terms_verdict (S5-07).',
            [('criterio', 'S5-07')], {'today': HOY, 'settings': {'assisted_reading_enabled': True}, 'catalog': cat, 'stores': stores}, q[:1], {'walmart': {'error': 'no_terms_verdict'}})


def ingesta():
    cat = catalogo([('a', 'g', 'abarrotes', [('a-1kg', 1000000)])])
    base = [precio('a-1kg', 't1', 1000, '2026-10-10')]
    casos = {'fecha-futura': ([precio('a-1kg', 't1', 1000, '2026-10-13')], 'future_date', 'una fecha posterior a hoy (2026-10-12) es error de captura (S-7)'),
             'presentacion-desconocida': ([precio('a-2kg', 't1', 1000, '2026-10-10')], 'unknown_presentation', 'la presentación no existe'),
             'tienda-desconocida': ([precio('a-1kg', 't9', 1000, '2026-10-10')], 'unknown_store', 'la tienda no existe'),
             'precio-duplicado': (base + base, 'duplicate_price', 'misma (presentación, tienda, fecha, fuente, promoción) dos veces')}
    for nombre, (pr, code, desc) in casos.items():
        fixture(f's5-ingesta-{nombre}', f'Validación común de la ingesta (§3, S5-10): {desc}. Toda consulta del caso sale con el error.',
                [('criterio', 'S5-10')], {'today': HOY, 'catalog': cat, 'stores': [tienda('t1')], 'prices': pr},
                [consulta('precio', 'current_price', presentation='a-1kg', store='t1')], {'precio': {'error': code}})


# ---------------------------------------------------------------- PROFECO y OCR (archivos de _datos/)
def profeco():
    DATOS.mkdir(parents=True, exist_ok=True)
    cab = 'producto,presentacion,marca,categoria,catalogo,precio,fecha_registro,cadena_comercial,giro,nombre_comercial,direccion,estado,municipio,latitud,longitud'

    def fila(prod, pres, marca, precio_, fecha_, cadena, nombre, muni):
        return f'{prod},{pres},{marca},ABARROTES,CATALOGO,{precio_},{fecha_},{cadena},SUPERMERCADO,{nombre},"CALLE 1, COL. UNO",CIUDAD DE MEXICO,{muni},19.45,-99.15'
    a = [fila('LECHE ENTERA', '1 L', 'SANTA CLARA', '24.50', '2026-08-15', 'BODEGA AURRERA', 'BODEGA AURRERA VALLEJO', 'AZCAPOTZALCO'),
         fila('LECHE ENTERA', '1 L', 'SANTA CLARA', '24.90', '15/08/2026', 'BODEGA AURRERA', 'BODEGA AURRERA VALLEJO', 'AZCAPOTZALCO'),
         fila('LECHE ENTERA', '1 L', 'SANTA CLARA', '24.50', 'ayer', 'BODEGA AURRERA', 'BODEGA AURRERA VALLEJO', 'AZCAPOTZALCO'),
         fila('LECHE ENTERA', '1 L', 'SANTA CLARA', '25.00', '2026-08-16', 'COMERCIAL OTRA', 'TIENDA OTRA', 'AZCAPOTZALCO'),
         fila('ARROZ SOS', '750 G', 'SOS', '22.00', '2026-08-17', 'BODEGA AURRERA', 'BODEGA AURRERA VALLEJO', 'AZCAPOTZALCO'),
         fila('FRIJOLES REFRITOS', '220 G', 'LA COSTEÑA', '12.00', '2026-08-20 10:30:00', 'BODEGA AURRERA', 'BODEGA AURRERA VALLEJO', 'AZCAPOTZALCO')]
    (DATOS / 'profeco-muestra-a.csv').write_bytes(b'\xef\xbb\xbf' + (cab + '\r\n' + '\r\n'.join(a) + '\r\n').encode('utf-8'))
    b = [fila('ATÚN DOLORES', '140 G', 'DOLORES', '21.00', '2026-08-31', 'WALMART', 'WALMART ACUEDUCTO', 'GUSTAVO A. MADERO'),
         'ATÚN DOLORES,140 G,DOLORES,ABARROTES,CATALOGO,21.00,2026-08-31,WALMART,SUPERMERCADO,WALMART ACUEDUCTO,"CALLE 2",CIUDAD DE MEXICO,GUSTAVO A. MADERO,19.45',
         fila('ATÚN DOLORES', '140 G', 'DOLORES', 'abc', '2026-08-31', 'WALMART', 'WALMART ACUEDUCTO', 'GUSTAVO A. MADERO'),
         fila('ATÚN DOLORES', '140 G', 'DOLORES', '21.00', '2026-08-31', 'WALMART', 'WALMART TOLUCA', 'TOLUCA')]
    (DATOS / 'profeco-muestra-b.csv').write_bytes((cab + '\n' + '\n'.join(b) + '\n').encode('cp1252'))
    (DATOS / 'profeco-vinculos.csv').write_text(
        'producto,presentacion,marca,presentacion_id,estado,confirmo,fecha_confirmacion\n'
        'LECHE ENTERA,1 L,SANTA CLARA,leche.santa-clara-entera-6x1l,confirmado,axel,2026-10-11\n'
        'FRIJOLES REFRITOS,220 G,LA COSTEÑA,frijoles-refritos.la-costena-220g,propuesto,,\n'
        'ATÚN DOLORES,140 G,DOLORES,atun-agua.dolores-140g,confirmado,ximena,2026-10-11\n', encoding='utf-8')
    stores = [tienda('bodega', profeco_name='BODEGA AURRERA'), tienda('bodega-vallejo', kind='branch', parent='bodega', profeco_name='BODEGA AURRERA VALLEJO', municipality='AZCAPOTZALCO'),
              tienda('walmart', profeco_name='WALMART'),
              tienda('walmart-tepeyac', kind='branch', parent='walmart', profeco_name='WALMART TEPEYAC', municipality='GUSTAVO A. MADERO')]
    d = 'oraculo/fixtures/s5/_datos/'
    args = lambda modo: dict(files=[d + 'profeco-muestra-a.csv', d + 'profeco-muestra-b.csv'], links_file=d + 'profeco-vinculos.csv', mode=modo)
    rej = [{'file': d + 'profeco-muestra-a.csv', 'row': 3}, {'file': d + 'profeco-muestra-b.csv', 'row': 2}, {'file': d + 'profeco-muestra-b.csv', 'row': 3}]
    e = {'produccion': {'read': 10, 'rejected': 3, 'filtered_out': 2, 'skipped': 2, 'loaded': 2, 'duplicates': 1,
                        'coverage': [{'chain': 'BODEGA AURRERA', 'municipality': 'AZCAPOTZALCO', 'price_count': 1}, {'chain': 'WALMART', 'municipality': 'GUSTAVO A. MADERO', 'price_count': 1}],
                        'rejections': rej},
         'revision': {'read': 10, 'rejected': 3, 'filtered_out': 2, 'skipped': 1, 'loaded': 3, 'duplicates': 1,
                      'coverage': [{'chain': 'BODEGA AURRERA', 'municipality': 'AZCAPOTZALCO', 'price_count': 2}, {'chain': 'WALMART', 'municipality': 'GUSTAVO A. MADERO', 'price_count': 1}],
                      'rejections': rej}}
    fixture('s5-profeco-muestra',
            'Cargador PROFECO sobre 2 CSV de muestra (ILUSTRATIVO; S5-05, §11): a = UTF-8 con BOM, 6 renglones; b = Windows-1252, 4 renglones. Renglón por renglón: a1 se carga; a2 es el mismo '
            'registro del mismo día en otro formato de fecha (15/08/2026) -> duplicado; a3 fecha ilegible -> rechazado; a4 cadena que no es de Axel -> filtrado; a5 sin liga -> omitido; '
            'a6 liga "propuesto" -> omitido en producción y cargado en revisión; b1 se carga (tienda = la cadena, porque ninguna sucursal se llama WALMART ACUEDUCTO); b2 con 14 columnas -> rechazado; '
            'b3 precio ilegible -> rechazado; b4 municipio fuera del filtro -> filtrado. El municipio de Walmart va en su sucursal walmart-tepeyac y no en la cadena: así las dos lecturas de §11 («municipios de las sucursales») dan lo mismo hasta que el integrador fije el texto (REVISION-S5-r1 H-01). Lectura 10 = 3 + 2 + 2 + 2 + 1 (producción) = 3 + 2 + 1 + 3 + 1 (revisión). '
            'Los motivos (`reason`) son texto libre y no se comparan.',
            [('criterio', 'S5-05'), ('manual', '§11 de S5-interfaz.md')],
            {'today': HOY, 'catalog_file': CAT, 'stores': stores},
            [consulta('produccion', 'profeco_import', **args('production')), consulta('revision', 'profeco_import', **args('review'))], e,
            prov('S5-profeco-almacen', 'almacenamiento, conteos y cobertura (opción A); ligas confirmadas: bloqueo S4-confirmacion-empates'))


def ocr():
    cab = 'level\tpage_num\tblock_num\tpar_num\tline_num\tword_num\tleft\ttop\twidth\theight\tconf\ttext'
    filas = [cab, '1\t1\t0\t0\t0\t0\t0\t0\t600\t800\t-1\t', '2\t1\t1\t0\t0\t0\t10\t10\t500\t300\t-1\t', '4\t1\t1\t1\t1\t0\t10\t10\t500\t30\t-1\t']

    def linea(n, palabras):
        for i, (t, c) in enumerate(palabras, 1):
            filas.append(f'5\t1\t1\t1\t{n}\t{i}\t{10 + 50 * i}\t{10 + 40 * n}\t45\t30\t{c}\t{t}')
    linea(1, [('2', '96'), ('LECHE', '95'), ('ENTERA', '93'), ('1L', '90'), ('48.00', '97')])
    linea(2, [('0.750', '91'), ('KG', '90'), ('MANZANA', '88'), ('31.50', '92')])
    linea(3, [('PAN', '80'), ('BLANCO', '81'), ('$45.00', '82')])
    linea(4, [('BOLSA', '60.4')])
    linea(5, [('TOTAL', '99'), ('1,234.50', '98')])
    linea(6, [('3', '90'), ('GALLETAS', '91'), ('MARIAS', '91'), ('12.00', '90')])
    filas.append('5\t1\t1\t1\t7\t1\t10\t300\t45\t30\t-1\t ')
    (DATOS / 'ocr-muestra.tsv').write_text('\n'.join(filas) + '\n', encoding='utf-8')
    L = lambda d_, q, a, c: {'description': d_, 'quantity_milli': q, 'amount_cents': a, 'confidence': c}
    e = {'ticket': {'lines': [L('LECHE ENTERA 1L', 2000, 4800, 94), L('MANZANA', 750, 3150, 90), L('PAN BLANCO', None, 4500, 81), L('BOLSA', None, None, 60),
                              L('TOTAL', None, 123450, 99), L('GALLETAS MARIAS', 3000, 1200, 91)]}}
    fixture('s5-ocr-parseo',
            'Parseo del TSV de Tesseract (ILUSTRATIVO; §12, S5-06; no corre Tesseract): palabras de nivel 5 agrupadas por bloque, párrafo y línea; importe = último número con dos decimales '
            '(48.00, $45.00, 1,234.50); cantidad = número inicial (2 -> 2000; 0.750 KG -> 750, sin la unidad); confianza = promedio HALF_UP (94.2 -> 94; 90.25 -> 90; 60.4 -> 60; '
            '362 / 4 = 90.5 -> 91). Una línea sin importe sale con importe nulo; una palabra vacía se ignora.',
            [('criterio', 'S5-06'), ('manual', '§12 de S5-interfaz.md')],
            {'today': HOY, 'catalog': catalogo([('a', 'g', 'abarrotes', [('a-1kg', 1000000)])])},
            [consulta('ticket', 'ocr_parse_tsv', tsv_file='oraculo/fixtures/s5/_datos/ocr-muestra.tsv')], e)


def profeco_filtro():
    """Bloqueo S5-profeco-filtro, opción A (PROVISIONAL): municipios solo de sucursales, comparación normalizada y
    sucursal buscada dentro de la cadena que casó (REVISION-S5-r2 R2-02). Esperado escrito a mano renglón por renglón."""
    cab = 'producto,presentacion,marca,categoria,catalogo,precio,fecha_registro,cadena_comercial,giro,nombre_comercial,direccion,estado,municipio,latitud,longitud'

    def fila(fecha_, cadena, nombre, muni):
        return f'ATÚN DOLORES,140 G,DOLORES,ABARROTES,CATALOGO,21.00,{fecha_},{cadena},SUPERMERCADO,{nombre},"CALLE 3",CIUDAD DE MEXICO,{muni},19.45,-99.15'
    filas = [fila('2026-08-01', 'WALMART', 'WALMART ACUEDUCTO', 'GUSTAVO A. MADERO'),         # f1 carga (cadena Wal-mart = WALMART; municipio de walmart-tepeyac); tienda walmart
             fila('2026-08-02', 'WAL-MART', 'WAL-MART TEPEYAC', 'GUSTAVO A. MADERO'),         # f2 carga; tienda walmart-tepeyac («Walmart Tepeyac» = «WAL-MART TEPEYAC»)
             fila('2026-08-03', 'BODEGA AURRERA', 'BODEGA AURRERA TLALNEPANTLA', 'TLALNEPANTLA DE BAZ'),  # f3 filtrado: ese municipio solo lo tiene la cadena
             fila('2026-08-04', 'WALMART', 'WALMART IZTAPALAPA', 'IZTAPALAPA'),               # f4 filtrado: ese municipio solo lo tiene el mercado
             fila('2026-08-05', 'SORIANA', 'SORIANA GAM', 'GUSTAVO A. MADERO'),               # f5 filtrado: no es cadena del catálogo
             fila('2026-08-06', 'WALMART', 'BODEGA AURRERA VALLEJO', 'AZCAPOTZALCO'),         # f6 carga; la sucursal es de otra cadena -> tienda walmart
             fila('2026-08-07', 'Bodega Aurrerá', 'Bodega Aurrera Vallejo', 'Azcapotzalco'),  # f7 carga; tienda bodega-vallejo (acentos y mayúsculas no cuentan)
             fila('2026-08-06', 'WALMART', 'WALMART AZCAPOTZALCO', 'AZCAPOTZALCO')]           # f8 duplicado de f6 (misma presentación, tienda walmart y día)
    (DATOS / 'profeco-filtro.csv').write_bytes((cab + '\n' + '\n'.join(filas) + '\n').encode('utf-8'))
    stores = [tienda('bodega', profeco_name='Bodega Aurrera', municipality='TLALNEPANTLA DE BAZ'),
              tienda('bodega-vallejo', kind='branch', parent='bodega', profeco_name='BODEGA AURRERA VALLEJO', municipality='Azcapotzalco'),
              tienda('walmart', profeco_name='Wal-mart'),
              tienda('walmart-tepeyac', kind='branch', parent='walmart', profeco_name='Walmart Tepeyac', municipality='Gustavo A. Madero'),
              tienda('tianguis', kind='market', name='Tianguis', municipality='IZTAPALAPA', visit_weekday=2)]
    d = 'oraculo/fixtures/s5/_datos/'
    e = {'read': 8, 'rejected': 0, 'filtered_out': 3, 'skipped': 0, 'loaded': 4, 'duplicates': 1,
         'coverage': [{'chain': 'Bodega Aurrerá', 'municipality': 'Azcapotzalco', 'price_count': 1},
                      {'chain': 'WAL-MART', 'municipality': 'GUSTAVO A. MADERO', 'price_count': 1},
                      {'chain': 'WALMART', 'municipality': 'AZCAPOTZALCO', 'price_count': 1},
                      {'chain': 'WALMART', 'municipality': 'GUSTAVO A. MADERO', 'price_count': 1}],
         'rejections': []}
    assert e['read'] == e['rejected'] + e['filtered_out'] + e['skipped'] + e['loaded'] + e['duplicates']
    assert e['loaded'] == sum(c['price_count'] for c in e['coverage'])
    fixture('s5-profeco-filtro',
            'Filtro de profeco_import (ILUSTRATIVO; §11, REVISION-S5-r2 R2-02) con la opción A del bloqueo S5-profeco-filtro: (1) el municipio de una cadena (Bodega en TLALNEPANTLA DE BAZ) '
            'o de un mercado (IZTAPALAPA) no entra al filtro, solo el de las sucursales; (2) cadena, municipio y nombre comercial se comparan normalizados («Wal-mart» = «WALMART» = «WAL-MART», '
            '«Bodega Aurrerá» = «BODEGA AURRERA», «Gustavo A. Madero» = «GUSTAVO A. MADERO»); (3) la sucursal del precio se busca solo dentro de la cadena que casó, así que f6 (WALMART con nombre '
            'BODEGA AURRERA VALLEJO) va a la cadena walmart y f8 es su duplicado. coverage[] agrupa por el texto del CSV tal como llega. Lectura 8 = 3 filtrados + 4 cargados + 1 duplicado.',
            [('bloqueo', BLOQ + 'S5-profeco-filtro.md'), ('manual', '§11 de S5-interfaz.md')],
            {'today': HOY, 'catalog_file': CAT, 'stores': stores},
            [consulta('filtro', 'profeco_import', files=[d + 'profeco-filtro.csv'], links_file=d + 'profeco-vinculos.csv', mode='production')],
            {'filtro': e},
            prov('S5-profeco-filtro', 'municipios solo de sucursales, comparación normalizada y sucursal dentro de su cadena (opción A)'))


def carritos():
    """Bloqueo S5-carritos-detalle, opción A (PROVISIONAL; REVISION-S5-r2 R2-03): sustituto con varias presentaciones en una tienda,
    missing_count y stale_count de compare_carts. Esperado con aritmética de enteros y Fraction escrita aquí."""
    sus = REP['sustitutos']
    assert [NOMBRE[x['bodega']][0] for x in sus] == ['arroz-grano-largo', 'sustituto-crema-cafe', 'leche']
    atun = 'atun-agua.dolores-140g'
    assert any(v[1] == atun for v in NOMBRE.values())
    tam = {v[1]: v[2] for v in NOMBRE.values()}
    lin = lambda lado, c: [{'presentation': NOMBRE[x[lado]][1], 'quantity_milli': x[c] * 1000, 'unit_price_cents': cts(x[lado + '_total']) // x[c]} for x in sus]
    a750, a1k = NOMBRE[sus[0]['bodega']][1], NOMBRE[sus[0]['walmart']][1]
    ln = lambda pres, n, u: {'presentation': pres, 'quantity_milli': n * 1000, 'unit_price_cents': u}
    b_lines = lin('bodega', 'cant_b') + [ln(atun, 6, 2100)]
    w_lines = lin('walmart', 'cant_w')[:2] + [ln(atun, 5, 2150)]
    t_lines = [ln(a750, 1, 1500), ln(a1k, 1, 5000), ln(atun, 5, 2000)]
    carts = [{'code': 'bodega', 'store': 'bodega', 'kind': 'expense', 'observed_date': '2026-09-01', 'source': 'manual', 'lines': b_lines},
             {'code': 'walmart', 'store': 'walmart', 'kind': 'cart', 'observed_date': '2026-10-10', 'source': 'manual', 'lines': w_lines},
             {'code': 'tercera', 'store': 'tercera', 'kind': 'cart', 'observed_date': '2026-10-10', 'source': 'manual', 'lines': t_lines}]
    tot = lambda ls: sum(x['quantity_milli'] * x['unit_price_cents'] for x in ls) // 1000
    assert (tot(b_lines), tot(w_lines), tot(t_lines)) == (42600, 29650, 16500)
    # arroz por unidad base, opción A (importe total ÷ contenido total): Bodega 6600/2.25 kg, Walmart 8400/3 kg, tercera 6500/1.75 kg
    pu = lambda ls: Fraction(sum(x['quantity_milli'] * x['unit_price_cents'] for x in ls), sum(x['quantity_milli'] * tam[x['presentation']] for x in ls))
    arroz = {'bodega': pu(b_lines[:1]), 'walmart': pu(w_lines[:1]), 'tercera': pu(t_lines[:2])}
    assert min(arroz, key=arroz.get) == 'walmart' and hu(arroz['tercera'] * 1000000) == 3714
    assert Fraction(1500 * 1000000, tam[a750]) == 2000  # la mejor presentación de la tercera (opción B) habría ganado con 20.00/kg
    st = lambda s, t_, ic, it, sc, pc, mc, stc, dif: {'store': s, 'total_cents': t_, 'identical_count': ic, 'identical_total_cents': it, 'substitute_count': sc,
                                                     'piece_count': pc, 'missing_count': mc, 'stale_count': stc, 'difference_vs_cheapest_cents': dif, 'assumes_membership': False}
    esperado = {'carritos': {
        'stores': [st('tercera', 16500, 1, 10000, 1, 7, 2, 0, 0), st('walmart', 29650, 1, 10750, 2, 9, 1, 0, 29650 - 16500),
                   st('bodega', 42600, 1, 10500, 2, 11, 0, 4, 42600 - 16500)],
        'identical_pieces': 5, 'winners': [{'store': 'tercera', 'count': 1}, {'store': 'walmart', 'count': 0}, {'store': 'bodega', 'count': 0}], 'ties': 0,
        'rows': [{'kind': 'identical', 'product': 'atun-agua', 'winner_store': 'tercera'},
                 {'kind': 'substitute', 'product': 'arroz-grano-largo', 'winner_store': 'walmart'},
                 {'kind': 'substitute', 'product': 'sustituto-crema-cafe', 'winner_store': 'bodega'},
                 {'kind': 'only_one', 'product': 'leche', 'winner_store': None},
                 {'kind': 'only_one', 'product': 'atun-agua', 'winner_store': None}]}}
    assert NOMBRE_PROD[atun] == 'atun-agua'
    fixture('s5-carritos-presentaciones',
            'compare_carts con tres canastas (ILUSTRATIVO; §7.2, REVISION-S5-r2 R2-03) y la opción A del bloqueo S5-carritos-detalle. Bodega (base, del 1-sep, vieja) = los 3 sustitutos de data.json + 6 atunes a 21.00; '
            'Walmart = arroz y crema de data.json + 5 atunes a 21.50 (sin leche); tercera = arroz 750 g a 15.00 + arroz 1 kg a 50.00 + 5 atunes a 20.00. (1) El sustituto del arroz lo gana Walmart '
            '(28.00/kg) porque la tercera cobró 65.00 por 1.75 kg (37.14/kg), aunque su 750 g sola saldría a 20.00/kg. (2) missing_count: tercera 2 (crema y leche), Walmart 1 (leche), Bodega 0. '
            '(3) stale_count de Bodega = 4 presentaciones agrupadas (el atún cuenta una vez aunque salga como idéntico y como «solo aquí»). winners[] va en el orden de stores[].',
            [('bloqueo', BLOQ + 'S5-carritos-detalle.md'), ('data.json', 'sustitutos'), ('manual', '§7.2 de S5-interfaz.md')],
            {'today': HOY, 'catalog_file': CAT, 'stores': TIENDAS_CANON + [tienda('tercera', name='Tercera tienda')], 'carts': carts},
            [consulta('carritos', 'compare_carts', base='bodega', others=['walmart', 'tercera'])], esperado,
            prov('S5-carritos-detalle', 'ganador del sustituto con varias presentaciones, missing_count y stale_count de carritos (opción A)'))


if __name__ == '__main__':
    canon()
    unidad()
    frescura()
    fuentes()
    canasta()
    empates()
    cotizacion()
    membresia()
    mercado()
    viejos()
    asistida()
    ingesta()
    profeco()
    profeco_filtro()
    ocr()
    carritos()
    print(f'{len(list(OUT.glob("*.json")))} fixtures en {OUT}')
