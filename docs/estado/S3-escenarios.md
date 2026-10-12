# S3 · Escenarios del comparador

Versión 1.0 (11-oct-2026, WF-1, agente `motor-S3-comparador-go`). Contrato: `docs/contratos/motores/S3-interfaz.md` §8 (ronda 4), `docs/criterios/S3.md` S3-12 a S3-14, decisiones C, K16, K17, K31 (7), K43 y K44, RN-11. Código: `backend/internal/motor_cred/comparador/`.

**Secuencia (S3-13, §10.0):** las secciones 1 a 4 (reglas, los 3 casos con sus restricciones e intención, y la tabla a mano completa) se escribieron **antes** de correr el comparador Go sobre los casos. La tabla sale de `backend/internal/motor_cred/comparador/testdata/tabla_mano.py`, un cálculo independiente (fracciones exactas de Python, disponible día por día, sin importar ni leer nada del motor Go) que también escribió las fixtures `testdata/caso-N.json` con su `esperado`, sin que nadie cambie un número. La sección 5 (tiempos) y la 6 (resultado de la corrida) se agregaron después. Desde la corrección `correccion-S3-r1` (11-oct-2026) la única fuente de la tabla son las fixtures de `oraculo/fixtures/s3/` (§7): las copias de `testdata/` se retiraron y `tabla_mano.py` ya no escribe, solo comprueba que esas fixtures son idénticas a la tabla.

## 1. Cómo se generan los escenarios (búsqueda exhaustiva, sin solver; S3-15)
- **Restricciones usadas (K16):** campo por campo, la de la comparación (`override`) si viene; si no, la global (`optimizer.*` de Ajustes › Créditos y condiciones). Si falta cualquiera de las 4 (`min_available_buffer_cents`, `max_live_purchases`, `zero_debt_by_date`, `date_window_days`) el resultado es `missing_constraints` con la lista en ese orden, `evaluated_count` 0, sin escenarios, sin óptimos y sin recomendación: la pantalla dice «faltan restricciones». No existe valor por defecto.
- **Créditos:** los de `credits`, en orden de `id`. Sin límite → `excluded` con `reason = no_limit` (PROVISIONAL `S3-objetivos-medida` (4)). Préstamo personal (`term_loan`): §8.2 no le da plazos, así que no aporta candidatos (SUPUESTO declarado en el reporte; no es un medio para pagar una compra planeada).
- **Fechas:** cada día de `desired_date − w` a `desired_date + w` (w = `date_window_days`), sin fechas anteriores a `today`.
- **Tienda:** la de la petición si trae `store_id` o `merchant_name`; si no, cada tienda de `baskets` (PROVISIONAL (1)); sin tienda ni canastas, un solo candidato sin tienda (SUPUESTO).
- **Plazo y promoción:** línea revolvente y crédito de tienda = plazos de `term_rates` ∪ plazos de las promociones aplicables; en cada plazo, «sin promoción» si el plazo tiene tasa, más cada promoción aplicable (comercio o tienda, plazo, compra mínima, vigencia; `ApplicablePromotions` del núcleo). Tarjeta = el cargo normal (plazo 1, sin promoción, 0.00) siempre, más cada promoción de meses sin intereses aplicable en su plazo. Un plazo sin tasa ni promoción no se genera (S3-06).
- **Cotización:** cada candidato se cotiza con `motor_cred.QuotePurchase` (ninguna fórmula duplicada). Si la cotización da error (sin tasa, promoción de tarjeta con tasa ≠ 0, tienda con datos K18, falta de primera fecha), el candidato no se genera ni cuenta en `evaluated_count`.

