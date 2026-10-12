# Revisión del contrato S3 · ronda 3 de 3 (agente revisor-contrato-S3-r3)

11-oct-2026 · WF-1 · Objeto: `docs/contratos/motores/S3-interfaz.md` v1.2 (115,607 bytes) · Vara: `docs/criterios/S3.md` 1.2, decisiones §C y §K; contrastes de consistencia con `oraculo/fixtures/FORMATO.md` (hoy 1.2), esquema 00001–00003, C1 por extractos y bloqueos abiertos. No corregí nada. Todas las cifras las recalculé con scripts propios (`fractions.Fraction`) en `datos/tmp/rev_s3_r3/` (ya borrados).

## Veredicto
**NO APROBADO.** 0 hallazgos altos, **4 medios**, 9 bajos. Los 4 hallazgos medios de la r2 (M1–M4) quedan **cerrados**. Los 4 medios nuevos son una pregunta que el implementador tendría que adivinar (N2) y tres inconsistencias con otros artefactos que el contrato no marca PROVISIONAL ni lleva a bloqueo (N1, N3 y N4). Ninguna cifra canon cambia con sus correcciones.

## Verificaciones propias
| Verificación | Resultado |
|---|---|
| (a) Hojas numéricas de `model.json` (script propio que recorre el árbol y casa cada hoja con los 42 patrones de §12 convertidos a expresiones regulares) | 596 en total; rama S3 **437** = `parametros` 9, `despensas` 144, `credito` 143, `sensibilidad` 45, `disposicion_por_despensa` 68, `interes_por_monto` 28; rama S2 159. 42 renglones, conteo por renglón = el declarado en los 42; **0 sin dueño, 0 con dos dueños**. Por fixture: `parametros` se reparte sin traslape entre `s3-kueski-despensa-1` (7) y `s3-series-a-b` (2). sha256 de `model.json` = `1a57066689dd…f1bfd4` (intacto) |
| 64 pagos canon (fecha, monto, capital, parte de Axel) con `sched`, `SplitTrunc`, `SplitHalfUp` y partes de §4 | 64 de 64 iguales a `model.json` |
| R-05 y R-06 con las reglas de §4 (d1 a 2 q, tasa 53450/309450) | total 362900, intereses 53450, tasa HALF_UP a 2 decimales = 17.27; Axel 32870 / 223170 / 111585; Karen 20580 / 139730 / 69865 = las 9 hojas de `data.json:financiamiento.dos_quincenas` (no tiene más hojas) |
| Línea de tiempo de §5.2, implementación propia | ejemplo 1 de §10.3: los 4 eventos, series [1000000, 1077363, 1154726] / [690550, 767913, 845276] y mínimos idénticos al `esperado`; canon: 14 + 14 valores = R-65/R-66, mínimos 809450 / 500000; desde 2026-12-16: 909450 / 600000; desde 2027-02-10: 809450 / 500000 |
| Desviación estándar de §8.4 | 0/100/200 → 81650; 20,000 series aleatorias: fórmula del contrato = HALF_UP(1000 σ) exacto en todas |
| Hallazgo 1 y 4 del contrato | pago de referencia desde total redondeado: 32 y 125869 contra 31 y 125868 (confirmado); binary64(8005000/309450) = 25.868476328970754 = `tasa_pct`, float(tasa)×100 = 25.86847632897075 (≠) |
| Cifras de fixtures | MSI 3 y tienda 3 de 100000 → 33333, 33333, 33334; 2028-02-10 → 02-29, 03-15, 03-30, 04-15; efectivo P1…P6 = `periodos[k].efectivo.total`; promoción 0 % de 309450: montos 77362 ×3 + 77364, capital 77363 ×3 + 77361, **intereses −1, −1, −1, +3** |
| C1 y esquema | orden de `missing_constraints` en `creditos-b.yaml` = §8.1; `comparator.timeout_seconds` sembrado en 5 (00003 l. 86); `optimizer.weights` 1,1,1,1 (00001 l. 1396); `scheduled_payment.interest_cents CHECK (>= 0)` (00001 l. 808) |

