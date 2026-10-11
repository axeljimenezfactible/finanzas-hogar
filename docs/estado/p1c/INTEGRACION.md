# Integración del flujo P1c · Aplicar respuestas de la puerta 1b · 11-oct-2026

Integrador: `integrador-p1c` (Opus). Base: commit `082235c` «[Puerta 1b] …» con las respuestas de Axel (K35 a K44 en `docs/plan/decisiones.md` §K). Este documento reúne lo que entregó cada artefacto, las verificaciones propias, las marcas PROVISIONAL que quedan, los hallazgos registrados, las propuestas para la puerta 2 y lo que pasa a WF-1.

## 1. Resultado de las revisiones

| Artefacto | Versión | Rondas | Veredicto final | Informe final |
|---|---|---|---|---|
| Modelo de datos (S1) | esquema 1.3 (migración 00003) | 3 | **APROBADO** (0 altos, 0 medios, 6 bajos) | `docs/estado/p1c/REVISION-S1-r3.md` |
| Arquitectura de información (U1) | 1.2 | 3 | **APROBADO** (0 altos, 0 medios, 10 bajos) | `docs/estado/p1c/REVISION-U1-r3.md` |
| Colorimetría (U2a) | 1.2 | 2 | **APROBADO** (0 altos, 0 medios, 3 bajos) | `docs/estado/p1c/REVISION-U2a-r2.md` |
| Contrato C1 | v1.2 (`info.version` 1.2.0) | 2 | **APROBADO** (0 altos, 0 medios; 3 bajos nuevos y 5 que siguen) | `docs/estado/p1c/REVISION-C1-r2.md` |
| Criterios de aceptación | 1.2 | 3 de 3 | **NO APROBADO** (1 alto, 2 medios, 4 bajos) | `docs/estado/p1c/REVISION-criterios-r3.md` |

Los criterios agotaron sus tres rondas. El hallazgo alto (A) no es una cifra mal: los criterios C1-22, U4-12, S7 y el `INDICE.md` §6 siguen describiendo los umbrales de alertas y los días antes del pago como si vivieran en `notification.params_json`, cuando el esquema 1.3 (corrección S1-r1, H-01) ya los siembra como claves `alerts.*` de `app_setting`, como pide K44. Con el esquema actual, las Verifica de C1-22 y U4-12 fallarían. Pasa a decisión del orquestador (§7, punto 1). Este integrador no cambió el contenido de ningún artefacto revisado.

## 2. Hecho por artefacto

**S1 · esquema 1.3** (`S1-esquema-1.3.md`, `correccion-S1-r1.md`, `correccion-S1-r2.md`)
- Migración nueva `backend/db/migrations/00003_puerta1b.sql` con Up y Down: `health_profile.goal_adjust_lose_kcal` y `goal_adjust_gain_kcal` (K36; NULL = valor de la casa; CHECK ≥ 0).
- 15 parámetros operativos como claves de `app_setting` con valor inicial sembrado y comentario que cita su decisión (K44): `auth.max_failed_attempts`, `auth.lockout_window_minutes`, `auth.password_min_length`, `session.max_age_days`, `price.stale_after_days`, `price.zero_weight_after_days`, `price.profeco_import_max_gb`, `comparator.timeout_seconds`, `attachment.max_size_mb`, `idempotency.retention_hours`, `health.goal_adjust_lose_kcal`, `health.goal_adjust_gain_kcal`, `alerts.budget_warning_pct`, `alerts.budget_limit_pct` y `alerts.payment_due_days_before`. Los tres de alertas salieron de `notification.params_json` (el Down los repone).
- Sin valor y marcados `@param_sin_valor`: las restricciones del optimizador (K16, las captura Axel) y los cuatro parámetros de avisos del bloqueo abierto `S1-parametros-avisos`.
- `reglas-negocio.md` 1.3 (29 reglas): RN-12 FIRME (K37 y K40), RN-17 FIRME (K36), RN-27 (K39), RN-29 nueva (inventario de parámetros por sección).
- `diccionario-datos.md` 1.3 regenerado con `gen_dict.py`: 55 tablas, 152 índices, 75 disparadores.
- `s1_verifica.py` con las pruebas K44, AVISOS13, PROG13, DDL y DOWN13.

