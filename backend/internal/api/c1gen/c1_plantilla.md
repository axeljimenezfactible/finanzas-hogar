# C1 · Contrato de datos y API · Finanzas Hogar

| | |
|---|---|
| Versión | **C1 v1.0** (10-oct-2026) · `info.version` 1.0.0 de `docs/contratos/openapi.yaml` |
| Autor | Agente S6a-contrato-c1 (flujo WF-0 «Cimientos», fase S6 parte a); corrección r1: agente correccion-C1-r1; corrección r2: agente correccion-C1-r2 |
| Estado | Propuesta para la puerta 1 (C1 + U1 + criterios) |
| Fuentes | `plan.md` §1, §3.4 (S6, S7), §3.5 (U1), §3.7; `decisiones.md`; esquema de S1, versión 1.1 (`backend/db/migrations/00001_esquema.sql`, `docs/contratos/diccionario-datos.md`, `docs/contratos/reglas-negocio.md`); `docs/design/U1-arquitectura.md` (199 necesidades N-xxx) |
| Archivo de máquina | `docs/contratos/openapi.yaml` (OpenAPI 3.1.1): {{N_RUTAS}} rutas, {{N_OPS}} operaciones, {{N_ESQ}} esquemas |
| Generador | Este documento, `openapi.yaml` y `backend/internal/api/routes.go` son archivos generados: no se editan a mano. Fuente: `backend/internal/api/c1gen/` (este texto sale de `c1_plantilla.md`; ver su `README.md`). Regenerar: `python3 -I backend/internal/api/c1gen/build.py`. Comprobar que nadie los editó a mano: `python3 -I backend/internal/api/c1gen/build.py --check`. Verificar C1-01 a C1-15, la matriz §10 y S6-03: `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py` |

Este documento fija **cómo** se piden los datos que U1 dice que cada pantalla **necesita**. El OpenAPI es la fuente de verdad de rutas y esquemas; aquí están las convenciones únicas, las métricas con su fórmula, la matriz pantalla → necesidad → endpoint y la cobertura de S1. Todo lo que no sale del plan, de las decisiones o del oráculo va marcado **PROVISIONAL** y apunta a su bloqueo en `docs/estado/bloqueos/`.

## Cambios

| Versión | Fecha | Cambio |
|---|---|---|
| C1 v1.0 | 10-oct-2026 | Primera versión (WF-0). Los cambios posteriores se publican como C1 v1.1, v1.2… (plan §3.7) con esta tabla actualizada y `info.version` igual. |
| C1 v1.0 · corrección r1 | 10-oct-2026 | Sin cambios de rutas, esquemas, convenciones ni matriz. El generador y el verificador pasan al repositorio (`backend/internal/api/c1gen/`) y las cabeceras de `openapi.yaml` y `routes.go` lo citan (hallazgo H-05 de la revisión r1). |
| C1 v1.0 · corrección r2 | 11-oct-2026 | Regenerado contra el esquema 1.1 de S1 (hallazgos R2-01 y R2-02 de la revisión r2; sigue siendo la entrega de WF-0 para la puerta 1). Lista de compras guardada: `generateShoppingList` (N-127), `createShoppingListItem` (N-152), `updateShoppingListItem` (N-151, N-153), `assignShoppingListStore` (N-165) y `getShoppingList` con la lista guardada y el filtro `is_bought` (N-150, N-073), todo PROVISIONAL por la opción A de `C1-lista-compras`. Presupuesto con vigencia (`valid_from_date`, `valid_to_date`; sin `period_id`) y 409 `budget_overlap`. Tiempos aptos de la receta como lista `meal_slot_ids` sobre `recipe_meal_slot` (sin `default_meal_slot_id`). Descripción del gasto opcional sin texto de relleno. Movimientos de inventario guardados en `inventory_movement`, con motivo en los ajustes. Ingreso esperado con `second_half_amount_cents`. `idempotency_key` en `x-tablas` de los POST de captura. Tope por quincena sin cifra. Columna nueva `category.kind`. |

## 1. Convenciones generales

