#!/usr/bin/env python3
"""Genera docs/contratos/C2-tokens.json (DTCG). Cadena unica de C2 (README de esta carpeta):
  DESIGN.md §3-§10.6 + C2-colores.json -> este script -> C2-tokens.json (NO se edita a mano)
  -> Style Dictionary -> tokens.css y tailwind.theme.css.
Entradas:
  - docs/contratos/C2-colores.json   (190 primitivos; se copian sin cambiar un hex. En los 8 roles
    y la accion critica, light/dark/surface-light/surface-dark se emiten como alias de los pasos
    700/400/100/900, que tienen el mismo hex: C2-10)
  - docs/design/DESIGN.md §7.1       (41 semanticos ui.*; se leen de la tabla)
  - los valores no de color de DESIGN.md §3-§6, §9.3-§9.4 y §10.6, transcritos abajo (si cambian en
    DESIGN, se cambian aqui y se corre `npm run check`; verificar_c2.py compara nombres con DESIGN)
Uso: python3 -I docs/design/tools/tokens/generar_c2_tokens.py
Salida determinista (mismo resultado en cada corrida)."""
import json, re, sys
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[4]
COLORES = RAIZ / "docs/contratos/C2-colores.json"
DESIGN = RAIZ / "docs/design/DESIGN.md"
SALIDA = RAIZ / "docs/contratos/C2-tokens.json"
EXT = "finanzas"  # espacio de nombres de $extensions propios


def hoja(valor, desc, tipo=None, **ext):
    t = {}
    if tipo:
        t["$type"] = tipo
    t["$value"] = valor
    t["$description"] = desc
    if ext:
        t["$extensions"] = {EXT: ext}
    return t


def grupo(desc, **hijos):
    g = {"$description": desc}
    g.update(hijos)
    return g


def colores():
    return json.loads(COLORES.read_text(encoding="utf-8"))["color"]


PASOS_ALIAS = {"light": "700", "dark": "400", "surface-light": "100", "surface-dark": "900"}


def alias_pasos(c):
    """C2-10: en los 8 roles y la accion critica, light/dark/surface-* son la misma idea que los
    pasos 700/400/100/900 (DESIGN §2.4: texto 700 en claro y 400 en oscuro; insignia 100 y 900).
    Se emiten como referencia a su paso; el hex no cambia (se verifica antes de reemplazar)."""
    grupos = [(f"color.role.{r}", c["role"][r]) for r in c["role"] if not r.startswith("$")]
    grupos.append(("color.critical", c["critical"]))
    n = 0
    for base, nodo in grupos:
        for hoja_, paso in PASOS_ALIAS.items():
            assert nodo[hoja_]["$value"].upper() == nodo[paso]["$value"].upper(), (base, hoja_, paso)
            nodo[hoja_]["$value"] = f"{{{base}.{paso}}}"
            n += 1
    assert n == 36, f"se esperaban 36 alias (9 grupos x 4), hay {n}"
    return c


def leer_hex(c, ruta, vistos=()):
    """Hex de un primitivo; sigue los alias de alias_pasos ({color.role.X.700})."""
    n = c
    for p in ruta.split(".")[1:]:
        n = n[p]
    v = n["$value"]
    if v.startswith("{"):
        assert ruta not in vistos, f"ciclo en {ruta}"
        return leer_hex(c, v.strip("{}"), vistos + (ruta,))
    return v.upper()


