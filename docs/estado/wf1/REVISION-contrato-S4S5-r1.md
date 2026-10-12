# Revisión del contrato S4/S5 · ronda 1 de 3

Revisor: `revisor-contrato-S4S5-r1` (WF-1, 11-oct-2026). Independiente: no escribí nada de lo que reviso y no corregí nada. Evalúo solo contra `docs/criterios/S4.md` 1.2, `docs/criterios/S5.md` 1.2 y decisiones §D, §E y §K.

Artefactos revisados: `docs/contratos/motores/S4-interfaz.md` 1.0, `docs/contratos/motores/S5-interfaz.md` 1.0, `docs/contratos/motores/productos-formato.md` 1.0 y el reporte `docs/estado/wf1/contrato-S4S5.md`. Para cotejar, leí por partes RN-12 a RN-20, extractos C1 (`compras-tiendas-a.yaml`, `compras-tiendas-b.yaml`, `salud.yaml`, `comidas.yaml`, `ajustes-c.yaml`, solo con grep), el esquema por grep (00001 a 00003), `oraculo/fixtures/FORMATO.md` y los 3 bloqueos nuevos.

## Veredicto

**NO APROBADO.** 0 hallazgos altos, 7 medios y 8 bajos. Todas las cifras canon y todos los ejemplos numéricos de los tres documentos los reproduje exactos con mis propios scripts. Ninguna cifra contradice el oráculo ni §K. Lo que impide aprobar es un tema de suficiencia y coherencia: el implementador siguiente tendría que adivinar en siete puntos, y en dos de ellos (M1 y M6) el contrato choca con C1.

## Verificación propia (scripts en `datos/tmp/rev-s4s5/`, `python3 -I`, `Decimal` y `Fraction`; ya borrados)

| Qué | Resultado mío | Contrato |
|---|---|---|
| K-01 Bodega = idénticos + sustitutos + atún | 277350 + 30000 + 2100 = **309450** = `totales.bodega` | igual |
| K-02 Walmart | 287400 + 35900 = **323300** = `totales.walmart` | igual |
| K-03 **+138.50** | 323300 − 309450 = **13850** = 10050 + 5900 − 2100 | igual |
| K-06 **+100.50** | 287400 − 277350 = **10050** | igual |
| K-07 piezas | 61 / 60 / 55 (= `articulos_*`) | igual |
| K-08 renglones | 33 / 3 | igual |
| unitario × cantidad = total | 0 renglones distintos | igual |
| Ganadores por precio unitario | 22 / 8 / 3 = `resultados_por_producto` | igual |
| Unidad base de los sustitutos | 2933 vs 2800 (W), 12692 vs 16406 (B), 3375 vs 2833 (W); mismo ganador que las notas | igual |
| S4 §8.3 costo ilustrativo | `items_identicos[21]`, `[16]`, `[15]`, `[0]` = espagueti 12.00, atún 21.00, elote 13.50, aceite de oliva 150.00 → 1200 + 2100 + 675 + 450 = **4425**; por porción 2213; para 4 porciones 8850 | igual |
| Frescura (K40) | 29 d (2026-11-08) peso 1; 30 d 60/61 = 0.9836; 31 d 59/61 = 0.9672; 60 d 30/61 = 0.4918; 90 d (2027-01-08) 0; 91 d no compite; `[0.9836, 0.4918, 0.0]` | igual |
| Mifflin-St Jeor (S4-07) | 1748.75 → × 1.2 = 2098.5 → 2099; 1330.25 → 1596.3 → 1596; −300 → 1798.5 → 1799; +250 → 1846.3 → 1846; RN-17 1648.75 × 31/20 = 2555.5625 → 2556; −500 → 2056 | igual |
| Fórmula MSJ en milésimas (`10·g + 625·mm − 5000·edad ± …`) | equivale a la de RN-17 | igual |
| IMC en décimas `HALF_UP(g·10000/mm²)` | 74000/2000 → 185; 100000 → 250; 120000 → 300; 72500/1750 → 237; 73990/2000 = 18.4975 exacto → 185 | igual |
| RN-15 | 2838000 − HALF_UP(946000·30/21 = 1351428.57) = **1486571** | igual |
| RN-16 | 10-oct-2026 es sábado (ISO 6); el siguiente martes es 2026-10-13 | igual |
| S5-09 ilustrativo | 100000 ÷ 12500 = 8; RN-13 13850 × 6 = 83100 | igual |
| Catálogo (e) | 36 productos, 39 presentaciones; los 39 nombres de `data.json` (33 + 3 + 3) existen tal cual y no hay nombres de más; los 39 renglones de §5.2 coinciden en orden, texto, piezas e importe con `cuentas` (axel, karen, conjunto) = **61 piezas** y **3,094.50**; precio unitario de cada renglón = `bodega_unit` (o total ÷ cantidad en sustitutos); piezas por presentación = `items_identicos[].cantidad` (atún 5 + 1 = 6); las 3 presentaciones solo de Walmart existen | igual |
| Esquema | columnas y enumerados citados existen (`store.membership_annual_fee_cents`, `assisted_reading_enabled`, `visit_weekday`, `price.source` con 5 valores, `nutrition_fact.source`, `tracking_mode`, `activity_level`, `goal`, `goal_adjust_*`, `ocr_draft_json`); claves `price.*` y `health.*` sembradas en 00001 y 00003 | igual |
| C1 | campos y enumerados citados existen (`bmi_class` 4 valores, `status` sufficient/reorder/out, `basket` shopping_list/expense, `verdict` 3 valores, `break_even_purchases` entero, `needed_milli` = plan − existencia + reposición) salvo lo de M1 y M6 | ver M1, M6 |
| Entorno | `tesseract 5.5.3`, idiomas `osd spa`; `unrar`, `bsdtar` y `sqlite3` presentes; `datos/semilla/tickets/` no existe (S5-06 «sin medir») | igual |
| Marcas | S4: 5 PROVISIONAL, todas con bloqueo; S5: 10, todas con bloqueo salvo «PROVISIONAL de S5-09» (etiqueta del propio criterio) y el tiempo de Tesseract (propuesta K31 con solicitud); los 3 bloqueos nuevos siguen el formato del protocolo §2 | correcto |