## Tabla de criterios
| Criterio | Estado | Evidencia propia |
|---|---|---|
| S3-01 (alta) | Cubierto | §11 `TestCanon85`, §11.1 y §10; 64 pagos y R-05/R-06 recalculados iguales |
| S3-02 (alta) | Cubierto | conteo propio 437, 0 sin dueño; `compara.py s3` con su línea de 437 (§10.2) |
| S3-03 (alta) | Cubierto | §5.2; series, mínimos e inicios desde otra fecha recalculados |
| S3-04 (alta) | Cubierto | §5.1 y `s3-feb-2027`; 16 fechas y 2028-02-29 |
| S3-05 (alta) | Cubierto | §5.2 (un evento por abono, el que completa repone, resto proyectado); ejemplo 1 reproducido |
| S3-06 (alta) | **Parcial (N2)** | §5.1 y §5.8 correctos en intereses y `ErrNoRate`; los decimales de la tasa cotizada no tienen entrada |
| S3-07 (alta) | **Parcial (N1)** | 4 fronteras y 0.00 correctos; la serie de la promoción en la línea choca con el esquema sin bloqueo |
| S3-08 (alta) | **Parcial (N3)** | §5.3 y §5.8; el dato «si incluye IVA» de K18 no se trata y su bloqueo abierto no se cita |
| S3-09 (alta) | Cubierto (PROVISIONAL `S3-convencion-prestamo`) | §5.4, §5.8; fixture con ajuste de centavos |
| S3-10 (media) | **Parcial (N1, N3)** | §5.5 y §5.8 cierran r2-M3; capital = monto sin decisión e IVA sin bloqueo |
| S3-11 (alta) | Cubierto | §6 (corte que sustituye, PROVISIONAL), persona responsable, orden y llave sin duplicados |
| S3-12 (alta) | Cubierto | §8.1–§8.3, §8.7, §11; C1 y 00003 contrastados |
| S3-13 (alta) | Cubierto en el contrato; ver N4 | §8.4 (fórmula verificada), §8.6, §10.1; la autoría de las fixtures del comparador queda abierta |
| S3-14 (alta) | Cubierto | §8.5: conjunto = `valid` (decisión C), empate = 1, ejemplo 100/200/400 |
| S3-15 (alta) | Cubierto | §1.5, §2 |
| S3-16 (alta) | Cubierto; ver N4 y B6 | §10.2; el autor del oráculo no se declara |
| S3-17 (alta) | Cubierto | §1.1–§1.3 |
| S3-18 (alta) | Cubierto | §11 (cobertura ≥ 85 % por paquete, gofmt, vet) |
| S3-19 (media) | Cubierto | §4 |
| S3-20 (media) | **Parcial (N1, N3)** | §14 y §14.1 completos para los 4 bloqueos S3; faltan dos dependencias sin bloqueo o sin cita |

## Estado de los hallazgos de la ronda 2
| Hallazgo r2 | Estado | Evidencia |
|---|---|---|
| M1 (FORMATO 1.1: R-05/R-06, `decimales`, `s3.go`) | **Cerrado** | (a) §9.3 `vista_dos_quincenas` y §10 con `id` `data.json:financiamiento.dos_quincenas`, `fuente` `data.json` y `tasa_pct` con `decimales`; mis 9 hojas coinciden. (b) `json:"decimales"` en §3 y §9.1. (c) §2 y §9.5: solo `backend/cmd/motores/s3.go` con `register("s3", evalS3)` y la firma `Evaluator`; `funcion` desconocida → `invalid`. Residuo nuevo: FORMATO ya es 1.2 (N4) |
| M2 (abonos → eventos) | **Cerrado** | §5.2 fija un evento por abono, el que completa lleva K e I, sobrepago 0/0, resto proyectado de un pago parcial no vencido y fecha de reposición; ejemplo 1 con `esperado` completo, reproducido por mi implementación |
| M3 (`QuotePurchase` por tipo) | **Cerrado** | §5.8: tasa 0 = abonos sin interés, precedencia sin `ErrNoRate` en tienda, tienda con interés `unsupported` PROVISIONAL; §8.2: el cargo normal siempre se genera; `term_loan` llama `Amortize` y aplica partes; sin `AnnualRate` → `ErrMissingCapture`; `Source` por tipo. Residuos nuevos: N1 y N2 |
| M4 (listas `provisional`) | **Cerrado** | §14.1 fija contenido y orden alfabético por bytes de las listas de `Purchase`, `CalendarEntry`, consultas de tarjeta, `Scenario` y `Result`; §10 marca `s3-comparador-caso-1…-3` y `-carga` |

