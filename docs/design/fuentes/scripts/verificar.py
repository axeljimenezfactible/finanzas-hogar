"""Verificación de U2b-1 (WF-1): fuentes, iconos y DESIGN.md secciones 1 a 6.

Uso: python3 -I docs/design/fuentes/scripts/verificar.py <raíz del repositorio>
Solo biblioteca estándar. Sale con código 1 si alguna comprobación falla.
- Fuentes: tamaño y sha256, firma woff2/ttf, ejes fvar, rasgos GSUB (tnum) y
  glifos de signos (+, U+2212) en cmap de los TTF de Inter.
- Licencias: OFL 1.1 (Inter) e ISC (Lucide).
- Iconos: símbolos del sprite, viewBox 24, sin trazo fijo; existen los nombres
  de la tabla «5.4 Mapa de iconos» de DESIGN.md.
- DESIGN.md: todo hex existe en C2-colores.json; toda ruta `color.*` existe;
  token seguido de hex coincide; contrastes de §2.8 recalculados (truncados a
  2 decimales, como contraste.csv); sin URLs de terceros; 10 encabezados.
"""
import hashlib, json, os, re, struct, sys

ROOT = sys.argv[1] if len(sys.argv) > 1 else "."
F = os.path.join(ROOT, "docs/design/fuentes")
DESIGN = os.path.join(ROOT, "docs/design/DESIGN.md")
C2 = os.path.join(ROOT, "docs/contratos/C2-colores.json")
fallas = []

def falla(msg):
    fallas.append(msg)
    print("FALLA:", msg)

def tablas_ttf(data):
    n = struct.unpack(">H", data[4:6])[0]
    t = {}
    for i in range(n):
        tag, _, off, ln = struct.unpack(">4sIII", data[12 + 16 * i:28 + 16 * i])
        t[tag.decode()] = (off, ln)
    return t

def ejes_fvar(data, t):
    off, _ = t["fvar"]
    _, _, arr, _, cnt, size, inst, _ = struct.unpack(">HHHHHHHH", data[off:off + 16])
    ejes = []
    for i in range(cnt):
        p = off + arr + i * size
        tag, mn, df, mx = struct.unpack(">4siii", data[p:p + 16])
        ejes.append((tag.decode(), mn / 65536, df / 65536, mx / 65536))
    return ejes, inst

def rasgos_gsub(data, t):
    off, _ = t["GSUB"]
    fl = struct.unpack(">H", data[off + 6:off + 8])[0]
    p = off + fl
    n = struct.unpack(">H", data[p:p + 2])[0]
    return sorted({data[p + 2 + 6 * i:p + 6 + 6 * i].decode() for i in range(n)})

def en_cmap(data, t, cps):
    off, _ = t["cmap"]
    n = struct.unpack(">H", data[off + 2:off + 4])[0]
    for i in range(n):
        pid, eid, so = struct.unpack(">HHI", data[off + 4 + 8 * i:off + 12 + 8 * i])
        s = off + so
        fmt = struct.unpack(">H", data[s:s + 2])[0]
        if fmt == 12:
            ng = struct.unpack(">I", data[s + 12:s + 16])[0]
            grupos = [struct.unpack(">III", data[s + 16 + 12 * g:s + 28 + 12 * g]) for g in range(ng)]
            return {cp: any(a <= cp <= b for a, b, _ in grupos) for cp in cps}
    return {cp: False for cp in cps}

print("== Fuentes ==")
for rel in ["inter/web/InterVariable.woff2", "inter/web/InterVariable-Italic.woff2",
            "inter/ttf/InterVariable.ttf", "inter/ttf/InterVariable-Italic.ttf",
            "inter/LICENSE.txt", "lucide/sprite.svg", "lucide/LICENSE", "lucide/package-1.55.0.json"]:
    p = os.path.join(F, rel)
    if not os.path.exists(p):
        falla(f"no existe {rel}")
        continue
    d = open(p, "rb").read()
    print(f"{rel}: {len(d)} B sha256 {hashlib.sha256(d).hexdigest()}")
    if rel.endswith(".woff2") and d[:4] != b"wOF2":
        falla(f"{rel} sin firma wOF2")
    if rel.endswith(".ttf"):
        if d[:4] != b"\x00\x01\x00\x00":
            falla(f"{rel} sin firma TrueType")
            continue
        t = tablas_ttf(d)
        ejes, inst = ejes_fvar(d, t)
        print("  ejes:", "; ".join(f"{a} {mn:g}–{mx:g} (inicial {df:g})" for a, mn, df, mx in ejes), f"· instancias {inst}")
        r = rasgos_gsub(d, t)
        print("  rasgos GSUB:", " ".join(r))
        if "tnum" not in r:
            falla(f"{rel} sin rasgo tnum")
        if not {"wght", "opsz"} <= {e[0] for e in ejes}:
            falla(f"{rel} sin ejes wght y opsz")
        g = en_cmap(d, t, [0x2B, 0x2212, 0x24, 0x2C, 0x2E])
        print("  glifos + − $ , . en cmap:", g)
        if not all(g.values()):
            falla(f"{rel} sin algún glifo de signo")