## Criterios

Estado: **cubierto** = el contrato fija todo lo que el criterio exige; **con observación** = cubierto, con un hallazgo bajo; **parcial** = le falta algo (hallazgo medio).

| Criterio | Sev. | Estado | Sección | Evidencia mía |
|---|---|---|---|---|
| S4-01 Recetario | alta | cubierto | S4 §3 | esquema `datos/semilla/recetario.schema.json`, llave `recetas`, campos con presentación comprable; el bloque JSON de §3.1 es válido |
| S4-02 Escalado exacto | alta | cubierto | §4, §13 | fracción reducida, 4 escalas × 60 recetas contra tabla `fractions` de otro agente; 250000 → 375000 |
| S4-03 Lista = Σ − existencia | alta | cubierto | §6 | 3 fixtures más sobrantes; la reposición está en RN-15 FIRME y su fórmula es PROVISIONAL con bloqueo `S4-reglas-lista` |
| S4-04 Inventario híbrido | alta | con observación (B4) | §7 | 1486571 verificado; agotado prevalece |
| S4-05 Plan por persona | alta | cubierto | §5 | Karen rechazada, «fuera» sin kcal, colación opcional, K38, K43 |
| S4-06 Costo | alta | cubierto | §8 | §8.3 reproducido desde `data.json`; «sin precio», nunca 0 |
| S4-07 Salud | alta | cubierto | §10 | todas las cifras de la fixture verificadas; factores en RN-17; la clasificación en ±0.05 es PROVISIONAL con bloqueo |
| S4-08 kcal con fuente | alta | parcial (M5) | §9 | fórmulas verificadas; falta cómo se registra la confirmación del empate |
| S4-09 Cinco recetas a mano | alta | parcial (M3) | §12, §15 | la fixture la elige el autor |
| S4-10 Cargas locales y licencias | media | cubierto | §9.1, §1 | PROCEDENCIA.md con las 3 fuentes y licencias; sin red |
| S4-11 Español | alta | parcial (M5) | §3.2, §9.1 | `nombre_fuente` solo en el CSV; falta el registro de confirmación (R9) |
| S4-12 Canal | media | cubierto | §11, §6.5 | canal por grupo; `ErrProductWithoutChannel` |
| S4-13 Recetas editables | media | parcial (M7) | §3.2, §13 | no queda claro dónde vive la prueba de alta, edición y borrado, ni la consulta del campo vacío |
| S4-14 IA apagada | media | con observación (B2) | §1, §13 | sin red ni variables de entorno; no cita las claves `ai.*` |
| S4-15 Calidad Go | alta | con observación (M7) | §1, §13 | sin flotantes, 85 %; riesgo de cobertura por `store.go` |
| S5-01 Comparativa | alta | cubierto (con M1 y M6) | §7.2, §3, §15 | K-01…K-08 y 22/8/3 verificados por mí |
| S5-02 Unidad base | alta | cubierto | §4 | 2933/2800, 12692/16406, 3375/2833 |
| S5-03 Mejor tienda | alta | cubierto | §7.3, §7.4 | desempate PROVISIONAL con bloqueo (RN no lo tiene; solicitud 6) |
| S5-04 Frescura y fuentes | alta | cubierto | §5, §3, §7.1 | tabla de fronteras verificada; 3 fuentes compiten, 2 solo referencia |
| S5-05 PROFECO | alta | con observación (B6) | §11 | conteo leídos = guardados + rechazados; almacén PROVISIONAL con bloqueo |
| S5-06 OCR ≥ 80 % | alta | parcial (M4) | §12 | tesseract 5.5.3 spa verificado; la carpeta de tickets choca con la verificación del criterio |
| S5-07 Lectura asistida | alta | cubierto | §13 | doble interruptor + veredicto por sitio con URL y fecha |
| S5-08 Sin red | alta | cubierto | §1, §13 | sin cliente HTTP en WF-1 |
| S5-09 Membresía | alta | con observación (B3) | §8 | 8 compras; `missing_fee` |
| S5-10 Trazabilidad | media | cubierto | §3, §6 | 5 fuentes; vigencia de promoción PROVISIONAL con bloqueo |
| S5-11 Mercados | media | cubierto | §9, §7.4 | 2026-10-13 verificado |
| S5-12 Precios viejos | media | cubierto | §10 | evaluado en la fecha de despensa |
| S5-13 Calidad Go | alta | cubierto | §1, §16 | sin flotantes; 85 % |

