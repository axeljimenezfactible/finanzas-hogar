#!/usr/bin/env python3
"""Extrae de raw/off-2026-10-11/*.csv.gz (volcado completo de Open Food Facts, ODbL + DbCL) solo las filas usadas.

Lectura en flujo (streaming) del .gz, sin cargarlo entero en memoria. Filtra por país México (countries_tags
contiene en:mexico) y por los códigos de barras que empates.tsv cita con fuente OFF. Escribe
off_filas_usadas.tsv, que es la entrada de OFF de construir_tabla.py.

Uso: python3 -I extraer_off.py
"""
import csv
import gzip
from pathlib import Path

AQUI = Path(__file__).resolve().parent
GZ = AQUI / "raw/off-2026-10-11/en.openfoodfacts.org.products.csv.gz"
COLS = ["code", "product_name", "brands", "quantity", "countries_tags", "energy-kcal_100g", "proteins_100g",
        "fat_100g", "carbohydrates_100g", "last_modified_t"]

csv.field_size_limit(10**9)
codigos = set()
with open(AQUI / "empates.tsv", encoding="utf-8", newline="") as f:
    for r in csv.reader(f, delimiter="\t"):
        if r and not r[0].startswith("#") and r[1].strip() == "OFF":
            codigos.add(r[2].strip())

filas = {}
with gzip.open(GZ, "rt", encoding="utf-8", newline="") as f:
    lector = csv.reader(f, delimiter="\t", quoting=csv.QUOTE_NONE)
    cab = next(lector)
    ix = {c: i for i, c in enumerate(cab)}
    pos = [ix[c] for c in COLS]
    for fila in lector:
        if len(fila) <= max(pos) or fila[ix["code"]] not in codigos:
            continue
        if "en:mexico" in fila[ix["countries_tags"]]:
            filas[fila[ix["code"]]] = [fila[i] for i in pos]

with open(AQUI / "off_filas_usadas.tsv", "w", encoding="utf-8", newline="") as f:
    w = csv.writer(f, delimiter="\t", lineterminator="\n")
    w.writerow(COLS)
    for c in sorted(filas):
        w.writerow(filas[c])
faltan = codigos - set(filas)
print("códigos pedidos:", len(codigos), "encontrados:", len(filas), "faltan:", sorted(faltan))
