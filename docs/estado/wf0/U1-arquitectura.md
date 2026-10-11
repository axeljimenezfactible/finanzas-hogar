# Reporte WF-0 · U1-arquitectura

Fase U1 «Arquitectura de información» · Flujo WF-0 «Cimientos» · 10-oct-2026 · Estado: **terminado** (sin cierre por tope).

## Hecho

- `docs/design/U1-arquitectura.md` completo, con las seis partes pedidas:
  - (a) Mapa del sitio y barra lateral con 9 entradas en el orden de la decisión, sus hijos y las 14 secciones de Ajustes; reglas de la barra (tocar el nombre o ▸, plegado por dispositivo, hijos dinámicos solo en Créditos); tabla de lo que se colapsa en móvil; 7 componentes globales (G-01 a G-07).
  - (b) Inventario de **27 pantallas** (P-00 a P-26) y 14 secciones de Ajustes (P-26.01 a P-26.14): propósito, subsecciones (* = por defecto), vistas de base de datos, componentes clave, acción primaria única, secundarias y estados vacío, carga y error. Además: estados comunes, matriz de vistas para las **31 entidades** de S1-01 (vistas, agrupación, filtros, filtro y orden por defecto, o el motivo de no tener lista) y tabla de edición en el lugar por pantalla.
  - (c) **199 necesidades de datos** (N-001 a N-315; 136 de lectura y 63 de escritura; numeración por bloques con huecos reservados), cada una con definición en palabras (fórmula, periodo, filtros, agrupación, orden) y frecuencia de actualización. Matriz pantalla ↔ necesidades con la columna «Endpoint C1» lista para C1. Tabla de entidades de §3.4.1 contra pantallas y lista de necesidades sin entidad explícita para S1.
  - (d) Los 6 flujos de la tarea (F-1 registrar gasto, F-2 semana de comidas, F-3 comparar canasta, F-4 crédito y calendario, F-5 reporte, F-6 peso) y 2 complementarios que piden los criterios (F-7 comparador de escenarios, F-8 ticket con OCR). Todos llevan pasos numerados, pantalla por paso y conteo de toques. La regla de conteo está escrita, y se demuestra en 6 puntos que F-1 tiene exactamente **3 toques**.
  - (e) Menú «/» con 27 filas de comandos (todas con destino) y 15 atajos sin duplicados. Cada atajo dice si coincide con Notion; los de Notion se verificaron en la ayuda oficial con Obscura. Incluye las reglas que evitan choques con el navegador y con AltGr.
  - (f) Reglas móvil primero, tabla elemento × punto de quiebre (360, 412, 768, 1024, 1440), explicación en palabras y tabla por pantalla con móvil 412 / escritorio 1440 / acción primaria / lo que se muestra bajo demanda.
- Dos bloqueos nuevos con el formato del protocolo: `U1-periodos-semana-quincena` y `U1-captura-fuera-de-casa`.
- Alineación con lo que otros agentes de WF-0 escribieron en paralelo: `docs/criterios/U1.md` y los bloqueos `criterios-*`. Ajusté N-091 y N-092 al hallazgo de `criterios-invariante-deuda`; la diferencia acumulada no se rotula como deuda viva. Las fronteras dudosas quedaron marcadas PROVISIONAL con referencia a su bloqueo (`criterios-frescura-peso`, `criterios-umbrales-ausentes`, `criterios-proyeccion-12m`).
- Verifiqué con el oráculo que el efectivo de las ventanas P1 a P6 coincide al centavo con los meses calendario oct-2026 a mar-2027. Va como observación para S1 y S2.

## Archivos

- `/home/deepsec/finanzas-hogar/docs/design/U1-arquitectura.md` (creado, ~123 KB, 934 líneas no vacías)
- `/home/deepsec/finanzas-hogar/docs/estado/wf0/U1-arquitectura.md` (este reporte)
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/U1-periodos-semana-quincena.md` (creado)
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/U1-captura-fuera-de-casa.md` (creado)
- Fuera del repositorio, solo en el directorio temporal de la sesión: el script `u1_verify.py` (copiado íntegro al final de este reporte) y la página de ayuda de Notion descargada.

## Verificación

Comandos ejecutados y salida resumida:

1. **Oráculo, periodos contra meses calendario**: `python3 -I` sobre `oraculo/model.json`. Efectivo P1 a P6 = 973.75 / 1,947.50 / 2,232.43 / 2,716.36 / 2,517.38 / 2,716.36, igual a la suma de pagos y papel de contado por mes calendario de oct-2026 a mar-2027 (6 de 6 OK). Suma 13,103.78 = canon; devengado 14,362.48 = canon. 30-oct-2026 es viernes; 10-oct-2026 y 10-abr-2027 son sábado.
2. **Atajos de Notion**: `obscura fetch --dump text https://www.notion.com/help/keyboard-shortcuts`, 11.6 KB, documentación oficial. Confirmados: Ctrl+K o Ctrl+P (buscar), Ctrl+[ y Ctrl+] (atrás y adelante), Ctrl+Shift+U (página padre), Ctrl+Shift+L (claro y oscuro), Ctrl+N (página nueva, solo en la app de escritorio), Ctrl+/ (editar filas seleccionadas), Ctrl+Enter, Esc, «/». Ctrl+\ no aparece en esa página, así que se marcó como «propio».
3. **`python3 -I u1_verify.py`** (script al final), resultado «TODO OK» y una línea PENDIENTE:
   - **U1-01**: 9 entradas de primer nivel en el orden exacto (Inicio … Ajustes). OK.
   - **U1-02**: 27 ids de pantalla únicos (P-00 a P-26), todos con nombre y sección; 199 necesidades N-xxx únicas; 0 referencias a N-xxx inexistentes. OK.
   - **U1-03**: PENDIENTE. `docs/contratos/openapi.yaml` y `C1-datos.md` no existían al cerrar (C1 se escribe en paralelo). Toda N-xxx está en la matriz 3.29 y toda pantalla y sección tiene fila; el script cruza N-xxx contra las operaciones de C1 en cuanto existan.
   - **U1-04**: F-1 suma 3 toques al guardar. La regla de conteo está en §4 y la demostración en F-1. OK.
   - **U1-05**: 8 flujos con tabla, todos los pasos numerados y con pantalla existente. OK, con la nota de alcance de «Solicitudes».
   - **U1-06**: la matriz §2.5 cubre las 31 entidades de S1-01: 23 con vista y filtro y orden por defecto, 8 con motivo (periodo, regla_reparto, presentacion, ingrediente_receta, tiempo_comida, perfil_salud, nutricion_ingrediente, destinatario). Revisión manual.
   - **U1-07**: 27 comandos, todos con destino; 15 atajos, 0 duplicados; edición en el lugar declarada por pantalla en §2.6. OK.
   - **U1-08**: 14 secciones de Ajustes con nombre y orden exactos, cada una con sus N-xxx en la matriz. OK.
   - **U1-09**: lista de 26 capacidades más abajo, todas con pantalla y N-xxx.
   - **U1-10**: ninguna pantalla propia de Karen salvo «Estado de cuenta de Karen» y ningún acceso para ella. Sus 56 menciones son reparto (decisión B), partes de compras y pagos con el crédito, participación, estado de cuenta y Personas y reparto. OK.
   - **U1-11**: §6.4 da móvil 412, escritorio 1440, acción primaria única y lo que se muestra bajo demanda para las 27 pantallas. Revisión manual.
   - **U1-12**: sin `diccionario-datos.md` ni C1 todavía. La tabla §3.30 y la lista de necesidades sin entidad quedan listas para el cruce.
   - **U1-13**: 77 nombres visibles revisados (pantallas, secciones, barra, comandos), 0 palabras en inglés. OK.
   - **U1-14**: hay sección Supuestos; dos bloqueos U1-*; lo provisional está marcado (§7 del documento).
   - **Cifras**: 30 apariciones de montos en $ (21 distintos), todas presentes en `model.json` o `data.json`, o derivadas de ahí (+$138.50 y +$100.50), o ejemplos de formato ($99,999.99). 0 fuera. No calculé cifras nuevas.
   - **Fechas**: todas las fechas con día de la semana y en formato dd-mmm-aaaa son válidas.

## Cobertura: decisiones A a G y la pantalla que las expone