def semanticos(c):
    """Parsea la tabla de DESIGN §7.1 y verifica cada hex contra C2-colores."""
    txt = DESIGN.read_text(encoding="utf-8").split("\n")
    ini = next(i for i, l in enumerate(txt) if l.startswith("### 7.1 "))
    fin = next(i for i, l in enumerate(txt) if l.startswith("### 7.2 "))
    filas = []
    for l in txt[ini:fin]:
        m = re.match(r"\| `(ui\.[^`]+)` \| `([^`]+)` · `(#[0-9A-Fa-f]{6})` \| `([^`]+)` · `(#[0-9A-Fa-f]{6})` \| ([^|]*)\| (.*) \|$", l)
        if m:
            filas.append(m.groups())
    assert len(filas) == 41, f"se esperaban 41 semanticos, hay {len(filas)}"
    ui = {}
    for nombre, cl, hl, cd, hd, alias, uso in filas:
        assert leer_hex(c, cl) == hl.upper(), (nombre, cl, hl)
        assert leer_hex(c, cd) == hd.upper(), (nombre, cd, hd)
        al = re.search(r"`(ui\.[^`]+)`", alias)
        vl, vd = (f"{{{al.group(1)}.light}}", f"{{{al.group(1)}.dark}}") if al else (f"{{{cl}}}", f"{{{cd}}}")
        p = nombre.split(".")[1:]
        n = ui
        for k in p[:-1]:
            n = n.setdefault(k, {})
        n[p[-1]] = {"$type": "color", "$description": uso.strip().capitalize() + " (DESIGN §7.1).",
                    "light": {"$value": vl}, "dark": {"$value": vd}}
    # Pares de rol, critico, persona, tienda y serie (DESIGN §10.3): alias de los primitivos
    roles = {"income": "ingreso", "expense": "gasto", "savings": "ahorro", "debt": "deuda y crédito",
             "warning": "advertencia", "error": "error", "success": "éxito", "info": "información"}
    ui["role"] = {}
    for r, es in roles.items():
        ui["role"][r] = {"$type": "color", "$description": f"Rol «{es}»: texto e icono (DESIGN §2.4).",
                         "light": {"$value": f"{{color.role.{r}.light}}"}, "dark": {"$value": f"{{color.role.{r}.dark}}"}}
        ui["role"][r + "-surface"] = {"$type": "color", "$description": f"Rol «{es}»: fondo de insignia (DESIGN §2.4).",
                         "light": {"$value": f"{{color.role.{r}.surface-light}}"}, "dark": {"$value": f"{{color.role.{r}.surface-dark}}"}}
    ui["critical"] = {"$type": "color", "$description": "Acción crítica: texto e icono (DESIGN §2.5).",
                      "light": {"$value": "{color.critical.light}"}, "dark": {"$value": "{color.critical.dark}"}}
    ui["critical-surface"] = {"$type": "color", "$description": "Acción crítica: fondo de insignia (DESIGN §2.5).",
                      "light": {"$value": "{color.critical.surface-light}"}, "dark": {"$value": "{color.critical.surface-dark}"}}
    ui["person"] = {k: {"$type": "color", "$description": f"Persona «{n}» en la paleta categórica (DESIGN §2.6).",
                        "light": {"$value": f"{{color.person.{k}.light}}"}, "dark": {"$value": f"{{color.person.{k}.dark}}"}}
                    for k, n in {"axel": "Axel", "ximena": "Ximena", "karen": "Karen", "shared": "Conjunto"}.items()}
    ui["store"] = {k: {"$type": "color", "$description": f"Tienda «{n}» en la paleta categórica (DESIGN §2.6).",
                       "light": {"$value": f"{{color.store.{k}.light}}"}, "dark": {"$value": f"{{color.store.{k}.dark}}"}}
                   for k, n in {"walmart": "Walmart", "bodega-aurrera": "Bodega Aurrera", "sams-club": "Sam's Club", "market": "Mercado"}.items()}
    ui["chart"] = {str(i): {"$type": "color", "$description": f"Serie {i} de gráficas, orden fijo (DESIGN §2.7, §8.6).",
                            "light": {"$value": f"{{color.chart.{i}.light}}"}, "dark": {"$value": f"{{color.chart.{i}.dark}}"}}
                   for i in range(1, 9)}
    ui["$type"] = "color"
    ui["$description"] = "Capa semántica: cada token es un par claro/oscuro que apunta a un primitivo (DESIGN §7.1, §10.3). Los componentes usan solo estos nombres."
    return ui


def px(n):
    return f"{n}px"


