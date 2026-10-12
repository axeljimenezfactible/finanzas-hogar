#!/usr/bin/env python3
"""Generador del recetario oficial datos/semilla/recetario.json (S4-interfaz.md §3; S4-01, S4-08).

Une los dos borradores de autoría (datos/semilla/recetario-a.json y recetario-b.json, 30 + 30 recetas), valida las
reglas cruzadas de §3.2 que el JSON Schema no expresa, quita los campos que §3.2 no define (`etiquetas` de los
borradores: no forman parte del recetario 1.0) y escribe los dos campos derivados de cada receta con
datos/semilla/nutricion/tabla.csv leída en modo `review` (renglones `propuesto` y `confirmado`, §9.1):
  - fuentes_kcal[]: un {producto_id, fuente, id_fuente} por ingrediente con renglón en tabla.csv, en el orden de ingredientes[];
  - kcal_por_porcion: HALF_UP(Σ exacto / (1000 × porciones_base)) en milésimas de kcal (§9.2 y §9.3), o null si algún
    ingrediente no tiene renglón o, en ml o pieza, no tiene gramos_por_unidad.
Aritmética exacta con fractions; un solo redondeo al final. Determinista; solo biblioteca estándar; no lee código Go.
Uso (desde la raíz del repositorio): python3 -I oraculo/rec_genera_recetario.py [--verifica]
  --verifica: no escribe; recalcula y compara contra el recetario.json existente (sale con 1 si difiere).
"""
import csv
import json
import sys
from decimal import Decimal
from fractions import Fraction
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[1]
SEMILLA = RAIZ / 'datos' / 'semilla'
BORRADORES = ('recetario-a.json', 'recetario-b.json')
CAMPOS = ('id', 'nombre', 'descripcion', 'tiempos', 'porciones_base', 'minutos_preparacion', 'minutos_coccion', 'origen',
          'ingredientes', 'pasos', 'kcal_por_porcion', 'fuentes_kcal')
CAMPOS_ING = ('producto_id', 'cantidad_milli', 'unidad', 'presentacion_id', 'texto', 'opcional')
SOLO_BORRADOR = ('etiquetas',)  # no están en §3.2: no pasan al recetario oficial
TOTAL = 60


def milli(texto, campo):
    """Decimal de texto (hasta 3 decimales) a milésimas enteras, sin flotantes."""
    v = Decimal(texto) * 1000
    if v != v.to_integral_value():
        raise SystemExit(f'{campo}: {texto} tiene más de 3 decimales')
    return int(v)


def hu(x):
    """ROUND_HALF_UP de una fracción no negativa."""
    x = Fraction(x)
    return (2 * x.numerator + x.denominator) // (2 * x.denominator)


def carga_catalogo():
    with open(SEMILLA / 'productos.json', encoding='utf-8') as fh:
        d = json.load(fh)
    cat = {}
    for p in d['productos']:
        g = p['gramos_por_unidad']
        cat[p['id']] = {'unidad': p['unidad_base'], 'gpum': None if g is None else milli(g, p['id'] + '.gramos_por_unidad'),
                        'presentaciones': {x['id'] for x in p['presentaciones']}}
    return cat


def carga_tabla():
    """tabla.csv en modo review: los dos estados cuentan; un renglón por producto."""
    with open(SEMILLA / 'nutricion' / 'tabla.csv', encoding='utf-8', newline='') as fh:
        filas = list(csv.DictReader(fh))
    tab, ingles = {}, set()
    for f in filas:
        if f['producto_id'] in tab:
            raise SystemExit(f"tabla.csv: producto repetido {f['producto_id']}")
        if f['estado'] not in ('propuesto', 'confirmado') or not f['id_fuente'] or f['fuente'] not in ('USDA', 'OFF', 'manual'):
            raise SystemExit(f"tabla.csv: renglón inválido {f['producto_id']}")
        tab[f['producto_id']] = {'kcal_milli': milli(f['kcal_100g'], f['producto_id'] + '.kcal_100g'), 'fuente': f['fuente'], 'id_fuente': f['id_fuente']}
        if f['nombre_fuente']:
            ingles.add(f['nombre_fuente'].casefold())
    return tab, ingles


def derivados(receta, cat, tab):
    """(kcal_por_porcion o None, fuentes_kcal) de §3.2 con las reglas de §9.2 y §9.3."""
    total, falta, fuentes = Fraction(0), False, []
    for i in receta['ingredientes']:
        fila = tab.get(i['producto_id'])
        if fila is None:
            falta = True
            continue
        fuentes.append({'producto_id': i['producto_id'], 'fuente': fila['fuente'], 'id_fuente': fila['id_fuente']})
        c = cat[i['producto_id']]
        if c['unidad'] == 'g':
            gramos_milli = Fraction(i['cantidad_milli'])
        elif c['gpum'] is None:
            falta = True
            continue
        else:
            gramos_milli = Fraction(i['cantidad_milli'] * c['gpum'], 1000)
        total += gramos_milli * fila['kcal_milli'] / 100000
    kcal = None if falta else hu(total / (1000 * receta['porciones_base']))
    return kcal, fuentes


