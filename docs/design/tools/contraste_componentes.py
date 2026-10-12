#!/usr/bin/env python3
"""U2d · Contraste de componentes (WCAG 2.x) para Finanzas Hogar.

Calcula el contraste de TODOS los pares componente x estado x tema que
DESIGN.md §10.5 enumera (158 pares x 2 temas = 316 combinaciones) y escribe
docs/design/contraste-componentes.csv.

Fuentes (solo lectura):
  - docs/contratos/C2-tokens.json: UNICA fuente de la resolucion y de los hex (C2-05). Cada lado
    del «Par» de §10.5 se resuelve por las referencias de C2: `ui.<x>` -> ui.<x>.light/.dark;
    «rol X» -> ui.role.X; «superficie de rol X» -> ui.role.X-surface; «persona/tienda/serie» ->
    ui.person|store|chart.*; «acción crítica» -> ui.critical. Se siguen los alias ui.* hasta el
    primer primitivo color.* (columna token_*) y de ahi hasta el hex (columna hex_*).
  - docs/design/DESIGN.md §10.5: solo la lista de pares (id, componente, estado, par, tipo).
  - Cruces secundarios (avisan y dan codigo 1 si no coinciden; no alimentan el calculo):
    columnas claro/oscuro de §10.5, tabla §7.1 y docs/contratos/C2-colores.json.
  - docs/design/contraste.csv (202 pares de colorimetria): solo para la columna de cruce,
    no se importa codigo de colorimetria.py.

Uso (desde cualquier carpeta):
  python3 -I docs/design/tools/contraste_componentes.py            # escribe el CSV
  python3 -I docs/design/tools/contraste_componentes.py --check    # el CSV en disco debe ser identico
  python3 -I docs/design/tools/contraste_componentes.py --independiente   # 10 razones con segundo calculo

Codigo de salida: 0 si todo cumple; 1 si hay pares bajo umbral o falla una
comprobacion; 2 si el DESIGN.md no tiene la forma esperada.
Solo biblioteca estandar. No cambia ningun token: si un par no cumple, se reporta.
"""
import csv
import io
import json
import math
import re
import sys
from decimal import Decimal, ROUND_DOWN, getcontext
from pathlib import Path

getcontext().prec = 60  # precision de sobra: el truncado a 2 decimales no depende de redondeos de coma flotante

RAIZ = Path(__file__).resolve().parents[3]
DESIGN = RAIZ / "docs/design/DESIGN.md"
TOKENS = RAIZ / "docs/contratos/C2-tokens.json"
COLORES = RAIZ / "docs/contratos/C2-colores.json"
CRUCE = RAIZ / "docs/design/contraste.csv"
SALIDA = RAIZ / "docs/design/contraste-componentes.csv"

UMBRAL_TEXTO = Decimal("4.5")
UMBRAL_INTERFAZ = Decimal("3.0")
ESPERADOS = {"pares": 158, "combinaciones": 316, "con_umbral": 304, "decorativas": 12, "no_en_csv": 50}
COLUMNAS = ["id", "componente", "estado", "tema", "par", "token_frente", "token_fondo",
            "hex_frente", "hex_fondo", "razon", "tipo", "umbral", "cumple", "cruce_contraste_csv"]


# ---------------------------------------------------------------- color (Decimal)
def lineal(canal8: int) -> Decimal:
    c = Decimal(canal8) / Decimal(255)
    if c <= Decimal("0.04045"):
        return c / Decimal("12.92")
    return ((c + Decimal("0.055")) / Decimal("1.055")) ** Decimal("2.4")


def luminancia(hexa: str) -> Decimal:
    r, g, b = (int(hexa[i:i + 2], 16) for i in (1, 3, 5))
    return Decimal("0.2126") * lineal(r) + Decimal("0.7152") * lineal(g) + Decimal("0.0722") * lineal(b)


def contraste(h1: str, h2: str) -> Decimal:
    a, b = luminancia(h1), luminancia(h2)
    alto, bajo = (a, b) if a >= b else (b, a)
    return (alto + Decimal("0.05")) / (bajo + Decimal("0.05"))


def truncar2(x: Decimal) -> str:
    return str(x.quantize(Decimal("0.01"), rounding=ROUND_DOWN))


