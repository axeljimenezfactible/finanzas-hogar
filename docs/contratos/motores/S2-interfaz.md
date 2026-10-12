# S2 · Interfaz del motor financiero

Versión 1.2 · 11-oct-2026 · WF-1 (1.0: agente `contrato-S2`; 1.1: agente `correccion-contrato-S2-r1`; 1.2: agente `correccion-contrato-S2-r2`) · Estado: PROPUESTA para el revisor de motores, que la evalúa contra `docs/criterios/S2.md` 1.2.

**Para quién.** Dos agentes que implementan en paralelo y sin ver el trabajo del otro: Go en `backend/internal/motor_fin` (paquete `motor_fin`, módulo `finanzas-hogar`) y Python en `oraculo/fin_*.py` (verificador independiente, S2-03). Los dos leen este archivo, `oraculo/fixtures/FORMATO.md` 1.2, `docs/criterios/S2.md` y las reglas que se citan. Aquí no hay código: tipos, firmas, contratos JSON y tablas. Piezas compartidas (FORMATO §4.1 y §5): el agente Go escribe además `backend/cmd/motores/main.go` (binario común con el registro de motores) y `cmd/motores/s2.go`; el agente Python escribe `oraculo/compara.py` (genérico, con las ramas de todos los motores). **Autoría de las fixtures (§5.0):** el agente Python de S2 escribe todos los archivos de `oraculo/fixtures/s2/` (las 27 fixtures y los 3 de `_datos/`); el agente Go no escribe en `oraculo/`. **`oraculo/verify3.py`** (solo la ruta de salida de su línea 89, S2-04) también queda en el alcance del agente Python de S2 (§0.5).

**Fuentes.** `docs/plan/decisiones.md` §B y §K (K2, K8, K14, K19, K20, K21, K39, K44); `docs/contratos/reglas-negocio.md` 1.3 (C-01…C-16, RN-01…RN-07, RN-18, RN-21, RN-22, RN-27, RN-28, RN-29); `docs/contratos/C1-datos.md` §9 (M-01…M-17) y los extractos `presupuesto-a.yaml` (`updateBudget`, parámetro `BudgetScope`) e `inicio.yaml`; `docs/plan/plan.md` §1.2 y §3.4 fila S2; `oraculo/model.json` (596 hojas numéricas, sha256 `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`) y `oraculo/model.py` (`sched`, `mk`). PROVISIONAL marca todo lo que depende de un bloqueo abierto (§8).

**Cambios de la 1.1** (hallazgos medios M1–M10 de `docs/estado/wf1/REVISION-contrato-S2-r1.md`; ninguna cifra canon cambia): M1, ids de consulta = ruta de `model.json` y vistas de cruce `vista_sensibilidad` y `vista_interes_por_monto` (§3.14, §5; FORMATO §2.1); M2, dueños de las piezas compartidas (arriba y FORMATO §4.1, §5, §6); M3, esquema de `_datos/insumo-s3.json` (§5.1); M4, partidas planificadas en la proyección con fórmula y marca `provisional` (§3.2, §3.12); M5, campos del flujo quincenal con fórmula y caso del 31-oct (§3.11, §5.2); M6, flujo de 16 quincenas y contraste con el valor estable (§3.11, §5); M7, regla de reparto efectiva ya resuelta (§2.7); M8, contenido y orden de cada lista (§3.2–§3.13); M9, fixture `s2-interes-por-monto` (§5); M10, `person` en ingresos y partidas, importe de un gasto a crédito (§2.4, §2.6, §2.7).

**Cambios de la 1.2** (hallazgos medios M1–M6 de `docs/estado/wf1/REVISION-contrato-S2-r2.md`; ninguna cifra canon cambia): M1, regla general de argumentos (§3) y una línea «Argumentos» por función con claves en inglés y obligatoriedad (§3.1–§3.14), más la función Python de `plan_contra_real` (§3.10) y la firma de `valor_estable` con `persons` (§3.6); M2, regla de contenido de los mapas `{code: cents}` (§1.1) y validación de personas (§1.4); M3, fórmula de cada campo de `presupuesto_mes` (§3.8) y de `month_summary` y de las metas sin monto o sin fecha (§3.13), con el bloqueo nuevo `S2-metas-resumen` (§8); M4, `projected` en `Expense` (§2.4); M5, autoría de `oraculo/fixtures/s2/`, trabajo del agente Go mientras falten y contenido completo de la semilla (§5.0, §5.3, §7); M6, corrección de `oraculo/verify3.py` línea 89 asignada al agente Python de S2 (encabezado, §0.5, §7).

## 0. Reglas de juego
1. **Puro y determinista.** Mismas entradas, mismas salidas; sin reloj, red, base de datos ni estado global. «Hoy» llega como parámetro (§2.1). La API (S6) arma las entradas desde SQLite y desde S3 (`motor_cred`); el motor no lee tablas.
2. **Frontera con S3.** S2 no calcula series de crédito: el calendario de pagos, sus partes por persona y el capital de cada pago llegan hechos por S3 (C-08, K14). S2 tampoco reparte gastos registrados: recibe las partes congeladas (`expense_share`, `credit_purchase_share`). Solo reparte (RN-03) las ocurrencias de partidas planificadas, que no guardan partes.
3. **Dinero en centavos enteros** (`int64` en Go, `int` en Python); ningún `float32`/`float64` en dinero (S2-13). Intermedios exactos: Go con enteros o `math/big.Rat`; Python con `fractions.Fraction` o `decimal.Decimal`.
4. **Independencia (S2-03).** El oráculo Python no lee salidas del motor Go ni comparte código con él. Método distinto recomendado: Go agrupa por periodo con aritmética entera; Python recorre el calendario día por día (`datetime.date`) y acumula con `Fraction`. Los dos redondean solo donde lo dice §1.3.
5. **Nada se escribe fuera de lugar (S2-04).** Ninguna ejecución modifica `oraculo/model.json` ni `oraculo/verify3.json` ni escribe fuera de `/home/deepsec/finanzas-hogar/`; temporales en `datos/tmp/`. **Corrección que exige S2-04, asignada al agente Python de S2:** hoy `oraculo/verify3.py` línea 89 escribe en `/home/deepsec/reporte-compra-latex/model/verify3.json`. Cambia solo esa ruta de salida a `datos/tmp/verify3.json`, armada desde la raíz del repositorio (`Path(__file__).resolve().parents[1]`) y creando la carpeta si falta; no toca el cálculo ni `oraculo/verify3.json`. Verifica: `grep -n reporte-compra-latex oraculo/*.py` → 0 líneas; `python3 -I oraculo/verify3.py` y luego `cmp datos/tmp/verify3.json oraculo/verify3.json` sin diferencias (sha256 `d86087f3797e08ceabd2c721bf423115345b318ae7fe5fa13b3400cd87b95041`; comprobado el 11-oct-2026 con una copia redirigida: iguales byte a byte); `sha256sum oraculo/model.json` sin cambio; borra el temporal.

## 1. Convenciones
### 1.1 Valores
| Clase | JSON | Go | Python | Ejemplo |
|---|---|---|---|---|
| Dinero | entero de centavos; clave con sufijo `_cents` (en las vistas del oráculo, §3.14, la clave de `model.json` sin sufijo) | `Cents` (`int64`) | `int` | `97375` = $973.75 |
| Fecha | `"AAAA-MM-DD"` | `Date` (año, mes, día; sin hora) | `datetime.date` | `"2027-02-28"` |
| Instante | RFC 3339 en UTC con «Z» (solo en entradas, K12) | `time.Time` | `datetime` con zona | `"2026-10-31T05:59:59Z"` |
| Mes | `"AAAA-MM"` | `Month` | `(año, mes)` | `"2026-10"` |
| Quincena | `"AAAA-MM-q1"` o `"AAAA-MM-q2"` (como `period` en C1) | `Fortnight` | `(año, mes, mitad)` | `"2026-10-q2"` |
| Semana | `"AAAA-Www"` (ISO 8601) | `Week` | `(año_iso, semana)` | `"2026-W42"` |
| Porcentaje | entero `_pct` donde C1 lo pide; décimas `_tenths_pct` donde lo pide RN-18 | `int64` | `int` | `88` |
| Conteo u ordinal | entero | `int` | `int` | `n_pagos` = 2 |
| Persona | `code` de `person` | `string` | `str` | `"axel"`, `"karen"` |

**Mapas `{code: cents}` de las salidas** (todo `by_person`, los `shares` de `cash.payments[]` en §3.2, `credit_payments.by_person` en §3.11 y los `by_person` de §3.6 y §3.12): llevan **siempre** al titular y a cada persona `external` de `persons`, en el orden de `persons`, con 0 si no tiene parte; **nunca** a la co-titular, que no tiene partes (§1.4). Los `by_category` llevan solo las categorías con importe distinto de 0 (activas o no), en orden de `code` (orden de bytes); quien lee una categoría ausente toma 0 (`vista_periodo`, `promedios` con `accrued_groceries`). Ejemplos con la semilla (§5.3): el componente `papel_contado` de `valor_estable` da `by_person` `{axel: 0, karen: 9950}`; cada pago de `s2-interes-por-monto` (compras solo del titular) da `shares` `{axel: …, karen: 0}`; `periodos {months: ["2026-10"]}` da `accrued.by_category` `{despensa: 194750}`, sin `papel`. Así `esperado` puede pedir al titular o a cualquier persona externa sin que el recorte (FORMATO §4, reglas 1 y 5) dé `__falta__`, y nunca pide a la co-titular ni una categoría en 0.

### 1.2 Calendario (RN-06, K8, C-05)
- Mes calendario: del día 1 al último. **Quincena 1**: días 1 a 15, paga el 15. **Quincena 2**: del 16 al último día, paga el 30 o el último día si el mes tiene menos de 30 (feb-2027 paga el 2027-02-28). El día 31 pertenece a la quincena 2; el día de pago no es el fin del periodo (S2-06).
- Semana de lunes a domingo (ISO). Días de un periodo = días reales del calendario (31, 30, 29 o 28).
- «Hoy» = fecha civil en `America/Mexico_City` del instante recibido, con la base de zonas del sistema (`time.LoadLocation`, `zoneinfo.ZoneInfo`), nunca un desfase fijo (C-05, S2-14).
- Etiqueta de un mes: `1-oct-2026 a 31-oct-2026` (meses `ene feb mar abr may jun jul ago sep oct nov dic`). Las etiquetas `10-oct-2026 a 9-nov-2026` de R-45…R-50 son ventanas del oráculo, no periodos del motor.

### 1.3 Redondeo (RN-01, RN-02, RN-07, K14, K39)
- `HALF_UP(a/b)`: cociente exacto redondeado a centavo; la mitad se aleja del cero, como `decimal.ROUND_HALF_UP` (con negativos, `HALF_UP(−x) = −HALF_UP(x)`).
- **Devengo de un costo C en n meses (RN-02):** meses 1…n−1 = HALF_UP(C/n); el último = C − (n−1)·HALF_UP(C/n). Por persona: cada persona que no es el titular lleva HALF_UP(su total/n) en los meses 1…n−1 y su resto en el último; el titular toma, cada mes, la parte del mes menos las de los demás. Ejemplo: despensa 1, C = 389500, n = 2 → 194750 y 194750; Karen HALF_UP(149972/2) = 74986; Axel 194750 − 74986 = 119764 (`periodos[0].devengado`).
- **Promedios (K20, K39, RN-27):** diario = HALF_UP(T/d); semanal = HALF_UP(7·T/d), un solo redondeo (no 7 × diario redondeado); mensual de un horizonte de k meses = HALF_UP(T/k).
- **Reparto de una ocurrencia planificada (RN-03):** parte por porcentaje = HALF_UP(monto·pct_bp/10000); parte fija = `fixed_cents`; la persona del resto absorbe la diferencia.
- **Valor estable:** cada partida aporta HALF_UP(total/ciclo) al mes (personas como en el devengo); la quincena es HALF_UP(x/2) por renglón (total y cada persona) por separado; la diferencia de centavo entre el total y la suma de personas se informa, no se corrige (R-40…R-42).
- Fuera de eso no hay redondeo: sumas y restas de centavos son exactas. Las series de pagos (base truncada, último ajusta, RN-07) son de S3 y llegan hechas.

