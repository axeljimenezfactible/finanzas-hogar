# Revisión S1 · esquema 1.2 · ronda 1 de 3 (flujo P1)

Revisor independiente: `revisor-S1-r1` · 11-oct-2026 · Evalúa contra `docs/criterios/S1.md` 1.1 y `docs/plan/decisiones.md` §K (K1–K3, K6, K8, K9, K11–K16, K18–K23). Revisé la migración `backend/db/migrations/00002_puerta1.sql` completa, los comentarios cambiados de `00001_esquema.sql` (`git diff HEAD`), `backend/db/tools/`, el diccionario 1.2 (por grep y sed), las reglas 1.2 (todo lo agregado en el diff) y el reporte `docs/estado/p1/S1-esquema-1.2.md`. No me basé en el reporte: reejecuté todo con mis propios scripts en `datos/tmp/rev-s1r1/`, que ya borré.

## Veredicto
**NO APROBADO.** Hay 1 hallazgo alto (H-01) y 2 medios (H-02 y H-03); los tres se corrigen en texto y comentarios, sin cambiar cifras del oráculo. La migración es sólida: aplica limpia, el Down deja el esquema 1.1 byte por byte, goose sube, baja y vuelve a subir, el DDL de 00001 no cambió, los scripts reproducen 389500 y 239528, el diccionario se regenera idéntico y las semillas nuevas solo llevan valores de §K.

## Verificación propia (a)–(j)
| # | Qué | Resultado |
|---|---|---|
| a | 00001+00002 con sqlite3 3.53.4 en `datos/tmp/rev-s1r1/rev.db` | `integrity_check` = ok · `foreign_key_check` = 0 filas · 55 tablas, 55 STRICT (`pragma_table_list.strict`), 0 sin STRICT · 114 índices con SQL (152 con los automáticos) · 69 disparadores |
| b | Up+Down de 00002 contra 00001 sola | `sqlite_master` (type, name, tbl_name, sql): idéntico, 1263 líneas, 111,282 bytes. `.dump` completo con datos: idéntico (18 INSERT) al normalizar las marcas de tiempo `strftime('now')` |
| c | goose v3.28.0: up / version / down-to 0 / version / up / version | rc 0 en los 6 pasos; versiones 2 → 0 → 2; 0 objetos fuera de goose tras down-to 0; esquema y datos tras el segundo up iguales a los de sqlite3 |
| d | DDL de 00001 contra `git show HEAD` | `sed 's/--.*//'`: diff vacío (1302 líneas). Mi tokenizador, que respeta las cadenas: 1302 = 1302, IGUALES. Datos sembrados de HEAD y de hoy: idénticos |
| e | `GOOSE=… python3 -I backend/db/tools/s1_verifica.py` | rc 0, «RESUMEN TODO OK»: KUESKI 389500 y 239528 sobre 1.2; NEG12: 24 negativas rechazadas y 10 positivas aceptadas; NEG 46/46; DOWN y GOOSE OK; no deja temporales. Por mi cuenta: Σ `despensas[0].pagos[*]` = 389500 (Axel 239528, Karen 149972), y 33 pruebas mías de las columnas nuevas (CHECK de K18 solo en tarjeta y crédito de tienda, `rate_decimals` ≥ 0, candidato de Telegram confirmado ⇔ con persona, `chat_id` único y numérico, FK, «Z» obligatoria, clave de 16 a 128 única con varios NULL, `goal` enumerado): 32 de 32 como se espera; una prueba extra de `credit_purchase` falló por otro CHECK y no la cuento |
| f | `python3 -I backend/db/tools/gen_dict.py` + `cmp` | idéntico a `docs/contratos/diccionario-datos.md` (173,798 bytes) |
| g | Cada K aplicada, con línea | Tabla «Aplicación de §K» abajo: 18 de 18 presentes; K8 incompleta en `shopping_list` (H-02) |
| h | PROVISIONAL que cite un bloqueo resuelto | 0. Las 11 marcas de las reglas citan 4 bloqueos sin «Respuesta de Axel»: `criterios-frescura-extremos`, `S1-fuentes-precio`, `S1-ajuste-objetivo-kcal`, `criterios-promedios-redondeo`; la migración tiene 0 y el diccionario solo la definición. Sin la palabra PROVISIONAL, un comentario sigue tratando como pendiente un bloqueo resuelto (H-02) |
| i | Semillas nuevas | La 00002 agrega exactamente 3 filas a `app_setting` (7 → 10) y ninguna otra tabla cambia: `projection.repeat_cadence_enabled` = false (K21), `session.max_age_days` = 90 (K11), `idempotency.retention_hours` = 24 (K6). `optimizer.*`: solo `optimizer.weights`, que viene de 1.1; las tres restricciones de K16 no tienen valor. Los factores 1.375, 1.55 y 1.725 salen de la opción B del bloqueo `S1-meta-kcal`, que la D elegida incluye («B más un ajuste…»). No hay valores inventados |
| j | Ejemplos de las reglas 1.2 | 28 de 28 reglas con ejemplo. Resolví las 62 citas con ruta contra `model.json` y `data.json`: todas existen y las 17 que llevan «= valor» coinciden. Las cuentas son exactas: 6,293.44 = 1,258.70 + 5,034.74 y 14,362.48 − 13,103.78 = 1,258.70; 17.27 % y 25.8685 % con HALF_UP; 1,034.74; 2,556 kcal; 88.4 %; 5,226.39 y 5,000.00 en `credito.eventos`; 13:20 de la Ciudad de México = 19:20Z. **Falla en dos reglas:** los ejemplos de RN-12 y RN-27 dan resultados que no están en el oráculo y no llevan el rótulo de ilustrativo (H-01) |

