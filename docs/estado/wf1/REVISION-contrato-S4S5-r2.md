# Revisión del contrato S4/S5 · ronda 2 de 3

Revisor: `revisor-contrato-S4S5-r2` (independiente; no escribió nada de lo revisado ni corrige nada) · 11-oct-2026 · WF-1.
Artefactos: `docs/contratos/motores/S4-interfaz.md` 1.1, `S5-interfaz.md` 1.1, `productos-formato.md` 1.0; reportes `docs/estado/wf1/contrato-S4S5.md` y `correccion-contrato-S4S5-r1.md`.
Base de evaluación: `docs/criterios/S4.md` 1.2, `docs/criterios/S5.md` 1.2, `docs/plan/decisiones.md` §D, §E y §K. Apoyos que leí por partes: `oraculo/fixtures/FORMATO.md` §1, §2 y §4; RN-12, RN-16 y RN-17; extracto C1 `compras-tiendas-a.yaml` (`BasketComparison`); migraciones 00001 y 00003 (claves de `app_setting`).

## Veredicto
**NO APROBADO.** Hay 0 hallazgos altos y 4 medios nuevos. Los 7 medios de la ronda 1 quedan **cerrados**. Las cifras canon (+138.50 y +100.50, K-01 a K-08 y ganadores 22 / 8 / 3) y el catálogo de los 61 artículos los recalculé yo y coinciden. Los 4 medios nuevos son preguntas que el implementador o el autor del oráculo tendrían que adivinar. Tres salen de la corrección r1 (rutas de `_datos/`, el modo del cotejo de kcal y la elección de la cotización) y uno es un hueco de RN-16 que el contrato no cubre.

## Verificación propia (comandos y salida resumida)
Scripts en `datos/tmp/rev-s4s5-r2/`, ya borrados. `python3 -I`, `Decimal` y `Fraction`, leyendo `datos/semilla/data.json` directamente.
- `canon.py`:
  - 0 renglones con `unitario × cantidad ≠ total`.
  - Bodega **309450** = idénticos 277350 + sustitutos 30000 + atún 2100. Walmart **323300** = 287400 + 35900.
  - W − B = **13850** (+138.50) = 10050 + 5900 − 2100. Idénticos 287400 − 277350 = **10050** (+100.50).
  - Piezas 61 / 60 / 55; renglones 33 / 3; ganadores Bodega 22, empates 8, Walmart 3 (= `resultados_por_producto`).
  - Unidad base: arroz 2933 / 2800, crema para café 12692 / 16406, leche 3375 / 2833, con el ganador de la nota.
  - S4 §8.3: los índices 21, 16, 15 y 0 son espagueti 12.00, atún 21.00, elote 13.50 y aceite 150.00. Resultado 4425 / 2213 / 8850.
  - Frescura: 60/61 (0.9836), 59/61, 30/61 (0.4918) y 0. Fechas 2026-11-08 (29 días), 2026-11-09 (30) y 2027-01-08 (90).
  - Salud: MSJ 1748.75, 1330.25 y 1648.75. Metas 2099, 1596, 1799, 1846, 2556 y 2056 (2555.5625 exacto).
  - IMC 185, 250, 300, 237 y 185 (73990 g da 18.4975 exacto). RN-15: 1486571. Factores 6/5, 11/8, 31/20, 69/40 y 19/10 = RN-17.
- `catalogo.py` (lee las tablas §5.1 y §5.2 de `productos-formato.md`):
  - 36 productos y 39 presentaciones. Los 39 nombres de `data.json` (33 idénticos + 3 + 3 sustitutos) existen tal cual. Todo id lleva el prefijo de su producto.
  - §5.2 tiene 39 renglones, **61 piezas** y **309450**. Son idénticos renglón a renglón a `cuentas.axel`, `cuentas.karen` y `cuentas.conjunto` (0 diferencias).
  - Cada renglón casa con el precio unitario de Bodega. Las piezas por presentación cuadran con `items_identicos[].cantidad` (atún 6 = 5 + 1). Unidades solo `g`, `ml` y `pieza`.
