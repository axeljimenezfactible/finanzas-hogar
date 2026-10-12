# Datos sintéticos de las maquetas · LEEME

Versión 1.4 (corrección r5) · 12-oct-2026 · decisión K45.1 (cifras sintéticas obligatorias) y criterio P3-12.

`datos-sinteticos.json` es el **único archivo de datos** de todos los artboards de la Pista 3. Todo está **inventado**: tiene la forma y la magnitud de lo real, pero ninguna cifra viene del reporte, de `data.json`, de `model.json` ni de Kueski. Los nombres **Axel, Ximena y Karen** sí aparecen (Q2-98 A); las tiendas y los bancos se llaman por su nombre público, con condiciones inventadas. Si Axel pide seudónimos, se cambian solo aquí.

## Reglas para los autores de artboards

1. **Lee de este archivo, no inventes cifras.** Si una pantalla necesita un dato que no está, anótalo en tu reporte como supuesto; no lo calcules a mano con otra magnitud.
2. **Dinero** en centavos enteros (`*_cents`). Donde una pantalla lo muestra hay un campo `*_fmt` hermano con el texto `$1,234.56` (negativos con signo menos `−$12.30`). Donde no hay `_fmt`, formatea con la misma regla: miles con coma, dos decimales.
3. **Hoy es el miércoles 21-oct-2026** (`meta.hoy`); no es la fecha real. Mes en curso 2026-10, quincena en curso del 16 al 31 oct (Q2), semana en curso 2026-W43 (lun 19 al dom 25 oct).
4. **Fechas** en ISO; en pantalla `21-oct-2026` o `Mié 21 oct 2026`. Los porcentajes de avance son enteros (ROUND_HALF_UP); las tasas son texto con sus decimales (`24.3120 %`).
5. **Las partes suman el total**: parte del hogar + parte de Karen = total; los pagos suman capital + intereses; los porcentajes salen de plan/gasto. `tools/ds/verificar_sinteticos.py` lo comprueba (57 comprobaciones).
6. **Color por persona y tienda**: usa `token_color` (`color.person.axel`, `color.store.walmart`, ...) y escribe siempre el nombre en texto (DESIGN §2.6).
7. **Sin datos no es cero**: los estados vacíos se dibujan con `estados_alternos` y con el texto «—»; los datos por defecto no los tienen.
8. **Nada real**: no copies cifras de `reporte-p01..p11.png` ni del reporte; no pegues credenciales ni correos.

## Esquema (llaves de primer nivel)

