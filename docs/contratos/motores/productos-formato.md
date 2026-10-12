# Formato del catálogo base de productos (`datos/semilla/productos.json`)

Versión 1.1 · 11-oct-2026 · WF-1, agente `contrato-S4S5`; 1.1 por `correccion-contrato-S4S5-r2` (hallazgo M-r2-4 de `docs/estado/wf1/REVISION-contrato-S4S5-r2.md`: solo §3, último punto). Define el archivo que escribirá otro agente (catálogo semilla) y que leen el cargador de la base, el recetario (S4), el comparador (S5), la tabla de nutrición y los vínculos de PROFECO. **Sin código.** Fuentes: decisiones §D y §E, K7; RN-15, RN-16, RN-20; esquema 1.3 (`product_group`, `product`, `presentation`, `store_product_group`); `datos/semilla/data.json`.

## 1. Estructura

```json
{
  "version": "1.0",
  "fecha": "2026-10-11",
  "grupos": [
    {"id": "abarrotes", "nombre": "Abarrotes", "canal": "cadena", "orden": 30}
  ],
  "productos": [
    {
      "id": "arroz-grano-largo",
      "nombre": "Arroz blanco grano largo",
      "categoria": "abarrotes",
      "unidad_base": "g",
      "gramos_por_unidad": null,
      "fuente_gramos": null,
      "vertical": "cadena",
      "sinonimos": ["arroz", "arroz blanco", "arroz grano largo"],
      "presentacion_tipica": "arroz-grano-largo.sos-750g",
      "presentaciones": [
        {"id": "arroz-grano-largo.sos-750g", "nombre": "Arroz SOS grano largo súper extra 750 g",
         "marca": "SOS", "contenido_milli": 750000, "granel": false, "codigo_barras": null,
         "nombres_ticket": ["Arroz SOS 750 g"]}
      ]
    }
  ]
}
```

## 2. Campos

| Campo | Regla | Columna |
|---|---|---|
| `grupos[].id` | `^[a-z0-9]+(-[a-z0-9]+)*$`, único | — (el cargador empata por nombre) |
| `grupos[].nombre` | es-MX, único sin distinguir mayúsculas | `product_group.name` |
| `grupos[].canal` | `mercado` o `cadena` | `product_group.purchase_channel` = `market` o `chain` |
| `grupos[].orden` | entero, orden en la lista de compras | `product_group.sort_order` |
| `productos[].id` | `^[a-z0-9]+(-[a-z0-9]+)*$`, único y **estable**: lo citan `recetario.json`, `nutricion/tabla.csv`, `profeco/vinculos.csv` y las fixtures | — (la base empata por `name`, que es único) |
| `.nombre` | genérico, en español de México, único sin distinguir mayúsculas; nunca el nombre en inglés de USDA ni una marca | `product.name` |
| `.categoria` | id de un `grupos[]` (la categoría **es** el grupo de compra) | `product.product_group_id` |
| `.unidad_base` | `g`, `ml` o `pieza` | `product.base_unit` = `g`, `ml`, `piece` |
| `.gramos_por_unidad` | decimal en texto con hasta 3 decimales (`"10.5"`), **obligatorio en `ml` y `pieza`** si el producto se usa en una receta (para kcal, `S4-interfaz.md` §9.2) y prohibido en `g` (CHECK del esquema); `null` = sin dato | `product.grams_per_unit_milli` (×1000) |
| `.fuente_gramos` | cita obligatoria cuando hay `gramos_por_unidad` (etiqueta, `fdc_id` de USDA, código de barras de OFF) | `product.notes` |
| `.vertical` | `mercado`, `cadena` o `ambas` (§3) | — (se deriva del grupo) |
| `.sinonimos` | arreglo de cadenas es-MX en minúsculas; para búsqueda, empate de renglones de ticket y OCR; puede ir vacío | — (los usa la búsqueda; no hay columna) |
| `.presentacion_tipica` | id de una de sus `presentaciones`: la que sugiere la receta y la lista cuando S5 no tiene precio | `recipe_ingredient.presentation_id` por omisión |
| `.presentaciones[]` | al menos una | `presentation` |
| `..id` | `<id del producto>.<slug>`, único en el archivo | — |
| `..nombre` | como se ve en tienda; **en los artículos de `data.json`, el texto exacto de `items_identicos[].producto` o `sustitutos[].bodega/walmart`** (S5 empareja por ese texto) | `presentation.name` |
| `..marca` | texto o `null` | `presentation.brand` |
| `..contenido_milli` | entero > 0 en milésimas de la unidad base (el **factor** de equivalencia): 750 g = 750000; 1 l = 1000000; 24 cubos = 24000 | `presentation.size_milli` |
| `..granel` | booleano; `true` = se vende por kilo o por unidad base | `presentation.is_bulk` |
| `..codigo_barras` | 8 a 14 dígitos o `null` (cruce con Open Food Facts) | `presentation.barcode` |
| `..nombres_ticket` | textos con que aparece en tickets o en `data.json: cuentas` (61 piezas en 39 renglones) | — (alias para empatar renglones) |

