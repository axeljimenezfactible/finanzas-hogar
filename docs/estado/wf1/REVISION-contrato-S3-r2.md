# Revisión del contrato S3 · ronda 2 de 3

Revisor: `revisor-contrato-S3-r2` (independiente; no escribió nada de lo revisado y no corrige nada) · 11-oct-2026 · WF-1
Objeto: `docs/contratos/motores/S3-interfaz.md` v1.1 (90,865 bytes, leído completo por partes), con el reporte `docs/estado/wf1/contrato-S3.md` y la corrección `docs/estado/wf1/correccion-contrato-S3-r1.md`.
Vara de medida: `docs/criterios/S3.md` 1.2, decisiones §C y §K. Para la coherencia entre artefactos usé `oraculo/fixtures/FORMATO.md` 1.1, RN-11 y RN-21 de `reglas-negocio.md`, el extracto C1 `creditos-b.yaml` (`Scenario`, `ScenarioResult`) y el esquema (`credit`, 00003).

## Veredicto
**NO APROBADO.** Hay 0 hallazgos altos, 4 medios y 7 bajos. De los 8 hallazgos medios de la ronda 1 cierro 7 (H2–H8). H1 queda **reabierto en parte**: el contrato se alineó a `FORMATO.md` 1.0, pero `FORMATO.md` ya está en 1.1 (17:55) desde antes de la corrección (18:06), y la 1.1 fija para S3 tres cosas que el contrato contradice (M1).

Todas las cifras que el contrato fija reproducen `model.json` y `data.json` al centavo. Lo comprobé con mis propios scripts: 296 comprobaciones y 0 fallas (abajo).

## Verificación propia (scripts en `datos/tmp/rev-s3-r2/`, borrados al terminar)
1. **(a) Hojas de S3** (`hojas.py`: recorre `model.json` y casa cada hoja numérica con mis propias expresiones de los 42 renglones de §12). sha256 `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` (intacto). Por rama: `parametros` 9, `despensas` 144, `credito` 143, `sensibilidad` 45, `disposicion_por_despensa` 68, `interes_por_monto` 28 = **437**; S2 159; total 596. Resultado: 42 renglones, 437 cubiertas, **0 sin dueño, 0 con dos dueños** y ningún renglón con un conteo distinto al declarado. La columna «Fixture» de §12 cruza con la tabla de §10: todas las hojas caen en un caso cuyo `id` es una ruta de `model.json` (`parametros` reparte sus 9 entre `s3-kueski-despensa-1` y `s3-series-a-b`).
2. **Reglas del contrato contra el canon** (`calc.py`, con `Fraction` y las reglas tal como están escritas). Comprobé:
   - las 4 despensas: capital, intereses, total, y en los 16 pagos la fecha (`sched`), el monto (`SplitTrunc`), la parte de Karen (`SplitHalfUp` de su total), la de Axel (el resto) y el capital (`SplitHalfUp`); además los intereses por persona y el pago por quincena;
   - los 17 eventos de `credito` (tipo, A y B), las series A y B, los mínimos 809450 y 500000;
   - **la propiedad de §5.2 (H5) con `From` en cada una de las 16 fechas de evento posteriores al ancla**: el inicio y la cola de eventos son idénticos al canon; los ejemplos del contrato dan 909450/600000 (16-dic-2026) y 809450/500000 (10-feb-2027);
   - `interes_por_monto` (28 hojas). Confirmo el hallazgo 1: desde el total redondeado salen 32 contra 31 y 125869 contra 125868;
   - `sensibilidad` (45 hojas y el booleano);
   - **las 30 filas S2 con las aserciones de §11.1**: R-37…R-44, P1…P6 (devengado, efectivo y diferencia acumulada), R-51, R-52 (por la suma de los periodos y por la ventana 10-oct a 10-abr), R-53, R-56…R-60;
   - R-05 y R-06 desde `data.json`: HALF_UP(100·53450/309450, 2) = 17.27; 362900 y 53450;
   - `tasa_pct` = binary64(8005000/309450), mientras que `float(tasa)×100` **no** coincide (confirmo el hallazgo 4);
   - los ejemplos de §3, la promoción al 0 % (77362 ×3 + 77364), MSI 3 y tienda 3 (33333, 33333, 33334), y 2028-02-10 → 2028-02-29, 03-15, 03-30, 04-15;
   - la fórmula entera de la desviación (§8.4): 0/100/200 → 81650, y **5,000 series aleatorias contra HALF_UP(1000·σ) exacto con 0 diferencias**;
   - la normalización 100/200/400 → 1, 2/3, 0.

   Salida: `comprobaciones 296 fallas 0`.
