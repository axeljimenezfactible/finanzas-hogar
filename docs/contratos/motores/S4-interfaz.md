# S4 · Interfaz del motor de recetas, comidas y salud (`backend/internal/recetas`)

Versión 1.4 · 11-oct-2026 · WF-1, agente `contrato-S4S5` (Opus); corregida por `correccion-contrato-S4S5-r1` (revisión `docs/estado/wf1/REVISION-contrato-S4S5-r1.md`) , por `correccion-contrato-S4S5-r2` (revisión `docs/estado/wf1/REVISION-contrato-S4S5-r2.md`) , por `correccion-contrato-S4S5-r3` (revisión `docs/estado/wf1/REVISION-contrato-S4S5-r3.md`) y por `correccion-contrato-S4S5-r4` (revisión `docs/estado/wf1/REVISION-contrato-S4S5-r4.md`). Documento de contrato, **sin código**: fija qué recibe y qué entrega el paquete, con qué aritmética y contra qué criterios se prueba. Lo implementan los agentes de WF-1 que nombra el plan §3.4 fila S4 (motor Opus; recetario y cargador de nutrición Sonnet).

**Fuentes que manda este documento (en este orden):** `docs/plan/decisiones.md` §D, K7, K8, K23, K25, K36, K38, K43 y K44; `docs/contratos/reglas-negocio.md` 1.3 (RN-14, RN-15, RN-16, RN-17, RN-19; convenciones C-01, C-02, C-05, C-06, C-11); esquema 1.3 (`backend/db/migrations/00001` a `00003`); contrato C1 v1.2 por extractos (`comidas.yaml`, `compras-tiendas-a.yaml`, `compras-tiendas-b.yaml`, `salud.yaml`); criterios `docs/criterios/S4.md` 1.2. Formato del catálogo: `docs/contratos/motores/productos-formato.md`. Precios: `docs/contratos/motores/S5-interfaz.md`.

**Marcas.** «PROVISIONAL (bloqueo X)» = regla que no está decidida en §A–§K ni en RN y queda con la opción indicada hasta la puerta 2; no cambia ninguna cifra canon. «ILUSTRATIVO» = número de una fixture que no viene del oráculo ni del reporte.

**Cambios de la versión 1.4 (11-oct-2026, ronda 4 de revisión):**
- M-r4-2: la cantidad a comprar es `quantity_milli` si la persona la fijó y `needed_milli` si no (C1: «NULL en calculated = comprar la necesidad»). `packages_to_buy` y `estimated_cents` salen de esa cantidad; el sobrante de paquete se sigue midiendo contra `needed_milli`, sin negativos. Un artículo `manual` (con producto) sigue la misma regla con su `quantity_milli` obligatoria y `needed_milli` nulo. Los artículos `calculated` y `manual` suman en `groups[]` y en la raíz, y `is_bought` no resta. `shopping_list` devuelve `items[].quantity_milli`, que la API pasa a S5 (`S5-interfaz.md` §15.1). Caso nuevo `s4-lista-cantidad-manual` (§6.1, §6.2, §6.5, §6.6, §12, §12.1, §12.2, §13, §15).

**Cambios de la versión 1.3 (11-oct-2026, ronda 3 de revisión):**
- M-r3-1: `shopping_list` entrega, por artículo, la recomendación y el precio que rigen (`recommended_store`, `recommended_presentation`, `price_store`, `price_observed_date`, `is_stale_price`, `estimated_cents`) y, en la raíz, el total `estimated_total_cents`. §12.2 trae la traducción a C1: `visit_date` → `visit_weekday`, `estimated_total_cents` → `total_cents` y `ShoppingList.total_cents` = Σ de los grupos. La opción A de `S4-asignacion-mercado` deja un artículo de mercado en una tienda distinta de la que recomienda S5, y eso exige cambiar C1 v1.2: solicitud de un campo calculado `assigned_store_id` y de la redacción de `recommended_store_id`, `estimated_price_cents` y `groups[].store_id`, PROVISIONAL y con una variante sin campo nuevo (§6.5, §12, §12.2, §13, §15).
- M-r3-2: la API pide las cotizaciones con `date` = «hoy» y con la regla de presentación y desempate de `S5-interfaz.md` §14.1 (§6.2 paso 5, §8.1, §12.1).

**Cambios de la versión 1.2 (11-oct-2026, ronda 2 de revisión):**
- M-r2-2: cada cotización de `quotes` lleva `recommended` (exactamente una `true` por producto con cotizaciones) y hay una sola regla de «cotización que rige» para `recipe_detail`, `week_plan` y `shopping_list` (§8.1, §6.2 paso 5, §12.1, §13).
- M-r2-3: el generador escribe `kcal_por_porcion` y `fuentes_kcal[]` con `tabla.csv` en modo `review`; el cargador del recetario los coteja solo en modo `review`; en `production` no coteja, carga las 60 recetas y las que tengan un dato sin confirmar salen «incompletas», sin `ErrKcalMismatch` (§3.2, §13, §14).
- M-r2-4: tienda de cada artículo, artículos de vertical `market` (varios mercados, ninguno, cadena más barata según S5), orden y total de `groups[]`: lo no deducible de RN-16 queda PROVISIONAL (bloqueo nuevo `S4-asignacion-mercado`) (§6.5, §11, §12, §12.1, §12.2, §13).

**Cambios de la versión 1.1 (11-oct-2026, ronda 1 de revisión):**
- M2: `Input` de las fixtures (§12.1), argumentos y salida de cada función consultable (§12.2) y oráculo `oraculo/s4_run_fixtures.py` para todos los casos (§12.3).
- M3: las 5 recetas de S4-09 las elige el revisor al revisar; la fixture del autor pasa a llamarse `s4-kcal-recetas-autor` y no sustituye a S4-09 (§12, §13, §15).
- M5: `tabla.csv` lleva `estado`, `confirmo` y `fecha_confirmacion`; en modo producción el cargador no carga lo no confirmado; quién confirma y cuándo es PROVISIONAL (bloqueo `S4-confirmacion-empates`) (§9.1).
- M7: la persistencia vive fuera de `internal/recetas/...`, en la capa de la API (S6); recetas de la casa con funciones puras y la consulta del campo de restricciones vacío (§1, §3.3, §13).

## 1. Paquete, límites y dependencias

| Tema | Regla |
|---|---|
| Ruta | `backend/internal/recetas` (módulo `finanzas-hogar`, Go 1.27). Subcomando `s4 dump` del binario común `backend/cmd/motores` (§12, FORMATO §4). |
| Dependencias | Solo biblioteca estándar. El núcleo es de **funciones puras** (sin base de datos, sin reloj, sin red): reciben «hoy» y los datos ya leídos. La lectura y escritura en SQLite (modernc.org/sqlite) vive **fuera** de `internal/recetas/...`: la hace la capa de la API (`backend/internal/api`, S6), que se mide sin umbral (S6-02, K25). Ningún archivo ni subpaquete de `internal/recetas` importa `database/sql` ni modernc.org/sqlite, así que todo lo que mide S4-15 lo cubren las pruebas del núcleo (§3.3). |
| Red | Ninguna llamada de red ni a USDA/OFF/IA en ejecución (S4-10, S4-14): `grep -rn "http" internal/recetas --include=*.go | grep -v _test.go` = 0. |
| Aritmética | Enteros `int64` en milésimas (C-11) y centavos (C-01); intermedios como fracción exacta (`math/big.Rat` o par numerador/denominador `int64` con verificación de desbordamiento). **Prohibido `float32`/`float64`** fuera de `_test.go` (S4-02, S4-15). |
| Redondeo | ROUND_HALF_UP (C-02) **una sola vez, al final de cada cálculo publicado**. Helper local `roundHalfUp(num, den)` (mitad se aleja de cero). Si `internal/dominio` publica el mismo helper, el integrador puede unificarlos (solicitud en el reporte). |
| Hoy y semana | «Hoy» lo calcula el llamador en America/Mexico_City (C-05) y lo pasa como fecha ISO. Semana lunes a domingo (K8). |
| Precios | S4 **no importa** `internal/precios`. Declara la interfaz `PriceSource` (§8.1) y la API conecta la implementación de S5. Así no hay ciclo de importación. |
| Cobertura | ≥ 85 % de sentencias en `internal/recetas` (K25, S4-15). |

**Archivos previstos** (nombres orientativos, uno por tema): `types.go`, `recipebook.go` (cargador y validación del recetario), `recipeedit.go` (recetas de la casa, §3.3), `scale.go`, `plan.go`, `shopping.go`, `inventory.go`, `cost.go`, `health.go`, `kcal.go`, `nutrition.go` (cargador de `tabla.csv`), `channel.go`, `errors.go`.

## 2. Tipos del dominio (Go, nombres en inglés corto)

