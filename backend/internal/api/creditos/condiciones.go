package creditos

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	mc "finanzas-hogar/internal/motor_cred"
	"finanzas-hogar/internal/motor_cred/comparador"
)

// dbErr traduce las reglas de la base: único → 409 already_exists; CHECK o FK → 422 en el campo dado.
func dbErr(err error, field, what string) error {
	switch {
	case err == nil:
		return nil
	case dominio.IsUnique(err):
		return dominio.Conflict("already_exists", "Ya existe "+what+" con esos datos. Edita el que ya está.")
	case dominio.IsConstraint(err):
		return dominio.Invalid(dominio.FieldIssue{Field: field, Code: "out_of_range", Message: "Los datos de " + what + " no cumplen las reglas. Revisa los campos."})
	}
	return err
}

// audited corre fn y escribe audit_log con el renglón antes y después (CONVENCIONES §8).
func audited(r *http.Request, reg *api.Registry, action, table string, id *int64, fn func(tx *sql.Tx) error) error {
	ctx := r.Context()
	return dominio.WithTx(ctx, reg.Deps().DB, func(tx *sql.Tx) error {
		var before json.RawMessage
		if *id != 0 {
			before, _ = dominio.RowJSON(ctx, tx, table, *id)
		}
		if err := fn(tx); err != nil {
			return err
		}
		var after json.RawMessage
		if action != "delete" {
			after, _ = dominio.RowJSON(ctx, tx, table, *id)
		}
		return dominio.Audit(ctx, tx, api.Actor(r), action, table, *id, before, after, api.RequestID(r))
	})
}

func validateCredit(ctx context.Context, q dominio.Querier, c dominio.CreditRow) error {
	var v api.Validator
	v.Check(strings.TrimSpace(c.Name) != "", "name", "required", "Escribe el nombre del crédito.")
	v.Check(oneOf(c.Kind, "revolving_line", "bank_card", "term_loan", "store_credit"), "kind", "invalid_format", "Elige un tipo de crédito válido.")
	v.Check(oneOf(c.PaymentFrequency, "weekly", "fortnightly", "monthly"), "payment_frequency", "invalid_format", "Elige semanal, quincenal o mensual.")
	v.Check(oneOf(c.AvailabilityReading, "A", "B"), "availability_reading", "invalid_format", "La lectura del disponible es A o B.")
	ok, err := dominio.Exists(ctx, q, "person", c.HolderPersonID)
	if err != nil {
		return err
	}
	v.Check(ok, "holder_person_id", "not_found", "No encontramos a esa persona titular.")
	limitKind := c.Kind == "revolving_line" || c.Kind == "bank_card"
	v.Check(!limitKind || c.CreditLimitCents != nil, "credit_limit_cents", "required", "Las líneas revolventes y las tarjetas necesitan su límite.")
	card := c.Kind == "bank_card"
	v.Check(!card || (c.StatementDay != nil && c.PaymentDueDay != nil), "statement_day", "required", "La tarjeta necesita día de corte y día límite de pago.")
	v.Check(card || (c.StatementDay == nil && c.PaymentDueDay == nil), "statement_day", "not_allowed", "El día de corte y de pago solo aplican a tarjetas.")
	v.Check(!card || c.PaymentFrequency == "monthly", "payment_frequency", "not_allowed", "La tarjeta se paga cada mes.")
	rev := c.Kind == "revolving_line"
	v.Check(rev || c.ReadingAAnchorDate == nil, "reading_a_anchor_date", "not_allowed", "La fecha ancla solo aplica a líneas revolventes.")
	v.Check(!rev || c.AvailabilityReading == "B" || c.ReadingAAnchorDate != nil, "reading_a_anchor_date", "required", "La lectura A necesita su fecha ancla.")
	v.Check((c.AnnualRateNum == nil) == (c.AnnualRateDen == nil), "annual_rate_den", "required", "La tasa anual lleva numerador y denominador juntos.")
	v.Check(c.AnnualRateDen == nil || *c.AnnualRateDen > 0, "annual_rate_den", "out_of_range", "El denominador de la tasa debe ser mayor que cero.")
	k18 := c.Kind == "bank_card" || c.Kind == "store_credit"
	v.Check(k18 || (c.RateBase == nil && c.MinPaymentPctBp == nil && c.MinPaymentFixedCents == nil), "rate_base", "not_allowed",
		"La base de la tasa y el pago mínimo solo aplican a tarjetas y créditos de tienda.")
	v.Check(c.RateBase == nil || oneOf(*c.RateBase, "daily", "monthly"), "rate_base", "invalid_format", "La base de la tasa es diaria o mensual.")
	// Q2-67 opción 1: vat_on_interest_bp con valor exige vat_included = false; la API no cambia el otro campo.
	v.Check(c.VatOnInterestBp == nil || (c.VatIncluded != nil && !*c.VatIncluded), "vat_on_interest_bp", "not_allowed",
		"El IVA aparte solo se captura si la tasa no trae IVA: marca «la tasa no incluye IVA» o deja vacío el IVA.")
	return v.Err()
}

