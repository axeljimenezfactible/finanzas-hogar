# REVISIÓN · criterios de la puerta 2 · ronda 1 de 3 (P2c)

Revisor: `revisor-criterios-r1` (independiente; no escribió nada de lo que revisa y no corrige nada). Fecha: 11-oct-2026, noche (reloj del equipo: 12-oct-2026 00:58).
Objeto: los 9 archivos que cambió `criterios-puerta2` (`docs/criterios/{S1,S2,S4,S5,S7,C1,U3,U4}.md` e `INDICE.md`) y su reporte `docs/estado/p2c/criterios-puerta2.md`.
Contra qué: `docs/criterios/tools/verif_final.sh`, `docs/estado/puerta2/PUERTA-2.md` §4 y §5, `docs/estado/puerta2/RESPUESTAS.md`, `docs/plan/decisiones.md` K46 a K52, los archivos de `docs/estado/bloqueos/` y los reportes de P2c. Todo lo de abajo lo volví a correr yo; no tomé nada de un reporte sin comprobarlo.

## Veredicto: **NO APROBADO** · 1 alto, 2 medios, 4 bajos

## Verificaciones (evidencia propia)

| # | Qué | Comando o script mío | Resultado |
|---|---|---|---|
| a | `make verify-criterios` (11 comprobaciones) | `make verify-criterios` y además `bash docs/criterios/tools/verif_final.sh` completo | rc 0. (1) TOTAL 134 S1–S8+WF4, incompletos: ninguno; (2) 18 archivos, 261 definidos, 0 sin §5; (3) 261 citados, 0 inexistentes; (4) 321 citas, 0 sin fuente; (5) 0 residuos N-xx; (6) 178 cifras canon, 0 diferencias; (7) SQL S1-09: UP OK, 54 tablas, integrity ok; (8) 12 cifras P1–P6 = oráculo; (9) matriz C1: 212 operaciones, 0 N-xxx sin operación; (10) 0 rutas /tmp; (11) TOTAL 261 (201 / 60), INDICE: 0 diferencias. **CUMPLE** |
| b | 261 criterios, 201 / 60 | `datos/tmp/rev-crit-r1/micuenta.py` (regex propia `^**ID (sev)`, sin los scripts de `tools/`) y huella md5 de «id + severidad» de cada archivo tocado contra `git show HEAD:` | 18 archivos, 261, alta 201, media 60, baja 0, 261 ids únicos. Ids y severidades idénticos a HEAD en S1, S2, S4, S5, S7, C1, U3 y U4 (8 de 8 huellas iguales). **CUMPLE** |
| c1 | Valores FIRMES nuevos contra la 00004 | base nueva con `goose up` (v3.28.0) en `datos/tmp/rev-crit-r1/r.db`; `valores.py` extrae de los 19 `.md` toda cita «`clave` = valor» de las 13 claves (10 de la puerta 2, comparador y frescura) y la compara con `app_setting` | Semilla 00004: `reports.daily_send_time` "21:30", `reports.weekly_weekday` 7, `reports.weekly_send_time` "19:00", `reports.monthly_day_of_month` 1, `reports.monthly_send_time` "08:00", `alerts.send_time` "09:00", `alerts.stale_prices_advance_days` 1, `alerts.meal_plan_missing_weekday` 7, `optimizer.date_window_days` 7, `ocr.timeout_seconds` 30 (todas con `user_account_id` NULL); comparador 5, frescura 30 y 90. 53 citas con clave en los criterios: 52 iguales a la semilla y 1 que no es valor inicial (S7-10 Verifica, `reports.monthly_day_of_month` = 31 como dato de prueba del mes corto). **CUMPLE** |
| c2 | Valores contra PUERTA-2 §4 (Q2-72 a Q2-77) y decisión F | misma tabla contra el oráculo escrito a mano desde §4 | semilla = §4 en las 13 claves (0 diferencias). **CUMPLE** |
| c3 | Rangos que declaran C1-22, U4-12, S5-04 | `rangos.sh`: `UPDATE` de prueba en la base temporal contra los disparadores de la 00004 | día de la semana 0 y 8 rechazados, 1 y 7 aceptados (las dos claves); día del mes 0 y 32 rechazados, 31 aceptado; anticipación, ventana y OCR −1 rechazado y 0 aceptado; comparador 0 y 31 rechazados, 1 y 30 aceptados; `price.stale_after_days` 0 rechazado, 1 aceptado; hora "24:00" y "9:00" rechazadas, "00:00" y "23:59" aceptadas. Coincide con lo escrito. `notification.params_json` tras la 00004: 7 filas, 0 con `send_time`, `advance_days`, `weekday` o `day_of_month`; `daily_summary` y `weekly_summary` = `{}` (lo que pide S7-04). **CUMPLE** |
| c4 | Control de C1-22 en `verify_c1.py --db` | base nueva con goose y `uv run … verify_c1.py --db` | 33 OK, `RESULTADO: TODO OK`, «claves de puerta 2 con rango distinto del aprobado: []» (ver B-1 por el texto citado). `extractos/INDICE.md`: C1 v1.3.0, 212 operaciones, 168 rutas, 255 esquemas; `getReportSettings`, `updateReportSettings`, `getAlertSettings`, `updateAlertSettings` en `ajustes-e.yaml`; `optimizer.date_window_days` en `OptimizerSettings.parameters` (`ajustes-d.yaml` l. 656). **CUMPLE** |
| d | S5-07 mide `docs/estado/S5-terminos.md` | lectura de S5-07 l. 58–61 y `grep -c "Veredicto" docs/estado/S5-terminos.md` | La Verifica usa la ruta absoluta `/home/deepsec/finanzas-hogar/docs/estado/S5-terminos.md` (no el reporte `docs/estado/p2c/S5-terminos.md`). Conteo = 3 (l. 13, 20, 27, los tres «no verificado»). El criterio no se marca cumplido y lo dice. Nota: el archivo está sin seguimiento en git (`??`), lo agrega el integrador. **CUMPLE** (ver B-2) |
| e | Ninguna cita a un bloqueo resuelto como si siguiera abierto | `citas_bloq.py`: estado de cada bloqueo por su última línea «Estado del bloqueo (puerta 2)» (o «Respuesta de Axel» si no la tiene) y contexto de ±160 caracteres de cada mención en los 19 `.md` | 0 casos en el texto vigente de los criterios. Los 13 avisos de mi script sobre `S1-parametros-avisos` «abierto» están todos en listas «Cambios de la versión 1.2 …» (historial: S7 l. 12–20, C1 l. 16 y 24, U4 l. 11–18, INDICE §11 l. 175–191), no en criterios. Los bloqueos citados como vigentes (INDICE §6, S2-16, S4-15, S5-13, C1-15) están ABIERTOS en su archivo; los citados como RESUELTOS lo están. **CUMPLE**, con la omisión inversa de H-1, M-1 y M-2 (bloqueos abiertos que no se citan) |
| f | Verificación global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | rc 0 en 21.6 s: «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». **CUMPLE** |
| g | Cada cambio contra la opción recomendada | lectura de PUERTA-2 §2.1–§2.4, §3.1–§3.26, §4, §5 y RESPUESTAS.md frente al diff (175 KB, por archivo) | Q2-11 B, Q2-72 a Q2-77 A, Q2-66 B, Q2-03/05/06/07/08 A, Q2-13 B, Q2-14 A, Q2-15 a Q2-18, Q2-21 a Q2-27, Q2-39 a Q2-43, Q2-65 A, Q2-67 (1), Q2-68 C, Q2-70 (a), Q2-86 A: aplicados como dice PUERTA-2. Sin recomendación (Q2-07 variante, Q2-12, Q2-48, Q2-50 a Q2-52, Q2-56 a Q2-58, Q2-01 caso 2, Q2-85, Q2-87 a Q2-90, Q2-94) quedan PROVISIONAL «en uso» o sin cambio. **Excepción: H-1** |

