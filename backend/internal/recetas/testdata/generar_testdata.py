"""Genera los datos de prueba de internal/recetas (5 recetas propias, catálogo mínimo y tabla nutricional).

Datos FICTICIOS de prueba (ILUSTRATIVO): no son el recetario ni la tabla reales. Los valores derivados
(kcal_por_porcion, fuentes_kcal y la tabla de escalado) se calculan aquí con fractions, de forma
independiente del código Go. Uso: python3 -I generar_testdata.py <carpeta testdata>
"""
import csv, json, sys
from decimal import Decimal, ROUND_HALF_UP
from fractions import Fraction as F

out = sys.argv[1]
G = [("frutas-verduras", "Frutas y verduras", "mercado", 10), ("carnes", "Carne y pollo", "mercado", 20),
     ("abarrotes", "Abarrotes", "cadena", 30), ("lacteos", "Lácteos y huevo", "cadena", 40)]
# id, nombre, grupo, unidad, gramos_por_unidad, [(presentación, nombre, contenido_milli)]
P = [("espagueti", "Espagueti", "abarrotes", "g", None, [("espagueti.la-moderna-250g", "Espagueti 250 g", 250000)]),
     ("atun", "Atún en agua", "abarrotes", "g", None, [("atun.lata-140g", "Atún lata 140 g", 140000)]),
     ("elote", "Granos de elote", "abarrotes", "g", None, [("elote.lata-220g", "Elote lata 220 g", 220000)]),
     ("aceite-oliva", "Aceite de oliva", "abarrotes", "ml", "0.913", [("aceite-oliva.500ml", "Aceite de oliva 500 ml", 500000)]),
     ("jitomate", "Jitomate saladet", "frutas-verduras", "g", None, [("jitomate.kilo", "Jitomate a granel", 1000000)]),
     ("huevo", "Huevo blanco", "lacteos", "pieza", "50", [("huevo.carton-12", "Huevo 12 piezas", 12000)]),
     ("tortilla", "Tortilla de maíz", "abarrotes", "g", None, [("tortilla.kilo", "Tortilla 1 kg", 1000000)]),
     ("arroz", "Arroz", "abarrotes", "g", None, [("arroz.1kg", "Arroz 1 kg", 1000000), ("arroz.750g", "Arroz 750 g", 750000)]),
     ("pollo", "Pierna de pollo", "carnes", "g", None, [("pollo.kilo", "Pollo a granel", 1000000)]),
     ("cebolla", "Cebolla blanca", "frutas-verduras", "g", None, [("cebolla.kilo", "Cebolla a granel", 1000000)]),
     ("sal", "Sal de mesa", "abarrotes", "g", None, [("sal.1kg", "Sal 1 kg", 1000000)]),
     ("frijol", "Frijol negro", "abarrotes", "g", None, [("frijol.1kg", "Frijol 1 kg", 1000000)]),
     ("aceite-soya", "Aceite de soya", "abarrotes", "ml", "0.921", [("aceite-soya.946ml", "Aceite de soya 946 ml", 946000)]),
     ("servilletas", "Servilletas", "abarrotes", "pieza", None, [("servilletas.paq-100", "Servilletas 100 piezas", 100000)])]
cat = {"version": "1.0", "fecha": "2026-10-11", "grupos": [{"id": a, "nombre": b, "canal": c, "orden": d} for a, b, c, d in G],
       "productos": [{"id": i, "nombre": n, "categoria": g, "unidad_base": u, "gramos_por_unidad": gp,
                      "presentacion_tipica": pr[0][0],
                      "presentaciones": [{"id": a, "nombre": b, "contenido_milli": c} for a, b, c in pr]}
                     for i, n, g, u, gp, pr in P]}
