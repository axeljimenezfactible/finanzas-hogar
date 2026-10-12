package gastos

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"slices"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	mc "finanzas-hogar/internal/motor_cred"
	"finanzas-hogar/internal/recetas"
)

// Confirmar un ticket (N-074, RN-23, K3). Una sola transacción crea:
//   - el gasto con sus renglones y la cuenta de cada renglón (mismo alta que createExpense, reparto RN-03);
//   - un precio por presentación con fuente ticket_ocr, la fecha del ticket y su tienda (decisión E);
//   - la entrada de inventario de cada producto con existencia por unidad (RN-15), con su inventory_movement
//     purchase ligado al renglón;
//   - con forma de pago crédito, la compra a crédito; sus partes y su serie de pagos salen de motor_cred.QuotePurchase.
// Al final el ticket queda confirmed, ligado al gasto y sin borrador.

// prefixed pone la ruta JSON del objeto anidado en los campos de un 422 (C1 §3: «ruta JSON del campo con error»).
func prefixed(err error, prefix string) error {
	e, ok := dominio.AsError(err)
	if !ok || len(e.Fields) == 0 {
		return err
	}
	out := *e
	out.Fields = make([]dominio.FieldIssue, len(e.Fields))
	for i, f := range e.Fields {
		f.Field = prefix + f.Field
		out.Fields[i] = f
	}
	return &out
}

// ticketUnitPrice es el precio por presentación del renglón: el capturado o, si falta, importe del renglón entre
// presentaciones con HALF_UP (PROVISIONAL, bloqueo §18).
func ticketUnitPrice(it dominio.ExpenseItemIn) int64 {
	if it.UnitPriceCents != nil {
		return *it.UnitPriceCents
	}
	return dominio.MulDivRound(it.LineTotalCents, 1000, it.QuantityMilli)
}

// stockDelta es lo que entra a la existencia en milésimas de la unidad base: el motor S4 multiplica presentaciones por
// contenido (recetas.PurchaseMovement) y, como la cantidad del renglón viene ×1000, se quita ese factor con HALF_UP
// (1.250 kg a granel = 1250 × 1 000 000 / 1000).
func stockDelta(product string, date string, quantityMilli, sizeMilli int64) (int64, error) {
	d, err := recetas.ParseDate(date)
	if err != nil {
		return 0, err
	}
	return dominio.RoundHalfUp(recetas.PurchaseMovement(product, d, quantityMilli, sizeMilli).DeltaMilli, 1000), nil
}

