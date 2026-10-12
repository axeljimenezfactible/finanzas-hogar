# S3 · Interfaz del motor de créditos

Versión 1.5 (11-oct-2026, WF-1 y puerta 2) · Autor: agente contrato-S3; corrección r1: agente correccion-contrato-S3-r1 (hallazgos medios H1–H8 de `docs/estado/wf1/REVISION-contrato-S3-r1.md`); corrección r2: agente correccion-contrato-S3-r2 (hallazgos medios M1–M4 de `docs/estado/wf1/REVISION-contrato-S3-r2.md`); corrección r3: agente correccion-contrato-S3-r3 (hallazgos medios N1–N4 de `docs/estado/wf1/REVISION-contrato-S3-r3.md`); corrección r4: agente correccion-contrato-S3-r4 (hallazgo medio P1 de `docs/estado/wf1/REVISION-contrato-S3-r4.md`); puerta 2 (ronda P2c): agente contratos-puerta2 (marcas de Q2-28 a Q2-38) y agente correccion-contratos-r1 (Q2-67 = 1 y hallazgos A-1, M-2 y M-3 de `docs/estado/p2c/REVISION-contratos-r1.md`) · **Contrato sin código**: fija paquetes, tipos, firmas, reglas de cálculo, fixtures, CLI, oráculo y pruebas que implementan los agentes «núcleo S3» (Opus), «optimizador S3» (Opus), «regresión S3» (Sonnet) e «integrador» (Opus). Lo revisa el «revisor de motores» contra `docs/criterios/S3.md` v1.2.

Fuentes (y nada más): decisiones §C y §K (K12–K19, K31.7, K43, K44); `docs/criterios/S3.md` 1.2; `reglas-negocio.md` 1.3 (Convenciones C-01…C-16, RN-06…RN-11, RN-21, RN-25, RN-26; RN-26 en 1.4: 5 datos con `vat_included`); esquema 1.3 (`00001` `credit`, `credit_term_rate`, `promotion`, `planned_item`, `credit_purchase`, `credit_purchase_share`, `scheduled_payment`, `scheduled_payment_share`, `credit_statement`, `credit_payment`, `person`; `00002` columnas K12/K18; `00003` claves `comparator.timeout_seconds` y `optimizer.*`; `00004` columna `credit.vat_included`, Q2-67 = 1); C1 v1.2 por extractos (`creditos-a.yaml`, `creditos-b.yaml`: `compareScenarios`, `ScenarioRequest`, `ScenarioResult` y `rate_pct` de la tasa por plazo y de la compra, «ausente si falta rate_decimals»; `ajustes-d.yaml`: `getOptimizerSettings`, `updateOptimizerSettings`); `oraculo/fixtures/FORMATO.md` 1.2 (autoría de cada archivo de la carpeta §1, forma común de fixtures §2, `id` de las consultas §2.1, razones con `decimales` §3, salida y regla 5 de los mapas por persona §4, registro en el binario común §4.1, comparación §5 y regla para S3 §6); `datos/semilla/data.json` (`financiamiento.dos_quincenas`: R-05 y R-06); `oraculo/model.py` (`sched`, `mk`, eventos, sensibilidad, `ipm`) y `oraculo/model.json` (sha256 `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`, intacto); el bloqueo `docs/estado/bloqueos/U4-iva-credito.md` (de U4; RESUELTO en la puerta 2, Q2-67 = 1; §5 y §5.3).

Marcas: **FIRME** = decidido; **PROVISIONAL (bloqueo X)** = depende de un bloqueo abierto; ninguna cifra canon es provisional. **Puerta 2 (11-oct-2026, ronda P2c, agente `contratos-puerta2`; solo marcas, sin cambiar reglas aprobadas):** Axel contestó Q2-28 a Q2-38 con la opción recomendada (`docs/estado/puerta2/RESPUESTAS.md`). Quedan FIRMES la tasa del periodo del préstamo (Q2-28 = (a)), el calendario de la tarjeta y el corte que sustituye pagos (Q2-32 = A), los objetivos del comparador con las propuestas * (Q2-33 = (b), Q2-34 a Q2-36 = A), el `Source` `quoted` del cargo normal (Q2-37 = (a)) y capital = monto con intereses 0 en los tres tipos, también en la línea (Q2-38 = (a)). Sigue PROVISIONAL (bloqueo abierto `docs/estado/bloqueos/S3-formula-tarjeta-pdf.md`) la fórmula de intereses y mínimo de tarjeta y tienda. La respuesta (a) a Q2-29 a Q2-31 (preguntas 1 a 3 de `S3-formula-tarjeta`, resuelto) pide que Axel entregue un estado de cuenta de tarjeta y uno de crédito de tienda, con los que S3 escribirá la fórmula; mientras no lleguen se usa la fórmula general (b). Una decisión sin resultado no cuenta como resuelta (`INDICE.md`). **Corrección r1 de P2c (Q2-67 = 1, `U4-iva-credito` resuelto):** la entrada nueva `VATIncluded` (`credit.vat_included`, 00004) y los 5 datos de K18 de «falta capturar» (RN-26 1.4) son FIRMES (§5, §5.3); dos lecturas que Q2-67 no cubre (préstamo con `VATIncluded` null; `VATIncluded` = false sin porcentaje) quedan PROVISIONALES con el bloqueo nuevo `docs/estado/bloqueos/S3-iva-lectura.md` (§5.4, §5.8). Los nombres de los bloqueos resueltos ya no entran en ninguna lista `Provisional` (§14.1). Bloqueos que abre este contrato: `docs/estado/bloqueos/S3-formula-tarjeta.md`, `S3-convencion-prestamo.md`, `S3-objetivos-medida.md`, desde la corrección r2 `S3-origen-cargo-normal.md` y, desde la corrección r3, `S3-serie-sin-intereses.md`; en la puerta 2 (corrección r1 de P2c), `S3-formula-tarjeta-pdf.md` y `S3-iva-lectura.md`, los dos ABIERTOS. El bloqueo `docs/estado/bloqueos/U4-iva-credito.md` (lo abrió U4) quedó RESUELTO en la puerta 2 (Q2-67 = 1) y ya no entra en ninguna lista `Provisional`. Qué bloqueo entra en cada salida y en cada fixture, y en qué orden: §14.1. El bloqueo `S1-parametros-avisos` (resuelto en la puerta 2, Q2-11 = B) no toca a S3 (los días antes del pago `[1, 0]` ya están decididos en `alerts.payment_due_days_before`).

## 1. Principios de implementación (FIRMES)
1. **Dinero:** `Cents` (`int64`, centavos MXN). Ningún `float32`/`float64` en dinero ni en tasas (S3-17). Intermedios con `math/big` (`big.Int`, `big.Rat`). Única excepción permitida: convertir una fracción a `float64` **solo** para comparar con las dos hojas no monetarias de `model.json` (`parametros.tasa`, `parametros.tasa_pct`) en el código de prueba (la CLI y las vistas emiten la fracción exacta, §9.3; `compara.py` hace su propia conversión), nunca en un cálculo.
2. **Redondeo (RN-01, K14):** `RoundHalfUp(num, den)` a centavo en el último paso. Los dos modos del oráculo: (a) **montos** de una serie: base = total / n **truncada** (ROUND_DOWN) y el último ajusta; (b) **capital**, **partes externas** y **devengo**: base HALF_UP y el último ajusta; el titular toma el resto del pago. **Pago de referencia** (solo cotizaciones y sensibilidad): HALF_UP(total **sin redondear** / n). Hallazgo verificado (script de §12): calcularlo con el total ya redondeado da 0.32 en vez de 0.31 (`interes_por_monto[0].pago_quincenal`) y 1,258.69 en vez de 1,258.68 (`interes_por_monto[5]`): **prohibido**.
3. **Tasas (C-10, K12, RN-25):** fracción exacta `Num/Den` + `Decimals` capturados; la pantalla redondea HALF_UP a esos decimales; sin decimales capturados (−1) no hay texto (§5.1); el cálculo usa la fracción.
4. **Fechas (C-03, C-05):** fecha civil sin hora (`Date`), zona America/Mexico_City decidida por quien llama; el motor no lee el reloj: recibe `Today`.
5. **Sin red, sin solver, sin base de datos:** el motor es puro (entradas → salidas). Ningún `net/http`, `database/sql`, HiGHS ni OR-Tools en `motor_cred` (S3-15). La API (WF-2) lee la base, arma las entradas y guarda las series ya calculadas (C-08).
6. **Determinismo:** misma entrada → misma salida byte a byte (orden canónico en todas las listas; ver cada sección).
7. **Errores con tipo** (`errors.Is`): `ErrInvalid` (entrada inválida → API 422 `validation_failed`), `ErrNoRate` («sin tasa», S3-06), `ErrMissingCapture` («falta capturar», K18), `ErrPromotionNotApplicable`, `ErrUnsupported` (frecuencia sin regla), `ErrTimeout` (comparador → API 503 `unavailable`, C1). Nunca se devuelve un número inventado en lugar de un error. La traducción a respuestas de la API es un dato (`ErrorResponses`, §8.7), sin `net/http`.

## 2. Paquetes y archivos
| Ruta | Contenido |
|---|---|
| `backend/internal/motor_cred/` (`package motor_cred`) | núcleo: `money.go` (Cents, Rate, redondeos), `date.go` (Date y quincenas), `series.go` (RN-07), `schedule.go` (calendarios por tipo), `revolving.go` (línea tipo Kueski, tasa, promociones, partes), `availability.go` (disponible A/B), `card.go`, `loan.go`, `store.go`, `state.go` (`Input`, `CreditState`, `PurchaseState` y decodificación estricta, §5.7), `calendar.go` (calendario unificado), `cashflow.go` (efectivo por ventana, §6.1), `debt.go` (deuda y compras vivas), `quote.go` (cotización por monto y sensibilidad), `views.go` (vistas con las claves de `model.json`, §9.3), `errors.go` (errores y tabla `ErrorResponses`, §8.7) |
| `backend/internal/motor_cred/comparador/` (`package comparador`) | `request.go`, `generate.go`, `constraints.go`, `objectives.go`, `score.go`, `compare.go`, `timeout.go` (`TimeoutFromSetting`, §8.7) |
| `backend/internal/motor_cred/*_test.go`, `comparador/*_test.go` | pruebas de §11 |
| `backend/internal/motor_cred/evaluar/` (`package evaluar`) | `evaluar.go`: `Evaluate(input json.RawMessage, fn string, args json.RawMessage) (out any, code string)`, el despacho de las consultas de §9.2 y §9.3 (importa `motor_cred` y `comparador`; por eso vive en su propio paquete y no hay ciclo); `evaluar_test.go` lo prueba sin el binario (§9.5) |
| `backend/cmd/motores/s3.go` (y, si se quiere, `s3_test.go`) | **lo único de S3 en el binario común** (`FORMATO.md` 1.2 §4.1): `func init() { register("s3", evalS3) }`, con `evalS3` de tipo `Evaluator` que llama a `evaluar.Evaluate` (§9.5). `main.go` y `main_test.go` son del agente Go de S2: S3 no los edita ni escribe un sustituto (corrige r2-M1 c) |
| `oraculo/cred_lib.py`, `oraculo/cred_escenarios.py`, `oraculo/cred_run_fixtures.py`, `oraculo/cred_genera_fixtures.py` | oráculo independiente (§10.2) y generador determinista de las fixtures que no son del comparador; `cred_run_fixtures.py` escribe `oraculo/salidas/s3.json`; los cuatro los escribe el agente de regresión (§10.0) |
| `oraculo/compara.py` | comparador común de `FORMATO.md` §5; S3 lo usa como `python3 -I oraculo/compara.py s3` y no escribe uno propio |
| `oraculo/fixtures/s3/<caso>.json` y `oraculo/fixtures/s3/_datos/*.json` | fixtures (§10) con la forma de `FORMATO.md` §2; autor de cada archivo, trabajo de cada agente mientras falten y secuencia del comparador en §10.0 (`FORMATO.md` 1.2 §1) |
| `docs/estado/S3-escenarios.md` | casos del comparador fijados antes de correr, su tabla a mano y los tiempos medidos (lo escribe el optimizador, que también transcribe esa tabla a las fixtures del comparador, §10.0) |

Nombre de paquete: `motor_cred` igual al directorio (go vet lo acepta). SUPUESTO: si S2 usa otra convención para `motor_fin`, el integrador unifica (solicitud en el reporte).

## 3. Tipos base (`money.go`, `date.go`)
```go
type Cents int64
type Rate struct {
    Num      int64 `json:"num"`
    Den      int64 `json:"den"`
    Decimals int   `json:"decimales"` // -1: sin dato de captura; clave de FORMATO.md 1.2 §3
}
func NewRate(num, den int64, decimals int) (Rate, error) // Den > 0, Num >= 0
func QuotedRate(interest, principal Cents, decimals int) Rate // intereses/capital (RN-08, rate_basis term_total); decimals = *PurchaseInput.RateDecimals, o -1 si es null (§5.1)
func (r Rate) Rat() *big.Rat
func (r Rate) Display() string            // "25.8685 %" (HALF_UP a Decimals; K12); Decimals = -1 → "" (§5.1)
func ApplyRate(c Cents, r Rate) Cents     // HALF_UP(c × Num / Den)
func RoundHalfUp(num, den *big.Int) Cents // a centavo; den > 0
func SplitTrunc(total Cents, n int) []Cents     // RN-07 montos
func SplitHalfUp(total Cents, n int) []Cents    // RN-07 capital y partes externas
func ReferencePayment(totalExact *big.Rat, n int) Cents // HALF_UP(total sin redondear / n)

type Date struct { Year int; Month time.Month; Day int }
func ParseDate(s string) (Date, error) // ISO AAAA-MM-DD; rechaza 2027-02-30
func (d Date) String() string          // ISO
func (d Date) AddDays(n int) Date
func (d Date) Compare(o Date) int
func DaysBetween(a, b Date) int
func FortnightOf(d Date) Fortnight     // RN-06/K8: Q1 días 1–15 paga el 15; Q2 del 16 al último, paga el 30 o el último
type Fortnight struct { Year int; Month time.Month; Half int; PayDay Date }
```
`SplitTrunc(389500, 4)` = 97375 ×4; `SplitTrunc(503474, 4)` = 125868 ×3 + 125870; `SplitHalfUp(309450, 4)` = 77363 ×3 + 77361; `SplitHalfUp(149972, 4)` = 37493 ×4 (todas en `model.json`, §12).
```go
type Reading string // "A" | "B" (credit.availability_reading); "" se trata como "B" (RN-09, S3-03)
```
JSON de los tipos base (§9.1): `Cents` = entero de centavos; `Rate` = `{"num": …, "den": …, "decimales": …}` (la clave de `FORMATO.md` 1.2 §3 y de S2; corrige r2-M1 b); `Date` = `"AAAA-MM-DD"`; `Reading` = `"A"`, `"B"` o `""`; un `*big.Rat` de salida (`TotalExact`) = `{"num": …, "den": …}` reducido.

## 4. Personas y partes (RN-03, RN-04, RN-08, C-09; FIRME)
```go
type ShareInput struct {
    PersonID  int64  `json:"person_id,omitempty"` // la API lo llena; en fixtures se omite (0)
    Code      string `json:"person"`              // person.code (axel, karen)
    External  bool   `json:"external"`            // person.participation = external
    Principal Cents  `json:"principal_cents"`     // lo que la persona compró
    Joint     *Cents `json:"joint_cents"`         // su mitad de lo conjunto (composición informativa; null = sin dato)
    PaperIn   *Cents `json:"paper_in_cents"`      // papel comprado dentro de la compra (informativo)
}
type PersonShare struct {
    PersonID     int64        `json:"person_id,omitempty"`
    Code         string       `json:"person"`
    External     bool         `json:"external"`
    Principal    Cents        `json:"principal_cents"`
    Interest     Cents        `json:"interest_cents"`
    Total        Cents        `json:"total_cents"`
    FirstPayment Cents        `json:"first_payment_cents"`
    Composition  *Composition `json:"composition"` // null si no vino Joint ni PaperIn
}
type Composition struct {
    Joint   Cents `json:"joint_cents"`
    Own     Cents `json:"own_cents"`      // Principal − Joint
    PaperIn Cents `json:"paper_in_cents"`
    Free    Cents `json:"free_cents"`     // Principal − PaperIn
}
```
Reglas: exactamente una parte no externa (el titular; la co-titular se registra a su nombre, C-09: la API suma su parte antes de llamar); Σ `Principal` = capital de la compra; si no, `ErrInvalid`. **Intereses por persona:** externa = HALF_UP(intereses × su capital / capital) (RN-08: proporcional al capital); el titular toma el resto. **Pagos por persona:** externa = `SplitHalfUp(su total, n)`; titular = monto del pago − Σ externas (absorbe centavos; si quedara negativo, `ErrInvalid`). `FirstPayment` = su parte del pago 1 (`pago_por_quincena`). Orden de partes: titular primero, luego externas por (`PersonID`, `Code`).
Canon: despensa 1, Karen 1,191.50 → intereses HALF_UP(800.50 × 1,191.50 / 3,094.50) = 308.22; Axel 492.28; pagos Karen 374.93 ×4 y Axel 598.82 ×4 (R-07…R-13). La «Karen fija» de `parametros.karen_fija` es esta misma parte externa (capital 1,191.50 constante en las 4 compras): no es un tope (K13).

