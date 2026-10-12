# Corrección k53 · ronda 1 (P2d · K53)

Agente `correccion-k53-r1` (Opus), 12-oct-2026, madrugada. Corrige los tres hallazgos que el orquestador me encargó de `docs/estado/p2c/REVISION-k53-r1.md`: A-1, M-1 y M-2. Tengo el mismo alcance que `k53-aplicar`. Reglas:
- No decidí ninguna pregunta de la puerta 4.
- No toqué los bajos B-1 a B-4, que no venían en el encargo.
- No hice `git add` ni commit. No usé Telegram, navegador ni WebFetch.
- No toqué `docs/design/sistema/`, `docs/design/maquetas/` ni `docs/estado/wfp3/`.

## Hecho

### A-1 (alto) · `docs/criterios/S5.md`

Sigue en la versión 1.4, con los mismos 13 criterios (10 de severidad alta y 3 media).

**S5-04**
- Antes decía «`site_assisted` solo es referencia… hasta que cada cadena tenga veredicto de términos» y lo cerraba con «Estos pesos y fuentes son FIRMES». Ese texto ya no está.
- Siguen FIRMES `manual`, `ticket_manual` y `ticket_ocr`, que compiten, y PROFECO, que nunca compite.
- `site_assisted`, la fuente del agente externo:
  - Desde K53 ya no aplica la condición del veredicto de K37.
  - Su peso es PREGUNTA de la puerta 4 y queda **PROVISIONAL** (bloqueo `S5-agente-externo`, pregunta a), con la PROPUESTA del orquestador: compite como `manual`. No es decisión de Axel.
  - HECHO: el motor de WF-1 y la fixture `s5-fuentes-compiten` todavía lo tratan como referencia. Se alinean en WF-3.
- «Verifica» con `-run Fuente`:
  - PROFECO aparece solo como referencia.
  - HECHO: la prueba verifica hoy el motor de WF-1, que trata `site_assisted` como referencia hasta WF-3. Cuando WF-3 alinee el motor con el peso vigente, cambian esa aserción y la salida esperada de la fixture.

**S5-10**
- La fuente «lectura asistida» pasa a «agente externo (`site_assisted`)».
- En «Verifica», `profeco` es solo referencia (K22, K37) y `site_assisted` es PROVISIONAL, con la misma cita y el mismo HECHO que en S5-04.

**Mismo defecto en el resto del archivo** (lo extendí; ver Supuesto 1):
- §3 «No aceptable»: decía «PROFECO o `site_assisted` compitiendo» y ahora dice «PROFECO compitiendo», con la nota de que el peso de `site_assisted` es PROVISIONAL.
- §4 «Umbrales»: el mismo ajuste.
- §5 «Trazabilidad»: K53 ahora apunta a S5-04, S5-07 y S5-10.

**Encabezado:** «Versión 1.4 … y corrección `correccion-k53-r1`», más una viñeta en «Cambios de la versión 1.4».

**Concordancia** con los textos que el revisor citó:
- `docs/criterios/INDICE.md`, fila K37: «su peso es PROVISIONAL».
- `S5-interfaz.md` 1.6 §3: «PROVISIONAL … **sí**, como `manual`», con el mismo HECHO del motor.

### M-1 (medio) · límites de lote sin fuente

**`docs/estado/bloqueos/S5-agente-externo.md`:** agregué la pregunta **(g), tamaño máximo de un lote**.
- HECHO: los valores que hoy trae C1 1.4.
- De dónde salen: no vienen de RESPUESTAS.md, de §K ni de C1 1.3. Los eligió `k53-aplicar`.
- Opciones:
  - (g1) dejar 500 productos, 500 precios y 100 recetas.
  - (g2) otros topes que dé Axel.
  - (g3) sin tope de renglones; solo «al menos 1».
- En «Mientras tanto»: (g1) como supuesto de `k53-aplicar`, no como decisión de Axel ni como PROPUESTA del orquestador.
- En «Impacto»: cambiarían las descripciones de C1, el 422 de WF-3 y sus pruebas. No toca el esquema de S1 ni las cifras canon.
- El bloqueo sigue ABIERTO y no contiene «Respuesta de Axel».

**`gen2_schemas.py`:** los `items` de `ExternalCatalogImport`, `ExternalPriceImport` y `ExternalRecipeImport` dicen ahora «De 1 a N … PROVISIONAL (bloqueo S5-agente-externo, pregunta g): tamaño máximo del lote, supuesto de k53-aplicar; no es decisión de Axel».

