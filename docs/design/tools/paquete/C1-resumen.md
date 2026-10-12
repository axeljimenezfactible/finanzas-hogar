# C1 · Resumen del contrato de datos de Finanzas Hogar

Resumen en lenguaje llano de lo que la aplicación guarda, calcula y muestra. Sirve para que el diseño de las pantallas use los nombres, los datos y los formatos reales. No contiene rutas, ni código, ni ejemplos de la interfaz de programación. Los valores de color, tipografía y espaciado están en `C2-tokens.json` y `DESIGN.pdf`; este resumen no los repite.

Versión del contrato: C1 1.2 (11-oct-2026). Versión del resumen: 1.0.

## 1. Qué es Finanzas Hogar

Una aplicación web de uso doméstico, instalable en el teléfono Android, para el hogar de Axel y Ximena y la participación de Karen. Lleva el presupuesto del mes, los gastos (con foto del ticket), los ingresos por quincena, el ahorro y las metas, los créditos y sus pagos, el plan de comidas con su costo y sus calorías, la lista de compras con la mejor tienda, el inventario de la despensa, la salud (peso y calorías) y los reportes mensuales. Todo está en pesos mexicanos y en español de México. Funciona en la red de casa; no depende de servicios de terceros para mostrarse.

Quién usa qué: Axel y Ximena tienen cuenta y ven lo mismo. Karen no tiene cuenta: aparece como persona que participa en una parte de la despensa y de ciertas compras, y recibe su estado de cuenta en PDF.

## 2. Módulos y pantallas

La barra lateral (en el teléfono, un panel que se abre desde el botón de menú; además, una barra inferior de cinco destinos) tiene este orden fijo. Entre paréntesis, el identificador de la pantalla.

| Módulo | Pantallas | Qué se hace ahí |
|---|---|---|
| Buscar y Registrar gasto | Búsqueda con atajo Ctrl+K, hoja «Registrar gasto» (P-03) | Encontrar cualquier cosa; registrar un gasto, un ingreso, un peso o un ticket en tres toques. |
| Inicio | Inicio (P-01) con las pestañas Resumen, Detalle, Tendencias y Calendario | Ver el estado del hogar en seis tarjetas y bajar al detalle. |
| Gastos | Todos los gastos (P-02), Gasto (P-04), Capturar ticket (P-05) | Consultar, filtrar y corregir gastos; foto del ticket, lectura local y confirmación. |
| Presupuesto | Presupuesto del mes (P-06), Flujo y proyección (P-07), Ingresos (P-08), Ahorro y metas (P-09) | Plan contra real por categoría y quincena; proyección a 12 meses; quincenas de ingreso; metas de ahorro. |
| Comidas | Plan de la semana (P-10), Recetario (P-11), Receta (P-12) | Asignar recetas por día y tiempo de comida; ver costo y calorías. |
| Compras y tiendas | Lista de compras (P-13), Comparar canasta (P-14), Productos y precios (P-15), Inventario (P-16) | Qué comprar y dónde; mejor tienda; precios con su frescura; qué hay en casa. |
| Créditos | Todas las líneas (P-17), una línea (P-18), una compra a crédito (P-19), Comparador de escenarios (P-20) | Deuda viva, disponible, calendario de pagos y la mejor forma de pagar una compra. |
| Salud | Salud (P-21) y la hoja Registrar peso (P-22) | Perfil, índice de masa corporal, meta de calorías y peso. |
| Reportes | Centro de reportes (P-23), Reporte mensual (P-24), Estado de cuenta de Karen (P-25) | Reportes en PDF, envíos y estado de cuenta de Karen. |
| Ajustes | Ajustes (P-26) con 14 secciones | Configurar el sistema (sección 3). |
| Fuera de la barra | Iniciar sesión (P-00) | Entrar con usuario y contraseña; sesión larga por dispositivo. |

Las pantallas de lista (Todos los gastos, Presupuesto, Ingresos, Plan de la semana, Productos y precios, Inventario, Créditos) se presentan como vistas de base de datos al estilo de Notion: tabla, tablero, calendario o galería, con filtros, orden y agrupación; la vista por defecto de cada una está marcada en el documento U3 y en DESIGN.

Las cuatro pestañas de Inicio (y de varias pantallas de lista) son siempre las mismas y en este orden: Resumen, Detalle, Tendencias, Calendario. Algunas pantallas no usan las cuatro.

## 3. Las 14 secciones de Ajustes

