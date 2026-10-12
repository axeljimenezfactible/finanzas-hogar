package semilla

// Gastos PROYECTADOS de la semilla REAL (WF-3b, agente semilla-c3, 12-oct-2026).
//
// DECISIÓN PROVISIONAL DEL ORQUESTADOR hasta la puerta 4 (bloqueo WF3a-bis-semilla-dev, punto 4 = opción a): la semilla
// real agrega los gastos 2 a 6 de oraculo/fixtures/s2/_datos/semilla-oraculo.json, que son los que el modelo canon
// (oraculo/model.json) asume:
//
//	gastos 2 a 4  las despensas 2 a 4 (10-dic-2026, 10-feb y 10-abr-2027; 4,000.00 a crédito, devengo a 2 meses),
//	              ligadas a su compra a crédito (credit_purchase.expense_id) para que su devengo se reparta en 2 meses;
//	gastos 5 y 6  los 2 papeles de contado de 199.00 de Karen (10-ene y 10-mar-2027; categoría «Papel»).
//
// Cómo se marca «proyectado»: el esquema NO tiene un campo para eso (expense.kind solo admite fixed, variable y
// upcoming, y «upcoming» significa «gasto venidero que ya ocurrió o quedó comprometido»; cambiarlo movería cifras de
// los motores). Por eso el gasto va como 'variable' —igual que en la fixture, que es lo que el oráculo cotejó— con la
// fecha futura tal cual, la descripción termina en «(proyectado)» y las notas lo dicen y citan el bloqueo. Ver
// WF3a-bis-semilla-dev.md, punto 4.
//
// Idempotente: cada gasto se busca por fecha, categoría, importe y medio de pago; la liga con la compra a crédito solo
// se escribe si falta (con audit_log de la actualización).

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"finanzas-hogar/internal/dominio"
)

// proyNota es la marca que llevan las notas de cada gasto proyectado.
const proyNota = "PROYECTADO (gasto futuro que asume el modelo canon, oraculo/model.json; aún no ocurre). PROVISIONAL hasta la puerta 4: bloqueo WF3a-bis-semilla-dev, punto 4. " +
	"El esquema no distingue «proyectado»: se identifica por su fecha futura y esta nota."

type oracleSeed struct {
	Categories []struct {
		Code     string `json:"code"`
		Name     string `json:"name"`
		Kind     string `json:"kind"`
		IsActive bool   `json:"is_active"`
	} `json:"categories"`
	Expenses []struct {
		ID            int64  `json:"id"`
		Kind          string `json:"kind"`
		Category      string `json:"category"`
		AmountCents   int64  `json:"amount_cents"`
		ExpenseDate   string `json:"expense_date"`
		AccrualMonths int64  `json:"accrual_months"`
		PaymentMethod string `json:"payment_method"`
		PaidBy        string `json:"paid_by"`
		Shares        []struct {
			Person      string `json:"person"`
			AmountCents int64  `json:"amount_cents"`
		} `json:"shares"`
		CreditPurchaseID *int64 `json:"credit_purchase_id"`
	} `json:"expenses"`
}

// oracleSeedPath resuelve semilla-oraculo.json: Options.OracleSeedPath, Dir/semilla-oraculo.json o
// Dir/../../oraculo/fixtures/s2/_datos/semilla-oraculo.json.
func (l *loader) oracleSeedPath() (string, error) {
	for _, c := range []string{l.opt.OracleSeedPath, filepath.Join(l.opt.Dir, "semilla-oraculo.json"),
		filepath.Join(l.opt.Dir, "..", "..", "oraculo", "fixtures", "s2", "_datos", "semilla-oraculo.json")} {
		if c == "" {
			continue
		}
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c, nil
		}
	}
	return "", errors.New("semilla: no encuentro oraculo/fixtures/s2/_datos/semilla-oraculo.json (los gastos proyectados 2 a 6 salen de ese archivo y no se inventan)")
}

// codeOf traduce el código de categoría de la fixture («despensa») al de la base («groceries»).
func codeOf(fx string) string {
	if fx == "despensa" {
		return "groceries"
	}
	return fx
}