3. **(b) Criterio → sección:** la tabla de §13 cubre S3-01…S3-20 y todas las secciones citadas existen.
4. **(d) K12–K18, K31, K43, K44:** el detalle está en la tabla de criterios. Lo que encontré son las marcas PROVISIONAL (M4) y los hallazgos bajos.
5. **(e) Casos del comparador:** §10.1 obliga a que cada caso traiga en la fixture los 4 valores de restricción (capturas de fixture, sin precarga). También obliga a ≥ 1 descartado por restricción, a uno sobre el tope y a la intención de los 4 óptimos distintos, verificados a mano antes de correr. El caso 3 cambia una restricción global con `Override`. Las cifras de los casos las fija el optimizador en `S3-escenarios.md` antes de correr, como pide S3-13. Riesgo de diseño: en el caso 2, el cargo normal de la tarjeta en la tienda de mejor U/C tiene 0.00 de intereses y el máximo rendimiento, así que los óptimos 1 y 4 coinciden salvo que el desempate por puntaje caiga fuera de esa tienda. El contrato lo reconoce («desempata por diseño»). No lo cuento como hallazgo, pero el revisor de S3-13 debe comprobarlo.

## Tabla de criterios
| Criterio | Estado | Evidencia propia |
|---|---|---|
| S3-01 (alta) | Cubierto; R-05/R-06 en desacuerdo con FORMATO 1.1 (M1) | §10, §11, §11.1; `calc.py`: despensas, eventos, series, mínimos, sensibilidad, `interes_por_monto`, R-05/R-06 y las 30 filas S2 cuadran |
| S3-02 (alta) | Cubierto | `hojas.py`: 437 / 0 sin dueño / 0 con dos dueños; `compara.py s3` en §10.2 |
| S3-03 (alta) | Cubierto (H5 cerrado) | propiedad de §5.2 en 16 fechas y 2 ejemplos; lectura omitida → B y bandera (§5.2, `s3-series-a-b`) |
| S3-04 (alta) | Cubierto | 16 fechas y 28-feb; 2028-02-29 |
| S3-05 (alta) | Parcial (M2) | §5.2 K15; la forma de los eventos con abonos registrados no está fijada |
| S3-06 (alta) | Cubierto | §5.1 orden cotizado → promoción → tabla → `ErrNoRate`; 400000 → 103474 |
| S3-07 (alta) | Cubierto | 4 fronteras en `s3-promocion-0`; 77362 ×3 + 77364 |
| S3-08 (alta) | Cubierto con PROVISIONAL `S3-formula-tarjeta` (bajo B3) | §5.3; `ErrMissingCapture` sin número |
| S3-09 (alta) | Cubierto con PROVISIONAL `S3-convencion-prestamo` (bajo B5) | §5.4 propiedades; fixture con ajuste de centavos |
| S3-10 (media) | Parcial (M3) | §5.5 no distingue «sin interés» de «sin datos» |
| S3-11 (alta) | Cubierto (H7 cerrado) | §6: el corte sustituye; llave sin duplicados; Σ por quincena |
| S3-12 (alta) | Cubierto con hallazgos (M3, M4) | §8.1–§8.3, §8.7, §11; `TimeoutFromSetting`, `ResponseFor`, clave sembrada 5 en `00003_puerta1b.sql` |
| S3-13 (alta) | Cubierto (H6 cerrado) | §8.4 (2) quincenas por `FortnightOf` con las de 0; §8.6 desempate con `OverCap` |
| S3-14 (alta) | Cubierto (H4 cerrado; bajo B7) | §8.5 declara y justifica el conjunto; empate = 1; 100/200/400 verificado |
| S3-15 (alta) | Cubierto | §1.5; §8.7 sin `net/http` |
| S3-16 (alta) | Cubierto | §10.2; sha256 intacto |
| S3-17 (alta) | Cubierto | §1.1–§1.3 (única excepción `float64`: comparar `tasa`/`tasa_pct` en la prueba) |
| S3-18 (alta) | Cubierto (bajo B4) | §11 |
| S3-19 (media) | Cubierto | §4: 30822/49228 y pagos 37493/59882 verificados |
| S3-20 (media) | Parcial (M4, bajo B1) | §14; bloqueos con referencias viejas |

## Hallazgos de esta ronda

