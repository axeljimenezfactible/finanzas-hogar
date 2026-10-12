# P2c · S1-puerta2 · Esquema 1.4 (migración 00004)

Agente S1-puerta2 (Opus), flujo P2c «Aplicación de la puerta 2», 11-oct-2026, noche. Reanudación tras el corte de las 22:35 (run wf_cbf48e23-39e): empecé de cero; no había archivos parciales míos.
Fuentes: `docs/estado/puerta2/RESPUESTAS.md` (Q2-11 B explícita; modo «lo recomendado + excepciones»), `docs/estado/puerta2/PUERTA-2.md` §3, §4 y §5, `docs/plan/decisiones.md` K46, K47 y K50, los bloqueos `S1-parametros-avisos`, `U4-iva-credito`, `U4-reparto-retroactivo`, `S3-serie-sin-intereses`, `U4-dia-mes-reporte` y `S2-ventanas`, y `backend/db/migrations/00003_puerta1b.sql`.

## Hecho

### Claves nuevas de `app_setting` (lista exacta para C1 y U4)
Todas del hogar (`user_account_id` NULL), sembradas por `backend/db/migrations/00004_puerta2.sql`. Las horas son texto JSON «HH:MM» de 24 h en America/Mexico_City; lo demás, enteros JSON.

| clave | valor inicial (value_json) | rango en el esquema | decisión | sección de Ajustes |
|---|---|---|---|---|
| `reports.daily_send_time` | `"21:30"` | hora 00:00 a 23:59 | Q2-72 (K50), lugar K47 | Reportes y frecuencias |
| `reports.weekly_weekday` | `7` | 1 a 7 (ISO 8601, 7 = domingo) | decisión F, lugar K47 | Reportes y frecuencias |
| `reports.weekly_send_time` | `"19:00"` | hora 00:00 a 23:59 | Q2-72 (K50), lugar K47 | Reportes y frecuencias |
| `reports.monthly_day_of_month` | `1` | 1 a 31 (mes corto: último día, Q2-66 B) | decisión F, Q2-66 B, lugar K47 | Reportes y frecuencias |
| `reports.monthly_send_time` | `"08:00"` | hora 00:00 a 23:59 | Q2-72 (K50), lugar K47 | Reportes y frecuencias |
| `alerts.send_time` | `"09:00"` | hora 00:00 a 23:59 | Q2-73 (K50), K43, lugar K47 | Notificaciones |
| `alerts.stale_prices_advance_days` | `1` | 0 o más | Q2-73 (K50), K43, lugar K47 | Notificaciones |
| `alerts.meal_plan_missing_weekday` | `7` | 1 a 7 | decisión F, lugar K47 | Notificaciones |
| `optimizer.date_window_days` | `7` | 0 o más | Q2-75 (K50) | Créditos y condiciones |
| `ocr.timeout_seconds` | `30` | 0 o más | Q2-76 (K50) | Datos y respaldos (PROPUESTA S1) |

Rangos nuevos sobre claves que ya existían (anotación `@param_rango`, sin repetir su `@param`): `comparator.timeout_seconds` de 1 a 30 (Q2-74; el valor sigue en 5), `price.stale_after_days` de 1 o más (Q2-77). La regla entre dos claves de Q2-77 (`price.zero_weight_after_days` ≥ `price.stale_after_days`) no va en el esquema: la valida el motor de precios. HECHO verificado: `backend/internal/precios/types.go` l. 107 a 112, `Settings.Validate` devuelve `ErrInvalidArgument` si `StaleAfterDays < 1` o `ZeroWeightAfterDays < StaleAfterDays`.

