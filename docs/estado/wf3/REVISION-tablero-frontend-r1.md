# REVISIÓN tablero-frontend · ronda 1 de 3 · revisor independiente

Agente `revisor-tablero-frontend-r1` (WF-3c, 12-oct-2026). Trabajo revisado: `front-tablero` (P-01 Resumen = tablero de control, K54; barómetro; iconos de sección de la barra lateral). Se evalúa **solo** contra `docs/criterios/U5.md` (U5-06, U5-07; U5-09 no aplica), `docs/design/U3-dashboard.md` 2.0 (criterios `docs/criterios/U3.md` 1.4), `docs/design/DESIGN.md` §7.6, §8, §8.5 (y §8.9, que §8 contiene) y `docs/criterios/C3.md` C3-03. No corregí nada: solo medí y leí.

**Dictamen: NO APROBADO** · 0 altos · **2 medios** · 5 bajos (aprobar pide 0 altos y 0 medios).

## 1. Cómo lo verifiqué (comandos propios, temporales en `datos/tmp/rev-tablero-r1/`, ya borrados)

| Paso | Comando | Resultado |
|---|---|---|
| tsc | `cd frontend && npx tsc --noEmit` | 0 errores (8.3 s) |
| build | `npx vite build --outDir ../datos/tmp/rev-tablero-r1/dist` (sin tocar `dist/` del servidor de desarrollo) | listo; 80 entradas precacheadas |
| C3-03 | `npm run verificar:rutas` | 212 rutas invocadas; 0 fuera de `openapi.yaml`; 0 sin registrar en `routes.go` (221) |
| Hosts | `node scripts/hosts-externos.mjs <dist temporal>` | 0 peticiones externas (3 literales que no son peticiones: licencia de Inter y workbox) |
| Playwright de Inicio | configuración propia (vite preview del build temporal, puerto 4183), proyectos 412×915 táctil y 1440×900; cada prueba del tablero recorre claro y oscuro = **4 combinaciones** | `pruebas/inicio`: **26 de 26** |
| Pruebas del revisor | `spec/revision.spec.ts` propia: barómetro 25/35/45 %, gráficas (SVG y HTML medidos), idioma, ISO, ≥ 44 px, iconos, tabulares, axe WCAG 2.0–2.2 A/AA, rejilla, estados carga/500/501/vacío | **10 de 10** corridas (5 por ancho); observaciones en JSON |
| Suite completa | `npx playwright test` con todo `frontend/pruebas` contra el build temporal (puerto 4184) | **253 passed, 4 skipped, 0 failed** (5.5 min) |
| Greps | hex, `rgb(`, `oklch(`, `dark:` y colores en `style` en Tablero, Barometro, GraficaInicio, Resumen, comun, mocks, textos y BarraLateral | 0 en los 8 archivos |

## 2. Fidelidad a U3 2.0 (elemento → presente)

