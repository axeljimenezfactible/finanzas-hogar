#!/usr/bin/env python3 -I
"""Genera docs/design/sistema/project/tokens.css desde tokens.json.

Formato «tokens.css as compiled» (referencia-tipo/format.md):
  :root, [data-theme="<primero>"] { --<color>; --<shadow> }
  [data-theme="<id>"] por cada tema siguiente (todos los colores y sombras re-declarados)
  :root { --<space>; --<radius>; --<otras familias>; --font-<clave> }
  .<estilo> por estilo tipográfico
  @font-face por type.fonts[]
Uso: python3 -I tokens_css.py [--check]   (--check regenera en memoria y verifica)
"""
import json, re, sys
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[2] / "sistema" / "project"
JSON = RAIZ / "tokens.json"
CSS = RAIZ / "tokens.css"
ALIAS = re.compile(r"^\{([A-Za-z0-9._-]+)\}$")
FAMILIAS_TOKENS = ["spacing", "radius", "fontWeight", "border", "focus", "size",
                   "layout", "icon", "motion", "breakpoint", "viz"]


def esc(nombre):
    return nombre.replace(".", "\\.")


def val(v):
    m = ALIAS.match(v) if isinstance(v, str) else None
    return f"var(--{esc(m.group(1))})" if m else v


def bloque(selector, decls):
    return selector + " {\n" + "".join(f"  --{esc(n)}: {v};\n" for n, v in decls) + "}\n"


def generar(t):
    temas = [x["id"] for x in t["color"]["themes"]]
    colores = t["color"]["tokens"]
    sombras = t["shadow"]["tokens"]
    out = ["/* Finanzas Hogar — generated from tokens.json */\n"]
    for i, tid in enumerate(temas):
        sel = f':root, [data-theme="{tid}"]' if i == 0 else f'[data-theme="{tid}"]'
        decls = [(c["name"], val(c["value"][tid])) for c in colores]
        decls += [(s["name"], val(s["value"][tid])) for s in sombras]
        out.append(bloque(sel, decls))
    decls = []
    for fam in FAMILIAS_TOKENS:
        if fam in t:
            decls += [(x["name"], val(x["value"])) for x in t[fam]["tokens"]]
    decls += [(f"font-{k}", v) for k, v in t["type"]["families"].items()]
    out.append(bloque(":root", decls))
    for g in t["type"]["groups"]:
        fam = g["family"]
        for s in g["styles"]:
            out.append(f".{esc(s['name'])} {{\n"
                       f"  font-family: var(--font-{fam});\n"
                       f"  font-size: {s['fontSize']};\n"
                       f"  line-height: {s['lineHeight']};\n"
                       f"  font-weight: {s['fontWeight']};\n"
                       f"  letter-spacing: {s['letterSpacing']};\n}}\n")
    for f in t["type"]["fonts"]:
        out.append("@font-face {\n"
                   f'  font-family: "{f["family"]}";\n'
                   f'  src: url("{f["file"]}") format("woff2");\n'
                   f'  font-weight: {f["weight"]};\n'
                   f'  font-style: {f["style"]};\n}}\n')
    return "\n".join(out)


def verificar(t, css):
    fallos = []
    if css.count("{") != css.count("}"):
        fallos.append("llaves desbalanceadas")
    if not css.startswith("/* Finanzas Hogar — generated from tokens.json */\n"):
        fallos.append("línea 1 incorrecta")
    cuerpo = re.sub(r"/\*.*?\*/", "", css, flags=re.S)
    profundidad = 0
    bloques = []  # (selector, texto)
    for m in re.finditer(r"([^{}]+)\{([^{}]*)\}", cuerpo):
        bloques.append((m.group(1).strip(), m.group(2)))
    if re.sub(r"([^{}]+)\{([^{}]*)\}", "", cuerpo).strip():
        fallos.append("texto fuera de bloques / anidamiento")
    sels = [b[0] for b in bloques]
    esperado = [':root, [data-theme="light"]', '[data-theme="dark"]', ':root']
    if sels[:3] != esperado:
        fallos.append(f"selectores iniciales {sels[:3]}")
    tema = {}
    for sel, txt in bloques[:3]:
        tema[sel] = dict(re.findall(r"--([^:\s]+):\s*([^;]+);", txt))
    claro, oscuro, raiz = (tema[s] for s in esperado)
    esc_n = lambda n: esc(n)
    for c in t["color"]["tokens"]:
        n = esc_n(c["name"])
        if n not in claro or n not in oscuro:
            fallos.append(f"color {n} falta en un tema")
        elif claro[n] != val(c["value"]["light"]) or oscuro[n] != val(c["value"]["dark"]):
            fallos.append(f"valor distinto {n}")
    for s in t["shadow"]["tokens"]:
        n = esc_n(s["name"])
        if n not in claro or n not in oscuro:
            fallos.append(f"sombra {n} falta")
    clases = set(re.findall(r"^\.([^\s{]+) \{", css, flags=re.M))
    for g in t["type"]["groups"]:
        for s in g["styles"]:
            if esc(s["name"]) not in clases:
                fallos.append(f"clase .{s['name']} falta")
    n_ff = css.count("@font-face")
    if n_ff != len(t["type"]["fonts"]):
        fallos.append("@font-face distinto de type.fonts")
    for f in t["type"]["fonts"]:
        if f'url("{f["file"]}") format("woff2")' not in css:
            fallos.append(f"src de {f['file']}")
    # conteo de variables = tokens
    otras = sum(len(t[f]["tokens"]) for f in FAMILIAS_TOKENS if f in t) + len(t["type"]["families"])
    esp_claro = len(t["color"]["tokens"]) + len(t["shadow"]["tokens"])
    if len(claro) != esp_claro or len(oscuro) != esp_claro:
        fallos.append(f"vars tema {len(claro)}/{len(oscuro)} != {esp_claro}")
    if len(raiz) != otras:
        fallos.append(f"vars :root {len(raiz)} != {otras}")
    # todo var(--x) usado existe
    existentes = set(claro) | set(oscuro) | set(raiz)
    for ref in set(re.findall(r"var\(--([^)]+)\)", css)):
        if ref not in existentes:
            fallos.append(f"var(--{ref}) no existe")
    resumen = dict(colores=len(t["color"]["tokens"]), sombras=len(t["shadow"]["tokens"]),
                   vars_tema=len(claro), vars_raiz=len(raiz), clases=len(clases),
                   font_face=n_ff, bytes=len(css.encode("utf-8")))
    return fallos, resumen


def main():
    t = json.loads(JSON.read_text(encoding="utf-8"))
    css = generar(t)
    if "--check" not in sys.argv:
        CSS.write_text(css, encoding="utf-8")
    existente = CSS.read_text(encoding="utf-8") if CSS.exists() else ""
    fallos, resumen = verificar(t, css)
    if existente != css:
        fallos.append("tokens.css en disco difiere del generado")
    print(resumen)
    for f in fallos:
        print("FALLO:", f)
    sys.exit(1 if fallos else 0)


if __name__ == "__main__":
    main()