def tipografia():
    esc = [("caption", 12, 16, 400, "metadatos, etiquetas de eje, pies de tabla, insignias (500 en etiquetas)"),
           ("small", 14, 20, 400, "interfaz densa: barra lateral, tablas, propiedades, menús, botones en escritorio"),
           ("body", 16, 24, 400, "texto de páginas y notas, campos de formulario, botones en móvil"),
           ("title-3", 20, 28, 600, "títulos de sección y de tarjeta, cifra secundaria"),
           ("title-2", 24, 32, 600, "cifra principal de una tarjeta de métrica, título de hoja o diálogo"),
           ("title-1", 32, 40, 700, "título de página")]
    f = {"$description": "Tipografía (DESIGN §3): Inter 4.1 variable autoalojada, seis tamaños, cifras tabulares en todo dinero."}
    f["family"] = {"$type": "fontFamily", "sans": hoja(
        ["Inter", "ui-sans-serif", "system-ui", "Noto Sans", "Roboto", "sans-serif", "Noto Color Emoji"],
        "Familia de la interfaz: Inter variable (SIL OFL 1.1, autoalojada) con respaldo de las fuentes del sistema (DESIGN §3.1).",
        faces=[{"style": "normal", "weight": "100 900", "src": "../design/fuentes/inter/web/InterVariable.woff2", "format": "woff2"},
               {"style": "italic", "weight": "100 900", "src": "../design/fuentes/inter/web/InterVariable-Italic.woff2", "format": "woff2"}],
        license="../design/fuentes/inter/LICENSE.txt",
        numeric="tabular-nums")}
    f["size"] = {"$type": "dimension"}
    f["line-height"] = {"$type": "dimension"}
    for n, s, lh, w, uso in esc:
        f["size"][n] = hoja(px(s), f"Tamaño {s} px: {uso} (DESIGN §3.2).")
        f["line-height"][n] = hoja(px(lh), f"Interlínea de {n}: {lh} px, múltiplo de 4 (DESIGN §3.2).")
    f["weight"] = {"$type": "fontWeight",
                   "regular": hoja(400, "Texto corrido y celdas (DESIGN §3.2)."),
                   "medium": hoja(500, "Etiquetas, elemento activo y botones (DESIGN §3.2)."),
                   "semibold": hoja(600, "Títulos de 20 y 24 y cifras principales (DESIGN §3.2)."),
                   "bold": hoja(700, "Solo el título de página (DESIGN §3.2).")}
    f["letter-spacing"] = hoja("0px", "Espaciado entre letras en 0: el tamaño óptico de Inter ya lo corrige (DESIGN §3.2).", "dimension")
    return f


def espaciado():
    usos = {0: "sin separación", 1: "icono y texto dentro de un control; entre insignias", 2: "entre controles relacionados; separación mínima entre objetivos táctiles",
            3: "relleno horizontal de botones y campos; entre elementos de una lista", 4: "margen lateral y medianil en móvil; relleno de tarjeta en móvil",
            6: "relleno de tarjeta desde 768; medianil desde 1024; entre bloques", 8: "margen lateral a 1024; entre secciones de una página",
            12: "margen lateral a 1440; aire sobre el título de página"}
    g = {"$type": "dimension", "$description": "Espaciado sobre retícula de 4 pt (DESIGN §4.1); en Tailwind 4, --spacing: 0.25rem."}
    for k, v in usos.items():
        g[str(k)] = hoja(px(k * 4), f"{k * 4} px: {usos[k]} (DESIGN §4.1).")
    return g


