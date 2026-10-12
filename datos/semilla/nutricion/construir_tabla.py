#!/usr/bin/env python3
"""Construye datos/semilla/nutricion/tabla.csv de forma reproducible.

Uso (desde cualquier carpeta):  python3 -I construir_tabla.py

Entradas (todas locales, sin red):
  raw/usda-foundation-2026-04-30/*.zip   USDA FoodData Central, Foundation Foods (CSV, 04/2026), CC0
  raw/usda-srlegacy-2018-04/*.zip        USDA FoodData Central, SR Legacy (CSV, 04/2018), CC0
  off_filas_usadas.tsv                   solo las filas de Open Food Facts que se usan (ODbL), ver extraer_off.py
  empates.tsv                            empate producto -> alimento de la fuente, escrito a mano, con confianza y nota
  ../productos.json                      catálogo (160 productos)

Salidas:
  tabla.csv           formato de docs/contratos/motores/S4-interfaz.md §9.1 (estado «propuesto» siempre)
  tabla_detalle.tsv   una fila por producto: conjunto de USDA, nutriente de energía usado, confianza y nota
  revisar.tsv         lista para Axel: empates con confianza media o baja y productos sin fuente

Reglas de lectura (ver PROCEDENCIA.md): kcal = nutriente 1008; si falta, 2048 (Atwater específico); si falta, 2047
(Atwater general). Cantidades por 100 g, redondeo ROUND_HALF_UP a 3 decimales, sin ceros a la derecha.
Un macro negativo (artefacto de «carbohidratos por diferencia» en algunos cortes de pollo) se deja vacío.
"""
import csv
import io
import json
import sys
import zipfile
from decimal import ROUND_HALF_UP, Decimal
from pathlib import Path

AQUI = Path(__file__).resolve().parent
FECHA = "2026-10-11"  # fecha de extracción (fija para que la salida sea idéntica en cada corrida)
ZIPS = {
    "Foundation 2026-04-30": AQUI / "raw/usda-foundation-2026-04-30/FoodData_Central_foundation_food_csv_2026-04-30.zip",
    "SR Legacy 2018-04": AQUI / "raw/usda-srlegacy-2018-04/FoodData_Central_sr_legacy_food_csv_2018-04.zip",
}
TIPOS_USDA = {"foundation_food", "sr_legacy_food"}
ENERGIA = ("1008", "2048", "2047")  # kcal; Atwater específico; Atwater general
MACROS = {"proteina": "1003", "grasa": "1004", "carbohidratos": "1005"}
CABECERA = ["producto_id", "kcal_100g", "proteina_g", "grasa_g", "carbohidratos_g", "fuente", "id_fuente",
            "nombre_fuente", "fecha", "estado", "confirmo", "fecha_confirmacion"]
FUENTES = {"USDA", "OFF", "manual", "SIN_FUENTE"}
CONFIANZAS = {"alta", "media", "baja"}


def fmt(valor):
    """Decimal -> texto con hasta 3 decimales (HALF_UP) y sin ceros sobrantes; negativo -> vacío."""
    if valor is None or valor < 0:
        return ""
    q = valor.quantize(Decimal("0.001"), rounding=ROUND_HALF_UP)
    s = format(q, "f")
    if "." in s:
        s = s.rstrip("0").rstrip(".")
    return s


def filas_zip(z, nombre):
    ruta = [n for n in z.namelist() if n.endswith("/" + nombre)][0]
    with z.open(ruta) as f:
        yield from csv.DictReader(io.TextIOWrapper(f, encoding="utf-8", newline=""))


def cargar_usda(ids):
    """ids: conjunto de fdc_id (texto). Devuelve {fdc_id: dict} desde Foundation y, si no está, SR Legacy."""
    res = {}
    pendientes = set(ids)
    for conjunto, ruta in ZIPS.items():
        if not pendientes:
            break
        z = zipfile.ZipFile(ruta)
        desc = {}
        for r in filas_zip(z, "food.csv"):
            if r["fdc_id"] in pendientes and r["data_type"] in TIPOS_USDA:
                desc[r["fdc_id"]] = r["description"]
        nut = {i: {} for i in desc}
        buscados = set(ENERGIA) | set(MACROS.values())
        for r in filas_zip(z, "food_nutrient.csv"):
            if r["fdc_id"] in nut and r["nutrient_id"] in buscados and r["amount"] != "":
                nut[r["fdc_id"]].setdefault(r["nutrient_id"], Decimal(r["amount"]))
        for i, d in desc.items():
            n = nut[i]
            usado = next((e for e in ENERGIA if e in n), None)
            res[i] = {
                "conjunto": conjunto,
                "nombre": d,
                "energia_id": usado,
                "kcal": n[usado] if usado else None,
                "proteina": n.get(MACROS["proteina"]),
                "grasa": n.get(MACROS["grasa"]),
                "carbohidratos": n.get(MACROS["carbohidratos"]),
            }
        pendientes -= set(desc)
    return res


def cargar_off():
    ruta = AQUI / "off_filas_usadas.tsv"
    res = {}
    if ruta.exists():
        with open(ruta, encoding="utf-8", newline="") as f:
            for r in csv.DictReader(f, delimiter="\t"):
                res[r["code"]] = r
    return res


