"""Oráculo S4 · escalado, cotización que rige, costo, kcal, inventario y salud.

Secciones de S4-interfaz.md: §4 escalado, §8 costo, §9 kcal, §7 inventario, §10 salud. Todo en Fraction; un redondeo al final.
Un valor «ausente / sin dato / sin precio» sale como null (supuesto S-O1 del reporte: el contrato dice «ausente» y «null» sin distinguirlos).
"""
from fractions import Fraction

from rec_comun import ErrorMotor, entero, half_up, piso


# ---- §4 escalado ----
def receta(d, code):
    if code not in d.recipes:
        raise ErrorMotor('invalid_argument', f'receta desconocida: {code}')
    return d.recipes[code]


def porciones(v, ruta='servings_milli'):
    if not isinstance(v, int) or isinstance(v, bool):
        raise ErrorMotor('invalid_argument', ruta)
    if v <= 0:
        raise ErrorMotor('invalid_servings', ruta)
    return v


def escala(d, code, serv):
    """[(ingrediente, cantidad escalada exacta en milésimas)] de TODOS los ingredientes (opcionales incluidos)."""
    r = receta(d, code)
    return [(i, Fraction(i['q'] * serv, r['base'] * 1000)) for i in r['ings']]


def scale_recipe(d, code, serv):
    serv = porciones(serv)
    return {'recipe': code, 'servings_milli': serv,
            'ingredients': [{'product': i['product'], 'unit': d.products[i['product']]['unit'], 'quantity': q} for i, q in escala(d, code, serv)]}


# ---- §8.1 cotización que rige ----
def cotiza_rige(d, producto, eleccion):
    qs = d.quotes.get(producto, [])
    if not qs:
        return None
    if eleccion is None or (eleccion['store'] is None and eleccion['pres'] is None):
        return next(q for q in qs if q['rec'])
    casan = [q for q in qs if (eleccion['store'] is None or q['store'] == eleccion['store'])
             and (eleccion['pres'] is None or q['pres'] == eleccion['pres'])]
    if not casan:
        return None
    if len(casan) == 1:
        return casan[0]
    rec = [q for q in casan if q['rec']]
    if len(rec) == 1:
        return rec[0]
    raise ErrorMotor('invalid_argument', f'{producto}: varias cotizaciones casan con la elección')


# ---- §8 costo y §9 kcal de una receta ----
def kcal_de(d, ing, cantidad):
    """Milésimas de kcal exactas del ingrediente, o None (sin dato)."""
    p = ing['product']
    n = d.nutrition.get((p, ing['presentation'])) or d.nutrition.get((p, None))
    if n is None:
        return None
    prod = d.products[p]
    if prod['unit'] == 'g':
        gramos = cantidad
    elif prod['gpum'] is None:
        return None
    else:
        gramos = cantidad * prod['gpum'] / 1000
    return Fraction(gramos) * n['kcal'] / 100000


def detalle(d, code, serv):
    """Cálculo exacto de una receta a `serv` porciones (milésimas): costos, kcal, faltantes y cotizaciones usadas."""
    r = receta(d, code)
    ings, usadas = [], []
    falta_precio, falta_nutri = [], []
    suma_c, suma_k = Fraction(0), Fraction(0)
    for i, q in escala(d, code, serv):
        if i['optional']:
            continue  # §5.4: opcional fuera de lista, costo y kcal
        p = i['product']
        cot = cotiza_rige(d, p, d.choices.get((p, 'calculated')))
        costo = None
        if cot is None:
            falta_precio.append(p)
        else:
            costo = q * Fraction(cot['price'], cot['size'])
            suma_c += costo
            usadas.append(cot)
        k = kcal_de(d, i, q)
        if k is None:
            falta_nutri.append(p)
        else:
            suma_k += k
        ings.append((p, q, costo, k))
    return {'recipe': code, 'base': r['base'], 'ings': ings, 'cost': None if falta_precio else suma_c,
            'kcal_milli': None if falta_nutri else suma_k, 'missing_prices': falta_precio, 'missing_nutrition': falta_nutri,
            'quotes': usadas}


def recipe_detail(d, code, serv):
    serv = porciones(serv)
    x = detalle(d, code, serv)
    n = Fraction(serv, 1000)  # porciones pedidas
    return {
        'recipe': code, 'servings_milli': serv,
        'ingredients': [{'product': p, 'quantity': q, 'cost_cents': None if c is None else half_up(c),
                         'kcal': None if k is None else half_up(k / 1000)} for p, q, c, k in x['ings']],
        'total_cost_cents': None if x['cost'] is None else half_up(x['cost']),
        'cost_per_serving_cents': None if x['cost'] is None else half_up(x['cost'] / n),
        'total_kcal': None if x['kcal_milli'] is None else half_up(x['kcal_milli'] / 1000),
        'kcal_per_serving': None if x['kcal_milli'] is None else half_up(x['kcal_milli'] / (1000 * n)),
        'missing_prices': x['missing_prices'], 'missing_nutrition': x['missing_nutrition']}


