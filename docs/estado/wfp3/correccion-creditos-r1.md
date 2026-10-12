# correccion-creditos-r1 · vuelta 2 · Créditos (P-17 a P-20)

Agente: `correccion-creditos-r1` · Fecha: 12-oct-2026 · Alcance: los 8 artboards de `docs/design/maquetas/lienzo/project/creditos/`, el fragmento `docs/design/maquetas/lienzo/fragmentos/creditos.json` y este reporte. Copia de la vuelta 1 (P3-06): `docs/design/salida/iter-1/creditos/` (8 artboards y el fragmento, tal como estaban antes de tocar nada).

## Hecho (9 hallazgos, todos de severidad media)

1. **Partes y filtro en agenda y calendario móvil.** P-17 m (Calendario): cada pago muestra «Hogar $X · pendiente» y «Karen $Y · pendiente» con su monto; los pagos que son solo de Karen o solo del hogar ya no dicen la otra parte. El dato de hogar y Karen viene del JSON (se copió de la agenda de escritorio). Fecha con día y año («Vie 30 oct 2026»). Botón «Filtrar» que abre una hoja con Línea, Familia (Hogar y Karen, Solo hogar, Solo Karen) y Estado (Todos, Pendientes, Pagados), con filtrado real y mensaje «Sin pagos con este filtro». P-18 m (Calendario): mismas partes con monto usando `r.hogar` y `r.karen`, y el mismo «Filtrar» (Línea, Familia, Estado).
2. **Serie de disponible en P-18 m › Detalle.** Los 20 eventos de `creditos.kueski.serie_disponible` (fecha, tipo, monto, capital repuesto, interés no repuesto, disponible A y B). Se ven los 4 primeros y el botón «Mostrar 16 más» (alterna a «Mostrar menos»). Los 10 primeros coinciden con los que ya tenía el escritorio.
3. **Gráfica A/B legible en móvil.** `viewBox="0 0 348 220"` con `width="100%"` en una tarjeta de 348 px: escala 1, texto de 12 px reales. Ejes en forma compacta («$5 mil», «$10 mil», «$15 mil»). Los valores finales salieron del SVG a una línea HTML («Al 15-abr-2027: lectura A $16,558.40 · lectura B $12,000.00»).
4. **«Marcar pago» por parte (F-4).** Hoja en móvil (P-18 m) y panel lateral en escritorio (P-17 e, P-18 e) con scrim, `role="dialog"`: compra de origen, monto programado, fecha programada, capital repuesto (de `pagos_programados`), casillas «Parte del hogar» y «Parte de Karen» con su monto y estado (separadas; solo aparece la parte que existe), fecha real «21-oct-2026» (hoy), «Guardar» y «Cancelar». El botón primario de P-18 abre la hoja del próximo pago pendiente; cada fila de los calendarios (P-18 m, P-18 e, agenda de P-17 e) trae «Marcar pago» o «Ver pago»; en la cuadrícula mensual de P-17 e los 6 pagos de octubre son botones que abren su vista previa (G-06). «Guardar» marca el pago como pagado en la fila.
5. **Promociones como franjas (P-17 e › Calendario).** Bodega Aurrera, del 19-oct al 18-nov-2026, se dibuja como franja continua sobre los días 19 a 31 de octubre y el 1 de noviembre, con rótulo «Bodega Aurrera · 4 quincenas · compra mínima $2,000.00» el 19 y «hasta el 18-nov-2026» el 26. Leyenda «Promoción vigente».
6. **Comparador completo.** «Compras vivas máx.» y «Diferencia entre quincenas» (`pago_max − pago_min`, calculada del JSON: E1 $0.00, E2 $0.02, E3 $0.01, E4 $0.02, E5 $0.01, E6 $0.00; compras vivas 3 y «No aplica» en la tarjeta) en la tabla de escritorio (3 columnas nuevas, con «Rendimiento por peso»), en las tarjetas de óptimos (m y e), en la recomendación y en la lista de escenarios móvil. «Rendimiento por peso» no está en el JSON: se muestra «Sin dato».
7. **Elementos que el JSON no trae (parcial).**
   - «Intereses pagados en el año» ya está dibujado en P-17 m y P-17 e con el valor «Sin dato» (no se inventó la cifra).
   - «Compras de la línea» ahora aparece en la tarjeta (Liverpool · laptop, meses sin intereses: total 12 × $1,174.92, pagos 5 de 12, saldo $8,224.44) y en la tienda (Colchón matrimonial, 22-ago-2026, total de los 10 abonos, 2 de 10, saldo $7,067.09), en P-18 m (tarjetas) y P-18 e (tabla). Todo se deriva de `tarjeta.msi` y `tienda.abonos`; no se inventó ninguna fila.
8. **Barra lateral canónica en los 4 artboards de escritorio.** Se reutilizó la anatomía de `shell/shell-barra-lateral-e.dc.html`: cuenta con avatar, botón «Ocultar barra lateral» (panel-left), Buscar y Registrar gasto con atajos, iconos de las secciones, chevrons, interruptor de tema Sistema/Claro/Oscuro. Créditos desplegado («Plegar Créditos», `aria-expanded="true"`) con hijos: Todas las líneas, las 4 líneas en orden alfabético (N-005), Registrar compra a crédito y Comparador de escenarios; la subpágina actual lleva `aria-current="page"` (en P-18 e sigue a la línea elegida). Las clases llevan prefijo `sb-` para no chocar con `.seg`/`.btn` de cada artboard.
9. **«Capturar» de 44 px.** En P-18 m › Tasas por plazo es un enlace propio con clase `tap` (min-height 44 px) en su renglón y apunta a Ajustes › Créditos y condiciones.

