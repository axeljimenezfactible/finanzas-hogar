# REVISIÓN · C3 v1.4 · cierre, ronda 1

Fecha: 2026-10-12 · Revisor: `revisor-c3-cierre-r1` (independiente; no redacté la matriz ni implementé nada) · Objeto: `docs/contratos/C3-integracion.md` **v1.4** (2026-10-12, `c3-cierre`), SHA-256 del contenido revisado `e341d182e08567d797ae503c3d3970b8d3542d650a632806559fc2df5507d163` (el mismo al empezar y al firmar) · Base: HEAD `0197dad` más el árbol de trabajo actual · Criterios: `docs/criterios/C3.md` 1.1, y nada más.

## Veredicto

**APROBADO y FIRMADO (ronda 1).** Con mis propias ejecuciones en una copia limpia se cumplen C3-01 a C3-11. Hay 0 hallazgos altos, 0 medios y 3 bajos. Ninguno de los bajos deja abierto un criterio. El bloque de firma está al final de `C3-integracion.md` (§9). El hash del contenido firmado se recalculó después de agregar el bloque y coincide.

## Cómo reejecuté (copia limpia)

- **Copia:** `rsync` del árbol de trabajo a `datos/tmp/rev-c3r1/repo`, sin `.git`, `node_modules` (enlace), `dist`, `backend/bin`, `datos/tmp` ni la base de desarrollo. El `diff -rq` contra el original salió vacío.
- **Binario y base:** compilé el binario en la copia. La base es nueva: `finanzas migrate`, `finanzas semilla` y `finanzas semilla --desarrollo`. Quedó con 6 gastos, 4 compras a crédito y 72 precios, el mismo estado que cita C3 §7.12.
- **Frontend:** `npm run build` en la copia.
- **Puertos:** solo en la copia, el `preview` pasó a 4183 y el servidor aislado a 8447, para no chocar con otros agentes. No cambié ninguna lógica de prueba.
- **Flujos y C3-09:** cada archivo de flujos corrió contra un servidor aislado nuevo (copia de la base nueva). `c3-09.spec.ts` también corrió en una copia aislada, no en el servidor de desarrollo.

## Criterios

| Criterio | Resultado | Evidencia propia |
|---|---|---|
| C3-01 (alta) | Cumple | `grep -c "^\| P-"` = 41 = T de U1 (27 P-xx + 14 P-26.xx, mi script); fila Global aparte; §1 define implementado, probado y capturado |
| C3-02 (alta) | Cumple | Script propio (`python3 -I`, U1 §2.1/§2.3/§3.29 + `openapi.yaml` con YAML + matriz). 0 pantallas sin fila. 0 filas cuyas operaciones difieran de `x-necesidades` ∩ U1. 221 operaciones: 5 sin fila, las 5 en la tabla de excepciones (4 de K53 y `getOpenApiDocument`). 0 sin excepción. 0 `N-xxx` de U1 sin operación ni Local. 0 `N-xxx` de C1 sin fila. El script del repositorio da lo mismo |
| C3-03 (alta) | Cumple | Script propio: 174 cadenas `/api/v1/…` distintas en `frontend/src`, 0 fuera de `openapi.yaml`, 0 sin ruta en `routes.go`; `openapi.yaml` = `routes.go` (221 = 221, 0 diferencias). Script del repositorio: 215 rutas invocadas, 0 y 0. `npm run build`: salida 0 (PWA con 80 entradas) |
| C3-04 (alta) | Cumple | `go test ./... -count=1 -v`: 33 `ok`, 5 sin pruebas, 0 FAIL, 1151 `--- PASS`, 0 FAIL, 0 SKIP; 221 `TestContrato_Operaciones/*` PASS y `TestContrato_Rutas` PASS. `npx playwright test --reporter=json`: 271 pasan, 4 omitidas (pruebas de escritorio con Ctrl+/ en el proyecto móvil), 0 fallan, 0 inestables. Con estos dos registros, `c3-integracion.py` (sin `--escribir`) sale con 0 y dice «Matriz del documento = matriz generada: sí» |
| C3-05 (alta) | Cumple | Script propio: 41 filas «capturado = sí» × 4 = 164 imágenes, 0 faltantes, todas PNG con ancho 412 (Android) o 1440 (escritorio), evidencia = las 4 rutas esperadas; 164 PNG en la carpeta; `U5-revision.md` ronda 2: 164 de 164 |
| C3-06 (alta) | Cumple | `grep -c "\| no \|"` = 1. Es la fila Global (capturado), con `docs/estado/bloqueos/C3-capturas-global.md`, citado en C3 §3 («decide Axel u orquestador») y en `c3-cierre.md` §Bloqueos |
| C3-07 (alta) | Cumple (firmado) | El bloque §9 trae nombre, fecha ISO, ronda 1, veredicto y SHA-256. Recalculado con el comando del bloque: `e341d182…d163` = bloque |
| C3-08 (alta) | Cumple | Contra el servidor real, con una copia aislada nueva por archivo: `flujos.spec.ts` 7 pasan, `pendientes.spec.ts` 4 pasan y `f7-f8.spec.ts` 2 pasan. F-1: 3 toques (K10), y el gasto aparece en K-1 y en P-02. F-2: 4. F-3 ruta B: 5 toques. F-3 ruta A según U1 §4: con la lista de F-2 se cumple la condición (c), la tarjeta queda en vacío y lleva a P-15. En F-3A3 hace 3 toques (tarjeta → Detalle → «Usar Bodega Aurrera · Vallejo»). F-4: 3 (y 4 al marcar pago). F-5: 3 (y 6 completo). F-6: 3. F-7 A/B: OK. F-8: OK. Son 8 de 8 flujos, con los 4 obligatorios. Ver B-1 sobre F-3A0 |
| C3-09 (alta) | Cumple | `c3-09.spec.ts` en copia aislada: **120 cifras, 0 diferencias**. Incluye K-4 próximo pago $973.75 (hogar $598.82, Karen $374.93), K-3 deuda viva de hoy **$3,895.00**, disponible A $10,000.00 y B $6,905.50, B-1 2 % y $9,201.18, P1–P6 y deuda viva al 10-abr-2027 $6,293.44. API (script propio con sesión): `getCreditTotals` `live_debt_cents` 389500; `getLiveDebt` 2026-10-12 = 389500 y 2027-04-10 = 629344; `getDashboardNextPayment` 97375 al 2026-10-30. Las tres coinciden con lo que calculé de `oraculo/model.json` |
| C3-10 (media) | Cumple | Leí 21 filas (Global, P-00 a P-02, P-11 a P-22 y P-26.10 a P-26.14): nombres y notas en es-MX. Las fechas de evidencia están en ISO (2026-10-12). Las fechas tipo «10-oct-2026» son textos que se ven en pantalla. Ver B-3 |
| C3-11 (media) | Cumple | Línea 3: «Versión v1.4 · 2026-10-12». `c3-cierre.md` tiene §Supuestos y §Bloqueos. Existen `C3-ajustes-tablero`, `C3-capturas-global`, `C3-flujos`, `C3-matriz` y `C3-ronda3-escalamiento` |

