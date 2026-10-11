# revisor-U2a-r1 · Reporte de agente (flujo P1, puerta 1)

Fecha: 11-oct-2026 · Rol: revisor independiente de la colorimetría 1.1 (U2a), ronda 1 de 3 · Estado: TERMINADO (no cerrado por tope) · Informe de revisión: `docs/estado/p1/REVISION-U2a-r1.md`.

## Hecho

- **Lecturas, por partes:** el protocolo completo; la sección K de decisiones; la respuesta de Axel en `U2-gasto-error` y `criterios-deltal`; el bloqueo `U2a-walmart-oscuro`; R3-04 en `REVISION-r3.md`; U2-02 a U2-09 en `docs/criterios/U2.md`; el reporte `U2a-1.1.md`; y de `colorimetria.md`, «Cambios 1.1», §2.7, §7.7 y §8.1 a §8.3, más búsquedas puntuales. Del diff de `colorimetria.py` leí solo lo cambiado (-U0).
- **`make verify-color` reejecutado:** pasa.
- **Implementación propia, sin importar `colorimetria.py`:**
  - conversión OKLab de Ottosson 2020, contraste WCAG 2.x, gris = Y^(1/3) y simulación de Machado 2009;
  - generación OKLCH → sRGB con reducción de croma, que reproduce los 7 hex de referencia;
  - réplica del barrido de ventanas de Karen y Bodega Aurrera, con prueba de la condición de Walmart.
- **Verificado:**
  - (a) a (f) del encargo;
  - U2-02 a U2-09, K26, K27 y R3-04;
  - las 202 filas de `contraste.csv` recalculadas.
- **Veredicto:** NO APROBADO. Hay 1 hallazgo medio (H-01: Bodega Aurrera y Ximena casi iguales en oscuro, sin bloqueo) y 3 bajos (H-02 a H-04).

## Archivos

- `docs/estado/p1/REVISION-U2a-r1.md` (nuevo)
- `docs/estado/p1/revisor-U2a-r1.md` (nuevo, este reporte)
- **Temporales:** `datos/tmp/rev-u2a-r1/` con `rev.py`, `lib.py` y `rev2.py`, borrados al terminar.
- `make verify-color` reescribió su árbol espejo en `datos/tmp/verify-wf0/`. Esa carpeta es del Makefile y ya existía, así que no la borré.

## Verificación

- **`make verify-color`:**
  - «VERIFICACIÓN OK»;
  - md, csv, json y html idénticos byte a byte;
  - contraste.csv: 202 filas, 0 con cumple=false.
- **`python3 -I datos/tmp/rev-u2a-r1/rev.py`:**
  - **Tokens:** 190 frente a 190, sin agregados ni quitados; 20 con `$value` distinto.
  - **Gasto:** 44.73° en claro y 45.03° en oscuro. A error, 25.26° y 24.93°; a advertencia, 30.69° y 30.72°.
  - **Karen y Bodega Aurrera:** 0 fallos frente a roles, acción crítica y marca. Entre ellas, 22.95° (ΔL 0.233) en claro y 22.90° (ΔL 0.104) en oscuro.
  - **Pares cercanos:** 6 en claro y 5 en oscuro, los mismos de §7.7, y ninguno es categórica o marca frente a rol.
  - **Escalonamiento, ΔL mínimo:** tiendas 0.1049 en claro y 0.1008 en oscuro; personas 0.1033 en claro y 0.1054 en oscuro.
  - **Gráficas, ΔL consecutivo mínimo:** 0.162 en claro y 0.198 en oscuro.
- **`python3 -I datos/tmp/rev-u2a-r1/rev2.py`:**
  - **CSV:** 202/202 filas con diferencia ≤ 0.01. 16 filas cambian frente a HEAD, ninguna nueva ni quitada. Mínimos: texto 4.73, interfaz 3.05.
  - **Escala de gasto:** \|ΔL\| máximo 0.0013 y croma máximo 0.1400.
  - **Estructurales:** se reproducen todos.
  - **Ventanas:** Karen en 331°–359°; Bodega Aurrera en 331°–335°, 338°–341°, 344°–348°, 352° y 355°–359°. Con Walmart en L 0.720, ninguna.
- **Comprobación final:**
  - `$extensions` solo cambia en los tokens con valor nuevo.
  - ΔE de Walmart 1.0 → 1.1: 0.39.
  - Bodega Aurrera–Ximena en oscuro: ΔE 1.33.
  - Gasto frente a la serie 3 en claro: ΔE 7.4.

## Supuestos

- U2-02 dice «las demás entradas (U2-08)», y lo leo como un punto con dos lecturas: separación estricta, o refuerzo declarado como admite U2-08. Por eso H-01 va como media (inconsistencia sin bloqueo) y no como alta. Si el orquestador adopta la lectura estricta, H-01 sube a alta.
- La simulación de daltonismo usa las matrices de Machado 2009 con severidad 1.0, en implementación propia; no ejecuté `culori`. Reproduce §8.3 al milésimo.
- Tolerancia de ±2° para el matiz medido, como fija U2-02 por la cuantización a hex.

## Bloqueos

- Ninguno escrito: el revisor no corrige. H-01 sugiere que U2a escriba `docs/estado/bloqueos/U2a-rosas-ximena-bodega.md`.

## Pendientes

- Ninguno dentro del alcance. La ronda 2 debe verificar la corrección de H-01 (y, si se atienden, de H-02 a H-04).

## Solicitudes fuera de alcance

- **`docs/criterios/U2.md` (U2-02):** aclarar si «conservar ≥ 20° o ΔL ≥ 0.15 frente a… las demás entradas» admite el refuerzo de U2-08 entre categóricas de grupos distintos. Motivo: hoy el criterio admite dos lecturas (H-01).
- **`docs/estado/bloqueos/INDICE.md`:** registrar `U2a-walmart-oscuro` (PROVISIONAL) y, si se crea, el bloqueo de H-01. Motivo: el índice lo lleva el orquestador.
- **`docs/estado/wf0/REVISION-r3.md`:** marcar R3-04 como atendido en la colorimetría 1.1. Motivo: lo cierra el integrador.

## Paquetes instalados

- Ninguno.

## Consumo aproximado

- 20 acciones de herramienta.
- Unos 140 KB leídos, incluido el diff de 40 KB de `colorimetria.py`.
- Unos 37 KB escritos: scripts temporales (unos 17 KB) y los dos informes (unos 20 KB).
