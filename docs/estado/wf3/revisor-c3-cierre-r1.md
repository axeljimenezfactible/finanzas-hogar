# revisor-c3-cierre-r1 · Reporte

Agente: `revisor-c3-cierre-r1` (WF-3a cierre; independiente) · 2026-10-12 · Objeto: `docs/contratos/C3-integracion.md` v1.4 contra `docs/criterios/C3.md` 1.1.

## Hecho
- HECHO: evalué C3 v1.4 contra C3-01 a C3-11 con mis propias ejecuciones en una copia limpia. Veredicto **aprobado, ronda 1**.
- HECHO: agregué el bloque de firma (§9) al final de `C3-integracion.md`. SHA-256 del contenido firmado: `e341d182e08567d797ae503c3d3970b8d3542d650a632806559fc2df5507d163`. Lo recalculé después de firmar y coincide.
- HECHO: escribí el informe `docs/estado/wf3/REVISION-c3-cierre-r1.md`, con 3 hallazgos bajos (B-1 a B-3) y ninguno alto ni medio.

## Archivos
- `docs/contratos/C3-integracion.md`: solo se agregó el bloque `<!-- c3:firma:inicio -->` … `<!-- c3:firma:fin -->`. El contenido anterior quedó intacto.
- `docs/estado/wf3/REVISION-c3-cierre-r1.md` (nuevo).
- `docs/estado/wf3/revisor-c3-cierre-r1.md` (este reporte).

## Verificación
Copia en `datos/tmp/rev-c3r1/repo` (rsync del árbol de trabajo; `diff -rq` vacío). Binario y `dist` compilados en la copia. Base nueva: `migrate`, `semilla` y `semilla --desarrollo`, con 6 gastos, 4 compras y 72 precios.
- **C3-01**: `grep -c "^| P-"` = 41; T de U1 = 41 (27 + 14); encabezado con las 3 definiciones.
- **C3-02**: script propio → pantallas sin fila 0; filas con operaciones ≠ `x-necesidades`∩U1 0; operaciones sin fila 5, todas en excepciones (4 de K53 + `getOpenApiDocument`), 0 sin justificar; N de U1 sin operación ni Local 0; N de C1 sin fila 0.
- **C3-03**: script propio → 174 cadenas `/api/v1`, 0 fuera de `openapi.yaml`, 0 sin registrar; `openapi.yaml` y `routes.go` iguales (221). `npm run build` → salida 0 (PWA, 80 entradas). `c3-integracion.py` → 215 rutas, 0 y 0.
- **C3-04**: `go test ./... -count=1 -v` → 33 ok, 5 sin pruebas, 1151 PASS, 0 FAIL, 0 SKIP; 221/221 `TestContrato_Operaciones` y `TestContrato_Rutas` PASS. `npx playwright test --reporter=json` → 271 pasan, 4 omitidas, 0 fallan (366 s). `python3 -I frontend/scripts/c3-integracion.py --pw-json <mío> --go-log <mío>` → salida 0 y «Matriz del documento = matriz generada: sí». La corrí sobre el repositorio, sin `--escribir`, antes y después de firmar.
- **C3-05**: script propio → 164/164 imágenes, 0 faltantes; anchos 412 y 1440; evidencia = 4 rutas por fila.
- **C3-06**: `grep -c "| no |"` = 1 (Global, bloqueo `C3-capturas-global`).
- **C3-07**: bloque §9; `python3 -I` del bloque → `e341d182…d163` = bloque.
- **C3-08**: servidor aislado nuevo por archivo (127.0.0.1:8447) → `flujos.spec.ts` 7 passed, `pendientes.spec.ts` 4 passed, `f7-f8.spec.ts` 2 passed; F-1 = 3 toques; F-3 ruta A según U1 §4 (condición (c) → P-15; F-3A3 = 3 toques); 8 de 8 flujos. F-3A0 sale verde con su paso en FALLA (B-1).
- **C3-09**: `c3-09.spec.ts` en copia aislada → «120 cifras leídas en pantalla · diferencias != 0: 0», con deuda viva de hoy $3,895.00 = oráculo. API con sesión (script propio; la contraseña se leyó del archivo y no se imprimió) → `live_debt_cents` 389500; `getLiveDebt` 389500 y 629344; próximo pago 97375 al 2026-10-30 = `model.json`.
- **C3-10**: 21 filas leídas, es-MX, fechas de evidencia en ISO.
- **C3-11**: «Versión v1.4 · 2026-10-12»; `c3-cierre.md` con §Supuestos y §Bloqueos; 5 bloqueos `C3-*`.

## Supuestos
1. «Copia limpia» = árbol de trabajo actual (HEAD `0197dad` + cambios sin confirmar), con binario, `dist` y base reconstruidos desde cero en la copia.
2. Cambié solo en la copia el puerto de `preview` (4173 → 4183) y el del servidor aislado (8444 → 8447), para no chocar con otros agentes. No cambié ninguna lógica de prueba.
3. Corrí C3-09 contra un servidor aislado con la base nueva y no contra 8443. Los datos son los mismos (misma semilla) y la corrida es más independiente.
4. El contenido firmado son todos los bytes antes de la línea `<!-- c3:firma:inicio -->`. Así lo dice el bloque, junto con el comando para comprobarlo.
5. Califiqué B-1 como baja, igual que R2-5 y R3-5. F-3 se cumple con pruebas que sí comprueban su resultado.

## Bloqueos
- Ninguno nuevo. Siguen abiertos para la puerta 4: `C3-capturas-global`, `C3-matriz`, `C3-flujos` punto 2 (condición (c), PROVISIONAL) y `C3-ajustes-tablero` (opción A, PROVISIONAL).

## Pendientes
- B-1: `frontend/pruebas-u5/flujos/flujos.spec.ts:85`. Agregar `expect(reg.ok)` con lo que pide U3 tarjeta 5 junto con la condición (c), o declararla de observación. C3 §4 deja de contarla.
- B-2: `docs/design/U1-arquitectura.md` §4 F-3 («Resultado esperado con la semilla») y `docs/design/U3-dashboard.md` tarjeta 5 («Ejemplo trazable»): cambiar a $3,359.00 (+$264.50), que es lo que fija K46.
- B-3: línea 3 de C3. Quitar «sin firma» en la próxima versión, que exige una firma nueva.

## Propuestas para Axel
- Ratificar en la puerta 4 la condición (c) de F-3 ruta A y el bloque «Parámetros del tablero» de P-26.02. Son las dos cosas PROVISIONALES sobre las que descansa la firma.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- 55 acciones, unos 125 KB leídos y unos 16 KB escritos. Borré los temporales de `datos/tmp/rev-c3r1` al terminar.
