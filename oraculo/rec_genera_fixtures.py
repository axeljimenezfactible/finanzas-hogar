#!/usr/bin/env python3
"""Genera las fixtures de oraculo/fixtures/s4/ (S4-interfaz.md §12). Autor único de esos archivos; determinista.

Los `esperado` de las fixtures con números sintéticos están escritos A MANO (aritmética en los comentarios de cada caso).
Los de recetas reales (kcal) se calculan aquí con una segunda vía, independiente de rec_calculo.py (solo `fractions`).
Uso: python3 -I oraculo/rec_genera_fixtures.py
"""
import csv
import json
from decimal import Decimal
from fractions import Fraction
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[1]
SALIDA = RAIZ / 'oraculo' / 'fixtures' / 's4'
BL_LISTA = 'docs/estado/bloqueos/S4-reglas-lista.md'
BL_MERCADO = 'docs/estado/bloqueos/S4-asignacion-mercado.md'


def Q(n, d=1):
    f = Fraction(n, d)
    return {'num': f.numerator, 'den': f.denominator}


def escribe(caso, descripcion, fuente, datos, consultas, esperado, provisional=None):
    fx = {'formato': 'fixture/1', 'motor': 's4', 'caso': caso, 'descripcion': descripcion,
          'fuente': [{'tipo': t, 'ruta': r} for t, r in fuente], 'provisional': provisional,
          'entradas': {'datos': datos, 'consultas': [{'id': i, 'funcion': f, 'argumentos': a} for i, f, a in consultas]},
          'esperado': esperado}
    with open(SALIDA / f'{caso}.json', 'w', encoding='utf-8') as fh:
        fh.write(json.dumps(fx, ensure_ascii=False, indent=2) + '\n')


# ---------- catálogo y recetas sintéticos (ILUSTRATIVOS) ----------
CAT = {
    'groups': [{'code': 'verduras', 'name': 'Verduras', 'purchase_channel': 'market'},
               {'code': 'frutas', 'name': 'Frutas', 'purchase_channel': 'market'},
               {'code': 'carnes', 'name': 'Carnes', 'purchase_channel': 'market'},
               {'code': 'abarrotes', 'name': 'Abarrotes', 'purchase_channel': 'chain'},
               {'code': 'lacteos', 'name': 'Lácteos y huevo', 'purchase_channel': 'chain'}],
    'products': [
        {'code': 'arroz', 'name': 'Arroz', 'base_unit': 'g', 'group': 'abarrotes', 'grams_per_unit_milli': None},
        {'code': 'frijol', 'name': 'Frijol', 'base_unit': 'g', 'group': 'abarrotes', 'grams_per_unit_milli': None},
        {'code': 'aceite-soya', 'name': 'Aceite de soya', 'base_unit': 'ml', 'group': 'abarrotes', 'grams_per_unit_milli': 921},
        {'code': 'huevo', 'name': 'Huevo', 'base_unit': 'piece', 'group': 'lacteos', 'grams_per_unit_milli': 50000},
        {'code': 'servilletas', 'name': 'Servilletas', 'base_unit': 'piece', 'group': 'abarrotes', 'grams_per_unit_milli': None},
        {'code': 'harina', 'name': 'Harina', 'base_unit': 'g', 'group': 'abarrotes', 'grams_per_unit_milli': None},
        {'code': 'jitomate', 'name': 'Jitomate', 'base_unit': 'g', 'group': 'verduras', 'grams_per_unit_milli': None},
        {'code': 'cebolla', 'name': 'Cebolla', 'base_unit': 'g', 'group': 'verduras', 'grams_per_unit_milli': None},
        {'code': 'fresa', 'name': 'Fresa', 'base_unit': 'g', 'group': 'frutas', 'grams_per_unit_milli': None},
        {'code': 'pollo', 'name': 'Pollo', 'base_unit': 'g', 'group': 'carnes', 'grams_per_unit_milli': None}],
    'presentations': [
        {'code': 'arroz.1kg', 'product': 'arroz', 'name': 'Arroz 1 kg', 'size_milli': 1000000},
        {'code': 'arroz.750g', 'product': 'arroz', 'name': 'Arroz 750 g', 'size_milli': 750000},
        {'code': 'frijol.1kg', 'product': 'frijol', 'name': 'Frijol 1 kg', 'size_milli': 1000000},
        {'code': 'aceite-soya.946ml', 'product': 'aceite-soya', 'name': 'Aceite de soya 946 ml', 'size_milli': 946000},
        {'code': 'huevo.carton-12', 'product': 'huevo', 'name': 'Huevo cartón de 12', 'size_milli': 12000},
        {'code': 'servilletas.420', 'product': 'servilletas', 'name': 'Servilletas 420 piezas', 'size_milli': 420000},
        {'code': 'harina.1kg', 'product': 'harina', 'name': 'Harina 1 kg', 'size_milli': 1000000},
        {'code': 'harina.500g', 'product': 'harina', 'name': 'Harina 500 g', 'size_milli': 500000},
        {'code': 'jitomate.kilo', 'product': 'jitomate', 'name': 'Jitomate por kilo', 'size_milli': 1000000},
        {'code': 'cebolla.kilo', 'product': 'cebolla', 'name': 'Cebolla por kilo', 'size_milli': 1000000},
        {'code': 'fresa.kilo', 'product': 'fresa', 'name': 'Fresa por kilo', 'size_milli': 1000000},
        {'code': 'pollo.kilo', 'product': 'pollo', 'name': 'Pollo por kilo', 'size_milli': 1000000}]}
PERSONAS = [{'code': 'persona-a', 'has_meal_plan': True}, {'code': 'persona-b', 'has_meal_plan': True},
            {'code': 'externa-c', 'has_meal_plan': False}]


def ing(prod, pres, q, opt=False):
    return {'product': prod, 'presentation': pres, 'quantity_milli': q, 'optional': opt}


def rec(code, name, slots, ings, base=2, origin='claude_seed'):
    return {'code': code, 'name': name, 'slots': slots, 'base_servings': base, 'prep_minutes': 10, 'cook_minutes': 20,
            'origin': origin, 'is_active': True, 'ingredients': ings}


R1 = rec('arroz-con-frijol', 'Arroz con frijol', ['lunch', 'dinner'],
         [ing('arroz', 'arroz.1kg', 450000), ing('frijol', 'frijol.1kg', 300000), ing('aceite-soya', 'aceite-soya.946ml', 15000)])
R2 = rec('huevo-con-frijol', 'Huevo con frijol', ['breakfast', 'snack'],
         [ing('huevo', 'huevo.carton-12', 4000), ing('frijol', 'frijol.1kg', 150000), ing('aceite-soya', 'aceite-soya.946ml', 10000)])


def quote(prod, store, pres, price, size, rec_=True, fecha='2026-10-10', stale=False):
    return {'product': prod, 'store': store, 'presentation': pres, 'price_cents': price, 'size_milli': size,
            'observed_date': fecha, 'is_stale': stale, 'recommended': rec_}


QW = [quote('arroz', 'walmart', 'arroz.1kg', 2800, 1000000), quote('frijol', 'walmart', 'frijol.1kg', 3000, 1000000),
      quote('aceite-soya', 'walmart', 'aceite-soya.946ml', 3990, 946000), quote('huevo', 'walmart', 'huevo.carton-12', 6000, 12000)]
TIENDAS_CADENA = [{'code': 'walmart', 'kind': 'chain', 'groups': [], 'next_visit': None},
                  {'code': 'bodega', 'kind': 'chain', 'groups': [], 'next_visit': None}]


def entrada(fecha, slot, persona, receta, serv, fuera=False):
    return {'date': fecha, 'slot': slot, 'person': persona, 'recipe': receta, 'servings_milli': serv, 'eating_out': fuera}


# Plan P (semana 2026-W42, lunes 2026-10-12): R1 suma 2000 + 2000 + 1000 = 5000 milésimas de porción (factor 5000/2000 = 5/2)
PLAN_P = [entrada('2026-10-12', 'lunch', 'persona-a', 'arroz-con-frijol', 2000),
          entrada('2026-10-12', 'lunch', 'persona-b', 'arroz-con-frijol', 2000),
          entrada('2026-10-13', 'breakfast', 'persona-a', 'huevo-con-frijol', 2000),
          entrada('2026-10-14', 'dinner', 'persona-b', 'arroz-con-frijol', 1000)]


def inv(code, prod=None, recipe=None, mode='stock', q=0, counted=None, rq=None, rp=None, level=None, expira=None):
    return {'code': code, 'product': prod, 'recipe': recipe, 'mode': mode, 'quantity_milli': q, 'counted_date': counted,
            'rate_quantity_milli': rq, 'rate_period_days': rp, 'reorder_level_milli': level, 'expires_date': expira}


def eleccion(prod, kind='calculated', store=None, pres=None, q=None, bought=False):
    return {'product': prod, 'kind': kind, 'store': store, 'presentation': pres, 'quantity_milli': q, 'is_bought': bought}


def datos(today='2026-10-11', **kw):
    d = {'today': today, 'catalog': CAT, 'recipes': [R1, R2], 'persons': PERSONAS, 'stores': TIENDAS_CADENA, 'meal_plan': PLAN_P,
         'inventory': [], 'list_choices': [], 'quotes': QW}
    d.update(kw)
    return d


def art(p, pres, needed, pk, left, est, recetas, *, kind='calculated', qty=None, reasons=('recipe',), store='walmart', bought=False,
        pstore='walmart', rstore='walmart', rpres=None, fecha='2026-10-10', stale=False):
    sin = est is None and pstore is None
    return {'product': p, 'kind': kind, 'needed_milli': needed, 'quantity_milli': qty, 'presentation': pres, 'packages_to_buy': pk,
            'leftover_milli': left, 'reasons': list(reasons), 'recipes': list(recetas), 'store': store, 'is_bought': bought,
            'recommended_store': rstore, 'recommended_presentation': pres if rpres is None else rpres, 'price_store': pstore,
            'price_observed_date': None if sin else fecha, 'is_stale_price': None if sin else stale, 'estimated_cents': est}


