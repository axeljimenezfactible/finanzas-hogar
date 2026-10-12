# Bloqueo S6-gastos-inicio

Agente api-gastos-inicio (WF-3a), 12-oct-2026. Cosas que el contrato y las reglas no fijan; cada una quedó PROVISIONAL en el código y se puede cambiar sin migración, salvo donde se dice.

## 1. Reparto RN-03 fuera del motor
- Qué falta: `motor_fin` reparte con `splitRule`, pero esa función no es exportada y el motor no se toca. El alta y el cambio de gastos necesitan repartir.
- Lo que hice (PROVISIONAL): `dominio.ExpenseSplitByRule` con la misma fórmula (porcentaje `ApplyBasisPoints` HALF_UP, fijo tal cual y el resto a la persona del resto). Agrupo por regla como dice RN-03. Probado con el ejemplo de la despensa 1 (Axel 190300 y Karen 119150).
- Opciones: A) exportar `motor_fin.SplitRule` (cambio en un motor aprobado: lo decide Axel) y llamarlo desde dominio; B) dejarlo en dominio y agregar al oráculo S2 un caso que lo compare.
- Impacto: ninguna cifra cambia; solo dónde vive la fórmula.

## 2. Confianza de los renglones del OCR
- Qué falta: `TicketScan.lines[].confidence` admite high, medium o low, pero el esquema 1.4 guarda solo `attachment.ocr_text` y no la confianza por renglón, y no hay umbrales.
- Lo que hice (PROVISIONAL): todos los renglones propuestos salen con `low`, así la pantalla P-05 no pliega nada como seguro.
- Opciones: A) columna nueva `attachment.ocr_lines_json` con la confianza del TSV (migración nueva) y umbrales en Ajustes; B) umbrales fijos de 85 o más para alta y de 60 a 84 para media; C) dejarlo en low.
- Impacto: P-05 pliega o no los renglones de confianza alta.

## 3. Promedio mensual del resumen
- Qué falta: «total ÷ meses del periodo» no dice cuántos meses tiene una quincena, una semana o un rango libre.
- Lo que hice (PROVISIONAL): cuento los meses calendario que toca el periodo y llamo a `motor_fin.MonthlyAverage`.
- Opciones: A) así; B) omitir `monthly_average_cents` si el periodo es menor que un mes; C) días ÷ 30.
- Impacto: solo esa cifra del resumen.

## 4. Periodo anterior equivalente del resumen
- Lo que hice (PROVISIONAL): para un mes, el mes anterior; para una quincena, la anterior; para un año, el anterior; para una semana o un rango, el mismo número de días justo antes.

## 5. Devengo y efectivo por gasto (`accruals`, `cash_flows`)
- Qué falta: para el devengo hay que llamar a `motor_fin.Periods`, y este exige en `params` los códigos de las categorías de despensa y de papel. Solo `groceries` está documentado en el esquema; el código de papel no está. Las compras a crédito necesitan además su plan de pagos.
- Lo que hice (PROVISIONAL): `accruals` no se manda. `cash_flows` solo se manda en gastos de contado, débito o transferencia: la fecha del gasto por el importe (RN-02).
- Pregunta: ¿cuál es el código de la categoría de papel? ¿Se lee de `app_setting` o es un código fijo?

## 6. Gastos que vienen de una compra a crédito
- Lo que hice (PROVISIONAL): `updateExpense` responde 409 `has_credit_purchase` si cambian el importe, la fecha o la forma de pago de un gasto con compra vinculada. `duplicateExpense` de un gasto a crédito responde 409 `has_credit_purchase`: la compra se registra en Créditos.

## 7. Filtro `category_id`
- Lo que hice (PROVISIONAL): coincidencia exacta, sin incluir subcategorías.
- Pregunta: ¿al filtrar por una categoría padre entran sus hijas?

