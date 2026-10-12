# P3-iter-2 · Créditos · revisión de la vuelta 2

Revisor: `revisor-creditos-r2` (independiente; no corrige) · Fecha: 12-oct-2026 · Criterio: P3-07 · Alcance: los 8 artboards de `docs/design/maquetas/lienzo/project/creditos/` (P-17 a P-20, móvil y escritorio) y `docs/design/maquetas/lienzo/fragmentos/creditos.json`.
Fuentes: `docs/criterios/P3-claude-design.md` 1.3 (P3-07, P3-08, P3-12), `docs/criterios/U5.md` (U5-04, U5-06, U5-07), `docs/design/DESIGN.md` (§7.7, §7.8, §7.9, §7.12, §8.5, §9.3, §10.5), `docs/design/U1-arquitectura.md` (§0.4, §2.1, §2.2, §3.19 a §3.22, F-4), `referencia-tipo/format.md` y `craft.md`, y `datos-sinteticos.json`.

## Veredicto

**No aprobado: 1 alto, 1 medio y 11 bajos.** Los 9 hallazgos medios de la vuelta 1 están cerrados en los archivos que se señalaron. Quedan dos residuos de datos (M-06 y M-07), que bajan a bajo porque los bloquea `datos-sinteticos.json`, fuera del área. Hay 2 hallazgos nuevos que impiden aprobar:

- **A-01 (alto, contraste).** En las barras de «deuda al cierre de cada mes», el relleno `--ui-chart-1` sobre la pista `--ui-viz-track` da **2.22:1 en oscuro**, por debajo de 3:1. Ese par no está en DESIGN §10.5, que para la barra de avance aprueba `ui.viz.meter` sobre `ui.viz.track` (C-102). El defecto viene de la vuelta 1, pero esa revisión no lo reportó.
- **M-01 (medio).** La agenda de P-17 móvil no deja tocar ningún pago (DESIGN §7.8). P-17 escritorio y P-18 sí lo permiten.

Forma, tokens, P3-12, enlaces, objetivos táctiles e idioma pasan el script sin fallas.

## Cierre de los hallazgos de la vuelta 1

| Id r1 | Hallazgo | Estado | Evidencia (render de `renderVals()` en node, varios estados) |
|---|---|---|---|
| M-01 | Partes del hogar y de Karen con monto, y «Filtrar» en la agenda de P-17 m y en el calendario de P-18 m | **Cerrado** | P-17 m › Calendario: «Vie 30 oct 2026 · Kueski Pay · Pago 2 de 2 · Hogar $572.76 · pendiente / Karen $381.84 · pendiente · $954.60». Los pagos solo de Karen muestran solo esa parte. El botón «Filtrar» abre una hoja con Línea, Familia y Estado que sí filtra, y si no hay coincidencias dice «Sin pagos con este filtro». P-18 m › Calendario muestra lo mismo con `r.hogar`/`r.karen` y su propio «Filtrar». |
| M-02 | Serie de disponible en P-18 m › Detalle | **Cerrado** | Muestra 4 eventos, cada uno con fecha, tipo, monto, descripción, capital repuesto, interés no repuesto y disponible A y B, más el botón «Mostrar 16 más». También aparece «Mínimo A $11,440.00 · mínimo B $6,881.60». |
| M-03 | Gráfica A/B ilegible en móvil | **Cerrado** | `viewBox="0 0 348 220"` con `width="100%"` dentro de una tarjeta de 348 px: escala 1 y texto de 12 px. Ejes compactos («$5 mil», «$10 mil», «$15 mil»). Los valores finales van en una línea HTML. |
| M-04 | «Marcar pago» por parte y vista previa (G-06) en P-18 m, P-18 e y P-17 e | **Cerrado** en los 3 archivos | Hoja o panel con `role="dialog"`, `aria-modal` y `aria-labelledby`. Contiene compra de origen, monto y fecha programados, capital repuesto, una casilla por parte con su monto y estado, fecha real «21-oct-2026», Guardar y Cancelar. Las filas de los calendarios traen «Marcar pago» o «Ver pago». En la cuadrícula de P-17 e los 6 pagos de octubre son botones. **Queda un residuo en P-17 m, que no estaba en la lista de r1: ver M-01 de esta vuelta.** |
| M-05 | Promociones como franjas en el mes de P-17 e | **Cerrado** | Bodega Aurrera se dibuja como franja del 19 al 31-oct y el 1-nov, con rótulo «4 quincenas · compra mínima $2,000.00» y «hasta el 18-nov-2026». Hay leyenda «Promoción vigente». |
| M-06 | Comparador incompleto | **Cerrado** en el markup; queda un residuo de datos (B-03) | «Compras vivas máx.» y «Diferencia entre quincenas» aparecen en las tarjetas, la tabla, la recomendación y la lista. Verifiqué E2, cuya diferencia es $0.02 (1,305.28 − 1,305.26 del JSON). «Rendimiento por peso» dice «sin dato» porque el JSON no lo trae. |
| M-07 | Elementos de U1 que el JSON no trae | **Parcial**, pasa a bajo (B-03) | «Intereses pagados en el año» ya está dibujado con «Sin dato». «Compras de la línea» ya aparece en la tarjeta (MSI Liverpool) y en la tienda (Colchón). Falta el resto de las columnas de N-202 en esas dos líneas. Bloquea el dueño de `datos-sinteticos.json`; P3-12 prohíbe inventar las cifras. |
| M-08 | Barra lateral reducida en los 4 artboards de escritorio | **Cerrado** | Tiene cuenta con avatar «A», «Ocultar barra lateral», «Buscar Ctrl+K», «Registrar gasto Alt+G», iconos y chevrons («Desplegar …») y Créditos desplegado con sus hijos. Abajo trae el interruptor «Tema: Sistema / Claro / Oscuro». No es `sticky` ni usa `100vh`. |
| M-09 | «Capturar» de unos 20 px en P-18 m | **Cerrado** | Es un enlace propio con clase `tap` (`min-height: var(--size-touch-min)` = 44 px) que abre `ajustes-p26-s05-creditos-condiciones-m`. |

