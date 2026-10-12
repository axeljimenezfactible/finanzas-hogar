# P2c · correccion-S1-r2 · Esquema 1.4 (migración 00004), ronda 2

Agente `correccion-S1-r2` (Opus), flujo P2c «Aplicación de la puerta 2», 11-oct-2026, noche.
- Corrige los hallazgos alto y medios de `docs/estado/p2c/REVISION-S1-r2.md` (H-01, N-01 y N-02) sobre el artefacto de S1-puerta2 (`docs/estado/p2c/S1-puerta2.md` y `docs/estado/p2c/correccion-S1-r1.md`).
- ALCANCE: `backend/db/migrations/00004_puerta2.sql`, `docs/contratos/diccionario-datos.md` y `docs/contratos/reglas-negocio.md`. No toqué los bajos (B-01 a B-06) ni nada «sin recomendación». No hice git add ni commit.
- La 00004 solo cambia en comentarios: el SQL, el CHECK y el Down son los mismos (verificación 3).

## Hecho

### H-01 (alta) · `make verify-c1` y `make verify-all` siguen fallando (fuera de ALCANCE; no lo corregí)
- HECHO: el arreglo va en `backend/internal/api/c1gen/gen2_schemas.py` (`_pv`), que es de C1. `git status` sigue sin cambios en `backend/internal/api/c1gen/`, y ningún reporte de P2c dice que un agente de C1 tomara la solicitud 1 de `S1-puerta2.md`.
- No cambié las horas a entero. K43, Q2-72 y Q2-73 las fijan como «HH:MM». Los `@param` de las l. 62, 64, 66 y 67 de la 00004 siguen con el tipo «hora HH:MM».
- Corrí la verificación otra vez y falla con el mismo mensaje (verificación 8). Va en solicitudes fuera de alcance con dos datos nuevos para C1: lo que sale de N-02 y el efecto del comentario de `vat_included` en el generador.

### N-01 (media) · La fila 5 del mapa de Ajustes ya trae `credit.vat_included`
- `docs/contratos/diccionario-datos.md` l. 2331, fila 5 «Créditos y condiciones (incluye promociones)», columna Destino:
  - antes: `… comisiones *_fee_cents, .vat_on_interest_bp, .rate_base, …`;
  - ahora: `` `.vat_included` (00004, Q2-67 opción 1: el sí o no de «¿la tasa incluye IVA?»; NULL = falta capturar), `.vat_on_interest_bp` (el porcentaje de IVA que se suma, solo con `.vat_included` = 0), `.rate_base` ``.
- HECHO: el mapa es texto fijo de la cola del diccionario, y `gen_dict.py` lo conserva. Lo edité ahí y regeneré, y el resultado es idempotente (verificación 4). El sufijo de parámetros K44 de la fila 5 no cambia.
- No apliqué el «de preferencia» de la revisión, que pide que S1-13 de `s1_verifica.py` compare por columna: ese archivo está fuera de mi ALCANCE (B-04). Hice el chequeo con un script propio fuera del repositorio (verificación 5). Con el diccionario de antes da `faltan = [credit.vat_included]` y rc=1; ahora da `[]`. La propuesta del cambio a `s1_verifica.py` va en solicitudes.

### N-02 (media) · Qué significa (`vat_included`, `vat_on_interest_bp`) = (0, NULL)
**HECHO:**
- El CHECK de la 00004 acepta (0, NULL). Lo comprobé con mi propia matriz (verificación 6).
- U4 P-26.05 captura el sí o no y el porcentaje por separado, así que el par ocurre.
- Ese caso ya es una pregunta abierta para Axel en el bloqueo **`S3-iva-lectura`, pregunta 2**. Lo abrió `correccion-contratos-r1` en esta misma ronda P2c: ABIERTO y **sin recomendación**.
  - Opción A: «falta capturar», y C1 lista `vat` en `missing_terms`. El bloqueo dice que eso «cambia RN-26 1.4».
  - Opción B: la API no deja guardar 0 sin porcentaje, con un CHECK o con una validación, y S3 responde `ErrInvalid`.
- `S3-interfaz.md` 1.5 §5.3 ya escribe `ErrMissingCapture` para ese par como lectura PROVISIONAL con ese bloqueo. RN-26 seguía diciendo `vat` «mientras `vat_included` sea NULL». De ahí la contradicción que señala la revisión.