## 8. Nombre del paquete
- **Cerrado** (correccion-backend-core-r1, 12-oct-2026, hallazgo M2): no se movió el paquete; CONVENCIONES-backend 1.1 §2 y §14 ya dicen `internal/api/gastos-inicio` con `package gastos`.
- El ALCANCE dice `backend/internal/api/gastos-inicio/`, y CONVENCIONES-backend.md §2 dice `internal/api/gastos`. Usé el directorio `gastos-inicio` con `package gastos` (`gastos.Register`). El integrador puede moverlo sin cambiar el código.

## 9. «Compras particulares» de Karen (`getPersonParticipation`, N-026 punto 2)
- **Qué falta**: U1 N-026 pide «sus compras particulares a crédito vivas», pero C1 no dice cuáles cuentan ni de quién son el total, los pagos hechos y el saldo.
- **Opciones**:
  - A (PROVISIONAL, implementada): compras con `status = active` en las que ella es la única persona con parte mayor que 0. El total es el de la compra. Pagos hechos son los pagos con su parte cubierta, de los pagos de la serie. Saldo = lo que le falta.
  - B: todas las compras vivas donde ella tiene parte, incluida la despensa común, con su total y su saldo.
- **Estado de una parte** (RN-04): está pagada si la suma de `credit_payment` con `payer_person_id` = ella y ese `scheduled_payment_id` llega al monto de su parte. Deuda viva = Σ de lo que falta en sus partes de compras no canceladas.
- **Impacto**: qué compras salen en la tarjeta D-7.

## 10. Filtro `period` en las tarjetas (`getDashboardMonthVsPlan`, `getDashboardFortnightFlow`)
- **Qué falta**: C1 acepta `period`, pero no dice qué mes o qué quincena se toma si el periodo es más largo.
- **PROVISIONAL**: se toma el mes o la quincena que contiene el inicio del periodo. Sin filtro, se usa hoy.
- `days_left`:
  - mes en curso: del día de hoy al último día, contando hoy;
  - mes pasado: 0;
  - mes futuro: todos sus días.
- Sin plan (`plan` 0 o ausente), `planned_cents`, `used_pct`, `remaining_cents` y `status` van ausentes. Los umbrales de `status` son `alerts.budget_warning_pct` y `budget_limit_pct` (K44), que lee `LoadFinEngine`.

## 11. `dominio.LoadFinEngine` falla con categorías del usuario (`code` NULL) · BLOQUEA EN USO REAL
- **Cerrado** (correccion-backend-core-r1, 12-oct-2026, hallazgo H1): `LoadFinEngine` lee `code` en `sql.NullString` y a las categorías sin código les da la clave interna `cat#<id>` (no `category_<id>`, que sí cabe en el CHECK `a-z0-9_` y podría chocar con un código real). Detalle en `S6-presupuesto-motor` punto 5.
- **Qué pasa**: `presupuesto_motor.go` (api-presupuesto) lee `category.code` en un `string`. El esquema permite `code` NULL en las categorías del usuario. Con una sola de ellas, toda operación que arma la entrada del motor S2 responde 500: mis 5 tarjetas y la tendencia, y también las de presupuesto.
- **Cambio pedido** (fuera de mi alcance):
  - escanear `code` en `sql.NullString`;
  - si es NULL, usar un código sintético estable y único, por ejemplo `fmt.Sprintf("category_%d", id)`, y guardarlo en `CategoryCode` y `CategoryID`.
- **Mientras tanto**: mis pruebas le ponen `code` a la categoría «Varios».

## 12. Búsqueda global (`search`)
- **Relevancia** (U1 N-002 no la define). PROVISIONAL: primero los títulos que empiezan con el texto y después los que lo contienen. Dentro de cada grupo, fecha descendente e id.
- **Secciones de Ajustes**:
  - `id` = clave de la pantalla de U4 (`P-26.01` a `P-26.14`);
  - `title` = nombre corto de la sección.
  - PROVISIONAL: frontend-base puede preferir la ruta.
- **«Nombres de página»** (U1 N-002) no tiene tipo en el enumerado de C1. No se buscan. PREGUNTA: ¿se agrega un tipo `page` en C1?
- **Ingresos**: solo `kind = extra` (U1: «ingresos extra (concepto)»).
- **Gastos**: título = descripción o, si no hay, tienda o categoría; subtítulo = categoría.
- **Compras a crédito**: título = descripción o tienda; subtítulo = la línea.
- **Tiendas**: subtítulo = municipio.