## Líneas de evidencia

- **P3-07**: revisé los 8 artboards contra DESIGN, U1 y los criterios estáticos de U5. Medí el contraste en los dos temas con los hex de `tokens.css`, que se genera de `tokens.json`. Hay 1 par de interfaz bajo umbral (A-01) y 0 pares de texto bajo umbral. Las tablas están abajo.
- **P3-08**: se cumplen el español de México, el dinero `$1,234.56` (salvo B-05 y B-07, que son bajos), las fechas `21-oct-2026` y «Vie 30 oct 2026» (salvo la tabla de P-18 e, B-07), la paleta índigo por tokens y la barra lateral colapsable en las 4 de escritorio. No se contradice ninguna decisión. Lo de Resumen ≤ 6 y las 14 secciones de Ajustes no aplica a este área.
- **P3-12**: el grep de las 21 cifras prohibidas, con y sin coma, da **0** en los 8 artboards y en el fragmento. El grep de credenciales también da 0. Fechas: 15 fechas sintéticas coinciden con el calendario real porque las dos son quincenales. 12 de ellas aparecen en Créditos, y ninguna va junto a un monto del reporte o de Kueski, porque esos montos no están. Montos: 223 distintos en los artboards. 195 están tal cual en el JSON (`_fmt` o `_cents`). Los 28 restantes salen de sumar o restar cifras del JSON: tendencias por línea de P-17 e, saldo Kueski de oct = $6,189.05 − $1,861.58 = $4,327.47, total de los MSI = 12 × $1,174.92 = $14,099.04 y total de la tienda $8,833.87. Muestra de 10 valores rastreados al JSON: $3,362.55 (`pagos_quincena_actual_hogar`), $6,881.60 (`disponible_b_hoy`), $1,145.53, $729.60, $5,213.40, $1,253.45, $1,383.13, $259.66, $2,800.61 y $10,421.18.
- **U5-04**: 0 pares de texto bajo 4.5:1 en claro y en oscuro (22 combinaciones de texto en reposo y 8 en hover). El par de interfaz `ui-chart-1` sobre `ui-viz-track` da 2.22 en oscuro (A-01). La pista sobre la tarjeta (1.12 y 1.45) y la retícula `ui-viz-grid` son decorativas según C-105 y §8.
- **U5-06**: tokens de color y de rol correctos. Cifras tabulares en tablas y en cifras sueltas; en texto corrido faltan (B-05). Objetivos ≥ 44 px en todos los controles móviles. Foco visible con `:focus-visible`. Barra lateral colapsable. Iconos de 14 y 16 px (B-06). La barra inferior móvil sigue sin iconos (B-08).
- **U5-07**: sin texto en inglés. «Kueski Pay» es nombre propio permitido. No hay emoji, barra de estado falsa ni teclado falso. No hay fecha ISO visible, salvo el `<input type="date">` de P-19 (B-09). Ningún negativo usa guion corto ante `$`. Todos los días de la semana que aparecen («Vie 30 oct 2026») coinciden con el calendario.

