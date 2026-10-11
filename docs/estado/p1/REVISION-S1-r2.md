# Revisión S1 · esquema 1.2 · ronda 2 de 3 (flujo P1)

Revisor independiente `revisor-S1-r2` · 11-oct-2026 · Evalúa solo contra `docs/criterios/S1.md` 1.1 y `docs/plan/decisiones.md` §K (K1–K3, K6, K8, K9, K11–K16, K18–K23). No corrige nada. Todas las verificaciones se reejecutaron con scripts propios en `datos/tmp/rev-s1-r2/` y la base `datos/tmp/rev.db` (borrados al terminar).

## Veredicto
**NO APROBADO.** 1 hallazgo alto (S1-03), 0 medios y 8 bajos (3 nuevos y 5 abiertos desde r1, uno de ellos en parte). El esquema, la migración 00002, el Down, goose, el diccionario regenerado y las cifras del oráculo pasan. Lo que falta es de rótulos en dos ejemplos de `reglas-negocio.md`: traen fechas que no están en el oráculo y no dicen que son ilustrativas. El verificador no lo ve porque solo revisa cifras de dos decimales.

## Verificaciones del encargo
| | Verificación | Resultado | Evidencia (comando propio y salida resumida) |
|---|---|---|---|
| a | 00001 + 00002 en `datos/tmp/rev.db` | Pasa | `sqlite3 -bail` con el Up de cada archivo: rc 0 y 0 · `integrity_check` = ok · `foreign_key_check` = 0 filas · `pragma_table_list`: 55 tablas, 55 STRICT, 0 sin STRICT · 114 índices propios + 38 automáticos = 152 · 71 disparadores |
| b | Down de 00002 = 00001 sola | Pasa | Volcado `type|name|tbl_name|sql` de `sqlite_master` antes de 00002 y después de su Down: `cmp` idéntico (1263 líneas). Además el `.dump` completo (esquema y datos, 102,696 bytes) sale idéntico. Down de 00001 deja 0 objetos |
| c | goose up / down-to 0 / up | Pasa | goose v3.28.0: up rc 0 (v2) · version 2 · down-to 0 rc 0 · version 0 · up rc 0 · version 2 · status con 00001 y 00002. `sqlite_master` de goose = el de sqlite3, salvo `goose_db_version` y `sqlite_sequence` |
| d | DDL de 00001 sin cambios contra `git show HEAD` | Pasa | Quité comentarios respetando literales y comparé tokens: HEAD 8599 = actual 8599, iguales. `git diff`: 14 inserciones y 13 borrados, todos de comentario |
| e | `python3 -I backend/db/tools/s1_verifica.py` | Pasa | rc 0, «RESUMEN TODO OK» (23 líneas). KUESKI: 389500 y 239528 (= 973.75 × 4 y 598.82 × 4). NEG 46/46 rechazadas, NEG12 26/26 rechazadas y 12/12 positivas aceptadas. S1-12: 164 valores, 0 diferencias |
| f | `gen_dict.py` regenera el diccionario | Pasa | `cmp` contra `docs/contratos/diccionario-datos.md`: idéntico (174,605 bytes); «columnas sin descripción: []» |
| g | Cada K aplicada, con línea | Pasa | Tabla «Trazabilidad K» abajo |
| h | Ninguna PROVISIONAL cita un bloqueo resuelto | Pasa | 00001 = 0, 00002 = 0, diccionario = 1 (la definición, L19), reglas = 11 líneas que citan solo `criterios-frescura-extremos`, `S1-fuentes-precio`, `S1-ajuste-objetivo-kcal` y `criterios-promedios-redondeo`: los 4 sin «Respuesta de Axel». Ninguna cita queda de `U1-periodos-semana-quincena` |
| i | Semillas nuevas solo con valores de §K | Pasa | `app_setting` de la 00002: `projection.repeat_cadence_enabled` = false (K21), `session.max_age_days` = 90 (K11), `idempotency.retention_hours` = 24 (K6). Claves `optimizer.*`: solo `optimizer.weights` de 1.1; las 3 de K16 sin valor. Factores 1.375, 1.55 y 1.725: opción B del bloqueo `S1-meta-kcal`, que la D elegida incluye («B más un ajuste según el objetivo») |
| j | Reglas 1.2 con cifras del oráculo exactas o rotuladas | **No pasa** | Script propio: 0 citas «ruta = valor» que no coinciden; RN-12 y RN-27 llevan la leyenda literal. Pero RN-06 L122 y RN-11 L175 traen fechas fuera del oráculo sin rótulo (R2-01) |

