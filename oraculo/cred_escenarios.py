#!/usr/bin/env python3
"""Comparador de escenarios del oráculo S3 por fuerza bruta propia (S3-interfaz.md §8; K16, K17).

Genera todos los candidatos (crédito × plazo × fecha ± días × tienda × promoción), aplica las 3
restricciones duras día por día, mide los 4 objetivos (desviación con math.isqrt), normaliza mín-máx
entre los válidos y elige 4 óptimos y la recomendación. Sin argumentos corre los casos del comparador
de oraculo/fixtures/s3 e imprime el resumen (aceptados con violación, óptimos distintos).

Tiempo máximo (§8.7; REVISION-S3-r2 M2, opción (a)): el oráculo no mide reloj. timeout_seconds = 0
responde timeout (determinista, fixture s3-comparador-tiempo) y un valor positivo calcula completo. El
tiempo real contra los 5 s (K43, configurable por K44) lo mide solo el motor en TestRestriccionesTiempo;
así la salida del oráculo no depende de qué tan cargado esté el equipo. Los ms del resumen son informativos.
"""
import copy
import datetime as dt
import json
import math
import os
import sys
import time
from fractions import Fraction as F

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import cred_lib as L  # noqa: E402

CONS = ('min_available_buffer_cents', 'max_live_purchases', 'zero_debt_by_date', 'date_window_days')
VIOL = (('min_available_buffer_cents', 'min_available_buffer'), ('max_live_purchases', 'max_live_purchases'),
        ('zero_debt_by_date', 'zero_debt_by_date'))
OBJ = (('lowest_interest', 'interest_cents', -1), ('even_low_payment', 'std_dev_milli_cents', -1),
       ('max_available', 'min_available_cents', 1), ('most_value_per_peso', 'value_per_peso_milli', 1))
REF = '__escenario__'
DAY = dt.timedelta(days=1)


