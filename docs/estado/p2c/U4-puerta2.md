# U4-puerta2 · Aplicación de la puerta 2 en U4, U3 y DESIGN

Agente `U4-puerta2` (Opus) · flujo P2c · 11-oct-2026, noche · base 1ccb468 más el árbol de trabajo de P2c (esquema 1.4 de S1, migración `00004_puerta2.sql`). Manda: `docs/estado/puerta2/RESPUESTAS.md` (modo «lo recomendado + excepciones», sin excepciones; Q2-11 = B), `docs/estado/puerta2/PUERTA-2.md` §2.4, §3.1, §3.22 a §3.26, §4 y §5, `docs/plan/decisiones.md` K46, K47, K50 y K51, `docs/estado/wf1/REVISION-U4-r5.md` y `docs/estado/p2c/S1-puerta2.md`.

## Hecho
**`docs/design/U4-ajustes-a.md` 1.4 → 1.5**
- **Q2-09 A (R5-1, media):** §0.2 tiene ahora una regla general para las 14 secciones: (1) toda lista o interruptor con vista previa la muestra al elegir, guarda con `Enter` o al salir y `Esc` vuelve (lista de ejemplos de P-26.02, P-26.04, P-26.13 y P-26.14); (2) `Ctrl+/` y el arrastre de tablero de un campo con vista previa la muestran antes de aplicar («Aplicar» o `Enter`; `Esc` devuelve la tarjeta); (3) cada sección dice qué muestra. P-26.04: vista previa al elegir, al soltar la tarjeta y con `Ctrl+/`. §0.2 «Confirmar» suma el alcance del reparto (P-26.12).
- **Bajos:** R5-2 (P-26.04 nombra los gastos fijos sin reparto propio que cambian, con la regla de la vista previa 3 (d)); R5-3 (P-26.05, vista previa 11 de promociones); R4-2 («cada cuánto» de 1 a 24, V-23, como `interval_count` de 00001).
- **Q2-67 opción 1:** fila «¿la tasa incluye IVA?» con `vat_included` (1, 0 y vacío = falta capturar) y `vat_on_interest_bp` solo con 0; «Falta capturar» con 5 datos (`vat`, RN-26 1.4); vista previa 4 (a) sin PROVISIONAL. Los dos casos que siguen abiertos citan `S3-iva-lectura` (fila y supuesto 9).
- **Q2-69 (a):** vista previa 8 sin PROVISIONAL (las comisiones solo se guardan como dato; RN 1.4).
- **K50 (Q2-74, Q2-75 y Q2-77):** §0.5 con los rangos de la puerta 2 y el tipo de hora; comparador de 1 a 30 s; `optimizer.date_window_days` = 7 como parámetro con «Vuelve a 7»; precio viejo ≥ 1 y peso cero ≥ precio viejo (lo validan el motor y la API, no el esquema). Mensajes nuevos V-23 a V-25. La tabla de parámetros de la parte a y la nota de las claves de la parte b están al día.

**`docs/design/U4-ajustes-b.md` 1.3 → 1.4**
- **Q2-11 B (K47):** P-26.08 y P-26.09 tienen una sola variante. Día y hora se leen con `getAlertSettings` y se escriben con `updateAlertSettings` (`parameters` y `reset_parameters`), con las 8 claves `reports.*` y `alerts.*` de la 00004; `updateNotification` solo cambia «Activo». Ya no hay marcas PROVISIONAL de `S1-parametros-avisos`.
- **Q2-72 y Q2-73 (K50):** valores iniciales con «Vuelve a…»: 21:30, domingo (7) a las 19:00, día 1 a las 08:00, avisos a las 09:00, precios viejos 1 día antes y plan sin armar el domingo (7).
- **Q2-66 B:** en el día del mes y en la vista previa: «en los meses que no tienen ese día sale el último día del mes».
- **Q2-68 C:** P-26.12 pregunta el alcance antes de guardar: «Solo a los gastos nuevos» (preseleccionada; supuesto 6) o «También a los gastos ya capturados», con «Confirmar». La escritura usa `updateCostShareRule` con un alcance que C1 v1.3 tiene que agregar (solicitud SC-U4-P2-1).
- **Q2-70 (a):** P-26.14 calcula los días restantes con piso(…) y el consumo semanal con HALF_UP(…), en milésimas. Ejemplo de RN-15: 32 días y 315.333 ml por semana. Lo comprobé con Python en centavos y milésimas enteros: 1,486,571; 32; 315,333.
- **Q2-09 A:** P-26.14 (lista de modo, arrastre y `Ctrl+/` con «{k} pasan a Reponer…») y P-26.13 citan la regla general.
- **Q2-76:** `ocr.timeout_seconds` (30 s) en P-26.11, con `getDataSettings`/`updateDataSettings`.
- **Q2-82 (K51):** interruptor «Atajos de una tecla» en P-26.10, encendido y local. DESIGN §9.6 dice «lo agregan U1 y U4».
- **R4-3:** las líneas 104, 209, 361 y 362 ya traen su cita o se reescribieron.