func creditFromCreate(in api.CreditCreate) dominio.CreditRow {
	return dominio.CreditRow{Name: strings.TrimSpace(in.Name), Kind: in.Kind, Lender: in.Lender, HolderPersonID: in.HolderPersonID,
		CreditLimitCents: in.CreditLimitCents, PaymentFrequency: in.PaymentFrequency, FirstPaymentMinDays: in.FirstPaymentMinDays,
		AvailabilityReading: deref(in.AvailabilityReading, "B"), AvailabilityConfirmed: deref(in.AvailabilityConfirmed, false),
		ReadingAAnchorDate: in.ReadingAAnchorDate, StatementDay: in.StatementDay, PaymentDueDay: in.PaymentDueDay,
		AnnualRateNum: in.AnnualRateNum, AnnualRateDen: in.AnnualRateDen, VatOnInterestBp: in.VatOnInterestBp,
		AnnualFeeCents: in.AnnualFeeCents, LateFeeCents: in.LateFeeCents, OpeningFeeCents: in.OpeningFeeCents,
		IsActive: deref(in.IsActive, true), Notes: in.Notes, RateBase: in.RateBase, MinPaymentPctBp: in.MinPaymentPctBp,
		MinPaymentFixedCents: in.MinPaymentFixedCents, AnnualRateDecimals: in.AnnualRateDecimals, VatIncluded: in.VatIncluded}
}

func writeCredit(reg *api.Registry, w http.ResponseWriter, r *http.Request, status int, id int64) error {
	s, err := snap(reg, r)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, status, s.view(s.byID[id]))
}

func createCredit(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.CreditCreate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		c := creditFromCreate(in)
		if err := validateCredit(r.Context(), reg.Deps().DB, c); err != nil {
			return err
		}
		var id int64
		err := audited(r, reg, "insert", "credit", &id, func(tx *sql.Tx) (err error) {
			id, err = dominio.InsertCredit(r.Context(), tx, c)
			return dbErr(err, "kind", "un crédito")
		})
		if err != nil {
			return err
		}
		return writeCredit(reg, w, r, http.StatusCreated, id)
	}
}

func updateCredit(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "credit_id")
		if err != nil {
			return err
		}
		var in api.CreditUpdate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		c, err := dominio.GetCredit(r.Context(), reg.Deps().DB, id)
		if err != nil {
			return notFound(err, msgCredit)
		}
		set(&c.Name, in.Name)
		c.Name = strings.TrimSpace(c.Name)
		set(&c.Kind, in.Kind)
		set(&c.HolderPersonID, in.HolderPersonID)
		set(&c.PaymentFrequency, in.PaymentFrequency)
		set(&c.AvailabilityReading, in.AvailabilityReading)
		set(&c.AvailabilityConfirmed, in.AvailabilityConfirmed)
		set(&c.IsActive, in.IsActive)
		nullable(&c.Lender, in.Lender)
		nullable(&c.CreditLimitCents, in.CreditLimitCents)
		nullable(&c.FirstPaymentMinDays, in.FirstPaymentMinDays)
		nullable(&c.ReadingAAnchorDate, in.ReadingAAnchorDate)
		nullable(&c.StatementDay, in.StatementDay)
		nullable(&c.PaymentDueDay, in.PaymentDueDay)
		nullable(&c.AnnualRateNum, in.AnnualRateNum)
		nullable(&c.AnnualRateDen, in.AnnualRateDen)
		nullable(&c.VatOnInterestBp, in.VatOnInterestBp)
		nullable(&c.AnnualFeeCents, in.AnnualFeeCents)
		nullable(&c.LateFeeCents, in.LateFeeCents)
		nullable(&c.OpeningFeeCents, in.OpeningFeeCents)
		nullable(&c.Notes, in.Notes)
		nullable(&c.RateBase, in.RateBase)
		nullable(&c.MinPaymentPctBp, in.MinPaymentPctBp)
		nullable(&c.MinPaymentFixedCents, in.MinPaymentFixedCents)
		nullable(&c.AnnualRateDecimals, in.AnnualRateDecimals)
		nullable(&c.VatIncluded, in.VatIncluded)
		if err := validateCredit(r.Context(), reg.Deps().DB, c); err != nil {
			return err
		}
		if err := audited(r, reg, "update", "credit", &id, func(tx *sql.Tx) error {
			return dbErr(dominio.UpdateCredit(r.Context(), tx, c), "kind", "un crédito")
		}); err != nil {
			return err
		}
		return writeCredit(reg, w, r, http.StatusOK, id)
	}
}