def nutrition_coverage(d):
    """§9.1: sobre los ingredientes de las recetas activas. Un producto cuenta una vez; su estado es el del dato genérico
    si existe y, si no, el del primer dato de entrada (supuesto S-O5)."""
    unicos = []
    for c in d.recipe_order:
        if d.recipes[c]['active']:
            for i in d.recipes[c]['ings']:
                if i['product'] not in unicos:
                    unicos.append(i['product'])
    conf = prop = 0
    faltan = []
    for p in unicos:
        n = d.nutrition.get((p, None)) or next((v for (pp, _), v in d.nutrition.items() if pp == p), None)
        if n is None:
            faltan.append(p)
        elif n['status'] == 'confirmed':
            conf += 1
        else:
            prop += 1
    return {'unique': len(unicos), 'with_source': conf + prop, 'confirmed': conf, 'proposed': prop, 'missing': sorted(faltan)}


# ---- §7 inventario ----
def existencia(it, dia):
    """Existencia estimada a `dia` (RN-15). Stock: lo vigente. Ritmo: contada − HALF_UP(ritmo × días / periodo), mínimo 0."""
    if it['mode'] == 'stock':
        return it['q']
    dias = max(0, (dia - it['counted']).days)
    return max(0, it['q'] - half_up(Fraction(it['rq'] * dias, it['rp'])))


def inventory_status(d):
    items, resumen = [], {'sufficient': 0, 'reorder': 0, 'out': 0}
    for it in d.inv:
        est = existencia(it, d.today)
        if est == 0:
            st = 'out'
        elif it['level'] is not None and est <= it['level']:
            st = 'reorder'
        else:
            st = 'sufficient'
        resumen[st] += 1
        rate = it['mode'] == 'consumption_rate'
        items.append({'code': it['code'], 'estimated_quantity_milli': est, 'status': st,
                      'weekly_consumption_milli': half_up(Fraction(it['rq'] * 7, it['rp'])) if rate else None,
                      'days_left': piso(Fraction(est * it['rp'], it['rq'])) if rate else None})
    return {'items': items, 'summary': resumen}


# ---- §10 salud ----
FACTOR = {'sedentary': Fraction(6, 5), 'light': Fraction(11, 8), 'moderate': Fraction(31, 20),
          'active': Fraction(69, 40), 'very_active': Fraction(19, 10)}


def edad(nac, hoy):
    return hoy.year - nac.year - ((hoy.month, hoy.day) < (nac.month, nac.day))


def ultimo_peso(d, persona):
    ws = [w for w in d.weights if w['person'] == persona]
    if not ws:
        return None
    return max(enumerate(ws), key=lambda t: (t[1]['date'], t[0]))[1]['g']


def clase_imc(tenths):
    return 'underweight' if tenths < 185 else 'normal' if tenths < 250 else 'overweight' if tenths < 300 else 'obesity'


def salud(d, persona):
    if persona not in d.person:
        raise ErrorMotor('invalid_argument', 'person')
    p = d.profiles.get(persona) or {'height': None, 'birth': None, 'sex': None, 'activity': None, 'goal': None,
                                    'adj_lose': None, 'adj_gain': None, 'target': None}
    peso = ultimo_peso(d, persona)
    años = edad(p['birth'], d.today) if p['birth'] else None
    bmi = half_up(Fraction(peso * 10000, p['height'] ** 2)) if peso and p['height'] else None
    faltan = [n for n, v in (('sex', p['sex']), ('birth_date', p['birth']), ('height_mm', p['height']), ('weight', peso),
                             ('activity', p['activity']), ('goal', p['goal'])) if v is None]
    exacta = None
    if not faltan:
        msj = 10 * peso + 625 * p['height'] - 5000 * años + (5000 if p['sex'] == 'male' else -161000)  # milésimas de kcal
        aj = 0
        if p['goal'] == 'lose':
            aj = -(p['adj_lose'] if p['adj_lose'] is not None else d.house_lose)
        elif p['goal'] == 'gain':
            aj = p['adj_gain'] if p['adj_gain'] is not None else d.house_gain
        exacta = (Fraction(msj) * FACTOR[p['activity']] + aj * 1000) / 1000  # kcal, antes del redondeo final
    sug = None if exacta is None else half_up(exacta)
    if p['target'] is not None:
        rige, est = p['target'], 'confirmed'
    elif sug is not None:
        rige, est = sug, 'suggested'
    else:
        rige, est = None, 'none'
    return {'person': persona, 'age': años, 'bmi_tenths': bmi, 'bmi_class': None if bmi is None else clase_imc(bmi),
            'suggested_kcal_exact': exacta, 'suggested_kcal_target': sug, 'daily_kcal_target': rige, 'target_status': est,
            'missing_inputs': faltan}


def health_profile(d, persona):
    return salud(d, persona)