### 1.4 Personas y «parte del hogar»
- El titular (`participation = holder`) recibe la parte del hogar Axel–Ximena (C-09); la co-titular no tiene partes. En las vistas del oráculo cada persona aparece con su `code` (`axel`, `karen` en la semilla).
- Una persona externa con `pays_share_directly = true` (Karen) tiene su parte de cada pago programado informada pero fuera del flujo del titular (RN-04). Sí entra en el efectivo total de un periodo, como en `periodos[k].efectivo.total`.
- Presupuesto y ahorro usan la **parte del hogar** (la del titular; C1-datos §9). El flujo quincenal usa el **efectivo del hogar**: la parte del titular en los pagos programados y lo que sale de su bolsillo en gastos de contado (§3.11). Periodos, horizonte y proyección dan además el total y cada persona.
- Validación de personas (C-09: el dinero del hogar Axel–Ximena va a nombre del titular): todo `person` o `paid_by` de un campo de dinero (`shares` de gastos, compras y pagos; `paid`; `embedded`; `share_rule`; `stable_items`) es el titular o una persona `external` de `persons`; un `code` que no está en `persons` o que es la co-titular es `invalid_argument`. Los ingresos tienen su regla (`income_not_holder`, §2.6). El `person` de una partida planificada de gasto es atribución (quién comió), no dinero, y puede ser cualquiera de `persons` (§2.7).

## 2. Tipos de entrada
La raíz es `Input` (en las fixtures, `entradas.datos`; ver §5 y `FORMATO.md`). Toda lista puede venir vacía. Un campo desconocido es error `unknown_field` en las dos implementaciones. Columna «Origen»: de dónde lo toma la API.

### 2.1 Parámetros · `Params` (`params`)
| Campo | Tipo | Oblig. | Significado | Origen |
|---|---|---|---|---|
| `today` | fecha | uno de `today`/`now` | Fecha de hoy ya resuelta | API |
| `now` | instante | uno de `today`/`now` | El motor obtiene `today` en `zone` (S2-14) | API |
| `zone` | texto | no | Solo `"America/Mexico_City"` (por omisión); otro valor es `invalid_zone` | decisión B |
| `holder` | persona | sí | Código del titular | `person` |
| `groceries_category` | texto | sí | Categoría que hace de «despensa» (ordinal `purchase_seq`, vistas del oráculo) | `category.code` (semilla de S1/S6) |
| `paper_category` | texto | sí | Categoría del papel (rama `papel`, `papel_contado`) | `category.code` |
| `repeat_cadence_enabled` | bool | no (`false`) | Interruptor de los meses 7–12 de la proyección | `app_setting projection.repeat_cadence_enabled` (K21, 00002) |
| `budget_warning_pct` | entero | no (80) | Umbral de aviso del presupuesto | `app_setting alerts.budget_warning_pct` (RN-22, K44, 00003) |
| `budget_limit_pct` | entero | no (100) | Umbral de excedido | `app_setting alerts.budget_limit_pct` |

### 2.2 Personas · `Person` (`persons`)
`code` (texto), `display_name` (texto), `participation` (`holder` · `co_holder` · `external`), `pays_share_directly` (bool). Exactamente un `holder`. Origen: `person`.

### 2.3 Categorías · `Category` (`categories`)
`code`, `name`, `kind` (`fixed` · `variable` · null), `is_active` (bool). Origen: `category`. El presupuesto solo cuenta categorías activas (M-01).

### 2.4 Gastos registrados · `Expense` (`expenses`)
Gastos fijos, variables y venideros ya capturados (`expense` + `expense_share`).
| Campo | Tipo | Oblig. | Significado |
|---|---|---|---|
| `id` | entero | sí | Identificador |
| `kind` | `fixed` · `variable` · `upcoming` | sí | Gasto fijo, variable o venidero |
| `category` | texto | sí | Código de categoría |
| `amount_cents` | dinero | sí | Importe. Con `credit` es el capital de la compra: debe ser igual a `principal_cents` de `credit_purchase_id` (si no, `sum_mismatch`), como la despensa 1: importe del gasto 3,094.50 = capital R-01 (RN-03, ejemplo) |
| `expense_date` | fecha | sí | Fecha de compra o consumo; inicia el devengo |
| `accrual_months` | entero 1–24 | no (1) | Meses de devengo (despensa bimestral = 2) |
| `payment_method` | `cash` · `debit` · `transfer` · `credit` | sí | Medio de pago |
| `paid_by` | persona o null | no | De qué bolsillo sale (null = titular) |
| `shares` | lista `{person, amount_cents}` | sí | Partes congeladas; suman `amount_cents` (si no, `shares_mismatch`) |
| `credit_purchase_id` | entero o null | con `credit` | Compra de S3 que lo financia (si falta, `missing_purchase`) |
| `planned_item_id` | entero o null | no | Partida planificada que cumple |
| `projected` | bool | no (`false`) | Solo `true` dentro de `cadence` (§2.10), como en `CreditPurchase` (§2.5); los gastos registrados lo omiten o lo traen en `false` |

Reglas (RN-02): de contado, débito o transferencia, el efectivo sale en `expense_date` del bolsillo de `paid_by`; el devengo es `amount_cents` en `accrual_months` meses desde el mes de `expense_date`, por persona según `shares`. Con `credit`, **la compra manda**: el devengo es su `total_cents` (capital + intereses) repartido con las partes de la compra (`shares[].total_cents`), en `accrual_months` meses desde el mes de `expense_date`; el efectivo son sus pagos programados. El gasto aporta fecha, categoría y meses de devengo.

### 2.5 Compras a crédito con su calendario (de S3) · `CreditPurchase` (`credit_purchases`)
| Campo | Tipo | Significado |
|---|---|---|
| `id`, `credit_id` | entero | Compra y línea |
| `expense_id` | entero o null | Gasto que la originó; si es null, la compra se devenga sola con `category` y `accrual_months` propios |
| `category` | texto o null | Solo cuando `expense_id` es null |
| `purchase_date` | fecha | Fecha de la compra |
| `principal_cents`, `interest_cents`, `total_cents` | dinero | Capital, intereses, total (= capital + intereses) |
| `accrual_months` | entero 1–24 | Solo cuando `expense_id` es null (por omisión 1) |
| `status` | `active` · `paid_off` · `cancelled` | Las canceladas se ignoran |
| `shares` | lista `{person, principal_cents, interest_cents, total_cents}` | Partes de la compra (`credit_purchase_share`); suman el total |
| `embedded` | lista `{category, person, amount_cents}` | Artículos marcados dentro de la compra (papel dentro de la despensa; `despensas[k].<persona>.papel_dentro`) |
| `payments` | lista `ScheduledPayment` | Calendario de S3 |
| `projected` | bool (`false`) | Solo `true` dentro de `cadence` (§2.10) |

`ScheduledPayment`: `seq`, `due_date`, `amount_cents`, `principal_cents`, `interest_cents`, `shares` (lista `{person, amount_cents}`), `paid` (lista `{person, amount_cents, paid_date}`; abonos registrados en `credit_payment`). S2 **verifica y no corrige** (S2-05): monto = capital + intereses, Σ partes = monto, Σ montos = total, Σ capital = capital de la compra; si algo falla, `sum_mismatch` con la ruta. **Ordinal `purchase_seq`:** posición 1, 2, … de la compra entre las de `groceries_category` ordenadas por (`purchase_date`, `id`); es la «despensa n» del oráculo (null en otras categorías).

### 2.6 Ingresos · `Income` (`incomes`)
`id`, `person` (código; obligatorio, `income.person_id`), `kind` (`fixed` = quincena, `extra`), `amount_cents` (> 0), `income_date`, `planned_item_id` (o null). Solo del titular (RN-05, decisión B): `person` distinto de `params.holder` es `income_not_holder`. Ximena no aporta ingresos.

### 2.7 Partidas planificadas · `PlannedItem` (`planned_items`)
Ingresos y gastos fijos o venideros del CSV de Axel (`planned_item`).
| Campo | Significado |
|---|---|
| `id`, `kind` (`income` · `fixed_expense` · `upcoming_expense`), `description` | Partida |
| `amount_cents`, `second_half_amount_cents` | Monto por ocurrencia; el de la quincena del 30 si difiere (solo `fortnightly`) |
| `frequency` (`weekly` · `fortnightly` · `monthly` · `yearly` · `once`), `interval_count` (1–24) | Periodicidad |
| `start_date`, `end_date` (o null) | Primera ocurrencia (su día es el ancla) y última fecha posible, inclusive |
| `category`, `payment_method`, `accrual_months` | Solo gastos |
| `person` | Código o null (`planned_item.person_id`, «a quien se atribuye»; null = hogar en general). En ingresos, solo el titular o null (otro es `income_not_holder`). En gastos se informa y S2 no la usa: las partes salen de `share_rule` y el efectivo de una ocurrencia sale del bolsillo del titular, como un `Expense` con `paid_by` null (§2.4), porque `planned_item` no tiene «quién paga» |
| `share_rule` | Regla **efectiva, ya resuelta por quien llama** con la prioridad de RN-03: la de la partida (`cost_share_rule_id`); si es NULL, la de su categoría (`category.default_cost_share_rule_id`); si tampoco, la regla por defecto (`cost_share_rule.is_default`, «Hogar», 100 % al titular). Forma `{parts: [{person, pct_bp} o {person, fixed_cents}], rest_person}`. Obligatoria en gastos (null es `invalid_argument`, igual que una parte fija mayor que el monto, RN-03); en ingresos va null. Ejemplo: una partida de Despensa sin regla propia llega con «Conjunto 50/50» de su categoría (`{parts: [{person: karen, pct_bp: 5000}], rest_person: axel}`), nunca como 100 % al titular |
| `credit_id` | Si no es null, la partida es de crédito: S2 la ignora (sus compras proyectadas llegan de S3 en `cadence`) |
| `is_active` | Las inactivas se ignoran |

**Ocurrencias** (comentarios de `planned_item` en 00001): `monthly` cada `interval_count` meses en el día de `start_date`, recortado al último día en meses cortos; `fortnightly` en los días de pago (15 y 30 o último), cada `interval_count` quincenas, con `second_half_amount_cents` en la del 30; `weekly` cada 7·`interval_count` días; `yearly` mismo mes y día, recortado; `once` solo `start_date`. Nunca después de `end_date`.

### 2.8 Presupuesto · `BudgetRow` (`budgets`) y `BudgetChange`
`BudgetRow`: `id`, `category` (código o null = renglón del total), `amount_cents`, `valid_from` (día 1), `valid_to` (último día de un mes o null = sin fin). Sin traslapes por categoría (`budget_overlap`).
`BudgetChange` (argumento de §3.9): `row_id`, `month`, `amount_cents`, `scope` (`onward` por omisión · `this_month`), con los nombres del parámetro `BudgetScope` de C1 (K2).

### 2.9 Metas y fondos · `SavingsGoal` (`savings_goals`) y `SavingsMovement` (`savings_movements`)
`SavingsGoal`: `id`, `name`, `kind` (`fund` · `goal`), `target_cents` y `target_date` (obligatorios en `goal`), `planned_monthly_cents` (o null), `start_date`, `is_active`. `SavingsMovement`: `goal_id`, `movement_date`, `kind` (`deposit` · `withdrawal`), `amount_cents` (> 0).