## 5. Los cuatro tipos de crédito
Entrada común (espejo de `credit`; la API la llena desde la base):
```go
type Kind string // "revolving_line" | "bank_card" | "term_loan" | "store_credit" (credit.kind)
type Credit struct {
    ID                  int64       `json:"id"`
    Name                string      `json:"name"`
    Kind                Kind        `json:"kind"`
    Holder              string      `json:"holder"`                     // person.code del titular
    HolderPersonID      int64       `json:"holder_person_id,omitempty"` // la API lo llena
    Limit               *Cents      `json:"limit_cents"`                // credit_limit_cents (obligatorio en línea y tarjeta)
    Frequency           string      `json:"frequency"`                  // "weekly" | "fortnightly" | "monthly"
    FirstPaymentMinDays *int        `json:"first_payment_min_days"`     // Kueski 15
    Reading             Reading     `json:"reading"`                    // availability_reading; "" = B (RN-09, S3-03)
    ReadingConfirmed    bool        `json:"reading_confirmed"`          // availability_confirmed
    ReadingAAnchor      *Date       `json:"reading_a_anchor"`           // reading_a_anchor_date (§5.2)
    StatementDay        *int        `json:"statement_day"`
    PaymentDueDay       *int        `json:"payment_due_day"`
    AnnualRate          *Rate       `json:"annual_rate"`                // K18
    RateBase            string      `json:"rate_base"`                  // "daily" | "monthly" | "" (K18)
    VATIncluded         *bool       `json:"vat_included"`               // credit.vat_included (00004, Q2-67 = 1): true = la tasa ya trae IVA; false = se suma VATOnInterestBP; null = falta capturar (RN-26 1.4, §5.3)
    VATOnInterestBP     *int        `json:"vat_on_interest_bp"`         // IVA que se suma (1600 = 16 %); solo con VATIncluded = false (CHECK de 00004)
    MinPaymentPctBP     *int        `json:"min_payment_pct_bp"`
    MinPaymentFixed     *Cents      `json:"min_payment_fixed_cents"`
    TermRates           []TermRate  `json:"term_rates"`
    Promotions          []Promotion `json:"promotions"`
}
type TermRate struct { // credit_term_rate
    TermCount int    `json:"term_count"`
    Rate      Rate   `json:"rate"`
    Source    string `json:"source"`     // quote | manual
    QuoteDate *Date  `json:"quote_date"`
}
type Promotion struct { // promotion
    ID          int64  `json:"id"`
    CreditID    *int64 `json:"credit_id"`
    Merchant    string `json:"merchant"`
    StoreID     *int64 `json:"store_id"`
    TermCount   int    `json:"term_count"`
    Rate        Rate   `json:"rate"`
    MinPurchase Cents  `json:"min_purchase_cents"`
    ValidFrom   Date   `json:"valid_from"`
    ValidTo     Date   `json:"valid_to"`
    Active      bool   `json:"active"`
}
type PurchaseInput struct {
    Credit           Credit       `json:"-"`                     // lo pone quien llama (en CreditState, el crédito del estado)
    Date             Date         `json:"date"`
    Principal        Cents        `json:"principal_cents"`
    TermCount        int          `json:"term_count"`
    QuotedInterest   *Cents       `json:"quoted_interest_cents"` // compra real: intereses que cotizó el acreedor (IVA incluido; interest_source quoted)
    RateDecimals     *int         `json:"rate_decimals"`         // credit_purchase.rate_decimals (00002, K12): decimales con que se cotizó la tasa; solo con QuotedInterest (§5.1)
    PromotionID      *int64       `json:"promotion_id"`          // compra con promoción elegida
    StoreID          *int64       `json:"store_id"`
    Merchant         string       `json:"merchant"`
    FirstPaymentDate *Date        `json:"first_payment_date"`    // capturada (préstamo, tienda); null = la calcula el motor
    Shares           []ShareInput `json:"shares"`
}
type Purchase struct {
    Principal        Cents              `json:"principal_cents"`
    Interest         Cents              `json:"interest_cents"`
    Total            Cents              `json:"total_cents"`
    TotalExact       *big.Rat           `json:"total_exact"`  // {"num", "den"}
    ReferencePayment Cents              `json:"reference_payment_cents"`
    Rate             Rate               `json:"rate"`
    RateDisplay      string             `json:"rate_display"` // Rate.Display(): "25.8685 %", "17.27 %" (K12)
    RateBasis        string             `json:"rate_basis"`
    Source           string             `json:"source"`       // quoted | term_table | promotion | amortization (interest_source); por tipo en §5.8
    PurchaseKind     string             `json:"purchase_kind"` // installment_plan | single_charge | loan_disbursement (credit_purchase.kind); §5.8
    PromotionID      *int64             `json:"promotion_id"`
    Payments         []ScheduledPayment `json:"payments"`
    Shares           []PersonShare      `json:"shares"`
    Provisional      []string           `json:"provisional"`  // bloqueos de los que depende esta salida; [] = FIRME; contenido y orden en §14.1
}
type ScheduledPayment struct {
    Seq       int            `json:"seq"`
    Due       Date           `json:"due_date"`
    Amount    Cents          `json:"amount_cents"`
    Principal Cents          `json:"principal_cents"`
    Interest  Cents          `json:"interest_cents"`
    Shares    []PaymentShare `json:"shares"`
}
type PaymentShare struct {
    PersonID int64  `json:"person_id,omitempty"`
    Code     string `json:"person"`
    Amount   Cents  `json:"amount_cents"`
}
func QuotePurchase(in PurchaseInput) (Purchase, error) // despacha por in.Credit.Kind (§5.8)
```
Invariantes de toda salida (las prueba `TestReparto` y la base las exige): `Total = Principal + Interest`; Σ `Amount` = `Total`; Σ `Principal` de pagos = capital; `Amount = Principal + Interest` en cada pago; Σ partes de cada pago = `Amount`; Σ pagos de cada persona = su `Total`.

### 5.1 Línea revolvente tipo Kueski (`revolving.go`, `schedule.go`; FIRME salvo la serie de capital sin intereses)
- **Intereses (RN-08, S3-06):** en este orden: (1) `QuotedInterest` → tal cual, tasa = `QuotedRate(QuotedInterest, capital, d)` con d = `RateDecimals` de la compra, o −1 si viene null (800.50 / 3,094.50 = 0.2586847632897075; con `rate_decimals` 4, «25.8685 %»); (2) `PromotionID` → debe ser aplicable (abajo) o `ErrPromotionNotApplicable`; intereses = `ApplyRate(capital, promo.Rate)` (0 % → 0.00), `Source = promotion`; (3) tasa de `TermRates` del plazo → `ApplyRate` (4,000.00 → 1,034.74), `Source = term_table`; (4) sin tasa para el plazo → `ErrNoRate` («sin tasa»; no se interpola ni se inventa para 6–12 quincenas). `TotalExact` = capital × (1 + tasa) sin redondear (para el pago de referencia).
- **Decimales de la tasa (K12, RN-25; corrige r3-N2):** `RateDecimals` es el espejo de `credit_purchase.rate_decimals` (00002 l. 123): los decimales con que el acreedor cotizó o Axel capturó la tasa de **esa** compra. Precedencia: (1) con `QuotedInterest`, `Rate.Decimals` = `RateDecimals` (null → −1); los decimales de un `TermRate` del mismo plazo **no** se usan aunque su fracción sea igual (la de d1, 80050/309450, es también la de la tabla a 4 q, pero su «25.8685 %» sale de su propio `rate_decimals` 4); (2) con promoción o `TermRate`, `Rate` es la de la promoción o la del `TermRate`, con sus `decimales`, y `RateDecimals` no se usa (la API puede mandar el que guardó). `RateDecimals` < 0 → `ErrInvalid` (CHECK `rate_decimals >= 0`, 00002). `Display()`: con `Decimals` ≥ 0, la fracción × 100 redondeada HALF_UP a `Decimals` decimales, con punto decimal y « %» («25.8685 %», «17.27 %»); con `Decimals` = −1 devuelve `""` y la compra sale con `RateDisplay` = `""`: el motor no inventa decimales y C1 omite `rate_pct` «si falta rate_decimals» (`creditos-a.yaml`, `creditos-b.yaml`). Los intereses no dependen de los decimales: con `RateDecimals` null, d1 da los mismos 80050 y `rate` = `{80050, 309450, −1}`.
- **Calendario (RN-06, `sched` del oráculo):** quincenal: candidatas = días de pago de cada quincena (15 y 30 o último día del mes); primer pago = primera candidata con `DaysBetween(compra, d) ≥ FirstPaymentMinDays` (15); los demás son las candidatas consecutivas, tantas como `TermCount` (el oráculo toma 4; la regla general genera las que hagan falta, hasta 12 quincenas). 10-oct-2026 → 30-oct, 15-nov, 30-nov, 15-dic; 10-feb-2027 → **28-feb**, 15-mar, 30-mar, 15-abr (S3-04).
- **Serie (RN-07):** montos `SplitTrunc(total, n)`, capital `SplitHalfUp(capital, n)`, interés del pago = monto − capital; partes por persona según §4. **Con intereses totales 0 y n ≥ 2** (promoción al 0 %, o cotización o tasa por plazo de 0) el **capital de cada pago es igual a su monto** en la línea, en los meses sin intereses y en la tienda (FIRME, puerta 2, Q2-38 = (a); bloqueo `S3-serie-sin-intereses` resuelto; corrige r3-N1): con 309450 a 4 q, montos y capital 77362 ×3 + 77364 e interés 0 en cada pago, que `scheduled_payment.interest_cents CHECK (interest_cents >= 0)` (00001 l. 808) admite. La 1.4 aplicaba a la línea RN-07 literal (opción b: capital 77363 ×3 + 77361 e intereses −1, −1, −1 y +3); el motor y el oráculo de la línea tienen que cambiar. Los montos, el interés total 0.00 y las partes por persona (salen de los montos, §4) son FIRMES (S3-07).
- **Promoción aplicable (RN-10, S3-07):** `Active`; `CreditID` nil o igual al crédito; comercio igual (por `StoreID` si ambos lo traen; si no, `Merchant` sin espacios a los lados y sin distinguir mayúsculas, como `COLLATE NOCASE`); `ValidFrom ≤ fecha ≤ ValidTo` (inclusive); capital ≥ `MinPurchase`; plazo = `TermCount`. Función pública `ApplicablePromotions(credit, date, principal, term, storeID, merchant) []Promotion` en orden de `ID`. El motor nunca precarga promociones: solo usa las que recibe.
- **Composición (RN-04):** si la parte trae `Joint`/`PaperIn`, el motor devuelve `Composition` (canon: Axel conjunto 851.00, propio 1,052.00, papel dentro 0.00, libre 1,903.00).

### 5.2 Disponible A/B y reposición (`availability.go`; FIRME salvo el aviso de Kueski)
```go
type CashEvent struct { // compra de contado del hogar: no mueve ninguna línea
    Date        Date   `json:"date"`
    Amount      Cents  `json:"amount_cents"`
    Person      string `json:"person"`      // quien la pagó (person.code)
    Description string `json:"description"`
}
type PaymentRecord struct { // credit_payment
    PurchaseRef   string `json:"purchase_ref"`
    Seq           int    `json:"seq"`
    Date          Date   `json:"date"`
    Payer         string `json:"payer"`                     // person.code
    PayerPersonID int64  `json:"payer_person_id,omitempty"`
    Amount        Cents  `json:"amount_cents"`
}
type TimelineInput struct { // no viaja en JSON: lo arma quien llama (consulta linea_tiempo, §9.2)
    State           CreditState // la línea con sus compras, pagos registrados y cortes (§5.7)
    Cash            []CashEvent // Input.Cash
    Today, From, To Date
}
type TimelineEvent struct {
    Date                Date   `json:"date"`
    Kind                string `json:"kind"` // start | purchase | cash | payment
    Description         string `json:"description"`
    PurchaseRef         string `json:"purchase_ref"`
    Seq                 int    `json:"seq"`
    Amount              Cents  `json:"amount_cents"`
    PrincipalRestored   Cents  `json:"principal_restored_cents"`
    InterestNotRestored Cents  `json:"interest_not_restored_cents"`
    AvailableA          Cents  `json:"available_a_cents"`
    AvailableB          Cents  `json:"available_b_cents"`
}
type Timeline struct {
    Events              []TimelineEvent `json:"events"`
    SeriesA             []Cents         `json:"series_a_cents"`
    SeriesB             []Cents         `json:"series_b_cents"`
    MinA                Cents           `json:"min_a_cents"`
    MinB                Cents           `json:"min_b_cents"`
    Reading             Reading         `json:"reading"`              // la que rige ("" → "B")
    PendingConfirmation bool            `json:"pending_confirmation"` // = !ReadingConfirmed: «pendiente confirmar con Kueski» (decisión C)
    AAnchor             Date            `json:"a_anchor"`             // ancla usada para la lectura A
    AAnchorDefaulted    bool            `json:"a_anchor_defaulted"`   // true si no había ancla y se usó From
}
func AvailabilityTimeline(in TimelineInput) (Timeline, error)
```
- Solo el capital usa y repone; los intereses nunca regresan; el contado no mueve la línea (RN-09, S3-05). Compras `cancelled` no cuentan.
- **Ancla de la lectura A (corrige r1-H5):** ancla = `ReadingAAnchor`. Sin ancla: con `Reading = "A"` → `ErrInvalid` (el esquema exige la fecha ancla con lectura A: `credit`, 00001); con B o vacío, A se calcula con ancla = `From` y `AAnchorDefaulted = true` (A es informativa mientras rige B). `From` anterior al ancla → `ErrInvalid`: la lectura A no existe antes de su ancla, y la API siempre pide `From` ≥ ancla (el tablero desde hoy; el comparador desde `Today`, §8.3).
- **Inicio en `From`:** B = límite − capital vivo a `From` (capital de las compras no canceladas con fecha ≤ `From`, menos el capital repuesto con fecha de reposición ≤ `From`). A = límite − Σ capital de las compras no canceladas con ancla < fecha ≤ `From` + Σ capital repuesto con ancla < fecha de reposición ≤ `From`: lo comprado hasta el ancla ya está dentro del límite libre y, desde el ancla, A se mueve con los mismos deltas que los eventos. La fecha de reposición es la de la regla K15 de abajo. Canon: ancla = `From` = 10-oct-2026 → A = 10,000.00, B = 6,905.50. Con el mismo canon y `From` = 16-dic-2026: A = 10,000.00 − 4,000.00 (despensa 2) + 3,094.50 (los 4 pagos de la despensa 1) = 9,094.50 y B = 10,000.00 − 4,000.00 = 6,000.00, que son el A y el B del canon tras el pago del 15-dic-2026.
- **Propiedad que se prueba** (`TestSeriesABDesdeOtraFecha`, §11): con los datos del canon y `From` = cada fecha de evento del canon posterior al ancla, el inicio vale el A y el B del canon tras el último evento con fecha ≤ `From`, y los eventos de (`From`, `To`] son idénticos a los del canon.
- **Eventos** con fecha en (`From`, `To`] (y el inicio): el inicio lleva `Amount`, `PrincipalRestored` e `InterestNotRestored` = 0 (como `credito.eventos[0]`); compra con fecha > `From` (monto = capital, `PrincipalRestored` = 0, A y B bajan el capital); contado (monto, sin mover A ni B); pago (un evento por abono registrado y uno por resto proyectado, con la regla de abajo). Inicio y contado llevan `PurchaseRef` "" y `Seq` 0; `Description` es texto informativo (en el contado, `CashEvent.Description`) y las fixtures no lo piden. A += delta sin tope; B = mín(límite, B + delta), con delta = −capital en la compra, `PrincipalRestored` en el pago y 0 en el contado. Orden: fecha, luego inicio < compra < contado < pago, luego `PurchaseRef` y `Seq`; entre los eventos de un mismo pago en la misma fecha, los abonos en su orden (abajo) y el resto proyectado al final.
- **Abonos registrados, resto proyectado y reposición (K15; corrige r2-M2).** Para cada pago programado p de una compra no cancelada (`CreditState.Quote`), con monto M, capital K e interés I = M − K, sus abonos son los `PaymentRecord` con su (`PurchaseRef`, `Seq`), en orden (`Date`, `Payer`, `Amount`), y S = Σ de sus montos.
  - **Un evento `payment` por abono**, con la fecha del abono, el `PurchaseRef` y el `Seq` del pago, `Amount` = monto del abono y `PrincipalRestored` = `InterestNotRestored` = 0, **salvo el abono que completa el pago** (el primero, en ese orden, con el acumulado ≥ M), que lleva `PrincipalRestored` = K e `InterestNotRestored` = I: el capital se repone entero cuando el pago se completa (K15) y la parte de Karen pagada por separado no repone nada antes. Los abonos posteriores al que completa (sobrepago) son eventos con 0 y 0.
  - **Resto proyectado:** si S < M (sin abonos o con abonos parciales) y el vencimiento de p ≥ `Today`, el resto se proyecta como **un** evento `payment` en el vencimiento, con `Amount` = M − S, `PrincipalRestored` = K e `InterestNotRestored` = I (lo ya abonado no cambia de fecha). Si el vencimiento es < `Today` y S < M, no hay proyección ni reposición (vencido sin completar); sus abonos siguen siendo eventos con 0 y 0.
  - **Fecha de reposición** de p = la del abono que lo completa o la del resto proyectado; es la que usa el inicio en `From` (arriba). Por pago, Σ `PrincipalRestored` es K o 0 y Σ `InterestNotRestored` es I o 0. Sin abonos registrados la regla da el canon: un evento por pago con M, K e I en su vencimiento (`credito.eventos`).
  - Ejemplo (`s3-reposicion-separada`, §10.3, con su `esperado` completo): pago 1 de la despensa 1 (M 97375, K 77363, I 20012) con el abono de Karen de 37493 el 29-oct (evento con 0 y 0) y el de Axel de 59882 el 31-oct (lo completa: 77363 y 20012); el pago 2, sin abonos, se proyecta completo el 15-nov.
- `SeriesA/B` = disponible después de cada evento de compra y de pago (cada abono y cada resto proyectado cuentan; sin inicio ni contado): 14 valores en el canon (R-65, R-66). `MinA/MinB` = mínimo de todos los eventos (8,094.50 / 5,000.00).
- Interruptor: `Reading` vacío = B y `PendingConfirmation` verdadero mientras `availability_confirmed = 0` (S3-03). La lectura A/B misma es la decisión pendiente con Kueski (decisión C): el motor siempre calcula ambas.

### 5.3 Tarjeta bancaria (`card.go`)
- **Calendario** (FIRME, puerta 2, Q2-32 = A; bloqueo `S3-formula-tarjeta` pregunta 4): corte = primera fecha ≥ compra con día `StatementDay` (recortado al fin de mes); pago = primera fecha > corte con día `PaymentDueDay` (recortado).
- **Cargo normal** (`single_charge`, 1 pago = total en la fecha de pago): intereses **0.00** si se paga el total a más tardar en la fecha de pago, con o sin datos capturados (FIRME, S3-08). Su `Source` es `quoted`, leído como «lo fija el banco en su estado de cuenta» (FIRME, puerta 2, Q2-37 = (a); §5.8).
- **Meses sin intereses** (promoción de la tarjeta con tasa 0 y plazo n): montos `SplitTrunc(total, n)` (FIRMES) en las fechas de pago de n cortes sucesivos; capital = monto (FIRME, puerta 2, Q2-38 = (a); corrige r3-N1).
- `func (c Credit) CardPeriodInterest(st StatementInput, paidCents Cents, paidDate Date) (Cents, error)` y `func (c Credit) CardMinimumPayment(st StatementInput) (Cents, error)`: si falta cualquiera de los **5** datos de K18 (`AnnualRate`, `RateBase`, `VATIncluded`, `MinPaymentPctBP`, `MinPaymentFixed`) → `ErrMissingCapture` («falta capturar»; K18, sin número). La lista de 5 datos y `VATIncluded` null = «falta capturar» son FIRMES (puerta 2, Q2-67 = 1; RN-26 1.4: `Credit.missing_terms` de C1 nombra `annual_rate`, `rate_base`, `vat`, `min_payment_pct_bp` y `min_payment_fixed_cents`; bloqueo `U4-iva-credito` resuelto; corrige r3-N3). `VATIncluded` = false sin `VATOnInterestBP` → `ErrMissingCapture` (PROVISIONAL, bloqueo `S3-iva-lectura` pregunta 2). Con datos: fórmula general PROVISIONAL (bloqueo `S3-formula-tarjeta-pdf`: Q2-29 a Q2-31 = (a), en espera de los estados de cuenta; mientras tanto, la opción (b) de `S3-formula-tarjeta`): base diaria = HALF_UP(Σ saldo diario × tasa anual / 360), base mensual = HALF_UP(saldo al corte × tasa anual / 12); IVA: se suma `VATOnInterestBP` con HALF_UP solo si `VATIncluded` = false, y con true no se suma (cuándo se suma: FIRME, Q2-67 = 1; cómo, parte de la fórmula PROVISIONAL); mínimo = máx(HALF_UP(saldo × pct / 10000), fijo), sin pasar del saldo. Toda salida calculada con la fórmula (interés o mínimo) lleva `S3-formula-tarjeta-pdf` en su lista `provisional` (§14.1); el 0.00 por pago total a tiempo es FIRME y lleva `[]`. Las dos funciones aceptan también un `store_credit` con los mismos datos K18 (§5.5, §5.8).
- `StatementInput` (de `credit_statement`; el motor puede devolver uno con `source = computed`), en el bloque siguiente.

```go
type StatementInput struct {
    StatementDate Date    `json:"statement_date"`
    DueDate       Date    `json:"due_date"`
    Balance       Cents   `json:"balance_cents"`        // pago para no generar intereses
    DailyBalances []Cents `json:"daily_balances_cents"` // base diaria; [] si no aplica
    Source        string  `json:"source"`               // bank | computed (credit_statement.source)
}
```


