#!/usr/bin/env python3
"""Oráculo del motor de créditos S3 (docs/contratos/motores/S3-interfaz.md §3–§7 y §9).

Método distinto del motor Go: fracciones exactas (fractions.Fraction), calendarios día por día.
Solo biblioteca estándar. No importa model.py (al importarse reescribe model.json) ni lee nada del
motor Go (S3-16). Dinero en centavos enteros; ROUND_HALF_UP alejándose de cero en el último paso.
"""
import calendar
import datetime as dt
import math
from fractions import Fraction as F

P_PRESTAMO = 'S3-convencion-prestamo'
P_TARJETA = 'S3-formula-tarjeta'
P_OBJ = 'S3-objetivos-medida'
P_CARGO = 'S3-origen-cargo-normal'
P_SERIE = 'S3-serie-sin-intereses'
P_IVA = 'U4-iva-credito'
KINDS = ('revolving_line', 'bank_card', 'term_loan', 'store_credit')
FREQS = ('weekly', 'fortnightly', 'monthly')


class CredError(Exception):
    def __init__(self, code, detalle=''):
        super().__init__(f'{code}: {detalle}')
        self.code = code


def falla(code, detalle=''):
    raise CredError(code, detalle)


# ---------------------------------------------------------------- dinero y tasas (§3)
def hu(x):
    """ROUND_HALF_UP a entero (centavo) de una fracción."""
    x = F(x)
    s = -1 if x < 0 else 1
    return s * math.floor(abs(x) + F(1, 2))


def split_trunc(total, n):
    base = total // n
    return [base] * (n - 1) + [total - base * (n - 1)]


def split_halfup(total, n):
    base = hu(F(total, n))
    return [base] * (n - 1) + [total - base * (n - 1)]


def frac(r):
    if not isinstance(r, dict) or not isinstance(r.get('num'), int) or not isinstance(r.get('den'), int):
        falla('invalid', 'tasa')
    if r['den'] <= 0 or r['num'] < 0:
        falla('invalid', 'tasa')
    return F(r['num'], r['den'])


def apply_rate(c, r):
    return hu(c * frac(r))


def display(r):
    """K12: HALF_UP a los decimales capturados; sin dato de captura (-1) no hay texto (hallazgo N2 de la r3)."""
    d = r.get('decimales', -1)
    if d is None or d < 0:
        return ''
    q = hu(frac(r) * 100 * 10 ** d)
    s = str(q).rjust(d + 1, '0')
    return (s[:-d] + '.' + s[-d:] if d > 0 else s) + ' %'


def rat_json(x):
    x = F(x)
    return {'num': x.numerator, 'den': x.denominator}


# ---------------------------------------------------------------- fechas (§3, RN-06)
def pd(s, opcional=True):
    if s is None:
        if opcional:
            return None
        falla('invalid', 'fecha nula')
    if not isinstance(s, str) or len(s) != 10 or s[4] != '-' or s[7] != '-':
        falla('invalid', f'fecha {s!r}')
    try:
        return dt.date(int(s[:4]), int(s[5:7]), int(s[8:]))
    except ValueError:
        falla('invalid', f'fecha {s!r}')


def iso(d):
    return d.isoformat()


def last_day(d):
    return calendar.monthrange(d.year, d.month)[1]


def is_payday(d):
    return d.day == 15 or d.day == min(30, last_day(d))


def fortnight(d):
    """(inicio, fin, día de pago) de la quincena de d (K8)."""
    if d.day <= 15:
        return d.replace(day=1), d.replace(day=15), d.replace(day=15)
    return d.replace(day=16), d.replace(day=last_day(d)), d.replace(day=min(30, last_day(d)))


def paydays_after(buy, min_days, n):
    out, d = [], buy
    while len(out) < n:
        d += dt.timedelta(days=1)
        if is_payday(d) and (d - buy).days >= min_days:
            out.append(d)
    return out


def dates_from_first(first, n, freq):
    if freq == 'weekly':
        return [first + dt.timedelta(days=7 * k) for k in range(n)]
    if freq == 'fortnightly':
        out, d = [first], first
        while len(out) < n:
            d += dt.timedelta(days=1)
            if is_payday(d):
                out.append(d)
        return out
    if freq == 'monthly':
        out = []
        for k in range(n):
            y = first.year + (first.month - 1 + k) // 12
            m = (first.month - 1 + k) % 12 + 1
            out.append(dt.date(y, m, min(first.day, calendar.monthrange(y, m)[1])))
        return out
    falla('invalid', f'frecuencia {freq!r}')


