# P3-iter-3 · Créditos · revisión de la vuelta 3

Revisor: `revisor-creditos-r3` (independiente, vuelta 3 de 3; no corrige) · Fecha: 12-oct-2026 · Criterio: P3-07 · Alcance: los 8 artboards de `docs/design/maquetas/lienzo/project/creditos/` (P-17 a P-20, móvil y escritorio) y `docs/design/maquetas/lienzo/fragmentos/creditos.json`.
Fuentes: `docs/criterios/P3-claude-design.md` 1.3 (P3-07, P3-08, P3-12), `docs/criterios/U5.md` (U5-04, U5-06, U5-07), `docs/design/DESIGN.md` (§7.7, §7.8, §7.10, §7.12, §8.5, §9.3, §10.5), `docs/design/U1-arquitectura.md` (§2.1, §2.2, §3.19 a §3.22, §6), `referencia-tipo/format.md` y `datos-sinteticos.json` 1.2. Las vueltas anteriores: `P3-iter-2-creditos.md` y `correccion-creditos-r2.md`.

## Veredicto

**Aprobado: 0 altos, 0 medios y 15 bajos.** Los dos hallazgos que bloqueaban en la vuelta 2 están cerrados:

- **A-01 (alto, contraste): cerrado.** Las barras de «deuda al cierre de cada mes» ya pintan el relleno con `--ui-viz-meter` sobre `--ui-viz-track`, que es el par C-102 de DESIGN §10.5. Lo recalculé: 4.51 en claro y 5.50 en oscuro.
- **M-01 (medio, agenda de P-17 móvil): cerrado.** Los 13 pagos de la agenda traen un botón «Marcar pago» o «Ver pago» de 44 px. El botón abre la hoja «Marcar pago» con el mismo contenido que P-18 m, y los datos cuadran 13 de 13 con `creditos.pagos_programados`.

Mi script pasa sin fallas en forma, tokens, P3-12, enlaces, objetivos táctiles, idioma y contraste. Los 11 bajos de la vuelta 2 siguen como estaban, porque no entraban en la corrección r2. Agrego 4 bajos nuevos (B-12 a B-15). Ninguno sube a medio.

## Cierre de los hallazgos de la vuelta 2

| Id r2 | Hallazgo | Estado | Evidencia (render propio de `renderVals()` en node y plantilla expandida en Python) |
|---|---|---|---|
| A-01 | Relleno `--ui-chart-1` sobre `--ui-viz-track`: 2.22:1 en oscuro | **Cerrado** | Los 4 archivos (`creditos-p17-lineas-m`, `-p17-lineas-e`, `-p18-credito-m`, `-p18-credito-e`) usan `<div style="height:100%;width:{{ m.w }};…;background:var(--ui-viz-meter)">` dentro de la pista `--ui-viz-track`. Claro: #66676e sobre #e5e6e9 = **4.51**. Oscuro: #d1d2d8 sobre #4d4e53 = **5.50**. Es igual a C-102. Hay 0 `background:var(--ui-chart-*)` en los 8 artboards. Un difflib contra la copia de `docs/design/salida/iter-2/creditos/` muestra que en P-17 e, P-18 m y P-18 e solo cambió ese token. La línea A/B de P-18 sigue con `--ui-chart-1` como trazo sobre la tarjeta: 7.23 en claro y 3.22 en oscuro, que es C-106. |
| M-01 | Los pagos de la agenda de P-17 m no se podían tocar | **Cerrado** | En P-17 m › Calendario, cada pago trae un `<button class="btn btn-s">`, que mide 44 px por `--size-control-coarse`. El nombre accesible es «Marcar pago de Kueski Pay, Vie 30 oct 2026». El botón abre una hoja `role="dialog"` con `aria-modal="true"` y `aria-labelledby="pay-t"`. La hoja muestra la compra de origen («Soriana · despensa chica»), el monto programado $954.60, la fecha programada «Vie 30 oct 2026», el capital repuesto $820.00, una casilla por parte con monto y estado («Parte del hogar $572.76 · pendiente», «Parte de Karen $381.84 · pendiente»), la fecha real «21-oct-2026» y los botones Guardar y Cancelar. Al pulsar Guardar, la fila pasa a «pagada», muestra el ícono de check y el botón cambia a «Ver pago». «Ver pago» abre la hoja con las casillas deshabilitadas. Con el filtro «Pagados» y nada marcado aparece «Sin pagos con este filtro.». Las filas de corte no llevan botón. Crucé fecha, línea, monto, partes y capital repuesto contra `pagos_programados` y coinciden **13 de 13**. Quedan residuos bajos: B-01, B-02 y B-12. |
| B-01 a B-11 | Bajos de la vuelta 2 | **Siguen**, porque la corrección r2 no los tenía en su alcance | B-01 y B-02 ahora también aplican a la hoja nueva de P-17 m. El detalle está en «Hallazgos». |

