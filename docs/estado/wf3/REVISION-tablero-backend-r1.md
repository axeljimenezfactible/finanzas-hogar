# REVISIÓN · tablero (backend) · ronda 1 de 3 · WF-3c · 12-oct-2026

Revisor independiente `revisor-tablero-backend-r1`. Evalúa el trabajo de `api-tablero` (reporte `docs/estado/wf3/api-tablero.md`) **solo** contra `docs/criterios/S6.md` (S6-01 a S6-04, S6-14 y S6-15), S2.md y S3.md (cifras del motor) y C1.md. No corrige nada. Todas las verificaciones se reejecutaron con comandos propios; los temporales quedaron en `datos/tmp/revisor-tablero-backend-r1/` y se borraron al terminar.

**Dictamen: NO APROBADO.** 0 altos, 2 medios y 4 bajos. Todas las cifras del motor coinciden con el oráculo y no hay verificaciones rotas. Los dos medios son reglas de RN-30 y RN-31 que el código no cumple fuera del caso por omisión.

## 1. Verificación por criterio

| Criterio | Resultado | Evidencia propia |
|---|---|---|
| S6-01 | CUMPLE | `make verify-wf3-go` EXIT 0: go1.27.2, go.mod 1.27; build y vet sin errores; `gofmt -l` vacío; `go mod tidy -diff` sin cambios; 33 paquetes OK y 5 sin pruebas. |
| S6-02 | CUMPLE | `make verify-wf3-cobertura` EXIT 0. Con umbral K25: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8. Sin umbral: api 82.0, auth 83.6, jobs 93.0, reportes 77.3 y telegram 80.0. Con `go test -cover`: api/tablero 87.6 %, api/gastos-inicio 82.0 %, cmd/finanzas 26.2 % e internal/api 88.8 %. Son las mismas cifras del reporte. |
| S6-03 | CUMPLE | `make verify-wf3-contrato` EXIT 0: «los 22 archivos coinciden byte a byte» y «221 operaciones en 176 rutas, 0 diferencias». Además `TestAreas_SinOperacionesEn501` da «221 de 221» y `TestContrato_Operaciones` da «221 de 221 operaciones cumplen»; las 5 operaciones nuevas y getDashboardMonthVsPlan pasan. La tabla operación → prueba del reporte está completa. |
| S6-04 | CUMPLE con observaciones (B-1, B-2 y B-4) | Las respuestas reales siguen la convención de C1: mes `AAAA-MM`, centavos enteros y ErrorBody con `code`, `message`, `fields` y `request_id`. 400 de formato y 422 de regla. Los mensajes están en español de México. |
| S6-14 | CUMPLE | `ss -ltnp :8443` muestra solo `127.0.0.1:8443` y `192.168.50.56:8443`. Por la IP de Tailscale (100.115.28.35) curl sale con código 7, conexión rechazada. HTTP/2: `%{http_version}` = 2 en loopback y en LAN, y `openssl s_client -alpn h2,http/1.1` da «ALPN protocol: h2». Con `--http1.1` responde 200 en HTTP/1.1. `TestEscucha_NegociaHTTP2`, `TestEscucha_SoloLANYLocal` y `TestEscucha_HTTPSConCertificado` pasan. |
| S6-15 | CUMPLE con observación (B-2) | Los identificadores están en inglés (`getDashboardDebtCapacity`, `incomeOf`, `scaleMaxPct`, `MonthHasData`). Los mensajes al usuario están en español de México. En `datos/dev/run/finanzas.log` (221 líneas) hay 0 coincidencias de cookie, password, contraseña o token. |
| S2 y S3 (cifras) | CUMPLE en el caso sembrado; ver M-2 | Ver §2. |
| C1 | CUMPLE con M-1 y B-1 | Esquemas y rutas iguales al contrato. Las reglas RN-31 y RN-30 se cumplen solo en parte (M-1 y M-2). |
| `make verify-all GOOSE=…` | EXIT 0 | «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». |

## 2. Cifras contra el oráculo y los motores (servidor de desarrollo, HTTP/2)

Peticiones con curl, la CA de mkcert y una sesión propia de axel. La contraseña de desarrollo se leyó del archivo sin imprimirse y la sesión se cerró al final (logout 204, luego 401). El binario `backend/bin/finanzas` (15:08) es posterior a todos los fuentes tocados.