## 2. Restricciones duras, medidas y puntaje
- **Ventana** de un escenario = [fecha de compra, su último pago].
- **Colchón (lectura del crédito, vacío = B):** `AvailabilityTimeline` del núcleo con `from = today`, `to` = último pago, las compras existentes más la del escenario. El disponible de un día es el que queda después del último evento con fecha ≤ ese día (el valor con que abre la ventana cuenta). `min_available_cents` = mínimo de la ventana; cumple si ≥ colchón.
- **Compras vivas:** máximo de compras vivas simultáneas del crédito usado en la ventana, contando la del escenario; una compra está viva desde su fecha hasta su último pago, ambos incluidos (RN-21); las canceladas no cuentan. Cumple si ≤ máximo.
- **Deuda en cero:** último vencimiento del crédito usado con la compra añadida; cumple si ≤ la fecha (un pago ese mismo día cumple).
- Si viola varias, `violated_constraint` es la primera en el orden colchón, vivas, deuda. El tope por persona **no** descarta: solo marca `over_cap` y desempata el objetivo 2.
- **Objetivo 1, intereses:** `interest_cents` de la cotización; menor es mejor.
- **Objetivo 2, pago parejo (K17; serie PROVISIONAL (2)):** quincenas consecutivas de la quincena de la compra a la de su último pago, incluidas las de total 0; total de cada quincena = suma de los pagos del calendario unificado (todas las compras no canceladas de todos los créditos más el escenario). `std_dev_milli_cents` = ⌊(⌊isqrt(4N) / n⌋ + 1) / 2⌋ con N = 10⁶ × (n·Σxᵢ² − (Σxᵢ)²): es HALF_UP(1000 × desviación estándar poblacional) exacto (0, 100 y 200 centavos → 81,650). Pico = `max_fortnight_payment_cents`; `fortnight_payment_spread_cents` = máximo − mínimo. `over_cap`: cada persona con tope cuyo pico de sus partes por quincena en esa serie pasa su tope. Menor desviación es mejor.
- **Objetivo 3, disponible:** `min_available_cents` de la ventana; mayor es mejor.
- **Objetivo 4, comprar más por peso (K17; PROVISIONAL (1)):** `value_per_peso_milli` = HALF_UP(monto × U × 100 / (C × (monto + intereses))), con U = unidades base de la canasta en milésimas y C = su costo en centavos en la tienda del escenario (milésimas de unidad base por peso gastado). Sin canasta de esa tienda, nulo. Mayor es mejor.
- **Normalización mín-máx (K17, RN-11):** solo entre los escenarios generados en el sentido de la decisión C, los que respetan las restricciones (`valid`); los descartados no mueven la escala. Menor-es-mejor: (máx − v)/(máx − mín); mayor-es-mejor: (v − mín)/(máx − mín), en fracción exacta (`big.Rat`).
- **Declaración de empate (S3-14):** si máximo = mínimo, o todos los valores son nulos, el objetivo vale **1** en todos los escenarios (no distingue: suma la misma constante a todos y el orden no cambia; se elige 1 porque todos tienen el mejor valor). Un nulo entre valores no nulos no puede ocurrir por construcción (todos los escenarios de una comparación comparten el conjunto de tiendas); si ocurriera, el nulo también vale 1.
- **Puntaje:** Σ wᵢ·nᵢ / Σ wᵢ con los pesos de la petición (`optimizer.weights`, enteros ≥ 0; 1, 1, 1, 1 → 0.25 cada uno, Σ = 1.00); Σ wᵢ = 0 o un peso negativo → `ErrInvalid`. `score_bp` = HALF_UP(10000 × puntaje); los descartados llevan `score_bp` nulo.
- **Óptimos:** el mejor valor de cada objetivo; empate → (solo objetivo 2, PROVISIONAL (2)) menos personas en `over_cap` → mayor puntaje exacto → llave canónica menor (`credit_id`, `term_count`, `purchase_date`, `store_id` nulo primero, `promotion_id` nulo primero). **Recomendación:** mayor puntaje exacto; empate → llave canónica. `valid` y `rejected` van en orden de llave canónica. Ningún válido → `no_valid_scenario` sin óptimos ni recomendación.
- **Listas `provisional` (§14.1):** escenario = `S3-objetivos-medida` ∪ la de su compra ∪ las de las entradas del calendario que forman su serie del objetivo 2; resultado = `S3-objetivos-medida` ∪ todas las de `valid` y `rejected` (`[]` con `missing_constraints`), sin repetidos y en orden por bytes.

## 3. Tiempo máximo y comportamiento al agotarse (K43, K44; S3-12)
`comparator.timeout_seconds` (valor sembrado 5 en `00003`, configurable en Ajustes › Créditos y condiciones) se convierte con `TimeoutFromSetting` (negativo → `ErrInvalid`; 0 vence de inmediato) y quien llama crea `context.WithTimeout`. `Compare` revisa el contexto antes de cada candidato y al terminar; si venció devuelve `motor_cred.ErrTimeout` **sin resultados parciales**, y `motor_cred.ResponseFor` lo traduce a **503 `unavailable`**. Prueba: `TestRestriccionesTiempo` (0 s → `ErrTimeout` y 503; 5 s → termina) y la fixture `testdata/tiempo.json` (`{"error": "timeout"}`).

## 4. Los 3 casos (fijados antes de correr; capturas de fixture, no valores de Axel)
Comunes: `today` 2026-11-01; un solo titular (`axel`, sin partes externas); líneas quincenales con 15 días mínimos al primer pago y lectura B; pesos 1, 1, 1, 1; tienda 1 «Tienda Barata» (mejor U/C) y tienda 2 «Tienda Cara» (peor U/C, la única con promoción). La diferencia de U/C entre tiendas (−33 % en los casos 1 y 3, −20 % en el 2) supera el efecto de los intereses de los plazos de 2 quincenas (−7.4 % a −14.7 %) y, en los casos 1 y 3, también el de 4 quincenas (−13 % y −20.6 %); en el caso 2 la 4 q de Línea Uno (−20.6 %) casi la iguala, y el óptimo (4) es el cargo normal sin intereses.

