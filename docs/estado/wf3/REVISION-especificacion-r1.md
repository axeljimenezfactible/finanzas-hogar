# REVISIÓN especificación r1 · WF-3c (K54, tablero de control) · 12-oct-2026

Revisor independiente `revisor-especificacion-r1` (ronda 1 de 3). No corrige. Evalúa el trabajo de `u3-v2` (U3 2.0, DESIGN 1.1, K54, criterios) y de `s1-c1-v15` (esquema 1.5, C1 v1.5, criterios C1 y S1) **solo** contra K54 (`finanzas-hogar-orq/k54-tablero.md`), `docs/criterios/U3.md`, `U2.md`, `C1.md`, `S1.md` y decisiones §A a §K. Cada verificación se volvió a correr con comandos y scripts propios en `datos/tmp/rev-esp-r1/` (borrados al terminar).

**Veredicto: NO APROBADO.** 2 altos, 4 medios y 5 bajos.

## 1. Resultado por punto de la tarea

| Punto | Resultado | Evidencia propia |
|---|---|---|
| (a) K54: palabras de Axel y las 4 decisiones | CUMPLE | Script `k54.py`: las dos citas largas de Axel (658 y 293 caracteres) están **idénticas** en `decisiones.md` K54 y en `k54-tablero.md`. Las decisiones 1 a 4 coinciden en contenido (opción A de cada una). La única cita que falta es la del «segundo problema» (lentitud), que no es parte de las 4 decisiones (B-1). |
| (b) U3 2.0 cubre lo pedido y cada bloque tiene operación en C1 1.5 | CUMPLE con M-1 | Ingreso con extra (K-2), créditos globales (K-3 y G-2), recordatorio N (K-4), barómetro con dígito y umbrales 30/40 (B-1), 6 gráficas (G-1 a G-6). Script `u3_02.py`: 22 métricas citadas (M-01 a M-22), todas en `C1-datos.md`; 34 `operationId` citados, todos en `openapi.yaml`; 41 campos de los bloques nuevos, todos en `openapi.yaml`. La desalineación de U3 §2.6 con C1 1.5 publicado va en M-1. |
| (c) Ninguna cifra inventada | CUMPLE | Script `montos.py` sobre el diff de U3: 18 montos. 16 están en `oraculo/model.json`. `$138.50` y `$3,233.00` vienen de S5 K-03 / `data.json` y ya estaban en 1.5. Los nuevos ($1,482.57, $1,767.51, $2,232.43, $1,947.50, $6,293.44, $5,034.74, $1,258.70) coinciden con `periodos[0..2]` y R-53 a R-55. Equivalencia de mes calendario con P1 a P3 comprobada: los pagos caen el 15, el 30 o el 28 y ningún evento cae entre los días 1 y 9. «27 %» va marcado SINTÉTICO. Los valores 25, 30, 31, 40 y 41 de A-18 son respuestas simuladas de prueba. En RN-30 y RN-31, los ejemplos con 4,000.00 y con 3,000.00 van marcados «ilustrativo» y su aritmética es correcta: 29.941 → 30, 37.064 → 37, margen 117.43, 55.811 → 56; 113, 112 y 117; 109.67 → 110. Lo no decidido lleva PROVISIONAL con su bloqueo (6 bloqueos nuevos). |
| (d) DESIGN §8.9: tokens y contraste | CUMPLE con M-3 y B-2 | 0 hex o `rgb()` en §8.9. Los 38 tokens citados de §8.2 a §8.9 existen en `tokens.css`, y los 55 de U3 también. Script `contraste.py` propio (resuelve `var()` de `tokens.css`, claro y oscuro): los 11 pares de la tabla coinciden **exactamente** (diferencia ±0.00). Pares que la tabla no lista, calculados aquí: `ui.viz.baseline` sobre tarjeta 3.47 / 6.09 (mínimo 3, cumple); aguja `ui.text.default` sobre `ui.viz.track` 9.64 / 7.42. La pista sobre la tarjeta da 1.12 / 1.45 (decorativa, igual que C-105). Las superficies de rol sobre la tarjeta dan de 1.00 a 1.04 (por eso no se usan). 0 pares bajo umbral. |
| (e) 00005: goose, `@param`, diccionario | CUMPLE con B-3 | goose v3.28.0 propio: `up` (versión 5) → `down` (versión 4) → esquema idéntico al de `up-to 4` (comparación de `sqlite_master`), 0 claves y 0 disparadores k54 → `down-to 0` → `up` (versión 5); integrity_check `ok`, foreign_key_check vacío; las 5 claves con 30, 40, 5, 3 y 10. Rangos: 12 de 12 pruebas propias (0, 101, `"40"`, 40.5, 21, 1, 51 y `null` rechazadas; 100, 20, 12 y 50 aceptadas). Los 5 `@param` tienen sus 7 campos y hay un `@param_rango`. `gen_dict.py` regenerado en temporal: **idéntico** a `diccionario-datos.md`. |
| (f) `make verify-all`, conteos de C1 y `verify_c1 --db` | **NO CUMPLE** (H-1) | `make verify-all GOOSE=…/goose` propio: **EXIT 2** en `verify-wf1-cobertura`. Antes pasan WF-0, C1 (22 archivos idénticos, `verify_c1.py --db` «TODO OK», 33 comprobaciones), S1 (TODO OK, K54 11/11 negativas y 5/5 positivas), U1 y criterios. Las etapas siguientes, corridas aparte (oráculos, recetario, tokens, U4 y paquete), dan EXIT 0. Conteos propios: 221 `operationId` únicos (216 en HEAD), 176 rutas, 267 esquemas, 221 entradas en `routes.go` (216 en HEAD), versión 1.5.0. Coinciden con el reporte. |
| (g) Criterios: 261 (201/60), ids, verify-criterios, sin «≤ 6 tarjetas» vigente | **NO CUMPLE** (H-2) | Script `ids.py`, HEAD contra la copia de trabajo: 261 = 201 / 60 en los dos y 0 ids o severidades cambiados. `verify-criterios` OK dentro de verify-all (261 definidos y 261 citados, 326 citas sin fuente faltante, 12 cifras canon). Pero `C1.md` §4 conserva «≤ 6 tarjetas de dashboard (C1-09)» como umbral vigente. |

