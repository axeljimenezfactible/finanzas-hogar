# Revisión C3 · ronda 1 · `revisor-c3-r1`

Fecha: 2026-10-12 · Revisor: `revisor-c3-r1` (independiente; no redactó la matriz ni implementó nada) · Objeto: `docs/contratos/C3-integracion.md` **v1.0** (2026-10-12, «sin firma»; 44 573 bytes, SHA-256 del archivo revisado `09e6b63c3fd9932c0bc6deef82563522e35ff4adb534f6859371c52195cad502`) sobre HEAD `6d691d0` con el árbol de trabajo de WF-3b · Criterios: `docs/criterios/C3.md` 1.1 (solo contra ellos).

## Veredicto

**DEVUELTO, sin firma.** 3 hallazgos altos, 0 medios, 5 bajos. C3-01 a C3-06 y C3-09 a C3-11 se cumplen con mis propias ejecuciones; **C3-08 no se cumple** (5 de 8 flujos de U1 pasan contra el servidor real: F-3 falla y F-7 y F-8, que cubren los obligatorios (b) y (d) de U1-05, no tienen prueba). Por eso no se agrega el bloque de firma de C3-07 y la matriz no se modificó.

## Criterios evaluados

| Criterio | Resultado | Evidencia (ejecución propia del 2026-10-12) |
|---|---|---|
| C3-01 | Cumple | `grep -c "^| P-"` = **41** = T de U1 (27 `P-xx` de §2.1 + 14 `P-26.xx` de §2.3) más 1 fila Global; §1 define implementado, probado y capturado |
| C3-02 | Cumple | Script propio (`python3 -I`): 41 pantallas, 0 sin fila, 0 filas sin pantalla; 211 `N-xxx` definidas en U1 §3 = 211 en §3.29 expandida; 0 sin operación ni Local (Local: N-003, N-301, N-317); 0 de C1 fuera de U1; la columna de operaciones de las 42 filas = `x-necesidades` ∩ U1 §3.29 (0 diferencias); 216 operaciones: 5 sin pantalla, las 5 en la tabla de excepciones (4 K53 + `getOpenApiDocument`); `openapi.yaml` = `routes.go` (216 = 216, mismo `operationId` por ruta) |
| C3-03 | Cumple | Script propio: el frontend (`frontend/src` sin `mocks/`, historias ni `esquema.gen.ts`) invoca 166 rutas distintas / 210 pares (método, ruta); **0 fuera de `openapi.yaml`, 0 sin registrar**; cierre de importaciones desde el archivo de entrada de cada fila: 0 necesidades sin llamada (P-01 N-054 vía P-03). No llamadas: las 5 excepciones y `listPeriods` (excepción PROVISIONAL, bloqueo `C3-matriz` 3c). `npm run build`: exit 0, 0 errores, PWA 82 entradas; árbol de git sin cambios tras el build |
| C3-04 | Cumple | `go test ./... -count=1 -v`: 32 paquetes `ok`, 0 FAIL, 1126 PASS, 0 SKIP; 216/216 `TestContrato_Operaciones/<operationId>` PASS y cada operación de cada fila tiene la suya. `npx playwright test` (JSON): **245 ejecutadas, 241 pasan, 4 omitidas, 0 fallan, 0 inestables** (5.2 min); las 69 pruebas que la matriz nombra (archivo:línea) pasaron todas. El script del redactor con mis dos registros: «Matriz del documento = matriz generada: sí» |
| C3-05 | Cumple | 41 filas «capturado = sí» → 4N = 164; **164 existen**, PNG válidos, ancho 412 (Android) / 1440 (escritorio), nombres citados = esperados; 164 PNG en la carpeta; `docs/estado/U5-revision.md` ronda 3: 164 de 164 |
| C3-06 | Cumple | `grep -c "| no |"` = **1** = filas con bloqueo citado (Global, capturado; `docs/estado/bloqueos/C3-capturas-global.md` existe y se menciona en §3 y en el reporte del redactor) |
| C3-07 | No se firma | Veredicto devuelto por C3-08; no se agrega bloque de firma |
| C3-08 | **No cumple** | Copias aisladas en `127.0.0.1:8444` (`servidor-prueba.sh`): `flujos.spec.ts` **6 de 7** (✘ `:122` F-3, paso 3 «Toque 1: tarjeta Mejor canasta», P-01 → P-14 no visible); `pendientes.spec.ts` 4 de 4. F-1 = 3 toques (K10). F-7 y F-8: ninguna prueba en `frontend/pruebas-u5/`, `pruebas/` ni `pruebas-integracion/`. Resultado: 5 de 8 (F-5 con reserva). Base de desarrollo intacta (expense 6, weight_entry 0, antes y después) |
| C3-09 | Cumple | `c3-09.spec.ts` contra `finanzas.local:8443`: 1/1, **111 cifras, 0 diferencias** con `model.json`; `humo-inicio.spec.ts` 1/1. API (solo lectura, sesión `axel`): próximo pago 97375 el 2026-10-30, partes 59882 / 37493, lectura B 690550; `getLiveDebt` con `date=2027-04-10` → **629344**. Las 18 cifras P1–P6 y las del crédito de §5 = `model.json` |
| C3-10 | Cumple | 10 filas leídas (Global, P-00 a P-08, P-26.05): nombres y notas en es-MX; fechas de evidencia ISO; un anglicismo («sin mocks», §4; hallazgo bajo) |
| C3-11 | Cumple | `grep -n "v1"` → línea 3 «Versión v1.0 · 2026-10-12» y línea 161 (historial); `ls docs/estado/bloqueos/C3-*` → 3 (`C3-capturas-global`, `C3-flujos`, `C3-matriz`); `docs/estado/wf3/c3-matriz.md` con §Supuestos y §Bloqueos |

