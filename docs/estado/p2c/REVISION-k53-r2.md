# Revisión independiente · k53-aplicar · ronda 2 (P2d · K53)

Revisor `revisor-k53-r2` (Opus), 12-oct-2026, madrugada. Evalúo el trabajo de `k53-aplicar` (reporte `docs/estado/p2c/k53-aplicar.md`) con la corrección `correccion-k53-r1` (reporte `docs/estado/p2c/correccion-k53-r1.md`). Lo contrasto con:
- `docs/estado/puerta2/RESPUESTAS.md`, sección «Puerta 2 bis»;
- `docs/plan/decisiones.md` §K (K37, K48 y K53);
- los criterios C1 y S5;
- `docs/plan/protocolo-agentes.md` §4.

No corregí nada. Solo leí, corrí las verificaciones y escribí este informe y mi reporte. No toqué `docs/design/sistema/`, `docs/design/maquetas/` ni `docs/estado/wfp3/`.

## Veredicto

**NO APROBADO: 0 altos, 1 medio y 7 bajos.**

- Los tres hallazgos de la ronda 1 (A-1, M-1 y M-2) están corregidos donde los señalé, y todas las verificaciones mecánicas pasan.
- **El medio:** la sección §1 «Objetivo» de `docs/criterios/S5.md` todavía pide mantener la lectura asistida «desactivada hasta que exista veredicto escrito por sitio». K53 retiró esa condición, y la corrección de A-1 la quitó de S5-04, S5-10, §3, §4 y §5, pero no de §1.

## Estado de los hallazgos de la ronda 1

| Id r1 | Qué pedía | Estado | Evidencia |
|---|---|---|---|
| A-1 (alto) | S5-04 y S5-10 sin la regla de K37 como FIRME; peso de `site_assisted` PROVISIONAL citando `S5-agente-externo` (a) | **CORREGIDO** en S5-04, S5-10, §3, §4 y §5; queda un residuo en §1 (M-1 de esta ronda) | S5-04, l. 51: «desde K53 ya no aplica la condición… PROVISIONAL (bloqueo `S5-agente-externo`, pregunta a) con la PROPUESTA del orquestador, que no es decisión de Axel». «FIRMES» solo cubre `manual`, `ticket_manual`, `ticket_ocr` y PROFECO. S5-10, l. 76 y 77: «agente externo (`site_assisted`)», PROVISIONAL, con el HECHO del motor |
| M-1 (medio) | Límites de lote declarados y cubiertos por una pregunta | **CORREGIDO** (queda un bajo, B-5) | `S5-agente-externo.md` trae la pregunta (g) con HECHO, opciones (g1) a (g3) e impacto, y en «Mientras tanto» la rotula «Supuesto de `k53-aplicar`… no es decisión de Axel». En `openapi.yaml`, los `items` de los 3 lotes y el 422 de `importExternalCatalog` citan la pregunta g |
| M-2 (medio) | La prosa de `balance_cents` cuadra con S2 1.3 §3.11 | **CORREGIDO** | Prosa nueva: «ingresos − gasto variable registrado − pagos de crédito − gastos fijos y venideros… El gasto variable registrado (variable_spent_cents de S2) no viene en esta fila». Coincide con S2-interfaz l. 257 y 260 a 262 |
| B-1 a B-4 (bajos) | — | **Sin cambios.** No venían en el encargo de la corrección y se repiten abajo como B-1 a B-4 | — |

## Verificación (comandos y salidas)