def card_dues(credit, buy, n):
    """PROVISIONAL S3-formula-tarjeta (pregunta 4): corte = primera fecha >= compra con día de corte
    (recortado); pago = primera fecha > corte con día de pago (recortado); cortes sucesivos."""
    sd, dd = credit.get('statement_day'), credit.get('payment_due_day')
    dues, d, first = [], buy, True
    while len(dues) < n:
        if not first:
            d += dt.timedelta(days=1)
        first = False
        while d.day != min(sd, last_day(d)):
            d += dt.timedelta(days=1)
        e = d + dt.timedelta(days=1)
        while e.day != min(dd, last_day(e)):
            e += dt.timedelta(days=1)
        dues.append(e)
    return dues


# ---------------------------------------------------------------- personas y partes (§4)
def order_shares(shares, principal):
    if not isinstance(shares, list) or not shares:
        falla('invalid', 'sin partes')
    hold = [s for s in shares if not s.get('external')]
    if len(hold) != 1:
        falla('invalid', 'titular')
    for s in shares:
        if not isinstance(s.get('principal_cents'), int) or s['principal_cents'] < 0:
            falla('invalid', 'parte')
    if sum(s['principal_cents'] for s in shares) != principal:
        falla('invalid', 'suma de partes')
    ext = sorted([s for s in shares if s.get('external')], key=lambda s: (s.get('person_id') or 0, s['person']))
    return hold + ext


def build_shares(shares, principal, interest, amounts):
    n = len(amounts)
    ext = shares[1:]
    ext_int = [hu(F(interest * s['principal_cents'], principal)) for s in ext]
    ext_parts = [split_halfup(s['principal_cents'] + i, n) for s, i in zip(ext, ext_int)]
    hold_parts = [amounts[k] - sum(p[k] for p in ext_parts) for k in range(n)]
    if any(x < 0 for x in hold_parts):
        falla('invalid', 'parte del titular negativa')
    ints = [interest - sum(ext_int)] + ext_int
    parts = [hold_parts] + ext_parts
    out = []
    for s, i, p in zip(shares, ints, parts):
        j, pi = s.get('joint_cents'), s.get('paper_in_cents')
        comp = None
        if j is not None or pi is not None:
            j, pi = j or 0, pi or 0
            comp = {'joint_cents': j, 'own_cents': s['principal_cents'] - j,
                    'paper_in_cents': pi, 'free_cents': s['principal_cents'] - pi}
        out.append({'person': s['person'], 'external': bool(s.get('external')),
                    'principal_cents': s['principal_cents'], 'interest_cents': i,
                    'total_cents': s['principal_cents'] + i, 'first_payment_cents': p[0],
                    'composition': comp})
    pay_sh = [[{'person': s['person'], 'amount_cents': p[k]} for s, p in zip(shares, parts)] for k in range(n)]
    return out, pay_sh


# ---------------------------------------------------------------- promociones (§5.1, RN-10)
def applicable_promotions(credit, date, principal, term, store_id, merchant):
    out = []
    for p in sorted(credit.get('promotions') or [], key=lambda p: p['id']):
        if not p.get('active'):
            continue
        if p.get('credit_id') is not None and p['credit_id'] != credit['id']:
            continue
        if p.get('store_id') is not None and store_id is not None:
            if p['store_id'] != store_id:
                continue
        elif (p.get('merchant') or '').strip().lower() != (merchant or '').strip().lower():
            continue
        if not (pd(p['valid_from'], False) <= date <= pd(p['valid_to'], False)):
            continue
        if principal < p['min_purchase_cents'] or p['term_count'] != term:
            continue
        out.append(p)
    return out


def _term_rate(credit, n):
    for t in credit.get('term_rates') or []:
        if t['term_count'] == n:
            return t
    return None


# ---------------------------------------------------------------- cotización (§5.1–§5.5, §5.8)
def _purchase(P, I, total_exact, n, rate, basis, source, pkind, promo_id, dates, amounts, principals, shares, prov):
    if shares is None:
        sh_out, pay_sh = [], [[] for _ in amounts]
    else:
        sh_out, pay_sh = build_shares(shares, P, I, amounts)
    pays = [{'seq': k + 1, 'due_date': iso(dates[k]), 'amount_cents': amounts[k],
             'principal_cents': principals[k], 'interest_cents': amounts[k] - principals[k],
             'shares': pay_sh[k]} for k in range(n)]
    ref = hu(F(total_exact) / n) if source != 'amortization' else None
    return {'principal_cents': P, 'interest_cents': I, 'total_cents': P + I,
            'total_exact': rat_json(total_exact), 'reference_payment_cents': ref,
            'rate': rate, 'rate_display': display(rate) if rate is not None else None,
            'rate_basis': basis, 'source': source, 'purchase_kind': pkind, 'promotion_id': promo_id,
            'payments': pays, 'shares': sh_out, 'provisional': sorted(set(prov))}


