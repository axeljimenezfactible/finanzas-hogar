# REVISION-datos-r3 · Datos sintéticos de las maquetas (ronda 3 de 3)

Revisor: `revisor-datos-r3` (independiente; no escribió ni corrige nada). Fecha: 11-oct-2026, noche.
Objeto: `docs/design/maquetas/datos-sinteticos.json` (252 KB, `meta.version` 1.2), `LEEME.md` 1.2, `docs/design/tools/ds/verificar_sinteticos.py` (57 comprobaciones) y `generar_sinteticos.py` (solo las líneas citadas), reportes `datos-sinteticos.md` y `correccion-datos-r2.md`.
Contra: P3-12 y P3-08 (`docs/criterios/P3-claude-design.md` 1.3), K45.1, U1 §2.1 y §3 (N-044, N-085, N-103, N-114, N-122, N-141, N-171, N-194, N-305), U3 §2 a §4 (6 tarjetas, D-1 a D-7, T-1 a T-4), U4 línea 624.

**Veredicto: NO APROBADO.** 0 altos, **1 medio**, 4 bajos. Los 2 medios de la ronda 2 (M-1 y M-2) están **cerrados**. Hay 1 medio nuevo: el criterio de `seis_mayores` contradice U1 N-044.

## 1. Hallazgos de la ronda 2

| # | Hallazgo r2 | Estado | Evidencia (comandos propios) |
|---|---|---|---|
| M-1 | Resumen del 20-oct con «3 gastos» y 2 en `gastos` | **Cerrado** | `reportes.avisos_enviados[1].texto` = «Resumen del 20-oct-2026: gastaste $951.30 (parte del hogar) en 2 gastos.» Recalculado: 2 gastos realizados ese día (Tianguis de la Narvarte · verdura, Consulta médica de Axel), Σ hogar $951.30. El del 19-oct: 3 gastos, $645.00, cuadra. Prueba negativa propia: copia con «en 3 gastos» → el verificador da `FALLA 34` y `RESULTADO: FALLA`. Las alertas del 80 % cuadran con la fecha de cruce (`mercado` pasa de 76.9 % a 94 % con el gasto del 20-oct, aviso el 21 a las 09:00; `despensa-comun` cruza el 18, aviso el 19). |
| M-2 (1) | P-08 N-103 sin serie | **Cerrado** | `ingresos.tendencia`: 6 meses (2026-05 a 2026-10), total = q15 + q30 + extras en los 6; octubre = `ingresos` (q15 $14,260.00 capturada, q30 $14,780.00 esperada, extras $2,090.00, total $31,130.00) y solo octubre lleva `incluye_esperado`. |
| M-2 (2) | P-02 N-044 y P-06 N-085 por categoría | **Cerrado, con M-1 nuevo** | `gastos_resumen.tendencia_por_categoria`: 6 meses × 13 categorías; Σ plan y Σ real de cada mes = `presupuesto.tendencia`; % ROUND_HALF_UP; octubre = `categorias` (plan y real) = `gastos_resumen.por_categoria.hogar_devengado_cents`; `otras_*` = Σ de las 7 restantes. La serie existe y cuadra; el criterio de las 6 mayores no sigue a U1 (ver M-1 nuevo). |
| M-2 (3) | P-15 N-171 sin historial | **Cerrado** | `productos[].historial_precios`: 105 pares presentación × tienda (los mismos de `precios`), 3 a 5 puntos, fechas ascendentes entre may-2026 y hoy, enteros, último punto = precio vigente (fecha, monto y fuente). 87 fechas distintas, ninguna en `model.json` ni `data.json`. |
| M-2 (4) | P-12 N-141 sin uso | **Cerrado** | `recetas[].uso`: veces = semana actual hasta hoy + 7 previas; la parte actual = casillas (día × tiempo) del plan hasta el 21-oct; próxima fecha = 1.ª fecha del plan desde hoy; ventana del 27-ago al 21-oct (56 días). 12 de 12 cuadran. |
| M-2 (5) | P-26.11 N-305 sin número de registros | **Cerrado** | `ajustes.datos_semilla[].registros`: 33 (= productos), 5 (= 3 compras de Kueski + 2 programadas), 10 (= 8 gastos fijos recurrentes + 2 quincenas del esperado), con `registros_texto` y `registros_estado`. |