## Tabla artboard × regla

F = forma del tipo · T = tokens · C = contraste · D = datos sintéticos · I = idioma y formatos · A = accesibilidad · L = enlaces

| Artboard | Cumple | Falla (Id) | Corrección para la vuelta 3 |
|---|---|---|---|
| `creditos-p17-lineas-m` | F, T, D, I, L, A (44 px); raíz de 412 px fija; partes y «Filtrar» | **A-01**, **M-01**, B-01, B-03, B-04, B-05, B-08 | Relleno de la barra con `--ui-viz-meter`; cada pago de la agenda abre la hoja «Marcar pago» |
| `creditos-p17-lineas-e` | F, T, D, I, L, A; página fluida con `expand: fill`; Estados vacío, carga y error; franjas; panel Marcar pago | **A-01**, B-02, B-03, B-04, B-05, B-06 | Relleno con `--ui-viz-meter`; enlaces de Ajustes en lugar de `#` |
| `creditos-p18-credito-m` | F, T, D, I, L, A; serie de disponible; gráfica a escala 1; hoja Marcar pago | **A-01**, B-01, B-02, B-05, B-07, B-10 | Relleno con `--ui-viz-meter`; hoja inferior con asa; filtros de compras |
| `creditos-p18-credito-e` | F, T, D, I, L, A; Estados con «Falta capturar · Capturar» (K18); compras de la tarjeta y de la tienda | **A-01**, B-02, B-03, B-05, B-06, B-07 | Relleno con `--ui-viz-meter`; ejes «$5 mil»; vista de mes en el calendario |
| `creditos-p19-compra-m` | F, T, C, D, I, L, A; `label for` en todos los campos | B-09 | Fecha como texto `12-oct-2026` |
| `creditos-p19-compra-e` | F, T, C, D, I, L, A; Estados «Calculando…», error y dato inválido; barra lateral canónica | B-06, B-09 | Igual que en móvil |
| `creditos-p20-comparador-m` | F, T, C, D, I, L, A; compras vivas máx. y diferencia entre quincenas | B-03 | Rendimiento por peso y pesos 25/25/25/25 cuando estén en el JSON |
| `creditos-p20-comparador-e` | F, T, C, D, I, L, A; Estados (ninguno, faltan restricciones, evaluando, no disponible) | B-03, B-04, B-06 | «Faltan restricciones · Capturar» → `ajustes-p26-s05-creditos-condiciones-e` |

## Hallazgos

### Alto

- **A-01 · Contraste de interfaz bajo umbral (U5-04, P3-07, DESIGN §10.5 C-102).** Las barras de «deuda al cierre de cada mes» (pestaña Tendencias) pintan el relleno `background:var(--ui-chart-1)` sobre la pista `background:var(--ui-viz-track)`. En oscuro, #797dcd sobre #4d4e53 da **2.22:1**, bajo 3:1; en claro da 6.48. El par no está en §10.5, que para relleno sobre pista aprueba `ui.viz.meter` sobre `ui.viz.track` (C-102: 4.51 en claro y 5.50 en oscuro). La barra es `aria-hidden` y el monto va en texto al lado, pero U5-04 solo exime los elementos que DESIGN declara decorativos, y la pista lo es (C-105) pero el relleno no. Está en `creditos-p17-lineas-m`, `creditos-p17-lineas-e`, `creditos-p18-credito-m` y `creditos-p18-credito-e` (el mismo `<div style="height:100%;width:{{ m.w }};…;background:var(--ui-chart-1)">`). Ya estaba en la vuelta 1 y no se reportó. **Corrección:** cambiar el relleno a `var(--ui-viz-meter)`, que es un cambio de un token en 4 archivos. Otra opción es quitar la pista, porque `ui-chart-1` sobre `ui-surface-card` da 7.23 y 3.22.