**C1 · v1.2** (`C1-v1.2.md`, `correccion-C1-r1.md`)
- `parse_schema.py` lee `@param`, `@param_sin_valor` y `@col_note` y compara cada valor inicial anotado con la siembra real.
- Cada clave `@param` se expone por la operación de settings de su sección, con `value`, `initial_value`, `is_initial`, `value_type` y `unit` (K44; «restablecer» sin cambio de esquema).
- `HealthProfile*` con el ajuste por objetivo y el valor de la casa (K36).
- Marcas quitadas: `getMealPlanWeek` (K38); 4 marcas de K43 (N-023, N-064, N-226, N-297 y N-300 pasan a FIRME); `weekly_average_cents` (K39).
- N-320 a N-327 con operación; `dismissTelegramChatCandidate` es nueva (N-324); N-289 sale de `updatePerson`.
- `send_time` y `advance_days` documentados en `params_json` con patrón y mínimo, sin valor por omisión (corrección r1, A-01).
- 210 operaciones en 167 rutas, 252 esquemas, matriz de 236 filas; 17 extractos, todos de 80 KB o menos.

**U1 · 1.2** (`U1-1.2.md`, `correccion-U1-r1.md`, `correccion-U1-r2.md`)
- K38 FIRME en §0.2 y §7: los domingos, P-10, P-13 y P-14 abren la semana siguiente.
- K43 en N-023, N-064, N-226, N-297, N-298 y N-300.
- K44 en §2.3: reglas comunes (valor inicial visible, «Restablecer al valor inicial») y tabla de 15 parámetros por sección, alineada con RN-29.
- K36: P-21 por persona y P-26.13 de la casa (N-326 y N-327).
- «Descartar chat detectado» (N-324) y la fila de K35 en §7.
- 211 necesidades (138 L + 73 E), 27 pantallas y 14 secciones de Ajustes; la regla de 3 toques sigue en 3.
- `verificar_u1.py` 1.2 con 7 comprobaciones nuevas.

**U2a · colorimetría 1.2** (`U2a-1.2.md`, `correccion-U2a-r1.md`)
- K41 y K42 con la opción 1, sin cambiar ningún valor: 190 tokens con el mismo nombre, hex y OKLCH que en HEAD.
- Se quitan las 5 marcas (Karen y Bodega Aurrera, claro y oscuro; Walmart oscuro).
- Regla nueva para `DESIGN.md`: no distinguir persona y tienda solo por color en una misma vista.
- `contraste.csv`, `paleta.html` y `paleta.png` salen idénticos byte a byte.

**Criterios · 1.2** (`criterios-1.2.md`, `correccion-criterios-r1.md`, `correccion-criterios-r2.md`)
- S2-08 FIRME con K39 (219.88 y 629.35); S5-04 FIRME con K40 y K37 (0.9836, 0.4918, 0; «91 días no compite»).
- S4-07 con K36; U2-02 y U2-08 con K41 y K42.
- U4-12 y U4-13 nuevos (K44 y K36); C1-22 y C1-23 nuevos; C1-15 define «bloqueo abierto = sin Respuesta de Axel».
- S3-12 y S6-10 con K43; U1-16 con K38; P3 1.2 con K35 (la Pista 3 la ejecuta Axel).
- `tools/conteo.py` nuevo. Total: 261 criterios (201 alta y 60 media) en 18 archivos.