### Medios
**M1 (S3-01, S3-02, S3-09) · El contrato se alineó a `FORMATO.md` 1.0 y contradice la 1.1, que rige para S3 sin excepciones.**
`FORMATO.md` 1.1 (17:55, anterior a la v1.1 del contrato, 18:06) fija para S3 tres cosas que el contrato no sigue:
- (a) **R-05 y R-06** van con `id` `data.json:financiamiento.dos_quincenas`, con una vista de claves literales (§2.1) y la razón con `decimales` (§3, §6). El contrato usa el `id` libre `dos-quincenas` con `cotizar_compra` (§10), así que `compara.py s3` no contrasta esas hojas con `data.json` (su línea sale «0 de 0») y solo las cruza `TestCanon85`. S2 sí usa `data.json:` y `decimales`.
- (b) La clave de la tasa con decimales capturados es `decimales` en FORMATO §3 y en S2; en el contrato es `decimals` (§3, §9.1).
- (c) `backend/cmd/motores/main.go` es del agente Go de S2 y nadie más lo edita; S3 agrega solo `backend/cmd/motores/s3.go` con `register("s3", evaluador)` y la firma `Evaluator` (FORMATO §4.1). El contrato dice en cambio que «el integrador une los dos subcomandos en un solo `main`» (§2) y no menciona `s3.go`, el evaluador ni el código de error de una `funcion` desconocida, que FORMATO §4.1 pide declarar.

Además, el contrato cita «FORMATO v1.0» en el encabezado, §9 y §10.
- Ruta: `docs/contratos/motores/S3-interfaz.md` §2, §3 (JSON de `Rate`), §9.1, §9.5, §10 (`s3-kueski-despensa-1`), contra `oraculo/fixtures/FORMATO.md` 1.1 §2.1, §3, §4.1, §6.
- Corrección sugerida:
  - (a) agregar `vista_dos_quincenas` (claves `total` e `intereses` en centavos y `tasa_pct` como `{num, den, decimales}`) y usarla con `id` `data.json:financiamiento.dos_quincenas` en `s3-kueski-despensa-1`; puede quedar `dos-quincenas` como cruce adicional;
  - (b) usar `decimales` en la salida que se compara con `data.json` o pedir al integrador que FORMATO admita una sola clave para todos los motores;
  - (c) §2 y §9.5: «S3 agrega `backend/cmd/motores/s3.go` con `func init() { register("s3", eval) }` contra la API de FORMATO §4.1; no edita `main.go`», y declarar el código de una `funcion` desconocida (p. ej. `invalid`);
  - citar FORMATO 1.1.

**M2 (S3-05, S3-03, S3-12) · La línea de tiempo no fija cómo los abonos registrados se vuelven eventos, ni qué pasa con un pago parcial que aún no vence.**
§5.2 dice que un pago repone su capital en la fecha en que queda completo y que «pagos con vencimiento ≥ `Today` sin registro se proyectan completos». Quedan dos preguntas sin respuesta:
- (1) Con registros (`PaymentRecord`), ¿hay un evento `payment` por registro (29-oct Karen 37493 sin reponer; 31-oct Axel 59882 que repone 77363) o uno solo por pago, fechado al completarse? ¿Con qué `Amount`, `Seq` e `InterestNotRestored`? FORMATO exige listas completas en `esperado`, así que el esperado de `s3-reposicion-separada` («repone 77363 el 31-oct, nada el 29 ni el 30») no basta: motor, oráculo y tabla a mano tienen que adivinar la misma lista.
- (2) Un pago con un abono parcial y vencimiento ≥ `Today` no es «sin registro» ni está completo, así que según el texto nunca se proyecta ni repone. Eso subestima el disponible de ahí en adelante, y con él el colchón y el objetivo 3 del comparador.
- Ruta: `S3-interfaz.md` §5.2 (Eventos, Reposición) y §10 (`s3-reposicion-separada`).
- Corrección sugerida: fijar la regla de eventos (p. ej. un evento por registro con `PrincipalRestored` = 0 salvo el que completa el pago, que repone todo el capital, e `InterestNotRestored` repartido igual) y la proyección del saldo de un pago parcial no vencido (p. ej. el resto se proyecta completo en su vencimiento). Poner la lista completa de eventos en el esperado del ejemplo 1.

