"""Tabla a mano de los 3 casos del comparador (S3-12, S3-13, S3-14; S3-interfaz.md §8, §10.1).

Cálculo independiente del motor Go: fracciones exactas, día por día, solo biblioteca estándar.
No importa ni lee nada del motor. La única fuente de la tabla son las fixtures
oraculo/fixtures/s3/s3-comparador-{caso-1,caso-2,caso-3,sin-restricciones,tiempo}.json (con sus
_datos/comparador-caso-N.json): este script ya no escribe copias; recalcula cada caso y comprueba
que la fixture es idéntica a la tabla a mano (sale con 1 si no). La tabla en Markdown va a la ruta
que reciba como argumento (para docs/estado/S3-escenarios.md).
Uso: python3 -I tabla_mano.py [<salida.md>]
"""
import json, math, sys
from datetime import date, timedelta
from fractions import Fraction as F
from pathlib import Path

D = date.fromisoformat
HERE = Path(__file__).resolve().parent
OBJ_MED = "S3-objetivos-medida"


def hu(fr):  # ROUND_HALF_UP de una fracción no negativa a entero
    fr = F(fr)
    return math.floor(fr + F(1, 2))


def last_day(y, m):
    return (date(y + (m == 12), m % 12 + 1, 1) - timedelta(days=1)).day


def clip(y, m, d):
    y, m = y + (m - 1) // 12, (m - 1) % 12 + 1
    return date(y, m, min(d, last_day(y, m)))


def fq(d):  # quincena: (año, mes, mitad)
    return (d.year, d.month, 1 if d.day <= 15 else 2)


def fq_next(q):
    y, m, h = q
    return (y, m, 2) if h == 1 else ((y + 1, 1, 1) if m == 12 else (y, m + 1, 1))


def pay_days_from(d):  # días de pago 15 y 30 (o último) desde d inclusive
    x = d
    while True:
        p = date(x.year, x.month, 15) if x.day <= 15 else clip(x.year, x.month, 30)
        if p >= x:
            yield p
            x = p + timedelta(days=1)
        else:
            x = date(x.year + (x.month == 12), x.month % 12 + 1, 1)


def sched_fortnight(d, min_days, n):
    out = []
    for p in pay_days_from(d):
        if out or (p - d).days >= min_days:
            out.append(p)
        if len(out) == n:
            return out


def card_dues(c, d, n):
    sd, pd = c["statement_day"], c["payment_due_day"]
    st = clip(d.year, d.month, sd)
    if st < d:
        st = clip(d.year, d.month + 1, sd)
    out = []
    for i in range(n):
        s = clip(st.year, st.month + i, sd)
        due = clip(s.year, s.month, pd)
        if due <= s:
            due = clip(s.year, s.month + 1, pd)
        out.append(due)
    return out


def rate(r):
    return F(r["num"], r["den"])


def promo_ok(p, c, d, amount, term, store_id, merchant):
    if not p["active"] or (p["credit_id"] is not None and p["credit_id"] != c["id"]):
        return False
    if p["store_id"] is not None and store_id is not None:
        if p["store_id"] != store_id:
            return False
    elif p["merchant"].strip().lower() != merchant.strip().lower():
        return False
    return D(p["valid_from"]) <= d <= D(p["valid_to"]) and amount >= p["min_purchase_cents"] and term == p["term_count"]