## Criterios de S1.md 1.1
| Criterio | Sev. | Estado | Evidencia |
|---|---|---|---|
| S1-01 Entidades | alta | Cumple | 31 de 31 (s1_verifica); el diccionario regenerado tiene las 55 tablas |
| S1-02 Migraciones | alta | Cumple | (a), (c) |
| S1-03 Ejemplo por regla | alta | **No cumple** | H-01: RN-12 y RN-27 sin rótulo, con cifras fuera del oráculo |
| S1-04 Sin huérfanas | alta | Cumple con observación | 55 de 55 con «Usado por». Revisé 5 contra U1: `telegram_chat_candidate` (N-299, N-316), `shopping_list` (N-150), `category` (N-282), `budget` (N-087) y `attachment` (P-05, N-075). La 00002 le cita a `attachment` una necesidad equivocada (H-03) |
| S1-05 Dinero exacto | alta | Cumple | 0 REAL; `min_payment_fixed_cents` INTEGER; 37 columnas `_cents`, todas INTEGER |
| S1-06 Fechas y zona | media | Cumple con observación | Los instantes nuevos exigen «Z» (GLOB + `julianday`). El grep literal da 1 coincidencia, en un comentario de la 00002 (H-07) |
| S1-07 Personas y acceso | alta | Cumple | s1_verifica S1-07; la 00002 no agrega columnas de permisos |
| S1-08 Reparto | alta | Cumple | Sin cambios en 1.2; s1_verifica S1-08 |
| S1-09 Periodos | alta | Cumple con observación | 18 filas, límites de K8 y comentarios FIRMES (K8) en `period`; falta K8 en `shopping_list.start_date` (H-02) |
| S1-10 Créditos | alta | Cumple | Cargué 2q 53450/309450 `rate_decimals` 2 y 4q 80050/309450 `rate_decimals` 4; leídas de vuelta y mostradas con HALF_UP: 17.27 % y 25.8685 % (K12) |
| S1-11 Comida y salud | media | Cumple | 12 de 12; `health_profile.goal` (K23) |
| S1-12 Ida y vuelta | alta | Cumple | 164 valores, 0 diferencias (s1_verifica) y mis sumas de (e) |
| S1-13 Ajustes | media | Cumple | 14 de 14; K4 en la fila 9, K21 en la 2, K23 en la 13, K13 en la 12, K18 en la 5 |
| S1-14 Sesiones | alta | Cumple | Fila 1 del mapa de Ajustes: `login_attempt` + `auth.max_failed_attempts` (5), `auth.lockout_window_minutes` (15, ventana deslizante por usuario e IP, K11) y `session.max_age_days` (90, K11). Sin secretos |
| S1-15 Integridad | media | Cumple | FK `telegram_chat_candidate.person_id` con índice; enumerados nuevos con CHECK; mi prueba de FK falla como debe |
| S1-16 Sin inventados | alta | Cumple | (h), (i); K6 y K7 sin marca |
| S1-17 Idioma y nombres | media | Cumple con observación | 656 identificadores en `[a-z0-9_]`; `rate_base` junto a `rate_basis` (H-06) |

## Aplicación de §K (grep con línea)
Archivos: 00001 = `backend/db/migrations/00001_esquema.sql`, 00002 = `…/00002_puerta1.sql`, D = `docs/contratos/diccionario-datos.md`, R = `docs/contratos/reglas-negocio.md`.

