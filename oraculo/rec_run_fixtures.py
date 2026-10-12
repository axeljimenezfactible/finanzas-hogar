#!/usr/bin/env python3
"""Ejecutor del oráculo S4: corre las fixtures y escribe oraculo/salidas/s4.json (FORMATO.md §4; S4-interfaz.md §12.3).

Uso: python3 -I oraculo/rec_run_fixtures.py [--fixtures oraculo/fixtures/s4] [--out oraculo/salidas/s4.json] [--caso <nombre>]
Implementa el contrato sin leer el código Go. Imprime una línea de resumen y una `DIF` por hoja de `esperado` que el oráculo no reproduce.
"""
import argparse
import json
import sys
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(Path(__file__).resolve().parent))

from rec_comun import ErrorMotor, a_json, fecha, semana_a_lunes  # noqa: E402
import rec_calculo as rc  # noqa: E402
import rec_datos as rd  # noqa: E402
import rec_plan as rp  # noqa: E402


def _args(args, permitidos):
    if not isinstance(args, dict):
        raise ErrorMotor('invalid_argument', 'argumentos')
    for k in args:
        if k not in permitidos:
            raise ErrorMotor('unknown_field', f'argumentos.{k}')
    for k in permitidos:
        if args.get(k) is None:
            raise ErrorMotor('invalid_argument', f'argumentos.{k}')
    return args


def calcula(d, q):
    """d: Datos cargados o ErrorMotor (entrada inválida: toda consulta del caso sale con ese código)."""
    try:
        if isinstance(d, ErrorMotor):
            raise d
        f, a = q['funcion'], q.get('argumentos', {})
        if f == 'scale_recipe':
            _args(a, ('recipe', 'servings_milli'))
            r = rc.scale_recipe(d, a['recipe'], a['servings_milli'])
        elif f == 'default_week':
            _args(a, ('date',))
            r = rp.default_week(d, fecha(a['date'], 'argumentos.date'))
        elif f == 'unplanned_slots':
            _args(a, ('week',))
            r = rp.unplanned_slots(d, a['week'])
        elif f == 'week_plan':
            _args(a, ('week',))
            r = rp.week_plan(d, a['week'])
        elif f == 'shopping_list':
            _args(a, ('week',))
            r = rp.shopping_list(d, a['week'])
        elif f == 'inventory_status':
            _args(a, ())
            r = rc.inventory_status(d)
        elif f == 'recipe_detail':
            _args(a, ('recipe', 'servings_milli'))
            r = rc.recipe_detail(d, a['recipe'], a['servings_milli'])
        elif f == 'nutrition_coverage':
            _args(a, ())
            r = rc.nutrition_coverage(d)
        elif f == 'health_profile':
            _args(a, ('person',))
            r = rc.health_profile(d, a['person'])
        else:
            raise ErrorMotor('invalid_argument', f'funcion {f}')
        return a_json(r)
    except ErrorMotor as e:
        return {'error': e.code}


def merge_patch(base, parche):
    if not isinstance(parche, dict):
        return parche
    out = dict(base) if isinstance(base, dict) else {}
    for k, v in parche.items():
        if v is None:
            out.pop(k, None)
        else:
            out[k] = merge_patch(out.get(k), v)
    return out


def resuelve_datos(datos, carpeta):
    if isinstance(datos, dict) and '$ref' in datos:
        with open(Path(carpeta) / datos['$ref'], encoding='utf-8') as fh:
            base = json.load(fh)
        return merge_patch(base, datos['parche']) if 'parche' in datos else base
    return datos


def recorta(calc, esp):
    if isinstance(esp, dict):
        if not isinstance(calc, dict):
            return '__tipo__'
        return {k: (recorta(calc[k], v) if k in calc else '__falta__') for k, v in esp.items()}
    if isinstance(esp, list):
        if not isinstance(calc, list):
            return '__tipo__'
        return [recorta(c, esp[i]) if i < len(esp) else '__sobra__' for i, c in enumerate(calc)]
    if isinstance(calc, (dict, list)):
        return '__tipo__'
    return calc


def hojas(x, ruta):
    if isinstance(x, dict):
        for k, v in x.items():
            yield from hojas(v, f'{ruta}.{k}')
    elif isinstance(x, list):
        yield f'{ruta}.#', len(x)
        for i, v in enumerate(x):
            yield from hojas(v, f'{ruta}[{i}]')
    else:
        yield ruta, x


def igual(a, b):
    if isinstance(a, bool) or isinstance(b, bool):
        return isinstance(a, bool) and isinstance(b, bool) and a == b
    return type(a) is type(b) and a == b


def lee_fixtures(carpeta, caso=None):
    res = []
    for p in sorted(Path(carpeta).glob('*.json')):
        if p.name.startswith('_') or (caso and p.stem != caso):
            continue
        with open(p, encoding='utf-8') as fh:
            fx = json.load(fh)
        if fx.get('formato') != 'fixture/1' or fx.get('motor') != 's4' or fx.get('caso') != p.stem:
            raise SystemExit(f'formato inválido: {p}')
        res.append((p, fx))
    return res


def corre(carpeta, caso=None):
    casos, difs, n_ok, n_hojas = {}, [], 0, 0
    fixtures = lee_fixtures(carpeta, caso)
    for p, fx in fixtures:
        try:
            d = rd.cargar(resuelve_datos(fx['entradas']['datos'], p.parent))
        except ErrorMotor as e:
            d = e
        salida, ok = {}, True
        for q in fx['entradas']['consultas']:
            esp = fx['esperado'].get(q['id'])
            calc = calcula(d, q)
            salida[q['id']] = recorta(calc, esp) if esp is not None else calc
            if esp is not None:
                obt = dict(hojas(salida[q['id']], q['id']))
                for r, v in hojas(esp, q['id']):
                    n_hojas += 1
                    if r not in obt or not igual(obt[r], v):
                        ok = False
                        difs.append(f"DIF {fx['caso']} {r} esperado={json.dumps(v, ensure_ascii=False)} oraculo={json.dumps(obt.get(r, '__falta__'), ensure_ascii=False)}")
        casos[fx['caso']] = salida
        n_ok += ok
    return fixtures, casos, difs, n_ok, n_hojas


def main(argv=None):
    ap = argparse.ArgumentParser(description='Oráculo S4: corre las fixtures y escribe la salida salida/1')
    ap.add_argument('--fixtures', default=str(RAIZ / 'oraculo' / 'fixtures' / 's4'))
    ap.add_argument('--out', default=str(RAIZ / 'oraculo' / 'salidas' / 's4.json'))
    ap.add_argument('--caso')
    a = ap.parse_args(argv)
    out = Path(a.out).resolve()
    if RAIZ not in out.parents:
        print(f'--out debe quedar dentro del repositorio: {out}', file=sys.stderr)
        return 2
    fixtures, casos, difs, n_ok, n_hojas = corre(a.fixtures, a.caso)
    out.parent.mkdir(parents=True, exist_ok=True)
    with open(out, 'w', encoding='utf-8') as fh:
        fh.write(json.dumps({'formato': 'salida/1', 'motor': 's4', 'casos': casos}, ensure_ascii=False, indent=2) + '\n')
    print(f'{len(fixtures)} fixtures, {n_ok} coinciden, {n_hojas} hojas de esperado, {len(difs)} diferencias')
    for x in difs:
        print(x)
    return 0 if not difs else 1


if __name__ == '__main__':
    sys.exit(main())