json.dump(cat, open(f"{out}/catalogo.json", "w"), ensure_ascii=False, indent=1)
unit = {i: u for i, _, _, u, _, _ in P}
gpu = {i: (F(Decimal(gp)) if gp else None) for i, _, _, _, gp, _ in P}
pres = {i: pr[0][0] for i, _, _, _, _, pr in P}
# producto, kcal_100g, fuente, id_fuente, nombre_fuente, estado
T = [("espagueti", "371", "USDA", "100001", "Pasta, dry, unenriched", "confirmado"),
     ("atun", "116", "OFF", "7501000000017", "Atun en agua", "confirmado"),
     ("elote", "81.5", "manual", "Etiqueta de prueba, lata 220 g", "", "confirmado"),
     ("aceite-oliva", "884", "USDA", "100002", "Oil, olive, salad or cooking", "confirmado"),
     ("jitomate", "18", "USDA", "100003", "Tomatoes, red, ripe, raw", "confirmado"),
     ("huevo", "143", "USDA", "100004", "Egg, whole, raw, fresh", "confirmado"),
     ("cebolla", "40", "USDA", "100005", "Onions, raw", "confirmado"),
     ("aceite-soya", "884", "USDA", "100006", "Oil, soybean", "confirmado"),
     ("tortilla", "218", "manual", "Tabla de prueba, página 1", "", "propuesto"),
     ("arroz", "365", "USDA", "100007", "Rice, white, long-grain, raw", "confirmado"),
     ("sal", "0", "USDA", "100008", "Salt, table", "confirmado"),
     ("pollo", "172.25", "USDA", "100009", "Chicken, drumstick, raw", "confirmado")]
with open(f"{out}/tabla.csv", "w", newline="") as f:
    w = csv.writer(f, lineterminator="\n")
    w.writerow("producto_id,kcal_100g,proteina_g,grasa_g,carbohidratos_g,fuente,id_fuente,nombre_fuente,fecha,estado,confirmo,fecha_confirmacion".split(","))
    for p, k, s, i, n, e in T:
        conf = ("axel", "2026-10-11") if e == "confirmado" else ("", "")
        w.writerow([p, k, "1.5" if p == "huevo" else "", "", "", s, i, n, "2026-10-11", e, *conf])
kcal = {p: F(Decimal(k)) for p, k, *_ in T}
src = {p: (s, i) for p, _, s, i, _, _ in T}
R = [("espagueti-con-atun", "Espagueti con atún y elote", ["lunch", "dinner"], [("espagueti", 250000), ("atun", 140000), ("elote", 110000), ("aceite-oliva", 15000)]),
     ("huevos-a-la-mexicana", "Huevos a la mexicana", ["breakfast", "dinner"], [("huevo", 4000), ("jitomate", 150000), ("cebolla", 40000), ("aceite-soya", 15000), ("tortilla", 200000)]),
     ("arroz-rojo", "Arroz rojo", ["lunch"], [("arroz", 200000), ("jitomate", 100000), ("cebolla", 30000), ("aceite-soya", 20000), ("sal", 3000)]),
     ("pollo-con-arroz", "Pollo con arroz", ["lunch", "dinner"], [("pollo", 500000), ("arroz", 150000), ("cebolla", 50000), ("sal", 4000)]),
     ("frijoles-de-la-olla", "Frijoles de la olla", ["lunch", "dinner", "snack"], [("frijol", 250000), ("cebolla", 40000), ("sal", 5000)])]
U = {"g": "g", "ml": "ml", "pieza": "pieza"}
recetas, escalado = [], []
for rid, nom, ts, ings in R:
    tot, falta = F(0), False
    for p, q in ings:
        if p not in kcal:
            falta = True
            continue
        g = F(q) if unit[p] == "g" else F(q) * gpu[p]
        tot += g * kcal[p] / 100  # milésimas de kcal: gramos_milli × kcal/100 g / 100
    kpp = None if falta else int((Decimal(tot.numerator) / Decimal(tot.denominator) / 2000).quantize(Decimal(1), ROUND_HALF_UP))
    recetas.append({"id": rid, "nombre": nom, "descripcion": "Receta de prueba.", "tiempos": ts, "porciones_base": 2,
                    "minutos_preparacion": 10, "minutos_coccion": 15, "origen": "claude",
                    "ingredientes": [{"producto_id": p, "cantidad_milli": q, "unidad": U[unit[p]], "presentacion_id": pres[p],
                                      "texto": "", "opcional": False} for p, q in ings],
                    "pasos": ["Paso único de prueba."], "kcal_por_porcion": kpp,
                    "fuentes_kcal": [{"producto_id": p, "fuente": src[p][0], "id_fuente": src[p][1]} for p, _ in ings if p in src]})
    for s in (4000, 1000, 3000, 6000):
        for p, q in ings:
            x = F(q) * s / 2000
            escalado.append({"recipe": rid, "servings_milli": s, "product": p, "num": x.numerator, "den": x.denominator})
json.dump({"version": "1.0", "fecha": "2026-10-11", "catalogo": "testdata/catalogo.json", "recetas": recetas},
          open(f"{out}/recetario.json", "w"), ensure_ascii=False, indent=1)
json.dump(escalado, open(f"{out}/escalado.json", "w"), indent=0)
print({r["id"]: r["kcal_por_porcion"] for r in recetas})