def formas():
    r = {"$type": "dimension", "$description": "Radios (DESIGN §4.4): un radio interior nunca es mayor que el de su contenedor.",
         "sm": hoja("6px", "Botones, campos, chips, elementos de menú e insignias (DESIGN §4.4)."),
         "md": hoja("10px", "Tarjetas, menús emergentes, avisos breves y vista previa (DESIGN §4.4)."),
         "lg": hoja("14px", "Diálogos, hojas inferiores y paneles grandes (DESIGN §4.4)."),
         "full": hoja("9999px", "Forma circular: avatares, puntos de categórica, interruptores; no es un paso de la escala (DESIGN §4.4).")}
    b = {"$type": "dimension", "$description": "Grosores de borde (DESIGN §4.5).",
         "width": {"default": hoja("1px", "Divisores, bordes de tarjeta y de campo (DESIGN §4.5)."),
                   "focus": hoja("2px", "Anillo de foco (DESIGN §4.5).")}}
    fo = {"$type": "dimension", "offset": hoja("2px", "Separación entre el control y el anillo de foco (DESIGN §4.5).")}
    capa = lambda y, bl, a: dict(color="{color.neutral.950}", offsetX="0px", offsetY=px(y), blur=px(bl), spread="0px")
    sh = {"$type": "shadow", "$description": "Sombras: dos niveles en claro; en oscuro no hay sombra y se usa borde de 1 px (DESIGN §4.5).",
          "0": hoja([dict(color="{color.neutral.950}", offsetX="0px", offsetY="0px", blur="0px", spread="0px")],
                    "Sin sombra: tarjetas en reposo, planas (DESIGN §4.5).", alpha=[0], css="none"),
          "1": hoja([capa(1, 2, 0), capa(1, 3, 0)],
                    "Nivel 1, solo claro: menús emergentes, tarjeta arrastrada, vista previa (DESIGN §4.5). Capas al 8 % y 6 % de neutral.950; en oscuro, borde 1 px de ui.border.divider.",
                    alpha=[0.08, 0.06], dark="none"),
          "2": hoja([capa(4, 12, 0), capa(12, 32, 0)],
                    "Nivel 2, solo claro: diálogos, hojas inferiores, menú «/», avisos (DESIGN §4.5). Capas al 10 % y 8 % de neutral.950; en oscuro, borde 1 px de ui.border.divider.",
                    alpha=[0.10, 0.08], dark="none")}
    ov = {"$type": "color", "$description": "Fondo atenuado bajo diálogos, hojas y panel ☰ (DESIGN §4.5): neutral.950 con transparencia.",
          "scrim": {"$description": "Velo: neutral.950 al 40 % en claro y 60 % en oscuro (DESIGN §4.5).",
                    "light": {"$value": "{color.neutral.950}", "$extensions": {EXT: {"alpha": 0.4}}},
                    "dark": {"$value": "{color.neutral.950}", "$extensions": {EXT: {"alpha": 0.6}}}}}
    return r, b, fo, sh, ov


