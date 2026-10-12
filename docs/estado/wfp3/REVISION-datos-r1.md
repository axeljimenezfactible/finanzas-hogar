# Revisión de datos sintéticos · ronda 1 de 3 · 11-oct-2026, noche

Revisor independiente `revisor-datos-r1`. No escribí nada de lo revisado y no corregí nada. Objeto: `docs/design/maquetas/datos-sinteticos.json` 1.0 (191 KB), `docs/design/maquetas/LEEME.md`, `docs/design/tools/ds/verificar_sinteticos.py` y el reporte `docs/estado/wfp3/datos-sinteticos.md`. Criterios: P3-12 y P3-08 de `docs/criterios/P3-claude-design.md` 1.3, K45.1, U1 §2 y §3 (necesidades por pantalla), U3 §2 (tarjetas) y §4 (tendencias), U4 P-26.09 y P-26.11.

**Veredicto: NO APROBADO** · 0 altos · 4 medios · 4 bajos. Lo de fondo (K45.1 y P3-12) está bien: ninguna cifra real, sumas correctas y fechas limpias. Lo que falla es la cobertura de 3 pantallas o secciones y lo que el LEEME dice del archivo.

## 1. Verificación repetida (mis propios comandos y scripts en `datos/tmp/rev-datos-r1/`, ya borrados)

| # | Qué | Resultado |
|---|---|---|
| a1 | `python3 -I docs/design/tools/ds/verificar_sinteticos.py` | 39 de 39 OK, `RESULTADO: OK`, salida 0 |
| a2 | Grep literal de la Verifica de P3-12 sobre el JSON y sobre `docs/design/maquetas/` | `0` y `0` |
| a3 | Las 19 cifras canon más 8 del reporte y de U3 que no están en el grep ($6,905.50 R-64, $10,000.00, $3,233.00 y $138.50 K-03, $1,197.64 M-01, $5,000.00, $4,000.00 y $199.00 de plan §1.2) en 4 formas: centavos enteros, flotante, texto con coma y texto sin coma | 0 coincidencias de dinero. La única «coincidencia» es `rate_den` = 1000000 (el denominador de la tasa, no es dinero) |
| a4 | Tasas 25.8685 y 17.27 | 0 (las sintéticas son 16.4150 % y 24.3120 %) |
| a5 | Credenciales (patrón de P3-12), `@`, IP, dominios, `http`, «factible», «tailscale» | 0 en todo |
| a6 | 16 sumas y porcentajes recalculados desde las partidas (script propio, no el del autor) | 16 de 16 OK: Σ plan $27,586.00 = resumen = tarjeta; real por categoría desde `gastos` (devengado de octubre, parte del hogar, sin venideros) = $20,525.78; % ROUND_HALF_UP 74 % y estado de las 13 categorías; ingreso del mes $31,130.00 y sin asignar $1,544.00; hogar + Karen = monto en los 33 gastos; flujo de la quincena $15,520.00 − $5,589.55 − $3,884.95 = $6,045.50; próximo pago desde `pagos_programados` (30-oct, 2 pagos, $1,861.58, sin vencidos); ahorro $41,730.00 = fondo + metas = movimientos por destino, meta más próxima «laptop» 42 %; canasta por tienda = Σ detalle, gana Bodega Aurrera $1,551.70 por $86.60 sobre Walmart; lista de compras $1,564.10 = Σ piezas × precio; plan semanal $1,062.20, kcal promedio 2542 y 1587, 2 casillas sin armar; IMC de 16 pesajes; deuda viva hoy $44,606.33 = Σ líneas = hogar + Karen; intereses Kueski = capital × tasa por plazo y pagos = total; 539 pares `_fmt`/`_cents` con la regla `$1,234.56` y `−$`; 0 montos `$NNNN.NN` sin separador y 0 negativos con guion |
| a7 | Comando de fechas de P3-12 (model.json ∪ data.json ∩ JSON) | 15 fechas: 2026-10-30, 2026-11-15, 2026-11-30, 2026-12-15, 2026-12-30, 2027-01-15, 2027-01-30, 2027-02-15, 2027-02-28, 2027-03-15, 2027-03-30, 2027-04-15, 2027-04-30, 2027-05-15, 2027-06-15. Todas en día 15, 30 o fin de febrero; 0 fechas sintéticas en día 9 o 10 desde el 1-oct-2026; 0 días de la semana mal puestos |
| a8 | Revisión de cada fecha coincidente contra los montos que la acompañan (canon y montos reales de esa misma fecha en model/data) | Ninguna va con una cifra canon ni de Kueski. Única coincidencia numérica: `ahorro_planificado_cents` $1,000.00 en las quincenas 2027-01-15, 02-15, 02-28 y 03-15 contra `capital_repuesto` $1,000.00 del modelo: es la aportación fija de ahorro de todas las quincenas, otro concepto y un número redondo. No es copia |
| c1 | Nombres de productos contra `data.json` | 0 nombres iguales. Comparten solo una palabra genérica (aceite, arroz, atún, leche, papel higiénico, yogur, café soluble); los sintéticos no llevan marca ni presentación de la compra real |
| c2 | Los 105 precios sintéticos contra los precios, totales y diferencias de `data.json` | 4 iguales, todos de producto distinto y número redondo: Zanahoria $21.00 (en data es el atún), Manzana roja $47.00 (servilletas), Calabacita $30.00 (catsup), Yogur natural $40.00 (duraznos). No es copia |
| c3 | Todos los `*_cents` ≥ $10 contra todos los números de model/data | 49 coincidencias, todas de números redondos ($1,000.00, $2,000.00, $3,000.00, $9,000.00) o precios chicos de c2. Ver B-4 |
| P3-08 | Lo que P3-08 pide que aplica a los datos | 6 tarjetas en `inicio.tarjetas`; dinero `$1,234.56` (a6); fechas en pantalla `21-oct-2026` / `Mié 21 oct 2026` (`formato`), 0 fechas con `/`; textos en español de México; nombres de persona solo Axel, Ximena y Karen (K45.1, Q2-98 A) |

