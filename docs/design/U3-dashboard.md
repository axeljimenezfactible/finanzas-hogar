# U3 · Dashboard (P-01 Inicio) · Finanzas Hogar

Versión 1.2 (11-oct-2026) · Agente U3-dashboard (WF-1); ronda 1 de corrección: agente `correccion-U3-r1`; ronda 2: agente `correccion-U3-r2` · Evalúa: `docs/criterios/U3.md` v1.2.

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

**Fuentes que rigen** (no se reabren): decisiones §A, §B, §C, §G y §K (K5, K8, K38, K41, K43, K44); U1 1.2 §0.2 (vocabulario, semana por defecto K38), §0.4 (formatos), §0.5 (personas), §1.3 (móvil), §2.4 (estados), §3.3 (P-01: N-020 a N-029), §3.9 (P-07), §6.1 y §6.2 (puntos de quiebre); contrato C1 v1.2 por sus extractos `inicio.yaml`, `presupuesto-a.yaml`, `presupuesto-b.yaml` (y, para los paneles que reutilizan necesidades de otras pantallas, `comidas.yaml`, `compras-tiendas-a.yaml`, `creditos-a.yaml`, `creditos-b.yaml`, `ajustes-b.yaml`, `reportes.yaml`); métricas M-01 a M-17 de `C1-datos.md` §9; DESIGN.md §1, §2.4, §2.6 a §2.10, §3, §4, §5.2 a §5.4, §7.6 (tarjetas de métrica), §7.7 (gráficas), §7.8 (estados de pago) y §8 (patrones de visualización). Donde U3 se aparta de DESIGN lo dice en el lugar, con su razón, y lo pide a DESIGN en §15.

**Reglas que este documento no repite en cada bloque:**
- La interfaz **no recalcula reglas de negocio**: muestra, filtra, ordena y suma lo que está a la vista (U1 §0.2, «Motor Sx»). Toda cifra viene de la operación C1 citada.
- **Sin datos no es cero**: si falta una entrada se muestra «—» con el motivo, nunca un monto en cero ni `0 %` (U1 §0.2).
- **Toda cifra de presupuesto, flujo y ahorro es la parte del hogar**; la parte de Karen es informativa y solo aparece donde ella participa (U1 §0.2, §0.5; decisiones §A).
- Dinero con el formato de U1 §0.4 y DESIGN §3.3: MXN, símbolo `$`, coma de miles, dos decimales, cifras tabulares (`tnum`), negativo con el signo menos tipográfico, color de rol **e** icono (nunca solo color). Fechas en pantalla `10-oct-2026`; con día de la semana donde importa (`Vie 30 oct 2026`). En la API las fechas van en ISO; en este documento solo aparecen en formato de pantalla.
- Colores, tipografía, espacios, radios y sombras **solo por nombre de token** de C2 / DESIGN (0 valores sueltos).
- Ejemplos con cifras: solo las del oráculo o la semilla con su referencia (R-nn de `S2.md`/`S3.md`, K-nn de `S5.md`); lo demás va como marcador entre corchetes, sin número: `[plan]`, `[n]`. Los meses oct-2026 a mar-2027 de los ejemplos son los periodos P1 a P6 del oráculo: C1-datos M-01 dice «P1 = oct-2026», `is_canonical` de C1 marca «oct-2026 a mar-2027» y `docs/criterios/S2.md` S2-09 (nota de equivalencia) explica que con la semilla los dos dan las mismas cifras, porque ningún evento cae entre los días 1 y 9 de un mes. Por eso «octubre de 2026» cita `periodos[0]` en la tarjeta 1, D-7, T-2, A-06 y A-11. Una prueba con eventos entre los días 1 y 9 toma su esperado de las ventanas que declara, nunca de R-45 a R-50 (S2-09).

---

## 1. Estructura de Inicio

### 1.1 Pestañas y página
- **Título de página**: «Inicio» (`font.size.title-1`, `font.weight.bold`), icono `house`.
- **Pestañas** (orden fijo, decisiones §G): «Resumen» · «Detalle» · «Tendencias» · «Calendario». Abre en «Resumen». La pestaña elegida se recuerda en el dispositivo (frecuencia «Local», U1 §0.3). En < 768 px la fila de pestañas se desplaza en horizontal y la activa queda siempre visible (U1 §1.3); en ≥ 768 px es una fila completa.
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
- **Destino**: la tarjeta entera es un enlace a la pantalla destino de su acción primaria (P-06, P-07, P-18, P-10, P-14 y P-09, §2), no a la pestaña Detalle; Detalle se abre con su pestaña. U3 se aparta aquí de DESIGN §7.6 («la tarjeta entera es enlace al Detalle») porque U3-01 pide una acción primaria única por tarjeta y exige que la de la tarjeta 4 lleve a P-10; solicitud 1 a DESIGN (§15).
- Toda la tarjeta mide al menos `size.touch.min` de alto y es un solo objetivo táctil con foco visible (`border.width.focus`, anillo `color.brand.600` / `color.brand.400`). Nombre accesible = título + valor principal + estado.
- Cada tarjeta hace **su propia petición** a C1: si una falla, las otras se ven normales.
- Ninguna tarjeta usa los umbrales 80 % y 100 % como constantes: el estado ya viene calculado en la respuesta (`status` de M-01). Si una implementación necesitara el umbral, lo lee de `getAlertSettings` (C1-22, K44).
- Los días antes del pago (`alerts.payment_due_days_before`, inicial 1 y 0, decisión F y K44) sí se usan, para el estado «Por vencer» de la tarjeta 3, D-3 y el Calendario (§1.4). Inicio los lee **una sola vez** con `getAlertSettings` (GET `/api/v1/settings/alerts`, campo `parameters.alerts.payment_due_days_before.value`) y nunca como constante; si esa lectura falla, ningún pago se muestra «Por vencer» (queda «Pendiente») y nada más cambia.

### 1.3 Estados comunes de las tarjetas (U1 §2.4)
- **Carga**: esqueleto estático con la forma final (etiqueta, bloque de cifra y una o dos líneas) en `ui.skeleton.on-card` (DESIGN §7.6; movimiento de DESIGN §6.2); nunca cifras provisionales. Si tarda más de 10 s: «Esto está tardando más de lo normal · Reintentar».
- **Error**: mensaje de la tarjeta (texto exacto en cada una) + botón secundario «Reintentar» (borde `color.neutral.600` / `color.neutral.400`), icono `circle-alert` con `color.role.error.light` / `color.role.error.dark`. La tarjeta conserva su título y su enlace.
- **Sin conexión** (G-07, K9): franja global «Sin conexión con el servidor de casa» (icono `wifi-off`; «· [n] en cola» si hay capturas guardadas en el teléfono). Cada tarjeta muestra lo último cargado con la línea «Cargado a las [HH:MM]» (`font.size.caption`) y su enlace sigue abriendo la pantalla destino en solo lectura. Si la tarjeta nunca se cargó en el teléfono: «Sin conexión · Esta página aún no se ha cargado en este teléfono». «Registrar gasto» sigue disponible y deja el gasto en la cola (U1 §2.4).

### 1.4 Estados de un pago (DESIGN §7.8)
La tarjeta 3, D-3 y el Calendario usan los estados de pago de DESIGN §7.8, siempre con icono y texto (nunca solo color). Un pago muestra un solo estado: el primero de esta tabla que se cumpla.