def _interest_line(credit, inp, date, P, n, allow_k18):
    """(1) cotizado, (2) promoción, (3) tabla por plazo, (4) sin tasa. Devuelve (I, total_exact, rate, source, promo)."""
    q, pid = inp.get('quoted_interest_cents'), inp.get('promotion_id')
    if q is not None:
        if not isinstance(q, int) or q < 0:
            falla('invalid', 'interés cotizado')
        # §5.1 (r3-N2): los decimales son los de la compra (rate_decimals), nunca los de la tabla
        dec = inp.get('rate_decimals')
        return q, F(P + q), {'num': q, 'den': P, 'decimales': -1 if dec is None else dec}, 'quoted', None
    if pid is not None:
        ap = [p for p in applicable_promotions(credit, date, P, n, inp.get('store_id'), inp.get('merchant'))
              if p['id'] == pid]
        if not ap:
            falla('promotion_not_applicable')
        r = ap[0]['rate']
        return apply_rate(P, r), P * (1 + frac(r)), r, 'promotion', pid
    t = _term_rate(credit, n)
    if t is not None:
        return apply_rate(P, t['rate']), P * (1 + frac(t['rate'])), t['rate'], 'term_table', None
    if allow_k18:
        if any(credit.get(k) in (None, '') for k in ('annual_rate', 'rate_base', 'min_payment_pct_bp',
                                                      'min_payment_fixed_cents')):
            falla('missing_capture')
        falla('unsupported', 'tienda con interés K18 sin regla de abonos')
    falla('no_rate')


def amortize(P, annual, vat_bp, n, freq, first, shares=None):
    if not isinstance(P, int) or P <= 0 or not isinstance(n, int) or n < 1 or freq not in FREQS or first is None:
        falla('invalid', 'amortizar')
    r = frac(annual) / {'monthly': 12, 'fortnightly': 24, 'weekly': 52}[freq]
    if vat_bp is not None:
        r *= 1 + F(vat_bp, 10000)
    pay = hu(F(P, n)) if r == 0 else hu(P * r / (1 - (1 + r) ** (-n)))
    bal, amounts, principals = P, [], []
    for k in range(n):
        i = hu(bal * r)
        a = pay - i if k < n - 1 else bal
        if a < 0:
            falla('invalid', 'pago menor que el interés')
        amounts.append(a + i)
        principals.append(a)
        bal -= a
    I = sum(amounts) - P
    out = _purchase(P, I, F(P + I), n, annual, 'annual', 'amortization', 'loan_disbursement', None,
                    dates_from_first(first, n, freq), amounts, principals, shares,
                    [P_PRESTAMO] + ([P_IVA] if vat_bp is None else []))
    out['reference_payment_cents'] = pay
    return out


def quote_purchase(credit, inp):
    kind = credit.get('kind')
    if kind not in KINDS:
        falla('invalid', 'tipo')
    date = pd(inp.get('date'), False)
    P, n = inp.get('principal_cents'), inp.get('term_count')
    if not isinstance(P, int) or P <= 0 or not isinstance(n, int) or n < 1:
        falla('invalid', 'capital o plazo')
    rd = inp.get('rate_decimals')
    if rd is not None and (not isinstance(rd, int) or isinstance(rd, bool) or rd < 0):
        falla('invalid', 'rate_decimals')
    shares = order_shares(inp.get('shares'), P)
    fpd = pd(inp.get('first_payment_date'))
    if fpd is not None and fpd < date:
        falla('invalid', 'primer pago antes de la compra')
    freq = credit.get('frequency')
    q, pid = inp.get('quoted_interest_cents'), inp.get('promotion_id')
    if kind == 'revolving_line':
        if credit.get('limit_cents') is None or freq not in FREQS:
            falla('invalid', 'línea sin límite o frecuencia')
        if freq == 'fortnightly' and fpd is None and credit.get('first_payment_min_days') is None:
            falla('invalid', 'sin días mínimos')
        if fpd is not None:
            dates = dates_from_first(fpd, n, freq)
        elif freq == 'fortnightly':
            dates = paydays_after(date, credit['first_payment_min_days'], n)
        else:
            falla('unsupported', 'frecuencia sin regla')
        I, te, rate, src, promo = _interest_line(credit, inp, date, P, n, False)
        return _purchase(P, I, te, n, rate, 'term_total', src, 'installment_plan', promo, dates,
                         split_trunc(P + I, n), split_halfup(P, n), shares, [P_SERIE] if I == 0 and n >= 2 else [])
    if kind == 'bank_card':
        if credit.get('limit_cents') is None or credit.get('statement_day') is None \
                or credit.get('payment_due_day') is None or q is not None:
            falla('invalid', 'tarjeta')
        if pid is not None:
            ap = [p for p in applicable_promotions(credit, date, P, n, inp.get('store_id'), inp.get('merchant'))
                  if p['id'] == pid]
            if not ap:
                falla('promotion_not_applicable')
            if frac(ap[0]['rate']) != 0:
                falla('unsupported', 'meses con intereses')
            am = split_trunc(P, n)
            return _purchase(P, 0, F(P), n, ap[0]['rate'], 'term_total', 'promotion', 'installment_plan', pid,
                             card_dues(credit, date, n), am, am, shares, [P_TARJETA] + ([P_SERIE] if n >= 2 else []))
        if n != 1:
            falla('no_rate')
        return _purchase(P, 0, F(P), 1, {'num': 0, 'den': 1, 'decimales': 0}, '', 'quoted', 'single_charge',
                         None, card_dues(credit, date, 1), [P], [P], shares, [P_TARJETA, P_CARGO])
    if kind == 'term_loan':
        if q is not None or pid is not None or freq not in FREQS or fpd is None:
            falla('invalid', 'préstamo')
        if credit.get('annual_rate') is None:
            falla('missing_capture')
        return amortize(P, credit['annual_rate'], credit.get('vat_on_interest_bp'), n, freq, fpd, shares)
    # store_credit
    if freq not in FREQS:
        falla('invalid', 'frecuencia')
    if fpd is None and not (freq == 'fortnightly' and credit.get('first_payment_min_days') is not None):
        falla('invalid', 'tienda sin primera fecha')
    dates = dates_from_first(fpd, n, freq) if fpd is not None else \
        paydays_after(date, credit['first_payment_min_days'], n)
    I, te, rate, src, promo = _interest_line(credit, inp, date, P, n, True)
    am = split_trunc(P + I, n)
    return _purchase(P, I, te, n, rate, 'term_total', src, 'installment_plan', promo, dates,
                     am, am if I == 0 else split_halfup(P, n), shares, [P_SERIE] if I == 0 and n >= 2 else [])