| Punto | Comando | Salida | Resultado |
|---|---|---|---|
| (a) cita literal | `python3 -I /tmp/revk53r2/a.py`: extrae de RESPUESTAS.md el literal «Ignoralo…» y lo busca | `literal chars 446`; `decisiones.md` True (1 vez); `bloqueos/S5-terminos.md` True; `bloqueos/S5-agente-externo.md` True | OK |
| (a) K37 y K48 sin borrar | `python3 -I /tmp/revk53r2/a2.py`: compara cada renglón con `git show HEAD:docs/plan/decisiones.md` | K37: `ahora empieza con HEAD completo: True`, solo agrega « **Sustituida en parte por K53…**». K48: `True`, solo agrega « **Sustituida por K53…**». `renglones quitados/reescritos respecto a HEAD: 3` (K37, K46 y K48) y K53 nuevo; K46 es de `p2c-remanentes` | OK |
| (b) S5-terminos | `git diff --stat` y `grep -n "Estado del bloqueo"` | `9 insertions(+)`, 0 renglones quitados; l. 42: `**Estado del bloqueo (puerta 2 bis): RESUELTO** por K53` | OK |
| (b) INDICE contra disco | `python3 -I /tmp/revk53r2/b.py`: última línea de estado de cada archivo | `archivos 85 resueltos 64 abiertos 21`. El Total del INDICE dice 85 / 64 / 21 | OK |
| (b) S5-agente-externo | `cat` | (a), (b) y (c) de la puerta 2 bis; (d), (e) y (f) técnicas; (g) de la r1. «Es PROPUESTA del orquestador… No es decisión de Axel»; la (e) va como «PROPUESTA de `k53-aplicar`»; `**Estado del bloqueo (P2d): ABIERTO**` | OK |
| (c) CHECK de `price.source` | `grep -rn "CHECK" backend/db/migrations/ \| grep -i "source\|origin"` | `00001_esquema.sql:711` lista `ticket_manual`, `ticket_ocr`, `site_assisted`, `profeco` y `manual`, sin `external_agent`; `:970`: `recipe.origin` lista `claude_seed` y `household` | OK |
| (c) S5-interfaz 1.6 | `grep -n -i "asistida\|assisted\|veredicto\|términos\|terms"` y `sed -n '68,76p;211,240p'` | §3, l. 74: `site_assisted` («agente externo»), con la fuente PROVISIONAL (d) y el peso PROVISIONAL (a), ambos con el HECHO del motor; coincide con el CHECK. §13: «La lectura asistida no aplica desde K53». Quedan restos sin marca en §15 (l. 283, 285, 291 y 311) y §16 (l. 338): B-1 | OK con un bajo |
| (d) C1 | `make verify-c1` | 22 `IDÉNTICO`; «los 22 archivos coinciden byte a byte con el generador»; `comprobaciones OK: 33`; `RESULTADO: TODO OK`; rc 0 | OK |
| (d) OpenAPI | `make verify-openapi` | `docs/contratos/openapi.yaml: OK`; rc 0 | OK |
| (d) generador | `python3 -I backend/internal/api/c1gen/build.py --check` | «los 22 archivos coinciden byte a byte con el generador» (incluye `extractos/ajustes-f.yaml`); rc 0 | OK |
| (d) conteos y seguridad | `uv run --with pyyaml python3 -I /tmp/revk53r2/d.py` | `version 1.4.0 ops 216 rutas 172 esquemas 262`. `serviceToken` = (`http`, `bearer`); seguridad global `cookieAuth`. Las 4 `/api/v1/external-agent/*`: `sec [{'serviceToken': []}]`, etiqueta «Compras y tiendas», `x-prov True`, `x-agente True`; los 3 POST con `IdempotencyKey`. `ops con serviceToken: 4`. `SavingsGoal.status`: `['in_progress', 'achieved', 'overdue']`. La fila de `CashFlowFortnights` no exige `planned_savings_cents` | OK, coincide con los reportes |
| (d) textos PROVISIONAL | `uv run --with pyyaml python3 -I /tmp/revk53r2/d2.py` | Cada operación cita la pregunta e. Precios cita a y d; recetas, c y f; el lote de catálogo, g; los `items` de los tres lotes, g | OK |
| (d) verify_c1 | `python3 -I …/verify_c1.py` (sin la base) | `C1-02 … operaciones sin N-xxx: 5` (las 4 y `GET /openapi.yaml`); `con serviceToken sin la marca: []`; `C1-13 … info.version: 1.4.0`; `C1-15 marcas PROVISIONAL en openapi.yaml: 18 … sin bloqueo abierto: []`. Mi corrida directa da `FALLA C1-14 jsonschema no disponible` porque falta el módulo en mi entorno; `make verify-c1` lo carga con `--with jsonschema` y pasa | OK |
| (e) criterios | `make verify-criterios` | `TOTAL 261 (alta 201 / media 60 / baja 0) \| archivos: 18`; `citas verificadas: 321 \| sin fuente: 0`; `TOTAL comprobadas: 178 \| con diferencias: 0`; `INDICE contra conteo propio: 0 diferencias`; rc 0 | OK |
| (e) S5-07 medible | Sus 5 comandos, con `bash -c` | dominios: `lineas=0`; interruptor `'false'`: 1; migraciones con `true`: 0; operaciones: 4; `7:  version: "1.4.0"`. En zsh, el primer comando falla con `no matches found: --include=*.go` y aun así sale «sin resultados» (B-2) | OK en bash, con un bajo |
| (f) todo | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | `rc=0`; `RESUMEN TODO OK`; s5: `43 fixtures, 43 coinciden`; `== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)` | OK |
| (g) alcance | `git status --porcelain -- docs/design/sistema docs/design/maquetas docs/estado/wfp3 \| wc -l`, y los archivos modificados contra los reportes de `k53-aplicar`, `correccion-k53-r1` y `p2c-remanentes` | `vedadas: 0`. Cada archivo modificado es de uno de los dos agentes; `p2c-remanentes` dice que no tocó c1gen, S5-interfaz, criterios, plan ni README. HEAD avanzó a `be34f91` con commits de la Pista 3, ajenos a este alcance | OK |