# ---------------------------------------------------------------- C2-tokens.json (unica fuente)
def cargar_tokens():
    """Ruta -> $value de cada hoja de C2-tokens.json (DTCG; $root cuenta como el grupo)."""
    raiz = json.loads(TOKENS.read_text(encoding="utf-8"))
    hojas = {}

    def recorre(n, ruta):
        if "$value" in n:
            hojas[".".join(ruta)] = n["$value"]
            return
        for k, v in n.items():
            if k == "$root":
                recorre(v, ruta)
            elif not k.startswith("$") and isinstance(v, dict):
                recorre(v, ruta + [k])
    for k, v in raiz.items():
        if not k.startswith("$") and isinstance(v, dict):
            recorre(v, [k])
    return hojas


def referencia(valor):
    m = re.fullmatch(r"\{([^{}]+)\}", valor) if isinstance(valor, str) else None
    return m.group(1) if m else None


def primitivo(hojas, token: str) -> str:
    """Sigue las referencias ui.* hasta el primer primitivo color.* (el nombre que va en el CSV)."""
    t, vistos = token, set()
    while t.startswith("ui."):
        if t in vistos or t not in hojas:
            raise ValueError(f"{token}: ciclo o hoja inexistente en {t}")
        vistos.add(t)
        r = referencia(hojas[t])
        if r is None:
            raise ValueError(f"{t}: un semántico debe ser referencia y vale {hojas[t]}")
        t = r
    if not t.startswith("color.") or t not in hojas:
        raise ValueError(f"{token}: no llega a un primitivo color.* ({t})")
    return t


def hex_de(hojas, token: str) -> str:
    """Hex final de una hoja de color, siguiendo alias (p. ej. color.role.income.light -> .700)."""
    t, vistos = token, set()
    while (r := referencia(hojas[t])) is not None:
        if t in vistos or r not in hojas:
            raise ValueError(f"{token}: ciclo o referencia rota en {t}")
        vistos.add(t)
        t = r
    v = hojas[t]
    if not re.fullmatch(r"#[0-9A-Fa-f]{6}", v):
        raise ValueError(f"{token}: valor no es #RRGGBB: {v}")
    return v.upper()


def cargar_colores():
    """Cruce secundario: ruta -> hex de C2-colores.json (la colorimetria)."""
    out = {}

    def recorre(n, ruta):
        for k, v in n.items():
            if k.startswith("$") or not isinstance(v, dict):
                continue
            if "$value" in v:
                out[".".join(ruta + [k])] = v["$value"].upper()
            else:
                recorre(v, ruta + [k])
    recorre(json.loads(COLORES.read_text(encoding="utf-8"))["color"], ["color"])
    return out


# ---------------------------------------------------------------- DESIGN.md
def celdas(linea: str):
    return [c.strip() for c in linea.strip().strip("|").split("|")]


def sin_tildes(s: str) -> str:
    return s.replace("`", "")


def leer_7_1():
    """ui.* -> (token claro, hex claro, token oscuro, hex oscuro) desde §7.1."""
    texto = DESIGN.read_text(encoding="utf-8").split("\n")
    ini = next(i for i, l in enumerate(texto) if l.startswith("### 7.1 "))
    fin = next(i for i, l in enumerate(texto) if i > ini and l.startswith("### 7.2 "))
    mapa = {}
    for l in texto[ini:fin]:
        if not l.startswith("| `ui."):
            continue
        c = celdas(l)
        ui = c[0].strip("`")
        m1 = re.fullmatch(r"`(color\.[\w.\-]+)` · `(#[0-9A-Fa-f]{6})`", c[1])
        m2 = re.fullmatch(r"`(color\.[\w.\-]+)` · `(#[0-9A-Fa-f]{6})`", c[2])
        if not (m1 and m2):
            raise ValueError(f"§7.1: fila ilegible: {l[:80]}")
        mapa[ui] = (m1.group(1), m1.group(2).upper(), m2.group(1), m2.group(2).upper())
    return mapa


