# Reglas de negocio · Finanzas Hogar

**Versión:** 1.0 · 10-oct-2026 · agente S1-modelo-datos (WF-0, fase S1)
**Fuentes:** `docs/plan/decisiones.md` (A a J), `docs/plan/plan.md` §1 a §3, `oraculo/model.json` (v3), `oraculo/model.py`, `datos/semilla/data.json`.
**Esquema:** `backend/db/migrations/00001_esquema.sql` · **Diccionario:** `docs/contratos/diccionario-datos.md`.

**Cómo leer las citas.** `model.json: ruta = valor` con índices desde 0 (`despensas[0]` es la despensa 1). `data.json: ruta`. Los ejemplos se escriben en pesos para leerlos; en la base van en centavos (3,094.50 = `309450`). Cada cifra tomada de `model.json` se recalculó con enteros de centavos en un verificador independiente (reporte `docs/estado/wf0/S1-modelo-datos.md`, sección Verificación). Cuando el reporte y el oráculo no traen el dato, el ejemplo dice **ilustrativo** y sus insumos no son datos de la casa.

**Estado de cada regla.** FIRME = sale de una decisión cerrada o del oráculo sin margen. PROVISIONAL = depende de un bloqueo abierto (se indica la ruta). Nada PROVISIONAL entra en cifras canon.

## Índice
| Regla | Tema | Estado |
|---|---|---|
| RN-01 | Redondeo a centavos | FIRME |
| RN-02 | Devengado contra efectivo | FIRME (último mes con centavos impares: PROVISIONAL) |
| RN-03 | Reparto configurable, 50/50 por defecto | FIRME |
| RN-04 | Karen: externa, paga directo a Kueski | FIRME (reposición por parte: PROVISIONAL) |
| RN-05 | Ingreso quincenal 15/30, solo de Axel | FIRME |
| RN-06 | Quincenas, fin de mes y calendario de pagos | FIRME (rangos de la quincena: PROVISIONAL) |
| RN-07 | Series de pagos y ajuste de centavos en el último pago | FIRME en lo que fija el oráculo; generalización PROVISIONAL |
| RN-08 | Tasa por compra con IVA incluido y tabla editable por plazo | FIRME |
| RN-09 | Disponible de la línea: lectura A/B, B por defecto | FIRME |
| RN-10 | Promociones manuales | FIRME |
| RN-11 | Comparador de escenarios: restricciones duras y objetivos | FIRME en definición; valores y alcance PROVISIONAL |
| RN-12 | Frescura de precios: 30 días | FIRME el umbral; peso PROVISIONAL |
| RN-13 | Membresía de Sam's: punto de equilibrio | FIRME |
| RN-14 | Tiempos de comida (4) y escalado de porciones | FIRME |
| RN-15 | Inventario híbrido | FIRME |
| RN-16 | Mercados con día de visita | FIRME |
| RN-17 | Salud: IMC y meta de calorías | IMC FIRME; fórmula de meta PROVISIONAL |
| RN-18 | Fondo y metas de ahorro | FIRME |
| RN-19 | Comidas fuera de casa | FIRME |
| RN-20 | Equivalencia por unidad base y sustitutos | FIRME |
| RN-21 | Deuda viva y compras vivas | FIRME |
| RN-22 | Presupuesto, real contra plan y alertas 80/100 % | FIRME |