**Por qué no lo escribí FIRME:** la revisión propone escribirlo como consecuencia de K18. Pero la pregunta exacta está abierta y sin recomendación, y K46 prohíbe decidir esas preguntas por Axel. Por eso lo marqué **PROVISIONAL, con la cita de `S3-iva-lectura` pregunta 2**, que es la alternativa que la propia revisión ofrece.

**Lectura en uso** (Supuesto 1): la de S3 1.5 y la de un esquema que admite el par.
- Falta el porcentaje, `Credit.missing_terms` lista `vat` y S3 responde «falta capturar» (`ErrMissingCapture`).
- No calcula con 0 % ni inventa un porcentaje (K18). En eso coinciden las dos opciones del bloqueo.
- Se dice qué pasa si Axel elige B. No cambié el CHECK.

Dónde quedó escrito (las mismas cuatro ideas en cada lugar):

| Lugar | Cambio |
|---|---|
| `docs/contratos/reglas-negocio.md` RN-26 (l. 278 y 279) | La lista de `missing_terms` remite al caso. Párrafo nuevo **«Caso `vat_included` = 0 sin porcentaje · PROVISIONAL» (bloqueo abierto `S3-iva-lectura`, pregunta 2, sin recomendación)** con HECHO, la lectura en uso, lo que cambia con B y «ninguna cifra canon cambia». RN-26 sigue FIRME en el índice, como RN-03 con su sub-viñeta PROVISIONAL |
| `reglas-negocio.md` «Cambios en 1.4» (l. 11) | Una frase de la corrección, sin la palabra PROVISIONAL, porque esa línea cita el bloqueo resuelto `S1-parametros-avisos` y el chequeo K44 la marcaría |
| 00004, comentario de la columna (l. 99) → diccionario l. 782 | Descripción: «Con 0 y vat_on_interest_bp NULL falta el porcentaje: PROVISIONAL (bloqueo abierto S3-iva-lectura, pregunta 2) `vat` también va en missing_terms y S3 dice «falta capturar»; no calcula con 0 % ni inventa un porcentaje (K18)». Restricción: «con 0 el CHECK admite el porcentaje en NULL» |
| 00004 `@used_new` (l. 93) → diccionario l. 751 | «con 0 y sin porcentaje también lo lista, PROVISIONAL por el bloqueo abierto `S3-iva-lectura`, pregunta 2, RN-26» |
| 00004 encabezado (l. 3 y 4, 21 a 23, 43 y 44) | Autoría de la corrección; en el apartado Q2-67, «acepta vat_included = 0 sin porcentaje … falta el porcentaje y es PROVISIONAL (bloqueo abierto S3-iva-lectura…)»; y en el apartado del Down, qué es ese par en 1.4 |
| Diccionario, invariante 16 (l. 2358, texto fijo de la cola) | HECHO del esquema, el **PROVISIONAL** con la lectura en uso y lo que cambia con B |

- Cada línea con PROVISIONAL cita `S3-iva-lectura` en la misma línea, junto con «bloqueo abierto». Así la cuenta S1-16 aunque su patrón no reconozca nombres que empiezan con `S3-` (ver Pendientes).
- **«Sumar el caso a K47»:** K47 vive en `backend/db/tools/s1_verifica.py`, fuera de mi ALCANCE (B-04). No lo edité. Probé el caso con mi script (verificación 6) y la línea exacta para K47 va en solicitudes.
- **«Avisar a C1 y S3»:** va en solicitudes fuera de alcance.

## Verificación
Mis scripts están en `/tmp/cs1r2/scripts/` (`fix.py`, `chk13.py`, `chk_ddl.py` y `chk_n02.py`), todos con `python3 -I`. Los respaldos de los tres archivos antes del cambio están en `/tmp/cs1r2/antes/`. Usé goose `/home/deepsec/finanzas-hogar-orq/gobin/goose` (v3.28.0) y sqlite3 3.53.4.