### Migración 00004 (esquema 1.3 → 1.4; 00001 a 00003 sin tocar)
- **(a) Q2-11 B (K47).** Las 8 claves `reports.*` y `alerts.*` de arriba. `weekday` y `day_of_month` salen de `notification.params_json` (`weekly_summary`, `meal_plan_missing` y `monthly_report` quedan con `{}`) conservando `updated_at`. Cada clave toma el valor que su aviso tenga al migrar, si es válido (en una base nueva, 7, 7 y 1); si no, el de la decisión F. `send_time` y `advance_days` nunca estuvieron en `params_json`: nacen como claves y el Down solo las borra. El Down repone `weekday` y `day_of_month` al principio de `params_json` con el valor vigente (con los iniciales, el texto de 1.3 byte por byte). Anotaciones `@param` (10), `@param_sin_valor_cerrado` (5: `optimizer.date_window_days`, `send_time`, `advance_days`, `weekday` y `day_of_month`) y `@resuelve S1-parametros-avisos`. RN-29 y el diccionario ya no dicen PROVISIONAL.
- **(b) Q2-74 a Q2-77.** `optimizer.date_window_days` = 7 y `ocr.timeout_seconds` = 30. Disparadores `trg_app_setting_param_range_ins` y `_upd`, con los rangos de una sola clave de la tabla de arriba, y `trg_app_setting_param_time_ins` y `_upd`, que piden la hora «HH:MM» con dos dígitos, de 00:00 a 23:59. Q2-77 queda documentada en RN-12, RN-29, en la invariante 17 del diccionario y en `@param_rango`.
- **(c) Q2-67 opción 1.** Columna `credit.vat_included INTEGER CHECK (vat_included IN (0, 1) AND (vat_included = 0 OR vat_on_interest_bp IS NULL))`. 1 = la tasa ya incluye IVA, 0 = se suma el de `vat_on_interest_bp` y NULL = falta capturar. `vat_on_interest_bp` conserva su significado y solo lleva valor con 0 (`@col_note`). Datos que ya existan: con porcentaje pasan a 0 y sin porcentaje a 1, el significado que les daba 1.3; `updated_at` se conserva (ver Supuestos). Para C1: `missing_terms` pasa a 5 datos (`annual_rate`, `rate_base`, `vat`, `min_payment_pct_bp` y `min_payment_fixed_cents`); está en RN-26, en la invariante 16 y en `@used_new`.
- **(d) Q2-66 B.** `reports.monthly_day_of_month` acepta de 1 a 31 en el esquema. La regla del mes corto (sale el último día del mes) está en RN-29 y en su `@param`, y la calcula S7.
- **(e) Q2-68 C.** Sin columna ni tabla. HECHO: `expense_share` ya congela las partes de cada gasto («Se congela al guardar para que cambiar una regla no altere gastos pasados», `00001_esquema.sql` l. 656). Con «solo los nuevos» no se toca nada guardado; con «también los ya capturados», el motor de reparto recalcula. La opción C lo asigna a C1 (parámetro de alcance) y al motor de reparto. Está en RN-03 y en la invariante 18 del diccionario.
- **(f) Q2-38 (a).** El `CHECK (interest_cents >= 0)` de `scheduled_payment` queda igual: con capital = monto, los intereses son 0 en cada pago. RN-07 gana la excepción «serie sin intereses», con un ejemplo ilustrativo: 3,094.50 a 4 quincenas paga 773.62 ×3 + 773.64 con intereses 0.

### Contratos regenerados o editados
- `docs/contratos/diccionario-datos.md` 1.4: lo regenera `python3 -I backend/db/tools/gen_dict.py` (el generador está en `backend/db/tools/`, no en `docs/contratos/tools/`, que no existe). Cambios:
  - generados: encabezado, «Cambios en 1.4», columna `credit.vat_included`, notas de `notification.params_json` y `credit.vat_on_interest_bp`, inventario K44 con 25 claves (3 sin valor) y semilla y avisos de la 00004;
  - texto fijo: filas 8 y 9 de «Ajustes: 14 secciones», sin `optimizer.date_window_days` en «sin valor», invariante 16 (IVA) y las nuevas 17 (Q2-77) y 18 (Q2-68).