## Convenciones (decisiones técnicas de S1)
- **C-01 Dinero.** Centavos enteros MXN en columnas `_cents`, tipo INTEGER en tablas STRICT: SQLite rechaza un `12.5` (probado). En pantalla: `$` con dos decimales y separador de miles (`$3,094.50`). Ningún cálculo de dinero usa flotantes; los intermedios usan fracciones exactas.
- **C-02 Redondeo.** ROUND_HALF_UP a centavo en el último paso de cada cálculo (RN-01), salvo la base de los montos de una serie de pagos, que el oráculo trunca (RN-07).
- **C-03 Fechas de negocio.** TEXT ISO `AAAA-MM-DD`, sufijo `_date`, validadas con `CHECK (date(x) IS x)`: se rechazan `2027-02-30` y `10-oct-2026` (probado). En pantalla `10-oct-2026`, y `Vie 30 oct 2026` donde importe el día de la semana.
- **C-04 Instantes.** TEXT ISO 8601 en UTC con «Z», sufijo `_at`, por defecto `strftime('%Y-%m-%dT%H:%M:%fZ','now')`. Se muestran convertidos a America/Mexico_City.
- **C-05 Zona horaria.** America/Mexico_City define «hoy», el corte de cada día, las quincenas y el horario de los jobs. Se usa la base de zonas horarias del sistema, nunca un desfase fijo.
- **C-06 Borrado lógico: no hay borrado lógico general.** Los catálogos (personas, cuentas, categorías, reglas de reparto, tiendas, productos, presentaciones, créditos, recetas, partidas planificadas, metas, destinatarios) se **archivan** con `is_active = 0`, y sus llaves foráneas `RESTRICT` impiden borrarlos si tienen historia (probado: no se puede borrar a Karen con partes registradas). Los movimientos (gastos, ingresos, pagos, planes de comida) se borran de verdad y su imagen anterior queda en `audit_log`. **Por qué:** una columna `deleted_at` obliga a filtrar en cada suma de dinero, y un solo olvido suma dinero borrado. La historia ya queda en `audit_log` y en los respaldos diarios.
- **C-07 Auditoría de cambios.** `audit_log` la escribe la API en la misma transacción de cada alta, cambio o baja: quién (`user_account_id`), cuándo, tabla, id y el renglón antes y después en JSON. Es inmutable (disparadores que rechazan UPDATE y DELETE, probado) y no tiene llave foránea para sobrevivir a cualquier baja. Los jobs escriben con usuario NULL y `request_id` = código del job.
- **C-08 Partes y series congeladas.** `expense_share`, `credit_purchase_share`, `scheduled_payment` y `scheduled_payment_share` se guardan ya calculadas al centavo cuando se registra la compra o el gasto. Cambiar una regla de reparto o una tasa afecta solo a lo nuevo. Para recalcular algo viejo hay que editarlo y el cambio queda en `audit_log`.
- **C-09 El hogar va a nombre del titular.** La parte de dinero del hogar Axel–Ximena se registra a nombre del titular, como en el oráculo, que la llama `axel`. A quién se atribuye un gasto (quién comió, quién gastó) va aparte, en `expense.person_id`.
- **C-10 Tasas como fracción exacta** `rate_num`/`rate_den`. Una tasa redondeada a cuatro decimales no reproduce el oráculo (RN-08).
- **C-11 Unidades enteras.** Porcentajes en puntos base (`_bp`, 0 a 10000 = 0 a 100 %). Cantidades físicas en milésimas de la unidad base (`_milli`). Peso en gramos, estatura en milímetros, kilocalorías en milésimas. Los rangos de los CHECK de salud (estatura 500 a 2500 mm, peso 20 a 400 kg, meta 500 a 10000 kcal) son solo para atrapar errores de captura.
- **C-12 Nombres.** Tablas y columnas en inglés, en singular y snake_case. Documentos y textos visibles en español de México. Hay claves estables (`code`) para lo que los motores buscan por nombre: personas, categorías, avisos y tiempos de comida.
- **C-13 Integridad.** `PRAGMA foreign_keys = ON` en cada conexión (lo activa S6 en el DSN, porque dentro de la transacción de goose el PRAGMA no surte efecto). STRICT en las 54 tablas. CHECK para enumerados, rangos, fechas y sumas dentro del renglón. Disparadores para reglas de las decisiones: Karen sin cuenta, nadie con cuenta, salud o plan de comidas pasa a externa, ingresos solo del titular, salud y comidas sin externos. Índices únicos: un solo titular y una sola co-titular (dos cuentas como máximo). Lo que SQLite no puede expresar está listado en el diccionario como «invariantes».
- **C-14 Sin credenciales en la base.** Solo hashes argon2id de contraseñas y SHA-256 de tokens de sesión. El token de Telegram y la llave de la API de Claude se leen al ejecutar y nunca se guardan.
- **C-15 Idempotencia de las capturas.** Cada POST de captura trae `Idempotency-Key` (C1 §1 punto 8). La API guarda la clave por cuenta en `idempotency_key` con la huella SHA-256 del cuerpo y la respuesta: si se repite con el mismo cuerpo devuelve la respuesta guardada; con otro cuerpo responde 409. La retención (24 h propuestas) es PROVISIONAL, bloqueo `C1-parametros-tecnicos`.
- **C-16 Descripción opcional del gasto.** `expense.description` admite NULL: el registro en 3 toques (U1 F-1, N-054) no la pide y la pantalla muestra la categoría. La API no inventa un texto de relleno.

---

## RN-01 Redondeo a centavos · FIRME
**Definición.** Todo monto se guarda en centavos enteros. Cuando un cálculo produce fracciones de centavo, se redondea con ROUND_HALF_UP (la mitad sube) en el último paso, nunca a la mitad del cálculo. La única excepción es la base de los montos de una serie de pagos (RN-07).
**Ejemplo.** Intereses de Karen en la despensa 1: 800.50 × 1,191.50 / 3,094.50 = 308.2196… → **308.22** (`model.json: despensas[0].karen.intereses = 308.22`).
**Esquema.** Todas las columnas `_cents`.

## RN-02 Devengado contra efectivo · FIRME (último mes con centavos impares: PROVISIONAL)
**Definición.**
- **Devengado:** el costo de un gasto es su importe más los intereses de su compra a crédito, si la hay. Se reparte en `accrual_months` meses calendario a partir del mes de `expense_date`, en partes iguales: cada mes HALF_UP(costo / meses), y el último mes ajusta para que la suma dé el costo exacto. El ajuste solo aparece con totales de centavos impares, que el oráculo no trae: PROVISIONAL, bloqueo `S1-redondeo-series`. Por persona, cada parte que no es el resto se reparte igual y el titular toma la diferencia.
- **Efectivo:** el dinero sale cuando se paga. De contado, débito o transferencia sale en `expense_date`, del bolsillo de `paid_by_person_id` (NULL = el titular). A crédito sale en cada `scheduled_payment.due_date`, según las partes del pago. La parte de quien paga directo al acreedor (Karen) se muestra, pero no sale del flujo del titular (RN-04).
- El presupuesto del hogar mide la parte del hogar. La de Karen es informativa, porque el plan se hace solo con el ingreso de Axel (decisión B) y ella paga lo suyo (decisión A).
- Los periodos P1 a P6 del oráculo van del día 10 al 9. Como las compras son el día 10 y los pagos los días 15 y 30, cada Pk da exactamente las mismas cifras que el mes calendario correspondiente, de oct-2026 a mar-2027 (verificado renglón por renglón). El sistema usa meses calendario (decisión B) sin mover ninguna cifra canon.

