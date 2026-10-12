# REVISION-especificacion-r2 · Revisión de la especificación del tablero (WF-3c, K54), ronda 2 de 3 · 12-oct-2026

Revisor independiente `revisor-especificacion-r2`. No corrijo. Evalúo el trabajo de `u3-v2` y `s1-c1-v15`, con la corrección `correccion-especificacion-r1`. La vara es K54 (`finanzas-hogar-orq/k54-tablero.md`), `docs/criterios/U3.md`, `U2.md`, `C1.md` y `S1.md`, y las decisiones §A–§K. Volví a correr todas las verificaciones con mis propios comandos, en `datos/tmp/rev-esp-r2/`; esa carpeta ya está borrada.

**Veredicto: NO APROBADO. Quedan 1 alto, 1 medio y 5 bajos nuevos.** De los hallazgos de la ronda 1, H-1 sigue abierto: `make verify-all` no pasa. H-2 y M-1 a M-4 quedan cerrados.

## 1. Resultado por punto de la tarea

| Punto | Resultado | Evidencia propia |
|---|---|---|
| (a) K54: palabras de Axel y las 4 decisiones | **CUMPLE** | Comparé con un script las dos citas de Axel de `k54-tablero.md` (658 y 293 caracteres) con la línea K54 de `decisiones.md`. Son idénticas. Las decisiones (1) a (4) están; §G se conserva («Dashboard minimalista…», línea 49) y el diff no borra nada. El «segundo problema» (lentitud) sigue sin registrarse; ver B-r1-1. |
| (b) U3 2.0 cubre lo que pidió Axel; cada bloque tiene su operación en C1 1.5 | **CUMPLE**, con M-N1 y B-N3 | Hay 4 K, 1 B, 6 G y 4 «Tarjeta», más §2.2, §2.3 y §2.4. Las 24 operaciones que cita §2 existen en los extractos (0 faltan). De los 112 nombres de campo citados, todos están en C1, salvo 5 que son valores de enum o claves de `model.json`. Los campos de §2.6 coinciden 1:1 con `DashboardDebtCapacity`, `DashboardMonthlyOverview.months[]`, `BudgetSignals.rows[].months[]`, `DashboardSettings`/`DashboardSettingsUpdate` (`reset_parameters`), `OverviewFrom`/`OverviewTo` (422) y los 14 campos nuevos. |
| (c) Ninguna cifra inventada | **CUMPLE** | Revisé los 23 montos de U3 contra `model.json` y `data.json`: 21 aparecen directos. Los otros 2 ($138.50 = 3,233.00 − 3,094.50, K-03 de S5, y $2,921.25, suma ya documentada) vienen sin cambio de U3 1.x. El 27 % lleva la marca SINTÉTICO. Los ejemplos ilustrativos de RN-30 y RN-31 están rotulados así y los recalculé: 30 %, 37 %, 117.43, 113/112/117 y 110. Todo lo no decidido es PROVISIONAL y tiene su bloqueo, y los 6 bloqueos están abiertos con qué falta, opciones e impacto. |
| (d) DESIGN §8.9 con tokens existentes y contraste propio en los dos temas | **CUMPLE**, con B-N5 | Los 25 tokens de §8.9 y U3 existen en `C2-tokens.json`, y `git diff` de C2, `tokens.css` y `tailwind.theme.css` sale vacío. Recalculé los 14 pares C-159 a C-172 con mi propia implementación: resuelvo las referencias de C2 y trunco a 2 decimales. Coinciden exactamente en claro y en oscuro. También recalculé otros 40 pares del CSV: 0 diferencias mayores de 0.01 y 0 bajo umbral. `contraste_componentes.py --check` da rc 0 y `--independiente` dice «10 razones coinciden». |
| (e) Migración 00005: goose up/down/up, `@param` completo y diccionario | **CUMPLE** | Con goose v3.28: `up` 0, `down-to 0` 0 (deja 0 objetos), `up` 0 y `down` de una sola migración 0 (quita las 5 claves y los 2 disparadores). Las 5 claves tienen su valor (30, 40, 5, 3, 10). Los disparadores rechazan `0` y `21` y aceptan `12`. Los 5 `@param` tienen sus 7 campos, más `@param_rango`. `gen_dict.py` regenera `diccionario-datos.md` idéntico byte a byte. |
| (f) `make verify-all` OK; conteos de C1; `verify_c1 --db` | **NO CUMPLE** (A-1) | `make verify-all`: **EXIT 2** en `verify-wf1-cobertura`. Con `make -k`, ese es el único objetivo que falla: oráculos, recetario, tokens, U4 y paquete pasan. C1: 22 archivos idénticos al generador, `verify_c1 --db` da TODO OK (33). Por extractos: 221 operaciones, 176 rutas y 267 esquemas, igual que `INDICE.md`. |
| (g) Criterios 261 (201/60), ids intactos, sin «≤ 6 tarjetas» vigente | **CUMPLE** | Script propio contra `git show HEAD`: 261 criterios (201 / 60 / 0) y 0 archivos con diferencias de id o de severidad. `make verify-criterios` da rc 0. En los criterios, «6 tarjetas» solo sale en citas históricas o de trazabilidad (§G, plan §3.5, listas de cambios). C1 §4 ya habla de los bloques del tablero. |