def quote(c, d, amount, term, promo_id, store_id, merchant, quoted=None):
    """Cotización propia (RN-06, RN-07, RN-08, §5.1, §5.3). None = no se genera."""
    if c["kind"] == "revolving_line":
        if quoted is not None:
            interest = quoted
        elif promo_id is not None:
            p = next(x for x in c["promotions"] if x["id"] == promo_id)
            interest = hu(amount * rate(p["rate"]))
        else:
            tr = [t for t in c["term_rates"] if t["term_count"] == term]
            if not tr:
                return None
            interest = hu(amount * rate(tr[0]["rate"]))
        total = amount + interest
        base = total // term
        amounts = [base] * (term - 1) + [total - base * (term - 1)]
        cb = hu(F(amount, term))
        caps = [cb] * (term - 1) + [amount - cb * (term - 1)]
        dues = sched_fortnight(d, c["first_payment_min_days"], term)
        prov = ["S3-serie-sin-intereses"] if interest == 0 and term >= 2 else []
    elif c["kind"] == "bank_card":
        interest = 0
        if promo_id is not None:
            base = amount // term
            amounts = [base] * (term - 1) + [amount - base * (term - 1)]
            caps = list(amounts)
            prov = ["S3-formula-tarjeta"] + (["S3-serie-sin-intereses"] if term >= 2 else [])
        else:
            amounts, caps = [amount], [amount]
            prov = ["S3-formula-tarjeta", "S3-origen-cargo-normal"]
        dues = card_dues(c, d, term)
    else:
        return None
    return {"date": d, "interest": interest, "pays": list(zip(dues, amounts, caps)), "prov": prov}


def existing(case):
    """Compras existentes cotizadas: {credit_id: [(ref, compra)]} (solo titular; sin abonos registrados)."""
    out = {}
    for cs in case["datos"]["credits"]:
        c = cs["credit"]
        out[c["id"]] = []
        for ps in cs["purchases"]:
            if ps["status"] == "cancelled":
                continue
            i = ps["input"]
            q = quote(c, D(i["date"]), i["principal_cents"], i["term_count"], i["promotion_id"], i["store_id"], i["merchant"], i["quoted_interest_cents"])
            out[c["id"]].append((ps["ref"], q, i["principal_cents"]))
    return out


def availability(c, purchases, today, P, L):
    """Disponible lectura B día por día (RN-09): mínimo en [P, L]."""
    limit = c["limit_cents"]
    b = limit
    for (_, q, k) in purchases:
        if q["date"] <= today:
            b -= k
        for (due, _, cap) in q["pays"]:
            if due >= today and due <= today:
                b += cap
    deltas = {}
    for (_, q, k) in purchases:
        if q["date"] > today:
            deltas[q["date"]] = deltas.get(q["date"], 0) - k
        for (due, _, cap) in q["pays"]:
            if due > today:
                deltas[due] = deltas.get(due, 0) + cap
    day, best = today, None
    while day <= L:
        if day > today:
            b = min(limit, b + deltas.get(day, 0))
        if day >= P:
            best = b if best is None else min(best, b)
        day += timedelta(days=1)
    return best


