# Revisión independiente · k53-aplicar · ronda 3 (P2d · K53)

Revisor `revisor-k53-r3` (Opus), 12-oct-2026, madrugada. Evalúo el trabajo de `k53-aplicar` (reporte `docs/estado/p2c/k53-aplicar.md`) con las correcciones `correccion-k53-r1` y `correccion-k53-r2`. Lo contrasto con:
- `docs/estado/puerta2/RESPUESTAS.md`, sección «Puerta 2 bis»;
- `docs/plan/decisiones.md` §K (K37, K48 y K53);
- los criterios C1 y S5;
- `docs/plan/protocolo-agentes.md` §4.

No corregí nada. Solo leí, corrí las verificaciones y escribí este informe y mi reporte. No toqué `docs/design/sistema/`, `docs/design/maquetas/` ni `docs/estado/wfp3/`, ni el alcance de `p2c-remanentes`.

## Veredicto

**NO APROBADO: 0 altos, 1 medio y 7 bajos.**

- El medio de la ronda 2 (§1 «Objetivo» de `S5.md`) está **CORREGIDO**. Todas las verificaciones mecánicas pasan: `verify-c1`, `verify-openapi`, `build.py --check`, `verify-criterios` y `verify-all`.
- **El medio nuevo (M-1):** C1 1.4 todavía describe dos campos con la regla anterior a K53:
  - `Store.assisted_reading_enabled` y `StoreUpdate.assisted_reading_enabled` dicen «Queda en 0 hasta que docs/estado/S5-terminos.md tenga veredicto para ese sitio».
  - `Price.source` describe `site_assisted` como «lectura asistida del sitio».

  Las dos descripciones salen de los comentarios de columna del esquema 1.4, que el generador copia tal cual. Chocan con K53 y con el resto del mismo contrato. Es el mismo tipo de residuo que el M-1 de la ronda 2, ahora en C1.

## Estado del hallazgo de la ronda 2

| Id r2 | Qué pedía | Estado | Evidencia |
|---|---|---|---|
| M-1 (medio) | §1 «Objetivo» de `S5.md` sin «desactivada hasta que exista veredicto escrito por sitio» | **CORREGIDO** | l. 36: «… de modo que reproduzca la comparativa del reporte y no lea sitios de tiendas (K53: la lectura asistida no aplica, el interruptor queda apagado y `docs/estado/S5-terminos.md` es registro histórico)». Está anotado en «Cambios de la versión 1.4» (l. 8). Quedan 3 menciones del texto viejo y ninguna pide esperar un veredicto: l. 8 cita lo que se retiró, l. 64 es el «Texto anterior, como registro» de S5-07 y l. 113 es la cita literal de §E que comprueba `citas.py` |
| B-1 a B-7 (bajos) | — | **Sin cambios.** No venían en el encargo de `correccion-k53-r2` | Se repiten abajo, con la evidencia de hoy |

## Verificación (comandos y salidas)

