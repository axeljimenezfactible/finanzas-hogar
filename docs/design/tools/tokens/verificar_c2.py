#!/usr/bin/env python3
"""Verificaciones de C2 (criterios C2-02 a C2-10 de docs/criterios/C2.md).
C2-04 incluye: hex iguales a C2-colores.json (siguiendo alias), tokens citados en DESIGN.md existentes y
ningún grupo sin $value citado como token. C2-10 incluye: 0 colores con el mismo hex sin alias.
Uso: python3 -I docs/design/tools/tokens/verificar_c2.py   (codigo 0 = todo bien)"""
import collections, json, re, sys
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[4]
C = RAIZ / "docs/contratos"
J = json.loads((C / "C2-tokens.json").read_text(encoding="utf-8"))
COL = json.loads((C / "C2-colores.json").read_text(encoding="utf-8"))
CSS = (C / "tokens.css").read_text(encoding="utf-8")
TW = (C / "tailwind.theme.css").read_text(encoding="utf-8")
DESIGN = (RAIZ / "docs/design/DESIGN.md").read_text(encoding="utf-8")
malos = []


def falla(msg):
    malos.append(msg)
    print("  FALLA:", msg)


# ---- recorrido DTCG ----
hojas = {}  # ruta -> (token, tipo efectivo)


def recorre(n, ruta, tipo):
    t = n.get("$type", tipo)
    if "$value" in n:
        hojas[".".join(ruta)] = (n, t)
        return
    for k, v in n.items():
        if k == "$root":
            recorre(v, ruta, t)
        elif not k.startswith("$"):
            recorre(v, ruta + [k], t)


for k, v in J.items():
    if not k.startswith("$"):
        recorre(v, [k], None)


def destino(ref):
    ref = ref.strip("{}")
    if ref in hojas:
        return hojas[ref]
    return None


def resuelve(valor, vistos=()):
    if isinstance(valor, str) and valor.startswith("{"):
        r = valor.strip("{}")
        if r in vistos:
            raise RuntimeError("ciclo " + r)
        d = hojas[r][0]
        return resuelve(d["$value"], vistos + (r,))
    return valor


def refs_de(v):
    if isinstance(v, str):
        return re.findall(r"\{([^}]+)\}", v)
    if isinstance(v, list):
        return [r for x in v for r in refs_de(x)]
    if isinstance(v, dict):
        return [r for x in v.values() for r in refs_de(x)]
    return []


print("C2-02 · formato DTCG")
sin_tipo = [r for r, (n, t) in hojas.items() if not t]
sin_valor = [r for r, (n, t) in hojas.items() if n.get("$value") in (None, "")]
rotas = [(r, x) for r, (n, t) in hojas.items() for x in refs_de(n["$value"]) if x not in hojas]
sin_desc = [r for r, (n, t) in hojas.items() if r.split(".")[0] != "color" and not n.get("$description") and not r.startswith("ui.")]
TIPOS = {"color", "dimension", "fontFamily", "fontWeight", "duration", "cubicBezier", "shadow", "number", "strokeStyle"}
fuera = [r for r, (n, t) in hojas.items() if t not in TIPOS]
for etiqueta, lista in (("sin tipo efectivo", sin_tipo), ("sin $value", sin_valor), ("referencias sin destino", rotas), ("tipo fuera de DTCG", fuera)):
    if lista:
        falla(f"{etiqueta}: {lista[:5]}")
cuenta = collections.Counter(t for n, t in hojas.values())
print(f"  {len(hojas)} tokens; 0 sin tipo, 0 sin $value, 0 referencias rotas; por tipo: {dict(sorted(cuenta.items()))}")

print("C2-03 · cobertura por categoría y temas")
grupos = collections.Counter()
for r in hojas:
    grupos[r.split(".")[0]] += 1
print("  por grupo:", dict(grupos))
# cada token de color existe en los dos temas: pares light/dark completos en ui y overlay
pares = collections.defaultdict(set)
for r in hojas:
    p = r.split(".")
    if p[0] in ("ui", "overlay") and p[-1] in ("light", "dark"):
        pares[".".join(p[:-1])].add(p[-1])
incompletos = [k for k, v in pares.items() if v != {"light", "dark"}]
if incompletos:
    falla(f"semánticos sin los dos temas: {incompletos}")
print(f"  {len(pares)} semánticos de color (ui + overlay) con claro y oscuro; 0 incompletos")
esperado = {"brand": 11, "neutral": 11}
for e, n in esperado.items():
    got = len([r for r in hojas if r.startswith(f"color.{e}.")])
    if got != n + (0 if e != "brand" else 0):
        falla(f"escala {e}: {got} pasos, se esperaban {n}")