## 2. Hallazgos de la ronda 1

| Id r1 | Estado | Comprobación |
|---|---|---|
| H-1 (alta) · `verify-all` roto, 5 operaciones en 501 | **REABIERTO → A-1** | No hay `backend/internal/api/tablero/` y ningún `.go` registra las 5 operaciones. Fallan las mismas 4 pruebas en 3 paquetes (ver A-1). El reporte de `s1-c1-v15` ya lo dice bien, pero la verificación sigue rota. |
| H-2 (alta) · «≤ 6 tarjetas» vigente en C1.md y C1-datos | **CERRADO** | C1-09 nombra los bloques del tablero y M-18 a M-22. C1 §4 dice «bloques del tablero de U3 2.0 (… K54) (C1-09)». `C1-datos.md:400` y `c1_plantilla.md:294` dicen «antes 6 tarjetas». |
| M-1 (media) · U3 §2.6 contra C1 1.5 | **CERRADO** (queda el bajo B-N1) | U3 §2.6 trae 422, `reset_parameters` y «publicadas en C1 v1.5». K-4 y §1.2 leen N de `getDashboardSettings`. «Fuentes que rigen» cita C1 v1.5. |
| M-2 (media) · U1 con 6 tarjetas vigentes | **CERRADO** | Las líneas 209, 216, 457 y 475, la fila P-01 de §6, F-4 paso 1 («Toca el pago en K-4») y la línea 1057 están al día. `verify-u1` da TODO OK. |
| M-3 (media) · pares del barómetro fuera del script | **CERRADO** | C-159 a C-172 están en §10.5 y en `contraste-componentes.csv` (344 filas, 0 falsas), con la marca de umbral (3.46/6.08) y la aguja sobre la pista (9.63/7.42). |
| M-4 (media) · Verifica de C1-13 con 1.4.0 | **CERRADO** | `grep -n -m1 'version: "1.5.0"' openapi.yaml` da la línea 7 con rc 0. |
| B-1 a B-5 (bajos) | B-1, B-2, B-3 y B-4 **siguen abiertos**; B-5 no lo volví a evaluar | B-1: K54 no menciona la lentitud. B-2: §8.9 cita `U3-tablero-barometro` punto 6 y no hay `U2-*.md`. B-3: el diccionario no tiene «Cambios en 1.5» y su fecha es 11-oct. B-4: el `INDICE.md` de bloqueos no lista los 6 nuevos. |

## 3. Hallazgos

