# Revisión U4 · ronda 2 de 3 (revisor-U4-r2)

Flujo WF-1 · 11-oct-2026 · Revisor independiente: no escribió nada de lo que revisa y no corrige nada. Evalúa solo contra `docs/criterios/U4.md` 1.2 y `docs/plan/decisiones.md` §A–§K. Objeto: `docs/design/U4-ajustes-a.md` 1.1 (P-26.01 a P-26.07) y `docs/design/U4-ajustes-b.md` 1.1 (P-26.08 a P-26.14), después de `correccion-U4-r1`.

## Veredicto
**NO APROBADO.** 1 hallazgo alto, 2 medios y 8 bajos (5 de ellos son los bajos H-4 a H-8 de la ronda 1, que siguen abiertos porque no se asignaron a la corrección). De los anteriores: H-1 está cerrado en los 7 casos que nombraba, pero U4-05 sigue sin cumplirse por otros ajustes con efecto numérico (R2-1). H-2 y H-3 quedan cerrados.

## Verificaciones propias
Scripts en `datos/tmp/revisor-U4-r2/` (`v.py`, `o.py`), corridos con `python3 -I` y borrados al terminar.
- **Secciones (U4-01):** separo las dos partes por `^## (P-26\.\d\d)` y salen 14 secciones, de P-26.01 a P-26.14, sin repetidas (7 en a y 7 en b). Las 14 tienen «Vista de base de datos», «Edición en el lugar», «Vista previa del efecto», «Estados», «Atajos», «412 px» y «1440 px». Los nombres coinciden con la decisión G.
- **operationId (U4-02):** junto los `operationId` de los 18 extractos (210, igual que `grep -c operationId: openapi.yaml`). Cito 100 en la parte a y 45 en la b: **0 inexistentes**. Las 73 operaciones de la etiqueta Ajustes (`ajustes-a.yaml` a `ajustes-e.yaml`) aparecen todas: **0 sin campo ni acción**. Por sección: entre 4 y 19 operaciones válidas. P-26.10 no tiene ninguna porque es local (U1 §0.3).
- **Parámetros K44 (U4-12):** cruzo las tablas contra el inventario `@param` de 00003 (líneas 38–52) y contra los `INSERT INTO app_setting` de 00001 (1393–1399), 00002 (143–149) y 00003 (80–86 y 139–149). Las 15 claves sembradas están en la sección que dice la columna «sección» del `@param`: 8 en la parte a (P-26.01, P-26.05, P-26.06) y 7 en la b (P-26.09, P-26.11, P-26.13). **0 claves sin sección y 0 parámetros sin clave.** Los valores iniciales 5, 15, 90, 12, 30, 90, 1, 5, 10, 24, 0, 0, 80, 100 y [1, 0] son iguales a los sembrados. Cada fila trae «Restablecer» o «Vuelve a …», y el valor vigente y el inicial van en la misma fila (§0.5 de a y tablas de b).
  - `send_time`, `advance_days`, `weekday` y `day_of_month` están en el `@param_sin_valor` de 00003 (líneas 56–59). Sus tablas por sección los marcan PROVISIONAL citando `S1-parametros-avisos`, que no tiene «Respuesta de Axel». `weekday` = 7 y `day_of_month` = 1 son FIRMES y tienen «Vuelve a 7/1». `send_time` y `advance_days` no tienen valor inicial y su «Restablecer» se activa cuando S7 proponga el valor.
  - «Mínimo de seguridad» aparece solo en SUPUESTOS (b, línea 388).