cat = {
    "color.brand/neutral 50-950": sum(1 for r in hojas if re.match(r"color\.(brand|neutral)\.\d+$", r)),
    "roles (8) + critical": sum(1 for r in hojas if re.match(r"color\.(role|critical)\.", r)),
    "personas y tiendas (8 × 2)": sum(1 for r in hojas if re.match(r"color\.(person|store)\.", r)),
    "series de gráficas (8 × 2)": sum(1 for r in hojas if r.startswith("color.chart.")),
    "tipografía": sum(1 for r in hojas if r.startswith("font.")),
    "espaciado": sum(1 for r in hojas if r.startswith("space.")),
    "radios": sum(1 for r in hojas if r.startswith("radius.")),
    "sombras + velo": sum(1 for r in hojas if r.startswith(("shadow.", "overlay."))),
    "movimiento": sum(1 for r in hojas if r.startswith("motion.")),
    "puntos de quiebre": sum(1 for r in hojas if r.startswith("breakpoint.")),
}
for k, v in cat.items():
    print(f"  {k}: {v}")
    if v == 0:
        falla("categoría vacía " + k)
for esp, ruta in ((("font.size", 6), None), (("space", 8), None), (("radius", 4), None), (("breakpoint", 5), None)):
    pre, n = esp
    got = len([r for r in hojas if r.startswith(pre + ".") and r.count(".") == 2 - (pre in ("space", "radius", "breakpoint"))])
for pre, n in (("font.size.", 6), ("space.", 8), ("radius.", 4), ("breakpoint.", 5), ("motion.duration.", 3), ("motion.easing.", 2), ("color.chart.", 16)):
    got = len([r for r in hojas if r.startswith(pre)])
    if got != n:
        falla(f"{pre}* : {got} tokens, se esperaban {n}")

print("C2-04 · C2 = C2-colores.json y DESIGN.md")


def hojas_color(n, ruta, out):
    for k, v in n.items():
        if k.startswith("$"):
            continue
        if "$value" in v:
            out[".".join(ruta + [k])] = v["$value"]
        else:
            hojas_color(v, ruta + [k], out)


orig = {}
hojas_color(COL["color"], ["color"], orig)
def hex_c2(r):
    """Hex final en C2-tokens.json, siguiendo alias (light/dark/surface-* de rol apuntan a su paso)."""
    return resuelve(hojas[r][0]["$value"]) if r in hojas else None


dif = [(r, v, hex_c2(r)) for r, v in orig.items() if hex_c2(r) != v]
extra = [r for r in hojas if r.startswith("color.") and r not in orig]
print(f"  {len(orig)} colores de C2-colores.json; diferencias de hex: {len(dif)}; colores de más en C2-tokens: {len(extra)}")
if dif or extra or len(orig) != 190:
    falla(f"colores distintos: {dif[:3]} {extra[:3]}")
# hex de las tablas de DESIGN.md
vistos = 0
bad = []
for nombre, hexv in re.findall(r"`((?:color)\.[A-Za-z0-9.\-]+)`\s*·\s*`(#[0-9A-Fa-f]{6})`", DESIGN):
    if nombre in orig:
        vistos += 1
        if orig[nombre].upper() != hexv.upper():
            bad.append((nombre, hexv, orig[nombre]))
print(f"  pares «token · hex» de DESIGN.md comparados: {vistos}; diferencias: {len(bad)}")
if bad:
    falla(f"hex de DESIGN.md distintos: {bad[:3]}")
citados = set(re.findall(r"`((?:color|ui|font|space|radius|shadow|overlay|border|focus|size|layout|icon|motion|breakpoint|viz)\.[A-Za-z0-9_.\-]*[A-Za-z0-9_])`", DESIGN))
todas_rutas = set(hojas)
grupos_rutas = {".".join(r.split(".")[:i]) for r in hojas for i in range(1, len(r.split(".")))}
ausentes = sorted(c for c in citados if c not in todas_rutas and c not in grupos_rutas and not re.search(r"\.N\.", c) and c != "color.chart.N.light")
print(f"  tokens citados en DESIGN.md: {len(citados)}; inexistentes en C2-tokens.json: {len(ausentes)} {ausentes[:6]}")
if ausentes:
    falla(f"tokens citados que no existen: {ausentes}")


def grupo_terminal(ruta):
    """Hijos de un grupo cuyos hijos son todos hojas con $value (y no es un par light/dark ni tiene $root)."""
    n = J
    for k in ruta.split("."):
        if not isinstance(n, dict) or k not in n:
            return None
        n = n[k]
    if not isinstance(n, dict) or "$value" in n or "$root" in n:
        return None
    hijos = [k for k in n if not k.startswith("$")]
    if hijos and set(hijos) != {"light", "dark"} and all(isinstance(n[k], dict) and "$value" in n[k] for k in hijos):
        return hijos
    return None