## Criterios S1
| Criterio | Estado | Evidencia |
|---|---|---|
| S1-01 (alta) | Cumple | Script propio: 31 de 31 entidades con tabla y sección `### `nombre`` en el diccionario; 55 tablas = 55 secciones |
| S1-02 (alta) | Cumple | Verificaciones a, b, c y d |
| S1-03 (alta) | **No cumple** | R2-01 (RN-06 y RN-11). Lo demás está bien: citas exactas y RN-12 y RN-27 rotulados |
| S1-04 (alta) | Cumple | 55/55 con «Usado por». Cinco que comprobé a mano en U1: `telegram_chat_candidate` (U1 L29 y L275: P-26.09, N-299, N-316); borrador de `attachment` (U1 L28, L238, L486-487: P-05, N-071, N-075, N-076); `health_profile.goal` (fila K23 de U1: P-21, P-26.13, N-230, N-311, N-313); `credit_term_rate.rate_decimals` (fila K12: N-204, N-286, N-287, N-289); `expense.idempotency_key` y `captured_at` (U1 L290: G-07, N-317) |
| S1-05 (alta) | Cumple | 0 columnas REAL/FLOAT/DOUBLE; 37 columnas `_cents` INTEGER, con `credit.min_payment_fixed_cents` |
| S1-06 (media) | Cumple | 0 fechas de pantalla en datos; la única coincidencia literal es un comentario (00002 L3, R2-08). Los instantes nuevos exigen «Z» |
| S1-07 (alta) | Cumple | Sin cambios en 1.2. El verificador rechaza la cuenta de Karen y no hay tablas de finanzas de su familia |
| S1-08 (alta) | Cumple | Sin cambios en 1.2. 60/50 lo rechaza la API (invariante 8, prueba de S6); despensa 1 sin pérdida |
| S1-09 (alta) | Cumple | 18 filas, 12/12 contra P1–P6 y ventana del 10 al 9 rechazada. Semana de lunes a domingo: probé `shopping_list` y rechaza martes 2026-10-13, domingo 2026-10-11 y el cambio a miércoles; acepta el lunes 2026-10-12, con fin = inicio + 6 |
| S1-10 (alta) | Cumple | 4 tipos. `rate_decimals` muestra 17.27 % y 25.8685 %; 1,034.74 sale de la fracción 80050/309450 |
| S1-11 (media) | Cumple | Sin cambios; 12/12 |
| S1-12 (alta) | Cumple | 164 valores comparados, 0 diferencias |
| S1-13 (media) | Cumple | 14/14 secciones; 46 referencias, ninguna inexistente; K21 en la fila 2 |
| S1-14 (alta) | Cumple | Diccionario L2270: `login_attempt` y `auth.*` (5; 15, ventana deslizante por usuario e IP, K11) y `session.max_age_days` (90, se renueva con el uso, K11). 0 tokens de Telegram en `backend/db`; la 00002 no tiene nombres con password, token ni secret |
| S1-15 (media) | Cumple | Llaves foráneas de las 5 tablas que elegí: `telegram_chat_candidate`→person RESTRICT; `attachment`→user_account y expense SET NULL; `weight_entry` y `health_profile`→person RESTRICT; `shopping_list_item`→store, presentation, product y shopping_list. Los índices de los 6 filtros los reporta el verificador |
| S1-16 (alta) | Cumple, con nota | Verificaciones h e i. Nota baja R2-03 (supuestos del reporte S1) |
| S1-17 (media) | Cumple | 656 identificadores en `[a-z0-9_]`; 36 filas del glosario, ninguna repetida; R2-07 sigue abierto (baja) |