## Líneas de evidencia

- **P3-07**: revisé los 8 artboards contra DESIGN, U1 (§2.1, §2.2, §3.19 a §3.22 y §6) y los criterios estáticos de U5. U3 y U4 no aplican a este área. Medí el contraste en el HTML renderizado con los tokens resueltos desde `tokens.json`, en los dos temas, en reposo y con `:hover`. Resultado: **0 pares de texto bajo 4.5:1 y 0 pares de interfaz bajo 3:1**, sin contar los que DESIGN declara decorativos. Hice 61,186 evaluaciones de texto y 14,612 de interfaz sobre 65 combinaciones distintas de token × fondo × tipo × estado.
- **P3-08**: no se contradice ninguna decisión. El texto está en español de México; el dinero va como `$1,234.56` (salvo los ejes de P-18 e, B-07); las fechas como `21-oct-2026` y «Vie 30 oct 2026», con los días de la semana revisados contra el calendario; la paleta índigo sale de los tokens; y la barra lateral es colapsable («Ocultar barra lateral» en las 4 de escritorio). Lo de Resumen ≤ 6 y las 14 secciones de Ajustes no aplica a este área.
- **P3-12**: el grep de §2 (21 cifras, con y sin coma) da **0** en los 8 artboards, en el fragmento y en `datos-sinteticos.json`. La variante en centavos sin punto también da 0, igual que el grep de credenciales. **Fechas**: el comando de P3-12 da 15 fechas que coinciden con `model.json` y `data.json` (2026-10-30, 2026-11-15, 2026-11-30, 2026-12-15, 2026-12-30 y diez de 2027), todas quincenales por construcción. Ninguna va junto a una cifra del reporte o de Kueski, porque esas cifras no están en los artboards. **Fechas visibles**: hay 78 distintas y las 78 están en el JSON. **Montos**: hay 237 distintos en el markup. 228 están literalmente en `datos-sinteticos.json` (`*_cents`, `*_fmt` o `por_linea_cents`). Los 9 restantes son derivados con aritmética que comprobé: $1,070.65 = $6,189.05 − $5,118.40; $2,152.22 = $15,981.90 − $13,829.68; $1,240.37 = $7,067.09 − $5,826.72; $14,099.04 = 12 × $1,174.92; $8,224.44 = 7 × $1,174.92; $8,833.87 = 9 × $883.39 + $883.36; y $0.01 y $0.02, las diferencias entre quincenas de P-20. **Muestra de 10** rastreada al JSON: $288.95 (`creditos.tienda.abonos[1].interes_cents`), $1,021.10 (`comparador.escenarios[1].intereses_cents`), $906.98 (`pagos_programados`, kc1 30-oct), $1,174.92 (`creditos.tarjeta.msi[0].mensualidad_cents`), $820.00 (`pagos_programados`, kc0 30-oct `capital_repuesto_cents`), $1,164.24 (`creditos.prestamo.amortizacion[6].capital_cents`), $5,205.53 (`creditos.tienda.abonos[2].saldo_cents`), $883.36 (`creditos.tienda.abonos[9].pago_cents`), $42,744.75 (`deuda_viva.serie_mensual[0].total_fmt`) y $708.88 (`creditos.tienda.abonos[5].capital_cents`).
- **U5-04**: el par más bajo de texto es `ui-text-secondary` sobre `ui-surface-hover` en hover (4.51 y 5.50, es C-138). En interfaz, los más bajos son `ui-chart-2` sobre `ui-surface-card` (3.15 en claro, C-107), `ui-chart-1` sobre `ui-surface-card` (3.22 en oscuro, C-106) y `ui-person-karen` sobre `ui-surface-raised` (3.38 en claro, C-052). Los pares bajo 3:1 que quedan son todos decorativos según DESIGN: la pista sobre la tarjeta o la página (C-105), la retícula (C-116), los esqueletos (C-143 y C-144) y el velo de la hoja. Por eso no cuentan. La franja de promoción (`ui-surface-selected` sobre la tarjeta, 1.00) se identifica por su borde superior de 3 px `ui-border-selected`, que da 5.16 y 7.94.
- **U5-06**: los tokens de color y de rol son correctos. Las cifras son tabulares en tablas y cifras sueltas, pero faltan en el texto corrido (B-05). En móvil, todos los controles miden ≥ 44 px (`min-height` resuelto con los tokens). Hay foco visible con `:focus-visible` en `a`, `button`, `summary`, `input` y `select`. La barra lateral es colapsable. Los iconos de 14 y 16 px siguen (B-06), y la barra inferior móvil sigue sin iconos (B-08).
- **U5-07**: no hay texto en inglés. «Error» aparece como etiqueta de estado y es palabra del español. «Kueski Pay» es nombre propio permitido. No hay emoji, ni barra de estado o teclado falsos, ni negativos con guion ante `$`. No hay fecha ISO en el texto visible, salvo el `value` de `<input type="date">` en P-19 (B-09). Los días de la semana que aparecen («Vie 30 oct 2026», «Dom 15 nov 2026», «Mar 17 nov 2026», «Mié 18 nov 2026», «Vie 20 nov 2026») coinciden con el calendario.