### Altos
- **A-1 · `make verify-all` da EXIT 2 (C1-02, S6-03; punto (f)). Es H-1 reabierto.**
  - Lo corrí aquí: `verify-wf1-cobertura` falla en 4 pruebas de 3 paquetes:
    - `cmd/finanzas`: `TestAreas_SinOperacionesEn501` (`areas_test.go:22`, «5 de 221 operaciones sin registrar»).
    - `internal/api`: `TestContrato_Operaciones`, con 5 subpruebas (`contrato_test.go:801`), y `TestContrato_Rutas` (`:873` y `:906`, «sin implementar 5»).
    - `internal/auth`: `TestPermisos_MismasOperacionesAmbasCuentas` (`auth_test.go:331`, Axel 501 y Ximena 501).
  - Faltan los manejadores de `getDashboardDebtCapacity`, `getDashboardMonthlyOverview`, `getBudgetSignals`, `getDashboardSettings` y `updateDashboardSettings`.
  - La corrección r1 dejó esto al agente `api-tablero`. Con eso el reporte queda al día, pero la verificación sigue rota.
  - **Corrección:** `api-tablero` registra las 5 operaciones con los motores S2 y S3, sin recalcular cifras:
    - Las 3 de Inicio y las 2 de Ajustes.
    - El PATCH valida con 422 cuando `debt.capacity_warning_pct` ≥ `debt.capacity_limit_pct`, cuando una clave llega a la vez en `parameters` y en `reset_parameters`, o cuando un valor está fuera de rango.
    - `monthly-overview` responde 422 si `from` > `to` o si pide más de 24 meses.
  - Después, `make verify-all GOOSE=…` tiene que dar EXIT 0.

### Medios
- **M-N1 · El indicador principal K-1 se apoya en una «M-01 ampliada» que C1-datos no amplía (C1-09; inconsistencia entre U3 y C1).**
  - U3 K-1 dice «Métrica C1: M-01 (N-020) ampliada» y define la fórmula de 8 campos nuevos: efectivo del mes, mes anterior (real y efectivo), las dos diferencias y los conteos de categorías en aviso y rebasadas.
  - En `C1-datos.md` §9, la fila M-01 (línea 404) solo define real, plan, % usado, restante, estado y días. Esos 8 campos solo aparecen en la lista de cambios (línea 42) y en las descripciones de `openapi.yaml`.
  - C1-09 exige que cada métrica que alimenta los bloques del tablero tenga en `C1-datos.md` su fórmula, periodo, filtros, tablas, endpoint y un ejemplo del oráculo.
  - Con K-2 y K-3 sí se hizo, con M-18 y M-19. Con K-1, que K54 nombra indicador principal («devengado y efectivo… contra el mes anterior»), no.
  - **Corrección:**
    - En `backend/internal/api/c1gen/c1_plantilla.md` §9, ampliar la fila M-01 o crear M-23 «Gasto del mes del tablero» con estos datos:
      - efectivo del hogar = `cash_household_cents` de M-11;
      - mes anterior = las mismas cifras del mes calendario anterior, ausentes sin datos;
      - diferencias = mes − mes anterior, con signo;
      - conteos = categorías con `status` `warning` o `exceeded` en M-09.
    - Ejemplo: oct-2026 con real **1,197.64** y efectivo **598.82**; nov-2026 con efectivo **1,197.64**, o sea **598.82** más que octubre (`periodos[0..1]`, R-13).
    - Regenerar con `build.py`.
    - Que U3 K-1 cite esa métrica.

### Bajos
- **B-N1 · U3 se quedó atrás en dos lugares.**
  - «Fuentes que rigen» (línea 27) sigue diciendo «métricas M-01 a M-17» cuando §2.6 ya publica M-18 a M-22.
  - §15 (línea 679) da como pendientes los cambios de U1 (P-01, §3.3, F-4) que ya hizo `correccion-especificacion-r1`.
  - **Corrección:** «M-01 a M-22» y marcar HECHO esa viñeta de §15.