**M3 (S3-10, S3-12, S3-06) · `QuotePurchase` no dice cómo trata cada tipo de crédito, y la generación de §8.2 deja fuera, al pie de la letra, el cargo normal de la tarjeta.**
- (a) **Crédito de tienda:** ningún campo de `Credit` distingue «abonos sin interés» de «sin datos K18». La fixture `s3-tienda-abonos` necesita los dos. §8.2 saca los plazos de la tienda de `TermRates` y §5.5 los saca de los datos K18, así que no se sabe si un crédito de tienda sin `TermRates` ni datos da `no_rate` (§5.1 regla 4) o `missing_capture` (§5.5), ni cómo un crédito de tienda con interés arma sus abonos: la fórmula de la tarjeta es por periodo de corte, no por compra. La variante `s3-tienda-abonos-interes` pide «intereses iguales al oráculo» sin algoritmo común.
- (b) **Tarjeta en el comparador:** §8.2 genera «sin promoción» solo «si el plazo tiene tasa» y dice que «un plazo sin tasa ni promoción no se genera». El cargo normal (plazo 1) no tiene `TermRate`, así que al pie de la letra no se genera, y el caso 2 depende de que compita.
- (c) **Préstamo:** `Amortize` no recibe partes, y no se dice que `QuotePurchase` lo llame para un `term_loan` ni qué error da sin `AnnualRate`. El calendario unificado necesita la persona responsable (S3-11).
- (d) `Purchase.Source` no tiene valor para el cargo normal ni para los abonos sin interés.
- Ruta: `S3-interfaz.md` §5 (`QuotePurchase`), §5.3–§5.5, §8.2.
- Corrección sugerida: una tabla de despacho de `QuotePurchase` por `Kind` con, para cada tipo, entradas obligatorias, regla de intereses, serie, `Source` y error con precedencia (`ErrNoRate` contra `ErrMissingCapture`). Definir cómo se captura una tienda sin interés (p. ej. `TermRate` con tasa 0/1 o una bandera) y su variante con interés (o declararla `ErrUnsupported` PROVISIONAL). En §8.2: «el cargo normal de la tarjeta (plazo 1, 0.00) siempre se genera».

**M4 (S3-20, S3-12, S3-13) · Las marcas PROVISIONAL del comparador no están fijadas.**
§14 dice que «toda salida que dependa» de un bloqueo lleva `Provisional`, y todas las salidas del comparador dependen de `S3-objetivos-medida`:
- (2) serie del objetivo 2 y su desempate, y (3) ventana del colchón y del objetivo 3, en todo escenario;
- (1) objetivo 4 cuando hay canasta, y (4) créditos sin límite;
- en el caso 2, además, el calendario de la tarjeta (`S3-formula-tarjeta`, pregunta 4).

Aun así, las fixtures `s3-comparador-caso-1…-3` y `-carga` no llevan `provisional` en §10, mientras que `s3-tarjeta-minimo`, `s3-prestamo-frances` y `s3-unificado` sí. El contrato tampoco dice qué lista y en qué orden lleva cada `Scenario.Provisional`, `Result.Provisional`, `Purchase.Provisional` o `CalendarEntry.Provisional`. Como el esperado del comparador es «`valid` completo (cada escenario con su llave y sus valores)», las tres fuentes tienen que adivinar esas listas. Y según el protocolo §2, nada provisional entra en cifras sin la marca.
- Ruta: `S3-interfaz.md` §8 (`Scenario`, `Result`), §10 (tabla, filas del comparador), §14.
- Corrección sugerida: marcar las 4 fixtures del comparador con `provisional` `{bloqueo: docs/estado/bloqueos/S3-objetivos-medida.md, que: …}` (y `S3-formula-tarjeta` en el caso 2). Fijar la regla de las listas: qué bloqueo entra en cada salida, sin repetidos y en orden alfabético; `Result.Provisional` = unión de las de sus escenarios.

