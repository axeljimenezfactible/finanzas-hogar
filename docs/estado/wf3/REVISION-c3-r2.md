# Revisión C3 · ronda 2 · `revisor-c3-r2`

Fecha: 2026-10-12 · Revisor: `revisor-c3-r2` (independiente; no redactó la matriz ni implementó nada) · Objeto: `docs/contratos/C3-integracion.md` **v1.1** (2026-10-12, «sin firma»; 47 182 bytes; SHA-256 del archivo revisado `c9e622b1599d33b87fe76f61d19db1dbf96c2938789020979066bac617232044`, igual antes y después de la revisión) sobre HEAD `6d691d0` con el árbol de trabajo de WF-3b · Criterios: `docs/criterios/C3.md` 1.1 (solo contra ellos).

## Veredicto

**DEVUELTO, sin firma.** 1 hallazgo alto, 0 medios y 6 bajos. C3-01 a C3-06 y C3-09 a C3-11 se cumplen con mis propias ejecuciones. **C3-08 no se cumple:** pasan 7 de 8 flujos de U1 contra el servidor real (F-5 con reserva). F-7 y F-8 ya tienen prueba y pasan, así que H-2 y H-3 de la ronda 1 quedan **cerrados**. F-3 sigue fallando con la semilla, así que **H-1 se reabre**. La propia matriz lo dice en §4: «C3-08 **no** se cumple hasta que Axel decida los puntos 2 y 3 de `C3-flujos`». No agrego el bloque de firma de C3-07 y no modifiqué la matriz.

## Hallazgos de la ronda 1

| Hallazgo r1 | Estado en r2 | Comprobación propia |
|---|---|---|
| H-1 (alta) F-3 falla con la semilla | **Reabierto** (sigue abierto; R2-1) | Corrí `flujos.spec.ts` en una copia nueva (`servidor-prueba.sh arrancar`, `127.0.0.1:8444`) y pasaron 6 de 7. Falló `:122`: en la ruta A, el toque 1 en la tarjeta «Mejor canasta» abre «Productos y precios» y no P-14. En la ruta B, la pantalla da Walmart · Acueducto de Guadalupe con «$264.50» de diferencia (= $3,359.00), y U1 pide $3,233.00 (+$138.50). Los puntos 2 y 3 de `C3-flujos.md` siguen sin decisión de Axel |
| H-2 (alta) F-7 sin prueba | **Cerrado** | `f7-f8.spec.ts:16` pasa en una copia nueva. Con la semilla, P-20 muestra «Faltan restricciones · Capturar» y no recomienda (K16, como pide U1 §4). Con restricciones de prueba llega a P-19 y guarda (5 toques con la semilla; 4 más hasta P-19) |
| H-3 (alta) F-8 sin prueba | **Cerrado** | `f7-f8.spec.ts:141` pasa en una copia nueva (28 s). El OCR local lee el ticket sintético y «Confirmar ticket» abre P-04. Toma 15 toques; C3-08 solo fija tope de toques para registrar gasto |
| H-4 (baja) deuda viva citada con `getDebtTrend` | Sigue abierto (R2-3) | `getLiveDebt` con `date=2027-04-10` devuelve `total_cents` **629344**, con las compras 125870 y 503474, iguales al detalle del oráculo (1,258.70 y 5,034.74). La matriz sigue citando `getDebtTrend` y el supuesto 6 sigue abierto |
| H-5 (baja) F-5 «con reserva» | Sigue abierto (R2-6) | `flujos.spec.ts:236` y `pendientes.spec.ts:13` pasan. El punto 4 de `C3-flujos` sigue sin decisión |
| H-6 (baja) lista de pruebas truncada | Sigue abierto (R2-4) | 23 filas terminan su lista con «…»; solo 60 nombres distintos son visibles |
| H-7 (baja) `listPeriods` sin llamada | Sigue abierto (R2-7) | Mi script lo confirma: es la única operación con pantalla que el frontend no llama. Sigue PROVISIONAL (bloqueo `C3-matriz` 3c) |
| H-8 (baja) «sin mocks» | **Cerrado** | `grep -n mocks` en la matriz no encuentra nada |

