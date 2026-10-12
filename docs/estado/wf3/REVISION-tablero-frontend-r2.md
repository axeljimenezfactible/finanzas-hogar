# REVISIÓN tablero-frontend · ronda 2 de 3 · revisor independiente

Agente `revisor-tablero-frontend-r2` (WF-3c, 12-oct-2026). Trabajo revisado: `front-tablero` con la corrección `correccion-tablero-frontend-r1` (P-01 Resumen = tablero de control, K54; barómetro). Se evalúa **solo** contra `docs/criterios/U5.md` (U5-06, U5-07; U5-09 no aplica), `docs/design/U3-dashboard.md` 2.0 (§2, §1.3, §7), `docs/design/DESIGN.md` §7.6, §8 (incluye §8.1 y §8.9) y §8.5, y `docs/criterios/C3.md` C3-03. No corregí nada: solo medí y leí.

**Dictamen: NO APROBADO** · 0 altos · **1 medio (nuevo)** · 6 bajos (aprobar pide 0 altos y 0 medios). Los 2 medios de la ronda 1 (H-1 y H-2) están **cerrados**.

## 1. Cómo lo verifiqué (comandos propios; temporales en `datos/tmp/rev-tablero-r2/`, ya borrados)

| Paso | Comando | Resultado |
|---|---|---|
| tsc | `cd frontend && npx tsc --noEmit` | 0 errores |
| build | `npx vite build --outDir ../datos/tmp/rev-tablero-r2/dist --emptyOutDir` (sin tocar `frontend/dist/`) | listo; 80 entradas precacheadas |
| Playwright de Inicio | configuración propia (`vite preview` del build temporal, puerto 4187), proyectos 412×915 táctil y 1440×900; las pruebas del tablero recorren claro y oscuro = **4 combinaciones** | `pruebas/inicio`: **32 de 32** |
| Pruebas del revisor | spec propia: R1 barómetro 25/35/45 %, R2 gráficas (aria, tablas, regiones, axe), R3 estados (carga, 500, 501, 404, señales 501, vacío), R4 detalle axe del vacío de B-1, R5 formas SVG y anillo de foco | **18 de 18** (R1 a R3 y R5 en las 4 combinaciones; R4 en los 2 anchos); observaciones en JSON |
| C3-03 | `npm run verificar:rutas` | 212 rutas invocadas; 0 fuera de `openapi.yaml`; 0 sin registrar en `routes.go` (221) |
| Hosts | `node scripts/hosts-externos.mjs <dist temporal>` | solo literales que no son peticiones (licencia de Inter, workbox, comentario de Tailwind) |
| Greps (DESIGN, C2) | hex, `rgb(`/`oklch(`/`hsl(`, `dark:` y colores en `style` en Tablero, GraficaInicio, Barometro, Resumen, comun, Tendencias, mocks y textos | 0 en los 8 archivos; `yAxis: [` (doble eje) = 0 |
| U5 sin servidor | `make verify-u5` | EXIT 0 (U5-04 bajo umbral = 0, U5-08 = 0, U5-11 = 0) |
| Global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | **EXIT 0** («Verificación ampliada OK») |

## 2. Hallazgos de la ronda 1

| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| H-1 (media) · «servidor anterior a C1 v1.5» | **CERRADO** | `Tablero.tsx:211-216` (B1) y `:245-257` (`esServidorViejo`, `ServidorViejo`, `EstadosGrafica soloNueva`); G1 y G6 con `soloNueva` (`:292`, `:312`). Con 501 y con 404 en `getDashboardDebtCapacity` y `getDashboardMonthlyOverview`, en las 4 combinaciones: B-1, G-1 y G-6 = `role="status"` «Disponible cuando el servidor se actualice», 12 px = `text-caption`, color = `--ui-role-info` sobre `--ui-role-info-surface`, icono 20×20, 0 `role="alert"` en `main`, 0 «Reintentar», sin `$` ni `%`; B-1 conserva «Ver créditos»; G-3 y K-4 siguen con sus datos; axe 0 serios (incluye `color-contrast`). |
| H-2 (media) · axe `scrollable-region-focusable` | **CERRADO** | `GraficaInicio.tsx:66-69` (`ContenedorTabla`), usado en `FiguraGrafica` y `FiguraBarras` (`Tablero.tsx:275`). Con las 6 gráficas pintadas (respuestas SINTÉTICAS): axe 0 serios en las 4 combinaciones, ocultas y con las 6 en «Ver como tabla». Oculta: 0 regiones. Visible: `role="region"`, nombre = título (G-3, G-1, G-6, G-2, G-5, G-4), `tabindex="0"`, recibe foco; con teclado (Shift+Tab) el anillo es `solid 2px` = `--ui-border-focus` en claro y oscuro. |
| H-3 a H-7 (bajos) | siguen abiertos (fuera del encargo de la corrección) | ver §7 |

