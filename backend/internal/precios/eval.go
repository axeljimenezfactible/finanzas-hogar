package precios

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"finanzas-hogar/internal/precios/asistida"
)

// Errores propios de las fixtures (S5-interfaz §15.2).
var (
	ErrUnknownField = errors.New("clave desconocida")
	ErrMissingField = errors.New("falta una clave obligatoria")
	ErrInvalidDate  = errors.New("fecha inválida")
)

// Code traduce un error al código de las fixtures (§15.2, §17): snake_case sin «Err».
func Code(err error) string {
	for _, c := range []struct {
		err  error
		code string
	}{
		{ErrUnknownField, "unknown_field"}, {ErrMissingField, "missing_field"}, {ErrInvalidDate, "invalid_date"},
		{ErrUnknownPresentation, "unknown_presentation"}, {ErrUnknownStore, "unknown_store"},
		{ErrFutureDate, "future_date"}, {ErrDuplicatePrice, "duplicate_price"},
		{ErrNotConfirmed, "not_confirmed"}, {ErrOverflow, "overflow"},
		{asistida.ErrNoTermsVerdict, "no_terms_verdict"},
	} {
		if errors.Is(err, c.err) {
			return c.code
		}
	}
	return "invalid_argument"
}

// RepoPath valida una ruta *_file: relativa a la raíz del repositorio, sin «..» y, si under no
// está vacío, dentro de esa carpeta (§15.1). Devuelve la ruta limpia.
func RepoPath(p, under string) (string, error) {
	c := filepath.Clean(p)
	if p == "" || !filepath.IsLocal(c) || strings.Contains(p, "..") ||
		(under != "" && !strings.HasPrefix(c, filepath.Clean(under)+string(filepath.Separator))) {
		return "", fmt.Errorf("%w: ruta %q", ErrInvalidArgument, p)
	}
	return c, nil
}

// strict decodifica rechazando claves desconocidas (unknown_field); null o vacío = {}.
func strict(raw json.RawMessage, v any) error {
	if len(bytes.TrimSpace(raw)) == 0 || string(bytes.TrimSpace(raw)) == "null" {
		raw = json.RawMessage("{}")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		if strings.Contains(err.Error(), "unknown field") {
			return fmt.Errorf("%w: %v", ErrUnknownField, err)
		}
		return fmt.Errorf("%w: %v", ErrInvalidArgument, err)
	}
	return nil
}

func need(ok bool, what string) error {
	if !ok {
		return fmt.Errorf("%w: %s", ErrMissingField, what)
	}
	return nil
}

func date(s *string, what string) (time.Time, error) {
	if s == nil {
		return time.Time{}, need(false, what)
	}
	t, err := ParseDate(*s)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s %q", ErrInvalidDate, what, *s)
	}
	return t, nil
}