- Claves `app_setting` que cita S5 §1: las cuatro existen. `price.stale_after_days` y `price.assisted_reading_enabled` vienen de 00001; `price.zero_weight_after_days` y `price.profeco_import_max_gb`, de 00003. `health.goal_adjust_*` vienen de 00003. RN-12 escribe la fórmula general `(D0 − edad)/(D0 − Dv + 1)` (K44) y RN-17, los 5 factores.
- Bloqueos nuevos (`S5-carrito-comparativa`, `S4-confirmacion-empates`) y vigentes (`S4-reglas-lista`, `S5-reglas-comparativa`, `S5-profeco-almacen`): todos existen con los 5 apartados del protocolo §2.
- Entorno: `tesseract` instalado. `goose` no está instalado, y S4 §3.3 dice que se instala. `datos/semilla/tickets/`, `recetario.json` y `productos.json` todavía no existen.

## Tabla de criterios
| Criterio | Sev. | Estado | Sección | Evidencia mía |
|---|---|---|---|---|
| S4-01 Recetario | alta | cubierto | §3 | llave `recetas`, esquema en `datos/semilla/recetario.schema.json`, validación cruzada |
| S4-02 Escalado exacto | alta | cubierto | §4 | fracción reducida; tabla independiente con `fractions` |
| S4-03 Lista = Σ − stock | alta | cubierto | §6 | 4 fixtures; reposición PROVISIONAL (`S4-reglas-lista`); RN-15 incluye la reposición |
| S4-04 Inventario | alta | cubierto (B4 r1 abierto) | §7 | 1486571 verificado |
| S4-05 Plan | alta | cubierto | §5 | Karen rechazada; domingo abre la semana siguiente (K38); huecos (K43) |
| S4-06 Costo | alta | cubierto, con observación (M-r2-2) | §8 | 4425 / 2213 / 8850 verificado; falta la regla para elegir la cotización entre varias |
| S4-07 Salud | alta | cubierto | §10 | MSJ, factores y metas verificados; K23 y K36 (por persona, casa como inicial, «mantener» = 0, un redondeo); clasificación PROVISIONAL |
| S4-08 kcal con fuente | alta | cubierto, con observación (M-r2-3) | §9 | dos conteos en modo `review`; el cotejo `kcal_por_porcion` choca con el modo `production` |
| S4-09 5 recetas a mano | alta | cubierto | §12 | las elige el revisor; la fixture del autor no lo sustituye. Verificación manual: no aplica a esta revisión (no hay recetario) |
| S4-10 Cargas locales | media | cubierto (B1 r1 abierto) | §9.1, §1 | `PROCEDENCIA.md` con 3 fuentes y licencias |
| S4-11 Español | alta | cubierto | §3.2, §9.1 | `estado` / `confirmo` / `fecha_confirmacion`; modo `production`; bloqueo `S4-confirmacion-empates` |
| S4-12 Canal | media | **incumplido (M-r2-4)** | §11, §6.5 | no dice qué mercado toca cuando hay varios, ni qué pasa si no hay ninguno |
| S4-13 Editables | media | cubierto (B-r2-4) | §3.3 | funciones puras, prueba sin base, consulta SQL con goose |
| S4-14 IA apagada | media | cubierto (B2 r1 abierto) | §1, §13 | |
| S4-15 Calidad | alta | cubierto | §1, §13 | nada de SQLite en `internal/recetas` |
| S5-01 Comparativa | alta | cubierto | §7.2, §3, §15 | +138.50 y +100.50 recalculados; `report_file` sin copias; `piece_count` e `identical_pieces` |
| S5-02 Unidad base | alta | cubierto | §4 | 2933/2800, 12692/16406, 3375/2833 |
| S5-03 Mejor tienda | alta | cubierto | §7.3, §7.4 | desempate PROVISIONAL (`S5-reglas-comparativa`) y solicitud para escribirlo en RN |
| S5-04 Frescura y fuentes | alta | cubierto | §5, §3 | K40 exacto; `manual`, `ticket_manual` y `ticket_ocr` compiten; `site_assisted` y PROFECO solo referencia (K22, K37) |
| S5-05 PROFECO | alta | cubierto, con observación (M-r2-1; B6 r1 abierto) | §11 | conteo leído = guardados + rechazados; ligas confirmadas |
| S5-06 OCR | alta | cubierto (sin medir) | §12 | `tickets/` solo con fotos; `tickets-esperado/` hermana |
| S5-07 Lectura asistida | alta | cubierto | §13 | 3 líneas por sitio; clave `price.assisted_reading_enabled` en 00001 |
| S5-08 Sin red | alta | cubierto | §1 | |
| S5-09 Membresía | alta | cubierto (B3 r1 abierto) | §8 | 8 compras; `missing_fee`; redondeo PROVISIONAL |
| S5-10 Trazabilidad | media | cubierto | §3, §6 | 5 fuentes; vigencia de promociones PROVISIONAL |
| S5-11 Mercados | media | cubierto | §9, §7.4 | 2026-10-13 |
| S5-12 Precios viejos | media | cubierto | §10 | |
| S5-13 Calidad | alta | cubierto | §1, §16 | |

