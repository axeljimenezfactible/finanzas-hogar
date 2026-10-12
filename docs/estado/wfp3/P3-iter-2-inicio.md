# P3 · Vuelta 2 · Área «inicio» · Informe de revisión (P3-07)

Revisor: `revisor-inicio-r2` (independiente, vuelta 2 de 3, no corrige). Fecha: 12-oct-2026.

**Objeto.** Los 8 artboards de `docs/design/maquetas/lienzo/project/inicio/`: P-01 Inicio en sus 4 subsecciones (Resumen, Detalle, Tendencias y Calendario), cada una en móvil 412 y escritorio 1440. También el fragmento `docs/design/maquetas/lienzo/fragmentos/inicio.json`.

**Criterios.**
- P3 1.3: P3-04, P3-06, P3-07, P3-08 y P3-12.
- U5-04, U5-06 y U5-07, en lo que se puede verificar en el markup.
- DESIGN.md: §2.4, §5.2, §7.4, §7.6, §7.7, §8.1, §8.5, §9.3 y §10.5, más `contraste.csv` y `contraste-componentes.csv`.
- U1 §2.1 y §2.2, filas de P-01.
- U3: §1 a §5 y §7.
- `referencia-tipo/format.md`, leído entero y tomado como dato.
- `datos-sinteticos.json` y `LEEME.md`.

**Método.** Escribí scripts propios en `datos/tmp/rev-inicio-r2/` y ya los borré:
- `forma.py`: forma del tipo, tokens, P3-12, nombres, raíz y enlaces del markup.
- `rv.js`: corre en `node` el `renderVals()` de cada artboard con un `DCLogic` de prueba.
- `texto.py`: dinero, cifras tabulares, inglés, fechas ISO, enlaces de `renderVals()` y montos contra el JSON.
- `contraste.py`: una cascada simplificada con especificidad que expande `sc-for` y `sc-if` con los valores reales de `renderVals()`. Resuelve los `var()` con `tokens.css` en claro y en oscuro.
- `cruce.py`: cruza cada par con las tablas oficiales.
- `tactil.py`: altos de los controles en móvil.
- `anid.py`: controles anidados o sin nombre.

Todos los scripts de Python corrieron con `python3 -I`. Comparé con `diff` contra la copia de la vuelta 1 (`docs/design/salida/iter-1/inicio/`). No usé navegador, red ni la herramienta Artifact.

**Veredicto: NO APROBADO.** Hay 1 hallazgo alto (A-2, nuevo), 0 medios y 9 bajos.
- El alto de la vuelta 1 (A-1) quedó **cerrado** en los artboards.
- A-2 es de la misma clase que A-1: el artboard contradice una decisión de la puerta 2 (Q2-17 C), y el origen está en el material compartido (`LEEME.md` regla 4 y `datos-sinteticos.json`), no en el criterio del autor.
- La corrección es corta: 4 textos y 4 atributos ARIA en 2 artboards, más el ajuste de datos.

## 1. Cierre de los hallazgos de la vuelta 1