## Criterios evaluados

| Criterio | Resultado | Evidencia (ejecución propia del 2026-10-12) |
|---|---|---|
| C3-01 | Cumple | `grep -c "^| P-"` da **41** = T de U1 (27 `P-xx` de §2.1 y 14 `P-26.xx` de §2.3), más 1 fila Global. §1 define implementado, probado y capturado |
| C3-02 | Cumple | Script propio (`python3 -I`, sin usar el del redactor). 41 pantallas: 0 sin fila y 0 filas sin pantalla. 211 `N-xxx` definidas en U1 §3, las mismas 211 de §3.29. Ninguna queda sin operación ni Local (las Local son N-003, N-301 y N-317) y C1 no tiene ninguna fuera de U1. La columna de operaciones coincide con `x-necesidades` ∩ §3.29 en las 42 filas (0 diferencias). De 216 operaciones, 5 no tienen pantalla y las 5 están en la tabla de excepciones (4 K53 y `getOpenApiDocument`). `openapi.yaml` = `routes.go`: 216 = 216, mismo par (método, ruta) → `operationId` |
| C3-03 | Cumple | Script propio sobre 116 archivos de `frontend/src` (sin `mocks/`, historias ni `esquema.gen.ts`): 164 rutas distintas, 208 pares (método, ruta) y 166 literales `/api/v1` o `/health`. **0 fuera de `openapi.yaml` y 0 sin registrar.** `downloadAttachment` se llama por `armarUrl` en `<img src>`. Cierre de importaciones por fila: solo falta N-054 en P-01, que U1 declara «vía P-03». No se llama solo `listPeriods` (excepción PROVISIONAL). `npm run build` termina con código 0 (PWA con 82 entradas) y el árbol de git queda igual |
| C3-04 | Cumple | `go test ./... -count=1 -v`: 32 paquetes `ok`, 0 FAIL, 1126 PASS, 0 SKIP. Las 216 `TestContrato_Operaciones/<operationId>` pasan y cada operación de cada fila tiene la suya. `npx playwright test --reporter=json`: **245 ejecutadas: 241 pasan, 4 omitidas, 0 fallan y 0 inestables** (5.2 min). Los 60 nombres (archivo:línea) visibles en la matriz pasaron todos. Cada fila `P-` tiene una prueba que pasa con su `P-xx` en el título; P-22 la tiene por la regla de §1 (F-6 / «Registrar peso»). Con mis dos registros, el script del redactor da «Matriz del documento = matriz generada: sí» |
| C3-05 | Cumple | 41 filas con capturado = sí, así que 4N = 164. **Las 164 existen**, son PNG válidos de 412 px (Android) o 1440 px (escritorio) y la evidencia cita exactamente las 4 esperadas por fila. La carpeta tiene 164 PNG |
| C3-06 | Cumple | `grep -c "| no |"` da **1**: la fila Global (capturado), con su bloqueo `C3-capturas-global.md` citado en §3 |
| C3-07 | No se firma | Veredicto devuelto por C3-08; no agrego el bloque de firma |
| C3-08 | **No cumple** | En copias aisladas nuevas: `flujos.spec.ts` **6 de 7** (✘ `:122` F-3); `f7-f8.spec.ts` **2 de 2**; `pendientes.spec.ts` **4 de 4**. F-1 se registra en 3 toques (K10). Resultado: **7 de 8** flujos (F-5 con reserva), con los 4 obligatorios de U1-05 cubiertos (F-1, F-2, F-7, F-8). **F-3 falla** (R2-1). Base de desarrollo intacta: 6 gastos y 0 pesos, antes y después |
| C3-09 | Cumple | `c3-09.spec.ts` contra `finanzas.local:8443`: 1 de 1, **111 cifras con 0 diferencias** frente a `model.json`. `humo-inicio.spec.ts`: 1 de 1. API (solo lectura, sesión `axel`): próximo pago 97375 el 2026-10-30, partes 59882 y 37493, lectura B 690550; `getLiveDebt` al 2027-04-10 da 629344 = 6,293.44 del oráculo |
| C3-10 | Cumple | Leí 13 filas (Global, P-00 a P-10 y P-26.14) y §3 a §8: nombres y notas en es-MX, fechas de evidencia en ISO. El anglicismo de r1 ya no está |
| C3-11 | Cumple | `grep -n "v1"` encuentra la línea 3 («Versión v1.1 · 2026-10-12») y las líneas 162 y 163 del historial. `ls docs/estado/bloqueos/C3-*` muestra 3 bloqueos (`C3-capturas-global`, `C3-flujos`, `C3-matriz`). §Supuestos y §Bloqueos están en `c3-matriz.md` y en `correccion-c3-r1.md` |

