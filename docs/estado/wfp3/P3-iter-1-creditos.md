# P3-iter-1 · Créditos · revisión de la vuelta 1

Agente: `revisor-creditos-r1` (revisor independiente, vuelta 1 de 3; no corrige) · Fecha: 12-oct-2026 · Alcance: los 8 artboards de `docs/design/maquetas/lienzo/project/creditos/` (P-17, P-18, P-19 y P-20, móvil y escritorio) y su fragmento `docs/design/maquetas/lienzo/fragmentos/creditos.json`.

Contra: `docs/criterios/P3-claude-design.md` 1.3 (P3-07, P3-08, P3-12), `docs/criterios/U5.md` (U5-04, U5-06 y U5-07 en lo que se ve en el markup), `docs/design/DESIGN.md` (§5.2, §7.6 a §7.10, §8.5, §9.3, §10.5), `docs/design/U1-arquitectura.md` (§0.4, §1.3, §2.1, §2.2, §2.4, §3.19 a §3.22, §6.4, F-4) y `docs/design/maquetas/referencia-tipo/format.md`. U3 (Inicio) y U4 (Ajustes) no aplican a esta área.

## Veredicto

**No aprobado.** 0 altos · 9 medios · 10 bajos. La forma, los tokens, el contraste y las cifras sintéticas están bien; faltan elementos secundarios de U1 (sobre todo en móvil), el flujo «Marcar pago» por parte (F-4) y la barra lateral colapsable.

## Texto con aspecto de instrucción (P3-04)

- En los 8 artboards y el fragmento: ninguno (solo texto de interfaz).
- `format.md` y `craft.md` (archivos del publicador del tipo) están escritos en imperativo; se leyeron como datos de forma, no como órdenes.
- El reporte del autor (`artboards-creditos.md`) pide «el orquestador debe medir las alturas»: es un pendiente, no una orden para este revisor.

## Líneas de evidencia

- **P3-04**: este revisor hizo 0 llamadas de navegador y 0 a Artifact. Solo usó Bash y lectura de archivos. 0 coincidencias de credenciales en `creditos/` (grep del patrón de P3-12).
- **P3-06**: 8 de 8 artboards (P-17, P-18, P-19 y P-20, cada una con `-m` y `-e`). Faltantes frente a U1 §2.1: ninguno. El fragmento lista exactamente los 8 archivos en `boards` y en `order`.
- **P3-07**: tabla artboard × regla abajo. Contraste recalculado con los hex de `tokens.json` en claro y oscuro: 26 pares extraídos del markup y 16 pares extra de interfaz o estados; **0 bajo umbral**.
- **P3-08**: lista de 8 puntos abajo. Falla «barra lateral colapsable». `grep -noE '\$[0-9]{4,}\.[0-9]{2}'` = 0.
- **P3-12**: grep de las 21 cifras canon (con y sin coma) = **0** en los 8 artboards y el fragmento. Credenciales = 0. Las fechas coinciden con `oraculo/model.json` en 15 fechas quincenales (30-oct-2026 … 15-jun-2027). Todas van con montos sintéticos ($954.60, $906.98, $803.27, $932.34…) y ninguna con $973.75, $598.82 o $374.93. Muestra de 10 valores contra `datos-sinteticos.json`, todos idénticos: deuda viva total $44,606.33 (`creditos.deuda_viva_total_cents` 4460633), pagos de la quincena $4,197.88 y hogar $3,362.55, pagos del mes $10,365.89, intereses del mes $1,107.07, Kueski deuda $6,189.05, disponible B $6,881.60 y A $11,440.00, próximo pago $1,861.58, Banorte $15,368.29 y disponible $30,856.15, Préstamo $15,981.90, Coppel $7,067.09 y abono $883.39, tasas 16.4150 % y 24.3120 %. Los montos que no están literales en el JSON (21 a 26 por archivo en P-17 e y P-18) son derivados aritméticos comprobados: deuda por línea al cierre de cada mes (p. ej. $15,981.90 − $1,452.91 = $14,528.99), intereses por pagar = deuda − capital ($1,070.65) y MSI por pagar = 7 × $1,174.92 = $8,224.44.
- **U5-04**: 0 pares bajo umbral (detalle abajo).
- **U5-06**: tokens de color y rol correctos y cifras tabulares en todo el dinero. Fallan: iconos de 14 y 16 px (B-01), barra lateral sin colapsar (M-08), barra inferior (B-04) y un objetivo táctil (M-09).
- **U5-07**: sin texto en inglés, sin fecha ISO visible salvo el valor del `<input type="date">` de P-19 (B-07) y sin emoji. Dinero `$1,234.56`, salvo los ejes de la gráfica (B-02). El signo U+2212 no aparece porque no hay negativos, y el guion corto ante `$` da 0.