## 3. Fidelidad a U3 2.0 (elemento → presente)

| Elemento U3 2.0 | Presente | Evidencia |
|---|---|---|
| Orden §2.4 (12 bloques) en 412 y 1440 | Sí | prueba A-01 (32/32) |
| K-1, K-2, K-3, K-4 con sus líneas, vacíos, errores y acciones | Sí | A-05, A-06, A-07, A-17; mis observaciones de vacío y 500 |
| B-1 completo (arco, tramos, marcas, aguja, dígito, «de [tope] % de tope», insignia, «Te caben» / «arriba del tope», «Sin créditos activos», vacío sin aguja, error) | Sí (con el hallazgo N-1 de accesibilidad en el vacío) | §4 |
| G-1 a G-6 con resumen, tabla, vacíos y reglas | Sí | §5 |
| Una petición por bloque; G-1 y G-6 de una sola | Sí | A-02 |
| Ningún umbral constante; la interfaz no calcula porcentajes ni estados | Sí | A-18; 75 % con escala 60: aguja al final y dígito «75 %» |
| Servidor anterior a C1 v1.5 (§2, reglas del tablero) | **Sí** (antes no) | §2, H-1 |
| Campo o señal nueva que falta: la línea no se pinta | Sí | `getBudgetSignals` 501: G-3 con barras y sin sugerencias (§7 «sin sugerencias si falla») |
| «Más del hogar» (2, 4, 5, 6) | Sí | prueba de Inicio |

## 4. Barómetro (DESIGN §8.9) con 25 %, 35 % y 45 % en las 4 combinaciones

| Valor | Estado | Valor pintado | Tramos | Insignia | Dígito | `aria-valuetext` |
|---|---|---|---|---|---|---|
| 25 % | normal | `--ui-role-success` (claro y oscuro) | success/warning/error en orden | «Verde · dentro del rango», color `--ui-role-success` sobre `-surface`, 1 icono | «25 %», 700, `tabular-nums`, `--ui-text-default` | «Capacidad de endeudamiento 25 % de 40 %: verde» |
| 35 % | warning | `--ui-role-warning` | ídem | «Amarillo · cerca del tope» | «35 %» | «… 35 % de 40 %: amarillo» |
| 45 % | exceeded | `--ui-role-error` | ídem | «Rojo · arriba del tope» | «45 %» | «… 45 % de 40 %: rojo» |

`aria-valuemin` 0 y `aria-valuemax` 60; aguja en `--ui-text-default`; marcas «30 %» y «40 %»; axe 0 serios con valor. **Sin dato** (semilla sin ingreso): `role="meter"` **sin** `aria-valuenow` → axe `aria-required-attr` (crítico) en las 4 combinaciones: hallazgo N-1.

## 5. Gráficas (§8, §8.1 y §8.5)

- **Un eje Y**: `ejes()` arma un solo `yAxis` y G-6 lo extiende; 0 `yAxis: [` en el código.
- **Barras ≤ 24 px**: G-1 6.4 px (1440) y 5.4 px (412) medidos en el SVG; G-2 y G-3 (HTML) 24 px.
- **Leyendas con refuerzo**: G-1 «Ingreso», «Gasto», «Pagos de crédito» con 2 `<pattern>` (series 2 y 3; 15 marcas con patrón); textos del SVG en `ui.text.secondary` o `ui.text.default`, nunca en el color de la serie.
- **Dinero en texto**: G-2 y G-3 con el valor en la fila; G-5 y G-4 con valores finales; G-1 a la tabla (columna < 72 px).
- **aria en español**: contenedor de ECharts «[título]. [resumen]» en G-1, G-6, G-5 y G-4; figuras con `figcaption`; tablas con `caption`.

## 6. U5-06, U5-07 y estados (f)

- En cada pantalla medida (barómetro ×3, gráficas, tablas visibles, 500, 501, 404, vacío) y en las 4 combinaciones: 0 controles < 44 × 44 con puntero grueso, 0 iconos fuera de norma, 0 montos sin `tabular-nums`, 0 textos en inglés (ayudante del proyecto y expresión propia más amplia sobre `main`, `aria-label` y `aria-valuetext`), 0 fechas ISO, sin desplazamiento horizontal.
- **Carga**: los 11 bloques con esqueleto (`bg-skeleton-on-card`); K-1 a K-4, B-1, G-1, G-2, G-3 y G-6 con título; **G-4 y G-5 sin título** (N-2, baja).
- **Error 500 en todas**: los 11 con mensaje en español y 1 «Reintentar»; G-4 y G-5 sin título (N-2); axe 0 serios.
- **501 / 404**: §2.
- **Vacío** (semilla): textos y acciones de U3 §7; axe 1 crítico (N-1).