## Tabla artboard × regla

F = forma del tipo · T = tokens · C = contraste · D = datos sintéticos · I = idioma y formatos · A = accesibilidad · L = enlaces

| Artboard | Cumple | Falla (Id) | Corrección sugerida (opcional; no bloquea) |
|---|---|---|---|
| `creditos-p17-lineas-m` | F, T, C, D, I, L, A (44 px). Raíz fija de 412 px. Agenda con «Marcar pago» y hoja; filtros; barras con C-102 | B-01, B-02, B-03, B-04, B-05, B-08, B-11, B-12, B-13 | Asa y primaria de 48 px en la hoja; nombres accesibles con la compra |
| `creditos-p17-lineas-e` | F, T, C, D, I, L, A. Página fluida con `expand: fill`; Estados vacío, carga y error; franjas; panel Marcar pago | B-02, B-03, B-04, B-05, B-06, B-12, B-13, B-14 | Enlaces a `ajustes-p26-s05-creditos-condiciones-e` en lugar de `#` |
| `creditos-p18-credito-m` | F, T, C, D, I, L, A. Variantes por tipo, serie de disponible, gráfica a escala 1, hoja Marcar pago | B-01, B-02, B-05, B-06, B-07, B-08, B-10, B-11, B-12, B-13 | Filtros de compras; «Ver como tabla» |
| `creditos-p18-credito-e` | F, T, C, D, I, L, A. Estados con «Falta capturar · Capturar» (K18); compras de la tarjeta y de la tienda | B-02, B-03, B-05, B-06, B-07, B-11, B-12, B-13, B-14, B-15 | Ejes «$5 mil»; vista de mes en Calendario |
| `creditos-p19-compra-m` | F, T, C, D, I, L, A. `label for` en todos los campos | B-06, B-08, B-09, B-13 | Fecha como texto `12-oct-2026` |
| `creditos-p19-compra-e` | F, T, C, D, I, L, A. Estados «Calculando…», error y dato inválido | B-06, B-09, B-13, B-14 | Igual que en móvil |
| `creditos-p20-comparador-m` | F, T, C, D, I, L, A. Compras vivas máx. y diferencia entre quincenas | B-03, B-06, B-08, B-11 | Rendimiento por peso y pesos iguales cuando estén en el JSON |
| `creditos-p20-comparador-e` | F, T, C, D, I, L, A. Estados: ninguno, faltan restricciones, evaluando y no disponible | B-03, B-04, B-06, B-11, B-13, B-14, B-15 | Tabla de escenarios sin desplazamiento a 1440 px |

