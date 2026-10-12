# REVISION-datos-r2 · Datos sintéticos de las maquetas (ronda 2 de 3)

Revisor: `revisor-datos-r2` (independiente; no escribió ni corrige nada). Fecha: 11-oct-2026, noche.
Objeto: `docs/design/maquetas/datos-sinteticos.json` (205 KB), `LEEME.md` 1.1, `docs/design/tools/ds/verificar_sinteticos.py`, reportes `datos-sinteticos.md` y `correccion-datos-r1.md`.
Contra: P3-12 y P3-08 (`docs/criterios/P3-claude-design.md` 1.3), K45.1, K50, U1 §2 y §3, U3 §2 (6 tarjetas), U4 P-26.05/.09/.11, C1 `openapi.yaml` (esquemas `Dashboard*` y `SavingsGoal`).

**Veredicto: NO APROBADO.** 0 altos, 2 medios, 5 bajos. Los 4 hallazgos medios de la ronda 1 quedan **cerrados**. Hay 2 medios nuevos: un texto de Reportes que no cuadra con Gastos y huecos de cobertura en vistas que piden series o conteos.

## 1. Hallazgos de la ronda 1

| # | Hallazgo r1 | Estado | Evidencia |
|---|---|---|---|
| 1 | P-04/P-05 sin renglones, OCR ni tickets por confirmar | **Cerrado** | `tickets.confirmados[0]` (g13): 12 renglones con producto, presentación, cantidad, precio unitario, importe, cuenta y partes; Σ = $2,918.40 = g13; `por_cuenta`. `lectura_pendiente`: 10 renglones con `texto_leido`, sugeridos, `confianza` (6 altas, 2 medias, 2 bajas: 9 y 10), diferencia −$15.90 = total leído − suma. `por_confirmar` 2, por fecha descendente; `ocr2` con borrador ($392.20 = Σ renglones). Índice `P-04/P-05` → `tickets`. |
| 2 | P-26 apuntaba a `notificaciones` inexistente; faltaban 4 alertas, N-305 y N-321 | **Cerrado** | `ajustes.notificaciones.alertas`: `payment_due` [1, 0], `budget_threshold` 80/100, `stale_prices` `advance_days` = 1, `meal_plan_missing` domingo (7), todas con activa, destinatarios y hora 09:00 (K50 punto 73). `ajustes.datos_semilla` 3 renglones. `parametros.comparador_max_s` = 5 y límites [1, 30] (K50 punto 74). Ninguna llave colgante en el índice. Lo que queda (número de registros de N-305 según U4) está en M-2. |
| 3 | Tarjetas de Inicio sin nombres de C1 | **Cerrado** | Contra el YAML de C1: `mes_vs_plan` = `DashboardMonthVsPlan` (7/7 campos), `flujo_quincenal` = `DashboardFortnightFlow` (7/7), `next_payment` (8/8) y `availability[]` (4/4), `comidas_semana` = `DashboardWeekMeals` (5/5), `mejor_canasta` = `DashboardBestBasket` (8/9; `expense_id` solo aplica con `basket` = `expense`), `ahorro` = `DashboardSavings` y `next_goal` con nombres de `SavingsGoal`, `status` `behind`. El LEEME ya lo describe así. |
| 4 | Calendario sin el 1 al 14 de octubre; LEEME con otro rango | **Cerrado** | 74 eventos del 2026-10-01 al 2027-04-17, como dice el LEEME. Octubre completo: renta 1, telefonías 6, internet 12, streaming 14, despensa 12, mercados 3, 6, 13, 17 y 20, pagos ya hechos y reporte del 1. Los 5 fijos y los 4 venideros de `gastos` están en el calendario. Lo que queda está en B-2. |

## 2. Hallazgos nuevos