## Hallazgos de esta ronda
| ID | Sev. | Criterio | Ruta | Hallazgo | Corrección sugerida |
|---|---|---|---|---|---|
| R2-01 | alta | S1-03 (j) | `docs/contratos/reglas-negocio.md` L122 (RN-06) y L175 (RN-11); `backend/db/tools/s1_verifica.py` L235-313 | Dos frases de ejemplo traen fechas que no están en el oráculo y no llevan rótulo. **RN-06** (tocada en 1.2 por K8): «El sábado 10-oct-2026 cae en la semana del lunes 5-oct-2026 al domingo 11-oct-2026, y el 31-oct-2026 en la quincena 2». Ni 2026-10-05, ni 2026-10-11, ni 2026-10-31 están en `model.json` o `data.json`, y S1-03 1.1 pide que el ejemplo de los periodos de K8 sea «ejemplo ilustrativo, sin dato en el oráculo». **RN-11**: «moverla al 16-dic-2026 la cumple». 2026-12-16 no está en el oráculo ni es el valor de una restricción, y el encabezado del ejemplo dice «las fechas y las demás cifras son del oráculo, con su ruta o con su cuenta». Es la misma afirmación que r1 señaló en RN-11. El bloque S1-03 del verificador solo revisa cifras de dos decimales, así que los dos casos salen OK. Las fechas están bien calculadas: no hay cifra mal, falta el rótulo | **RN-06:** anteponer «Ilustrativo (sin dato en el oráculo):» a la frase de semana y quincena, o usar el ejemplo del criterio (31-oct-2026 en la quincena 2, del 2026-10-16 al 2026-10-31; semana del miércoles 14-oct-2026, del 2026-10-12 al 2026-10-18). **RN-11:** rotular como ilustrativo el escenario del 16-dic-2026, o decir en el encabezado «las fechas de compras y pagos son del oráculo». **Verificador:** extender S1-03 a fechas (ISO y dd-mmm-aaaa) fuera del oráculo en frases sin rótulo y comprobar que la versión actual da FALLA |
| R2-02 | baja | S1-03 | reglas L175 (RN-11); s1_verifica.py, bloque S1-03 | 1,814.50 no está en el oráculo. Pasa como «cuenta citada» (`financiamiento.dos_quincenas.total / 2`), una excepción que agregó el verificador y que S1-03 no menciona. La cuenta es exacta y da lo mismo que 1,115.85 + 698.65 | Citar las dos cifras del oráculo (`data.json: financiamiento.dos_quincenas.axel.por_quincena = 1,115.85` y `.karen.por_quincena = 698.65`), o pedir que criterios 1.2 admita por escrito la «cuenta citada» |
| R2-03 | baja | S1-16 | `docs/estado/p1/S1-esquema-1.2.md` L35-36 | Supuestos trae dos cambios de esquema: K12 aplicado también a `credit.annual_rate_decimals` y `promotion.rate_decimals`, y columnas nuevas que admiten NULL. S1-16 pide que ningún supuesto cambie el esquema. En el fondo están respaldados: K12 es general, K18 hace que Axel capture la tasa anual y los NULL los impone ALTER TABLE ADD COLUMN | Pasarlos a «Hecho» como aplicación de K12 y nota técnica. Si el integrador no los acepta como K12, abrir bloqueo |
| R2-04 | baja | Evidencia | S1-esquema-1.2.md L17, L26, L30-31 | El reporte S1 quedó con cifras anteriores a la corrección: 69 disparadores (hoy 71), diccionario de 173,798 bytes (hoy 174,605) y la línea S1-03 sin el conteo aparte de ilustrativos. El dato vigente está en `correccion-S1-r1.md` | El integrador actualiza esas líneas o enlaza `correccion-S1-r1.md` (es la solicitud 1 de ese reporte) |
| R2-05 | baja | K18 (H-04 de r1, abierto) | `backend/db/migrations/00002_puerta1.sql` L112-115 | Sin cambio: la base no exige que los datos de K18 vayan todos juntos, y `vat_on_interest_bp` NULL quiere decir «ya trae IVA». Lo cubre el invariante 16 en la API | La de r1: disparador de todo o nada, o columna explícita, con su prueba |
| R2-06 | baja | K22 (H-05 de r1, abierto) | reglas L181 (RN-12) | RN-12 sigue recomendando la lectura (b) aunque K22 dice «entre 30 y 90 días decae en línea recta de 1 a 0». El valor provisional en uso es el (a), así que ninguna cifra está mal | Quitar la recomendación del contrato y dejar la discusión en el bloqueo abierto |
| R2-07 | baja | S1-17 (H-06 de r1, abierto) | 00002 L112 | `credit.rate_base` (`daily`/`monthly`) convive con `credit_purchase.rate_basis`, de otro significado | Renombrar antes del commit, en Up, Down, reglas y glosario, y regenerar el diccionario |
| R2-08 | baja | S1-06 (H-07 de r1, abierto) | 00002 L3 | La cabecera escribe «11-oct-2026», que es formato de pantalla; el grep literal de S1-06 da 1 coincidencia | Escribir 2026-10-11 |
| R2-09 | baja | S1-03 (H-08 de r1, abierto en parte) | reglas: RN-13, RN-14, RN-16 a RN-19, RN-22 a RN-24 y RN-26; frases de RN-04, RN-07, RN-10, RN-11, RN-15 y RN-28 | Ya está rotulada la frase de RN-07 «la serie real sería…». Siguen 10 ejemplos y 6 frases sin la leyenda literal «(sin dato en el oráculo)» | Usar la leyenda literal y poner el motivo después |

