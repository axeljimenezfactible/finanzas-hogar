#!/usr/bin/env python3
"""Ejecutor del oráculo S5: corre las fixtures y escribe oraculo/salidas/s5.json (FORMATO.md §4; S5-interfaz.md §15.3).

Uso: python3 -I oraculo/pre_run_fixtures.py [--fixtures oraculo/fixtures/s5] [--out oraculo/salidas/s5.json] [--caso <nombre>]
Implementa el contrato sin leer el código Go. Corre desde la raíz del repositorio (las rutas *_file son relativas a ella).
Imprime una línea de resumen y una `DIF` por hoja de `esperado` que el oráculo no reproduce.
"""
import argparse
import json
import os
import sys
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(Path(__file__).resolve().parent))

from pre_comun import ErrorMotor, fecha  # noqa: E402
import pre_calculo as pc  # noqa: E402
import pre_datos as pd  # noqa: E402
import pre_profeco as pp  # noqa: E402


def _args(a, oblig, opc=()):
    if not isinstance(a, dict):
        raise ErrorMotor('invalid_argument', 'argumentos')
    for k in a:
        if k not in oblig and k not in opc:
            raise ErrorMotor('unknown_field', f'argumentos.{k}')
    for k in oblig:
        if a.get(k) is None:
            raise ErrorMotor('invalid_argument', f'argumentos.{k}')
    return a


def calcula(d, q):
    """d: Datos cargados o ErrorMotor (entrada inválida: toda consulta del caso sale con ese código)."""
    try:
        if isinstance(d, ErrorMotor):
            raise d
        f, a = q['funcion'], q.get('argumentos', {})
        if f == 'base_unit_price':
            _args(a, ('presentation', 'price_cents'))
            r = pc.base_unit_price(d, a['presentation'], a['price_cents'])
        elif f == 'freshness':
            _args(a, ('observed_date', 'source'))
            if a['source'] not in pd.FUENTES:
                raise ErrorMotor('invalid_argument', 'source')
            r = pc.freshness(d, fecha(a['observed_date'], 'argumentos.observed_date'), a['source'])
        elif f == 'current_price':
            _args(a, ('presentation', 'store'))
            r = pc.current_price(d, a['presentation'], a['store'])
        elif f == 'quote_for':
            _args(a, ('product',), ('store', 'presentation'))
            r = pc.quote_for(d, a['product'], a.get('store'), a.get('presentation'))
        elif f == 'compare_carts':
            _args(a, ('base', 'others'))
            r = pc.compare_carts(d, a['base'], a['others'])
        elif f == 'compare_basket':
            _args(a, ())
            r = pc.compare_basket(d)
        elif f == 'membership':
            _args(a, ('store',), ('savings_per_purchase_cents',))
            r = pc.membership(d, a['store'], a.get('savings_per_purchase_cents'))
        elif f == 'next_visit':
            _args(a, ('store', 'from'))
            r = pc.next_visit(d, a['store'], fecha(a['from'], 'argumentos.from'))
        elif f == 'stale_prices_for_plan':
            _args(a, ('products', 'assigned_stores', 'grocery_date'))
            r = pc.stale_prices_for_plan(d, a['products'], a['assigned_stores'], fecha(a['grocery_date'], 'argumentos.grocery_date'))
        elif f == 'profeco_import':
            _args(a, ('files', 'links_file', 'mode'))
            r = pp.profeco_import(d, a['files'], a['links_file'], a['mode'])
        elif f == 'ocr_parse_tsv':
            _args(a, ('tsv_file',))
            r = pp.ocr_parse_tsv(a['tsv_file'])
        elif f == 'assisted_can_enable':
            _args(a, ('store',))
            r = pc.assisted_can_enable(d, a['store'])
        else:
            raise ErrorMotor('invalid_argument', f'funcion {f}')
        return json.loads(json.dumps(r))
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
        if fx.get('formato') != 'fixture/1' or fx.get('motor') != 's5' or fx.get('caso') != p.stem:
            raise SystemExit(f'formato inválido: {p}')
        res.append((p, fx))
    return res


def corre(carpeta, caso=None):
    casos, difs, n_ok, n_hojas = {}, [], 0, 0
    fixtures = lee_fixtures(carpeta, caso)
    for p, fx in fixtures:
        try:
            d = pd.cargar(resuelve_datos(fx['entradas']['datos'], p.parent))
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
    ap = argparse.ArgumentParser(description='Oráculo S5: corre las fixtures y escribe la salida salida/1')
    ap.add_argument('--fixtures', default=str(RAIZ / 'oraculo' / 'fixtures' / 's5'))
    ap.add_argument('--out', default=str(RAIZ / 'oraculo' / 'salidas' / 's5.json'))
    ap.add_argument('--caso')
    a = ap.parse_args(argv)
    out = Path(a.out).resolve()
    if RAIZ not in out.parents:
        print(f'--out debe quedar dentro del repositorio: {out}', file=sys.stderr)
        return 2
    carpeta = str(Path(a.fixtures).resolve())
    os.chdir(RAIZ)
    fixtures, casos, difs, n_ok, n_hojas = corre(carpeta, a.caso)
    out.parent.mkdir(parents=True, exist_ok=True)
    with open(out, 'w', encoding='utf-8') as fh:
        fh.write(json.dumps({'formato': 'salida/1', 'motor': 's5', 'casos': casos}, ensure_ascii=False, indent=2) + '\n')
    print(f'{len(fixtures)} fixtures, {n_ok} coinciden, {n_hojas} hojas de esperado, {len(difs)} diferencias')
    for x in difs:
        print(x)
    return 0 if not difs else 1


if __name__ == '__main__':
    sys.exit(main())
