# Corrección del contrato S4/S5 · ronda 3 (WF-1, 11-oct-2026)

Agente: `correccion-contrato-S4S5-r3`. Informe atendido: `docs/estado/wf1/REVISION-contrato-S4S5-r3.md` (NO APROBADO: 0 altos, 2 medios). Corregí los dos medios, M-r3-1 y M-r3-2. No toqué los bajos (B-r3-1 y B-r3-2 ni los abiertos de r1 y r2) ni reabrí decisiones de §A–§K. `S4-interfaz.md` y `S5-interfaz.md` pasan a la versión 1.3. `productos-formato.md` 1.1 no cambia: su §3 ya remite a S4 §6.5 y sigue siendo cierto.

## Hecho
- **M-r3-2 (S5 §14.1, regla única de `QuoteFor`).**
  - Fecha: `date` = «hoy» en America/Mexico_City (`today`, C-05), también para semanas futuras. La única excepción es §10, que evalúa en la fecha de despensa.
  - Candidatos: pares (tienda, presentación) con precio que compite, en las tiendas que compiten según §7.4. Con tienda elegida, con presentación elegida o con las dos, el mismo criterio dentro de esa elección.
  - Orden: (1) menor precio por unidad base ajustado, como fracción; (2) menor precio del paquete; (3) orden de `stores`; (4) orden del catálogo (en la API, menor id). Los pasos 2 a 4 son PROVISIONAL (bloqueo `S5-reglas-comparativa`, pregunta 2) y no dependen de la lista, porque `QuoteFor` no la recibe.
  - `recommended = true` solo sin elección.
- **M-r3-2 (S5 §7.3, §7.4, §10, §14, §15, §15.1, §15.2, §16, §18).**
  - `CompareBasket` recomienda cada renglón con `QuoteFor` sin elección y cotiza la tienda recomendada con esa misma presentación. En las demás tiendas sigue el menor costo de paquetes. `estimated_price_cents` = techo(necesidad / contenido) × precio del paquete. Resultado: una sola presentación recomendada.
  - El empate «más renglones, menos viejos, menor id» queda solo para el orden de tiendas por canasta.
  - `StalePricesForPlan` usa `QuoteFor` en la fecha de despensa, con la tienda que eligió la persona o la recomendada.
  - Consulta nueva `quote_for`, con las mismas claves que `quotes` de S4. Caso nuevo `s5-cotizacion-recomendada` y prueba `TestMejorTiendaCotizacion`.
  - `shopping_list.recommended_presentation` de S5 = `items[].presentation` de S4.
- **M-r3-2 (S4 §6.2 paso 5, §8.1, §12.1).** Las cotizaciones se piden con `date` = `today` y la regla de S5 §14.1. Con solo tienda elegida, la presentación la elige esa regla dentro de la tienda. Si una cotización por elección repite tienda y presentación de la recomendada, se manda solo la recomendada.
- **M-r3-1 (S4 §6.5, §12, §12.2, §13, §15).**
  - `shopping_list` entrega por artículo `recommended_store`, `recommended_presentation`, `price_store`, `price_observed_date`, `is_stale_price` y `estimated_cents`, y en la raíz `estimated_total_cents` = Σ de los grupos. Regla: sin elección en un eje, la recomendación de ese eje es la que rige. Con elección, es la de la cotización recomendada y solo informa.
  - §6.5 paso 3 dice que la opción A exige cambiar C1 v1.2, con el ejemplo de la uva.
  - §12.2 trae la tabla de traducción a C1: `visit_date` → `visit_weekday`, `estimated_total_cents` → `total_cents`, `ShoppingList.total_cents` = Σ, `items[].store` → `assigned_store_id` (nuevo). También trae la solicitud a C1 y la variante sin campo nuevo.
  - Fixtures `s4-canal-mercados` y `s4-costo-cotizacion-elegida` ampliadas; prueba `TestListaPrecioYRecomendacion`.
- **Bloqueos.**
  - `S4-asignacion-mercado`: en las opciones 2A y 3A anoté que exigen cambiar C1. Impacto y «Cómo seguí» citan la solicitud (lo pide M-r3-1).
  - `S5-reglas-comparativa`: precisión dentro de 2A (desempate por renglón sin la lista) para que el bloqueo no contradiga §14.1.

## Archivos
- `docs/contratos/motores/S5-interfaz.md` (1.2 → 1.3, 49.6 KB)
- `docs/contratos/motores/S4-interfaz.md` (1.2 → 1.3, 64.4 KB)
- `docs/estado/bloqueos/S4-asignacion-mercado.md` (agregados marcados «corrección r3»)
- `docs/estado/bloqueos/S5-reglas-comparativa.md` (agregados marcados «corrección r3»)
- `docs/estado/wf1/correccion-contrato-S4S5-r3.md` (este reporte)

