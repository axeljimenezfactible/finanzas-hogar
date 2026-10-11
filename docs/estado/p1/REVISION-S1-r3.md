# Revisión S1 · esquema 1.2 · ronda 3 de 3 (flujo P1)

Fecha: 2026-10-11 · Agente: revisor-S1-r3 (independiente: no escribí nada de lo revisado y no corrijo nada).
Base de evaluación: `docs/criterios/S1.md` versión 1.1 y `docs/plan/decisiones.md` §K (K1–K3, K6, K8, K9, K11–K16, K18–K23).
Objeto: `backend/db/migrations/00002_puerta1.sql` (leída entera), comentarios cambiados de `00001_esquema.sql` (`git diff HEAD`), `backend/db/tools/` (`s1_verifica.py`, `gen_dict.py`), `docs/contratos/diccionario-datos.md` 1.2 y `docs/contratos/reglas-negocio.md` 1.2 (por secciones), y el reporte `docs/estado/p1/S1-esquema-1.2.md`.
No aprobé por lo que dicen los reportes: cada verificación la reejecuté con mis propios comandos y scripts en `datos/tmp/` (ya borrados).

## Veredicto
**APROBADO.** 0 hallazgos altos y 0 medios.
- Hay 9 bajos: 1 nuevo (R3-01, sensibilidad del verificador) y 8 que siguen abiertos desde r2 (R2-02 a R2-09, sin cambio).
- El hallazgo alto de r2 (R2-01, S1-03) queda **cerrado**.

## Verificaciones del encargo
| # | Qué se pidió | Cómo lo comprobé | Resultado |
|---|---|---|---|
| a | 00001 + 00002 en `datos/tmp/rev.db` | goose v3.28.0 `up`; luego `PRAGMA integrity_check`, `foreign_key_check` y `PRAGMA table_list` | OK: `ok`, 0 filas, 55 tablas STRICT. La única tabla no STRICT es `goose_db_version`, que crea goose |
| b | Down de 00002 = 00001 sola | b1 = `up-to 1`; b2 = `up` y luego `down`; comparé `sqlite_master` (type, name, tbl_name, sql) | OK: igual en orden de rowid y ordenado (232 objetos). Los datos de todas las tablas son iguales sin `created_at`/`updated_at` (app_setting 7, meal_slot 4, notification 7). El `iterdump` completo solo difiere en la hora de las semillas, que sale de cada corrida |
| c | goose up / down-to 0 / up | goose v3.28.0 en `datos/tmp/c.db` | OK: códigos 0, 0, 0; versiones 2, 0, 2; tablas 55, 0, 55; `integrity_check` ok; esquema igual al de (a) |
| d | DDL de 00001 sin cambios | Comparé `git show HEAD:…/00001_esquema.sql` con el archivo actual. Un lexer propio quita `--` y `/* */` sin tocar los literales | OK: Up y Down iguales sin comentarios, anotaciones `+goose` iguales, `sqlite_master` igual sin comentarios y semillas iguales (sin marcas de tiempo). Las 20 líneas `+`/`−` de `git diff -U0` son todas de comentario |
| e | `s1_verifica.py` | `GOOSE=<goose v3.28.0> python3 -I backend/db/tools/s1_verifica.py` | OK: rc 0, «RESUMEN TODO OK». KUESKI: suma de pagos de la despensa 1 = **389500** y parte de Axel = **239528**. NEG: 46/46 rechazadas y 8/8 positivas. NEG12: 26/26 y 12/12. S1-12: 164 valores con 0 diferencias. DOWN: 232 objetos. GOOSE: [0, 0, 0] |
| f | `gen_dict.py` regenera el diccionario | `python3 -I backend/db/tools/gen_dict.py datos/tmp/dic-regen.md` y `cmp` | OK: idéntico (174,605 bytes). 55 tablas, 152 índices, 71 disparadores, 0 columnas sin descripción |
| g | Cada K aplicada, con línea | `grep` de `\bK<n>` (sin dígito después) en 00001, 00002, diccionario y reglas, y lectura de cada línea | OK: tabla «Trazabilidad K» abajo |
| h | Ningún PROVISIONAL cita un bloqueo resuelto | Script propio: toma cada línea con la marca, extrae los bloqueos citados y busca «Respuesta de Axel» en cada uno | OK: 00001 tiene 0, 00002 tiene 0, el diccionario 1 (L19, la definición de la marca) y las reglas 11. En las reglas, L9 es la definición y las otras 10 citan solo bloqueos abiertos: `S1-ajuste-objetivo-kcal`, `S1-fuentes-precio`, `criterios-frescura-extremos` y `criterios-promedios-redondeo`. Los cuatro tienen el formato del protocolo §2 y no tienen respuesta. 0 marcas citan un bloqueo resuelto. Ningún valor de K31 aparece en el esquema ni en el contrato |
| i | Ninguna cifra inventada en semillas nuevas | Revisé los INSERT de la 00002 y `app_setting` en la base | OK: `projection.repeat_cadence_enabled` = false (K21), `session.max_age_days` = 90 (K11) e `idempotency.retention_hours` = 24 (K6). Las 3 claves `optimizer.*` de restricciones no están en ningún INSERT ni en la base (K16). El diccionario las lista en «Claves de app_setting previstas sin valor». Las semillas de 00001 son iguales a las de HEAD |
| j | Reglas 1.2: ejemplos exactos del oráculo o rotulados | Escáner propio de los 28 ejemplos, más 4 mutaciones del verificador en una copia temporal | OK con bajas: detalle abajo |