**(b) Suficiencia:** motor S4 y S5: casi completa (M1, M6, M7). Recetario: suficiente. Nutrición y vínculos PROFECO: falta la confirmación humana (M5). OCR: falta la carpeta de transcripciones (M4). Oráculos y fixtures: insuficiente (M2, M3).
**(c) Canon:** +138.50 y +100.50 reproducidos exactos (tabla de arriba).
**(d) §K:** K7 (cortes del IMC, unidad base HALF_UP por kg, l o pieza, equilibrio en compras al año, estado del inventario), K22 (30 días; PROFECO nunca compite; ticket y OCR compiten), K23 (MSJ × 5 factores de 1.2 a 1.9 ± ajuste, confirmación) y K24 (5 tickets, renglón correcto, «sin medir») están bien reflejados; ver M4 sobre la verificación de K24.
**(e) Catálogo:** cubre las 61 piezas y los 39 renglones, verificado.

## Hallazgos

### Medios

**M1 · El modo de comparación del canon no cabe en C1 y supone una compra en Walmart.** `S5-interfaz.md` §3 («dos compras registradas, una por tienda… fuente `ticket_manual`»), §7.2 («Compara dos o más compras registradas») y §14 (`CompareExpenses(compras[], date)`). En C1, `getBasketComparison` con `basket = expense` recibe **una sola** `expense_id` (`compras-tiendas-a.yaml`, parámetros `basket` y `expense_id`). Además, `data.json` dice «Precios del carrito ya con descuentos e IVA» y solo la compra de Bodega se registró y se repartió (`cuentas` = 61 piezas = 3,094.50). Si se siembra Walmart como compra con ticket, el sistema registra un gasto de 3,233.00 que no ocurrió, con una fuente (`ticket_manual`) que no consta. Nada de esto aparece como supuesto, bloqueo ni solicitud. El canon no cambia, porque `manual` y `ticket_manual` compiten igual (K37). Lo que cambia es la operación y la semilla (protocolo §2: «un supuesto que cambie… endpoints es un bloqueo»).
*Corrección:* modelar el canon como una compra registrada (Bodega, `expense`) contra el carrito de Walmart, con precios `manual` del 2026-10-10 y las cantidades del carrito, que solo existe como entrada del motor y de la fixture y no como gasto. Otra opción es abrir la pregunta en `S5-reglas-comparativa` y pedir a C1, como solicitud, un segundo parámetro (`compare_expense_id` o un carrito). En cualquier caso, declarar la fuente de cada semilla según K37.

