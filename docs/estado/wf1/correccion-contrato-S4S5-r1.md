# correccion-contrato-S4S5-r1 · reporte final (WF-1, 11-oct-2026)

Corrige los 7 hallazgos medios de `docs/estado/wf1/REVISION-contrato-S4S5-r1.md` (no había altos). No reabre decisiones de §A–§K ni cambia cifras canon. Los 8 bajos (B1–B8) quedan fuera de esta ronda, como pidió la tarea.

## Hecho
- **M1 (S5-01, C1):** el canon pasa a ser **la compra registrada de Bodega contra el carrito de Walmart** (`CompareCarts`, tipo nuevo `Cart`). El carrito de Walmart usa precios `manual` del 2026-10-10 y sus propias cantidades, y existe solo como entrada del motor y de la fixture: no se siembra ningún gasto de 3,233.00. Se declara la fuente de cada semilla según K37. La parte de la API queda PROVISIONAL con el bloqueo nuevo `S5-carrito-comparativa`: cómo arma `getBasketComparison` (`basket = expense`, una sola `expense_id`) la canasta de las otras tiendas, la fuente de la compra de Bodega y si los precios de Walmart se siembran en `price`. Secciones: S5 §2, §3, §7.2, §14, §15.
- **M2 (FORMATO §0 y §4):** se agregan las tablas del `Input` raíz (claves, tipos, `_cents`/`_milli`, obligatoriedad, `unknown_field`/`missing_field`, entidades por código, rutas `*_file`) y la tabla de argumentos y salida de cada función consultable (S4 §12.1–§12.2 y S5 §15.1–§15.2). Se nombran los oráculos `oraculo/s4_run_fixtures.py` y `oraculo/s5_run_fixtures.py`: cubren **todos** los casos, escriben `oraculo/salidas/<motor>.json`, no leen el código Go y los escribe otro agente. Se dice qué casos se verifican de forma independiente (S4 §12.3 y S5 §15.3). `s5-comparativa-canon` lee `data.json` por la clave `report_file`, sin copias.
- **M3 (S4-09):** las 5 recetas las elige el revisor al revisar, entre las 60. La fixture del autor pasa a llamarse `s4-kcal-recetas-autor` (prueba `TestKcalRecetasAutor`) y no sustituye a S4-09. El revisor puede dejar las suyas en `s4-kcal-revisor` (S4 §12, §12.3, §13, §15).
- **M4 (S5-06, K24):** `datos/semilla/tickets/` guarda solo las 5 fotos (`.jpg`, `.jpeg` o `.png`), así que `ls | wc -l` = 5. Las transcripciones van en la carpeta hermana `datos/semilla/tickets-esperado/NN.csv`. No van en `tickets/esperado/` porque `ls` contaría la subcarpeta y daría 6 (S5 §12, §16, §18).
- **M5 (S4-08, S4-11, R9):** `tabla.csv` y `vinculos.csv` suman las columnas `estado` (`propuesto`/`confirmado`), `confirmo` y `fecha_confirmacion`. Los cargadores tienen dos modos: `production`, que no carga lo `propuesto`, y `review`, que lo usa en pruebas, fixtures y medición de WF-1 y lo cuenta aparte. S4-08 se reporta con dos conteos. Quién confirma y cuándo queda PROVISIONAL con el bloqueo nuevo `S4-confirmacion-empates` (opción A: Axel en la puerta 2). Secciones: S4 §2, §9.1, §13, §15 y S5 §11, §16.
- **M6 (S5-01 K-07, C1):** `stores[].identical_count` = renglones (33), como el propio C1 resume el canon («33 idénticos, 3 sustitutos»). K-07 sale en `stores[].piece_count` (61 y 60) y en `identical_pieces` (55). Se pide a C1 que cambie la descripción de `identical_count` (S5 §7.2 paso 6, §15.2, §16).
- **M7 (S4-13, S4-15):** ningún archivo ni subpaquete de `internal/recetas` toca SQLite: la persistencia vive en `internal/api` (S6), que S6-02 mide sin umbral. §3.3 nueva: funciones puras `ValidateRecipe`, `ApplyRecipeEdit` y `RemoveRecipe` (archiva o borra, según C-06), `TestRecetaEditable` sin base, `TestRestriccionesVacias`, y la consulta SQL del campo `health_profile.dietary_restrictions` vacío, que corre el agente del motor S4 sobre una base temporal hecha con goose (S4 §1, §3.3, §13, §15).

