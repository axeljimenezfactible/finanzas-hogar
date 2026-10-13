# correccion-backend-cierre-r1 · WF-3c bis · 12-oct-2026

Agente `correccion-backend-cierre-r1` (Opus). HEAD `0197dad`, sobre el diff sin commit de `backend-cierre`. Corrige los tres medios de `docs/estado/wf3/REVISION-backend-cierre-r1.md` (M-1, M-2 y M-3) y, de paso, los bajos B-2, B-3 y B-4. No toqué motores ni cifras de motor, ni el frontend, ni git, ni el servidor de desarrollo.

## Hecho
- **HECHO · M-1 (salida ii de la revisión, PROVISIONAL).** `states` (`backend/internal/api/creditos/motor.go`) ya no filtra `credit.is_active`: `getCreditTotals` (K-3, P-17) y `getDebtTrend` sin `credit_id` (G-5) cuentan todas las líneas, activas o archivadas, y con `credit_id` cuentan las líneas nombradas aunque estén archivadas. Así la deuda viva de hoy es la de `getLiveDebt` hoy, cada cierre de G-5 es la de `getLiveDebt` ese día y los pagos del hogar del mes de K-3 son los del barómetro. Es lo que dicen RN-21 (FIRME), RN-30, U3 K-3 («de todas las líneas»), C1 (`DashboardDebtCapacity.payments_cents` = `CreditTotals.month_household_cents`) y las decisiones 1 y 2 del orquestador. La pregunta para Axel queda en el **punto 5** (nuevo) de `S6-tablero-lineas-archivadas`. Ninguna cifra del motor cambia: solo qué líneas se le pasan.
- **HECHO · M-2 (C1-15).** `month_pending_household_cents` (`gen2_schemas.py`) dice su regla y la marca: «…su parte del hogar menos lo que ya abonaron personas del hogar a ese pago, nunca negativa; los abonos de Karen no la cubren (RN-04); un pago de corte de tarjeta pendiente cuenta su parte del hogar completa (PROVISIONAL, bloqueo S6-deuda-viva-hoy, punto 3)». También marqué lo que aplica M-1: `CreditTotals.live_debt_cents`, `CreditTotals.month_total_cents` y `DebtTrend.months[].live_debt_cents` citan `S6-deuda-viva-hoy` (puntos 1 y 2) y `S6-tablero-lineas-archivadas` (punto 5), y lo mismo hacen M-17 y M-19 de `c1_plantilla.md`. La sección v1.5.1 tiene una viñeta con lo de esta corrección. Regeneré con `build.py`. `S6-deuda-viva-hoy` pasa de 0 citas a 3 en `openapi.yaml` y a 4 en `C1-datos.md`. `info.version` sigue en 1.5.1 y no cambia ningún tipo de `types_gen.go`.
- **HECHO · M-3 (C1-13).** Queda como solicitud fuera de alcance 1: el archivo es `docs/criterios/C1.md`.
- **HECHO · bajos.** B-2: `S6-deuda-viva-hoy.md` ahora nombra `LiveDebtAt`. B-3: en la tabla «Cambios» de `c1_plantilla.md`, la fila 1.5.1 quedó después de la 1.5. B-4: RN-21, `S6-deuda-viva-hoy` (2a) y el comentario de `getDebtTrend` dicen «parte de la deuda de hoy menos los pagos que vencen antes del cierre». B-1: la tabla de operaciones contra pruebas va abajo, en Verificación.
- **HECHO · documentos.** RN-21 (PROVISIONAL): las dos cifras cuentan todas las líneas, con cita del punto 5. RN-30: el paréntesis ya dice que es la misma cifra que `month_household_cents`. También actualicé la nota de corrección del encabezado de `reglas-negocio.md`, `S6-tablero-lineas-archivadas.md` (la pregunta pasó al punto 5 y el punto 3 remite a él), `S6-deuda-viva-hoy.md` (nuevo párrafo «Qué líneas cuentan») y la fila de las archivadas en `INDICE.md`.

## Archivos
- Go: `backend/internal/api/creditos/motor.go` (`states`), `backend/internal/api/creditos/creditos.go` (comentarios de `getCreditTotals` y `getDebtTrend`) y `backend/internal/api/creditos/creditos_test.go` (prueba nueva `TestOp_getCreditTotals_LineaArchivada`).
- Contrato: `backend/internal/api/c1gen/gen2_schemas.py` y `backend/internal/api/c1gen/c1_plantilla.md`. Lo regenerado: `docs/contratos/openapi.yaml`, `docs/contratos/C1-datos.md` y `docs/contratos/extractos/*`. `routes.go` y `types_gen.go` no cambian respecto de `backend-cierre`.
- Documentos: `docs/contratos/reglas-negocio.md`, `docs/estado/bloqueos/S6-tablero-lineas-archivadas.md`, `docs/estado/bloqueos/S6-deuda-viva-hoy.md`, `docs/estado/bloqueos/INDICE.md` y este reporte.

## Verificación
- **M-1, prueba nueva** `TestOp_getCreditTotals_LineaArchivada`. Monta Kueski y una segunda línea igual, cada una con la despensa del 10-oct, y archiva la segunda con `PATCH /credits/{id}` `is_active=false`. Resultados:
  - Archivar no cambia el cuerpo de `getCreditTotals` ni el de `getDebtTrend`.
  - `live_debt_cents` 779000 = `getLiveDebt` del 2026-10-12.
  - `month_total_cents` 194750 (2 × 97375), con `month_household_cents` = `payments_cents` y `month_total_cents` = `payments_total_cents` de `getDashboardDebtCapacity`.
  - `month_payment_count` 2 = `payment_count`, y `month_pending_household_cents` = `month_household_cents` (nada pagado).
  - Los 12 cierres de G-5 = `getLiveDebt` del último día de cada mes; octubre da 584250 = 2 × 292125.
  - Con `credit_id` de la línea archivada, octubre da 292125.
