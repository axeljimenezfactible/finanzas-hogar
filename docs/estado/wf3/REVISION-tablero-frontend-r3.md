# REVISIÓN tablero-frontend · ronda 3 de 3 · revisor independiente

Agente `revisor-tablero-frontend-r3` (WF-3c, 12-oct-2026). Trabajo revisado: `front-tablero` con las correcciones `correccion-tablero-frontend-r1` y `correccion-tablero-frontend-r2` (P-01 Resumen = tablero de control, K54; barómetro). Se evalúa **solo** contra `docs/criterios/U5.md` (U5-06, U5-07; U5-09 no aplica), `docs/design/U3-dashboard.md` 2.0 (§2, §1.3, §2.4, B-1), `docs/design/DESIGN.md` §7.6, §8 (con §8.1 y §8.9) y §8.5, y `docs/criterios/C3.md` C3-03. No corregí nada: solo medí y leí.

**Dictamen: APROBADO** · 0 altos · 0 medios · 8 bajos (2 nuevos: N-3 y N-4; 6 que siguen de rondas anteriores). El medio de la ronda 2 (N-1) está **cerrado**.

## 1. Cómo lo verifiqué (comandos propios; temporales en `datos/tmp/rev-tablero-r3/`, ya borrados)

| Paso | Comando | Resultado |
|---|---|---|
| tsc | `cd frontend && npx tsc --noEmit -p .` | EXIT 0, 0 errores |
| build | `npx vite build --outDir ../datos/tmp/rev-tablero-r3/dist --emptyOutDir` (sin tocar `frontend/dist/`) | EXIT 0; 80 entradas precacheadas (2,526 KiB) |
| Playwright de Inicio | configuración propia: `vite preview` del build temporal en el puerto 4188; proyectos `movil` 412×915 táctil y `escritorio` 1440×900; las pruebas del tablero recorren claro y oscuro = **4 combinaciones** | `frontend/pruebas/inicio`: **32 de 32** |
| Pruebas del revisor | spec propia: R1 barómetro 25/35/45 %, R2 sin dato (semilla), R3 gráficas sintéticas, R4 estados (carga, 500 en todo, 501), R5 volcado SVG de G-1 y G-6, R6 árbol accesible de B-1, R7 variantes de B-1 y rejilla 360/412/768/1440, R8 eje y aguja | R1, R2, R3: 12 de 12 (4 combinaciones cada una); R5 a R8: 7 de 7; R4: 0 de 4 **solo** por axe en carga (N-3, baja, fuera del alcance); el resto de R4 medido y correcto |
| C3-03 | `npm run verificar:rutas` | 212 rutas invocadas; 0 fuera de `openapi.yaml`; 0 sin registrar en `routes.go` (221) |
| Hosts | `node scripts/hosts-externos.mjs ../datos/tmp/rev-tablero-r3/dist` | solo `literal-no-peticion` (w3.org, react.dev, licencia de Inter, workbox, comentario de Tailwind) |
| Greps (DESIGN, C2) | hex `#…`, `rgb(`/`hsl(`/`oklch(`, `dark:`, `yAxis: [`, color en `style`, fecha ISO literal en Tablero, GraficaInicio, Barometro, Resumen, comun, Tendencias, mocks y textos | **0** en los 8 archivos |
| Global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | **EXIT 0** («Verificación ampliada OK») |

## 2. Hallazgo de la ronda 2