## 2. Hallazgos nuevos

### M-1 (media) · `seis_mayores` usa «mayor plan de octubre»; U1 N-044 fija «las 6 categorías con más gasto del rango»
- U1 (`docs/design/U1-arquitectura.md` línea 490), N-044: «suma de la parte del hogar devengada por categoría; **las 6 categorías con más gasto del rango** y «Otras»».
- `generar_sinteticos.py` línea 669: `SEIS_MAYORES = [...] sorted(CAT, key=lambda c: (-c["plan_cents"], c["id"]))[:6]  # las 6 de mayor plan` (plan de octubre). El LEEME lo repite («`seis_mayores` (las de mayor plan)») y `correccion-datos-r2.md` lo declara como supuesto con una cita que U1 no dice («U1 dice «las 6 mayores»»).
- Con los datos del archivo el conjunto cambia. Gasto real (parte del hogar devengada) de mayo a octubre: vivienda $43,200.00, super $18,932.68, transporte $14,515.88, despensa común $11,674.80, mercado $10,742.59, **restaurantes $9,623.35**, **servicios $9,574.00**. Según U1: [vivienda, super, transporte, despensa-comun, mercado, **restaurantes**]. En el archivo: [vivienda, super, **servicios**, transporte, despensa-comun, mercado].
- Efecto: P-02 › Tendencias dibujaría «Servicios» como serie propia y metería «Restaurantes» en «Otras», contra U1. El verificador no lo detecta: con una copia donde `seis_mayores` va por gasto real (y `otras_*` recalculadas) también da `RESULTADO: OK` (comprobación 36 solo revisa que `otras_*` sea el complemento).
- No es una cifra del reporte ni un problema de forma; es una inconsistencia entre el archivo de datos y U1 que cambia una vista (media).

### B-1 (baja) · Tendencias que se arman sumando, sin serie lista
- P-09 › Tendencias, N-114 («por mes: ahorrado al cierre y **aportado contra planificado**»): `ahorro.tendencia` solo trae `ahorrado_cierre_cents`. El aportado por mes sale de sumar `ahorro.movimientos` (comprobado: el cierre de cada mes = Σ acumulada de movimientos; octubre $1,750.00 = `aportado_mes_cents`), pero el planificado de mayo a septiembre no está (solo el de octubre, $2,000.00).
- P-01 › T-3 y P-17 › Tendencias, N-194 («intereses por mes, pagados y por pagar»): `deuda_viva.serie_mensual[].intereses_mes_cents` no separa lo pagado de lo pendiente (octubre mezcla las dos cosas).
- Corrección: sumar `aportado_cents` y `planificado_cents` a `ahorro.tendencia` y separar los intereses, o una regla en el LEEME (como la de 12 meses) que diga cómo dibujarlos.

### B-2 (baja) · La regla de «menos de 12 meses» del LEEME no cita su fuente
El LEEME dice «dibuja los 6 meses anteriores sin cifra, como estado vacío». U1 N-044, N-085 y N-103 dicen «del primer mes con datos al mes en curso (máximo 12)», es decir, 6 meses sin huecos; U3 §4 (ventana fija de 12 meses y subtítulo «Últimos 12 meses») apoya la regla del LEEME. No es inventada, pero U1 y U3 no dicen lo mismo: conviene citar U3 §4 en el LEEME y marcarla como supuesto, para que los autores no la lean como contraria a U1.

### B-3 (baja) · Bajas de la ronda 2 que siguen abiertas
- B-1 r2: campos sin cifras (N-060, N-128, N-274, N-294, N-317, N-011).
- B-2 r2: 18 títulos del calendario con ids internos («Pago tarjeta · tc-estado», «Pago kueski · kc0», «Pago prestamo · pre1»); horizonte disparejo de enero a abril de 2027.
- B-3 r2: `disponible_b_cents` = 900000 ($9,000.00) el 2026-12-15 y el 2027-02-15, fechas que coinciden con `model.json`, donde $9,000.00 también aparece.
- B-4 r2: 105 `precio_unidad_base_cents` flotantes (p. ej. 3.37) sin `_fmt`.
- B-5 r2: «Pendiente de importar» contra «Pendiente» de U4; `person_id` y `store_id` como texto. (La versión ya está alineada: 1.2.)

