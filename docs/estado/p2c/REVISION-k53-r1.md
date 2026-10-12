# Revisión independiente · k53-aplicar · ronda 1 (P2d · K53)

Revisor `revisor-k53-r1` (Opus), 12-oct-2026, madrugada. Evalúa el trabajo de `k53-aplicar` (reporte `docs/estado/p2c/k53-aplicar.md`). Lo contrasta con `docs/estado/puerta2/RESPUESTAS.md` («Puerta 2 bis»), `docs/plan/decisiones.md` §K (K37, K48, K53), los criterios C1 y S5 y `docs/plan/protocolo-agentes.md` §4. No corregí nada: solo leí, corrí verificaciones y escribí este informe y mi reporte. No toqué `docs/design/sistema/`, `docs/design/maquetas/` ni `docs/estado/wfp3/`.

## Veredicto

**NO APROBADO: 1 alto, 2 medios y 4 bajos.**

Todas las verificaciones mecánicas pasan: C1, OpenAPI, criterios y `verify-all`. Lo que falla es de fondo:
- **El alto:** el criterio S5-04 sigue diciendo FIRME la regla de K37 que K53 sustituyó.
- **Medio 1:** hay límites de lote que no salen de ninguna fuente.
- **Medio 2:** la prosa del saldo quincenal (I-6) no cuadra con S2 1.3 §3.11.

## Verificación (comandos y salidas)

| Punto | Comando | Salida | Resultado |
|---|---|---|---|
| (a) cita literal | `python3 -I`: extrae de RESPUESTAS.md el literal «Ignoralo…» (446 caracteres) y lo busca en otros archivos | `decisiones.md` True (1 vez); `bloqueos/S5-terminos.md` True; `bloqueos/S5-agente-externo.md` True | OK |
| (a) K37 y K48 sin borrar | `python3 -I`: compara cada renglón contra `git show HEAD:docs/plan/decisiones.md` | K37: `prefix kept: True`, solo agrega «Sustituida en parte por K53…». K48: `prefix kept: True`, solo agrega «Sustituida por K53…». K46 cambió, pero lo cambió `p2c-remanentes`: su texto lo firma | OK |
| (b) S5-terminos | `git diff docs/estado/bloqueos/S5-terminos.md` | solo inserciones: «Respuesta de Axel (puerta 2 bis…)» con el literal y «**Estado del bloqueo (puerta 2 bis): RESUELTO** por K53» | OK |
| (b) INDICE de bloqueos | `ls … \| grep -v INDICE \| wc -l` y lectura del último «Estado» de cada archivo con `python3 -I` | 85 archivos: 64 RESUELTOS; 6 ABIERTO con «Estado del bloqueo»; 15 abiertos en otro formato (13 de P2c y 2 de la Pista 3), o sea 21 abiertos. La tabla dice 85 / 64 / 21 | OK |
| (b) S5-agente-externo | `cat` | (a) peso, (b) qué escribe, (c) receta en borrador; además (d), (e) y (f) técnicas. «Es PROPUESTA del orquestador… No es decisión de Axel»; «Estado del bloqueo (P2d): ABIERTO» | OK |
| (b) hechos del esquema que cita el bloqueo | `grep` y `awk` en `00001_esquema.sql` | `price.source` CHECK (`ticket_manual`, `ticket_ocr`, `site_assisted`, `profeco`, `manual`), l. 711; `recipe.origin` CHECK (`claude_seed`, `household`), l. 970; `idempotency_key.user_account_id` NOT NULL; `audit_log.user_account_id` admite NULL | OK |
| (c) S5-interfaz 1.6 | `grep -n -i "asistida\|veredicto\|términos"` y `sed -n` de §3, §6 y §13 | §3: `site_assisted` («agente externo» desde K53), con PROVISIONAL (d) y (a) y el HECHO del motor; §13: «La lectura asistida no aplica desde K53». Quedan restos en §15 y §16 (B-1) | OK con un bajo |
| (c) CHECK de `price.source` | `grep -rn "source.*CHECK\|CHECK.*source" backend/db/migrations/` | l. 711: no existe `external_agent`. El contrato usa `site_assisted` y no inventa un valor | OK |
| (d) C1 | `make verify-c1` | 22 `IDÉNTICO`; «los 22 archivos coinciden byte a byte con el generador»; `comprobaciones OK: 33`; `RESULTADO: TODO OK` | OK |
| (d) OpenAPI | `make verify-openapi` | `docs/contratos/openapi.yaml: OK` | OK |
| (d) conteos y seguridad | `uv run --with pyyaml python3 -I` sobre `openapi.yaml` | `version 1.4.0 ops 216 rutas 172 esquemas 262`. `serviceToken` es http bearer. Seguridad global: `cookieAuth`. Las 4 `/api/v1/external-agent/*` llevan `security [{'serviceToken': []}]`, etiqueta «Compras y tiendas», `x-provisional = True` y `x-agente-externo`. Ninguna otra operación usa `serviceToken`. Los 3 POST llevan `IdempotencyKey`; todas responden 401 | OK |
| (d) verify_c1 directo | `verify_c1.py --db` sobre una base migrada temporal | `operaciones sin N-xxx: 5` (las 4 y `GET /openapi.yaml`); `con serviceToken sin la marca: []`; `C1-15 marcas PROVISIONAL en openapi.yaml: 15 … sin bloqueo abierto: []`; 33 OK | OK, coincide con el reporte |
| (e) criterios | `make verify-criterios` | `TOTAL 261 (alta 201 / media 60 / baja 0)`; `TOTAL comprobadas: 178 \| con diferencias: 0`; `INDICE contra conteo propio: 0 diferencias` | OK |
| (e) S5-07 medible | los 5 comandos de su «Verifica», con `bash` | dominios: 0 líneas; interruptor `'false'`: 1; migraciones con `true`: 0; operaciones: 4; `7:  version: "1.4.0"` | OK (el último no dice qué espera: B-2) |
| (f) todo | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | `rc=0`; `RESUMEN TODO OK`; `== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)` | OK |
| (g) alcance | `git status --porcelain`, contra los reportes de los dos agentes; `git status --porcelain docs/design/sistema docs/design/maquetas docs/estado/wfp3 \| wc -l` | Cada archivo modificado es de `k53-aplicar` o de `p2c-remanentes`. 0 cambios en las carpetas vedadas. HEAD avanzó a `be34f91` [P3e], de la Pista 3, y no toca este alcance | OK |

