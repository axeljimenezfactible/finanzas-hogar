package semilla

import (
	"fmt"
	"math/big"
	"sort"

	"finanzas-hogar/internal/dominio"
	mc "finanzas-hogar/internal/motor_cred"
)

// halfUp es ROUND_HALF_UP de num/den con enteros no negativos.
func halfUp(num, den int64) int64 {
	n := new(big.Int).Mul(big.NewInt(num), big.NewInt(2))
	n.Add(n, big.NewInt(den))
	return n.Div(n, big.NewInt(2*den)).Int64()
}

// kueski: la línea revolvente Kueski del titular con sus tasas por plazo (2 y 4 quincenas). Los términos salen de
// data.json (intereses a 2 quincenas) y de oraculo/model.json (límite, capital e intereses a 4 quincenas); la tasa se
// guarda como fracción exacta con los decimales capturados (K12: 17.27 % y 25.8685 %).
func (l *loader) kueski() error {
	d, err := readData(l.opt.Dir)
	if err != nil {
		return err
	}
	m, err := readModel(l.opt.Dir, l.opt.ModelPath)
	if err != nil {
		return err
	}
	d1 := m.Despensas[0]
	cid, _, err := l.ensure("créditos", "credit", `SELECT id FROM credit WHERE name = ? COLLATE NOCASE`, a("Kueski"),
		`INSERT INTO credit (name, kind, lender, holder_person_id, credit_limit_cents, payment_frequency, first_payment_min_days,
		   availability_reading, availability_confirmed, reading_a_anchor_date, notes)
		 VALUES ('Kueski', 'revolving_line', 'Kueski', ?, ?, 'fortnightly', 15, 'B', 0, ?, ?)`,
		l.id("person:axel"), int64(m.Parametros.Limite), d1.Fecha,
		"Línea revolvente del titular. Lectura del disponible B por defecto, pendiente de confirmar con Kueski (decisión C); la ancla de la lectura A es la fecha de la primera despensa.")
	if err != nil {
		return err
	}
	l.ids["credit:kueski"] = cid
	rates := []struct {
		n        int
		num, den int64
		dec      int
		note     string
	}{
		{2, int64(d.Financiamiento.DosQuincenas.Intereses), int64(d1.Capital), 2, "Intereses a 2 quincenas de data.json (financiamiento.dos_quincenas) sobre el capital de la despensa 1."},
		{len(d1.Pagos), int64(d1.Intereses), int64(d1.Capital), 4, "Intereses a 4 quincenas de oraculo/model.json (despensas[0]) sobre el capital de la despensa 1."},
	}
	for _, r := range rates {
		if _, _, err := l.ensure("tasas por plazo", "credit_term_rate", `SELECT id FROM credit_term_rate WHERE credit_id = ? AND term_count = ?`, a(cid, r.n),
			`INSERT INTO credit_term_rate (credit_id, term_count, rate_num, rate_den, source, rate_decimals, notes) VALUES (?, ?, ?, ?, 'manual', ?, ?)`,
			cid, r.n, r.num, r.den, r.dec, r.note); err != nil {
			return err
		}
	}
	return nil
}

type itemLine struct {
	pres       int64
	name       string
	qty, total int64
	byAccount  map[string]int64 // cuenta → importe (centavos)
	qtyByAcc   map[string]int64
	rule, bp   int64 // regla de reparto elegida y puntos base de Karen
}