def leer_10_5():
    texto = DESIGN.read_text(encoding="utf-8").split("\n")
    ini = next(i for i, l in enumerate(texto) if l.startswith("### 10.5 "))
    fin = next(i for i, l in enumerate(texto) if i > ini and l.startswith("### 10.6 "))
    filas = []
    for l in texto[ini:fin]:
        if not re.match(r"\| C-\d{3} \|", l):
            continue
        c = celdas(l)
        if len(c) != 8:
            raise ValueError(f"§10.5: {len(c)} columnas en {c[0]}")
        cl = [t.strip("`") for t in c[5].split(" / ")]
        os_ = [t.strip("`") for t in c[6].split(" / ")]
        if len(cl) != 2 or len(os_) != 2:
            raise ValueError(f"§10.5: par ilegible en {c[0]}")
        filas.append({"id": c[0], "componente": c[1], "estado": c[2], "par": sin_tildes(c[3]),
                      "par_crudo": c[3], "tipo": c[4], "claro": cl, "oscuro": os_,
                      "en_csv_dice": c[7]})
    return filas


def semantico_de(txt: str):
    """Nombre del semántico de C2-tokens.json (sin .light/.dark) para un lado del «Par» de §10.5."""
    txt = txt.strip()
    m = re.fullmatch(r"`(ui\.[\w.\-]+)`", txt)
    if m:
        return m.group(1)
    m = re.fullmatch(r"superficie de rol (\w+)", txt)
    if m:
        return f"ui.role.{m.group(1)}-surface"
    m = re.fullmatch(r"rol (\w+)", txt)
    if m:
        return f"ui.role.{m.group(1)}"
    m = re.fullmatch(r"persona ([\w\-]+)", txt)
    if m:
        return f"ui.person.{m.group(1)}"
    m = re.fullmatch(r"tienda ([\w\-]+)", txt)
    if m:
        return f"ui.store.{m.group(1)}"
    m = re.fullmatch(r"serie (\d)", txt)
    if m:
        return f"ui.chart.{m.group(1)}"
    if txt == "acción crítica":
        return "ui.critical"
    if txt == "superficie de acción crítica":
        return "ui.critical-surface"
    return None


def umbral_de(tipo: str):
    if tipo == "decorativo":
        return None
    if tipo == "texto 4.5":
        return UMBRAL_TEXTO
    if tipo in ("interfaz 3", "texto grande 3"):
        return UMBRAL_INTERFAZ
    raise ValueError(f"tipo desconocido: {tipo}")


def cargar_cruce():
    d = {}
    with CRUCE.open(encoding="utf-8", newline="") as f:
        for r in csv.DictReader(f):
            d[(r["primer_plano"], r["fondo"], r["tema"])] = r["ratio"]
    return d