**Integrador**
- `Makefile`: etiquetas 1.3 y 1.2. `verify-s1` comprueba que `s1_verifica.py` aplicó cada migración del comodín. `verify-sql` muestra el orden del Down. `verify-criterios` exige además «INDICE contra conteo propio: 0 diferencias» (solicitud de `criterios-1.2.md`) y rotula el subtotal del paso 1.
- `README.md`: estado y tabla de versiones.
- `docs/estado/bloqueos/INDICE.md`: conteos y «Bloqueos nuevos de P1c».
- Este documento y el commit.

## 3. Verificaciones del integrador
- `make verify-all GOOSE=…/gobin/goose`, **dos corridas: las dos pasan** (exit 0, 10.2 s y 10.1 s) y sus salidas son idénticas línea a línea (`diff` vacío). Resumen:
  - **SQL:** Up de 00001, 00002 y 00003; 55 tablas (55 STRICT), 152 índices, 75 disparadores; `foreign_key_check` 0, `integrity_check` ok. Down en orden inverso (00003, 00002, 00001) deja 0 objetos y Up otra vez da 55 tablas. goose v3.28.0: up, down-to 0 y up sin errores.
  - **Go 1.27.2:** build, vet y gofmt sin errores.
  - **Colorimetría:** 202 pares, 0 con `cumple=false`; los 4 archivos regenerados son idénticos byte a byte.
  - **OpenAPI:** `openapi.yaml: OK`.
  - **C1:** `build.py --check` da 21 de 21 archivos idénticos; `verify_c1.py --db` da 29 comprobaciones OK y «TODO OK».
  - **S1:** `s1_verifica.py` da «RESUMEN TODO OK»: S1-01 a S1-17, KUESKI, K12, NEG (46 de 46), NEG12 (26 de 26), K44, AVISOS13, PROG13, DDL, DOWN13, DOWN y GOOSE. Aplicó las 3 migraciones del comodín.
  - **U1:** `verificar_u1.py` da 63 OK, 5 INFO, 0 PENDIENTE y «TODO OK».
  - **Criterios:** 261 criterios (201 / 60 / 0), 305 citas sin fallas, 178 cifras canon sin diferencias, 12 de 12 cifras del calendario, matriz de 236 filas con 0 N-xxx ausentes, «INDICE contra conteo propio: 0 diferencias»; las 11 comprobaciones dan lo esperado.
- **Prueba negativa del Makefile:** `make verify-s1 MIGRATIONS="…/00001_esquema.sql …/00009_falsa.sql"` falla con «FALLA: s1_verifica.py no aplicó 00009_falsa.sql» (exit 2). Después, `make clean` deja «Limpio».
- `python3 -I docs/criterios/tools/conteo.py` da 261 (201 / 60 / 0) en 18 archivos y 0 diferencias con `INDICE.md`. Los conteos del índice de criterios ya estaban reconciliados y no se tocaron.
- Conteo de bloqueos por archivo (sección «Respuesta de Axel», o «Resolución» en `criterios-nxx`): 43 archivos, 42 resueltos y 1 abierto (`S1-parametros-avisos`).

## 4. Marcas PROVISIONAL y bloqueos que citan (`grep -c PROVISIONAL`, 11-oct-2026)