| Id | Decisión (resumen) | Pantallas | Cómo se expone |
|---|---|---|---|
| A.1 | Axel y Ximena co-titulares con los mismos permisos; usuario y contraseña; sesión larga por dispositivo; cierre de sesiones en Ajustes; argon2id; 5 intentos en 15 min | P-00, P-26.01 | Entrada y bloqueo con hora (N-010, N-011); sesiones de ambos y su cierre (N-271, N-272); sin diferencias de permisos (§0.5). argon2id no tiene exposición en pantalla (S6) |
| A.2 | Karen participante externa sin cuenta: despensa común y sus compras con Kueski; paga directo; su parte es informativa y repone la línea; no se modelan las finanzas de su familia; estado de cuenta PDF que Axel reenvía | P-01 (Detalle), P-03, P-04, P-05, P-17, P-18, P-19, P-25, P-26.12 | Reparto por gasto y por renglón; partes en compras y pagos; su parte se marca pagada por separado y eso repone la línea (N-210); bloque «Participación de Karen» (N-026); estado de cuenta con «Compartir» (N-260, N-261, F-5). Sin pantalla ni acceso para Karen |
| A.3 | Plan de comidas y salud solo para Axel y Ximena | P-10, P-21, P-22, P-26.13 | Las personas fijas del plan y de salud son Axel y Ximena |
| A.4 | Avisos de Telegram a Axel y Ximena; Ximena registra su chat; destinatarios por aviso en Ajustes | P-26.08, P-26.09, P-23 | Destinatarios por reporte y por aviso (N-296, N-298); estado del registro de chat (N-299); historial (N-247) |
| B.1 | MXN, $, centavos exactos, ROUND_HALF_UP, último pago ajusta; zona America/Mexico_City; fechas `10-oct-2026` y `Vie 30 oct 2026`; ISO en archivos; solo es-MX | Todas | Formatos en §0.4; «último pago ajusta centavos» en N-206 y N-216; fechas con día en calendarios y pagos; 0 nombres en inglés |
| B.2 | Mes calendario como base; quincena 15 y 30 como subperiodo; despensa bimestral devengada aparte | P-01, P-06, P-07 | Mes en curso (N-020, N-080); quincenas (N-021, N-082, N-090); despensa devengada (N-083). Días de cada quincena PROVISIONAL (bloqueo) |
| B.3 | Ingresos solo de Axel, quincenales 15 y 30, capturados por él, más extras | P-08, P-26.03, P-07 | Quincena esperada o capturada y extras (N-100 a N-105); no se pide persona |
| B.4 | Reparto con % por gasto o categoría, 50/50 por defecto; qué se comparte con Karen se fija por gasto; despensa conjunta 50/50; el crédito se reparte según lo que cada quien compró | P-03, P-04, P-05, P-19, P-26.04, P-26.12 | Reparto por defecto por categoría (N-281) y del conjunto (N-309); reparto por gasto y por renglón (N-054, N-061, N-071); partes de la compra según los renglones (N-215) |
| B.5 | Fondo de ahorro y metas con fecha, aportación mensual planificada, avance en el dashboard | P-09, P-01, P-26.13 | Tarjeta «Ahorro» (N-025); metas con avance y estado (N-111); definición en N-312 |
| B.6 | Comer fuera es un gasto normal en «comida fuera» con persona; el plan marca «fuera» sin receta ni kcal | P-03, P-10 | Gasto con persona (N-054); casilla «fuera» (N-124) sin costo ni kcal (N-023, N-121) |
| C.1 | 4 tipos de crédito; Axel tiene varios | P-17, P-18, P-26.05 | Variantes por tipo (N-202 a N-207); alta por tipo (N-283) |
| C.2 | Kueski: tasa por compra (intereses con IVA cotizados) más tabla editable por plazo (2 q = 17.27 %, 4 q = 25.8685 %); promociones 0 % en un calendario | P-19, P-18, P-26.05, P-17, P-20 | Intereses cotizados por compra (N-218); tabla de tasas (N-204, N-286); calendario de promociones (N-287) y franjas en el calendario unificado (N-192); promociones en escenarios (N-226) |
| C.3 | Interruptor A/B, B por defecto, con aviso «pendiente confirmar» | P-26.05, P-01, P-17, P-18 | Interruptor con vista previa del efecto (N-284, N-285); aviso en tarjeta, líneas y detalle (N-022, N-193, N-201) |
| C.4 | Comparador de escenarios: 3 restricciones duras, 4 objetivos, 4 óptimos lado a lado, recomendación ponderada con pesos iguales, tope por persona que no es restricción dura, búsqueda exhaustiva | P-20, P-26.05 | N-225, N-226 y flujo F-7; valores por defecto en N-288 |
| C.5 | Regresión Kueski-Despensa contra `model.json` | P-07, P-14, P-18, P-19 | Las cifras salen del motor; las canon citadas en el documento se comprobaron con el script |
| D.1 | Recetario inicial más las de Axel y Ximena, editable; campo de restricciones por persona vacío | P-11, P-12, P-26.07, P-26.13 | N-135 a N-144 con origen; restricciones vacías en el perfil (N-311) |
| D.2 | Tiempos: desayuno, comida, cena y colaciones; porciones base para 2 con escalado | P-10, P-12 | Cuadrícula de 4 tiempos (N-120); escalado (N-140); Ambos = 2 porciones (N-123) |
| D.3 | Perfil por persona, IMC, meta de kcal, kcal por receta y por día; macronutrientes después | P-21, P-22, P-12, P-10, P-26.13 | N-230 a N-233, N-140, N-121; no hay pantalla de macronutrientes |
| D.4 | Fuentes nutricionales autorizadas (USDA, Open Food Facts, tablas mexicanas manuales, etiquetas) | P-12, P-15, P-26.07 | Fuente por ingrediente (N-140, N-173); bases locales (N-294, N-295) |
| D.5 | Inventario híbrido; sobrantes como stock; mercados dados de alta en Ajustes con día de visita y qué se compra ahí; cadenas Walmart, Bodega Aurrera y Sam's | P-16, P-13, P-26.14, P-26.06 | N-180 a N-184, N-314; lista agrupada por vertical con día de visita (N-150); mercados (N-290) |
| E.1 | Fuentes de precio: manual, OCR con confirmación humana, lectura asistida desactivada hasta veredicto, PROFECO | P-05, P-15, P-26.06 | OCR con confirmación (N-070 a N-076, F-8); captura manual (N-174); fuentes y veredicto por sitio (N-291, N-292); importación PROFECO (N-293) |
| E.2 | Precio viejo a los 30 días, que pesa menos | P-15, P-14, P-13 | Marca de viejo (N-170, N-171), peso menor en la comparativa (N-161), frescura en la lista (N-150). Frontera PROVISIONAL |
| E.3 | Sam's: comparar como si tuvieran membresía y calcular el punto de equilibrio | P-14, P-26.06 | N-164; costo anual en N-290 |
| E.4 | Unidad base, sustitutos, mejor tienda por ingrediente y por canasta | P-14, P-15, P-13 | N-161 a N-163, N-170, N-150 |
| F.1 | Resumen diario y semanal por Telegram, mensual PDF el día 1, vista web mensual, estado de cuenta de Karen bajo demanda | P-23, P-24, P-25, P-26.08 | N-245 a N-262, N-296 |
| F.2 | Alertas: pago próximo, presupuesto 80 % y 100 %, precios viejos antes de la despensa, plan sin armar | P-26.09, P-23, P-06, P-01, P-10 | Configuración (N-298); historial (N-247); estados 80/100 (N-020, N-084); casillas sin armar (N-023, N-121) |
| F.3 | Bot MultiX, texto plano, token solo en tiempo de ejecución; prefijo de desarrollo | P-26.09, P-23 | Registro del chat sin mostrar el token (N-299); texto plano en el historial (N-247). El prefijo no se expone (es del orquestador) |
| G.1 | Fiel a Notion: barra lateral con jerarquía, bloques y propiedades, 4 vistas, menú «/», atajos, edición en el lugar, claro y oscuro | Todas, G-01 a G-07 | §1, §2.5, §2.6 y §5 del documento |
| G.2 | Android y escritorio Linux; PWA instalable con mkcert; móvil primero; registrar gasto en ≤ 3 toques | Todas, P-26.01 | F-1 con 3 toques; §6; «Este dispositivo» (N-274) |
| G.3 | Tema según el sistema con interruptor; saturación apagada; densidad cómoda | P-26.10, G-01 | Tema (N-301), Ctrl+Shift+L, /tema. Saturación y densidad son valores fijos de U2, sin control en pantalla |
| G.4 | Paleta nueva; matiz C | — | Sin exposición en U1: lo define U2 |
| G.5 | Dashboard minimalista: Resumen (≤ 6 tarjetas), Detalle, Tendencias, Calendario | P-01 y las demás pantallas con subsecciones | Exactamente 6 tarjetas (N-020 a N-025); regla de subsecciones en §2 |
| G.6 | Ajustes con 14 secciones | P-26.01 a P-26.14 | §2.3 y N-270 a N-315 |
| G.7 | Tipografía e iconos autoalojados; cifras tabulares; WCAG 2.2 AA | Todas | Cifras tabulares (§0.4); objetivos táctiles de 24 px (§6.1). Tipografía y contraste: U2 |