| # | Valor | API | Oráculo o motor | ¿Coincide? |
|---|---|---|---|---|
| 1 | Pagos del hogar, oct, nov y dic-2026 (`payments_cents`) | 59,882 · 119,764 · 148,257 | `model.json periodos[0..2].efectivo.axel` = 598.82 · 1,197.64 · 1,482.57 | Sí |
| 2 | Pagos con Karen, oct, nov y dic (`payments_total_cents`) | 97,375 · 194,750 · 223,243 | `periodos[0..2].efectivo.total` = 973.75 · 1,947.50 · 2,232.43 | Sí |
| 3 | Devengado y efectivo de oct en M-21 (`accrued_cents`, `cash_cents`) | 119,764 · 59,882 | `periodos[0].devengado.axel` 1,197.64 y `efectivo.axel` 598.82. Además `getDashboardMonthVsPlan` da nov 119,764 y dic 176,751 (= 1,197.64 y 1,767.51 del oráculo) | Sí |
| 4 | Ratio oct, nov y dic (`ratio_pct`, `ratio_total_pct`) | 2/4 · 5/8 · 6/9 | Recalculado con ROUND_HALF_UP: 59,882 ÷ 2,450,000 = 2.44 → 2; 97,375 ÷ 2,450,000 = 3.97 → 4; 4.89 → 5; 7.95 → 8; 6.05 → 6; 9.11 → 9. Ingreso 2,450,000 = `getIncomeMonth.total_cents` | Sí |
| 5 | Margen (`headroom_cents`) | 920,118 · 860,236 · 831,743 | 40 × 2,450,000 ÷ 100 − pagos = 980,000 − 59,882, − 119,764 y − 148,257 | Sí |
| 6 | Plan de oct (`planned_cents`) | 200,000 | `getBudgetMonth.planned_spending_cents` 200,000 y `getDashboardMonthVsPlan.planned_cents` 200,000 | Sí |
| 7 | Diferencias de M-01 en dic (campos de v1.5) | 56,987 · 28,493 | 176,751 − 119,764 y 148,257 − 119,764 | Sí |

**Umbrales desde `app_setting` (d).**
- Con `PATCH` 1/2 → 200:
  - octubre pasa a `warning`, con escala 10 y margen −10,882 (= 49,000 − 59,882);
  - noviembre pasa a `exceeded`, con margen −70,764;
  - monthly-overview da `warning`, `exceeded` y `exceeded`.
- Amarillo 45 con tope 2 → 422 `out_of_range`, y no se guardó.
- `"4"` → 400 (ver B-1).
- `reset_parameters` → 200 y `is_initial` en true; octubre vuelve a `normal`, 30/40 y escala 60. La base de desarrollo quedó en sus valores iniciales; la bitácora conserva los cambios.

**Sin cifras a mano en producción (e).**
- Ingreso, devengado, efectivo, plan y % usado salen de `motor_fin` (FortnightFlow, PlanVsActual y BudgetMonth).
- Los pagos salen de `motor_cred.UnifiedCalendar` (por `creditos.PaymentsByMonth`).
- El manejador solo divide con `dominio.PercentOf` y `MulDivRound` (ROUND_HALF_UP), como piden RN-30 y M-20.
- `scaleMaxPct` cumple «múltiplo de 10 ≥ 1.5 × tope»: lo comprobé con 1, 33, 34, 40 y 45.

## 3. Hallazgos

### M-1 (media) · `getBudgetSignals` cuenta meses que no están cerrados cuando `period` es posterior al mes en curso
- **Dónde:** `backend/internal/api/tablero/tablero.go`, líneas 284-299. La ventana es `m.Add(-i)` sin tope en el último mes cerrado.
- **Qué dice la regla:** RN-31 y el esquema `BudgetSignals` de C1 v1.5: «se miran los `budget.signal_months` meses **cerrados** anteriores al mes pedido; **el mes en curso no cuenta**».
- **Comprobado en el servidor de desarrollo (hoy es 12-oct-2026):**
  - `GET /api/v1/budget/signals?period=2026-11` mira 2026-08, 2026-09 y **2026-10**, que es el mes en curso: la categoría 1 aparece con 60 %.
  - `?period=2027-01` mira 2026-10 (en curso), 2026-11 y 2026-12 (proyección). Devuelve **`"signal":"cut"`** para la categoría 1 (60 %, 60 % y 88 %).
  - `?period=2027-02` mira solo meses de proyección y también devuelve `cut`.
  - Es decir: sugiere «recortar presupuesto» con gasto que todavía no ocurre. No hay prueba con `period` futuro.
- **Qué hacer:**
  - Cortar la racha en todo mes ≥ el mes en curso, porque no está cerrado. Si se quiere mostrar, que no cuente.
  - O bien responder 422 a un `period` posterior al mes siguiente al en curso.
  - Cualquiera de las dos se marca PROVISIONAL en el bloqueo `C1-tablero` (punto 5, formato de `period`) o en un bloqueo nuevo `S6-tablero-senal-periodo.md`, con sus opciones.
  - Agregar una prueba con `period` futuro (y con `period` = mes siguiente) que espere `none` o 422.

### M-2 (media) · Los pagos de crédito del tablero (S3) dejan fuera las líneas archivadas; el efectivo de S2 no
- **Dónde:** `backend/internal/api/creditos/tablero.go`, líneas 26-31. `PaymentsByMonth` usa `s.states(nil)`, que filtra `c.IsActive` (`api/creditos/motor.go`, líneas 205-214).
- **Por qué es un problema:**
  - `dominio.LoadFinEngine` (S2) carga todos los `scheduled_payment` sin mirar `credit.is_active`.
  - RN-30 (PROVISIONAL) dice «pagos = parte del hogar de los pagos programados con fecha en el mes, **de todas las líneas**, pagados o no».
  - Que se cuenten solo las líneas activas es un supuesto que cambia cifras. No está en §Supuestos ni en un bloqueo, y no tiene prueba.
