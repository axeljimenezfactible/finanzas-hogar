# artboards-ajustes-b · Pista 3, fase 2 (lienzo) · P-26.08 a P-26.14

Fecha: 12-oct-2026. Carpeta: `docs/design/maquetas/lienzo/project/ajustes/`. Fragmento: `docs/design/maquetas/lienzo/fragmentos/ajustes-b.json` (14 tableros, `order` por sección, móvil antes que escritorio).

## Hecho (14 artboards, `ajustes-p26-sNN-<slug>-{m,e}.dc.html`)
| Sección | Archivos | Contenido principal |
|---|---|---|
| P-26.08 Reportes y frecuencias | `s08-reportes` | Tabla de 3 reportes (activo, día, hora, destinatarios, próximo envío, último resultado con el error real del diario), vista previa «Antes → Después» del próximo envío, 5 parámetros `reports.*` con valor inicial, «Bajo demanda» (mensual ahora, estado de cuenta de Karen, vista web), Estados |
| P-26.09 Notificaciones | `s09-notificaciones` | 4 alertas (la de precios viejos apagada), umbrales 80/100 con vista previa (`$22,068.80 → $23,448.10`), 6 parámetros `alerts.*`, chats confirmados, chat detectado «Ximena (tablet)» con Confirmar/Descartar y diálogo crítico, aviso de prueba, Estados |
| P-26.10 Apariencia | `s10-apariencia` | Tema segmentado, muestra clara y oscura, atajos de una tecla con vista previa, valores fijos, nota de estados |
| P-26.11 Datos y respaldos | `s11-datos-respaldos` | Respaldos (5), diálogo de restauración paso 2 (`RESTAURAR`), importar CSV, datos semilla (2 cargadas, 1 pendiente), IA opcional apagada, estado del servicio, parámetros (adjunto, idempotencia, OCR) con vista previa, Estados |
| P-26.12 Personas y reparto | `s12-personas-reparto` | Axel/Ximena/Karen, tope por quincena vacío, 3 repartos con barra etiquetada, vista previa con el gasto de $2,918.40, diálogo de alcance (Q2-68 C), regla de resto, Estados |
| P-26.13 Salud y metas de ahorro | `s13-salud-metas` | Fichas de Axel y Ximena (IMC, meta sugerida/confirmada, ajustes K36), ajuste de la casa, metas (fondo + Laptop atrasada + Viaje en ruta), vista previa (3) del total reservado, Estados |
| P-26.14 Inventario | `s14-inventario` | 14 artículos (modo, nivel, ritmo, existencia, estado con icono y texto, días restantes), resumen 1/5/8, vista previa de estado, Estados |

Todas las cifras salen de `datos-sinteticos.json` (el generador las lee del archivo). Las únicas derivadas son aritmética sobre esos datos: 80 % y 85 % del plan de $27,586.00; 40 % de $2,918.40 (gasto g13); $2,000.00 − $500.00; −100 kcal; nivel 300 g de la pasta.

## Verificación (script propio, `python3 -I`, en `/tmp/ajb/verif.py`; 14 archivos, 0 fallas)
Por archivo: doctype, `lang="es-MX"`, charset, `<title>`, `support.js`, `tokens.css`, `x-dc`/`helmet`, bloque `data-dc-script` con `class Component extends DCLogic`; JSON de `data-props` parsea y ancho 412/1440; `html.parser` sin etiquetas sin cerrar; 0 ids repetidos; 0 hex ni `rgb()` fuera de `data-props`; 0 cifras prohibidas P3-12 (con y sin coma y en centavos); todos los `var(--x)` existen en `sistema/project/tokens.css` (generado de `tokens.json`); sin `innerHTML`, `iframe`, red, emoji ni «lorem»; móvil con `width:412px` y `min-height:915px`; escritorio con raíz sin ancho; `tabular-nums`; todo `input`/`select` con `label` o `aria-label`. Fragmento: JSON válido y lista exactamente los 14 archivos.
No se abrió navegador (regla K45): no hay verificación visual.

## Supuestos
1. Alturas del `$preview` y del fragmento: estimadas por conteo de texto y controles (sin navegador). Móviles de 915 a 3,356 px; escritorio de 1,000 a 3,102 px. El orquestador debe re-medirlas.
2. Claves del fragmento con extensión: `ajustes/<archivo>.dc.html`.
3. El JSON no trae fecha de nacimiento: P-26.13 muestra la edad (`edad` en años cumplidos) y deja el campo de fecha vacío con marcador `dd-mmm-aaaa`. Tampoco trae `chat_id_masked`: la terminación del chat se dibuja «•••». Tampoco hay restauración probada, huella, tipo de respaldo ni caducidad: se muestra «—» y «Diario» (todos los respaldos del JSON son diarios).
4. Valores iniciales de `reports.*`, `alerts.*`, `ocr.timeout_seconds` (30 s) y `health.goal_adjust_*` (0) vienen de U4-ajustes-b (Q2-72, Q2-73, K50, K36); el JSON los trae solo en parte. Vistas previas con valor tecleado de ejemplo (hora 20:00, umbral 85 %, adjunto 15 MB, Karen 40 %, ajuste para bajar 500 kcal) son ilustrativas y no están guardadas.
5. Barra lateral de la aplicación con los 9 iconos de `BarraLateral/preview.html`; la columna interna de las 14 secciones va solo con texto (sin iconos) para no pasar el tope de peso. Móvil sin barra inferior: arriba solo «‹ Ajustes» (la lista de 14 secciones es de `ajustes-p26-indice`).
6. Enlaces a `s01` a `s07` apuntan a los archivos de `ajustes-a` que existían al generar (los 7 de escritorio); si cambian de nombre, hay que regenerar. Enlaces a pantallas de otras áreas siguen la lista de la tarea; botones secundarios sin destino van a `#`.
7. Iconos: Lucide en un `<symbol>` por archivo y `<use>`; estado siempre con icono y texto.
8. DESIGN §7 se leyó a través de los README de `sistema/project/components/` y de `contraste-componentes.csv` (pares insignia de rol, campo, tabla); no se leyeron U3 (no aplica a Ajustes) ni todas las filas de U1 §2.5.

## Bloqueos
Ninguno.

## Pendientes
- Revisión visual y re-medición de alturas.
- Ajustes-b quedó sin `renderVals()` con datos (el contenido es texto literal por tratarse de pantallas estáticas); si se quiere editar por tweaks, hay que mover las listas.
- Iconos de la columna interna de secciones y barra inferior móvil, si el revisor los pide.

## Consumo aproximado
Leído ~165 KB (U4-ajustes-b secciones 8 a 14 y §0 de la parte a, format.md, LEEME, previews de componentes, consultas al JSON). Escrito por herramientas ~57 KB (generador `/tmp/ajb/*.py` ~48 KB, verificador y este reporte). Los 14 artboards generados pesan ~224 KB en disco más 5 KB del fragmento: pasan el tope de 150 KB si se cuenta el resultado y no lo escrito a mano; se redujo con poda de CSS. ~40 acciones.