### 2.10 Cadencia proyectada · `Cadence` (`cadence`) — PROVISIONAL (`docs/estado/bloqueos/S2-proyeccion-interruptor.md`)
`{purchases: [CreditPurchase con projected = true], expenses: [Expense con projected = true]}`. Las arma quien llama (S3 calcula las series de las despensas proyectadas); S2 solo las usa si `repeat_cadence_enabled` y solo en los meses 7–12 de la proyección. Una fecha anterior al mes 7 es `cadence_before_month_7`. Monto, fechas y pagador del papel alterno quedan en el bloqueo.

### 2.11 Partidas del valor estable · `StableItem` (`stable_items`)
`{label, cycle_months, total_cents, shares: [{person, total_cents}]}`. En la semilla: `{label: "despensa", cycle_months: 2, total_cents: 503474, shares: [karen 149972, axel 353502]}` (`despensas[1]`, R-17) y `{label: "papel_contado", cycle_months: 2, total_cents: 19900, shares: [karen 19900]}` (`parametros.papel`, `papel[2]`). C1 v1.2 no tiene endpoint que consuma el valor estable (no es una métrica M-xx): se expone para la regresión (R-37…R-44) y para los reportes.

### 2.12 Ventanas declaradas · `Window` (argumento de §3.2)
`{label, start, end}` (inclusive). Solo para fixtures con eventos entre los días 1 y 9 de un mes (S2-09): el esperado se calcula con esas ventanas, nunca con R-45…R-50.

## 3. Funciones y salidas
Cada función tiene firma Go (paquete `motor_fin`), función Python (módulo `oraculo/fin_*.py`) y nombre JSON (`funcion` de una consulta en las fixtures, §5). Toda función devuelve `(resultado, error)` en Go y lanza `ErrorMotor(code, detalle)` en Python; los códigos están en §3.15. Las salidas JSON usan claves en inglés con `_cents` para dinero, salvo las vistas del oráculo (§3.14).

**Argumentos de una consulta** (`argumentos` en las fixtures, FORMATO §2): un objeto JSON con las claves en inglés que declara la línea «Argumentos» de cada función (`{}` si no lleva). Una clave desconocida es `unknown_field`; falta una obligatoria, llegan las dos alternativas de «exactamente uno» (o ninguna), o un valor tiene tipo o rango inválido: `invalid_argument`. Los datos salen siempre de `entradas.datos` (`Input`); los argumentos solo dicen qué calcular. El evaluador de `cmd/motores/s2.go` los traduce a la firma Go y `fin_run_fixtures.py` a la función Python.

### 3.1 Calendario · `calendario` (S2-06, S2-14)
- Go: `Today(now time.Time, zone string) (Date, error)`; `MonthBounds(m Month) (start, end Date)`; `FortnightOf(d Date) Fortnight`; `FortnightBounds(f Fortnight) (start, end, pay Date)`; `WeekOf(d Date) Week`; `WeekBounds(w Week) (monday, sunday Date)`.
- Python `fin_calendario.py`: `hoy`, `limites_mes`, `quincena_de`, `limites_quincena`, `semana_de`, `limites_semana`.
- Argumentos (exactamente uno, obligatorio) `{date}` (fecha) o `{now}` (instante; la zona sale de `params.zone`) → `{today, month, month_start, month_end, month_label, days_in_month, fortnight, fortnight_start, fortnight_end, pay_date, week, week_start, week_end}`.
- Ejemplos: 2027-02-20 → `2027-02-q2`, 2027-02-16 a 2027-02-28, paga 2027-02-28; 2026-10-31 → `2026-10-q2`, 2026-10-16 a 2026-10-31, paga 2026-10-30; 2026-10-14 → `2026-W42`, 2026-10-12 a 2026-10-18; `now` 2026-10-31T05:59:59Z → today 2026-10-30, y 2026-10-31T06:00:00Z → 2026-10-31 (UTC−6 todo el año).

### 3.2 Periodos: devengado y efectivo · `periodos` (S2-02, S2-06, S2-09, M-08, M-11)
- Go: `Periods(in Input, spec PeriodSpec) ([]PeriodSummary, error)`, con `PeriodSpec{Months []Month; Windows []Window}` (uno de los dos). Python `fin_devengo.py`: `periodos(datos, meses=None, ventanas=None)`.
- Argumentos (exactamente uno, lista no vacía) `{months: [Month]}` o `{windows: [Window]}` → lista de resúmenes de periodo, uno por elemento (abajo).
- Salida por periodo:
  - `label`, `start`, `end`, `days`;
  - `accrued`: `{total_cents, by_person: {code: cents}, by_category: {code: cents}}`;
  - `cash`: `{total_cents, payments_cents, cash_expenses_cents, by_person, by_category, payment_count, payments: [{purchase_id, purchase_seq, seq, due_date, amount_cents, shares: {code: cents}}]}` (pagos ordenados por `due_date`, `purchase_id`, `seq`);
  - `cumulative_diff_cents`: devengado acumulado − efectivo acumulado desde el primer evento de las entradas hasta `end` (incluye periodos anteriores no pedidos);
  - `covering_groceries_seq`: `purchase_seq` de la compra de despensa más reciente que devenga en el periodo (null si ninguna).
- Devengado: suma del devengo mensual (§1.3) de cada gasto y compra que cae en el mes. Con ventanas, un costo con n meses se devenga en n ventanas consecutivas a partir de la que contiene su fecha (así lo hace el oráculo con la despensa del día 10). Efectivo: pagos con `due_date` en el periodo (monto completo y partes por persona) más gastos no crédito con `expense_date` en el periodo (a `paid_by` o al titular). `periodos` solo ve registros: los ingresos y las partidas planificadas no entran aquí; las ocurrencias planificadas entran a `accrued` y `cash` solo dentro de la proyección (§3.12, con su fórmula) y al flujo (§3.11).
- Orden: un elemento por mes de `Months` o por ventana de `Windows`, en el orden recibido.
- Invariantes que se prueban (S2-05): Σ `by_person` = `total_cents` en devengado y en efectivo; Σ `payments[].shares` = `amount_cents`; `payments_cents + cash_expenses_cents = total_cents`.

### 3.3 Horizonte · `horizonte` (R-51…R-60)
- Go: `Horizon(in Input, spec HorizonSpec) (HorizonSummary, error)`, `HorizonSpec{Start, End Date; Months []Month}`. Python: `horizonte(datos, inicio, fin, meses)`.
- Argumentos `{start, end, months}` (los tres obligatorios; `start ≤ end`; `months` lista no vacía de `Month`) → `HorizonSummary` (abajo).
- Salida: `payments_by_purchase: [{purchase_id, purchase_seq, in_horizon_cents, complete}]` (compras no canceladas con al menos un pago en [Start, End], ordenadas por (`purchase_date`, `purchase_id`); `complete` = todos sus pagos caen dentro), `paper_cash_cents` (gastos no crédito de `paper_category` en [Start, End]), `cash_total_cents` (todo el efectivo con fecha en [Start, End]), `sum_of_periods_cash_cents` y `accrued_total_cents` (Σ de `periodos(Months)`), `live_debt` (§3.4 con `cut = End`), `monthly_avg_cash_cents = HALF_UP(cash_total/len(Months))`, `monthly_avg_accrued_cents = HALF_UP(accrued_total/len(Months))`, `final_diff_cents = accrued_total − cash_total`.
- Semilla: Start 2026-10-10, End 2027-04-10, Months 2026-10…2027-03 → 1310378 efectivo, 1436248 devengado, promedios 218396 y 239375, diferencia 125870.

### 3.4 Deuda viva e identidad · `deuda_viva` (S2-07, RN-21, K19, M-12)
- Go: `LiveDebt(in Input, cut Date) (LiveDebtResult, error)`. Python: `deuda_viva(datos, corte)`.
- Argumentos `{cut}` (fecha, obligatorio) → `LiveDebtResult` (abajo).
- Salida: `{cut, live_debt_cents, by_purchase: [{purchase_id, purchase_seq, pending_cents, started}], accrued_to_date_cents, cash_to_date_cents, cumulative_diff_cents, unaccrued_balance_cents, identity_holds, identity_without_balance_holds}`.
- Definiciones: deuda viva = Σ de los pagos con `due_date > cut` de las compras no canceladas con `purchase_date ≤ cut` (RN-21); `by_purchase` lleva esas compras con `pending_cents > 0`, ordenadas por (`purchase_date`, `purchase_id`); `started` = la compra tiene algún pago con `due_date ≤ cut`. Devengado acumulado = devengo de los **meses cerrados** (fin de mes ≤ cut); efectivo acumulado = todo el efectivo con fecha ≤ cut; saldo no devengado = Σ costo (RN-02) de gastos y compras con fecha ≤ cut (un gasto a crédito y su compra cuentan una vez) − devengado acumulado. Identidad (K19): `live_debt = cumulative_diff + unaccrued_balance`; `identity_without_balance_holds` dice si además `live_debt = cumulative_diff` (solo con saldo 0).
- Por qué «meses cerrados»: es la única lectura que da las 7 filas de S2-07 con las cifras de `model.json`; el 10-abr-2027 la despensa 4 (comprada ese día) no ha entrado al devengado: 629344 = 125870 + 503474. Las 7 filas, calculadas hoy desde `model.json` (cierre: deuda = diferencia + saldo): 2026-10-31: 292125 = 97375 + 194750; 2026-11-30: 97375 = 97375 + 0; 2026-12-31: 377606 = 125869 + 251737; 2027-01-31: 125870 = 125870 + 0; 2027-02-28: 377606 = 125869 + 251737; 2027-03-31: 125870 = 125870 + 0; 2027-04-10: 629344 = 125870 + 503474.
- M-12 por línea y por compra es de S3; S2 entrega el total y la identidad.

### 3.5 Promedios · `promedios` y `promedio_mensual` (S2-08, K20, K39, RN-27)
- Go: `Averages(total Cents, start, end Date) (Average, error)` y `MonthlyAverage(total Cents, months int) (Cents, error)`. Python: `promedios(total, inicio, fin)`, `promedio_mensual(total, meses)`.
- Argumentos de `promedios` (exactamente una de las dos formas, con todas sus claves): `{total_cents, start, end}` (`start ≤ end`) o `{month, basis}` con `basis` ∈ `cash` · `accrued` · `accrued_groceries`; con `month`, `start`/`end` son los límites del mes y el total sale de `periodos {months: [month]}`: `cash` → `cash.total_cents`, `accrued` → `accrued.total_cents`, `accrued_groceries` → `accrued.by_category[groceries_category]` (0 si falta, §1.1) → `{days, daily_cents, weekly_cents}`. Argumentos de `promedio_mensual`: `{total_cents, month_count}` (los dos obligatorios; `month_count` entero ≥ 1, el `months` de Go y el `meses` de Python) → `{monthly_cents}`.
- Regresión: oct-2026 efectivo 97375/31 → 3141 y 21988; nov-2026 194750/30 → 6492 y 45442; feb-2027 251738/28 → 8991 y 62935; devengado de la despensa de oct-2026 194750/31 → 6282 y 43976. Las lecturas con el diario redondeado (21987, 62937, 43974) son incorrectas (K39). Mensuales: 1310378/6 → 218396 (R-56); 1436248/6 → 239375 (R-57).

