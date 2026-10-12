# U3 · Dashboard: tablero de control (P-01 Inicio) · Finanzas Hogar

Versión 2.0 (12-oct-2026) · Agente `u3-v2` (WF-3c, decisión K54 de Axel, 12-oct-2026): **«Tablero de control»**. La pestaña «Resumen» deja las 6 tarjetas y pasa a ser un tablero: el gasto real del mes como indicador principal (K-1), ingreso del mes con su extra (K-2), créditos globales (K-3), pagos próximos (K-4), el barómetro de capacidad de endeudamiento (B-1, DESIGN §8.9) y seis gráficas de análisis (G-1 a G-6), con el desglose en cadena hasta el historial del precio por unidad base (§2.2) y la tabla «Qué miro para prevenir» (§2.3). Detalle (§3), Tendencias (§4) y Calendario (§5) se conservan; G-4 y G-5 son T-2 y T-3. Las tarjetas 2, 4, 5 y 6 siguen sin cambio en «Más del hogar» (§2.5, PROVISIONAL); la 1 pasa a K-1 y la 3 a K-3 y K-4.

**Cambios de la versión 2.0:**
- §2 nuevo (sustituye «Resumen: 6 tarjetas»): K-1 a K-4, B-1, §2.1 y G-1 a G-6, §2.2 desglose en cadena, §2.3 «Qué miro para prevenir», §2.4 orden en 412 y rejilla en 1440, §2.5 «Más del hogar», §2.6 operaciones de C1 (3 nuevas, un par de Ajustes, 14 campos nuevos en 3 esquemas y 5 claves de `app_setting`) para C1 v1.5.
- PROVISIONAL con bloqueo para la puerta 4: ingreso neto, pagos que cuentan, redondeo y escala del barómetro (`docs/estado/bloqueos/U3-tablero-barometro.md`); N y rango de K-4 (`U3-tablero-recordatorio.md`); horizonte de las gráficas (`U3-tablero-horizonte.md`); regla «subir o recortar presupuesto» (`U3-tablero-senal-presupuesto.md`); tarjetas que se conservan (`U3-tablero-mas-del-hogar.md`).
- §1, §3, §6 a §15: las referencias a las tarjetas 1 y 3 pasan a K-1, K-3 y K-4; filtros, estados, tamaños, Karen, accesibilidad, criterios A-xx, capturas, trazabilidad, supuestos y dependencias al día. Ninguna cifra canon cambia; las cifras nuevas de los ejemplos salen de `oraculo/model.json` (`periodos[0..2]`) o son marcadores.

- Corrección r1 (12-oct-2026, agente `correccion-especificacion-r1`, hallazgo M-1 de `docs/estado/wf3/REVISION-especificacion-r1.md`): §2.6 queda alineada con C1 v1.5 tal como se publicó (ventana inválida de `getDashboardMonthlyOverview`: 422 `validation_failed`; el PATCH de `updateDashboardSettings` lleva `parameters` y `reset_parameters`); las operaciones, los 14 campos y M-18 a M-22 pasan de PROPUESTA a publicados en C1 v1.5; K-4 lee N de `getDashboardSettings` una sola vez (§1.2) y la interfaz no fija N; «Fuentes que rigen» cita C1 v1.5. Ninguna cifra cambia.
- Corrección r2 (12-oct-2026, agente `correccion-especificacion-r2`, hallazgo M-N1 de `docs/estado/wf3/REVISION-especificacion-r2.md`): K-1 cita la fila M-01 de C1-datos §9, que ya define los 8 campos nuevos de `DashboardMonthVsPlan` con su fórmula, periodo, tablas y ejemplo del oráculo (598.82 de efectivo en oct-2026; en nov-2026, 0.00 de diferencia en lo gastado y 598.82 más en efectivo). Las cifras y la fórmula de K-1 no cambian.

Versión 1.5 (11-oct-2026) · Agente `correccion-U4-r2` (P2c, hallazgos M-1 y M-2 de `docs/estado/p2c/REVISION-U4-r2.md`): (M-1) en D-2 el ahorro planificado deja de ser FIRME: el reparto mitad y mitad sí lo es (Q2-13 B), pero el valor es **PROVISIONAL** por el bloqueo abierto `docs/estado/bloqueos/S2-ahorro-quincenal.md`; mientras llegue `null`, la fila muestra «—»; §15 lo pasa a «Abiertos» y lo quita de los resueltos. (M-2) D-6 muestra el avance con un decimal (`progress_tenths_pct`), porque Q2-17 C deja el entero solo en la tarjeta de Inicio y D-6 es la métrica M-16 «Metas»; §15 pide a C1-datos alinear M-16. Ninguna cifra de canon cambia. Versión 1.4 (11-oct-2026) · Agente `correccion-U4-r1` (P2c, hallazgo M-4 de `docs/estado/p2c/REVISION-U4-r1.md`): se quitan las marcas PROVISIONAL de bloqueos que la puerta 2 dejó RESUELTOS, `S2-presupuesto-total` (Q2-19 B y Q2-20 A), `S2-metas` (Q2-16 A, Q2-17 C y Q2-18 A), `S2-proyeccion-interruptor` (Q2-21 a Q2-24 C) y Q2-13 B y Q2-14 A de `S2-flujo-quincenal`; siguen PROVISIONAL, citando sus bloqueos abiertos, el gasto variable estimado (`S2-flujo-quincenal`, Q2-12) y la aportación de una meta vencida sin cumplir (`S2-metas-vencida`). La tarjeta 6 y D-6 usan los tres estados de Q2-18 A. Versión 1.3 (11-oct-2026) · Agente `U4-puerta2` (P2c, puerta 2): Q2-65 A (T-2 con los acumulados de C1, una sola variante; bloqueo `docs/estado/bloqueos/U3-t2-acumulados.md` RESUELTO) y Q2-86 A (la tarjeta 5 nombra la segunda tienda con `second_store_id`, solicitud a C1 aprobada); ninguna cifra cambia. Versión 1.2 (11-oct-2026) · Agente U3-dashboard (WF-1); ronda 1 de corrección: agente `correccion-U3-r1`; ronda 2: agente `correccion-U3-r2` · Evalúa: `docs/criterios/U3.md` v1.2.

**Cambios de la versión 1.2** (hallazgo alto H-1 y medios H-2 y H-3 de `docs/estado/wf1/REVISION-U3-r2.md`):
- H-1: en T-2 las dos capas tienen el mismo origen que `cumulative_difference_cents` (C1 M-11): capa 1 = efectivo acumulado, cima = devengado acumulado. C1 v1.2 no trae esos dos acumulados; U3 los pide (§15) y deja las capas, sus etiquetas finales y su resumen **PROVISIONAL** con dos variantes (`docs/estado/bloqueos/U3-t2-acumulados.md`). Queda prohibido sumar `cash_cents` desde el primer mes visible y apilar encima `cumulative_difference_cents`. A-11 suma una corrida con reloj en nov-2026.
- H-2: los ejemplos de oct-2026 a mar-2027 son los periodos P1 a P6 del oráculo (C1-datos M-01; `is_canonical`; `S2.md` S2-09, nota de equivalencia). Se retira la frase de T-2 que lo negaba; R-51 y R-52 comprueban las capas en mar-2027.
- H-3: cada gráfica lleva su resumen de una frase y su `<figcaption>` (DESIGN §7.7), con marcadores y solo cifras de C1; sin variaciones calculadas y, en T-2, sin rotular la deuda viva.

**Cambios de la versión 1.1** (hallazgos medios H-1 a H-7 de `docs/estado/wf1/REVISION-U3-r1.md`):
- H-1: la tarjeta sigue DESIGN §7.6 (etiqueta `font.size.small` en `ui.text.secondary`, relleno `space.6` desde 768 px); las gráficas cargan sin ejes falsos (§7.7); nuevo §1.4 con los estados de pago de DESIGN §7.8, incluidos «Por vencer» y «Parcial» (K15); el destino de cada tarjeta se mantiene y se pide a DESIGN (§15).
- H-2: §4 adopta DESIGN §8.1, §8.3, §8.5 y §8.7 (sin contorno en T-1; en T-2, capa 2 en cero con «Pagado por adelantado» y línea «Hoy»; etiquetado selectivo; Ximena con cuadrado y trazo 8 · 4); se retira el supuesto 5.
- H-3 y H-4: D-3 hace dos peticiones explícitas (vencidos, y de hoy a hoy + 60) y su ejemplo trae los 3 pagos que caben en el rango.
- H-5: la tabla de T-2 entra en §9 y en A-13.
- H-6: la fila «Total» de T-2 suma las 12 filas visibles, sin cifra de canon.
- H-7: D-4 y D-5 escriben su fórmula, periodo y filtros mientras C1-datos §9 no tenga su métrica (solicitud en §15).

**Fuentes que rigen** (no se reabren): decisiones §A, §B, §C, §G y §K (K5, K8, K38, K41, K43, K44 y, desde 2.0, K54); U1 1.2 §0.2 (vocabulario, semana por defecto K38), §0.4 (formatos), §0.5 (personas), §1.3 (móvil), §2.4 (estados), §3.3 (P-01: N-020 a N-029), §3.9 (P-07), §6.1 y §6.2 (puntos de quiebre); contrato C1 v1.5 por sus extractos `inicio.yaml`, `presupuesto-a.yaml`, `presupuesto-b.yaml` (y, para los paneles que reutilizan necesidades de otras pantallas, `comidas.yaml`, `compras-tiendas-a.yaml`, `creditos-a.yaml`, `creditos-b.yaml`, `ajustes-b.yaml`, `reportes.yaml`); métricas M-01 a M-17 de `C1-datos.md` §9; DESIGN.md §1, §2.4, §2.6 a §2.10, §3, §4, §5.2 a §5.4, §7.6 (tarjetas de métrica), §7.7 (gráficas), §7.8 (estados de pago) y §8 (patrones de visualización). Donde U3 se aparta de DESIGN lo dice en el lugar, con su razón, y lo pide a DESIGN en §15.

**Reglas que este documento no repite en cada bloque:**
- La interfaz **no recalcula reglas de negocio**: muestra, filtra, ordena y suma lo que está a la vista (U1 §0.2, «Motor Sx»). Toda cifra viene de la operación C1 citada.
- **Sin datos no es cero**: si falta una entrada se muestra «—» con el motivo, nunca un monto en cero ni `0 %` (U1 §0.2).
- **Toda cifra de presupuesto, flujo y ahorro es la parte del hogar**; la parte de Karen es informativa y solo aparece donde ella participa (U1 §0.2, §0.5; decisiones §A).
- Dinero con el formato de U1 §0.4 y DESIGN §3.3: MXN, símbolo `$`, coma de miles, dos decimales, cifras tabulares (`tnum`), negativo con el signo menos tipográfico, color de rol **e** icono (nunca solo color). Fechas en pantalla `10-oct-2026`; con día de la semana donde importa (`Vie 30 oct 2026`). En la API las fechas van en ISO; en este documento solo aparecen en formato de pantalla.
- Colores, tipografía, espacios, radios y sombras **solo por nombre de token** de C2 / DESIGN (0 valores sueltos).
- Ejemplos con cifras: solo las del oráculo o la semilla con su referencia (R-nn de `S2.md`/`S3.md`, K-nn de `S5.md`); lo demás va como marcador entre corchetes, sin número: `[plan]`, `[n]`. Los meses oct-2026 a mar-2027 de los ejemplos son los periodos P1 a P6 del oráculo: C1-datos M-01 dice «P1 = oct-2026», `is_canonical` de C1 marca «oct-2026 a mar-2027» y `docs/criterios/S2.md` S2-09 (nota de equivalencia) explica que con la semilla los dos dan las mismas cifras, porque ningún evento cae entre los días 1 y 9 de un mes. Por eso «octubre de 2026» cita `periodos[0]` en K-1, D-7, T-2, A-06 y A-11 (y `periodos[1]` y `periodos[2]`, noviembre y diciembre, en K-1, B-1, G-1 y G-6). Una prueba con eventos entre los días 1 y 9 toma su esperado de las ventanas que declara, nunca de R-45 a R-50 (S2-09).

---

## 1. Estructura de Inicio

### 1.1 Pestañas y página
- **Título de página**: «Inicio» (`font.size.title-1`, `font.weight.bold`), icono `house`.
- **Pestañas** (orden fijo, decisiones §G): «Resumen» · «Detalle» · «Tendencias» · «Calendario». Abre en «Resumen», que desde U3 2.0 es el tablero de control de §2 (K54; la etiqueta de la pestaña no cambia, decisiones §G). La pestaña elegida se recuerda en el dispositivo (frecuencia «Local», U1 §0.3). En < 768 px la fila de pestañas se desplaza en horizontal y la activa queda siempre visible (U1 §1.3); en ≥ 768 px es una fila completa.
- **Acción primaria de la página**: «Registrar gasto» (U1 §6.4, fila P-01). Es el **único** botón con relleno de marca de la página (DESIGN P2, §2.9): en < 768 px es el «+» de la barra inferior (icono `plus`); en ≥ 768 px, el botón «+ Registrar gasto» de la barra superior. Abre P-03 (flujo F-1, 3 toques, K10).
- Las tarjetas **no** llevan botón con relleno: su acción primaria única es un **enlace** (texto subrayado con `color.brand.700` / `color.brand.300` e icono `chevron-right`) y toda la tarjeta es zona de toque que lleva al mismo destino.

### 1.2 Anatomía común de una tarjeta de métrica
Componente «Tarjeta de métrica» (DESIGN §7.6). Partes, de arriba abajo:

| Parte | Contenido | Tokens |
|---|---|---|
| Etiqueta (título) | Icono de 20 px (`size.icon.sm`, DESIGN §5.2; el mapa de §5.4 prevé iconos en tarjetas) + título en español, sin dos puntos | etiqueta `font.size.small` en `ui.text.secondary` (DESIGN §7.6); el icono toma el color del texto que acompaña (DESIGN §5.3) |
| Valor principal | Una cifra o texto corto; nunca se trunca ni se parte (U1 §6.1 regla 5) | `font.size.title-2`, `font.weight.semibold`, `tnum`; neutro `color.neutral.900` / `color.neutral.100` salvo cifra de rol (DESIGN §2.9) |
| Valor secundario | Una línea a 360 px; una línea más desde 412 px (U1 §6.2) | `font.size.small`, `color.neutral.700` / `color.neutral.300` |
| Insignias de estado | Icono + texto (nunca solo color) | insignia de rol: texto `color.role.<rol>.light` sobre `color.role.<rol>.surface-light` (oscuro: `.dark` sobre `.surface-dark`) |
| Pie | Enlace de la acción primaria | `font.size.small`, `font.weight.medium`, alto de zona `size.touch.min` |

- Superficie: `color.neutral.100` sobre fondo `color.neutral.50` (oscuro: `color.neutral.900` sobre `color.neutral.950`), `shadow.0`, borde `border.width.default`, radio `radius.md`, relleno interior `space.4` (< 768 px) y `space.6` (≥ 768 px), como DESIGN §7.6. Separación entre tarjetas: el medianil de DESIGN §4.2, `space.4` (< 1024 px) y `space.6` (≥ 1024 px).
- **Destino**: la tarjeta entera es un enlace a la pantalla destino de su acción primaria (K-1 → P-06, K-2 → P-08, K-3 y B-1 → P-17, K-4 → pestaña Calendario con cada fila a P-18; en «Más del hogar», P-07, P-10, P-14 y P-09; §2), no a la pestaña Detalle; Detalle se abre con su pestaña. U3 se aparta aquí de DESIGN §7.6 («la tarjeta entera es enlace al Detalle») porque U3-01 pide una acción primaria única por tarjeta y exige que la de la tarjeta 4 lleve a P-10; solicitud 1 a DESIGN (§15).
- Toda la tarjeta mide al menos `size.touch.min` de alto y es un solo objetivo táctil con foco visible (`border.width.focus`, anillo `color.brand.600` / `color.brand.400`). Nombre accesible = título + valor principal + estado.
- Cada tarjeta hace **su propia petición** a C1: si una falla, las otras se ven normales.
- Ninguna tarjeta usa los umbrales 80 % y 100 % como constantes: el estado ya viene calculado en la respuesta (`status` de M-01). Si una implementación necesitara el umbral, lo lee de `getAlertSettings` (C1-22, K44).
- Los días antes del pago (`alerts.payment_due_days_before`, inicial 1 y 0, decisión F y K44) sí se usan, para el estado «Por vencer» de K-4, D-3 y el Calendario (§1.4). Inicio los lee **una sola vez** con `getAlertSettings` (GET `/api/v1/settings/alerts`, campo `parameters.alerts.payment_due_days_before.value`) y nunca como constante; si esa lectura falla, ningún pago se muestra «Por vencer» (queda «Pendiente») y nada más cambia.
- N de «Pagos próximos» (K-4, K54) tampoco es constante: Inicio lo lee **una sola vez** con `getDashboardSettings` (GET `/api/v1/settings/dashboard`, C1 v1.5, campo `parameters["dashboard.payment_reminder_count"].value`, valor inicial 5), igual que `getAlertSettings`. Si esa lectura falla, K-4 muestra su estado de error (§7) y los demás bloques siguen; el valor 5 y esta conducta son PROVISIONALES (bloqueo `docs/estado/bloqueos/U3-tablero-recordatorio.md`).

### 1.3 Estados comunes de las tarjetas (U1 §2.4)
- **Carga**: esqueleto estático con la forma final (etiqueta, bloque de cifra y una o dos líneas) en `ui.skeleton.on-card` (DESIGN §7.6; movimiento de DESIGN §6.2); nunca cifras provisionales. Si tarda más de 10 s: «Esto está tardando más de lo normal · Reintentar».
- **Error**: mensaje de la tarjeta (texto exacto en cada una) + botón secundario «Reintentar» (borde `color.neutral.600` / `color.neutral.400`), icono `circle-alert` con `color.role.error.light` / `color.role.error.dark`. La tarjeta conserva su título y su enlace.
- **Sin conexión** (G-07, K9): franja global «Sin conexión con el servidor de casa» (icono `wifi-off`; «· [n] en cola» si hay capturas guardadas en el teléfono). Cada tarjeta muestra lo último cargado con la línea «Cargado a las [HH:MM]» (`font.size.caption`) y su enlace sigue abriendo la pantalla destino en solo lectura. Si la tarjeta nunca se cargó en el teléfono: «Sin conexión · Esta página aún no se ha cargado en este teléfono». «Registrar gasto» sigue disponible y deja el gasto en la cola (U1 §2.4).