### Lista de capacidades (U1-09)

| # | Capacidad | Pantalla | Necesidades |
|---|---|---|---|
| 1 | Dashboard Resumen / Detalle / Tendencias / Calendario | P-01 | N-020 a N-029 |
| 2 | Participación de Karen (despensa común y crédito, sin cuenta) | P-01, P-19, P-25 | N-026, N-215, N-260 |
| 3 | Comparador de escenarios con 4 objetivos y recomendación | P-20 | N-225, N-226 |
| 4 | Calendario de promociones | P-26.05, P-17 | N-287, N-192 |
| 5 | Interruptor A/B con aviso «pendiente confirmar» | P-26.05, P-01, P-17, P-18 | N-284, N-285, N-022, N-193, N-201 |
| 6 | Plan de comidas con 4 tiempos y «fuera» | P-10 | N-120, N-123, N-124 |
| 7 | Lista de compras | P-13 | N-150 a N-153 |
| 8 | Recetas editables | P-11, P-12 | N-135, N-136, N-142 |
| 9 | Perfil de salud, IMC y kcal | P-21, P-22 | N-230 a N-233, N-240 |
| 10 | Precios con frescura | P-15 | N-170, N-171 |
| 11 | Membresía de Sam's (punto de equilibrio) | P-14 | N-164 |
| 12 | Mercados con día de visita | P-26.06, P-13 | N-290, N-150 |
| 13 | Captura de ticket con OCR y confirmación | P-05 | N-070 a N-076 |
| 14 | Estado de cuenta de Karen en PDF | P-25 | N-260, N-261 |
| 15 | Reporte mensual PDF y vista web mensual | P-23, P-24 | N-248, N-256, N-257 |
| 16 | Metas de ahorro | P-09, P-26.13 | N-111, N-312 |
| 17 | Ingresos quincenales y extras | P-08 | N-100 a N-105 |
| 18 | Tema claro, oscuro o del sistema | P-26.10 | N-301 |
| 19 | Sesiones por dispositivo | P-26.01 | N-271, N-272 |
| 20 | Inventario híbrido y sobrantes | P-16, P-26.14 | N-180, N-184, N-314 |
| 21 | Créditos de 4 tipos con calendario unificado | P-17, P-18 | N-192, N-202 a N-207 |
| 22 | Presupuesto mensual con alertas 80 % y 100 % | P-06 | N-081, N-084 |
| 23 | Flujo por quincena y proyección de 12 meses | P-07 | N-090, N-091 |
| 24 | Respaldos con integridad y restauración | P-26.11 | N-302, N-303 |
| 25 | IA opcional desactivable | P-26.11 | N-306, N-307 |
| 26 | Notificaciones con destinatarios por aviso | P-26.09 | N-298, N-299, N-300 |