| Llave | Qué trae | Pantallas |
|---|---|---|
| `meta`, `formato`, `indice_pantallas` | hoy, mes, semana, quincena; formatos; qué llaves usa cada pantalla | todas |
| `personas`, `cuentas_y_formas_de_pago` | Axel (titular), Ximena (co-titular), Karen (externa, sin cuenta); reparto Hogar/Karen/Conjunto; formas de pago | todas |
| `ingresos` | esperado (q15 y q30), quincenas (una capturada, una esperada), extras, total del mes; **`tendencia`** (P-08 › Tendencias, N-103): 6 meses, de mayo a octubre de 2026, cada uno con `quincena_15_cents`, `quincena_30_cents`, `extras_cents` y `total_cents`; octubre lleva `incluye_esperado` (la quincena del 30 aún no se captura) | P-08, P-07 |
| `categorias` | 13 categorías con plan, real del mes, diferencia, % usado, estado (normal/aviso/excedido), fecha de cruce de 80 y 100 % | P-06, P-03, P-26.02/.04 |
| `presupuesto` | resumen del mes, Q1 y Q2, despensa devengada, alertas, tendencia, gastos fijos recurrentes, despensa programada | P-06, P-26.02 |
| `gastos`, `gastos_resumen`, `reparto` | **`gastos_resumen.tendencia_por_categoria`** (P-02 › Tendencias N-044; P-06 N-085 al elegir una categoría): 6 meses con plan, real y % usado de las 13 categorías (cada mes suma el plan y el real de `presupuesto.tendencia`; octubre = `categorias`), `seis_mayores` (las 6 con más gasto del hogar devengado de mayo a octubre, U1 N-044) y `otras_*_cents` (las 7 restantes, para la serie «Otras»); 33 gastos (fijos, variables, venideros; contado, Kueski, tarjeta); partes hogar/Karen; `devengado` (por mes) y `efectivo` (fechas en que sale el dinero; `null` = el mismo día) | P-02, P-03, P-04, P-05 |
| `tickets` | **P-04**: ticket confirmado del gasto `g13` (`confirmados[0]`: 11 productos y el renglón «Ajuste del ticket», cada uno con producto, presentación, cantidad, precio unitario, importe, cuenta y partes hogar/Karen; totales `por_cuenta`). **P-05**: `lectura_pendiente` (OCR de un ticket de Walmart con `texto_leido`, producto y cuenta sugeridos, `confianza` alta/media/baja —renglones 9 y 10 bajos— y la `diferencia_cents` contra el total leído) y `por_confirmar` (2 lecturas, por fecha descendente; `ocr2` tiene `borrador` en `borradores[0]` con sus correcciones). Solo `g13` trae renglones; los demás gastos con `tiene_ticket` muestran el adjunto sin renglones | P-04, P-05 |
| `ahorro` | fondo, 2 metas (una en ruta, una atrasada), 19 movimientos, tendencia | P-09 |
| `creditos` | 4 líneas (Kueski Pay revolvente, Banorte Oro tarjeta, préstamo personal, crédito Coppel); compras Kueski y programadas; `pagos_programados` (hogar/Karen, estado por parte); serie de disponible A/B; tasas por plazo; ciclo de la tarjeta; amortización y abonos | P-17, P-18, P-19 |
| `promociones`, `comparador` | 3 promociones (una vigente hoy); compra planeada, restricciones, 6 escenarios, óptimos, descartados con motivo | P-20, P-26.05 |
| `tiendas` | Walmart, Bodega Aurrera, Sam's Club (membresía en evaluación), Mercado Portales (sábado), Tianguis de la Narvarte (martes) | P-26.06, P-13, P-14 |
| `productos` | 33 productos con presentaciones y 105 precios por tienda y fecha (22 viejos, de 30 días o más; `fuente` ticket/manual/profeco); **`historial_precios`** (P-15 › Tendencias, N-171): por cada par presentación × tienda, de 3 a 5 puntos `{fecha, precio_cents, fuente}` en orden de fecha, desde mayo de 2026; el último es el precio vigente de `precios`. El precio por unidad base se calcula: `precio_cents` ÷ `cantidad_base` de la presentación (no hay flotantes en los puntos) | P-15 |
| `canasta` | comparativa de la lista de la semana: total por tienda, idénticos, sustitutos, sin precio, ganadores, detalle por artículo, membresía de Sam's | P-14, Inicio |
| `recetas`, `plan_semana` | 12 recetas con ingredientes, costo y kcal; **`uso`** de cada receta (P-12 N-141 y el orden de N-122): `veces_8_semanas` (= `en_semana_actual_hasta_hoy` + `en_7_semanas_previas`; la ventana va del 27-ago al 21-oct-2026) y `proxima_fecha` (la primera del plan semanal desde hoy, o `null` = «—» si no está en el plan); plan 7 días × 4 tiempos × 2 personas con «fuera» y 2 casillas sin armar | P-10, P-11, P-12 |
| `lista_compras`, `inventario` | lista agrupada por tienda o mercado; 14 artículos de inventario (stock por unidad y ritmo de consumo; resumen 1 agotado / 5 reponer / 8 suficiente) y movimientos | P-13, P-16, P-26.14 |
| `salud` | perfiles, IMC, meta de kcal sugerida y confirmada, 8 semanas de peso, kcal de hoy, tendencia de 12 semanas | P-21, P-22 |
| `reportes`, `reporte_mensual`, `estado_cuenta_karen` | reportes generados, próximos envíos, avisos enviados (uno falló), programación (21:30, dom 19:00, día 1 08:00; avisos 09:00); reporte de octubre (parcial); estado de cuenta de Karen | P-23, P-24, P-25 |
| `proyeccion`, `deuda_viva` | 24 quincenas de efectivo (en las 24, `variable_estimado_cents` y `ahorro_planificado_cents` van en `null` y el acumulado parte de 0: ver «Notas de lectura (r5)»); 12 meses devengado contra efectivo; serie de deuda viva y deuda a una fecha | P-07, Inicio › Tendencias |
| `calendario_hogar`, `participacion_karen` | 74 eventos del 1-oct-2026 al 17-abr-2027, por tipo: octubre completo (gastos fijos ya pagados el 1, 6, 12 y 14; despensa del 12; visitas a mercado realizadas del 3, 6, 13, 17 y 20; pagos de crédito ya pagados; reporte del 1) y lo que viene hasta las fechas meta de ahorro; los ya hechos llevan `estado` («realizado», «realizada», «pagada», «enviado»); bloque «Participación de Karen» | P-01 |
| `inicio` | las 6 tarjetas, segundo nivel y apuntadores a las tendencias. Las tarjetas usan los nombres de campo de C1: `mes_vs_plan` (`month`, `actual_cents`, ...), `flujo_quincenal`, `proximo_pago.next_payment` (`due_date`, `credit_id`, `amount_cents`, `household_share_cents`, `external_share_cents`, `days_left`, `payment_count`, `is_overdue`; `dia_semana` y los `_fmt` son de presentación), `comidas_semana`, `mejor_canasta` y `ahorro.next_goal` (nombres de `SavingsGoal`: `name`, `target_cents`, `target_date`, `saved_cents`, `progress_pct`, `planned_monthly_cents`, `required_monthly_cents`, `status` `on_track`/`behind`; más `id`, `kind` y `etiqueta_estado`). Son nombres de C1 los de las tarjetas; `segundo_nivel` y `ahorro.metas` (P-09) usan nombres en español | P-01 |
| `ajustes`, `estados_alternos`, `recientes` | usuarios, sesiones, parámetros (incluye `comparador_max_s` = 5 s, K50), respaldos, IA, PROFECO, Telegram; **`ajustes.notificaciones`** (P-26.09: las 4 alertas `payment_due`, `budget_threshold`, `stale_prices`, `meal_plan_missing`, cada una con activa, destinatarios, anticipación y hora 09:00; pago próximo [1, 0], precios viejos 1 día —K50—, umbrales 80 y 100, plan sin armar los domingos; `stale_prices` viene apagada para dibujar ese estado); **`ajustes.datos_semilla`** (P-26.11: 3 renglones con estado de carga y su número de registros: `registros`, `registros_texto` y `registros_estado`; los nombres no llevan cifras); variante «pago vencido» (con los nombres de C1); recientes | P-26, G-03 |