def grp(v, store, visita, total):
    return {'vertical': v, 'store': store, 'visit_date': visita, 'estimated_total_cents': total}


SEM = {'week': '2026-W42'}
FUENTE_LISTA = [('criterio', 'S4-03'), ('regla', 'RN-15')]
A, B_, F, H = 'arroz-con-frijol', 'huevo-con-frijol', None, None


def casos_escalado():
    rc = rec('escala-prueba', 'Receta de escalado', ['lunch'],
             [ing('harina', 'harina.1kg', 250000), ing('aceite-soya', 'aceite-soya.946ml', 7000), ing('huevo', 'huevo.carton-12', 1000)])
    cs, esp = [], {}
    for sid, serv in (('x2', 4000), ('mitad', 1000), ('x3-2', 3000), ('x3', 6000), ('fraccion', 333)):
        cs.append((sid, 'scale_recipe', {'recipe': 'escala-prueba', 'servings_milli': serv}))
    f = lambda n: Fraction(n)
    def ings(k):  # k = servings_milli / 2000, exacto
        return [{'product': 'harina', 'unit': 'g', 'quantity': Q(250000 * k)}, {'product': 'aceite-soya', 'unit': 'ml', 'quantity': Q(7000 * k)},
                {'product': 'huevo', 'unit': 'piece', 'quantity': Q(1000 * k)}]
    # 2 -> 4: x2 (500000, 14000, 2000); 2 -> 1: x1/2 (125000, 3500, 500); 2 -> 3: x3/2 (375000, 10500, 1500, RN-14);
    # 2 -> 6: x3; 333/2000: 250000*333/2000 = 41625; 7000*333/2000 = 2331/2 (1165.5); 1000*333/2000 = 333/2000
    for sid, serv, k in (('x2', 4000, Fraction(2)), ('mitad', 1000, Fraction(1, 2)), ('x3-2', 3000, Fraction(3, 2)),
                         ('x3', 6000, Fraction(3)), ('fraccion', 333, Fraction(333, 2000))):
        esp[sid] = {'recipe': 'escala-prueba', 'servings_milli': serv, 'ingredients': ings(k)}
    esp['x3-2']['ingredients'][0]['quantity'] = Q(375000)
    cs += [('cero', 'scale_recipe', {'recipe': 'escala-prueba', 'servings_milli': 0}),
           ('negativa', 'scale_recipe', {'recipe': 'escala-prueba', 'servings_milli': -1000})]
    esp['cero'] = {'error': 'invalid_servings'}
    esp['negativa'] = {'error': 'invalid_servings'}
    d = {'today': '2026-10-11', 'catalog': CAT, 'recipes': [rc], 'persons': PERSONAS}
    escribe('s4-escalado-2a4', 'RN-14, S4-02: receta de 2 porciones con 250000 (g), 7000 (ml) y 1000 (pieza) escalada a 4, 1, 3, 6 y 0.333 porciones. '
            'Exacto y sin redondeo: x2, x1/2, x3/2 (375000 g; 7 ml -> 10500 = 21000/2; 1 pieza -> 1500), x3 y la fracción 333/2000 '
            '(250000*333/2000 = 41625; 7000*333/2000 = 2331/2; 1000*333/2000 = 333/2000). Porciones <= 0: invalid_servings. ILUSTRATIVO.',
            [('criterio', 'S4-02'), ('regla', 'RN-14'), ('manual', 'cuentas en descripcion')], d, cs, esp)


def casos_lista():
    # Totales de P (cuenta a mano): arroz 450000*5/2 = 1125000; frijol 300000*5/2 + 150000 = 900000; aceite 15000*5/2 + 10000 = 47500; huevo 4000.
    base_txt = ('Plan P de 2026-W42: arroz-con-frijol 2000+2000+1000 porciones (x5/2) y huevo-con-frijol 2000 (x1). Totales: arroz 1125000, '
                'frijol 750000+150000 = 900000, aceite 37500+10000 = 47500, huevo 4000. Precios Walmart 2026-10-10: arroz 2800/1 kg, frijol 3000/1 kg, '
                'aceite 3990/946 ml, huevo 6000/cartón de 12. ILUSTRATIVO.')
    it_aceite = lambda n, pk, left, est=3990: art('aceite-soya', 'aceite-soya.946ml', n, pk, left, est, [A, B_])
    sin = [it_aceite(47500, 1, 898500),
           art('arroz', 'arroz.1kg', 1125000, 2, 875000, 5600, [A]),            # techo(1125000/1000000) = 2; 2000000 - 1125000 = 875000; 2 x 2800
           art('frijol', 'frijol.1kg', 900000, 1, 100000, 3000, [A, B_]),
           art('huevo', 'huevo.carton-12', 4000, 1, 8000, 6000, [B_])]
    esp = {'lista': {'week': '2026-W42', 'items': sin, 'groups': [grp('chain', 'walmart', None, 18590)],
                     'estimated_total_cents': 18590, 'missing_prices': []}}  # 3990 + 5600 + 3000 + 6000
    escribe('s4-lista-sin-existencia', 'S4-03: sin existencia. ' + base_txt + ' Paquetes: aceite 1 (sobra 898500), arroz 2 (sobra 875000), frijol 1 (100000), huevo 1 (8000). Total 18590.',
            FUENTE_LISTA + [('manual', 'cuentas en descripcion')], datos(), [('lista', 'shopping_list', SEM)], esp)

    d = datos(inventory=[inv('inv-arroz', 'arroz', q=300000), inv('inv-frijol', 'frijol', q=150000)])
    par = [it_aceite(47500, 1, 898500), art('arroz', 'arroz.1kg', 825000, 1, 175000, 2800, [A]),
           art('frijol', 'frijol.1kg', 750000, 1, 250000, 3000, [A, B_]), art('huevo', 'huevo.carton-12', 4000, 1, 8000, 6000, [B_])]
    esp = {'lista': {'items': par, 'groups': [grp('chain', 'walmart', None, 15790)], 'estimated_total_cents': 15790, 'missing_prices': []}}
    escribe('s4-lista-existencia-parcial', 'S4-03: existencia parcial. ' + base_txt + ' Existencia: arroz 300000 y frijol 150000. Neta: arroz 825000 (1 paquete, sobra 175000), '
            'frijol 750000 (1 paquete, sobra 250000). Total 2800 + 3000 + 3990 + 6000 = 15790.',
            FUENTE_LISTA + [('manual', 'cuentas en descripcion')], d, [('lista', 'shopping_list', SEM)], esp)

    d = datos(inventory=[inv('inv-arroz', 'arroz', q=2000000), inv('inv-frijol', 'frijol', q=900000), inv('inv-aceite', 'aceite-soya', q=100000),
                         inv('inv-huevo', 'huevo', q=1000)])
    esp = {'lista': {'items': [art('huevo', 'huevo.carton-12', 3000, 1, 9000, 6000, [B_])], 'groups': [grp('chain', 'walmart', None, 6000)],
                     'estimated_total_cents': 6000, 'missing_prices': []}}
    escribe('s4-lista-existencia-mayor', 'S4-03: existencia mayor o igual que la necesidad = 0 a comprar y sin artículo (sin negativos). ' + base_txt +
            ' Existencia: arroz 2000000 (> 1125000), frijol 900000 (= 900000), aceite 100000 (> 47500) y huevo 1000: solo queda el huevo (4000 - 1000 = 3000; 1 paquete, sobra 9000).',
            FUENTE_LISTA + [('manual', 'cuentas en descripcion')], d, [('lista', 'shopping_list', SEM)], esp)

    d = datos(inventory=[inv('sob-a', recipe='arroz-con-frijol', q=1000), inv('sob-b', recipe='arroz-con-frijol', q=2000, expira='2026-10-11')])
    items = [it_aceite(40000, 1, 906000), art('arroz', 'arroz.1kg', 900000, 1, 100000, 2800, [A]),
             art('frijol', 'frijol.1kg', 750000, 1, 250000, 3000, [A, B_]), art('huevo', 'huevo.carton-12', 4000, 1, 8000, 6000, [B_])]
    esp = {'lista': {'items': items, 'groups': [grp('chain', 'walmart', None, 15790)], 'estimated_total_cents': 15790, 'missing_prices': []}}
    escribe('s4-lista-sobrantes', 'S4-03, §6.4: un sobrante registrado de 1 porción de arroz-con-frijol cubre la primera entrada del plan (2000 -> 1000 porciones); '
            'otro de 2 porciones caducó el 2026-10-11 y no cubre ninguna. ' + base_txt + ' R1 efectivo = 1000+2000+1000 = 4000 (x2): arroz 900000, frijol 600000+150000 = 750000, '
            'aceite 30000+10000 = 40000, huevo 4000. Total 2800 + 3000 + 3990 + 6000 = 15790.',
            FUENTE_LISTA + [('manual', 'cuentas en descripcion')], d, [('lista', 'shopping_list', SEM)], esp)

    # Cantidad fijada y artículos manual: el ejemplo ILUSTRATIVO de §6.6
    r = rec('base-cantidad', 'Receta base de cantidad', ['lunch'], [ing('arroz', 'arroz.1kg', 500000), ing('frijol', 'frijol.1kg', 1200000)])
    d = datos(recipes=[r], meal_plan=[entrada('2026-10-12', 'lunch', 'persona-a', 'base-cantidad', 2000)],
              list_choices=[eleccion('arroz', q=3000000), eleccion('frijol', q=500000), eleccion('aceite-soya', 'manual', q=1892000, bought=True),
                            eleccion('servilletas', 'manual', q=1000)],
              quotes=QW[:3])
    items = [art('arroz', 'arroz.1kg', 500000, 3, 2500000, 8400, ['base-cantidad'], qty=3000000),
             art('frijol', 'frijol.1kg', 1200000, 1, 0, 3000, ['base-cantidad'], qty=500000),
             art('aceite-soya', 'aceite-soya.946ml', None, 2, 1892000, 7980, [], kind='manual', qty=1892000, reasons=('manual',), bought=True),
             art('servilletas', None, None, None, None, None, [], kind='manual', qty=1000, reasons=('manual',), store=None, pstore=None, rstore=None)]
    items[3]['recommended_presentation'] = None
    esp = {'lista': {'items': items, 'groups': [grp('chain', 'walmart', None, 19380), grp('chain', None, None, 0)], 'estimated_total_cents': 19380,
                     'missing_prices': ['servilletas']}}
    escribe('s4-lista-cantidad-manual', 'S4-03, S4-06, C1 quantity_milli, §6.6 (ILUSTRATIVO): arroz necesita 500000 y la persona fijó 3000000 (3 paquetes, sobra 3000000-500000 = 2500000, 3 x 2800 = 8400); '
            'frijol necesita 1200000 y fijó 500000 (1 paquete, sobra max(0, 1000000-1200000) = 0, 3000); aceite de soya manual comprado de 1892000 = 2 paquetes de 946000, 7980, suma aunque is_bought; '
            'servilletas manual sin cotización: paquetes, sobrante y precio nulos y en missing_prices. Grupo (chain, walmart) = 8400 + 3000 + 7980 = 19380; (chain, null) = 0; total de la raíz 19380.',
            [('criterio', 'S4-03'), ('criterio', 'S4-06'), ('regla', 'RN-15'), ('manual', '§6.6 del contrato')], d, [('lista', 'shopping_list', SEM)], esp)

    # Reposición por ritmo (PROVISIONAL, S4-reglas-lista pregunta 2 y 1)
    d = datos(meal_plan=[entrada('2026-10-13', 'lunch', 'persona-a', 'arroz-con-frijol', 2000)],
              inventory=[inv('inv-frijol-ritmo', 'frijol', mode='consumption_rate', q=1000000, counted='2026-10-08', rq=700000, rp=7, level=300000),
                         inv('sob-r2', recipe='huevo-con-frijol', q=3000)])
    # lunes 10-12 (4 días): 1000000 - 700000*4/7 = 600000; domingo 10-18 (10 días): 1000000 - 1000000 = 0 <= 300000 -> R = 700000.
    # frijol: 300000 + 700000 - 600000 = 400000
    items = [art('aceite-soya', 'aceite-soya.946ml', 15000, 1, 931000, 3990, [A]), art('arroz', 'arroz.1kg', 450000, 1, 550000, 2800, [A]),
             art('frijol', 'frijol.1kg', 400000, 1, 600000, 3000, [A], reasons=('recipe', 'replenishment'))]
    esp = {'lista': {'items': items, 'groups': [grp('chain', 'walmart', None, 9790)], 'estimated_total_cents': 9790, 'missing_prices': []},
           'plan': {'leftovers': [{'recipe': 'huevo-con-frijol', 'servings_milli': 3000}]}}
    escribe('s4-lista-reposicion-ritmo', 'PROVISIONAL (S4-reglas-lista, preguntas 1 y 2, opción A). Plan: arroz-con-frijol 2000 el 2026-10-13. Frijol por ritmo: contado 1000000 el 2026-10-08, '
            'ritmo 700000 cada 7 días, nivel 300000. Existencia al lunes 10-12 (4 días) = 1000000 - 400000 = 600000; al domingo 10-18 (10 días) = 0 <= 300000, así que R = 700000. '
            'Neta frijol = 300000 + 700000 - 600000 = 400000 (motivos recipe y replenishment). Un sobrante registrado de 3 porciones de huevo-con-frijol que el plan no consume sale en leftovers. '
            'Arroz 450000, aceite 15000. Total 2800 + 3000 + 3990 = 9790. ILUSTRATIVO.',
            [('criterio', 'S4-03'), ('regla', 'RN-15'), ('manual', 'cuentas en descripcion')], d,
            [('lista', 'shopping_list', SEM), ('plan', 'week_plan', SEM)], esp,
            provisional={'bloqueo': BL_LISTA, 'que': 'reposición por ritmo (pregunta 2) y sobrantes previstos (pregunta 1), opción A'})