| Tipo | Campos (unidad) | Origen en el esquema |
|---|---|---|
| `Product` | `ID int64`, `Code string` (id del catálogo), `Name string`, `BaseUnit` (`g`, `ml`, `piece`), `GramsPerUnitMilli *int64` (ml o pieza → gramos ×1000), `GroupID *int64`, `Channel` (`market`, `chain`) | `product`, `product_group.purchase_channel` |
| `Presentation` | `ID`, `ProductID`, `Name`, `SizeMilli int64` (contenido en milésimas de la unidad base) | `presentation.size_milli` |
| `Recipe` | `ID`, `Code string` (id del recetario), `Name`, `Slots []SlotCode`, `BaseServings int` (2), `PrepMinutes`, `CookMinutes *int`, `Origin` (`claude_seed`, `household`), `Steps []string`, `Ingredients []Ingredient`, `IsActive bool` | `recipe`, `recipe_meal_slot` |
| `Ingredient` | `ProductID`, `PresentationID *int64` (presentación comprable sugerida), `QuantityMilli int64` (para `BaseServings`), `DisplayText string`, `Optional bool`, `SortOrder int` | `recipe_ingredient` |
| `SlotCode` | `breakfast`, `lunch`, `dinner`, `snack` (los 4 renglones de `meal_slot`; `snack` es opcional) | `meal_slot.code` |
| `PlanEntry` | `ID`, `Date`, `Slot`, `PersonID`, `RecipeID *int64`, `ServingsMilli *int64`, `EatingOut bool` | `meal_plan_entry` |
| `Person` | `ID`, `Code`, `HasMealPlan bool` (titular o co-titular; las externas no) | `person` + disparadores de C-13 |
| `InventoryItem` | `ID`, `ProductID *int64`, `RecipeID *int64`, `Mode` (`stock`, `consumption_rate`), `QuantityMilli`, `CountedDate *Date`, `RateQuantityMilli *int64`, `RatePeriodDays *int`, `ReorderLevelMilli *int64`, `ExpiresDate *Date` | `inventory_item` |
| `Movement` | `ItemID`, `Date`, `Kind` (`purchase`, `recipe_leftover`, `plan_consumption`, `adjustment`), `Reason *` (`count`, `consumption`, `waste`), `DeltaMilli int64` | `inventory_movement` |
| `NutritionFact` | `ProductID`, `PresentationID *int64`, `KcalMilliPer100g int64`, `ProteinMg/FatMg/CarbsMg *int64` (por 100 g), `Source` (`usda_fdc`, `open_food_facts`, `manual`, más `mx_table` y `store_label` que admite el esquema), `SourceRef string`, `License string`, `ObservedDate`, `Proposed bool` (renglón `propuesto` de `tabla.csv`; solo existe en modo `review`, §9.1) | `nutrition_fact` |
| `HealthProfile` | `PersonID`, `HeightMM *int`, `BirthDate *Date`, `Sex *` (`female`, `male`), `Activity *` (5 niveles), `Goal *` (`lose`, `maintain`, `gain`), `AdjustLoseKcal *int`, `AdjustGainKcal *int`, `DailyKcalTarget *int`, `DietaryRestrictions *string` (vacío) | `health_profile` (00001 + 00002 + 00003) |
| `WeightEntry` | `PersonID`, `Date`, `WeightG int` | `weight_entry` |
| `Settings` | `HouseAdjustLoseKcal int` (0), `HouseAdjustGainKcal int` (0) | `app_setting health.goal_adjust_lose_kcal`, `health.goal_adjust_gain_kcal` (00003, K36, K44) |

## 3. Recetario semilla `datos/semilla/recetario.json` (S4-01)

**Ubicación del esquema:** `datos/semilla/recetario.schema.json` (JSON Schema 2020-12). Lo escribe el agente del recetario con exactamente los campos de abajo; S4-01 lo valida con `python3 -I -m jsonschema` o con la prueba Go `TestRecetarioValido`. **Llave de conteo:** `recetas` (`len(d['recetas'])` = 60).

### 3.1 Estructura

```json
{
  "version": "1.0",
  "fecha": "2026-10-11",
  "catalogo": "datos/semilla/productos.json",
  "recetas": [
    {
      "id": "espagueti-con-atun",
      "nombre": "Espagueti con atún y elote",
      "descripcion": "Pasta corta de diario, lista en 20 minutos.",
      "tiempos": ["lunch", "dinner"],
      "porciones_base": 2,
      "minutos_preparacion": 10,
      "minutos_coccion": 12,
      "origen": "claude",
      "ingredientes": [
        {"producto_id": "espagueti", "cantidad_milli": 250000, "unidad": "g",
         "presentacion_id": "espagueti.la-moderna-250g", "texto": "1 paquete de 250 g", "opcional": false}
      ],
      "pasos": ["Cuece la pasta en agua con sal 10 minutos.", "..."],
      "kcal_por_porcion": null,
      "fuentes_kcal": []
    }
  ]
}
```

### 3.2 Campos y reglas

| Campo | Tipo y regla | Columna |
|---|---|---|
| `id` | cadena `^[a-z0-9]+(-[a-z0-9]+)*$`, única en el archivo; estable (las fixtures y la revisión la citan) | ninguna: la base no tiene `code` en `recipe`; el cargador empata por `nombre` (único, NOCASE) |
| `nombre` | es-MX, 3 a 80 caracteres, único sin distinguir mayúsculas; nunca un nombre en inglés (S4-11) | `recipe.name` |
| `descripcion` | opcional, ≤ 200 caracteres | `recipe.description` |
| `tiempos` | arreglo no vacío, sin repetidos, de `breakfast`, `lunch`, `dinner`, `snack` | `recipe_meal_slot` |
| `porciones_base` | entero; **2** en todo el recetario de Claude (decisiones D); el esquema de la base admite 1 a 50 para recetas de la casa | `recipe.base_servings` |
| `minutos_preparacion` | entero ≥ 0, obligatorio | `recipe.prep_minutes` |
| `minutos_coccion` | entero ≥ 0, obligatorio (0 si no se cuece) | `recipe.cook_minutes` |
| `origen` | `claude` (recetario inicial) o `propia` (Axel o Ximena) | `recipe.origin` = `claude_seed` o `household` |
| `ingredientes[]` | arreglo no vacío; un `producto_id` una sola vez por receta (`UNIQUE (recipe_id, product_id)`) | `recipe_ingredient` |
| `.producto_id` | id que existe en `productos.json` | `product_id` |
| `.cantidad_milli` | entero > 0 en **milésimas de la unidad base** del producto (250 g = 250000; media pieza = 500; 15 ml = 15000). Prohibido «al gusto»: la sal o las especias llevan gramos exactos | `quantity_milli` |
| `.unidad` | `g`, `ml` o `pieza`; debe ser igual a la `unidad_base` del producto (redundancia que el validador cruza) | — |
| `.presentacion_id` | obligatorio: presentación comprable del mismo producto en `productos.json` | `presentation_id` |
| `.texto` | opcional, cómo lo diría la receta («1/2 taza») | `display_text` |
| `.opcional` | booleano; **siempre `false` en el recetario de Claude** (ver §5.4) | `is_optional` |
| `pasos` | arreglo no vacío de cadenas es-MX | `recipe.instructions` (unidos con salto de línea y numerados al mostrar) |
| `kcal_por_porcion` | entero o `null`; **derivado, solo para revisión**: lo escribe el generador del recetario con `tabla.csv` leída en modo `review` (renglones `propuesto` y `confirmado`, §9.1), las reglas de §9.2 y la fórmula `kcal_per_serving` de §9.3 con `porciones_base`; `null` si algún ingrediente no tiene renglón en `tabla.csv` o, en `ml` o `pieza`, no tiene `grams_per_unit_milli`. El cargador lo coteja **solo en modo `review`** (párrafo «Cargador del recetario y modos») | ninguna (se calcula) |
| `fuentes_kcal[]` | derivado, solo para revisión: un `{producto_id, fuente, id_fuente}` por cada ingrediente que tiene renglón en `tabla.csv` (de cualquier `estado`), en el orden de `ingredientes[]`, copiado tal cual del CSV (`fuente` = `USDA`, `OFF` o `manual`); se coteja como `kcal_por_porcion`, solo en modo `review` | ninguna (vive en `nutrition_fact`) |

**Validación cruzada del cargador** (además del JSON Schema): todo `producto_id` y `presentacion_id` existe y la presentación es del producto; `unidad` = unidad base; 60 recetas exactas en la semilla; ningún nombre visible igual a un nombre en inglés de USDA (`nombre_fuente` de `tabla.csv`) (S4-11). Las recetas de la casa usan las mismas reglas por la API (`createRecipe`, `updateRecipe`, `deleteRecipe`) a través de las funciones puras de §3.3; el borrado es archivo lógico (`is_active = 0`) cuando el plan las usa (C-06; S4-13).

**Cargador del recetario y modos** `LoadRecipeBook(path, catalog, nutrition, mode)`, con el mismo `mode` que `LoadNutritionCSV` (§9.1) y la tabla nutricional cargada en ese mismo modo:
- **`review`** (pruebas Go, fixtures, `TestRecetarioValido` y la revisión de WF-1): recalcula `kcal_por_porcion` y `fuentes_kcal[]` de cada receta con las reglas de la tabla de arriba; si alguno no coincide con lo escrito (incluidos `null` contra un número y un elemento de más, de menos o en otro orden), devuelve `ErrKcalMismatch` y esa receta no se carga. Así, el número que ve el revisor (S4-09) es el que calcula el motor.
- **`production`** (carga a la base de la app): **no coteja** e ignora los dos campos derivados, que no tienen columna. Toda receta válida se carga aunque tenga ingredientes con renglón `propuesto`: sus kcal se calculan al consultar con lo que entró a `nutrition_fact` (solo lo `confirmado`, §9.1) y, mientras falte un dato, la receta sale «incompleta» (`kcal` ausente y `missing_nutrition`, §9.3); nunca `ErrKcalMismatch`. La semilla de producción carga las 60 recetas antes y después de la puerta 2 (`TestRecetarioProduccion`).
- Confirmar un renglón (`propuesto` → `confirmado`) no cambia `kcal_100g`, `fuente` ni `id_fuente`, así que el cotejo en `review` sigue valiendo. Si al confirmar se corrige un valor, el generador vuelve a escribir los dos campos; el cotejo detecta un recetario que no se regeneró.

### 3.3 Recetas de la casa y restricciones alimentarias (S4-13)
- **Funciones puras** (`recipeedit.go`, sin base): `ValidateRecipe(r, catalog)` aplica las reglas de §3.2 a una receta nueva o editada (`origin = household`; `porciones_base` de 1 a 50, el rango del esquema); `ApplyRecipeEdit(old, edit, catalog)` devuelve la receta editada y validada sin cambiar su código ni su origen (todo el recetario es editable, decisiones D); `RemoveRecipe(r, usedByPlan)` devuelve `archive` (`is_active = 0`) si alguna entrada del plan la usa y `delete` si no tiene historia (C-06: los catálogos con historia se archivan y sus llaves `RESTRICT` impiden borrarlos).
- **`TestRecetaEditable`** prueba esas funciones sobre un recetario en memoria: alta de una receta de la casa válida (y rechazo de una inválida), edición, y borrado (archivo si el plan la usa; borrado si no). No abre SQLite.
- **Persistencia:** insertar, editar y borrar en SQLite lo hacen `createRecipe`, `updateRecipe` y `deleteRecipe` en la capa de la API (`backend/internal/api`, S6), con su prueba por operación (S6-03). No hay `store.go` ni subpaquete con base dentro de `internal/recetas`.
- **Restricciones:** el núcleo no lee ni aplica `DietaryRestrictions` (decisiones D: «sin restricciones alimentarias»). **`TestRestriccionesVacias`**: el plan, la lista y las kcal calculados con perfiles de restricciones vacías o nulas son iguales a los calculados sin esos perfiles (nada se filtra).
- **Consulta del campo vacío (S4-13):** la corre el agente del motor S4 en su verificación y la anota en la línea `S4-13`. Base temporal con el comando de S6-03, `goose -dir backend/db/migrations sqlite3 datos/tmp/s4-13.db up` (si `goose` no está instalado, lo instala y lo reporta), más la semilla de S1 si ya existe (el reporte dice si la aplicó). Esperado: `sqlite3 datos/tmp/s4-13.db "SELECT COUNT(*) FROM pragma_table_info('health_profile') WHERE name = 'dietary_restrictions'"` = 1 (el campo existe) y `sqlite3 datos/tmp/s4-13.db "SELECT COUNT(*) FROM health_profile WHERE COALESCE(dietary_restrictions, '') <> ''"` = 0 (está vacío). La base temporal se borra al terminar.

