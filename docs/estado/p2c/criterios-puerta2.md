# P2c · criterios-puerta2 (Sonnet) · 11-oct-2026, noche

Alcance: `docs/criterios/{S1,S2,S4,S5,S7,C1,U3,U4}.md`, `docs/criterios/INDICE.md` y este reporte. No toqué P3 (ya está en 1.3), U2, S3, S6, S8 ni WF4 (sus marcas K31 son de la puerta 4). No hice `git add` ni `commit`, no usé navegador, WebFetch ni Telegram. Manda `docs/estado/puerta2/RESPUESTAS.md` (modo «lo recomendado + excepciones», sin excepciones; Q2-11 = B, Q2-64 = A, Q2-96 = A, y la puerta 2 bis). Ids y severidades intactos: **261 criterios (201 alta / 60 media), 18 archivos**.

## Hecho
Versiones nuevas: S1 1.2, S2 1.3, S4 1.3, S5 1.3, S7 1.3, C1 1.3, U3 1.3, U4 1.3; INDICE 1.3 (§6, §10 y §12 nuevo). Cada archivo trae su lista «Cambios de la versión …».

- **Valores K31/K44 aprobados pasan de PROVISIONAL a FIRME con clave y valor (Q2-72 a Q2-77, K50).**
  - Reportes: `reports.daily_send_time` = `"21:30"`, `reports.weekly_weekday` = 7, `reports.weekly_send_time` = `"19:00"`, `reports.monthly_day_of_month` = 1, `reports.monthly_send_time` = `"08:00"`.
  - Avisos: `alerts.send_time` = `"09:00"`, `alerts.stale_prices_advance_days` = 1, `alerts.meal_plan_missing_weekday` = 7.
  - Otros: `optimizer.date_window_days` = 7, `ocr.timeout_seconds` = 30; comparador de 1 a 30 s (sigue en 5); `price.stale_after_days` ≥ 1 y `price.zero_weight_after_days` ≥ ese valor (lo validan el motor y la API, no el esquema).
  - Dónde quedó: S1-16 (lista de qué ya no lleva marca), C1-22 (tabla de las diez claves con rangos y operaciones; verificación con `verify_c1.py --db`), C1-23, U4-12 (lista mínima de 16 a 18 parámetros), S5-04 y S5-06, S7-04, S7-08 y S7-10.
- **Q2-11 B (K47): S7, C1-22 y U4-12 leen los parámetros de avisos en `app_setting`.** Se quitó la excepción de `notification.params_json` y la marca sobre el lugar; S7 cita `getReportSettings` y `updateReportSettings` (`reports.*`) y `getAlertSettings` y `updateAlertSettings` (`alerts.*`) de C1 1.3. Las pruebas de configuración de S7 cambian el valor por esas operaciones y comprueban que `params_json` ya no trae `send_time`, `advance_days`, `weekday` ni `day_of_month`. S7-10 agrega el mes corto (Q2-66 B): sale el último día. Los contenidos de los resúmenes siguen PROVISIONALES hasta la puerta 4 (K31).
- **S5-07: NO CUMPLE en el fondo, no se marcó cumplido.**
  - `docs/estado/S5-terminos.md` existe con 3 veredictos «no verificado» (desafío anti-robot en las 3 portadas; ninguna página legal leída); `grep -c "Veredicto"` = 3 solo cumple lo mecánico.
  - Nota en S5-07, §3, §4, cambios e INDICE: se escaló a Axel (B él lee y dicta, C pendiente, D descartar) y la lectura asistida sigue apagada.
  - La respuesta de la puerta 2 bis (K53, aplicación sin lectura de tiendas; se registra en P2d y aún no está en `decisiones.md`) queda citada como contexto; el criterio, su id y su severidad no cambian.