## 2. Criterios evaluados

| Criterio | Veredicto | Nota |
|---|---|---|
| U3-01 | CUMPLE | `grep -c`: K-[1-4] = 4, B-1 = 1, G-[1-6] = 6, Tarjeta = 4; §2.2, §2.3 y §2.4 presentes; umbrales leídos de la respuesta; tarjeta 4 con K38 (domingo y entre semana); tarjeta 5 «$138.50 menos que Walmart». |
| U3-02 | CUMPLE (ver M-1) | 0 métricas ni operaciones faltantes en C1 1.5. |
| U3-03 | CUMPLE | §7: 25 filas × vacío, carga, error y sin conexión, sin celdas vacías. |
| U3-04 | CUMPLE | §2.4 (apilado a 360 y 412, rejilla a 768, 1024 y 1440, orden del foco) y §8. |
| U3-05 | CUMPLE | G-4 = T-2 y G-5 = T-3; 0 donas; 2,921.25, 3,895.00, 13,103.78 y 14,362.48 sin cambio. |
| U3-06 | CUMPLE | 37 líneas con «Karen»; las nuevas son informativas (K-3, K-4 y B-1 desde 768). |
| U3-07 | CUMPLE | K-3: A $10,000.00 y B $6,905.50 con «Pendiente confirmar». K-4: 30-oct por $973.75. |
| U3-08 | CUMPLE | Ver (c). |
| U3-09 | CUMPLE | 0 fechas ISO en U3; textos nuevos en es-MX; las palabras en inglés que aparecen son marcadores (`[warning]`) o citas literales de Axel. |
| U3-10 | CUMPLE | §6 con filas para K-1 a K-4, B-1, G-1 a G-6 y T2 a T6. |
| U3-11 | CUMPLE | 0 hex o `rgb()`; tokens de rol existentes; forma, patrón, etiqueta y tabla en las gráficas. |
| U3-12 | CUMPLE (ver B-4) | §11, §14 y §15; 5 `U3-tablero-*` más `U3-t2-acumulados`. |
| U2-04 | CUMPLE en valores, falta cobertura (M-3) | 0 pares bajo umbral; los pares del barómetro no están en la tabla de §10.5 que calcula el script. |
| U2-12, U2-13 | CUMPLE | §8.9 con anatomía, reglas, tamaños y accesibilidad; §8.2 con dos filas nuevas. |
| U2-14 | Bajo (B-2) | La numeración PROVISIONAL de DESIGN quedó en un bloqueo `U3-*`, no en uno `U2-*`. |
| C1-01, C1-02, C1-21, C1-22 | CUMPLE | Validador OK, 5 operaciones con N-xxx PROVISIONAL (bloqueo C1-tablero), 22 extractos de 80 000 bytes o menos (el mayor, `ajustes-c.yaml`, tiene 78 329), K44 con 30 claves y cada una con su lectura y su cambio. |
| C1-09 | **NO CUMPLE** (H-2) | Sigue con «6 tarjetas». |
| C1-13 | Medio (M-4) | El comando de su Verifica busca 1.4.0. |
| S1-02, S1-03, S1-13, S1-16 | CUMPLE | Ver (e) y la salida de `s1_verifica.py` («RESUMEN TODO OK»). |