### Medio

- **M-01 · Los pagos de la agenda de P-17 móvil no se pueden tocar (DESIGN §7.8; U1 N-192 y N-210).** En la pestaña Calendario de `creditos-p17-lineas-m`, cada pago es un `<li class="card">` sin botón ni enlace. No abre la vista previa (G-06) ni ofrece «Marcar pago», aunque §7.8 dice «Tocar un pago abre su vista previa con «Marcar pago» como primaria» para P-17 y P-18. P-17 escritorio, P-18 móvil y P-18 escritorio sí lo tienen. **Corrección:** agregar a cada fila un botón «Marcar pago» o «Ver pago» de 44 px, o hacer toda la fila un botón con nombre completo. Debe abrir la misma hoja de P-18 m (compra de origen, capital repuesto, partes, fecha real, Guardar y Cancelar) usando `pagos_programados`.

### Bajos

- **B-01 · Hoja móvil fuera de la anatomía de DESIGN §7.12.** Las hojas «Marcar pago» y «Filtrar» de P-18 m y P-17 m cuelgan de la barra superior (`top: var(--layout-topbar-height-mobile)`) y tienen el radio abajo. Les falta el asa de 32 × 4 px, y la primaria mide 44 px en lugar de 48 px fija abajo (`--size-control-primary-sheet`). La hoja inferior pide radio solo arriba.
- **B-02 · Detalles de N-210 en la hoja «Marcar pago»** (P-17 e, P-18 m y P-18 e). El monto solo se muestra («Monto programado») y no se puede editar, aunque N-210 dice «monto (el programado por defecto)». Las partes pagadas salen `checked disabled`, así que no se pueden desmarcar («Permite desmarcar»). En un pago ya hecho, como el del 15-oct, la fecha real propuesta es la de hoy (21-oct-2026) y no la registrada.
- **B-03 · Huecos del JSON** (residuo de M-06 y M-07 de la vuelta 1; dueño: `datos-sinteticos.json`). Faltan «Intereses pagados en el año» (dice «Sin dato» en P-17 m y e) y «Rendimiento por peso» (dice «sin dato» en P-20 m y e). Los pesos son 40/30/30/0 en lugar de 4 objetivos iguales (N-225, era B-08). Banorte Oro dice «Compras vivas 2» en P-17, pero P-18 lista solo 1 (los MSI). A las compras de la tarjeta y de la tienda les faltan capital, intereses, promoción y partes (N-202). En la tabla de compras de la tarjeta de P-18 e, la columna «Fecha» dice «Meses sin intereses». **Corrección:** el dueño de los datos agrega `creditos.totales.intereses_anio_*`, `creditos.tarjeta.compras[]` completas, el rendimiento por peso de cada escenario y pesos 25/25/25/25; luego se sustituyen los «Sin dato».
- **B-04 · Enlaces a `#` con destino existente** (residuo de B-05 de la vuelta 1): «Confirmar en Ajustes › Créditos y condiciones» en P-17 m y P-17 e, «Agregar en Ajustes › Créditos» en Estados de P-17 e y «Faltan restricciones · Capturar» en P-20 e. `ajustes/ajustes-p26-s05-creditos-condiciones-{m,e}.dc.html` existe, y P-18 ya enlaza ahí.
- **B-05 · Dinero sin `tabular-nums` en texto corrido.** Aparece en «Lectura A: $11,440.00 · pendiente de confirmar» (`caption sec`, P-17 m y e), en «Monto original $18,000.00» (`dd`, P-17 e), en la franja «compra mínima $2,000.00» (P-17 e), en la `figcaption` de la gráfica A/B (P-18 m y e) y en las etiquetas `<text>` del SVG. **Corrección:** agregar la clase `tn` o `font-variant-numeric` al `<g>`.
- **B-06 · Iconos de 14 y 16 px** (era B-01; DESIGN §5.2 pide 20 o 24). En el markup hay 14 px en P-17 e ×13, P-17 m ×1, P-18 e ×2, P-18 m ×3, P-19 m ×1, P-20 e ×1 y P-20 m ×3, y 16 px en P-17 e ×3, P-18 e ×2, P-19 e ×3 y P-20 e ×2.
- **B-07 · Calendario y gráfica de P-18 escritorio.** Los ejes dicen «$0», «$5,000», «$10,000» y «$15,000» en lugar de «$5 mil» (§8.5, era B-02). No hay «Ver como tabla» (§7.7) en P-18 m ni en P-18 e. La pestaña Calendario es una tabla con fechas «15-oct-2026» sin día de la semana; §7.8 pide vista de mes desde 768 px y U1 §0.4 pide «Jue 15 oct 2026». Los datos y las acciones sí están completos.
- **B-08 · Barra inferior móvil** (era B-04): no tiene iconos de destino, usa `ui-surface-card` en lugar de `ui-surface-raised` y sus destinos miden 44 en lugar de 56 px (§7.10). Pasa en los 4 artboards `-m`.
- **B-09 · `<input type="date" value="2026-10-12">` en P-19 m y e** (era B-07). El navegador lo muestra en formato local; P-20 y la hoja ya usan texto «12-oct-2026».
- **B-10 · Faltan filtros de compras en P-18 m.** N-202 pide filtros «vivas o liquidadas, familia». P-18 e los tiene («Vivas», «Liquidadas», «Familia: todas») y P-18 m no.
- **B-11 · Alturas y cifras derivadas.** Los `$preview` y `h` del fragmento son estimados: P-18 m de 1800 con la serie desplegada y la hoja abierta puede quedar corto, P-18 e de 1300 parece bajo para la barra lateral desplegada, y P-17 e, P-20 m y P-20 e también son estimados. El orquestador debe medirlos. Además, 28 montos son derivados a mano y no están en el JSON (tendencias por línea de P-17 e, totales de MSI y de tienda). Revisé que cuadran, pero conviene subirlos a `datos-sinteticos.json` para que Reportes e Inicio usen los mismos.