### 3.6 Valor estable · `valor_estable` (R-37…R-44)
- Go: `StableValue(items []StableItem, persons []Person) (StableValueResult, error)` (el titular y las claves de los mapas salen de `persons`, §1.1). Python: `valor_estable(partidas, personas)`.
- Argumentos `{items}` (lista de `StableItem`, opcional; si falta, se usan los `stable_items` de las entradas) → salida (abajo).
- Salida: `{month: {total_cents, by_person}, fortnight: {total_cents, by_person, sum_persons_cents, cent_diff_cents}, components: [{label, monthly_cents, by_person}]}` con §1.3; `components` va en el orden de entrada de las partidas. Semilla: despensa 251737 (Karen 74986, Axel 176751) + papel 9950 (Karen) → mes 261687 (Axel 176751, Karen 84936); quincena 130844 (Axel 88376, Karen 42468), suma 130844, diferencia 0.
- Cruce con S3: con las partidas de `sensibilidad[k]` (total de la despensa de N, Karen 149972, papel 19900) reproduce `costo_mensual_total`, `axel_mensual` y `karen_mensual` de las 5 filas (verificado hoy con `Decimal`: 198753, 230220, 261687, 293154, 324621). Esas 15 hojas son de S3; S2 solo las confirma con `vista_sensibilidad` (§3.14) en la fixture `s2-sensibilidad`, ids `sensibilidad[0]`…`sensibilidad[4]`.

### 3.7 Papel · `papel` (R-68…R-73)
- Go: `PaperPlan(in Input) ([]PaperPackage, error)`. Python: `papel(datos)`.
- Argumentos `{}` → lista de paquetes (abajo).
- Paquetes: los `embedded` de `paper_category` en compras (forma `in_groceries`, fecha = `purchase_date`, con su `purchase_seq`) y los gastos no crédito de `paper_category` (forma `cash`), ordenados por fecha; en la misma fecha, primero los `embedded` por (`purchase_id`, posición en `embedded`) y luego los de contado por `id` del gasto; `package` = 1, 2, … Salida por paquete: `{package, date, person, person_name, form, purchase_seq, amount_cents}`. Los paquetes `cash` son los que entran a `papel_contado` del periodo y del horizonte.

### 3.8 Presupuesto del mes · `presupuesto_mes` (S2-11, RN-22, M-01, M-09)
- Go: `BudgetMonth(in Input, m Month) (BudgetMonthResult, error)`. Python `fin_presupuesto.py`: `presupuesto_mes(datos, mes)`.
- Argumentos `{month}` (obligatorio) → `BudgetMonthResult` (abajo).
- Plan de una categoría = `amount_cents` del renglón cuya vigencia contiene el mes (null si no hay). Real = devengado de la **parte del hogar** en la categoría en el mes (RN-22, RN-02). Salida: `{month, days_in_month, days_left, categories: [{category, plan_cents, actual_cents, remaining_cents, used_pct, warning, limit, state}], total: {plan_cents, actual_cents, remaining_cents, used_pct, warning, limit, state, declared_total_row_cents}, accrued_household_cents, cash_household_cents}`.
- `categories` lleva **todas las categorías activas** de las entradas (con o sin renglón; sin renglón, `state` `no_plan`), ordenadas por `code` (orden de bytes); una categoría inactiva no aparece ni suma al total (M-01).
- `used_pct = HALF_UP(actual·100/plan)` entero (M-01). **Banderas con comparación exacta en enteros** (S2-11, «frontera exacta»): `warning ⇔ actual·100 ≥ budget_warning_pct·plan`; `limit ⇔ actual·100 ≥ budget_limit_pct·plan`. `state`: `normal`, `warning` (≥ aviso y < excedido), `exceeded` (≥ excedido) o `no_plan` (sin renglón). Por eso 79 999 de 100 000 muestra `used_pct` 80 con `state` `normal` (ver §8, observación para U1/U3).
- `total.plan_cents` = Σ planes de las categorías activas (M-01). PROVISIONAL (`S2-presupuesto-total`): qué hacer con el renglón del total (`category` null), que se informa en `declared_total_row_cents`, y con un plan de 0 (hoy: `used_pct` null, banderas en `false`, `state` `no_plan`).
- `days_left`: días del mes desde hoy contando hoy (M-01); 0 si el mes ya pasó; `days_in_month` si es futuro.
- **Fórmula de cada campo** (mes m; «parte del hogar» = `by_person[holder]`, §1.4; «devengado» y «efectivo» = los de `periodos {months: [m]}`, solo registros, §3.2):
  - Categoría con renglón: `plan_cents` = su monto; `actual_cents` = devengado de la parte del hogar en la categoría; `remaining_cents = plan_cents − actual_cents` (negativo si se rebasó, como en C1); `used_pct`, `warning`, `limit` y `state` con las fórmulas de arriba.
  - Categoría sin renglón (`no_plan`): `actual_cents` igual que arriba; `plan_cents`, `remaining_cents` y `used_pct` en null (C1: `used_pct` «ausente sin plan»); `warning` y `limit` en `false`.
  - `total.actual_cents` = Σ `actual_cents` de **todas** las filas de `categories`, con o sin renglón (M-01: «real = Σ parte del hogar devengada en el mes»); `total.plan_cents` = Σ de los `plan_cents` no nulos (0 si ninguna categoría tiene renglón, que es el caso plan 0, PROVISIONAL); `total.remaining_cents = total.plan_cents − total.actual_cents`; `total.used_pct`, `warning`, `limit` y `state` con las mismas fórmulas sobre los totales.
  - `accrued_household_cents` = devengado del mes en `by_person[holder]` y `cash_household_cents` = efectivo del mes en `by_person[holder]`, de **todas** las categorías, activas o no (lectura literal del «real» de M-01; difiere de `total.actual_cents` solo si una categoría inactiva tiene gasto en el mes). El efectivo excluye la parte que paga directo una persona externa (RN-04). Con la semilla en 2026-10: 119764 (= `periodos[0].devengado.axel`, el ejemplo de M-01) y 59882 (= `periodos[0].efectivo.axel`).
- Ejemplo (RN-22): plan 200000 en dic-2026, real del hogar 176751 (`periodos[2].devengado.axel`) → `used_pct` 88, `warning` true, `limit` false.

### 3.9 Cambio de presupuesto con alcance · `aplicar_cambio_presupuesto` (K2, RN-22)
- Go: `ApplyBudgetChange(rows []BudgetRow, ch BudgetChange) (BudgetChangeResult, error)`. Python: `aplicar_cambio_presupuesto(renglones, cambio)`.
- Argumentos `{row_id, month, amount_cents, scope}` = los campos de `BudgetChange` (§2.8): `row_id`, `month` y `amount_cents` obligatorios; `scope` opcional (`onward` por omisión) → `BudgetChangeResult` (abajo). Los renglones salen de `entradas.datos.budgets` (en Python, `renglones` es esa lista), nunca de los argumentos; un `row_id` que no está en ellos es `invalid_argument`.
- Salida: `{rows, effective_row_id, deleted_row_ids}` (renglones ordenados por categoría, primero el del total con `category` null y luego por `code`, y dentro de cada una por `valid_from`; `deleted_row_ids` ascendentes; ids nuevos = máximo + 1, + 2, …). Reglas de `updateBudget` en C1: `onward` cierra el renglón el último día del mes anterior, abre uno desde `month` sin fin y elimina los posteriores de la misma categoría; `this_month` deja el monto nuevo solo ese mes y el anterior sigue desde el mes siguiente hasta su fin original, sin tocar los posteriores; si el renglón empieza en `month`, el monto cambia en ese renglón (con `onward` queda sin fin y se eliminan los posteriores; con `this_month` el resto de la vigencia queda en otro renglón con el monto anterior). `month` fuera de la vigencia → `out_of_range`.
- Ejemplo (RN-22): 200000 desde 2026-10-01 sin fin; cambio a 250000 en 2026-12: `onward` → oct–nov 200000 y desde dic 250000 sin fin; `this_month` → oct–nov 200000, dic 250000, desde ene-2027 200000.

### 3.10 Plan contra real por mes · `plan_contra_real` (M-14)
- Go: `PlanVsActual(in Input, from Month, n int, category string) ([]PlanActualRow, error)` (1 ≤ n ≤ 12; `category` vacío = total). Python `fin_presupuesto.py`: `plan_contra_real(datos, desde, n, categoria=None)`.
- Argumentos `{from, count, category}`: `from` (`Month`) y `count` (1–12, la `n` de Go y de Python) obligatorios; `category` opcional (ausente o null = total) → lista de `{month, plan_cents, actual_cents}`, una fila por mes desde `from`, consecutivos y en orden. Cada fila copia `plan_cents` y `actual_cents` de lo que da `presupuesto_mes {month}` para ese mes: de la fila de `categories` de esa categoría (null en `plan_cents` si no tiene renglón) o, sin `category`, de `total`. Una `category` que no está entre las activas de las entradas es `invalid_argument`.

### 3.11 Flujo por quincena · `flujo_quincenal` (S2-10, M-02, M-10)
- Go: `FortnightFlow(in Input, from Fortnight, count int) ([]FortnightRow, error)` (1 ≤ count ≤ 24). Python `fin_flujo.py`: `flujo_quincenal(datos, desde, n)`.
- Argumentos `{from, count}` (`Fortnight` y entero 1–24, la `n` de Python; los dos obligatorios) → lista de filas (abajo).
- **Bases** (quincena q = [`start`, `end`], hoy T = `params.today`): «efectivo del hogar» de un gasto registrado no crédito = `amount_cents` si `paid_by` es null o el titular, 0 si lo paga otra persona (§2.4); de una ocurrencia planificada de gasto = su monto completo (sale del bolsillo del titular, §2.7); de un pago programado = la parte del titular. Es la misma base que el efectivo por persona de §3.2 (`paid_by`) y que el «contado + partes del hogar pagadas» de M-02. Un gasto con `credit` no cuenta como gasto: su efectivo son los pagos (`credit_payments`). Abono vigente = elemento de `paid` con `paid_date ≤ T`. Ocurrencia «ligada» = tiene captura según la regla PROVISIONAL del final. Solo cuentan fechas dentro de q; `cadence` no entra al flujo.
- **Campos por quincena** (centavos salvo fechas, banderas y marcas):

| Campo | Fórmula |
|---|---|
| `fortnight`, `start`, `end`, `pay_date` | §1.2 |
| `income.fixed_cents` · `income.extra_cents` | Σ `incomes` `fixed` · `extra` con `income_date` en q (capturados, sin importar T) |
| `income.expected_cents` | Σ ocurrencias de partidas `income` en q **no ligadas** (M-02: la esperada cuando no hay captura) |
| `income.income_is_expected` | `expected_cents > 0` |
| `income.total_cents` | `fixed_cents + extra_cents + expected_cents` |
| `credit_payments.by_person` | Por persona, Σ de su parte en los pagos programados con `due_date` en q (compras no canceladas) |
| `credit_payments.household_cents` | `by_person[holder]` |
| `credit_payments.external_informative_cents` | Σ partes de las personas con `pays_share_directly` (Karen, RN-04); nunca entran a `spent` ni a `committed` |
| `credit_payments.household_paid_cents` | Σ por pago de mín(parte del titular, abonos vigentes del titular a ese pago) |
| `credit_payments.household_pending_cents` | `household_cents − household_paid_cents` |
| `fixed_upcoming.planned_cents` | Σ montos completos de las ocurrencias de partidas `fixed_expense` y `upcoming_expense` en q (ligadas o no) |
| `fixed_upcoming.paid_cents` | Σ efectivo del hogar de gastos registrados `fixed` y `upcoming` no crédito con fecha en q y ≤ T |
| `fixed_upcoming.pending_cents` | Σ montos de esas ocurrencias **no ligadas** + Σ efectivo del hogar de gastos registrados `fixed` y `upcoming` no crédito con fecha en q y > T |
| `variable_spent_cents` | Σ efectivo del hogar de gastos registrados `variable` no crédito con fecha en q y ≤ T |
| `variable_estimated_cents` · `savings_planned_cents` | `null` mientras `S2-flujo-quincenal` siga abierto (métodos sin decidir); no entran a ningún saldo |
| `spent_cents` | `variable_spent_cents + fixed_upcoming.paid_cents + credit_payments.household_paid_cents` (M-02: contado + partes del hogar pagadas hasta hoy) |
| `committed_cents` | `credit_payments.household_pending_cents + fixed_upcoming.pending_cents` (M-02) |
| `balance_cents` | `income.total_cents − spent_cents − committed_cents` |
| `cumulative_balance_cents` | Σ `balance_cents` desde la primera quincena de la consulta (`from`) hasta esta, inclusive; parte de 0 (el modelo no tiene saldo inicial de cuentas; M-10 empieza en la Q1 del mes en curso) |
| `provisional` | `null`, o `["S2-flujo-quincenal"]` si en q alguna ocurrencia tiene una captura con su `planned_item_id` dentro de la ventana de la liga (ahí la regla PROVISIONAL decide si se cuenta) |