| Archivo | Líneas | Bloqueos citados en esas líneas | Lectura |
|---|---|---|---|
| `backend/db/migrations/00001_esquema.sql` | 0 | — | — |
| `backend/db/migrations/00002_puerta1.sql` | 0 | — | — |
| `backend/db/migrations/00003_puerta1b.sql` | 6 | `S1-parametros-avisos` (abierto) ×5 | La 6.ª línea continúa la misma nota |
| `docs/contratos/diccionario-datos.md` | 7 | `S1-parametros-avisos` ×6 | La 7.ª es la leyenda (línea 20) |
| `docs/contratos/reglas-negocio.md` | 6 | `S1-parametros-avisos` ×5 | La 6.ª es la leyenda (línea 9) |
| `docs/design/U1-arquitectura.md` | 4 | `U1-lista-vigente-domingo` (resuelto) ×1 | Ninguna marca activa. Las 4 son historial: 3 dicen «se quitan las marcas»; la línea 39 está en «Cambios de la versión 1.1» (hallazgo bajo H-05 de U1) |
| `docs/contratos/openapi.yaml` | 3 | `S1-parametros-avisos` ×3 | — |
| `docs/contratos/C1-datos.md` | 13 | `S1-parametros-avisos` ×7; `C1-lista-compras` (resuelto) ×1 | La de `C1-lista-compras` está en la tabla de cambios (fila de v1.0 r2, historial); las otras 5 son leyenda, títulos de §13 y texto histórico |
| `docs/design/colorimetria.md` | 0 | — | Las 5 marcas de K41 y K42 se quitaron |
| `docs/contratos/C2-colores.json` | 0 | — | Ídem |
| `docs/criterios/` (19 archivos con INDICE) | 80 | `S1-parametros-avisos` ×13 (C1 3, INDICE 3, S7 4, U4 3) | Los resueltos que aparecen (`C1-valores-k31`, `S1-ajuste-objetivo-kcal`, `U1-lista-vigente-domingo`, `criterios-promedios-redondeo`, `S1-fuentes-precio`, `criterios-frescura-extremos`, `U2a-rosas-ximena-bodega`, `criterios-nxx`) salen solo como «RESUELTO, se quita la marca» o dentro de una Verifica que exige 0. Las otras 67 líneas no citan ningún archivo: son reglas de los criterios (p. ej. «valores K31 PROVISIONAL hasta su puerta» o «sin valor inventado sin marca»), no marcas sobre un artefacto |
| **Total** | **119** | | |

Conclusión: **toda marca PROVISIONAL activa cita solo el bloqueo abierto `S1-parametros-avisos`**, y ninguna cita un bloqueo con «Respuesta de Axel». `s1_verifica.py` S1-16 da lo mismo para la parte de S1: 19 líneas (migración 6, diccionario 7, reglas 6). Script: conteo por archivo con las 43 rutas de bloqueo como patrón (temporal en el scratchpad, ya borrado).

## 5. Hallazgos que quedaron registrados

**Criterios 1.2 (no aprobado; `REVISION-criterios-r3.md`)**
- **A (alta).** C1-22, U4-12, S7 (S7-04, S7-08, S7-10, §5) y `INDICE.md` §6 (líneas 68, 96 y 174) contradicen K44 y el esquema 1.3: tratan los umbrales 80 y 100 % y los días antes [1, 0] como si vivieran en `params_json`, y prohíben sembrar `alerts.*`. Las Verifica de C1-22 y U4-12 fallan con el esquema actual.
- **B (media).** P3-06: el conteo por iteración no cuadra con K35, porque `siguiente-prompt.md` vive en `iter-N/` y Claude Design puede exportar un solo PDF.
- **C (media).** La Verifica de P3-04 no tiene ruta, y el texto que busca sale también en un proceso correcto.
- **D a G (bajas).**
  - D: falsos positivos en las regex de P3-01, P3-03 y P3-12.
  - E: U4-12 atribuye a K44 la sección de `send_time` y `advance_days`.
  - F: P3-11 no está en la puerta 3 y no comprueba `salida/final/`.
  - G: S6-10 no define cuántos bytes tiene un MB.

**S1 (bajas; `REVISION-S1-r3.md`)**
- B-01: el Down no repone `params_json` byte por byte si la casa cambió valores.
- B-02: la descripción de `notification.params_json` empieza con datos que ya no aplican.
- B-03: la línea de versión de las reglas mezcla 1.3 y 1.2.
- B-04: `optimizer.weights` no está en el inventario de K44.
- B-05: el reporte S1 describe el bloqueo antes de que se redujera.
- B-06: RN-27 llama «periodo de octubre» a la ventana del oráculo.