1. **Base**: `https://finanzas.local:8443` (solo red local, decisión H). Recursos bajo `/api/v1`; fuera de él solo `GET /health` y `GET /openapi.yaml`.
2. **Formato**: JSON UTF-8 (`application/json`), salvo adjuntos y PDF (binario) y cargas de archivo (`multipart/form-data`).
3. **Nombres**: rutas, parámetros y propiedades en inglés corto, **iguales al `nombre_codigo` del diccionario de S1**. Descripciones y mensajes en español de México.
4. **Extensiones del contrato**: `x-tablas` (tablas de S1 que lee o escribe cada operación), `x-necesidades` (necesidades N-xxx de U1 que cubre), `x-tabla` (tabla de un esquema de entidad), `x-calculado` (campo calculado por la API o un motor; no es columna), `x-entrada` (dato de una acción que no se guarda, p. ej. la contraseña o la confirmación de una restauración), `x-columna` (columna a la que corresponde una entrada en lista, p. ej. `plan_dates` → `plan_date`), `x-app-setting` (clave de `app_setting`) y `x-provisional` (depende de un bloqueo abierto).
5. **Booleanos**: las columnas 0/1 de SQLite se exponen como `true`/`false`.
6. **Valores nulos**: en las respuestas, una columna NULL **se omite** (campo ausente = sin valor; la pantalla muestra «—», U1 §0.2 «sin datos no es cero»). Los esquemas de lectura marcan como `required` solo las columnas NOT NULL.
7. **Altas** (`POST`) aceptan las columnas escribibles; `id`, `created_at` y `updated_at` los pone el servidor. **Cambios** (`PATCH`) siguen JSON Merge Patch (RFC 7396): solo se cambian los campos enviados y `null` borra un valor que admite NULL. Reemplazos completos de listas (orden de categorías, destinatarios, grupos de producto) usan `PUT` con `{"ids": [...]}`.
8. **Idempotencia**: todo `POST` de captura exige el encabezado `Idempotency-Key` (16 a 128 caracteres `[A-Za-z0-9_-]`, UUID v4 recomendado). Repetir la clave con el mismo cuerpo devuelve la respuesta original con `Idempotent-Replayed: true`; con otro cuerpo, 409 `idempotency_key_reused`. La API guarda cada clave por cuenta en la tabla `idempotency_key` de S1, con la huella SHA-256 del cuerpo y la respuesta (S1 C-15), y la declara en `x-tablas` de cada operación que la exige; un job la purga por antigüedad. Retención de 24 h: **PROVISIONAL** (bloqueo `C1-parametros-tecnicos`). Lo exigen {{N_IDEM}} operaciones. No la llevan: las consultas por `POST` que no guardan nada (`previewCreditPurchase`, `compareScenarios`, `previewPlannedItemsImport`, `suggestRecipeWithAi`), las acciones de sesión (`login`, `logout`, `changeOwnPassword`, `revokeSessions`), el aviso de prueba (`sendTestNotification`) y la restauración, que lleva su propia confirmación escrita (`restoreBackup`).
9. **Identificador de petición**: toda respuesta lleva `X-Request-Id`; el mismo valor va en `Error.request_id` y en `audit_log.request_id` de los cambios que haga la petición.
10. **Cálculos**: la API no recalcula reglas en el cliente. Las cifras marcadas «motor Sx» las calcula ese motor con las reglas de `reglas-negocio.md`; la interfaz solo muestra, filtra, ordena y suma lo visible (U1 §0.2).

## 2. Autenticación, sesiones y permisos

