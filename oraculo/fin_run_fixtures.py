#!/usr/bin/env python3
"""Ejecutor del oráculo S2: corre las fixtures y escribe oraculo/salidas/s2.json (FORMATO.md §4; S2-interfaz.md §6).

Uso: python3 -I oraculo/fin_run_fixtures.py [--fixtures oraculo/fixtures/s2] [--out oraculo/salidas/s2.json] [--caso <nombre>]
Rutas explícitas relativas al directorio actual; por omisión, relativas a la raíz del repositorio. No escribe fuera del repo.
"""
import argparse
import copy
import json
import sys
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(Path(__file__).resolve().parent))

from fin_comun import ErrorMotor, a_json, como_datos, fecha, mes, obj, quincena  # noqa: E402
import fin_calendario as fc  # noqa: E402
import fin_devengo as dv  # noqa: E402
import fin_flujo as ff  # noqa: E402
import fin_metas as fm  # noqa: E402
import fin_presupuesto as fp  # noqa: E402
import fin_proyeccion as fpy  # noqa: E402
import fin_vistas as fv  # noqa: E402


# ---- argumentos (§3, línea «Argumentos» de cada función) -----------------
def _claves(args, permitidas, obligatorias=()):
    if not isinstance(args, dict):
        raise ErrorMotor('invalid_argument', 'argumentos')
    for k in args:
        if k not in permitidas:
            raise ErrorMotor('unknown_field', f'argumentos.{k}')
    for k in obligatorias:
        if args.get(k) is None:
            raise ErrorMotor('invalid_argument', f'argumentos.{k}')


def _uno(args, a, b):
    if (args.get(a) is None) == (args.get(b) is None):
        raise ErrorMotor('invalid_argument', f'argumentos: exactamente uno de {a}/{b}')


def _entero(v, ruta, lo=None, hi=None):
    if not isinstance(v, int) or isinstance(v, bool) or (lo is not None and v < lo) or (hi is not None and v > hi):
        raise ErrorMotor('invalid_argument', ruta)
    return v


def _meses(v, ruta):
    if not isinstance(v, list) or not v:
        raise ErrorMotor('invalid_argument', ruta)
    return [mes(x, f'{ruta}[{i}]') for i, x in enumerate(v)]


def _ventanas(v):
    if not isinstance(v, list) or not v:
        raise ErrorMotor('invalid_argument', 'argumentos.windows')
    res = []
    for i, w in enumerate(v):
        obj(w, 'Window', f'argumentos.windows[{i}]')
        a, b = fecha(w.get('start'), f'windows[{i}].start'), fecha(w.get('end'), f'windows[{i}].end')
        if b < a or not isinstance(w.get('label'), str):
            raise ErrorMotor('invalid_argument', f'argumentos.windows[{i}]')
        res.append((w['label'], a, b))
    return res


def _rango(args):
    a, b = fecha(args['start'], 'argumentos.start'), fecha(args['end'], 'argumentos.end')
    if b < a:
        raise ErrorMotor('invalid_argument', 'start > end')
    return a, b


