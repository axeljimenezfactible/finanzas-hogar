#!/usr/bin/env python3
"""Ejecutor del oráculo S3 (FORMATO.md §4): corre las consultas de oraculo/fixtures/s3 con cred_lib y
cred_escenarios y escribe oraculo/salidas/s3.json (salida/1, recortada con «esperado» como plantilla).

Uso: python3 -I oraculo/cred_run_fixtures.py [--fixtures oraculo/fixtures/s3] [--out oraculo/salidas/s3.json]
     [--caso <nombre>]. Código 0 si escribió la salida; 1 si no pudo leer o escribir; 2 por uso o formato.
"""
import argparse
import copy
import json
import math
import os
import sys

BASE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, BASE)
import cred_lib as L  # noqa: E402


def merge_patch(target, patch):
    """JSON Merge Patch (RFC 7396)."""
    if not isinstance(patch, dict):
        return copy.deepcopy(patch)
    out = copy.deepcopy(target) if isinstance(target, dict) else {}
    for k, v in patch.items():
        if v is None:
            out.pop(k, None)
        else:
            out[k] = merge_patch(out.get(k), v)
    return out


def load_fixture(path):
    with open(path, encoding='utf-8') as fh:
        fx = json.load(fh)
    caso = os.path.basename(path)[:-5]
    if fx.get('formato') != 'fixture/1' or fx.get('motor') != 's3' or fx.get('caso') != caso:
        raise ValueError(f'{path}: formato, motor o caso inválido')
    datos = fx['entradas']['datos']
    if isinstance(datos, dict) and '$ref' in datos:
        with open(os.path.join(os.path.dirname(path), datos['$ref']), encoding='utf-8') as fh:
            base = json.load(fh)
        datos = merge_patch(base, datos['parche']) if 'parche' in datos else base
    fx['_datos'] = datos
    ids = [q['id'] for q in fx['entradas']['consultas']]
    if len(ids) != len(set(ids)):
        raise ValueError(f'{path}: id repetido')
    return fx


def trim(calc, exp):
    if isinstance(exp, dict):
        if not isinstance(calc, dict):
            return '__tipo__'
        return {k: (trim(calc[k], v) if k in calc else '__falta__') for k, v in exp.items()}
    if isinstance(exp, list):
        if not isinstance(calc, list):
            return '__tipo__'
        return [trim(c, exp[i]) if i < len(exp) else '__sobra__' for i, c in enumerate(calc)]
    if isinstance(calc, (dict, list)):
        return '__tipo__'
    return calc


def run_case(fx):
    res = {}
    esp = fx.get('esperado') or {}
    for q in fx['entradas']['consultas']:
        out, code = L.evaluate(fx['_datos'], q['funcion'], q.get('argumentos') or {})
        if code:
            res[q['id']] = {'error': code}
        elif q['id'] in esp:
            res[q['id']] = trim(out, esp[q['id']])
        else:
            res[q['id']] = out
    return res


# ------------------------------------------------------------------ segunda vía en float (como verify3.py)
def _r2(x):
    """HALF_UP a 2 decimales en flotante (segunda vía; nunca entra en una salida)."""
    return math.floor(x * 100 + 0.5 + 1e-9) / 100


