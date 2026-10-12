# motor-S3-nucleo-go · reporte final (WF-1, reanudación 11-oct-2026)

Agente: motor-S3-nucleo-go (Opus). Contrato: `docs/contratos/motores/S3-interfaz.md` v1.4 (ronda 4). Criterios: `docs/criterios/S3.md` 1.2.

## Hecho
- Revisé los 12 parciales del intento anterior (compilaban y pasaban vet; sin pruebas) y los ajusté al contrato v1.4:
  - `PurchaseInput.RateDecimals` (K12, §5.1): con `QuotedInterest` los decimales salen de la compra (null → −1); los de un `TermRate` del mismo plazo ya no se usan; un valor < 0 da `ErrInvalid`.
  - `Rate.Display()` con `Decimals` = −1 devuelve `""` (antes inventaba 4 decimales).
  - Listas `Provisional` de §14.1, sin repetidos y en orden por bytes (`provList`): `S3-serie-sin-intereses` cuando los intereses son 0 y n ≥ 2 (línea, MSI, tienda); `U4-iva-credito` en el interés y el mínimo con fórmula y en el préstamo sin IVA; `CardMinimum` nuevo, que devuelve la lista del mínimo.
  - Quité el ayudante `RatFloat`, que no se usaba.
- `canontest/` (paquete de apoyo, solo lo usan las pruebas): `Seed` = seedFromModel y `Sensitivity` = sensitivityFromModel (§10.0). Leen `model.json` y `data.json` al correr y no tienen cifras canon escritas.
- `evaluar/evaluar.go`: `Evaluate` con las 13 funciones de §9.2 y las 7 vistas de §9.3; una `funcion` desconocida da `invalid` y un campo desconocido da `unknown_field`.
- `cmd/motores/s3.go`: `register("s3", evalS3)` contra la API de `main.go`, que ya existía y no toqué.
- Pruebas con los nombres de §11: 29 en `motor_cred`, `TestEvaluar` y `TestOraculoFixtures` en `evaluar`, y `TestCanontest`.

## Archivos
- Modificados (parciales): `backend/internal/motor_cred/{money,state,revolving,card,loan,store,views,errors}.go`
- Sin cambios: `backend/internal/motor_cred/{date,availability,calendar,debt}.go`
- Nuevos: `backend/internal/motor_cred/canontest/canontest.go`, `canontest/canontest_test.go`, `backend/internal/motor_cred/evaluar/evaluar.go`, `evaluar/evaluar_test.go`, `backend/internal/motor_cred/{canon_test,engine_test,timeline_test}.go`, `backend/cmd/motores/s3.go`, este reporte

## Verificación
- `gofmt -l internal/motor_cred cmd/motores`: vacío. `go vet ./internal/motor_cred/... ./cmd/motores/`: limpio. `go test ./cmd/motores/`: ok.
- `go test ./internal/motor_cred/ -count=1 -cover -v`: 29 de 29 PASS, 0 SKIP, cobertura **93.7 %**. `canontest`: 90.0 %. `evaluar`: 90.5 %; `TestEvaluar` PASS.
- `TestOraculoFixtures` FALLA, como pide el contrato (§10.0: falla con nombre y autor, nunca `t.Skip`), solo porque faltan los 9 archivos del optimizador S3: `_datos/comparador-caso-{1,2,3}.json` y `s3-comparador-{caso-1,caso-2,caso-3,sin-restricciones,tiempo,carga}.json`. Las 21 fixtures del agente de regresión, más `_datos/kueski-semilla.json`, corren sin una sola diferencia contra su `esperado`.
- S3-01: `TestCanon85` da 85 de 85 subpruebas R-01…R-85 PASS (las 30 filas S2 con las aserciones de §11.1).
- S3-02: `TestCanonHojas437` da «437 hojas comparadas» (parametros 9, despensas 144, credito 143, sensibilidad 45, disposicion 68, interes_por_monto 28). `kueski-semilla.json` es igual, como JSON, al `Input` de seedFromModel, y los argumentos de `s3-sensibilidad` y `s3-interes-por-monto` son iguales al insumo de sensitivityFromModel.
- S3-03: 28 de 28 valores; con la lectura omitida sale «B» y la bandera en verdadero; 16 fechas de inicio probadas en `TestSeriesABDesdeOtraFecha`.
- S3-04: 16 de 16 fechas, incluido el 28-feb-2027; además el 29-feb-2028.
- S3-05: el capital repuesto suma el capital de la compra; K15 con la lista completa del ejemplo 1; abono parcial y pago vencido.
- S3-06/S3-07/S3-08/S3-09/S3-10/S3-11/S3-19/S3-20: `Tasa*`, `PromocionFronteras`, `Tarjeta*`, `PrestamoFrances`, `Tienda*`, `Unificado`, `Reparto`, `ProvisionalListas` y `DespachoCotizacion`: PASS.
- CLI: `go build -o ../datos/tmp/motores ./cmd/motores && datos/tmp/motores s3 dump --fixtures oraculo/fixtures/s3 --out datos/tmp/motor-s3.json` termina con código 0: `salida/1`, `s3`, 21 casos. Las 12 consultas con error son las esperadas (`no_rate` ×4, `promotion_not_applicable` ×2, `missing_capture` ×3, `invalid` ×2, `unsupported` ×1).
- `python3 -I oraculo/compara.py s3`: oráculo/esperado, motor/esperado y oráculo/motor dan 893 de 893 cada uno; model.json (ramas de s3) 437 de 437; data.json 9 de 9; provisionales 695 de 695; **diferencias: 0**.
- `sha256sum oraculo/model.json` = `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` antes y después (S3-16).
- S3-15: 0 líneas con `net/http` fuera de las pruebas; `highs/ortools` en go.mod y go.sum: 0. S3-17: 0 `float32/float64` fuera de las pruebas; en las pruebas, solo `ratioFloat` y `modelFloat` de `canon_test.go`, para las dos hojas flotantes `tasa` y `tasa_pct` (§1.1). `t.Skip`: 0.
- Temporales borrados: `datos/tmp/motores` y `datos/tmp/motor-s3.json`.

