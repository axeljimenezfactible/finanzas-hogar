#!/usr/bin/env python3
"""Verifica docs/design/maquetas/datos-sinteticos.json (P3-12 y consistencia). Uso: python3 -I verificar_sinteticos.py [raiz_repo]
Termina con «RESULTADO: OK» o sale con código 1."""
import sys, re, json, calendar, math
from datetime import date
RAIZ = sys.argv[1] if len(sys.argv) > 1 else "/home/deepsec/finanzas-hogar"
JS = RAIZ + "/docs/design/maquetas/datos-sinteticos.json"
TXT = open(JS, encoding="utf-8").read(); D = json.loads(TXT); HOY = date.fromisoformat(D["meta"]["hoy"])
fallos, n = [], [0]
def ok(cond, msg):
    n[0] += 1; print(("OK   " if cond else "FALLA"), f"{n[0]:02d}", msg)
    if not cond: fallos.append(msg)
def rh(a, b): return (2 * a + b) // (2 * b)
PROHIB = r"3,?094\.50|800\.50|598\.82|374\.93|5,?034\.74|1,?258\.68|883\.75|6,?293\.44|973\.75|1,?258\.70|883\.77|2,?395\.28|1,?499\.72|2,?616\.87|1,?767\.51|849\.36|8,?094\.50|14,?362\.48|13,?103\.78|25\.8685|17\.27"
CENTS = [309450, 80050, 59882, 37493, 503474, 125868, 88375, 629344, 97375, 125870, 88377, 239528, 149972, 261687, 176751, 84936, 809450, 1436248, 1310378]
ok(len(re.findall(PROHIB, TXT)) == 0, "grep de las cifras prohibidas de P3-12 (con/sin coma) = 0 coincidencias")
nums = {int(x) for x in re.findall(r"(?<![\d.])-?\d+(?![\d.]|\w)", TXT)}
ok(not (nums & set(CENTS)), "ningún valor en centavos coincide con las cifras canon (309450, 80050, 59882, ...)")
ok(len(re.findall(r"(?i)(password|contraseña) *[:=] *[^ *]{6,}|bot[_-]?token *[:=] *[^ *]{6,}|api[_-]?key *[:=] *[^ *]{6,}|[0-9]{8,}:[A-Za-z0-9_-]{30,}|@[a-z0-9-]+\.[a-z]{2,}", TXT)) == 0, "sin credenciales ni correos")
# ---- presupuesto y categorías
C = D["categorias"]; G = D["gastos"]; R = [g for g in G if g["estado"] == "realizado"]
ok(sum(c["plan_cents"] for c in C) == D["presupuesto"]["resumen"]["plan_gasto_cents"], "suma del plan por categoría = plan total del presupuesto")
def hog(g): return sum(x["hogar_cents"] for x in (g["devengado"] or [dict(mes=g["fecha"][:7], hogar_cents=g["hogar_cents"])]) if x["mes"] == "2026-10")
ok(all(sum(hog(g) for g in R if g["categoria_id"] == c["id"]) == c["real_cents"] for c in C) and sum(c["real_cents"] for c in C) == D["presupuesto"]["resumen"]["real_cents"], "real por categoría = Σ parte del hogar devengada de los gastos; suma = real total")
ok(all(c["usado_pct"] == rh(c["real_cents"] * 100, c["plan_cents"]) and c["estado"] == ("excedido" if c["usado_pct"] >= 100 else "aviso" if c["usado_pct"] >= 80 else "normal") and c["diferencia_cents"] == c["plan_cents"] - c["real_cents"] for c in C), "porcentaje entero (ROUND_HALF_UP), estado y diferencia coinciden con plan/real")
ok(any(80 <= c["usado_pct"] < 100 for c in C) and any(c["usado_pct"] >= 100 for c in C), "hay categorías en aviso (80 a 99 %) y excedidas (100 % o más)")
rs = D["presupuesto"]["resumen"]; ok(rs["sin_asignar_cents"] == rs["ingreso_mes_cents"] - rs["plan_gasto_cents"] - rs["ahorro_planificado_cents"] and rs["ingreso_mes_cents"] == D["ingresos"]["total_mes_cents"], "sin asignar = ingreso − plan − ahorro planificado; ingreso del mes = quincenas + extras")
# ---- gastos
ok(len(G) >= 25 and {"fijo", "variable", "venidero"} <= {g["tipo"] for g in G}, f"{len(G)} gastos (≥ 25) con fijos, variables y venideros")
ok(all(g["hogar_cents"] + g["karen_cents"] == g["monto_cents"] for g in G), "en todos los gastos parte del hogar + parte de Karen = monto")
ok(any(g["reparto"].get("hogar_pct") == 50 for g in G) and any(g["reparto"].get("hogar_pct") not in (None, 50) for g in G), "hay reparto Conjunto 50/50 y otro con porcentaje distinto")
ok(all(sum(x["hogar_cents"] + x["karen_cents"] for x in g["devengado"]) >= g["monto_cents"] for g in G if g["devengado"]) and any(g["efectivo"] for g in G), "gastos a crédito: devengado (con intereses) en meses distintos al efectivo (pagos)")
ok(abs(D["gastos_resumen"]["suma_hogar_cents"] + D["gastos_resumen"]["suma_karen_cents"] - D["gastos_resumen"]["suma_montos_cents"]) == 0, "totales del filtro: hogar + Karen = suma de montos")
# ---- créditos
K = D["creditos"]; PG = K["pagos_programados"]
for c in K["kueski"]["compras"] + K["kueski"]["compras_programadas"]:
    ps = [p for p in PG if p["compra_id"] == c["id"]]
    assert c["capital_cents"] + c["intereses_cents"] == c["total_cents"]
    assert sum(p["monto_cents"] for p in ps) == c["total_cents"] and sum(p["hogar_cents"] for p in ps) == c["hogar"]["total_cents"] and sum(p["karen_cents"] for p in ps) == c["karen"]["total_cents"], c["id"]
    assert c["hogar"]["capital_cents"] + c["karen"]["capital_cents"] == c["capital_cents"] and sum(p["capital_repuesto_cents"] for p in ps) == c["capital_cents"], c["id"]