## Supuestos

Ninguno cambia cifras, esquema ni tokens. Los que reducen o fijan la forma de C1 van marcados PROVISIONAL para revisarse en la puerta 1.

1. **Conteo de toques**: un toque es una pulsación en un control; teclear valores (monto, peso), desplazar y los diálogos del sistema no cuentan. Coincide con la opción (a) del bloqueo `criterios-toques`. Con la lectura estricta (cada tecla cuenta), F-1 suma 3 toques más las teclas del monto.
2. **Preferencias por dispositivo** (**PROVISIONAL**): la vista elegida, filtros, orden, tema, barra lateral plegada y recientes viven en el dispositivo y no se comparten entre usuarios. Si Axel prefiere vistas compartidas como en Notion, es un cambio que se suma a C1 (v1.1).
3. **Barra inferior en móvil** con 5 destinos (Inicio, Comidas, «+», Compras, Buscar): adaptación móvil de U1; Notion sigue siendo la única referencia visual.
4. **Semana por defecto**: los domingos, el plan de comidas y la lista de compras abren la semana siguiente, en línea con el aviso del domingo de «plan sin armar».
5. **Cuenta sugerida por renglón en el OCR**: la que tuvo ese producto en la última compra confirmada; si no hay, el reparto de la categoría. Siempre se confirma a mano.

## Bloqueos

- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/U1-periodos-semana-quincena.md`: inicio de la semana y días de cada quincena (incluido el 31). En U1 quedó PROVISIONAL: lunes a domingo, y quincenas 1–15 y 16–fin de mes.
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/U1-captura-fuera-de-casa.md`: si la PWA guarda en cola lo capturado fuera de casa. En U1 quedó PROVISIONAL: capturas deshabilitadas sin conexión.
- Referenciados, de otros agentes: `criterios-frescura-peso`, `criterios-umbrales-ausentes`, `criterios-invariante-deuda`, `criterios-proyeccion-12m`, `criterios-toques` (U1 escribe la regla), `criterios-nxx` (U1 lo resuelve; ver Solicitudes).

## Pendientes

Nada de lo asignado a U1 queda pendiente. Dependencias de otros agentes:

1. C1: mapear las 199 necesidades (N-001 a N-315) a endpoints y re-ejecutar `u1_verify.py` para cerrar U1-03. N-003 y N-301 son locales y no requieren endpoint.
2. U3: especificar las 6 tarjetas a partir de N-020 a N-026.
3. U4: campos y validaciones de P-26.01 a P-26.14; valores por defecto no decididos (hora del resumen diario, anticipación de alertas, holgura del comparador, precisión del peso).
4. U5: verificar en Android a 412 × 915 que el teclado numérico abre con el toque 1 de F-1 y que la hoja cabe sin desplazar a 360 px.
5. Axel: responder los dos bloqueos U1.

## Solicitudes fuera de alcance

1. **`docs/criterios/U1.md`** (agente de criterios o integrador):
   - (a) U1-02, U1-04, U1-05 y U1-08 usan «N-xx» como id de pantalla. La tarea de U1 define N-xx como **necesidad de datos** («tabla de NECESIDADES DE DATOS con id N-xx … el contrato C1 mapeará cada N-xx a endpoints»), así que U1 usa **P-xx para pantallas** y **N-xxx para necesidades**. Esto resuelve `criterios-nxx` con su opción (b). Propuesta: leer «pantalla N-xx» como «pantalla P-xx» y «N-xx con datos» como N-xxx en U1-03.
   - (b) U1-05 nombra como obligatorios el comparador y el ticket con OCR, pero la tarea de U1 fija otros 6 flujos. U1 entrega los 6 de la tarea y agrega F-7 y F-8. Propuesta: aceptar 6 + 2 complementarios.
2. **`docs/estado/bloqueos/criterios-toques.md`** (integrador): la regla de conteo ya está escrita en U1 §4 (opción a); se puede cerrar.
3. **S1** (`diccionario-datos.md`, migraciones): necesidades sin entidad explícita en §3.4.1. Son: renglones de ticket con su cuenta (N-061, N-074); lista de compras con marcas de comprado y artículos manuales (N-150 a N-153); lecturas OCR sin confirmar (N-075, N-076); gastos fijos recurrentes (N-276); despensa programada (N-277); estado de pago por parte, hogar y Karen (N-210); destino de un movimiento de ahorro (N-112, N-113); reporte generado (N-245), si `job` no alcanza.
4. **S1 y S2** (`reglas-negocio.md`): el efectivo de P1 a P6 de `model.json` coincide al centavo con oct-2026 a mar-2027. Con un devengo que asigne cada despensa al mes de compra y al siguiente, la proyección por mes calendario reproduce $14,362.48 y $13,103.78. Además, usar la misma respuesta del bloqueo `U1-periodos-semana-quincena`.
5. **C1** (`C1-datos.md`): escribir el mapeo con una fila por necesidad que contenga el id y la operación (por ejemplo «N-020 … GET /api/…»), para que el script lo lea. P-26.10 no requiere endpoint.
6. **S6**: las sesiones de ambos usuarios son visibles y cualquiera de los dos las puede cerrar (co-titulares, N-271); el certificado raíz público de la CA local se descarga desde N-274; el comportamiento del bloqueo sigue `criterios-login-bloqueo`.
7. **Integrador**: si se quiere re-ejecutar la verificación sin copiar el bloque de abajo, guardar el script en una ruta del repositorio (por ejemplo `docs/design/verificar_u1.py`).