## Hallazgos

### Altos

Ninguno.

### Medios

Ninguno.

### Bajos (no bloquean; se arrastran a la fase siguiente o a la referencia visual)

- **B-01 · Hoja móvil fuera de la anatomía de DESIGN §7.12** (P-17 m, incluida la hoja «Marcar pago» nueva, y P-18 m). Las hojas cuelgan de la barra superior (`top: var(--layout-topbar-height-mobile)`) y tienen el radio abajo. No tienen el asa de 32 × 4 px (`--size-sheet-handle-*`), y la primaria mide 44 px en lugar de 48 px (`--size-control-primary-sheet`) fija abajo.
- **B-02 · N-210 incompleto en «Marcar pago»** (P-17 m, P-17 e, P-18 m y P-18 e). El monto solo se muestra y no se puede editar. Las partes pagadas salen `checked disabled`, así que no se pueden desmarcar. En «Ver pago», la hoja sigue titulada «Marcar pago», «Guardar» sigue activo y la fecha real propuesta es la de hoy, no la registrada.
- **B-03 · Huecos del JSON** (dueño: `datos-sinteticos.json`). Los «Intereses pagados en el año» dicen «Sin dato» (P-17) y el «Rendimiento por peso» dice «sin dato» (P-20). Los pesos son 40/30/30/0 en lugar de 4 iguales (N-225). A las compras de la tarjeta y de la tienda les faltan columnas de N-202, y en P-18 e la columna «Fecha» dice «Meses sin intereses». Los intereses del mes no separan lo pagado de lo por pagar (N-194). Banorte Oro dice «Compras vivas 2», pero P-18 lista 1.
- **B-04 · `href="#"` aunque el destino existe.** Pasa en «Confirmar en Ajustes › Créditos y condiciones» (P-17 m y e), «Agregar en Ajustes › Créditos» (Estados de P-17 e) y «Faltan restricciones · Capturar» (P-20 e). El destino es `../ajustes/ajustes-p26-s05-creditos-condiciones-{m,e}.dc.html`. No son enlaces rotos, porque `#` no lleva a ningún lado, y no están en la barra lateral ni en un botón principal.
- **B-05 · Dinero sin `tabular-nums` en texto corrido.** Aparece en «Lectura A: $11,440.00 · pendiente de confirmar» (`div.caption.sec`, P-17 m y e), «Préstamo de $18,000.00 · desembolso 18-jun-2026» (`dd`, P-17 e), «compra mínima $2,000.00» (`div.promo`, P-17 e) y en la `figcaption` y los `<text>` de la gráfica A/B (P-18 m y e).
- **B-06 · Iconos de 14 y 16 px.** DESIGN §5.2 pide 20 o 24 px. Pasa en P-17 e y m, P-18 e y m, P-19 e y m y P-20 e y m.
- **B-07 · P-18.** En escritorio, los ejes dicen «$0», «$5,000», «$10,000» y «$15,000» en lugar de «$5 mil» (§8.5), que P-18 m ya usa. No hay «Ver como tabla» (§7.7) ni en m ni en e. La pestaña Calendario de P-18 e es una tabla con «15-oct-2026» sin día de la semana ni vista de mes (§7.8 y U1 §0.4).
- **B-08 · Barra inferior móvil** (los 4 artboards `-m`). No tiene iconos, usa `ui-surface-card` en lugar de `ui-surface-raised` y sus destinos miden 44 px en lugar de 56 (§7.10).
- **B-09 · `<input type="date" value="2026-10-12">` en P-19 m y e.** El navegador muestra la fecha en su formato local, mientras que el resto de las pantallas usa «12-oct-2026».
- **B-10 · P-18 m sin filtros de compras.** N-202 pide «vivas o liquidadas, familia». P-18 e sí los tiene.
- **B-11 · Alturas estimadas y montos derivados.** Los `h` del fragmento (P-17 m 3000, P-18 m 1800, P-18 e 1300, P-17 e, P-20 m y e) no se midieron en un navegador. Además hay 9 montos derivados que no están en el JSON; su aritmética cuadra (ver P3-12).
- **B-12 (nuevo) · Nombres accesibles repetidos en los botones de pago.** En P-17 m, «Marcar pago de Kueski Pay, Vie 30 oct 2026» aparece 2 veces, y lo mismo pasa el 15-nov y el 30-nov: son dos compras distintas el mismo día. En P-17 e (13 botones) y en P-18 m y e (16 botones) dicen solo «Marcar pago». **Corrección:** incluir el número de pago o la compra, por ejemplo «Marcar pago 2 de 2 · Soriana · Kueski Pay, Vie 30 oct 2026».
- **B-13 (nuevo) · Detalles del formato del tipo.** Hay 19 `<sc-for>` sin `hint-placeholder-count` (P-17 m 5, P-17 e 1, P-18 m 8, P-18 e 4, P-20 e 1), aunque `format.md` pide poner siempre los `hint-*`. El `<form onsubmit="return false">` de P-19 m y e es un manejador en texto que el formato no usa; los eventos van como `onSubmit="{{ fn }}"`. Hoy no tiene efecto, porque el formulario no tiene botón de envío y tiene 7 campos.
- **B-14 (nuevo) · Token de esqueleto equivocado.** En la tarjeta «Estados» de P-17 e, P-18 e, P-19 e y P-20 e, `.sk` usa `--ui-skeleton-on-card` sobre un bloque `--ui-surface-page`, pero C-143 pide `--ui-skeleton-on-page` sobre la página. Es decorativo y no cambia el veredicto de contraste.
- **B-15 (nuevo) · Tabla ancha en escritorio.** «Todos los escenarios» de P-20 e tiene `min-width:1560px` (14 columnas) dentro de un contenedor de 1120 px, así que se desplaza horizontalmente aun a 1440 px. U1 §6.1, regla 3, solo admite desplazamiento horizontal en tablas con columnas agregadas por la persona. «Compras de la línea» de P-18 e (`min-width:1100px`) probablemente también se desplaza; hay que medirlo. **Corrección:** mover las columnas secundarias a la vista previa o a una segunda línea por fila.

