#!/usr/bin/env python3
"""Genera docs/design/maquetas/datos-sinteticos.json (Pista 3, K45.1). Todo es INVENTADO.
Uso: python3 -I generar_sinteticos.py <salida.json>. Determinista: sin azar ni reloj."""
import sys, json, math, hashlib, calendar
from datetime import date, timedelta
HOY = date(2026, 10, 21)
def rh(n, d):  # ROUND_HALF_UP entero (positivos)
    return (2 * n + d) // (2 * d)
def iso(d): return d.isoformat()
def D(s): return date.fromisoformat(s)
def fm(c):
    s = f"${abs(c)//100:,}.{abs(c)%100:02d}"
    return ("−" + s) if c < 0 else s
def qd(y, m): return [date(y, m, 15), date(y, m, min(30, calendar.monthrange(y, m)[1]))]
def qlist(start, n):  # primeras n fechas de quincena >= start
    out, y, m = [], start.year, start.month
    while len(out) < n:
        out += [d for d in qd(y, m) if d >= start]; y, m = (y + (m == 12), m % 12 + 1)
    return out[:n]
def addm(y, m, k): t = y * 12 + m - 1 + k; return t // 12, t % 12 + 1
def split(total, pct):  # parte del hogar (redondeo half up) y parte de Karen (resto)
    h = rh(total * pct, 100); return h, total - h
def h8(s): return int(hashlib.sha256(s.encode()).hexdigest()[:8], 16)
def amort(P, rate_m, n):  # sistema francés; el último pago ajusta centavos
    pm = int(round(P * rate_m / (1 - (1 + rate_m) ** -n)))
    out, s = [], P
    for i in range(1, n + 1):
        it = int(round(s * rate_m)); pay = pm if i < n else s + it
        out.append(dict(num=i, pago_cents=pay, capital_cents=pay - it, interes_cents=it, saldo_cents=s - (pay - it))); s -= pay - it
    return out
PER = [dict(id="axel", nombre="Axel", rol="Titular", tiene_cuenta=True, usuario="axel", token_color="color.person.axel"),
       dict(id="ximena", nombre="Ximena", rol="Co-titular", tiene_cuenta=True, usuario="ximena", token_color="color.person.ximena"),
       dict(id="karen", nombre="Karen", rol="Participante externa", tiene_cuenta=False, usuario=None, token_color="color.person.karen")]
# ---------------- ingresos
ING = dict(esperado=dict(q15_cents=1426000, q30_cents=1478000, rige_desde="2026-10-01"),
  quincenas=[dict(fecha="2026-10-15", estado="capturada", monto_cents=1426000, capturada_por="axel", capturada_el="2026-10-15T18:42:00-06:00"),
             dict(fecha="2026-10-30", estado="esperada", monto_cents=1478000, capturada_por=None, capturada_el=None)],
  extras=[dict(fecha="2026-10-04", concepto="Venta de bicicleta usada", monto_cents=135000),
          dict(fecha="2026-10-18", concepto="Reembolso de seguro dental", monto_cents=74000)])
ING_MES = sum(q["monto_cents"] for q in ING["quincenas"]) + sum(e["monto_cents"] for e in ING["extras"])
# ---------------- categorías: id, nombre, tipo, reparto, plan
CATS = [("vivienda","Vivienda","fijo","hogar",720000),("servicios","Servicios y suscripciones","fijo","hogar",260000),
 ("despensa-comun","Despensa común","variable","conjunto50",200000),("super","Súper de la semana","variable","hogar",320000),
 ("mercado","Frutas y verduras (mercado)","variable","hogar",180000),("transporte","Transporte","variable","hogar",260000),
 ("salud","Salud","variable","hogar",130000),("restaurantes","Restaurantes y café","variable","hogar",180000),
 ("ocio","Ocio","variable","hogar",90000),("mantenimiento","Hogar y mantenimiento","variable","hogar",110000),
 ("personal","Ropa y personal","variable","hogar",100000),("mascotas","Mascotas","variable","hogar",100000),
 ("seguros","Seguros","fijo","hogar",108600)]
# ---------------- gastos: (fecha, comercio, cat, monto, persona, reparto(h=100,k=0,c<pct>), pago, tipo, ticket, tienda)
G = [("2026-10-01","Renta departamento","vivienda",720000,"axel","h","contado","fijo",0,None),
("2026-10-02","Walmart · súper semanal","super",114235,"ximena","h","contado","variable",1,"walmart"),
("2026-10-03","Gasolina Pemex","transporte",78000,"axel","h","contado","variable",1,None),
("2026-10-03","Mercado Portales · frutas y verduras","mercado",41250,"ximena","h","contado","variable",1,"mercado-portales"),
("2026-10-05","Farmacia del Ahorro","salud",28640,"ximena","h","contado","variable",1,None),
("2026-10-06","Telefonía móvil Axel","servicios",39900,"axel","h","contado","fijo",0,None),
("2026-10-06","Telefonía móvil Ximena","servicios",34900,"ximena","h","contado","fijo",0,None),
("2026-10-06","Tianguis de la Narvarte · verdura","mercado",20500,"ximena","h","contado","variable",1,"tianguis-narvarte"),
("2026-10-07","Petco · croquetas y arena","mascotas",51860,"ximena","h","contado","variable",1,None),
("2026-10-08","Café Punta del Cielo","restaurantes",9200,"axel","h","contado","variable",0,None),
("2026-10-11","Taquería · comida del domingo","restaurantes",54000,"ximena","h","contado","variable",0,None),
("2026-10-12","Internet del hogar","servicios",59900,"axel","h","contado","fijo",0,None),
("2026-10-12","Walmart · despensa bimestral","despensa-comun",291840,"axel","c50","kueski","variable",1,"walmart"),
("2026-10-13","Tianguis de la Narvarte · verdura","mercado",33050,"ximena","h","contado","variable",1,"tianguis-narvarte"),
("2026-10-14","Streaming de video y música","servicios",21900,"axel","h","contado","fijo",0,None),
("2026-10-14","Bodega Aurrera · súper semanal","super",96820,"ximena","h","contado","variable",1,"bodega-aurrera"),
("2026-10-16","Gasolina Pemex","transporte",82000,"axel","h","tarjeta","variable",1,None),
("2026-10-17","Mercado Portales · frutas y verduras","mercado",43620,"ximena","h","contado","variable",1,"mercado-portales"),
("2026-10-17","Cine y palomitas","ocio",46500,"axel","h","tarjeta","variable",0,None),
("2026-10-18","Sam's Club · comida familiar","despensa-comun",124000,"axel","c60","contado","variable",1,"sams-club"),
("2026-10-18","Ferretería · focos y cinta","mantenimiento",38990,"axel","h","contado","variable",1,None),
("2026-10-19","Mercado Libre · audífonos de Karen","personal",138000,"axel","k","kueski","variable",0,None),
("2026-10-19","Uber (dos viajes)","transporte",18600,"axel","h","contado","variable",0,None),
("2026-10-19","Ropa · playera de Ximena","personal",45900,"ximena","h","tarjeta","variable",1,None),
("2026-10-20","Tianguis de la Narvarte · verdura","mercado",30130,"ximena","h","contado","variable",1,"tianguis-narvarte"),
("2026-10-20","Consulta médica de Axel","salud",65000,"axel","h","contado","variable",0,None),
("2026-10-21","Veterinario · vacuna","mascotas",58000,"ximena","h","tarjeta","variable",1,None),
("2026-10-21","Café y pan","restaurantes",13450,"ximena","h","contado","variable",0,None),
("2026-10-21","Bodega Aurrera · lácteos y reposición","super",41135,"ximena","h","contado","variable",1,"bodega-aurrera"),
("2026-10-23","Agua (recibo)","servicios",23840,"axel","h","contado","venidero",0,None),
("2026-10-24","Luz (recibo)","servicios",68430,"axel","h","contado","venidero",0,None),
("2026-10-26","Consulta con nutriólogo de Ximena","salud",85000,"ximena","h","contado","venidero",0,None),
("2026-10-27","Seguro del auto","seguros",108600,"axel","h","contado","venidero",0,None)]
# ---------------- créditos: compras a Kueski (capital, partes, plazo, tasa por plazo)
TASA = {2: 164150, 4: 243120}  # millonésimas de fracción: 16.4150 % y 24.3120 %
KUE = [("kc0","2026-09-27","Soriana · despensa chica",164000,60,2,"2026-10-15"),
       ("kc1","2026-10-12","Walmart · despensa bimestral",291840,50,4,"2026-10-30"),
       ("kc2","2026-10-19","Mercado Libre · audífonos de Karen",138000,0,2,"2026-11-15"),
       ("kc3","2026-12-12","Walmart · despensa bimestral (programada)",300000,50,4,"2026-12-30"),
       ("kc4","2027-02-12","Walmart · despensa bimestral (programada)",300000,50,4,"2027-02-28")]
def reparte(tot, n):  # n pagos, el último ajusta centavos
    p = rh(tot, n); return [p] * (n - 1) + [tot - p * (n - 1)]
def mk_compra(cid, f, com, cap, hpct, n, p1):
    inte = rh(cap * TASA[n], 10**6); tot = cap + inte
    caph, capk = split(cap, hpct); inth, intk = split(inte, hpct) if cap else (0, 0)
    # el interés se reparte en proporción al capital de cada parte
    inth = rh(inte * caph, cap); intk = inte - inth
    ph, pk = reparte(caph + inth, n), reparte(capk + intk, n)
    ch, ck = reparte(caph, n), reparte(capk, n)
    fechas = qlist(D(p1), n)
    pagos = []
    for i in range(n):
        pagos.append(dict(num=i + 1, fecha=iso(fechas[i]), monto_cents=ph[i] + pk[i], hogar_cents=ph[i], karen_cents=pk[i],
            capital_repuesto_cents=ch[i] + ck[i], hogar_estado="pendiente", karen_estado="pendiente"))
    return dict(id=cid, linea_id="kueski", comercio=com, fecha=f, capital_cents=cap, intereses_cents=inte, total_cents=tot, plazo_quincenas=n,
        tasa_pct=f"{TASA[n]/10000:.4f} %", promocion_id=None, hogar=dict(capital_cents=caph, intereses_cents=inth, total_cents=caph + inth),
        karen=dict(capital_cents=capk, intereses_cents=intk, total_cents=capk + intk), pagos=pagos, programada=cid in ("kc3", "kc4"))
KC = [mk_compra(*k) for k in KUE]
for p in KC[0]["pagos"][:1]:  # 15-oct de la compra anterior: ambas partes ya marcadas
    p.update(hogar_estado="pagada", karen_estado="pagada", hogar_pagada_el="2026-10-16", karen_pagada_el="2026-10-16")