def compare(case):
    rq = case["request"]
    today = D(case["datos"]["today"])
    g, o = rq["global"], rq["override"]
    keys = ["min_available_buffer_cents", "max_live_purchases", "zero_debt_by_date", "date_window_days"]
    used = {k: (o[k] if o[k] is not None else g[k]) for k in keys}
    missing = [k for k in keys if used[k] is None]
    if missing:
        return {"status": "missing_constraints", "missing_constraints": missing, "constraints_used": used,
                "evaluated_count": 0, "valid": [], "optima": None, "recommendation": None,
                "rejected": [], "excluded": [], "provisional": []}
    ex = existing(case)
    amount = rq["amount_cents"]
    w = used["date_window_days"]
    desired = D(rq["desired_date"])
    if rq["store_id"] is not None or rq["merchant_name"]:
        stores = [(rq["store_id"], rq["merchant_name"])]
    else:
        stores = [(b["store_id"], b["merchant_name"]) for b in rq["baskets"]]
    baskets = rq["baskets"]
    caps = {pc["person"]: pc["cap_cents"] for pc in rq["caps"]}
    cands, excluded = [], []
    credits = sorted(case["datos"]["credits"], key=lambda cs: cs["credit"]["id"])
    for cs in credits:
        c = cs["credit"]
        if c["limit_cents"] is None:
            excluded.append({"credit_id": c["id"], "reason": "no_limit"})
            continue
        for k in range(-w, w + 1):
            d = desired + timedelta(days=k)
            if d < today:
                continue
            for (sid, mer) in stores:
                opts = []
                if c["kind"] == "bank_card":
                    opts.append((1, None))
                    for p in sorted(c["promotions"], key=lambda p: p["id"]):
                        if promo_ok(p, c, d, amount, p["term_count"], sid, mer) and p["rate"]["num"] == 0:
                            opts.append((p["term_count"], p["id"]))
                elif c["kind"] == "revolving_line":
                    terms = {t["term_count"] for t in c["term_rates"]}
                    for p in c["promotions"]:
                        if promo_ok(p, c, d, amount, p["term_count"], sid, mer):
                            terms.add(p["term_count"])
                    for t in sorted(terms):
                        if any(x["term_count"] == t for x in c["term_rates"]):
                            opts.append((t, None))
                        for p in sorted(c["promotions"], key=lambda p: p["id"]):
                            if promo_ok(p, c, d, amount, t, sid, mer):
                                opts.append((t, p["id"]))
                for (t, pid) in opts:
                    q = quote(c, d, amount, t, pid, sid, mer)
                    if q is not None:
                        cands.append((c, t, d, sid, mer, pid, q))
    valid, rejected = [], []
    for (c, t, d, sid, mer, pid, q) in cands:
        P, L = d, q["pays"][-1][0]
        mine = ex[c["id"]] + [("_escenario", q, amount)]
        min_av = availability(c, mine, today, P, L)
        live = 0
        day = P
        while day <= L:
            live = max(live, sum(1 for (_, x, _) in mine if x["date"] <= day <= x["pays"][-1][0]))
            day += timedelta(days=1)
        zero = max(x["pays"][-1][0] for (_, x, _) in mine)
        # serie del objetivo 2: quincenas de fq(P) a fq(L) con las de 0
        qs, x = [], fq(P)
        while True:
            qs.append(x)
            if x == fq(L):
                break
            x = fq_next(x)
        tot = {qq: 0 for qq in qs}
        provs = {OBJ_MED} | set(q["prov"])
        allp = [(cc, xx) for cc in ex for (_, xx, _) in ex[cc]] + [(c["id"], q)]
        for (_, xx) in allp:
            hit = False
            for (due, amt, _) in xx["pays"]:
                if fq(due) in tot:
                    tot[fq(due)] += amt
                    hit = True
            if hit:
                provs |= set(xx["prov"])
        xs = [tot[qq] for qq in qs]
        n = len(xs)
        N = 10**6 * (n * sum(v * v for v in xs) - sum(xs) ** 2)
        std = ((math.isqrt(4 * N) // n) + 1) // 2
        over = []
        if "axel" in caps and max(xs) > caps["axel"]:  # todas las compras de los casos son solo del titular axel
            over.append({"person": "axel", "peak_cents": max(xs), "cap_cents": caps["axel"]})
        bq = None
        for b in baskets:
            if (sid is not None and b["store_id"] == sid) or (sid is None and b["merchant_name"].strip().lower() == mer.strip().lower()):
                bq = b
                break
        vpp = None if bq is None else hu(F(amount * bq["base_units_milli"] * 100, bq["cost_cents"] * (amount + q["interest"])))
        sc = {"credit_id": c["id"], "term_count": t, "purchase_date": P.isoformat(), "store_id": sid,
              "promotion_id": pid, "interest_cents": q["interest"], "max_fortnight_payment_cents": max(xs),
              "fortnight_payment_spread_cents": max(xs) - min(xs), "std_dev_milli_cents": std,
              "min_available_cents": min_av, "max_live_purchases": live, "zero_debt_date": zero.isoformat(),
              "value_per_peso_milli": vpp, "score_bp": None, "over_cap": over, "provisional": sorted(provs)}
        viol = None
        if min_av < used["min_available_buffer_cents"]:
            viol = "min_available_buffer"
        elif live > used["max_live_purchases"]:
            viol = "max_live_purchases"
        elif zero > D(used["zero_debt_by_date"]):
            viol = "zero_debt_by_date"
        if viol:
            rejected.append({"scenario": sc, "violated_constraint": viol})
        else:
            valid.append(sc)
    key = lambda s: (s["credit_id"], s["term_count"], s["purchase_date"],
                     (0, 0) if s["store_id"] is None else (1, s["store_id"]),
                     (0, 0) if s["promotion_id"] is None else (1, s["promotion_id"]))
    valid.sort(key=key)
    rejected.sort(key=lambda r: key(r["scenario"]))
    prov = {OBJ_MED}
    for s in valid + [r["scenario"] for r in rejected]:
        prov |= set(s["provisional"])
    res = {"status": "ok" if valid else "no_valid_scenario", "missing_constraints": [], "constraints_used": used,
           "evaluated_count": len(valid) + len(rejected), "valid": valid, "optima": None, "recommendation": None,
           "rejected": rejected, "excluded": excluded, "provisional": sorted(prov)}
    if not valid:
        return res
    wts = rq["weights"]
    W = [wts["lowest_interest"], wts["even_low_payment"], wts["max_available"], wts["most_value_per_peso"]]
    fields = [("interest_cents", -1), ("std_dev_milli_cents", -1), ("min_available_cents", 1), ("value_per_peso_milli", 1)]

    def norm(field, sign):
        vals = [s[field] for s in valid if s[field] is not None]
        if not vals or max(vals) == min(vals):
            return [F(1)] * len(valid)
        hi, lo = max(vals), min(vals)
        return [F(1) if s[field] is None else (F(hi - s[field], hi - lo) if sign < 0 else F(s[field] - lo, hi - lo)) for s in valid]
    ns = [norm(f, sg) for (f, sg) in fields]
    exact = []
    for i, s in enumerate(valid):
        sc = sum(F(W[j]) * ns[j][i] for j in range(4)) / sum(W)
        exact.append(sc)
        s["score_bp"] = hu(sc * 10000)
    order = list(range(len(valid)))

    def best(j):
        f, sg = fields[j]
        def k(i):
            v = valid[i][f]
            v = (-sg * v) if v is not None else math.inf
            tie = len(valid[i]["over_cap"]) if j == 1 else 0
            return (v, tie, -exact[i], key(valid[i]))
        return valid[min(order, key=k)]
    names = ["lowest_interest", "even_low_payment", "max_available", "most_value_per_peso"]
    res["optima"] = {names[j]: best(j) for j in range(4)}
    res["recommendation"] = valid[min(order, key=lambda i: (-exact[i], key(valid[i])))]
    return res


# ---------------------------------------------------------------- casos (capturas de fixture, no de Axel)
def R(num, den, dec):
    return {"num": num, "den": den, "decimales": dec}


T2, T4 = R(53450, 309450, 2), R(80050, 309450, 4)  # tabla semilla por plazo (decisión C)


def credit(id_, name, kind, limit, term_rates=(), promotions=(), **kw):
    c = {"id": id_, "name": name, "kind": kind, "holder": "axel", "limit_cents": limit,
         "frequency": "fortnightly" if kind == "revolving_line" else "monthly",
         "first_payment_min_days": 15 if kind == "revolving_line" else None,
         "reading": "B", "reading_confirmed": False, "reading_a_anchor": "2026-11-01",
         "statement_day": None, "payment_due_day": None, "annual_rate": None, "rate_base": "",
         "vat_on_interest_bp": None, "min_payment_pct_bp": None, "min_payment_fixed_cents": None,
         "term_rates": [{"term_count": t, "rate": r, "source": "manual", "quote_date": None} for (t, r) in term_rates],
         "promotions": list(promotions)}
    c.update(kw)
    return c


def promo(id_, credit_id, merchant, store_id, term, frm, to, minimum):
    return {"id": id_, "credit_id": credit_id, "merchant": merchant, "store_id": store_id, "term_count": term,
            "rate": R(0, 1, 0), "min_purchase_cents": minimum, "valid_from": frm, "valid_to": to, "active": True}


def purchase(ref, d, amount, term, quoted=None):
    return {"ref": ref, "status": "active", "input": {
        "date": d, "principal_cents": amount, "term_count": term, "quoted_interest_cents": quoted,
        "promotion_id": None, "store_id": None, "merchant": "", "first_payment_date": None,
        "shares": [{"person": "axel", "external": False, "principal_cents": amount, "joint_cents": None, "paper_in_cents": None}]}}


def state(c, purchases):
    return {"credit": c, "purchases": purchases, "payments": [], "statements": []}


BASKETS = [{"store_id": 1, "merchant_name": "Tienda Barata", "base_units_milli": 60000, "cost_cents": 200000},
           {"store_id": 2, "merchant_name": "Tienda Cara", "base_units_milli": 40000, "cost_cents": 200000}]


def request(amount, desired, glob, over=None, caps=(), baskets=BASKETS):
    nul = {"min_available_buffer_cents": None, "max_live_purchases": None, "zero_debt_by_date": None, "date_window_days": None}
    return {"amount_cents": amount, "desired_date": desired, "store_id": None, "merchant_name": "",
            "shares": [{"person": "axel", "external": False, "principal_cents": amount, "joint_cents": None, "paper_in_cents": None}],
            "global": glob, "override": over or nul,
            "weights": {"lowest_interest": 1, "even_low_payment": 1, "max_available": 1, "most_value_per_peso": 1},
            "caps": [{"person": "axel", "cap_cents": c} for c in caps], "baskets": baskets}


def G(buf, live, zero, w):
    return {"min_available_buffer_cents": buf, "max_live_purchases": live, "zero_debt_by_date": zero, "date_window_days": w}


T6 = R(35, 100, 0)  # tasa de fixture para 6 quincenas (captura de fixture, no de Axel)
L2RATES = [(2, R(8, 100, 0)), (4, R(15, 100, 0))]


def case13_credits(extra=()):
    return [state(credit(1, "Línea Uno", "revolving_line", 1500000, [(2, T2), (4, T4), (6, T6)],
                         [promo(1, 1, "Tienda Cara", 2, 4, "2026-11-28", "2026-11-29", 100000)]),
                  [purchase("L1-f", "2027-01-20", 250000, 1, quoted=0)]),
            state(credit(2, "Línea Dos", "revolving_line", 600000, L2RATES),
                  [purchase("L2-a", "2026-10-25", 116550, 2, quoted=9318), purchase("L2-b", "2026-11-01", 100000, 4)])] + list(extra)


CASES = {
    # Caso 1: dos líneas tipo Kueski.
    "caso-1": {"datos": {"today": "2026-11-01", "cash": [], "credits": case13_credits()},
               "request": request(200000, "2026-11-30", G(300000, 2, "2027-02-15", 1), caps=[120000])},
    # Caso 2: línea + tarjeta con meses sin intereses; el cargo normal (0.00) compite con la promoción.
    "caso-2": {"datos": {"today": "2026-11-01", "cash": [], "credits": [
        state(credit(1, "Línea Uno", "revolving_line", 800000, [(2, T2), (4, T4), (6, T6)]),
              [purchase("L1-a", "2026-10-20", 200000, 4), purchase("L1-b", "2026-10-25", 100000, 2),
               purchase("L1-f", "2027-01-20", 260000, 2)]),
        state(credit(3, "Tarjeta", "bank_card", 650000, statement_day=5, payment_due_day=20,
                     promotions=[promo(2, 3, "Tienda Cara", 2, 2, "2026-11-29", "2026-11-30", 200000)]),
              [])]},
        "request": request(300000, "2026-11-30", G(320000, 2, "2027-02-28", 1), caps=[240000],
                           baskets=[{**BASKETS[0], "cost_cents": 300000, "base_units_milli": 90000},
                                    {**BASKETS[1], "cost_cents": 300000, "base_units_milli": 72000}])},
    # Caso 3: tope por persona capturado, Override que cambia restricciones globales y un crédito sin límite.
    "caso-3": {"datos": {"today": "2026-11-01", "cash": [], "credits": case13_credits(
        [state(credit(4, "Crédito sin límite", "revolving_line", None, [(2, T2)]), [])])},
        "request": request(180000, "2026-11-30", G(350000, 2, "2026-12-31", 1),
                           over=G(300000, None, "2027-02-15", None), caps=[100000])},
}


def k(s):
    return f'{s["credit_id"]}/{s["term_count"]}q/{s["purchase_date"]}/t{s["store_id"]}/p{s["promotion_id"]}'


def money(c):
    return "—" if c is None else f"{c / 100:,.2f}"


def checks(res):
    out = []
    viol = {r["violated_constraint"] for r in res["rejected"]}
    out.append(("un descartado por cada restricción", viol == {"min_available_buffer", "max_live_purchases", "zero_debt_by_date"}))
    out.append(("uno válido sobre el tope", any(s["over_cap"] for s in res["valid"])))
    op = res["optima"]
    out.append(("4 óptimos distintos", len({k(op[n]) for n in op}) == 4))
    out.append(("(1) intereses = la promoción", op["lowest_interest"]["promotion_id"] is not None))
    return out


def table(res):
    rows = ["| Llave (crédito/plazo/fecha/tienda/promoción) | Intereses | Desv. (milésimas de centavo) | Pico | Disponible mín. | Vivas | Deuda en cero | Unid. por peso (milésimas) | Puntaje (pb) | Sobre tope | Resultado |",
            "|---|---:|---:|---:|---:|---:|---|---:|---:|---|---|"]
    allsc = [(s, "válido") for s in res["valid"]] + [(r["scenario"], "descartado: " + r["violated_constraint"]) for r in res["rejected"]]
    for s, st in allsc:
        oc = ", ".join(f'{o["person"]} {money(o["peak_cents"])} > {money(o["cap_cents"])}' for o in s["over_cap"]) or "—"
        rows.append(f'| {k(s)} | {money(s["interest_cents"])} | {s["std_dev_milli_cents"]:,} | {money(s["max_fortnight_payment_cents"])} | '
                    f'{money(s["min_available_cents"])} | {s["max_live_purchases"]} | {s["zero_debt_date"]} | '
                    f'{"—" if s["value_per_peso_milli"] is None else format(s["value_per_peso_milli"], ",")} | '
                    f'{"—" if s["score_bp"] is None else format(s["score_bp"], ",")} | {oc} | {st} |')
    return "\n".join(rows)


PROV_Q = ("serie y desempate del objetivo 2, ventana del colchón y del objetivo 3, tienda y canasta del objetivo 4 "
          "y créditos sin límite (preguntas 1 a 4)")
OTROS = {"S3-formula-tarjeta": "docs/estado/bloqueos/S3-formula-tarjeta.md (calendario de la tarjeta, pregunta 4)",
         "S3-origen-cargo-normal": "docs/estado/bloqueos/S3-origen-cargo-normal.md (`source` del cargo normal)",
         "S3-serie-sin-intereses": "docs/estado/bloqueos/S3-serie-sin-intereses.md (serie de capital de la promoción 0 % y de MSI)"}


def fixture(name, desc, datos, args, esperado, prov):
    return {"formato": "fixture/1", "motor": "s3", "caso": f"s3-comparador-{name}", "descripcion": desc,
            "fuente": [{"tipo": "manual", "ruta": "docs/estado/S3-escenarios.md"}, {"tipo": "criterio", "ruta": "S3-12"},
                       {"tipo": "criterio", "ruta": "S3-13"}, {"tipo": "criterio", "ruta": "S3-14"}],
            "provisional": prov, "entradas": {"datos": datos, "consultas": [{"id": name, "funcion": "comparar", "argumentos": args}]},
            "esperado": {name: esperado}}


FIX = HERE.parents[4] / "oraculo" / "fixtures" / "s3"


def check(fx):
    """Compara la fixture recalculada con la de oraculo/fixtures/s3 (resuelve $ref); no escribe nada."""
    p = FIX / f"{fx['caso']}.json"
    if not p.exists():
        print("    NO ", f"falta oraculo/fixtures/s3/{p.name}")
        return False
    cur = json.loads(p.read_text(encoding="utf-8"))
    d = cur["entradas"]["datos"]
    if isinstance(d, dict) and "$ref" in d:
        cur["entradas"]["datos"] = json.loads((FIX / d["$ref"]).read_text(encoding="utf-8"))
    ok = cur == fx
    print("   ", "OK " if ok else "NO ", f"oraculo/fixtures/s3/{p.name} = tabla a mano")
    return ok


def main():
    md, good = [], True
    for name, case in CASES.items():
        res = compare(case)
        cs = checks(res)
        print(name, res["status"], "evaluados", res["evaluated_count"], "válidos", len(res["valid"]),
              "descartados", len(res["rejected"]), "excluidos", res["excluded"])
        for (t, ok) in cs:
            print("   ", "OK " if ok else "NO ", t)
        for n, s in res["optima"].items():
            print("    óptimo", n, k(s), s["score_bp"])
        print("    recomendación", k(res["recommendation"]), res["recommendation"]["score_bp"])
        others = [OTROS[p] for p in res["provisional"] if p in OTROS]
        prov = {"bloqueo": "docs/estado/bloqueos/S3-objetivos-medida.md",
                "que": PROV_Q + ("; otros bloqueos: " + ", ".join(others) if others else "")}
        esp = {key: res[key] for key in ("status", "evaluated_count", "valid", "rejected", "optima", "recommendation", "excluded", "provisional")}
        fx = fixture(name, f"Comparador, {name}: tabla a mano de docs/estado/S3-escenarios.md (S3-12, S3-13, S3-14).",
                     case["datos"], {"timeout_seconds": 5, "request": case["request"]}, esp, prov)
        good &= check(fx)
        md.append(f"### Tabla a mano · {name}\n\n`status` {res['status']}, `evaluated_count` {res['evaluated_count']} "
                  f"({len(res['valid'])} válidos, {len(res['rejected'])} descartados), excluidos: "
                  f"{', '.join(str(e['credit_id']) + ' ' + e['reason'] for e in res['excluded']) or 'ninguno'}.\n\n" + table(res) +
                  "\n\nÓptimos: " + "; ".join(f"{n} → `{k(s)}`" for n, s in res["optima"].items()) +
                  f". Recomendación → `{k(res['recommendation'])}` ({res['recommendation']['score_bp']:,} pb). "
                  f"`provisional` del resultado: {', '.join(res['provisional'])}.\n")
    c1 = CASES["caso-1"]
    nul = {"min_available_buffer_cents": None, "max_live_purchases": None, "zero_debt_by_date": None, "date_window_days": None}
    rq = {**c1["request"], "global": nul, "override": nul}
    res = compare({"datos": c1["datos"], "request": rq})
    esp = {key: res[key] for key in ("status", "missing_constraints", "evaluated_count", "valid", "rejected", "optima", "recommendation", "provisional")}
    fx = fixture("sin-restricciones", "Caso 1 sin restricciones capturadas: «faltan restricciones» y 0 recomendaciones (K16).",
                 c1["datos"], {"timeout_seconds": 5, "request": rq}, esp, None)
    good &= check(fx)
    fx = fixture("tiempo", "Caso 1 con timeout_seconds 0: vence de inmediato, 503 unavailable (K43, K44).",
                 c1["datos"], {"timeout_seconds": 0, "request": c1["request"]}, {"error": "timeout"}, None)
    good &= check(fx)
    if len(sys.argv) > 1:
        Path(sys.argv[1]).write_text("\n".join(md))
    return good


if __name__ == "__main__":
    sys.exit(0 if main() else 1)