| # | Qué | Comando | Salida | Estado |
|---|---|---|---|---|
| 1 | Reemplazos exactos (cada texto aparece una vez) | `python3 -I /tmp/cs1r2/scripts/fix.py /home/deepsec/finanzas-hogar` | `00004_puerta2.sql OK 5 cambios` · `diccionario-datos.md OK 2 cambios` · `reglas-negocio.md OK 3 cambios`. En el primer intento a dos líneas nuevas del encabezado les faltaba `--`, y `gen_dict.py` falló con `near "(": syntax error`. Las corregí y no quedan líneas sin `--` fuera del SQL | OK |
| 2 | Qué líneas cambiaron | `diff` contra `/tmp/cs1r2/antes/` | 00004: `3c3,4`, `20c21,23`, `40c43,44`, `89c93`, `95c99` (solo comentarios). Reglas: `11c11`, `278c278,279`. Diccionario: `751c751`, `782c782`, `2331c2331`, `2358c2358`. Tamaños: 22,815 → 23,696 B, 68,536 → 69,687 B y 196,676 → 197,887 B | OK |
| 3 | La 00004 solo cambia en comentarios | `chk_ddl.py` (00001 a 00003 + 00004 de antes y de ahora, en memoria) | `Up objetos 286 286 iguales True` · `Down objetos 282 282 iguales True` · `código SQL sin comentarios idéntico: True` | OK |
| 4 | Diccionario regenerado e idempotente | `python3 -I backend/db/tools/gen_dict.py`, otra vez a `/tmp/cs1r2/dic2.md` y `cmp` | «parámetros sin fila en «Ajustes: 14 secciones»: []» · «versión 1.4 · parámetros 25 · sin valor 3 · valores iniciales que no coinciden con lo sembrado: []» · «tablas 55 … índices 152 disparadores 79» · «columnas sin descripción: []» · rc=0 · `CMP_IDENTICO` | OK |
| 5 | N-01: cada columna de un `@alter` con `@used`/`@used_new` «Ajustes «X»» está en la fila X | `python3 -I /tmp/cs1r2/scripts/chk13.py /home/deepsec/finanzas-hogar` | Antes: `faltan = [('00004_puerta2.sql', 'credit.vat_included', 'Créditos y condiciones')]`, rc=1. Ahora: 8 columnas revisadas (4 de crédito de la 00002, 3 de salud de la 00002 y la 00003, y `vat_included`), `faltan = []`, rc=0. Informativo, no es de N-01: 5 tablas de la 00001 cuyo `@used` nombra una sección que no las cita por nombre (`cost_share_rule` en Categorías, que sí trae `category.default_cost_share_rule_id`; `product_group`, `product` y `presentation`) | OK |
| 6 | N-02: matriz del CHECK, marcas y textos | `python3 -I /tmp/cs1r2/scripts/chk_n02.py /home/deepsec/finanzas-hogar` | (1) tarjeta `bank_card`, 11 casos con UPDATE y rollback, `aceptado (0, NULL)=True`, `discrepancias=[]`: acepta (0, NULL), (0, 1600), (0, 0), (1, NULL) y (NULL, NULL); rechaza (1, 1600), (NULL, 1600), (1, 0), (NULL, 0), (2, NULL) y (0, 10001). (2) `S3-iva-lectura ABIERTO sin respuesta=True`; 8 líneas PROVISIONAL del par (00004 l. 23, 44, 93 y 99; diccionario 751, 782 y 2358; reglas 279), todas con la cita. (3) invariante 16, nota de `vat_included` y RN-26 dicen lo mismo = True, True y True. rc=0 | OK |
| 7 | S1 | `make verify-s1 GOOSE=…` | rc=0, `RESUMEN TODO OK`. S1-13 OK. S1-16: «líneas con PROVISIONAL: 12 (migración 4, diccionario 5, reglas 3) · … líneas sin bloqueo en su bloque=0 []». K44: «marcas PROVISIONAL que citan un bloqueo resuelto=[]». PROG14: «diccionario K44 con las 10 y sin PROVISIONAL=True · RN-29 … sin PROVISIONAL … =True». K47: «negativas: 19, rechazadas 19 … positivas: 11, aceptadas 11». DDL, DOWN14 («sqlite_master igual al esquema 1.3 = True (244 objetos, diferencias=0)») y AVISOS14 OK | OK |
| 8 | Verificación global | `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | rc=2. SQL: 4 Up, `tablas=55 (STRICT=55) · índices=152 · disparadores=79 · foreign_key_check=0 filas · integrity_check=ok`, el Down deja 0 objetos y goose hace up, down-to 0 y up sin errores. WF-0 OK. S1 `RESUMEN TODO OK`. U1 `RESULTADO: TODO OK`. Criterios: «INDICE contra conteo propio: 0 diferencias». WF-1 «verificación completa OK». **C1: `gen2: tipo de parámetro desconocido: hora HH:MM` / `make: *** [Makefile:130: verify-c1] Error 1`**. Es el único destino que falla | FALLA solo C1 (H-01, fuera de ALCANCE) |
| 9 | 00001 a 00003 intactas; c1gen y herramientas sin cambios míos | `git diff --stat HEAD -- 00001… 00002… 00003…` · `git status --short backend/internal/api/c1gen/` · `git diff --stat -- backend/db/tools/` | Vacío. Vacío. `gen_dict.py` +44 y `s1_verifica.py` +110: los mismos números que vio la revisión r2, porque no los toqué | OK |

## Supuestos
1. **La lectura en uso de (0, NULL) es la opción A de `S3-iva-lectura` pregunta 2, marcada PROVISIONAL. No es una decisión.**
   - HECHOS que la sostienen: el esquema en uso acepta el par, y la opción B pide cambiarlo. `S3-interfaz.md` 1.5 ya usa `ErrMissingCapture`. C1 v1.3 todavía no existe, así que nada de C1 está en uso.
   - Con RN-26 sin tocar, C1 diría «sin datos faltantes» y S3 «falta capturar» para la misma tarjeta, que es la contradicción de N-02.
   - No recomiendo A ni B a Axel. El bloqueo sigue ABIERTO y sin recomendación.
2. Las dos opciones del bloqueo coinciden en que S3 no calcula con 0 % ni inventa un porcentaje. Esa parte la escribí como consecuencia de K18 («no inventa tasa ni mínimo»), igual que lo hace `S3-interfaz.md`. Va dentro del párrafo PROVISIONAL y no como regla FIRME aparte.
3. No toqué `backend/db/tools/s1_verifica.py` ni `gen_dict.py`, aunque la revisión sugiere sumar el caso a K47 y endurecer S1-13. Están fuera de mi ALCANCE literal, y B-04 ya señala que se editaron fuera de él. Probé lo mismo con mis scripts (verificaciones 5 y 6) y dejo los cambios exactos en solicitudes.
4. La celda de `credit.vat_on_interest_bp` (diccionario l. 770) sigue empezando con la frase de 1.3, «NULL = la tasa o la cotización ya incluye IVA». Es B-01 (baja) y no estaba en mi encargo. No la toqué. El resto de esa celda remite a `vat_included` y a RN-26.

## Bloqueos
- No abrí ninguno nuevo. Cité el que ya existe: **`docs/estado/bloqueos/S3-iva-lectura.md`, pregunta 2** (ABIERTO, sin recomendación, P2c).

## Pendientes
- H-01: `make verify-c1`, y con él `make verify-all`, pasan cuando un agente de C1 tome la solicitud 1 de `S1-puerta2.md` y cierre FALLA K44 y C1-15 de `verify_c1.py`. Después, la ronda 3 corre `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`.
- El patrón de bloqueos de S1-16 y K44 en `s1_verifica.py` (`(?:S1|C1|U1|U2|S6|criterios)-…`) no reconoce nombres que empiezan con `S3-` ni con `U4-`.
  - Hoy S1-16 da OK porque mis líneas dicen «bloqueo abierto» después de PROVISIONAL, pero informa `bloqueos abiertos citados: ['S1-reparto-edicion-gasto']`, sin `S3-iva-lectura`.
  - El cambio va en la solicitud 3.
- Bajos sin cambios: B-01 a B-04 de r1 y B-05 y B-06 de r2.

## Solicitudes fuera de alcance
1. **Orquestador → agente de C1 (H-01 y N-02):**
   - Tomar la solicitud 1 de `docs/estado/p2c/S1-puerta2.md`: tipo «hora HH:MM» en `backend/internal/api/c1gen/gen2_schemas.py` `_pv`, como string con patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$`, las 10 claves con sus rangos, `credit.vat_included` y `vat` en `Credit.missing_terms`. Regenerar y cerrar FALLA K44 y C1-15 de `verify_c1.py`. Las horas no pasan a entero (K43, Q2-72, Q2-73).
   - Dato nuevo de N-02: mientras `S3-iva-lectura` pregunta 2 siga abierta, `Credit.missing_terms` lista `vat` también con `vat_included` = 0 y `vat_on_interest_bp` NULL. Es PROVISIONAL (RN-26 1.4) y la API no responde 422 por ese par. Con (1, %) y (NULL, %) el CHECK sí rechaza (r1), y C1 puede responder 422 antes de llegar a la base.
   - Efecto en el generador: el comentario de `credit.vat_included` en la 00004 ahora contiene «PROVISIONAL … S3-iva-lectura».
     - `gen1_base.py` (l. 85 a 90) marcará `x-provisional: true` en el campo y citará `docs/estado/bloqueos/S3-iva-lectura.md`.
     - La marca se refiere solo al par (0, NULL). Si C1 prefiere no marcar el campo entero, la nota de la columna en la 00004 se puede reescribir sin la palabra PROVISIONAL en una ronda de S1. Pedirlo.