**Detalle de (j).**
- **Escáner propio.** Recorre 95 frases con 62 fechas y 170 cifras de dos decimales. Las compara con las 29 fechas y las 248 cifras que hay en `model.json` y `data.json`.
  - Fuera del oráculo y en una frase sin la palabra «ilustrativo» solo queda 1,814.50 de RN-11, una cuenta citada (R2-02).
  - Las otras banderas no son defectos: frases que empiezan con «Ilustrativo» en mayúscula (RN-07 y RN-10, que siguen en R2-09) y un «30-oct» de RN-09 sin año, que es la fecha 2026-10-30 del oráculo.
  - 24 citas con ruta se resolvieron contra el oráculo sin ninguna discrepancia. La de RN-02 parecía fallar porque mi script corta la frase en el «:» de la celda de la tabla, pero 973.75 está en esa misma celda.
- **Hechos de RN-06 y RN-11, recalculados:**
  - el 10-oct-2026 es sábado, el 5-oct-2026 es lunes y el 11-oct-2026 es domingo;
  - del 10-oct al 15-oct hay 5 días y al 30-oct 20; del 10-feb al 28-feb-2027 hay 18;
  - `despensas[0].pagos` caen el 2026-10-30, 2026-11-15, 2026-11-30 y 2026-12-15; `despensas[1].fecha` es 2026-12-10 y `despensas[3].pagos[3].fecha` es 2027-06-15;
  - `disponible_B` es 5,226.39 el 2026-12-10 y 5,000.00 el 2027-02-10 y el 2027-04-10;
  - (800.50 − 667.50) / (800.50 − 534.50) = 1/2.
- **Mutaciones del verificador:**
  - quitar el rótulo de RN-06 da FALLA en S1-03 (2026-10-05, 2026-10-11 y 2026-10-31);
  - quitar el de RN-11 da FALLA (2026-12-16);
  - cambiar 6,293.44 por 6,293.45 en RN-21 da FALLA;
  - cambiar 5,226.39 por 5,226.40 dentro de la frase rotulada de RN-11 da **OK** (es R3-01).