### 1.4 Estados de un pago (DESIGN §7.8)
K-4, D-3 y el Calendario usan los estados de pago de DESIGN §7.8, siempre con icono y texto (nunca solo color). Un pago muestra un solo estado: el primero de esta tabla que se cumpla.

| Estado | Cuándo (dato de C1) | Icono y color | Texto |
|---|---|---|---|
| Pagado | `status` = `paid` (todas las partes pagadas) | `circle-check`, `color.role.success.*` | «Pagado» |
| Vencido | `status` = `overdue` (en `getDashboardNextPayment`, `is_overdue` = true) | `calendar-x`, `color.role.warning.*` | «Vencido» y la línea «Sin marcar como pagado» (U1 N-022, C1 M-03) |
| Por vencer | pendiente, y los días que faltan para el pago son uno de los de `alerts.payment_due_days_before` (§1.2) | `calendar-clock`, `color.role.warning.*` | «Vence mañana» (1), «Vence hoy» (0); otro día configurado en Ajustes: «Vence en [d] días» |
| Parcial | pendiente, con una parte de `shares[]` pagada por completo (`paid_cents` = `amount_cents`) y otra no; el capital se repone solo cuando el pago se completa (K15) | `circle-dot`, `ui.icon.default` | «Falta la parte de Karen» (o «Falta la parte del hogar» si la que falta es la del hogar) |
| Pendiente | cualquier otro pago con alguna parte sin pagar | `circle-dashed`, `ui.icon.default` | «Pendiente» |

- Un pago «Vencido» o «Por vencer» que ya tiene una parte pagada agrega en su línea secundaria «Falta la parte de Karen» (o «del hogar»).
- Qué estados puede mostrar cada bloque: K-4 y D-3 traen `shares[]` y solo pagos con alguna parte sin pagar, así que nunca muestran «Pagado» (sí «Parcial»); el Calendario muestra los cinco (§5).
- «Sin marcar como pagado» acompaña a «Vencido» porque U1 N-022 y C1 M-03 piden esa marca y DESIGN §7.8 solo pone «Vencido». El orden de prioridad y los textos «Falta la parte del hogar» y «Vence en [d] días» tampoco están en DESIGN (supuestos 6 y 7; solicitud 2 a DESIGN, §15).

---

## 2. Resumen: tablero de control (K54)

**Por qué cambia** (decisiones K54, 12-oct-2026): Axel pidió «un dashboard de un sistema de control completo» con ingreso mensual y extra, créditos globales, recordatorio de pagos, un barómetro de capacidad de endeudamiento «con el digito de cuanto se lleva» y «graficas de analisis». Y fijó el propósito del sistema: «saber cuanto se esta gastando realmente mes con mes y ver los desgloces a detalle incluso de cuando costo el kg de cebollas, para ver si aumentamos presupuesto o lo recortamos». Por eso la pestaña «Resumen» deja las 6 tarjetas de U3 1.x y pasa a ser un **tablero de control**. Se retira la regla «≤ 6 tarjetas» (K54 decisión 3). Las pestañas, su orden y la acción primaria de página no cambian (§1.1).

**Indicador principal: el gasto real del mes** (K-1), devengado y en efectivo, contra el plan y contra el mes anterior, con desglose en cadena hasta el precio por unidad base de un producto (§2.2) y la sugerencia «subir o recortar presupuesto» por categoría (G-3).

**Bloques del tablero** (ids fijos; el orden de apilado y la rejilla están en §2.4):

| Id | Bloque | Pregunta que contesta | Operación C1 |
|---|---|---|---|
| K-1 | Gasto del mes (principal) | ¿Cuánto gastamos realmente y cómo vamos contra el plan y contra el mes pasado? | `getDashboardMonthVsPlan` (M-01) + 8 campos nuevos de C1 v1.5 |
| K-2 | Ingreso del mes | ¿Cuánto entró este mes y hubo extra? | `getIncomeMonth` (N-100) + 4 campos nuevos de C1 v1.5 |
| K-3 | Créditos | ¿A cuánto ascienden los créditos y cuánto toca pagar este mes? | `getCreditTotals` (M-12 a hoy) + 2 campos nuevos; `getDashboardNextPayment` (M-04) |
| K-4 | Pagos próximos | ¿Qué pagos vienen y cuál está por vencer o vencido? | `listScheduledPayments` (M-13) |
| B-1 | Barómetro de capacidad de endeudamiento | ¿Qué parte del ingreso se va en pagos de crédito? | **NUEVA** `getDashboardDebtCapacity` |
| G-1 | Ingreso, gasto y pagos por mes | ¿Entra más de lo que gastamos y pagamos? | **NUEVA** `getDashboardMonthlyOverview` |
| G-2 | Deuda viva por línea | ¿En qué línea está la deuda? | `getLiveDebt` (M-12) |
| G-3 | Presupuesto por categoría | ¿Qué categoría se pasa y cuál conviene subir o recortar? | `getBudgetByCategory` (M-09) + **NUEVA** `getBudgetSignals` |
| G-4 | Devengado contra efectivo | ¿Lo que se consume se está pagando? | `getAccrualVsCash` (M-11) = T-2 de §4 |
| G-5 | Deuda viva proyectada | ¿Cuándo baja la deuda? | `getDebtTrend` (M-17) = T-3 de §4 |
| G-6 | Tendencia del barómetro | ¿La capacidad va a empeorar en los próximos meses? | **NUEVA** `getDashboardMonthlyOverview` (la misma petición de G-1) |
| 2.5 | Más del hogar | Quincena, comidas, canasta y ahorro (tarjetas 2, 4, 5 y 6 de U3 1.x sin cambio) | las de cada tarjeta; **PROVISIONAL** |

**Reglas del tablero** (además de las del inicio del documento y de §1):
- Cada bloque hace **su propia petición** y tiene sus estados (§1.3, §7): si uno falla, los demás se ven normales. G-1 y G-6 comparten una petición; si falla, los dos muestran su error.
- La interfaz **no calcula** porcentajes, diferencias ni estados: llegan de C1 (motores S2 y S3). Solo cuenta, ordena y suma lo que está a la vista (U1 §0.2).
- Umbrales: los del presupuesto (80 % y 100 %) se leen una vez de `getAlertSettings` (§1.2). Los del barómetro (30 % y 40 %) llegan en la respuesta de `getDashboardDebtCapacity` y de `getDashboardMonthlyOverview` (`warning_pct`, `limit_pct`, de las claves de `app_setting` `debt.capacity_warning_pct` y `debt.capacity_limit_pct`, K54 y K44). Ningún umbral es constante en el código.
- K-1 a K-4 usan la «Tarjeta de métrica» (§1.2; K-4 con una lista compacta). B-1 usa el barómetro de DESIGN §8.9. Las gráficas G-1 a G-6 siguen todas las reglas de §4 (orden fijo de series, forma y patrón, etiquetas de dinero en texto, `<figure>` con `<figcaption>`, resumen sin variaciones calculadas, «Ver como tabla»).
- Los campos y operaciones **NUEVOS** de §2.6 están **publicados en C1 v1.5** (agente `s1-c1-v15`). Si un servidor anterior a C1 v1.5 responde sin un campo nuevo, la línea que lo usa no se pinta y el bloque muestra lo demás; un bloque que depende solo de una operación nueva muestra «Disponible cuando el servidor se actualice» (`font.size.caption`, `info`), nunca una cifra inventada.

### K-1 · Gasto del mes (indicador principal)
- **Métrica C1**: M-01 (N-020), ampliada en C1 v1.5 para K-1: la fila M-01 de C1-datos §9 trae la fórmula, el periodo (mes y mes calendario anterior), las tablas y el ejemplo del oráculo de los 8 campos nuevos (C1-09; corrección r2) · `getDashboardMonthVsPlan` (GET `/api/v1/dashboard/month-vs-plan`) · esquema `DashboardMonthVsPlan`: `month`, `actual_cents`, `planned_cents`, `used_pct`, `remaining_cents`, `status` (`normal`, `warning`, `exceeded`), `days_left` y, **nuevos en C1 v1.5** (§2.6): `cash_cents`, `previous_month`, `previous_actual_cents`, `previous_cash_cents`, `actual_difference_cents`, `cash_difference_cents`, `categories_warning_count`, `categories_exceeded_count`.
- **Fórmula** (motor S2): gastado (devengado) = M-01 real = Σ parte del hogar devengada en el mes; pagado en efectivo = efectivo del hogar del mes, la misma cifra que `cash_household_cents` de M-11 para ese mes (contado + partes del hogar de pagos de crédito con fecha en el mes); plan, % usado, restante y estado como M-01 (FIRME desde la puerta 2, Q2-19 B y Q2-20 A); mes anterior = las mismas dos cifras del mes calendario anterior, ausentes si el sistema no tiene datos de ese mes (sin datos no es cero); diferencias = mes − mes anterior, con signo, calculadas por el motor; conteos = categorías del mes con `status` `warning` o `exceeded` en M-09.
- **Periodo**: mes en curso (`period` = `2026-10`).
- **Contenido** (icono `wallet`):
  - Título: «Gasto del mes».
  - Valor principal: `[gastado]` con la etiqueta «gastado en [mes]» (neutro).
  - Valor secundario (360): «[n] % de [plan] · te quedan [restante]»; si `remaining_cents` < 0: «−[monto] rebasado» con `color.role.warning.*` e icono `triangle-alert` (DESIGN §3.3). Barra de avance de DESIGN §7.6 con `role="meter"` y `aria-valuetext` «[n] % del plan».
  - Desde 412 suma: «Pagado en efectivo [efectivo]» y la comparación «[diferencia] más que [mes anterior]» (icono `arrow-up-right`, `color.role.expense.*`), «[diferencia] menos que [mes anterior]» (icono `arrow-down-left`, `ui.text.secondary`) o «Igual que [mes anterior]»; sin mes anterior: «Sin datos de [mes anterior]». Las palabras «más» y «menos» van siempre en el texto (DESIGN §7.6).
  - Insignias: `warning` → «Vas en [n] %»; `exceeded` → «Rebasaste el plan» (rol advertencia, `triangle-alert`, como la tarjeta 1 de U3 1.x); si `categories_exceeded_count` > 0, «[n] categorías rebasadas» (rol advertencia), que lleva a G-3.
- **Acción primaria única**: enlace «Ver presupuesto» → P-06 Presupuesto del mes (primer paso del desglose de §2.2).
- **Vacío**: sin `planned_cents`: «—» en el % y «Aún no hay presupuesto para [mes]» con «Llevas [gastado] gastado»; la acción cambia a «Capturar presupuesto» → P-26.02.
- **Error**: «No se pudo cargar el gasto del mes. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable** (oct-2026 = P1, nov-2026 = P2; regla de ejemplos al inicio): reloj en octubre de 2026: gastado **$1,197.64** (`periodos[0].devengado.axel`, M-01) y pagado en efectivo **$598.82** (`periodos[0].efectivo.axel`: la parte del hogar del pago del 30 de octubre, R-13); «Sin datos de septiembre». Reloj en noviembre de 2026: gastado **$1,197.64** y efectivo **$1,197.64** (`periodos[1].devengado.axel` y `periodos[1].efectivo.axel`); contra octubre: «Igual que octubre» en lo gastado y «**$598.82** más que octubre» en efectivo. Plan: marcador `[plan]` (lo captura Axel).

### K-2 · Ingreso del mes
- **Métrica C1**: N-100 (P-08) · `getIncomeMonth` (GET `/api/v1/incomes/month`) · esquema `IncomeMonth`: `month`, `fortnights[]` (`pay_date`, `status`: `expected` o `captured`, `expected_cents`, `amount_cents`), `extras[]`, `captured_cents`, `expected_pending_cents`, `extras_cents`, `total_cents` y, **nuevos en C1 v1.5**: `previous_month`, `previous_fortnights_cents`, `previous_extras_cents`, `fortnights_difference_cents`. Métrica M-18 «Ingreso del mes» de C1-datos §9 (C1 v1.5, §2.6).
- **Fórmula** (motor S2): quincenas del mes = `captured_cents` + `expected_pending_cents` (capturadas, y las esperadas de P-26.03 que aún no se capturan); extra = `extras_cents`; total = `total_cents`; mes anterior y diferencia de quincenas, del motor.
- **Periodo**: mes en curso (`period` = `2026-10`).
- **Contenido** (icono `banknote`):
  - Título: «Ingreso del mes».
  - Valor principal: «+[quincenas]» con `color.role.income.*` e icono `arrow-down-left` (DESIGN §3.3), etiqueta «de quincenas».
  - Valor secundario (360): «Extra este mes +[extra]» o, con `extras_cents` = 0, «Sin extra este mes».
  - Desde 412 suma: «[n] de [m] quincenas capturadas» (conteo de `fortnights[]` a la vista) y «[diferencia] más que [mes anterior]» o «menos que», o «Sin datos de [mes anterior]».
  - Insignia si alguna quincena llega con `status` = `expected`: «Ingreso esperado» (rol información, `info`).
- **Acción primaria única**: enlace «Ver ingresos» → P-08 Ingresos.
- **Vacío** (sin quincena capturada ni esperada en P-26.03 y sin extras): «—» y «Falta tu ingreso quincenal»; la acción cambia a «Capturar ingreso» → P-26.03.
- **Error**: «No se pudo cargar el ingreso del mes. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo**: marcadores `[quincenas]` y `[extra]`; el oráculo no trae ingresos.

### K-3 · Créditos
- **Métrica C1**: M-12 a hoy (deuda viva) y M-04 (disponible) · `getCreditTotals` (GET `/api/v1/credits/totals`) · esquema `CreditTotals`: `live_debt_cents`, `month_total_cents`, `month_household_cents`, `month_external_cents` y, **nuevos en C1 v1.5**, `month_payment_count` y `month_pending_household_cents`; y `getDashboardNextPayment` (GET `/api/v1/dashboard/next-payment`) para `availability[]` (`credit_id`, `availability_reading`, `available_cents`, `availability_confirmed`). Nombres de línea: `listCredits`.
- **Fórmula** (motor S3): deuda viva total = capital + intereses por pagar de todas las líneas (M-12, a hoy); pagos del mes = Σ pagos programados con fecha en el mes en curso de todas las líneas (calendario unificado, M-13), total y partes; cuántos = número de esos pagos; falta pagar = parte del hogar aún sin pagar de esos pagos; disponible = M-04 con la **lectura activa**.
- **Qué lectura muestra** (U3-07, decisiones §C): como la tarjeta 3 de U3 1.x, solo la lectura activa del interruptor (por defecto **B**) con su etiqueta «Lectura B»; la otra se ve lado a lado en Detalle › D-3 y en P-18.
- **Contenido** (icono `credit-card`):
  - Título: «Créditos».
  - Valor principal: `[deuda viva]` con `color.role.debt.*` y la etiqueta «deuda viva total» (rol con icono, DESIGN §2.4).
  - Valor secundario (360): «Pagos de [mes]: [total del mes] · [n] pagos».
  - Desde 412 suma: «Hogar [parte del hogar] · Karen [parte de Karen]» (Karen solo si `month_external_cents` > 0) y «Falta pagar [falta pagar]».
  - Bloque «Disponible hoy», una fila por línea de `availability[]`: «Disponible hoy [monto] · Lectura [A/B]»; si `availability_confirmed` = false, insignia **«Pendiente confirmar»** (rol información, `info`) con nombre accesible «Reposición de capital pendiente de confirmar con Kueski». Con más de una línea revolvente: la primera y «+[n] líneas».
- **Acción primaria única**: enlace «Ver créditos» → P-17 Todas las líneas.
- **Vacío**: sin créditos: «Aún no hay créditos» y la acción «Agregar crédito» → P-26.05; sin pagos en el mes: «Sin pagos en [mes]» y el resto igual.
- **Error**: «No se pudieron cargar los créditos. Revisa que estés en la red de casa.» · «Reintentar». Si solo falla `getDashboardNextPayment`, el bloque «Disponible hoy» dice «—» con «Reintentar» y el resto se ve.
- **Ejemplo trazable** (U3-07): hoy 10-oct-2026: «Pagos de octubre: **$973.75** · 1 pago», «Hogar **$598.82** · Karen **$374.93**» (R-13; `periodos[0].efectivo.n_pagos` = 1); «Disponible hoy **$6,905.50** · Lectura B · Pendiente confirmar» y con el interruptor en A **$10,000.00** (R-64). Deuda viva de hoy: marcador `[deuda viva]`; la cifra de control del motor es la deuda viva al 10-abr-2027, **$6,293.44** (R-53).

### K-4 · Pagos próximos
- **Métrica C1**: M-13 (N-192) · `listScheduledPayments` (GET `/api/v1/scheduled-payments`: `id`, `due_date`, `amount_cents`, `credit_id`, `shares[]` con `person_id`, `amount_cents`, `paid_cents`, `status`; `status`) en **dos peticiones**, como D-3 (§3): (1) vencidos, solo si `getDashboardNextPayment` trae `is_overdue` = true (la respuesta de K-3 se reutiliza); (2) próximos: `status` = `pending` y `period` = `[hoy]..[hoy + 30]`. Las dos con orden `due_date` y `limit` = **N**. Días de «Por vencer»: `getAlertSettings` (§1.2).
- **N y rango** (**PROVISIONAL**, bloqueo `docs/estado/bloqueos/U3-tablero-recordatorio.md`): N = 5 filas en total (primero los vencidos), rango de hoy a hoy + 30 días («los pagos del mes» de Axel, también cerca del fin de mes). N = `parameters["dashboard.payment_reminder_count"].value` de `getDashboardSettings` (C1 v1.5; clave de `app_setting` con valor inicial 5, K44), leído una sola vez, como `getAlertSettings` en §1.2; la interfaz no fija N.
- **Contenido** (icono `calendar-clock`): título «Pagos próximos»; una fila por pago: fecha con día de la semana («Vie 30 oct 2026»), línea, monto total y estado de §1.4 («Vencido» con «Sin marcar como pagado», «Por vencer», «Parcial» con «Falta la parte de Karen», «Pendiente»; «Pagado» no llega). Desde 412, debajo: «Hogar [parte] · Karen [parte]». Cada fila es un enlace a P-18 de su línea (F-4). Si hay más pagos que N: «y [n] más en Calendario».
- **Acción primaria única**: enlace «Ver calendario» → pestaña Calendario de Inicio (§5).
- **Vacío**: «Sin pagos en los próximos 30 días»; sin créditos: «Aún no hay créditos» · «Agregar crédito» → P-26.05.
- **Error**: «No se pudieron cargar los pagos. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable**: hoy 10-oct-2026, una fila: «Vie 30 oct 2026 · Kueski · **$973.75** · Pendiente», «Hogar **$598.82** · Karen **$374.93**» (R-13); el pago del 15 de noviembre (R-14) cae fuera de los 30 días. Jue 29 oct 2026: «Vence mañana» (§1.4, A-17).

