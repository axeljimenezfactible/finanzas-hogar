# P2c · REVISION-S1-r1 · Esquema 1.4 (migración 00004)

Revisor independiente `revisor-S1-r1` (ronda 1 de 3), flujo P2c «Aplicación de la puerta 2», 11-oct-2026, noche. No escribí nada de lo que reviso y no corregí nada.
Evalué contra `docs/criterios/S1.md` (1.1), `docs/plan/decisiones.md` (K44, K46 a K50), `docs/estado/puerta2/RESPUESTAS.md` y `PUERTA-2.md` §3 y §4.
Artefactos revisados:
- `backend/db/migrations/00004_puerta2.sql` (completa);
- `docs/contratos/diccionario-datos.md` 1.4 (secciones nuevas);
- `docs/contratos/reglas-negocio.md` 1.4 (RN-03, RN-07, RN-11, RN-12, RN-26 y RN-29);
- el reporte `docs/estado/p2c/S1-puerta2.md`;
- el diff de `backend/db/tools/s1_verifica.py`.

## Veredicto
**NO APROBADO**: 1 hallazgo alto, 2 medios y 4 bajos.
- El contenido del esquema cumple las respuestas de Axel. Valores, rangos, traspaso de `params_json`, Down, 00001 a 00003 y diccionario están verificados con comandos míos.
- El alto es la verificación global: `make verify-c1`, y con él `make verify-all`, falla por el tipo nuevo «hora HH:MM» de la 00004. Para cerrarlo hace falta un cambio de C1 que está fuera del ALCANCE de S1.
- Los dos medios son el CHECK incompleto de `credit.vat_included` y una pregunta abierta de Q2-68 C que no tiene bloqueo.

## Verificaciones (evidencia propia)
Mis scripts estuvieron en `datos/tmp/rev-s1r1/` (`chk.py` y `chk2.py`, con `python3 -I`) y ya los borré. Usé goose `/home/deepsec/finanzas-hogar-orq/gobin/goose` (v3.28.0) y sqlite3 3.53.4.