## Criterios S1 (versión 1.1)
| ID | Sev. | Estado | Evidencia (propia) |
|---|---|---|---|
| S1-01 | alta | Cumple | Las 31 entidades de §3.4.1 tienen fila en el mapa del diccionario y su tabla existe (por ejemplo, ingreso → `planned_item` + `income` y respaldo → `backup`). Las 55 tablas tienen fila en el mapa (45 de entidades y 10 complementarias) y su sección `###` |
| S1-02 | alta | Cumple | Verificaciones (a) y (c) |
| S1-03 | alta | Cumple | Verificación (j) y las mutaciones. R2-01 cerrado. Quedan bajas: R2-02, R2-09 y R3-01 |
| S1-04 | alta | Cumple | 55 de 55 secciones tienen «Usado por» lleno. Comprobé 5 al azar (semilla 20261011) contra U1 y los motores: `notification_subscription` (Ajustes «Notificaciones», 3 veces en U1), `inventory_item` («Inventario», 9), `income` («Ingresos», 10), `expense_share` (filtros persona y familia de C1; S2 y S7) y `backup` («Datos y respaldos», 2; S6 y S8). Las 5 tienen consumidor real |
| S1-05 | alta | Cumple | 0 columnas REAL, FLOAT o DOUBLE. Las 37 columnas `*_cents` son INTEGER: una sola convención |
| S1-06 | media | Cumple | 0 fechas en formato de pantalla en datos o en SQL. El grep literal encuentra 1, en un comentario de 00002 L3 (R2-08, baja) |
| S1-07 | alta | Cumple | Las cuentas de Axel y Ximena se aceptan. La de Karen la rechaza `trg_user_account_not_external_ins`. También se rechazan pasar a externa a Ximena, que tiene cuenta, y un hash sin `$argon2id$`. `user_account` no tiene columna de rol |
| S1-08 | alta | Cumple (pendiente de S6) | La base acepta partes de 60 % + 50 %. El rechazo está documentado (invariante 8 del diccionario, con la prueba comprometida para S6), como admite el criterio. La despensa 1 cabe sin pérdida (S1-12) |
| S1-09 | alta | Cumple | Se recuperan 18 filas (6 meses y 12 quincenas). Ninguna quincena sale de los límites de K8. La quincena 2 de oct va de 2026-10-16 a 2026-10-31 y paga el 2026-10-30; la de feb va de 2027-02-16 a 2027-02-28 y paga el 2027-02-28. La semana del miércoles 2026-10-14 va de 2026-10-12 a 2026-10-18. El CHECK rechaza el mes 2026-10-10 a 2026-11-09 y una quincena 2 que termina el 30. Equivalencia calculada por mí: 12 de 12 cifras contra `periodos[0..5]` y ningún evento fechado entre los días 1 y 9 |
| S1-10 | alta | Cumple | Hay 4 tipos de crédito. Las fracciones 1069/6189 y 1601/6189 (= 53450/309450 y 80050/309450), con `rate_decimals` 2 y 4, se muestran como 17.27 y 25.8685. La lectura por defecto es B y `availability_confirmed` = 0. `rate_base` se rechaza en `revolving_line` |
| S1-11 | media | Cumple | El DDL de 00001 no cambió y la 00002 suma `health_profile.goal`. La línea S1-11 de `s1_verifica` da 12 de 12 puntos |
| S1-12 | alta | Cumple | `s1_verifica`, corrido por mí: 4 despensas, 16 pagos, 17 eventos, series A/B de 14/14, 164 valores y 0 diferencias |
| S1-13 | media | Cumple | Las 14 secciones tienen destino. Los 60 destinos existen como columna o como ajuste sembrado. Los únicos que no son columna son el CSV de gastos fijos y las 3 claves `optimizer.*` que K16 deja sin semilla, y están documentadas |
| S1-14 | alta | Cumple | `password_hash` tiene CHECK argon2id. En las migraciones el grep solo encuentra nombres de columna y texto: `password_hash`, `password_changed_at`, `token_hash` y la palabra «token» en comentarios. El patrón `[0-9]{8,}:AA` no aparece en `backend/db`. El diccionario declara K11 en L2270: `login_attempt`, 5 intentos, ventana deslizante de 15 min por usuario e IP y sesión de 90 días renovable con el uso |
| S1-15 | media | Cumple | `foreign_key_list`: expense 6, credit_purchase 3, scheduled_payment 1, price 4 y telegram_chat_candidate 1. Índices de los 6 filtros: `ix_expense_person_date`, `ix_expense_share_person`, `ix_expense_category_date`, `ix_expense_store_date`, `ix_expense_date` e `ix_credit_purchase_credit_status` |
| S1-16 | alta | Cumple | Verificaciones (h) e (i). R:36 y R:216 (el título de RN-17) contienen «IMC» y la marca, pero la marca es del ajuste por objetivo de K23 (bloqueo abierto); la línea del IMC (L218) no la lleva. El reporte tiene SUPUESTOS. Sigue R2-03 (baja) |
| S1-17 | media | Cumple | El glosario tiene 37 filas, sin códigos ni términos repetidos. Las 55 tablas están en snake_case ASCII. Sigue R2-07 (baja) |