**`gen3_ops.py`:** en `importExternalCatalog`, el 422 por lote vacío o de más de 500 renglones cita la pregunta (g), con el mismo texto.

**`c1_plantilla.md`:**
- Una viñeta nueva en «Lo PROVISIONAL del agente externo».
- La fila «Agente externo (K53)» de §13 nombra el tamaño de los lotes.
- Una fila de Cambios nueva, «C1 v1.4 · corrección r1 (P2d)». Solo cambian textos; `info.version` sigue en 1.4.0.

### M-2 (medio) · prosa de `balance_cents` (I-6)

En `gen2_schemas.py`, `CashFlowFortnights.rows.balance_cents` dice ahora: «Saldo de la quincena = ingresos − gasto variable registrado − pagos de crédito − gastos fijos y venideros (S2 1.3 §3.11: income.total_cents − spent_cents − committed_cents, con spent_cents = variable_spent_cents + fixed_upcoming.paid_cents + credit_payments.household_paid_cents). El gasto variable registrado (variable_spent_cents de S2) no viene en esta fila. No resta el ahorro planificado…».

La fórmula sale de S2-interfaz 1.3 §3.11, líneas 257 a 262. No agregué ningún campo nuevo: queda como PROPUESTA, abajo.

### Regeneración

`python3 -I backend/internal/api/c1gen/build.py` reescribió los 22 archivos: `openapi.yaml`, `routes.go`, `C1-datos.md`, `extractos/INDICE.md` y 18 extractos. El texto nuevo quedó en `openapi.yaml` (4 citas de la pregunta g y la prosa nueva del saldo), en `C1-datos.md`, `presupuesto-a.yaml` y `compras-tiendas-b.yaml`.

## Verificación (comandos y salidas)