// ticketCredit valida la compra a crédito del ticket y la cotiza con el motor S3. Las reglas son las de
// createCreditPurchase (api-creditos) más las del ticket: el capital es el importe del gasto y la fecha, la del ticket.
func ticketCredit(ctx context.Context, q dominio.Querier, today dominio.Date, in dominio.ExpenseIn, cp api.CreditPurchaseCreate) (mc.Purchase, dominio.CreditRow, error) {
	const f = "credit_purchase."
	s, err := loadCredInput(ctx, q, today)
	if err != nil {
		return mc.Purchase{}, dominio.CreditRow{}, err
	}
	var c dominio.CreditRow
	found := false
	for _, x := range s.credits {
		if x.ID == cp.CreditID {
			c, found = x, true
		}
	}
	var v api.Validator
	v.Check(found, f+"credit_id", "not_found", "No encontramos ese crédito.")
	v.Check(slices.Contains([]string{"installment_plan", "single_charge", "loan_disbursement"}, cp.Kind), f+"kind", "invalid_format", "Elige un tipo de compra válido.")
	v.Check(cp.PrincipalCents == in.AmountCents, f+"principal_cents", "invalid_sum", "El capital de la compra debe ser el total del ticket.")
	v.Check(cp.TermCount >= 1 && cp.TermCount <= 360, f+"term_count", "out_of_range", "El número de pagos va de 1 a 360.")
	v.Check(cp.Kind != "single_charge" || cp.TermCount == 1, f+"term_count", "out_of_range", "Un cargo normal es de un solo pago.")
	v.Check(cp.InterestCents == nil || *cp.InterestCents >= 0, f+"interest_cents", "out_of_range", "Los intereses no pueden ser negativos.")
	v.Check(cp.PaymentFrequency == nil || !found || *cp.PaymentFrequency == c.PaymentFrequency, f+"payment_frequency", "not_allowed", "La frecuencia de los pagos es la del crédito.")
	v.Check(cp.Expense == nil && cp.ExpenseID == nil, f+"expense", "not_allowed", "El gasto de la compra es el del ticket; no mandes expense ni expense_id.")
	date, err := v.Date(cp.PurchaseDate, f+"purchase_date")
	if err != nil {
		return mc.Purchase{}, c, err
	}
	v.Check(date.String() == in.ExpenseDate, f+"purchase_date", "not_allowed", "La fecha de la compra es la del ticket.")
	var first *mc.Date
	if cp.FirstPaymentDate != nil {
		d, err := v.Date(*cp.FirstPaymentDate, f+"first_payment_date")
		if err != nil {
			return mc.Purchase{}, c, err
		}
		first = api.Ptr(cdateOf(d))
	}
	var shares []mc.ShareInput
	var sum int64
	for _, sh := range cp.Shares {
		p, ok := s.persons[sh.PersonID]
		v.Check(ok, f+"shares", "not_found", "Una de las personas de las partes no existe.")
		v.Check(sh.PrincipalCents >= 0, f+"shares", "out_of_range", "El capital de cada parte no puede ser negativo.")
		sum += sh.PrincipalCents
		shares = append(shares, mc.ShareInput{PersonID: sh.PersonID, Code: p.Code, External: p.Participation == "external", Principal: mc.Cents(sh.PrincipalCents)})
	}
	v.Check(len(cp.Shares) == 0 || sum == cp.PrincipalCents, f+"shares", "invalid_sum", "Las partes deben sumar el monto de la compra.")
	if err := v.Err(); err != nil {
		return mc.Purchase{}, c, err
	}
	if shares == nil { // sin partes: todo al titular del crédito (la regla de api-creditos)
		h := s.persons[c.HolderPersonID]
		shares = []mc.ShareInput{{PersonID: c.HolderPersonID, Code: h.Code, External: h.Participation == "external", Principal: mc.Cents(cp.PrincipalCents)}}
	}
	pin := mc.PurchaseInput{Credit: s.credit(c), Date: cdateOf(date), Principal: mc.Cents(cp.PrincipalCents), TermCount: int(cp.TermCount),
		QuotedInterest: centsPtr(cp.InterestCents), RateDecimals: intPtr(cp.RateDecimals), PromotionID: cp.PromotionID, FirstPaymentDate: first, Shares: shares}
	for _, p := range s.promos {
		if cp.PromotionID != nil && p.ID == *cp.PromotionID {
			pin.Merchant, pin.StoreID = p.MerchantName, p.StoreID
		}
	}
	p, err := mc.QuotePurchase(pin)
	if err != nil {
		return p, c, credMotorErr(err)
	}
	if cp.InterestSource != nil && *cp.InterestSource != p.Source {
		return p, c, dominio.Invalid(dominio.FieldIssue{Field: f + "interest_source", Code: "not_allowed",
			Message: "El origen de los intereses no coincide con el cálculo (" + p.Source + ")."})
	}
	return p, c, nil
}

// credMotorErr traduce un error del motor S3 a la respuesta de C1, como api-creditos (422 con el campo, 503 por tiempo).
func credMotorErr(err error) error {
	issue := func(field, code, msg string) error {
		return dominio.Invalid(dominio.FieldIssue{Field: "credit_purchase." + field, Code: code, Message: msg})
	}
	switch {
	case errors.Is(err, mc.ErrTimeout), errors.Is(err, context.DeadlineExceeded):
		return dominio.Unavailable("El cálculo del crédito tardó más del tiempo máximo. Intenta de nuevo.")
	case errors.Is(err, mc.ErrNoRate):
		return issue("term_count", "not_found", "Este crédito no tiene tasa capturada para ese plazo. Captúrala en Ajustes o manda los intereses cotizados.")
	case errors.Is(err, mc.ErrMissingCapture):
		return issue("credit_id", "required", "Faltan datos del estado de cuenta de este crédito (tasa, base, IVA o pago mínimo). Captúralos en Ajustes.")
	case errors.Is(err, mc.ErrPromotionNotApplicable):
		return issue("promotion_id", "not_allowed", "La promoción no aplica a esta compra (comercio, compra mínima, plazo o vigencia).")
	case errors.Is(err, mc.ErrUnsupported):
		return issue("credit_id", "not_allowed", "Este crédito todavía no tiene regla de intereses para este caso.")
	case errors.Is(err, mc.ErrInvalid):
		return issue("credit_id", "out_of_range", "Los datos no cuadran para el cálculo del crédito. Revisa montos, partes, plazo y fechas.")
	}
	return err
}