## 4. Escalado exacto de porciones (RN-14, S4-02)

- **Fórmula:** cantidad escalada = `quantity_milli × servings_milli / (base_servings × 1000)`, con `servings_milli` = porciones × 1000 (1 porción = 1000; media = 500).
- El resultado se guarda como **fracción exacta reducida** (`Qty{Num, Den}`); si es entera, `Den = 1`. Nunca se redondea dentro del motor: el redondeo a milésimas solo es para mostrar (HALF_UP) y el redondeo a paquetes solo ocurre en la lista de compras (§6).
- La unidad no cambia. Porciones admitidas: cualquier entero ≥ 1 y fracciones de porción en milésimas.
- **Ejemplo (RN-14):** 250 g (250000) para 2 → 3 porciones = 375000 exacto. **Ejemplo ILUSTRATIVO de fracción:** 7 ml (7000) para 2 → 3 porciones = 21000/2 = 10500 exacto; 1 pieza (1000) → 3 porciones = 1500 (pieza y media).
- **Prueba:** `TestEscalado*`: las 60 recetas × todos sus ingredientes × (2→4, 2→1, 2→3, 2→6) contra una tabla independiente calculada con `fractions.Fraction` en Python (`oraculo/fixtures/s4/s4-escalado-recetario.json`, generada por el agente del recetario); 0 diferencias. 2→4 debe dar exactamente el doble en todos.

## 5. Plan semanal por persona (RN-14, RN-19, K8, K38, K43; S4-05)

### 5.1 Celdas
- Una celda es día × tiempo × persona. Cada renglón (`PlanEntry`) es receta + porciones, o «fuera».
- **Solo personas con plan** (Axel y Ximena). Una entrada para una persona externa (Karen) devuelve `ErrPersonWithoutPlan`; la base además lo rechaza con su disparador.
- La colación (`snack`) es opcional: vacía no cuenta como hueco; llena tiene costo y kcal como cualquier tiempo.
- **«Fuera» (RN-19):** sin receta, sin porciones, sin kcal y sin costo de recetas; el gasto lo registra la API en la categoría `eating_out` con la persona (no es trabajo de S4). S4 lo reporta como `eating_out_slots` (C1 `TodayKcal`).
- Una celda con «fuera» no admite además una receta: `ErrEatingOutConflict` (lectura de RN-19 «ese tiempo se marca fuera»; supuesto S-1 del reporte, no cambia cifras).
- Varias recetas en la misma celda están permitidas (índice `ux_meal_plan_entry_slot`); cada una con sus porciones.

### 5.2 Semana por omisión y huecos
- `DefaultWeek(today)`: el lunes de la semana de `today`; **si `today` es domingo, el lunes siguiente** (K38, FIRME). Toda escritura de la lista lleva semana explícita (C1 R3-01).
- `UnplannedSlots(week)`: número de pares día × tiempo obligatorio (desayuno, comida, cena) en que **alguna** persona con plan no tiene receta ni «fuera»; máximo 21 (K43, C1 `unplanned_slots`).

### 5.3 Salida (C1 `MealPlanWeek`)
| Campo C1 | Cálculo |
|---|---|
| `cells[].cost_cents` | §8.2; `null` si falta algún precio (no 0) |
| `cells[].kcal` | §9.3; `null` en «fuera» o si falta un dato nutricional |
| `cells[].price_store_id`, `price_observed_date` | tienda de la base de precio si todos los ingredientes salieron de la misma; si son varias, `price_store_id` ausente y la fecha es la **más vieja** de los precios usados (supuesto S-2, no cambia cifras) |
| `days[].cost_cents` | Σ de `cells[].cost_cents` del día (S4-06: «Σ de las celdas»); `null` si alguna celda es `null` |
| `days[].persons[].kcal` | Σ de `cells[].kcal` de la persona en el día; «fuera» no suma |
| `days[].persons[].daily_kcal_target`, `difference_kcal` | meta que rige (§10.4) y kcal − meta |
| `week_cost_cents` | Σ de `days[].cost_cents` |
| `unplanned_slots` | §5.2 |
| `leftovers[]` | §6.4 (PROVISIONAL, bloqueo `S4-reglas-lista`) |

### 5.4 Ingredientes opcionales
El recetario de Claude no lleva opcionales. En recetas de la casa, un ingrediente opcional **no entra** en la lista, el costo ni las kcal y se muestra aparte como «opcional» (supuesto S-3; ninguna cifra canon lo usa).

## 6. Lista de compras neteada con inventario híbrido (RN-15, K1, S4-03)

### 6.1 Entradas
Semana (lunes), entradas del plan de esa semana, recetas, catálogo (productos, presentaciones, grupos), inventario (`InventoryItem` de productos y de sobrantes), presentaciones, tiendas y cantidades elegidas por personas en `shopping_list_item` (se conservan), los artículos `manual` (§6.6) y la recomendación de S5 (§8.1).

### 6.2 Cálculo por producto `p` (todo en milésimas de la unidad base, fracciones exactas)
1. **Necesidad de recetas** `N(p)` = Σ de cantidades escaladas (§4) de `p` en las entradas con receta, **después de netear sobrantes** (6.4). Ingredientes opcionales fuera (§5.4).
2. **Existencia** `E(p)`: en `stock`, `quantity_milli` vigente; en `consumption_rate`, la existencia estimada al **lunes** de la semana con la fórmula de RN-15: `contada − HALF_UP(rate_quantity_milli × días_desde_conteo / rate_period_days)`, nunca menor que 0.
3. **Reposición por ritmo** `R(p)`: PROVISIONAL (bloqueo `S4-reglas-lista`, pregunta 2), opción A: si el artículo es `consumption_rate`, tiene `reorder_level_milli` y su existencia estimada al **domingo** de la semana es ≤ ese nivel, `R(p)` = consumo estimado de una semana = `HALF_UP(rate_quantity_milli × 7 / rate_period_days)`; si no, 0. Motivo `replenishment`.
4. **Necesidad neta** `needed_milli(p)` = `max(0, N(p) + R(p) − E(p))`, llevada a milésima entera **hacia arriba** (nunca se compra de menos; diferencia < 0.001 de unidad).
5. **Paquetes** `packages_to_buy` = techo(cantidad a comprar / `size_milli`), con cantidad a comprar = `quantity_milli` si la persona la fijó y `needed_milli` si no (§6.6), de la presentación que rige: la elegida por la persona; si no eligió presentación, la de la **cotización que rige** (§8.1; sin ninguna elección, la recomendada por S5; si eligió solo tienda, la que `QuoteFor` elige dentro de esa tienda, `S5-interfaz.md` §14.1); si no hay cotización, la `presentacion_id` del ingrediente.
6. **Sobrante de paquete previsto** (`leftover_milli`) = `max(0, packages_to_buy × size_milli − needed_milli)`: lo que queda en existencia después de que el plan consume su necesidad, incluido lo que la persona decidió comprar de más; si fijó menos que la necesidad, 0 (sin negativos, S4-03; §6.6). Al confirmar la compra entra como existencia (movimiento `purchase` por el paquete completo; el consumo del plan lo descuenta), y la lista siguiente lo resta en el paso 2.
7. Productos con `needed_milli = 0` y sin elección de persona no aparecen como artículo calculado (no hay cantidades negativas).

### 6.3 Persistencia (la hace la API con lo que devuelve S4)
Un artículo `calculated` por producto (`ux_shopping_list_item_calculated`), con `needed_milli`, `reasons` (`recipe`, `replenishment`) y `recipe_ids`. Regenerar actualiza `needed_milli` y **conserva** `quantity_milli`, `store_id`, `presentation_id`, `is_bought` y los artículos `manual` (C1 `generateShoppingList`, invariante 13 de S1).

### 6.4 Sobrantes de receta
- **Consumo de sobrantes registrados (FIRME, S4-03):** un `InventoryItem` con `recipe_id` y porciones > 0 cubre primero las entradas del plan de esa receta en orden de fecha y tiempo, solo si `expires_date` es nulo o ≥ la fecha de la entrada; las porciones cubiertas no generan ingredientes en `N(p)`.
- **Sobrante previsto (`MealPlanWeek.leftovers`)**: PROVISIONAL (bloqueo `S4-reglas-lista`, pregunta 1), opción A: se cocina exactamente lo planeado, así que el plan no genera sobrantes nuevos; `leftovers` lista los sobrantes **registrados** que el plan de la semana no alcanza a consumir (receta y porciones restantes). Los sobrantes reales se registran con `createInventoryItem` (movimiento `recipe_leftover`).