ok(True, "compras Kueski: capital + intereses = total; los pagos suman el total (hogar y Karen por separado) y el capital repuesto suma el capital")
ok(all(p["monto_cents"] == p["hogar_cents"] + p["karen_cents"] for p in PG), "cada pago programado: hogar + Karen = monto")
for k, key in (("prestamo", "amortizacion"), ("tienda", "abonos")):
    L = K[k][key]; cap = K[k].get("monto_original_cents", K[k].get("total_cents"))
    assert sum(r["capital_cents"] for r in L) == cap and all(r["pago_cents"] == r["capital_cents"] + r["interes_cents"] for r in L) and L[-1]["saldo_cents"] == 0, k
ok(True, "préstamo y crédito de tienda: los pagos suman capital + intereses, el capital suma el monto y el saldo final es 0")
ev = K["kueski"]["serie_disponible"]["eventos"]; ok(K["kueski"]["serie_disponible"]["min_b_cents"] == min(e["disponible_b_cents"] for e in ev) and K["kueski"]["serie_disponible"]["min_a_cents"] == min(e["disponible_a_cents"] for e in ev) and K["kueski"]["disponible_hoy"]["B_cents"] < K["kueski"]["disponible_hoy"]["A_cents"], "serie de disponible: mínimos A y B coinciden con los eventos; B < A hoy")
ok(K["deuda_viva_total_cents"] == sum(l["deuda_viva_cents"] for l in K["lineas"]) and {l["tipo"] for l in K["lineas"]} == {"revolvente", "tarjeta", "prestamo", "tienda"}, "deuda viva total = Σ líneas; hay revolvente, tarjeta, préstamo y tienda")
pr = D["inicio"]["tarjetas"]["proximo_pago"]["next_payment"]; ok(pr["amount_cents"] == pr["household_share_cents"] + pr["external_share_cents"] and pr["payment_count"] == len([p for p in PG if p["fecha"] == pr["due_date"] and p["hogar_estado"] != "pagada" and p["compra_id"] not in ("kc3", "kc4")]), "tarjeta Próximo pago: monto = hogar + Karen y cuenta los pagos de ese día")
ok({"due_date", "credit_id", "amount_cents", "household_share_cents", "external_share_cents", "days_left", "payment_count", "is_overdue"} <= set(pr) and pr["days_left"] == (date.fromisoformat(pr["due_date"]) - HOY).days and pr["is_overdue"] is False, "tarjeta Próximo pago: usa los nombres de DashboardNextPayment (C1) y los días que faltan salen de la fecha")
pv = D["estados_alternos"]["pago_vencido"]; ok(pv["is_overdue"] and pv["days_left"] == (date.fromisoformat(pv["due_date"]) - HOY).days < 0, "estado alterno pago vencido: nombres de C1 y días vencidos coherentes con hoy")
fq = D["inicio"]["tarjetas"]["flujo_quincenal"]; ok(fq["closing_balance_cents"] == fq["income_cents"] - fq["spent_cents"] - fq["committed_cents"], "tarjeta Flujo quincenal: saldo al cierre = ingresos − hechos − comprometido")
mv = D["inicio"]["tarjetas"]["mes_vs_plan"]; ok(mv["month"] == D["meta"]["mes_en_curso"] and "mes" not in mv and mv["used_pct"] == rh(mv["actual_cents"] * 100, mv["planned_cents"]) and mv["remaining_cents"] == mv["planned_cents"] - mv["actual_cents"] and mv["actual_cents"] == rs["real_cents"], "tarjeta Mes contra plan coincide con el presupuesto")
# ---- ahorro
A = D["ahorro"]; ok(all(sum(m["monto_cents"] * (1 if m["tipo"] == "aportación" else -1) for m in A["movimientos"] if m["destino"] == d) == v for d, v in [("fondo", A["fondo"]["ahorrado_cents"])] + [(m["id"], m["ahorrado_cents"]) for m in A["metas"]]) and A["ahorrado_total_cents"] == A["fondo"]["ahorrado_cents"] + sum(m["ahorrado_cents"] for m in A["metas"]), "ahorro: movimientos por destino = ahorrado; total = fondo + metas")
ok(all(m["avance_pct"] == rh(m["ahorrado_cents"] * 100, m["monto_meta_cents"]) for m in A["metas"]) and {m["estado"] for m in A["metas"]} == {"en_ruta", "atrasada"}, "metas: % de avance coincide y hay una en ruta y una atrasada")
ng = D["inicio"]["tarjetas"]["ahorro"]["next_goal"]; mg = [m for m in A["metas"] if m["id"] == ng["id"]][0]
ok({"name", "target_cents", "target_date", "saved_cents", "progress_pct", "planned_monthly_cents", "required_monthly_cents", "status"} <= set(ng) and ng["status"] == {"en_ruta": "on_track", "atrasada": "behind"}[mg["estado"]] and (ng["name"], ng["target_cents"], ng["target_date"], ng["saved_cents"], ng["progress_pct"]) == (mg["nombre"], mg["monto_meta_cents"], mg["fecha_meta"], mg["ahorrado_cents"], mg["avance_pct"]) and ng["target_date"] == min(m["fecha_meta"] for m in A["metas"]), "tarjeta Ahorro: next_goal usa los nombres de SavingsGoal (C1), 'behind' = atrasada y es la meta más próxima")
# ---- tickets (P-04, P-05)
T = D["tickets"]; g13 = [g for g in G if g["id"] == "g13"][0]; tk = T["confirmados"][0]
def imp_ok(r): return r["importe_cents"] == rh(round(r["cantidad"] * 1000) * r["precio_unitario_cents"], 1000)
ok(tk["gasto_id"] == "g13" and 8 <= len(tk["renglones"]) <= 12 and sum(r["importe_cents"] for r in tk["renglones"]) == tk["total_cents"] == g13["monto_cents"] and all(imp_ok(r) for r in tk["renglones"] if not r["es_ajuste"]) and sum(1 for r in tk["renglones"] if r["es_ajuste"]) == 1 and tk["renglones"][-1]["es_ajuste"], "ticket confirmado de g13: 8 a 12 renglones (el último es el ajuste), importes = cantidad × precio, la suma da el monto del gasto")
ok(all(r["hogar_cents"] + r["karen_cents"] == r["importe_cents"] for r in tk["renglones"]) and sum(r["hogar_cents"] for r in tk["renglones"]) == tk["hogar_cents"] == g13["hogar_cents"] and sum(r["karen_cents"] for r in tk["renglones"]) == tk["karen_cents"] == g13["karen_cents"] and sum(c["importe_cents"] for c in tk["por_cuenta"]) == tk["total_cents"] and all(r["cuenta"] for r in tk["renglones"]), "ticket confirmado: cada renglón tiene cuenta; hogar + Karen = importe; las partes suman las del gasto g13; totales por cuenta suman el total")
oc = T["lectura_pendiente"]; lows = [r for r in oc["renglones"] if r["confianza"] == "baja"]
ok(6 <= len(oc["renglones"]) <= 10 and 1 <= len(lows) <= 2 and {r["confianza"] for r in oc["renglones"]} >= {"alta", "media", "baja"} and all(r["texto_leido"] and r["producto_id"] and r["cuenta_sugerida"] for r in oc["renglones"]) and all(imp_ok(r) for r in oc["renglones"]) and oc["suma_renglones_cents"] == sum(r["importe_cents"] for r in oc["renglones"]) and oc["total_leido_cents"] == oc["suma_renglones_cents"] + oc["diferencia_cents"] != oc["suma_renglones_cents"] and oc["bajos"] == [r["orden"] for r in lows], "lectura OCR pendiente: 6 a 10 renglones con texto leído, producto y cuenta sugeridos y confianza (1 o 2 bajos); la diferencia contra el total cuadra")
PC = T["por_confirmar"]; bor = T["borradores"][0]
ok(len(PC) == 2 and sum(1 for x in PC if x["tiene_borrador"]) == 1 and [x["fecha_foto"] for x in PC] == sorted((x["fecha_foto"] for x in PC), reverse=True) and PC[0]["id"] == oc["id"] and bor["id"] in {x["id"] for x in PC if x["tiene_borrador"]} and PC[1]["total_cents"] == bor["borrador"]["total_cents"] == sum(r["importe_cents"] for r in bor["borrador"]["renglones"]) != bor["total_propuesto_cents"] and all(imp_ok(r) for r in bor["borrador"]["renglones"]), "tickets por confirmar: 2 lecturas por fecha descendente, una con borrador (su total es el del borrador y suma sus renglones)")
# ---- ajustes: notificaciones, datos semilla, parámetros
AJ = D["ajustes"]; NF = AJ["notificaciones"]["alertas"]; nb = {a["id"]: a for a in NF}
ok([a["id"] for a in NF] == ["payment_due", "budget_threshold", "stale_prices", "meal_plan_missing"] and all(set(a) >= {"activa", "destinatarios", "hora", "anticipacion"} and set(a["destinatarios"]) <= {"axel", "ximena"} and a["destinatarios"] for a in NF) and any(not a["activa"] for a in NF) and nb["payment_due"]["anticipacion"]["dias_antes"] == [1, 0] and nb["stale_prices"]["anticipacion"]["advance_days"] == 1 and nb["budget_threshold"]["umbrales"]["aviso_pct"] == AJ["parametros"]["umbral_aviso_pct"] == 80 and nb["budget_threshold"]["umbrales"]["tope_pct"] == AJ["parametros"]["umbral_limite_pct"] == 100 and nb["meal_plan_missing"]["anticipacion"]["weekday"] == 7 and all(a["hora"] == D["reportes"]["programados"]["avisos_hora"] == "09:00" for a in NF), "Ajustes › Notificaciones: 4 alertas (payment_due, budget_threshold, stale_prices, meal_plan_missing) con activa, destinatarios, anticipación y hora 09:00; pago [1, 0], precios viejos 1 día, umbrales 80 y 100, plan los domingos; una apagada")
ok(len(AJ["datos_semilla"]) == 3 and all(x["estado"] for x in AJ["datos_semilla"]) and not re.search(r"\d", " ".join(x["nombre"] for x in AJ["datos_semilla"])) and AJ["parametros"]["comparador_max_s"] == 5, "Ajustes › Datos semilla: 3 renglones con estado y sin cifras; parámetro comparador_max_s = 5 s")
def existe(path):
    o = D
    for k in path.split("."):
        if k.endswith("[]"):  # lista de objetos: cada elemento debe traer lo que sigue
            k = k[:-2]
            if not (isinstance(o, dict) and isinstance(o.get(k), list) and o[k]): return False
            rest = path.split(".")[path.split(".").index(k + "[]") + 1:]
            return all(existe_en(x, rest) for x in o[k])
        if not (isinstance(o, dict) and k in o): return False
        o = o[k]
    return True
