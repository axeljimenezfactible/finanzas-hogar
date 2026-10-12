# P2c · REVISION-S1-r3 · Esquema 1.4 (migración 00004), ronda 3 de 3

Revisor independiente `revisor-S1-r3`, flujo P2c «Aplicación de la puerta 2», 11-oct-2026, noche. No escribí nada de lo que reviso y no corregí nada.
- **Qué revisé:**
  - `backend/db/migrations/00004_puerta2.sql`, completa (221 líneas);
  - `docs/contratos/diccionario-datos.md` 1.4: las partes nuevas y las que cambió `correccion-S1-r2`;
  - `docs/contratos/reglas-negocio.md` 1.4: RN-03, RN-07, RN-11, RN-12, RN-26 y RN-29;
  - los reportes `docs/estado/p2c/S1-puerta2.md`, `correccion-S1-r1.md` y `correccion-S1-r2.md`.
- **Contra qué:** `docs/criterios/S1.md`, las decisiones K44 y K46 a K50 (más §F y K52), PUERTA-2 §3 y §4, y `RESPUESTAS.md`.
- **Herramientas:**
  - goose v3.28.0 (`/home/deepsec/finanzas-hogar-orq/gobin/goose`) y sqlite3 3.53.4, en la CLI y en el módulo de Python;
  - mis scripts `datos/tmp/rev-s1-r3/chk.py` (esquema, semilla, Down y CHECK) y `chk2.py` (marcas PROVISIONAL y mapa de Ajustes), corridos con `python3 -I`. Los borré al terminar.
- **Base:** HEAD 1ccb468 más el árbol de trabajo a las 23:30. `correccion-S1-r2` dejó la 00004 a las 23:26 y el diccionario y las reglas a las 23:25 y 23:26.

## Veredicto
**NO APROBADO: 1 alto, 0 medios y 7 bajos.**
- **Alto que sigue abierto: H-01.**
  - `make verify-c1`, y con él `make verify-all`, siguen fallando con `gen2: tipo de parámetro desconocido: hora HH:MM`.
  - Si se salta ese error, `verify_c1.py` da FALLA K44 y FALLA C1-15.
  - El arreglo está en C1, fuera del ALCANCE de S1: `backend/internal/api/c1gen/` no tiene cambios en el árbol.
  - Esta es la última ronda, así que el cierre pasa al orquestador o al integrador (ver «Solicitudes»).
- **Cerrados con mis pruebas:**
  - N-01: la fila 5 del mapa de Ajustes ya trae `.vat_included`.
  - N-02: el par (0, NULL) dice lo mismo en RN-26, en la invariante 16, en la nota de la columna, en `@used_new` y en el encabezado. Es PROVISIONAL y cita el bloqueo abierto `S3-iva-lectura`, pregunta 2, que no tiene recomendación.
- **El artefacto de S1 cumple sus criterios:**
  - pasan las verificaciones (a) a (f) del encargo, `make verify-s1` y `make verify-sql`;
  - las 10 claves nuevas tienen el valor de PUERTA-2 §4 y de la decisión F;
  - las marcas PROVISIONAL citan bloqueos abiertos.