Extras baratos del mismo toque (no estaban pedidos): el texto «(programada) (programada)» ya no se duplica (B-06); «Editar condiciones», «Confirmar en Ajustes» y «Capturar» enlazan a `ajustes/ajustes-p26-s05-creditos-condiciones-{m,e}` en vez de `#` (B-05); se quitó el atributo `class` duplicado de P-20 m (B-07). No se tocó el resto de bajos (B-01 iconos de 14/16 px preexistentes, B-02 ejes del escritorio, B-03, B-04, B-08, B-10).

Archivos modificados (7 de 8; P-19 m no cambió): `creditos-p17-lineas-m`, `-p17-lineas-e`, `-p18-credito-m`, `-p18-credito-e`, `-p19-compra-e` (solo barra lateral), `-p20-comparador-m`, `-p20-comparador-e`; fragmento `creditos.json` (solo `h` de P-17 e, P-20 m y P-20 e).

## Verificación (script propio, `python3 -I`, más una prueba en node de `renderVals()`)

- Forma, 8 de 8: `<!doctype html>`, `lang="es-MX"`, `<title>`, `./support.js`, `../ds/finanzas-hogar/tokens.css`, un `x-dc`, `helmet`+`style`, `data-dc-script`, `class Component extends DCLogic`, `renderVals()`; JSON de `data-props` parsea y su `$preview` es igual a `w`×`h` del fragmento.
- `html.parser`: 0 etiquetas sin cerrar ni huérfanas. Ids únicos y cada `input` con su `label for`.
- Móvil: raíz de 412 px fija con `box-sizing`; escritorio: sin ancho en px en la raíz, `expand: fill`.
- 0 hex literales, 0 `rgb()`/`hsl()` fuera de `data-props`; todas las `var(--x)` existen en `docs/design/sistema/project/tokens.css`.
- Cifras prohibidas de P3-12 (la lista del grep de `docs/criterios/P3-claude-design.md`, con y sin coma): 0 en los 8 artboards y en el fragmento. 0 credenciales, 0 emoji, 0 iframes, `innerHTML` o URL externas, 0 «-$» (signo menos).
- Enlaces `<a href>` relativos: todos existen en el lienzo. Fragmento: JSON válido, lista exactamente los 8 archivos en `boards` y `order`.
- `renderVals()` ejecutado en node con un `DCLogic` mínimo en 3 a 6 estados por archivo (pestañas, línea, hoja de pago abierta, filtros, serie expandida, y los manejadores `markPay`, `saveSheet`, `closeSheet`, `toggleSer`, `openFil`, `clearFil`): 0 errores y 0 variables del markup sin resolver.
- Resultado del script: `RESULTADO: OK`.

## Supuestos

- «Hoy» es el 21-oct-2026 (`meta.hoy`), así que la fecha real propuesta en la hoja es «21-oct-2026».
- Las casillas «Parte del hogar» y «Parte de Karen» alternan su estado en la hoja, pero «Guardar» marca el pago completo como pagado (el prototipo no lee el DOM); la separación por parte queda dibujada, no persistida.
- «Compra de origen» de la tarjeta/estado de cuenta y del préstamo se redactó con datos del JSON («Compras del ciclo del estado de cuenta», «Préstamo de $18,000.00 · desembolso 18-jun-2026»); el capital repuesto sale tal cual de `pagos_programados`.
- Total de la compra de la tienda ($8,833.87) y total de los MSI ($14,099.04) son sumas del JSON (abonos y mensualidades), no cifras nuevas.
- El filtro «Familia» se entiende como Hogar / Karen (como en el filtro de escritorio «Familia: hogar y Karen»).
- Los `$preview`/`h` nuevos son estimados (no hay navegador): P-17 e 2600 a 2700, P-20 m 2400 a 2650, P-20 e 1850 a 1910.

## Bloqueos

- Hallazgo 7, parte de datos: `docs/design/maquetas/datos-sinteticos.json` no tiene `creditos.totales.intereses_anio_*` ni `creditos.tarjeta.compras` ni la fila de compra de la tienda con capital, intereses y partes; ese archivo no está en mi alcance y no lo toqué. Hasta que el dueño de los datos los agregue, «Intereses pagados en el año» dice «Sin dato» y las compras de la tarjeta y de la tienda salen solo de lo derivable. Tampoco existe el «rendimiento de compra por peso» (hallazgo 6): «Sin dato».

## Pendientes

- El dueño de `datos-sinteticos.json` agrega `intereses_anio_*`, `tarjeta.compras` (con comercio, fecha, capital, intereses, partes) y el rendimiento por peso del comparador; después se sustituyen los «Sin dato» y las filas derivadas.
- El orquestador debe medir las alturas reales (P-18 m con la hoja abierta y la serie expandida es más alta que el `$preview`; P-17 e, P-20 m y P-20 e son estimados).
- Sin hacer (bajos de la revisión): B-01, B-02 (ejes y «Ver como tabla»), B-03, B-04 (iconos de la barra inferior móvil), B-08 (pesos 40/30/30 en el JSON), B-10, y el año en el calendario de P-18 e. `is_interactive` sigue como estaba (decide el orquestador).

## Consumo aproximado

Lectura ≈ 300 KB (informe de la vuelta 1, 8 artboards por secciones, barra lateral de shell, JSON por consultas); escritura ≈ 100 KB (scripts de parche y este reporte; los artboards crecieron ≈ 105 KB en disco por la barra lateral, las hojas y la serie); ≈ 34 acciones de herramienta.