- Pagado: la parte de una persona en un pago está pagada si sus abonos vigentes la cubren; el capital solo se repone al completarse el pago (K15, asunto de S3). La parte de Karen va en `by_person` y en `external_informative_cents`, nunca en `spent` ni `committed` (RN-04). Orden: una fila por quincena desde `from`, consecutivas.
- **S2-10, 16 quincenas** (`s2-flujo-quincenal`, desde 2026-10-q2 hasta 2027-06-q1, un pago por quincena): despensa 1 → 2026-10-q2…2026-12-q1 con hogar 59882 y Karen 37493 (R-13…R-16); despensa 2 → 2026-12-q2…2027-02-q1 con 88375 ×3 y 88377 (R-18…R-21); despensa 3 → 2027-02-q2…2027-04-q1 (R-23…R-26); despensa 4 → 2027-04-q2…2027-06-q1 (R-28…R-31); Karen 37493 en las 16, en `external_informative_cents`.
- **Contraste con el valor estable** (S2-10, R-37…R-44; lo hace `TestFlujoQuincenal` con dos salidas de S2 de la misma fixture, sin campo nuevo): (1) ciclo: para cada despensa 2, 3 y 4, Σ `household_cents` de sus 4 quincenas = 2 × `valor_estable.mes.axel` (353502 = 2 × 176751) y Σ `by_person.karen` = 2 × `componentes.karen_mitad_despensa` (149972 = 2 × 74986), con diferencia 0; (2) por quincena: `valor_estable.quincena.axel` (88376, R-41) contra `household_cents` de esas 12 quincenas: la prueba informa «valor estable − pago» (+1 en 9 quincenas, −1 en 3) y no lo corrige (nota (d) de S2.md §6: son dos cifras canon distintas). La despensa 1 (capital distinto) queda fuera del contraste.
- PROVISIONAL (`S2-flujo-quincenal`): cómo se liga una ocurrencia planificada con lo capturado (lectura de trabajo: un gasto o ingreso con ese `planned_item_id` en la misma quincena, o en el mismo mes si la frecuencia es mensual o mayor, la da por cumplida). Esa es la «ventana de la liga» de la marca `provisional`.

### 3.12 Proyección de 12 meses · `proyeccion` (S2-09, RN-28, K21, M-11)
- Go: `Projection(in Input, from Month) (ProjectionResult, error)` (`from` vacío = mes de hoy). Python: `proyeccion(datos, desde=None)`.
- Argumentos `{from}` (`Month`, opcional; ausente = mes de hoy) → `ProjectionResult` (abajo).
- Salida: `{repeat_cadence_enabled, months: [12 × {index, month, accrued, cash, income_cents, planned_expense_cents, net_cash_cents, cumulative_diff_cents, provisional}]}`.
- **Fórmulas por mes m** (registros = gastos, ingresos y compras de las entradas, con todos sus pagos; ocurrencias = las de §2.7 de partidas activas sin `credit_id` con fecha en m y **no ligadas**, regla PROVISIONAL de §3.11):
  - `accrued` y `cash`: como en §3.2, sobre registros + ocurrencias de gasto + `cadence` (solo meses 7–12 y solo con el interruptor encendido). Una ocurrencia de gasto se trata como un gasto de contado con `paid_by` null: efectivo = su monto completo en su fecha, al titular; devengo RN-02 de su monto en `accrual_months` meses desde su mes, por persona con las partes de su `share_rule` (RN-03, §1.3); `by_category` con su `category`.
  - `income_cents` = Σ `incomes` con fecha en m + Σ ocurrencias de partidas `income` en m (la de la quincena del 30 con `second_half_amount_cents` si lo trae).
  - `planned_expense_cents` = Σ montos completos de las ocurrencias de gasto en m (total, no la parte del hogar); ya están dentro de `cash.total_cents` y se informan aparte.
  - `net_cash_cents = income_cents − cash.by_person[holder]` (efectivo del hogar: incluye las ocurrencias de gasto, que salen del titular, y excluye la parte de Karen).
  - `cumulative_diff_cents`: como en §3.2 (devengado acumulado − efectivo acumulado desde el primer evento), con ocurrencias y cadencia incluidas.
  - `provisional`: `null` o lista ordenada de bloqueos: `"S2-flujo-quincenal"` si en m alguna ocurrencia tiene una captura con su `planned_item_id` dentro de la ventana de la liga (§3.11); `"S2-proyeccion-interruptor"` en los meses 7–12 con el interruptor encendido.
- Contenido: meses 1–12 = registros + ocurrencias + `cadence` según las fórmulas. Interruptor apagado (por omisión, K21): ninguna despensa nueva ni papel alterno; en la semilla, abr–sep 2027 traen solo los pagos ya programados: el último de la despensa 3 (2027-04-15, R-26) y los de la despensa 4 (2027-04-30, 05-15, 05-30, 06-15; R-28…R-31), más el devengo de la despensa 4 (abr y may) y las partidas de una fixture de CSV. Encendido: `provisional` incluye `"S2-proyeccion-interruptor"` en los meses 7–12; apagado, no.
- Meses 1–6 = P1–P6 (oct-2026…mar-2027): 12 de 12 cifras de devengado y efectivo iguales (R-45…R-50), porque la semilla no tiene eventos entre los días 1 y 9; se prueba con la semilla sin partidas (`s2-proyeccion-7-12-defecto`). Con partidas (`s2-proyeccion-csv`) cada mes suma sus ocurrencias con las fórmulas de arriba. Sin `planned_items` la proyección es válida (decisión H).

### 3.13 Metas y fondos · `metas` (S2-12, RN-18, M-07, M-16)
- Go: `SavingsGoals(in Input) (SavingsResult, error)`. Python `fin_metas.py`: `metas(datos)`.
- Argumentos `{}` → `SavingsResult` (abajo).
- Salida: `{month_summary: {saved_cents, contributed_this_month_cents, planned_this_month_cents}, goals: [{goal_id, kind, saved_cents, target_cents, progress_tenths_pct, progress_pct, months_left, needed_monthly_cents, plan: [{month, cents}], state}]}`.
- `goals` lleva las metas activas (`is_active`), ordenadas por `goal_id`; `plan` va por mes ascendente.
- Firmes: ahorrado = Σ aportaciones − Σ retiros con fecha ≤ hoy (RN-18, M-07); `progress_tenths_pct = HALF_UP(saved·1000/target)` (décimas, RN-18: 314672 de 629344 → 500 = 50.0 %); `progress_pct = HALF_UP(saved·100/target)` (entero, M-07/M-16); sin tope, una meta cumplida puede pasar de 100. Plan de aportaciones: con n meses, base HALF_UP(restante/n) y la última ajusta, de modo que Σ `plan` = `target − saved` al centavo (S2-12); meta cumplida → `plan` vacío y `needed_monthly_cents` 0. `planned_this_month_cents` = Σ `planned_monthly_cents` de las metas activas.
- **`month_summary`** (mes de hoy, el «mes en curso» de M-07), sobre las metas de `goals` (activas; fondos y metas, M-07 «fondo y metas»): `saved_cents` = Σ `saved_cents` de `goals`; `contributed_this_month_cents` = Σ `amount_cents` de los movimientos `deposit` de esas metas con `movement_date` en el mes de hoy y ≤ hoy (las «aportaciones» de M-07; los retiros no se restan aquí porque ya restan en `saved_cents`); `planned_this_month_cents` = Σ de los `planned_monthly_cents` no nulos de esas metas. PROVISIONAL (`S2-metas-resumen`): si el resumen cuenta solo las metas activas (lectura de trabajo) o también las inactivas con saldo.
- **Sin monto o sin fecha meta** (un `fund` puede no traerlos, §2.9): sin `target_cents`, `target_cents`, `progress_tenths_pct` y `progress_pct` van en null; sin `target_cents` o sin `target_date`, `needed_monthly_cents` va en null y `plan` vacío; sin `target_date`, `months_left` va en null. `saved_cents` se calcula siempre.
- PROVISIONAL (`S2-metas`): cómo se cuentan los `months_left`, qué porcentaje muestra la pantalla (décimas de RN-18 o entero de C1) y los valores de `state`.

### 3.14 Vistas del oráculo (claves de `model.json`) · `vista_periodo`, `vista_horizonte`, `vista_valor_estable`, `vista_papel`, `vista_sensibilidad`, `vista_interes_por_monto`
Adaptadores finos sobre §3.2, §3.3–3.5, §3.6 y §3.7 que emiten las claves literales de `model.json`, en centavos, para las fixtures canon. En una fixture, el `id` de una consulta de vista es la ruta de `model.json` que reproduce (FORMATO §2.1): `periodos[0]`, `valor_estable`, `horizonte`, `papel`, `sensibilidad[k]`, `interes_por_monto[k]`. Exigen la forma de la semilla (un titular y una sola persona externa; partidas `despensa` y `papel_contado`); otra forma es `oracle_view_shape`. Las personas salen con su `code`.
- Argumentos (Python `fin_vistas.py`: una función por vista, con `datos` y estos argumentos): `vista_periodo {month}` (obligatorio); `vista_horizonte {start, end, months}` (los tres obligatorios, como §3.3); `vista_valor_estable {items}` (opcional, como §3.6); `vista_papel {}`; `vista_sensibilidad {items}` (obligatorio); `vista_interes_por_monto {purchase_id}` (obligatorio; un id que no está en `credit_purchases` es `invalid_argument`).
- `vista_periodo(mes)` → `despensa_que_cubre` = `covering_groceries_seq`; `devengado.{despensa_total = accrued.by_category[groceries_category], papel_contado = devengado de gastos no crédito de paper_category, total, <persona>}`; `efectivo.{pagos = payments_cents, papel = efectivo de gastos no crédito de paper_category, total, <persona>, n_pagos = payment_count, detalle_pagos[{despensa = purchase_seq, fecha, monto, <persona>}]}`; `diferencia_acumulada`; y `etiqueta`, `inicio`, `fin` con el mes calendario (texto; no se comparan con `model.json`).
- `vista_horizonte(inicio, fin, meses)` → `pagos_despensa_<seq>` (compra completa en el horizonte) o `pagos_despensa_<seq>_en_horizonte` (parcial); `papel_contado`; `total_efectivo`; `suma_periodos_efectivo`; `total_devengado`; `deuda_viva_10_abr` (el sufijo es el día y el mes abreviado de `fin`); `deuda_viva_detalle.{despensa_<seq>_pendiente si started, si no despensa_<seq>}` (solo compras con pendiente > 0); `promedio_mensual_efectivo`; `promedio_mensual_devengado`; `diferencia_final`.
- `vista_valor_estable(partidas)` → `mes.{total, axel, karen}`; `quincena.{total, axel, karen, suma_axel_karen, diferencia_centavo}`; `componentes.{mitad_despensa, papel_contado_mes, karen_mitad_despensa, axel_mitad_despensa}`.
- `vista_papel()` → lista `{paquete, compra, quien = person_name, forma = "dentro de despensa <seq>" o "de contado", monto}`.
- `vista_sensibilidad({items})` (cruce, id `sensibilidad[k]`) → `{costo_mensual_total = month.total_cents, axel_mensual = month.by_person[holder], karen_mensual = month.by_person[la persona externa]}` de `valor_estable` con esas partidas: `despensa` de ciclo 2 con `total_cents` = `sensibilidad[k].total` y partes karen 149972 (`parametros.karen_fija.total`) y axel `sensibilidad[k].axel_total`, más `papel_contado` 19900 de Karen. 15 hojas de S3 que S2 confirma.
- `vista_interes_por_monto({purchase_id})` (cruce, id `interes_por_monto[k]`) → `{monto = principal_cents, intereses = total − monto, total, pago_quincenal = amount_cents del pago con seq 1}`, donde `total` = Σ del devengado de la compra en todos sus meses (§3.2) y debe ser igual a Σ del efectivo de sus pagos (si no, `sum_mismatch`). 28 hojas de S3 que S2 confirma por devengado y efectivo.