## Verificaciones (evidencia propia)
| # | Qué | Comando | Salida | Estado |
|---|---|---|---|---|
| a1 | SQL limpio | `make verify-sql GOOSE=…` | rc=0. 4 Up aplicados; `tablas=55 (STRICT=55) · índices=152 · disparadores=79 · foreign_key_check=0 filas · integrity_check=ok`. El Down en orden inverso deja 0 objetos. Up otra vez: 55. goose up, down-to 0 y up sin errores | OK |
| a2 | S1 | `make verify-s1 GOOSE=…` y `s1_verifica.py` directo (líneas completas) | rc=0, `RESUMEN TODO OK`, «aplicó las 4 migraciones del comodín». Pasan S1-13 (25 de 25 claves `@param` en la fila de su sección), S1-16 (12 líneas PROVISIONAL), PROG14, K47 (19 negativas y 11 positivas), DDL, DOWN14 (244 objetos, 0 diferencias) y GOOSE | OK |
| a3 | goose con datos | `goose up-to 3`; persona y 2 tarjetas de 1.3 (1600 y NULL, `updated_at` fijo); `up`; `down-to 3`; `up` | rc 0/0/0/0. Tras up: `T1\|0\|1600\|2026-01-01T00:00:00.000Z` y `T2\|1\|N\|2026-01-02T00:00:00.000Z`, versión 4, 30 claves, params `{}` e integrity ok. Tras down-to 3: `T1\|1600`, `T2\|N` con su `updated_at`, 20 claves, `{"weekday": 7}`, `{"day_of_month": 1}`, `{"weekday": 7}` y sin columna `vat_included`. Tras up: ok, 30 claves | OK |
| b1 | `@param` de la 00004 | `chk.py` B-param | 10 `@param` con sus 7 campos llenos y ninguna clave sobrante ni faltante. Tipo: «hora HH:MM» en las 4 horas y «entero» en el resto. El valor del `@param` es igual al sembrado, a PUERTA-2 §4 y a la decisión F: daily `"21:30"`, weekly 7 y `"19:00"`, monthly 1 y `"08:00"` (Q2-72); alerts `"09:00"` y 1 (Q2-73) y 7 (F, «domingo»); `optimizer.date_window_days` 7 (Q2-75) y `ocr.timeout_seconds` 30 (Q2-76). La decisión citada y la sección coinciden. `diferencias=[]` | OK |
| b2 | Claves que ya existían y descripciones | `chk.py` B-otras y B-desc | `comparator.timeout_seconds` 5 (Q2-74: «mantener 5 s»), `price.stale_after_days` 30 y `price.zero_weight_after_days` 90. `app_setting`: 20 claves en 1.3 y 30 en 1.4. 0 claves nuevas sin descripción | OK |
| b3 | Rangos y formato (UPDATE e INSERT, con rollback) | `chk.py` B-rangos: 35 casos en los dos modos | 24 aceptables y 46 no aceptables, 0 discrepancias. Rechazados: comparador 0, 31, `"5"` y 5.0; precio viejo 0; día de la semana 0 y 8; día del mes 0, 32 y 1.5; anticipación, ventana y OCR −1; horas `"24:00"`, `"29:00"`, `"9:00"`, `"21:60"`, `2130`, `"21:30 "`, `"21:3"`, `21:30` sin comillas y `"2a:00"`. Aceptados: 1 y 30, día 31, 0 en anticipación y ventana, `"00:00"` y `"23:59"` | OK |
| b4 | RN-29 contra el `@param` | lectura de RN-29 (l. 293 en adelante) | Las 10 filas tienen el mismo valor, rango, decisión, sección y migración. Q2-66 B está escrito en «Día del mes en un mes corto». Ejemplo: (120 − 60) / 91 = 0.6593 y 30 / 61 = 0.4918, correctos | OK |
| c1 | `params_json` tras Up (en la misma base) | `chk.py` C-up y C-semilla | Ninguna fila trae `send_time`, `advance_days`, `weekday` ni `day_of_month`. Pasan a `{}` las tres que los tenían, y `budget_threshold` conserva `{"scopes": ["category", "total"]}`. `updated_at` conservado = True. El Up no tiene `json_set`, `json_insert`, `"send_time"` ni `"advance_days"` en código | OK |
| c2 | Down sobre la misma base | `chk.py` C-down | `sqlite_master` igual a 1.3 (244 objetos). `notification` (code, params_json y updated_at) igual byte por byte. `app_setting` igual | OK |
| c3 | Down y Up con valores que no son los iniciales | `chk.py` C-down-valores y C-up-valores | Down con 5, 3 y 31 → `{"weekday": 5}`, `{"weekday": 3}` y `{"day_of_month": 31}`. Up desde 1.3 con `{"weekday": 4, "extra": 1}`, 29 y 9 (inválido) → 4, 29 y 7 (el de F); params `{"extra":1}`, `{}` y `{}` | OK |
| d1 | `vat_on_interest_bp` sin cambio | `chk.py` D-cols y D-ddl | En `PRAGMA table_xinfo(credit)` las columnas de 1.3 son idénticas, con `vat_on_interest_bp` INTEGER y NULL admitido. La única nueva es `vat_included`, al final. La línea DDL de `vat_on_interest_bp` es idéntica en 1.3 y en 1.4 | OK |
| d2 | Matriz del CHECK (U4-iva-credito, opción 1) | `chk.py` D-check: 13 casos, con UPDATE y con INSERT | Acepta (0, 1600), (0, 0), (1, NULL), (0, NULL) y (NULL, NULL). Rechaza (NULL, 1600), (1, 1600), (1, 0), (NULL, 0), (2, NULL), (−1, NULL), (0, 10001) y (0, −1). 0 discrepancias | OK |
| d3 | Traspaso de 1.3 y Down | `chk.py` D-traspaso y D-down | 1600 → (0, 1600); NULL → (1, NULL); 0 → (0, 0). `updated_at` conservado; tabla temporal = 0; integrity ok; fk vacío. Down: columnas de 1.3 y porcentajes de vuelta | OK |
| e | 00001 a 00003 intactas | `git diff --stat HEAD -- 00001… 00002… 00003…` y `git diff … \| wc -l` | vacío; 0 líneas | OK |
| f | Diccionario = generador | `python3 -I backend/db/tools/gen_dict.py datos/tmp/rev-s1-r3/dic.md` (dos veces) y `cmp` | «parámetros sin fila en «Ajustes: 14 secciones»: []» · «versión 1.4 · parámetros 25 · sin valor 3 · valores iniciales que no coinciden con lo sembrado: []» · 55 tablas, 152 índices y 79 disparadores · «columnas sin descripción: []». `CMP_IDENTICO` contra el repositorio; idempotente | OK |
| g1 | C1 | `make verify-c1 GOOSE=…` | `gen2: tipo de parámetro desconocido: hora HH:MM` / `make: *** [Makefile:130: verify-c1] Error 1`. rc=2 | **FALLA** (H-01) |
| g2 | C1 sin el error de gen2 | base con las 4 Up (awk del Makefile) y `uv run … verify_c1.py --db` | rc=1. **FALLA K44**: 10 claves sin lectura, cambio y restablecer en C1 (las 10 de la 00004). **FALLA C1-15**: 6 marcas PROVISIONAL de `openapi.yaml` (`Notification.params_json`, `send_time` y `advance_days`) y líneas de `C1-datos.md` (P-26.08 y P-26.09) sin bloqueo abierto, porque citan `S1-parametros-avisos`, ya resuelto. `RESULTADO: 2 falla(s)` | **FALLA** (H-01) |
| g3 | Verificación global | `make -k verify-all GOOSE=…` | rc=2. Pasan WF-0 (SQL, Go, colorimetría y OpenAPI), S1 (`RESUMEN TODO OK`), U1 (`RESULTADO: TODO OK`), criterios («INDICE contra conteo propio: 0 diferencias») y WF-1 («verificación completa OK»). Falla solo verify-c1. `git status`: `backend/internal/api/c1gen/` sin cambios (`gen2_schemas.py` es de las 16:47) | **FALLA** (H-01) |
| h | Marcas PROVISIONAL (S1-16) | `chk2.py`: cada línea con PROVISIONAL de las migraciones, del diccionario y de las reglas, cruzada con `docs/estado/bloqueos/` | 00004 l. 23, 44, 93 y 99; diccionario l. 751, 782 y 2358; y reglas l. 279 citan `S3-iva-lectura` (sin «Respuesta de Axel», **Estado: ABIERTO**). Diccionario l. 2360 y reglas l. 103 citan `S1-reparto-edicion-gasto` (ABIERTO). En la 00003, l. 56 a 62 citan `S1-parametros-avisos` (RESUELTO): son historia por `@resuelve` de la 00004. Diccionario l. 21, reglas l. 9 y 00003 l. 11 son definiciones o la continuación de un bloque. Ninguna línea de 1.4 cita solo bloqueos resueltos | OK (ver B-07) |
| i | Mapa de Ajustes (S1-13, N-01) | `sed -n 2331p`; `chk2.py`, que busca cada `@param` en su fila y cada columna de un `@alter` con «Ajustes «X»» en la fila X | La fila 5 trae `` `.vat_included` (00004, Q2-67 opción 1: …; NULL = falta capturar), `.vat_on_interest_bp` (… solo con `.vat_included` = 0) ``. Las 10 claves están en su fila. Mi script marcó también `attachment.ocr_draft_*` (00002) en «Notificaciones», pero es un falso positivo de mi lector: el `@used_new` de la l. 32 de la 00002 no nombra Ajustes. Lo revisé a mano | OK |
| j | N-02: coherencia | lectura de 00004 l. 19 a 23, 37 a 44, 93 y 99; diccionario l. 751, 770, 782 y 2358; RN-26 l. 278 y 279; `S3-iva-lectura.md`; `S3-interfaz.md` l. 116, 258, 311 y 885 | Todos dicen lo mismo de (0, NULL): falta el porcentaje, `Credit.missing_terms` lista `vat` y S3 responde «falta capturar» (`ErrMissingCapture`), sin 0 % y sin inventar (K18). Es PROVISIONAL y cita `S3-iva-lectura` pregunta 2 (opciones A y B, «Sin recomendación», ABIERTO). Coincide con lo que S3 1.5 ya usa. Con B cambia el esquema o C1. No decide por Axel (K46) | OK |