## Verificación
- `python3 -I datos/tmp/corr-s4s5-r3/verif.py` (solo biblioteca estándar, `Decimal` y `Fraction`; lee `data.json` y los `.md`): **84 OK, 0 fallas**.
  - **Canon S5-01**, calculado desde los renglones:
    - idénticos 277350 / 287400 (+10050);
    - totales 309450 / 323300 (+13850), iguales a `totales`;
    - piezas 55 / 61 / 60, renglones 33 / 3, ganadores 22 / 8 / 3;
    - el texto de S5 cita cada cifra.
  - **Unidad base S5-02:** 2933 / 2800, 12692 / 16406 y 3375 / 2833.
  - **Otras cifras sin cambio:** S4 §8.3 da 4425 / 2213 / 8850; MSJ 1748750 y 1330250; RN-15 1486571; K40 [0.9836, 0.4918, 0].
  - **Caso nuevo `s5-cotizacion-recomendada`** (implementación independiente de §14.1 y §7.3):
    - arroz sin elección → Walmart 1 kg a 2800;
    - con `store` = Bodega → 750 g a 2200;
    - con `presentation` = 750 g → empate en (1) y (2), gana Walmart por (3);
    - frijol → Bodega 500 g a 1500, por (2);
    - por kg: 2800, 2933, 2933, 3000, 3000 y 3000;
    - `compare_basket` (arroz, 500 g): Walmart 1 kg a 2800 y Bodega 500 g a 1500; `cheapest_store` = Bodega.
    - Con la regla vieja de §7.3, Walmart habría elegido 750 g a 2200. Era el hallazgo, y ya no ocurre.
  - **Estructura:**
    - bloques JSON válidos y cercas pares;
    - versión 1.3 con la cita de la revisión r3 en S4 y S5;
    - S4 §15 cubre S4-01 a S4-15 y S5 §18 cubre S5-01 a S5-13;
    - cada `-run` de `docs/criterios/S4.md` y `S5.md` encuentra pruebas, solo de su criterio, sin mezclas (incluidas `TestListaPrecioYRecomendacion` y `TestMejorTiendaCotizacion`);
    - todo bloqueo citado existe (salvo `S4-meta-ajustes`, hipotético según S4-07 1.2);
    - no queda ninguna frase vieja: «la que ya tiene asignados más renglones de la lista», «la cotización que rige en la tienda asignada (o la recomendada)», «con tienda o presentación: la de esa elección» y «la recomendación se sigue viendo en la comparativa»;
    - los dos bloqueos conservan los 5 apartados del protocolo §2.
- **Criterios tocados y evidencia:**
  - S4-03: §6.2 paso 5, §6.5 (precio y recomendación) y §12.2 (traducción).
  - S4-06: §8.1, fecha y regla de S5 §14.1.
  - S4-12: §6.5 paso 3 y §12.2 (impacto en C1, PROVISIONAL).
  - S5-03: §7.3, §7.4, §14.1, `quote_for` y `s5-cotizacion-recomendada`.
  - S5-12: §10 con §14.1.
- **Código de otros agentes** (informativo; no es de este artefacto y no lo toqué):
  - `go vet ./internal/precios/... ./internal/recetas/...`: `internal/recetas/types.go:229: undefined: Nutrition`, así que `recetas` no compila.
  - `gofmt -l` marca `internal/precios/seed.go`, `types.go`, `internal/recetas/num.go` y `types.go`.
  - `go test -cover`: `precios/ocr` 99.3 %; `precios`, `asistida` y `profeco` sin pruebas (0.0 %).
  - `oraculo/fixtures/s4/` solo tiene `_datos/` vacío; `oraculo/fixtures/s5/` y `oraculo/salidas/s4.json` y `s5.json` no existen. Por eso `compara.py s4` y `compara.py s5` aún no se pueden correr.

## Supuestos
- S-r3-1: la fecha de `QuoteFor` es «hoy» también para el costo de una semana futura. Sale de C-05, de S5-04 (la edad se mide contra «hoy»), de `ErrFutureDate` y del contraste con S5-12, que sí pide la fecha de despensa. No cambia ninguna cifra canon.
- S-r3-2: `estimated_price_cents` es el costo de los paquetes para la necesidad, no el precio de un paquete. Así cuadra con el total de grupo de S4 (Σ paquetes × precio). En el ejemplo del revisor coincide, porque es un paquete.
- S-r3-3: cuando la persona eligió en un eje, la recomendación de ese eje solo informa y sale de la cotización recomendada. Así se cumple literalmente «rige mientras X esté ausente» de C1.