**U1 (bajas; `REVISION-U1-r3.md`)**
- H-r3-01: la sección «Para C1» de `U1-1.2.md` conserva 2 viñetas viejas; rige la versión «corregida r2» de `correccion-U1-r2.md`.
- H-r3-02: N-300 da por hecho `params_json`.
- H-r3-03 y H-06: redacción.
- Vienen de la ronda 2:
  - H-r2-02: alcance de `send_time`.
  - H-02: falta la decisión A o F en la columna «Decisión».
  - H-03: «Decidido (K43)» para un valor que K43 no da.
  - H-04: ubicaciones propuestas que no se marcan como propuesta.
  - H-05: historial de la 1.1 con «queda PROVISIONAL».
  - H-07: `weekday` y `day_of_month` como parámetros de K44.

**U2a (bajas; `REVISION-U2a-r2.md`)**
- B-01: encabezado «medidos en el hex».
- B-02: el refuerzo de K41 está redactado distinto en §7.7.
- B-03: §12 todavía propone lo que K42 ya decidió.

**C1 (bajas; `REVISION-C1-r2.md`)**
- N-01: los encabezados generados de `openapi.yaml` y `routes.go` dicen v1.1 y v1.0.
- N-02: falta `minimum: 0` en `goal_adjust_*_kcal` de Create y Update.
- N-03: N-205 remite a un bloqueo resuelto (citar K18).
- B-01: textos de C1.md fuera del alcance de C1.
- B-02: `x-provisional` de `NotificationUpdate` cubre también `is_enabled`.
- B-03: la descripción de `Notification.params_json` se contradice.
- B-04: un mensaje de `verify_c1.py` dice que la regla del domingo es PROVISIONAL.
- B-05: falta un comando en la verificación de `C1-v1.2.md`.

**Integrador (bajas, nuevas; no cambian cifras, esquema ni endpoints)**
- I-01: la cabecera de `docs/contratos/C1-datos.md` dice «Estado: Propuesta para la puerta 1» y «Fuentes: … esquema de S1, versión 1.1»; con v1.2 debería decir puerta 1b y esquema 1.3. Sale de `backend/internal/api/c1gen/c1_plantilla.md`: se corrige ahí y se regenera, junto con N-01.
- I-02: S4-07 (`docs/criterios/S4.md`) nombra el bloqueo `S4-meta-ajustes`, que no existe. Es un bloqueo que S4 abriría si los factores intermedios de actividad no tienen fuente. No es marca activa, pero conviene decirlo así («S4 abre `docs/estado/bloqueos/S4-meta-ajustes.md` si…»).
- I-03: dos recomendaciones opuestas sobre `S1-parametros-avisos`: S1 recomienda la B y la corrección de criterios r1 la A. Dos propuestas distintas para el tope del adjunto: S1 propone de 1 a 25 MB y criterios 1.2 propone 10 MB como máximo. Las dos van juntas a Axel en §8.
- I-04: el encabezado de `docs/estado/bloqueos/INDICE.md` («32 archivos, 31 abiertos») describe WF-0. Lo dejé como historia y puse el conteo vigente al inicio de su «Resumen».

## 6. Bloqueos nuevos de P1c (1, abierto)
- `docs/estado/bloqueos/S1-parametros-avisos.md`: dónde viven `send_time`, `advance_days`, `weekday` (7) y `day_of_month` (1).
  - Opciones: A, en `notification.params_json`; B, en claves de `app_setting`; C, mixta.
  - No cambia cifras.
  - Lo dejó S1 esquema 1.3; `correccion-S1-r1` lo redujo y `correccion-S1-r2` lo amplió.
  - Registrado en `docs/estado/bloqueos/INDICE.md` («Bloqueos nuevos de P1c»).
- El integrador no abrió bloqueos.