**Ejemplo.** Despensa 1, costo 3,895.00 (`model.json: despensas[0].total`), bimestral (`accrual_months` = 2):

| mes | devengado | efectivo |
|---|---|---|
| oct-2026 (P1) | 1,947.50 (`periodos[0].devengado.despensa_total`) | 973.75: pago del 30-oct (`periodos[0].efectivo.total`) |
| nov-2026 (P2) | 1,947.50 (`periodos[1].devengado.despensa_total`) | 1,947.50: pagos del 15-nov y 30-nov (`periodos[1].efectivo.total`) |

En octubre, el devengado de Karen es HALF_UP(1,499.72 / 2) = 749.86 y el de Axel 1,947.50 − 749.86 = 1,197.64 (`periodos[0].devengado.karen` / `.axel`). Al cierre de octubre la diferencia acumulada es 973.75 (`periodos[0].diferencia_acumulada`). En el horizonte 10-oct-2026 a 10-abr-2027 el devengado suma 14,362.48 y el efectivo 13,103.78 (`horizonte.total_devengado`, `horizonte.total_efectivo`).
**Esquema.** `expense.expense_date`, `accrual_months`, `payment_method`, `paid_by_person_id`; `credit_purchase.interest_cents`; `scheduled_payment.due_date`; `scheduled_payment_share`.

## RN-03 Reparto configurable, 50/50 por defecto · FIRME
**Definición.**
- Cada gasto dice de quién es: del hogar Axel–Ximena, de Karen, o conjunto con su porcentaje (decisión B). Se expresa con una **regla de reparto** (`cost_share_rule`). Cada parte es un porcentaje (`pct_bp`) o un monto fijo (`fixed_cents`), y la **persona del resto** recibe lo demás y absorbe los centavos.
- **Prioridad:** regla del renglón → regla del gasto → regla de la categoría → regla por defecto. La regla por defecto es «Hogar» (100 % al titular), porque Karen solo entra en la despensa común y en sus compras con Kueski (decisión A).
- La regla «Conjunto 50/50» es Karen 5000 pb y el resto al titular. Es la regla de la categoría Despensa (decisión B: «la despensa base es conjunto 50/50 por defecto»).
- **Cálculo:** se agrupan los renglones por regla y se suma cada grupo, como hace el reporte con el «conjunto». En cada grupo, una parte por porcentaje vale HALF_UP(suma × `pct_bp` / 10000) y una parte fija vale `fixed_cents` (si supera la suma, se rechaza). El resto recibe la suma menos las partes. Al final se suman las partes por persona y se guardan congeladas en `expense_share`.
- Los renglones «mixto» de `data.json` (una parte de las piezas para cada quien) se reparten por cantidad con una regla de un solo uso (`is_template = 0`).

**Ejemplo.** Los renglones «conjunto» de la despensa 1 suman 1,702.00 (`data.json: cuentas.conjunto.total`) → Karen 851.00 y Axel 851.00 (`data.json: cuentas.conjunto.mitad = 851.0`; `model.json: parametros.karen_fija.conjunto = 851.0`). Totales del gasto: Axel 1,052.00 + 851.00 = **1,903.00** y Karen 340.50 + 851.00 = **1,191.50** (`model.json: despensas[0].axel.capital`, `despensas[0].karen.capital`), que suman 3,094.50, el importe del gasto. Renglón mixto: elote, 5 piezas por 67.50, 3 son de Karen → regla de un solo uso Karen 6000 pb → 40.50 (`data.json: cuentas.karen.lineas[4].importe = 40.5`).
**Esquema.** `cost_share_rule`, `cost_share_rule_part`, `category.default_cost_share_rule_id`, `expense.cost_share_rule_id`, `expense_item.cost_share_rule_id`, `expense_share`.

## RN-04 Karen: participante externa que paga directo a Kueski · FIRME (reposición por parte: PROVISIONAL)
**Definición.**
- Karen (`participation = external`) no tiene cuenta (disparador) ni entra al plan de comidas o de salud (disparadores). Solo tiene partes en la despensa común y en sus compras con Kueski. No se modelan las finanzas de su familia: no tiene ingresos ni gastos propios en el sistema.
- `pays_share_directly = 1`: ella paga su parte de cada pago programado directo a Kueski. Para el flujo de efectivo del titular esa parte es **informativa**: aparece en el calendario y en su estado de cuenta, pero no cuenta como salida de Axel.
- Cuando ella paga se registra un `credit_payment` con `payer_person_id` = Karen, y ese pago **repone la línea** con la parte de capital de lo que pagó. Al cubrirse el pago completo, se repone el capital del pago programado. Cómo se divide ese capital entre su parte y la de Axel cuando pagan por separado no está en el oráculo: PROVISIONAL, bloqueo `S1-reposicion-parte-externa`.
- Su tope por quincena es un objetivo del optimizador, no una restricción dura (decisión C).
- El sistema genera su estado de cuenta en PDF (`generated_report` con `kind = external_statement_pdf`) y Axel lo reenvía a mano.