### 3.15 Errores
`unknown_field`, `invalid_date`, `invalid_zone`, `invalid_argument`, `shares_mismatch`, `sum_mismatch`, `missing_purchase`, `income_not_holder`, `budget_overlap`, `out_of_range`, `cadence_before_month_7`, `oracle_view_shape`. Go: `type Error struct{ Code, Detail string }`; Python: `ErrorMotor(code, detalle)`. Una fixture puede esperar un error: `{"<consulta>": {"error": "<code>"}}`. Una `funcion` desconocida en una consulta es `invalid_argument` (FORMATO §4.1).

## 4. Las 596 hojas de `model.json`: 159 de S2 y 437 de S3
Cada hoja de S2 tiene función (vista), fixture y prueba. Las llaves se expanden: `{a,b}` enumera y `{0..5}` es un rango; `n` es el número de hojas del renglón. Además de la prueba nombrada, toda hoja la cubre `TestOraculoHojas/<ruta>` (Go) y la línea `model.json (ramas de s2)` de `compara.py s2`.

| # | Hojas | n | Vista ← función | Fixture | Prueba Go |
|---|---|---|---|---|---|
| 1 | `periodos[0].{despensa_que_cubre,devengado.{despensa_total,papel_contado,total,axel,karen},efectivo.{pagos,papel,total,axel,karen,n_pagos},diferencia_acumulada}` | 13 | `vista_periodo` 2026-10 ← `periodos` | `s2-periodos-p1` | `TestCanon/R-45` |
| 2 | `periodos[0].efectivo.detalle_pagos[0].{despensa,monto,axel,karen}` | 4 | ídem | `s2-periodos-p1` | `TestSumaPartes` |
| 3 | `periodos[1].{despensa_que_cubre,devengado.{despensa_total,papel_contado,total,axel,karen},efectivo.{pagos,papel,total,axel,karen,n_pagos},diferencia_acumulada}` | 13 | `vista_periodo` 2026-11 | `s2-periodos-p2` | `TestCanon/R-46` |
| 4 | `periodos[1].efectivo.detalle_pagos[{0,1}].{despensa,monto,axel,karen}` | 8 | ídem | `s2-periodos-p2` | `TestSumaPartes` |
| 5 | `periodos[2].{despensa_que_cubre,devengado.{despensa_total,papel_contado,total,axel,karen},efectivo.{pagos,papel,total,axel,karen,n_pagos},diferencia_acumulada}` | 13 | `vista_periodo` 2026-12 | `s2-periodos-p3` | `TestCanon/R-47` |
| 6 | `periodos[2].efectivo.detalle_pagos[{0,1}].{despensa,monto,axel,karen}` | 8 | ídem | `s2-periodos-p3` | `TestSumaPartes` |
| 7 | `periodos[3].{despensa_que_cubre,devengado.{despensa_total,papel_contado,total,axel,karen},efectivo.{pagos,papel,total,axel,karen,n_pagos},diferencia_acumulada}` | 13 | `vista_periodo` 2027-01 | `s2-periodos-p4` | `TestCanon/R-48` |
| 8 | `periodos[3].efectivo.detalle_pagos[{0,1}].{despensa,monto,axel,karen}` | 8 | ídem | `s2-periodos-p4` | `TestSumaPartes` |
| 9 | `periodos[4].{despensa_que_cubre,devengado.{despensa_total,papel_contado,total,axel,karen},efectivo.{pagos,papel,total,axel,karen,n_pagos},diferencia_acumulada}` | 13 | `vista_periodo` 2027-02 | `s2-periodos-p5` | `TestCanon/R-49` |
| 10 | `periodos[4].efectivo.detalle_pagos[{0,1}].{despensa,monto,axel,karen}` | 8 | ídem | `s2-periodos-p5` | `TestQuincenaFinDeMes/feb-2027` |
| 11 | `periodos[5].{despensa_que_cubre,devengado.{despensa_total,papel_contado,total,axel,karen},efectivo.{pagos,papel,total,axel,karen,n_pagos},diferencia_acumulada}` | 13 | `vista_periodo` 2027-03 | `s2-periodos-p6` | `TestCanon/R-50` |
| 12 | `periodos[5].efectivo.detalle_pagos[{0,1}].{despensa,monto,axel,karen}` | 8 | ídem | `s2-periodos-p6` | `TestSumaPartes` |
| 13 | `valor_estable.mes.{total,axel,karen}` | 3 | `vista_valor_estable` ← `valor_estable` | `s2-valor-estable` | `TestCanon/R-37…R-39` |
| 14 | `valor_estable.quincena.{total,axel,karen,suma_axel_karen,diferencia_centavo}` | 5 | ídem | `s2-valor-estable` | `TestCanon/R-40…R-42`, `TestFlujoQuincenal` |
| 15 | `valor_estable.componentes.{mitad_despensa,papel_contado_mes,karen_mitad_despensa,axel_mitad_despensa}` | 4 | ídem | `s2-valor-estable` | `TestCanon/R-43`, `TestCanon/R-44` |
| 16 | `horizonte.{pagos_despensa_1,pagos_despensa_2,pagos_despensa_3_en_horizonte,papel_contado}` | 4 | `vista_horizonte` ← `horizonte` | `s2-horizonte` | `TestCanon/R-59`, `TestCanon/R-60` |
| 17 | `horizonte.{total_efectivo,suma_periodos_efectivo,total_devengado,diferencia_final}` | 4 | ídem | `s2-horizonte` | `TestCanon/R-51`, `R-52`, `R-58` |
| 18 | `horizonte.{deuda_viva_10_abr,deuda_viva_detalle.{despensa_3_pendiente,despensa_4}}` | 3 | `vista_horizonte` ← `deuda_viva` | `s2-horizonte` | `TestCanon/R-53…R-55`, `TestInvarianteDeudaViva/2027-04-10` |
| 19 | `horizonte.{promedio_mensual_efectivo,promedio_mensual_devengado}` | 2 | `vista_horizonte` ← `promedio_mensual` | `s2-horizonte` | `TestCanon/R-56`, `R-57`, `TestPromedios` |
| 20 | `papel[{0..5}].{paquete,monto}` | 12 | `vista_papel` ← `papel` | `s2-papel` | `TestCanon/R-68…R-73` |

Total: 122 + 12 + 13 + 12 = **159**. Las **437** de S3 (`parametros` 9, `despensas` 144, `credito` 143, `sensibilidad` 45, `disposicion_por_despensa` 68, `interes_por_monto` 28) llegan como entrada o no las usa S2: `despensas[0..3]` (fechas, montos, partes, capital, `papel_dentro`) forman las `credit_purchases` de la semilla; `parametros.papel` es el monto de los gastos de papel de contado; 15 hojas de `sensibilidad` (`costo_mensual_total`, `axel_mensual`, `karen_mensual`) y las 28 de `interes_por_monto` se confirman con `s2-sensibilidad` y `s2-interes-por-monto` (cruces, 43 hojas) sin pasar a ser de S2; `credito`, `disposicion_por_despensa` y el resto de `parametros` no entran al motor financiero (las filas S3 de la regresión llegan en `_datos/insumo-s3.json`, §5.1). Ninguna hoja queda sin dueño (S2-02 + S3-01).

## 5. Fixtures (`oraculo/fixtures/s2/`, formato en `oraculo/fixtures/FORMATO.md`)
Datos compartidos: `_datos/semilla-oraculo.json` = `Input` armado solo con hojas de S3 (`despensas[0..3]` → 4 gastos de despensa a crédito con `accrual_months` 2 y sus 4 compras; papel de contado de `papel[2]` y `papel[4]`, de Karen; personas `axel` titular y `karen` externa), `today` 2026-10-11, interruptor apagado, sin abonos (`paid` vacío en todos los pagos) y sin ingresos ni partidas planificadas. Contenido completo en §5.3; autoría de todos los archivos de esta carpeta en §5.0. `_datos/insumo-s3.json`: las 55 filas S3 de la regresión con el esquema de §5.1. `_datos/interes-por-monto.json`: las 7 compras de `s2-interes-por-monto` (§5.2).
| Caso | Consultas | Fuente | Criterio |
|---|---|---|---|
| `s2-periodos-p1` … `p6` | `vista_periodo` 2026-10 … 2027-03 | `model.json:periodos[0..5]` | S2-01, S2-02, S2-09 |
| `s2-valor-estable` · `s2-horizonte` · `s2-papel` | vistas de §3.14 | `model.json:valor_estable` · `horizonte` · `papel` | S2-01, S2-02 |
| `s2-sensibilidad` | `vista_sensibilidad` ×5, ids `sensibilidad[0]`…`[4]` | `model.json:sensibilidad[0..4]` (cruce con S3, 15 hojas) | S2-03 |
| `s2-interes-por-monto` | `vista_interes_por_monto` ×7, ids `interes_por_monto[0]`…`[6]` (§5.2) | `model.json:interes_por_monto[0..6]` (cruce con S3, 28 hojas) | S2-03 |
| `s2-deuda-viva-identidad` | `deuda_viva` en los 7 cierres | S2-07 (§3.4) | S2-07 |
| `s2-fin-de-mes-feb-2027` · `s2-mes-30-dias` · `s2-mes-31-dias` · `s2-semana` | `calendario`, `periodos`, `flujo_quincenal` | S2-06; el gasto del 2026-10-31 se calcula a mano (§5.2) | S2-06 |
| `s2-zona-medianoche` | `calendario` con `now` | S2-14 | S2-14 |
| `s2-promedios` | `promedios` ×4, `promedio_mensual` ×2 | K39, RN-27 | S2-08 |
| `s2-ventanas-dia-5` | `periodos` con ventanas y un gasto del 2026-11-05 | S2-09 | S2-09 |
| `s2-proyeccion-7-12-defecto` · `s2-proyeccion-7-12-cadencia` · `s2-proyeccion-csv` · `s2-liga-provisional` | `proyeccion` apagada con la semilla sin partidas; encendida (PROVISIONAL, `S2-proyeccion-interruptor`); apagada con partidas de CSV sin capturas ligadas; `proyeccion` y `flujo_quincenal` con una captura ligada (PROVISIONAL, `S2-flujo-quincenal`) (§5.2) | K21, RN-28 | S2-09 |
| `s2-flujo-quincenal` | `flujo_quincenal` desde 2026-10-q2, 16 quincenas (hasta 2027-06-q1), y `valor_estable` para el contraste de §3.11 (ids libres) | `model.json:despensas[0..3].pagos`, `valor_estable` | S2-10 |
| `s2-presupuesto-scope` · `s2-alertas-80-100` | `aplicar_cambio_presupuesto`; `presupuesto_mes` en 79 %, 79.999 %, 80 %, 99 %, 100 % y 101 % | K2, RN-22, S2-11 | S2-11 |
| `s2-metas` | `metas`: meses exactos, residuo de centavos, meta cumplida (PROVISIONAL) | RN-18, S2-12 | S2-12 |

