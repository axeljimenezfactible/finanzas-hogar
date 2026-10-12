# Corrección del contrato S4/S5 · ronda 4 (WF-1, 11-oct-2026)

Agente: `correccion-contrato-S4S5-r4`. Revisión atendida: `docs/estado/wf1/REVISION-contrato-S4S5-r4.md` (0 altos, 2 medios). Corregí los dos medios (M-r4-1 y M-r4-2). No toqué los bajos (B-r4-1 a B-r4-3) ni nada fuera de esos dos hallazgos, y no reabrí decisiones de §A–§K. Ninguna cifra canon cambia.

## Hecho
- **M-r4-1 (S5-03), en `S5-interfaz.md` 1.4.** El total de cada tienda en la comparativa por canasta usa ahora una sola regla: el menor costo de paquetes completos para cubrir la cantidad a comprar, también en la tienda recomendada del renglón (§7.3, párrafo nuevo «Total por tienda: una sola regla»).
  - La recomendación por renglón (`items[].recommended_*`, `estimated_price_cents`, `is_stale_price`) sigue siendo `QuoteFor` sin elección, igual que en S4 (§7.3, «Recomendación por renglón»).
  - Se dice que la línea de la tienda recomendada en `prices[]` puede mostrar otra presentación, marcada como sustituto. En ese caso, lo que el renglón suma al total de esa tienda es menor que `estimated_price_cents`.
  - Se ajustaron el ejemplo de §7.3, la fila `CompareBasket` de §14, la viñeta `recommended` de §14.1 («con la misma regla en **todas** las tiendas, también la recomendada»), la fila de `compare_basket` en §15.2 y §16/§18.
  - Caso `s5-cotizacion-recomendada`: Walmart cotiza 750 g a 2200 (sustituto) y Bodega 500 g a 1500 (sustituto). Totales: Bodega 1500 y Walmart 2200 (+700); `cheapest_store` = Bodega.
  - Caso nuevo `s5-canasta-regla-unica`, con el contraejemplo de la revisión más un artículo `manual` con `quantity_milli`: Walmart **9430** y Bodega **9460**, `cheapest_store` = Walmart. Con la regla vieja, Walmart habría salido en 10780 y la canasta habría señalado a Bodega. Prueba nueva: `TestMejorTiendaCanastaReglaUnica`.
- **M-r4-2 (S4-03, S4-06), en `S4-interfaz.md` 1.4 y `S5-interfaz.md` 1.4.**
  - S4 §6.6 (nuevo):
    - cantidad a comprar = `quantity_milli` si no es nula; si es nula, `needed_milli` (lo dice C1: «NULL en calculated = comprar la necesidad»);
    - `packages_to_buy` y `estimated_cents` salen de esa cantidad; la presentación y la tienda no cambian;
    - `leftover_milli` = max(0, paquetes × contenido − `needed_milli`);
    - un artículo `manual` lleva `needed_milli` = `null`, `reasons` = [`manual`] y `recipes` vacía, y su `quantity_milli` es obligatoria;
    - un artículo `manual` sin presentación elegida ni cotización sale con `null` y va en `missing_prices[]`;
    - los artículos `calculated` y `manual` suman en `groups[]` y en la raíz, y `is_bought` no resta;
    - un artículo de texto libre no entra en S4.
  - S4 §6.1, §6.2 (pasos 5 y 6) y §6.5 (total de `groups[]`) remiten a §6.6.
  - S4 §12.1: `list_choices.quantity_milli` es obligatoria en `manual`. §12.2: `shopping_list` devuelve `items[].quantity_milli`, y la traducción a C1 lo pasa a S5.
  - S4: caso nuevo `s4-lista-cantidad-manual` (§12) con el ejemplo de §6.6 (arroz 3 paquetes / 2500000 / 8400; frijol 1 / 0 / 3000; aceite `manual` comprado 2 / 1892000 / 7980; servilletas sin precio; grupo Walmart y raíz **19380**). Pruebas nuevas: `TestListaCantidadFijadaYManual` (S4-03) y `TestCostoCantidadFijada` (S4-06).
  - S5 §7.3 y §15.1: `shopping_list` trae `quantity_milli` (o `null`) y `needed_milli` puede ser `null`; si los dos son `null`, `invalid_argument`. La cantidad a comprar se calcula con la regla de S4 §6.6 y la usan `estimated_price_cents`, `items[].quantity_milli` y el total de cada tienda.
- Versión 1.3 → 1.4 en los dos contratos, cada uno con su bloque «Cambios de la versión 1.4». A la entrada M-r3-2 de S5 se le agregó una nota de que la 1.4 cambia cómo se cotiza la tienda recomendada.

## Archivos
- `docs/contratos/motores/S4-interfaz.md` (1.3 → 1.4; §6.6 nueva, más §6.1, §6.2, §6.5, §12, §12.1, §12.2, §13 y §15).
- `docs/contratos/motores/S5-interfaz.md` (1.3 → 1.4; §7.3, §14, §14.1, §15, §15.1, §15.2, §16 y §18).
- `docs/contratos/motores/productos-formato.md`: sin cambios (ningún hallazgo lo pide).
- `docs/estado/wf1/correccion-contrato-S4S5-r4.md` (este reporte).