- **B-N2 · Valores fijos en el nombre accesible de G-3.**
  - Dice «[categoría]: 3 meses seguidos arriba del plan por más de 10 %», pero los dos valores son configurables (`months_considered` y `margin_pct` llegan en la respuesta). §2 pide que ningún umbral sea constante.
  - **Corrección:** «[categoría]: [meses] meses seguidos arriba del plan por más de [margen] %».
- **B-N3 · K-3 no nombra M-19.** La línea «Métrica C1» de K-3 y la fila K-3 de la tabla de §2 citan M-12 y M-04, pero no M-19 «Créditos del mes». M-19 es la que define `month_payment_count` y `month_pending_household_cents`, y hoy solo aparece en §2.6. **Corrección:** sumar M-19 en las dos.
- **B-N4 · Etiquetas de versión viejas.** `Makefile` (líneas 8, 66 y 149) y `backend/db/tools/s1_verifica.py:190` dicen «esquema 1.4», pero ya corre la 00005 y el esquema es 1.5. **Corrección:** cambiar a «esquema 1.5».
- **B-N5 · La punta de la aguja puede tocar el arco de escala.**
  - DESIGN §8.9 dice que la aguja va «del centro al arco de escala».
  - Si la punta toca los tramos, ese par no está en §10.5. Mi cálculo da `ui.text.default` sobre los tramos de `ui.role.*`: de 1.99 a 2.27 en claro y de 1.67 a 1.84 en oscuro.
  - La aguja se sigue distinguiendo, porque la mayor parte de su largo queda sobre la tarjeta y la pista.
  - **Corrección:** que la aguja termine en el borde interior del arco de escala, o que lleve un anillo de 2 px de `ui.surface.card` en la punta. Si no, listar el par como decorativo con su razón.

## 4. Criterios evaluados
- **U3:**
  - U3-01: K/B/G/Tarjeta = 4/1/6/4; §2.2 a §2.4 presentes.
  - U3-02: con M-N1 y B-N3.
  - U3-03: los 25 bloques de §7 tienen vacío, carga, error y sin conexión.
  - U3-04 y U3-10: los bloques están en §2.4, §6 y §8.
  - U3-08: ver (c).
  - U3-09: 0 fechas ISO.
  - U3-11: 0 hex o `rgb(`.
  - U3-12: 5 bloqueos `U3-tablero-*`.
- **U2:** U2-04 cumple (14 pares más 40 al azar recalculados; 0 bajo umbral). U2-14 queda con B-r1-2 abierto.
- **C1:**
  - C1-02: las 5 operaciones nuevas cuelgan de N-022, N-027, N-028, N-275 y N-278, PROVISIONAL con `C1-tablero`. El backend no cumple; ver A-1.
  - C1-09: M-N1.
  - C1-13: cumple.
  - C1-21: 22 extractos, la suma da 221/176/267.
- **S1:** S1-02 (goose), S1-16 (PROVISIONAL citado) y K54 de `s1_verifica`: TODO OK.

## 5. Comandos propios (resumen)
- `make verify-all GOOSE=…`: EXIT 2. Con `make -k verify-all`, solo hay 1 «Error» (Makefile:240).
- goose directo: `up`, `down-to 0` y `up` en una base temporal, todos con rc 0. Más `down` de una migración, sondas de los disparadores y `gen_dict.py` comparado con `cmp`: idéntico.
- Scripts `python3 -I`:
  - comparación literal de K54;
  - montos de U3 contra `model.json` y `data.json`;
  - operaciones y campos de U3 §2 contra los extractos;
  - contraste propio desde `C2-tokens.json` (14 + 5 pares) y desde el CSV (40 pares);
  - ids y severidades de los criterios contra HEAD;
  - conteos de C1 desde los extractos.
- `contraste_componentes.py --check` y `--independiente`: rc 0.
- `make verify-criterios`: rc 0.
- `grep -n -m1 'version: "1.5.0"'`: línea 7.