## Hallazgos

### ALTO

**A-1 · S5-04 y S5-10 dicen FIRME la regla de K37 que K53 sustituyó**
- **Archivo:** `docs/criterios/S5.md` (lo editó `k53-aplicar`; versión 1.4).
- **HECHO, S5-04:** dice «`site_assisted` solo es referencia (como PROFECO) hasta que cada cadena tenga veredicto de términos en `docs/estado/S5-terminos.md`». Y cierra: «**Estos pesos y fuentes son FIRMES, sin marca PROVISIONAL**».
- **HECHO, S5-10:** dice «solo es referencia (K37: `site_assisted` y `profeco`)» y nombra la fuente «lectura asistida».
- **Qué dice K53 (§K):** la condición «hasta que cada cadena tenga veredicto» ya no aplica, y el peso de `site_assisted` es PREGUNTA de la puerta 4 (bloqueo `S5-agente-externo`, (a)). Lo mismo dice la nota que `k53-aplicar` puso al final de K37.
- **Contradice a otros archivos:**
  - la fila K37 de `docs/criterios/INDICE.md` («su peso es PROVISIONAL»);
  - `S5-interfaz.md` 1.6 §3 («PROVISIONAL … **sí**, como `manual`»).
- **Por qué es alto:** el criterio que gobierna la verificación de S5 contesta en firme la pregunta (a) de la puerta 4 y contradice §K.
- **Corrección sugerida:**
  - En S5-04 y S5-10, quitar como vigente la condición del veredicto.
  - Marcar el peso de `site_assisted` como PROVISIONAL, citando `S5-agente-externo` pregunta (a) y la PROPUESTA del orquestador («compite como `manual`»).
  - Dejar como HECHO que `-run Fuente` verifica el motor de WF-1, que lo trata como referencia hasta WF-3.
  - En S5-10, «lectura asistida» pasa a «agente externo (`site_assisted`)».
  - Anotarlo en «Cambios de la versión 1.4».
  - Correr `make verify-criterios`.

### MEDIOS

**M-1 · Límites de lote sin fuente ni pregunta**
- **Archivos:** `backend/internal/api/c1gen/gen2_schemas.py` y `gen3_ops.py`; de ahí salen `openapi.yaml`, `C1-datos.md` y los extractos.
- **HECHO:**
  - `ExternalPriceImport.items` dice «De 1 a 500 precios».
  - `ExternalRecipeImport.items` dice «De 1 a 100 recetas».
  - `importExternalCatalog` dice «un lote vacío o de más de 500 renglones responde 422».
- **De dónde no salen:**
  - No los fija RESPUESTAS.md, §K ni C1 1.3. Fuera del generador ningún documento los tiene: `grep -rln "500 renglones\|De 1 a 500\|1 a 100 recetas"` solo encuentra `gen2_schemas.py` y `gen3_ops.py`.
  - El reporte no los declara en «Supuestos» (son 7 y ninguno los menciona).
  - `S5-agente-externo` no tiene una pregunta para ellos. El `x-provisional` de las operaciones cita solo las preguntas (a) a (f).
- **Por qué es medio:** son valores elegidos por el agente, PROVISIONAL sin bloqueo que los cubra.
- **Corrección sugerida:**
  - Declararlos como supuesto.
  - Agregar en `S5-agente-externo` la pregunta (g), tamaño máximo de lote, con opciones e impacto.
  - Citarla en el texto PROVISIONAL de los tres esquemas o de las tres operaciones.
  - Regenerar con `build.py` y correr `make verify-c1`.