| Comando | Salida | Resultado |
|---|---|---|
| `python3 -I /tmp/correccion-k53-r1/edit.py` (cada reemplazo exige 1 sola aparición) | `OK docs/criterios/S5.md 9 reemplazos`; `OK …/S5-agente-externo.md 3`; `OK …/gen2_schemas.py 4`; `OK …/gen3_ops.py 1`; `OK …/c1_plantilla.md 2`; `OK fila de Cambios insertada en la línea 30` | OK |
| `python3 -I backend/internal/api/c1gen/build.py` | 22 «escrito …», rc 0 | OK |
| `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | `rc=0`; C1: 22 `IDÉNTICO`, «los 22 archivos coinciden byte a byte con el generador», `comprobaciones OK: 33`, `RESULTADO: TODO OK`; criterios: `TOTAL 261 (alta 201 / media 60 / baja 0)`, `citas verificadas: 321 \| sin fuente: 0`, `TOTAL comprobadas: 178 \| con diferencias: 0`, `INDICE contra conteo propio: 0 diferencias`; S1: `RESUMEN TODO OK`; U1: `RESULTADO: TODO OK`; `== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)` | OK |
| `verify_c1.py --db` directo (base migrada temporal) | `C1-15 marcas PROVISIONAL en openapi.yaml: 18` (eran 15; suben 3 por los `items` de los tres lotes) `· sin bloqueo abierto: [] · líneas PROVISIONAL en C1-datos.md: 17 · sin bloqueo: []`; `C1-02 … del agente externo: 4 … con serviceToken sin la marca: []`; `RESULTADO: TODO OK` | OK |
| `grep -c "Respuesta de Axel" docs/estado/bloqueos/S5-agente-externo.md` | `0`: el bloqueo sigue contando como abierto para C1-15 | OK |
| `grep -n "hasta que cada cadena\|FIRMES" docs/criterios/S5.md` | Solo quedan en la viñeta de cambios de la 1.4 (l. 7), en el historial de la 1.2 (l. 22) y en S5-04 como cita de la condición que «ya no aplica» | OK |
| `go test ./internal/precios/ -run Fuente -count=1` | `ok finanzas-hogar/internal/precios`. `TestFuenteReferencia` usa `SourceSiteAssisted` (l. 472): HECHO del motor de WF-1 | OK |
| `ls oraculo/fixtures/s5/ \| grep -c s5-fuentes-compiten` | `1` | OK |

## Supuestos

1. **A-1, extendido a §3, §4 y §5 de `S5.md`.** Esas secciones repetían la regla anterior («`site_assisted` … solo referencia» y «… compitiendo» como no aceptable). Si las dejaba, el archivo seguía contestando en firme la pregunta (a). Es el mismo defecto en el mismo archivo y no toqué otra cosa.
2. **Texto de lo PROVISIONAL en A-1.** Uso la PROPUESTA del orquestador («compite como `manual`»), igual que K53 y `S5-interfaz` 1.6 §3. No resuelvo la PREGUNTA del revisor sobre cuál PROVISIONAL rige (ver Pendientes).
3. **M-1, valores.** Los límites 500, 500 y 100 se quedan como estaban y se declaran **supuesto de `k53-aplicar`**, la opción (g1). No es decisión de Axel ni PROPUESTA del orquestador. No elegí valores nuevos.
4. **M-1, `c1_plantilla.md`.** También agregué la pregunta (g) en `c1_plantilla.md`, para que `C1-datos.md` la liste junto con las demás. Fue una viñeta, la fila de §13 y una fila de Cambios. `info.version` no cambia porque solo cambian textos.
5. **M-2.** Solo corregí la prosa. No agregué `variable_spent_cents` a la fila: es PROPUESTA y la decisión es de C1 y S2.

## Bloqueos

- No abrí bloqueos nuevos.
- `S5-agente-externo` sigue **ABIERTO**, ahora con las preguntas (a) a (g).

## Pendientes

- **INDICE de bloqueos.** En `docs/estado/bloqueos/INDICE.md`, la fila de `S5-agente-externo` (l. 525) lista las preguntas (a) a (f) y le falta la (g). Esa fila la dio de alta `p2c-remanentes`, así que no la toqué para no invadir su alcance. Que la actualice el integrador o quien lleve el INDICE.
- **PREGUNTA del revisor, para el orquestador.**
  - Lo PROVISIONAL vigente para (a) es «compite como `manual`», la PROPUESTA del encargo de P2d.
  - La regla de RESPUESTAS.md para lo no decidido es otra: la opción en uso, que en el motor y la fixture es «referencia».
  - Hoy `S5.md`, K53 y `S5-interfaz` usan la primera y dejan la segunda como HECHO del motor. Si el orquestador elige la otra, cambian esos tres textos.
- **Pista 3, fuera de mi alcance.** `docs/design/tools/ds/verificar_sinteticos.py` l. 157 dice «saldo = ingresos − pagos de crédito − fijos y venideros», sin el gasto variable registrado. Hay que revisar si sus datos sintéticos tienen gasto variable en alguna quincena.
- **Bajos B-1 a B-4 de la revisión.** No venían en este encargo:
  - `S5-interfaz` §15 y §16;
  - S5-08 y la salida esperada del último comando de S5-07;
  - `c1_plantilla.md` l. 157;
  - la precisión del reporte de `k53-aplicar`.
- **Reporte de `k53-aplicar`.** `docs/estado/p2c/k53-aplicar.md` l. 71 conserva la prosa anterior del saldo. Es el registro de esa ronda y no lo edité.

## Propuestas para Axel (puerta 4)

- **(g) Tamaño máximo de un lote del agente externo:**
  - (g1) 500 productos, 500 precios y 100 recetas, como está hoy;
  - (g2) otros topes;
  - (g3) sin tope de renglones.
- Siguen abiertas (a) a (f) de `S5-agente-externo`. La más visible es (a): ¿los precios del agente compiten como `manual` o son solo referencia, como PROFECO?
- **PROPUESTA técnica (para C1 y S2, no para Axel):** agregar `variable_spent_cents` a la fila de `CashFlowFortnights`. Así `balance_cents` = `income_cents` − `variable_spent_cents` − `credit_payments_cents` − `fixed_and_upcoming_cents` cuadra con la fila sola.

## Consumo aproximado

- **Lectura:** unos 100 KB, de 400 permitidos (revisión, bloqueo, RESPUESTAS, extractos de §K, `S5.md`, `gen2`/`gen3` por partes, la plantilla, los logs).
- **Escritura a mano:** unos 10 KB de texto nuevo más este reporte (unos 9 KB), de 150 permitidos. Además, `build.py` regeneró 22 archivos.
- **Acciones:** unas 25, de 60 permitidas.
