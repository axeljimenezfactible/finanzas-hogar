# Reporte · `revisor-c3-tablero` (WF-3c, ronda 1)

Fecha: 2026-10-12 · Agente: `revisor-c3-tablero` (revisor independiente; no redacté la matriz ni implementé nada) · Base: HEAD `5a089da` con el árbol de trabajo de WF-3c · Objeto: `docs/contratos/C3-integracion.md` v1.3 (SHA-256 `8b03300d85c0eb775af714f2c49df91af698b2750b08ae0ccc6a1043122d691b`, sin cambios antes y después) · Criterios: `docs/criterios/C3.md` 1.1. Sin `git add` ni commit. Sin Telegram.

## Hecho
- Revisé C3 v1.3 contra `C3.md` completo con mis propios scripts y ejecuciones. **Veredicto: DEVUELTO, sin firma.** 4 altos (H-1 C3-08 flujos F-1, F-3 y F-4; H-2 `c3-09.spec.ts` falla; H-3 `updateDashboardSettings` sin pantalla, C3-03; H-4 deuda viva de K-3 y P-17 distinta del oráculo y de `getLiveDebt`), 1 medio (H-5) y 4 bajos (H-6 a H-9). Detalle: `docs/estado/wf3/REVISION-c3-tablero.md`.
- No agregué bloque de firma (C3-07) y no modifiqué la matriz, las pruebas ni el código.

## Archivos
- Nuevos: `docs/estado/wf3/REVISION-c3-tablero.md`, `docs/estado/wf3/revisor-c3-tablero.md`.
- Temporales (borrados al terminar): `datos/tmp/revc3/` (registros, script propio `rev.py`, sonda Playwright de solo lectura), `datos/tmp/u5-flujos/` y `datos/tmp/playwright-u5-flujos/` (copias que crean `servidor-prueba.sh` y la configuración de flujos).
- Respaldados y restaurados sin cambios: `frontend/pruebas-u5/flujos/resultados/` y `frontend/pruebas-u5/flujos/capturas/` (las pruebas de flujos los reescriben).

## Verificación
- **C3-01 · cumple.** Script propio: T de U1 = 41 (27 + 14); `grep -c "^| P-"` = 41; 42 filas con Global; §1 define los 3 estados.
- **C3-02 · cumple.** Script propio sobre U1 §2.1, §2.3 y §3.29, `openapi.yaml` y la matriz: 0 pantallas sin fila, 211 `N-xxx` (3 Local) y 0 sin operación, 0 `N-xxx` de C1 sin fila, 221 operaciones y 5 sin pantalla con excepción justificada, 0 filas con operaciones distintas a las derivadas. `c3-integracion.py`: «Matriz del documento = matriz generada: sí». Baja H-8 (P-01 sin `getCreditTotals` ni `getDashboardSettings` porque U1 no se actualizó con U3 2.0).
- **C3-03 · no cumple.** `cd frontend && npm run build` → EXIT 0. Rutas: 214 invocadas, 0 fuera de `openapi.yaml`, 0 sin registrar (221 en `routes.go` = `openapi.yaml`). `python3 -I frontend/scripts/c3-integracion.py --pw-json datos/tmp/revc3/pw.json --go-log datos/tmp/revc3/go.log` → **EXIT 1**: «Operaciones con pantalla que el frontend no llama: 2 ['listPeriods', 'updateDashboardSettings']; sin excepción justificada o sin cobertura: 1 ['updateDashboardSettings']» (H-3).
- **C3-04 · cumple.** `cd backend && go test ./... -count=1 -v` → 33 `ok`, 0 `FAIL`, 1147 `--- PASS`, 1 `--- SKIP`, 221/221 `TestContrato_Operaciones/*` PASS. `cd frontend && npx playwright test --reporter=json` → 263: 259 pasan, 4 omitidas, 0 fallan (348 s). 60 pruebas citadas visibles: todas pasan. Baja H-6 (23 filas con lista truncada).
- **C3-05 · cumple.** 164/164 PNG citados existen, son PNG y miden 412 o 1440 px de ancho según dispositivo; 164 en la raíz de `docs/design/capturas-u5/`.
- **C3-06 · cumple en la forma.** `grep -c "| no |"` = 1 (Global, bloqueo `C3-capturas-global.md` citado).
- **C3-07 · no se firma** (devuelto). `sha256sum` de la matriz igual antes y después.
- **C3-08 · no cumple.** `bash frontend/pruebas-u5/flujos/servidor-prueba.sh arrancar` y `FH_BASE=https://127.0.0.1:8444 npx playwright test -c pruebas-u5/flujos/playwright.config.ts <archivo>`, una copia nueva por archivo: `flujos.spec` 4 passed y 3 failed (F-1 `:13` partida con `tarjeta-1`; F-3 `:122` toque 1 de la ruta A lleva a «Productos y precios»; F-4 `:222` «tarjeta Próximo pago»: `locator.tap` agota 15 s); `pendientes.spec` 3 passed y 1 failed (F-4b `:57`, mismo selector); `f7-f8.spec` 2 passed. F-3A3 (`:117`, ruta A con la condición (c)) **pasa** en 3 toques, contra lo que dice la matriz (H-5). Base de desarrollo: 6 gastos, 4 compras a crédito y 0 pesos, antes y después.
- **C3-09 · no cumple.** `c3-09.spec.ts` del repositorio falla en la línea 31 (`tarjeta-3`). Sonda propia de solo lectura en `https://finanzas.local:8443` (`NODE_EXTRA_CA_CERTS` con la CA de mkcert, sesión de `axel`, contraseña leída del archivo de desarrollo sin imprimirla): 57 cifras, 55 iguales a `model.json`: A $10,000.00, B $6,905.50, mínimos $8,094.50 y $5,000.00, próximo pago $973.75 del 30-oct-2026 con Hogar $598.82 y Karen $374.93, deuda viva al 2027-04-10 $6,293.44 (P-07 y `getLiveDebt`), P1–P6 (42 cifras) y el barómetro (2 % de 40 %, pagos $598.82, margen $9,201.18; fórmula de U3 §B-1 recalculada con ROUND_HALF_UP). **2 diferencias:** K-3 y P-17 «deuda viva total» $18,999.22 (`getCreditTotals`) contra $3,895.00 del oráculo a hoy, de `getLiveDebt?date=2026-10-12` y de P-07 (H-4).
- **C3-10 · cumple.** 12 filas leídas: español de México, fechas de evidencia ISO.
- **C3-11 · cumple.** `grep -n "v1"` → «Versión v1.3 · 2026-10-12»; `c3-tablero.md` con §Supuestos y §Bloqueos; `ls docs/estado/bloqueos/C3-*` → 5 archivos.