| | Caso 1 · dos líneas tipo Kueski | Caso 2 · línea + tarjeta con MSI | Caso 3 · tope, `override` y crédito sin límite |
|---|---|---|---|
| Compra | 2,000.00 el 30-nov-2026 ± 1 día | 3,000.00 el 30-nov-2026 ± 1 día | 1,800.00 el 30-nov-2026 ± 1 día |
| Créditos | 1 «Línea Uno» límite 15,000.00, tasas 2 q 17.27 %, 4 q 25.8685 %, 6 q 35 % (fixture); 2 «Línea Dos» límite 6,000.00, 2 q 8 %, 4 q 15 % | 1 «Línea Uno» límite 8,000.00, 2 q, 4 q y 6 q como en el caso 1; 3 «Tarjeta» límite 6,500.00, corte día 5, pago día 20 | los dos del caso 1 más 4 «Crédito sin límite» (sin `limit_cents`: se excluye `no_limit`) |
| Promoción 0 % | id 1, Línea Uno, Tienda Cara, 4 q, mínimo 1,000.00, del 28 al 29-nov-2026 | id 2, Tarjeta, MSI 2 meses, Tienda Cara, mínimo 2,000.00, del 29 al 30-nov-2026 | la del caso 1 |
| Compras existentes | Línea Uno: L1-f 20-ene-2027, 2,500.00 a 1 pago con intereses cotizados 0.00 (paga 15-feb-2027); Línea Dos: L2-a 25-oct-2026, 1,165.50 a 2 q con intereses cotizados 93.18 (pagos 629.34 el 15 y el 30-nov), L2-b 1-nov-2026, 1,000.00 a 4 q (287.50 del 30-nov al 15-ene) | Línea Uno: L1-a 20-oct-2026, 2,000.00 a 4 q (15-nov a 30-dic), L1-b 25-oct-2026, 1,000.00 a 2 q (15 y 30-nov), L1-f 20-ene-2027, 2,600.00 a 2 q (15 y 28-feb-2027) | las del caso 1 |
| Canastas (U, C) | tienda 1: 60.000 unidades, 2,000.00; tienda 2: 40.000, 2,000.00 | tienda 1: 90.000, 3,000.00; tienda 2: 72.000, 3,000.00 | las del caso 1 |
| Restricciones | global: colchón 3,000.00, vivas 2, deuda en cero 15-feb-2027, ± 1 día; sin `override` | global: colchón 3,200.00, vivas 2, deuda en cero 28-feb-2027, ± 1 día | global: colchón 3,500.00, vivas 2, deuda en cero 31-dic-2026, ± 1 día; `override`: colchón 3,000.00 y deuda en cero 15-feb-2027 (vivas y ventana salen de la global) |
| Tope por persona | axel 1,200.00 | axel 2,400.00 | axel 1,000.00 |
| Descartes por diseño | colchón: Línea Dos el 29-nov (antes de que L2-a y L2-b repongan capital el 30); vivas: Línea Dos el 30-nov (L2-a sigue viva ese día); deuda: 6 q (terminan después del 15-feb) | colchón: Línea Uno el 29-nov; vivas: Línea Uno el 30-nov (L1-a, L1-b y la nueva); deuda: 6 q el 1-dic | como el caso 1; con las globales el colchón descartaría Línea Dos el 1-dic y la deuda en cero todo lo de 4 q: el `override` cambia el conjunto |
| Intención de los 4 óptimos | (1) la promoción (0.00); (2) Línea Uno a 4 q, que rellena las quincenas bajas que dejan L2-a y L2-b; (3) Línea Uno a 2 q, cuya ventana no alcanza la compra L1-f del 20-ene (la de 4 q sí); (4) Tienda Barata con el plazo de menor tasa (Línea Dos 2 q el 1-dic) | (1) MSI: empata en 0.00 con el cargo normal y gana por diseño en el puntaje, porque reparte el pago (desviación mucho menor que el pago único del cargo normal); (2) Línea Uno a 4 q el 1-dic; (3) Línea Uno a 2 q el 1-dic (la tarjeta tiene menos holgura); (4) el cargo normal en Tienda Barata (0.00 de intereses y mejor U/C) | como el caso 1 con otro monto: (1) promoción; (2) 4 q; (3) 2 q; (4) Línea Dos 2 q en Tienda Barata, que solo es válida gracias al `override` del colchón |

Los 3 casos tienen ≥ 2 créditos con límite, ≥ 2 plazos con tasa, una sola promoción 0 % aplicable solo en algunas fechas y solo en la tienda de peor U/C, compras existentes que hacen desigual el perfil por quincena, al menos un descartado por cada restricción y al menos un válido sobre el tope (lo comprueba `tabla_mano.py`).

## 4.1 Tablas a mano
Montos en pesos (la fixture los trae en centavos). Llave = crédito/plazo/fecha/tienda/promoción. «Desv.» = `std_dev_milli_cents`; «Pico» = `max_fortnight_payment_cents`; «Unid. por peso» = `value_per_peso_milli`.

### Tabla a mano · caso-1

`status` ok, `evaluated_count` 31 (17 válidos, 14 descartados), excluidos: ninguno.

