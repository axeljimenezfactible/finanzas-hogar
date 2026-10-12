# Revisión del contrato S3 · ronda 4 (agente revisor-contrato-S3-r4)

11-oct-2026 · WF-1 · Artefacto: `docs/contratos/motores/S3-interfaz.md` v1.3 (916 líneas, 129,503 bytes) · Se evalúa solo contra `docs/criterios/S3.md` 1.2 y `docs/plan/decisiones.md` §C y §K. No corregí nada; todos los números de abajo salen de mis propios comandos (script temporal `datos/tmp/rev-s3-r4/verifica.py`, borrado al terminar).

## Veredicto
**NO APROBADO.** 0 hallazgos altos, **1 medio** (P1, residuo de r3-N4) y 13 bajos (4 nuevos y 9 de la ronda 3 que siguen abiertos). De los 4 medios de la ronda 3: N1, N2 y N3 quedan **cerrados**; N4 queda cerrado en lo que pedía (autoría, secuencia del comparador, FORMATO 1.2, regla 5), salvo el residuo P1.

## Verificaciones propias
| Qué | Comando o método | Resultado |
|---|---|---|
| `model.json` intacto | `sha256sum oraculo/model.json` | `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` (igual al de S3-16) |
| (a) Conteo de hojas | recorrido propio de `model.json` (numéricas, sin booleanos) | parametros 9, despensas 144, credito 143, sensibilidad 45, disposicion_por_despensa 68, interes_por_monto 28 = **437**; S2 = 159; total 596 |
| (a) Cruce con la tabla §12 | los 42 renglones convertidos a expresiones regulares | **0 sin dueño, 0 con dos dueños**, los 42 conteos declarados coinciden (suma 437) |
| Reglas del contrato contra el canon | reimplementación propia con `fractions.Fraction` de §3, §4, §5.1, §5.2 y §7 (sched, SplitTrunc, SplitHalfUp, partes, composición, línea de tiempo, sensibilidad, cotización por monto) | **437 de 437 hojas iguales**, además fechas de los 16 pagos y fecha y tipo de los 17 eventos |
| R-05 y R-06 (`data.json`) | d1 a 2 q con la tasa 53450/309450 | 362900, 53450, axel 32870/223170/111585, karen 20580/139730/69865 y 17.27 con `decimales` 2: iguales a §10 y a `data.json` |
| K12 · `Display()` | HALF_UP a los decimales capturados | «25.8685» (4) y «17.27» (2) |
| Hallazgos 1 y 4 del contrato | pago de referencia desde el total redondeado; binary64(tasa) × 100 | el primero falla justo en `interes_por_monto[0]` y `[5]`; el segundo no da `tasa_pct` (sí lo da binary64(8,005,000/309,450)) |
| Efectivo, horizonte y deuda (§6.1, §7) | ventanas `periodos[k]` de `model.json` | 97375, 194750, 223243, 271636, 251738, 271636; horizonte 1310378 con 39800 de contado; deuda viva 629344 = 125870 + 503474 |
| Ejemplos de §3 y del bloqueo | aritmética entera | `SplitTrunc(389500,4)` 97375×4; `SplitTrunc(503474,4)` 125868×3+125870; `SplitHalfUp(309450,4)` 77363×3+77361; `SplitHalfUp(149972,4)` 37493×4; promoción 0 %: intereses por pago −1, −1, −1, +3; opción c: +1, −1, −1, +1 y −1, 0, +1; 100000 en 3 iguales con ambas reglas; KT 149972; 4,000.00 → 103474 |
| Calendario | `sched` propio | 10-feb-2027 → 28-feb, 15-mar, 30-mar, 15-abr; 10-feb-2028 → 29-feb, 15-mar, 30-mar, 15-abr |
| Objetivo 2 (§8.4) | ⌊(⌊isqrt(4N)/n⌋+1)/2⌋ contra `Decimal` de 60 dígitos | 0/100/200 → 81650; **0 discrepancias en 20,000 series aleatorias** de 1 a 14 quincenas |
| S3-14 | (máx − v)/(máx − mín) | 100/200/400 → 1, 2/3, 0 |
| Ejemplo 1 de §10.3 | K15 a mano | M 97375, K 77363, I 20012; A 1000000, 1077363, 1154726; B 690550, 767913, 845276; pago 2 el 15-nov: igual al `esperado` |
| Esquema | `sed`/`grep` en las migraciones | 00001 l. 808 `CHECK (interest_cents >= 0)`; 00002 `credit_purchase.rate_decimals` con `CHECK (rate_decimals IS NULL OR (rate_decimals >= 0 AND rate_num IS NOT NULL))`; 00003 l. 86 `comparator.timeout_seconds` = 5 y l. 54–55 `optimizer.*` sin valor; 00001 l. 448 `vat_on_interest_bp` NULL = incluye IVA; l. 441 `availability_reading` NOT NULL DEFAULT 'B' |
| C1 por extractos | `grep` | `rate_pct` «ausente si falta rate_decimals» (`creditos-a.yaml` l. 739, `creditos-b.yaml` l. 872); `missing_constraints` con las 4 claves (`creditos-b.yaml` l. 1446–1455) |
| Bloqueos | lectura | `S3-serie-sin-intereses.md` con las 5 secciones del protocolo §2; `U4-iva-credito.md` existe y el contrato lo cita en encabezado, §5, §5.3–§5.5, §5.8, §14 y §14.1 |