### 6.5 Agrupación por canal (S4-12)
Cada artículo hereda la vertical de su grupo (`product_group.purchase_channel`); un producto sin grupo es error de catálogo (`ErrProductWithoutChannel`) y S4-12 cuenta 0. La tienda de cada artículo (`items[].store`) se fija así, en este orden:
1. **Tienda elegida** por la persona (`store` no nulo en `list_choices`; para el artículo calculado, la elección del artículo `calculated` del producto): esa, en cualquier vertical (C1: «el elegido»).
2. **Vertical `chain`** sin tienda elegida: la tienda de la cotización que rige (§8.1: la recomendada, `recommended = true`, si la persona no eligió nada; la de la presentación elegida si eligió solo presentación); sin cotización, `store` nulo y el producto en `missing_prices[]`.
3. **Vertical `market`** sin tienda elegida: PROVISIONAL (bloqueo `S4-asignacion-mercado`), opción A. RN-16 (FIRME) dice que el renglón «se asigna a un mercado que surte su grupo, en su día de visita», pero no dice cuál entre varios, qué pasa si no hay ninguno ni si gana una cadena más barata:
   - candidatos: los `stores` con `kind = market` cuyo `groups[]` incluye el grupo del producto (`store_product_group`);
   - con uno o más: el de `next_visit` más cercana (la fecha la calcula S5 §9; §12.1); en empate, el primero en el orden de `stores` (la API los entrega por `id` ascendente);
   - sin ninguno (el caso del arranque: «Axel aún no da de alta sus mercados», RN-16): `store` nulo dentro de la vertical `market`; la pantalla dice «sin mercado dado de alta»;
   - que S5 recomiende una cadena más barata para ese producto (S5 §7.4, `productos-formato.md` §3) **no mueve** el artículo: `items[].store` sigue siendo el mercado (o nulo), y la recomendación de S5 por precio sale en `items[].recommended_store` (C1 `recommended_store_id`) y en la comparativa (`compare_basket`). La persona puede elegir esa cadena en la lista, y entonces rige por el paso 1;
   - **esto exige cambiar C1 v1.2.** C1 dice que `recommended_store_id` «rige mientras `store_id` esté ausente» y que `groups[].store_id` es «el elegido o el recomendado». Con esta opción, y también sin ningún mercado (`store` nulo), la tienda asignada y la recomendada pueden ser distintas. Ejemplo: uva verde sin elección, un tianguis de martes surte el grupo y S5 recomienda Bodega Aurrera; el artículo va en el grupo del tianguis y su precio es el de Bodega. La traducción PROVISIONAL y la solicitud a C1 están en §12.2.

La tienda solo decide el agrupamiento: el precio y la presentación del artículo salen siempre de la cotización que rige (§8.1), aunque sea de otra tienda.

**Precio y recomendación por artículo** (de aquí salen los campos calculados de C1 `ShoppingListItem`, §12.2):
- `price_store`, `price_observed_date` e `is_stale_price`: la tienda, la fecha y la marca de viejo de la cotización que rige (§8.1). `estimated_cents` = `packages_to_buy × price_cents` de esa cotización. Los cuatro son `null` si el artículo está «sin precio».
- `recommended_store` = `price_store` si la persona no eligió tienda; si eligió, la tienda de la cotización recomendada (solo informativa; `null` si no hay).
- `recommended_presentation` = `presentation` (la que rige, §6.2 paso 5) si la persona no eligió presentación; si eligió, la de la cotización recomendada (solo informativa; `null` si no hay).

Así, mientras la persona no elija, la recomendación es exactamente la tienda y la presentación con que S4 calcula paquetes y precio: una sola presentación recomendada por artículo (`S5-interfaz.md` §7.3 y §14.1).

**`groups[]`** (C1 `ShoppingList.groups`): un grupo por cada par distinto (`vertical`, `store`) de los artículos; `visit_date` = `next_visit` de la tienda si es mercado (si no, `null`); `estimated_total_cents` = Σ `packages_to_buy × price_cents` de la cotización que rige de sus artículos con precio, `calculated` y `manual`, comprados o no (§6.6) (entero exacto: los paquetes son de la presentación de esa cotización y de la cantidad a comprar, §6.2 paso 5); los artículos sin precio no suman y se listan en `missing_prices[]`. Orden: primero `market` y luego `chain`; dentro de cada vertical, por `visit_date` ascendente (nulo al final) y luego por el orden de `stores`, con `store` nulo al final. El orden y el total con artículos sin precio son parte de la opción A del mismo bloqueo (pregunta 4). El total de la lista, `estimated_total_cents` en la raíz, es Σ de los `estimated_total_cents` de los grupos (C1 `ShoppingList.total_cents`); los nombres de C1 están en la traducción de §12.2.

### 6.6 Cantidad fijada por una persona y artículos `manual` (S4-03, S4-06; C1 `quantity_milli`)
C1 v1.2 define `ShoppingListItem.quantity_milli` como la «cantidad a comprar que fijó una persona (N-152, N-153)… NULL en calculated = comprar la necesidad», y `ShoppingListItemCreate` la exige en todo artículo `manual`. De ahí salen estas reglas, iguales para los dos tipos de artículo:
- **Cantidad a comprar** = `quantity_milli` si no es nula; si es nula, `needed_milli`. Un artículo `manual` siempre trae `quantity_milli` (si viene `null`, `invalid_argument`).
- `packages_to_buy` = techo(cantidad a comprar / `size_milli`) de la presentación que rige (§6.2 paso 5), y `estimated_cents` = `packages_to_buy × price_cents` de la cotización que rige (§6.5). La cantidad fijada cambia los paquetes y el costo, pero no la presentación ni la tienda: esas siguen §6.2 paso 5, §6.5 y §8.1.
- `leftover_milli` = `max(0, packages_to_buy × size_milli − needed_milli)` (§6.2 paso 6), con `needed_milli` = 0 en un artículo `manual`. Es lo que queda en existencia después del consumo del plan.
- **Artículo `calculated` con `quantity_milli`:** conserva su `needed_milli`, `reasons[]` y `recipes[]`. Si su necesidad baja a 0, sigue apareciendo, porque la cantidad fijada es una elección (§6.2 paso 7).
- **Artículo `manual`** (un renglón `kind = manual` de `list_choices`, siempre con producto del catálogo): es un artículo propio aunque el producto también tenga artículo `calculated`. Lleva `needed_milli` = `null` (C1: no tiene necesidad calculada), `reasons[]` = [`manual`] (C1 `reasons`) y `recipes[]` vacía. Su tienda y su grupo siguen §6.5 con la vertical de su producto, y su cotización que rige es la de su propia elección (§8.1). Si no eligió presentación y no hay cotización, no hay presentación de respaldo, porque no viene de un ingrediente: `presentation`, `packages_to_buy`, `leftover_milli` y los cuatro campos de precio salen `null`, y el producto va en `missing_prices[]`.
- **Totales.** Los artículos `calculated` y `manual` con precio suman su `estimated_cents` en `groups[].estimated_total_cents` y, por esa vía, en `estimated_total_cents` de la raíz. `is_bought` no cambia ningún cálculo: es una marca (C1: «No crea gasto ni precio hasta confirmar el ticket»), y U1 N-150 pide «Totales por grupo y total general» sin excluir lo comprado. Ni C1 ni U1 tienen un subtotal de lo que falta por comprar.
- **Texto libre.** Un artículo `manual` de texto libre (C1 `free_text`) no tiene producto, grupo ni precio: no entra en `list_choices` y S4 no lo calcula.

Ejemplo ILUSTRATIVO (caso `s4-lista-cantidad-manual`). Cotizaciones recomendadas de Walmart del 2026-10-10: arroz 1 kg a 2800, frijol 1 kg a 3000 y aceite de soya 946 ml a 3990. Las servilletas (por pieza, grupo de vertical `chain`) no tienen cotización.

| Artículo | `needed_milli` | `quantity_milli` | Cantidad a comprar | `packages_to_buy` | `leftover_milli` | `estimated_cents` |
|---|---|---|---|---|---|---|
| arroz, `calculated` | 500000 | 3000000 | 3000000 | 3 | 2500000 | 8400 |
| frijol, `calculated` | 1200000 | 500000 | 500000 | 1 | 0 | 3000 |
| aceite de soya, `manual`, `is_bought` = true | `null` | 1892000 | 1892000 | 2 | 1892000 | 7980 |
| servilletas, `manual`, sin cotización | `null` | 1000 | 1000 | `null` | `null` | `null` (en `missing_prices[]`) |

Grupos: (`chain`, `walmart`) = 8400 + 3000 + 7980 = **19380**; (`chain`, `null`) = 0, porque solo tiene las servilletas sin precio. Total de la raíz: **19380**. Sin la cantidad fijada, el arroz llevaría 1 paquete (2800) y un sobrante de 500000, y el frijol 2 paquetes (6000) y un sobrante de 800000.

## 7. Inventario híbrido (RN-15, K7; S4-04)

| Salida C1 (`InventoryItem`) | Cálculo |
|---|---|
| `estimated_quantity_milli` | `stock`: `quantity_milli`. `consumption_rate`: fórmula de RN-15 a «hoy». Ejemplo RN-15: 2838000 contados, ritmo 946000 cada 21 días, 30 días después → 2838000 − HALF_UP(946000 × 30 / 21) = 2838000 − 1351429 = **1486571** |
| `status` | `out` si la existencia estimada es 0; si no, `reorder` si hay `reorder_level_milli` y existencia ≤ nivel; si no, `sufficient`. **Agotado prevalece** (con nivel ≥ 0, toda existencia 0 también cumple «≤ nivel»; si reponer ganara, agotado nunca aparecería) |
| `weekly_consumption_milli` | `consumption_rate`: `HALF_UP(rate_quantity_milli × 7 / rate_period_days)`; `stock`: ausente |
| `days_left` | `consumption_rate`: piso(`estimada × rate_period_days / rate_quantity_milli`) días completos cubiertos (supuesto S-4); `stock`: ausente |
| Sobrantes (`recipe_id`) | siempre `stock`, en porciones × 1000 |

**Constructores de movimientos** (la API los guarda en la misma transacción que cambia la existencia, U1 N-181 y N-183): `PurchaseMovement` (+paquetes × contenido), `PlanConsumptionMovement` (− cantidad escalada, sin bajar de 0; el faltante se reporta), `LeftoverMovement` (+porciones), `AdjustmentMovement` (motivo `count`, `consumption` o `waste`; la diferencia contra la cantidad nueva). Qué operación de C1 dispara el consumo del plan lo define C1/U1; S4 solo entrega el constructor. `InventorySummary` = conteo por `status`.

## 8. Costo por receta, día y semana (S4-06)