| Hallazgo r2 | Estado | Evidencia |
|---|---|---|
| **N-1 (media)** · barómetro sin dato con `role="meter"` sin `aria-valuenow` (axe `aria-required-attr` crítico) | **CERRADO** | `Barometro.tsx:57-59`: sin dato, `role="img"` con `aria-label` = «Capacidad de endeudamiento sin dato, tope 40 %» y ningún `aria-value*`. En la semilla (sin ingreso), en las 4 combinaciones: 0 `meter`, 1 imagen con ese nombre exacto, único atributo ARIA = `aria-label`; 3 tramos, marcas «30 %» y «40 %», 0 aguja, 0 valor, dígito «—», «Falta tu ingreso del mes para calcular la capacidad» y «Capturar ingreso» → `/ajustes/03` (P-26.03). axe WCAG 2.0–2.2 A/AA: 0 serios ni críticos en la página; axe **con todas las reglas** sobre B-1: 0 infracciones. Con valor, el `meter` sigue (R1). La prueba de la semilla (`tablero.escritorio.spec.ts:42-46`) ahora exige 0 `meter`, la imagen y axe 0 serios. Lo que DESIGN §8.9 no dice (rol sin dato) quedó PROVISIONAL con `docs/estado/bloqueos/WF3c-barometro-sin-dato.md` (opciones A hecha, B, C descartada), citado en el comentario del componente. |
| N-2, H-3 a H-7 (bajos) | siguen abiertos | fuera del encargo de la corrección r2; `Tablero.tsx` y `mocks.ts` no cambiaron desde la ronda 2 (ver §7) |

## 3. Fidelidad a U3 2.0 (elemento → presente)

| Elemento U3 2.0 | Presente | Evidencia (mía) |
|---|---|---|
| Orden §2.4 (12 bloques) en 412 y 1440; foco = orden de 360 | Sí | A-01 en 32/32; DOM en orden |
| Rejilla 768: K-1 + B-1 (4+4), K-4 a todo el ancho, K-2 + K-3, gráficas una por fila | Sí | R7 (cajas medidas) |
| Rejilla 1440: K-1 + B-1 (6+6); K-2, K-3, K-4 (4+4+4); G-3+G-1, G-6+G-2, G-5+G-4; «Más del hogar» a todo el ancho | Sí | R7 (mismo `top` por fila, anchos 548/357/1120) |
| K-1 a K-4 con líneas, vacíos, errores y acciones | Sí | 32/32; R4 (textos de error de U3 §2 por bloque) |
| B-1: arco, tramos, marcas, aguja, dígito, «de [tope] % de tope», insignia con icono y texto | Sí | R1 |
| B-1: 360 una línea («Pagos … de un ingreso de …»); 412 suma «Te caben …» o «−$7,000.00 arriba del tope»; 768 suma «Con la parte de Karen: [n] %» | Sí | R7 (360, 412, 1440) |
| B-1: valor > escala (75 % con escala 60) → aguja al final (x2 = 182) y dígito «75 %» sin recorte a 360 | Sí | R7 |
| B-1: sin créditos → «0 %» y «Sin créditos activos»; `income_is_expected` → «Ingreso esperado» | Sí | R7 |
| B-1 vacío y error | Sí | R2; R4 («No se pudo calcular la capacidad de endeudamiento. Revisa que estés en la red de casa.» · «Reintentar» · «Ver créditos») |
| G-1 a G-6 con resumen, `figcaption`, tabla y vacíos | Sí | R3; 32/32 |
| Servidor anterior a C1 v1.5 (501) | Sí | R4: B-1, G-1, G-6 = `status` «Disponible cuando el servidor se actualice», 0 «Reintentar» |
| Ningún umbral constante; la interfaz no calcula el estado | Sí | A-18; R1: 25 % con `status` = `exceeded` pinta el tramo rojo (obedece a C1) |
| «Más del hogar» (2, 4, 5, 6) | Sí | 32/32 |

## 4. Barómetro (DESIGN §8.9) con 25 %, 35 % y 45 % en las 4 combinaciones (R1)

Comparación del color calculado de cada pieza contra el del token en el mismo tema (`var(--ui-…)` resuelto en la página):

| Valor | Estado C1 | Valor pintado | Tramos | Insignia | Dígito | `aria-valuetext` |
|---|---|---|---|---|---|---|
| 25 % | normal | = `--ui-role-success` | success, warning, error en orden | «Verde · dentro del rango»: texto = `--ui-role-success`, fondo = `--ui-role-success-surface`, 1 icono `aria-hidden` | «25 %», 32 px, 700, `tabular-nums`, = `--ui-text-default` | «Capacidad de endeudamiento 25 % de 40 %: verde» |
| 35 % | warning | = `--ui-role-warning` | ídem | «Amarillo · cerca del tope» (warning sobre warning-surface) | «35 %» | «… 35 % de 40 %: amarillo» |
| 45 % | exceeded | = `--ui-role-error` | ídem | «Rojo · arriba del tope» (error sobre error-surface) | «45 %» | «… 45 % de 40 %: rojo» |