## Hallazgos

### MEDIO

**M-1 · `docs/criterios/S5.md` §1 «Objetivo» conserva la condición del veredicto que K53 retiró**

- **HECHO.** La l. 35 (§1, sin cambios desde HEAD; el diff de `S5.md` no tiene ningún bloque en §1) pide: «… los documentos `docs/estado/S5-terminos.md` y `docs/estado/S5-cobertura-profeco.md`, de modo que reproduzca la comparativa del reporte y **mantenga la lectura asistida de sitios desactivada hasta que exista veredicto escrito por sitio**».
- **Con qué choca:**
  - K53 (1): «Se retira del plan la "lectura asistida": dejan de aplicar el interruptor …, el veredicto de términos …».
  - Las palabras de Axel: «por eso la regla no funciona».
  - El mismo archivo: S5-07 (l. 63: «No aplica desde K53»), §3 (l. 95: «Cualquier lectura de sitios de tiendas desde la aplicación (K53) … no aceptable») y S5-04 («ya no aplica la condición … hasta que cada cadena tenga veredicto»).
- **Por qué importa.** El encabezado de `S5.md` dice que el revisor de WF-1 y WF-3 evalúa **solo** contra este archivo, y §1 todavía deja abierta la vía «veredicto → activar».
- **Por qué es medio y no alto:**
  - No decide ninguna pregunta de la puerta 4.
  - No rompe ninguna verificación.
  - En la práctica pide lo mismo que K53, el interruptor apagado.
  - Es un residuo de A-1: la corrección revisó §3, §4 y §5 «por el mismo defecto» y no revisó §1.
- **Corrección sugerida:**
  - En §1, cambiar la cláusula por algo como: «… y no lea sitios de tiendas (K53: la lectura asistida no aplica; el interruptor queda apagado; `docs/estado/S5-terminos.md` es registro histórico)».
  - Anotarlo en «Cambios de la versión 1.4».
  - Correr `make verify-criterios`.

### BAJOS

- **B-1 (de la r1, sin cambios) · `S5-interfaz.md` 1.6 §15 y §16.** Estos renglones siguen sin la marca «no aplica desde K53; se retira en WF-3»:
  - l. 283: `settings.assisted_reading_enabled`;
  - l. 285: `stores[].assisted_reading`;
  - l. 291: `terms` («veredictos ya leídos de `docs/estado/S5-terminos.md`»);
  - l. 311: `assisted_can_enable` con `no_terms_verdict`;
  - l. 338: S5-07 → `TestLecturaAsistidaSinVeredicto`, aunque el «Verifica» de S5-07 ya no corre esa prueba.

  **Corrección:** marcarlos como HECHO del harness de WF-1 que se retira en WF-3.
- **B-2 (de la r1, ampliado) · `S5.md`, S5-08 y S5-07.**
  - S5-08 (l. 68 y 69) sigue diciendo «salvo el módulo de lectura asistida». HECHO: el grep de red en `internal/precios` da 0 líneas.
  - El último comando de S5-07 no dice qué salida espera.
  - Nuevo: el primer grep de S5-07 (`--include=*.go` sin comillas) en zsh falla con `no matches found` y deja la tubería vacía, o sea un «sin resultados» falso. En bash sí funciona. El patrón también aparece en otros criterios (S2, S3, S4, S6 y S7), así que es convención previa.

  **Corrección:** `--include='*.go'`; «= `7:  version: "1.4.0"`»; S5-08 «sin llamadas de red (K53)».