Los nombres de campo siguen a C1 (`docs/contratos/openapi.yaml`) en las 6 tarjetas de Inicio y en `estados_alternos.pago_vencido` (`used_pct`, `closing_balance_cents`, `due_date`, `availability_reading`, ...). En el resto del archivo son nombres en español de este archivo. En `credit_id` el valor es el id de la línea de este archivo (`kueski`), no un entero.

**Estados con menos de 12 meses.** Las series mensuales traen 6 meses (mayo a octubre de 2026): `presupuesto.tendencia`, `ahorro.tendencia`, `ingresos.tendencia` y `gastos_resumen.tendencia_por_categoria`. En T-1 y en las demás vistas de «Últimos 12 meses» dibuja los 6 meses anteriores sin cifra, como estado vacío, y no inventes valores. (Supuesto: U3 §4 fija la ventana de 12 meses con el subtítulo «Últimos 12 meses», y U1 N-044, N-085 y N-103 dicen «del primer mes con datos al mes en curso (máximo 12)», que sería de 6 meses sin huecos; las dos lecturas se concilian dibujando la ventana fija de U3 con los meses sin datos en estado vacío. Si Axel prefiere la lectura de U1, se acorta la ventana a 6 meses.) El historial de precios (`productos[].historial_precios`) llega hasta mayo de 2026 y es de 3 a 5 puntos por tienda: dibuja solo esos puntos, sin interpolar ni completar.