## Archivos
- `docs/contratos/motores/S4-interfaz.md` (1.0 → 1.1; 48,076 bytes; 101 líneas de diff)
- `docs/contratos/motores/S5-interfaz.md` (1.0 → 1.1; 40,460 bytes; 92 líneas de diff)
- `docs/estado/bloqueos/S5-carrito-comparativa.md` (nuevo)
- `docs/estado/bloqueos/S4-confirmacion-empates.md` (nuevo)
- `docs/estado/wf1/correccion-contrato-S4S5-r1.md` (este reporte)
- `docs/contratos/motores/productos-formato.md`: **sin cambios** (ningún hallazgo medio lo toca).

## Verificación
- **Cifras canon y ejemplos:** script `python3 -I` con `Decimal` y `Fraction` (en `datos/tmp/corr-s4s5/`, ya borrado), con el modelado nuevo de M1 (Bodega = compra de 39 renglones; Walmart = carrito):
  - K-01 = 309450 ✓; K-02 = 323300 ✓; K-03 = 13850 = 10050 + 5900 − 2100 ✓; K-04/K-05 = 277350 / 287400 ✓; K-06 = 10050 ✓.
  - K-07: `piece_count` 61 / 60 e `identical_pieces` 55 = `articulos_*` ✓. K-08: `identical_count` 33, `substitute_count` 3 ✓. Ganadores 22 / 8 / 3 ✓.
  - Unidad base: 2933/2800, 12692/16406, 3375/2833 ✓. Todos los `unitario × cantidad = total` y todos los `total ÷ cantidad` de los sustitutos son exactos.
  - Opción A del bloqueo (Walmart cotizado con las cantidades de Bodega, ILUSTRATIVO): 335900 ≠ 323300, lo que confirma que la API no reproduce el canon con esa opción.
  - S4 §8.3: 4425 / 2213 / 8850 ✓. Frescura: 29 d → 1; 30 → 60/61 (0.9836); 31 → 59/61; 60 → 30/61 (0.4918); 90 → 0; 91 → no compite ✓.
  - MSJ: 1748.75, 1330.25 y 1648.75 → metas 2099, 1596, 1799, 1846, 2556 y 2056 ✓. IMC: 185, 250, 300, 237 y 185 ✓. RN-15: 1486571 ✓.
- **Estructura** (script `estructura.py`, ya borrado): 0 tablas con columnas desparejas en S4 y en S5; el bloque JSON de S4 §3.1 es válido. No queda ningún nombre viejo (`compare_expenses`, `CompareExpenses`, `s4-kcal-cinco-recetas`, `TestKcalCincoRecetas`, `NN.esperado.csv`, `../../../datos`, «dos compras registradas»). Cada PROVISIONAL nuevo cita un bloqueo que existe. Los 3 PROVISIONAL sin archivo de bloqueo que quedan ya venían de la versión 1.0 y el revisor los aceptó: «PROVISIONAL de S5-09», «mismo bloqueo» de PROFECO y la propuesta K31 del tiempo de Tesseract.
- **Entorno:** `tesseract 5.5.3`, idiomas `osd spa` ✓. `datos/semilla/tickets/` no existe, así que S5-06 sigue «sin medir». `oraculo/compara.py`, `oraculo/fixtures/s4/` y `oraculo/fixtures/s5/` todavía no existen. El artefacto es un contrato sin código: no hay pruebas Go, cobertura ni `compara.py` que correr; se verifican sus cifras y su estructura.
- **Criterios tocados → evidencia:**
  - S4-08 → §9.1 (dos conteos en modo `review`).
  - S4-09 → §12, que fija que el revisor elige las 5 recetas.
  - S4-11 → §9.1 (modo `production`) y `TestNutricionSinConfirmar`.
  - S4-13 → §3.3 (funciones puras, prueba y consulta SQL).
  - S4-15 → §1, sin base en `internal/recetas`.
  - S5-01 → §3, §7.2, §15 (`CompareCarts`, `piece_count`, `identical_pieces`, `report_file`).
  - S5-05 → §11 (ligas confirmadas).
  - S5-06 → §12 (carpetas separadas).
  - FORMATO §0 y §4 → S4 §12.1–§12.3 y S5 §15.1–§15.3.
  - El resto de los criterios sigue con el mismo estado que les dio el revisor.