# ---------------- otras líneas de crédito
PRE = amort(1800000, 0.299 / 12, 15)
PRE_F = [date(2026 + (6 + i) // 12, (6 + i) % 12 + 1, 18) for i in range(15)]  # 18-jul-2026 a 18-sep-2027
TIE = amort(699000, 0.045, 10)
TIE_F = [date(2026 + (8 + i) // 12, (8 + i) % 12 + 1, 20) for i in range(10)]  # 20-sep-2026 a 20-jun-2027
for i, p in enumerate(PRE): p["fecha"] = iso(PRE_F[i]); p["estado"] = "pagada" if PRE_F[i] <= HOY else "pendiente"
for i, p in enumerate(TIE): p["fecha"] = iso(TIE_F[i]); p["estado"] = "pagada" if TIE_F[i] <= HOY else "pendiente"
MSI = 117492  # mensualidad; 12 en total
PAGOS = []  # pagos programados de TODAS las líneas (lista plana, hogar/Karen)
def addp(linea, compra, f, h, k, cap, hest, kest, cfecha, num=None, hp=None, kp=None):
    PAGOS.append(dict(id=f"{compra}-{len(PAGOS)+1}", linea_id=linea, compra_id=compra, fecha=f, monto_cents=h + k, hogar_cents=h, karen_cents=k,
        capital_repuesto_cents=cap, hogar_estado=hest, karen_estado=kest, compra_fecha=cfecha, num=num))
for c in KC:
    for p in c["pagos"]:
        addp("kueski", c["id"], p["fecha"], p["hogar_cents"], p["karen_cents"], p["capital_repuesto_cents"], p["hogar_estado"], p["karen_estado"], c["fecha"], p["num"])
for p in PRE: addp("prestamo", "pre1", p["fecha"], p["pago_cents"], 0, p["capital_cents"], p["estado"], "n/a", "2026-06-18", p["num"])
for p in TIE: addp("tienda", "tie1", p["fecha"], p["pago_cents"], 0, p["capital_cents"], p["estado"], "n/a", "2026-08-22", p["num"])
addp("tarjeta", "tc-estado", "2026-10-06", 521340, 0, 521340, "pagada", "n/a", "2026-09-16")
addp("tarjeta", "tc-estado", "2026-11-06", 610485, 0, 610485, "pendiente", "n/a", "2026-10-17")
addp("tarjeta", "tc-estado", "2026-12-06", 103900 + MSI, 0, 103900 + MSI, "pendiente", "n/a", "2026-10-19")
for i in range(6): addp("tarjeta", "tc-msi", iso(date(2027 + (i) // 12, i + 1, 6)), MSI, 0, MSI, "pendiente", "n/a", "2026-06-18")
for p in PAGOS: p["estado_global"] = "pagada" if p["hogar_estado"] == "pagada" and p["karen_estado"] in ("pagada", "n/a") else "pendiente"
PAGOS.sort(key=lambda p: (p["fecha"], p["linea_id"]))
def hoy_estado(p):
    if p["estado_global"] == "pagada": return "pagada"
    return "sin_marcar" if D(p["fecha"]) < HOY else "pendiente"
# ---------------- gastos derivados
GAS = []
for i, g in enumerate(G):
    f, com, cat, m, per, rep, pago, tipo, tk, tienda = g
    if rep == "h": h, k, rp = m, 0, dict(modo="hogar")
    elif rep == "k": h, k, rp = 0, m, dict(modo="karen")
    else: pc = int(rep[1:]); h, k = split(m, pc); rp = dict(modo="conjunto", hogar_pct=pc, karen_pct=100 - pc)
    d = dict(id=f"g{i+1:02d}", fecha=f, comercio=com, categoria_id=cat, monto_cents=m, persona_id=per, reparto=rp, hogar_cents=h, karen_cents=k,
        forma_pago=pago, tipo=tipo, tiene_ticket=bool(tk), tienda_id=tienda, estado="venidero" if tipo == "venidero" else "realizado", compra_credito_id=None)
    d["devengado"] = [dict(mes=f[:7], hogar_cents=h, karen_cents=k)]
    d["efectivo"] = [dict(fecha=f, hogar_cents=h, karen_cents=k)] if pago == "contado" else []
    if pago == "kueski":
        c = KC[1] if "despensa" in com else KC[2]; d["compra_credito_id"] = c["id"]
        ht, kt = c["hogar"]["total_cents"], c["karen"]["total_cents"]
        if c["id"] == "kc1": d["devengado"] = [dict(mes="2026-10", hogar_cents=rh(ht, 2), karen_cents=rh(kt, 2)), dict(mes="2026-11", hogar_cents=ht - rh(ht, 2), karen_cents=kt - rh(kt, 2))]
        else: d["devengado"] = [dict(mes="2026-10", hogar_cents=ht, karen_cents=kt)]
        d["efectivo"] = [dict(fecha=p["fecha"], hogar_cents=p["hogar_cents"], karen_cents=p["karen_cents"]) for p in c["pagos"]]
    if pago == "tarjeta": d["efectivo"] = [dict(fecha="2026-11-06" if f <= "2026-10-17" else "2026-12-06", hogar_cents=h, karen_cents=k)]
    GAS.append(d)
REAL = [g for g in GAS if g["estado"] == "realizado"]
def hog_mes(g, mes="2026-10"): return sum(x["hogar_cents"] for x in g["devengado"] if x["mes"] == mes)
def est(p): return "excedido" if p >= 100 else "aviso" if p >= 80 else "normal"
CAT = []
for (cid, nom, tipo, rep, plan) in CATS:
    gs = [g for g in REAL if g["categoria_id"] == cid]; real = sum(hog_mes(g) for g in gs); p = rh(real * 100, plan)
    # fechas de cruce de 80 y 100 %
    acc, c80, c100 = 0, None, None
    for g in sorted(gs, key=lambda g: g["fecha"]):
        acc += hog_mes(g)
        if c80 is None and acc * 100 >= 80 * plan: c80 = g["fecha"]
        if c100 is None and acc * 100 >= 100 * plan: c100 = g["fecha"]
    CAT.append(dict(id=cid, nombre=nom, tipo=tipo, reparto_defecto=dict(modo="conjunto", hogar_pct=50, karen_pct=50) if rep == "conjunto50" else dict(modo="hogar"),
        orden=len(CAT) + 1, activa=True, plan_cents=plan, real_cents=real, diferencia_cents=plan - real, usado_pct=p, estado=est(p), num_gastos=len(gs),
        cruzo_80=c80 if p >= 80 else None, cruzo_100=c100 if p >= 100 else None, venidero_cents=sum(g["hogar_cents"] for g in GAS if g["categoria_id"] == cid and g["estado"] == "venidero")))
PLAN = sum(c["plan_cents"] for c in CAT); REALT = sum(c["real_cents"] for c in CAT)
AHO_PLAN = 200000
def qi(f): return (f[:7], 1 if int(f[8:]) <= 15 else 2)
# ---------------- tiendas
TIENDAS = [dict(id="walmart", nombre="Walmart", tipo="cadena", token_color="color.store.walmart", visita=None, requiere_membresia=False, tiene_membresia=False, ultimo_ticket="2026-10-02"),
 dict(id="bodega-aurrera", nombre="Bodega Aurrera", tipo="cadena", token_color="color.store.bodega-aurrera", visita=None, requiere_membresia=False, tiene_membresia=False, ultimo_ticket="2026-10-21"),
 dict(id="sams-club", nombre="Sam's Club", tipo="cadena", token_color="color.store.sams-club", visita=None, requiere_membresia=True, tiene_membresia=False, membresia_costo_anual_cents=65000, evaluacion="en evaluación (la casa no tiene membresía)", ultimo_ticket="2026-10-18"),
 dict(id="mercado-portales", nombre="Mercado Portales", tipo="mercado", token_color="color.store.market", visita=dict(dia_semana=6, nombre="sábado"), compra_ahi=["verduras", "frutas", "lacteos-huevo"], ultimo_ticket="2026-10-17"),
 dict(id="tianguis-narvarte", nombre="Tianguis de la Narvarte", tipo="tianguis", token_color="color.store.market", visita=dict(dia_semana=2, nombre="martes"), compra_ahi=["verduras", "frutas", "carnes"], ultimo_ticket="2026-10-20")]
FRESCAS = {"walmart": ["2026-10-02"], "bodega-aurrera": ["2026-10-14", "2026-10-21"], "sams-club": ["2026-10-18"], "mercado-portales": ["2026-10-03", "2026-10-17"], "tianguis-narvarte": ["2026-10-06", "2026-10-13", "2026-10-20"]}
VIEJAS = ["2026-08-22", "2026-09-12", "2026-09-18"]
FAC = {"walmart": 1.0, "bodega-aurrera": 0.985, "sams-club": 0.88, "mercado-portales": 0.93, "tianguis-narvarte": 0.92}
# id, nombre, marca, grupo, unidad, kcal (por 100 g/ml; por pieza si unidad pza), [(pres, etiqueta, cantidad_base, tiendas W B S M)], precio de la presentación 0 en Walmart (cents)
PROD = [("jitomate","Jitomate saladet","","verduras","g",18,[("kg","1 kg",1000,"WBM")],3400),("cebolla","Cebolla blanca","","verduras","g",40,[("kg","1 kg",1000,"WBM"),("b3","bolsa 3 kg",3000,"S")],2850),
("aguacate","Aguacate Hass","","frutas","g",160,[("kg","1 kg",1000,"WBM")],7900),("limon","Limón con semilla","","frutas","g",30,[("kg","1 kg",1000,"WBM")],3200),
("papa","Papa blanca","","verduras","g",77,[("kg","1 kg",1000,"WBM"),("b5","bolsa 5 kg",5000,"S")],2700),("zanahoria","Zanahoria","","verduras","g",41,[("kg","1 kg",1000,"WBM")],2100),
("platano","Plátano tabasco","","frutas","g",89,[("kg","1 kg",1000,"WBM")],2450),("manzana","Manzana roja","","frutas","g",52,[("kg","1 kg",1000,"WBM"),("b25","bolsa 2.5 kg",2500,"S")],4800),
("espinaca","Espinaca","","verduras","g",23,[("mj","manojo 250 g",250,"WBM")],1800),("calabacita","Calabacita","","verduras","g",17,[("kg","1 kg",1000,"WBM")],3100),
("tomate-verde","Tomate verde","","verduras","g",32,[("kg","1 kg",1000,"WBM")],2900),("serrano","Chile serrano","","verduras","g",32,[("kg","1 kg",1000,"WBM")],5500),
("pollo","Pechuga de pollo","","carnes","g",165,[("kg","1 kg",1000,"WB"),("c2","charola 2 kg",2000,"S")],10900),("res-molida","Carne molida de res","","carnes","g",250,[("kg","1 kg",1000,"WB"),("c2","paquete 2 kg",2000,"S")],14200),
("huevo","Huevo blanco","San Juan","lacteos-huevo","pza",72,[("d12","12 piezas",12,"WBM"),("c30","30 piezas",30,"S")],4300),("atun","Atún en agua","Dolores","abarrotes","g",116,[("l140","lata 140 g",140,"WB"),("p6","6 latas de 140 g",840,"S")],1950),
("frijol","Frijol negro","Verde Valle","abarrotes","g",340,[("b900","bolsa 900 g",900,"WB")],3600),("leche","Leche entera","Lala","lacteos-huevo","ml",61,[("l1","1 L",1000,"WB"),("c12","caja de 12 L",12000,"S")],2650),
("panela","Queso panela","Fud","lacteos-huevo","g",280,[("p400","400 g",400,"WB")],6200),("yogur","Yogur natural","Yoplait","lacteos-huevo","g",61,[("g900","900 g",900,"WB"),("g27","3 × 900 g",2700,"S")],4100),
("arroz","Arroz blanco","Morelos","abarrotes","g",360,[("b900","bolsa 900 g",900,"WB"),("b4","bolsa 4 kg",4000,"S")],2700),("pasta","Pasta espagueti","La Moderna","abarrotes","g",371,[("p200","200 g",200,"WB")],950),
("aceite","Aceite vegetal","1-2-3","abarrotes","ml",884,[("l1","1 L",1000,"WB"),("l3","3 L",3000,"S")],4800),("tortilla","Tortilla de maíz","","abarrotes","g",218,[("kg","1 kg",1000,"WBM")],2300),
("pan","Pan de caja integral","Bimbo","abarrotes","g",250,[("p680","680 g",680,"WB")],4400),("avena","Avena en hojuelas","Quaker","abarrotes","g",389,[("b800","bolsa 800 g",800,"WB")],3800),
("cafe","Café soluble","Nescafé","abarrotes","g",353,[("f200","frasco 200 g",200,"WBS")],10900),("azucar","Azúcar estándar","Zulka","abarrotes","g",387,[("kg","1 kg",1000,"WB"),("b5","bolsa 5 kg",5000,"S")],3150),
("papel","Papel higiénico","Regio","limpieza","pza",None,[("r12","12 rollos",12,"WB"),("r30","30 rollos",30,"S")],7800),("detergente","Detergente líquido","Ariel","limpieza","ml",None,[("l3","3 L",3000,"WB"),("l6","6 L",6000,"S")],11900),
("jabon-trastes","Jabón para trastes","Axión","limpieza","ml",None,[("m750","750 ml",750,"WB")],3400),("cloro","Cloro","Cloralex","limpieza","ml",None,[("l1","1 L",1000,"WB")],1550),
("croquetas","Croquetas para perro","Pedigree","mascotas","g",None,[("k4","4 kg",4000,"WB"),("k15","15 kg",15000,"S")],16800)]
SID = {"W": "walmart", "B": "bodega-aurrera", "S": "sams-club"}
PRODS, PREC = [], {}  # PREC[(prod, pres, tienda)] = dict
for (pid, nom, marca, grupo, un, kc, pres, p0) in PROD:
    pu0 = p0 / pres[0][2]; lst = []
    for (qid, et, qb, tt) in pres:
        tiendas = [SID[c] for c in tt if c in SID]
        if "M" in tt:
            mk = [t for t in ("mercado-portales", "tianguis-narvarte") if h8(pid + t) % 4]; tiendas += mk or ["mercado-portales"]
        for t in tiendas:
            h = h8(pid + qid + t); ruido = 0.96 + (h % 9) / 100.0
            precio = int(round(pu0 * qb * FAC[t] * ruido / 10.0)) * 10
            r = h % 100
            if pid == "pasta": r = 5
            if r < 14: fecha, fuente = VIEJAS[h % 3], "manual"
            elif r < 22 and t in SID.values(): fecha, fuente = "2026-10-02", "profeco"
            elif r < 55: fecha, fuente = "2026-10-19", "manual"
            else: fecha, fuente = FRESCAS[t][h % len(FRESCAS[t])], "ticket"
            dias = (HOY - D(fecha)).days
            e = dict(tienda_id=t, presentacion_id=qid, precio_cents=precio, precio_unidad_base_cents=round(precio / qb, 4), fecha=fecha, fuente=fuente, dias=dias, viejo=dias >= 30, solo_referencia=fuente == "profeco")
            PREC[(pid, qid, t)] = e; lst.append(e)
    if pid != "pasta" and not [e for e in lst if not e["viejo"] and not e["solo_referencia"]]:  # todo producto con al menos un precio fresco (salvo la pasta, ejemplo de precio viejo)
        e = lst[0]; e.update(fecha=FRESCAS[e["tienda_id"]][0], fuente="ticket", viejo=False, solo_referencia=False, dias=(HOY - D(FRESCAS[e["tienda_id"]][0])).days)
    fr = [e for e in lst if not e["viejo"] and not e["solo_referencia"]]
    mejor = min(fr, key=lambda e: e["precio_unidad_base_cents"]) if fr else None
    q = {x[0]: x[2] for x in pres}
    PRODS.append(dict(id=pid, nombre=nom, marca=marca or None, grupo=grupo, unidad_base=un, kcal_por_100_g_ml=kc if un != "pza" else None, kcal_por_pieza=kc if un == "pza" else None,
        nutricion_fuente=None if kc is None else "USDA Foundation Foods 04/2026 (ejemplo)", presentaciones=[dict(id=x[0], etiqueta=x[1], cantidad_base=x[2]) for x in pres],
        precios=lst, frescura="fresco" if fr else ("viejo" if lst else "sin_precio"),
        mejor_precio_fresco=dict(tienda_id=mejor["tienda_id"], presentacion_id=mejor["presentacion_id"], precio_unidad_base_cents=mejor["precio_unidad_base_cents"], fecha=mejor["fecha"]) if mejor else None))
PBY = {p["id"]: p for p in PRODS}
def ref(pid):  # precio de referencia por unidad base: el más barato fresco entre tiendas
    m = PBY[pid]["mejor_precio_fresco"] or min(PBY[pid]["precios"], key=lambda e: e["precio_unidad_base_cents"])  # sin fresco: el más barato aunque sea viejo
    return (m["precio_unidad_base_cents"], m["tienda_id"], m["fecha"])
# ---------------- recetas
REC = [("r01","Avena con plátano y canela",["desayuno","colacion"],5,10,[("avena",80),("leche",400),("platano",240),("azucar",10)]),
("r02","Huevos a la mexicana con tortillas",["desayuno"],10,10,[("huevo",6),("jitomate",200),("cebolla",60),("serrano",20),("aceite",15),("tortilla",200)]),
("r03","Pollo asado con arroz y ensalada",["comida","cena"],15,35,[("pollo",500),("arroz",160),("jitomate",150),("limon",40),("zanahoria",100),("aceite",15)]),
("r04","Pasta con atún y jitomate",["comida","cena"],10,20,[("pasta",200),("atun",280),("jitomate",300),("cebolla",60),("aceite",15)]),
("r05","Frijoles de la olla con aguacate",["comida"],10,90,[("frijol",250),("cebolla",60),("tortilla",250),("aguacate",150),("serrano",10)]),
("r06","Sopa de verduras",["cena","comida"],15,30,[("papa",300),("zanahoria",200),("calabacita",250),("cebolla",60),("jitomate",150),("aceite",10)]),
("r07","Picadillo de res con papa",["comida"],15,35,[("res-molida",400),("papa",300),("zanahoria",150),("jitomate",200),("cebolla",60),("aceite",15)]),
("r08","Quesadillas de panela con salsa verde",["cena"],10,15,[("tortilla",300),("panela",200),("tomate-verde",200),("serrano",20),("cebolla",40)]),
("r09","Yogur con manzana y avena",["desayuno","colacion"],5,0,[("yogur",500),("manzana",300),("avena",60)]),
("r10","Ensalada de atún con pan",["cena","colacion"],10,0,[("atun",280),("aguacate",150),("jitomate",200),("limon",60),("pan",100)]),
("r11","Tacos de pollo con guacamole",["cena"],15,20,[("pollo",400),("tortilla",300),("aguacate",200),("limon",40),("cebolla",40)]),
("r12","Espinacas con huevo y panela",["desayuno"],10,10,[("espinaca",250),("huevo",4),("panela",100),("tortilla",200)])]
def kc_unit(pid): p = PBY[pid]; return p["kcal_por_pieza"] if p["unidad_base"] == "pza" else p["kcal_por_100_g_ml"] / 100.0
RECETAS, RBY = [], {}
for (rid, nom, tiempos, prep, coc, ing) in REC:
    L, tc, tk = [], 0, 0
    for (pid, q) in ing:
        pu, t, f = ref(pid); c = int(round(pu * q)); k = int(round(kc_unit(pid) * q)); tc += c; tk += k
        pr = PBY[pid]["presentaciones"][0]
        L.append(dict(producto_id=pid, producto=PBY[pid]["nombre"], cantidad=q, unidad=PBY[pid]["unidad_base"], presentacion_sugerida=pr["etiqueta"],
            piezas_a_comprar=math.ceil(q / pr["cantidad_base"]), costo_cents=c, precio_tienda_id=t, precio_fecha=f, kcal=k, kcal_fuente="USDA (ejemplo)"))
    r = dict(id=rid, nombre=nom, tiempos_aptos=tiempos, porciones_base=2, preparacion_min=prep, coccion_min=coc, tiempo_total_min=prep + coc, origen="recetario inicial" if rid < "r09" else "Axel y Ximena",
        ingredientes=L, costo_total_cents=tc, costo_porcion_cents=rh(tc, 2), kcal_total=tk, kcal_porcion=rh(tk, 2),
        pasos=["Lava y corta los ingredientes.", "Cocina según el tiempo indicado.", "Sirve y ajusta la sal al gusto."])
    RECETAS.append(r); RBY[rid] = r
# ---------------- plan semanal (semana lun 19 a dom 25 oct 2026, ISO 2026-W43)
SEM0 = date(2026, 10, 19); DIAS = ["Lun", "Mar", "Mié", "Jue", "Vie", "Sáb", "Dom"]
DES = ["r01","r02","r09","r12","r01","r02","r09"]; COM = ["r03","r05","r07","r04","r05","r03","F"]; CEN = ["r08","r10","r11","r06","r08","r11",None]; COL = ["r09",None,"r01",None,"r10",None,None]
PLANSEM, SEMCOST, KCAL = [], {"axel": [0]*7, "ximena": [0]*7}, {"axel": [0]*7, "ximena": [0]*7}
def celda(x): return None if x is None else ("fuera" if x == "F" else x)
for i in range(7):
    row = dict(fecha=iso(SEM0 + timedelta(i)), dia=DIAS[i], tiempos={})
    for tn, arr in (("desayuno", DES), ("comida", COM), ("cena", CEN), ("colacion", COL)):
        row["tiempos"][tn] = {}
        for per in ("axel", "ximena"):
            v = celda(arr[i])
            if per == "axel" and tn == "comida" and i in (0, 2): v = "fuera"
            if per == "ximena" and tn == "cena" and i == 4: v = "fuera"
            if per == "ximena" and tn == "cena" and i == 5: v = None
            if v in (None, "fuera"): c = dict(estado="vacio" if v is None else "fuera")
            else:
                r = RBY[v]; por = 2 if per == "axel" and tn in ("comida", "cena") else 1; c = dict(estado="receta", receta_id=v, porciones=por, costo_cents=r["costo_porcion_cents"] * por, kcal=r["kcal_porcion"] * por)
                SEMCOST[per][i] += c["costo_cents"]; KCAL[per][i] += c["kcal"]
            row["tiempos"][tn][per] = c
    PLANSEM.append(row)
SIN_ARMAR = sum(1 for r in PLANSEM for tn in ("desayuno", "comida", "cena") if any(r["tiempos"][tn][p]["estado"] == "vacio" for p in ("axel", "ximena")))
COSTO_SEM = sum(sum(v) for v in SEMCOST.values())
# ---------------- inventario
INV = [("arroz","unidad",650,"g",400,None),("aceite","unidad",700,"ml",300,None),("azucar","unidad",900,"g",300,None),("frijol","unidad",0,"g",300,None),("pasta","unidad",400,"g",400,None),
 ("atun","unidad",280,"g",280,None),("cafe","unidad",120,"g",100,None),("leche","ritmo",1500,"ml",None,3500),("huevo","ritmo",6,"pza",None,18),("tortilla","ritmo",700,"g",None,1400),
 ("papel","ritmo",5,"pza",None,2),("croquetas","ritmo",2600,"g",None,900),("jabon-trastes","ritmo",400,"ml",None,150),("detergente","ritmo",1200,"ml",None,500)]
INVL, STK, RIT = [], {}, {}
for (pid, modo, ex, un, um, cons) in INV:
    dias = None if modo == "unidad" else round(ex / (cons / 7.0), 1)
    estado = "agotado" if ex == 0 else ("reponer" if (modo == "unidad" and ex <= um) or (modo == "ritmo" and dias < 7) else "suficiente")
    if modo == "unidad": STK[pid] = ex
    else: RIT[pid] = (cons, ex, dias)
    h = h8("inv" + pid)
    INVL.append(dict(producto_id=pid, producto=PBY[pid]["nombre"], modo="stock por unidad" if modo == "unidad" else "ritmo de consumo", existencia=ex, unidad=un, umbral_reposicion=um,
        consumo_semana=cons, dias_restantes=dias, estado=estado, ultima_entrada=dict(fecha=["2026-10-02", "2026-10-14", "2026-10-17", "2026-10-19"][h % 4], origen=["compra", "sobrante", "ajuste"][h % 3 if h % 7 else 0])))
ORD = {"agotado": 0, "reponer": 1, "suficiente": 2}
INVL.sort(key=lambda x: (ORD[x["estado"]], x["producto"]))
INV_MOV = [dict(producto_id="leche", fecha="2026-10-14", tipo="entrada por compra", cantidad=6000, unidad="ml", origen="Bodega Aurrera"), dict(producto_id="leche", fecha="2026-10-19", tipo="consumo del plan", cantidad=-1500, unidad="ml", origen="plan semanal"),
 dict(producto_id="arroz", fecha="2026-10-17", tipo="sobrante de receta", cantidad=150, unidad="g", origen="Pollo asado con arroz"), dict(producto_id="frijol", fecha="2026-10-20", tipo="ajuste", cantidad=-120, unidad="g", origen="conteo")]
# ---------------- lista de compras de la semana y canasta comparada
NEC, MOT = {}, {}
for row in PLANSEM:
    for tn in row["tiempos"].values():
        for per in tn.values():
            if per["estado"] == "receta":
                for ing in RBY[per["receta_id"]]["ingredientes"]:
                    NEC[ing["producto_id"]] = NEC.get(ing["producto_id"], 0) + ing["cantidad"] * per["porciones"] / 2.0; MOT.setdefault(ing["producto_id"], set()).add(per["receta_id"])
NETOS = {}
for pid, n in NEC.items():
    net = n - STK.get(pid, 0)
    if pid in RIT and RIT[pid][2] < 7: net = max(net, RIT[pid][0] - RIT[pid][1])
    if net > 0: NETOS[pid] = (net, "recetas" if pid not in RIT else "recetas y reposición", sorted(MOT[pid]))
for pid, (cons, ex, dias) in RIT.items():
    if pid not in NETOS and dias < 7: NETOS[pid] = (cons - ex, "reposición", [])
NETOS["detergente"] = (2000.0, "agregado a mano", [])
def cands(pid, net, solo_fresco=True):
    out = []
    for p in PBY[pid]["presentaciones"]:
        for t in [e["tienda_id"] for e in PBY[pid]["precios"] if e["presentacion_id"] == p["id"]]:
            e = PREC[(pid, p["id"], t)]
            if e["solo_referencia"] or (solo_fresco and e["viejo"]): continue
            pz = math.ceil(net / p["cantidad_base"]); out.append((pz * e["precio_cents"], t, p, e, pz))
    return out
ITEMS = []
for pid, (net, mot, recs) in sorted(NETOS.items()):
    c = cands(pid, net) or cands(pid, net, False); ok = [x for x in c if x[2]["cantidad_base"] <= 2.2 * net] or c
    cost, t, p, e, pz = min(ok, key=lambda x: (x[0], x[2]["cantidad_base"]))
    ITEMS.append(dict(producto_id=pid, producto=PBY[pid]["nombre"], cantidad_necesaria=round(net), unidad=PBY[pid]["unidad_base"], presentacion_id=p["id"], presentacion=p["etiqueta"], piezas=pz, tienda_id=t,
        precio_cents=e["precio_cents"], subtotal_cents=cost, precio_fecha=e["fecha"], precio_viejo=e["viejo"], motivo=mot, recetas=recs, comprado=False))
for i, it in enumerate([x for x in ITEMS if x["tienda_id"] == "walmart"][:2] + [x for x in ITEMS if x["tienda_id"] == "tianguis-narvarte"][:1]): it["comprado"] = True
LISTA = []
for t in TIENDAS:
    its = [x for x in ITEMS if x["tienda_id"] == t["id"]]
    if its: LISTA.append(dict(tienda_id=t["id"], tienda=t["nombre"], dia_visita=t.get("visita"), articulos=len(its), total_cents=sum(x["subtotal_cents"] for x in its), items=its))
CAN = {t["id"]: dict(tienda_id=t["id"], total_cents=0, identicos_cents=0, identicos=0, sustitutos=0, sin_precio=0, viejos=0, cubiertos=0) for t in TIENDAS}
DET = []
for it in ITEMS:
    pid, net = it["producto_id"], it["cantidad_necesaria"]; row = dict(producto_id=pid, producto=it["producto"], cantidad=net, unidad=it["unidad"], presentacion=it["presentacion"], piezas=it["piezas"], precios={})
    for t in TIENDAS:
        tid = t["id"]; e = PREC.get((pid, it["presentacion_id"], tid)); c = CAN[tid]
        if e and not e["solo_referencia"]: cost, tipo, ee = it["piezas"] * e["precio_cents"], "idéntico", e
        else:
            pq = {p["id"]: p["cantidad_base"] for p in PBY[pid]["presentaciones"]}
            alt = [x for x in PBY[pid]["precios"] if x["tienda_id"] == tid and not x["solo_referencia"] and pq[x["presentacion_id"]] <= 2.2 * net]  # un sustituto solo si su tamaño es comparable
            if alt:
                ee = min(alt, key=lambda x: math.ceil(net / pq[x["presentacion_id"]]) * x["precio_cents"]); cost, tipo = math.ceil(net / pq[ee["presentacion_id"]]) * ee["precio_cents"], "sustituto"
            else: c["sin_precio"] += 1; row["precios"][tid] = dict(tipo="sin precio"); continue
        c["total_cents"] += cost; c["cubiertos"] += 1; c["viejos"] += ee["viejo"]
        if tipo == "idéntico": c["identicos"] += 1; c["identicos_cents"] += cost
        else: c["sustitutos"] += 1
        row["precios"][tid] = dict(tipo=tipo, costo_cents=cost, fecha=ee["fecha"], fuente=ee["fuente"], viejo=ee["viejo"], presentacion_id=ee["presentacion_id"])
    cs = {k: v["costo_cents"] for k, v in row["precios"].items() if "costo_cents" in v}; mn = min(cs.values())
    ganan = [k for k, v in cs.items() if v == mn]; row["ganadora"] = ganan[0] if len(ganan) == 1 else None; row["empate"] = ganan if len(ganan) > 1 else []
    row["diferencia_cents"] = max(cs.values()) - mn; DET.append(row)
for c in CAN.values(): c["cobertura_pct"] = rh(c["cubiertos"] * 100, len(ITEMS))
COMP = sorted([c for c in CAN.values() if c["cobertura_pct"] >= 90], key=lambda c: c["total_cents"])
GAN = {t["id"]: sum(1 for r in DET if r["ganadora"] == t["id"]) for t in TIENDAS}; EMP = sum(1 for r in DET if r["empate"])
SAMS_ITEMS = [r for r in DET if "costo_cents" in r["precios"].get("sams-club", {})]
AHO_CAN = sum(min(v["costo_cents"] for k, v in r["precios"].items() if k != "sams-club" and "costo_cents" in v) - r["precios"]["sams-club"]["costo_cents"] for r in SAMS_ITEMS if any("costo_cents" in v for k, v in r["precios"].items() if k != "sams-club"))
# ---------------- ahorro
def movs(dest, ini_fecha, total, regulares, extra=()):
    L = [(f, "aportación", m, "") for f, m in regulares] + list(extra); res = total - sum(m if t == "aportación" else -m for (_, t, m, _) in L)
    return [(ini_fecha, "aportación", res, "Aportación inicial")] + L
MV = []
for dest, rows in (("fondo", movs("fondo", "2026-05-04", 2148000, [(f, 80000) for f in ("2026-06-04", "2026-07-04", "2026-08-04", "2026-09-04", "2026-10-02")], [("2026-08-21", "retiro", 120000, "Reparación del coche")])),
                   ("laptop", movs("laptop", "2026-05-18", 615000, [(f, 50000) for f in ("2026-06-04", "2026-07-04", "2026-08-04", "2026-09-04")] + [("2026-10-16", 25000)])),
                   ("viaje", movs("viaje", "2026-05-04", 1410000, [(f, 70000) for f in ("2026-06-04", "2026-07-04", "2026-08-04", "2026-09-04", "2026-10-02")]))):
    for (f, t, m, n) in rows: MV.append(dict(fecha=f, tipo=t, monto_cents=m, destino=dest, nota=n or None))
MV.sort(key=lambda x: (x["fecha"], x["destino"]), reverse=True)
for i, m in enumerate(MV): m["id"] = f"mv{len(MV)-i:02d}"
def ahorrado(d=None, hasta="9999"): return sum((x["monto_cents"] if x["tipo"] == "aportación" else -x["monto_cents"]) for x in MV if (d is None or x["destino"] == d) and x["fecha"] <= hasta)
APORT_OCT = sum(x["monto_cents"] for x in MV if x["tipo"] == "aportación" and x["fecha"].startswith("2026-10"))
METAS = []
for (mid, nom, mt, fm_, plan) in (("laptop", "Laptop nueva", 1450000, "2027-01-17", 50000), ("viaje", "Viaje a Oaxaca", 1800000, "2027-04-17", 70000)):
    sv = ahorrado(mid); mm = (int(fm_[:4]) - 2026) * 12 + int(fm_[5:7]) - 10; nec = -(-(mt - sv) // mm)
    METAS.append(dict(id=mid, nombre=nom, monto_meta_cents=mt, fecha_meta=fm_, ahorrado_cents=sv, avance_pct=rh(sv * 100, mt), aportacion_planificada_cents=plan, aportacion_necesaria_cents=nec, estado="en_ruta" if plan >= nec else "atrasada", etiqueta_estado="En ruta" if plan >= nec else "Atrasada"))
METAS.sort(key=lambda m: m["fecha_meta"])
AHO = dict(ahorrado_total_cents=ahorrado(), fondo=dict(ahorrado_cents=ahorrado("fondo"), aportacion_planificada_cents=80000), metas=METAS, aportado_mes_cents=APORT_OCT, planificado_mes_cents=AHO_PLAN, movimientos=MV,
    tendencia=[dict(mes=m, ahorrado_cierre_cents=ahorrado(None, m + "-31")) for m in ("2026-05", "2026-06", "2026-07", "2026-08", "2026-09", "2026-10")])
# ---------------- salud
def sal(pid, nom, h, edad, sexo, niv, fac, obj, aj, kg, conf, hist):
    bmr = 10 * kg + 6.25 * h - 5 * edad + (5 if sexo == "H" else -161); sug = int(bmr * fac - (aj if obj == "bajar" else 0) + 0.5)
    imc = lambda k: round(k / (h / 100.0) ** 2, 1); cl = lambda v: "bajo peso" if v < 18.5 else "normal" if v < 25 else "sobrepeso" if v < 30 else "obesidad"
    return dict(persona_id=pid, nombre=nom, altura_cm=h, edad=edad, sexo=sexo, nivel_actividad=niv, factor_actividad=fac, objetivo=obj, ajuste_objetivo_kcal=aj, peso_actual_kg=kg, peso_fecha=hist[-1][0], imc=imc(kg), imc_clasificacion=cl(imc(kg)),
        meta_kcal_sugerida=sug, meta_kcal_confirmada=conf, requiere_confirmar=conf != sug, historial_peso=[dict(fecha=f, peso_kg=w, imc=imc(w)) for f, w in hist])
SEM8 = [iso(date(2026, 8, 31) + timedelta(7 * i)) for i in range(8)]
SAL = dict(perfiles=[sal("axel", "Axel", 178, 36, "H", 3, 1.55, "bajar", 400, 84.6, 2364, list(zip(SEM8, [86.9, 86.5, 86.1, 85.8, 85.2, 85.0, 84.9, 84.6]))),
                     sal("ximena", "Ximena", 163, 33, "M", 2, 1.375, "mantener", 0, 62.8, 1900, list(zip(SEM8, [63.4, 63.2, 63.3, 63.0, 62.9, 63.1, 62.9, 62.8])))],
          casa=dict(ajuste_bajar_kcal=0, ajuste_subir_kcal=0), kcal_hoy=[dict(persona_id=p, planeadas=KCAL[p][2], dia="Mié 21 oct 2026") for p in ("axel", "ximena")])
for p, k in zip(SAL["perfiles"], SAL["kcal_hoy"]): k.update(meta=p["meta_kcal_confirmada"], diferencia=k["planeadas"] - p["meta_kcal_confirmada"])
SEM12 = [iso(date(2026, 7, 27) + timedelta(7 * i)) for i in range(12)]
SAL["tendencia_semanal"] = [dict(semana_inicio=f, costo_cents=70000 + h8("c" + f) % 2600 * 10, kcal_axel=2150 + h8("a" + f) % 380, kcal_ximena=1700 + h8("x" + f) % 330) for f in SEM12]
# ---------------- flujo por quincena (efectivo, parte del hogar)
FIJ = [("Renta", "vivienda", 720000, 1), ("Telefonía móvil Axel", "servicios", 39900, 6), ("Telefonía móvil Ximena", "servicios", 34900, 6), ("Internet del hogar", "servicios", 59900, 12), ("Streaming", "servicios", 21900, 14),
       ("Agua (recibo)", "servicios", 23840, 23), ("Luz (recibo)", "servicios", 68430, 24), ("Seguro del auto", "seguros", 108600, 27)]
FIJOS = [dict(concepto=c, categoria_id=k, monto_cents=m, dia_del_mes=d, reparto=dict(modo="hogar"), activo=True) for (c, k, m, d) in FIJ]
VAR_EST, AHO_Q = 580000, 100000
def pag_q(y, m, q, solo=None):
    lo, hi = (1, 15) if q == 1 else (16, 31)
    return [p for p in PAGOS if p["fecha"][:7] == f"{y}-{m:02d}" and lo <= int(p["fecha"][8:]) <= hi]
def con_q(y, m, q, desde=1, hasta=31, venid=None):
    lo, hi = max(desde, 1 if q == 1 else 16), min(hasta, 15 if q == 1 else 31)
    return sum(g["hogar_cents"] for g in GAS if g["forma_pago"] == "contado" and g["fecha"][:7] == f"{y}-{m:02d}" and lo <= int(g["fecha"][8:]) <= hi and (venid is None or (g["estado"] == "venidero") == venid))
SALDO_INI = 1200000  # efectivo disponible al arrancar octubre (inventado)
FQ = []
for k in range(24):
    y, m = addm(2026, 10, k // 2); q = k % 2 + 1; idx = k // 2; pc = sum(p["hogar_cents"] for p in pag_q(y, m, q)); fx = sum(f["monto_cents"] for f in FIJOS if (f["dia_del_mes"] <= 15) == (q == 1))
    if k == 0: ing = ING["quincenas"][0]["monto_cents"] + ING["extras"][0]["monto_cents"]
    elif k == 1: ing = ING["quincenas"][1]["monto_cents"] + ING["extras"][1]["monto_cents"]
    else: ing = ING["esperado"]["q15_cents"] if q == 1 else ING["esperado"]["q30_cents"]
    var = VAR_EST if idx < 6 else 0; ah = AHO_Q if idx < 6 else 0
    if k < 2: var, fx = 0, con_q(y, m, q)  # octubre: lo real y lo comprometido ya está en «contado» (incluye fijos y venideros)
    FQ.append(dict(quincena=f"{y}-{m:02d}-q{q}", inicio=iso(date(y, m, 1 if q == 1 else 16)), fin=iso(date(y, m, 15 if q == 1 else calendar.monthrange(y, m)[1])), mes_n=idx + 1, ingresos_cents=ing, ingreso_esperado=k >= 1,
        pagos_credito_cents=pc, gasto_contado_o_fijo_cents=fx, variable_estimado_cents=var, ahorro_planificado_cents=ah, base="real/comprometido" if k < 2 else ("fijos + variable estimado" if idx < 6 else "solo compromisos contratados")))
acc = SALDO_INI
for r in FQ:
    r["saldo_cents"] = r["ingresos_cents"] - r["pagos_credito_cents"] - r["gasto_contado_o_fijo_cents"] - r["variable_estimado_cents"] - r["ahorro_planificado_cents"]; acc += r["saldo_cents"]; r["saldo_acumulado_cents"] = acc
# tarjeta de Inicio · flujo de la quincena en curso (16 a 31 oct)
HECHOS = con_q(2026, 10, 2, 16, 21) + sum(p["hogar_cents"] for p in pag_q(2026, 10, 2) if p["hogar_estado"] == "pagada" and int(p["fecha"][8:]) <= 21)
COMPR = con_q(2026, 10, 2, 22, 31) + sum(p["hogar_cents"] for p in pag_q(2026, 10, 2) if p["hogar_estado"] != "pagada")
ING_Q2 = FQ[1]["ingresos_cents"]; CIERRE = ING_Q2 - HECHOS - COMPR
FQ[1].update(pagos_credito_cents=sum(p["hogar_cents"] for p in pag_q(2026, 10, 2)), gasto_contado_o_fijo_cents=con_q(2026, 10, 2)); FQ[1]["saldo_cents"] = ING_Q2 - FQ[1]["pagos_credito_cents"] - FQ[1]["gasto_contado_o_fijo_cents"] - FQ[1]["ahorro_planificado_cents"]
acc = SALDO_INI
for r in FQ: acc += r["saldo_cents"]; r["saldo_acumulado_cents"] = acc
# ---------------- proyección 12 meses (devengado contra efectivo)
DESP = {"2026-10": [("kc1", 0)], "2026-11": [("kc1", 1)], "2026-12": [("kc3", 0)], "2027-01": [("kc3", 1)], "2027-02": [("kc4", 0)], "2027-03": [("kc4", 1)]}
KCB = {c["id"]: c for c in KC}; PROY = []; ad = ae = 0
for i in range(12):
    y, m = addm(2026, 10, i); mes = f"{y}-{m:02d}"; fixm = sum(f["monto_cents"] for f in FIJOS); dh = dk = 0
    for cid, half in DESP.get(mes, []):
        c = KCB[cid]; ht, kt = c["hogar"]["total_cents"], c["karen"]["total_cents"]; dh += rh(ht, 2) if half == 0 else ht - rh(ht, 2); dk += rh(kt, 2) if half == 0 else kt - rh(kt, 2)
    if i == 0: dev_h = REALT + sum(g["hogar_cents"] for g in GAS if g["estado"] == "venidero"); dev_k = dk
    else: dev_h = fixm + (2 * VAR_EST if i < 6 else 0) + dh; dev_k = dk
    pm = [p for p in PAGOS if p["fecha"][:7] == mes]
    pre = [p for p in pm if p["compra_id"] in ("pre1", "tie1", "tc-msi")]  # compras anteriores al sistema: se devengan al pagarse
    dev_h += sum(p["hogar_cents"] for p in pre)
    ef_h = (sum(g["hogar_cents"] for g in GAS if g["forma_pago"] == "contado" and g["fecha"][:7] == mes) if i == 0 else fixm + (2 * VAR_EST if i < 6 else 0)) + sum(p["hogar_cents"] for p in pm)
    ef_k = sum(p["karen_cents"] for p in pm) + (sum(g["karen_cents"] for g in GAS if g["forma_pago"] == "contado" and g["fecha"][:7] == mes) if i == 0 else 0)
    ad += dev_h + dev_k; ae += ef_h + ef_k
    PROY.append(dict(mes=mes, mes_n=i + 1, devengado_hogar_cents=dev_h, devengado_karen_cents=dev_k, devengado_total_cents=dev_h + dev_k, efectivo_hogar_cents=ef_h, efectivo_karen_cents=ef_k, efectivo_total_cents=ef_h + ef_k,
        diferencia_acumulada_cents=ad - ae, base="octubre: real + venideros" if i == 0 else ("fijos + variable estimado + compromisos" if i < 6 else "solo compromisos contratados (interruptor apagado)")))
# ---------------- deuda viva
def deuda(f):
    ps = [p for p in PAGOS if p["fecha"] > f and p["compra_fecha"] <= f]; por = {}
    for p in ps: por[p["linea_id"]] = por.get(p["linea_id"], 0) + p["monto_cents"]
    return dict(fecha=f, total_cents=sum(por.values()), hogar_cents=sum(p["hogar_cents"] for p in ps), karen_cents=sum(p["karen_cents"] for p in ps), por_linea_cents=por)
DEU = []
for i in range(12):
    y, m = addm(2026, 10, i); f = iso(date(y, m, calendar.monthrange(y, m)[1])); d = deuda(f)
    pm = [p for p in PAGOS if p["fecha"][:7] == f[:7]]; d.update(mes=f[:7], intereses_mes_cents=sum(p["monto_cents"] - p["capital_repuesto_cents"] for p in pm)); DEU.append(d)
# ---------------- disponible A/B de la línea revolvente
LIM, ANCLA, COLCHON = 1200000, "2026-10-12", 415000
def serie(compras):
    ev = []; b0 = LIM - sum(c["capital_cents"] - sum(p["capital_repuesto_cents"] for p in c["pagos"] if p["fecha"] <= ANCLA) for c in compras if c["fecha"] <= ANCLA)
    a, b = LIM, b0; out = [dict(fecha=ANCLA, tipo="start", descripcion="Arranque de la lectura A", monto_cents=0, capital_repuesto_cents=0, interes_no_repuesto_cents=0, disponible_a_cents=a, disponible_b_cents=b)]
    for c in compras:
        if c["fecha"] > ANCLA: ev.append((c["fecha"], 0, "purchase", c["comercio"], c["capital_cents"], 0, 0, -c["capital_cents"]))
        for p in c["pagos"]:
            if p["fecha"] > ANCLA: ev.append((p["fecha"], 1, "payment", f"Pago {p['num']} · {c['comercio']}", p["monto_cents"], p["capital_repuesto_cents"], p["monto_cents"] - p["capital_repuesto_cents"], p["capital_repuesto_cents"]))
    for (f, _, t, d, mo, cr, inr, dl) in sorted(ev):
        a += dl; b += dl; out.append(dict(fecha=f, tipo=t, descripcion=d, monto_cents=mo, capital_repuesto_cents=cr, interes_no_repuesto_cents=inr, disponible_a_cents=a, disponible_b_cents=b))
    return out
SER = serie(KC); ULT = [e for e in SER if e["fecha"] <= iso(HOY)][-1]
# ---------------- promociones y comparador de escenarios
PROMOS = [dict(id="pr1", linea_id="kueski", comercio="Walmart", plazo_quincenas=2, tasa_promocional_pct="0 %", etiqueta_tasa="sin intereses", compra_minima_cents=150000, inicio="2026-11-02", fin="2026-11-20", fuente="Aviso en la app de la línea"),
 dict(id="pr2", linea_id="kueski", comercio="Bodega Aurrera", plazo_quincenas=4, tasa_promocional_pct="9.5000 %", etiqueta_tasa="9.5000 %", compra_minima_cents=200000, inicio="2026-10-19", fin="2026-11-18", fuente="Aviso en la app de la línea"),
 dict(id="pr3", linea_id="tarjeta", comercio="Liverpool", plazo_quincenas=None, plazo_meses=6, tasa_promocional_pct="0 %", etiqueta_tasa="6 meses sin intereses", compra_minima_cents=250000, inicio="2026-11-12", fin="2026-11-22", fuente="Folleto del banco")]
for p in PROMOS: p["vigente_hoy"] = p["inicio"] <= iso(HOY) <= p["fin"]
def primer_pago(f): return iso(qlist(D(f) + timedelta(15), 1)[0])
def escenario(eid, f, n, promo):
    cap = 420000; c = mk_compra("cx", f, "Walmart · licuadora y horno", cap, 100, n, primer_pago(f))
    if promo: c["intereses_cents"] = 0; c["total_cents"] = cap; pg = reparte(cap, n); [p.update(monto_cents=pg[i], hogar_cents=pg[i], karen_cents=0, capital_repuesto_cents=pg[i]) for i, p in enumerate(c["pagos"])]
    s = serie([x for x in KC] + [c]); mb = min(e["disponible_b_cents"] for e in s)
    vivas = sum(1 for x in KC + [c] if x["fecha"] <= f and x["pagos"][-1]["fecha"] >= f)
    ok_col = mb >= COLCHON; ok_viv = vivas <= 4
    return dict(id=eid, linea_id="kueski", fecha=f, plazo_quincenas=n, promocion_id="pr1" if promo else None, intereses_cents=c["intereses_cents"], total_cents=c["total_cents"],
        pago_max_cents=max(p["monto_cents"] for p in c["pagos"]), pago_min_cents=min(p["monto_cents"] for p in c["pagos"]), disponible_minimo_b_cents=mb, compras_vivas_max=vivas,
        deuda_cero_fecha=c["pagos"][-1]["fecha"], valido=ok_col and ok_viv, motivo_descarte=None if ok_col and ok_viv else ("Deja el disponible por debajo del colchón mínimo" if not ok_col else "Rebasa el máximo de compras vivas"))
ESC = [escenario("e1", "2026-11-05", 2, True), escenario("e2", "2026-11-05", 4, False), escenario("e3", "2026-11-23", 2, False), escenario("e4", "2026-12-12", 4, False), escenario("e5", "2026-11-05", 2, False)]
ESC.append(dict(id="e6", linea_id="tarjeta", fecha="2026-11-13", plazo_meses=6, promocion_id="pr3", intereses_cents=0, total_cents=420000, pago_max_cents=70000, pago_min_cents=70000, disponible_minimo_b_cents=None, compras_vivas_max=None, deuda_cero_fecha="2027-05-06", valido=True, motivo_descarte=None))
VAL = [e for e in ESC if e["valido"]]
OPT = dict(menos_intereses=min(VAL, key=lambda e: e["intereses_cents"])["id"], menor_pago_por_quincena=min(VAL, key=lambda e: e["pago_max_cents"])["id"],
           mayor_disponible_minimo=max([e for e in VAL if e["disponible_minimo_b_cents"] is not None], key=lambda e: e["disponible_minimo_b_cents"])["id"], deuda_en_cero_mas_pronto=min(VAL, key=lambda e: e["deuda_cero_fecha"])["id"])
# ---------------- líneas de crédito
def vivos(lid): return [p for p in PAGOS if p["linea_id"] == lid and p["estado_global"] != "pagada"]
LINEAS = [dict(id="kueski", nombre="Kueski Pay", tipo="revolvente", institucion="Kueski", limite_cents=LIM, frecuencia_pagos="quincenal (15 y 30 o último día)", primer_pago_min_dias=15, lectura_activa="B", lectura_confirmada=False, ancla_lectura_a=ANCLA, activa=True,
   comision_anual_cents=0, comision_tardio_cents=25000, comision_apertura_cents=0),
 dict(id="tarjeta", nombre="Banorte Oro", tipo="tarjeta", institucion="Banorte", limite_cents=3800000, tasa_anual_pct="42.5000 %", base_tasa="diaria", incluye_iva=True, pago_minimo_pct="6.00 %", pago_minimo_fijo_cents=30000, dia_corte=17, dia_pago=6,
   comision_anual_cents=89000, comision_tardio_cents=45000, comision_apertura_cents=0, activa=True),
 dict(id="prestamo", nombre="Préstamo personal", tipo="prestamo", institucion="Santander", limite_cents=1800000, tasa_anual_pct="29.9000 %", frecuencia_pagos="mensual (día 18)", activa=True),
 dict(id="tienda", nombre="Crédito Coppel", tipo="tienda", institucion="Coppel", limite_cents=699000, tasa_anual_pct="54.0000 %", base_tasa="mensual", incluye_iva=True, pago_minimo_pct="12.00 %", pago_minimo_fijo_cents=50000, frecuencia_pagos="mensual (día 20)", activa=True)]
for l in LINEAS:
    vs = vivos(l["id"]); l["deuda_viva_cents"] = sum(p["monto_cents"] for p in vs if p["compra_fecha"] <= iso(HOY)); sig = min(vs, key=lambda p: p["fecha"]) if vs else None
    l["proximo_pago"] = None if not sig else dict(fecha=sig["fecha"], monto_cents=sum(p["monto_cents"] for p in vs if p["fecha"] == sig["fecha"]))
    l["compras_vivas"] = len({p["compra_id"] for p in vs if p["compra_fecha"] <= iso(HOY)})
    l["disponible_cents"] = (ULT["disponible_b_cents"] if l["id"] == "kueski" else l["limite_cents"] - 714385 if l["id"] == "tarjeta" else None)
LINEAS[0]["disponible_a_hoy_cents"] = ULT["disponible_a_cents"]; LINEAS[0]["disponible_b_hoy_cents"] = ULT["disponible_b_cents"]
DEUDA_TOT = sum(l["deuda_viva_cents"] for l in LINEAS)
def pend(cond): return [p for p in PAGOS if cond(p)]
PEND = sorted([p for p in PAGOS if p["estado_global"] != "pagada" and p["compra_fecha"] <= iso(HOY)], key=lambda p: p["fecha"]); F0 = PEND[0]["fecha"]; DIA0 = [p for p in PEND if p["fecha"] == F0]
PROX = dict(due_date=F0, dia_semana=DIAS[D(F0).weekday()], credit_id=DIA0[0]["linea_id"], amount_cents=sum(p["monto_cents"] for p in DIA0), household_share_cents=sum(p["hogar_cents"] for p in DIA0), external_share_cents=sum(p["karen_cents"] for p in DIA0), days_left=(D(F0) - HOY).days, payment_count=len(DIA0), is_overdue=False)  # nombres de DashboardNextPayment (C1); dia_semana es de presentación
# ---------------- Karen
KPAG = [p for p in PAGOS if p["karen_cents"] > 0 and p["compra_fecha"] <= iso(HOY)]
KAR = dict(despensa_devengada_mes_cents=sum(x["karen_cents"] for g in GAS if g["categoria_id"] == "despensa-comun" for x in g["devengado"] if x["mes"] == "2026-10"),
  compras_particulares=[dict(compra_id="kc2", comercio=KC[2]["comercio"], fecha=KC[2]["fecha"], total_cents=KC[2]["karen"]["total_cents"], pagos_hechos=0, pagos_total=2, saldo_cents=KC[2]["karen"]["total_cents"])],
  pagos_del_mes=[dict(fecha=p["fecha"], compra_id=p["compra_id"], monto_cents=p["karen_cents"], estado="pagada" if p["karen_estado"] == "pagada" else "pendiente") for p in KPAG if p["fecha"].startswith("2026-10")],
  deuda_viva_cents=sum(p["karen_cents"] for p in KPAG if p["karen_estado"] != "pagada"),
  nota="Karen paga directo a Kueski su parte; su parte es informativa y repone la línea cuando paga")
EC = dict(periodo="2026-10", compras_particulares=KAR["compras_particulares"], partes_despensa=[dict(compra_id=c["id"], fecha=c["fecha"], capital_cents=c["karen"]["capital_cents"], intereses_cents=c["karen"]["intereses_cents"], total_cents=c["karen"]["total_cents"], pago_por_quincena_cents=c["pagos"][0]["karen_cents"]) for c in KC[:2]],
  pagado_periodo_cents=sum(p["karen_cents"] for p in KPAG if p["karen_estado"] == "pagada" and p["fecha"].startswith("2026-10")), saldo_pendiente_cents=KAR["deuda_viva_cents"], nota=KAR["nota"])
# ---------------- calendario del hogar (1-oct-2026 a 17-abr-2027: octubre completo con lo ya realizado, y lo que viene)
CAL = []
def ev(f, t, tit, monto=None, **k): CAL.append(dict(fecha=f, tipo=t, titulo=tit, monto_cents=monto, **k))
for (y, m) in ((2026, 10), (2026, 11), (2026, 12)):
    for d_, mo in zip(qd(y, m), (1426000, 1478000)): ev(iso(d_), "ingreso", "Quincena de Axel", mo, estado="capturada" if iso(d_) < iso(HOY) else "esperada")
    for f in FIJOS:
        fd = date(y, m, f["dia_del_mes"])
        if fd > HOY: ev(iso(fd), "gasto_fijo", f["concepto"], f["monto_cents"], venidero=y == 2026 and m == 10)
        elif (y, m) == (2026, 10): ev(iso(fd), "gasto_fijo", f["concepto"], f["monto_cents"], venidero=False, estado="realizado")  # gastos fijos de octubre ya pagados (coinciden con gastos[])
    ev(iso(date(y, m, 17)), "corte_tarjeta", "Corte de Banorte Oro"); ev(iso(date(y, m, 6)), "corte_tarjeta", "Fecha límite de pago de Banorte Oro") if date(y, m, 6) > HOY else None
for p in PAGOS:
    if date(2026, 10, 1) <= D(p["fecha"]) <= date(2026, 12, 31) and p["compra_fecha"] <= iso(HOY) or (p["compra_id"] == "kc3" and p["fecha"] <= "2026-12-31"): ev(p["fecha"], "pago_credito", f"Pago {p['linea_id']} · {p['compra_id']}", p["monto_cents"], hogar_cents=p["hogar_cents"], karen_cents=p["karen_cents"], estado=hoy_estado(p))
for g in GAS:
    if g["estado"] == "venidero" and g["tipo"] == "venidero" and g["categoria_id"] == "salud": ev(g["fecha"], "gasto_venidero", g["comercio"], g["monto_cents"])
ev("2026-12-12", "despensa", "Despensa bimestral programada", 300000)
for g in GAS:  # lo ya realizado en octubre: la despensa a crédito y las visitas a mercado con ticket (fechas de gastos[])
    if g["estado"] == "realizado" and g["categoria_id"] == "despensa-comun" and g["forma_pago"] == "kueski": ev(g["fecha"], "despensa", "Despensa bimestral (Walmart)", g["monto_cents"], estado="realizada", gasto_id=g["id"])
    if g["estado"] == "realizado" and g["tienda_id"] in ("mercado-portales", "tianguis-narvarte"): ev(g["fecha"], "visita_mercado", "Visita a " + ("Mercado Portales" if g["tienda_id"] == "mercado-portales" else "Tianguis de la Narvarte"), None, estado="realizada", gasto_id=g["id"])
d_ = HOY + timedelta(1)
while d_ <= date(2026, 11, 8):
    if d_.weekday() == 5: ev(iso(d_), "visita_mercado", "Visita a Mercado Portales")
    if d_.weekday() == 1: ev(iso(d_), "visita_mercado", "Visita al Tianguis de la Narvarte")
    d_ += timedelta(1)
for p in PROMOS: ev(p["inicio"], "promocion", f"Promoción {p['comercio']} ({p['etiqueta_tasa']}): inicia", None, fin=p["fin"])
for m in METAS: ev(m["fecha_meta"], "meta_ahorro", f"Fecha meta: {m['nombre']}", m["monto_meta_cents"])
for f in ("2026-10-01", "2026-11-01", "2026-12-01"): ev(f, "reporte", "Envío del reporte mensual (08:00)", estado="enviado" if f < iso(HOY) else "programado")
TORD = ["ingreso", "pago_credito", "gasto_fijo", "gasto_venidero", "corte_tarjeta", "despensa", "visita_mercado", "promocion", "meta_ahorro", "reporte"]
CAL.sort(key=lambda e: (e["fecha"], TORD.index(e["tipo"])))
# ---------------- ajustes, notificaciones y reportes
def cg(f): return sum(g["hogar_cents"] for g in GAS if g["fecha"] == f and g["estado"] == "realizado")
def cn(f): return sum(1 for g in GAS if g["fecha"] == f and g["estado"] == "realizado")  # conteo de gastos del día: sale de gastos, igual que el monto
def txt_dia(f, etiqueta): return f"Resumen del {etiqueta}: gastaste {fm(cg(f))} (parte del hogar) en {cn(f)} " + ("gasto" if cn(f) == 1 else "gastos") + "."
AVISOS = [dict(enviado="2026-10-21T09:00:00-06:00", tipo="alerta de presupuesto al 80 %", destinatarios=["axel", "ximena"], resultado="entregado", texto=f"Frutas y verduras (mercado) va en 94 % de su plan de octubre: llevas {fm(CAT[4]['real_cents'])} de {fm(CAT[4]['plan_cents'])}."),
 dict(enviado="2026-10-20T21:30:00-06:00", tipo="resumen diario", destinatarios=["axel", "ximena"], resultado="falló: sin conexión con Telegram", texto=txt_dia("2026-10-20", "20-oct-2026")),
 dict(enviado="2026-10-19T21:30:00-06:00", tipo="resumen diario", destinatarios=["axel", "ximena"], resultado="entregado", texto=txt_dia("2026-10-19", "19-oct-2026")),
 dict(enviado="2026-10-19T09:00:00-06:00", tipo="alerta de presupuesto al 80 %", destinatarios=["axel", "ximena"], resultado="entregado", texto=f"Despensa común va en 83 % de su plan de octubre: {fm(CAT[2]['real_cents'])} de {fm(CAT[2]['plan_cents'])}."),
 dict(enviado="2026-10-18T19:00:00-06:00", tipo="resumen semanal", destinatarios=["axel", "ximena"], resultado="entregado", texto="Resumen de la semana del 12 al 18 de octubre."),
 dict(enviado="2026-10-18T09:00:00-06:00", tipo="alerta de plan sin armar", destinatarios=["ximena"], resultado="entregado", texto="Te faltan 2 casillas por armar en el plan de la semana del 19 al 25 de octubre."),
 dict(enviado="2026-10-14T09:00:00-06:00", tipo="alerta de pago próximo", destinatarios=["axel", "ximena"], resultado="entregado", texto="Mañana vence el pago de Kueski Pay de la compra anterior."),
 dict(enviado="2026-10-02T09:00:00-06:00", tipo="alerta de presupuesto al 100 %", destinatarios=["axel", "ximena"], resultado="entregado", texto="Vivienda llegó al 100 % de su plan de octubre."),
 dict(enviado="2026-10-01T08:00:00-06:00", tipo="reporte mensual", destinatarios=["axel", "ximena"], resultado="entregado", texto="Reporte mensual de septiembre de 2026 listo.")]
REPORTES = [dict(id="rp1", tipo="mensual PDF", periodo="2026-09", generado="2026-10-01T08:00:00-06:00", origen="programado", estado="listo", envio=dict(a=["axel", "ximena"], cuando="2026-10-01T08:01:00-06:00", resultado="entregado")),
 dict(id="rp2", tipo="estado de cuenta de Karen", periodo="2026-09", generado="2026-10-02T12:20:00-06:00", origen="manual (Axel)", estado="listo", envio=None),
 dict(id="rp3", tipo="mensual PDF", periodo="2026-08", generado="2026-09-01T08:00:00-06:00", origen="programado", estado="error: base de datos ocupada", envio=None)]
AJ = dict(usuarios=[dict(id="axel", nombre="Axel", usuario="axel", rol="Co-titular", ultimo_acceso="2026-10-21T07:55:00-06:00"), dict(id="ximena", nombre="Ximena", usuario="ximena", rol="Co-titular", ultimo_acceso="2026-10-20T22:10:00-06:00")],
 sesiones=[dict(dispositivo="Pixel 8 · Chrome", usuario="axel", inicio="2026-09-22T19:30:00-06:00", ultimo_uso="2026-10-21T07:55:00-06:00", vence="2026-12-21", actual=True), dict(dispositivo="iPhone · Safari", usuario="ximena", inicio="2026-09-23T08:12:00-06:00", ultimo_uso="2026-10-20T22:10:00-06:00", vence="2026-12-22", actual=False),
   dict(dispositivo="Laptop · Firefox", usuario="axel", inicio="2026-10-03T21:00:00-06:00", ultimo_uso="2026-10-17T20:40:00-06:00", vence="2027-01-01", actual=False)],
 parametros=dict(intentos_fallidos=5, ventana_min=15, sesion_dias=90, password_largo_min=12, umbral_aviso_pct=80, umbral_limite_pct=100, pago_proximo_dias_antes=[1, 0], precio_viejo_dias=30, precio_peso_cero_dias=90, adjunto_max_mb=10, idempotencia_horas=24, profeco_max_gb=1),
 reportes_programados=dict(resumen_diario=dict(activo=True, hora="21:30"), resumen_semanal=dict(activo=True, dia="domingo", hora="19:00"), mensual_pdf=dict(activo=True, dia_del_mes=1, hora="08:00", enviar_telegram=True, guardar=True), destinatarios=["axel", "ximena"], avisos_hora="09:00"),
 proximos_envios=[dict(tipo="resumen diario", cuando="2026-10-21T21:30:00-06:00"), dict(tipo="resumen semanal", cuando="2026-10-25T19:00:00-06:00"), dict(tipo="reporte mensual", cuando="2026-11-01T08:00:00-06:00")],
 telegram=dict(confirmados=[dict(persona_id="axel", confirmado_el="2026-09-22"), dict(persona_id="ximena", confirmado_el="2026-09-22")], detectados=[dict(nombre_telegram="Ximena (tablet)", detectado="2026-10-20T18:05:00-06:00")]),
 respaldos=[dict(fecha=f"2026-10-{d}T03:00:00-06:00", tamano_mb=round(2.3 + i * 0.1, 1), integrity_check="ok", ubicacion="disco local") for i, d in enumerate(("21", "20", "19", "18", "17"))],
 proximo_respaldo="2026-10-22T03:00:00-06:00", ia=[dict(proveedor="Modelo local (Ollama)", activo=False), dict(proveedor="API de Claude", activo=False, llave="no configurada")], servicio=dict(version="1.0.0-demo", en_marcha_desde="2026-10-12", base_de_datos_mb=3.1),
 profeco=dict(ultimo_corte="2026-10-02", registros=48210, cobertura=[dict(cadena="Walmart", pct=92), dict(cadena="Bodega Aurrera", pct=88), dict(cadena="Sam's Club", pct=71)]), tema="sistema")
# Ajustes › Notificaciones (P-26.09, N-298/N-300/N-323), Datos semilla (P-26.11, N-305) y tiempo del comparador (P-26.05, N-321)
AJ["notificaciones"] = dict(nota="Las 4 alertas de la decisión F; valores de K44, K47 y K50 (anticipación de precios viejos 1 día, hora de avisos 09:00). Los chats de Telegram están en ajustes.telegram",
  alertas=[dict(id="payment_due", nombre="Pago próximo", activa=True, destinatarios=["axel", "ximena"], hora="09:00", anticipacion=dict(clave="alerts.payment_due_days_before", dias_antes=[1, 0], valor_inicial=[1, 0], texto="1 día antes y el mismo día")),
    dict(id="budget_threshold", nombre="Presupuesto al 80 % y al 100 %", activa=True, destinatarios=["axel", "ximena"], hora="09:00", anticipacion=None,
         umbrales=dict(aviso_pct=80, tope_pct=100, claves=["alerts.budget_warning_pct", "alerts.budget_limit_pct"], valores_iniciales=[80, 100]), regla="Mira cada categoría y el total"),
    dict(id="stale_prices", nombre="Precios viejos antes del día de despensa", activa=False, destinatarios=["axel"], hora="09:00", anticipacion=dict(clave="advance_days", advance_days=1, valor_inicial=1, texto="1 día antes del día de despensa")),
    dict(id="meal_plan_missing", nombre="Plan de comidas sin armar", activa=True, destinatarios=["ximena"], hora="09:00", anticipacion=dict(clave="weekday", weekday=7, dia="domingo", valor_inicial=7, texto="Se revisa el domingo"), regla="Casillas sin armar = día × tiempo obligatorio sin receta ni «fuera»")],
  destinatarios_posibles=["axel", "ximena"], aviso_de_prueba=dict(ultimo="2026-10-12T20:30:00-06:00", a="axel", resultado="entregado"))
AJ["datos_semilla"] = [dict(id="compra-semilla", nombre="Compra de despensa de la semilla", estado="Cargada", cargada_el="2026-10-12T20:02:00-06:00"),
  dict(id="credito-kueski", nombre="Crédito Kueski", estado="Cargada", cargada_el="2026-10-12T20:03:00-06:00"),
  dict(id="gastos-fijos-ingresos", nombre="Gastos fijos e ingresos (archivo CSV)", estado="Pendiente de importar", cargada_el=None)]
AJ["parametros"]["comparador_max_s"] = 5; AJ["parametros"]["comparador_max_s_limites"] = [1, 30]
# ---------------- tickets: renglones (N-061), lectura OCR propuesta (N-071) y tickets por confirmar (N-076)
CUENTA_C50 = "Conjunto 50 %"
def reng(orden, pid, qid, tienda, cant, cuenta, **extra):
    pr, e = PBY[pid], PREC[(pid, qid, tienda)]; pu = e["precio_cents"]; imp = rh(round(cant * 1000) * pu, 1000)
    h, k = (imp, 0) if cuenta == "Hogar" else (0, imp) if cuenta == "Karen" else split(imp, 50)
    pres = [x for x in pr["presentaciones"] if x["id"] == qid][0]
    return dict(orden=orden, producto_id=pid, producto=pr["nombre"], presentacion_id=qid, presentacion=pres["etiqueta"], cantidad=cant, unidad="kg" if cant != int(cant) else "pza",
        precio_unitario_cents=pu, importe_cents=imp, cuenta=cuenta, hogar_cents=h, karen_cents=k, es_ajuste=False, **extra)
def ajuste(orden, imp, cuenta, **extra):  # renglón «Ajuste del ticket» (descuentos, redondeo)
    h, k = (imp, 0) if cuenta == "Hogar" else (-rh(-imp * 50, 100), imp + rh(-imp * 50, 100)) if imp < 0 else split(imp, 50)
    return dict(orden=orden, producto_id=None, producto="Ajuste del ticket", presentacion_id=None, presentacion=None, cantidad=1, unidad="pza", precio_unitario_cents=imp, importe_cents=imp, cuenta=cuenta, hogar_cents=h, karen_cents=k, es_ajuste=True, **extra)
def por_cuenta(rs):
    out = {}
    for r in rs: c = out.setdefault(r["cuenta"], dict(cuenta=r["cuenta"], importe_cents=0, hogar_cents=0, karen_cents=0)); c["importe_cents"] += r["importe_cents"]; c["hogar_cents"] += r["hogar_cents"]; c["karen_cents"] += r["karen_cents"]
    return list(out.values())
G13 = [g for g in GAS if g["id"] == "g13"][0]
R13 = [reng(1, "arroz", "b900", "walmart", 6, CUENTA_C50), reng(2, "frijol", "b900", "walmart", 4, CUENTA_C50), reng(3, "aceite", "l1", "walmart", 6, CUENTA_C50), reng(4, "atun", "l140", "walmart", 12, CUENTA_C50),
       reng(5, "papel", "r12", "walmart", 5, CUENTA_C50), reng(6, "detergente", "l3", "walmart", 2, CUENTA_C50), reng(7, "cafe", "f200", "walmart", 2, CUENTA_C50), reng(8, "leche", "l1", "walmart", 30, CUENTA_C50),
       reng(9, "huevo", "d12", "walmart", 6, CUENTA_C50), reng(10, "azucar", "kg", "walmart", 4, CUENTA_C50), reng(11, "pasta", "p200", "walmart", 11, CUENTA_C50)]
R13.append(ajuste(12, G13["monto_cents"] - sum(r["importe_cents"] for r in R13), CUENTA_C50))
TK1 = dict(id="tk1", gasto_id="g13", tienda_id="walmart", sucursal="Sucursal de ejemplo", fecha=G13["fecha"], total_cents=G13["monto_cents"], forma_pago="kueski", compra_credito_id="kc1", estado="confirmado",
    adjunto=dict(tipo="imagen", nombre="ticket-walmart-2026-10-12.jpg", tamano_mb=1.8), confirmado_por="axel", confirmado_el="2026-10-12T19:40:00-06:00", renglones=R13,
    suma_productos_cents=sum(r["importe_cents"] for r in R13 if not r["es_ajuste"]), ajuste_cents=R13[-1]["importe_cents"], por_cuenta=por_cuenta(R13),
    hogar_cents=sum(r["hogar_cents"] for r in R13), karen_cents=sum(r["karen_cents"] for r in R13), registra_precios="Cada renglón registró su precio con fuente «ticket» y la fecha del ticket (N-074)")
PROD13 = {r["producto_id"] for r in R13 if r["producto_id"]}
def sug(pid): return (CUENTA_C50, "última compra confirmada") if pid in PROD13 else ("Hogar", "reparto de la categoría Súper de la semana")
def ocr(orden, texto, pid, qid, tienda, cant, conf, pu=None, **extra):
    r = reng(orden, pid, qid, tienda, cant, sug(pid)[0]); r.pop("cuenta"); r.pop("hogar_cents"); r.pop("karen_cents"); r.pop("es_ajuste")
    cs, orig = sug(pid)
    if pu is not None: r["precio_unitario_cents"] = pu; r["importe_cents"] = rh(round(cant * 1000) * pu, 1000)
    r.update(texto_leido=texto, cuenta_sugerida=cs, origen_cuenta=orig, confianza=conf, **extra)
    h, k = (r["importe_cents"], 0) if cs == "Hogar" else split(r["importe_cents"], 50); r.update(hogar_cents=h, karen_cents=k)
    return r
RO = [ocr(1, "JITOMATE SALADET KG 0.850", "jitomate", "kg", "walmart", 0.85, "alta"), ocr(2, "CEBOLLA BCA KG 0.600", "cebolla", "kg", "walmart", 0.6, "alta"),
      ocr(3, "PECH POLLO KG 1.250", "pollo", "kg", "walmart", 1.25, "media"), ocr(4, "HUEVO BLANCO 12 PZ", "huevo", "d12", "walmart", 1, "alta"),
      ocr(5, "LECHE ENT 1L X6", "leche", "l1", "walmart", 6, "alta"), ocr(6, "ACEITE VEG 1-2-3 1L", "aceite", "l1", "walmart", 1, "alta"),
      ocr(7, "JAB TRAST AXN 75O", "jabon-trastes", "m750", "walmart", 1, "media"), ocr(8, "AGUACATE HAS KG 0.480", "aguacate", "kg", "walmart", 0.48, "alta"),
      ocr(9, "FRJ N6RO VV 9O0G", "frijol", "b900", "walmart", 1, "baja", alternativas=[dict(producto_id="arroz", presentacion_id="b900", producto="Arroz blanco")], motivo_confianza="El texto leído no coincide con ningún producto; se sugiere el más parecido"),
      ocr(10, "T0RTIL L4 MA1Z KG", "tortilla", "kg", "walmart", 1, "baja", pu=2230, motivo_confianza="El precio leído «22.3O» es ambiguo (la letra O puede ser un cero)")]
SUMA_O = sum(r["importe_cents"] for r in RO); TOT_O = SUMA_O - 1590
OCR1 = dict(id="ocr1", estado="pendiente de confirmar", tienda_id="walmart", sucursal="Sucursal de ejemplo", fecha="2026-10-21", foto=dict(nombre="foto-ticket-walmart-2026-10-21.jpg", tomada="2026-10-21T18:42:00-06:00", tamano_mb=2.4), motor="OCR local (Tesseract, español)",
    categoria_sugerida_id="super", total_leido_cents=TOT_O, suma_renglones_cents=SUMA_O, diferencia_cents=TOT_O - SUMA_O, diferencia_texto="El total del ticket es menor que la suma de los renglones: se propone un renglón «Ajuste del ticket» (descuentos, redondeo)",
    renglones=RO, bajos=[r["orden"] for r in RO if r["confianza"] == "baja"], tiene_borrador=False,
    nota="Nada se guarda sin confirmación humana (decisión E): confirmar crea gasto, renglones, precios y entradas de inventario")
R2 = [reng(1, "leche", "l1", "bodega-aurrera", 6, "Hogar"), reng(2, "arroz", "b900", "bodega-aurrera", 2, CUENTA_C50), reng(3, "frijol", "b900", "bodega-aurrera", 2, "Hogar"), reng(4, "pasta", "p200", "bodega-aurrera", 4, "Hogar"), reng(5, "atun", "l140", "bodega-aurrera", 4, "Hogar")]
SUMA_2 = sum(r["importe_cents"] for r in R2)
OCR2 = dict(id="ocr2", estado="pendiente de confirmar", tienda_id="bodega-aurrera", sucursal="Sucursal de ejemplo", fecha="2026-10-19", foto=dict(nombre="foto-ticket-bodega-2026-10-19.jpg", tomada="2026-10-19T20:05:00-06:00", tamano_mb=1.9), tiene_borrador=True,
    total_propuesto_cents=SUMA_2 - PREC[("atun", "l140", "bodega-aurrera")]["precio_cents"], borrador=dict(guardado_el="2026-10-19T20:15:00-06:00", total_cents=SUMA_2, renglones=R2,
    correcciones=[dict(renglon=5, campo="cantidad", de=3, a=4), dict(renglon=2, campo="cuenta", de="Hogar", a=CUENTA_C50)], por_cuenta=por_cuenta(R2)))
TICKETS = dict(nota="Datos de P-04 (renglones del gasto g13) y de P-05 (lectura propuesta y tickets por confirmar). Solo g13 trae renglones capturados; los demás gastos con tiene_ticket muestran el adjunto sin renglones.",
    confirmados=[TK1], lectura_pendiente=OCR1,
    por_confirmar=[dict(id=o["id"], fecha_foto=o["fecha"], tienda_id=o["tienda_id"], total_cents=(o["borrador"]["total_cents"] if o["tiene_borrador"] else o["total_leido_cents"]), tiene_borrador=o["tiene_borrador"]) for o in (OCR1, OCR2)],
    borradores=[OCR2])
# ---------------- armado
def addfmt(o, skip=("precios", "tiempos", "ingredientes", "historial_peso")):
    if isinstance(o, dict):
        for k, v in list(o.items()):
            if k in skip: continue
            if k.endswith("_cents") and isinstance(v, int) and not isinstance(v, bool): o[k[:-6] + "_fmt"] = fm(v)
            else: addfmt(v, skip)
    elif isinstance(o, list):
        for x in o: addfmt(x, skip)
MESES_H = [("2026-05", 2540000, 2388310), ("2026-06", 2540000, 2602475), ("2026-07", 2610000, 2447920), ("2026-08", 2610000, 2731860), ("2026-09", 2610000, 2512640), ("2026-10", PLAN, REALT)]
# ---------------- series sintéticas para las vistas de Tendencias y para P-12 / P-26.11 (ronda 2: M-2)
def prop(total, pesos):  # reparte un entero en partes proporcionales (mayor resto); la suma es exacta
    sp = sum(pesos); base = [total * w // sp for w in pesos]; resto = total - sum(base)
    orden = sorted(range(len(pesos)), key=lambda i: (-(total * pesos[i] % sp), i))
    for i in orden[:resto]: base[i] += 1
    return base
MESES_ANT = [m for (m, _, _) in MESES_H[:-1]]
# ingresos por mes (P-08 › Tendencias, N-103): antes de rige_desde la quincena era distinta; octubre es el del archivo
ING_TEND = []
for (m, pl, rl) in MESES_H:
    if m == "2026-10":
        ING_TEND.append(dict(mes=m, quincena_15_cents=ING["quincenas"][0]["monto_cents"], quincena_30_cents=ING["quincenas"][1]["monto_cents"], extras_cents=sum(e["monto_cents"] for e in ING["extras"]), total_cents=ING_MES, incluye_esperado=True)); continue
    h = h8("ing" + m); f = 0.955 + (h % 20) / 1000.0
    q15 = int(round(ING["esperado"]["q15_cents"] * f / 1000.0)) * 1000; q30 = int(round(ING["esperado"]["q30_cents"] * f / 1000.0)) * 1000
    ex = 0 if h % 3 == 0 else (60000 + (h >> 4) % 9 * 10000)
    ING_TEND.append(dict(mes=m, quincena_15_cents=q15, quincena_30_cents=q30, extras_cents=ex, total_cents=q15 + q30 + ex, incluye_esperado=False))
# gasto del hogar por categoría y mes (P-02 › Tendencias N-044 y P-06 N-085): cada mes suma el plan y el real del mes de presupuesto.tendencia
TPC_MESES = []
for (m, pl, rl) in MESES_H:
    if m == "2026-10":
        cats = [dict(categoria_id=c["id"], plan_cents=c["plan_cents"], real_cents=c["real_cents"], usado_pct=c["usado_pct"]) for c in CAT]
    else:
        fijas = [c for c in CAT if c["tipo"] == "fijo"]; vars_ = [c for c in CAT if c["tipo"] == "variable"]
        plan_v = prop(pl - sum(c["plan_cents"] for c in fijas), [c["plan_cents"] for c in vars_]); plan = {c["id"]: c["plan_cents"] for c in fijas}; plan.update({c["id"]: v for c, v in zip(vars_, plan_v)})
        real = {}
        for c in fijas:  # lo fijo se paga igual; los seguros no se cobran en ningún mes de esta serie (como en octubre)
            real[c["id"]] = 0 if c["real_cents"] == 0 else (c["plan_cents"] if c["id"] == "vivienda" else int(round(c["plan_cents"] * (0.55 + (h8(m + c["id"]) % 12) / 100.0) / 100.0)) * 100)
        pesos = [int(c["plan_cents"] * (88 + h8(m + c["id"]) % 25)) for c in vars_]
        for c, v in zip(vars_, prop(rl - sum(real.values()), pesos)): real[c["id"]] = v
        cats = [dict(categoria_id=c["id"], plan_cents=plan[c["id"]], real_cents=real[c["id"]], usado_pct=rh(real[c["id"]] * 100, plan[c["id"]])) for c in CAT]
    TPC_MESES.append(dict(mes=m, plan_cents=pl, real_cents=rl, por_categoria=cats))
SEIS_MAYORES = [c["id"] for c in sorted(CAT, key=lambda c: (-c["plan_cents"], c["id"]))[:6]]  # las 6 de mayor plan; el resto va en «Otras»
for t in TPC_MESES: t["otras_real_cents"] = sum(c["real_cents"] for c in t["por_categoria"] if c["categoria_id"] not in SEIS_MAYORES); t["otras_plan_cents"] = sum(c["plan_cents"] for c in t["por_categoria"] if c["categoria_id"] not in SEIS_MAYORES)
TENDENCIA_POR_CATEGORIA = dict(nota="6 meses (mayo a octubre de 2026), no 12: los meses anteriores a mayo se dibujan como estado vacío «—». Cada mes suma el plan y el real de presupuesto.tendencia", seis_mayores=SEIS_MAYORES, meses=TPC_MESES)
# historial de precios por producto × presentación × tienda (P-15 › Tendencias, N-171): el último punto es el precio vigente de «precios»
def dia_ok(d): return d.day not in (9, 10, 15, 28, 30, 31)
for P_ in PRODS:
    ser = []
    for e in P_["precios"]:
        h = h8("hist" + P_["id"] + e["presentacion_id"] + e["tienda_id"]); pts = [dict(fecha=e["fecha"], precio_cents=e["precio_cents"], fuente=e["fuente"])]
        d = D(e["fecha"]); pr = e["precio_cents"]
        for k in range(2 + h % 3):
            d = d - timedelta(11 + (h >> (3 * k + 2)) % 17)
            while not dia_ok(d): d -= timedelta(1)
            if d < date(2026, 5, 2): break
            pr = max(10, int(round(pr * (0.94 + ((h >> (4 * k + 1)) % 15) / 100.0) / 10.0)) * 10)
            pts.insert(0, dict(fecha=iso(d), precio_cents=pr, fuente=e["fuente"] if e["fuente"] == "profeco" else ["ticket", "manual"][(h >> k) % 2]))
        ser.append(dict(tienda_id=e["tienda_id"], presentacion_id=e["presentacion_id"], puntos=pts))
    P_["historial_precios"] = ser
# uso de cada receta (P-12, N-141 y orden de N-122): veces en el plan en las últimas 8 semanas (hasta hoy) y próxima fecha planeada
for r in RECETAS:
    slots = [(d["fecha"], tn) for d in PLANSEM for tn in d["tiempos"] if any(d["tiempos"][tn][p].get("receta_id") == r["id"] for p in ("axel", "ximena"))]
    hasta_hoy = sum(1 for (f, _) in slots if f <= iso(HOY)); prev = h8("uso" + r["id"]) % 6; fut = sorted(f for (f, _) in slots if f >= iso(HOY))
    r["uso"] = dict(ventana_inicio=iso(HOY - timedelta(55)), ventana_fin=iso(HOY), en_semana_actual_hasta_hoy=hasta_hoy, en_7_semanas_previas=prev, veces_8_semanas=hasta_hoy + prev, proxima_fecha=fut[0] if fut else None)
FIJ_MES_Q1, FIJ_MES_Q2 = sum(f["monto_cents"] for f in FIJOS if f["dia_del_mes"] <= 15), sum(f["monto_cents"] for f in FIJOS if f["dia_del_mes"] > 15)
VARP = sum(c["plan_cents"] for c in CAT if c["tipo"] == "variable"); PLQ1 = FIJ_MES_Q1 + rh(VARP, 2)
REQ1 = sum(hog_mes(g) for g in REAL if g["fecha"] <= "2026-10-15"); REQ2 = REALT - REQ1
DIAS_FALTAN = calendar.monthrange(2026, 10)[1] - HOY.day + 1
dif = sorted(DET, key=lambda r: -r["diferencia_cents"])[:5]
MV_CARD = dict(month="2026-10", actual_cents=REALT, planned_cents=PLAN, used_pct=rh(REALT * 100, PLAN), remaining_cents=PLAN - REALT, status=est(rh(REALT * 100, PLAN)), days_left=DIAS_FALTAN)
KM = METAS[0]; GANA = COMP[0]; SEG = COMP[1]
NEXT_GOAL = dict(id=KM["id"], name=KM["nombre"], kind="goal", target_cents=KM["monto_meta_cents"], target_date=KM["fecha_meta"], saved_cents=KM["ahorrado_cents"], progress_pct=KM["avance_pct"],
  planned_monthly_cents=KM["aportacion_planificada_cents"], required_monthly_cents=KM["aportacion_necesaria_cents"], status={"en_ruta": "on_track", "atrasada": "behind"}[KM["estado"]], etiqueta_estado=KM["etiqueta_estado"])  # nombres de SavingsGoal (C1)
INICIO = dict(hoy=dict(fecha=iso(HOY), dia="Mié 21 oct 2026"),
 tarjetas=dict(mes_vs_plan=MV_CARD,
  flujo_quincenal=dict(start_date="2026-10-16", end_date="2026-10-31", income_cents=ING_Q2, income_is_expected=True, spent_cents=HECHOS, committed_cents=COMPR, closing_balance_cents=CIERRE),
  proximo_pago=dict(next_payment=PROX, availability=[dict(credit_id="kueski", availability_reading="B", available_cents=ULT["disponible_b_cents"], available_otra_lectura=dict(lectura="A", available_cents=ULT["disponible_a_cents"]), availability_confirmed=False)]),
  comidas_semana=dict(week_start=iso(SEM0), week_end=iso(SEM0 + timedelta(6)), cost_cents=COSTO_SEM, unplanned_slots=SIN_ARMAR,
    persons=[dict(person_id=p, avg_kcal_per_day=rh(sum(KCAL[p]), sum(1 for k in KCAL[p] if k)), days_with_plan=sum(1 for k in KCAL[p] if k), daily_kcal_target=[x for x in SAL["perfiles"] if x["persona_id"] == p][0]["meta_kcal_confirmada"]) for p in ("axel", "ximena")]),
  mejor_canasta=dict(basket="shopping_list", store_id=GANA["tienda_id"], total_cents=GANA["total_cents"], segunda_tienda_id=SEG["tienda_id"], difference_vs_second_cents=SEG["total_cents"] - GANA["total_cents"], covered_count=GANA["cubiertos"], item_count=len(ITEMS), stale_price_count=GANA["viejos"], assumes_membership=False),
  ahorro=dict(saved_total_cents=AHO["ahorrado_total_cents"], contributed_month_cents=APORT_OCT, planned_month_cents=AHO_PLAN, next_goal=NEXT_GOAL)),
 segundo_nivel=dict(categorias_mayor_uso=sorted([dict(id=c["id"], nombre=c["nombre"], usado_pct=c["usado_pct"], estado=c["estado"]) for c in CAT], key=lambda c: -c["usado_pct"])[:5],
   proximos_5_pagos=[dict(fecha=p["fecha"], linea_id=p["linea_id"], monto_cents=p["monto_cents"], hogar_cents=p["hogar_cents"], karen_cents=p["karen_cents"]) for p in PEND[:5]],
   comidas_hoy={tn: {per: PLANSEM[2]["tiempos"][tn][per] for per in ("axel", "ximena")} for tn in ("desayuno", "comida", "cena", "colacion")}, articulos_mayor_diferencia=[dict(producto=r["producto"], diferencia_cents=r["diferencia_cents"], ganadora=r["ganadora"]) for r in dif]),
 tendencias=dict(plan_vs_real=[dict(mes=m, plan_cents=p, real_cents=r, usado_pct=rh(r * 100, p)) for (m, p, r) in MESES_H], devengado_vs_efectivo=PROY, deuda_viva=DEU, costo_semanal_comidas=SAL["tendencia_semanal"]),
 karen=KAR)
GRES = dict(periodo="2026-10", num_gastos=len(REAL), suma_montos_cents=sum(g["monto_cents"] for g in REAL), suma_hogar_cents=sum(g["hogar_cents"] for g in REAL), suma_karen_cents=sum(g["karen_cents"] for g in REAL),
  por_categoria=sorted([dict(categoria_id=c["id"], hogar_devengado_cents=c["real_cents"], num_gastos=c["num_gastos"], pct_del_total=round(c["real_cents"] * 100.0 / REALT, 1)) for c in CAT], key=lambda c: -c["hogar_devengado_cents"]),
  anterior=dict(periodo="2026-09", hogar_cents=2512640, diferencia_cents=REALT - 2512640), venideros_cents=sum(g["hogar_cents"] for g in GAS if g["estado"] == "venidero"))
PRES = dict(resumen=dict(mes="2026-10", ingreso_mes_cents=ING_MES, plan_gasto_cents=PLAN, ahorro_planificado_cents=AHO_PLAN, sin_asignar_cents=ING_MES - PLAN - AHO_PLAN, real_cents=REALT, restante_cents=PLAN - REALT, usado_pct=MV_CARD["used_pct"], estado=MV_CARD["status"], dias_que_faltan=DIAS_FALTAN),
  quincenas=[dict(quincena="2026-10-q1", real_cents=REQ1, plan_cents=PLQ1), dict(quincena="2026-10-q2", real_cents=REQ2, plan_cents=PLAN - PLQ1)],
  despensa_devengada=[dict(compra_id="kc1", fecha="2026-10-12", total_con_intereses_cents=KC[1]["total_cents"], devengado_mes_cents=KC[1]["hogar"]["total_cents"] // 2 + KC[1]["karen"]["total_cents"] // 2 + (KC[1]["hogar"]["total_cents"] % 2 + KC[1]["karen"]["total_cents"] % 2), hogar_mes_cents=rh(KC[1]["hogar"]["total_cents"], 2), karen_mes_cents=rh(KC[1]["karen"]["total_cents"], 2))],
  alertas=[dict(categoria_id=c["id"], nombre=c["nombre"], usado_pct=c["usado_pct"], umbral=100 if c["usado_pct"] >= 100 else 80, cruzo=c["cruzo_100"] or c["cruzo_80"]) for c in CAT if c["usado_pct"] >= 80],
  tendencia=[dict(mes=m, plan_cents=p, real_cents=r, usado_pct=rh(r * 100, p)) for (m, p, r) in MESES_H], gastos_fijos_recurrentes=FIJOS, despensa_programada=dict(monto_base_cents=300000, dia_del_mes=12, frecuencia="bimestral", linea_id="kueski", plazo_quincenas=4, proxima_fecha="2026-12-12"),
  proyeccion_meses_7_a_12=dict(repetir_despensa_y_papel=False), metodo_variable=dict(texto="Estimado fijo de ejemplo por quincena (no es el método real del motor S2)", por_quincena_cents=VAR_EST))
CRED = dict(lineas=LINEAS, deuda_viva_total_cents=DEUDA_TOT, totales=dict(pagos_quincena_actual_cents=sum(p["monto_cents"] for p in pag_q(2026, 10, 2)), pagos_quincena_actual_hogar_cents=sum(p["hogar_cents"] for p in pag_q(2026, 10, 2)), pagos_mes_cents=sum(p["monto_cents"] for p in PAGOS if p["fecha"].startswith("2026-10")),
   intereses_mes_cents=DEU[0]["intereses_mes_cents"]), aviso_ab=dict(lectura_activa="B", texto="Reposición de capital pendiente de confirmar con Kueski", confirmada=False),
  kueski=dict(compras=KC[:3], compras_programadas=KC[3:], tasas_por_plazo=[dict(plazo_quincenas=2, tasa_pct="16.4150 %", rate_num=164150, rate_den=1000000), dict(plazo_quincenas=4, tasa_pct="24.3120 %", rate_num=243120, rate_den=1000000)],
    serie_disponible=dict(eventos=SER, min_a_cents=min(e["disponible_a_cents"] for e in SER), min_b_cents=min(e["disponible_b_cents"] for e in SER)), disponible_hoy=dict(A_cents=ULT["disponible_a_cents"], B_cents=ULT["disponible_b_cents"], minimo_A_cents=min(e["disponible_a_cents"] for e in SER), minimo_B_cents=min(e["disponible_b_cents"] for e in SER)),
    semantica="B = límite − capital vivo; A = límite libre desde la fecha ancla, las compras posteriores restan y los pagos suman encima (ilustrativo; el motor S3 manda)"),
  tarjeta=dict(ciclo=dict(fecha_corte="2026-10-17", fecha_limite_pago="2026-11-06", saldo_al_corte_cents=610485, pago_sin_intereses_cents=610485, pago_minimo_cents=max(rh(610485 * 6, 100), 30000), saldo_actual_cents=714385),
    msi=[dict(comercio="Liverpool · laptop", mensualidad_cents=MSI, meses_restantes=7, meses_total=12)]),
  prestamo=dict(monto_original_cents=1800000, desembolso="2026-06-18", amortizacion=PRE), tienda=dict(compra="Colchón matrimonial", fecha="2026-08-22", total_cents=699000, abonos=TIE),
  pagos_programados=[dict(p, estado_hoy=hoy_estado(p)) for p in PAGOS if p["compra_fecha"] <= iso(HOY) or p["compra_id"] in ("kc3", "kc4")], promociones=PROMOS)
CMP = dict(compra=dict(descripcion="Licuadora y horno en Walmart", monto_cents=420000, quien="Hogar 100 %", fecha_deseada="2026-11-05", tolerancia_dias=7),
  restricciones=dict(colchon_minimo_cents=COLCHON, linea_colchon="kueski", max_compras_vivas=4, fecha_deuda_cero=None), pesos=dict(intereses=40, pago_por_quincena=30, disponible=30), tope_pago_quincena=dict(axel_cents=None, ximena_cents=None),
  escenarios=ESC, optimos=OPT, recomendacion="e1", descartados=[dict(id=e["id"], motivo=e["motivo_descarte"]) for e in ESC if not e["valido"]])
SAMSC = AHO_CAN; SAMSANUAL = SAMSC * 12; MEMB = 65000
GRES["tendencia_por_categoria"] = TENDENCIA_POR_CATEGORIA
CANASTA = dict(semana="2026-W43", origen="lista de compras de la semana", articulos=len(ITEMS), por_tienda=list(CAN.values()), tiendas_comparables=[c["tienda_id"] for c in COMP], ranking_nota="Se ordena por total solo entre tiendas con cobertura de 90 % o más",
  diferencia_contra_la_mas_barata=[dict(tienda_id=c["tienda_id"], total_cents=c["total_cents"] - COMP[0]["total_cents"], identicos_cents=c["identicos_cents"] - [x for x in COMP if x["tienda_id"] == COMP[0]["tienda_id"]][0]["identicos_cents"]) for c in COMP],
  ganadores=GAN, empates=EMP, detalle=sorted(DET, key=lambda r: -r["diferencia_cents"]),
  membresia_sams=dict(costo_anual_cents=MEMB, ahorro_por_canasta_cents=SAMSC, compras_al_ano_supuestas=12, ahorro_anual_cents=SAMSANUAL, punto_equilibrio_canastas=None if SAMSC <= 0 else -(-MEMB // SAMSC),
    veredicto="No conviene con esta canasta: Sam's no ahorra contra la mejor tienda sin Sam's" if SAMSC <= 0 else "Conviene si compras la canasta con la frecuencia supuesta"))
INDICE = {"P-00": [], "P-01": ["inicio", "calendario_hogar", "participacion_karen"], "P-02/P-03": ["gastos", "gastos_resumen", "gastos_resumen.tendencia_por_categoria", "reparto", "categorias"], "P-04/P-05": ["gastos", "tickets", "reparto", "categorias", "productos", "tiendas"], "P-06": ["presupuesto", "categorias", "gastos_resumen.tendencia_por_categoria"], "P-07": ["proyeccion", "deuda_viva", "presupuesto.metodo_variable"],
 "P-08": ["ingresos", "ingresos.tendencia"], "P-09": ["ahorro"], "P-10/P-11/P-12": ["recetas", "recetas[].uso", "plan_semana"], "P-13": ["lista_compras", "inventario"], "P-14": ["canasta"], "P-15": ["productos", "productos[].historial_precios", "tiendas"], "P-16": ["inventario"], "P-17/P-18/P-19": ["creditos", "promociones"], "P-20": ["comparador"],
 "P-21/P-22": ["salud"], "P-23/P-24/P-25": ["reportes", "reporte_mensual", "estado_cuenta_karen"], "P-26": ["ajustes", "ajustes.notificaciones", "ajustes.datos_semilla", "ajustes.parametros", "tiendas", "categorias"]}
# número de registros de cada renglón de Datos semilla (P-26.11, N-305): conteos de este mismo archivo
DS = {x["id"]: x for x in AJ["datos_semilla"]}
_nk = len(CRED["kueski"]["compras"]) + len(CRED["kueski"]["compras_programadas"]); _nf = len(FIJOS) + 2
DS["compra-semilla"].update(registros=len(PRODS), registros_texto=f"{len(PRODS)} productos con sus precios", registros_estado="cargados")
DS["credito-kueski"].update(registros=_nk, registros_texto=f"{_nk} compras de Kueski Pay (hechas y programadas)", registros_estado="cargados")
DS["gastos-fijos-ingresos"].update(registros=_nf, registros_texto=f"{_nf} registros por importar (gastos fijos y las 2 quincenas esperadas)", registros_estado="por importar")
DATOS = dict(
 meta=dict(titulo="Datos sintéticos de Finanzas Hogar", version="1.2", todo_inventado=True, hoy=iso(HOY), dia_semana="miércoles", zona_horaria="America/Mexico_City", moneda="MXN", mes_en_curso="2026-10", semana_en_curso=dict(iso="2026-W43", inicio=iso(SEM0), fin=iso(SEM0 + timedelta(6))), quincena_en_curso=dict(id="2026-10-q2", inicio="2026-10-16", fin="2026-10-31"),
   aviso="Cifras y fechas inventadas con la forma y magnitud de las reales (K45.1). No usar como dato real. Ver LEEME.md."),
 formato=dict(dinero="Centavos enteros en campos *_cents; el texto $1,234.56 va en el campo *_fmt hermano o se arma con la misma regla (miles con coma, 2 decimales, signo menos «−$»).", fecha_iso="AAAA-MM-DD", fecha_pantalla="21-oct-2026 · Mié 21 oct 2026", porcentaje="entero ROUND_HALF_UP (80 %); las tasas llevan sus decimales como texto"),
 indice_pantallas=INDICE,
 personas=PER, cuentas_y_formas_de_pago=dict(reparto=[dict(id="hogar", nombre="Hogar", hogar_pct=100), dict(id="karen", nombre="Karen", hogar_pct=0), dict(id="conjunto", nombre="Conjunto", hogar_pct=50, nota="50/50 por defecto; otros % se capturan por gasto", token_color="color.person.shared")],
   formas_de_pago=[dict(id="contado", nombre="De contado")] + [dict(id=l["id"], nombre=l["nombre"], tipo=l["tipo"]) for l in LINEAS]),
 ingresos=dict(ING, total_mes_cents=ING_MES, tendencia=ING_TEND), categorias=CAT, presupuesto=PRES, gastos=GAS, gastos_resumen=GRES,
 reparto=dict(regla="La parte del hogar se redondea ROUND_HALF_UP y la de Karen es el resto: las partes siempre suman el total", ejemplo_50_50="g13", ejemplo_distinto="g20", defecto_conjunto=dict(hogar_pct=50, karen_pct=50)),
 ahorro=AHO, creditos=CRED, promociones=PROMOS, comparador=CMP, tiendas=TIENDAS, productos=PRODS, canasta=CANASTA, recetas=RECETAS,
 plan_semana=dict(semana="2026-W43", dias=PLANSEM, totales=dict(costo_semana_cents=COSTO_SEM, costo_por_dia_cents=[SEMCOST["axel"][i] + SEMCOST["ximena"][i] for i in range(7)], kcal_por_dia=KCAL, casillas_sin_armar=SIN_ARMAR, meta_kcal={p["persona_id"]: p["meta_kcal_confirmada"] for p in SAL["perfiles"]})),
 lista_compras=dict(semana="2026-W43", total_cents=sum(g["total_cents"] for g in LISTA), grupos=LISTA), inventario=dict(articulos=INVL, resumen={e: sum(1 for x in INVL if x["estado"] == e) for e in ORD}, movimientos=INV_MOV),
 salud=SAL, reportes=dict(generados=REPORTES, proximos_envios=AJ["proximos_envios"], avisos_enviados=AVISOS, programados=AJ["reportes_programados"]),
 reporte_mensual=dict(mes="2026-10", parcial=True, meses_disponibles=[m for (m, _, _) in MESES_H], resumen=dict(ingreso_cents=ING_MES, real_cents=REALT, plan_cents=PLAN, ahorro_aportado_cents=APORT_OCT), presupuesto_por_categoria=[dict(id=c["id"], plan_cents=c["plan_cents"], real_cents=c["real_cents"], usado_pct=c["usado_pct"]) for c in CAT],
   creditos=dict(pagos_mes_cents=CRED["totales"]["pagos_mes_cents"], deuda_viva_cierre_cents=DEU[0]["total_cents"], intereses_cents=DEU[0]["intereses_mes_cents"]), comidas=dict(costo_semana_cents=COSTO_SEM), karen=KAR),
 estado_cuenta_karen=EC, tickets=TICKETS, proyeccion=dict(saldo_inicial_cents=SALDO_INI, flujo_por_quincena=FQ, meses=PROY, supuestos=dict(ingreso_esperado_q15_cents=ING["esperado"]["q15_cents"], ingreso_esperado_q30_cents=ING["esperado"]["q30_cents"], variable_estimado_por_quincena_cents=VAR_EST, meses_7_a_12="solo compromisos contratados, fijos e ingresos")),
 deuda_viva=dict(serie_mensual=DEU, a_fecha=[deuda("2026-10-21"), deuda("2026-12-31")], por_defecto_fin_de_horizonte=DEU[-1]["fecha"]), calendario_hogar=CAL, participacion_karen=KAR, inicio=INICIO, ajustes=AJ,
 estados_alternos=dict(pago_vencido=dict(due_date="2026-10-15", credit_id="kueski", parte="karen", amount_cents=38184, days_left=-6, is_overdue=True, texto="Sin marcar como pagado · venció hace 6 días", nota="Variante de la tarjeta Próximo pago (mismos nombres de campo de C1); los datos por defecto no la tienen"),
   sin_plan=dict(texto="Aún no hay presupuesto para octubre"), sin_precio=dict(texto="—")),
 recientes=[dict(tipo="gasto", id="g29", titulo="Bodega Aurrera · lácteos y reposición"), dict(tipo="receta", id="r03", titulo="Pollo asado con arroz y ensalada"), dict(tipo="crédito", id="kueski", titulo="Kueski Pay"), dict(tipo="pantalla", id="P-06", titulo="Presupuesto del mes"),
   dict(tipo="producto", id="aguacate", titulo="Aguacate Hass"), dict(tipo="meta", id="laptop", titulo="Laptop nueva"), dict(tipo="pantalla", id="P-14", titulo="Comparar canasta"), dict(tipo="gasto", id="g13", titulo="Walmart · despensa bimestral")])
def compactar(D_):
    D_["inicio"].pop("karen"); D_["inicio"]["tendencias"] = dict(ver="plan_vs_real = presupuesto.tendencia; devengado_vs_efectivo = proyeccion.meses; deuda_viva = deuda_viva.serie_mensual; costo_semanal_comidas = salud.tendencia_semanal")
    for c in D_["creditos"]["kueski"]["compras"] + D_["creditos"]["kueski"]["compras_programadas"]: c.pop("pagos"); c["pagos"] = "ver creditos.pagos_programados (compra_id = id)"
    for p in D_["creditos"]["pagos_programados"]:
        for k in ("id", "estado_global", "compra_fecha", "hogar_pagada_el", "karen_pagada_el"): p.pop(k, None)
    for g in D_["gastos"]:
        if g["forma_pago"] == "contado": g["efectivo"] = None
        if len(g["devengado"]) == 1 and g["forma_pago"] != "kueski": g["devengado"] = None
        g.pop("compra_credito_id") if g["compra_credito_id"] is None else None
    for p in D_["productos"]:
        for e in p["precios"]:
            e["precio_unidad_base_cents"] = round(e["precio_unidad_base_cents"], 2); e.pop("solo_referencia")
    for r in D_["canasta"]["detalle"]:
        for v in r["precios"].values():
            for k in ("fuente", "presentacion_id"): v.pop(k, None)
    for r in D_["recetas"]:
        r["kcal_fuente"] = "USDA Foundation Foods 04/2026 (ejemplo)"; r["pasos"] = ["Prepara los ingredientes.", "Cocina según el tiempo indicado.", "Sirve."]
        for i in r["ingredientes"]:
            for k in ("producto", "kcal_fuente", "presentacion_sugerida"): i.pop(k)
    for m in D_["creditos"]["tienda"]["abonos"] + D_["creditos"]["prestamo"]["amortizacion"]: pass
def main(out):
    compactar(DATOS)
    for k in ("inicio", "categorias", "gastos_resumen", "ahorro", "estado_cuenta_karen", "participacion_karen", "deuda_viva", "reporte_mensual", "tickets"): addfmt(DATOS[k])
    addfmt(DATOS["presupuesto"]["resumen"]); addfmt(DATOS["creditos"]["lineas"]); addfmt(DATOS["creditos"]["totales"]); addfmt(DATOS["creditos"]["tarjeta"]); addfmt(DATOS["creditos"]["kueski"]["disponible_hoy"]); addfmt(DATOS["creditos"]["kueski"]["compras"])
    addfmt(DATOS["comparador"]); addfmt(DATOS["canasta"]["por_tienda"]); addfmt(DATOS["canasta"]["membresia_sams"]); addfmt(DATOS["lista_compras"]["grupos"]); addfmt(DATOS["proyeccion"]["meses"]); addfmt(DATOS["gastos"][:3])
    with open(out, "w", encoding="utf-8") as f:
        f.write("{\n" + ",\n".join(json.dumps(k, ensure_ascii=False) + ": " + json.dumps(v, ensure_ascii=False, separators=(",", ":")) for k, v in DATOS.items()) + "\n}\n")
if __name__ == "__main__": main(sys.argv[1])
