#!/usr/bin/env python3
"""WF-3a · front-variantes: contraste y conformidad de lo medido en las historias de Storybook.

Entrada: medidas.json (lo escribe scripts/contraste-variantes.mjs: color de primer plano y fondo efectivo de cada
elemento renderizado, en hex, con los nombres de token que les corresponden, más tamaños de área táctil y problemas
de la corrida).
La razón de contraste se calcula con la fórmula de docs/design/tools/contraste_componentes.py (contraste(), Decimal,
truncada a 2 decimales): el módulo se IMPORTA, no se reescribe. El umbral es el de WCAG 2.x que usa DESIGN.md §10.5
(texto 4.5; texto grande e interfaz 3).

Clasificación de las fallas: las que coinciden con un hallazgo conocido de frontend/.storybook/hallazgos.json
(componente que no sigue DESIGN.md y cuya corrección pertenece a otro agente) se informan como HALLAZGO y no cambian
el código de salida; cualquier otra es FALLA y lo vuelve 1. Un hallazgo que ya no se reproduce se avisa como CERRADO.

Uso: python3 -I contraste-variantes.py <medidas.json> <carpeta de salida>
Salida: <carpeta>/contraste-variantes.csv, <carpeta>/resumen.json; código 0 si no hay FALLAS ni problemas, 1 si las hay.
"""
import csv
import importlib.util
import json
import sys
from collections import Counter, defaultdict
from decimal import Decimal
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[2]
FORMULA = RAIZ / "docs/design/tools/contraste_componentes.py"
CATALOGO = RAIZ / "docs/design/contraste-componentes.csv"
HALLAZGOS = RAIZ / "frontend/.storybook/hallazgos.json"

spec = importlib.util.spec_from_file_location("contraste_componentes", FORMULA)
cc = importlib.util.module_from_spec(spec)
spec.loader.exec_module(cc)  # solo define funciones y constantes; main() queda bajo __main__


def catalogo():
    """(hex_frente, hex_fondo, tema) -> razón de DESIGN.md §10.5 (docs/design/contraste-componentes.csv)."""
    if not CATALOGO.exists():
        return {}
    with CATALOGO.open(encoding="utf-8", newline="") as f:
        return {(r["hex_frente"], r["hex_fondo"], r["tema"]): r["razon"] for r in csv.DictReader(f)}


def evaluar(m):
    """Convierte una medición en fila con razón, umbral, cumple y motivo."""
    t = m["tipo"]
    base = {**m, "razon": "", "umbral_txt": "", "cumple": "true", "motivo": ""}
    if t == "foco-ausente":
        return {**base, "cumple": "false", "motivo": "foco sin anillo visible de 2 px"}
    if t == "foco-transicion":
        return {**base, "cumple": "false", "motivo": "el anillo de foco debe aparecer sin transición (§7.0, §4.5)"}
    if t == "deshabilitado-estilo":
        dif = m.get("diferencias", [])
        return {**base, "cumple": "false" if dif else "true", "motivo": "deshabilitado sin los tokens de §7.0: " + "; ".join(dif) if dif else ""}
    if t == "fuera-de-tokens":
        return {**base, "razon": str(cc.truncar2(cc.contraste(m["fg"], m["bg"]))), "umbral_txt": "n/a", "cumple": "false", "motivo": "color fuera de los tokens de C2"}
    razon = cc.contraste(m["fg"], m["bg"])  # misma fórmula que U2d
    umbral = Decimal(str(m["umbral"]))
    cumple, motivo = razon >= umbral, ""
    if not m["tokFg"] or not m["tokBg"]:  # un color que no es ningún token --ui-* de C2 no está permitido
        cumple, motivo = False, "color fuera de los tokens de C2"
    elif not cumple:
        motivo = "bajo umbral"
    if t == "foco" and (m.get("grosor", 0) < 2 or m.get("separacion", 0) != 2):
        cumple, motivo = False, "el anillo debe medir 2 px y separarse 2 px (§7.0)"
    return {**base, "razon": cc.truncar2(razon), "umbral_txt": str(umbral), "cumple": "true" if cumple else "false", "motivo": motivo}


def componente(titulo):
    return titulo.split(" › ")[0].replace("Componentes/", "")


def coincide(h, f):
    return (componente(f["titulo"]) in h["componentes"] and f["tipo"] in h["tipos"]
            and (not h.get("estados") or f["estado"] in h["estados"])
            and (not h.get("motivo") or h["motivo"] in f["motivo"])
            and (not h.get("elemento_contiene") or h["elemento_contiene"] in f["elemento"]))