## Fidelidad: pantalla → elementos de U1

| Pantalla | Presentes | Faltantes o parciales |
|---|---|---|
| P-17 (N-190 a N-194, §2.2) | Tarjeta por línea con nombre, tipo, límite o monto original, deuda viva y disponible. El disponible es la lectura B activa y la A «pendiente de confirmar» en la revolvente, «Límite menos saldo» en la tarjeta y «No aplica» en préstamo y tienda. Cada tarjeta trae próximo pago (fecha y monto) y compras vivas. Totales: deuda viva, quincena y mes con hogar y Karen, e intereses del mes. Aviso A/B fijo («Reposición de capital pendiente de confirmar con Kueski»). Calendario: agenda de 60 días en m y mes en e, con partes, estado por parte, corte, promociones y filtros (línea, familia y estado). En m, cada pago abre «Marcar pago». Tendencias de 12 meses: en e por línea; en m total + 1 dato, como permite U1 §6.2. Primaria «Registrar compra a crédito», más «Simular compra» y «Ver calendario». Estados vacío, carga y error (e) | Intereses del año e intereses pagados/por pagar (B-03) |
| P-18 (N-200 a N-210) | Variante por tipo: revolvente (A/B, serie, mínimos $11,440.00 y $6,881.60, tasas por plazo), tarjeta (ciclo, MSI, pago mínimo), préstamo (amortización) y tienda (abonos). Compras de la línea, calendario con partes, Marcar pago por parte, «Editar condiciones» que abre P-26.05 y «Falta capturar · Capturar» (K18). Estados «Sin compras vivas», «Sin pagos pendientes», esqueleto y Reintentar (e) | Filtros de compras en m (B-10); vista de mes y «Ver como tabla» (B-07); desmarcar y monto editable (B-02); columnas de N-202 (B-03) |
| P-19 (N-215 a N-218) | Línea, comercio, fecha, capital, intereses con IVA, plazo, promoción, partes por persona, vista previa del cálculo, pagos con capital repuesto y estado por parte, gasto vinculado, Editar y Eliminar. Estados «Calculando…», «No se pudo calcular el calendario» y dato inválido (e) | Fecha ISO en el `input` (B-09) |
| P-20 (N-225, N-226) | Datos de la compra, restricciones «solo para esta comparación», topes de Axel y Ximena (no duros), pesos, 4 óptimos, recomendación con puntaje, tabla con intereses, pago máximo, diferencia, disponible mínimo, compras vivas máx. y deuda en cero, descartados con motivo, total evaluado, «Usar escenario recomendado» que abre P-19, «Cambiar restricciones» y «Cambiar pesos». Estados (e) | Rendimiento por peso y pesos iguales (B-03); tabla ancha (B-15) |