## 13. Tarjeta «Próximo pago y disponible» (`getDashboardNextPayment`, N-022) · continuación 2

- **Qué falta**: C1 no dice (a) qué hacer si varios pagos del mismo día son de líneas distintas (`credit_id` es uno solo), (b) si `payment_count` y la suma cuentan los pagos ya liquidados de ese día, (c) qué mostrar si el motor S3 no da disponible para una línea, (d) qué hacer con un pago programado sin partes.
- **PROVISIONAL aplicado**: (a) se suman todos los pagos pendientes del día y `credit_id` es el del primero (fecha, id); (b) solo cuentan los pagos con alguna parte pendiente; (c) la línea se omite de `availability` y se anota un aviso en el registro; (d) sin partes, todo es del hogar y cuenta como pendiente. Las compras `cancelled` no cuentan.
- **Duplicado de código**: la entrada del motor S3 se arma igual que `internal/api/creditos/motor.go` (`snapshot.state`), que no está exportada. **PROPUESTA**: que api-creditos exporte `func CreditStates(ctx, q, today) (map[int64]mc.CreditState, error)` (o lo mueva a `internal/dominio`) y que `gastos` lo use; así hay un solo armado.
- **Impacto**: bajo; las cifras salen del motor y de las partes guardadas.

## 14. Calendario del hogar (`listCalendarEvents`, N-029) · continuación 2

- **HECHO**: `income` (ingresos recibidos, `captured`), `credit_payment` (pago programado con partes y estado `paid`, `pending` u `overdue`), `fixed_expense` (`captured`) y `upcoming_expense` (`pending` u `overdue`) de `expense`, `card_statement` (corte capturado, monto = saldo al corte), `market_visit` (con `precios.NextVisit`), `promotion` (franja `valid_from`..`valid_to`), `savings_goal` (fecha meta de metas activas), `savings_movement` y `monthly_report` (día `reports.monthly_day_of_month`, recortado al último día en meses cortos; `ref_table = notification`).
- **PREGUNTA (no se emite todavía)**: los eventos **esperados** de partidas planificadas (`income` y gastos fijos o venideros con `status = expected`, y `planned_grocery`) necesitan las ocurrencias de `planned_item`, que calcula `motor_fin` (`occurrencesIn`, no exportada). Recalcularlas aquí sería repetir el motor. Opciones: A) exportar en `motor_fin` una función de solo lectura `PlannedOccurrences(in Input, from, to Date) []Occurrence` sin cambiar cifras (requiere OK de Axel por ser motor aprobado); B) dejar el calendario solo con lo capturado. **Recomendación: A.**
- **PREGUNTA**: `meal_plan` (estados `complete`, `incomplete`, `empty`) no tiene tabla en `x-tablas` de C1. Opciones: A) leer `meal_plan_entry` por día y marcar el estado por día; B) quitar el tipo de C1. Por ahora no se emite.
- **PROVISIONAL**: el filtro `family` solo distingue eventos con partes (pagos de crédito, gastos fijos y venideros); los eventos sin partes cuentan como del hogar. `credit_id` deja solo los eventos ligados a una línea (pagos, cortes y promociones de esa línea). Los cortes proyectados por `statement_day` sin estado de cuenta capturado no se emiten.
- **Impacto**: medio; el calendario funciona con lo capturado, pero no muestra lo planeado hasta resolver A.

## 15. Operaciones que dependen de código no exportado de otras áreas · continuación 2