## Hallazgos nuevos
### Medios
**N1 (media · S3-07, S3-10, S3-20, K14) · Serie con intereses 0: interés negativo en la línea y «capital = monto» sin decisión.** §5.1 aplica RN-07 (montos truncados, capital HALF_UP) también a la promoción al 0 %: con 309450 a 4 q salen montos 77362 ×3 + 77364 y capital 77363 ×3 + 77361, es decir, **intereses −1, −1, −1 y +3** (lo recalculé). `scheduled_payment.interest_cents CHECK (interest_cents >= 0)` (00001 l. 808) rechaza esa serie al guardarla (C-08), y es el caso más común de Kueski (RN-10: promociones casi siempre 0 % a 4 q). A la vez, §5.3 (meses sin intereses) y §5.5 (abonos sin interés) usan «capital = monto» marcado FIRME, que no está en K14 ni en RN-07 (RN-07 FIRME: capital con base HALF_UP en toda serie) y que con otros importes difiere de RN-07. El contrato lo deja como «Hallazgo 5 … no se corrige aquí» con una propuesta, sin bloqueo y sin marca PROVISIONAL: dos reglas distintas para el mismo caso y un choque con el esquema que nadie decide.
- Ruta: `docs/contratos/motores/S3-interfaz.md` §5.1 (Serie), §5.3 (Meses sin intereses), §5.5, §14 (Hallazgo 5); `backend/db/migrations/00001_esquema.sql` l. 808; `docs/contratos/reglas-negocio.md` RN-07 y RN-10.
- Corrección sugerida: abrir `docs/estado/bloqueos/S3-serie-sin-intereses.md` (opciones: (a) con intereses 0, capital = monto en todos los tipos; (b) RN-07 literal y el CHECK admite negativos por pago; (c) RN-07 literal y el ajuste va al primer pago). Marcar PROVISIONAL citando ese bloqueo la serie de capital de la promoción en la línea (§5.1), de MSI (§5.3) y de la tienda sin interés (§5.5); sumarlo a §14 y al orden de §14.1. Los montos (77362 ×3 + 77364; 33333, 33333, 33334) y las cifras canon no cambian.

**N2 (media · S3-06, K12) · Los decimales de la tasa cotizada no tienen entrada.** `QuotedRate(interest, principal, decimals)` pide los decimales, pero `PurchaseInput` no trae ningún campo para ellos (el esquema sí: `credit_purchase.rate_decimals`, 00002 l. 123), y `Rate.Display()` no dice qué hace con `Decimals` = −1 («sin dato de captura»). Aun así §5.1 dice «25.8685 % con 4 decimales» para d1 y la fixture `s3-tasa-plazo` espera «25.8685 %» para «d1 con cotización». Motor, oráculo y fixture tienen que adivinar de dónde salen los 4 decimales (¿del `TermRate` del mismo plazo?, ¿un valor fijo?) y qué muestra una tasa sin decimales.
- Ruta: `docs/contratos/motores/S3-interfaz.md` §3 (`QuotedRate`, `Display`), §5 (`PurchaseInput`), §5.1 (1), §10 (semilla d1, `s3-tasa-plazo`).
- Corrección sugerida: agregar `RateDecimals *int json:"rate_decimals"` a `PurchaseInput` (espejo de `credit_purchase.rate_decimals`), obligatorio con `QuotedInterest` para mostrar la tasa, poner 4 en d1 de la semilla y fijar `Display()` con −1 (por ejemplo, la fracción sin redondear o error), con el orden de precedencia si también hay `TermRate` del plazo.