def segunda_via(fx_dir, casos):
    """Recalcula en float, desde la semilla y la receta de §10, las series y los totales del canon y los compara a
    2 decimales con la salida exacta del oráculo. Devuelve (coinciden, total, diferencias)."""
    with open(os.path.join(fx_dir, '_datos', 'kueski-semilla.json'), encoding='utf-8') as fh:
        sem = json.load(fh)
    cr = sem['credits'][0]
    lim = cr['credit']['limit_cents'] / 100
    r4 = [t for t in cr['credit']['term_rates'] if t['term_count'] == 4][0]['rate']
    tasa = r4['num'] / r4['den']
    pares, desp = [], []
    for k, pst in enumerate(cr['purchases']):
        i = pst['input']
        P, n = i['principal_cents'] / 100, i['term_count']
        I = i['quoted_interest_cents'] / 100 if i['quoted_interest_cents'] is not None else _r2(P * tasa)
        T = _r2(P + I)
        base = math.floor(T / n * 100 + 1e-9) / 100
        montos = [base] * (n - 1) + [_r2(T - base * (n - 1))]
        cb = _r2(P / n)
        cap = [cb] * (n - 1) + [_r2(P - cb * (n - 1))]
        kc = [s['principal_cents'] for s in i['shares'] if s['external']][0] / 100
        ki = _r2(I * kc / P)
        kt = _r2(kc + ki)
        kb = _r2(kt / n)
        kp = [kb] * (n - 1) + [_r2(kt - kb * (n - 1))]
        ax = [_r2(m - x) for m, x in zip(montos, kp)]
        o = casos.get('s3-kueski-despensa-%d' % (k + 1), {}).get('despensas[%d]' % k, {})
        desp.append((pst['input']['date'], P, cap, i['date']))
        f = [(I, o.get('intereses')), (T, o.get('total')), (ki, o.get('karen', {}).get('intereses')),
             (kt, o.get('karen', {}).get('total')), (_r2(I - ki), o.get('axel', {}).get('intereses')),
             (_r2(P - kc + I - ki), o.get('axel', {}).get('total'))]
        for j, row in enumerate(o.get('pagos', [])):
            f += [(montos[j], row.get('monto')), (cap[j], row.get('capital_repuesto')),
                  (kp[j], row.get('karen')), (ax[j], row.get('axel'))]
        pares += f
    # series A y B: eventos de compra y de pago del 10-oct-2026 al 10-abr-2027 (lectura B con tope)
    import cred_lib as LL
    hoy, hasta = sem['today'], '2027-04-10'
    ev = []
    for fecha, P, cap, _ in desp:
        if hoy < fecha <= hasta:
            ev.append((fecha, 1, -P))
        q = LL.paydays_after(LL.pd(fecha), cr['credit']['first_payment_min_days'], len(cap))
        ev += [(LL.iso(d), 3, c) for d, c in zip(q, cap) if hoy < LL.iso(d) <= hasta]
    ev.sort()
    A, B = lim, _r2(lim - sum(P for fecha, P, _, _ in desp if fecha <= hoy))
    sa, sb, ma, mb = [], [], A, B
    for _, _, d in ev:
        A, B = _r2(A + d), _r2(min(lim, B + d))
        sa.append(A)
        sb.append(B)
        ma, mb = min(ma, A), min(mb, B)
    oc = casos.get('s3-series-a-b', {}).get('credito', {})
    cs = oc.get('comparacion_series', {})
    pares += list(zip(sa, cs.get('calculado_A', []))) + list(zip(sb, cs.get('calculado_B', [])))
    pares += [(ma, oc.get('minimo_A')), (mb, oc.get('minimo_B'))]
    if len(sa) != len(cs.get('calculado_A', [])):
        pares.append((len(sa), None))
    # sensibilidad e interés por monto
    for fila in casos.get('s3-sensibilidad', {}).get('sensibilidad', []):
        N, kc = fila['N'] / 100, 1191.50
        t = N * (1 + tasa)
        kt = _r2(kc + _r2(kc * tasa))
        pares += [(_r2(N * tasa), fila['intereses']), (_r2(N + _r2(N * tasa)), fila['total']),
                  (_r2(t / 4), fila['pago_quincenal']), (_r2(N + _r2(N * tasa) - kt), fila['axel_total']),
                  (_r2(t / 2 + 99.50), fila['costo_mensual_total']), (_r2((t - kt) / 2), fila['axel_mensual']),
                  (_r2(kt / 2 + 99.50), fila['karen_mensual'])]
    for fila in casos.get('s3-interes-por-monto', {}).get('interes_por_monto', []):
        m = fila['monto'] / 100
        pares += [(_r2(m * tasa), fila['intereses']), (_r2(m + _r2(m * tasa)), fila['total']),
                  (_r2(m * (1 + tasa) / 4), fila['pago_quincenal'])]
    difs = [(a, b) for a, b in pares if not isinstance(b, int) or abs(a - b / 100) >= 0.005]
    return len(pares) - len(difs), len(pares), difs


def main(argv=None):
    ap = argparse.ArgumentParser()
    ap.add_argument('--fixtures', default=os.path.join(BASE, 'fixtures', 's3'))
    ap.add_argument('--out', default=os.path.join(BASE, 'salidas', 's3.json'))
    ap.add_argument('--caso')
    a = ap.parse_args(argv)
    try:
        names = sorted(n for n in os.listdir(a.fixtures) if n.endswith('.json') and not n.startswith('_'))
    except OSError as e:
        print(f'no pude leer {a.fixtures}: {e}', file=sys.stderr)
        return 1
    if a.caso:
        names = [n for n in names if n[:-5] == a.caso]
        if not names:
            print(f'caso {a.caso} no existe', file=sys.stderr)
            return 2
    casos = {}
    for n in names:
        try:
            fx = load_fixture(os.path.join(a.fixtures, n))
        except (ValueError, KeyError, json.JSONDecodeError) as e:
            print(f'formato inválido: {e}', file=sys.stderr)
            return 2
        casos[fx['caso']] = run_case(fx)
    try:
        os.makedirs(os.path.dirname(os.path.abspath(a.out)), exist_ok=True)
        with open(a.out, 'w', encoding='utf-8') as fh:
            json.dump({'formato': 'salida/1', 'motor': 's3', 'casos': casos}, fh, ensure_ascii=False, indent=2)
            fh.write('\n')
    except OSError as e:
        print(f'no pude escribir {a.out}: {e}', file=sys.stderr)
        return 1
    print(f'oráculo s3 · {len(casos)} fixtures · {sum(len(v) for v in casos.values())} consultas → {a.out}')
    if not a.caso and os.path.exists(os.path.join(a.fixtures, '_datos', 'kueski-semilla.json')):
        ok, tot, difs = segunda_via(a.fixtures, casos)
        print(f'segunda vía (float, 2 decimales): {ok} de {tot} cifras coinciden')
        for x, y in difs:
            print(f'  DIF-FLOAT float={x} oraculo={y}')
    return 0


if __name__ == '__main__':
    sys.exit(main())