# ---------------------------------------------------------------- estado (§5.7)
def check_input(datos):
    ids = [c['credit']['id'] for c in datos.get('credits') or []]
    if len(ids) != len(set(ids)):
        falla('invalid', 'crédito repetido')


def quote_state(st):
    credit = st['credit']
    refs = [p['ref'] for p in st.get('purchases') or []]
    if len(refs) != len(set(refs)):
        falla('invalid', 'ref repetida')
    for r in st.get('payments') or []:
        if r['purchase_ref'] not in refs:
            falla('invalid', 'abono sin compra')
    if st.get('statements') and credit.get('kind') != 'bank_card':
        falla('invalid', 'cortes fuera de tarjeta')
    ps = sorted(st.get('purchases') or [], key=lambda p: (p['input']['date'], p['ref']))
    return [(p, quote_purchase(credit, p['input'])) for p in ps]


def live_quotes(st):
    return [(p, q) for p, q in quote_state(st) if p['status'] != 'cancelled']


def find_credit(datos, cid):
    for c in datos.get('credits') or []:
        if c['credit']['id'] == cid:
            return c
    falla('invalid', f'crédito {cid}')


# ---------------------------------------------------------------- disponible A/B (§5.2, K15)
KORD = {'start': 0, 'purchase': 1, 'cash': 2, 'payment': 3}


