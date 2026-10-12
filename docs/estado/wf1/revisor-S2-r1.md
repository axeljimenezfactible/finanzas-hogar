# revisor-S2-r1 · Reporte final (WF-1, 11-oct-2026)

## Hecho
- Revisé de forma independiente el motor financiero S2 (`backend/internal/motor_fin`, `backend/cmd/motores`) y su oráculo (`oraculo/fin_*.py`, `oraculo/compara.py`, `oraculo/fixtures/s2/`). Usé solo `docs/criterios/S2.md` 1.2, decisiones §B y §K y `S2-interfaz.md` 1.2.
- Reejecuté las verificaciones (a) a (g) del encargo con comandos propios. Escribí 3 scripts propios en `datos/tmp/rev-s2/`:
  - una comparación de las 159 hojas contra `model.json`;
  - un generador de fixtures propias con esperado calculado a mano;
  - una revisión de la propiedad Σ partes = total en toda la salida del motor.
- Probé el motor y el oráculo con 4 sondas de borde fuera de las fixtures: una ventana posterior al primer evento, una ocurrencia quincenal con ancla fuera del día de pago, `accrual_months` 0 e `interval_count` 0.
- Informe: `docs/estado/wf1/REVISION-S2-r1.md`. Veredicto: **NO APROBADO** (0 altos, 3 medios, 2 bajos).

## Archivos
- Creados: `docs/estado/wf1/REVISION-S2-r1.md` y este reporte.
- No modifiqué nada del alcance revisado.
- `oraculo/salidas/s2.json` lo reescribió `fin_run_fixtures.py` al reejecutarlo, con el mismo sha256 (`ba2cbcb1…dce7`).

## Verificación
- (a) `go vet` y `gofmt -l` en `motor_fin` y `cmd/motores`: limpios. `go test ./internal/motor_fin/... -count=1 -cover`: PASS, **88.8 %**; 298 PASS, 0 FAIL, 0 SKIP. `cmd/motores`: PASS, 78.4 %. A nivel módulo, `gofmt -l .` y `go vet ./...` fallan solo en `internal/precios` e `internal/recetas`, que son de otros agentes (B1).
- (b) `fin_run_fixtures.py`: 27/27. `go run ./cmd/motores s2 dump`: exit 0, salida idéntica byte a byte a la del oráculo. `compara.py s2`: 2318/2318 en los tres cruces, 159/159, 43/43, 867/867 provisionales, exit 0. Mi script: motor 159/159 y oráculo 159/159 contra `model.json`.
- (c) Independencia: sin lecturas cruzadas. `trim` usa `esperado` solo como forma.
- (d) Los sha256 de `model.json` (`1a570666…bfd4`) y `verify3.json` (`d86087f3…5041`) no cambiaron en ningún momento. Hubo 0 escrituras fuera del repositorio: `verify3.json` de `reporte-compra-latex` sigue con fecha del 10-oct.
- (e) Criterios S2-01…S2-16: la tabla con evidencia propia está en el informe. CUMPLEN 14, contando S2-03 y S2-12 que llevan observación; PARCIALES 2 (S2-05 y S2-16).
- (f) K2, K8, K19, K20, K21 y K39: la fixture propia `rev-propio` coincide con mis cuentas a mano en sus 248 hojas, con 0 diferencias en el motor y en el oráculo. `rev-cadencia-apagada` coincide en 326 hojas con la proyección por defecto. El comando `Decimal` del criterio imprime `219.88 629.35`.
- (g) Ninguna cifra fue ajustada.

## Supuestos
- Para M3 tomé la Verifica de S2-05 al pie de la letra («prueba de propiedad sobre todas las fixtures»). La propiedad sí se cumple, así que el hallazgo es de cobertura de la prueba y no de cifras, y lo califiqué como medio.

## Bloqueos
Ninguno nuevo; un revisor no los abre. La corrección de M1 sugiere abrir `docs/estado/bloqueos/S2-ventanas.md`.

## Pendientes
Ninguno de esta ronda. La ronda 2 debe verificar M1, M2 y M3.

## Solicitudes fuera de alcance
1. Dueños de `internal/precios` e `internal/recetas`: dejar limpios `gofmt` y `go vet` (B1).
2. Interfaz S2 §2.12 y §3.2: decidir qué pasa con un evento anterior a la primera ventana y con uno en un hueco entre ventanas (M1).
3. `reglas-negocio.md` RN-18: agregar la fórmula del plan de aportaciones al cerrar `S2-metas` (B2).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: S2 no tiene valores K31 pendientes.

## Consumo aproximado
- 44 acciones de herramienta.
- Leídos: ~190 KB (protocolo 6, decisiones §B y §K 16, criterios 35, interfaz por partes ~45, reportes 13, salidas y código por partes ~75).
- Escritos: ~35 KB (informe 13, este reporte 4, scripts y fixtures temporales ~18, ya borrados).