def existe_en(o, ks):
    for k in ks:
        if not (isinstance(o, dict) and k in o): return False
        o = o[k]
    return True
ok(all(existe(l) for v in D["indice_pantallas"].values() for l in v) and "P-04/P-05" in D["indice_pantallas"] and "tickets" in D["indice_pantallas"]["P-04/P-05"], "índice de pantallas: todas las llaves que cita existen (ninguna colgante); P-04 y P-05 apuntan a tickets")
# ---- ronda 2: textos de Reportes, series de Tendencias, uso de recetas y registros de Datos semilla
import re as _re
def monto_fmt(c): return ("−$" if c < 0 else "$") + f"{abs(c) // 100:,}.{abs(c) % 100:02d}"
MES = {"ene": 1, "feb": 2, "mar": 3, "abr": 4, "may": 5, "jun": 6, "jul": 7, "ago": 8, "sep": 9, "oct": 10, "nov": 11, "dic": 12}
res = [a for a in D["reportes"]["avisos_enviados"] if a["tipo"] == "resumen diario"]; bien = len(res) >= 2
for a in res:
    m = _re.match(r"Resumen del (\d{1,2})-([a-z]{3})-(\d{4}): gastaste (\$[\d,]+\.\d{2}) \(parte del hogar\) en (\d+) gastos?\.$", a["texto"])
    if not m: bien = False; print("     resumen diario con formato inesperado:", a["texto"]); continue
    f = f"{m.group(3)}-{MES[m.group(2)]:02d}-{int(m.group(1)):02d}"; dia = [g for g in R if g["fecha"] == f]
    esperado = (monto_fmt(sum(g["hogar_cents"] for g in dia)), len(dia)); leido = (m.group(4), int(m.group(5)))
    if esperado != leido: bien = False; print("     resumen del", f, "dice", leido, "y gastos da", esperado)
