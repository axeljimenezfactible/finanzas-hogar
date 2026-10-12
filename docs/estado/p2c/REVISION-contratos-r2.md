# P2c · REVISION-contratos-r2 · 11-oct-2026, noche

Revisor independiente, ronda 2 de 3 (`revisor-contratos-r2`). No escribí nada de lo revisado y no corregí nada.
Objeto: lo que escribieron `contratos-puerta2` y `correccion-contratos-r1`: `S2-interfaz.md` 1.3, `S3-interfaz.md` 1.5, `S4-interfaz.md` 1.5 y `S5-interfaz.md` 1.5; `S2-ventanas.md` pregunta 4; los bloqueos nuevos `S3-formula-tarjeta-pdf`, `S3-iva-lectura` y `S4-empates-sin-confirmar`; los reportes `contratos-puerta2.md` y `correccion-contratos-r1.md`.
Contra qué: `docs/criterios/S2.md` a `S5.md`; PUERTA-2 §2 y §3; decisiones §B a §E y §K (K46); `RESPUESTAS.md`.
Cada verificación la volví a correr con mis scripts en `datos/tmp/revisor-contratos-r2/`, que ya borré.

## Veredicto
**NO APROBADO.** Quedan 1 hallazgo alto, 2 medios y 4 bajos.
- El alto es la verificación global, que sigue rota (`verify-c1`). Está fuera del alcance de los contratos.
- Los dos medios piden poco: dar de alta 3 bloqueos en INDICE y K46, y agregar 7 comentarios de código a la lista para WF-3.
- Ninguna cifra canon cambió y ninguna respuesta «sin recomendación» quedó decidida.

## Verificaciones (evidencia propia)
| # | Verificación | Resultado | Evidencia |
|---|---|---|---|
| a | Q2-01 a Q2-08 aplicadas según su corrección sugerida | CUMPLE | Q2-01: S2 l. 151 («contiguas, ascendentes y disjuntas»), `S2-ventanas.md` l. 90 (pregunta 4 con 1A–1D y 2A–2C, y sus cifras 50000, 97375 y 292125); la consulta A + C está en la lista de cambios. Q2-02: S5 l. 10 y l. 126 (regla de 3 canastas o más). Q2-03: S5 l. 131 (paso 8, PROVISIONAL, `S5-carritos-granel` ABIERTO). Q2-04: S5 l. 189 («1,021.00» = 102100) y l. 194–195. Q2-05 y Q2-08: S5 l. 166–167 y l. 291–292 (completas y luego incompletas). Q2-06: S4 l. 10 y l. 254–255 (`choice_store`, `choice_presentation`). Q2-07: S4 l. 11 y l. 185 (regla del motor PROVISIONAL, `S4-presentacion-respaldo` ABIERTO) |
| b | Ninguna «sin recomendación» decidida | CUMPLE | Q2-01 caso 2 → 2A en uso, PROVISIONAL (`S2-ventanas` ABIERTO). Variante de Q2-07 → regla del motor (`S4-presentacion-respaldo` ABIERTO). Q2-56 a Q2-58 → A en uso (S5 l. 194, `S5-profeco-filtro` ABIERTO). Q2-90 → valores `propuesto`, PROVISIONALES (S4 l. 293, `S4-empates-sin-confirmar` ABIERTO; `tabla.csv`: 150 de 150 renglones `propuesto`, contados por mí). Ver B-3 |
| c | Cada PROVISIONAL cita un bloqueo abierto | CUMPLE en los contratos; el registro está incompleto (M-1) | `prov.py` revisa cada renglón con «provisional» de S2 a S5 contra la línea «Estado del bloqueo (puerta 2)» o «Estado:» del archivo. **0 marcas activas** citan un bloqueo no ABIERTO. Las que aparecen son bitácora o negaciones: S3 l. 7, 258, 273, 315, 790, 885, 893–901; S4 l. 20, 26, 31, 293 y 346; S5 l. 14, 15, 22, 28, 34 y 235. Ninguna marca PROVISIONAL va sin bloqueo (solo encabezados: S2 l. 458 y S3 l. 884) |
| d | Las cifras canon no cambian | CUMPLE | `nums.py` compara `HEAD` contra el árbol de trabajo en S2 a S5 y **no desaparece ninguna cifra**. Las nuevas son de ejemplos: S2 4,000.00, 199.00 y 100000; S4 2250, 1350 y 2300 (ejemplo r5); S5 1,021.00, 102100 y 0.00. En `verify-all`: oráculos S2 30/30, S3 906/906 (canon 437/437), S4 26/26 y S5 29/29, todos con 0 diferencias; `model.json` sha256 `1a57066689dd…` intacto |
| e | «Cambios que deben implementar los motores» cubre lo que cambió | PARCIAL (M-2) | **Fixtures:** `fix.py` encuentra 25 cuyo `provisional` cita un bloqueo RESUELTO y las 25 están cubiertas: S2, 2 en `contratos-puerta2.md`; S3, 9 por §14.1 y su lista; S4, 5; S5, 9. Además, `s3where.py` comprueba que el `S3-formula-tarjeta` de `s3-unificado`, `-sin-corte` y `comparador-caso-2` viene del calendario o de un corte `bank`, que son FIRMES por Q2-32. **Comentarios de código:** `code.py` encuentra 62 sitios con PROVISIONAL y un bloqueo resuelto; faltan 7 (M-2) |
| f | `make verify-all GOOSE=…/gobin/goose` | FALLA (A-1) | rc=2. WF-0 OK. **`verify-c1`: «gen2: tipo de parámetro desconocido: hora HH:MM»** (Makefile:130). S1 TODO OK, U1 63 OK, criterios 178 con 0 diferencias, INDICE con 0 diferencias contra su conteo propio, K25 (motor_cred 93.7 %, motor_fin 89.2 %, recetas 93.2 %, precios 96.1 %), WF-1 OK |