| Estado | Cuándo (dato de C1) | Icono y color | Texto |
|---|---|---|---|
| Pagado | `status` = `paid` (todas las partes pagadas) | `circle-check`, `color.role.success.*` | «Pagado» |
| Vencido | `status` = `overdue` (en la tarjeta 3, `is_overdue` = true) | `calendar-x`, `color.role.warning.*` | «Vencido» y la línea «Sin marcar como pagado» (U1 N-022, C1 M-03) |
| Por vencer | pendiente, y los días que faltan para el pago son uno de los de `alerts.payment_due_days_before` (§1.2) | `calendar-clock`, `color.role.warning.*` | «Vence mañana» (1), «Vence hoy» (0); otro día configurado en Ajustes: «Vence en [d] días» |
| Parcial | pendiente, con una parte de `shares[]` pagada por completo (`paid_cents` = `amount_cents`) y otra no; el capital se repone solo cuando el pago se completa (K15) | `circle-dot`, `ui.icon.default` | «Falta la parte de Karen» (o «Falta la parte del hogar» si la que falta es la del hogar) |
| Pendiente | cualquier otro pago con alguna parte sin pagar | `circle-dashed`, `ui.icon.default` | «Pendiente» |

- Un pago «Vencido» o «Por vencer» que ya tiene una parte pagada agrega en su línea secundaria «Falta la parte de Karen» (o «del hogar»).
- Qué estados puede mostrar cada bloque: la tarjeta 3 trae `is_overdue` y `days_left` pero no lo pagado por parte, así que muestra «Vencido», «Por vencer» o ninguno, nunca «Parcial»; D-3 trae `shares[]` y solo pagos con alguna parte sin pagar, así que nunca muestra «Pagado»; el Calendario muestra los cinco (§5).
- «Sin marcar como pagado» acompaña a «Vencido» porque U1 N-022 y C1 M-03 piden esa marca y DESIGN §7.8 solo pone «Vencido». El orden de prioridad y los textos «Falta la parte del hogar» y «Vence en [d] días» tampoco están en DESIGN (supuestos 6 y 7; solicitud 2 a DESIGN, §15).

---

## 2. Resumen: 6 tarjetas, orden fijo

Orden fijo (U1 §3.3): 1 Mes contra plan → 2 Flujo de la quincena → 3 Próximo pago y disponible → 4 Comidas de la semana → 5 Mejor canasta → 6 Ahorro y metas. Ninguna tarjeta adicional (decisiones §G, U3-01). Los nombres son los de la lista del plan en español sin abreviaturas («Mes contra plan» = «Mes vs plan»; «Ahorro y metas» = «Ahorro»).

### Tarjeta 1 · Mes contra plan
- **Métrica C1**: M-01 (N-020) · `getDashboardMonthVsPlan` (GET `/api/v1/dashboard/month-vs-plan`) · esquema `DashboardMonthVsPlan`, campos `month`, `actual_cents`, `planned_cents`, `used_pct`, `remaining_cents`, `status` (`normal`, `warning`, `exceeded`), `days_left`.
- **Fórmula** (motor S2, M-01): real = Σ parte del hogar devengada en el mes; plan = Σ presupuesto vigente del mes de las categorías activas (RN-22); % usado = real ÷ plan × 100, entero ROUND_HALF_UP; restante = plan − real; estado según la respuesta; días que faltan contando hoy.
- **Periodo**: mes en curso (periodo base, decisiones §B). Petición: `period` = mes en curso (`2026-10`).
- **Contenido** (icono `wallet`):
  - Título: «Mes contra plan».
  - Valor principal: `[restante]` con la etiqueta «te quedan» (neutro). Si `remaining_cents` < 0: «−[monto] rebasado» con `color.role.warning.*` e icono `triangle-alert` (DESIGN §3.3, negativo sin rol propio).
  - Valor secundario (360): «[real] de [plan] · [n] %». Desde 412 suma: «Faltan [d] días del mes».
  - Insignias: `warning` → «Vas en [n] %» (rol advertencia, `triangle-alert`); `exceeded` → «Rebasaste el plan» (rol advertencia, `triangle-alert`); `normal` → sin insignia.
- **Acción primaria única**: enlace «Ver presupuesto» → P-06 Presupuesto del mes.
- **Vacío** (sin `planned_cents`): valor principal «—»; secundario «Aún no hay presupuesto para octubre» (mes de `month`) y «Llevas [real] gastado»; la acción cambia a «Capturar presupuesto» → P-26.02 (Ajustes › Presupuesto).
- **Error**: «No se pudo cargar el presupuesto del mes. Revisa que estés en la red de casa.» · «Reintentar».
- **PROVISIONAL** (`docs/estado/bloqueos/S2-presupuesto-total.md`): qué monto es el plan cuando existe un renglón de presupuesto total del mes, y qué % y estado lleva una categoría con plan de cero. La tarjeta muestra `planned_cents`, `used_pct` y `status` tal como lleguen; si `used_pct` no llega, muestra «—» sin insignia. El diseño no cambia con la respuesta.
- **Ejemplo trazable**: octubre de 2026, real (parte del hogar devengada) **$1,197.64** (`periodos[0].devengado.axel`, M-01; P1 = oct-2026); plan: marcador `[plan]` (lo captura Axel). Hoy 10-oct-2026: «Faltan 22 días del mes».

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
- **PROVISIONAL** (`docs/estado/bloqueos/S2-flujo-quincenal.md`): la regla para dar por cumplida una partida planificada (no contar dos veces un fijo ya pagado) afecta a «Por pagar»; la tarjeta muestra `committed_cents` tal como llegue. La parte del hogar del pago del 30 de octubre (R-13) no depende de esto.
- **Ejemplo trazable**: hoy 16-oct-2026, quincena «Del 16 al 31 oct»; «Por pagar» incluye la parte del hogar del pago del 30 de octubre, **$598.82** (R-13, M-02); ingresos: marcador `[ingresos]`.

### Tarjeta 3 · Próximo pago y disponible
- **Métrica C1**: M-03 y M-04 (N-022) · `getDashboardNextPayment` (GET `/api/v1/dashboard/next-payment`) · esquema `DashboardNextPayment`: `next_payment` (`due_date`, `credit_id`, `amount_cents`, `household_share_cents`, `external_share_cents`, `days_left`, `payment_count`, `is_overdue`) y `availability[]` (`credit_id`, `availability_reading`, `available_cents`, `availability_confirmed`). Nombre de la línea: `name` de `listCredits` (Credit).
- **Fórmula**: próximo pago = pago con alguna parte pendiente y fecha mínima igual o posterior a hoy entre todas las líneas; si hay pendientes con fecha anterior a hoy, primero el más antiguo; pagos del mismo día se suman (M-03). Disponible = disponible de hoy de cada línea revolvente con la **lectura activa** (motor S3, M-04).
- **Periodo**: de hoy en adelante (sin parámetros). El próximo pago sale del **calendario unificado** de pagos de todas las líneas (M-13).
- **Qué lectura muestra** (U3-07, decisiones §C): la tarjeta muestra **solo la lectura activa** del interruptor de Ajustes › Créditos y condiciones (por defecto **B**), con su etiqueta «Lectura B» (o «Lectura A»). **La otra lectura** se ve en Detalle › D-3 «Pagos y disponible», que pone las dos lecturas lado a lado con `previewAvailability` (M-04), y en P-18; la tarjeta no la muestra para no duplicar la cifra principal.
- **Contenido** (icono `calendar-clock`):
  - Título: «Próximo pago».
  - Valor principal: `[monto total]` (neutro).
  - Valor secundario (360): «[Día] [d] [mes] [año] · [línea]» con día de la semana (p. ej. «Vie 30 oct 2026 · Kueski»). Desde 412 suma: «Hogar [parte del hogar] · Karen [parte de Karen]» (Karen solo si `external_share_cents` > 0) y «Faltan [d] días» («Hoy», «Mañana»).
  - Si `payment_count` > 1: «[n] pagos este día».
  - Estado del pago (§1.4): si `is_overdue` = true, insignia «Vencido» (`calendar-x`, rol advertencia) y la línea «Sin marcar como pagado · venció hace [d] días»; si no, y `days_left` es uno de los días de `alerts.payment_due_days_before`, insignia «Vence mañana» o «Vence hoy» (`calendar-clock`, rol advertencia) en lugar de «Faltan [d] días». Sin «Parcial»: `next_payment` no trae lo pagado por parte (se ve en D-3).
  - Bloque «Disponible hoy» (una fila por línea de `availability[]`): «Disponible hoy [monto] · Lectura [A/B]»; si `availability_confirmed` = false, insignia **«Pendiente confirmar»** (rol información, icono `info`) con nombre accesible «Reposición de capital pendiente de confirmar con Kueski». Con más de una línea revolvente, la tarjeta muestra la del próximo pago y «+[n] líneas» (el resto en D-3).