def medidas():
    S = lambda v, d: hoja(px(v), d)
    size = {"$type": "dimension", "$description": "Medidas de control, icono y componente (DESIGN §4.6, §5.2, §10.6)."}
    size["touch"] = {"min": S(44, "Área de toque mínima de todo control táctil, aunque el dibujo sea menor (DESIGN §4.6).")}
    size["control"] = {"coarse": hoja("{size.touch.min}", "Altura de botones, campos, filas y chips con puntero grueso; alias de size.touch.min (DESIGN §4.6)."),
                       "primary-sheet": S(48, "Botón primario de hojas inferiores en móvil (DESIGN §4.6)."),
                       "fine": S(36, "Altura visual en escritorio con puntero fino; no es objetivo táctil (DESIGN §4.6).")}
    size["icon"] = {"sm": S(20, "Icono junto a texto de 14 y 16 px (DESIGN §5.2)."), "md": S(24, "Icono junto a títulos o solo en barras (DESIGN §5.2).")}
    size["bottombar"] = {"item": {"height": hoja("{layout.bottombar.height}", "Alto de cada destino de la barra inferior: el de la barra, sin la zona segura; alias de layout.bottombar.height (DESIGN §4.6, §9.3)."),
                                  "min-width": S(64, "Ancho mínimo de cada destino de la barra inferior; en 360 px caben 5 destinos de 72 px (DESIGN §4.6, §9.3).")},
                         "add": S(40, "Botón «+» de la barra inferior (DESIGN §7.10)."),
                         "indicator-width": S(56, "Ancho de la píldora del destino activo (DESIGN §7.10)."),
                         "indicator-height": S(32, "Alto de la píldora del destino activo (DESIGN §7.10).")}
    size["sidebar"] = {"item": S(32, "Elemento de la barra lateral y de la lista de Ajustes con puntero fino; con puntero grueso, size.touch.min (DESIGN §7.9, §7.13).")}
    size["chip"] = {"height": S(24, "Chip de persona, tienda o filtro (DESIGN §7.4).")}
    size["badge"] = {"height": S(20, "Insignia de rol (DESIGN §7.4).")}
    size["avatar"] = S(24, "Avatar con anillo de 2 px (DESIGN §7.4).")
    size["mark"] = {"category": S(10, "Marca de forma de persona o tienda (DESIGN §7.4).")}
    size["checkbox"] = S(20, "Casilla, dentro de un área de 44 px (DESIGN §7.3).")
    size["switch"] = {"width": S(36, "Ancho del interruptor (DESIGN §7.3)."), "height": S(20, "Alto del interruptor (DESIGN §7.3).")}
    size["meter"] = {"height": S(8, "Barra de avance (DESIGN §7.6).")}
    size["sheet"] = {"handle-width": S(32, "Ancho del asa de la hoja inferior (DESIGN §7.12)."), "handle-height": S(4, "Alto del asa de la hoja inferior (DESIGN §7.12).")}
    size["chart"] = {"height-compact": S(200, "Alto de gráfica en menos de 768 px (DESIGN §7.7)."), "height": S(280, "Alto de gráfica desde 768 px (DESIGN §7.7).")}
    size["calendar"] = {"cell": S(96, "Celda del mes a 1024 px (DESIGN §7.8)."), "cell-wide": S(120, "Celda del mes a 1440 px (DESIGN §7.8).")}
    size["agenda"] = {"row": S(56, "Fila de agenda y de pagos (DESIGN §7.8).")}
    size["viz"] = {"bar-max": S(24, "Grosor máximo de una barra (DESIGN §8.1)."), "bar-radius": S(4, "Radio de las esquinas de una barra (DESIGN §8.1)."),
                   "line": S(2, "Grosor de línea de serie (DESIGN §8.1)."), "marker": S(10, "Tamaño del marcador de serie (DESIGN §8.1)."),
                   "gap": S(2, "Separación entre marcas de una gráfica (DESIGN §8.1).")}
    lay = {"$type": "dimension", "$description": "Medidas de estructura (DESIGN §4.3, §10.6)."}
    lay["sidebar"] = {"width": S(240, "Barra lateral fija desde 1024 (DESIGN §4.3)."), "overlay-width": S(280, "Barra lateral superpuesta a 768 (DESIGN §4.3)."),
                      "mobile-width": hoja(85, "Panel ☰ en 360 y 412: 85 % del ancho de pantalla (DESIGN §4.3).", "number", unit="%")}
    lay["topbar"] = {"height": {"mobile": S(56, "Barra superior en menos de 768 px (DESIGN §4.3)."), "desktop": S(48, "Barra superior desde 768 px (DESIGN §4.3).")}}
    lay["bottombar"] = {"height": hoja("56px", "Barra inferior; el consumidor suma env(safe-area-inset-bottom) (DESIGN §4.3).", "dimension")}
    lay["reading"] = {"max-width": S(720, "Máximo de lectura de páginas de texto desde 1440 (DESIGN §4.3).")}
    lay["preview"] = {"width": {"desktop": S(400, "Vista previa lateral a 1024 (DESIGN §4.3)."), "wide": S(480, "Vista previa lateral a 1440 (DESIGN §4.3).")}}
    lay["dialog"] = {"width": {"short": S(480, "Diálogo corto de confirmación desde 768 (DESIGN §4.3)."), "long": S(640, "Diálogo largo de formularios desde 768 (DESIGN §4.3).")}}
    lay["command"] = {"width": S(640, "Panel del menú «/» y búsqueda desde 768 (DESIGN §4.3).")}
    lay["board"] = {"column-width": S(272, "Columna del tablero (DESIGN §7.5.3).")}
    lay["toast"] = {"max-width": hoja("{layout.dialog.width.short}", "Aviso breve; alias del diálogo corto (DESIGN §10.6).", "dimension")}
    icon = {"$description": "Trazo de iconos Lucide (DESIGN §5.2).",
            "stroke": {"$root": hoja("1.5px", "Trazo dibujado de 1.5 px en los dos tamaños (DESIGN §5.2).", "dimension"),
                       "user-units": {"$type": "number", "$description": "stroke-width en unidades del viewBox de 24 que da 1.5 px dibujados (DESIGN §5.2).",
                                      "md": hoja(1.5, "A 24 px de icono: 1.5 × 24 / 24."), "sm": hoja(1.8, "A 20 px de icono: 1.5 × 24 / 20.")}}}
    return size, lay, icon