ok(bien, "avisos enviados › resumen diario: el monto (parte del hogar) y el conteo de gastos de cada resumen coinciden con los gastos realizados de ese día (P3-12: Reportes cuadra con Gastos)")
PT = D["presupuesto"]["tendencia"]; IT = D["ingresos"]["tendencia"]
ok([x["mes"] for x in IT] == [x["mes"] for x in PT] and all(x["total_cents"] == x["quincena_15_cents"] + x["quincena_30_cents"] + x["extras_cents"] for x in IT) and IT[-1]["total_cents"] == D["ingresos"]["total_mes_cents"] and (IT[-1]["quincena_15_cents"], IT[-1]["quincena_30_cents"]) == tuple(q["monto_cents"] for q in D["ingresos"]["quincenas"]) and IT[-1]["incluye_esperado"] and not any(x["incluye_esperado"] for x in IT[:-1]), "ingresos.tendencia (P-08 N-103): mismos meses que presupuesto.tendencia; total = quincenas + extras; octubre = ingresos del mes")
TP = D["gastos_resumen"]["tendencia_por_categoria"]; ids = [c["id"] for c in C]; bien = [t["mes"] for t in TP["meses"]] == [x["mes"] for x in PT]
for t, pt in zip(TP["meses"], PT):
    pc = t["por_categoria"]
    bien &= [c["categoria_id"] for c in pc] == ids and sum(c["plan_cents"] for c in pc) == t["plan_cents"] == pt["plan_cents"] and sum(c["real_cents"] for c in pc) == t["real_cents"] == pt["real_cents"] and all(c["usado_pct"] == rh(c["real_cents"] * 100, c["plan_cents"]) for c in pc if c["plan_cents"])
    bien &= t["otras_real_cents"] == sum(c["real_cents"] for c in pc if c["categoria_id"] not in TP["seis_mayores"]) and t["otras_plan_cents"] == sum(c["plan_cents"] for c in pc if c["categoria_id"] not in TP["seis_mayores"])
