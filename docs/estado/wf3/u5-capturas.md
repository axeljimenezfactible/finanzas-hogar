# WF-3b · u5-capturas · Reporte

Agente `u5-capturas` (Sonnet) · 12-oct-2026 · HEAD `6d691d0` (+ cambios sin commit de otros agentes en `frontend/src`; el binario que sirve la app es el compilado antes de esta corrida). Criterios: U5-01, U5-02, U5-03, U5-04, U5-08, U5-10. **No corregí la aplicación**: las fallas van abajo con regla, nodo y pantalla.

## Hecho
- **Lista T por script** (`apoyo.ts`/`verificar_capturas.py` leen `docs/design/U1-arquitectura.md`): 27 pantallas `P-xx` (P-00 a P-26) + 14 secciones `P-26.01` a `P-26.14` = **T = 41, 4T = 164**.
- **Playwright 1.60 contra la aplicación real** (servidor Go en `https://finanzas.local:8443`, base `datos/finanzas-dev.db` con la semilla; certificado mkcert aceptado sin `ignoreHTTPSErrors`). 4 proyectos: `android-claro`, `android-oscuro` (412×915, táctil, `isMobile`) y `escritorio-claro`, `escritorio-oscuro` (1440×900); escala 1; `colorScheme` emulado, **sin tocar el interruptor**. Inicia sesión **una vez por contexto** (formulario real de P-00; la contraseña se lee de `~/.config/finanzas-hogar/dev-semilla.txt` en ejecución y no se imprime ni se guarda). P-00 se captura en un contexto sin sesión.
- Por cada una de las 164 capturas, en la misma carga: PNG de página completa → axe-core → contraste en el DOM → hosts pedidos → tema real. P-03 y P-22 son hojas: se abren con el botón «Registrar gasto» (desde Inicio) y «Registrar peso» (desde Salud); su contraste se mide solo dentro del diálogo. Ids de P-04, P-12, P-18 y P-19 se leen del API con la sesión (compra de $3,094.50, primera receta, crédito Kueski, primera compra).
- Prueba del interruptor (U5-10) en los dos tamaños, 7 pasos en Android y 10 en escritorio.
- Comando (desde `frontend/`): `npx playwright test -c playwright.u5.config.ts` (todo) · `--grep axe` (U5-03) · `--grep contraste` (U5-04) · `--grep tema` (U5-10). Última corrida completa desde cero: **170 passed (7.8 min)**, 0 fallas: 164 capturas + 4 comprobaciones de T + 2 del interruptor.

## Archivos
- Nuevos: `frontend/playwright.u5.config.ts`; `frontend/pruebas-u5/{apoyo.ts, capturas.spec.ts, tema.spec.ts, cierre.ts, contraste-dom.js, verificar_capturas.py, resumen_u5.py}`.
- Evidencia: `docs/design/capturas-u5/<P-xx>-<android|escritorio>-<claro|oscuro>.png` (164, 24 MB con los JSON), `axe/<mismo nombre>.json` (164), `contraste/<mismo nombre>.json` (164), `hosts.json`, `evidencia/{manifiesto.json, cifras-oraculo.json, tema-android-claro.json, tema-escritorio-claro.json}`.
- Bloqueos: `docs/estado/bloqueos/U5-contraste-ui-controles.md`, `docs/estado/bloqueos/U5-conteo-carpeta.md`. Este reporte.
- Temporales borrados (`datos/tmp/playwright-u5-capturas` y los logs de las corridas). Los JSON de contraste y axe son salida del script, no escritura manual.