def casos_inventario():
    d = datos(meal_plan=[], recipes=[R1, R2], inventory=[
        inv('frijol-ritmo', 'frijol', mode='consumption_rate', q=2838000, counted='2026-09-11', rq=946000, rp=21, level=500000),
        inv('arroz-stock', 'arroz', q=400000, level=400000),
        inv('aceite-stock', 'aceite-soya', q=0, level=200000),
        inv('sobrante-r1', recipe='arroz-con-frijol', q=2000),
        inv('agotado-ritmo', 'huevo', mode='consumption_rate', q=1000000, counted='2026-10-01', rq=1000000, rp=5, level=None),
        inv('huevo-ritmo', 'huevo', mode='consumption_rate', q=12000, counted='2026-10-08', rq=6000, rp=7, level=9500),
        inv('mitad-ritmo', 'servilletas', mode='consumption_rate', q=1000, counted='2026-10-10', rq=1, rp=2, level=None)])
    # frijol: 2838000 - HALF_UP(946000*30/21 = 1351428.57) = 1486571 (ejemplo RN-15); semanal HALF_UP(946000*7/21 = 315333.33) = 315333;
    #   días: piso(1486571*21/946000 = 32.99999) = 32 (31217991 < 31218000 = 33*946000)
    # arroz: 400000 = nivel -> reorder (igualdad). aceite: 0 -> out (agotado prevalece aunque 0 <= nivel). sobrante: stock 2000, sin nivel -> sufficient
    # leche-ritmo: 10 días x 200000 = 2000000 > 1000000 -> 0 (nunca negativo) -> out; semanal 1000000*7/5 = 1400000; días 0
    # huevo-ritmo: 12000 - HALF_UP(6000*3/7 = 2571.43) = 9429 <= 9500 -> reorder; semanal 6000; días piso(9429*7/6000 = 11.0005) = 11
    # sal-ritmo: 1 día: 1000 - HALF_UP(1*1/2 = 0.5) = 999 (la mitad sube) -> sufficient; semanal HALF_UP(3.5) = 4; días 999*2/1 = 1998
    it = lambda c, e, s, w=None, dl=None: {'code': c, 'estimated_quantity_milli': e, 'status': s, 'weekly_consumption_milli': w, 'days_left': dl}
    esp = {'inv': {'items': [it('frijol-ritmo', 1486571, 'sufficient', 315333, 32), it('arroz-stock', 400000, 'reorder'), it('aceite-stock', 0, 'out'),
                             it('sobrante-r1', 2000, 'sufficient'), it('agotado-ritmo', 0, 'out', 1400000, 0), it('huevo-ritmo', 9429, 'reorder', 6000, 11),
                             it('mitad-ritmo', 999, 'sufficient', 4, 1998)],
                   'summary': {'sufficient': 3, 'reorder': 2, 'out': 2}}}
    escribe('s4-inventario-modos', 'S4-04, RN-15, K7 (hoy 2026-10-11): ritmo con el ejemplo de RN-15 (2838000 contados, 946000 cada 21 días, 30 días después: 2838000 - 1351429 = 1486571; '
            'semanal 315333; días completos 32); stock en igualdad con el nivel (reorder); stock en 0 con nivel (out prevalece); sobrante de receta en porciones x 1000; ritmo que se agota sin negativos; '
            'ritmo que baja a reorder; mitad que sube (999). ILUSTRATIVO salvo el ejemplo de RN-15.',
            [('criterio', 'S4-04'), ('regla', 'RN-15'), ('decision', 'K7'), ('manual', 'cuentas en descripcion')], d, [('inv', 'inventory_status', {})], esp)


# ---------- plan semanal ----------
NUT_SINT = [{'product': 'arroz', 'presentation': None, 'kcal_milli_per_100g': 360000, 'source': 'manual', 'source_ref': 'ILUSTRATIVO', 'status': 'confirmed'},
            {'product': 'frijol', 'presentation': None, 'kcal_milli_per_100g': 340000, 'source': 'manual', 'source_ref': 'ILUSTRATIVO', 'status': 'confirmed'},
            {'product': 'aceite-soya', 'presentation': None, 'kcal_milli_per_100g': 884000, 'source': 'manual', 'source_ref': 'ILUSTRATIVO', 'status': 'confirmed'},
            {'product': 'huevo', 'presentation': None, 'kcal_milli_per_100g': 143000, 'source': 'manual', 'source_ref': 'ILUSTRATIVO', 'status': 'confirmed'}]


def perfil(persona, **kw):
    p = {'person': persona, 'height_mm': None, 'birth_date': None, 'sex': None, 'activity': None, 'goal': None, 'adjust_lose_kcal': None,
         'adjust_gain_kcal': None, 'daily_kcal_target': None, 'dietary_restrictions': None}
    p.update(kw)
    return p