**`docs/design/U3-dashboard.md` 1.2 → 1.3**
- **Q2-65 A:** T-2 con `cumulative_accrued_cents` y `cumulative_cash_cents` de C1 v1.3, en una sola variante; se retira la variante (b). La fuente de capas, el pago adelantado, el resumen, las etiquetas, el ejemplo, A-11, la trazabilidad, el supuesto 9 y §15 quedan sin PROVISIONAL. Ninguna cifra cambia: 2,921.25, 3,895.00, 13,103.78 y 14,362.48 siguen como estaban.
- **Q2-86 A:** la tarjeta 5 dice «[diferencia] menos que [segunda tienda]» con `second_store_id`, y si el campo no llega vuelve al texto de antes. El ejemplo queda «$138.50 menos que Walmart».

**`docs/design/DESIGN.md`:** solo se quitaron las marcas de Q2-81 (§2, 135°), Q2-82 (tabla de WCAG 2.1.4 y §9.6) y Q2-84 (supuesto 9, que cubre también los 2 s de Q2-83). Son 4 líneas (`diff` contra HEAD: 8 líneas `<`/`>`). Ningún valor, token, hex ni C2 cambió. Q2-78, Q2-79, Q2-80 y Q2-83 no tenían marca en DESIGN, así que no se tocaron.

**`docs/design/U4-ajustes.md`:** lo regeneré con el procedimiento de `INTEGRACION.md` §2.5: encabezado e índice, parte a, línea en blanco y parte b. En el encabezado actualicé las versiones, el estado de revisión y Q2-11 B; el índice queda igual.

## Verificación
- `make verify-wf1-u4` → «U4-ajustes.md termina con parte a + parte b sin cambios: True · secciones P-26: 14 de 14».
- Script propio sobre `docs/contratos/extractos/*.yaml` (C1 1.2; `openapi.yaml` sigue en `version: "1.2.0"`, así que **C1 1.3 todavía no existe**):
  - 210 operationId, 73 con la etiqueta Ajustes;
  - el U4 unido cita 109 (parte a 100, parte b 45), con **0 inexistentes** y **0 operaciones de Ajustes sin cita**;
  - 7 + 7 secciones P-26;
  - tablas: parte a con 178 filas y parte b con 117, con 0 filas de columnas distintas y 0 celdas vacías;
  - 0 hex y 0 fechas ISO en las partes;
  - `params_json`: 0 en a y 1 en b (la frase que dice que ya no se usa);
  - DESIGN: «PROPUESTA» = 0 y K51 = 4.
- Los campos nuevos de C1 v1.3 (`vat_included`, `second_store_id`, `cumulative_*` y el alcance del reparto) se citan como campos y no como operaciones. **Nota:** cuando exista C1 1.3, hay que repetir la comprobación de operationId y de campos.
- `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` (registro en /tmp/u4p2/verify-all.log):
  - pasan: SQL (goose up, down-to 0 y up), Go, colorimetría, OpenAPI, WF-0, S1, U1 («TODO OK»), Criterios, Go test, Recetario, Tokens C2, **U4** y Paquete;
  - **falla verify-c1** con «gen2: tipo de parámetro desconocido: hora HH:MM», como S1 ya esperaba;
  - **falla verify-wf1-oraculos** por 16 diferencias en S2: el motor Go no emite `provisional` en `s2-liga-provisional` ni en `s2-proyeccion-7-12-cadencia`. Ninguna de las dos fallas es de U4.