## Bloqueos
- Ninguno nuevo.
- Actualizados: `docs/estado/bloqueos/S4-asignacion-mercado.md` (impacto en C1 de las opciones 2A y 3A) y `docs/estado/bloqueos/S5-reglas-comparativa.md` (desempate por renglón dentro de 2A).
- Sigue abierto, sin tocar: `docs/estado/bloqueos/S1-parametros-avisos.md`.

## Pendientes
- Ninguno de los hallazgos alta o media.
- Bajos que siguen abiertos y no se pidieron:
  - B-r3-1: `recipe_book_file` sin tabla nutricional en `review` (S4 §3.2, §12.1).
  - B-r3-2: cómo encuentran la raíz del repositorio las pruebas Go (S5 §15.1, §16; S4 §12.1).
  - Los de r1 y r2.

## Solicitudes fuera de alcance
1. **C1** (fuente de `docs/contratos/openapi.yaml` que genera `backend/internal/api/c1gen/build.py`; extracto `docs/contratos/extractos/compras-tiendas-a.yaml`). Motivo: M-r3-1. Con la opción A de `S4-asignacion-mercado`, la tienda asignada y la recomendada pueden ser distintas, y C1 v1.2 no tiene dónde publicar las dos.
   - `ShoppingListItem.assigned_store_id`: campo nuevo, entero, `x-calculado`, «Tienda o mercado donde se agrupa el artículo (S4 §6.5); ausente = sin tienda».
   - `recommended_store_id`: «Recomendación de S5 por precio (`QuoteFor`); rige precio y presentación mientras `store_id` esté ausente; en vertical `market` no rige la agrupación».
   - `estimated_price_cents`: «Costo estimado: `packages_to_buy` × precio del paquete de la cotización que rige (S4 §8.1); en vertical `market` puede ser de otra tienda».
   - `ShoppingList.groups[].store_id`: «Tienda o mercado asignado (`assigned_store_id` de sus artículos)».
   - Opcional: `groups[].visit_date`.
   - Variante sin campo nuevo: `recommended_store_id` = la tienda asignada en artículos sin tienda elegida, y solo cambia la redacción de `estimated_price_cents`.
   - Si la puerta 2 elige 2B y 3C, la solicitud se cancela.
2. **Motor S5** (`backend/internal/precios/basket.go`, agente implementador). Alinear con S5 1.3:
   - `QuoteFor` con solo presentación: hoy `case store == ""` ignora `pres`.
   - Desempate por renglón de §14.1, en vez de `breakTie(assigned, stale)`.
   - `CompareBasket`: recomendación = `QuoteFor` sin elección; la tienda recomendada se cotiza con esa presentación; `estimated_price_cents` = techo × precio; el desempate de canasta cuenta renglones recomendados.
   - `StalePricesForPlan` con `QuoteFor` en la fecha de despensa.
   - Consulta `quote_for` en `s5 dump` y `TestMejorTiendaCotizacion`.
3. **Motor S4 y su oráculo** (`backend/internal/recetas`, `oraculo/rec_*.py`): campos nuevos de `shopping_list` (§6.5), `estimated_total_cents` en la raíz y `TestListaPrecioYRecomendacion`.
4. **Autor de fixtures y oráculo de S5** (`oraculo/fixtures/s5/`, `oraculo/s5_run_fixtures.py`): caso `s5-cotizacion-recomendada` y consulta `quote_for`, con las cifras de S5 §14.1.
5. **Integrador:** edité dos archivos de bloqueo que no están en la lista literal de mi alcance. `S4-asignacion-mercado` lo pide M-r3-1; `S5-reglas-comparativa` evita que el bloqueo contradiga §14.1. Los dos agregados van marcados «corrección r3».

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Ninguna de K31.
- Dato para la puerta 2: la opción A de las preguntas 2 y 3 de `S4-asignacion-mercado` obliga a cambiar C1 (campo `assigned_store_id`); con 2B y 3C, C1 v1.2 queda igual.

## Consumo aproximado
- Unas 33 acciones de herramienta.
- Unos 200 KB leídos (protocolo, revisión, decisiones §D, §E y §K, criterios S4 y S5, partes de S4, S5, C1 y bloqueos, y 6 KB de código).
- Unos 45 KB escritos: parches ~25 KB, verificación ~7 KB y este reporte ~11 KB. Los scripts reescribieron completos S4 (64 KB), S5 (50 KB) y los dos bloqueos.
- Temporales en `datos/tmp/corr-s4s5-r3/`: borrados al terminar.