## Coherencia con K7, K22, K23 y K24
- **K7:** coherente. Precio por 1 kg, 1 l o 1 pieza con HALF_UP (S5 §4). IMC de la OMS con cortes 18.5 / 25 / 30 (S4 §10.1). Equilibrio de Sam's en compras al año (S5 §8). Estado `out` / `reorder` / `sufficient` (S4 §7).
- **K22:** coherente. Viejo a los 30 días o más; PROFECO nunca compite (S5 §3 y §6; S4 §8.1). Con K37, `manual` compite y `site_assisted` es solo referencia. Con K40/K44 aplica la recta exacta con parámetros de `app_setting` (RN-12).
- **K23:** coherente. Mifflin-St Jeor × factor de 5 niveles (RN-17) ± ajuste por objetivo. La persona confirma o escribe su meta (S4 §10.4). Con K36: por persona, NULL toma el valor de la casa, iniciales 0 y 0, «mantener» no ajusta y hay un solo redondeo.
- **K24:** coherente. 5 fotos en `datos/semilla/tickets/` (sin fijar formato) y transcripción en `tickets-esperado/`. El renglón correcto exige descripción, cantidad e importe exactos. Mientras falten las fotos, el resultado es «sin medir». Todo OCR pasa por confirmación humana.

## Hallazgos de esta ronda

### Medios
- **M-r2-1 · Suficiencia del oráculo y del motor (S5-05, S5-06): las rutas de `_datos/` no tienen con qué resolverse.**
  - Ruta: `docs/contratos/motores/S5-interfaz.md` §15.1 (línea 216) y §15.2 (`profeco_import`, `ocr_parse_tsv`).
  - Problema: los argumentos `files[]`, `links_file` y `tsv_file` van «relativos a la carpeta del caso (`_datos/`), como el `$ref`». Pero el evaluador de FORMATO §4.1 es `func(input json.RawMessage, q Query)`: no recibe la carpeta de fixtures. Además `main.go` es del agente de S2 y «nadie más lo edita». El `$ref` lo resuelve `main.go`; estos argumentos los tendría que resolver el evaluador de S5, y no sabe dónde está la carpeta. El implementador tendría que adivinar: leer `os.Args` o fijar `oraculo/fixtures/s5/` en el código.
  - Corrección: declarar esas rutas relativas a la raíz del repositorio (`oraculo/fixtures/s5/_datos/<archivo>`), con la misma regla que los `*_file` de §15.1. Otra salida: pedir al integrador que FORMATO §4.1 pase la carpeta al evaluador.
- **M-r2-2 · Suficiencia (S4-03, S4-06): falta cómo se elige la cotización cuando hay varias.**
  - Ruta: `S4-interfaz.md` §12.1 (`quotes`), §8.1 y §6.2 paso 5.
  - Problema: S4 usa «la que recomienda S5» o «la mejor tienda por ingrediente de S5» cuando la persona no eligió tienda ni presentación. `quotes` es una lista plana `{product, store, presentation, …}` sin marca de cuál es la recomendada. Una fixture que traiga la cotización de la tienda elegida y la de la mejor tienda para el mismo producto deja sin regla a `recipe_detail`, `week_plan` y `shopping_list`. El oráculo y el motor Go podrían elegir distinto: la primera de la lista, la más barata o repetir la lógica con peso de S5 §7.1, que S4 no debe duplicar.
  - Corrección: agregar `recommended` (booleano, exactamente una por producto) a cada cotización de `quotes` y decir que las demás solo sirven para atender `list_choices`. Otra salida: fijar «una cotización sin tienda elegida por producto».