| Llave (crédito/plazo/fecha/tienda/promoción) | Intereses | Desv. (milésimas de centavo) | Pico | Disponible mín. | Vivas | Deuda en cero | Unid. por peso (milésimas) | Puntaje (pb) | Sobre tope | Resultado |
|---|---:|---:|---:|---:|---:|---|---:|---:|---|---|
| 1/2q/2026-11-29/t1/pNone | 345.45 | 25,615,415 | 1,460.23 | 13,000.00 | 1 | 2027-02-15 | 26 | 7,517 | axel 1,460.23 > 1,200.00 | válido |
| 1/2q/2026-11-29/t2/pNone | 345.45 | 25,615,415 | 1,460.23 | 13,000.00 | 1 | 2027-02-15 | 17 | 5,642 | axel 1,460.23 > 1,200.00 | válido |
| 1/2q/2026-11-30/t1/pNone | 345.45 | 25,615,415 | 1,460.23 | 13,000.00 | 1 | 2027-02-15 | 26 | 7,517 | axel 1,460.23 > 1,200.00 | válido |
| 1/2q/2026-11-30/t2/pNone | 345.45 | 25,615,415 | 1,460.23 | 13,000.00 | 1 | 2027-02-15 | 17 | 5,642 | axel 1,460.23 > 1,200.00 | válido |
| 1/2q/2026-12-01/t1/pNone | 345.45 | 55,282,787 | 1,460.23 | 13,000.00 | 1 | 2027-02-15 | 26 | 6,683 | axel 1,460.23 > 1,200.00 | válido |
| 1/2q/2026-12-01/t2/pNone | 345.45 | 55,282,787 | 1,460.23 | 13,000.00 | 1 | 2027-02-15 | 17 | 4,808 | axel 1,460.23 > 1,200.00 | válido |
| 1/4q/2026-11-29/t1/pNone | 517.37 | 11,499,600 | 916.84 | 12,000.00 | 2 | 2027-02-15 | 24 | 6,410 | — | válido |
| 1/4q/2026-11-29/t2/pNone | 517.37 | 11,499,600 | 916.84 | 12,000.00 | 2 | 2027-02-15 | 16 | 4,744 | — | válido |
| 1/4q/2026-11-29/t2/p1 | 0.00 | 13,739,153 | 916.84 | 12,000.00 | 2 | 2027-02-15 | 20 | 8,014 | — | válido |
| 1/4q/2026-11-30/t1/pNone | 517.37 | 11,499,600 | 916.84 | 12,000.00 | 2 | 2027-02-15 | 24 | 6,410 | — | válido |
| 1/4q/2026-11-30/t2/pNone | 517.37 | 11,499,600 | 916.84 | 12,000.00 | 2 | 2027-02-15 | 16 | 4,744 | — | válido |
| 1/4q/2026-12-01/t1/pNone | 517.37 | 100,381,009 | 3,129.35 | 11,500.00 | 2 | 2027-02-15 | 24 | 3,782 | axel 3,129.35 > 1,200.00 | válido |
| 1/4q/2026-12-01/t2/pNone | 517.37 | 100,381,009 | 3,129.35 | 11,500.00 | 2 | 2027-02-15 | 16 | 2,115 | axel 3,129.35 > 1,200.00 | válido |
| 2/2q/2026-12-01/t1/pNone | 160.00 | 50,911,688 | 1,367.50 | 3,250.00 | 2 | 2027-01-15 | 28 | 5,618 | axel 1,367.50 > 1,200.00 | válido |
| 2/2q/2026-12-01/t2/pNone | 160.00 | 50,911,688 | 1,367.50 | 3,250.00 | 2 | 2027-01-15 | 19 | 3,743 | axel 1,367.50 > 1,200.00 | válido |
| 2/4q/2026-12-01/t1/pNone | 300.00 | 99,437,795 | 3,075.00 | 3,250.00 | 2 | 2027-02-15 | 26 | 3,160 | axel 3,075.00 > 1,200.00 | válido |
| 2/4q/2026-12-01/t2/pNone | 300.00 | 99,437,795 | 3,075.00 | 3,250.00 | 2 | 2027-02-15 | 17 | 1,285 | axel 3,075.00 > 1,200.00 | válido |
| 1/6q/2026-11-29/t1/pNone | 700.00 | 81,249,821 | 2,950.00 | 11,499.99 | 2 | 2027-02-28 | 22 | — | axel 2,950.00 > 1,200.00 | descartado: zero_debt_by_date |
| 1/6q/2026-11-29/t2/pNone | 700.00 | 81,249,821 | 2,950.00 | 11,499.99 | 2 | 2027-02-28 | 15 | — | axel 2,950.00 > 1,200.00 | descartado: zero_debt_by_date |
| 1/6q/2026-11-30/t1/pNone | 700.00 | 81,249,821 | 2,950.00 | 11,499.99 | 2 | 2027-02-28 | 22 | — | axel 2,950.00 > 1,200.00 | descartado: zero_debt_by_date |
| 1/6q/2026-11-30/t2/pNone | 700.00 | 81,249,821 | 2,950.00 | 11,499.99 | 2 | 2027-02-28 | 15 | — | axel 2,950.00 > 1,200.00 | descartado: zero_debt_by_date |
| 1/6q/2026-12-01/t1/pNone | 700.00 | 86,436,688 | 2,950.00 | 11,166.66 | 2 | 2027-03-15 | 22 | — | axel 2,950.00 > 1,200.00 | descartado: zero_debt_by_date |
| 1/6q/2026-12-01/t2/pNone | 700.00 | 86,436,688 | 2,950.00 | 11,166.66 | 2 | 2027-03-15 | 15 | — | axel 2,950.00 > 1,200.00 | descartado: zero_debt_by_date |
| 2/2q/2026-11-29/t1/pNone | 160.00 | 21,244,316 | 1,367.50 | 2,417.25 | 3 | 2027-01-15 | 28 | — | axel 1,367.50 > 1,200.00 | descartado: min_available_buffer |
| 2/2q/2026-11-29/t2/pNone | 160.00 | 21,244,316 | 1,367.50 | 2,417.25 | 3 | 2027-01-15 | 19 | — | axel 1,367.50 > 1,200.00 | descartado: min_available_buffer |
| 2/2q/2026-11-30/t1/pNone | 160.00 | 21,244,316 | 1,367.50 | 3,250.00 | 3 | 2027-01-15 | 28 | — | axel 1,367.50 > 1,200.00 | descartado: max_live_purchases |
| 2/2q/2026-11-30/t2/pNone | 160.00 | 21,244,316 | 1,367.50 | 3,250.00 | 3 | 2027-01-15 | 19 | — | axel 1,367.50 > 1,200.00 | descartado: max_live_purchases |
| 2/4q/2026-11-29/t1/pNone | 300.00 | 12,225,904 | 916.84 | 2,417.25 | 3 | 2027-01-30 | 26 | — | — | descartado: min_available_buffer |
| 2/4q/2026-11-29/t2/pNone | 300.00 | 12,225,904 | 916.84 | 2,417.25 | 3 | 2027-01-30 | 17 | — | — | descartado: min_available_buffer |
| 2/4q/2026-11-30/t1/pNone | 300.00 | 12,225,904 | 916.84 | 3,250.00 | 3 | 2027-01-30 | 26 | — | — | descartado: max_live_purchases |
| 2/4q/2026-11-30/t2/pNone | 300.00 | 12,225,904 | 916.84 | 3,250.00 | 3 | 2027-01-30 | 17 | — | — | descartado: max_live_purchases |

Óptimos: lowest_interest → `1/4q/2026-11-29/t2/p1`; even_low_payment → `1/4q/2026-11-29/t1/pNone`; max_available → `1/2q/2026-11-29/t1/pNone`; most_value_per_peso → `2/2q/2026-12-01/t1/pNone`. Recomendación → `1/4q/2026-11-29/t2/p1` (8,014 pb). `provisional` del resultado: S3-objetivos-medida, S3-serie-sin-intereses.

### Tabla a mano · caso-2

`status` ok, `evaluated_count` 26 (12 válidos, 14 descartados), excluidos: ninguno.