def principal(entrada: Path, salida: Path, parcial: bool = False) -> int:
    datos = json.loads(entrada.read_text(encoding="utf-8"))
    if not datos["historias"]:
        print("Sin historias: nada que comprobar.")
        return 1
    cat = catalogo()
    conocidos = json.loads(HALLAZGOS.read_text(encoding="utf-8")) if HALLAZGOS.exists() else []
    filas = [evaluar(m) for m in datos["medidas"]]
    # Áreas táctiles: 44 px con puntero grueso (§4.6); con puntero fino se exige el mínimo de WCAG 2.2 (24 px).
    for t in datos["tamanos"]:
        minimo = 44 if t["contexto"] == "movil" else 24
        malo = t["ancho"] < minimo or t["alto"] < minimo
        filas.append({**t, "tipo": "tamano", "estado": "reposo", "fg": "", "bg": "", "tokFg": [], "tokBg": [], "razon": "", "umbral_txt": f"{minimo} px",
                      "cumple": "false" if malo else "true", "motivo": f"área táctil de {t['ancho']} x {t['alto']} px, mínimo {minimo} px" if malo else ""})
    # Cruce con el catálogo de U2d: mismo par de hex y tema -> misma razón (la fórmula es una sola).
    en_cat = difiere = 0
    for f in filas:
        if f["tipo"] in ("texto", "icono", "borde-fuera", "borde-dentro", "marca", "marca-serie", "foco"):
            k = (f["fg"], f["bg"], f["tema"])
            if k in cat:
                en_cat += 1
                difiere += cat[k] != f["razon"]
    # CSV deduplicado.
    cols = ["historia", "titulo", "tema", "contexto", "estado", "tipo", "elemento", "fg", "bg", "tokFg", "tokBg", "razon", "umbral_txt", "cumple", "motivo"]
    vistos, unicas = set(), []
    for f in filas:
        k = (f["historia"], f["tema"], f["contexto"], f["estado"], f["tipo"], f["elemento"], f["fg"], f["bg"])
        if k not in vistos:
            vistos.add(k)
            unicas.append(f)
    salida.mkdir(parents=True, exist_ok=True)
    with (salida / "contraste-variantes.csv").open("w", encoding="utf-8", newline="") as fh:
        w = csv.writer(fh, lineterminator="\n")
        w.writerow(cols)
        for f in unicas:
            w.writerow([("/".join(f[c]) if isinstance(f[c], list) else f[c]) for c in cols])

    fallas = [f for f in unicas if f["cumple"] == "false"]
    por_hallazgo = defaultdict(list)
    nuevas = []
    for f in fallas:
        h = next((h for h in conocidos if coincide(h, f)), None)
        (por_hallazgo[h["id"]] if h else nuevas).append(f)
    pares = {(f["tema"], f["tipo"], f["fg"], f["bg"]) for f in unicas if f["fg"]}
    por_tema, por_estado, por_tipo = (Counter(f[c] for f in unicas) for c in ("tema", "estado", "tipo"))
    medicion_contraste = [f for f in unicas if f["razon"] and f["umbral_txt"] not in ("n/a", "")]

    print(f"historias: {datos['historias']} · páginas renderizadas (historia x tema x contexto): {datos['paginas']}")
    print(f"mediciones únicas: {len(unicas)} (claro {por_tema['claro']}, oscuro {por_tema['oscuro']}) · con razón de contraste: {len(medicion_contraste)} · pares de color distintos: {len(pares)}")
    print("por estado: " + ", ".join(f"{k} {v}" for k, v in sorted(por_estado.items())))
    print("por tipo: " + ", ".join(f"{k} {v}" for k, v in sorted(por_tipo.items())))
    print(f"cruce con contraste-componentes.csv (U2d): {en_cat} mediciones con el mismo par de hex en C-001..C-158, {difiere} con razón distinta")
    print(f"FALLAS nuevas (no son hallazgos conocidos): {len(nuevas)}")
    grupo = defaultdict(list)
    for f in nuevas:
        grupo[(componente(f["titulo"]), f["estado"], f["tipo"], "/".join(f["tokFg"][:2]) or f["fg"], "/".join(f["tokBg"][:2]) or f["bg"], f["motivo"])].append(f)
    for k, g in sorted(grupo.items()):
        peor = min((Decimal(x["razon"]) for x in g if x["razon"]), default=None)
        print(f"  FALLA · {k[0]} · {k[1]} · {k[2]} · {k[3]} sobre {k[4]}{f' = {peor}' if peor is not None else ''} · {k[5]} · {len(g)} mediciones · p. ej. {'; '.join(sorted({x['elemento'] for x in g})[:2])}")
    print(f"hallazgos conocidos (frontend/.storybook/hallazgos.json): {len([h for h in conocidos if por_hallazgo[h['id']]])} de {len(conocidos)} se reproducen")
    for h in conocidos:
        g = por_hallazgo[h["id"]]
        estado = f"{len(g)} mediciones en {'/'.join(sorted({x['tema'] for x in g}))}" if g else ("no evaluable (corrida parcial)" if parcial else "CERRADO (ya no se reproduce)")
        print(f"  HALLAZGO {h['id']} · {', '.join(h['componentes'])} · {h['resumen']} · {estado}")
    print(f"problemas de la corrida (consola, hosts externos, desborde, render): {len(datos['problemas'])}")
    for p in datos["problemas"]:
        print(f"  PROBLEMA · {p.get('titulo', p.get('contexto'))} · {p.get('tema', '')} · {p.get('contexto', '')}: {p['problema']}")
    ok = not (nuevas or datos["problemas"] or difiere)
    (salida / "resumen.json").write_text(json.dumps({
        "historias": datos["historias"], "paginas": datos["paginas"], "mediciones": len(unicas), "con_razon": len(medicion_contraste),
        "pares": len(pares), "por_tema": por_tema, "por_estado": por_estado, "por_tipo": por_tipo, "en_catalogo": en_cat,
        "catalogo_difiere": difiere, "fallas_nuevas": len(nuevas), "hallazgos": {k: len(v) for k, v in por_hallazgo.items()},
        "problemas": len(datos["problemas"])}, ensure_ascii=False, indent=1), encoding="utf-8")
    print("RESULTADO:", "OK (sin fallas nuevas)" if ok else "CON FALLAS")
    return 0 if ok else 1


if __name__ == "__main__":
    pos = [a for a in sys.argv[1:] if not a.startswith("--")]
    if len(pos) != 2:
        print(__doc__)
        sys.exit(2)
    sys.exit(principal(Path(pos[0]), Path(pos[1]), "--parcial" in sys.argv))