## 2. Cobertura: pantalla → llave del JSON

| Pantalla (U1 §2.1, §2.3) | Llaves | Estado |
|---|---|---|
| P-00 Iniciar sesión | `ajustes.parametros` (intentos 5, ventana 15 min) | Completo (el índice dice `[]`; alcanza) |
| P-01 Inicio · Resumen | `inicio.tarjetas` (6) | Completo; nombres C1 incompletos (M-3) |
| P-01 · Detalle | `inicio.segundo_nivel`, `participacion_karen` | Completo |
| P-01 · Tendencias T-1..T-4 | `presupuesto.tendencia` (6), `proyeccion.meses` (12), `deuda_viva.serie_mensual` (12), `salud.tendencia_semanal` (12) | T-1 con 6 de 12 meses (B-1) |
| P-01 · Calendario | `calendario_hogar` (58) | **Parcial**: empieza el 15-oct (M-4) |
| P-02 Todos los gastos | `gastos` (33), `gastos_resumen`, `categorias` | Completo |
| P-03 Registrar gasto | `categorias`, `personas`, `cuentas_y_formas_de_pago`, `tiendas` | Completo |
| P-04 Gasto | `gastos` (propiedades, `devengado`, `efectivo`) | **Falta N-061 renglones del ticket** (M-1) |
| P-05 Capturar ticket | — | **Sin datos**: N-071, N-076 (M-1) |
| P-06 Presupuesto del mes | `presupuesto`, `categorias` | Completo |
| P-07 Flujo y proyección | `proyeccion`, `deuda_viva`, `presupuesto.metodo_variable` | Completo |
| P-08 Ingresos | `ingresos` | Completo |
| P-09 Ahorro y metas | `ahorro` | Completo |
| P-10 / P-11 / P-12 | `plan_semana`, `recetas`, `productos` | Completo |
| P-13 Lista de compras | `lista_compras`, `inventario` | Completo |
| P-14 Comparar canasta | `canasta` | Completo |
| P-15 Productos y precios | `productos` (105 precios, 22 viejos), `tiendas` | Completo; precio por unidad sin formato (B-2) |
| P-16 Inventario | `inventario` | Completo |
| P-17 / P-18 / P-19 | `creditos` (4 tipos, serie A/B, tasas, ciclo, amortización, abonos, `pagos_programados`), `promociones` | Completo |
| P-20 Comparador | `comparador` | Completo |
| P-21 / P-22 Salud | `salud` | Completo |
| P-23 / P-24 / P-25 | `reportes`, `reporte_mensual`, `estado_cuenta_karen` | Completo |
| P-26.01 Cuenta y roles | `ajustes.usuarios`, `sesiones`, `parametros` | Completo |
| P-26.02 / .03 / .04 | `presupuesto`, `ingresos.esperado`, `categorias` | Completo |
| P-26.05 Créditos y condiciones | `creditos.lineas` (tasa, base, IVA, pago mínimo), `promociones` | Falta el tiempo máximo del comparador, 5 s (N-321, K50) (M-2) |
| P-26.06 Tiendas y precios | `tiendas`, `canasta.membresia_sams`, `ajustes.profeco`, `parametros` | Completo |
| P-26.07 Recetas | `recetas` (`porciones_base`, `tiempos_aptos`), `productos[].nutricion_fuente` | Suficiente para maquetar |
| P-26.08 Reportes | `ajustes.reportes_programados` (21:30 · dom 19:00 · día 1 08:00, K50), `proximos_envios` | Completo |
| P-26.09 Notificaciones | `ajustes.telegram`, `parametros` (umbrales 80/100, [1, 0]) | **Faltan las 4 alertas; el puntero `notificaciones` no existe** (M-2) |
| P-26.10 Apariencia | `ajustes.tema` | Completo |
| P-26.11 Datos y respaldos | `ajustes.respaldos`, `ia`, `servicio`, `parametros` | Falta N-305 «Datos semilla» (M-2) |
| P-26.12 / .13 / .14 | `personas`, `reparto`, `comparador.tope_pago_quincena`, `salud`, `ahorro`, `inventario` | Completo |