- **Comprobado con una prueba temporal propia** (`go test -overlay`, sin tocar el árbol). Usa el fixture de `tablero_datos_test.go` y después `UPDATE credit SET is_active = 0`. Octubre, antes y después:
  - `credit_payments_cents`: 650,000 → **0**;
  - `debt_ratio_pct`: 33 → **0**;
  - `has_credits`: true → **false**, que el contrato presenta como «0 % cifra real»;
  - `cash_cents` (S2), en la misma fila: **700,000** antes y después, porque sí incluye esos 6,500.00.
- **Efecto:** cuando Axel archive una tarjeta ya liquidada, su historia desaparece de G-1 y G-6 y del barómetro de los meses pasados, mientras el efectivo sigue contándola.
- **Qué hacer:**
  - Abrir un bloqueo PROVISIONAL, por ejemplo `S6-tablero-lineas-archivadas.md`, con dos opciones: (a) contar todas las líneas, activas o archivadas, como dice RN-30; (b) solo las activas, corrigiendo el texto de RN-30, M-20 y M-21.
  - Aplicar la opción PROVISIONAL.
  - Agregar una prueba con una línea archivada que compare `credit_payments_cents` contra la parte de crédito del efectivo de S2.
  - Si `api/creditos/` queda fuera del ALCANCE de api-tablero, va como solicitud fuera de alcance con esta ruta y este cambio.

### B-1 (baja) · Error de tipo en `updateDashboardSettings`: 400 en la API, 422 en el texto de C1
- `"debt.capacity_warning_pct":"4"` responde 400 `bad_request`. La descripción de la operación y de `DashboardSettingsUpdate` (`docs/contratos/extractos/ajustes-e.yaml`, líneas 323 y 1453) dice «de otro tipo responde 422 validation_failed».
- El agente siguió C1 §3 y CONVENCIONES §4 («el formato es 400») y lo dejó como PREGUNTA 2.
- La misma frase ya existe en 8 operaciones de settings (`ajustes-c`, `-d` y `-e`), así que la inconsistencia es anterior a este trabajo.
- Falta abrir el bloqueo o la solicitud para alinear el texto en el generador (`backend/internal/api/c1gen/gen3_ops.py`).

### B-2 (baja) · Dos mensajes 400 que confunden en estas operaciones
- El 400 de tipo incorrecto dice «Tipo incorrecto: los importes son enteros en centavos y las fechas, texto AAAA-MM-DD.», pero aquí se trata de porcentajes y conteos. Es el mensaje genérico de `reg.DecodeJSON`.
- `GET /api/v1/dashboard/debt-capacity?period=oct` responde «Usa, por ejemplo, 2026-10 o 2026-10-q1.», pero esta operación rechaza `2026-10-q1` con 422. Es el mensaje genérico de `api.QueryPeriod`.
- Sugerencia: un mensaje propio en `monthParam` y en `updateDashboardSettings`, o una solicitud a backend-base.

### B-3 (baja) · `docs/estado/bloqueos/INDICE.md` no lista `S6-tablero-sin-datos`
Tampoco lista `C1-tablero` ni `U3-tablero-*`. El agente lo dejó en Pendientes. Hay que agregar las filas.

### B-4 (baja) · Algunas pruebas de 4xx usan `c.Raw`, que no valida el esquema
- Están en `tablero_test.go` (líneas 92, 162 y 171) y en `tablero_datos_test.go` (líneas 220, 223 y 318).
- El reporte (S6-04) afirma que «cada `c.Do` valida … también los 4xx propios: 400 de formato». No es exacto para estos casos.
- La cobertura de contrato no se pierde, porque `TestContrato_Operaciones` sí valida los 4xx de las 221. Basta con cambiar a `c.Do`.

## 4. Lo que está bien
- HTTP/2 por ALPN con un cambio mínimo y una prueba que falla sin `NextProtos`. La escucha sigue solo en LAN y loopback.
- Las cinco operaciones responden, sus cifras salen de S2 y S3 y coinciden con el oráculo en todos los valores que recalculé. Los umbrales se leen de `app_setting` en cada petición.
- La regla «sin datos» quedó en una sola función PROVISIONAL, con su bloqueo y sus tres opciones.
- Las pruebas con fixture cubren las fronteras de redondeo de RN-30 y RN-31, el cambio de color por ajustes, la bitácora y el 422 sin guardar.

## 5. Para la ronda 2
Cerrar M-1 y M-2 (código o bloqueo PROVISIONAL, más una prueba cada uno). B-1 a B-4 son opcionales.