| # | Qué | Comando | Resultado | Estado |
|---|---|---|---|---|
| a1 | `make verify-s1` | `make verify-s1 GOOSE=…/gobin/goose` | `RESUMEN TODO OK`, rc=0. Pasan S1-01 a S1-17, KUESKI, K12, NEG, NEG12, K44, AVISOS13, PROG14, K47, AVISOS14, DDL, DOWN14, DOWN13, DOWN y GOOSE («versión final 4») | OK |
| a2 | `make verify-sql` | `make verify-sql GOOSE=…` | tablas=55 (STRICT=55), índices=152, disparadores=79, foreign_key_check=0 filas, integrity_check=ok. El Down inverso deja 0 objetos y el Up otra vez da 55 tablas. goose up, down-to 0 y up sin errores; rc=0 | OK |
| a3 | goose up / down-to 0 / up con mi script | `chk.py` A-GOOSE | rc=[0, 0, 0], tablas=[55, 0, 55], version=4, integrity=ok, fk_check=0 | OK |
| b1 | `@param` completo | `grep '^-- @param ' 00004 \| awk -F' \\| ' '{print NF}'` | 10 líneas, todas con 7 campos (clave, tipo, unidad, valor, decisión, sección y migración) | OK |
| b2 | Valores contra PUERTA-2 §4 y la decisión F | `chk.py` B-VALORES | 13 claves comparadas (las 10 nuevas, más `comparator.timeout_seconds` = 5, `price.stale_after_days` = 30 y `price.zero_weight_after_days` = 90 sin cambio): 0 diferencias. Q2-72: "21:30", 7 + "19:00", 1 + "08:00". Q2-73: "09:00" y 1. Q2-75: 7. Q2-76: 30. F: domingo = 7 y día 1. Claves del hogar: 20 en 1.3 y 30 en 1.4 | OK |
| b3 | Rangos (Q2-74, Q2-77, Q2-66 B, ISO 8601, HH:MM) | `chk.py` B-RANGOS | 21 negativas, 21 rechazadas: comparador 0, 31 y "5"; precio viejo 0; semana 0 y 8; día 0 y 32; −1 en anticipación, ventana y OCR; 7.5; 7.0; "24:00", "9:00", "12:60", "29:00", 2130, "21:30:00" y "ab:cd". 13 positivas, 13 aceptadas: 1 y 30; día 29 y 31; "00:00" y "23:59"; 0 en anticipación, ventana y OCR; `zero_weight` = 10 la acepta el esquema, porque es regla del motor | OK |
| b4 | Q2-77 en el motor | `sed -n '100,115p' backend/internal/precios/types.go` | `Validate`: `StaleAfterDays < 1 \|\| ZeroWeightAfterDays < StaleAfterDays` → `ErrInvalidArgument`. RN-12 l. 191: peso (D0 − edad) / (D0 − Dv + 1); con 30 y 90 da (90 − d)/61, que es K40 | OK |
| c1 | `params_json` sin los 4 parámetros | `chk.py` C-PARAMS | Claves en 1.3: {day_of_month, scopes, weekday}. En 1.4: {scopes}. `weekly_summary`, `monthly_report` y `meal_plan_missing` quedan en `{}`. `send_time` y `advance_days` no estaban en 1.3, así que el Down no tiene que reponerlos | OK |
| c2 | El Down repone el estado 1.3 | `chk.py` E-DOWN | sqlite_master, app_setting, notification (con `updated_at`) y credit quedan iguales a 1.3: 0 secciones distintas | OK |
| c3 | El Down usa el valor vigente | `chk.py` E-DOWNVAL | Con 3, 31 y 5 en las claves, el Down da `{"weekday": 3}`, `{"day_of_month": 31}` y `{"weekday": 5}` | OK |
| c4 | El Up toma el valor de 1.3 | `chk.py` C-UPVAL | 3 y 31 pasan a las claves; un 9 inválido da 7; `{"day_of_month": 31, "x": 1}` queda en `{"x":1}` | OK |
| c5 | `updated_at` conservado en el Up | `chk2.py` | notification True y credit True (fechas fijadas a mano antes de migrar) | OK |
| d1 | `vat_on_interest_bp` sin cambio | `chk.py` D-VATBP y diff de caracteres en `chk2` | Definición idéntica. El DDL de `credit` en 1.4 es el de 1.3 más 106 caracteres agregados (la columna `vat_included`) y 0 borrados | OK |
| d2 | Migración de créditos de 1.3 | `chk.py` D-BACKFILL | Con 1600 pasa a 0 y con NULL pasa a 1, que es el significado de 1.3. `updated_at` igual | OK |
| d3 | CHECK de la opción 1 | `chk.py` D-CHECK | (1, 1600) rechaza; vat_included 2 y −1 rechazan; (0, 1600), (1, NULL), (0, NULL) y (NULL, NULL) aceptan; **(NULL, 1600) acepta** | **M-01** |
| e | 00001 a 00003 intactas | `git diff --stat HEAD -- 00001… 00002… 00003…` y `sha256sum` contra `git show HEAD:` | diff vacío (rc=0). Los 3 sha256 son iguales a HEAD: c6f6979e…, 05874b55… y 7c0e9cff… | OK |
| f | Diccionario regenerado | `python3 -I backend/db/tools/gen_dict.py datos/tmp/rev-s1r1/dic.md && cmp …` | `CMP IDENTICO`; sha256 bab3528a…e1b5 en ambos. Salida: versión 1.4, 25 parámetros, 3 sin valor, 0 valores que no coinciden, «parámetros sin fila en Ajustes: []» y 0 columnas sin descripción | OK |
| g | Verificación global | `make verify-c1 GOOSE=…` | `gen2: tipo de parámetro desconocido: hora HH:MM`, `make: *** [Makefile:130: verify-c1] Error 1`, rc=2. Por eso `make verify-all` también falla | **H-01** |
| h | S1-06 y S1-14 | `grep -n -i "oct-2026\|ene-20"`, `grep -i "password\|token\|secret\|apikey"` y `"[0-9]\{8,\}:AA"` en la 00004 | Solo hay una fecha, en el comentario del encabezado (l. 3), y ningún dato. 0 secretos | OK |
| i | Nada «sin recomendación» decidido | grep de Q2-12, 48, 50 a 52, 56 a 58, 85, 87 a 91 y 94 en la 00004 y en las líneas nuevas del diccionario y de reglas | 0 citas. El reporte las deja sin decidir (§g) | OK |
| j | S1-16 al pie de la letra | `grep -rn PROVISIONAL backend/db diccionario reglas` | 5 líneas de la 00003 (l. 11, 56 a 59 y 62) citan `S1-parametros-avisos`, que ya tiene «Respuesta de Axel». Diccionario y reglas: solo la leyenda | **B-03** |
| k | RN-07, ejemplo de la serie sin intereses (Q2-38 (a)) | aritmética | 3,094.50 / 4 = 773.625 → base truncada 773.62; 773.62 × 3 = 2,320.86; 3,094.50 − 2,320.86 = 773.64. Coincide con PUERTA-2 §3.12 | OK |