### B-1 · Barómetro de capacidad de endeudamiento
- **Métrica C1**: M-20 «Capacidad de endeudamiento» (C1 v1.5) · **NUEVA** `getDashboardDebtCapacity` (GET `/api/v1/dashboard/debt-capacity`, `period` opcional, por omisión el mes en curso) · esquema nuevo `DashboardDebtCapacity` (campos exactos en §2.6): `month`, `income_cents`, `income_is_expected`, `extras_cents`, `payments_cents`, `payments_total_cents`, `payment_count`, `ratio_pct`, `ratio_total_pct`, `warning_pct`, `limit_pct`, `scale_max_pct`, `status`, `headroom_cents`, `has_credits`.
- **Fórmula** (motores S2 y S3; K54 decisión 1): capacidad usada = pagos del mes de todos los créditos ÷ ingreso mensual neto × 100. `payments_cents` = Σ parte del hogar de los pagos programados con fecha en el mes, de todas las líneas, pagados o no (la misma cifra que `month_household_cents` de K-3). `income_cents` = quincenas del mes, capturadas más esperadas sin capturar, **sin extras** (**PROVISIONAL**). `ratio_pct` = entero con ROUND_HALF_UP (protocolo §6). Estado: `normal` (verde) si `ratio_pct` ≤ `warning_pct`; `warning` (amarillo) si `ratio_pct` ≤ `limit_pct`; `exceeded` (rojo) si `ratio_pct` > `limit_pct`. Así se lee K54 decisión 2, «verde hasta 30 %, amarillo de 30 a 40 %, rojo arriba de 40 %»: 30 exacto es verde y 40 exacto es amarillo. `headroom_cents` = `limit_pct` × `income_cents` ÷ 100 (ROUND_HALF_UP) − `payments_cents`: lo que todavía cabe en pagos al mes antes del rojo; es negativo si ya se pasó. `scale_max_pct` = el múltiplo de 10 igual o mayor que 1.5 × `limit_pct` (60 con los valores iniciales).
- **PROVISIONAL** (bloqueo `docs/estado/bloqueos/U3-tablero-barometro.md`, puerta 4): qué es «ingreso neto» (sin extras), qué pagos cuentan (parte del hogar; `ratio_total_pct` con la parte de Karen solo se informa), el redondeo a entero y la escala.
- **Umbrales** (K54 decisión 2, K44): `debt.capacity_warning_pct` = 30 y `debt.capacity_limit_pct` = 40, configurables en Ajustes; llegan en la respuesta.
- **Contenido** (icono `credit-card`; componente «Barómetro», DESIGN §8.9):
  - Título: «Capacidad de endeudamiento».
  - Arco de 0 % a `scale_max_pct` con tres tramos: verde hasta `warning_pct` (`ui.role.success`), amarillo hasta `limit_pct` (`ui.role.warning`) y rojo hasta el final (`ui.role.error`). Marcas en texto «[warning] %» y «[limit] %» en los límites de tramo. Aguja en `ratio_pct`; si `ratio_pct` > `scale_max_pct`, la aguja queda al final y el dígito dice la cifra exacta.
  - Dígito: «[n] %» (`font.size.title-1`, `font.weight.bold`, `tnum`, `ui.text.default`) y debajo «de [limit] % de tope».
  - Insignia de estado, siempre con icono y texto: `normal` → «Verde · dentro del rango» (`circle-check`, rol éxito); `warning` → «Amarillo · cerca del tope» (`triangle-alert`, rol advertencia); `exceeded` → «Rojo · arriba del tope» (`circle-alert`, rol error).
  - Valor secundario (360): «Pagos [pagos] de un ingreso de [ingreso]».
  - Desde 412 suma: «Te caben [headroom] más al mes antes del rojo» o, con `headroom_cents` < 0, «−[monto] arriba del tope» (rol advertencia, `triangle-alert`, DESIGN §3.3); si `income_is_expected`, insignia «Ingreso esperado» (`info`). Desde 768: «Con la parte de Karen: [n] %» (`ratio_total_pct`, informativo, §9).
  - Texto accesible: `role="meter"`, `aria-valuemin` = 0, `aria-valuemax` = `scale_max_pct`, `aria-valuenow` = `ratio_pct` y `aria-valuetext` «Capacidad de endeudamiento [n] % de [limit] %: [verde, amarillo o rojo]». Ejemplo **SINTÉTICO** (no sale del oráculo, que no trae ingresos): «Capacidad de endeudamiento 27 % de 40 %: verde».
- **Acción primaria única**: enlace «Ver créditos» → P-17 Todas las líneas.
- **Vacío**: sin `income_cents`: arco sin aguja, dígito «—», «Falta tu ingreso del mes para calcular la capacidad» y la acción «Capturar ingreso» → P-26.03. Sin créditos (`has_credits` = false): el dígito muestra el `ratio_pct` que llega (0 %, cifra real: no hay pagos) con la línea «Sin créditos activos».
- **Error**: «No se pudo calcular la capacidad de endeudamiento. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable** (oct-2026 a dic-2026 = P1 a P3): pagos del hogar **$598.82** en octubre (1 pago), **$1,197.64** en noviembre (2 pagos) y **$1,482.57** en diciembre (2 pagos) (`periodos[i].efectivo.axel` y `n_pagos`; en P1 a P3 el efectivo del hogar son solo pagos de crédito); con la parte de Karen, **$973.75**, **$1,947.50** y **$2,232.43** (`periodos[i].efectivo.total`). Ingreso: marcador `[ingreso]`; por eso el porcentaje del ejemplo es `[n]`.

### 2.1 Gráficas de análisis (G-1 a G-6)
Reglas comunes: las de §4 (DESIGN §7.7 y §8): orden fijo de series con forma, patrón y etiqueta; etiquetado selectivo (DESIGN §8.5); un solo eje Y; `<figure>` con `<figcaption>` y resumen de una frase sin variaciones calculadas; «Ver como tabla» con las mismas cifras; esqueleto sin ejes falsos; alto `size.chart.height-compact` en < 768 y `size.chart.height` desde 768. **0 donas** en Inicio (G-2 usa barras; §4). G-4 y G-5 son T-2 y T-3 de §4, que **se conservan** en Tendencias y se citan aquí sin cambio. **Horizonte de las gráficas** (**PROVISIONAL**, bloqueo `docs/estado/bloqueos/U3-tablero-horizonte.md`): G-1, los 11 meses anteriores y el mes en curso; G-6, los 6 anteriores, el en curso y los 5 siguientes (la proyección es lo que previene); G-4 y G-5, las ventanas de T-2 y T-3. Para G-1 y G-6 la PWA hace **una** petición a `getDashboardMonthlyOverview` con `from` = mes en curso − 11 y `to` = mes en curso + 5 (17 meses) y cada gráfica muestra su tramo.

### G-1 · Ingreso, gasto y pagos por mes
- **Métrica C1**: M-21 «Resumen mensual del tablero» (C1 v1.5) · **NUEVA** `getDashboardMonthlyOverview` (GET `/api/v1/dashboard/monthly-overview`, `from` y `to` obligatorios) · `months[]`: `month`, `is_projection`, `income_total_cents`, `income_fortnights_cents`, `income_extras_cents`, `income_is_expected`, `accrued_cents`, `cash_cents`, `planned_cents`, `credit_payments_cents`, `credit_payments_total_cents`, `debt_ratio_pct`, `debt_status` (§2.6).
- **Fórmula y periodo**: por mes: ingreso = quincenas + extras (K-2); gasto = devengado del hogar (M-01 real); pagos de crédito = parte del hogar de los pagos con fecha en el mes (B-1). Los 11 meses anteriores y el en curso.
- **Patrón**: **columnas agrupadas**, 3 series (DESIGN §8.2, máximo 3): 1 «Ingreso» = `color.chart.1`, liso; 2 «Gasto» = `color.chart.2`, rayas a 45°; 3 «Pagos de crédito» = `color.chart.3`, rayas a 135°; 2 px de separación, sin contorno. Valores en texto sobre las columnas del mes en curso; los demás, en la tabla. Un mes sin ingreso o sin gasto no lleva esa columna (sin datos no es cero) y la tabla dice «—».
- **Nota fija** bajo la leyenda: «Los pagos de crédito pagan compras que ya cuentan en el gasto de su mes; no se suman al gasto».
- **Resumen** (DESIGN §7.7): «En [mes en curso]: ingreso [monto], gasto [monto] y pagos de crédito [monto]».
- **Vacío**: «Aún no hay meses con ingreso ni gasto» · «Capturar ingreso» → P-26.03.
- **Ejemplo trazable**: gasto **$1,197.64** en octubre y en noviembre de 2026 y **$1,767.51** en diciembre (`periodos[0..2].devengado.axel`); pagos de crédito **$598.82**, **$1,197.64** y **$1,482.57** (`periodos[0..2].efectivo.axel`); ingreso: marcador `[ingreso]`.

### G-2 · Deuda viva por línea
- **Métrica C1**: M-12 (N-092) · `getLiveDebt` (GET `/api/v1/debt/live`, `date` = hoy) · `LiveDebt`: `date`, `total_cents`, `credits[]` (`credit_id`, `total_cents`); nombres con `listCredits`.
- **Patrón**: **barras horizontales** (DESIGN §8.2, «¿cuánto en cada…?»), orden descendente por monto, una serie `color.chart.1` sin patrón, valor en texto al final de cada barra. No es dona: con una sola línea (hoy solo Kueski) un anillo completo no informa y con más de 5 líneas DESIGN §8.4 pide barras; Inicio sigue con 0 donas.
- **Resumen**: «Deuda viva hoy: [total]; la mayor, [línea] con [monto]».
- **Vacío**: «No hay deudas vivas» (con `circle-check`).
- **Ejemplo**: a la fecha de control del motor, 10-abr-2027, Kueski **$6,293.44** (R-53 = **$1,258.70** R-54 + **$5,034.74** R-55, M-12); hoy, marcador `[deuda viva]`.

### G-3 · Presupuesto por categoría (con la sugerencia subir o recortar)
- **Métrica C1**: M-09 (N-081) · `getBudgetByCategory` (`period` = mes en curso; `rows[]`: `category_id`, `planned_cents`, `actual_cents`, `difference_cents`, `used_pct`, `status`, `expense_count`) · umbrales `alerts.budget_warning_pct` y `alerts.budget_limit_pct` de `getAlertSettings` (§1.2) · sugerencia: M-22 (C1 v1.5) · **NUEVA** `getBudgetSignals` (GET `/api/v1/budget/signals`, `period`) · `BudgetSignals`: `month`, `months_considered`, `margin_pct`, `rows[]` (`category_id`, `signal`: `raise`, `cut` o `none`, `months[]` con `month`, `planned_cents`, `actual_cents`, `used_pct`). Nombres con `listCategories`.
- **Regla de la sugerencia** (**PROVISIONAL**, bloqueo `docs/estado/bloqueos/U3-tablero-senal-presupuesto.md`; la calcula el motor S2, no la interfaz): con los **3** meses cerrados anteriores al mes pedido (`budget.signal_months`, valor inicial 3), una categoría lleva `raise` («subir presupuesto») si en los 3 su `used_pct` pasó de 100 + **10** (`budget.signal_margin_pct`, valor inicial 10), y `cut` («recortar presupuesto») si en los 3 quedó abajo de 100 − 10. Un mes sin plan de esa categoría corta la racha. El mes en curso no cuenta porque no ha cerrado.
- **Patrón**: **barras horizontales en % del plan** (DESIGN §8.2, fila de 1.1): una barra por categoría con plan, largo = `used_pct`, serie `color.chart.1`; con `status` `warning` o `exceeded`, la barra va en `ui.role.warning` con `triangle-alert` y el texto «Vas en [n] %» o «Rebasada». Marcas verticales en `alerts.budget_warning_pct` (1 px, trazo 4 · 4, `ui.viz.baseline`) y en `alerts.budget_limit_pct` (2 px, `ui.viz.baseline`), cada una con su «[n] %» en texto. El eje va de 0 % al mayor entre 100 % y el mayor `used_pct` a la vista. Orden descendente por `used_pct` (la interfaz ordena lo visible). Etiqueta de cada barra: «[categoría] · [real] de [plan] · [n] %». Con `signal` ≠ `none`, insignia de rol información (`info`): «Sugerencia: subir presupuesto» o «Sugerencia: recortar presupuesto», con nombre accesible «[categoría]: 3 meses seguidos arriba del plan por más de 10 %» (o «abajo»).
- **Tamaño**: a 360 y 412, las 8 categorías de mayor % y el enlace «Ver todas» → P-06; desde 768, todas. Las categorías con gasto y sin plan no llevan barra: línea «[n] categorías con gasto y sin plan» · «Capturar presupuesto» → P-26.02.
- **Desglose**: tocar una barra abre P-02 con `category_id` y `period` = mes (§2.2).
- **Resumen**: «En [mes]: [n] categorías arriba de [warning] % y [m] rebasadas» (conteo de lo que está a la vista).
- **Vacío**: «Aún no hay presupuesto por categoría» · «Capturar presupuesto» → P-26.02. Sin `getBudgetSignals` (error o C1 sin publicarla), las barras se ven sin insignias de sugerencia.
- **Ejemplo**: marcadores (`[categoría]`, `[real]`, `[plan]`, `[n]`): el plan lo captura Axel (M-09).

### G-4 · Devengado contra efectivo
Es **T-2** de §4 sin cambio (M-11, `getAccrualVsCash`, capas `cumulative_cash_cents` y `cumulative_accrued_cents`, Q2-65 A), con su resumen, su nota de proyección y sus ejemplos (R-45 a R-52). En el tablero no lleva filtros. Usa la misma petición que la pestaña Tendencias.

### G-5 · Deuda viva proyectada
Es **T-3** de §4 sin filtro de crédito (M-17, `getDebtTrend`), con su resumen y su ejemplo (R-53 a R-55). En el tablero no lleva filtros.

### G-6 · Tendencia del barómetro
- **Métrica C1**: M-21 (C1 v1.5) · `getDashboardMonthlyOverview` (la petición de G-1) · `months[]`.`debt_ratio_pct`, `debt_status`, `is_projection`; `warning_pct` y `limit_pct` de la respuesta.
- **Fórmula y periodo**: por mes, la misma capacidad de B-1; meses posteriores al en curso (`is_projection` = true) con los pagos programados (S3) y las quincenas esperadas de P-26.03. Los 6 meses anteriores, el en curso y los 5 siguientes (**PROVISIONAL**, horizonte).
- **Patrón**: **línea** (DESIGN §8.2, tiempo): una serie «Capacidad de endeudamiento», `color.chart.1`, 2 px, marcador círculo; meses de proyección con trazo 8 · 4 y la nota «Proyección con pagos programados e ingreso esperado». Líneas de referencia horizontales de 1 px en `ui.viz.baseline` en `warning_pct` y `limit_pct`, con su «[n] %» en texto al final; línea vertical «Hoy» (DESIGN §8.3). Bajo cada punto, el estado no va solo en color: lo dicen la tabla y la descripción emergente («Diciembre de 2026, [n] %: amarillo»). Eje Y de 0 % al mayor entre `scale_max_pct` y el mayor `debt_ratio_pct` a la vista. Un mes sin ingreso no lleva punto y corta la línea; la tabla dice «—» con «Falta ingreso».
- **Resumen**: «[mes en curso]: [n] % ([estado]); el mayor de la proyección es [mes] con [n] %» (el mayor de lo visible).
- **Vacío**: «Falta tu ingreso para calcular la capacidad» · «Capturar ingreso» → P-26.03.
- **Ejemplo**: pagos del hogar por mes de B-1 (**$598.82**, **$1,197.64**, **$1,482.57**); con ingreso `[ingreso]`, porcentajes `[n]`.

### 2.2 Desglose en cadena (del tablero al precio del kilo)
Responde «cuánto costó el kg de cebollas» (K54) sin pantallas nuevas: cada paso **enlaza** una pantalla que ya existe (U1 §2.1) con sus filtros.

| Paso | Dónde | Qué toca la persona | Qué ve | Necesidad U1 |
|---|---|---|---|---|
| 1 | Inicio › Resumen (K-1 o G-3) | Una barra de G-3, o «Ver presupuesto» en K-1 | Categorías del mes con su semáforo contra el plan (P-06, o P-02 con la categoría) | N-020, N-081 |
| 2 | P-02 Todos los gastos | Llega con `category_id` y `period` = mes | Gastos del mes de esa categoría y sus totales («Resumen del periodo» con el mes anterior) | N-040, N-041, N-042 |
| 3 | P-04 Gasto | Toca un gasto | Renglones del ticket: producto, presentación, cantidad, precio unitario e importe | N-061 |
| 4 | P-15 Productos y precios | Toca el producto de un renglón | Precios del producto por tienda y presentación, con precio por unidad base, fecha y fuente | N-171 |
| 5 | P-15 › Tendencias | Toca «Tendencias» | Historial del precio por unidad base por fecha y tienda (p. ej. «Cebolla · [precio] por kg» en cada tienda) | N-171 |

Toques desde Inicio hasta el historial del precio: 4 (barra, gasto, producto y pestaña). El paso 4 pide que el renglón de P-04 enlace a P-15 con el producto elegido; si U1 o la implementación no lo tienen, va como solicitud (§15). Ejemplo: marcador `[precio]`; la semilla no trae cebolla.

### 2.3 Qué miro para prevenir
La tabla que pidió Axel («a mi me sirven bastante para irme previniendo»): pregunta del mes, dónde se contesta y qué señal mirar. Las acciones son sugerencias de lectura, no reglas del sistema.

