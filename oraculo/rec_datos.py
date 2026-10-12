"""Oráculo S4 · lectura y validación del `Input` (S4-interfaz.md §12.1) y de las semillas (§3, §9.1).

Todo se cita por código. Un campo desconocido es unknown_field; uno obligatorio ausente, missing_field.
"""
import csv
import json
from decimal import Decimal, InvalidOperation
from pathlib import Path

from rec_comun import (ErrorMotor, booleano, entero, enum, fecha, lista, objeto, texto)

RAIZ = Path(__file__).resolve().parents[1]
SLOTS = ('breakfast', 'lunch', 'dinner', 'snack')
UNIDADES = ('g', 'ml', 'piece')
FUENTES = ('usda_fdc', 'open_food_facts', 'manual', 'mx_table', 'store_label')
ACTIVIDADES = ('sedentary', 'light', 'moderate', 'active', 'very_active')


class Datos:
    pass


def _milli(txt, ruta, tope_dec=3):
    """Decimal en texto -> milésimas enteras, sin flotantes (hasta tope_dec decimales)."""
    try:
        d = Decimal(txt)
    except (InvalidOperation, TypeError):
        raise ErrorMotor('invalid_argument', ruta)
    m = d * 1000
    if m != m.to_integral_value() or d < 0:
        raise ErrorMotor('invalid_argument', ruta)
    return int(m)


def _ruta_repo(p, ruta):
    texto(p, ruta)
    q = Path(p)
    return q if q.is_absolute() else RAIZ / q


# ---- catálogo ----
def _catalogo_de_archivo(ruta):
    with open(ruta, encoding='utf-8') as fh:
        d = json.load(fh)
    grupos = [{'code': g['id'], 'name': g['nombre'], 'purchase_channel': {'mercado': 'market', 'cadena': 'chain'}[g['canal']]}
              for g in d['grupos']]
    unidad = {'g': 'g', 'ml': 'ml', 'pieza': 'piece'}
    productos, pres = [], []
    for p in d['productos']:
        gpu = p.get('gramos_por_unidad')
        productos.append({'code': p['id'], 'name': p['nombre'], 'base_unit': unidad[p['unidad_base']], 'group': p['categoria'],
                          'grams_per_unit_milli': None if gpu is None else _milli(gpu, f"productos.{p['id']}.gramos_por_unidad")})
        for x in p['presentaciones']:
            pres.append({'code': x['id'], 'product': p['id'], 'name': x['nombre'], 'size_milli': x['contenido_milli']})
    return {'groups': grupos, 'products': productos, 'presentations': pres}


def _catalogo(raw, d):
    g = raw.get('catalog')
    f = raw.get('catalog_file')
    if (g is None) == (f is None):
        raise ErrorMotor('missing_field' if g is None else 'invalid_argument', 'catalog/catalog_file: uno de los dos')
    cat = g if g is not None else _catalogo_de_archivo(_ruta_repo(f, 'catalog_file'))
    objeto(cat, 'catalog', ('groups', 'products', 'presentations'))
    d.groups, d.products, d.pres = {}, {}, {}
    for i, x in enumerate(lista(cat['groups'], 'catalog.groups')):
        r = f'catalog.groups[{i}]'
        objeto(x, r, ('code', 'name', 'purchase_channel'))
        c = texto(x['code'], r + '.code')
        if c in d.groups:
            raise ErrorMotor('invalid_argument', r + '.code repetido')
        d.groups[c] = {'channel': enum(x['purchase_channel'], r + '.purchase_channel', ('market', 'chain'))}
    for i, x in enumerate(lista(cat['products'], 'catalog.products')):
        r = f'catalog.products[{i}]'
        objeto(x, r, ('code', 'name', 'base_unit', 'group', 'grams_per_unit_milli'))
        c = texto(x['code'], r + '.code')
        if c in d.products:
            raise ErrorMotor('invalid_argument', r + '.code repetido')
        gr = x['group']
        if gr is not None and gr not in d.groups:
            raise ErrorMotor('invalid_argument', r + '.group')
        d.products[c] = {'unit': enum(x['base_unit'], r + '.base_unit', UNIDADES), 'group': gr,
                         'gpum': entero(x['grams_per_unit_milli'], r + '.grams_per_unit_milli', 1, nulo=True)}
    for i, x in enumerate(lista(cat['presentations'], 'catalog.presentations')):
        r = f'catalog.presentations[{i}]'
        objeto(x, r, ('code', 'product', 'name', 'size_milli'))
        c = texto(x['code'], r + '.code')
        if c in d.pres:
            raise ErrorMotor('invalid_argument', r + '.code repetido')
        if x['product'] not in d.products:
            raise ErrorMotor('unknown_product', r + '.product')
        d.pres[c] = {'product': x['product'], 'size': entero(x['size_milli'], r + '.size_milli', 1)}