## Hallazgos

### H-01 · ALTA · sigue abierto · Verificación global rota (Q2-97 y K52 piden cerrar con `make verify-all`)
- **Ruta:**
  - `backend/db/migrations/00004_puerta2.sql` l. 66, 68, 70 y 71: `@param … | hora HH:MM | …`. En r1 y r2 eran las l. 62, 64, 66 y 67; se movieron por el encabezado nuevo.
  - Quien lo consume: `backend/internal/api/c1gen/gen2_schemas.py` l. 240 a 244 (`_pv`). Solo conoce «entero» y «lista de enteros», y lo demás termina en `SystemExit('gen2: tipo de parámetro desconocido: ' + typ)`.
- **Evidencia:** g1, g2 y g3.
  - `gen2_schemas.py` sigue sin cambios.
  - Ningún reporte de `docs/estado/p2c/` dice que un agente de C1 tomara la solicitud 1 de `S1-puerta2.md`. Están S1, contratos, bloqueos, S5-términos y sus revisiones.
  - Sin el error de gen2 quedan FALLA K44 (10 claves) y FALLA C1-15 (6 marcas en `openapi.yaml` y líneas de `C1-datos.md` que citan `S1-parametros-avisos`, ya resuelto).
- **Impacto:** el commit «[Puerta 2]» no puede cerrar con `make verify-all` (Q2-97 A, K52).
  - Es lo único que impide aprobar.
  - No es un defecto del esquema: las horas como texto «HH:MM» las piden K43, Q2-72 y Q2-73, y K47 las pone en `app_setting`.
