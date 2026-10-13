# revisor-backend-cierre-r1 · WF-3c bis · 12-oct-2026

Revisor independiente (Opus). Evalúa `backend-cierre` solo contra S6-03, S6-04, C1.md, S2.md, S3.md y las decisiones 1, 2 y 3 de WF-3c bis. No corrige, no hace commit y no envía Telegram. El informe completo está en `docs/estado/wf3/REVISION-backend-cierre-r1.md`.

## Hecho
- **Dictamen: NO APROBADO.** Quedan 0 altos, 3 medios (M-1, M-2 y M-3) y 4 bajos.
- Reejecuté las verificaciones de Make, la cobertura, `build.py --check` y una sonda de solo lectura contra el servidor de desarrollo. Corrí también una prueba propia con overlay: una línea inactiva con pagos programados.
- Las cifras del motor y del oráculo coinciden en la base sembrada. Los 12 cierres de `getDebtTrend` son iguales a `getLiveDebt` al último día de cada mes.

## Archivos
- `docs/estado/wf3/REVISION-backend-cierre-r1.md` (informe) y este reporte. No toqué nada más del repositorio.
- Temporales en `datos/tmp/rev-bc/` (`zz_revisor_test.go` con `overlay.json`, la sonda, `valida.py` y logs). Los borré al terminar.

## Verificación
- **(a)** `make verify-wf3-go` EXIT 0 (33 OK y 5 sin pruebas). `make verify-wf3-contrato` EXIT 0. `build.py --check` EXIT 0 (22 archivos idénticos). Cobertura con `go test -cover`: creditos 80.9 %, presupuesto 77.7 %, tablero 87.8 % e internal/api 88.8 %.
- **S6-03:** 221 operaciones en 176 rutas, 0 diferencias; «221 de 221 operaciones cumplen». CUMPLE, con B-1: falta la tabla operación → prueba en el reporte.
- **S6-04:** centavos enteros, `AAAA-MM`, `ErrorBody` y mensajes en español de México. 400 de B-6 comprobado en vivo. CUMPLE.
- **(b)** Servidor: `live_debt_cents` 389500 = `getLiveDebt` hoy 389500. G-5 empieza en 292125 al cierre de octubre (2,921.25 del oráculo). 629344 al 10-abr-2027. M-19: 1 pago y 59882 por pagar. K-2 de noviembre coincide con `getIncomeMonth(2026-10)`. `scale_max_pct` 60, requerido, y las respuestas validan con jsonschema.
- **(b) Prueba propia, línea inactiva.** El barómetro y el resumen mensual la cuentan (750000, 3 pagos, 38 %), y a1 funciona. Pero `getCreditTotals` y `getDebtTrend` no la cuentan: 1125000 contra 1325000 de `getLiveDebt` hoy, y 375000 contra 475000 al 31-oct (M-1).
- **(c)** `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` EXIT 0, «Verificación ampliada OK».
- **(d)** RN-21 tiene un párrafo PROVISIONAL que cita `S6-deuda-viva-hoy` (abierto, en `INDICE.md`), y RN-30 cita `S6-tablero-lineas-archivadas`. OK.
- **C1-13:** `grep -n -m1 'version: "1.5.0"' openapi.yaml` sale con exit 1 (M-3). **C1-15:** `S6-deuda-viva-hoy` tiene 0 citas en C1 (M-2).

## Supuestos
- PREGUNTA: tomé «la deuda viva proyectada arranca ahí» como «parte de 389500 y baja con el pago del 30-oct», que es la PROPUESTA de `backend-cierre` (B-4).
- La prueba de la línea inactiva usa datos falsos insertados en una base de prueba. Archivar una línea con deuda viva sí se puede por la API (`condiciones.go:142`, sin guarda), así que el caso es alcanzable.

## Bloqueos
- Ninguno nuevo. M-1 propone ampliar `S6-tablero-lineas-archivadas` o `S6-deuda-viva-hoy`; lo decide el orquestador.

## Pendientes
- Corrección de M-1, M-2 y M-3 y una ronda 2 de esta revisión.

## Propuestas para Axel
- Puerta 4: decidir si la deuda viva de K-3, P-17 y G-5 incluye las líneas archivadas que todavía tienen pagos pendientes. Así la cifra sería la misma que la de `getLiveDebt` y P-07.

## Paquetes instalados
Ninguno. `uv` usó pyyaml y jsonschema en un entorno efímero, como en `verify-all`.

## Consumo aproximado
Unas 50 acciones, unos 190 KB leídos y unos 25 KB escritos.