| Elemento U3 2.0 | Presente | Evidencia |
|---|---|---|
| Orden §2.4 (12 bloques) en el DOM en 412 y 1440 | Sí | `[data-bloque]` = K-1, B-1, K-4, K-2, K-3, G-3, G-1, G-6, G-2, G-5, G-4, mas |
| Rejilla 1440: K-1 + B-1 (6 + 6); K-2, K-3, K-4 (4 + 4 + 4); gráficas 6 + 6 por pares; «Más del hogar» a lo ancho; mismo alto por fila | Sí | medido: K-1 y B-1 548 × 486; K-2, K-3, K-4 357 × 303 (x 280, 661, 1043); G-3/G-1, G-6/G-2, G-5/G-4 en pares de 548 |
| K-1: valor, «gastado en [mes]», «[n] % de [plan] · te quedan», barra `role="meter"` «[n] % del plan», efectivo, comparación con palabras e icono, «Sin datos de [mes]», insignias, vacío → P-26.02, error, acción → P-06 | Sí (detalle en H-4) | código K1 y pruebas A-06 (oct y nov) |
| K-2: «+[quincenas]» de rol ingreso, extra / «Sin extra», «[n] de [m] quincenas capturadas», comparación, «Ingreso esperado», vacío → P-26.03, error | Sí | código K2; seed: «Falta tu ingreso quincenal» |
| K-3: deuda viva de rol, «Pagos de [mes]: [total] · [n] pagos», Hogar · Karen, «Falta pagar», «Disponible hoy … · Lectura B», «Pendiente confirmar», «+[n] líneas», vacíos, error | Sí («Reintentar» parcial: PROVISIONAL, bloqueo `WF3c-tablero-interfaz` punto 2) | prueba A-05 |
| K-4: N de `getDashboardSettings`, vencidos primero, fecha con día, línea, monto, estado, Hogar · Karen desde 412, fila → P-18, «Ver calendario», vacíos, error | Sí («y [n] más»: PROVISIONAL, punto 1) | A-05, A-17 «Vence mañana» |
| B-1: arco 0 → `scale_max_pct`, tramos de rol, marcas «[n] %», aguja, dígito, «de [tope] % de tope», insignia icono + palabra, «Te caben» / «arriba del tope», «Con la parte de Karen» desde 768, vacío sin aguja, «Sin créditos activos», error | Sí | §3 |
| G-1: columnas agrupadas, 3 series en orden con patrón, nota fija, resumen, tabla, vacío | Sí (separación 2 px: PROVISIONAL, punto 5) | §4 |
| G-2: barras horizontales, orden descendente, valor al final, resumen, vacío con `circle-check` | Sí | §4 |
| G-3: barras en % del plan, marcas 80 % (1 px, discontinua) y 100 % (2 px), rol advertencia con `triangle-alert` y texto, orden descendente, 8 en móvil + «Ver todas», «sin plan», sugerencias, fila → P-02 con `categoria` y `periodo` | Sí | prueba A-20, A-21 |
| G-4 = T-2 y G-5 = T-3 sin filtros | Sí | figcaption |
| G-6: línea, proyección 8 · 4, referencias 30 % y 40 % con texto, «Hoy», eje de 0 al mayor entre escala y dato, «— · Falta ingreso», resumen | Sí | §4 |
| Una petición por bloque; G-1 y G-6 de una sola (`from` = mes − 11, `to` = mes + 5) | Sí | prueba A-02 |
| Ningún umbral constante; la interfaz no calcula porcentajes, diferencias ni estados | Sí | lectura de Tablero.tsx y Barometro.tsx |
| «Disponible cuando el servidor se actualice» para un bloque que depende solo de una operación nueva (U3 §2, reglas del tablero) | **No** | hallazgo H-1 |
| «Más del hogar» con las tarjetas 2, 4, 5 y 6 (PROVISIONAL) | Sí | `h2` + `tarjeta-2,4,5,6` |
| §2.3 «Qué miro para prevenir» | No es bloque de la interfaz (no está en la tabla de bloques ni en el orden §2.4); sin hallazgo | — |

## 3. Barómetro (DESIGN §8.9) con 25 %, 35 % y 45 % (respuesta de C1 simulada, ingreso SINTÉTICO), en las 4 combinaciones

| Valor | Estado de C1 | Valor pintado = token | Tramos = `ui.role.success/warning/error` | Insignia (color, fondo, icono) | Dígito | `aria-valuetext` |
|---|---|---|---|---|---|---|
| 25 % | normal | `--ui-role-success` (claro `rgb(32,123,62)`, oscuro `rgb(126,207,144)`) | sí, 8 px | `ui.role.success` sobre `-surface`, 1 icono, «Verde · dentro del rango» | «25 %», 32 px, 700, `tabular-nums`, `ui.text.default` | «Capacidad de endeudamiento 25 % de 40 %: verde» |
| 35 % | warning | `--ui-role-warning` | sí | «Amarillo · cerca del tope» | «35 %» | «… 35 % de 40 %: amarillo» |
| 45 % | exceeded | `--ui-role-error` | sí | «Rojo · arriba del tope» | «45 %» | «… 45 % de 40 %: rojo» |

`aria-valuemin` 0, `aria-valuenow` = valor, `aria-valuemax` = 60 (`scale_max_pct`); aguja en `ui.text.default`; marcas «30 %» y «40 %» en HTML; sin dato: dígito «—» y sin aguja. axe `color-contrast`: 0 en las 4 combinaciones.