## Hallazgos

### H-1 · alta · U4-05 da por firme la preselección del diálogo de alcance, que es un bloqueo abierto sin recomendación
- **Ruta:** `docs/criterios/U4.md` l. 51 (cuerpo de U4-05) y l. 52 (Verifica).
- **HECHO:** U4-05 dice «Un cambio de reparto (P-26.12) pregunta antes de guardar el alcance «Solo a los gastos nuevos» (preseleccionado) o «También a los gastos ya capturados» (Q2-68 = C)» y la Verifica exige «P-26.12 trae las dos opciones de alcance con «Solo a los gastos nuevos» preseleccionada». `docs/estado/bloqueos/U4-reparto-preseleccion.md` (ABIERTO, abierto en P2c por `correccion-U4-r1`) dice: «Q2-68 C, literal […] No fija preselección», da tres opciones (A, B, C) y «Este agente no recomienda: solo abre la pregunta»; se pregunta en la puerta 4. `docs/design/U4-ajustes-b.md` 1.5 l. 247 y l. 392 marca la preselección **PROVISIONAL** citando ese bloqueo.
- **Por qué es alto:** el criterio decide una pregunta abierta y sin recomendación, y además se la atribuye a Q2-68 = C, que no la fija. Un revisor de U4 que siga U4-05 trataría la opción A como requisito firme y reprobaría un cambio a B o C que Axel elija en la puerta 4; el artefacto U4 y el criterio quedan en contradicción.
- **Corrección sugerida:** en U4-05 escribir «la preselección es PROVISIONAL (bloqueo `docs/estado/bloqueos/U4-reparto-preseleccion.md`, abierto; en uso: «Solo a los gastos nuevos»)» y en la Verifica «P-26.12 trae las dos opciones de alcance y la preselección marcada PROVISIONAL con cita del bloqueo»; agregar `U4-reparto-preseleccion` (U4-05) a la lista de INDICE §6 y a la lista «Cambios de la versión 1.3» de U4.

