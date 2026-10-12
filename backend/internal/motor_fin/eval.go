package motor_fin

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ---------- Vistas del oráculo (§3.14): claves literales de model.json, en centavos ----------

func (in *Input) shape() (string, error) {
	ext := ""
	n := 0
	for _, p := range in.Persons {
		if p.Participation == "external" {
			ext, n = p.Code, n+1
		}
	}
	if n != 1 {
		return "", errf(CodeOracleViewShape, "las vistas exigen un titular y una sola persona externa")
	}
	return ext, nil
}

func (in *Input) viewPeriod(m Month) (OMap, error) {
	ext, err := in.shape()
	if err != nil {
		return nil, err
	}
	ps, _ := records(in).periods(PeriodSpec{Months: []Month{m}})
	x, h := ps[0], in.holder
	det := []OMap{}
	for _, py := range x.Cash.Payments {
		det = append(det, OMap{{"despensa", py.PurchaseSeq}, {"fecha", py.DueDate}, {"monto", py.AmountCents},
			{h, py.Shares.Get(h)}, {ext, py.Shares.Get(ext)}})
	}
	dev := OMap{{"despensa_total", x.Accrued.ByCategory.Get(in.Params.GroceriesCategory)}, {"papel_contado", x.Accrued.paper},
		{"total", x.Accrued.TotalCents}, {h, x.Accrued.ByPerson.Get(h)}, {ext, x.Accrued.ByPerson.Get(ext)}}
	efe := OMap{{"pagos", x.Cash.PaymentsCents}, {"papel", x.Cash.paper}, {"total", x.Cash.TotalCents},
		{h, x.Cash.ByPerson.Get(h)}, {ext, x.Cash.ByPerson.Get(ext)}, {"n_pagos", x.Cash.PaymentCount}, {"detalle_pagos", det}}
	return OMap{{"etiqueta", x.Label}, {"inicio", x.Start}, {"fin", x.End}, {"despensa_que_cubre", x.CoveringGroceriesSeq},
		{"devengado", dev}, {"efectivo", efe}, {"diferencia_acumulada", x.CumulativeDiffCents}}, nil
}

func seqKey(seq *int, id int64) string {
	if seq == nil {
		return fmt.Sprintf("compra_%d", id)
	}
	return fmt.Sprintf("despensa_%d", *seq)
}

func (in *Input) viewHorizon(spec HorizonSpec) (OMap, error) {
	if _, err := in.shape(); err != nil {
		return nil, err
	}
	h, err := Horizon(*in, spec)
	if err != nil {
		return nil, err
	}
	var o OMap
	for _, x := range h.PaymentsByPurchase {
		k := "pagos_" + seqKey(x.PurchaseSeq, x.PurchaseID)
		if !x.Complete {
			k += "_en_horizonte"
		}
		o = append(o, KV{k, x.InHorizonCents})
	}
	det := OMap{}
	for _, x := range h.LiveDebt.ByPurchase {
		k := seqKey(x.PurchaseSeq, x.PurchaseID)
		if x.Started {
			k += "_pendiente"
		}
		det = append(det, KV{k, x.PendingCents})
	}
	suf := fmt.Sprintf("%d_%s", spec.End.Day, monthAbbr[spec.End.Month-1])
	return append(o, KV{"papel_contado", h.PaperCashCents}, KV{"total_efectivo", h.CashTotalCents},
		KV{"suma_periodos_efectivo", h.SumOfPeriodsCashCents}, KV{"total_devengado", h.AccruedTotalCents},
		KV{"deuda_viva_" + suf, h.LiveDebt.LiveDebtCents}, KV{"deuda_viva_detalle", det},
		KV{"promedio_mensual_efectivo", h.MonthlyAvgCashCents}, KV{"promedio_mensual_devengado", h.MonthlyAvgAccruedCents},
		KV{"diferencia_final", h.FinalDiffCents}), nil
}