## Contraste recalculado (U5-04): hex de `tokens.json`, referencias resueltas por tema

| Par (primer plano / fondo) | Tipo | Claro | Oscuro | Par de §10.5 |
|---|---|---|---|---|
| ui-text-secondary / ui-surface-hover (hover) | texto | 4.51 | 5.50 | C-138 |
| ui-role-success / ui-surface-card | texto e icono | 4.74 | 6.42 | C-121 |
| ui-role-success / ui-role-success-surface | texto | 4.81 | 6.19 | C-072 |
| ui-text-secondary / ui-surface-card | texto | 5.04 | 7.97 | C-095, C-117 |
| ui-role-warning / ui-role-warning-surface | texto | 5.12 | 6.09 | C-070 |
| ui-text-link / ui-role-warning-surface | texto | 5.15 | 8.01 | cumple 4.5 |
| ui-text-link / ui-surface-card | texto | 5.16 | 7.94 | C-030 |
| ui-role-debt / ui-role-debt-surface | texto | 5.29 | 6.05 | C-069 |
| ui-text-on-primary / ui-action-primary-bg (y hover) | texto | 5.77 (8.49) | 10.47 (12.66) | C-001, C-004 |
| **ui-viz-meter / ui-viz-track** (relleno de barra, A-01) | interfaz | **4.51** | **5.50** | **C-102** |
| ui-chart-1 / ui-surface-card (trazo A) | interfaz | 7.23 | 3.22 | C-106 |
| ui-chart-2 / ui-surface-card (trazo B) | interfaz | 3.15 | 7.73 | C-107 |
| ui-viz-baseline / ui-surface-card | interfaz | 3.47 | 6.09 | C-115 |
| ui-person-karen / ui-surface-raised (punto) | interfaz | 3.38 | 4.68 | C-052 |
| ui-border-control / ui-surface-card | interfaz | 3.47 | 6.09 | C-011 |
| ui-border-selected / ui-surface-card (franja de promoción) | interfaz | 5.16 | 7.94 | cumple 3 |
| ui-border-focus / ui-surface-hover (el peor fondo) | interfaz | 3.19 | 4.16 | cumple 3 |
| ui-viz-track / ui-surface-card · ui-viz-grid / ui-surface-card | decorativo | 1.12 | 1.45 | C-105, C-116 |
| ui-viz-track / ui-surface-page (pista en móvil) | decorativo | 1.20 | 1.90 | análogo a C-105 |
| ui-skeleton-on-card / ui-surface-page | decorativo | 1.20 | 1.90 | C-143 (ver B-14) |