### 5.4 Préstamo personal a plazo fijo (`loan.go`)
`func Amortize(principal Cents, annual Rate, vatBP *int, n int, frequency string, first Date) (Purchase, error)`: tabla francesa. Tasa del periodo r = anual / 12, / 24 o / 52 y, con IVA, r × (1 + IVA) — FIRME (puerta 2, Q2-28 = (a); bloqueo `S3-convencion-prestamo` resuelto); `vatBP` es el IVA que se suma y nil = no se suma (FIRME). Qué `vatBP` recibe un `term_loan` lo decide `QuotePurchase` con `VATIncluded` (§5.8): true → nil y false → `VATOnInterestBP` (FIRMES, puerta 2, Q2-67 = 1); null → nil, que lee el dato faltante como en 1.3, «la tasa ya trae IVA» (PROVISIONAL, en uso, bloqueo `S3-iva-lectura` pregunta 1); false sin `VATOnInterestBP` → `ErrMissingCapture` (PROVISIONAL, `S3-iva-lectura` pregunta 2). La consulta `amortizar` recibe `vat_bp` directo: null = no se suma IVA (FIRME). Pago fijo P = HALF_UP(capital × r / (1 − (1 + r)^−n)) con `big.Rat` (potencia entera exacta); interés del periodo = HALF_UP(saldo × r); amortización = P − interés; **último pago** = saldo + su interés (ajusta centavos). FIRMES las propiedades de S3-09: Σ amortización = capital, saldo final 0.00, Σ pagos = capital + Σ intereses. Fechas: primera capturada; luego mismo día del mes recortado al fin de mes (regla del ancla de `planned_item.start_date`), quincenas consecutivas (RN-06) o +7 días. `TotalExact` = `Total` y `ReferencePayment` = P. `Amortize` no recibe partes: la consulta `amortizar` devuelve `Shares` y las partes de cada pago vacías (`[]`); `QuotePurchase` la llama para un `term_loan` y le aplica las partes de §4 (§5.8).

### 5.5 Crédito de tienda con abonos (`store.go`)
Abonos sin interés (tasa 0 en `TermRates`, §5.8): montos `SplitTrunc(total, n)` (el último ajusta) y saldo final 0.00 (FIRMES); capital = monto, como los meses sin intereses de §5.3 (FIRME, puerta 2, Q2-38 = (a); corrige r3-N1). Con intereses o comisión (datos K18): mismos 5 datos y mismo error que la tarjeta (K18; FIRME, puerta 2, Q2-67 = 1, §5.3) y misma fórmula PROVISIONAL por periodo (`CardPeriodInterest`, `CardMinimumPayment`; bloqueo `S3-formula-tarjeta-pdf`: Q2-31 = (a), pendiente del estado de cuenta de crédito de tienda que entrega Axel); los abonos de una compra con esa fórmula no tienen regla y `QuotePurchase` devuelve `ErrUnsupported` (§5.8). Fechas como el préstamo (primera capturada) o, si el crédito es quincenal sin fecha capturada, la regla de RN-06. Cómo se captura «sin interés», precedencia de errores y tienda con interés: §5.8.

### 5.6 Frecuencias sin regla
Línea revolvente semanal o mensual sin `FirstPaymentDate`: `ErrUnsupported` (RN-06 solo define la quincenal; ningún crédito de la casa lo necesita hoy). SUPUESTO declarado en el reporte.

### 5.7 Estado de un crédito y entradas del motor (`state.go`; FIRME; corrige r1-H2)
```go
type PurchaseState struct {
    Ref    string        `json:"ref"`    // único dentro del crédito; la API pone el id de credit_purchase
    Status string        `json:"status"` // active | paid_off | cancelled (credit_purchase.status)
    Input  PurchaseInput `json:"input"`  // lo capturado; Input.Credit lo pone el motor con el crédito del estado
}
type CreditState struct {
    Credit     Credit           `json:"credit"`
    Purchases  []PurchaseState  `json:"purchases"`
    Payments   []PaymentRecord  `json:"payments"`   // credit_payment: abonos registrados
    Statements []StatementInput `json:"statements"` // credit_statement; solo tarjeta ([] en los demás)
}
type Input struct { // entradas.datos de las fixtures (FORMATO §2); la API arma lo mismo desde la base
    Today   Date          `json:"today"`
    Credits []CreditState `json:"credits"`
    Cash    []CashEvent   `json:"cash"` // compras de contado del hogar (papel)
}
func (s CreditState) Quote() ([]Purchase, error) // QuotePurchase de cada compra, en orden (Input.Date, Ref)
func DecodeInput(r io.Reader) (Input, error)      // json.Decoder con DisallowUnknownFields: campo desconocido → ErrInvalid, código unknown_field
```
Reglas: `Ref` repetida dentro de un crédito, `PaymentRecord.PurchaseRef` sin compra, `Statements` en un crédito que no es tarjeta o `Credit.ID` repetido → `ErrInvalid`. Las compras `cancelled` se calculan (para mostrarlas) pero no cuentan en disponible, calendario, deuda ni compras vivas. «Despensa n» de las vistas (§9.3) = posición n de la compra en el orden (`Input.Date`, `Ref`) dentro de su crédito. El contado vive en `Input.Cash` y no en `CreditState` porque no pertenece a ningún crédito: la línea de tiempo lo lista como evento informativo y el efectivo por ventana (§6.1) lo suma.

### 5.8 Despacho de `QuotePurchase` por tipo de crédito (FIRME salvo lo marcado; corrige r2-M3)
`QuotePurchase` aplica el renglón de `in.Credit.Kind` (un `Kind` desconocido es `ErrInvalid`). Comunes a los cuatro: `Principal` > 0, `TermCount` ≥ 1, `Shares` válidas (§4), `RateDecimals` ≥ 0 si viene (§5.1) y, si viene, `FirstPaymentDate` ≥ `Date`; si no, `ErrInvalid`. Los errores se revisan en el orden de la última columna y se devuelve el primero que aplica; nunca un número en su lugar.

| `Kind` | Entradas obligatorias | Intereses (en orden) | Serie y fechas | `Source` · `PurchaseKind` · `RateBasis` · `Rate` | Errores, en orden de precedencia |
|---|---|---|---|---|---|
| `revolving_line` | `Limit`, `Frequency`; quincenal sin `FirstPaymentDate`: `FirstPaymentMinDays` | §5.1: (1) `QuotedInterest`; (2) `PromotionID`; (3) `TermRate` del plazo; (4) error | §5.1 (RN-06, RN-07) y §5.6 | `quoted`, `promotion` o `term_table` · `installment_plan` · `term_total` · la de la regla aplicada | `ErrInvalid` (entradas) → `ErrUnsupported` (§5.6) → `ErrPromotionNotApplicable` → `ErrNoRate`. Nunca `ErrMissingCapture`: la línea no usa datos K18 |
| `bank_card` | `Limit`, `StatementDay`, `PaymentDueDay`; `QuotedInterest` null | **Cargo normal** (`TermCount` 1 sin `PromotionID`): 0.00, porque se paga el total a más tardar en la fecha de pago (S3-08). **Meses sin intereses** (`PromotionID`): la promoción debe ser aplicable (§5.1) y de tasa 0 | Cargo normal: un pago = total = capital en la fecha de pago de su corte; MSI: `SplitTrunc(total, n)` en las fechas de pago de n cortes sucesivos, capital = monto (FIRME, puerta 2, Q2-38 = (a); §5.3). Fechas FIRMES (puerta 2, Q2-32 = A; §5.3) | Cargo normal: `quoted` (FIRME, puerta 2, Q2-37 = (a)) · `single_charge` · `""` (sin tasa: la API guarda la tasa de la compra en NULL) · `{0, 1, 0}`; MSI: `promotion` · `installment_plan` · `term_total` · la de la promoción | `ErrInvalid` (entradas, o `QuotedInterest` presente) → `ErrPromotionNotApplicable` → `ErrUnsupported` (promoción de tarjeta con tasa ≠ 0: los meses con intereses no tienen regla) → `ErrNoRate` (`TermCount` > 1 sin promoción). Nunca `ErrMissingCapture`: los datos K18 solo los usan `CardPeriodInterest` y `CardMinimumPayment` (sus 5 datos, FIRMES por Q2-67 = 1, §5.3) |
| `term_loan` | `AnnualRate`, `Frequency`, `FirstPaymentDate`; `QuotedInterest` y `PromotionID` null; `Limit` opcional | `Amortize(Principal, AnnualRate, vatBP, TermCount, Frequency, FirstPaymentDate)` con `vatBP` según `VATIncluded` (§5.4: true → nil y false → `VATOnInterestBP`, FIRMES por Q2-67 = 1; null → nil, PROVISIONAL `S3-iva-lectura` pregunta 1; tasa del periodo FIRME, Q2-28 = (a)) | La de `Amortize`, con las partes de §4 aplicadas a cada pago | `amortization` · `loan_disbursement` · `annual` · `AnnualRate` | `ErrInvalid` (entradas, sin `FirstPaymentDate`, o con `QuotedInterest` o `PromotionID`) → `ErrMissingCapture` (sin `AnnualRate`: la tasa del préstamo es captura del producto y no tiene valor por defecto; o `VATIncluded` = false sin `VATOnInterestBP`: PROVISIONAL, `S3-iva-lectura` pregunta 2) |
| `store_credit` | `Frequency`; `FirstPaymentDate`, salvo quincenal con `FirstPaymentMinDays` (regla de RN-06, §5.5); `Limit` opcional | (1) `QuotedInterest` (compra real, RN-08); (2) `PromotionID` (§5.1); (3) `TermRate` del plazo (RN-08, compra planeada): **tasa 0 = abonos sin interés**, tasa > 0 como en la línea; (4) error | §5.5: montos `SplitTrunc(total, n)`; capital por pago de RN-07 y, con intereses 0, capital = monto (FIRME, Q2-38 = (a), §5.5) | `quoted`, `promotion` o `term_table` · `installment_plan` · `term_total` · la de la regla aplicada | `ErrInvalid` (entradas) → `ErrPromotionNotApplicable` → en (4): si falta cualquiera de los 5 datos K18 (`AnnualRate`, `RateBase`, `VATIncluded`, `MinPaymentPctBP`, `MinPaymentFixed`) → `ErrMissingCapture` («falta capturar», S3-10; FIRME, Q2-67 = 1; con `VATIncluded` = false sin `VATOnInterestBP`, PROVISIONAL `S3-iva-lectura` pregunta 2); con los cinco → `ErrUnsupported` PROVISIONAL (`S3-formula-tarjeta-pdf`, Q2-31 = (a)). Nunca `ErrNoRate` |

- **Tienda sin interés frente a «sin datos» (S3-10):** «abonos sin interés» se captura como un renglón de `credit_term_rate` del crédito de tienda con su plazo y tasa 0 (`rate_num` 0, `rate_den` 1, `source` `manual`), que el esquema ya admite (`rate_num ≥ 0`, 00001); no hace falta columna nueva. «Sin datos» es no tener `TermRate` del plazo, ni promoción, ni cotización, ni los datos K18: `missing_capture`. Así un campo distingue los dos casos y la fixture `s3-tienda-abonos` trae los dos.
- **Tienda con interés (datos K18):** la fórmula de §5.3 da el interés de un periodo de corte, no los abonos de una compra, y ninguna fuente dice cómo se arman esos abonos (`S3-formula-tarjeta` pregunta 3: Q2-31 = (a), en espera del estado de cuenta; bloqueo abierto `S3-formula-tarjeta-pdf`). Mientras tanto, `QuotePurchase` de esa compra devuelve `ErrUnsupported` (PROVISIONAL, `S3-formula-tarjeta-pdf`) y nunca un número. Los intereses y el mínimo de un periodo de ese crédito sí se calculan con `CardPeriodInterest` y `CardMinimumPayment` (consultas `interes_tarjeta` y `minimo_tarjeta` con su `credit_id`; mismos datos K18, mismo `ErrMissingCapture` y misma fórmula PROVISIONAL): es el algoritmo común de motor y oráculo para S3-10 (fixture `s3-tienda-abonos-interes`). Una compra de tienda con intereses cotizados (1) o con tasa por plazo (3) no necesita la fórmula. Un crédito con una compra `unsupported` no entra al calendario unificado (`CreditState.Quote` devuelve el error).
- **Préstamo dentro de `QuotePurchase` (S3-09, S3-11):** se llama `Amortize` y a su serie se le aplican las partes de §4 (intereses por persona proporcionales al capital; cada externa `SplitHalfUp(su total, n)` y el titular el resto), así el calendario unificado trae la persona responsable de cada pago. La consulta `amortizar` llama `Amortize` sola y devuelve las partes vacías.
- **`Source` del cargo normal:** el esquema exige `interest_source` y ninguno de sus cuatro valores describe un cargo normal; el motor emite `quoted`, leído como «lo fija el banco en su estado de cuenta», sin cambio de esquema (FIRME, puerta 2, Q2-37 = (a); bloqueo `S3-origen-cargo-normal` resuelto). Los 0.00 del cargo normal son FIRMES.
- **`TotalExact` y `ReferencePayment`:** con `QuotedInterest`, `TotalExact` = capital + intereses cotizados; con promoción o `TermRate`, capital × (1 + tasa) sin redondear (§5.1); en la tarjeta, `Total`; en el préstamo, §5.4. `ReferencePayment` = HALF_UP(`TotalExact` / n), salvo en el préstamo (P).

## 6. Calendario unificado de pagos (`calendar.go`; FIRME salvo la tarjeta con cortes; lo consumen S2 y S7)
```go
type CalendarEntry struct {
    Date          Date           `json:"date"`
    CreditID      int64          `json:"credit_id"`
    CreditKind    Kind           `json:"credit_kind"`
    Source        string         `json:"source"`         // scheduled | statement
    PurchaseRef   string         `json:"purchase_ref"`   // "" en un corte
    Seq           int            `json:"seq"`            // 0 en un corte
    StatementDate *Date          `json:"statement_date"` // solo en un corte
    Amount        Cents          `json:"amount_cents"`
    Principal     Cents          `json:"principal_cents"`
    Interest      Cents          `json:"interest_cents"`
    Shares        []PaymentShare `json:"shares"`         // persona responsable de cada parte (Karen paga directo la suya)
    Status        string         `json:"status"`         // pending | partial | paid (según PaymentRecord)
    Provisional   []string       `json:"provisional"` // §14.1: la de su compra; en un corte computed, S3-formula-tarjeta-pdf
}
type FortnightTotal struct {
    PayDay   Date           `json:"pay_day"`
    Total    Cents          `json:"total_cents"`
    ByPerson []PaymentShare `json:"by_person"` // orden por person
}
type Calendar struct {
    Entries     []CalendarEntry  `json:"entries"`
    ByFortnight []FortnightTotal `json:"by_fortnight"`
}
func UnifiedCalendar(credits []CreditState, from, to Date) (Calendar, error)
```
- **Qué entra:** de cada crédito, los pagos programados (`Source = scheduled`) de sus compras no canceladas (`CreditState.Quote`) con fecha en [from, to] y, solo en tarjeta, sus cortes (`Source = statement`) con `DueDate` en [from, to].
- **Tarjeta con cortes: el corte sustituye, no se suma (corrige r1-H7; FIRME desde la puerta 2, porque el calendario de la tarjeta es FIRME por Q2-32 = A).** Un corte de una tarjeta (`StatementInput`, de banco o calculado) **sustituye** a los pagos programados de esa tarjeta cuyo vencimiento es su `DueDate` (cargo normal y mensualidades de meses sin intereses que caen en ese pago): esos pagos no entran y entra una sola entrada del corte con `Date = DueDate`, `Amount = Balance`, `Principal` = Σ capital de los pagos sustituidos e `Interest = Amount − Principal`; partes: cada persona externa con Σ de sus partes en los pagos sustituidos y el titular (`Holder`) con el resto. Un `Interest` o una parte del titular negativos → `ErrInvalid` (el corte no puede ser menor que lo que incluye). Sin corte para un vencimiento, entran los pagos programados. Así el mismo dinero nunca aparece dos veces. Las entradas de un corte `computed` (saldo calculado con la fórmula de §5.3) llevan `Provisional = ["S3-formula-tarjeta-pdf"]`; las de un corte `bank` llevan `[]`.
- **Estado:** `Status` por (`PurchaseRef`, `Seq`) con Σ `PaymentRecord`: 0 → pending, menor que `Amount` → partial, ≥ → paid; en un corte, con Σ de los registros de los pagos que sustituye.
- **Orden:** `Date`, `CreditID`, `Source` (scheduled antes que statement), `PurchaseRef`, `Seq`. **Sin duplicados:** llave (`CreditID`, `Source`, `PurchaseRef`, `Seq`, `StatementDate`); una repetida (la misma compra o el mismo pago dos veces, dos cortes con la misma fecha) es `ErrInvalid`.
- **Por quincena:** `ByFortnight` agrupa por `FortnightOf(Date).PayDay` (RN-06) las quincenas con alguna entrada; Σ `ByFortnight.Total` = Σ `Amount` y en cada quincena Σ `ByPerson` = `Total` (S3-11). S2 lo usa para el efectivo por quincena y por mes; S7 para «pago próximo» con `alerts.payment_due_days_before` = `[1, 0]` (decisión F, K44). El objetivo 2 del comparador (§8.4) usa además las quincenas sin entradas, con total 0.

### 6.1 Efectivo por ventana (`cashflow.go`; FIRME; corrige r1-H3)
```go
type CashWindow struct {
    From         Date           `json:"from"`
    To           Date           `json:"to"`
    Payments     Cents          `json:"payments_cents"` // Σ Amount del calendario con fecha en [From, To]
    Cash         Cents          `json:"cash_cents"`     // Σ CashEvent.Amount con fecha en [From, To]
    Total        Cents          `json:"total_cents"`    // Payments + Cash
    PaymentCount int            `json:"payment_count"`
    ByPerson     []PaymentShare `json:"by_person"`      // partes de los pagos más cada contado a su Person; orden por person
}
func CashFlow(cal Calendar, cash []CashEvent, from, to Date) CashWindow // cal = UnifiedCalendar(…, from, to)
```
Es el insumo de efectivo que S3 entrega a S2 (S2 lo recalcula y lo confirma, S2-01). Canon con las ventanas del 10 al 9 del oráculo (`periodos[k].inicio`…`fin`): P1 973.75, P2 1,947.50, P3 2,232.43, P4 2,716.36 (con el papel de contado del 10-ene-2027), P5 2,517.38 y P6 2,716.36 (con el del 10-mar-2027) = `periodos[k].efectivo.total` (R-45…R-50); en [10-oct-2026, 10-abr-2027], 13,103.78 con 398.00 de contado (R-52, R-60).

