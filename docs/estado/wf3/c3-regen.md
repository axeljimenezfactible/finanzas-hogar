# c3-regen · WF-3a bis · matriz C3 v0.95

Fecha: 2026-10-12 · Agente: c3-regen (Sonnet). Alcance: `docs/contratos/C3-integracion.md`, `docs/estado/bloqueos/C3-matriz.md`, `frontend/scripts/c3-integracion.py`, este reporte.

## Hecho
- **HECHO** · Compilé de nuevo la interfaz (`npm run build`), corrí Go (`TestContrato_Operaciones|TestContrato_Rutas`) y Playwright completo (JSON) y regeneré la matriz con `--escribir`.
- **HECHO** · Resultado: filas con implementado = no: **0** (antes 7); probado sí 42/42; capturado sí 0/42 (WF-3b). Operaciones sin llamada: **1** (`listPeriods`, antes 7). Operaciones sin pantalla: 5, todas justificadas (4 de K53 y `getOpenApiDocument`).
- **HECHO** · Decisiones del orquestador anotadas punto por punto en `docs/estado/bloqueos/C3-matriz.md` (sección final). Fila P-26 «Local (contenedor)»; `queryView` y `listIncomes` ya se llaman, por lo que no necesitan excepción; `listPeriods` queda como excepción técnica PROVISIONAL con su cobertura (N-080 `getBudgetMonth`, N-082 `getBudgetFortnights`, N-100 `getIncomeMonth`/`listIncomes`).
- **HECHO** · `c3-integracion.py`: nueva tabla `EXCEPCIONES_NO_LLAMADA` y comprobación de cobertura. Una operación no llamada sin excepción escrita, o con una necesidad que ninguna otra operación llamada cubra, hace fallar el script (rc 1). La tabla al pie de la matriz ahora trae «necesidades cubiertas por otra operación» y «motivo». Texto de la fila P-26 ajustado.
- **HECHO** · `C3-integracion.md`: versión 0.95 (sigue sin firma), §1 con la decisión de P-26, §3 reescrita (solo Capturado = no, más las decisiones PROVISIONALES).

## Archivos
- `/home/deepsec/finanzas-hogar/docs/contratos/C3-integracion.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C3-matriz.md`
- `/home/deepsec/finanzas-hogar/frontend/scripts/c3-integracion.py`
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/c3-regen.md`

## Verificación
- `npm run build` rc 0 (precache 81 entradas).
- `go test ./internal/api -run 'TestContrato_Operaciones|TestContrato_Rutas' -count=1 -v` → ok; 216 `--- PASS: TestContrato_Operaciones/<operationId>`.
- `npx playwright test --reporter=json` (2 corridas completas): 235 ejecutadas, **230 pasan, 4 omitidas, 1 falla** en ambas (ver Bloqueos). Aislado, `gastos-presupuesto/correccion-bis.escritorio.spec.ts` pasa 11/11 (2 veces).
- `npm run test:integracion` (servidor real finanzas.local:8443): 1 passed.
- **C3-01** · `grep -c "^| P-" docs/contratos/C3-integracion.md` = 41 (27 P-xx + 14 P-26.xx, con P-26 contenedor) más la fila Global = 42 filas.
- **C3-02** · script: «pantallas U1: 41 (P-xx 27 + P-26.xx 14); filas: 42», «pantallas sin fila: 0 []; filas sin operación ni Local: 0 []», «operaciones: 216; sin pantalla: 5 (excepciones: 5, de ellas K53: 4; sin justificar: 0 [])». Con `--escribir` y sin él: «Matriz del documento = matriz generada: sí», rc 0.
- **C3-03** · «210 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (routes.go: 216; openapi: 216 operaciones)».
- **C3-04** · Go 216/216 PASS; Playwright 230 de 235 (4 omitidas, 1 con falla; ninguna fila de la matriz depende solo de ella).
- **C3-06** · `grep -c "| no |"` = 42, todas por «Capturado» (WF-3b); implementado y probado sin «no».
- **C3-10** · nombres y notas en es-MX, fechas ISO.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` rc 0 («Verificación ampliada OK»).

## Supuestos
- Los agentes de área ya habían terminado (las fuentes no cambiaron desde el build de las 07:36 salvo archivos generados).
- La decisión 3(c) se aplicó solo a `listPeriods`; `queryView` y `listIncomes` ya tienen llamada.

## Bloqueos
- `docs/estado/bloqueos/C3-matriz.md` (todas las decisiones siguen PROVISIONALES hasta la puerta 4).
- Prueba fallida repetible en la corrida completa (ver Pendientes).

## Pendientes
- Corregir `frontend/pruebas/gastos-presupuesto/correccion-bis.escritorio.spec.ts:79` («P-02 Ctrl+/ …»: espera 2 parches y recibe 0 en la corrida completa; pasa aislada): dependencia de orden o de estado. Antes de firmar C3 (C3-04).
- WF-3b: capturas U5 (columna «capturado»), C3-08 con F-1 a F-6 contra el servidor real, C3-09 completo, firma C3-07 por revisor independiente.

## Propuestas para Axel
- Aceptar `listPeriods` como excepción técnica (o quitarla de C1 `x-necesidades`) y la fila P-26 «Local (contenedor)» en la puerta 4.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Acciones: ~25 · KB leídos: ~150 · KB escritos: ~20. Temporales de `datos/tmp/c3/` borrados al terminar.