**N3 (media · S3-08, S3-10, S3-20, K18) · Bloqueo abierto `U4-iva-credito` sin citar.** K18 pide cinco datos por producto, entre ellos «si incluye IVA», y «sin datos muestra falta capturar». El contrato lee `VATOnInterestBP` nil como «la tasa ya trae IVA» y su `ErrMissingCapture` solo revisa 4 datos (§5.3, §5.8), así que un producto sin ese dato capturado calcula como si incluyera IVA. El bloqueo `docs/estado/bloqueos/U4-iva-credito.md` (abierto, escrito antes de la v1.2 del contrato) pregunta justo eso y su opción 2 «cambia el significado de NULL en esquema, C1 y motor S3»; el contrato lo trata como FIRME y no lo nombra en §5.3, §5.8, §14 ni §14.1.
- Ruta: `docs/contratos/motores/S3-interfaz.md` §5 (`VATOnInterestBP`), §5.3, §5.8 (renglones `bank_card` y `store_credit`), §14, §14.1.
- Corrección sugerida: marcar PROVISIONAL citando `U4-iva-credito` el tratamiento del IVA y la lista de datos de `ErrMissingCapture`; sumar el bloqueo al orden de §14.1 (`S3-convencion-prestamo` < `S3-formula-tarjeta` < `S3-objetivos-medida` < `S3-origen-cargo-normal` < `U4-iva-credito`) y a las listas de `interes_tarjeta`, `minimo_tarjeta` y de las fixtures `s3-tarjeta-minimo` y `s3-tienda-abonos-interes` («otros bloqueos»).

**N4 (media · S3-13, S3-16; FORMATO 1.2 §1 y §4) · El contrato cita FORMATO 1.1, pero rige la 1.2, y no declara la autoría de fixtures ni del oráculo.** `FORMATO.md` 1.2 (11-oct-2026, anterior a la v1.2 del contrato) exige que cada interfaz declare quién escribe cada archivo de `oraculo/fixtures/<motor>/` (casos y `_datos/`) y qué hace el otro agente mientras falten; «un archivo sin autor declarado es un hueco de la interfaz» (S2 lo hace en `S2-interfaz.md` §5.0). El contrato solo dice que el `esperado` lo escribe el agente de regresión (§10) y que `S3-escenarios.md` lo escribe el optimizador (§2). No asigna los archivos de casos ni `_datos/kueski-semilla.json` ni `_datos/comparador-caso-N.json`, no asigna `oraculo/cred_lib.py`, `cred_escenarios.py` y `cred_run_fixtures.py` (§2 los lista sin autor), no dice qué hacen el núcleo y el optimizador mientras falten las fixtures (`TestCanonHojas437`, `TestOraculoFixtures` y `TestOptimos` las leen) y no fija la secuencia «tabla a mano del optimizador → fixture del comparador» que pide S3-13 («fijados antes de correr»). Tampoco declara la regla 5 de §4 (qué claves de persona trae cada objeto `<persona>` de las vistas).
- Ruta: `docs/contratos/motores/S3-interfaz.md` encabezado, §2, §9, §10 y §10.1 (contra `oraculo/fixtures/FORMATO.md` 1.2 §1 y §4 regla 5).
- Corrección sugerida: citar FORMATO 1.2 y agregar un §10.0 «Autoría» con un renglón por archivo (casos, `_datos/`, `cred_*.py`, `S3-escenarios.md`) con su agente (núcleo, optimizador o regresión), qué hace cada uno mientras falten (p. ej., el núcleo prueba con datos en línea en `_test.go` sin cifras canon escritas) y la secuencia del comparador; y en §9.3 decir que `<persona>` toma las claves de todas las partes de la compra (titular y externas).

