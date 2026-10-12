"""Oráculo S5 · carga y validación del `Input` (S5-interfaz.md §15.1, §3, §6).

Claves desconocidas -> unknown_field; ausentes -> missing_field. Rutas `*_file` relativas a la raíz del repositorio.
Todo se lee de data.json y productos.json con json (Decimal para los importes), sin flotantes.
"""
import json
from decimal import Decimal
from pathlib import Path

from pre_comun import (ErrorMotor, booleano, entero, enum, fecha, instante, lista, objeto, texto)

RAIZ = Path(__file__).resolve().parents[1]
FUENTES = ('ticket_manual', 'ticket_ocr', 'manual', 'site_assisted', 'profeco')
COMPITEN = ('ticket_manual', 'ticket_ocr', 'manual')
UNIDADES = ('g', 'ml', 'piece')
DATOS_S5 = RAIZ / 'oraculo' / 'fixtures' / 's5' / '_datos'


class Datos:
    pass


def ruta_repo(p, ruta, dentro=None):
    """Ruta relativa a la raíz del repositorio. `dentro`: carpeta que debe contenerla (sin `..`, no absoluta)."""
    texto(p, ruta)
    q = Path(p)
    if q.is_absolute() or '..' in q.parts:
        raise ErrorMotor('invalid_argument', ruta)
    r = RAIZ / q
    if dentro is not None and dentro not in r.parents:
        raise ErrorMotor('invalid_argument', ruta)
    return r


def centavos(v, ruta):
    """Decimal de data.json -> centavos enteros, exactos."""
    c = Decimal(v) * 100
    if c != c.to_integral_value():
        raise ErrorMotor('invalid_argument', ruta)
    return int(c)


def _milli_txt(txt, ruta):
    c = Decimal(txt) * 1000
    if c != c.to_integral_value() or c < 0:
        raise ErrorMotor('invalid_argument', ruta)
    return int(c)


def _catalogo_de_archivo(ruta):
    with open(ruta, encoding='utf-8') as fh:
        d = json.load(fh)
    grupos = [{'code': g['id'], 'name': g['nombre'], 'purchase_channel': {'mercado': 'market', 'cadena': 'chain'}[g['canal']]}
              for g in d['grupos']]
    unidad = {'g': 'g', 'ml': 'ml', 'pieza': 'piece'}
    prods, pres = [], []
    for p in d['productos']:
        gpu = p.get('gramos_por_unidad')
        prods.append({'code': p['id'], 'name': p['nombre'], 'base_unit': unidad[p['unidad_base']], 'group': p['categoria'],
                      'grams_per_unit_milli': None if gpu is None else _milli_txt(gpu, p['id'])})
        for x in p['presentaciones']:
            pres.append({'code': x['id'], 'product': p['id'], 'name': x['nombre'], 'size_milli': x['contenido_milli']})
    return {'groups': grupos, 'products': prods, 'presentations': pres}


def _catalogo(raw, d):
    g, f = raw.get('catalog'), raw.get('catalog_file')
    if (g is None) == (f is None):
        raise ErrorMotor('missing_field' if g is None else 'invalid_argument', 'catalog/catalog_file: uno de los dos')
    cat = g if g is not None else _catalogo_de_archivo(ruta_repo(f, 'catalog_file'))
    objeto(cat, 'catalog', ('groups', 'products', 'presentations'))
    d.groups, d.products, d.pres = {}, {}, {}
    d.pres_names = {}
    for i, x in enumerate(lista(cat['groups'], 'catalog.groups')):
        r = f'catalog.groups[{i}]'
        objeto(x, r, ('code', 'name', 'purchase_channel'))
        c = texto(x['code'], r + '.code')
        if c in d.groups:
            raise ErrorMotor('invalid_argument', r + '.code repetido')
        d.groups[c] = enum(x['purchase_channel'], r + '.purchase_channel', ('market', 'chain'))
    for i, x in enumerate(lista(cat['products'], 'catalog.products')):
        r = f'catalog.products[{i}]'
        objeto(x, r, ('code', 'name', 'base_unit', 'group', 'grams_per_unit_milli'))
        c = texto(x['code'], r + '.code')
        if c in d.products:
            raise ErrorMotor('invalid_argument', r + '.code repetido')
        if x['group'] is not None and x['group'] not in d.groups:
            raise ErrorMotor('invalid_argument', r + '.group')
        d.products[c] = {'unit': enum(x['base_unit'], r + '.base_unit', UNIDADES), 'group': x['group'], 'idx': i, 'pres': []}
    for i, x in enumerate(lista(cat['presentations'], 'catalog.presentations')):
        r = f'catalog.presentations[{i}]'
        objeto(x, r, ('code', 'product', 'name', 'size_milli'))
        c = texto(x['code'], r + '.code')
        if c in d.pres:
            raise ErrorMotor('invalid_argument', r + '.code repetido')
        if x['product'] not in d.products:
            raise ErrorMotor('invalid_argument', r + '.product')
        d.pres[c] = {'code': c, 'product': x['product'], 'size': entero(x['size_milli'], r + '.size_milli', 1), 'idx': i}
        d.pres_names[x['name']] = c
        d.products[x['product']]['pres'].append(c)