- `docs/contratos/reglas-negocio.md` 1.4: «Cambios en 1.4», RN-03 (alcance Q2-68), RN-07 (serie sin intereses), RN-11 (ventana de 7 días), RN-12 (rango Q2-77), RN-26 (`vat_included` y 5 datos) y RN-29 (FIRME, 10 filas nuevas, rangos, mes corto y horas; sin PROVISIONAL).
- Herramientas de S1 (ver Supuestos 1):
  - `backend/db/tools/gen_dict.py`: anotaciones `@param_rango`, `@param_sin_valor_cerrado` y `@col_note` (la última manda), y textos de 1.4;
  - `backend/db/tools/s1_verifica.py`: checks nuevos PROG14 (sustituye a PROG13), K47, AVISOS14 y DOWN14; S1-16 y K44 tratan como historia las marcas de la 00003 que cita un bloqueo con `@resuelve`; NEG12 solo cuenta las 3 claves de K16; DDL compara también la 00003 contra HEAD.

### (g) Lo que NO toca el esquema y por qué (PUERTA-2 §3 y §5)
- Q2-12 (sin respuesta), Q2-13 a Q2-20 y Q2-25 a Q2-27: reglas del motor S2 y de C1, sin columnas.
- Q2-21 a Q2-24 (C): el valor inicial A se siembra como partidas planificadas y se cambia con B. La fila 2 de Ajustes ya trae `planned_item.interval_count`, `.credit_id`, `.term_count`, `.accrual_months` y el reparto con `cost_share_rule_part.fixed_cents`. Quién siembra, S2 o S8, no lo decido.
- Q2-28 a Q2-36: fórmulas del motor S3; Q2-33 (b) pide `store_id` en `Scenario` de C1, no en el esquema.
- Q2-37 (a): `quoted`, sin cambio de esquema según el archivo.
- Q2-39 a Q2-42 (A): solo C1. Q2-40 A y Q2-41 A dicen «exige cambiar C1» y `shopping_list_item.store_id` ya existe.
- Q2-43 (A): sin esquema ni C1, según el archivo.
- Q2-44 a Q2-46 (A): la opción C de Q2-44 habría tocado el esquema, pero no fue la elegida.
- Q2-47, Q2-49, Q2-53 a Q2-55, Q2-59 a Q2-61 y Q2-63 (A): las opciones con tabla o columna nueva (Q2-53 B, Q2-54 B, Q2-60 B y Q2-61 B) no se eligieron. PROFECO sigue fuera de la base (Q2-53 A).
- Q2-62 (B): se calcula con el historial y no lleva clave nueva (la clave era la A).
- Q2-48, Q2-50 a Q2-52 y Q2-56 a Q2-58: sin respuesta. Siguen como están y no las decido.
- Q2-64: nada en esquema (el archivo lo dice).
- Q2-65 (A): campos de C1 y S2.
- Q2-69 y Q2-70 (a): las comisiones siguen como dato y las fórmulas son de motor.
- Q2-71 (1): cero cambios.
- §5: Q2-78 a Q2-84 son de diseño (Q2-82 se guarda en el dispositivo y Q2-84 no es parámetro). Q2-86 (A) es solo C1. Q2-92 usa `store.municipality`, que ya existe. Q2-93 carga datos del catálogo, sin DDL. Q2-85, Q2-87 a Q2-91 y Q2-94 no tienen respuesta. Las opciones de Q2-87 (b) y Q2-88 (c), que tocarían el esquema, no se aplican.

## Verificación
- Antes de cambiar nada, `make verify-s1 GOOSE=…/goose` daba FALLAS: [S1-16, K44]. La causa: `S1-parametros-avisos` ya tenía «Respuesta de Axel» en el árbol y la 00003, que no se edita, sigue citándolo como PROVISIONAL.
- `make verify-sql GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`:
  - Up de las 4 migraciones con sqlite3 3.53.4: tablas=55 (STRICT=55), índices=152, disparadores=79, foreign_key_check=0 filas, integrity_check=ok;
  - Down en orden inverso (00004,00003,00002,00001) deja 0 objetos; Up otra vez: 55 tablas;
  - goose v3.28.0: up, down-to 0 y up otra vez sin errores.