## 7. Pendientes para WF-1 (con ruta y dueño)
1. **Criterios 1.2** (lo decide el orquestador; lo hace el agente de criterios de WF-1).
   - Hallazgo A: `docs/criterios/C1.md` (C1-22, líneas 9 y 116-117), `U4.md` (U4-12, líneas 7 y 68-69), `S7.md` (líneas 7, 35, 51, 59 y 109) e `INDICE.md` (líneas 68, 96 y 174). Los umbrales y los días antes del pago pasan a `alerts.*` de `app_setting`. La excepción PROVISIONAL queda solo para los cuatro de `S1-parametros-avisos`.
   - B y C: `P3-claude-design.md` (P3-06 y P3-04). D a G como bajos.
   - C1-13 a «C1 v1.2 / 1.2.0», C1-21 a 210 operaciones y la excepción de C1-12 para claves expuestas (`revisor-C1-r2.md`, B-01).
   - `U1.md`: criterios para K44, K36 y «Descartar chat detectado» (solicitud 2 de `U1-1.2.md`).
   - `U3.md`: sin K38 ni K44.
   - S7-04, S7-08 y S7-10: horas y anticipaciones (pendiente de `criterios-1.2.md`).
2. **C1** (próxima corrección: `backend/internal/api/c1gen/`, regenerar con `build.py`).
   - N-01: `gen4_build.py` L111 y `gen5_go.py` L15, la versión desde una constante.
   - N-02, N-03, B-02, B-03 y B-04 (`verify_c1.py`).
   - B-05: `docs/estado/p1c/C1-v1.2.md`.
   - I-01: cabecera de `c1_plantilla.md`.
3. **S1** (próxima migración, 00004).
   - B-01 a B-06.
   - En `00003_puerta1b.sql` L53-54, «lo propone S7 (WF-1)» pasa a «S7 (WF-3, K31, puerta 4)» (solicitud de `revisor-criterios-r3.md`). La 00003 no se reescribe si ya se aplicó en algún lado; si no, es un comentario.
   - Siguen de P1 §7.2 y §7.7: el tercer estado del IVA (K18) y renombrar `credit.rate_base` o `credit_purchase.rate_basis` (R2-07).
   - Si Axel elige la opción B: claves de avisos en `app_setting`, con un Down que las reponga.
4. **U1** (versión siguiente o WF-1): H-r3-01 (dejar una sola sección «Para C1» en `docs/estado/p1c/U1-1.2.md`), H-r3-02, H-r3-03, H-r2-02, H-02 a H-05 y H-07 en `docs/design/U1-arquitectura.md`.
5. **U2 (WF-1, `docs/design/DESIGN.md`).**
   - Reglas de K41: no distinguir persona y tienda solo por color en una misma vista; el refuerzo de U2-08 entre categóricas, también entre grupos distintos.
   - Opcional: los bajos B-01 a B-03 en `docs/design/tools/colorimetria.py`, y regenerar.
6. **U4 (WF-1, `docs/design/U4-ajustes.md`).**
   - Campo, validación de tipo y «restablecer» de cada parámetro de U1 §2.3 y RN-29, más N-326 y N-327.
   - De P1 §7.6: `ProjectionSettings` (K21), restricciones con «faltan restricciones» (K16), objetivo de peso (K23) y tope por persona (K13).
7. **S7 (WF-3).**
   - Los avisos de presupuesto y de pago próximo leen `alerts.*` de `app_setting`.
   - S7 propone `send_time` y `advance_days` (K31, puerta 4).
8. **Pista 3 (K35).**
   - Al cierre de WF-1, el agente paquete-claude-design (con revisor propio) deja `docs/design/entrada/`: `LEEME.md`, `prompt.md`, `DESIGN.pdf` y `.docx`, `U3-dashboard.pdf`, `U4-ajustes.pdf`, `C1-resumen.pdf`, `C2-tokens.json`, `tokens.css`, `paleta.png`, `reporte-p01..p11.png` y `MANIFIESTO.md` con sha256.
   - Crea también `docs/design/salida/` y `salida/final/` (P3-11, hallazgo F).
   - El orquestador no abre claude.ai/design.