| Llave (crédito/plazo/fecha/tienda/promoción) | Intereses | Desv. (milésimas de centavo) | Pico | Disponible mín. | Vivas | Deuda en cero | Unid. por peso (milésimas) | Puntaje (pb) | Sobre tope | Resultado |
|---|---:|---:|---:|---:|---:|---|---:|---:|---|---|
| 1/2q/2026-12-01/t1/pNone | 518.18 | 72,777,051 | 2,388.44 | 4,000.00 | 2 | 2027-02-28 | 26 | 7,205 | — | válido |
| 1/2q/2026-12-01/t2/pNone | 518.18 | 72,777,051 | 2,388.44 | 4,000.00 | 2 | 2027-02-28 | 20 | 5,842 | — | válido |
| 1/4q/2026-12-01/t1/pNone | 776.05 | 65,464,463 | 2,468.56 | 3,900.00 | 2 | 2027-02-28 | 24 | 5,636 | axel 2,468.56 > 2,400.00 | válido |
| 1/4q/2026-12-01/t2/pNone | 776.05 | 65,464,463 | 2,468.56 | 3,900.00 | 2 | 2027-02-28 | 19 | 4,500 | axel 2,468.56 > 2,400.00 | válido |
| 3/1q/2026-11-29/t1/pNone | 0.00 | 129,827,012 | 3,629.35 | 3,500.00 | 1 | 2026-12-20 | 30 | 5,597 | axel 3,629.35 > 2,400.00 | válido |
| 3/1q/2026-11-29/t2/pNone | 0.00 | 129,827,012 | 3,629.35 | 3,500.00 | 1 | 2026-12-20 | 24 | 4,233 | axel 3,629.35 > 2,400.00 | válido |
| 3/1q/2026-11-30/t1/pNone | 0.00 | 129,827,012 | 3,629.35 | 3,500.00 | 1 | 2026-12-20 | 30 | 5,597 | axel 3,629.35 > 2,400.00 | válido |
| 3/1q/2026-11-30/t2/pNone | 0.00 | 129,827,012 | 3,629.35 | 3,500.00 | 1 | 2026-12-20 | 24 | 4,233 | axel 3,629.35 > 2,400.00 | válido |
| 3/1q/2026-12-01/t1/pNone | 0.00 | 150,000,500 | 3,629.35 | 3,500.00 | 1 | 2026-12-20 | 30 | 5,000 | axel 3,629.35 > 2,400.00 | válido |
| 3/1q/2026-12-01/t2/pNone | 0.00 | 150,000,500 | 3,629.35 | 3,500.00 | 1 | 2026-12-20 | 24 | 3,636 | axel 3,629.35 > 2,400.00 | válido |
| 3/2q/2026-11-29/t2/p2 | 0.00 | 72,997,808 | 2,129.35 | 3,500.00 | 1 | 2027-01-20 | 24 | 5,914 | — | válido |
| 3/2q/2026-11-30/t2/p2 | 0.00 | 72,997,808 | 2,129.35 | 3,500.00 | 1 | 2027-01-20 | 24 | 5,914 | — | válido |
| 1/2q/2026-11-29/t1/pNone | 518.18 | 55,282,787 | 2,388.44 | 3,000.00 | 3 | 2027-02-28 | 26 | — | — | descartado: min_available_buffer |
| 1/2q/2026-11-29/t2/pNone | 518.18 | 55,282,787 | 2,388.44 | 3,000.00 | 3 | 2027-02-28 | 20 | — | — | descartado: min_available_buffer |
| 1/2q/2026-11-30/t1/pNone | 518.18 | 55,282,787 | 2,388.44 | 4,000.00 | 3 | 2027-02-28 | 26 | — | — | descartado: max_live_purchases |
| 1/2q/2026-11-30/t2/pNone | 518.18 | 55,282,787 | 2,388.44 | 4,000.00 | 3 | 2027-02-28 | 20 | — | — | descartado: max_live_purchases |
| 1/4q/2026-11-29/t1/pNone | 776.05 | 28,197,387 | 1,573.36 | 3,000.00 | 3 | 2027-02-28 | 24 | — | — | descartado: min_available_buffer |
| 1/4q/2026-11-29/t2/pNone | 776.05 | 28,197,387 | 1,573.36 | 3,000.00 | 3 | 2027-02-28 | 19 | — | — | descartado: min_available_buffer |
| 1/4q/2026-11-30/t1/pNone | 776.05 | 28,197,387 | 1,573.36 | 4,000.00 | 3 | 2027-02-28 | 24 | — | — | descartado: max_live_purchases |
| 1/4q/2026-11-30/t2/pNone | 776.05 | 28,197,387 | 1,573.36 | 4,000.00 | 3 | 2027-02-28 | 19 | — | — | descartado: max_live_purchases |
| 1/6q/2026-11-29/t1/pNone | 1,050.00 | 58,244,895 | 2,199.55 | 3,000.00 | 3 | 2027-02-28 | 22 | — | — | descartado: min_available_buffer |
| 1/6q/2026-11-29/t2/pNone | 1,050.00 | 58,244,895 | 2,199.55 | 3,000.00 | 3 | 2027-02-28 | 18 | — | — | descartado: min_available_buffer |
| 1/6q/2026-11-30/t1/pNone | 1,050.00 | 58,244,895 | 2,199.55 | 3,900.00 | 3 | 2027-02-28 | 22 | — | — | descartado: max_live_purchases |
| 1/6q/2026-11-30/t2/pNone | 1,050.00 | 58,244,895 | 2,199.55 | 3,900.00 | 3 | 2027-02-28 | 18 | — | — | descartado: max_live_purchases |
| 1/6q/2026-12-01/t1/pNone | 1,050.00 | 67,202,705 | 2,199.55 | 3,400.00 | 2 | 2027-03-15 | 22 | — | — | descartado: zero_debt_by_date |
| 1/6q/2026-12-01/t2/pNone | 1,050.00 | 67,202,705 | 2,199.55 | 3,400.00 | 2 | 2027-03-15 | 18 | — | — | descartado: zero_debt_by_date |