- **Acción primaria única**: enlace «Ver crédito» → P-18 de la línea del próximo pago (sin próximo pago: la primera línea de `availability[]`).
- **Vacío**: sin `next_payment`: valor principal «Sin pagos pendientes» y se mantiene el bloque «Disponible hoy»; sin créditos dados de alta (`availability` vacío y sin pago): «Aún no hay créditos» y la acción cambia a «Agregar crédito» → P-26.05.
- **Error**: «No se pudieron cargar los pagos. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable** (U3-07): hoy 10-oct-2026 con lectura B: valor principal **$973.75**, «Vie 30 oct 2026 · Kueski», «Hogar $598.82 · Karen $374.93» (R-13), «Faltan 20 días»; «Disponible hoy **$6,905.50** · Lectura B · Pendiente confirmar» (R-64). Con el interruptor en A: «Disponible hoy **$10,000.00** · Lectura A» (R-64).

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
- **Métrica C1**: M-06 (N-024) · `getDashboardBestBasket` (GET `/api/v1/dashboard/best-basket`) · esquema `DashboardBestBasket`: `basket` (`shopping_list` o `expense`), `expense_id`, `store_id`, `total_cents`, `difference_vs_second_cents`, `covered_count`, `item_count`, `stale_price_count`, `assumes_membership`. Nombre de tienda: `listStores`.
- **Fórmula** (motor S5): lista de compras vigente (o la última compra con renglones si está vacía) comparada por tienda; ganadora = menor total; diferencia contra la segunda; cobertura = artículos con precio ÷ artículos; precios viejos usados; Sam's Club como si tuvieran membresía.
- **Periodo**: la **lista de compras vigente** (U1 §0.2): la de la semana por defecto según K38 (los domingos, la siguiente). El parámetro `week` es obligatorio (R3-01) y la PWA manda esa semana (Dom 11 oct 2026 → `2026-W42`; Sáb 10 oct 2026 → `2026-W41`).
- **Contenido** (icono `shopping-cart`):
  - Título: «Mejor canasta».
  - Valor principal: «[tienda] [total]» con punto del color de la tienda (`color.store.<tienda>.*`) y su nombre en texto (DESIGN §2.6 regla 1).
  - Valor secundario (360): «[diferencia] menos que la siguiente tienda». Desde 412 suma: «[covered] de [items] artículos con precio» y, si `stale_price_count` > 0, insignia «[n] precios viejos» (rol advertencia, `triangle-alert`).
  - Si `assumes_membership` = true: insignia «Como si tuvieran membresía» (rol información, `info`).
  - Si `basket` = `expense`: línea «Tu lista está vacía; se compara la compra del [fecha]».
- **Acción primaria única**: enlace «Comparar canasta» → P-14 (abre la semana según K38).
- **Vacío** (sin `store_id`: ningún artículo con precio, o sin lista ni compras con renglones): «—» y «Aún no hay precios para comparar»; la acción cambia a «Registrar precios» → P-15.
- **Error**: «No se pudo comparar la canasta. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable**: compra semilla del 10-oct-2026 (`basket` = `expense`): «Bodega Aurrera **$3,094.50**» (R-01), «**$138.50** menos que la siguiente tienda» (K-03 de `S5.md`: Walmart **$3,233.00**, `data.json` `totales.walmart`). C1 no entrega el nombre de la segunda tienda en esta operación; se ve en D-5 (ver solicitud del reporte).

### Tarjeta 6 · Ahorro y metas
- **Métrica C1**: M-07 (N-025) · `getDashboardSavings` (GET `/api/v1/dashboard/savings`) · esquema `DashboardSavings`: `saved_total_cents`, `contributed_month_cents`, `planned_month_cents`, `next_goal` (`SavingsGoal`: `name`, `target_cents`, `target_date`, `saved_cents`, `progress_pct`, `status` = `on_track`, `behind`, `achieved`).
- **Fórmula** (motor S2): ahorrado = Σ aportaciones − Σ retiros (fondo y metas); aportado del mes contra planificado del mes (Σ aportación mensual planificada activa); meta más próxima por fecha con % = ahorrado ÷ monto meta × 100, entero.
- **Periodo**: mes en curso (sin parámetros).
- **Contenido** (icono `piggy-bank`):
  - Título: «Ahorro y metas».
  - Valor principal: `[ahorrado]` con `color.role.savings.*` y la etiqueta «ahorrado» (rol con icono, DESIGN §2.4).
  - Valor secundario (360): «Este mes [aportado] de [planificado]». Desde 412 suma la meta más próxima: «[nombre] · [n] % · [fecha meta]» e insignia de estado: «En ruta» (rol información, `info`), «Atrasada» (rol advertencia, `triangle-alert`), «Cumplida» (rol éxito, `circle-check`).
- **Acción primaria única**: enlace «Ver ahorro» → P-09 Ahorro y metas.
- **Vacío** (sin fondo ni metas): «—» y «Aún no tienes fondo ni metas de ahorro»; la acción cambia a «Crear meta» → P-26.13 (Ajustes › Salud y metas de ahorro). Sin `next_goal` pero con fondo: se omite la línea de la meta.
- **Error**: «No se pudo cargar el ahorro. Revisa que estés en la red de casa.» · «Reintentar».
- **PROVISIONAL** (`docs/estado/bloqueos/S2-metas.md`): estado de la meta y forma del % (entero en C1, un decimal en RN-18). La tarjeta usa `progress_pct` (entero) y los tres valores de `status` de C1 (`on_track`, `behind`, `achieved`); si el bloqueo cambia los valores, cambian solo las insignias.
- **Ejemplo**: marcadores; el oráculo no trae ahorro (M-07).

---

## 3. Detalle

