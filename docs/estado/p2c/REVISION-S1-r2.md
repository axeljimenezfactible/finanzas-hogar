# P2c · REVISION-S1-r2 · Esquema 1.4 (migración 00004), ronda 2 de 3

Revisor independiente `revisor-S1-r2`, flujo P2c «Aplicación de la puerta 2», 11-oct-2026, noche. No escribí nada de lo que reviso y no corregí nada.
- Artefacto: `backend/db/migrations/00004_puerta2.sql`, `docs/contratos/diccionario-datos.md` 1.4 (partes nuevas), `docs/contratos/reglas-negocio.md` 1.4 (RN tocadas), `docs/estado/p2c/S1-puerta2.md` y `docs/estado/p2c/correccion-S1-r1.md`.
- Contra qué: `docs/criterios/S1.md`, decisiones K44 y K46 a K50, PUERTA-2 §3 y §4 y RESPUESTAS.md.
- Herramientas: goose v3.28.0 (`/home/deepsec/finanzas-hogar-orq/gobin/goose`), sqlite3 3.53.4 (CLI y módulo de Python) y mi script `datos/tmp/rev-s1-r2/chk.py` (`python3 -I`). Lo borré al terminar.

## Veredicto
**NO APROBADO: 1 alto, 2 medios y 5 bajos.**
- **Alto que sigue abierto:** H-01. `make verify-c1`, y con él `make verify-all`, siguen fallando con `gen2: tipo de parámetro desconocido: hora HH:MM`.
- **Cerrados:** M-01 (CHECK de `vat_included`) y M-02 (bloqueo `S1-reparto-edicion-gasto` con sus marcas PROVISIONAL). Los comprobé con pruebas propias.
- **Medios nuevos:**
  - N-01: la tabla «Ajustes: 14 secciones» (S1-13) no tiene `credit.vat_included`.
  - N-02: el esquema acepta (`vat_included`, `vat_on_interest_bp`) = (0, NULL) y ningún contrato de S1 dice qué significa. RN-26 y la invariante 16 llevan a resultados distintos.
- **Lo que sí pasa con mis comandos:**
  - verificaciones (a) a (f) del encargo;
  - valores y rangos de las 10 claves nuevas;
  - traspaso y Down de `params_json`;
  - traspaso de créditos de 1.3;
  - 00001 a 00003 intactas;
  - diccionario idéntico al de `gen_dict.py`.

