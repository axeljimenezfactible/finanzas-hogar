# artboards-ajustes-a — reporte

Fecha de trabajo: 12-oct-2026. Agente Sonnet, Pista 3 fase 2 (lienzo), área `ajustes`, parte a.

## Hecho
16 artboards (móvil `-m` y escritorio `-e`) en `docs/design/maquetas/lienzo/project/ajustes/`, fragmento `docs/design/maquetas/lienzo/fragmentos/ajustes-a.json` (16 `boards`, 16 `order`; móvil y luego escritorio de cada pantalla):
- `ajustes-p26-indice-{m,e}`: índice con las 14 secciones (en escritorio, la columna interna de 14 más tarjetas con un resumen de cada una sacado del JSON sintético).
- `ajustes-p26-s01-cuenta-roles`, `s02-presupuesto`, `s03-ingresos`, `s04-categorias`, `s05-creditos-condiciones`, `s06-tiendas-precios`, `s07-recetas` en `-m` y `-e`.
- Cada sección sigue U4 P-26.01 a P-26.07: campos con su etiqueta, tablas (en móvil, filas apiladas), interruptores, parámetros K44 con «Valor inicial» y «Restablecer», vistas previas «Antes → Después», diálogo de confirmación (alcance del presupuesto, cierre de sesiones) y avisos (pendiente de confirmar, faltan restricciones, sin veredicto V-19). El escritorio trae la tarjeta «Estados» (vacío, carga, error al cargar) con el texto de U4 de cada sección.
- Datos: todos de `datos-sinteticos.json` (el generador los lee y los copia como texto literal). Barra lateral, botones y tarjetas enlazan por `../<area>/<archivo>.dc.html` al artboard del mismo ancho.
- Los artboards los genera un script propio (`/tmp/ajgen/{lib,sec_a,sec_b,build,verify}.py`, fuera del proyecto); no se escribieron a mano uno por uno.

## Verificación (`python3 -I /tmp/ajgen/verify.py`, resultado: 16 archivos, 0 fallas)
- Líneas obligatorias presentes (doctype, `lang="es-MX"`, `support.js`, `tokens.css`, `x-dc`, `helmet`, `DCLogic`); `data-props` parsea y su `$preview` coincide con el fragmento (412 o 1440; alto igual).
- `html.parser`: 0 etiquetas sin cerrar, 0 ids repetidos; cada `label for` apunta a un id; cada campo tiene `label` o `aria-label`; cada `<use>` tiene su `<symbol>`.
- 0 hex literales fuera de `data-props`; todos los `var(--x)` existen en `tokens.css` (generado de `tokens.json`); 0 cifras prohibidas de P3-12 (con y sin coma, y en centavos); 0 `{{`, iframes, `innerHTML` ni URLs; 0 emoji; 0 negativos con guion.
- Raíz móvil `class="m"` (ancho `var(--breakpoint-phone)` = 412 px, mínimo 915); raíz de escritorio `pg` sin ancho en px. Fragmento: `expand: "fill"` solo en `-e`, `is_interactive: false`, `page: "ajustes"`; lista exactamente mis 16 archivos.
- No se renderizó nada en navegador (K45): las alturas del fragmento son estimadas y a propósito altas.

## Supuestos
1. Ranuras de archivos de la parte b: usé los nombres `s08-reportes-frecuencias`, `s09-notificaciones`, `s10-apariencia`, `s11-datos-respaldos`, `s12-personas-reparto`, `s13-salud-metas`, `s14-inventario` en los enlaces de la columna interna y de «Ir a Datos y respaldos» / «Personas y reparto». Si el agente b usa otros nombres, esos enlaces rompen: revisar al unir.
2. Valores tecleados de ejemplo en las vistas previas (no vienen del JSON): sesión 60 días, Súper de la semana $3,400.00 (total $27,786.00), quincena del 30 de $14,950.00 (esperado $29,210.00). El ejemplo de reparto usa el gasto g13 del JSON ($2,918.40, 50/50) en lugar del «$199.00 del papel», que no existe en los datos.
3. Pesos del comparador: el JSON trae 3 (40, 30, 30); U4 pide 4. Dibujé el cuarto («Comprar más por peso») en 0 y los % normalizados 40/30/30/0.
4. «Deuda en cero antes de» está vacía en el JSON; por U4 K16 dibujé «Faltan restricciones» aunque P-20 muestre una recomendación (incongruencia de datos, no de maqueta).
5. «Falta capturar»: todas las líneas del JSON tienen sus condiciones, así que se dibuja el estado «todo capturado»; el aviso «Falta capturar» no aparece en datos.
6. Bases nutricionales: conteo de productos del catálogo por `nutricion_fuente` (28 USDA, 5 manual) y fecha de carga 12-oct-2026 (la de la semilla); Open Food Facts sin datos. «Clave que usan los cálculos»: `groceries` en Despensa común y `eating_out` en Restaurantes y café (lectura mía de decisión B).
7. Origen y fecha de cotización de las tasas por plazo: «Tecleada» y «—» (el JSON no los trae). «App instalada: Sí» es dato del navegador, no del JSON.
8. Sin barra inferior en móvil (G-05) ni barra superior de escritorio completa: el móvil lleva «‹ Ajustes» y el escritorio, migas y acción primaria.

## Bloqueos
Ninguno.

## Pendientes
- Medir las alturas reales del fragmento (`h`) al publicar y ajustarlas si algún artboard queda corto.
- Revisar con el agente b los nombres de archivo de las secciones 08 a 14 (supuesto 1).
- El calendario de promociones del escritorio es tabla (con el selector Calendario · Tabla · Tablero), el móvil es lista por mes; no hay calendario mensual dibujado. Vista tablero de Categorías y de Tiendas: solo se dibuja el selector.
- `support.js` y `ds/finanzas-hogar/tokens.css` (con `fonts/`) los publica el orquestador.

## Consumo aproximado
- Leído: ~260 KB (U4 P-26.01 a .07 y convenciones, U1 §1.1/§2.3, DESIGN §5.4/§7, format.md, craft.md, LEEME y fragmentos del JSON sintético). Escrito con herramientas: ~62 KB de scripts del generador + este reporte + fragmento (los 16 artboards generados suman ~362 KB en disco, producidos por el script). ~35 acciones de herramienta.