// proyectados agrega los gastos 2 a 6 de semilla-oraculo.json y liga las despensas 2 a 4 con su compra a crédito.
func (l *loader) proyectados() error {
	p, err := l.oracleSeedPath()
	if err != nil {
		return err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	var f oracleSeed
	if err := json.Unmarshal(b, &f); err != nil {
		return fmt.Errorf("semilla: %s: %w", filepath.Base(p), err)
	}
	m, err := readModel(l.opt.Dir, l.opt.ModelPath)
	if err != nil {
		return err
	}
	// Categoría «Papel» (la única de la fixture que la semilla real no traía).
	for _, c := range f.Categories {
		if codeOf(c.Code) != "papel" {
			continue
		}
		id, _, err := l.ensure("categorías", "category", `SELECT id FROM category WHERE code = ?`, a("papel"),
			`INSERT INTO category (parent_id, code, name, kind, default_cost_share_rule_id, sort_order, is_active) VALUES (NULL, 'papel', ?, ?, ?, 3, ?)`,
			c.Name, c.Kind, l.id("rule:karen"), c.IsActive)
		if err != nil {
			return err
		}
		l.ids["category:papel"] = id
	}
	if l.ids["category:papel"] == 0 {
		return errors.New("semilla: semilla-oraculo.json ya no trae la categoría «papel»")
	}
	sort.Slice(f.Expenses, func(i, j int) bool { return f.Expenses[i].ID < f.Expenses[j].ID })
	cid := l.id("credit:kueski")
	n := 0
	for _, e := range f.Expenses {
		if e.ID < 2 { // el gasto 1 es la despensa real del 10-oct, que ya cargó compra10Oct
			continue
		}
		n++
		cat, ok := l.ids["category:"+codeOf(e.Category)]
		if !ok {
			return fmt.Errorf("semilla: el gasto %d de semilla-oraculo.json usa la categoría «%s», que no existe", e.ID, e.Category)
		}
		var paidBy any
		if e.PaidBy != "" {
			pid, ok := l.ids["person:"+e.PaidBy]
			if !ok {
				return fmt.Errorf("semilla: el gasto %d lo paga «%s», que no es una persona sembrada", e.ID, e.PaidBy)
			}
			paidBy = pid
		}
		var sum int64
		for _, s := range e.Shares {
			sum += s.AmountCents
		}
		if sum != e.AmountCents {
			return fmt.Errorf("semilla: el gasto %d de semilla-oraculo.json: sus partes suman %d y el importe es %d", e.ID, sum, e.AmountCents)
		}
		if err := l.cotejaProyectado(e.ID, e.Category, e.AmountCents, e.ExpenseDate, e.Shares, m); err != nil {
			return err
		}
		var desc string
		var rule any
		switch e.Category {
		case "despensa":
			desc = fmt.Sprintf("Despensa programada del %s (proyectado)", e.ExpenseDate)
		default:
			desc = fmt.Sprintf("Papel de contado del %s (proyectado)", e.ExpenseDate)
			rule = l.id("rule:karen")
		}
		eid, created, err := l.ensure("gastos proyectados", "expense",
			`SELECT id FROM expense WHERE expense_date = ? AND category_id = ? AND amount_cents = ? AND payment_method = ?`, a(e.ExpenseDate, cat, e.AmountCents, e.PaymentMethod),
			`INSERT INTO expense (kind, description, category_id, amount_cents, expense_date, accrual_months, payment_method, paid_by_person_id, person_id, cost_share_rule_id, store_id, notes)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULL, ?, NULL, ?)`,
			e.Kind, desc, cat, e.AmountCents, e.ExpenseDate, e.AccrualMonths, e.PaymentMethod, paidBy, rule, proyNota)
		if err != nil {
			return err
		}
		if created {
			for _, s := range e.Shares {
				pid, ok := l.ids["person:"+s.Person]
				if !ok {
					return fmt.Errorf("semilla: el gasto %d reparte a «%s», que no es una persona sembrada", e.ID, s.Person)
				}
				if _, err := l.insert("expense_share", `INSERT INTO expense_share (expense_id, person_id, amount_cents) VALUES (?, ?, ?)`, eid, pid, s.AmountCents); err != nil {
					return err
				}
			}
		}
		if e.CreditPurchaseID == nil {
			continue
		}
		// La despensa se liga a su compra a crédito (misma fecha e importe); si ya estaba ligada no se toca.
		var pid int64
		var cur sql.NullInt64
		if err := l.tx.QueryRowContext(l.ctx, `SELECT id, expense_id FROM credit_purchase WHERE credit_id = ? AND purchase_date = ? AND kind = 'installment_plan' AND principal_cents = ?`,
			cid, e.ExpenseDate, e.AmountCents).Scan(&pid, &cur); err != nil {
			return fmt.Errorf("semilla: el gasto %d no encuentra su compra a crédito del %s por %d: %w", e.ID, e.ExpenseDate, e.AmountCents, err)
		}
		if cur.Valid {
			if cur.Int64 != eid {
				return fmt.Errorf("semilla: la compra a crédito %d ya está ligada al gasto %d y no al %d", pid, cur.Int64, eid)
			}
			continue
		}
		before, err := dominio.RowJSON(l.ctx, l.tx, "credit_purchase", pid)
		if err != nil {
			return err
		}
		if _, err := l.tx.ExecContext(l.ctx, `UPDATE credit_purchase SET expense_id = ? WHERE id = ?`, eid, pid); err != nil {
			return fmt.Errorf("semilla: ligar la compra a crédito %d con su gasto: %w", pid, err)
		}
		after, err := dominio.RowJSON(l.ctx, l.tx, "credit_purchase", pid)
		if err != nil {
			return err
		}
		if err := l.audit("update", "credit_purchase", pid, before, after); err != nil {
			return err
		}
		l.rep.count("compras a crédito ligadas a su gasto", true)
	}
	if n != 5 {
		return fmt.Errorf("semilla: semilla-oraculo.json trae %d gastos después del primero y se esperaban 5 (despensas 2 a 4 y 2 papeles)", n)
	}
	l.rep.Notices = append(l.rep.Notices, "Gastos PROYECTADOS (despensas 2 a 4 y 2 papeles de contado de 199.00, fechas futuras) sembrados porque el modelo canon los asume; el esquema no distingue «proyectado» y se marcan por su descripción y notas. PROVISIONAL hasta la puerta 4 (bloqueo WF3a-bis-semilla-dev, punto 4).")
	return nil
}

// cotejaProyectado compara el gasto proyectado con oraculo/model.json: una despensa con su capital y su reparto, un
// papel con la fecha y el monto de contado de Karen. Una diferencia es un HALLAZGO y la carga se detiene.
func (l *loader) cotejaProyectado(id int64, cat string, cents int64, date string, shares []struct {
	Person      string `json:"person"`
	AmountCents int64  `json:"amount_cents"`
}, m *modelFile) error {
	if cat == "despensa" {
		for _, d := range m.Despensas {
			if d.Fecha != date {
				continue
			}
			want := map[string]int64{"axel": int64(d.Axel.Capital), "karen": int64(d.Karen.Capital)}
			if int64(d.Capital) != cents {
				return fmt.Errorf("semilla: HALLAZGO: el gasto %d (%s) trae %d y la despensa %d de oraculo/model.json %d", id, date, cents, d.N, int64(d.Capital))
			}
			for _, s := range shares {
				if want[s.Person] != s.AmountCents {
					return fmt.Errorf("semilla: HALLAZGO: el gasto %d (%s) reparte %d a %s y model.json %d", id, date, s.AmountCents, s.Person, want[s.Person])
				}
			}
			l.rep.Checks = append(l.rep.Checks, fmt.Sprintf("gasto proyectado %d (despensa %d, %s) = oraculo/model.json", id, d.N, date))
			return nil
		}
		return fmt.Errorf("semilla: HALLAZGO: el gasto %d (%s) no corresponde a ninguna despensa de oraculo/model.json", id, date)
	}
	var papel []struct {
		Compra string `json:"compra"`
		Forma  string `json:"forma"`
		Monto  money  `json:"monto"`
	}
	if err := json.Unmarshal(m.Papel, &papel); err != nil {
		return fmt.Errorf("semilla: model.json papel: %w", err)
	}
	for _, p := range papel {
		if p.Compra == date && strings.HasPrefix(p.Forma, "de contado") && int64(p.Monto) == cents {
			l.rep.Checks = append(l.rep.Checks, fmt.Sprintf("gasto proyectado %d (papel de contado, %s) = oraculo/model.json", id, date))
			return nil
		}
	}
	return fmt.Errorf("semilla: HALLAZGO: el gasto %d (papel %s por %d) no está como papel de contado en oraculo/model.json", id, date, cents)
}
