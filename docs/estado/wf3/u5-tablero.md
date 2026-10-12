# WF-3a · u5-tablero · Reporte

Agente `u5-tablero` (Sonnet) · 12-oct-2026 · HEAD `5a089da`. Criterios: U5-01, U5-03, U5-04, U5-10 (y U5-02 sobre el tablero). **No corregí nada de la aplicación ni de los motores**: las diferencias van abajo con regla, nodo y pantalla.

## Hecho
- `npm run build` (frontend) y `go build ./...` + `go build -o backend/bin/finanzas ./cmd/finanzas`: sin errores. Servidor reiniciado con `ops/dev-stop.sh` y `ops/dev-run.sh` (`/health` ok en https://finanzas.local:8443).
- Regeneré TODAS las capturas con `npx playwright test -c playwright.u5.config.ts` desde `frontend/`: **174 passed (9.0 min)** (164 capturas + axe + contraste + tema + Tendencias), con la barra lateral nueva (iconos de sección visibles en la captura de escritorio de P-01).
- Nuevo `frontend/pruebas-u5/subsecciones.spec.ts`: P-01 en sus 4 subsecciones (Resumen, Detalle, Tendencias, Calendario) × 4 combinaciones con captura, axe, contraste, tema, hosts y desborde. Salida en `docs/design/capturas-u5/evidencia/inicio-subsecciones/` (16 PNG + 16 JSON; subcarpeta, no cuenta en U5-01).
- Nuevo `frontend/pruebas-u5/tablero-cifras.spec.ts`: guarda el texto visible del Resumen y las respuestas crudas de `/api/v1/dashboard/*` de la misma carga en `evidencia/tablero-cifras/<combinación>.json` (4 archivos).
- Lectura de cifras en pantalla contra API y oráculo (abajo). Temporales (`datos/tmp/arc.mjs`, `desb.mjs`, `api.mjs`, config temporal, logs y `playwright-u5-*`) borrados.

## Archivos
- Regenerados: `docs/design/capturas-u5/*.png` (164), `axe/*.json` (164), `contraste/*.json` (164), `hosts.json`, `evidencia/manifiesto.json`, `evidencia/cifras-oraculo.json`, `evidencia/tema-*.json`, `evidencia/tendencias/*`.
- Nuevos: `frontend/pruebas-u5/subsecciones.spec.ts`, `frontend/pruebas-u5/tablero-cifras.spec.ts`, `docs/design/capturas-u5/evidencia/inicio-subsecciones/*`, `docs/design/capturas-u5/evidencia/tablero-cifras/*`, este reporte.
- Nota: `frontend/playwright.u5.config.ts` no incluye aún los dos specs nuevos (su `testMatch` es solo de los viejos); ver solicitudes fuera de alcance. Los corrí con una config temporal ya borrada.

## Verificación
- **U5-01** · OK: `verificar_capturas.py` → `T = 41 · 4T = 164`; `PNG correctos: 164 de 164 · faltantes: 0 · mal: 0`; `PNG en la carpeta: 164` (el desborde de 3 px de P-26.13 Android de la ronda anterior ya no existe).
- **U5-03** · OK en las 164: `axe-core: 164 JSON, violaciones serias o críticas = 0`. P-01 en las 4 subsecciones × 4 combinaciones: 0 violaciones (16 de 16).
- **U5-04** · FALLA en lectura literal: `24466 pares evaluados, bajo umbral = 2` (P-01 android-oscuro y P-01 escritorio-oscuro, 1 cada una; las otras 162 en 0). Es el mismo par: `a[data-prueba="tarjeta-B-1"] … div[data-prueba="barometro"] > svg.block.h-auto`, tipo `icono`, fg `#000000` sobre `#36363A` = 1.746 (umbral 3). Diagnóstico (medido con `getComputedStyle`): el `<svg>` raíz del barómetro es `aria-hidden`, no tiene `stroke` y su `fill` es el inicial negro, que no se pinta; `contraste-dom.js` (líneas 91-96) toma ese `fill` como color del ícono. Los trazos reales sí pasan: arcos verde/ámbar/rojo 4.74/5.13/5.40 (claro) y 6.42/6.04/5.82 (oscuro) sobre la tarjeta; aguja 10.76; marcas 3.47 y 6.09; solo la pista gris del arco es 1.12 (claro) y 1.45 (oscuro), decorativa (la información la lleva el arco de valor, el dígito y la insignia). En claro el par da negro sobre claro y pasa por casualidad. **Falso positivo de la medición, no defecto de la app**; no cambié el script.
- **U5-04 en subsecciones** · mismo par en Resumen oscuro (2 de 16); Detalle, Tendencias y Calendario en 0 bajo umbral.
- **U5-10** · OK: en las 164 capturas y en las 16 de subsecciones el `background-color` del `body` = token de C2 y sin `data-theme` (`tema.spec.ts` 2 pruebas pasan).
- **U5-08** · OK: `hosts.json` = `finanzas.local:8443` en 164 capturas; 0 hosts externos. **U5-11** · OK: 41 filas P- con sus 4 rutas.
- `make verify-u5` termina con error solo por el punto U5-04 de arriba (los demás U5-01, 03, 08, 11 en verde).
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` · exit 0 (WF-1 completo OK).
- **U5-02 · cifras del tablero (P-01 Resumen) contra API y oráculo** (4 combinaciones, mismos 61 importes y mismas respuestas):
  - Próximo pago: pantalla «Vie 30 oct 2026 · Kueski · $973.75 · Hogar $598.82 · Karen $374.93»; API `next-payment` 97375 / 59882 / 37493; oráculo P1 total 973.75, axel 598.82, Karen 374.93. **Coincide.**
  - Barómetro: pantalla «2 %», «de 40 % de tope», «Verde», «Con la parte de Karen: 4 %», «Te caben $9,201.18»; API `ratio_pct` 2, `ratio_total_pct` 4, `headroom_cents` 920118; RN-30: 598.82/24,500 = 2.44 → 2; 973.75/24,500 = 3.97 → 4; 40 % × 24,500 − 598.82 = 9,201.18. Serie mensual nov 5 %, dic 6 %, ene a mar 7 %: coincide con oráculo (1,197.64, 1,482.57, 1,767.50/1,767.52). **Coincide** (el oráculo no trae ingresos; el 24,500 es del plan sembrado).
  - Disponible hoy $6,905.50 (Lectura B) = oráculo `disponible_B`; Gasto del mes $1,197.64 y efectivo $598.82 = oráculo P1 devengado/efectivo axel. **Coinciden.**
  - Devengado contra efectivo: oct-2026 $1,947.50/$973.75, dic $2,517.37/$2,232.43, acumulado a mar-2027 $14,362.48 devengado y $13,103.78 efectivo = `horizonte.total_devengado` y `total_efectivo`. **Coinciden.**
  - Deuda viva proyectada: mar-2027 $6,293.44 = `deuda_viva_10_abr` (y `/debt/live?date=2027-04-10` = 629344). **Coincide.**
  - **DIFERENCIA 1 (media): dos «deuda viva» distintas en la misma pantalla y fecha.** K-3 Créditos muestra «$18,999.22 deuda viva total» (API `/credits/totals` `live_debt_cents` 1899922 = 3,895.00 + 3 × 5,034.74, es decir la despensa 1 más las despensas 2, 3 y 4 que la semilla marca «aún no se compra»), mientras «Deuda viva por línea» muestra «Deuda viva hoy: $3,895.00» (`/debt/live?date=2026-10-12` = 389500, igual a `pagos_despensa_1` del oráculo). Con el oráculo, la deuda viva al 12-oct-2026 es 3,895.00 y 18,999.22 no existe en él.
  - **DIFERENCIA 2 (media): «Deuda viva proyectada» desde oct-2026.** La primera fila es $18,025.47 (= 18,999.22 − 973.75) y baja a $16,077.97, $13,845.54…; con el oráculo, al cierre de oct-2026 son 2,921.25 (3,895.00 − 973.75) hasta la compra de la despensa 2 (10-dic-2026). La serie solo coincide con el oráculo desde que todas las compras ya ocurrieron (mar/abr-2027). Mismo origen que la diferencia 1: se cuentan compras futuras de la semilla (bloqueo existente `S6-semilla-despensas-futuras`).
  - Texto: el subtítulo de «Deuda viva proyectada» dice «Últimos 12 meses» y los meses son futuros (oct-2026 a sep-2027): debería decir «Próximos». Baja.
  - Visual (en la captura escritorio oscuro): en «Tendencia del barómetro» la etiqueta de serie «Capacidad de endeudamiento» pisa la etiqueta del eje «30 %»; `tendencias.spec.ts` solo mide P-01 Tendencias y no lo detecta en el Resumen. Baja.

## Hallazgos de la aplicación (sin corregir)
1. **P-01 Detalle en Android desborda en horizontal: 475 px en vez de 412** (claro y oscuro). Causa: `span.inline-flex.h-6.w-max` (insignia de `pantallas/inicio/comun.tsx:25`) con el texto «Reposición de capital pendiente de confirmar con Kueski» en una sola línea de 442 px; arrastra la barra inferior fija a 475. Severidad alta (U5-06: sin desborde horizontal); no la ve U5-01 porque la captura de P-01 es el Resumen.
2. Diferencias 1 y 2 y los dos textos/visual bajos de la sección U5-02.

## Supuestos
1. P-01 «4 subsecciones» se capturan en `evidencia/inicio-subsecciones/` y no en la raíz de `capturas-u5/` para no alterar el conteo de U5-01 (164).
2. Fecha del sistema 12-oct-2026: «hoy» = 12-oct (la semilla de desarrollo lleva compras de oráculo en fechas futuras).
3. Los falsos positivos de contraste se determinan midiendo `getComputedStyle` de los elementos del barómetro, sin tocar `contraste-dom.js`.
4. La contraseña de desarrollo se leyó de `~/.config/finanzas-hogar/dev-semilla.txt` en ejecución y no se imprimió ni se guardó.

## Bloqueos
Ninguno nuevo. Referencia: `S6-semilla-despensas-futuras` (diferencias 1 y 2); `U5-contraste-ui-controles` (sin cambios).

## Pendientes
- Decisión sobre U5-04: corregir `contraste-dom.js` para ignorar el `fill` inicial del `<svg>` raíz cuando sus hijos pintan con su propio `stroke` (y medir esos hijos), o marcar el barómetro como excepción; luego repetir `--grep contraste`.
- Corregir el desborde de la insignia en P-01 Detalle Android y decidir cuál es la «deuda viva» canónica de K-3 (WF-3c).
- Añadir los dos specs nuevos al `testMatch` de `playwright.u5.config.ts` (fuera de alcance).

## Propuestas para Axel (PROPUESTA, no decidido)
| Qué | Resultado | Costo | Riesgo |
|---|---|---|---|
| K-3 y «Deuda viva proyectada» solo con compras hechas a la fecha (como `/debt/live`) | Una sola cifra de deuda viva ($3,895.00 hoy), alineada al oráculo | Backend: `credits/totals` y `monthly-overview` con fecha | Bajo; hay que decidir si «deuda viva total» incluye lo programado |
| Permitir salto de línea en la insignia (`w-max` → `max-w-full` con texto envolvente) | P-01 Detalle sin desborde en teléfono | Una clase | Muy bajo |
| Ignorar el `<svg>` raíz sin trazo propio en el script de contraste | U5-04 en 0 sin excepciones | Un ajuste de script | Bajo |

## Paquetes instalados
Ninguno.

## Consumo aproximado
Acciones ≈ 40 · KB leídos ≈ 190 · KB escritos ≈ 12 (más ≈ 25 MB de evidencia regenerada por el script).