func termRateView(t dominio.TermRateRow) api.CreditTermRate {
	return api.CreditTermRate{ID: t.ID, CreditID: t.CreditID, TermCount: t.TermCount, RateNum: t.RateNum, RateDen: t.RateDen,
		Source: t.Source, QuoteDate: t.QuoteDate, Notes: t.Notes, CreatedAt: api.APITime(t.CreatedAt), UpdatedAt: api.APITime(t.UpdatedAt),
		RateDecimals: t.RateDecimals, RatePct: ratePct(t.RateNum, t.RateDen, t.RateDecimals)}
}

func validateTermRate(t dominio.TermRateRow) error {
	var v api.Validator
	v.Check(t.TermCount >= 1 && t.TermCount <= 120, "term_count", "out_of_range", "El plazo va de 1 a 120 pagos.")
	v.Check(t.RateNum >= 0, "rate_num", "out_of_range", "El numerador de la tasa no puede ser negativo.")
	v.Check(t.RateDen > 0, "rate_den", "out_of_range", "El denominador de la tasa debe ser mayor que cero.")
	v.Check(oneOf(t.Source, "quote", "manual"), "source", "invalid_format", "El origen es cotización o captura manual.")
	v.Check(t.Source != "quote" || t.QuoteDate != nil, "quote_date", "required", "Una tasa de cotización necesita la fecha de la cotización.")
	v.Check(t.RateDecimals == nil || *t.RateDecimals >= 0, "rate_decimals", "out_of_range", "Los decimales no pueden ser negativos.")
	return v.Err()
}

func listCreditTermRates(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "credit_id")
		if err != nil {
			return err
		}
		p, err := api.ParsePage(r, []string{"term_count"}, "term_count")
		if err != nil {
			return err
		}
		if ok, err := dominio.Exists(r.Context(), reg.Deps().DB, "credit", id); err != nil || !ok {
			return notFound(firstErr(err, sql.ErrNoRows), msgCredit)
		}
		ts, err := dominio.QueryTermRates(r.Context(), reg.Deps().DB, "credit_id = ?",
			p.OrderBy(map[string]string{"term_count": "term_count", "id": "id"})+" "+p.LimitOffset(), id)
		if err != nil {
			return err
		}
		items := make([]api.CreditTermRate, 0, len(ts))
		for _, t := range ts {
			items = append(items, termRateView(t))
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

func firstErr(a, b error) error {
	if a != nil {
		return a
	}
	return b
}

func createCreditTermRate(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		cid, err := api.PathID(r, "credit_id")
		if err != nil {
			return err
		}
		var in api.CreditTermRateCreate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		if ok, err := dominio.Exists(r.Context(), reg.Deps().DB, "credit", cid); err != nil || !ok {
			return notFound(firstErr(err, sql.ErrNoRows), msgCredit)
		}
		t := dominio.TermRateRow{CreditID: cid, TermCount: in.TermCount, RateNum: in.RateNum, RateDen: in.RateDen, Source: in.Source,
			QuoteDate: in.QuoteDate, Notes: in.Notes, RateDecimals: in.RateDecimals}
		if err := validateTermRate(t); err != nil {
			return err
		}
		var id int64
		if err := audited(r, reg, "insert", "credit_term_rate", &id, func(tx *sql.Tx) (err error) {
			id, err = dominio.InsertTermRate(r.Context(), tx, t)
			return dbErr(err, "term_count", "una tasa para ese plazo")
		}); err != nil {
			return err
		}
		t, err = dominio.GetTermRate(r.Context(), reg.Deps().DB, id)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusCreated, termRateView(t))
	}
}