Total: **27 fixtures** (las 24 de la 1.0 más `s2-interes-por-monto`, `s2-proyeccion-csv` y `s2-liga-provisional`). Las consultas de vista usan como `id` la ruta de `model.json` que reproducen; las demás, nombres libres (FORMATO §2.1).

### 5.0 Autoría de `oraculo/fixtures/s2/` y orden de trabajo (alcance disjunto, protocolo §3)
- **Autor único: el agente Python de S2.** Escribe las 27 fixtures de la tabla y los 3 archivos de `_datos/` (`semilla-oraculo.json`, §5.3; `insumo-s3.json`, §5.1; `interes-por-monto.json`, §5.2) con un script propio, `oraculo/fin_genera_fixtures.py`, que lee `model.json` y `datos/semilla/data.json` con FORMATO §3 y escribe de forma determinista. Nadie más los edita; un cambio se pide al integrador.
- **Origen de `esperado`.** Consulta con `id` ruta: el valor convertido de `model.json` o `data.json` (FORMATO §2.1, §6). Consulta de nombre libre: las cifras que fijan esta interfaz y `S2.md` (citadas en `fuente`) y, para las demás claves, el cálculo del oráculo, avisado en `descripcion` («resto calculado por el oráculo»). Si el motor difiere, se documentan ambos valores (protocolo §6); no se ajusta ninguno.
- **El agente Go mientras falten.** No escribe en `oraculo/` ni crea sustitutos. (1) Arma en sus pruebas el mismo `Input` con la receta de §5.3 leyendo `model.json` (ayudante de prueba `seedFromModel`) y prueba con él sus funciones, `TestOraculoHojas` y las filas S2 de `TestCanon`, más las cifras de esta interfaz; (2) lo que lee esos archivos (las filas S3 de `TestCanon`, `TestOraculoFixtures` y las pruebas que recorren todas las fixtures) falla con «falta oraculo/fixtures/s2/<archivo> (autor: agente Python de S2)», nunca con `t.Skip` (S2 §3); (3) lo anota como pendiente y el integrador corre §6 y §7 completos al juntar. Cuando `semilla-oraculo.json` existe, `TestCanon` exige además que sea igual, como JSON, al `Input` de `seedFromModel`: una diferencia es hallazgo de uno de los dos.
- **El agente Python mientras falte el motor:** `python3 -I oraculo/compara.py s2 --sin-motor` (FORMATO §5).

### 5.1 `_datos/insumo-s3.json` (filas S3 de `TestCanon`)
Lo escribe el agente Python de S2 con un script que convierte las hojas de `model.json` y `data.json` de la tabla de abajo; representa lo que S3 entrega y S2 no calcula. Lo lee `TestCanon` (Go). Forma:
```json
{
  "formato": "insumo-s3/1",
  "descripcion": "Filas S3 de la regresión tal como S3 las entrega; S2 no las calcula",
  "filas": {
    "R-01": {"despensas[0].capital": 309450},
    "R-04": {"parametros.tasa_pct": {"num": 80050, "den": 309450}},
    "R-05": {"data.json:financiamiento.dos_quincenas.tasa_pct": {"num": 53450, "den": 309450, "decimales": 2}},
    "R-13": {"despensas[0].pagos[0].fecha": "2026-10-30", "despensas[0].pagos[0].monto": 97375, "despensas[0].pagos[0].axel": 59882, "despensas[0].pagos[0].karen": 37493, "despensas[0].pagos[0].capital_repuesto": 77363},
    "R-65": {"credito.comparacion_series.esperado_A": [1077363, 1154726, 1232089, 832089, 909450, 1009450, 1109450, 1209450, 809450, 909450, 1009450, 1109450, 1209450, 809450]}
  }
}
```
- `filas` lleva exactamente las **55 filas S3** (R-01…R-36, R-61…R-67, R-74…R-85), de R-01 a R-85; cada una es un objeto `{llave: valor}` con las llaves de la tabla, en ese orden.
- **Llave** = ruta de una hoja de `model.json` con la sintaxis de FORMATO §3, o `data.json:` + ruta (R-05, R-06). Los objetos van siempre aplanados a sus hojas; una lista de escalares (R-65, R-66) va como lista JSON, en su orden.
- **Valores:** dinero en centavos (FORMATO §3: `int(valor × 100)`, exacto); fechas y textos literales (`fecha`, `tipo`); R-04 razón `{num, den}` (FORMATO §3: el binary64 correctamente redondeado de 100·num/den es `tasa_pct`); R-05 `{num, den, decimales}`: HALF_UP(100·53450/309450) a 2 decimales = 17.27, el valor guardado (nota (c) de S2.md §6).

| Fila | Llaves |
|---|---|
| R-01 · R-02 · R-03 | `despensas[0].capital` · `despensas[0].intereses` · `despensas[0].total` |
| R-04 | `parametros.tasa_pct` |
| R-05 | `data.json:financiamiento.dos_quincenas.tasa_pct` |
| R-06 | `data.json:financiamiento.dos_quincenas.total`, `data.json:financiamiento.dos_quincenas.intereses` |
| R-07 · R-08 · R-09 | `despensas[0].axel.total` · `.pago_por_quincena` · `.intereses` |
| R-10 · R-11 · R-12 | `despensas[0].karen.total` · `.pago_por_quincena` · `.intereses` |
| R-13…R-16 | `despensas[0].pagos[i].{fecha,monto,axel,karen,capital_repuesto}`, i = 0…3 (una i por fila) |
| R-17 · R-22 · R-27 | `despensas[k].{capital,intereses,total}`, k = 1 · 2 · 3 |
| R-18…R-21 · R-23…R-26 · R-28…R-31 | `despensas[k].pagos[i].{fecha,monto,axel,karen}`, k = 1 · 2 · 3, i = 0…3 |
| R-32…R-36 | `parametros.karen_fija.total` · `.pago` · `.capital` · `.intereses` · `.conjunto` (una por fila) |
| R-61 · R-62 · R-63 | `credito.minimo_A` · `credito.minimo_B` · `parametros.limite_credito` |
| R-64 | `credito.eventos[0].disponible_A`, `credito.eventos[0].disponible_B` |
| R-65 · R-66 | `credito.comparacion_series.esperado_A` · `.esperado_B` (lista de 14) |
| R-67 | `credito.eventos[i].{fecha,tipo,monto,disponible_A,disponible_B}`, i = 0…16 (85 llaves) |
| R-74…R-78 | `sensibilidad[k].{intereses,total,pago_quincenal,axel_total}`, k = 0…4 |
| R-79…R-85 | `interes_por_monto[k].{intereses,total,pago_quincenal}`, k = 0…6 |

`TestCanon` con cada fila S3, en este orden y sin `SKIP`: (1) cada llave contra `model.json` o `data.json` con FORMATO §3; (2) las que S2 usa como entrada, contra sus entradas: R-01…R-03 y R-17, R-22, R-27 = `principal_cents`, `interest_cents`, `total_cents` de las compras 1–4 de `semilla-oraculo.json`; R-07…R-12 = `shares[].total_cents` e `interest_cents` de la compra 1 y la parte de cada persona en su pago 1; R-13…R-31 = `payments[i]` (`due_date`, `amount_cents`, `shares`, y `principal_cents` = `capital_repuesto` en R-13…R-16); R-74…R-78 = partidas de `s2-sensibilidad`; R-79…R-85 = compras de `_datos/interes-por-monto.json`; (3) R-04…R-06, R-32…R-36 y R-61…R-67 solo pasan (1): S2 las recibe y no las usa.

### 5.2 Detalle de casos nuevos o precisados
- **`s2-mes-31-dias`** (S2-06): `{"$ref": "_datos/semilla-oraculo.json", "parche": {…}}` con `params.today` = `"2026-10-31"` y `expenses` y `categories` reemplazadas (Merge Patch reemplaza listas completas) por las de la semilla más una categoría `variable` y un gasto de fixture (`fuente` `manual`, no es dato de Axel): `amount_cents` 50000, `expense_date` 2026-10-31, `kind` `variable`, `payment_method` `cash`, `paid_by` null, `shares` [titular 50000]. Se comparan: `calendario {date: 2026-10-31}` → `fortnight` 2026-10-q2, 2026-10-16 a 2026-10-31, `pay_date` 2026-10-30; `flujo_quincenal` desde 2026-10-q1, 2 quincenas → q1: `variable_spent_cents` 0 y `spent_cents` 0; q2: `variable_spent_cents` 50000, `spent_cents` 50000 y, sin abonos, `household_pending_cents` 59882 (pago del 2026-10-30, R-13) dentro de `committed_cents`; `periodos` 2026-10 → `cash.cash_expenses_cents` 50000 y `cash.total_cents` 147375 (97375 + 50000). Un motor que cierre la quincena 2 el día 30 da 0 en q2 y falla.
- **`s2-interes-por-monto`** (S2-03, cruce): `_datos/interes-por-monto.json` = personas y `params` de la semilla y 7 compras (`id` k + 1, `expense_id` null, `category` = `groceries_category`, `accrual_months` 2, `purchase_date` 2026-10-10, partes solo del titular) con `principal_cents` = `interes_por_monto[k].monto`, `interest_cents` y `total_cents` de `model.json`, y la serie que entregaría S3 (no la calcula S2): 4 pagos el 2026-10-30, 2026-11-15, 2026-11-30 y 2026-12-15 (calendario de la despensa 1, RN-06), montos `pago_quincenal` en los tres primeros y `total − 3·pago_quincenal` en el cuarto (RN-07), capital HALF_UP(capital/4) con el último ajustando (K14). Consultas `interes_por_monto[0]`…`[6]` con `vista_interes_por_monto {purchase_id}`.
- **`s2-proyeccion-7-12-defecto`**: la semilla tal cual (sin partidas): meses 1–6 = P1–P6 (12 de 12) y meses 7–12 con solo los pagos programados y el devengo de la despensa 4; `provisional` null en los 12.
- **`s2-proyeccion-csv`**: semilla + partidas de fixture (`fuente` `manual`): un ingreso quincenal del titular con `second_half_amount_cents`, un gasto fijo mensual con regla «Hogar» y un gasto venidero `once` con «Conjunto 50/50» y `accrual_months` 2; sin capturas ligadas, así que `provisional` es null y el caso es firme. Se comparan `income_cents`, `planned_expense_cents`, `accrued`, `cash`, `net_cash_cents` y `cumulative_diff_cents` de los 12 meses.
- **`s2-liga-provisional`** (`provisional` = `{bloqueo: docs/estado/bloqueos/S2-flujo-quincenal.md, que: "liga ocurrencia ↔ captura"}`): semilla + el ingreso quincenal de `s2-proyeccion-csv` + un `income` `fixed` capturado el 2026-10-14 con su `planned_item_id` (la ocurrencia es el 2026-10-15: la lectura de trabajo la liga y la opción «fecha exacta» no). Se comparan `flujo_quincenal` desde 2026-10-q1 (2 quincenas: `provisional` `["S2-flujo-quincenal"]` en q1 y null en q2; `expected_cents` 0 en q1 con la lectura de trabajo) y `proyeccion` (marca solo en 2026-10).