- `getDashboardWeekMeals` (N-023): el costo y las kcal de la semana salen de `recetas.Engine.WeekPlanFor`, cuya entrada arma `internal/api/comidas` (`area.loadMotor`, ~100 líneas, no exportada). **PROPUESTA**: que api-comidas-salud exporte `func WeekPlan(ctx, q, today, monday dominio.Date) (*recetas.WeekPlan, []dominio.MealPersonRow, []dominio.HealthProfileRow, error)`.
- `getDashboardBestBasket` (N-024): la canasta y su comparación salen de `internal/api/comprasa` (`handlers.buildList`, `world.book.CompareBasket`, `membership`; no exportadas). **PROPUESTA**: que api-compras-a exporte `func BestBasket(ctx, q, today, monday dominio.Date) (api.DashboardBestBasket, error)` con la misma regla de `getBasketComparison`.
- Mientras no existan, las dos siguen en 501 (no se duplica el armado de los motores S4 y S5).
- **Estado tras correccion-backend-core-r1 (12-oct-2026, hallazgo H2)**:
  - `getDashboardWeekMeals`: HECHO. `comidas.WeekMeals(ctx, reg, monday)` (archivo nuevo `internal/api/comidas/inicio.go`, área api-comidas-salud) y el manejador en `inicio.go`; semana del filtro `period` o la semana en curso (K8). Prueba `TestOp_getDashboardWeekMeals` (en `internal/api/comidas/inicio_test.go`). Qué cuenta como «día con plan»: PROVISIONAL, `S6-comidas-calculos` punto 14.
  - `getDashboardBestBasket`: **sigue en 501**. Necesita una exportación de `internal/api/comprasa` (api-compras-a), fuera del alcance de la corrección; pedida en `solicitudes_fuera_de_alcance` de `docs/estado/wf3/correccion-backend-core-r1.md` con la firma y el cambio exacto. Mientras tanto fallan a propósito `TestAreas_SinOperacionesEn501`, `TestPermisos_MismasOperacionesAmbasCuentas` y `TestContrato_Rutas` («sin implementar 1»).
  - **Continuación 1 de correccion-backend-core-r1 (12-oct-2026)**: la exportación de `internal/api/comprasa` sigue sin existir; `getDashboardBestBasket` sigue en 501 y es lo único que deja en rojo `make verify-all` (Makefile:233) y `make verify-wf3` (Makefile:326 y :331). No se duplicó el armado de la lista y de la comparación en este paquete (unas 300 líneas de `comprasa`: `loadWorld`, `buildList`, `CompareBasket`, `membership`), porque la tarjeta debe dar las mismas cifras que `getBasketComparison`. Vale la solicitud de `docs/estado/wf3/correccion-backend-core-r1.md` (firma `comprasa.WeekBasketComparison`); sustituye a la PROPUESTA `BestBasket` de arriba.