def timeline(st, cash, today, frm, to, reading=None):
    credit = st['credit']
    limit = credit.get('limit_cents')
    if limit is None or to < frm:
        falla('invalid', 'línea sin límite o rango')
    rd = (credit.get('reading') or '') if reading is None else reading
    if rd not in ('', 'A', 'B'):
        falla('invalid', 'lectura')
    eff = rd or 'B'
    anchor = pd(credit.get('reading_a_anchor'))
    if eff == 'A' and anchor is None:
        falla('invalid', 'lectura A sin ancla')
    defaulted = anchor is None
    if defaulted:
        anchor = frm
    if frm < anchor:
        falla('invalid', 'desde antes del ancla')
    quotes = live_quotes(st)
    recs = st.get('payments') or []
    pev, rest = [], []
    for p, q in quotes:
        for sp in q['payments']:
            M, K = sp['amount_cents'], sp['principal_cents']
            rs = sorted([r for r in recs if r['purchase_ref'] == p['ref'] and r['seq'] == sp['seq']],
                        key=lambda r: (r['date'], r['payer'], r['amount_cents']))
            acc, done = 0, False
            for j, r in enumerate(rs):
                acc += r['amount_cents']
                d = pd(r['date'], False)
                if not done and acc >= M:
                    done = True
                    pev.append((d, p['ref'], sp['seq'], j, r['amount_cents'], K, M - K))
                    rest.append((d, K))
                else:
                    pev.append((d, p['ref'], sp['seq'], j, r['amount_cents'], 0, 0))
            due = pd(sp['due_date'], False)
            if acc < M and due >= today:
                pev.append((due, p['ref'], sp['seq'], 10 ** 9, M - acc, K, M - K))
                rest.append((due, K))
    pdates = [(pd(p['input']['date'], False), p['input']['principal_cents']) for p, _ in quotes]
    B = limit - sum(c for d, c in pdates if d <= frm) + sum(k for d, k in rest if d <= frm)
    A = limit - sum(c for d, c in pdates if anchor < d <= frm) + sum(k for d, k in rest if anchor < d <= frm)
    raw = []
    for p, _ in quotes:
        d = pd(p['input']['date'], False)
        if frm < d <= to:
            raw.append(((d, 1, p['ref'], 0, 0), 'purchase', p['ref'], 0, p['input']['principal_cents'],
                        -p['input']['principal_cents'], 0, f"Compra {p['ref']}"))
    for j, c in enumerate(cash or []):
        d = pd(c['date'], False)
        if frm < d <= to:
            raw.append(((d, 2, '', 0, j), 'cash', '', 0, c['amount_cents'], 0, 0, c.get('description', '')))
    for d, ref, seq, sub, amt, K, I in pev:
        if frm < d <= to:
            raw.append(((d, 3, ref, seq, sub), 'payment', ref, seq, amt, K, I, f'Pago {ref} {seq}'))
    raw.sort(key=lambda e: e[0])
    ev = [{'date': iso(frm), 'kind': 'start', 'description': 'Disponible al inicio', 'purchase_ref': '',
           'seq': 0, 'amount_cents': 0, 'principal_restored_cents': 0, 'interest_not_restored_cents': 0,
           'available_a_cents': A, 'available_b_cents': B}]
    sa, sb = [], []
    for key, kind, ref, seq, amt, delta, I, desc in raw:
        A += delta
        B = min(limit, B + delta)
        ev.append({'date': iso(key[0]), 'kind': kind, 'description': desc, 'purchase_ref': ref, 'seq': seq,
                   'amount_cents': amt, 'principal_restored_cents': delta if kind == 'payment' else 0,
                   'interest_not_restored_cents': I, 'available_a_cents': A, 'available_b_cents': B})
        if kind in ('purchase', 'payment'):
            sa.append(A)
            sb.append(B)
    return {'events': ev, 'series_a_cents': sa, 'series_b_cents': sb,
            'min_a_cents': min(e['available_a_cents'] for e in ev),
            'min_b_cents': min(e['available_b_cents'] for e in ev),
            'reading': eff, 'pending_confirmation': not credit.get('reading_confirmed'),
            'a_anchor': iso(anchor), 'a_anchor_defaulted': defaulted}


# ---------------------------------------------------------------- calendario unificado (§6) y efectivo (§6.1)
def _status(paid, amount):
    return 'pending' if paid == 0 else ('partial' if paid < amount else 'paid')


def unified_calendar(states, frm, to):
    entries, keys = [], set()
    for st in states:
        credit = st['credit']
        cid, kind = credit['id'], credit['kind']
        recs = st.get('payments') or []
        sched = [(p['ref'], sp, q['provisional']) for p, q in live_quotes(st) for sp in q['payments']]
        sdates = [s['statement_date'] for s in st.get('statements') or []]
        sdues = [s['due_date'] for s in st.get('statements') or []]
        if len(sdates) != len(set(sdates)) or len(sdues) != len(set(sdues)):
            falla('invalid', 'cortes repetidos')
        stm = {}
        for s in st.get('statements') or []:
            due = pd(s['due_date'], False)
            if frm <= due <= to:
                stm[due] = s

        def paid(ref, seq):
            return sum(r['amount_cents'] for r in recs if r['purchase_ref'] == ref and r['seq'] == seq)
        for ref, sp, prov in sched:
            due = pd(sp['due_date'], False)
            if frm <= due <= to and due not in stm:
                entries.append({'date': iso(due), 'credit_id': cid, 'credit_kind': kind, 'source': 'scheduled',
                                'purchase_ref': ref, 'seq': sp['seq'], 'statement_date': None,
                                'amount_cents': sp['amount_cents'], 'principal_cents': sp['principal_cents'],
                                'interest_cents': sp['interest_cents'], 'shares': sp['shares'],
                                'status': _status(paid(ref, sp['seq']), sp['amount_cents']),
                                'provisional': list(prov)})
        for due, s in stm.items():
            subs = [(ref, sp) for ref, sp, _ in sched if pd(sp['due_date'], False) == due]
            bal = s['balance_cents']
            princ = sum(sp['principal_cents'] for _, sp in subs)
            ext = {}
            for _, sp in subs:
                for sh in sp['shares']:
                    if sh['person'] != credit['holder']:
                        ext[sh['person']] = ext.get(sh['person'], 0) + sh['amount_cents']
            hold = bal - sum(ext.values())
            if bal - princ < 0 or hold < 0:
                falla('invalid', 'corte menor que lo que incluye')
            shares = [{'person': credit['holder'], 'amount_cents': hold}] + \
                     [{'person': k, 'amount_cents': ext[k]} for k in sorted(ext)]
            entries.append({'date': iso(due), 'credit_id': cid, 'credit_kind': kind, 'source': 'statement',
                            'purchase_ref': '', 'seq': 0, 'statement_date': s['statement_date'],
                            'amount_cents': bal, 'principal_cents': princ, 'interest_cents': bal - princ,
                            'shares': shares,
                            'status': _status(sum(paid(r, sp['seq']) for r, sp in subs), bal),
                            'provisional': [P_TARJETA]})
    for e in entries:
        k = (e['credit_id'], e['source'], e['purchase_ref'], e['seq'], e['statement_date'])
        if k in keys:
            falla('invalid', 'entrada duplicada')
        keys.add(k)
    entries.sort(key=lambda e: (e['date'], e['credit_id'], 0 if e['source'] == 'scheduled' else 1,
                                e['purchase_ref'], e['seq']))
    byf = {}
    for e in entries:
        pay = iso(fortnight(pd(e['date']))[2])
        t = byf.setdefault(pay, {'total': 0, 'p': {}})
        t['total'] += e['amount_cents']
        for sh in e['shares']:
            t['p'][sh['person']] = t['p'].get(sh['person'], 0) + sh['amount_cents']
    return {'entries': entries,
            'by_fortnight': [{'pay_day': k, 'total_cents': v['total'],
                              'by_person': [{'person': p, 'amount_cents': v['p'][p]} for p in sorted(v['p'])]}
                             for k, v in sorted(byf.items())]}


