// Paquete precios: motor de precios y tiendas (S5). Núcleo puro, sin red ni base de datos:
// ingesta por fuente, equivalencia por unidad base, frescura, precio vigente, comparativa
// por ingrediente y por canasta, membresía de Sam's y mercados con día de visita.
// Contrato: docs/contratos/motores/S5-interfaz.md 1.5; reglas RN-12, RN-13, RN-16, RN-20.
// Dinero en centavos int64; comparaciones con fracciones exactas; HALF_UP solo al publicar.
package precios

import (
	"errors"
	"fmt"
	"math/big"
	"time"
)

// Errores del motor (S5-interfaz §17). Code los traduce al código de las fixtures.
var (
	ErrUnknownPresentation = errors.New("presentación desconocida")
	ErrUnknownStore        = errors.New("tienda desconocida")
	ErrFutureDate          = errors.New("fecha posterior a hoy")
	ErrDuplicatePrice      = errors.New("precio repetido")
	ErrNotConfirmed        = errors.New("OCR sin confirmación humana")
	ErrOverflow            = errors.New("desbordamiento aritmético")
	ErrInvalidArgument     = errors.New("argumento inválido")
	// ErrUnknownProduct no está en §17: sale como invalid_argument.
	ErrUnknownProduct = fmt.Errorf("%w: producto desconocido", ErrInvalidArgument)
)

// Fuentes de precio (price.source). Las tres primeras compiten con el mismo peso (K22, K37);
// site_assisted y profeco solo son referencia (K22: PROFECO nunca compite).
const (
	SourceTicketManual = "ticket_manual"
	SourceTicketOCR    = "ticket_ocr"
	SourceManual       = "manual"
	SourceSiteAssisted = "site_assisted"
	SourceProfeco      = "profeco"
)

// Competes dice si una fuente compite en la comparativa; ok = false si la fuente no existe.
func Competes(source string) (competes, ok bool) {
	switch source {
	case SourceTicketManual, SourceTicketOCR, SourceManual:
		return true, true
	case SourceSiteAssisted, SourceProfeco:
		return false, true
	}
	return false, false
}

// Tipos de tienda (store.kind).
const (
	KindChain  = "chain"
	KindBranch = "branch"
	KindMarket = "market"
)

// Group es un grupo de compra (product_group); PurchaseChannel = market o chain.
type Group struct {
	Code, Name, PurchaseChannel string
}

// Product es un producto genérico con su unidad base (g, ml, piece).
type Product struct {
	Code, Name, BaseUnit, Group string
}

// Presentation es una presentación con su contenido en milésimas de la unidad base.
type Presentation struct {
	Code, Product, Name string
	SizeMilli           int64
}

// Store es una cadena, sucursal o mercado (store + store_product_group).
type Store struct {
	Code, Kind, Parent, Name, ProfecoName, Municipality string
	VisitWeekday                                        int // 1 = lunes … 7 = domingo; 0 = sin día
	RequiresMembership, HasMembership, AssistedReading  bool
	MembershipFeeCents                                  *int64
	Groups                                              []string
}

// AssumesMembership: Sam's se compara como si tuvieran membresía (RN-13, decisiones E).
func (s Store) AssumesMembership() bool { return s.RequiresMembership && !s.HasMembership }

// PriceObs es una observación de precio (tabla price).
type PriceObs struct {
	Presentation, Store string
	PriceCents          int64
	ObservedDate        time.Time
	Source              string
	IsPromotional       bool
	SourceRef           string
	CreatedAt           time.Time
}

// Settings son los parámetros de app_setting price.* (K44); nunca constantes en código.
type Settings struct {
	StaleAfterDays         int  // price.stale_after_days, inicial 30 (K22)
	ZeroWeightAfterDays    int  // price.zero_weight_after_days, inicial 90 (K40)
	AssistedReadingEnabled bool // price.assisted_reading_enabled, inicial false (S5-07)
}

// DefaultSettings devuelve los valores iniciales sembrados en app_setting (00001, 00003).
func DefaultSettings() Settings {
	return Settings{StaleAfterDays: 30, ZeroWeightAfterDays: 90}
}

// Validate exige 1 ≤ Dv ≤ D0 para que la recta (D0 − edad)/(D0 − Dv + 1) tenga sentido.
func (s Settings) Validate() error {
	if s.StaleAfterDays < 1 || s.ZeroWeightAfterDays < s.StaleAfterDays {
		return fmt.Errorf("%w: días de frescura %d y %d", ErrInvalidArgument, s.StaleAfterDays, s.ZeroWeightAfterDays)
	}
	return nil
}

// Frac es una fracción exacta no negativa (num/den, den > 0), reducida.
type Frac struct {
	Num int64 `json:"num"`
	Den int64 `json:"den"`
}

// NewFrac reduce n/d; d debe ser > 0.
func NewFrac(n, d int64) Frac {
	g := gcd(abs64(n), d)
	if g == 0 {
		return Frac{0, 1}
	}
	return Frac{n / g, d / g}
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

// cmpRat compara a/b contra c/d (b, d > 0) con multiplicación cruzada exacta (math/big).
func cmpRat(a, b, c, d int64) int {
	l := new(big.Int).Mul(big.NewInt(a), big.NewInt(d))
	r := new(big.Int).Mul(big.NewInt(c), big.NewInt(b))
	return l.Cmp(r)
}

// cmpRatBig compara (a·b)/(c·d) contra (e·f)/(g·h) (denominadores > 0) sin desbordar.
func cmpRatBig(a, b, c, d, e, f, g, h int64) int {
	l := new(big.Int).Mul(big.NewInt(a), big.NewInt(b))
	l.Mul(l, new(big.Int).Mul(big.NewInt(g), big.NewInt(h)))
	r := new(big.Int).Mul(big.NewInt(e), big.NewInt(f))
	r.Mul(r, new(big.Int).Mul(big.NewInt(c), big.NewInt(d)))
	return l.Cmp(r)
}

// halfUp devuelve round(n/d) con ROUND_HALF_UP para n ≥ 0 y d > 0, sin desbordar.
func halfUp(n, d *big.Int) (int64, error) {
	q := new(big.Int).Mul(n, big.NewInt(2))
	q.Add(q, d)
	q.Quo(q, new(big.Int).Mul(d, big.NewInt(2)))
	if !q.IsInt64() {
		return 0, ErrOverflow
	}
	return q.Int64(), nil
}

// mul multiplica con detección de desbordamiento.
func mul(a, b int64) (int64, error) {
	r := new(big.Int).Mul(big.NewInt(a), big.NewInt(b))
	if !r.IsInt64() {
		return 0, ErrOverflow
	}
	return r.Int64(), nil
}

// ParseDate lee una fecha ISO (AAAA-MM-DD) como medianoche UTC (fecha de calendario).
func ParseDate(s string) (time.Time, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err == nil && t.Year() < 1 {
		// Año 0000: invalid_date, igual que el oráculo (date.fromisoformat). PROVISIONAL, bloqueo
		// S5-precio-vigente-captura pregunta 3.
		return time.Time{}, fmt.Errorf("año 0000: %q", s)
	}
	return t, err
}

// AgeDays son los días de calendario entre observed y today (negativo si observed es futura).
func AgeDays(observed, today time.Time) int {
	o := time.Date(observed.Year(), observed.Month(), observed.Day(), 0, 0, 0, 0, time.UTC)
	t := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	return int(t.Sub(o).Hours() / 24)
}

// isoWeekday: 1 = lunes … 7 = domingo.
func isoWeekday(t time.Time) int { return (int(t.Weekday())+6)%7 + 1 }