### M-1 · media · C1-15 enumera las marcas vigentes de C1 v1.3 y omite `S1-reparto-edicion-gasto`
- **Ruta:** `docs/criterios/C1.md` l. 104 (Verifica de C1-15) e `INDICE.md` l. 68 (§6, lista de bloqueos que sostienen marcas).
- **HECHO:** la Verifica dice «las marcas vigentes en C1 v1.3 citan `C1-hora-reportes-job`, `C1-reparto-alcance` o `S3-iva-lectura` (abiertos)». `docs/contratos/C1-datos.md` l. 1044 tiene «… PROVISIONAL | `S1-reparto-edicion-gasto` (abierto)» (fila «Edición de un gasto anterior después de «solo los nuevos»»), y `docs/estado/bloqueos/S1-reparto-edicion-gasto.md` l. 33 dice «Estado: ABIERTO». La misma marca está en `reglas-negocio.md` l. 103.
- **Impacto:** la lista se lee como cerrada; el revisor de C1 tendría que adivinar si esa marca es un defecto (no lo es: el bloqueo está abierto).
- **Corrección sugerida:** agregar `S1-reparto-edicion-gasto` a la lista de C1-15 y a INDICE §6, o redactar la lista como ejemplos («entre ellas …») y dejar la regla general (bloqueo sin «Respuesta de Axel» o con «Estado del bloqueo: ABIERTO»).

### M-2 · media · S5-13 dice que las marcas de S5 citan «solo» cuatro bloqueos, y el motor y el oráculo citan dos más, abiertos en P2c
- **Ruta:** `docs/criterios/S5.md` l. 83 (S5-13) e `INDICE.md` l. 68 (§6).
- **HECHO:** S5-13: «Las marcas PROVISIONAL vigentes en S5 1.5 citan solo `S5-carrito-comparativa`, `S5-carritos-detalle`, `S5-carritos-granel` y `S5-profeco-filtro` (bloqueos abiertos)». El código que S5-13 evalúa trae marcas PROVISIONAL con otros dos bloqueos ABIERTOS: `backend/internal/precios/pricing.go` l. 215 y `oraculo/pre_comun.py` l. 41 y 58 (`S5-precio-vigente-captura`), `pricing.go` l. 227 y `compare.go` l. 76 y 105 (`S5-motor-oraculo-reglas`, 8 reglas «en uso: el motor»). Los dos bloqueos los abrieron `correccion-motor-S5-r1` y `-r2` en P2c y terminan en «Estado: ABIERTO».
- **Impacto:** «S5 1.5» es el contrato, que no los cita, pero S5-13 es el criterio de calidad del motor; con «solo» el revisor de S5 tendría que adivinar si esas marcas del código son defectos. INDICE §6 tampoco los lista.
- **Corrección sugerida:** en S5-13 decir «las marcas del contrato S5 1.5 citan … y las del motor y el oráculo además `S5-motor-oraculo-reglas` y `S5-precio-vigente-captura` (abiertos, P2c)»; agregarlos a INDICE §6.