def dec(texto):
    return Decimal(texto) if texto not in ("", None) else None


def main():
    catalogo = json.load(open(AQUI.parent / "productos.json", encoding="utf-8"))
    orden = [p["id"] for p in catalogo["productos"]]
    nombres = {p["id"]: p["nombre"] for p in catalogo["productos"]}

    empates = {}
    with open(AQUI / "empates.tsv", encoding="utf-8", newline="") as f:
        for n, r in enumerate(csv.reader(f, delimiter="\t"), 1):
            if not r or r[0].startswith("#"):
                continue
            if len(r) != 5:
                sys.exit(f"empates.tsv línea {n}: se esperaban 5 columnas y hay {len(r)}")
            pid, fuente, idf, conf, nota = [c.strip() for c in r]
            if pid in empates:
                sys.exit(f"empates.tsv línea {n}: producto repetido {pid}")
            if pid not in nombres:
                sys.exit(f"empates.tsv línea {n}: producto inexistente {pid}")
            if fuente not in FUENTES or conf not in CONFIANZAS:
                sys.exit(f"empates.tsv línea {n}: fuente o confianza inválida")
            empates[pid] = (fuente, idf, conf, nota)
    faltan = [p for p in orden if p not in empates]
    if faltan:
        sys.exit("productos sin renglón en empates.tsv: " + ", ".join(faltan))

    usda = cargar_usda({v[1] for v in empates.values() if v[0] == "USDA"})
    off = cargar_off()

    tabla, detalle, revisar = [], [], []
    conteo = {"USDA": 0, "OFF": 0, "manual": 0, "SIN_FUENTE": 0}
    conjuntos = {}
    for pid in orden:
        fuente, idf, conf, nota = empates[pid]
        conteo[fuente] += 1
        if fuente == "SIN_FUENTE":
            detalle.append([pid, nombres[pid], fuente, "", "", "", conf, nota])
            revisar.append([pid, nombres[pid], "sin fuente", "", "", "", "", nota])
            continue
        if fuente == "USDA":
            d = usda.get(idf)
            if d is None:
                sys.exit(f"{pid}: fdc_id {idf} no está en Foundation ni en SR Legacy")
            if d["kcal"] is None:
                sys.exit(f"{pid}: fdc_id {idf} no trae energía")
            kcal, prot, gra, car = d["kcal"], d["proteina"], d["grasa"], d["carbohidratos"]
            nombre_fuente, conjunto, energia = d["nombre"], d["conjunto"], d["energia_id"]
            conjuntos[conjunto] = conjuntos.get(conjunto, 0) + 1
        elif fuente == "OFF":
            r = off.get(idf)
            if r is None:
                sys.exit(f"{pid}: código {idf} no está en off_filas_usadas.tsv")
            kcal, prot, gra, car = dec(r["energy-kcal_100g"]), dec(r["proteins_100g"]), dec(r["fat_100g"]), dec(r["carbohydrates_100g"])
            nombre_fuente = (r["product_name"] + " (" + r["brands"] + ")") if r["brands"] else r["product_name"]
            conjunto, energia = "Open Food Facts 2026-10-11", "energy-kcal_100g"
            conjuntos[conjunto] = conjuntos.get(conjunto, 0) + 1
        else:
            sys.exit(f"{pid}: la fuente manual aún no se usa en esta tabla")
        tabla.append([pid, fmt(kcal), fmt(prot), fmt(gra), fmt(car), fuente, idf, nombre_fuente, FECHA, "propuesto", "", ""])
        detalle.append([pid, nombres[pid], fuente, idf, conjunto, energia, conf, nota])
        if conf != "alta":
            revisar.append([pid, nombres[pid], conf, fuente, idf, nombre_fuente, fmt(kcal), nota])

    with open(AQUI / "tabla.csv", "w", encoding="utf-8", newline="") as f:
        w = csv.writer(f, lineterminator="\n")
        w.writerow(CABECERA)
        w.writerows(tabla)
    with open(AQUI / "tabla_detalle.tsv", "w", encoding="utf-8", newline="") as f:
        w = csv.writer(f, delimiter="\t", lineterminator="\n")
        w.writerow(["producto_id", "nombre_es", "fuente", "id_fuente", "conjunto", "energia", "confianza", "nota"])
        w.writerows(detalle)
    with open(AQUI / "revisar.tsv", "w", encoding="utf-8", newline="") as f:
        w = csv.writer(f, delimiter="\t", lineterminator="\n")
        w.writerow(["producto_id", "nombre_es", "motivo", "fuente", "id_fuente", "nombre_fuente", "kcal_100g", "nota"])
        w.writerows(revisar)

    print("productos del catálogo:", len(orden))
    print("renglones en tabla.csv:", len(tabla))
    print("por fuente:", conteo)
    print("por conjunto:", conjuntos)
    conf = {}
    for d in detalle:
        conf[d[6]] = conf.get(d[6], 0) + 1
    print("por confianza:", conf)
    print("en revisar.tsv:", len(revisar))


if __name__ == "__main__":
    main()