## Supuestos
- La deuda viva «a hoy» del oráculo es el total de la compra 1 (`despensas[0].total` = 3,895.00), porque al 2026-10-12 es la única compra con fecha pasada y su primer pago es el 2026-10-30. La propia API (`getLiveDebt` con fecha de hoy) da lo mismo. Si Axel decide que las compras futuras sembradas cuentan como deuda viva de hoy, H-4 cambia de dueño (oráculo y criterio), no de severidad.
- F-3 ruta A se evalúa con la condición (c) PROVISIONAL del orquestador (`C3-flujos` punto 2): su evidencia es F-3A3, que pasa. Aun así `flujos.spec.ts:122` es una prueba del repositorio que falla, y C3-08 pide que pasen las pruebas de los flujos.
- «Probado» se evalúa con la regla PROVISIONAL de `C3-matriz` punto 4 (título de la prueba y API simulada).

## Bloqueos
- Sin bloqueos nuevos. Siguen abiertos y los cito: `docs/estado/bloqueos/C3-ajustes-tablero.md` (PREGUNTA, H-3), `C3-flujos.md` punto 2 (PROVISIONAL (c), la ratifica Axel), `C3-capturas-global.md`, `C3-matriz.md`.

## Pendientes
- H-1 y H-2: corregir `flujos.spec.ts`, `pendientes.spec.ts` y `c3-09.spec.ts` para U3 2.0 (dueño `u5-flujos` u otro que designe el orquestador) y reejecutar.
- H-3: decisión sobre `updateDashboardSettings`.
- H-4: respuesta sobre la deuda viva de `getCreditTotals` y, si se corrige, reejecutar C3-09.
- H-5 a H-9 en la siguiente regeneración de C3; luego nueva revisión y firma.

## Propuestas para Axel
| Opción | Resultado | Costo | Riesgo |
|---|---|---|---|
| A. PROPUESTA (la recomendada): corregir H-1 a H-4 (pruebas de flujos y de cifras para el nuevo tablero, apartado «Parámetros del tablero», deuda viva a la fecha de hoy) y pedir una ronda 2 de revisión | C3 firmado con evidencia del repositorio | Medio: 3 pruebas, 1 sección de Ajustes y 1 cálculo del servidor | Bajo |
| B. Firmar con H-3 y H-4 como excepciones con fecha | Se llega antes a la puerta 4 | Bajo ahora | Alto: Inicio y Créditos muestran una deuda 4.9 veces mayor que la real a hoy |

## Solicitudes fuera de alcance
- `frontend/pruebas-u5/flujos/flujos.spec.ts:18-19,41-42` y `:226`, `pendientes.spec.ts:62`, `c3-09.spec.ts:31-37`: selectores de U3 2.0 (H-1, H-2).
- Servidor de `getCreditTotals` (`live_debt_cents`): misma regla de fecha que `getLiveDebt` (H-4, tras la respuesta).
- `docs/design/U1-arquitectura.md` §3.3 y §3.29: N-191 y N-275 en la reutilización de P-01 (H-8).

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Acciones: 45. KB leídos: ~120 (C3 completo ~59, C3.md 11, bloqueos ~8, U1 y U3 por secciones ~10, pruebas y registros ~25, salidas de scripts ~10). KB escritos: ~22 (2 reportes ~19, script y sonda temporales ~6).