- **Estado tras correccion-backend-areas-r2 (12-oct-2026, hallazgo H-1 de REVISION-backend-areas-r2)**: `getDashboardBestBasket` **HECHO**; ya no hay operaciones en 501 (216 de 216).
  - api-compras-a: `comprasa.BestBasket(ctx, q dominio.Querier, today, monday dominio.Date) (api.DashboardBestBasket, error)` en `internal/api/comprasa/comprasa.go` (la firma del hallazgo; sustituye a la solicitud `WeekBasketComparison` de correccion-backend-core-r1) y `comprasa.WeekParam(r)` (misma lectura de `week` que el área: falta → 422, otro formato → 400). Para no duplicar, `lista.go` separa `comparableList`, `expenseCart` y `firstComplete`, que ahora usan también `getBasketComparison` y `compareExpense` (mismo comportamiento: sus pruebas siguen en PASS). Las cifras salen tal cual de `precios.CompareBasket` / `CompareExpense` (motor S5); ninguna se calcula aparte.
  - api-gastos-inicio: manejador `getDashboardBestBasket` en `inicio.go` y `reg.Handle` en `gastos.go`. Pruebas: `TestOp_getDashboardBestBasket` (`gastos-inicio/mejor_canasta_test.go`: 200 validado contra openapi.yaml en vacío, por compra y por lista, cada uno igual a `getBasketComparison`; 422, 400, 401) y `TestBestBasket_*` (`comprasa/mejor_canasta_test.go`, una con la semilla real).
  - HECHO (regla del motor, no decisión nueva): ganadora = `cheapest_store_id` de S5 §7.4 (la primera tienda completa); si ninguna tienda está completa no hay `store_id` y la tarjeta queda vacía, igual que la comparativa.
  - **PROVISIONAL 15.1 · «lista vacía»**: la lista de la semana se considera vacía si no tiene artículos con producto y cantidad por comprar (los de solo texto no se pueden comparar); entonces se usa la última compra con renglones (`expense_date` DESC, `id` DESC). Opciones: (a) esta; (b) solo si no hay lista o no tiene renglones. Impacto: con (b), una lista de puro texto deja la tarjeta vacía en lugar de comparar la última compra.
  - **PROVISIONAL 15.2 · `item_count`**: el mismo número que `listBaskets` para esa canasta (lista: todos sus renglones; compra: sus renglones). Opciones: (a) esta; (b) solo los artículos comparables. Impacto: con (b), «[covered] de [items]» siempre diría N de N, porque la ganadora es una tienda completa.
  - **PROVISIONAL 15.3 · `covered_count`**: `identical_count + substitute_count` de la ganadora (un sustituto cuenta como «con precio»). Opciones: (a) esta; (b) solo idénticos. Impacto: con (b) baja la cobertura cuando la ganadora usa sustitutos.
  - **PROVISIONAL 15.4 · segunda tienda** (PREGUNTA para Axel): `second_store_id` y `difference_vs_second_cents` son la tienda siguiente del orden de S5 y su `difference_vs_cheapest_cents`. S5 §7.4 deja esa diferencia en null si la tienda está incompleta. En ese caso la tarjeta no trae segunda tienda (U3: «Ausente si no hay segunda tienda»). Opciones: (a) esta, que es la del motor; (b) restar los totales aunque la segunda esté incompleta, que sería una cifra que S5 no da. Impacto: el ejemplo de U3 («$138.50 menos que Walmart», 61 contra 60 artículos en `data.json`) solo sale si Walmart queda completa en la comparación.
  - **PROVISIONAL 15.5 · compra que no se puede comparar**: puede pasar porque la compra no tiene tienda activa, porque no tiene renglones del catálogo, o porque el motor responde «argumento inválido» (p. ej. «sin canastas que comparar»: ninguna otra tienda tiene precios). En esos casos la tarjeta responde 200 vacía (`basket` = `expense`, `expense_id`, sin `store_id`): es el vacío de U3, «Aún no hay precios para comparar». `getBasketComparison` sigue respondiendo 422 en esos casos. Opciones: (a) esta; (b) el mismo 422. Impacto: con (b), Inicio mostraría el error «Revisa que estés en la red de casa», y eso pasa hoy con la semilla real: no trae precios, que los carga el agente externo (K53).

## 16. Vistas de base de datos (`queryView`, C1 §8) · continuación 3

HECHO: las 23 entidades responden en los modos de C1 §8, con los `group_by` y `date_field` de su tabla. Otro valor de `group_by` o `date_field` responde 422, igual que un modo que la entidad no tiene. El formato malo responde 400. Una entidad fuera del enumerado responde 404. Lo siguiente lo decidí yo y queda **PROVISIONAL** hasta que Axel lo confirme:

1. **Periodo por defecto**: C1 dice «el periodo por defecto de cada operación», pero no le da uno a `queryView`. Sin `period` no se filtra por fecha. Opciones: (A) sin filtro, lo actual; (B) el mes en curso, como en `listExpenses`. Impacto: con (B), una tabla de gastos sin periodo solo trae el mes.
2. **Tablero sin `group_by`**: se agrupa por la primera propiedad permitida en orden alfabético. Opciones: (A) eso; (B) 422 si falta `group_by`; (C) fijar una propiedad por entidad.
3. **Filtros globales que no aplican a la entidad** (p. ej. `person_id` en `store`): no filtran. `category_id` y `family` solo aplican a `expense`, con las mismas reglas que `listExpenses`. Opciones: (A) ignorarlos, lo actual; (B) responder 422.
4. **`status`**: solo lo tienen `credit_purchase`, `scheduled_payment`, `savings_goal`, `inventory_item`, `notification_delivery`, `job` (con su `last_status`) y `backup`. En las demás entidades responde 422, igual que un valor fuera de su enumerado.
5. **Columnas del tablero**:
   - con enumerado, salen todas las columnas, aunque tengan 0 registros, en el orden del enumerado;
   - con id, solo salen los grupos que tienen registros, ordenados por etiqueta;
   - un valor NULL va en la columna `""` con la etiqueta «Sin valor».
   - `total_cents` y `share_bp` solo van si la entidad tiene importe principal. `share_bp` es el total del grupo ÷ el total de la vista, con HALF_UP.