Óptimos: lowest_interest → `3/2q/2026-11-29/t2/p2`; even_low_payment → `1/4q/2026-12-01/t1/pNone`; max_available → `1/2q/2026-12-01/t1/pNone`; most_value_per_peso → `3/1q/2026-11-29/t1/pNone`. Recomendación → `1/2q/2026-12-01/t1/pNone` (7,205 pb). `provisional` del resultado: S3-formula-tarjeta, S3-objetivos-medida, S3-origen-cargo-normal, S3-serie-sin-intereses.

### Tabla a mano · caso-3

`status` ok, `evaluated_count` 31 (17 válidos, 14 descartados), excluidos: 4 no_limit.

| Llave (crédito/plazo/fecha/tienda/promoción) | Intereses | Desv. (milésimas de centavo) | Pico | Disponible mín. | Vivas | Deuda en cero | Unid. por peso (milésimas) | Puntaje (pb) | Sobre tope | Resultado |
|---|---:|---:|---:|---:|---:|---|---:|---:|---|---|
| 1/2q/2026-11-29/t1/pNone | 310.91 | 20,087,254 | 1,342.96 | 13,200.00 | 1 | 2027-02-15 | 26 | 7,692 | axel 1,342.96 > 1,000.00 | válido |
| 1/2q/2026-11-29/t2/pNone | 310.91 | 20,087,254 | 1,342.96 | 13,200.00 | 1 | 2027-02-15 | 17 | 5,817 | axel 1,342.96 > 1,000.00 | válido |
| 1/2q/2026-11-30/t1/pNone | 310.91 | 20,087,254 | 1,342.96 | 13,200.00 | 1 | 2027-02-15 | 26 | 7,692 | axel 1,342.96 > 1,000.00 | válido |
| 1/2q/2026-11-30/t2/pNone | 310.91 | 20,087,254 | 1,342.96 | 13,200.00 | 1 | 2027-02-15 | 17 | 5,817 | axel 1,342.96 > 1,000.00 | válido |
| 1/2q/2026-12-01/t1/pNone | 310.91 | 49,754,626 | 1,342.96 | 13,200.00 | 1 | 2027-02-15 | 26 | 6,839 | axel 1,342.96 > 1,000.00 | válido |
| 1/2q/2026-12-01/t2/pNone | 310.91 | 49,754,626 | 1,342.96 | 13,200.00 | 1 | 2027-02-15 | 17 | 4,964 | axel 1,342.96 > 1,000.00 | válido |
| 1/4q/2026-11-29/t1/pNone | 465.63 | 12,370,748 | 916.84 | 12,050.00 | 2 | 2027-02-15 | 24 | 6,372 | — | válido |
| 1/4q/2026-11-29/t2/pNone | 465.63 | 12,370,748 | 916.84 | 12,050.00 | 2 | 2027-02-15 | 16 | 4,705 | — | válido |
| 1/4q/2026-11-29/t2/p1 | 0.00 | 14,998,625 | 916.84 | 12,050.00 | 2 | 2027-02-15 | 20 | 7,963 | — | válido |
| 1/4q/2026-11-30/t1/pNone | 465.63 | 12,370,748 | 916.84 | 12,050.00 | 2 | 2027-02-15 | 24 | 6,372 | — | válido |
| 1/4q/2026-11-30/t2/pNone | 465.63 | 12,370,748 | 916.84 | 12,050.00 | 2 | 2027-02-15 | 16 | 4,705 | — | válido |
| 1/4q/2026-12-01/t1/pNone | 465.63 | 99,293,294 | 3,066.43 | 11,600.00 | 2 | 2027-02-15 | 24 | 3,756 | axel 3,066.43 > 1,000.00 | válido |
| 1/4q/2026-12-01/t2/pNone | 465.63 | 99,293,294 | 3,066.43 | 11,600.00 | 2 | 2027-02-15 | 16 | 2,090 | axel 3,066.43 > 1,000.00 | válido |
| 2/2q/2026-12-01/t1/pNone | 144.00 | 45,820,519 | 1,259.50 | 3,450.00 | 2 | 2027-01-15 | 28 | 5,765 | axel 1,259.50 > 1,000.00 | válido |
| 2/2q/2026-12-01/t2/pNone | 144.00 | 45,820,519 | 1,259.50 | 3,450.00 | 2 | 2027-01-15 | 19 | 3,890 | axel 1,259.50 > 1,000.00 | válido |
| 2/4q/2026-12-01/t1/pNone | 270.00 | 98,482,562 | 3,017.50 | 3,450.00 | 2 | 2027-02-15 | 26 | 3,157 | axel 3,017.50 > 1,000.00 | válido |
| 2/4q/2026-12-01/t2/pNone | 270.00 | 98,482,562 | 3,017.50 | 3,450.00 | 2 | 2027-02-15 | 17 | 1,282 | axel 3,017.50 > 1,000.00 | válido |
| 1/6q/2026-11-29/t1/pNone | 630.00 | 81,201,604 | 2,905.00 | 11,600.00 | 2 | 2027-02-28 | 22 | — | axel 2,905.00 > 1,000.00 | descartado: zero_debt_by_date |
| 1/6q/2026-11-29/t2/pNone | 630.00 | 81,201,604 | 2,905.00 | 11,600.00 | 2 | 2027-02-28 | 15 | — | axel 2,905.00 > 1,000.00 | descartado: zero_debt_by_date |
| 1/6q/2026-11-30/t1/pNone | 630.00 | 81,201,604 | 2,905.00 | 11,600.00 | 2 | 2027-02-28 | 22 | — | axel 2,905.00 > 1,000.00 | descartado: zero_debt_by_date |
| 1/6q/2026-11-30/t2/pNone | 630.00 | 81,201,604 | 2,905.00 | 11,600.00 | 2 | 2027-02-28 | 15 | — | axel 2,905.00 > 1,000.00 | descartado: zero_debt_by_date |
| 1/6q/2026-12-01/t1/pNone | 630.00 | 86,019,724 | 2,905.00 | 11,300.00 | 2 | 2027-03-15 | 22 | — | axel 2,905.00 > 1,000.00 | descartado: zero_debt_by_date |
| 1/6q/2026-12-01/t2/pNone | 630.00 | 86,019,724 | 2,905.00 | 11,300.00 | 2 | 2027-03-15 | 15 | — | axel 2,905.00 > 1,000.00 | descartado: zero_debt_by_date |
| 2/2q/2026-11-29/t1/pNone | 144.00 | 16,153,147 | 1,259.50 | 2,617.25 | 3 | 2027-01-15 | 28 | — | axel 1,259.50 > 1,000.00 | descartado: min_available_buffer |
| 2/2q/2026-11-29/t2/pNone | 144.00 | 16,153,147 | 1,259.50 | 2,617.25 | 3 | 2027-01-15 | 19 | — | axel 1,259.50 > 1,000.00 | descartado: min_available_buffer |
| 2/2q/2026-11-30/t1/pNone | 144.00 | 16,153,147 | 1,259.50 | 3,450.00 | 3 | 2027-01-15 | 28 | — | axel 1,259.50 > 1,000.00 | descartado: max_live_purchases |
| 2/2q/2026-11-30/t2/pNone | 144.00 | 16,153,147 | 1,259.50 | 3,450.00 | 3 | 2027-01-15 | 19 | — | axel 1,259.50 > 1,000.00 | descartado: max_live_purchases |
| 2/4q/2026-11-29/t1/pNone | 270.00 | 13,341,151 | 916.84 | 2,617.25 | 3 | 2027-01-30 | 26 | — | — | descartado: min_available_buffer |
| 2/4q/2026-11-29/t2/pNone | 270.00 | 13,341,151 | 916.84 | 2,617.25 | 3 | 2027-01-30 | 17 | — | — | descartado: min_available_buffer |
| 2/4q/2026-11-30/t1/pNone | 270.00 | 13,341,151 | 916.84 | 3,450.00 | 3 | 2027-01-30 | 26 | — | — | descartado: max_live_purchases |
| 2/4q/2026-11-30/t2/pNone | 270.00 | 13,341,151 | 916.84 | 3,450.00 | 3 | 2027-01-30 | 17 | — | — | descartado: max_live_purchases |