const msgTermRate = "No encontramos esa tasa por plazo."

func updateCreditTermRate(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "credit_term_rate_id")
		if err != nil {
			return err
		}
		var in api.CreditTermRateUpdate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		t, err := dominio.GetTermRate(r.Context(), reg.Deps().DB, id)
		if err != nil {
			return notFound(err, msgTermRate)
		}
		set(&t.TermCount, in.TermCount)
		set(&t.RateNum, in.RateNum)
		set(&t.RateDen, in.RateDen)
		set(&t.Source, in.Source)
		nullable(&t.QuoteDate, in.QuoteDate)
		nullable(&t.Notes, in.Notes)
		nullable(&t.RateDecimals, in.RateDecimals)
		if err := validateTermRate(t); err != nil {
			return err
		}
		if err := audited(r, reg, "update", "credit_term_rate", &id, func(tx *sql.Tx) error {
			return dbErr(dominio.UpdateTermRate(r.Context(), tx, t), "term_count", "una tasa para ese plazo")
		}); err != nil {
			return err
		}
		t, err = dominio.GetTermRate(r.Context(), reg.Deps().DB, id)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, termRateView(t))
	}
}

func deleteCreditTermRate(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "credit_term_rate_id")
		if err != nil {
			return err
		}
		if err := audited(r, reg, "delete", "credit_term_rate", &id, func(tx *sql.Tx) error {
			return dominio.DeleteRow(r.Context(), tx, "credit_term_rate", id)
		}); err != nil {
			return notFound(err, msgTermRate)
		}
		return api.NoContent(w)
	}
}

func promoView(p dominio.PromotionRow, today string) api.Promotion {
	return api.Promotion{ID: p.ID, CreditID: p.CreditID, MerchantName: p.MerchantName, StoreID: p.StoreID, TermCount: p.TermCount,
		RateNum: p.RateNum, RateDen: p.RateDen, MinPurchaseCents: p.MinPurchaseCents, ValidFrom: p.ValidFrom, ValidTo: p.ValidTo,
		SourceNote: p.SourceNote, IsActive: p.IsActive, CreatedAt: api.APITime(p.CreatedAt), UpdatedAt: api.APITime(p.UpdatedAt),
		RateDecimals: p.RateDecimals, IsCurrent: api.Ptr(p.IsActive && p.ValidFrom <= today && today <= p.ValidTo),
		RatePct: ratePct(p.RateNum, p.RateDen, p.RateDecimals)}
}

func validatePromo(ctx context.Context, q dominio.Querier, p dominio.PromotionRow) error {
	var v api.Validator
	v.Check(strings.TrimSpace(p.MerchantName) != "", "merchant_name", "required", "Escribe el comercio de la promoción.")
	v.Check(p.TermCount >= 1 && p.TermCount <= 120, "term_count", "out_of_range", "El plazo va de 1 a 120 pagos.")
	v.Check(p.RateNum >= 0 && p.RateDen > 0, "rate_den", "out_of_range", "La tasa necesita numerador de 0 o más y denominador mayor que cero.")
	v.Check(p.MinPurchaseCents >= 0, "min_purchase_cents", "out_of_range", "La compra mínima no puede ser negativa.")
	v.Check(p.ValidTo >= p.ValidFrom, "valid_to", "out_of_range", "El fin de la vigencia no puede ser antes del inicio.")
	for _, ref := range []struct {
		id          *int64
		table, f, m string
	}{{p.CreditID, "credit", "credit_id", "No encontramos ese crédito."}, {p.StoreID, "store", "store_id", "No encontramos esa tienda."}} {
		if ref.id != nil {
			ok, err := dominio.Exists(ctx, q, ref.table, *ref.id)
			if err != nil {
				return err
			}
			v.Check(ok, ref.f, "not_found", ref.m)
		}
	}
	return v.Err()
}