## Verificación
- **U5-01 · PARCIAL (162 de 164)**. `python3 -I frontend/pruebas-u5/verificar_capturas.py` → `T = 41 · 4T = 164`; `PNG correctos: 162 de 164 · faltantes: 0 · mal: 2`. Los 2 son `P-26.13-android-claro.png` y `P-26.13-android-oscuro.png`: **415×1732 en vez de 412** (desborda 3 px a lo ancho: el selector de vista Tabla/Tablero/Calendario/Galería llega a 415; la tabla interna llega a 1117 px). Es defecto de la app, no de la captura. Los 164 existen, no están en blanco (desviación estándar > 0 en todas) y las 160 restantes miden 412 o 1440 de ancho. `ls docs/design/capturas-u5/ | wc -l` da **168** (164 PNG + `axe`, `contraste`, `evidencia`, `hosts.json`; `ls *.png | wc -l` = 164): bloqueo `U5-conteo-carpeta`. Nota de método: en Android la captura es de página completa y la barra inferior fija aparece a media imagen (artefacto de `fullPage`).
- **U5-02 · OK**. Capturas generadas por el comando anterior desde la aplicación real. 5 cifras visibles = `oraculo/model.json` (todas existen allí, `cifras-oraculo.json`): disponible inicial A **$10,000.00** (P-17, P-18); B **$6,905.50** (P-01, P-03, P-17, P-18, P-26.05); pago del 30-oct-2026 **$973.75** (P-01, P-17, P-18, P-24); compra del 10-oct **$3,094.50** (P-01, P-02, P-04, P-14); deuda viva al 10-abr-2027 **$6,293.44** (P-09, P-26.13); parte de Karen **$1,191.50** (P-04, P-19, P-25). Ninguna pantalla cayó en el error de página ni quedó cargando al capturar.
- **U5-03 · FALLA (20 violaciones serias/críticas = 5 reglas × 4 combinaciones, las mismas en las 4)**. Comando literal de U5-03 → `20`. Por combinación (serias/críticas · otras · incompletos): android-claro 5 · 0 · 15; android-oscuro 5 · 0 · 15; escritorio-claro 5 · 0 · 43; escritorio-oscuro 5 · 0 · 43. Lista (JSON en `axe/`):
  1. **P-06** `aria-meter-name` (serious), 3 nodos `div[role=meter][aria-valuetext="60 % del plan"]` (y «123 %», «116 %»): las barras de avance no tienen nombre accesible.
  2. **P-09** `aria-meter-name` (serious), 2 nodos `div[aria-valuetext="700 % del plan"]` y `"50 % del plan"`.
  3. **P-08** `button-name` (**critical**), 1 nodo `th:nth-child(5) > button.inline-flex.h-9` (botón de orden de una columna, sin texto).
  4. **P-26.05** `nested-interactive` (serious), 1 nodo: `<button>` dentro de `<button>` («Kueski»).
  5. **P-26.06** `nested-interactive` (serious), 4 nodos: igual, en la columna de nombre de la tabla (Acueducto de Guadalupe, Bodega Aurrera, …).
  Incompletos revisados por tipo (suman 1,186 nodos): `color-contrast` 1,158 («parcialmente cubierto por otro elemento» 1,066, sobre todo encabezados de tabla pegajosos y la barra fija; «contiene una imagen» 80; «contenido demasiado corto» 12): todos son elementos con texto que mi script del DOM sí evaluó (0 de texto bajo umbral, ver U5-04); `aria-hidden-focus` 24 (guardas de foco de Radix, `tabindex=0` bajo `aria-hidden`, en las hojas P-03/P-22; revisar a mano); `aria-prohibited-attr` 4 (P-26.10, `div[aria-label="Muestra de tarjeta"]` sin rol).
- **U5-04 · FALLA solo en la lectura literal (536 = 4 × 134 pares de UI bajo 3:1); texto e íconos 0 bajo umbral**. Pares evaluados (elementos: texto + íconos + controles con límite visible) y bajo umbral, por combinación:

  | Combinación | Pantallas | Pares evaluados | Bajo umbral |
  |---|---|---|---|
  | android-claro | 41 | 4,550 | 134 |
  | android-oscuro | 41 | 4,564 | 134 |
  | escritorio-claro | 41 | 5,764 | 134 |
  | escritorio-oscuro | 41 | 5,778 | 134 |

  Los 536 son de tipo `ui`, subtipo **`boton`** (408: tarjetas-botón y chips con borde de separación, `#E5E6E9` sobre `#FAFAFB` = 1.196 en claro, `#4D4E53` sobre `#222325` = 1.895 en oscuro; P-06, P-07, P-09, P-10, P-11, P-17, P-18) o **`seleccion`** (128: pestañas de vista Tabla/Galería y chips cuyo estado se marca con relleno, `#F0F1FF` sobre `#FAFAFB` = 1.075 claro y `#303258` sobre `#222325` = 1.29 oscuro; P-02, P-03, P-04, P-08, P-11, P-12, P-15, P-19, P-20, P-22 a P-25, P-26.01 y otras de Ajustes). **0 pares `campo`** (inputs, selects, casillas, interruptores) bajo umbral. Todos esos controles llevan texto propio; si la regla excluye los botones con texto del 3:1 el resultado de «seleccion» queda como único tema (estado solo por relleno). Cada par lista razón, fg, bg, nodo y texto; las razones se recalcularon de los hex con otra implementación en Python: 2,538 pares únicos, diferencia máxima 0.0005. Hay 0 pares «indeterminados» (imagen o degradado). Fuera de medición en el DOM: marcas de las gráficas ECharts (el texto de las gráficas SÍ entra) y `::placeholder`. Bloqueo `U5-contraste-ui-controles` pide la regla.