## Hallazgos

| # | Sev. | Criterio | Hallazgo | Corrección sugerida |
|---|---|---|---|---|
| R2-1 (= H-1 r1, reabierto) | alta | C3-08 | F-3 «Comparar canasta» **falla** contra el servidor real con la semilla (`flujos.spec.ts:122`). Ruta A: el toque 1 en «Mejor canasta» abre «Productos y precios» porque faltan precios de los ingredientes del plan. Ruta B: la pantalla da Walmart +$264.50 ($3,359.00) y U1 da $3,233.00 (+$138.50). Hay además un síntoma que no se ve en la cuenta «6 de 7»: en `flujos.spec.ts:70` (F-3A0, semilla tal cual), la tarjeta dice «se compara la compra del 10-oct-2026», pero P-14 abre con la lista vacía de la semana 42 y su registro marca el paso como FALLA | Que Axel decida los puntos 2 y 3 de `C3-flujos.md`: precios de los ingredientes en la semilla o aclaración en U1, y cifra de U1 o del motor S5. Después, reejecutar en una copia aislada |
| R2-2 | baja | C3-08 | §4 (fila F-8) y `correccion-c3-r1.md` citan `f7-f8.spec.ts:139`, pero la prueba de F-8 está en la **línea 141** (así lo reporta Playwright en mi ejecución) | Citar `f7-f8.spec.ts:141` |
| R2-3 (= H-4) | baja | C3-09 | §5 cita la deuda viva con `getDebtTrend` y el supuesto 6 pide confirmarla | Citar `getLiveDebt` con `date=2027-04-10` (`total_cents` 629344, lo comprobé) y cerrar el supuesto 6 |
| R2-4 (= H-6) | baja | C3-04 | La columna «probado» trunca las pruebas Playwright con «…» en 23 filas | Listar todas o enlazar la lista completa que calcula el script |
| R2-5 | baja | C3-08 | `flujos.spec.ts:70` (F-3A0) pasa aunque su único paso queda en FALLA: `Flujo.paso` atrapa el error y la prueba no comprueba `reg.ok` al final. La cuenta «6 de 7» de §4 y §6 incluye esa prueba verde | Agregar `expect(reg.ok)` al final o declararla de observación y sacarla de la cuenta; anotar el síntoma en `C3-flujos` |
| R2-6 (= H-5) | baja | C3-08 | F-5 «pasa con reserva»: usa oct-2026 en lugar de «Mensual · sep-2026» y toma 6 toques para generar y descargar (U1: 5) | Que Axel decida el punto 4 de `C3-flujos` antes de la firma |
| R2-7 (= H-7) | baja | C3-02/C3-03 | `listPeriods` aparece en P-06 y P-08 aunque ninguna pantalla lo llama; es una excepción PROVISIONAL (bloqueo `C3-matriz` 3c) | Que Axel la confirme en la puerta 4 o que C1 la quite de `x-necesidades` de N-080, N-082 y N-100 |

## Nota para el integrador (no es hallazgo de la matriz)
- `frontend/pruebas-u5/flujos/` y `docs/design/capturas-u5/` siguen **sin seguimiento en git** (`??`). Si no se agregan al commit, la evidencia de §4, §5 y de la fila Global no queda en el repositorio.
- Mis corridas no dejaron rastro en la evidencia: respaldé `frontend/pruebas-u5/flujos/capturas/` y `resultados/` antes de correr y las restauré después. `git status` dio lo mismo antes y después, y la matriz no cambió (mismo SHA-256).