## Verificaciones (evidencia propia)
| # | Qué | Comando | Salida | Estado |
|---|---|---|---|---|
| a1 | SQL limpio | `make verify-sql GOOSE=…` | 4 Up aplicados; `tablas=55 (STRICT=55) · índices=152 · disparadores=79 · foreign_key_check=0 filas · integrity_check=ok`; Down en orden inverso deja 0 objetos; goose up, down-to 0 y up sin errores. rc=0 | OK |
| a2 | S1 | `make verify-s1 GOOSE=…` | `RESUMEN TODO OK`, «aplicó las 4 migraciones del comodín». K47, AVISOS14, PROG14, DOWN14, DDL y GOOSE OK. rc=0 | OK |
| a3 | goose con datos | `goose … up-to 3`; 2 créditos de 1.3 (1600 y NULL, `updated_at` fijo); `up`; `down-to 3`; `up` | rc 0/0/0/0. Tras up: `T1\|0\|1600\|2026-01-01T00:00:00.000Z`, `T2\|1\|\|2026-01-02T00:00:00.000Z`; versión 4; `sqlite_temp_master` = 0. Tras down-to 3: porcentajes y `updated_at` de 1.3, y `params_json` `{"weekday": 7}`, `{"day_of_month": 1}`, `{"weekday": 7}`. Tras up: integrity ok, fk vacío, 30 claves | OK |
| b1 | `@param` de la 00004 | `chk.py` B-campos, B-claves y B-valores | 10 `@param` con sus 7 campos llenos. Claves = las 10 esperadas. Valor del `@param` = valor sembrado = PUERTA-2 §4 y decisión F: daily "21:30", weekly 7 y "19:00", monthly 1 y "08:00", alerts "09:00", 1 y 7, `optimizer.date_window_days` 7, `ocr.timeout_seconds` 30. Decisión citada: Q2-72, Q2-73, F, Q2-75 y Q2-76. Tipo «hora HH:MM» solo en las 4 horas. `user_account_id` NULL y descripción presente. Diferencias = [] | OK |
| b2 | Claves que ya existían | `chk.py` | `comparator.timeout_seconds` 5 (Q2-74, «mantener 5 s»), `price.stale_after_days` 30 y `price.zero_weight_after_days` 90. `app_setting`: 20 claves en 1.3 y 30 en 1.4 | OK |
| b3 | Rangos y formato (INSERT y UPDATE, con rollback) | `chk.py` B-rangos | 52 aceptables y 68 no aceptables, 0 discrepancias. Ejemplos: comparador 0, 31, "5" y 5.0 rechazados, 1 y 30 aceptados; día de la semana 0 y 8 rechazados; día del mes 32 y 1.5 rechazados; horas "24:00", "9:00", "21:60", "29:00", 2130, "21:30 " y "21:3" rechazadas, "00:00" y "23:59" aceptadas | OK |
| b4 | Regla entre dos claves (Q2-77) en el motor | `grep -A12 'Validate' backend/internal/precios/types.go` | l. 109: `if s.StaleAfterDays < 1 \|\| s.ZeroWeightAfterDays < s.StaleAfterDays` → `ErrInvalidArgument`, como dicen RN-12 y la 00004 | OK |
| c1 | `params_json` tras Up | `chk.py` C-up | Ninguna fila trae `send_time`, `advance_days`, `weekday` ni `day_of_month`: budget_threshold `{"scopes": ["category", "total"]}` y las demás `{}`. `updated_at` conservado = True | OK |
| c2 | Down (copia de una base 1.3 → Up → Down) | `chk.py` C-down y C-down-esquema | `notification` (code, params_json y updated_at) igual byte por byte a 1.3. `sqlite_master` = 1.3 (244 objetos). `app_setting` igual | OK |
| c3 | Down y Up con valores que no son los iniciales | `chk.py` C-down-valores y C-up-valores | Down con 5, 3 y 31 → `{"weekday": 5}`, `{"weekday": 3}` y `{"day_of_month": 31}`. Up desde 1.3 con `{"weekday": 4, "extra": 1}`, 29 y 9 (inválido) → `reports.weekly_weekday` 4, `monthly_day_of_month` 29 y `meal_plan_missing_weekday` 7 (el de F); `params_json` `{"extra":1}` y `{}` | OK |
| d1 | `vat_on_interest_bp` sin cambio | `chk.py` D-cols y D-ddl | `PRAGMA table_xinfo(credit)`: las columnas de 1.3 son iguales (`vat_on_interest_bp` INTEGER, NULL admitido) y solo se agrega `vat_included` al final. La línea DDL de `vat_on_interest_bp` es idéntica | OK |
| d2 | Matriz del CHECK (U4-iva-credito opción 1) | `chk.py` D-check, 13 casos × INSERT/UPDATE | Rechaza con «CHECK constraint failed»: (NULL, 1600), (1, 1600), (1, 0), (NULL, 0), (2, NULL), (−1, NULL) y (0, 10001). Acepta (0, 1600), (0, 0), (1, NULL), (0, NULL) y (NULL, NULL). La única «discrepancia» fue mi caso ('1', NULL): la tabla STRICT convierte el texto '1' en el entero 1 (`typeof` = integer). Es lo normal de SQLite y no es defecto | OK (ver N-02 sobre (0, NULL)) |
| d3 | Traspaso de 1.3 con datos (Python) | `chk.py` D-traspaso y D-down | (1600) → (0, 1600); NULL → (1, NULL); (0) → (0, 0). `updated_at` conservado; tabla temporal borrada (0). Down: columnas de 1.3 iguales | OK |
| e | 00001 a 00003 intactas | `git diff --stat HEAD -- 00001… 00002… 00003…` | vacío, 0 líneas de diff | OK |
| f | Diccionario = generador | `python3 -I backend/db/tools/gen_dict.py datos/tmp/rev-s1-r2/dic1.md`; `cmp` contra el repositorio | «versión 1.4 · parámetros 25 · sin valor 3 · valores iniciales que no coinciden con lo sembrado: []», 55 tablas, 152 índices, 79 disparadores, 0 columnas sin descripción. `CMP_REPO_IDENTICO` | OK |
| g1 | C1 | `make verify-c1 GOOSE=…` | `gen2: tipo de parámetro desconocido: hora HH:MM` / `make: *** [Makefile:130: verify-c1] Error 1`. rc=2 | **FALLA** (H-01) |
| g2 | Verificación global | `make -k verify-all GOOSE=…` | rc=2. Pasan WF-0 (SQL, Go, colorimetría y OpenAPI), S1 (`RESUMEN TODO OK`), U1 (`RESULTADO: TODO OK`), criterios y WF-1 («verificación completa OK»). Falla solo verify-c1, con el mismo mensaje. `git status`: `backend/internal/api/c1gen/` sin cambios | **FALLA** (H-01) |
| h | Marcas PROVISIONAL (S1-16) | `grep -rn PROVISIONAL backend/db diccionario reglas` | Diccionario l. 2360 (invariante 18) y RN-03 l. 103 citan `S1-reparto-edicion-gasto` (ABIERTO, l. 33). Las de la 00003 citan `S1-parametros-avisos`, ya resuelto; la 00004 las declara historia con `@resuelve` (B-03 de r1) | OK |
| i | Mapa de Ajustes (S1-13) | `grep -c vat_included` en el diccionario; `sed -n 2331p` | `vat_included` aparece en 4 líneas (10, 770, 782 y 2358) y **no** en la fila 5 «Créditos y condiciones» (l. 2331), que sí trae `.vat_on_interest_bp`. Las 10 claves nuevas sí están en sus filas | **FALLA** (N-01) |

