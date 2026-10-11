#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Colorimetría del candidato C · Índigo 280° · Finanzas Hogar (U2a, WF-0).

Script reproducible que usa solo la biblioteca estándar de Python 3. Al ejecutarlo
regenera, a partir de las reglas de este archivo:

  docs/design/colorimetria.md      documento con método, escalas, roles y tablas calculadas
  docs/design/contraste.csv        todos los pares de contraste (WCAG 2.x) en claro y oscuro
  docs/contratos/C2-colores.json   tokens DTCG de color (parte de color del contrato C2)
  docs/design/paleta.html          muestrario estático (sin CDN ni fuentes externas)
  docs/design/paleta.png           captura con chromium --headless=new (si chromium existe)

Uso:
  python3 docs/design/tools/colorimetria.py                  # genera todo y verifica
  python3 docs/design/tools/colorimetria.py --sin-png        # genera sin la captura PNG
  python3 docs/design/tools/colorimetria.py --solo-verificar # verifica los archivos existentes

Fuentes de las fórmulas (detalle y citas en colorimetria.md §2):
  * OKLab/OKLCH: Björn Ottosson, «A perceptual color space for image processing»,
    https://bottosson.github.io/posts/oklab/ (matrices M1 y M2 publicadas ahí).
  * sRGB: IEC 61966-2-1:1999 (función de transferencia con umbrales 0.04045 y 0.0031308).
  * Luminancia relativa y razón de contraste: WCAG 2.2 (https://www.w3.org/TR/WCAG22/),
    definiciones «relative luminance» y «contrast ratio».
  * Daltonismo: G. M. Machado, M. M. Oliveira y L. A. F. Fernandes (2009), «A
    Physiologically-based Model for Simulation of Color Vision Deficiency», IEEE TVCG
    15(6):1291-1298, doi:10.1109/TVCG.2009.113. Matrices de severidad 1.0 de la Tabla 1 de
    https://www.inf.ufrgs.br/~oliveira/pubs_files/CVD_Simulation/CVD_Simulation.html
    (consultada el 2026-10-10), aplicadas en RGB lineal (el artículo, p. 1294, asume RGB
    lineal; así lo hace el paquete R colorspace >= 2.1-0 con simulate_cvd(linear = TRUE)).

El script no usa red ni dependencias externas. Chromium solo se invoca, si existe, para la
captura PNG.
"""
import argparse
import csv
import html as _html
import io
import json
import math
import re
import shutil
import subprocess
import sys
from decimal import Decimal, ROUND_DOWN
from itertools import combinations, permutations, product
from pathlib import Path

RAIZ = Path(__file__).resolve().parents[3]
SALIDA = {
    "md": RAIZ / "docs" / "design" / "colorimetria.md",
    "csv": RAIZ / "docs" / "design" / "contraste.csv",
    "json": RAIZ / "docs" / "contratos" / "C2-colores.json",
    "html": RAIZ / "docs" / "design" / "paleta.html",
    "png": RAIZ / "docs" / "design" / "paleta.png",
}
FECHA = "2026-10-10"
VERSION = "1.1"
FECHA_11 = "2026-10-11"   # puerta 1: respuestas K26 y K27 de Axel y hallazgo R3-04 (sección «Cambios 1.1»)
SCRIPT_REL = "docs/design/tools/colorimetria.py"

# =============================================================================
# 1. Conversiones de color (Ottosson, IEC 61966-2-1, WCAG 2.x)
# =============================================================================


def srgb_a_lineal(c):
    """sRGB codificado [0,1] -> sRGB lineal (IEC 61966-2-1)."""
    return c / 12.92 if c <= 0.04045 else ((c + 0.055) / 1.055) ** 2.4


def lineal_a_srgb(c):
    """sRGB lineal [0,1] -> sRGB codificado (IEC 61966-2-1)."""
    return 12.92 * c if c <= 0.0031308 else 1.055 * c ** (1.0 / 2.4) - 0.055


def _raiz_cubica(x):
    return math.copysign(abs(x) ** (1.0 / 3.0), x)


def lineal_a_oklab(rgb):
    """sRGB lineal -> OKLab (matrices M1 y M2 de Ottosson)."""
    r, g, b = rgb
    l = 0.4122214708 * r + 0.5363325363 * g + 0.0514459929 * b
    m = 0.2119034982 * r + 0.6806995451 * g + 0.1073969566 * b
    s = 0.0883024619 * r + 0.2817188376 * g + 0.6299787005 * b
    l_, m_, s_ = _raiz_cubica(l), _raiz_cubica(m), _raiz_cubica(s)
    return (0.2104542553 * l_ + 0.7936177850 * m_ - 0.0040720468 * s_,
            1.9779984951 * l_ - 2.4285922050 * m_ + 0.4505937099 * s_,
            0.0259040371 * l_ + 0.7827717662 * m_ - 0.8086757660 * s_)


def oklab_a_lineal(lab):
    """OKLab -> sRGB lineal (inversas de M2 y M1 de Ottosson)."""
    L, a, b = lab
    l_ = L + 0.3963377774 * a + 0.2158037573 * b
    m_ = L - 0.1055613458 * a - 0.0638541728 * b
    s_ = L - 0.0894841775 * a - 1.2914855480 * b
    l, m, s = l_ ** 3, m_ ** 3, s_ ** 3
    return (4.0767416621 * l - 3.3077115913 * m + 0.2309699292 * s,
            -1.2684380046 * l + 2.6097574011 * m - 0.3413193965 * s,
            -0.0041960863 * l - 0.7034186147 * m + 1.7076147010 * s)


def oklch_a_oklab(L, C, h):
    hr = math.radians(h)
    return (L, C * math.cos(hr), C * math.sin(hr))


def oklab_a_oklch(lab):
    L, a, b = lab
    return (L, math.hypot(a, b), math.degrees(math.atan2(b, a)) % 360.0)


def en_gamut(rgb, eps=1e-9):
    return all(-eps <= v <= 1.0 + eps for v in rgb)


def recortar_croma(L, C, h, iteraciones=60):
    """Mayor croma <= C que cabe en el gamut sRGB, por bisección sobre el croma (L y h fijos)."""
    if en_gamut(oklab_a_lineal(oklch_a_oklab(L, C, h))):
        return C
    bajo, alto = 0.0, C
    for _ in range(iteraciones):
        medio = (bajo + alto) / 2.0
        if en_gamut(oklab_a_lineal(oklch_a_oklab(L, medio, h))):
            bajo = medio
        else:
            alto = medio
    return bajo


def lineal_a_hex(rgb):
    """sRGB lineal -> hex de 8 bits, redondeo ROUND_HALF_UP de cada canal."""
    canales = []
    for v in rgb:
        v = min(1.0, max(0.0, v))
        canales.append(int(math.floor(lineal_a_srgb(v) * 255.0 + 0.5)))
    return "#%02X%02X%02X" % tuple(canales)


def hex_a_srgb(hx):
    hx = hx.lstrip("#")
    return tuple(int(hx[i:i + 2], 16) / 255.0 for i in (0, 2, 4))


def hex_a_lineal(hx):
    return tuple(srgb_a_lineal(c) for c in hex_a_srgb(hx))


def hex_a_oklab(hx):
    return lineal_a_oklab(hex_a_lineal(hx))


def hex_a_oklch(hx):
    return oklab_a_oklch(hex_a_oklab(hx))


class Color:
    """Color definido en OKLCH (L, croma pedido, matiz) y recortado al gamut sRGB."""

    __slots__ = ("l", "c", "h", "c_pedido", "hex")

    def __init__(self, l, c_pedido, h):
        self.l, self.h, self.c_pedido = l, h, c_pedido
        self.c = recortar_croma(l, c_pedido, h)
        self.hex = lineal_a_hex(oklab_a_lineal(oklch_a_oklab(l, self.c, h)))

    @property
    def recortado(self):
        return self.c < self.c_pedido - 1e-9


def luminancia(hx):
    """Luminancia relativa WCAG 2.x: 0.2126 R + 0.7152 G + 0.0722 B (lineales)."""
    r, g, b = hex_a_lineal(hx)
    return 0.2126 * r + 0.7152 * g + 0.0722 * b


def contraste(a, b):
    ya, yb = luminancia(a), luminancia(b)
    if ya < yb:
        ya, yb = yb, ya
    return (ya + 0.05) / (yb + 0.05)


def fmt_ratio(r):
    """Razón truncada (ROUND_DOWN) a dos decimales: nunca se muestra mayor de lo que es."""
    return str(Decimal(r).quantize(Decimal("0.01"), rounding=ROUND_DOWN))


# --- Daltonismo: Machado, Oliveira y Fernandes (2009), severidad 1.0 (Tabla 1) -------------
MACHADO_2009 = {
    "protanopia": ((0.152286, 1.052583, -0.204868),
                   (0.114503, 0.786281, 0.099216),
                   (-0.003882, -0.048116, 1.051998)),
    "deuteranopia": ((0.367322, 0.860646, -0.227968),
                     (0.280085, 0.672501, 0.047413),
                     (-0.011820, 0.042940, 0.968881)),
    "tritanopia": ((1.255528, -0.076749, -0.178779),
                   (-0.078411, 0.930809, 0.147602),
                   (0.004733, 0.691367, 0.303900)),
}
TIPOS_CVD = ("protanopia", "deuteranopia", "tritanopia")
VISTAS = ("gris",) + TIPOS_CVD
NOMBRE_VISTA = {"normal": "Visión normal", "gris": "Escala de grises", "protanopia": "Protanopia",
                "deuteranopia": "Deuteranopia", "tritanopia": "Tritanopia"}


def simular(hx, tipo, espacio="lineal"):
    """Devuelve sRGB lineal simulado. espacio='lineal' (método del script) o 'gamma'
    (matriz sobre sRGB codificado, como filterDeficiency* de culori 4.0.2; solo para robustez)."""
    m = MACHADO_2009[tipo]
    v = hex_a_lineal(hx) if espacio == "lineal" else hex_a_srgb(hx)
    out = tuple(min(1.0, max(0.0, m[i][0] * v[0] + m[i][1] * v[1] + m[i][2] * v[2])) for i in range(3))
    return out if espacio == "lineal" else tuple(srgb_a_lineal(x) for x in out)


def hex_simulado(hx, tipo):
    return lineal_a_hex(simular(hx, tipo))


def l_gris(hx):
    """L (OKLab) del gris de igual luminancia: para un gris, L = Y^(1/3)."""
    return luminancia(hx) ** (1.0 / 3.0)


def hex_gris(hx):
    y = luminancia(hx)
    return lineal_a_hex((y, y, y))


def hex_vista(hx, vista):
    if vista == "normal":
        return hx
    if vista == "gris":
        return hex_gris(hx)
    return hex_simulado(hx, vista)


def l_vista(hx, vista, espacio="lineal"):
    if vista == "normal":
        return hex_a_oklab(hx)[0]
    if vista == "gris":
        return l_gris(hx)
    return lineal_a_oklab(simular(hx, vista, espacio))[0]


def delta_e(a, b, tipo=None, espacio="lineal"):
    """Distancia euclidiana en OKLab ×100 (opcionalmente bajo simulación)."""
    if tipo is None:
        return 100.0 * math.dist(hex_a_oklab(a), hex_a_oklab(b))
    return 100.0 * math.dist(lineal_a_oklab(simular(a, tipo, espacio)),
                             lineal_a_oklab(simular(b, tipo, espacio)))


def dif_matiz(a, b):
    d = abs(a - b) % 360.0
    return min(d, 360.0 - d)


# =============================================================================
# 2. Reglas de diseño (datos)
# =============================================================================
PASOS = (50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950)
L_PASO = dict(zip(PASOS, (0.985, 0.962, 0.925, 0.865, 0.785, 0.695, 0.605, 0.515, 0.425, 0.335, 0.255)))
# Perfil 0.10 -> 1.0 -> 0.45 del croma máximo, explícito por paso (reproduce el plan §4.3).
PERFIL_CROMA = dict(zip(PASOS, (0.10, 0.25, 0.40, 0.65, 0.85, 1.00, 1.00, 0.90, 0.75, 0.60, 0.45)))
MATIZ_BASE = 280.0
MATIZ_OPUESTO = 100.0
CROMA_MARCA = 0.11
CROMA_NEUTRO = 0.012
CROMA_ROL = 0.14
UMBRAL_TEXTO = 4.5
UMBRAL_UI = 3.0
UMBRAL_DL = 0.10          # criterio U2 del plan §3.5 (más estricto que el 0.08 del encargo)
UMBRAL_DL_ENCARGO = 0.08
SEPARACION_MATIZ = 20.0   # ° respecto a roles (plan §4.3, decisiones §I)
SEPARACION_L = 0.15       # o diferencia de L respecto a roles
ZONA_RESERVADA = SEPARACION_MATIZ   # ° a cada lado del opuesto 100° que solo usa color.critical (decisiones §I)
CROMA_VISIBLE = 0.02      # croma desde el que un token cuenta como cromático (el neutro llega a 0.012)
CVD_OBJETIVO, CVD_PISO, NORMAL_PISO = 8.0, 6.0, 15.0   # ΔE OKLab×100 (skill dataviz)
CROMA_PISO_GRAFICAS = 0.10

PLAN_C = {
    "marca": "F9FAFF F0F1FF E1E4FF C9CFFF ADB3F4 9095E0 7579C2 5D609F 46487B 303258 1E203A".split(),
    "neutro": "FAFAFB F2F2F4 E5E6E9 D1D2D8 B7B8C0 9B9CA4 808189 66676E 4D4E53 36363A 222325".split(),
}

ROLES = [
    {"codigo": "income", "nombre": "Ingreso", "matiz": 135.0,
     "origen": "Propuesta U2a con la opción que prevé decisiones §I (135°): a 15° de éxito y a 5° del acento dividido 130°"},
    {"codigo": "expense", "nombre": "Gasto", "matiz": 45.0,
     "origen": "K26 (puerta 1): terracota 40–45°; 45° es el grado del rango más separado de error y advertencia (Cambios 1.1)"},
    {"codigo": "savings", "nombre": "Ahorro", "matiz": 200.0, "origen": "decisiones §I"},
    {"codigo": "debt", "nombre": "Deuda / crédito", "matiz": 310.0, "origen": "decisiones §I"},
    {"codigo": "warning", "nombre": "Advertencia", "matiz": 75.0, "origen": "decisiones §I"},
    {"codigo": "error", "nombre": "Error", "matiz": 20.0, "origen": "decisiones §I (rojo)"},
    {"codigo": "success", "nombre": "Éxito", "matiz": 150.0, "origen": "decisiones §I (reasignado desde 280°)"},
    {"codigo": "info", "nombre": "Información", "matiz": 240.0, "origen": "decisiones §I"},
]
CRITICO = {"codigo": "critical", "nombre": "Acción crítica", "matiz": MATIZ_OPUESTO,
           "origen": "decisiones §I: eje opuesto reservado para confirmar acciones críticas"}

# Valores del plan §4.3 (candidato C) calculados a L 0.56 (claro) y L 0.74 (oscuro).
PLAN_ROLES = {
    "income": ("#218A45", "#6FC082", "150° en el plan"), "expense": ("#B84D49", "#ED8C84", ""),
    "savings": ("#00858A", "#24C1C9", ""), "debt": ("#8D5AB1", "#C295E3", ""),
    "warning": ("#9C6900", "#D79F4C", ""), "error": ("#857500", "#BDAC4A", "100° en el plan"),
    "success": ("", "", "sin valor: 280° se reasignó"), "info": ("#007CB7", "#5CB4EF", ""),
}

# Categóricas: (código, nombre, matiz, croma pedido, L claro, L oscuro, motivo)
PERSONAS = [
    ("axel", "Axel", 256.0, 0.12, 0.31, 0.605,
     "Azul a 256° (6° más violeta que el 250° del plan para conservar croma ≥ 0.10 en L 0.31); como "
     "queda a 16° de información, se separa por luminosidad (ΔL ≥ 0.15 frente al rol en ambos temas)."),
    ("ximena", "Ximena", 350.0, 0.12, 0.54, 0.835,
     "Conserva 350° del plan (30° de error, 40° de deuda)."),
    ("karen", "Karen", 333.0, 0.12, 0.642, 0.725,
     "Versión 1.1 (K26): el 50° de la 1.0 quedaba a 5° de gasto 45°. Con su L (la más alta del grupo en claro, "
     "la última que da 3:1) se separa de los roles por matiz, y sus ventanas factibles quedan entre deuda 310° y "
     "error 20°; 333° queda a 23° de deuda y de Bodega Aurrera, lo más lejos de Ximena (350°) que permite el "
     "reparto, y de ella se separa por ΔL ≥ 0.10. L en claro 0.642 (antes 0.640) para dar margen al escalonamiento. "
     "PROVISIONAL (bloqueo `U2a-rosas-ximena-bodega`): frente a Ximena queda bajo ≥ 20° o ΔL ≥ 0.15 en ambos temas "
     "y se distingue por el refuerzo de §7.7 (Cambios 1.1, punto 2)."),
    ("shared", "Conjunto", 280.0, 0.045, 0.42, 0.94,
     "Lo compartido lleva el matiz de la marca con croma bajo (casi neutro), como el croma 0.03 "
     "del plan; 280° queda a 30° de deuda y 40° de información. Es la única excepción deliberada al "
     "piso de croma 0.10: se lee como «de ambos», no como una persona más."),
]
TIENDAS = [
    ("walmart", "Walmart", 262.0, 0.12, 0.515, 0.716,
     "220° del plan quedaba a 20.0° exactos de ahorro e información (frágil al cuantizar a 8 bits); "
     "se mueve al azul de 262° (22° de información). PROVISIONAL en oscuro (bloqueo `U2a-walmart-oscuro`): "
     "L 0.716 en vez de 0.720 para que Bodega Aurrera quepa en el escalonamiento de las tiendas (Cambios 1.1)."),
    ("bodega-aurrera", "Bodega Aurrera", 356.0, 0.11, 0.41, 0.829,
     "Versión 1.1 (K26): el 50° de la 1.0 quedaba a 5° de gasto 45°. Con su L (entre Mercado y Walmart) se separa "
     "de los roles por matiz en ambos temas: 356° queda a 24° de error y a 23° de Karen, de la que en oscuro no se "
     "separa por luminosidad. Croma 0.11 y L 0.829 en oscuro: el rosa pierde luminancia en gris y con 0.12 no cabe "
     "en el escalonamiento de las tiendas (Cambios 1.1). PROVISIONAL (bloqueo `U2a-rosas-ximena-bodega`): frente a "
     "Ximena queda bajo ≥ 20° o ΔL ≥ 0.15 en ambos temas (en oscuro, a simple vista el mismo color) y se distingue "
     "solo por el refuerzo de §7.7 (Cambios 1.1, punto 2)."),
    ("sams-club", "Sam's Club", 175.0, 0.12, 0.62, 0.93,
     "Conserva 175° del plan (25° de éxito y de ahorro); va en la L más alta de cada tema, donde el "
     "verde azulado conserva croma ≥ 0.10."),
    ("market", "Mercado", 28.0, 0.12, 0.30, 0.615,
     "Terracota a 28° (2° menos que el plan, para quedar a más de 20° de Bodega Aurrera); como queda "
     "a 3° de gasto y 8° de error, se separa por luminosidad (ΔL ≥ 0.15 frente a los roles en ambos temas)."),
]
# Propuesta categórica del plan §4.3 (valor único por color): hex, matiz, L.
PLAN_CATEGORICAS = {"axel": ("#0F74C5", 250, 0.55), "ximena": ("#D375A4", 350, 0.68), "karen": ("#E0AE57", 80, 0.78),
                    "shared": ("#495766", 250, 0.45), "walmart": ("#2FB5D8", 220, 0.72),
                    "bodega-aurrera": ("#FFA658", 60, 0.80), "sams-club": ("#11957C", 175, 0.60), "market": ("#9C4438", 30, 0.50)}
# Valores PROVISIONALES de la versión 1.1 (protocolo §2): token → nota que se agrega a su descripción.
PROVISIONAL = {"color.store.walmart.dark": " · PROVISIONAL: bloqueo U2a-walmart-oscuro (L 0.720 → 0.716)"}
# Ronda 1 de revisión U2a (H-01): Karen y Bodega Aurrera quedan frente a Ximena bajo «≥ 20° o ΔL ≥ 0.15» y solo las
# separa el refuerzo de §7.7; U2-02 admite dos lecturas. PROVISIONAL mientras Axel decide el bloqueo (opciones en
# docs/estado/bloqueos/U2a-rosas-ximena-bodega.md). rosas_ximena() aborta si estas marcas no coinciden con el cálculo.
BLOQUEO_ROSAS = "U2a-rosas-ximena-bodega"
ROSAS_XIMENA = (("person", "karen"), ("store", "bodega-aurrera"))
PROVISIONAL.update({"color.%s.%s.%s" % (g, c, s): " · PROVISIONAL: bloqueo %s (bajo ≥ 20° o ΔL ≥ 0.15 frente a Ximena)"
                    % BLOQUEO_ROSAS for g, c in ROSAS_XIMENA for s in ("light", "dark")})
ORDEN_PERSONAS = ["axel", "ximena", "karen", "shared"]
ORDEN_TIENDAS = ["walmart", "bodega-aurrera", "sams-club", "market"]

# Gráficas: orden de uso, matiz, clase de luminosidad (A alta / B baja) y nombre.
# Ronda 1 de revisión (H-03): la serie 6 era «mostaza» 100°, el matiz que decisiones §I reserva para
# acción crítica; pasa a «azul cielo» 225°, el primer candidato del barrido de §8.8 (barrido_serie). Con
# ese matiz, la búsqueda de orden de §8.2 (ranking_orden) pone óxido en la serie 3 y magenta en la 5.
GRAFICAS = [
    (1, 280.0, "B", "índigo"), (2, 180.0, "A", "verde azulado"), (3, 45.0, "B", "óxido"),
    (4, 55.0, "A", "naranja"), (5, 345.0, "B", "magenta"), (6, 225.0, "A", "azul cielo"),
    (7, 310.0, "B", "violeta"), (8, 140.0, "A", "verde"),
]
SERIE_CORREGIDA = 6       # serie que se movió en la ronda 1 (§8.8)
PASO_BARRIDO = 5          # ° entre matices candidatos del barrido de §8.8
FIJAS_ORDEN = (1, 7, 8)   # series con posición fija en la búsqueda de orden (§8.2)
# Orden de la ronda 0 por nombre (None = la serie corregida, con su matiz nuevo), para compararlo en §8.2.
ORDEN_RONDA0 = ("índigo", "verde azulado", "magenta", "naranja", "óxido", None, "violeta", "verde")
L_CLASE = {("claro", "A"): 0.615, ("claro", "B"): 0.44, ("oscuro", "A"): 0.84, ("oscuro", "B"): 0.62}
CROMA_GRAFICA = {"claro": 0.13, "oscuro": 0.12}

TEMAS = ("claro", "oscuro")
SUPERFICIE = {"claro": {"fondo": "color.neutral.50", "tarjeta": "color.neutral.100"},
              "oscuro": {"fondo": "color.neutral.950", "tarjeta": "color.neutral.900"}}
SUFIJO = {"claro": "light", "oscuro": "dark"}

# =============================================================================
# 3. Construcción del modelo
# =============================================================================


def escala(h, croma_max):
    return {p: Color(L_PASO[p], croma_max * PERFIL_CROMA[p], h) for p in PASOS}


class Modelo:
    def __init__(self):
        self.tokens = {}       # ruta -> dict(hex, l, c, h, desc)
        self.grupos = {}       # ruta de grupo -> descripción
        self.marca = escala(MATIZ_BASE, CROMA_MARCA)
        self.neutro = escala(MATIZ_BASE, CROMA_NEUTRO)
        self.esc_rol = {r["codigo"]: escala(r["matiz"], CROMA_ROL) for r in ROLES + [CRITICO]}
        self.paso_rol = {}     # codigo -> {"claro": paso, "oscuro": paso, "historial": {...}}
        self.cat = {"person": {}, "store": {}}
        self.graf = {}
        self._registrar_todo()

    # -- registro de tokens ---------------------------------------------------
    def reg(self, ruta, color, desc=None):
        if isinstance(color, Color):
            self.tokens[ruta] = {"hex": color.hex, "l": color.l, "c": color.c, "h": color.h, "desc": desc}
        else:
            self.tokens[ruta] = {"hex": color[0], "l": color[1], "c": color[2], "h": color[3], "desc": desc}

    def hex(self, ruta):
        return self.tokens[ruta]["hex"]

    def sup(self, tema, cual):
        return self.hex(SUPERFICIE[tema][cual])

    def _elegir_paso(self, codigo, tema, esc=None):
        esc = self.esc_rol[codigo] if esc is None else esc
        if tema == "claro":
            candidatos = (600, 700, 800, 900)
            fondos = (self.neutro[50].hex, self.neutro[100].hex, esc[100].hex)
        else:
            candidatos = (400, 300, 200, 100)
            fondos = (self.neutro[950].hex, self.neutro[900].hex, esc[900].hex)
        historial = []
        for p in candidatos:
            minimo = min(contraste(esc[p].hex, f) for f in fondos)
            historial.append((p, minimo))
            if minimo >= UMBRAL_TEXTO:
                return p, historial
        raise SystemExit("Ningún paso de %s cumple 4.5:1 en tema %s" % (codigo, tema))

    def _registrar_todo(self):
        self.grupos["color"] = "Parte de color del contrato C2 (U2a, WF-0). Generado por %s." % SCRIPT_REL
        self.reg("color.base.white", ("#FFFFFF", 1.0, 0.0, 0.0), "Blanco puro (texto sobre marca 700 y error 700 en claro)")
        self.grupos["color.brand"] = "Marca · matiz 280° · croma máx. 0.11 · L fija por paso"
        for p in PASOS:
            self.reg("color.brand.%d" % p, self.marca[p])
        self.grupos["color.neutral"] = "Neutro tintado con la marca · matiz 280° · croma máx. 0.012"
        for p in PASOS:
            self.reg("color.neutral.%d" % p, self.neutro[p])
        self.grupos["color.role"] = "Roles semánticos: light/dark = texto e icono por tema; surface-* = fondo de insignia"
        for r in ROLES + [CRITICO]:
            cod = r["codigo"]
            pc, hc = self._elegir_paso(cod, "claro")
            po, ho = self._elegir_paso(cod, "oscuro")
            self.paso_rol[cod] = {"claro": pc, "oscuro": po, "historial": {"claro": hc, "oscuro": ho}}
            base = "color.critical" if cod == "critical" else "color.role.%s" % cod
            self.grupos[base] = "%s · matiz %d° · %s" % (r["nombre"], r["matiz"], r["origen"])
            esc = self.esc_rol[cod]
            self.reg(base + ".light", esc[pc], "Paso %d (texto/icono en claro)" % pc)
            self.reg(base + ".dark", esc[po], "Paso %d (texto/icono en oscuro)" % po)
            self.reg(base + ".surface-light", esc[100], "Paso 100 (fondo de insignia en claro)")
            self.reg(base + ".surface-dark", esc[900], "Paso 900 (fondo de insignia en oscuro)")
            for p in PASOS:
                self.reg("%s.%d" % (base, p), esc[p])
        for grupo, lista, titulo in (("person", PERSONAS, "Personas"), ("store", TIENDAS, "Tiendas")):
            self.grupos["color." + grupo] = "%s · paleta categórica con L escalonada por tema" % titulo
            for cod, nombre, h, c, lc, lo, motivo in lista:
                cc, co = Color(lc, c, h), Color(lo, c, h)
                self.cat[grupo][cod] = {"nombre": nombre, "h": h, "c": c, "claro": cc, "oscuro": co, "motivo": motivo}
                self.grupos["color.%s.%s" % (grupo, cod)] = "%s · matiz %d°" % (nombre, h)
                self.reg("color.%s.%s.light" % (grupo, cod), cc, "%s en claro (L %.3f)%s" % (
                    nombre, lc, PROVISIONAL.get("color.%s.%s.light" % (grupo, cod), "")))
                self.reg("color.%s.%s.dark" % (grupo, cod), co, "%s en oscuro (L %.3f)%s" % (
                    nombre, lo, PROVISIONAL.get("color.%s.%s.dark" % (grupo, cod), "")))
        self.grupos["color.chart"] = "Series de gráficas en orden fijo de uso (1 a 8); variante por tema"
        for n, h, clase, nombre in GRAFICAS:
            self.graf[n] = {"h": h, "clase": clase, "nombre": nombre}
            self.grupos["color.chart.%d" % n] = "Serie %d · %s · matiz %d° · clase %s" % (n, nombre, h, clase)
            for tema in TEMAS:
                col = Color(L_CLASE[(tema, clase)], CROMA_GRAFICA[tema], h)
                self.graf[n][tema] = col
                self.reg("color.chart.%d.%s" % (n, SUFIJO[tema]), col,
                         "Serie %d (%s) en %s" % (n, nombre, tema))


# =============================================================================
# 4. Pares de contraste
# =============================================================================


def construir_pares(M):
    pares = []

    def par(fp, fo, tema, uso, umbral):
        a, b = M.hex(fp), M.hex(fo)
        r = contraste(a, b)
        pares.append({"primer_plano": fp, "fondo": fo, "tema": tema, "uso": uso, "ratio": r,
                      "umbral": umbral, "cumple": r >= umbral, "fp_hex": a, "fo_hex": b})

    T, U = UMBRAL_TEXTO, UMBRAL_UI
    for tema in TEMAS:
        f, t = SUPERFICIE[tema]["fondo"], SUPERFICIE[tema]["tarjeta"]
        sfx = SUFIJO[tema]
        claro = tema == "claro"
        n = lambda p: "color.neutral.%d" % p
        b = lambda p: "color.brand.%d" % p
        e = lambda p: "color.role.error.%d" % p
        k = lambda p: "color.critical.%d" % p
        texto, secundario, deshab = (n(900), n(700), n(700)) if claro else (n(100), n(300), n(400))
        for sup, nom in ((f, "fondo"), (t, "tarjeta")):
            par(texto, sup, tema, "Texto normal sobre %s" % nom, T)
            par(secundario, sup, tema, "Texto secundario sobre %s" % nom, T)
            par(deshab, sup, tema, "Texto deshabilitado sobre %s" % nom, T)
        par(deshab, n(100) if claro else n(900), tema, "Texto de control deshabilitado sobre su relleno", T)
        par(secundario, f, tema, "Marcador de posición (placeholder) dentro de campo", T)
        par(texto, f, tema, "Texto capturado dentro de campo", T)
        enl, enl_h = (b(700), b(800)) if claro else (b(300), b(200))
        for sup, nom in ((f, "fondo"), (t, "tarjeta")):
            par(enl, sup, tema, "Enlace (reposo) sobre %s" % nom, T)
            par(enl_h, sup, tema, "Enlace (hover) sobre %s" % nom, T)
        # Botón primario
        if claro:
            botones = [("primario", b(700), b(800), "color.base.white"),
                       ("destructivo", e(700), e(800), "color.base.white"),
                       ("de acción crítica", k(700), k(800), "color.base.white")]
        else:
            botones = [("primario", b(300), b(200), b(950)),
                       ("destructivo", e(300), e(200), e(950)),
                       ("de acción crítica", k(300), k(200), k(950))]
        for nombre, relleno, relleno_h, txt in botones:
            par(txt, relleno, tema, "Botón %s (reposo): texto sobre relleno" % nombre, T)
            par(txt, relleno_h, tema, "Botón %s (hover): texto sobre relleno" % nombre, T)
            for sup, nom in ((f, "fondo"), (t, "tarjeta")):
                par(relleno, sup, tema, "Botón %s (reposo): relleno contra %s" % (nombre, nom), U)
                par(relleno_h, sup, tema, "Botón %s (hover): relleno contra %s" % (nombre, nom), U)
        # Botón secundario
        borde = n(600) if claro else n(400)
        hover_sec = n(200) if claro else n(800)
        for sup, nom in ((f, "fondo"), (t, "tarjeta")):
            par(texto, sup, tema, "Botón secundario (reposo): texto sobre %s" % nom, T)
            par(borde, sup, tema, "Botón secundario (reposo y hover): borde contra %s" % nom, U)
        par(texto, hover_sec, tema, "Botón secundario (hover): texto sobre relleno de hover", T)
        # Campos y foco
        foco = b(600) if claro else b(400)
        for sup, nom in ((f, "fondo"), (t, "tarjeta")):
            par(borde, sup, tema, "Borde de campo contra %s" % nom, U)
            par(foco, sup, tema, "Anillo de foco contra %s" % nom, U)
        # Superficie de marca (30 %): elemento activo de barra lateral, fila seleccionada
        activo = b(100) if claro else b(900)
        par(texto, activo, tema, "Elemento activo / fila seleccionada: texto sobre superficie de marca", T)
        par(enl, activo, tema, "Elemento activo: etiqueta e icono de marca sobre superficie de marca", T)
        # Roles
        for r in ROLES + [CRITICO]:
            base = "color.critical" if r["codigo"] == "critical" else "color.role.%s" % r["codigo"]
            fg = "%s.%s" % (base, sfx)
            par(fg, "%s.surface-%s" % (base, sfx), tema, "Insignia de %s: texto sobre su fondo" % r["nombre"].lower(), T)
            for sup, nom in ((f, "fondo"), (t, "tarjeta")):
                par(fg, sup, tema, "Cifra o texto de %s sobre %s" % (r["nombre"].lower(), nom), T)
        # Gráficas, personas y tiendas (objetos gráficos: 3:1)
        for nnum in range(1, 9):
            for sup, nom in ((f, "fondo"), (t, "tarjeta")):
                par("color.chart.%d.%s" % (nnum, sfx), sup, tema, "Serie %d de gráfica contra %s" % (nnum, nom), U)
        for grupo, orden, etiqueta in (("person", ORDEN_PERSONAS, "Persona"), ("store", ORDEN_TIENDAS, "Tienda")):
            for cod in orden:
                nombre = M.cat[grupo][cod]["nombre"]
                for sup, nom in ((f, "fondo"), (t, "tarjeta")):
                    par("color.%s.%s.%s" % (grupo, cod, sfx), sup, tema,
                        "%s %s (punto, anillo o serie) contra %s" % (etiqueta, nombre, nom), U)
    return pares


# =============================================================================
# 5. Métricas
# =============================================================================


def metricas_graficas(M, espacio="lineal"):
    res = {}
    for tema in TEMAS:
        cols = [M.graf[n][tema].hex for n in range(1, 9)]
        cons = []
        for i in range(7):
            a, b = cols[i], cols[i + 1]
            fila = {"par": "%d-%d" % (i + 1, i + 2)}
            for v in ("normal",) + VISTAS:
                fila["dl_" + v] = abs(l_vista(a, v, espacio) - l_vista(b, v, espacio))
            fila["de_normal"] = delta_e(a, b)
            for tcv in TIPOS_CVD:
                fila["de_" + tcv] = delta_e(a, b, tcv, espacio)
            cons.append(fila)
        todos = {"normal": min(delta_e(a, b) for a, b in combinations(cols, 2))}
        for tcv in TIPOS_CVD:
            todos[tcv] = min(delta_e(a, b, tcv, espacio) for a, b in combinations(cols, 2))
        dl_todos = min(abs(l_gris(a) - l_gris(b)) for a, b in combinations(cols, 2))
        res[tema] = {"cols": cols, "cons": cons, "todos": todos, "dl_todos_gris": dl_todos,
                     "min_dl": min(min(f["dl_" + v] for v in VISTAS) for f in cons),
                     "min_de_cvd": min(min(f["de_protanopia"], f["de_deuteranopia"]) for f in cons),
                     "min_de_tritan": min(f["de_tritanopia"] for f in cons),
                     "min_de_normal": min(f["de_normal"] for f in cons)}
    return res


def separacion_roles(M, hx, tema):
    """Para cada rol (y acción crítica): Δmatiz y ΔL; devuelve la peor holgura y detalle."""
    sfx = SUFIJO[tema]
    L, C, H = hex_a_oklch(hx)
    filas = []
    for r in ROLES + [CRITICO]:
        base = "color.critical" if r["codigo"] == "critical" else "color.role.%s" % r["codigo"]
        rh = M.hex("%s.%s" % (base, sfx))
        rl, rc, rhh = hex_a_oklch(rh)
        dm, dl = dif_matiz(H, rhh), abs(L - rl)
        ok = dm >= SEPARACION_MATIZ or dl >= SEPARACION_L
        filas.append({"rol": r["nombre"], "codigo": r["codigo"], "dm": dm, "dl": dl, "ok": ok,
                      "holgura": max(dm / SEPARACION_MATIZ, dl / SEPARACION_L)})
    return filas


def metricas_categoricas(M):
    out = {}
    for grupo, orden in (("person", ORDEN_PERSONAS), ("store", ORDEN_TIENDAS)):
        out[grupo] = {}
        for tema in TEMAS:
            cols = {cod: M.cat[grupo][cod][tema].hex for cod in orden}
            dl_ok = min(abs(hex_a_oklab(cols[a])[0] - hex_a_oklab(cols[b])[0]) for a, b in combinations(orden, 2))
            dl_g = min(abs(l_gris(cols[a]) - l_gris(cols[b])) for a, b in combinations(orden, 2))
            adj = list(zip(orden, orden[1:]))
            sep = {cod: separacion_roles(M, cols[cod], tema) for cod in orden}
            out[grupo][tema] = {
                "cols": cols, "dl_oklab": dl_ok, "dl_gris": dl_g,
                "adj_normal": min(delta_e(cols[a], cols[b]) for a, b in adj),
                "adj_cvd": min(min(delta_e(cols[a], cols[b], "protanopia"), delta_e(cols[a], cols[b], "deuteranopia")) for a, b in adj),
                "adj_tritan": min(delta_e(cols[a], cols[b], "tritanopia") for a, b in adj),
                "sep": sep, "sep_ok": all(f["ok"] for v in sep.values() for f in v),
                "contraste_min": min(min(contraste(c, M.sup(tema, "fondo")), contraste(c, M.sup(tema, "tarjeta"))) for c in cols.values()),
                "croma_min": min(hex_a_oklch(cols[c])[1] for c in orden if c != "shared"),
            }
    return out


def ventana_valor_unico(M):
    """Rango de luminancia que cumple 3:1 contra las cuatro superficies a la vez."""
    y_t100, y_t900 = luminancia(M.neutro[100].hex), luminancia(M.neutro[900].hex)
    y_f50, y_f950 = luminancia(M.neutro[50].hex), luminancia(M.neutro[950].hex)
    y_max = min((y_t100 + 0.05), (y_f50 + 0.05)) / UMBRAL_UI - 0.05
    y_min = UMBRAL_UI * (max(y_t900, y_f950) + 0.05) - 0.05
    return {"y_min": y_min, "y_max": y_max, "l_min": y_min ** (1 / 3), "l_max": y_max ** (1 / 3),
            "l_claro_max": y_max ** (1 / 3), "l_oscuro_min": y_min ** (1 / 3)}


def intervalos_perfil():
    """Intervalo de fracción de croma por paso que sigue reproduciendo los 22 valores del plan."""
    res = {}
    for i, p in enumerate(PASOS):
        ok = []
        for k in range(0, 651):
            fr = k / 500.0
            if (Color(L_PASO[p], CROMA_MARCA * fr, MATIZ_BASE).hex[1:] == PLAN_C["marca"][i]
                    and Color(L_PASO[p], CROMA_NEUTRO * fr, MATIZ_BASE).hex[1:] == PLAN_C["neutro"][i]):
                ok.append(fr)
        res[p] = (min(ok), max(ok)) if ok else None
    return res


def distancias_roles(M, tema):
    sfx = SUFIJO[tema]
    lista = ROLES + [CRITICO]
    filas = []
    for a, b in combinations(lista, 2):
        ba = "color.critical" if a["codigo"] == "critical" else "color.role.%s" % a["codigo"]
        bb = "color.critical" if b["codigo"] == "critical" else "color.role.%s" % b["codigo"]
        ha, hb = M.hex("%s.%s" % (ba, sfx)), M.hex("%s.%s" % (bb, sfx))
        filas.append({"a": a["nombre"], "b": b["nombre"], "dm": dif_matiz(a["matiz"], b["matiz"]),
                      "normal": delta_e(ha, hb), "protanopia": delta_e(ha, hb, "protanopia"),
                      "deuteranopia": delta_e(ha, hb, "deuteranopia"), "tritanopia": delta_e(ha, hb, "tritanopia")})
    return filas


REFUERZO_ROLES = {
    frozenset(("expense", "error")): "Contexto y forma: gasto solo como cifra con signo «−»; error solo como mensaje con icono de alerta y texto (K26).",
    frozenset(("income", "success")): "Ingreso siempre con «+» y flecha de entrada en cifras; éxito siempre con palomita y texto en avisos.",
}
# Hallazgo R3-04 (versión 1.1): la marca entra al barrido de pares cercanos de §7.7 con el paso de texto de cada tema.
PASO_MARCA = {"claro": 700, "oscuro": 400}
REFUERZO_MARCA = ("Marca frente a categórica (R3-04): la marca solo vive en controles (botón primario, enlace subrayado, anillo "
                  "de foco, elemento activo) y nunca marca a una persona o tienda; la categórica siempre va como punto, anillo de "
                  "avatar o marca de serie junto a su nombre en texto neutro (forma y etiqueta).")
# Versión 1.0 de lo que cambia en la 1.1 (K26 y bloqueo U2a-walmart-oscuro): matiz de gasto; matiz, croma, L claro, L oscuro.
ANTES_11 = {"expense": 25.0,
            ("person", "karen"): (50.0, 0.12, 0.64, 0.725),
            ("store", "bodega-aurrera"): (50.0, 0.12, 0.41, 0.83),
            ("store", "walmart"): (262.0, 0.12, 0.515, 0.72)}
GASTO_RANGO = range(40, 46)   # K26: terracota 40–45°, grados enteros

ESTRUCTURALES = [  # descripción, primer plano, fondo, (mín, máx) del plan, texto del plan
    ("Texto 900 / fondo 50", "color.neutral.900", "color.neutral.50", (11.4, 11.4), "11.4"),
    ("Texto secundario 700 / fondo 50", "color.neutral.700", "color.neutral.50", (5.3, 5.3), "5.3"),
    ("Neutro 600 / tarjeta 100 (solo texto grande o iconos; aquí no se usa para texto)", "color.neutral.600", "color.neutral.100", (3.45, 3.45), "3.45"),
    ("Blanco / marca 600 (no se usa para texto de botón)", "color.base.white", "color.brand.600", (3.7, 4.0), "3.7–4.0 (3.98 en decisiones §I)"),
    ("Marca 700 / fondo 50", "color.brand.700", "color.neutral.50", (5.1, 5.5), "5.1–5.5"),
    ("Blanco / marca 700 (texto del botón primario)", "color.base.white", "color.brand.700", (5.8, 5.8), "5.8 (criterios U2-07)"),
    ("Texto 100 / fondo 950", "color.neutral.100", "color.neutral.950", (14.1, 14.1), "14.1"),
    ("Texto 300 / tarjeta 900", "color.neutral.300", "color.neutral.900", (7.9, 7.9), "7.9"),
    ("Marca 300 / fondo 950", "color.brand.300", "color.neutral.950", (10.4, 10.7), "10.4–10.7"),
    ("Marca 950 / marca 300 (texto del botón primario en oscuro)", "color.brand.950", "color.brand.300", (10.4, 10.4), "10.4 (criterios U2-07)"),
    ("Divisor 200 / fondo 50 (decorativo)", "color.neutral.200", "color.neutral.50", (1.2, 1.2), "1.2"),
    ("Divisor 800 / fondo 950 (decorativo)", "color.neutral.800", "color.neutral.950", (1.9, 1.9), "1.9"),
    ("Borde de campo 600 / fondo 50", "color.neutral.600", "color.neutral.50", (4.2, 99.0), "≥ 4.2"),
]


def entidades(M, tema):
    sfx = SUFIJO[tema]
    out = [("rol", r["codigo"], r["nombre"], M.hex("%s.%s" % (base_rol(r["codigo"]), sfx))) for r in ROLES + [CRITICO]]
    for grupo, orden, et in (("person", ORDEN_PERSONAS, "persona"), ("store", ORDEN_TIENDAS, "tienda")):
        out += [(et, cod, M.cat[grupo][cod]["nombre"], M.cat[grupo][cod][tema].hex) for cod in orden]
    return out


def entidades_con_marca(M, tema):
    """Entidades de §7.7 más la marca en su paso de texto (700 en claro, 400 en oscuro; hallazgo R3-04)."""
    p = PASO_MARCA[tema]
    return entidades(M, tema) + [("marca", "brand", "Marca %d" % p, M.hex("color.brand.%d" % p))]


def pares_cercanos(M):
    """Pares (roles, acción crítica, personas, tiendas y marca) con Δ matiz < 20° y ΔL < 0.15, con su refuerzo."""
    filas = []
    for tema in TEMAS:
        for (ga, ca, na, ha), (gb, cb, nb, hb) in combinations(entidades_con_marca(M, tema), 2):
            La, _, Ha = hex_a_oklch(ha)
            Lb, _, Hb = hex_a_oklch(hb)
            dm, dl = dif_matiz(Ha, Hb), abs(La - Lb)
            if dm >= SEPARACION_MATIZ or dl >= SEPARACION_L:
                continue
            if "marca" in (ga, gb):
                ref = "SIN RESOLVER" if "rol" in (ga, gb) else REFUERZO_MARCA
            elif ga == "rol" and gb == "rol":
                ref = REFUERZO_ROLES.get(frozenset((ca, cb)), "Icono y texto propios de cada rol; nunca comparten componente.")
            elif ga == gb:
                ref = "Mismo grupo: orden fijo, etiqueta directa y ΔL ≥ 0.10 dentro del grupo."
            elif "rol" in (ga, gb):
                ref = "SIN RESOLVER"
            else:
                ref = "Grupos distintos (persona y tienda): nunca comparten gráfica ni leyenda; en tablas van en columnas distintas y con su nombre."
            if "ximena" in (ca, cb) and {ca, cb} - {"ximena"} <= {c for _, c in ROSAS_XIMENA}:
                ref += " **PROVISIONAL**: bloqueo `%s` (¿basta este refuerzo frente a «≥ 20° o ΔL ≥ 0.15» de K26?)." % BLOQUEO_ROSAS
            filas.append({"tema": tema, "a": "%s (%s)" % (na, ga), "b": "%s (%s)" % (nb, gb), "dm": dm, "dl": dl, "refuerzo": ref})
    return filas


def _refs_roles_marca(M, tema):
    """(L, h) medidos del texto de los 8 roles, acción crítica y la marca (paso de texto) en un tema."""
    hexes = [M.hex("%s.%s" % (base_rol(r["codigo"]), SUFIJO[tema])) for r in ROLES + [CRITICO]]
    hexes.append(M.hex("color.brand.%d" % PASO_MARCA[tema]))
    return [hex_a_oklch(hx)[0::2] for hx in hexes]


def _holgura(L, H, refs):
    """Holgura mínima de la regla ≥ 20° o ΔL ≥ 0.15 frente a una lista de (L, h); ≥ 1 cumple."""
    return min(max(dif_matiz(H, rh) / SEPARACION_MATIZ, abs(L - rl) / SEPARACION_L) for rl, rh in refs)


def barrido_gasto(M):
    """K26: para cada grado entero de 40° a 45°, distancias medidas del texto de gasto a error, advertencia y categóricas.
    Gana el mayor Δ matiz mínimo a error y advertencia (ambos temas); empates, por ΔE normal y luego por ΔE con daltonismo."""
    filas = []
    for h in GASTO_RANGO:
        esc = escala(float(h), CROMA_ROL)
        f = {"h": h, "temas": {}}
        for tema in TEMAS:
            p, _ = M._elegir_paso("expense", tema, esc)
            hx = esc[p].hex
            L, _, H = hex_a_oklch(hx)
            d = {}
            for cod in ("error", "warning"):
                ref = M.hex("color.role.%s.%s" % (cod, SUFIJO[tema]))
                d[cod] = (dif_matiz(H, hex_a_oklch(ref)[2]), delta_e(hx, ref), min(delta_e(hx, ref, t) for t in TIPOS_CVD))
            cat = min((_holgura(L, H, [hex_a_oklch(c)[0::2]]), n) for g, _, n, c in entidades(M, tema) if g != "rol")
            f["temas"][tema] = {"paso": p, "hex": hx, "h": H, "d": d, "cat": cat}
        for k, i in (("dm", 0), ("de", 1), ("cvd", 2)):
            f[k] = min(f["temas"][t]["d"][c][i] for t in TEMAS for c in ("error", "warning"))
        f["dm_error"] = min(f["temas"][t]["d"]["error"][0] for t in TEMAS)
        f["dm_warning"] = min(f["temas"][t]["d"]["warning"][0] for t in TEMAS)
        f["cat"] = min(f["temas"][t]["cat"] for t in TEMAS)
        filas.append(f)
    mejor = max(filas, key=lambda f: (round(f["dm"], 6), round(f["de"], 6), round(f["cvd"], 6)))
    return filas, mejor


def ventanas_matiz(M, grupo, cod, otros=None, propio=None):
    """Matices enteros en los que la categórica `cod`, con su croma y sus L finales (o `propio` = (croma, {tema: L})),
    cumple las reglas de §7: ≥ 20° o ΔL ≥ 0.15 frente a roles, acción crítica y marca; contraste ≥ 3:1 sobre fondo y
    tarjeta; croma ≥ 0.10; fuera de la zona reservada; ΔL ≥ 0.10 (OKLab y gris) frente a su grupo. `otros` =
    {tema: {código: hex}} sustituye el color de otros miembros. Devuelve los matices y sus ventanas contiguas."""
    e = M.cat[grupo][cod]
    c_ped, ls = propio if propio else (e["c"], {t: e[t].l for t in TEMAS})
    orden = ORDEN_PERSONAS if grupo == "person" else ORDEN_TIENDAS
    refs = {t: _refs_roles_marca(M, t) for t in TEMAS}
    vecinos = {t: [((otros or {}).get(t) or {}).get(o, M.cat[grupo][o][t].hex) for o in orden if o != cod] for t in TEMAS}
    vecinos = {t: [(hex_a_oklab(v)[0], l_gris(v)) for v in vs] for t, vs in vecinos.items()}
    sup = {t: (M.sup(t, "fondo"), M.sup(t, "tarjeta")) for t in TEMAS}
    ok = []
    for h in range(360):
        bien = True
        for t in TEMAS:
            hx = Color(ls[t], c_ped, float(h)).hex
            L, C, H = hex_a_oklch(hx)
            g = l_gris(hx)
            if not (dist_reserva(hx, float(h)) >= ZONA_RESERVADA and C >= CROMA_PISO_GRAFICAS
                    and min(contraste(hx, x) for x in sup[t]) >= UMBRAL_UI and _holgura(L, H, refs[t]) >= 1.0
                    and all(abs(L - vl) >= UMBRAL_DL and abs(g - vg) >= UMBRAL_DL for vl, vg in vecinos[t])):
                bien = False
                break
        if bien:
            ok.append(h)
    vent = []
    for h in ok:
        if vent and h == vent[-1][1] + 1:
            vent[-1][1] = h
        else:
            vent.append([h, h])
    return ok, vent


def eleccion_karen_bodega(M, vk, vb):
    """K26: reparto de Karen y Bodega Aurrera en sus ventanas, entre deuda y error. Criterios, en orden: (1) el mayor
    hueco mínimo de matiz entre deuda, las dos y error; (2) Karen lo más lejos posible de Ximena (su grupo); (3) la mayor
    holgura mínima de las dos frente a roles, acción crítica y marca. El par cumple la regla entre sí en ambos temas."""
    hd = next(r["matiz"] for r in ROLES if r["codigo"] == "debt")
    he = next(r["matiz"] for r in ROLES if r["codigo"] == "error") + 360.0
    ek, eb = M.cat["person"]["karen"], M.cat["store"]["bodega-aurrera"]
    refs = {t: _refs_roles_marca(M, t) for t in TEMAS}

    def medir(e, h):
        out = {}
        for t in TEMAS:
            L, _, H = hex_a_oklch(Color(e[t].l, e["c"], float(h)).hex)
            out[t] = (L, H, _holgura(L, H, refs[t]))
        return out
    mk, mb = {a: medir(ek, a) for a in vk}, {b: medir(eb, b) for b in vb}
    desp = lambda h: h + 360.0 if h < 180 else float(h)
    cands = []
    for a in vk:
        for b in vb:
            pa, pb = desp(a), desp(b)
            if not (hd < pa < he and hd < pb < he):
                continue
            if any(dif_matiz(mk[a][t][1], mb[b][t][1]) < SEPARACION_MATIZ and abs(mk[a][t][0] - mb[b][t][0]) < SEPARACION_L for t in TEMAS):
                continue
            lo, hi = sorted((pa, pb))
            huecos = (lo - hd, hi - lo, he - hi)
            cands.append({"a": a, "b": b, "huecos": huecos, "min": min(huecos), "xim": dif_matiz(a, M.cat["person"]["ximena"]["h"]),
                          "hol": min(min(mk[a][t][2], mb[b][t][2]) for t in TEMAS)})
    cands.sort(key=lambda c: (-c["min"], -c["xim"], -round(c["hol"], 9), c["a"], c["b"]))
    return cands


def tokens_antes_11(M):
    """Tokens de la versión 1.0 que cambian en la 1.1, recalculados con el mismo método a partir de ANTES_11."""
    antes = {}
    esc = escala(ANTES_11["expense"], CROMA_ROL)
    base = "color.role.expense"
    antes[base + ".light"] = esc[M._elegir_paso("expense", "claro", esc)[0]]
    antes[base + ".dark"] = esc[M._elegir_paso("expense", "oscuro", esc)[0]]
    antes[base + ".surface-light"], antes[base + ".surface-dark"] = esc[100], esc[900]
    for p in PASOS:
        antes["%s.%d" % (base, p)] = esc[p]
    for clave, v in ANTES_11.items():
        if clave != "expense":
            h, c, lc, lo = v
            antes["color.%s.%s.light" % clave] = Color(lc, c, h)
            antes["color.%s.%s.dark" % clave] = Color(lo, c, h)
    return {r: col for r, col in antes.items() if col.hex != M.hex(r)}


def dist_reserva(hx, h_diseno=None):
    """Distancia de matiz (°) al eje opuesto: la menor entre la medida en el hex y la de diseño."""
    d = dif_matiz(hex_a_oklch(hx)[2], MATIZ_OPUESTO)
    return d if h_diseno is None else min(d, dif_matiz(h_diseno, MATIZ_OPUESTO))


def cercanos_reserva(M, n=5):
    """Colores de identidad (roles, personas, tiendas y series; sin acción crítica) más cercanos al opuesto."""
    dist = {}
    for tema in TEMAS:
        filas = [(g, nombre, hx) for g, cod, nombre, hx in entidades(M, tema) if cod != "critical"]
        filas += [("serie %d" % k, M.graf[k]["nombre"], M.graf[k][tema].hex) for k in range(1, 9)]
        for g, nombre, hx in filas:
            dist[(g, nombre)] = min(dist.get((g, nombre), 360.0), dist_reserva(hx))
    return sorted(dist.items(), key=lambda kv: (kv[1], kv[0]))[:n]


MOTIVOS_DESCARTE = {"croma": "croma final < 0.10", "zona": "matiz medido dentro de la zona reservada",
                    "contraste": "contraste < 3:1", "dl": "ΔL con una vecina < 0.10",
                    "de": "ΔE de vecinas bajo el piso"}


def _controles_serie(M, n, h, clase):
    """Controles de §8 para la serie n con matiz h y las otras siete en su lugar, en ambos temas.
    Devuelve (motivos de descarte, métricas)."""
    motivos = set()
    met = {"dh": 360.0, "dl": 9.0, "de": 999.0, "hex": {}, "croma": {}}
    for tema in TEMAS:
        col = Color(L_CLASE[(tema, clase)], CROMA_GRAFICA[tema], h)
        cols = [M.graf[k][tema].hex for k in range(1, 9)]
        cols[n - 1] = col.hex
        met["hex"][tema], met["croma"][tema] = col.hex, col.c
        if col.c < CROMA_PISO_GRAFICAS:
            motivos.add("croma")
        if dist_reserva(col.hex, h) < ZONA_RESERVADA:
            motivos.add("zona")
        if min(contraste(col.hex, M.sup(tema, s)) for s in ("fondo", "tarjeta")) < UMBRAL_UI:
            motivos.add("contraste")
        for i in (n - 2, n - 1):
            if i < 0 or i > 6:
                continue
            a, b = cols[i], cols[i + 1]
            dl_lin = min(abs(l_vista(a, v) - l_vista(b, v)) for v in VISTAS)
            dl_gam = min(abs(l_vista(a, v, "gamma") - l_vista(b, v, "gamma")) for v in VISTAS)
            de_n = delta_e(a, b)
            de_c = min(delta_e(a, b, "protanopia"), delta_e(a, b, "deuteranopia"))
            met["dl"], met["de"] = min(met["dl"], dl_lin), min(met["de"], de_n)
            if min(dl_lin, dl_gam) < UMBRAL_DL:
                motivos.add("dl")
            if de_n < NORMAL_PISO or de_c < CVD_OBJETIVO:
                motivos.add("de")
        hm = hex_a_oklch(col.hex)[2]
        met["dh"] = min([met["dh"]] + [dif_matiz(hm, hex_a_oklch(cols[j])[2]) for j in range(8) if j != n - 1])
    return motivos, met


def barrido_serie(M):
    """Barrido de §8.8: matices múltiplos de PASO_BARRIDO para la serie SERIE_CORREGIDA (con su clase) y las
    otras siete en su lugar. Ordena los que pasan todos los controles por la mayor distancia de matiz a las
    otras series y, después, por el mayor ΔL con las vecinas."""
    n = SERIE_CORREGIDA
    clase = M.graf[n]["clase"]
    pasan, descartados, motivos_cuenta, en_zona = [], [], {}, 0
    for h in range(0, 360, PASO_BARRIDO):
        if dif_matiz(h, MATIZ_OPUESTO) < ZONA_RESERVADA:
            en_zona += 1
            continue
        motivos, met = _controles_serie(M, n, float(h), clase)
        met["h"] = h
        for m in motivos:
            motivos_cuenta[m] = motivos_cuenta.get(m, 0) + 1
        if motivos:
            descartados.append((met, motivos))
        else:
            pasan.append(met)
    pasan.sort(key=lambda d: (-d["dh"], -d["dl"], d["h"]))
    descartados.sort(key=lambda dm: (-dm[0]["dh"], dm[0]["h"]))
    return {"pasan": pasan, "descartados": descartados, "motivos": motivos_cuenta, "en_zona": en_zona,
            "evaluados": len(pasan) + len(descartados)}


def _metricas_orden(M, orden):
    """Métricas de un orden (lista de números de serie vigentes) en ambos temas."""
    dl, de, t_n, t_c = 9.0, 999.0, 999.0, 999.0
    tres, par12 = {}, {}
    for tema in TEMAS:
        cols = [M.graf[k][tema].hex for k in orden]
        for i in range(7):
            dl = min(dl, min(abs(l_vista(cols[i], v) - l_vista(cols[i + 1], v)) for v in VISTAS))
            de = min(de, delta_e(cols[i], cols[i + 1]))
        pares3 = list(combinations(range(3), 2))
        tres[tema] = min((delta_e(cols[i], cols[j]), orden[i], orden[j]) for i, j in pares3)
        t_n = min(t_n, tres[tema][0])
        t_c = min([t_c] + [min(delta_e(cols[i], cols[j], "protanopia"), delta_e(cols[i], cols[j], "deuteranopia")) for i, j in pares3])
        par12[tema] = delta_e(cols[0], cols[1])
    return {"orden": list(orden), "dl": dl, "de": de, "tres_n": t_n, "tres_cvd": t_c, "tres": tres, "par12": par12,
            "clave": (dl, de, t_n)}


def ranking_orden(M):
    """Búsqueda exhaustiva de §8.2 con los matices vigentes: permutaciones que alternan clases y dejan las
    series FIJAS_ORDEN en su lugar. Clave: mínimo ΔL consecutivo (gris y tres simulaciones, ambos temas),
    después mínimo ΔE normal entre vecinas y, para desempatar, mínimo ΔE normal entre todos los pares de las
    tres primeras series (las únicas que se usan en dispersión, mapas y pequeños múltiplos, §8.5)."""
    libres = [k for k in range(1, 9) if k not in FIJAS_ORDEN]
    pos = {}
    for k in libres:
        pos.setdefault(M.graf[k]["clase"], []).append(k)
    clases = sorted(pos)
    lista = []
    for perms in product(*[list(permutations(pos[c])) for c in clases]):
        orden = list(range(1, 9))
        for c, perm in zip(clases, perms):
            for destino, origen in zip(pos[c], perm):
                orden[destino - 1] = origen
        lista.append(_metricas_orden(M, orden))
    lista.sort(key=lambda d: (tuple(-x for x in d["clave"]), d["orden"]))
    por_nombre = {M.graf[k]["nombre"]: k for k in range(1, 9)}
    r0 = _metricas_orden(M, [SERIE_CORREGIDA if n is None else por_nombre[n] for n in ORDEN_RONDA0])
    act = _metricas_orden(M, list(range(1, 9)))
    puesto = lambda d: 1 + sum(1 for x in lista if x["clave"] > d["clave"])
    empates = [x for x in lista if x["orden"] != act["orden"] and x["dl"] == act["dl"] and x["de"] == act["de"]]
    return {"lista": lista, "actual": act, "ronda0": r0, "puesto_actual": puesto(act), "puesto_ronda0": puesto(r0),
            "empates": empates}


# =============================================================================
# 6. Escritura: CSV y JSON
# =============================================================================


def escribir_csv(pares):
    buf = io.StringIO()
    w = csv.writer(buf, lineterminator="\n")
    w.writerow(["primer_plano", "fondo", "tema", "uso", "ratio", "umbral", "cumple", "primer_plano_hex", "fondo_hex"])
    for p in pares:
        w.writerow([p["primer_plano"], p["fondo"], p["tema"], p["uso"], fmt_ratio(p["ratio"]),
                    "%.1f" % p["umbral"], "true" if p["cumple"] else "false", p["fp_hex"], p["fo_hex"]])
    SALIDA["csv"].write_text(buf.getvalue(), encoding="utf-8")


def tokens_json(M):
    raiz = {"color": {"$description": M.grupos["color"]}}
    for ruta, t in M.tokens.items():
        partes = ruta.split(".")
        nodo = raiz
        for i, parte in enumerate(partes[:-1]):
            sub = ".".join(partes[:i + 1])
            if parte not in nodo:
                nodo[parte] = {}
                if sub in M.grupos and sub != "color":
                    nodo[parte]["$description"] = M.grupos[sub]
            nodo = nodo[parte]
        hoja = {"$type": "color", "$value": t["hex"]}
        if t["desc"]:
            hoja["$description"] = t["desc"]
        hoja["$extensions"] = {"oklch": {"l": round(t["l"], 6), "c": round(t["c"], 6), "h": round(t["h"], 2)}}
        nodo[partes[-1]] = hoja
    return raiz


def escribir_json(M):
    SALIDA["json"].write_text(json.dumps(tokens_json(M), ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


# =============================================================================
# 7. Escritura: Markdown
# =============================================================================


def tabla(encabezados, filas):
    lineas = ["| " + " | ".join(encabezados) + " |", "|" + "|".join("---" for _ in encabezados) + "|"]
    for f in filas:
        lineas.append("| " + " | ".join(str(x) for x in f) + " |")
    return "\n".join(lineas)


def r2(x):
    return fmt_ratio(x)


def base_rol(cod):
    return "color.critical" if cod == "critical" else "color.role.%s" % cod


def rosas_ximena(M):
    """Ronda 1 de revisión U2a (H-01): Karen y Bodega Aurrera frente a Ximena en cada tema, el par persona–tienda más
    parecido de la 1.0 y la comprobación de que ningún matiz de sus ventanas las separa de Ximena en ambos temas.
    Aborta si los tokens marcados PROVISIONAL con BLOQUEO_ROSAS no son exactamente los que quedan bajo la regla."""
    if getattr(M, "_rosas", None):
        return M._rosas
    ex = M.cat["person"]["ximena"]

    def lejos(a, b):
        (La, _, Ha), (Lb, _, Hb) = hex_a_oklch(a), hex_a_oklch(b)
        return dif_matiz(Ha, Hb) >= SEPARACION_MATIZ or abs(La - Lb) >= SEPARACION_L
    R = {"pares": {}, "bajo": set()}
    for g, c in ROSAS_XIMENA:
        e = M.cat[g][c]
        for t in TEMAS:
            (La, _, Ha), (Lx, _, Hx) = hex_a_oklch(e[t].hex), hex_a_oklch(ex[t].hex)
            R["pares"][(c, t)] = {"hex": e[t].hex, "hx": ex[t].hex, "dm": dif_matiz(Ha, Hx), "dl": abs(La - Lx),
                                  "de": delta_e(e[t].hex, ex[t].hex)}
            if not lejos(e[t].hex, ex[t].hex):
                R["bajo"].add("color.%s.%s.%s" % (g, c, SUFIJO[t]))
        v, _ = ventanas_matiz(M, g, c)
        sale = [h for h in v if all(lejos(Color(e[t].l, e["c"], float(h)).hex, ex[t].hex) for t in TEMAS)]
        if sale:
            raise SystemExit("Bloqueo %s: %s se separa de Ximena en %s; revisar el texto y el bloqueo" % (BLOQUEO_ROSAS, e["nombre"], sale))
    marcados = {r for r, nota in PROVISIONAL.items() if BLOQUEO_ROSAS in nota}
    if R["bajo"] != marcados:
        raise SystemExit("Bloqueo %s: bajo la regla frente a Ximena %s; marcados PROVISIONAL %s" % (BLOQUEO_ROSAS, sorted(R["bajo"]), sorted(marcados)))
    antes = tokens_antes_11(M)
    h10 = lambda r: antes[r].hex if r in antes else M.hex(r)
    ref = []
    for t in TEMAS:
        for a in ORDEN_PERSONAS:
            for b in ORDEN_TIENDAS:
                ha, hb = h10("color.person.%s.%s" % (a, SUFIJO[t])), h10("color.store.%s.%s" % (b, SUFIJO[t]))
                if not lejos(ha, hb):
                    ref.append((delta_e(ha, hb), M.cat["person"][a]["nombre"], M.cat["store"][b]["nombre"], t))
    R["ref10"] = min(ref)
    M._rosas = R
    return R


def _rosas_cifras(M):
    R = rosas_ximena(M)
    p = R["pares"]
    return R, p[("bodega-aurrera", "oscuro")], p[("bodega-aurrera", "claro")], p[("karen", "claro")], p[("karen", "oscuro")]


def md_rosas_ximena(M):
    """Párrafo de «Cambios 1.1», punto 2: Karen y Bodega Aurrera frente a Ximena (bloqueo BLOQUEO_ROSAS)."""
    R, bo, bc, kc, ko = _rosas_cifras(M)
    de0, a0, b0, t0 = R["ref10"]
    return ("**Frente a Ximena (ronda 1 de revisión, H-01; PROVISIONAL).** En oscuro, Bodega Aurrera (`%s`) y Ximena (`%s`) quedan a %.1f°, ΔL %.3f y **ΔE %.1f**: a simple vista son el mismo color. En claro quedan a %.1f°, ΔL %.3f y ΔE %.1f. Karen y Ximena quedan a %.1f° y %.1f°, con ΔL %.3f y %.3f y ΔE %.1f y %.1f (claro y oscuro). Bodega Aurrera–Ximena es una colisión nueva de la 1.1: en la 1.0, el par persona–tienda más parecido era %s–%s en %s, con ΔE %.1f. Con la L y el croma de cada una, ningún matiz de sus ventanas (tabla de arriba) separa a Karen ni a Bodega Aurrera de Ximena por ≥ 20° o ΔL ≥ 0.15 en ambos temas; el script lo comprueba. §7.7 declara el refuerzo, y U2-08 lo admite. Pero U2-02 y la respuesta de Axel en `U2-gasto-error` piden moverlas «para conservar separación ≥ 20° o ΔL ≥ 0.15» frente a las demás entradas. Como el criterio admite dos lecturas, el bloqueo `docs/estado/bloqueos/%s.md` lleva las opciones para Axel. Mientras decide, van **PROVISIONAL** los %d tokens que quedan bajo la regla frente a Ximena (%s), con la nota en su descripción en `C2-colores.json`.\n" % (
        bo["hex"], bo["hx"], bo["dm"], bo["dl"], bo["de"], bc["dm"], bc["dl"], bc["de"], kc["dm"], ko["dm"], kc["dl"], ko["dl"],
        kc["de"], ko["de"], a0, b0, t0, de0, BLOQUEO_ROSAS, len(R["bajo"]), ", ".join("`%s`" % r for r in sorted(R["bajo"]))))


def hallazgo_rosas_ximena(M):
    """Hallazgo de §12 para el bloqueo BLOQUEO_ROSAS (ronda 1 de revisión U2a, H-01)."""
    R, bo, bc, kc, ko = _rosas_cifras(M)
    return ("Ronda 1 de revisión U2a (H-01): en oscuro, Bodega Aurrera y Ximena quedan a %.1f°, ΔL %.3f y ΔE %.1f (a simple vista, el mismo color; en la 1.0, el par persona–tienda más parecido tenía ΔE %.1f), y en claro a %.1f°, ΔL %.3f y ΔE %.1f. Karen y Ximena quedan a %.1f° y %.1f°, con ΔL %.3f y %.3f. Con las L actuales, ningún matiz de sus ventanas las separa de Ximena por ≥ 20° o ΔL ≥ 0.15. §7.7 declara el refuerzo, pero U2-02 admite dos lecturas: bloqueo `%s`, con %d tokens PROVISIONALES (Cambios 1.1, punto 2)." % (
        bo["dm"], bo["dl"], bo["de"], R["ref10"][0], bc["dm"], bc["dl"], bc["de"], kc["dm"], ko["dm"], kc["dl"], ko["dl"], BLOQUEO_ROSAS, len(R["bajo"])))


def md_cambios_11(M, G):
    """Sección «Cambios 1.1» (puerta 1): K26, K27 y hallazgo R3-04, con valores anteriores y nuevos y el porqué."""
    o = []
    ap = o.append
    h_gasto = next(r["matiz"] for r in ROLES if r["codigo"] == "expense")
    filas_g, mejor = barrido_gasto(M)
    if mejor["h"] != h_gasto:
        raise SystemExit("K26: el barrido de gasto elige %d° y ROLES tiene %s°" % (mejor["h"], h_gasto))
    ek, eb, ew = M.cat["person"]["karen"], M.cat["store"]["bodega-aurrera"], M.cat["store"]["walmart"]
    vk, wk = ventanas_matiz(M, "person", "karen")
    vb, wb = ventanas_matiz(M, "store", "bodega-aurrera")
    aw, ab = ANTES_11[("store", "walmart")], ANTES_11[("store", "bodega-aurrera")]
    w_antes = {"oscuro": {"walmart": Color(aw[3], aw[1], aw[0]).hex}}
    _, wb0 = ventanas_matiz(M, "store", "bodega-aurrera", w_antes)
    _, wb00 = ventanas_matiz(M, "store", "bodega-aurrera", w_antes, (ab[1], {"claro": ab[2], "oscuro": ab[3]}))
    cands = eleccion_karen_bodega(M, vk, vb)
    if not cands or (cands[0]["a"], cands[0]["b"]) != (ek["h"], eb["h"]):
        raise SystemExit("K26: el reparto de Karen y Bodega Aurrera no coincide con PERSONAS/TIENDAS")
    desp = lambda h: h + 360 if h < 180 else h
    hd = next(r["matiz"] for r in ROLES if r["codigo"] == "debt")
    he = next(r["matiz"] for r in ROLES if r["codigo"] == "error")
    if any(not (hd < desp(h) < he + 360) for h in vk + vb) or wb0 or wb00:
        raise SystemExit("Cambios 1.1: las ventanas cambiaron; revisar el texto de la sección")
    l_rol = {t: L_PASO[M.paso_rol["expense"][t]] for t in TEMAS}
    if any(abs(e[t].l - l_rol[t]) >= SEPARACION_L for e in (ek, eb) for t in TEMAS):
        raise SystemExit("Cambios 1.1: Karen o Bodega Aurrera ya se separan de los roles por luminosidad; revisar el texto")
    antes = tokens_antes_11(M)
    vtxt = lambda w: ", ".join(("%d°" % a) if a == b else "%d°–%d°" % (a, b) for a, b in w) or "ninguna"

    def fmt(hx):
        L, C, H = hex_a_oklch(hx)
        return "`%s` · %.1f° · L %.3f · C %.3f" % (hx, H, L, C)

    ap("## Cambios 1.1 (puerta 1, %s)\n" % FECHA_11)
    ap("Aplica las respuestas de Axel en la puerta 1 (`docs/plan/decisiones.md` §K): **K26** (gasto pasa a terracota 40–45° y error se queda en 20°; Karen y Bodega Aurrera, que estaban en 50°, se mueven en la paleta categórica; la marca entra al barrido de pares cercanos, hallazgo R3-04 de la revisión 3) y **K27** (paleta de gráficas: ΔL ≥ 0.10 entre series consecutivas con forma, patrón y etiqueta de refuerzo; todos los pares se reportan como información). Quedan resueltos los bloqueos `U2-gasto-error` y `criterios-deltal`. No cambian el método (L por paso, perfil de croma, umbrales), la marca, el neutro, los demás roles ni las 8 series de gráficas. Todo lo regenera el script y su verificación final lo comprueba.\n")
    ap("### Parámetros que cambian\n")
    filas = [["Gasto (`color.role.expense`)", "matiz %d°" % ANTES_11["expense"], "matiz %d°" % h_gasto, "K26"]]
    for clave, (h, c, lc, lo) in [(k, v) for k, v in ANTES_11.items() if k != "expense"]:
        e = M.cat[clave[0]][clave[1]]
        filas.append(["%s (`color.%s.%s`)" % ((e["nombre"],) + clave), "%d° · croma %.2f · L %.3f / %.3f" % (h, c, lc, lo),
                      "%d° · croma %.2f · L %.3f / %.3f" % (e["h"], e["c"], e["claro"].l, e["oscuro"].l),
                      "**PROVISIONAL**, bloqueo `U2a-walmart-oscuro`" if clave == ("store", "walmart") else
                      ("K26; **PROVISIONAL** frente a Ximena, bloqueo `%s`" % BLOQUEO_ROSAS) if clave in ROSAS_XIMENA else "K26"])
    ap(tabla(["Elemento", "1.0: matiz · croma pedido · L claro / oscuro", "1.1", "Origen"], filas))
    ap("\n### Tokens que cambian (medidos en el hex)\n")
    ap(tabla(["Token", "1.0: hex · matiz · L · croma", "1.1: hex · matiz · L · croma"], [["`%s`" % r, fmt(col.hex), fmt(M.hex(r))] for r, col in antes.items()]))
    ap("\nLos otros %d tokens no cambian de valor. Tablas regeneradas: §5.3 (escala de gasto), §6.1, §6.5, §6.6, §7.1 a §7.5, §7.7, §10, §12 y el Anexo A, además de `contraste.csv`, `C2-colores.json`, `paleta.html` y `paleta.png`.\n" % (len(M.tokens) - len(antes)))

    ap("### 1. Gasto: el grado de 40° a 45° más separado (K26)\n")
    ap("Para cada grado entero del rango el script construye la escala de gasto con el mismo método (L por paso, perfil de croma, croma máx. %.2f), elige el paso de texto con la misma regla de 4.5:1 y mide en el hex el Δ matiz a error (20°) y a advertencia (75°), el ΔE normal y el ΔE mínimo con las tres simulaciones de daltonismo frente a esos dos roles, y la holgura mínima frente a las 8 categóricas (máx(Δ matiz / 20°, ΔL / 0.15); cumple con ≥ 1). Gana el mayor Δ matiz mínimo; a igualdad, el mayor ΔE normal y después el mayor ΔE con daltonismo.\n" % CROMA_ROL)
    filas = []
    for f in filas_g:
        tc, to = f["temas"]["claro"], f["temas"]["oscuro"]
        filas.append([("**%d°**" if f is mejor else "%d°") % f["h"], "%d · `%s` · %.1f°" % (tc["paso"], tc["hex"], tc["h"]),
                      "%d · `%s` · %.1f°" % (to["paso"], to["hex"], to["h"]), "%.1f°" % f["dm_error"], "%.1f°" % f["dm_warning"],
                      "%.1f" % f["de"], "%.1f" % f["cvd"], "%.2f (%s)" % f["cat"]])
    ap(tabla(["Matiz", "Claro: paso · hex · matiz medido", "Oscuro: paso · hex · matiz medido", "Δ matiz a error (mín.)", "Δ matiz a advertencia (mín.)", "ΔE normal mín.", "ΔE daltonismo mín.", "Holgura frente a categóricas"], filas))
    adv_lejos = all(f["dm_warning"] > f["dm_error"] for f in filas_g)
    ap("\n**Se elige %d°.** %sEn %d° gasto queda a %.1f° de error y a %.1f° de advertencia, con ΔE normal mínimo %.1f y ΔE mínimo con daltonismo %.1f, los mayores del rango. Frente a las categóricas la holgura va de %.2f a %.2f en todo el rango y la fija %s, así que no decide. El paso de texto sigue siendo 700 en claro y 400 en oscuro, con 4.5:1 (§6.1).\n" % (
        mejor["h"], "Dentro del rango, subir el matiz aleja a gasto de error y lo acerca a advertencia, pero advertencia queda siempre más lejos que error, así que el Δ matiz mínimo crece con el matiz y el mejor grado es el extremo. " if adv_lejos else "",
        mejor["h"], mejor["dm_error"], mejor["dm_warning"], mejor["de"], mejor["cvd"], min(f["cat"][0] for f in filas_g), max(f["cat"][0] for f in filas_g), mejor["cat"][1]))

    ap("### 2. Karen y Bodega Aurrera (K26)\n")
    ap("Con gasto en %d° el rango cálido queda cerrado: entre error 20°, gasto %d° y advertencia 75° no hay ningún matiz a ≥ 20° de los tres, y de %d° a %d° está la zona reservada de acción crítica. Karen ocupa la L más alta de las personas en claro (la última que da 3:1) y Bodega Aurrera la L entre Mercado y Walmart; con esas L no llegan a ΔL 0.15 frente al texto de los roles (L %.3f en claro y %.3f en oscuro) en ningún tema, así que tienen que separarse por matiz de todos los roles, de acción crítica y de la marca. El script recorre los 360 matices enteros con la L y el croma finales de cada una y conserva los que cumplen todas las reglas de §7 (≥ 20° o ΔL ≥ 0.15 frente a roles, acción crítica y marca; contraste ≥ 3:1 sobre fondo y tarjeta; croma ≥ 0.10; fuera de la zona reservada; ΔL ≥ 0.10 en L OKLab y en gris frente a su grupo):\n" % (
        h_gasto, h_gasto, MATIZ_OPUESTO - ZONA_RESERVADA, MATIZ_OPUESTO + ZONA_RESERVADA, l_rol["claro"], l_rol["oscuro"]))
    ap(tabla(["Color", "L claro / oscuro · croma pedido", "Matices que cumplen todo", "Con Walmart en L %.3f (1.0)" % aw[3]], [
        ["Karen", "%.3f / %.3f · %.2f" % (ek["claro"].l, ek["oscuro"].l, ek["c"]), vtxt(wk), "igual (otro grupo)"],
        ["Bodega Aurrera", "%.3f / %.3f · %.2f" % (eb["claro"].l, eb["oscuro"].l, eb["c"]), vtxt(wb),
         "%s; con L %.3f y croma %.2f de la 1.0: %s" % (vtxt(wb0), ab[3], ab[1], vtxt(wb00))]]))
    dl_kb = abs(hex_a_oklab(ek["oscuro"].hex)[0] - hex_a_oklab(eb["oscuro"].hex)[0])
    if dl_kb >= SEPARACION_L:
        raise SystemExit("Cambios 1.1: Karen y Bodega Aurrera ya se separan por luminosidad en oscuro; revisar el texto")
    c0 = cands[0]
    ap("\nTodas las ventanas quedan entre deuda (%d°) y error (%d° = %d°), y en oscuro Karen y Bodega Aurrera tienen ΔL %.3f (< 0.15), así que también necesitan ≥ 20° de matiz entre sí. Se reparten los %d° que hay de deuda a error en tres huecos con el mayor hueco mínimo; a igualdad, Karen lo más lejos posible de Ximena (%d°, su mismo grupo) y después la mayor holgura frente a roles, acción crítica y marca. Resultado: **Karen %d° y Bodega Aurrera %d°** (huecos de %d°, %d° y %d°). Los cinco mejores repartos:\n" % (
        hd, he, he + 360, dl_kb, he + 360 - hd, M.cat["person"]["ximena"]["h"], c0["a"], c0["b"], c0["huecos"][0], c0["huecos"][1], c0["huecos"][2]))
    ap(tabla(["Karen", "Bodega Aurrera", "Huecos deuda → · → · → error", "Hueco mínimo", "Δ Karen–Ximena", "Holgura mín. frente a roles y marca"],
             [["%d°" % c["a"], "%d°" % c["b"], "%d° · %d° · %d°" % c["huecos"], "%d°" % c["min"], "%d°" % c["xim"], "%.2f" % c["hol"]] for c in cands[:5]]))
    ap("\nKaren queda a %d° de Ximena: dentro de las personas se distinguen por ΔL ≥ 0.10 en ambos temas (§7.3), por el orden fijo de la leyenda y por la etiqueta directa. Bodega Aurrera queda cerca de Ximena, pero son de grupos distintos y nunca comparten gráfica ni leyenda (§7.7).\n" % c0["xim"])
    ap(md_rosas_ximena(M))
    so = [M.cat["store"][c]["oscuro"].hex for c in ORDEN_TIENDAS]
    sup_o = (M.sup("oscuro", "fondo"), M.sup("oscuro", "tarjeta"))
    w0, mk_, sm_ = Color(aw[3], aw[1], aw[0]), M.cat["store"]["market"]["oscuro"], M.cat["store"]["sams-club"]["oscuro"]
    ap("### 3. Walmart en oscuro (PROVISIONAL)\n")
    ap("En oscuro las cuatro tiendas se escalonan con ΔL ≥ 0.10 en L OKLab y en gris entre Mercado (L %.3f, contraste %.2f, el mínimo del grupo) y Sam's Club (L %.3f, croma %.3f, cerca del piso de 0.10). El rosa pierde luminancia en gris (en Bodega Aurrera, gris = L − %.3f), y con Walmart en L %.3f ningún matiz entero cumple para Bodega Aurrera (última columna de la tabla del punto 2). Walmart baja a L %.3f (`%s` → `%s`, ΔE %.1f, imperceptible) y Bodega Aurrera usa croma %.2f y L %.3f en oscuro; el escalonamiento de las tiendas en oscuro queda en ΔL mínimo %.3f (OKLab) y %.3f (gris). Como es un token que K26 no nombra, va **PROVISIONAL** con el bloqueo `docs/estado/bloqueos/U2a-walmart-oscuro.md` (opciones para Axel) y la nota en la descripción del token en `C2-colores.json`.\n" % (
        mk_.l, min(contraste(mk_.hex, x) for x in sup_o), sm_.l, hex_a_oklch(sm_.hex)[1], eb["oscuro"].l - l_gris(eb["oscuro"].hex), aw[3],
        ew["oscuro"].l, w0.hex, ew["oscuro"].hex, delta_e(w0.hex, ew["oscuro"].hex), eb["c"], eb["oscuro"].l,
        min(abs(hex_a_oklab(a)[0] - hex_a_oklab(b)[0]) for a, b in combinations(so, 2)), min(abs(l_gris(a) - l_gris(b)) for a, b in combinations(so, 2))))
    mp = [f for f in pares_cercanos(M) if "(marca)" in f["a"] + f["b"]]
    ap("### 4. La marca en el barrido de pares cercanos (R3-04)\n")
    ap("§7.7 incluye ahora la marca en su paso de texto (%d en claro, %d en oscuro). Ningún rol ni acción crítica queda cerca de la marca (el script lo verifica); quedan cerca estas categóricas, que se refuerzan con forma y etiqueta: la marca solo vive en controles (botón primario, enlace subrayado, anillo de foco, elemento activo) y la categórica siempre va como punto, anillo de avatar o marca de serie junto a su nombre en texto neutro.\n" % (PASO_MARCA["claro"], PASO_MARCA["oscuro"]))
    ap(tabla(["Tema", "Color A", "Color B", "Δ matiz", "ΔL"], [[f["tema"], f["a"], f["b"], "%.1f°" % f["dm"], "%.3f" % f["dl"]] for f in mp]))
    ap("\n### 5. Paleta de gráficas (K27)\n")
    ap("No cambia ningún valor de las 8 series. K27 fija la lectura del criterio: el umbral duro es ΔL ≥ 0.10 entre series consecutivas en el orden de uso (mínimo en gris y en las tres simulaciones: %.3f en claro y %.3f en oscuro, §8.3) y los 28 pares se reportan como información (ΔL mínimo en gris %.3f en claro y %.3f en oscuro, §8.5). El color nunca va solo: cada serie lleva forma de marcador, patrón de trazo o relleno y etiqueta directa. Las tablas de simulación de daltonismo de las series (§8.4) no cambian porque las series no cambian; los ΔE con daltonismo de las categóricas (§7.3) se recalculan con Karen y Bodega Aurrera nuevas.\n" % (
        G["claro"]["min_dl"], G["oscuro"]["min_dl"], G["claro"]["dl_todos_gris"], G["oscuro"]["dl_todos_gris"]))
    return o


def generar_md(M, pares, G, Gg, C, V, I):
    o = []
    ap = o.append
    n50, n100, n900, n950 = M.neutro[50].hex, M.neutro[100].hex, M.neutro[900].hex, M.neutro[950].hex
    total = len(pares)
    n_texto = sum(1 for p in pares if p["umbral"] == UMBRAL_TEXTO)
    n_ui = total - n_texto
    peor_t = min((p for p in pares if p["umbral"] == UMBRAL_TEXTO), key=lambda p: p["ratio"])
    peor_u = min((p for p in pares if p["umbral"] == UMBRAL_UI), key=lambda p: p["ratio"])
    plan_ok = sum(M.marca[p].hex[1:] == PLAN_C["marca"][i] for i, p in enumerate(PASOS)) + \
        sum(M.neutro[p].hex[1:] == PLAN_C["neutro"][i] for i, p in enumerate(PASOS))
    plan_rol_cr = [min(contraste(v[0], n50), contraste(v[0], n100)) for v in PLAN_ROLES.values() if v[0]]
    plan_rol_bajo = sum(1 for x in plan_rol_cr if x < UMBRAL_TEXTO)
    p600 = [M.paso_rol[r["codigo"]]["historial"]["claro"][0][1] for r in ROLES + [CRITICO]]
    cat_cl = sum(1 for v in PLAN_CATEGORICAS.values() if min(contraste(v[0], n50), contraste(v[0], n100)) < UMBRAL_UI)
    cat_os = sum(1 for v in PLAN_CATEGORICAS.values() if min(contraste(v[0], n950), contraste(v[0], n900)) < UMBRAL_UI)

    ap("# Colorimetría · candidato C «Índigo 280°»\n")
    ap("> Documento generado por `%s` (Python 3, solo biblioteca estándar). No se edita a mano: se cambia la regla en el script y se regenera. Fase U2a del flujo WF-0, %s. Es la base de color de `DESIGN.md` (U2, WF-1) y la parte de color del contrato C2 (`docs/contratos/C2-colores.json`); el resto de los tokens de C2 los hace U2b en WF-1. Queda sujeto a la puerta 2 de aprobación. Ronda 1 de revisión (2026-10-11): corrige el hallazgo H-03 (la serie 6 de gráficas usaba el matiz reservado del eje opuesto; §3 y §8.8). **Versión %s (puerta 1, %s):** aplica las respuestas K26 y K27 de Axel y el hallazgo R3-04; lo que cambia, con valores anteriores y nuevos, está en «Cambios 1.1».\n" % (SCRIPT_REL, FECHA, VERSION, FECHA_11))
    ap("## Índice\n")
    ap("0. Cambios 1.1 (puerta 1): gasto, Karen, Bodega Aurrera, marca en pares cercanos y paleta de gráficas\n1. Qué se genera y cómo regenerarlo\n2. Método\n3. Matiz base, eje opuesto y armonía complementaria dividida\n4. Regla 60-30-10\n5. Escalas de 11 pasos\n6. Roles semánticos\n7. Paleta categórica: personas y tiendas\n8. Paleta de gráficas (8 series)\n9. Reglas de uso por componente\n10. Contraste (resumen de `contraste.csv`)\n11. Tokens (`C2-colores.json`)\n12. Hallazgos, decisiones derivadas del cálculo y bloqueos\n\nAnexo A. Tokens: ruta → hex (lista que verifica el script)\n")

    o.extend(md_cambios_11(M, G))

    # 1 ------------------------------------------------------------------------------
    ap("## 1. Qué se genera y cómo regenerarlo\n")
    ap(tabla(["Archivo", "Contenido"], [
        ["`docs/design/colorimetria.md`", "este documento"],
        ["`docs/design/contraste.csv`", "%d pares de contraste (%d de texto, %d de interfaz y gráficas), claro y oscuro" % (total, n_texto, n_ui)],
        ["`docs/contratos/C2-colores.json`", "%d tokens DTCG de color" % len(M.tokens)],
        ["`docs/design/paleta.html`", "muestrario estático (sin CDN, sin fuentes externas, sin JavaScript)"],
        ["`docs/design/paleta.png`", "captura 1400×3000 con `chromium --headless=new`"],
    ]))
    ap("\nComando: `python3 %s` (agrega `--sin-png` para omitir la captura y `--solo-verificar` para revisar los archivos existentes sin regenerarlos). La salida es determinista: dos ejecuciones dan archivos idénticos. Al final el script verifica, leyendo los archivos ya escritos, que ningún par del CSV quede bajo su umbral, que los hex del JSON coincidan con el Anexo A de este documento, que la paleta de gráficas mantenga ΔL ≥ %.2f entre series consecutivas en gris y en las tres simulaciones (con las dos variantes de §8.6) y que sus tres primeras series pasen la prueba de todos los pares (§8.5), que las categóricas cumplan su escalonamiento y su separación de los roles, que ningún token fuera de `color.critical` use la zona reservada del eje opuesto (§3), que marca y neutro reproduzcan los valores de partida del plan y, desde la versión 1.1, que gasto esté en 40–45° a ≥ 20° de error, que la marca no quede cerca de ningún rol y que Karen y Bodega Aurrera cumplan la regla de separación entre sí; si algo falla, termina con código 1.\n" % (SCRIPT_REL, UMBRAL_DL))

    # 2 ------------------------------------------------------------------------------
    ap("## 2. Método\n")
    ap("Resume y concreta el plan §4.2. Cada color se define en OKLCH (luminosidad L, croma C, matiz h en grados) y se convierte a sRGB de 8 bits.\n")
    ap("### 2.1 Conversión\n")
    ap("OKLCH → OKLab (a = C·cos h, b = C·sin h) → LMS no lineal (inversa de M2) → LMS (cubo) → sRGB lineal (inversa de M1) → sRGB codificado (función de transferencia de IEC 61966-2-1: 12.92·c si c ≤ 0.0031308; 1.055·c^(1/2.4) − 0.055 en otro caso) → hex con redondeo ROUND_HALF_UP por canal. Las matrices M1 y M2 son las publicadas por Björn Ottosson en «A perceptual color space for image processing» (bottosson.github.io/posts/oklab). La conversión inversa (hex → OKLab) se usa para medir.\n")
    ap("### 2.2 Luminosidad fija por paso\n")
    ap(tabla(["Paso"] + [str(p) for p in PASOS], [["L"] + ["%.3f" % L_PASO[p] for p in PASOS]]))
    ap("\nTodas las escalas (marca, neutro, roles y acción crítica) usan la misma L en cada paso, así que dos colores del mismo paso tienen la misma luminosidad percibida aunque cambie el matiz.\n")
    ap("### 2.3 Perfil de croma por paso\n")
    ap(tabla(["Paso"] + [str(p) for p in PASOS], [
        ["Fracción del croma máximo"] + ["%.2f" % PERFIL_CROMA[p] for p in PASOS],
        ["Intervalo que reproduce el plan"] + ["%.3f–%.3f" % I[p] for p in PASOS],
    ]))
    ap("\nCroma máximo por familia: marca %.3f, neutro tintado %.3f, roles y acción crítica %.2f (plan §4.2: «saturación apagada»). El plan define el perfil solo por sus extremos y su pico (0.10 → 1.0 → 0.45); la tabla explícita de arriba reproduce los **%d de 22** valores de partida del plan §4.3 (marca y neutro del candidato C). La segunda fila muestra, por paso, el intervalo de fracciones que también los reproduce: en los pasos 50–300 la marca queda recortada por el gamut y el neutro se pierde en la cuantización, así que esos pasos solo quedan acotados por intervalo y se eligieron 0.10, 0.25, 0.40 y 0.65 dentro de él. Esto importa en las escalas de rol, cuyos pasos claros no siempre se recortan.\n" % (CROMA_MARCA, CROMA_NEUTRO, CROMA_ROL, plan_ok))
    ap("### 2.4 Recorte al gamut sRGB\n")
    ap("Si el color pedido cae fuera de sRGB, se busca por bisección (60 iteraciones) el mayor croma que cabe, manteniendo L y h; un color está dentro si sus tres canales lineales quedan en [0, 1] con tolerancia 1e-9. Por eso las tablas muestran «croma pedido» y «croma final». El matiz y la luminosidad nunca se alteran para encajar.\n")
    ap("### 2.5 Contraste (WCAG 2.x)\n")
    ap("Luminancia relativa Y = 0.2126·R + 0.7152·G + 0.0722·B con R, G, B lineales (umbral 0.04045 de IEC 61966-2-1, el que cita la nota de WCAG 2.2; con valores de 8 bits da lo mismo que 0.03928). Razón = (Y₁ + 0.05) / (Y₂ + 0.05) con Y₁ ≥ Y₂. Umbrales: texto ≥ 4.5:1 (criterio 1.4.3, incluye texto deshabilitado por regla del proyecto) e interfaz, bordes de campo, foco, series y marcadores categóricos ≥ 3:1 (criterio 1.4.11). Las razones se truncan a dos decimales (nunca se redondean hacia arriba) y el cumplimiento se evalúa con el valor sin truncar.\n")
    ap("### 2.6 Simulación de daltonismo\n")
    ap("Modelo de Machado, Oliveira y Fernandes (2009), «A Physiologically-based Model for Simulation of Color Vision Deficiency», IEEE TVCG 15(6):1291-1298, doi:10.1109/TVCG.2009.113. Se usan las matrices de severidad 1.0 (dicromacia) de la Tabla 1 publicada por los autores en inf.ufrgs.br/~oliveira/pubs_files/CVD_Simulation/CVD_Simulation.html (consultada el %s; coinciden con las que trae culori 4.0.2):\n" % FECHA)
    filas_m = []
    for tcv in TIPOS_CVD:
        m = MACHADO_2009[tcv]
        filas_m.append([NOMBRE_VISTA[tcv]] + ["%.6f %.6f %.6f" % fila for fila in m])
    ap(tabla(["Tipo", "Fila R", "Fila G", "Fila B"], filas_m))
    ap("\nSe aplican sobre sRGB **lineal** y el resultado se recorta a [0, 1]. El artículo asume RGB lineal (p. 1294) y así lo hace el paquete R colorspace desde la versión 2.1-0 (`simulate_cvd(linear = TRUE)`). Ojo: `filterDeficiencyProt/Deuter/Trit` de culori 4.0.2 aplican la matriz sobre sRGB codificado (`src/deficiency.js`); por eso el script también mide la paleta de gráficas con esa variante, como prueba de robustez (§8.6).\n")
    ap("### 2.7 Escala de grises y distancias\n")
    ap("El gris equivalente de un color es el gris de igual luminancia relativa Y; su luminosidad OKLab es exactamente L_gris = Y^(1/3) (para un gris, l = m = s = Y). La distancia entre colores es ΔE = 100·‖ΔOKLab‖ (euclidiana). Los umbrales de distinción entre series vecinas son los de la skill dataviz que el plan §4.1 adopta: ΔE ≥ %.0f bajo protanopia y deuteranopia (piso %.0f solo con codificación secundaria) y ΔE ≥ %.0f en visión normal; además ΔL ≥ %.2f entre series consecutivas en gris y en cada simulación (criterio U2 del plan §3.5; el encargo pide ≥ %.2f).\n" % (CVD_OBJETIVO, CVD_PISO, NORMAL_PISO, UMBRAL_DL, UMBRAL_DL_ENCARGO))

    # 3 ------------------------------------------------------------------------------
    ap("## 3. Matiz base, eje opuesto y armonía complementaria dividida\n")
    adv = next(r for r in ROLES if r["codigo"] == "warning")
    fmt_cerca = lambda it: "%s (%s) a %.1f°" % (it[0][1].lower() if it[0][0] == "rol" else it[0][1], it[0][0], it[1])
    zona_txt = ("- **Zona reservada %d°–%d°.** Ningún otro token con croma visible (≥ %.2f; el neutro llega a %.3f) usa un matiz a menos de %d° del opuesto: ni roles, ni personas y tiendas, ni series de gráficas. El ancho es la separación mínima entre matices distintos que decisiones §I fija para roles y categóricas, y no puede ser mayor: advertencia (%d°, también decidida) queda a %d° del opuesto (%.1f° medidos en el hex). El script lo verifica en cada token, con el matiz de diseño y con el medido en el hex. Colores de identidad más cercanos al opuesto: %s. Hasta la ronda 1 de revisión la serie %d de gráficas usaba 100° (hallazgo H-03); ahora usa %d° (§8.8).\n"
                % (MATIZ_OPUESTO - ZONA_RESERVADA, MATIZ_OPUESTO + ZONA_RESERVADA, CROMA_VISIBLE, CROMA_NEUTRO, ZONA_RESERVADA,
                   adv["matiz"], dif_matiz(adv["matiz"], MATIZ_OPUESTO), dict(cercanos_reserva(M, 99)).get(("rol", adv["nombre"]), 0.0), "; ".join(fmt_cerca(x) for x in cercanos_reserva(M)),
                   SERIE_CORREGIDA, M.graf[SERIE_CORREGIDA]["h"]))
    ap("- **Matiz base 280° (índigo)**, aprobado por Axel (decisiones §I). Lleva la marca, las superficies tintadas (neutro con croma 0.012, apenas perceptible) y el acento principal. Lectura sobria y moderna, tipo Notion/Linear.\n- **Eje opuesto 100° (oliva amarillo)** = 280° − 180°. Es el matiz de mayor contraste frente a la marca, por eso se **reserva para «acción crítica» de confirmación** (decisiones §I): un diálogo que pide confirmar algo irreversible es lo único que debe interrumpir. No se usa para error, que conserva la convención del rojo (20°).\n" + zona_txt + "- **Complementario dividido:** los dos vecinos del opuesto a ±30°, 70° y 130°, son los acentos que acompañan a la marca sin chocar con ella. Se asignan a los dos estados que más se ven en una app de finanzas además del gasto: **advertencia 75°** (decidida; 5° por encima del vecino exacto 70°) e **ingreso 135°** (la opción que decisiones §I prevé para U2; 5° por encima del vecino exacto 130°, el mismo desplazamiento que advertencia, §6.4). Así la armonía queda base 280° + acentos 75°/135° + opuesto 100° en reserva, y se evita el choque directo rojo/verde porque gasto y error (rojos) y éxito (verde) no forman parte de la armonía: son convenciones semánticas, con su propio matiz y siempre con icono y texto.\n")
    filas_h = []
    for r in ROLES + [CRITICO]:
        filas_h.append([r["nombre"], "%d°" % r["matiz"], "%d°" % dif_matiz(r["matiz"], MATIZ_BASE), "%d°" % dif_matiz(r["matiz"], MATIZ_OPUESTO), r["origen"]])
    filas_h.insert(0, ["Marca", "280°", "0°", "180°", "decisiones §I"])
    ap(tabla(["Familia", "Matiz", "Distancia a la base", "Distancia al opuesto", "Origen"], filas_h))
    ap("")

    # 4 ------------------------------------------------------------------------------
    ap("## 4. Regla 60-30-10\n")
    ap("Proporción orientativa del área visible de cada pantalla (no se mide en píxeles; guía la composición):\n")
    ap(tabla(["Proporción", "Papel", "Qué cubre", "Claro", "Oscuro"], [
        ["60 %", "Superficies", "fondo de página, tarjetas, divisores, texto neutro", "neutral 50 (fondo), 100 (tarjeta), 200 (divisor), 700/900 (texto)", "neutral 950 (fondo), 900 (tarjeta), 800 (divisor), 300/100 (texto)"],
        ["30 %", "Acento", "marca en pasos de superficie: barra lateral, elemento activo, fila seleccionada, encabezados de vista, chips de filtro", "brand 50–200", "brand 900–950"],
        ["10 %", "Énfasis", "una acción primaria por pantalla, enlaces, foco, roles semánticos, categóricas, series de gráficas, acción crítica", "brand 700/800, roles paso 700, categóricas y series claras", "brand 300/200, roles paso 400, categóricas y series oscuras"],
    ]))
    ap("\nEl neutro está tintado con el matiz base (croma ≤ 0.012), así que el 90 % de la pantalla ya pertenece a la familia índigo y el 10 % de énfasis es lo único con saturación visible.\n")

    # 5 ------------------------------------------------------------------------------
    ap("## 5. Escalas de 11 pasos\n")
    for titulo, esc, plan, cmax in (("5.1 Marca (280°, croma máx. 0.11)", M.marca, PLAN_C["marca"], CROMA_MARCA),
                                    ("5.2 Neutro tintado (280°, croma máx. 0.012)", M.neutro, PLAN_C["neutro"], CROMA_NEUTRO)):
        ap("### %s\n" % titulo)
        filas = []
        for i, p in enumerate(PASOS):
            c = esc[p]
            filas.append([p, "%.3f" % c.l, "%.4f" % c.c_pedido, "%.4f" % c.c + (" (recortado)" if c.recortado else ""), "%.0f°" % c.h, "`%s`" % c.hex, "`#%s`" % plan[i], "sí" if c.hex[1:] == plan[i] else "**no**"])
        ap(tabla(["Paso", "L", "Croma pedido", "Croma final", "Matiz", "Hex", "Plan §4.3", "Coincide"], filas))
        ap("")
    ap("### 5.3 Escalas de rol y de acción crítica (croma máx. 0.14)\n")
    filas = []
    for r in ROLES + [CRITICO]:
        esc = M.esc_rol[r["codigo"]]
        filas.append(["%s (%d°)" % (r["nombre"], r["matiz"])] + ["`%s`" % esc[p].hex for p in PASOS])
    ap(tabla(["Familia"] + [str(p) for p in PASOS], filas))
    ap("\nCada rol tiene su escala completa en el JSON (`color.role.<rol>.50` … `.950`, y `color.critical.*`), porque los estados hover usan pasos vecinos (800 en claro, 200 en oscuro).\n")

    # 6 ------------------------------------------------------------------------------
    ap("## 6. Roles semánticos\n")
    ap("### 6.1 Valores finales\n")
    ap("Texto e icono del rol: paso 700 en claro y paso 400 en oscuro. Fondo de insignia: paso 100 en claro y paso 900 en oscuro. Contraste mínimo exigido: 4.5:1, porque los roles se usan como texto (cifras con signo, mensajes, insignias).\n")
    filas = []
    for r in ROLES + [CRITICO]:
        cod = r["codigo"]
        b = base_rol(cod)
        lc, lo = M.hex(b + ".light"), M.hex(b + ".dark")
        sl, so = M.hex(b + ".surface-light"), M.hex(b + ".surface-dark")
        filas.append([r["nombre"], "`%s`" % b, "%d°" % r["matiz"],
                      "%d · `%s`" % (M.paso_rol[cod]["claro"], lc), "%s / %s / %s" % (r2(contraste(lc, n50)), r2(contraste(lc, n100)), r2(contraste(lc, sl))),
                      "%d · `%s`" % (M.paso_rol[cod]["oscuro"], lo), "%s / %s / %s" % (r2(contraste(lo, n950)), r2(contraste(lo, n900)), r2(contraste(lo, so)))])
    ap(tabla(["Rol", "Token", "Matiz", "Claro: paso · hex", "Contraste fondo 50 / tarjeta 100 / insignia", "Oscuro: paso · hex", "Contraste fondo 950 / tarjeta 900 / insignia"], filas))
    ap("\n«Acción crítica» no es uno de los 8 roles: es el énfasis del eje opuesto (§3) y vive en `color.critical`.\n")
    ap("### 6.2 Ajuste de paso: 600 → 700 en claro\n")
    ap("El encargo pide el paso 600 sobre fondo 50 en claro. Con la L fija del paso 600 (0.605) ningún matiz llega a 4.5:1, así que se baja un paso (700, L 0.515) y se documenta aquí. En oscuro el paso 400 cumple sin cambios.\n")
    filas = []
    for r in ROLES + [CRITICO]:
        cod = r["codigo"]
        h = M.paso_rol[cod]["historial"]["claro"]
        esc = M.esc_rol[cod]
        filas.append([r["nombre"], "`%s`" % esc[600].hex, r2(h[0][1]), "`%s`" % esc[M.paso_rol[cod]["claro"]].hex, r2(h[-1][1]), "%d → %d" % (600, M.paso_rol[cod]["claro"])])
    ap(tabla(["Rol", "Paso 600", "Contraste mínimo (fondo, tarjeta, insignia)", "Paso elegido", "Contraste mínimo", "Cambio"], filas))
    ap("")
    ap("### 6.3 Comparación con los valores del plan §4.3\n")
    ap("El plan calculó los roles a L 0.56 en claro y L 0.74 en oscuro (sus contrastes son contra marca 50 `#F9FAFF`). Recalculados contra las superficies definitivas, varios quedan bajo 4.5:1 para texto, lo que confirma el ajuste anterior (%d de %d con valor quedan bajo 4.5:1 en al menos una superficie; rango %s–%s:1):\n" % (plan_rol_bajo, len(plan_rol_cr), r2(min(plan_rol_cr)), r2(max(plan_rol_cr))))
    filas = []
    for r in ROLES:
        pc, po, nota = PLAN_ROLES[r["codigo"]]
        if not pc:
            filas.append([r["nombre"], "—", "—", "—", "—", nota])
            continue
        filas.append([r["nombre"], "`%s`" % pc, "%s / %s" % (r2(contraste(pc, n50)), r2(contraste(pc, n100))), "`%s`" % po, "%s / %s" % (r2(contraste(po, n950)), r2(contraste(po, n900))), nota or "—"])
    ap(tabla(["Rol", "Plan claro", "Contraste fondo 50 / tarjeta 100", "Plan oscuro", "Contraste fondo 950 / tarjeta 900", "Nota"], filas))
    ap("")
    dist_c = {(f["a"], f["b"]): f for f in distancias_roles(M, "claro")}
    dist_o = {(f["a"], f["b"]): f for f in distancias_roles(M, "oscuro")}
    ie_c, ie_o = dist_c[("Ingreso", "Éxito")], dist_o[("Ingreso", "Éxito")]
    ge_c, ge_o = dist_c[("Gasto", "Error")], dist_o[("Gasto", "Error")]
    h_gasto = next(r["matiz"] for r in ROLES if r["codigo"] == "expense")
    antes11 = tokens_antes_11(M)
    ge0_c = delta_e(antes11["color.role.expense.light"].hex, M.hex("color.role.error.light"))
    ge0_o = delta_e(antes11["color.role.expense.dark"].hex, M.hex("color.role.error.dark"))
    gm_c, gm_o = (dif_matiz(hex_a_oklch(M.hex("color.role.expense.%s" % x))[2], hex_a_oklch(M.hex("color.role.error.%s" % x))[2]) for x in ("light", "dark"))
    ap("### 6.4 Ingreso frente a éxito (propuesta)\n")
    ap("Decisiones §I deja ingreso en 150° (misma familia que éxito, distinguido por icono y contexto) «o U2 propone 135°». **Propuesta U2a: ingreso 135°**, la opción ya prevista, así que no se reabre nada. Razones: (1) separa ingreso de éxito 15° en vez de 0°; (2) queda a 5° del acento dividido 130° (100° + 30°), el mismo desplazamiento que advertencia tiene respecto a 70°, de modo que ingreso y advertencia forman el par de acentos de la armonía (§3); (3) queda a 35° del opuesto 100° (acción crítica) y a 60° de advertencia. Distancias calculadas entre ingreso y éxito (texto de rol): claro ΔE normal %.1f, protanopia %.1f, deuteranopia %.1f, tritanopia %.1f; oscuro %.1f, %.1f, %.1f, %.1f. Lado a lado apenas se distinguen en visión normal (ΔE ≈ 3, cerca del umbral de percepción) y no se distinguen con daltonismo, por eso el color nunca va solo: **ingreso** siempre lleva signo «+» y flecha de entrada en cifras y listas; **éxito** siempre lleva palomita y texto («Guardado», «Pagado»), y vive en avisos y estados, no en cifras.\n" % (ie_c["normal"], ie_c["protanopia"], ie_c["deuteranopia"], ie_c["tritanopia"], ie_o["normal"], ie_o["protanopia"], ie_o["deuteranopia"], ie_o["tritanopia"]))
    ap("### 6.5 Gasto frente a error\n")
    ap("Versión 1.1 (K26; bloqueo `U2-gasto-error` resuelto con la opción 2): gasto pasa de %d° a %d° (terracota) y error se queda en 20°. Medidos en el hex quedan a %.1f° de matiz en claro y %.1f° en oscuro, con ΔE normal %.1f y %.1f (en la 1.0, con gasto en %d°, eran %.1f y %.1f); el grado sale del barrido de «Cambios 1.1». Con daltonismo rojo-verde siguen cerca (ΔE protanopia %.1f / %.1f, deuteranopia %.1f / %.1f), así que se conserva la regla de uso: el **gasto** aparece como cifra con signo «−» en listas y tablas; el **error** solo como mensaje de validación o aviso con icono de alerta y texto. Nunca comparten componente.\n" % (
        ANTES_11["expense"], h_gasto, gm_c, gm_o, ge_c["normal"], ge_o["normal"], ANTES_11["expense"], ge0_c, ge0_o, ge_c["protanopia"], ge_o["protanopia"], ge_c["deuteranopia"], ge_o["deuteranopia"]))
    ap("### 6.6 Pares de roles cercanos (ΔE normal < 10)\n")
    ap("Todos los roles comparten L y croma, así que la distancia depende solo del ángulo de matiz. Se listan los pares con ΔE normal < 10 (los que pueden confundirse a simple vista):\n")
    filas = []
    for tema, dist in (("claro", dist_c), ("oscuro", dist_o)):
        for f in sorted(dist.values(), key=lambda x: x["normal"]):
            if f["normal"] < 10.0:
                filas.append([tema, "%s – %s" % (f["a"], f["b"]), "%.0f°" % f["dm"], "%.1f" % f["normal"], "%.1f" % f["protanopia"], "%.1f" % f["deuteranopia"], "%.1f" % f["tritanopia"]])
    ap(tabla(["Tema", "Par", "Δ matiz", "ΔE normal", "ΔE protanopia", "ΔE deuteranopia", "ΔE tritanopia"], filas))
    ap("\nTodos estos pares viven en contextos distintos (cifra, aviso, diálogo de confirmación) y todos llevan icono y texto; ningún significado depende solo del color (criterio 1.4.1).\n")

    # 7 ------------------------------------------------------------------------------
    ap("## 7. Paleta categórica: personas y tiendas\n")
    ap("Reglas (plan §4.2–4.3 y encargo): luminosidad escalonada con ΔL ≥ %.2f entre todos los pares de cada grupo (medida en L OKLab y en L del gris equivalente), separación ≥ %.0f° de matiz **o** ΔL ≥ %.2f respecto a cada rol del mismo tema (también se exige frente a acción crítica), y contraste ≥ 3:1 sobre fondo y tarjeta. Un solo valor por color no puede cumplir 3:1 en los dos temas y a la vez escalonarse: la ventana común es L_gris ∈ [%.3f, %.3f] (ancho %.3f < %.2f). Por eso cada color tiene variante clara y oscura con el **mismo matiz** (la identidad se conserva) y distinta L.\n" % (UMBRAL_DL, SEPARACION_MATIZ, SEPARACION_L, V["l_min"], V["l_max"], V["l_max"] - V["l_min"], UMBRAL_DL))
    for grupo, orden, titulo in (("person", ORDEN_PERSONAS, "7.1 Personas"), ("store", ORDEN_TIENDAS, "7.2 Tiendas")):
        ap("### %s\n" % titulo)
        filas = []
        for cod in orden:
            d = M.cat[grupo][cod]
            cc, co = d["claro"], d["oscuro"]
            filas.append([d["nombre"], "`color.%s.%s`" % (grupo, cod), "%d°" % d["h"], "%.3f" % d["c"],
                          "`%s` · L %.3f · gris %.3f · %s" % (cc.hex, cc.l, l_gris(cc.hex), r2(min(contraste(cc.hex, n50), contraste(cc.hex, n100)))),
                          "`%s` · L %.3f · gris %.3f · %s" % (co.hex, co.l, l_gris(co.hex), r2(min(contraste(co.hex, n950), contraste(co.hex, n900)))),
                          d["motivo"]])
        ap(tabla(["Nombre", "Token", "Matiz", "Croma pedido", "Claro: hex · L · L gris · contraste mín.", "Oscuro: hex · L · L gris · contraste mín.", "Motivo"], filas))
        ap("")
    ap("### 7.3 Escalonamiento y distinción dentro de cada grupo\n")
    filas = []
    for grupo, titulo in (("person", "Personas"), ("store", "Tiendas")):
        for tema in TEMAS:
            d = C[grupo][tema]
            filas.append([titulo, tema, "%.3f" % d["dl_oklab"], "%.3f" % d["dl_gris"], "%.1f" % d["adj_normal"], "%.1f" % d["adj_cvd"], "%.1f" % d["adj_tritan"], r2(d["contraste_min"]), "%.3f" % d["croma_min"]])
    ap(tabla(["Grupo", "Tema", "ΔL mín. (OKLab, todos los pares)", "ΔL mín. (gris, todos los pares)", "ΔE vecinos normal", "ΔE vecinos protan/deutan", "ΔE vecinos tritan", "Contraste mín.", "Croma mín. (sin Conjunto)"], filas))
    ap("\nOrden fijo de leyenda y apilado: personas Axel, Ximena, Karen, Conjunto; tiendas Walmart, Bodega Aurrera, Sam's Club, Mercado (los «vecinos» de la tabla son los contiguos en ese orden). Con cuatro niveles separados ≥ 0.10 no cabe la banda de L de la skill dataviz (§8.7): el escalonamiento ocupa casi todo el rango que permite el contraste, así que algunos colores quedan muy oscuros en claro (L 0.30–0.31) o muy claros en oscuro (L 0.93–0.94). El piso de croma 0.10 se cumple en todos salvo Conjunto, que es bajo a propósito.\n")
    ap("### 7.4 Separación respecto a los roles\n")
    ap("Para cada color, el rol más exigente (menor holgura) y la regla con la que cumple:\n")
    filas = []
    for grupo, orden in (("person", ORDEN_PERSONAS), ("store", ORDEN_TIENDAS)):
        for cod in orden:
            for tema in TEMAS:
                sep = C[grupo][tema]["sep"][cod]
                peor = min(sep, key=lambda f: f["holgura"])
                regla = "Δ matiz ≥ 20°" if peor["dm"] >= SEPARACION_MATIZ else "ΔL ≥ 0.15"
                filas.append([M.cat[grupo][cod]["nombre"], tema, peor["rol"], "%.1f°" % peor["dm"], "%.3f" % peor["dl"], regla, "sí" if all(f["ok"] for f in sep) else "**no**"])
    ap(tabla(["Color", "Tema", "Rol más cercano", "Δ matiz", "ΔL", "Cumple por", "Cumple con todos"], filas))
    ap("")
    ap("### 7.5 De la propuesta del plan a los valores finales\n")
    plan_cat = PLAN_CATEGORICAS
    filas = []
    for grupo, orden in (("person", ORDEN_PERSONAS), ("store", ORDEN_TIENDAS)):
        for cod in orden:
            hx, h, L = plan_cat[cod]
            d = M.cat[grupo][cod]
            filas.append([d["nombre"], "`%s` %d° L %.2f" % (hx, h, L), "%s / %s" % (r2(min(contraste(hx, n50), contraste(hx, n100))), r2(min(contraste(hx, n950), contraste(hx, n900)))),
                          "%d° · `%s` / `%s`" % (d["h"], d["claro"].hex, d["oscuro"].hex)])
    ap(tabla(["Nombre", "Plan (valor único)", "Contraste del plan claro / oscuro", "Final: matiz · claro / oscuro"], filas))
    ap("\nColisiones resueltas: las dos que lista decisiones §I (Karen 80° y Bodega 60° frente a advertencia 75°) y tres más que aparecen al aplicar la regla a todos los roles: Axel 250° frente a información 240°, Mercado 30° frente a gasto 25° y error 20°, y Walmart 220° justo en el límite de ahorro 200° e información 240°. Además, %d valores del plan no llegaban a 3:1 sobre las superficies claras y %d sobre las oscuras. En la versión 1.1 (K26) gasto pasa a terracota, Karen y Bodega Aurrera salen de 50° y Walmart ajusta su L en oscuro (PROVISIONAL); frente a Ximena, Karen y Bodega Aurrera quedan PROVISIONALES (bloqueo `U2a-rosas-ximena-bodega`). Detalle en «Cambios 1.1».\n" % (cat_cl, cat_os))
    ap("### 7.6 Reglas de uso de las categóricas\n")
    ap("- Nunca llevan texto encima: el nombre va en texto neutro junto a un punto, anillo de avatar o marca de serie del color.\n- La misma persona o tienda conserva su matiz en ambos temas; solo cambia la luminosidad.\n- En gráficas por persona o por tienda se usan estos colores (el color sigue a la entidad) y no la paleta de series; con leyenda siempre y etiqueta directa cuando hay ≤ 4 elementos.\n- Mercados adicionales que Axel dé de alta usan el color de Mercado; si una gráfica necesita distinguir varios mercados, se usan las series de gráficas en su orden.\n")

    ap("### 7.7 Pares cercanos en todo el sistema y su refuerzo\n")
    cercanos = pares_cercanos(M)
    if any(f["refuerzo"] == "SIN RESOLVER" for f in cercanos):
        raise SystemExit("§7.7: hay pares sin refuerzo (categórica o marca frente a rol)")
    ap("Revisión de todos los pares entre los 8 roles, acción crítica, las 4 personas, las 4 tiendas y la marca (paso %d en claro y %d en oscuro; hallazgo R3-04, versión 1.1), en cada tema, con matiz y L medidos en el hex final. Se listan los que quedan con Δ matiz < 20° **y** ΔL < 0.15, con el refuerzo que evita confundirlos (ninguno es categórica ni marca frente a rol):\n" % (PASO_MARCA["claro"], PASO_MARCA["oscuro"]))
    filas = [[f["tema"], f["a"], f["b"], "%.1f°" % f["dm"], "%.3f" % f["dl"], f["refuerzo"]] for f in cercanos]
    ap(tabla(["Tema", "Color A", "Color B", "Δ matiz", "ΔL", "Refuerzo"], filas))
    ap("")

    # 8 ------------------------------------------------------------------------------
    ap("## 8. Paleta de gráficas (8 series)\n")
    ap("### 8.1 Orden de uso\n")
    filas = []
    for n in range(1, 9):
        g = M.graf[n]
        cc, co = g["claro"], g["oscuro"]
        filas.append([n, g["nombre"], "%d°" % g["h"], g["clase"],
                      "`%s` · L %.3f · C %.3f · %s" % (cc.hex, cc.l, cc.c, r2(min(contraste(cc.hex, n50), contraste(cc.hex, n100)))),
                      "`%s` · L %.3f · C %.3f · %s" % (co.hex, co.l, co.c, r2(min(contraste(co.hex, n950), contraste(co.hex, n900))))])
    ap(tabla(["Serie", "Nombre", "Matiz", "Clase", "Claro: hex · L · C · contraste mín.", "Oscuro: hex · L · C · contraste mín."], filas))
    ap("\nTokens `color.chart.<n>.light` y `color.chart.<n>.dark`. Las series se asignan en este orden fijo, nunca se reciclan: una novena serie se agrupa en «Otros» (neutral 600 en claro, neutral 400 en oscuro, que cumplen 3:1 como los bordes de campo).\n")
    ap("### 8.2 Regla de luminosidad\n")
    nombres_clase = {c: ", ".join(M.graf[k]["nombre"] for k in range(1, 9) if M.graf[k]["clase"] == c) for c in ("A", "B")}
    R = ranking_orden(M)
    nom_orden = lambda orden: ", ".join(M.graf[k]["nombre"] for k in orden)
    f1, f7, f8 = (M.graf[k]["nombre"] for k in FIJAS_ORDEN)
    act, r0 = R["actual"], R["ronda0"]
    cambio = [k for k in range(1, 9) if r0["orden"][k - 1] != k]
    if len(cambio) == 2:
        cambio_txt = "pone %s en la serie %d y %s en la %d" % (M.graf[cambio[0]]["nombre"], cambio[0], M.graf[cambio[1]]["nombre"], cambio[1])
    else:
        cambio_txt = "cambia las series %s" % ", ".join(str(k) for k in cambio)
    orden_txt = ("En la ronda 1 de revisión (hallazgo H-03) cambió el matiz de la serie %d (§8.8) y el script repite la búsqueda con los matices vigentes (%d órdenes posibles; tabla siguiente). Como tercer criterio, solo para desempatar, se agrega el ΔE en visión normal entre todos los pares de las tres primeras series, las únicas que se usan en dispersión, mapas y pequeños múltiplos (§8.5). El orden vigente queda en el puesto %d."
                 % (SERIE_CORREGIDA, len(R["lista"]), R["puesto_actual"]))
    if cambio:
        orden_txt += (" Respecto a la ronda 0 %s: con el orden anterior y el matiz nuevo, el mínimo ΔL consecutivo bajaría a %.3f (puesto %d) y las tres primeras series quedarían con ΔE normal mínimo %.1f entre todos sus pares, bajo el piso %.0f de la skill dataviz; con el orden vigente son %.3f y %.1f."
                      % (cambio_txt, r0["dl"], R["puesto_ronda0"], r0["tres_n"], NORMAL_PISO, act["dl"], act["tres_n"]))
    for x in R["empates"]:
        orden_txt += " Empata en los dos primeros criterios con el orden %s, que pierde en el tercero (%.1f frente a %.1f)." % (nom_orden(x["orden"]), x["tres_n"], act["tres_n"])
    ap("Las series alternan dos clases de luminosidad para que las vecinas se separen en gris y con daltonismo: clase A (alta) L %.3f en claro y %.2f en oscuro; clase B (baja) L %.2f en claro y %.2f en oscuro. Las impares son B y las pares A. Los matices de clase B (%s) conservan croma en luminosidad baja y los de clase A (%s) en luminosidad alta, por eso la clase de cada matiz es la misma en los dos temas. Los límites salen del contraste: en claro la clase A queda bajo el máximo que da 3:1 sobre la tarjeta 100 y en oscuro la clase B sobre el mínimo que da 3:1 sobre la tarjeta 900. Croma pedido %.2f en claro y %.2f en oscuro; todas las series quedan con croma final ≥ %.2f (piso de la skill dataviz para que el matiz siga identificando). El orden se eligió en la ronda 0 por búsqueda exhaustiva entre las permutaciones que alternan clases, con %s fijo en la serie 1 (gráficas de una sola serie) y %s y %s al final (para que las primeras cinco tengan matices bien distintos y el verde no sugiera «positivo» en gráficas comunes), maximizando el mínimo ΔL consecutivo en las cuatro vistas y, después, el ΔE en visión normal entre vecinas. %s\n" % (L_CLASE[("claro", "A")], L_CLASE[("oscuro", "A")], L_CLASE[("claro", "B")], L_CLASE[("oscuro", "B")], nombres_clase["B"], nombres_clase["A"], CROMA_GRAFICA["claro"], CROMA_GRAFICA["oscuro"], CROMA_PISO_GRAFICAS, f1, f7, f8, orden_txt))
    filas = []
    for i, d in enumerate(R["lista"][:4], 1):
        nota = ", ".join(n for n, ok in (("vigente", d["orden"] == act["orden"]), ("ronda 0 con la serie %d nueva" % SERIE_CORREGIDA, d["orden"] == r0["orden"])) if ok)
        filas.append([i, nom_orden(d["orden"]), "%.3f" % d["dl"], "%.1f" % d["de"], "%.1f" % d["tres_n"], "%.1f" % d["tres_cvd"], nota])
    ap(tabla(["Puesto", "Orden (series 1 a 8)", "ΔL consecutivo mín.", "ΔE normal mín. entre vecinas", "ΔE normal mín. entre las 3 primeras", "ΔE protan/deutan mín. entre las 3 primeras", "Nota"], filas))
    ap("\nLos cuatro primeros de %d órdenes; cada valor es el mínimo de los dos temas.\n" % len(R["lista"]))
    ap("### 8.3 Distancias entre series consecutivas\n")
    for tema in TEMAS:
        d = G[tema]
        ap("**Tema %s** (series sobre fondo %s y tarjeta %s)\n" % (tema, M.sup(tema, "fondo"), M.sup(tema, "tarjeta")))
        filas = []
        for f in d["cons"]:
            filas.append([f["par"], "%.3f" % f["dl_gris"], "%.3f" % f["dl_protanopia"], "%.3f" % f["dl_deuteranopia"], "%.3f" % f["dl_tritanopia"],
                          "%.1f" % f["de_normal"], "%.1f" % f["de_protanopia"], "%.1f" % f["de_deuteranopia"], "%.1f" % f["de_tritanopia"]])
        ap(tabla(["Par", "ΔL gris", "ΔL protanopia", "ΔL deuteranopia", "ΔL tritanopia", "ΔE normal", "ΔE protanopia", "ΔE deuteranopia", "ΔE tritanopia"], filas))
        ap("\nMínimo ΔL consecutivo en gris y simulaciones: **%.3f** (umbral %.2f; el encargo pide %.2f). ΔE mínimo entre vecinas: normal %.1f (piso %.0f), protan/deutan %.1f (objetivo %.0f), tritan %.1f.\n" % (d["min_dl"], UMBRAL_DL, UMBRAL_DL_ENCARGO, d["min_de_normal"], NORMAL_PISO, d["min_de_cvd"], CVD_OBJETIVO, d["min_de_tritan"]))
    ap("### 8.4 Hex simulados y luminosidad por serie\n")
    for tema in TEMAS:
        ap("**Tema %s** (hex · L; en gris, L = Y^(1/3))\n" % tema)
        filas = []
        for n, hx in enumerate(G[tema]["cols"], 1):
            filas.append([n, "`%s` · %.3f" % (hx, hex_a_oklab(hx)[0]), "`%s` · %.3f" % (hex_gris(hx), l_gris(hx))] +
                         ["`%s` · %.3f" % (hex_simulado(hx, t), l_vista(hx, t)) for t in TIPOS_CVD])
        ap(tabla(["Serie", "Normal", "Gris", "Protanopia", "Deuteranopia", "Tritanopia"], filas))
        ap("")
    ap("ΔL en gris entre los 28 pares (informativo por K27: el umbral duro es ΔL ≥ 0.10 entre series consecutivas, §8.3, con forma del marcador, patrón de trazo o relleno y etiqueta directa como refuerzo; bloqueo `criterios-deltal` resuelto):\n")
    for tema in TEMAS:
        cols = G[tema]["cols"]
        filas = []
        for i in range(8):
            filas.append([i + 1] + ["%.3f" % abs(l_gris(cols[i]) - l_gris(cols[j])) if j < i else "" for j in range(8)])
        ap("**Tema %s**\n" % tema)
        ap(tabla(["Serie"] + [str(j + 1) for j in range(8)], filas))
        n_bajo = sum(1 for a, b in combinations(cols, 2) if abs(l_gris(a) - l_gris(b)) < UMBRAL_DL)
        ap("\nPares con ΔL gris < %.2f: %d de 28 (todos no consecutivos).\n" % (UMBRAL_DL, n_bajo))
    ap("### 8.5 Pares no consecutivos\n")
    ap("El criterio U2 del plan §3.5 dice «ΔL ≥ 0.10 entre series». Entre **todas** las series es imposible con 8 series: harían falta 7 saltos de 0.10, un rango de 0.70 en L, y el contraste de 3:1 deja como mucho L_gris ≤ %.3f en claro (rango [0, %.3f] aun aceptando negro) y L_gris ≥ %.3f en oscuro (rango [%.3f, 1], %.3f). Por eso el criterio se aplica, como pide el encargo, a series **consecutivas** del orden de uso, que son las que se tocan en barras apiladas, barras agrupadas y leyendas. Mínimos entre todos los pares (informativo): claro ΔE normal %.1f y ΔL gris %.3f; oscuro ΔE normal %.1f y ΔL gris %.3f. Consecuencia de uso: en dispersión, mapas o pequeños múltiplos (donde cualquier par puede quedar junto) se usan como máximo las 3 primeras series, y toda gráfica con 2 o más series lleva leyenda y, con ≤ 4, etiqueta directa; en impresión en gris o modo de alto contraste se activa la textura de trama.\n" % (V["l_claro_max"], V["l_claro_max"], V["l_oscuro_min"], V["l_oscuro_min"], 1 - V["l_oscuro_min"], G["claro"]["todos"]["normal"], G["claro"]["dl_todos_gris"], G["oscuro"]["todos"]["normal"], G["oscuro"]["dl_todos_gris"]))
    tres = {tema: (min(delta_e(a, b) for a, b in combinations(G[tema]["cols"][:3], 2)),
                   min(min(delta_e(a, b, "protanopia"), delta_e(a, b, "deuteranopia")) for a, b in combinations(G[tema]["cols"][:3], 2))) for tema in TEMAS}
    ap("Las tres primeras series pasan la prueba de todos los pares de la skill dataviz (`--pairs all`): ΔE mínimo entre ellas de %.1f en visión normal y %.1f con protanopia o deuteranopia en claro, y de %.1f y %.1f en oscuro (pisos %.0f y %.0f). El script lo verifica.\n"
       % (tres["claro"][0], tres["claro"][1], tres["oscuro"][0], tres["oscuro"][1], NORMAL_PISO, CVD_OBJETIVO))
    ap("### 8.6 Robustez frente a la implementación de la simulación\n")
    ap("Con la variante que aplica las matrices sobre sRGB codificado (la de culori 4.0.2), el mínimo ΔL consecutivo es %.3f en claro y %.3f en oscuro: la paleta también pasa el umbral %.2f si U2b o U5 simulan con culori sin linealizar. Recomendación para ellos: convertir a `lrgb` antes de aplicar la matriz, para coincidir con este documento.\n" % (Gg["claro"]["min_dl"], Gg["oscuro"]["min_dl"], UMBRAL_DL))
    ap("### 8.7 Banda de luminosidad de la skill dataviz\n")
    ap("La skill propone L 0.43–0.77 en claro y 0.48–0.67 en oscuro, calibrada para su superficie oscura `#1a1a19`. Con la tarjeta 900 de este sistema (L 0.335), 3:1 exige L ≥ %.3f, así que su banda oscura dejaría solo [%.2f, 0.67], insuficiente para alternar con ΔL ≥ 0.10. Se adoptan sus demás controles (croma ≥ 0.10, ΔE de vecinas con daltonismo ≥ 8 y en visión normal ≥ 15, contraste ≥ 3:1) y se reemplaza la banda por los límites de contraste de §8.2.\n" % (V["l_oscuro_min"], V["l_oscuro_min"]))
    ap("### 8.8 Matiz reservado y elección de la serie %d\n" % SERIE_CORREGIDA)
    B = barrido_serie(M)
    g6 = M.graf[SERIE_CORREGIDA]
    elegido = next((i for i, d in enumerate(B["pasan"]) if d["h"] == g6["h"]), None)
    desc_txt = "; ".join("%s: %d" % (MOTIVOS_DESCARTE[m], B["motivos"][m]) for m in ("croma", "zona", "contraste", "dl", "de") if m in B["motivos"])
    ap("Decisiones §I reserva el eje opuesto 100° para «acción crítica» de confirmación, y la zona reservada de §3 (%d°–%d°) vale también para las series. Hasta la ronda 1 de revisión la serie %d era «mostaza» 100° (`#978500` en claro, `#DDCC6B` en oscuro), dentro de esa zona y sin excepción declarada (hallazgo H-03 de la revisión independiente de WF-0). Para elegir su nuevo matiz, el script barre los múltiplos de %d° fuera de la zona (%d candidatos; %d quedan dentro de ella) con la clase de la serie (%s) y las otras siete series en su lugar, y repite los controles de esta sección en ambos temas: croma final ≥ %.2f, matiz medido en el hex también fuera de la zona, contraste ≥ 3:1 sobre fondo y tarjeta, ΔL ≥ %.2f con las vecinas (series %d y %d) en gris y en las tres simulaciones con las dos variantes de §8.6, y ΔE de vecinas ≥ %.0f en visión normal y ≥ %.0f con protanopia y deuteranopia. Pasan %d. Descartes por motivo (un matiz puede tener varios): %s.\n"
       % (MATIZ_OPUESTO - ZONA_RESERVADA, MATIZ_OPUESTO + ZONA_RESERVADA, SERIE_CORREGIDA, PASO_BARRIDO, B["evaluados"], B["en_zona"], g6["clase"],
          CROMA_PISO_GRAFICAS, UMBRAL_DL, SERIE_CORREGIDA - 1, SERIE_CORREGIDA + 1, NORMAL_PISO, CVD_OBJETIVO, len(B["pasan"]), desc_txt))
    ap("Entre los que pasan se elige el que más se aleja en matiz de las otras siete series: una paleta categórica identifica cada serie por una familia de matiz distinta (control 1 de la skill dataviz, ocho familias en orden fijo), y las series de la misma clase solo se distinguen por matiz cuando no son vecinas (§8.5). El desempate es el mayor ΔL con las vecinas. Los cinco primeros:\n")
    filas = []
    for i, d in enumerate(B["pasan"][:5], 1):
        filas.append([i, "%d°" % d["h"], "`%s` · C %.3f" % (d["hex"]["claro"], d["croma"]["claro"]), "`%s` · C %.3f" % (d["hex"]["oscuro"], d["croma"]["oscuro"]),
                      "%.1f°" % d["dh"], "%.3f" % d["dl"], "%.1f" % d["de"], "**elegido**" if d["h"] == g6["h"] else ""])
    ap(tabla(["Puesto", "Matiz", "Claro: hex · C", "Oscuro: hex · C", "Δ matiz mín. a las otras series", "ΔL mín. con vecinas", "ΔE normal mín. con vecinas", "Elección"], filas))
    if elegido is None:
        eleg_txt = "La serie %d usa %d°, que no está entre los candidatos que pasan: revisar el barrido." % (SERIE_CORREGIDA, g6["h"])
    else:
        eleg_txt = "Se usa el puesto %d, %d° («%s»): queda a %.1f° de la serie más cercana en matiz y a %.1f° del opuesto." % (
            elegido + 1, g6["h"], g6["nombre"], B["pasan"][elegido]["dh"], min(dist_reserva(g6[t].hex, g6["h"]) for t in TEMAS))
    if B["descartados"]:
        dd, mm = B["descartados"][0]
        eleg_txt += " El descartado con más separación de matiz es %d° (%.1f°), que falla por %s." % (
            dd["h"], dd["dh"], ", ".join(MOTIVOS_DESCARTE[m] for m in ("croma", "zona", "contraste", "dl", "de") if m in mm))
    matiz_rol = {r["codigo"]: r["matiz"] for r in ROLES}
    roles_cerca = []
    for tema in TEMAS:
        hm = hex_a_oklch(g6[tema].hex)[2]
        roles_cerca += [(dif_matiz(hm, hex_a_oklch(M.hex("%s.%s" % (base_rol(r["codigo"]), SUFIJO[tema])))[2]), r["nombre"].lower()) for r in ROLES]
    d_rol, n_rol = min(roles_cerca)
    d_vio, d_ver = dif_matiz(M.graf[7]["h"], matiz_rol["debt"]), dif_matiz(M.graf[8]["h"], matiz_rol["income"])
    eleg_txt += (" Como el resto de las series, la %d no se separa de los roles: queda a %.0f° de %s (rol); la serie 7 (%s) %s y la 8 (%s) queda a %.0f° de ingreso. Los roles viven en cifras, insignias y avisos, siempre con signo, icono o texto; las series viven en gráficas con leyenda, y cuando una serie significa un estado (por ejemplo, presupuesto excedido) usa el token de rol y no el de serie (regla de la skill dataviz)."
                 % (SERIE_CORREGIDA, d_rol, n_rol, M.graf[7]["nombre"], "comparte matiz con deuda" if d_vio < 0.5 else "queda a %.0f° de deuda" % d_vio, M.graf[8]["nombre"], d_ver))
    ap("\n" + eleg_txt + "\n")

    # 9 ------------------------------------------------------------------------------
    ap("## 9. Reglas de uso por componente\n")
    def mn(prefijo, tema):
        vals = [p["ratio"] for p in pares if p["tema"] == tema and p["uso"].startswith(prefijo)]
        return r2(min(vals)) if vals else "—"
    filas = [
        ["Fondo de página", "neutral 50", "neutral 950", "—"],
        ["Tarjeta", "neutral 100", "neutral 900", "—"],
        ["Divisor decorativo", "neutral 200", "neutral 800", "exento (decorativo)"],
        ["Texto normal", "neutral 900", "neutral 100", "%s / %s" % (mn("Texto normal", "claro"), mn("Texto normal", "oscuro"))],
        ["Texto secundario", "neutral 700", "neutral 300", "%s / %s" % (mn("Texto secundario", "claro"), mn("Texto secundario", "oscuro"))],
        ["Texto deshabilitado", "neutral 700 + relleno neutral 100", "neutral 400 + relleno neutral 900", "%s / %s" % (mn("Texto de", "claro"), mn("Texto de", "oscuro"))],
        ["Marcador de posición", "neutral 700", "neutral 300", "%s / %s" % (mn("Marcador", "claro"), mn("Marcador", "oscuro"))],
        ["Enlace (subrayado) · hover", "brand 700 · brand 800", "brand 300 · brand 200", "%s / %s" % (mn("Enlace", "claro"), mn("Enlace", "oscuro"))],
        ["Botón primario · hover", "relleno brand 700 · 800, texto blanco", "relleno brand 300 · 200, texto brand 950", "%s / %s" % (mn("Botón primario", "claro"), mn("Botón primario", "oscuro"))],
        ["Botón secundario · hover", "borde neutral 600, texto neutral 900 · relleno neutral 200", "borde neutral 400, texto neutral 100 · relleno neutral 800", "%s / %s" % (mn("Botón secundario", "claro"), mn("Botón secundario", "oscuro"))],
        ["Botón destructivo · hover", "relleno error 700 · 800, texto blanco", "relleno error 300 · 200, texto error 950", "%s / %s" % (mn("Botón destructivo", "claro"), mn("Botón destructivo", "oscuro"))],
        ["Botón de acción crítica · hover", "relleno critical 700 · 800, texto blanco", "relleno critical 300 · 200, texto critical 950", "%s / %s" % (mn("Botón de acción", "claro"), mn("Botón de acción", "oscuro"))],
        ["Borde de campo", "neutral 600", "neutral 400", "%s / %s" % (mn("Borde de campo", "claro"), mn("Borde de campo", "oscuro"))],
        ["Anillo de foco (2 px, separación 2 px)", "brand 600", "brand 400", "%s / %s" % (mn("Anillo", "claro"), mn("Anillo", "oscuro"))],
        ["Elemento activo, fila seleccionada", "relleno brand 100, texto neutral 900, icono brand 700", "relleno brand 900, texto neutral 100, icono brand 300", "%s / %s" % (mn("Elemento activo", "claro"), mn("Elemento activo", "oscuro"))],
        ["Insignia de rol", "texto rol 700 sobre rol 100", "texto rol 400 sobre rol 900", "%s / %s" % (mn("Insignia", "claro"), mn("Insignia", "oscuro"))],
        ["Cifra de rol (con signo e icono)", "rol 700", "rol 400", "%s / %s" % (mn("Cifra", "claro"), mn("Cifra", "oscuro"))],
        ["Series de gráficas", "`color.chart.<n>.light`", "`color.chart.<n>.dark`", "%s / %s" % (mn("Serie", "claro"), mn("Serie", "oscuro"))],
        ["Personas y tiendas", "`color.person|store.*.light`", "`color.person|store.*.dark`", "%s / %s" % (r2(min(C["person"]["claro"]["contraste_min"], C["store"]["claro"]["contraste_min"])), r2(min(C["person"]["oscuro"]["contraste_min"], C["store"]["oscuro"]["contraste_min"])))],
    ]
    ap(tabla(["Elemento", "Claro", "Oscuro", "Contraste mínimo claro / oscuro"], filas))
    ap("\nNotas: (1) «texto 950» del botón primario en oscuro se lee como brand 950, el 950 de la misma escala del relleno; con neutral 950 daría %s:1, también válido. (2) WCAG 1.4.3 exime al texto de controles deshabilitados, pero el encargo exige 4.5:1 a todo texto; en claro no hay paso entre 600 (%s:1 sobre fondo, insuficiente) y 700, así que el deshabilitado usa el mismo neutral 700 que el secundario y se reconoce por el relleno neutral 100, el borde punteado y el cursor. (3) Los divisores y bordes de tarjeta (neutral 200 en claro, 800 en oscuro) son decorativos: los campos y botones tienen su propio borde de 3:1. (4) Los enlaces van subrayados para no depender del color frente al texto normal.\n" % (r2(contraste(M.neutro[950].hex, M.marca[300].hex)), r2(contraste(M.neutro[600].hex, n50))))

    # 10 -----------------------------------------------------------------------------
    ap("## 10. Contraste (resumen de `contraste.csv`)\n")
    ap("Columnas: `primer_plano`, `fondo` (rutas de token), `tema`, `uso`, `ratio` (truncada a 2 decimales), `umbral`, `cumple`, y además `primer_plano_hex` y `fondo_hex` para volver a medir sin el JSON.\n")
    filas = []
    for tema in TEMAS:
        for umbral, nombre in ((UMBRAL_TEXTO, "texto"), (UMBRAL_UI, "interfaz y gráficas")):
            sub = [p for p in pares if p["tema"] == tema and p["umbral"] == umbral]
            filas.append([tema, nombre, len(sub), sum(1 for p in sub if p["cumple"]), r2(min(p["ratio"] for p in sub))])
    ap(tabla(["Tema", "Tipo", "Pares", "Cumplen", "Razón mínima"], filas))
    ap("\nTotal: **%d pares, %d cumplen, %d bajo umbral.** Par de texto más ajustado: %s (%s, %s:1). Par de interfaz más ajustado: %s (%s, %s:1).\n" % (total, sum(1 for p in pares if p["cumple"]), sum(1 for p in pares if not p["cumple"]), peor_t["uso"], peor_t["tema"], r2(peor_t["ratio"]), peor_u["uso"], peor_u["tema"], r2(peor_u["ratio"])))

    ap("**Pares estructurales del plan §4.3, recalculados** (tolerancia ±0.1; si la diferencia es mayor se documentan ambos valores):\n")
    filas = []
    for desc, fp, fo, (lo, hi), txt in ESTRUCTURALES:
        r = contraste(M.hex(fp), M.hex(fo))
        dif = 0.0 if lo <= r <= hi else min(abs(r - lo), abs(r - hi))
        filas.append([desc, "`%s` / `%s`" % (M.hex(fp), M.hex(fo)), txt, r2(r), "dentro" if dif <= 0.1 else "**difiere %.2f**" % dif])
    ap(tabla(["Par", "Hex", "Plan", "Medido", "Tolerancia"], filas))
    ap("\nDonde la diferencia supera 0.1 se documentan ambos valores y rige el medido. En particular, el borde de campo neutral 600 sobre fondo 50 da %s:1, no ≥ 4.2:1 como indicaba el plan; cumple el 3:1 que exige el borde de un componente (criterio 1.4.11).\n" % r2(contraste(M.neutro[600].hex, n50)))

    # 11 -----------------------------------------------------------------------------
    ap("## 11. Tokens (`C2-colores.json`)\n")
    ap("Formato DTCG con `$value` en hex sRGB, como pide el encargo. Cada token lleva `$type: \"color\"`, `$value` y `$extensions.oklch` con `l`, `c` (6 decimales) y `h` de diseño, antes de cuantizar a 8 bits; `c` es el croma ya recortado al gamut. Con 6 decimales, convertir esos valores a hex reproduce exactamente `$value` (con 4 decimales un canal de `color.role.warning.50` cruza el medio punto de redondeo). Grupos:\n")
    ap(tabla(["Grupo", "Tokens", "Contenido"], [
        ["`color.base.white`", 1, "blanco puro para texto sobre rellenos 700 en claro"],
        ["`color.brand.50…950`", 11, "marca 280°"],
        ["`color.neutral.50…950`", 11, "neutro tintado 280°"],
        ["`color.role.<rol>.{light,dark,surface-light,surface-dark,50…950}`", 8 * 15, "8 roles: income, expense, savings, debt, warning, error, success, info"],
        ["`color.critical.{light,dark,surface-light,surface-dark,50…950}`", 15, "acción crítica, eje opuesto 100°"],
        ["`color.person.<persona>.{light,dark}`", 8, "axel, ximena, karen, shared (conjunto)"],
        ["`color.store.<tienda>.{light,dark}`", 8, "walmart, bodega-aurrera, sams-club, market (mercado)"],
        ["`color.chart.<1…8>.{light,dark}`", 16, "series de gráficas en orden de uso"],
    ]))
    ap("\n`light`/`dark` de cada rol duplican el paso elegido (700/400) y `surface-*` los pasos 100/900, para que U2b mapee variables por tema sin conocer la regla. Los nombres de código de roles, personas y tiendas están en inglés (protocolo §8); si el diccionario de datos de S1 fija otros códigos, se agregan alias en WF-1 sin cambiar valores.\n")

    # 12 -----------------------------------------------------------------------------
    ap("## 12. Hallazgos, decisiones derivadas del cálculo y bloqueos\n")
    ap("**Hallazgos**\n")
    hallazgos = [
        "Los 22 valores de partida del plan §4.3 (marca y neutro) se reproducen exactamente con la tabla de perfil de §2.3; el plan no la explicitaba.",
        "Los valores de rol del plan (L 0.56) quedan entre %s y %s:1 sobre las superficies claras; %d de %d quedan bajo 4.5:1 en al menos una de las dos. Se usan pasos de escala (claro 700, oscuro 400); el paso 600 que pedía el encargo da %s–%s:1 (§6.2)." % (r2(min(plan_rol_cr)), r2(max(plan_rol_cr)), plan_rol_bajo, len(plan_rol_cr), r2(min(p600)), r2(max(p600))),
        "Versión 1.1 (K26): gasto pasa de %d° a %d° y deja de confundirse con error 20° (ΔE normal %.1f en claro y %.1f en oscuro; en la 1.0 eran %.1f y %.1f). Bloqueo `U2-gasto-error` resuelto (opción 2)." % (ANTES_11["expense"], h_gasto, ge_c["normal"], ge_o["normal"], ge0_c, ge0_o),
        "Además de las dos colisiones listadas en decisiones §I, la regla de separación detecta tres más en la propuesta categórica del plan: Axel–información, Mercado–gasto/error y Walmart en el límite exacto de 20°. Las cinco quedan resueltas (§7.5).",
        "%d colores categóricos del plan no llegaban a 3:1 sobre las superficies claras y %d sobre las oscuras; no existe un valor único que cumpla 3:1 en ambos temas y además se escalone (ventana L_gris [%.3f, %.3f]). Personas, tiendas y series tienen variante por tema (§7, §8)." % (cat_cl, cat_os, V["l_min"], V["l_max"]),
        "«ΔL ≥ 0.10 entre series» (plan §3.5) solo es alcanzable entre series consecutivas; entre todos los pares es matemáticamente imposible con 8 series (§8.5). K27 fija esa lectura: umbral duro entre series consecutivas y los 28 pares como información (bloqueo `criterios-deltal` resuelto).",
        "culori 4.0.2 aplica las matrices de Machado sobre sRGB codificado, no lineal; la paleta pasa con ambas variantes (§8.6).",
        "Pares de roles con ΔE normal < 15 (ingreso–éxito, advertencia–acción crítica, entre otros, §6.6): ninguno comparte componente y todos llevan icono y texto.",
        "Ronda 1 de revisión (H-03): la serie %d de gráficas («mostaza», 100°) usaba el matiz que decisiones §I reserva para acción crítica y el documento no declaraba la excepción. Pasa a «%s» %d° (§8.8); con ese matiz, la búsqueda de orden de §8.2 %s. El script verifica la zona reservada %d°–%d° en todos los tokens (§3)." % (
            SERIE_CORREGIDA, M.graf[SERIE_CORREGIDA]["nombre"], M.graf[SERIE_CORREGIDA]["h"], cambio_txt if cambio else "conserva el orden de la ronda 0", MATIZ_OPUESTO - ZONA_RESERVADA, MATIZ_OPUESTO + ZONA_RESERVADA),
    ]
    if cambio and r0["tres_n"] < NORMAL_PISO:
        p_os = r0["tres"]["oscuro"]
        a_, b_ = M.graf[p_os[1]], M.graf[p_os[2]]
        par_txt = "entre %s y %s%s" % (a_["nombre"], b_["nombre"], (", ambas de clase %s y con la misma luminosidad" % a_["clase"]) if a_["clase"] == b_["clase"] else "")
        hallazgos.append("Encontrado al corregir H-03: con el orden de la ronda 0, las tres primeras series (%s), que §8.5 reserva para dispersión, mapas y pequeños múltiplos, no llegaban al piso de ΔE %.0f entre todos sus pares en visión normal (%.1f en claro y %.1f en oscuro, %s). Con el orden vigente dan %.1f y %.1f, y el script lo verifica (§8.5)." % (
            nom_orden(r0["orden"][:3]), NORMAL_PISO, r0["tres"]["claro"][0], r0["tres"]["oscuro"][0], par_txt, act["tres"]["claro"][0], act["tres"]["oscuro"][0]))
    mp = [f for f in cercanos if "(marca)" in f["a"] + f["b"]]
    hallazgos += [
        "Versión 1.1 (K26): con gasto en %d° no queda hueco cálido a ≥ 20° de error, gasto y advertencia, y la L de Karen y de Bodega Aurrera las obliga a separarse de los roles por matiz; sus ventanas factibles quedan entre deuda y error y van a %d° y %d° (Cambios 1.1)." % (h_gasto, M.cat["person"]["karen"]["h"], M.cat["store"]["bodega-aurrera"]["h"]),
        "Versión 1.1 (R3-04): con la marca en el barrido de §7.7 quedan cerca %s; todos con refuerzo de forma y etiqueta, y ningún rol cerca de la marca." % "; ".join("%s y %s en %s (%.1f°, ΔL %.3f)" % (f["a"], f["b"], f["tema"], f["dm"], f["dl"]) for f in mp),
        "Versión 1.1: en oscuro, Bodega Aurrera en rosa no cabe en el escalonamiento ΔL ≥ 0.10 de las tiendas con Walmart en L %.3f (ningún matiz entero cumple); Walmart pasa a L %.3f, PROVISIONAL con el bloqueo `U2a-walmart-oscuro`, y Bodega Aurrera usa croma %.2f (Cambios 1.1)." % (ANTES_11[("store", "walmart")][3], M.cat["store"]["walmart"]["oscuro"].l, M.cat["store"]["bodega-aurrera"]["c"]),
        hallazgo_rosas_ximena(M),
    ]
    for i, h in enumerate(hallazgos, 1):
        ap("%d. %s" % (i, h))
    ap("\n**Decisiones derivadas del cálculo** (no son supuestos: las fuerza un requisito explícito)\n")
    ap("- Paso 700 en claro para los roles (el 600 incumple 4.5:1).\n- Variantes `light`/`dark` para personas, tiendas y series (un valor único no cumple 3:1 en ambos temas con escalonamiento).\n- Criterio de ΔL aplicado a series consecutivas, con umbral 0.10 (más estricto que el 0.08 del encargo).\n- Texto deshabilitado en claro con neutral 700 (el 600 incumple 4.5:1).\n- Versión 1.1: Bodega Aurrera con croma %.2f y Karen con L %.3f en claro (márgenes del escalonamiento de su grupo; Cambios 1.1).\n" % (M.cat["store"]["bodega-aurrera"]["c"], M.cat["person"]["karen"]["claro"].l))
    ap("**Propuestas sujetas a la puerta 2:** ingreso 135° (§6.4, opción prevista en decisiones §I), la resolución de colisiones categóricas (§7.5; en la versión 1.1, gasto 45°, Karen y Bodega Aurrera en sus matices nuevos, «Cambios 1.1») y la paleta de gráficas corregida en la ronda 1 (serie %d en %d°%s; §8.2 y §8.8).\n" % (SERIE_CORREGIDA, M.graf[SERIE_CORREGIDA]["h"], (", con %s en la serie %d y %s en la %d" % (M.graf[cambio[0]]["nombre"], cambio[0], M.graf[cambio[1]]["nombre"], cambio[1])) if len(cambio) == 2 else ""))
    ap("**Bloqueos:** `docs/estado/bloqueos/U2a-walmart-oscuro.md` (versión 1.1: Walmart en oscuro con L %.3f en vez de %.3f, PROVISIONAL; «Cambios 1.1») y `docs/estado/bloqueos/%s.md` (ronda 1 de revisión, H-01: Karen y Bodega Aurrera bajo ≥ 20° o ΔL ≥ 0.15 frente a Ximena; %d tokens PROVISIONALES; «Cambios 1.1», punto 2). Resueltos en la puerta 1: `U2-gasto-error` (K26) y `criterios-deltal` (K27).\n" % (M.cat["store"]["walmart"]["oscuro"].l, ANTES_11[("store", "walmart")][3], BLOQUEO_ROSAS, len(rosas_ximena(M)["bajo"])))

    # Anexo ---------------------------------------------------------------------------
    ap("## Anexo A. Tokens: ruta → hex\n")
    ap("Lista completa y verificable de `C2-colores.json` (el script compara esta tabla con el JSON en ambos sentidos).\n")
    filas = [["`%s`" % ruta, "`%s`" % t["hex"], "%.6f" % t["l"], "%.6f" % t["c"], "%.2f" % t["h"]] for ruta, t in M.tokens.items()]
    ap(tabla(["Token", "Hex", "L", "C", "h"], filas))
    ap("")
    return "\n".join(o) + "\n"


# =============================================================================
# 8. Escritura: muestrario HTML
# =============================================================================
CSS = """
*{box-sizing:border-box}
body{margin:0;background:#F2F2F4;color:#36363A;font:13px/1.4 system-ui,-apple-system,"Segoe UI",Roboto,"Noto Sans",Ubuntu,sans-serif}
.wrap{max-width:1380px;margin:0 auto;padding:14px 16px 24px}
h1{font-size:20px;margin:0 0 2px}
h2{font-size:15px;margin:16px 0 6px}
h3{font-size:13px;margin:10px 0 6px}
p.nota{margin:0 0 8px;color:#66676E}
.mono{font-family:ui-monospace,"SFMono-Regular",Menlo,Consolas,monospace}
.escala{display:grid;grid-template-columns:150px repeat(11,1fr);gap:3px;margin-bottom:3px}
.escala .et{font-weight:600;font-size:12px;display:flex;align-items:center}
.sw{height:30px;border-radius:4px;padding:2px 4px;font:10px/1.15 ui-monospace,Menlo,Consolas,monospace;display:flex;flex-direction:column;justify-content:center}
.temas{display:grid;grid-template-columns:1fr 1fr;gap:12px}
.panel{border-radius:10px;padding:12px}
.tarjeta{border-radius:8px;padding:10px;margin-top:8px}
.fila{display:flex;flex-wrap:wrap;gap:6px;align-items:center;margin:4px 0}
.btn{display:inline-block;padding:6px 12px;border-radius:6px;font-weight:600;font-size:12px;border:1px solid transparent}
.campo{display:inline-block;width:190px;padding:6px 8px;border-radius:6px;font-size:12px}
.badge{display:inline-block;padding:2px 8px;border-radius:999px;font-size:11px;font-weight:600}
.punto{display:inline-block;width:11px;height:11px;border-radius:50%;margin-right:4px;vertical-align:-1px}
.chip{display:inline-flex;align-items:center;margin-right:10px;font-size:12px}
.cifra{font-variant-numeric:tabular-nums;font-weight:600;margin-right:14px}
.sim{display:grid;grid-template-columns:96px repeat(8,1fr) 8px repeat(4,1fr) 8px repeat(4,1fr) 64px;gap:2px;align-items:center;margin-bottom:2px}
.sim .c{height:20px;border-radius:3px}
.sim .et{font-size:11px}
.sim .v{font:10px ui-monospace,Menlo,monospace;text-align:right}
.tabla{columns:4;column-gap:14px;font-size:10px;line-height:1.3}
.tabla div{break-inside:avoid;display:flex;gap:4px;border-bottom:1px solid #E5E6E9;padding:1px 0}
.tabla .u{flex:1}
.tabla .r{font-family:ui-monospace,Menlo,monospace;white-space:nowrap}
.ok{color:#207B3E;font-weight:700}
.mal{color:#A34449;font-weight:700}
"""


def _texto_sobre(hx):
    return "#FFFFFF" if contraste("#FFFFFF", hx) >= contraste("#222325", hx) else "#222325"


def _sw(hx, etiqueta):
    return '<div class="sw" style="background:%s;color:%s"><b>%s</b><span>%s</span></div>' % (hx, _texto_sobre(hx), _html.escape(str(etiqueta)), hx)


def _svg_series(cols, fondo, texto2):
    anchos = [18, 15, 14, 13, 12, 10, 9, 9]
    x, w_total, partes = 0.0, 600.0, []
    for c, a in zip(cols, anchos):
        w = w_total * a / 100.0
        partes.append('<rect x="%.1f" y="0" width="%.1f" height="18" rx="3" fill="%s"/>' % (x, max(0.0, w - 2), c))
        x += w
    alturas = [92, 78, 70, 62, 54, 46, 38, 30]
    for i, (c, h) in enumerate(zip(cols, alturas)):
        bx = 12 + i * 74
        partes.append('<rect x="%d" y="%d" width="48" height="%d" rx="3" fill="%s"/>' % (bx, 128 - h, h, c))
        partes.append('<text x="%d" y="142" font-size="11" fill="%s" text-anchor="middle">%d</text>' % (bx + 24, texto2, i + 1))
    return '<svg viewBox="0 0 600 146" width="100%%" role="img" aria-label="Series de gráfica" style="display:block;height:auto;background:%s">%s</svg>' % (fondo, "".join(partes))


def _panel_tema(M, tema, G, C):
    sfx = SUFIJO[tema]
    claro = tema == "claro"
    f, t = M.sup(tema, "fondo"), M.sup(tema, "tarjeta")
    nh = lambda p: M.neutro[p].hex
    bh = lambda p: M.marca[p].hex
    eh = lambda p: M.esc_rol["error"][p].hex
    kh = lambda p: M.esc_rol["critical"][p].hex
    txt, sec, des = (nh(900), nh(700), nh(700)) if claro else (nh(100), nh(300), nh(400))
    enl, enl_h = (bh(700), bh(800)) if claro else (bh(300), bh(200))
    borde, foco, hover_sec, activo = (nh(600), bh(600), nh(200), bh(100)) if claro else (nh(400), bh(400), nh(800), bh(900))
    o = ['<div class="panel" style="background:%s;color:%s">' % (f, txt)]
    o.append('<h2 style="margin-top:0">Tema %s <span class="mono" style="color:%s;font-weight:400">fondo %s · tarjeta %s</span></h2>' % (tema, sec, f, t))
    o.append('<div class="fila"><span>Texto normal %s:1</span><span style="color:%s">· secundario %s:1</span><span style="color:%s">· deshabilitado %s:1</span></div>' % (
        r2(contraste(txt, f)), sec, r2(contraste(sec, f)), des, r2(contraste(des, f))))
    o.append('<div class="fila"><span style="color:%s;text-decoration:underline">Enlace en reposo</span><span style="color:%s;text-decoration:underline">Enlace en hover</span></div>' % (enl, enl_h))
    if claro:
        botones = [(bh(700), "#FFFFFF", "Primario"), (bh(800), "#FFFFFF", "Primario hover"), (eh(700), "#FFFFFF", "Eliminar"), (eh(800), "#FFFFFF", "Eliminar hover"), (kh(700), "#FFFFFF", "Confirmar"), (kh(800), "#FFFFFF", "Confirmar hover")]
    else:
        botones = [(bh(300), bh(950), "Primario"), (bh(200), bh(950), "Primario hover"), (eh(300), eh(950), "Eliminar"), (eh(200), eh(950), "Eliminar hover"), (kh(300), kh(950), "Confirmar"), (kh(200), kh(950), "Confirmar hover")]
    fila = ['<span class="btn" style="background:%s;color:%s">%s</span>' % (b, c, e) for b, c, e in botones]
    fila.append('<span class="btn" style="background:transparent;color:%s;border-color:%s">Secundario</span>' % (txt, borde))
    fila.append('<span class="btn" style="background:%s;color:%s;border-color:%s">Secundario hover</span>' % (hover_sec, txt, borde))
    o.append('<div class="fila">%s</div>' % "".join(fila))
    o.append('<div class="fila"><span class="campo" style="background:%s;border:1px solid %s;color:%s">Monto, p. ej. 1,250.00</span><span class="campo" style="background:%s;border:1px solid %s;color:%s;outline:2px solid %s;outline-offset:2px">$830.50</span><span class="campo" style="background:%s;border:1px dashed %s;color:%s">Deshabilitado</span></div>' % (
        f, borde, sec, f, borde, txt, foco, nh(100) if claro else nh(900), borde, des))
    o.append('<div class="tarjeta" style="background:%s">' % t)
    o.append('<div class="fila"><span style="background:%s;color:%s;padding:4px 10px;border-radius:6px;font-weight:600"><span style="color:%s">&#9632;</span> Presupuesto (activo)</span><span style="color:%s">Gastos</span><span style="color:%s">Créditos</span></div>' % (activo, txt, enl, sec, sec))
    insignias = []
    for r in ROLES + [CRITICO]:
        b = base_rol(r["codigo"])
        insignias.append('<span class="badge" style="background:%s;color:%s">%s</span>' % (M.hex(b + ".surface-" + sfx), M.hex(b + "." + sfx), _html.escape(r["nombre"])))
    o.append('<div class="fila">%s</div>' % "".join(insignias))
    cifras = [("income", "+$12,500.00 quincena"), ("expense", "−$830.50 despensa"), ("savings", "+$1,000.00 fondo"), ("debt", "$2,431.06 Kueski")]
    o.append('<div class="fila">%s</div>' % "".join('<span class="cifra" style="color:%s">%s</span>' % (M.hex("color.role.%s.%s" % (c, sfx)), _html.escape(tx)) for c, tx in cifras))
    for grupo, orden in (("person", ORDEN_PERSONAS), ("store", ORDEN_TIENDAS)):
        o.append('<div class="fila">%s</div>' % "".join('<span class="chip"><span class="punto" style="background:%s"></span>%s</span>' % (M.cat[grupo][cod][tema].hex, _html.escape(M.cat[grupo][cod]["nombre"])) for cod in orden))
    cols = G[tema]["cols"]
    o.append(_svg_series(cols, t, sec))
    o.append('<div class="fila">%s</div>' % "".join('<span class="chip"><span class="punto" style="background:%s"></span>%d %s</span>' % (cols[n - 1], n, _html.escape(M.graf[n]["nombre"])) for n in range(1, 9)))
    o.append('</div></div>')
    return "\n".join(o)


def _bloque_simulacion(M, tema, G):
    f = M.sup(tema, "fondo")
    txt = M.neutro[900].hex if tema == "claro" else M.neutro[100].hex
    sec = M.neutro[700].hex if tema == "claro" else M.neutro[300].hex
    cols = G[tema]["cols"]
    pers = [M.cat["person"][c][tema].hex for c in ORDEN_PERSONAS]
    tien = [M.cat["store"][c][tema].hex for c in ORDEN_TIENDAS]
    o = ['<div class="panel" style="background:%s;color:%s">' % (f, txt)]
    o.append('<h3 style="margin-top:0">Tema %s: series 1–8 · personas · tiendas</h3>' % tema)
    for v in ("normal",) + VISTAS:
        celdas = ['<div class="et">%s</div>' % NOMBRE_VISTA[v]]
        celdas += ['<div class="c" style="background:%s"></div>' % hex_vista(c, v) for c in cols]
        celdas.append("<div></div>")
        celdas += ['<div class="c" style="background:%s"></div>' % hex_vista(c, v) for c in pers]
        celdas.append("<div></div>")
        celdas += ['<div class="c" style="background:%s"></div>' % hex_vista(c, v) for c in tien]
        if v == "normal":
            dato = "ΔE %.1f" % G[tema]["min_de_normal"]
        else:
            dato = "ΔL %.3f" % min(f_["dl_" + v] for f_ in G[tema]["cons"])
        celdas.append('<div class="v" style="color:%s">%s</div>' % (sec, dato))
        o.append('<div class="sim">%s</div>' % "".join(celdas))
    o.append('<p class="nota" style="color:%s;margin-top:4px">Columna derecha: mínimo entre series consecutivas (ΔL en gris y simulaciones; ΔE en visión normal).</p></div>' % sec)
    return "\n".join(o)


def generar_html(M, pares, G, C):
    o = ['<!doctype html>', '<html lang="es-MX">', '<head>', '<meta charset="utf-8">',
         '<meta name="viewport" content="width=device-width, initial-scale=1">',
         '<title>Paleta Índigo 280°</title>', '<style>%s</style>' % CSS, '</head>', '<body>', '<div class="wrap">']
    o.append('<h1>Paleta Índigo 280° · Finanzas Hogar</h1>')
    o.append('<p class="nota">Matiz base 280° · eje opuesto 100° (acción crítica) · acentos divididos 75°/135° · generado por <span class="mono">%s</span> (%s). %d tokens · %d pares de contraste, %d bajo umbral.</p>' % (
        SCRIPT_REL, FECHA, len(M.tokens), len(pares), sum(1 for p in pares if not p["cumple"])))
    o.append('<h2>Escalas de 11 pasos (L fija por paso)</h2>')
    filas = [("Marca 280°", M.marca), ("Neutro tintado", M.neutro)] + [("%s %d°" % (r["nombre"], r["matiz"]), M.esc_rol[r["codigo"]]) for r in ROLES + [CRITICO]]
    for nombre, esc in filas:
        o.append('<div class="escala"><div class="et">%s</div>%s</div>' % (_html.escape(nombre), "".join(_sw(esc[p].hex, p) for p in PASOS)))
    o.append('<h2>Componentes, roles, categóricas y gráficas en cada tema</h2>')
    o.append('<div class="temas">%s%s</div>' % (_panel_tema(M, "claro", G, C), _panel_tema(M, "oscuro", G, C)))
    o.append('<h2>Gris y simulaciones de daltonismo (Machado 2009, severidad 1.0, RGB lineal)</h2>')
    o.append('<div class="temas">%s%s</div>' % (_bloque_simulacion(M, "claro", G), _bloque_simulacion(M, "oscuro", G)))
    o.append('<h2>Tabla de contrastes (%d pares)</h2>' % len(pares))
    o.append('<div class="tabla">')
    for p in pares:
        marca = '<span class="ok">✓</span>' if p["cumple"] else '<span class="mal">✗</span>'
        o.append('<div><span class="u">%s · %s</span><span class="r">%s/%s</span>%s</div>' % (p["tema"], _html.escape(p["uso"]), r2(p["ratio"]), "%.1f" % p["umbral"], marca))
    o.append('</div>')
    o.append('</div></body></html>')
    return "\n".join(o) + "\n"


# =============================================================================
# 9. Captura PNG
# =============================================================================


def capturar_png():
    exe = shutil.which("chromium") or shutil.which("chromium-browser")
    if not exe:
        print("· chromium no encontrado: se omite paleta.png")
        return False
    url = SALIDA["html"].resolve().as_uri()
    destino = "--screenshot=%s" % SALIDA["png"]
    intentos = [[exe, "--headless=new", destino, "--window-size=1400,3000", url],
                [exe, "--headless=new", "--disable-gpu", "--hide-scrollbars", destino, "--window-size=1400,3000", url]]
    for i, cmd in enumerate(intentos, 1):
        try:
            if SALIDA["png"].exists():
                SALIDA["png"].unlink()
            subprocess.run(cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, timeout=180, check=False)
        except (OSError, subprocess.TimeoutExpired) as exc:
            print("· intento %d de chromium falló: %s" % (i, exc))
            continue
        if SALIDA["png"].exists() and SALIDA["png"].stat().st_size > 0:
            print("· paleta.png generada (intento %d)" % i)
            return True
        print("· intento %d de chromium no produjo PNG" % i)
    print("· chromium falló tras 2 intentos: se omite paleta.png")
    return False


# =============================================================================
# 10. Verificación (lee los archivos escritos)
# =============================================================================


def _aplanar(nodo, ruta, out):
    for k, v in nodo.items():
        if k.startswith("$"):
            continue
        r = "%s.%s" % (ruta, k) if ruta else k
        if isinstance(v, dict) and "$value" in v:
            out[r] = v
        elif isinstance(v, dict):
            _aplanar(v, r, out)


def verificar():
    err, info = [], []
    # --- CSV
    with open(SALIDA["csv"], newline="", encoding="utf-8") as fh:
        filas = list(csv.DictReader(fh))
    falsos = [r for r in filas if r["cumple"] != "true"]
    if falsos:
        err.append("contraste.csv: %d filas con cumple distinto de true" % len(falsos))
    for r in filas:
        real = contraste(r["primer_plano_hex"], r["fondo_hex"])
        if not (0 <= real - float(r["ratio"]) < 0.01 + 1e-9):
            err.append("contraste.csv: ratio no coincide con los hex en «%s» (%s)" % (r["uso"], r["tema"]))
        if real < float(r["umbral"]):
            err.append("contraste.csv: «%s» (%s) bajo umbral" % (r["uso"], r["tema"]))
    info.append("contraste.csv: %d filas, %d con cumple=false" % (len(filas), len(falsos)))
    # --- JSON y Anexo A del Markdown
    datos = json.loads(SALIDA["json"].read_text(encoding="utf-8"))
    tok = {}
    _aplanar(datos, "", tok)
    hx = {}
    for ruta, t in tok.items():
        if t.get("$type") != "color" or not re.fullmatch(r"#[0-9A-F]{6}", t.get("$value", "")):
            err.append("JSON: token mal formado %s" % ruta)
            continue
        ok = t.get("$extensions", {}).get("oklch", {})
        if not all(k in ok for k in ("l", "c", "h")):
            err.append("JSON: falta $extensions.oklch en %s" % ruta)
        hx[ruta] = t["$value"]
    md = SALIDA["md"].read_text(encoding="utf-8")
    anexo = dict(re.findall(r"^\| `(color\.[^`]+)` \| `(#[0-9A-F]{6})` \|", md, re.M))
    if anexo != hx:
        faltan = set(hx) ^ set(anexo)
        difieren = [k for k in set(hx) & set(anexo) if hx[k] != anexo[k]]
        err.append("JSON y colorimetria.md difieren: %d rutas no comunes, %d hex distintos" % (len(faltan), len(difieren)))
    info.append("C2-colores.json: %d tokens; Anexo A de colorimetria.md: %d filas; coinciden: %s" % (len(hx), len(anexo), "sí" if anexo == hx else "NO"))
    hex_md = set(re.findall(r"#[0-9A-F]{6}", md))
    sin_md = [k for k, v in hx.items() if v not in hex_md]
    if sin_md:
        err.append("hex del JSON ausentes en colorimetria.md: %s" % sin_md[:5])
    # --- Marca y neutro frente al plan §4.3
    rep = sum(hx["color.brand.%d" % p][1:] == PLAN_C["marca"][i] for i, p in enumerate(PASOS)) + \
        sum(hx["color.neutral.%d" % p][1:] == PLAN_C["neutro"][i] for i, p in enumerate(PASOS))
    info.append("plan §4.3 (marca y neutro) reproducido: %d/22" % rep)
    if rep != 22:
        err.append("marca/neutro no reproducen el plan §4.3 (%d/22)" % rep)
    sup = {tema: (hx[SUPERFICIE[tema]["fondo"]], hx[SUPERFICIE[tema]["tarjeta"]]) for tema in TEMAS}
    # --- Roles
    for r in ROLES + [CRITICO]:
        b = base_rol(r["codigo"])
        for tema in TEMAS:
            s = SUFIJO[tema]
            fg = hx["%s.%s" % (b, s)]
            m = min(contraste(fg, x) for x in sup[tema] + (hx["%s.surface-%s" % (b, s)],))
            if m < UMBRAL_TEXTO:
                err.append("rol %s %s bajo 4.5:1 (%.2f)" % (r["codigo"], tema, m))
    # --- Gráficas
    for tema in TEMAS:
        s = SUFIJO[tema]
        cols = [hx["color.chart.%d.%s" % (n, s)] for n in range(1, 9)]
        min_dl = {v: min(abs(l_vista(cols[i], v) - l_vista(cols[i + 1], v)) for i in range(7)) for v in VISTAS}
        min_cvd = min(min(delta_e(cols[i], cols[i + 1], "protanopia"), delta_e(cols[i], cols[i + 1], "deuteranopia")) for i in range(7))
        min_n = min(delta_e(cols[i], cols[i + 1]) for i in range(7))
        min_cr = min(contraste(c, x) for c in cols for x in sup[tema])
        croma = min(hex_a_oklch(c)[1] for c in cols)
        min_dl_gamma = min(abs(l_vista(cols[i], v, "gamma") - l_vista(cols[i + 1], v, "gamma")) for i in range(7) for v in VISTAS)
        tres_n = min(delta_e(a, b) for a, b in combinations(cols[:3], 2))
        tres_c = min(min(delta_e(a, b, "protanopia"), delta_e(a, b, "deuteranopia")) for a, b in combinations(cols[:3], 2))
        info.append("gráficas %s: ΔL consecutivo mín. gris %.3f · protan %.3f · deutan %.3f · tritan %.3f (variante culori %.3f); ΔE vecinas CVD %.1f, normal %.1f; 3 primeras entre todos sus pares ΔE normal %.1f, CVD %.1f; contraste mín. %s; croma mín. %.3f" % (
            tema, min_dl["gris"], min_dl["protanopia"], min_dl["deuteranopia"], min_dl["tritanopia"], min_dl_gamma, min_cvd, min_n, tres_n, tres_c, fmt_ratio(min_cr), croma))
        if min(min_dl.values()) < UMBRAL_DL:
            err.append("gráficas %s: ΔL consecutivo %.3f < %.2f" % (tema, min(min_dl.values()), UMBRAL_DL))
        if min_dl_gamma < UMBRAL_DL:
            err.append("gráficas %s: ΔL consecutivo con la variante culori %.3f < %.2f" % (tema, min_dl_gamma, UMBRAL_DL))
        if tres_n < NORMAL_PISO or tres_c < CVD_OBJETIVO:
            err.append("gráficas %s: las 3 primeras series no pasan la prueba de todos los pares (ΔE %.1f / %.1f)" % (tema, tres_n, tres_c))
        if min_cvd < CVD_OBJETIVO or min_n < NORMAL_PISO:
            err.append("gráficas %s: ΔE de vecinas bajo objetivo" % tema)
        if min_cr < UMBRAL_UI:
            err.append("gráficas %s: contraste bajo 3:1" % tema)
        if croma < CROMA_PISO_GRAFICAS - 0.0015:
            err.append("gráficas %s: croma bajo el piso" % tema)
    # --- Categóricas
    for grupo, orden in (("person", ORDEN_PERSONAS), ("store", ORDEN_TIENDAS)):
        for tema in TEMAS:
            s = SUFIJO[tema]
            cols = {c: hx["color.%s.%s.%s" % (grupo, c, s)] for c in orden}
            dlo = min(abs(hex_a_oklab(cols[a])[0] - hex_a_oklab(cols[b])[0]) for a, b in combinations(orden, 2))
            dlg = min(abs(l_gris(cols[a]) - l_gris(cols[b])) for a, b in combinations(orden, 2))
            cr_min = min(contraste(c, x) for c in cols.values() for x in sup[tema])
            fallas_sep = []
            for c, v in cols.items():
                L, _, H = hex_a_oklch(v)
                for r in ROLES + [CRITICO]:
                    rl, _, rh = hex_a_oklch(hx["%s.%s" % (base_rol(r["codigo"]), s)])
                    if dif_matiz(H, rh) < SEPARACION_MATIZ and abs(L - rl) < SEPARACION_L:
                        fallas_sep.append("%s-%s" % (c, r["codigo"]))
            info.append("%s %s: ΔL mín. OKLab %.3f, gris %.3f; contraste mín. %s; colisiones con roles: %d" % (grupo, tema, dlo, dlg, fmt_ratio(cr_min), len(fallas_sep)))
            if dlo < UMBRAL_DL or dlg < UMBRAL_DL:
                err.append("%s %s: escalonamiento < %.2f" % (grupo, tema, UMBRAL_DL))
            if cr_min < UMBRAL_UI:
                err.append("%s %s: contraste bajo 3:1" % (grupo, tema))
            if fallas_sep:
                err.append("%s %s: colisiones %s" % (grupo, tema, fallas_sep))
    # --- Zona reservada del eje opuesto (decisiones §I): solo color.critical usa MATIZ_OPUESTO ± ZONA_RESERVADA
    revisados, dentro, cercano = 0, [], (360.0, "")
    for ruta, t in tok.items():
        if ruta.startswith("color.critical.") or ruta not in hx:
            continue
        oklch = t.get("$extensions", {}).get("oklch", {})
        if oklch.get("c", 0.0) < CROMA_VISIBLE:
            continue
        revisados += 1
        d = dist_reserva(hx[ruta], oklch.get("h", 0.0))
        cercano = min(cercano, (d, ruta))
        if d < ZONA_RESERVADA:
            dentro.append(ruta)
    info.append("zona reservada %d°–%d° (acción crítica): %d tokens cromáticos fuera de color.critical, %d dentro; el más cercano %s a %.1f°" % (
        MATIZ_OPUESTO - ZONA_RESERVADA, MATIZ_OPUESTO + ZONA_RESERVADA, revisados, len(dentro), cercano[1], cercano[0]))
    if dentro:
        err.append("zona reservada del eje opuesto usada fuera de color.critical: %s" % dentro)
    # --- Versión 1.1 (K26 y R3-04): gasto en 40–45° y lejos de error, marca lejos de los roles, Karen frente a Bodega Aurrera
    h_gasto = next(r["matiz"] for r in ROLES if r["codigo"] == "expense")
    if not GASTO_RANGO[0] <= h_gasto <= GASTO_RANGO[-1]:
        err.append("gasto fuera de 40–45° (K26): %s°" % h_gasto)
    for tema in TEMAS:
        s = SUFIJO[tema]
        lch = {r: hex_a_oklch(hx[r]) for r in hx}
        ge = dif_matiz(lch["color.role.expense.%s" % s][2], lch["color.role.error.%s" % s][2])
        mk = lch["color.brand.%d" % PASO_MARCA[tema]]
        cerca = [r["codigo"] for r in ROLES + [CRITICO] if dif_matiz(mk[2], lch["%s.%s" % (base_rol(r["codigo"]), s)][2]) < SEPARACION_MATIZ
                 and abs(mk[0] - lch["%s.%s" % (base_rol(r["codigo"]), s)][0]) < SEPARACION_L]
        ka, kb = lch["color.person.karen.%s" % s], lch["color.store.bodega-aurrera.%s" % s]
        dm_kb, dl_kb = dif_matiz(ka[2], kb[2]), abs(ka[0] - kb[0])
        info.append("versión 1.1 %s: gasto–error %.1f°; marca %d cerca de roles: %d; Karen–Bodega Aurrera Δ matiz %.1f°, ΔL %.3f" % (
            tema, ge, PASO_MARCA[tema], len(cerca), dm_kb, dl_kb))
        if ge < SEPARACION_MATIZ:
            err.append("gasto y error a menos de %d° en %s" % (SEPARACION_MATIZ, tema))
        if cerca:
            err.append("marca cerca de roles en %s: %s" % (tema, cerca))
        if dm_kb < SEPARACION_MATIZ and dl_kb < SEPARACION_L:
            err.append("Karen y Bodega Aurrera sin separación en %s" % tema)
    # --- HTML sin recursos externos
    h = SALIDA["html"].read_text(encoding="utf-8")
    externos = re.findall(r'(?:src|href)\s*=\s*"(?:https?:)?//', h) + re.findall(r"@import|url\(", h)
    info.append("paleta.html: %d referencias externas" % len(externos))
    if externos:
        err.append("paleta.html tiene recursos externos")
    return err, info


# =============================================================================
# 11. Principal
# =============================================================================


def main():
    ap = argparse.ArgumentParser(description="Colorimetría del candidato C (Índigo 280°)")
    ap.add_argument("--sin-png", action="store_true", help="no generar paleta.png")
    ap.add_argument("--solo-verificar", action="store_true", help="solo verificar los archivos existentes")
    args = ap.parse_args()
    if not args.solo_verificar:
        M = Modelo()
        pares = construir_pares(M)
        G = metricas_graficas(M)
        Gg = metricas_graficas(M, "gamma")
        C = metricas_categoricas(M)
        V = ventana_valor_unico(M)
        I = intervalos_perfil()
        for clave in ("md", "csv", "json", "html"):
            SALIDA[clave].parent.mkdir(parents=True, exist_ok=True)
        escribir_csv(pares)
        escribir_json(M)
        SALIDA["md"].write_text(generar_md(M, pares, G, Gg, C, V, I), encoding="utf-8")
        SALIDA["html"].write_text(generar_html(M, pares, G, C), encoding="utf-8")
        print("· escritos: %s" % ", ".join(str(SALIDA[k].relative_to(RAIZ)) for k in ("md", "csv", "json", "html")))
        if not args.sin_png:
            capturar_png()
    err, info = verificar()
    for linea in info:
        print("· " + linea)
    if err:
        print("VERIFICACIÓN FALLIDA:")
        for e in err:
            print("  ✗ " + e)
        sys.exit(1)
    print("VERIFICACIÓN OK")


if __name__ == "__main__":
    main()