## Verificación
- `python3 -I datos/tmp/corr-s4s5-r4/verif.py` (solo biblioteca estándar, `Fraction` y `Decimal`; lee `data.json`, los dos contratos y los criterios): **98 OK, 0 FAIL**.
  - **Canon S5-01** desde los renglones de `data.json`: idénticos 277350 / 287400 (+10050); totales 309450 / 323300 (+13850 = 10050 + 5900 − 2100), iguales a `totales`; piezas 61 / 60 / 55; renglones 33 / 3; ganadores 22 / 8 / 3.
  - **S5-02:** 2933 / 2800, 12692 / 16406, 3375 / 2833. **K40:** 0.9836 / 0.4918 / 0.0000.
  - **S4 §8.3:** 4425 / 2213 / 8850. **S4-07:** MSJ 1748.75 y 1330.25; con el factor 1.2 y los ajustes: 2099 / 1596 / 1799 / 1846. **RN-15:** 1486571.
  - **Implementación independiente de S5 §14.1 y §7.3 1.4:**
    - `s5-cotizacion-recomendada`:
      - `quote_for` sin elección → Walmart 1 kg a 2800; con Bodega → 750 g a 2200; con 750 g → Walmart; frijol → Bodega 500 g;
      - `compare_basket` → Walmart 2200 (750 g, sustituto) y Bodega 1500 (500 g, sustituto), +700, `cheapest_store` = Bodega;
      - `estimated_price_cents` = 2800.
    - `s5-canasta-regla-unica`:
      - totales 9430 / 9460 y `cheapest_store` = Walmart; con la regla 1.3 serían 10780 / 9460;
      - por unidad base: 2900 y 2920 por kg; 4218 y 4228 por litro;
      - el contraejemplo de la revisión, solo con arroz, da 1450 / 1460 y Walmart.
  - **S4 §6.6:** los 4 artículos, el grupo (`chain`, `walmart`) = 19380, el grupo (`chain`, `null`) = 0 y la raíz 19380. Sin cantidad fijada: arroz 1 paquete / 500000 y frijol 2 paquetes (6000) / 800000.
  - **Cifras en el texto:** cada cifra de arriba aparece en el contrato que la cita.
  - **Estructura:**
    - cercas pares; S4 con 1 bloque JSON válido, S5 sin bloques;
    - versión 1.4 con la cita de la revisión r4 en los dos;
    - S4 §15 cubre S4-01 a S4-15 y S5 §18 cubre S5-01 a S5-13;
    - cada `-run` de `docs/criterios/S4.md` (Canal, Costo, Escalado, Inventario, Kcal, Lista, Plan, Salud) y de `S5.md` (Comparativa, Frescura, Fuente, LecturaAsistida, MejorTienda, Membresia, Mercado, PreciosViejos, Unidad) encuentra pruebas, todas de su criterio. Las nuevas caen en S4-03, S4-06 y S5-03;
    - todo bloqueo citado existe (salvo `S4-meta-ajustes`, que es hipotético según S4-07 1.2);
    - no quedan frases viejas: «no recomendadas», «Walmart cotiza 1 kg (idéntico)», «aunque otra de la misma tienda cubra la necesidad con menos costo», «techo(`needed_milli / size_milli`)», «Entrada: la lista de compras de S4 (producto y `needed_milli`)».
- **Criterios tocados y evidencia:**
  - S4-03: §6.2 pasos 5 y 6, §6.5 (`groups[]`), §6.6, §12.1, §12.2, caso `s4-lista-cantidad-manual` y `TestListaCantidadFijadaYManual`.
  - S4-06: §6.6 (`estimated_cents` con la cantidad fijada) y `TestCostoCantidadFijada`.
  - S5-03: §7.3 (una sola regla), §14.1, los casos `s5-cotizacion-recomendada` y `s5-canasta-regla-unica`, y `TestMejorTiendaCanastaReglaUnica`.
  - Los demás criterios no cambian.
- **Código de otros agentes** (informativo; no es de este artefacto y no lo toqué):
  - `go vet ./internal/precios/... ./internal/recetas/...`: `internal/recetas/types.go:229:13: undefined: Nutrition`.
  - `gofmt -l` marca `internal/precios/seed.go`, `types.go`, `internal/recetas/num.go` y `types.go`.
  - `go test -cover`: `precios/ocr` 99.3 %; `precios`, `asistida` y `profeco` 0.0 % (sin pruebas); `recetas` no compila.
  - `oraculo/fixtures/s5/`, `oraculo/s4_run_fixtures.py` y `oraculo/s5_run_fixtures.py` no existen, y `oraculo/fixtures/s4/` solo tiene `_datos/`. Por eso `compara.py s4` y `compara.py s5` todavía no se pueden correr.
- Temporales en `datos/tmp/corr-s4s5-r4/` (scripts de edición y verificación, y copia previa de los dos contratos): borrados al terminar.