func listPromotions(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"valid_from"}, "valid_from")
		if err != nil {
			return err
		}
		where, args := []string{"1 = 1"}, []any{}
		for _, f := range []string{"credit_id", "store_id"} {
			ids, err := api.QueryIDs(r, f)
			if err != nil {
				return err
			}
			if len(ids) > 0 {
				where = append(where, f+" IN ("+strings.TrimSuffix(strings.Repeat("?, ", len(ids)), ", ")+")")
				for _, id := range ids {
					args = append(args, id)
				}
			}
		}
		if m := strings.TrimSpace(r.URL.Query().Get("merchant_name")); m != "" {
			where, args = append(where, "merchant_name = ? COLLATE NOCASE"), append(args, m)
		}
		per, err := api.QueryPeriod(r)
		if err != nil {
			return err
		}
		on, err := api.QueryDate(r, "valid_on")
		if err != nil {
			return err
		}
		if per != nil {
			where, args = append(where, "valid_from <= ? AND valid_to >= ?"), append(args, per.End.String(), per.Start.String())
		}
		if on != nil {
			where, args = append(where, "valid_from <= ? AND valid_to >= ?"), append(args, on.String(), on.String())
		}
		t := today(reg).String()
		if per == nil && on == nil { // por defecto: vigentes y futuras
			where, args = append(where, "valid_to >= ?"), append(args, t)
		}
		ps, err := dominio.QueryPromotions(r.Context(), reg.Deps().DB, strings.Join(where, " AND "),
			p.OrderBy(map[string]string{"valid_from": "valid_from", "id": "id"})+" "+p.LimitOffset(), args...)
		if err != nil {
			return err
		}
		items := make([]api.Promotion, 0, len(ps))
		for _, x := range ps {
			items = append(items, promoView(x, t))
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

func writePromo(reg *api.Registry, w http.ResponseWriter, r *http.Request, status int, id int64) error {
	p, err := dominio.GetPromotion(r.Context(), reg.Deps().DB, id)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, status, promoView(p, today(reg).String()))
}

func createPromotion(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.PromotionCreate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		p := dominio.PromotionRow{CreditID: in.CreditID, MerchantName: strings.TrimSpace(in.MerchantName), StoreID: in.StoreID,
			TermCount: in.TermCount, RateNum: deref(in.RateNum, 0), RateDen: deref(in.RateDen, 1), MinPurchaseCents: deref(in.MinPurchaseCents, 0),
			ValidFrom: in.ValidFrom, ValidTo: in.ValidTo, SourceNote: in.SourceNote, IsActive: deref(in.IsActive, true), RateDecimals: in.RateDecimals}
		if err := validatePromo(r.Context(), reg.Deps().DB, p); err != nil {
			return err
		}
		var id int64
		if err := audited(r, reg, "insert", "promotion", &id, func(tx *sql.Tx) (err error) {
			id, err = dominio.InsertPromotion(r.Context(), tx, p)
			return dbErr(err, "merchant_name", "una promoción")
		}); err != nil {
			return err
		}
		return writePromo(reg, w, r, http.StatusCreated, id)
	}
}

const msgPromo = "No encontramos esa promoción."

func updatePromotion(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "promotion_id")
		if err != nil {
			return err
		}
		var in api.PromotionUpdate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		p, err := dominio.GetPromotion(r.Context(), reg.Deps().DB, id)
		if err != nil {
			return notFound(err, msgPromo)
		}
		nullable(&p.CreditID, in.CreditID)
		set(&p.MerchantName, in.MerchantName)
		p.MerchantName = strings.TrimSpace(p.MerchantName)
		nullable(&p.StoreID, in.StoreID)
		set(&p.TermCount, in.TermCount)
		set(&p.RateNum, in.RateNum)
		set(&p.RateDen, in.RateDen)
		set(&p.MinPurchaseCents, in.MinPurchaseCents)
		set(&p.ValidFrom, in.ValidFrom)
		set(&p.ValidTo, in.ValidTo)
		nullable(&p.SourceNote, in.SourceNote)
		set(&p.IsActive, in.IsActive)
		nullable(&p.RateDecimals, in.RateDecimals)
		if err := validatePromo(r.Context(), reg.Deps().DB, p); err != nil {
			return err
		}
		if err := audited(r, reg, "update", "promotion", &id, func(tx *sql.Tx) error {
			return dbErr(dominio.UpdatePromotion(r.Context(), tx, p), "merchant_name", "una promoción")
		}); err != nil {
			return err
		}
		return writePromo(reg, w, r, http.StatusOK, id)
	}
}