## 4. Gráficas (§8 y §8.5)

- **Un eje Y** en todas (`ejes()` arma un solo `yAxis`; G-6 lo sobrescribe sin agregar otro).
- **Barras ≤ 24 px**: G-1 6.4 px (1440) y 5.4 px (412); G-2 y G-3 (HTML) 24 px de alto (`size.viz.bar-max`).
- **Leyenda con refuerzo**: G-1 «Ingreso», «Gasto», «Pagos de crédito» con 2 patrones SVG (series 2 y 3); G-6, una serie con su nombre.
- **Etiquetas de dinero en texto**: G-2 y G-3 con el valor en la barra; en G-1 la columna mide < 72 px y el valor pasa a la tabla (§8.5 lo permite); «Ver como tabla» en las 6.
- **aria en español**: el contenedor de ECharts dice «[título]. [resumen]» (y además es `aria-hidden`); figuras con `figcaption`; tabla con `caption`.

## 5. U5-06 y U5-07 (412 táctil y 1440, claro y oscuro)

0 controles < 44 × 44 con puntero grueso; 0 iconos fuera de norma (20/24 px, trazo 1.5 dibujado; incluida la barra lateral); 0 montos sin `tabular-nums`; 0 textos en inglés (ayudante del proyecto y expresión propia sobre `main` y el texto de los SVG); 0 fechas ISO; sin desplazamiento horizontal (412 = 412, 1440 = 1440). axe: 0 `color-contrast`; **1 regla seria** (hallazgo H-2).

## 6. Estados (f)

| Estado | Resultado |
|---|---|
| Carga | los 11 bloques con esqueleto (`skeleton-on-card`) y título visible; sin ejes falsos |
| Error 500 en todas | los 11 con su mensaje en español sin jerga y «Reintentar»; los demás no caen |
| Vacío | K-2, B-1, G-1, G-3, G-6, G-2 y K-4 con su texto y acción de U3 |
| Servidor anterior a C1 v1.5 (501 en las operaciones nuevas) | **B-1, G-1 y G-6 dicen «Revisa que estés en la red de casa» y ofrecen «Reintentar»** (H-1) |

## 7. Hallazgos

### H-1 (media) · Estado «servidor anterior a C1 v1.5» ausente
- **Criterio**: U3 2.0 §2, reglas del tablero («un bloque que depende solo de una operación nueva muestra «Disponible cuando el servidor se actualice» (`font.size.caption`, `info`)»); U3-03.
- **Qué pasa**: `textos.tablero.servidorViejo` está en `frontend/src/pantallas/inicio/textos.ts:40` y nada lo usa. Con 501 en `getDashboardDebtCapacity` y `getDashboardMonthlyOverview`, B-1 muestra «No se pudo calcular la capacidad de endeudamiento. Revisa que estés en la red de casa.» con «Reintentar», y G-1 y G-6 «No se pudo cargar la gráfica. Revisa que estés en la red de casa.»: le dice a la persona que revise la red cuando el servidor no tiene la operación.
- **Corrección**: en `Tablero.tsx` (B1 y `EstadosGrafica` de G-1 y G-6), si la respuesta es 501 (o 404 de ruta inexistente), pintar `T.servidorViejo` en `text-caption` con el rol información y sin «Reintentar» ni cifra; prueba de Playwright con 501 en las dos operaciones.

### H-2 (media) · axe serio `scrollable-region-focusable` en el tablero (4 combinaciones)
- **Criterio**: U5-06 (accesibilidad de la lista visual) e inconsistencia con el artefacto `docs/design/capturas-u5/axe/P-01-*.json` (0 serios); U5-03 pide 0 serios y fallará en la próxima captura de P-01.
- **Qué pasa**: 6 nodos por combinación: el envoltorio de la tabla oculta `div.overflow-x-auto.sr-only` de `FiguraBarras` (nuevo, `Tablero.tsx`: G-2, G-3) y de `FiguraGrafica` (`GraficaInicio.tsx`: G-1, G-6, G-4, G-5). `overflow-x-auto` gana sobre el `overflow: hidden` de `sr-only` y deja una región desplazable sin foco. Antes el Resumen no tenía figuras; el tablero las trae al primer plano.
- **Corrección**: en los dos componentes, `className={comoTabla ? 'overflow-x-auto' : 'sr-only'}` (sin desplazamiento cuando está oculta) y, cuando se ve, `tabIndex={0}` con `role="region"` y `aria-label` del título; volver a correr axe en las 4 combinaciones.