## 7. Deuda viva, compras vivas, cotización por monto y sensibilidad
```go
type DebtLine struct {
    CreditID     int64  `json:"credit_id"`
    PurchaseRef  string `json:"purchase_ref"`  // "" en las líneas por crédito
    Pending      Cents  `json:"pending_cents"` // Σ montos con vencimiento posterior a AsOf
    PaymentCount int    `json:"payment_count"` // pagos pendientes
}
type DebtSummary struct {
    AsOf       Date       `json:"as_of"`
    Total      Cents      `json:"total_cents"`
    ByPurchase []DebtLine `json:"by_purchase"` // solo Pending > 0; orden CreditID y (Input.Date, Ref)
    ByCredit   []DebtLine `json:"by_credit"`   // solo Pending > 0; orden CreditID
}
func LiveDebt(credits []CreditState, asOf Date) (DebtSummary, error) // RN-21
func LivePurchases(state CreditState, at Date) ([]string, error)      // Refs vivas en esa fecha
type AmountQuote struct {
    Amount           Cents `json:"amount_cents"`
    Interest         Cents `json:"interest_cents"`
    Total            Cents `json:"total_cents"`
    ReferencePayment Cents `json:"reference_payment_cents"`
}
func QuoteByAmount(amount Cents, r Rate, n int) AmountQuote // interes_por_monto
type SensitivityInput struct {
    Amounts                []Cents    `json:"amounts_cents"`
    Rate                   Rate       `json:"rate"`
    TermCount              int        `json:"term_count"`
    External               ShareInput `json:"external"`                       // Karen 1,191.50
    HolderCurrentPrincipal Cents      `json:"holder_current_principal_cents"` // 1,903.00
    AccrualMonths          int        `json:"accrual_months"`                 // 2: despensa bimestral
    CashMonthly            Cents      `json:"cash_monthly_cents"`             // 99.50 = 199.00 / 2 (dato de S2, valor estable)
}
type SensitivityRow struct {
    Amount             Cents `json:"amount_cents"`            // N
    Interest           Cents `json:"interest_cents"`          // intereses
    Total              Cents `json:"total_cents"`             // total
    ReferencePayment   Cents `json:"reference_payment_cents"` // pago_quincenal
    HolderTotal        Cents `json:"holder_total_cents"`      // axel_total
    HolderPrincipal    Cents `json:"holder_principal_cents"`  // axel_capital
    HolderBelowCurrent bool  `json:"holder_below_current"`    // axel_menor_que_actual
    MonthlyPurchase    Cents `json:"monthly_purchase_cents"`  // HALF_UP(t / AccrualMonths), sin contado (insumo de R-43)
    MonthlyTotal       Cents `json:"monthly_total_cents"`     // costo_mensual_total
    HolderMonthly      Cents `json:"holder_monthly_cents"`    // axel_mensual
    ExternalMonthly    Cents `json:"external_monthly_cents"`  // karen_mensual
}
func Sensitivity(in SensitivityInput) ([]SensitivityRow, error) // sensibilidad
```
- **Deuda viva (RN-21, K19; FIRME):** Σ montos de pagos programados con vencimiento **posterior** a `asOf` de compras no canceladas. Al 10-abr-2027: 6,293.44 = 1,258.70 (despensa 3) + 5,034.74 (despensa 4) (cifras de S2: `horizonte.deuda_viva_*`; S3 las produce y S2 las confirma con la identidad de K19).
- **Compra viva (RN-11, RN-21; FIRME):** desde su fecha hasta su último pago, **inclusive** (lectura literal de «desde… hasta» y del ejemplo de RN-11: con máximo 1, la despensa 2 el 10-dic-2026 viola y el 16-dic-2026 cumple).
- **Cotización por monto:** `Interest = ApplyRate(amount, r)`, `Total = amount + Interest`, `ReferencePayment = HALF_UP(amount × (1 + r) / n)` desde el total exacto (§1.2). Las 7 filas de `interes_por_monto` (montos 1.00 … 5,000.00, tasa de 4 quincenas, n = 4).
- **Sensibilidad:** por monto N, con t = N × (1 + tasa) exacto y KT = total de la parte externa (HALF_UP de su capital × tasa más su capital = 1,499.72): `Interest` = ApplyRate(N); `Total` = N + intereses; `ReferencePayment` = HALF_UP(t / n); `HolderTotal` = total − KT; `HolderPrincipal` = N − capital externo; `HolderBelowCurrent` = `HolderPrincipal` < `HolderCurrentPrincipal` (booleano, no es hoja numérica); `MonthlyPurchase` = HALF_UP(t / AccrualMonths); `MonthlyTotal` = `MonthlyPurchase` + `CashMonthly`; `HolderMonthly` = HALF_UP((t − KT) / AccrualMonths); `ExternalMonthly` = HALF_UP(KT / AccrualMonths) + `CashMonthly`. Verificado contra las 45 hojas (§12). Nota: `HolderMonthly` + `ExternalMonthly` puede diferir en 0.01 de `MonthlyTotal`; el canon no la tiene y no se fuerza.

## 8. Comparador (`motor_cred/comparador`)
```go
type Request struct {
    Today       mc.Date          `json:"today"`    // en fixtures sale de Input.Today
    Amount      mc.Cents         `json:"amount_cents"`
    DesiredDate mc.Date          `json:"desired_date"`
    StoreID     *int64           `json:"store_id"`
    Merchant    string           `json:"merchant_name"`
    Shares      []mc.ShareInput  `json:"shares"`
    Global      Constraints      `json:"global"`   // Ajustes › Créditos y condiciones (optimizer.*)
    Override    Constraints      `json:"override"` // las de la petición (C1 ScenarioRequest); K16
    Weights     Weights          `json:"weights"`  // optimizer.weights; enteros ≥ 0
    Credits     []mc.CreditState `json:"credits"`  // créditos activos (§5.7); en fixtures sale de Input.Credits
    Caps        []PersonCap      `json:"caps"`     // person.fortnight_payment_cap_cents (K13; puede no haber)
    Baskets     []BasketQuote    `json:"baskets"`  // de S5 (K17), precios vigentes; FIRME, Q2-33 = (b)
}
type Constraints struct { // null = sin capturar
    MinAvailableBuffer *mc.Cents `json:"min_available_buffer_cents"`
    MaxLivePurchases   *int      `json:"max_live_purchases"`
    ZeroDebtBy         *mc.Date  `json:"zero_debt_by_date"`
    DateWindowDays     *int      `json:"date_window_days"`
}
type Weights struct {
    LowestInterest   int64 `json:"lowest_interest"`
    EvenLowPayment   int64 `json:"even_low_payment"`
    MaxAvailable     int64 `json:"max_available"`
    MostValuePerPeso int64 `json:"most_value_per_peso"`
}
type PersonCap struct {
    Person   string   `json:"person"`
    PersonID int64    `json:"person_id,omitempty"`
    Cap      mc.Cents `json:"cap_cents"`
}
type BasketQuote struct { // U y C de la canasta en esa tienda
    StoreID        int64    `json:"store_id"`
    Merchant       string   `json:"merchant_name"`
    BaseUnitsMilli int64    `json:"base_units_milli"`
    Cost           mc.Cents `json:"cost_cents"`
}
type Scenario struct {
    CreditID               int64        `json:"credit_id"`
    TermCount              int          `json:"term_count"`
    PurchaseDate           mc.Date      `json:"purchase_date"`
    StoreID                *int64       `json:"store_id"`
    PromotionID            *int64       `json:"promotion_id"`
    Interest               mc.Cents     `json:"interest_cents"`
    MaxFortnightPayment    mc.Cents     `json:"max_fortnight_payment_cents"`
    FortnightPaymentSpread mc.Cents     `json:"fortnight_payment_spread_cents"`
    StdDevMilliCents       int64        `json:"std_dev_milli_cents"`
    MinAvailable           mc.Cents     `json:"min_available_cents"`
    MaxLivePurchases       int          `json:"max_live_purchases"`
    ZeroDebtDate           mc.Date      `json:"zero_debt_date"`
    ValuePerPesoMilli      *int64       `json:"value_per_peso_milli"` // null sin BasketQuote de la tienda
    ScoreBP                *int64       `json:"score_bp"`             // null en los descartados
    OverCap                []PersonPeak `json:"over_cap"`
    Provisional            []string     `json:"provisional"` // §14.1: los de su compra y los de la serie del objetivo 2
}
type PersonPeak struct {
    Person string   `json:"person"`
    Peak   mc.Cents `json:"peak_cents"` // su pago máximo por quincena en la serie del objetivo 2
    Cap    mc.Cents `json:"cap_cents"`
}
type Rejected struct {
    Scenario           Scenario `json:"scenario"`
    ViolatedConstraint string   `json:"violated_constraint"` // min_available_buffer | max_live_purchases | zero_debt_by_date (C1)
}
type ExcludedCredit struct {
    CreditID int64  `json:"credit_id"`
    Reason   string `json:"reason"` // no_limit (FIRME, Q2-36 = A)
}
type Optima struct {
    LowestInterest   Scenario `json:"lowest_interest"`
    EvenLowPayment   Scenario `json:"even_low_payment"`
    MaxAvailable     Scenario `json:"max_available"`
    MostValuePerPeso Scenario `json:"most_value_per_peso"`
}
type Result struct { // espejo de C1 ScenarioResult, más excluded y provisional
    Status             string           `json:"status"`              // ok | missing_constraints | no_valid_scenario
    MissingConstraints []string         `json:"missing_constraints"` // min_available_buffer_cents, max_live_purchases, zero_debt_by_date, date_window_days (orden de C1)
    ConstraintsUsed    Constraints      `json:"constraints_used"`
    EvaluatedCount     int              `json:"evaluated_count"` // candidatos enumerados = len(Valid) + len(Rejected)
    Valid              []Scenario       `json:"valid"`
    Optima             *Optima          `json:"optima"`         // null si Status ≠ ok
    Recommendation     *Scenario        `json:"recommendation"` // null si Status ≠ ok
    Rejected           []Rejected       `json:"rejected"`
    Excluded           []ExcludedCredit `json:"excluded"`
    Provisional        []string         `json:"provisional"` // §14.1: unión de las de sus escenarios; [] con missing_constraints
}
func Compare(ctx context.Context, req Request) (Result, error)
```
`Scenario` lleva los campos de C1 `Scenario` más `store_id`, `std_dev_milli_cents`, `over_cap` y `provisional`, y `Result` lleva además `excluded` y `provisional` (solicitud a C1 en el reporte del contrato). Los descartados llevan sus valores calculados y `score_bp` = null. Las listas `provisional` de `Scenario` y `Result` siguen §14.1 (corrige r2-M4).

**8.1 Restricciones (K16; FIRME).** `constraints_used` = Override si viene, si no Global, campo por campo. Si falta cualquiera de las 4 (incluida `DateWindowDays`, como dice C1): `status = missing_constraints`, lista de faltantes en el orden de C1, `evaluated_count = 0`, listas vacías, sin óptimos ni recomendación (pantalla: «faltan restricciones»). No hay valores por defecto. `DateWindowDays` tiene propuesta para Axel (reporte), no valor.
**8.2 Generación (búsqueda exhaustiva, FIRME salvo lo marcado).** Créditos: los de `Credits` con límite (FIRME, puerta 2, Q2-36 = A; sin límite → `excluded` con `reason = no_limit`). Fechas: cada día de `DesiredDate − w` a `DesiredDate + w`, sin fechas anteriores a `Today`. Tienda: la de la petición si trae `StoreID` o `Merchant`; si no, cada `BasketQuote` (FIRME, puerta 2, Q2-33 = (b)). Plazos y promociones por tipo (§5.8; corrige r2-M3 b): línea y tienda = plazos de `TermRates` ∪ plazos de promociones aplicables y, en cada plazo, «sin promoción» si el plazo tiene `TermRate` (cuenta la tasa 0 de los abonos sin interés de una tienda) más cada promoción aplicable (§5.1); **tarjeta: el cargo normal (plazo 1, sin promoción, 0.00) siempre se genera**, en cada fecha y tienda, más cada promoción de meses sin intereses aplicable en su plazo. Un plazo de línea o tienda sin `TermRate` ni promoción **no se genera** (S3-06), y una tienda que solo tiene datos K18 no aporta plazos (su cotización es `unsupported` o `missing_capture`, §5.8). Cada escenario se calcula con `QuotePurchase` del núcleo (ninguna fórmula duplicada); un candidato cuya `QuotePurchase` devuelve error no se genera ni cuenta en `evaluated_count`.
**8.3 Evaluación de restricciones (FIRME; la ventana, por Q2-35 = A).** Ventana = vida de la compra del escenario [fecha, último pago]. Disponible de la línea usada: `AvailabilityTimeline` (§5.2) con `From = Today` (≥ ancla), `To` = último pago del escenario, las compras existentes más la del escenario y la lectura del crédito (`Reading`, vacío = B); el disponible de cada día de la ventana es el del último evento (o el del inicio) con fecha ≤ ese día, así que el valor con que abre la ventana cuenta y, con `Today` posterior al ancla, A ya no vuelve al límite (corrige r1-H5). `min_available_cents` = mínimo de la ventana. Colchón: `min_available_cents` ≥ colchón. Compras vivas: máximo de compras vivas simultáneas del crédito usado en la ventana, contando la del escenario, ≤ máximo. Deuda en cero: último vencimiento del crédito usado con la compra añadida ≤ `ZeroDebtBy` (pagos ese mismo día cumplen). Si viola varias, `violated_constraint` = la primera en el orden de C1 (colchón, vivas, deuda). El tope por persona **no** descarta: marca `OverCap` (S3-12).
**8.4 Objetivos (K17).** (1) `interest_cents`, menor es mejor (FIRME). (2) **Pago parejo (FIRME, Q2-34 = A; serie fijada en r1-H6):** quincenas = todas las quincenas consecutivas de `FortnightOf(fecha de compra del escenario)` a `FortnightOf(fecha de su último pago)`, inclusive, **incluidas las de total 0**; total de cada quincena = Σ `Amount` de las entradas del calendario unificado (§6: compras existentes no canceladas de todos los créditos de `Credits` más el escenario) cuya `FortnightOf(Date)` es esa quincena; es la serie del hogar (todas las personas). Con n quincenas y totales xᵢ en centavos: N = 10⁶ × (n·Σxᵢ² − (Σxᵢ)²), entero ≥ 0, y `StdDevMilliCents` = ⌊(⌊isqrt(4N) / n⌋ + 1) / 2⌋, que es HALF_UP(1000 × desviación estándar poblacional) exacto (Go con `big.Int.Sqrt`; el oráculo con `math.isqrt`). Ejemplo: totales 0, 100 y 200 centavos → N = 6 × 10¹⁰, isqrt(4N) = 489,897, ⌊489,897 / 3⌋ = 163,299 y `StdDevMilliCents` = 81,650 (81.650 centavos); se muestra en centavos HALF_UP. `max_fortnight_payment_cents` = máx xᵢ y `fortnight_payment_spread_cents` = máx − mín de la misma serie (con las de 0). `OverCap`: para cada persona con `PersonCap`, su pico = máximo por quincena de sus partes (`ByPerson`) en esas mismas quincenas; si pasa su tope, entra con `Peak` y `Cap` (orden por person). Menor `StdDevMilliCents` es mejor. (3) `min_available_cents` en la ventana, mayor es mejor. (4) `value_per_peso_milli` = HALF_UP(monto × U × 100 / (C × (monto + intereses))) — FIRME, Q2-33 = (b); sin `BasketQuote` de esa tienda el valor es nulo y el objetivo no distingue; mayor es mejor.
**8.5 Normalización y puntaje (K17, RN-11; FIRME).** **Conjunto (corrige r1-H4): los escenarios generados en el sentido de la decisión C**, que dice «genera escenarios (crédito × plazo × fecha ± días × promoción) que respetan las restricciones duras»: los generados son los que respetan las 3 restricciones, es decir, `valid`. Los descartados (`rejected`) son candidatos que el motor enumera para poder decir qué restricción violan (cuentan en `evaluated_count`), pero no son escenarios generados y no entran al mínimo ni al máximo: un candidato que no se puede elegir no mueve la escala de los que sí. Solicitud a RN-11 y a C1 (`Scenario.score_bp`) para que digan «entre los escenarios generados (los que respetan las restricciones duras, decisión C)». Menor-es-mejor (máx − v)/(máx − mín); mayor-es-mejor (v − mín)/(máx − mín), en `big.Rat`. **Empate (máx = mín) o todos nulos: el objetivo vale 1 en todos los escenarios** (declaración de S3-14: el orden no cambia con cualquier constante, porque suma lo mismo a todos; se elige 1 porque todos son el mejor valor). Puntaje = Σ wᵢ·nᵢ / Σ wᵢ (pesos 1,1,1,1 → 0.25 c/u; Σ wᵢ = 0 → `ErrInvalid`); `score_bp` = HALF_UP(10000 × puntaje). Ejemplo de S3-14: intereses 100, 200, 400 → 1, 2/3, 0.
**8.6 Óptimos, recomendación y desempates (FIRME).** Óptimo de cada objetivo = mejor valor; empate → (solo en el objetivo 2; FIRME, puerta 2, Q2-34 = A) menos personas en `OverCap`, como dice el bloqueo `S3-objetivos-medida` (2): el pico contra el tope «solo marca y desempata» → mayor puntaje exacto → llave canónica menor (`credit_id`, `term_count`, `purchase_date`, `store_id` nil primero, `promotion_id` nil primero). Recomendación = mayor puntaje exacto; empate → llave canónica (el tope no desempata la recomendación: el bloqueo lo pone en el objetivo 2). Listas `valid` y `rejected` en orden de llave canónica. Ningún válido → `no_valid_scenario`.
**8.7 Tiempo máximo (K43, K44; FIRME; corrige r1-H8).**
```go
// comparador/timeout.go
func TimeoutFromSetting(seconds int64) (time.Duration, error) // comparator.timeout_seconds → duración; < 0 → ErrInvalid; 0 = vence de inmediato
// motor_cred/errors.go: la traducción a la respuesta de C1 es un dato, sin net/http (S3-15)
type ErrorResponse struct { Err error; Status int; Code string }
var ErrorResponses = []ErrorResponse{
    {ErrInvalid, 422, "validation_failed"},
    {ErrTimeout, 503, "unavailable"},
}
func ResponseFor(err error) (status int, code string, ok bool) // primer renglón con errors.Is; ok = false si no está en la tabla
```
La API (WF-2) lee `comparator.timeout_seconds` de `app_setting` (inicial 5, 00003), llama `TimeoutFromSetting`, crea `context.WithTimeout` y traduce el error con `ResponseFor`; la CLI (§9.5) y las pruebas hacen lo mismo con el argumento `timeout_seconds` de la consulta `comparar`. `Compare` revisa `ctx.Err()` antes de cada escenario y al terminar; si venció devuelve `ErrTimeout` sin resultados parciales. La tabla solo trae las dos respuestas que C1 fija para el motor; los demás errores los traduce la API según C1 (WF-2). Medición: `TestRestriccionesTiempo` corre los 3 casos y `s3-comparador-carga` con 5 s y escribe con `t.Logf` una línea `S3-TIEMPO <caso> <ms> ms` por caso; el optimizador copia esas líneas a `docs/estado/S3-escenarios.md` y, si algún caso pasa de 5,000 ms, abre `docs/estado/bloqueos/S3-tiempo-optimizador.md` y no cambia el valor. Qué parte de S3-12 se prueba aquí y cuál en WF-2: §11. Estimación: ≤ 6 créditos × 7 plazos × 15 fechas × 6 tiendas × 2 promociones ≈ 7,600 escenarios de ≤ 12 pagos, muy por debajo de 5 s (lo mide el caso de carga).