func deletePromotion(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "promotion_id")
		if err != nil {
			return err
		}
		if _, err := dominio.GetPromotion(r.Context(), reg.Deps().DB, id); err != nil {
			return notFound(err, msgPromo)
		}
		n, err := dominio.CountWhere(r.Context(), reg.Deps().DB, "credit_purchase", "promotion_id = ?", id)
		if err != nil {
			return err
		}
		if n > 0 {
			return dominio.Conflict("in_use", "Hay compras que usan esta promoción. Desactívala en lugar de eliminarla.")
		}
		if err := audited(r, reg, "delete", "promotion", &id, func(tx *sql.Tx) error {
			return dominio.DeleteRow(r.Context(), tx, "promotion", id)
		}); err != nil {
			return notFound(err, msgPromo)
		}
		return api.NoContent(w)
	}
}

// Valores del comparador (K16, K43, K44).
var optimizerParams = []string{"comparator.timeout_seconds", "optimizer.date_window_days"}

type constraints struct {
	buffer, maxLive *int64
	zeroDebt        *string
	weights         api.OptimizerWeights
}

func readConstraints(ctx context.Context, q dominio.Querier) (constraints, error) {
	var c constraints
	for key, dst := range map[string]any{"optimizer.min_available_buffer_cents": &c.buffer, "optimizer.max_live_purchases": &c.maxLive,
		"optimizer.zero_debt_by_date": &c.zeroDebt, "optimizer.weights": &c.weights} {
		raw, found, err := dominio.SettingRaw(ctx, q, key)
		if err != nil {
			return c, err
		}
		if found && len(raw) > 0 {
			if err := json.Unmarshal(raw, dst); err != nil {
				return c, err
			}
		}
	}
	return c, nil
}

func intParam(ctx context.Context, q dominio.Querier, key string) (api.SettingParameterInteger, error) {
	p, err := dominio.ReadParameter(ctx, q, key)
	if err != nil {
		return api.SettingParameterInteger{}, err
	}
	out := api.SettingParameterInteger{Key: p.Key, ValueType: p.ValueType, Unit: p.Unit, IsInitial: p.IsInitial, Value: api.Null[int64](), InitialValue: api.Null[int64]()}
	if len(p.Value) > 0 {
		_ = json.Unmarshal(p.Value, &out.Value)
	}
	if len(p.InitialValue) > 0 {
		_ = json.Unmarshal(p.InitialValue, &out.InitialValue)
	}
	return out, nil
}

func getOptimizerSettings(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return writeOptimizer(reg, w, r)
	}
}