## Fidelidad: pantalla → elementos de U1

| Pantalla | Presentes | Faltantes o parciales |
|---|---|---|
| P-17 (N-190 a N-194, §2.2) | Tarjeta por línea con tipo, límite o monto original, deuda viva, disponible (B activa y A «pendiente de confirmar», límite − saldo, «No aplica») y próximo pago. Totales de la quincena y del mes con hogar y Karen, intereses del mes y aviso A/B fijo. Calendario: agenda de 60 días (m) y mes (e) con partes, corte y promociones. Filtros de línea, familia y estado. Tendencias de 12 meses. Primaria «Registrar compra a crédito»; «Simular compra» y «Ver calendario». Estados vacío, carga y error (e). | Intereses del año = «Sin dato» (B-03). Pago de la agenda móvil sin vista previa (M-01). |
| P-18 (N-200 a N-210) | Variante por tipo: revolvente (A/B, serie, tasas por plazo), tarjeta (estado del ciclo, MSI, pago mínimo con fórmula), préstamo (amortización) y tienda (abonos). Compras de la línea. Calendario con partes. Marcar pago por parte. Editar condiciones abre P-26.05. Estado «Falta capturar · Capturar» (K18). | Columnas de N-202 en tarjeta y tienda (B-03). Filtros de compras en móvil (B-10). Vista de mes en escritorio (B-07). Monto editable y desmarcar (B-02). |
| P-19 (N-215 a N-218) | Línea, comercio, fecha, capital, intereses (IVA incluido), plazo, promoción, capital por persona y vista previa del cálculo con partes y pago por quincena. Pagos con capital repuesto y estado por parte. Gasto vinculado. Modo detalle o alta. Estados «Calculando…», error y dato inválido (e). | Fecha en formato ISO (B-09). |
| P-20 (N-225, N-226) | Datos de la compra, restricciones duras «solo para esta comparación», topes de Axel y Ximena (no son duros), pesos, 4 óptimos lado a lado, recomendación con puntaje, tabla o lista con intereses, pago máximo, diferencia, disponible mínimo, compras vivas máx., deuda en cero y descartados con motivo. Total evaluado. Estados (e). | Rendimiento por peso = «sin dato». Pesos 40/30/30/0 (B-03). |

## Contraste recalculado (U5-04): hex de `tokens.css`, generado de `tokens.json`

Pares de texto, que salen del DOM renderizado en hasta 60 estados por artboard: 0 bajo 4.5. Los más justos:

| Par (primer plano / fondo) | Tipo | Claro | Oscuro |
|---|---|---|---|
| ui-role-success / ui-surface-card | texto | 4.74 | 6.42 |
| ui-role-success / ui-role-success-surface | texto | 4.81 | 6.19 |
| ui-text-secondary / ui-surface-card | texto | 5.04 | 7.97 |
| ui-role-warning / ui-role-warning-surface | texto | 5.12 | 6.09 |
| ui-text-secondary / ui-role-success-surface | texto | 5.12 | 7.69 |
| ui-text-link / ui-role-warning-surface | texto | 5.15 | 8.01 |
| ui-text-link / ui-surface-card | texto | 5.16 | 7.94 |
| ui-role-debt / ui-role-debt-surface | texto | 5.29 | 6.05 |
| ui-text-on-primary / ui-action-primary-bg (y -hover) | texto | 5.77 (8.49) | 10.47 (12.66) |
| ui-border-control / ui-surface-card | interfaz | 3.47 | 6.09 |
| ui-chart-2 / ui-surface-card (trazo B) | interfaz | 3.15 | 7.73 |
| ui-chart-1 / ui-surface-card (trazo A) | interfaz | 7.23 | 3.22 |
| ui-viz-baseline / ui-surface-card | interfaz | 3.47 | 6.09 |
| **ui-chart-1 / ui-viz-track (relleno de barra)** | **interfaz** | **6.48** | **2.22 (A-01)** |
| ui-viz-track / ui-surface-card · ui-viz-grid / ui-surface-card | decorativo (C-105, §8) | 1.12 | 1.45 |

El hover se calculó con la cascada completa. El aparente `a:hover` = `ui-text-link-hover` sobre el primario no ocurre, porque `.btn-p:hover` gana por especificidad.

## Verificación (scripts propios en `datos/tmp/rev-cred2/`, `python3 -I` y node; borrados al final)

1. **Forma**, 8 de 8: `<!doctype html>`, `lang="es-MX"`, un `<title>` propio y único, `<script src="./support.js">` y `<link … ../ds/finanzas-hogar/tokens.css>` exactos, un solo `x-dc` con `helmet` y `style` al inicio, un `data-dc-script` con JSON válido y `$preview` igual a `w`×`h` del fragmento, `class Component extends DCLogic` y `renderVals()`. Nombres válidos y stems únicos. Todos los huecos son rutas con puntos (0 expresiones). `html.parser` no encontró etiquetas sin cerrar ni cierres huérfanos. Raíz móvil `width:412px` con `border-box`. Escritorio sin ancho en px, con `max-width`, `flex-wrap` y contenido `flex:999 1 560px`, y `expand: fill` en el fragmento. El fragmento lista exactamente los 8 archivos en `boards` y en `order`.
2. **Ejecución**: `renderVals()` corrió en node con un `DCLogic` mínimo. Exploré estados por BFS llamando a todos los manejadores: 312 estados en P-17 m, más de 400 en P-17 e y P-18 (se cortó en el tope), y 1 a 2 en P-19 y P-20. Rendericé la plantilla (`sc-if`, `sc-for`) en hasta 60 estados por archivo: 0 errores y 0 huecos sin resolver.
3. **Tokens**: 0 hex y 0 `rgb()`/`hsl()` fuera de `data-props`. Todas las `var(--x)` existen; `--font-sans` sale de `type.families.sans` de `tokens.json`.
4. **Contraste**: cascada CSS propia (selectores, especificidad, herencia y `:hover`) sobre el DOM renderizado. Fondo efectivo = ancestro más cercano con fondo. Umbral 4.5 para texto, 3 para texto grande e interfaz.
5. **Accesibilidad**: todo `button`, `a` y `summary` tiene nombre. Todo `input`/`select` tiene `label`. No hay interactivos anidados. Los SVG son `aria-hidden` o `role="img"` con nombre. En móvil, todo control mide ≥ 44 px (`min-height` resuelto con los tokens).
6. **P3-12**: el grep de §2 da 0. La intersección de fechas con `model.json` y `data.json` da 15 fechas (12 en Créditos), ninguna junto a un monto prohibido. Montos: 195 de 223 están literalmente en el JSON y los 28 restantes son derivados.
7. **Enlaces**: todo `href` relativo existe en el lienzo y tiene el mismo ancho. Los `#` están listados en B-04.
