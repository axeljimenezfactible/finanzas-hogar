package dominio_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

func p64(v int64) *int64    { return &v }
func pstr(v string) *string { return &v }

// TestCreditosRepos recorre los repositorios de créditos (api-creditos) contra la base migrada.
func TestCreditosRepos(t *testing.T) {
	ctx := context.Background()
	db, fx := testutil.NewDB(t)
	c := dominio.CreditRow{Name: "Kueski Pay", Kind: "revolving_line", HolderPersonID: fx.AxelPersonID, CreditLimitCents: p64(1000000),
		PaymentFrequency: "fortnightly", FirstPaymentMinDays: p64(15), AvailabilityReading: "B", IsActive: true}
	id, err := dominio.InsertCredit(ctx, db, c)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dominio.InsertCredit(ctx, db, c); !dominio.IsUnique(err) {
		t.Fatalf("nombre repetido: %v", err)
	}
	c, err = dominio.GetCredit(ctx, db, id)
	if err != nil || c.Name != "Kueski Pay" || *c.CreditLimitCents != 1000000 || c.VatIncluded != nil {
		t.Fatalf("%+v %v", c, err)
	}
	c.Notes = pstr("línea")
	if err := dominio.UpdateCredit(ctx, db, c); err != nil {
		t.Fatal(err)
	}
	if cs, err := dominio.QueryCredits(ctx, db, "", "ORDER BY c.id"); err != nil || len(cs) != 1 || *cs[0].Notes != "línea" {
		t.Fatalf("%+v %v", cs, err)
	}
	if _, err := dominio.GetCredit(ctx, db, 999); !errors.Is(err, sql.ErrNoRows) {
		t.Fatal(err)
	}
	c.ID = 999
	if err := dominio.UpdateCredit(ctx, db, c); !errors.Is(err, sql.ErrNoRows) {
		t.Fatal(err)
	}
	tr := dominio.TermRateRow{CreditID: id, TermCount: 4, RateNum: 80050, RateDen: 309450, Source: "quote", QuoteDate: pstr("2026-10-10"), RateDecimals: p64(4)}
	trID, err := dominio.InsertTermRate(ctx, db, tr)
	if err != nil {
		t.Fatal(err)
	}
	tr, err = dominio.GetTermRate(ctx, db, trID)
	if err != nil || tr.RateNum != 80050 {
		t.Fatalf("%+v %v", tr, err)
	}
	tr.Notes = pstr("cotización")
	if err := dominio.UpdateTermRate(ctx, db, tr); err != nil {
		t.Fatal(err)
	}
	if _, err := dominio.GetTermRate(ctx, db, 999); !errors.Is(err, sql.ErrNoRows) {
		t.Fatal(err)
	}
	pr := dominio.PromotionRow{CreditID: &id, MerchantName: "Chedraui", TermCount: 4, RateDen: 1, ValidFrom: "2026-10-30", ValidTo: "2026-10-30", IsActive: true}
	prID, err := dominio.InsertPromotion(ctx, db, pr)
	if err != nil {
		t.Fatal(err)
	}
	pr, err = dominio.GetPromotion(ctx, db, prID)
	if err != nil || pr.MerchantName != "Chedraui" {
		t.Fatalf("%+v %v", pr, err)
	}
	pr.IsActive = false
	if err := dominio.UpdatePromotion(ctx, db, pr); err != nil {
		t.Fatal(err)
	}
	if _, err := dominio.GetPromotion(ctx, db, 999); !errors.Is(err, sql.ErrNoRows) {
		t.Fatal(err)
	}
	pu := dominio.PurchaseRow{CreditID: id, Kind: "installment_plan", PurchaseDate: "2026-10-10", PrincipalCents: 309450, InterestCents: 80050,
		TotalCents: 389500, TermCount: 4, PaymentFrequency: "fortnightly", FirstPaymentDate: "2026-10-30", InterestSource: "quoted",
		InterestIncludesVat: true, Status: "active"}
	puID, err := dominio.InsertPurchase(ctx, db, pu)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dominio.InsertPurchaseShare(ctx, db, dominio.PurchaseShareRow{CreditPurchaseID: puID, PersonID: fx.AxelPersonID,
		PrincipalCents: 309450, InterestCents: 80050, TotalCents: 389500}); err != nil {
		t.Fatal(err)
	}
	spID, err := dominio.InsertScheduled(ctx, db, dominio.ScheduledRow{CreditPurchaseID: puID, Seq: 1, DueDate: "2026-10-30",
		AmountCents: 97375, PrincipalCents: 77363, InterestCents: 20012})
	if err != nil {
		t.Fatal(err)
	}
	if err := dominio.InsertScheduledShare(ctx, db, spID, fx.AxelPersonID, 97375); err != nil {
		t.Fatal(err)
	}
	payID, err := dominio.InsertPayment(ctx, db, dominio.PaymentRow{CreditID: id, ScheduledPaymentID: &spID, PayerPersonID: fx.AxelPersonID,
		PaidDate: "2026-10-12", AmountCents: 1000})
	if err != nil {
		t.Fatal(err)
	}
	pu, err = dominio.GetPurchase(ctx, db, puID)
	if err != nil || pu.TotalCents != 389500 {
		t.Fatalf("%+v %v", pu, err)
	}
	pu.Description = pstr("Despensa 1")
	if err := dominio.UpdatePurchaseText(ctx, db, pu); err != nil {
		t.Fatal(err)
	}
	if _, err := dominio.GetPurchase(ctx, db, 999); !errors.Is(err, sql.ErrNoRows) {
		t.Fatal(err)
	}
	if ps, err := dominio.QueryPurchases(ctx, db, "", "ORDER BY p.id"); err != nil || len(ps) != 1 || *ps[0].Description != "Despensa 1" {
		t.Fatalf("%+v %v", ps, err)
	}
	if sh, err := dominio.PurchaseShares(ctx, db, puID); err != nil || len(sh[puID]) != 1 {
		t.Fatalf("%+v %v", sh, err)
	}
	if sps, err := dominio.QueryScheduled(ctx, db, "", "ORDER BY sp.id"); err != nil || len(sps) != 1 || sps[0].CreditID != id {
		t.Fatalf("%+v %v", sps, err)
	}
	if ss, err := dominio.ScheduledShares(ctx, db, 0); err != nil || ss[spID][0].PaidCents != 1000 {
		t.Fatalf("%+v %v", ss, err)
	}
	if ms, err := dominio.QueryPayments(ctx, db, "m.id = ?", payID); err != nil || len(ms) != 1 || ms[0].PurchaseID != puID || ms[0].Seq != 1 {
		t.Fatalf("%+v %v", ms, err)
	}
	if _, err := db.Exec(`INSERT INTO credit (name, kind, holder_person_id, credit_limit_cents, payment_frequency, statement_day, payment_due_day)
 VALUES ('Tarjeta', 'bank_card', ?, 100, 'monthly', 5, 25)`, fx.AxelPersonID); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`INSERT INTO credit_statement (credit_id, statement_date, due_date, balance_cents, source) VALUES (2, '2026-10-05', '2026-10-25', 50, 'bank')`); err != nil {
		t.Fatal(err)
	}
	if st, err := dominio.QueryStatements(ctx, db); err != nil || len(st) != 1 || st[0].BalanceCents != 50 {
		t.Fatalf("%+v %v", st, err)
	}
	if ps, err := dominio.CreditPersons(ctx, db); err != nil || ps[fx.KarenPersonID].Participation != "external" {
		t.Fatalf("%+v %v", ps, err)
	}
	if n, err := dominio.CountWhere(ctx, db, "credit_purchase", "promotion_id IS NULL"); err != nil || n != 1 {
		t.Fatalf("%d %v", n, err)
	}
	if ok, err := dominio.Exists(ctx, db, "credit", id); err != nil || !ok {
		t.Fatal(ok, err)
	}
	if err := dominio.DeleteRow(ctx, db, "credit_payment", payID); err != nil {
		t.Fatal(err)
	}
	if err := dominio.DeleteRow(ctx, db, "credit_payment", payID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatal(err)
	}
	if _, err := dominio.QueryCredits(ctx, db, "no_such_col = 1", ""); err == nil {
		t.Fatal("una consulta mala debe fallar")
	}
}