## Estado de los hallazgos de la ronda 1
- **S1-03 (alta, H-01 de r1): reabierto en parte, ahora R2-01.** RN-12 está cerrado: leyenda literal; pesos 1 y 0.5 con la lectura (a) y 60/61 ≈ 0.98 y 30/61 ≈ 0.49 con la (b), comprobados; fechas 8-nov, 9-nov, 9-dic-2026 y 8-ene-2027 correctas. RN-27 está cerrado: leyenda literal; 1,947.50 / 31 = 62.82; 439.76 y 439.74 correctos. RN-11: el 1,814.50 ya trae su cuenta (queda R2-02, baja), pero el encabezado sigue diciendo que todas las fechas son del oráculo y trae 16-dic-2026. Contar aparte los ilustrativos: hecho en `correccion-S1-r1.md`; `S1-esquema-1.2.md` sigue sin actualizar (R2-04).
- **S1-09 / K8 / h (media, H-02 de r1): cerrado.** 00001 L1111 dice «Lunes que abre la semana del plan (K8…)», el diccionario L1774 lo replica y la 00002 L69-88 lo exige con dos disparadores. Mis pruebas: martes, domingo y cambio a miércoles rechazados; lunes aceptado. No queda ninguna cita de `U1-periodos-semana-quincena` en el esquema, el diccionario ni las reglas. El bloqueo abierto `U1-lista-vigente-domingo` solo toca la semana que se muestra por defecto, no el esquema.
- **S1-04 (media, H-03 de r1): cerrado.** 00002 L32 y diccionario L1107: «U1 P-05 (N-071, N-075 y N-076)», como U1 L28 (K3 → P-05, N-071, N-075, N-076, F-8). `telegram_chat_candidate`: «U1 P-26.09: N-299 … N-316», como U1 L29. No queda ninguna cita de N-087 para el ticket.
- **H-04 a H-08 (bajas):** H-04, H-05, H-06 y H-07 siguen abiertos (R2-05 a R2-08). H-08 está abierto en parte (R2-09).

## Trazabilidad K (verificación g)
| K | Esquema (00001 = M1, 00002 = M2) | Diccionario | Reglas |
|---|---|---|---|
| K1 | M1 L1106-1107, L1122-1123 | L82-83 | L204 (RN-15) |
| K2 | M1 L299 | L530 | L246 (RN-22) |
| K3 | M2 L29-34 | L1107, L1125-1126 | L251 (RN-23) |
| K6 | M1 L161, L176; M2 L145 | L337, L344, L363 | L62 (C-15) |
| K8 | M1 L186, L194, L200, L1111; M2 L70-88 | L374, L384, L1774, L1795-1796 | L118-119, L122 (RN-06) |
| K9 | M2 L91-105 | L999-1000, L1901 | L256 (RN-24) |
| K11 | M2 L144 | L2251, L2270 | (sin regla; S1-14 lo pide en el diccionario) |
| K12 | M2 L115, L119, L123, L127 | L770, L814, L855 | L51 (C-04, «Z» solo en entradas), L261 (RN-25) |
| K13 | M1 L39 | L143 | L105 (RN-04) |
| K14 | — (no cambia el esquema) | — | L74 (RN-02), L132 (RN-07, igual a la opción A del bloqueo) |
| K15 | — | — | L104 (RN-04) |
| K16 | M2 L10, L138-140 (sin semilla) | L2258-2260 | L165, L175 (RN-11) |
| K18 | M2 L111-114 | L740, L767-768 | L266 (RN-26) |
| K19 | — | — | L241-242 (RN-21; 6,293.44 = 1,258.70 + 5,034.74, citas exactas) |
| K20 | — | — | L271-272 (RN-27) |
| K21 | M2 L143 | L2251, L2271 | L276 (RN-28) |
| K22 | — (`price.stale_after_days` = 30 desde 1.1) | L2263 | L180-182 (RN-12) |
| K23 | M1 L1169; M2 L134 | L1863, L1869 | L217 (RN-17) |