`role="meter"` con nombre «Capacidad de endeudamiento», `aria-valuemin` 0, `aria-valuemax` 60, `aria-valuenow` = valor; aguja = `--ui-text-default`, 2 px (`size.viz.line`); valor 8 px (`size.meter.height`); marcas «30 %» y «40 %» en `tabular-nums` y `--ui-text-secondary`; axe 0 serios; 0 textos en inglés; 0 controles < 44 px en `movil`. Sin dato: §2.

## 5. Gráficas (§8, §8.1 y §8.5) (R3, R5)

- **Un eje Y**: `ejes()` arma un solo `yAxis` (0 `yAxis: [`); en el SVG, las etiquetas de dinero de G-1 y de % de G-6 están todas del lado izquierdo; las de la derecha de G-6 son las etiquetas de las líneas de referencia «30 %» y «40 %» (`markLine`), no un segundo eje.
- **Barras ≤ 24 px**: G-1 5.4 px (412) y 6.4 px (1440) medidas en el SVG; G-2 y G-3 (HTML) 24 px.
- **Leyendas con refuerzo**: G-1 «Ingreso», «Gasto», «Pagos de crédito» con 2 `<pattern>` (series 2 y 3).
- **Dinero en texto**: resúmenes con montos (G-1 «En oct-2026: ingreso $20,000.00, gasto $1,197.64 y pagos de crédito $598.82»; G-2, G-4, G-5 con valores).
- **aria en español**: `figcaption` en las 6; descripción de ECharts «[título]. [resumen]» en G-1, G-6, G-5 y G-4; enlaces de G-3 «Ver gastos de …» y sugerencias en español.
- Rellenos negros en el SVG de G-6: son `clipPath` y una capa con `fill-opacity="0"` y `stroke-width="0"` de ECharts (invisibles), no colores fuera de tokens.
- axe 0 serios en las 4 combinaciones con las 6 gráficas.

## 6. U5-06, U5-07 y estados (f) (R1 a R4)

- En cada estado medido (barómetro ×3, sin dato, gráficas, carga, 500 en todo, 501) y en las 4 combinaciones: **0 controles < 44 × 44 en `movil`** (en escritorio con puntero fino rige `size.control.fine` = 36 px), 0 iconos fuera de norma, 0 montos sin `tabular-nums`, 0 textos en inglés (ayudante del proyecto más una expresión propia sobre `aria-label` y `aria-valuetext`), 0 fechas ISO visibles ni en ARIA, sin desplazamiento horizontal.
- **Carga**: los 11 bloques con esqueleto (`bg-skeleton-on-card`); 9 con título; **G-5 y G-4 sin título** (N-2, sigue). axe da `aria-prohibited-attr` (serio) en el `Cargando` común (N-3, nuevo, baja, fuera del alcance).
- **Error 500 en todo**: los 11 con mensaje de U3 en español, 1 `alert` y 1 «Reintentar» cada uno; G-5 y G-4 sin título (N-2); axe 0 serios.
- **501**: §3. **Vacío** (semilla): §2; axe 0 serios.

## 7. Hallazgos

### N-3 (baja, nuevo, fuera del alcance de `front-tablero`) · Esqueleto de carga con `aria-label` sin rol (axe `aria-prohibited-attr`)
- **Criterio**: U3 §1.3 y DESIGN §7.6 (estado de carga); no entra en U5-03, que mide las capturas ya cargadas.
- **Qué pasa**: el componente común `Cargando` (`frontend/src/componentes/Estados.tsx:15`, de frontend-base, en HEAD) pinta `<div aria-busy="true" aria-label="Cargando…">` sin rol. Con las operaciones de Inicio sin responder, axe 4.13 (WCAG 2.0–2.2 A/AA) da `aria-prohibited-attr` (serio) en 14 nodos de P-01 (K-1, B-1, K-2, K-3, G-3, G-1, G-6, G-2, G-5, G-4 y tarjetas 2, 4, 5, 6), en las 4 combinaciones. Lo usan todas las pantallas.
- **Corrección** (para el dueño de `Estados.tsx`): `role="status"` en ese `div` (permite el nombre), o quitar `aria-label` y agregar `<span className="sr-only">{textos.app.cargando}</span>`.