## Hallazgos
**A-1 (alta) · verificación global rota.** `make verify-all` termina con rc=2 porque `verify-c1` falla con `gen2: tipo de parámetro desconocido: hora HH:MM`. Lo provoca el tipo nuevo de `00004_puerta2.sql` y de `diccionario-datos.md`, que `backend/internal/api/c1gen/gen2_schemas.py` (`_pv`) no conoce. Es el A-2 de la ronda 1: está fuera del alcance de los contratos y sigue igual.
- Ruta: `backend/internal/api/c1gen/gen2_schemas.py`
- Corrección: que un agente con alcance en C1, o el integrador, aplique la solicitud 1 de `S1-puerta2.md`: el tipo «hora HH:MM» como string con patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$`. Después, `build.py --check` y `make verify-all`.

**M-1 (media) · tres bloqueos ABIERTOS citados por los contratos no están en el registro.** `S3-formula-tarjeta-pdf`, `S3-iva-lectura` y `S4-empates-sin-confirmar` tienen 0 menciones en `INDICE.md` y 0 en `decisiones.md` (K46). Los totales del INDICE (75 / 63 / 12; fila P2c 4 / 0 / 4) no los cuentan. Los bloqueos resueltos de los que salen tampoco remiten a ellos: `S3-formula-tarjeta.md`, `S4-confirmacion-empates.md` y `U4-iva-credito.md`. Si WF-3 toma el INDICE como lista de lo abierto, no los ve. Esto cierra en el contrato el M-3 de la ronda 1, pero mueve la inconsistencia al registro. Fuera de los contratos falta también `S1-reparto-edicion-gasto`, con 0 menciones.
- Ruta: `docs/estado/bloqueos/INDICE.md`
- Corrección: que el integrador o el agente de bloqueos aplique la solicitud 2 de `correccion-contratos-r1.md`:
  - dar de alta los 3 bloqueos (y `S1-reparto-edicion-gasto`);
  - ajustar los totales y el conteo de K46;
  - poner una línea de remisión en los 3 bloqueos resueltos.

**M-2 (media) · la lista para WF-3 deja fuera 7 comentarios PROVISIONAL del código (M-4 de la ronda 1, cerrado en parte).** `correccion-contratos-r1.md` promete los comentarios «con archivo y línea», pero no trae estos:
- `backend/internal/recetas/health.go` l. 31–32 (`S4-reglas-lista`, pregunta 3);
- `backend/internal/precios/compare.go` l. 362–363 (`S5-reglas-comparativa`, pregunta 2);
- `backend/internal/precios/profeco/import.go` l. 29–31 y l. 37–38 (`S5-profeco-almacen`, preguntas 3 y 1);
- `backend/internal/motor_cred/card.go` l. 9 (`S3-formula-tarjeta`, pregunta 4);
- `backend/internal/motor_cred/calendar.go` l. 55–56 (corte que sustituye pagos, `S3-formula-tarjeta`);
- `oraculo/cred_lib.py` l. 141 (`S3-formula-tarjeta`, pregunta 4).

La receta del reporte tampoco los encuentra. Su `grep` exige PROVISIONAL y el bloqueo en el mismo renglón, y en los 4 primeros archivos no hay ninguna otra coincidencia donde «revisar los bloques vecinos». Además, su patrón excluye `S3-formula-tarjeta`. Con la única regla escrita para ese bloqueo («→ `-pdf` en interés, mínimo y corte `computed`»), WF-3 tendría que adivinar qué hacer con los tres últimos, que son FIRMES por Q2-32 y no pasan a `-pdf`.
- Ruta: `docs/estado/p2c/correccion-contratos-r1.md`
- Corrección:
  - agregar los 7 sitios con su destino. Todos pasan a FIRME: `health.go` por Q2-46 = A; `compare.go` por Q2-60 = A; `import.go` l. 29–31 por Q2-55 = A y l. 37–38 por Q2-53 = A; `card.go` l. 9, `calendar.go` y `cred_lib.py` por Q2-32 = A;
  - agregar `card.go` l. 89 (pasa a `-pdf`) y las pruebas que hoy esperan listas viejas: `engine_test.go` l. 117, 135, 227, 243, 286 y 411; `timeline_test.go` l. 441 y 461; `comparador_test.go` l. 390;
  - cambiar la receta por un `grep -rn -A2 -B2 PROVISIONAL` y dejar de excluir `S3-formula-tarjeta` (sin `-pdf`).

**B-1 (baja) · `S2-ventanas.md` l. 85.** La línea «Qué falta» de Q2-01 dice que las opciones «las redacta quien aplique Q2-01» y que «esta línea se sustituye por la pregunta completa». La pregunta ya está escrita desde l. 90, pero la línea no se sustituyó. Corrección: remitir a «Pregunta 4» (l. 90).

**B-2 (baja) · bajos de la ronda 1 sin atender.**
- B-1: S2 l. 354, fila de `s2-liga-provisional` junto a la marca de `S2-flujo-quincenal`.
- B-2: S2 l. 359 dice «27 fixtures», pero en `oraculo/fixtures/s2` hay 30.
- B-3: `contratos-puerta2.md` l. 57, «el motor las agrega en orden del plan».

Corrección: la que da `REVISION-contratos-r1.md` l. 55–59.

**B-3 (baja) · `S3-iva-lectura` pregunta 2.** La lectura PROVISIONAL (`ErrMissingCapture` con `VATIncluded` = false y sin porcentaje) no es «la opción en uso», porque en 1.3 ese estado no existía. Es razonable: no da número (RN-26) y cita el bloqueo ABIERTO. Pero el bloqueo debería decirlo así. Corrección: en «Cómo seguí», agregar «no hay opción en uso; se tomó la que no da número (RN-26), PROVISIONAL».

**B-4 (baja) · S4 l. 20 (bitácora M-r3-1).** Dice «PROVISIONAL y con una variante sin campo nuevo», pero desde la puerta 2 rige «Sin variante» (S4 l. 446). Es historia, pero puede confundir. Corrección: añadir «(superado en la puerta 2: sin variante, §12.2)».

## Estado de los hallazgos de la ronda 1
| Ronda 1 | Estado | Evidencia |
|---|---|---|
| A-1 (alta) · S3 y Q2-67 = 1 | CERRADO | S3 l. 116 trae `VATIncluded *bool`. l. 258 y §5.8 (l. 310–312) piden 5 datos. «4 datos» y «lista de 4» ya no aparecen (`grep`). El préstamo con null queda PROVISIONAL y cita `S3-iva-lectura` (ABIERTO). La semilla Kueski lleva `"vat_included": true` (l. 679). Las pruebas están en l. 782 y 789. La lista de cambios lo incluye |
| A-2 (alta) · `verify-c1` | SIGUE ABIERTO → A-1 | mismo error en mi corrida |
| M-1 (media) · marcas que citaban bloqueos resueltos en S4 y S5 | CERRADO | S4 l. 293, 342, 346 y 362; S5 l. 235 y 244. `prov.py` da 0 marcas activas |
| M-2 (media) · S3 §5.8 `bank_card` y l. 594 | CERRADO | l. 310: capital = monto, fechas y `quoted`, todos FIRMES; l. 595 remite a §14.1 |
| M-3 (media) · `S3-formula-tarjeta` RESUELTO y citado | CERRADO en el contrato (cita `S3-formula-tarjeta-pdf`, ABIERTO); el registro queda en M-1 | S3 l. 7, 258, 315, 885 y 901 |
| M-4 (media) · lista de cambios incompleta | PARCIAL → M-2 | Están las 14 fixtures y las de S3; faltan 7 comentarios |