def casos_plan():
    plan = [entrada('2026-10-12', 'lunch', 'persona-a', 'arroz-con-frijol', 2000), entrada('2026-10-12', 'lunch', 'persona-b', 'arroz-con-frijol', 1000),
            entrada('2026-10-12', 'dinner', 'persona-a', None, None, fuera=True),
            entrada('2026-10-13', 'breakfast', 'persona-a', 'huevo-con-frijol', 2000), entrada('2026-10-13', 'breakfast', 'persona-b', 'huevo-con-frijol', 2000),
            entrada('2026-10-13', 'snack', 'persona-b', 'huevo-con-frijol', 1000)]
    d = datos(meal_plan=plan, nutrition=NUT_SINT,
              health_profiles=[perfil('persona-a', daily_kcal_target=2000),
                               perfil('persona-b', height_mm=1650, birth_date='1998-03-15', sex='female', activity='sedentary', goal='maintain')],
              weights=[{'person': 'persona-b', 'date': '2026-10-01', 'weight_g': 60000}])
    # R1 a 2000: arroz 450000*2800/1e6 = 1260; frijol 300000*3000/1e6 = 900; aceite 15000*3990/946000 = 63.2664; suma 2223.2664 -> 2223; a 1000: 1111.6332 -> 1112
    # R2 a 2000: huevo 4000*6000/12000 = 2000; frijol 150000*3000/1e6 = 450; aceite 10000*3990/946000 = 42.1776; suma 2492.1776 -> 2492; a 1000: 1246.0888 -> 1246
    # kcal R1 a 2000: arroz 450 g * 3.6 = 1620; frijol 300 g * 3.4 = 1020; aceite 15 ml * 0.921 = 13.815 g * 8.84 = 122.1246; suma 2762.1246 -> 2762; a 1000: 1381.0623 -> 1381
    # kcal R2 a 2000: huevo 4 pzas * 50 g = 200 g * 1.43 = 286; frijol 150 g * 3.4 = 510; aceite 10 ml * 0.921 = 9.21 g * 8.84 = 81.4164; suma 877.4164 -> 877; a 1000: 438.7082 -> 439
    # metas: persona-a confirmada 2000; persona-b sugerida 1596 (mujer 28 años, 60 kg, 165 cm, sedentary, maintain: 1330.25 * 6/5 = 1596.30)
    c = lambda f, s, p, r, sv, co, k, ps='walmart', pd='2026-10-10': {'date': f, 'slot': s, 'person': p, 'recipe': r, 'servings_milli': sv, 'eating_out': r is None,
                                                                  'cost_cents': co, 'kcal': k, 'price_store': ps if r else None, 'price_observed_date': pd if r else None}
    celdas = [c('2026-10-12', 'lunch', 'persona-a', 'arroz-con-frijol', 2000, 2223, 2762), c('2026-10-12', 'lunch', 'persona-b', 'arroz-con-frijol', 1000, 1112, 1381),
              c('2026-10-12', 'dinner', 'persona-a', None, None, None, None),
              c('2026-10-13', 'breakfast', 'persona-a', 'huevo-con-frijol', 2000, 2492, 877), c('2026-10-13', 'breakfast', 'persona-b', 'huevo-con-frijol', 2000, 2492, 877),
              c('2026-10-13', 'snack', 'persona-b', 'huevo-con-frijol', 1000, 1246, 439)]
    vacio = lambda f: {'date': f, 'cost_cents': 0, 'persons': [{'person': 'persona-a', 'kcal': 0, 'daily_kcal_target': 2000, 'difference_kcal': -2000},
                                                                 {'person': 'persona-b', 'kcal': 0, 'daily_kcal_target': 1596, 'difference_kcal': -1596}]}
    dias = [{'date': '2026-10-12', 'cost_cents': 3335, 'persons': [{'person': 'persona-a', 'kcal': 2762, 'daily_kcal_target': 2000, 'difference_kcal': 762},
                                                                    {'person': 'persona-b', 'kcal': 1381, 'daily_kcal_target': 1596, 'difference_kcal': -215}]},
            {'date': '2026-10-13', 'cost_cents': 6230, 'persons': [{'person': 'persona-a', 'kcal': 877, 'daily_kcal_target': 2000, 'difference_kcal': -1123},
                                                                    {'person': 'persona-b', 'kcal': 1316, 'daily_kcal_target': 1596, 'difference_kcal': -280}]}]
    dias += [vacio(f'2026-10-{x}') for x in (14, 15, 16, 17, 18)]
    # huecos: de los 21 pares día x tiempo obligatorio, solo están completos el almuerzo del 12 (las dos personas) y el desayuno del 13 -> 19; la cena del 12 tiene a persona-a "fuera" pero falta persona-b
    esp = {'plan': {'week': '2026-W42', 'cells': celdas, 'days': dias, 'week_cost_cents': 9565, 'unplanned_slots': 19, 'eating_out_slots': 1, 'leftovers': []},
           'semana-domingo': {'week': '2026-W42', 'monday': '2026-10-12'}, 'semana-lunes': {'week': '2026-W42', 'monday': '2026-10-12'},
           'semana-sabado': {'week': '2026-W41', 'monday': '2026-10-05'}, 'huecos-w42': {'unplanned_slots': 19}, 'huecos-w43': {'unplanned_slots': 21}}
    escribe('s4-plan-celdas', 'S4-05, RN-14, RN-19, K38, K43 (hoy domingo 2026-10-11): celdas con receta, «fuera» (sin receta, kcal ni costo), colación vacía de persona-a y llena de persona-b, '
            'costos y kcal por celda, día y semana con metas (persona-a confirmada 2000; persona-b sugerida 1596). Costos (centavos): R1 a 2000 = 1260+900+63.27 = 2223, a 1000 = 1112; R2 a 2000 = 2000+450+42.18 = 2492, a 1000 = 1246; '
            'día 12 = 2223+1112 = 3335, día 13 = 2492+2492+1246 = 6230, semana 9565. kcal: R1 2762 y 1381; R2 877 y 439; persona-b el día 13 = 877+439 = 1316. '
            'Huecos: 19 de 21 (completos: comida del 12 y desayuno del 13). El domingo abre la semana siguiente (K38). ILUSTRATIVO.',
            [('criterio', 'S4-05'), ('regla', 'RN-19'), ('decision', 'K38'), ('decision', 'K43'), ('manual', 'cuentas en descripcion')], d,
            [('plan', 'week_plan', SEM), ('semana-domingo', 'default_week', {'date': '2026-10-11'}), ('semana-lunes', 'default_week', {'date': '2026-10-12'}),
             ('semana-sabado', 'default_week', {'date': '2026-10-10'}), ('huecos-w42', 'unplanned_slots', SEM), ('huecos-w43', 'unplanned_slots', {'week': '2026-W43'})], esp)
    escribe('s4-plan-karen-rechazada', 'S4-05, decisiones A: una entrada de plan para una persona externa (sin plan) es person_without_plan. ILUSTRATIVO.',
            [('criterio', 'S4-05'), ('decision', 'A')], datos(meal_plan=[entrada('2026-10-12', 'lunch', 'externa-c', 'arroz-con-frijol', 2000)]),
            [('plan', 'week_plan', SEM)], {'plan': {'error': 'person_without_plan'}})
    escribe('s4-plan-fuera-con-receta', 'S4-05, RN-19: una celda «fuera» que además trae receta es eating_out_conflict. ILUSTRATIVO.',
            [('criterio', 'S4-05'), ('regla', 'RN-19')],
            datos(meal_plan=[dict(entrada('2026-10-12', 'lunch', 'persona-a', 'arroz-con-frijol', None, fuera=True))]),
            [('plan', 'week_plan', SEM)], {'plan': {'error': 'eating_out_conflict'}})


# ---------- cotización que rige ----------
def casos_cotizacion():
    cot = [quote('arroz', 'walmart', 'arroz.1kg', 2800, 1000000), quote('arroz', 'bodega', 'arroz.750g', 2200, 750000, rec_=False),
           quote('frijol', 'walmart', 'frijol.1kg', 3000, 1000000), quote('frijol', 'bodega', 'frijol.1kg', 3200, 1000000, rec_=False, fecha='2026-10-08'),
           quote('aceite-soya', 'walmart', 'aceite-soya.946ml', 3990, 946000),
           quote('harina', 'walmart', 'harina.1kg', 3000, 1000000), quote('harina', 'walmart', 'harina.500g', 1700, 500000, rec_=False)]
    h = ing('harina', 'harina.1kg', 100000)
    c1 = rec('combo', 'Combo con aceite', ['lunch'], [ing('arroz', 'arroz.1kg', 500000), ing('frijol', 'frijol.1kg', 250000), ing('aceite-soya', 'aceite-soya.946ml', 15000), h])
    c2 = rec('combo2', 'Combo sin aceite', ['lunch'], [ing('arroz', 'arroz.1kg', 500000), ing('frijol', 'frijol.1kg', 250000), h])
    d = datos(recipes=[c1, c2], quotes=cot,
              meal_plan=[entrada('2026-10-12', 'lunch', 'persona-a', 'combo2', 2000), entrada('2026-10-13', 'dinner', 'persona-b', 'combo', 2000)],
              list_choices=[eleccion('frijol', store='bodega'), eleccion('aceite-soya', store='bodega'), eleccion('harina', pres='harina.500g')])
    # arroz sin elección: rige la recomendada (Walmart 1 kg, 2800). frijol con tienda elegida (bodega, 3200). aceite con tienda elegida sin cotización: sin precio (no cae a la recomendada).
    # harina con presentación elegida (500 g, 1700). combo2 a 2000: 1400 + 800 + 340 = 2540 (por porción 1270); a 3000: 2100 + 1200 + 510 = 3810 (1270)
    q = lambda n, a: (n, 'recipe_detail', a)
    item = lambda p, c, k=None: {'product': p, 'cost_cents': c}
    esp = {'combo2-2': {'ingredients': [item('arroz', 1400), item('frijol', 800), item('harina', 340)], 'total_cost_cents': 2540, 'cost_per_serving_cents': 1270, 'missing_prices': []},
           'combo2-3': {'ingredients': [item('arroz', 2100), item('frijol', 1200), item('harina', 510)], 'total_cost_cents': 3810, 'cost_per_serving_cents': 1270, 'missing_prices': []},
           'combo-2': {'ingredients': [item('arroz', 1400), item('frijol', 800), item('aceite-soya', None), item('harina', 340)], 'total_cost_cents': None,
                       'cost_per_serving_cents': None, 'missing_prices': ['aceite-soya']},
           'plan': {'cells': [{'cost_cents': 2540, 'price_store': None, 'price_observed_date': '2026-10-08'}, {'cost_cents': None}],
                    'days': [{'cost_cents': 2540}, {'cost_cents': None}, {'cost_cents': 0}, {'cost_cents': 0}, {'cost_cents': 0}, {'cost_cents': 0}, {'cost_cents': 0}],
                    'week_cost_cents': None},
           'lista': {'items': [
               art('aceite-soya', 'aceite-soya.946ml', 15000, 1, 931000, None, ['combo'], store='bodega', pstore=None, rstore='walmart'),
               art('arroz', 'arroz.1kg', 1000000, 1, 0, 2800, ['combo', 'combo2']),
               art('frijol', 'frijol.1kg', 500000, 1, 500000, 3200, ['combo', 'combo2'], store='bodega', pstore='bodega', fecha='2026-10-08'),
               art('harina', 'harina.500g', 200000, 1, 300000, 1700, ['combo', 'combo2'], rpres='harina.1kg')],
               'groups': [grp('chain', 'walmart', None, 4500), grp('chain', 'bodega', None, 3200)], 'estimated_total_cents': 7700, 'missing_prices': ['aceite-soya']}}
    # recommended_store de frijol y aceite (tienda elegida): la tienda de la recomendada = walmart; de arroz y harina (sin tienda elegida) = price_store = walmart
    escribe('s4-costo-cotizacion-elegida', 'S4-06, §8.1, §6.5 (ILUSTRATIVO): sin elección rige la recomendada (arroz Walmart 1 kg, 2800; la cotización de Bodega 750 g se ignora); tienda elegida con cotización (frijol, Bodega 3200, observada el 2026-10-08); '
            'tienda elegida SIN cotización (aceite en Bodega: sin precio, no cae a la recomendada de Walmart); presentación elegida (harina 500 g, 1700 contra 3000 de 1 kg). combo2 a 2000: 1400+800+340 = 2540 (por porción 1270); '
            'a 3000: 2100+1200+510 = 3810. Celda con varias tiendas: price_store ausente y fecha la más vieja (2026-10-08). Lista: paquetes arroz 1 (1000000/1000000), frijol 1 (sobra 500000), '
            'harina 200000/500000 = 1 (sobra 300000); grupos walmart 2800+1700 = 4500, bodega 3200; total 7700.',
            [('criterio', 'S4-06'), ('criterio', 'S4-03'), ('manual', 'cuentas en descripcion')], d,
            [q('combo2-2', {'recipe': 'combo2', 'servings_milli': 2000}), q('combo2-3', {'recipe': 'combo2', 'servings_milli': 3000}), q('combo-2', {'recipe': 'combo', 'servings_milli': 2000}),
             ('plan', 'week_plan', SEM), ('lista', 'shopping_list', SEM)], esp)
    dos = [quote('arroz', 'walmart', 'arroz.1kg', 2800, 1000000), quote('arroz', 'bodega', 'arroz.750g', 2200, 750000)]
    escribe('s4-costo-dos-recomendadas', 'S4-06, §8.1: dos cotizaciones recomendadas del mismo producto es invalid_argument en toda consulta. ILUSTRATIVO.',
            [('criterio', 'S4-06'), ('manual', 'cuentas en descripcion')], datos(quotes=dos, meal_plan=[]),
            [('detalle', 'recipe_detail', {'recipe': 'arroz-con-frijol', 'servings_milli': 2000})], {'detalle': {'error': 'invalid_argument'}})