## 9. JSON, consultas, vistas, registro y CLI `s3 dump` (`oraculo/fixtures/FORMATO.md` 1.2; corrige r1-H1, r2-M1 y r3-N4)
**9.1 Valores** (mismas convenciones que `S2-interfaz.md` §1.1 y `FORMATO.md` §3): dinero = entero de centavos con sufijo `_cents` en las salidas genéricas; las vistas (§9.3) usan las claves literales de `model.json`, sin sufijo; fecha `"AAAA-MM-DD"`; persona = `code` de `person` (en Go los `PersonID` son opcionales y valen 0 en fixtures); tasa (entrada y salidas genéricas) = `{"num", "den", "decimales"}` (`FORMATO.md` 1.2 §3); razón de salida en las vistas de `model.json` = `{"num", "den"}` exacta; en la vista de `data.json` (`vista_dos_quincenas`) = `{"num", "den", "decimales"}`. Las llaves son las etiquetas `json` de §3–§8; un puntero ausente vale null y una lista ausente vale `[]`. La decodificación es estricta (`DecodeInput`, §5.7): campo desconocido → `unknown_field`.
**Códigos de error en las salidas** (`{"error": "<código>"}`, `FORMATO.md` §2 y §4): `ErrInvalid` → `invalid` (`unknown_field` si fue un campo desconocido), `ErrNoRate` → `no_rate`, `ErrMissingCapture` → `missing_capture`, `ErrPromotionNotApplicable` → `promotion_not_applicable`, `ErrUnsupported` → `unsupported`, `ErrTimeout` → `timeout`. **Una `funcion` desconocida → `invalid`** (la declaración que pide `FORMATO.md` 1.2 §4.1); `argumentos` que no casan con la función → `invalid` (`unknown_field` si traen un campo desconocido).

**9.2 Funciones de consulta** (`funcion` de cada consulta de `FORMATO.md` §2; la CLI y el oráculo implementan las mismas; un `credit_id` busca el crédito en `datos.credits`):
| `funcion` | Go | `argumentos` | Salida |
|---|---|---|---|
| `cotizar_compra` | `QuotePurchase` | `credit_id`, `purchase` (`PurchaseInput`) | `Purchase` |
| `promociones_aplicables` | `ApplicablePromotions` | `credit_id`, `date`, `principal_cents`, `term_count`, `store_id`, `merchant` | `{"promotion_ids": […]}` |
| `linea_tiempo` | `AvailabilityTimeline` | `credit_id`, `today`, `from`, `to`; `reading` opcional (sustituye `Credit.Reading` solo en esa consulta; la variante «omitida» manda `""`) | `Timeline` |
| `calendario_unificado` | `UnifiedCalendar` | `from`, `to` (todos los créditos de `datos`) | `Calendar` |
| `efectivo_ventana` | `CashFlow` | `from`, `to` (calendario de `datos` en esa ventana y `datos.cash`) | `CashWindow` |
| `deuda_viva` | `LiveDebt` | `as_of` | `DebtSummary` |
| `compras_vivas` | `LivePurchases` | `credit_id`, `at` | `{"refs": […]}` |
| `cotizar_monto` | `QuoteByAmount` | `amount_cents`, `rate`, `term_count` | `AmountQuote` |
| `sensibilidad` | `Sensitivity` | `SensitivityInput` | `{"rows": […]}` |
| `interes_tarjeta`, `minimo_tarjeta` | `CardPeriodInterest`, `CardMinimumPayment` | `credit_id`, `statement` (`StatementInput`); `paid_cents` y `paid_date` solo en el interés | `{"interest_cents", "provisional"}`, `{"minimum_cents", "provisional"}` (§14.1); el `credit_id` es de una tarjeta o de un crédito de tienda con datos K18 (§5.8) |
| `amortizar` | `Amortize` | `principal_cents`, `annual_rate`, `vat_bp`, `term_count`, `frequency`, `first_date` | `Purchase` sin partes (§5.4) |
| `comparar` | `Compare` | `timeout_seconds` (el valor de `comparator.timeout_seconds`) y `request` (`Request` sin `today` ni `credits`, que salen de `datos`) | `Result` |
| `vista_*` | §9.3 | ver §9.3 | claves literales de `model.json` o de `data.json` |

**9.3 Vistas con las claves de `model.json` y de `data.json`** (`views.go`; para las fixtures canon, R-05 y R-06 y `TestCanonHojas437`; `FORMATO.md` 1.2 §2.1). Dinero en centavos sin sufijo. `tasa` y `tasa_pct` salen como la misma razón `{"num": 80050, "den": 309450}` (la de la despensa 1) y `compara.py` aplica el ×100 de `tasa_pct` en aritmética exacta antes de redondear a binary64 (`FORMATO.md` §3). Verificado: binary64(8,005,000 / 309,450) = 25.868476328970754 = `parametros.tasa_pct`, mientras que binary64(`tasa`) × 100 en flotante **no** coincide; por eso el ×100 va en la fracción (solicitud al integrador para escribirlo así en `FORMATO.md`).
| `funcion` (argumentos) | Reproduce | Claves y origen |
|---|---|---|
| `vista_parametros` (`credit_id`) | `parametros` | `tasa`, `tasa_pct` (razón de la despensa 1); `papel` (monto del primer `datos.cash`); `limite_credito` (`Limit`); `karen_fija.{capital, intereses, total, pago, conjunto}` = parte externa de la despensa 1 (`Principal`, `Interest`, `Total`, `FirstPayment`, `Composition.Joint`) |
| `vista_despensa` (`credit_id`, `n`) | `despensas[n−1]` | `n`, `capital`, `intereses`, `total`; `pagos[*].{fecha, monto, axel, karen, capital_repuesto}`; `<persona>.{capital, conjunto, propio, papel_dentro, libre, intereses, total, pago_por_quincena}` con `<persona>` = `code` |
| `vista_disposicion` (`credit_id`, `n`) | `disposicion_por_despensa[n−1]` | `despensa` (= n) y `<persona>.{…}` como arriba |
| `vista_credito` (`credit_id`, `today`, `from`, `to`) | `credito` | `eventos[*].{fecha, tipo, monto, capital_repuesto, interes_no_repuesto, disponible_A, disponible_B}` con `tipo` = inicio, compra, papel o pago para `Kind` = start, purchase, cash o payment; `comparacion_series.{esperado_A, calculado_A}` = `SeriesA`, `comparacion_series.{esperado_B, calculado_B}` = `SeriesB`, `coinciden` (booleano); `minimo_A`, `minimo_B` |
| `vista_sensibilidad` (`SensitivityInput`) | `sensibilidad` | lista de `{N, intereses, total, pago_quincenal, axel_total, axel_capital, axel_menor_que_actual, costo_mensual_total, axel_mensual, karen_mensual}`; `N` en centavos (dinero según `FORMATO.md` §3) |
| `vista_interes_por_monto` (`amounts_cents`, `rate`, `term_count`) | `interes_por_monto` | lista de `{monto, intereses, total, pago_quincenal}` |
| `vista_dos_quincenas` (`credit_id`, `purchase`) | `data.json:financiamiento.dos_quincenas` (R-05, R-06) | `QuotePurchase` de `purchase` (corrige r2-M1 a): `total` (`Total`), `intereses` (`Interest`), `tasa_pct` = `{"num", "den", "decimales"}` (`Rate` con sus decimales capturados; `compara.py` aplica el ×100 y HALF_UP a `decimales` y compara con 17.27, `FORMATO.md` 1.2 §3) y `<persona>.{intereses, total, por_quincena}` (`PersonShare.Interest`, `Total`, `FirstPayment`) con `<persona>` = `code` |

**Claves de persona (`FORMATO.md` 1.2 §4 regla 5; corrige r3-N4):** cada objeto `<persona>` de `vista_despensa`, `vista_disposicion` y `vista_dos_quincenas`, y las claves de persona de `pagos[*]` (`axel` y `karen` en el canon), llevan **una clave por cada parte de la compra**, el titular y cada externa, con su `code`, siempre, aunque su importe sea 0; nunca la de una persona que no es parte de la compra. Orden: el de §4 (titular primero, luego externas por (`PersonID`, `Code`)). Así `esperado` pide las claves de todas las partes y una que falte en una implementación sale `__falta__`. Las claves de `vista_parametros` (`karen_fija` = la parte externa de la despensa 1, §4) y de `vista_sensibilidad` son literales de `model.json` y no dependen de los `code`.

Los textos de `model.json` (`descripcion`, `nota`, `cubre`) no se emiten ni se comparan.

**9.4 Frontera S3 → S2:** cada compra calculada se entrega con la forma exacta de `S2-interfaz.md` §2.5 `CreditPurchase` (`purchase_date`, `principal_cents`, `interest_cents`, `total_cents`, `status`, `shares`, `embedded` con el papel dentro, `payments` con `seq`, `due_date`, `amount_cents`, `principal_cents`, `interest_cents`, `shares`, `paid`). S2 verifica y no corrige.

**9.5 Registro en el binario común y CLI (`FORMATO.md` 1.2 §4 y §4.1; corrige r2-M1 c).** S3 agrega **solo** `backend/cmd/motores/s3.go`, contra la API que fija `main.go` (`Query`, `Evaluator`, `register`):
```go
package main

import (
	"encoding/json"

	"finanzas-hogar/internal/motor_cred/evaluar"
)

func init() { register("s3", evalS3) }

// evalS3 es el Evaluator de S3 (FORMATO.md 1.2 §4.1): input es entradas.datos ya resuelto.
func evalS3(input json.RawMessage, q Query) (out any, code string) {
	return evaluar.Evaluate(input, q.Func, q.Args)
}
```
`evaluar.Evaluate` decodifica `input` con `DecodeInput` (§5.7), despacha por `funcion` (§9.2, §9.3) y devuelve la salida, con dinero en enteros de centavos, o un código de §9.1; **una `funcion` desconocida devuelve `code` = `invalid`**. S3 **no edita `main.go` ni `main_test.go`** (son del agente Go de S2): argumentos, lista de fixtures, `$ref` y `parche`, validación de `formato` y de `motor`, recorte con `esperado` como plantilla (`FORMATO.md` §4), escritura determinista y códigos de salida (0 si escribió el archivo, 1 si no pudo leer o escribir, 2 por uso incorrecto o subcomando sin motor registrado) son de `main.go`; si S3 necesita un cambio ahí, lo pide al integrador. Si `main.go` aún no existe cuando termina el núcleo, S3 deja `s3.go` contra esta API, prueba `evaluar.Evaluate` desde `evaluar_test.go`, no crea un `main` sustituto y lo anota como pendiente; el integrador compila el binario al juntar. Corrida: `cd /home/deepsec/finanzas-hogar/backend && go build -o ../datos/tmp/motores ./cmd/motores && cd .. && datos/tmp/motores s3 dump --fixtures oraculo/fixtures/s3 --out datos/tmp/motor-s3.json [--caso <nombre>]` → `{"formato": "salida/1", "motor": "s3", "casos": {…}}`, con `--out` bajo `datos/tmp/`. No hay modo `--canon`: el árbol de las 6 ramas lo arma `TestCanonHojas437` con las vistas.

## 10. Fixtures y oráculo (`oraculo/fixtures/FORMATO.md` 1.2; corrige r1-H1, r2-M1, r2-M2, r2-M4, r3-N2, r3-N4 y r4-P1)
**10.0 Autoría de `oraculo/fixtures/s3/` y del oráculo, y orden de trabajo** (`FORMATO.md` 1.2 §1; alcance disjunto, protocolo §3; corrige r3-N4 y r4-P1). Nadie edita un archivo de otro autor; un cambio se pide al integrador.

| Archivo | Autor | Cómo lo escribe |
|---|---|---|
| `_datos/kueski-semilla.json` y todos los casos de la tabla de abajo salvo los `s3-comparador-*` (canon, R-05 y R-06, cruces con S2 y sin canon; cuáles llevan `provisional`: §14.1) | agente de **regresión** S3 | con `oraculo/cred_genera_fixtures.py`, determinista: las hojas canon y las de `data.json` las convierte leyendo solo `model.json` y `data.json` (`FORMATO.md` §3); la semilla, con la receta de «Datos compartidos» (abajo); las cifras sin canon, de esta interfaz, `S3.md` y las decisiones (citadas en `fuente`), o a mano con la cuenta en `descripcion` |
| `docs/estado/S3-escenarios.md`, `_datos/comparador-caso-1.json`, `-2.json`, `-3.json` y los 6 casos `s3-comparador-caso-1`, `-caso-2`, `-caso-3`, `-sin-restricciones`, `-tiempo` y `-carga` | **optimizador** S3 | en la secuencia de abajo; `-sin-restricciones` y `-tiempo` son variantes del caso 1 (`$ref` a `_datos/comparador-caso-1.json`) |
| `oraculo/cred_lib.py`, `oraculo/cred_escenarios.py`, `oraculo/cred_run_fixtures.py`, `oraculo/cred_genera_fixtures.py` | agente de **regresión** S3 | §10.2; solo biblioteca estándar, `python3 -I`; nunca lee la salida del motor |
| `oraculo/compara.py` | agente Python de S2 (`FORMATO.md` §5) | S3 lo usa y no lo edita |

**Secuencia del comparador (S3-13, «se fijan por escrito … antes de correr»):** (1) el optimizador escribe en `S3-escenarios.md` los 3 casos (entradas, los 4 valores de restricción, la intención de cada óptimo, §10.1) y su tabla a mano completa (`status`, `evaluated_count`, `valid` con llaves, valores y listas `provisional`, `rejected` con `violated_constraint`, los 4 `optima` distintos y `recommendation`), sin haber corrido su comparador sobre esos casos; (2) transcribe esa tabla a `_datos/comparador-caso-N.json` y a los `esperado` de los `s3-comparador-*` sin cambiar un número; (3) recién entonces corre su comparador, y el agente de regresión corre `cred_escenarios.py`; una diferencia con la tabla es hallazgo y se documentan ambos valores (protocolo §6), nunca se corrige la tabla para que cuadre; (4) agrega a `S3-escenarios.md`, en una sección aparte, los tiempos medidos (`S3-TIEMPO`) y el comportamiento al agotarse el tiempo (S3-12).

**Mientras falten archivos de otro autor** (los tres agentes trabajan en paralelo):
- **Núcleo S3:** no escribe en `oraculo/` ni crea sustitutos. (1) Arma en sus pruebas el `Input` de la semilla con un ayudante de prueba `seedFromModel` que sigue la receta de «Datos compartidos» leyendo `model.json` y `data.json` al correr (capital, intereses, partes, fechas, límite y contado salen de sus hojas; lo que no es hoja, los 15 días, la lectura `"B"` y los decimales 2 y 4 de K12, sale de esta interfaz) y, para `sensibilidad` e `interes_por_monto`, cuyas entradas no están en el `Input`, un segundo ayudante de prueba `sensitivityFromModel` que sigue la receta de «Insumo de sensibilidad y de cotización por monto» (abajo); con los dos corre `TestCanon85`, `TestCanonHojas437` y las demás pruebas con cifras canon (incluidas `TestTasaPorMonto` y `TestTasaSensibilidad`), sin leer archivos de otro autor; las pruebas sin canon usan datos en línea en `_test.go` con las cifras de esta interfaz, nunca cifras canon escritas; (2) lo que lee archivos de `oraculo/fixtures/s3/` (`TestOraculoFixtures`) falla con «falta oraculo/fixtures/s3/<archivo> (autor: <agente>)», nunca con `t.Skip`; (3) lo anota como pendiente y el integrador corre §10.2 completo al juntar. Cuando `_datos/kueski-semilla.json` existe, `TestCanonHojas437` exige además que sea igual, como JSON, al `Input` de `seedFromModel`, y cuando existen `s3-sensibilidad` y `s3-interes-por-monto`, que los argumentos de sus consultas sean iguales al insumo de `sensitivityFromModel`: una diferencia es hallazgo de uno de los dos.
- **Optimizador S3:** sus fixtures no esperan a nadie (secuencia de arriba). Antes de tenerlas prueba en `_test.go` con los ejemplos de esta interfaz (0/100/200 → 81,650, §8.4; 100/200/400 → 1, 2/3 y 0, §8.5); `TestRestricciones`, `TestOptimos` y `TestRestriccionesTiempo` leen sus fixtures y fallan con el mensaje de arriba si faltan. Si el núcleo aún no existe, lo anota como pendiente y no escribe un sustituto.
- **Regresión S3:** sin motor, `python3 -I oraculo/compara.py s3 --sin-motor` (`FORMATO.md` §5). Sin las fixtures del comparador, prueba `cred_escenarios.py` con los mismos ejemplos de §8.4 y §8.5; sin argumentos (§10.2) escribe «falta oraculo/fixtures/s3/<archivo> (autor: optimizador S3)» por cada caso que falte y termina con código 1. Un número de un `esperado` del optimizador que el oráculo no reproduce es hallazgo con ambos valores; nunca edita la fixture ajena.

**Forma:** la de `FORMATO.md` §2, sin extensiones: un archivo por caso en `oraculo/fixtures/s3/<caso>.json` con `formato` = `"fixture/1"`, `motor` = `"s3"`, `caso` (= nombre del archivo), `descripcion`, `fuente` (lista `{tipo, ruta}`), `provisional` (`null` o `{"bloqueo": "docs/estado/bloqueos/<archivo>.md", "que": "<qué parte>"}`), `entradas` (`datos` = `Input` de §5.7, en línea o `{"$ref": "_datos/<archivo>.json"}`, y `consultas` con las funciones de §9.2) y `esperado` (`{id: salida}`: subconjunto de claves o `{"error": "<código>"}`). Un `parche` (JSON Merge Patch, RFC 7396) reemplaza listas completas, así que las variantes que cambian una compra, un pago o un corte llevan sus `datos` en línea. Se retira la forma de la versión 1.0 de este contrato (carpeta plana, `{"$model": …}`, `cred_run_fixtures.py --escribir`, `cred_compara.py`, `cmd/finanzas … --fixture|--canon`).
**De dónde sale cada `esperado`** (lo escribe el autor que fija §10.0 **antes** de correr motor y oráculo; nunca se copia de la salida de uno de ellos):
- **Canon** (`fuente` `model.json`; `id` = ruta de `model.json`: `parametros`, `despensas[0]`, `credito`, `sensibilidad`, …): las hojas de `model.json` convertidas con `FORMATO.md` §3 por un script que solo lee `model.json`. `compara.py s3` las vuelve a comparar con `model.json` («model.json (ramas de s3): 437 de 437 hojas coinciden») y `TestCanon85` y `TestCanonHojas437` leen `model.json` al correr (S3-01).
- **R-05 y R-06** (corrige r2-M1 a; `FORMATO.md` 1.2 §2.1, §3 y §6): una consulta con `id` `data.json:financiamiento.dos_quincenas` y `funcion` `vista_dos_quincenas` (§9.3), con `fuente` `{"tipo": "data.json", "ruta": "financiamiento.dos_quincenas"}`. La compra es d1 a 2 quincenas **sin** cotización, así que los intereses salen de la tabla (53450/309450 con `decimales` 2: 53450 exactos) y la tasa trae sus decimales capturados (K12). `esperado` = las 9 hojas de `datos/semilla/data.json` convertidas por un script que solo lee `data.json`: `total` 362900, `intereses` 53450, `tasa_pct` `{"num": 53450, "den": 309450, "decimales": 2}` (se compara con 17.27 por la regla de `decimales`), `axel.{intereses, total, por_quincena}` 32870, 223170, 111585 y `karen.{intereses, total, por_quincena}` 20580, 139730, 69865 (las 9 salen con las reglas de §4). `compara.py s3` las contrasta con `data.json` en cada corrida (su línea `data.json:` deja de salir «0 de 0») y `TestCanon85` lee `data.json` al correr. La consulta libre `dos-quincenas` (`cotizar_compra` con la misma compra) queda como cruce adicional.
- **Cruces con S2** (`s3-deuda-viva-10-abr-2027`): `id` libres con funciones genéricas (`deuda_viva`, `efectivo_ventana`, `compras_vivas`), `fuente` `criterio` S3-01 y `regla` RN-21; los valores son hojas de las ramas de S2 de `model.json` en centavos. No entran a las líneas de `model.json` de `compara.py s3` (esas hojas son de S2, `S2-interfaz.md` §4); su cruce con `model.json` lo hace `TestCanon85` (§11.1).
- **Sin canon:** valores de un criterio, una regla o una decisión (`fuente` `criterio`, `regla` o `decision`; p. ej. S3-07: 0.00 y 77362 ×3 + 77364) o cálculo a mano explicado en `descripcion` (`fuente` `manual`). Las PROVISIONALES llevan `provisional` con su bloqueo y no entran en cifras canon.
- **Comparador:** la tabla a mano de `docs/estado/S3-escenarios.md`, fijada antes de correr (S3-13) y transcrita a la fixture por el mismo optimizador sin cambiar un número (secuencia de §10.0), con `fuente` `manual` y `criterio` S3-12 y S3-13: `esperado` trae `status`, `evaluated_count`, `valid` **completo** (cada escenario con su llave, sus valores y su lista `provisional`, §14.1), `rejected` completo (llave y `violated_constraint`), `optima` y `recommendation`. Así `compara.py s3` compara esperado, oráculo y motor en todo el conjunto («mismo conjunto, mismos valores», S3-12). Por eso cada caso se diseña chico (decenas de candidatos, lo que cabe en una tabla a mano); el tiempo se mide aparte con `s3-comparador-carga`.