Para comprobar mi implementación, comparé con la tabla de la vuelta 2 los 14 pares que tenemos en común y coinciden en los dos temas (±0.01). Sigue siendo cierto que antes de la corrección `ui-chart-1` sobre `ui-viz-track` daba 6.48 y 2.22.

## Verificación (scripts propios en `datos/tmp/rev-cred3/`, `python3 -I` y node; borrados al final)

1. **Forma**, 8 de 8: tienen `<!doctype html>`, `lang="es-MX"`, `meta charset`, un `<title>` propio y único, `<script src="./support.js">` y `<link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css">` exactos, un solo `x-dc` con `helmet><style>` al inicio, un `data-dc-script` con JSON válido, `class Component extends DCLogic` y `renderVals()`. El `$preview` coincide con `w`×`h` del fragmento. `html.parser` no encontró etiquetas sin cerrar ni cierres huérfanos, y los ids son únicos. Los nombres son válidos, y hay 0 stems duplicados entre los 90 artboards del lienzo. Todos los huecos son rutas con puntos (0 expresiones) y todos los `on*` son `{{ fn }}`, salvo `onsubmit` (B-13). La raíz móvil es `width:412px;min-height:915px;box-sizing:border-box`. La de escritorio no tiene ancho en px, usa `flex-wrap`, contenido `flex:999 1 560px;min-width:0` y `max-width:1120px`; ninguna usa `sticky` ni `100vh`, y todas tienen `expand: fill` en el fragmento. Las tablas anchas van en `overflow-x:auto`. El fragmento es JSON válido y `boards` y `order` listan exactamente los 8 archivos.
2. **Ejecución**: corrí `renderVals()` en node con un `DCLogic` mínimo y exploré estados por BFS llamando a todos los manejadores: 150 estados en P-17 y P-18, 2 en P-19 y 1 en P-20, con 0 errores. Expandí la plantilla (`sc-if`, `sc-for` y huecos) en Python en 40 estados por archivo, más secuencias dirigidas en P-17 m: Calendario → Marcar pago → Guardar → Ver pago, y Filtrar → Pagados. Quedaron 0 huecos sin resolver. Solo 4 nodos de plantilla no se cubrieron: 3 variantes de la casilla de la hoja y el «Sin pagos con este filtro.» de P-18 m, que usa el mismo `p.small.sec` que el de P-17 m, que sí se renderizó.
3. **Tokens**: 0 hex, 0 `rgb()` y 0 `hsl()` fuera de `data-props`. Los SVG solo usan `currentColor` y `none`. Las `var(--x)` distintas de cada archivo (59, 50, 61, 53, 52, 45, 50 y 38 en P-17 e, P-17 m, P-18 e, P-18 m, P-19 e, P-19 m, P-20 e y P-20 m) existen en `tokens.json`, que tiene los mismos nombres que `tokens.css` (0 diferencias).
4. **Contraste**: escribí una cascada CSS propia (selectores, especificidad, herencia y `:hover` de cada elemento con regla hover). El fondo efectivo es el del ancestro más cercano que tenga fondo. Los umbrales son 4.5 para texto, 3 para texto de 24 px o más e interfaz (iconos, trazos, rellenos sin texto y bordes de controles).
5. **Accesibilidad**: todo `button`, `a` y `summary` tiene nombre. Todo `input` y `select` tiene `label` (`for`, `label` envolvente o `aria-label`). No hay interactivos anidados en ningún estado. Los SVG son `aria-hidden` o tienen nombre. En móvil, todo control mide ≥ 44 px; las únicas excepciones son las tarjetas-enlace de las líneas, que miden mucho más.
6. **P3-12**: el grep oficial y la variante en centavos dan 0. Corrí el comando de fechas de P3-12, y rastreé los montos y las fechas como se describe arriba.
7. **Enlaces**: cada `href` relativo existe en el lienzo y tiene el mismo ancho (`-m` → `-m` y `-e` → `-e`). Los `#` son Buscar, Editar y Reintentar, que no tienen artboard, y los 3 de B-04.