### M-1 (media) · Un resumen diario de Reportes no cuadra con Gastos (P3-12: «que Inicio, Detalle y Reportes cuadren entre sí»)
- `reportes.avisos_enviados[1].texto` = «Resumen del 20-oct-2026: gastaste $951.30 (parte del hogar) en **3 gastos**». En `gastos` hay **2** gastos con fecha 2026-10-20: «Tianguis de la Narvarte · verdura» y «Consulta médica de Axel». La suma de la parte del hogar, $951.30, sí cuadra; el conteo no. El texto del 19-oct (3 gastos, $645.00) sí cuadra.
- Causa: `docs/design/tools/ds/generar_sinteticos.py` línea 540 calcula el monto con `cg('2026-10-20')`, pero el «en 3 gastos» está fijo en el código. El verificador no revisa los textos de los avisos.
- Efecto: P-23 (avisos enviados) diría 3 gastos, y P-02 y el calendario de ese día mostrarían 2.
- Corrección: calcular el conteo como el monto (líneas 540 y 541) y sumar al verificador una comprobación que recalcule monto y conteo de cada resumen diario.

### M-2 (media) · Cobertura: vistas y campos que piden cifras o conteos que el archivo no trae (U1 §2.1 y §3; K45.1: el autor no inventa)
El índice y el LEEME no dicen qué hacer con estas vistas. La regla de «menos de 12 meses» solo cubre `presupuesto.tendencia` (T-1).
- **P-08 › Tendencias, N-103 (ingreso por mes, hasta 12 meses):** `ingresos` solo trae octubre (esperado, 2 quincenas y 2 extras). No hay serie de mayo a septiembre, aunque `presupuesto.tendencia` y `ahorro.tendencia` empiezan en mayo.
- **P-02 › Tendencias, N-044 (gasto del hogar por categoría y mes, las 6 mayores más «Otras») y P-06 › Tendencias, N-085 «al elegir una categoría»:** `presupuesto.tendencia` solo trae plan y real totales por mes. `gastos_resumen.por_categoria` es solo de octubre.
- **P-15 › Tendencias, N-171 (historial de precio por unidad base, por fecha y tienda):** 105 precios en 105 pares producto × tienda, ninguno con más de una fecha. No hay historial que dibujar.
- **P-12, N-141 «Uso de la receta» (veces en el plan en las últimas 8 semanas y próxima fecha planeada):** las 12 recetas no traen ese dato. Es contenido de la página principal de P-12, así que se dibuja. También falta para el orden de N-122 (veces usadas).
- **P-26.11, N-305:** U4 (`U4-ajustes.md` línea 624) pide para cada uno de los 3 renglones «Cargada» o «Pendiente» **y su número de registros**. `ajustes.datos_semilla` trae el estado y la fecha de carga, pero no el número de registros. El LEEME dice «sin cifras» y no explica por qué; un conteo sintético no es una cifra del reporte.
- Corrección (cualquiera de las dos): (a) sumar las series y campos sintéticos (`ingresos.tendencia`, `gastos_resumen.tendencia_por_categoria`, historial en `productos[].precios`, `recetas[].uso`, `datos_semilla[].registros`) con su comprobación; o (b) escribir en el LEEME una regla explícita, como la de T-1, que diga qué vistas se dibujan como estado vacío con «—» o no se dibujan en esta ronda, y apuntarlas en el índice.

### B-1 (baja) · Campos sin cifras que faltan (se pueden dibujar con «—» o con el texto fijo de U1, sin inventar montos)
- P-04, N-060: los gastos no traen «quién lo creó y cuándo», «quién lo cambió por última vez y cuándo» ni `nota`. Para g13 se puede deducir de `tickets.confirmados[0].confirmado_por`/`confirmado_el`.
- P-10 › Calendario, N-128: no está el estado del plan por día del mes (completo, incompleto o vacío). Solo existe la semana W43.
- P-26.01, N-274 («Este dispositivo»: app instalada, certificado de la CA local); P-26.07, N-294 (sin llave propia en `ajustes`; se deduce de `recetas` y `productos[].nutricion_fuente`); estado global «Sin conexión», N-317 (n en cola), que no está en `estados_alternos`; P-00, N-011 (hora de desbloqueo).

### B-2 (baja) · Calendario: horizonte disparejo y títulos con ids internos
- Octubre a diciembre están completos. De enero a abril de 2027 solo hay las 2 fechas meta, aunque `creditos.pagos_programados` trae 19 pagos en esos meses y los ingresos y gastos fijos siguen. La frase del LEEME «y lo que viene hasta las fechas meta de ahorro» hace pensar que el calendario está completo hasta abril.
- La agenda móvil de 30 días (21-oct a 20-nov, N-029) no trae las visitas a mercado del sábado 14 ni del martes 17-nov. El 10-nov se omite por la regla de fechas.
- Títulos con ids internos: «Pago kueski · kc0», «Pago tarjeta · tc-estado», «Pago prestamo · pre1». También hay «Visita a» y «Visita al» Tianguis de la Narvarte mezclados.