**Textos que repiten cifras.** Los resúmenes diarios de `reportes.avisos_enviados` (P-23) calculan su monto y su conteo de gastos con los gastos realizados de ese día; el verificador lo comprueba. Si cambias un gasto y regeneras, el texto se actualiza solo.

## Cifras que deben coincidir entre pantallas

- Mes contra plan (Inicio) = `presupuesto.resumen` = suma de `categorias[].real_cents`: real $20,525.78 de $27,586.00, 74 %, faltan 11 días.
- Próximo pago (Inicio) = primer pago pendiente de `creditos.pagos_programados`: vie 30 oct 2026, 2 pagos ese día, hogar y Karen por separado.
- Disponible hoy de Kueski Pay: lectura B $6,881.60 (activa, pendiente de confirmar) y lectura A $11,440.00.
- Mejor canasta (Inicio) = primera tienda comparable de `canasta`; la lista de compras (P-13) la desglosa por tienda.
- Ahorrado total (Inicio) = `ahorro.ahorrado_total_cents`; la meta más próxima es la atrasada.
- La deuda viva por mes (Tendencias) es `deuda_viva.serie_mensual`; la de la fecha elegida, `deuda_viva.a_fecha`.

## Fechas que coinciden con el calendario real (inevitables)

La regla de quincena (ingresos y pagos el 15 y el 30, o el último día de febrero) hace que 15 fechas ISO coincidan con `oraculo/model.json`: `2026-10-30`, `2026-11-15`, `2026-11-30`, `2026-12-15`, `2026-12-30`, `2027-01-15`, `2027-01-30`, `2027-02-15`, `2027-02-28`, `2027-03-15`, `2027-03-30`, `2027-04-15`, `2027-04-30`, `2027-05-15`, `2027-06-15`. Ninguna va con un monto del reporte ni de Kueski. El resto de las fechas usa otros días (nunca el 9 ni el 10 después del 1-oct-2026). Las series con historia (peso desde 31-ago, comidas desde 27-jul, ahorro desde may-2026, tendencia mensual, préstamo, crédito de tienda y meses sin intereses) empiezan antes de oct-2026, porque una gráfica necesita pasado.

## Lo ilustrativo (no es regla de negocio)

El motor real manda (S2 financiero, S3 créditos, S4 recetas, S5 precios). Este archivo solo imita: la semántica de las lecturas A/B, el estado de una meta, el pago mínimo de la tarjeta (el mayor entre el porcentaje y el fijo), el gasto variable estimado (monto fijo por quincena), el orden de tiendas (solo con cobertura de 90 % o más), los sustitutos (solo presentaciones de tamaño comparable), la frecuencia de compra en Sam's (12 al año) y el devengo de préstamo, tienda y meses sin intereses (al pagarse). El «real» del mes no incluye los gastos venideros hasta que se marquen como pagados.

## Cómo se regenera y se comprueba

```
python3 -I /home/deepsec/finanzas-hogar/docs/design/tools/ds/generar_sinteticos.py /home/deepsec/finanzas-hogar/docs/design/maquetas/datos-sinteticos.json
python3 -I /home/deepsec/finanzas-hogar/docs/design/tools/ds/verificar_sinteticos.py
```

El generador es determinista (sin azar ni reloj): cambiar un dato y regenerar mantiene las sumas. El verificador hace el grep de P3-12 (cifras prohibidas con y sin coma, y en centavos), busca credenciales, comprueba sumas y porcentajes, compara fechas con `model.json` y `data.json`, y termina con `RESULTADO: OK`.

## Notas de lectura (ronda 3)