**Datos compartidos** `oraculo/fixtures/s3/_datos/kueski-semilla.json` (el `Input` del canon): `today` 2026-10-10; un crédito id 1, `revolving_line`, «Kueski», `holder` axel, `limit_cents` 1000000, `fortnightly`, `first_payment_min_days` 15, `reading` "B", `reading_confirmed` false, `reading_a_anchor` 2026-10-10, `term_rates` 2 q = 53450/309450 (`decimales` 2) y 4 q = 80050/309450 (`decimales` 4), ambas `manual` sin `quote_date`, sin promociones; compras `d1` (2026-10-10, 309450, 4, `quoted_interest_cents` 80050, `rate_decimals` 4, K12: «25.8685 %»; axel 190300 con `joint_cents` 85100 y `paper_in_cents` 0; karen externa 119150 con 85100 y 19900) y `d2`, `d3`, `d4` (2026-12-10, 2027-02-10, 2027-04-10; 400000, 4, sin cotización ni `rate_decimals`, salen de la tabla; axel 280850 con 85100 y 19900; karen 119150 con 85100 y 0), todas `active`; `payments` y `statements` vacíos; `cash`: 2027-01-10 y 2027-03-10, 19900, karen, «Papel de contado».

**Insumo de sensibilidad y de cotización por monto** (ayudante de prueba `sensitivityFromModel` del núcleo y argumentos de las consultas de `s3-sensibilidad` y `s3-interes-por-monto`, que escribe el agente de regresión con la misma receta; corrige r4-P1). Lo que es hoja sale de `model.json` al correr, convertido con `FORMATO.md` §3 (dinero × 100 exacto); lo que no es hoja sale de esta interfaz. Las cifras entre paréntesis comprueban la receta y no se escriben en las pruebas:
- `vista_sensibilidad` (`SensitivityInput`, §7): `amounts_cents` = `sensibilidad[*].N` en su orden (300000 a 500000 de 50000 en 50000); `rate` = la tasa de 4 q de la semilla, `{"num": despensas[0].intereses, "den": despensas[0].capital, "decimales": 4}` (80050/309450, la misma de `term_rates` 4 q y de d1; los decimales salen de K12, como en la semilla); `term_count` = el plazo de d1 en la semilla, que es el número de `despensas[0].pagos` (4); `external` = `{"person": "karen", "external": true, "principal_cents": parametros.karen_fija.capital, "joint_cents": null, "paper_in_cents": null}` (119150; `Sensitivity` solo usa el capital externo, §7, y la composición informativa va en null); `holder_current_principal_cents` = `despensas[0].axel.capital` (190300); `accrual_months` = 2 (despensa bimestral, §7; no es hoja); `cash_monthly_cents` = HALF_UP(`parametros.papel` / 2) (19900 / 2 = 9950; el «/ 2» es el del valor estable, «199/2» de `valor_estable.componentes.formula` y R-44; no es hoja).
- `vista_interes_por_monto` (§9.3): `amounts_cents` = `interes_por_monto[*].monto` en su orden (100, 10000, 100000, 200000, 309450, 400000, 500000); `rate` y `term_count` = los de `vista_sensibilidad` (tasa de 4 q de la semilla, 4).

Verificado con un script propio (reporte `docs/estado/wf1/correccion-contrato-S3-r4.md`): con este insumo y las reglas de §7, las 45 hojas de `sensibilidad` y las 28 de `interes_por_monto` coinciden (73 de 73), y la fila `Amount` = 400000 da R-37…R-43 con las aserciones de §11.1 (7 de 7). `seedFromModel` y `sensitivityFromModel` leen solo `model.json` y `data.json`; si `_datos/kueski-semilla.json`, `s3-sensibilidad` o `s3-interes-por-monto` faltan, las pruebas canon no fallan por eso (solo `TestOraculoFixtures`, con el mensaje de autor de arriba).

| Caso (`oraculo/fixtures/s3/`) | `datos` · consultas | `esperado` · `fuente` · criterio |
|---|---|---|
| `s3-kueski-despensa-1` | semilla · `parametros` (`vista_parametros`), `despensas[0]` (`vista_despensa` n 1), `disposicion_por_despensa[0]` (`vista_disposicion` n 1), `data.json:financiamiento.dos_quincenas` (`vista_dos_quincenas`: d1 a 2 q sin cotización, con `quoted_interest_cents` y `rate_decimals` null, tasa de la tabla 53450/309450 con `decimales` 2) y el cruce libre `dos-quincenas` (`cotizar_compra` con la misma compra) | hojas de `parametros.{tasa, tasa_pct, karen_fija.*}`, `despensas[0]`, `disposicion_por_despensa[0]` (`model.json`); `data.json:financiamiento.dos_quincenas`: sus 9 hojas convertidas (R-05, R-06; `fuente` `data.json`); `dos-quincenas`: `interest_cents` 53450, `total_cents` 362900, `source` `term_table`, `rate_display` «17.27 %» · S3-01, S3-04, S3-06, S3-19 |
| `s3-kueski-despensa-2`, `-3`, `-4` | semilla · `despensas[k]`, `disposicion_por_despensa[k]` (n = k + 1) | `model.json` · S3-01, S3-04, S3-19 |
| `s3-feb-2027` | semilla · `cotizar_compra` el 2027-02-10 y el 2028-02-10 (4 q) | fechas 2027-02-28, 03-15, 03-30, 04-15 (`despensas[2].pagos`) y 2028-02-29, 03-15, 03-30, 04-15 (a mano con RN-06) · `criterio` S3-04, `manual` · S3-04 |
| `s3-series-a-b` | semilla · `credito` (`vista_credito`: `today` y `from` 2026-10-10, `to` 2027-04-10), `parametros` (`vista_parametros`), `lectura-omitida` (`linea_tiempo` con `reading` "") | `credito` y `parametros.{papel, limite_credito}` (`model.json`); omitida: `reading` "B", `pending_confirmation` true, `series_b_cents` · S3-03, S3-05 |
| `s3-series-desde-otra-fecha` | semilla · `linea_tiempo` con `from` 2026-12-16 y con `from` 2027-02-10 | inicio y eventos = los del canon tras el último evento ≤ `from` (§5.2): 9,094.50 / 6,000.00 y 8,094.50 / 5,000.00 · valores de `credito.eventos` con `id` libres, `manual` · S3-03 |
| `s3-reposicion-separada` | en línea (ejemplo 1 de §10.3) · `linea_tiempo` del 2026-10-10 al 2026-11-15 | la lista **completa** del `esperado` del ejemplo 1 (§10.3): inicio; abono de Karen del 29-oct (0 y 0); abono de Axel del 31-oct que completa el pago 1 (repone 77363, interés 20012); pago 2 proyectado el 15-nov; series, mínimos y bandera (§5.2, corrige r2-M2) · `decision` K15, `manual` · S3-05 |
| `s3-deuda-viva-10-abr-2027` | semilla · `deuda-viva` (`deuda_viva` `as_of` 2027-04-10), `efectivo-p1`…`efectivo-p6` (`efectivo_ventana` del 10 al 9), `efectivo-horizonte` (2026-10-10 a 2027-04-10), `vivas` (`compras_vivas` `at` 2027-04-10) | 629344 con d3 125870 y d4 503474; 97375, 194750, 223243, 271636, 251738, 271636; 1310378 con `cash_cents` 39800; `["d3", "d4"]` · `criterio` S3-01, `regla` RN-21 · S3-01, S3-11 |
| `s3-sensibilidad` | sin crédito · `sensibilidad` (`vista_sensibilidad`: montos 300000 a 500000 de 50000 en 50000, tasa 80050/309450, n 4, externa karen 119150, titular actual 190300, `accrual_months` 2, `cash_monthly_cents` 9950; = insumo de `sensitivityFromModel`, arriba) | `sensibilidad` (`model.json`) · S3-01, S3-02 |
| `s3-interes-por-monto` | sin crédito · `interes_por_monto` (`vista_interes_por_monto`: 100, 10000, 100000, 200000, 309450, 400000, 500000; tasa 4 q; n 4; = insumo de `sensitivityFromModel`, arriba) | `interes_por_monto` (`model.json`) · S3-01, S3-02 |
| `s3-tasa-plazo` | semilla · `cotizar_compra` a 6, 8, 10 y 12 q; 400000 a 4 q; d1 con cotización (`rate_decimals` 4) y la misma con `rate_decimals` null (corrige r3-N2) | `{"error": "no_rate"}` ×4; `interest_cents` 103474; d1: `rate_display` «25.8685 %»; sin `rate_decimals`: `interest_cents` 80050, `rate` `{"num": 80050, "den": 309450, "decimales": -1}` y `rate_display` `""` (§5.1) · `criterio` S3-06, `decision` K12 · S3-06 |
| `s3-promocion-0` | en línea: la línea de la semilla con una promoción de fixture («Chedraui», 4 q, 0/1, mínimo 200000, del 2026-10-28 al 2026-10-30) · `cotizar_compra` y `promociones_aplicables` el 28, el 30, el 27 y con 199999 | 0 el 28 y el 30; el 27 y con 199999: sin promoción aplicable, `promotion_not_applicable` si se pide y tasa normal sin ella; 309450 al 0 % → montos 77362 ×3 + 77364 e `interest_cents` 0, sin pedir capital ni interés por pago (FIRMES desde la puerta 2, Q2-38 = (a): capital = monto e interés 0 por pago; la fixture puede pedirlos) · `criterio` S3-07 · S3-07 |
| `s3-tarjeta-sin-intereses` | en línea: dos tarjetas de fixture (límite 2000000, corte 5, pago 25), una con datos K18 y otra sin ellos · `cotizar_compra` (cargo normal; MSI 3 de 100000), `interes_tarjeta` pagado completo un día antes y el mismo día | 0 en los 4; 33333 ×2 + 33334; solo hojas FIRMES (`interest_cents`, `total_cents`, montos y la lista `provisional`), sin fechas ni `source` de la tarjeta ni capital e interés por pago de MSI (FIRMES desde la puerta 2, Q2-32, Q2-37 y Q2-38; la fixture puede pedirlos) · `criterio` S3-08 · S3-08 |
| `s3-tarjeta-minimo` (PROVISIONAL `S3-formula-tarjeta-pdf`) | en línea: tarjeta con datos de fixture (455/1000 dec 1, diaria, `vat_included` false con IVA 1600, 500 bp, fijo 20000) y la misma sin datos (`vat_included` null) · `interes_tarjeta` un día después, `minimo_tarjeta` | a mano en `descripcion`; sin datos `{"error": "missing_capture"}` · `manual` · S3-08 |
| `s3-prestamo-frances` (FIRME desde la puerta 2, Q2-28 = (a)) | sin crédito · `amortizar` 1000000 a 12 mensual, 24/100 anual, IVA 1600; 100000 a 7 (obliga ajuste) | a mano; propiedades de S3-09 · `manual` · S3-09 |
| `s3-tienda-abonos` | en línea: crédito de tienda con `term_rates` 3 abonos a tasa 0 (`{"num": 0, "den": 1, "decimales": 0}`, abonos sin interés, §5.8) y otro sin `term_rates` ni datos K18 · `cotizar_compra` 100000 en 3 abonos | montos 33333, 33333, 33334 e `interest_cents` 0, `source` `term_table`, sin pedir capital ni interés por pago (FIRMES desde la puerta 2, Q2-38 = (a): capital = monto); sin datos `{"error": "missing_capture"}` · `criterio` S3-10 · S3-10 |
| `s3-tienda-abonos-interes` (PROVISIONAL `S3-formula-tarjeta-pdf`, Q2-31) | en línea: crédito de tienda con datos K18 de fixture (`vat_included` false con IVA 1600) y sin `term_rates` · `interes_tarjeta` y `minimo_tarjeta` con un `statement` de fixture pagado un día después de su fecha de pago; `cotizar_compra` 100000 en 3 abonos | intereses y mínimo con la fórmula de §5.3 (a mano en `descripcion`); `cotizar_compra`: `{"error": "unsupported"}` (§5.8) · `manual` · S3-10 |
| `s3-unificado-sin-corte` (calendario de la tarjeta FIRME desde la puerta 2, Q2-32 = A; `provisional` según §14.1 por las demás entradas) | en línea: Kueski (d1), tarjeta con cargo normal y MSI 3, préstamo y tienda sin interés (tasa 0, §5.8) · `calendario_unificado` | orden, Σ por quincena = Σ entradas, 0 duplicados, los pagos programados de la tarjeta · `manual`, `criterio` S3-11 · S3-11 |
| `s3-unificado` (corte `bank` que sustituye pagos, FIRME desde la puerta 2, Q2-32 = A; `provisional` según §14.1 por las demás entradas) | lo mismo más un corte `bank` cuyo `due_date` es el del cargo normal y la primera mensualidad · `calendario_unificado` | el corte sustituye a esos 2 pagos: Σ calendario = Σ programados no sustituidos + `balance_cents` · `manual` · S3-11 |
| `s3-unificado-duplicado` | lo mismo con dos cortes de la misma fecha | `{"error": "invalid"}` · `criterio` S3-11 · S3-11 |
| `s3-comparador-caso-1`…`-3` (FIRMES desde la puerta 2, Q2-32 a Q2-38; `provisional` solo si §14.1 da una lista no vacía para sus créditos; corrige r2-M4 y r3-N1) | §10.1 y ejemplo 2 de §10.3 · `comparar` con `timeout_seconds` 5 | tabla a mano de `S3-escenarios.md`, con las listas `provisional` de cada escenario y del resultado según §14.1 · S3-12, S3-13, S3-14 |
| `s3-comparador-sin-restricciones` | caso 1 con `global` y `override` vacíos | `missing_constraints` (las 4), `evaluated_count` 0, sin óptimos ni recomendación · `criterio` S3-12 · S3-12 |
| `s3-comparador-tiempo` | caso 1 con `timeout_seconds` 0 | `{"error": "timeout"}` (503 `unavailable` por `ResponseFor`, §8.7) · `criterio` S3-12, `decision` K43 · S3-12 |
| `s3-comparador-carga` (FIRME desde la puerta 2, Q2-33 y Q2-36) | rejilla sintética (6 créditos con límite, 7 plazos con tasa, `date_window_days` 7, 6 tiendas, una promoción por crédito) · `comparar` con 5 | `status` y `evaluated_count` (cuenta a mano de la rejilla); da el tiempo de §8.7 · `manual` · S3-12 |

**10.1 Casos del comparador** (el optimizador los fija con cifras en `docs/estado/S3-escenarios.md` **antes** de correr el motor y verifica a mano que los 4 óptimos son distintos; no se ajustan después). Cada caso trae en la fixture los 4 valores de restricción (capturas de fixture, no de Axel) y debe tener: ≥ 2 créditos con límite; ≥ 2 plazos con tasa; una sola promoción 0 % aplicable solo en algunas fechas y en la tienda de peor U/C; ≥ 2 tiendas cuya diferencia de U/C supere el efecto de los intereses; compras existentes que hagan desigual el perfil por quincena; al menos un escenario descartado por cada restricción y uno marcado sobre el tope. Intención: (1) intereses → el escenario con la promoción; (2) parejo → el plazo y la fecha que rellenan las quincenas bajas de lo existente; (3) disponible → la línea con más holgura; (4) rendimiento → la tienda de mejor U/C con el plazo de menor tasa. Caso 1: dos líneas tipo Kueski; caso 2: línea + tarjeta con MSI (el cargo normal de 0.00 compite con la promoción: el optimizador desempata por diseño, no por suerte); caso 3: tope por persona capturado y `Override` que cambia una restricción global. Para que la tabla a mano quepa, cada caso enumera decenas de candidatos.

**10.2 Oráculo** (método distinto: `fractions.Fraction`, día por día; solo biblioteca estándar, `python3 -I`; no importa Go ni lee su salida, S3-16; temporales en `datos/tmp/`): `cred_lib.py` (fórmulas propias de §3–§7 y las vistas de §9.3), `cred_escenarios.py` (búsqueda exhaustiva propia: candidatos, conjunto factible, objetivos con `math.isqrt`, normalización, 4 óptimos y recomendación; sin argumentos corre solo los casos del comparador) y `cred_run_fixtures.py --fixtures oraculo/fixtures/s3 --out oraculo/salidas/s3.json` (corre todas las consultas de §9.2 y escribe la salida de `FORMATO.md` §4). La comparación de las cuatro fuentes la hace `oraculo/compara.py s3` (`FORMATO.md` §5): esperado, oráculo y motor y, en las fixtures canon, `model.json`. Así el oráculo nunca lee al motor, el motor nunca lee al oráculo y un tercero compara. La «diferencia de conjuntos motor/oráculo = 0» de S3-12 es la línea «oráculo contra motor» de `compara.py s3` en las fixtures del comparador. Orden de corrida (S3-02, S3-09, S3-12):
```
cd /home/deepsec/finanzas-hogar
python3 -I oraculo/cred_run_fixtures.py --fixtures oraculo/fixtures/s3 --out oraculo/salidas/s3.json
cd backend && go build -o ../datos/tmp/motores ./cmd/motores && cd .. && datos/tmp/motores s3 dump --fixtures oraculo/fixtures/s3 --out datos/tmp/motor-s3.json
python3 -I oraculo/compara.py s3
```
Salida esperada de la última línea de cobertura: «model.json (ramas de s3): 437 de 437 hojas coinciden» (S3-02) y «diferencias: 0», sin `SIN-DUEÑO`; las PROVISIONALES salen en su línea `provisionales: N de N`.