Óptimos: lowest_interest → `1/4q/2026-11-29/t2/p1`; even_low_payment → `1/4q/2026-11-29/t1/pNone`; max_available → `1/2q/2026-11-29/t1/pNone`; most_value_per_peso → `2/2q/2026-12-01/t1/pNone`. Recomendación → `1/4q/2026-11-29/t2/p1` (7,963 pb). `provisional` del resultado: S3-objetivos-medida, S3-serie-sin-intereses.

## 5. Tiempos medidos (después de correr; `TestRestriccionesTiempo`, 5 s por caso)
Comando: `cd backend && go test ./internal/motor_cred/comparador/ -run TestRestriccionesTiempo -count=3 -v | grep S3-TIEMPO` (11-oct-2026, Go 1.27.2, equipo de desarrollo de 16 hilos; el comparador corre en un solo hilo).

| Caso | Candidatos | Corrida 1 | Corrida 2 | Corrida 3 |
|---|---:|---:|---:|---:|
| caso-1 | 31 | 0 ms | 0 ms | 0 ms |
| caso-2 | 26 | 0 ms | 1 ms | 0 ms |
| caso-3 | 31 | 0 ms | 0 ms | 1 ms |
| carga (6 créditos × 7 plazos × 15 fechas × 6 tiendas + 1 promoción por crédito) | 3,870 | 111 ms | 123 ms | 118 ms |

Ningún caso pasa de 5,000 ms: no se abre `docs/estado/bloqueos/S3-tiempo-optimizador.md` y el valor de 5 s (K43, FIRME, configurable por K44) no cambia. El peor caso medido usa el 2.5 % del tiempo máximo. La rejilla de carga vive en la fixture `oraculo/fixtures/s3/s3-comparador-carga.json` (volcada con `json.Marshal` de la antigua `carga()` de `comparador_test.go`, que se retiró, sin cambiar un número; §7); su cuenta a mano es 6 × (15 × 6 × 7 + 15) = 3,870 y la comprueban la prueba, el oráculo y `compara.py s3`.

## 6. Resultado de la corrida contra la tabla a mano
- `go test ./internal/motor_cred/comparador/ -count=1 -cover -v`: los 3 casos reproducen **toda** la tabla a mano (`status`, `evaluated_count`, `valid` completo con valores, puntajes y listas `provisional`, `rejected` con su restricción, `excluded`, los 4 `optima` y `recommendation`): 0 diferencias, así que no hay hallazgo de cifras que documentar. Cobertura del paquete 96.0 %.
- Óptimos obtenidos (iguales a la intención de §4): caso 1: intereses `1/4q/2026-11-29/t2/p1`, parejo `1/4q/2026-11-29/t1`, disponible `1/2q/2026-11-29/t1`, por peso `2/2q/2026-12-01/t1`, recomendación la promoción (8,014 pb). Caso 2: intereses MSI `3/2q/2026-11-29/t2/p2` (empata en 0.00 con los cargos normales y gana por puntaje: 5,914 contra 5,597 del mejor cargo normal), parejo `1/4q/2026-12-01/t1`, disponible `1/2q/2026-12-01/t1`, por peso el cargo normal `3/1q/2026-11-29/t1`, recomendación `1/2q/2026-12-01/t1` (7,205 pb). Caso 3: los mismos papeles que el caso 1 con 1,800.00 (recomendación la promoción, 7,963 pb); Línea Dos el 1-dic solo es válida por el `override` del colchón.
- **Hallazgo para el bloqueo `S3-objetivos-medida` (pregunta 2, serie del objetivo 2):** con la serie de §8.4, una compra cuyo único pago cae en la misma quincena de la compra (p. ej. un cargo normal de tarjeta comprado el 1-dic con pago el 15-dic) tiene una serie de una sola quincena y desviación 0: «pago parejo» la califica como la más pareja aunque sea un pago único grande. Los casos se diseñaron para no caer en eso (corte el 5 y pago el 20), pero conviene que Axel lo vea al decidir la medida (por ejemplo, medir sobre una ventana fija de quincenas o sumar el pico contra el tope).