9. **Lo que quedaba de `docs/estado/p1/INTEGRACION.md` §7.**
   - Atendidos en P1c: el punto 3 (C1 v1.2 aplicó «Para C1»: N-289 fuera de `updatePerson` y N-324 nueva), el 4 (S4-07 y S5-04 citan sus bloqueos; C1-15 y las notas 8 y 10 del índice) y el 8 (Axel respondió los 8 bloqueos).
   - El punto 1, según el mensaje del commit `3f2533b`, lo atendió la ronda 4 de C1 v1.1.
   - Siguen: el 2 (IVA), el 5 (copias de los verificadores en los anexos de `docs/estado/wf0/`), el 6 (U4 y U2) y el 7 (R3-04 en `docs/estado/wf0/REVISION-r3.md` §6, `rate_base`, cifras finales de los reportes de P1).

## 8. Propuestas para Axel (puerta 2)
**A. Bloqueo `S1-parametros-avisos`.** Dónde viven `send_time`, `advance_days`, `weekday` y `day_of_month`.
- A: `params_json`, K43 al pie de la letra.
- B: `app_setting`, coherente con K44 y con los umbrales.
- C: mixta.
- S1 recomienda la **B**; la corrección de criterios r1 recomendó la **A**, porque no cambia esquema ni contrato.

**B. Mínimos y máximos de los parámetros (K44).** El esquema y la interfaz no tienen topes hasta que se aprueben.

| Parámetro | Valor inicial | Propuesta de S1 | Otras propuestas | Justificación |
|---|---|---|---|---|
| `auth.max_failed_attempts` | 5 | 3 a 10 | ≥ 3 (criterios 1.2); ≥ 1 (U1) | Menos de 3 bloquea a la casa por un error de dedo; más de 10 debilita el freno a quien adivina contraseñas |
| `auth.lockout_window_minutes` | 15 | 5 a 60 | ≥ 1 (U1) | Con menos de 5 casi no frena; con más de 60 un error deja fuera demasiado tiempo |
| `auth.password_min_length` | 12 | 12 a 64 | — | No bajar de K6; NIST SP 800-63B rev. 4 pide 15 si la contraseña es el único factor y admitir 64 |
| `session.max_age_days` | 90 | 1 a 365 | ≤ 365 (criterios 1.2) | Un teléfono perdido no debe quedar con sesión abierta más de un año |
| `price.stale_after_days` / `price.zero_weight_after_days` | 30 / 90 | 1 a 365, con precio viejo ≤ peso cero | U1: peso cero > precio viejo | Sin esa relación, la recta de RN-12 no existe; S1 y U1 difieren en si se admite la igualdad |
| `price.profeco_import_max_gb` | 1 | 1 a 4 | — | Disco del servidor de casa |
| `comparator.timeout_seconds` | 5 | 1 a 30 | ≥ 1 (U1) | Más de 30 s deja la respuesta HTTP esperando sin aviso |
| `attachment.max_size_mb` | 10 | 1 a 25, con 1 MB = 1,048,576 bytes | ≤ 10 como tope de la interfaz (criterios 1.2) | Una foto de ticket pesa unos cuantos MB; más de 25 solo llena el disco |
| `idempotency.retention_hours` | 24 | 1 a 168 | — | La cola del teléfono no depende de esto (sus claves no se purgan) |
| `health.goal_adjust_lose_kcal` / `gain_kcal` | 0 / 0 | 0 a 1,000 | — | Más de 1,000 kcal al día no se recomienda sin supervisión médica |
| `alerts.budget_warning_pct` / `alerts.budget_limit_pct` | 80 / 100 | Aviso de 1 a 99; tope de 100 a 200 (corrección S1-r1) | U1: aviso < tope | Si no, los dos avisos llegan juntos |
| `alerts.payment_due_days_before` | [1, 0] | 1 a 3 elementos, sin repetidos, de 0 a 14 días; sin lista vacía (corrección S1-r1, C1 v1.2 y revisor C1-r1) | — | Para apagar el aviso ya está `is_enabled` |