6. **PREGUNTA, etiquetas en español de los enumerados**: C1 no las trae. Las propuse a partir del diccionario de datos, por ejemplo «Venidero», «Línea revolvente», «Precio viejo», «Reponer» o «Agotado». Si U1 tiene otras, hay que pasarlas a `vistas.go` (`lb*`). Los días del calendario se escriben así: «lunes 12 de octubre de 2026».
7. **`record`**:
   - lleva las columnas de la tabla: NULL ausente, 0/1 como booleano e instantes en −06:00;
   - suma los valores que calcula un motor: `status`, `saved_cents`, `live_debt_cents`, `freshness` e `is_stale`;
   - nunca lleva `password_hash`, `token_hash`, `idempotency_key` ni `ocr_draft_json`;
   - no lleva las propiedades anidadas de `components.schemas.<Entidad>`, como `Expense.shares`. Para eso está el GET del recurso.
8. **`product.freshness`**: sale de `precios.Freshness` aplicado al último precio de cualquier presentación. api-compras-a la calcula con `QuoteFor`, el mejor precio, y el resultado puede cambiar. **PROPUESTA**: que api-compras-a exporte `ProductFreshness(ctx, q, today)`.
9. **`scheduled_payment.status`**: copia la regla de api-creditos (`pagos.go`, `scheduledView`). Si lo abonado cubre el monto, va «paid»; si la fecha ya pasó, «overdue»; si no, «pending». **PROPUESTA**: exportarla desde `internal/dominio`.
10. **Calendario de tiendas**: un mercado aparece en cada fecha de su `visit_weekday` (ISO 8601) dentro del `period` o, si no hay, del mes en curso.
11. **`q` fuera de gastos**: busca sin mayúsculas ni acentos en el título y en cualquier propiedad de texto del registro.
12. **Títulos** sin columna de nombre:
    - `backup`: «Respaldo <id>»;
    - `scheduled_payment`: «<compra> · pago <n>»;
    - `credit_term_rate`: «<crédito> · <n> pagos»;
    - `meal_plan_entry` sin receta: «Comer fuera» o «Sin receta».

## 17. Lectura de tickets (`createTicketScan`, N-075, K3) · continuación 3

HECHO:
- La foto se guarda con `adjuntos.Store`, con la misma carpeta y el mismo límite `attachment.max_size_mb` que `uploadExpenseAttachment`.
- Se crea un `attachment` `ticket_photo` sin gasto y con `ocr_status = pending`, más su `audit_log`. La respuesta es 202 con `TicketScan`.
- El job `ticket_ocr` lee la cola con `ocr.Engine.Recognize`; solo hace la llamada. Guarda `ocr_text` (el texto crudo de cada renglón), `ocr_engine = "tesseract 5.5.3 spa"` y `ocr_status`, que queda `done` o `failed`, también con `audit_log`.

**PROVISIONAL**:
1. **Horario del job**: cada 5 minutos, y la subida lo pide «ahora» con `Jobs.RunNow`, así que la lectura empieza en el siguiente tick. Es un valor operativo de K31 (puerta 4). Justificación: si el servidor se reinicia, ninguna foto se queda sin leer más de 5 minutos.
2. **PREGUNTA**: si el OCR falla, por ejemplo porque Tesseract no está, el ticket queda `failed` y se captura a mano. C1 no tiene una operación para reintentar. Opciones: (A) dejarlo así; (B) reintentar N veces antes de marcarlo `failed`.
3. **Tipos admitidos**: solo imágenes, JPG, PNG o WebP. Un PDF responde 415 y se sugiere adjuntarlo al gasto. El formulario solo admite el campo `file`; cualquier otro responde 400.