def cash_flow(cal, cash, frm, to):
    ents = [e for e in cal['entries'] if frm <= pd(e['date']) <= to]
    cs = [c for c in cash or [] if frm <= pd(c['date'], False) <= to]
    per = {}
    for e in ents:
        for sh in e['shares']:
            per[sh['person']] = per.get(sh['person'], 0) + sh['amount_cents']
    for c in cs:
        per[c['person']] = per.get(c['person'], 0) + c['amount_cents']
    pay = sum(e['amount_cents'] for e in ents)
    ca = sum(c['amount_cents'] for c in cs)
    return {'from': iso(frm), 'to': iso(to), 'payments_cents': pay, 'cash_cents': ca, 'total_cents': pay + ca,
            'payment_count': len(ents),
            'by_person': [{'person': p, 'amount_cents': per[p]} for p in sorted(per)]}


# ---------------------------------------------------------------- deuda y compras vivas (§7, RN-21)
def live_debt(states, as_of):
    byp, byc, tot = [], [], 0
    for st in sorted(states, key=lambda s: s['credit']['id']):
        cid, csum, ccnt = st['credit']['id'], 0, 0
        for p, q in live_quotes(st):
            pend = [sp['amount_cents'] for sp in q['payments'] if pd(sp['due_date']) > as_of]
            if sum(pend) > 0:
                byp.append({'credit_id': cid, 'purchase_ref': p['ref'], 'pending_cents': sum(pend),
                            'payment_count': len(pend)})
            csum += sum(pend)
            ccnt += len(pend)
        if csum > 0:
            byc.append({'credit_id': cid, 'purchase_ref': '', 'pending_cents': csum, 'payment_count': ccnt})
        tot += csum
    return {'as_of': iso(as_of), 'total_cents': tot, 'by_purchase': byp, 'by_credit': byc}


def live_refs(st, at):
    return [p['ref'] for p, q in live_quotes(st)
            if pd(p['input']['date']) <= at <= pd(q['payments'][-1]['due_date'])]


# ---------------------------------------------------------------- cotización por monto y sensibilidad (§7)
def quote_by_amount(amount, r, n):
    i = apply_rate(amount, r)
    return {'amount_cents': amount, 'interest_cents': i, 'total_cents': amount + i,
            'reference_payment_cents': hu(amount * (1 + frac(r)) / n)}


def sensitivity(si):
    r, n, m = frac(si['rate']), si['term_count'], si['accrual_months']
    ext = si['external']['principal_cents']
    kt = ext + hu(ext * r)
    cm = si['cash_monthly_cents']
    rows = []
    for N in si['amounts_cents']:
        t = N * (1 + r)
        i = hu(N * r)
        rows.append({'amount_cents': N, 'interest_cents': i, 'total_cents': N + i,
                     'reference_payment_cents': hu(t / n), 'holder_total_cents': N + i - kt,
                     'holder_principal_cents': N - ext,
                     'holder_below_current': N - ext < si['holder_current_principal_cents'],
                     'monthly_purchase_cents': hu(t / m), 'monthly_total_cents': hu(t / m) + cm,
                     'holder_monthly_cents': hu((t - kt) / m), 'external_monthly_cents': hu(F(kt, m)) + cm})
    return rows