# ---------------------------------------------------------------- generacion
def generar():
    hojas = cargar_tokens()
    hex_col = cargar_colores()
    ui7 = leer_7_1()
    filas = leer_10_5()
    cruce = cargar_cruce()
    avisos = []  # cruces secundarios que no coinciden con C2-tokens.json (hallazgos, no se corrige nada)

    # Cruce secundario 1: §7.1 contra C2-tokens.json (token primitivo y hex por tema)
    for ui, (tc, hc, to, ho) in ui7.items():
        for tema, t, h in (("light", tc, hc), ("dark", to, ho)):
            hoja_ui = f"{ui}.{tema}"
            if hoja_ui not in hojas:
                avisos.append(f"§7.1 {ui}: no existe {hoja_ui} en C2-tokens.json")
                continue
            prim = primitivo(hojas, hoja_ui)
            if prim != t or hex_de(hojas, prim) != h:
                avisos.append(f"§7.1 {hoja_ui}: dice {t} {h} y C2-tokens.json resuelve {prim} {hex_de(hojas, prim)}")

    salida = []
    usados = set()
    for f in filas:
        izq, der = f["par_crudo"].split(" sobre ", 1)
        sem_f, sem_b = semantico_de(izq), semantico_de(der)
        for lado, sem, txt in (("frente", sem_f, izq), ("fondo", sem_b, der)):
            if sem is None:
                raise ValueError(f"{f['id']}: el lado {lado} «{txt.strip()}» no corresponde a ningún semántico de C2-tokens.json")
        umbral = umbral_de(f["tipo"])
        for idx, tema, sufijo in ((0, "claro", "light"), (1, "oscuro", "dark")):
            tf, tb = primitivo(hojas, f"{sem_f}.{sufijo}"), primitivo(hojas, f"{sem_b}.{sufijo}")
            usados.update((tf, tb))
            # Cruce secundario 2: lo que dicen las columnas de §10.5
            for lado, tok, tabla in (("frente", tf, f["claro" if idx == 0 else "oscuro"][0]),
                                     ("fondo", tb, f["claro" if idx == 0 else "oscuro"][1])):
                if tok != tabla:
                    avisos.append(f"{f['id']} {tema}: {lado} de «{f['par']}» resuelve a {tok} en C2-tokens.json y §10.5 dice {tabla}")
            hf, hb = hex_de(hojas, tf), hex_de(hojas, tb)
            r = contraste(hf, hb)
            razon = truncar2(r)
            if umbral is None:
                cumple, um = "exento", "exento"
            else:
                cumple = "true" if r >= umbral else "false"  # se compara el valor exacto, no el truncado
                um = str(umbral)
            previo = cruce.get((tf, tb, tema))
            if previo is None:
                cr = "nuevo"
            elif previo == razon:
                cr = "igual"
            else:
                cr = f"difiere (contraste.csv {previo})"
            salida.append({"id": f["id"], "componente": f["componente"], "estado": f["estado"], "tema": tema,
                           "par": f["par"], "token_frente": tf, "token_fondo": tb, "hex_frente": hf,
                           "hex_fondo": hb, "razon": razon, "tipo": f["tipo"], "umbral": um,
                           "cumple": cumple, "cruce_contraste_csv": cr, "_exacto": r,
                           "_dice": f["en_csv_dice"].split(" · ")[idx]})
    # Cruce secundario 3: el hex de cada primitivo usado, contra C2-colores.json
    for t in sorted(usados):
        if hex_col.get(t) != hex_de(hojas, t):
            avisos.append(f"{t}: C2-tokens.json resuelve {hex_de(hojas, t)} y C2-colores.json tiene {hex_col.get(t)}")
    return salida, avisos


def a_texto_csv(filas) -> str:
    buf = io.StringIO()
    w = csv.writer(buf, lineterminator="\n")
    w.writerow(COLUMNAS)
    for r in filas:
        w.writerow([r[c] for c in COLUMNAS])
    return buf.getvalue()


def resumen(filas, avisos):
    pares = len({r["id"] for r in filas})
    con = [r for r in filas if r["umbral"] != "exento"]
    dec = [r for r in filas if r["umbral"] == "exento"]
    bajo = [r for r in con if r["cumple"] == "false"]
    nuevos = [r for r in filas if r["cruce_contraste_csv"] == "nuevo"]
    iguales = [r for r in filas if r["cruce_contraste_csv"] == "igual"]
    difieren = [r for r in filas if r["cruce_contraste_csv"].startswith("difiere")]
    # lo que DESIGN.md dice sobre «en contraste.csv» contra lo que se encuentra
    contradice = [r for r in filas if (r["_dice"] == "sí") != (r["cruce_contraste_csv"] != "nuevo")]
    return {"pares": pares, "combinaciones": len(filas), "con_umbral": len(con), "decorativas": len(dec),
            "bajo": bajo, "nuevos": nuevos, "iguales": iguales, "difieren": difieren,
            "contradice": contradice, "avisos": avisos}


# ---------------------------------------------------------------- segundo calculo, escrito aparte
def independiente_razon(hex1: str, hex2: str):
    """Segunda implementacion, deliberadamente distinta: coma flotante, umbral 0.03928 de WCAG 2.0,
    luminancia como producto punto y calculo de canales en una lista. No comparte codigo con lineal()."""
    def canal(c):
        s = c / 255.0
        return s / 12.92 if s <= 0.03928 else math.pow((s + 0.055) / 1.055, 2.4)
    def lum(h):
        v = [canal(int(h[i:i + 2], 16)) for i in (1, 3, 5)]
        return sum(k * x for k, x in zip((0.2126, 0.7152, 0.0722), v)), v
    (la, va), (lb, vb) = lum(hex1), lum(hex2)
    return (max(la, lb) + 0.05) / (min(la, lb) + 0.05), la, lb, va, vb