## Paquetes instalados

Ninguno.

## Consumo aproximado

- Unas 47 acciones de herramienta.
- Leído, unos 75 KB: protocolo (4.6 KB), decisiones (11 KB), secciones del plan §0–§2, §3.4, §3.5, §3.7 (~22 KB), `docs/criterios/U1.md` (11.6 KB), 10 bloqueos de criterios (~11 KB), extracto de S1-01, ayuda de Notion (11.6 KB); `model.json` y `data.json` solo mediante scripts.
- Escrito, unos 150 KB: documento de diseño (~123 KB), este reporte (~24 KB), dos bloqueos (4.3 KB).

## Anexo: script de verificación `u1_verify.py`

Uso: guardarlo y ejecutar `python3 -I u1_verify.py /home/deepsec/finanzas-hogar`. Salida de la última ejecución:

```
OK    U1-01 barra lateral: 9 entradas de primer nivel en orden: ['Inicio', 'Gastos', 'Presupuesto', 'Comidas', 'Compras y tiendas', 'Créditos', 'Salud', 'Reportes', 'Ajustes']
OK    U1-02 ids de pantalla únicos: 27 (P-00..P-26)
OK    U1-02 cada pantalla con nombre y sección
OK    U1-08 Ajustes: 14 secciones con nombre y orden exactos
OK    N-xxx definidas y únicas: 199 (L=136, E=63)
OK    Referencias a N-xxx inexistentes: []
OK    Toda N-xxx aparece en la matriz 3.29; faltan: []
OK    Toda pantalla y sección de Ajustes está en la matriz; faltan: []
OK    F-1: 4 pasos numerados, todos con pantalla existente (sin pantalla: [], inexistentes: [])
OK    U1-04 F-1 registrar gasto: toques acumulados al guardar = **3**
OK    F-2: 8 pasos numerados, todos con pantalla existente (sin pantalla: [], inexistentes: [])
OK    F-3: 6 pasos numerados, todos con pantalla existente (sin pantalla: [], inexistentes: [])
OK    F-4: 4 pasos numerados, todos con pantalla existente (sin pantalla: [], inexistentes: [])
OK    F-5: 5 pasos numerados, todos con pantalla existente (sin pantalla: [], inexistentes: [])
OK    F-6: 4 pasos numerados, todos con pantalla existente (sin pantalla: [], inexistentes: [])
OK    F-7: 7 pasos numerados, todos con pantalla existente (sin pantalla: [], inexistentes: [])
OK    F-8: 7 pasos numerados, todos con pantalla existente (sin pantalla: [], inexistentes: [])
OK    U1-05 flujos con tabla: 8 (6 de la tarea + 2 complementarios)
OK    U1-07 comandos «/»: 27 filas, sin destino: []
OK    U1-07 atajos: 15 filas, duplicados: []
OK    Cifras en $ (30 apariciones, 21 distintas) presentes en model.json/data.json o ejemplos de formato; fuera: []
OK    Fechas con día de la semana y fechas dd-mmm-aaaa válidas; erróneas: []
OK    U1-13 nombres visibles sin palabras en inglés (77 revisados); coincidencias: []
OK    U1-10 ninguna pantalla propia de Karen salvo «Estado de cuenta de Karen»
Menciones de «Karen» en el documento: 56
PENDIENTE U1-03: docs/contratos/openapi.yaml o C1-datos.md aún no existen (C1 se escribe en paralelo en WF-0); re-ejecutar cuando existan

RESULTADO: TODO OK
```