Revisé las decisiones una por una:
- Q2-11 B (K47): 8 claves `reports.*` y `alerts.*`, sacadas de `params_json`, y el Down las repone.
- Q2-66 B: día de 1 a 31 en el esquema y la regla del mes corto en RN-29.
- Q2-67 opción 1: columna nueva y `vat_on_interest_bp` sin cambio. El CHECK es parcial (M-01).
- Q2-68 C: sin DDL, en RN-03 y en la invariante 18. Queda una pregunta sin bloqueo (M-02).
- Q2-38 (a): sin DDL y en RN-07.
- Q2-72 a Q2-77 (K50): valores y rangos cumplen.
- Ninguna respuesta de Axel queda contradicha.

## Hallazgos

**H-01 (alta) · Verificación global rota: `make verify-c1` y `make verify-all` fallan por el tipo nuevo «hora HH:MM».**
- Ruta: `backend/db/migrations/00004_puerta2.sql` l. 58, 60, 62 y 63 (`@param … | hora HH:MM | …`). Quien la consume es `backend/internal/api/c1gen` (`gen2_schemas.py`, `_pv`), que solo conoce «entero» y «lista de enteros».
- Evidencia: verificación g. El reporte lo reconoce («Es esperado») y lo deja como solicitud 1 a C1. Ningún reporte de P2c (`contratos-puerta2.md` tampoco) toma esa solicitud.
- Impacto: con el árbol actual, el commit «[Puerta 2]» (Q2-97, K52) no puede cerrar con `make verify-all`.
- Corrección sugerida: no la puede hacer S1 solo, porque está fuera de su ALCANCE.
  - El orquestador asigna en esta ronda P2c la solicitud 1 de `S1-puerta2.md` a un agente de C1: tipo «hora HH:MM» en `_pv` (string con patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$`), las 10 claves en sus secciones con sus rangos, `credit.vat_included` y `vat` en `Credit.missing_terms`; luego regenera.
  - Después, S1-r2 vuelve a correr `make verify-all GOOSE=…`.
  - No sirve cambiar el tipo de las horas a entero: K43 y Q2-72/Q2-73 las fijan como «HH:MM».

**M-01 (media) · El CHECK de `credit.vat_included` acepta un porcentaje de IVA con `vat_included` NULL, cosa que la opción 1 aprobada no permite.**
- Ruta: `backend/db/migrations/00004_puerta2.sql` l. 86. Contradice la l. 77 (`@col_note`) y `docs/contratos/diccionario-datos.md` l. 770.
- Evidencia: verificación d3. `INSERT … (vat_included, vat_on_interest_bp) = (NULL, 1600)` se acepta. `vat_included IN (0, 1)` da NULL con un valor NULL, así que todo el CHECK da NULL y pasa.
- Lo que dice lo aprobado:
  - la opción 1 literal (`docs/estado/bloqueos/U4-iva-credito.md`, Q2-67): «`vat_on_interest_bp` solo cuando `vat_included = 0`»;
  - la `@col_note` (l. 77): «solo lleva valor cuando `vat_included` = 0».
- Impacto: C1 y S3 tendrían que adivinar qué hacer con (NULL, 1600): ¿falta `vat`, o se suma el 16 %? Además K47 de `s1_verifica.py` no prueba ese caso.
- Corrección sugerida:
  - `CHECK (vat_included IN (0, 1) AND (vat_on_interest_bp IS NULL OR vat_included IS 0))`;
  - una negativa nueva en K47, (NULL, 1600) → rechazo;
  - repetir DOWN14 (el DROP COLUMN sigue valiendo, porque el CHECK es de la propia columna).

**M-02 (media) · Una pregunta abierta de Q2-68 C se queda en «Pendientes» del reporte, sin bloqueo y sin marca en RN-03.**
- Ruta: `docs/estado/p2c/S1-puerta2.md` §Pendientes, punto 2. Afecta a `docs/contratos/reglas-negocio.md` RN-03 y a la invariante 18 del diccionario.
- La pregunta: si después de un cambio de regla con «solo los nuevos» se edita el importe de un gasto anterior, ¿con qué partes se recalcula? `expense_share` solo guarda `amount_cents` (00001 l. 658 a 666), sin porcentaje ni versión de la regla, y `updateCostShareRule` reemplaza las partes de la regla.
- Impacto: la REGLA SUPREMA pide abrir un bloqueo para lo no decidido. Sin bloqueo, C1 y el motor de reparto (WF-3) tendrían que adivinar, y la respuesta podría pedir una columna nueva, o sea, tocar el esquema.
- Corrección sugerida:
  - abrir `docs/estado/bloqueos/S1-reparto-edicion-gasto.md` con qué falta, las opciones y el impacto (C1, motor de reparto y quizá migración). Opciones, sin decidir: recalcular con la proporción de las partes congeladas, guardar el porcentaje por parte en `expense_share` (con migración) o usar la regla vigente;
  - citarlo en RN-03 y en la invariante 18 como PROVISIONAL para ese caso.

**B-01 (baja) · Diccionario, `credit.vat_on_interest_bp` (l. 770): la celda sigue empezando con la frase de 1.3, «NULL = la tasa o la cotización ya incluye IVA (Kueski)», y luego la corrige con «Desde 1.4…».**
- En 1.4, ese NULL puede significar «incluye» (con 1), «falta capturar» (con NULL) o «falta el porcentaje» (con 0).
- Corrección sugerida: que la `@col_note` sustituya la frase, o que `gen_dict.py` la marque como «(1.3)». El comentario del DDL de la 00001 no se edita.

**B-02 (baja) · `ocr.timeout_seconds` queda en «Datos y respaldos» como «propuesta S1», sin bloqueo y sin pregunta en PUERTA-2.**
- Ruta: 00004 l. 67, RN-29 y diccionario l. 2303. `docs/estado/wf1/ocr.md` (pendientes, punto 3) la dejaba «a decidir en la puerta 2» y la puerta no la preguntó.
- Está rotulada y figura en las Propuestas para Axel, así que se puede rastrear.
- Corrección sugerida: que el integrador la lleve a la lista de excepciones de Axel o a WF-3, o abrir `U4-seccion-ocr.md` si se quiere tratar como PROVISIONAL formal.

**B-03 (baja) · S1-16 al pie de la letra: el grep de `backend/db` todavía encuentra 5 marcas PROVISIONAL en la 00003 que citan un bloqueo ya resuelto (`S1-parametros-avisos`).**
- Cómo se resolvió: con `-- @resuelve` en la 00004 (l. 53) y la excepción «historia» de `s1_verifica.py`. Es correcto, porque la 00003 no se edita (punto e).
- El problema: el texto de S1-16 en `docs/criterios/S1.md` 1.1 no prevé historia inmutable.
- Corrección sugerida: solicitud a criterios para que la Verifica de S1-16 excluya las marcas de migraciones cerradas cuyo bloqueo declara una migración posterior con `@resuelve`.

**B-04 (baja) · `gen_dict.py` y `s1_verifica.py` se editaron fuera del ALCANCE literal (lo declara el Supuesto 1).**
- Revisé el diff de `s1_verifica.py`: amplía los chequeos (PROG14, K47, AVISOS14, DOWN14, DDL sobre 00001 a 00003). NEG12 se ciñe a las 3 claves de K16 porque `optimizer.date_window_days` ya se siembra. La excepción «historia» solo cubre bloqueos que una migración posterior marca con `@resuelve`.
- Mis verificaciones propias (a3 a f) no dependen de esa herramienta y coinciden con ella.
- Corrección sugerida: el integrador ratifica la edición. Aparte, el Makefile dice «esquema 1.3» en `verify-s1` (l. 8 y 60). Es cosmético y está en la solicitud 7.

## Estado de los anteriores
Es la primera ronda de revisión de S1 en P2c y no hay hallazgos previos de este artefacto. El bloqueo de P1c `S1-parametros-avisos` (Q2-11) queda resuelto en el esquema: `@resuelve`, las 8 claves y `params_json` sin los cuatro parámetros (verificaciones b2, c1 a c4).