### B-3 (baja, B-4 de r1 sigue abierto) · $9,000.00 en fechas coincidentes
- `creditos.kueski.serie_disponible.eventos[10]` (2026-12-15) y `[15]` (2027-02-15) llevan `disponible_b_cents` = 900000. Las dos fechas coinciden con `model.json`, y $9,000.00 aparece en `model.json` (`credito.eventos[9].disponible_B`, del 2027-01-30, y en las series `esperado_B` y `calculado_B`).
- No es copia: el par fecha-monto es distinto y el valor sale de la aritmética sintética (límite $12,000 − capital $3,000). Aun así, la frase del LEEME «Ninguna va con un monto del reporte ni de Kueski» no es literal, y la comprobación 50 del verificador no compara montos en las fechas coincidentes.

### B-4 (baja, B-2 de r1 sigue abierto) · `precio_unidad_base_cents` es flotante (p. ej. 3.37) y no tiene `_fmt`
Va contra la regla 2 del LEEME («dinero en centavos enteros»). El autor de P-15 tiene que decidir cómo mostrar el precio por unidad base ($/kg o $/L).

### B-5 (baja) · Detalles de consistencia
- `meta.version` = «1.0» y el LEEME dice «Versión 1.1».
- `person_id` («axel») y `store_id` («bodega-aurrera») son texto, y C1 los define como enteros. El LEEME solo aclara el caso de `credit_id`.
- `datos_semilla` dice «Pendiente de importar» y U4 dice «Pendiente».

## 3. Verificaciones repetidas con comandos propios

| Verificación | Resultado |
|---|---|
| `python3 -I docs/design/tools/ds/verificar_sinteticos.py` | 51 OK, `RESULTADO: OK` |
| Regeneración determinista hacia `datos/tmp/revdatos-r2/regen.json` + `cmp` | idéntica |
| Grep P3-12 (21 cifras, con y sin coma) sobre `docs/design/maquetas/` | **0** (`salida/iter-*` y `salida/final` aún no existen) |
| Patrones de credencial de P3-12 | 0; ningún correo en el JSON |
| P3-08: `\$[0-9]{4,}\.[0-9]{2}` en maquetas | 0; cada `*_fmt` coincide con su `*_cents` hermano (`$1,234.56`, `−$`; 1,851 campos `*_cents` recorridos); 5 textos libres con dinero, todos bien formados; fechas en pantalla `20-oct-2026` |
| Cifras canon en centavos (comparación de todos los `*_cents` con los 215 valores ≥ $10 de `data.json` y `model.json`) | 0 cifras canon; 14 coincidencias de valores redondos entre conceptos distintos (p. ej. $21.00 zanahoria contra $21.00 atún real; $2,000.00 de plan de categoría contra 2000 de `interes_por_monto`) |
| 26 sumas y porcentajes propios (`chk.py`) | 26 OK: Σplan y Σreal de las categorías = resumen; % ROUND_HALF_UP, estado y diferencia por categoría; tarjeta Mes contra plan (11 días contando hoy); ingreso del mes y «sin asignar»; flujo de la quincena ($15,520.00 = $14,780.00 + $740.00; saldo = ingresos − hechos − comprometido); hogar + Karen en 33 gastos y en todos los pagos; `gastos_resumen` (29 realizados, $24,413.20); ticket g13; lectura OCR; borrador; próximo pago (30-oct, 2 pagos de Kueski, 9 días); ahorro (fondo + metas = $41,730.00; % 42 y 78); movimientos por destino; canasta por tienda; mejor canasta (diferencia $86.60); plan semanal ($1,062.20); lista de compras; deuda viva; compras Kueski; IMC. (S9 falló en la primera corrida por el redondeo al par de Python; con ROUND_HALF_UP cuadra: 0.85 × 3370 = 2864.5 → 2865.) |
| Fechas contra `oraculo/model.json` y `data.json` (comando de P3-12) | 15 coincidencias, todas día 15, 30 o 28-feb (regla de quincena), las mismas que lista el LEEME; la única con un monto que también existe en `model.json` es B-3 |
| Copia de `data.json` | 0 nombres de producto idénticos (42 reales contra 33 sintéticos); 5 marcas en común (Cloralex, Dolores, La Moderna, Nescafé, Yoplait), con presentación y precio distintos (p. ej. Yoplait real 750 g a $75.00; sintético 900 g a $39.80/$40.00) |
| Conteos que afirma el LEEME | 33 productos, 105 precios, 22 viejos (≥ 30 días); 14 de inventario; 19 movimientos; 3 promociones (1 vigente); 4 líneas; 6 escenarios; 12 recetas; 33 gastos (5 fijos, 24 variables, 4 venideros); 74 eventos: todo coincide |
| Cifras «que deben coincidir» del LEEME (más la deuda de Karen) | $20,525.78 de $27,586.00, 74 %, 11 días; vie 30-oct con 2 pagos; B $6,881.60 y A $11,440.00; deuda de Karen $3,802.33 = `deuda_viva.a_fecha` de hoy; `reporte_mensual.resumen` = presupuesto: todo coincide |
| Avisos enviados contra datos | alertas 80 % (Despensa común $1,650.98 de $2,000.00 = 83 %; Frutas $1,685.50 de $1,800.00 = 94 %) y 100 % (Vivienda) cuadran con `categorias`; resumen del 19-oct cuadra; **resumen del 20-oct no** (M-1) |