## Supuestos
1. **`leftover_milli`** sigue la fórmula del paso 6 de §6.2, que ya existía (paquetes × contenido − `needed_milli`), y ahora lleva `max(0, …)`. Es nuevo que pueda dar negativo, cuando la persona fija menos que la necesidad, y S4-03 prohíbe cantidades negativas. En un artículo `manual` se toma `needed_milli` = 0. Es una cifra del motor que no es canon y C1 no la publica.
2. **`is_bought` no resta del total.** C1 llama a `total_cents` «Total estimado», U1 N-150 pide «Totales por grupo y total general» y ninguno de los dos excluye lo comprado. Un subtotal de lo que falta por comprar no existe en C1 ni en U1; lo dejo como propuesta y no lo invento.
3. **Cada renglón `manual` es un artículo propio**, aunque el producto tenga además un artículo `calculated`. C1 solo hace único el artículo `calculated` por producto (`ux_shopping_list_item_calculated`). En `compare_basket` cada artículo es un renglón y se cotiza por separado.
4. **Un artículo `manual` sin presentación elegida ni cotización** sale con `null` y va en `missing_prices[]`, que es la convención «sin dato» de §12.2. No hay presentación de respaldo porque no viene de un ingrediente.

## Bloqueos
- Ninguno nuevo. Siguen vigentes, sin cambios: `docs/estado/bloqueos/S4-asignacion-mercado.md`, `S4-reglas-lista.md`, `S4-confirmacion-empates.md`, `S5-reglas-comparativa.md`, `S5-carrito-comparativa.md`, `S5-profeco-almacen.md` y `C1-campos-calculados.md`. El abierto de la puerta 2 es `S1-parametros-avisos.md`, que no toca este artefacto.

## Pendientes
- Los bajos de la revisión r4 (B-r4-1 `StalePricesForPlan` con presentación elegida, B-r4-2 «menor id» en §7.3, B-r4-3 variante sin mercado), los bajos anteriores (B-r3-1, B-r3-2) y los de r1 y r2. No los tocan los hallazgos de esta ronda.
- Artículos `manual` de texto libre (C1 `free_text`): S4 no los calcula (§6.6). Falta decidir en qué grupo de `ShoppingList.groups` los pone la API, porque no tienen vertical ni tienda recomendada. Le toca a U1, C1 o S6, no al motor.
- Fixtures y oráculos (agentes de implementación):
  - escribir `oraculo/fixtures/s4/s4-lista-cantidad-manual.json` (§6.6) y `oraculo/fixtures/s5/s5-canasta-regla-unica.json` (§15);
  - escribir el esperado de `compare_basket` en `s5-cotizacion-recomendada` con la regla 1.4;
  - implementar `quantity_milli` en `oraculo/s4_run_fixtures.py` y `oraculo/s5_run_fixtures.py` cuando existan.

## Solicitudes fuera de alcance
1. **`backend/internal/precios/basket.go` (`ListItem`, `CompareBasket`).** Hay que alinearlo con S5 1.4:
   - `ListItem` gana `QuantityMilli *int64`, y `NeededMilli` puede ser nulo (artículo `manual`); hoy rechaza `NeededMilli <= 0`;
   - la cantidad a comprar es `QuantityMilli` si no es nula y, si no, `NeededMilli`; si las dos son nulas, `invalid_argument`;
   - `options` y `QuantityMilli` de la salida usan esa cantidad;
   - la recomendación por renglón es `QuoteFor` sin elección (§14.1), no `winners` + `breakTie` con `assigned` (esa es la regla anterior a la 1.3).

   El total de cada tienda con el menor costo de paquetes ya coincide con la regla única. Motivo: M-r4-1 y M-r4-2; si no se alinea, el motor y el oráculo divergen.
2. **Al implementador de la lista de compras de `internal/recetas`** (todavía no hay código): §6.6 completa. `list_choices.quantity_milli` cambia `packages_to_buy`, `estimated_cents` y `leftover_milli`; los artículos `manual` y los comprados suman en `groups[]` y en la raíz. Motivo: M-r4-2.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- (Opcional, puerta 2) Un subtotal «falta por comprar» en la lista, que no cuente lo marcado como comprado. Hoy C1 y U1 solo tienen el total estimado de la lista completa. El contrato lo deja así y no agrega el campo. Sumarlo sería un campo calculado nuevo en C1 (`ShoppingList`), sin efecto en ninguna cifra canon.

## Consumo aproximado
- Acciones de herramienta: 30.
- Leído: ~190 KB (protocolo 6, revisión r4 16, §K 15, criterios S4 y S5 37, secciones de S4 y S5 ~45, extracto C1 ~12, U1, RN y bloqueos ~8, reporte r3 ~4, `data.json` y código Go ~8, revisiones propias de salida ~10).
- Escrito: ~45 KB de texto propio (scripts de edición ~13, §6.6 4.5, párrafo de §7.3 2.7, verificación 9, este reporte ~12). Los scripts reescribieron en su lugar `S4-interfaz.md` (70 KB) y `S5-interfaz.md` (54 KB); lo nuevo en ellos son unos 14 KB.