## 7. Integración con el oráculo (correccion-S3-r1, 11-oct-2026; REVISION-S3-r1 H1, H2 y H3)
- **Fixtures (§10.0 de la interfaz):** `oraculo/fixtures/s3/_datos/comparador-caso-{1,2,3}.json` (el `entradas.datos` de cada caso) y `s3-comparador-{caso-1,caso-2,caso-3,sin-restricciones,tiempo,carga}.json`. Los cinco primeros son las fixtures de `testdata/` copiadas sin cambiar un número, con `datos` = `{"$ref": "_datos/comparador-caso-N.json"}`; `-sin-restricciones` y `-tiempo` apuntan al caso 1 (sus `datos` eran idénticos). `-carga` lleva la rejilla en línea y su `esperado` es solo `status` «ok» y `evaluated_count` 3,870 (cuenta a mano de §5); PROVISIONAL `S3-objetivos-medida` como los casos 1 a 3. `python3 -I backend/internal/motor_cred/comparador/testdata/tabla_mano.py` recalcula la tabla y da «OK» en las 5 fixtures de la tabla.
- **Vía oficial:** la consulta «comparar» de `evaluar.Evaluate` llama `comparador.Query` (antes devolvía `unsupported`), así que `motores s3 dump` y `compara.py s3` incluyen el comparador.
- **Corrida:** `cred_escenarios.py` termina con código 0 (aceptados con violación = 0; 4 óptimos distintos en cada caso; `sin-restricciones` → `missing_constraints`; `tiempo` → `timeout`; `carga` 3,870 evaluados). `compara.py s3` sobre 27 fixtures: 906 de 906 hojas en las tres direcciones, 437 de 437 de `model.json`, 9 de 9 de `data.json`, 2,718 de 2,718 provisionales (ahí entran los casos 1 a 3 y la carga), diferencias 0. Las salidas completas de oráculo y motor en los 6 casos del comparador son idénticas (1,864 hojas). Tiempo del motor en esta corrida: caso-1 1 ms, caso-2 1 ms, caso-3 1 ms, carga 137 ms.
- **Observación (resuelta en correccion-S3-r2, ver §8):** el oráculo de Python (fracciones, día por día) tardaba unos 4.5 s en la carga y aplicaba el mismo límite de 5 s de reloj (§8.7), así que en un equipo cargado podía responder `timeout` en esa fixture. El límite de 5 s es del motor (S3-12) y no cambia.

## 8. Tiempo en el oráculo y salida guardada (correccion-S3-r2, 11-oct-2026; REVISION-S3-r2 M1 y M2)
- **Comportamiento al agotarse el tiempo (sin cambio en el motor):** el motor sigue igual: lee `comparator.timeout_seconds` (inicial 5 s, K43, configurable por K44), revisa el contexto antes de cada escenario y al terminar, y si vence devuelve `ErrTimeout` → 503 `unavailable` sin resultados parciales (§8.7 de la interfaz). El tiempo real contra los 5 s lo mide **solo** `TestRestriccionesTiempo` (líneas `S3-TIEMPO` de §6).
- **Oráculo (M2, opción (a) de la revisión):** `oraculo/cred_escenarios.py` ya no mide reloj. `timeout_seconds` = 0 responde `timeout` (determinista; fixture `s3-comparador-tiempo`) y un valor positivo calcula completo. Así la salida del oráculo no depende de qué tan cargado esté el equipo. Los milisegundos que imprime su resumen son informativos. Solicitud al integrador del contrato: dejarlo escrito en §8.7 y §10.0 de `S3-interfaz.md` (reporte `docs/estado/wf1/correccion-S3-r2.md`).
- **Prueba con el equipo cargado:** con dos procesos ocupando el mismo núcleo (`taskset -c 5`), la carga tardó 7,395 ms en el oráculo y aun así dio `ok` con 3,870 evaluados; la salida tuvo el mismo sha256 que la corrida sin carga y `compara.py s3` dio 0 diferencias.
- **Salida guardada (M1):** `oraculo/salidas/s3.json` se regeneró con `python3 -I oraculo/cred_run_fixtures.py --fixtures oraculo/fixtures/s3 --out oraculo/salidas/s3.json`: 27 fixtures, 79 consultas, segunda vía 174 de 174; sha256 `b84c0db7fa80d19dffbd9f9d1834de52b5bef34f9c8f561727e0b81f2b1bbd8f`, igual en dos corridas. `python3 -I oraculo/compara.py s3 --motor datos/tmp/motor-s3.json` con el oráculo por defecto: 906 de 906 en las tres direcciones, 437 de 437 de `model.json`, 9 de 9 de `data.json`, 2,718 de 2,718 provisionales, diferencias 0, código 0. Tiempo del motor en esta corrida: caso-1 0 ms, caso-2 0 ms, caso-3 0 ms, carga 127 ms.