- **S2 (contrato 1.3):** S2-09 y S2-16 citan la regla de ventanas contiguas como **PROVISIONAL (`S2-ventanas`, pregunta 4, caso 1)** y el caso 2 (ventanas que no duran un mes; 97375 contra 292125) como **PROVISIONAL sin decidir**. S2-09 quita la condición del bloqueo de proyección (RESUELTO, `provisional` null). S2-10 anota Q2-13 B y Q2-14 A como FIRMES y `variable_estimated_cents` (Q2-12) y `savings_planned_cents` (`S2-ahorro-quincenal`) como PROVISIONALES en `null`. S2-12 pone `months_left`, décimas, tres estados y meta vencida PROVISIONAL (`S2-metas-vencida`). S2-14 suma que un instante sin «Z» da `invalid_date`. S2-16 enumera las cuatro marcas vigentes.
- **S4 (contrato 1.5):** S4-06 con la cotización por elección (`choice_store`, `choice_presentation`; ejemplo del arroz, Walmart 1 kg a 2800); S4-03 con el respaldo de presentación PROVISIONAL (`S4-presentacion-respaldo`) y el orden fijo de listas; S4-12 con las reglas de mercado FIRMES (Q2-39 a Q2-42: tienda asignada, `visit_date`, total solo con precio, orden de grupos); S4-10 con quién confirma los empates (FIRME) y los 150 renglones `propuesto` PROVISIONALES (`S4-empates-sin-confirmar`; comprobé que `tabla.csv` tiene 150); S4-15 con las marcas vigentes.
- **S5 (contrato 1.5):** S5-03 con una sola «mejor tienda» (completas primero, `null` en incompletas) y el granel PROVISIONAL (`S5-carritos-granel`); S5-04 con S y Z configurables; S5-05 con la lectura PROFECO como regla técnica y filtro y liga normalizada PROVISIONALES «en uso A»; S5-06 con `ocr.timeout_seconds`; S5-12 sin productos repetidos; S5-13 con las marcas vigentes.
- **U4-05:** la regla general de §0.2 (Q2-09 A, cierra R5-1), el alcance del reparto (Q2-68 C) y la fórmula de P-26.14 (Q2-70 a) entran al criterio y a su Verifica.
- **U3:** U3-05 con T-2 en una sola variante (acumulados de C1 1.3, Q2-65 A) y U3-01 con la segunda tienda de la tarjeta 5 (`second_store_id`, Q2-86 A).
- **C1:** C1-13 pasa a C1 v1.3 (`info.version` 1.3.0); C1-21 a 212 operaciones (168 rutas, 255 esquemas); C1-15 nombra los bloqueos abiertos que sostienen marcas en `openapi.yaml`.
- **INDICE:** versiones de §10 al día; §6 con `S1-parametros-avisos` RESUELTO y la lista de bloqueos que hoy sostienen marcas; §12 nuevo con el origen de cada cambio y las PROPUESTAS.

## Verificación
- `make verify-criterios` (antes y después de editar): 11 comprobaciones esperadas; salida final: `criterios definidos: 261 | IDs distintos citados: 261 | citados que no existen: {}`, `citas verificadas: 321 | sin fuente: 0`, `TOTAL 261 (alta 201 / media 60 / baja 0) | archivos: 18`, `INDICE contra conteo propio: 0 diferencias`, `residuos «N-xx» y rutas /tmp en los criterios: 0 y 0; las 11 comprobaciones dan lo esperado`.
- `python3 -I docs/criterios/tools/conteo.py`: versiones S1 1.2, S2 1.3, S4 1.3, S5 1.3, S7 1.3, U3 1.3, U4 1.3, C1 1.3; sin diferencias con INDICE §2 y §10.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → `Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)`, rc 0 (registro en `/tmp/p2c-criterios/verify-all.log`).
- Comprobaciones de hechos citados: `grep -c Veredicto docs/estado/S5-terminos.md` = 3; `grep -c ",propuesto" datos/semilla/nutricion/tabla.csv` = 150; `info.version` de `openapi.yaml` = "1.3.0"; `cumulative_accrued_cents` aparece 4 veces en `openapi.yaml`.
- Reemplazos con conteo exacto (scripts en `/tmp/p2c-criterios/`, fuera del repo; cada uno aborta sin escribir si un texto no aparece las veces esperadas).
- Los §5 nuevos no llevan citas «…» nuevas (los K46 a K52 se citan sin comillas), así que `citas.py` no tiene nada nuevo que verificar.

