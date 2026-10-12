# correccion-contrato-S3-r1 · reporte final (WF-1, 11-oct-2026)
Objeto: `docs/contratos/motores/S3-interfaz.md` v1.0 → **v1.1**. Corrige los 8 hallazgos medios (H1–H8) de `docs/estado/wf1/REVISION-contrato-S3-r1.md`. No reabre decisiones §A–§K; ninguna cifra canon cambia.

## Hecho
- **H1 (S3-01, S3-02, S3-09) · fixtures, CLI y oráculo alineados a `oraculo/fixtures/FORMATO.md` v1.0:** carpeta `oraculo/fixtures/s3/<caso>.json` y `_datos/`; llaves `formato, motor, caso, descripcion, fuente [{tipo, ruta}], provisional {bloqueo, que}, entradas {datos, consultas}, esperado`; CLI `backend/cmd/motores s3 dump --fixtures --out` con salida `salida/1` y errores dentro (código 0/1/2); oráculo `cred_run_fixtures.py` → `oraculo/salidas/s3.json`; comparación `python3 -I oraculo/compara.py s3`. Se retiraron `$model`, `--escribir`, `cred_compara.py`, `cmd/finanzas … --fixture|--canon` y la carpeta plana (§2, §9, §10, §11). Nuevas: funciones de consulta (§9.2), vistas con claves de `model.json` (§9.3), códigos de error de salida (§9.1), procedencia de cada `esperado` (§10). R-05 y R-06 con fuente `data.json` (`datos/semilla/data.json`, `financiamiento.dos_quincenas`) y sin el SUPUESTO «contra plan §1.2».
- **H2 (S3-03, S3-11, S3-12, S3-13) · tipos declarados con llave JSON:** `Reading`, `CreditState`, `PurchaseState`, `Input`, `DecodeInput` (§5.7); etiquetas `json` en `ShareInput`, `PersonShare`, `Composition`, `Credit`, `TermRate`, `Promotion`, `PurchaseInput`, `Purchase`, `ScheduledPayment`, `PaymentShare`, `CashEvent`, `PaymentRecord`, `TimelineEvent`, `Timeline`, `StatementInput`; `CalendarEntry`, `FortnightTotal`, `Calendar`, `CashWindow`; `DebtLine`, `DebtSummary`, `AmountQuote`, `SensitivityInput`, `SensitivityRow`; `Request`, `Constraints`, `Weights`, `PersonCap`, `BasketQuote`, `Scenario`, `PersonPeak`, `Rejected`, `ExcludedCredit`, `Optima`, `Result`. Dos ejemplos JSON de `entradas` (§10.3): línea de tiempo con datos en línea y comparador.
- **H3 (S3-01) · `TestCanon85` con aserción exacta para las 30 filas S2:** función nueva `CashFlow` (efectivo por ventana, §6.1) y tabla R-nn → función S3 → aserción (§11.1), sin `t.Skip`.
- **H4 (S3-14, K17) · conjunto de la normalización declarado:** «generados» en el sentido de la decisión C («genera escenarios … que respetan las restricciones duras») = `valid`; los descartados cuentan en `evaluated_count` pero no en mín/máx (§8.5, §14). Solicitudes a RN-11 y C1.
- **H5 (S3-03, S3-12) · lectura A con `From` distinto del ancla:** A(From) = límite − Σ capital de compras con ancla < fecha ≤ From + Σ capital repuesto con ancla < reposición ≤ From; sin ancla con lectura A → `ErrInvalid` (CHECK de `credit`); con B, ancla = From y `AAnchorDefaulted`; From < ancla → `ErrInvalid` (§5.2). El comparador usa `From = Today` y toma el valor vigente en cada día de la ventana (§8.3). Prueba nueva `TestSeriesABDesdeOtraFecha` y fixture `s3-series-desde-otra-fecha`.
- **H6 (S3-13) · objetivo 2 unívoco:** quincenas consecutivas de `FortnightOf(compra)` a `FortnightOf(último pago)`, inclusive, con las de total 0; fórmula entera exacta `StdDevMilliCents` = ⌊(⌊isqrt(4N)/n⌋ + 1)/2⌋ con N = 10⁶(nΣx² − (Σx)²) y ejemplo; desempate del objetivo 2 con `OverCap`, alineado con el bloqueo `S3-objetivos-medida` (2) (§8.4, §8.6). Prueba nueva `TestObjetivoParejo`.
- **H7 (S3-11) · sin doble conteo de tarjeta:** el corte sustituye a los pagos programados de su vencimiento; sin corte, entran los programados; PROVISIONAL con `S3-formula-tarjeta` (§6); `TestUnificado` y fixtures `s3-unificado`, `s3-unificado-sin-corte`, `s3-unificado-duplicado`.
- **H8 (S3-12, K43, K44, S3-15) · tiempo máximo probable en un motor puro:** `TimeoutFromSetting(seconds)` y la tabla `ErrorResponses` (`ErrTimeout` → 503 `unavailable`, `ErrInvalid` → 422 `validation_failed`) como dato sin `net/http`, con `ResponseFor` (§8.7); pruebas `TestTiempoClaveSembrada` (lee el valor sembrado de 00003 como texto) y `TestRestriccionesTiempo` (0 → timeout y 503; 5 → fin, líneas `S3-TIEMPO <caso> <ms> ms` que el optimizador copia a `S3-escenarios.md`); caso de carga `s3-comparador-carga`; división explícita de lo que se verifica en S3 y en WF-2 (§11).

