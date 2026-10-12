package motor_cred

// Vistas con las claves literales de model.json y data.json (§9.3). Dinero en centavos sin sufijo.

type obj = map[string]any

func ratio(r Rate) obj { return obj{"num": r.Num, "den": r.Den} }

func quotedCredit(in Input, id int64) (CreditState, []quotedPurchase, error) {
	st, err := in.CreditByID(id)
	if err != nil {
		return CreditState{}, nil, err
	}
	qs, err := st.quoted()
	return st, qs, err
}

func nth(qs []quotedPurchase, n int) (quotedPurchase, error) {
	if n < 1 || n > len(qs) {
		return quotedPurchase{}, invalid("despensa %d inexistente", n)
	}
	return qs[n-1], nil
}

func personView(s PersonShare) obj {
	o := obj{"capital": s.Principal, "intereses": s.Interest, "total": s.Total, "pago_por_quincena": s.FirstPayment}
	if s.Composition != nil {
		o["conjunto"], o["propio"] = s.Composition.Joint, s.Composition.Own
		o["papel_dentro"], o["libre"] = s.Composition.PaperIn, s.Composition.Free
	}
	return o
}

// ViewParametros reproduce parametros.
func ViewParametros(in Input, creditID int64) (obj, error) {
	st, qs, err := quotedCredit(in, creditID)
	if err != nil {
		return nil, err
	}
	d1, err := nth(qs, 1)
	if err != nil {
		return nil, err
	}
	o := obj{"tasa": ratio(d1.Purchase.Rate), "tasa_pct": ratio(d1.Purchase.Rate), "limite_credito": st.Credit.Limit}
	if len(in.Cash) > 0 {
		o["papel"] = in.Cash[0].Amount
	}
	for _, s := range d1.Purchase.Shares {
		if s.External {
			k := obj{"capital": s.Principal, "intereses": s.Interest, "total": s.Total, "pago": s.FirstPayment}
			if s.Composition != nil {
				k["conjunto"] = s.Composition.Joint
			}
			o["karen_fija"] = k
			break
		}
	}
	return o, nil
}

// ViewDespensa reproduce despensas[n−1].
func ViewDespensa(in Input, creditID int64, n int) (obj, error) {
	_, qs, err := quotedCredit(in, creditID)
	if err != nil {
		return nil, err
	}
	q, err := nth(qs, n)
	if err != nil {
		return nil, err
	}
	p := q.Purchase
	pagos := make([]obj, len(p.Payments))
	for i, sp := range p.Payments {
		po := obj{"fecha": sp.Due.String(), "monto": sp.Amount, "capital_repuesto": sp.Principal}
		for _, sh := range sp.Shares {
			po[sh.Code] = sh.Amount
		}
		pagos[i] = po
	}
	o := obj{"n": n, "fecha": q.Input.Date.String(), "capital": p.Principal, "intereses": p.Interest, "total": p.Total, "pagos": pagos}
	for _, s := range p.Shares {
		o[s.Code] = personView(s)
	}
	return o, nil
}

// ViewDisposicion reproduce disposicion_por_despensa[n−1].
func ViewDisposicion(in Input, creditID int64, n int) (obj, error) {
	_, qs, err := quotedCredit(in, creditID)
	if err != nil {
		return nil, err
	}
	q, err := nth(qs, n)
	if err != nil {
		return nil, err
	}
	o := obj{"despensa": n, "fecha": q.Input.Date.String()}
	for _, s := range q.Purchase.Shares {
		o[s.Code] = personView(s)
	}
	return o, nil
}

var tipoEvento = map[string]string{"start": "inicio", "purchase": "compra", "cash": "papel", "payment": "pago"}

// ViewCredito reproduce credito.
func ViewCredito(in Input, creditID int64, today, from, to Date) (obj, error) {
	st, err := in.CreditByID(creditID)
	if err != nil {
		return nil, err
	}
	tl, err := AvailabilityTimeline(TimelineInput{State: st, Cash: in.Cash, Today: today, From: from, To: to})
	if err != nil {
		return nil, err
	}
	evs := make([]obj, len(tl.Events))
	for i, e := range tl.Events {
		evs[i] = obj{"fecha": e.Date.String(), "tipo": tipoEvento[e.Kind], "monto": e.Amount,
			"capital_repuesto": e.PrincipalRestored, "interes_no_repuesto": e.InterestNotRestored,
			"disponible_A": e.AvailableA, "disponible_B": e.AvailableB}
	}
	return obj{"eventos": evs, "minimo_A": tl.MinA, "minimo_B": tl.MinB,
		"comparacion_series": obj{"esperado_A": tl.SeriesA, "calculado_A": tl.SeriesA,
			"esperado_B": tl.SeriesB, "calculado_B": tl.SeriesB, "coinciden": true}}, nil
}

// ViewSensibilidad reproduce sensibilidad.
func ViewSensibilidad(si SensitivityInput) ([]obj, error) {
	rows, err := Sensitivity(si)
	if err != nil {
		return nil, err
	}
	out := make([]obj, len(rows))
	for i, r := range rows {
		out[i] = obj{"N": r.Amount, "intereses": r.Interest, "total": r.Total, "pago_quincenal": r.ReferencePayment,
			"axel_total": r.HolderTotal, "axel_capital": r.HolderPrincipal, "axel_menor_que_actual": r.HolderBelowCurrent,
			"costo_mensual_total": r.MonthlyTotal, "axel_mensual": r.HolderMonthly, "karen_mensual": r.ExternalMonthly}
	}
	return out, nil
}

// ViewInteresPorMonto reproduce interes_por_monto.
func ViewInteresPorMonto(amounts []Cents, r Rate, n int) ([]obj, error) {
	if !r.valid() || n < 1 {
		return nil, invalid("tasa o plazo")
	}
	out := make([]obj, len(amounts))
	for i, a := range amounts {
		q := QuoteByAmount(a, r, n)
		out[i] = obj{"monto": q.Amount, "intereses": q.Interest, "total": q.Total, "pago_quincenal": q.ReferencePayment}
	}
	return out, nil
}

// ViewDosQuincenas reproduce data.json:financiamiento.dos_quincenas (R-05, R-06).
func ViewDosQuincenas(in Input, creditID int64, pin PurchaseInput) (obj, error) {
	st, err := in.CreditByID(creditID)
	if err != nil {
		return nil, err
	}
	pin.Credit = st.Credit
	p, err := QuotePurchase(pin)
	if err != nil {
		return nil, err
	}
	o := obj{"total": p.Total, "intereses": p.Interest,
		"tasa_pct": obj{"num": p.Rate.Num, "den": p.Rate.Den, "decimales": p.Rate.Decimals}}
	for _, s := range p.Shares {
		o[s.Code] = obj{"intereses": s.Interest, "total": s.Total, "por_quincena": s.FirstPayment}
	}
	return o, nil
}