bien &= [(c["categoria_id"], c["plan_cents"], c["real_cents"]) for c in TP["meses"][-1]["por_categoria"]] == [(c["id"], c["plan_cents"], c["real_cents"]) for c in C] and len(TP["seis_mayores"]) == 6 and set(TP["seis_mayores"]) <= set(ids)
ok(bien, "gastos_resumen.tendencia_por_categoria (P-02 N-044, P-06 N-085): 6 meses; cada mes suma el plan y el real de presupuesto.tendencia; % ROUND_HALF_UP; octubre = categorias; «Otras» = las 7 que no son de las 6 mayores")
PR = D["productos"]; pp = [e for p in PR for e in p["precios"]]; bien = True; npar = 0
for p in PR:
    HP = {(h["tienda_id"], h["presentacion_id"]): h["puntos"] for h in p["historial_precios"]}; bien &= len(HP) == len(p["historial_precios"]) == len(p["precios"])
    for e in p["precios"]:
        pts = HP.get((e["tienda_id"], e["presentacion_id"])); npar += 1
        if not pts: bien = False; continue
        fs_ = [x["fecha"] for x in pts]
        bien &= len(pts) >= 3 and fs_ == sorted(set(fs_)) and fs_[0] >= "2026-05-01" and (pts[-1]["fecha"], pts[-1]["precio_cents"], pts[-1]["fuente"]) == (e["fecha"], e["precio_cents"], e["fuente"]) and all(isinstance(x["precio_cents"], int) and x["precio_cents"] > 0 for x in pts)
