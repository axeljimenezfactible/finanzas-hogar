# artboards-salud-reportes · Pista 3 fase 2 (lienzo) · 12-oct-2026

## Hecho
8 artboards (móvil 412 px y escritorio PAGE fluida) en `docs/design/maquetas/lienzo/project/salud-reportes/` y fragmento `docs/design/maquetas/lienzo/fragmentos/salud-reportes.json` (8 boards, `order` completo, escritorio con `"expand":"fill"`).

| Pantalla | Móvil | Escritorio |
|---|---|---|
| P-21 Salud | `salud-p21-salud-m` (Axel y Ximena en pestañas; tabs con radios reales, sin JS; `is_interactive: true`) | `salud-p21-salud-e` (dos columnas: perfil, IMC, meta de kcal con ajuste por objetivo N-327 y «Confirmar», kcal de hoy, gráfica de peso, historial, tarjeta «Estados») |
| P-23 Centro de reportes | `reportes-p23-centro-m` (lista de reportes, avisos plegados) | `reportes-p23-centro-e` (generar reporte, próximos envíos 21:30 / dom 19:00 / día 1 08:00 / avisos 09:00, tabla de generados, avisos enviados, «Estados») |
| P-24 Reporte mensual | `reportes-p24-mensual-m` | `reportes-p24-mensual-e` (7 secciones plegables a ancho de lectura, índice lateral, «Estados») |
| P-25 Estado de cuenta de Karen | `reportes-p25-karen-m` | `reportes-p25-karen-e` (vista previa en página, movimientos plegados, estados anteriores, «Estados») |

Enlaces: barra lateral, barra inferior, «Registrar peso» (P-22), P-23 ↔ P-24 ↔ P-25 y ☰ → `shell-barra-lateral-m`. «Editar perfil», «Buscar», subsecciones y filtros van a `#` (el artboard destino no tiene nombre fijo en la lista).

## Verificación (script propio `python3 -I`, /tmp/pista3_sr/verify.py)
- Líneas obligatorias, `<title>`, `lang="es-MX"`, `support.js`, `tokens.css`, `x-dc`, `helmet`, `data-dc-script`: OK en los 8.
- `data-props` parsea; `$preview` 412 o 1440: OK.
- `html.parser`: sin etiquetas sin cerrar ni cierres inesperados: OK.
- 0 hex literales fuera de `data-props`; 0 `rgb()/hsl()`; 0 `{{`; sin iframes, innerHTML ni red; sin emoji.
- 0 cifras prohibidas de P3-12 (con y sin coma).
- Todos los `var(--x)` existen en `tokens.css`/`tokens.json`.
- Cada importe `$` de los archivos existe en el JSON sintético (campo `_fmt` o algún `*_cents` formateado). Sin signo menos ASCII antes de `$`.
- Fragmento válido; lista exactamente los 8 archivos.
- Nota: no se abrió ningún navegador (K45), por lo que no hay captura ni medición de alto.

## Supuestos
1. Los alto `h` son estimaciones por cálculo (sin navegador): móvil 2200 / 2300 / 3400 / 2000, escritorio 2400 / 2300 / 3100 / 2000. La raíz móvil lleva `min-height` = h y la barra inferior queda al fondo; si el contenido real excede h, el orquestador debe subir `h` del fragmento.
2. P-24 muestra octubre 2026 (parcial, cifras al 21-oct), el único mes del JSON; «Parcial» y «Generar PDF» (no hay PDF de octubre); «Compartir» queda `aria-disabled` hasta generar el PDF. U1 pide «último cerrado» por defecto, pero no hay datos de septiembre.
3. «Despensa y precios» usa la canasta de la semana 2026-W43 (no hay comparativa mensual); peso al inicio y al cierre usa el primer y último registro de octubre (5 y 19 oct); kcal por día, la semana del 12 al 18 oct contra la meta confirmada.
4. P-25: «pendientes posteriores» se arman con `pagos_programados` (parte de Karen, compras kc0 a kc2; las programadas kc3 y kc4 no se incluyen); suman el `saldo_pendiente` del JSON ($3,802.33). Pago por quincena de kc2: $803.27 (el segundo, $803.26).
5. P-21: «Ajuste para bajar» de Axel = 400 kcal (`ajuste_objetivo_kcal`, objetivo bajar); «para subir» vacío con «El de la casa: 0 kcal». Ximena (mantener) con ambos vacíos. «Confirmar 1,816 kcal» solo en Ximena (`requiere_confirmar`). El nivel de actividad se muestra como «Nivel n de 5 · factor» (el JSON no trae nombres de nivel). No se dibujan la gráfica de 12 semanas de kcal (N-233) ni la vista Detalle/Tendencias: están detrás de las pestañas.
6. P-23: la fila de error (rp3, «base de datos ocupada») ofrece «Reintentar envío» como pide U1 §2.2. El estado de cuenta de Karen aparece «Sin envío · se comparte a mano» (decisión A de U1). Los avisos están abiertos en escritorio y plegados en móvil (U1 §6.4). No se dibuja la vista de calendario de envíos.
7. Estilos: componentes en clases dentro de `<helmet><style>` (podadas por archivo y minificadas) y diseño puntual en línea; iconos Lucide por `<symbol>` + `<use>` dentro del mismo archivo. Los tamaños de letra usan las clases `.title-*`, `.body`, `.small`, `.caption` de `tokens.css` (no hay variables de tamaño de letra).
8. El texto «Sin datos en abr-2026» (P-24, estado vacío) usa un mes anterior al primero con datos, para no contradecir el selector.
9. Excepción al fragmento de plantilla: P-21 móvil lleva `"is_interactive": true` porque sus pestañas funcionan.

## Bloqueos
Ninguno.

## Pendientes
- Medir el alto real de los 8 artboards cuando el orquestador los abra en el lienzo y corregir `h`.
- Enlazar «Editar perfil» (P-26.13) y las subsecciones cuando existan sus artboards.
- Revisar en el lienzo el tema oscuro de los marcadores de persona (Ximena y Karen son casi el mismo color; se distinguen por forma y nombre).

## Consumo aproximado
- Leído: ~180 KB (U1 por secciones, DESIGN por secciones, README, format.md, craft.md, LEEME, partes del JSON, tokens).
- Escrito: ~167 KB de artboards + ~2 KB de fragmento + ~5 KB de este reporte = ~174 KB, por encima del tope de 150 KB: cada archivo debe repetir CSS (~5 KB), la hoja de iconos (~3 KB) y la barra lateral o inferior (~3 KB). Además ~56 KB de scripts generadores y de verificación en /tmp/pista3_sr (fuera del repositorio).
- Acciones: ~30 de 60.