func writeOptimizer(reg *api.Registry, w http.ResponseWriter, r *http.Request) error {
	ctx, q := r.Context(), reg.Deps().DB
	c, err := readConstraints(ctx, q)
	if err != nil {
		return err
	}
	out := api.OptimizerSettings{MissingConstraints: []string{}, MinAvailableBufferCents: c.buffer, MaxLivePurchases: c.maxLive,
		ZeroDebtByDate: c.zeroDebt, Weights: c.weights}
	for _, m := range []struct {
		missing bool
		name    string
	}{{c.buffer == nil, "min_available_buffer_cents"}, {c.maxLive == nil, "max_live_purchases"}, {c.zeroDebt == nil, "zero_debt_by_date"}} {
		if m.missing {
			out.MissingConstraints = append(out.MissingConstraints, m.name)
		}
	}
	if out.Parameters.ComparatorTimeoutSeconds, err = intParam(ctx, q, optimizerParams[0]); err != nil {
		return err
	}
	if out.Parameters.OptimizerDateWindowDays, err = intParam(ctx, q, optimizerParams[1]); err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

// optimizerPatch es OptimizerSettingsUpdate con null distinguible: null o vacío responde 422 (K44).
type optimizerPatch struct {
	MinAvailableBufferCents api.Nullable[int64]                `json:"min_available_buffer_cents,omitzero"`
	MaxLivePurchases        api.Nullable[int64]                `json:"max_live_purchases,omitzero"`
	ZeroDebtByDate          api.Nullable[string]               `json:"zero_debt_by_date,omitzero"`
	Weights                 api.Nullable[api.OptimizerWeights] `json:"weights,omitzero"`
	Parameters              map[string]json.RawMessage         `json:"parameters,omitzero"`
	ResetParameters         []string                           `json:"reset_parameters,omitzero"`
}

func updateOptimizerSettings(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in optimizerPatch
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		var v api.Validator
		for _, n := range []struct {
			set, valid bool
			f          string
		}{{in.MinAvailableBufferCents.Set, in.MinAvailableBufferCents.Valid, "min_available_buffer_cents"},
			{in.MaxLivePurchases.Set, in.MaxLivePurchases.Valid, "max_live_purchases"},
			{in.ZeroDebtByDate.Set, in.ZeroDebtByDate.Valid, "zero_debt_by_date"}, {in.Weights.Set, in.Weights.Valid, "weights"}} {
			v.Check(!n.set || n.valid, n.f, "required", "Este valor no puede quedar vacío; escribe uno.")
		}
		v.Check(!in.MinAvailableBufferCents.Valid || in.MinAvailableBufferCents.Value >= 0, "min_available_buffer_cents", "out_of_range", "El colchón no puede ser negativo.")
		v.Check(!in.MaxLivePurchases.Valid || in.MaxLivePurchases.Value >= 1, "max_live_purchases", "out_of_range", "El máximo de compras vivas es 1 o más.")
		if in.ZeroDebtByDate.Valid {
			if _, err := v.Date(in.ZeroDebtByDate.Value, "zero_debt_by_date"); err != nil {
				return err
			}
		}
		if wt := in.Weights.Value; in.Weights.Valid {
			v.Check(wt.LowestInterest >= 0 && wt.EvenLowPayment >= 0 && wt.MaxAvailable >= 0 && wt.MostValuePerPeso >= 0 &&
				wt.LowestInterest+wt.EvenLowPayment+wt.MaxAvailable+wt.MostValuePerPeso > 0, "weights", "out_of_range",
				"Los pesos son enteros de 0 o más y al menos uno mayor que cero.")
		}
		if err := v.Err(); err != nil {
			return err
		}
		ctx := r.Context()
		err := dominio.WithTx(ctx, reg.Deps().DB, func(tx *sql.Tx) error {
			for key, val := range map[string]struct {
				ok bool
				v  any
			}{"optimizer.min_available_buffer_cents": {in.MinAvailableBufferCents.Valid, in.MinAvailableBufferCents.Value},
				"optimizer.max_live_purchases": {in.MaxLivePurchases.Valid, in.MaxLivePurchases.Value},
				"optimizer.zero_debt_by_date":  {in.ZeroDebtByDate.Valid, in.ZeroDebtByDate.Value},
				"optimizer.weights":            {in.Weights.Valid, in.Weights.Value}} {
				if !val.ok {
					continue
				}
				raw, _ := json.Marshal(val.v)
				if err := dominio.SetSetting(ctx, tx, key, raw, api.Actor(r), api.RequestID(r)); err != nil {
					return err
				}
			}
			if len(in.Parameters) == 0 && len(in.ResetParameters) == 0 {
				return nil
			}
			return dominio.ApplyParams(ctx, tx, dominio.ParamChange{Set: in.Parameters, Reset: in.ResetParameters}, optimizerParams, api.Actor(r), api.RequestID(r))
		})
		if err != nil {
			return err
		}
		return writeOptimizer(reg, w, r)
	}
}