## Supuestos
- S-11 (S5 §7.2): para `piece_count`, un renglón a granel (cantidad no entera) cuenta como 1 pieza. Ninguna cifra canon lo usa.
- Las fuentes `manual` de la semilla del 10-oct-2026 no son supuesto: quedan PROVISIONAL en el bloqueo `S5-carrito-comparativa` (preguntas 2 y 3). Con `manual` o con `ticket_manual` el canon es el mismo (K37).
- En el `Input`, las entidades se citan por código y las rutas `*_file` se leen desde la raíz del repositorio (desde ahí corren `dump` y el oráculo, FORMATO §4). Es convención de las fixtures; no cambia el esquema ni C1.

## Bloqueos
- `docs/estado/bloqueos/S5-carrito-comparativa.md` (nuevo; M1)
- `docs/estado/bloqueos/S4-confirmacion-empates.md` (nuevo; M5)
- Siguen vigentes los de la versión 1.0: `S4-reglas-lista`, `S5-reglas-comparativa` y `S5-profeco-almacen`.

## Pendientes
- Hallazgos bajos B1–B8 de la revisión r1, sin tocar por instrucción. B5 queda de hecho resuelto, porque S5 ya no usa `$ref` fuera de `_datos/` sino `report_file`. B8 es del reporte `contrato-S4S5.md`, que está fuera de mi alcance.

## Solicitudes fuera de alcance
1. **C1** (`BasketComparison.stores[].identical_count`, vía `backend/internal/api/c1gen`): cambiar la descripción a «Renglones idénticos.» y, si Axel quiere K-07 en pantalla, agregar `stores[].piece_count` e `identical_pieces`. Motivo: M6.
2. **C1 `getBasketComparison`:** solo si en la puerta 2 se elige la opción B o C del bloqueo `S5-carrito-comparativa`, agregar el segundo parámetro (`compare_expense_id` o carrito guardado). Motivo: M1.
3. **`docs/estado/bloqueos/INDICE.md`:** indexar `S5-carrito-comparativa` y `S4-confirmacion-empates`.
4. **Reporte `docs/estado/wf1/contrato-S4S5.md`:** su solicitud 2 (FORMATO §2: admitir `$ref` fuera de `_datos/`) se puede retirar, porque S5 1.1 lee `data.json` con la clave `report_file` del `Input`.
5. **Semilla de S1 y S6:** no sembrar ningún gasto de Walmart del 10-oct-2026 (S5 §3). Si los precios del carrito entran a `price`, se decide en la pregunta 3 del bloqueo `S5-carrito-comparativa`.
6. **S6 (`internal/api`):** `createRecipe`, `updateRecipe` y `deleteRecipe` persisten con las funciones puras de S4 §3.3 (`ValidateRecipe`, `ApplyRecipeEdit`, `RemoveRecipe`) y tienen su prueba por operación (S6-03).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- **`S5-carrito-comparativa`:**
  - Pregunta 1: opción A, no cambiar C1. El canon ya se reproduce en el motor y en la fixture; comparar la compra contra los precios vigentes es lo útil en el día a día. Si quiere ver +138.50 en pantalla, opción C.
  - Pregunta 2: según exista o no ticket impreso de Bodega.
  - Pregunta 3: opción A, la app arranca sin precios de Walmart. Un precio de carrito que no se compró es dato blando, y el del 10-oct-2026 ya es viejo el 9-nov-2026.
- **`S4-confirmacion-empates`:** opción A, Axel confirma en la puerta 2 sobre la lista de renglones `propuesto`. No toca el esquema ni C1, y nada entra sin confirmar.
- Sin valores K31 nuevos.

## Consumo aproximado
- 37 acciones de herramienta.
- ~185 KB leídos.
- ~63 KB de contenido nuevo escrito: scripts de parche y bloques de texto ~43 KB, bloqueos 5.7 KB, scripts de verificación ~5 KB, este reporte ~8 KB. Los parches reescribieron S4 (48 KB) y S5 (40 KB) completos con `python3 -I`.
