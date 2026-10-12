# WF-3a · semilla · Reporte

Agente `semilla` (Sonnet) · 12-oct-2026 · HEAD «[Puerta 2 bis] K53» (sin commit). Convenciones: `CONVENCIONES-backend.md`. Criterios que tocan: C3-09 (cifras con la semilla) y S6 (subcomando `semilla`).

## Hecho (HECHO)
- `internal/semilla`: cargador idempotente de `datos/semilla/` a una base migrada, **una sola transacción** (todo o nada), cada alta con `audit_log` (`request_id = semilla`, actor nulo).
- Personas Axel (titular), Ximena (co-titular), Karen (externa, `pays_share_directly`), y las 2 cuentas (argon2id). Contraseñas: **nunca en el repositorio**; ver «Contraseñas iniciales».
- Reglas de reparto Hogar (por defecto), Karen y Conjunto 50/50 (RN-03); categorías Despensa (`groceries`, Conjunto 50/50) y Comida fuera (`eating_out`). Las demás categorías las captura Axel (U4 P-26.04).
- Tiendas de `data.json` (Bodega Aurrera / Vallejo y Walmart / Acueducto de Guadalupe), 9 grupos, 160 productos, 186 presentaciones, 60 recetas con 561 ingredientes y tiempos (validados con `recetas.LoadCatalogFile` y `LoadRecipeBook`).
- La compra del 10-oct-2026: 1 gasto de $3,094.50 (crédito, devengo bimestral, categoría Despensa), **36 renglones / 61 piezas** con cuenta por artículo (axel → Hogar, karen → Karen, conjunto → Conjunto 50/50 y 3 renglones **mixtos** con regla de un solo uso por cantidad: frijoles 50 %, salchicha 50 %, elote 60 %), partes congeladas Axel $1,903.00 y Karen $1,191.50.
- Crédito Kueski (línea revolvente, límite $10,000, quincenal, 15 días, lectura B pendiente de confirmar) con sus tasas a 2 y 4 quincenas (fracción exacta y decimales capturados) y las **4 despensas** con sus 16 pagos y sus partes, **calculados por `motor_cred.QuotePurchase`** y cotejados renglón por renglón con `oraculo/model.json` (un cambio de un centavo detiene la carga con «HALLAZGO» y los dos valores).
- PROFECO: solo la muestra, con `profeco.Importer` en modo producción. Tabla nutricional: solo renglones `confirmado` (hoy 0 de 150). `gastos_fijos_ingresos.csv`: no existe; aviso e importador enganchable (`DefaultImportPlannedCSV`). Promociones: `data.json` no trae.
- Prueba de cruce C3-09 como último paso de la carga y como `semilla.Verify`: disponible inicial A $10,000.00 / B $6,905.50, pago del 30-oct-2026 $973.75 y deuda viva al 10-abr-2027 $6,293.44, calculados con `AvailabilityTimeline` y `LiveDebt` sobre la base sembrada.
- Subcomando `finanzas semilla` instalado (`cmd/finanzas/semilla.go`, una línea: `seedRunner = seed.Run`).

## Archivos
Creados: `backend/internal/semilla/{semilla,personas,catalogo,datos,compra,otros,verificar}.go`, `backend/internal/semilla/semilla_test.go`, `backend/cmd/finanzas/semilla.go`, `backend/cmd/finanzas/semilla_test.go`, `docs/estado/wf3/semilla.md`, bloqueos `docs/estado/bloqueos/S6-semilla-despensas-futuras.md`, `S6-semilla-sucursales-municipio.md`, `S6-semilla-catalogo-alias.md`. No se modificó `main.go` ni migraciones ni motores.

## Contraseñas iniciales (cómo se dan y cómo las cambia Axel)
1. Antes de `finanzas semilla`, exportar `FINANZAS_SEED_PASSWORD_AXEL` y `FINANZAS_SEED_PASSWORD_XIMENA` (mínimo 12 caracteres, `auth.password_min_length`), **o** poner un archivo FUERA del repositorio, modo `600`, con renglones `axel=…` y `ximena=…` y exportar `FINANZAS_SEED_PASSWORD_FILE=/ruta`. El programa rechaza un archivo legible por otros o dentro de un repositorio git. Sin valores, la carga falla sin sembrar nada. Nunca se imprimen ni van a `audit_log` (solo el hash).
2. Solo se usan para cuentas **nuevas**; sembrar de nuevo no las cambia.
3. Axel las cambia al entrar: Ajustes → Cuenta y roles (P-26.01) → «Cambiar mi contraseña» (`changeOwnPassword`), y Ximena igual con la suya.
4. En pruebas se usa `testutil.FakeCredential`.

## Verificación
```
cd backend && go test -count=1 -race -v ./internal/semilla/   → 9 pruebas PASS (≈38 s con -race)
cd backend && go test -count=1 -cover ./internal/semilla ./cmd/finanzas → ok · 79.4 % de internal/semilla
go build ./... · go vet ./internal/semilla ./cmd/... · gofmt -l internal/semilla cmd → sin salida
finanzas semilla -db datos/tmp/semilla.db (1.ª vez) → 3 personas, 2 cuentas, 160 productos, 186 presentaciones, 60 recetas, 1 gasto, 36 renglones, 4 compras, 16 pagos; 6 «cotejo OK»
finanzas semilla (2.ª vez, otras contraseñas) → 0 nuevos en todas las entidades; mismas 1211 filas de audit_log
sqlite3: 16 pagos programados; primeros: 2026-10-30 97375, 2026-11-15 97375
```
- **C3-09** (`TestSembrarCargaTodoYCruzaConElOraculo`): A 1,000,000 / B 690,550 centavos; primer pago 97,375 el 2026-10-30; deuda viva al 2027-04-10 629,344 = oráculo.
- **Idempotencia** (`TestSembrarDosVecesNoDuplica`): conteo de **todas** las tablas igual tras la 2.ª carga; el hash de la contraseña no cambia; sin contraseñas pasa si las cuentas ya existen.
- Todo o nada (`TestSinContrasenasNoSiembraNada`, `TestOraculoDistintoEsHallazgoYNoSiembra`: 973.76 en vez de 973.75 → «HALLAZGO … 97376», 0 filas). Contraseña nunca en claro (`TestContrasenasInicialesNuncaEnClaro`). Archivo de claves (`TestArchivoDeContrasenas`). Nutrición solo confirmada (`TestNutricionSoloConfirmada`). CSV de partidas (`TestPartidasCSVAvisaOImporta`).
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: las migraciones (goose up/down/up) y `go build`/`go vet` pasan; **se detiene en `gofmt -l`** por dos archivos de OTROS agentes, no míos: `backend/internal/api/sistema/adjuntos.go` y `backend/internal/reportes/vistas.go` (ver solicitudes). Mis archivos están formateados.