- `make verify-s1 GOOSE=…/goose` da RESUMEN TODO OK y «s1_verifica.py aplicó las 4 migraciones del comodín». Entre ellos:
  - PROG14 OK: params_json `{}` en los 3 avisos, 10 claves con su valor, 5 `@param_sin_valor` cerrados, diccionario y RN-29 sin PROVISIONAL;
  - K47 OK: 17 negativas rechazadas y 11 positivas aceptadas. Las negativas incluyen semana 8, día 32, «24:00», «9:00», comparador 31 y 0, precio viejo 0 y `vat_included` = 1 con porcentaje;
  - AVISOS14 OK: semana 5 y día 31 pasan a las claves; un 9 inválido da 7; `vat_included` de créditos de 1.3 queda [0, 1]; el Down repone semana 3; `updated_at` conservado;
  - DOWN14 OK: el esquema 1.3 vuelve idéntico (244 objetos) con sus 20 claves, avisos byte por byte y créditos;
  - DOWN13, DOWN y GOOSE OK (versión final 4);
  - DDL OK: 00001 a 00003 iguales a git HEAD;
  - S1-03 OK: ejemplo ilustrativo de RN-07; S1-16 OK: 0 líneas sin bloqueo.
- `python3 -I backend/db/tools/gen_dict.py` da versión 1.4, 25 parámetros, 3 sin valor y 0 valores que no coinciden. Hay 55 tablas, 152 índices, 79 disparadores y 0 columnas sin descripción. Una segunda corrida a otro archivo da `cmp` idéntico: es idempotente.
- `make verify-u1` (RESULTADO: TODO OK) y `make verify-criterios`: salida 0.
- `make -k verify-all` (antes del último ajuste de anotaciones): WF-0, S1, U1, Criterios y WF-1 OK; verify-c1 FALLA.
- Hoy `make verify-c1` falla con «gen2: tipo de parámetro desconocido: hora HH:MM». Es esperado: el generador de C1 solo conoce «entero» y «lista de enteros», y C1 todavía no expone las claves nuevas (ver solicitudes). Usé «entero» en los demás parámetros para que el único tipo nuevo sea la hora.

## Supuestos
1. Edité las herramientas de S1 (`gen_dict.py` y `s1_verifica.py`), que no están en la lista literal de ALCANCE. Sin eso, `make verify-s1` no podía pasar:
   - PROG13 exigía `weekday` en `params_json`;
   - S1-16 y K44 marcan la historia inmutable de la 00003;
   - `gen_dict.py` no conocía 1.4.
   Las dos herramientas son de S1 y la tarea pide «make verify-s1 debe pasar».
2. Los nombres de las claves los puse yo, con los prefijos que pidió la tarea (`reports.*` y `alerts.*`, una clave por valor). `alerts.send_time` es una sola hora para «los avisos con hora», como la plantea U4b (`docs/design/U4-ajustes-b.md` l. 81 y 114) y como la da Q2-73 («hora de avisos 09:00», en singular). Qué avisos la usan lo fija S7.
3. Sobre `credit.vat_included` en filas que ya existen, la opción 1 no dice qué hacer con los NULL de 1.3. Usé el significado que 1.3 daba a esos NULL («la tasa o la cotización ya incluye IVA») y la razón de PUERTA-2 para recomendar la opción 1 («sin cambiar el significado de lo que ya hay»). No hay base de producción, así que hoy no afecta ninguna fila. La alternativa (dejarlas NULL para obligar a capturar) es una línea del UPDATE.
4. Rangos en el esquema: solo los aprobados de una sola clave (Q2-74, Q2-77, Q2-66) y los de formato (día 1 a 7 y «HH:MM» con dos dígitos). Q2-74 no dice dónde se aplica el límite de 1 a 30; lo puse también en el esquema porque es trivial y es de seguridad. Para la anticipación, la ventana de fechas y el OCR no hay rango aprobado: solo se pide 0 o más.
5. Secciones de Ajustes:
   - reportes y avisos: la de U4b (`U4-ajustes-b.md` l. 384, U4-12);
   - ventana de fechas: la de U4a (`U4-ajustes-a.md` l. 257);
   - `ocr.timeout_seconds` en «Datos y respaldos»: PROPUESTA S1, junto a la IA opcional. `docs/estado/wf1/ocr.md` l. 57 la dejó «a decidir en la puerta 2» y la puerta no la preguntó.