## 7. Hallazgos

### N-1 (media, nuevo) · Barómetro sin dato: `role="meter"` sin `aria-valuenow` (axe crítico)
- **Criterio**: U5-06 (barómetro de DESIGN §8.9 y accesibilidad del Resumen); DESIGN §8.9 Accesibilidad (`role="meter"` con `aria-valuenow`); inconsistencia con `docs/design/capturas-u5/axe/P-01-*.json` (0 serios): U5-03 fallará en la próxima captura, porque `datos/semilla/` no trae ingresos y B-1 queda sin dato en el servidor de desarrollo.
- **Qué pasa**: `frontend/src/componentes/Barometro.tsx:54` pone `aria-valuenow={conValor ? valor : undefined}`. Sin ingreso, axe 4.13 (WCAG 2.0–2.2 A/AA) da `aria-required-attr` (impacto crítico) en `<div role="meter" aria-label="Capacidad de endeudamiento" aria-valuemin="0" aria-valuemax="60" aria-valuetext="Capacidad de endeudamiento sin dato, tope 40 %">`: «Required ARIA attribute not present: aria-valuenow». Pasa en 412 y 1440, claro y oscuro. La prueba nueva de axe usa solo respuestas con valor (27 %) y la prueba de la semilla (`tablero.escritorio.spec.ts`, bloque de B-1 sin dato) da por bueno el `meter` sin valor.
- **Corrección**: sin valor, no usar `role="meter"`: `role="img"` con `aria-label={textoAccesible}` (mismo texto «… sin dato, tope 40 %»), o un contenedor sin rol con ese texto en `sr-only`. Con valor sigue el `meter` actual. Ajustar la prueba de la semilla (buscar por `role="img"` o por el texto) y correr axe también en el vacío de la semilla en las 4 combinaciones. DESIGN §8.9 no dice qué rol usar sin dato: ver PROPUESTA.

### N-2 (baja, nuevo) · G-4 y G-5 sin título en carga y en error
- **Criterio**: U3 §1.3 («La tarjeta conserva su título y su enlace»); consistencia con G-1, G-2, G-3 y G-6 de la misma pantalla.
- **Qué pasa**: `Tablero.tsx:67-68` usa `T3`/`T2` de `Tendencias.tsx` (`Estados` sin título, líneas 57 y 93). En carga el bloque no tiene texto; con 500 dice «No se pudo cargar la gráfica. Revisa que estés en la red de casa.» sin decir cuál. La ronda 1 dio el título por visible en los 11 bloques: no lo está en estos dos.
- **Corrección**: pasar el título a `T2`/`T3` para su esqueleto y su falla (o envolverlos en el tablero como hace `EstadosGrafica`).

### Bajos de la ronda 1 que siguen abiertos (sin cambios en el código)
- **H-3**: nombres accesibles incompletos (`Resumen.tsx:38` sin comparación de K-1 ni disponible de K-3; `Tablero.tsx:399` `aria-label={T.g3.abrir(n)}` tapa estado y sugerencia).
- **H-4**: K-1 rebasado sigue con «te quedan —» (`Tablero.tsx:98`).
- **H-5**: marcador de G-6 de 8 px sin anillo (`Tablero.tsx:336` y `:339`; §8.1 pide 10 px con anillo de 2 px de `ui.surface.card`).
- **H-6**: G-3 sigue con `w ?? 0` en el resumen (`Tablero.tsx:386`).
- **H-7**: `mocks.ts:22` (`payment_count` 1 y `ratio_total_pct` n + 2) y «(ver solicitudes)» en `front-tablero.md:36`.

## 8. PREGUNTA / PROPUESTA
- **PROPUESTA (DESIGN §8.9, Accesibilidad)**: agregar «sin dato: `role="img"` con el texto accesible, sin `aria-valuenow`», para que N-1 no dependa de una decisión de quien implementa. La solicitud va a quien edita DESIGN; no cambia tokens ni cifras.
- **PREGUNTA (puerta 4, sin hallazgo)**: con un servidor anterior a C1 v1.5, `getDashboardSettings` responde 501 y K-4 muestra «No se pudieron cargar los pagos. Revisa que estés en la red de casa.» con «Reintentar». Es lo que manda U3 (línea 65, PROVISIONAL, `U3-tablero-recordatorio.md`), pero el texto culpa a la red. ¿Axel prefiere que K-4 use 5 o «Disponible cuando el servidor se actualice» en ese caso?