| Id | Sev. v1 | Hallazgo | Estado | Evidencia (vuelta 2) |
|---|---|---|---|---|
| A-1 | alta | Insignias «Atrasada» y «En ruta» en las metas, contra Q2-18 A | **Cerrado** en los artboards. En los datos sigue abierto (ver B-11) | El `diff` contra `iter-1` muestra exactamente 6 líneas cambiadas, y nada más cambió en los 8 artboards ni en `inicio.json`. Las líneas son resumen-m l. 107, resumen-e l. 144, detalle-m l. 150 y 156, y detalle-e l. 189 y 195. Las 6 insignias dicen ahora «En curso», con `bd bi` (rol información) y el icono `info` (círculo con i), como pide U3 tarjeta 6. «Atrasada» y «En ruta» aparecen 0 veces en Inicio. Las dos metas tienen fecha futura (17-ene-2027 y 17-abr-2027) y no están cumplidas, así que «En curso» es el estado correcto |
| B-1 | baja | Trazo de 1.5 a 20 px (sale de 1.25 px) | Abierto | `.ic{…stroke-width:1.5…}` sigue en los 8 helmets |
| B-2 | baja | «Pendiente» de D-3 sin `circle-dashed` | **Cerrado: falso positivo de la vuelta 1** | Cada fila de D-3 empieza con un círculo con `stroke-dasharray:3 2.2` en `--ui-icon-default` (detalle-m l. 112 y su equivalente en escritorio), sin cambios desde la vuelta 1. El estado lleva icono y texto (U3 §1.4). El icono está al inicio de la fila y no junto al texto. Juntarlos es opcional |
| B-3 | baja | Marca «viejo» en D-5 | Abierto | Sin cambios. Falta la decisión |
| B-4 | baja | «Meta 2,364» y «Meta 1,900» en lugar de «Meta de Axel» y «Meta de Ximena» | Abierto | tendencias-m y tendencias-e siguen con «Meta 2,364», «Meta 1,900» y «Meta (mismo color, punteada)» |
| B-5 | baja | Los ejemplos de error no traen el enlace de la tarjeta | Abierto | resumen-e l. 169-175 y detalle-e l. 240-246: título, mensaje y «Reintentar», sin «Ver presupuesto» ni «Ver todo» |
| B-6 | baja | `aria-label` en un `<span>` sin rol («Pendiente confirmar») | Abierto | resumen-m l. 83 y resumen-e l. 120, sin cambios |
| B-7 | baja | Alturas `h` estimadas sin navegador | Abierto (orquestador) | Sin cambios. Se mide al publicar |
| B-8 | baja | Consumo del autor en la vuelta 1 por encima del tope | No aplica en esta vuelta | La corrección declara unas 10 acciones y unos 4 KB modificados |

## 2. Hallazgos de esta vuelta

| Id | Sev. | Regla | Artboards (líneas) | Hallazgo | Corrección |
|---|---|---|---|---|---|
| **A-2** | **alta** | E y D: U3 §3, fila D-6; decisión **Q2-17 C** (`docs/estado/bloqueos/S2-metas.md`, RESUELTO); `S2-interfaz.md` 1.3 l. 295 | detalle-m l. 151-152 y 157-158; detalle-e l. 190-191 y 196-197 | D-6 «Metas» muestra el avance en entero: «42 % · $6,150.00 de $14,500.00» y «78 % · $14,100.00 de $18,000.00», con `aria-valuenow="42"`/`"78"` y `aria-valuetext="42 % de la meta"`. U3 D-6 pide «el % en texto, **con un decimal**» (`progress_tenths_pct`, ROUND_HALF_UP). Q2-17 C dice: «un decimal en Metas y entero en la tarjeta de Inicio». U3 aclara que D-6 es la métrica M-16 «Metas» y que el entero (`progress_pct`, M-07) va **solo** en la tarjeta 6. La tarjeta 6 del Resumen está bien con «42 %». **Origen**: `LEEME.md` regla 4 dice «Los porcentajes de avance son enteros (ROUND_HALF_UP)», y `datos-sinteticos.json` solo trae `ahorro.metas[].avance_pct` = 42 y 78, sin décimas. El autor siguió la guía de datos, pero la pantalla contradice la decisión | **Autor**: «42.4 % · $6,150.00 de $14,500.00» y «78.3 % · $14,100.00 de $18,000.00». Las cuentas son 615000 ÷ 1450000 = 42.41 % y 1410000 ÷ 1800000 = 78.33 %, con ROUND_HALF_UP a décimas. También `aria-valuenow="42.4"`/`"78.3"` y `aria-valuetext="42.4 % de la meta"`/`"78.3 % de la meta"`, en los dos artboards. **Orquestador**: agregar las décimas al JSON (p. ej. `avance_decimas_pct`: 424 y 783, o `progress_tenths_pct` con el nombre de C1). Corregir `LEEME.md` regla 4: «entero en la tarjeta de Inicio y un decimal en Metas». Revisar P-09 y P-26.13 (otras áreas): P-26.13 móvil muestra «42 %» y «78 %» |
| B-9 | baja | U5-06 y DESIGN §5.2 («Solo dos tamaños»; `size.icon.sm` = 20 px para insignias) | los 8 (l. 28) | `.bd .ic{width:var(--space-4);height:var(--space-4)}` dibuja los iconos de las insignias a 16 px. DESIGN dice que «un icono de 16 px a 1.5 px pierde detalle». Además, usa un token de espaciado como tamaño de icono | `.bd .ic{width:var(--size-icon-sm);height:var(--size-icon-sm)}`, o quitar la regla: `.ic` ya mide 20 px. La insignia mide 20 px de alto (`--size-badge-height`), así que el icono cabe |
| B-10 | baja | U5-07 (lista automática de palabras en inglés: `Error`) | resumen-e l. 170, detalle-e l. 241, tendencias-e l. 143, calendario-e l. 146 | La etiqueta del ejemplo de estado dice «Error». Es español válido, pero el script de U5-07 la contaría como inglés | «Con error» o «Falla al cargar», junto a «Vacío» y «Carga» |
| B-11 | baja (para Inicio) | P3-12 («todos los autores usan el mismo archivo… para que cuadren entre sí») y Q2-18 A | `datos-sinteticos.json` l. 13 y 33; `LEEME.md` l. 29 y 56 | El JSON sigue con `ahorro.metas[].estado` = `atrasada`/`en_ruta`, `etiqueta_estado` = «Atrasada»/«En ruta» y `inicio.tarjetas.ahorro.next_goal.status` = `behind`. `LEEME.md` dice «una en ruta, una atrasada» y «la meta más próxima es la atrasada». Inicio ya va por delante del JSON (dice «En curso»). En cambio, `presupuesto-p09-ahorro-{m,e}` y `ajustes-p26-s13-salud-metas-{m,e}` todavía traen «Atrasada»/«En ruta» en sus datos. El lienzo no cuadra entre áreas: la misma meta sale «En curso» en Inicio y «Atrasada» en P-09 | Orquestador: `estado` → `en_curso`/`in_progress` y `etiqueta_estado` → «En curso» en las dos metas; `next_goal.status` → `in_progress`; corregir `LEEME.md` l. 29 y 56; volver a correr `verificar_sinteticos.py`. Para Inicio es baja porque sus artboards ya cumplen. Para P-09 y P-26.13 es la misma contradicción de Q2-18 A |