## Bloqueos
Ninguno nuevo. `S1-parametros-avisos` queda resuelto en el esquema (`@resuelve`, Q2-11 B).

## Pendientes
- `make verify-c1` pasa cuando C1 se regenera (solicitud 1).
- PREGUNTA para C1 y el motor de reparto (WF-3), de Q2-68 C: si después de un cambio con «solo los nuevos» se edita el importe de un gasto anterior, ¿con qué partes se recalcula? `expense_share` guarda montos y no la versión de la regla. No toca el esquema hasta que se decida.
- Q2-21 a Q2-24 (C): confirmar con S2 y S8 que el valor inicial A se siembra como partidas planificadas, sin parámetro nuevo.

## Solicitudes fuera de alcance
1. **C1** (`backend/internal/api/c1gen`, `openapi.yaml`, `C1-datos.md` §11.7):
   - agregar el tipo «hora HH:MM» en `gen2_schemas.py` `_pv` (string con patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$`);
   - exponer las 10 claves en sus secciones y los rangos (1 a 7, 1 a 31, 1 a 30 y 1 o más);
   - quitar `weekday`, `day_of_month`, `send_time` y `advance_days` de `NotificationUpdate`;
   - agregar `credit.vat_included` (C1 v1.3) y `vat` en `Credit.missing_terms` (5 datos);
   - regenerar.
2. **U4**:
   - P-26.08 y P-26.09 leen y escriben las claves de `app_setting` (variante única);
   - P-26.05 «¿la tasa incluye IVA?» se guarda en `vat_included`;
   - P-26.12 pregunta el alcance.
3. **Motor S3**: `missing_terms` con `vat`; `vat_included` decide si se suma `vat_on_interest_bp`; serie sin intereses (capital = monto) en la línea (Q2-38).
4. **S8 (WF-4)**: la semilla de Kueski carga `vat_included` = 1 (decisión C).
5. **S7 (WF-3)**: lee horas y días de `app_setting`; con 29, 30 o 31 en un mes corto, el reporte sale el último día (Q2-66 B).
6. **Criterios**: C1-22, U4-12 y S7 dejan de marcar estos cuatro parámetros como PROVISIONAL.
7. **Integrador**: el Makefile dice «esquema 1.3» en `verify-s1` (l. 8, 60 y 139). Es cosmético.
8. **Bloqueos**: el contexto de `U4-reparto-retroactivo.md` dice que los gastos «no guardan una copia de las partes», pero `expense_share` sí las congela.

## Propuestas para Axel
- Confirmar que «Tiempo máximo del OCR» vive en Ajustes › «Datos y respaldos», junto a la IA opcional. La otra opción que dejó `ocr.md` es «Tiendas y precios».
- Confirmar el supuesto 3: los créditos que ya existan sin porcentaje de IVA pasan como «incluye IVA», en lugar de «falta capturar».
- Si quiere rangos de seguridad para el OCR, la anticipación de precios viejos o la ventana de fechas, hay que aprobarlos; hoy solo se pide 0 o más.

## Consumo aproximado
- Acciones de herramienta: 46.
- KB leídos: ~260.
- KB escritos de mi autoría: ~75 (migración 22, este reporte 12, parches de las herramientas ~18, reglas ~10 y texto fijo del diccionario ~6). Además, el generador reescribe el diccionario completo (190 KB) con un cambio neto de ~12 KB.