### 8.1 Entrada de precios desde S5 (`PriceSource`)
S4 recibe de S5, por producto y fecha de referencia, cotizaciones vigentes que compiten según S5 (`S5-interfaz.md` §5 a §7 y §14 `QuoteFor`): `Quote{ProductID, StoreID, PresentationID, PriceCents, SizeMilli, ObservedDate, IsStale, Recommended}` o «sin precio». Por producto llega **exactamente una** recomendada (`Recommended = true`: `QuoteFor` sin tienda ni presentación, el par tienda-presentación de menor precio por unidad base ajustado con el desempate de `S5-interfaz.md` §14.1, que no depende de la lista) y, además, la de cada tienda o presentación que la persona eligió en la lista (`QuoteFor` con esa elección; con solo tienda, la misma regla elige la presentación dentro de esa tienda). La API las pide con `date` = «hoy» en America/Mexico_City (`today`, S5 §14.1), también para una semana futura. Si una cotización por elección coincide en tienda y presentación con la recomendada, la API manda solo la recomendada. S4 no compara precios ni repite el peso de S5 §7.1. **Cotización que rige**, una sola regla para `recipe_detail`, `week_plan` y `shopping_list`:
1. Si la persona eligió tienda o presentación para el producto (`list_choices` con `store` o `presentation` no nulos; para costos de recetas, del plan y del artículo calculado cuenta la elección del artículo `calculated`; un artículo `manual` usa la suya): la cotización cuya tienda es la elegida (si eligió tienda) y cuya presentación es la elegida (si eligió presentación). Si casan varias, la recomendada si está entre ellas; si siguen varias, `invalid_argument`. Si no casa ninguna, el producto queda «sin precio»: no se cae a la recomendada, porque la persona eligió (§8.2).
2. Si no eligió nada: la recomendada.
3. Sin cotizaciones del producto: «sin precio».

Las cotizaciones que no son la recomendada solo sirven para atender las elecciones del paso 1; si ninguna elección las pide, se ignoran. PROFECO y `site_assisted` nunca llegan aquí como base de costo (K22, K37).

### 8.2 Fórmulas
- Costo exacto de un ingrediente = `cantidad_escalada_milli × PriceCents / SizeMilli` (fracción exacta; precio por unidad base sin redondear: RN-20).
- `RecipeDetail.ingredients[].cost_cents` = HALF_UP del costo exacto del ingrediente (para mostrar).
- `RecipeDetail.total_cost_cents` = HALF_UP(Σ costos exactos) — un solo redondeo; la pantalla no suma los renglones mostrados.
- `cost_per_serving_cents` = HALF_UP(Σ exacto ÷ porciones pedidas).
- Celda del plan = HALF_UP(costo exacto de la receta para `base_servings` × `servings_milli` / (`base_servings` × 1000)).
- Día = Σ celdas; semana = Σ días (enteros, sin redondeo adicional).
- **Sin precio:** el ingrediente queda «sin precio» (`cost_cents` ausente, nunca 0.00) y el total de la receta, la celda, el día y la semana quedan sin valor, con la lista de productos sin precio (`ErrNoPrice` no es error fatal: es estado).

### 8.3 Ejemplo ILUSTRATIVO con precios reales del ticket del 10-oct-2026 (Bodega Aurrera)
Receta de 2 porciones: espagueti 250 g (Spaghetti La Moderna 250 g a 12.00, `data.json: items_identicos[21]`), atún 140 g (21.00, `[16]`), granos de elote 110 g (lata 220 g a 13.50, `[15]`), aceite de oliva 15 ml (500 ml a 150.00, `[0]`): 1200 + 2100 + 675 + 450 = **4425** centavos ($44.25); por porción HALF_UP(2212.5) = **$22.13**; para 4 porciones **$88.50**. (Verificado con `fractions.Fraction`; la receta es ilustrativa, los precios son del ticket.)

## 9. Kilocalorías con fuente citada (RN-17, S4-08, S4-10, S4-11)

### 9.1 Tabla `datos/semilla/nutricion/tabla.csv`
UTF-8 sin BOM, separador coma, encabezado exacto:

`producto_id,kcal_100g,proteina_g,grasa_g,carbohidratos_g,fuente,id_fuente,nombre_fuente,fecha,estado,confirmo,fecha_confirmacion`

| Columna | Regla | Columna en la base |
|---|---|---|
| `producto_id` | id de `productos.json`; uno por renglón y **único** en el archivo (el valor preferido) | `nutrition_fact.product_id`, `is_preferred = 1` |
| `kcal_100g` | decimal con punto, hasta 3 decimales, ≥ 0, obligatorio; se lee como texto y se pasa a milésimas sin flotantes (`884` → 884000; `52.5` → 52500) | `kcal_milli_per_100g` |
| `proteina_g`, `grasa_g`, `carbohidratos_g` | opcionales (vacío = sin dato), hasta 3 decimales → mg | `protein_mg_per_100g`, `fat_mg_per_100g`, `carbs_mg_per_100g` |
| `fuente` | `USDA`, `OFF` o `manual` | `usda_fdc`, `open_food_facts`, `manual` |
| `id_fuente` | USDA: `fdc_id` (dígitos); OFF: código de barras 8 a 14 dígitos; manual: cita textual obligatoria (tabla y página, o etiqueta y presentación) | `source_ref` |
| `nombre_fuente` | nombre tal como viene en la fuente (puede estar en inglés); **solo procedencia**, nunca se muestra (S4-11, riesgo R9) | `notes` del renglón nutricional no existe: queda solo en el CSV |
| `fecha` | ISO, fecha de extracción o captura | — (procedencia) |
| `estado` | `propuesto` (empate que sugirió el agente) o `confirmado` (una persona lo revisó contra la fuente); obligatorio (R9, S4-11) | — (procedencia: `nutrition_fact` no tiene columna de estado; a la base solo entra lo confirmado) |
| `confirmo` | código de la persona que confirmó (`axel` o `ximena`); vacío si `propuesto`, obligatorio si `confirmado` | — (procedencia) |
| `fecha_confirmacion` | ISO; vacía si `propuesto`; obligatoria si `confirmado` y no anterior a `fecha` | — (procedencia) |

**Licencias** (S4-10) en `datos/semilla/nutricion/PROCEDENCIA.md`: USDA FoodData Central Foundation Foods 04/2026, CC0, carga local del CSV; Open Food Facts, ODbL + DbCL, volcado local filtrado; tablas mexicanas (SMAE) solo por captura `manual` citando la fuente (licencia sin verificar). `source_license` se llena con `CC0` u `ODbL` y queda nulo en `manual`. Todo empate producto ↔ fuente es **asistido con confirmación humana** (R9); ninguno entra automático: el agente de nutrición escribe sus empates como `propuesto` y solo una persona los pasa a `confirmado`. **Quién confirma y cuándo:** PROVISIONAL (bloqueo `S4-confirmacion-empates`), opción A: Axel en la puerta 2, sobre la lista de renglones `propuesto` que entrega el agente de nutrición.

**Cargador** `LoadNutritionCSV(path, catalog, mode)`: rechaza encabezado distinto, producto inexistente, decimales con más de 3 cifras, `id_fuente` vacío, duplicados, `estado` fuera de sus dos valores y `confirmo` o `fecha_confirmacion` incoherentes con el `estado`. `mode`: **`production`** (carga a la base de la app): un renglón `propuesto` **no se carga** (se rechaza con motivo `unconfirmed`) y su producto queda «sin dato»; **`review`** (pruebas Go, fixtures y medición de WF-1): carga los dos estados y marca los propuestos. Devuelve `Coverage{Unique, WithSource, Confirmed, Proposed, Missing []ProductCode}` sobre los ingredientes del recetario: `WithSource = Confirmed + Proposed = Unique − len(Missing)` y 0 filas con fuente vacía (S4-08). S4-08 se reporta en modo `review` con los dos conteos (con valor y fuente; de ellos, confirmados); S4-11 («ninguno entra automático») lo garantiza el modo `production`.

### 9.2 Gramos de un ingrediente
`g`: la cantidad escalada. `ml` o `pieza`: cantidad × `grams_per_unit_milli` / 1000 (densidad o peso promedio del catálogo, con su fuente en `productos.json`). Sin `grams_per_unit_milli` en ml o pieza → «sin dato» (no 0).

### 9.3 Fórmulas (milésimas de kcal, exactas; un redondeo al final)
- kcal exactas del ingrediente = gramos_milli × `kcal_milli_per_100g` / 100000 → milésimas de kcal.
- Dato nutricional que rige: el preferido de la presentación del ingrediente si existe; si no, el preferido genérico del producto (supuesto S-5).
- `total_kcal` = HALF_UP(Σ exacto / 1000); `kcal_per_serving` = HALF_UP(Σ exacto / (1000 × porciones pedidas)) (S4-08: Σ (kcal/100 g × cantidad/100) / porciones).
- Celda = HALF_UP(Σ exacto de la receta base × `servings_milli` / (`base_servings` × 1000 × 1000)); día = Σ celdas (RN-17).
- `% de la meta` = HALF_UP(kcal del día × 1000 / meta) / 10 con 1 decimal (RN-17: 1,850 contra 2,000 = 92.5 %).
- Un ingrediente sin dato deja la receta como «incompleta»: `kcal` ausente y lista `missing_nutrition`; nunca 0 ni valor inventado.

## 10. Salud: IMC y meta de calorías (RN-17, K7, K23, K36; S4-07)

### 10.1 Edad e IMC
- Edad = años cumplidos a «hoy» (America/Mexico_City) desde `birth_date`.
- `bmi_tenths` = HALF_UP(`weight_g` × 10000 / `height_mm`²) con el último peso (`weight_entry` más reciente). Ejemplo RN-17: 72500 g y 1750 mm → 236.73 → **237** (23.7).
- `bmi_class` (K7: 18.5 / 25 / 30): `underweight` < 18.5 ≤ `normal` < 25.0 ≤ `overweight` < 30.0 ≤ `obesity`. Sobre qué valor se clasifica (exacto o el redondeado a un decimal) es PROVISIONAL (bloqueo `S4-reglas-lista`, pregunta 3), opción A: **sobre `bmi_tenths`**, el valor que se muestra. Solo difiere en la franja de ±0.05 de un corte (73990 g y 2000 mm = 18.4975 exacto, se muestra 18.5: A → `normal`, B → `underweight`).