# un grupo de hojas citado como si fuera un token con valor (p. ej. layout.dialog.width en vez de .short/.long)
como_token = sorted(c for c in citados if c in grupos_rutas and grupo_terminal(c))
print(f"  grupos sin $value citados como token en DESIGN.md: {len(como_token)} {como_token[:6]}")
if como_token:
    falla(f"DESIGN.md cita grupos como tokens (citar sus hojas): {[(c, grupo_terminal(c)) for c in como_token]}")

print("C2-05 · contraste de los pares de componente básicos (la tabla completa de 158 pares es de U2d)")


def lum(h):
    h = h.lstrip("#")
    c = [int(h[i:i + 2], 16) / 255 for i in (0, 2, 4)]
    c = [x / 12.92 if x <= 0.03928 else ((x + 0.055) / 1.055) ** 2.4 for x in c]
    return 0.2126 * c[0] + 0.7152 * c[1] + 0.0722 * c[2]


def ratio(a, b):
    la, lb = sorted((lum(a), lum(b)), reverse=True)
    return int((la + 0.05) / (lb + 0.05) * 100) / 100


def ui(nombre, tema):
    return resuelve(hojas[f"ui.{nombre}.{tema}"][0]["$value"])


def sem(nombre, tema):
    return ui(nombre, tema)


pares_c = []
for n in ("text.default", "text.secondary", "text.link", "icon.default"):
    for f in ("surface.page", "surface.card", "surface.raised"):
        pares_c.append((n, f, 4.5))
pares_c += [("text.on-primary", "action.primary.bg", 4.5), ("text.on-destructive", "action.destructive.bg", 4.5), ("text.on-critical", "action.critical.bg", 4.5),
            ("text.on-primary", "action.primary.bg-hover", 4.5), ("text.on-primary", "action.primary.bg-pressed", 4.5),
            ("text.on-destructive", "action.destructive.bg-hover", 4.5), ("text.on-critical", "action.critical.bg-hover", 4.5),
            ("border.control", "surface.field", 3.0), ("border.control", "surface.page", 3.0), ("border.control", "surface.card", 3.0), ("border.focus", "surface.page", 3.0),
            ("border.focus", "surface.card", 3.0), ("border.error", "surface.field", 3.0), ("text.default", "surface.selected", 4.5), ("text.disabled", "surface.disabled", 3.0)]
for r in ("income", "expense", "savings", "debt", "warning", "error", "success", "info"):
    pares_c += [(f"role.{r}", "surface.page", 4.5), (f"role.{r}", "surface.card", 4.5), (f"role.{r}", f"role.{r}-surface", 4.5)]
pares_c += [("critical", "surface.page", 4.5), ("critical", "critical-surface", 4.5)]
for p_ in ("axel", "ximena", "karen", "shared"):
    pares_c += [(f"person.{p_}", "surface.page", 3.0), (f"person.{p_}", "surface.card", 3.0)]
for s_ in ("walmart", "bodega-aurrera", "sams-club", "market"):
    pares_c += [(f"store.{s_}", "surface.page", 3.0), (f"store.{s_}", "surface.card", 3.0)]
for i in range(1, 9):
    pares_c += [(f"chart.{i}", "surface.card", 3.0), (f"chart.{i}", "surface.page", 3.0)]
bajo = []
total = 0
for a, b, u in pares_c:
    for tema in ("light", "dark"):
        total += 1
        r = ratio(sem(a, tema), sem(b, tema))
        if r < u:
            bajo.append((a, b, tema, r, u))
print(f"  {total} combinaciones (par × tema) calculadas desde C2-tokens.json; bajo umbral: {len(bajo)}")
for x in bajo:
    print("   bajo umbral:", x)
if bajo:
    falla(f"{len(bajo)} combinaciones bajo umbral (ver arriba; se reporta, no se cambia ningún token)")

print("C2-06 · tokens.css")
claro = re.search(r":root \{(.*?)\n\}", CSS, re.S).group(1)
m_osc = re.search(r"@media \(prefers-color-scheme: dark\) \{\s*:root:not\(\[data-theme=\"light\"\]\) \{(.*?)\n  \}\n\}", CSS, re.S).group(1)
a_osc = re.search(r":root\[data-theme=\"dark\"\] \{(.*?)\n\}", CSS, re.S).group(1)
vn = lambda b: set(re.findall(r"(--ui-[a-z0-9\-]+):", b))
v_claro, v_m, v_a = vn(claro), vn(m_osc), vn(a_osc)
print(f"  --ui-* claro {len(v_claro)} · oscuro (media) {len(v_m)} · oscuro (data-theme) {len(v_a)}; diferencias claro/oscuro: {len(v_claro ^ v_m)} y {len(v_claro ^ v_a)}")
if v_claro != v_m or v_claro != v_a:
    falla("conjuntos --ui-* distintos entre temas")