func str(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type inputJSON struct {
	Today    *string `json:"today"`
	Settings *struct {
		StaleAfterDays         *int  `json:"stale_after_days"`
		ZeroWeightAfterDays    *int  `json:"zero_weight_after_days"`
		AssistedReadingEnabled *bool `json:"assisted_reading_enabled"`
	} `json:"settings"`
	CatalogFile *string `json:"catalog_file"`
	Catalog     *struct {
		Groups []struct {
			Code, Name      string
			PurchaseChannel string `json:"purchase_channel"`
		} `json:"groups"`
		Products []struct {
			Code, Name        string
			BaseUnit          string  `json:"base_unit"`
			Group             *string `json:"group"`
			GramsPerUnitMilli *int64  `json:"grams_per_unit_milli"`
		} `json:"products"`
		Presentations []struct {
			Code, Product, Name string
			SizeMilli           int64 `json:"size_milli"`
		} `json:"presentations"`
	} `json:"catalog"`
	Stores []struct {
		Code, Kind         *string
		Parent             *string  `json:"parent"`
		Name               *string  `json:"name"`
		ProfecoName        *string  `json:"profeco_name"`
		Municipality       *string  `json:"municipality"`
		VisitWeekday       *int     `json:"visit_weekday"`
		RequiresMembership bool     `json:"requires_membership"`
		HasMembership      bool     `json:"has_membership"`
		MembershipFeeCents *int64   `json:"membership_fee_cents"`
		AssistedReading    bool     `json:"assisted_reading"`
		Groups             []string `json:"groups"`
	} `json:"stores"`
	Prices []struct {
		Presentation, Store *string
		PriceCents          *int64  `json:"price_cents"`
		ObservedDate        *string `json:"observed_date"`
		Source              *string `json:"source"`
		IsPromotional       bool    `json:"is_promotional"`
		SourceRef           *string `json:"source_ref"`
		CreatedAt           *string `json:"created_at"`
	} `json:"prices"`
	Carts []struct {
		Code, Store, Kind *string
		ObservedDate      *string `json:"observed_date"`
		Source            *string `json:"source"`
		Lines             []struct {
			Presentation   *string `json:"presentation"`
			QuantityMilli  *int64  `json:"quantity_milli"`
			UnitPriceCents *int64  `json:"unit_price_cents"`
		} `json:"lines"`
	} `json:"carts"`
	ReportFile   *string `json:"report_file"`
	ShoppingList []struct {
		Product                 *string `json:"product"`
		NeededMilli             *int64  `json:"needed_milli"`
		QuantityMilli           *int64  `json:"quantity_milli"`
		RecommendedPresentation *string `json:"recommended_presentation"`
	} `json:"shopping_list"`
	GroceryPurchases []struct {
		Date, Store *string
	} `json:"grocery_purchases"`
	Terms []struct {
		Site     *string `json:"site"`
		Verdict  *string `json:"verdict"`
		URL      *string `json:"url"`
		ReadDate *string `json:"read_date"`
	} `json:"terms"`
}

// Input es entradas.datos de una fixture S5 ya validada (§15.1).
type Input struct {
	Book      *Book
	Carts     map[string]Cart
	List      []ListItem
	Purchases []Purchase
	Terms     []asistida.Term
}

// ParseInput valida y arma el Input de §15.1. Las rutas *_file se abren relativas al
// directorio actual, que es la raíz del repositorio (CLI de §15).
func ParseInput(raw json.RawMessage) (*Input, error) {
	var in inputJSON
	if err := strict(raw, &in); err != nil {
		return nil, err
	}
	today, err := date(in.Today, "today")
	if err != nil {
		return nil, err
	}
	s := DefaultSettings()
	if in.Settings != nil {
		if v := in.Settings.StaleAfterDays; v != nil {
			s.StaleAfterDays = *v
		}
		if v := in.Settings.ZeroWeightAfterDays; v != nil {
			s.ZeroWeightAfterDays = *v
		}
		if v := in.Settings.AssistedReadingEnabled; v != nil {
			s.AssistedReadingEnabled = *v
		}
	}
	var cat *Catalog
	switch {
	case in.CatalogFile != nil && in.Catalog != nil:
		return nil, fmt.Errorf("%w: catalog_file y catalog", ErrInvalidArgument)
	case in.CatalogFile != nil:
		p, err := RepoPath(*in.CatalogFile, "")
		if err != nil {
			return nil, err
		}
		if cat, err = LoadCatalogFile(p); err != nil {
			return nil, err
		}
	case in.Catalog != nil:
		var gs []Group
		var ps []Product
		var pr []Presentation
		for _, g := range in.Catalog.Groups {
			gs = append(gs, Group{g.Code, g.Name, g.PurchaseChannel})
		}
		for _, p := range in.Catalog.Products {
			ps = append(ps, Product{p.Code, p.Name, p.BaseUnit, str(p.Group)})
		}
		for _, p := range in.Catalog.Presentations {
			pr = append(pr, Presentation{p.Code, p.Product, p.Name, p.SizeMilli})
		}
		if cat, err = NewCatalog(gs, ps, pr); err != nil {
			return nil, err
		}
	default:
		return nil, need(false, "catalog_file o catalog")
	}
	var stores []Store
	for _, x := range in.Stores {
		if err := need(x.Code != nil && x.Kind != nil, "stores[].code y kind"); err != nil {
			return nil, err
		}
		st := Store{Code: *x.Code, Kind: *x.Kind, Parent: str(x.Parent), Name: str(x.Name),
			ProfecoName: str(x.ProfecoName), Municipality: str(x.Municipality),
			RequiresMembership: x.RequiresMembership, HasMembership: x.HasMembership,
			MembershipFeeCents: x.MembershipFeeCents, AssistedReading: x.AssistedReading, Groups: x.Groups}
		if x.VisitWeekday != nil {
			st.VisitWeekday = *x.VisitWeekday
		}
		stores = append(stores, st)
	}
	book, err := NewBook(cat, stores, nil, s, today)
	if err != nil {
		return nil, err
	}
	for _, x := range in.Prices {
		if err := need(x.Presentation != nil && x.Store != nil && x.PriceCents != nil && x.Source != nil,
			"prices[]: presentation, store, price_cents y source"); err != nil {
			return nil, err
		}
		od, err := date(x.ObservedDate, "prices[].observed_date")
		if err != nil {
			return nil, err
		}
		var ca time.Time
		if x.CreatedAt != nil {
			// §15.1: instante RFC 3339. time.Parse guarda nanosegundos (PROVISIONAL, bloqueo S5-precio-vigente-captura
			// pregunta 2) y no acepta «t», «z» ni el segundo 60 (pregunta 3).
			if ca, err = time.Parse(time.RFC3339, *x.CreatedAt); err != nil || !rfc3339Offset(*x.CreatedAt) {
				return nil, fmt.Errorf("%w: created_at %q", ErrInvalidDate, *x.CreatedAt)
			}
		}
		if err := book.AddPrice(PriceObs{*x.Presentation, *x.Store, *x.PriceCents, od, *x.Source,
			x.IsPromotional, str(x.SourceRef), ca}); err != nil {
			return nil, err
		}
	}
	out := &Input{Book: book, Carts: map[string]Cart{}}
	if in.ReportFile != nil {
		if len(in.Carts) > 0 {
			return nil, fmt.Errorf("%w: report_file y carts", ErrInvalidArgument)
		}
		p, err := RepoPath(*in.ReportFile, "")
		if err != nil {
			return nil, err
		}
		r, err := LoadReport(p, cat)
		if err != nil {
			return nil, err
		}
		out.Carts[r.Bodega.Code], out.Carts[r.Walmart.Code] = r.Bodega, r.Walmart
	}
	for _, x := range in.Carts {
		if err := need(x.Code != nil && x.Store != nil && x.Kind != nil && x.Source != nil,
			"carts[]: code, store, kind y source"); err != nil {
			return nil, err
		}
		if _, dup := out.Carts[*x.Code]; dup || (*x.Kind != "expense" && *x.Kind != "cart") {
			return nil, fmt.Errorf("%w: canasta %q", ErrInvalidArgument, *x.Code)
		}
		od, err := date(x.ObservedDate, "carts[].observed_date")
		if err != nil {
			return nil, err
		}
		c := Cart{Code: *x.Code, Store: *x.Store, Kind: *x.Kind, ObservedDate: od, Source: *x.Source}
		for _, l := range x.Lines {
			if err := need(l.Presentation != nil && l.QuantityMilli != nil && l.UnitPriceCents != nil,
				"carts[].lines[]: presentation, quantity_milli y unit_price_cents"); err != nil {
				return nil, err
			}
			c.Lines = append(c.Lines, CartLine{*l.Presentation, *l.QuantityMilli, *l.UnitPriceCents})
		}
		out.Carts[c.Code] = c
	}
	for _, x := range in.ShoppingList {
		if err := need(x.Product != nil, "shopping_list[].product"); err != nil {
			return nil, err
		}
		out.List = append(out.List, ListItem{*x.Product, x.NeededMilli, x.QuantityMilli, str(x.RecommendedPresentation)})
	}
	for _, x := range in.GroceryPurchases {
		d, err := date(x.Date, "grocery_purchases[].date")
		if err != nil {
			return nil, err
		}
		if err := need(x.Store != nil, "grocery_purchases[].store"); err != nil {
			return nil, err
		}
		out.Purchases = append(out.Purchases, Purchase{d, *x.Store})
	}
	for _, x := range in.Terms {
		if err := need(x.Site != nil, "terms[].site"); err != nil {
			return nil, err
		}
		out.Terms = append(out.Terms, asistida.Term{Site: *x.Site, Verdict: str(x.Verdict), URL: str(x.URL), ReadDate: str(x.ReadDate)})
	}
	return out, nil
}

type priceOut struct {
	PriceCents   int64  `json:"price_cents"`
	ObservedDate string `json:"observed_date"`
	Source       string `json:"source,omitempty"`
}

type quoteOut struct {
	Product      string `json:"product"`
	Store        string `json:"store"`
	Presentation string `json:"presentation"`
	PriceCents   int64  `json:"price_cents"`
	SizeMilli    int64  `json:"size_milli"`
	ObservedDate string `json:"observed_date"`
	IsStale      bool   `json:"is_stale"`
	Recommended  bool   `json:"recommended"`
}

func iso(t time.Time) string { return t.Format(time.DateOnly) }

// Eval calcula una consulta de fixture (§15.2) con las funciones del núcleo. profeco_import y
// ocr_parse_tsv las atiende backend/cmd/motores/s5.go con los subpaquetes profeco y ocr.
func Eval(raw json.RawMessage, fn string, args json.RawMessage) (any, error) {
	in, err := ParseInput(raw)
	if err != nil {
		return nil, err
	}
	b := in.Book
	switch fn {
	case "base_unit_price":
		var a struct {
			Presentation *string `json:"presentation"`
			PriceCents   *int64  `json:"price_cents"`
		}
		if err := strict(args, &a); err != nil {
			return nil, err
		}
		if err := need(a.Presentation != nil && a.PriceCents != nil, "presentation y price_cents"); err != nil {
			return nil, err
		}
		p, ok := b.Catalog.Presentation(*a.Presentation)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownPresentation, *a.Presentation)
		}
		c, per, exact, err := BaseUnitPrice(*a.PriceCents, p.SizeMilli, b.Catalog.unitOf(p))
		if err != nil {
			return nil, err
		}
		return struct {
			BaseUnitPriceCents int64  `json:"base_unit_price_cents"`
			Per                string `json:"per"`
			Exact              Frac   `json:"exact"`
		}{c, per, exact}, nil
	case "freshness":
		var a struct {
			ObservedDate *string `json:"observed_date"`
			Source       *string `json:"source"`
		}
		if err := strict(args, &a); err != nil {
			return nil, err
		}
		od, err := date(a.ObservedDate, "observed_date")
		if err != nil {
			return nil, err
		}
		if err := need(a.Source != nil, "source"); err != nil {
			return nil, err
		}
		w, err := Freshness(od, b.Today, b.Settings, *a.Source)
		if err != nil {
			return nil, err
		}
		return struct {
			AgeDays       int    `json:"age_days"`
			IsStale       bool   `json:"is_stale"`
			Weight        Frac   `json:"weight"`
			WeightDisplay string `json:"weight_display"`
			Competes      bool   `json:"competes"`
		}{w.AgeDays, w.IsStale, w.Frac, w.Display(), w.Competes}, nil
	case "current_price":
		var a struct {
			Presentation *string `json:"presentation"`
			Store        *string `json:"store"`
		}
		if err := strict(args, &a); err != nil {
			return nil, err
		}
		if err := need(a.Presentation != nil && a.Store != nil, "presentation y store"); err != nil {
			return nil, err
		}
		if _, ok := b.Catalog.Presentation(*a.Presentation); !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownPresentation, *a.Presentation)
		}
		if _, ok := b.Store(*a.Store); !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownStore, *a.Store)
		}
		cur := b.CurrentPrice(*a.Presentation, *a.Store, b.Today)
		out := struct {
			Price      *priceOut  `json:"price"`
			References []priceOut `json:"references"`
			Promotion  *priceOut  `json:"promotion"`
		}{References: []priceOut{}}
		if cur.Price != nil {
			out.Price = &priceOut{cur.Price.PriceCents, iso(cur.Price.ObservedDate), cur.Price.Source}
		}
		for _, r := range cur.References {
			out.References = append(out.References, priceOut{r.PriceCents, iso(r.ObservedDate), r.Source})
		}
		if cur.Promotion != nil {
			out.Promotion = &priceOut{PriceCents: cur.Promotion.PriceCents, ObservedDate: iso(cur.Promotion.ObservedDate)}
		}
		return out, nil
	case "quote_for":
		var a struct {
			Product      *string `json:"product"`
			Store        *string `json:"store"`
			Presentation *string `json:"presentation"`
		}
		if err := strict(args, &a); err != nil {
			return nil, err
		}
		if err := need(a.Product != nil, "product"); err != nil {
			return nil, err
		}
		q, ok, err := b.QuoteFor(*a.Product, str(a.Store), str(a.Presentation), b.Today)
		if err != nil {
			return nil, err
		}
		out := struct {
			Quote *quoteOut `json:"quote"`
		}{}
		if ok {
			out.Quote = &quoteOut{q.Product, q.Store, q.Presentation, q.PriceCents, q.SizeMilli,
				iso(q.ObservedDate), q.IsStale, q.Recommended}
		}
		return out, nil
	case "compare_carts":
		var a struct {
			Base   *string  `json:"base"`
			Others []string `json:"others"`
		}
		if err := strict(args, &a); err != nil {
			return nil, err
		}
		if err := need(a.Base != nil, "base"); err != nil {
			return nil, err
		}
		base, ok := in.Carts[*a.Base]
		if !ok {
			return nil, fmt.Errorf("%w: canasta %q", ErrInvalidArgument, *a.Base)
		}
		var others []Cart
		for _, c := range a.Others {
			o, ok := in.Carts[c]
			if !ok {
				return nil, fmt.Errorf("%w: canasta %q", ErrInvalidArgument, c)
			}
			others = append(others, o)
		}
		return b.CompareCarts(base, others)
	case "compare_basket":
		if err := strict(args, &struct{}{}); err != nil {
			return nil, err
		}
		return b.CompareBasket(in.List)
	case "membership":
		var a struct {
			Store   *string `json:"store"`
			Savings *int64  `json:"savings_per_purchase_cents"`
		}
		if err := strict(args, &a); err != nil {
			return nil, err
		}
		if err := need(a.Store != nil, "store"); err != nil {
			return nil, err
		}
		return b.Membership(*a.Store, a.Savings, in.List, in.Purchases)
	case "next_visit":
		var a struct {
			Store *string `json:"store"`
			From  *string `json:"from"`
		}
		if err := strict(args, &a); err != nil {
			return nil, err
		}
		from, err := date(a.From, "from")
		if err != nil {
			return nil, err
		}
		if err := need(a.Store != nil, "store"); err != nil {
			return nil, err
		}
		st, ok := b.Store(*a.Store)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownStore, *a.Store)
		}
		d, err := NextVisit(st, from)
		if err != nil {
			return nil, err
		}
		return struct {
			Date string `json:"date"`
		}{iso(d)}, nil
	case "stale_prices_for_plan":
		var a struct {
			Products       []string          `json:"products"`
			AssignedStores map[string]string `json:"assigned_stores"`
			GroceryDate    *string           `json:"grocery_date"`
		}
		if err := strict(args, &a); err != nil {
			return nil, err
		}
		gd, err := date(a.GroceryDate, "grocery_date")
		if err != nil {
			return nil, err
		}
		return b.StalePricesForPlan(a.Products, a.AssignedStores, gd)
	case "assisted_can_enable":
		var a struct {
			Store *string `json:"store"`
		}
		if err := strict(args, &a); err != nil {
			return nil, err
		}
		if err := need(a.Store != nil, "store"); err != nil {
			return nil, err
		}
		st, ok := b.Store(*a.Store)
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrUnknownStore, *a.Store)
		}
		if err := asistida.CanEnable(asistida.SiteFor(st.Code, st.Name), b.Settings.AssistedReadingEnabled,
			st.AssistedReading, in.Terms); err != nil {
			return nil, err
		}
		return struct {
			CanEnable bool `json:"can_enable"`
		}{true}, nil
	}
	return nil, fmt.Errorf("%w: función %q", ErrInvalidArgument, fn)
}

// rfc3339Offset revisa la zona de un instante que time.Parse(time.RFC3339) ya aceptó: RFC 3339 §5.6 la define como «Z» o
// ±time-hour:time-minute, con hora 00-23 y minuto 00-59, y time.Parse también deja pasar +24:00 y +05:60.
func rfc3339Offset(s string) bool {
	if strings.HasSuffix(s, "Z") {
		return true
	}
	z := s[len(s)-6:] // time.Parse ya exigió «Z» o ±hh:mm con dígitos
	return (z[1]-'0')*10+(z[2]-'0') <= 23 && (z[4]-'0')*10+(z[5]-'0') <= 59
}