# ---------------------------------------------------------------- tarjeta y tienda con datos K18 (§5.3; PROVISIONAL)
def _k18(credit):
    if credit.get('kind') not in ('bank_card', 'store_credit'):
        falla('invalid', 'no es tarjeta ni tienda')


def _missing(credit):
    return any(credit.get(k) in (None, '') for k in ('annual_rate', 'rate_base', 'min_payment_pct_bp',
                                                     'min_payment_fixed_cents'))


def card_interest(credit, st, paid, paid_date):
    _k18(credit)
    if paid is not None and paid_date is not None and paid >= st['balance_cents'] \
            and pd(paid_date) <= pd(st['due_date'], False):
        return {'interest_cents': 0, 'provisional': []}
    if _missing(credit):
        falla('missing_capture')
    r = frac(credit['annual_rate'])
    if credit['rate_base'] == 'daily':
        base = hu(sum(st.get('daily_balances_cents') or []) * r / 360)
    elif credit['rate_base'] == 'monthly':
        base = hu(st['balance_cents'] * r / 12)
    else:
        falla('invalid', 'base')
    vat = credit.get('vat_on_interest_bp')
    iva = hu(F(base * vat, 10000)) if vat is not None else 0
    return {'interest_cents': base + iva, 'provisional': [P_TARJETA, P_IVA]}


def card_minimum(credit, st):
    _k18(credit)
    if _missing(credit):
        falla('missing_capture')
    bal = st['balance_cents']
    m = max(hu(F(bal * credit['min_payment_pct_bp'], 10000)), credit['min_payment_fixed_cents'])
    return {'minimum_cents': min(m, bal), 'provisional': [P_TARJETA, P_IVA]}


# ---------------------------------------------------------------- vistas con claves de model.json y data.json (§9.3)
def _purchases_of(st):
    return quote_state(st)


def _persona(ps):
    c = ps['composition'] or {}
    return {'capital': ps['principal_cents'], 'conjunto': c.get('joint_cents'), 'propio': c.get('own_cents'),
            'papel_dentro': c.get('paper_in_cents'), 'libre': c.get('free_cents'),
            'intereses': ps['interest_cents'], 'total': ps['total_cents'],
            'pago_por_quincena': ps['first_payment_cents']}


def vista_despensa(st, n):
    lst = _purchases_of(st)
    if not 1 <= n <= len(lst):
        falla('invalid', 'n')
    _, q = lst[n - 1]
    out = {'n': n, 'capital': q['principal_cents'], 'intereses': q['interest_cents'], 'total': q['total_cents'],
           'pagos': []}
    for sp in q['payments']:
        row = {'fecha': sp['due_date'], 'monto': sp['amount_cents']}
        for sh in sp['shares']:
            row[sh['person']] = sh['amount_cents']
        row['capital_repuesto'] = sp['principal_cents']
        out['pagos'].append(row)
    for ps in q['shares']:
        out[ps['person']] = _persona(ps)
    return out


def vista_disposicion(st, n):
    d = vista_despensa(st, n)
    out = {'despensa': n}
    for k, v in d.items():
        if isinstance(v, dict):
            out[k] = v
    return out


def vista_parametros(datos, st):
    _, q = _purchases_of(st)[0]
    ext = [s for s in q['shares'] if s['external']][0]
    cash = datos.get('cash') or []
    r = {'num': q['rate']['num'], 'den': q['rate']['den']}
    return {'tasa': r, 'tasa_pct': dict(r), 'papel': cash[0]['amount_cents'] if cash else None,
            'limite_credito': st['credit'].get('limit_cents'),
            'karen_fija': {'capital': ext['principal_cents'], 'intereses': ext['interest_cents'],
                           'total': ext['total_cents'], 'pago': ext['first_payment_cents'],
                           'conjunto': (ext['composition'] or {}).get('joint_cents')}}


def vista_credito(datos, st, today, frm, to):
    tl = timeline(st, datos.get('cash'), today, frm, to)
    tipo = {'start': 'inicio', 'purchase': 'compra', 'cash': 'papel', 'payment': 'pago'}
    return {'eventos': [{'fecha': e['date'], 'tipo': tipo[e['kind']], 'monto': e['amount_cents'],
                         'capital_repuesto': e['principal_restored_cents'],
                         'interes_no_repuesto': e['interest_not_restored_cents'],
                         'disponible_A': e['available_a_cents'], 'disponible_B': e['available_b_cents']}
                        for e in tl['events']],
            'comparacion_series': {'esperado_A': tl['series_a_cents'], 'esperado_B': tl['series_b_cents'],
                                   'calculado_A': tl['series_a_cents'], 'calculado_B': tl['series_b_cents'],
                                   'coinciden': True},
            'minimo_A': tl['min_a_cents'], 'minimo_B': tl['min_b_cents']}