**10.3 Ejemplos de `entradas`** (corrige r1-H2). Ejemplo 1, línea de tiempo con datos en línea (`s3-reposicion-separada`, completo con su `esperado`; corrige r2-M2):
```json
{
  "formato": "fixture/1",
  "motor": "s3",
  "caso": "s3-reposicion-separada",
  "descripcion": "Despensa 1 con la parte de Karen pagada aparte el 29-oct y la de Axel el 31-oct: el capital del pago 1 (773.63) se repone el 31-oct y no antes (K15); el pago 2, sin abonos, se proyecta completo el 15-nov (§5.2)",
  "fuente": [{"tipo": "decision", "ruta": "K15"}, {"tipo": "criterio", "ruta": "S3-05"}, {"tipo": "manual", "ruta": "descripcion"}],
  "provisional": null,
  "entradas": {
    "datos": {
      "today": "2026-10-10",
      "credits": [
        {
          "credit": {
            "id": 1, "name": "Kueski", "kind": "revolving_line", "holder": "axel",
            "limit_cents": 1000000, "frequency": "fortnightly", "first_payment_min_days": 15,
            "reading": "B", "reading_confirmed": false, "reading_a_anchor": "2026-10-10",
            "statement_day": null, "payment_due_day": null, "annual_rate": null, "rate_base": "",
            "vat_included": true, "vat_on_interest_bp": null, "min_payment_pct_bp": null, "min_payment_fixed_cents": null,
            "term_rates": [
              {"term_count": 2, "rate": {"num": 53450, "den": 309450, "decimales": 2}, "source": "manual", "quote_date": null},
              {"term_count": 4, "rate": {"num": 80050, "den": 309450, "decimales": 4}, "source": "manual", "quote_date": null}
            ],
            "promotions": []
          },
          "purchases": [
            {
              "ref": "d1",
              "status": "active",
              "input": {
                "date": "2026-10-10", "principal_cents": 309450, "term_count": 4,
                "quoted_interest_cents": 80050, "rate_decimals": 4, "promotion_id": null, "store_id": null, "merchant": "",
                "first_payment_date": null,
                "shares": [
                  {"person": "axel", "external": false, "principal_cents": 190300, "joint_cents": 85100, "paper_in_cents": 0},
                  {"person": "karen", "external": true, "principal_cents": 119150, "joint_cents": 85100, "paper_in_cents": 19900}
                ]
              }
            }
          ],
          "payments": [
            {"purchase_ref": "d1", "seq": 1, "date": "2026-10-29", "payer": "karen", "amount_cents": 37493},
            {"purchase_ref": "d1", "seq": 1, "date": "2026-10-31", "payer": "axel", "amount_cents": 59882}
          ],
          "statements": []
        }
      ],
      "cash": []
    },
    "consultas": [
      {"id": "linea", "funcion": "linea_tiempo", "argumentos": {"credit_id": 1, "today": "2026-10-10", "from": "2026-10-10", "to": "2026-11-15"}}
    ]
  },
  "esperado": {
    "linea": {
      "events": [
        {"date": "2026-10-10", "kind": "start", "purchase_ref": "", "seq": 0, "amount_cents": 0, "principal_restored_cents": 0, "interest_not_restored_cents": 0, "available_a_cents": 1000000, "available_b_cents": 690550},
        {"date": "2026-10-29", "kind": "payment", "purchase_ref": "d1", "seq": 1, "amount_cents": 37493, "principal_restored_cents": 0, "interest_not_restored_cents": 0, "available_a_cents": 1000000, "available_b_cents": 690550},
        {"date": "2026-10-31", "kind": "payment", "purchase_ref": "d1", "seq": 1, "amount_cents": 59882, "principal_restored_cents": 77363, "interest_not_restored_cents": 20012, "available_a_cents": 1077363, "available_b_cents": 767913},
        {"date": "2026-11-15", "kind": "payment", "purchase_ref": "d1", "seq": 2, "amount_cents": 97375, "principal_restored_cents": 77363, "interest_not_restored_cents": 20012, "available_a_cents": 1154726, "available_b_cents": 845276}
      ],
      "series_a_cents": [1000000, 1077363, 1154726],
      "series_b_cents": [690550, 767913, 845276],
      "min_a_cents": 1000000,
      "min_b_cents": 690550,
      "reading": "B",
      "pending_confirmation": true,
      "a_anchor": "2026-10-10",
      "a_anchor_defaulted": false
    }
  }
}
```
Cuentas del `esperado` (a mano, §5.2): inicio en `From` = ancla = 2026-10-10, A = 1,000,000 (d1 tiene fecha = ancla) y B = 1,000,000 − 309,450 = 690,550; el abono de Karen no repone; el de Axel completa el pago 1 (37,493 + 59,882 = 97,375 = M) y suma K = 77,363 a A y a B; el pago 2 (vencimiento 15-nov ≥ `Today`, sin abonos) se proyecta completo y suma otros 77,363; el pago 3 (30-nov) queda fuera de `To`.
Ejemplo 2, comparador (`s3-comparador-caso-1`, solo `entradas`; `_datos/comparador-caso-1.json` es un `Input` como el del ejemplo 1 con las dos líneas y sus compras existentes; las cifras de aquí ilustran la forma y las del caso las fija `S3-escenarios.md`):
```json
{
  "entradas": {
    "datos": {"$ref": "_datos/comparador-caso-1.json"},
    "consultas": [
      {
        "id": "caso-1",
        "funcion": "comparar",
        "argumentos": {
          "timeout_seconds": 5,
          "request": {
            "amount_cents": 400000,
            "desired_date": "2026-11-20",
            "store_id": null,
            "merchant_name": "",
            "shares": [
              {"person": "axel", "external": false, "principal_cents": 280850, "joint_cents": null, "paper_in_cents": null},
              {"person": "karen", "external": true, "principal_cents": 119150, "joint_cents": null, "paper_in_cents": null}
            ],
            "global": {"min_available_buffer_cents": 150000, "max_live_purchases": 2, "zero_debt_by_date": "2027-06-30", "date_window_days": 2},
            "override": {"min_available_buffer_cents": null, "max_live_purchases": null, "zero_debt_by_date": null, "date_window_days": null},
            "weights": {"lowest_interest": 1, "even_low_payment": 1, "max_available": 1, "most_value_per_peso": 1},
            "caps": [],
            "baskets": [
              {"store_id": 1, "merchant_name": "Tienda A", "base_units_milli": 52000, "cost_cents": 400000},
              {"store_id": 2, "merchant_name": "Tienda B", "base_units_milli": 61000, "cost_cents": 400000}
            ]
          }
        }
      }
    ]
  }
}
```

## 11. Pruebas (nombres que casan con `-run` de `docs/criterios/S3.md`)
| Prueba Go | Paquete | Qué verifica | Criterio |
|---|---|---|---|
| `TestCanon85` | motor_cred | 85 subpruebas `R-01`…`R-85`, cada una con su aserción exacta y 0 `t.Skip`; lee `oraculo/model.json` y `datos/semilla/data.json` (R-05, R-06) al correr. Filas S3: la salida de su función, con el `Input` de `seedFromModel` y el insumo de `sensitivityFromModel` (§10.0), contra la ruta de `S3.md` §6. Filas S2: tabla §11.1 | S3-01 |
| `TestCanonHojas437` | motor_cred | las 6 ramas armadas con las vistas de §9.3 contra las 437 hojas de §12, leídas de `model.json`: `parametros`, `despensas`, `credito` y `disposicion_por_despensa` sobre el `Input` de `seedFromModel`; `sensibilidad` (`vista_sensibilidad`) e `interes_por_monto` (`vista_interes_por_monto`) sobre el insumo de `sensitivityFromModel` (§10.0); si existe `_datos/kueski-semilla.json`, exige además que sea igual a ese `Input`, y si existen `s3-sensibilidad` y `s3-interes-por-monto`, que los argumentos de sus consultas sean iguales a ese insumo (§10.0) | S3-01, S3-02 |
| `TestSeriesAB` | motor_cred | 28 valores; lectura omitida → B y bandera | S3-03 |
| `TestSeriesABDesdeOtraFecha` | motor_cred | `From` = cada fecha de evento del canon posterior al ancla: inicio = A y B del canon tras el último evento ≤ `From` y eventos siguientes idénticos; `From` < ancla → `ErrInvalid`; lectura A sin ancla → `ErrInvalid` (§5.2) | S3-03 |
| `TestCalendarioKueski` | motor_cred | 16 fechas (R-13…R-31) y 28-feb | S3-04 |
| `TestReposicionCapital`, `TestReposicionParteSeparada`, `TestReposicionAbonoParcial` | motor_cred | Σ capital repuesto = capital; K15 con la lista completa del ejemplo 1 de §10.3 (un evento por abono; el que completa repone K e I); el resto de un pago parcial no vencido se proyecta en su vencimiento y uno vencido no repone (§5.2) | S3-05 |
| `TestTasaPorCompra`, `TestTasaSinCaptura`, `TestTasaPorMonto`, `TestTasaSensibilidad` | motor_cred | RN-08, `ErrNoRate`, `interes_por_monto`, `sensibilidad`; decimales de la tasa cotizada: `rate_decimals` 4 → «25.8685 %», null → `rate_display` `""` con los mismos intereses (§5.1) | S3-06, S3-02 |
| `TestPromocionFronteras` | motor_cred | 4 fronteras | S3-07 |
| `TestTarjetaPagoTotal`, `TestTarjetaMSI`, `TestTarjetaSinCaptura`, `TestTarjetaMinimo` | motor_cred | 0.00, MSI, «falta capturar» (cada uno de los 5 datos de K18, incluido `VATIncluded` null; Q2-67 = 1), fórmula (PROVISIONAL, `S3-formula-tarjeta-pdf`) | S3-08 |
| `TestPrestamoFrances` | motor_cred | cuadre al centavo | S3-09 |
| `TestTiendaAbonos`, `TestTiendaSinCaptura`, `TestTiendaConInteres` | motor_cred | tasa 0 = abonos sin interés con capital = monto (FIRME, Q2-38 = (a)); sin `TermRates` ni datos K18 → `ErrMissingCapture`; con datos K18 → `ErrUnsupported` en `QuotePurchase` e interés y mínimo por periodo con la fórmula (PROVISIONAL, `S3-formula-tarjeta-pdf`) (§5.8) | S3-10 |
| `TestUnificado` | motor_cred | orden, Σ por quincena = Σ entradas, Σ `ByPerson` = total, 0 duplicados; sin corte, los pagos programados de la tarjeta; con corte, el corte sustituye a los pagos de su vencimiento y Σ calendario = Σ programados no sustituidos + Σ `Balance` (FIRME, Q2-32 = A); duplicado → `ErrInvalid` (§6) | S3-11 |
| `TestEfectivoVentana` | motor_cred | `CashFlow` en las 6 ventanas del canon y en el horizonte (§6.1) | S3-01 (R-45…R-52, R-60) |
| `TestDeudaViva`, `TestCompraViva` | motor_cred | RN-21 al 10-abr-2027; intervalo inclusivo | S3-01 (R-53…R-55), S3-12 |
| `TestReparto` | motor_cred | invariantes de §5 y R-07…R-16 | S3-19 |
| `TestDespachoCotizacion` | motor_cred | la tabla de §5.8 renglón por renglón: entradas obligatorias, `Source`, `PurchaseKind`, `RateBasis` y precedencia de errores por tipo; préstamo dentro de `QuotePurchase` con las partes de §4; sin `AnnualRate` → `ErrMissingCapture`; préstamo con `VATIncluded` true, false y null (§5.4) | S3-06, S3-08, S3-09, S3-10 |
| `TestProvisionalListas` | motor_cred y comparador | listas de §14.1: contenido por salida (desde la puerta 2 solo entran `S3-formula-tarjeta-pdf` y `S3-iva-lectura`; ya no entran `S3-convencion-prestamo`, `S3-formula-tarjeta`, `S3-objetivos-medida`, `S3-origen-cargo-normal`, `S3-serie-sin-intereses` ni `U4-iva-credito`), sin repetidos y en orden alfabético; `Result.Provisional` = unión de las de sus escenarios; `[]` con `missing_constraints` | S3-20 |
| `TestEvaluar` | evaluar | despacho de §9.2 y §9.3; `funcion` desconocida → `invalid`; `vista_dos_quincenas` contra `data.json` leído al correr | S3-01, S3-09 |
| `TestOraculoFixtures` | ambos | el motor contra `esperado` en todas las fixtures de `oraculo/fixtures/s3/` (las lee con la forma de `FORMATO.md`; no lee la salida del oráculo); un archivo de §10.0 que falte la hace fallar con su nombre y su autor, nunca `t.Skip` | S3-09, S3-12, S3-13 |
| `TestRestricciones`, `TestRestriccionesFaltan` | comparador | 0 aceptados con violación; conjunto = `esperado` (tabla a mano); «faltan restricciones» con 0 recomendaciones; el cargo normal de la tarjeta siempre se genera (§8.2) | S3-12 |
| `TestTiempoClaveSembrada` | comparador | lee como texto `backend/db/migrations/00003_puerta1b.sql` (sin `database/sql`), toma el valor sembrado de `comparator.timeout_seconds` y comprueba `TimeoutFromSetting(valor)` = 5 s; `TimeoutFromSetting(-1)` → `ErrInvalid` | S3-12 (K43, K44) |
| `TestRestriccionesTiempo` | comparador | caso 1 con `timeout_seconds` 0 → `ErrTimeout` y `ResponseFor` = (503, `unavailable`); los 3 casos y `s3-comparador-carga` con 5 terminan sin error y escriben `S3-TIEMPO <caso> <ms> ms` con `t.Logf` | S3-12 |
| `TestOptimos` | comparador | 3 casos × 4 óptimos distintos = `esperado`; objetivo 2 con las quincenas de §8.4 | S3-13 |
| `TestObjetivoParejo` | comparador | ejemplo 0/100/200 → 81,650; una quincena de total 0 dentro de la ventana cuenta en n; empate en desviación → gana el de menos `OverCap` (§8.6) | S3-13 |
| `TestPuntaje` | comparador | pesos Σ = 1.00; 100/200/400 → 1, 2/3, 0; empate → 1; pesos editados cambian la recomendación; un descartado con valor extremo no mueve la normalización (§8.5) | S3-14 |

Ninguna prueba lleva cifras canon escritas en el código. Cobertura ≥ 85 % en cada paquete (K25, S3-18); `gofmt -l` vacío y `go vet` limpio. S3-15/S3-16/S3-17 se verifican con los `grep` y `sha256sum` del criterio.
**S3-12 en un motor puro (corrige r1-H8).** Se prueba en S3: el valor sembrado de la clave (`TestTiempoClaveSembrada`), la conversión pura (`TimeoutFromSetting`), el 503 `unavailable` como dato (`ResponseFor`), el vencimiento con 0 y el fin con 5 (`TestRestriccionesTiempo`) y el tiempo medido de cada caso (líneas `S3-TIEMPO`, que el optimizador copia a `S3-escenarios.md`). Queda para WF-2 (API, C1-23): que cada petición a `compareScenarios` lea el valor vigente de `comparator.timeout_seconds` en `app_setting` (también después de cambiarlo en Ajustes) y que la respuesta HTTP sea 503 `unavailable`; solicitud para que `S3.md` S3-12 lo diga así.

### 11.1 `TestCanon85`: aserción de las 30 filas S2 (corrige r1-H3 y r4-P1)
Cada subprueba lee su valor esperado de `model.json` al correr; los valores S3 salen de las funciones con el `Input` de `seedFromModel` y el insumo de sensibilidad de §10.0 (`sensitivityFromModel`), sin leer archivos de otro autor; si existen `_datos/kueski-semilla.json` y `s3-sensibilidad`, exige además igualdad, como `TestCanonHojas437`. Ventanas k = 0…5: [`periodos[k].inicio`, `periodos[k].fin`], leídas de `model.json`. Donde la cifra usa una regla de S2 (devengo, quincena del valor estable, promedio), la prueba la aplica con funciones de S3 (`SplitHalfUp`, `RoundHalfUp`) como cruce; la cifra la produce `motor_fin` y la confirma S2-01, y aquí la subprueba pasa o falla por su propia aserción, nunca con `t.Skip`. Comprobado con un script propio (reporte `docs/estado/wf1/correccion-contrato-S3-r1.md`): las 30 cuadran al centavo.
| Filas | Función S3 | Aserción exacta |
|---|---|---|
| R-37, R-38, R-39 | `Sensitivity`, fila `Amount` = 400000 | `MonthlyTotal`, `HolderMonthly`, `ExternalMonthly` = `valor_estable.mes.{total, axel, karen}` (261687, 176751, 84936) |
| R-40, R-41, R-42 | la misma fila | `RoundHalfUp(x, 2)` de esos tres = `valor_estable.quincena.{total, axel, karen}` (130844, 88376, 42468; la quincena es la mitad HALF_UP del mes, `S2-interfaz.md` §3.6) |
| R-43 | la misma fila | `MonthlyPurchase` = `valor_estable.componentes.mitad_despensa` (251737) |
| R-44 | `vista_parametros` | `RoundHalfUp(papel, 2)` = `valor_estable.componentes.papel_contado_mes` (9950) |
| R-45…R-50 | `CashFlow` en la ventana k y `CreditState.Quote` | efectivo: `Total` = `periodos[k].efectivo.total`; devengado: Σ de la parte que cae en la ventana k de `SplitHalfUp(Total_j, 2)` de cada compra j (la ventana de su fecha y la siguiente, `S2-interfaz.md` §3.2) + `Cash` = `periodos[k].devengado.total`; diferencia: Σ_{i ≤ k} (devengado_i − efectivo_i) = `periodos[k].diferencia_acumulada` |
| R-51, R-52 | lo mismo | Σ_{k=0..5} devengado = `horizonte.total_devengado`; `CashFlow(2026-10-10, 2027-04-10).Total` = Σ_k efectivo = `horizonte.total_efectivo` |
| R-53, R-54, R-55 | `LiveDebt` al 2027-04-10 | `Total`, `Pending` de d3 y de d4 = `horizonte.deuda_viva_10_abr`, `horizonte.deuda_viva_detalle.despensa_3_pendiente`, `.despensa_4` |
| R-56, R-57 | de R-52 y R-51 | `RoundHalfUp(total, 6)` = `horizonte.promedio_mensual_efectivo` y `.promedio_mensual_devengado` |
| R-58 | lo mismo | R-51 − R-52 (de S3) = `horizonte.diferencia_final` |
| R-59 | `UnifiedCalendar` del 2026-10-10 al 2027-04-10 | Σ `Amount` de las entradas de d3 = `horizonte.pagos_despensa_3_en_horizonte` |
| R-60 | `CashFlow(2026-10-10, 2027-04-10)` | `Cash` = `horizonte.papel_contado` |
| R-68…R-73 | `CreditState.Quote` (composición) e `Input.Cash` | paquetes = cada parte con `PaperIn` > 0 (fecha de la compra, persona, forma «dentro de despensa n», monto) más cada `CashEvent` (fecha, persona, «de contado», monto), en orden de fecha: `paquete` = posición, `compra` = fecha, `quien` = persona sin distinguir mayúsculas, `forma` y `monto` = `papel[k]` |

## 12. Las 437 hojas numéricas de S3 (S3-02): función, fixture y prueba
Patrones: `[*]` = todos los índices; `{a,b}` = alternativas. «Eco» = valor de entrada que la salida devuelve igual (se compara igual). Conteo verificado con un script temporal que recorre `model.json` y casa cada hoja con exactamente un renglón: **437 cubiertas, 0 sin dueño, 0 con dos dueños** (reporte `docs/estado/wf1/contrato-S3.md`).