- **Otras claves sembradas citadas:** `optimizer.weights` (1, 1, 1, 1), `ai.local_model_enabled` y `ai.claude_api_enabled` (false), `price.assisted_reading_enabled` (false) y `projection.repeat_cadence_enabled` (false). Todas coinciden con la migración.
- **Cifras del oráculo** (Decimal, ROUND_HALF_UP):
  - Las series A y B de 14 valores, formateadas como `$x,xxx.xx`, aparecen literales en la parte a.
  - `minimo_A` 8094.5 y `minimo_B` 5000.0; `eventos[0]` `disponible_A` 10000.0 y `disponible_B` 6905.5.
  - A − B = 3,094.50 en las 14 posiciones, lo que sostiene el «+ Δ» del límite.
  - `despensas[0]`: 4 × 973.75 = 3,895.00 = 3,094.50 + 800.50; 3,094.50 × 25.8685 % = 800.50; 199 ÷ 2 = 99.50.
  - Meta de calorías: 1,748.75 × 1.2 = 2,098.50; 2,098.50 − 300 = 1,798.50, que con HALF_UP da 1,799.
  - RN-15: 946000 × 30 ÷ 21 → 1351429 milésimas, y 2838000 − 1351429 = 1486571.
  - RN-18: 3,146.72 ÷ 6,293.44 = 50.0 %.
  - K40: (90 − 30) ÷ 61 = 0.9836 y 30 ÷ 61 = 0.4918; 10-oct-2026 + 30 días = 9-nov-2026 y + 90 días = 8-ene-2027.
  - **0 diferencias.**
- **Texto y tokens (U4-09, U4-11):** 0 hex, 0 `rgb(` y 0 fechas ISO en las dos partes. Cito 21 tokens y todos existen en `C2-colores.json` o en `DESIGN.md`, incluido `size.control.primary-sheet` (DESIGN, línea 408). «Se ve(n) al guardar» = 0. «Suman más de 100» = 0, «deben sumar 100» = 0 y V-16 aparece 1 vez en a y 2 en b. `neutral.700` = 0 y «No se guardó» = 0.
- **PROVISIONAL:** 5 líneas en a, todas con su cita. 24 líneas en b, de las cuales 5 no citan bloqueo en la misma línea:
  - 26: remite a «Escritura de día y hora», que sí lo cita.
  - 104: redacción del texto de ayuda.
  - 209: hora del respaldo, cita K31 en su fila.
  - 361 y 362: tabla consolidada (H-7).
- **U1:** comparo con U1 §2.3 (tabla de las 14 y reglas K44), §2.4 (estados), §2.5.a y §2.5.b (vistas, filtros y orden), §2.6 (P-26.01 a P-26.14 con guardado automático; A/B y restaurar confirman), §5.2 (comandos) y §5.3 (atajos). DESIGN §7.13 da origen a los 2 s de «Guardado», al token `color.role.success.light`/`.dark` y al estado deshabilitado sin conexión.
- **Contratos y bloqueos:**
  - `CreditPurchaseCreate` (`creditos-b.yaml` 888–935): tiene `payment_frequency` y `first_payment_date`, y no tiene `first_payment_min_days`.
  - `progress_pct` es «entero ROUND_HALF_UP» (C1-datos, línea 136, y extracto).
  - Bloqueos abiertos que tocan U4: `S2-metas` y `S3-formula-tarjeta` (sin respuesta) y `S3-convencion-prestamo`.

