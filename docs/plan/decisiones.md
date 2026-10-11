# Decisiones cerradas de Axel · Finanzas Hogar (10-oct-2026)

Fuente: 17 rondas de preguntas en la Etapa A (registro completo en `plan.md` §2) más la aprobación del plan. **Ningún agente reabre estas decisiones.** Lo que no esté aquí ni en el plan se trata como bloqueo (ver `protocolo-agentes.md` §2).

## A. Personas y acceso
- Usuarios con cuenta: **Axel** (titular) y **Ximena** (esposa, **co-titular**, mismos permisos). Login con usuario y contraseña, sesión larga por dispositivo (meses), cierre de sesiones desde Ajustes, hash fuerte (argon2id), límite de intentos (5 en 15 min).
- **Karen** (madre de Axel): **participante externa**, sin cuenta ni acceso. Entra solo en (a) la despensa en común y (b) sus compras particulares hechas con el crédito Kueski. **Paga directo a Kueski su parte**: su parte es informativa para el flujo de efectivo de Axel y repone la línea cuando ella paga. No se modelan las finanzas de su familia. Recibe un **estado de cuenta PDF** que el sistema genera y Axel reenvía a mano.
- Plan de comidas y salud: solo Axel y Ximena (2 adultos).
- Avisos de Telegram: Axel y Ximena (ella registra su chat en el bot MultiX una vez; destinatarios configurables por aviso en Ajustes).

## B. Dinero, periodos y reparto
- Moneda **MXN**, símbolo `$`, dos decimales, centavos exactos (ROUND_HALF_UP; el último pago de una serie ajusta centavos). Zona `America/Mexico_City`. Fechas en pantalla `10-oct-2026` (con día de la semana donde importe: `Vie 30 oct 2026`); en archivos y API ISO `2026-10-10`. Idioma **es-MX** únicamente.
- Periodo base del presupuesto: **mes calendario**; la **quincena (15 y 30, o último día del mes)** es subperiodo; la despensa bimestral se **devenga** aparte (mitad por mes), como en el modelo v3.
- Ingresos: **solo los de Axel**, **quincenales (15 y 30)**, capturados por él; más **ingresos extra** manuales. Ximena no aporta ingresos al modelo.
- Reparto del gasto conjunto: **porcentaje configurable por gasto o categoría**, **50/50 por defecto**. Qué se comparte con Karen: **configurable por gasto** (cada gasto dice si es del hogar Axel-Ximena, de Karen, o conjunto y en qué %). La despensa base es conjunto 50/50 por defecto; el crédito se reparte según lo que cada quien compró.
- **Ahorro:** fondo de ahorro y **metas con fecha** (aportación mensual planificada, avance en el dashboard).
- Comidas fuera de casa: **gasto normal** en categoría «comida fuera», con persona; el plan de comidas marca el tiempo como «fuera» sin receta ni calorías.

## C. Créditos
- Tipos a modelar: **línea revolvente tipo Kueski** (límite, tasa por compra, repone capital), **tarjeta bancaria** (corte, pago, mínimo, meses sin intereses), **préstamo personal a plazo fijo** (amortización), **crédito de tienda con abonos**. Axel tiene varios de distintos tipos.
- Kueski: **no hay tabla pública de tasas por plazo** (verificado en fuentes oficiales). La tasa se captura **por compra** (Kueski cotiza total e intereses al pagar; los intereses ya incluyen IVA) y además existe una **tabla editable por plazo** con dos puntos conocidos: 2 quincenas = 17.27 %, 4 quincenas = 25.8685 %. Promociones 0 % se capturan a mano en un **calendario de promociones** (comercio, plazo, compra mínima, vigencia).
- Reposición de capital pendiente de confirmar con Kueski: **interruptor A/B en Ajustes, por defecto B** (B: $10,000 es el límite total; A: $10,000 libres hoy y los pagos suman encima). Aviso visible «pendiente confirmar».
- Optimizador «mayor aprovechamiento»: **comparador de escenarios**. Para una compra planeada genera escenarios (crédito × plazo × fecha ± días × promoción) que respetan las **restricciones duras**: colchón mínimo de disponible, máximo de compras vivas a la vez, deuda en cero antes de una fecha. Evalúa **4 objetivos**: menores intereses, pago por quincena parejo y bajo tope, máximo disponible, comprar más por peso (tienda más barata y promociones). Muestra los 4 óptimos lado a lado y una recomendación con puntaje ponderado por defecto (pesos iguales normalizados). El tope de pago por persona **no** es restricción dura. Sin solver externo: búsqueda exhaustiva.
- Regresión obligatoria: Kueski-Despensa reproduce las cifras de `oraculo/model.json` (ver `plan.md` §1.2).