// insertTicketPurchase guarda la compra cotizada con sus partes y su serie (las mismas tablas y columnas que
// createCreditPurchase), ligada al gasto del ticket, con audit_log.
func insertTicketPurchase(r *http.Request, tx *sql.Tx, cp api.CreditPurchaseCreate, p mc.Purchase, c dominio.CreditRow, expenseID int64) (int64, error) {
	ctx := r.Context()
	includesVat := true
	if cp.InterestIncludesVat != nil {
		includesVat = *cp.InterestIncludesVat
	}
	row := dominio.PurchaseRow{CreditID: c.ID, ExpenseID: &expenseID, Kind: cp.Kind, Description: cp.Description, PurchaseDate: cp.PurchaseDate,
		PrincipalCents: int64(p.Principal), InterestCents: int64(p.Interest), TotalCents: int64(p.Total), TermCount: cp.TermCount,
		PaymentFrequency: c.PaymentFrequency, FirstPaymentDate: p.Payments[0].Due.String(), InterestSource: p.Source,
		InterestIncludesVat: includesVat, PromotionID: p.PromotionID, Status: "active", Notes: cp.Notes, RateDecimals: cp.RateDecimals}
	if p.RateBasis != "" && p.Rate.Den > 0 {
		row.RateNum, row.RateDen, row.RateBasis = api.Ptr(p.Rate.Num), api.Ptr(p.Rate.Den), api.Ptr(p.RateBasis)
	} else {
		row.RateDecimals = nil
	}
	id, err := dominio.InsertPurchase(ctx, tx, row)
	if err != nil {
		return 0, err
	}
	for _, sh := range p.Shares {
		if _, err := dominio.InsertPurchaseShare(ctx, tx, dominio.PurchaseShareRow{CreditPurchaseID: id, PersonID: sh.PersonID,
			PrincipalCents: int64(sh.Principal), InterestCents: int64(sh.Interest), TotalCents: int64(sh.Total)}); err != nil {
			return 0, err
		}
	}
	for _, sp := range p.Payments {
		spID, err := dominio.InsertScheduled(ctx, tx, dominio.ScheduledRow{CreditPurchaseID: id, Seq: int64(sp.Seq), DueDate: sp.Due.String(),
			AmountCents: int64(sp.Amount), PrincipalCents: int64(sp.Principal), InterestCents: int64(sp.Interest)})
		if err != nil {
			return 0, err
		}
		for _, sh := range sp.Shares {
			if err := dominio.InsertScheduledShare(ctx, tx, spID, sh.PersonID, int64(sh.Amount)); err != nil {
				return 0, err
			}
		}
	}
	after, err := dominio.RowJSON(ctx, tx, "credit_purchase", id)
	if err != nil {
		return 0, err
	}
	return id, dominio.Audit(ctx, tx, api.Actor(r), "insert", "credit_purchase", id, nil, after, api.RequestID(r))
}