**Propósito**: el segundo nivel de cada tarjeta, sin salir de Inicio (U1 N-027), más el bloque «Participación de Karen» (U1 N-026). Responde «¿de dónde sale la cifra de la tarjeta?» con lo mínimo; lo completo vive en la pantalla destino.
**Contenido**: 7 paneles en el orden de las tarjetas (D-1 a D-6) y al final D-7. Cada panel es un bloque con título (`font.size.title-3`), contenido y un enlace «Ver todo» a su pantalla destino (enlace, no botón con relleno). Componentes: «Tabla compacta» y «Lista de propiedades» (DESIGN §7).
**Filtros**: barra de vista de Detalle (en < 768 px, botón «Filtrar» con icono `list-filter` que abre una hoja inferior y chips con los filtros activos bajo el título, U1 §1.3). Cada panel aplica solo los filtros que su operación admite (§6); un filtro que un panel no admite **no** lo cambia y el panel lo dice con la línea «Este bloque no usa el filtro [nombre]» (`font.size.caption`).
**Periodo por defecto**: mes en curso; la quincena es subperiodo (selector «Mes · 1.ª quincena · 2.ª quincena», valores `period` = `2026-10`, `2026-10-q1`, `2026-10-q2`; K8).
**Estados**: los de §1.3 por panel (cada panel hace su petición), con los textos de §7.

| Panel | Contenido | Métrica · operación C1 (campos) | Vacío |
|---|---|---|---|
| D-1 Presupuesto por categoría | Las 5 categorías con mayor % usado: nombre, barra horizontal real contra plan (patrón «barras para comparar categorías»), «[real] de [plan] · [n] %» e insignia de estado como en la tarjeta 1 | M-09 (N-081) · `getBudgetByCategory` (`month`, `rows[]`: `category_id`, `planned_cents`, `actual_cents`, `difference_cents`, `used_pct`, `status`, `expense_count`); nombres con `listCategories` | «Aún no hay presupuesto por categoría» · «Capturar presupuesto» → P-26.02 |
| D-2 Quincena en curso | Fila de la quincena en curso como lista de propiedades: ingresos (marca «esperado» si `income_is_expected`), pagos de crédito (parte del hogar), fijos y venideros, gasto variable estimado, ahorro planificado, saldo de la quincena y saldo acumulado | M-10 (N-090) · `getCashFlowFortnights` (`rows[]`: `start_date`, `end_date`, `pay_date`, `income_cents`, `income_is_expected`, `credit_payments_cents`, `fixed_and_upcoming_cents`, `variable_estimate_cents`, `planned_savings_cents`, `balance_cents`, `cumulative_balance_cents`, `is_actual`). Gasto variable estimado y ahorro planificado: **PROVISIONAL** (`docs/estado/bloqueos/S2-flujo-quincenal.md`) | «Falta tu ingreso quincenal» · «Capturar ingreso» → P-26.03 |
| D-3 Pagos y disponible | (a) Pagos de todas las líneas en dos grupos, cada uno en orden de fecha: primero **los vencidos** (del más antiguo al más reciente) y luego **los próximos 5 pendientes** de hoy a hoy + 60. Cada fila: fecha con día de la semana («Vie 30 oct 2026»), línea, monto total, «Hogar [parte] · Karen [parte]» y estado de §1.4 («Vencido» con «Sin marcar como pagado», «Por vencer», «Parcial» con «Falta la parte de Karen», o «Pendiente»; «Pagado» no llega a este panel). Si hay más de 5 vencidos, la línea «Hay más pagos vencidos» y el enlace «Ver todo». (b) Disponible de cada línea revolvente con **las dos lecturas lado a lado** («Lectura A [monto]» · «Lectura B [monto]», la activa con la insignia «Activa») y el mínimo del horizonte de cada lectura; aviso «Reposición de capital pendiente de confirmar con Kueski» mientras no se confirme | M-13 (N-192) · `listScheduledPayments` (`id`, `due_date`, `amount_cents`, `credit_id`, `shares[]`: `person_id`, `amount_cents`, `paid_cents`, `status`; `status`) en **dos peticiones**, las dos con orden `due_date`, `limit` = 5 y los filtros de persona, familia y crédito de la barra (§6): (1) **vencidos**: `status` = `overdue` y `period` = `[fecha del vencido más antiguo]..[ayer]` (rango inclusivo de C1), donde esa fecha es `next_payment.due_date` de `getDashboardNextPayment` cuando `is_overdue` = true (M-03 pone primero el vencido más antiguo; se reutiliza la respuesta de la tarjeta 3 o se pide); con `is_overdue` = false no hay vencidos y esta petición no se hace; (2) **próximos**: `status` = `pending` y `period` = `[hoy]..[hoy + 60]` (M-13 «hoy..+60»). Días de «Por vencer»: `getAlertSettings` (§1.2) · M-04 (N-285, N-201) · `previewAvailability` por línea revolvente | «Sin pagos pendientes» (b se sigue mostrando); sin créditos: «Aún no hay créditos» · «Agregar crédito» → P-26.05 |
| D-4 Comidas de hoy | Casillas de hoy por tiempo (desayuno, comida, cena, colación) y persona (Axel, Ximena): receta o «Fuera», kcal y costo | N-121 · `getMealPlanWeek` (`cells[]`: `plan_date`, `meal_slot_id`, `person_id`, `recipe_id`, `is_eating_out`, `cost_cents`, `kcal`), con `period` = semana en curso; la PWA muestra las casillas cuyo `plan_date` es hoy. **Métrica pendiente en C1-datos §9** (solicitud a C1, §15); mientras tanto rige esta definición: «Comidas de hoy» (N-121 vía N-027) = casillas de `cells[]` con `plan_date` = hoy, por tiempo y persona, con `recipe_id` o `is_eating_out`, `kcal` y `cost_cents` tal como llegan (sin sumas); periodo: la semana ISO de lunes a domingo que contiene hoy (K8; también en domingo, porque el panel muestra hoy y K38 solo cambia la semana que abre P-10); filtro: `person_id` (S) | «Hoy no hay comidas en el plan» · «Armar plan» → P-10 (regla K38) |
| D-5 Diferencias de precio | Los 5 artículos con mayor diferencia de precio entre tiendas de la canasta de la tarjeta 5: artículo, tienda ganadora (punto + nombre), diferencia, marca «viejo» si el precio es viejo; y el total por tienda (orden fijo de tiendas, DESIGN §2.6) | N-162 · `getBasketComparison` (`stores[]`: `store_id`, `total_cents`, `difference_vs_cheapest_cents`, `stale_count`, `assumes_membership`; renglones con `product_id`, `winner_store_id`, `difference_cents`, `prices[]`.`is_stale`), con el mismo `week` (y `basket`, `expense_id`) de la tarjeta 5. **Métrica pendiente en C1-datos §9** (solicitud a C1, §15); mientras tanto rige esta definición: «Diferencias de precio» (N-162 vía N-027) = los 5 renglones con mayor `difference_cents`, en orden descendente (la interfaz ordena lo que está a la vista, U1 §0.2), y `stores[]` con `total_cents` tal como llega; periodo: la lista vigente (`week` según K38, como la tarjeta 5); filtro: tienda en la vista (V) | «Aún no hay precios para comparar» · «Registrar precios» → P-15 |
| D-6 Metas | Cada meta activa: nombre, barra de avance (`color.role.savings.*` con el % en texto), «[ahorrado] de [meta]», aportación mensual necesaria, fecha meta y estado (insignias de la tarjeta 6) | M-16 (N-111) · `listSavingsGoals` (`name`, `target_cents`, `target_date`, `saved_cents`, `progress_pct`, `required_monthly_cents`, `status`), con `status` activo. Aportación mensual necesaria y estado: **PROVISIONAL** (`docs/estado/bloqueos/S2-metas.md`) | «Aún no tienes metas de ahorro» · «Crear meta» → P-26.13 |
| D-7 Participación de Karen | Ver §3.1 | M-08 (N-026) · `getPersonParticipation` | Ver §3.1 |