## 3. Tabla artboard × regla

Claves:
- **F**: forma del tipo.
- **T**: tokens.
- **C**: contraste (U5-04 y P3-07).
- **D**: datos sintéticos (P3-12).
- **E**: fidelidad a U1 y U3.
- **A**: accesibilidad (U5-06 y craft).
- **I**: idioma y formatos (U5-07).
- **R**: raíz móvil fija o escritorio fluido.
- **L**: enlaces.

| Artboard | F | T | C (texto / UI, por tema) | D | E | A | I | R | L | Falladas (Id) | Corrección para la vuelta 3 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| `inicio-p01-resumen-m` | ✔ | ✔ (B-1, B-9) | ✔ 12 / 3 | ✔ | ✔ 6 tarjetas; A-1 cerrado | ✔ (B-6) | ✔ | ✔ 412 × 1840 | ✔ 14 | B-1, B-6, B-7, B-9 | Opcional: trazo 1.8, iconos de insignia de 20 px y nombre accesible de «Pendiente confirmar» |
| `inicio-p01-resumen-e` | ✔ | ✔ (B-1, B-9) | ✔ 14 / 3 | ✔ | ✔ 6 tarjetas; Estados ✔ (B-5) | ✔ (B-6) | ✔ (B-10) | ✔ fluida, `fill` | ✔ 15 + 8 | B-1, B-5, B-6, B-9, B-10 | Lo mismo que en móvil; «Ver presupuesto» en el ejemplo de error; «Con error» |
| `inicio-p01-detalle-m` | ✔ | ✔ (B-1, B-9) | ✔ 13 / 10 | ✔ | **✘ A-2** (B-3) | ✔ | ✔ | ✔ 412 × 4250 | ✔ 17 | **A-2**; B-1, B-3, B-7, B-9 | **D-6 con un decimal: «42.4 %» y «78.3 %», también en `aria-valuenow` y `aria-valuetext`** |
| `inicio-p01-detalle-e` | ✔ | ✔ (B-1, B-9) | ✔ 14 / 10 | ✔ | **✘ A-2** (B-3, B-5) | ✔ | ✔ (B-10) | ✔ fluida, `fill` | ✔ 18 + 8 | **A-2**; B-1, B-3, B-5, B-9, B-10 | **Lo mismo que en móvil**; «Ver todo» en el ejemplo de error |
| `inicio-p01-tendencias-m` | ✔ | ✔ (B-1, B-9) | ✔ 7 / 21 | ✔ | ✔ (B-4) | ✔ | ✔ | ✔ 412 × 2450 | ✔ 9 | B-1, B-4, B-7, B-9 | «Meta de Axel» y «Meta de Ximena» |
| `inicio-p01-tendencias-e` | ✔ | ✔ (B-1, B-9) | ✔ 9 / 21 | ✔ | ✔ (B-4) | ✔ | ✔ (B-10) | ✔ fluida, `fill` | ✔ 9 + 8 | B-1, B-4, B-9, B-10 | Lo mismo; «Con error» |
| `inicio-p01-calendario-m` | ✔ | ✔ (B-1, B-9) | ✔ 9 / franjas | ✔ | ✔ | ✔ | ✔ | ✔ 412 × 3050 | ✔ 14 + 6 | B-1, B-7, B-9 | Opcional |
| `inicio-p01-calendario-e` | ✔ | ✔ (B-1, B-9) | ✔ 11 / franjas | ✔ | ✔ | ✔ | ✔ (B-10) | ✔ fluida, `fill` | ✔ 14 + 8 | B-1, B-9, B-10 | «Con error» |