def valida(recetas, cat, ingles):
    """Reglas cruzadas de §3.2 que el esquema no puede expresar; devuelve la lista de errores."""
    err, ids, noms = [], set(), set()
    for r in recetas:
        rid = r.get('id')
        if rid in ids:
            err.append(f'{rid}: id repetido')
        if r['nombre'].casefold() in noms:
            err.append(f'{rid}: nombre repetido ({r["nombre"]})')
        if r['nombre'].casefold() in ingles:
            err.append(f'{rid}: nombre igual a un nombre_fuente de tabla.csv')
        ids.add(rid)
        noms.add(r['nombre'].casefold())
        vistos = set()
        for i in r['ingredientes']:
            p = i['producto_id']
            if p in vistos:
                err.append(f'{rid}: {p} repetido')
            vistos.add(p)
            if p not in cat:
                err.append(f'{rid}: producto {p} no existe')
                continue
            if i['presentacion_id'] not in cat[p]['presentaciones']:
                err.append(f'{rid}: presentación {i["presentacion_id"]} no es de {p}')
            if i['unidad'] != cat[p]['unidad']:
                err.append(f'{rid}: unidad {i["unidad"]} distinta de la base de {p}')
            if 'al gusto' in i.get('texto', '').casefold():
                err.append(f'{rid}: «al gusto» en {p}')
    if len(recetas) != TOTAL:
        err.append(f'el recetario trae {len(recetas)} recetas y deben ser {TOTAL}')
    return err


def arma():
    cat = carga_catalogo()
    tab, ingles = carga_tabla()
    recetas = []
    for nombre in BORRADORES:
        with open(SEMILLA / nombre, encoding='utf-8') as fh:
            d = json.load(fh)
        extra = set(d) - {'version', 'fecha', 'catalogo', 'recetas'}
        if extra:
            raise SystemExit(f'{nombre}: llaves de raíz fuera de §3.1: {sorted(extra)}')
        for r in d['recetas']:
            ajenos = set(r) - set(CAMPOS) - set(SOLO_BORRADOR)
            if ajenos:
                raise SystemExit(f'{nombre}: {r.get("id")}: campos fuera de §3.2: {sorted(ajenos)}')
            n = {k: r[k] for k in CAMPOS if k in r and k not in ('kcal_por_porcion', 'fuentes_kcal')}
            n['ingredientes'] = [{k: i[k] for k in CAMPOS_ING if k in i} for i in r['ingredientes']]
            n['kcal_por_porcion'], n['fuentes_kcal'] = derivados(n, cat, tab)
            recetas.append(n)
    err = valida(recetas, cat, ingles)
    if err:
        raise SystemExit('errores:\n' + '\n'.join(err))
    return {'version': '1.0', 'fecha': '2026-10-11', 'catalogo': 'datos/semilla/productos.json', 'recetas': recetas}


def formatea(libro):
    """JSON legible y compacto: un campo por renglón y cada ingrediente, fuente o paso en un solo renglón."""
    def j(x):
        return json.dumps(x, ensure_ascii=False)
    L = ['{'] + [f'  {j(k)}: {j(libro[k])},' for k in ('version', 'fecha', 'catalogo')] + ['  "recetas": [']
    for n, r in enumerate(libro['recetas']):
        L.append('    {')
        campos = list(r.items())
        for m, (k, v) in enumerate(campos):
            coma = ',' if m < len(campos) - 1 else ''
            if k in ('ingredientes', 'pasos', 'fuentes_kcal') and v:
                L.append(f'      {j(k)}: [')
                L += [f'        {j(x)}' + (',' if t < len(v) - 1 else '') for t, x in enumerate(v)]
                L.append('      ]' + coma)
            else:
                L.append(f'      {j(k)}: {j(v)}{coma}')
        L.append('    }' + (',' if n < len(libro['recetas']) - 1 else ''))
    texto = '\n'.join(L + ['  ]', '}']) + '\n'
    assert json.loads(texto) == libro
    return texto


def main():
    libro = arma()
    texto = formatea(libro)
    ruta = SEMILLA / 'recetario.json'
    nulas = [r['id'] for r in libro['recetas'] if r['kcal_por_porcion'] is None]
    if '--verifica' in sys.argv[1:]:
        igual = ruta.exists() and ruta.read_text(encoding='utf-8') == texto
        print(f'recetario.json {"coincide" if igual else "DIFIERE"} con lo regenerado')
        sys.exit(0 if igual else 1)
    ruta.write_text(texto, encoding='utf-8')
    print(f'{len(libro["recetas"])} recetas en datos/semilla/recetario.json; kcal nulas: {len(nulas)} {nulas}')


if __name__ == '__main__':
    main()