## Hallazgos de esta ronda
| ID | Sev. | Criterio | Ruta | Hallazgo | Corrección sugerida |
|---|---|---|---|---|---|
| R3-01 | baja | S1-03 (Verifica) | `backend/db/tools/s1_verifica.py` L330-339 (bloque S1-03, L238-364) | El verificador no compara con el oráculo las cifras ni las fechas de una frase rotulada «ilustrativo», aunque lleven ruta. La mutación 5,226.39 → 5,226.40 en RN-11 L177 (la frase del colchón ilustrativo, con ruta `credito.eventos`) da S1-03 OK y RESUMEN TODO OK. Las cifras actuales están bien: las comprobé a mano contra `credito.eventos` | En las frases rotuladas, comparar también las cifras que llevan ruta (buscando el valor en la lista citada) o que están en el oráculo, y contar como ilustrativas solo las demás. Sumar esta mutación a la autoprueba |

## Siguen abiertos desde r2 (todos bajos)
| ID | Criterio | Ruta | Estado hoy | Corrección sugerida |
|---|---|---|---|---|
| R2-02 | S1-03 | reglas L177 (RN-11) | 1,814.50 no está en el oráculo. Pasa como cuenta citada (`financiamiento.dos_quincenas.total / 2`). La cuenta es exacta: 3,629.00 / 2 = 1,115.85 + 698.65 | Citar 1,115.85 + 698.65 (`data.json: financiamiento.dos_quincenas.axel.por_quincena` y `.karen.por_quincena`), o que criterios 1.2 admita la cuenta citada |
| R2-03 | S1-16 | `S1-esquema-1.2.md` L35-36 | En Supuestos siguen dos puntos de esquema: K12 en `credit.annual_rate_decimals` y `promotion.rate_decimals`, y las columnas NULL por ALTER TABLE. Los respaldan K12 y la restricción técnica | Pasarlos a «Hecho» |
| R2-04 | Evidencia | `S1-esquema-1.2.md` L14, L16, L17, L26, L30, L31 | Cifras de antes de las correcciones: 13 líneas, 24 negativas y 10 positivas, 69 disparadores, 173,798 bytes y S1-03 sin los conteos de ilustrativos ni de fechas. Hoy son 20 líneas `+`/`−` de comentario, NEG12 26 y 12, 71 disparadores y 174,605 bytes; los conteos están en `correccion-S1-r1.md` y `correccion-S1-r2.md` | El integrador actualiza esas líneas o enlaza los dos reportes de corrección |
| R2-05 | K18 | 00002 L112-115; diccionario L2301 | La base no exige que los datos de K18 vayan juntos, y `vat_on_interest_bp` NULL quiere decir «ya trae IVA». Lo cubre el invariante 16 en la API | Disparador de todo o nada, o columna explícita de IVA, con su prueba |
| R2-06 | K22 | reglas L183 (RN-12) | Sigue «la que S1 recomienda» (la lectura b) dentro del contrato. El valor provisional en uso es la lectura (a), así que ninguna cifra está mal | Dejar la recomendación solo en el bloqueo |
| R2-07 | S1-17 | 00002 L112 frente a 00001 L756 | `credit.rate_base` (daily/monthly) convive con `credit_purchase.rate_basis` (term_total/annual), que significa otra cosa | Renombrar antes del commit (Up, Down, reglas y glosario) y regenerar el diccionario |
| R2-08 | S1-06 | 00002 L3 | La cabecera dice «11-oct-2026», que es formato de pantalla | Escribir 2026-10-11 |
| R2-09 | S1-03 | reglas | 10 ejemplos completos no traen la leyenda literal «(sin dato en el oráculo)»: L192 RN-13, L197 RN-14, L213 RN-16, L224 RN-17, L229 RN-18, L234 RN-19, L249 RN-22, L254 RN-23, L259 RN-24 y L269 RN-26. Tampoco la traen las frases con «ilustrativo» de L110 RN-04, L136 RN-07, L161 RN-10, L177 RN-11 («tercer escenario») y L279 RN-28. Todas dicen «ilustrativo» y ninguna se presenta como canon | Usar la leyenda literal y poner el motivo después |

