# revisor-especificacion-r3 · Reporte · 12-oct-2026

Agente `revisor-especificacion-r3` (WF-3c, ronda 3 de 3). Revisor independiente: no corrijo nada. El informe completo está en `docs/estado/wf3/REVISION-especificacion-r3.md`.

No hice `git add` ni `git commit`, no mandé Telegram, no usé sudo y no instalé paquetes. No toqué nada fuera de estos dos reportes.

## Hecho
- **Veredicto: NO APROBADO.** Quedan 0 altos, 1 medio (M-R3-1) y bajos (B-R3-1, B-R3-2 y los que siguen de r1 y r2).
- **Cerrados los dos hallazgos de r2:**
  - A-1 / H-1: las 5 operaciones están registradas y `make verify-all` da EXIT 0.
  - M-N1: la fila M-01 de C1-datos §9 ya trae los 8 campos con su ejemplo del oráculo.
- **Nuevo M-R3-1:** U3 G-6 usa `scale_max_pct`, pero ese campo no está en `DashboardMonthlyOverview`, la única operación de G-6.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/REVISION-especificacion-r3.md` (informe)
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/revisor-especificacion-r3.md` (este reporte)

## Verificación
- **Punto (a) · K54:** un script compara las citas «…» de `k54-tablero.md` con la línea K54 de `decisiones.md`.
  - Resultado: 658 y 293 caracteres, idénticas. La de la lentitud no está: es el bajo B-r1-1.
  - `git diff --numstat docs/plan/decisiones.md`: 1 línea agregada y 0 borradas.
- **Punto (b) · U3-01 · U3-02 · C1-09:** un script cruza los identificadores de cada bloque K, B y G de U3 §2 con los extractos de sus operaciones. 0 faltantes reales.
  - Revisando esquema por esquema: `DashboardMonthlyOverview` (`inicio.yaml`, de 1499 a 1573) solo trae `warning_pct`, `limit_pct` y `months`, y G-6 (U3:229) usa `scale_max_pct`. Eso es M-R3-1.
- **Punto (c) · U3-08:**
  - Los montos de U3 que no aparecen en `model.json` ni en `data.json` son números de sección (P-26.xx, §3.11, §3.13) o vienen de U3 1.x (138.50, 2,921.25) y ya se aceptaron en r2.
  - Los 3 montos nuevos frente a HEAD están en `periodos[2]`.
  - Recalculé RN-30 (30, 37, 117.43, 56) y RN-31 (113, 112, 117, 110).
- **Punto (d) · U2:** contraste propio de C-159 a C-172 desde `C2-tokens.json`, truncado a 2 decimales.
  - Los 28 valores coinciden con DESIGN §8.9 y con el CSV, y ninguno queda bajo su mínimo.
  - Los 51 tokens citados existen y hay 0 colores sueltos en las líneas nuevas.
  - `C2-tokens.json` y `tokens.css` sin cambios.
- **Punto (e) · S1:** goose con una base temporal:
  - `up`, `down`, `up`, `down-to 0` y `up`: todo con rc 0, versión 5.
  - El disparador rechaza 0 y 21 y acepta 12.
  - `gen_dict.py` a un archivo temporal es idéntico a `diccionario-datos.md` según `cmp`.
- **Punto (f) · C1-02 · S6-03:**
  - `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: **EXIT 0**. C1 da 22 archivos idénticos y `verify_c1 --db` da «comprobaciones OK: 33, TODO OK».
  - Mis conteos: 221 operaciones, 176 rutas y 267 esquemas; `version: "1.5.0"` en la línea 7.
  - `go test -count=1 -v` de las 4 pruebas de A-1: PASS.
  - `go test -cover`: `api/tablero` 71.2 % y `api/creditos` 80.4 %.
  - `make -k verify-wf3-contrato verify-wf3-secretos`: rc 0, «221 operaciones en 176 rutas, 0 diferencias».
- **Punto (g):** script contra `git show HEAD`: 261 criterios (201 / 60 / 0) y 0 diferencias de id o de severidad. `make verify-criterios` da rc 0. «6 tarjetas» solo aparece como cita histórica.

## Supuestos
1. Califiqué M-R3-1 como medio porque es una «inconsistencia entre artefactos» (U3 G-6 contra C1 M-21) y no cambia ninguna cifra.
2. Califiqué B-R3-2 (los 14 campos opcionales sin llenar) como bajo porque toca a la implementación del backend (S6), no a U3, C1 ni S1. Aun así, lo debe ver el integrador antes del frontend P-01.
3. En `datos/tmp/` hay carpetas de otros agentes (`rev-rem-r1/`, `corr-s5-r*`, `playwright-*` y otras) y la que usan los objetivos de `make` (`verify-wf3/`). No son mías y no las toqué.

## Bloqueos
No abrí ninguno. Siguen abiertos `C1-tablero` y los cinco `U3-tablero-*`.

## Pendientes
- **M-R3-1:** se corrige en U3 G-6, cambiando a «al mayor entre `limit_pct` y el mayor `debt_ratio_pct`». La otra opción es agregar `scale_max_pct` a `DashboardMonthlyOverview` (c1gen, `build.py`, `tablero.go`, U3 §2.6, M-21).
- B-R3-1 y B-R3-2.
- Los bajos de r1 y r2 que siguen abiertos (B-N1 a B-N5, B-r1-1, B-r1-3 y B-r1-4).

## Propuestas para Axel
1. **PROPUESTA:** que la escala de G-6 sea la misma que la del barómetro B-1 (0 a `scale_max_pct`, hoy 60 %), para que el tablero lea igual en el medidor y en la tendencia. Eso pide la opción B de M-R3-1.

## Paquetes instalados
Ninguno.

## Consumo aproximado
46 acciones. Unos 150 KB leídos, por secciones: U3 §2 y §2.6, DESIGN por diff, extractos con `awk` y `grep`, y nada de `openapi.yaml` completo en contexto. Unos 17 KB escritos (estos dos reportes). Borré `datos/tmp/rev-esp-r3/`.