- **Que la prueba sí detecta el defecto:** con `go test -overlay` y el `motor.go` anterior, la prueba FALLA: «archivar cambió las cifras… serie antes 584250 … después 292125».
- `go test -count=1 -cover`: api/creditos **82.5 %** (antes 80.9), api/tablero 87.8, api/presupuesto 77.7 e internal/api 88.8, todo OK. `TestOp_getCreditTotals_DeudaVivaHoy` (389500, 292125, 97375) y `TestOp_Dashboard_LineaArchivada/opcion_a_PROVISIONAL` siguen en verde.
- **S6-01:** `make verify-wf3-go` EXIT 0 (gofmt sin archivos, `mod tidy -diff` sin cambios, 33 paquetes OK y 5 sin pruebas).
- **S6-03:** `make verify-wf3-contrato` EXIT 0: «los 22 archivos coinciden byte a byte con el generador» y «S6-03: 221 operaciones en 176 rutas, 0 diferencias». `python3 -I backend/internal/api/c1gen/build.py --check`: «RESULTADO: los 22 archivos coinciden byte a byte con el generador».
- **C1-15:** `verify_c1.py` da «OK C1-15 toda marca PROVISIONAL remite a un bloqueo abierto». Hay 46 marcas en `openapi.yaml`, ninguna sin bloqueo, y los bloqueos citados incluyen `S6-deuda-viva-hoy` y `S6-tablero-lineas-archivadas`. Sin Make, con `python3 -I`, C1-14 sale FALLA solo porque falta el módulo `jsonschema`. Dentro de `verify-all` pasa.
- **C1-13:** `verify_c1.py` da «OK C1-13 … info.version: 1.5.1». La línea literal del Verifica de `docs/criterios/C1.md` (`version: "1.5.0"`) sigue saliendo con exit 1 hasta que se aplique la solicitud 1.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: EXIT 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». S1-16 OK.
- **Base de desarrollo** (consulta de solo lectura con `sqlite3 mode=ro`): una sola línea, Kueski, activa. El cambio no mueve ninguna cifra del servidor: 389500 hoy, 292125 al cierre de oct y 629344 al 10-abr-2027. No reinicié el servidor, no abrí sesión y no escribí en la base.
- **B-1, operaciones contra pruebas** (S6-03):

| Operación | Prueba |
|---|---|
| `getCreditTotals` | `TestOp_getCreditTotals`, `TestOp_getCreditTotals_DeudaVivaHoy`, `TestOp_getCreditTotals_LineaArchivada` |
| `getDebtTrend` | `TestOp_getDebtTrend`, `TestOp_getCreditTotals_DeudaVivaHoy`, `TestOp_getCreditTotals_LineaArchivada` |
| `getDashboardDebtCapacity` | `TestOp_getCreditTotals_LineaArchivada` (igualdad con K-3) y las de `api/tablero` |

## Supuestos
- **PROPUESTA aplicada como PROVISIONAL (M-1):** elegí la salida (ii) y no la (i). Es la única que cumple a la vez la decisión 1 del orquestador («la misma semántica de `getLiveDebt?date=hoy`») y la decisión 2 («de TODAS las líneas, como dice C1») sin cambiar textos FIRMES. Extendí la (ii) a las cifras del mes de `getCreditTotals` porque C1 ya dice que `payments_cents` del barómetro es «la misma cifra que `CreditTotals.month_household_cents`». Si el orquestador prefiere la (i), se revierte con una línea en `states` y los textos citados.
- Con `credit_id`, `getDebtTrend` cuenta las líneas nombradas aunque estén archivadas, para que la suma por línea sea igual al total sin filtro (punto 5, PROVISIONAL).

## Bloqueos
- Ninguno nuevo. El punto 5 (nuevo) de `S6-tablero-lineas-archivadas` pregunta a Axel en la puerta 4 si archivar una línea con pagos pendientes la saca de la deuda viva, de «Créditos del mes» y de G-5.

## Solicitudes fuera de alcance
1. `docs/criterios/C1.md`, líneas 6, 9, 106 y 107: cambiar «1.5.0» por «1.5.1» (en la 107, `version: "1.5.0"` por `version: "1.5.1"`) y nombrar la revisión «C1 v1.5.1 (escala de G-6, WF-3c bis)». Así la Verifica literal de C1-13 vuelve a encontrar la línea 7 de `openapi.yaml`.
2. Sigue en pie la solicitud 1 de `backend-cierre.md`: en `backend/internal/api/gastos-inicio/vistas.go:235` (`deriveCredits`), usar `creditos.LiveDebtAt(...)` en vez de `mc.LiveDebt(...)`. Este cambio no la afecta, porque `deriveCredits` ya arma cada línea con `ci.state(c)`.

## Pendientes
- Hay que reiniciar el servidor de desarrollo con el binario nuevo (`ops/dev-stop.sh`, `go build`, `ops/dev-run.sh`) antes de las capturas siguientes. Con la base de desarrollo actual las cifras son las mismas, así que no hay prisa. No lo reinicié para no cortar a otros agentes que lo estén usando.

## Propuestas para Axel
- Puerta 4: confirmar el punto 5 de `S6-tablero-lineas-archivadas`. Si archivas una línea que todavía tiene pagos, su deuda y sus pagos siguen sumando en Inicio y en Créditos hasta que termines de pagarla.

## Paquetes instalados
Ninguno.

## Consumo aproximado
29 acciones, unos 95 KB leídos y unos 22 KB escritos a mano (sin contar lo que regeneró `build.py`).