ok(bien and npar == len(pp), f"productos[].historial_precios (P-15 N-171): {npar} pares producto × presentación × tienda con 3 a 5 fechas ascendentes, enteros en centavos; el último punto es el precio vigente de «precios»")
sem = {(d["fecha"], tn): {pe["receta_id"] for pe in d["tiempos"][tn].values() if pe.get("receta_id")} for d in D["plan_semana"]["dias"] for tn in d["tiempos"]}; bien = True
for r in D["recetas"]:
    u = r["uso"]; usados = [f for (f, tn), rs_ in sem.items() if r["id"] in rs_]; fut = sorted(f for f in usados if f >= D["meta"]["hoy"])
    bien &= u["en_semana_actual_hasta_hoy"] == sum(1 for f in usados if f <= D["meta"]["hoy"]) and u["veces_8_semanas"] == u["en_semana_actual_hasta_hoy"] + u["en_7_semanas_previas"] and u["proxima_fecha"] == (fut[0] if fut else None) and (date.fromisoformat(u["ventana_fin"]) - date.fromisoformat(u["ventana_inicio"])).days == 55
ok(bien, "recetas[].uso (P-12 N-141): veces en 8 semanas = las del plan de esta semana hasta hoy + 7 semanas previas; la próxima fecha sale del plan semanal")
ds = AJ["datos_semilla"]; ok(all(isinstance(x["registros"], int) and x["registros"] > 0 and x["registros_texto"].startswith(str(x["registros"])) for x in ds) and ds[0]["registros"] == len(D["productos"]) and ds[1]["registros"] == len(K["kueski"]["compras"]) + len(K["kueski"]["compras_programadas"]), "Ajustes › Datos semilla (P-26.11 N-305): cada renglón trae su número de registros, tomado de conteos de este mismo archivo")
# ---- calendario del hogar
CAL = D["calendario_hogar"]; fe = [e["fecha"] for e in CAL]
fij_oct = {(g["fecha"], g["monto_cents"]) for g in R if g["tipo"] == "fijo"}
ok(fe == sorted(fe) and fe[0] == "2026-10-01" and len(fij_oct) == 5 and fij_oct <= {(e["fecha"], e["monto_cents"]) for e in CAL if e["tipo"] == "gasto_fijo"} and [f for f in ("2026-10-01", "2026-10-06", "2026-10-12", "2026-10-14") if f not in fe] == [], "calendario del hogar: empieza el 2026-10-01 e incluye los gastos fijos ya realizados de octubre (renta, telefonía, internet, streaming)")
ok({e["fecha"] for e in CAL if e["tipo"] == "visita_mercado" and e.get("estado") == "realizada"} == {g["fecha"] for g in R if g["tienda_id"] in ("mercado-portales", "tianguis-narvarte")}, "calendario del hogar: las visitas a mercado realizadas coinciden con los gastos con ticket de mercado y tianguis")
# ---- comidas
P = D["plan_semana"]; cel = [c for d in P["dias"] for t in d["tiempos"].values() for c in t.values()]
ok(len(P["dias"]) == 7 and len(cel) == 7 * 4 * 2 and any(c["estado"] == "fuera" for c in cel), "plan semanal: 7 días × 4 tiempos × 2 personas con casillas «fuera»")
RB = {r["id"]: r for r in D["recetas"]}; ok(sum(c["costo_cents"] for c in cel if c["estado"] == "receta") == P["totales"]["costo_semana_cents"] and all(c["costo_cents"] == RB[c["receta_id"]]["costo_porcion_cents"] * c["porciones"] for c in cel if c["estado"] == "receta"), "plan semanal: el costo de las casillas suma el total y sale del costo por porción de la receta")
ok(len(RB) >= 8 and all(r["costo_total_cents"] == sum(i["costo_cents"] for i in r["ingredientes"]) and r["kcal_total"] == sum(i["kcal"] for i in r["ingredientes"]) and r["costo_porcion_cents"] == rh(r["costo_total_cents"], 2) for r in RB.values()), f"{len(RB)} recetas (≥ 8): costo y kcal totales = Σ ingredientes; por porción = total ÷ 2")
L = D["lista_compras"]; ok(L["total_cents"] == sum(g["total_cents"] for g in L["grupos"]) and all(g["total_cents"] == sum(i["subtotal_cents"] for i in g["items"]) and all(i["subtotal_cents"] == i["piezas"] * i["precio_cents"] for i in g["items"]) for g in L["grupos"]), "lista de compras: subtotales = piezas × precio; grupos por tienda suman el total")
CN = D["canasta"]; ok(all(c["total_cents"] == sum(r["precios"][c["tienda_id"]].get("costo_cents", 0) for r in CN["detalle"]) for c in CN["por_tienda"]) and sum(CN["ganadores"].values()) + CN["empates"] == CN["articulos"], "canasta: total por tienda = Σ del detalle; ganadores + empates = artículos")
ok(all(c["identicos"] + c["sustitutos"] + c["sin_precio"] == CN["articulos"] for c in CN["por_tienda"]), "canasta: idénticos + sustitutos + sin precio = artículos en cada tienda")
# ---- productos y precios
PR = D["productos"]; pp = [e for p in PR for e in p["precios"]]
ok(len(PR) >= 30 and sum(1 for e in pp if e["viejo"]) >= 5 and all(e["viejo"] == ((HOY - date.fromisoformat(e["fecha"])).days >= 30) and e["dias"] == (HOY - date.fromisoformat(e["fecha"])).days for e in pp), f"{len(PR)} productos (≥ 30); precios viejos = 30 días o más contra hoy; días coherentes")
ok(len(D["tiendas"]) == 5 and any(t["tipo"] == "mercado" for t in D["tiendas"]) and any(t["tipo"] == "tianguis" for t in D["tiendas"]) and all(t["visita"] for t in D["tiendas"] if t["tipo"] in ("mercado", "tianguis")), "tiendas: Walmart, Bodega Aurrera, Sam's Club y 2 mercados con día de visita")
# ---- salud
S = D["salud"]["perfiles"]; ok(all(len(p["historial_peso"]) == 8 and all(h["imc"] == round(h["peso_kg"] / (p["altura_cm"] / 100.0) ** 2, 1) for h in p["historial_peso"]) for p in S), "salud: 8 semanas de peso por persona; IMC recalculado coincide")
# ---- proyección y deuda
PQ = D["proyeccion"]; ok(len(PQ["meses"]) == 12 and len(PQ["flujo_por_quincena"]) == 24, "proyección: 12 meses y 24 quincenas")
acc = PQ["saldo_inicial_cents"]; good = True
for r in PQ["flujo_por_quincena"]:
    s = r["ingresos_cents"] - r["pagos_credito_cents"] - r["gasto_contado_o_fijo_cents"] - r["variable_estimado_cents"] - r["ahorro_planificado_cents"]; acc += s; good &= (s == r["saldo_cents"] and acc == r["saldo_acumulado_cents"])