### 10.2 Meta sugerida (FIRME)
- MSJ en milésimas de kcal: `10 × weight_g + 625 × height_mm − 5000 × edad + 5000` (hombre) o `− 161000` (mujer).
- Factor de actividad exacto: `sedentary` 6/5, `light` 11/8, `moderate` 31/20, `active` 69/40, `very_active` 19/10 (RN-17, columna `activity_level`; los 5 niveles con fuente en RN-17, sin bloqueo `S4-meta-ajustes`).
- Ajuste: `lose` → − (`AdjustLoseKcal` de la persona o, si es NULL, `HouseAdjustLoseKcal`); `gain` → + (`AdjustGainKcal` o `HouseAdjustGainKcal`); `maintain` → 0 aunque los ajustes valgan algo. Valores de la casa iniciales 0 y 0 (K36, K44).
- `suggested_kcal_target` = HALF_UP((MSJ × factor ± ajuste × 1000) / 1000), un solo redondeo al final.
- Falta cualquier insumo (sexo, nacimiento, estatura, peso, actividad u objetivo) → sin sugerencia y lista `missing_inputs`.

### 10.3 Casos de las fixtures (personas ficticias, ILUSTRATIVO; nunca datos de Axel o Ximena)
| Caso | MSJ | × factor (exacto) | Meta sugerida |
|---|---|---|---|
| Hombre 30 años, 80 kg, 175 cm, `sedentary`, `maintain` | 1,748.75 | 2,098.50 | **2,099** |
| Mujer 28 años, 60 kg, 165 cm, `sedentary`, `maintain` | 1,330.25 | 1,596.30 | **1,596** |
| Mismo hombre, `lose`, ajustes de la casa 0 y 0 | 1,748.75 | 2,098.50 − 0 | 2,099 |
| Misma mujer, `gain`, ajustes 0 y 0 | 1,330.25 | 1,596.30 + 0 | 1,596 |
| Hombre `lose` con bajar = 300 (fixture, no canon) | | 1,798.50 | **1,799** |
| Mujer `gain` con subir = 250 (fixture) | | 1,846.30 | **1,846** |
| `maintain` con ambos ajustes = 300 | | 2,098.50 | 2,099 |
| RN-17: hombre 35 años, 72.5 kg, 175 cm, `moderate`, `maintain` | 1,648.75 | 2,555.5625 | **2,556**; con `lose` y 500 propios: **2,056** |

Los criterios S4-07 citan los valores antes del redondeo final (2,098.50; 1,596.30; 1,798.50; 1,846.30); la salida entera es la de la última columna (HALF_UP). La prueba verifica ambos.

**IMC de las fixtures** (estatura 2000 mm, ILUSTRATIVO): 74000 g → 18.5 `normal`; 100000 g → 25.0 `overweight`; 120000 g → 30.0 `obesity`; 72500 g y 1750 mm → 23.7 `normal`.

### 10.4 Meta que rige
`daily_kcal_target` confirmado o escrito por la persona → `TargetStatus = confirmed`; si es NULL y hay sugerencia → la sugerida con `TargetStatus = suggested` (la pantalla la muestra como «sugerida»); sin ninguna → `none`. Una meta escrita reemplaza a la sugerida. Rango 500 a 10000 del esquema (C-11).

## 11. Canal de compra por ingrediente (S4-12)
El canal sale del grupo del producto (`product_group.purchase_channel`: `market` o `chain`). `vertical` = `ambas` en `productos.json` no se guarda: significa que el grupo se puede ofrecer a mercados en Ajustes (`store_product_group`), y S5 deja competir a los mercados que declaran el grupo (ver `productos-formato.md` §3). La tienda y el grupo de cada artículo en la lista siguen §6.5 (vertical `market`: PROVISIONAL, bloqueo `S4-asignacion-mercado`); que una cadena compita en S5 por un producto de mercado no cambia su vertical ni su grupo. Prueba `TestCanal*`: lista agrupada por canal, asignación de mercados de §6.5 y conteo de ingredientes sin canal = 0.

## 12. Fixtures `oraculo/fixtures/s4/` y CLI `s4 dump`

**Formato:** el común de `oraculo/fixtures/FORMATO.md` 1.0 (`fixture/1`, agente `contrato-S2`): un archivo por caso `oraculo/fixtures/s4/<caso>.json` con `formato`, `motor` = `s4`, `caso` (`s4-<tema>[-<variante>]`), `descripcion`, `fuente` (lista de `{tipo, ruta}` con `tipo` ∈ `data.json`, `regla`, `criterio`, `decision`, `manual`; los casos ILUSTRATIVOS usan `manual` y lo dicen en `descripcion`), `provisional` (`null`, o `{"bloqueo": "docs/estado/bloqueos/S4-reglas-lista.md", "que": …}` en los casos que tocan sus preguntas), `entradas.datos` (catálogo, recetas, plan, inventario, perfiles, ajustes de la casa, «hoy» y cotizaciones de precio, en línea o por `$ref` a `oraculo/fixtures/s4/_datos/`; claves exactas en §12.1) y `entradas.consultas` (`{id, funcion, argumentos}`), y `esperado` = `{id: salida}`. Valores (FORMATO §3): dinero en centavos enteros, cantidades en milésimas enteras, kcal enteras, fechas ISO, semana `AAAA-Www`, fracciones exactas como `{"num": n, "den": d}` reducidas; claves de salida en inglés con `_cents` y `_milli`.

**Funciones consultables** (`funcion`): `scale_recipe` (§4), `default_week` y `unplanned_slots` (§5.2), `week_plan` (§5.3), `shopping_list` (§6), `inventory_status` (§7), `recipe_detail` (§8 y §9: costo y kcal escalados), `nutrition_coverage` (§9.1), `health_profile` (§10: edad, IMC, clase, meta sugerida y meta que rige). Argumentos y salida de cada una en §12.2.

| Caso | Fuente | Qué fija |
|---|---|---|
| `s4-escalado-2a4` | RN-14, S4-02 | receta con 250000 g, 7000 ml y 1000 pieza → ×2, ×1/2, ×3/2, ×3 exactos |
| `s4-escalado-recetario` | S4-02 | las 60 recetas × 4 escalas (la genera el agente del recetario con Python `fractions`) |
| `s4-lista-sin-existencia` | S4-03 | Σ ingredientes de 2 recetas en 3 celdas; paquetes y sobrante de paquete |
| `s4-lista-existencia-parcial` | S4-03 | mismo plan con existencia parcial de 2 productos |
| `s4-lista-existencia-mayor` | S4-03 | existencia mayor que la necesidad → 0 a comprar, sin negativos |
| `s4-lista-sobrantes` | S4-03 | sobrante registrado de 1 porción cubre la primera entrada; uno caducado no cubre |
| `s4-lista-cantidad-manual` | S4-03, S4-06, C1 `quantity_milli`, ILUSTRATIVO | el ejemplo de §6.6: cantidad fijada mayor y menor que la necesidad (paquetes, sobrante y costo), artículo `manual` comprado que suma, artículo `manual` sin cotización (`null` y `missing_prices[]`), `groups[]` y total de la raíz |
| `s4-inventario-modos` | RN-15, K7 | `stock`, `consumption_rate` (1486571 del ejemplo de RN-15), sobrante; `reorder` en igualdad y `out` en 0 |
| `s4-plan-celdas` | RN-14, RN-19, K43 | receta, «fuera» (kcal ausente), colación vacía y llena, Karen rechazada, huecos = conteo esperado |
| `s4-costo-cotizacion-elegida` | S4-06, §8.1, ILUSTRATIVO | cotizaciones con `recommended`: sin elección rige la recomendada; tienda elegida con cotización y sin ella («sin precio», no cae a la recomendada); presentación elegida; dos recomendadas del mismo producto → `invalid_argument`; en `shopping_list`, `recommended_store`, `recommended_presentation`, `price_store` y `estimated_cents` sin elección, con solo tienda y con solo presentación elegidas (§6.5) |
| `s4-canal-mercados` | S4-12, RN-16, ILUSTRATIVO; `provisional` (`S4-asignacion-mercado`) | dos mercados que surten el grupo (gana la `next_visit` más cercana; empate por el orden de `stores`), ninguno (`store` nulo en `market`), cadena elegida por la persona para un producto de mercado, orden y `estimated_total_cents` de `groups[]` con un artículo sin precio; un artículo de mercado cuyo `recommended_store` (la cadena de la cotización recomendada) no es su `store`; `estimated_total_cents` de la raíz = Σ de los grupos |
| `s4-costo-recetas` | S4-06, data.json | 3 recetas con precios de Bodega del 10-oct (incluye §8.3: 4425 / 2213 / 8850), 1 día, 1 semana, un ingrediente sin precio |
| `s4-kcal-recetas-autor` | S4-08 | 5 recetas que elige el **autor**, con los valores de `tabla.csv` (modo `review`) y su `id_fuente`; el esperado lo calcula `s4_run_fixtures.py` con `fractions`. Es prueba propia del autor y **no** sustituye a S4-09 |
| `s4-salud-perfiles` | S4-07, RN-17 | las filas de §10.3 y los 4 IMC |

**S4-09 (las 5 recetas las elige el revisor):** el revisor de recetas y precios elige, al revisar, 5 de las 60 recetas del recetario (no las elige el autor ni las fija este contrato), recalcula a mano desde los valores de la fuente citada en `tabla.csv` las kcal por porción y el escalado 2→4, y lo anota receta por receta en `docs/estado/wf1/<revisor>.md`, sección «S4 verificación manual». Si quiere dejarlo ejecutable, agrega sus 5 recetas en una fixture propia `oraculo/fixtures/s4/s4-kcal-revisor.json` (fuente `{"tipo": "criterio", "ruta": "S4-09"}`), que corre con el mismo `dump` y `compara.py`.