- **Esquema**: `cookieAuth`, cookie `fh_session` **HttpOnly; Secure; SameSite=Strict; Path=/**. El token en claro solo viaja en la cookie; la base guarda su SHA-256 (`device_session.token_hash`). Ninguna respuesta incluye hashes, tokens, `chat_id` ni llaves (C1-12).
- **Operaciones públicas**: solo `POST /api/v1/auth/login` y `GET /health` (sin datos). Todo lo demás, incluidos adjuntos, PDF y `GET /openapi.yaml`, exige sesión; sin sesión, vencida o cerrada responde **401 `unauthorized`**.
- **Inicio de sesión** (`login`, N-010): usuario y contraseña (argon2id en la base). Error de credenciales: 401 `invalid_credentials` con el mensaje genérico «Usuario o contraseña incorrectos». **Límite (decisión A): 5 intentos fallidos en 15 minutos** (`app_setting auth.max_failed_attempts = 5`, `auth.lockout_window_minutes = 15`). El intento siguiente responde **429 `too_many_attempts`**, mensaje «Demasiados intentos. Vuelve a intentar a las 21:45.» y `retry_at` (instante con -06:00), **aunque la contraseña sea correcta**; así se cubre también el «estado de bloqueo» (N-011) sin una consulta pública que revele usuarios. **PROVISIONAL** (bloqueo existente `criterios-login-bloqueo`, opción c: S6 propone y Axel aprueba en la puerta 1): ventana deslizante de 15 minutos; se cuentan los fallos por usuario y por IP (lo que llegue primero a 5); `retry_at` = hora del fallo más antiguo de la ventana + 15 minutos.
- **Sesión larga por dispositivo** (decisión A, «meses»): **PROVISIONAL** 90 días desde el último uso, renovable con cada petición (`device_session.expires_at`, `app_setting session.max_age_days` sin valor sembrado; bloqueo `criterios-login-bloqueo`). `logout` cierra la del dispositivo; Ajustes lista las sesiones de **ambos** usuarios y cualquiera puede cerrarlas (`listSessions`, `revokeSession`, `revokeSessions`, N-271 y N-272). Una sesión cerrada responde 401 en la siguiente petición.
- **Permisos**: Axel y Ximena son co-titulares con **los mismos permisos**; ninguna operación filtra ni oculta datos según quién inició sesión (U1 §0.5). El usuario de la sesión solo da valores por defecto (persona del gasto, del peso) y el «creado por» de `audit_log`. **No existe operación para crear cuentas**: Karen no tiene cuenta ni acceso (decisión A); `person` no se crea desde la API.
- **CSRF**: con `SameSite=Strict` (**PROVISIONAL**, bloqueo `C1-parametros-tecnicos`) el navegador no envía la cookie en peticiones que vienen de otro sitio.
- **Contraseña mínima** del cambio de contraseña (N-273): **PROVISIONAL** 12 caracteres (bloqueo `C1-parametros-tecnicos`).

## 3. Modelo de errores (esquema único)

Todas las respuestas 4xx y 5xx de todas las operaciones usan `components/schemas/Error` mediante `components/responses/*` (más `default` para cualquier otro código):

```json
{"code": "validation_failed", "message": "Revisa los campos marcados.", "fields": [{"field": "amount_cents", "code": "out_of_range", "message": "El monto debe ser mayor que cero."}], "request_id": "01JABCDEF0123456789"}
```

- `code`: estable, en inglés, `^[a-z][a-z0-9_]*$`. `message`: español de México, lenguaje simple (qué pasó y qué hacer). `fields`: detalle por campo (nombre_codigo). `request_id`: siempre. `retry_at`: solo en 429.

| HTTP | Uso | Códigos estables |
|---|---|---|
| 400 | Petición mal formada: JSON inválido, parámetro con formato incorrecto (p. ej. `period=oct-2026`, fecha en formato de pantalla `10-oct-2026`) | `bad_request` |
| 401 | Sin sesión, vencida o cerrada; credenciales incorrectas en el inicio de sesión | `unauthorized`, `invalid_credentials` |
| 403 | Acción prohibida sobre el registro (precio de PROFECO) | `forbidden`, `read_only_source` |
| 404 | El registro o la ruta no existe (también una persona no externa en el estado de cuenta) | `not_found` |
| 405 | Método no admitido en una ruta existente | `method_not_allowed` |
| 409 | Conflicto con el estado: ya existe, en uso, gasto con compra a crédito, receta planeada, casilla con «fuera», vigencia de presupuesto traslapada, IA apagada, lectura asistida sin veredicto, clave de idempotencia reutilizada | `conflict`, `already_exists`, `in_use`, `has_credit_purchase`, `has_payments`, `recipe_planned`, `slot_conflict`, `budget_overlap`, `ai_disabled`, `assisted_reading_not_allowed`, `idempotency_key_reused` |
| 413 | Archivo demasiado grande | `payload_too_large` |
| 415 | Tipo de archivo no admitido | `unsupported_media_type` |
| 422 | Validación de negocio o de campos (con `fields`) | `validation_failed`, `invalid_current_password`, `not_revolving`, `not_a_card`, `external_person` |
| 429 | Límite de intentos de inicio de sesión (con `retry_at`) | `too_many_attempts` |
| 500 / 501 / 503 | Error interno / operación aún no implementada (esqueleto WF-0) / servicio o archivo no disponible | `internal_error`, `not_implemented`, `unavailable` |

**Validación de dinero**: un importe que no es entero (`973.75`, `"973.75"`) es 400 `bad_request`; un entero fuera de regla (cero o negativo donde se pide mayor que 0, partes que no suman el total) es 422 `validation_failed` con `fields[].code = out_of_range` o `invalid_sum`. **Validación de fechas**: una fecha que no es `AAAA-MM-DD` válida (`2026-02-30`, `10-oct-2026`) es 400; una fecha válida fuera de regla (pago antes de la compra, fin antes del inicio) es 422.

## 4. Fechas y dinero

- **Fechas de negocio** (`Date`, sufijo `_date`, y `valid_from`/`valid_to`): ISO 8601 `2026-10-10`, calendario de `America/Mexico_City`. **Instantes** (`DateTime`, sufijo `_at`): ISO 8601 con desplazamiento explícito de `America/Mexico_City`, hoy `-06:00` (`2026-10-10T21:30:00-06:00`); la base los guarda en UTC con «Z» y la API los convierte. **Mes** (`Month`): `2026-10`. Ningún formato de pantalla (`10-oct-2026`, `Vie 30 oct 2026`) aparece en la API; la interfaz los forma.
- **Quincenas**: los días de pago son el **15 y el 30**; en meses más cortos el «30» es el **último día** (en febrero, el 28: `2027-02-28`, R-23). Rangos para agrupar: Q1 del 1 al 15 y Q2 del 16 al último día — **PROVISIONAL** (bloqueo `U1-periodos-semana-quincena`, opción C, la misma que usan U1 y S1). **Semana**: lunes a domingo, **PROVISIONAL** (mismo bloqueo).
- **Dinero** (`MoneyCents`): entero exacto de **centavos MXN**, sufijo `_cents`, nunca `number` de punto flotante. `97375` = $973.75. Redondeo **ROUND_HALF_UP**; en una serie de pagos **el último pago ajusta centavos** (5,034.74 → 1,258.68 ×3 + 1,258.70, R-18 a R-21). Los dos modos de base de la serie que usa el oráculo y su generalización están en `S1-redondeo-series` (las series se guardan ya calculadas, el contrato solo las expone). En pantalla: `$`, separador de miles y dos decimales.
- **Porcentajes** (`BasisPoints`, sufijo `_bp`): 10000 = 100 %. Los porcentajes de avance calculados (`used_pct`, `progress_pct`) son enteros ROUND_HALF_UP. **Tasas**: fracción exacta `rate_num`/`rate_den` (S1; la tasa redondeada no reproduce el oráculo) y, para mostrar, `rate_pct` en texto con 4 decimales (`25.8685`). **Cantidades** (`Milli`, sufijo `_milli`): milésimas de la unidad base (g, ml, pieza) o de porción.

## 5. Paginación y orden (una sola convención)

- **Colección** = lista de registros de una entidad. Toda colección (`{{N_LISTAS}}` operaciones `GET`) usa **los mismos parámetros**: `cursor` (opaco; sin valor = primera página), `limit` (por defecto **50**, como pide U1; máximo **200**, PROVISIONAL `C1-parametros-tecnicos`) y `sort` (propiedades separadas por coma, `-` = descendente; el `id` cierra siempre el orden para que el cursor sea estable). **Misma envoltura** `Page`: `{"items": [...], "next_cursor": "…" | null}`; **`next_cursor: null` indica que no hay más**. El orden por defecto de cada colección está en su descripción (tomado de U1 §2.5.b).
- **Agregados y métricas** (tarjetas, resúmenes, series de 12 meses o 24 quincenas, comparativas) son objetos calculados con tamaño acotado; no se paginan.

## 6. Filtros globales (definición única)

Los seis filtros del plan (§3.7) tienen **una sola definición** en `components/parameters` y todas las operaciones los usan por `$ref`:

| Filtro | Parámetro | Tipo y formato | Significado |
|---|---|---|---|
| persona | `person_id` | ids separados por coma | `person.id`; en gastos, `expense.person_id`; en partes de dinero, la persona de la parte |
| familia | `family` | `household` · `external` · `joint` | U1 §0.2: Hogar (partes a nombre del titular) o Karen (partes de persona externa); `joint` = con ambas partes. Semántica incluyente **PROVISIONAL** (bloqueo `C1-filtro-familia`) |
| categoría | `category_id` | ids separados por coma | `category.id` |
| tienda | `store_id` | ids separados por coma | `store.id` (en promociones, el comercio ligado) |
| periodo | `period` | `2026-10` mes · `2026-10-q1`/`q2` quincena · `2026-W42` semana · `2026` año · `2026-10-10..2026-11-09` rango | mes calendario como base; quincena y semana PROVISIONALES (`U1-periodos-semana-quincena`) |
| crédito | `credit_id` | ids separados por coma | `credit.id`; en gastos, la forma de pago con esa línea |

Además se reutilizan `q` (texto) y los parámetros de paginación. Los filtros propios de una entidad (`status`, `kind`, `source`, etc.) se documentan en cada operación.

**Matriz recurso × filtro** (operaciones que admiten cada filtro):

{{FILTROS}}

## 7. Recursos: las 31 entidades de plan §3.4.1

Cada entidad tiene su recurso (lectura y, cuando el usuario la edita, escritura) o declara qué motor u operación la usa. Las vistas siguen U1 §2.5.a y se piden con `GET /api/v1/views/{entity}` (§8).

{{ENTIDADES}}

## 8. Vistas de base de datos (tabla, tablero, calendario, galería)

`GET /api/v1/views/{entity}?mode=table|board|calendar|gallery` (decisión G) aplica los filtros globales, `q`, `status`, `sort` y la paginación únicos, y responde `ViewPage` = `Page` + `mode` + `groups`. En **tablero**, `group_by` elige la columna y cada grupo trae `count`, `total_cents` (Σ del importe principal) y `share_bp`; en **calendario**, `date_field` elige la fecha y los grupos son días; en **galería**, cada registro trae `image_attachment_id` cuando hay foto. La vista elegida, filtros y columnas visibles se recuerdan en el dispositivo (U1 §2.5 y su supuesto 2); no hay vistas guardadas compartidas en esta versión.

{{VISTAS}}

## 9. Métricas del dashboard (fórmula, periodo, filtros, tablas, endpoint, ejemplo del oráculo)

Las 6 tarjetas de Inicio › Resumen (U1 N-020 a N-025, U3) y las métricas de sus subsecciones. Periodo por defecto: **mes calendario**; la quincena es subperiodo. «Parte del hogar» = parte a nombre del titular; toda cifra de presupuesto, flujo y ahorro usa la parte del hogar (U1 §0.2). Los ejemplos usan cifras del oráculo (`oraculo/model.json`, identificadores R-nn de `docs/criterios/S3.md`) o un marcador cuando el oráculo no trae el dato. En la API van en centavos.

| Id | Tarjeta o subsección | Fórmula | Periodo | Filtros | Tablas de origen | Endpoint | Ejemplo (oráculo) |
|---|---|---|---|---|---|---|---|
| M-01 | Mes vs plan (N-020) | real = Σ parte del hogar devengada en el mes (motor S2); plan = Σ `budget.amount_cents` del renglón vigente en el mes (RN-22) de las categorías activas; % usado = real ÷ plan × 100 entero ROUND_HALF_UP; restante = plan − real; estado normal < 80, aviso ≥ 80, excedido ≥ 100; días que faltan contando hoy | mes en curso (`period`) | `period` | budget, category, period, expense, expense_share, credit_purchase | GET `/api/v1/dashboard/month-vs-plan` | oct-2026: real (parte del hogar devengada) **1,197.64** = `periodos[0].devengado.axel` (P1 = oct-2026, hallazgo 4 de S1); plan: marcador (lo captura Axel) |
| M-02 | Flujo quincenal (N-021) | ingresos = quincena capturada (o esperada, marcada) + extras de la quincena; egresos hechos = efectivo pagado hasta hoy (contado + partes del hogar pagadas); comprometido = partes del hogar pendientes de la quincena + fijos y venideros sin pagar; saldo al cierre = ingresos − egresos − comprometido | quincena en curso | `period` (`2026-10-q2`) | income, planned_item, expense, expense_share, scheduled_payment, scheduled_payment_share, credit_payment | GET `/api/v1/dashboard/fortnight-flow` | Q2 oct-2026: comprometido incluye la parte del hogar del pago del 30-oct **598.82** (R-13); ingreso: marcador |
| M-03 | Próximo pago (N-022 a) | pago con alguna parte pendiente y fecha mínima ≥ hoy en todas las líneas; si hay pendientes con fecha < hoy, primero el más antiguo («sin marcar como pagado»); pagos del mismo día se suman | hoy en adelante | — | scheduled_payment, scheduled_payment_share, credit_payment, credit | GET `/api/v1/dashboard/next-payment` | **30-oct-2026, 973.75** (hogar 598.82, Karen 374.93) (R-13) |
| M-04 | Disponible A/B (N-022 b, N-201, N-285) | disponible de hoy de cada línea revolvente con la lectura activa (B por defecto) y aviso «pendiente confirmar»; vista previa con ambas lecturas y el mínimo del horizonte (motor S3) | hoy y horizonte | `credit_id` | credit, credit_purchase, scheduled_payment, credit_payment | GET `/api/v1/dashboard/next-payment`, GET `/api/v1/credits/{credit_id}/availability-preview` | 10-oct-2026: **A 10,000.00 · B 6,905.50** (R-64); mínimos **8,094.50 · 5,000.00** (R-61, R-62) |
| M-05 | Comidas de la semana (N-023) | costo = Σ costo de casillas con receta («fuera» no suma); kcal promedio por día = Σ kcal de la persona ÷ días con al menos un tiempo asignado; casillas sin armar = pares día × tiempo obligatorio sin receta ni «fuera» para alguna persona (máx. 21; PROVISIONAL `criterios-umbrales-ausentes`) | semana en curso (lunes a domingo, PROVISIONAL) | `period` (`2026-W42`) | meal_plan_entry, meal_slot, recipe, recipe_ingredient, nutrition_fact, price, presentation, health_profile | GET `/api/v1/dashboard/week-meals` | marcador (el oráculo no trae comidas) |
| M-06 | Mejor canasta (N-024) | lista de compras vigente (o la última compra con renglones) comparada por tienda (motor S5); ganadora = menor total; diferencia contra la segunda; cobertura = artículos con precio ÷ artículos; precios viejos usados; Sam's como si tuvieran membresía | lista vigente | — | shopping_list, shopping_list_item, meal_plan_entry, recipe_ingredient, inventory_item, price, presentation, product, store, expense, expense_item | GET `/api/v1/dashboard/best-basket` | compra semilla 10-oct-2026: Bodega Aurrera **3,094.50** (R-01) contra Walmart 3,233.00 (`data.json totales.walmart`), diferencia **138.50** (plan §3.4 fila S5) |
| M-07 | Ahorro (N-025) | ahorrado = Σ aportaciones − Σ retiros (fondo y metas); aportado del mes contra planificado del mes (Σ `planned_monthly_cents` activos); meta más próxima con % = ahorrado ÷ monto meta × 100 entero | mes en curso | — | savings_goal, savings_movement | GET `/api/v1/dashboard/savings` | marcador (el oráculo no trae ahorro) |
| M-08 | Participación de Karen (N-026) | su parte devengada de la despensa del mes; sus compras vivas; sus partes de pagos del mes con estado; deuda viva = Σ partes pendientes | mes en curso | — | person, expense, expense_share, credit_purchase, credit_purchase_share, scheduled_payment, scheduled_payment_share, credit_payment | GET `/api/v1/persons/{person_id}/participation` | oct-2026: parte devengada **749.86** (`periodos[0].devengado.karen`); su parte del pago del 30-oct **374.93** (R-11) |
| M-09 | Presupuesto por categoría (N-081, N-027) | por categoría: plan, real (parte del hogar devengada), plan − real, % usado y estado (umbrales de M-01) | mes (`period`) | `period` | budget, category, period, expense, expense_share | GET `/api/v1/budget/categories` | marcador (plan por categoría lo captura Axel) |
| M-10 | Flujo por quincena (N-090) | por quincena (24 hacia adelante): ingresos, partes del hogar de pagos, fijos y venideros, variable estimado (PROVISIONAL `criterios-promedios`), ahorro planificado, saldo y saldo acumulado; efectivo | desde Q1 del mes en curso | — | income, planned_item, expense, expense_share, scheduled_payment, scheduled_payment_share, credit_payment, savings_goal | GET `/api/v1/cash-flow/fortnights` | Q2 oct-2026: pagos de crédito, parte del hogar **598.82** (R-13) |
| M-11 | Devengado contra efectivo (N-091, N-028) | por mes: devengado y efectivo (totales y partes); diferencia acumulada = devengado acumulado − efectivo acumulado (no se rotula deuda viva; `criterios-invariante-deuda`) | mes en curso a +11 | — | expense, expense_share, planned_item, scheduled_payment, scheduled_payment_share, credit_purchase | GET `/api/v1/cash-flow/accrual-vs-cash` | oct-2026: devengado **1,947.50**, efectivo **973.75**, diferencia acumulada 973.75 (R-45); totales P1–P6 14,362.48 / 13,103.78 (R-51, R-52) |
| M-12 | Deuda viva a una fecha (N-092) | Σ de lo que queda por pagar después de la fecha (capital + intereses) de las compras hechas hasta esa fecha; total, por línea y por compra (motor S3) | fecha elegida | `date` | credit_purchase, scheduled_payment, credit | GET `/api/v1/debt/live` | 10-abr-2027: **6,293.44** (R-53 = 1,258.70 R-54 + 5,034.74 R-55) |
| M-13 | Calendario unificado de pagos (N-192, N-027) | pagos de todas las líneas en el rango con partes y estado; promociones vigentes y cortes de tarjeta | mes visible o hoy..+60 | `period`, `credit_id`, `family` | scheduled_payment, scheduled_payment_share, credit_payment, credit_purchase, promotion, credit_statement | GET `/api/v1/scheduled-payments`, GET `/api/v1/calendar/events` | pagos de la despensa 1: 30-oct, 15-nov, 30-nov y 15-dic-2026 por **973.75** (R-13 a R-16) |
| M-14 | Plan contra real por mes (N-085, N-028) | por mes (máx. 12): plan total y real total (o de una categoría) | 12 meses | `category_id` | budget, category, expense, expense_share | GET `/api/v1/budget/monthly-trend` | marcador (plan lo captura Axel) |
| M-15 | Costo semanal y kcal (N-129, N-233, N-028) | por semana (12): costo de casillas con receta y kcal promedio por día por persona contra su meta | 12 semanas | `person_id` | meal_plan_entry, recipe, recipe_ingredient, nutrition_fact, price, health_profile | GET `/api/v1/meal-plan/trends` | marcador |
| M-16 | Metas (N-111, N-027) | por meta: ahorrado, % de avance entero, aportación mensual necesaria = (monto meta − ahorrado) ÷ meses que faltan (motor S2) y estado | hoy | `status` | savings_goal, savings_movement | GET `/api/v1/savings-goals` | marcador |
| M-17 | Deuda viva proyectada (N-194, N-209) | por mes (12): deuda viva al cierre por línea; intereses pagados y por pagar | 12 meses | `credit_id` | credit, credit_purchase, scheduled_payment, credit_payment | GET `/api/v1/credits/debt-trend` | punto de control: deuda viva al 10-abr-2027 **6,293.44** (R-53); los meses sin cifra canon llevan la marca de `criterios-proyeccion-12m` |

## 10. Matriz pantalla P-xx → necesidad N-xxx → endpoint(s)

Una fila por necesidad y pantalla (las necesidades reutilizadas aparecen en cada pantalla que las usa, como en U1 §3.29). Estado: **OK** = cubierta; **PROVISIONAL** = cubierta en parte, con bloqueo; **BLOQUEO** = sin endpoint por un bloqueo abierto (hoy ninguna); **Local** = vive en el dispositivo (U1 la declara sin API).

{{MATRIZ}}

## 11. Cobertura

### 11.1 Resumen
{{RESUMEN}}

### 11.2 Operaciones sin necesidad de U1 (excepciones técnicas, una por una)
{{SIN_N}}

### 11.3 Necesidades sin endpoint
{{N_SIN_EP}}

### 11.4 Tablas de S1 → operaciones ({{N_TABLAS_OP}} de {{N_TABLAS}})
{{TABLAS}}

`idempotency_key` no tiene pantalla propia: es infraestructura de S6. La leen y escriben en la misma transacción todas las operaciones que exigen `Idempotency-Key` (§1 punto 8; S1 C-15), por eso figura en su `x-tablas`, y la purga por antigüedad la hace un job de S6 (retención PROVISIONAL, `C1-parametros-tecnicos`). `inventory_movement`, `recipe_meal_slot`, `shopping_list` y `shopping_list_item` son detalle de entidades del plan (inventario, receta y lista de compras) y tienen sus propias operaciones.

### 11.5 Las 14 secciones de Ajustes → necesidades → operaciones
{{AJUSTES}}

### 11.6 Capacidades exigidas (C1-11) → operaciones
{{CAPACIDADES}}

## 12. Índice de operaciones

{{OPERACIONES}}

## 13. PROVISIONAL y bloqueos

| Tema | Qué queda PROVISIONAL en C1 | Bloqueo |
|---|---|---|
| Filtro familia | valores `household`/`external`/`joint` y semántica incluyente | `C1-filtro-familia` (C1) |
| Lista de compras | que las marcas de comprado, los artículos a mano y la tienda o cantidad elegidas se guarden en el servidor y se compartan entre los dos teléfonos (opción A): `getShoppingList` (N-150, N-073), `generateShoppingList` (N-127), `createShoppingListItem` (N-152), `updateShoppingListItem` (N-151, N-153) y `assignShoppingListStore` (N-165). Si Axel elige B o C en la puerta 1, esas operaciones cambian o salen | `C1-lista-compras` (C1), igual a `S1-lista-compras` |
| Huecos U1 ↔ S1 que siguen abiertos | qué vigencia escribe un cambio de presupuesto (solo ese mes, o ese y los siguientes; N-087, N-278), si «Guardar para después» guarda las correcciones del OCR (N-075) y cómo se registra el chat con el bot (N-299). Los otros cuatro huecos y la idempotencia ya los resolvió el esquema 1.1 de S1 (§14) | `C1-huecos-u1-s1` (C1), igual a `S1-huecos-negocio` |
| Parámetros técnicos | `limit` máximo 200, retención de `Idempotency-Key` (24 h; la tabla ya existe), `SameSite=Strict`, contraseña mínima, tamaño máximo de la carga de PROFECO | `C1-parametros-tecnicos` (C1) |
| Campos calculados de motores | clasificación del IMC, precio por unidad base, punto de equilibrio de la membresía, rendimiento por peso y puntaje del comparador | `C1-campos-calculados` (C1) |
| Inicio de sesión y sesión | conteo por usuario e IP, ventana deslizante, duración 90 días renovable | `criterios-login-bloqueo` (existente) |
| Quincena y semana | Q1 1–15, Q2 16–fin; semana lunes a domingo | `U1-periodos-semana-quincena` (existente) |
| Adjuntos y tiempos | tamaño máximo 10 MB de adjuntos; hora de resúmenes y anticipación de alertas; tiempo máximo del comparador | `criterios-umbrales-ausentes` (existente) |
| Frescura | frontera ≥ 30 días y peso de precio viejo | `criterios-frescura-peso`, `S1-frescura-peso` (existentes) |
| Promedios y variable estimado | definición de promedios y del gasto variable estimado | `criterios-promedios` (existente) |
| Comparador | restricciones sin valor; normalización, rendimiento por peso, puntaje | `S1-restricciones-optimizador`, `criterios-optimizador` (existentes) |
| Pago por parte | capital que repone la parte de Karen pagada por separado | `S1-reposicion-parte-externa` (existente) |
| Meta de kcal | escala de actividad | `S1-meta-kcal` (existente) |
| Meses sin cifra canon | marca de los meses fuera de P1–P6 | `criterios-proyeccion-12m` (existente) |

## 14. Inconsistencias S1 ↔ U1 encontradas (sin tocar sus archivos)

En C1 v1.0 (WF-0) C1 tomó en cada caso la lectura más conservadora y la marcó PROVISIONAL, y pidió las correcciones a S1 en `docs/estado/wf0/S6a-contrato-c1.md`. El esquema 1.1 de S1 dio dónde guardar siete de los diez puntos: cinco quedan resueltos, y la lista de compras y la vigencia que escribe un cambio de presupuesto esperan a Axel. La corrección r2 regeneró C1 contra ese esquema. La columna «Estado» dice qué queda abierto.

{{INCONSISTENCIAS}}