## Archivos
- Modificado: `/home/deepsec/finanzas-hogar/docs/contratos/motores/S3-interfaz.md` (46,597 → 90,865 bytes; §12 idéntico byte a byte).
- Creado: `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-contrato-S3-r1.md` (este reporte).
- Temporales en `datos/tmp/corr-s3r1/` (partes, scripts de ensamblado y verificación, copia de v1.0): borrados al terminar.

## Verificación
El artefacto es un contrato sin código (`backend/internal/` solo tiene `api`; no hay `motor_cred` que compilar ni pruebas Go que correr). Se verificó así:
- `python3 -I datos/tmp/corr-s3r1/verif.py oraculo/model.json datos/semilla/data.json` (implementación propia con `Fraction` de las reglas escritas: `sched`, `SplitTrunc`, `SplitHalfUp`, partes, línea de tiempo, devengo, efectivo por ventana) → **86 comprobaciones, 0 fallas**: las 4 despensas (totales y 16 pagos con fecha, monto, Axel, Karen y capital) iguales a `model.json`; **las 30 filas S2 de §11.1** (R-37…R-44 desde `Sensitivity` N = 400000; R-45…R-50 devengado, efectivo, diferencia acumulada, pagos y papel por ventana; R-51, R-52, R-56…R-60; R-53…R-55 con `LiveDebt`; R-68…R-73 paquete, fecha, persona, forma y monto); R-05/R-06 desde `data.json` («17.27», 53450, 362900); los 17 eventos de `credito` con A y B; **propiedad de H5** para las 16 fechas de evento posteriores al ancla y los dos ejemplos (909450/600000 y 809450/500000); fórmula de la desviación: ejemplo 0/100/200 → 81650 y 2,000 series aleatorias contra la definición exacta de HALF_UP(1000·√var); `tasa` y `tasa_pct` binary64 desde la fracción (y `float(tasa) × 100` **no** coincide: hallazgo 4 de §14).
- `python3 -I datos/tmp/corr-s3r1/estructura.py …` → 14 secciones en orden; **12 bloques Go aceptados por `gofmt -e`** (etiquetas `json` bien formadas) y **2 bloques JSON válidos**; §12 sin cambios; 23 nombres de caso con la forma `s3-<tema>[-<variante>]` de FORMATO §1; los 14 tipos que pedía H2 declarados; restos de la forma vieja solo en la frase que la retira (`$model` 1, `cred_compara` 1, `cmd/finanzas` 1, `--escribir` 1, `--canon` 2: la frase de retiro y «no hay modo `--canon`»).
- `sha256sum oraculo/model.json` → `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` (intacto, S3-16).
- Criterios cubiertos por esta corrección y su evidencia en el contrato: S3-01 (§6.1, §10, §11.1; 30 filas S2 verificadas), S3-02 (§10.2 `compara.py s3`, §12 intacto), S3-03 (§5.2 ancla; propiedad verificada), S3-09 (§10.2 `cred_run_fixtures.py` con FORMATO), S3-11 (§6 sustitución del corte), S3-12 (§8.3, §8.7, §11), S3-13 (§8.4 (2), §8.6), S3-14 (§8.5), S3-15 (§8.7 sin `net/http`).