**CLI** (FORMATO §4): `cd backend && go build -o ../datos/tmp/motores ./cmd/motores && cd .. && datos/tmp/motores s4 dump --fixtures oraculo/fixtures/s4 --out datos/tmp/motor-s4.json [--caso <caso>]`. Escribe `{"formato": "salida/1", "motor": "s4", "casos": {…}}` con cada salida recortada a las claves de `esperado`, determinista, sin red ni base de datos; códigos de salida 0, 1 y 2 de FORMATO §4. La comparación es `python3 -I oraculo/compara.py s4` (FORMATO §5). El oráculo es el de §12.3; el subcomando `s4` vive en el binario común `backend/cmd/motores` que define FORMATO §4.

### 12.1 `Input` de S4 (raíz de `entradas.datos`)
Claves en inglés y snake_case; dinero en `_cents`, cantidades en `_milli` (milésimas de la unidad base, C-11), fechas ISO. Toda lista puede venir vacía. Las entidades se citan por su **código** (texto: los `id` de `productos.json` y `recetario.json`, o el que escriba la fixture), nunca por id numérico; la API traduce códigos e ids. Un campo desconocido, en la raíz o dentro de cualquier objeto, es error `unknown_field` en el motor Go y en el oráculo; un campo obligatorio ausente es `missing_field`. Las rutas (`*_file`) son relativas a la raíz del repositorio, desde donde se corren `dump` y el oráculo (FORMATO §4), y se leen sin copiarlas.

| Clave | Tipo | Oblig. | Contenido | Origen en la API |
|---|---|---|---|---|
| `today` | fecha | sí | «Hoy» ya resuelto en America/Mexico_City (C-05) | API |
| `catalog_file` | ruta | uno de los dos | `datos/semilla/productos.json`, leído con las reglas de `productos-formato.md` | semilla |
| `catalog` | objeto | uno de los dos | `groups[]` {`code`, `name`, `purchase_channel` (`market` · `chain`)}; `products[]` {`code`, `name`, `base_unit` (`g` · `ml` · `piece`), `group` (código o `null`), `grams_per_unit_milli` (entero o `null`)}; `presentations[]` {`code`, `product`, `name`, `size_milli`} | `product_group`, `product`, `presentation` |
| `recipe_book_file` | ruta | al menos uno | `datos/semilla/recetario.json` (§3) | semilla |
| `recipes` | lista | al menos uno | {`code`, `name`, `slots[]` (`breakfast` · `lunch` · `dinner` · `snack`), `base_servings`, `prep_minutes`, `cook_minutes`, `origin` (`claude_seed` · `household`), `is_active`, `ingredients[]` {`product`, `presentation`, `quantity_milli`, `optional`}}; si viene con `recipe_book_file` se suma a sus recetas y un código repetido es `invalid_argument` | `recipe`, `recipe_meal_slot`, `recipe_ingredient` |
| `nutrition_file` | ruta | no (uno u otro) | `datos/semilla/nutricion/tabla.csv` (§9.1), cargada en modo `review` | semilla |
| `nutrition` | lista | no (uno u otro) | {`product`, `presentation` (o `null`), `kcal_milli_per_100g`, `source` (`usda_fdc` · `open_food_facts` · `manual` · `mx_table` · `store_label`), `source_ref`, `status` (`proposed` · `confirmed`)} | `nutrition_fact` con `is_preferred = 1` |
| `persons` | lista | sí | {`code`, `has_meal_plan`} | `person` |
| `stores` | lista | no | {`code`, `kind` (`chain` · `branch` · `market`), `groups[]` (códigos de grupo que surte), `next_visit` (fecha, obligatoria en `market`; `null` en cadenas y sucursales; si no, `invalid_argument`. En la API la da S5 `NextVisit(store, from)` con `from` = el más tardío entre `today` y el lunes de la semana de la lista: PROVISIONAL, bloqueo `S4-asignacion-mercado`, pregunta 4)}; el orden de la lista desempata la asignación de mercados (§6.5): la API los entrega por `id` ascendente | `store`, `store_product_group` |
| `meal_plan` | lista | no | {`date`, `slot`, `person`, `recipe` (o `null`), `servings_milli` (o `null`), `eating_out`} | `meal_plan_entry` |
| `inventory` | lista | no | {`code`, `product` (o `null`), `recipe` (o `null`; sobrante en porciones × 1000), `mode` (`stock` · `consumption_rate`), `quantity_milli`, `counted_date`, `rate_quantity_milli`, `rate_period_days`, `reorder_level_milli`, `expires_date`}; los cinco últimos pueden ser `null` | `inventory_item` |
| `list_choices` | lista | no | elecciones de la persona en la lista: {`product`, `kind` (`calculated` · `manual`), `store` (o `null`), `presentation` (o `null`), `quantity_milli` (o `null`: cantidad a comprar que fijó la persona, §6.6; en `manual` es obligatoria y `null` es `invalid_argument`), `is_bought`}; cada renglón `manual` es un artículo propio de la lista | `shopping_list_item` |
| `quotes` | lista | no | cotizaciones que entrega S5 por `PriceSource` (§8.1), ya resueltas: {`product`, `store`, `presentation`, `price_cents`, `size_milli`, `observed_date`, `is_stale`, `recommended` (booleano)}; por producto con cotizaciones, **exactamente una** con `recommended = true` (la mejor tienda por ingrediente de S5) y a lo más una por (`product`, `store`, `presentation`); si no, `invalid_argument`. Las demás solo atienden `list_choices` (§8.1). Un producto sin cotización queda «sin precio» | S5 `QuoteFor` con `date` = `today` (S5 §14.1; mismas claves que `quote_for`, S5 §15.2) |
| `health_profiles` | lista | no | {`person`, `height_mm`, `birth_date`, `sex`, `activity`, `goal`, `adjust_lose_kcal`, `adjust_gain_kcal`, `daily_kcal_target`, `dietary_restrictions`}; todo salvo `person` puede ser `null` | `health_profile` |
| `weights` | lista | no | {`person`, `date`, `weight_g`} | `weight_entry` |
| `settings` | objeto | no | {`house_adjust_lose_kcal`, `house_adjust_gain_kcal`}; ausente = 0 y 0 | `app_setting health.goal_adjust_lose_kcal` y `health.goal_adjust_gain_kcal` (K36, K44) |

### 12.2 Funciones consultables: argumentos y salida
Las salidas usan claves en inglés; las entidades salen por su código con la clave sin `_id` (`product`, `presentation`, `recipe`, `person`, `store`) y la API las traduce a los `*_id` de C1. Las cantidades escaladas son fracciones `{"num": n, "den": d}` reducidas (FORMATO §3). Un valor `null` es un estado que no es error (sin precio, sin dato, sin sugerencia). Los errores de §14 salen como `{"error": "<código>"}` con el nombre en snake_case sin `Err` (`ErrPersonWithoutPlan` → `person_without_plan`), además de `unknown_field`, `missing_field`, `invalid_date` e `invalid_argument`.

| `funcion` | `argumentos` | Salida |
|---|---|---|
| `scale_recipe` | `recipe`, `servings_milli` | {`recipe`, `servings_milli`, `ingredients[]` {`product`, `unit`, `quantity` (fracción)}} (§4) |
| `default_week` | `date` | {`week` (`AAAA-Www`), `monday`} (§5.2) |
| `unplanned_slots` | `week` | {`unplanned_slots`} (§5.2) |
| `week_plan` | `week` | {`week`, `cells[]` {`date`, `slot`, `person`, `recipe`, `servings_milli`, `eating_out`, `cost_cents`, `kcal`, `price_store`, `price_observed_date`}, `days[]` {`date`, `cost_cents`, `persons[]` {`person`, `kcal`, `daily_kcal_target`, `difference_kcal`}}, `week_cost_cents`, `unplanned_slots`, `eating_out_slots`, `leftovers[]` {`recipe`, `servings_milli`}} (§5.3) |
| `shopping_list` | `week` | {`week`, `items[]` {`product`, `kind`, `needed_milli` (`null` en `manual`), `quantity_milli` (la elección de la persona o `null`), `presentation`, `packages_to_buy`, `leftover_milli`, `reasons[]`, `recipes[]`, `store`, `is_bought`, `recommended_store`, `recommended_presentation`, `price_store`, `price_observed_date`, `is_stale_price`, `estimated_cents`}, `groups[]` {`vertical`, `store`, `visit_date`, `estimated_total_cents`}, `estimated_total_cents`, `missing_prices[]`} (§6; tienda, orden y total de `groups[]`, y precio y recomendación por artículo, en §6.5; cantidad fijada y artículos `manual` en §6.6; traducción a C1 abajo) |
| `inventory_status` | ninguno (usa `today`) | {`items[]` {`code`, `estimated_quantity_milli`, `status`, `weekly_consumption_milli`, `days_left`}, `summary` {`sufficient`, `reorder`, `out`}} (§7) |
| `recipe_detail` | `recipe`, `servings_milli` | {`recipe`, `servings_milli`, `ingredients[]` {`product`, `quantity` (fracción), `cost_cents`, `kcal`}, `total_cost_cents`, `cost_per_serving_cents`, `total_kcal`, `kcal_per_serving`, `missing_prices[]`, `missing_nutrition[]`} (§8.2, §9.3) |
| `nutrition_coverage` | ninguno | {`unique`, `with_source`, `confirmed`, `proposed`, `missing[]`} (§9.1) |
| `health_profile` | `person` | {`person`, `age`, `bmi_tenths`, `bmi_class`, `suggested_kcal_exact` (fracción, antes del redondeo final), `suggested_kcal_target`, `daily_kcal_target`, `target_status` (`confirmed` · `suggested` · `none`), `missing_inputs[]`} (§10) |

**Traducción de `shopping_list` a C1 v1.2** (la hace la API; `ShoppingList` y `ShoppingListItem` de `docs/contratos/extractos/compras-tiendas-a.yaml`; los códigos pasan a ids). Las filas marcadas son PROVISIONAL (bloqueo `S4-asignacion-mercado`): con la opción A de sus preguntas 2 y 3, un artículo de vertical `market` sin tienda elegida queda en un mercado (o sin tienda) mientras S5 recomienda otra tienda por precio, y C1 v1.2 no tiene dónde publicar las dos.

