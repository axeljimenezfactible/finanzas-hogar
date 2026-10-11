# U1 · Arquitectura de información · Finanzas Hogar

| | |
|---|---|
| Fase | U1 «Arquitectura de información» (Pista 2, flujo WF-0) |
| Agente | U1-arquitectura (v1.0) · U1-1.1 (v1.1) · U1-1.2 (v1.2) |
| Versión | **1.2** (11-oct-2026): respuestas de la puerta 1b (`decisiones.md` K35 a K44) sobre la 1.1 (respuestas de la puerta 1, §K, y hallazgos R2-05 y R3-01) |
| Fecha | 10-oct-2026 (v1.0) · 11-oct-2026 (v1.1 y v1.2) |
| Estado | v1.0 aprobada en la puerta 1 (K34); v1.1 cerrada en el flujo P1; v1.2 para revisión |
| Fuentes | `plan.md` §1 (HECHOS), §2, §3.4 (S2 a S7 y 3.4.1), §3.5 (U1, U3, U4, U5), §3.7; `decisiones.md` A a G (y H solo para IA opcional y PWA); en v1.1, también §K y `docs/criterios/U1.md` 1.1; en v1.2, K35 a K44 |
| Fuera de alcance | Colores, tipografías, espaciados y tamaños (U2, contrato C2). Endpoints y esquemas (C1, S1, S6). Campos y validaciones finas de Ajustes (U4). |

Este documento dice **qué pantallas existen, qué muestra cada una y qué datos necesita**. C1 dice **cómo** se piden esos datos: cada necesidad `N-xxx` de la sección 3 se mapea a uno o más endpoints en C1.


## Cambios de la versión 1.2 (11-oct-2026)

Aplica las respuestas de Axel en la puerta 1b que tocan a U1 (`decisiones.md` K35 a K44) y decide el pendiente «Descartar chat detectado» (`docs/estado/p1/INTEGRACION.md` §7, punto 7). No cambian la barra lateral, las 27 pantallas, las 14 secciones de Ajustes, el camino de 3 toques de F-1 ni ninguna cifra del oráculo. Necesidades: **211** (antes 203): 138 de lectura y 73 de escritura; las nuevas son N-320 a N-327 (números a continuación de la última, sin reusar).

| # | Cambio | Origen | Dónde |
|---|---|---|---|
| 1 | Semana por defecto: los domingos, P-10, P-13 y P-14 abren la semana que empieza el lunes; «Cambiar semana» siempre disponible. Pasa a **FIRME** y se quitan sus marcas PROVISIONAL | K38 | §0.2, §7 (N-120, N-150 y N-160 remiten a §0.2) |
| 2 | Pasan a **FIRMES** la definición de «casillas sin armar» (N-023) y los valores de C1 que eran propuesta: adjunto de hasta 10 MB (N-064), comparador de hasta 5 s (N-226), hora `send_time` de reportes y avisos y anticipación `advance_days` del aviso de precios viejos (N-297, N-298, N-300; la del pago próximo es `alerts.payment_due_days_before`, valor inicial [1, 0], decisión F, y se edita en N-323) | K43 | N-023, N-064, N-226, N-297, N-298, N-300, §7 |
| 3 | Todo parámetro operativo existe y se cambia en la sección de Ajustes a la que pertenece, sin sección nueva: valor inicial visible, validación «sin valor inventado» y «Restablecer al valor inicial»; una necesidad por grupo de parámetros | K44 (ubicación por sección: lectura PROPUESTA del orquestador; Axel la ve en la puerta 2) | §0.2 (precio viejo), §2.3 (P-26.01, P-26.05, P-26.06, P-26.09, P-26.11, P-26.13 y tabla «Parámetros operativos»), §2.6, N-320 a N-323, N-325, N-326, §7 |
| 4 | Ajuste de la meta de kcal por objetivo: valores de la casa en Ajustes › Salud y metas de ahorro (0 y 0 kcal) y ajuste por persona en P-21 (valor inicial «el de la casa», «Restablecer»); «mantener» no ajusta | K36 | §2.2 (P-21), §2.3 (P-26.13), §2.6 (P-21), N-230, N-311, N-326, N-327, §7 |
| 5 | «Descartar chat detectado» existe: borra un chat detectado que nadie va a confirmar | U1 (INTEGRACION §7, punto 7) | §2.3 (P-26.09), N-299, N-324 |
| 6 | La Pista 3 (Claude Design) la ejecuta Axel con el paquete de `docs/design/entrada/` | K35 | §7 |
| 7 | Matriz §3.29 con N-320 a N-327 y conteos rehechos por script. Verificación: `python3 -I docs/design/tools/verificar_u1.py` | — | §3.29 |
| 8 | Corrección r1 (H-01 de `docs/estado/p1c/REVISION-U1-r1.md`): la anticipación del aviso de pago próximo se nombra con su clave del esquema, `days_before` (lista de días antes del pago, valor inicial [1, 0], decisión F), y `advance_days` queda solo para el aviso de precios viejos, como en `00003_puerta1b.sql` | Revisión r1 | §2.3 (tabla de parámetros), N-298, N-300 |
| 9 | Corrección r2 (H-r2-01 de `docs/estado/p1c/REVISION-U1-r2.md`): coherencia con el esquema 1.3 tal como quedó tras la corrección S1-r1 (H-01 de `docs/estado/p1c/REVISION-S1-r1.md`). Los umbrales del presupuesto y los días antes del pago próximo son claves de `app_setting` (`alerts.budget_warning_pct` 80, `alerts.budget_limit_pct` 100 y `alerts.payment_due_days_before` [1, 0]; decisión F, K44) y ya no van en `notification.params_json`; los lee, guarda y restablece N-323. N-300 se queda con `send_time` y `advance_days`, cuya ubicación decide el bloqueo abierto `S1-parametros-avisos`. «Para C1» corregida: `docs/estado/p1c/correccion-U1-r2.md` | Revisión r2 | §2.3 (tabla de parámetros), N-298, N-300, N-323, §7 |

## Cambios de la versión 1.1 (11-oct-2026)

Aplica las respuestas de Axel en la puerta 1 que tocan a U1 (`decisiones.md` §K) y los hallazgos R2-05 y R3-01 de la revisión de WF-0, con los criterios U1 1.1 (U1-04, U1-11, U1-15, U1-16). No cambian la barra lateral, las 27 pantallas, las 14 secciones de Ajustes ni ninguna cifra del oráculo. Necesidades: **203** (antes 199): 138 de lectura y 65 de escritura; las nuevas son N-316 a N-319 (números a continuación de la última, sin reusar).

| # | Cambio | Origen | Dónde |
|---|---|---|---|
| 1 | Semana de lunes a domingo; quincenas del 1 al 15 y del 16 al último día del mes: **FIRME** (se quitan las marcas PROVISIONAL) | K8 | §0.2, N-120, §7 |
| 2 | «Semana por defecto» y «lista de compras vigente» definidas en un solo lugar; toda lectura y escritura de la lista lleva la semana que muestra la pantalla. La regla del domingo queda PROVISIONAL (bloqueo `U1-lista-vigente-domingo`) | R3-01, U1-16 | §0.2, §2.5 (plan_comida), N-073, N-120, N-127, N-150 a N-153, N-160, N-165, F-2, §7 |
| 3 | Filtro «familia» con tres valores incluyentes: Hogar (con lo conjunto), Karen (con lo conjunto) y Conjunto | K5 | §0.2, N-040 |
| 4 | Estado «Sin conexión»: gastos y peso se siguen capturando a una cola en el teléfono que se envía sola y sin duplicar al volver a la red de casa; las demás pantallas quedan en solo lectura o con aviso. El camino de 3 toques no cambia | K9, U1-15 | §1.4 (G-07), §2.4, N-004, N-006, N-054, N-055, N-240, N-317 (nueva, Local), F-1, F-6, §7 |
| 5 | Regla de conteo de toques: opción (a), confirmada por Axel | K10, U1-04 | §4, F-1 |
| 6 | El cambio de monto del presupuesto lleva alcance «Solo este mes» o «De este mes en adelante» (preseleccionado) | K2 | §2.2 (P-06), N-087, N-275, N-278 |
| 7 | «Guardar para después» conserva las correcciones del ticket como borrador | K3 | §2.2 (P-05), N-071, N-075, N-076, F-8 |
| 8 | Chat de Telegram: el bot detecta el chat cuando la persona le escribe y ella lo confirma en Ajustes › Notificaciones | K4 | §2.3 (P-26.09), N-299, N-316 (nueva) |
| 9 | Tasas tal como se capturaron, con sus decimales; los cálculos usan la fracción exacta | K12 | §0.4, N-204, N-286, N-287, N-289 |
| 10 | Tope de pago por quincena opcional por persona, en Ajustes › Personas y reparto (sale de Créditos y condiciones) | K13 | §2.3 (P-26.12), N-225, N-288, N-309, N-310 |
| 11 | Restricciones duras globales del comparador en Ajustes › Créditos y condiciones, ajustables en cada comparación y sin valor por defecto; P-20 avisa «Faltan restricciones» y no recomienda | K16 | §2.2 (P-20), §2.3 (P-26.05), N-225, N-226, N-288, N-289, F-7 |
| 12 | Interruptor de proyección «Repetir la despensa bimestral y el papel» en Ajustes › Presupuesto, apagado por defecto; su regla (N-318) rige los meses 7 a 12 de P-07, por quincena en N-090 y por mes en N-091 | K21 | §2.3 (P-26.02), N-090, N-091, N-093, N-318 y N-319 (nuevas), §7 |
| 13 | Meta de kcal sugerida con objetivo (bajar, mantener o subir de peso) que la persona confirma o cambia | K23 | §2.2 (P-21), §2.3 (P-26.13), N-230, N-311, N-313 |
| 14 | Objetivo táctil de al menos 44 × 44 px; los 24 × 24 px de WCAG 2.5.8 quedan solo como piso | R2-05, U1-11 | §6.1 punto 4 |
| 15 | Se quitan otras marcas que cerraron las respuestas: precio viejo a los 30 días o más (K22) y relación entre diferencia acumulada y deuda viva (K19). Sigue PROVISIONAL la definición de «casillas sin armar» de N-023 hasta que Axel apruebe el valor (K31) | K19, K22, K31 | §0.2, N-023, N-091, §7 |
| 16 | Matriz §3.29: N-313 en P-21; N-316 a N-319 en sus secciones; N-317 en Global, P-03 y P-22. Verificación: `python3 -I docs/design/tools/verificar_u1.py` | — | §3.29 |
| 17 | Tarjeta bancaria y crédito de tienda: por línea, Axel captura desde su estado de cuenta la tasa anual (tal como se capturó, K12), la base diaria o mensual, si incluye IVA, el porcentaje de pago mínimo y el mínimo fijo, sin valor por defecto. Mientras falte alguno, la línea muestra «Falta capturar» (abre P-26.05 en esa línea) y S3 no calcula sus intereses ni su pago mínimo con valores supuestos (RN-26) | K18, K12 | §2.2 (P-18), §2.3 (P-26.05), §2.4, N-200, N-201, N-205, N-283, N-289, §7 |

---

## 0. Convenciones

### 0.1 Identificadores

| Prefijo | Qué es | Ejemplo |
|---|---|---|
| `P-xx` | Pantalla | P-01 Inicio |
| `P-26.xx` | Sección de Ajustes (14) | P-26.05 Créditos y condiciones |
| `G-xx` | Componente global, presente en todas las pantallas | G-03 Menú «/» |
| `N-xxx` | Necesidad de datos. Tipo **L** (lectura) o **E** (escritura). Numeración por bloques de pantalla con huecos reservados | N-020 |
| `F-x` | Flujo clave | F-1 Registrar gasto |

### 0.2 Vocabulario de datos (rige todas las `N-xxx`)

- **Hoy**: fecha local en `America/Mexico_City`. Toda necesidad «Diaria» cambia a las 00:00 de esa zona.
- **Mes en curso**: el mes calendario que contiene hoy. Es el periodo base del presupuesto (decisión B).
- **Quincena** (**FIRME**, K8): Q1 del día 1 al 15; Q2 del día 16 al último día del mes. El ingreso de Axel y los pagos caen el 15 y el 30; en febrero, el último día (28 o 29).
- **Semana** (**FIRME**, K8): de lunes a domingo (ISO 8601). Toda necesidad «por semana» usa esta semana.
- **Hogar**: Axel y Ximena. Equivale a la cuenta `axel` de la semilla (`data.json`) y a la parte «Axel» del modelo v3. **Familia** (filtro de N-040 y de las vistas que usan sus filtros; K5): tres valores incluyentes: **Hogar** = gastos con reparto Hogar o Conjunto; **Karen** = Karen o Conjunto; **Conjunto** = solo los de reparto Conjunto.
- **Reparto** de un gasto, renglón o compra: Hogar (100 %), Karen (100 %) o Conjunto con un % para el hogar y el resto para Karen (50/50 por defecto, decisión B). La semilla trae además renglones «mixto»; S1 define su representación.
- **Parte del hogar / parte de Karen**: lo que corresponde a cada uno según el reparto, con la regla de `reglas-negocio.md` (S1), que garantiza que la suma de las partes es el total. **Toda cifra de presupuesto, flujo y ahorro usa la parte del hogar.** La parte de Karen es informativa (decisión A) y solo aparece donde ella participa.
- **Devengado y efectivo**: según `reglas-negocio.md` (S1, S2). El presupuesto (real contra plan) se mide en devengado, porque la despensa bimestral se devenga a mitad por mes (decisión B). El flujo por quincena se mide en efectivo.
- **Pago pendiente**: pago programado con al menos una parte (hogar o Karen) sin marcar como pagada. Cada parte se marca por separado porque cada quien paga directo a Kueski (decisión A).
- **Precio viejo**: precio cuya fecha tiene 30 días o más respecto a hoy (decisión E y K22: a los 30 días o más; el precio del 10-oct-2026 ya es viejo el 9-nov-2026). Si no, **fresco**. Los 30 días y los 90 en que el peso de un precio llega a cero (K40) son parámetros de Ajustes › Tiendas y precios con esos valores iniciales (K44, N-322).
- **Semana por defecto y lista de compras vigente** (única definición; R3-01, U1-16): P-10, P-13 y P-14 abren por defecto la **semana en curso** (lunes a domingo, K8) y, **los domingos, la semana siguiente**, la que empieza el lunes (**FIRME**, K38; «Cambiar semana» siempre disponible). La **lista de compras vigente** es la lista de esa semana. El cambio ocurre a las 00:00 de `America/Mexico_City`; una pantalla abierta conserva la semana que muestra hasta que la persona cambie de semana o vuelva a entrar. **Toda lectura y escritura de la lista lleva la semana (`week` en C1, obligatoria) que la pantalla muestra**: ninguna acción actúa sobre una semana que no se ve. El plan de la semana (P-10) recibe su semana en el filtro `period` de C1 y, si no lo manda, C1 aplica esta misma regla (`getMealPlanWeek`, K38).
- **Motor Sx**: la cifra la calcula ese motor (S2 financiero, S3 créditos, S4 recetas y comidas, S5 precios y tiendas, S7 reportes). La interfaz no recalcula reglas de negocio: muestra, filtra, ordena y suma lo que está a la vista.
- **Sin datos no es cero**: si falta una entrada (sin plan, sin precio, sin peso, sin meta), se muestra «—» con el motivo, nunca `$0.00` ni `0 %`.

### 0.3 Frecuencias de actualización

| Clave | Significado |
|---|---|
| Inmediata | Se calcula al consultar y refleja cualquier captura anterior. Después de guardar algo, la pantalla abierta vuelve a pedir sus necesidades. |
| Diaria | Inmediata y además cambia a las 00:00 de `America/Mexico_City` aunque nadie capture (días que faltan, precios viejos, «hoy»). |
| Al generar | Resultado de un trabajo (OCR, PDF, importación PROFECO, respaldo, comparador de escenarios). Queda con su fecha y hora de generación. |
| Programada | Ocurre a la hora configurada en P-26.08 o P-26.09 (resumen diario de noche, semanal el domingo, mensual el día 1). |
| Local | Vive en el dispositivo y no requiere API: tema, barra lateral plegada, vista elegida, recientes. |

### 0.4 Formatos en pantalla (decisión B)

- **Dinero**: `$1,258.68`. MXN, dos decimales, cifras tabulares y alineadas a la derecha en tablas.
- **Fechas**: `10-oct-2026`. Donde importa el día de la semana (calendarios, pagos, plan de comidas): `Vie 30 oct 2026`. En archivos y API, ISO `2026-10-10`.
- **Porcentajes de avance** (presupuesto, metas): entero con ROUND_HALF_UP (`80 %`). **Tasas**: tal como se capturaron, con el mismo número de decimales (`17.27 %`, `25.8685 %`); los cálculos usan la fracción exacta (K12).
- **Idioma**: español de México en todo texto visible.

### 0.5 Reglas de personas

- Axel y Ximena son co-titulares: **ven lo mismo y pueden hacer lo mismo**. Ninguna pantalla, filtro o permiso oculta datos según quién inició sesión. El usuario de la sesión solo se usa como valor por defecto (persona del gasto, persona del peso) y para «creado por».
- **Karen aparece solo como participante**: reparto de gastos y renglones de la despensa común, partes de sus compras con el crédito Kueski y de sus pagos, el bloque «Participación de Karen» de Inicio › Detalle, su estado de cuenta (P-25) y Ajustes › Personas y reparto. No hay pantallas de sus finanzas ni de su familia, ni ingresos suyos.
- Plan de comidas y salud: solo Axel y Ximena (decisión A).

---

## 1. Mapa del sitio y barra lateral (a)

### 1.1 Árbol del sitio

```
Finanzas Hogar  (G-01 barra lateral)
├─ Buscar                         → G-04 (Ctrl+K)
├─ Registrar gasto                → P-03 (hoja)
├─ Inicio                         → P-01  [Resumen* · Detalle · Tendencias · Calendario]
├─ Gastos ▸
│  ├─ Todos los gastos            → P-02  [Resumen · Detalle* · Tendencias · Calendario]
│  │   └─ (cada gasto)            → P-04
│  └─ Capturar ticket             → P-05
├─ Presupuesto ▸
│  ├─ Presupuesto del mes         → P-06  [Resumen* · Detalle · Tendencias · Calendario]
│  ├─ Flujo y proyección          → P-07  [Resumen* · Detalle · Tendencias]
│  ├─ Ingresos                    → P-08  [Resumen · Detalle* · Tendencias · Calendario]
│  └─ Ahorro y metas              → P-09  [Resumen* · Detalle · Tendencias · Calendario]
├─ Comidas ▸
│  ├─ Plan de la semana           → P-10  [Resumen · Detalle* · Tendencias · Calendario]
│  └─ Recetario                   → P-11
│      └─ (cada receta)           → P-12
├─ Compras y tiendas ▸
│  ├─ Lista de compras            → P-13
│  ├─ Comparar canasta            → P-14  [Resumen* · Detalle]
│  ├─ Productos y precios         → P-15  [Resumen · Detalle* · Tendencias]
│  └─ Inventario                  → P-16  [Resumen · Detalle*]
├─ Créditos ▸
│  ├─ Todas las líneas            → P-17  [Resumen* · Detalle · Tendencias · Calendario]
│  ├─ (cada línea, por su nombre) → P-18  [Resumen* · Detalle · Tendencias · Calendario]
│  │   └─ (cada compra)           → P-19
│  └─ Comparador de escenarios    → P-20
├─ Salud                          → P-21  [Resumen* · Detalle · Tendencias]
├─ Reportes ▸
│  ├─ Centro de reportes          → P-23
│  ├─ Reporte mensual             → P-24
│  └─ Estado de cuenta de Karen   → P-25
└─ Ajustes ▸                      → P-26
   ├─ 01 Cuenta y roles
   ├─ 02 Presupuesto
   ├─ 03 Ingresos
   ├─ 04 Categorías
   ├─ 05 Créditos y condiciones (incluye promociones)
   ├─ 06 Tiendas y precios (incluye mercados y membresía Sam's)
   ├─ 07 Recetas
   ├─ 08 Reportes y frecuencias
   ├─ 09 Notificaciones
   ├─ 10 Apariencia
   ├─ 11 Datos y respaldos
   ├─ 12 Personas y reparto
   ├─ 13 Salud y metas de ahorro
   └─ 14 Inventario

Fuera de la barra: P-00 Iniciar sesión · P-22 Registrar peso (hoja que se abre desde P-21, «+» o /peso).
* = subsección que abre por defecto.
```

### 1.2 Jerarquía y reglas de la barra lateral (G-01)