Ejemplos trazables de Detalle: D-3 (a) con hoy 10-oct-2026, rango de hoy a hoy + 60 (hasta el 9-dic-2026, inclusivo): «Vie 30 oct 2026», «Dom 15 nov 2026» y «Lun 30 nov 2026», cada uno por **$973.75** con «Hogar $598.82 · Karen $374.93» y estado «Pendiente» (R-13 a R-15, M-13); no hay vencidos (la tarjeta 3 trae `is_overdue` = false) y el pago del 15-dic-2026 (R-16) cae fuera del rango. D-3 (b) el 10-oct-2026: «Lectura A **$10,000.00**» · «Lectura B **$6,905.50**» (R-64); mínimo del horizonte **$8,094.50** y **$5,000.00** (R-61, R-62). D-5 con la compra semilla: Bodega Aurrera **$3,094.50** (R-01) y Walmart **$3,233.00** (`totales.walmart`), diferencia **$138.50** (K-03).

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
- **Métrica C1**: M-11 (N-091) · `getAccrualVsCash` (GET `/api/v1/cash-flow/accrual-vs-cash`) · `projection_mode`, `months[]`: `month`, `accrued_cents`, `cash_cents`, `accrued_household_cents`, `accrued_external_cents`, `cash_household_cents`, `cash_external_cents`, `cumulative_difference_cents`, `is_canonical`.
- **Fórmula y periodo**: mes en curso a 11 meses después; diferencia acumulada = devengado acumulado − efectivo acumulado (motor S2). **No se rotula «deuda viva»** (K19, N-091).
- **Patrón**: **área apilada devengado contra efectivo**, acumulada. Orden fijo de capas, de abajo arriba: 1 «Efectivo acumulado» = `color.chart.1`, relleno sólido, borde superior continuo con marcador círculo; 2 «Diferencia acumulada» (la capa «devengado sin pagar» de DESIGN §8.3; el rótulo sigue a C1 M-11) = `color.chart.2`, relleno con trama diagonal, borde superior discontinuo con marcador cuadrado. La orilla superior de la pila es el devengado acumulado y lleva la etiqueta «Devengado acumulado». **Las dos capas tienen el mismo origen** que `cumulative_difference_cents`, que C1 define como `periodos[i].diferencia_acumulada` de `model.json`: se acumula desde P1 aunque `months[]` empiece en el mes en curso. Capa 1 = efectivo acumulado desde ese origen; cima = devengado acumulado desde ese origen; capa 2 = cima − capa 1 = `cumulative_difference_cents` tal como llega. **Prohibido** sumar `cash_cents` desde el primer mes visible y apilar encima `cumulative_difference_cents`: mezcla dos orígenes y, desde el segundo mes, la cima deja de ser el devengado acumulado.
- **Fuente de las capas · PROVISIONAL** (`docs/estado/bloqueos/U3-t2-acumulados.md`): C1 v1.2 solo trae `cumulative_difference_cents` con ese origen, y la interfaz no puede reconstruir los otros dos acumulados con lo que está a la vista.
  - Variante (a), la que U3 pide (solicitud a C1 en §15): C1 agrega en `months[]` `cumulative_accrued_cents` y `cumulative_cash_cents`, con el mismo origen que `cumulative_difference_cents` (así `cumulative_difference_cents` = `cumulative_accrued_cents` − `cumulative_cash_cents`). Capa 1 = `cumulative_cash_cents`; cima = `cumulative_accrued_cents`; capa 2 = `cumulative_difference_cents`. Etiquetas «Devengado acumulado» y «Efectivo acumulado».
  - Variante (b), sin cambiar C1: todo relativo a la ventana. Capa 1 = Σ `cash_cents` de los meses visibles hasta ese mes; cima = Σ `accrued_cents` de los meses visibles hasta ese mes; capa 2 = cima − capa 1. Leyenda y etiquetas: «Efectivo en estos 12 meses», «Diferencia en estos 12 meses» y «Devengado en estos 12 meses». La tabla muestra además `cumulative_difference_cents` con su propio rótulo, «Diferencia acumulada desde el inicio». Un tramo con capa 2 negativa vale cero y lleva «En estos meses se pagó más de lo devengado», porque puede ser deuda de antes de la ventana y no un pago adelantado.
  - No dependen del bloqueo: la tabla mensual (devengado, efectivo, partes del hogar y de Karen, y `cumulative_difference_cents`), la fila «Total», la línea «Hoy», los meses 7 a 12 y los estados de §7.
- **Pago adelantado** (DESIGN §8.3): si un mes trae `cumulative_difference_cents` negativo (se pagó antes de devengar), la capa 2 **vale cero en ese tramo** y el texto «Pagado por adelantado» lo dice junto al tramo; nunca se dibuja una capa negativa y la gráfica sigue siendo un área apilada (en ese tramo la cima de la pila es el efectivo acumulado). La tabla muestra la diferencia tal como llega, con su signo menos. Esta regla es la de la variante (a); en la variante (b) rige el texto de su bloque.
- Línea vertical «Hoy» de 1 px en `ui.viz.baseline` con su etiqueta (DESIGN §8.3; única retícula vertical, §8.1).
- **Resumen** (DESIGN §7.7), con las cifras de las etiquetas finales: «Hasta [último mes], con proyección desde [séptimo mes]: devengado acumulado [monto] y efectivo acumulado [monto]». No nombra la diferencia como deuda ni como «debes» (K19, N-091, M-11). PROVISIONAL con las capas (`docs/estado/bloqueos/U3-t2-acumulados.md`); en la variante (b): «En estos 12 meses, hasta [último mes]: devengado [monto] y efectivo [monto]».
- Etiquetas en texto al final de la pila (DESIGN §8.3 y §8.5): «Devengado acumulado» y «Efectivo acumulado», con su valor final (en la variante (a), `cumulative_accrued_cents` y `cumulative_cash_cents` del último mes; PROVISIONAL con las capas). El subtítulo da solo el periodo y **no** rotula la deuda viva (K19, N-091, M-11): U3 se aparta aquí del ejemplo de subtítulo de DESIGN §8.3; solicitud 3 a DESIGN (§15).
- Meses 7 a 12: trazo discontinuo y la nota de lo que incluye la proyección según `projection_mode` (K21, N-093). Con el interruptor encendido, esos meses son **PROVISIONAL** (`docs/estado/bloqueos/S2-proyeccion-interruptor.md`) y la nota lo dice: «Proyección con despensas repetidas: por confirmar»; con el interruptor apagado (por omisión) no dependen del bloqueo. La tabla trae también las partes del hogar y de Karen de cada mes (`accrued_household_cents`, `accrued_external_cents`, `cash_household_cents`, `cash_external_cents`); la de Karen es informativa, de la despensa en común (§9).
- **Vacío**: «Aún no hay gastos ni pagos para comparar» · «Registrar gasto» (acción de página). **Ejemplo trazable** (oct-2026 = P1 y mar-2027 = P6, regla de ejemplos al inicio; C1-datos M-01, `is_canonical`): fila oct-2026, devengado **$1,947.50**, efectivo **$973.75** y diferencia acumulada **$973.75** (R-45, M-11); fila nov-2026, **$1,947.50**, **$1,947.50** y **$973.75** (R-46). Capas en la variante (a), que no cambian con el mes en curso porque el origen no se mueve: en oct-2026, efectivo acumulado **$973.75** y devengado acumulado **$1,947.50** (R-45); en nov-2026, efectivo acumulado **$2,921.25** (= 973.75 + 1,947.50, efectivo de R-45 y R-46) y devengado acumulado **$3,895.00** (= 1,947.50 + 1,947.50, devengado de R-45 y R-46); en mar-2027, **$13,103.78** (R-52) y **$14,362.48** (R-51), con diferencia acumulada **$1,258.70** (R-50). **Fila «Total» de la tabla**: suma de las 12 filas visibles (la interfaz suma lo que está a la vista, U1 §0.2), sin cifra de canon: R-51 y R-52 son los totales de P1 a P6 (seis filas), no de doce.

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
| T1 Mes contra plan | — | — | — | — | S · mes en curso | — |
| T2 Flujo de la quincena | — | — | — | — | S · quincena en curso (K8) | — |
| T3 Próximo pago y disponible | — | — | — | — | — · de hoy en adelante | — |
| T4 Comidas de la semana | — | — | — | — | S · semana en curso (K8) | — |
| T5 Mejor canasta | — | — | — | — | — · lista vigente (`week` según K38) | — |
| T6 Ahorro y metas | — | — | — | — | — · mes en curso | — |
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

