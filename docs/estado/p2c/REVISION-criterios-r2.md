# Revisión de criterios · ronda 2 de 3 (P2c) · `revisor-criterios-r2`

Fecha: 11-oct-2026, noche. Revisor independiente: no escribí nada de lo que reviso y no corregí nada. Evalúo solo contra `docs/criterios/tools/verif_final.sh`, `docs/estado/puerta2/PUERTA-2.md` §4 y §5 y los reportes de P2c (`criterios-puerta2.md`, `correccion-criterios-r1.md`). Árbol revisado: HEAD 4079cee (P3a) con los cambios de P2c sin commit; `git diff --stat -- docs/criterios/` = 9 archivos (C1, INDICE, S1, S2, S4, S5, S7, U3, U4), 164 líneas nuevas y 85 quitadas. Mis scripts y la base temporal estuvieron en `datos/tmp/rev-crit-r2/` y los borré al terminar.

## Veredicto: **NO APROBADO** · 0 altos, 1 medio, 2 bajos

Los tres hallazgos de la ronda 1 (H-1, M-1 y M-2) están **cerrados**. El medio es nuevo y viene de cómo se redactó el cierre de M-1. En C1-15 la regla de «bloqueo abierto» ahora pide «Estado: ABIERTO» **y** que no haya «Respuesta de Axel». Con esa regla, las marcas que S2-16 y S5-13 exigen pasan a ser defectos. Además, dos de los cuatro bloqueos que C1-15 da como ejemplo no tienen esa línea de estado.

## Verificaciones (con mis propios comandos)