01 Cuenta y roles · 02 Presupuesto · 03 Ingresos · 04 Categorías · 05 Créditos y condiciones (incluye promociones) · 06 Tiendas y precios (incluye mercados y membresía de Sam's) · 07 Recetas · 08 Reportes y frecuencias · 09 Notificaciones · 10 Apariencia · 11 Datos y respaldos · 12 Personas y reparto · 13 Salud y metas de ahorro · 14 Inventario.

Cada sección edita en el lugar con guardado automático y muestra una vista previa del efecto antes de confirmar cambios que mueven cifras. Todo valor que el sistema usa para decidir (por ejemplo cuándo avisar que un presupuesto va al 80 % o al 100 %, cuántos días antes de un pago avisar, qué tan viejo es un precio) vive en su sección de Ajustes con un valor inicial; nada está fijo en el programa.

## 4. Entidades principales

El sistema maneja 31 entidades. Agrupadas por tema:

| Tema | Entidades |
|---|---|
| Personas y acceso | persona; cuenta de usuario; sesión por dispositivo |
| Dinero que entra y sale | gasto (con sus renglones, su reparto entre personas y sus adjuntos); regla de reparto; categoría; ingreso y partidas planeadas; periodo (mes y quincena) |
| Ahorro | meta de ahorro; movimiento de ahorro |
| Créditos | crédito (línea revolvente, tarjeta, préstamo o crédito de tienda); tasa por plazo; compra a crédito; pago programado (con la parte de cada persona); promoción |
| Compras | tienda o mercado; producto; presentación comprable; precio; inventario y sus movimientos |
| Comidas y salud | receta y sus ingredientes; tiempo de comida (desayuno, comida, cena y colación); plan de comidas; nutrición por ingrediente; perfil de salud; peso histórico |
| Avisos y sistema | notificación; destinatario; tarea programada; respaldo |

Ideas que importa conocer al diseñar:

- **El hogar y Karen.** Un gasto o un pago puede repartirse entre personas. La «parte del hogar» es lo que cuenta para el presupuesto; la parte de Karen se muestra aparte y alimenta su estado de cuenta.
- **Devengado y efectivo.** Una compra a crédito se «devenga» cuando ocurre, pero el efectivo sale cuando se pagan las parcialidades. Las pantallas de flujo muestran las dos lecturas y nunca las mezclan en una cifra.
- **Estados de un pago.** Pendiente, por vencer, vencido, parcial y pagado. Cada estado lleva icono y texto; nunca solo color.
- **Frescura de un precio.** Un precio viejo pesa menos al comparar tiendas y se marca como «viejo».
- **Lo que falta capturar** se dice con esas palabras («falta capturar la tasa de esta tarjeta»), no con un error técnico.

## 5. Métricas del dashboard

Las seis tarjetas del Resumen de Inicio (orden fijo, ninguna tarjeta adicional):

| # | Tarjeta | Qué responde | Cifras principales |
|---|---|---|---|
| 1 | Mes contra plan | ¿Cuánto me queda del presupuesto de este mes? | Lo que te queda, real contra plan, porcentaje usado, días que faltan, estado (normal, aviso o rebasado). |
| 2 | Flujo de la quincena | ¿Cuánto entra y cuánto sale en esta quincena? | Ingresos, egresos ya hechos, comprometido y saldo de la quincena. |
| 3 | Próximo pago y disponible | ¿Qué pago sigue y cuánto puedo gastar? | Fecha, monto y línea del próximo pago; disponible de cada línea revolvente con el aviso «pendiente confirmar» cuando aplica. |
| 4 | Comidas de la semana | ¿Cuánto cuesta lo que voy a comer? | Costo de la semana, calorías promedio por día y por persona, casillas sin armar. |
| 5 | Mejor canasta | ¿Dónde me conviene comprar? | Tienda ganadora, diferencia contra la segunda, cobertura de precios. |
| 6 | Ahorro y metas | ¿Cómo voy con mis metas? | Total ahorrado, aportado contra planificado del mes, meta más próxima con su porcentaje. |

Pestaña Detalle: la misma información con más renglones, más la participación de Karen (su parte de la despensa del mes, sus compras vivas, sus partes de pagos con estado y su deuda viva).

Pestaña Tendencias (cuatro gráficas): plan contra real por mes (12 meses); devengado contra efectivo; deuda viva proyectada (12 meses, por línea); costo semanal y calorías de las comidas (12 semanas).

Pestaña Calendario: ingresos, pagos de crédito, gastos fijos y venideros, fechas de corte de tarjeta, despensa programada, días de visita a mercados y promociones vigentes.

Otras métricas que alimentan pantallas: presupuesto por categoría; flujo por quincena (24 quincenas hacia adelante); devengado contra efectivo por mes; deuda viva a una fecha; calendario unificado de pagos; plan contra real por mes de una categoría; costo semanal y calorías; avance de cada meta con su aportación mensual necesaria; deuda viva proyectada.

## 6. Formatos

- **Dinero**: pesos mexicanos con símbolo, separador de miles y dos decimales, por ejemplo `$1,258.68`. Las cifras son exactas al centavo; nunca se redondean para mostrarlas ni se abrevian (nada de «1.2 mil»). Cifras tabulares y alineadas a la derecha en tablas. Los egresos llevan signo menos y flecha de salida; los ingresos, signo más y flecha de entrada.
- **Fechas**: en pantalla `10-oct-2026`; donde importa el día de la semana (calendarios, pagos, plan de comidas), `Vie 30 oct 2026`. Zona horaria de la Ciudad de México.
- **Porcentajes de avance** (presupuesto, metas): números enteros con espacio antes del signo, por ejemplo `80 %`. Las **tasas de interés** se muestran como se capturaron, con los mismos decimales (`17.27 %`).
- **Idioma**: español de México, trato de tú, frases cortas, sin signos de exclamación ni emojis en mensajes del sistema. El éxito es discreto: «Guardado · Deshacer».
- **Personas y tiendas**: siempre con su nombre en texto junto al color; nunca se distinguen solo por color.
- **Estados vacíos y de error**: dicen qué pasó y qué hacer, por ejemplo «Sin conexión con el servidor de casa» o «Aún no hay presupuesto para octubre».

## 7. Qué no cubre este resumen

El detalle de cada campo, las reglas de negocio y los cálculos exactos están en documentos internos que no forman parte de este paquete. Para maquetar no hacen falta: lo necesario está aquí, en `prompt.md`, en `DESIGN.pdf`, `U3-dashboard.pdf` y `U4-ajustes.pdf`.