func (in *Input) viewStable(items []StableItem) (OMap, error) {
	ext, err := in.shape()
	if err != nil {
		return nil, err
	}
	sv, err := StableValue(items, in.Persons)
	if err != nil {
		return nil, err
	}
	var desp, pap *StableComponent
	for i := range sv.Components {
		switch sv.Components[i].Label {
		case "despensa":
			desp = &sv.Components[i]
		case "papel_contado":
			pap = &sv.Components[i]
		}
	}
	if desp == nil || pap == nil {
		return nil, errf(CodeOracleViewShape, "faltan las partidas despensa y papel_contado")
	}
	h := in.holder
	f := sv.Fortnight
	return OMap{
		{"mes", OMap{{"total", sv.Month.TotalCents}, {h, sv.Month.ByPerson.Get(h)}, {ext, sv.Month.ByPerson.Get(ext)}}},
		{"quincena", OMap{{"total", f.TotalCents}, {h, f.ByPerson.Get(h)}, {ext, f.ByPerson.Get(ext)},
			{"suma_" + h + "_" + ext, f.SumPersonsCents}, {"diferencia_centavo", f.CentDiffCents}}},
		{"componentes", OMap{{"mitad_despensa", desp.MonthlyCents}, {"papel_contado_mes", pap.MonthlyCents},
			{ext + "_mitad_despensa", desp.ByPerson.Get(ext)}, {h + "_mitad_despensa", desp.ByPerson.Get(h)}}},
	}, nil
}

func (in *Input) viewPaper() ([]OMap, error) {
	if _, err := in.shape(); err != nil {
		return nil, err
	}
	pk, err := PaperPlan(*in)
	if err != nil {
		return nil, err
	}
	out := []OMap{}
	for _, x := range pk {
		forma := "de contado"
		if x.Form == "in_groceries" {
			forma = fmt.Sprintf("dentro de despensa %d", *x.PurchaseSeq)
		}
		out = append(out, OMap{{"paquete", x.Package}, {"compra", x.Date}, {"quien", x.PersonName}, {"forma", forma}, {"monto", x.AmountCents}})
	}
	return out, nil
}

func (in *Input) viewSensitivity(items []StableItem) (OMap, error) {
	ext, err := in.shape()
	if err != nil {
		return nil, err
	}
	sv, err := StableValue(items, in.Persons)
	if err != nil {
		return nil, err
	}
	return OMap{{"costo_mensual_total", sv.Month.TotalCents}, {in.holder + "_mensual", sv.Month.ByPerson.Get(in.holder)},
		{ext + "_mensual", sv.Month.ByPerson.Get(ext)}}, nil
}

func (in *Input) viewInterest(id int64) (OMap, error) {
	if _, err := in.shape(); err != nil {
		return nil, err
	}
	l := records(in)
	var cp *CreditPurchase
	for _, x := range l.purchs {
		if x.ID == id {
			cp = x
		}
	}
	if cp == nil {
		return nil, errf(CodeInvalidArgument, "purchase_id %d no está en credit_purchases", id)
	}
	var acc, cash Cents
	for i := range l.costs {
		if c := &l.costs[i]; c.purchase == id {
			for k := 0; k < c.months; k++ {
				t, _ := c.slot(k, in.holder)
				acc += t
			}
		}
	}
	var first any
	for _, py := range cp.Payments {
		cash += py.AmountCents
		if py.Seq == 1 {
			first = py.AmountCents
		}
	}
	if acc != cash {
		return nil, errf(CodeSumMismatch, "compra %d: devengado %d ≠ efectivo %d", id, acc, cash)
	}
	return OMap{{"monto", cp.PrincipalCents}, {"intereses", acc - cp.PrincipalCents}, {"total", acc}, {"pago_quincenal", first}}, nil
}

// ---------- Evaluador de consultas (§3, FORMATO §2) ----------

// Calendar es la salida de la consulta calendario (§3.1).
type Calendar struct {
	Today          Date      `json:"today"`
	Month          Month     `json:"month"`
	MonthStart     Date      `json:"month_start"`
	MonthEnd       Date      `json:"month_end"`
	MonthLabel     string    `json:"month_label"`
	DaysInMonth    int       `json:"days_in_month"`
	Fortnight      Fortnight `json:"fortnight"`
	FortnightStart Date      `json:"fortnight_start"`
	FortnightEnd   Date      `json:"fortnight_end"`
	PayDate        Date      `json:"pay_date"`
	Week           Week      `json:"week"`
	WeekStart      Date      `json:"week_start"`
	WeekEnd        Date      `json:"week_end"`
}