**Validación del cargador:** ids únicos y con el patrón; toda `categoria` existe; `presentacion_tipica` pertenece al producto; `contenido_milli > 0`; `gramos_por_unidad` según la unidad; los 39 nombres de presentación de `data.json` (33 idénticos + 3 sustitutos de Bodega + 3 de Walmart) existen tal cual; los 39 textos de `cuentas` aparecen en algún `nombres_ticket`. El catálogo es editable desde Ajustes (`createProduct`, `updateProduct`, `createPresentation`); el archivo solo siembra.

## 3. Vertical y canal
- `mercado`: el grupo debe tener `canal = mercado` (verduras, fruta y carne no especializada, decisiones §D).
- `cadena`: el grupo debe tener `canal = cadena`.
- `ambas`: el grupo da el canal de la lista (cualquiera de los dos) y el producto **también** se puede comprar en el otro. La base no guarda «ambas»: basta con que Axel declare el grupo en los mercados que lo surten (`store_product_group`), y S5 deja competir a esos mercados (`S5-interfaz.md` §7.4). El cargador lo usa solo para sugerir grupos al dar de alta un mercado.
- Una cadena con precio compite por cualquier producto, también de vertical `mercado` (la uva verde del 10-oct se compró en Bodega Aurrera). Competir en S5 no mueve el artículo de grupo en la lista de compras: su tienda y su grupo siguen `S4-interfaz.md` §6.5 (vertical `mercado`: PROVISIONAL, bloqueo `S4-asignacion-mercado`).

## 4. Grupos iniciales (propuesta del formato, editables en Ajustes)

| id | Nombre | Canal |
|---|---|---|
| `frutas-verduras` | Frutas y verduras | mercado |
| `carnes` | Carne y pollo (no especializada) | mercado |
| `abarrotes` | Abarrotes | cadena |
| `lacteos` | Lácteos y huevo | cadena |
| `salchichoneria` | Salchichonería | cadena |
| `panaderia` | Panadería | cadena |
| `limpieza` | Limpieza del hogar | cadena |
| `higiene` | Higiene personal | cadena |
| `papel` | Papel y desechables | cadena |

Los productos de las 60 recetas que no están en `data.json` (verduras, carne, huevo, tortilla, especias…) los agrega el agente del catálogo con este mismo formato; las verduras y la carne no especializada van en grupos de canal `mercado`.

## 5. Los artículos de `data.json` en el catálogo (verificado)

Tabla generada y verificada con `python3 -I` contra `data.json`: los 39 renglones de `cuentas` suman **61 piezas** y **3,094.50**; cada renglón empata con una presentación cuyo precio unitario coincide con `bodega_unit` (o `total ÷ cantidad` en sustitutos); las piezas por presentación cuadran con `items_identicos[].cantidad` (atún: 6 = 5 idénticas + la sexta lata). **Notas de contenido** a confirmar en etiqueta por el agente del catálogo: «Pasta dental Colgate Triple Acción 2 x 80 ml» = 160 ml; «Yoghurt bebible Alpura multi sabor 10 pzas 220 g» = 10 × 220 g; «Queso petit suisse Danonino fresa 6 pzas 42 g» = 6 × 42 g; «Uva verde por pieza de 1 kg aprox» = 1 kg nominal; piezas: rollos (papel higiénico), hojas dobles (toallas), servilletas y cubos (caldo). La leche de Bodega no dice «entera»; por eso el producto genérico es «Leche».