**Ejemplo.** Pago de la despensa 1 del 30-oct-2026: 973.75 = Axel 598.82 + Karen 374.93 (`model.json: despensas[0].pagos[0].monto`, `.axel`, `.karen`). Del flujo de Axel salen 598.82; los 374.93 los paga Karen. Con el pago completo se reponen 773.63 de capital (`despensas[0].pagos[0].capital_repuesto`) y el disponible en lectura B queda en 7,679.13 (`credito.comparacion_series.esperado_B[0]`). El tope de Karen es 374.93 (`parametros.karen_fija.pago`), y su papel de contado queda «fuera del cap de 374.93» (`papel[2].nota`).
**Esquema.** `person.participation`, `pays_share_directly`, `fortnight_payment_cap_cents`; `scheduled_payment_share`; `credit_payment.payer_person_id`; `generated_report`.

## RN-05 Ingreso quincenal 15/30, solo de Axel · FIRME
**Definición.** Solo se registran ingresos del titular (disparadores en `income` y `planned_item`); Ximena no aporta ingresos al modelo. El sueldo es una partida planificada quincenal (`planned_item` con `kind = income` y `frequency = fortnightly`) cuyas ocurrencias caen en los días de pago de la quincena: el 15, y el 30 o el último día si el mes es más corto. El monto de la quincena del 15 va en `amount_cents`; si el del 30 es distinto va en `second_half_amount_cents` (U1 N-279), y un cambio de montos rige desde una fecha (`start_date`; la partida anterior se cierra con `end_date`). Cada quincena Axel captura el monto real (`income` con `kind = fixed`, ligado a la partida). Los ingresos extra se capturan a mano (`kind = extra`). No hay corrimiento por fin de semana: la fecha es la decidida, igual que en el oráculo, donde el 15-nov-2026 cae en domingo y se mantiene (`model.json: despensas[0].pagos[1].fecha = 2026-11-15`).
**Ejemplo.** El reporte y el oráculo no traen montos de ingreso (el CSV de Axel aún no existe), así que el ejemplo es de fechas: las quincenas de febrero de 2027 pagan el 15-feb-2027 y el 28-feb-2027, el mismo calendario que usa Kueski (`model.json: despensas[1].pagos[3].fecha = 2027-02-15`, `despensas[2].pagos[0].fecha = 2027-02-28`).
**Esquema.** `income`, `planned_item` (`amount_cents`, `second_half_amount_cents`), `period.pay_date`; disparadores `trg_income_holder_only_*` y `trg_planned_item_income_holder_*`.

## RN-06 Quincenas, fin de mes y calendario de pagos · FIRME (rangos de la quincena: PROVISIONAL)
**Definición.**
- **Mes calendario:** del día 1 al último. **Quincena 1:** días 1 a 15, paga el 15. **Quincena 2:** del día 16 al último día del mes, paga el 30, o el último día si el mes tiene menos de 30 (28 o 29 de febrero). En meses de 31 días se paga el 30 y el día 31 pertenece a la quincena 2. El CHECK de `period` lo valida (probado: rechaza una quincena 1 que paga el 14). Los días de pago son FIRMES (oráculo). Que la quincena abarque del 1 al 15 y del 16 al fin de mes es la opción C del bloqueo `U1-periodos-semana-quincena`; la otra lectura es «de pago a pago». Queda PROVISIONAL, igual que en U1. Si se elige «de pago a pago», cambian el CHECK de `period` (migración nueva) y las cifras por quincena, pero no las mensuales ni las fechas de pago.
- **Calendario de una compra quincenal (Kueski):** las fechas candidatas son los días de pago de las quincenas. El primer pago es la primera candidata que queda al menos `first_payment_min_days` (15) días después de la compra, y los demás son las candidatas consecutivas.

**Ejemplo.** Compra del 10-oct-2026: el 15-oct queda a 5 días (no cuenta) y el 30-oct a 20 días, así que los pagos son 30-oct, 15-nov, 30-nov y 15-dic-2026 (`model.json: despensas[0].pagos[0..3].fecha`). Compra del 10-feb-2027: el primer pago es el 28-feb-2027, porque «en febrero el 30 es el 28» (`despensas[2].pagos[0].fecha`).
**Esquema.** `period` (CHECK), `credit.first_payment_min_days`, `credit.payment_frequency`, `scheduled_payment.due_date`.

## RN-07 Series de pagos y ajuste de centavos en el último pago · FIRME en lo que fija el oráculo; generalización PROVISIONAL
**Definición.**
- **Monto de cada pago:** base = total / n truncado a centavos (ROUND_DOWN). Los pagos 1 a n−1 valen la base y el último vale total − (n−1) × base. Así lo calcula el oráculo (`model.py`, líneas 55 a 57); con HALF_UP no se reproduce.
- **Capital que repone cada pago:** base = HALF_UP(capital / n), y el último ajusta. Los intereses de cada pago son monto − capital.
- **Partes por persona en cada pago:** cada persona que no es el resto lleva su parte (Karen: su total entre n), y el titular toma el monto menos las partes de los demás, así absorbe los centavos.
- **Pago quincenal «de referencia»** (cotizaciones y tabla de sensibilidad): HALF_UP(total sin redondear / n). Es informativo y puede diferir en un centavo de la serie real.
- El oráculo no fija la base de las partes por persona cuando el total de alguien no se divide exacto entre n, ni el devengo con centavos impares: queda PROVISIONAL, bloqueo `S1-redondeo-series`.