Las tarjetas del Resumen no llevan filtros (son la vista de la casa); el filtro vive en Detalle, Tendencias y Calendario. El enlace «Ver todos los gastos con estos filtros» de Detalle lleva los seis filtros a P-02.

---

## 7. Estados por tarjeta y panel

Reglas comunes en §1.3 (carga con esqueleto y aviso a los 10 s; error con «Reintentar»; sin conexión con lo último cargado y su hora). Celdas = texto exacto o regla.

| Bloque | Vacío (qué se ve · acción) | Carga | Error (texto · acción) | Sin conexión |
|---|---|---|---|---|
| T1 | «Aún no hay presupuesto para [mes]» · «Capturar presupuesto» | esqueleto de tarjeta | «No se pudo cargar el presupuesto del mes. Revisa que estés en la red de casa.» · «Reintentar» | último dato + «Cargado a las [HH:MM]» |
| T2 | «Falta tu ingreso quincenal» · «Capturar ingreso» | esqueleto de tarjeta | «No se pudo cargar el flujo de la quincena. …» · «Reintentar» | ídem |
| T3 | «Sin pagos pendientes» / «Aún no hay créditos» · «Agregar crédito» | esqueleto de tarjeta | «No se pudieron cargar los pagos. …» · «Reintentar» | ídem |
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
| Tarjetas del Resumen | 1 columna, orden fijo 1 a 6, relleno de tarjeta `space.4`. A 360: título, valor principal y 1 línea secundaria; a 412 suma 1 dato secundario (U1 §6.2). La cifra principal nunca se trunca: si no cabe, se quita primero lo secundario (prueba de ancho de U1 §6.1 regla 5) | 768: 2 columnas × 3 filas (8 columnas, margen `space.6`, relleno de tarjeta `space.6`, DESIGN §7.6); 1024: 3 columnas × 2 filas (12 columnas, margen `space.8`, barra lateral fija de `layout.sidebar.width`) | 3 columnas × 2 filas con medianil `space.6` y relleno de tarjeta `space.6`; todas las líneas secundarias visibles |
| Acción primaria de página | «+» de la barra inferior (`size.bottombar.item`) | 768: «+ Registrar gasto» en la barra superior (`layout.topbar.height`); igual en 1024 | igual que 1024 |
| Pestañas | fila con desplazamiento horizontal, activa visible | fila completa | fila completa |
| Detalle | paneles a 1 columna; tablas compactas de 2 líneas (título y cifra; 1 propiedad debajo; a 412, 1 más); filtros en hoja inferior con chips | 768: 1 columna con tablas de hasta 5 columnas; 1024: 2 columnas de paneles | 2 columnas de paneles; D-7 a todo el ancho al final |
| Tendencias | 1 gráfica por fila, alto reducido, leyenda debajo, máximo 6 series, valor al tocar; etiqueta del periodo actual en texto | 1 gráfica por fila con leyenda a la derecha y valor al pasar el puntero | 2 gráficas por fila (T-1 con T-2; T-3 con T-4) |
| Calendario | agenda de hoy a hoy + 30 (412: mes compacto opcional) | mes | mes con panel del día seleccionado |
| Tipografía | `font.size.title-2` cifra; `font.size.small` en `ui.text.secondary` para la etiqueta de tarjeta y el título de gráfica (DESIGN §7.6, §7.7); `font.size.title-3` título de panel de Detalle; `font.size.small` secundario; `font.size.caption` metadatos | igual | igual |
| Alto de controles | `size.control.coarse` (puntero grueso) | `size.control.coarse` si es táctil; `size.control.fine` con puntero fino | `size.control.fine` (puntero fino) |

Teléfono en horizontal: rige el ancho resultante (915 px usa la regla de 768, U1 §6.1 regla 7).

---

## 9. Participación de Karen (reglas)

- Karen **no tiene cuenta y no ve el dashboard** (decisiones §A). Solo Axel y Ximena, co-titulares, ven Inicio, y ven lo mismo (U1 §0.5).
- Aparece **solo en lo que comparte**: (a) la despensa en común y (b) sus compras con el crédito Kueski y sus pagos. En Inicio eso es: la parte de Karen en la tarjeta 3 (próximo pago), en D-3 y en los pagos del Calendario (incluido el estado «Parcial · Falta la parte de Karen», §1.4); el panel D-7; las columnas de su parte en la tabla de T-2 («Ver como tabla»: devengado y efectivo de su parte de la despensa en común, `accrued_external_cents` y `cash_external_cents`, informativas); el valor «Karen» del filtro de familia (K5), que incluye lo conjunto.
- Su parte es **informativa** para el flujo de Axel: no suma a las tarjetas 1, 2 y 6, que usan la parte del hogar (U1 §0.2).
- **No hay indicadores de las finanzas de su familia**: ni ingresos, ni presupuesto, ni ahorro, ni gráficas propias; las gráficas por persona (T-4) solo tienen a Axel y Ximena (plan de comidas y salud: decisiones §A).
- Su estado de cuenta PDF es una acción del hogar desde P-25 (`generateExternalStatement`) y **Axel lo reenvía a mano** (decisiones §A, F).

---

## 10. Color, tipografía y accesibilidad

- Solo tokens de C2 / DESIGN por nombre: fondo y tarjeta `color.neutral.50` / `color.neutral.100` (oscuro `color.neutral.950` / `color.neutral.900`); texto `color.neutral.900` / `color.neutral.700` (oscuro `color.neutral.100` / `color.neutral.300`); enlaces `color.brand.700` / `color.brand.300`; roles `color.role.{income,expense,savings,debt,warning,error,success,info}.{light,dark,surface-light,surface-dark}`; personas `color.person.{axel,ximena,karen,shared}.*`; tiendas `color.store.{walmart,bodega-aurrera,sams-club,market}.*`; series `color.chart.1` a `color.chart.8` y «Otros» `color.neutral.600` / `color.neutral.400`. Ningún color nuevo (DESIGN §2.10 regla 4).
- El color nunca es el único portador de significado: cifras con signo e icono (DESIGN §3.3), estados con icono y texto, series con forma, patrón y etiqueta, personas y tiendas con nombre (K41). Ingreso y éxito no comparten componente: ingreso solo en cifras con «+», éxito solo en insignias con palomita (DESIGN §2.4).
- Cada gráfica tiene «Ver como tabla» con las mismas cifras, orden fijo de series y etiquetas de dinero en texto (§4).
- Contraste: los pares de U3 son los de DESIGN §2.8 y §2.9 (202 de 202 cumplen); U3 no crea pares nuevos.
- Movimiento: las cifras no se animan (DESIGN §6.4); con movimiento reducido las gráficas aparecen sin transición.
- Teclado (escritorio): las pestañas se recorren con flechas; cada tarjeta es un enlace enfocable en el orden 1 a 6.