```python
"""Verificación de U1 (docs/design/U1-arquitectura.md). Uso: python3 -I u1_verify.py [raiz_repo]"""
import re, json, sys, os, datetime as dt
ROOT = sys.argv[1] if len(sys.argv) > 1 else '/home/deepsec/finanzas-hogar'
doc = open(os.path.join(ROOT, 'docs/design/U1-arquitectura.md'), encoding='utf-8').read()
L = doc.split('\n')
fallas = []
def ok(c, msg):
    print(('OK   ' if c else 'FALLA') + ' ' + msg)
    if not c: fallas.append(msg)
def section(start, end):
    a = doc.index(start); b = doc.index(end, a + len(start)); return doc[a:b]
def rows(text):
    return [[c.strip() for c in l.strip().strip('|').split('|')] for l in text.split('\n') if l.startswith('| ') and not l.startswith('|---')]

# U1-01 barra lateral
tree = section('```\nFinanzas Hogar', '\n```')
top = []
for l in tree.split('\n'):
    if l.startswith('├─ ') or l.startswith('└─ '):
        t = re.split(r'\s{2,}|→', l[3:])[0].replace('▸', '').strip(); top.append(t)
secc = [t for t in top if t not in ('Buscar', 'Registrar gasto')]
esperadas = ['Inicio', 'Gastos', 'Presupuesto', 'Comidas', 'Compras y tiendas', 'Créditos', 'Salud', 'Reportes', 'Ajustes']
ok(secc == esperadas, f'U1-01 barra lateral: {len(secc)} entradas de primer nivel en orden: {secc}')

# U1-02 inventario de pantallas
inv = rows(section('### 2.1 Qué es cada pantalla', '### 2.2'))[1:]
pids = [r[0] for r in inv]
ok(len(pids) == len(set(pids)), f'U1-02 ids de pantalla únicos: {len(pids)} ({pids[0]}..{pids[-1]})')
ok(all(len(r) >= 3 and r[1] and r[2] for r in inv), 'U1-02 cada pantalla con nombre y sección')
ajs = rows(section('### 2.3 Ajustes: las 14 secciones', '### 2.4'))[1:]
nombres14 = ['Cuenta y roles', 'Presupuesto', 'Ingresos', 'Categorías', 'Créditos y condiciones (incluye promociones)', "Tiendas y precios (incluye mercados y membresía Sam's)", 'Recetas', 'Reportes y frecuencias', 'Notificaciones', 'Apariencia', 'Datos y respaldos', 'Personas y reparto', 'Salud y metas de ahorro', 'Inventario']
ok([r[1] for r in ajs] == nombres14 and [r[0] for r in ajs] == [f'P-26.{i:02d}' for i in range(1, 15)], f'U1-08 Ajustes: {len(ajs)} secciones con nombre y orden exactos')
todas_p = set(pids) | {r[0] for r in ajs}

# Necesidades N-xxx
nec = section('## 3. Necesidades de datos', '### 3.29')
defs = re.findall(r'^\| (N-\d{3}) \| ([LE]) \|', nec, re.M)
nids = [d[0] for d in defs]
ok(len(nids) == len(set(nids)), f'N-xxx definidas y únicas: {len(nids)} (L={sum(1 for d in defs if d[1]=="L")}, E={sum(1 for d in defs if d[1]=="E")})')
nset = set(nids)
refs = set(re.findall(r'N-\d{3}', doc))
ok(refs <= nset, f'Referencias a N-xxx inexistentes: {sorted(refs - nset)}')
# Matriz 3.29
mat = rows(section('### 3.29', '### 3.30'))[1:]
def expand(cell):
    out = set()
    for a, b in re.findall(r'(N-\d{3}) a (N-\d{3})', cell):
        out |= {n for n in nset if a <= n <= b}
    return out | set(re.findall(r'N-\d{3}', cell))
en_matriz = set(); pant_matriz = set()
for r in mat:
    pant_matriz |= set(re.findall(r'P-\d{2}(?:\.\d{2})?', r[0]))
    en_matriz |= expand(r[1]) | expand(r[2])
ok(nset <= en_matriz, f'Toda N-xxx aparece en la matriz 3.29; faltan: {sorted(nset - en_matriz)}')
falt_p = sorted((todas_p - {'P-26'}) - pant_matriz)
ok(not falt_p, f'Toda pantalla y sección de Ajustes está en la matriz; faltan: {falt_p}')

# Flujos
fl = section('## 4. Flujos clave', '### Resumen de toques')
bloques = re.split(r'\n(?=### F-\d|\*\*F-\d)', fl)
ids_g = set(re.findall(r'^\| (G-\d{2}) \|', doc, re.M))
nflujos = 0
for b in bloques:
    m = re.match(r'(?:### |\*\*)(F-\d)', b)
    if not m: continue
    nflujos += 1
    pasos = [r for r in rows(b) if r and r[0].isdigit()]
    sinp = [r[0] for r in pasos if not re.search(r'[PG]-\d{2}', r[1])]
    malos = [x for r in pasos for x in re.findall(r'[PG]-\d{2}(?:\.\d{2})?', r[1]) if x not in todas_p | ids_g]
    ok(pasos and not sinp and not malos, f'{m.group(1)}: {len(pasos)} pasos numerados, todos con pantalla existente (sin pantalla: {sinp}, inexistentes: {malos})')
    if m.group(1) == 'F-1':
        ultimo = pasos[-1][-1]
        ok(ultimo.strip('*') == '3', f'U1-04 F-1 registrar gasto: toques acumulados al guardar = {ultimo}')
ok(nflujos == 8, f'U1-05 flujos con tabla: {nflujos} (6 de la tarea + 2 complementarios)')

# Menú / y atajos
cmd = rows(section('### 5.2 Lista de comandos', '### 5.3'))[1:]
sin_dest = [r[0] for r in cmd if not r[3] or r[3] == '—']
ok(not sin_dest, f'U1-07 comandos «/»: {len(cmd)} filas, sin destino: {sin_dest}')
ata = rows(section('### 5.3 Atajos', '### 5.4'))[1:]
claves = [r[0] for r in ata]
dup = sorted({c for c in claves if claves.count(c) > 1})
ok(not dup, f'U1-07 atajos: {len(claves)} filas, duplicados: {dup}')

# Cifras contra el oráculo y la semilla
def nums(o, acc):
    if isinstance(o, dict): [nums(v, acc) for v in o.values()]
    elif isinstance(o, list): [nums(v, acc) for v in o]
    elif isinstance(o, (int, float)) and not isinstance(o, bool): acc.add(round(float(o), 2))
    return acc
m = json.load(open(os.path.join(ROOT, 'oraculo/model.json')))
d = json.load(open(os.path.join(ROOT, 'datos/semilla/data.json')))
canon = nums(m, set()) | nums(d, set())
t = d['totales']
canon |= {round(t['walmart'] - t['bodega'], 2), round(t['walmart_identicos'] - t['bodega_identicos'], 2)}
formato = {99999.99, 0.0}
cifras = re.findall(r'\$(\d{1,3}(?:,\d{3})*\.\d{2})', doc)
fuera = sorted({c for c in cifras if round(float(c.replace(',', '')), 2) not in canon | formato})
ok(not fuera, f'Cifras en $ ({len(cifras)} apariciones, {len(set(cifras))} distintas) presentes en model.json/data.json o ejemplos de formato; fuera: {fuera}')

# Fechas
mes = dict(ene=1, feb=2, mar=3, abr=4, may=5, jun=6, jul=7, ago=8, sep=9, oct=10, nov=11, dic=12)
dsem = ['Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb', 'Dom']
malas = []
for wd, dd, mm, yy in re.findall(r'\b(Lun|Mar|Mié|Jue|Vie|Sáb|Dom) (\d{1,2}) (ene|feb|mar|abr|may|jun|jul|ago|sep|oct|nov|dic) (\d{4})', doc):
    if dsem[dt.date(int(yy), mes[mm], int(dd)).weekday()] != wd: malas.append(f'{wd} {dd} {mm} {yy}')
for dd, mm, yy in re.findall(r'\b(\d{1,2})-(ene|feb|mar|abr|may|jun|jul|ago|sep|oct|nov|dic)-(\d{4})', doc):
    try: dt.date(int(yy), mes[mm], int(dd))
    except ValueError: malas.append(f'{dd}-{mm}-{yy}')
ok(not malas, f'Fechas con día de la semana y fechas dd-mmm-aaaa válidas; erróneas: {malas}')

# Idioma de nombres visibles
visibles = [r[1] for r in inv] + [r[1] for r in ajs] + secc + [r[2] for r in cmd]
ingles = re.compile(r'\b(Dashboard|Settings|Home|Save|Board|Gallery|Table|Calendar|Search|Login|Logout|Report|Budget|Credit|Shopping|Inventory|Health|Expense|Income|Recipe|Store)\b', re.I)
hits = sorted({v for v in visibles if ingles.search(v)})
ok(not hits, f'U1-13 nombres visibles sin palabras en inglés ({len(visibles)} revisados); coincidencias: {hits}')

# Karen
ok(not any('Karen' in r[1] and 'Estado de cuenta' not in r[1] for r in inv), 'U1-10 ninguna pantalla propia de Karen salvo «Estado de cuenta de Karen»')
print('Menciones de «Karen» en el documento:', doc.count('Karen'))

# Cruce con C1 (U1-03), si ya existe
oa = os.path.join(ROOT, 'docs/contratos/openapi.yaml'); c1 = os.path.join(ROOT, 'docs/contratos/C1-datos.md')
if os.path.exists(oa) and os.path.exists(c1):
    y = open(oa, encoding='utf-8').read(); ops = set(); cur = None
    for l in y.split('\n'):
        mp = re.match(r'^  (/\S*):\s*$', l)
        if mp: cur = mp.group(1); continue
        mm_ = re.match(r'^    (get|post|put|patch|delete):', l)
        if mm_ and cur: ops.add(f'{mm_.group(1).upper()} {cur}')
    c1t = open(c1, encoding='utf-8').read(); mapa = {}
    for l in c1t.split('\n'):
        ns = re.findall(r'N-\d{3}', l); os_ = re.findall(r'\b(GET|POST|PUT|PATCH|DELETE)\s+`?(/[^\s`|]*)', l)
        for n in ns: mapa.setdefault(n, set()).update(f'{a} {b}' for a, b in os_)
    locales = {'N-003', 'N-301'}
    sin_ep = sorted(n for n in nset - locales if not mapa.get(n))
    usados = set().union(*mapa.values()) if mapa else set()
    sin_n = sorted(ops - usados)
    ok(not sin_ep, f'U1-03 N-xxx sin endpoint en C1: {sin_ep}')
    ok(not sin_n, f'U1-03 operaciones de openapi.yaml sin N-xxx: {sin_n}')
else:
    print('PENDIENTE U1-03: docs/contratos/openapi.yaml o C1-datos.md aún no existen (C1 se escribe en paralelo en WF-0); re-ejecutar cuando existan')
print('\nRESULTADO:', 'TODO OK' if not fallas else f'{len(fallas)} falla(s)')
```