- **U5-08 · OK**. `hosts.json`: **1 host único, `finanzas.local:8443`**, registrado por petición (`context.on('request')`, http/https/ws) en las 164 capturas; cada prueba además exige 0 hosts fuera de `finanzas.local`/`localhost`/`127.0.0.1`.
- **U5-10 · OK**. En las 164 capturas el `background-color` computado del `body` = token de C2 (`ui.surface.page`: claro `#FAFAFB` = `rgb(250, 250, 251)`; oscuro `{color.neutral.950}` = `#222325` = `rgb(34, 35, 37)`), sin `data-theme` en `<html>` y con `prefers-color-scheme` igual al `colorScheme` emulado (manifiesto: 0 desajustes). Interruptor (P-26.10 en los dos tamaños y el botón de la barra lateral en escritorio): sistema claro/oscuro sin tocar nada; «Claro» con sistema oscuro → claro; «Oscuro» con sistema claro → oscuro; «Sistema» vuelve a seguir al sistema; «Claro» persiste al recargar (`fh-tema`). 7 y 10 pasos en `evidencia/tema-*.json`.

## Supuestos
1. Una captura por pantalla en su subsección por omisión (Resumen, Detalle, etc. según U1); las otras subsecciones no se capturan (U5-01 pide por `P-xx`).
2. Escala 1: el tamaño del PNG es el del viewport en ancho (412 / 1440) y la altura de página completa (≥ 915 / ≥ 900).
3. P-00 se captura sin sesión; el resto con la cuenta de desarrollo `axel`. P-03 y P-22 se capturan con la hoja abierta sobre Inicio y Salud.
4. Ids de P-04, P-12, P-18, P-19: la compra del 10-oct-2026, la primera receta, el crédito Kueski y su primera compra.
5. Contraste: el fondo efectivo se compone por ascendencia hasta una capa opaca; los canales se redondean a 8 bits antes de la razón; opacidad de ancestros se aplica al texto; grande = ≥ 24 px o ≥ 18.66 px en negrita; se excluyen controles deshabilitados (WCAG 1.4.11 los exime) y controles sin borde ni relleno visible.
6. Los PNG se toman con `animations: 'disabled'` tras esperar red en reposo, ausencia de `aria-busy` y 0.9 s para las gráficas.
7. `hosts.json` y el manifiesto se juntan en `globalTeardown` (`cierre.ts`) y se fusionan con lo previo si se corre un subconjunto.

## Bloqueos
`docs/estado/bloqueos/U5-contraste-ui-controles.md` (qué cuenta como componente de UI en U5-04) y `docs/estado/bloqueos/U5-conteo-carpeta.md` (`ls | wc -l` de U5-01 da 168). Hallazgos de la aplicación (no son bloqueos; los atienden los agentes de corrección): los 5 de axe, el desborde de P-26.13 en Android y los 128 de «seleccion».

## Pendientes
- Corrección de los 5 hallazgos de axe, del desborde de 3 px en P-26.13 Android y decisión sobre los pares de UI; luego re-correr `npx playwright test -c playwright.u5.config.ts` (≈ 8 min) y `python3 -I frontend/pruebas-u5/resumen_u5.py`.
- U5-11 (marcar «capturado» en la matriz C3-05 con las 4 rutas) queda para el integrador; U5-07 y U5-09 y U5-13 son de otros agentes.
- Revisar a mano los 24 `aria-hidden-focus` y los 4 `aria-prohibited-attr` incompletos.

## Propuestas para Axel (PROPUESTA, no decidido)
| Qué | Resultado | Costo | Riesgo |
|---|---|---|---|
| Dar nombre a las barras de avance de P-06 y P-09 (`aria-label` con la categoría) | Se van 2 de las 5 fallas de axe | Bajo | Ninguno |
| Quitar el botón anidado en P-26.05 y P-26.06 y dar nombre al botón de orden de P-08 | Se van las otras 3 fallas (una es crítica) | Bajo | Ninguno |
| Subir el borde de chips y tarjetas-botón a `border-control` y marcar la pestaña activa con borde o subrayado ≥ 3:1 | Quedan 0 pares de UI bajo 3:1 sea cual sea la regla | Medio (toca C2/DESIGN) | Cambia el aspecto en las 4 combinaciones |
| Aceptar la lectura B del bloqueo (botones con texto fuera del 3:1) | Cierra U5-04 sin tocar diseño | Ninguno | Solo vale si Axel y el revisor la aceptan |
Recomendación: corregir las 5 de axe ahora y decidir los chips en la puerta 4. Observación visual (no medida): en escritorio la etiqueta «Axel» del encabezado de la barra lateral se empalma con el ícono de plegar (P-17 escritorio oscuro).

## Paquetes instalados
Ninguno (Playwright 1.60 y `@axe-core/playwright` ya estaban; `identify` de ImageMagick y `numpy` ya existían).

## Consumo aproximado
≈ 38 acciones de herramienta · ≈ 120 KB leídos · ≈ 45 KB escritos a mano (scripts ≈ 31 KB, bloqueos ≈ 4 KB, reporte ≈ 10 KB); la corrida generó además ≈ 24 MB de PNG y JSON.