| K | Dónde quedó | Estado |
|---|---|---|
| K1 | 00001 L1106-1107 y L1122-1123; D L81-82 y L1766; R L202 (RN-15) | Aplicada |
| K2 | 00001 L299 (`budget.valid_to_date`); D L529; R L244 (RN-22, «de este mes en adelante» preseleccionado) | Aplicada |
| K3 | 00002 L32-33; D L1124-1125; R L249-251 (RN-23) | Aplicada (cita de U1 equivocada, H-03) |
| K6 | 00001 L161 y L176; 00002 L123; D L336, L362 y L2248; R L60 (C-15) | Aplicada |
| K8 | 00001 L186, L194 y L200; D L373, L383, L391 y L2305; R L116-117 (RN-06) | Incompleta: 00001 L1111 y D L1773 (H-02) |
| K9 | 00002 L73-83; D L998-999; R L254-256 (RN-24) | Aplicada |
| K11 | 00002 L122; D L2248, L2260 y L2267 | Aplicada |
| K12 | 00002 L93, L97, L101 y L105; D L769, L813 y L854; R L49 (C-04) y L259-261 (RN-25) | Aplicada |
| K13 | 00001 L39; D L142, L2271 y L2278; R L103-106 (RN-04) | Aplicada |
| K14 | R L72 (RN-02) y L130 (RN-07) | Aplicada |
| K15 | R L102 (RN-04) | Aplicada |
| K16 | 00002 L9 y L116-118 (sin semilla); D L2255-2257 y L2308; R L163-168 (RN-11) | Aplicada |
| K18 | 00002 L90-92; D L766-768 y L2298; R L264-266 (RN-26) | Aplicada (observación H-04) |
| K19 | R L239-240 (RN-21) | Aplicada |
| K20 | R L269-270 (RN-27) | Aplicada (rótulo, H-01) |
| K21 | 00002 L121; D L2248 y L2268; R L274-276 (RN-28) | Aplicada |
| K22 | D L2260; R L178-183 (RN-12) | Aplicada (rótulo, H-01; observación H-05) |
| K23 | 00001 L1169; 00002 L112; D L1860, L1866 y L2279; R L215-224 (RN-17) | Aplicada |