### Bajos
- **B1 (S3-08):** §5.3 dice que `CardPeriodInterest` sin cualquiera de los 4 datos da `ErrMissingCapture`, pero `s3-tarjeta-sin-intereses` espera 0.00 para la tarjeta sin datos pagada a tiempo. Escribir que «pago total a tiempo» (`paid_cents` ≥ `balance_cents` y `paid_date` ≤ `due_date`) se revisa antes y da 0.00 con `[]`.
- **B2:** `S3-formula-tarjeta.md` (Impacto: «§5.2 y §5.4», `Card.InterestForPeriod`, `StoreCredit`, «la variante con interés de `s3-tienda-abonos`») y `S3-convencion-prestamo.md` (`Loan.Amortize`, «§5.3») citan secciones y nombres de la v1.0; hoy son §5.3/§5.5, `CardPeriodInterest`, `Amortize` (§5.4) y `s3-tienda-abonos-interes`.
- **B3 (S3-02):** las vistas no dicen si emiten las hojas de texto `despensas[*].fecha`, `disposicion_por_despensa[*].fecha` y `parametros.horizonte`; si el script del `esperado` las copia, salen `__falta__`. Listarlas como emitidas o excluidas (hoy solo se excluyen `descripcion`, `nota` y `cubre`).
- **B4 (S3-12):** casos sin regla en §8.2: un `term_loan` con `Limit` en `Credits` (no hay plazos para préstamo), petición sin tienda y sin `baskets`, `BasketQuote` cuando la petición trae solo `merchant_name` (r2-B6 sigue abierto) y la `Ref` de la compra del escenario (no debe chocar con una existente).
- **B5 (S3-12):** no se fija qué imprime `python3 -I oraculo/cred_escenarios.py` sin argumentos, y S3-12 espera de él «escenarios aceptados con violación = 0».
- **B6 (S3-16):** el contrato no menciona `oraculo/verify3.py` línea 89; `S2-interfaz.md` §0.5 lo asigna al agente Python de S2. Citarlo en §13 (S3-16) para que el revisor de motores sepa dónde verificarlo.
- **B7 (S3-05):** un `PaymentRecord` con `seq` fuera de 1…`term_count` no tiene regla (¿`ErrInvalid`?).
- **B8 (S3-11):** §5.8 dice que un crédito con una compra `unsupported` «no entra al calendario unificado (`CreditState.Quote` devuelve el error)»: no queda claro si `UnifiedCalendar` omite ese crédito o falla entero.
- **B9 (S3-13):** `value_per_peso_milli` redondea a milésimas enteras por peso (el ejemplo 2 da 13, 10 y 15), lo que produce empates frecuentes; y en el caso 2 el cargo normal de 0.00, que siempre se genera, empata con la promoción en el objetivo 1 y con la mejor tienda en el objetivo 4. Si el escenario de mayor puntaje entre los de interés 0 está en la mejor tienda, los óptimos 1 y 4 coinciden; el optimizador debe diseñar el caso para evitarlo (nota para `S3-escenarios.md`).

## Suficiencia para implementar (c)
Con N1–N4 corregidos, núcleo, comparador y oráculo se implementan sin preguntas: tipos con llaves JSON, despacho por tipo, línea de tiempo con abonos, calendario, desviación estándar exacta, normalización, desempates, tiempo máximo y listas `provisional` están fijados y los recalculé. Preguntas que hoy quedarían: de dónde salen los decimales de la tasa cotizada (N2); qué serie de capital usar con intereses 0 y qué hacer con el CHECK (N1); si falta el dato de IVA, «falta capturar» o calcular (N3); quién escribe cada fixture, `_datos` y el oráculo, y en qué orden (N4); y los bajos B1, B4, B7 y B8.

## Coherencia con K12–K18 y K31 (d)
K12: tasa como fracción con decimales (falta la entrada de la tasa cotizada, N2). K13: tope opcional por persona, solo marca y desempata. K14: dos modos aplicados; «capital = monto» de §5.3/§5.5 se aparta sin decisión (N1). K15: §5.2 completo. K16: globales y por comparación, sin valores por defecto, «faltan restricciones», colchón en la línea usada, vivas y deuda en el crédito usado. K17: mín-máx entre generados (`valid`, decisión C), desviación estándar y pico, unidades base por peso (PROVISIONAL `S3-objetivos-medida`). K18: datos por producto y `ErrMissingCapture`, salvo el IVA (N3). K31: `optimizer.date_window_days` propuesto en 7 días, sin valor sembrado. K43/K44: 5 s desde `app_setting`, configurable, 503 `unavailable` como dato.

## Los 3 casos del comparador (e)
§10.1 exige en cada caso los 4 valores de restricción en la fixture, ≥ 2 créditos con límite, ≥ 2 plazos con tasa, una promoción 0 % acotada, ≥ 2 tiendas, compras existentes desiguales, un descartado por restricción y uno sobre el tope, 4 óptimos distintos verificados a mano antes de correr, y la intención de cada óptimo; el ejemplo 2 trae las restricciones explícitas. Cumple en el contrato; las cifras las fija `S3-escenarios.md` (ver B9).