| # | Qué | Cómo | Resultado |
|---|---|---|---|
| a | `make verify-criterios` (11 comprobaciones) | `make verify-criterios` (1.4 s) | `criterios definidos: 261`, `IDs distintos citados: 261`, `citados que no existen: {}`; `citas verificadas: 321`, `sin fuente: 0`; canon `TOTAL comprobadas: 178`, `con diferencias: 0`; `cifras comparadas: 12`, `diferencias: 0`; `N-xxx sin operacion ni Local: 0`; matriz C1 236 filas y 211 N-xxx, 0 ausentes; `TOTAL 261 (alta 201 / media 60 / baja 0)`, 18 archivos; `INDICE contra conteo propio: 0 diferencias`; `residuos «N-xx» y rutas /tmp: 0 y 0; las 11 comprobaciones dan lo esperado`. **Cumple.** |
| b | 261 criterios, 201/60 | `cuenta_rev.py`, que es mío: usa una regex propia, lee los 18 archivos y compara cada uno contra `git show HEAD:` | 18 archivos, 261 criterios (201 alta, 60 media, 0 baja) y 261 ids únicos. **0 cambios de id o de severidad contra HEAD.** Ningún renglón tiene una severidad escrita fuera de alta, media o baja. **Cumple.** |
| c1 | Valores FIRMES nuevos contra la 00004 | `goose up` de 00001 a 00004 en una base temporal (`goose: successfully migrated database to version: 4`) y `SELECT` de `app_setting` del hogar | `reports.daily_send_time` "21:30", `reports.weekly_weekday` 7, `reports.weekly_send_time` "19:00", `reports.monthly_day_of_month` 1, `reports.monthly_send_time` "08:00", `alerts.send_time` "09:00", `alerts.stale_prices_advance_days` 1, `alerts.meal_plan_missing_weekday` 7, `optimizer.date_window_days` 7, `ocr.timeout_seconds` 30. Siguen igual `comparator.timeout_seconds` 5 y `price.stale_after_days`/`zero_weight_after_days` 30/90. Son 10 claves nuevas en el `INSERT` de la 00004 (l. 158–180), las mismas «diez claves» que dice C1-22. |
| c2 | Lo que dicen los criterios contra la semilla | `firmes.py`, que es mío: busca `` `clave` = valor `` en los 19 `.md` | 35 menciones en total. 34 coinciden con la semilla y la restante es el escenario de prueba de S7-10 (`monthly_day_of_month` = 31 en noviembre, que sale el 30). Mi script marcó 10 «diferencias» al principio, pero 9 las causé yo: había puesto en 0 `ocr.timeout_seconds` y la anticipación al probar los rangos. Revisé cada una a mano. **Cumple.** |
| c3 | Rangos de C1-22 contra los disparadores de 00004 | `UPDATE` de prueba | `comparator.timeout_seconds` rechaza 0 y 31 y acepta 30. `price.stale_after_days` rechaza 0. `reports.weekly_weekday` rechaza 8. `reports.monthly_day_of_month` rechaza 32. `ocr.timeout_seconds` y `alerts.stale_prices_advance_days` aceptan 0 («0 o más», como dice C1-22). Las horas "24:00" y "9:00" se rechazan por el patrón HH:MM. **Coincide** con C1-22, con Q2-74 y Q2-77 y con K50 (`decisiones.md` l. 154). |
| c4 | Criterios contra PUERTA-2 §4 | Lectura de PUERTA-2 l. 291–302 y de `decisiones.md` l. 151 (K47) y l. 154 (K50) | Q2-72: 21:30, domingo 19:00 y día 1 a las 08:00. Q2-73: 1 día y 09:00. Q2-74: 5 s, con límites de 1 a 30. Q2-75: 7. Q2-76: 30. Q2-77: S ≥ 1 y Z ≥ S, validado en el motor y la API. Lo que dicen C1-22, U4-12, S7-04, S7-08, S7-10, S5-04, S5-06 y S1-16 coincide. **Cumple.** |
| c5 | PUERTA-2 §5: las preguntas sin recomendación | `grep` de Q2-78 a Q2-98 en los criterios | Q2-85 (U3 l. 8) y Q2-87 a Q2-90 y Q2-94 (S4 l. 10 y l. 124) dicen «sin respuesta, pasan a WF-3». Q2-91 tiene 0 menciones y «Express» o «Wal-mart» también 0, así que ningún criterio decide la pregunta. Q2-86 = A en U3-01. Ningún criterio decide algo «sin recomendación». **Cumple.** |
| d | S5-07 mide `docs/estado/S5-terminos.md` | `ls`, `grep -c "Veredicto"` y `go test ./internal/precios/... -run LecturaAsistida -count=1 -v` | La Verifica usa la ruta absoluta `/home/deepsec/finanzas-hogar/docs/estado/S5-terminos.md`, de 4.5 KB. No usa el reporte `docs/estado/p2c/S5-terminos.md`, de 2.2 KB. `Veredicto` aparece 3 veces (l. 13, 20 y 27) y las 3 dicen «no verificado». `TestLecturaAsistidaSinVeredicto` PASS. Hay 0 dominios de tiendas en `.go` fuera de las pruebas. El criterio declara «NO CUMPLE en el fondo» y no se marca cumplido. **Cumple.** |
| e | Ningún bloqueo resuelto citado como abierto | `bloq.py` y `citas_e.py`, que son míos. Toman el estado de cada uno de los 84 bloqueos de la última línea «Estado del bloqueo (puerta 2)» o, si no la tiene, de «Respuesta de Axel». Después revisan la oración de cada mención | **0 casos en el texto vigente.** Hubo 17 avisos y ninguno es defecto. 14 están en listas históricas («Cambios de la versión 1.2» e INDICE §11). S2 l. 90 dice bien que las preguntas 1 a 3 de `S2-ventanas` están resueltas. S3 l. 50 y la ruta de `S3-tiempo-optimizador.md` (S3-12) son condicionales: «si…, abre el bloqueo». `criterios-nxx` está cerrado con la opción (b) y no tiene respuesta de Axel, como dice INDICE. **Cumple.** Encontré el problema contrario, que se describe en M-1 abajo. |
| f | Verificación global | `make verify-all GOOSE=…/goose` (21.6 s) | rc 0, `Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)`. Las 5 líneas con «error» o «FAIL» que encontró grep dicen «sin errores» o «aceptadas por error []». |

## Hallazgos