En **L**, el primer número cuenta los `href` distintos del markup. El segundo cuenta los distintos de `renderVals()`, ya evaluado en `node`.

Lo que cumplen los 8 artboards (todo medido por script):

- **Forma (format.md y encargo).**
  - La línea 1 es `<!doctype html>`. Llevan `lang="es-MX"`, `meta charset="utf-8"` y un `<title>` propio por pantalla.
  - `./support.js` va antes de `../ds/finanzas-hogar/tokens.css`, dentro de `<head>`.
  - Hay un solo `<x-dc>`, que abre con `<helmet><style>`. El `body` usa tokens para `margin:0`, fondo, color y fuente, y hay reglas `a{}` y `a:hover{}`.
  - `<script type="text/x-dc" data-dc-script data-props='…'>` va después de `</x-dc>`. Su JSON parsea y `$preview` es igual a `w`×`h` del fragmento en los 8. Lleva `class Component extends DCLogic` y `renderVals()`.
  - `html.parser`: 0 etiquetas sin cerrar y 0 cierres inesperados.
  - Los huecos son solo rutas (0 expresiones). Los `sc-for` y `sc-if` traen sus `hint-*`. `renderVals()` corre sin error en `node` en los 8.
  - Nombres válidos. Los 90 stems del lienzo son únicos.
- **Raíz.**
  - Móvil: `width:412px;box-sizing:border-box`, `min-height` igual a `h` (≥ 915) y sin `expand`.
  - Escritorio: sin ancho en px, `flex-wrap:wrap`, barra lateral `flex:1 1 var(--layout-sidebar-width)` y `main` con `flex:999 1 560px;min-width:0`. Sin `sticky`, `100vh` ni `position:fixed`. `max-width` en el contenido y `"expand":"fill"` en el fragmento.
- **Tokens.**
  - 0 hex fuera de `data-props`, y 0 `rgb()`/`hsl()`.
  - Entre 54 y 70 `var(--x)` distintos por archivo, todos en `tokens.css`.
  - En `tokens.json` están todos salvo `font-sans`, que es `type.families.sans`: según su `meta.note`, los nombres de `tokens.json` son los de `tokens.css`.
  - El único `color-mix` es el área de T-2, que DESIGN §8.1 permite.
- **Red.** 0 iframes, 0 `innerHTML`, 0 `fetch` y 0 URL externas (la única es el `xmlns` de SVG).
- **Controles.**
  - 0 `div`/`span` clicables, 0 controles anidados (ningún `button` dentro de `<a>`), 0 controles sin nombre y 0 `<a>` sin `href`.
  - Los 207 `<svg>` llevan `aria-hidden="true"` o `role="img"` con `aria-label`.
  - Inicio no tiene campos de formulario.
- **Contenido.** 0 emoji, 0 barras de estado o teclados falsos y 0 fechas ISO visibles.
- **Dinero.**
  - Hay 281 apariciones de dinero. Fuera de las etiquetas de eje, todas tienen la forma `$1,234.56`, y hay 0 negativos con guion (todos con U+2212).
  - Todas llevan `tabular-nums`, sea por `.n`, `.ax` o `.lb`.
  - El compacto («$10 mil», «$500», «$0») aparece solo en etiquetas de eje, con la clase `.ax`, como pide DESIGN §8.5.