**M-2 · I-6: la prosa de `balance_cents` omite el gasto variable registrado**
- **Archivo:** `gen2_schemas.py` l. 124 (`CashFlowFortnights.rows`); de ahí salen `openapi.yaml` y `presupuesto-*.yaml`.
- **HECHO, texto nuevo:** «Saldo de la quincena = ingresos − pagos de crédito − gastos fijos y venideros (S2 1.3 §3.11: income.total_cents − spent_cents − committed_cents)».
- **HECHO, S2 1.3 §3.11:** `spent_cents = variable_spent_cents + fixed_upcoming.paid_cents + credit_payments.household_paid_cents`. Así que el saldo **sí** resta el gasto variable ya registrado.
- **La falla:** la fila de C1 no trae `variable_spent`. Quien siga la prosa calcula `income_cents − credit_payments_cents − fixed_and_upcoming_cents`. Eso no es `balance_cents` en cuanto hay un gasto variable registrado en la quincena. La prosa contradice la fórmula que trae entre paréntesis.
- **Corrección sugerida:**
  - Prosa: «ingresos − gasto variable registrado − pagos de crédito − gastos fijos y venideros».
  - Si se quiere que la fila cuadre por sí sola, PROPUESTA: un campo `variable_spent_cents`. Es de C1/S2 y queda a decisión del integrador.
  - Regenerar.

### BAJOS

- **B-1 · `S5-interfaz.md` 1.6 conserva la lectura asistida sin marca.**
  - §15 trae `settings.assisted_reading_enabled`, `stores[].assisted_reading` y `terms` («veredictos ya leídos de S5-terminos.md») sin «no aplica desde K53; se retira en WF-3».
  - §16 asocia S5-07 con `TestLecturaAsistidaSinVeredicto`, pero el «Verifica» de S5-07 ya no usa `-run`.
- **B-2 · `S5.md`, S5-08 y S5-07.**
  - S5-08 todavía dice «`internal/precios` no hace llamadas de red salvo el módulo de lectura asistida» y verifica «solo en el módulo de lectura asistida». Hoy hay 0 llamadas y K53 retira el módulo.
  - En S5-07, el último comando (`grep -n -m1 "version: \"1.4.0\""`) no dice qué salida espera. Por ejemplo: «una línea, la 7».
- **B-3 · Catálogo de errores.** `c1_plantilla.md` l. 157 sigue listando «lectura asistida sin veredicto» entre los casos de 409. Debe decir «interruptor de lectura asistida en true (no aplica desde K53)».
- **B-4 · Precisión del reporte de `k53-aplicar`.**
  - Dice «20 extractos más su `INDICE.md`». Son 18 extractos y el INDICE: con `openapi.yaml`, `routes.go` y `C1-datos.md` suman los 22 archivos de `build.py --check`.
  - No menciona que `extractos/ajustes-f.yaml` es nuevo (sin seguimiento en git).
  - Dice que «los tres POST son idempotentes por llave natural». Pero `importExternalRecipes` rechaza una receta con nombre existente: sale `rejected`, no `unchanged`. Eso no casa con `ExternalImportResult.unchanged_count` («repetir el lote no duplica nada»). No duplica, pero el resultado no es el mismo.

## Lo que está bien (HECHO)

- K53 trae el literal idéntico. K37 y K48 solo reciben la marca al final; K46 es de `p2c-remanentes`.
- Lo PROVISIONAL va rotulado «PROPUESTA del orquestador, no decisión de Axel». Las preguntas (a), (b) y (c) de la puerta 2 bis quedan abiertas.
- `S5-terminos` está RESUELTO con su línea de estado. El INDICE cuenta 85 / 64 / 21 y coincide con el disco.
- No se inventó `external_agent` en `price.source`. El texto reconoce que el esquema 1.4 no guarda el borrador ni el origen de la receta.
- C1 1.4.0 sale solo del generador:
  - 216 operaciones, 172 rutas y 262 esquemas;
  - `serviceToken` solo en las 4 operaciones marcadas, con `verify_c1.py` reforzado;
  - I-2 (`in_progress`, `achieved`, `overdue`) coincide con S2 1.3 §3.13.
- `verify-all` está en verde.
- El conflicto «compite» contra la fixture `s5-fuentes-compiten` está documentado como HECHO y en el impacto del bloqueo.

## PREGUNTA (para el orquestador, no para Axel)

- El encargo de P2d fija como PROVISIONAL «compite como `manual`». La regla de RESPUESTAS.md para lo no decidido es «PROVISIONAL con la opción que ya está en uso», y la opción en uso en el motor y la fixture es «referencia». Esa elección es del orquestador, no de `k53-aplicar`, y quedó rotulada como PROPUESTA. No la cuento como hallazgo. Pero A-1 se resuelve más limpio si el orquestador confirma cuál de las dos es la PROVISIONAL vigente.