## 4. Cobertura: pantalla → llave

| Pantalla | Llaves | Estado |
|---|---|---|
| P-00 | ninguna (`ajustes.usuarios`) | OK (B-1: hora de desbloqueo) |
| P-01 Resumen | `inicio.tarjetas` (6) | OK, U3 tarjetas 1 a 6 con nombres de C1; variante `estados_alternos.pago_vencido` |
| P-01 Detalle / Tendencias / Calendario | `inicio.segundo_nivel`, `participacion_karen` / `presupuesto.tendencia`, `proyeccion.meses`, `deuda_viva.serie_mensual`, `salud.tendencia_semanal` / `calendario_hogar` | OK (T-1 con regla de 6 meses; B-2) |
| P-02 / P-03 | `gastos`, `gastos_resumen`, `reparto`, `categorias`, `cuentas_y_formas_de_pago` | OK; Tendencias por categoría: M-2 |
| P-04 / P-05 | `gastos`, `tickets` | OK (N-061, N-071, N-075, N-076); N-060: B-1 |
| P-06 / P-07 | `presupuesto`, `categorias` / `proyeccion`, `deuda_viva`, `presupuesto.metodo_variable` | OK; P-06 por categoría en Tendencias: M-2 |
| P-08 | `ingresos` | Detalle y Calendario OK; Tendencias: M-2 |
| P-09 | `ahorro` | OK |
| P-10 / P-11 / P-12 | `plan_semana`, `recetas` | OK; P-10 Calendario: B-1; P-12 N-141: M-2 |
| P-13 / P-14 / P-16 | `lista_compras`, `inventario` / `canasta` / `inventario` | OK |
| P-15 | `productos`, `tiendas` | Detalle OK; Tendencias: M-2; B-4 |
| P-17 / P-18 / P-19 / P-20 | `creditos`, `promociones` / `comparador` | OK |
| P-21 / P-22 | `salud` | OK |
| P-23 / P-24 / P-25 | `reportes`, `reporte_mensual`, `estado_cuenta_karen` | OK salvo M-1 |
| P-26 (14 secciones) | `ajustes.*`, `presupuesto`, `ingresos.esperado`, `categorias`, `creditos`, `promociones`, `comparador.restricciones`, `tiendas`, `personas`, `salud`, `ahorro.metas`, `inventario` | .02 a .06, .08 a .10 y .12 a .14 OK; .01 y .07: B-1; .11 sin número de registros: M-2 |
| G-03 | `recientes` | OK |

## 5. LEEME (claridad para los autores)
Es claro: 8 reglas, esquema por llave con pantallas, cifras que deben coincidir, fechas inevitables, qué es ilustrativo y cómo regenerar y comprobar. Lo que falta está en M-2 (una regla para las vistas sin serie) y en B-2, B-3 y B-5 (rango del calendario, una frase sobre las fechas que no es literal y la versión).
