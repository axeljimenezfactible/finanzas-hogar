# c3-tablero · WF-3c · matriz C3 v1.3 (sin firma)

Fecha: 2026-10-12 · Agente: c3-tablero (Sonnet). Alcance: `docs/contratos/C3-integracion.md`, `frontend/scripts/c3-integracion.py` (solo para el pendiente de pantalla y las subsecciones de P-01), `docs/estado/bloqueos/C3-ajustes-tablero.md` (nuevo), nota en `docs/estado/bloqueos/C3-flujos.md`, este reporte.

## Hecho
- **HECHO** · Pruebas de hoy: `go test ./... -count=1 -v` (33 paquetes ok, 1147 PASS, 0 FAIL, 221/221 `TestContrato_Operaciones`), `npm run build` rc 0 (80 entradas), `npx playwright test` completo (API simulada): 263 ejecutadas, 259 pasan, 4 omitidas, 0 fallan.
- **HECHO** · Matriz regenerada con `--escribir` (C1 v1.5.0, 221 operaciones, 5 nuevas): implementado sí 42/42, probado sí 42/42, capturado sí 41/42 (la fila Global sigue en «no» con su bloqueo). P-01: 21 necesidades con llamada y ruta registrada, 23 operaciones con su subprueba Go, 13 pruebas Playwright, 4 capturas más las 16 de las subsecciones (citadas aparte). Cada fila «sí» trae sus 4 rutas.
- **HECHO** · C3 pasó a v1.3 sin firma: encabezado y fuentes, §3 reescrita, §4 con la reejecución de los flujos, §5 con la actualización de C3-09, §6 con los comandos de hoy, §7 supuestos 7 a 9, §8 historial. F-3 lleva la aclaración de U1 (condición (c) de la ruta A, PROVISIONAL del orquestador).
- **HECHO** · Hallazgo: `updateDashboardSettings` (C1 v1.5, N-278) no la llama ninguna pantalla. No se admitió como excepción (el orquestador solo autorizó `listPeriods`): el script la muestra como «PANTALLA PENDIENTE» y sale con 1. Bloqueo nuevo `C3-ajustes-tablero` con 3 opciones.
- **HECHO** · Hallazgo: con el Inicio de U3 2.0, **no pasan** F-1 (`flujos.spec.ts:13`), F-3 (`:122`), F-4 (`:222` y `pendientes.spec.ts:57`), F-3A3 (`pendientes.spec.ts:117`) ni `c3-09.spec.ts:28`. Causa comprobada en copias de las pruebas: selectores `tarjeta-1`/`tarjeta-3` que ya no existen y una preparación de datos de F-3A3 que no contempla los precios de Walmart de la semilla. Con eso corregido en copia, F-1 (3 toques), F-4 (3) y F-3 ruta A (3) pasan, y 107 de 110 cifras de C3-09 coinciden con el oráculo. **Eso no es evidencia del repositorio: C3-08 y C3-09 quedan sin cumplir.**

## Archivos
- `/home/deepsec/finanzas-hogar/docs/contratos/C3-integracion.md`
- `/home/deepsec/finanzas-hogar/frontend/scripts/c3-integracion.py`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C3-ajustes-tablero.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C3-flujos.md` (actualización al pie)
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/c3-tablero.md`