### B-4 (baja) · Datos semilla: el CSV «pendiente» y lo que ya se ve en otras pantallas
`ajustes.datos_semilla[2]` dice que el CSV de gastos fijos e ingresos está «Pendiente de importar» con 10 registros, mientras P-06 ya muestra 8 gastos fijos recurrentes activos y P-08 el esperado q15/q30 vigente desde el 1-oct (q15 capturada el 15-oct). Se puede leer como captura a mano, pero ningún texto lo dice. Corrección: dejarlo «Cargada», o una nota en el LEEME que diga que esos datos se capturaron a mano.

## 3. Verificaciones repetidas con comandos propios

Scripts en `datos/tmp/revdatos-r3/` (borrados al final).

| Verificación | Resultado |
|---|---|
| `python3 -I docs/design/tools/ds/verificar_sinteticos.py` | 57 OK, `RESULTADO: OK` |
| Regeneración hacia `datos/tmp/revdatos-r3/regen.json` + `cmp` | idéntica (determinista) |
| Prueba negativa: resumen del 20-oct con «3 gastos» | `FALLA 34`, `RESULTADO: FALLA` (la comprobación nueva sí detecta) |
| Prueba negativa: `seis_mayores` por gasto real | `RESULTADO: OK` (el verificador no revisa el criterio; ver M-1) |
| Grep P3-12, las 21 cifras con y sin coma (patrón literal del criterio) sobre el JSON | **0**; sobre todo `docs/design/maquetas/` (incluye el LEEME): **0**. `salida/iter-*` y `salida/final` aún no existen |
| Las mismas 21 cifras en centavos sin separador (309450, 80050, 59882, …, 258685, 1727) | 0 |
| Patrones de credencial de P3-12 en `maquetas/` | 0; 0 correos (`@`) en el JSON |
| P3-08: 848 pares `*_fmt`/`*_cents` | 0 diferencias con `$1,234.56` / `−$`; 0 `$` sin coma de miles, 0 sin dos decimales, 0 guion ASCII antes de `$`; fechas en pantalla `20-oct-2026`, 0 formatos ajenos |
| 22 sumas y porcentajes propios (S1 a S22) | 21 OK, 1 FALLA (S4 = M-1). S1 Σ categorías por mes = `presupuesto.tendencia` (6 meses) y % ROUND_HALF_UP; S2 octubre = `categorias`; S3 `otras_*`; S5 `categorias.real` = `gastos_resumen.por_categoria` ($20,525.78; la diferencia con `suma_hogar` $21,078.00 es el devengo de g13: $1,459.20 − $906.98 = $552.22); S6 ingresos; S7 historial; S8 uso de recetas; S9 datos semilla; S10 resúmenes diarios; S14 tarjeta Mes contra plan ($20,525.78 de $27,586.00, 74 %); S15 ahorro (fondo + metas = $41,730.00 = cierre de octubre; cierres = Σ movimientos); S16 aportado $1,750.00 y planificado $2,000.00; S17 metas 42 % y 78 %; S18 % de `presupuesto.tendencia`; S19 próximo pago 30-oct (2 pagos, $1,861.58, hogar $1,026.25); S20 50 pagos: monto = hogar + Karen; S21 33 gastos: monto = hogar + Karen; S22 solo octubre incluye lo esperado |
| Fechas contra `oraculo/model.json` y `data.json` (comando de P3-12) | 15 coincidencias, todas día 15, 30 o 28-feb (regla de quincena), las mismas que lista el LEEME; 0 fechas en día 9 o 10 desde el 1-oct-2026; 0 coincidencias en las 87 fechas nuevas del historial |
| Copia de `data.json` y `model.json` | 0 nombres de producto idénticos; 5 marcas en común (Cloralex, Dolores, La Moderna, Nescafé, Yoplait) y **ninguno** de sus precios ni puntos de historial coincide con el precio real de ese producto. 67 de 2,473 campos `*_cents` coinciden con algún valor real ≥ $10, siempre valores redondos entre conceptos distintos ($14.00 a $105.00 en precios de verduras y abarrotes; $1,000.00, $2,000.00, $3,000.00, $9,000.00); 18 de ellos en las llaves nuevas (p. ej. $105.00 de yogur 3 × 900 g en Sam's contra $105.00 de 5 latas de atún real); ninguna cifra canon |

## 4. Cobertura: pantalla → llave

| Pantalla (U1 §2.1) | Llaves | Estado |
|---|---|---|
| P-00 | ninguna (`ajustes.usuarios`) | OK (B-3: hora de desbloqueo) |
| P-01 Resumen | `inicio.tarjetas` (6: `mes_vs_plan`, `flujo_quincenal`, `proximo_pago`, `comidas_semana`, `mejor_canasta`, `ahorro`) | OK, U3 tarjetas 1 a 6 |
| P-01 Detalle | `inicio.segundo_nivel` (D-1 `categorias_mayor_uso`, D-3 `proximos_5_pagos`, D-4 `comidas_hoy`, D-5 `articulos_mayor_diferencia`), D-2 `proyeccion`, D-6 `ahorro.metas`, D-7 `participacion_karen` | OK |
| P-01 Tendencias / Calendario | T-1 `presupuesto.tendencia` (+ `tendencia_por_categoria` con filtro), T-2 `proyeccion.meses`, T-3 `deuda_viva.serie_mensual`, T-4 `salud.tendencia_semanal` / `calendario_hogar` | OK (B-1: intereses sin separar; B-2) |
| P-02 / P-03 | `gastos`, `gastos_resumen`, `gastos_resumen.tendencia_por_categoria`, `reparto`, `categorias` | OK; Tendencias: **M-1** (las 6 mayores) |
| P-04 / P-05 | `gastos`, `tickets` | OK (B-3: N-060) |
| P-06 | `presupuesto`, `categorias`, `tendencia_por_categoria` (N-085 por categoría) | OK |
| P-07 | `proyeccion`, `deuda_viva`, `presupuesto.metodo_variable` | OK |
| P-08 | `ingresos`, `ingresos.tendencia` | OK |
| P-09 | `ahorro` | OK; Tendencias N-114: B-1 |
| P-10 / P-11 / P-12 | `plan_semana`, `recetas`, `recetas[].uso` | OK (N-122 y N-141); P-10 Calendario N-128: B-3 |
| P-13 / P-14 / P-16 | `lista_compras`, `inventario` / `canasta` / `inventario` | OK |
| P-15 | `productos`, `productos[].historial_precios`, `tiendas` | OK (B-3: flotantes) |
| P-17 / P-18 / P-19 / P-20 | `creditos`, `promociones` / `comparador` | OK (B-1: N-194) |
| P-21 / P-22 | `salud` | OK (meta de kcal por semana = la de `perfiles`) |
| P-23 / P-24 / P-25 | `reportes`, `reporte_mensual`, `estado_cuenta_karen` | OK (M-1 de r2 cerrado) |
| P-26 (14 secciones) | `ajustes.*` (incluye `notificaciones`, `datos_semilla`, `parametros`), `tiendas`, `categorias`, y las llaves de dominio | OK (.11 con registros; B-4; .01 y .07: B-3) |
| G-03 | `recientes` | OK |

`indice_pantallas` cubre P-00 a P-26 en 18 entradas, sin llaves colgantes (comprobación 33).

## 5. LEEME (claridad para los autores)
Claro y al día: versión 1.2 = `meta.version`, 57 comprobaciones, las llaves nuevas con su pantalla y su N-xxx, la regla de 12 meses ampliada a todas las series, el historial sin interpolar y la nota de textos que repiten cifras. Falta: corregir la definición de `seis_mayores` (M-1), citar la fuente de la regla de 12 meses (B-2) y decir cómo dibujar N-114 y N-194 (B-1).

## 6. Corrección sugerida para M-1 (para quien corrija; este revisor no toca nada)
1. `generar_sinteticos.py` línea 669: ordenar por Σ `real_cents` (parte del hogar devengada) de los 6 meses de `TPC_MESES`, con desempate fijo (orden de P-26.04 o id), y recalcular `otras_*`. Resultado esperado: vivienda, super, transporte, despensa-comun, mercado, restaurantes.
2. Verificador, comprobación 36: recalcular `seis_mayores` con ese criterio.
3. LEEME: «`seis_mayores` (las 6 con más gasto del hogar devengado de mayo a octubre, U1 N-044)». Quitar la cita «U1 dice «las 6 mayores»» del supuesto de `correccion-datos-r2.md` en el reporte de la corrección siguiente.