## 3. Hallazgos

### Altos

**H-1 (alta) · `make verify-all` está roto (EXIT 2).**
- **Qué falla.** C1 v1.5 publica 5 operaciones nuevas (`getDashboardDebtCapacity`, `getDashboardMonthlyOverview`, `getBudgetSignals`, `getDashboardSettings` y `updateDashboardSettings`) que ningún área registra. Responden 501, y `verify-wf1-cobertura` falla en **4 pruebas de 3 paquetes**:
  - `cmd/finanzas`: `TestAreas_SinOperacionesEn501`.
  - `internal/api`: `TestContrato_Operaciones`, con 5 subpruebas, y `TestContrato_Rutas`.
  - `internal/auth`: `TestPermisos_MismasOperacionesAmbasCuentas` (Axel 501 y Ximena 501 en las 5).
- **El reporte de `s1-c1-v15` se queda corto.** Solo cita `TestAreas_SinOperacionesEn501` y la da por «esperada». La regla de la tarea es que verify-all no se rompa.
- **Corrección.** Antes de integrar WF-3c, el agente backend `api-tablero` registra los 5 manejadores (con los motores S2 y S3, sin recalcular cifras) en `gastos-inicio` y `ajustes`, y verify-all vuelve a dar EXIT 0. Si no da tiempo, el integrador no cierra WF-3c con C1 v1.5 sin esos manejadores. Además, el reporte `docs/estado/wf3/s1-c1-v15.md` §Verificación debe listar las 4 pruebas de los 3 paquetes.

**H-2 (alta) · `C1.md` y `C1-datos.md` conservan «≤ 6 tarjetas» como regla vigente (contradice K54, decisión 3).**
- **Dónde queda la regla.** `s1-c1-v15` subió `C1.md` a 1.5 sin tocar estos puntos:
  - `docs/criterios/C1.md` §4, línea 163: «≤ 6 tarjetas de dashboard (C1-09)».
  - C1-09, línea 88: «Cada métrica que alimenta las 6 tarjetas de U3 … Cubre: mes vs plan, flujo quincenal, … ahorro», sin M-18 a M-22.
  - `docs/contratos/C1-datos.md` §9, línea 400, salida del generador (`c1_plantilla.md`): «Las 6 tarjetas de Inicio › Resumen (U1 N-020 a N-025, U3)».
- **Corrección.**
  - `C1.md` §4: cambiar por «bloques del tablero de U3 2.0 (K-1 a K-4, B-1, G-1 a G-6 y Más del hogar; K54) (C1-09)».
  - C1-09: cambiar por «Cada métrica que alimenta los bloques del tablero de U3 2.0 (K54) y sus subsecciones …» y sumar «M-18 a M-22 del tablero» a «Cubre».
  - `backend/internal/api/c1gen/c1_plantilla.md`, introducción de §9: cambiar por «Los bloques del tablero de Inicio › Resumen (U3 2.0, K54; antes 6 tarjetas)». Después, regenerar con `build.py`.
  - Mantener los mismos ids y severidades (261 = 201 / 60).

### Medios