## Hallazgos

**H-01 (alta, sigue abierto) · Verificación global rota: `make verify-c1` y `make verify-all` fallan por el tipo «hora HH:MM».**
- Ruta: `backend/db/migrations/00004_puerta2.sql` l. 62, 64, 66 y 67 (`@param … | hora HH:MM | …`). Quien lo consume es `backend/internal/api/c1gen/gen2_schemas.py` l. 240 (`_pv`), que solo conoce «entero» y «lista de enteros».
- Evidencia: g1 y g2. `gen2_schemas.py` sigue sin cambios en el árbol, y ningún reporte de `docs/estado/p2c/` dice que un agente de C1 tomara la solicitud 1 de `S1-puerta2.md`.
- Impacto: el commit «[Puerta 2]» (Q2-97, K52) no puede cerrar con `make verify-all`.
- Corrección sugerida (no cambia):
  - el orquestador asigna en P2c la solicitud 1 de `S1-puerta2.md` a un agente de C1: tipo «hora HH:MM» en `_pv`, como string con patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$`, las 10 claves con sus rangos, `credit.vat_included` y `vat` en `Credit.missing_terms`; luego regenera;
  - REVISION-bloqueos-r3 avisa que, sin el error de gen2, `verify_c1.py` da además FALLA K44 y FALLA C1-15: ese agente tiene que cerrarlos también;
  - después, `make verify-all GOOSE=…` en la ronda 3;
  - las horas no pasan a entero (K43, Q2-72 y Q2-73).

**N-01 (media) · S1-13: la fila «Créditos y condiciones» del mapa de Ajustes no tiene `credit.vat_included`.**
- Ruta: `docs/contratos/diccionario-datos.md` l. 2331 (sección «Ajustes: 14 secciones → campos → tabla.columna», l. 2322).
- Evidencia: verificación i.
  - La fila trae `.vat_on_interest_bp` y `.rate_base`, pero no la columna nueva.
  - La propia columna dice `@used_new Ajustes «Créditos y condiciones» (P-26.05, «¿la tasa incluye IVA?»)` (00004 l. 89).
  - U4 P-26.05 (`docs/design/U4-ajustes-a.md` l. 217) tiene el campo «¿la tasa incluye IVA? y, si no, IVA que se suma | Opción sí o no y porcentaje».
- Impacto:
  - S1-13 pide que cada campo de Ajustes tenga su `tabla.columna`. La parte «sí o no» del campo ya no vive en `vat_on_interest_bp` sino en `vat_included`, y el mapa no lo dice;
  - C1 y U4 tendrían que deducir el destino;
  - el chequeo S1-13 de `s1_verifica.py` no lo detecta.
- Corrección sugerida:
  - agregar `.vat_included` (00004, Q2-67 opción 1) a la fila 5, junto a `.vat_on_interest_bp`, en la fuente que `gen_dict.py` conserva, y regenerar;
  - de preferencia, que S1-13 de `s1_verifica.py` compruebe que cada columna con `@used`/`@used_new` «Ajustes «X»» aparece en la fila X.

**N-02 (media) · (`vat_included`, `vat_on_interest_bp`) = (0, NULL) lo acepta el esquema, pero ningún contrato de S1 dice qué significa, y RN-26 y la invariante 16 llevan a resultados distintos.**
- Ruta:
  - `backend/db/migrations/00004_puerta2.sql` l. 95: el CHECK acepta (0, NULL), según d2;
  - `docs/contratos/reglas-negocio.md` RN-26 (l. 278): «`vat` mientras `vat_included` sea NULL»;
  - `docs/contratos/diccionario-datos.md` invariante 16 (l. 2358): «con su porcentaje en `vat_on_interest_bp` cuando no lo incluye … si falta alguno, S3 muestra «falta capturar»»;
  - nota de la columna, l. 782: «solo lleva valor con 0». Dice cuándo puede tener valor, no que deba tenerlo.
- HECHO: el encabezado de la 00004 (l. 39 y 40) reconoce que ese par existe («un crédito con `vat_included` = 0 y sin porcentaje vuelve a leerse "incluye IVA" en 1.3»), pero no dice qué es en 1.4. U4 captura «sí o no» y el porcentaje por separado, así que el par puede ocurrir: se elige «no» y todavía no se escribe el porcentaje.
- Impacto:
  - Con RN-26, `missing_terms` no lista `vat` y S3 calcula con el porcentaje en NULL. `S3-interfaz.md` l. 116 lee NULL como «la tasa ya trae IVA», o sea, suma 0 % de IVA en una tarjeta que, según Axel, no lo incluye.
  - Con la invariante 16, falta un dato y debería decir «falta capturar».
  - C1 y S3 tendrían que adivinar entre esas dos lecturas o rechazar el par con 422. K18 prohíbe inventar.
- Corrección sugerida (sin decidir por Axel): escribir en RN-26, en la nota de `vat_included` y en la invariante 16 que con `vat_included` = 0 y `vat_on_interest_bp` NULL falta el porcentaje y `missing_terms` lista `vat`, y sumar ese caso a K47.
  - Eso sale de K18 («no inventa tasa ni mínimo»). No hace falta columna.
  - Si el corrector lo considera una decisión nueva, que lo marque PROVISIONAL con un bloqueo (por ejemplo, el mismo `U4-iva-credito` reabierto en ese punto, o uno nuevo) en lugar de dejarlo implícito.
  - En los dos casos, avisar a C1 y a S3 por solicitud.

**B-05 (baja, nuevo) · `S1-puerta2.md` §Hecho (c) (l. 29) sigue citando el CHECK de la ronda 0, `(vat_included = 0 OR vat_on_interest_bp IS NULL)`, que aceptaba (NULL, 1600).**
- El CHECK vigente está en `correccion-S1-r1.md` y en la 00004 l. 95.
- Corrección: una nota en esa línea que remita a la corrección, o el texto vigente.

**B-06 (baja, nuevo) · `docs/estado/bloqueos/INDICE.md` todavía no lista `S1-reparto-edicion-gasto`.**
- Es la solicitud 2 de `correccion-S1-r1.md`, fuera del ALCANCE de S1.
- Corrección: el integrador o el agente de bloqueos lo agrega en «Bloqueos nuevos de P2c» y actualiza el conteo.

## Estado de los anteriores (REVISION-S1-r1)
| Hallazgo | Estado | Evidencia propia |
|---|---|---|
| H-01 (alta) verificación global rota | **ABIERTO** | g1 y g2: el mismo mensaje y `gen2_schemas.py` sin cambios |
| M-01 (media) CHECK de `vat_included` | **CERRADO** | d2: (NULL, 1600) se rechaza con INSERT y con UPDATE. d3 y a3: el traspaso de 1.3 y goose con datos funcionan con el CHECK nuevo; la tabla temporal se borra. a2: K47 tiene 19 negativas |
| M-02 (media) Q2-68 C sin bloqueo | **CERRADO** | `docs/estado/bloqueos/S1-reparto-edicion-gasto.md` ABIERTO, con HECHO, qué falta, opciones A, B y C sin recomendar, impacto y la opción en uso. PROVISIONAL en RN-03 l. 103 y en la invariante 18 l. 2360, citando el bloqueo. Queda solo la entrada en INDICE (B-06) |
| B-01 (baja) celda de `vat_on_interest_bp` con la frase de 1.3 | sigue (l. 770) | `grep -n vat_on_interest_bp` |
| B-02 (baja) sección del OCR «propuesta S1» | sigue | `@param ocr.timeout_seconds` «sección: propuesta S1» |
| B-03 (baja) marcas históricas de la 00003 frente a la letra de S1-16 | sigue | h |
| B-04 (baja) `gen_dict.py` y `s1_verifica.py` editados fuera del ALCANCE | sigue; ahora con 3 cambios más de `correccion-S1-r1` (Supuesto 1) | `git diff --stat`: `s1_verifica.py` +110/− |

## Solicitudes fuera de alcance
1. **Orquestador → agente de C1:** H-01 (solicitud 1 de `S1-puerta2.md` y FALLA K44 y C1-15 de `verify_c1.py`). Con N-02, `Credit.missing_terms` incluye `vat` también con (0, NULL), o 422.
2. **Orquestador → S3 (contratos):** `S3-interfaz.md` l. 116, 257, 272, 884, 903 y 919 siguen marcando PROVISIONAL con `U4-iva-credito`, que ya está resuelto (K46, Q2-67 opción 1). No es de S1, pero le afecta N-02.
3. **Integrador:** B-06 (INDICE) y ratificar B-04.