# ---------- canal y mercados (PROVISIONAL) ----------
def casos_canal():
    guiso = rec('guiso', 'Guiso de prueba', ['lunch'], [ing('jitomate', 'jitomate.kilo', 200000), ing('fresa', 'fresa.kilo', 100000), ing('pollo', 'pollo.kilo', 300000),
                                                          ing('arroz', 'arroz.1kg', 200000), ing('frijol', 'frijol.1kg', 100000), ing('cebolla', 'cebolla.kilo', 80000)])
    stores = [{'code': 'tianguis-jueves', 'kind': 'market', 'groups': ['verduras', 'frutas'], 'next_visit': '2026-10-15'},
              {'code': 'mercado-martes', 'kind': 'market', 'groups': ['verduras'], 'next_visit': '2026-10-13'},
              {'code': 'mercado-del-barrio', 'kind': 'market', 'groups': ['frutas'], 'next_visit': '2026-10-15'}] + TIENDAS_CADENA
    cot = [quote('arroz', 'walmart', 'arroz.1kg', 2800, 1000000), quote('frijol', 'walmart', 'frijol.1kg', 3000, 1000000),
           quote('jitomate', 'walmart', 'jitomate.kilo', 2200, 1000000), quote('cebolla', 'bodega', 'cebolla.kilo', 1500, 1000000),
           quote('pollo', 'walmart', 'pollo.kilo', 12000, 1000000)]
    d = datos(recipes=[guiso], stores=stores, quotes=cot, meal_plan=[entrada('2026-10-12', 'lunch', 'persona-a', 'guiso', 2000)],
              list_choices=[eleccion('cebolla', store='bodega')])
    items = [art('arroz', 'arroz.1kg', 200000, 1, 800000, 2800, ['guiso']),
             art('cebolla', 'cebolla.kilo', 80000, 1, 920000, 1500, ['guiso'], store='bodega', pstore='bodega', rstore='bodega'),
             art('fresa', 'fresa.kilo', 100000, 1, 900000, None, ['guiso'], store='tianguis-jueves', pstore=None, rstore=None),
             art('frijol', 'frijol.1kg', 100000, 1, 900000, 3000, ['guiso']),
             art('jitomate', 'jitomate.kilo', 200000, 1, 800000, 2200, ['guiso'], store='mercado-martes'),
             art('pollo', 'pollo.kilo', 300000, 1, 700000, 12000, ['guiso'], store=None)]
    grupos = [grp('market', 'mercado-martes', '2026-10-13', 2200), grp('market', 'tianguis-jueves', '2026-10-15', 0), grp('market', 'bodega', None, 1500),
              grp('market', None, None, 12000), grp('chain', 'walmart', None, 5800)]
    esp = {'lista': {'items': items, 'groups': grupos, 'estimated_total_cents': 21500, 'missing_prices': ['fresa']}}
    escribe('s4-canal-mercados', 'PROVISIONAL (S4-asignacion-mercado, opción A). S4-12, RN-16 (ILUSTRATIVO). Jitomate (verduras) lo surten tianguis-jueves (visita 2026-10-15) y mercado-martes (2026-10-13): gana el martes. '
            'Fresa (frutas): tianguis-jueves y mercado-del-barrio con la misma visita 2026-10-15: gana el primero en `stores`. Pollo (carnes): ningún mercado lo surte, tienda nula en la vertical market. '
            'Cebolla: la persona eligió Bodega (cadena) para un producto de mercado. Jitomate: la recomendada de S5 es Walmart, pero el artículo sigue en el mercado. '
            'Grupos: market (martes 2200), (jueves 0, fresa sin precio), (bodega 1500, sin día), (nula 12000), luego chain walmart 2800+3000 = 5800. Total 2200+0+1500+12000+5800 = 21500.',
            [('criterio', 'S4-12'), ('regla', 'RN-16'), ('manual', 'cuentas en descripcion')], d, [('lista', 'shopping_list', SEM)], esp,
            provisional={'bloqueo': BL_MERCADO, 'que': 'asignación de mercados, orden y total de groups[]'})


# ---------- recetas reales del recetario oficial (datos/semilla/recetario.json) ----------
def _json(ruta):
    with open(RAIZ / ruta, encoding='utf-8') as fh:
        return json.load(fh, parse_float=Decimal)


def real(code, archivo):
    """Receta del recetario oficial; `archivo` solo dice de qué borrador (recetario-a o recetario-b) viene."""
    for r in _json('datos/semilla/recetario.json')['recetas']:
        if r['id'] == code:
            return {'code': code, 'name': r['nombre'], 'slots': r['tiempos'], 'base_servings': r['porciones_base'], 'prep_minutes': r['minutos_preparacion'],
                    'cook_minutes': r['minutos_coccion'], 'origin': 'claude_seed', 'is_active': True,
                    'ingredients': [ing(i['producto_id'], i['presentacion_id'], i['cantidad_milli'], i['opcional']) for i in r['ingredientes']]}
    raise SystemExit(f'receta {code} ({archivo}) no está en datos/semilla/recetario.json')


REALES = [('huevos-a-la-mexicana', 'recetario-a'), ('huevo-con-jamon-y-frijoles', 'recetario-a'), ('avena-con-manzana-y-canela', 'recetario-a'),
          ('yogur-griego-con-fresas-y-avena', 'recetario-a'), ('espagueti-con-atun-y-elote', 'recetario-b')]
TABLA_CSV = 'datos/semilla/nutricion/tabla.csv'
BL_EMPATES = 'docs/estado/bloqueos/S4-confirmacion-empates.md'
FUENTE = {'USDA': 'usda_fdc', 'OFF': 'open_food_facts', 'manual': 'manual'}
ESTADO = {'propuesto': 'proposed', 'confirmado': 'confirmed'}


def tabla():
    """tabla.csv en modo review (§9.1): {producto: {kcal (kcal por 100 g, fracción exacta), fuente, id_fuente, estado}}."""
    with open(RAIZ / TABLA_CSV, encoding='utf-8', newline='') as fh:
        return {f['producto_id']: {'kcal': Fraction(Decimal(f['kcal_100g'])), 'fuente': f['fuente'], 'id_fuente': f['id_fuente'], 'estado': f['estado']}
                for f in csv.DictReader(fh)}


def nut(prod, tab):
    """Renglón en línea copiado tal cual de tabla.csv: mismo valor, fuente, id_fuente y estado (nada se marca confirmado si no lo está)."""
    t = tab[prod]
    return {'product': prod, 'presentation': None, 'kcal_milli_per_100g': int(t['kcal'] * 1000), 'source': FUENTE[t['fuente']],
            'source_ref': t['id_fuente'], 'status': ESTADO[t['estado']]}


def cita(prods, tab):
    return '; '.join(f"{p} {tab[p]['kcal']} kcal/100 g ({tab[p]['fuente']} {tab[p]['id_fuente']}, {tab[p]['estado']})" for p in prods)


def provisional_empates(tab, prods):
    if all(tab[p]['estado'] == 'confirmado' for p in prods):
        return None
    return {'bloqueo': BL_EMPATES, 'que': 'los renglones de tabla.csv de estas recetas están en estado propuesto; si al confirmarlos se corrige un valor, se regeneran el recetario y esta fixture'}


def catalogo_real():
    d = _json('datos/semilla/productos.json')
    res = {}
    for p in d['productos']:
        g = p['gramos_por_unidad']
        res[p['id']] = {'unit': {'g': 'g', 'ml': 'ml', 'pieza': 'piece'}[p['unidad_base']], 'gpum': None if g is None else int(Decimal(g) * 1000)}
    return res