| S4 (`shopping_list`) | C1 v1.2 | Nota |
|---|---|---|
| `week` | `start_date`, `end_date` | lunes y domingo de la semana (K8) |
| `items[].kind` | `origin` | |
| `items[].needed_milli`, `packages_to_buy`, `reasons[]` | mismos nombres | `needed_milli` ausente en `manual`; `packages_to_buy` sale de la cantidad a comprar (§6.6) |
| `items[].quantity_milli` | `quantity_milli` | la elección guardada, sin cambio (§6.3); la API la pasa a S5 en `shopping_list` (`S5-interfaz.md` §15.1) |
| `items[].recipes[]` | `recipe_ids` | |
| `items[].presentation` | ninguno | es la presentación que rige: `presentation_id` si la persona eligió; si no, `recommended_presentation_id`. `presentation_id` es la elección guardada y se conserva (§6.3) |
| `items[].recommended_presentation` | `recommended_presentation_id` | rige mientras `presentation_id` esté ausente, como dice C1 |
| `items[].recommended_store` | `recommended_store_id` | PROVISIONAL: en vertical `market` puede ser distinta de `items[].store`, y entonces **no rige la agrupación** |
| `items[].store` | `assigned_store_id` (campo calculado nuevo: solicitud a C1) | PROVISIONAL. No es `store_id`, que es la elección guardada y se conserva (§6.3) |
| `items[].estimated_cents` | `estimated_price_cents` | C1 dice «en la tienda asignada»; es el de la cotización que rige, que en vertical `market` puede ser de otra tienda (PROVISIONAL, solicitud de redacción) |
| `items[].price_observed_date`, `is_stale_price` | mismos nombres | |
| `items[].price_store`, `leftover_milli`, `missing_prices[]` | ninguno | un artículo sin precio sale sin `estimated_price_cents` |
| `groups[].vertical` | `groups[].purchase_channel` | |
| `groups[].store` | `groups[].store_id` | la tienda asignada; C1 dice «el elegido o el recomendado» (PROVISIONAL, solicitud de redacción) |
| `groups[].visit_date` | `groups[].visit_weekday` | día ISO de `visit_date` (1 = lunes … 7 = domingo), que es el `visit_weekday` del mercado; la fecha no tiene campo en C1 (solicitud opcional: `visit_date`) |
| `groups[].estimated_total_cents` | `groups[].total_cents` | |
| artículos con el mismo (`vertical`, `store`) | `groups[].items` | |
| `estimated_total_cents` | `ShoppingList.total_cents` | Σ `groups[].estimated_total_cents` |

**Solicitud a C1** (detalle en `docs/estado/wf1/correccion-contrato-S4S5-r3.md`):
- campo calculado `ShoppingListItem.assigned_store_id` = `items[].store`;
- `groups[].store_id` = «la tienda asignada (`assigned_store_id` de sus artículos)»;
- `recommended_store_id` = «recomendación de S5 por precio (`QuoteFor`); rige el precio y la presentación mientras `store_id` esté ausente; en vertical `market` no rige la agrupación»;
- `estimated_price_cents` = «costo estimado con la cotización que rige (S4 §8.1)».

**Variante sin campo nuevo** (si C1 no lo acepta): para un artículo sin tienda elegida, la API publica en `recommended_store_id` la tienda asignada (`items[].store`), y la cadena más barata solo se ve en `getBasketComparison`; la redacción de `estimated_price_cents` cambia igual. Con las opciones 2B y 3C del bloqueo, `items[].store` = `recommended_store` en todo artículo sin tienda elegida, y C1 v1.2 no cambia.

### 12.3 Oráculo y verificación independiente
- `oraculo/s4_run_fixtures.py` (Python 3.14 con `-I`, solo biblioteca estándar: `json`, `csv`, `fractions`, `decimal`, `datetime`) implementa las funciones de §12.2 desde este contrato, **sin leer el código Go**, y escribe `oraculo/salidas/s4.json` (FORMATO §4) con **todas** las consultas de **todos** los casos de `oraculo/fixtures/s4/`; `python3 -I oraculo/compara.py s4` lo enfrenta con el motor y con `esperado` (FORMATO §5). Lo escribe un agente distinto del que implementa `internal/recetas`. También genera la tabla de `s4-escalado-recetario` y los esperados de `s4-kcal-recetas-autor`.
- Los casos ILUSTRATIVOS (`fuente` `manual`) escriben en `descripcion` el cálculo a mano de sus cifras (por ejemplo, §8.3 y §10.3).
- S4-09 no lo cubre ninguna fixture del autor: lo hace el revisor con 5 recetas que él elige (párrafo S4-09 de §12).

## 13. Pruebas (nombres que casan con los `-run` de los criterios)

| Criterio | Pruebas Go | Fixtures |
|---|---|---|
| S4-01 | `TestRecetarioValido` (60 recetas, esquema, validación cruzada y cotejo de kcal en modo `review`), `TestRecetarioProduccion` (en modo `production` cargan las 60 sin cotejo aunque haya renglones `propuesto`; esas recetas salen incompletas, sin `ErrKcalMismatch`) | `recetario.json` |
| S4-02 | `TestEscalado2a4`, `TestEscaladoRecetario`, `TestEscaladoFracciones` | `s4-escalado-*` |
| S4-03 | `TestListaSinExistencia`, `TestListaExistenciaParcial`, `TestListaExistenciaMayor`, `TestListaSobrantes`, `TestListaConservaElecciones`, `TestListaPrecioYRecomendacion` (§6.5: precio y recomendación por artículo y total de la raíz), `TestListaCantidadFijadaYManual` (§6.6: cantidad fijada, artículos `manual` e `is_bought` en los totales) | `s4-lista-*`, `s4-costo-cotizacion-elegida` |
| S4-04 | `TestInventarioModos`, `TestInventarioEstado`, `TestInventarioMovimientos` | `s4-inventario-modos` |
| S4-05 | `TestPlanCeldas`, `TestPlanFuera`, `TestPlanKarenRechazada`, `TestPlanCasillasSinArmar`, `TestPlanSemanaDomingo` | `s4-plan-celdas` |
| S4-06 | `TestCostoReceta`, `TestCostoDiaSemana`, `TestCostoSinPrecio`, `TestCostoCotizacionQueRige` (§8.1), `TestCostoCantidadFijada` (§6.6: `estimated_cents` y totales de la lista con cantidad fijada) | `s4-costo-recetas`, `s4-costo-cotizacion-elegida`, `s4-lista-cantidad-manual` |
| S4-07 | `TestSaludIMC`, `TestSaludMeta`, `TestSaludAjustePersonaYCasa`, `TestSaludMetaConfirmada` | `s4-salud-perfiles` |
| S4-08 | `TestKcalReceta`, `TestKcalSinDato`, `TestKcalRecetasAutor`, `TestKcalCargadorCSV` (incluye `estado`, `confirmo` y `fecha_confirmacion`), `TestKcalCotejoRevision` (`kcal_por_porcion` o `fuentes_kcal[]` distintos → `ErrKcalMismatch` solo en modo `review`) | `s4-kcal-recetas-autor` |
| S4-09 | ninguna: lo verifica el revisor a mano con 5 recetas que él elige (§12) | opcional `s4-kcal-revisor`, del revisor |
| S4-11 | `TestNombresEnEspanol` (ningún nombre visible igual a `nombre_fuente`), `TestNutricionSinConfirmar` (en modo `production` un renglón `propuesto` no se carga) | `tabla.csv` |
| S4-12 | `TestCanalAgrupacion`, `TestCanalSinGrupo`, `TestCanalMercados` (§6.5) | catálogo, `s4-canal-mercados` |
| S4-13 | `TestRecetaEditable` (funciones puras de §3.3: alta, edición, borrado o archivo, sin base), `TestRestriccionesVacias` (§3.3) y la consulta SQL del campo vacío de §3.3 | — |
| S4-14 | toda la suite corre sin variables de IA definidas; `TestSinIA` verifica que el paquete no lee variables de entorno | — |
| S4-15 | `go vet`, `gofmt -l`, cobertura ≥ 85 % | — |

Ningún nombre de prueba de otro tema contiene `Plan`, `Lista`, `Costo`, `Salud` o `Kcal` (para que los `-run` de los criterios no mezclen).

## 14. Errores del paquete
`ErrPersonWithoutPlan`, `ErrEatingOutConflict`, `ErrUnknownProduct`, `ErrUnknownPresentation`, `ErrUnitMismatch`, `ErrProductWithoutChannel`, `ErrKcalMismatch` (solo en modo `review`, §3.2), `ErrInvalidServings` (≤ 0), `ErrNotMonday`, `ErrOverflow`. Estados que no son error: «sin precio», «sin dato nutricional», «sin sugerencia» (con la lista de insumos faltantes).

## 15. Verificación · criterio → sección

| Criterio | Sección |
|---|---|
| S4-01 Recetario completo y válido | §3 (cotejo de kcal solo en `review`, §3.2) |
| S4-02 Escalado exacto | §4, §13 |
| S4-03 Lista = Σ ingredientes − stock | §6 (precio y recomendación por artículo en §6.5; cantidad fijada por una persona y artículos `manual` en §6.6), §12.2 (traducción a C1) |
| S4-04 Inventario híbrido | §7, §6.2 |
| S4-05 Plan semanal por persona | §5 |
| S4-06 Costo por receta, día y semana | §8 (cotización que rige, §8.1; fecha, presentación y desempate de `QuoteFor` en `S5-interfaz.md` §14.1); costo de la lista con cantidad fijada en §6.6 |
| S4-07 Perfil, IMC y meta | §10 |
| S4-08 kcal con fuente citada | §9 (dos conteos en modo `review`) |
| S4-09 Cinco recetas a mano | §12 (las elige el revisor; `s4-kcal-recetas-autor` no lo sustituye), §9.3 |
| S4-10 Cargas locales y licencias | §9.1, §1 |
| S4-11 Español sin nombres en inglés | §3.2, §9.1 (confirmación y modo `production`) |
| S4-12 Canal por ingrediente | §11, §6.5 (mercados: PROVISIONAL, `S4-asignacion-mercado`), §12.2 (impacto en C1 y solicitud) |
| S4-13 Recetas editables | §3.3, §13 |
| S4-14 IA apagada | §1, §13 |
| S4-15 Calidad Go y supuestos | §1 (sin base en `internal/recetas`), §13 |