print("== Licencias ==")
lic_i = open(os.path.join(F, "inter/LICENSE.txt"), encoding="utf-8").read()
lic_l = open(os.path.join(F, "lucide/LICENSE"), encoding="utf-8").read()
print("Inter OFL 1.1:", "SIL Open Font License, Version 1.1" in lic_i)
print("Lucide ISC:", lic_l.startswith("ISC License"), "· MIT de Feather incluida:", "Cole Bemis" in lic_l)
if "SIL Open Font License, Version 1.1" not in lic_i:
    falla("licencia OFL ausente")
if not lic_l.startswith("ISC License"):
    falla("licencia ISC ausente")

print("== Iconos ==")
sp = open(os.path.join(F, "lucide/sprite.svg"), encoding="utf-8").read()
ids = re.findall(r'<symbol id="([^"]+)" viewBox="([^"]+)"', sp)
nombres = {i for i, _ in ids}
print(f"símbolos {len(ids)} · viewBox distintos {sorted({v for _, v in ids})} · stroke-width fijos {len(re.findall('stroke-width', sp))}")
if "@license lucide-static v1.55.0 - ISC" not in sp:
    falla("sprite sin cabecera de licencia")
txt = open(DESIGN, encoding="utf-8").read() if os.path.exists(DESIGN) else ""
m = re.search(r"### 5\.4 Mapa de iconos(.*?)(\n### |\n## )", txt, re.S)
usados = []
if m:
    for fila in m.group(1).splitlines():
        if fila.startswith("|") and not fila.startswith("|---"):
            celdas = fila.split("|")
            if len(celdas) > 2:
                usados += re.findall(r"`([a-z0-9]+(?:-[a-z0-9]+)*)`", celdas[2])
    faltan = [u for u in usados if u not in nombres]
    print(f"iconos citados en 5.4: {len(usados)} · faltan en el sprite: {faltan}")
    if faltan:
        falla(f"iconos inexistentes {faltan}")
else:
    falla("DESIGN.md sin tabla 5.4 Mapa de iconos")

print("== DESIGN.md ==")
c2 = json.load(open(C2, encoding="utf-8"))
tok = {}
def recorre(n, ruta):
    if isinstance(n, dict) and "$value" in n:
        tok[ruta] = n["$value"].upper()
    elif isinstance(n, dict):
        for k, v in n.items():
            if not k.startswith("$"):
                recorre(v, f"{ruta}.{k}" if ruta else k)
recorre(c2, "")
hexes = set(tok.values())
cit = [h.upper() for h in re.findall(r"#[0-9A-Fa-f]{6}\b", txt)]
ajenos = sorted({h for h in cit if h not in hexes})
print(f"tokens en C2 {len(tok)} · hex citados {len(cit)} ({len(set(cit))} distintos) · fuera de C2 {ajenos}")
if ajenos:
    falla(f"hex que no están en C2: {ajenos}")
rutas = [r for r in re.findall(r"`(color\.[a-z0-9.\-]+)`", txt)]
malas = sorted({r for r in rutas if r not in tok})
print(f"rutas de token citadas {len(rutas)} · inexistentes {malas}")
if malas:
    falla(f"rutas inexistentes: {malas}")
pares = re.findall(r"`(color\.[a-z0-9.\-]+)`[^`|\n]{0,8}`(#[0-9A-Fa-f]{6})`", txt)
dif = [(r, h) for r, h in pares if tok.get(r) != h.upper()]
print(f"pares token→hex {len(pares)} · no coinciden {dif}")
if dif:
    falla(f"token y hex no coinciden: {dif}")

def lum(h):
    c = [int(h[i:i + 2], 16) / 255 for i in (1, 3, 5)]
    c = [x / 12.92 if x <= 0.04045 else ((x + 0.055) / 1.055) ** 2.4 for x in c]
    return 0.2126 * c[0] + 0.7152 * c[1] + 0.0722 * c[2]

def ratio(a, b):
    la, lb = sorted([lum(a), lum(b)], reverse=True)
    return (la + 0.05) / (lb + 0.05)

m = re.search(r"### 2\.8 [^\n]*\n(.*?)(\n### |\n## )", txt, re.S)
n_ok = 0
if m:
    for fila in m.group(1).splitlines():
        r = re.findall(r"`(color\.[a-z0-9.\-]+)`", fila)
        v = re.search(r"\|\s*\*{0,2}(\d+\.\d\d)\*{0,2}\s*\|\s*[^|]*\|?\s*$", fila)
        if len(r) >= 2 and v:
            calc = int(ratio(tok[r[0]], tok[r[1]]) * 100) / 100
            ok = abs(calc - float(v.group(1))) < 0.005
            n_ok += ok
            print(f"  {r[0]} / {r[1]}: documento {v.group(1)} · recalculado {calc:.2f} {'ok' if ok else 'DIFIERE'}")
            if not ok:
                falla(f"contraste {r[0]}/{r[1]}")
    print(f"contrastes de 2.8 recalculados que coinciden: {n_ok}")
else:
    falla("DESIGN.md sin tabla 2.8")
terceros = re.findall(r"fonts\.googleapis|fonts\.gstatic|cdn\.|unpkg|jsdelivr", txt)
print("referencias a terceros:", len(terceros))
if terceros:
    falla("referencias a servicios de terceros")
h2 = re.findall(r"^## (.+)$", txt, re.M)
print("encabezados ##:", len(h2), "·", " | ".join(x[:28] for x in h2))
print("RESULTADO:", "OK" if not fallas else f"{len(fallas)} fallas")
sys.exit(1 if fallas else 0)