### Bajos
- **B1 (S3-20):** los bloqueos citan secciones y nombres de la v1.0. `S3-formula-tarjeta.md` habla de «§5.2 y §5.4», `Card.InterestForPeriod`, `Card.MinimumPayment` y `StoreCredit`, cuando hoy son §5.3/§5.5, `CardPeriodInterest` y `CardMinimumPayment`. `S3-convencion-prestamo.md` habla de «§5.3» y `Loan.Amortize`, cuando hoy es §5.4 `Amortize`. Corrección: actualizar las referencias (solicitud al dueño de los bloqueos).
- **B2 (S3-01):** §11.1 escribe `RoundHalfUp(x, 2)` con centavos y un entero, pero la firma de §3 es `RoundHalfUp(num, den *big.Int)`. Corrección: escribir «HALF_UP(x / 2)» o la llamada con `big.Int`.
- **B3 (S3-08, PROVISIONAL):** `CardPeriodInterest(st, paidCents, paidDate)` no dice qué hacen `paidCents` y `paidDate` cuando no se cumple el 0.00 (abono parcial a tiempo, total tardío), ni si el IVA se calcula sobre el interés ya redondeado. Corrección: «sin pago total a tiempo, intereses = fórmula sobre el estado completo; IVA = HALF_UP(interés redondeado × bp / 10000)», o lo que se decida.
- **B4 (S3-01, S3-09):** `-run Canon` casa también con `TestCanonHojas437`, así que la salida no muestra «85 de 85» limpio. Además, S3-09 corre `cred_run_fixtures.py` sin argumentos y el contrato no da valores por omisión a `--fixtures` y `--out`. Corrección: declarar esos valores por omisión y precisar el conteo esperado, o renombrar la prueba de hojas.
- **B5 (S3-03, lectura A):** la regla «compras con ancla < fecha ≤ `From`» deja fuera de A una compra del escenario hecha el mismo día del ancla, cuando `Today` = ancla. Corrección: aclarar que la compra del escenario siempre resta.
- **B6 (S3-13, PROVISIONAL (1)):** no se dice cómo se elige el `BasketQuote` cuando la petición trae solo `merchant_name`. Corrección: misma regla que las promociones (`StoreID` y, si falta, el nombre sin distinguir mayúsculas).
- **B7 (S3-14):** el conjunto «generados = `valid`» está declarado y justificado con la decisión C, pero S3-12 dice «genera … y descarta», y RN-11 y C1 `Scenario.score_bp` siguen diciendo «entre los escenarios generados» sin aclararlo. Corrección: que el integrador aplique la solicitud ya escrita y que Axel lo vea en la puerta 2.

## Estado de los hallazgos de la ronda 1
| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| H1 fixtures, CLI y oráculo contra FORMATO; R-05/R-06 sin fuente | **Reabierto en parte (M1)** | Carpeta, llaves `fixture/1`, `salida/1`, `cred_run_fixtures.py`, `compara.py s3` y fuente `data.json`: cerrado. `id` de R-05/R-06, `decimales` y `s3.go`/`main.go`: chocan con FORMATO 1.1 |
| H2 tipos sin declarar | Cerrado | `CreditState`, `PurchaseState`, `Input`, `Result`, `Scenario`, `SensitivityRow`, `DebtLine`, `PersonPeak`, `StatementInput`, `Reading` con etiquetas `json` (§3–§8); `DecodeInput` estricto |
| H3 TestCanon85 sin aserción en las 30 filas S2 | Cerrado | `CashFlow` (§6.1) y la tabla de §11.1; las 30 cuadran en mi script |
| H4 conjunto de la normalización | Cerrado (bajo B7) | §8.5 lo declara y lo justifica; solicitud a RN-11 y C1 |
| H5 inicio de A con `From` > ancla | Cerrado | Fórmula nueva verificada en las 16 fechas y en los 2 ejemplos |
| H6 objetivo 2 no unívoco; `over_cap` | Cerrado | `FortnightOf`, incluye totales 0; desempate con `OverCap` en §8.6; `TestObjetivoParejo` |
| H7 doble conteo del corte | Cerrado | §6: el corte sustituye, partes y llave; PROVISIONAL |
| H8 503 y app_setting en un motor puro | Cerrado | `TimeoutFromSetting`, `ErrorResponses`/`ResponseFor`, `TestTiempoClaveSembrada` (lee `00003_puerta1b.sql`, valor 5), parte de WF-2 declarada con solicitud a S3.md |

## (c) Preguntas que todavía tendría que adivinar el implementador
1. Eventos de pago con abonos registrados y proyección de un pago parcial no vencido (M2).
2. Tienda «sin interés» frente a «sin datos»; precedencia `no_rate`/`missing_capture`; abonos con interés (M3).
3. El cargo normal de la tarjeta en la generación del comparador; el préstamo dentro de `QuotePurchase` (partes, error); `Source` del cargo normal y de los abonos (M3).
4. Contenido y orden de las listas `Provisional`; `provisional` en las fixtures del comparador (M4).
5. `id`/vista de R-05 y R-06, `decimales`, evaluador en `s3.go` y código de una `funcion` desconocida (M1).
6. Abono parcial o tardío en `CardPeriodInterest` e IVA (B3); `BasketQuote` solo con nombre de comercio (B6); valores por omisión de `cred_run_fixtures.py` (B4).