1. **Estructura tipo Notion**: encabezado con el nombre del espacio («Finanzas Hogar») y el usuario de la sesión (menú: Cuenta, Apariencia, Cerrar sesión); dos acciones fijas (Buscar, Registrar gasto); la lista de secciones con triángulo ▸ para desplegar subpáginas; pie con el interruptor de tema.
2. **Orden fijo** de secciones, el de la decisión: Inicio, Gastos, Presupuesto, Comidas, Compras y tiendas, Créditos, Salud, Reportes, Ajustes. No se reordenan ni se ocultan (Axel y Ximena ven lo mismo).
3. **Tocar el nombre** de una sección abre su primera subpágina (Gastos → Todos los gastos, Presupuesto → Presupuesto del mes, Comidas → Plan de la semana, Compras y tiendas → Lista de compras, Créditos → Todas las líneas, Reportes → Centro de reportes, Ajustes → 01 Cuenta y roles en escritorio o la lista de secciones en móvil). **Tocar ▸** solo despliega o pliega.
4. **Desplegado por defecto**: solo la sección de la página abierta. Ajustes siempre empieza plegado. El estado plegado o desplegado se recuerda por dispositivo (Local).
5. **Hijos dinámicos**: solo Créditos lista sus líneas por nombre (son pocas; N-005). Gastos, recetas, productos y compras son registros de base de datos y se abren desde su vista, no desde la barra.
6. La página activa va resaltada; las migas de pan de la barra superior (G-02) siguen la misma jerarquía.
7. **Colapsable en escritorio** (botón «‹‹» o Ctrl+\): al colapsar se oculta y reaparece como panel flotante al llevar el puntero al borde izquierdo.

### 1.3 Qué se colapsa en móvil (< 768 px)

| Elemento | Escritorio (≥ 1024 px) | Móvil (< 768 px) |
|---|---|---|
| Barra lateral | Fija y colapsable | Oculta. Se abre como panel desde ☰ (arriba a la izquierda). No se abre deslizando desde el borde, porque choca con el gesto «atrás» de Android |
| Navegación principal | Barra lateral | Barra inferior fija (G-05): Inicio · Comidas · «+» · Compras · Buscar |
| Subpáginas | Visibles al desplegar | Dentro del panel; solo la sección actual desplegada |
| Migas de pan | Completas | «‹ Página padre» |
| Pestañas Resumen / Detalle / Tendencias / Calendario | Fila de pestañas | Fila con desplazamiento horizontal; la activa siempre visible |
| Acciones secundarias | Botones en el encabezado | Menú «⋯» del encabezado; la primaria sigue visible |
| Propiedades de un registro | Todas | Las 4 primeras + «Mostrar N más» |
| Filtros, orden y agrupación | Barra de la vista | Botón «Filtrar» que abre una hoja inferior; chips con los filtros activos bajo el título |
| Vista previa de un registro (G-06) | Panel lateral derecho | Página completa con «‹» |
| Hojas y diálogos | Ventana modal | Hoja inferior a pantalla completa |

### 1.4 Componentes globales

| Id | Componente | Qué hace |
|---|---|---|
| G-01 | Barra lateral | Navegación jerárquica (1.2) |
| G-02 | Barra superior de página | Migas, título, propiedades clave, acción primaria, «⋯» con secundarias; en escritorio, también el acceso a búsqueda |
| G-03 | Menú «/» | Comandos para registrar, ir a, cambiar vista y herramientas (sección 5) |
| G-04 | Búsqueda | Ctrl+K o Ctrl+P en escritorio; pestaña «Buscar» en móvil. Un solo panel: resultados de búsqueda y, si el texto empieza con «/», comandos |
| G-05 | Barra inferior (solo < 768 px) | Cinco destinos; «+» abre P-03 |
| G-06 | Vista previa de registro | Abre un registro sin salir de la vista: panel lateral en escritorio, página completa en móvil |
| G-07 | Avisos breves y estados globales | Aviso «Guardado · Deshacer» (5 s); franja «Sin conexión con el servidor de casa» con el número de capturas en cola (N-317, §2.4); sesión vencida → P-00; registro o página inexistente |

---

## 2. Inventario de pantallas (b)

**Regla de subsecciones** (decisión G, dashboard minimalista): **Resumen** reúne como máximo 6 tarjetas con las cifras clave; **Detalle** es la base de datos de la página con sus vistas tabla, tablero o galería (o el contenido principal si no es base de datos); **Tendencias** son gráficas en el tiempo; **Calendario** es la vista calendario de la misma base de datos. El detalle se pide bajo demanda: tocar una tarjeta lleva a su detalle. No todas las pantallas tienen las cuatro.

### 2.1 Qué es cada pantalla

| Id | Nombre | Sección | Propósito | Subsecciones (* por defecto) | Vistas de base de datos | Componentes clave |
|---|---|---|---|---|---|---|
| P-00 | Iniciar sesión | — | Entrar con usuario y contraseña; sesión larga por dispositivo | — | — | Formulario usuario y contraseña, mensaje de bloqueo |
| P-01 | Inicio | Inicio | Ver el estado del hogar en 6 tarjetas y bajar al detalle | Resumen* · Detalle · Tendencias · Calendario | Calendario (eventos del hogar) | 6 tarjetas en orden fijo: Mes vs plan, Flujo quincenal, Próximo pago y disponible, Comidas de la semana, Mejor canasta, Ahorro. Bloque «Participación de Karen» (Detalle). 4 gráficas (Tendencias). Calendario del hogar |
| P-02 | Todos los gastos | Gastos | Consultar, filtrar y corregir gastos | Resumen · Detalle* · Tendencias · Calendario | Tabla*, tablero (por categoría), calendario, galería (tickets) | Barra de vista con filtros (persona, familia, categoría, tienda, forma de pago, tipo, periodo), fila de totales, edición en el lugar |
| P-03 | Registrar gasto | Hoja global | Registrar un gasto en ≤ 3 toques | — | — | Selector de tipo (Gasto · Ingreso · Peso · Ticket), monto con teclado numérico, 6 chips de categoría + «Más», «Más opciones» (fecha, persona, reparto, forma de pago, tienda, nota, adjunto), Guardar |
| P-04 | Gasto | Gastos | Ver y editar un gasto, su reparto y su ticket | — | Tabla (renglones) | Propiedades editables, bloque de reparto con partes hogar y Karen, renglones del ticket, adjunto, nota (bloques) |
| P-05 | Capturar ticket | Gastos | Foto del ticket → OCR local → confirmar → crea gasto, precios e inventario | — (pasos 1 a 4) | Tabla (renglones) | Cámara o galería, foto visible, renglones editables con nivel de confianza, emparejado con catálogo, cuenta por renglón, forma de pago, tickets por confirmar |
| P-06 | Presupuesto del mes | Presupuesto | Plan contra real por categoría y por quincena | Resumen* · Detalle · Tendencias · Calendario | Tabla* (categorías), tablero (normal, aviso ≥ 80 %, excedido ≥ 100 %) | Selector de mes, barras de avance, alertas 80/100, despensa devengada |
| P-07 | Flujo y proyección | Presupuesto | Efectivo por quincena, proyección de 12 meses devengado contra efectivo y deuda viva a una fecha | Resumen* · Detalle · Tendencias | Tabla (quincenas y meses) | Tabla de quincenas, gráfica devengado contra efectivo, selector de fecha para deuda viva, supuestos visibles |
| P-08 | Ingresos | Presupuesto | Capturar las quincenas de Axel y registrar extras | Resumen · Detalle* · Tendencias · Calendario | Tabla*, calendario | Filas de quincena (esperada o capturada), extras |
| P-09 | Ahorro y metas | Presupuesto | Fondo, metas con fecha y aportaciones | Resumen* · Detalle · Tendencias · Calendario | Galería* (metas), tabla (movimientos), tablero (por estado), calendario (fechas meta) | Barras de avance, estado en ruta, atrasada o cumplida |
| P-10 | Plan de la semana | Comidas | Asignar recetas por día, tiempo y persona; ver costo y kcal | Resumen · Detalle* · Tendencias · Calendario | Calendario* (semana: días × tiempos), tabla (asignaciones), tablero (por tiempo), galería (recetas de la semana) | Cuadrícula semanal, selector de semana, selector de receta con fila de días, totales de costo y kcal, casillas sin armar |
| P-11 | Recetario | Comidas | Buscar y crear recetas | — | Galería*, tabla, tablero (por tiempo de comida) | Tarjetas con costo y kcal por porción, filtros |
| P-12 | Receta | Comidas | Ingredientes exactos, escalado, costo y kcal | — | Tabla (ingredientes) | Selector de porciones, ingredientes con presentación comprable, fuente nutricional por ingrediente, pasos (bloques) |
| P-13 | Lista de compras | Compras y tiendas | Qué comprar y dónde, neteado con el inventario | — | Tabla*, tablero (por tienda o mercado) | Grupos por vertical (mercados y cadenas), casillas de comprado, totales por tienda |
| P-14 | Comparar canasta | Compras y tiendas | Mejor tienda para una canasta y por artículo; evaluación de la membresía de Sam's | Resumen* · Detalle | Tabla (artículos × tiendas) | Selector de canasta, totales por tienda, idénticos, sustitutos y sin precio, ganadores, membresía Sam's con punto de equilibrio |
| P-15 | Productos y precios | Compras y tiendas | Catálogo, precios por tienda, frescura y fuentes | Resumen · Detalle* · Tendencias | Tabla*, tablero (por frescura), galería | Precio por unidad base, marca de viejo, fuente, historial |
| P-16 | Inventario | Compras y tiendas | Qué hay en casa (stock por unidad y ritmo de consumo) | Resumen · Detalle* | Tabla*, tablero (por estado), galería | Estado suficiente, reponer o agotado; días restantes; movimientos |
| P-17 | Todas las líneas | Créditos | Deuda viva, disponible y calendario unificado de pagos | Resumen* · Detalle · Tendencias · Calendario | Galería* (líneas), tabla, calendario (pagos y promociones) | Tarjeta por línea, aviso A/B «pendiente confirmar», calendario unificado |
| P-18 | Crédito | Créditos | Estado de una línea según su tipo y su calendario | Resumen* · Detalle · Tendencias · Calendario | Tabla (compras, amortización o abonos), calendario (pagos) | Variante por tipo (revolvente, tarjeta, préstamo, tienda), serie de disponible A/B, tasas por plazo, marcar pago |
| P-19 | Compra a crédito | Créditos | Alta y detalle de una compra: intereses cotizados, plazo, partes y pagos | — | Tabla (pagos) | Formulario con vista previa del cálculo, partes por persona, pagos con estado por parte |
| P-20 | Comparador de escenarios | Créditos | Elegir línea, plazo, fecha y promoción para una compra planeada | — | Tabla (escenarios) | Datos de la compra, restricciones duras, pesos, 4 óptimos lado a lado, recomendación, descartados con motivo |
| P-21 | Salud | Salud | Perfil, IMC, meta de kcal, kcal del día y peso | Resumen* · Detalle · Tendencias | Tabla (pesos) | Dos columnas Axel y Ximena, gráfica de peso, kcal de hoy contra meta |
| P-22 | Registrar peso | Hoja (desde P-21 o «+») | Capturar peso en ≤ 3 toques | — | — | Persona, campo peso, último peso, Guardar |
| P-23 | Centro de reportes | Reportes | Reportes generados, próximos envíos y avisos enviados | — | Tabla*, calendario (envíos) | Lista de PDF, estado de envío por Telegram, historial de avisos |
| P-24 | Reporte mensual | Reportes | Vista web del mes con las mismas cifras del PDF | — | — | Selector de mes, secciones del reporte, descargar o generar PDF |
| P-25 | Estado de cuenta de Karen | Reportes | Vista previa y PDF de la parte de Karen | — | Tabla (movimientos) | Selector de periodo, compras y pagos de Karen, saldo, generar PDF, compartir |
| P-26 | Ajustes | Ajustes | Configurar el sistema en 14 secciones (2.3) | — | Según sección (2.3) | Columna interna con las 14 secciones (escritorio) o lista (móvil), edición en el lugar con guardado automático, vista previa del efecto |

### 2.2 Acciones y estados por pantalla

La acción primaria es única por pantalla y siempre visible (en móvil, en el encabezado o en la barra inferior). La edición en el lugar guarda sola, como en Notion; por eso en páginas de registro la primaria es la acción de negocio siguiente y no «Guardar».

| Id | Acción primaria | Secundarias | Vacío | Carga | Error |
|---|---|---|---|---|---|
| P-00 | Entrar | Ninguna (la recuperación de contraseña no se hace en pantalla; la define S6 en el runbook) | — | Botón en espera, campos bloqueados | «Usuario o contraseña incorrectos». Bloqueo: «Demasiados intentos. Vuelve a intentar a las HH:MM». Sin servidor: franja G-07 |
| P-01 | Registrar gasto | Foto de ticket; Capturar quincena (solo los días 15 y 30 si falta); cambiar subsección | Cada tarjeta se vacía por separado con texto guía y la acción que la llena (p. ej. «Sin presupuesto del mes · Definir presupuesto») | Esqueleto de las 6 tarjetas con su forma; nunca `$0.00` provisional | Por tarjeta: «No se pudo calcular · Reintentar»; las demás siguen visibles |
| P-02 | Registrar gasto | Foto de ticket; Filtrar; Ordenar; Agrupar; cambiar vista | «Sin gastos en oct-2026 · Registrar gasto». Con filtros: «Ningún gasto coincide · Quitar filtros» | Esqueleto de 8 filas | Fila de error con Reintentar. Edición fallida: el valor vuelve al anterior con aviso |
| P-03 | Guardar | Cambiar tipo (Ingreso, Peso, Ticket); Más opciones; Cerrar | No aplica (formulario) | Guardar en espera; la hoja no se cierra hasta confirmar | Lo capturado se conserva. Mensaje junto al campo (monto ≤ 0, falta categoría) o general «No se guardó · Reintentar» |
| P-04 | Adjuntar ticket (si no tiene) o Ver ticket | Duplicar; Marcar como pagado (solo venideros); Eliminar (con confirmación); Abrir compra a crédito | — | Esqueleto de propiedades | «Este gasto ya no existe» (lo borró el otro usuario) + Volver |
| P-05 | Confirmar ticket | Repetir foto; Agregar renglón; Guardar para después (conserva las correcciones como borrador, K3); Cancelar | «No hay tickets por confirmar» | «Leyendo el ticket…» con la foto visible; esqueleto de renglones | OCR sin texto: «No se pudo leer · Repetir foto o capturar a mano»; la foto se conserva |
| P-06 | Ajustar presupuesto (con alcance «Solo este mes» o «De este mes en adelante», K2) | Cambiar mes; Ver gastos de la categoría | «Sin presupuesto definido · Ajustar presupuesto» | Esqueleto de tabla | Reintentar |
| P-07 | Simular compra (abre P-20) | Elegir fecha de deuda viva; Ver supuestos | «Faltan datos para proyectar: ingreso esperado (Ajustes › Ingresos)» | Esqueleto | Reintentar |
| P-08 | Capturar quincena | Registrar ingreso extra | «Sin ingreso esperado · Ir a Ajustes › Ingresos» | Esqueleto | Reintentar |
| P-09 | Registrar aportación | Nueva meta (abre P-26.13); Registrar retiro | «Aún no hay metas · Nueva meta» | Esqueleto de tarjetas | Reintentar |
| P-10 | Generar lista de compras | Copiar semana anterior; Cambiar semana; Marcar «fuera»; Vaciar casilla | «Semana sin armar · Copiar semana anterior o elige recetas» | Esqueleto de cuadrícula | La casilla que no se guardó vuelve a vacía con aviso |
| P-11 | Nueva receta | Filtrar; Sugerir con IA (solo si está activa en P-26.11) | «Sin recetas que coincidan» | Esqueleto de galería | Reintentar |
| P-12 | Agregar al plan | Editar en el lugar; Duplicar; Eliminar | — | Esqueleto | Ingrediente sin precio: «—» y motivo. Sin kcal: «sin fuente». Ninguno bloquea la receta |
| P-13 | Terminar compra (abre P-05 prellenada) | Comparar tiendas; Agregar artículo; Cambiar semana | «Nada que comprar: el inventario cubre el plan» o «Sin plan esta semana · Planear» | Esqueleto | Reintentar |
| P-14 | Usar esta tienda (nombre de la ganadora) | Cambiar canasta; Simular pago con crédito (abre P-20) | «La canasta no tiene artículos» o «Sin precios para comparar · Capturar precios» | Esqueleto | Reintentar |
| P-15 | Capturar precio | Foto de ticket; Nuevo producto; Filtrar | «Sin productos» | Esqueleto | Reintentar |
| P-16 | Ajustar existencia | Registrar sobrante; Filtrar | «Inventario vacío: se llena al confirmar tickets» | Esqueleto | Reintentar |
| P-17 | Registrar compra a crédito | Simular compra; Ver calendario | «Sin créditos · Agregar en Ajustes › Créditos» | Esqueleto de tarjetas | Reintentar |
| P-18 | Marcar pago | Registrar compra; Simular compra; Editar condiciones (abre P-26.05) | «Sin compras vivas» (revolvente) o «Sin pagos pendientes». Tarjeta bancaria o crédito de tienda sin alguno de los cinco datos de su estado de cuenta (N-283): «Falta capturar · Capturar» (abre P-26.05 en esa línea) en lugar de sus intereses y su pago mínimo (K18) | Esqueleto | Reintentar |
| P-19 | Guardar (alta) o Marcar pago (detalle) | Editar; Eliminar; Ver gasto vinculado | — | Vista previa «Calculando…» | Dato inválido: junto al campo. Motor sin resultado: «No se pudo calcular el calendario» y no guarda |
| P-20 | Usar escenario recomendado | Usar otro óptimo; Cambiar restricciones (solo para esta comparación, K16); Cambiar pesos | «Ningún escenario cumple las restricciones» + la restricción que más descarta. Sin restricciones capturadas: «Faltan restricciones · Capturar» y no hay recomendación (K16) | «Evaluando N escenarios…» | Reintentar |
| P-21 | Registrar peso | Editar perfil (abre P-26.13); Confirmar meta sugerida (K23); Ajuste por objetivo de cada persona, en el lugar, con «Restablecer» (K36, N-327) | «Completa el perfil de salud» | Esqueleto | Reintentar |
| P-22 | Guardar | Cambiar persona | — | Guardar en espera | Junto al campo |
| P-23 | Generar reporte mensual | Estado de cuenta de Karen; Ver reporte web; Filtrar | «Aún no hay reportes (el primero se genera el día 1)» | Esqueleto | Fila en error con motivo y «Reintentar envío» |
| P-24 | Descargar PDF (o Generar PDF si no existe) | Cambiar mes; Compartir | «Sin datos en sep-2026» | Esqueleto por sección | Por sección, con Reintentar |
| P-25 | Generar PDF | Cambiar periodo; Compartir el último PDF | «Karen no tiene compras ni pagos en el periodo» | Esqueleto | Reintentar |
| P-26 | Según sección (2.3) | Según sección | Según sección | Esqueleto de formulario | El campo que no se guardó vuelve a su valor anterior con aviso |

### 2.3 Ajustes: las 14 secciones (decisión G)

En escritorio, Ajustes muestra una columna interna con las 14 secciones y el contenido a la derecha. En móvil, una lista de secciones y cada sección como página. U4 especifica campos, validaciones y vista previa del efecto.