### M-1 · media · C1-15 e INDICE §6 definen «bloqueo abierto» de forma que contradice S2-16 y S5-13 y no se cumple en los bloqueos que citan
- **Ruta:** `docs/criterios/C1.md` l. 103 (cuerpo de C1-15) y l. 104 (Verifica); `docs/criterios/INDICE.md` l. 68 (§6, «Estado = …»).
- **HECHO 1.** Hoy conviven tres reglas distintas:
  - El cuerpo de C1-15 dice: «"Bloqueo abierto" significa un bloqueo […] que no tiene sección "Respuesta de Axel"; uno que la tiene está RESUELTO y su marca se quita».
  - La Verifica de C1-15 corre `grep -rn "PROVISIONAL\|x-provisional"` sobre **todo** `docs/contratos/`, incluido `motores/`, y pide «0 marcas que citen un bloqueo con "Respuesta de Axel"». También dice que los bloqueos que cita son «abiertos, "Estado: ABIERTO" y sin "Respuesta de Axel"; la regla general es esa».
  - INDICE §6 dice: «"RESUELTO" si trae la sección "Respuesta de Axel"».
  - En cambio, la Verifica de S2-16 (S2.md l. 91) mide **por pregunta**: «cada marca vigente cita uno de esos bloqueos y una pregunta sin "Respuesta de Axel"».
- **HECHO 2.** Cuatro bloqueos tienen «Respuesta de Axel» y aun así siguen abiertos en parte:
  - `S2-ventanas`: respuesta en l. 56; «Estado del bloqueo (puerta 2): ABIERTO» en l. 88 (pregunta 4).
  - `S2-flujo-quincenal`: l. 8 y l. 32 (Q2-12).
  - `S5-carrito-comparativa`: l. 14 y l. 38 (Q2-48).
  - `S5-terminos`: l. 22 y l. 33.

  Hay marcas vigentes que los citan: `docs/contratos/motores/S2-interfaz.md` l. 155 y l. 177 (`S2-ventanas` pregunta 4), l. 258 (`S2-flujo-quincenal`, Q2-12), y `S5-interfaz.md` l. 31 y l. 74 (`S5-carrito-comparativa`, Q2-48). S2-09, S2-10, S2-16 y S5-13 **exigen** esas marcas, e INDICE §6 las enumera como vigentes. Si se aplica C1-15 al pie de la letra, `grep -L "Respuesta de Axel"` no lista esos archivos, la Verifica falla y el cuerpo manda quitar la marca.
- **HECHO 3.** C1-15 afirma que `C1-hora-reportes-job`, `C1-reparto-alcance`, `S3-iva-lectura` y `S1-reparto-edicion-gasto` tienen «Estado: ABIERTO». Al contar líneas con «estado:» con `grep -c -i`, `C1-hora-reportes-job.md` da 0 y `C1-reparto-alcance.md` da 0, y `S5-precio-vigente-captura.md` (citado en S5-13) también da 0. Con «la regla general es esa», las marcas que esos dos bloqueos de C1 sostienen tampoco la cumplen.
- **Impacto:** S2-16 y C1-15 dan respuestas opuestas para la misma marca. El revisor de C1 (y el de S2 o S5, si se guía por INDICE §6) tendría que adivinar si esas marcas son defectos. No cambia ninguna cifra ni ningún id.
- **Corrección sugerida:** usar una sola regla, por pregunta, en el cuerpo y la Verifica de C1-15 y en INDICE §6. Una marca es vigente si cita un bloqueo sin «Respuesta de Axel», o una pregunta que el bloqueo declara ABIERTA (en su «Estado» o en «Estado del bloqueo (puerta 2): ABIERTO»). Un bloqueo con la respuesta completa está RESUELTO. La Verifica debe medirlo así: `grep -L "Respuesta de Axel"` **o** `grep -n "Estado del bloqueo (puerta 2): ABIERTO"`, más la pregunta citada. Hay que quitar de C1-15 la afirmación «"Estado: ABIERTO"» sobre los cuatro bloqueos, o pedir al dueño de los bloqueos que agregue esa línea a `C1-hora-reportes-job`, `C1-reparto-alcance` y `S5-precio-vigente-captura`. Esto último está fuera del alcance de criterios. No cambian ids ni severidades.

### B-1 · baja · (sigue de r1) C1-22 cita la salida de `verify_c1.py` con otras palabras
- **Ruta:** `docs/criterios/C1.md` l. 132. **HECHO:** el criterio dice «claves de la puerta 2 con un rango distinto del aprobado: []» y el script imprime «claves de puerta 2 con rango distinto del aprobado: {k47rng}» (`verify_c1.py` l. 555). **Corrección:** copiar el texto exacto.