// CalendarOf arma la salida de calendario para una fecha.
func CalendarOf(d Date) Calendar {
	m := d.MonthOf()
	f := FortnightOf(d)
	fs, fe, pay := FortnightBounds(f)
	w := WeekOf(d)
	ws, we := WeekBounds(w)
	return Calendar{d, m, m.First(), m.Last(), m.Label(), m.Days(), f, fs, fe, pay, w, ws, we}
}

func asErr(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return errf(CodeInvalidArgument, "%v", err)
}

func missing(fn string) error {
	return errf(CodeInvalidArgument, "%s: faltan argumentos obligatorios o sobran alternativas", fn)
}

// Eval evalúa una consulta (funcion y argumentos) sobre las entradas JSON.
func Eval(input []byte, fn string, args []byte) (any, *Error) {
	in, err := DecodeInput(input)
	if err != nil {
		return nil, asErr(err)
	}
	out, err := dispatch(&in, fn, args)
	if err != nil {
		return nil, asErr(err)
	}
	return out, nil
}

type monthArg struct {
	Month *Month `json:"month"`
}

type horizonArg struct {
	Start  *Date    `json:"start"`
	End    *Date    `json:"end"`
	Months *[]Month `json:"months"`
}

func (a horizonArg) spec(fn string) (HorizonSpec, error) {
	if a.Start == nil || a.End == nil || a.Months == nil || len(*a.Months) == 0 {
		return HorizonSpec{}, missing(fn)
	}
	return HorizonSpec{*a.Start, *a.End, *a.Months}, nil
}

type itemsArg struct {
	Items *[]StableItem `json:"items"`
}