- **B-3 (de la r1, sin cambios) · `c1_plantilla.md` l. 159.** El catálogo de 409 sigue diciendo «lectura asistida sin veredicto». **Corrección:** «interruptor de lectura asistida en true (no aplica desde K53)».
- **B-4 (de la r1, sin cambios) · Precisión del reporte de `k53-aplicar`.**
  - Dice «20 extractos»: son 18 más el INDICE.
  - No menciona que `ajustes-f.yaml` es nuevo.
  - Llama idempotentes por llave natural a las recetas, aunque una receta repetida sale `rejected`, no `unchanged`.
  - En la l. 71 conserva la prosa vieja del saldo.

  **Corrección:** una nota de errata en el reporte del integrador; el reporte de ronda no se reescribe.
- **B-5 (nuevo) · La pregunta (g) no se propagó a los resúmenes.**
  - La fila de `S5-agente-externo` en `docs/estado/bloqueos/INDICE.md` (l. 525) enumera (a) a (f). La fila es de `p2c-remanentes`, y el INDICE avisa que «las opciones … completos están en el archivo».
  - K53 en `decisiones.md` (l. 157) también enumera solo (a) a (f).

  **Corrección:** agregar «(g) tamaño máximo de un lote (supuesto de `k53-aplicar`)» en la fila del INDICE y en K53. Que lo haga el integrador.
- **B-6 (nuevo) · `docs/criterios/INDICE.md` no refleja la corrección de A-1.**
  - La fila K53 (l. 170) lista solo `S5.md` (S5-07).
  - La versión 1.4 (l. 3) dice «K53 aplicada a S5-07 y a C1 v1.4».
  - En cambio, `S5.md` §5 (l. 115) asocia K53 a S5-04, S5-07 y S5-10.

  **Corrección:** «`S5.md` (S5-04, S5-07, S5-10)» y la misma mención en l. 3.
- **B-7 (nuevo) · Resúmenes de §E sin la nota de K53.**
  - `README.md` l. 29 dice «lectura asistida de sitios desactivada hasta tener veredicto de términos», y la l. 5 del mismo README dice que se retiró.
  - `docs/criterios/INDICE.md` l. 105 dice lo mismo («lectura asistida desactivada hasta veredicto»).
  - Las dos filas resumen §E, que K53 decidió no editar.

  **Corrección:** agregar «(sustituido por K53)» en las dos filas.

## Lo que está bien (HECHO)

- **K53:** trae el literal de Axel idéntico (446 caracteres) en `decisiones.md` y en los dos bloqueos. K37 y K48 conservan su texto completo y solo reciben la marca final.
- **Bloqueos:** `S5-terminos` está RESUELTO con su línea de estado, y el INDICE (85 / 64 / 21) coincide con el disco. `S5-agente-externo` está ABIERTO con (a) a (g); lo PROVISIONAL va como PROPUESTA del orquestador y la (e) como PROPUESTA de `k53-aplicar`, nunca como decisión de Axel.
- **S5-interfaz 1.6:** usa `site_assisted`, coherente con el CHECK cerrado. No inventa `external_agent` y reconoce que el esquema 1.4 no guarda el borrador ni el origen de la receta.
- **C1 1.4.0:**
  - sale solo del generador, con 216 operaciones, 172 rutas y 262 esquemas;
  - `serviceToken` está solo en las 4 operaciones marcadas;
  - lleva `x-provisional` en las 4;
  - I-2 e I-6 cuadran con S2 1.3;
  - los límites de lote van PROVISIONAL con la pregunta g.
- **Criterios y verificación:** siguen 261 criterios (201 alta y 60 media), `verify-all` está en verde y no hay cambios en las carpetas vedadas.

## PREGUNTA (para el orquestador; sin cambios desde la r1)

- Lo PROVISIONAL de (a) es «compite como `manual`», la PROPUESTA del encargo de P2d. El motor y la fixture usan «referencia».
- RESPUESTAS.md no fija cuál de las dos rige para K53. La regla «la opción en uso» aplica a las 28 preguntas sin recomendación, no a K53.
- No lo cuento como hallazgo: está rotulado como PROPUESTA y documentado como HECHO en S5.md, S5-interfaz y el bloqueo.
