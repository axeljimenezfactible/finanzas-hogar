#!/usr/bin/env python3
"""Comparación genérica de fixtures de motores (oraculo/fixtures/FORMATO.md §5). Dueño: agente Python de S2.

Uso: python3 -I oraculo/compara.py <motor> [--fixtures D] [--oraculo F] [--motor F] [--modelo F] [--data F] [--sin-motor]
Compara, hoja por hoja, esperado, oráculo y motor; y las consultas con id ruta contra model.json o data.json.
Si el archivo del motor no existe y no se pasó --sin-motor, avisa y compara sin motor. Solo lee; no escribe nada.
Código de salida: 0 todo coincide; 1 diferencias u hojas sin dueño; 2 archivo faltante o formato inválido.
"""
import argparse
import json
import re
import sys
from decimal import Decimal
from fractions import Fraction
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[1]
RAMAS = {'s2': ('periodos', 'valor_estable', 'horizonte', 'papel'),
         's3': ('parametros', 'despensas', 'credito', 'sensibilidad', 'disposicion_por_despensa', 'interes_por_monto'),
         's4': (), 's5': ()}
RAICES = ('parametros', 'periodos', 'despensas', 'credito', 'valor_estable', 'sensibilidad', 'papel',
          'disposicion_por_despensa', 'horizonte', 'interes_por_monto')
CONTEOS = {'n', 'n_pagos', 'despensa', 'despensa_que_cubre', 'paquete'}
RAZONES = {'tasa', 'tasa_pct'}
FALTA = '__falta__'


class Formato(Exception):
    pass


def pasos(ruta):
    return [(n if n else int(i)) for n, i in re.findall(r'([^.\[\]]+)|\[(\d+)\]', ruta)]


def nodo(raiz, ruta):
    x = raiz
    for k in pasos(ruta):
        x = x[k]
    return x


def es_razon(x):
    return isinstance(x, dict) and {'num', 'den'} <= set(x) <= {'num', 'den', 'decimales'}


def hojas(x, ruta):
    if isinstance(x, dict) and not es_razon(x):
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


def clave_final(ruta):
    k = pasos(ruta)
    return k[-1] if k and isinstance(k[-1], str) else None


def valor_modelo(v, ruta):
    """Hoja de model.json/data.json en las unidades de las fixtures (FORMATO §3); None si no se compara."""
    if isinstance(v, (bool, str)) or v is None or isinstance(v, (dict, list)):
        return None
    k = clave_final(ruta)
    if k in CONTEOS:
        return int(v)
    if k in RAZONES:
        return ('razon', k, float(v), Decimal(v))
    x = Decimal(v) * 100
    if x != x.to_integral_value():
        raise Formato(f'{ruta}: {v} no es un número exacto de centavos')
    return int(x)


def coincide_modelo(mv, x):
    if isinstance(mv, tuple):
        _, k, f, dec = mv
        if not es_razon(x):
            return False
        fr = Fraction(x['num'], x['den']) * (100 if k == 'tasa_pct' else 1)
        if 'decimales' in x:
            q = Decimal(1).scaleb(-x['decimales'])
            return (Decimal(fr.numerator) / Decimal(fr.denominator)).quantize(q, rounding='ROUND_HALF_UP') == dec
        return float(fr) == f
    return igual(x, mv)


def carga(ruta, que):
    try:
        with open(ruta, encoding='utf-8') as fh:
            return json.load(fh, parse_float=Decimal) if que in ('modelo', 'data') else json.load(fh)
    except FileNotFoundError:
        raise Formato(f'falta {que}: {ruta}') from None
    except json.JSONDecodeError as e:
        raise Formato(f'{que} no es JSON válido: {ruta}: {e}') from None


def raiz_de(i):
    if i.startswith('data.json:'):
        return 'data', i[len('data.json:'):]
    if re.match(r'^(%s)($|[.\[])' % '|'.join(RAICES), i):
        return 'modelo', i
    if re.fullmatch(r'[a-z0-9][a-z0-9_-]*', i):
        return None, i
    raise Formato(f'id inválido: {i}')