## Supuestos
- Lectura declarada, no supuesto: «escenarios generados» (K17, S3-14) = los que respetan las restricciones duras, por el texto de la decisión C; queda escrita y justificada en §8.5 con solicitud a RN-11 y C1.
- Diseño de interfaz sin cifras ni esquema: nombres de las funciones de consulta (§9.2), códigos de error de salida (`invalid`, `no_rate`, …), campos nuevos de salida (`Purchase.rate_display`, `SensitivityRow.monthly_purchase_cents`, `Timeline.a_anchor`/`a_anchor_defaulted`, `Result.excluded`) y de entrada (`Credit.holder`, `CashEvent.person`, argumento `reading` de `linea_tiempo`), regla «sin ancla con lectura B → ancla = From marcada» y «From < ancla → `ErrInvalid`».
- La regla del corte que sustituye pagos (reparto de capital, interés y partes) va PROVISIONAL con el bloqueo existente `S3-formula-tarjeta`; no entra en cifras canon.
- «Decenas de candidatos por caso» es una pauta de diseño de los casos de prueba, no un parámetro operativo.

## Bloqueos
- Ninguno nuevo. Siguen abiertos (sin cambio de archivo): `docs/estado/bloqueos/S3-formula-tarjeta.md` (ahora también cubre la sustitución del corte en el calendario), `S3-objetivos-medida.md`, `S3-convencion-prestamo.md`. `S1-parametros-avisos` no toca a S3.

## Pendientes
- Hallazgos bajos de la r1 que no estaban en el encargo: H9 (iii) `Rate.Display()` con `Decimals = -1` y (iv) datos de la fixture «100000 a 7» de `s3-prestamo-frances`; H10 (cifras de los casos del comparador y desempate del caso 2, a cargo de `S3-escenarios.md`); H11 (`oraculo/verify3.py:89` sin dueño en §2); H12 (solicitudes a `reglas-negocio.md`). H9 (i), (ii) y (v) quedaron resueltos de paso en §5.2, §9.3 y §3.
- Revisión r2 del contrato por el revisor de motores.

## Solicitudes fuera de alcance
1. `docs/contratos/reglas-negocio.md` RN-11: decir que la normalización mín-máx es «entre los escenarios generados (los que respetan las restricciones duras, decisión C)»; motivo: H4.
2. `docs/contratos/openapi.yaml` (C1): `Scenario.score_bp` con la misma aclaración; `Scenario` + `store_id`, `std_dev_milli_cents` (milésimas de centavo, la medida de K17; sustituye el nombre `fortnight_payment_stddev_cents` de la solicitud 1 de `contrato-S3.md`), `over_cap` [`{person, peak_cents, cap_cents}`] y `provisional`; `ScenarioResult` + `excluded` [`{credit_id, reason}`] y `provisional`; motivo: H2, H4, K17.
3. `docs/estado/bloqueos/S3-objetivos-medida.md` (2): escribir el conjunto exacto de quincenas de §8.4 y que `over_cap` desempata solo el óptimo del objetivo 2; motivo: H6 (que bloqueo y contrato digan lo mismo).
4. `docs/estado/bloqueos/S3-formula-tarjeta.md`: añadir la pregunta del corte que sustituye pagos en el calendario unificado (reparto de capital, interés y partes); motivo: H7.
5. `docs/criterios/S3.md`: S3-12 «Verifica» con la división de §11 (en S3: valor sembrado, `TimeoutFromSetting`, `ResponseFor`, 0 y 5 s, líneas `S3-TIEMPO`; en WF-2: lectura en vivo de `app_setting` y HTTP 503, C1-23); S3-02 con `python3 -I oraculo/compara.py s3`; S3-09 con `cred_run_fixtures.py --fixtures oraculo/fixtures/s3 --out oraculo/salidas/s3.json`; motivo: H1, H8.
6. `oraculo/fixtures/FORMATO.md` (integrador): en §3, que el ×100 de `tasa_pct` se aplica a la fracción exacta antes de redondear a binary64 (verificado: binary64(8005000/309450) = `model.json`; `float(tasa)·100` no); confirmar el dueño de `oraculo/compara.py` y el binario único `backend/cmd/motores` con `s2` y `s3`; motivo: H1.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- Sin propuestas nuevas. Se mantiene la del contrato v1.0: `optimizer.date_window_days` = 7 días (K31, PROVISIONAL hasta la puerta 2).

## Consumo aproximado
- Acciones de herramienta: 35.
- Leído: ~178 KB (protocolo, FORMATO, revisión, decisiones §C y §K, criterios S3, S3-interfaz v1.0 completo por partes, S2-interfaz por partes, extracto C1 de escenarios, bloqueos S3, salidas de verificación).
- Escrito por mí: ~97 KB (partes de las secciones ~72 KB, scripts ~13 KB, progreso y reporte ~12 KB). El archivo final de 90.9 KB lo ensambló un script con esas partes y el texto sin cambios de v1.0.