def elegir_10(filas):
    """Deterministico: 4 con menor margen sobre su umbral, 2 decorativos y 4 repartidos en el CSV."""
    con = [r for r in filas if r["umbral"] != "exento"]
    estrechos = sorted(con, key=lambda r: (r["_exacto"] / Decimal(r["umbral"]), r["id"], r["tema"]))[:4]
    dec = [r for r in filas if r["umbral"] == "exento"]
    decs = [dec[0], dec[-1]]
    usados = {(r["id"], r["tema"]) for r in estrechos + decs}
    repartidos = []
    for k in (40, 120, 200, 280):
        i = k
        while (filas[i]["id"], filas[i]["tema"]) in usados or filas[i]["umbral"] == "exento":
            i += 1
        repartidos.append(filas[i])
        usados.add((filas[i]["id"], filas[i]["tema"]))
    return estrechos + decs + repartidos


def verificar_independiente(filas):
    elegidas = elegir_10(filas)
    ok = True
    print("id    tema    frente   fondo    razon_csv  independiente  dif   L_frente   L_fondo")
    for r in elegidas:
        ind, la, lb, _, _ = independiente_razon(r["hex_frente"], r["hex_fondo"])
        csvr = float(r["razon"])
        buen = (csvr - 1e-9) <= ind < (csvr + 0.01 + 1e-9)  # el CSV trunca: csv <= real < csv + 0.01
        ok &= buen
        print(f"{r['id']}  {r['tema']:7} {r['hex_frente']} {r['hex_fondo']}  {r['razon']:>8}  {ind:13.4f}  "
              f"{ind - csvr:+.4f}  {la:.5f}  {lb:.5f}  {'OK' if buen else 'DIFIERE'}")
    print(f"independiente: {len(elegidas)} razones comprobadas, {'todas coinciden' if ok else 'HAY DIFERENCIAS'}")
    return ok


# ---------------------------------------------------------------- principal
def main(argv):
    try:
        filas, avisos = generar()
    except (ValueError, StopIteration, KeyError) as e:
        print(f"ERROR de forma en DESIGN.md o C2-tokens.json: {e!r}", file=sys.stderr)
        return 2
    texto = a_texto_csv(filas)
    s = resumen(filas, avisos)
    forma_ok = (s["pares"], s["combinaciones"], s["con_umbral"], s["decorativas"]) == (
        ESPERADOS["pares"], ESPERADOS["combinaciones"], ESPERADOS["con_umbral"], ESPERADOS["decorativas"])
    codigo = 0
    if "--check" in argv:
        en_disco = SALIDA.read_text(encoding="utf-8") if SALIDA.exists() else None
        identico = en_disco == texto
        print(f"check: CSV en disco {'IDENTICO' if identico else 'DISTINTO'} al regenerado")
        codigo = 0 if identico else 1
    elif "--independiente" not in argv:
        SALIDA.write_text(texto, encoding="utf-8")
        print(f"escrito {SALIDA.relative_to(RAIZ)} ({len(texto.encode())} bytes)")
    print(f"pares: {s['pares']} · combinaciones: {s['combinaciones']} · con umbral: {s['con_umbral']} · "
          f"decorativas exentas: {s['decorativas']}")
    print(f"bajo umbral: {len(s['bajo'])}")
    for r in s["bajo"]:
        print(f"  BAJO {r['id']} {r['tema']} {r['componente']} / {r['estado']}: {r['razon']} < {r['umbral']}")
    print(f"cruce con contraste.csv: igual {len(s['iguales'])} · nuevo {len(s['nuevos'])} · difiere {len(s['difieren'])}")
    print(f"combinaciones cuyo «en contraste.csv» de DESIGN.md no coincide con lo encontrado: {len(s['contradice'])}")
    print(f"cruces secundarios (§7.1, columnas de §10.5, C2-colores.json) que no coinciden con C2-tokens.json: {len(s['avisos'])}")
    for a in s["avisos"]:
        print("  AVISO", a)
    if not forma_ok:
        print(f"AVISO: forma distinta de la esperada {ESPERADOS}")
        codigo = max(codigo, 1)
    if "--independiente" in argv:
        if not verificar_independiente(filas):
            codigo = max(codigo, 1)
    if s["bajo"] or s["difieren"] or s["avisos"]:
        codigo = max(codigo, 1)
    return codigo


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
