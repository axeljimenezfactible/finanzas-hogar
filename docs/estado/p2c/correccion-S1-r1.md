# P2c · correccion-S1-r1 · Esquema 1.4 (migración 00004)

Agente `correccion-S1-r1` (Opus), flujo P2c «Aplicación de la puerta 2», 11-oct-2026, noche.
- Corrige los hallazgos alto y medios de `docs/estado/p2c/REVISION-S1-r1.md` sobre el artefacto de S1-puerta2 (`docs/estado/p2c/S1-puerta2.md`).
- No toqué los bajos (B-01 a B-04) ni nada «sin recomendación». No hice git add ni commit.

## Hecho

### H-01 (alta) · `make verify-c1` y `make verify-all` siguen fallando (fuera de ALCANCE, no lo corregí)
- HECHO: el arreglo está en `backend/internal/api/c1gen/gen2_schemas.py` (`_pv`), que es de C1. La corrección sugerida por la revisión lo deja fuera del ALCANCE de S1: el orquestador asigna a un agente de C1 la solicitud 1 de `S1-puerta2.md`.
- No cambié las horas a entero, porque K43 y Q2-72/Q2-73 las fijan como «HH:MM».
- Lo verifiqué otra vez: falla con el mismo mensaje (ver Verificación). Va en solicitudes fuera de alcance, con un dato nuevo para C1: la regla del CHECK corregido (M-01).

### M-01 (media) · CHECK de `credit.vat_included`: corregido
`backend/db/migrations/00004_puerta2.sql`:
- CHECK nuevo, tal como lo sugirió la revisión: `CHECK (vat_included IN (0, 1) AND (vat_on_interest_bp IS NULL OR vat_included IS 0))`. Ahora (NULL, 1600) se rechaza, igual que (1, 1600).
- HECHO que obligó a ajustar también el traspaso: SQLite 3.37 o posterior prueba el CHECK de una columna nueva contra las filas que ya existen. Un crédito de 1.3 con porcentaje queda en (NULL, 1600) durante el `ALTER`, así que con el CHECK nuevo la migración fallaba. Prueba: con sqlite3 3.53.4, `ALTER TABLE c ADD COLUMN vat_included … CHECK (…)` sobre una fila (1600) da `CHECK constraint failed`, rc=1.
- Cómo quedó: antes del `ALTER`, el porcentaje se aparta en la tabla temporal `p2_credit_vat_13` (`CREATE TEMP TABLE … AS SELECT id, vat_on_interest_bp`) y se deja en NULL. Después del `ALTER` vuelve, junto con `vat_included` = 0, y la tabla temporal se borra.
- Lo que no cambia:
  - el significado de 1.3: con porcentaje pasa a 0 y sin porcentaje a 1, el Supuesto 3 de S1-puerta2;
  - `updated_at`: se conserva con el mismo truco de dos pasos (`~`) del implementador;
  - el Down: el `DROP COLUMN` sigue valiendo.
- Comentarios actualizados: encabezado (líneas 3, 18 a 20 y 33 a 40), el comentario del paso y la regla de la columna («con 1 o NULL, `vat_on_interest_bp` es NULL (solo lleva valor con 0)»).
- `@col_note` (línea 80) y diccionario: ya decían «solo lleva valor cuando `vat_included` = 0». Ahora el esquema coincide con ellos.
- `backend/db/tools/s1_verifica.py` (ver Supuesto 1):
  - K47 suma 2 negativas: `iva_sin_dato_con_porcentaje` (INSERT (NULL, 1600)) e `iva_sin_dato_con_porcentaje_al_editar` (UPDATE a (NULL, 1600)). Ahora son 19 negativas, antes 17;
  - la positiva de 1.2 `tarjeta_con_datos_k18` lleva `vat_included` = 0 cuando la columna existe (`VI18`). Sin eso, el CHECK corregido la rechazaba, porque insertaba (NULL, 1600).
- DOWN14 corrido otra vez: OK.
- `docs/contratos/diccionario-datos.md` regenerado con `gen_dict.py`. La fila de `credit.vat_included` (l. 782) dice «con 1 o NULL, vat_on_interest_bp es NULL (solo lleva valor con 0)».

### M-02 (media) · Pregunta de Q2-68 C sin bloqueo: bloqueo abierto y marcas PROVISIONAL
- Nuevo `docs/estado/bloqueos/S1-reparto-edicion-gasto.md` (ABIERTO), con HECHO, la pregunta, las opciones sin decidir, el impacto y cómo seguí.
- HECHOS que cita:
  - `expense_share` guarda solo `amount_cents`;
  - `updateCostShareRule` reemplaza las partes;
  - `updateExpense` de C1 v1.2 dice «recalcula las partes», o sea, con la regla vigente (`docs/contratos/extractos/gastos.yaml`);
  - S2 no reparte gastos registrados (`S2-interfaz.md`, l. 17).