- **M-r2-3 · El cotejo de `kcal_por_porcion` choca con el modo `production` (S4-01, S4-08; lo introdujo la corrección M5).**
  - Ruta: `S4-interfaz.md` §3.2 (fila `kcal_por_porcion` y `fuentes_kcal[]`) contra §9.1 (`LoadNutritionCSV` en `production`).
  - Problema: el generador escribe `kcal_por_porcion` desde `tabla.csv`, y el cargador del recetario lo recalcula; si no coincide, `ErrKcalMismatch` y la receta «no se carga». En `production` los renglones `propuesto` no se cargan, así que hasta la puerta 2 toda receta con un ingrediente sin confirmar recalcula `null` (incompleta) contra un número escrito. Resultado: no carga, y la semilla de producción tendría menos de 60 recetas. El contrato no dice con qué modo se escribe ni con cuál se coteja.
  - Corrección: decir que el cotejo corre solo en modo `review`, o que en `production` una receta con kcal incompleta se carga con `kcal` ausente y `missing_nutrition`, sin `ErrKcalMismatch`. Además, fijar que el generador escribe `kcal_por_porcion` en modo `review` (o solo con renglones `confirmado`) y que `fuentes_kcal[]` se coteja igual.
- **M-r2-4 · S4-12 / RN-16: falta qué mercado toca cuando hay varios o ninguno.**
  - Ruta: `S4-interfaz.md` §6.5 (y `groups[]` de `shopping_list` en §12.2).
  - Problema: los artículos de vertical `market` «se asignan a un mercado que surte su grupo en su próximo día de visita». No dice cuál si varios mercados surten el grupo (decisiones §D: «varios»), ni qué pasa si ninguno lo surte. Ese segundo caso es el del arranque: RN-16 dice «Axel aún no da de alta sus mercados». Tampoco dice qué pasa si S5 recomienda una cadena más barata para un producto de mercado: `productos-formato.md` §3 y S5 §7.4 dejan competir a las cadenas por la uva. `groups[].store`, `visit_date` y la agrupación que exige S4-12 quedan a criterio del implementador y del oráculo.
  - Corrección: fijar la regla. Por ejemplo: el mercado con la `next_visit` más cercana y, en empate, el de menor código. Sin mercado, el artículo va al grupo `market` con `store` nulo y la leyenda «sin mercado dado de alta», o pasa a la cadena que recomienda S5. Si la regla no se deduce de RN-16, que quede PROVISIONAL con un bloqueo (puede sumarse a `S4-reglas-lista`).

### Bajos
- **B-r2-1** `S5-interfaz.md` §7.2 y §15.2: en modo `compare_carts` no se definen `missing_count`, `stale_count` ni `difference_vs_cheapest_cents`. Por ejemplo, ¿la sexta lata de atún cuenta como faltante en Walmart? En `compare_basket` tampoco se definen `piece_count` ni `identical_pieces`. *Corrección:* una línea por clave y modo.
- **B-r2-2** `S5-interfaz.md` §12 (`ocr_parse_tsv`): faltan detalles del TSV de Tesseract.
  - La confianza viene con decimales y vale −1 en renglones que no son palabra.
  - No dice con qué separador se une `descripcion`.
  - No dice qué hacer con un número inicial que no es cantidad (por ejemplo, un código de barras).
  - Tampoco fija el formato de `cantidad` en `tickets-esperado/NN.csv`.
  - *Corrección:* fijarlos para que el oráculo y Go den lo mismo en `s5-ocr-parseo`.