## Verificación
- **C3-01** · `grep -c "^| P-" docs/contratos/C3-integracion.md` = 41 (27 P-xx + 14 P-26.xx, con P-26 contenedor) más la fila Global = 42 filas; T de U1 = 41.
- **C3-02** · script: «pantallas sin fila: 0 []; filas sin operación ni Local: 0 []»; «operaciones: 221; sin pantalla: 5 (excepciones: 5, de ellas K53: 4; sin justificar: 0)»; necesidades «de U1 sin operación ni Local: 0; de C1 sin fila en U1: 0». **Pero** «Operaciones con pantalla que el frontend no llama: 2 [listPeriods, updateDashboardSettings]; sin excepción justificada: 1 [updateDashboardSettings]» y rc 1.
- **C3-03** · «214 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (routes.go: 221; openapi: 221 operaciones)»; `npm run build` rc 0. «Matriz del documento = matriz generada: sí».
- **C3-04** · Go 221/221 `TestContrato_Operaciones` y 1147 PASS; Playwright (API simulada) 259 de 263, 4 omitidas, 0 fallan. Servidor real: ver C3-08.
- **C3-05** · «capturas: 164 esperadas (4 × 41 P-); filas con 4/4: 41; filas en «no»: [Global]». Además 16 de 16 de las subsecciones de P-01 en `docs/design/capturas-u5/evidencia/inicio-subsecciones/`.
- **C3-06** · `grep -c "| no |"` = 1 (Global, con `C3-capturas-global.md`). Pendiente abierto aparte: `updateDashboardSettings` (bloqueo `C3-ajustes-tablero`).
- **C3-08** · NO se cumple: con las pruebas del repositorio pasan 5 de 8 flujos (F-2, F-5, F-6, F-7 `f7-f8.spec.ts:16`, F-8 `:141`); fallan F-1, F-3 y F-4. Entre los 4 obligatorios de U1-05 falla F-1.
- **C3-09** · NO se cumple con evidencia del repositorio: `c3-09.spec.ts` falla en `tarjeta-3`. Copia adaptada: 110 cifras leídas, 107 coinciden, 3 sin poder leerse (tarjeta nueva). Próximo pago 973.75 y partes 598.82 / 374.93 se vieron en pantalla en F-4; ratio del barómetro según `u5-tablero.md`.
- **C3-10** · nombres y notas en es-MX, fechas ISO. **C3-11** · versión v1.3 y fecha; supuestos 7 a 9; bloqueos `C3-ajustes-tablero` y `C3-flujos`.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` rc 0 («Verificación ampliada OK»).

## Supuestos
- `updateDashboardSettings` no es excepción técnica (no hay autorización); lo decide Axel u orquestador.
- «Probado» conserva la regla PROVISIONAL del punto 4 de `C3-matriz`.
- Los diagnósticos con copias no alteran ninguna prueba ni la matriz.
- No reejecuté `humo-inicio.spec.ts` ni `tablero-cifras`.

## Bloqueos
- `docs/estado/bloqueos/C3-ajustes-tablero.md` (nuevo, PREGUNTA).
- `docs/estado/bloqueos/C3-flujos.md` (actualización: pruebas viejas con U3 2.0).

## Pendientes
- Pruebas de `frontend/pruebas-u5/flujos/` a corregir (dueño `u5-flujos`) y reejecutar; luego regenerar C3 y firma por el revisor independiente.
- Apartado «Parámetros del tablero» en P-26.02 (o decisión).
- Diferencias de cifras del tablero reportadas por `u5-tablero` (K-3 $18,999.22 contra $3,895.00; «Deuda viva proyectada» desde $18,025.47) y el desborde a 475 px de P-01 Detalle en Android por una insignia `w-max`: no son de C3; quedan para el agente del tablero.

## Propuestas para Axel
- Elegir la opción A de `C3-ajustes-tablero` (apartado en Ajustes) en lugar de una excepción.
- Ratificar la condición (c) de F-3 ruta A en la puerta 4.
- P-14: no mostrar $0.00 ni diferencia negativa para una tienda sin ningún precio.

## Solicitudes fuera de alcance
- `frontend/pruebas-u5/flujos/flujos.spec.ts` líneas 18, 19, 41 y 42: cambiar `tarjeta-1` por el selector de «Gasto del mes» del tablero (la tarjeta cuyo texto cambia al guardar el gasto); línea 226: `tarjeta-3` por `tarjeta-K-4` (es un `section`; tocar su enlace interior). Líneas 74-146 (F-3A0/F-3 ruta A): no esperar P-14 con 3 artículos con precio en una lista de decenas; aplicar la condición (c) (tarjeta lleva a «Productos y precios»).
- `frontend/pruebas-u5/flujos/pendientes.spec.ts` línea 62: `tarjeta-3` por `tarjeta-K-4`; línea 135: `l.prices.length === 0` por `!l.prices.some((p: any) => p.store_id === 2)`.
- `frontend/pruebas-u5/flujos/c3-09.spec.ts` líneas 31-39: leer «Próximo pago», Hogar y Karen de `tarjeta-K-4`; quitar de ahí «Disponible hoy, lectura B» (ya se lee en Detalle).
- `frontend/src/pantallas/ajustes-a/S02Presupuesto.tsx` (o una sección nueva de P-26.02): apartado «Parámetros del tablero» que llame PATCH `/api/v1/settings/dashboard` (opción A del bloqueo).
- P-01 Detalle en Android desborda a 475 px por una insignia `w-max` (de `u5-tablero`).

## Paquetes instalados
Ninguno.

## Consumo aproximado
Acciones: ~40 · KB leídos: ~190 · KB escritos: ~30. Temporales de `datos/tmp/c3/` y `datos/tmp/u5-flujos/` borrados al terminar.