2. **Orquestador → S3 (contratos y motor):**
   - RN-26 1.4, la invariante 16 y la nota de `vat_included` ya dicen lo mismo que `S3-interfaz.md` 1.5 §5.3 para (false, nil): `ErrMissingCapture`, PROVISIONAL `S3-iva-lectura` pregunta 2.
   - Revisar que la fixture `s3-tarjeta-minimo` y §14.1 nombren ese caso si hace falta.
   - Sigue vigente la solicitud 2 de `REVISION-S1-r2.md`, que pide quitar las marcas a `U4-iva-credito` (resuelto) de `S3-interfaz.md` l. 116, 257, 272, 884, 903 y 919, si todavía quedan.
3. **Integrador o dueño de `backend/db/tools/s1_verifica.py`:**
   - (a) En K47 (l. 653 a 655), sumar la positiva PROVISIONAL `('iva_no_incluido_sin_porcentaje', "UPDATE credit SET vat_included = 0, vat_on_interest_bp = NULL WHERE id = (SELECT min(id) FROM credit)")` con un comentario que cite `S3-iva-lectura` pregunta 2. Si Axel elige B, pasa a negativa.
   - (b) En S1-13, comprobar que cada columna de un `@alter` cuyo `@used`/`@used_new` nombra Ajustes «X» aparezca en la fila X. La lógica está en `/tmp/cs1r2/scripts/chk13.py`.
   - (c) En S1-16 y K44 (`PATB`, l. 59, y el patrón de S1-16, l. 261), aceptar `S2`–`S8` y `U3`/`U4` como prefijos de bloqueo.