| Pregunta | Dónde | Señal que previene | Qué hacer |
|---|---|---|---|
| ¿Cuánto gastamos realmente este mes? | K-1 | % del plan, «Vas en [n] %», «Rebasaste el plan» | Bajar a G-3 y a la categoría |
| ¿Vamos peor que el mes pasado? | K-1 (desde 412) | «[diferencia] más que [mes anterior]» | Ver qué categoría creció (G-3, §2.2) |
| ¿Qué categoría se pasa? | G-3 | Barra amarilla, «Rebasada», marcas de 80 % y 100 % | Desglose §2.2 hasta el gasto y el producto |
| ¿Subimos o recortamos el presupuesto? | G-3 | «Sugerencia: subir presupuesto» o «recortar» (3 meses seguidos) | Ajustar el plan en P-26.02 |
| ¿Subió el precio de algo (el kg de cebolla)? | §2.2 → P-15 › Tendencias | Historial por unidad base y por tienda | Comparar tienda en P-14 |
| ¿Cuánto del ingreso se va en créditos? | B-1 | Color y dígito del barómetro; «Te caben [monto] más» | Con amarillo o rojo, no tomar otra compra a crédito; simularla antes en P-20 |
| ¿Va a empeorar en los próximos meses? | G-6 | Puntos de proyección arriba de la línea de 30 % o 40 % | Planear antes del mes que se pasa |
| ¿Qué pagos vienen? | K-4 | «Vence mañana», «Vence hoy», «Vencido» | Pagar y marcar en P-18 |
| ¿Cuánto debemos en total y cuándo baja? | K-3 y G-5 | Deuda viva total y su curva | Ver la línea en P-18 |
| ¿Ya entró la quincena? ¿Hubo extra? | K-2 | «Ingreso esperado», «Extra este mes» | Capturar la quincena en P-08 |
| ¿Gastamos más de lo que entra? | G-1 | Columna de gasto más alta que la de ingreso | Revisar G-3 y el plan |
| ¿Lo que se consume se está pagando? | G-4 | Capa «diferencia» que crece mes con mes | Ver pagos (K-4) y la deuda (G-5) |

### 2.4 Orden y rejilla (móvil 412 primero, escritorio 1440)
- **360 y 412** (1 columna, orden de apilado): 1 K-1 · 2 B-1 · 3 K-4 · 4 K-2 · 5 K-3 · 6 G-3 · 7 G-1 · 8 G-6 · 9 G-2 · 10 G-5 · 11 G-4 · 12 «Más del hogar» (tarjetas 2, 4, 5 y 6). Primero el propósito (gasto real), luego el «panel máximo» de Axel y lo que pide acción (pagos). A 360, cada bloque con título, valor principal y una línea secundaria; desde 412, una línea más (U1 §6.2). La cifra principal y el dígito del barómetro nunca se truncan (U1 §6.1 regla 5).
- **768** (8 columnas): K-1 y B-1 en una fila (4 + 4); K-4 a todo el ancho; K-2 y K-3 (4 + 4); gráficas, una por fila; «Más del hogar», 2 × 2.
- **1024 y 1440** (12 columnas; en 1024 con la barra lateral fija): fila 1, K-1 (6) y B-1 (6); fila 2, K-2 (4), K-3 (4) y K-4 (4); filas 3 a 5, gráficas de dos en dos (6 + 6): G-3 con G-1, G-6 con G-2, G-5 con G-4; fila 6, «Más del hogar» (1440: 4 tarjetas de 3 columnas; 1024: 2 × 2). Medianil `space.6`, relleno de tarjeta `space.6`; todas las tarjetas de una fila con el mismo alto (DESIGN §7.6).
- **Teclado y lector**: el orden del foco es el de 360 en todos los anchos (el DOM sigue ese orden; la rejilla no lo cambia).

### 2.5 Más del hogar (tarjetas que se conservan)
**PROVISIONAL** (bloqueo `docs/estado/bloqueos/U3-tablero-mas-del-hogar.md`): K54 retira la regla de 6 tarjetas pero no pide quitar de Inicio la quincena, las comidas, la canasta ni el ahorro. Por eso las tarjetas 2, 4, 5 y 6 de U3 1.x siguen **sin cambio de contenido** al final del tablero, bajo el título de sección «Más del hogar» (`font.size.title-3`). Así siguen en pie los flujos que entran por ellas (F-2 y F-3 ruta A de U1 §4) y sus pruebas. La tarjeta 1 pasa a K-1. La tarjeta 3 se reparte entre K-3 (disponible A/B) y K-4 (próximo pago con estado). Cada tarjeta conserva su texto, su operación, sus estados y sus ejemplos:

### Tarjeta 2 · Flujo de la quincena
- **Métrica C1**: M-02 (N-021) · `getDashboardFortnightFlow` (GET `/api/v1/dashboard/fortnight-flow`) · esquema `DashboardFortnightFlow`, campos `start_date`, `end_date`, `income_cents`, `income_is_expected`, `spent_cents`, `committed_cents`, `closing_balance_cents`.
- **Fórmula** (motor S2, efectivo, parte del hogar): ingresos = quincena capturada o esperada + extras de la quincena; egresos hechos = salidas de efectivo hasta hoy (contado + partes del hogar pagadas); comprometido = partes del hogar pendientes de la quincena + fijos y venideros sin pagar; saldo al cierre = ingresos − egresos hechos − comprometido.
- **Periodo**: quincena en curso con los límites de K8 (del 1 al 15 y del 16 al último día; el 31 cae en la segunda). Petición: `period` = quincena en curso (`2026-10-q1` o `2026-10-q2`).
- **Contenido** (icono `banknote`):
  - Título: «Flujo de la quincena»; subtítulo «Del [día] al [día] [mes]» con las fechas de `start_date` y `end_date` (p. ej. «Del 16 al 31 oct»).
  - Valor principal: `[saldo al cierre]` con la etiqueta «al cierre». Negativo: «−[monto] faltante» con `color.role.warning.*` e icono `triangle-alert`.
  - Valor secundario (360): «Entra +[ingresos] · Sale −[egresos + comprometido]» (ingreso con `arrow-down-left` y `color.role.income.*`; salida con `arrow-up-right` y `color.role.expense.*`). Desde 412 separa: «Ya salió −[egresos] · Por pagar −[comprometido]».
  - Insignia si `income_is_expected` = true: «Ingreso esperado» (rol información, icono `info`).
- **Acción primaria única**: enlace «Ver flujo» → P-07 Flujo y proyección.
- **Vacío** (sin quincena capturada ni esperada en P-26.03): ingresos «—» con «Falta tu ingreso quincenal»; la acción cambia a «Capturar ingreso» → P-26.03 (Ajustes › Ingresos). Egresos y comprometido se siguen mostrando.
- **Error**: «No se pudo cargar el flujo de la quincena. Revisa que estés en la red de casa.» · «Reintentar».
- **FIRME** desde la puerta 2 (Q2-14 A; bloqueo `docs/estado/bloqueos/S2-flujo-quincenal.md`, que sigue abierto solo por Q2-12): una partida planificada se da por cumplida con un gasto o ingreso con el mismo `planned_item_id` en la misma quincena (o en el mismo mes si su frecuencia es mensual o mayor), así un fijo ya pagado no cuenta dos veces en «Por pagar»; la tarjeta muestra `committed_cents` tal como llegue. La parte del hogar del pago del 30 de octubre (R-13) no depende de esto.
- **Ejemplo trazable**: hoy 16-oct-2026, quincena «Del 16 al 31 oct»; «Por pagar» incluye la parte del hogar del pago del 30 de octubre, **$598.82** (R-13, M-02); ingresos: marcador `[ingresos]`.

### Tarjeta 4 · Comidas de la semana
- **Métrica C1**: M-05 (N-023) · `getDashboardWeekMeals` (GET `/api/v1/dashboard/week-meals`) · esquema `DashboardWeekMeals`: `week_start`, `week_end`, `cost_cents`, `persons[]` (`person_id`, `avg_kcal_per_day`, `days_with_plan`, `daily_kcal_target`), `unplanned_slots`. Nombres: `display_name` de `listPersons`.
- **Fórmula** (motor S4): costo = Σ costo de las casillas con receta («fuera» no suma); kcal promedio por día de cada persona = Σ kcal asignadas ÷ días con al menos un tiempo asignado; casillas sin armar = pares día × tiempo obligatorio (desayuno, comida, cena) sin receta ni «fuera» para alguna persona, máximo 21 (K43, FIRME).
- **Periodo**: **semana en curso** de lunes a domingo (K8), como dice U1 N-023; K38 **no** cambia el contenido de la tarjeta. La PWA manda siempre `period` explícito con la semana ISO en curso (p. ej. `2026-W41`) para no depender del valor por omisión de la operación.
- **Contenido** (icono `utensils`):
  - Título: «Comidas de la semana»; subtítulo «Lun [d] al Dom [d] [mes]».
  - Valor principal: `[costo]` con la etiqueta «costo de la semana».
  - Valor secundario (360): «[n] casillas sin armar» (si es 0: «Semana completa» con `circle-check`). Desde 412, una línea por persona en orden fijo Axel, Ximena (DESIGN §2.6), con punto `color.person.axel.*` / `color.person.ximena.*` y el nombre en texto: «Axel · [kcal] kcal al día de [meta] · [n] de 7 días». Sin meta: «sin meta». Sin días con plan: «—».
- **Acción primaria única**: enlace «Ver plan» → **P-10 Plan de la semana**, que abre por defecto la semana según **K38** (U1 §0.2): los domingos, la semana que empieza el lunes siguiente; el resto de los días, la semana en curso; «Cambiar semana» siempre disponible en P-10.
  - Caso domingo: hoy **Dom 11 oct 2026** → la tarjeta muestra «Lun 5 al Dom 11 oct» (semana `2026-W41`, N-023) y «Ver plan» abre P-10 en «Lun 12 al Dom 18 oct» (`2026-W42`, K38).
  - Caso entre semana: hoy **Mié 14 oct 2026** → la tarjeta muestra «Lun 12 al Dom 18 oct» (`2026-W42`) y P-10 abre la misma semana.
- **Vacío** (`unplanned_slots` = 21 y sin casillas con receta): valor principal «—»; secundario «Aún no hay plan para esta semana»; la acción cambia a «Armar plan» → P-10 (misma regla K38).
- **Error**: «No se pudo cargar el plan de comidas. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo**: marcadores (`[costo]`, `[kcal]`, `[meta]`, `[n]`); el oráculo no trae comidas (M-05).

### Tarjeta 5 · Mejor canasta
- **Métrica C1**: M-06 (N-024) · `getDashboardBestBasket` (GET `/api/v1/dashboard/best-basket`) · esquema `DashboardBestBasket`: `basket` (`shopping_list` o `expense`), `expense_id`, `store_id`, `total_cents`, `difference_vs_second_cents`, `covered_count`, `item_count`, `stale_price_count`, `assumes_membership` y, desde C1 v1.3, `second_store_id` (Q2-86 A). Nombre de tienda: `listStores`.
- **Fórmula** (motor S5): lista de compras vigente (o la última compra con renglones si está vacía) comparada por tienda; ganadora = menor total; diferencia contra la segunda; cobertura = artículos con precio ÷ artículos; precios viejos usados; Sam's Club como si tuvieran membresía.
- **Periodo**: la **lista de compras vigente** (U1 §0.2): la de la semana por defecto según K38 (los domingos, la siguiente). El parámetro `week` es obligatorio (R3-01) y la PWA manda esa semana (Dom 11 oct 2026 → `2026-W42`; Sáb 10 oct 2026 → `2026-W41`).
- **Contenido** (icono `shopping-cart`):
  - Título: «Mejor canasta».
  - Valor principal: «[tienda] [total]» con punto del color de la tienda (`color.store.<tienda>.*`) y su nombre en texto (DESIGN §2.6 regla 1).
  - Valor secundario (360): «[diferencia] menos que [segunda tienda]», con el nombre de `second_store_id` (Q2-86 A); si C1 aún no lo entrega, «[diferencia] menos que la siguiente tienda». Desde 412 suma: «[covered] de [items] artículos con precio» y, si `stale_price_count` > 0, insignia «[n] precios viejos» (rol advertencia, `triangle-alert`).
  - Si `assumes_membership` = true: insignia «Como si tuvieran membresía» (rol información, `info`).
  - Si `basket` = `expense`: línea «Tu lista está vacía; se compara la compra del [fecha]».
- **Acción primaria única**: enlace «Comparar canasta» → P-14 (abre la semana según K38).
- **Vacío** (sin `store_id`: ningún artículo con precio, o sin lista ni compras con renglones): «—» y «Aún no hay precios para comparar»; la acción cambia a «Registrar precios» → P-15.
- **Error**: «No se pudo comparar la canasta. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable**: compra semilla del 10-oct-2026 (`basket` = `expense`): «Bodega Aurrera **$3,094.50**» (R-01), «**$138.50** menos que la siguiente tienda» (K-03 de `S5.md`: Walmart **$3,233.00**, `data.json` `totales.walmart`). Con `second_store_id` (C1 v1.3, Q2-86 A): «$138.50 menos que Walmart».

### Tarjeta 6 · Ahorro y metas
- **Métrica C1**: M-07 (N-025) · `getDashboardSavings` (GET `/api/v1/dashboard/savings`) · esquema `DashboardSavings`: `saved_total_cents`, `contributed_month_cents`, `planned_month_cents`, `next_goal` (`SavingsGoal`: `name`, `target_cents`, `target_date`, `saved_cents`, `progress_pct`, `status`: tres estados por Q2-18 A, `in_progress`, `achieved` y `overdue`, como en `docs/contratos/motores/S2-interfaz.md` 1.3 §3.13; C1 v1.3 todavía declara `on_track`, `behind` y `achieved`, solicitud a C1 en §15).
- **Fórmula** (motor S2): ahorrado = Σ aportaciones − Σ retiros (fondo y metas); aportado del mes contra planificado del mes (Σ aportación mensual planificada activa); meta más próxima por fecha con % = ahorrado ÷ monto meta × 100, entero.
- **Periodo**: mes en curso (sin parámetros).
- **Contenido** (icono `piggy-bank`):
  - Título: «Ahorro y metas».
  - Valor principal: `[ahorrado]` con `color.role.savings.*` y la etiqueta «ahorrado» (rol con icono, DESIGN §2.4).
  - Valor secundario (360): «Este mes [aportado] de [planificado]». Desde 412 suma la meta más próxima: «[nombre] · [n] % · [fecha meta]» e insignia de estado: «En curso» (`in_progress`; rol información, `info`), «Vencida» (`overdue`; rol advertencia, `triangle-alert`), «Cumplida» (`achieved`; rol éxito, `circle-check`).
- **Acción primaria única**: enlace «Ver ahorro» → P-09 Ahorro y metas.
- **Vacío** (sin fondo ni metas): «—» y «Aún no tienes fondo ni metas de ahorro»; la acción cambia a «Crear meta» → P-26.13 (Ajustes › Salud y metas de ahorro). Sin `next_goal` pero con fondo: se omite la línea de la meta.
- **Error**: «No se pudo cargar el ahorro. Revisa que estés en la red de casa.» · «Reintentar».
- **FIRME** desde la puerta 2 (bloqueo `docs/estado/bloqueos/S2-metas.md` RESUELTO): la tarjeta usa `progress_pct`, entero (Q2-17 C: entero en la tarjeta de Inicio y un decimal en Metas), y los tres estados de Q2-18 A (en curso, cumplida y vencida; no hay «atrasada»). La tarjeta no muestra la aportación necesaria, así que no depende del bloqueo abierto `S2-metas-vencida`.
- **Ejemplo**: marcadores; el oráculo no trae ahorro (M-07).

### 2.6 Operaciones de C1 del tablero (publicadas en C1 v1.5 por el agente `s1-c1-v15`)
Todo dinero en centavos enteros (`MoneyCents`), porcentajes enteros con ROUND_HALF_UP (protocolo §6), meses `Month` (`AAAA-MM`). «Ausente» = el campo no llega cuando falta el dato (sin datos no es cero). Ninguna cifra la calcula la interfaz.

**Operaciones nuevas**

| operationId | Ruta | Parámetros | Esquema de respuesta (200) y campos exactos | Motor |
|---|---|---|---|---|
| `getDashboardDebtCapacity` | GET `/api/v1/dashboard/debt-capacity` | `period` (Month, opcional; por omisión el mes en curso) | `DashboardDebtCapacity`: `month` (Month); `income_cents` (MoneyCents, ausente sin quincena capturada ni esperada); `income_is_expected` (boolean); `extras_cents` (MoneyCents, informativo, no entra al cociente); `payments_cents` (MoneyCents, parte del hogar de los pagos del mes); `payments_total_cents` (MoneyCents, con la parte de Karen); `payment_count` (integer); `ratio_pct` (integer, ausente sin `income_cents` o con 0); `ratio_total_pct` (integer, con `payments_total_cents`; ausente igual); `warning_pct` (integer, `x-app-setting` `debt.capacity_warning_pct`); `limit_pct` (integer, `x-app-setting` `debt.capacity_limit_pct`); `scale_max_pct` (integer); `status` (`normal`, `warning`, `exceeded`; ausente sin `ratio_pct`); `headroom_cents` (MoneyCents con signo; ausente sin `income_cents`); `has_credits` (boolean) | S2 (ingreso) y S3 (pagos) |
| `getDashboardMonthlyOverview` | GET `/api/v1/dashboard/monthly-overview` | `from` y `to` (Month, obligatorios; `from` ≤ `to`, máximo 24 meses; si no, 422 `validation_failed`) | `DashboardMonthlyOverview`: `warning_pct`, `limit_pct` (integer, mismas claves); `months[]`: `month` (Month); `is_projection` (boolean: mes posterior al en curso); `income_total_cents`, `income_fortnights_cents`, `income_extras_cents` (MoneyCents; ausentes sin ingreso capturado ni esperado); `income_is_expected` (boolean); `accrued_cents` (MoneyCents, devengado del hogar, M-01 real; ausente en proyección y en meses sin datos); `cash_cents` (MoneyCents, efectivo del hogar, M-11; ausente igual); `planned_cents` (MoneyCents, plan del mes; ausente sin plan); `credit_payments_cents` y `credit_payments_total_cents` (MoneyCents); `debt_ratio_pct` (integer, ausente sin ingreso); `debt_status` (`normal`, `warning`, `exceeded`; ausente sin `debt_ratio_pct`) | S2 y S3 |
| `getBudgetSignals` | GET `/api/v1/budget/signals` | `period` (Month, opcional; por omisión el mes en curso) | `BudgetSignals`: `month` (Month); `months_considered` (integer, `x-app-setting` `budget.signal_months`); `margin_pct` (integer, `x-app-setting` `budget.signal_margin_pct`); `rows[]`: `category_id` (integer); `signal` (`raise`, `cut`, `none`); `months[]`: `month`, `planned_cents` (ausente sin plan), `actual_cents`, `used_pct` (ausente sin plan) | S2 |
| `getDashboardSettings` y `updateDashboardSettings` | GET y PATCH `/api/v1/settings/dashboard` | PATCH: `parameters` y `reset_parameters`, como `updateAlertSettings` (K44, C1 §11.7) | `DashboardSettings`: `parameters` con `debt.capacity_warning_pct` (30), `debt.capacity_limit_pct` (40), `dashboard.payment_reminder_count` (5), `budget.signal_months` (3) y `budget.signal_margin_pct` (10), cada uno `SettingParameterInteger`; validación: 1 ≤ warning < limit ≤ 100, 1 ≤ N ≤ 20, 2 ≤ meses ≤ 12, 1 ≤ margen ≤ 50 (rangos **PROVISIONALES**) | — |