## Hallazgos (todos bajos)

| Id | Sev. | Criterio | Hallazgo | Corrección sugerida |
|---|---|---|---|---|
| B-1 (= R2-5, R3-5) | baja | C3-08 | `flujos.spec.ts:71` (F-3A0) sale en verde aunque su único paso queda en FALLA, porque falta `expect(reg.ok)` (línea 85). En mi corrida el registro dice «F-3A0 … toques 1 · FALLA»: P-14 abre la semana 42 vacía y no aparece «más barata». Además, la tarjeta se leyó como «Mejor canasta — $3,094.50 … se compara la compra del —», porque la tienda y la fecha aún no habían cargado. C3 §4 la cuenta entre las que pasan («7 de 7», fila F-3). F-3 se cumple igual por `:123` y `pendientes.spec.ts:119`, que sí comprueban su resultado | Agregar `expect(reg.ok)` con lo que pide U3 tarjeta 5 junto con la condición (c): P-14 abre la semana y dice «La canasta no tiene artículos». Esperar a que la tarjeta termine de cargar antes de leerla. O declararla de observación y sacarla de la cuenta de §4 |
| B-2 | baja | C3-08 (fuera de C3) | U1 §4 F-3 («Resultado esperado con la semilla») y U3 tarjeta 5 («Ejemplo trazable») siguen diciendo Walmart $3,233.00 (+$138.50). Lo que fija K46 (Q2-47 = A), y lo que dan la pantalla y la API, es Acueducto de Guadalupe $3,359.00 (+$264.50). Ya se pidió en C3 v1.2 (`C3-flujos` punto 3) | Que el dueño de U1 y U3 actualice esas dos líneas a la cifra de K46 |
| B-3 | baja | C3-10/C3-11 | La línea 3 de C3 todavía dice «sin firma» y «esta versión no lo trae». Ahora está firmada (§9). No lo cambié porque forma parte del contenido firmado | En la próxima versión, actualizar la línea 3 y volver a firmar |

## Lo que no hice

- No cambié ninguna prueba, ningún código ni el contenido de la matriz. Solo agregué el bloque de firma al final.
- No toqué la base de desarrollo ni el servidor de 8443.
- Las capturas que generaron los flujos quedaron en la copia y se borraron con ella.