ok(good, "flujo por quincena: saldo = ingresos − egresos; el acumulado encadena desde el saldo inicial")
dif = 0; good = True
for m in PQ["meses"]: dif += m["devengado_total_cents"] - m["efectivo_total_cents"]; good &= (m["devengado_total_cents"] == m["devengado_hogar_cents"] + m["devengado_karen_cents"] and m["diferencia_acumulada_cents"] == dif)
ok(good, "devengado contra efectivo: total = hogar + Karen; la diferencia acumulada encadena mes a mes")
ok(all(d["total_cents"] == sum(d["por_linea_cents"].values()) == d["hogar_cents"] + d["karen_cents"] for d in D["deuda_viva"]["serie_mensual"]) and D["deuda_viva"]["serie_mensual"][-1]["total_cents"] == 0, "deuda viva: total = Σ líneas = hogar + Karen; llega a 0 al cierre del horizonte")
# ---- reportes y avisos
Rp = D["reportes"]["programados"]; ok((Rp["resumen_diario"]["hora"], Rp["resumen_semanal"]["hora"], Rp["resumen_semanal"]["dia"], Rp["mensual_pdf"]["dia_del_mes"], Rp["mensual_pdf"]["hora"], Rp["avisos_hora"]) == ("21:30", "19:00", "domingo", 1, "08:00", "09:00"), "reportes: diario 21:30, semanal domingo 19:00, mensual día 1 08:00; avisos 09:00")
# ---- fechas contra model.json
fm = set(re.findall(r"\d{4}-\d{2}-\d{2}", open(RAIZ + "/oraculo/model.json", encoding="utf-8").read())) | set(re.findall(r"\d{4}-\d{2}-\d{2}", open(RAIZ + "/datos/semilla/data.json", encoding="utf-8").read()))
fs = set(re.findall(r"\d{4}-\d{2}-\d{2}", TXT)); co = sorted(fm & fs)
def regla(f): y, m, d = map(int, f.split("-")); return d == 15 or d == min(30, calendar.monthrange(y, m)[1])
print("     coincidencias con model.json/data.json:", co)
ok(all(regla(f) for f in co), "las fechas que coinciden con el calendario real son solo de la regla de quincena (día 15, o día 30 / último día del mes), documentadas como inevitables")
ok(not [f for f in fs if f[8:] in ("09", "10") and f > "2026-10-01"] and "2026-10-10" not in fs, "ninguna fecha cae en un día 9 o 10 después del 1-oct-2026 ni en el 10-oct-2026")
print("RESULTADO: OK" if not fallos else "RESULTADO: FALLA (" + "; ".join(fallos) + ")"); sys.exit(1 if fallos else 0)