## P3-08 · lista de 8 puntos (área créditos)

| # | Punto | Resultado |
|---|---|---|
| 1 | Resumen ≤ 6 tarjetas | No aplica (P-17 tiene 4 tarjetas de métrica) |
| 2 | 14 secciones de Ajustes | No aplica |
| 3 | Barra lateral colapsable | **Falla**: la barra de escritorio no tiene «‹‹» ni chevrons (M-08) |
| 4 | Español de México | Cumple |
| 5 | Dinero `$1,234.56` tabular | Cumple en cifras y tablas; ejes `$5,000` en lugar de `$5 mil` (B-02) |
| 6 | Fechas `10-oct-2026` | Cumple; agenda «Vie 30 oct» sin año (B-09); input date (B-07) |
| 7 | Paleta C2 (índigo 280°) | Cumple: 0 hex, 0 rgb/hsl y 35 tokens `ui-*`, todos existen |
| 8 | Notion como única referencia | Cumple |

## Tabla artboard × regla

Reglas: F = forma del tipo (format.md y la orden), T = tokens y 0 hex, C = contraste (U5-04), D = datos sintéticos (P3-12), U1 = fidelidad, A = accesibilidad (U5-06 y §9.3), I = idioma y formatos (U5-07), L = enlaces y fragmento.

| Artboard | Cumple | Falla (ID) | Corrección para la vuelta 2 |
|---|---|---|---|
| `creditos-p17-lineas-m` | F, T, C, D, I, L; raíz 412 fija | U1 (M-01, M-07), B-04, B-05, B-06, B-09 | Montos por parte y «Filtrar» en la agenda; intereses del año cuando estén en el JSON |
| `creditos-p17-lineas-e` | F, T, C, D, I, L; página fluida con `expand: fill`; Estados vacío, carga y error | U1 (M-05, M-07), M-04, M-08, B-03, B-06 | Promociones como franjas en el mes; Marcar pago (G-06); barra lateral canónica |
| `creditos-p18-credito-m` | F, T, C, D, I, L | U1 (M-01, M-02, M-04, M-07), A (M-03, M-09), B-01, B-02, B-06, B-10 | Partes en el calendario; serie de disponible en Detalle; gráfica legible a 412; «Capturar» de 44 px; hoja Marcar pago |
| `creditos-p18-credito-e` | F, T, C, D, I, L; Estados con «Falta capturar · Capturar» (K18) | U1 (M-04, M-07), M-08, B-02, B-03, B-06 | Hoja Marcar pago por parte; compras de la tarjeta y de la tienda; barra lateral canónica |
| `creditos-p19-compra-m` | F, T, C, D, L; form con `label for` en los 10 campos | B-07 (input date, onsubmit) | Fecha como texto `12-oct-2026` o con formato; quitar `onsubmit` |
| `creditos-p19-compra-e` | F, T, C, D, L; Estados Calculando, error y dato inválido | M-08, B-07 | Igual que móvil; barra lateral canónica |
| `creditos-p20-comparador-m` | F, T, C, D, I, L; 4 óptimos apilados, tabla plegada | U1 (M-06), B-07 (class duplicado), B-08 | Compras vivas máx., diferencia entre quincenas y rendimiento por peso |
| `creditos-p20-comparador-e` | F, T, C, D, I, L; Estados (ninguno, faltan restricciones, evaluando, no disponible) | U1 (M-06), M-08, B-03 | Igual que móvil; barra lateral canónica |

## Hallazgos

### Medios