**Ejemplo.** Despensa 2, total 5,034.74 en 4 pagos: 5,034.74 / 4 = 1,258.685 → base 1,258.68. Pagos **1,258.68, 1,258.68, 1,258.68 y 1,258.70** (`model.json: despensas[1].pagos[0..3].monto`). Con HALF_UP saldría 1,258.69 ×3 + 1,258.67, que no coincide. Partes: Karen 374.93 en cada pago y Axel 883.75 ×3 + 883.77 (`despensas[1].pagos[*].axel`). Capital de la despensa 1: 3,094.50 / 4 = 773.625 → 773.63 ×3 + 773.61 (`despensas[0].pagos[*].capital_repuesto`), con intereses de 200.12 ×3 + 200.14. Referencia contra serie: con 5,000.00 el pago de referencia es 1,573.36 (`interes_por_monto[6].pago_quincenal`), mientras que la serie real sería 1,573.35 ×3 + 1,573.37.
**Esquema.** `scheduled_payment` (CHECK monto = capital + intereses), `scheduled_payment_share`. Los motores guardan la serie ya calculada (C-08).

## RN-08 Tasa por compra con IVA incluido y tabla editable por plazo · FIRME
**Definición.**
- **Compra real:** los intereses son los que Kueski cotiza al pagar, que ya incluyen IVA, y se guardan tal cual en `credit_purchase.interest_cents` (`interest_source = quoted`). Su tasa es la fracción intereses / capital (`rate_basis = term_total`).
- **Compra planeada o escenario:** intereses = HALF_UP(capital × `rate_num` / `rate_den`) con la tasa del plazo en `credit_term_rate`, o con la de la promoción. La tabla es editable y tiene una tasa vigente por plazo.
- Las tasas se guardan como fracción exacta (C-10). Kueski no publica una tabla por plazo (plan §1.6); los dos puntos conocidos salen de cotizaciones reales del 10-oct-2026: 2 quincenas = 53450/309450 (17.27 %) y 4 quincenas = 80050/309450 (25.8685 %).
- **Intereses por persona:** proporcionales a su capital, con HALF_UP, y el resto al titular.

**Ejemplo.** Tasa de 4 quincenas = 800.50 / 3,094.50 = 0.258684763… (`model.json: parametros.tasa`). Para 4,000.00: HALF_UP(4,000.00 × 80050 / 309450) = **1,034.74** (`interes_por_monto[5].intereses`, `despensas[1].intereses`). Para 1,000.00 da 258.68 (`interes_por_monto[2].intereses`); con 25.8685 % tecleado saldría 258.69, y con 5,000.00 saldría 1,293.43 en vez de 1,293.42 (`interes_por_monto[6].intereses`). Reparto de los 800.50: Karen 308.22 y Axel 492.28 (`despensas[0].karen.intereses`, `despensas[0].axel.intereses`). A 2 quincenas: intereses 534.50 sobre 3,094.50, total 3,629.00 (`data.json: financiamiento.dos_quincenas`).
**Esquema.** `credit_term_rate`; `credit_purchase` (`interest_cents`, `interest_source`, `rate_num`/`rate_den`/`rate_basis`, `interest_includes_vat`); `credit_purchase_share`.

## RN-09 Disponible de la línea: lectura A/B, B por defecto · FIRME
**Definición.**
- Solo el capital usa y repone la línea. Los intereses nunca regresan, y una compra de contado no la mueve (`model.json: credito.nota`).
- **Lectura B** (por defecto, la conservadora): el límite es el tope total. El disponible inicial es el límite menos el capital vivo. Cada compra resta su capital y cada pago suma su capital repuesto, sin pasar del límite.
- **Lectura A:** el límite está libre en la fecha ancla (`reading_a_anchor_date`; lo comprado hasta ese día ya cuenta dentro) y los pagos suman encima, sin tope. El disponible inicial es el límite; las compras posteriores al ancla restan su capital y los pagos suman su capital repuesto.
- El interruptor vive en Ajustes (`credit.availability_reading`). Mientras `availability_confirmed = 0` se muestra «pendiente confirmar con Kueski».

**Ejemplo.** El 10-oct-2026, tras la despensa 1: A = 10,000.00 y B = 10,000.00 − 3,094.50 = 6,905.50 (`model.json: credito.eventos[0].disponible_A` / `disponible_B`). Tras el pago del 30-oct: A = 10,773.63 y B = 7,679.13 (`credito.comparacion_series.esperado_A[0]` / `esperado_B[0]`). Mínimos del horizonte: A 8,094.50 y **B 5,000.00** (`credito.minimo_A`, `credito.minimo_B`). En las 14 cifras de cada serie, A − B = 3,094.50.
**Esquema.** `credit` (`credit_limit_cents`, `availability_reading` con default `B`, `availability_confirmed`, `reading_a_anchor_date`), `credit_purchase.principal_cents`, `scheduled_payment.principal_cents`.