**Campos nuevos en operaciones existentes** (todos opcionales en el esquema, para no romper P-06, P-08 ni P-17):

| Operación · esquema | Campos nuevos |
|---|---|
| `getDashboardMonthVsPlan` · `DashboardMonthVsPlan` | `cash_cents` (MoneyCents, efectivo del hogar del mes = M-11 `cash_household_cents`); `previous_month` (Month); `previous_actual_cents`, `previous_cash_cents` (MoneyCents, ausentes sin datos de ese mes); `actual_difference_cents`, `cash_difference_cents` (MoneyCents con signo, mes − mes anterior; ausentes si falta el mes anterior); `categories_warning_count`, `categories_exceeded_count` (integer) |
| `getIncomeMonth` · `IncomeMonth` | `previous_month` (Month); `previous_fortnights_cents`, `previous_extras_cents` (MoneyCents, ausentes sin datos); `fortnights_difference_cents` (MoneyCents con signo) |
| `getCreditTotals` · `CreditTotals` | `month_payment_count` (integer); `month_pending_household_cents` (MoneyCents, parte del hogar sin pagar de los pagos del mes) |

**Claves nuevas de `app_setting`** (esquema 1.5, migración 00005 con `@param`; K54 y K44): `debt.capacity_warning_pct` = 30 y `debt.capacity_limit_pct` = 40 (FIRMES, K54 decisión 2); `dashboard.payment_reminder_count` = 5, `budget.signal_months` = 3 y `budget.signal_margin_pct` = 10 (**PROVISIONALES**, bloqueos de §2). Dónde se editan en Ajustes lo decide U4 (§15).

**Métricas de C1-datos §9 (publicadas en C1 v1.5)**: M-18 «Ingreso del mes» (K-2), M-19 «Créditos del mes» (K-3: pagos del mes, cuántos y falta pagar), M-20 «Capacidad de endeudamiento» (B-1), M-21 «Resumen mensual del tablero» (G-1 y G-6) y M-22 «Sugerencia de presupuesto» (G-3), con las fórmulas de sus bloques; las reglas del barómetro y de la sugerencia son RN-30 y RN-31 de `reglas-negocio.md`.

---

## 3. Detalle

**Propósito**: el segundo nivel de los bloques del tablero (D-1 a D-6 siguen el orden de las 6 tarjetas de U3 1.x), sin salir de Inicio (U1 N-027), más el bloque «Participación de Karen» (U1 N-026). Responde «¿de dónde sale la cifra de la tarjeta?» con lo mínimo; lo completo vive en la pantalla destino.
**Contenido**: 7 paneles en el orden de las tarjetas (D-1 a D-6) y al final D-7. Cada panel es un bloque con título (`font.size.title-3`), contenido y un enlace «Ver todo» a su pantalla destino (enlace, no botón con relleno). Componentes: «Tabla compacta» y «Lista de propiedades» (DESIGN §7).
**Filtros**: barra de vista de Detalle (en < 768 px, botón «Filtrar» con icono `list-filter` que abre una hoja inferior y chips con los filtros activos bajo el título, U1 §1.3). Cada panel aplica solo los filtros que su operación admite (§6); un filtro que un panel no admite **no** lo cambia y el panel lo dice con la línea «Este bloque no usa el filtro [nombre]» (`font.size.caption`).
**Periodo por defecto**: mes en curso; la quincena es subperiodo (selector «Mes · 1.ª quincena · 2.ª quincena», valores `period` = `2026-10`, `2026-10-q1`, `2026-10-q2`; K8).
**Estados**: los de §1.3 por panel (cada panel hace su petición), con los textos de §7.