4. **Agente de bloqueos:**
   - En `docs/estado/bloqueos/S3-iva-lectura.md` «Cómo seguí», anotar que S1 1.4 (RN-26, nota de `vat_included` e invariante 16) escribe la lectura de la pregunta 2 como PROVISIONAL con la cita.
   - Pendiente de r2: B-06, que es dar de alta `S1-reparto-edicion-gasto` en `INDICE.md`.

## Propuestas para Axel
- PREGUNTA (ya está en `S3-iva-lectura`, pregunta 2, sin recomendación): si marcas «no incluye IVA» en una tarjeta o crédito de tienda y todavía no escribes el porcentaje, ¿qué hace la app?
  - A) La guarda y la tarjeta dice «falta capturar» hasta que pongas el porcentaje.
  - B) No te deja guardar «no» sin el porcentaje.
  - Mientras contestas, la app usa A como PROVISIONAL. Ninguna cifra cambia.

## Consumo aproximado
- Acciones de herramienta: ~36.
- KB leídos: ~115.
  - Revisión r2: 14.
  - Reporte r1: 12.
  - Encabezado y Up de la 00004: ~13.
  - Diccionario por secciones: ~25.
  - RN-26 y encabezado de reglas: ~4.
  - Bloqueos `U4-iva-credito` y `S3-iva-lectura`: ~6.
  - Fragmentos de S3, gen_dict, s1_verifica y c1gen: ~15.
  - Salidas de verificación: ~26.
- KB escritos de mi autoría: ~31.
  - Este reporte: ~13.
  - Cambios netos en los tres archivos: ~3.3.
  - Scripts temporales en `/tmp/cs1r2/scripts`: ~10.
  - Respaldos y registros en `/tmp`: no cuentan para el repositorio.
  - Además, el generador reescribe el diccionario completo (198 KB), con un cambio neto de ~1.2 KB.