# ---- recetas ----
def _recetario_de_archivo(ruta):
    with open(ruta, encoding='utf-8') as fh:
        d = json.load(fh)
    unidad = {'g': 'g', 'ml': 'ml', 'pieza': 'piece'}
    res = []
    for r in d['recetas']:
        res.append({'code': r['id'], 'name': r['nombre'], 'slots': r['tiempos'], 'base_servings': r['porciones_base'],
                    'prep_minutes': r['minutos_preparacion'], 'cook_minutes': r['minutos_coccion'],
                    'origin': {'claude': 'claude_seed', 'propia': 'household'}[r['origen']], 'is_active': True,
                    'ingredients': [{'product': i['producto_id'], 'presentation': i['presentacion_id'], 'quantity_milli': i['cantidad_milli'],
                                     'optional': i['opcional'], '_unidad': unidad[i['unidad']]} for i in r['ingredientes']]})
    return res


def _recetas(raw, d):
    rl = raw.get('recipes')
    rf = raw.get('recipe_book_file')
    if rl is None and rf is None:
        raise ErrorMotor('missing_field', 'recipes/recipe_book_file: al menos uno')
    todas = []
    if rf is not None:
        todas += _recetario_de_archivo(_ruta_repo(rf, 'recipe_book_file'))
    if rl is not None:
        todas += lista(rl, 'recipes')
    d.recipes, d.recipe_order = {}, []
    for i, x in enumerate(todas):
        r = f'recipes[{i}]'
        objeto(x, r, ('code', 'name', 'slots', 'base_servings', 'prep_minutes', 'cook_minutes', 'origin', 'is_active', 'ingredients'))
        c = texto(x['code'], r + '.code')
        if c in d.recipes:
            raise ErrorMotor('invalid_argument', r + '.code repetido')
        ings, vistos = [], set()
        for j, y in enumerate(lista(x['ingredients'], r + '.ingredients')):
            ri = f'{r}.ingredients[{j}]'
            objeto(y, ri, ('product', 'presentation', 'quantity_milli', 'optional'), ('_unidad',))
            p = y['product']
            if p not in d.products:
                raise ErrorMotor('unknown_product', ri + '.product')
            if p in vistos:
                raise ErrorMotor('invalid_argument', ri + '.product repetido')
            vistos.add(p)
            pr = y['presentation']
            if pr not in d.pres or d.pres[pr]['product'] != p:
                raise ErrorMotor('unknown_presentation', ri + '.presentation')
            if y.get('_unidad') is not None and y['_unidad'] != d.products[p]['unit']:
                raise ErrorMotor('unit_mismatch', ri)
            ings.append({'product': p, 'presentation': pr, 'q': entero(y['quantity_milli'], ri + '.quantity_milli', 1),
                         'optional': booleano(y['optional'], ri + '.optional')})
        d.recipes[c] = {'name': texto(x['name'], r + '.name'),
                        'slots': [enum(s, r + '.slots', SLOTS) for s in lista(x['slots'], r + '.slots')],
                        'base': entero(x['base_servings'], r + '.base_servings', 1, 50),
                        'origin': enum(x['origin'], r + '.origin', ('claude_seed', 'household')),
                        'active': booleano(x['is_active'], r + '.is_active'), 'ings': ings}
        d.recipe_order.append(c)