func confirmTicketScan(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "attachment_id")
		if err != nil {
			return err
		}
		var b api.TicketConfirmRequest
		if err := reg.DecodeJSON(r, &b); err != nil {
			return err
		}
		var v api.Validator
		d, err := v.Date(b.Expense.ExpenseDate, "expense.expense_date")
		if err != nil {
			return err
		}
		ctx, db, clock := r.Context(), reg.Deps().DB, reg.Deps().Clock
		t, err := dominio.GetTicketScan(ctx, db, id)
		if err != nil {
			return err
		}
		if t.OCRStatus == "confirmed" || t.ExpenseID.Valid {
			return dominio.Conflict("conflict", "Este ticket ya está confirmado. Abre su gasto para hacer cambios.")
		}
		in, err := newExpenseIn(reg, b.Expense, d)
		if err != nil {
			return err
		}
		credit := in.PaymentMethod == "credit"
		if err := prefixed(validateExpense(r, db, in, credit), "expense."); err != nil {
			return err
		}
		var presIDs []int64
		for _, it := range in.Items {
			if it.PresentationID != nil && !slices.Contains(presIDs, *it.PresentationID) {
				presIDs = append(presIDs, *it.PresentationID)
			}
		}
		pres, err := dominio.GastosTicketPresentations(ctx, db, presIDs)
		if err != nil {
			return err
		}
		v.Check(len(pres) == len(presIDs), "expense.items", "not_found", "Una de las presentaciones del ticket no existe en el catálogo. Elige otra o deja el renglón sin ligar.")
		v.Check(len(presIDs) == 0 || in.StoreID != nil, "expense.store_id", "required", "Elige la tienda del ticket: los precios se guardan por tienda.")
		v.Check(!credit || b.CreditPurchase != nil, "credit_purchase", "required", "Con forma de pago crédito, captura los datos de la compra a crédito.")
		v.Check(credit || b.CreditPurchase == nil, "credit_purchase", "not_allowed", "Los datos de la compra a crédito solo van con forma de pago crédito.")
		if err := v.Err(); err != nil {
			return err
		}
		var quote mc.Purchase
		var cr dominio.CreditRow
		if credit {
			if quote, cr, err = ticketCredit(ctx, db, dominio.Today(clock), in, *b.CreditPurchase); err != nil {
				return err
			}
		}
		now, actor, reqID := dominio.FormatDBTime(clock.Now()), api.Actor(r), api.RequestID(r)
		out := api.TicketConfirmResult{PriceIDs: []int64{}, InventoryItemIDs: []int64{}}
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			eid, err := dominio.SaveExpense(ctx, tx, 0, in, true)
			if err != nil {
				return err
			}
			if err := audit(r, tx, "insert", eid, nil); err != nil {
				return err
			}
			out.ExpenseID = eid
			lineIDs, err := dominio.GastosExpenseItemIDs(ctx, tx, eid)
			if err != nil {
				return err
			}
			priced := map[int64]bool{}
			for _, it := range in.Items {
				if it.PresentationID == nil {
					continue
				}
				pr, eiID := pres[*it.PresentationID], lineIDs[it.LineNo]
				if !priced[pr.ID] { // un precio por presentación y ticket: el del primer renglón (PROVISIONAL, bloqueo §18)
					pid, err := dominio.GastosSaveTicketPrice(ctx, tx, dominio.TicketPrice{PresentationID: pr.ID, StoreID: *in.StoreID,
						PriceCents: ticketUnitPrice(it), ExpenseItemID: eiID, AttachmentID: id, ObservedDate: in.ExpenseDate}, now, actor, reqID)
					if err != nil {
						return err
					}
					priced[pr.ID] = true
					out.PriceIDs = append(out.PriceIDs, pid)
				}
				invID, ok, err := dominio.GastosStockItem(ctx, tx, pr.ProductID)
				if err != nil || !ok { // sin existencia por unidad: no se suma (PREGUNTA, bloqueo §18)
					if err != nil {
						return err
					}
					continue
				}
				delta, err := stockDelta(pr.ProductName, in.ExpenseDate, it.QuantityMilli, pr.SizeMilli)
				if err != nil {
					return err
				}
				if err := dominio.GastosTicketStockIn(ctx, tx, invID, in.ExpenseDate, delta, eiID, now, actor, reqID); err != nil {
					return err
				}
				if !slices.Contains(out.InventoryItemIDs, invID) {
					out.InventoryItemIDs = append(out.InventoryItemIDs, invID)
				}
			}
			if credit {
				cpID, err := insertTicketPurchase(r, tx, *b.CreditPurchase, quote, cr, eid)
				if err != nil {
					return err
				}
				out.CreditPurchaseID = &cpID
			}
			before, err := dominio.RowJSON(ctx, tx, "attachment", id)
			if err != nil {
				return err
			}
			ok, err := dominio.GastosConfirmTicket(ctx, tx, id, eid, now)
			if err != nil {
				return err
			}
			if !ok {
				return dominio.Conflict("conflict", "Este ticket ya está confirmado. Abre su gasto para hacer cambios.")
			}
			after, err := dominio.RowJSON(ctx, tx, "attachment", id)
			if err != nil {
				return err
			}
			return dominio.Audit(ctx, tx, actor, "update", "attachment", id, before, after, reqID)
		})
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusCreated, out)
	}
}