- **M-01 · Agenda y calendario móviles sin montos por parte** (N-192, N-208, F-4 paso 2, DESIGN §7.8). En `creditos-p17-lineas-m` (pestaña Calendario), los renglones de `agenda` no traen campos `hogar` ni `karen`; solo dicen «Hogar pendiente · Karen pendiente». En `creditos-p18-credito-m`, `cal` sí trae `r.hogar` y `r.karen`, pero el markup no los pinta. F-4 espera «Vie 30 oct 2026 · $X · hogar $Y · Karen $Z · pendiente». Tampoco hay botón «Filtrar» (U1 §1.3; filtros línea, familia y estado de N-192). Corrección: agregar a cada fila el chip de la otra persona con su monto (o «Hogar $453.49 · pendiente / Karen $453.49 · pendiente») y un botón «Filtrar».
- **M-02 · P-18 móvil sin la serie de disponible** (N-203; U1 §6.4 «Serie A/B y tasas en Detalle»). `creditos-p18-credito-m` › Detalle solo muestra «mínimo A · mínimo B». Corrección: lista de eventos (fecha, tipo, monto, capital repuesto, interés no repuesto, disponible A y B) con las primeras 4 y «Mostrar 16 más».
- **M-03 · Texto de la gráfica ilegible en móvil** (DESIGN §7.7 ejes en `font.size.caption` de 12 px; U5-06). En `creditos-p18-credito-m` › Tendencias, el `<svg viewBox="0 0 640 220" width="100%">` vive en una tarjeta de unos 348 px y escala a 0.54, así que las etiquetas de 12 unidades se ven de unos 6.5 px. Corrección: `viewBox` con el ancho real del móvil (unos 348) y texto de 12, o etiquetas de eje y valores finales en HTML fuera del SVG.
- **M-04 · «Marcar pago» por parte no está dibujado** (N-210; F-4 pasos 3 y 4; DESIGN §7.8 «Tocar un pago abre su vista previa (G-06) con Marcar pago»). En P-18 (m, e), `markPay` solo cambia a la pestaña Calendario, y en P-17 el calendario no abre ninguna vista previa. Corrección: un `sc-if` con la hoja (móvil) o el panel lateral (escritorio) del pago: compra de origen, capital repuesto, casillas «Parte del hogar» y «Parte de Karen», fecha real (hoy) y monto, y «Guardar».
- **M-05 · El calendario mensual de escritorio no tiene promociones** (N-192 «promociones vigentes como franjas»; U1 §2.1 «calendario (pagos y promociones)»). En `creditos-p17-lineas-e` › Calendario, la cuadrícula de octubre solo tiene pagos y el corte. Bodega Aurrera (19-oct-2026 a 18-nov-2026) debería verse como franja. Corrección: dibujar la franja sobre los días de vigencia, con comercio, plazo y compra mínima.
- **M-06 · Salida del comparador incompleta** (N-226). En `creditos-p20-comparador-m` y `-e` faltan compras vivas máximas (el JSON trae `compras_vivas_max` y `renderVals` expone `vivas`, pero no se usa), la diferencia entre quincenas (`pago_max − pago_min`, derivable) y el rendimiento de compra por peso (no está en el JSON). Corrección: agregar las columnas «Compras vivas máx.» y «Diferencia entre quincenas» a la tabla y a las tarjetas. Para el rendimiento, pedir el dato a `datos-sinteticos.json` o marcarlo «sin dato».
- **M-07 · Elementos de U1 que faltan porque el JSON no los trae.** N-191 «intereses pagados en el año» no está en P-17 (m, e). N-202 «Compras de la línea» solo aparece en la revolvente: la tarjeta (`creditos-p18-credito-*`) solo muestra MSI y la tienda solo abonos y una tarjeta «Compras vivas». El autor no los inventó (correcto según P3-12); el hueco está en `datos-sinteticos.json` (`creditos.totales` sin intereses del año; `creditos.tarjeta` sin `compras`). Corrección: el dueño de los datos agrega `intereses_anio_*` y `tarjeta.compras`, y el autor los dibuja.
- **M-08 · Barra lateral de escritorio no colapsable** (P3-08 «barra lateral colapsable»; DESIGN §7.9). Los 4 artboards `-e` usan una barra reducida: sin botón «‹‹» (`panel-left`), sin iconos de sección (solo Créditos tiene), sin chevrons, avatar ni interruptor de tema. La canónica `shell/shell-barra-lateral-e.dc.html` sí trae «Ocultar barra lateral». Corrección: reutilizar la anatomía de la barra de shell con Créditos desplegado y «Todas las líneas» o la línea actual con `aria-current="page"`.
- **M-09 · Objetivo táctil de 20 px en móvil** (DESIGN §9.3). En `creditos-p18-credito-m` › Detalle › Tasas por plazo, el enlace «Capturar» va dentro de «Otros plazos: sin capturar · Capturar» con `line-height` de 20 px y sin `min-height`. Corrección: clase `tap` (44 px) o botón propio.

### Bajos