def kcal_independiente(receta, kcal100, cat, serv_milli):
    """Segunda vía (solo fractions): milésimas de kcal exactas de la receta para `serv_milli`."""
    tot = Fraction(0)
    for i in receta['ingredients']:
        cant = Fraction(i['quantity_milli'] * serv_milli, receta['base_servings'] * 1000)
        c = cat[i['product']]
        gramos = cant if c['unit'] == 'g' else cant * c['gpum'] / 1000
        tot += gramos * kcal100[i['product']] * 1000 / 100000
    return tot


def hu(x):  # HALF_UP de una fracción no negativa
    x = Fraction(x)
    return (2 * x.numerator + x.denominator) // (2 * x.denominator)


def casos_kcal():
    cat = catalogo_real()
    tab = tabla()
    kcal100 = {p: t['kcal'] for p, t in tab.items()}
    recetas = [real(c, a) for c, a in REALES]
    prods = []
    for r in recetas:
        for i in r['ingredients']:
            if i['product'] not in prods:
                prods.append(i['product'])
    consultas, esp = [], {}
    for r in recetas:
        for serv, tag in ((2000, '2'), (4000, '4')):
            k = kcal_independiente(r, kcal100, cat, serv)
            idc = f"{r['code']}-{tag}"
            consultas.append((idc, 'recipe_detail', {'recipe': r['code'], 'servings_milli': serv}))
            esp[idc] = {'recipe': r['code'], 'servings_milli': serv,
                        'ingredients': [{'product': i['product'], 'quantity': Q(Fraction(i['quantity_milli'] * serv, r['base_servings'] * 1000))} for i in r['ingredients']],
                        'total_kcal': hu(k / 1000), 'kcal_per_serving': hu(k / 1000 / Fraction(serv, 1000)), 'missing_nutrition': []}
    consultas.append(('cobertura', 'nutrition_coverage', {}))
    conf = sum(1 for p in prods if tab[p]['estado'] == 'confirmado')
    esp['cobertura'] = {'unique': len(prods), 'with_source': len(prods), 'confirmed': conf, 'proposed': len(prods) - conf, 'missing': []}
    d = {'today': '2026-10-11', 'catalog_file': 'datos/semilla/productos.json', 'recipes': recetas, 'persons': PERSONAS, 'nutrition_file': TABLA_CSV}
    escribe('s4-kcal-recetas-autor', 'S4-08 (prueba del AUTOR, no sustituye a S4-09): 5 recetas REALES de datos/semilla/recetario-a.json (huevos-a-la-mexicana, huevo-con-jamon-y-frijoles, avena-con-manzana-y-canela, '
            'yogur-griego-con-fresas-y-avena) y recetario-b.json (espagueti-con-atun-y-elote), a 2 y a 4 porciones (escalado 2->4 exacto), con nutrition_file = ' + TABLA_CSV + ' en modo review. Valores e id_fuente de tabla.csv: ' + cita(prods, tab) + '. '
            'Esperado calculado por una segunda vía independiente (fractions): kcal = Σ (gramos × kcal/100 g / 100), gramos de ml y pieza = cantidad × gramos_por_unidad; un solo redondeo HALF_UP al final.',
            [('criterio', 'S4-08'), ('regla', 'RN-17'), ('manual', 'segunda vía en rec_genera_fixtures.py'), ('manual', TABLA_CSV)], d, consultas, esp,
            provisional_empates(tab, prods))

    # Sin dato: falta jamón de pavo; la tortilla está solo propuesta (se carga en modo review)
    r1, r2 = recetas[1], recetas[3]
    prods2 = []
    for r in (r1, r2):
        for i in r['ingredients']:
            if i['product'] not in prods2:
                prods2.append(i['product'])
    con_dato = [p for p in prods2 if p != 'jamon-pavo']
    nutri = [nut(p, tab) for p in con_dato]
    conf2 = sum(1 for p in con_dato if tab[p]['estado'] == 'confirmado')
    k2 = kcal_independiente(r2, kcal100, cat, 2000)
    d = {'today': '2026-10-11', 'catalog_file': 'datos/semilla/productos.json', 'recipes': [r1, r2], 'persons': PERSONAS, 'nutrition': nutri,
         'meal_plan': [entrada('2026-10-12', 'lunch', 'persona-a', 'huevo-con-jamon-y-frijoles', 2000)]}
    esp = {'incompleta': {'total_kcal': None, 'kcal_per_serving': None, 'missing_nutrition': ['jamon-pavo'],
                          'total_cost_cents': None, 'missing_prices': ['huevo', 'jamon-pavo', 'frijoles-refritos', 'tortilla-maiz', 'aceite-soya', 'sal']},
           'completa': {'total_kcal': hu(k2 / 1000), 'kcal_per_serving': hu(k2 / 1000 / 2), 'missing_nutrition': []},
           'cobertura': {'unique': len(prods2), 'with_source': len(con_dato), 'confirmed': conf2, 'proposed': len(con_dato) - conf2, 'missing': ['jamon-pavo']},
           'plan': {'cells': [{'kcal': None, 'cost_cents': None}],
                    'days': [{'persons': [{'person': 'persona-a', 'kcal': None}, {'person': 'persona-b', 'kcal': 0}]}, {}, {}, {}, {}, {}, {}]}}
    escribe('s4-kcal-sin-dato', 'S4-08, S4-11, §9.1 y §9.3: sin dato nutricional de jamon-pavo, huevo-con-jamon-y-frijoles queda incompleta (kcal nula, nunca 0, y missing_nutrition); la celda y el día de esa persona quedan sin kcal; '
            'sin cotizaciones el costo es nulo y missing_prices lista los seis productos. yogur-griego-con-fresas-y-avena sí se calcula. Los renglones nutricionales van en línea, copiados tal cual de ' + TABLA_CSV +
            ' (valor, fuente, id_fuente y estado): ' + cita(con_dato, tab) + '. ILUSTRATIVO solo en el renglón de jamon-pavo, que se omite a propósito para simular un producto sin dato '
            '(tabla.csv cubre los 150 productos de comida del catálogo, y nutrition_file junto con nutrition es invalid_argument por §12.1). Cobertura: únicos = ' + str(len(prods2)) +
            ', con fuente ' + str(len(con_dato)) + ' (confirmados ' + str(conf2) + ', propuestos ' + str(len(con_dato) - conf2) + ') y 1 sin dato.',
            [('criterio', 'S4-08'), ('criterio', 'S4-11'), ('manual', 'segunda vía en rec_genera_fixtures.py'), ('manual', TABLA_CSV)], d,
            [('incompleta', 'recipe_detail', {'recipe': 'huevo-con-jamon-y-frijoles', 'servings_milli': 2000}),
             ('completa', 'recipe_detail', {'recipe': 'yogur-griego-con-fresas-y-avena', 'servings_milli': 2000}),
             ('cobertura', 'nutrition_coverage', {}), ('plan', 'week_plan', SEM)], esp, provisional_empates(tab, con_dato))