## Tabla de criterios
| Criterio | Estado | Evidencia propia |
|---|---|---|
| S3-01 (alta) | Cumple, con P1 | §4, §5.1, §6.1, §7, §9.3, §11 y §11.1; reproduje R-05…R-16, R-45…R-55 y las 9 hojas de `data.json`. P1: §11.1 y §10.0 se contradicen sobre de dónde salen las entradas de `TestCanon85` |
| S3-02 (alta) | Cumple, con P1 | 437 hojas contadas y cruzadas (0 sin dueño); 437 de 437 reproducidas con las reglas. P1: `seedFromModel` no arma las entradas de `sensibilidad` ni de `interes_por_monto` (73 hojas) |
| S3-03 (alta) | Cumple | §5.2: 14 + 14 valores y mínimos 809450 / 500000 reproducidos; `""` → B y `pending_confirmation` |
| S3-04 (alta) | Cumple | §5.1 y `s3-feb-2027`: 16 fechas y 28-feb reproducidos |
| S3-05 (alta) | Cumple | §5.2 (K15): ejemplo 1 reproducido |
| S3-06 (alta) | Cumple | §3, §5, §5.1, §5.8: `RateDecimals`, `Display()` con −1, `ErrNoRate` para 6–12 q, 103474 |
| S3-07 (alta) | Cumple | §5.1: 0.00 e intereses FIRMES; serie de capital PROVISIONAL con su bloqueo (baja Q1) |
| S3-08 (alta) | Cumple, con B1 abierto | §5.3, §5.8; lista de 4 datos PROVISIONAL (`U4-iva-credito`) |
| S3-09 (alta) | Cumple | §5.4 (propiedades FIRMES, tasa del periodo PROVISIONAL) |
| S3-10 (media) | Cumple | §5.5, §5.8 (tasa 0 = abonos sin interés; datos K18 → `unsupported` PROVISIONAL) |
| S3-11 (alta) | Cumple, con B8 abierto | §6 |
| S3-12 (alta) | Cumple, con B4 y B5 abiertos | §8.1–§8.3, §8.7, §10.0, §11; 5 s desde `app_setting`, 503 como dato |
| S3-13 (alta) | Cumple, con B9 y Q2 | §8.4, §8.6, §10.0 (secuencia), §10.1 |
| S3-14 (alta) | Cumple | §8.5 (empate = 1, conjunto = `valid`), ejemplo reproducido |
| S3-15 (alta) | Cumple | §1.5, §2 |
| S3-16 (alta) | Cumple (B6) | §10.0, §10.2; `model.json` intacto; `grep reporte-compra-latex oraculo/*.py` = 0 líneas hoy |
| S3-17 (alta) | Cumple | §1.1–§1.3 |
| S3-18 (alta) | Cumple | §11 |
| S3-19 (media) | Cumple | §4: 308.22 / 492.28 y 374.93 / 598.82 reproducidos |
| S3-20 (media) | Cumple | §14 y §14.1; 5 bloqueos propios más `U4-iva-credito` |

## Estado de los hallazgos de la ronda 3
| Hallazgo | Estado | Evidencia |
|---|---|---|
| N1 (serie con intereses 0, CHECK, «capital = monto») | **Cerrado** | bloqueo `S3-serie-sin-intereses.md` (formato del protocolo §2, opciones a/b/c verificadas arriba); PROVISIONAL en §5.1, §5.3, §5.5, §5.8, `s3-promocion-0`, `s3-tarjeta-sin-intereses`, `s3-tienda-abonos`, fixtures del comparador, §14 y §14.1. Residuo menor: Q1 |
| N2 (decimales de la tasa cotizada) | **Cerrado** | `PurchaseInput.RateDecimals` (espejo de 00002), precedencia en §5.1, `Display()` con −1 → `""`, d1 con `rate_decimals` 4 en la semilla y el ejemplo 1, variante null en `s3-tasa-plazo`, `TestTasaPorCompra` |
| N3 (IVA y «falta capturar») | **Cerrado** | `U4-iva-credito` citado y la lectura de NULL y la lista de 4 datos marcadas PROVISIONAL (§5, §5.3–§5.5, §5.8, §14.1, fixtures `s3-tarjeta-minimo` y `s3-tienda-abonos-interes`) |
| N4 (FORMATO 1.2, autoría, mientras falten) | **Cerrado salvo P1** | §10.0 con autor por archivo, secuencia y trabajo de cada agente; FORMATO 1.2 en las citas vigentes; regla 5 en §9.3. Queda la contradicción con §11.1 y las entradas de sensibilidad (P1) |