## D. Comida, salud e inventario
- Recetas: **recetario inicial escrito por Claude** (cocina casera mexicana, ingredientes exactos con presentación comprable, porciones, tiempos) **+ las de Axel y Ximena**, editable. Sin restricciones alimentarias (campo por persona existe, vacío).
- Tiempos: **desayuno, comida, cena y colaciones** (cuarto tiempo opcional con costo y calorías). Porciones base para 2 con escalado.
- Salud: perfil por persona (altura, peso con historial, edad, sexo, actividad), **IMC**, **meta de calorías por día** y **calorías por receta y por día** contra la meta. Macronutrientes quedan para después (modelo de datos previsto).
- Datos nutricionales autorizados: USDA FoodData Central (CC0; carga local del CSV Foundation Foods 04/2026), Open Food Facts (ODbL; volcado local filtrado), tablas mexicanas (captura manual citando fuente), etiquetas de productos de las tiendas cruzadas en internet (captura manual asistida).
- Inventario **híbrido**: stock por unidad para básicos y no perecederos; ritmo de consumo para el resto; sobrantes de recetas entran como stock. Dos verticales de compra: verduras y carne no especializada en **mercados o tianguis** (varios; Axel los da de alta en Ajustes con nombre, día de visita y qué se compra ahí); el resto en Walmart, Bodega Aurrera y Sam's Club.

## E. Precios y tiendas
- Fuentes autorizadas: **captura manual desde ticket**, **foto del ticket con OCR local** (Tesseract 5.5.3 spa, siempre con confirmación humana), **lectura asistida de los sitios de las tiendas** (Obscura o Claude in Chrome, baja frecuencia, **desactivada hasta tener veredicto escrito de términos por sitio** en `docs/estado/S5-terminos.md`), **PROFECO Quién es Quién en los Precios** (datos abiertos mensuales; metadatos y diccionario en `datos/semilla/profeco/`).
- Frescura: un precio se marca viejo a los **30 días** y pesa menos en la comparativa.
- Sam's Club: **no tienen membresía; evaluarla**: comparar como si la tuvieran y calcular el punto de equilibrio anual contra su costo.
- Equivalencia por unidad base (g, ml, pieza); sustitutos marcados; mejor tienda por ingrediente y por canasta.

## F. Reportes y avisos
- Reportes por defecto: **resumen diario** (Telegram, noche), **resumen semanal** (Telegram, domingo), **reporte mensual PDF en LaTeX** (día 1, enviado por Telegram y guardado), **vista web mensual**. Estado de cuenta de Karen en PDF bajo demanda.
- Alertas: **pago próximo** (día anterior y el mismo día), **presupuesto al 80 % y 100 %** (categoría y total), **precios viejos antes del día de despensa**, **plan de comidas sin armar** (domingo).
- Canal: bot MultiX (`~/.openclaw/openclaw.json`, token solo en runtime), texto plano sin parse_mode. Avisos de desarrollo del orquestador con prefijo `[Finanzas-Hogar]`.