### H-3 (baja) · Nombres accesibles incompletos
- **Criterio**: DESIGN §7.6 Accesibilidad («nombre completo («Gasto del mes, $8,240.10, 12 % más que septiembre»)»); U3 K-3 y G-3 (nombre accesible de las insignias).
- **Qué pasa**: el `aria-label` de la tarjeta-enlace (`Resumen.tsx`, `Cargada`) une título, valor, etiqueta, estado y acción: K-1 no lleva la comparación ni el % del plan; K-3 no lleva «Disponible hoy … Pendiente confirmar». En G-3, el `aria-label` «Ver gastos de [categoría]» de cada fila tapa la etiqueta, «Rebasada» / «Vas en» y la sugerencia (siguen en la tabla oculta).
- **Corrección**: sumar al nombre la línea de comparación (K-1) y del disponible (K-3), o usar `aria-describedby`; en G-3, `aria-label` = «[etiqueta]. [estado]. [sugerencia]. Ver gastos de [categoría]».

### H-4 (baja) · K-1 rebasado dice «te quedan —»
- **Criterio**: U3 2.0 K-1 («si `remaining_cents` < 0: «−[monto] rebasado»»).
- **Qué pasa**: `Tablero.tsx` K1 pinta «[n] % de [plan] · te quedan —» y debajo «−[monto] rebasado».
- **Corrección**: con rebasado, «[n] % de [plan]» sin «te quedan» y la línea «−[monto] rebasado».

### H-5 (baja) · Marcador de G-6
- **Criterio**: DESIGN §8.1 (marcador 10 px con anillo de 2 px de `ui.surface.card`).
- **Qué pasa**: `symbolSize: 8` sin anillo en las dos series de G-6.
- **Corrección**: `symbolSize: 10`, `itemStyle: { color: c1, borderColor: token('--ui-surface-card'), borderWidth: 2 }`.

### H-6 (baja) · G-3 con `getAlertSettings` caído dice «arriba de 0 %»
- **Criterio**: U3 G-3 (umbrales de `alerts.*`, ninguno constante).
- **Qué pasa**: `T.g3.resumen(…, w ?? 0, …)` en `Tablero.tsx` escribe 0 % si no llegaron los umbrales.
- **Corrección**: sin `w`, resumen sin la cláusula del umbral (o «—»).

### H-7 (baja) · Mocks y reporte
- **Criterio**: inconsistencia menor entre artefactos.
- **Qué pasa**: `frontend/src/pantallas/inicio/mocks.ts` dice «coherentes entre sí», pero `ratio_total_pct` = n + 2 da 8 y 9 en noviembre y diciembre cuando sus pagos dan 10 y 11 (194,750 y 223,243 de 2,000,000, ROUND_HALF_UP), y `payment_count` es 1 donde U3 B-1 dice 2 pagos. Solo `npm run dev:mocks`, marcado SINTÉTICO. El reporte `front-tablero.md` remite a «(ver solicitudes)» y no tiene esa sección.
- **Corrección**: escribir los valores literales (10, 11 y 2 pagos) y agregar la sección de solicitudes al reporte o quitar la referencia.

## 8. PREGUNTA / PROPUESTA
- **PREGUNTA** (puerta 4, sin hallazgo): §2.3 «Qué miro para prevenir» queda solo en U3; si Axel la quiere visible en Inicio, hace falta un bloque y un lugar en §2.4.
- **PROPUESTA**: los 7 puntos de `WF3c-tablero-interfaz` y los 4 iconos de `WF3c-barra-iconos-seccion` están bien marcados PROVISIONAL; sin hallazgo.