// lineas arma los renglones del ticket desde data.json: uno por producto (36), con la cantidad y el importe de
// Bodega Aurrera y la cuenta de cada pieza (axel, karen, conjunto o mixto). La cuenta sale de data.json: cuentas.
func (l *loader) lineas(d *dataFile) ([]*itemLine, error) {
	byPres := map[int64]*itemLine{}
	var order []*itemLine
	get := func(name string) (*itemLine, error) {
		id, ok := l.pres[name]
		if !ok {
			return nil, fmt.Errorf("semilla: data.json: «%s» no está en productos.json (ni como presentación ni como nombre de ticket)", name)
		}
		if x, ok := byPres[id]; ok {
			return x, nil
		}
		x := &itemLine{pres: id, name: name, byAccount: map[string]int64{}, qtyByAcc: map[string]int64{}}
		byPres[id] = x
		order = append(order, x)
		return x, nil
	}
	type expect struct{ qty, total int64 }
	want := map[int64]expect{}
	var wantOrder []int64
	add := func(name string, q int64, t money) error {
		id, ok := l.pres[name]
		if !ok {
			return fmt.Errorf("semilla: data.json: «%s» no está en productos.json", name)
		}
		if _, dup := want[id]; !dup {
			wantOrder = append(wantOrder, id)
		}
		e := want[id]
		want[id] = expect{e.qty + q, e.total + int64(t)}
		_, err := get(name)
		return err
	}
	for _, it := range d.ItemsIdenticos {
		if err := add(it.Producto, it.Cantidad, it.BodegaTotal); err != nil {
			return nil, err
		}
	}
	for _, s := range d.Sustitutos {
		if err := add(s.Bodega, s.CantB, s.BodegaTotal); err != nil {
			return nil, err
		}
	}
	for _, acc := range []string{"axel", "karen", "conjunto"} {
		for _, ln := range d.Cuentas[acc].Lineas {
			x, err := get(ln.Producto)
			if err != nil {
				return nil, err
			}
			x.qty += ln.Cant
			x.total += int64(ln.Importe)
			x.byAccount[acc] += int64(ln.Importe)
			x.qtyByAcc[acc] += ln.Cant
		}
	}
	// Cotejo: las cuentas deben cuadrar con los artículos de Bodega. La única diferencia admitida es la sexta lata de
	// atún (atun_sexta_pieza_bodega): 1 pieza que data.json lista aparte de items_identicos.
	extra := 0
	for _, id := range wantOrder {
		x, e := byPres[id], want[id]
		if x.qty == e.qty && x.total == e.total {
			continue
		}
		if x.qty == e.qty+1 && x.total == e.total+int64(d.AtunSexta) && extra == 0 {
			extra++
			continue
		}
		return nil, fmt.Errorf("semilla: data.json: «%s» no cuadra: cuentas %d pzas / %d centavos contra artículos %d / %d", x.name, x.qty, x.total, e.qty, e.total)
	}
	var qty, total int64
	for _, x := range order {
		qty += x.qty
		total += x.total
	}
	if qty != d.Totales.ArticulosBodega || total != int64(d.Totales.Bodega) {
		return nil, fmt.Errorf("semilla: data.json: la compra suma %d piezas y %d centavos; totales dice %d y %d", qty, total, d.Totales.ArticulosBodega, int64(d.Totales.Bodega))
	}
	return order, nil
}

// ruleFor elige la regla de reparto del renglón: la de su cuenta o, si las piezas son de dos cuentas (mixto), una
// regla de un solo uso por cantidad (RN-03): Karen = su importe / importe del renglón.
func (l *loader) ruleFor(x *itemLine) (ruleID int64, karenBP int64, err error) {
	switch len(x.byAccount) {
	case 1:
		for acc := range x.byAccount {
			return l.id("rule:" + map[string]string{"axel": "hogar", "karen": "karen", "conjunto": "conjunto"}[acc]), map[string]int64{"axel": 0, "karen": 10000, "conjunto": 5000}[acc], nil
		}
	case 2:
		if x.byAccount["axel"] > 0 && x.byAccount["karen"] > 0 {
			bp := halfUp(x.byAccount["karen"]*10000, x.total)
			if halfUp(x.total*bp, 10000) != x.byAccount["karen"] {
				return 0, 0, fmt.Errorf("semilla: «%s»: el reparto mixto no se expresa en puntos base exactos", x.name)
			}
			id, err := l.insert("cost_share_rule", `INSERT INTO cost_share_rule (name, remainder_person_id, is_template, is_default, notes) VALUES (?, ?, 0, 0, ?)`,
				"Mixto: "+x.name, l.id("person:axel"), "Regla de un solo uso por cantidad (renglón mixto de data.json).")
			if err != nil {
				return 0, 0, err
			}
			if _, err := l.insert("cost_share_rule_part", `INSERT INTO cost_share_rule_part (rule_id, person_id, pct_bp) VALUES (?, ?, ?)`, id, l.id("person:karen"), bp); err != nil {
				return 0, 0, err
			}
			l.rep.count("reglas de un solo uso (mixto)", true)
			return id, bp, nil
		}
	}
	return 0, 0, fmt.Errorf("semilla: «%s»: combinación de cuentas no prevista en data.json", x.name)
}