def _tiendas(raw, d):
    d.stores = {}
    for i, x in enumerate(lista(raw.get('stores', []), 'stores')):
        r = f'stores[{i}]'
        objeto(x, r, ('code', 'kind', 'name'), ('parent', 'profeco_name', 'municipality', 'visit_weekday', 'requires_membership',
                                                  'has_membership', 'membership_fee_cents', 'assisted_reading', 'groups'))
        c = texto(x['code'], r + '.code')
        if c in d.stores:
            raise ErrorMotor('invalid_argument', r + '.code repetido')
        groups = lista(x.get('groups', []), r + '.groups')
        for gcode in groups:
            if gcode not in d.groups:
                raise ErrorMotor('invalid_argument', r + '.groups')
        wd = entero(x.get('visit_weekday'), r + '.visit_weekday', 1, nulo=True)
        if wd is not None and wd > 7:
            raise ErrorMotor('invalid_argument', r + '.visit_weekday')
        d.stores[c] = {
            'code': c, 'idx': i, 'kind': enum(x['kind'], r + '.kind', ('chain', 'branch', 'market')),
            'parent': texto(x.get('parent'), r + '.parent', nulo=True),  # cadena de una sucursal (filtro PROFECO, bloqueo S5-profeco-filtro)
            'profeco_name': texto(x.get('profeco_name'), r + '.profeco_name', nulo=True),
            'municipality': texto(x.get('municipality'), r + '.municipality', nulo=True), 'visit_weekday': wd,
            'requires': booleano(x.get('requires_membership', False), r + '.requires_membership'),
            'has': booleano(x.get('has_membership', False), r + '.has_membership'),
            'fee': entero(x.get('membership_fee_cents'), r + '.membership_fee_cents', 0, nulo=True),
            'assisted': booleano(x.get('assisted_reading', False), r + '.assisted_reading'), 'groups': list(groups)}


def _precios(raw, d):
    d.prices, vistos = [], set()
    for i, x in enumerate(lista(raw.get('prices', []), 'prices')):
        r = f'prices[{i}]'
        objeto(x, r, ('presentation', 'store', 'price_cents', 'observed_date', 'source', 'is_promotional', 'created_at'), ('source_ref',))
        if x['presentation'] not in d.pres:
            raise ErrorMotor('unknown_presentation', r)
        if x['store'] not in d.stores:
            raise ErrorMotor('unknown_store', r)
        p = {'pres': x['presentation'], 'store': x['store'], 'cents': entero(x['price_cents'], r + '.price_cents', 0),
             'date': fecha(x['observed_date'], r + '.observed_date'), 'source': enum(x['source'], r + '.source', FUENTES),
             'promo': booleano(x['is_promotional'], r + '.is_promotional'), 'created': instante(x['created_at'], r + '.created_at'),
             'ref': x.get('source_ref'), 'n': i}
        if p['date'] > d.today:
            raise ErrorMotor('future_date', r)
        clave = (p['pres'], p['store'], p['date'], p['source'], p['promo'])
        if clave in vistos:
            raise ErrorMotor('duplicate_price', r)
        vistos.add(clave)
        d.prices.append(p)


def _linea_carrito(x, r, d):
    objeto(x, r, ('presentation', 'quantity_milli', 'unit_price_cents'))
    if x['presentation'] not in d.pres:
        raise ErrorMotor('unknown_presentation', r)
    return {'pres': x['presentation'], 'qty': entero(x['quantity_milli'], r + '.quantity_milli', 1),
            'unit': entero(x['unit_price_cents'], r + '.unit_price_cents', 0)}


def _carritos(raw, d):
    d.carts = {}
    for i, x in enumerate(lista(raw.get('carts', []), 'carts')):
        r = f'carts[{i}]'
        objeto(x, r, ('code', 'store', 'kind', 'observed_date', 'source', 'lines'))
        c = texto(x['code'], r + '.code')
        if c in d.carts:
            raise ErrorMotor('invalid_argument', r + '.code repetido')
        if x['store'] not in d.stores:
            raise ErrorMotor('unknown_store', r)
        d.carts[c] = {'code': c, 'store': x['store'], 'kind': enum(x['kind'], r + '.kind', ('expense', 'cart')),
                      'date': fecha(x['observed_date'], r + '.observed_date'), 'source': enum(x['source'], r + '.source', FUENTES),
                      'lines': [_linea_carrito(y, f'{r}.lines[{j}]', d) for j, y in enumerate(lista(x['lines'], r + '.lines'))]}