**M2 · Las fixtures no tienen definido su `Input` ni sus argumentos, y S5 no nombra su oráculo.** `S4-interfaz.md` §12 y `S5-interfaz.md` §15. `FORMATO.md` §0 exige que «cada motor declara en su interfaz… sus tipos de entrada», y §4 que cada motor nombre su oráculo. S2 lo hace con tablas `Input`/`Params` (claves JSON, tipo, si es obligatorio, origen, `unknown_field`). S4 y S5 solo dan tipos Go (§2) y una lista en prosa («catálogo, recetas, plan, inventario…»). No fijan las claves raíz del JSON, los nombres en snake_case de cada campo, ni los `argumentos` de `scale_recipe`, `shopping_list`, `compare_expenses`, `freshness`, etc. S5 tampoco nombra un verificador independiente; S4 nombra `oraculo/s4_run_fixtures.py`, pero solo para 2 casos. El autor de las fixtures y el de `cmd/motores` tendrían que adivinar las mismas claves, y `compara.py` fallaría por diferencias de nombre.
*Corrección:* agregar en cada interfaz una tabla «`Input` (raíz de `entradas.datos`)» con claves, tipos, unidades (`_cents`, `_milli`), obligatoriedad y error de campo desconocido, y una tabla de `argumentos` y forma de salida por función. Nombrar `oraculo/s5_run_fixtures.py`, o decir qué casos verifica otra implementación independiente y cuáles salen de cálculo a mano.

**M3 · S4-09: las cinco recetas las elige el autor.** `S4-interfaz.md` §12, caso `s4-kcal-cinco-recetas` («5 recetas del recetario… el revisor lo rehace a mano»), y §15 (S4-09 → esa fixture). El criterio S4-09 (alto) dice: «El revisor elige 5 recetas (no las elige el autor)». Tal como está, el revisor solo rehace las que escogió quien implementa.
*Corrección:* decir que las 5 recetas de S4-09 las elige el revisor al revisar, entre las 60, y que `s4-kcal-cinco-recetas` es una prueba propia del autor que no sustituye a S4-09. Otra opción: que el revisor agregue sus 5 a la fixture.

**M4 · S5-06: la carpeta de tickets choca con la verificación del criterio.** `S5-interfaz.md` §12 pone las fotos `NN.jpg` y las transcripciones `NN.esperado.csv` en la misma carpeta `datos/semilla/tickets/`. S5-06 verifica `ls datos/semilla/tickets/ | wc -l` = 5 imágenes; con las transcripciones al lado da 10. Además, el contrato fija `.jpg` sin que K24 lo diga.
*Corrección:* llevar las transcripciones a `datos/semilla/tickets/esperado/NN.csv` (o a `oraculo/fixtures/s5/ocr/`) y aceptar cualquier imagen que Tesseract lea (jpg o png); o pedir que el criterio cuente solo imágenes.

**M5 · Confirmación humana de empates (S4-11, R9) y vínculos PROFECO sin registro.** `S4-interfaz.md` §9.1 («Todo empate producto ↔ fuente es asistido con confirmación humana; ninguno entra automático») y `S5-interfaz.md` §11 (`vinculos.csv` «con confirmación humana»). Ni `tabla.csv` ni `vinculos.csv` tienen dónde constar la confirmación (quién y cuándo), y no se dice quién confirma en WF-1. El agente de nutrición no sabe si debe incluir los empates sin confirmar (y arriesgar S4-11) o dejarlos fuera (y no llegar al 100 % de S4-08).
*Corrección:* agregar columnas `estado` (propuesto o confirmado), `confirmo` y `fecha_confirmacion` a ambos CSV, y que el cargador rechace en producción lo no confirmado. Si no está decidido quién confirma ni cuándo (p. ej. Axel en la puerta 2), abrir un bloqueo y marcar PROVISIONAL lo que dependa de eso.