**C. Lo demás para la puerta 2.**
- **«Restablecer».** C1 v1.2 ya expone `initial_value`, sin cambio de esquema. La alternativa es una columna `app_setting.default_value_json` en una 00004. Confirmar.
- **Sección de cada parámetro** (lectura PROPUESTA de K44). Confirmar la tabla de U1 §2.3, que coincide con RN-29, incluidas las tres que K44 no ubica: retención de idempotencia (criterios 1.2 propone Datos y respaldos), largo de contraseña y corte PROFECO. Decidir también dónde van `send_time` y `advance_days` (hallazgo E).
- **`advance_days`.** K43 dice que «nace» con valor, pero §A–§K no lo da. Confirmar que lo propone S7 (K31, puerta 4) o darlo.
- **Qué no es parámetro.** Confirmar que las constantes de fórmula y los límites de formato del contrato no son configurables (RN-29). Decidir si `optimizer.weights` entra al inventario de K44 (S1 B-04).
- **K35, PROPUESTA del orquestador.** WF-3 no espera a la Pista 3: si `docs/design/referencia-visual.md` llega antes de F1–F4 se incorpora; si llega después, se usa en U5 y en las correcciones.

## 9. Paquetes instalados
Ningún reporte de P1c declara paquetes instalados. El integrador tampoco instaló nada.

## 10. Consumo por agente (de sus reportes; cifras aproximadas)

| Agente | Artefacto | Acciones | KB leídos | KB escritos |
|---|---|---|---|---|
| S1-esquema-1.3 | S1 | 38 | 205 | 60 |
| correccion-S1-r1 | S1 | 37 | 175 | 48 |
| correccion-S1-r2 | S1 | 27 | 95 | 15 |
| revisor-S1-r1 / r2 / r3 | S1 | 28 / 28 / 25 | 125 / 205 / 140 | 25 / 17 / 27 |
| U1-1.2 | U1 | 28 | 170 | 60 |
| correccion-U1-r1 / r2 | U1 | 23 / 34 | 100 / 95 | 32 / 38 |
| revisor-U1-r1 / r2 / r3 | U1 | 27 / 33 / 36 | 155 / 175 / 190 | 22 / 22 / 25 |
| U2a-1.2 | U2a | 40 | 150 | 30 |
| correccion-U2a-r1 | U2a | 19 | 75 | 12 |
| revisor-U2a-r1 / r2 | U2a | 33 / 34 | 140 / 140 | 30 / 26 |
| criterios-1.2 | criterios | 58 | 190 | 75 |
| correccion-criterios-r1 / r2 | criterios | 22 / 14 | 130 / 75 | 18 / 12 |
| revisor-criterios-r1 / r2 / r3 | criterios | 47 / 39 / 33 | 175 / 160 / 150 | 19 / 17 / 18 |
| C1-v1.2 | C1 | 46 | 180 | 55 |
| correccion-C1-r1 | C1 | 29 | 115 | 14 |
| revisor-C1-r1 / r2 | C1 | 35 / 39 | 205 / 185 | 17 / 26 |
| integrador-p1c | integración | ~52 | ~175 | ~40 |
| **Total (27 agentes)** | | **~904** | **~4,075** | **~800** |

Por artefacto: S1 183 acciones (945 KB leídos, 192 KB escritos); U1 181 (885, 199); U2a 126 (505, 98); criterios 213 (880, 159); C1 149 (685, 112). Ningún agente pasó los topes del protocolo §4 (máximo leído 205 KB de 400; máximo escrito 75 KB de 150; máximo 58 acciones de 60, `criterios-1.2`). No cuentan los archivos que reescriben los generadores (`build.py`, `gen_dict.py`, `colorimetria.py`).