- PREGUNTA: después de un cambio con «solo los nuevos», si se edita un gasto anterior, ¿con qué partes se recalcula?
- Opciones, sin decidir:
  - A: proporción de las partes congeladas. Escala las partes fijas como si fueran porcentaje y no sabe quién es la persona del resto.
  - B: guardar cómo se calculó cada parte. Pide una migración posterior a la 00004, C1 y el reparto de la API.
  - C: la regla vigente, que es la opción en uso. El gasto editado pasa al reparto nuevo.
  - En A y B también falta decidir cómo se sabe que un gasto quedó fuera del cambio.
- `docs/contratos/reglas-negocio.md` RN-03: sub-viñeta **PROVISIONAL** que cita `S1-reparto-edicion-gasto`, con la opción en uso (`updateExpense` de C1 v1.2 recalcula con la regla vigente) y el aviso de que la respuesta podría pedir una columna.
- `docs/contratos/diccionario-datos.md`, invariante 18: la misma marca **PROVISIONAL** con la cita. Está en el texto fijo de la cola, que `gen_dict.py` conserva (ver la verificación de idempotencia).

## Verificación
Mis scripts están en `/tmp/cs1r1/scripts/`, fuera del repositorio: `fix_mig.py`, `fix_verifica.py`, `fix_docs.py` y `chk.py`, todos con `python3 -I`. Respaldos de los archivos antes del cambio: `/tmp/cs1r1/*.antes`. Usé goose `/home/deepsec/finanzas-hogar-orq/gobin/goose` (v3.28.0) y sqlite3 3.53.4.

| # | Qué | Comando | Salida | Estado |
|---|---|---|---|---|
| 1 | El CHECK sugerido tal cual, sobre filas de 1.3 | `sqlite3 t.db "ALTER TABLE c ADD COLUMN vat_included INTEGER CHECK (…)"` con una fila (1600) | `Error … CHECK constraint failed`, rc=1 | motivo del apartado temporal |
| 2 | Up sobre 1.3 con datos | `chk.py` A-UP13 | (1600, 1600, 0, True), (None, None, 1, True) y (1234, 1234, 0, True): el porcentaje vuelve, `vat_included` = 0/1 y `updated_at` sin cambio. `temp_restantes= 0` | OK |
| 3 | Matriz del CHECK (INSERT y UPDATE) | `chk.py` B-CHECK | Rechaza (NULL, 1600), (1, 1600), (2, NULL) y (−1, NULL). Acepta (0, 1600), (1, NULL), (0, NULL) y (NULL, NULL). Con UPDATE da lo mismo: True | OK |
| 4 | Down con datos | `chk.py` C-DOWN | Las columnas de 1.3 quedan iguales: True. Sin `vat_included` | OK |
| 5 | goose con datos (tabla temporal dentro de la transacción) | `chk.py` D-GOOSE: up-to 3, créditos, up, down-to 3, up | rc [0, 0, 0, 0]. Tras up: (1600, 0, '2026-01-01T00:00:00.000Z') y (None, 1, '2026-01-02T00:00:00.000Z'). Versión 4. Tras down-to 3: (1600, …) y (None, …). integrity ok, fk 0 | OK |
| 6 | `make verify-s1 GOOSE=…` | rc=0 | `RESUMEN TODO OK` | OK |
| 7 | `s1_verifica.py` directo (líneas completas) | `GOOSE=… python3 -I backend/db/tools/s1_verifica.py` | K47: «negativas: 19, rechazadas 19, aceptadas por error [] · positivas: 11, aceptadas 11». NEG12: «positivas: 12, aceptadas 12, rechazadas por error []». AVISOS14: «vat_included de los créditos de 1.3 (con 1600 y sin porcentaje) = [0, 1] … updated_at conservado=True». DOWN14: «sqlite_master igual al esquema 1.3 = True (244 objetos, diferencias=0) … créditos con sus columnas de 1.3 iguales = True». S1-16: «líneas con PROVISIONAL: 4 (migración 0, diccionario 2, reglas 2) · bloqueos abiertos citados: ['S1-reparto-edicion-gasto'] · líneas sin bloqueo=0». K44: «marcas PROVISIONAL que citan un bloqueo resuelto=[]». PROG14 OK. DDL OK. `RESUMEN TODO OK`, rc=0 | OK |
| 8 | Diccionario regenerado e idempotente | `python3 -I backend/db/tools/gen_dict.py`, otra vez a `/tmp/cs1r1/dic2.md` y `cmp` | versión 1.4, 25 parámetros, 3 sin valor, 0 valores que no coinciden, 55 tablas, 152 índices, 79 disparadores, 0 columnas sin descripción. `CMP_IDENTICO` (la invariante 18 editada sobrevive). diff contra el anterior: solo l. 782 (`vat_included`) y l. 2360 (invariante 18) | OK |
| 9 | 00001 a 00003 intactas | `git diff --stat HEAD -- 00001… 00002… 00003…` | vacío, rc=0 | OK |
| 10 | Verificación global | `make -k verify-all GOOSE=…` | rc=2. SQL: 55 tablas (STRICT=55), 152 índices, 79 disparadores, fk 0, integrity ok; Down a 0 objetos; goose up, down-to 0 y up sin errores. WF-0 OK. S1 `RESUMEN TODO OK`. U1 `RESULTADO: TODO OK`. Criterios: «INDICE contra conteo propio: 0 diferencias». WF-1 «verificación completa OK». **C1: `gen2: tipo de parámetro desconocido: hora HH:MM` / `make: *** [Makefile:130: verify-c1] Error 1`** | FALLA solo C1 (H-01, fuera de ALCANCE) |
| 11 | Nadie más escribe créditos (NULL, porcentaje) en SQLite | `grep -rln vat_on_interest_bp` (py, go y sql) y `grep` en `verify_c1.py` | Los motores y oráculos de S3 lo usan en memoria (JSON) y no en SQLite. `verify_c1.py` no inserta créditos | OK |