## Tabla de criterios
| Criterio | Estado | Evidencia propia |
|---|---|---|
| U4-01 (alta) | Cumple | 14 secciones con los nombres de la decisión G. Cada una con tabla Campo · Tipo · Valor por defecto · Validación · Operación · N-xxx. Las dos partes usan formatos distintos (H-8, baja). |
| U4-02 (alta) | Cumple | 145 citas de `operationId`, 0 inexistentes. 73 de 73 operaciones de Ajustes con campo o acción. Los campos sin operación son locales o calculados y lo dicen: «— (Local)», «se compara en el dispositivo». |
| U4-03 (alta) | Cumple | Lista marcada completa. **Cuenta:** co-titulares, sesiones, Karen sin cuenta. **Presupuesto:** K2 preseleccionado. **Ingresos:** 15 y 30, extras, Ximena sin ingresos. **Categorías:** «Comida fuera». **Créditos:** 4 tipos, 17.27 % y 25.8685 %, tasa por compra, promociones con comercio, plazo, mínimo y vigencia, A/B = B con «pendiente confirmar», K12, K16 con «faltan restricciones», K18 con «falta capturar», K21. **Tiendas:** mercados con día y qué se compra, Sam's con costo anual, 30 días, lectura asistida bloqueada por V-19. **Recetas:** editables, restricciones vacías. **Reportes:** diario, semanal, mensual, Karen, vista web. **Notificaciones:** 4 alertas, destinatarios, candidatos K4. **Apariencia:** tema, saturación y densidad fijas. **Datos:** respaldo, integridad, restaurar. **Personas:** 50/50, Karen, tope K13. **Salud:** perfil, IMC, objetivo K23, metas. **Inventario:** modo y sobrantes. **IA:** apagada. |
| U4-04 (alta) | Cumple | Muestra de 15 valores por defecto, todos con origen en §A–§K o en la semilla de la migración: A/B = B (C), IA apagada (H), lectura asistida apagada (E), tema Sistema (G), Conjunto 50/50 (B), 30 días (E, K22), 80 y 100 % (F, K44), 5 intentos en 15 min (A, K11), tasas 17.27 y 25.8685 (C), «De este mes en adelante» (K2), meses 7 a 12 apagado (K21), restricciones sin valor (K16), tope vacío (K13), 90 días (K40), 0 y 0 kcal (K36). También 10 MB y 5 s (K43), 12 caracteres (K6), 1 GB (K6), 24 h (K6), pesos 1, 1, 1, 1 (C), bimestral y devengo de 2 meses (B), porciones 2 (D). «Hogar» como reparto por defecto viene de C1 (`is_default`, `ajustes-b.yaml` 749–751). |
| U4-05 (alta) | **No cumple** | Los 7 casos de H-1 (r1) ya tienen vista previa definida. Quedan ajustes con efecto numérico sin vista previa (R2-1). La vista previa del pago mínimo y la de metas dependen de bloqueos abiertos que no se citan (R2-2 y R2-3). El caso A/B obligatorio sí cumple, con cifras exactas (R-61, R-62, R-64, R-65, R-66). |
| U4-06 (alta) | Cumple | Vistas: tabla, tablero, calendario (promociones, tiendas, metas) y galería (líneas, tiendas, recetario, metas). Atajos: `Alt+N`, `Ctrl+Enter`, `Ctrl+/`, `↑`/`↓` con `Enter`, `Esc`, `Ctrl+K`, `Ctrl+\`, `Ctrl+Shift+L`, `?` y `←`/`→`, todos en U1 §5.3; 0 nuevos y 0 duplicados. Menú «/» y barra lateral en §0.7. Siguen abiertos los bajos H-5 y H-6. |
| U4-07 (alta) | Cumple | Catálogo V-01 a V-22. Muestra de 15 mensajes sin inglés ni códigos: V-01, V-03, V-05, V-07, V-08, V-12, V-14, V-16, V-17, V-18, V-19, V-20 y V-21, más «Escribe la hora con dos dígitos, como hh:mm» y «Escribe un número entero de kcal, 0 o más». Dinero `$1,234.56`; tasas con hasta 4 decimales y como se capturaron (K12). Una sola regla de reparto (V-16, RN-03, S1-08). |
| U4-08 (alta) | Cumple | Las acciones críticas usan §0.4 con `shield-alert` y `color.critical.light`/`.dark`, y el botón usa `color.critical.700`/`.300`. Son: cerrar sesiones (todas y todas salvo la actual), restaurar (doble confirmación con `RESTAURAR`), activar cada IA, activar la lectura asistida (general y por sitio), eliminar tasa, promoción o receta, quitar monto y descartar chat. Cada diálogo dice qué pasará. |
| U4-09 (media) | Cumple | 0 hex y 0 `rgb(`. 412 px y 1440 px en las 14 secciones. §0.7: `size.touch.min` de 44 px, `size.control.coarse`/`fine`, foco visible con `ui.border.focus`; nunca solo color. Texto viejo en la parte b (R2-6, baja). |
| U4-10 (media) | Cumple | Un solo bloque de estados (a §0.3) para las 14: vacío, carga (aviso a los 10 s), guardando, guardado (2 s, DESIGN §7.13), restablecido, error de validación, error al guardar, error al cargar y sin conexión. Las 14 secciones tienen «Estados» y remiten a §0.3. La confirmación es la misma en las 14 (§0.2). |
| U4-11 (media) | **No cumple** | 0 fechas ISO. SUPUESTOS en las dos partes y 4 bloqueos `U4-*` que existen. Pero dos vistas previas presentan como firmes cálculos que están en bloqueos abiertos de otros artefactos, sin PROVISIONAL ni cita (R2-2), y un bloqueo nuevo repite otro abierto con opciones distintas (R2-3). |
| U4-12 (alta) | Cumple | Cruce 15 de 15 con `app_setting` y su sección. Los 16 de la lista mínima tienen valor inicial (o «sin valor inicial: lo propone S7, PROVISIONAL»), validación, «restablecer» y decisión. Los 4 de `params_json` se comprueban aparte (arriba). La tabla consolidada de la parte b sigue incompleta (H-7, baja). |
| U4-13 (alta) | Cumple | P-26.13: dos valores de la casa (`health.goal_adjust_*`, 0 y 0, entero ≥ 0, «Vuelve a 0 kcal») y dos por persona (vacío = el de la casa, y lo dice). «Mantener» no tiene campo. Desglose con 2,098.50 − 300 = 1,798.50 → 1,799. «Usar el de la casa» funciona como «restablecer». N-326 y N-327. |

## Hallazgos

**R2-1 · alta · U4-05 · Quedan ajustes con efecto numérico sin vista previa antes de confirmar.**
La corrección cubrió los 7 casos de H-1, pero quedan estos cuatro con efecto numérico:
1. **P-26.13, Estatura, Fecha de nacimiento y Sexo** (`docs/design/U4-ajustes-b.md`, tabla «Perfil de salud» y «Vista previa del efecto»). Estos tres datos cambian el IMC y la meta sugerida (RN-17, Mifflin-St Jeor), pero la vista previa solo dice que el desglose se recalcula «al cambiar actividad, objetivo o ajuste». El IMC (K7) no tiene vista previa.
2. **P-26.05, «Tarjeta: día de corte y día límite de pago»** (`docs/design/U4-ajustes-a.md`, Bloque 1). Cambian `next_payment_date`, el «próximo pago» que el mismo bloque muestra en solo lectura. No tienen vista previa.
3. **P-26.05, «Préstamo: tasa anual»**. Cambia la amortización de los préstamos (RN-07 y RN-08, motor S3). No tiene vista previa.
4. **P-26.05, vista previa 4, condiciones de tarjeta o tienda**. Solo cubre el momento de completar los datos («al completar…»). No define el «Antes → Después» cuando se cambia una tasa, una base, el IVA o un mínimo que ya estaba capturado.

U4-05 exige 0 ajustes sin definir.
**Corrección sugerida.** Definir una vista previa en línea (§0.2) para cada uno:
1. Perfil: IMC y meta sugerida, antes → después, calculados en el navegador con K7 y RN-17 al cambiar cualquier dato del perfil.
2. Días de corte y de pago: próxima fecha de corte y próxima fecha de pago, antes → después.
3. Préstamo: PROVISIONAL citando `docs/estado/bloqueos/S3-convencion-prestamo.md`, con texto «Vista previa pendiente».
4. Tarjeta o tienda: el mismo cálculo de la vista previa 4, antes → después, sujeto a R2-2.

**R2-2 · media · U4-05 y U4-11 · La vista previa del pago mínimo de tarjeta o tienda usa como firme una fórmula que está bloqueada.**
`docs/design/U4-ajustes-a.md` P-26.05, vista previa 4, dice: «se muestra el pago mínimo del saldo actual según la fórmula general (K18)». Pero el bloqueo abierto `docs/estado/bloqueos/S3-formula-tarjeta.md`, pregunta (2), dice que la fórmula del pago mínimo no está decidida: puede ser el mayor entre % × saldo y el mínimo fijo, o % × saldo + intereses + IVA con piso en el fijo. La pregunta (3) deja abierto lo mismo para el crédito de tienda. U4 muestra la cifra sin la marca PROVISIONAL y sin citar el bloqueo. Quien implemente no sabe qué fórmula usar en la pantalla.
**Corrección sugerida.** Marcar el pago mínimo de la vista previa 4 como PROVISIONAL citando `S3-formula-tarjeta`. Mientras no haya respuesta, mostrar «Vista previa pendiente» (§0.2) y no una cifra.

**R2-3 · media · U4-05 y U4-11 · La vista previa de metas no cuadra con el bloqueo abierto `S2-metas` ni con C1.**
`docs/design/U4-ajustes-b.md` P-26.13 («Vista previa del efecto», metas) y `docs/estado/bloqueos/U4-vista-previa-formulas.md` (pregunta 2) tienen tres problemas:
1. **Bloqueo repetido con otras opciones.** La pregunta 2 hace la misma pregunta que el bloqueo abierto `docs/estado/bloqueos/S2-metas.md` (1): cómo se cuentan los meses que faltan en la aportación necesaria. Las opciones no coinciden:
   - U4: (c) es «días ÷ 30.4375».
   - S2: C es «meses completos entre hoy y `target_date`».
   - Con dos bloqueos, Axel puede contestar dos veces de forma distinta. Además, el bloqueo nuevo dice que C1 «solo dice “aportación mensual necesaria para llegar a target_date”», pero C1-datos M-16 (línea 353) da la forma «(monto meta − ahorrado) ÷ meses que faltan».
2. **Avance con un decimal, como si fuera firme.** La vista previa (1) muestra el avance «con un decimal (RN-18)», por ejemplo 50.0 %. En C1, `progress_pct` es «entero ROUND_HALF_UP» (C1-datos, línea 136). La misma pregunta está abierta en `S2-metas` (2). La columna «Ahorrado, avance…» de la misma tabla lee `progress_pct`, que es entero. El resultado es un «Antes → Después» con un formato distinto del valor guardado, sin PROVISIONAL.
3. **Tablero por estado sin marca.** El tablero «por estado (en ruta, atrasada, cumplida)» depende de `S2-metas` (3), que tampoco se cita.

**Corrección sugerida.**
1. Citar `S2-metas` en P-26.13 para los meses que faltan, el formato del avance y los estados.
2. Quitar la pregunta 2 de `U4-vista-previa-formulas` o redactarla con las mismas opciones que `S2-metas`, remitiendo a ese bloqueo.
3. Marcar PROVISIONAL el número de decimales del avance en la vista previa y los valores del tablero por estado.

**R2-4 · baja · U4-05 · La vista previa 7 de P-26.05 queda PROVISIONAL sin necesidad, y cita una solicitud en vez de un bloqueo.**
`CreditPurchaseCreate` ya acepta `first_payment_date` (`creditos-b.yaml`, cerca de la línea 923). El navegador puede calcular el primer pago con RN-06 y los días tecleados, como ya hace con las próximas fechas del gasto fijo, y mandarlo a `previewCreditPurchase` sin cambiar C1 (SC-U4-1). Además, la marca cita una solicitud del reporte `correccion-U4-r1`, no un archivo de bloqueo.
**Corrección sugerida.** Calcular el «Después» en el navegador y pasar `first_payment_date`. Si se mantiene la solicitud, dejarla como opcional.

**R2-5 · baja · tarea (g) · El mismo campo se edita en dos secciones, contra la frontera declarada.**
`docs/design/U4-ajustes-a.md` «Fronteras con la parte b» dice que las restricciones alimentarias «se editan aquí (P-26.07)… el resto del perfil, en P-26.13». Pero `docs/design/U4-ajustes-b.md` P-26.13 también incluye «Restricciones alimentarias» como campo editable (`updateHealthProfile`, `dietary_restrictions`).
**Corrección sugerida.** Elegir un lugar para editarlas y, en el otro, dejar el campo en solo lectura con un enlace. O declarar en las dos partes que se editan en los dos lugares.

**R2-6 · baja · U4-09 y U4-11 · Texto viejo en la parte b.**
- `docs/design/U4-ajustes-b.md`, línea 3: dice que «la §7 de componentes aún no existe», pero la parte a cita DESIGN §7.0 y §7.12 a §7.15, que sí existen.
- «Cobertura de criterios U4 en la parte b», U4-11: no menciona el bloqueo nuevo `U4-vista-previa-formulas`.

**Corrección sugerida.** Actualizar las dos frases.

**R2-7 · baja · U4-11 · Sigue abierto H-4 (r1).** P-26.01 «Estados» todavía dice «Contraseña cambiada. Tus otras sesiones siguen abiertas.», un comportamiento que `changeOwnPassword` no fija.
**Corrección sugerida.** La de la ronda 1: quitar la frase o declararla en SUPUESTOS con una solicitud a C1.

**R2-8 · baja · U4-06 · Sigue abierto H-5 (r1).** Estas vistas se apartan de U1 §2.5.a sin declararlo en SUPUESTOS:
- P-26.07: tablero «por origen».
- P-26.14: tablero por modo y sin galería, contra tablero por Estado y galería «Sí».
- P-26.12: tabla «Repartos reutilizables», cuando `regla_reparto` no tiene vista de lista.
- P-26.05: tabla de líneas, contra «Galería (líneas)» de §2.3.

**Corrección sugerida.** La de la ronda 1: declarar la elección en SUPUESTOS y pedir a U1 que alinee §2.3 y §2.5.a.

**R2-9 · baja · U4-06 · Sigue abierto H-6 (r1).**
- `/promocion` (U1 §5.2 → P-26.05) no aparece en a §0.7 ni en P-26.05.
- `/reporte` sigue atado a la acción de P-26.08, cuando U1 lo manda a P-23 (N-248).

**Corrección sugerida.** Citar `/promocion` y alinear el destino de `/reporte` con U1, o pedir el cambio a U1.

**R2-10 · baja · U4-12 · Sigue abierto H-7 (r1).** En la tabla «Parámetros operativos de la parte b», las filas 5, 15 y 16 dicen «lugar PROVISIONAL» sin citar `S1-parametros-avisos`. La tabla no tiene columnas de validación ni de «restablecer».
**Corrección sugerida.** Citar el bloqueo en cada fila y, al unir las partes, dejar una sola tabla de 7 columnas.

**R2-11 · baja · U4-01 · Sigue abierto H-8 (r1).** Las dos partes siguen con formatos distintos:
- Encabezados: `## P-26.01 · …` en a y `## P-26.08 …` en b.
- Tablas: 6 columnas en a y 7 en b.
- Tokens: a usa el token de grupo `color.role.warning` en P-26.05, y b usa tokens de hoja.

**Corrección sugerida.** Unificar al integrar.

## Lo que está bien (sin hallazgo)
- Las 73 operaciones de Ajustes están cubiertas y no hay ningún `operationId` inexistente.
- La tabla A/B usa las cifras exactas del oráculo y el «+ Δ» del límite se demuestra con A − B constante.
- La regla del reparto es una sola (V-16 y RN-03).
- Hay un solo bloque de convenciones, y su origen es DESIGN §7.13 y §7.14.
- Los 16 parámetros de U4-12 y las 15 claves de `app_setting` están en su sección. K36 está completo. No hay hex, fechas ISO ni mínimos de seguridad inventados.
- Las propuestas de K31 están escritas en el reporte de U4b, con su justificación.

## Estado de los anteriores (REVISION-U4-r1)
| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| H-1 (alta, U4-05) | **Cerrado en sus 7 casos, pero U4-05 sigue sin cumplirse (ver R2-1)** | P-26.05 tiene las vistas previas 5 (límite), 6 (frecuencia), 7 (días mínimos, PROVISIONAL) y 8 (comisiones, PROVISIONAL con bloqueo). P-26.02 tiene la (4), despensa programada, con las cifras de `despensas[0]`. P-26.13 y P-26.14 tienen vista previa en línea con RN-18 y RN-15. «Se ve(n) al guardar» = 0. |
| H-2 (media, U4-10 y U4-09) | **Cerrado** | a §0 es el único bloque de convenciones para las 14 secciones. «Guardado» usa `color.role.success.light`/`.dark` durante 2 s (DESIGN §7.13, línea 931). Hay un solo texto de error. «Sin conexión», el aviso a los 10 s, `shield-alert`, el foco visible y `size.control.*` están en a §0. La parte b remite a a §0. `neutral.700` = 0. |
| H-3 (media, U4-07) | **Cerrado** | V-16 es igual en a §0.6 y en b P-26.12. La persona del resto no tiene campo. «Suman más de 100» = 0. Ya no hay recálculo de la parte del hogar como campo editable. |
| H-4 a H-8 (bajas) | Abiertos, sin cambio | No se asignaron a `correccion-U4-r1`. Pasan a R2-7 a R2-11. |
