# P2c · correccion-contratos-r2 · 11-oct-2026, noche

Objeto: los hallazgos alto y medios de `docs/estado/p2c/REVISION-contratos-r2.md` (A-1, M-1 y M-2), sobre lo que escribieron `contratos-puerta2` y `correccion-contratos-r1`.
Alcance: `S2-interfaz.md` a `S5-interfaz.md`, `S2-ventanas.md`, los bloqueos nuevos, `oraculo/fixtures/FORMATO.md` y este reporte.
Fuentes: `RESPUESTAS.md`; `PUERTA-2.md` §3.9 (Q2-32, l. 162), §3.15 (Q2-46, l. 205–206), §3.18 (Q2-53 y Q2-55, l. 228–231) y §3.20 (Q2-60, l. 245–251); `S3-interfaz.md` 1.5 §5.3 (l. 258), §11 (l. 789–790) y §14.1 (l. 900–930); los bloqueos `S3-formula-tarjeta`, `S4-reglas-lista`, `S5-reglas-comparativa` y `S5-profeco-almacen` (secciones «### Q2-…», todas RESUELTO); el código citado, leído renglón por renglón.
No hice `git add` ni `commit`. No toqué la Pista 3 ni los hallazgos bajos (B-1 a B-4), porque la orden pide solo el alto y los medios.

## Hecho
- **Ningún archivo de mi ALCANCE cambió.** Los tres hallazgos piden cambios fuera de él (`gen2_schemas.py`, `INDICE.md`, `decisiones.md`, tres bloqueos resueltos y la lista para WF-3). Los contratos ya citan bien los bloqueos (verificación c). El sha256 de los 9 archivos de mi alcance es el mismo antes y después.
- **A-1 (alta) · `verify-c1`. No lo corregí:** está fuera de mi alcance. Mi corrida da el mismo error que la revisión: `gen2: tipo de parámetro desconocido: hora HH:MM` (Makefile:130). Va en la solicitud 1.
- **M-1 (media) · registro de los bloqueos.**
  - En mi alcance. Los 3 bloqueos nuevos (`S3-formula-tarjeta-pdf`, `S3-iva-lectura` y `S4-empates-sin-confirmar`) ya terminan con «- **Estado:** ABIERTO». Es el mismo formato de los 4 de `contratos-puerta2` que el INDICE ya cuenta, así que no hizo falta tocarlos. Ninguno trae «Respuesta de Axel», así que `verify_c1.py` (l. 322–323) también los cuenta como abiertos.
  - Fuera de alcance. Preparé el cambio exacto para el integrador en `datos/tmp/correccion-contratos-r2/aplica-registro.py`: 11 reemplazos, cada uno con su conteo esperado. Sin argumentos solo ensaya y no escribe. El texto completo va también en el Anexo, porque `datos/tmp` lo ignora git y se puede borrar. Esto es lo que cambia:
    - **INDICE**: la fila P2c pasa de 4 / 0 / 4 a 8 / 0 / 8 y el total, de 75 / 63 / 12 a 79 / 63 / 16. El párrafo de conteo gana una frase. La regla de las herramientas pasa de «67 resueltos y 8 abiertos» a «67 y 12». La sección «Bloqueos nuevos de P2c» gana 4 filas: los 3 de arriba y `S1-reparto-edicion-gasto`.
    - **K46**: «conteo de 75» pasa a «75 al cerrar las respuestas y 79 con las correcciones».
    - **Remisión**: una línea antes del «Estado del bloqueo» final de `S3-formula-tarjeta.md`, `S4-confirmacion-empates.md` y `U4-iva-credito.md`.
  - Conteo en disco al 11-oct-2026, noche: **79 archivos**, 63 RESUELTO y 16 abiertos. Las 4 que faltan en el INDICE son las que señala la revisión, con 0 menciones cada una.
- **M-2 (media) · lista para WF-3.** Va completa en «Lista para WF-3 (complemento r2)». Esa sección **sustituye la receta de `correccion-contratos-r1.md`** y complementa sus listas y las de `contratos-puerta2.md`. Trae:
  - los 7 sitios de la revisión con su destino FIRME, comprobado en PUERTA-2 y en cada bloqueo;
  - `card.go` l. 89, que pasa a `-pdf`;
  - las 9 pruebas que nombra la revisión y 4 más que encontró la receta nueva: `engine_test.go` l. 171 y 183, `comparador_test.go` l. 394 y `timeline_test.go` l. 460;
  - constantes y generadores con nombres de bloqueos resueltos. Uno no estaba en ninguna lista: `backend/internal/motor_cred/comparador/testdata/tabla_mano.py`, con 8 renglones.
  - Corrijo dos números de renglón de la revisión: el comentario de `import.go` sobre `RecordSink` está en l. 36–37, no en 37–38, y el de `compare.go` va de l. 362 a 364.

## Verificación
- `sha256sum` de `S2-interfaz.md` a `S5-interfaz.md`, `S2-ventanas.md`, `FORMATO.md` y los 3 bloqueos nuevos, antes (`datos/tmp/correccion-contratos-r2/sha-antes.txt`) y después: **9 de 9 iguales** (`sha256sum -c`).
- `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP=…/datos/tmp/correccion-contratos-r2/vw`: **rc=2**. El log está en `datos/tmp/correccion-contratos-r2/verify-all.log`.
  - WF-0: «verificación completa OK».
  - **verify-c1: «gen2: tipo de parámetro desconocido: hora HH:MM», `make: *** [Makefile:130: verify-c1] Error 1`** (A-1, fuera de alcance).
  - S1: «RESUMEN TODO OK». U1: 63 OK y 0 pendientes.
  - Criterios: 178 comprobadas con 0 diferencias, 261 definidos y el INDICE de criterios con 0 diferencias.
  - K25: motor_cred 93.7 %, motor_fin 89.2 %, recetas 93.2 % y precios 96.1 %.
  - Oráculos con 0 diferencias: S2 30/30; S3 27 fixtures (906/906, canon 437/437); S4 26/26; S5 29/29.
  - `model.json` sha256 `1a57066689dd…` intacto. WF-1: «verificación completa OK».
- `prov.py` (mío) cruza cada renglón con «provisional» de S2 a S5 contra el «Estado» de cada bloqueo. Los renglones que citan un bloqueo no ABIERTO son exactamente los de bitácora o de negación que lista la revisión:
  - S3 l. 7, 258, 273, 315, 790, 885, 893, 894, 896, 898 y 901;
  - S4 l. 20, 26, 31, 293 y 346;
  - S5 l. 14, 15, 22, 28, 34 y 235;
  - S2: ninguno.
  - Los 8 bloqueos citados como abiertos están ABIERTO: `S3-formula-tarjeta-pdf`, `S3-iva-lectura`, `S4-empates-sin-confirmar`, `S4-presentacion-respaldo`, `S5-profeco-filtro`, `S5-carritos-granel`, `S2-ventanas` y `S1-reparto-edicion-gasto`.
  - Citas en los contratos: `S3-formula-tarjeta-pdf` 21 (S3), `S3-iva-lectura` 15 (S3) y `S4-empates-sin-confirmar` 4 (S4).
- Las recetas de §E, al 11-oct-2026, noche:
  - **receta 1: 40 renglones**. Coincide con mi `scan.py`: 40 renglones únicos, 41 menciones, porque `rec_plan.py` l. 3 nombra dos bloqueos. Incluye los 7 sitios y `card.go` l. 89.
  - **receta 2: 106 renglones** en `backend/internal` y `oraculo`. Son 98 de `scan.py` más los 8 de `tabla_mano.py`, que `scan.py` no veía porque solo leía `.go` en `backend`.
- `python3 -I datos/tmp/correccion-contratos-r2/aplica-registro.py` (ensayo): «reemplazos: 11 todos con el conteo esperado».
  - INDICE: −5 / +9 renglones. `decisiones.md`: −1 / +1. Cada bloqueo resuelto: +2.
  - «ensayo: no se escribió nada». El md5 de INDICE y `decisiones.md` no cambió.

## Supuestos
1. «Bloqueos nuevos», en mi alcance, son los archivos creados en P2c por la cadena de contratos. `INDICE.md`, `decisiones.md` y los bloqueos RESUELTOS no lo son. La corrección que sugiere la revisión para M-1 ya los asigna al «integrador o agente de bloqueos».
2. No edité `correccion-contratos-r1.md`, que es el reporte de otra ronda. La lista y la receta corregidas viven en este reporte, que lo sustituye en esa parte (solicitud 3).
3. Destinos de las pruebas: los tomé de la tabla de S3 §14.1, sin regla nueva. Hay dos casos que dependen del dato de entrada del caso, y en los dos doy ambas salidas tal como las fija el contrato:
   - `engine_test.go` l. 243 depende de `VATIncluded`;
   - `timeline_test.go` l. 461, préstamo `l1`, depende de su `vat_included`.

## Bloqueos
- Ninguno nuevo.
- A-1 sigue: `make verify-all` está roto en `verify-c1` hasta que alguien con alcance en C1 aplique la solicitud 1.

## Pendientes
- Las solicitudes fuera de alcance 1 a 3: integrador o agente de C1.
- Los bajos de la revisión, que no pidió esta orden:
  - B-1: `S2-ventanas.md` l. 85;
  - B-2: los bajos de la ronda 1 de S2;
  - B-3: «Cómo seguí» de `S3-iva-lectura`;
  - B-4: S4 l. 20.
- Ronda 3 del revisor.

## Lista para WF-3 (complemento r2; sustituye la receta de r1)
**A. Los 7 sitios de M-2: comentario PROVISIONAL de un bloqueo resuelto → FIRME.** Se quita «PROVISIONAL, bloqueo …» y se cita la respuesta.

| Archivo y renglones | Cita hoy | Pasa a | Fuente |
|---|---|---|---|
| `backend/internal/recetas/health.go` l. 31–32 | `S4-reglas-lista`, pregunta 3, opción A (IMC con el valor que se muestra) | FIRME, Q2-46 = A | PUERTA-2 l. 205–206 «Recomendada: A en las tres»; `S4-reglas-lista.md` «### Q2-46», RESUELTO |
| `backend/internal/precios/compare.go` l. 362–364 | `S5-reglas-comparativa`, pregunta 2, opción A (empates de `rankStores`) | FIRME, Q2-60 = A | PUERTA-2 §3.20: «A en 1 a 3 y 4b, B en 4a»; «### Q2-60», RESUELTO |
| `backend/internal/precios/profeco/import.go` l. 29–31 | `S5-profeco-almacen`, pregunta 3, opción A (el primero gana) | FIRME, Q2-55 = A | PUERTA-2 l. 230–231 «Recomendada: A en las tres»; «### Q2-55», RESUELTO |
| `backend/internal/precios/profeco/import.go` l. 36–37 | `S5-profeco-almacen`, pregunta 1, opción A (`qqp.sqlite` aparte) | FIRME, Q2-53 = A | lo mismo; «### Q2-53» |
| `backend/internal/motor_cred/card.go` l. 8–9 | `S3-formula-tarjeta`, pregunta 4 (`cardDueDates`) | FIRME, Q2-32 = A | PUERTA-2 l. 162 «Recomendada: A»; S3 l. 7: «el calendario de la tarjeta y el corte que sustituye pagos (Q2-32 = A)» |
| `backend/internal/motor_cred/calendar.go` l. 55–56 | `S3-formula-tarjeta` (el corte sustituye los pagos) | FIRME, Q2-32 = A | lo mismo |
| `oraculo/cred_lib.py` l. 141–142 | `S3-formula-tarjeta`, pregunta 4 (`card_dues`) | FIRME, Q2-32 = A | lo mismo |

**B. Pasa a `-pdf`.** `backend/internal/motor_cred/card.go` l. 89 dice «Fórmula general PROVISIONAL (bloqueo S3-formula-tarjeta, opción b)». Sigue PROVISIONAL y pasa a citar `S3-formula-tarjeta-pdf` (S3 l. 258 y §14.1).

**C. Pruebas que esperan listas viejas** (destino según S3 §14.1; orden alfabético por bytes)

| Archivo y renglón | Hoy | Pasa a |
|---|---|---|
| `motor_cred/engine_test.go` l. 117 («cargo normal») | [`S3-formula-tarjeta`, `S3-origen-cargo-normal`] | `[]` (`Purchase` de tarjeta) |
| `engine_test.go` l. 135 («MSI») | [`S3-formula-tarjeta`, `S3-serie-sin-intereses`] | `[]` |
| `engine_test.go` l. 171 («interés») · lo encontró la receta 2 | [`S3-formula-tarjeta`, `U4-iva-credito`] | [`S3-formula-tarjeta-pdf`] (interés con la fórmula) |
| `engine_test.go` l. 183 («mínimo») · receta 2 | lo mismo | [`S3-formula-tarjeta-pdf`] |
| `engine_test.go` l. 227 («préstamo con IVA») | [`S3-convencion-prestamo`] | `[]` (`VATIncluded` = false con su porcentaje) |
| `engine_test.go` l. 243 («préstamo sin IVA») | [`S3-convencion-prestamo`, `U4-iva-credito`] | `VATIncluded` nil → [`S3-iva-lectura`]; true → `[]` (§11 `TestDespachoCotizacion` pide true, false y null) |
| `engine_test.go` l. 286 («tienda») y l. 411 («promoción») | [`S3-serie-sin-intereses`] | `[]` |
| `motor_cred/timeline_test.go` l. 441 («corte», `Source` `bank` en l. 427) | [`S3-formula-tarjeta`] | `[]` (corte `bank`) |
| `timeline_test.go` l. 460–461 (`TestProvisionalListas`) | `c1`, `c2`, `l1` y `t1` con nombres viejos | `c1`, `c2` y `t1`: `[]`. `l1`: `[]` si el préstamo lleva `vat_included` = false (r1 lo pide para el préstamo de `s3-unificado`); [`S3-iva-lectura`] solo si llega null. `d1`: sigue `{}` |
| `motor_cred/comparador/comparador_test.go` l. 390 (entrada) y l. 394 (lo esperado) · la l. 394 la encontró la receta 2 | [`S3-formula-tarjeta`] | [`S3-formula-tarjeta-pdf`] en los dos (la entrada es un corte `computed`) |

**D. Constantes y generadores** (receta 2; no llevan PROVISIONAL en el renglón)
- `motor_cred/revolving.go` l. 11–15 y `oraculo/cred_lib.py` l. 13–18 (constantes de nombres). Desde la puerta 2 solo entran `S3-formula-tarjeta-pdf` y `S3-iva-lectura` (§14.1). Las constantes de los resueltos dejan de usarse.
- `motor_cred/comparador/request.go` l. 131 (`provObjectives`): `S3-objetivos-medida` deja de entrar en `Scenario` y `Result` (`contratos-puerta2.md`).
- **`motor_cred/comparador/testdata/tabla_mano.py`** l. 18, 113, 120, 123, 452–454 y 496: la tabla a mano de las fixtures del comparador. No estaba en ninguna lista. Las listas por compra pasan a `[]`; `OTROS` y el `provisional` de l. 496 pasan a null, porque las fixtures `s3-comparador-*` van en null desde la puerta 2 (§14.1, tabla de fixtures).
- `oraculo/cred_genera_fixtures.py` l. 22–26, 380, 415–417, 511–513, 521, 550, 600, 603 y 639: listas y textos de las fixtures de S3. `s3-tarjeta-minimo` y `s3-tienda-abonos-interes` pasan a `-pdf`; las demás, a null (§14.1 y la lista de r1).
- `oraculo/pre_genera_fixtures.py`: l. 270, 324, 403, 426 y 458 pasan a `provisional` null (las 8 fixtures de S5 de r1). l. 596 (`s5-profeco-muestra`) pasa a `S5-profeco-filtro`, con el `que` de r1.
- `oraculo/rec_genera_fixtures.py`: l. 16–17 (`BL_LISTA`, `BL_MERCADO`) dejan de usarse; l. 422 (`BL_EMPATES`) pasa a `docs/estado/bloqueos/S4-empates-sin-confirmar.md` (r1).
- S2: ya los cubre `contratos-puerta2.md` («Quitar los comentarios PROVISIONAL…» y «`provisional` siempre `null`…»); aquí van con archivo y línea.
  - `motor_fin/input.go` l. 214;
  - `plan.go` l. 30, 43 y 415;
  - `projection.go` l. 25, 113 y 127;
  - `engine_test.go` l. 299;
  - `oraculo/fin_metas.py` l. 3 y 5;
  - `fin_presupuesto.py` l. 20;
  - `fin_proyeccion.py` l. 4 y 12;
  - `fin_genera_fixtures.py` l. 333, 350, 351 y 451.

**E. Receta (sustituye la de `correccion-contratos-r1.md`; no excluye `S3-formula-tarjeta`, sí su `-pdf`)**
```
P='S3-(convencion-prestamo|objetivos-medida|origen-cargo-normal|serie-sin-intereses)|S3-formula-tarjeta([^-]|$)|U4-iva-credito|S4-(asignacion-mercado|confirmacion-empates|reglas-lista)|S5-(reglas-comparativa|profeco-almacen)|S2-(metas|presupuesto-total|proyeccion-interruptor)([^-]|$)|S2-metas-resumen'
# 1: comentarios PROVISIONAL, con 2 renglones de contexto (hoy 40 renglones)
grep -rn -B2 -A2 PROVISIONAL --include='*.go' --include='*.py' backend/internal oraculo | grep -E "$P"
# 2: toda mención de un bloqueo resuelto: listas, constantes, pruebas y generadores (hoy 106 renglones)
grep -rnE "$P" --include='*.go' --include='*.py' backend/internal oraculo
```
El patrón es el de los bloqueos RESUELTOS de S2 a S5 y `U4-iva-credito` al 11-oct-2026, noche. `S2-metas-vencida`, `S2-flujo-quincenal` y `S2-ventanas` siguen ABIERTOS y quedan fuera. Si un bloqueo cambia de estado, cambia el patrón.

## Solicitudes fuera de alcance
1. **A-1 (agente de C1 o integrador).** Enseñar a `backend/internal/api/c1gen/gen2_schemas.py` (`_pv`) el tipo «hora HH:MM»: string con patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$` (solicitud 1 de `S1-puerta2.md`). Después, `build.py --check` y `make verify-all`.
2. **M-1 (integrador).**
   - Correr `python3 -I datos/tmp/correccion-contratos-r2/aplica-registro.py`. Si dice «11 todos con el conteo esperado», correrlo con `--aplicar`.
   - Si `datos/tmp` ya no existe, aplicar a mano el Anexo.
   - Después, revisar que la fila P2c del INDICE dé 8 / 0 / 8 y el total, 79 / 63 / 16.
3. **M-2 (integrador).** Unir «Lista para WF-3 (complemento r2)» a las de `contratos-puerta2.md` y `correccion-contratos-r1.md` en la INTEGRACION de P2c. Conviene una línea en `correccion-contratos-r1.md` que remita aquí, porque su receta queda sustituida.
4. **Observación (integrador; no es de los contratos).** `backend/cmd/finanzas/main.go` l. 6 cita `criterios-bind`, que está RESUELTO, en un «Mientras tanto». La receta 2 no lo incluye.

## Propuestas para Axel
Ninguna nueva. Siguen las de r1: contestar por id `S3-iva-lectura` y Q2-90 (`S4-empates-sin-confirmar`), y entregar los PDF de `S3-formula-tarjeta-pdf`.

## Anexo · texto exacto de la solicitud 2 (el mismo que el script)
- **INDICE, fila del resumen:**
  - antes: `| Flujo P2c (…) | 4 | 0 | 4 | abiertos por contratos-puerta2 al aplicar Q2-03, Q2-07, Q2-13, Q2-16 y Q2-18 |`
  - después: `| 8 | 0 | 8 |` con «4 abiertos por `contratos-puerta2` al aplicar Q2-03, Q2-07, Q2-13, Q2-16 y Q2-18; 3 por `correccion-contratos-r1` (Q2-29 a Q2-31 = (a) sin los PDF, Q2-67 = 1 y Q2-90 sin respuesta) y 1 por `correccion-S1-r1` (Q2-68 = C)».
- **INDICE, fila Total:** `| **Total** | **79** | **63** | **16** | |`.
- **INDICE, al final del párrafo «Conteo…»**, después de «…la fuente es la carpeta, no este párrafo.»: «**Conteo al 11-oct-2026, noche (solicitud de `correccion-contratos-r2`, hallazgo M-1 de `docs/estado/p2c/REVISION-contratos-r2.md`, con la carpeta en disco):** 79 archivos, **63 resueltos y 16 abiertos**: los 75 de arriba, más 4 que abrieron las correcciones de P2c (`S3-formula-tarjeta-pdf`, `S3-iva-lectura` y `S4-empates-sin-confirmar` de `correccion-contratos-r1`; `S1-reparto-edicion-gasto` de `correccion-S1-r1`). Las 16 abiertas: las 6 de WF-1, las 8 de P2c y las 2 de la Pista 3.»
- **INDICE, regla de las herramientas:**
  - «cuentan 67 resueltos y 8 abiertos en la carpeta de hoy» → «… y 12 abiertos …»;
  - «los 4 de P2c y los 2 de la Pista 3» → «los 8 de P2c …».
- **INDICE, entrada de «Bloqueos nuevos de P2c»:** «Los 4 están **ABIERTOS**» → «Después, las correcciones de P2c abrieron 4 más: 3 `correccion-contratos-r1` (`docs/estado/p2c/correccion-contratos-r1.md`) y 1 `correccion-S1-r1` (`docs/estado/p2c/correccion-S1-r1.md`); esos 4 los dio de alta el integrador a pedido de `correccion-contratos-r2`. Los 8 están **ABIERTOS**».
- **INDICE, 4 filas nuevas** después de la de `S5-carritos-granel.md`. Columnas: Archivo · Nació de · Pregunta · Opciones · Qué queda PROVISIONAL · Fase. El texto completo de cada celda está en `ROWS` del script.
  - `S3-formula-tarjeta-pdf.md`
    - Nació de: Q2-29 a Q2-31 = (a), decisión sin resultado.
    - Pregunta: entregar los PDF de tarjeta y de crédito de tienda.
    - Opciones: 1, entrega; 2, excepción a (b) o (c). Sin recomendación.
    - Queda PROVISIONAL: la fórmula general (b), los cortes `computed` y dos fixtures. FIRMES: $0.00, MSI, abonos sin interés y calendario (Q2-32 = A).
    - Fase: S3 (S3-08, S3-10) y WF-3.
  - `S3-iva-lectura.md`
    - Nació de: Q2-67 = 1, dos casos que no cubre.
    - Pregunta 1: préstamo con `vat_included` NULL (A no suma IVA, en uso; B `ErrMissingCapture`; C la API no lo deja guardar).
    - Pregunta 2: 0 sin porcentaje (A `ErrMissingCapture` y `vat` en `missing_terms`; B la API no lo deja guardar y el motor responde `ErrInvalid`). Sin recomendación.
    - Queda PROVISIONAL: en la pregunta 1, la A con [`S3-iva-lectura`]; en la 2, `ErrMissingCapture`.
    - Fase: S3, C1, U4 y WF-3.
  - `S4-empates-sin-confirmar.md`
    - Nació de: Q2-43 = A, y Q2-90 sin respuesta.
    - Pregunta: confirmar o corregir los empates de `revisar.tsv` y, cuando exista, las ligas PROFECO.
    - Opciones: A confirmar, B corregir. Sin recomendación.
    - Queda PROVISIONAL: los 150 renglones `propuesto` de `tabla.csv` y las fixtures `s4-kcal-*`. En `production` no se carga ningún `propuesto`.
    - Fase: S4 (S4-08), S5 y WF-3.
  - `S1-reparto-edicion-gasto.md`
    - Nació de: Q2-68 = C (`U4-reparto-retroactivo`); lo abrió `correccion-S1-r1`.
    - Pregunta: con qué partes se recalcula un gasto anterior que se edita después de un cambio «solo los nuevos».
    - Opciones: A proporción congelada, B guardar cómo se calculó, C regla vigente. Sin recomendación.
    - Queda PROVISIONAL: la C, en RN-03 y en la invariante 18.
    - Fase: C1, S1, S7, U4 (P-02) y WF-3.
- **K46** (`docs/plan/decisiones.md` l. 120): «ya está en el conteo de 75 del INDICE)» → «ya está en el conteo del INDICE: 75 archivos al cerrar las respuestas y 79 (63 resueltos y 16 abiertos) con los 4 que abrieron después las correcciones de P2c, `S3-formula-tarjeta-pdf`, `S3-iva-lectura`, `S4-empates-sin-confirmar` y `S1-reparto-edicion-gasto`)».
- **Remisiones**, antes de la línea final «**Estado del bloqueo (puerta 2): RESUELTO**»:
  - `S3-formula-tarjeta.md`: «Remisión (P2c, 11-oct-2026, noche): la respuesta (a) a las preguntas 1 a 3 (Q2-29 a Q2-31) espera los estados de cuenta en PDF; lo que sigue abierto vive en `docs/estado/bloqueos/S3-formula-tarjeta-pdf.md` (ABIERTO, lo abrió `correccion-contratos-r1`). La pregunta 4 (Q2-32 = A: calendario de la tarjeta y corte que sustituye pagos) es FIRME.»
  - `S4-confirmacion-empates.md`: «Remisión (P2c, 11-oct-2026, noche): quién confirma y cuándo es FIRME (Q2-43 = A). La confirmación de los valores `propuesto` (Q2-90, sin respuesta) vive en `docs/estado/bloqueos/S4-empates-sin-confirmar.md` (ABIERTO, lo abrió `correccion-contratos-r1`).»
  - `U4-iva-credito.md`: «Remisión (P2c, 11-oct-2026, noche): Q2-67 = 1 es FIRME. Los dos casos que no cubre (préstamo con `vat_included` NULL; `vat_included` = 0 sin porcentaje) viven en `docs/estado/bloqueos/S3-iva-lectura.md` (ABIERTO, lo abrió `correccion-contratos-r1`).»

## Consumo aproximado
34 acciones de herramienta · unos 125 KB leídos · unos 35 KB escritos: este reporte, unos 18 KB; `aplica-registro.py`, unos 9 KB; `scan.py` y `prov.py`, unos 4 KB; sin contar el log de `verify-all`.