| Punto | Comando | Salida | Resultado |
|---|---|---|---|
| (a) cita literal | `python3 -I /tmp/revk53r3/a.py`: extrae de RESPUESTAS.md el literal «Ignoralo…» y lo busca | `literal chars 446`; `decisiones.md True 1`; `bloqueos/S5-terminos.md True 1`; `bloqueos/S5-agente-externo.md True 1` | OK |
| (a) K37 y K48 sin borrar | `python3 -I /tmp/revk53r3/a2.py`: compara cada renglón con `git show HEAD:docs/plan/decisiones.md` | K37: `actual empieza con HEAD completo: True`, solo agrega « **Sustituida en parte por K53 (P2d, 12-oct-2026):** …». K48: `True`, solo agrega « **Sustituida por K53 (P2d, 12-oct-2026):** …». Cambian 3 renglones de HEAD (K37, K46 y K48) y se agrega K53. K46 es de `p2c-remanentes` (su reporte, l. 45 y 80) | OK |
| (b) S5-terminos | `git diff --stat` y `grep -n "Estado del bloqueo"` | `9 insertions(+)` en el bloqueo y `2 ++` en `docs/estado/S5-terminos.md`, sin renglones quitados. l. 42: `**Estado del bloqueo (puerta 2 bis): RESUELTO** por K53`. El registro del 11-oct-2026 solo recibe la nota de cierre «No aplica desde K53» | OK |
| (b) INDICE contra disco | `python3 -I /tmp/revk53r3/b2.py`: última línea «Estado…» de cada archivo | `archivos 85 {'RESUELTO': 64, 'ABIERTO': 21, '??': 0}`. El Total del INDICE (l. 27) dice 85 / 64 / 21. La fila WF-1 dice 26 / 21 / 5 y la fila P2d 1 / 0 / 1 | OK |
| (b) S5-agente-externo | `cat` | Trae (a), (b) y (c) de la puerta 2 bis, (d), (e) y (f) técnicas y (g) de la r1. Rótulos: «Es PROPUESTA del orquestador … No es decisión de Axel»; la (e) va como «PROPUESTA de `k53-aplicar`»; la (g) como «Supuesto de `k53-aplicar` … no es decisión de Axel». Cierra con `**Estado del bloqueo (P2d): ABIERTO**` | OK |
| (c) CHECK de `price.source` | `grep -rn "CHECK" backend/db/migrations/ \| grep -i "source\|origin"` | `00001_esquema.sql:711`: `source IN ('ticket_manual', 'ticket_ocr', 'site_assisted', 'profeco', 'manual')`, sin `external_agent`. `:970`: `recipe.origin IN ('claude_seed', 'household')` | OK |
| (c) S5-interfaz 1.6 | `grep -n -i "asistida\|assisted\|veredicto\|terms\|agente externo\|external"` y `sed -n '68,78p;114,122p;211,240p'` | §3, l. 74: `site_assisted` («agente externo»), con la fuente PROVISIONAL (pregunta d) y el peso PROVISIONAL (pregunta a), los dos con el HECHO del motor. Coincide con el CHECK. §13: «La lectura asistida no aplica desde K53». Quedan restos sin marca en §15 y §16 (B-1) | OK, con un bajo |
| (d) C1 | `make verify-c1` | 22 `IDÉNTICO`; «los 22 archivos coinciden byte a byte con el generador»; `comprobaciones OK: 33`; `RESULTADO: TODO OK`; rc 0 | OK |
| (d) OpenAPI | `make verify-openapi` | `docs/contratos/openapi.yaml: OK`; rc 0 | OK |
| (d) generador | `python3 -I backend/internal/api/c1gen/build.py --check` | «los 22 archivos coinciden byte a byte con el generador»; rc 0 | OK |
| (d) conteos y seguridad | `uv run --no-project --with pyyaml python3 -I /tmp/revk53r3/d.py` | `version 1.4.0 ops 216 rutas 172 esquemas 262`. `serviceToken` = (`http`, `bearer`); seguridad global `cookieAuth`. `ops con serviceToken`: solo las 4 `/api/v1/external-agent/*`. Las 4 llevan `sec [{'serviceToken': []}]`, etiqueta «Compras y tiendas», `x-provisional` y `x-agente-externo`. Los 3 POST llevan `IdempotencyKey`. `SavingsGoal.status` = `['in_progress', 'achieved', 'overdue']`. `CashFlowFortnights` sin `planned_savings_cents` obligatorio | OK, coincide con el reporte |
| (d) citas de lo PROVISIONAL | `… /tmp/revk53r3/d3.py` | Las 4 operaciones citan `S5-agente-externo`, pregunta e. Además, el catálogo cita g; los precios, a y d; las recetas, c y f. Los 3 esquemas de lote van PROVISIONAL (catálogo b y g, precios a, d y g, recetas c, f y g). `serviceToken` lleva `x-provisional: true` y cita la pregunta e | OK |
| (d) verify_c1 | `uv run --no-project --with pyyaml --with jsonschema python3 -I …/verify_c1.py` | rc 0, 33 `OK`. `C1-02 … operaciones sin N-xxx: 5` (las 4 y `GET /openapi.yaml`), `con serviceToken sin la marca: []`; `C1-13 … info.version: 1.4.0`; `C1-15 marcas PROVISIONAL en openapi.yaml: 18 … sin bloqueo abierto: []` | OK |
| (d) lectura asistida en C1 | `… /tmp/revk53r3/d2.py` y `grep -n -i "veredicto\|asistida" docs/contratos/openapi.yaml` | `updateStore` y `updatePriceSettings`: 409 `assisted_reading_not_allowed` «porque ningún veredicto la habilita». `Store.assisted_reading_verdict`, `PriceSources`, `PriceSettings` y `PriceSettingsUpdate`: «no aplica desde K53». **Pero** l. 14076 (`Store`) y l. 14212 (`StoreUpdate`): «Queda en 0 hasta que docs/estado/S5-terminos.md tenga veredicto para ese sitio»; l. 14747 (`Price.source`): «… lectura asistida del sitio …» | **M-1** |
| (e) criterios | `make verify-criterios GOOSE=…` | `TOTAL 261 (alta 201 / media 60 / baja 0) \| archivos: 18`; `citas verificadas: 321 \| sin fuente: 0`; `TOTAL comprobadas: 178 \| con diferencias: 0`; `cifras comparadas: 12 \| diferencias: 0`; `INDICE contra conteo propio: 0 diferencias`; rc 0 | OK |
| (e) S5-07 medible | Sus 5 comandos, con `bash -c` | dominios en `.go` sin pruebas: `0 lineas`; interruptor `'false'`: `1`; migraciones con `true`: `0 archivos`; operaciones: `4`; `7:  version: "1.4.0"`. Cada comando tiene una salida esperada concreta, salvo el último (B-2) | OK, con un bajo |
| (f) todo | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` (log de 147 renglones) | `rc=0`. C1 `TODO OK`, S1 `RESUMEN TODO OK`, U1 `comprobaciones OK: 63`. Oráculos: s2 `32 coinciden`, s3 27 fixtures, s4 `26 coinciden` y s5 `43 fixtures, 43 coinciden`, todos con 0 diferencias. `pre_verify.py: 15 comprobaciones OK`; `verificar_c2.py: RESULTADO: OK`. Al final: `== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)` | OK |
| (g) alcance | `git status --porcelain -- docs/design/sistema docs/design/maquetas docs/estado/wfp3 \| wc -l` y cada archivo modificado contra los reportes de `k53-aplicar`, las dos correcciones y `p2c-remanentes` | `0` en las carpetas vedadas. Cada archivo modificado es de uno de los dos agentes. `p2c-remanentes` declara que no tocó S5-interfaz, C1/c1gen, criterios, plan ni README, y en `decisiones.md` solo tocó K46. HEAD avanzó a `be34f91` con commits de la Pista 3, ajenos a este alcance | OK |

## Hallazgos

### MEDIO

**M-1 · C1 1.4 conserva en tres campos la regla de la lectura asistida que K53 retiró**

- **HECHO.** `docs/contratos/openapi.yaml`, generado por `build.py`:
  - l. 14076, `Store.assisted_reading_enabled`, y l. 14212, `StoreUpdate.assisted_reading_enabled`: «Lectura asistida del sitio de la cadena. **Queda en 0 hasta que docs/estado/S5-terminos.md tenga veredicto para ese sitio**».
  - l. 14747, `Price.source` (enum con `site_assisted`): «Fuente: ticket capturado a mano, ticket por OCR ya confirmado, **lectura asistida del sitio**, PROFECO u otra captura manual».
  - Las copias están en los extractos `ajustes-c.yaml` (l. 750 y 886), `compras-tiendas-a.yaml` (l. 1277) y `compras-tiendas-b.yaml` (l. 971).
  - Origen: los comentarios de columna de `backend/db/migrations/00001_esquema.sql` l. 326 y l. 711. `gen1_base.col_schema` los copia tal cual. HEAD ya traía las dos descripciones y P2d no las cambió.
- **Con qué choca:**
  - K53 (1): «dejan de aplicar el interruptor … que se quedan en false, el veredicto de términos …». Axel: «por eso la regla no funciona».
  - El mismo contrato:
    - `updateStore` (l. 4855): «true responde siempre 409 `assisted_reading_not_allowed`, porque ningún veredicto la habilita».
    - `Store.assisted_reading_verdict` (l. 14122): «ningún veredicto habilita la lectura asistida».
    - `importExternalPrices`: «source = site_assisted, que la documentación llama «agente externo»».
  - `C1-datos.md` l. 48 y el reporte de `k53-aplicar` dicen que `Store` ya describe la lectura asistida como «no aplica desde K53». Es cierto del campo calculado `assisted_reading_verdict`, no de `assisted_reading_enabled`, y `StoreUpdate` no lleva ninguna nota.
  - `S5-interfaz.md` 1.6 §3 (l. 74): `site_assisted` = «agente externo».
- **Por qué importa.** WF-3 construye la API y el frontend con C1. El cuerpo de `PATCH /stores/{id}` (`StoreUpdate`) le dice al implementador que el campo se activa cuando exista un veredicto, que es justo la vía «veredicto → activar» que K53 cerró y que la ronda 2 ya corrigió en `S5.md` §1. Y la fuente que C1 1.4 usa para el agente externo sigue rotulada como lectura del sitio.
- **Por qué es medio y no alto:**
  - No decide ninguna pregunta de la puerta 4.
  - No rompe ninguna verificación.
  - En la práctica pide lo mismo que K53: interruptor apagado y 409.
  - K53 dice que P2d no cambia el esquema de S1, así que el comentario de la columna puede quedarse. Pero el generador sí es del alcance de `k53-aplicar`, y el encargo pedía que `Store` dijera «no aplica desde K53».
- **Corrección sugerida (solo por el generador):**
  1. En `backend/internal/api/c1gen/gen2_schemas.py`, después de crear `S['Store']` y `S['StoreUpdate']`, sustituir la descripción de `assisted_reading_enabled` por algo como «Lectura asistida del sitio de la cadena. No aplica desde K53: la aplicación no lee sitios de tiendas; queda en false y true responde 409 `assisted_reading_not_allowed`. La columna sigue en el esquema 1.4 hasta que S1 la retire.»
  2. En `Price.source`, cambiar la descripción por algo como «… `site_assisted` = agente externo desde K53 (PROVISIONAL, bloqueo `S5-agente-externo`, pregunta d) …». Es la única copia en `openapi.yaml` (`grep` da 1 renglón).
  3. Regenerar con `build.py` y correr `make verify-c1 verify-openapi` y `build.py --check`.
  4. Anotarlo en los cambios de C1 v1.4 (`c1_plantilla.md`) y agregar a los pendientes de S1 (WF-3) los comentarios de `00001_esquema.sql` l. 326, l. 711 y l. 1399 («Además cada cadena necesita su veredicto de términos (decisión E)»).

### BAJOS (de la ronda 2, sin cambios)

- **B-1 · `S5-interfaz.md` 1.6 §15 y §16 sin la marca de K53.** Siguen igual:
  - l. 283, `settings.assisted_reading_enabled`;
  - l. 285, `stores[].assisted_reading`;
  - l. 291, `terms` («veredictos ya leídos de `docs/estado/S5-terminos.md`»);
  - l. 311, `assisted_can_enable` con `no_terms_verdict`;
  - l. 338, S5-07 → `TestLecturaAsistidaSinVeredicto`.

  Hay restos parecidos en l. 271 (fixture `s5-lectura-asistida`, «interruptor sin veredicto rechazado») y l. 346 (`ErrNoTermsVerdict`). **Corrección:** marcarlos como HECHO del harness de WF-1 que se retira en WF-3.
- **B-2 · `S5.md`, S5-08 y S5-07.**
  - S5-08 (l. 69) sigue diciendo «salvo el módulo de lectura asistida».
  - El último comando de S5-07 no dice qué salida espera.
  - En zsh, `--include=*.go` va sin comillas.

  **Corrección:** S5-08 «sin llamadas de red (K53)»; agregar «= `7:  version: "1.4.0"`»; escribir `--include='*.go'`.
- **B-3 · `c1_plantilla.md` → `C1-datos.md` l. 159.** El catálogo de 409 sigue diciendo «lectura asistida sin veredicto». Va junto con M-1.
- **B-4 · Precisión del reporte de `k53-aplicar`.**
  - Dice «20 extractos»: son 18 más el INDICE.
  - No dice que `ajustes-f.yaml` es nuevo.
  - Llama idempotentes por llave natural a las recetas, aunque una receta repetida sale `rejected`.
  - Afirma que `Store` dice «no aplica desde K53»; ver M-1.

  **Corrección:** una errata en el reporte del integrador.
- **B-5 · La pregunta (g) no está en los resúmenes.** La fila de `S5-agente-externo` en `docs/estado/bloqueos/INDICE.md` (l. 525) y K53 en `decisiones.md` enumeran solo (a) a (f). **Corrección:** que el integrador agregue la (g).
- **B-6 · `docs/criterios/INDICE.md` no refleja A-1 de la r1.** La fila K53 (l. 170) dice `S5.md` (S5-07) y la l. 3 dice «K53 aplicada a S5-07 y a C1 v1.4». En cambio, `S5.md` §5 asocia K53 a S5-04, S5-07 y S5-10. **Corrección:** listar «S5-04, S5-07, S5-10».
- **B-7 · Resúmenes de §E sin la nota de K53.** `README.md` l. 29 y `docs/criterios/INDICE.md` l. 105 dicen «lectura asistida … desactivada hasta (tener) veredicto». **Corrección:** agregar «(sustituido por K53)».

## Lo que está bien (HECHO)

- **K53:**
  - Trae el literal de Axel idéntico (446 caracteres) en `decisiones.md` y en los dos bloqueos.
  - K37 y K48 conservan su texto completo y solo reciben la marca final.
  - Lo PROVISIONAL va como PROPUESTA del orquestador, no como decisión de Axel.
- **Bloqueos:**
  - `S5-terminos` está RESUELTO con su línea de estado.
  - El INDICE (85 / 64 / 21) coincide con el disco.
  - `S5-agente-externo` está ABIERTO, con (a) a (g) y opciones e impacto.
- **S5-interfaz 1.6:**
  - Usa `site_assisted`, coherente con el CHECK cerrado.
  - No inventa `external_agent` como valor guardado.
  - Reconoce como HECHO que el motor y la fixture van detrás del contrato hasta WF-3.
- **C1 1.4.0:**
  - Sale solo del generador (22 archivos idénticos).
  - Tiene 216 operaciones, 172 rutas y 262 esquemas.
  - `serviceToken` está solo en las 4 operaciones del agente externo, con `x-provisional` y citas a las preguntas.
  - Los 3 POST llevan `Idempotency-Key`.
  - I-2 e I-6 cuadran.
- **Criterios:**
  - S5-07 es medible y cada comando da hoy lo esperado.
  - §1 de `S5.md` ya está alineado.
  - Siguen 261 criterios (201 alta y 60 media).
  - `verify-all` está en verde y las carpetas vedadas no tienen cambios.

## PREGUNTA (para el orquestador; sin cambios desde la r1)

- Lo PROVISIONAL de (a) es «compite como `manual`», la PROPUESTA del encargo de P2d. El motor y la fixture usan «referencia».
- RESPUESTAS.md no fija cuál de las dos rige mientras Axel decide. La regla «la opción en uso» es de las 28 preguntas sin recomendación, no de K53.
- No lo cuento como hallazgo: está rotulado como PROPUESTA y documentado como HECHO en `S5.md`, `S5-interfaz` y el bloqueo.