### 5.3 Contenido completo de `_datos/semilla-oraculo.json`
Receta exacta: la siguen el script del agente Python y el ayudante `seedFromModel` del agente Go. `c(x)` = centavos con FORMATO §3; `D[k]` = `despensas[k]`, k = 0…3; `PA` = `papel`; «por p» = para p en `axel`, `karen`, en ese orden. Van las 12 claves de `Input`, en este orden:

| Clave | Contenido |
|---|---|
| `params` | `{today: "2026-10-11", holder: "axel", groceries_category: "despensa", paper_category: "papel", repeat_cadence_enabled: false}`; `zone` y los umbrales, por omisión |
| `persons` | `{code: "axel", display_name: "Axel", participation: "holder", pays_share_directly: false}`, `{code: "karen", display_name: "Karen", participation: "external", pays_share_directly: true}`; sin co-titular (no tiene partes, §1.4) |
| `categories` | `{code: "despensa", name: "Despensa", kind: "variable", is_active: true}`, `{code: "papel", name: "Papel", kind: "variable", is_active: true}` |
| `expenses` 1–4 (despensas) | id k + 1: `{id, kind: "variable", category: "despensa", amount_cents: c(D[k].capital), expense_date: D[k].fecha, accrual_months: 2, payment_method: "credit", paid_by: null, shares: [{person: p, amount_cents: c(D[k].p.capital)} por p], credit_purchase_id: k + 1, planned_item_id: null}` |
| `expenses` 5–6 (papel de contado) | id 5 = `PA[2]`, id 6 = `PA[4]`: `{id, kind: "variable", category: "papel", amount_cents: c(parametros.papel) = 19900, expense_date: PA[i].compra, accrual_months: 1, payment_method: "cash", paid_by: "karen", shares: [{person: "karen", amount_cents: 19900}], credit_purchase_id: null, planned_item_id: null}` (`karen` = `PA[i].quien` en minúsculas) |
| `credit_purchases` 1–4 | id k + 1: `{id, credit_id: 1, expense_id: k + 1, purchase_date: D[k].fecha, principal_cents: c(D[k].capital), interest_cents: c(D[k].intereses), total_cents: c(D[k].total), status: "active", shares: [{person: p, principal_cents: c(D[k].p.capital), interest_cents: c(D[k].p.intereses), total_cents: c(D[k].p.total)} por p], embedded: [{category: "papel", person: p, amount_cents: c(D[k].p.papel_dentro)} por p con papel_dentro > 0], payments: [por i = 0…3, con x = D[k].pagos[i]: {seq: i + 1, due_date: x.fecha, amount_cents: c(x.monto), principal_cents: c(x.capital_repuesto), interest_cents: c(x.monto) − c(x.capital_repuesto), shares: [{person: p, amount_cents: c(x.p)} por p], paid: []}], projected: false}`; `category` y `accrual_months` se omiten (`expense_id` no es null, §2.5) |
| `incomes`, `planned_items`, `budgets`, `savings_goals`, `savings_movements` | `[]` |
| `cadence` | `{purchases: [], expenses: []}` |
| `stable_items` | Las 2 de §2.11: `{label: "despensa", cycle_months: 2, total_cents: c(D[1].total) = 503474, shares: [{person: "karen", total_cents: c(D[1].karen.total) = 149972}, {person: "axel", total_cents: c(D[1].axel.total) = 353502}]}` y `{label: "papel_contado", cycle_months: 2, total_cents: c(parametros.papel) = 19900, shares: [{person: "karen", total_cents: 19900}]}` |

- Los códigos `despensa` y `papel` son identificadores de la fixture: S1/S6 aún no siembran categorías y ninguna cifra canon depende del código. El `kind` `variable` de las categorías y de los 6 gastos no cambia ninguna cifra: los 4 gastos a crédito no cuentan como gasto en el flujo (§3.11) y los 2 de papel los paga Karen (0 en el efectivo del hogar).
- `embedded` queda en `[{papel, karen, 19900}]` en la compra 1 y `[{papel, axel, 19900}]` en las compras 2, 3 y 4 (`papel[0]`, `[1]`, `[3]`, `[5]`).
- Comprobado el 11-oct-2026 (script temporal de `correccion-contrato-S2-r2`): la receta pasa las validaciones de §2.4 y §2.5 (28 de 28) y, con las reglas de §1.3, §3.2, §3.3, §3.4, §3.6 y §3.7, reproduce las **159 de 159** hojas S2 de `model.json`.
- Los casos que cambian algo lo hacen con `parche` (FORMATO §2); un gasto o una categoría que agregue un caso lleva el siguiente id libre (7 en adelante).

## 6. Cruce motor ↔ oráculo
1. Oráculo: `python3 -I oraculo/fin_run_fixtures.py --fixtures oraculo/fixtures/s2 --out oraculo/salidas/s2.json` → «N fixtures, N coinciden, 0 diferencias» contra `esperado` (S2-03).
2. Motor: `cd backend && go build -o ../datos/tmp/motores ./cmd/motores && cd .. && datos/tmp/motores s2 dump --fixtures oraculo/fixtures/s2 --out datos/tmp/motor-s2.json` (rutas relativas al directorio actual).
3. Comparación: `python3 -I oraculo/compara.py s2` → `model.json (ramas de s2): 159 de 159 hojas coinciden` (S2-02) y `model.json (cruces con otras ramas): 43 de 43 hojas coinciden` (15 de `sensibilidad` y 28 de `interes_por_monto`, S2-03), más las líneas de fixtures y la lista de diferencias. Formato de las salidas, de las líneas y códigos de salida: `FORMATO.md` §4 y §5.

## 7. Pruebas y calidad
| Prueba Go (`go test ./internal/motor_fin/... -run <patrón> -count=1 -v`) | Criterio | Qué comprueba |
|---|---|---|
| `TestCanon` (85 subpruebas `R-01`…`R-85`) | S2-01 | Lee `oraculo/model.json` y `datos/semilla/data.json` (R-05, R-06) al correr; las 30 filas S2 contra las salidas del motor; las 55 filas S3 con los pasos (1)–(3) de §5.1 sobre `_datos/insumo-s3.json`; además, `semilla-oraculo.json` = `seedFromModel` (§5.0, §5.3); 0 `SKIP` (si falta un archivo de `oraculo/fixtures/s2/`, falla con su nombre y su autor) |
| `TestOraculoHojas` (159 subpruebas) | S2-02 | Cada hoja de §4 contra `model.json` |
| `TestOraculoFixtures` (una por fixture) | S2-03 | Salidas contra `esperado` |
| `TestSumaPartes` | S2-05 | Invariantes de §3.2, §2.5 y §3.11 en todas las fixtures |
| `TestQuincenaFinDeMes` (`feb-2027`, `mes-30-dias`, `mes-31-dias`, `semana`) | S2-06 | Límites de K8 y el gasto del 2026-10-31 con `today` 2026-10-31 (campos de §5.2) |
| `TestInvarianteDeudaViva` (7 cierres) | S2-07 | Identidad de K19 y en qué cierres vale sin saldo (P2, P4, P6) |
| `TestPromedios` | S2-08 | Mensual R-56/R-57, diarios y semanales de §3.5 |
| `TestProyeccionDefecto`, `TestProyeccionCadencia` | S2-09 | `s2-proyeccion-7-12-defecto` (meses 1–6 = P1–P6; 7–12 solo pagos programados) y `s2-proyeccion-csv` (ocurrencias con las fórmulas de §3.12); los dos modos de K21; la marca `provisional` de `s2-liga-provisional` |
| `TestFlujoQuincenal` | S2-10 | 16 quincenas: 59882/37493 (R-13…R-16) y 88375 ×3 + 88377 con 37493 (R-18…R-21, R-23…R-26, R-28…R-31); contraste con el valor estable de §3.11 (ciclos con diferencia 0 y diferencias por quincena informadas) |
| `TestPresupuestoUmbrales`, `TestPresupuestoScope` | S2-11 | Fronteras exactas; alcance de K2 |
| `TestMetasAhorro` | S2-12 | Σ plan = restante; avance |
| `TestZonaHoy` | S2-14 | Medianoche local |

Las pruebas encuentran la raíz del repositorio subiendo hasta hallar `oraculo/model.json`; ningún valor esperado se pega en el código (S2 §3). S2-13: `grep -rn "float32\|float64" backend/internal/motor_fin --include=*.go | grep -v _test.go` vacío. S2-15: `go vet` y `gofmt -l` limpios y cobertura ≥ 85 % (K25). S2-04: sha256 de `model.json` igual antes y después, y `grep -n reporte-compra-latex oraculo/*.py` vacío tras la corrección de `verify3.py` (§0.5).

## 8. Lo que esta interfaz no decide (S2-16)
| Tema | Dónde queda | Qué sale PROVISIONAL |
|---|---|---|
| Cadencia de los meses 7–12 con el interruptor encendido (monto y fecha de cada despensa, mes y pagador del papel alterno) | `docs/estado/bloqueos/S2-proyeccion-interruptor.md` | `cadence`, meses 7–12 encendidos, `s2-proyeccion-7-12-cadencia` |
| Variable estimado de M-10, ahorro planificado por quincena, liga ocurrencia ↔ captura | `docs/estado/bloqueos/S2-flujo-quincenal.md` | `variable_estimated_cents` y `savings_planned_cents` (null hasta decidir), la regla de la liga de §3.11 y la marca `provisional` de §3.11 y §3.12 donde se aplica, `s2-liga-provisional` |
| «Meses que faltan», porcentaje en pantalla (RN-18 décimas contra C1 entero), estados de una meta | `docs/estado/bloqueos/S2-metas.md` | `months_left`, `needed_monthly_cents`, `plan`, `state`, `s2-metas` |
| Conjunto de metas del resumen del mes (solo activas o también inactivas con saldo) | `docs/estado/bloqueos/S2-metas-resumen.md` | `month_summary.saved_cents` y `contributed_this_month_cents` (lectura de trabajo: las activas de `goals`); `s2-metas` ya es PROVISIONAL |
| Renglón del total del presupuesto y plan de 0 | `docs/estado/bloqueos/S2-presupuesto-total.md` | `declared_total_row_cents`, el caso plan 0 |
| Dónde viven `send_time`, `advance_days`, `weekday`, `day_of_month` | `S1-parametros-avisos` (puerta 2) | Nada de S2: los umbrales 80/100 ya son `alerts.*` (K44) |

Observación para U1/U3 (no bloquea S2): `used_pct` entero puede mostrar 80 con `state` `normal` entre 79.5 % y 79.99 %. Fuera de S2: M-03, M-04, M-12 por línea, M-13 y M-17 (S3); M-05 y M-15 (S4); M-06 (S5).

## 9. Criterio → sección
| Criterio | Sección | Criterio | Sección |
|---|---|---|---|
| S2-01 | §4, §5, §5.1, §5.3, §7 `TestCanon` | S2-09 | §3.2, §3.12, §2.4, §2.7, §2.10, §2.12, §5.2 |
| S2-02 | §4, §6, `FORMATO.md` §2.1, §5 | S2-10 | §3.11, §3.6, §2.6, §2.7, §5 |
| S2-03 | §0.4, §1.1 (mapas), §3 (argumentos), §3.14, §5, §5.0, §5.3, §6, `FORMATO.md` §1, §2.1, §4, §4.1 | S2-11 | §3.8, §3.9, §3.10, §2.1 |
| S2-04 | §0.5 (corrección de `verify3.py`), §7 | S2-12 | §3.13, §2.9 |
| S2-05 | §1.3, §2.5, §3.2, §7 | S2-13 | §0.3, §1.1, §7 |
| S2-06 | §1.2, §3.1, §3.11, §5.2 | S2-14 | §1.2, §2.1, §3.1 |
| S2-07 | §3.4 | S2-15 | §7 |
| S2-08 | §1.3, §3.5 | S2-16 | §8 |