## G. Interfaz
- **Fiel a Notion, aplicación en el navegador**: barra lateral colapsable con jerarquía, bloques y propiedades, vistas de base de datos (tabla, tablero, calendario, galería), menú `/`, atajos, edición en el lugar, claro y oscuro. **Notion es la única referencia visual.**
- Dispositivos: **Android** (Axel y Ximena) y escritorio Linux. **PWA instalable** con HTTPS por certificado propio (mkcert). Móvil primero: registrar gasto en ≤ 3 toques.
- Tema por defecto **según el sistema**, con interruptor. Saturación **apagada**, densidad **cómoda**.
- **Paleta nueva completa** (no se parte del reporte). Matiz base: **ver sección I**.
- Dashboard minimalista: Resumen (≤ 6 tarjetas) / Detalle / Tendencias / Calendario.
- Ajustes, **14 secciones**: Cuenta y roles · Presupuesto · Ingresos · Categorías · Créditos y condiciones (incluye promociones) · Tiendas y precios (incluye mercados y membresía Sam's) · Recetas · Reportes y frecuencias · Notificaciones · Apariencia · Datos y respaldos · **Personas y reparto** · **Salud y metas de ahorro** · **Inventario**.
- Tipografía e iconos **autoalojados**; cifras tabulares para dinero; WCAG 2.2 AA mínimo (texto ≥ 4.5:1, UI ≥ 3:1, calculado).

## H. Infraestructura y operación
- Corre en **la máquina de casa** (Arch, LAN 192.168.50.x, mDNS activo), **solo red local**. Dirección propuesta `https://finanzas.local:8443`.
- **Go** + **SQLite (un archivo)** + **servicio systemd nativo** (usuario propio, arranca con la máquina). Respaldos **solo locales** (copia diaria versionada en otra carpeta o disco, `integrity_check`, restauración probada).
- IA opcional: **modelo local (Ollama) y API de Claude**, ambos **desactivados por defecto** y desactivables en Ajustes, solo para OCR de tickets y sugerencias de recetas. El núcleo es determinista.
- Instalar paquetes del sistema: **autorizado sin lista previa**; se reporta lo instalado.
- Datos semilla: compra del 10-oct (`datos/semilla/data.json`), crédito Kueski con 4 despensas y pagos (`oraculo/model.json`), y `datos/semilla/gastos_fijos_ingresos.csv` que entrega Axel (plantilla en `plan.md` §8.2; si no existe, el importador queda listo). Sin históricos.
- Repositorio `~/finanzas-hogar`, git local sin remoto. Sin fecha límite: calidad primero. Sin tope total de tokens; tope de 300K por agente.

## I. Pila y matiz (aprobación del plan)
- **Pila A aprobada:** Go 1.27 + `net/http` + modernc.org/sqlite + goose; React 19 + Vite 8 + TypeScript + Tailwind 4 + shadcn/ui 4 (Radix) + TanStack Table 9 + dnd-kit 6 + ECharts 6 + vite-plugin-pwa 2; **BlockNote solo para notas** (aislado tras interfaz propia); Style Dictionary 5 + culori 4; Storybook 10; Playwright 1.64 + axe-core; mkcert; Tesseract; LaTeX (lualatex) local. Oráculo verificador en Python (`oraculo/`). Sin HiGHS ni OR-Tools.
- **Matiz base:** PENDIENTE DE RESPUESTA DE AXEL (candidatos A 155°, B 195°, C 280° en `plan.md` §4.3). Hasta que se registre aquí, U2a no fija escalas de marca.

## J. Puertas de aprobación en la Etapa B
1. Tras WF-0: C1 + U1 + criterios. 2. Tras WF-1: DESIGN.md + C2 + motores. 3. Antes de Pista 3: Axel inicia sesión en claude.ai/design y comparte `docs/design/entrada/`. 4. Tras WF-3: C3 + U5. Entre puertas no se pregunta salvo bloqueos.