## 4. Fidelidad (pantalla → U1 y U3)

| Bloque | Elemento | m | e |
|---|---|---|---|
| Página | «Inicio» con `house`. Pestañas Resumen · Detalle · Tendencias · Calendario en orden fijo (`role="tablist"`). Única acción con relleno, «Registrar gasto»: «+» en móvil o botón en la barra superior en escritorio. «Foto de ticket» en «⋯» | ✔ | ✔ |
| Resumen | **Exactamente 6 tarjetas-enlace**, en el orden de U3 §2, a P-06, P-07, P-18, P-10, P-14 y P-09. En escritorio hay además 3 ejemplos de estado, que no son tarjetas del Resumen | ✔ 6 | ✔ 6 |
| T1 a T5 | Contenido, iconos (`wallet`, `banknote`, `calendar-clock`, `utensils` y `shopping-cart`), flechas `arrow-down-left`/`arrow-up-right`, insignias «Ingreso esperado» (info), «Pendiente confirmar» (info) y «2 precios viejos» (advertencia), y enlaces «Ver …» | ✔ | ✔ |
| T6 | `piggy-bank`; «$41,730.00» en rol ahorro con «ahorrado»; «Este mes $1,750.00 de $2,000.00»; «Laptop nueva · 42 % · 17-ene-2027» con el entero que pide Q2-17 C; insignia «En curso» (info, `info`) | ✔ (A-1 cerrado) | ✔ (A-1 cerrado) |
| Resumen | Estados de U1 §2.2 y U3 §1.3: vacío con «—», «Aún no hay presupuesto para octubre», «Llevas $20,525.78 gastado» y «Capturar presupuesto» → P-26.02; carga con esqueleto en `ui-skeleton-on-card` y el aviso de 10 s; error con el texto exacto, `circle-alert` en `ui-role-error` y «Reintentar» | n/a | ✔ (B-5) |
| Detalle | D-1 a D-5 y D-7 como en la vuelta 1; enlace al pie «Ver todos los gastos con estos filtros» → P-02 | ✔ | ✔ |
| D-6 | Nombre, barra en `--ui-role-savings`, «[ahorrado] de [meta]», aportación necesaria, fecha meta e insignia «En curso» | ✔ | ✔ |
| D-6 | **% de avance con un decimal** (Q2-17 C) | **✘ A-2** | **✘ A-2** |
| Tendencias | T-1 a T-4 en `<figure>`/`<figcaption>`, con resumen, leyenda, «Ver como tabla», etiquetado selectivo y formato completo en las etiquetas. Las sumas de T-2 ($242,520.65 y $257,740.56) cuadran con Σ `proyeccion.meses[].devengado_total_cents` y `efectivo_total_cents`, que recalculé | ✔ (B-4) | ✔ (B-4) |
| Calendario | Agenda, mes de lunes a domingo, franjas de quincena (continua `--ui-border-selected` y punteada `--ui-viz-baseline`), panel del día y filtros | ✔ | ✔ |
| Estados | Fila «Estados» con vacío, carga y error en los 4 artboards de escritorio, con los textos de U3 §7 | n/a | ✔ |

## 5. Contraste (U5-04 y P3-07)

- **Pares evaluados**: 108 combinaciones distintas (tema × primer plano × fondo), que son 54 por tema.
  - 35 por tema salen de la cascada: 22 de texto y 13 de UI (marcas de gráfica, puntos y cuadros de persona o tienda, medidores). Contando por clase de marca, son 32 usos de UI por tema.
  - 19 por tema son pares extra que la cascada no ve: estados `:hover`, foco, iconos de solo icono, bordes de control, error, `--ui-role-savings` sobre el carril y franjas del calendario.
- **Bajo umbral: 0** en claro y 0 en oscuro.
  - La rejilla (`--ui-viz-grid`), el carril (`--ui-viz-track`) y el halo de 2 px de los marcadores (trazo `--ui-surface-card`, DESIGN §8.1 «separadas por 2 px del color de la superficie») son decorativos y quedan fuera del umbral.