def casos_costo_recetas():
    c = 'datos/semilla/productos.json'
    esp_i = rec('espagueti-ilustrativo', 'Espagueti con atún (ilustrativa)', ['lunch', 'dinner'],
                [ing('espagueti', 'espagueti.la-moderna-250g', 250000), ing('atun-agua', 'atun-agua.dolores-140g', 140000),
                 ing('elote-grano', 'elote-grano.herdez-220g', 110000), ing('aceite-oliva', 'aceite-oliva.nutrioli-oli-500ml', 15000)], origin='household')
    yog = rec('yogur-con-avena-ilustrativa', 'Yogur con avena (ilustrativa)', ['breakfast'],
              [ing('yogur-griego', 'yogur-griego.yoplait-750g', 300000), ing('avena', 'avena.3-minutos-400g', 40000)], origin='household')
    fri = rec('frijoles-con-jamon-ilustrativa', 'Frijoles con jamón (ilustrativa)', ['dinner'],
              [ing('frijoles-refritos', 'frijoles-refritos.la-costena-220g', 220000), ing('jamon-pavo', 'jamon-pavo.sabori-400g', 80000),
               ing('aceite-soya', 'aceite-soya.nutrioli-946ml', 10000)], origin='household')
    real_ = real('espagueti-con-atun-y-elote', 'recetario-b')
    b = lambda p, pres, price, size: quote(p, 'bodega', pres, price, size)
    cot = [b('espagueti', 'espagueti.la-moderna-250g', 1200, 250000), b('atun-agua', 'atun-agua.dolores-140g', 2100, 140000),
           b('elote-grano', 'elote-grano.herdez-220g', 1350, 220000), b('aceite-oliva', 'aceite-oliva.nutrioli-oli-500ml', 15000, 500000),
           b('yogur-griego', 'yogur-griego.yoplait-750g', 7500, 750000), b('avena', 'avena.3-minutos-400g', 3600, 400000),
           b('frijoles-refritos', 'frijoles-refritos.la-costena-220g', 1200, 220000), b('jamon-pavo', 'jamon-pavo.sabori-400g', 6200, 400000),
           b('aceite-soya', 'aceite-soya.nutrioli-946ml', 4800, 946000)]
    plan = [entrada('2026-10-12', 'lunch', 'persona-a', 'espagueti-ilustrativo', 2000), entrada('2026-10-12', 'lunch', 'persona-b', 'espagueti-ilustrativo', 2000),
            entrada('2026-10-12', 'dinner', 'persona-a', 'frijoles-con-jamon-ilustrativa', 2000), entrada('2026-10-12', 'dinner', 'persona-b', 'yogur-con-avena-ilustrativa', 2000),
            entrada('2026-10-13', 'breakfast', 'persona-a', 'yogur-con-avena-ilustrativa', 2000),
            entrada('2026-10-19', 'lunch', 'persona-a', 'espagueti-con-atun-y-elote', 2000)]
    d = {'today': '2026-10-11', 'catalog_file': c, 'recipes': [esp_i, yog, fri, real_], 'persons': PERSONAS,
         'stores': [{'code': 'bodega', 'kind': 'chain', 'groups': [], 'next_visit': None}], 'meal_plan': plan, 'quotes': cot}
    # espagueti: 1200 + 2100 + 110000*1350/220000 = 675 + 15000*15000/500000 = 450 -> 4425; por porción HALF_UP(2212.5) = 2213; a 4 porciones 8850 (2212.5 -> 2213)
    # yogur: 300000*7500/750000 = 3000 + 40000*3600/400000 = 360 -> 3360 (1680 por porción)
    # frijoles: 220000*1200/220000 = 1200 + 80000*6200/400000 = 1240 + 10000*4800/946000 = 50.7399 -> 2490.7399 -> 2491 (1245.37 -> 1245)
    # día 12: 4425 + 4425 + 2491 + 3360 = 14701; día 13: 3360; semana 18061. Semana 43: la receta real tiene 5 ingredientes sin precio: todo nulo.
    ei = lambda p, co: {'product': p, 'cost_cents': co}
    esp = {'espagueti-2': {'ingredients': [ei('espagueti', 1200), ei('atun-agua', 2100), ei('elote-grano', 675), ei('aceite-oliva', 450)], 'total_cost_cents': 4425,
                           'cost_per_serving_cents': 2213, 'missing_prices': []},
           'espagueti-4': {'ingredients': [ei('espagueti', 2400), ei('atun-agua', 4200), ei('elote-grano', 1350), ei('aceite-oliva', 900)], 'total_cost_cents': 8850,
                           'cost_per_serving_cents': 2213, 'missing_prices': []},
           'yogur-2': {'total_cost_cents': 3360, 'cost_per_serving_cents': 1680, 'missing_prices': []},
           'frijoles-2': {'ingredients': [ei('frijoles-refritos', 1200), ei('jamon-pavo', 1240), ei('aceite-soya', 51)], 'total_cost_cents': 2491,
                          'cost_per_serving_cents': 1245, 'missing_prices': []},
           'real-2': {'ingredients': [ei('espagueti', 1200), ei('atun-agua', 2100), ei('elote-grano', 675), ei('jitomate-saladet', None), ei('cebolla-blanca', None), ei('ajo', None),
                                      ei('aceite-oliva', 450), ei('sal', None), ei('oregano', None)], 'total_cost_cents': None, 'cost_per_serving_cents': None,
                      'missing_prices': ['jitomate-saladet', 'cebolla-blanca', 'ajo', 'sal', 'oregano']},
           'semana-42': {'days': [{'date': '2026-10-12', 'cost_cents': 14701}, {'date': '2026-10-13', 'cost_cents': 3360}] + [{'cost_cents': 0}] * 5, 'week_cost_cents': 18061},
           'semana-43': {'cells': [{'cost_cents': None}], 'days': [{'date': '2026-10-19', 'cost_cents': None}] + [{'cost_cents': 0}] * 6, 'week_cost_cents': None}}
    escribe('s4-costo-recetas', 'S4-06, §8.2, §8.3 (precios de Bodega Aurrera del ticket del 2026-10-10, data.json items_identicos[21], [16], [15], [0]; el resto ILUSTRATIVO): espagueti-ilustrativo 2 porciones = '
            '1200+2100+675+450 = 4425, por porción HALF_UP(2212.5) = 2213, a 4 porciones 8850 (§8.3); yogur con avena = 3000+360 = 3360; frijoles con jamón = 1200+1240+50.74 = 2491. '
            'La receta REAL espagueti-con-atun-y-elote (recetario-b) tiene 5 ingredientes sin precio: costo y totales nulos, nunca 0. Día 12 = 4425+4425+2491+3360 = 14701, día 13 = 3360, semana 42 = 18061; '
            'la semana 43 trae la receta incompleta y queda sin valor.',
            [('criterio', 'S4-06'), ('data.json', 'items_identicos'), ('manual', 'cuentas en descripcion')], d,
            [('espagueti-2', 'recipe_detail', {'recipe': 'espagueti-ilustrativo', 'servings_milli': 2000}), ('espagueti-4', 'recipe_detail', {'recipe': 'espagueti-ilustrativo', 'servings_milli': 4000}),
             ('yogur-2', 'recipe_detail', {'recipe': 'yogur-con-avena-ilustrativa', 'servings_milli': 2000}), ('frijoles-2', 'recipe_detail', {'recipe': 'frijoles-con-jamon-ilustrativa', 'servings_milli': 2000}),
             ('real-2', 'recipe_detail', {'recipe': 'espagueti-con-atun-y-elote', 'servings_milli': 2000}), ('semana-42', 'week_plan', SEM), ('semana-43', 'week_plan', {'week': '2026-W43'})], esp)


# ---------- salud ----------
def casos_salud():
    def p(code, nac, sexo, alto, peso, act, meta, **kw):
        return perfil(code, birth_date=nac, sex=sexo, height_mm=alto, activity=act, goal=meta, **kw), {'person': code, 'date': '2026-10-01', 'weight_g': peso}
    H, M = '1996-10-11', '1998-03-15'
    filas = [p('h30-mantener', H, 'male', 1750, 80000, 'sedentary', 'maintain'), p('h30-bajar-casa', H, 'male', 1750, 80000, 'sedentary', 'lose'),
             p('h30-bajar-propio', H, 'male', 1750, 80000, 'sedentary', 'lose', adjust_lose_kcal=100),
             p('h30-mantener-ajustes', H, 'male', 1750, 80000, 'sedentary', 'maintain', adjust_lose_kcal=300, adjust_gain_kcal=300),
             p('m28-mantener', M, 'female', 1650, 60000, 'sedentary', 'maintain'), p('m28-subir-casa', M, 'female', 1650, 60000, 'sedentary', 'gain'),
             p('h35-rn17', '1991-06-01', 'male', 1750, 72500, 'moderate', 'maintain'),
             p('h35-rn17-bajar', '1991-06-01', 'male', 1750, 72500, 'moderate', 'lose', adjust_lose_kcal=500),
             p('h29-vispera', '1996-10-12', 'male', 1750, 80000, 'sedentary', 'maintain'),
             p('meta-escrita', H, 'male', 1750, 80000, 'sedentary', 'maintain', daily_kcal_target=1900),
             p('imc-bajo', None, None, 2000, 70000, None, None), p('imc-18-5', None, None, 2000, 74000, None, None),
             p('imc-25-0', None, None, 2000, 100000, None, None), p('imc-30-0', None, None, 2000, 120000, None, None)]
    perfiles = [f[0] for f in filas] + [perfil('sin-datos')]
    pesos = [{'person': 'h30-mantener', 'date': '2026-01-05', 'weight_g': 85000}] + [f[1] for f in filas]
    personas = [{'code': pr['person'], 'has_meal_plan': True} for pr in perfiles]
    # MSJ (milésimas de kcal) = 10*peso_g + 625*estatura_mm - 5000*edad + 5000 (hombre) o -161000 (mujer); meta = MSJ * factor +- ajuste*1000, un solo redondeo
    # h30: 800000 + 1093750 - 150000 + 5000 = 1748750 -> x6/5 = 2098500 -> 2098.5 = 4197/2 -> 2099
    # m28: 600000 + 1031250 - 140000 - 161000 = 1330250 -> x6/5 = 1596300 -> 15963/10 -> 1596
    # casa 300/250: h30 bajar 2098.5 - 300 = 1798.5 -> 1799; m28 subir 1596.3 + 250 = 1846.3 -> 1846; h30 con 100 propios 1998.5 -> 1999; mantener ignora ajustes -> 2099
    # h35 (RN-17): 725000 + 1093750 - 175000 + 5000 = 1648750 x31/20 = 2555562.5 -> 40889/16 -> 2556; bajar 500: 2055.5625 -> 32889/16 -> 2056
    # h29 (víspera de cumpleaños): 800000 + 1093750 - 145000 + 5000 = 1753750 x6/5 = 2104500 -> 4209/2 -> 2105
    # IMC = HALF_UP(peso_g * 10000 / estatura_mm^2): 80000 -> 800000000/3062500 = 261.22 -> 261; 72500 -> 236.73 -> 237; 2000 mm: 70000 -> 175; 74000 -> 185; 100000 -> 250; 120000 -> 300
    s = lambda code, edad, imc, clase, ex, sug, rige, est: {'person': code, 'age': edad, 'bmi_tenths': imc, 'bmi_class': clase,
                                                           'suggested_kcal_exact': None if ex is None else Q(*ex), 'suggested_kcal_target': sug, 'daily_kcal_target': rige, 'target_status': est}
    esp = {'h30-mantener': s('h30-mantener', 30, 261, 'overweight', (4197, 2), 2099, 2099, 'suggested'),
           'h30-bajar-casa': s('h30-bajar-casa', 30, 261, 'overweight', (3597, 2), 1799, 1799, 'suggested'),
           'h30-bajar-propio': s('h30-bajar-propio', 30, 261, 'overweight', (3997, 2), 1999, 1999, 'suggested'),
           'h30-mantener-ajustes': s('h30-mantener-ajustes', 30, 261, 'overweight', (4197, 2), 2099, 2099, 'suggested'),
           'm28-mantener': s('m28-mantener', 28, 220, 'normal', (15963, 10), 1596, 1596, 'suggested'),
           'm28-subir-casa': s('m28-subir-casa', 28, 220, 'normal', (18463, 10), 1846, 1846, 'suggested'),
           'h35-rn17': s('h35-rn17', 35, 237, 'normal', (40889, 16), 2556, 2556, 'suggested'),
           'h35-rn17-bajar': s('h35-rn17-bajar', 35, 237, 'normal', (32889, 16), 2056, 2056, 'suggested'),
           'h29-vispera': s('h29-vispera', 29, 261, 'overweight', (4209, 2), 2105, 2105, 'suggested'),
           'meta-escrita': s('meta-escrita', 30, 261, 'overweight', (4197, 2), 2099, 1900, 'confirmed'),
           'imc-bajo': s('imc-bajo', None, 175, 'underweight', None, None, None, 'none'), 'imc-18-5': s('imc-18-5', None, 185, 'normal', None, None, None, 'none'),
           'imc-25-0': s('imc-25-0', None, 250, 'overweight', None, None, None, 'none'), 'imc-30-0': s('imc-30-0', None, 300, 'obesity', None, None, None, 'none'),
           'sin-datos': s('sin-datos', None, None, None, None, None, None, 'none')}
    # m28: 60000*10000/1650^2 = 600000000/2722500 = 220.39 -> 220
    consultas = [(p_['person'], 'health_profile', {'person': p_['person']}) for p_ in perfiles]
    base = {'today': '2026-10-11', 'catalog': {'groups': [], 'products': [], 'presentations': []}, 'recipes': [], 'persons': personas,
            'health_profiles': perfiles, 'weights': pesos}
    d1 = dict(base, settings={'house_adjust_lose_kcal': 300, 'house_adjust_gain_kcal': 250})
    escribe('s4-salud-perfiles', 'S4-07, RN-17, K7, K23, K36 (personas ficticias; hoy 2026-10-11; casa: bajar 300, subir 250 como fixture, no canon). Cálculo a mano en `fuente`/descripción: '
            'hombre 30 años, 80 kg, 175 cm, sedentary: MSJ 1748.75, x1.2 = 2098.50 -> 2099 (mantener, con ajustes 300/300 también 2099); bajar con la casa 300: 1798.50 -> 1799; bajar con 100 propios: 1998.50 -> 1999; '
            'mujer 28 años, 60 kg, 165 cm: MSJ 1330.25, x1.2 = 1596.30 -> 1596; subir con la casa 250: 1846.30 -> 1846; RN-17: hombre 35 años, 72.5 kg, 175 cm, moderate: 1648.75 x1.55 = 2555.5625 -> 2556, bajar 500 propios 2056; '
            'víspera de cumpleaños (29 años): 2104.50 -> 2105 (HALF_UP); meta escrita 1900 reemplaza a la sugerida (2099); sin insumos: sin sugerencia (none). '
            'IMC (estatura 2000 mm): 70 kg 17.5 underweight; 74 kg 18.5 normal; 100 kg 25.0 overweight; 120 kg 30.0 obesity; 72.5 kg y 1750 mm 23.7 normal; 80 kg y 1750 mm 26.1; el peso vigente es el más reciente (80 kg, no 85).',
            [('criterio', 'S4-07'), ('regla', 'RN-17'), ('decision', 'K36'), ('manual', 'cuentas en la descripcion')], d1, consultas, esp)
    d0 = dict(base, health_profiles=[x for x in perfiles if x['person'] in ('h30-bajar-casa', 'm28-subir-casa')],
              persons=[x for x in personas if x['code'] in ('h30-bajar-casa', 'm28-subir-casa')], weights=[x for x in pesos if x['person'] in ('h30-bajar-casa', 'm28-subir-casa')])
    escribe('s4-salud-ajuste-inicial', 'S4-07, K36, K44: con los ajustes de la casa en su valor inicial (`settings` ausente = 0 y 0), el hombre con objetivo «bajar» sale 2098.50 -> 2099 y la mujer con objetivo «subir» '
            '1596.30 -> 1596, iguales al caso sin ajuste. ILUSTRATIVO.', [('criterio', 'S4-07'), ('decision', 'K36'), ('manual', 'cuentas en s4-salud-perfiles')], d0,
            [('h30-bajar-casa', 'health_profile', {'person': 'h30-bajar-casa'}), ('m28-subir-casa', 'health_profile', {'person': 'm28-subir-casa'})],
            {'h30-bajar-casa': esp['h30-mantener'] | {'person': 'h30-bajar-casa'}, 'm28-subir-casa': esp['m28-mantener'] | {'person': 'm28-subir-casa'}})
    dfr = dict(base, health_profiles=[perfil('franja', height_mm=2000)], persons=[{'code': 'franja', 'has_meal_plan': True}],
               weights=[{'person': 'franja', 'date': '2026-10-01', 'weight_g': 73990}])
    escribe('s4-salud-imc-franja', 'PROVISIONAL (S4-reglas-lista, pregunta 3, opción A): 73990 g y 2000 mm = 18.4975 exacto, se muestra 18.5 (bmi_tenths = 185) y se clasifica sobre el valor mostrado: normal. '
            'Con la opción B sería underweight. ILUSTRATIVO.', [('criterio', 'S4-07'), ('manual', 'cuentas en descripcion')], dfr,
            [('franja', 'health_profile', {'person': 'franja'})], {'franja': {'bmi_tenths': 185, 'bmi_class': 'normal'}},
            provisional={'bloqueo': BL_LISTA, 'que': 'clasificación del IMC sobre el valor mostrado (pregunta 3), opción A'})