def evalua(datos, funcion, args):
    """Traduce una consulta a la función Python y devuelve la salida (con fechas ya en ISO)."""
    if args is None:
        args = {}
    if funcion == 'calendario':
        _claves(args, ('date', 'now'))
        _uno(args, 'date', 'now')
        r = fc.calendario(datos, date=args.get('date'), now=args.get('now'))
    elif funcion == 'periodos':
        _claves(args, ('months', 'windows'))
        _uno(args, 'months', 'windows')
        D = como_datos(datos)
        r = dv.periodos(D, meses=_meses(args['months'], 'argumentos.months')) if args.get('months') is not None else dv.periodos(D, ventanas=_ventanas(args['windows']))
    elif funcion in ('horizonte', 'vista_horizonte'):
        _claves(args, ('start', 'end', 'months'), ('start', 'end', 'months'))
        a, b = _rango(args)
        m = _meses(args['months'], 'argumentos.months')
        r = dv.horizonte(datos, a, b, m) if funcion == 'horizonte' else fv.vista_horizonte(datos, a, b, m)
    elif funcion == 'deuda_viva':
        _claves(args, ('cut',), ('cut',))
        r = dv.deuda_viva(datos, fecha(args['cut'], 'argumentos.cut'))
    elif funcion == 'promedios':
        _claves(args, ('total_cents', 'start', 'end', 'month', 'basis'))
        if 'month' in args or 'basis' in args:
            if any(k in args for k in ('total_cents', 'start', 'end')):
                raise ErrorMotor('invalid_argument', 'argumentos: dos formas a la vez')
            _claves(args, ('month', 'basis'), ('month', 'basis'))
            if args['basis'] not in ('cash', 'accrued', 'accrued_groceries'):
                raise ErrorMotor('invalid_argument', 'argumentos.basis')
            como_datos(datos)
            r = dv.promedios_mes(datos, mes(args['month'], 'argumentos.month'), args['basis'])
        else:
            _claves(args, ('total_cents', 'start', 'end'), ('total_cents', 'start', 'end'))
            como_datos(datos)
            a, b = _rango(args)
            r = dv.promedios(_entero(args['total_cents'], 'argumentos.total_cents'), a, b)
    elif funcion == 'promedio_mensual':
        _claves(args, ('total_cents', 'month_count'), ('total_cents', 'month_count'))
        como_datos(datos)
        r = dv.promedio_mensual(_entero(args['total_cents'], 'argumentos.total_cents'), _entero(args['month_count'], 'argumentos.month_count', 1))
    elif funcion in ('valor_estable', 'vista_valor_estable'):
        _claves(args, ('items',))
        items = args.get('items')
        if items is not None and not isinstance(items, list):
            raise ErrorMotor('invalid_argument', 'argumentos.items')
        r = dv.valor_estable_de(datos, items) if funcion == 'valor_estable' else fv.vista_valor_estable(datos, items)
    elif funcion in ('papel', 'vista_papel', 'metas'):
        _claves(args, ())
        r = {'papel': dv.papel, 'vista_papel': fv.vista_papel, 'metas': fm.metas}[funcion](datos)
    elif funcion in ('presupuesto_mes', 'vista_periodo'):
        _claves(args, ('month',), ('month',))
        m = mes(args['month'], 'argumentos.month')
        r = fp.presupuesto_mes(datos, m) if funcion == 'presupuesto_mes' else fv.vista_periodo(datos, m)
    elif funcion == 'aplicar_cambio_presupuesto':
        _claves(args, ('row_id', 'month', 'amount_cents', 'scope'), ('row_id', 'month', 'amount_cents'))
        D = como_datos(datos)
        if args.get('scope') not in (None, 'onward', 'this_month'):
            raise ErrorMotor('invalid_argument', 'argumentos.scope')
        r = fp.aplicar_cambio_presupuesto(D.budgets, dict(row_id=_entero(args['row_id'], 'argumentos.row_id'), month=mes(args['month'], 'argumentos.month'),
                                                             amount_cents=_entero(args['amount_cents'], 'argumentos.amount_cents', 0), scope=args.get('scope')))
    elif funcion == 'plan_contra_real':
        _claves(args, ('from', 'count', 'category'), ('from', 'count'))
        r = fp.plan_contra_real(datos, mes(args['from'], 'argumentos.from'), _entero(args['count'], 'argumentos.count', 1, 12), args.get('category'))
    elif funcion == 'flujo_quincenal':
        _claves(args, ('from', 'count'), ('from', 'count'))
        r = ff.flujo_quincenal(datos, quincena(args['from'], 'argumentos.from'), _entero(args['count'], 'argumentos.count', 1, 24))
    elif funcion == 'proyeccion':
        _claves(args, ('from',))
        r = fpy.proyeccion(datos, mes(args['from'], 'argumentos.from') if args.get('from') is not None else None)
    elif funcion == 'vista_sensibilidad':
        _claves(args, ('items',), ('items',))
        if not isinstance(args['items'], list):
            raise ErrorMotor('invalid_argument', 'argumentos.items')
        r = fv.vista_sensibilidad(datos, args['items'])
    elif funcion == 'vista_interes_por_monto':
        _claves(args, ('purchase_id',), ('purchase_id',))
        r = fv.vista_interes_por_monto(datos, _entero(args['purchase_id'], 'argumentos.purchase_id'))
    else:
        raise ErrorMotor('invalid_argument', f'funcion desconocida: {funcion}')
    return a_json(r)


def calcula(datos, consulta):
    try:
        return evalua(datos, consulta.get('funcion'), consulta.get('argumentos'))
    except ErrorMotor as e:
        return {'error': e.code}


# ---- fixtures, $ref, parche (RFC 7396) y recorte (FORMATO §2, §4) --------
def merge_patch(base, parche):
    if not isinstance(parche, dict):
        return copy.deepcopy(parche)
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
        if fx.get('formato') != 'fixture/1' or fx.get('motor') != 's2' or fx.get('caso') != p.stem:
            raise SystemExit(f'formato inválido: {p}')
        res.append((p, fx))
    return res


def corre(carpeta, caso=None):
    casos, difs, n_ok = {}, [], 0
    fixtures = lee_fixtures(carpeta, caso)
    for p, fx in fixtures:
        datos = resuelve_datos(fx['entradas']['datos'], p.parent)
        salida, ok = {}, True
        for q in fx['entradas']['consultas']:
            esp = fx['esperado'].get(q['id'])
            calc = calcula(datos, q)
            salida[q['id']] = recorta(calc, esp) if esp is not None else calc
            if esp is not None:
                obt = dict(hojas(salida[q['id']], q['id']))
                for r, v in hojas(esp, q['id']):
                    if r not in obt or not igual(obt[r], v):
                        ok = False
                        difs.append(f"DIF {fx['caso']} {r} esperado={json.dumps(v, ensure_ascii=False)} oraculo={json.dumps(obt.get(r, '__falta__'), ensure_ascii=False)}")
        casos[fx['caso']] = salida
        n_ok += ok
    return fixtures, casos, difs, n_ok


def main(argv=None):
    ap = argparse.ArgumentParser(description='Oráculo S2: corre las fixtures y escribe la salida salida/1')
    ap.add_argument('--fixtures', default=str(RAIZ / 'oraculo' / 'fixtures' / 's2'))
    ap.add_argument('--out', default=str(RAIZ / 'oraculo' / 'salidas' / 's2.json'))
    ap.add_argument('--caso')
    a = ap.parse_args(argv)
    out = Path(a.out).resolve()
    if RAIZ not in out.parents:
        print(f'--out debe quedar dentro del repositorio: {out}', file=sys.stderr)
        return 2
    fixtures, casos, difs, n_ok = corre(a.fixtures, a.caso)
    out.parent.mkdir(parents=True, exist_ok=True)
    with open(out, 'w', encoding='utf-8') as fh:
        fh.write(json.dumps({'formato': 'salida/1', 'motor': 's2', 'casos': casos}, ensure_ascii=False, indent=2) + '\n')
    print(f'{len(fixtures)} fixtures, {n_ok} coinciden, {len(difs)} diferencias')
    for d in difs:
        print(d)
    return 0 if not difs else 1


if __name__ == '__main__':
    sys.exit(main())