## Supuestos
1. Dónde quedó cada cosa se tomó de los reportes de P2c (S1, contratos, C1, U4, motores) y de los contratos 1.3/1.5; no releí los motores ni el esquema, solo `S2-interfaz.md`, `S4-interfaz.md` y `S5-interfaz.md` en los puntos citados, `C1-datos.md` (claves y rangos) y `U4-ajustes-b.md` (claves de reportes y avisos).
2. «Lista mínima de U4-12 de 16 a 18»: cuenté `optimizer.date_window_days` y `ocr.timeout_seconds` como (17) y (18) y agrupé las cuatro horas en (4), los dos días de la semana en (15) y la anticipación en (5). Es contenido del criterio, no un criterio nuevo.
3. El contenido de los resúmenes y del reporte mensual sigue PROVISIONAL hasta la puerta 4: K50 aprobó horas y anticipación, no contenidos.
4. El S4-10 es donde puse la confirmación de empates (la nutrición vive en S4-10 y S4-08 solo pide fuente citada).
5. S5-07 conserva su texto de fondo: K53 no está en `decisiones.md` (el último es K52), así que no lo apliqué como decisión registrada.

## Bloqueos
No abrí ninguno. Los que sostienen marcas en los criterios están en `docs/estado/bloqueos/` (lista en INDICE §6): `S2-ventanas` (pregunta 4), `S2-flujo-quincenal`, `S2-ahorro-quincenal`, `S2-metas-vencida`, `S4-presentacion-respaldo`, `S4-empates-sin-confirmar`, `S5-carrito-comparativa`, `S5-carritos-detalle`, `S5-carritos-granel`, `S5-profeco-filtro`, `S5-terminos`, `C1-hora-reportes-job`, `C1-reparto-alcance` y `S3-iva-lectura`.

## Pendientes
- Repasar las listas «Cambios de la versión …» de C1, U4 y S7 si el revisor de P2c pide menos historial (conservé el historial anterior intacto).
- Cuando P2d registre K53: reformular o retirar S5-07 y revisar S5-04, S5-08 y S5-10 (referencias a `site_assisted` y a la lectura asistida) con la firma de Axel.
- Cuando WF-3 resuelva las preguntas sin recomendación (Q2-12, Q2-48, Q2-50 a Q2-52, Q2-56 a Q2-58, Q2-85, Q2-87 a Q2-91, Q2-94 y la pregunta 4 de `S2-ventanas`), quitar las marcas correspondientes en S2, S4, S5 y C1-15.

## Propuestas para Axel
Nada de esto está añadido a los criterios; ids y conteo siguen en 261.
1. **S5-07 y K53:** al registrar K53 en P2d, reformular S5-07 como «la aplicación no lee sitios de tiendas ni activa la lectura asistida» o retirarlo con tu firma. Hoy sigue NO CUMPLE.
2. **Criterio nuevo en C1 para el agente externo de precios** (importación de productos, presentaciones, precios con origen y fecha, y recetas con credencial de servicio propia), cuando P2d la defina. No lo escribí porque la API no existe.
3. **Criterios propios para Q2-67 (IVA del crédito), Q2-69 (a) (comisiones solo como dato) y Q2-82 (atajos de una tecla):** hoy los cubren U4-03, U4-04 y U4-07 sin criterio propio. Tú dices si los quieres.

## Solicitudes fuera de alcance
1. **`docs/design/DESIGN.md` l. 933 («Avisos y frecuencias»)** sigue PROVISIONAL por `S1-parametros-avisos`, ya RESUELTO con Q2-11 B (lo anotó U4-puerta2). Le toca al dueño de DESIGN o al integrador.
2. **U4-12 exige que `ocr.timeout_seconds` aparezca en Datos y respaldos** con el valor 30 y «Vuelve a 30»; `U4-ajustes-b.md` l. 361 y l. 387 ya lo traen, pero su nota dice «PROPUESTA de S1 que la puerta 2 no preguntó» para la sección (no para el valor, que es Q2-76). Está de acuerdo con lo que escribí; solo conviene que quien revise U4 lo sepa.
3. **U1** (N-297, N-300 y tabla K44 de §2.3) todavía dice que `send_time` y `advance_days` viven en `params_json` (lo pidió C1-puerta2); S7-04 a S7-10 y C1-22 ya no lo aceptan, así que U1 tiene que cambiarse a `reports.*` y `alerts.*`.
4. **Bloqueos INDICE:** el índice de `docs/estado/bloqueos/INDICE.md` es de otro agente; no lo toqué.

## Consumo aproximado
Acciones de herramienta: unas 58. KB leídos: unos 190 (criterios por tramos, contratos 1.3/1.5 en los puntos citados, reportes de P2c, `PUERTA-2.md` §4 y §5, `RESPUESTAS.md`, `decisiones.md` K31 a K52). KB escritos: unos 75 (diferencias de 8 archivos de criterios más este reporte).