---

## 11. Paso a implementación (F2, WF-3)

### 11.1 Componentes de DESIGN que usa U3
«Tarjeta de métrica» (§1.2), «Pestañas», «Insignia de rol», «Enlace», «Botón primario» (solo «Registrar gasto»), «Botón secundario» («Reintentar»), «Esqueleto», «Franja sin conexión», «Barra de vista con filtros» y «Hoja inferior de filtros», «Chip de filtro», «Tabla compacta», «Lista de propiedades», «Barra de avance», «Gráfica» (barras agrupadas, barras horizontales, área apilada, líneas) con «Ver como tabla», «Calendario» (mes y agenda) y «Vista previa de registro». Los nombres siguen DESIGN §7 y §8; donde U3 se aparta de su contenido lo dice en el lugar (destino de la tarjeta en §1.2, «Sin marcar como pagado» y prioridad de estados en §1.4, subtítulo de T-2, línea de meta de T-4b, resumen de las gráficas sin variaciones calculadas en §4) y lo pide a DESIGN en §15.

### 11.2 Criterios de aceptación para F2
Pruebas con Playwright 1.64 + axe-core, reloj simulado en `America/Mexico_City` y datos de la semilla y el oráculo.

| Id | Criterio | Cómo se prueba |
|---|---|---|
| A-01 | El Resumen tiene exactamente 6 tarjetas en el orden de §2 | contar tarjetas a 360, 412, 768, 1024 y 1440 |
| A-02 | Cada tarjeta llama solo a su operación de §2 con los parámetros indicados; Inicio lee `getAlertSettings` una sola vez (§1.2); D-3 hace sus dos peticiones de `listScheduledPayments` (§3) y el Calendario la suya junto a `listCalendarEvents` (§5) | registro de peticiones de red por tarjeta y panel |
| A-03 | Tarjeta 4 y K38: con reloj en Dom 11 oct 2026 la tarjeta pide `period` = `2026-W41` y «Ver plan» abre P-10 en `2026-W42`; con Mié 14 oct 2026, ambas en `2026-W42` | dos corridas con fecha simulada |
| A-04 | Tarjeta 5 pide `week` de la lista vigente (Dom 11 oct 2026 → `2026-W42`) | petición de red |
| A-05 | Tarjeta 3 con la semilla el 10-oct-2026: $973.75, «Vie 30 oct 2026», «Hogar $598.82 · Karen $374.93», «Disponible hoy $6,905.50 · Lectura B», insignia «Pendiente confirmar»; con lectura A, $10,000.00 | captura y texto |
| A-06 | Tarjeta 1 en octubre de 2026 (= P1, regla de ejemplos al inicio) muestra real $1,197.64; sin presupuesto, «—» y «Capturar presupuesto» | texto |
| A-07 | Cada bloque muestra vacío, carga, error y sin conexión con los textos de §7; el aviso de 10 s aparece en la carga lenta | respuestas simuladas (vacía, lenta, 500, sin red) |
| A-08 | Una sola acción con relleno de marca en la página («Registrar gasto»); cada tarjeta, un solo enlace | conteo de botones primarios y enlaces por tarjeta |
| A-09 | Objetivos táctiles ≥ 44 × 44 px y sin desplazamiento horizontal de página en los 5 anchos | medición de cajas y del ancho de documento |
| A-10 | La cifra principal más larga cabe completa en la tarjeta a 360 px | prueba de ancho de U1 §6.1 regla 5 |
| A-11 | Gráficas con orden fijo de series, marcador o patrón distinto por serie, etiquetas de dinero en texto y «Ver como tabla» con las mismas cifras; T-2 con reloj en oct-2026: fila oct-2026 con $1,947.50, $973.75 y $973.75 (R-45); T-2 con reloj en nov-2026 (la ventana empieza en nov-2026): fila nov-2026 con $1,947.50, $1,947.50 y $973.75 (R-46); capas de la variante (a) (PROVISIONAL, `docs/estado/bloqueos/U3-t2-acumulados.md`) en nov-2026, efectivo acumulado $2,921.25 y devengado acumulado $3,895.00, y en mar-2027, $13,103.78 y $14,362.48 (R-52, R-51), las mismas con los dos relojes; en la variante (b), con reloj en nov-2026 las capas de mar-2027 valen R-52 menos el efectivo de R-45 y R-51 menos el devengado de R-45; la fila «Total» de T-2 suma las 12 filas visibles; con un mes de diferencia negativa la capa 2 vale cero con «Pagado por adelantado»; ninguna marca con contorno y ninguna gráfica con un número en cada punto; cada gráfica en `<figure>` con `<figcaption>` (título y resumen de §4), resumen sin porcentajes ni variaciones calculadas y, en T-2, sin la palabra «deuda» | inspección del DOM y de la tabla; T-2 con dos corridas de reloj simulado (oct-2026 y nov-2026) |
| A-12 | Calendario: franjas de quincena 1 al 15 y 16 al último día; febrero de 2027 con pago «Dom 28 feb 2027»; semana de lunes a domingo | fecha simulada en octubre de 2026 y febrero de 2027 |
| A-13 | Karen: solo en tarjeta 3, D-3, D-7, Calendario, las columnas de su parte en la tabla de T-2 y el filtro de familia; ningún indicador de su familia | revisión de textos con «Karen» |
| A-14 | Formato: dinero con `tnum` y signo menos tipográfico; fechas `10-oct-2026` o con día de la semana; 0 textos en inglés | axe-core y búsqueda de textos |
| A-15 | Tema claro y oscuro con los tokens de §10; axe-core sin violaciones AA | axe-core en los dos temas |
| A-16 | D-3 el 10-oct-2026 con la semilla: 3 pagos («Vie 30 oct 2026», «Dom 15 nov 2026», «Lun 30 nov 2026»), ninguno vencido y sin petición de vencidos | petición de red y texto |
| A-17 | Estados de §1.4 con el pago del 30 de octubre: el Jue 29 oct 2026 «Vence mañana»; el Vie 30 oct 2026 «Vence hoy»; el Sáb 31 oct 2026 sin marcar, «Vencido» con «Sin marcar como pagado»; con solo la parte de Karen pagada, «Falta la parte del hogar»; sin respuesta de `getAlertSettings`, ningún «Por vencer» | reloj simulado y respuestas simuladas |

---

## 12. Capturas para U5

Cada captura en claro y en oscuro, a **412 × 915** y **1440 × 900** (U1 §6.1), con reloj simulado:
1. Resumen con datos de la semilla el 10-oct-2026 (tarjetas 1, 3 y 5 con cifras; 2, 4 y 6 con su vacío).
2. Resumen el Dom 11 oct 2026 y el Mié 14 oct 2026 (tarjeta 4 y destino P-10, K38).
3. Resumen en carga (esqueletos), con error en una tarjeta y sin conexión con «· [n] en cola».
4. Resumen a 360 × 800 con la prueba de ancho de la cifra.
5. Detalle completo, incluido D-7 «Participación de Karen» con la nota fija.
6. Tendencias con T-1 a T-4 y una gráfica en «Ver como tabla».
7. Calendario: agenda (412) y mes con panel del día (1440) de octubre de 2026; mes de febrero de 2027.