| Id | Sección | Contiene | Vistas | Acción primaria |
|---|---|---|---|---|
| P-26.01 | Cuenta y roles | Usuarios Axel y Ximena (co-titulares, mismos permisos), cambio de la contraseña propia, sesiones de ambos por dispositivo con cierre, «Este dispositivo» (instalar la app, certificado de la CA local), parámetros de acceso: intentos de inicio de sesión, su ventana, días de sesión y largo mínimo de contraseña (K44, N-320) | Tabla (sesiones) | Cerrar otras sesiones |
| P-26.02 | Presupuesto | Monto mensual por categoría, gastos fijos recurrentes, despensa bimestral programada (monto base, día, línea y plazo con que se paga), interruptor de proyección de los meses 7 a 12 «Repetir la despensa bimestral y el papel» (apagado por defecto, K21) | Tabla | Nuevo gasto fijo |
| P-26.03 | Ingresos | Ingreso quincenal esperado de Axel (quincena del 15 y del 30) | Tabla | — (edición en el lugar) |
| P-26.04 | Categorías | Nombre, tipo, reparto por defecto, orden (arrastrar), activa | Tabla*, tablero (por reparto) | Nueva categoría |
| P-26.05 | Créditos y condiciones (incluye promociones) | Líneas por tipo y sus condiciones (tarjeta bancaria y crédito de tienda: tasa anual tal como se capturó, base diaria o mensual, si incluye IVA, porcentaje de pago mínimo y mínimo fijo, tomados de su estado de cuenta; los captura Axel, sin valor por defecto, y mientras falte alguno la línea dice «Falta capturar», K12 y K18), interruptor A/B con vista previa, tabla de tasas por plazo, calendario de promociones, restricciones duras globales del comparador (las captura Axel; sin valor por defecto, K16) y sus pesos, tiempo máximo del comparador (K43 y K44, N-321) | Galería (líneas), tabla (tasas), calendario* (promociones) | Nuevo crédito |
| P-26.06 | Tiendas y precios (incluye mercados y membresía Sam's) | Cadenas y sucursales, mercados y tianguis (días de visita, qué se compra ahí), membresía de Sam's (costo anual), fuentes de precio (manual, OCR, lectura asistida bloqueada hasta veredicto, PROFECO), parámetros de precios: días para precio viejo, días para peso cero y tamaño máximo del corte PROFECO (K6, K22, K40 y K44, N-322) | Tabla*, galería | Nueva tienda o mercado |
| P-26.07 | Recetas | Porciones base, tiempos de comida, bases nutricionales locales (USDA, Open Food Facts, manuales), recetario inicial | Tabla | Recargar base nutricional |
| P-26.08 | Reportes y frecuencias | Resumen diario (noche), semanal (domingo), mensual PDF (día 1), vista web mensual, destinatarios | Tabla | Generar reporte mensual ahora |
| P-26.09 | Notificaciones | Las 4 alertas (activa, destinatarios, anticipación), chats de Axel y de Ximena en el bot MultiX: el bot detecta el chat cuando la persona le escribe y ella lo confirma aquí o lo descarta (K4, N-324), umbrales del aviso de presupuesto (K44, N-323) | Tabla | Enviar aviso de prueba |
| P-26.10 | Apariencia | Tema: sistema (por defecto), claro u oscuro | — | — (un solo control) |
| P-26.11 | Datos y respaldos | Respaldos (lista, integridad, restauración), importar `gastos_fijos_ingresos.csv`, estado de la semilla, IA opcional (interruptores), estado del servicio, parámetros de datos: tamaño máximo de un adjunto y retención de las claves de idempotencia (K43 y K44, N-325) | Tabla | Respaldar ahora |
| P-26.12 | Personas y reparto | Axel, Ximena, Karen (externa, sin cuenta), reparto por defecto del conjunto (50/50), tope de pago por quincena por persona (opcional, K13) | Tabla | — (edición en el lugar) |
| P-26.13 | Salud y metas de ahorro | Perfiles de salud de Axel y Ximena (objetivo de peso y meta de kcal sugerida para confirmar, K23; ajuste por objetivo de cada persona, K36, N-327), ajuste por objetivo de la casa para bajar y para subir (K36 y K44, N-326), metas de ahorro con fecha, aportación planificada al fondo | Tabla, galería (metas) | Nueva meta |
| P-26.14 | Inventario | Modo por producto (stock por unidad o ritmo de consumo), umbral de reposición, regla de sobrantes | Tabla*, tablero (por modo) | — (edición en el lugar) |

Estados de Ajustes: vacío por sección con su acción (p. ej. «Sin créditos · Nuevo crédito»); carga con esqueleto del formulario; error con el campo de vuelta a su valor anterior. El interruptor A/B y la restauración de un respaldo muestran su efecto antes de confirmar.

**Parámetros operativos (K44).** Todo valor numérico operativo existe y se cambia en Ajustes, incluidos los de seguridad. Cada uno vive en la sección a la que pertenece, sin sección nueva: siguen las 14 de la decisión G (ubicación por sección: lectura PROPUESTA del orquestador en K44; Axel la ve en la puerta 2; U4 especifica los campos). Reglas comunes:
- **Valor inicial visible**: junto a cada campo se lee «Valor inicial: …», el decidido en `decisiones.md` (columna «Decisión»). Con los valores iniciales el sistema reproduce las cifras canon del oráculo.
- **Sin valor inventado**: un campo vacío o no válido no se guarda y vuelve a su valor anterior con el aviso (estado de Ajustes de arriba). El sistema nunca pone un valor que nadie decidió: el único que ofrece es el inicial. Los mínimos y máximos de los valores de seguridad los propone S1 y se aprueban en la puerta 2 (K44); hasta entonces no hay topes.
- **«Restablecer al valor inicial»**: una acción por parámetro, activa cuando el valor vigente difiere del inicial.
- **Sin valor decidido**: si un parámetro no tiene valor inicial decidido (p. ej. `send_time` y `advance_days`, cuyas claves fijó K43), el campo dice «Sin valor inicial», «Restablecer» no aparece y nada se rellena con un valor inventado; el valor se propone según K31.
- **Una necesidad por grupo**, no una por parámetro: lee cada valor vigente con su valor inicial y guarda o restablece los del grupo.

| Sección | Parámetro | Valor inicial | Decisión | Necesidad |
|---|---|---|---|---|
| P-26.01 Cuenta y roles | Intentos fallidos de inicio de sesión permitidos en la ventana | 5 | K44 | N-320 |
| P-26.01 Cuenta y roles | Ventana de los intentos (deslizante, por usuario e IP) | 15 min | K11, K44 | N-320 |
| P-26.01 Cuenta y roles | Duración de la sesión (renovable con el uso) | 90 días | K11, K44 | N-320 |
| P-26.01 Cuenta y roles | Largo mínimo de una contraseña nueva | 12 caracteres | K6, K44 | N-320 |
| P-26.05 Créditos y condiciones | Tiempo máximo del comparador (si se pasa, responde «no disponible») | 5 s | K43, K44 | N-321 |
| P-26.06 Tiendas y precios | Días para que un precio sea viejo | 30 días | K22, K44 | N-322 |
| P-26.06 Tiendas y precios | Días en que el peso de un precio viejo llega a cero y deja de competir | 90 días | K40, K44 | N-322 |
| P-26.06 Tiendas y precios | Tamaño máximo de un corte PROFECO importado | 1 GB | K6, K44 | N-322 |
| P-26.08 Reportes y frecuencias | Hora de cada reporte (`send_time`) | Sin valor inicial decidido (K31) | K43, K44 | N-297 (por reporte) |
| P-26.09 Notificaciones | Umbrales del aviso de presupuesto, categoría y total (`alerts.budget_warning_pct` y `alerts.budget_limit_pct`) | 80 % y 100 % | K44 | N-323 |
| P-26.09 Notificaciones | Anticipación del aviso de pago próximo (`alerts.payment_due_days_before`: lista de días antes del pago) | [1, 0]: 1 día antes y el mismo día (decisión F; N-298) | K44 | N-323 |
| P-26.09 Notificaciones | Anticipación del aviso de precios viejos antes del día de despensa (`advance_days`, días enteros ≥ 0) y hora de los avisos que la llevan (`send_time`) | Sin valor inicial decidido (K31) | K43, K44 | N-300 (por aviso) |
| P-26.11 Datos y respaldos | Tamaño máximo de un adjunto de gasto | 10 MB | K43, K44 | N-325 |
| P-26.11 Datos y respaldos | Retención de las claves de idempotencia | 24 h | K6, K44 | N-325 |
| P-26.13 Salud y metas de ahorro | Ajuste por objetivo de la casa: bajar y subir | 0 kcal y 0 kcal | K36, K44 | N-326 |

El ajuste por objetivo **de cada persona** (K36) no es un parámetro de la casa: se edita en P-21 (N-327) y en su ficha de P-26.13; vacío = usa el de la casa. Los demás parámetros que inventaríe S1 en el esquema 1.3 (K44) siguen estas mismas reglas en la sección a la que pertenecen; si alguno no cabe en un grupo de esta tabla, U1 suma su necesidad en la versión siguiente.

### 2.4 Estados comunes

- **Carga**: esqueleto con la forma final del contenido. Si tarda más de 10 s: «Esto está tardando más de lo normal · Reintentar». Nunca cifras provisionales.
- **Vacío**: dice qué falta y ofrece la acción que lo llena; nunca una pantalla en blanco.
- **Error**: lenguaje simple (qué pasó y qué hacer) + Reintentar. Lo capturado no se pierde. La edición en el lugar que falla regresa al valor anterior.
- **Falta capturar** (K18): una tarjeta bancaria o un crédito de tienda al que le falta alguno de los cinco datos de su estado de cuenta (N-283) no muestra intereses ni pago mínimo calculados con valores supuestos; en su lugar dice «Falta capturar · Capturar», que abre P-26.05 en esa línea (RN-26). En P-18 lo detallan N-200, N-201 y N-205.
- **Sin conexión** (G-07, K9): franja fija «Sin conexión con el servidor de casa» y, si hay capturas guardadas en el teléfono, «· n en cola», donde n es el número de capturas pendientes (N-317). Lo ya cargado sigue visible con su hora de carga. **Gastos y peso se siguen capturando** (P-03 pestaña Gasto y P-22, con las categorías de la última carga): «Guardar» deja la captura en la cola del teléfono con su clave de idempotencia y su fecha y hora de captura (C1-19), y el aviso dice «Gasto guardado en el teléfono · Se enviará al volver a la red de casa · Deshacer» (o «Peso guardado…»). Al volver la conexión (N-004), la cola se envía sola, en orden de captura y sin duplicar (N-054, N-240); al terminar: «n capturas enviadas». Si un envío falla, la captura sigue en la cola con «No se pudo enviar · Revisar» y la franja lo cuenta; nunca se pierde ni se borra sin confirmación (N-317). **Nada más entra a la cola** (K9). Las demás pantallas sin conexión: **solo lectura** con lo ya cargado en P-01, P-02, P-04, P-06 a P-18, P-21 y P-23 a P-26, y sus escrituras (acción primaria que guarda, edición en el lugar, marcar comprado, generar) muestran «Requiere conexión con el servidor de casa»; **aviso «Requiere conexión»** en lugar del contenido en P-00, P-05, P-19, P-20 y en las pestañas Ingreso y Ticket de P-03; una página que nunca se cargó en el teléfono muestra «Sin conexión · Esta página aún no se ha cargado en este teléfono». Registrar un gasto sin conexión sigue siendo de 3 toques (F-1).
- **Sesión vencida** o cerrada desde otro dispositivo: P-00 con aviso; al entrar, vuelve a la página donde estaba.
- **Registro borrado** por el otro usuario: «Ya no existe» + Volver. Si un registro cambió mientras estaba abierto, aviso «Este registro cambió · Recargar» (S6 define la detección).

### 2.5 Vistas de base de datos por entidad

Las 31 entidades de §3.4.1 (lista de S1-01). La vista elegida, sus filtros, orden, agrupación y columnas visibles se recuerdan por dispositivo (Local); no hay vistas guardadas compartidas en esta versión. Las entidades sin vista de lista propia dicen por qué.

**2.5.a Vistas disponibles** (Sí* = vista por defecto)

| Entidad | Dónde vive | Tabla | Tablero (agrupa por) | Calendario (fecha) | Galería (tarjeta) | Sin vista de lista: motivo |
|---|---|---|---|---|---|---|
| persona | P-26.12 | Sí* | — | — | — | |
| cuenta_usuario | P-26.01 | Sí* | — | — | — | |
| sesion_dispositivo | P-26.01 | Sí* | — | — | — | |
| periodo | — | — | — | — | — | Dato derivado (mes, quincena, semana); aparece como selector de periodo en P-02, P-06, P-07, P-08 y P-10 |
| ingreso | P-08 | Sí* | — | Fecha | — | |
| categoria | P-26.04 | Sí* | Reparto por defecto | — | — | |
| gasto | P-02 | Sí* | Categoría | Fecha del gasto | Tickets: foto, comercio, total | |
| regla_reparto | — | — | — | — | — | Se edita como propiedad de la categoría (P-26.04) y del gasto (P-03, P-04); el % por defecto vive en P-26.12 |
| meta_ahorro | P-09, P-26.13 | Sí (P-26.13) | Estado (en ruta, atrasada, cumplida) | Fecha meta | Sí* (P-09): avance | |
| movimiento_ahorro | P-09 | Sí* | — | Fecha | — | |
| credito | P-17, P-26.05 | Sí (P-26.05) | Tipo | — | Sí* (P-17): deuda, disponible, próximo pago | |
| tasa_por_plazo | P-26.05, P-18 | Sí* | — | — | — | |
| compra_credito | P-18 | Sí* | Estado (viva, liquidada) | Fecha de compra | — | |
| pago_programado | P-17, P-18, P-19 | Sí | Estado (pendiente, pagado) | Sí*: fecha de pago (subsección Calendario) | — | |
| promocion | P-26.05; franjas en P-17 | Sí | Línea | Sí*: vigencia de inicio a fin | — | |
| tienda | P-26.06 | Sí* | Tipo (cadena, mercado) | Días de visita (recurrentes) | Sí | |
| producto | P-15 | Sí* | Frescura (fresco, viejo, sin precio) | — | Sí | |
| presentacion | — | — | — | — | — | Tabla dentro de la página del producto (P-15) |
| precio | P-15 (precios del producto) | Sí* | Frescura | — | — | |
| receta | P-11 | Sí | Tiempo de comida | — | Sí* | |
| ingrediente_receta | — | — | — | — | — | Tabla dentro de la receta (P-12), en el orden de la receta |
| tiempo_comida | — | — | — | — | — | Catálogo fijo de 4 (desayuno, comida, cena, colación; decisión D): son las filas del plan y se muestran en P-26.07 |
| plan_comida | P-10 | Sí (asignaciones) | Tiempo de comida | Sí*: semana | Recetas de la semana | |
| inventario_item | P-16, P-26.14 | Sí* | Estado (suficiente, reponer, agotado) | — | Sí | |
| perfil_salud | — | — | — | — | — | Una ficha por persona (2) en P-21 y P-26.13 |
| peso_historico | P-21 | Sí* | — | — | — | La serie va en Tendencias |
| nutricion_ingrediente | — | — | — | — | — | Propiedad del producto (P-15) y columna de los ingredientes (P-12); las cargas se ven en P-26.07 |
| notificacion | P-23 (avisos enviados), P-26.09 (configuración) | Sí* | Tipo | Fecha de envío | — | |
| destinatario | — | — | — | — | — | Se elige dentro de cada aviso y reporte (P-26.08, P-26.09) |
| job | P-23, P-05, P-26.06, P-26.11 | Sí* | Tipo | Fecha | — | |
| respaldo | P-26.11 | Sí* | — | — | — | |

**2.5.b Filtros y orden por defecto** (los filtros usan el vocabulario de C1: persona, familia, categoría, tienda, periodo, crédito; se agregan estado y texto donde se indica)

| Entidad | Filtros disponibles | Filtro por defecto | Orden por defecto |
|---|---|---|---|
| persona | — | Ninguno (Axel, Ximena, Karen) | Axel, Ximena, Karen |
| cuenta_usuario | — | Ninguno (2 cuentas) | Nombre |
| sesion_dispositivo | Persona | Sesiones abiertas de ambos usuarios | Último uso descendente |
| ingreso | Periodo | Mes en curso | Fecha ascendente |
| categoria | Estado (activa, archivada), texto | Activas | Orden manual (arrastrar) |
| gasto | Persona, familia, categoría, tienda, periodo, crédito (forma de pago), tipo, con ticket, texto | Mes en curso | Fecha descendente, luego captura descendente |
| meta_ahorro | Estado | Activas | Fecha meta ascendente |
| movimiento_ahorro | Periodo, destino | Mes en curso | Fecha descendente |
| credito | Tipo, estado | Activas | Fecha del próximo pago ascendente |
| tasa_por_plazo | Crédito | Línea abierta | Plazo ascendente |
| compra_credito | Crédito, familia, periodo, estado | Vivas | Fecha descendente |
| pago_programado | Crédito, familia, periodo, estado | Con alguna parte pendiente, de hoy en adelante | Fecha ascendente |
| promocion | Crédito, periodo, tienda (comercio) | Vigentes y futuras | Inicio ascendente |
| tienda | Tipo, estado | Activas | Nombre |
| producto | Tienda, categoría, frescura, fuente, texto | Ninguno | Nombre |
| precio | Tienda, periodo, fuente | Precios del producto abierto | Fecha descendente |
| receta | Tiempo, origen, texto, rango de costo y kcal | Ninguno | Nombre |
| plan_comida | Persona, periodo (semana) | Semana por defecto (§0.2) | Día y tiempo |
| inventario_item | Categoría, modo, estado, texto | Ninguno | Agotado, reponer, suficiente; luego nombre |
| peso_historico | Persona, periodo | Ambas personas | Fecha descendente |
| notificacion | Periodo, tipo | Mes en curso | Fecha descendente |
| job | Periodo, tipo, estado | Según la pantalla (P-23: reportes; P-05: lecturas sin confirmar) | Fecha descendente |
| respaldo | Periodo | Ninguno | Fecha descendente |

### 2.6 Edición en el lugar por pantalla

Como en Notion: se toca el valor, se cambia y se guarda solo. Lo calculado nunca se edita en el lugar.

| Pantalla | Se edita en el lugar | No se edita en el lugar (y por qué) |
|---|---|---|
| P-02 | En tabla y tablero: descripción, categoría, monto, fecha, persona, reparto y tipo; arrastrar una tarjeta a otra columna cambia su categoría (N-047) | Forma de pago a crédito (abre P-19 porque genera pagos); ticket adjunto (en P-04) |
| P-04 | Todas las propiedades, los renglones (producto, cantidad, precio, cuenta) y la nota (N-062) | Compra a crédito vinculada (en P-19); partes, devengo y efectivo (calculados) |
| P-05 | Encabezado y renglones propuestos, antes de confirmar | Después de confirmar se edita en P-04 |
| P-06 | Plan mensual de cada categoría, tras «Ajustar presupuesto» (N-087) | Real, restante y % (calculados) |
| P-08 | Monto, fecha y nota de quincenas capturadas y de extras (N-105) | Monto esperado (en P-26.03) |
| P-09 | Fecha, monto, destino y nota de cada movimiento (N-116) | Ahorrado, avance y estado (calculados); definición de metas (en P-26.13) |
| P-10 | Casillas: receta, porciones y «fuera» (N-123 a N-125) | Costo y kcal (calculados) |
| P-11 | En tabla: nombre y tiempos aptos (N-142) | Ingredientes y pasos (en P-12); costo y kcal (calculados) |
| P-12 | Propiedades, ingredientes (producto, cantidad, unidad, presentación), pasos y notas (N-142) | Costo, kcal y fuentes (calculados) |
| P-13 | Comprado, cantidad y tienda (N-151, N-153) | Artículos que vienen del plan (se cambian en P-10) |
| P-15 | Nombre, marca, categoría, presentaciones y sustitutos (N-175); precios capturados a mano o por ticket (N-176) | Precios de PROFECO (vienen del corte importado) |
| P-16 | Existencia, con motivo (N-183) | Estado y días restantes (calculados) |
| P-17, P-18 | — | Todo se calcula o viene de P-19 y P-26.05; «Marcar pago» es una acción (N-210) |
| P-19 | Comercio y nota; capital, intereses, plazo y reparto recalculan con vista previa (N-216) antes de guardar (N-219) | Pagos (se marcan con N-210) |
| P-20 | Datos de entrada: monto, fecha, holgura, restricciones, tope y pesos | Resultados (calculados) |
| P-21 | Registros del historial de peso (N-242); ajuste por objetivo de cada persona, para bajar y para subir, en kcal (N-327, K36) | Perfil (en P-26.13); IMC y kcal (calculados) |
| P-23, P-24, P-25 | — | Reportes generados (son copias con fecha) |
| P-26.01 a P-26.14 | Todos los campos de configuración, con guardado automático; los parámetros operativos muestran su valor inicial y «Restablecer al valor inicial» (K44, §2.3) | El interruptor A/B y la restauración de un respaldo piden confirmar después de mostrar su efecto |

---

## 3. Necesidades de datos por pantalla (c)

**Cómo leer**: cada fila es una necesidad que C1 debe cubrir con uno o más endpoints. **Tipo L** = lectura (métrica o lista); **Tipo E** = escritura (crear, cambiar, borrar o lanzar un trabajo). La definición dice qué se calcula (en palabras), sobre qué periodo, con qué filtros, cómo se agrupa y en qué orden. Cuando una cifra la calcula un motor, se dice cuál; la regla exacta vive en `reglas-negocio.md` y la pantalla solo la muestra. Una necesidad se define una vez y otras pantallas la **reutilizan** por su id.

### 3.1 Globales

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-001 | L | Sesión actual | Usuario de la sesión (Axel o Ximena), nombre visible y rol (co-titular). Se usa para valores por defecto y «creado por» | Inmediata |
| N-002 | L | Búsqueda global | Coincidencias de un texto, sin distinguir mayúsculas ni acentos, en: gastos (descripción, comercio, categoría), ingresos extra (concepto), recetas (nombre, ingredientes), productos (nombre, marca, presentación), tiendas y mercados, líneas de crédito, compras a crédito (comercio), metas de ahorro, secciones de Ajustes y nombres de página. Agrupado por tipo; dentro de cada grupo, por relevancia y luego por fecha descendente; máximo 5 por grupo con «Ver todos». Texto vacío: recientes (N-003) | Inmediata |
| N-003 | L | Recientes | Últimas 8 páginas o registros abiertos en este dispositivo | Local |
| N-004 | L | Estado del servidor | Si el servidor responde y la hora de la última respuesta. Alimenta la franja «Sin conexión» (G-07); mientras hay error se vuelve a consultar cada 30 s; al volver la respuesta se envía la cola (N-317) | Inmediata |
| N-005 | L | Líneas para la barra lateral | Nombre e id de cada línea de crédito activa, en orden alfabético (hijos dinámicos de Créditos) | Inmediata |
| N-006 | E | Cerrar sesión | Cierra la sesión de este dispositivo; si hay capturas en cola (N-317), antes avisa cuántas se quedarían sin enviar | — |
| N-317 | L | Cola sin conexión | Capturas de gasto (P-03) y de peso (P-22) guardadas en el teléfono mientras no hay conexión (K9; solo esos dos tipos): tipo, datos capturados, fecha y hora de captura, clave de idempotencia y estado (en cola, enviando, no se pudo enviar con su motivo). Alimenta el conteo «n en cola» de G-07. Al volver la conexión (N-004) se envían en orden de captura con N-054 o N-240, que reciben la clave y la fecha de captura; la enviada sale de la cola; la que falla se queda con «Revisar» (abre la captura para corregirla y reenviarla) o «Descartar» (con confirmación). Si la sesión venció, se envían después de iniciar sesión | Local |

### 3.2 P-00 Iniciar sesión

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-010 | E | Iniciar sesión | Usuario y contraseña. Crea una sesión larga en este dispositivo (meses, decisión A). Responde: sesión creada; error genérico «usuario o contraseña incorrectos»; o bloqueo con hora de desbloqueo tras 5 intentos fallidos en 15 min (S6 fija si el conteo es por usuario o por dispositivo) | — |
| N-011 | L | Estado de bloqueo | Para quien intenta entrar: bloqueado sí o no y hora local de desbloqueo | Inmediata |

### 3.3 P-01 Inicio

**Resumen (6 tarjetas, orden fijo).** Cada tarjeta lleva a su detalle: Mes vs plan → P-06; Flujo quincenal → P-07; Próximo pago → P-18 de esa línea; Comidas → P-10; Mejor canasta → P-14; Ahorro → P-09.

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-020 | L | Tarjeta «Mes vs plan» | Mes en curso. **Real** = suma de la parte del hogar de todo lo devengado en el mes (motor S2). **Plan** = suma del presupuesto mensual de las categorías activas (P-26.02). **% usado** = real ÷ plan × 100, entero. **Restante** = plan − real (negativo si se rebasó). **Estado**: normal si % < 80; aviso si ≥ 80; excedido si ≥ 100 (mismos umbrales que la alerta, decisión F). **Días que faltan** para terminar el mes, contando hoy. Sin plan: «—» | Diaria |
| N-021 | L | Tarjeta «Flujo quincenal» | Quincena en curso, en efectivo y parte del hogar (motor S2). **Ingresos** = quincena de Axel (capturada; si no, la esperada de P-26.03 con la marca «esperada») + extras con fecha en la quincena. **Egresos hechos** = salidas de efectivo con fecha en la quincena hasta hoy (gastos de contado + partes del hogar ya marcadas como pagadas). **Comprometido** = partes del hogar de pagos pendientes con fecha en la quincena + gastos venideros y fijos de la quincena aún no pagados. **Saldo al cierre** = ingresos − egresos hechos − comprometido | Diaria |
| N-022 | L | Tarjeta «Próximo pago y disponible» | **(a) Próximo pago**: entre todas las líneas, el pago pendiente con la fecha más cercana igual o posterior a hoy. Si existe algún pago pendiente con fecha anterior a hoy, se muestra primero el más antiguo con la marca «sin marcar como pagado». Campos: fecha con día de la semana, línea, monto total, parte del hogar, parte de Karen y días que faltan; si varios pagos caen el mismo día, se suman y se indica cuántos son. **(b) Disponible**: para cada línea revolvente con interruptor A/B (hoy, Kueski), el disponible de hoy según la lectura activa (motor S3), con su etiqueta («Lectura B») y el aviso «pendiente confirmar» | Diaria |
| N-023 | L | Tarjeta «Comidas de la semana» | Semana en curso (costo y kcal: motor S4). **Costo** = suma del costo de las casillas con receta; las marcadas «fuera» no suman, porque su gasto se registra como gasto normal. **kcal promedio por día** de cada persona = suma de kcal asignadas a esa persona ÷ número de días de la semana con al menos un tiempo asignado a esa persona (se muestra «n de 7 días»), contra su meta diaria (P-26.13). **Casillas sin armar** = número de pares día × tiempo obligatorio (desayuno, comida, cena) en que alguna de las 2 personas no tiene receta ni «fuera»; máximo 21; las colaciones son opcionales y no cuentan. Definición **FIRME** (K43), la misma del aviso «plan sin armar»; es una regla, no un parámetro | Inmediata |
| N-024 | L | Tarjeta «Mejor canasta» | Lista de compras vigente (§0.2, N-150) comparada por el motor S5: tienda con menor total (Sam's Club se compara como si tuvieran membresía y lleva esa marca), diferencia contra la segunda tienda, cobertura (artículos con precio en esa tienda ÷ artículos de la lista) y número de precios viejos usados. Si la lista está vacía, usa la canasta de la última compra registrada con renglones e indica cuál | Diaria |
| N-025 | L | Tarjeta «Ahorro» | **Ahorrado** = suma de aportaciones − retiros de fondo y metas. **Aportado este mes** contra **planificado del mes** (suma de aportaciones mensuales planificadas del fondo y de las metas activas). **Meta más próxima** por fecha: nombre, % de avance (ahorrado ÷ monto meta, entero), fecha meta y estado (en ruta, atrasada o cumplida; motor S2) | Diaria |

**Detalle.**

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-026 | L | Participación de Karen | Solo despensa común y línea revolvente donde ella compra. (1) Su parte devengada de la despensa en el mes en curso (motor S2). (2) Sus compras particulares a crédito vivas: comercio, fecha, total, pagos hechos de cuántos, saldo. (3) Sus partes de pagos con fecha en el mes en curso: fecha, monto, estado (pagada o pendiente). (4) Su deuda viva total en la línea = suma de sus partes pendientes. Nota fija: «Karen paga directo a Kueski su parte; su parte es informativa y repone la línea cuando paga» | Diaria |
| N-027 | L | Segundo nivel de las tarjetas | Reutiliza: N-081 (las 5 categorías con mayor % usado), N-090 (filas de la quincena en curso), N-192 (próximos 5 pagos, todas las líneas), N-121 (comidas de hoy por tiempo y persona), N-162 (los 5 artículos con mayor diferencia de precio), N-111 (metas) | Según la necesidad reutilizada |

**Tendencias y Calendario.**

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-028 | L | Gráficas de Inicio | Reutiliza cuatro series: N-085 (real contra plan por mes), N-091 (devengado contra efectivo, 12 meses), N-194 (deuda viva proyectada por mes), N-129 (costo semanal de comidas, 12 semanas) | Inmediata |
| N-029 | L | Calendario del hogar | Eventos del rango visible (el mes; en móvil, agenda de hoy a 30 días). Por día: ingresos programados (15 y 30), pagos de crédito (monto total, partes y estado), gastos fijos y venideros, fechas de corte de tarjetas, despensa programada, días de visita a mercados, vigencia de promociones (franja de inicio a fin), fechas meta de ahorro y envío del reporte mensual (día 1). Filtro por tipo de evento. Orden: fecha y, dentro del día, el orden de tipos de esta lista | Diaria |

### 3.4 P-02 Todos los gastos

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-040 | L | Lista de gastos | Gastos con fecha en el periodo (por defecto el mes en curso; también quincena, semana o rango). Filtros combinables: persona (Axel, Ximena), familia (Hogar, Karen o Conjunto, incluyentes según §0.2: Hogar y Karen incluyen lo conjunto, K5), categoría (varias), tienda o comercio, forma de pago (de contado o una línea de crédito), tipo (fijo, variable, venidero), con ticket (sí o no) y texto. Columnas: fecha, descripción o comercio, categoría, monto, parte del hogar, parte de Karen, persona, forma de pago, tipo, ticket. Orden por defecto: fecha descendente y, a igual fecha, captura descendente; cualquier columna puede ordenar. Páginas de 50 | Inmediata |
| N-041 | L | Totales del filtro | Sobre exactamente el conjunto de N-040 (todas sus páginas): número de gastos, suma de montos, suma de la parte del hogar y suma de la parte de Karen | Inmediata |
| N-042 | L | Resumen del periodo | Con el periodo y los filtros de N-040: suma de la parte del hogar; la misma suma del periodo anterior equivalente (mes, quincena o semana anterior) y la diferencia; promedios diario, semanal y mensual de gasto (motor S2); categoría con mayor gasto | Inmediata |
| N-043 | L | Gasto por categoría | Con el periodo y los filtros: por categoría, suma de la parte del hogar, número de gastos y % del total; orden descendente por monto. Alimenta el tablero (columnas por categoría con su total en la cabecera) | Inmediata |
| N-044 | L | Tendencia mensual | Por mes calendario, del primer mes con datos al mes en curso (máximo 12): suma de la parte del hogar devengada por categoría; las 6 categorías con más gasto del rango y «Otras». Respeta los filtros persona, familia y forma de pago | Inmediata |
| N-045 | L | Calendario de gastos | Mes visible: por día, suma de la parte del hogar y lista de gastos (descripción, monto); los venideros llevan su marca | Inmediata |
| N-046 | L | Galería de tickets | Gastos del periodo con ticket adjunto: miniatura, comercio, fecha y total; orden por fecha descendente | Inmediata |
| N-047 | E | Editar en el lugar | Cambia una propiedad desde la tabla o el tablero (categoría, monto, fecha, persona, reparto, tipo). Arrastrar una tarjeta a otra columna del tablero cambia su categoría. Usa N-062 | — |

### 3.5 P-03 Registrar gasto

Las pestañas Ingreso y Peso de la misma hoja usan N-102 y N-240; la pestaña Ticket abre P-05.

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-050 | L | Categorías sugeridas | Las 6 categorías activas con más gastos registrados (conteo, de ambos usuarios) en los últimos 90 días contando hoy; los empates se resuelven por el orden de P-26.04. Si hay menos de 6 con historial, se completan con las primeras categorías activas en el orden de P-26.04 | Inmediata |
| N-051 | L | Todas las categorías | Categorías activas con su reparto por defecto, en el orden de P-26.04 (para «Más» y para P-26.04) | Inmediata |
| N-052 | L | Valores por defecto | Fecha = hoy. Persona = usuario de la sesión. Forma de pago = «De contado». Reparto = el de la categoría elegida; si la categoría no define uno, «Hogar» (Karen solo participa en la despensa común y en sus compras a crédito, decisión A). Tipo = variable, o venidero si la fecha es posterior a hoy | Inmediata |
| N-053 | L | Líneas para «forma de pago» | Líneas de crédito activas (nombre y tipo), en orden alfabético. Elegir una revolvente o de tienda lleva a P-19 prellenada con monto, fecha y comercio | Inmediata |
| N-054 | E | Crear gasto | Obligatorios: monto (mayor que 0, en centavos), categoría, fecha, persona y forma de pago. Opcionales: descripción o comercio, tienda, reparto distinto al de la categoría (Hogar, Karen o Conjunto con % del hogar), adjunto y nota. Devuelve el gasto con sus partes calculadas. Acepta la clave de idempotencia y la fecha y hora de captura de una captura hecha sin conexión (N-317, K9): reenviar la misma clave no crea otro gasto | — |
| N-055 | E | Deshacer | Durante 5 s después de guardar, «Deshacer» elimina el gasto recién creado (usa N-063); sin conexión, lo quita de la cola del teléfono (N-317) | — |

### 3.6 P-04 Gasto

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-060 | L | Gasto | Todas sus propiedades: fecha, monto, categoría, persona, descripción o comercio, tienda, forma de pago (con la línea y la compra a crédito vinculadas), tipo, reparto y partes calculadas (hogar, Karen), devengo (mes o meses a los que se asigna y monto en cada uno, motor S2), efectivo (fecha o fechas en que sale el dinero: la del gasto si es de contado, los pagos si es a crédito), adjunto, nota, quién lo creó y cuándo, quién lo cambió por última vez y cuándo | Inmediata |
| N-061 | L | Renglones del ticket | Si el gasto viene de un ticket, por renglón: producto, presentación, cantidad, precio unitario, importe y cuenta (Hogar, Karen, Conjunto con %, o «mixto» según S1). Totales por cuenta y partes del hogar y de Karen. Si la suma de importes difiere del total del ticket, la diferencia aparece como renglón «Ajuste del ticket» (descuentos, redondeo). Orden: el del ticket | Inmediata |
| N-062 | E | Actualizar gasto | Cambia una o varias propiedades o renglones (edición en el lugar con guardado automático) y recalcula las partes | — |
| N-063 | E | Eliminar gasto | Con confirmación. Un gasto con compra a crédito vinculada no se elimina desde aquí: primero se elimina la compra en P-19 | — |
| N-064 | E | Adjuntar o quitar ticket | Imagen (JPG, PNG o WebP) o PDF de hasta 10 MB (**FIRME**, K43); el máximo se cambia en Ajustes › Datos y respaldos (K44, N-325). Un archivo más grande no se adjunta y el aviso dice el máximo vigente | — |
| N-065 | E | Marcar venidero como pagado | Para un gasto venidero: fecha real (hoy por defecto) y monto real (el previsto por defecto); pasa a gasto realizado | — |
| N-066 | E | Duplicar gasto | Crea una copia con fecha de hoy | — |

### 3.7 P-05 Capturar ticket

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-070 | E | Leer ticket | Sube la foto (cámara o galería) y lanza el OCR local (Tesseract, español). Devuelve el identificador de la lectura | Al generar |
| N-071 | L | Resultado de la lectura | Encabezado propuesto (tienda y sucursal, fecha, total) y renglones propuestos en el orden del ticket: texto leído, producto y presentación sugeridos del catálogo, cantidad, precio unitario, importe, cuenta sugerida (la que tuvo ese producto en la última compra confirmada; si no hay, el reparto de la categoría del gasto) y confianza por renglón (alta, media, baja). La IA opcional (P-26.11) solo cambia cómo se propone; nada se guarda sin confirmación humana (decisión E). Si la lectura tiene borrador (N-075), devuelve el borrador con las correcciones guardadas en lugar de la propuesta original (K3) | Al generar |
| N-072 | L | Buscar en catálogo | Productos y presentaciones que coinciden con un texto, para corregir un renglón (N-002 restringida a productos) | Inmediata |
| N-073 | L | Prellenado desde la lista | Si se llega desde P-13 «Terminar compra»: los artículos marcados como comprados (producto, cantidad, tienda) en la lista de la semana que mostraba P-13 (semana explícita, §0.2) como renglones propuestos | Inmediata |
| N-074 | E | Confirmar ticket | En una sola operación: crea el gasto con sus renglones y la cuenta de cada renglón; registra cada precio (fuente «ticket», fecha del ticket, tienda); suma las entradas de inventario de los productos con stock; y, si la forma de pago es crédito, crea la compra a crédito con sus intereses cotizados y plazo (N-218) | — |
| N-075 | E | Guardar para después | Guarda la lectura sin confirmar **con las correcciones hechas hasta ese momento como borrador** (encabezado y renglones; K3); no crea gasto, precio ni inventario. Volver a guardar reemplaza el borrador. Al abrirla desde N-076, P-05 muestra el borrador (N-071) | — |
| N-076 | L | Tickets por confirmar | Lecturas guardadas sin confirmar: fecha de la foto, tienda y total (los del borrador si lo hay; si no, los propuestos) y si tiene borrador (K3); orden por fecha descendente | Inmediata |

### 3.8 P-06 Presupuesto del mes

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-080 | L | Resumen del mes | Mes elegido (por defecto el en curso). **Ingreso del mes** = quincenas capturadas o esperadas + extras. **Plan de gasto** = suma del presupuesto de las categorías. **Ahorro planificado** = suma de aportaciones mensuales planificadas. **Sin asignar** = ingreso − plan de gasto − ahorro planificado. **Real** = suma de la parte del hogar devengada (motor S2). **Restante** = plan − real. % usado y estado con los umbrales de N-020 | Inmediata |
| N-081 | L | Presupuesto por categoría | Por categoría activa: plan del mes, real (parte del hogar devengada), diferencia (plan − real), % usado, estado y número de gastos. Orden por defecto: % usado descendente; alternativos: nombre, plan, real. Agrupación opcional por tipo (fijo o variable). El tablero agrupa en normal, aviso (≥ 80 %) y excedido (≥ 100 %) | Inmediata |
| N-082 | L | Quincenas del mes | Para Q1 y Q2 del mes elegido: real (parte del hogar, por fecha del gasto) y la parte del plan que corresponde a cada quincena según el motor S2 (los gastos fijos cuentan en la quincena de su fecha) | Inmediata |
| N-083 | L | Despensa devengada | Despensas bimestrales que devengan en el mes (motor S2, mitad por mes): fecha de compra, total con intereses, monto devengado en el mes, parte del hogar y parte de Karen | Inmediata |
| N-084 | L | Alertas activas | Categorías y total con % usado ≥ 80 o ≥ 100 en el mes en curso, con la fecha en que cruzaron cada umbral | Inmediata |
| N-085 | L | Plan contra real por mes | Por mes calendario, del primer mes con datos al mes en curso (máximo 12): plan total y real total; al elegir una categoría, sus valores | Inmediata |
| N-086 | L | Calendario del mes | N-029 limitado a eventos de dinero (ingresos, pagos, gastos fijos y venideros) | Diaria |
| N-087 | E | Ajustar presupuesto | Cambia el monto mensual de una categoría con su **alcance** (K2): «Solo este mes» o «De este mes en adelante», con «De este mes en adelante» preseleccionado. «Este mes» es el mes que muestra P-06; la pantalla manda la categoría, el mes, el monto y el alcance. Es el mismo dato que P-26.02 (una sola fuente) | — |

### 3.9 P-07 Flujo y proyección

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-090 | L | Flujo por quincena | Motor S2, efectivo, parte del hogar. Desde la Q1 del mes en curso y 24 quincenas hacia adelante, cada una con su saldo de la quincena y el saldo acumulado. **Meses 1 a 6** (el mes en curso es el 1, como en N-091): ingresos (capturados o esperados), pagos de crédito (partes del hogar), gastos fijos y venideros, gasto variable estimado (método del motor S2, visible en N-093) y aportación planificada a ahorro; las quincenas ya transcurridas del mes en curso muestran lo real. **Meses 7 a 12**: rige N-318 (K21, RN-28), igual que en N-091: con el interruptor apagado (por defecto), solo compromisos contratados (pagos programados de compras ya hechas, partes del hogar), gastos fijos e ingresos del CSV, así que esas quincenas no llevan gasto variable estimado, aportación planificada a ahorro ni gastos venideros; encendido, se agrega la cadencia bimestral de despensa a crédito y el papel alterno. La pantalla dice en esas quincenas qué incluye la proyección (N-093) | Diaria |
| N-091 | L | Devengado contra efectivo, 12 meses | Motor S2. Por mes calendario, del mes en curso a 11 meses después: devengado total, efectivo total, sus partes del hogar y de Karen, y **diferencia acumulada** = devengado acumulado − efectivo acumulado (la misma cifra que `periodos[i].diferencia_acumulada` de `model.json`). No se rotula «deuda viva»: la deuda viva total es N-092 y se relaciona con esta diferencia así (K19): deuda viva total = diferencia acumulada + saldo no devengado. **Meses 7 a 12** (K21): por defecto solo compromisos contratados, gastos fijos e ingresos del CSV; con el interruptor de N-318 encendido se agrega la cadencia bimestral de despensa a crédito y el papel alterno. La pantalla dice en esos meses qué incluye la proyección (N-093) | Inmediata |
| N-092 | L | Deuda viva a una fecha | Motor S3. Para una fecha elegida (por defecto, el fin del horizonte proyectado): **deuda viva total** = lo que queda por pagar después de esa fecha (capital + intereses) de las compras a crédito hechas hasta esa fecha; total, por línea y por compra. Canon semilla: $6,293.44 al 10-abr-2027 (`model.json`; incluye la despensa 4, comprada ese día) | Inmediata |
| N-093 | L | Supuestos de la proyección | Ingreso quincenal esperado (P-26.03), gastos fijos considerados, método y monto del gasto variable estimado (S2), despensas futuras programadas (monto base, fechas, línea y plazo; P-26.02) y el estado del interruptor de los meses 7 a 12 (N-318, K21) | Inmediata |

### 3.10 P-08 Ingresos

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-100 | L | Ingresos del mes | Mes elegido. Quincena del 15 y quincena del 30 (o último día del mes) con estado «esperada» (sin capturar; monto de P-26.03) o «capturada» (monto, quién la capturó y cuándo); ingresos extra (fecha, concepto, monto). Orden por fecha. Totales: capturado, esperado pendiente, extras y total | Diaria |
| N-101 | E | Capturar quincena | Fecha de la quincena (15 o 30 / último día), monto (prellenado con el esperado) y nota. Solo ingresos de Axel (decisión B): no se pide persona | — |
| N-102 | E | Registrar ingreso extra | Fecha (hoy por defecto), concepto, monto y nota | — |
| N-105 | E | Corregir o eliminar ingreso | Cambia monto, fecha, concepto o nota de una quincena capturada o de un extra (edición en el lugar), o lo elimina con confirmación; una quincena eliminada vuelve a «esperada» | — |
| N-103 | L | Ingreso por mes | Por mes, del primer mes con datos al mes en curso (máximo 12): quincenas + extras | Inmediata |
| N-104 | L | Calendario de ingresos | Fechas de quincena (esperadas y capturadas) y extras del mes visible | Diaria |

### 3.11 P-09 Ahorro y metas

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-110 | L | Fondo y totales | Ahorrado total y por destino (fondo general y cada meta) = aportaciones − retiros; aportado en el mes en curso contra planificado del mes | Inmediata |
| N-111 | L | Metas | Por meta activa: nombre, monto meta, fecha meta, ahorrado, % de avance (entero), aportación mensual planificada, aportación mensual necesaria para llegar a la fecha (motor S2) y estado (en ruta, atrasada o cumplida; motor S2). Orden por fecha meta ascendente | Diaria |
| N-112 | L | Movimientos de ahorro | Fecha, tipo (aportación o retiro), monto, destino y nota. Filtros: destino y periodo. Orden por fecha descendente | Inmediata |
| N-113 | E | Registrar movimiento | Fecha (hoy por defecto), tipo, monto, destino (fondo o una meta) y nota | — |
| N-116 | E | Corregir o eliminar movimiento | Cambia fecha, monto, destino o nota (edición en el lugar), o lo elimina con confirmación | — |
| N-114 | L | Tendencia de ahorro | Por mes: ahorrado al cierre y aportado contra planificado | Inmediata |
| N-115 | L | Calendario de ahorro | Fechas meta y movimientos del mes visible | Inmediata |

### 3.12 P-10 Plan de la semana

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-120 | L | Plan de la semana | Semana elegida (al abrir, la «semana por defecto» de §0.2). Por día (lunes a domingo, K8) × tiempo (desayuno, comida, cena, colación) × persona (Axel, Ximena): receta y porciones, «fuera» o vacío. Por casilla: costo (motor S4, con la base de precio usada: tienda y fecha) y kcal por persona (motor S4) | Inmediata |
| N-121 | L | Totales del plan | Por día y por semana: costo (suma de casillas con receta; «fuera» no tiene costo de receta); kcal por día de cada persona y diferencia contra su meta; casillas sin armar (definición de N-023); porciones sobrantes previstas por receta (motor S4; entran al inventario como sobrante). «Comidas de hoy» es el día de hoy de esta necesidad | Inmediata |
| N-122 | L | Recetas para una casilla | Recetas aptas para el tiempo elegido. Orden: veces usadas en las últimas 8 semanas (descendente) y luego nombre. Con costo y kcal por porción; filtro por texto | Inmediata |
| N-123 | E | Asignar receta | Día o días (la fila de días del selector aplica la receta a otros días del mismo tiempo), tiempo, persona (Ambos por defecto, o solo Axel o solo Ximena), receta y porciones por persona (1 por defecto; con Ambos son 2, la porción base de la decisión D) | — |
| N-124 | E | Marcar «fuera» | Día, tiempo y persona (Ambos por defecto). Sin receta ni kcal (decisión B) | — |
| N-125 | E | Vaciar casilla | Quita la receta o la marca «fuera» | — |
| N-126 | E | Copiar semana anterior | Copia las casillas de la semana previa solo en las casillas vacías de la semana elegida | — |
| N-127 | E | Generar lista de compras | Produce o actualiza la lista de compras de la semana que muestra P-10 (N-150); la semana va explícita (§0.2, R3-01) | — |
| N-128 | L | Calendario mensual | Por día del mes visible: plan completo, incompleto o vacío | Inmediata |
| N-129 | L | Tendencias de comidas | Últimas 12 semanas con plan: costo semanal y kcal promedio por día de cada persona contra su meta | Inmediata |

### 3.13 P-11 Recetario

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-135 | L | Recetas | Nombre, tiempos aptos, porciones base, tiempo total (preparación + cocción), costo por porción (motor S4, con su base de precio), kcal por porción (motor S4), origen (recetario inicial o de Axel y Ximena) y foto si existe. Filtros: tiempo, origen, texto (nombre o ingrediente), rango de costo y de kcal. Orden por defecto: nombre; alternativos: costo, kcal, más usadas en 8 semanas | Inmediata |
| N-136 | E | Crear receta | Nombre, tiempos aptos, porciones base (2 por defecto), tiempos de preparación y cocción, ingredientes (producto, cantidad, unidad, presentación comprable), pasos (bloques de notas) y origen = Axel y Ximena | — |
| N-137 | E | Sugerir receta con IA | Solo si la IA está activa (desactivada por defecto, decisión H). Devuelve un borrador que se revisa antes de guardar | Al generar |

### 3.14 P-12 Receta

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-140 | L | Receta escalada | Para las porciones elegidas (por defecto, las base), por ingrediente: producto, cantidad escalada, unidad, presentación comprable sugerida y cuántas piezas comprar, costo (motor S4, con tienda y fecha del precio) y kcal con su fuente (USDA, Open Food Facts o manual con cita). Totales: costo total, costo por porción, kcal total y kcal por porción. Pasos y notas | Inmediata |
| N-141 | L | Uso de la receta | Veces en el plan en las últimas 8 semanas y próxima fecha planeada | Inmediata |
| N-142 | E | Actualizar receta | Propiedades, ingredientes y pasos (edición en el lugar) | — |
| N-143 | E | Duplicar receta | Copia con el nombre «(copia)» y origen Axel y Ximena | — |
| N-144 | E | Eliminar receta | Con confirmación; si aparece en planes de hoy en adelante, lo dice y pide confirmar | — |

«Agregar al plan» usa N-123.

### 3.15 P-13 Lista de compras

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-150 | L | Lista de compras de la semana | Lista de la semana que manda la pantalla (por defecto, la lista vigente de §0.2). **Artículos** = ingredientes del plan escalados − inventario disponible + reposición de básicos por ritmo de consumo + artículos agregados a mano (motor S4). Por artículo: producto, cantidad necesaria (unidad base), presentación y piezas a comprar, tienda o mercado asignado (el que recomienda S5 o el elegido), precio estimado con su frescura, motivo (recetas que lo piden, reposición o manual) y comprado sí o no. Agrupado por vertical (mercados y tianguis, cada uno con su día de visita; cadenas, por tienda), luego por categoría de producto; orden alfabético. Totales por grupo y total general | Diaria |
| N-151 | E | Marcar comprado | Marca o desmarca un artículo de la lista de la semana que muestra P-13 (semana explícita, §0.2). No crea gasto ni precio hasta confirmar el ticket | — |
| N-152 | E | Agregar artículo | A la lista de la semana que muestra P-13 (semana explícita, §0.2): producto del catálogo o texto libre, cantidad y tienda | — |
| N-153 | E | Cambiar tienda o cantidad de un artículo | En la lista de la semana que muestra P-13 (semana explícita, §0.2): asigna otra tienda o mercado, o cambia la cantidad a comprar | — |

«Terminar compra» abre P-05 con N-073.

### 3.16 P-14 Comparar canasta

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-160 | L | Canastas disponibles | La lista de compras de la semana elegida (por defecto la vigente, §0.2) y las compras registradas con renglones (fecha, tienda, total), por fecha descendente | Inmediata |
| N-161 | L | Comparativa por tienda | Motor S5, canasta elegida, equivalencias por unidad base. Por tienda (Walmart, Bodega Aurrera, Sam's Club como si tuvieran membresía, y mercados con precio): total, artículos idénticos (mismo producto y presentación), sustitutos (presentación distinta comparada por unidad base), sin precio y precios viejos (pesan menos según S5). Orden por total ascendente. Diferencia contra la más barata, en total y en idénticos | Diaria |
| N-162 | L | Detalle por artículo | Por artículo: cantidad, precio en cada tienda con fecha y fuente, precio por unidad base, tienda ganadora, diferencia, marca de sustituto (con su nota) y marca de viejo. Orden por diferencia descendente | Diaria |
| N-163 | L | Ganadores | Número de artículos en que gana cada tienda y número de empates | Diaria |
| N-164 | L | Membresía Sam's | Motor S5: ahorro de la canasta en Sam's contra la mejor tienda sin Sam's, ahorro anual estimado (con la frecuencia de compra que defina S5), costo anual de la membresía (P-26.06), punto de equilibrio y veredicto con la cifra | Diaria |
| N-165 | E | Usar esta tienda | Asigna la tienda elegida a los artículos que tienen precio en ella, en la lista de la semana que se está comparando (semana explícita; por defecto la vigente, §0.2) | — |

Cifras canon que esta pantalla debe poder mostrar con la compra semilla del 10-oct-2026 (criterio S5 y `data.json`): Bodega Aurrera $3,094.50 contra Walmart $3,233.00, diferencia +$138.50; en idénticos $2,773.50 contra $2,874.00, diferencia +$100.50; 33 idénticos y 3 sustitutos; Bodega gana 22, empates 8, Walmart gana 3. «Simular pago con crédito» abre P-20 con el total de la canasta.

### 3.17 P-15 Productos y precios

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-170 | L | Catálogo | Por producto: nombre, marca, categoría de producto, presentaciones (unidad base y factor), mejor precio fresco por unidad base y su tienda, último precio (fecha), frescura (fresco, viejo o sin precio) y sustitutos. Filtros: tienda, categoría, frescura, fuente, texto. Orden: nombre. El tablero agrupa por frescura | Diaria |
| N-171 | L | Precios de un producto | Por tienda y presentación: precio, precio por unidad base, fecha, fuente (manual, ticket, sitio o PROFECO), días desde la fecha y marca de viejo. Historial (Tendencias): precio por unidad base por fecha y tienda | Diaria |
| N-172 | L | Resumen de precios | Número de productos con precio fresco, viejo y sin precio; por tienda, fecha del precio más reciente | Diaria |
| N-173 | L | Nutrición del producto | kcal por 100 g o 100 ml y su fuente con cita | Inmediata |
| N-174 | E | Capturar precio | Producto y presentación, tienda, precio y fecha (hoy por defecto); fuente = manual | — |
| N-175 | E | Crear o editar producto | Producto, presentaciones (unidad base y factor), sustitutos y nutrición manual con cita | — |
| N-176 | E | Corregir o eliminar precio | Corrige precio, fecha o tienda de un precio capturado a mano o por ticket, o lo elimina con confirmación. Los precios de PROFECO no se editan | — |

### 3.18 P-16 Inventario

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-180 | L | Inventario | Por artículo: producto, modo (stock por unidad o ritmo de consumo, P-26.14), existencia y unidad; si es por ritmo, consumo estimado por semana y días restantes (motor S4); última entrada (fecha y origen: compra, sobrante o ajuste); estado (suficiente, reponer o agotado; motor S4). Filtros: modo, estado, texto. Orden: agotado, reponer, suficiente y luego nombre. El tablero agrupa por estado | Diaria |
| N-181 | L | Movimientos de un artículo | Fecha, tipo (entrada por compra, sobrante de receta, consumo del plan, ajuste), cantidad y origen; orden por fecha descendente | Inmediata |
| N-182 | L | Resumen | Conteo de artículos por estado | Diaria |
| N-183 | E | Ajustar existencia | Cantidad nueva o diferencia, y motivo (conteo, consumo, merma) | — |
| N-184 | E | Registrar sobrante | Producto, cantidad y receta de origen (el motor S4 también lo propone desde el plan) | — |

### 3.19 P-17 Todas las líneas (Créditos)

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-190 | L | Líneas | Por línea activa: nombre, tipo, límite o monto original, deuda viva (capital + intereses por pagar), disponible (revolvente: lectura activa A o B con su etiqueta; tarjeta: límite − saldo; préstamo y tienda: no aplica), próximo pago (fecha y monto) y compras vivas. Orden por fecha del próximo pago ascendente | Diaria |
| N-191 | L | Totales | Deuda viva total; pagos con fecha en la quincena en curso y en el mes en curso (total, parte del hogar, parte de Karen); intereses pagados en el mes y en el año calendario | Diaria |
| N-192 | L | Calendario unificado | Motor S3. Pagos de todas las líneas en el rango (el mes visible; en móvil, agenda de hoy a 60 días): fecha, línea, compra, monto, parte del hogar, parte de Karen y estado de cada parte (pagada, pendiente, o sin marcar con fecha pasada). Promociones vigentes como franjas (comercio, plazo, compra mínima, vigencia) y fechas de corte de tarjetas. Orden por fecha. Filtros: línea, familia, estado | Diaria |
| N-193 | L | Aviso A/B | Lectura activa (A o B; B por defecto) y texto fijo «Reposición de capital pendiente de confirmar con Kueski» | Inmediata |
| N-194 | L | Tendencias de deuda | Motor S3. Por mes (12 meses): deuda viva proyectada al cierre, por línea; intereses por mes (pagados y por pagar) | Inmediata |

### 3.20 P-18 Crédito (una línea)

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-200 | L | Línea | Propiedades según el tipo (nombre, tipo, límite o monto, tasa tal como se capturó, fechas de corte y de pago, comisiones, notas; K12) y, si es revolvente tipo Kueski, la lectura A/B activa. Tarjeta bancaria y crédito de tienda: los cinco datos de su estado de cuenta (N-283, K18), cada uno como se capturó o «Falta capturar» | Inmediata |
| N-201 | L | Resumen de la línea | Deuda viva, intereses por pagar, próximo pago (fecha, monto, partes) y compras vivas. Revolvente: disponible de hoy con lectura A y con lectura B (la activa destacada, la otra en segundo plano) y aviso N-193. Tarjeta bancaria y crédito de tienda (K18): los intereses por pagar y el pago mínimo salen de la fórmula general de S3 con los cinco datos de N-283 (RN-26); si falta alguno, las cifras que dependen de ellos dicen «Falta capturar · Capturar» (abre P-26.05 en esa línea) y no se calculan con valores supuestos | Diaria |
| N-202 | L | Compras de la línea | Revolvente, tarjeta y tienda. Por compra: comercio, fecha, capital, intereses cotizados (IVA incluido), total, plazo, promoción, partes (hogar, Karen), pagos hechos de cuántos y saldo. Filtros: vivas o liquidadas, familia. Orden por fecha descendente | Inmediata |
| N-203 | L | Serie de disponible | Solo revolvente, motor S3. Eventos en orden de fecha: fecha, tipo (inicio, compra, pago), descripción, monto, capital repuesto, interés no repuesto, disponible con lectura A y con lectura B; mínimos del horizonte para A y para B. Canon semilla (`model.json`): mínimo A $8,094.50 y mínimo B $5,000.00; series A y B de 14 valores | Inmediata |
| N-204 | L | Tasas por plazo | Solo revolvente. Plazo en quincenas y tasa tal como se capturó (K12): 2 quincenas = 17.27 %, 4 quincenas = 25.8685 %; los demás plazos vacíos hasta capturarse en P-26.05 | Inmediata |
| N-205 | L | Estado de la tarjeta | Ciclo actual: fecha de corte, fecha límite de pago, saldo al corte, pago para no generar intereses, pago mínimo y saldo actual; compras a meses sin intereses (comercio, mensualidad, meses restantes). El pago mínimo sale de la fórmula general de S3 con los cinco datos de N-283 (K18, RN-26); si falta alguno, dice «Falta capturar · Capturar» (abre P-26.05 en esa línea) y no se calcula con valores supuestos | Diaria |
| N-206 | L | Amortización | Solo préstamo. Por pago: número, fecha, pago, capital, interés, saldo y estado; el último pago ajusta centavos | Inmediata |
| N-207 | L | Abonos | Solo crédito de tienda. Por abono: número, fecha, monto, saldo y estado | Inmediata |
| N-208 | L | Calendario de la línea | N-192 filtrado por la línea | Diaria |
| N-209 | L | Tendencias de la línea | Revolvente: disponible A y B por fecha. Todas: deuda viva al cierre de cada mes | Inmediata |
| N-210 | E | Marcar pago | Para un pago programado, marca por separado la parte del hogar y la parte de Karen como pagadas (cada quien paga directo a Kueski, decisión A), con fecha real (hoy por defecto) y monto (el programado por defecto). La reposición de capital de la línea la calcula S3. Permite desmarcar | — |

### 3.21 P-19 Compra a crédito

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-215 | L | Compra | Línea, comercio, fecha, capital, intereses cotizados (IVA incluido), total, plazo, promoción aplicada, partes por persona (hogar y Karen: capital, intereses, total y pago por quincena; motor S3, según lo que cada quien compró), pagos (fecha, monto, parte del hogar, parte de Karen, capital repuesto, estado de cada parte) y gasto vinculado | Inmediata |
| N-216 | L | Vista previa del cálculo | Con los datos capturados y sin guardar: partes y calendario de pagos que resultarían (motor S3), con el último pago ajustando centavos, y si la promoción aplica (comercio, compra mínima, vigencia) | Inmediata (en cada cambio) |
| N-217 | L | Promociones aplicables | Promociones vigentes en la fecha de compra para el comercio y la línea, con plazo y compra mínima | Diaria |
| N-218 | E | Crear compra | Línea, comercio, fecha, capital, intereses cotizados (de la cotización de Kueski) o, si no hay cotización, plazo con la tasa de la tabla de P-26.05; plazo; promoción; reparto por persona (desde los renglones del ticket o capital por persona capturado); gasto vinculado (lo crea si no existe) | — |
| N-219 | E | Editar o eliminar compra | Con confirmación; eliminar quita sus pagos programados | — |

Canon semilla (`model.json`, despensa 1): capital $3,094.50, intereses $800.50 a 4 quincenas, total $3,895.00; hogar (parte «Axel» del modelo) $2,395.28 con $598.82 por quincena; Karen $1,499.72 con $374.93 por quincena; pagos de $973.75 el 30-oct, 15-nov, 30-nov y 15-dic-2026.

### 3.22 P-20 Comparador de escenarios

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-225 | L | Valores por defecto | Restricciones duras globales de P-26.05 (N-288, K16): colchón mínimo de disponible en la línea que se usa, máximo de compras vivas a la vez en el crédito que se usa y fecha opcional de deuda en cero en el crédito que se usa; se pueden cambiar para esta comparación sin tocar las globales. No tienen valor por defecto: si no están capturadas, P-20 muestra «Faltan restricciones» y no recomienda. Tope de pago por quincena de cada persona, de P-26.12 (N-309; opcional y no es restricción dura, K13), y pesos de los 4 objetivos (iguales y normalizados) | Inmediata |
| N-226 | L | Resultado del comparador | Motor S3, búsqueda exhaustiva sin solver. **Entrada**: monto, fecha deseada ± días, comercio o tienda, quién compra (Hogar, Karen o Conjunto con %), restricciones, tope y pesos. **Salida**: escenarios válidos (línea × plazo × fecha × promoción) con intereses totales, pago por quincena (máximo y diferencia entre quincenas), disponible mínimo resultante, compras vivas máximas, fecha en que la deuda queda en cero y rendimiento de compra por peso (tienda más barata y promociones); el óptimo de cada uno de los 4 objetivos, lado a lado; la recomendación por puntaje ponderado con su puntaje; los descartados con la restricción que violan; y el total evaluado. Si faltan restricciones duras, no hay recomendación y la salida lo indica (K16). Tiempo máximo del cálculo: 5 s (**FIRME**, K43; se cambia en Ajustes › Créditos y condiciones, K44, N-321); si se pasa, P-20 avisa «Comparador no disponible» sin resultados parciales | Al generar (en cada cambio de entrada) |

«Usar escenario» abre P-19 prellenada; no escribe nada hasta guardar ahí.

### 3.23 P-21 Salud

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-230 | L | Perfiles | Por persona (Axel, Ximena): altura, peso actual (último registro y su fecha), edad, sexo, actividad, IMC con su clasificación (motor S4), objetivo (bajar, mantener o subir de peso), meta de calorías por día confirmada o escrita por la persona y la sugerida por S4 (Mifflin-St Jeor × factor de actividad de 5 niveles, 1.2 a 1.9, menos el ajuste de bajar o más el de subir; «mantener» no ajusta; RN-17, K36), con «Confirmar» si difieren (K23); el ajuste por objetivo de la persona, para bajar y para subir, en kcal: su valor o, si está vacío, «El de la casa: n kcal» (N-326, K36); sin meta confirmada, N-023, N-231 y N-233 muestran «—» con «Meta sin confirmar» (§0.2) y restricciones alimentarias (vacías hoy, decisión D) | Inmediata |
| N-231 | L | Calorías de hoy | Por persona: kcal planeadas hoy (suma de las casillas de hoy con receta), meta y diferencia; los tiempos «fuera» se marcan sin kcal | Diaria |
| N-232 | L | Historial de peso | Registros (fecha, peso, IMC a esa fecha) por persona, por fecha descendente. En Tendencias: serie de peso por persona | Inmediata |
| N-233 | L | kcal promedio por semana | Últimas 12 semanas: kcal promedio por día de cada persona (definición de N-023) contra su meta | Inmediata |
| N-327 | E | Ajuste por objetivo de una persona | En la ficha de cada persona (Axel, Ximena), dos campos en kcal (≥ 0): ajuste para bajar y ajuste para subir (K36). Valor inicial: el de la casa (N-326; 0 y 0 kcal al empezar), que se muestra junto al campo; vacío = usa el de la casa. «Restablecer» vacía el campo y vuelve a usar el de la casa. Sin valor inventado: un valor no válido no se guarda y el campo vuelve al anterior. La meta sugerida de N-230 se recalcula con el ajuste del objetivo elegido; la meta confirmada no cambia sola (K23) | — |

«Editar perfil» abre P-26.13 (N-311, N-313). El ajuste por objetivo de cada persona se edita aquí, en el lugar (N-327), y también en su ficha de P-26.13.

### 3.24 P-22 Registrar peso

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-240 | E | Registrar peso | Persona (por defecto el usuario de la sesión), fecha (hoy) y peso en kg (la precisión la fija S1). Acepta la clave de idempotencia y la fecha y hora de captura de una captura hecha sin conexión (N-317, K9): reenviar la misma clave no duplica | — |
| N-241 | L | Último peso | Último registro de la persona elegida (peso y fecha), para mostrar «antes: … el …» | Inmediata |
| N-242 | E | Corregir o eliminar peso | Desde el historial de P-21: corrige fecha o peso de un registro, o lo elimina con confirmación | — |

### 3.25 P-23 Centro de reportes

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-245 | L | Reportes generados | Tipo (mensual PDF o estado de cuenta de Karen), periodo, fecha de generación, origen (programado o manual, y quién), estado (listo, o error con motivo) y envío por Telegram (a quién, cuándo, resultado). Filtros: tipo, periodo. Orden por fecha descendente | Inmediata |
| N-246 | L | Próximos envíos | Resumen diario, semanal y mensual: siguiente fecha y hora, y destinatarios (P-26.08, P-26.09) | Diaria |
| N-247 | L | Avisos enviados | Fecha y hora, tipo (resumen diario, semanal o mensual; alerta de pago próximo, de presupuesto al 80 %, al 100 %, de precios viejos o de plan sin armar), destinatarios, resultado de la entrega y texto enviado (texto plano). Filtros: tipo, periodo. Orden por fecha descendente | Inmediata |
| N-248 | E | Generar reporte mensual | Mes; opción de enviarlo por Telegram a los destinatarios configurados. El resultado aparece en N-245 | Al generar |
| N-249 | L | Descargar archivo | El PDF guardado de un reporte | — |
| N-250 | E | Reintentar envío | Reenvía por Telegram un reporte cuyo envío falló | — |

### 3.26 P-24 Reporte mensual

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-255 | L | Meses disponibles | Del primer mes con datos al mes en curso; el mes en curso lleva la marca «parcial» | Inmediata |
| N-256 | L | Reporte del mes | Mes elegido (por defecto el último cerrado). Secciones con las mismas cifras que el PDF (S7 cruza PDF contra API): resumen (ingreso, real contra plan, ahorro), presupuesto por categoría, flujo por quincena, créditos (pagos del mes, deuda viva al cierre, intereses), despensa y precios (comparativa del mes), comidas y salud (costo de comidas, kcal promedio contra meta, peso al inicio y al cierre) y participación de Karen. Se calcula al consultar; el PDF es una copia con su fecha de generación | Inmediata |
| N-257 | L | PDF del mes | Si existe: fecha de generación y archivo (N-249) | Inmediata |

«Generar PDF» usa N-248.

### 3.27 P-25 Estado de cuenta de Karen

| Id | Tipo | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|
| N-260 | L | Estado de cuenta | Periodo elegido (por defecto el mes en curso; también un mes anterior o un rango). Compras particulares de Karen a crédito (comercio, fecha, capital, intereses, total, pago por quincena); su parte de cada despensa común (fecha, su capital en conjunto y propio, intereses, total); sus partes de pagos con fecha en el periodo y las pendientes posteriores (fecha, monto, estado); total pagado en el periodo; saldo pendiente total. Nota fija sobre el pago directo a Kueski | Inmediata |
| N-261 | E | Generar PDF | Periodo. Guarda el archivo y lo lista en N-245; desde ahí se descarga o se comparte con la hoja de compartir de Android | Al generar |
| N-262 | L | Estados anteriores | N-245 filtrado a «estado de cuenta de Karen» | Inmediata |

### 3.28 P-26 Ajustes

| Id | Tipo | Sección | Métrica o lista | Definición | Actualización |
|---|---|---|---|---|---|
| N-270 | L | 01 | Usuarios | Axel y Ximena: nombre visible, usuario, rol «co-titular» y último acceso | Inmediata |
| N-271 | L | 01 | Sesiones | De ambos usuarios (co-titulares, ven lo mismo): dispositivo y navegador, inicio, último uso y si es la sesión actual. Orden por último uso descendente | Inmediata |
| N-272 | E | 01 | Cerrar sesiones | Una sesión, todas las de un usuario salvo la actual, o todas | — |
| N-273 | E | 01 | Cambiar contraseña | Solo la propia: contraseña actual y nueva | — |
| N-274 | L | 01 | Este dispositivo | Si la app está instalada (dato del navegador, Local) y el certificado raíz público de la CA local para descargar, con instrucciones (S6) | Inmediata |
| N-320 | E | 01 | Parámetros de acceso | Intentos fallidos de inicio de sesión permitidos en la ventana (valor inicial 5), ventana deslizante por usuario e IP (valor inicial 15 min) duración de la sesión, renovable con el uso (valor inicial 90 días) (K11, K44), y largo mínimo de una contraseña nueva (valor inicial 12 caracteres; K6; lo usa N-273). Lee cada valor vigente con su valor inicial; guarda o restablece (reglas de §2.3, «Parámetros operativos») | — |
| N-275 | L | 02 | Presupuesto mensual | Por categoría activa: monto del mes en curso y, si es de «Solo este mes» (K2), esa marca y el monto que rige desde el mes siguiente; total del mes en curso | Inmediata |
| N-276 | L | 02 | Gastos fijos recurrentes | Concepto, categoría, monto, día del mes, reparto y activo; orden por día | Inmediata |
| N-277 | L | 02 | Despensa programada | Monto base, día de la despensa y frecuencia bimestral, línea y plazo con que se paga, próxima fecha | Inmediata |
| N-278 | E | 02 | Editar presupuesto, fijos y despensa | Crear, editar y desactivar. Cambiar el monto del presupuesto de una categoría pide el alcance de N-087 («De este mes en adelante» preseleccionado; aquí «este mes» es el mes en curso; K2) | — |
| N-318 | L | 02 | Proyección de los meses 7 a 12 | Interruptor «Repetir la despensa bimestral y el papel»: apagado (por defecto) = los meses 7 a 12 de P-07 llevan solo compromisos contratados, gastos fijos e ingresos del CSV; encendido = se agrega la cadencia bimestral de despensa a crédito y el papel alterno (K21) | Inmediata |
| N-319 | E | 02 | Cambiar el interruptor de proyección | Enciende o apaga N-318; P-07 se recalcula al volver a pedir sus necesidades | — |
| N-279 | L | 03 | Ingreso esperado | Monto esperado de la quincena del 15 y de la del 30 (pueden diferir) y desde cuándo rige | Inmediata |
| N-280 | E | 03 | Editar ingreso esperado | Montos y fecha desde la que rigen | — |
| N-281 | L | 04 | Categorías | Nombre, tipo (fijo o variable), reparto por defecto (Hogar, Karen o Conjunto con %), orden, activa y número de gastos | Inmediata |
| N-282 | E | 04 | Editar categorías | Crear, renombrar, reordenar (arrastrar), cambiar reparto y archivar. Una categoría con gastos no se borra: se archiva | — |
| N-283 | L | 05 | Créditos | Líneas con sus condiciones según el tipo (límite, tasa, días de corte y de pago, comisiones, plazo, abonos) y si están activas. Tarjeta bancaria y crédito de tienda (K18): los cinco datos que Axel toma de su estado de cuenta, sin valor por defecto: tasa anual tal como se capturó, con sus decimales (K12); base diaria o mensual; si la tasa incluye IVA (si no lo incluye, el IVA que se suma a los intereses); porcentaje de pago mínimo, y mínimo fijo. Mientras falte alguno, la línea muestra «Falta capturar» y cuáles faltan | Inmediata |
| N-284 | L | 05 | Lectura A/B | Valor del interruptor (B por defecto) y aviso «pendiente confirmar» | Inmediata |
| N-285 | L | 05 | Vista previa A/B | Motor S3: disponible de hoy y mínimo del horizonte con A y con B, para ver el efecto antes de cambiar | Inmediata |
| N-286 | L | 05 | Tasas por plazo | Plazo (quincenas) y tasa (%) tal como se capturó, con sus decimales (K12). Valores conocidos: 2 = 17.27 %, 4 = 25.8685 %. Los cálculos usan la fracción exacta | Inmediata |
| N-287 | L | 05 | Promociones | Comercio, línea, plazo, tasa promocional por todo el plazo tal como se capturó («sin intereses» si es 0; K12), compra mínima, vigencia (inicio y fin), fuente o nota y si está vigente hoy. Vista calendario por vigencia; tabla ordenada por inicio | Diaria |
| N-288 | L | 05 | Valores del comparador | Restricciones duras globales (K16): colchón mínimo (medido en la línea que se usa), máximo de compras vivas (en el crédito que se usa) y fecha de deuda en cero (opcional; en el crédito que se usa). Las captura Axel y no tienen valor por defecto: mientras falten, se muestra «Falta capturar». Pesos de los 4 objetivos. El tope por persona está en P-26.12 (N-309, K13) | Inmediata |
| N-289 | E | 05 | Editar créditos y condiciones | Líneas (en tarjeta bancaria y crédito de tienda, los cinco datos del estado de cuenta de N-283, sin valor por defecto; si se guarda sin alguno, la línea queda en «Falta capturar», K18), interruptor A/B, tasas (se guardan con los decimales tecleados, K12; también la tasa anual), promociones, y restricciones duras y pesos del comparador (K16) | — |
| N-321 | E | 05 | Tiempo máximo del comparador | Segundos que puede tardar el cálculo de N-226 antes de responder «no disponible» (valor inicial 5 s; K43, K44). Lee el valor vigente con su valor inicial; guarda o restablece (reglas de §2.3) | — |
| N-290 | L | 06 | Tiendas | Cadenas y sucursales (Walmart, Bodega Aurrera, Sam's Club), mercados y tianguis (nombre, días de visita, qué se compra ahí) y membresía (Sam's: costo anual; «no la tienen», se evalúa) | Inmediata |
| N-291 | L | 06 | Fuentes de precio | Captura manual (siempre activa); OCR local (Tesseract, español); lectura asistida por sitio (desactivada; solo se puede activar si `docs/estado/S5-terminos.md` tiene veredicto favorable para ese sitio, que se muestra); PROFECO (último corte importado, fecha, registros y cobertura por cadena y municipio) | Inmediata |
| N-292 | E | 06 | Editar tiendas y fuentes | Crear y editar tiendas y mercados; activar la lectura asistida de un sitio solo si su veredicto lo permite | — |
| N-293 | E | 06 | Importar corte PROFECO | Archivo local (RAR o CSV de datos abiertos); devuelve registros cargados y cobertura | Al generar |
| N-322 | E | 06 | Parámetros de precios | Días para que un precio sea viejo (valor inicial 30; K22) y días en que su peso llega a cero y deja de competir (valor inicial 90; K40). S5 aplica RN-12 con los valores vigentes. Tamaño máximo de un corte PROFECO importado (valor inicial 1 GB; K6; lo usa N-293). Lee cada valor vigente con su valor inicial; guarda o restablece (reglas de §2.3; K44) | — |
| N-294 | L | 07 | Recetario y nutrición | Porciones base (2), tiempos de comida (desayuno, comida, cena, colación), número de recetas por origen y bases nutricionales locales (USDA Foundation Foods 04/2026, Open Food Facts filtrado, manuales) con fecha de carga y número de alimentos | Inmediata |
| N-295 | E | 07 | Recargar base nutricional | Desde el archivo local; devuelve el conteo | Al generar |
| N-296 | L | 08 | Reportes programados | Resumen diario (activo, hora de la noche), semanal (domingo, hora), mensual PDF (día 1, hora, enviar por Telegram, guardar) y sus destinatarios (Axel, Ximena) | Inmediata |
| N-297 | E | 08 | Editar reportes programados | Activo, hora (`send_time`, `HH:MM` en `America/Mexico_City`; **FIRME**, K43) con «Restablecer al valor inicial» cuando lo hay (K44, §2.3), y destinatarios | — |
| N-298 | L | 09 | Alertas | Pago próximo (día anterior y mismo día; anticipación `alerts.payment_due_days_before`: lista de días antes del pago, valor inicial [1, 0], decisión F; se edita en N-323), presupuesto al 80 % y al 100 % (categoría y total; umbrales en N-323, K44), precios viejos antes del día de despensa (anticipación `advance_days`, días enteros ≥ 0, solo de este aviso), plan de comidas sin armar (domingo): activa, destinatarios, anticipación y hora (`send_time`, `HH:MM` en `America/Mexico_City`) donde aplica, cada una con su valor inicial donde lo haya (claves `send_time` y `advance_days` **FIRMES**, K43; K44). los días antes del pago próximo y los umbrales son claves de `app_setting` desde el esquema 1.3 (00003, K44) y ya no van en `notification.params_json`; dónde se guardan `send_time` y `advance_days` lo decide el bloqueo abierto `S1-parametros-avisos` y la pantalla no cambia | Inmediata |
| N-299 | L | 09 | Registro de chat | Por persona (Axel, Ximena): chat confirmado sí o no y fecha. **Chats detectados** por el bot MultiX cuando alguien le escribe y aún sin confirmar (K4): nombre que muestra Telegram y fecha y hora del mensaje. Instrucciones para escribirle al bot y confirmar aquí (el texto exacto lo fija S7 según cómo detecte el chat). Cada chat detectado se confirma (N-316) o se descarta (N-324). Nunca se muestra el token | Inmediata |
| N-300 | E | 09 | Editar alertas y probar | Activa, destinatarios, anticipación y hora donde aplica: `advance_days` del aviso de precios viejos y `send_time`, en `params_json` (claves `send_time` y `advance_days` **FIRMES**, K43), con «Restablecer al valor inicial» cuando lo hay (K44, §2.3); la anticipación del pago próximo (`alerts.payment_due_days_before`, valor inicial [1, 0], decisión F) y los umbrales del presupuesto se editan en la misma pantalla con N-323 («Restablecer» vuelve a [1, 0] y a 80 % y 100 %); «Enviar aviso de prueba» a una persona | — |
| N-316 | E | 09 | Confirmar chat detectado | Asigna un chat detectado (N-299) a una persona (por defecto el usuario de la sesión); desde entonces los avisos de esa persona van a ese chat (K4) | — |
| N-323 | E | 09 | Umbrales del presupuesto y anticipación del pago próximo | Porcentajes del presupuesto (por categoría y total) en que salta el aviso (`alerts.budget_warning_pct` y `alerts.budget_limit_pct`, valores iniciales 80 % y 100 %; K44) y días antes del pago en que salta el aviso de pago próximo (`alerts.payment_due_days_before`, valor inicial [1, 0]: 1 día antes y el mismo día; decisión F, K44). Lee cada valor vigente con su valor inicial; guarda o restablece (reglas de §2.3). Son claves de `app_setting` desde el esquema 1.3 (00003) y ya no van en `notification.params_json`; la pantalla no cambia | — |
| N-324 | E | 09 | Descartar chat detectado | Borra, con confirmación, un chat detectado (N-299) que nadie va a confirmar; no toca los chats ya confirmados (S1 1.2: descartar borra el candidato) | — |
| N-301 | L | 10 | Tema | Sistema (por defecto), claro u oscuro; por dispositivo | Local |
| N-302 | L | 11 | Respaldos | Lista (fecha y hora, tamaño, resultado de `integrity_check`, ubicación), último respaldo correcto, próximo programado (diario) y fecha de la última restauración probada | Inmediata |
| N-303 | E | 11 | Respaldar y restaurar | «Respaldar ahora»; «Restaurar» pide doble confirmación (acción crítica) | Al generar |
| N-304 | E | 11 | Importar gastos fijos e ingresos | `gastos_fijos_ingresos.csv` (plantilla de `plan.md` §8.2) con vista previa de lo que se creará y confirmación | Al generar |
| N-305 | L | 11 | Datos semilla | Estado de carga de la compra del 10-oct-2026, del crédito Kueski y del CSV de gastos fijos e ingresos | Inmediata |
| N-306 | L | 11 | IA opcional | Modelo local (Ollama) y API de Claude: activa sí o no (desactivadas por defecto), usos permitidos (OCR de tickets, sugerencias de recetas) y disponible sí o no. La llave nunca se muestra; solo «configurada» | Inmediata |
| N-307 | E | 11 | Activar o desactivar IA | Por proveedor | — |
| N-308 | L | 11 | Estado del servicio | Versión, en marcha desde y tamaño de la base de datos | Inmediata |
| N-325 | E | 11 | Parámetros de datos | Tamaño máximo de un adjunto de gasto (valor inicial 10 MB; K43; lo usa N-064) y retención de las claves de idempotencia (valor inicial 24 h; K6; las usan N-054 y N-240). Lee cada valor vigente con su valor inicial; guarda o restablece (reglas de §2.3; K44) | — |
| N-309 | L | 12 | Personas y reparto | Axel (titular, usuario), Ximena (co-titular, usuaria) y Karen (participante externa, sin cuenta), con lo que hace cada quien en el sistema; reparto por defecto del conjunto (50/50 Hogar y Karen); tope de pago por quincena por persona, opcional y vacío mientras Axel no lo capture (K13; solo lo usa el comparador, N-225) | Inmediata |
| N-310 | E | 12 | Editar personas y reparto | Nombres visibles, % por defecto del conjunto y tope por quincena de cada persona (se puede dejar vacío, K13) | — |
| N-311 | L | 13 | Perfiles de salud | Por persona (Axel, Ximena): altura, fecha de nacimiento o edad (S1), sexo, actividad (5 niveles), objetivo (bajar, mantener o subir de peso), meta de kcal sugerida (motor S4: Mifflin-St Jeor × factor de actividad de 1.2 a 1.9, menos el ajuste de bajar o más el de subir; «mantener» no ajusta; RN-17, K36), ajuste por objetivo de la persona (bajar y subir en kcal; vacío = el de la casa, N-327, K36), meta confirmada o escrita por la persona (K23) y restricciones (vacías) | Inmediata |
| N-312 | L | 13 | Metas y fondo | Metas (nombre, monto, fecha, aportación mensual planificada, activa) y aportación mensual planificada al fondo general | Inmediata |
| N-313 | E | 13 | Editar perfiles, metas y fondo | Crear, editar y desactivar. En el perfil: elegir el objetivo y confirmar la meta sugerida o escribir otra (K23; también desde P-21) | — |
| N-326 | E | 13 | Ajuste por objetivo de la casa | Dos valores en kcal (≥ 0): ajuste para bajar y ajuste para subir (valores iniciales 0 kcal y 0 kcal; K36, K44). Son el valor inicial del ajuste de cada persona que lo tenga vacío (N-327). Lee cada valor vigente con su valor inicial; guarda o restablece (reglas de §2.3) | — |
| N-314 | L | 14 | Configuración de inventario | Por producto: modo (stock por unidad o ritmo de consumo), unidad y umbral de reposición o ritmo inicial. Regla fija: los sobrantes de recetas entran como stock (decisión D) | Inmediata |
| N-315 | E | 14 | Editar inventario | Modo, unidad y umbrales | — |

### 3.29 Matriz pantalla ↔ necesidades (para C1)

C1 completa la última columna con el o los endpoints de cada necesidad. Una necesidad reutilizada se lista en cada pantalla que la usa. En v1.2 hay **211** necesidades (N-001 a N-327; 138 de lectura y 73 de escritura); N-317 es Local (vive en el teléfono), igual que N-003 y N-301.

| Pantalla | Lectura | Escritura | Endpoint C1 |
|---|---|---|---|
| Global (G-01 a G-07) | N-001, N-002, N-003, N-004, N-005, N-317 | N-006 | Lo llena C1 |
| P-00 | N-011 | N-010 | Lo llena C1 |
| P-01 | N-020 a N-029, y por reutilización N-081, N-085, N-090, N-091, N-111, N-121, N-129, N-162, N-192, N-194 | N-054 (vía P-03) | Lo llena C1 |
| P-02 | N-040 a N-046 | N-047, N-062, N-063 | Lo llena C1 |
| P-03 | N-050 a N-053, N-317 | N-054, N-055, N-102, N-240 | Lo llena C1 |
| P-04 | N-060, N-061 | N-062 a N-066 | Lo llena C1 |
| P-05 | N-071, N-072, N-073, N-076 | N-070, N-074, N-075 | Lo llena C1 |
| P-06 | N-080 a N-086 | N-087 | Lo llena C1 |
| P-07 | N-090 a N-093 | — | Lo llena C1 |
| P-08 | N-100, N-103, N-104 | N-101, N-102, N-105 | Lo llena C1 |
| P-09 | N-110, N-111, N-112, N-114, N-115 | N-113, N-116 | Lo llena C1 |
| P-10 | N-120, N-121, N-122, N-128, N-129 | N-123 a N-127 | Lo llena C1 |
| P-11 | N-135 | N-136, N-137 | Lo llena C1 |
| P-12 | N-140, N-141 | N-123, N-142, N-143, N-144 | Lo llena C1 |
| P-13 | N-150 | N-151, N-152, N-153 | Lo llena C1 |
| P-14 | N-160 a N-164 | N-165 | Lo llena C1 |
| P-15 | N-170 a N-173 | N-174, N-175, N-176 | Lo llena C1 |
| P-16 | N-180, N-181, N-182 | N-183, N-184 | Lo llena C1 |
| P-17 | N-190 a N-194 | — | Lo llena C1 |
| P-18 | N-200 a N-209 | N-210 | Lo llena C1 |
| P-19 | N-215, N-216, N-217 | N-210, N-218, N-219 | Lo llena C1 |
| P-20 | N-225, N-226 | — | Lo llena C1 |
| P-21 | N-230 a N-233 | N-242, N-313, N-327 | Lo llena C1 |
| P-22 | N-241, N-317 | N-240 | Lo llena C1 |
| P-23 | N-245, N-246, N-247, N-249 | N-248, N-250 | Lo llena C1 |
| P-24 | N-255, N-256, N-257, N-249 | N-248 | Lo llena C1 |
| P-25 | N-260, N-262, N-249 | N-261 | Lo llena C1 |
| P-26.01 | N-270, N-271, N-274 | N-272, N-273, N-320 | Lo llena C1 |
| P-26.02 | N-275, N-276, N-277, N-318 | N-278, N-319 | Lo llena C1 |
| P-26.03 | N-279 | N-280 | Lo llena C1 |
| P-26.04 | N-281 | N-282 | Lo llena C1 |
| P-26.05 | N-283 a N-288 | N-289, N-321 | Lo llena C1 |
| P-26.06 | N-290, N-291 | N-292, N-293, N-322 | Lo llena C1 |
| P-26.07 | N-294 | N-295 | Lo llena C1 |
| P-26.08 | N-296 | N-297, N-248 | Lo llena C1 |
| P-26.09 | N-298, N-299 | N-300, N-316, N-323, N-324 | Lo llena C1 |
| P-26.10 | N-301 (Local, sin API) | — | No requiere |
| P-26.11 | N-302, N-305, N-306, N-308 | N-303, N-304, N-307, N-325 | Lo llena C1 |
| P-26.12 | N-309 | N-310 | Lo llena C1 |
| P-26.13 | N-311, N-312 | N-313, N-326, N-327 | Lo llena C1 |
| P-26.14 | N-314 | N-315 | Lo llena C1 |

### 3.30 Entidades de §3.4.1 y las pantallas que las usan (para S1)

| Entidad | Pantallas |
|---|---|
| persona | P-26.12, P-02 (filtro), P-21, P-25 |
| cuenta_usuario | P-00, P-26.01 |
| sesion_dispositivo | P-26.01 |
| periodo (mes, quincena) | P-01, P-06, P-07, P-08 |
| ingreso | P-08, P-07, P-06, P-26.03 |
| categoria | P-26.04, P-03, P-02, P-06 |
| gasto | P-02, P-03, P-04, P-05, P-06, P-01 |
| regla_reparto | P-26.12, P-26.04, P-03, P-04 |
| meta_ahorro, movimiento_ahorro | P-09, P-26.13, P-01 |
| credito, tasa_por_plazo | P-17, P-18, P-26.05 |
| compra_credito | P-19, P-18, P-25, P-01 |
| pago_programado | P-17, P-18, P-19, P-01, P-07, P-25 |
| promocion | P-26.05, P-17, P-19, P-20 |
| tienda | P-26.06, P-13, P-14, P-15 |
| producto, presentacion | P-15, P-05, P-12, P-13 |
| precio | P-15, P-14, P-13, P-05 |
| receta, ingrediente_receta | P-11, P-12, P-10 |
| tiempo_comida | P-10, P-26.07 |
| plan_comida | P-10, P-01, P-21 |
| inventario_item | P-16, P-13, P-26.14 |
| perfil_salud, peso_historico | P-21, P-22, P-26.13 |
| nutricion_ingrediente | P-12, P-15, P-26.07 |
| notificacion, destinatario | P-26.09, P-23 |
| job | P-05, P-23, P-26.06, P-26.11 |
| respaldo | P-26.11 |

**Necesidades sin entidad explícita en §3.4.1** (S1 decide si son tablas, columnas o datos derivados): renglones de ticket con su cuenta (N-061, N-074); lista de compras con marcas de comprado y artículos manuales (N-150 a N-153); lecturas OCR sin confirmar (N-075, N-076); gastos fijos recurrentes (N-276); despensa programada (N-277); estado de pago por parte, hogar y Karen (N-210); destino de un movimiento de ahorro, fondo o meta (N-112, N-113); reporte generado con archivo, periodo y envío (N-245), si `job` no alcanza.

---

## 4. Flujos clave (d)

**Cómo se cuentan los toques en móvil** (regla (a), confirmada por Axel el 11-oct-2026, K10). Un toque es una pulsación sobre un control de la interfaz (botón, chip, fila, pestaña, casilla): abrir, elegir, guardar. No cuentan: teclear valores en el teclado (monto, peso), los valores por defecto, desplazar la pantalla, ni los diálogos del sistema operativo (permiso de cámara). Punto de partida: la app instalada, abierta en P-01 con la sesión ya iniciada (sesión larga por dispositivo). Ancho de referencia: 412 px (U5); el camino crítico es el mismo en 360 px. Los montos de los ejemplos marcados «(ejemplo)» son ilustrativos; los demás salen de `oraculo/model.json` o `datos/semilla/data.json`.

### F-1 Registrar gasto (desde Inicio) · criterio ≤ 3 toques

| Paso | Pantalla | Qué hace la persona | Qué hace el sistema | Toques acumulados |
|---|---|---|---|---|
| 1 | P-01 → P-03 | Toca «+» en la barra inferior (o «Registrar gasto», la acción primaria de Inicio) | Abre la hoja con el tipo «Gasto», el monto enfocado y el teclado numérico abierto. Fecha = hoy, persona = usuario de la sesión, forma de pago = «De contado» (N-052); los valores por defecto no cuentan como toques (K10). Muestra las 6 categorías sugeridas (N-050) | 1 |
| 2 | P-03 | Teclea el monto, por ejemplo `185.50` (ejemplo) | Valida en vivo que sea mayor que 0. Teclear no cuenta (K10) | 1 |
| 3 | P-03 | Toca un chip de categoría | Marca la categoría y aplica su reparto por defecto (N-051) | 2 |
| 4 | P-03 → P-01 | Toca «Guardar» | Crea el gasto (N-054), cierra la hoja y muestra «Gasto guardado · Deshacer» durante 5 s (N-055). Inicio vuelve a pedir N-020 y N-021 | **3** |

**Demostración de que el camino crítico tiene exactamente 3 toques:**

1. El «+» está siempre visible en la barra inferior de Inicio y no está dentro de ningún menú (toque 1).
2. El foco del monto y el teclado numérico se abren dentro del mismo toque 1: el campo se enfoca en el manejador del toque. Si el navegador no abre el teclado, la hoja muestra su propio teclado numérico, sin toques extra. U5 lo verifica en Android a 412 × 915.
3. De los campos obligatorios de N-054 (monto, categoría, fecha, persona, forma de pago), tres tienen valor por defecto (fecha, persona, forma de pago) y el reparto sale de la categoría. El monto se teclea, y teclear no cuenta como toque.
4. La categoría es el único campo obligatorio que requiere toque: un chip entre 6 (toque 2).
5. «Guardar» es un solo toque y no pide confirmación (toque 3). El error se corrige con «Deshacer» o editando en P-04, sin diálogo previo.
6. No hay que desplazar la hoja: con el teclado abierto, el monto, los 6 chips y «Guardar» caben en el área visible desde 360 px (requisito de diseño para U2; lo verifica U5).

Por lo tanto: 1 (abrir) + 1 (categoría) + 1 (guardar) = **3 toques** ≤ 3. Es la regla que Axel confirmó (K10): teclear y los valores por defecto no cuentan.

**Variantes (fuera del camino crítico):** categoría que no está entre las 6: «+», «Más», elegir de la lista y «Guardar» (4 toques). Cambiar persona, reparto, fecha o forma de pago: «Más opciones» (+1) y +1 por campo. Pagar con una línea de crédito: «Más opciones» → «Forma de pago» → línea, y continúa en P-19 con los datos prellenados. En escritorio: Alt+G abre P-03 y Ctrl+Enter guarda.

**Variante «Sin conexión» (K9, U1-15):** P-01 muestra la franja «Sin conexión con el servidor de casa» (G-07). Los pasos y los toques son los mismos: «+», categoría (de las ya cargadas) y «Guardar» en P-03 = **3 toques**. En el paso 4 el gasto se guarda en la cola del teléfono con su clave de idempotencia y su fecha y hora de captura (N-317); el aviso dice «Gasto guardado en el teléfono · Se enviará al volver a la red de casa · Deshacer» y la franja suma uno a «en cola». Al volver la conexión se envía solo y sin duplicar (N-054); si falla, sigue en la cola con «No se pudo enviar · Revisar». Sin conexión, «Forma de pago» solo ofrece «De contado»: pagar con una línea de crédito requiere conexión (P-19).

### F-2 Planear la semana de comidas

| Paso | Pantalla | Qué hace la persona | Qué hace el sistema | Toques |
|---|---|---|---|---|
| 1 | P-01 → P-10 | Toca «Comidas» en la barra inferior | Abre el plan en Detalle en la semana por defecto (§0.2; N-120, N-121) | 1 |
| 2 | P-10 | Opcional: «⋯» → «Copiar semana anterior» | Llena solo las casillas vacías (N-126) | +2 |
| 3 | P-10 | Toca una casilla vacía (por ejemplo, Lun · Comida) | Abre el selector de recetas de ese tiempo, con la fila de días arriba y el día de la casilla ya marcado (N-122) | +1 |
| 4 | P-10 | Opcional: toca otros días en la fila de días | Marca en qué otros días va la misma receta | +1 por día |
| 5 | P-10 | Toca una receta | Asigna la receta a Ambos, 1 porción por persona, en los días marcados, y cierra el selector (N-123). La casilla muestra costo y kcal | +1 |
| 6 | P-10 | Opcional: para comer fuera, casilla → «Fuera de casa» | Marca «fuera» sin receta ni kcal (N-124) | +2 |
| 7 | P-10 | Opcional: toca la pestaña «Resumen» | Costo de la semana, kcal por día contra meta y casillas sin armar (N-121) | +1 |
| 8 | P-10 → P-13 | Toca «Generar lista de compras» | Genera la lista neteada con el inventario (N-127) y la abre (N-150) | +1 |

**Conteo:** toques = 1 (entrar) + Σ por receta distinta [1 casilla + días extra + 1 receta] + 2 por cada casilla «fuera» + 1 (generar), y +1 si se revisa el Resumen. Ejemplo: 21 casillas obligatorias cubiertas con 7 recetas distintas de 3 días cada una = 1 + 7 × (1 + 2 + 1) + 1 = **30 toques**. Copiando la semana anterior y cambiando 4 casillas (una receta nueva cada una, sin días extra) = 1 + 2 + 4 × 2 + 1 = **12 toques**.

### F-3 Comparar tienda para una canasta

**Ruta A, la lista de esta semana (desde Inicio):**

| Paso | Pantalla | Qué hace la persona | Qué hace el sistema | Toques |
|---|---|---|---|---|
| 1 | P-01 → P-14 | Toca la tarjeta «Mejor canasta» | Abre Comparar canasta con la lista vigente (§0.2), en Resumen: total por tienda, diferencia, cobertura, precios viejos y membresía Sam's (N-160, N-161, N-163, N-164) | 1 |
| 2 | P-14 | Toca la pestaña «Detalle» | Precio de cada artículo en cada tienda, ganador y marcas de sustituto y de viejo (N-162) | 2 |
| 3 | P-14 → P-13 | Toca «Usar Bodega Aurrera» (la primaria lleva el nombre de la ganadora) | Asigna la tienda a la lista (N-165) y abre la lista agrupada por esa tienda | **3** |

**Ruta B, otra canasta (por ejemplo, la compra del 10-oct-2026):**

| Paso | Pantalla | Qué hace la persona | Qué hace el sistema | Toques |
|---|---|---|---|---|
| 1 | P-01 → P-13 | Toca «Compras» en la barra inferior | Abre la lista de compras | 1 |
| 2 | P-13 → P-14 | «⋯» → «Comparar tiendas» | Abre Comparar canasta con la lista de la semana que muestra P-13 (§0.2) | 3 |
| 3 | P-14 | Toca el selector de canasta y elige «Compra del 10-oct-2026» | Compara esa canasta (N-160, N-161) | **5** |

Resultado esperado con la semilla: Bodega Aurrera $3,094.50 contra Walmart $3,233.00 (+$138.50); en idénticos +$100.50; 33 idénticos y 3 sustitutos.

### F-4 Ver un crédito y su calendario de pagos

| Paso | Pantalla | Qué hace la persona | Qué hace el sistema | Toques |
|---|---|---|---|---|
| 1 | P-01 → P-18 | Toca la tarjeta «Próximo pago» | Abre la línea de ese pago en Resumen: deuda viva, disponible A y B con aviso «pendiente confirmar», próximo pago (N-200, N-201) | 1 |
| 2 | P-18 | Toca la pestaña «Calendario» | Agenda de pagos de la línea (N-208). Con la semilla, la primera fila es «Vie 30 oct 2026 · $973.75 · hogar $598.82 · Karen $374.93 · pendiente» | **2** |
| 3 | P-18 → G-06 | Toca un pago | Vista previa: compra de origen, capital repuesto, estado de cada parte | 3 |
| 4 | G-06 | Opcional: «Marcar pago» → marca la parte del hogar o la de Karen → «Guardar» | Marca la parte como pagada y recalcula el disponible (N-210, motor S3) | 6 |

Desde la barra lateral: ☰ → Créditos → la línea → «Calendario» = 4 toques.

### F-5 Ver o generar un reporte

| Paso | Pantalla | Qué hace la persona | Qué hace el sistema | Toques |
|---|---|---|---|---|
| 1 | P-01 → G-01 | Toca ☰ | Abre la barra lateral como panel | 1 |
| 2 | G-01 → P-23 | Toca «Reportes» | Abre el Centro de reportes (N-245, N-246, N-247) | 2 |
| 3 | P-23 → P-24 | Toca la fila «Mensual · sep-2026» | Abre la vista web del mes (N-255, N-256) | **3** (ver) |
| 4 | P-24 | Toca «Descargar PDF» si ya existe, o «Generar PDF» si no | Descarga el PDF (N-249) o lo genera con «Generando…» (N-248) | 4 |
| 5 | P-24 | Si se generó: toca «Descargar» o «Compartir» | Descarga, o abre la hoja de compartir de Android | 5 (generar y obtener) |

**Estado de cuenta de Karen:** ☰ (1) → «Reportes» (2) → «Estado de cuenta de Karen» (3) → periodo por defecto el mes en curso → «Generar PDF» (4, N-261) → «Compartir» (5). Axel lo reenvía a mano desde la hoja de compartir de Android (decisión A); el sistema no se lo manda a Karen.

### F-6 Capturar peso

| Paso | Pantalla | Qué hace la persona | Qué hace el sistema | Toques |
|---|---|---|---|---|
| 1 | P-01 → P-03 | Toca «+» en la barra inferior | Abre la hoja de registro con el tipo «Gasto» | 1 |
| 2 | P-03 → P-22 | Toca «Peso» en el selector de tipo | Muestra el campo de peso enfocado con teclado numérico; persona = usuario de la sesión; fecha = hoy; último peso de esa persona (N-241) | 2 |
| 3 | P-22 | Teclea el peso | No cuenta como toque | 2 |
| 4 | P-22 → P-01 | Toca «Guardar» | Registra el peso (N-240), muestra «Peso guardado · Deshacer» y actualiza IMC e historial en P-21 | **3** |

Variante: peso de la otra persona, +1 (chip «Ximena» o «Axel»). Desde Salud: ☰ → «Salud» → «Registrar peso» → «Guardar» = 4 toques.

**Variante «Sin conexión» (K9, U1-15):** mismos 3 toques. «Guardar» en P-22 deja el peso en la cola del teléfono con su clave de idempotencia y su fecha y hora de captura (N-317); el aviso dice «Peso guardado en el teléfono · Se enviará al volver a la red de casa · Deshacer». «Antes: …» (N-241) muestra el último peso cargado o el último en cola de esa persona. Al volver la conexión se envía sin duplicar (N-240) y entonces se actualizan el IMC y el historial de P-21; si falla, sigue en la cola con «No se pudo enviar · Revisar».

### Flujos complementarios

Los 6 flujos anteriores son los que pide la tarea de U1. Los criterios de U1 (U1-05) también piden como obligatorios el comparador de escenarios (decisión C) y el ticket con OCR (decisión E); se agregan como F-7 y F-8 para que ambas listas queden cubiertas (ver «Solicitudes fuera de alcance» del reporte).

**F-7 Planear una compra a crédito con el comparador**

| Paso | Pantalla | Qué hace la persona | Qué hace el sistema | Toques |
|---|---|---|---|---|
| 1 | P-01 → G-01 | Toca ☰ | Abre la barra lateral | 1 |
| 2 | G-01 → P-17 | Toca «Créditos» | Abre Todas las líneas | 2 |
| 3 | P-17 → P-20 | «⋯» → «Simular compra» | Abre el comparador con las restricciones globales y los pesos de P-26.05 y el tope por persona de P-26.12 (N-225). Si faltan restricciones, muestra «Faltan restricciones · Capturar» (abre P-26.05) y no recomienda (K16) | 4 |
| 4 | P-20 | Teclea el monto; opcional: cambia fecha deseada, holgura en días o quién compra | Evalúa todos los escenarios en cada cambio (N-226) | 4 (+1 por dato cambiado) |
| 5 | P-20 | Revisa los 4 óptimos (apilados en móvil) y la recomendación con su puntaje | Muestra también los descartados y la restricción que violan. Solo hay recomendación con las restricciones capturadas (paso 3, K16) | 4 |
| 6 | P-20 → P-19 | Toca «Usar escenario recomendado» (u «Usar» en otro óptimo) | Abre P-19 prellenada con línea, plazo, fecha, promoción y monto; no guarda nada | 5 |
| 7 | P-19 | Si la compra es hoy: teclea los intereses que cotizó Kueski y toca «Guardar» | Muestra la vista previa (N-216) y crea la compra con sus pagos (N-218) | **6** |

**Sin restricciones capturadas** (K16): mientras Axel no capture en P-26.05 las restricciones globales (no tienen valor por defecto, N-225), en el paso 3 P-20 muestra «Faltan restricciones · Capturar» (abre P-26.05) y no recomienda; por eso no hay recomendación en el paso 5 ni «Usar escenario recomendado» en el paso 6. Es lo que verá el primer uso real.

Atajo desde una canasta: tarjeta «Mejor canasta» (1) → «⋯» → «Simular pago con crédito» (3, con el total de la canasta ya puesto) → «Usar escenario recomendado» (4) → «Guardar» (5). Si faltan restricciones, el atajo se queda en P-20 con «Faltan restricciones · Capturar» (abre P-26.05) y sin recomendación (K16).

**F-8 Foto de ticket → OCR → confirmación → precios guardados**

| Paso | Pantalla | Qué hace la persona | Qué hace el sistema | Toques |
|---|---|---|---|---|
| 1 | P-01 → P-03 | Toca «+» | Abre la hoja de registro | 1 |
| 2 | P-03 → P-05 | Toca «Ticket» en el selector de tipo | Abre la cámara del teléfono | 2 |
| 3 | P-05 | Toma la foto y la acepta en la cámara | Sube la foto y lee el ticket con OCR local, «Leyendo el ticket…» (N-070) | 4 |
| 4 | P-05 | Revisa los renglones propuestos; los de confianza baja aparecen desplegados | Propone producto, cantidad, precio y la cuenta de cada renglón: la que tuvo ese producto en la última compra confirmada o, si no hay, el reparto de la categoría (N-071) | 4 |
| 5 | P-05 | Opcional: corrige un renglón (tocarlo y elegir el producto correcto) o cambia su cuenta | Busca en el catálogo (N-072) | +2 por corrección, +2 por cuenta cambiada |
| 6 | P-05 | Opcional: si pagó con crédito, elige la línea y teclea intereses cotizados y plazo | Prepara la compra a crédito | +2 |
| 7 | P-05 → P-04 | Toca «Confirmar ticket» | En una sola operación crea gasto, renglones, precios con fuente «ticket», entradas de inventario y, si aplica, la compra a crédito (N-074); abre el gasto creado | **5** sin correcciones |

Nada se guarda sin el paso 7 (confirmación humana, decisión E). Si la lectura falla, se conserva la foto y se ofrece capturar a mano. «Guardar para después» (N-075) deja la lectura con las correcciones hechas como borrador en «Tickets por confirmar» (N-076); al abrirla, P-05 vuelve con esas correcciones (K3).

### Resumen de toques en móvil

| Flujo | Camino más corto | Toques |
|---|---|---|
| F-1 Registrar gasto | Inicio → «+» → categoría → Guardar | **3** |
| F-2 Planear la semana | Comidas → copiar semana anterior → 4 cambios → generar lista | 12 (semana nueva completa: 30) |
| F-3 Comparar canasta | Tarjeta «Mejor canasta» → Detalle → Usar tienda | 3 (solo ver: 1) |
| F-4 Ver crédito y calendario | Tarjeta «Próximo pago» → Calendario | 2 |
| F-5 Ver o generar reporte | ☰ → Reportes → mes | 3 para ver; 5 para generar y descargar |
| F-6 Capturar peso | «+» → Peso → Guardar | 3 |
| F-7 Compra a crédito con comparador | ☰ → Créditos → Simular compra → Usar escenario → Guardar | 6 (desde una canasta: 5) |
| F-8 Ticket con OCR | «+» → Ticket → foto → aceptar → Confirmar | 5 sin correcciones |

---

## 5. Menú «/» y atajos de teclado (e)

### 5.1 Comportamiento del menú «/» (G-03)

1. **Cómo se abre**: en escritorio, tecleando «/» cuando el foco no está en un campo de texto; también desde la búsqueda (Ctrl+K) si el texto empieza con «/». En móvil, con la pestaña «Buscar» de la barra inferior: es el mismo panel de búsqueda y comandos.
2. **Dentro de las notas** (bloques de BlockNote en gastos, recetas y otras páginas), «/» abre el menú de bloques del editor, no este menú.
3. **Búsqueda del comando**: por nombre y por alias, sin distinguir mayúsculas ni acentos («credito» encuentra «crédito»). Primero aparecen los comandos de la pantalla actual y luego el resto por grupo. ↑ y ↓ mueven la selección, Enter ejecuta, Esc cierra.
4. **Nada se guarda desde el menú**: cada comando abre la pantalla u hoja indicada con sus valores por defecto; se guarda con la acción de esa pantalla.
5. Los comandos de vista solo aparecen en pantallas con base de datos.

### 5.2 Lista de comandos

| Comando | Alias | Hace | Abre | Dónde aparece |
|---|---|---|---|---|
| /gasto | g, gastar | Registrar gasto | P-03 | Siempre |
| /ticket | t, foto, ocr | Capturar ticket | P-05 | Siempre |
| /ingreso | extra | Registrar ingreso extra | Hoja de P-08 | Siempre |
| /quincena | q, sueldo | Capturar quincena de Axel | Hoja de P-08 | Siempre |
| /credito | kueski, msi, compra a credito | Registrar compra a crédito | P-19 | Siempre |
| /pago | abono | Marcar pago de crédito | Lista de pagos pendientes → N-210 | Siempre |
| /ahorro | aportar, retiro | Registrar movimiento de ahorro | Hoja de P-09 | Siempre |
| /peso | p, bascula | Registrar peso | P-22 | Siempre |
| /precio | — | Capturar precio | Hoja de P-15 | Siempre |
| /receta | — | Nueva receta | P-12 nueva | Siempre |
| /comida | plan, menu | Asignar receta al plan | Selector de P-10 | Siempre |
| /lista | comprar | Agregar artículo a la lista de compras | Hoja de P-13 | Siempre |
| /existencia | stock, sobrante | Ajustar existencia o registrar sobrante | Hoja de P-16 | Siempre |
| /promocion | promo | Nueva promoción | P-26.05 | Siempre |
| /comparar | canasta | Comparar canasta | P-14 | Siempre |
| /simular | escenarios | Comparador de escenarios | P-20 | Siempre |
| /reporte | pdf | Generar reporte mensual | P-23 (N-248) | Siempre |
| /karen | estado de cuenta | Estado de cuenta de Karen | P-25 | Siempre |
| /respaldo | backup | Respaldar ahora | P-26.11 (N-303) | Siempre |
| /ir | (nombre de la página) | Ir a cualquier página: Inicio, Gastos, Presupuesto, Flujo y proyección, Ingresos, Ahorro y metas, Plan de la semana, Recetario, Lista de compras, Comparar canasta, Productos y precios, Inventario, Créditos, cada línea por su nombre, Comparador, Salud, Reportes, Reporte mensual, Estado de cuenta de Karen, Ajustes y cada una de sus 14 secciones. Basta teclear el nombre | La página | Siempre |
| /tabla, /tablero, /galeria | — | Cambiar la vista de la base de datos | La pantalla actual en esa vista (Local) | Pantallas con base de datos |
| /resumen, /detalle, /tendencias, /calendario | — | Cambiar de subsección | La pantalla actual en esa subsección | Pantallas con subsecciones |
| /filtrar, /ordenar, /agrupar | — | Abrir los controles de la vista | Controles de la vista actual (Local) | Pantallas con base de datos |
| /tema | claro, oscuro | Alternar tema (sistema, claro, oscuro) | N-301 (Local) | Siempre |
| /barra | — | Mostrar u ocultar la barra lateral | G-01 (Local) | Escritorio |
| /atajos | ayuda | Ver la lista de atajos | Hoja de atajos (G-03) | Siempre |
| /salir | cerrar sesion | Cerrar sesión | N-006 | Siempre |

### 5.3 Atajos de teclado (escritorio)

«Notion» indica que el atajo y su acción coinciden con la ayuda oficial de Notion (notion.com/help/keyboard-shortcuts, consultada el 10-oct-2026); «propio» indica un atajo de esta app.

| Atajo | Acción | Dónde | Origen |
|---|---|---|---|
| `/` | Abrir el menú «/» | Fuera de campos de texto | Notion («/» abre el menú), adaptado a comandos |
| `Ctrl+K` o `Ctrl+P` | Buscar o saltar a una página reciente | Siempre | Notion |
| `Ctrl+[` / `Ctrl+]` | Atrás / adelante (también funcionan `Alt+←` / `Alt+→` del navegador) | Siempre | Notion |
| `Ctrl+Shift+U` | Subir a la página padre | Siempre | Notion |
| `Ctrl+Shift+L` | Alternar claro y oscuro | Siempre | Notion |
| `Ctrl+\` | Mostrar u ocultar la barra lateral (también el botón «‹‹» o /barra) | Escritorio | Propio |
| `Alt+G` | Registrar gasto desde cualquier pantalla | Siempre | Propio |
| `Alt+N` | Nuevo registro en la base de datos de la página actual | Pantallas con base de datos | Propio (Notion usa `Ctrl+N`, que el navegador reserva) |
| `Ctrl+Enter` | Guardar la hoja o el formulario abierto; con una fila seleccionada, abrirla | Hojas, formularios y tablas | Notion (`Ctrl+Enter` abre o modifica el bloque actual), adaptado |
| `Ctrl+/` | Editar a la vez las propiedades de las filas seleccionadas | Tablas y tableros | Notion |
| `Esc` | Cerrar menú, hoja o vista previa; quitar la selección | Siempre | Notion |
| `↑` / `↓` y `Enter` | Moverse entre filas o resultados y abrir el seleccionado | Tablas, menú, búsqueda | Notion |
| `←` / `→` | Semana o mes anterior / siguiente | Plan de la semana y calendarios (fuera de campos de texto) | Propio |
| `Ctrl+F` | Buscar dentro de la página (del navegador) | Siempre | Notion |
| `?` | Ver la lista de atajos | Fuera de campos de texto | Propio |

### 5.4 Reglas para que los atajos no choquen

- No se usan combinaciones que el navegador reserva o usa: `Ctrl+N`, `Ctrl+T`, `Ctrl+W`, `Ctrl+Tab` y sus variantes con Shift, `Ctrl+1` a `Ctrl+9`, `Alt+1` a `Alt+9`, `Alt+F`, `Alt+E`.
- Los atajos de una sola tecla (`/`, `?`, flechas) solo funcionan cuando el foco no está en un campo de texto.
- Se usa el Alt izquierdo; AltGr se ignora para no chocar con `@` y otros caracteres del teclado latinoamericano. Donde `\` o `[` requieren AltGr, quedan el botón, el comando /barra y `Alt+←` / `Alt+→`.
- La lista completa se ve con `?` o /atajos, y cada botón con atajo lo muestra en su descripción emergente.
- En Android no hay atajos: el «+» de la barra inferior y la pestaña «Buscar» cubren lo mismo.

---

## 6. Móvil primero: puntos de quiebre (f)

### 6.1 Reglas generales

1. Se diseña primero a 360 px y se amplía. Referencias de U5: Android 412 × 915 y escritorio 1440 × 900.
2. Rangos: **360** = hasta 411 px · **412** = de 412 a 767 · **768** = de 768 a 1023 · **1024** = de 1024 a 1439 · **1440** = desde 1440.
3. Ninguna pantalla tiene desplazamiento horizontal de la página en ningún ancho. Solo se desplazan horizontalmente los componentes que lo declaran: fila de pestañas, tablero y tabla con columnas agregadas por la persona.
4. **Objetivo táctil: al menos 44 × 44 px** (R2-05). Los 24 × 24 px de WCAG 2.2 AA (criterio 2.5.8) son solo el piso que nada puede bajar, no el objetivo de diseño; U2 fija los tamaños exactos a partir de 44 px.
5. Las cifras de dinero nunca se truncan ni se parten en dos líneas: si no caben, se quita información secundaria, no la cifra. La tarjeta más estrecha debe mostrar `$99,999.99` completo a 360 px.
6. Con el teclado en pantalla, las hojas se ajustan al área visible y su acción primaria queda arriba del teclado.
7. Teléfono en horizontal: se aplica la regla del ancho resultante (por ejemplo, 915 px usa la de 768).

### 6.2 Qué cambia en cada punto de quiebre

| Elemento | 360 | 412 | 768 | 1024 | 1440 |
|---|---|---|---|---|---|
| Barra lateral | Panel desde ☰, 85 % del ancho, con fondo atenuado | Igual que 360 | Panel superpuesto desde ☰, ancho fijo (U2) | Fija y colapsable (Ctrl+\) | Fija; páginas de texto con ancho máximo de lectura |
| Barra inferior (G-05) | Sí, 5 destinos | Sí | No; «+ Registrar» pasa a la barra superior | No | No |
| Tarjetas del Resumen (≤ 6) | 1 columna, orden fijo | 1 columna; cada tarjeta suma 1 dato secundario | 2 columnas (3 filas) | 3 columnas (2 filas) | 3 columnas (2 filas) con más aire |
| Pestañas de subsección | Con desplazamiento horizontal | Con desplazamiento horizontal | Fila completa | Fila completa | Fila completa |
| Tabla | Filas compactas de 2 líneas (título y cifra; 1 propiedad debajo); el resto en la vista previa | Filas compactas + 1 propiedad más | Hasta 5 columnas | Completa, columnas configurables | Completa |
| Tablero | 1 columna visible; se cambia deslizando o con el selector de grupo | Igual que 360 | 2 columnas visibles | 3 a 4 | Todas las que quepan |
| Calendario | Agenda por día | Agenda; mes compacto opcional | Mes | Mes | Mes con detalle del día seleccionado |
| Galería | 1 columna (tickets en 2) | 2 columnas | 3 | 3 a 4 | 4 a 5 |
| Plan de la semana | Lista por día: hoy desplegado, los demás plegados | Igual que 360 | 4 días visibles con desplazamiento | 7 días × 4 tiempos | 7 días × 4 tiempos, con costo y kcal en cada casilla |
| Comparador (4 óptimos) | Apilados | Apilados | 2 × 2 | 2 × 2 | 4 en una fila |
| Vista previa de registro | Página completa | Página completa | Panel central | Panel lateral derecho | Panel lateral derecho más ancho |
| Hojas y diálogos | Hoja inferior a pantalla completa | Hoja inferior | Diálogo centrado | Diálogo centrado | Diálogo centrado |
| Menú «/» y búsqueda | Hoja a pantalla completa | Igual que 360 | Panel centrado | Panel centrado | Panel centrado |
| Gráficas | Alto reducido, leyenda debajo, máximo 6 series, valor al tocar | Igual que 360 | Leyenda a la derecha, valor al pasar el puntero | Igual que 768 | Igual que 768 |
| Propiedades de un registro | 4 + «Mostrar N más» | 4 + «Mostrar N más» | Todas, en 1 columna | Todas | Todas, en 2 columnas |
| Migas de pan | «‹ Padre» | «‹ Padre» | Completas, recortadas al centro si no caben | Completas | Completas |

### 6.3 En palabras

- **360**: el mínimo. Una columna, barra inferior con «+», todo lo secundario en «⋯» u hojas. F-1 cabe en 3 toques sin desplazar.
- **412**: el teléfono de referencia. Igual que 360, con un dato secundario más por tarjeta y por fila, y galerías de 2 columnas.
- **768**: tableta vertical o ventana mediana. Desaparece la barra inferior y aparecen las 2 columnas de tarjetas, el calendario mensual y los diálogos centrados; la barra lateral sigue superpuesta.
- **1024**: escritorio pequeño. La barra lateral se fija, las tablas se muestran completas y la vista previa abre de lado, como en Notion.
- **1440**: escritorio de referencia. Mismo esquema con más aire: plan semanal con costo y kcal en cada casilla, comparador con los 4 óptimos en una fila y propiedades en 2 columnas.

### 6.4 Por pantalla: móvil (412 px) y escritorio (1440 px)

«Bajo demanda» es lo que no se muestra de entrada y aparece al tocar, desplegar o abrir (detalle progresivo).

| Pantalla | Móvil 412 | Escritorio 1440 | Acción primaria | Bajo demanda |
|---|---|---|---|---|
| P-00 | Formulario a todo el ancho | Formulario centrado de ancho fijo | Entrar | — |
| P-01 | 6 tarjetas en 1 columna; barra inferior con «+» | 6 tarjetas en 3 × 2; barra lateral fija | Registrar gasto | Detalle de cada tarjeta al tocarla; Detalle, Tendencias y Calendario en pestañas |
| P-02 | Filas compactas (descripción y monto; categoría y fecha debajo); filtros en hoja | Tabla completa con fila de totales; vista previa lateral | Registrar gasto | Propiedades extra en la vista previa; Resumen y Tendencias en pestañas |
| P-03 | Hoja inferior a pantalla completa con teclado numérico | Diálogo centrado | Guardar | «Más opciones» |
| P-04 | Página completa; 4 propiedades + «Mostrar N más» | Panel lateral o página; propiedades en 2 columnas | Adjuntar ticket o Ver ticket | Renglones plegados con su total |
| P-05 | Pasos en página completa: foto, renglones como tarjetas, confirmar | Foto a la izquierda, tabla de renglones a la derecha | Confirmar ticket | Renglones de confianza alta plegados |
| P-06 | Resumen en 1 columna; categorías como filas con barra | Resumen en 3 × 2 y tabla de categorías con barras | Ajustar presupuesto | Quincenas y despensa devengada al desplegar |
| P-07 | Quincenas como lista; gráfica debajo | Tabla de 24 quincenas y gráfica lado a lado | Simular compra | Supuestos de la proyección plegados |
| P-08 | Lista de quincenas y extras | Tabla con calendario a un lado | Capturar quincena | Tendencias en pestaña |
| P-09 | Metas como tarjetas en 1 columna | Galería de 4 columnas | Registrar aportación | Movimientos al abrir una meta |
| P-10 | Lista por día (hoy desplegado) | Cuadrícula de 7 días × 4 tiempos con costo y kcal por casilla | Generar lista de compras | Totales del día al desplegarlo |
| P-11 | Galería de 2 columnas | Galería de 5 columnas | Nueva receta | Filtros en hoja |
| P-12 | Porciones, ingredientes y pasos apilados | Ingredientes y pasos en 2 columnas | Agregar al plan | Fuente nutricional al tocar un ingrediente |
| P-13 | Grupos plegables por tienda con casillas grandes | Tabla agrupada por tienda con totales | Terminar compra | Motivo del artículo al tocarlo |
| P-14 | Totales por tienda apilados; artículos en pestaña Detalle | Tabla artículos × tiendas | Usar esta tienda | Membresía Sam's plegada en Resumen |
| P-15 | Filas compactas (producto y mejor precio) | Tabla completa | Capturar precio | Historial al abrir el producto |
| P-16 | Filas compactas ordenadas por estado | Tabla con tablero opcional | Ajustar existencia | Movimientos al abrir el artículo |
| P-17 | Líneas como tarjetas apiladas | Galería de 3 columnas | Registrar compra a crédito | Calendario en pestaña |
| P-18 | Resumen y pestañas | Resumen y tabla de compras lado a lado | Marcar pago | Serie A/B y tasas en Detalle |
| P-19 | Formulario con vista previa debajo | Formulario y vista previa lado a lado | Guardar (alta) o Marcar pago (detalle) | Calendario de pagos plegado |
| P-20 | Datos arriba; 4 óptimos apilados | Datos a la izquierda; 4 óptimos en una fila | Usar escenario recomendado | Tabla de todos los escenarios y descartados, plegada |
| P-21 | Axel y Ximena en pestañas | Dos columnas | Registrar peso | Historial en Tendencias |
| P-22 | Hoja inferior | Diálogo | Guardar | — |
| P-23 | Lista de reportes | Tabla con columnas de envío | Generar reporte mensual | Avisos enviados en sección plegada |
| P-24 | Secciones apiladas | Secciones con ancho de lectura | Descargar PDF | Cada sección plegable |
| P-25 | Vista previa en página | Vista previa con ancho de lectura | Generar PDF | Movimientos detallados plegados |
| P-26 | Lista de 14 secciones; cada sección en su página | Columna interna de secciones y formulario a la derecha | Según sección (2.3) | Vista previa del efecto (A/B, restauración) |
---

## 7. Abierto y dependencias

| Tema | Estado | Dónde |
|---|---|---|
| Inicio de semana y días de cada quincena | **FIRME** (K8): lunes a domingo; 1–15 y 16–último día del mes | Bloqueo `U1-periodos-semana-quincena`, resuelto |
| Captura sin conexión fuera de casa | **Decidido** (K9): cola en el teléfono para gastos y peso (§2.4, F-1, F-6, N-317) | Bloqueo `U1-captura-fuera-de-casa`, resuelto |
| Frontera de «precio viejo» (30 días o más) | **Decidido** (K22) | Bloqueo `criterios-frescura-peso`, resuelto |
| Definición de «casillas sin armar» (N-023) | **FIRME** (K43): pares día × tiempo obligatorio sin receta ni «fuera», máximo 21; es una regla, no un parámetro | Bloqueos `criterios-umbrales-ausentes` y `C1-valores-k31`, resueltos |
| Semana por defecto los domingos (la siguiente) para P-10, P-13 y P-14 | **FIRME** (K38): los domingos abren la semana que empieza el lunes; «Cambiar semana» siempre disponible (§0.2; N-120, N-150 y N-160 remiten ahí) | Bloqueo `U1-lista-vigente-domingo`, resuelto |
| Relación entre diferencia acumulada y deuda viva (N-091, N-092) | **Decidido** (K19): dos cifras distintas; deuda viva total = diferencia acumulada + saldo no devengado | Bloqueo `criterios-invariante-deuda`, resuelto |
| Meses 7 a 12 de la proyección (N-090, N-091) | **Decidido** (K21): interruptor N-318, apagado por defecto | Bloqueo `criterios-proyeccion-12m`, resuelto |
| Tarjeta bancaria y crédito de tienda: tasa anual, base, IVA, porcentaje de pago mínimo y mínimo fijo | **Decidido** (K18): los captura Axel por línea desde su estado de cuenta, sin valor por defecto; mientras falte alguno, «Falta capturar» (§2.4, P-26.05, P-18) | Bloqueo `criterios-credito-formulas`, resuelto; `reglas-negocio.md` RN-26 |
| Ajuste de la meta de kcal por objetivo | **Decidido** (K36): valores de la casa 0 y 0 kcal en P-26.13 (N-326); por persona en P-21 (N-327), vacío = el de la casa; «mantener» no ajusta | Bloqueo `S1-ajuste-objetivo-kcal`, resuelto; RN-17 |
| Valores de C1 que eran propuesta: adjunto, comparador, hora y anticipación de avisos | **Decidido** (K43): 10 MB (N-064), 5 s (N-226), `send_time` y `advance_days` (N-297, N-298, N-300) | Bloqueo `C1-valores-k31`, resuelto |
| Parámetros operativos configurables | **Decidido** (K44): cada uno en su sección de Ajustes, con valor inicial y «Restablecer al valor inicial» (§2.3, N-320 a N-323, N-325, N-326). La ubicación por sección es lectura PROPUESTA del orquestador y Axel la ve en la puerta 2; los mínimos y máximos de los de seguridad los propone S1 para esa puerta. Los umbrales del presupuesto y los días antes del pago próximo son claves `alerts.*` de `app_setting` desde el esquema 1.3 (00003; N-323); dónde se guardan `send_time` y `advance_days` lo decide el bloqueo abierto `S1-parametros-avisos` (no cambia la pantalla) | `decisiones.md` K44; `docs/design/U4-ajustes.md` |
| Diseño visual con Claude Design (Pista 3) | **Decidido** (K35): lo ejecuta Axel con el paquete de `docs/design/entrada/` que se entrega al cierre de WF-1; lo que descargue va en `docs/design/salida/`. El orquestador no abre la herramienta. U1 no entra como archivo en el paquete: sus pantallas llegan por U3, U4 y DESIGN.md | `decisiones.md` K35 |
| Reglas de devengo, reparto, promedios, costo de receta, meta de kcal, IMC, frescura ponderada | Las fijan S1, S2, S4 y S5 en `reglas-negocio.md`; U1 solo muestra el resultado | `docs/contratos/reglas-negocio.md` |
| Campos, validaciones y valores por defecto de Ajustes (hora del resumen diario, anticipación de alertas) | U4; los parámetros operativos siguen las reglas comunes de §2.3 (K44) | `docs/design/U4-ajustes.md` |
| Especificación de las 6 tarjetas | U3, a partir de N-020 a N-026 | `docs/design/U3-dashboard.md` |
| Endpoints de cada necesidad | C1 llena la columna de 3.29 | `docs/contratos/C1-datos.md`, `openapi.yaml` |

Observación para S1 y S2: el efectivo de las ventanas P1 a P6 de `model.json` (del día 10 al día 9) coincide al centavo con los meses calendario oct-2026 a mar-2027, porque las compras caen el día 10 y los pagos el 15 y el 30. Con un devengo que asigne cada despensa al mes de compra y al siguiente, la proyección por mes calendario de N-091 reproduce los totales canon de $14,362.48 (devengado) y $13,103.78 (efectivo).