## Supuestos
1. En DESIGN, «marcas PROVISIONAL» de Q2-78 a Q2-84 = las marcas «PROPUESTA… lo decide Axel en la puerta 2» y la condicional del supuesto 9. El «PROVISIONAL» de la línea 933 es de Q2-11 y no se tocó (ver solicitudes).
2. P-26.12 preselecciona «Solo a los gastos nuevos», en analogía con K2 y por la razón de Q2-68 C («no cambiar sin querer estados de cuenta ya entregados a Karen»).
3. Las claves `reports.*` van por `getAlertSettings`/`updateAlertSettings`, como la variante B de U4 1.3. La operación final la fija C1 v1.3.
4. «Atajos de una tecla» usa N-301, porque U1 §3.28 no tiene una necesidad propia para él.
5. Numeración: parte a 1.5, parte b 1.4, U3 1.3. DESIGN conserva su encabezado.

## Bloqueos
No abrí ninguno. Sigue abierto y citado `docs/estado/bloqueos/S3-iva-lectura.md` (préstamo con IVA vacío; «No» sin porcentaje), sin recomendación, que pasa a WF-3. Sin respuesta y sin tocar: Q2-85 (sin recomendación).

## Pendientes
- Revisor de P2c: revisar U4 a 1.5, b 1.4 y U3 1.3.
- Siguen abiertos los bajos R3-2 a R3-6 de U4 (no se pidieron).
- Repetir la comprobación de operationId cuando exista C1 v1.3.
- Si `S3-iva-lectura` pregunta 1 = B, P-26.05 suma el IVA del préstamo.

## Solicitudes fuera de alcance
1. **C1 v1.3:**
   - SC-U4-P2-1, el alcance en `updateCostShareRule` (Q2-68 C), con el conteo de gastos ya capturados para la vista previa;
   - `second_store_id` en `DashboardBestBasket` (Q2-86 A);
   - `cumulative_accrued_cents` y `cumulative_cash_cents` en `AccrualVsCash.months[]` (Q2-65 A);
   - claves `reports.*` y `alerts.*` en `getAlertSettings`/`updateAlertSettings`, `optimizer.date_window_days` en `parameters`, `ocr.timeout_seconds` en `getDataSettings`, `vat_included` y el tipo de hora (todo esto ya lo pide S1).
2. **DESIGN.md l. 933 («Avisos y frecuencias»):** sigue PROVISIONAL por `S1-parametros-avisos`, ya RESUELTO con Q2-11 B. Le toca al dueño de DESIGN o al integrador; mi alcance era solo Q2-78 a Q2-84.
3. **U1 §2.3 y §3.28:** Apariencia pasa a tener dos controles (Q2-82, K51).
4. **Motor S2 / integrador:** las 16 diferencias de `verify-wf1-oraculos`.
5. **Criterios U3/U4:** quitar las menciones a las variantes de T-2 y a la PROVISIONAL de U4-12, si las tienen.

## Propuestas para Axel
- Confirmar que el alcance del reparto preselecciona «Solo a los gastos nuevos» (supuesto 2).
- Confirmar que `ocr.timeout_seconds` va en «Datos y respaldos» (la misma propuesta que hizo S1).

## Consumo aproximado
Unas 50 acciones de herramienta y unos 250 KB leídos, todo por partes con `sed -n` y `grep`. Escrito: unos 12 KB de texto nuevo neto en los 5 archivos de diseño, más un script de 30 KB en /tmp y este reporte. Los cinco archivos se reescribieron completos con el script, con respaldo en /tmp/u4p2/bak. Diferencias reales:
- numstat (+ − archivo): 4	4	docs/design/DESIGN.md
- numstat (+ − archivo): 17	20	docs/design/U3-dashboard.md
- numstat (+ − archivo): 33	22	docs/design/U4-ajustes-a.md
- numstat (+ − archivo): 45	39	docs/design/U4-ajustes-b.md
- numstat (+ − archivo): 80	63	docs/design/U4-ajustes.md