**M-1 (media) · U3 2.0 §2.6 (y B-1 y K-4) no está alineado con C1 v1.5 tal como se publicó.**
- **Las cuatro diferencias:**
  - **Código de error de la ventana:** U3 §2.6 dice que una ventana inválida de `getDashboardMonthlyOverview` responde «400 `validation_failed`»; C1 responde 422 (C1 §3 y supuesto 4 de `s1-c1-v15`).
  - **Campo para restablecer:** U3 dice que el PATCH lleva «`parameters` y `reset`»; C1 usa `reset_parameters`.
  - **Estado de lo nuevo:** U3 sigue marcando las operaciones, los campos y M-18 a M-22 como PROPUESTA «mientras C1 no los publique», y K-4 dice «mientras C1 no la publique, la interfaz usa 5». C1 1.5 ya los publica.
  - **De dónde sale N en K-4:** K-4 no dice que N se lee de `getDashboardSettings`, así que el 5 queda como constante en el código. Eso choca con la regla de §2 «Ningún umbral es constante en el código».
- **Detalle menor:** «Fuentes que rigen» todavía cita «contrato C1 v1.2».
- **Corrección** (agente de U3):
  - §2.6: 422 y `reset_parameters`; las operaciones y M-18 a M-22 pasan a «publicadas en C1 v1.5».
  - K-4: N = `parameters["dashboard.payment_reminder_count"].value` de `getDashboardSettings`, leído una vez, como `getAlertSettings` en §1.2.
  - «Fuentes que rigen»: C1 v1.5.

**M-2 (media) · U1 sigue con el Resumen de 6 tarjetas como regla vigente: los artefactos no son coherentes con K54 ni con U3 2.0.**
- **Dónde está la regla:** `docs/design/U1-arquitectura.md`:
  - Línea 209: «Resumen reúne como máximo 6 tarjetas».
  - Línea 216: P-01, «6 tarjetas en orden fijo».
  - Línea 457: «Resumen (6 tarjetas, orden fijo)».
  - Línea 1169: P-01, «6 tarjetas … 3 × 2».
  - Línea 1215.
  - F-4 paso 1: «Toca la tarjeta «Próximo pago»».
- **Quién lo dejó así.** `u3-v2` lo dejó bien registrado como solicitud fuera de su alcance (U3 §15), pero sigue abierto.
- **Corrección.** El orquestador asigna la solicitud al dueño de U1:
  - §2.1: «Resumen = tablero de control de U3 2.0 (K54)».
  - Filas de P-01 en §2.1, §2.2 y §6.
  - §3.3: necesidades propias del tablero (así se cierra el punto 3 de C1-tablero).
  - F-4 paso 1: «Toca el pago en K-4».

**M-3 (media) · Los pares de contraste del barómetro quedan fuera de la tabla que calcula el script (U2-04).**
- **Qué pasa.** DESIGN §8.9 trae 11 pares calculados con un script ad hoc de `u3-v2`, ya borrado, y deja su numeración para «cuando U2 regenere la tabla». La tabla §10.5 y `contraste_componentes --check` no los cubren. U2-04 pide que la tabla calculada por script cubra **todos** los pares usados.
- **Lo que falta además.** La tabla de §8.9 omite dos pares que el componente usa:
  - Las marcas de umbral, `ui.viz.baseline` sobre `ui.surface.card`: 3.47 en claro y 6.09 en oscuro.
  - La aguja, `ui.text.default` sobre `ui.viz.track`: 9.64 y 7.42.
- **Valores.** Todos los valores de §8.9 coinciden con el recálculo propio y ninguno queda bajo umbral.
- **Corrección.** Agregar a §10.5 las filas C-xxx del barómetro: tramos de rol sobre la tarjeta, valor sobre la pista, insignias, dígito, marcas y aguja, en claro y oscuro. Regenerar y comprobar con `contraste_componentes --check` (dentro de `make verify-wf1-tokens`).

**M-4 (media) · El comando de C1-13 busca la versión anterior.**
- **Qué falla.** `C1.md` C1-13 ya dice «C1 v1.5 (`info.version` 1.5.0)», pero su Verifica sigue con `grep -n -m1 "version: \"1.4.0\"" …/openapi.yaml`. Corrido aquí: sin coincidencias (rc = 1); el contrato dice `version: "1.5.0"` en la línea 7. Leído al pie de la letra, el criterio falla, aunque `verify_c1.py` sí exige 1.5.0.
- **Corrección.** Cambiar a `grep -n -m1 "version: \"1.5.0\""`.