# ---- nutrición ----
_ENC = 'producto_id,kcal_100g,proteina_g,grasa_g,carbohidratos_g,fuente,id_fuente,nombre_fuente,fecha,estado,confirmo,fecha_confirmacion'


def _nutricion_de_csv(ruta, d):
    with open(ruta, encoding='utf-8', newline='') as fh:
        filas = list(csv.reader(fh))
    if not filas or ','.join(filas[0]) != _ENC:
        raise ErrorMotor('invalid_argument', 'nutrition_file: encabezado')
    fuente = {'USDA': 'usda_fdc', 'OFF': 'open_food_facts', 'manual': 'manual'}
    res = []
    for n, f in enumerate(filas[1:], 2):
        if len(f) != 12:
            raise ErrorMotor('invalid_argument', f'nutrition_file:{n}')
        p, kcal, _, _, _, fu, idf, _, _, est, conf, fconf = f
        if p not in d.products or fu not in fuente or not idf or est not in ('propuesto', 'confirmado'):
            raise ErrorMotor('invalid_argument', f'nutrition_file:{n}')
        if (est == 'confirmado') != bool(conf and fconf):
            raise ErrorMotor('invalid_argument', f'nutrition_file:{n}')
        res.append({'product': p, 'presentation': None, 'kcal_milli_per_100g': _milli(kcal, f'nutrition_file:{n}'),
                    'source': fuente[fu], 'source_ref': idf, 'status': 'proposed' if est == 'propuesto' else 'confirmed'})
    return res


def _nutricion(raw, d):
    nl, nf = raw.get('nutrition'), raw.get('nutrition_file')
    if nl is not None and nf is not None:
        raise ErrorMotor('invalid_argument', 'nutrition/nutrition_file: uno u otro')
    filas = _nutricion_de_csv(_ruta_repo(nf, 'nutrition_file'), d) if nf is not None else (lista(nl, 'nutrition') if nl is not None else [])
    d.nutrition = {}
    for i, x in enumerate(filas):
        r = f'nutrition[{i}]'
        objeto(x, r, ('product', 'presentation', 'kcal_milli_per_100g', 'source', 'source_ref', 'status'))
        p, pr = x['product'], x['presentation']
        if p not in d.products:
            raise ErrorMotor('unknown_product', r + '.product')
        if pr is not None and (pr not in d.pres or d.pres[pr]['product'] != p):
            raise ErrorMotor('unknown_presentation', r + '.presentation')
        if (p, pr) in d.nutrition:
            raise ErrorMotor('invalid_argument', r + ' repetido')
        enum(x['source'], r + '.source', FUENTES)
        texto(x['source_ref'], r + '.source_ref')
        d.nutrition[(p, pr)] = {'kcal': entero(x['kcal_milli_per_100g'], r + '.kcal_milli_per_100g', 0),
                                'status': enum(x['status'], r + '.status', ('proposed', 'confirmed'))}