| Panel | Contenido | Métrica · operación C1 (campos) | Vacío |
|---|---|---|---|
| D-1 Presupuesto por categoría | Las 5 categorías con mayor % usado: nombre, barra horizontal real contra plan (patrón «barras para comparar categorías»), «[real] de [plan] · [n] %» e insignia de estado como en K-1 | M-09 (N-081) · `getBudgetByCategory` (`month`, `rows[]`: `category_id`, `planned_cents`, `actual_cents`, `difference_cents`, `used_pct`, `status`, `expense_count`); nombres con `listCategories` | «Aún no hay presupuesto por categoría» · «Capturar presupuesto» → P-26.02 |
| D-2 Quincena en curso | Fila de la quincena en curso como lista de propiedades: ingresos (marca «esperado» si `income_is_expected`), pagos de crédito (parte del hogar), fijos y venideros, gasto variable estimado, ahorro planificado, saldo de la quincena y saldo acumulado | M-10 (N-090) · `getCashFlowFortnights` (`rows[]`: `start_date`, `end_date`, `pay_date`, `income_cents`, `income_is_expected`, `credit_payments_cents`, `fixed_and_upcoming_cents`, `variable_estimate_cents`, `planned_savings_cents`, `balance_cents`, `cumulative_balance_cents`, `is_actual`). Gasto variable estimado: **PROVISIONAL** (`docs/estado/bloqueos/S2-flujo-quincenal.md`, Q2-12, abierta); ahorro planificado: el reparto mitad y mitad es FIRME (Q2-13 B: HALF_UP en la quincena 1 y el resto en la 2); el valor es **PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/S2-ahorro-quincenal.md`: falta la base, meta por meta o la suma del mes, qué metas cuentan y si resta del saldo). El motor S2 lo entrega en `null` y no lo resta de `balance_cents` ni de `cumulative_balance_cents` (`docs/contratos/motores/S2-interfaz.md` 1.3 §3.11); mientras llegue `null`, la fila muestra «—», nunca un monto en cero (regla «Sin datos no es cero»). C1 v1.3 declara `planned_savings_cents` entero obligatorio; solicitud a C1 en §15 | «Falta tu ingreso quincenal» · «Capturar ingreso» → P-26.03 |
| D-3 Pagos y disponible | (a) Pagos de todas las líneas en dos grupos, cada uno en orden de fecha: primero **los vencidos** (del más antiguo al más reciente) y luego **los próximos 5 pendientes** de hoy a hoy + 60. Cada fila: fecha con día de la semana («Vie 30 oct 2026»), línea, monto total, «Hogar [parte] · Karen [parte]» y estado de §1.4 («Vencido» con «Sin marcar como pagado», «Por vencer», «Parcial» con «Falta la parte de Karen», o «Pendiente»; «Pagado» no llega a este panel). Si hay más de 5 vencidos, la línea «Hay más pagos vencidos» y el enlace «Ver todo». (b) Disponible de cada línea revolvente con **las dos lecturas lado a lado** («Lectura A [monto]» · «Lectura B [monto]», la activa con la insignia «Activa») y el mínimo del horizonte de cada lectura; aviso «Reposición de capital pendiente de confirmar con Kueski» mientras no se confirme | M-13 (N-192) · `listScheduledPayments` (`id`, `due_date`, `amount_cents`, `credit_id`, `shares[]`: `person_id`, `amount_cents`, `paid_cents`, `status`; `status`) en **dos peticiones**, las dos con orden `due_date`, `limit` = 5 y los filtros de persona, familia y crédito de la barra (§6): (1) **vencidos**: `status` = `overdue` y `period` = `[fecha del vencido más antiguo]..[ayer]` (rango inclusivo de C1), donde esa fecha es `next_payment.due_date` de `getDashboardNextPayment` cuando `is_overdue` = true (M-03 pone primero el vencido más antiguo; se reutiliza la respuesta de K-3 o se pide); con `is_overdue` = false no hay vencidos y esta petición no se hace; (2) **próximos**: `status` = `pending` y `period` = `[hoy]..[hoy + 60]` (M-13 «hoy..+60»). Días de «Por vencer»: `getAlertSettings` (§1.2) · M-04 (N-285, N-201) · `previewAvailability` por línea revolvente | «Sin pagos pendientes» (b se sigue mostrando); sin créditos: «Aún no hay créditos» · «Agregar crédito» → P-26.05 |
| D-4 Comidas de hoy | Casillas de hoy por tiempo (desayuno, comida, cena, colación) y persona (Axel, Ximena): receta o «Fuera», kcal y costo | N-121 · `getMealPlanWeek` (`cells[]`: `plan_date`, `meal_slot_id`, `person_id`, `recipe_id`, `is_eating_out`, `cost_cents`, `kcal`), con `period` = semana en curso; la PWA muestra las casillas cuyo `plan_date` es hoy. **Métrica pendiente en C1-datos §9** (solicitud a C1, §15); mientras tanto rige esta definición: «Comidas de hoy» (N-121 vía N-027) = casillas de `cells[]` con `plan_date` = hoy, por tiempo y persona, con `recipe_id` o `is_eating_out`, `kcal` y `cost_cents` tal como llegan (sin sumas); periodo: la semana ISO de lunes a domingo que contiene hoy (K8; también en domingo, porque el panel muestra hoy y K38 solo cambia la semana que abre P-10); filtro: `person_id` (S) | «Hoy no hay comidas en el plan» · «Armar plan» → P-10 (regla K38) |
| D-5 Diferencias de precio | Los 5 artículos con mayor diferencia de precio entre tiendas de la canasta de la tarjeta 5: artículo, tienda ganadora (punto + nombre), diferencia, marca «viejo» si el precio es viejo; y el total por tienda (orden fijo de tiendas, DESIGN §2.6) | N-162 · `getBasketComparison` (`stores[]`: `store_id`, `total_cents`, `difference_vs_cheapest_cents`, `stale_count`, `assumes_membership`; renglones con `product_id`, `winner_store_id`, `difference_cents`, `prices[]`.`is_stale`), con el mismo `week` (y `basket`, `expense_id`) de la tarjeta 5. **Métrica pendiente en C1-datos §9** (solicitud a C1, §15); mientras tanto rige esta definición: «Diferencias de precio» (N-162 vía N-027) = los 5 renglones con mayor `difference_cents`, en orden descendente (la interfaz ordena lo que está a la vista, U1 §0.2), y `stores[]` con `total_cents` tal como llega; periodo: la lista vigente (`week` según K38, como la tarjeta 5); filtro: tienda en la vista (V) | «Aún no hay precios para comparar» · «Registrar precios» → P-15 |
| D-6 Metas | Cada meta activa: nombre, barra de avance (`color.role.savings.*` con el % en texto, con un decimal), «[ahorrado] de [meta]», aportación mensual necesaria, fecha meta y estado (insignias de la tarjeta 6) | M-16 (N-111) · `listSavingsGoals` (`name`, `target_cents`, `target_date`, `saved_cents`, `progress_tenths_pct`, `required_monthly_cents`, `status`), con `status` activo. Avance con un decimal: `progress_tenths_pct` (décimas de punto porcentual con ROUND_HALF_UP, C1 v1.3), por Q2-17 C («un decimal en Metas y entero en la tarjeta de Inicio»): D-6 es la métrica M-16 «Metas» y el entero (`progress_pct`, M-07) queda solo en la tarjeta 6; sin `target_cents` el campo llega ausente y se muestra «—». C1-datos M-16 todavía dice «% de avance entero»; solicitud a C1 en §15. Aportación mensual necesaria con los meses del mes de hoy al de la fecha meta, los dos incluidos (Q2-16 A), y estado con los tres valores de Q2-18 A (`in_progress`, `achieved` y `overdue`; insignias de la tarjeta 6); bloqueo `docs/estado/bloqueos/S2-metas.md` RESUELTO. La aportación de una meta vencida sin cumplir es **PROVISIONAL** (`docs/estado/bloqueos/S2-metas-vencida.md`, abierto) | «Aún no tienes metas de ahorro» · «Crear meta» → P-26.13 |
| D-7 Participación de Karen | Ver §3.1 | M-08 (N-026) · `getPersonParticipation` | Ver §3.1 |

Ejemplos trazables de Detalle: D-3 (a) con hoy 10-oct-2026, rango de hoy a hoy + 60 (hasta el 9-dic-2026, inclusivo): «Vie 30 oct 2026», «Dom 15 nov 2026» y «Lun 30 nov 2026», cada uno por **$973.75** con «Hogar $598.82 · Karen $374.93» y estado «Pendiente» (R-13 a R-15, M-13); no hay vencidos (`getDashboardNextPayment` trae `is_overdue` = false) y el pago del 15-dic-2026 (R-16) cae fuera del rango. D-3 (b) el 10-oct-2026: «Lectura A **$10,000.00**» · «Lectura B **$6,905.50**» (R-64); mínimo del horizonte **$8,094.50** y **$5,000.00** (R-61, R-62). D-5 con la compra semilla: Bodega Aurrera **$3,094.50** (R-01) y Walmart **$3,233.00** (`totales.walmart`), diferencia **$138.50** (K-03).

Al pie de Detalle, enlace «Ver todos los gastos con estos filtros» → P-02, que recibe persona, familia, categoría, tienda, periodo y crédito (`getExpenseSummary` y `listExpenses` los admiten todos).

### 3.1 Participación de Karen (D-7)
- **Métrica C1**: M-08 (N-026) · `getPersonParticipation` (GET `/api/v1/persons/{person_id}/participation`) · esquema `PersonParticipation`: `month`, `accrued_groceries_cents`, `purchases[]` (`description`, `purchase_date`, `total_cents`, `payments_made`, `payments_total`, `balance_cents`), `payment_shares[]` (`due_date`, `amount_cents`, `status`), `live_debt_cents`. El `person_id` es el de la persona con `participation` externa en `listPersons`; con otra persona C1 responde 404.
- **Periodo**: mes en curso.
- **Contenido**, título «Participación de Karen» con punto `color.person.karen.*` y el nombre en texto:
  1. «Su parte de la despensa este mes: [monto]».
  2. «Sus compras a crédito»: tabla compacta con comercio, fecha, total, «Pagos [hechos] de [total]» y saldo.
  3. «Sus pagos de este mes»: fecha con día de la semana, monto de su parte y estado («Pagado» con `circle-check`; «Pendiente»).
  4. «Lo que debe en la línea: [monto]».
  5. Nota fija (texto exacto de U1 N-026): «Karen paga directo a Kueski su parte; su parte es informativa y repone la línea cuando paga».
- **Acción del panel**: enlace «Estado de cuenta de Karen» → P-25, donde se genera el PDF (`generateExternalStatement`) que **Axel reenvía a mano** (decisiones §A, F). Karen no tiene cuenta ni ve el dashboard.
- **Vacío**: «Karen no tiene compras ni pagos este mes»; la nota fija se mantiene.
- **Error**: «No se pudo cargar la participación de Karen. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable**: octubre de 2026, «Su parte de la despensa este mes: **$749.86**» (`periodos[0].devengado.karen`, M-08; P1 = oct-2026); su parte del pago del 30 de octubre, **$374.93** (R-11, R-13).

---

## 4. Tendencias

**Propósito**: ver cómo van las cifras de las tarjetas en el tiempo, con las cuatro series de U1 N-028.
**Contenido**: 4 gráficas en orden fijo T-1 a T-4, cada una con título, subtítulo de periodo, resumen de una frase debajo del título (DESIGN §7.7; su texto está en cada bloque), leyenda, etiquetas de dinero en texto y la alternativa «Ver como tabla» (enlace que cambia la gráfica por una tabla con las mismas cifras; la tabla es también el contenido para lectores de pantalla). Componentes: «Gráfica» y «Tabla compacta» (DESIGN §7) con los patrones de DESIGN §8. Biblioteca: ECharts 6 (decisiones §I).
**Filtros y periodos**: cada gráfica tiene su ventana fija de C1 (12 meses, 12 semanas); los filtros que admite están en su bloque y en §6.
**Estados**: §7.
**Reglas de todas las gráficas** (DESIGN §2.7, §2.10; K27, K41):
- Series en **orden fijo**; el color de una serie nunca va solo: lleva **forma de marcador**, **patrón de trazo o relleno** y **etiqueta directa** en texto.
- Paleta de series `color.chart.1` a `color.chart.8` (`.light` / `.dark`) en ese orden; las gráficas por persona usan la categórica `color.person.*` (DESIGN §2.6 regla 2); personas y tiendas nunca comparten gráfica ni leyenda.
- Etiquetas de dinero siempre en texto, nunca solo al tocar, con el **etiquetado selectivo de DESIGN §8.5** (nunca un número en cada punto): columnas con el valor sobre la columna si cabe con aire; líneas con el último valor y los extremos; área con los valores finales. Formato completo en etiquetas y leyendas; el compacto solo en ejes. Una etiqueta que no cabe no se recorta: se mueve fuera de la marca o pasa a «Ver como tabla». El valor al tocar (< 768) o al pasar el puntero (≥ 768) es un extra (U1 §6.2).
- Marcas de DESIGN §8.1: barras y capas separadas por 2 px del color de la superficie, **nunca un borde dibujado alrededor de una marca**; un solo eje Y; textos en `ui.text.default` o `ui.text.secondary`, nunca en el color de la serie; sin retícula vertical salvo «Hoy».
- Ejes y rejilla `color.neutral.700` / `color.neutral.300` para texto (`font.size.caption`) y divisores decorativos `color.neutral.200` / `color.neutral.800`; cifras con `tnum`.
- Sin donas en Dashboard. Si una implementación futura la necesitara, solo para composición de ≤ 5 partes (U2-13); hoy: 0 donas.
- Los meses sin cifra del oráculo no llevan marca de canon; los meses 7 a 12 se rotulan como proyección (N-093).
- **Resumen y `<figcaption>`** (DESIGN §7.7): cada gráfica va en `<figure>` con `<figcaption>`, que trae el título y el resumen de una frase de su bloque. El resumen va debajo del título, en `ui.text.secondary`. Solo enuncia cifras que llegan de C1 (o sumas de lo que está a la vista), con su mes o su semana. **No calcula variaciones ni porcentajes** («subió 12 %»), porque la interfaz no recalcula (U1 §0.2). En eso U3 se aparta del ejemplo de DESIGN §7.7 (solicitud 5 a DESIGN, §15). Si falta una de sus cifras, esa parte dice «—» con el motivo. En carga, vacío o error no hay resumen: manda el estado de §7. Con «Ver como tabla», el resumen se queda.

### T-1 · Plan contra real por mes
- **Métrica C1**: M-14 (N-085) · `getBudgetMonthlyTrend` (GET `/api/v1/budget/monthly-trend`) · `months[]`: `month`, `planned_cents`, `actual_cents`; filtro `category_id`.
- **Fórmula y periodo**: por mes (máx. 12): plan total y real total (parte del hogar devengada), o de una categoría.
- **Patrón**: **barras agrupadas** por mes (comparar dos cantidades en cada mes). Orden fijo de series: 1 «Real» = `color.chart.1`, relleno sólido; 2 «Plan» = `color.chart.2` con el patrón de la serie 2 (rayas a 45° cada 6 px, DESIGN §8.6) y **sin contorno** (DESIGN §8.1); las dos barras de cada mes se separan 2 px con el color de la superficie. Leyenda: «Real», «Plan».
- **Subtítulo**: «Últimos 12 meses · [categoría o «Todas las categorías»]».
- **Resumen** (DESIGN §7.7): «En [último mes]: real [monto] de un plan de [monto]», con `actual_cents` y `planned_cents` del último mes de `months[]`. Sin plan en ese mes: «En [último mes]: real [monto]; aún sin plan».
- **Vacío**: «Aún no hay presupuesto capturado» · «Capturar presupuesto» → P-26.02. **Ejemplo**: marcador (el plan lo captura Axel; M-14).

### T-2 · Devengado contra efectivo
- **Métrica C1**: M-11 (N-091) · `getAccrualVsCash` (GET `/api/v1/cash-flow/accrual-vs-cash`) · `projection_mode`, `months[]`: `month`, `accrued_cents`, `cash_cents`, `accrued_household_cents`, `accrued_external_cents`, `cash_household_cents`, `cash_external_cents`, `cumulative_difference_cents`, `is_canonical` y, desde C1 v1.3, `cumulative_accrued_cents` y `cumulative_cash_cents` (Q2-65 A).
- **Fórmula y periodo**: mes en curso a 11 meses después; diferencia acumulada = devengado acumulado − efectivo acumulado (motor S2). **No se rotula «deuda viva»** (K19, N-091).
- **Patrón**: **área apilada devengado contra efectivo**, acumulada. Orden fijo de capas, de abajo arriba: 1 «Efectivo acumulado» = `color.chart.1`, relleno sólido, borde superior continuo con marcador círculo; 2 «Diferencia acumulada» (la capa «devengado sin pagar» de DESIGN §8.3; el rótulo sigue a C1 M-11) = `color.chart.2`, relleno con trama diagonal, borde superior discontinuo con marcador cuadrado. La orilla superior de la pila es el devengado acumulado y lleva la etiqueta «Devengado acumulado». **Las dos capas tienen el mismo origen** que `cumulative_difference_cents`, que C1 define como `periodos[i].diferencia_acumulada` de `model.json`: se acumula desde P1 aunque `months[]` empiece en el mes en curso. Capa 1 = efectivo acumulado desde ese origen; cima = devengado acumulado desde ese origen; capa 2 = cima − capa 1 = `cumulative_difference_cents` tal como llega. **Prohibido** sumar `cash_cents` desde el primer mes visible y apilar encima `cumulative_difference_cents`: mezcla dos orígenes y, desde el segundo mes, la cima deja de ser el devengado acumulado.
- **Fuente de las capas** (Q2-65 A; bloqueo `docs/estado/bloqueos/U3-t2-acumulados.md` RESUELTO): C1 v1.3 agrega en `months[]` `cumulative_accrued_cents` y `cumulative_cash_cents`, con el mismo origen que `cumulative_difference_cents` (así `cumulative_difference_cents` = `cumulative_accrued_cents` − `cumulative_cash_cents`). Capa 1 = `cumulative_cash_cents`; cima = `cumulative_accrued_cents`; capa 2 = `cumulative_difference_cents`. Etiquetas «Devengado acumulado» y «Efectivo acumulado». Se retira la variante (b), relativa a la ventana.
- **Pago adelantado** (DESIGN §8.3): si un mes trae `cumulative_difference_cents` negativo (se pagó antes de devengar), la capa 2 **vale cero en ese tramo** y el texto «Pagado por adelantado» lo dice junto al tramo; nunca se dibuja una capa negativa y la gráfica sigue siendo un área apilada (en ese tramo la cima de la pila es el efectivo acumulado). La tabla muestra la diferencia tal como llega, con su signo menos.
- Línea vertical «Hoy» de 1 px en `ui.viz.baseline` con su etiqueta (DESIGN §8.3; única retícula vertical, §8.1).
- **Resumen** (DESIGN §7.7), con las cifras de las etiquetas finales: «Hasta [último mes], con proyección desde [séptimo mes]: devengado acumulado [monto] y efectivo acumulado [monto]». No nombra la diferencia como deuda ni como «debes» (K19, N-091, M-11).
- Etiquetas en texto al final de la pila (DESIGN §8.3 y §8.5): «Devengado acumulado» y «Efectivo acumulado», con su valor final (`cumulative_accrued_cents` y `cumulative_cash_cents` del último mes). El subtítulo da solo el periodo y **no** rotula la deuda viva (K19, N-091, M-11): U3 se aparta aquí del ejemplo de subtítulo de DESIGN §8.3; solicitud 3 a DESIGN (§15).
- Meses 7 a 12: trazo discontinuo y la nota de lo que incluye la proyección según `projection_mode` (K21, N-093). Con el interruptor encendido, esos meses repiten las despensas con el patrón del oráculo como valor inicial sembrado, que cambia con partidas planificadas (Q2-21 a Q2-24 C; bloqueo `docs/estado/bloqueos/S2-proyeccion-interruptor.md` RESUELTO), y la nota lo dice: «Proyección con despensas repetidas»; con el interruptor apagado (por omisión) no se repiten. La tabla trae también las partes del hogar y de Karen de cada mes (`accrued_household_cents`, `accrued_external_cents`, `cash_household_cents`, `cash_external_cents`); la de Karen es informativa, de la despensa en común (§9).
- **Vacío**: «Aún no hay gastos ni pagos para comparar» · «Registrar gasto» (acción de página). **Ejemplo trazable** (oct-2026 = P1 y mar-2027 = P6, regla de ejemplos al inicio; C1-datos M-01, `is_canonical`): fila oct-2026, devengado **$1,947.50**, efectivo **$973.75** y diferencia acumulada **$973.75** (R-45, M-11); fila nov-2026, **$1,947.50**, **$1,947.50** y **$973.75** (R-46). Capas, que no cambian con el mes en curso porque el origen no se mueve: en oct-2026, efectivo acumulado **$973.75** y devengado acumulado **$1,947.50** (R-45); en nov-2026, efectivo acumulado **$2,921.25** (= 973.75 + 1,947.50, efectivo de R-45 y R-46) y devengado acumulado **$3,895.00** (= 1,947.50 + 1,947.50, devengado de R-45 y R-46); en mar-2027, **$13,103.78** (R-52) y **$14,362.48** (R-51), con diferencia acumulada **$1,258.70** (R-50). **Fila «Total» de la tabla**: suma de las 12 filas visibles (la interfaz suma lo que está a la vista, U1 §0.2), sin cifra de canon: R-51 y R-52 son los totales de P1 a P6 (seis filas), no de doce.

### T-3 · Deuda viva proyectada
- **Métrica C1**: M-17 (N-194) · `getDebtTrend` (GET `/api/v1/credits/debt-trend`) · `months[]`: `month`, `live_debt_cents`, `credits[]` (`credit_id`, `live_debt_cents`), `interest_paid_cents`, `interest_pending_cents`; filtro `credit_id`.
- **Fórmula y periodo**: por mes (12), deuda viva al cierre por línea; intereses pagados y por pagar (motor S3).
- **Patrón**: **líneas** (tendencia). Sin filtro: una línea «Deuda viva total» = `color.chart.1`, continua, marcador círculo. Con filtro de crédito: una línea por línea elegida en el orden de `listCredits`, series `color.chart.1`, `color.chart.2`, … con marcadores círculo, cuadrado, triángulo, rombo y trazos continuo, discontinuo, punteado, en ese orden. Intereses pagados y por pagar solo en la tabla.
- **Resumen** (DESIGN §7.7): sin filtro, «Deuda viva al cierre de [primer mes]: [monto]; al cierre de [último mes]: [monto]», con `live_debt_cents` del primer y del último mes de `months[]`. Con filtro de crédito: «Al cierre de [último mes]: [línea] [monto] · [línea] [monto]», con `credits[]` del último mes en el orden de la leyenda. Aquí sí se rotula la deuda viva, porque es la métrica de la gráfica (M-17).
- **Vacío**: «No hay deudas vivas» (con `circle-check`). **Ejemplo trazable**: la cifra de control del motor es la deuda viva al 10-abr-2027, **$6,293.44** (R-53 = **$1,258.70** R-54 + **$5,034.74** R-55, M-17); los puntos mensuales son «al cierre» de cada mes y los entrega C1 tal cual.

### T-4 · Costo semanal y kcal de comidas
- **Métrica C1**: M-15 (N-129) · `getMealPlanTrends` (GET `/api/v1/meal-plan/trends`) · `weeks[]`: `week_start`, `cost_cents`, `persons[]` (`person_id`, `avg_kcal_per_day`, `daily_kcal_target`); filtro `person_id`.
- **Fórmula y periodo**: por semana (12, lunes a domingo, K8): costo de casillas con receta; kcal promedio por día por persona contra su meta (motor S4).
- **Patrón**: dos gráficas de **líneas** apiladas en vertical, mismo eje de semanas:
  - T-4a «Costo por semana»: una línea `color.chart.1`, continua, marcador círculo.
  - T-4b «kcal al día por persona»: orden fijo Axel, Ximena; «Axel» = `color.person.axel.*`, continua, círculo; «Ximena» = `color.person.ximena.*`, trazo 8 · 4 y marcador cuadrado (la segunda persona usa el marcador y el trazo de la serie 2, DESIGN §8.7); la meta de cada persona, línea del mismo color con trazo de puntos 2 · 4, sin marcadores y con la etiqueta «Meta de [nombre]» (DESIGN §8 no fija la línea de meta: supuesto 8 y solicitud 4 a DESIGN, §15). Etiqueta directa con el nombre al final de cada línea (≤ 4 elementos, DESIGN §2.6 regla 2).
- **Resumen** (DESIGN §7.7), uno por gráfica, con la última semana de `weeks[]`: T-4a «Semana Lun [d] al Dom [d mmm]: costo [monto]» (`cost_cents`); T-4b «Semana Lun [d] al Dom [d mmm]: Axel [n] kcal al día con meta de [n]; Ximena [n] kcal al día con meta de [n]» (`avg_kcal_per_day`, `daily_kcal_target`).
- **Vacío**: «Aún no hay semanas con plan de comidas» · «Armar plan» → P-10 (regla K38). **Ejemplo**: marcadores (M-15).

---

## 5. Calendario

**Propósito**: ver en un solo lugar qué entra y qué sale por día: pagos del calendario unificado de créditos y eventos de despensa, con las quincenas marcadas.
- **Métrica C1**: M-13 (N-029) · `listCalendarEvents` (GET `/api/v1/calendar/events`) · `CalendarEventPage` de `CalendarEvent`: `event_type`, `date`, `end_date` (solo promociones), `title`, `amount_cents`, `household_share_cents`, `external_share_cents`, `status` (`pending`, `paid`, `overdue`, `expected`, `captured`, `complete`, `incomplete`, `empty`), `credit_id`, `ref_table`, `ref_id`. Filtros: `period`, `credit_id`, `family`, `event_type`.
- **Tipos que muestra Inicio** (U1 N-029, orden dentro del día): ingresos (`income`), pagos de crédito (`credit_payment`), gastos fijos (`fixed_expense`) y venideros (`upcoming_expense`), fechas de corte (`card_statement`), despensa programada (`planned_grocery`), días de visita a mercados (`market_visit`), promociones (`promotion`, franja de inicio a fin), fechas meta de ahorro (`savings_goal`) y envío del reporte mensual (`monthly_report`). `meal_plan` y `savings_movement` son de otros calendarios y no se piden aquí.
- **Iconos por tipo** (DESIGN §5.4): ingreso `arrow-down-left` con «+» y `color.role.income.*`; pago de crédito `credit-card`; fijo y venidero `arrow-up-right` con «−» y `color.role.expense.*`; corte `calendar-clock`; despensa programada `shopping-cart`; mercado `store` con punto `color.store.market.*`; promoción `calendar` con franja; meta `piggy-bank`; reporte `file-chart-column`. Estados de los pagos de crédito: los cinco de §1.4 (DESIGN §7.8): «Pendiente», «Por vencer», «Vencido» con «Sin marcar como pagado», «Parcial» con «Falta la parte de Karen» y «Pagado». Ingresos esperados: «Esperado» (`info`, rol información). «Parcial» no viene en `CalendarEvent.status`: el Calendario pide además `listScheduledPayments` con el mismo `period`, `credit_id` y `family` y `status` = `pending`, y une cada evento `credit_payment` cuyo `ref_table` es `scheduled_payment` con el pago cuyo `id` es `ref_id` (M-13 cita las dos operaciones). Si esa petición falla, cada pago conserva el estado de `CalendarEvent` y nada más cambia.
- **Pago de crédito**: «[línea] · [monto total]» y, debajo, «Hogar [parte] · Karen [parte]» (Karen solo en pagos de la línea donde ella compra).
- **Quincenas (K8)**: la vista marca cada mes con dos franjas: «1.ª quincena · del 1 al 15» y «2.ª quincena · del 16 al [último día]»; el 31 cae en la segunda. El **día de pago** (15 y 30, o el último día del mes) es un dato aparte: lo marcan los eventos de ingreso y de pago con su fecha, no la franja. Febrero: segunda quincena del 16 al 28 (o 29) y pago el último día, p. ej. «Dom 28 feb 2027». La semana va de lunes a domingo (K8): la rejilla del mes empieza en lunes.
- **Vistas**:

| Ancho | Vista | Rango por defecto (`period`) |
|---|---|---|
| 360 y 412 | Agenda por día: encabezado de día «Vie 30 oct 2026», eventos en el orden de tipos; separador al cambiar de quincena; en 412, mes compacto opcional arriba (puntos por día) | de hoy a hoy + 30 días (rango inclusivo); «Ver 30 días más» amplía a hoy + 60 |
| 768 y 1024 | Mes: rejilla de lunes a domingo con las franjas de quincena; hasta 3 eventos por día y «+[n] más» | mes en curso (`2026-10`) |
| 1440 | Mes con el detalle del día seleccionado en panel derecho (lista completa del día con montos y partes) | mes en curso |

- **Filtros**: tipo de evento (selector múltiple «Tipos»), crédito, familia (K5: «Hogar» incluye lo conjunto; «Karen» incluye lo conjunto; «Conjunto» solo lo conjunto; sin valor: todos), periodo (navegación «‹ mes anterior · mes siguiente ›» y «Hoy»).
- **Acción del panel**: tocar un evento abre su registro (vista previa G-06: página completa en < 768 px, panel en ≥ 1024 px); un pago abre P-19 de su compra o P-18 de su línea.
- **Vacío**: «No hay eventos en este periodo»; si no hay ningún crédito, ingreso ni despensa programada: «Aún no hay pagos ni ingresos programados» · «Capturar ingreso» → P-26.03.
- **Ejemplo trazable**: octubre a diciembre de 2026, pagos de la despensa 1 el «Vie 30 oct 2026», «Dom 15 nov 2026», «Lun 30 nov 2026» y «Mar 15 dic 2026», cada uno por **$973.75** (R-13 a R-16), con «Hogar $598.82 · Karen $374.93».

---

## 6. Filtros y periodos por tarjeta y panel

Filtros de C1 (`componentes-comunes.yaml`): persona `person_id`, familia `family` (K5), categoría `category_id`, tienda `store_id`, periodo `period`, crédito `credit_id`. Claves: **S** = filtro en el servidor (la operación lo admite); **V** = filtra lo que está a la vista (U1 §0.2: la interfaz filtra y ordena lo que muestra, sin recalcular); **—** = no aplica (la operación no lo admite y la cifra no se puede filtrar sin recalcular; el bloque lo dice, §3). Periodo por defecto: **mes calendario**, con la quincena como subperiodo (decisiones §B).

| Tarjeta o panel | Persona | Familia | Categoría | Tienda | Periodo (por defecto) | Crédito |
|---|---|---|---|---|---|---|
| K-1 Gasto del mes | — | — | — | — | S · mes en curso | — |
| K-2 Ingreso del mes | — | — | — | — | S · mes en curso | — |
| K-3 Créditos | — | — | — | — | — · hoy y mes en curso | — |
| K-4 Pagos próximos | — | — | — | — | S · vencidos hasta ayer y de hoy a hoy + 30 (PROVISIONAL) | — |
| B-1 Barómetro | — | — | — | — | S · mes en curso | — |
| G-1 Ingreso, gasto y pagos | — | — | — | — | S · 11 meses anteriores y el en curso (PROVISIONAL) | — |
| G-2 Deuda viva por línea | — | — | — | — | S · hoy | — |
| G-3 Presupuesto por categoría | — | — | V | — | S · mes en curso (sugerencia: 3 meses cerrados) | — |
| G-4 = T-2 · G-5 = T-3 | — | — | — | — | — · las de T-2 y T-3 | — |
| G-6 Tendencia del barómetro | — | — | — | — | S · 6 anteriores, el en curso y 5 siguientes (PROVISIONAL) | — |
| T2 Flujo de la quincena (Más del hogar) | — | — | — | — | S · quincena en curso (K8) | — |
| T4 Comidas de la semana (Más del hogar) | — | — | — | — | S · semana en curso (K8) | — |
| T5 Mejor canasta (Más del hogar) | — | — | — | — | — · lista vigente (`week` según K38) | — |
| T6 Ahorro y metas (Más del hogar) | — | — | — | — | — · mes en curso | — |
| D-1 Presupuesto por categoría | — | — | V | — | S · mes (o quincena) | — |
| D-2 Quincena en curso | — | — | — | — | — · quincena en curso | — |
| D-3 Pagos y disponible | S | S | — | — | S · vencidos hasta ayer y de hoy a hoy + 60 | S |
| D-4 Comidas de hoy | S | — | — | — | S · semana en curso | — |
| D-5 Diferencias de precio | — | — | — | V | — · lista vigente (`week`) | — |
| D-6 Metas | — | — | — | — | — · hoy | — |
| D-7 Participación de Karen | — | — | — | — | — · mes en curso | — |
| T-1 Plan contra real | — | — | S | — | — · 12 meses | — |
| T-2 Devengado contra efectivo | — | — | — | — | — · mes en curso a + 11 | — |
| T-3 Deuda viva proyectada | — | — | — | — | — · 12 meses | S |
| T-4 Costo semanal y kcal | S | — | — | — | — · 12 semanas | — |
| C-1 Calendario | — | S | — | — | S · mes (agenda: hoy a + 30) | S |

Los bloques del tablero no llevan filtros de la barra (son la vista de la casa; G-3 solo ordena y recorta lo visible); el filtro vive en Detalle, Tendencias y Calendario. El enlace «Ver todos los gastos con estos filtros» de Detalle lleva los seis filtros a P-02.

---

## 7. Estados por tarjeta y panel

Reglas comunes en §1.3 (carga con esqueleto y aviso a los 10 s; error con «Reintentar»; sin conexión con lo último cargado y su hora). Celdas = texto exacto o regla.

| Bloque | Vacío (qué se ve · acción) | Carga | Error (texto · acción) | Sin conexión |
|---|---|---|---|---|
| K-1 | «Aún no hay presupuesto para [mes]» · «Capturar presupuesto»; sin mes anterior, «Sin datos de [mes]» | esqueleto de tarjeta con barra | «No se pudo cargar el gasto del mes. Revisa que estés en la red de casa.» · «Reintentar» | último dato + «Cargado a las [HH:MM]» |
| K-2 | «Falta tu ingreso quincenal» · «Capturar ingreso» | esqueleto de tarjeta | «No se pudo cargar el ingreso del mes. …» · «Reintentar» | ídem |
| K-3 | «Aún no hay créditos» · «Agregar crédito»; «Sin pagos en [mes]» | esqueleto de tarjeta + bloque de disponible | «No se pudieron cargar los créditos. …» · «Reintentar» | ídem |
| K-4 | «Sin pagos en los próximos 30 días» / «Aún no hay créditos» · «Agregar crédito» | 3 filas esqueleto | «No se pudieron cargar los pagos. …» · «Reintentar» | ídem |
| B-1 | «Falta tu ingreso del mes para calcular la capacidad» · «Capturar ingreso»; sin créditos, 0 % con «Sin créditos activos» | arco esqueleto en `ui.skeleton.on-card` sin aguja ni marcas falsas | «No se pudo calcular la capacidad de endeudamiento. …» · «Reintentar» | ídem |
| G-1 | «Aún no hay meses con ingreso ni gasto» · «Capturar ingreso» | ídem T-1 | ídem T-1 | ídem |
| G-2 | «No hay deudas vivas» (sin acción) | ídem T-1 | ídem T-1 | ídem |
| G-3 | «Aún no hay presupuesto por categoría» · «Capturar presupuesto» | ídem T-1 | ídem T-1; sin sugerencias si falla `getBudgetSignals` | ídem |
| G-4, G-5 | los de T-2 y T-3 | ídem T-1 | ídem T-1 | ídem |
| G-6 | «Falta tu ingreso para calcular la capacidad» · «Capturar ingreso» | ídem T-1 | ídem T-1 (comparte petición con G-1) | ídem |
| T2 | «Falta tu ingreso quincenal» · «Capturar ingreso» | esqueleto de tarjeta | «No se pudo cargar el flujo de la quincena. …» · «Reintentar» | ídem |
| T4 | «Aún no hay plan para esta semana» · «Armar plan» | esqueleto de tarjeta | «No se pudo cargar el plan de comidas. …» · «Reintentar» | ídem |
| T5 | «Aún no hay precios para comparar» · «Registrar precios» | esqueleto de tarjeta | «No se pudo comparar la canasta. …» · «Reintentar» | ídem |
| T6 | «Aún no tienes fondo ni metas de ahorro» · «Crear meta» | esqueleto de tarjeta | «No se pudo cargar el ahorro. …» · «Reintentar» | ídem |
| D-1 | «Aún no hay presupuesto por categoría» · «Capturar presupuesto» | 5 filas esqueleto | «No se pudo cargar el presupuesto por categoría. …» · «Reintentar» | ídem |
| D-2 | «Falta tu ingreso quincenal» · «Capturar ingreso» | lista esqueleto de 7 propiedades | «No se pudo cargar la quincena. …» · «Reintentar» | ídem |
| D-3 | «Sin pagos pendientes» · «Agregar crédito» si no hay créditos | 5 filas esqueleto + bloque de disponible | «No se pudieron cargar los pagos. …» · «Reintentar» | ídem |
| D-4 | «Hoy no hay comidas en el plan» · «Armar plan» | 4 filas esqueleto | «No se pudieron cargar las comidas de hoy. …» · «Reintentar» | ídem |
| D-5 | «Aún no hay precios para comparar» · «Registrar precios» | 5 filas esqueleto | «No se pudo comparar la canasta. …» · «Reintentar» | ídem |
| D-6 | «Aún no tienes metas de ahorro» · «Crear meta» | 3 filas esqueleto | «No se pudieron cargar las metas. …» · «Reintentar» | ídem |
| D-7 | «Karen no tiene compras ni pagos este mes» (nota fija visible) | bloque esqueleto | «No se pudo cargar la participación de Karen. …» · «Reintentar» | ídem |
| T-1 | «Aún no hay presupuesto capturado» · «Capturar presupuesto» | esqueleto del alto final en `ui.skeleton.on-card`, sin ejes falsos (DESIGN §7.7) | «No se pudo cargar la gráfica. …» · «Reintentar» | ídem |
| T-2 | «Aún no hay gastos ni pagos para comparar» · «Registrar gasto» | ídem T-1 | ídem T-1 | ídem |
| T-3 | «No hay deudas vivas» (sin acción) | ídem T-1 | ídem T-1 | ídem |
| T-4 | «Aún no hay semanas con plan de comidas» · «Armar plan» | ídem T-1 | ídem T-1 | ídem |
| C-1 | «No hay eventos en este periodo» / «Aún no hay pagos ni ingresos programados» · «Capturar ingreso» | rejilla o agenda esqueleto | «No se pudo cargar el calendario. …» · «Reintentar» | ídem; meses nunca cargados: «Sin conexión · Esta página aún no se ha cargado en este teléfono» |

«…» abrevia en esta tabla la segunda frase común: «Revisa que estés en la red de casa.» Los mensajes de error nunca muestran códigos ni texto en inglés; con 401 se va a P-00 con aviso (sesión vencida, U1 §2.4) y con 404 de un registro, «Ya no existe · Volver».

---

## 8. Diseño por tamaño

Móvil primero: se diseña a 360 px y se amplía; referencias de U5: Android 412 × 915 y escritorio 1440 × 900 (U1 §6.1). Rejilla de DESIGN §4.2. Ningún ancho tiene desplazamiento horizontal de la página; solo la fila de pestañas se desplaza. Todo objetivo táctil ≥ `size.touch.min` (44 × 44 px) con `space.2` entre objetivos contiguos.

| Elemento | 360 y 412 (4 columnas, margen `space.4`) | 768 y 1024 | 1440 (12 columnas, margen `space.12`) |
|---|---|---|---|
| Tablero (Resumen) | 1 columna en el orden de apilado de §2.4 (K-1, B-1, K-4, K-2, K-3, G-3, G-1, G-6, G-2, G-5, G-4, «Más del hogar»), relleno de tarjeta `space.4`. A 360: título, valor principal y 1 línea secundaria; a 412 suma 1 dato (U1 §6.2). La cifra principal y el dígito del barómetro nunca se truncan (U1 §6.1 regla 5) | 768: K-1 con B-1 y K-2 con K-3 (4 + 4), K-4 y cada gráfica a todo el ancho, «Más del hogar» 2 × 2 (relleno `space.6`); 1024: rejilla de 1440 con la barra lateral fija de `layout.sidebar.width` y «Más del hogar» 2 × 2 | K-1 (6) y B-1 (6); K-2, K-3 y K-4 (4 cada uno); gráficas de dos en dos (6 + 6); «Más del hogar» en 4 tarjetas de 3 columnas; medianil `space.6`; todas las líneas secundarias visibles |
| Acción primaria de página | «+» de la barra inferior (`size.bottombar.item`) | 768: «+ Registrar gasto» en la barra superior (`layout.topbar.height`); igual en 1024 | igual que 1024 |
| Pestañas | fila con desplazamiento horizontal, activa visible | fila completa | fila completa |
| Detalle | paneles a 1 columna; tablas compactas de 2 líneas (título y cifra; 1 propiedad debajo; a 412, 1 más); filtros en hoja inferior con chips | 768: 1 columna con tablas de hasta 5 columnas; 1024: 2 columnas de paneles | 2 columnas de paneles; D-7 a todo el ancho al final |
| Tendencias y gráficas del tablero | 1 gráfica por fila, alto reducido, leyenda debajo, máximo 6 series, valor al tocar; etiqueta del periodo actual en texto | 1 gráfica por fila con leyenda a la derecha y valor al pasar el puntero | 2 gráficas por fila (T-1 con T-2; T-3 con T-4) |
| Calendario | agenda de hoy a hoy + 30 (412: mes compacto opcional) | mes | mes con panel del día seleccionado |
| Tipografía | `font.size.title-2` cifra; `font.size.small` en `ui.text.secondary` para la etiqueta de tarjeta y el título de gráfica (DESIGN §7.6, §7.7); `font.size.title-3` título de panel de Detalle; `font.size.small` secundario; `font.size.caption` metadatos | igual | igual |
| Alto de controles | `size.control.coarse` (puntero grueso) | `size.control.coarse` si es táctil; `size.control.fine` con puntero fino | `size.control.fine` (puntero fino) |

Teléfono en horizontal: rige el ancho resultante (915 px usa la regla de 768, U1 §6.1 regla 7).

---

## 9. Participación de Karen (reglas)

- Karen **no tiene cuenta y no ve el dashboard** (decisiones §A). Solo Axel y Ximena, co-titulares, ven Inicio, y ven lo mismo (U1 §0.5).
- Aparece **solo en lo que comparte**: (a) la despensa en común y (b) sus compras con el crédito Kueski y sus pagos. En Inicio eso es: la parte de Karen en K-3 (pagos del mes), en K-4 (pagos próximos), en la línea informativa «Con la parte de Karen: [n] %» de B-1 (desde 768), en la tarjeta 2 cuando aplica, en D-3 y en los pagos del Calendario (incluido el estado «Parcial · Falta la parte de Karen», §1.4); el panel D-7; las columnas de su parte en la tabla de T-2 («Ver como tabla»: devengado y efectivo de su parte de la despensa en común, `accrued_external_cents` y `cash_external_cents`, informativas); el valor «Karen» del filtro de familia (K5), que incluye lo conjunto.
- Su parte es **informativa** para el flujo de Axel: no suma a K-1, K-2, al dígito de B-1, a G-1, G-3, G-6 ni a las tarjetas 2 y 6, que usan la parte del hogar (U1 §0.2).
- **No hay indicadores de las finanzas de su familia**: ni ingresos, ni presupuesto, ni ahorro, ni gráficas propias; las gráficas por persona (T-4) solo tienen a Axel y Ximena (plan de comidas y salud: decisiones §A).
- Su estado de cuenta PDF es una acción del hogar desde P-25 (`generateExternalStatement`) y **Axel lo reenvía a mano** (decisiones §A, F).

---

## 10. Color, tipografía y accesibilidad

- Solo tokens de C2 / DESIGN por nombre: fondo y tarjeta `color.neutral.50` / `color.neutral.100` (oscuro `color.neutral.950` / `color.neutral.900`); texto `color.neutral.900` / `color.neutral.700` (oscuro `color.neutral.100` / `color.neutral.300`); enlaces `color.brand.700` / `color.brand.300`; roles `color.role.{income,expense,savings,debt,warning,error,success,info}.{light,dark,surface-light,surface-dark}`; personas `color.person.{axel,ximena,karen,shared}.*`; tiendas `color.store.{walmart,bodega-aurrera,sams-club,market}.*`; series `color.chart.1` a `color.chart.8` y «Otros» `color.neutral.600` / `color.neutral.400`. Ningún color nuevo (DESIGN §2.10 regla 4).
- El color nunca es el único portador de significado: cifras con signo e icono (DESIGN §3.3), estados con icono y texto, series con forma, patrón y etiqueta, personas y tiendas con nombre (K41). Ingreso y éxito no comparten componente: ingreso solo en cifras con «+», éxito solo en insignias con palomita (DESIGN §2.4).
- Cada gráfica tiene «Ver como tabla» con las mismas cifras, orden fijo de series y etiquetas de dinero en texto (§4).
- Contraste: los pares de U3 son los de DESIGN §2.8 y §2.9 (202 de 202 cumplen) y, desde 2.0, los del barómetro de DESIGN §8.9 (`ui.role.success`, `ui.role.warning` y `ui.role.error` sobre `ui.surface.card` y `ui.viz.track`, y el texto de las insignias sobre su superficie: todos ≥ 3:1 de interfaz o ≥ 4.5:1 de texto en claro y en oscuro). El rojo del barómetro es la excepción de K54 a DESIGN §2.4 («error, solo en validación y avisos»).
- Movimiento: las cifras no se animan (DESIGN §6.4); con movimiento reducido las gráficas aparecen sin transición.
- Teclado (escritorio): las pestañas se recorren con flechas; cada bloque del tablero es enfocable en el orden de apilado de §2.4 en todos los anchos.

---

## 11. Paso a implementación (F2, WF-3)

### 11.1 Componentes de DESIGN que usa U3
«Tarjeta de métrica» (§1.2), «Barómetro» (DESIGN §8.9, B-1), «Pestañas», «Insignia de rol», «Enlace», «Botón primario» (solo «Registrar gasto»), «Botón secundario» («Reintentar»), «Esqueleto», «Franja sin conexión», «Barra de vista con filtros» y «Hoja inferior de filtros», «Chip de filtro», «Tabla compacta», «Lista de propiedades», «Barra de avance», «Gráfica» (barras agrupadas, columnas agrupadas, barras horizontales, barras horizontales en % del plan, área apilada, líneas con líneas de referencia) con «Ver como tabla», «Calendario» (mes y agenda) y «Vista previa de registro». Los nombres siguen DESIGN §7 y §8; donde U3 se aparta de su contenido lo dice en el lugar (destino de la tarjeta en §1.2, «Sin marcar como pagado» y prioridad de estados en §1.4, subtítulo de T-2, línea de meta de T-4b, resumen de las gráficas sin variaciones calculadas en §4) y lo pide a DESIGN en §15.

### 11.2 Criterios de aceptación para F2
Pruebas con Playwright 1.64 + axe-core, reloj simulado en `America/Mexico_City` y datos de la semilla y el oráculo.

| Id | Criterio | Cómo se prueba |
|---|---|---|
| A-01 | El Resumen es el tablero de §2: K-1 a K-4, B-1, G-1 a G-6 y «Más del hogar» con las tarjetas 2, 4, 5 y 6, en el orden de apilado de §2.4 a 360 y 412 y en la rejilla de §2.4 a 768, 1024 y 1440; el orden del foco es el de 360 en todos los anchos | contar bloques y leer su orden en el DOM en los 5 anchos |
| A-02 | Cada tarjeta llama solo a su operación de §2 con los parámetros indicados; Inicio lee `getAlertSettings` y `getDashboardSettings` una sola vez cada una (§1.2); D-3 hace sus dos peticiones de `listScheduledPayments` (§3) y el Calendario la suya junto a `listCalendarEvents` (§5) | registro de peticiones de red por tarjeta y panel |
| A-03 | Tarjeta 4 y K38: con reloj en Dom 11 oct 2026 la tarjeta pide `period` = `2026-W41` y «Ver plan» abre P-10 en `2026-W42`; con Mié 14 oct 2026, ambas en `2026-W42` | dos corridas con fecha simulada |
| A-04 | Tarjeta 5 pide `week` de la lista vigente (Dom 11 oct 2026 → `2026-W42`) | petición de red |
| A-05 | K-3 y K-4 con la semilla el 10-oct-2026: K-4 con la fila «Vie 30 oct 2026 · Kueski · $973.75 · Pendiente» y «Hogar $598.82 · Karen $374.93»; K-3 con «Pagos de octubre: $973.75 · 1 pago», «Disponible hoy $6,905.50 · Lectura B» e insignia «Pendiente confirmar»; con lectura A, $10,000.00 | captura y texto |
| A-06 | K-1 en octubre de 2026 (= P1, regla de ejemplos al inicio) muestra gastado $1,197.64, pagado en efectivo $598.82 y «Sin datos de septiembre»; en noviembre de 2026, «Igual que octubre» y «$598.82 más que octubre» en efectivo; sin presupuesto, «—» y «Capturar presupuesto» | texto con dos relojes simulados |
| A-07 | Cada bloque muestra vacío, carga, error y sin conexión con los textos de §7; el aviso de 10 s aparece en la carga lenta | respuestas simuladas (vacía, lenta, 500, sin red) |
| A-08 | Una sola acción con relleno de marca en la página («Registrar gasto»); cada tarjeta, un solo enlace | conteo de botones primarios y enlaces por tarjeta |
| A-09 | Objetivos táctiles ≥ 44 × 44 px y sin desplazamiento horizontal de página en los 5 anchos | medición de cajas y del ancho de documento |
| A-10 | La cifra principal más larga y el dígito del barómetro caben completos a 360 px | prueba de ancho de U1 §6.1 regla 5 |
| A-11 | Gráficas con orden fijo de series, marcador o patrón distinto por serie, etiquetas de dinero en texto y «Ver como tabla» con las mismas cifras; T-2 con reloj en oct-2026: fila oct-2026 con $1,947.50, $973.75 y $973.75 (R-45); T-2 con reloj en nov-2026 (la ventana empieza en nov-2026): fila nov-2026 con $1,947.50, $1,947.50 y $973.75 (R-46); capas (Q2-65 A) en nov-2026, efectivo acumulado $2,921.25 y devengado acumulado $3,895.00, y en mar-2027, $13,103.78 y $14,362.48 (R-52, R-51), las mismas con los dos relojes; la fila «Total» de T-2 suma las 12 filas visibles; con un mes de diferencia negativa la capa 2 vale cero con «Pagado por adelantado»; ninguna marca con contorno y ninguna gráfica con un número en cada punto; cada gráfica en `<figure>` con `<figcaption>` (título y resumen de §4), resumen sin porcentajes ni variaciones calculadas y, en T-2, sin la palabra «deuda» | inspección del DOM y de la tabla; T-2 con dos corridas de reloj simulado (oct-2026 y nov-2026) |
| A-12 | Calendario: franjas de quincena 1 al 15 y 16 al último día; febrero de 2027 con pago «Dom 28 feb 2027»; semana de lunes a domingo | fecha simulada en octubre de 2026 y febrero de 2027 |
| A-13 | Karen: solo en K-3, K-4, la línea informativa de B-1 (desde 768), D-3, D-7, Calendario, las columnas de su parte en la tabla de T-2 y el filtro de familia; ningún indicador de su familia | revisión de textos con «Karen» |
| A-14 | Formato: dinero con `tnum` y signo menos tipográfico; fechas `10-oct-2026` o con día de la semana; 0 textos en inglés | axe-core y búsqueda de textos |
| A-15 | Tema claro y oscuro con los tokens de §10; axe-core sin violaciones AA | axe-core en los dos temas |
| A-16 | D-3 el 10-oct-2026 con la semilla: 3 pagos («Vie 30 oct 2026», «Dom 15 nov 2026», «Lun 30 nov 2026»), ninguno vencido y sin petición de vencidos | petición de red y texto |
| A-17 | Estados de §1.4 con el pago del 30 de octubre: el Jue 29 oct 2026 «Vence mañana»; el Vie 30 oct 2026 «Vence hoy»; el Sáb 31 oct 2026 sin marcar, «Vencido» con «Sin marcar como pagado»; con solo la parte de Karen pagada, «Falta la parte del hogar»; sin respuesta de `getAlertSettings`, ningún «Por vencer» | reloj simulado y respuestas simuladas |
| A-18 | B-1 con respuestas simuladas: `ratio_pct` 30 → «Verde · dentro del rango» (`circle-check`); 31 y 40 → «Amarillo · cerca del tope» (`triangle-alert`); 41 → «Rojo · arriba del tope» (`circle-alert`); sin `income_cents`, dígito «—» y «Capturar ingreso»; `aria-valuetext` «Capacidad de endeudamiento [n] % de 40 %: [estado]»; con `warning_pct` = 25 en la respuesta, la marca del tramo se mueve a 25 % (ningún umbral constante) | respuestas simuladas y lectura del árbol de accesibilidad |
| A-19 | G-1 y G-6 salen de una sola petición a `getDashboardMonthlyOverview` (`from` = mes en curso − 11, `to` = mes en curso + 5); G-1 con 3 series en columnas, la nota fija de pagos y sin columna en meses sin datos; G-6 con las líneas de referencia en `warning_pct` y `limit_pct`, trazo 8 · 4 en proyección y la línea «Hoy» | registro de red e inspección del DOM |
| A-20 | G-3: barras en % del plan, marcas en `alerts.budget_warning_pct` y `alerts.budget_limit_pct` leídas de `getAlertSettings`, orden descendente por `used_pct`, insignia «Sugerencia: subir presupuesto» o «recortar» solo con `signal` = `raise` o `cut`; sin `getBudgetSignals`, barras sin insignias | respuestas simuladas |
| A-21 | Desglose de §2.2: barra de G-3 → P-02 con `category_id` y `period` → P-04 → P-15 con el producto → pestaña «Tendencias», en 4 toques | recorrido con la semilla |
| A-22 | Ningún bloque calcula porcentajes, diferencias ni estados: todos llegan de C1; con un campo nuevo de C1 v1.5 ausente, su línea no se pinta y el resto del bloque se ve | respuestas sin los campos nuevos |

---

## 12. Capturas para U5

Cada captura en claro y en oscuro, a **412 × 915** y **1440 × 900** (U1 §6.1), con reloj simulado:
1. Tablero con datos de la semilla el 10-oct-2026: K-1, K-3 y K-4 con cifras (A-05, A-06); K-2, B-1, G-1 y G-6 con su vacío por falta de ingreso; G-3 con su vacío si no hay plan.
2. Tablero con un ingreso **sintético** de prueba (marcado así en la captura) para ver el barómetro en verde, amarillo y rojo (A-18) y G-6 con proyección.
3. Tablero en carga (esqueletos), con error en un bloque y sin conexión con «· [n] en cola».
4. Tablero a 360 × 800 con la prueba de ancho de la cifra principal y del dígito del barómetro.
4b. «Más del hogar» el Dom 11 oct 2026 y el Mié 14 oct 2026 (tarjeta 4 y destino P-10, K38).
5. Detalle completo, incluido D-7 «Participación de Karen» con la nota fija.
6. Tendencias con T-1 a T-4 y una gráfica en «Ver como tabla».
7. Calendario: agenda (412) y mes con panel del día (1440) de octubre de 2026; mes de febrero de 2027.

---

## 13. Trazabilidad con los criterios U3

| Criterio | Dónde se cumple |
|---|---|
| U3-01 · tablero de control (K54) | §2 (K-1 a K-4, B-1, §2.1 y G-1 a G-6, §2.2, §2.3, §2.4), §2.5 con las tarjetas 2, 4, 5 y 6 (Tarjeta 4 con K38, U1 §0.2, casos domingo y entre semana; Tarjeta 5 con `second_store_id`), §11.2 A-01, A-03 y A-18 a A-21 |
| U3-02 · métricas C1 | §2 a §5 (M-xx, `operationId`, campos, fórmula y periodo en cada bloque); §2.6 con las operaciones y campos NUEVOS publicados en C1 v1.5 y las métricas M-18 a M-22 de C1-datos §9; verificación en el reporte |
| U3-03 · estados | §1.3, §7 (25 filas de bloques × vacío, carga, error, sin conexión), cada bloque de §2 con su vacío y su error |
| U3-04 · móvil primero | §2.4 (orden de apilado en 412 y rejilla en 1440), §8 (360/412, 768/1024, 1440), §1.2 (táctil 44 px) |
| U3-05 · Detalle, Tendencias, Calendario | §3, §4, §5 (propósito, contenido, métricas, filtros, estados); quincenas K8 y febrero en §5; 0 donas (§4); T-2 con las dos capas del mismo origen (Q2-65 A, acumulados de C1 v1.3) |
| U3-06 · Karen | §3.1, §9 |
| U3-07 · próximo pago y A/B | K-3 (lectura activa y «Pendiente confirmar»), K-4 (próximo pago del calendario unificado), D-3 |
| U3-08 · cifras trazables | ejemplos con R-nn o K-nn (oct-2026 a mar-2027 = P1 a P6, regla de ejemplos al inicio); marcadores sin número; verificación en el reporte |
| U3-09 · formatos | reglas al inicio, textos de §2 a §7 |
| U3-10 · filtros | §6 |
| U3-11 · tokens y gráficas | §1.2, §2.1 y G-1 a G-6, B-1 (DESIGN §8.9, tokens de rol existentes), §4 (con el resumen y el `<figcaption>` de DESIGN §7.7), §10 |
| U3-12 · implementación y supuestos | §11, §14, §15, el bloqueo `docs/estado/bloqueos/U3-t2-acumulados.md` (RESUELTO, Q2-65 A) y los cinco `docs/estado/bloqueos/U3-tablero-*.md` (PROVISIONAL de U3 2.0) |

---

## 14. Supuestos

1. Textos de rangos cortos en tarjetas («Del 16 al 31 oct», «Lun 12 al Dom 18 oct») y de meses en gráficas («oct-2026»): derivados del formato de U1 §0.4; no cambian cifras, esquema, operaciones ni tokens.
2. Hora de carga sin conexión como «[HH:MM]» en 24 horas (la misma forma que `send_time`, K43); solo texto.
3. Las tarjetas del Resumen no llevan filtros (vista de la casa); los filtros viven en Detalle, Tendencias y Calendario (§6).
4. K-3 (antes la tarjeta 3), con varias líneas revolventes, muestra la primera de `availability[]` y «+[n] líneas» (las demás en D-3).
5. (Retirado en la versión 1.1: T-2 sigue DESIGN §8.3, con la capa 2 en cero en los tramos de pago adelantado, y ya no cambia a dos líneas.)
6. Prioridad entre estados de pago (Pagado, Vencido, Por vencer, Parcial, Pendiente) y la línea «Falta la parte de …» en un pago vencido o por vencer (§1.4): DESIGN §7.8 da los estados pero no su prioridad; solo orden de presentación y texto.
7. Textos «Falta la parte del hogar» y «Vence en [d] días» (§1.4), para los casos que DESIGN §7.8 no nombra; solo texto.
8. Línea de meta de T-4b con el color de la persona y trazo de puntos 2 · 4 sin marcadores, distinta de los trazos de Axel (continuo) y Ximena (8 · 4); DESIGN §8 no fija la línea de meta; solo dibujo.
9. Redacción de los resúmenes de una frase de T-1 a T-4 (§4): DESIGN §7.7 pide el resumen pero no su texto; U3 lo arma solo con cifras de C1, sin comparaciones; solo texto.
10. (2.0) La pestaña sigue llamándose «Resumen» (decisiones §G); K54 cambia su contenido, no su nombre.
11. (2.0) La comparación contra el mes anterior se dice en dinero con «más que», «menos que» o «Igual que» y los iconos `arrow-up-right` y `arrow-down-left` del mapa de DESIGN §5.4; sin porcentaje de variación (la interfaz no calcula, U1 §0.2).
12. (2.0) La sugerencia «subir o recortar presupuesto» va con rol información (`info`) y la palabra «Sugerencia»: es una lectura del motor, no una alerta.
13. (2.0) Orden de apilado de §2.4 (gasto, barómetro, pagos, ingreso, créditos, gráficas, «Más del hogar»): sigue el propósito de K54; solo orden.
14. (2.0) Nota fija de G-1 («Los pagos de crédito pagan compras que ya cuentan en el gasto de su mes») para que las tres columnas no se lean como sumables; solo texto.

Los supuestos 1 a 9 no cambian cifras, esquema, operaciones ni tokens. Los 10 a 14 tampoco; lo que sí cambia operaciones o claves (campos nuevos de §2.6) se publicó en C1 v1.5 y en el esquema 1.5 y lo que Axel no decidió es PROVISIONAL con bloqueo.

## 15. Dependencias y abiertos
- **U3 2.0 (K54)**, dueños y pedidos:
  - C1 v1.5 (`s1-c1-v15`), HECHO: las 3 operaciones nuevas, el par `getDashboardSettings` / `updateDashboardSettings`, los 14 campos nuevos de `DashboardMonthVsPlan`, `IncomeMonth` y `CreditTotals` (§2.6) y las métricas M-18 a M-22 en C1-datos §9. Backend de las 5 operaciones HECHO en `backend/internal/api/tablero/` (corrección r2, hallazgo A-1: motores S2 y S3, sin recalcular cifras; `make verify-all` da EXIT 0). Falta que los manejadores de área llenen los 14 campos opcionales nuevos de `DashboardMonthVsPlan`, `IncomeMonth` y `CreditTotals` (hoy se omiten y la interfaz cae a su regla de §2) (H-1 de `docs/estado/wf3/REVISION-especificacion-r1.md`). Con un 501 o un campo ausente, los bloques muestran lo que ya existe (regla de §2).
  - Esquema 1.5, HECHO (migración 00005 con `@param` y bloqueo): `debt.capacity_warning_pct` = 30, `debt.capacity_limit_pct` = 40 (FIRMES, K54), `dashboard.payment_reminder_count` = 5, `budget.signal_months` = 3 y `budget.signal_margin_pct` = 10 (PROVISIONALES).
  - U4: en qué sección de Ajustes se editan las 5 claves (propuesta: las dos de `debt.*` en P-26.05 «Créditos y condiciones», las dos de `budget.signal_*` en P-26.02 «Presupuesto» y la de `dashboard.*` en «Notificaciones»), con vista previa de su efecto.
  - U1 (fuera del alcance de `u3-v2`): §2.1 y §2.2 fila P-01 («6 tarjetas»), §3.3 (N-020 a N-029 más las necesidades nuevas del tablero), §6 fila P-01 y F-4 paso 1 («Toca la tarjeta «Próximo pago»» pasa a «Toca el pago en K-4»). F-3 ruta A ya se aclaró (C3-flujos punto 2, opción c).
  - P-04 → P-15: el renglón del ticket enlaza a P-15 con su producto (paso 4 de §2.2).
  - DESIGN 1.1 §8.9 «Barómetro» y la fila nueva de §8.2 ya existen (agente `u3-v2`).
  - Frontend P-01 nuevo, capturas U5 de §12 y C3 regenerada: WF-3c.
- DESIGN §7.6 a §7.8 y §8 ya existen y U3 los sigue desde la versión 1.1. Solicitudes a DESIGN (U2b), con su razón (detalle en `docs/estado/wf1/correccion-U3-r1.md`):
  1. §7.6: la tarjeta entera lleva a la pantalla destino de su acción primaria, no a la pestaña Detalle (U3-01: acción primaria única y la tarjeta 4 a P-10); y la anatomía suma el icono de 20 px junto a la etiqueta (§5.2; §5.4 ya prevé iconos en tarjetas).
  2. §7.8: la línea «Sin marcar como pagado» junto a «Vencido» (U1 N-022, C1 M-03), la prioridad entre estados y los textos «Falta la parte del hogar» y «Vence en [d] días» (U3 §1.4, supuestos 6 y 7).
  3. §8.3: quitar del subtítulo el rótulo «Deuda viva …» o justificarlo frente a K19, N-091 y M-11 («no se rotula deuda viva»).
  4. §8: fijar la línea de meta en una gráfica de líneas (U3 T-4b usa el color de la persona y trazo de puntos 2 · 4 sin marcadores, supuesto 8).
  5. §7.7: el ejemplo de resumen («El gasto subió 12 % respecto a septiembre») pide calcular una variación que la interfaz no hace (U1 §0.2). U3 usa resúmenes que solo enuncian cifras de C1 (§4). Pedido: cambiar el ejemplo o aclarar que la variación tiene que llegar calculada desde C1.
- C1-datos §9: dos métricas para D-4 y D-5, «Comidas de hoy» (N-121 vía N-027) y «Diferencias de precio» (N-162 vía N-027), con la fórmula, el periodo y los filtros que U3 §3 escribe mientras tanto (hallazgo H-7).
- C1 (`getAccrualVsCash`, esquema `AccrualVsCash.months[]` en `openapi.yaml`, y C1-datos M-11): agregar `cumulative_accrued_cents` y `cumulative_cash_cents`, con el mismo origen que `cumulative_difference_cents` (Q2-65 A, aprobada en la puerta 2; bloqueo `docs/estado/bloqueos/U3-t2-acumulados.md` RESUELTO; C1 v1.3).
- C1: `second_store_id` en `DashboardBestBasket` para que la tarjeta 5 nombre la segunda tienda (Q2-86 A: solicitud aprobada; C1 v1.3). Mientras no llegue, «menos que la siguiente tienda» y la tienda se ve en D-5; no bloquea.
- C1 (`SavingsGoal.status`, que sale de `backend/internal/api/c1gen/gen2_schemas.py`): C1 v1.3 todavía declara `on_track`, `behind` y `achieved` («en ruta, atrasada o cumplida»). Q2-18 A dejó tres estados, `in_progress`, `achieved` y `overdue` (`docs/contratos/motores/S2-interfaz.md` 1.3 §3.13), que son los que usan la tarjeta 6 y D-6.
- C1-datos §9, M-16 «Metas» (N-111, N-027): dice «% de avance entero»; con Q2-17 C la métrica de Metas va con un decimal (`progress_tenths_pct`, ya en C1 v1.3) y el entero queda en M-07 (tarjeta 6). El bloqueo `docs/estado/bloqueos/S2-metas.md` lo pide en «Qué cambia» («alinear RN-18 con M-07 y M-16»). D-6 ya usa `progress_tenths_pct`.
- C1 (`getCashFlowFortnights`, `rows[].planned_savings_cents`): C1 v1.3 lo declara `MoneyCents` obligatorio y describe `balance_cents` como «income − egresos − ahorro». Con el bloqueo abierto `docs/estado/bloqueos/S2-ahorro-quincenal.md`, el motor S2 entrega el ahorro en `null` y no lo resta del saldo (S2-interfaz 1.3 §3.11). Pedido: que C1 marque el campo como PROVISIONAL con ese bloqueo, admita `null` y no dé por hecha la pregunta 3 en la descripción de `balance_cents`. U3 D-2 muestra «—» mientras llegue `null`.
- Bloqueos de S2 que tocan cifras del dashboard (U3 solo muestra lo que llegue y marca PROVISIONAL lo que depende de un bloqueo abierto):
  - Abiertos: `docs/estado/bloqueos/S2-flujo-quincenal.md`, solo Q2-12 (gasto variable estimado de D-2); `docs/estado/bloqueos/S2-ahorro-quincenal.md` (valor del ahorro planificado de D-2: base, metas que cuentan y si resta del saldo; el reparto de Q2-13 B ya es FIRME); y `docs/estado/bloqueos/S2-metas-vencida.md` (aportación de una meta vencida sin cumplir en D-6).
  - RESUELTOS en la puerta 2, sin marca: `docs/estado/bloqueos/S2-presupuesto-total.md` (K-1, antes tarjeta 1; Q2-19 B y Q2-20 A), Q2-13 B y Q2-14 A de `S2-flujo-quincenal` (tarjeta 2; Q2-13 B fija solo el reparto del ahorro planificado, cuyo valor en D-2 sigue abierto, arriba), `docs/estado/bloqueos/S2-metas.md` (tarjeta 6 con el avance entero y D-6 con un decimal; Q2-16 A, Q2-17 C y Q2-18 A) y `docs/estado/bloqueos/S2-proyeccion-interruptor.md` (T-2, meses 7 a 12 con el interruptor encendido; Q2-21 a Q2-24 C).
  - Ninguno cambia la estructura, los textos fijos ni las operaciones de U3.
- Bloqueo propio de U3: `docs/estado/bloqueos/U3-t2-acumulados.md` (capas de T-2, sus etiquetas finales y su resumen): RESUELTO en la puerta 2 con Q2-65 A, una sola variante. El bloqueo `S1-parametros-avisos` (RESUELTO, Q2-11 B) no toca el dashboard: U3 no lee ni escribe `send_time`, `advance_days`, `weekday` ni `day_of_month`; los días de «Por vencer» son `alerts.payment_due_days_before`, ya decididos (decisión F, K44).