## RN-10 Promociones manuales · FIRME
**Definición.** Las promociones (casi todas a 0 % y 4 quincenas) se capturan a mano en el calendario de promociones: comercio, crédito, plazo, tasa (0/1 = 0 %), compra mínima y vigencia (inclusive). Una promoción aplica a una compra si coincide el comercio, la fecha cae en la vigencia, el importe llega a la compra mínima y el plazo es el de la promoción. Con promoción, los intereses salen de su tasa (`interest_source = promotion`) y la serie sigue RN-07.
**Ejemplo.** En el tablero oficial de 2026, Chedraui da 0 % a 4 quincenas solo el 30-oct-2026, con compra mínima de $2,000 a $2,500 (plan §1.6; Walmart, Bodega Aurrera y Sam's no aparecen en el tablero). Ilustrativo con cifras del oráculo: si una compra de 3,094.50 tuviera 0 % a 4 quincenas, los intereses serían 0.00 en lugar de 800.50 (`model.json: despensas[0].intereses`) y los pagos 773.62 ×3 + 773.64 (RN-07).
**Esquema.** `promotion`, `credit_purchase.promotion_id`.

## RN-11 Comparador de escenarios: restricciones duras y objetivos · FIRME en definición; valores y alcance PROVISIONAL
**Definición.**
- Para una compra planeada (`planned_item` con crédito) se generan escenarios: crédito × plazo × fecha ± días × promoción, por búsqueda exhaustiva y sin solver externo.
- **Restricciones duras** (un escenario que viola una se descarta):
  1. **Colchón mínimo:** el disponible nunca baja del colchón en ninguna fecha del horizonte.
  2. **Máximo de compras vivas a la vez:** una compra está viva desde su fecha hasta su último pago (RN-21).
  3. **Deuda en cero antes de una fecha:** después de esa fecha no queda ningún pago programado.

  El tope de pago por persona **no** es restricción dura.
- **Objetivos** (se muestran los 4 óptimos lado a lado): menores intereses; pago por quincena parejo y bajo el tope; máximo disponible; comprar más por peso (tienda más barata y promociones). La recomendación sale de un puntaje ponderado con pesos iguales normalizados (`app_setting optimizer.weights`).
- Los valores de las restricciones los da Axel. Si son globales, por crédito o por compra, y si el colchón se mide en la línea usada o en la suma de líneas, queda PROVISIONAL: bloqueo `S1-restricciones-optimizador`. Cómo se mide cada objetivo y cómo se normalizan antes de ponderar está en el bloqueo `criterios-optimizador` (PROVISIONAL). El tiempo máximo de la búsqueda está en `criterios-umbrales-ausentes`.

**Ejemplo** (los valores de las restricciones son hipotéticos; las fechas y cifras son del oráculo). Con un colchón de 5,500.00 en lectura B, el plan del oráculo se descarta porque el disponible baja a 5,226.39 el 10-dic-2026 y a 5,000.00 el 10-feb-2027 y el 10-abr-2027 (`model.json: credito.eventos`, `credito.minimo_B`); con un colchón de 5,000.00 pasa. Con un máximo de 1 compra viva, comprar la despensa 2 el 10-dic-2026 viola la regla, porque la despensa 1 sigue viva hasta su pago del 15-dic-2026 (`despensas[0].pagos[3].fecha`); moverla al 16-dic-2026 la cumple. Con «deuda en cero antes del 1-may-2027», la despensa 4 se descarta porque paga hasta el 15-jun-2027 (`despensas[3].pagos[3].fecha`).
**Esquema.** `planned_item`, `credit`, `credit_term_rate`, `promotion`, `credit_purchase.status`, `person.fortnight_payment_cap_cents`, `app_setting` (`optimizer.*`).

## RN-12 Frescura de precios: 30 días · FIRME el umbral; peso PROVISIONAL
**Definición.** La edad de un precio son los días entre `observed_date` y hoy (calendario de America/Mexico_City). Un precio es **viejo** cuando su edad es de 30 días o más (`app_setting price.stale_after_days = 30`). Se marca en pantalla y pesa menos en la comparativa; cuánto menos no está decidido (PROVISIONAL, bloqueo `S1-frescura-peso`, que también pide confirmar el «30 o más» y coincide con `criterios-frescura-peso`). La comparativa usa el último precio regular (`is_promotional = 0`) de cada presentación en cada tienda.
**Ejemplo.** Los precios del ticket del 10-oct-2026 (`data.json: fecha = 2026-10-10`) están frescos hasta el 8-nov-2026 (29 días) y son viejos desde el 9-nov-2026 (30 días).
**Esquema.** `price.observed_date`, `is_promotional`, índice `ix_price_latest`; `app_setting price.stale_after_days`.

## RN-13 Membresía de Sam's: punto de equilibrio · FIRME
**Definición.** Hoy no tienen membresía (`store.has_membership = 0`, `requires_membership = 1`), y los precios de Sam's se comparan como si la tuvieran. El ahorro anual es la suma, en 12 meses, de (precio en la mejor otra tienda − precio en Sam's) × cantidad comprada, solo en los productos donde Sam's es más barato y con las cantidades del plan de compras. La membresía conviene si el ahorro anual es igual o mayor que su costo anual (`membership_annual_fee_cents`, que Axel captura en Ajustes). El punto de equilibrio es el gasto anual en Sam's en el que el ahorro iguala ese costo.
**Ejemplo ilustrativo** (el reporte no compara contra Sam's ni trae el costo de la membresía). Entre las dos tiendas del reporte, la canasta más barata ahorra 138.50 por despensa (`data.json: totales.walmart` 3,233.00 − `totales.bodega` 3,094.50). Con 6 despensas al año (bimestral) son 831.00. Si Sam's ahorrara eso frente a la mejor otra tienda, una membresía de hasta 831.00 al año se pagaría sola.
**Esquema.** `store` (`requires_membership`, `has_membership`, `membership_annual_fee_cents`), `price`, `presentation`.

## RN-14 Tiempos de comida (4) y escalado de porciones · FIRME
**Definición.** El día tiene 4 tiempos: desayuno, comida, cena y colación (cuarto tiempo opcional, con costo y calorías como los demás). Cada renglón del plan es día × tiempo × persona × receta × porciones, o «fuera». Solo aplica a Axel y Ximena (disparador). Una receta puede servir para varios tiempos (`recipe_meal_slot`, U1 N-122, N-135 y N-136). Las recetas se escriben para 2 porciones y se escalan: cantidad = cantidad base × porciones / `base_servings`, exacta en milésimas.
**Ejemplo ilustrativo** (sin datos en el reporte ni en el oráculo). Una receta para 2 porciones con 250 g de arroz (`quantity_milli` 250000) servida para 3 porciones lleva 375 g (375000).
**Esquema.** `meal_slot` (4 renglones de referencia), `recipe_meal_slot`, `meal_plan_entry`, `recipe.base_servings`, `recipe_ingredient.quantity_milli`.

## RN-15 Inventario híbrido · FIRME
**Definición.**
- **Existencia por unidad** (`stock`), para básicos y no perecederos: la existencia es la del último conteo, más compras, menos consumos.
- **Ritmo de consumo** (`consumption_rate`), para el resto: existencia estimada = existencia contada − HALF_UP(cantidad del ritmo × días desde el conteo / días del ritmo), nunca menor que 0.
- Los **sobrantes** de recetas entran como existencia por unidad, contados en porciones (`inventory_item` con `recipe_id`).
- Cada entrada o salida queda en `inventory_movement` con su tipo (compra, sobrante de receta, consumo del plan o ajuste) y, en los ajustes, su motivo (conteo, consumo o merma), en la misma transacción que cambia la existencia (U1 N-181 y N-183).
- La lista de compras es la suma de ingredientes del plan menos la existencia más la reposición por ritmo (S4). Se guarda por semana (`shopping_list`) solo lo que no se calcula: artículos agregados a mano, tienda o cantidad elegidas y marcas de comprado (`shopping_list_item`). Regenerar la lista actualiza la necesidad calculada y conserva lo elegido. **PROVISIONAL**: que las marcas se compartan entre los dos teléfonos es la opción A del bloqueo `C1-lista-compras`, pendiente de Axel en la puerta 1.

**Ejemplo** (la cantidad es del reporte; el ritmo es ilustrativo). La despensa 1 trae 3 botellas de aceite de soya de 946 ml (`data.json: cuentas.conjunto.lineas[1]`, cant 3) → 2,838 ml (2838000 milésimas). Con un ritmo de una botella cada 21 días, a los 30 días la existencia estimada sería 2,838 − 1,351.429 = 1,486.571 ml.
**Esquema.** `inventory_item`, `inventory_movement`, `shopping_list`, `shopping_list_item` (PROVISIONAL).

## RN-16 Mercados con día de visita · FIRME
**Definición.** Hay dos verticales de compra. Verduras y carne no especializada se compran en mercados o tianguis: `store` con `kind = market`, cada uno con su `visit_weekday` (1 = lunes … 7 = domingo) y los grupos que se compran ahí en `store_product_group`. El resto se compra en Walmart, Bodega Aurrera y Sam's Club (`kind = chain` o `branch`). La vertical de cada producto la da su grupo (`product_group.purchase_channel`). Un renglón de la lista con vertical mercado se asigna a un mercado que surte su grupo, en su día de visita.
**Ejemplo ilustrativo** (Axel aún no da de alta sus mercados). Para un tianguis con `visit_weekday` = 2 (martes), después del sábado 10-oct-2026 la siguiente visita es el martes 13-oct-2026.
**Esquema.** `store.visit_weekday` (obligatorio en mercados), `product_group`, `store_product_group`, `product.product_group_id`.

## RN-17 Salud: IMC y meta de calorías · IMC FIRME; fórmula de meta PROVISIONAL
**Definición.**
- **IMC** = peso (kg) / estatura (m)², con el último peso registrado. En enteros: `weight_g` × 1000 / `height_mm`², mostrado con un decimal (HALF_UP).
- **Meta de calorías por día:** la captura cada persona (`daily_kcal_target`). Si además el sistema la sugiere con una fórmula, cuál fórmula y con qué escala de actividad queda PROVISIONAL (bloqueo `S1-meta-kcal`, que coincide con `criterios-meta-kcal`).
- **Calorías del día** = suma de las recetas planeadas (kcal por porción × porciones); los tiempos «fuera» no suman. Se comparan contra la meta en porcentaje.
- **Kcal de una receta** = suma, por ingrediente, de gramos × kcal por 100 g / 100. Los gramos salen de la cantidad en unidad base × `grams_per_unit_milli` cuando la unidad es ml o pieza. Cada ingrediente cita su fuente.
- Los macronutrientes tienen columnas previstas en `nutrition_fact` y se llenan después.

**Ejemplo ilustrativo** (el reporte y el oráculo no traen datos de salud). 72.5 kg y 1.75 m → 72500 × 1000 / 1750² = 23.67 → IMC **23.7**. Un día con 1,850 kcal planeadas contra una meta de 2,000 kcal = 92.5 %.
**Esquema.** `health_profile`, `weight_entry`, `nutrition_fact`, `product.grams_per_unit_milli`.

## RN-18 Fondo y metas de ahorro · FIRME
**Definición.** Hay fondos (sin fecha obligatoria) y metas con monto y fecha (`savings_goal`). Avance = aportaciones − retiros (`savings_movement`). Porcentaje = avance / monto objetivo, con un decimal. La aportación mensual planificada (`planned_monthly_cents`) se reserva en el plan del mes y el dashboard la muestra junto al avance.
**Ejemplo ilustrativo** (con una cifra del oráculo como objetivo). Una meta de 6,293.44 (el monto de `model.json: horizonte.deuda_viva_10_abr`) para el 10-abr-2027, con aportaciones por 3,146.72, va en 50.0 %.
**Esquema.** `savings_goal` (CHECK: una meta lleva monto y fecha), `savings_movement`.

## RN-19 Comidas fuera de casa · FIRME
**Definición.** Comer fuera es un gasto normal en la categoría «comida fuera» (`category.code = eating_out`), con la persona que comió (`expense.person_id`). En el plan de comidas ese tiempo se marca «fuera» (`is_eating_out = 1`): sin receta, sin porciones y sin calorías (CHECK), y no suma al costo de recetas ni a las calorías del día.
**Ejemplo ilustrativo.** Ximena come fuera el lunes 12-oct-2026: se registra un gasto de comida fuera con `person_id` = Ximena y un renglón del plan (12-oct, comida, Ximena, fuera). Ese día sus calorías solo cuentan desayuno, cena y colación.
**Esquema.** `expense` (`category_id`, `person_id`), `meal_plan_entry.is_eating_out`.

## RN-20 Equivalencia por unidad base y sustitutos · FIRME
**Definición.** Cada producto tiene unidad base (g, ml o pieza) y cada presentación su contenido en milésimas de esa unidad. Precio por unidad base = precio / contenido; se compara con fracciones exactas y se muestra por kilo, litro o pieza con HALF_UP a centavos. Dos presentaciones del mismo producto son **sustitutas** y en la comparativa se marcan así; «idéntico» significa la misma presentación.
**Ejemplo.** Arroz: 3 × 750 g por 66.00 = $29.33/kg contra 3 × 1 kg por 84.00 = $28.00/kg (`data.json: sustitutos[0].nota` = «$29.33/kg vs $28.00/kg: Walmart más barato»).
**Esquema.** `product.base_unit`, `presentation.size_milli`, `price`.

## RN-21 Deuda viva y compras vivas · FIRME
**Definición.** La deuda viva a una fecha es la suma de los pagos programados posteriores a esa fecha de las compras no canceladas. Una compra está viva mientras tiene algún pago pendiente (`credit_purchase.status = active`). Su relación con devengado y efectivo (el invariante de S2) está en el bloqueo `criterios-invariante-deuda`. Aquí «deuda viva» es la deuda viva total, 6,293.44 al 10-abr-2027, el nombre que usan las opciones a y b de ese bloqueo.
**Ejemplo.** Al 10-abr-2027 la deuda viva es **6,293.44** = 1,258.70 pendiente de la despensa 3 + 5,034.74 de la despensa 4 (`model.json: horizonte.deuda_viva_10_abr`, `horizonte.deuda_viva_detalle`).
**Esquema.** `scheduled_payment.due_date`, `credit_purchase.status`, `credit_payment`.

## RN-22 Presupuesto, real contra plan y alertas 80/100 % · FIRME
**Definición.** El presupuesto es mensual, por categoría y total (`budget`), con vigencia: un monto rige desde un mes (`valid_from_date`) hasta otro inclusive (`valid_to_date`) o sin fin (NULL), y el monto de un mes es el del renglón cuya vigencia lo contiene. Un monto de un solo mes es un renglón que empieza y termina en ese mes. Si un cambio hecho en Ajustes rige solo para ese mes o para los siguientes hasta otro cambio es **PROVISIONAL** (bloqueo `S1-huecos-negocio`); el esquema admite las dos lecturas. El real del mes es el devengado de la parte del hogar en cada categoría (RN-02). Hay alerta al llegar al 80 % y al 100 % (mayor o igual), por categoría y en el total (`notification` `budget_threshold`).
**Ejemplo ilustrativo** (el oráculo no trae montos de presupuesto). Con un presupuesto de despensa de 2,000.00 en dic-2026, el devengado del hogar es 1,767.51 (`model.json: periodos[2].devengado.axel`) = 88.4 %, así que salta la alerta del 80 %.
**Esquema.** `budget` (`valid_from_date`, `valid_to_date`), `expense`, `expense_share`, `notification`.