### Bajos

- **B-1 · Falta el «segundo problema» en K54.** `decisiones.md` K54 no registra lo que Axel reportó a la misma hora: «entro a las secciones y cargan bastante lento…», con el diagnóstico de mDNS y HTTP/2 que `k54-tablero.md` sí recoge. No es una de las 4 decisiones. Propuesta: una línea en K54 o una decisión aparte, para que la corrección de `tls.go` tenga fuente.
- **B-2 · Numeración del barómetro.** DESIGN mete el barómetro como §8.9 y `k54-tablero.md` dice «§8.5». La diferencia está justificada y marcada PROVISIONAL, pero en el bloqueo `U3-tablero-barometro` (punto 6), no en `docs/estado/bloqueos/U2-<tema>.md`, como pide U2-14. Propuesta: abrir `U2-barometro-numeracion.md`, o citarlo también desde allí.
- **B-3 · Encabezado del diccionario.** `diccionario-datos.md` 1.5 conserva la fecha «11-oct-2026» en el encabezado, no tiene la línea «Cambios en 1.5» (00005 y K54) y el Autor no nombra a `s1-c1-v15`. Corrección: sumar la línea 1.5 en `gen_dict.py`, que genera el encabezado, y regenerar.
- **B-4 · Índice de bloqueos.** `docs/estado/bloqueos/INDICE.md` no lista los 6 bloqueos nuevos (`U3-tablero-*` y `C1-tablero`). Los dos reportes lo señalan como pendiente fuera de su alcance. Le toca al integrador.
- **B-5 · Vacío de la tarjeta 5.** En U3 §2.5, la condición de vacío de la tarjeta 5 («sin `store_id`: ningún artículo con precio, o sin lista…») no nombra el caso «ninguna tienda con todos los precios» (S5 §7.4, Q2-60 = A). De ese caso depende la nueva «Condición de la ruta A» de U1 F-3. Propuesta: sumarlo a la lista de causas.

## 4. Lo que está bien (HECHO comprobado)
- **K54.** Las citas son literales y las 4 decisiones son fieles. «Más del hogar» y los demás puntos no decididos van PROVISIONAL, cada uno con su bloqueo, y cada bloqueo trae opciones e impacto.
- **U3 2.0.**
  - Responde al propósito de Axel: el gasto real del mes es el indicador principal, el desglose en cadena llega hasta P-15 › Tendencias en 4 toques y la tabla «Qué miro para prevenir» está completa.
  - La interfaz no calcula nada.
  - El estado del barómetro llega de C1 con umbrales que no son constantes.
- **00005.**
  - Solo agrega 5 claves y 2 disparadores, sin cambiar DDL anterior; el Down deja el 1.4 idéntico.
  - Las verificaciones K54 del script S1 son nuevas y sus pruebas negativas son buenas.
  - El cambio de PROG14 en `s1_verifica.py` es acotado: no debilita la comprobación de la 00004.
- **C1 v1.5.** Lo generó `c1gen` y es idéntico byte a byte. Los 14 campos nuevos son opcionales y no rompen P-06, P-08 ni P-17. `x-app-setting` está en los 4 campos de umbral.

## 5. Comandos propios (resumen)
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP=…/datos/tmp/rev-esp-r1/vwf0`: EXIT 2. Falla `verify-wf1-cobertura`: 3 paquetes en FAIL, el resto en ok.
- `make verify-wf1-oraculos verify-wf1-recetario verify-wf1-tokens verify-wf1-u4 verify-wf1-paquete`: EXIT 0.
- goose up / down / `up-to 4` / comparación de `sqlite_master` / `down-to 0` / up: todo rc 0; esquema tras el down = 1.4.
- Scripts propios en `python3 -I`: `k54.py`, `contraste.py`, `montos.py`, `u3_02.py` e `ids.py`; también la regeneración de `gen_dict.py` con `cmp` y los conteos de openapi y `routes.go`.