- **Datos semilla (P-26.11).** `ajustes.datos_semilla[2]` (gastos fijos e ingresos) dice «Pendiente de importar» con 10 registros, mientras P-06 ya muestra los 8 gastos fijos y P-08 el esperado q15/q30 vigente desde el 1-oct. Se leen como captura a mano; el CSV sigue sin cargarse y por eso la fila aún no dice «Cargada».
- **Sumas que se arman con los datos.** `ahorro.tendencia` trae solo `ahorrado_cierre_cents`: el aportado de cada mes sale de sumar `ahorro.movimientos` por mes (el planificado solo existe para octubre: `ahorro.planificado_mes_cents`; en los demás meses no lo dibujes). `deuda_viva.serie_mensual[].intereses_mes_cents` no separa pagados y por pagar: dibuja el total del mes, sin inventar la separación.

## Notas de lectura (r4, 12-oct-2026)

- **Quincena en curso (16 al 31 oct) · `proyeccion.flujo_por_quincena[1]`.** *(Reemplazada por la nota r5, abajo: la regla de null ahora vale para las 24 quincenas y el acumulado ya no es $12,827.78.)*
- **Inventario (P-16 y P-26.14) · K7.** Nivel de reposición vacío (`umbral_reposicion: null`) significa «Sin aviso» y nunca da «reponer». Regla del generador: `agotado` si la existencia es 0; `reponer` si hay nivel y la existencia es igual o menor que el nivel; `suficiente` en el resto. Huevo, leche y tortilla ahora traen nivel (huevo 12 pza, leche 2000 ml, tortilla 1000 g) y siguen en «reponer», con nivel visible. Croquetas, jabón para trastes, detergente y papel conservan nivel vacío («Sin aviso») y son «suficiente». `inventario.resumen` queda en 1 / 5 / 8.

## Notas de lectura (r5, 12-oct-2026) · flujo por quincena según S2-interfaz 1.3 §3.11

- **Regla (hallazgo A-3, reabierto en la iteración 4 de Inicio).** En `proyeccion.flujo_por_quincena`, **`variable_estimado_cents` y `ahorro_planificado_cents` van en `null` en las 24 quincenas** (§3.11: PROVISIONAL, bloqueos `S2-flujo-quincenal` y `S2-ahorro-quincenal` abiertos; «no entra a ningún saldo»). Se eligió `null` y no «dato informativo» porque es lo que dice el contrato campo por campo. La pantalla muestra «—», nunca `$0.00`. El supuesto de 5,800 por quincena sigue solo como texto informativo en `proyeccion.supuestos.variable_estimado_por_quincena_cents` y `presupuesto.metodo_variable`; no entra a ninguna fila de esta serie.
- **Saldo.** `saldo_cents` = `ingresos_cents` − `pagos_credito_cents` − `gasto_contado_o_fijo_cents` en todas las filas (igual que `balance_cents` de §3.11). `saldo_acumulado_cents` = Σ de `saldo_cents` desde la primera fila (2026-10-q1), **parte de 0**. `proyeccion.saldo_inicial_cents` se conserva como llave pero vale **0** (el modelo no tiene saldo inicial de cuentas); ninguna pantalla debe mostrarlo como «saldo inicial».
- **Valores nuevos (centavos, saldo / acumulado).** 2026-10-q1: −421772 / −421772 · **2026-10-q2 (en curso): 604550 / 182778 ($6,045.50 / $1,827.78)** · 2026-11-q1: −106434 / 76344 · 2026-11-q2: 998151 / 1074495 · 2026-12-q1: 282659 / 1357154. El saldo de 2026-10-q2 sigue igual a `inicio.tarjetas.flujo_quincenal.closing_balance_cents` (604550).
- **No cambió** ninguna otra llave ni valor del archivo. Las etiquetas `base` («fijos + variable estimado», etc.) conservan su texto porque P-07 las traduce por clave; ya no describen una resta. Los 12 meses de `proyeccion.meses` (devengado contra efectivo, §3.12) no cambian.
