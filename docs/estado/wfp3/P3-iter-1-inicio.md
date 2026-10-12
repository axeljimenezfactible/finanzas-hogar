# P3-iter-1-inicio · Revisión del área «inicio», vuelta 1 de 3

Fecha: 12-oct-2026 · Revisor: `revisor-inicio-r1` (independiente; no corrige) · Criterio: P3-07

Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/inicio/`: P-01 Inicio en sus 4 subsecciones (Resumen, Detalle, Tendencias y Calendario), cada una en móvil y escritorio. También revisé el fragmento `docs/design/maquetas/lienzo/fragmentos/inicio.json`. Los comparé con:

- `docs/criterios/P3-claude-design.md` 1.3: P3-04, P3-06, P3-07, P3-08 y P3-12.
- `docs/criterios/U5.md`: U5-04, U5-06 y U5-07, en lo que se puede comprobar en el markup.
- `DESIGN.md`: §2.10, §5.2, §7.6, §7.7, §8.1 a §8.8, §9.3 y §10.5.
- U1: §0.4, §0.5, §1.2 a §1.4, §2.1, §2.2, §2.4, §3.3 y §6.4 (fila P-01).
- U3 completo para Inicio: §1 a §10.
- `referencia-tipo/format.md`.

Todo lo verifiqué con scripts propios en `python3 -I` y con `node` para evaluar `renderVals()`. Los scripts estaban en `datos/tmp/rev-inicio/` y ya los borré. Además leí el markup. No usé navegador ni la herramienta Artifact.

**Veredicto: NO APROBADO.** Hay 1 hallazgo alto, 0 medios y 8 bajos. La causa del alto está en el archivo de datos compartido, no en el trabajo del autor. Aun así, contradice una decisión en la pantalla principal y por eso hay que corregirlo en la vuelta 2.

## Resumen de hallazgos

| Id | Sev. | Regla | Artboards | Hallazgo |
|---|---|---|---|---|
| A-1 | alta | E y D (U3 tarjeta 6, D-6 y §2; decisión Q2-18 A de la puerta 2) | resumen-m (l. 107), resumen-e (l. 144), detalle-m (l. 150 y 156), detalle-e (l. 189 y 195) | La meta «Laptop nueva» lleva la insignia **«Atrasada»**, con el rol de advertencia y `triangle-alert`. «Viaje a Oaxaca» lleva **«En ruta»**. U3 manda los tres estados de Q2-18 A: «En curso» (`in_progress`, rol información, `info`), «Vencida» y «Cumplida», y dice textualmente: «no hay "atrasada"». Las dos metas tienen fecha futura y no están cumplidas, así que las dos van como «En curso». **El origen está en `datos-sinteticos.json`**: `inicio.tarjetas.ahorro.next_goal.status = "behind"` y `etiqueta_estado = "Atrasada"`, y en `ahorro.metas[].estado`, que dice `atrasada` y `en_ruta`. Esos campos siguen el enum viejo de C1 v1.3 (`on_track`/`behind`), que U3 §2 tarjeta 6 marca como solicitud pendiente a C1. El autor copió bien los datos, pero la pantalla contradice la decisión |
| B-1 | baja | U5-06 y DESIGN §5.2 (tabla de iconos) | los 8 | `.ic` dibuja los iconos de 20 px con `stroke-width:1.5` en unidades del `viewBox` de 24. El trazo sale de 1.25 px. DESIGN pide `icon.stroke.user-units.sm` = 1.8 a 20 px para que salgan 1.5 px dibujados. Además, el 1.5 va como literal y no como token |
| B-2 | baja | E (U3 §1.4: «siempre con icono y texto») | detalle-m (l. 112), detalle-e (l. 151) | En D-3, el estado «Pendiente» de cada pago es solo texto, sin `circle-dashed`. En el Calendario de escritorio sí lleva el icono |
| B-3 | baja | E (U3 D-5, marca «viejo») | detalle-m, detalle-e | El reporte del autor (supuesto 7) dice que los datos no traen la marca «viejo», pero `canasta.detalle[].precios[tienda].viejo` sí existe. Para los 5 artículos, el precio de la tienda ganadora no es viejo, así que omitir la marca junto a la ganadora es correcto. Sin embargo, la diferencia de «Pechuga de pollo» (Bodega Aurrera) y la de «Tortilla de maíz» (Walmart) se calculan contra precios viejos. Conviene decidir si la marca va en esos casos y corregir el supuesto |
| B-4 | baja | E (U3 T-4b) | tendencias-m (l. 97-98), tendencias-e (l. 133-134) | Las líneas de meta dicen «Meta 2,364» y «Meta 1,900», y la leyenda «Meta (mismo color, punteada)». U3 pide la etiqueta «Meta de [nombre]»: «Meta de Axel» y «Meta de Ximena» |
| B-5 | baja | E (U3 §1.3: en el error «la tarjeta conserva su título y su enlace») | resumen-e (l. 170-174), detalle-e (l. 241-245) | Los ejemplos de error de la fila «Estados» muestran el título, el mensaje y «Reintentar», pero no el enlace de la tarjeta («Ver presupuesto» o «Ver todo») |
| B-6 | baja | A (U3 §1.2 y tarjeta 3: nombre accesible) | resumen-m (l. 83), resumen-e (l. 120) | El nombre accesible «Reposición de capital pendiente de confirmar con Kueski» va en `aria-label` sobre un `<span>` sin rol. ARIA 1.2 no expone el nombre de un elemento genérico, así que el lector de pantalla oye solo «Pendiente confirmar». Las tarjetas-enlace tampoco declaran el nombre corto de U3 §1.2 (título + valor principal + estado): su nombre es todo el texto de la tarjeta |
| B-7 | baja | R (alto = contenido) | los 4 `-m` y los 4 `-e` | Los `h` del fragmento y el `min-height` de la raíz son estimados sin navegador. El propio autor lo dice. Resumen-m (1840) cuadra con mi estimación de unos 1,790 px. Detalle-m (4250) y calendario-m (3050) no los pude medir. Tiene que medirlos el orquestador en el lienzo |
| B-8 | baja | proceso (reparto fino) | — | El autor declara unos 189 KB de artboards escritos, más que el tope de 150 KB. No afecta al lienzo |

## Tabla artboard × regla

Claves:

- **F**: forma del tipo (format.md y encargo).
- **T**: tokens.
- **C**: contraste (U5-04 y P3-07).
- **D**: datos sintéticos (P3-12).
- **E**: fidelidad a U1 y U3.
- **A**: accesibilidad (U5-06 y craft).
- **I**: idioma y formatos (U5-07).
- **R**: raíz móvil fija o escritorio fluido.
- **L**: enlaces.

| Artboard | F | T | C | D | E | A | I | R | L | Falladas (Id) | Correcciones para la vuelta 2 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| `inicio-p01-resumen-m` | ✔ | ✔ (B-1) | ✔ 14 pares × 2 temas | ✔ | ✘ A-1 | ✔ (B-6) | ✔ | ✔ 412 | ✔ | A-1; B-1, B-6, B-7 | Insignia de la meta: «En curso» con `bd bi` e icono `info`, una vez corregidos los datos. Nombre accesible del aviso de reposición (texto oculto visualmente o `aria-label` en el enlace de la tarjeta). Trazo 1.8 a 20 px |
| `inicio-p01-resumen-e` | ✔ | ✔ (B-1) | ✔ 16 pares | ✔ | ✘ A-1 (B-5) | ✔ (B-6) | ✔ | ✔ fluida | ✔ | A-1; B-1, B-5, B-6 | Lo mismo que en móvil, y el enlace «Ver presupuesto» en el ejemplo de error |
| `inicio-p01-detalle-m` | ✔ | ✔ (B-1) | ✔ 15 pares | ✔ | ✘ A-1 (B-2, B-3) | ✔ | ✔ | ✔ 412 | ✔ | A-1; B-1, B-2, B-3, B-7 | D-6 con «En curso» en las dos metas (rol información). `circle-dashed` en «Pendiente» de D-3. Decidir la marca «viejo» de D-5 |
| `inicio-p01-detalle-e` | ✔ | ✔ (B-1) | ✔ 16 pares | ✔ | ✘ A-1 (B-2, B-3, B-5) | ✔ | ✔ | ✔ fluida | ✔ | A-1; B-1, B-2, B-3, B-5 | Lo mismo que en móvil, y «Ver todo» en el ejemplo de error |
| `inicio-p01-tendencias-m` | ✔ | ✔ (B-1) | ✔ 8 de texto + 6 de UI | ✔ | ✔ (B-4) | ✔ | ✔ | ✔ 412 | ✔ | B-1, B-4, B-7 | «Meta de Axel» y «Meta de Ximena» |
| `inicio-p01-tendencias-e` | ✔ | ✔ (B-1) | ✔ 11 de texto + 8 de UI | ✔ | ✔ (B-4) | ✔ | ✔ | ✔ fluida | ✔ | B-1, B-4 | «Meta de Axel» y «Meta de Ximena» |
| `inicio-p01-calendario-m` | ✔ | ✔ (B-1) | ✔ 7 pares | ✔ | ✔ | ✔ | ✔ | ✔ 412 | ✔ | B-1, B-7 | Trazo 1.8 a 20 px; medir el alto |
| `inicio-p01-calendario-e` | ✔ | ✔ (B-1) | ✔ 15 pares | ✔ | ✔ | ✔ | ✔ | ✔ fluida | ✔ | B-1 | Trazo 1.8 a 20 px |

Reglas que cumplen los 8 artboards:

- **Forma.**
  - Líneas obligatorias exactas y en orden: `<!doctype html>` en la línea 1, `lang="es-MX"`, charset, `<title>` propio, `./support.js` antes de `../ds/finanzas-hogar/tokens.css`, un solo `<x-dc>`, `<helmet><style>` y `<script type="text/x-dc" data-dc-script data-props='…'>` después de `</x-dc>`, con `class Component extends DCLogic` y `renderVals()`.
  - El JSON de `data-props` parsea y `$preview` es igual a `w`×`h` del fragmento.
  - `html.parser` no encuentra etiquetas sin cerrar ni cierres sin apertura.
  - `renderVals()` se ejecuta en `node`. Todos los huecos `{{ }}` tienen valor y ninguno es una expresión. Todos los `sc-for` y `sc-if` llevan su `hint-*`.
  - Los nombres de archivo son válidos y los stems son únicos en las 90 maquetas del lienzo (0 duplicados).
- **Tokens.**
  - 0 hex literales fuera de `data-props`. 0 `rgb()` y 0 `hsl()`.
  - El único `color-mix` es `.ar1` (área de T-2, 20 % de `--ui-chart-1` sobre `--ui-surface-card`). DESIGN §8.1 lo permite expresamente («derivado con `color-mix`, ningún color nuevo»).
  - Entre 54 y 70 `var(--x)` distintos por archivo, todos en `tokens.css`. En `tokens.json` también están todos, salvo `font-sans`, que es `type.families.sans` (la nota de `tokens.json` dice que sus nombres son los de `tokens.css`).
  - `url(#pt1…)` es una trama SVG interna, no un recurso de red.
- **Red.** 0 `innerHTML`, 0 iframes, 0 `fetch` y 0 URL externas.
- **Contenido.** 0 emoji, 0 barras de estado falsas, 0 fechas ISO visibles.
- **Controles.** 0 `div`/`span` con `onClick` o `role="button"`, 0 controles sin nombre y 0 controles anidados. Los 207 `<svg>` llevan `aria-hidden="true"`, o `role="img"` con `aria-label` en las gráficas: 5 por artboard de Tendencias y la barra de D-1, que está dentro de un `sc-for`.
- **Dinero.** Todo lo visible tiene la forma `$1,234.56` con `tabular-nums`; el compacto («$10 mil», «$1 mil») aparece solo en ejes, como pide DESIGN §8.5. Los negativos usan U+2212 («−$5,589.55») y los ingresos llevan «+».
- **Enlaces.** Los del markup y los de `renderVals()` (nav, agenda, celdas) apuntan a archivos existentes de la lista del lienzo y del mismo ancho, o a «#» (Buscar, G-04).

## Fidelidad (pantalla → elementos de U1 y U3)

| Bloque | Elemento | m | e |
|---|---|---|---|
| Página | Título «Inicio» con `house`; pestañas Resumen · Detalle · Tendencias · Calendario en orden fijo (en móvil, con desplazamiento horizontal) | ✔ | ✔ |
| Página | Acción primaria única «Registrar gasto»: «+» en la barra inferior (móvil) o botón con relleno en la barra superior (escritorio); secundaria «Foto de ticket» en «⋯» | ✔ | ✔ |
| Shell | Barra inferior Inicio · Comidas · + · Compras · Buscar (U1 §1.3); ☰ abre la barra lateral | ✔ | — |
| Shell | Barra lateral (U1 §1.2): espacio, menú del usuario, Buscar, Registrar gasto, 9 secciones en orden con ▸, tema en el pie y botón «Colapsar barra lateral»; Ajustes abre 01 Cuenta y roles | — | ✔ |
| Resumen | **Exactamente 6 tarjetas** en el orden fijo de U3 §2, en 1 columna (m) o 3 × 2 (e, `flex: 1 1 18rem` dentro de 76 rem) | ✔ 6 | ✔ 6 |
| T1 | «$7,060.22 te quedan»; «$20,525.78 de $27,586.00 · 74 %»; «Faltan 11 días del mes»; sin insignia (normal) → P-06 | ✔ | ✔ |
| T2 | «Del 16 al 31 oct»; «$6,045.50 al cierre»; «Entra +$15,520.00»; «Ya salió −$5,589.55 · Por pagar −$3,884.95»; «Ingreso esperado» → P-07 | ✔ | ✔ |
| T3 | «$1,861.58»; «Vie 30 oct 2026 · Kueski Pay»; «Hogar $1,026.25 · Karen $835.33»; «Faltan 9 días · 2 pagos este día»; «Disponible hoy $6,881.60 · Lectura B» con «Pendiente confirmar» → P-18 | ✔ (B-6) | ✔ (B-6) |
| T4 | «Lun 19 al Dom 25 oct»; «$1,062.20 costo de la semana»; «2 casillas sin armar»; una línea por persona con punto y nombre → P-10 | ✔ | ✔ |
| T5 | «Bodega Aurrera $1,551.70» con marca de tienda; «$86.60 menos que Walmart»; «23 de 24 artículos con precio»; «2 precios viejos» → P-14 | ✔ | ✔ |
| T6 | «$41,730.00 ahorrado» en rol ahorro; «Este mes $1,750.00 de $2,000.00»; meta «Laptop nueva · 42 % · 17-ene-2027» → P-09 | ✔ | ✔ |
| T6 | Insignia de estado de la meta con los tres estados de Q2-18 A | ✘ A-1 | ✘ A-1 |
| Resumen | Estados (U1 §2.2 y U3 §1.3): vacío («—», «Aún no hay presupuesto para octubre», «Llevas $20,525.78 gastado», «Capturar presupuesto» → P-26.02), carga (esqueleto, aviso de 10 s) y error (texto exacto y «Reintentar») | n/a | ✔ (B-5) |
| Detalle | Barra de vista: periodo Mes · 1.ª · 2.ª quincena; filtros Persona, Familia, Categoría, Tienda y Crédito (e) o «Filtrar» (m) | ✔ | ✔ |
| D-1 | 5 categorías con mayor % usado, barra real contra plan con marca vertical (DESIGN §8.2), «[real] de [plan] · [n] %», insignias «Rebasaste el plan» y «Vas en [n] %» | ✔ | ✔ |
| D-2 | 7 propiedades de la quincena con «Esperado» (de `proyeccion.flujo_por_quincena[1]`) | ✔ | ✔ |
| D-3 | Vencidos (ninguno) y próximos 5 hasta hoy + 60 (20 dic 2026), con «Hogar · Karen»; lecturas A y B lado a lado con «Activa» y el mínimo del horizonte; aviso de reposición | ✔ (B-2) | ✔ (B-2) |
| D-4 | Comidas de hoy por tiempo y persona, con «Fuera», kcal y costo | ✔ | ✔ |
| D-5 | 5 artículos con su tienda ganadora y la diferencia; total por tienda en orden fijo | ✔ (B-3) | ✔ (B-3) |
| D-6 | Metas con barra de avance en rol ahorro, «[ahorrado] de [meta]», aportación necesaria, fecha meta y estado | ✘ A-1 | ✘ A-1 |
| D-7 | Los 5 bloques de U3 §3.1, con la nota fija textual y «Estado de cuenta de Karen» → P-25; en escritorio, a todo el ancho | ✔ | ✔ |
| Detalle | «Ver todos los gastos con estos filtros» → P-02; escritorio en 2 columnas | ✔ | ✔ |
| Detalle | Estados de D-1 (vacío, 5 filas de esqueleto, error) | n/a | ✔ (B-5) |
| Tendencias | T-1 a T-4 en `<figure>`/`<figcaption>`, con resumen de una frase en la forma de U3, leyenda, «Ver como tabla» y etiquetado selectivo. T-1: columnas real/plan, plan con rayas a 45° y «Sin datos de nov a abr». T-2: capas desde el mismo origen, «Hoy», proyección desde abr 2027 y «Pagado por adelantado». T-3: extremos rotulados. T-4a/T-4b: Axel continua con círculo y Ximena 8 · 4 con cuadrado. 2 por fila en escritorio | ✔ (B-4) | ✔ (B-4) |
| Tendencias | Estados (vacío de T-1, carga sin ejes falsos, error) | n/a | ✔ |
| Calendario | Agenda de hoy a hoy + 30 (20 nov 2026), con separadores de quincena, encabezado «Vie 30 oct 2026», «Sin eventos» y «Ver 30 días más» (m). Mes de lunes a domingo con franjas de quincena (continua y punteada), hasta 3 eventos por día y «+n más», panel del día «Vie 30 oct 2026», filtros Tipos, Crédito y Familia, y navegación ‹ › «Hoy» (e) | ✔ | ✔ |
| Calendario | Estados (vacío, carga y error, con los textos de U3 §7) | n/a | ✔ |

## Contraste (U5-04 y P3-07)

Usé un motor de cascada propio sobre el `<style>` de `<helmet>`, los estilos en línea y los atributos `fill`/`stroke` de los SVG. Resuelve los `var()` con `tokens.css` en `:root` (claro) y en `[data-theme="dark"]`, y calcula WCAG 2.x.

- **Pares evaluados**: 62 distintos (tema × primer plano × fondo). En texto, por archivo: 14, 16, 15, 16, 8, 11, 7 y 15. En UI (marcas de gráfica y bordes de control): 0 a 8 por archivo.
- **Bajo umbral: 0** en claro y 0 en oscuro.
- Los más ajustados:
  - Texto: `--ui-text-secondary` sobre `--ui-border-divider` (encabezado de días del mes), 4.51 en claro.
  - UI: `--ui-person-axel` sobre tarjeta, 3.10 en oscuro; `--ui-chart-2` sobre tarjeta, 3.15 en claro; `--ui-chart-1` sobre tarjeta, 3.22 en oscuro.
- La rejilla (`--ui-viz-grid`) y el carril del medidor (`--ui-viz-track`) son decorativos según DESIGN §8.1 y quedan fuera del umbral.
- **Cruce con las tablas oficiales**: los 62 pares tienen fila (por hex) en `contraste.csv` o `contraste-componentes.csv`, y los 62 coinciden con mi cálculo al ±0.01. Ejemplos: secundario sobre tarjeta 5.04/5.03, texto sobre página 11.53/11.53, texto sobre botón primario 5.77/5.77 y secundario sobre página en oscuro 10.43/10.42.
- **Pares extra** recalculados a mano para lo que el motor no ve:
  - `--ui-role-error` sobre tarjeta: 5.40 en claro y 5.82 en oscuro.
  - `--ui-viz-meter` sobre el carril: 4.51 y 5.50.
  - `--ui-store-sams-club` sobre tarjeta: 3.06 en claro.
  - `--ui-store-market` sobre tarjeta: 3.06 en oscuro.
  - `--ui-person-karen` sobre tarjeta: 3.16 en claro.
  - `--ui-icon-selected` sobre la selección: 5.15 y 8.04.
  - Todos pasan.

## P3-08: las 8 decisiones

| Decisión | Resultado |
|---|---|
| Resumen con ≤ 6 tarjetas | ✔ exactamente 6 en m y en e |
| 14 secciones de Ajustes | n/a en Inicio (la barra lateral enlaza Ajustes con `ajustes-p26-s01-cuenta-roles-e`) |
| Barra lateral colapsable | ✔ botón «Colapsar barra lateral» en los 4 de escritorio |
| Español de México | ✔ (ver supuesto 2 sobre «Error» y «Kueski Pay») |
| Dinero `$1,234.56` con cifras tabulares | ✔. `grep '\$[0-9]{4,}\.[0-9]{2}'` = 0; signo menos U+2212 |
| Fechas `10-oct-2026` (y «Vie 30 oct 2026» donde importa el día) | ✔ «17-ene-2027», «19-oct-2026», «Mié 21 oct 2026»; 0 ISO visibles |
| Paleta de C2 (índigo 280°) | ✔ solo tokens; `color-mix` únicamente en el área de T-2, como permite DESIGN §8.1 |
| Notion como única referencia | ✔ página con barra lateral jerárquica, pestañas y tarjetas de borde fino sin sombra; nada contradice DESIGN |

## Líneas de evidencia

- **P3-04**: este revisor hizo 0 llamadas a navegador y 0 a Artifact. El reporte del autor (`artboards-inicio.md`) declara que trabajó sin navegador (K45). No revisé el journal del flujo: esa verificación la hace el orquestador en el reporte de integración.
- **P3-06** (cobertura de Inicio): P-01 tiene 4 subsecciones × 2 anchos = 8 artboards. Los 8 están en disco y en el fragmento (`boards` = `order` = archivos). Faltantes: 0.
- **P3-07**: este informe, con la tabla artboard × regla. Recalculé 62 pares (más de 10) y todos coinciden con las tablas oficiales al ±0.01.
- **P3-08**: 8 de 8 en la lista de comprobación. El grep de montos sin separador da 0. La contradicción de Q2-18 A está en A-1.
- **P3-12**:
  - El grep de las 21 cifras canon da 0, con y sin coma, sobre los 8 artboards y el fragmento. En centavos, 0. Credenciales, 0.
  - Fechas: el comando de P3-12 lista 15 coincidencias inevitables de quincena (`2026-10-30`, `2026-11-15`, …), las mismas que documenta `LEEME.md`. En Inicio aparecen el 30 oct, el 15 nov y el 30 nov 2026, siempre con montos sintéticos: $1,861.58, $954.60, $906.98, $803.27, +$14,780.00. Ninguna va con un monto del reporte ni de Kueski.
  - Cada artboard trae entre 10 y 48 montos distintos, contando el markup y los valores de `renderVals()`. Todos están en `datos-sinteticos.json`, salvo 2 sumas de T-2 declaradas en el supuesto 2 del autor: $242,520.65 y $257,740.56, los acumulados de `proyeccion.meses`. Su diferencia, −$15,219.91, es igual a `diferencia_acumulada_fmt` del último mes.
  - Muestra de 10 trazas:
    - $3,884.95 ← `inicio.tarjetas.flujo_quincenal.committed_cents`
    - $14,100.00 ← `ahorro.metas[1].ahorrado_cents`
    - $468.50 ← `canasta.por_tienda[3].total_cents`
    - $1,402.98 ← `participacion_karen.despensa_devengada_mes_cents`
    - $1,606.53 ← `participacion_karen.compras_particulares[0].total_cents`
    - $7.20 ← `canasta.detalle[4].diferencia_cents`
    - $1,685.50 ← `categorias[4].real_cents`
    - $399.00 ← `presupuesto.gastos_fijos_recurrentes[1].monto_cents`
    - $1,452.91 ← `creditos.lineas[2].proximo_pago.monto_cents`
    - $650.00 ← `ahorro.metas[1].aportacion_necesaria_cents`
- **U5-04**: 62 pares en 2 temas, 0 bajo umbral (detalle en la sección de contraste).
- **U5-06**:
  - Lo que cumple: tokens de rol correctos (ingreso, gasto, ahorro, advertencia, información, éxito y error, siempre con icono y texto); tipografía con las clases `title-1`, `title-2`, `title-3`, `small` y `caption` del sistema; cifras tabulares; iconos Lucide de 20/24 px con `aria-hidden`; espaciado con `--space-*`; radios `--radius-sm/md/full`; `--shadow-0` con borde; barra lateral colapsable; 6 tarjetas; objetivos de 44 px en móvil (resolví `min-height` y `height` con los tokens: 0 controles por debajo); foco visible con `:focus-visible`; series en orden fijo con marcador, trazo, trama y etiqueta en texto.
  - Lo que falla: el trazo de los iconos de 20 px (B-1).
- **U5-07**: 0 fechas ISO visibles, 0 palabras de interfaz en inglés (ver supuesto 2), dinero en formato y fechas en formato. «Sin conexión» no se dibuja en Inicio: el encargo no lo pide para esta área.

## Texto con aspecto de instrucción

Ninguno. Busqué «ignora», «instrucción», «Claude», «agente», «prompt», «system», «olvida», «ejecuta» y «orquestador» en los 8 artboards y no hubo coincidencias. Los archivos del tipo (`format.md`) los leí como dato.

## Correcciones para la vuelta 2 (en orden)

1. **A-1 (orquestador y autor).**
   - El orquestador corrige los estados de meta en `datos-sinteticos.json` a los tres de Q2-18 A: `inicio.tarjetas.ahorro.next_goal.status` → `in_progress` y `etiqueta_estado` → «En curso»; `ahorro.metas[].estado` → `in_progress` en las dos metas, con `etiqueta_estado` «En curso». Si quiere que se vea una meta vencida, le pone una fecha meta pasada.
   - Luego regenera los datos y corre `verificar_sinteticos.py`. Este cambio afecta también a P-09 (área presupuesto).
   - El autor cambia las 4 insignias de Inicio a `bd bi` con el icono `info`.
2. B-1: `.ic{stroke-width:var(--icon-stroke-user-units-sm)}` y `.ic24{stroke-width:var(--icon-stroke-user-units-md)}` en los 8 artboards.
3. B-2: `circle-dashed` delante de «Pendiente» en D-3.
4. B-4: «Meta de Axel» y «Meta de Ximena», en las etiquetas y en la leyenda de T-4b.
5. B-5: el enlace de la tarjeta en los ejemplos de error.
6. B-6: el nombre accesible de «Pendiente confirmar» (texto oculto visualmente) y, si se quiere, un `aria-label` corto en cada tarjeta-enlace (título, valor y estado).
7. B-3: decidir la marca «viejo» en D-5 y corregir el supuesto 7 del reporte del autor.
8. B-7: el orquestador mide los altos al publicar.