- **B-01** Iconos de 14 px (28 usos) y 16 px (10 usos); DESIGN §5.2 pide solo 20 y 24 px. Además, `stroke-width` 1.5 en el `viewBox` de 24 también a 20 px (DESIGN pide 1.8 a 20 px). Todas las rutas.
- **B-02** Ejes de dinero `$0`, `$5,000`, `$10,000` y `$15,000` en lugar del compacto «$5 mil» (DESIGN §8.5), y sin «Ver como tabla» (§7.7). Están en `creditos-p18-credito-m` y `-e`.
- **B-03** Las insignias «Vacío», «Carga» (Estados) y «Próxima» (promociones) usan el rol deuda (`b-debt`) para algo que no es deuda; conviene una insignia neutra. Están en P-17 e, P-18 e, P-19 e y P-20 e.
- **B-04** La barra inferior móvil no tiene iconos (`house`, `utensils`, `plus`, `shopping-cart`, `search`), usa fondo `ui-surface-card` en lugar de `ui-surface-raised` y sus destinos miden 44 px en lugar de 56 (DESIGN §7.10). Los 4 artboards `-m`.
- **B-05** El prototipo tiene enlaces incompletos (no rotos). «Abrir menú» es un `<button>` sin destino aunque existe `../shell/shell-barra-lateral-m.dc.html`. «Editar condiciones», «Confirmar en Ajustes › Créditos y condiciones», «Capturar» y «Faltan restricciones · Capturar» van a `#`, aunque ya existe `../ajustes/ajustes-p26-s05-creditos-condiciones-{m,e}.dc.html`.
- **B-06** Textos incoherentes:
  - «Walmart · despensa bimestral (programada) (programada)», duplicado en el calendario de P-18 (m, e).
  - «Hogar pendiente» en pagos que son 100 % de Karen (Mercado Libre · audífonos de Karen), en P-17 y P-18.
  - «Límite $18,000.00» en el renglón del préstamo de P-17 m › Detalle, aunque es el monto original.
  - «54.0000 % mensual» para Coppel, cuando el JSON da tasa anual con base mensual.
- **B-07** Detalles de markup:
  - `<li class="card" style="…" class="small">` lleva `class` duplicado en `creditos-p20-comparador-m`; el navegador descarta el segundo.
  - `onsubmit="return false"` está como texto en P-19 (m, e); el runtime solo admite manejadores `{{ fn }}`.
  - `<input type="date" value="2026-10-12">` en P-19 se muestra según el idioma del navegador y no como «12-oct-2026» (U5-07).
- **B-08** El puntaje del comparador (76, 70, 50 y 43) es derivado: el JSON no lo trae (supuesto 3 del autor). Además, el JSON tiene 3 pesos (40/30/30) y N-225 pide 4 objetivos iguales por defecto. Hay que corregirlo en `datos-sinteticos.json`, no en el artboard.
- **B-09** La agenda de P-17 dice «Vie 30 oct» sin año y el calendario de P-18 dice «15-oct-2026» sin día de la semana (U1 §0.4: «Vie 30 oct 2026» en calendarios y pagos).
- **B-10** Las pestañas van en 16 px y peso medio también en reposo (DESIGN §7.10: `font.size.small`, medio solo la activa). En P-18 m › Detalle se ven todas las propiedades en lugar de «4 + Mostrar N más», y las acciones secundarias quedan visibles en lugar de ir en «⋯» (U1 §1.3). El sprite de iconos es un segundo nodo de primer nivel dentro de `<x-dc>` (format.md muestra una sola raíz). Los `$preview` y `h` son estimados (P-18 m de 1800 puede quedar corto con 5 compras). `is_interactive: true` en P-17, P-18 y P-19 sigue a format.md pero difiere de la plantilla de la orden (`false`); lo decide el orquestador.

## Fidelidad: pantalla → elementos de U1