def movimiento_quiebres_viz():
    m = {"$description": "Movimiento (DESIGN §6): 150–250 ms; la curva enfática solo para entrar.",
         "duration": {"$type": "duration",
                      "fast": hoja("150ms", "Cambios de estado: hover, presionado, color; salidas (DESIGN §6.1)."),
                      "base": hoja("200ms", "Menús emergentes, aviso que entra, indicador de pestaña (DESIGN §6.1)."),
                      "slow": hoja("250ms", "Diálogos, hojas inferiores, panel ☰, vista previa (DESIGN §6.1).")},
         "easing": {"$type": "cubicBezier",
                    "standard": hoja([0.2, 0, 0, 1], "Cambios dentro de la pantalla y salidas (DESIGN §6.1)."),
                    "emphasized": hoja([0.05, 0.7, 0.1, 1], "Entradas de lo que aparece sobre la pantalla (DESIGN §6.1).")}}
    q = {"$type": "dimension", "$description": "Puntos de quiebre con min-width, móvil primero (DESIGN §9.4). No sirven dentro de @media por ser variables.",
         "base": hoja("360px", "Base sin consulta, hasta 411 px (DESIGN §9.4)."), "phone": hoja("412px", "412–767 px (DESIGN §9.4)."),
         "tablet": hoja("768px", "768–1023 px (DESIGN §9.4)."), "desktop": hoja("1024px", "1024–1439 px (DESIGN §9.4)."),
         "wide": hoja("1440px", "Desde 1440 px (DESIGN §9.4).")}
    dash = [None, [8, 4], [2, 4], [8, 4, 2, 4], [4, 4], [12, 4], [8, 3, 2, 3, 2, 3], [2, 2]]
    marca = ["circle", "square", "triangle", "diamond", "triangle-down", "cross", "x", "ring"]
    pat = [dict(kind="solid"), dict(kind="stripes", angle=45, step=6), dict(kind="stripes", angle=135, step=6), dict(kind="dots", size=2, step=6),
           dict(kind="grid-diagonal", step=8), dict(kind="stripes", angle=45, step=10), dict(kind="stripes", angle=135, step=10), dict(kind="dots", size=3, step=8)]
    viz = {"$description": "Forma, trazo y patrón de cada serie de gráficas (DESIGN §8.6, K27): el color nunca va solo.", "series": {"$type": "strokeStyle"}}
    for i in range(8):
        v = "solid" if dash[i] is None else {"dashArray": [px(d) for d in dash[i]], "lineCap": "butt"}
        viz["series"][str(i + 1)] = {"dash": hoja(v, f"Trazo de la serie {i + 1} (DESIGN §8.6).", marker=marca[i], pattern=pat[i])}
    return m, q, viz


def main():
    c = alias_pasos(colores())
    r, b, fo, sh, ov = formas()
    size, lay, icon = medidas()
    m, q, viz = movimiento_quiebres_viz()
    raiz = {
        "$description": "C2 v1.0 · Contrato de tokens de Finanzas Hogar (11-oct-2026). Única fuente de valores de diseño, en formato DTCG. Los 190 primitivos de color vienen de C2-colores.json sin cambiar un hex (en los 8 roles y la acción crítica, light, dark y surface-* son alias de los pasos 700, 400, 100 y 900); el resto sale de DESIGN.md. Lo genera generar_c2_tokens.py: no se edita a mano. Compilación: docs/design/tools/tokens/README.md.",
        "$extensions": {EXT: {"contrato": "C2", "version": "1.0", "fecha": "2026-10-11",
                              "cambios": ["1.0 (11-oct-2026): primera versión. Primitivos de color 1.2 (K26, K41, K42), 41 semánticos ui.* más pares de rol, crítico, persona, tienda y serie, tipografía, espaciado, radios, bordes, sombras, velo, medidas, movimiento, puntos de quiebre y trazos de series.",
                                          "1.0 · corrección de la revisión U2 r1 (11-oct-2026, antes de la puerta 2; ningún hex cambia): light, dark y surface-* de los 8 roles y de la acción crítica pasan a alias de los pasos 700, 400, 100 y 900 (C2-10); size.bottombar.item queda como height (alias de layout.bottombar.height, 56px) y min-width 64px, en lugar de min-width 56px y min-height 64px (DESIGN §4.6, §9.3)."]}},
        "color": c, "ui": semanticos(c), "font": tipografia(), "space": espaciado(), "radius": r, "border": b, "focus": fo,
        "shadow": sh, "overlay": ov, "size": size, "layout": lay, "icon": icon, "motion": m, "breakpoint": q, "viz": viz}
    SALIDA.write_text(json.dumps(raiz, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print("escrito", SALIDA.relative_to(RAIZ), SALIDA.stat().st_size, "B")


if __name__ == "__main__":
    sys.exit(main())