### B-1 · baja · C1-22 cita una salida de `verify_c1.py` con otras palabras
- **Ruta:** `docs/criterios/C1.md` l. 132 (Verifica de C1-22).
- **HECHO:** el criterio pide «claves de la puerta 2 con un rango distinto del aprobado: []»; el script imprime «claves de puerta 2 con rango distinto del aprobado: []» (`verify_c1.py` l. 555; mi corrida da `[]`). Un `grep` literal del texto del criterio no encuentra la línea.
- **Corrección sugerida:** copiar el texto exacto de la salida.

### B-2 · baja · S5-07 y la puerta 2 bis: queda para P2d
- **Ruta:** `docs/criterios/S5.md` l. 58–61 e `INDICE.md` l. 68 y §12.
- **HECHO:** RESPUESTAS.md (puerta 2 bis) dice que con K53 «interruptor, veredicto de términos y S5-07 dejan de aplicar» y que «el bloqueo S5-terminos se cierra con esta respuesta», pero también «se registra en P2d». Hoy `decisiones.md` llega a K52 (0 menciones de K53), K48 y `bloqueos/S5-terminos.md` dicen ABIERTO, y el criterio conserva S5-07 con la nota «NO CUMPLE en el fondo» y la PROPUESTA de reformularlo. Es coherente con lo registrado y está declarado; no lo cuento como defecto de P2c. Además, el renglón «Estado al 11-oct-2026 (P2c)» mete un estado de evaluación dentro del texto del criterio.
- **Corrección sugerida:** que P2d, al registrar K53, retire o reformule S5-07 (con la firma de Axel), quite `S5-terminos` de INDICE §6 y mueva la nota de estado a §12 o al reporte; si no, el revisor de S5 reprobará un alto que Axel ya retiró.

### B-3 · baja · «Q2-01 caso 1 = A» se puede leer como respuesta de Axel
- **Ruta:** `docs/criterios/S2.md` l. 6 y l. 127; `INDICE.md` l. 208.
- **HECHO:** Q2-01 = A autoriza escribir la pregunta 4 con sus dos casos y poner la regla PROVISIONAL (PUERTA-2 §2.1); el caso 1 sigue sin respuesta (`bloqueos/INDICE.md` l. 426: «sin respuesta; Q2-01 = A autoriza escribir la pregunta, no la cierra»). El cuerpo de S2-09 y S2-16 lo dice bien (contigüidad PROVISIONAL); solo la etiqueta confunde.
- **Corrección sugerida:** «Q2-01 = A (autoriza la pregunta 4 y la regla PROVISIONAL; casos 1 y 2 sin respuesta)».

### B-4 · baja · U2 sigue llamando «propuesta» al ingreso en 135° (fuera del alcance declarado)
- **Ruta:** `docs/criterios/U2.md` l. 27 y l. 51.
- **HECHO:** U2-02 y U2-08 dicen «ingreso 150° (se admite 135° como propuesta de U2 …)»; Q2-81 = A (K51) ratificó 135°. No contradice (lo admite), pero ya no es propuesta. El reporte declara U2 fuera de su alcance.
- **Corrección sugerida:** en la ronda siguiente de criterios o por el integrador: «ingreso 135° (Q2-81 = A, K51)».

## Estado de los anteriores
Primera ronda de revisión de los criterios en P2c: no hay hallazgos anteriores de este flujo que seguir. Los de las revisiones de criterios de P1 (r1 a r4) quedaron en el historial de §11 y no se reabren.

## Notas (no son hallazgos)
- La ubicación de `ocr.timeout_seconds` en «Datos y respaldos» es propuesta de S1 que la puerta 2 no preguntó; U4-12 lo dice y cae en la regla de K44 («la sección que U4 proponga, con su justificación en SUPUESTOS»).
- Las solicitudes fuera de alcance del reporte (DESIGN.md l. 933, U1 con `params_json`) son correctas y siguen abiertas para sus dueños.
- No apliqué ninguna corrección: soy revisor.