### 5.1 Productos y presentaciones (36 productos, 39 presentaciones)

| Producto (`id`) | Nombre es-MX | Categoría | Unidad base | Vertical | Presentaciones (`id`: contenido en milésimas · nombre exacto de `data.json`) |
|---|---|---|---|---|---|
| `aceite-oliva` | Aceite de oliva | Abarrotes | ml | cadena | `aceite-oliva.nutrioli-oli-500ml`: 500000 · Aceite de oliva Nutrioli Oli 500 ml |
| `aceite-soya` | Aceite de soya | Abarrotes | ml | cadena | `aceite-soya.nutrioli-946ml`: 946000 · Aceite de soya Nutrioli 946 ml |
| `frijoles-refritos` | Frijoles refritos bayos | Abarrotes | g | cadena | `frijoles-refritos.la-costena-220g`: 220000 · Frijoles bayos La Costeña refritos 220 g |
| `duraznos-almibar` | Duraznos en almíbar | Abarrotes | g | cadena | `duraznos-almibar.la-costena-425g`: 425000 · Duraznos La Costeña cubitos en almíbar 425 g |
| `cafe-soluble` | Café soluble | Abarrotes | g | cadena | `cafe-soluble.nescafe-verarica-180g`: 180000 · Café soluble Nescafé Reserva Mexicana Verarica 180 g |
| `jamon-pavo` | Jamón de pavo | Salchichonería | g | cadena | `jamon-pavo.sabori-400g`: 400000 · Jamón de pavo Sabori Parma virginia premium 400 g |
| `salchicha-pavo` | Salchicha de pavo | Salchichonería | g | cadena | `salchicha-pavo.sabori-500g`: 500000 · Salchicha de pavo Sabori Parma 500 g |
| `avena` | Avena en hojuelas | Abarrotes | g | cadena | `avena.3-minutos-400g`: 400000 · Avena 3 Minutos con sémola de trigo 400 g |
| `blanqueador` | Blanqueador con cloro | Limpieza del hogar | ml | cadena | `blanqueador.cloralex-3750ml`: 3750000 · Blanqueador Cloralex El Rendidor 3.75 l |
| `caldo-pollo-cubos` | Caldo de pollo en cubos | Abarrotes | pieza | cadena | `caldo-pollo-cubos.knorr-24-cubos`: 24000 · Caldo de pollo Knorr 24 cubos |
| `pasta-dental` | Pasta dental | Higiene personal | ml | cadena | `pasta-dental.colgate-2x80ml`: 160000 · Pasta dental Colgate Triple Acción 2 x 80 ml |
| `jabon-manos` | Jabón líquido para manos | Higiene personal | ml | cadena | `jabon-manos.grisi-450ml`: 450000 · Jabón líquido Grisi frutos rojos 450 ml |
| `acondicionador` | Acondicionador para cabello | Higiene personal | ml | cadena | `acondicionador.pantene-170ml`: 170000 · Acondicionador Pantene 3 Minute Miracle Repara 170 ml |
| `desinfectante-aerosol` | Desinfectante en aerosol | Limpieza del hogar | g | cadena | `desinfectante-aerosol.lysol-354g`: 354000 · Desinfectante Lysol Crisp Linen aerosol 354 g |
| `detergente-polvo` | Detergente en polvo | Limpieza del hogar | g | cadena | `detergente-polvo.roma-1kg`: 1000000 · Detergente en polvo Roma multiusos 1 kg |
| `elote-grano` | Granos de elote | Abarrotes | g | cadena | `elote-grano.herdez-220g`: 220000 · Granos de elote Herdez 220 g |
| `atun-agua` | Atún en agua | Abarrotes | g | cadena | `atun-agua.dolores-140g`: 140000 · Atún Dolores aleta amarilla en agua con aceite 140 g (5 pzas) |
| `uva-verde` | Uva verde | Frutas y verduras | g | mercado | `uva-verde.pieza-1kg`: 1000000 · Uva verde por pieza de 1 kg aprox |
| `jabon-trastes` | Jabón líquido para trastes | Limpieza del hogar | ml | cadena | `jabon-trastes.salvo-1200ml`: 1200000 · Jabón para trastes Salvo Power Clean 1.2 l |
| `pan-caja` | Pan de caja | Panadería | g | cadena | `pan-caja.oroweat-mantequilla-680g`: 680000 · Pan Oroweat mantequilla 680 g |
| `papel-higienico` | Papel higiénico | Papel y desechables | pieza | cadena | `papel-higienico.cottonelle-30-rollos`: 30000 · Papel higiénico Kleenex Cottonelle Elegance 30 rollos |
| `espagueti` | Espagueti | Abarrotes | g | cadena | `espagueti.la-moderna-250g`: 250000 · Spaghetti La Moderna 250 g |
| `queso-manchego` | Queso manchego | Lácteos y huevo | g | cadena | `queso-manchego.nochebuena-200g`: 200000 · Queso manchego NocheBuena rebanadas 200 g |
| `mayonesa-chipotle` | Mayonesa con chipotle | Abarrotes | g | cadena | `mayonesa-chipotle.mccormick-310g`: 310000 · Mayonesa McCormick con chilpotle squeeze 310 g |
| `salsa-soya` | Salsa de soya | Abarrotes | ml | cadena | `salsa-soya.maggi-290ml`: 290000 · Soya Maggi tradicional 290 ml |
| `catsup` | Catsup | Abarrotes | g | cadena | `catsup.heinz-397g`: 397000 · Catsup Heinz 397 g |
| `mayonesa` | Mayonesa | Abarrotes | g | cadena | `mayonesa.mccormick-420g`: 420000 · Aderezo de mayonesa McCormick 420 g |
| `salsa-inglesa` | Salsa inglesa | Abarrotes | ml | cadena | `salsa-inglesa.crosse-blackwell-290ml`: 290000 · Salsa Inglesa Crosse & Blackwell 290 ml |
| `servilletas` | Servilletas de papel | Papel y desechables | pieza | cadena | `servilletas.petalo-420`: 420000 · Servilletas Pétalo 420 pzas |
| `toallas-papel` | Toallas de papel | Papel y desechables | pieza | cadena | `toallas-papel.petalo-mega-jumbo-270`: 270000 · Servitoalla Pétalo Mega Jumbo 270 hojas dobles |
| `yogur-griego` | Yogur griego natural | Lácteos y huevo | g | cadena | `yogur-griego.yoplait-750g`: 750000 · Yoghurt Yoplait Griego natural sin azúcar 750 g |
| `yogur-bebible` | Yogur bebible | Lácteos y huevo | g | cadena | `yogur-bebible.alpura-10x220g`: 2200000 · Yoghurt bebible Alpura multi sabor 10 pzas 220 g |
| `petit-suisse` | Queso petit suisse | Lácteos y huevo | g | cadena | `petit-suisse.danonino-6x42g`: 252000 · Queso petit suisse Danonino fresa 6 pzas 42 g |
| `arroz-grano-largo` | Arroz blanco grano largo | Abarrotes | g | cadena | `arroz-grano-largo.sos-750g`: 750000 · Arroz SOS grano largo súper extra 750 g<br>`arroz-grano-largo.sos-uruguayo-1kg`: 1000000 · Arroz SOS grano largo uruguayo 1 kg |
| `sustituto-crema-cafe` | Sustituto de crema para café | Abarrotes | g | cadena | `sustituto-crema-cafe.coffee-mate-pack-780g`: 780000 · Pack Coffee Mate 640 g + 140 g (780 g)<br>`sustituto-crema-cafe.coffee-mate-640g`: 640000 · Coffee Mate en polvo 640 g |
| `leche` | Leche | Lácteos y huevo | ml | cadena | `leche.santa-clara-4x1l-taza`: 4000000 · Leche Santa Clara 4 pack 1 l + taza de regalo<br>`leche.santa-clara-entera-6x1l`: 6000000 · Leche Santa Clara entera caja 6 pzas 1 l |