## Hallazgos nuevos
### Medios
- **P1 (S3-01, S3-02; residuo de r3-N4; FORMATO 1.2 §1).** §10.0 (l. 602) dice que el núcleo, sin archivos de otro autor, corre `TestCanon85` y `TestCanonHojas437` con `seedFromModel`, que solo arma el `Input` (créditos y contado). Pero §11.1 (l. 797) dice que en `TestCanon85` «los valores S3 salen de las funciones con `_datos/kueski-semilla.json` y, en sensibilidad, con el insumo de `s3-sensibilidad`», dos archivos del agente de regresión; y §11 (l. 768) arma «las 6 ramas … sobre el `Input` de `seedFromModel`», cuando `vista_sensibilidad` recibe un `SensitivityInput` (montos, externa, capital actual del titular, `accrual_months`, `cash_monthly_cents`) y `vista_interes_por_monto` recibe montos, tasa y plazo, que no están en `Input`. Como las pruebas no pueden llevar cifras canon escritas, el implementador tiene que adivinar de qué hojas salen (73 de las 437 hojas y R-37…R-43) o hacer que `TestCanon85` falle hasta que exista la fixture ajena.
  - Corrección sugerida: en §10.0, agregar a la receta de `seedFromModel` (o a un segundo ayudante `sensitivityFromModel`) las entradas exactas: montos = `sensibilidad[*].N`, tasa = la de 4 q de la semilla, n = 4, externa = `parametros.karen_fija.capital` (karen), capital actual del titular = `despensas[0].axel.capital`, `accrual_months` = 2 (de esta interfaz), `cash_monthly_cents` = HALF_UP(`parametros.papel` / 2); montos de `interes_por_monto` = `interes_por_monto[*].monto`. En §11.1, cambiar la frase de l. 797 por «con el `Input` de `seedFromModel` y el insumo de sensibilidad de §10.0; si existen `_datos/kueski-semilla.json` y `s3-sensibilidad`, exige además igualdad, como `TestCanonHojas437`».

### Bajos (nuevos)
- **Q1 (S3-07, S3-20, K14).** Los intereses por pago negativos no pasan solo con intereses 0: con RN-07 salen también con intereses totales positivos de unos centavos, cuando 2r ≥ n y I < n − r (r = capital mod n). Ejemplos verificados: 309450 a 4 q con I = 1 → −1, −1, −1, +4; 309450 a 12 q con I = 1…4. Hallé 14 casos con I entre 1 y 12 centavos en 16 combinaciones. §5.1, §14.1 y el bloqueo solo marcan «`Interest` = 0», así que esas series salen FIRMES y el CHECK de 00001 l. 808 las rechaza. Es casi imposible en la práctica (intereses de centavos). Corrección: definir el disparador como «RN-07 da algún interés por pago < 0» en §5.1, §14.1 y `S3-serie-sin-intereses.md`.
- **Q2 (S3-12, S3-13).** §10.1 pide en **cada** caso «uno marcado sobre el tope», pero solo el caso 3 lleva «tope por persona capturado», y el ejemplo 2 (caso 1) trae `caps: []`. Corrección: decir que la marca sobre el tope se exige solo en el caso 3, y que los casos 1 y 2, sin tope, cubren el «sin tope capturado» de S3-13. Otra opción es pedir topes en los tres.
- **Q3 (S3-12, caso 2).** `TimelineInput.State` dice «con sus compras, pagos registrados y cortes», pero ninguna regla de §5.2 usa los cortes. Para la tarjeta del caso 2 (colchón y objetivo 3), falta escribir que el disponible se mueve con los pagos programados y los abonos, y que los cortes no lo cambian (o lo contrario).
- **Q4 (S3-20, trazabilidad).** El reporte `docs/estado/wf1/contrato-S3.md` sigue describiendo la v1.0 (46.6 KB, `cred_compara.py`, `$model`, `cmd/finanzas`, 19 fixtures, 3 bloqueos). La v1.3 mide 129.5 KB, y su estado real solo está en `correccion-contrato-S3-r1..r3.md`. Además, `S3-formula-tarjeta.md` («Cómo seguí») sigue diciendo que los meses sin intereses y los abonos sin interés «quedan FIRMES», aunque su serie de capital ahora es PROVISIONAL (`S3-serie-sin-intereses`); esto amplía B2. Corrección: nota al inicio de `contrato-S3.md` que remita a la v1.3 y a los reportes de corrección, y actualizar esa frase del bloqueo.