## Supuestos
1. Las despensas 2 a 4 se siembran como compras a crédito activas sin gasto (lo exigen C3-09 y `decisiones.md` «4 despensas y pagos»): PROVISIONAL, bloqueo `S6-semilla-despensas-futuras`.
2. Tasas a 2 y 4 quincenas con `source = manual` (como `canontest.Seed`), no `quote`; la lectura de disponible queda B sin confirmar (decisión C: «por defecto B, pendiente confirmar»).
3. `recipe.instructions` = los pasos unidos con salto de línea, sin numerar (C1 los devuelve como un solo texto).
4. Un renglón del ticket por producto (36), no por línea de `cuentas` (39): lo que `data.json` reparte entre dos cuentas se vuelve un renglón mixto con regla de un solo uso por cantidad (RN-03). `unit_price_cents` solo cuando el importe divide exacto.
5. La sexta lata de atún (`atun_sexta_pieza_bodega`) es la única diferencia admitida entre `items_identicos` y `cuentas`; cualquier otra detiene la carga.
6. Sin papel de contado (`model.json: papel`): no es compra a crédito y su gasto lo captura Axel; no entra a esta semilla.
7. Idempotencia por llave natural (código, nombre, fecha + importe): si Axel renombra o edita algo sembrado y se vuelve a sembrar, puede reaparecer. La semilla es de arranque, no de sincronización.
8. `profeco_name` de las cadenas («Bodega Aurrera», «Wal-mart») viene de `cadena_comercial` de la muestra: PROVISIONAL.

## Bloqueos
`docs/estado/bloqueos/S6-semilla-despensas-futuras.md` · `S6-semilla-sucursales-municipio.md` · `S6-semilla-catalogo-alias.md`. Referencias: `S4-empates-sin-confirmar` (150 renglones de nutrición «propuesto» no entran, Q2-43 = A), `S3-iva-lectura` (sin cambio).

## Pendientes
- `gastos_fijos_ingresos.csv` (lo entrega Axel) y su importación en la carga: falta que Ajustes exporte el importador (solicitud 1).
- Precios de PROFECO: 0 cargados hasta que haya municipio en las sucursales y `profeco/vinculos.csv` confirmado.
- Kcal de las recetas: 0 hasta que Axel confirme `nutricion/tabla.csv` (Q2-90).
- Cobertura 79.4 % (sin umbral K25 para este paquete; `semilla` no está en `K25_RAICES`).

## Solicitudes fuera de alcance
1. `backend/internal/api/ajustes/importar.go`: exportar el importador de partidas, p. ej. `func ImportCSV(ctx context.Context, db dominio.Querier, today dominio.Date, b []byte) (int, error)` (mismo cuerpo que `importPlannedItems`); y en `backend/cmd/finanzas/semilla.go` agregar `seed.DefaultImportPlannedCSV = func(ctx context.Context, db dominio.Querier, path string) (int, error) { b, err := os.ReadFile(path); if err != nil { return 0, err }; return ajustes.ImportCSV(ctx, db, dominio.Today(dominio.SystemClock{}), b) }`. Hasta entonces la semilla solo avisa.
2. `gofmt -w backend/internal/api/sistema/adjuntos.go backend/internal/reportes/vistas.go` (dueños: auth-sistema y s7-pdf): rompen `make verify-all` en `gofmt -l`.
3. `docs/estado/bloqueos/INDICE.md`: agregar los 3 bloqueos `semilla-*`.
4. `ops/runbook.md` (WF-4): documentar el orden «migrar → exportar contraseñas → `finanzas semilla` → cambiar contraseñas en Ajustes».
5. Opcional: copiar `oraculo/model.json` a `datos/semilla/model.json` (plan §3: copias); mientras no exista, la semilla lo busca en `../../oraculo/model.json`.

## Propuestas para Axel (PROPUESTA, no decidido)
- Confirmar el reparto de despensas futuras (A, B o C del bloqueo) antes de capturar la partida «Despensa» en Ajustes.
- Dictar municipio y estado de las dos sucursales para que la muestra de PROFECO cargue algo.
- Revisar y confirmar `nutricion/tabla.csv` (150 renglones) para que las recetas tengan kcal.
- Guardar las contraseñas iniciales en un administrador de contraseñas y cambiarlas en el primer inicio de sesión.

## Paquetes instalados
Ninguno (sin `sudo`; ningún módulo de Go nuevo: solo paquetes del propio proyecto).

## Consumo aproximado
≈ 56 acciones de herramienta · ≈ 330 KB leídos · ≈ 75 KB escritos (código ≈ 52 KB, bloqueos ≈ 6 KB, reporte ≈ 9 KB, pruebas ≈ 14 KB).