| # | Patrón | Hojas | Función | Fixture | Prueba |
|---|---|---|---|---|---|
| 1 | `parametros.tasa` | 1 | `QuotedRate` (fracción → float solo al comparar) | kueski-despensa-1 | CanonHojas437, TasaPorCompra |
| 2 | `parametros.tasa_pct` | 1 | `QuotedRate` × 100 | kueski-despensa-1 | CanonHojas437, TasaPorCompra |
| 3 | `parametros.papel` | 1 | eco de `CashEvent.Amount` | series-a-b | CanonHojas437 |
| 4 | `parametros.limite_credito` | 1 | eco de `Credit.Limit` | series-a-b | CanonHojas437, SeriesAB |
| 5 | `parametros.karen_fija.{capital,intereses,total,pago,conjunto}` | 5 | `PersonShare` externa y `Composition.Joint` | kueski-despensa-1 | Reparto |
| 6 | `despensas[*].n` | 4 | eco (orden de compra) | kueski-despensa-1..4 | CanonHojas437 |
| 7 | `despensas[*].capital` | 4 | `Purchase.Principal` | kueski-despensa-1..4 | Canon85 |
| 8 | `despensas[*].intereses` | 4 | `Purchase.Interest` (cotizado; tabla) | kueski-despensa-1..4 | TasaPorCompra |
| 9 | `despensas[*].total` | 4 | `Purchase.Total` | kueski-despensa-1..4 | Canon85 |
| 10 | `despensas[*].pagos[*].monto` | 16 | `SplitTrunc` | kueski-despensa-1..4 | Canon85 |
| 11 | `despensas[*].pagos[*].axel` | 16 | parte del titular (resto) | kueski-despensa-1..4 | Reparto |
| 12 | `despensas[*].pagos[*].karen` | 16 | `SplitHalfUp` de la externa | kueski-despensa-1..4 | Reparto |
| 13 | `despensas[*].pagos[*].capital_repuesto` | 16 | `SplitHalfUp` del capital | kueski-despensa-1..4 | ReposicionCapital |
| 14 | `despensas[*].{axel,karen}.capital` | 8 | eco `ShareInput.Principal` | kueski-despensa-1..4 | Reparto |
| 15 | `despensas[*].{axel,karen}.conjunto` | 8 | eco `Composition.Joint` | kueski-despensa-1..4 | Reparto |
| 16 | `despensas[*].{axel,karen}.propio` | 8 | `Composition.Own` | kueski-despensa-1..4 | Reparto |
| 17 | `despensas[*].{axel,karen}.papel_dentro` | 8 | eco `Composition.PaperIn` | kueski-despensa-1..4 | Reparto |
| 18 | `despensas[*].{axel,karen}.libre` | 8 | `Composition.Free` | kueski-despensa-1..4 | Reparto |
| 19 | `despensas[*].{axel,karen}.intereses` | 8 | `PersonShare.Interest` | kueski-despensa-1..4 | Reparto |
| 20 | `despensas[*].{axel,karen}.total` | 8 | `PersonShare.Total` | kueski-despensa-1..4 | Reparto |
| 21 | `despensas[*].{axel,karen}.pago_por_quincena` | 8 | `PersonShare.FirstPayment` | kueski-despensa-1..4 | Reparto |
| 22 | `credito.eventos[*].monto` | 17 | `AvailabilityTimeline` | series-a-b | SeriesAB |
| 23 | `credito.eventos[*].capital_repuesto` | 17 | `AvailabilityTimeline` | series-a-b | ReposicionCapital |
| 24 | `credito.eventos[*].interes_no_repuesto` | 17 | `AvailabilityTimeline` | series-a-b | ReposicionCapital |
| 25 | `credito.eventos[*].disponible_A` | 17 | `AvailabilityTimeline` | series-a-b | SeriesAB |
| 26 | `credito.eventos[*].disponible_B` | 17 | `AvailabilityTimeline` | series-a-b | SeriesAB |
| 27 | `credito.comparacion_series.esperado_A[*]` | 14 | `Timeline.SeriesA` | series-a-b | SeriesAB |
| 28 | `credito.comparacion_series.esperado_B[*]` | 14 | `Timeline.SeriesB` | series-a-b | SeriesAB |
| 29 | `credito.comparacion_series.calculado_A[*]` | 14 | `Timeline.SeriesA` | series-a-b | SeriesAB |
| 30 | `credito.comparacion_series.calculado_B[*]` | 14 | `Timeline.SeriesB` | series-a-b | SeriesAB |
| 31 | `credito.minimo_A` | 1 | `Timeline.MinA` | series-a-b | SeriesAB |
| 32 | `credito.minimo_B` | 1 | `Timeline.MinB` | series-a-b | SeriesAB |
| 33 | `sensibilidad[*].N` | 5 | eco | sensibilidad | TasaSensibilidad |
| 34 | `sensibilidad[*].{intereses,total}` | 10 | `Sensitivity` (`ApplyRate`) | sensibilidad | TasaSensibilidad |
| 35 | `sensibilidad[*].pago_quincenal` | 5 | `ReferencePayment` (total exacto) | sensibilidad | TasaSensibilidad |
| 36 | `sensibilidad[*].{axel_total,axel_capital}` | 10 | `Sensitivity` | sensibilidad | TasaSensibilidad |
| 37 | `sensibilidad[*].{costo_mensual_total,axel_mensual,karen_mensual}` | 15 | `Sensitivity` (desde t exacto) | sensibilidad | TasaSensibilidad |
| 38 | `disposicion_por_despensa[*].despensa` | 4 | eco | kueski-despensa-1..4 | CanonHojas437 |
| 39 | `disposicion_por_despensa[*].{axel,karen}.{capital,conjunto,propio,papel_dentro,libre,intereses,total,pago_por_quincena}` | 64 | lo mismo que 14–21 (vista de `despensas[k]`) | kueski-despensa-1..4 | Reparto |
| 40 | `interes_por_monto[*].monto` | 7 | eco | interes-por-monto | TasaPorMonto |
| 41 | `interes_por_monto[*].{intereses,total}` | 14 | `QuoteByAmount` | interes-por-monto | TasaPorMonto |
| 42 | `interes_por_monto[*].pago_quincenal` | 7 | `ReferencePayment` (total exacto) | interes-por-monto | TasaPorMonto |

Por rama: `parametros` 9, `despensas` 144, `credito` 143, `sensibilidad` 45, `disposicion_por_despensa` 68, `interes_por_monto` 28 = **437**. Las hojas booleanas (`credito.comparacion_series.coinciden`, `sensibilidad[*].axel_menor_que_actual`) y de texto no cuentan; `axel_menor_que_actual` igual se devuelve y se prueba en `TestTasaSensibilidad`.

## 13. Criterio → sección
| Criterio | Sección | | Criterio | Sección |
|---|---|---|---|---|
| S3-01 | §4, §5.1, §6.1, §7, §9.3 (`vista_dos_quincenas`), §10, §11 `TestCanon85` y §11.1 | | S3-11 | §6 (el corte sustituye), §6.1 |
| S3-02 | §12, §10.2 `compara.py s3`, `TestCanonHojas437` | | S3-12 | §8.1–§8.3 (cargo normal en §8.2), §8.7 (`TimeoutFromSetting`, `ResponseFor`), §10.1, §11 |
| S3-03 | §5.2 (ancla de A) | | S3-13 | §8.4 (serie del objetivo 2), §8.6, §10.0 (secuencia del comparador), §10.1 |
| S3-04 | §5.1 calendario, `s3-feb-2027` | | S3-14 | §8.5 (conjunto: generados según la decisión C; empate: vale 1) |
| S3-05 | §5.2 abonos, resto proyectado y reposición (K15) | | S3-15 | §1.5, §2 |
| S3-06 | §5.1 intereses y decimales (`RateDecimals`, `Display`), `ErrNoRate`, §5.8 | | S3-16 | §10.0 (autor del oráculo), §10.2 |
| S3-07 | §5.1 promoción aplicable y serie sin intereses (FIRME, Q2-38) | | S3-17 | §1.1–§1.3 |
| S3-08 | §5.3, §5.8 | | S3-18 | §11 |
| S3-09 | §5.4, §5.8, §9.5 | | S3-19 | §4 |
| S3-10 | §5.5, §5.8 | | S3-20 | §14, §14.1 y reporte |

## 14. PROVISIONAL, hallazgos y propuestas
- **PROVISIONAL:** fórmula de intereses y mínimo de tarjeta y tienda, y cotización de una tienda con interés por datos K18 (`ErrUnsupported`) (bloqueo `S3-formula-tarjeta-pdf`: `S3-formula-tarjeta` preguntas 1 a 3, Q2-29 a Q2-31 = (a), pendiente de los estados de cuenta que entrega Axel); préstamo con `VATIncluded` null y `VATIncluded` = false sin porcentaje (`S3-iva-lectura` preguntas 1 y 2). **FIRMES desde la puerta 2:** tasa del periodo del préstamo (Q2-28), calendario de la tarjeta y corte que sustituye pagos (Q2-32), canasta y tienda del objetivo 4, serie y desempate del objetivo 2, ventana del colchón y créditos sin límite (Q2-33 a Q2-36), `Source` del cargo normal (Q2-37), serie de capital con intereses 0 (Q2-38) y entrada `VATIncluded` con los 5 datos de «falta capturar» (Q2-67 = 1, `U4-iva-credito`). La lectura A/B es la decisión pendiente con Kueski y se muestra con `PendingConfirmation` (decisión C). Toda salida que dependa de lo PROVISIONAL lleva `Provisional` con el nombre del bloqueo; qué lista lleva cada salida y cada fixture: §14.1.
- **Declaración (S3-14, K17):** la normalización mín-máx es entre los escenarios generados en el sentido de la decisión C, los que respetan las restricciones duras (`valid`); los descartados no entran (§8.5).
- **Hallazgo 1 (FIRME en este contrato):** el pago de referencia sale del total sin redondear; desde el total redondeado no reproduce `interes_por_monto[0]` (0.32 contra 0.31) ni `[5]` (1,258.69 contra 1,258.68).
- **Hallazgo 2:** `parametros.karen_fija.intereses` (308.22) coincide con HALF_UP(800.50 × 1,191.50 / 3,094.50) y con HALF_UP(1,191.50 × tasa); la regla de RN-08 (proporcional a los intereses) reproduce las 8 hojas de intereses por persona.
- **Hallazgo 3 (C1):** `Scenario.fortnight_payment_spread_cents` (máx − mín) no es la medida de K17 (desviación estándar); C1 no expone la desviación, la tienda del escenario ni la marca sobre el tope (solicitud en el reporte).
- **Hallazgo 4 (`FORMATO.md` §3):** `parametros.tasa_pct` (25.868476328970754) es binary64(8,005,000 / 309,450) y no binary64(`tasa`) × 100 en flotante: el ×100 se aplica en la fracción (§9.3).
- **Propuesta K31 (`optimizer.date_window_days`, sin valor sembrado):** 7 días; ver el reporte.
- **Cambios de la versión 1.1 (corrección r1):** H1 → §2, §9, §10 y §11 alineados a `FORMATO.md` v1.0 y fuente `data.json` para R-05 y R-06; H2 → tipos con llaves JSON en §3–§8, §5.7 y ejemplos en §10.3; H3 → §6.1 y §11.1; H4 → §8.5; H5 → §5.2 y §8.3; H6 → §8.4 (2) y §8.6; H7 → §6; H8 → §8.7 y §11. Ninguna cifra canon cambia.
- **Hallazgo 5 → bloqueo `docs/estado/bloqueos/S3-serie-sin-intereses.md` (corrige r3-N1):** en la línea revolvente, una promoción al 0 % con los dos modos de RN-07 (§5.1) da montos 77362 ×3 + 77364 y capital 77363 ×3 + 77361 para 309450, es decir, intereses por pago −1, −1, −1 y +3, que `scheduled_payment.interest_cents CHECK (interest_cents >= 0)` (00001 l. 808) rechaza al guardar; los meses sin intereses (§5.3) y los abonos sin interés de tienda (§5.5) usan capital = monto, regla que no está en K14 ni en RN-07. Con los montos de RN-07 fijos, la única serie de capital con intereses ≥ 0 en cada pago es capital = monto. Las tres series quedan PROVISIONALES con el bloqueo (opciones a, b y c) hasta que Axel decida; ningún monto ni cifra canon cambia.
- **Cambios de la versión 1.2 (corrección r2):** M1 → §2 (`evaluar/` y `s3.go`; `main.go` es de S2), §3 y §9.1 (`decimales`), §9.3 (`vista_dos_quincenas`), §9.5 (registro, `Evaluator` y `invalid` para una `funcion` desconocida) y §10 (`id` `data.json:financiamiento.dos_quincenas`), con `FORMATO.md` 1.1; M2 → §5.2 (un evento por abono, el que completa repone, resto proyectado de un pago parcial no vencido) y `esperado` completo del ejemplo 1 de §10.3; M3 → §5.8 (despacho por tipo, tienda sin interés con tasa 0, tienda con interés `unsupported`, préstamo con partes, `Source` y `PurchaseKind`; bloqueo nuevo `S3-origen-cargo-normal`), §5.3–§5.5 y §8.2 (el cargo normal siempre se genera); M4 → §14.1 y marcas `provisional` en las fixtures del comparador y del calendario unificado. Ninguna cifra canon cambia.

- **Cambios de la versión 1.3 (corrección r3):** N1 → bloqueo nuevo `S3-serie-sin-intereses` y marcas PROVISIONAL en §5.1 (serie), §5.3 (MSI), §5.5, §5.8, las fixtures `s3-promocion-0`, `s3-tarjeta-sin-intereses` y `s3-tienda-abonos` (no piden capital ni interés por pago), §14 y §14.1; N2 → `RateDecimals` en `PurchaseInput` (§5), precedencia de decimales y `Display()` con −1 (§3, §5.1), `rate_decimals` 4 en d1 de la semilla y del ejemplo 1, variante sin decimales en `s3-tasa-plazo` (§10) y `TestTasaPorCompra` (§11); N3 → bloqueo abierto `U4-iva-credito` citado en el encabezado, §5 (`VATOnInterestBP`), §5.3, §5.4, §5.5, §5.8, §14 y §14.1 (`interes_tarjeta`, `minimo_tarjeta`, préstamo sin IVA y fixtures `s3-tarjeta-minimo` y `s3-tienda-abonos-interes`); N4 → `FORMATO.md` 1.2 en todas las citas vigentes, §10.0 (autoría de cada archivo, trabajo de cada agente mientras falten y secuencia del comparador), §2, §9.3 (claves de persona, regla 5 de §4) y §11 (`seedFromModel`, archivo faltante). Ninguna cifra canon cambia.
- **Cambios de la versión 1.4 (corrección r4):** P1 → §10.0 (segundo ayudante `sensitivityFromModel` con la receta exacta del insumo de `vista_sensibilidad` y `vista_interes_por_monto`, igualdad con `s3-sensibilidad` y `s3-interes-por-monto` cuando existen, y la tabla de casos que lo cita), §11 (`TestCanon85` y `TestCanonHojas437`: qué ramas salen de cada ayudante) y §11.1 (los valores S3 ya no salen de archivos del agente de regresión). Ninguna cifra canon ni regla de cálculo cambia; los bajos Q1–Q4 y B1–B9 de la revisión r4 no se tocan en esta corrección.
- **Cambios de la versión 1.5 (puerta 2, ronda P2c):** agente `contratos-puerta2`: marcas de Q2-28 y Q2-32 a Q2-38 (FIRMES; encabezado, §5.1, §5.3, §5.4, §5.5, §5.8, §6, §8, fixtures, pruebas, §14 y §14.1). Agente `correccion-contratos-r1` (`docs/estado/p2c/REVISION-contratos-r1.md`): A-1 → Q2-67 = 1 (entrada `VATIncluded` en §5, 5 datos de K18 en §5.3, §5.5 y §5.8, `vatBP` del préstamo en §5.4 y §5.8, semilla Kueski con `vat_included` true, pruebas de §11 y listas de §14.1; bloqueo nuevo `S3-iva-lectura`); M-2 → fila `bank_card` de §5.8 FIRME (Q2-32, Q2-37 y Q2-38) y §10 sin «PROVISIONALES»; M-3 → bloqueo nuevo `S3-formula-tarjeta-pdf`, que sustituye a `S3-formula-tarjeta` en todas las marcas y listas `Provisional` de la fórmula. Ninguna cifra canon cambia.

### 14.1 Listas `Provisional` (regla FIRME; corrige r2-M4)
Cada salida con campo `Provisional` lleva los nombres de los bloqueos de los que depende (el nombre del archivo de `docs/estado/bloqueos/` sin `.md`), **sin repetidos y en orden alfabético por bytes**; `[]` = FIRME. Desde la puerta 2 solo pueden entrar los bloqueos abiertos `S3-formula-tarjeta-pdf` < `S3-iva-lectura`; los demás bloqueos de S3 y `U4-iva-credito` se resolvieron (Q2-28, Q2-32 a Q2-38 y Q2-67) y ya no entran en ninguna lista.

| Salida | Bloqueos que entran |
|---|---|
| `Purchase.Provisional` (`QuotePurchase`, `amortizar`) | `term_loan` con `VATIncluded` null: `S3-iva-lectura`; en todo lo demás, `[]` (desde la puerta 2: tarjeta, préstamo con `VATIncluded` true o false, `amortizar` con o sin `vat_bp` —es el IVA que se suma—, línea, MSI y tienda, con o sin intereses 0) |
| `CalendarEntry.Provisional` | entrada `scheduled`: la `Provisional` de su compra; entrada `statement`: `S3-formula-tarjeta-pdf` si el corte es `computed`; `[]` si es `bank` |
| `provisional` de `interes_tarjeta` y `minimo_tarjeta` | interés 0.00 por pago total a tiempo: `[]` (FIRME); interés con la fórmula y todo mínimo: `S3-formula-tarjeta-pdf` (la lista de datos de `ErrMissingCapture` y la lectura de `VATIncluded` son FIRMES desde Q2-67 = 1, §5.3) |
| `Scenario.Provisional` (en `valid`, `rejected`, `optima` y `recommendation`) | la `Provisional` de la compra del escenario ∪ las `Provisional` de las entradas del calendario unificado que forman la serie del objetivo 2 (§8.4) |
| `Result.Provisional` | `status` `missing_constraints`: `[]` (K16 es FIRME y no se calcula nada); `ok` o `no_valid_scenario`: la unión de las `Provisional` de todos los escenarios de `valid` y `rejected` |

`Timeline`, `CashWindow`, `DebtSummary` y las vistas no tienen campo `Provisional`; si una fixture las usa con compras cuya `Provisional` no es vacía, la marca va en la fixture (abajo).

**Fixtures.** Una fixture cuyo `esperado` trae alguna hoja que depende de un bloqueo lleva `provisional` = `{"bloqueo": "docs/estado/bloqueos/<archivo>.md", "que": "<qué parte>"}`. `FORMATO.md` 1.2 §2 admite un solo objeto: va el bloqueo principal y `que` termina con «otros bloqueos: » y la ruta de cada uno de los demás que aparecen en las listas `provisional` de su `esperado`, en el orden de arriba (solicitud al integrador para que `FORMATO.md` admita una lista). Una fixture que solo trae hojas FIRMES de salidas con lista no vacía (p. ej. `s3-tarjeta-sin-intereses`: los 0.00 y los montos de MSI) va con `provisional` null y su `esperado` no pide fechas ni `source` de la tarjeta.

| Fixture | `provisional` |
|---|---|
| `s3-comparador-caso-1`, `s3-comparador-caso-2`, `s3-comparador-caso-3`, `s3-comparador-carga` | null desde la puerta 2 (Q2-32 a Q2-38), salvo que §14.1 dé una lista no vacía para alguno de sus créditos (un préstamo con `VATIncluded` null: `S3-iva-lectura`); entonces ese bloqueo |
| `s3-comparador-sin-restricciones`, `s3-comparador-tiempo` | null (K16 y K43 FIRMES; no se calcula ningún escenario) |
| `s3-unificado-sin-corte`, `s3-unificado` | null desde la puerta 2 (calendario de la tarjeta y corte `bank` que sustituye pagos, Q2-32 = A), salvo `S3-iva-lectura` si su préstamo trae `vat_included` null |
| `s3-tarjeta-minimo`, `s3-tienda-abonos-interes` | `S3-formula-tarjeta-pdf` (fórmula de intereses y mínimo: Q2-29 a Q2-31 = (a), en espera de los estados de cuenta), sin «otros bloqueos» |
| `s3-prestamo-frances` | null desde la puerta 2 (Q2-28 = (a); `amortizar` con `vat_bp` null no suma IVA, FIRME por Q2-67 = 1, §5.4) |
| las demás | null |