// compra10Oct: el gasto de la despensa del 10-oct-2026 (61 piezas en 36 renglones) con crédito Kueski, devengo bimestral
// y reparto por renglón (axel, karen, conjunto o mixto). Devuelve el id del gasto y las partes congeladas (centavos).
func (l *loader) compra10Oct(d *dataFile, total int64) (int64, map[string]int64, error) {
	store := l.id("store:" + d.TiendaA)
	cat := l.id("category:groceries")
	if id, ok, err := l.find(`SELECT id FROM expense WHERE expense_date = ? AND store_id = ? AND category_id = ? AND amount_cents = ?`, d.Fecha, store, cat, total); err != nil {
		return 0, nil, err
	} else if ok {
		l.rep.count("gasto: despensa 10-oct-2026", false)
		return id, nil, nil
	}
	lines, err := l.lineas(d)
	if err != nil {
		return 0, nil, err
	}
	eid, err := l.insert("expense", `INSERT INTO expense (kind, description, category_id, amount_cents, expense_date, accrual_months, payment_method, paid_by_person_id, person_id, cost_share_rule_id, store_id)
		VALUES ('variable', ?, ?, ?, ?, 2, 'credit', NULL, NULL, NULL, ?)`,
		"Despensa del 10-oct-2026 en "+d.TiendaA, cat, total, d.Fecha, store)
	if err != nil {
		return 0, nil, err
	}
	l.rep.count("gasto: despensa 10-oct-2026", true)
	for i, x := range lines {
		rule, bp, err := l.ruleFor(x)
		if err != nil {
			return 0, nil, err
		}
		x.rule, x.bp = rule, bp
		var unit any
		if x.total%x.qty == 0 {
			unit = x.total / x.qty
		}
		if _, err := l.insert("expense_item", `INSERT INTO expense_item (expense_id, line_no, description, presentation_id, quantity_milli, unit_price_cents, line_total_cents, cost_share_rule_id)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, eid, i+1, x.name, x.pres, x.qty*1000, unit, x.total, rule); err != nil {
			return 0, nil, err
		}
		l.rep.count("renglones del ticket", true)
	}
	// RN-03: se agrupan los renglones por regla y se parte la SUMA de cada grupo (el 50/50 de la suma de «conjunto»,
	// no de cada renglón); Karen recibe HALF_UP(suma × pct), y el resto es del titular, a nombre de quien va el hogar.
	karen := int64(0)
	groups := map[int64]int64{}
	bpOf := map[int64]int64{}
	for _, x := range lines {
		groups[x.rule] += x.total
		bpOf[x.rule] = x.bp
	}
	ids := make([]int64, 0, len(groups))
	for r := range groups {
		ids = append(ids, r)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	for _, r := range ids {
		karen += halfUp(groups[r]*bpOf[r], 10000)
	}
	shares := map[string]int64{"karen": karen, "axel": total - karen}
	for _, p := range []string{"axel", "karen"} {
		if _, err := l.insert("expense_share", `INSERT INTO expense_share (expense_id, person_id, amount_cents) VALUES (?, ?, ?)`, eid, l.id("person:"+p), shares[p]); err != nil {
			return 0, nil, err
		}
	}
	return eid, shares, nil
}

func mcDate(s string) mc.Date {
	d, _ := mc.ParseDate(s)
	return d
}

// despensas: la despensa 1 (compra real del 10-oct-2026 con su gasto) y las despensas 2 a 4 programadas del oráculo,
// cada una con su compra a crédito, sus partes y su calendario de pagos TAL COMO LOS CALCULA motor_cred (QuotePurchase);
// después se cotejan contra oraculo/model.json: si una cifra difiere es un hallazgo y la carga se detiene.
func (l *loader) despensas() error {
	d, err := readData(l.opt.Dir)
	if err != nil {
		return err
	}
	m, err := readModel(l.opt.Dir, l.opt.ModelPath)
	if err != nil {
		return err
	}
	if m.Despensas[0].Fecha != d.Fecha {
		return fmt.Errorf("semilla: la despensa 1 del oráculo (%s) y la compra de data.json (%s) tienen fechas distintas", m.Despensas[0].Fecha, d.Fecha)
	}
	eid, shares, err := l.compra10Oct(d, int64(m.Despensas[0].Capital))
	if err != nil {
		return err
	}
	if shares != nil {
		d1 := m.Despensas[0]
		if shares["axel"] != int64(d1.Axel.Capital) || shares["karen"] != int64(d1.Karen.Capital) {
			return fmt.Errorf("semilla: HALLAZGO: el reparto de la compra (axel %d, karen %d) difiere del oráculo (axel %d, karen %d)", shares["axel"], shares["karen"], int64(d1.Axel.Capital), int64(d1.Karen.Capital))
		}
		l.rep.Checks = append(l.rep.Checks, fmt.Sprintf("reparto de la compra del 10-oct = oráculo (Axel %s, Karen %s)", dominio.FormatPesos(shares["axel"]), dominio.FormatPesos(shares["karen"])))
	}
	cid := l.id("credit:kueski")
	var limit, holder int64
	var minDays int
	if err := l.tx.QueryRowContext(l.ctx, `SELECT credit_limit_cents, holder_person_id, first_payment_min_days FROM credit WHERE id = ?`, cid).Scan(&limit, &holder, &minDays); err != nil {
		return err
	}
	cr := l.motorCredit(cid, limit, holder, minDays, m)
	for k, ds := range m.Despensas {
		var found int64
		if id, ok, err := l.find(`SELECT id FROM credit_purchase WHERE credit_id = ? AND purchase_date = ? AND kind = 'installment_plan' AND principal_cents = ?`, cid, ds.Fecha, int64(ds.Capital)); err != nil {
			return err
		} else if ok {
			l.rep.count("compras a crédito (despensas)", false)
			found = id
		}
		if found != 0 {
			continue
		}
		in := mc.PurchaseInput{Credit: cr, Date: mcDate(ds.Fecha), Principal: mc.Cents(ds.Capital), TermCount: len(ds.Pagos),
			Shares: []mc.ShareInput{
				{PersonID: l.id("person:axel"), Code: "axel", Principal: mc.Cents(ds.Axel.Capital)},
				{PersonID: l.id("person:karen"), Code: "karen", External: true, Principal: mc.Cents(ds.Karen.Capital)}}}
		if k == 0 { // la compra real: Kueski cotizó los intereses el día de la compra
			q, dec := mc.Cents(ds.Intereses), 4
			in.QuotedInterest, in.RateDecimals = &q, &dec
		}
		p, err := mc.QuotePurchase(in)
		if err != nil {
			return fmt.Errorf("semilla: motor_cred no pudo cotizar la despensa %d: %w", ds.N, err)
		}
		if err := cotejaDespensa(ds, p); err != nil {
			return err
		}
		var expense any
		var note any
		if k == 0 {
			expense = eid
		} else {
			note = "Despensa programada del oráculo (aún no se compra). PROVISIONAL: bloqueo S6-semilla-despensas-futuras."
		}
		pid, err := l.insert("credit_purchase", `INSERT INTO credit_purchase (credit_id, expense_id, kind, description, purchase_date, principal_cents, interest_cents, total_cents,
			term_count, payment_frequency, first_payment_date, interest_source, rate_num, rate_den, rate_basis, rate_decimals, interest_includes_vat, status, notes)
			VALUES (?, ?, 'installment_plan', ?, ?, ?, ?, ?, ?, 'fortnightly', ?, ?, ?, ?, ?, ?, 1, 'active', ?)`,
			cid, expense, fmt.Sprintf("Despensa %d", ds.N), ds.Fecha, int64(p.Principal), int64(p.Interest), int64(p.Total), len(p.Payments),
			p.Payments[0].Due.String(), p.Source, p.Rate.Num, p.Rate.Den, p.RateBasis, nullDec(p.Rate.Decimals), note)
		if err != nil {
			return err
		}
		l.rep.count("compras a crédito (despensas)", true)
		for _, sh := range p.Shares {
			if _, err := l.insert("credit_purchase_share", `INSERT INTO credit_purchase_share (credit_purchase_id, person_id, principal_cents, interest_cents, total_cents) VALUES (?, ?, ?, ?, ?)`,
				pid, l.id("person:"+sh.Code), int64(sh.Principal), int64(sh.Interest), int64(sh.Total)); err != nil {
				return err
			}
		}
		for _, pay := range p.Payments {
			sid, err := l.insert("scheduled_payment", `INSERT INTO scheduled_payment (credit_purchase_id, seq, due_date, amount_cents, principal_cents, interest_cents) VALUES (?, ?, ?, ?, ?, ?)`,
				pid, pay.Seq, pay.Due.String(), int64(pay.Amount), int64(pay.Principal), int64(pay.Interest))
			if err != nil {
				return err
			}
			l.rep.count("pagos programados", true)
			for _, sh := range pay.Shares {
				if _, err := l.insert("scheduled_payment_share", `INSERT INTO scheduled_payment_share (scheduled_payment_id, person_id, amount_cents) VALUES (?, ?, ?)`,
					sid, l.id("person:"+sh.Code), int64(sh.Amount)); err != nil {
					return err
				}
			}
		}
		l.rep.Checks = append(l.rep.Checks, fmt.Sprintf("despensa %d: intereses, total, %d pagos y partes = oraculo/model.json (motor_cred)", ds.N, len(p.Payments)))
	}
	return nil
}

func nullDec(d int) any {
	if d < 0 {
		return nil
	}
	return d
}

// motorCredit arma el crédito Kueski para el motor con lo que acaba de quedar en la base (tasas incluidas).
func (l *loader) motorCredit(id, limit, holder int64, minDays int, m *modelFile) mc.Credit {
	lim := mc.Cents(limit)
	md := minDays
	anchor := mcDate(m.Despensas[0].Fecha)
	c := mc.Credit{ID: id, Name: "Kueski", Kind: mc.KindRevolving, Holder: "axel", HolderPersonID: holder, Limit: &lim, Frequency: "fortnightly",
		FirstPaymentMinDays: &md, Reading: "B", ReadingAAnchor: &anchor, Promotions: []mc.Promotion{}}
	rows, err := l.tx.QueryContext(l.ctx, `SELECT term_count, rate_num, rate_den, coalesce(rate_decimals, -1), source FROM credit_term_rate WHERE credit_id = ? ORDER BY term_count`, id)
	if err != nil {
		return c
	}
	defer rows.Close()
	for rows.Next() {
		var n, num, den, dec int64
		var src string
		if rows.Scan(&n, &num, &den, &dec, &src) == nil {
			c.TermRates = append(c.TermRates, mc.TermRate{TermCount: int(n), Rate: mc.Rate{Num: num, Den: den, Decimals: int(dec)}, Source: src})
		}
	}
	return c
}

// cotejaDespensa compara una compra cotizada por el motor con la del oráculo. Un renglón distinto es un HALLAZGO
// (protocolo §6): se informan los dos valores y no se ajusta ninguno.
func cotejaDespensa(ds modelDespensa, p mc.Purchase) error {
	bad := func(what string, got, want any) error {
		return fmt.Errorf("semilla: HALLAZGO en la despensa %d: %s = %v en motor_cred y %v en oraculo/model.json", ds.N, what, got, want)
	}
	if int64(p.Interest) != int64(ds.Intereses) {
		return bad("intereses", int64(p.Interest), int64(ds.Intereses))
	}
	if int64(p.Total) != int64(ds.Total) {
		return bad("total", int64(p.Total), int64(ds.Total))
	}
	if len(p.Payments) != len(ds.Pagos) {
		return bad("número de pagos", len(p.Payments), len(ds.Pagos))
	}
	for i, pay := range p.Payments {
		w := ds.Pagos[i]
		if pay.Due.String() != w.Fecha || int64(pay.Amount) != int64(w.Monto) {
			return bad(fmt.Sprintf("pago %d", i+1), fmt.Sprintf("%s %d", pay.Due, int64(pay.Amount)), fmt.Sprintf("%s %d", w.Fecha, int64(w.Monto)))
		}
		for _, sh := range pay.Shares {
			want := map[string]money{"axel": w.Axel, "karen": w.Karen}[sh.Code]
			if int64(sh.Amount) != int64(want) {
				return bad(fmt.Sprintf("pago %d de %s", i+1, sh.Code), int64(sh.Amount), int64(want))
			}
		}
	}
	for _, sh := range p.Shares {
		want := map[string]modelPerson{"axel": ds.Axel, "karen": ds.Karen}[sh.Code]
		if int64(sh.Principal) != int64(want.Capital) || int64(sh.Interest) != int64(want.Intereses) || int64(sh.Total) != int64(want.Total) {
			return bad("parte de "+sh.Code, fmt.Sprintf("%d/%d/%d", sh.Principal, sh.Interest, sh.Total), fmt.Sprintf("%d/%d/%d", want.Capital, want.Intereses, want.Total))
		}
	}
	return nil
}