def _report(ruta, d):
    """Arma los carritos `bodega` (expense) y `walmart` (cart) del 2026-10-10 desde data.json (S5-interfaz §15.1, §3)."""
    with open(ruta, encoding='utf-8') as fh:
        rep = json.load(fh, parse_float=Decimal)
    for t in ('bodega', 'walmart'):
        if t not in d.stores:
            raise ErrorMotor('unknown_store', f'report_file necesita la tienda {t}')
    fch = fecha(rep['fecha'], 'report_file.fecha')

    def pres(nombre):
        if nombre not in d.pres_names:
            raise ErrorMotor('unknown_presentation', nombre)
        return d.pres_names[nombre]
    lb, lw = [], []
    atun = None
    for i, it in enumerate(rep['items_identicos']):
        p = pres(it['producto'])
        q = it['cantidad'] * 1000
        lb.append({'pres': p, 'qty': q, 'unit': centavos(it['bodega_unit'], f'items_identicos[{i}].bodega_unit')})
        lw.append({'pres': p, 'qty': q, 'unit': centavos(it['walmart_unit'], f'items_identicos[{i}].walmart_unit')})
        if it['producto'].startswith('Atún'):
            atun = p
    lb.append({'pres': atun, 'qty': 1000, 'unit': centavos(rep['atun_sexta_pieza_bodega'], 'atun_sexta_pieza_bodega')})
    for i, s in enumerate(rep['sustitutos']):
        lb.append({'pres': pres(s['bodega']), 'qty': s['cant_b'] * 1000, 'unit': centavos(Decimal(s['bodega_total']) / s['cant_b'], f'sustitutos[{i}]')})
        lw.append({'pres': pres(s['walmart']), 'qty': s['cant_w'] * 1000, 'unit': centavos(Decimal(s['walmart_total']) / s['cant_w'], f'sustitutos[{i}]')})
    d.carts['bodega'] = {'code': 'bodega', 'store': 'bodega', 'kind': 'expense', 'date': fch, 'source': 'manual', 'lines': lb}
    d.carts['walmart'] = {'code': 'walmart', 'store': 'walmart', 'kind': 'cart', 'date': fch, 'source': 'manual', 'lines': lw}


def _lista_compras(raw, d):
    d.shopping = []
    for i, x in enumerate(lista(raw.get('shopping_list', []), 'shopping_list')):
        r = f'shopping_list[{i}]'
        objeto(x, r, ('product', 'needed_milli', 'quantity_milli', 'recommended_presentation'))
        if x['product'] not in d.products:
            raise ErrorMotor('invalid_argument', r + '.product')
        nec = entero(x['needed_milli'], r + '.needed_milli', 1, nulo=True)
        cant = entero(x['quantity_milli'], r + '.quantity_milli', 1, nulo=True)
        if nec is None and cant is None:
            raise ErrorMotor('invalid_argument', r)
        rec = x['recommended_presentation']
        if rec is not None and rec not in d.pres:
            raise ErrorMotor('unknown_presentation', r)
        d.shopping.append({'product': x['product'], 'qty': cant if cant is not None else nec, 'rec': rec})


def cargar(raw):
    d = Datos()
    objeto(raw, 'datos', ('today',), ('settings', 'catalog_file', 'catalog', 'stores', 'prices', 'carts', 'report_file',
                                      'shopping_list', 'grocery_purchases', 'terms'))
    d.today = fecha(raw['today'], 'today')
    s = objeto(raw.get('settings', {}), 'settings', (), ('stale_after_days', 'zero_weight_after_days', 'assisted_reading_enabled'))
    d.dv = entero(s.get('stale_after_days', 30), 'settings.stale_after_days', 1)
    d.d0 = entero(s.get('zero_weight_after_days', 90), 'settings.zero_weight_after_days', 1)
    if d.d0 < d.dv:
        raise ErrorMotor('invalid_argument', 'settings')
    d.assisted_on = booleano(s.get('assisted_reading_enabled', False), 'settings.assisted_reading_enabled')
    _catalogo(raw, d)
    _tiendas(raw, d)
    _precios(raw, d)
    _carritos(raw, d)
    if 'report_file' in raw:
        if raw.get('carts'):
            raise ErrorMotor('invalid_argument', 'report_file no se combina con carts')
        _report(ruta_repo(raw['report_file'], 'report_file'), d)
    _lista_compras(raw, d)
    d.grocery = []
    for i, x in enumerate(lista(raw.get('grocery_purchases', []), 'grocery_purchases')):
        r = f'grocery_purchases[{i}]'
        objeto(x, r, ('date', 'store'))
        if x['store'] not in d.stores:
            raise ErrorMotor('unknown_store', r)
        d.grocery.append({'date': fecha(x['date'], r + '.date'), 'store': x['store']})
    d.terms = []
    for i, x in enumerate(lista(raw.get('terms', []), 'terms')):
        r = f'terms[{i}]'
        objeto(x, r, ('site', 'verdict', 'url', 'read_date'))
        d.terms.append({'site': enum(x['site'], r + '.site', ('walmart', 'bodega_aurrera', 'sams_club')),
                        'verdict': x['verdict'], 'url': x['url'], 'read_date': x['read_date']})
    return d