U3, tarjeta → llave: 1 `mes_vs_plan` · 2 `flujo_quincenal` · 3 `proximo_pago` · 4 `comidas_semana` · 5 `mejor_canasta` · 6 `ahorro`. Las 6 traen todos los datos de su contenido, vacío aparte (`estados_alternos`). Las 4 tendencias y D-7 «Participación de Karen» también tienen datos.

## 3. Hallazgos

### M-1 (media · cobertura) · P-04 y P-05 sin datos de ticket
U1 §3.6 N-061 pide para P-04 los renglones del ticket: producto, presentación, cantidad, precio unitario, importe, cuenta, totales por cuenta y la diferencia contra el total. U1 §3.7 N-071 y N-076 piden para P-05 la lectura propuesta (encabezado, texto leído, producto sugerido, cuenta sugerida, nivel de confianza) y los tickets por confirmar, con su borrador. El JSON no tiene ninguna llave de renglones, de lectura ni de confianza. `gastos[]` trae solo `tiene_ticket: true`. El índice manda P-04 y P-05 a `gastos`, así que el autor de esos artboards tendría que inventar los renglones, y eso va contra la regla 1 del LEEME.
**Corrección:** agregar en el generador una llave `tickets` con: (a) un ticket confirmado de un gasto con `tiene_ticket` (por ejemplo `g13`), con 8 a 12 renglones cuya suma dé el monto, la cuenta de cada renglón y las partes hogar y Karen; (b) una lectura OCR pendiente con 6 a 10 renglones propuestos con `texto_leido`, producto sugerido, `confianza`, uno o dos de confianza baja y una diferencia contra el total; (c) dos «tickets por confirmar», uno con borrador. Después, regenerar, sumar la comprobación al verificador y apuntar P-04 y P-05 a `tickets` en el índice y en el LEEME.

### M-2 (media · cobertura e inconsistencia) · Ajustes › Notificaciones incompleta y puntero colgante
`indice_pantallas["P-26"]` incluye `"notificaciones"`, pero esa llave no existe ni en la raíz ni en `ajustes`. U4 P-26.09 pide las 4 alertas de la decisión F (`payment_due`, `budget_threshold`, `stale_prices`, `meal_plan_missing`), cada una con activa, destinatarios, anticipación y hora, y además `advance_days` de precios viejos = 1 (K50). Nada de eso está en el JSON; hay solo los umbrales, los días antes del pago y `avisos_hora`. Faltan también, en P-26.11, N-305 «Datos semilla» (estado de carga) y, en P-26.05, N-321 «Tiempo máximo del comparador» = 5 s (K50).
**Corrección:** agregar `ajustes.notificaciones`: 4 filas con los valores de la decisión F y de K50 (precios viejos 1 día antes, avisos a las 09:00), destinatarios Axel y Ximena y una alerta apagada para el estado. Agregar también `ajustes.datos_semilla` (3 renglones con estado de carga, sin cifras) y `ajustes.parametros.comparador_max_s = 5`. Si no, corregir el puntero del índice.