def casos_errores():
    """Errores de entrada (§12.1) y de catálogo (§14, S4-12): cada uno en su propia fixture porque invalida todo el conjunto de datos."""
    vacio = {'groups': [], 'products': [], 'presentations': []}
    base = {'today': '2026-10-11', 'catalog': vacio, 'recipes': [], 'persons': []}
    consulta = [('r', 'inventory_status', {})]
    sin_hoy = {k: v for k, v in base.items() if k != 'today'}
    escribe('s4-entrada-campo-ausente', '§12.1: un campo obligatorio ausente (today) es missing_field. ILUSTRATIVO.', [('criterio', 'S4-15'), ('manual', '§12.1')], sin_hoy, consulta,
            {'r': {'error': 'missing_field'}})
    escribe('s4-entrada-campo-desconocido', '§12.1: un campo desconocido dentro de un objeto (persons[0].apodo) es unknown_field. ILUSTRATIVO.', [('criterio', 'S4-15'), ('manual', '§12.1')],
            dict(base, persons=[{'code': 'persona-a', 'has_meal_plan': True, 'apodo': 'x'}]), consulta, {'r': {'error': 'unknown_field'}})
    escribe('s4-entrada-fecha-invalida', '§12.1 y §12.2: 2026-02-30 no es fecha: invalid_date. ILUSTRATIVO.', [('criterio', 'S4-15'), ('manual', '§12.1')],
            dict(base, today='2026-02-30'), consulta, {'r': {'error': 'invalid_date'}})
    mal = rec('receta-mala', 'Receta con producto inexistente', ['lunch'], [ing('no-existe', 'arroz.1kg', 1000)])
    escribe('s4-entrada-producto-desconocido', '§14: una receta con un producto que no está en el catálogo es unknown_product. ILUSTRATIVO.', [('criterio', 'S4-01'), ('manual', '§14')],
            {'today': '2026-10-11', 'catalog': CAT, 'recipes': [mal], 'persons': PERSONAS}, consulta, {'r': {'error': 'unknown_product'}})
    cat = json.loads(json.dumps(CAT))
    cat['products'].append({'code': 'ajo', 'name': 'Ajo', 'base_unit': 'g', 'group': None, 'grams_per_unit_milli': None})
    cat['presentations'].append({'code': 'ajo.kilo', 'product': 'ajo', 'name': 'Ajo por kilo', 'size_milli': 1000000})
    r = rec('con-ajo', 'Receta con ajo', ['lunch'], [ing('ajo', 'ajo.kilo', 10000)])
    escribe('s4-canal-sin-grupo', 'S4-12, §6.5: un producto sin grupo no tiene canal: la lista falla con product_without_channel; el detalle de la receta, que no usa el canal, sí se calcula '
            '(10000 x 2000/2000 = 10000 g). ILUSTRATIVO.', [('criterio', 'S4-12'), ('manual', '§6.5')],
            {'today': '2026-10-11', 'catalog': cat, 'recipes': [r], 'persons': PERSONAS, 'meal_plan': [entrada('2026-10-12', 'lunch', 'persona-a', 'con-ajo', 2000)]},
            [('lista', 'shopping_list', SEM), ('detalle', 'recipe_detail', {'recipe': 'con-ajo', 'servings_milli': 2000})],
            {'lista': {'error': 'product_without_channel'}, 'detalle': {'ingredients': [{'product': 'ajo', 'quantity': Q(10000)}]}})


def caso_escalado_recetario():
    """Tabla de las 60 recetas × 4 escalas (S4-02): solo si ya existe datos/semilla/recetario.json (el integrador lo arma de recetario-a y recetario-b)."""
    ruta = RAIZ / 'datos' / 'semilla' / 'recetario.json'
    if not ruta.exists():
        print('s4-escalado-recetario: PENDIENTE (no existe datos/semilla/recetario.json)')
        return
    libro = _json('datos/semilla/recetario.json')['recetas']
    cs, esp = [], {}
    for r in libro:
        for serv in (4000, 2000 // 2, 3000, 6000):
            idc = f"{r['id']}-{serv}"
            cs.append((idc, 'scale_recipe', {'recipe': r['id'], 'servings_milli': serv}))
            unidad = {'g': 'g', 'ml': 'ml', 'pieza': 'piece'}
            esp[idc] = {'recipe': r['id'], 'servings_milli': serv,
                        'ingredients': [{'product': i['producto_id'], 'unit': unidad[i['unidad']], 'quantity': Q(Fraction(i['cantidad_milli'] * serv, r['porciones_base'] * 1000))}
                                        for i in r['ingredientes']]}
    # nutrition_file: el motor carga recipe_book_file en modo review y coteja kcal_por_porcion y fuentes_kcal con tabla.csv (§3.2)
    d = {'today': '2026-10-11', 'catalog_file': 'datos/semilla/productos.json', 'recipe_book_file': 'datos/semilla/recetario.json', 'nutrition_file': TABLA_CSV, 'persons': PERSONAS}
    escribe('s4-escalado-recetario', f'S4-02: las {len(libro)} recetas de datos/semilla/recetario.json a 4, 1, 3 y 6 porciones ({len(cs)} consultas), con fractions de Python. 2 -> 4 es exactamente el doble. '
            'El recetario se carga en modo review con ' + TABLA_CSV + '.',
            [('criterio', 'S4-02')], d, cs, esp)


def main():
    SALIDA.mkdir(parents=True, exist_ok=True)
    for f in (casos_escalado, casos_lista, casos_inventario, casos_plan, casos_cotizacion, casos_canal, casos_kcal, casos_costo_recetas, casos_salud, casos_errores, caso_escalado_recetario):
        f()
    print(f'{len(list(SALIDA.glob("*.json")))} fixtures en {SALIDA.relative_to(RAIZ)}')


if __name__ == '__main__':
    import sys
    sys.path.insert(0, str(Path(__file__).resolve().parent))
    main()