---

## 13. Trazabilidad con los criterios U3

| Criterio | Dónde se cumple |
|---|---|
| U3-01 · 6 tarjetas | §2 (6 bloques «Tarjeta»), Tarjeta 4 (K38, U1 §0.2, casos domingo y entre semana), §11.2 A-01 y A-03 |
| U3-02 · métricas C1 | §2 a §5 (M-xx, `operationId`, campos, fórmula y periodo en cada bloque); verificación en el reporte |
| U3-03 · estados | §1.3, §7 (18 bloques × vacío, carga, error, sin conexión) |
| U3-04 · móvil primero | §8 (360/412, 768/1024, 1440), §1.2 (táctil 44 px) |
| U3-05 · Detalle, Tendencias, Calendario | §3, §4, §5 (propósito, contenido, métricas, filtros, estados); quincenas K8 y febrero en §5; 0 donas (§4); T-2 con las dos capas del mismo origen (PROVISIONAL, `docs/estado/bloqueos/U3-t2-acumulados.md`) |
| U3-06 · Karen | §3.1, §9 |
| U3-07 · próximo pago y A/B | Tarjeta 3, D-3 |
| U3-08 · cifras trazables | ejemplos con R-nn o K-nn (oct-2026 a mar-2027 = P1 a P6, regla de ejemplos al inicio); marcadores sin número; verificación en el reporte |
| U3-09 · formatos | reglas al inicio, textos de §2 a §7 |
| U3-10 · filtros | §6 |
| U3-11 · tokens y gráficas | §1.2, §4 (con el resumen y el `<figcaption>` de DESIGN §7.7), §10 |
| U3-12 · implementación y supuestos | §11, §14, §15 y el bloqueo `docs/estado/bloqueos/U3-t2-acumulados.md` |

---

## 14. Supuestos

1. Textos de rangos cortos en tarjetas («Del 16 al 31 oct», «Lun 12 al Dom 18 oct») y de meses en gráficas («oct-2026»): derivados del formato de U1 §0.4; no cambian cifras, esquema, operaciones ni tokens.
2. Hora de carga sin conexión como «[HH:MM]» en 24 horas (la misma forma que `send_time`, K43); solo texto.
3. Las tarjetas del Resumen no llevan filtros (vista de la casa); los filtros viven en Detalle, Tendencias y Calendario (§6).
4. La tarjeta 3, con varias líneas revolventes, muestra la del próximo pago y «+[n] líneas» (las demás en D-3).
5. (Retirado en la versión 1.1: T-2 sigue DESIGN §8.3, con la capa 2 en cero en los tramos de pago adelantado, y ya no cambia a dos líneas.)
6. Prioridad entre estados de pago (Pagado, Vencido, Por vencer, Parcial, Pendiente) y la línea «Falta la parte de …» en un pago vencido o por vencer (§1.4): DESIGN §7.8 da los estados pero no su prioridad; solo orden de presentación y texto.
7. Textos «Falta la parte del hogar» y «Vence en [d] días» (§1.4), para los casos que DESIGN §7.8 no nombra; solo texto.
8. Línea de meta de T-4b con el color de la persona y trazo de puntos 2 · 4 sin marcadores, distinta de los trazos de Axel (continuo) y Ximena (8 · 4); DESIGN §8 no fija la línea de meta; solo dibujo.
9. Redacción de los resúmenes de una frase de T-1 a T-4 (§4) y de los rótulos de la variante (b) de T-2: DESIGN §7.7 pide el resumen pero no su texto; U3 lo arma solo con cifras de C1, sin comparaciones; solo texto.

Ninguno cambia cifras, esquema, operaciones ni tokens.

## 15. Dependencias y abiertos
- DESIGN §7.6 a §7.8 y §8 ya existen y U3 los sigue desde la versión 1.1. Solicitudes a DESIGN (U2b), con su razón (detalle en `docs/estado/wf1/correccion-U3-r1.md`):
  1. §7.6: la tarjeta entera lleva a la pantalla destino de su acción primaria, no a la pestaña Detalle (U3-01: acción primaria única y la tarjeta 4 a P-10); y la anatomía suma el icono de 20 px junto a la etiqueta (§5.2; §5.4 ya prevé iconos en tarjetas).
  2. §7.8: la línea «Sin marcar como pagado» junto a «Vencido» (U1 N-022, C1 M-03), la prioridad entre estados y los textos «Falta la parte del hogar» y «Vence en [d] días» (U3 §1.4, supuestos 6 y 7).
  3. §8.3: quitar del subtítulo el rótulo «Deuda viva …» o justificarlo frente a K19, N-091 y M-11 («no se rotula deuda viva»).
  4. §8: fijar la línea de meta en una gráfica de líneas (U3 T-4b usa el color de la persona y trazo de puntos 2 · 4 sin marcadores, supuesto 8).
  5. §7.7: el ejemplo de resumen («El gasto subió 12 % respecto a septiembre») pide calcular una variación que la interfaz no hace (U1 §0.2). U3 usa resúmenes que solo enuncian cifras de C1 (§4). Pedido: cambiar el ejemplo o aclarar que la variación tiene que llegar calculada desde C1.
- C1-datos §9: dos métricas para D-4 y D-5, «Comidas de hoy» (N-121 vía N-027) y «Diferencias de precio» (N-162 vía N-027), con la fórmula, el periodo y los filtros que U3 §3 escribe mientras tanto (hallazgo H-7).
- C1 (`getAccrualVsCash`, esquema `AccrualVsCash.months[]` en `openapi.yaml`, y C1-datos M-11): agregar `cumulative_accrued_cents` y `cumulative_cash_cents`, con el mismo origen que `cumulative_difference_cents` (variante (a) de T-2; bloqueo `docs/estado/bloqueos/U3-t2-acumulados.md`).
- C1: la tarjeta 5 no puede nombrar la segunda tienda (`DashboardBestBasket` no trae su `store_id`); hoy se ve en D-5. Solicitud en el reporte; no bloquea.
- Bloqueos de S2 que tocan cifras del dashboard (U3 solo muestra lo que llegue y marca PROVISIONAL): `docs/estado/bloqueos/S2-presupuesto-total.md` (tarjeta 1), `docs/estado/bloqueos/S2-flujo-quincenal.md` (tarjeta 2, D-2), `docs/estado/bloqueos/S2-metas.md` (tarjeta 6, D-6), `docs/estado/bloqueos/S2-proyeccion-interruptor.md` (T-2, meses 7 a 12 con el interruptor encendido). Ninguno cambia la estructura, los textos fijos ni las operaciones de U3.
- Bloqueo propio de U3: `docs/estado/bloqueos/U3-t2-acumulados.md` (capas de T-2, sus etiquetas finales y su resumen; PROVISIONAL con dos variantes; se decide en la puerta 2). El bloqueo abierto `S1-parametros-avisos` no toca el dashboard: U3 no lee ni escribe `send_time`, `advance_days`, `weekday` ni `day_of_month`; los días de «Por vencer» son `alerts.payment_due_days_before`, ya decididos (decisión F, K44).