def vista_sensibilidad(si):
    return [{'N': r['amount_cents'], 'intereses': r['interest_cents'], 'total': r['total_cents'],
             'pago_quincenal': r['reference_payment_cents'], 'axel_total': r['holder_total_cents'],
             'axel_capital': r['holder_principal_cents'], 'axel_menor_que_actual': r['holder_below_current'],
             'costo_mensual_total': r['monthly_total_cents'], 'axel_mensual': r['holder_monthly_cents'],
             'karen_mensual': r['external_monthly_cents']} for r in sensitivity(si)]


def vista_dos_quincenas(st, inp):
    q = quote_purchase(st['credit'], inp)
    out = {'total': q['total_cents'], 'intereses': q['interest_cents'], 'tasa_pct': dict(q['rate'])}
    for ps in q['shares']:
        out[ps['person']] = {'intereses': ps['interest_cents'], 'total': ps['total_cents'],
                             'por_quincena': ps['first_payment_cents']}
    return out


# ---------------------------------------------------------------- despacho de consultas (§9.2)
def evaluate(datos, fn, a):
    """Devuelve (salida, código). Códigos de §9.1."""
    try:
        check_input(datos)
        today = pd(datos.get('today'))
        sts = datos.get('credits') or []
        if fn == 'cotizar_compra':
            st = find_credit(datos, a['credit_id'])
            return quote_purchase(st['credit'], a['purchase']), ''
        if fn == 'promociones_aplicables':
            st = find_credit(datos, a['credit_id'])
            ps = applicable_promotions(st['credit'], pd(a['date'], False), a['principal_cents'], a['term_count'],
                                       a.get('store_id'), a.get('merchant'))
            return {'promotion_ids': [p['id'] for p in ps]}, ''
        if fn == 'linea_tiempo':
            st = find_credit(datos, a['credit_id'])
            return timeline(st, datos.get('cash'), pd(a['today'], False), pd(a['from'], False),
                            pd(a['to'], False), a.get('reading')), ''
        if fn == 'calendario_unificado':
            return unified_calendar(sts, pd(a['from'], False), pd(a['to'], False)), ''
        if fn == 'efectivo_ventana':
            frm, to = pd(a['from'], False), pd(a['to'], False)
            return cash_flow(unified_calendar(sts, frm, to), datos.get('cash'), frm, to), ''
        if fn == 'deuda_viva':
            return live_debt(sts, pd(a['as_of'], False)), ''
        if fn == 'compras_vivas':
            return {'refs': live_refs(find_credit(datos, a['credit_id']), pd(a['at'], False))}, ''
        if fn == 'cotizar_monto':
            return quote_by_amount(a['amount_cents'], a['rate'], a['term_count']), ''
        if fn == 'sensibilidad':
            return {'rows': sensitivity(a)}, ''
        if fn == 'interes_tarjeta':
            st = find_credit(datos, a['credit_id'])
            return card_interest(st['credit'], a['statement'], a.get('paid_cents'), a.get('paid_date')), ''
        if fn == 'minimo_tarjeta':
            st = find_credit(datos, a['credit_id'])
            return card_minimum(st['credit'], a['statement']), ''
        if fn == 'amortizar':
            return amortize(a['principal_cents'], a['annual_rate'], a.get('vat_bp'), a['term_count'],
                            a['frequency'], pd(a['first_date'])), ''
        if fn == 'comparar':
            import cred_escenarios
            return cred_escenarios.compare(datos, a['timeout_seconds'], a['request']), ''
        if fn == 'vista_parametros':
            return vista_parametros(datos, find_credit(datos, a['credit_id'])), ''
        if fn == 'vista_despensa':
            return vista_despensa(find_credit(datos, a['credit_id']), a['n']), ''
        if fn == 'vista_disposicion':
            return vista_disposicion(find_credit(datos, a['credit_id']), a['n']), ''
        if fn == 'vista_credito':
            return vista_credito(datos, find_credit(datos, a['credit_id']), pd(a['today'], False),
                                 pd(a['from'], False), pd(a['to'], False)), ''
        if fn == 'vista_sensibilidad':
            return vista_sensibilidad(a), ''
        if fn == 'vista_interes_por_monto':
            return [{'monto': x['amount_cents'], 'intereses': x['interest_cents'], 'total': x['total_cents'],
                     'pago_quincenal': x['reference_payment_cents']}
                    for x in (quote_by_amount(m, a['rate'], a['term_count']) for m in a['amounts_cents'])], ''
        if fn == 'vista_dos_quincenas':
            return vista_dos_quincenas(find_credit(datos, a['credit_id']), a['purchase']), ''
        return None, 'invalid'
    except CredError as e:
        return None, e.code
    except (KeyError, TypeError) as e:
        return None, 'invalid'