if re.findall(r"(--ui-[a-z0-9\-]+): ([^;]+);", m_osc) != re.findall(r"(--ui-[a-z0-9\-]+): ([^;]+);", a_osc):
    falla("el oscuro de @media y el de data-theme difieren")
for pat, nombre in ((r"prefers-color-scheme: dark", "prefers-color-scheme: dark"), (r"data-theme", "data-theme"), (r"prefers-reduced-motion", "prefers-reduced-motion")):
    n = len(re.findall(pat, CSS))
    print(f"  {nombre}: {n}")
    if n < 1:
        falla("falta " + nombre)
prefijos = {p: len(re.findall(r"^\s*--" + p, CSS, re.M)) for p in ("color-", "space-", "radius-", "shadow-", "duration-", "font-")}
print("  variables por prefijo:", prefijos)
if 0 in prefijos.values():
    falla("prefijo sin variables")
sueltos = re.findall(r"#[0-9A-Fa-f]{3,8}\b", re.sub(r"--color-[a-z0-9\-]+:\s*#[0-9A-Fa-f]{6};", "", CSS))
if sueltos:
    falla(f"colores sueltos fuera de los primitivos: {sueltos[:5]}")

print("C2-08 · autoalojado")
url = re.findall(r"https?://[^\s\"')]+", CSS + TW + (C / "C2-tokens.json").read_text(encoding="utf-8"))
print(f"  URL a terceros en tokens.css, tailwind.theme.css y C2-tokens.json: {len(url)}")
if url:
    falla(f"URL encontradas: {url[:3]}")
for src in re.findall(r'url\("([^"]+)"\)', CSS):
    ok = (C / src).resolve().exists()
    print(f"  fuente {src}: {'existe' if ok else 'NO EXISTE'}")
    if not ok:
        falla("fuente ausente " + src)
lic = (C / hojas["font.family.sans"][0]["$extensions"]["finanzas"]["license"]).resolve()
print(f"  licencia {lic.relative_to(RAIZ)}: {'existe' if lic.exists() else 'NO EXISTE'}")
if not lic.exists():
    falla("licencia ausente")

print("C2-10 · huérfanos, duplicados y nombres")
sin_cita = sorted(r for r in hojas if not r.startswith(("color.", "ui.")) and not any(r == c or r.startswith(c + ".") or c.startswith(r + ".") or re.sub(r"\.\d+\.dash$", ".N.dash", r) == c for c in citados))
print(f"  tokens no de color/ui sin cita literal en DESIGN.md: {len(sin_cita)} {sin_cita[:8]}")
dup = collections.defaultdict(list)
for r, (n, t) in hojas.items():
    if r.startswith(("color.", "ui.")) or r.startswith("viz."):
        continue
    v = n["$value"]
    if isinstance(v, str) and v.startswith("{"):
        continue
    dup[(t, json.dumps(v, sort_keys=True))].append(r)
dups = {k: v for k, v in dup.items() if len(v) > 1}
print(f"  grupos de valor igual sin alias entre tokens que no son de color: {len(dups)} (informativo: misma medida con distinta idea de diseño, Supuesto 12 de U2c)")
for k, v in sorted(dups.items(), key=lambda x: x[0][1]):
    print(f"   {k[0]} {k[1]}: {', '.join(v)}")
# Colores (C2-10): todo hex literal es único; lo que repite un valor es alias. Los ui.* y overlay.* son siempre referencias.
dupc = collections.defaultdict(list)
literales_sem = []
for r, (n, t) in hojas.items():
    if t != "color":
        continue
    v = n["$value"]
    es_ref = isinstance(v, str) and v.startswith("{")
    if r.startswith(("ui.", "overlay.")) and not es_ref:
        literales_sem.append(r)
    if not es_ref:
        dupc[v.upper()].append(r)
dupsc = {k: v for k, v in dupc.items() if len(v) > 1}
alias_col = sum(1 for r, (n, t) in hojas.items() if r.startswith("color.") and str(n["$value"]).startswith("{"))
print(f"  colores: {sum(1 for r in hojas if r.startswith('color.'))} primitivos ({alias_col} alias de su paso); hex literales repetidos sin alias: {len(dupsc)}; semánticos con valor literal: {len(literales_sem)}")
for k, v in sorted(dupsc.items()):
    print(f"   {k}: {', '.join(v)}")
if dupsc or literales_sem:
    falla(f"colores con el mismo valor sin alias: {list(dupsc.items())[:3]}; semánticos literales: {literales_sem[:3]}")
noms = [r for r in hojas if re.search(r"[A-Z_ ]", r)]
if noms:
    falla(f"nombres con mayúsculas o espacios: {noms[:3]}")
print("RESULTADO:", "OK" if not malos else f"{len(malos)} FALLAS")
sys.exit(1 if malos else 0)