**M6 · `identical_count` y las piezas por tienda.** `S5-interfaz.md` §7.2, pasos 5 y 6: `identical_count` = renglones (33). C1 lo describe como «Artículos idénticos.», y en `data.json` y en S5-01 «artículos» son piezas (`articulos_identicos` = 55; «61 contra 60 artículos»). Además, K-07 (61 y 60 piezas por tienda) no tiene clave de salida: solo se nombra `identical_pieces` (55), que C1 no tiene. El implementador de la API y el de las fixtures tendrían que adivinar qué número va en `identical_count` y con qué clave sale K-07.
*Corrección:* nombrar las claves de salida del motor y del volcado para K-07 (p. ej. `stores[].piece_count` y `identical_pieces`) y pedir a C1, como solicitud, que la descripción de `identical_count` diga «renglones idénticos» (o alinear el contrato a piezas).

**M7 · Dónde vive la persistencia (S4-13 y S4-15).** `S4-interfaz.md` §1 permite «un archivo `store.go` aparte que no usan las pruebas del núcleo». Si ese archivo vive en `internal/recetas`, cuenta en la cobertura del 85 % (S4-15, alto) sin pruebas. §13 pide `TestRecetaEditable` (alta, edición, archivo) y `TestRestriccionesVacias` en un paquete que no toca la base. S4-13 pide insertar, editar y borrar una receta y una consulta que confirme el campo vacío. No queda claro si esas pruebas usan SQLite ni dónde va la consulta.
*Corrección:* fijar la persistencia fuera del paquete medido (subpaquete `internal/recetas/store` con sus propias pruebas, o la capa de S6). Decir que `TestRecetaEditable` prueba las funciones puras de validación, cambio y archivo de recetas de la casa, y quién hace la consulta de `health_profile.dietary_restrictions` vacío (S1 en la semilla o S6).

### Bajos

- **B1** `S4-interfaz.md` §9.1: `fuente` de `tabla.csv` solo admite `USDA`, `OFF` y `manual`. El esquema admite además `mx_table` y `store_label`, y decisiones §D distingue tablas mexicanas y etiquetas. *Corrección:* mapear SMAE a `mx_table` y etiqueta a `store_label`, o justificar por qué se juntan.
- **B2** `S4-interfaz.md` §1 y §13 (S4-14): no cita los interruptores que el criterio revisa, `ai.local_model_enabled` y `ai.claude_api_enabled` (00001, `false`). *Corrección:* citarlos y que `TestSinIA` (o la prueba de S1) verifique su valor inicial.
- **B3** `S5-interfaz.md` §8: sigue bien a K7 y S5-09 (equilibrio en compras al año), pero RN-13 dice «el punto de equilibrio es el gasto anual en Sam's» y el contrato no pide alinearlo. *Corrección:* solicitud para corregir RN-13 a «compras al año» (K7).
- **B4** `S4-interfaz.md` §7: «agotado prevalece» se decide por razonamiento. S4-04 dice que la prioridad «la declare `reglas-negocio.md`», y ni RN-15 ni C1 la declaran. El resultado coincide con lo que el criterio espera. *Corrección:* solicitud para escribir la prioridad en RN-15.
- **B5** `S5-interfaz.md` §15: `{"$ref": "../../../datos/semilla/data.json"}` contradice `FORMATO.md` §2 (solo `_datos/`; ya está pedido en la solicitud 2). Además, el texto dice «relativo a la raíz del repositorio» y la ruta es relativa a la carpeta del caso. *Corrección:* corregir el texto y esperar la solicitud.
- **B6** `S5-interfaz.md` §11: S5-05 consulta `MAX(fecha_registro)`, y ese campo se guarda como texto en un formato que el diccionario no da («Datetime 8»), así que el máximo puede salir mal. El contrato ya guarda `fecha` en ISO. *Corrección:* decir que la verificación usa `MAX(fecha)`, o guardar `fecha_registro` normalizada.
- **B7** `productos-formato.md` §5.1: la presentación `atun-agua.dolores-140g` se llama «…140 g (5 pzas)» y su contenido es 140000 (una lata). El nombre exacto se exige para empatar, pero en pantalla parece un paquete de 5. *Corrección:* nombre visible limpio y el texto de `data.json` como alias para empatar.
- **B8** Reporte `contrato-S4S5.md`, Verificación: dice «bloques JSON válidos (3 de 3)» y encontré 2 (S4 §3.1 y productos §1; S5 no trae ninguno). Ambos son válidos. *Corrección:* ajustar la cifra.

## Estado de los anteriores
Ronda 1: no hay revisiones previas de este artefacto.