## Supuestos
1. `canontest/` no está en la tabla de §2. Hace falta para que `motor_cred` y `evaluar` compartan seedFromModel sin ciclo de importación. Las pruebas canon van en el paquete externo `motor_cred_test`.
2. `TestOraculoFixtures` vive en `evaluar`, junto al despacho. El contrato dice «ambos»: el comparador hará la suya.
3. La consulta `comparar` responde `unsupported` mientras no exista `motor_cred/comparador/`; nunca devuelve un número.
4. No escribí fixtures en `oraculo/fixtures/s3/`. La tarea lo permitía solo si la carpeta estaba vacía; el contrato v1.4 §10.0 prohíbe al núcleo escribir en `oraculo/`, la carpeta ya tenía `_datos/` y el agente de regresión las escribió después.
5. §5.6: una línea semanal o mensual sin `FirstPaymentDate` da `ErrUnsupported`.
6. Dos cortes de una misma tarjeta con el mismo `due_date` también dan `ErrInvalid` (lectura estricta de «sin duplicados», §6).
7. Agregué `CardInterest` y `CardMinimum`, que devuelven además la lista provisional; las firmas del contrato (`CardPeriodInterest`, `CardMinimumPayment`) se conservan.
Ninguno cambia cifras, esquema ni endpoints.

## Bloqueos
Ninguno nuevo. Las salidas llevan las marcas de los bloqueos del contrato: `S3-formula-tarjeta`, `S3-convencion-prestamo`, `S3-origen-cargo-normal`, `S3-serie-sin-intereses` y `U4-iva-credito`.

## Pendientes
- `backend/internal/motor_cred/evaluar/evaluar.go`, entrada `"comparar"`: importar `finanzas-hogar/internal/motor_cred/comparador` y llamar `Compare` con `timeout_seconds` y `request` (§9.2) cuando ese paquete exista (integrador u optimizador S3).
- `TestOraculoFixtures` pasará cuando el optimizador escriba sus 9 archivos y se conecte `comparar`.
- El integrador corre §10.2 completo al juntar (`cred_run_fixtures.py`, `s3 dump`, `compara.py s3`).

## Solicitudes fuera de alcance
- `docs/contratos/motores/S3-interfaz.md` §2 y §11: agregar `motor_cred/canontest/` (paquete de apoyo de pruebas) y anotar que `TestOraculoFixtures` del núcleo está en `evaluar`. Motivo: alcance disjunto y sin ciclo de importación.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: el núcleo no tiene valores operativos de K31. El tiempo del comparador es del optimizador.

## Consumo aproximado
~42 acciones de herramienta; ~285 KB leídos; ~100 KB escritos.