- **B-r2-3** El extracto C1 `BasketComparison` (descripción) afirma que la API da el canon 309450 contra 323300. Con la opción A PROVISIONAL de `S5-carrito-comparativa`, la API no lo reproduce, y el contrato no lo pide como solicitud. *Corrección:* agregar la solicitud condicionada: con A, que la descripción de C1 diga que el canon es del motor.
- **B-r2-4** `S4-interfaz.md` §3.3: sin la semilla de S1, la consulta `COUNT(*) … WHERE COALESCE(dietary_restrictions,'') <> ''` = 0 sale cierta en una tabla vacía y no prueba nada. *Corrección:* insertar antes dos perfiles con el campo nulo, o exigir la semilla.
- **B-r2-5** `S5-interfaz.md` §15.1 (`report_file`): la sexta lata (`atun_sexta_pieza_bodega`) no trae nombre en `data.json`. Conviene citar la presentación `atun-agua.dolores-140g`. Además, §3 habla de «los 39 renglones» de la compra y `report_file` arma 37 (33 + 1 + 3). Son equivalentes una vez agrupados por presentación; conviene decirlo.
- **B-r2-6** `S5-interfaz.md` §8: «últimos 365 días» no dice si incluye los extremos (`today − 365 < fecha ≤ today`).
- Siguen abiertos de r1 (la corrección no los tocó, por instrucción): **B1, B2, B3, B4, B6, B7 y B8**. **B5** queda resuelto, porque S5 1.1 usa `report_file` en lugar del `$ref` fuera de `_datos/`.

## Estado de los hallazgos de la ronda 1
| # | Hallazgo r1 | Estado | Evidencia |
|---|---|---|---|
| 1 | S5-01 / C1: canon modelado como dos compras; siembra de Walmart | **cerrado** | S5 §3 y §7.2: Bodega es la compra registrada (`expense`) y Walmart un carrito (`cart`) que **no** se siembra como gasto. `CompareCarts` reproduce K-01 a K-08. La API con una sola `expense_id` queda PROVISIONAL (bloqueo `S5-carrito-comparativa`, con los 5 apartados). La fuente de cada semilla está declarada (K37). Canon recalculado por mí |
| 2 | Input de fixtures y oráculos sin declarar | **cerrado** (ver M-r2-1 y M-r2-2) | S4 §12.1–§12.3 y S5 §15.1–§15.3: claves, tipos, obligatoriedad, `unknown_field` / `missing_field`, argumentos y salida por función. `s4_run_fixtures.py` y `s5_run_fixtures.py` cubren todos los casos, los escribe otro agente y no leen Go |
| 3 | S4-09: las 5 recetas las elegía el autor | **cerrado** | S4 §12 (párrafo S4-09), §13 y §15: las elige el revisor; `s4-kcal-recetas-autor` no lo sustituye |
| 4 | S5-06: fotos y transcripciones en la misma carpeta | **cerrado** | S5 §12: `tickets/` solo con 5 fotos (`.jpg`, `.jpeg` o `.png`) y `tickets-esperado/NN.csv` hermana, así que `ls \| wc -l` = 5 |
| 5 | Confirmación humana sin campo ni responsable | **cerrado** (ver M-r2-3) | `tabla.csv` y `vinculos.csv` con `estado`, `confirmo` y `fecha_confirmacion`; modos `production` y `review`; S4-08 con dos conteos; quién confirma PROVISIONAL (bloqueo `S4-confirmacion-empates`) |
| 6 | `identical_count` renglones contra «artículos» | **cerrado** | Verifiqué en el extracto C1 que `BasketComparison` resume el canon como «33 idénticos, 3 sustitutos» y que `items` es «Por artículo» (renglón). `identical_count` = 33 casa con eso. K-07 sale en `piece_count` (61 / 60) e `identical_pieces` (55), y hay una solicitud a C1 para la descripción |
| 7 | S4-13 / S4-15: `store.go` y pruebas con base | **cerrado** | S4 §1: ningún archivo de `internal/recetas` importa SQLite. §3.3: funciones puras, `TestRecetaEditable` sin base, `TestRestriccionesVacias` y la consulta SQL con goose, que corre el agente del motor |