### 5.2 Las 61 piezas de `data.json` (39 renglones de `cuentas`) con su nombre en el catálogo

Orden y textos de `data.json: cuentas.axel.lineas`, `cuentas.karen.lineas` y `cuentas.conjunto.lineas`. El texto del renglón entra como `nombres_ticket` de la presentación.

| # | Cuenta | Renglón en `data.json` | Piezas | Importe | Producto del catálogo | Presentación | Unidad |
|---|---|---|---|---|---|---|---|
| 1 | axel | Duraznos La Costeña cubitos 425 g | 1 | 40.00 | `duraznos-almibar` Duraznos en almíbar | `duraznos-almibar.la-costena-425g` | g |
| 2 | axel | Jamón de pavo Sabori 400 g | 1 | 62.00 | `jamon-pavo` Jamón de pavo | `jamon-pavo.sabori-400g` | g |
| 3 | axel | Salchicha de pavo Sabori 500 g | 1 | 40.00 | `salchicha-pavo` Salchicha de pavo | `salchicha-pavo.sabori-500g` | g |
| 4 | axel | Avena 3 Minutos 400 g | 1 | 36.00 | `avena` Avena en hojuelas | `avena.3-minutos-400g` | g |
| 5 | axel | Lysol Crisp Linen 354 g | 1 | 139.00 | `desinfectante-aerosol` Desinfectante en aerosol | `desinfectante-aerosol.lysol-354g` | g |
| 6 | axel | Frijoles refritos La Costeña | 2 | 24.00 | `frijoles-refritos` Frijoles refritos bayos | `frijoles-refritos.la-costena-220g` | g |
| 7 | axel | Granos de elote Herdez | 2 | 27.00 | `elote-grano` Granos de elote | `elote-grano.herdez-220g` | g |
| 8 | axel | Atún Dolores 140 g | 6 | 126.00 | `atun-agua` Atún en agua | `atun-agua.dolores-140g` | g |
| 9 | axel | Uva verde 1 kg | 1 | 87.00 | `uva-verde` Uva verde | `uva-verde.pieza-1kg` | g |
| 10 | axel | Leche Santa Clara 4 pack + taza | 1 | 135.00 | `leche` Leche | `leche.santa-clara-4x1l-taza` | ml |
| 11 | axel | Pan Oroweat mantequilla 680 g | 1 | 92.00 | `pan-caja` Pan de caja | `pan-caja.oroweat-mantequilla-680g` | g |
| 12 | axel | Queso manchego NocheBuena 200 g | 1 | 61.00 | `queso-manchego` Queso manchego | `queso-manchego.nochebuena-200g` | g |
| 13 | axel | Yoghurt Yoplait Griego 750 g | 1 | 75.00 | `yogur-griego` Yogur griego natural | `yogur-griego.yoplait-750g` | g |
| 14 | axel | Yoghurt Alpura bebible 10 pzas | 1 | 108.00 | `yogur-bebible` Yogur bebible | `yogur-bebible.alpura-10x220g` | g |
| 15 | karen | Frijoles refritos La Costeña | 2 | 24.00 | `frijoles-refritos` Frijoles refritos bayos | `frijoles-refritos.la-costena-220g` | g |
| 16 | karen | Salchicha de pavo Sabori 500 g | 1 | 40.00 | `salchicha-pavo` Salchicha de pavo | `salchicha-pavo.sabori-500g` | g |
| 17 | karen | Papel higiénico Cottonelle 30 rollos | 1 | 199.00 | `papel-higienico` Papel higiénico | `papel-higienico.cottonelle-30-rollos` | pieza |
| 18 | karen | Danonino fresa 6 pzas | 1 | 37.00 | `petit-suisse` Queso petit suisse | `petit-suisse.danonino-6x42g` | g |
| 19 | karen | Granos de elote Herdez | 3 | 40.50 | `elote-grano` Granos de elote | `elote-grano.herdez-220g` | g |
| 20 | conjunto | Aceite de oliva Nutrioli 500 ml | 1 | 150.00 | `aceite-oliva` Aceite de oliva | `aceite-oliva.nutrioli-oli-500ml` | ml |
| 21 | conjunto | Aceite de soya Nutrioli 946 ml | 3 | 144.00 | `aceite-soya` Aceite de soya | `aceite-soya.nutrioli-946ml` | ml |
| 22 | conjunto | Arroz SOS 750 g | 3 | 66.00 | `arroz-grano-largo` Arroz blanco grano largo | `arroz-grano-largo.sos-750g` | g |
| 23 | conjunto | Pack Coffee Mate | 1 | 99.00 | `sustituto-crema-cafe` Sustituto de crema para café | `sustituto-crema-cafe.coffee-mate-pack-780g` | g |
| 24 | conjunto | Nescafé Verarica 180 g | 1 | 112.00 | `cafe-soluble` Café soluble | `cafe-soluble.nescafe-verarica-180g` | g |
| 25 | conjunto | Cloralex 3.75 l | 2 | 116.00 | `blanqueador` Blanqueador con cloro | `blanqueador.cloralex-3750ml` | ml |
| 26 | conjunto | Caldo Knorr 24 cubos | 2 | 112.00 | `caldo-pollo-cubos` Caldo de pollo en cubos | `caldo-pollo-cubos.knorr-24-cubos` | pieza |
| 27 | conjunto | Spaghetti La Moderna | 2 | 24.00 | `espagueti` Espagueti | `espagueti.la-moderna-250g` | g |
| 28 | conjunto | Jabón Grisi 450 ml | 2 | 124.00 | `jabon-manos` Jabón líquido para manos | `jabon-manos.grisi-450ml` | ml |
| 29 | conjunto | Pasta dental Colgate | 1 | 35.00 | `pasta-dental` Pasta dental | `pasta-dental.colgate-2x80ml` | ml |
| 30 | conjunto | Pantene 3 Minute Miracle | 2 | 180.00 | `acondicionador` Acondicionador para cabello | `acondicionador.pantene-170ml` | ml |
| 31 | conjunto | Detergente Roma 1 kg | 1 | 35.00 | `detergente-polvo` Detergente en polvo | `detergente-polvo.roma-1kg` | g |
| 32 | conjunto | Jabón Salvo 1.2 l | 2 | 122.00 | `jabon-trastes` Jabón líquido para trastes | `jabon-trastes.salvo-1200ml` | ml |
| 33 | conjunto | Mayonesa McCormick chipotle | 1 | 44.00 | `mayonesa-chipotle` Mayonesa con chipotle | `mayonesa-chipotle.mccormick-310g` | g |
| 34 | conjunto | Soya Maggi | 1 | 39.00 | `salsa-soya` Salsa de soya | `salsa-soya.maggi-290ml` | ml |
| 35 | conjunto | Catsup Heinz | 1 | 29.00 | `catsup` Catsup | `catsup.heinz-397g` | g |
| 36 | conjunto | Mayonesa McCormick 420 g | 1 | 49.00 | `mayonesa` Mayonesa | `mayonesa.mccormick-420g` | g |
| 37 | conjunto | Salsa Inglesa Crosse & Blackwell | 1 | 70.00 | `salsa-inglesa` Salsa inglesa | `salsa-inglesa.crosse-blackwell-290ml` | ml |
| 38 | conjunto | Servilletas Pétalo 420 | 1 | 50.00 | `servilletas` Servilletas de papel | `servilletas.petalo-420` | pieza |
| 39 | conjunto | Servitoalla Pétalo | 3 | 102.00 | `toallas-papel` Toallas de papel | `toallas-papel.petalo-mega-jumbo-270` | pieza |
| | **Total** | **39 renglones** | **61** | **3,094.50** | | | |

Presentaciones de Walmart que no están en `cuentas` (sustitutos de `data.json: sustitutos[*].walmart`): `arroz-grano-largo.sos-uruguayo-1kg` (Arroz SOS grano largo uruguayo 1 kg), `sustituto-crema-cafe.coffee-mate-640g` (Coffee Mate en polvo 640 g) y `leche.santa-clara-entera-6x1l` (Leche Santa Clara entera caja 6 pzas 1 l).