func compareScenarios(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in api.ScenarioRequest
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		var v api.Validator
		v.PositiveCents(in.AmountCents, "amount_cents")
		desired, err := v.Date(in.DesiredDate, "desired_date")
		if err != nil {
			return err
		}
		ctx, q := r.Context(), reg.Deps().DB
		s, err := snap(reg, r)
		if err != nil {
			return err
		}
		shares := checkShares(&v, s, in.AmountCents, in.Shares)
		if err := v.Err(); err != nil {
			return err
		}
		if shares == nil { // sin partes: todo al titular (como en CreditPurchaseCreate); PROVISIONAL S6-creditos
			for _, p := range s.persons {
				if p.Participation == "holder" && (shares == nil || p.ID < shares[0].PersonID) {
					shares = []mc.ShareInput{s.shareInput(p.ID, in.AmountCents)}
				}
			}
		}
		c, err := readConstraints(ctx, q)
		if err != nil {
			return err
		}
		window, err := dominio.SettingInt(ctx, q, "optimizer.date_window_days")
		if err != nil {
			return err
		}
		secs, err := dominio.SettingInt(ctx, q, "comparator.timeout_seconds")
		if err != nil {
			return err
		}
		timeout, err := comparador.TimeoutFromSetting(secs)
		if err != nil {
			return motorErr(err, "comparator.timeout_seconds")
		}
		req := comparador.Request{Today: mdateOf(s.today), Amount: mc.Cents(in.AmountCents), DesiredDate: mdateOf(desired),
			StoreID: in.StoreID, Merchant: deref(in.MerchantName, ""), Shares: shares, Credits: s.states(nil),
			Global: comparador.Constraints{MinAvailableBuffer: centsPtr(c.buffer), MaxLivePurchases: intPtr(c.maxLive),
				ZeroDebtBy: mdatePtr(c.zeroDebt), DateWindowDays: intPtr(&window)},
			Override: comparador.Constraints{MinAvailableBuffer: centsPtr(in.MinAvailableBufferCents), MaxLivePurchases: intPtr(in.MaxLivePurchases),
				ZeroDebtBy: mdatePtr(in.ZeroDebtByDate), DateWindowDays: intPtr(in.DateWindowDays)},
			Weights: comparador.Weights(deref(in.Weights, c.weights))}
		for _, p := range s.persons {
			if p.CapCents != nil {
				req.Caps = append(req.Caps, comparador.PersonCap{Person: p.Code, PersonID: p.ID, Cap: mc.Cents(*p.CapCents)})
			}
		}
		cctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		res, err := comparador.Compare(cctx, req)
		if err != nil {
			return motorErr(err, "amount_cents")
		}
		return api.WriteJSON(w, http.StatusOK, scenarioResult(res))
	}
}

func scenario(x comparador.Scenario) api.Scenario {
	return api.Scenario{CreditID: x.CreditID, TermCount: int64(x.TermCount), PurchaseDate: x.PurchaseDate.String(), StoreID: x.StoreID,
		PromotionID: x.PromotionID, InterestCents: int64(x.Interest), MaxFortnightPaymentCents: int64(x.MaxFortnightPayment),
		FortnightPaymentSpreadCents: int64(x.FortnightPaymentSpread), MinAvailableCents: int64(x.MinAvailable),
		MaxLivePurchases: int64(x.MaxLivePurchases), ZeroDebtDate: x.ZeroDebtDate.String(), ValuePerPesoMilli: x.ValuePerPesoMilli, ScoreBp: x.ScoreBP}
}

func i64(p *int) *int64 {
	if p == nil {
		return nil
	}
	return api.Ptr(int64(*p))
}

func scenarioResult(res comparador.Result) api.ScenarioResult {
	out := api.ScenarioResult{Status: res.Status, MissingConstraints: append([]string{}, res.MissingConstraints...),
		EvaluatedCount: int64(res.EvaluatedCount), Valid: []api.Scenario{}, Rejected: []api.ScenarioResultRejectedItem{}}
	cu := res.ConstraintsUsed
	out.ConstraintsUsed = &api.ScenarioResultConstraintsUsed{MaxLivePurchases: i64(cu.MaxLivePurchases), DateWindowDays: i64(cu.DateWindowDays)}
	if cu.MinAvailableBuffer != nil {
		out.ConstraintsUsed.MinAvailableBufferCents = api.Ptr(int64(*cu.MinAvailableBuffer))
	}
	if cu.ZeroDebtBy != nil {
		out.ConstraintsUsed.ZeroDebtByDate = api.Ptr(cu.ZeroDebtBy.String())
	}
	for _, x := range res.Valid {
		out.Valid = append(out.Valid, scenario(x))
	}
	for _, x := range res.Rejected {
		out.Rejected = append(out.Rejected, api.ScenarioResultRejectedItem{Scenario: scenario(x.Scenario), ViolatedConstraint: x.ViolatedConstraint})
	}
	if o := res.Optima; o != nil {
		a, b, c, d := scenario(o.LowestInterest), scenario(o.EvenLowPayment), scenario(o.MaxAvailable), scenario(o.MostValuePerPeso)
		out.Optima = api.ScenarioResultOptima{LowestInterest: &a, EvenLowPayment: &b, MaxAvailable: &c, MostValuePerPeso: &d}
	}
	if res.Recommendation != nil {
		rec := scenario(*res.Recommendation)
		out.Recommendation = &rec
	}
	return out
}