def std_milli(xs):
    """§8.4: HALF_UP(1000 x desviación estándar poblacional) exacto con math.isqrt."""
    n = len(xs)
    N = 10 ** 6 * (n * sum(x * x for x in xs) - sum(xs) ** 2)
    return (math.isqrt(4 * N) // n + 1) // 2


def normaliza(vals, sign):
    """§8.5: mín-máx entre los válidos; empate o todos nulos = 1; sign < 0 menor es mejor.
    SUPUESTO (no ocurre en los casos): un nulo con otros no nulos vale 0."""
    nn = [v for v in vals if v is not None]
    if not nn or max(nn) == min(nn):
        return [F(1)] * len(vals)
    lo, hi = min(nn), max(nn)
    return [F(0) if v is None else (F(hi - v, hi - lo) if sign < 0 else F(v - lo, hi - lo)) for v in vals]


def key(s):
    return (s['credit_id'], s['term_count'], s['purchase_date'],
            (0,) if s['store_id'] is None else (1, s['store_id']),
            (0,) if s['promotion_id'] is None else (1, s['promotion_id']))


def _options(c, d, amount, sid, mer):
    kind = c['kind']
    promos = [p for p in sorted(c.get('promotions') or [], key=lambda p: p['id'])
              if p in L.applicable_promotions(c, d, amount, p['term_count'], sid, mer)]
    if kind in ('revolving_line', 'store_credit'):
        rates = {t['term_count'] for t in c.get('term_rates') or []}
        out = []
        for t in sorted(rates | {p['term_count'] for p in promos}):
            if t in rates:
                out.append((t, None))
            out += [(t, p['id']) for p in promos if p['term_count'] == t]
        return out
    if kind == 'bank_card':
        return [(1, None)] + [(p['term_count'], p['id']) for p in promos]
    return []


def _evaluate(datos, st, d, sid, mer, inp, q, used, req, today):
    c = st['credit']
    st2 = copy.deepcopy(st)
    st2['purchases'] = list(st2.get('purchases') or []) + [{'ref': REF, 'status': 'active', 'input': inp}]
    last = L.pd(q['payments'][-1]['due_date'])
    tl = L.timeline(st2, datos.get('cash'), today, today, last)
    k = 'available_a_cents' if tl['reading'] == 'A' else 'available_b_cents'
    evs = [(L.pd(e['date']), e[k]) for e in tl['events']]
    lives = [(L.pd(p['input']['date']), L.pd(qq['payments'][-1]['due_date'])) for p, qq in L.live_quotes(st2)]
    min_av, live_max, day = None, 0, d
    while day <= last:
        v = [x for e, x in evs if e <= day][-1]
        min_av = v if min_av is None else min(min_av, v)
        live_max = max(live_max, sum(1 for a, b in lives if a <= day <= b))
        day += DAY
    zero = max(b for _, b in lives)
    others = [s for s in datos.get('credits') or [] if s['credit']['id'] != c['id']] + [st2]
    f0, f1 = L.fortnight(d), L.fortnight(last)
    cal = L.unified_calendar(others, f0[0], f1[1])
    fns, cur = [], f0
    while cur[2] <= f1[2]:
        fns.append(cur[2])
        cur = L.fortnight(cur[1] + DAY)
    tot = {f: 0 for f in fns}
    per = {f: {} for f in fns}
    prov = {L.P_OBJ} | set(q['provisional'])
    for e in cal['entries']:
        f = L.fortnight(L.pd(e['date']))[2]
        if f in tot:
            tot[f] += e['amount_cents']
            prov |= set(e['provisional'])
            for sh in e['shares']:
                per[f][sh['person']] = per[f].get(sh['person'], 0) + sh['amount_cents']
    xs = [tot[f] for f in fns]
    std = std_milli(xs)
    over = []
    for cap in sorted(req.get('caps') or [], key=lambda c: c['person']):
        peak = max(per[f].get(cap['person'], 0) for f in fns)
        if peak > cap['cap_cents']:
            over.append({'person': cap['person'], 'peak_cents': peak, 'cap_cents': cap['cap_cents']})
    vpp = None
    for b in req.get('baskets') or []:
        if (sid is not None and b['store_id'] == sid) or \
                (sid is None and b['merchant_name'].strip().lower() == (mer or '').strip().lower()):
            A = inp['principal_cents']
            vpp = L.hu(F(A * b['base_units_milli'] * 100, b['cost_cents'] * (A + q['interest_cents'])))
            break
    s = {'credit_id': c['id'], 'term_count': inp['term_count'], 'purchase_date': L.iso(d), 'store_id': sid,
         'promotion_id': inp['promotion_id'], 'interest_cents': q['interest_cents'],
         'max_fortnight_payment_cents': max(xs), 'fortnight_payment_spread_cents': max(xs) - min(xs),
         'std_dev_milli_cents': std, 'min_available_cents': min_av, 'max_live_purchases': live_max,
         'zero_debt_date': L.iso(zero), 'value_per_peso_milli': vpp, 'score_bp': None, 'over_cap': over,
         'provisional': sorted(prov)}
    viol = None
    if min_av < used['min_available_buffer_cents']:
        viol = 'min_available_buffer'
    elif live_max > used['max_live_purchases']:
        viol = 'max_live_purchases'
    elif zero > L.pd(used['zero_debt_by_date'], False):
        viol = 'zero_debt_by_date'
    return s, viol


def compare(datos, timeout_seconds, req):
    if not isinstance(timeout_seconds, int) or isinstance(timeout_seconds, bool) or timeout_seconds < 0:
        L.falla('invalid', 'timeout_seconds')
    if timeout_seconds == 0:   # §8.7: 0 = vence de inmediato; con un valor positivo no se mide reloj (M2)
        L.falla('timeout')
    L.check_input(datos)
    today = L.pd(datos.get('today'), False)
    g, o = req.get('global') or {}, req.get('override') or {}
    used = {k: (o.get(k) if o.get(k) is not None else g.get(k)) for k in CONS}
    missing = [k for k in CONS if used[k] is None]
    if missing:
        return {'status': 'missing_constraints', 'missing_constraints': missing, 'constraints_used': used,
                'evaluated_count': 0, 'valid': [], 'optima': None, 'recommendation': None, 'rejected': [],
                'excluded': [], 'provisional': []}
    w = req['weights']
    ws = [w['lowest_interest'], w['even_low_payment'], w['max_available'], w['most_value_per_peso']]
    if any(x < 0 for x in ws) or sum(ws) == 0:
        L.falla('invalid', 'pesos')
    amount, desired, win = req['amount_cents'], L.pd(req['desired_date'], False), used['date_window_days']
    dates = [desired + k * DAY for k in range(-win, win + 1) if desired + k * DAY >= today]
    if req.get('store_id') is not None or (req.get('merchant_name') or ''):
        stores = [(req.get('store_id'), req.get('merchant_name') or '')]
    else:
        stores = [(b['store_id'], b['merchant_name']) for b in req.get('baskets') or []] or [(None, '')]
    valid, rejected, excluded = [], [], []
    for st in sorted(datos.get('credits') or [], key=lambda s: s['credit']['id']):
        c = st['credit']
        if c.get('limit_cents') is None:
            excluded.append({'credit_id': c['id'], 'reason': 'no_limit'})
            continue
        for d in dates:
            for sid, mer in stores:
                for t, pid in _options(c, d, amount, sid, mer):
                    inp = {'date': L.iso(d), 'principal_cents': amount, 'term_count': t,
                           'quoted_interest_cents': None, 'promotion_id': pid, 'store_id': sid, 'merchant': mer,
                           'first_payment_date': None, 'shares': req['shares']}
                    try:
                        q = L.quote_purchase(c, inp)
                    except L.CredError:
                        continue
                    s, viol = _evaluate(datos, st, d, sid, mer, inp, q, used, req, today)
                    if viol is None:
                        valid.append(s)
                    else:
                        rejected.append({'scenario': s, 'violated_constraint': viol})
    valid.sort(key=key)
    rejected.sort(key=lambda r: key(r['scenario']))
    prov = {L.P_OBJ}
    for s in valid + [r['scenario'] for r in rejected]:
        prov |= set(s['provisional'])
    out = {'status': 'ok' if valid else 'no_valid_scenario', 'missing_constraints': [], 'constraints_used': used,
           'evaluated_count': len(valid) + len(rejected), 'valid': valid, 'optima': None,
           'recommendation': None, 'rejected': rejected, 'excluded': excluded, 'provisional': sorted(prov)}
    if not valid:
        return out
    cols = [normaliza([s[f] for s in valid], sign) for _, f, sign in OBJ]
    norm = {i: [c[i] for c in cols] for i in range(len(valid))}
    score = {i: sum(F(wi) * x for wi, x in zip(ws, norm[i])) / sum(ws) for i in range(len(valid))}
    for i, s in enumerate(valid):
        s['score_bp'] = L.hu(10000 * score[i])
    idx = range(len(valid))
    optima = {}
    for name, f, sign in OBJ:
        vals = [valid[i][f] for i in idx if valid[i][f] is not None]
        best = (min(vals) if sign < 0 else max(vals)) if vals else None
        cand = [i for i in idx if valid[i][f] == best] if vals else list(idx)
        if name == 'even_low_payment':
            cand.sort(key=lambda i: (len(valid[i]['over_cap']), -score[i], key(valid[i])))
        else:
            cand.sort(key=lambda i: (-score[i], key(valid[i])))
        optima[name] = valid[cand[0]]
    rec = sorted(idx, key=lambda i: (-score[i], key(valid[i])))[0]
    out['optima'] = optima
    out['recommendation'] = valid[rec]
    return out


CASOS = ('s3-comparador-caso-1', 's3-comparador-caso-2', 's3-comparador-caso-3',
         's3-comparador-sin-restricciones', 's3-comparador-tiempo', 's3-comparador-carga')


def autoprueba():
    """Ejemplos de la interfaz (§8.4 y §8.5) mientras falten las fixtures del optimizador (§10.0)."""
    ok = std_milli([0, 100, 200]) == 81650
    ok &= normaliza([100, 200, 400], -1) == [F(1), F(2, 3), F(0)]
    ok &= normaliza([5, 5, 5], -1) == [F(1)] * 3 and normaliza([None, None], 1) == [F(1)] * 2
    w = [1, 1, 1, 1]
    ok &= L.hu(10000 * sum(F(x) * y for x, y in zip(w, [F(1), F(2, 3), F(0), F(1)])) / sum(w)) == 6667
    print('autoprueba §8.4 (0/100/200 -> 81650) y §8.5 (100/200/400 -> 1, 2/3, 0; empate -> 1): '
          + ('bien' if ok else 'MAL'))
    return ok


def _resumen(fx_dir):
    import cred_run_fixtures as R
    total_bad, faltan = 0, 0
    for caso in CASOS:
        ruta = os.path.join(fx_dir, caso + '.json')
        if not os.path.exists(ruta):
            print(f'falta oraculo/fixtures/s3/{caso}.json (autor: optimizador S3)')
            faltan += 1
            continue
        fx = R.load_fixture(ruta)
        for qy in fx['entradas']['consultas']:
            t0 = time.monotonic()
            out, code = L.evaluate(fx['_datos'], qy['funcion'], qy['argumentos'])
            ms = round((time.monotonic() - t0) * 1000)
            if code:
                print(f'{caso}: error {code} ({ms} ms)')
                continue
            used = out['constraints_used']
            bad = sum(1 for s in out['valid'] if s['min_available_cents'] < used['min_available_buffer_cents']
                      or s['max_live_purchases'] > used['max_live_purchases']
                      or s['zero_debt_date'] > used['zero_debt_by_date'])
            total_bad += bad
            op = out['optima'] or {}
            ks = [key(op[k]) for k in op]
            print(f"{caso}: {out['status']}, evaluados {out['evaluated_count']}, válidos {len(out['valid'])}, "
                  f"descartados {len(out['rejected'])}, aceptados con violación {bad}, "
                  f"óptimos distintos {len(set(ks))} de {len(ks)}, {ms} ms")
    print(f'escenarios aceptados con violación = {total_bad}')
    return total_bad, faltan


if __name__ == '__main__':
    base = os.path.dirname(os.path.abspath(__file__))
    bien = autoprueba()
    malos, faltan = _resumen(os.path.join(base, 'fixtures', 's3'))
    sys.exit(0 if bien and not malos and not faltan else 1)