def compara(a):
    carpeta = Path(a.fixtures)
    if not carpeta.is_dir():
        raise Formato(f'falta la carpeta de fixtures: {carpeta}')
    modelo, data = carga(a.modelo, 'modelo'), carga(a.data, 'data')
    fuentes = {'modelo': modelo, 'data': data}
    ora = carga(a.oraculo, 'oráculo')
    if ora.get('formato') != 'salida/1' or ora.get('motor') != a.motor_nombre:
        raise Formato(f'salida del oráculo con formato o motor inválido: {a.oraculo}')
    mot, aviso = None, None
    if not a.sin_motor:
        if Path(a.motor).exists():
            mot = carga(a.motor, 'motor')
            if mot.get('formato') != 'salida/1' or mot.get('motor') != a.motor_nombre:
                raise Formato(f'salida del motor con formato o motor inválido: {a.motor}')
        else:
            aviso = f'sin motor (no existe {a.motor})'
    cont = {k: [0, 0] for k in ('oe', 'me', 'om', 'prov')}
    modelo_hojas = {}  # (fuente, ruta) -> [coincide, es_rama_propia]
    difs = []
    fixtures = sorted(p for p in carpeta.glob('*.json') if not p.name.startswith('_'))
    for p in fixtures:
        fx = carga(p, 'fixture')
        if fx.get('formato') != 'fixture/1' or fx.get('motor') != a.motor_nombre or fx.get('caso') != p.stem:
            raise Formato(f'fixture inválida: {p.name}')
        caso, prov = fx['caso'], fx.get('provisional') is not None
        ids = [q['id'] for q in fx['entradas']['consultas']]
        if len(ids) != len(set(ids)):
            raise Formato(f'{caso}: ids repetidos')
        for qid in ids:
            fuente, ruta_base = raiz_de(qid)
            if fuente:
                try:
                    nodo(fuentes[fuente], ruta_base)
                except (KeyError, IndexError, TypeError):
                    raise Formato(f'{caso}: el id {qid} no existe en {fuente}') from None
        for qid, esp in fx['esperado'].items():
            fuente, ruta_base = raiz_de(qid)
            o = dict(hojas(ora['casos'].get(caso, {}).get(qid, FALTA), qid))
            m = dict(hojas(mot['casos'].get(caso, {}).get(qid, FALTA), qid)) if mot else None
            for r, v in hojas(esp, qid):
                ov = o.get(r, FALTA)
                mv = m.get(r, FALTA) if m is not None else None
                bien = igual(ov, v)
                if prov:
                    cont['prov'][0] += bien
                    cont['prov'][1] += 1
                else:
                    cont['oe'][0] += bien
                    cont['oe'][1] += 1
                    if m is not None:
                        cont['me'][0] += igual(mv, v)
                        cont['me'][1] += 1
                        cont['om'][0] += igual(ov, mv)
                        cont['om'][1] += 1
                linea = None
                if not bien or (m is not None and not (igual(mv, v) and igual(ov, mv))):
                    linea = f'DIF {caso} {r} esperado={json.dumps(v, ensure_ascii=False)} oraculo={json.dumps(ov, ensure_ascii=False)}'
                    linea += f' motor={json.dumps(mv, ensure_ascii=False)}' if m is not None else ''
                if fuente and not prov:
                    rr = (ruta_base + r[len(qid):]) if fuente == 'modelo' else ruta_base + r[len(qid):]
                    if r.endswith('.#'):
                        try:
                            largo = len(nodo(fuentes[fuente], rr[:-2]))
                        except (KeyError, IndexError, TypeError):
                            largo = None
                        if largo != v and linea is None:
                            linea = f'DIF {caso} {r} esperado={v} modelo={largo}'
                    else:
                        try:
                            hv = nodo(fuentes[fuente], rr)
                        except (KeyError, IndexError, TypeError):
                            hv = None
                        mvv = valor_modelo(hv, rr)
                        if mvv is not None:
                            ok = coincide_modelo(mvv, v) and coincide_modelo(mvv, ov) and (m is None or coincide_modelo(mvv, mv))
                            clave = (fuente, rr)
                            propia = fuente == 'modelo' and pasos(rr)[0] in RAMAS.get(a.motor_nombre, ())
                            prev = modelo_hojas.get(clave, [True, propia])
                            modelo_hojas[clave] = [prev[0] and ok, propia]
                            if not ok and linea is None:
                                linea = f'DIF {caso} {r} esperado={json.dumps(v)} oraculo={json.dumps(ov)} modelo={hv}'
                if linea:
                    difs.append(linea)
    propias_total = [r for rama in RAMAS.get(a.motor_nombre, ()) for r, v in hojas(modelo[rama], rama) if valor_modelo(v, r) is not None and not r.endswith('.#')]
    sin_dueno = [r for r in propias_total if ('modelo', r) not in modelo_hojas]
    propias_ok = sum(1 for (f, r), (ok, prop) in modelo_hojas.items() if prop and ok)
    cruces = [(ok) for (f, r), (ok, prop) in modelo_hojas.items() if f == 'modelo' and not prop]
    datas = [(ok) for (f, r), (ok, prop) in modelo_hojas.items() if f == 'data']
    print(f'compara {a.motor_nombre} · {len(fixtures)} fixtures')
    print(f"oráculo contra esperado: {cont['oe'][0]} de {cont['oe'][1]} hojas coinciden")
    if mot is None:
        print(f'motor contra esperado: {aviso or "sin motor (--sin-motor)"}')
        print(f'oráculo contra motor: {aviso or "sin motor (--sin-motor)"}')
    else:
        print(f"motor contra esperado: {cont['me'][0]} de {cont['me'][1]} hojas coinciden")
        print(f"oráculo contra motor: {cont['om'][0]} de {cont['om'][1]} hojas coinciden")
    print(f'model.json (ramas de {a.motor_nombre}): {propias_ok} de {len(propias_total)} hojas coinciden')
    print(f'model.json (cruces con otras ramas): {sum(cruces)} de {len(cruces)} hojas coinciden')
    print(f'data.json: {sum(datas)} de {len(datas)} hojas coinciden')
    if cont['prov'][1]:
        print(f"provisionales: {cont['prov'][0]} de {cont['prov'][1]} hojas coinciden")
    print(f'diferencias: {len(difs)}')
    for d in difs:
        print(d)
    for r in sin_dueno:
        print(f'SIN-DUEÑO {r}')
    return 0 if not difs and not sin_dueno else 1


def main(argv=None):
    ap = argparse.ArgumentParser(description='Compara esperado, oráculo y motor de las fixtures de un motor')
    ap.add_argument('motor_nombre', metavar='motor', choices=sorted(RAMAS))
    ap.add_argument('--fixtures')
    ap.add_argument('--oraculo')
    ap.add_argument('--motor')
    ap.add_argument('--modelo', default=str(RAIZ / 'oraculo' / 'model.json'))
    ap.add_argument('--data', default=str(RAIZ / 'datos' / 'semilla' / 'data.json'))
    ap.add_argument('--sin-motor', action='store_true')
    a = ap.parse_args(argv)
    a.fixtures = a.fixtures or str(RAIZ / 'oraculo' / 'fixtures' / a.motor_nombre)
    a.oraculo = a.oraculo or str(RAIZ / 'oraculo' / 'salidas' / f'{a.motor_nombre}.json')
    a.motor = a.motor or str(RAIZ / 'datos' / 'tmp' / f'motor-{a.motor_nombre}.json')
    try:
        return compara(a)
    except Formato as e:
        print(f'compara {a.motor_nombre}: formato inválido o archivo faltante: {e}', file=sys.stderr)
        return 2


if __name__ == '__main__':
    sys.exit(main())