# ---- resto ----
def _resto(raw, d):
    d.persons = []
    for i, x in enumerate(lista(raw['persons'], 'persons')):
        objeto(x, f'persons[{i}]', ('code', 'has_meal_plan'))
        d.persons.append({'code': texto(x['code'], f'persons[{i}].code'), 'plan': booleano(x['has_meal_plan'], f'persons[{i}].has_meal_plan')})
    d.person = {p['code']: p for p in d.persons}
    d.stores = []
    for i, x in enumerate(lista(raw.get('stores', []), 'stores')):
        r = f'stores[{i}]'
        objeto(x, r, ('code', 'kind', 'groups', 'next_visit'))
        kind = enum(x['kind'], r + '.kind', ('chain', 'branch', 'market'))
        nv = x['next_visit']
        if (kind == 'market') != (nv is not None):
            raise ErrorMotor('invalid_argument', r + '.next_visit')
        for g in lista(x['groups'], r + '.groups'):
            if g not in d.groups:
                raise ErrorMotor('invalid_argument', r + '.groups')
        d.stores.append({'code': texto(x['code'], r + '.code'), 'kind': kind, 'groups': list(x['groups']),
                         'visit': fecha(nv, r + '.next_visit') if nv is not None else None})
    d.store = {s['code']: s for s in d.stores}
    d.plan = []
    for i, x in enumerate(lista(raw.get('meal_plan', []), 'meal_plan')):
        r = f'meal_plan[{i}]'
        objeto(x, r, ('date', 'slot', 'person', 'recipe', 'servings_milli', 'eating_out'))
        d.plan.append({'date': fecha(x['date'], r + '.date'), 'slot': enum(x['slot'], r + '.slot', SLOTS), 'person': x['person'],
                       'recipe': x['recipe'], 'serv': x['servings_milli'], 'out': booleano(x['eating_out'], r + '.eating_out'), 'idx': i})
    d.inv = []
    for i, x in enumerate(lista(raw.get('inventory', []), 'inventory')):
        r = f'inventory[{i}]'
        objeto(x, r, ('code', 'product', 'recipe', 'mode', 'quantity_milli'),
               ('counted_date', 'rate_quantity_milli', 'rate_period_days', 'reorder_level_milli', 'expires_date'))
        if (x['product'] is None) == (x['recipe'] is None):
            raise ErrorMotor('invalid_argument', r + ': product o recipe')
        mode = enum(x['mode'], r + '.mode', ('stock', 'consumption_rate'))
        it = {'code': texto(x['code'], r + '.code'), 'product': x['product'], 'recipe': x['recipe'], 'mode': mode,
              'q': entero(x['quantity_milli'], r + '.quantity_milli', 0),
              'counted': fecha(x['counted_date'], r + '.counted_date') if x.get('counted_date') is not None else None,
              'rq': entero(x.get('rate_quantity_milli'), r + '.rate_quantity_milli', 1, nulo=True),
              'rp': entero(x.get('rate_period_days'), r + '.rate_period_days', 1, nulo=True),
              'level': entero(x.get('reorder_level_milli'), r + '.reorder_level_milli', 0, nulo=True),
              'expires': fecha(x['expires_date'], r + '.expires_date') if x.get('expires_date') is not None else None}
        if x['product'] is not None and x['product'] not in d.products:
            raise ErrorMotor('unknown_product', r + '.product')
        if x['recipe'] is not None and (x['recipe'] not in d.recipes or mode != 'stock'):
            raise ErrorMotor('invalid_argument', r + '.recipe')
        if mode == 'consumption_rate' and None in (it['counted'], it['rq'], it['rp']):
            raise ErrorMotor('missing_field', r + ': ritmo requiere counted_date, rate_quantity_milli y rate_period_days')
        d.inv.append(it)
    d.choices = {}
    for i, x in enumerate(lista(raw.get('list_choices', []), 'list_choices')):
        r = f'list_choices[{i}]'
        objeto(x, r, ('product', 'kind', 'store', 'presentation', 'quantity_milli', 'is_bought'))
        p, k = x['product'], enum(x['kind'], r + '.kind', ('calculated', 'manual'))
        if p not in d.products:
            raise ErrorMotor('unknown_product', r + '.product')
        if x['presentation'] is not None and (x['presentation'] not in d.pres or d.pres[x['presentation']]['product'] != p):
            raise ErrorMotor('unknown_presentation', r + '.presentation')
        if x['store'] is not None and x['store'] not in d.store:
            raise ErrorMotor('invalid_argument', r + '.store')
        if (p, k) in d.choices:
            raise ErrorMotor('invalid_argument', r + ' repetido')
        if k == 'manual' and x['quantity_milli'] is None:
            raise ErrorMotor('invalid_argument', r + '.quantity_milli: obligatoria en manual')
        d.choices[(p, k)] = {'store': x['store'], 'pres': x['presentation'], 'q': entero(x['quantity_milli'], r + '.quantity_milli', 0, nulo=True),
                             'bought': booleano(x['is_bought'], r + '.is_bought')}
    d.quotes = {}
    vistos = set()
    for i, x in enumerate(lista(raw.get('quotes', []), 'quotes')):
        r = f'quotes[{i}]'
        objeto(x, r, ('product', 'store', 'presentation', 'price_cents', 'size_milli', 'observed_date', 'is_stale', 'recommended'))
        p = x['product']
        if p not in d.products:
            raise ErrorMotor('unknown_product', r + '.product')
        if x['presentation'] not in d.pres or d.pres[x['presentation']]['product'] != p:
            raise ErrorMotor('unknown_presentation', r + '.presentation')
        clave = (p, x['store'], x['presentation'])
        if clave in vistos:
            raise ErrorMotor('invalid_argument', r + ' repetida')
        vistos.add(clave)
        d.quotes.setdefault(p, []).append({'store': texto(x['store'], r + '.store'), 'pres': x['presentation'],
                                           'price': entero(x['price_cents'], r + '.price_cents', 0), 'size': entero(x['size_milli'], r + '.size_milli', 1),
                                           'observed': fecha(x['observed_date'], r + '.observed_date'), 'stale': booleano(x['is_stale'], r + '.is_stale'),
                                           'rec': booleano(x['recommended'], r + '.recommended')})
    for p, qs in d.quotes.items():
        if sum(q['rec'] for q in qs) != 1:
            raise ErrorMotor('invalid_argument', f'quotes: {p} debe tener exactamente una recomendada')
    d.profiles = {}
    for i, x in enumerate(lista(raw.get('health_profiles', []), 'health_profiles')):
        r = f'health_profiles[{i}]'
        objeto(x, r, ('person',), ('height_mm', 'birth_date', 'sex', 'activity', 'goal', 'adjust_lose_kcal', 'adjust_gain_kcal',
                                   'daily_kcal_target', 'dietary_restrictions'))
        if x['person'] not in d.person or x['person'] in d.profiles:
            raise ErrorMotor('invalid_argument', r + '.person')
        d.profiles[x['person']] = {
            'height': entero(x.get('height_mm'), r + '.height_mm', 1, nulo=True),
            'birth': fecha(x['birth_date'], r + '.birth_date') if x.get('birth_date') is not None else None,
            'sex': enum(x.get('sex'), r + '.sex', ('female', 'male'), nulo=True),
            'activity': enum(x.get('activity'), r + '.activity', ACTIVIDADES, nulo=True),
            'goal': enum(x.get('goal'), r + '.goal', ('lose', 'maintain', 'gain'), nulo=True),
            'adj_lose': entero(x.get('adjust_lose_kcal'), r + '.adjust_lose_kcal', 0, nulo=True),
            'adj_gain': entero(x.get('adjust_gain_kcal'), r + '.adjust_gain_kcal', 0, nulo=True),
            'target': entero(x.get('daily_kcal_target'), r + '.daily_kcal_target', 500, 10000, nulo=True)}
    d.weights = []
    for i, x in enumerate(lista(raw.get('weights', []), 'weights')):
        r = f'weights[{i}]'
        objeto(x, r, ('person', 'date', 'weight_g'))
        if x['person'] not in d.person:
            raise ErrorMotor('invalid_argument', r + '.person')
        d.weights.append({'person': x['person'], 'date': fecha(x['date'], r + '.date'), 'g': entero(x['weight_g'], r + '.weight_g', 1)})
    s = raw.get('settings')
    if s is None:
        s = {'house_adjust_lose_kcal': 0, 'house_adjust_gain_kcal': 0}
    objeto(s, 'settings', ('house_adjust_lose_kcal', 'house_adjust_gain_kcal'))
    d.house_lose = entero(s['house_adjust_lose_kcal'], 'settings.house_adjust_lose_kcal', 0)
    d.house_gain = entero(s['house_adjust_gain_kcal'], 'settings.house_adjust_gain_kcal', 0)


def cargar(raw):
    """Input crudo -> Datos validados. Lanza ErrorMotor (unknown_field, missing_field, ...)."""
    objeto(raw, 'datos', ('today', 'persons'),
           ('catalog_file', 'catalog', 'recipe_book_file', 'recipes', 'nutrition_file', 'nutrition', 'stores', 'meal_plan',
            'inventory', 'list_choices', 'quotes', 'health_profiles', 'weights', 'settings'))
    d = Datos()
    d.today = fecha(raw['today'], 'today')
    _catalogo(raw, d)
    _recetas(raw, d)
    _nutricion(raw, d)
    _resto(raw, d)
    return d