### M-3 (media · inconsistencia entre archivos) · Las tarjetas de Inicio no usan todas los nombres de C1, aunque el LEEME lo dice
El LEEME dice: «`inicio` · las 6 tarjetas con los nombres de campo de C1 (`DashboardMonthVsPlan`, `DashboardFortnightFlow`, `DashboardNextPayment`, ...)». Contra `docs/contratos/openapi.yaml`:
- `inicio.tarjetas.proximo_pago.next_payment` trae `fecha, linea_id, monto_cents, hogar_cents, karen_cents, dias_que_faltan, num_pagos, vencido`; C1 usa `due_date, credit_id, amount_cents, household_share_cents, external_share_cents, days_left, payment_count, is_overdue`.
- `inicio.tarjetas.ahorro.next_goal` trae `nombre, monto_meta_cents, fecha_meta, ahorrado_cents, avance_pct, aportacion_planificada_cents, aportacion_necesaria_cents, estado: "atrasada"`; `SavingsGoal` usa `name, target_cents, target_date, saved_cents, progress_pct, planned_monthly_cents, required_monthly_cents, status: "behind"`.
- `inicio.tarjetas.mes_vs_plan.mes` debería ser `month`.
`flujo_quincenal`, `comidas_semana`, `mejor_canasta` y `availability[]` sí siguen C1.
**Corrección:** renombrar en el generador esas llaves a las de C1, conservando los `_fmt` hermanos y `etiqueta_estado`, y regenerar. Si se prefiere no tocarlas, el LEEME tiene que decir qué tarjetas siguen C1 y cuáles no.

### M-4 (media · cobertura e inconsistencia) · `calendario_hogar` no cubre octubre y el LEEME da otro rango
El LEEME dice «58 eventos (22-oct a 31-dic)». El archivo va del 15-oct-2026 al 17-abr-2027. U1 N-029 pide los eventos «del rango visible (el mes)». La vista de octubre queda vacía del 1 al 14, aunque en `gastos` hay gastos fijos con esas fechas: renta el 1-oct, telefonía el 6-oct, internet el 12-oct y streaming el 14-oct.
**Corrección:** generar eventos desde el 1-oct-2026 (gastos fijos realizados, la quincena del 15 y las visitas a mercado) y poner en el LEEME el rango real y el número de eventos.

### B-1 (baja · cobertura) · T-1 con 6 de 12 meses
U3 §4 fija la ventana de T-1 en «Últimos 12 meses». `presupuesto.tendencia` trae de mayo a octubre de 2026. Es un estado válido (meses sin cifra), pero el LEEME no le dice al autor que dibuje 6 meses vacíos.
**Corrección:** decirlo en el LEEME o ampliar la serie a 12 meses.

### B-2 (baja · formato) · `precio_unidad_base_cents` es flotante y no tiene regla de pantalla
La regla 2 del LEEME dice «centavos enteros», pero `productos[].precios[].precio_unidad_base_cents` vale 3.37 (centavos por gramo). No hay `_fmt` y, con la regla general, saldría «$0.03». P-15 muestra el «precio por unidad base».
**Corrección:** agregar un `precio_unidad_fmt` («$33.70/kg», «$21.00/pza») o una regla en el LEEME.

### B-3 (baja · estilo) · «P-3-12»
El LEEME y el verificador escriben «P-3-12»; el id es «P3-12».

### B-4 (baja · prudencia) · $9,000.00 de disponible B coincide con la serie real de Kueski
`creditos.kueski.serie_disponible.eventos[10]` y `[15]` (pagos del 15-dic-2026 y del 15-feb-2027) dan lectura B $9,000.00. Sale de la aritmética sintética: límite $12,000.00 menos la compra programada de $3,000.00. `model.json` `credito.eventos[9]` y `[15]` (30-ene y 30-mar-2027) también dan $9,000.00. Las fechas son otras, no está en la lista de P3-12 y no es canon, pero es el mismo concepto.
**Corrección (opcional):** subir la compra programada a $3,150.00 o un monto parecido para que la lectura B no caiga en $9,000.00.

## 4. LEEME para los autores (d)
Es claro y útil: tiene reglas numeradas, el esquema con la pantalla de cada llave, las cifras que deben cuadrar entre pantallas, las fechas inevitables, qué es ilustrativo y cómo regenerar y comprobar. Las cifras que da (mes $20,525.78 de $27,586.00, 74 %, 11 días; próximo pago 30-oct con 2 pagos; B $6,881.60 y A $11,440.00; 39 comprobaciones) coinciden con el JSON. Le faltan las correcciones de M-3 y M-4 y lo que dicen B-1 y B-2.

## 5. Observaciones fuera de los criterios evaluados (no cuentan para el veredicto)
- `docs/design/tools/ds/` (verificador y generador) no se ha agregado a git y está fuera de las 4 rutas permitidas de la Pista 3. El reporte del autor declara el generador como extra; el verificador estaba asignado. Lo decide el orquestador.
- El autor declara 290 KB escritos contra un tope de 150 KB, con el exceso justificado en su reporte.