### N-4 (baja, nuevo) · Eje de la aguja de 16 px en lugar de `size.viz.marker` (10 px)
- **Criterio**: DESIGN §8.9, fila «Aguja» («eje de `size.viz.marker` con anillo de 2 px de `ui.surface.card`»).
- **Qué pasa**: `Barometro.tsx:75` dibuja el eje con `r={5}` en unidades del `viewBox` (200); el SVG se pinta a 320 px, así que el eje mide **16 px** de diámetro (medido en 412 y 1440). El anillo sí es 2 px (`non-scaling-stroke`).
- **Corrección**: radio en función del ancho pintado (`r = 5 × 200 / ancho`) o un punto HTML encima del SVG (como las marcas) con `size-[var(--size-viz-marker)]` y anillo de 2 px `surface-card`.

### Bajos que siguen abiertos (sin cambios en el código desde la ronda 2)
- **N-2**: G-4 y G-5 sin título en carga y en error (`Tablero.tsx:67-68` con `T3`/`T2` de `Tendencias.tsx:57` y `:93`). Remedido en R4: «No se pudo cargar la gráfica. Revisa que estés en la red de casa.» sin decir cuál.
- **H-3**: nombres accesibles incompletos (`Resumen.tsx:38`; `Tablero.tsx:399`). Nota de esta ronda: sin dato, el enlace de B-1 se llama «Capacidad de endeudamiento · — · Capturar ingreso»: el «—» no dice «sin dato» como sí lo dice la imagen.
- **H-4**: K-1 rebasado con «te quedan —» (`Tablero.tsx:98`).
- **H-5**: marcador de G-6 de 8 px sin anillo (`Tablero.tsx:336` y `:339`, `symbolSize: 8`).
- **H-6**: G-3 con `w ?? 0` en el resumen (`Tablero.tsx:386`).
- **H-7**: `mocks.ts:22` (`payment_count` 1 y `ratio_total_pct` n + 2) y «(ver solicitudes)» en `front-tablero.md:36`.

## 8. HECHO / PREGUNTA / PROPUESTA

- **HECHO**: N-1 cerrado; 0 altos y 0 medios; tsc, build, 32/32, C3-03 y `verify-all` en verde.
- **HECHO (observación, fuera del alcance de `front-tablero`)**: `docs/design/capturas-u5/P-01-*.png` y `axe/P-01-*.json` son de las 11:34–11:41; el tablero llegó después (reporte de `front-tablero` 15:32, `Tablero.tsx` sin versionar). La evidencia U5 de P-01 todavía muestra el Resumen anterior; hay que regenerarla con `make verify-u5-servidor` (dueño: WF-3b o el integrador). Por mi R2, el vacío de la semilla da axe 0 serios, así que U5-03 no debería fallar por B-1.
- **PREGUNTA (bloqueo `WF3c-barometro-sin-dato.md`)**: ¿se fija en DESIGN §8.9, Accesibilidad, «Sin dato: `role="img"` con el texto accesible, sin `aria-valuenow`» (opción A, la que quedó) o se prefiere B (texto en `sr-only`)?
- **PREGUNTA (puerta 4, de la ronda 2, sin hallazgo)**: con servidor anterior a C1 v1.5, K-4 culpa a la red ante un 501 de `getDashboardSettings`; ¿se prefiere «Disponible cuando el servidor se actualice»?
- **PROPUESTA**: corregir N-3 en `Estados.tsx` (una línea) beneficia a todas las pantallas; N-4 y los bajos de `Tablero.tsx` se pueden juntar en una sola pasada.