## Estado de los hallazgos de la ronda 2
- **R2-01 (alta, S1-03, verificación j): CERRADO.**
  - **RN-06, L124:** la frase dice ahora «Ejemplo ilustrativo (sin dato en el oráculo) de los periodos (K8): el sábado 10-oct-2026 cae en la semana del lunes 5-oct-2026 al domingo 11-oct-2026, y el 31-oct-2026 en la quincena 2 de octubre, que paga el 30-oct-2026». El cálculo del primer pago (15-oct-2026) va en una frase aparte con la misma leyenda. Los pagos siguen citados (`despensas[0].pagos[0..3].fecha`, comprobado).
  - **RN-11, L177:** el encabezado dice ahora «las fechas de compras y pagos y las demás cifras son del oráculo… salvo lo rotulado como ilustrativo». El escenario del 16-dic-2026 y la fecha «deuda en cero antes del 1-may-2027» llevan la leyenda literal.
  - Las fechas están bien calculadas.
  - El verificador ya revisa fechas en tres formatos (ISO, dd-mmm-aaaa y dd-mmm) y su autoprueba pasa (5 casos). Mis mutaciones que quitan cada rótulo dan FALLA.
- **R2-02 a R2-09 (bajas):** siguen abiertas y sin cambio. Quedaron fuera del encargo de `correccion-S1-r2`.

## Trazabilidad K (verificación g)
Abreviaturas: 00001 = M1, 00002 = M2, diccionario = D, reglas = R.

| K | Esquema | Diccionario | Reglas |
|---|---|---|---|
| K1 | M1 L1106-1107 y L1122-1123 (lista compartida, sin marca) | D L82-83, L1767-1768 y L1801 | R L206 (RN-15) |
| K2 | M1 L299 (`budget.valid_to_date`) | D L530 | R L248 (RN-22) |
| K3 | M2 L29-34 | D L1125-1126 | R L253 (RN-23) |
| K6 | M1 L161 y L176 (24 h); M2 L137 y L145 | D L337, L363, L2251, L2263 y L2306 | R L64 (C-15) |
| K8 | M1 L186, L194, L200 y L1111; M2 L70-88 (lunes por disparador, probado) | D L374, L384, L392, L1774 y L1795 | R L120-121 y L124 (RN-06); L209 (RN-15) |
| K9 | M2 L91-105 | D L999-1000, L1901-1902 y L2300 | R L258 (RN-24) |
| K11 | M2 L144 (90 días); `auth.lockout_window_minutes` = 15 desde 1.1 | D L2251, L2263 y L2270 | No hay regla; S1-14 lo pide en el diccionario |
| K12 | M2 L108-127 | D L770, L814, L855, L1223 y L2274 | R L53 (C-04, «Z» solo en entradas) y L263 (RN-25) |
| K13 | M1 L39 | D L143 y L2281 | R L107 (RN-04) |
| K14 | No cambia el esquema | — | R L76 (RN-02) y L134 (RN-07) |
| K15 | — | — | R L106 (RN-04) |
| K16 | M2 L10 y L138-140 (sin semilla) | D L2258-2260 y L2311 | R L167 y L177 (RN-11) |
| K18 | M2 L108-114 | D L740, L767-769 y L2301 | R L268 (RN-26) |
| K19 | — | — | R L243 (RN-21: 6,293.44 = 1,258.70 + 5,034.74, citas exactas) |
| K20 | — | — | R L273 (RN-27) |
| K21 | M2 L143 | D L2251 y L2271 | R L278 (RN-28) |
| K22 | `price.stale_after_days` = 30 desde 1.1 | D L2263 | R L182-184 (RN-12) |
| K23 | M1 L1169 (factores de la opción B de `S1-meta-kcal`, resuelto); M2 L130-134 | D L1863, L1869 y L2282 | R L219 (RN-17) |