func dispatch(in *Input, fn string, args []byte) (any, error) {
	dec := func(v any) error { return decodeStrict(args, v) }
	switch fn {
	case "calendario":
		var a struct {
			Date *Date      `json:"date"`
			Now  *time.Time `json:"now"`
		}
		if err := dec(&a); err != nil {
			return nil, err
		}
		if (a.Date == nil) == (a.Now == nil) {
			return nil, missing(fn)
		}
		if a.Now != nil {
			d, err := Today(*a.Now, in.Params.Zone)
			if err != nil {
				return nil, err
			}
			return CalendarOf(d), nil
		}
		return CalendarOf(*a.Date), nil
	case "periodos":
		var a struct {
			Months  *[]Month  `json:"months"`
			Windows *[]Window `json:"windows"`
		}
		if err := dec(&a); err != nil {
			return nil, err
		}
		if (a.Months == nil) == (a.Windows == nil) {
			return nil, missing(fn)
		}
		var s PeriodSpec
		if a.Months != nil {
			s.Months = *a.Months
		} else {
			s.Windows = *a.Windows
		}
		return records(in).periods(s)
	case "horizonte", "vista_horizonte":
		var a horizonArg
		if err := dec(&a); err != nil {
			return nil, err
		}
		s, err := a.spec(fn)
		if err != nil {
			return nil, err
		}
		if fn == "horizonte" {
			return Horizon(*in, s)
		}
		return in.viewHorizon(s)
	case "deuda_viva":
		var a struct {
			Cut *Date `json:"cut"`
		}
		if err := dec(&a); err != nil {
			return nil, err
		}
		if a.Cut == nil {
			return nil, missing(fn)
		}
		return records(in).liveDebt(*a.Cut), nil
	case "promedios":
		var a struct {
			Total *Cents  `json:"total_cents"`
			Start *Date   `json:"start"`
			End   *Date   `json:"end"`
			Month *Month  `json:"month"`
			Basis *string `json:"basis"`
		}
		if err := dec(&a); err != nil {
			return nil, err
		}
		f1 := a.Total != nil && a.Start != nil && a.End != nil && a.Month == nil && a.Basis == nil
		f2 := a.Month != nil && a.Basis != nil && a.Total == nil && a.Start == nil && a.End == nil
		switch {
		case f1:
			return Averages(*a.Total, *a.Start, *a.End)
		case f2:
			ps, _ := records(in).periods(PeriodSpec{Months: []Month{*a.Month}})
			var t Cents
			switch *a.Basis {
			case "cash":
				t = ps[0].Cash.TotalCents
			case "accrued":
				t = ps[0].Accrued.TotalCents
			case "accrued_groceries":
				t = ps[0].Accrued.ByCategory.Get(in.Params.GroceriesCategory)
			default:
				return nil, errf(CodeInvalidArgument, "basis %q", *a.Basis)
			}
			return Averages(t, a.Month.First(), a.Month.Last())
		}
		return nil, missing(fn)
	case "promedio_mensual":
		var a struct {
			Total *Cents `json:"total_cents"`
			N     *int   `json:"month_count"`
		}
		if err := dec(&a); err != nil {
			return nil, err
		}
		if a.Total == nil || a.N == nil {
			return nil, missing(fn)
		}
		v, err := MonthlyAverage(*a.Total, *a.N)
		if err != nil {
			return nil, err
		}
		return map[string]Cents{"monthly_cents": v}, nil
	case "valor_estable", "vista_valor_estable", "vista_sensibilidad":
		var a itemsArg
		if err := dec(&a); err != nil {
			return nil, err
		}
		items := in.StableItems
		if a.Items != nil {
			if err := in.checkStable(*a.Items); err != nil {
				return nil, err
			}
			items = *a.Items
		} else if fn == "vista_sensibilidad" {
			return nil, missing(fn)
		}
		switch fn {
		case "valor_estable":
			return StableValue(items, in.Persons)
		case "vista_valor_estable":
			return in.viewStable(items)
		}
		return in.viewSensitivity(items)
	case "papel", "vista_papel", "metas":
		var a struct{}
		if err := dec(&a); err != nil {
			return nil, err
		}
		switch fn {
		case "papel":
			return PaperPlan(*in)
		case "vista_papel":
			return in.viewPaper()
		}
		return SavingsGoals(*in)
	case "presupuesto_mes", "vista_periodo":
		var a monthArg
		if err := dec(&a); err != nil {
			return nil, err
		}
		if a.Month == nil {
			return nil, missing(fn)
		}
		if fn == "vista_periodo" {
			return in.viewPeriod(*a.Month)
		}
		return BudgetMonth(*in, *a.Month)
	case "aplicar_cambio_presupuesto":
		var a struct {
			RowID  *int64  `json:"row_id"`
			Month  *Month  `json:"month"`
			Amount *Cents  `json:"amount_cents"`
			Scope  *string `json:"scope"`
		}
		if err := dec(&a); err != nil {
			return nil, err
		}
		if a.RowID == nil || a.Month == nil || a.Amount == nil {
			return nil, missing(fn)
		}
		ch := BudgetChange{RowID: *a.RowID, Month: *a.Month, AmountCents: *a.Amount}
		if a.Scope != nil {
			ch.Scope = *a.Scope
			if ch.Scope == "" {
				return nil, errf(CodeInvalidArgument, "scope vacío")
			}
		}
		return ApplyBudgetChange(in.Budgets, ch)
	case "plan_contra_real":
		var a struct {
			From     *Month  `json:"from"`
			Count    *int    `json:"count"`
			Category *string `json:"category"`
		}
		if err := dec(&a); err != nil {
			return nil, err
		}
		if a.From == nil || a.Count == nil {
			return nil, missing(fn)
		}
		cat := ""
		if a.Category != nil {
			cat = *a.Category
		}
		return PlanVsActual(*in, *a.From, *a.Count, cat)
	case "flujo_quincenal":
		var a struct {
			From  *Fortnight `json:"from"`
			Count *int       `json:"count"`
		}
		if err := dec(&a); err != nil {
			return nil, err
		}
		if a.From == nil || a.Count == nil {
			return nil, missing(fn)
		}
		return FortnightFlow(*in, *a.From, *a.Count)
	case "proyeccion":
		var a struct {
			From *Month `json:"from"`
		}
		if err := dec(&a); err != nil {
			return nil, err
		}
		var m Month
		if a.From != nil {
			m = *a.From
		}
		return Projection(*in, m)
	case "vista_interes_por_monto":
		var a struct {
			ID *int64 `json:"purchase_id"`
		}
		if err := dec(&a); err != nil {
			return nil, err
		}
		if a.ID == nil {
			return nil, missing(fn)
		}
		return in.viewInterest(*a.ID)
	}
	return nil, errf(CodeInvalidArgument, "funcion desconocida %q", strings.TrimSpace(fn))
}