## 18. Confirmar ticket (`confirmTicketScan`, N-074, RN-23, K3) · continuación 4

HECHO:
- Todo pasa en una sola transacción:
  - el gasto, con el mismo alta que `createExpense` (reparto RN-03, `audit_log`);
  - un precio por presentación, con fuente `ticket_ocr`, la fecha del ticket, la tienda del gasto, `expense_item_id` y `attachment_id`;
  - en los productos con existencia por unidad, se suma lo comprado: cantidad × contenido con `recetas.PurchaseMovement`, más un `inventory_movement` `purchase` ligado al renglón;
  - con `payment_method = credit`, la compra a crédito: `motor_cred.QuotePurchase`, partes, pagos programados y sus partes, ligada al gasto;
  - al final, el ticket queda `confirmed` con su `expense_id` y el borrador vacío (RN-23).
- Ya confirmado = 409. Cualquier 4xx deja la base igual que antes.

**PROVISIONAL** (cada punto se puede cambiar sin tocar el esquema):
1. **Sin tienda no hay precios**: `price.store_id` es obligatorio. Si un renglón trae presentación y el gasto no trae `store_id`, la respuesta es 422 en `expense.store_id`. La otra opción es confirmar sin guardar precios.
2. **Precio sin `unit_price_cents`**: se toma el importe del renglón entre las presentaciones, con HALF_UP (5800 ÷ 2 = 2900; 3125 ÷ 1.25 = 2500).
3. **Misma presentación dos veces en el ticket**: se guarda un solo precio, el del primer renglón. La llave UNIQUE de `price` no admite dos.
4. **Ese día ya había precio `ticket_ocr` regular de esa presentación en esa tienda** (otro ticket): se actualiza con el de este ticket, con `audit_log` `update`, y su id va en `price_ids`. Es la misma regla de `CBUpsertSourcedPrice` de api-compras-b.
5. **Granel**: la cantidad viene ×1000 (1.250 kg = 1250). Lo que entra a la existencia es `HALF_UP(PurchaseMovement(cantidad, contenido) / 1000)`; con piezas enteras es exacto (3 × 946 ml = 2838000).
6. **Compra a crédito del ticket**: valida lo mismo que `createCreditPurchase` (api-creditos), y además:
   - el capital tiene que ser igual al total del ticket (`invalid_sum`);
   - `purchase_date` tiene que ser la fecha del ticket;
   - `credit_purchase.expense` y `credit_purchase.expense_id` no se aceptan, porque el gasto es el del ticket;
   - sin `shares`, todo el capital va al titular del crédito, como en api-creditos.
   El armado de la entrada del motor está repetido en `proximo_pago.go` (`credInput.credit`), así que sigue en pie la PROPUESTA del §13 de exportarlo desde api-creditos.
7. **Estados que se pueden confirmar**: `pending`, `done` y `failed`. `failed` es el ticket que se captura a mano cuando el OCR falla (§17.2). Si se confirma uno en `pending`, el job ya no lo toca, porque `GastosSetTicketOCR` solo actualiza los `pending`.
8. Los errores de campo llevan la ruta JSON (`expense.items`, `credit_purchase.principal_cents`), como dice C1 §3.

**PREGUNTA** (para Axel):
- a) **Productos sin renglón de inventario** o en **ritmo de consumo**. C1 dice «productos con stock», y `TicketConfirmResult.inventory_item_ids` dice «creados o actualizados». Hoy solo se suma a los que ya tienen `inventory_item` en `stock`; a los demás no se les crea nada. Opciones:
  - (A) así;
  - (B) crear un `inventory_item` `stock` para cada producto comprado que no tenga uno;
  - (C) en ritmo de consumo, tomar la compra como un conteo nuevo (`counted_date` = fecha del ticket).
  Impacto: la lista de compras (S4) y la tarjeta de inventario.
- b) **Reparto de la compra a crédito**: las partes por persona llegan como las mande el cliente; no se sacan de la cuenta de cada renglón del gasto. ¿Se derivan de los renglones (Karen = sus renglones) o se quedan como captura aparte (hoy)?