| Pantalla | Presentes | Faltantes |
|---|---|---|
| P-17 | N-190 (todas las líneas, ordenadas por próximo pago, disponible A/B con etiqueta); N-191 deuda viva, quincena y mes con partes, intereses del mes; N-193 aviso fijo; N-194 (12 meses por línea en e; total en m); galería de 3 columnas (e) y tarjetas apiladas (m); calendario unificado con estado y corte; acciones primaria y secundarias; Estados vacío, carga y error (e) | Intereses del año (M-07); franjas de promoción en el mes (M-05); montos por parte y filtro en móvil (M-01) |
| P-18 | Variante por tipo (revolvente, tarjeta, préstamo, tienda); N-200 condiciones; N-201 tarjetas; N-203 serie (e, 10 de 20 + «Mostrar 10 más»); N-204 tasas; N-205 ciclo y MSI; N-206 y N-207; N-208; N-209 gráfica A/B y deuda por mes; Estados con K18 (e) | N-203 en móvil (M-02); N-210 / F-4 (M-04); N-202 en tarjeta y tienda (M-07); partes en el calendario móvil (M-01) |
| P-19 | N-215 compra, partes y pagos con estado por parte, gasto vinculado; N-216 vista previa; N-217 promociones; alta y detalle; Estados Calculando y error (e) | — |
| P-20 | Datos de la compra, restricciones solo para esta comparación (K16), topes como no duros, 4 pesos, recomendación con puntaje, 4 óptimos (apilados en m, en fila en e), tabla de escenarios con descartados y motivo, total evaluado; Estados (e) | Compras vivas máx., diferencia entre quincenas y rendimiento por peso (M-06) |

## Contraste recalculado (U5-04), hex de `tokens.json`

| Par (primer plano / fondo) | Tipo | Claro | Oscuro |
|---|---|---|---|
| ui-text-default / ui-surface-card | texto | 10.76 | 10.76 |
| ui-text-secondary / ui-surface-card | texto | 5.04 | 7.97 |
| ui-text-secondary / ui-surface-page | texto | 5.40 | 10.43 |
| ui-text-link / ui-role-warning-surface | texto | 5.15 | 8.01 |
| ui-text-link / ui-surface-card | texto | 5.16 | 7.94 |
| ui-text-on-primary / ui-action-primary-bg | texto | 5.77 | 10.47 |
| ui-role-debt / ui-role-debt-surface | texto | 5.29 | 6.05 |
| ui-role-success / ui-surface-card | texto | 4.74 | 6.42 |
| ui-role-warning / ui-role-warning-surface | texto | 5.12 | 6.09 |
| ui-text-secondary / ui-surface-hover | texto | 4.51 | 5.50 |
| ui-border-control / ui-surface-card | interfaz | 3.47 | 6.09 |
| ui-border-focus / ui-surface-card | interfaz | 3.56 | 6.03 |
| ui-chart-2 / ui-surface-card | interfaz | 3.15 | 7.73 |
| ui-person-karen / ui-surface-card | interfaz | 3.16 | 4.68 |

En total, 26 pares del markup y 16 extra (estados hover, foco, selección, series y personas): 0 bajo umbral. El más justo es `ui-text-secondary` sobre `ui-surface-hover` (4.51) en el hover de los botones terciarios.

## Verificación (scripts propios, `python3 -I` y node, en `datos/tmp/rev-cred/`, borrados al final)

- Forma: doctype, `lang="es-MX"`, `title` propio, `./support.js`, `../ds/finanzas-hogar/tokens.css`, un `x-dc`, `helmet` con `style`, `data-dc-script` con JSON válido, `class Component extends DCLogic` y `renderVals()`: 8/8. `$preview` = `w × h` del fragmento: 8/8. `html.parser`: 0 etiquetas sin cerrar o huérfanas. Nombres `creditos-pNN-…-{m,e}.dc.html` válidos y stems únicos.
- Raíz móvil `width:412px; box-sizing:border-box`: 4/4. Escritorio: raíz flex-wrap sin ancho en px, `main` `flex:999 1 560px` con contenedor `max-width:1120px`, `expand: fill` en 4/4.
- Tokens: 0 hex y 0 `rgb()` o `hsl()` fuera de `data-props`; todas las `var(--x)` existen en `tokens.json`.
- Huecos: `renderVals()` ejecutado en node con un `DCLogic` mínimo, en 1 a 25 estados por archivo (pestañas × líneas en P-18). 0 huecos sin resolver, con alcance de `sc-for` anidados.
- Enlaces: todos los `href` relativos existen en el lienzo y son del mismo ancho; el resto es `#`.
- Objetivos táctiles móviles, por estilo calculado (clases del `helmet` + estilo en línea + tokens): 1 control < 44 px (M-09). Los `summary` usan `details>summary` con 44 px.
- Idioma: 0 palabras en inglés visibles (los tipos `start`, `payment` y `purchase` del JSON ya vienen como Inicio, Pago y Compra), 0 emoji y 46 fechas con día de la semana correctas.