### B-2 · baja · (sigue de r1, B-3) «Q2-01 caso 1 = A» se puede leer como respuesta de Axel
- **Ruta:** `docs/criterios/S2.md` l. 6 e `INDICE.md` l. 208. **Corrección:** escribir «Q2-01 = A (autoriza la pregunta 4 y la regla PROVISIONAL; los casos 1 y 2 siguen sin respuesta)».

## Estado de los hallazgos de la ronda 1

| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| H-1 (alta) · U4-05 decidía la preselección | **CERRADO** | U4.md l. 51 dice: «Q2-68 = C, que no fija cuál aparece preseleccionada; **la preselección es PROVISIONAL** (bloqueo `…/U4-reparto-preseleccion.md`, abierto, se pregunta en la puerta 4; en uso: "Solo a los gastos nuevos")». La Verifica (l. 52) pide «la preselección marcada PROVISIONAL con cita del bloqueo […] (0 requisitos firmes sobre cuál opción va preseleccionada)». El bloqueo sigue ABIERTO y sin recomendación. `U4-ajustes-b.md` l. 247 y l. 392 marcan PROVISIONAL. Criterio y artefacto ya coinciden. |
| M-1 (media) · C1-15 omitía `S1-reparto-edicion-gasto` | **CERRADO** en lo que se pidió | C1.md l. 104 e INDICE l. 68 ya lo nombran. `C1-datos.md` l. 1044, `reglas-negocio.md` l. 103 y `diccionario-datos.md` l. 2360 lo citan, y el bloqueo dice «Estado: ABIERTO» y no tiene respuesta. La redacción nueva de la regla general abre el **M-1 de esta ronda**. |
| M-2 (media) · S5-13 «solo» cuatro bloqueos | **CERRADO** | S5.md l. 83 separa las marcas del contrato S5 1.5 de las del motor y el oráculo, y nombra `S5-motor-oraculo-reglas` y `S5-precio-vigente-captura`. Mi mapeo da 14 renglones con el primero (`pricing.go`, `compare.go`, `profeco/parse.go`…) y 5 con el segundo (`pricing.go`, `eval.go`, `pre_calculo.py`, `pre_comun.py`). Los dos están sin respuesta. INDICE §6 también los lista. |
| B-1 (baja) | sigue | Ver B-1 arriba. |
| B-2 (baja) · S5-07 y K53 | sigue, toca a P2d | `decisiones.md` llega a K52, K53 tiene 0 entradas y los criterios la citan 6 veces como «se registra en P2d». |
| B-3 (baja) | sigue | Ver B-2 arriba. |
| B-4 (baja) · U2 135° | sigue, fuera del alcance de P2c | No se tocó U2. |

## Notas para el integrador (fuera del alcance de los criterios; no son hallazgos de esta revisión)
1. `backend/internal/precios/ocr/types.go` l. 13–15 todavía dice «PROVISIONAL: la clave `ocr.timeout_seconds` todavía no existe en app_setting». Ya existe: la siembra la 00004 con 30 (Q2-76 = A, K50, FIRME). Bajo S5-06 y S5-13 es una marca vieja.
2. `oraculo/pre_genera_fixtures.py` l. 254 tiene «PROVISIONAL, bloqueo S5-reglas-comparativa pregunta 3». Ese bloqueo está RESUELTO (Q2-59 a Q2-63).
3. El motor S4 no cambió en P2c (`git status` sin cambios en `backend/internal/recetas`). Sigue con marcas PROVISIONAL de `S4-asignacion-mercado` (`shopping.go` l. 243, `oraculo/rec_plan.py` l. 3 y l. 151, `rec_genera_fixtures.py` l. 395) y de `S4-reglas-lista` (`plan.go` l. 66, `shopping.go` l. 113, `rec_genera_fixtures.py` l. 225 y l. 235). Los dos bloqueos están RESUELTOS y S4-15 dice que «no dejan marca». El revisor de S4 en WF-3 lo va a reprobar.
4. `extras.go` l. 89–90 («sin costo capturado (PROVISIONAL, S5-09)») no cita un bloqueo, pero es lo que pide el texto de S5-09. No es un defecto.
5. `S2-interfaz.md` l. 354 (fixture `s2-liga-provisional`) cita `S2-flujo-quincenal` sin decir la pregunta. La Verifica de S2-16 pide la pregunta. Toca al contrato S2.