## Supuestos
1. **`backend/db/tools/s1_verifica.py` está fuera de mi ALCANCE literal.** Lo edité porque la corrección de M-01 pide expresamente «agregar en K47 la negativa (NULL, 1600)» y porque, sin ajustar la positiva `tarjeta_con_datos_k18`, `make verify-s1` fallaba con el CHECK corregido. Fueron 3 cambios mínimos, marcados con «correccion-S1-r1, M-01». El integrador los ratifica, como en B-04.
2. **`docs/estado/bloqueos/S1-reparto-edicion-gasto.md` también está fuera de mi ALCANCE literal.** Lo creé porque la corrección de M-02 lo pide con ese nombre. No edité `docs/estado/bloqueos/INDICE.md` (solicitud 2).
3. **El apartado temporal del porcentaje es técnico y no es una decisión.** Conserva lo que el implementador ya había fijado para los créditos de 1.3 (Supuesto 3 de S1-puerta2, propuesta para Axel). La tabla temporal no es STRICT ni entra en el diccionario: `gen_dict.py` solo lee `CREATE TABLE <nombre> (`, y la tabla se borra en la misma migración.
4. **La «opción en uso» que marqué PROVISIONAL en RN-03 y en la invariante 18 es la C.** Es un HECHO de C1 v1.2 (`updateExpense` recalcula las partes) y no una decisión mía. El bloqueo no recomienda ninguna opción.

## Bloqueos
- Nuevo: `docs/estado/bloqueos/S1-reparto-edicion-gasto.md` (ABIERTO). Se decide en WF-3 o cuando Axel conteste.

## Pendientes
- H-01: `make verify-c1`, y con él `make verify-all`, pasan cuando C1 tome la solicitud 1 de `S1-puerta2.md`. Después, S1-r2 corre otra vez `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`.
- Los bajos B-01 a B-04 de `REVISION-S1-r1.md` siguen como estaban, porque no estaban en mi encargo. B-01: la celda de `credit.vat_on_interest_bp` sigue empezando con la frase de 1.3.

## Solicitudes fuera de alcance
1. **Orquestador → agente de C1 (H-01):**
   - tomar la solicitud 1 de `docs/estado/p2c/S1-puerta2.md`: tipo «hora HH:MM» en `gen2_schemas.py` `_pv`, como string con patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$`, más las 10 claves con sus rangos, `credit.vat_included` y `vat` en `Credit.missing_terms`, y regenerar;
   - dato nuevo de esta ronda: el esquema rechaza `vat_on_interest_bp` con valor salvo con `vat_included` = 0, o sea (1, %) y (NULL, %). La validación de `Credit` en C1 puede reflejarlo para responder 422 y no un error de la base.
2. **Integrador o agente de bloqueos:** agregar `S1-reparto-edicion-gasto` a `docs/estado/bloqueos/INDICE.md`, en «Bloqueos nuevos de P2c». La carpeta tiene ahora 76 archivos sin contar el índice; el último conteo del índice daba 75.
3. **Integrador:** ratificar los cambios a `backend/db/tools/s1_verifica.py` (Supuesto 1).
4. **C1, U4 (P-02, edición en el lugar) y el reparto de la API (WF-3):** leer la marca PROVISIONAL de RN-03 y de la invariante 18 mientras se resuelve `S1-reparto-edicion-gasto`.

## Propuestas para Axel
- PREGUNTA, en el bloqueo `S1-reparto-edicion-gasto`: después de cambiar una regla de reparto con «solo los nuevos», si se edita un gasto anterior, ¿con qué partes se recalcula? Las opciones son A (proporción congelada), B (guardar cómo se calculó cada parte, con migración) o C (regla vigente, lo de hoy). No recomiendo ninguna.

## Consumo aproximado
- Acciones de herramienta: ~42.
- KB leídos: ~120 (revisión 13, reporte S1 16, migración 22, fragmentos de herramientas, DDL, RN-03, diccionario y extractos ~35, y salidas de verificación ~30).
- KB escritos de mi autoría: ~27 (este reporte ~12, bloqueo 4.5, cambios en migración ~2, s1_verifica ~1, RN-03 e invariante 18 ~1, y scripts temporales ~7). Además, el generador reescribe el diccionario completo (196 KB) con un cambio neto de ~1 KB.