- **Mínimos**:
  - Texto: `--ui-role-success` sobre su superficie da 4.81 en claro; `--ui-role-income` sobre tarjeta, 4.80.
  - UI: `--ui-store-sams-club` sobre tarjeta, 3.06 en claro; `--ui-store-market`, 3.06 en oscuro; `--ui-person-axel`, 3.10 en oscuro; `--ui-chart-2`, 3.15 en claro; `--ui-chart-1`, 3.22 en oscuro.
- **Cruce con las tablas oficiales**: 104 de 108 tienen fila, por hex y tema, en `contraste.csv` o `contraste-componentes.csv`, y coinciden con mi cálculo al ±0.01, truncado como en DESIGN §10.5. Las 4 sin fila:
  - El halo de los marcadores (decorativo).
  - `--ui-role-savings` sobre `--ui-viz-track` (relleno del medidor de D-6): 4.34 en claro y 4.43 en oscuro, sobre el umbral de 3:1.
- **Muestra de 10 pares (claro · oscuro)**:
  - texto sobre página: 11.53 · 14.07
  - secundario sobre tarjeta: 5.04 · 7.97
  - enlace sobre tarjeta: 5.16 · 7.94
  - texto sobre botón primario: 5.77 · 10.47
  - info sobre su superficie: 4.96 · 6.12
  - advertencia sobre su superficie: 5.12 · 6.09
  - gasto sobre tarjeta: 5.33 · 5.90
  - ahorro sobre tarjeta: 4.84 · 6.42
  - Ximena sobre tarjeta: 4.85 · 6.93
  - línea base sobre tarjeta: 3.47 · 6.09

## 6. P3-08: decisiones

| Decisión | Resultado |
|---|---|
| Resumen con ≤ 6 tarjetas | ✔ exactamente 6 en m y en e |
| 14 secciones de Ajustes | n/a en Inicio (la barra lateral enlaza Ajustes con `ajustes-p26-s01-cuenta-roles-e`) |
| Barra lateral colapsable | ✔ botón «Colapsar barra lateral» en los 4 de escritorio |
| Español de México | ✔. Supuesto: «Kueski Pay» es nombre propio. «Error» es español (B-10) |
| Dinero `$1,234.56` tabular | ✔ 281 apariciones; fuera de los ejes compactos, todas con formato completo; U+2212 en los negativos |
| Fechas `10-oct-2026` | ✔ «17-ene-2027», «17-abr-2027» y «19-oct-2026»; con día de la semana, «Vie 30 oct 2026»; 0 ISO visibles |
| Paleta de C2 (índigo 280°) | ✔ solo tokens; un `color-mix` permitido |
| Notion como única referencia | ✔ barra lateral jerárquica, pestañas y tarjetas de borde fino con `--shadow-0` |
| Decisiones de la puerta 2 en Inicio | Q2-18 A ✔ (A-1 cerrado). **Q2-17 C ✘ en D-6 (A-2)** |

## 7. Líneas de evidencia

- **P3-04**: este revisor hizo 0 llamadas a navegador y 0 a Artifact. El journal del flujo lo revisa el orquestador.
- **P3-06**: los 8 artboards de P-01 (4 subsecciones × 2 anchos) están en disco y en el fragmento: `boards` = `order` = disco, sin duplicados, con `title`, `page:"inicio"` e `is_interactive:false`, y sin `x`/`y`. Faltantes: 0. La copia de la vuelta 1 está en `docs/design/salida/iter-1/inicio/` (9 archivos).
- **P3-07**: este informe. Recalculé 108 combinaciones de contraste y 104 coinciden con las tablas oficiales.
- **P3-08**: 8 de 8 en la lista de la sección 6. Contradicción de Q2-17 C en A-2.
- **P3-12**:
  - El grep oficial de las 21 cifras canon, con y sin coma, sobre los 8 artboards y `inicio.json` da **0**. En centavos, 0. Credenciales, 0.
  - El comando de fechas lista 15 coincidencias inevitables de quincena (`2026-10-30`, `2026-11-15`, `2026-11-30`, …). En Inicio, el 30 oct, el 15 nov y el 30 nov 2026 van siempre con montos sintéticos ($1,861.58, $954.60, $906.98, $803.27 y +$14,780.00), nunca con cifras del reporte ni de Kueski.
  - Todos los montos visibles de los 8 artboards (markup y `renderVals()`) están en `datos-sinteticos.json`, salvo dos grupos. Los ejes compactos. Y las 2 sumas de T-2, que recalculé: $242,520.65 = Σ `proyeccion.meses[].devengado_total_cents` y $257,740.56 = Σ `efectivo_total_cents`.
  - Muestra de 10 trazas:
    - $7,060.22 ← `presupuesto.resumen.restante_cents`
    - $6,045.50 ← `inicio.tarjetas.flujo_quincenal.closing_balance_cents`
    - $1,861.58 ← `creditos.lineas[0].proximo_pago.monto_cents`
    - $6,881.60 ← `creditos.lineas[0].disponible_cents`
    - $11,440.00 ← `creditos.lineas[0].disponible_a_hoy_cents`
    - $1,062.20 ← `plan_semana.totales.costo_semana_cents`
    - $1,551.70 ← `canasta.por_tienda[1].total_cents`
    - $41,730.00 ← `ahorro.ahorrado_total_cents`
    - $2,783.34 ← `ahorro.metas[0].aportacion_necesaria_cents`
    - $42,744.75 ← `reporte_mensual.creditos.deuda_viva_cierre_cents`
  - Fechas: 17-ene-2027 ← `ahorro.metas[0].fecha_meta`; «Lun 19 al Dom 25 oct» ← `meta.semana_en_curso`; «Vie 30 oct 2026» ← `inicio.tarjetas.proximo_pago.next_payment.due_date`.