## Hallazgos
| ID | Sev. | Criterio | Ruta | Hallazgo | Corrección sugerida |
|---|---|---|---|---|---|
| H-01 | alta | S1-03; (j) | `docs/contratos/reglas-negocio.md` L183 (RN-12), L270 (RN-27), L173 (RN-11); `backend/db/tools/s1_verifica.py` (bloque S1-03); `docs/estado/p1/S1-esquema-1.2.md` | El ejemplo de RN-12 da pesos de 1, 0.5, 60/61 ≈ 0.98 y 30/61 ≈ 0.49, y el de RN-27 da 62.82, 439.76 y 439.74. Ninguna de esas cifras está en el oráculo, y los dos ejemplos empiezan con «**Ejemplo.**», sin rótulo de ilustrativo. S1-03 nombra la frescura como regla sin dato en el oráculo, con ejemplo ilustrativo, y espera «0 ejemplos con cifra que no exista en el oráculo; 0 ejemplos ilustrativos sin el rótulo». En RN-11, que es nuevo en 1.2, el ejemplo dice que «las fechas y cifras son del oráculo», pero trae 1,814.50, que se deriva de 3,629.00 / 2 y no está en el oráculo. El S1-03 de `s1_verifica.py` solo revisa las rutas citadas, por eso lo da OK, y el reporte no cuenta aparte los ejemplos ilustrativos | Rotular los ejemplos de RN-12 y RN-27 con «Ejemplo ilustrativo (sin dato en el oráculo)» y dejar citadas sus entradas del oráculo (fecha 2026-10-10 y 1,947.50). En RN-11, citar la derivación (`data.json: financiamiento.dos_quincenas.total` 3,629.00 / 2) o marcar esa frase como ilustrativa. Que el S1-03 de `s1_verifica.py` marque las cifras de dos decimales que no estén en `model.json` ni en `data.json` dentro de los ejemplos sin rótulo, y que cuente los ilustrativos; poner ese conteo en el reporte |
| H-02 | media | S1-09; (g) K8; (h) | `backend/db/migrations/00001_esquema.sql` L1111; `docs/contratos/diccionario-datos.md` L1773 | `shopping_list.start_date` sigue diciendo «Qué día empieza la semana lo fija U1 (bloqueo U1-periodos-semana-quincena)». Ese bloqueo está RESUELTO: K8 fija la semana de lunes a domingo, y RN-06 y la nota U1/S2 (D L2305) ya lo aplican. Los comentarios de `period` se actualizaron y este no; además nada pide que la lista empiece en lunes | Cambiar solo el comentario de 00001 L1111 a «Lunes que abre la semana del plan (K8: semana de lunes a domingo)» y regenerar el diccionario con `gen_dict.py`. Si se quiere hacer cumplir: un disparador en la 00002 (`strftime('%w', NEW.start_date) <> '1'` → ABORT, con su DROP en el Down) o el invariante 17 del diccionario (lo valida la API) |
| H-03 | media | S1-04 (trazabilidad) | `backend/db/migrations/00002_puerta1.sql` L31; `docs/contratos/diccionario-datos.md` L1106 | `@used_new` de las columnas del borrador cita «U1 N-087 (confirmación del ticket)». En U1, tanto en 1.0 como en 1.1, N-087 es «Ajustar presupuesto», y el propio diccionario lo usa así en L521. El borrador de K3 es N-075 «Guardar para después», en P-05 «Confirmar ticket»; U1 L28 asigna K3 a N-071, N-075, N-076 y F-8 | Cambiar a «U1 P-05 (N-071, N-075 y N-076)» en 00002 L31 y regenerar el diccionario. Conviene sumar N-316 «Confirmar chat detectado» al `@used` de `telegram_chat_candidate` (L43) |
| H-04 | baja | K18 | `00002_puerta1.sql` L90-93; D L2298 (invariante 16); R L264 (RN-26) | «Si incluye IVA» va en `vat_on_interest_bp`, donde NULL = «la tasa ya trae IVA», así que esa columna no tiene estado de «falta capturar». Lo cubre el invariante 16 (la API guarda juntos los datos de K18), pero la base no lo exige: cada columna admite NULL por separado. Si se colara una captura parcial, S3 supondría que el IVA ya viene incluido | Hacer que la base exija todo o nada en `bank_card` y `store_credit` (disparador en la 00002 sobre `rate_base`, `min_payment_pct_bp`, `min_payment_fixed_cents` y `annual_rate_num`), o una columna booleana explícita con NULL = falta capturar; agregar una prueba a NEG12 |
| H-05 | baja | K22 | R L179 (RN-12) | RN-12 recomienda la lectura (b), (90 − edad)/61, y dice que es «la recta de la opción C». La respuesta de Axel en `S1-frescura-peso` y K22 dicen «de 30 a 90 días decae en línea recta de 1 a 0», es decir, que el día 30 pesa 1, que es la lectura (a). RN-12 usa (a) como valor provisional, así que ninguna cifra está mal | Quitar la recomendación del contrato y dejar la discusión en el bloqueo abierto `criterios-frescura-extremos`, o citar K22 textual |
| H-06 | baja | S1-17 | 00002 L90 (`credit.rate_base`); RN-25 y RN-26; invariante 16 | La columna nueva `rate_base` (`daily`/`monthly`) convive con `credit_purchase.rate_basis` (`term_total`/`annual`, de 1.1): son conceptos distintos con nombres casi iguales, y «base» no es el término en inglés | Antes del commit, renombrar (por ejemplo `rate_accrual_basis`) en Up, Down, reglas y glosario, y regenerar el diccionario |
| H-07 | baja | S1-06 | 00002 L3 | La cabecera escribe «11-oct-2026», formato de pantalla. El grep literal de S1-06 da 1 coincidencia; es un comentario y el criterio solo cuenta datos | Escribir 2026-10-11 |
| H-08 | baja | S1-03 | R (10 ejemplos «Ejemplo ilustrativo …»); R L132 (RN-07) | Ningún ejemplo usa la leyenda literal del criterio, «Ejemplo ilustrativo (sin dato en el oráculo)»: hay 0 apariciones contra 10 con motivos distintos. En RN-07 sigue sin rótulo «la serie real sería 1,573.35 ×3 + 1,573.37», que viene de 1.0 y no está en el oráculo | Usar la leyenda literal, conservando el motivo después, para que el grep del criterio encuentre los ejemplos; rotular la frase de RN-07 |

## Estado de los hallazgos anteriores (WF-0 r3, heredados que tocan S1)
| ID | Estado | Evidencia |
|---|---|---|
| R2-03 (baja) tope de Karen | Cerrado | 00001 L39 y D: «No hay un tope fijo… (K13)»; RN-04: 374.93 «no un tope» |
| R2-04 (baja) versión de reglas | Cerrado | `reglas-negocio.md` 1.2 con «Cambios en 1.2» |
| R2-06 (baja) «Z» y decimales de tasa | Cerrado en S1 (C-04 y `rate_decimals`, K12); el lado de C1 le toca a C1 | R L49; 00002 L93-105 |
| R2-07 (baja) ilustrativos en criterios | Cerrado en S1.md 1.1; las reglas todavía no cumplen el rótulo | H-01, H-08 |
| R2-08 (baja) scripts en anexos | Cerrado en S1: `backend/db/tools/gen_dict.py` y `s1_verifica.py`, los dos corren (rc 0, cmp idéntico). Que el reporte de WF-0 enlace a esos archivos es una solicitud del agente al integrador | (e), (f) |
| R3-02 (baja) límites técnicos sin marca | Cerrado por K6 | 0 marcas en las líneas de K6 |