- **Corrección sugerida** (fuera del ALCANCE de S1; no corresponde a otra ronda de S1). El orquestador asigna en P2c un agente de C1 que:
  - agregue el tipo «hora HH:MM» a `_pv`, como string con patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$`, y no como entero;
  - exponga las 10 claves con lectura, cambio y restablecer, con sus rangos (K44 de `verify_c1.py`);
  - quite o rehaga las 6 marcas PROVISIONAL de `openapi.yaml` y las de `C1-datos.md` que citan `S1-parametros-avisos`. `Notification.params_json` ya no lleva `send_time` ni `advance_days` (K47). Con eso cierra C1-15;
  - agregue `credit.vat_included` y `vat` en `Credit.missing_terms`, con el caso (0, NULL) como PROVISIONAL `S3-iva-lectura` (solicitud 1 de `correccion-S1-r2.md`);
  - regenere y corra `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`, que debe dar rc=0.

### B-01 · baja · sigue · La celda de `credit.vat_on_interest_bp` empieza con la frase de 1.3
- **Ruta:** `docs/contratos/diccionario-datos.md` l. 770, que sale del comentario DDL de la 00001 y de la `@col_note` de la 00004 l. 85.
- **Qué dice:** empieza con «NULL = la tasa o la cotización ya incluye IVA (Kueski)». Más adelante remite a `vat_included` y a RN-26, pero no menciona el par (0, NULL), que en 1.4 significa «falta el porcentaje».
- **Corrección:** en la `@col_note` de la 00004, decir que en 1.4 NULL depende de `vat_included` (1 = ya incluye; NULL = falta capturar; 0 = falta el porcentaje, PROVISIONAL `S3-iva-lectura`) y regenerar.

### B-02 · baja · sigue · El OCR tiene la sección «propuesta S1»
- **Ruta:** `@param ocr.timeout_seconds` (00004 l. 75) y RN-29: «Datos y respaldos (propuesta S1)».
- La sección no la aprobó nadie. Se resuelve cuando Axel la vea (U4 o la puerta 4).

### B-03 · baja · sigue · Las marcas «provisional» de la 00003 citan un bloqueo resuelto
- **Ruta:** 00003 l. 56 a 62.
- Quedan como historia por `@resuelve S1-parametros-avisos` (00004 l. 61), porque la 00003 no se edita. Cumple el espíritu de S1-16, pero no su letra.

### B-04 · baja · sigue · Herramientas editadas fuera del ALCANCE
- **Ruta:** `backend/db/tools/gen_dict.py` (+44) y `s1_verifica.py` (+110), los mismos números que vio r2.
- El integrador las ratifica.

### B-05 · baja · sigue · `S1-puerta2.md` cita el CHECK de la ronda 0
- **Ruta:** `docs/estado/p2c/S1-puerta2.md` §Hecho (c), l. 29: `(vat_included = 0 OR vat_on_interest_bp IS NULL)`.
- El CHECK vigente está en la 00004 l. 99 y en `correccion-S1-r1.md`.
- **Corrección:** una nota que remita a la corrección.

### B-06 · baja · sigue y se amplía · `INDICE.md` no lista 2 bloqueos abiertos que cita S1 1.4
- **Ruta:** `docs/estado/bloqueos/INDICE.md`. `grep` de `S1-reparto-edicion-gasto` y de `S3-iva-lectura` da 0 coincidencias. La carpeta tiene 79 archivos sin contar el INDICE.
- «Bloqueos nuevos de P2c» (l. 398) y K46 cuentan 4 nuevos (`S2-ahorro-quincenal`, `S2-metas-vencida`, `S4-presentacion-respaldo` y `S5-carritos-granel`).
- S1 1.4 cita `S3-iva-lectura` en 8 líneas y `S1-reparto-edicion-gasto` en 2.
- `S3-iva-lectura.md`, en «Cómo seguí», tampoco dice que S1 1.4 escribe su pregunta 2 como PROVISIONAL (solicitud 4 de `correccion-S1-r2`).
- **Corrección** (integrador o agente de bloqueos): dar de alta los dos en el INDICE y en el conteo de K46, y completar «Cómo seguí».

### B-07 · baja · nuevo · Huecos de `s1_verifica.py` que hoy cubren mis scripts
- **Ruta:** `backend/db/tools/s1_verifica.py`.
- **S1-16:** informa «bloqueos abiertos citados: ['S1-reparto-edicion-gasto']» y no reconoce `S3-iva-lectura`, porque el patrón solo admite S1, C1, U1, U2, S6 y criterios. Da OK porque las líneas dicen «bloqueo abierto», no porque compruebe el archivo. Mi `chk2.py` sí lo comprueba: h.
- **K47:** no tiene la positiva (0, NULL); la cubre mi d2.
- **S1-13:** no compara por columna; la cubre mi i.
- **Corrección:** las solicitudes 3 (a), (b) y (c) de `correccion-S1-r2.md`, a cargo del dueño de la herramienta.

## Estado de los anteriores
| Hallazgo | Estado | Evidencia propia |
|---|---|---|
| H-01 (alta) verificación global rota | **ABIERTO** | g1, g2 y g3: el mismo mensaje; `gen2_schemas.py` sin cambios; además FALLA K44 y FALLA C1-15 |
| N-01 (media) S1-13: falta `credit.vat_included` en la fila 5 | **CERRADO** | i: la fila 5 trae `.vat_included` y aclara qué guarda `.vat_on_interest_bp`. f: el diccionario es idéntico al generador e idempotente. a2: S1-13 OK |
| N-02 (media) (0, NULL) sin significado y RN-26 contra la invariante 16 | **CERRADO** | j: las 5 fuentes de S1 dicen lo mismo, como PROVISIONAL con la cita de `S3-iva-lectura` pregunta 2 (ABIERTO, sin recomendación, que no se decide por Axel). d2: el CHECK acepta el par y da la misma matriz de aceptados y rechazados que en r2 (la corrección dice que solo tocó comentarios; no lo comparé contra la versión de r2, que no está en git). h: las 8 líneas citan el bloqueo abierto |
| M-01 (media, r1) CHECK de `vat_included` | sigue CERRADO | d2 y d3 |
| M-02 (media, r1) Q2-68 C sin bloqueo | sigue CERRADO | h: RN-03 l. 103 e invariante 18 l. 2360 citan `S1-reparto-edicion-gasto` (ABIERTO) |
| B-01 a B-04 (r1) y B-05 y B-06 (r2) | siguen | ver arriba; B-06 se amplía a `S3-iva-lectura` |

## Solicitudes fuera de alcance
1. **Orquestador → agente de C1 (H-01):** lo de la corrección sugerida de H-01. Sin eso no hay commit «[Puerta 2]».
2. **Integrador o agente de bloqueos (B-06):** dar de alta `S1-reparto-edicion-gasto` y `S3-iva-lectura` en INDICE y en K46, y completar «Cómo seguí» de `S3-iva-lectura`.
3. **Dueño de `s1_verifica.py` (B-07):** solicitudes 3 (a) a (c) de `correccion-S1-r2.md`.
4. **Integrador (B-04):** ratificar los cambios a `gen_dict.py` y `s1_verifica.py`.