### Bajos de la ronda 3 que siguen abiertos (sin cambios en la v1.3)
- **B1 (S3-08):** §5.3 no dice que en `CardPeriodInterest` el «pago total a tiempo» se revisa antes que los datos K18. S3-08 y `s3-tarjeta-sin-intereses` lo resuelven (0.00 sin datos), pero el texto se contradice.
- **B2:** `S3-formula-tarjeta.md` y `S3-convencion-prestamo.md` citan secciones y nombres de la v1.0.
- **B3 (S3-02):** las vistas no dicen si emiten `despensas[*].fecha`, `disposicion_por_despensa[*].fecha` ni `parametros.horizonte`.
- **B4 (S3-12):** §8.2 no tiene regla para un `term_loan` con `Limit`, para una petición sin tienda ni `baskets`, para `BasketQuote` cuando solo viene `merchant_name`, ni para la `Ref` de la compra del escenario.
- **B5 (S3-12):** no está fijada la salida de `cred_escenarios.py` sin argumentos cuando las fixtures existen («escenarios aceptados con violación = 0»).
- **B6 (S3-16):** §13 no cita `oraculo/verify3.py` l. 89. Hoy esa línea ya no escribe fuera del repositorio (lo comprobé con `grep`).
- **B7 (S3-05):** no hay regla para un `PaymentRecord` con `seq` fuera de 1…`term_count`.
- **B8 (S3-11):** no se sabe si `UnifiedCalendar` omite un crédito con una compra `unsupported` o falla entero.
- **B9 (S3-13):** `value_per_peso_milli` en milésimas enteras empata seguido, y en el caso 2 el cargo normal empata con la promoción. Es una nota de diseño para `S3-escenarios.md`.

## Suficiencia para implementar (c)
Con lo que hay, núcleo, comparador y oráculo se implementan sin preguntas en tipos, despacho por tipo, línea de tiempo con abonos, calendario, desviación estándar exacta, normalización, desempates, tiempo máximo, listas `provisional`, autoría y secuencia de fixtures: los recalculé. Preguntas que quedarían: de dónde salen el `SensitivityInput` y los montos de `interes_por_monto` sin fixtures, y si `TestCanon85` lee archivos de regresión (P1, media); el orden «pago total a tiempo» frente a «falta capturar» (B1); los casos sin regla de §8.2 (B4); `seq` fuera de rango (B7); el crédito `unsupported` en el calendario (B8); si cada caso del comparador lleva tope (Q2); y si los cortes de la tarjeta mueven el disponible (Q3).

## Coherencia con K12–K18 y K31 (d)
- **K12:** cumple. La tasa es una fracción con sus decimales capturados; ahora también la cotizada (`RateDecimals`), y sin decimales no hay texto, igual que C1.
- **K13:** el tope por persona es opcional, solo marca y desempata el objetivo 2.
- **K14:** los dos modos, partes externas HALF_UP con el titular tomando el resto, y pago de referencia solo en cotizaciones. La serie con intereses 0 queda PROVISIONAL con bloqueo (Q1 amplía el disparador).
- **K15:** §5.2 completo.
- **K16:** restricciones globales y ajustables, sin valor por defecto, aviso «faltan restricciones»; colchón en la línea usada, vivas y deuda en el crédito usado.
- **K17:** mín-máx entre los generados (`valid`, decisión C), desviación estándar con pico contra el tope, y unidades base por peso (PROVISIONAL `S3-objetivos-medida`).
- **K18:** datos por producto y `ErrMissingCapture`; la lista de 4 datos y la lectura del IVA son PROVISIONALES (`U4-iva-credito`).
- **K31:** `optimizer.date_window_days` propuesto en 7 días, sin valor sembrado.
- **K43 y K44:** 5 s leídos de `app_setting`, configurables, con 503 `unavailable` como dato.

## Los 3 casos del comparador (e)
§10.1 pide en cada caso los 4 valores de restricción en la fixture (capturas de fixture), ≥ 2 créditos con límite, ≥ 2 plazos con tasa, una promoción 0 % acotada en la tienda de peor U/C, ≥ 2 tiendas, compras existentes desiguales, un descartado por restricción, la intención de cada uno de los 4 óptimos y 4 óptimos distintos verificados a mano antes de correr. §10.0 fija la secuencia: tabla a mano, luego transcripción sin cambiar números, y solo entonces se corre. El ejemplo 2 trae las restricciones explícitas. Esto cumple en el contrato; las cifras las fija `docs/estado/S3-escenarios.md` (todavía no existe, lo escribe el optimizador). Ver Q2 y B9.

## Nota
`backend/internal/motor_cred/` y `oraculo/fixtures/s3/_datos/` ya existen sin commit (implementación en curso). No los evalué: esta revisión es del contrato.