## Hallazgos

| # | Sev. | Criterio | Hallazgo | Corrección sugerida |
|---|---|---|---|---|
| H-1 | alta | C3-08 | F-3 «Comparar canasta» **falla** contra el servidor real con la semilla (reproducido: `flujos.spec.ts:122`, ruta A paso 3; ruta B con cifra de Walmart distinta de U1, H-4) | Que Axel decida `C3-flujos` puntos 2 y 3 (semilla con precios de ingredientes o aclaración de U1; cifra de U1 o del motor S5) y reejecutar en copia aislada |
| H-2 | alta | C3-08 | F-7 «Compra a crédito con el comparador» (obligatorio (b) de U1-05) **sin prueba** contra el servidor real | Prueba en `frontend/pruebas-u5/flujos/`: P-01 → G-01 → P-17 → P-20 → P-19, copia aislada; con la semilla, P-20 debe mostrar «Faltan restricciones · Capturar» si faltan las de P-26.05 (K16) |
| H-3 | alta | C3-08 | F-8 «Ticket con OCR» (obligatorio (d) de U1-05) **sin prueba** contra el servidor real | Prueba P-01 → P-03 → P-05 → P-04 con una foto de ticket de prueba y el OCR local, copia aislada |
| H-4 | baja | C3-09 | La cifra de API de la deuda viva cita `getDebtTrend`, donde 629344 sale en el mes «2027-03» (el mes «2027-04» trae 377606); el supuesto 6 pide confirmarlo | Citar `getLiveDebt` con `date=2027-04-10` (devuelve `total_cents` 629344, comprobado) y cerrar el supuesto 6 |
| H-5 | baja | C3-08 | F-5 «pasa con reserva»: oct-2026 en lugar de «Mensual · sep-2026» y 6 toques para generar y descargar (U1: 5) | Decidir `C3-flujos` punto 4 antes de la firma |
| H-6 | baja | C3-04 | La columna «probado» trunca las pruebas Playwright con «…» (solo 69 nombres visibles); el cruce de nombres solo puede hacerse con los visibles | Listar todas o enlazar la lista completa que calcula el script |
| H-7 | baja | C3-02/C3-03 | `listPeriods` figura en P-06 y P-08 sin que ninguna pantalla lo llame; se admite como excepción PROVISIONAL del orquestador (bloqueo `C3-matriz` 3c) | Que Axel confirme en la puerta 4 o que C1 la saque de `x-necesidades` de N-080, N-082 y N-100 |
| H-8 | baja | C3-10 | «sin mocks» (§4) es anglicismo; el resto del documento dice «API simulada» | «sin API simulada» |

## Nota para el integrador (no es hallazgo de la matriz)
- `frontend/pruebas-u5/flujos/` está **sin seguimiento en git** (`??`) y la matriz cita de ahí las pruebas de §4 y §5 y las capturas `U5-13-sin-conexion-*.jpg` de la fila Global: si no se agrega al commit, esa evidencia no queda en el repositorio.
- Mis corridas de `flujos.spec.ts` y `pendientes.spec.ts` reescribieron las capturas de paso `.jpg` de esa carpeta (las genera `apoyo.ts` en cada corrida; mismas pantallas y misma fecha). No corrí `captura-sin-conexion.spec.ts`: las `U5-13-*` quedaron intactas.