- **U5-04**: 54 pares por tema, 0 bajo umbral (sección 5).
- **U5-06**:
  - Cumple:
    - Tokens de rol correctos y siempre con icono y texto.
    - Escala tipográfica del sistema (`title-1/2/3`, `small`, `caption`) y cifras tabulares.
    - Espaciado `--space-*`, radios `--radius-sm/md/full` y `--shadow-0` con borde.
    - Barra lateral colapsable y 6 tarjetas.
    - Objetivos ≥ 44 px en móvil: `tactil.py` resolvió `height`/`min-height` de los 70 controles de los 4 móviles. Hay 0 por debajo; los únicos sin alto declarado son las tarjetas-enlace `.cl`, de más de 100 px.
    - `:focus-visible` con `--ui-border-focus`.
    - Series en orden fijo con marcador, trazo y etiqueta.
  - Falla (baja):
    - B-1: trazo de 1.5 a 20 px.
    - B-9: iconos de insignia de 16 px.
- **U5-07**:
  - 0 fechas ISO visibles.
  - Dinero y fechas en formato.
  - Lista de inglés: solo «Error», que es español (B-10).
  - «Sin conexión» no se dibuja en Inicio porque el encargo no lo pide para esta área.

## 8. Texto con aspecto de instrucción

Ninguno. Busqué «ignora», «instrucción», «Claude», «agente», «prompt», «system», «olvida», «ejecuta» y «orquestador» en los 8 artboards y en `inicio.json`: 0 coincidencias. `format.md` y el resto del material del tipo los leí como dato.

## 9. Correcciones para la vuelta 3 (en orden)

1. **A-2 (orquestador y autor).**
   - Orquestador: agregar las décimas de avance de las metas a `datos-sinteticos.json` (424 y 783) y corregir `LEEME.md` regla 4.
   - Autor: en `inicio-p01-detalle-m` (l. 151-152 y 157-158) y `inicio-p01-detalle-e` (l. 190-191 y 196-197), «42.4 %» y «78.3 %» en el texto, en `aria-valuenow` y en `aria-valuetext`. La tarjeta 6 del Resumen se queda con «42 %».
2. **B-11 (orquestador).** Estados de meta en `datos-sinteticos.json` y `LEEME.md` según Q2-18 A, y aviso a las áreas presupuesto (P-09) y ajustes (P-26.13), que aún dicen «Atrasada»/«En ruta». Luego `verificar_sinteticos.py`.
3. Bajas opcionales del autor:
   - B-9: iconos de insignia a `--size-icon-sm`.
   - B-1: trazo `var(--icon-stroke-user-units-sm)` a 20 px.
   - B-4: «Meta de Axel» y «Meta de Ximena».
   - B-5: el enlace de la tarjeta en los ejemplos de error.
   - B-10: «Con error».
   - B-6: nombre accesible de «Pendiente confirmar».
   - B-3: decidir la marca «viejo».
4. B-7: el orquestador mide los altos al publicar.
