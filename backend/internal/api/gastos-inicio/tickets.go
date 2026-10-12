package gastos

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/precios/ocr"
)

// lowConfidence es la confianza de los renglones propuestos mientras la lectura no guarde la confianza por renglón
// (esquema 1.4 solo guarda ocr_text). PROVISIONAL, bloqueo S6-gastos-inicio §2.
const lowConfidence = "low"

// toTicketScan arma TicketScan: renglones propuestos con el analizador del motor de OCR (precios/ocr) sobre ocr_text,
// sin los renglones de resumen (TOTAL, IVA…), y el borrador K3 si existe.
func toTicketScan(t dominio.TicketScanRow) (api.TicketScan, error) {
	out := api.TicketScan{ID: t.ID, ExpenseID: nullInt(t.ExpenseID), Kind: t.Kind, OriginalName: nullStr(t.OriginalName), MimeType: t.MimeType,
		SizeBytes: t.SizeBytes, SHA256: t.SHA256, OCRStatus: t.OCRStatus, OCREngine: nullStr(t.OCREngine), OCRText: nullStr(t.OCRText),
		UploadedByUserID: nullInt(t.UploadedByUserID), CreatedAt: api.APITime(t.CreatedAt), UpdatedAt: api.APITime(t.UpdatedAt),
		Lines: []api.TicketScanLinesItem{}}
	n := int64(0)
	for _, raw := range strings.Split(t.OCRText.String, "\n") {
		if strings.TrimSpace(raw) == "" {
			continue
		}
		l := ocr.ParseLine(raw, ocr.ProfileGeneric)
		if l.EsResumen {
			continue
		}
		n++
		text := strings.TrimSpace(l.Descripcion)
		if text == "" {
			text = strings.TrimSpace(raw)
		}
		out.Lines = append(out.Lines, api.TicketScanLinesItem{LineNo: n, Text: text, QuantityMilli: l.CantidadMilli, LineTotalCents: l.ImporteCents, Confidence: lowConfidence})
	}
	if t.DraftJSON.Valid {
		var d api.TicketScanDraft
		if err := json.Unmarshal([]byte(t.DraftJSON.String), &d); err != nil {
			return out, err
		}
		if d.Lines == nil {
			d.Lines = []api.TicketScanDraftLinesItem{}
		}
		out.HasDraft, out.Draft, out.DisplayStoreID, out.DisplayTotalCents = true, &d, d.StoreID, d.TotalCents
		out.OCRDraftSavedAt = api.Ptr(api.APITime(t.DraftSavedAt.String))
	}
	return out, nil
}

func listTicketScans(reg *api.Registry) api.HandlerFunc {
	cols := map[string]string{"id": "a.id", "created_at": "a.created_at", "updated_at": "a.updated_at", "ocr_status": "a.ocr_status"}
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"id", "created_at", "updated_at", "ocr_status"}, "-created_at")
		if err != nil {
			return err
		}
		rows, err := dominio.ListTicketScans(r.Context(), reg.Deps().DB, p.OrderBy(cols), p.LimitOffset())
		if err != nil {
			return err
		}
		items := make([]api.TicketScan, 0, len(rows))
		for _, t := range rows {
			ts, err := toTicketScan(t)
			if err != nil {
				return err
			}
			ts.Lines, ts.Draft = []api.TicketScanLinesItem{}, nil // la cola no trae los renglones (N-076)
			items = append(items, ts)
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

func writeTicket(reg *api.Registry, w http.ResponseWriter, r *http.Request, id int64) error {
	t, err := dominio.GetTicketScan(r.Context(), reg.Deps().DB, id)
	if err != nil {
		return err
	}
	ts, err := toTicketScan(t)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, ts)
}

func getTicketScan(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "attachment_id")
		if err != nil {
			return err
		}
		return writeTicket(reg, w, r, id)
	}
}

func saveTicketScanDraft(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "attachment_id")
		if err != nil {
			return err
		}
		ctx, db := r.Context(), reg.Deps().DB
		if _, err := dominio.GetTicketScan(ctx, db, id); err != nil {
			return err
		}
		var d api.TicketScanDraft
		if err := reg.DecodeJSON(r, &d); err != nil {
			return err
		}
		var v api.Validator
		v.Required(d.Lines != nil, "lines") // TicketScanDraft.required = [lines] en C1 (hallazgo H3 de REVISION-backend-core-r1)
		if d.ExpenseDate != nil {
			if _, err := v.Date(*d.ExpenseDate, "expense_date"); err != nil {
				return err
			}
		}
		v.Check(d.TotalCents == nil || *d.TotalCents >= 0, "total_cents", "out_of_range", "El total no puede ser negativo.")
		v.Check(d.PaymentMethod == nil || slices.Contains(paymentMethods, *d.PaymentMethod), "payment_method", "invalid_format", "La forma de pago debe ser efectivo, débito, transferencia o crédito.")
		for _, ref := range []struct {
			id         *int64
			table, fld string
		}{{d.StoreID, "store", "store_id"}, {d.CategoryID, "category", "category_id"}, {d.CostShareRuleID, "cost_share_rule", "cost_share_rule_id"}} {
			if ref.id != nil {
				ok, err := dominio.ExpenseRefExists(ctx, db, ref.table, *ref.id)
				if err != nil {
					return err
				}
				v.Check(ok, ref.fld, "not_found", "No existe el registro elegido en "+ref.fld+".")
			}
		}
		seen := map[int64]bool{}
		for _, l := range d.Lines {
			v.Check(l.LineNo >= 1 && !seen[l.LineNo], "lines", "invalid_format", "Cada renglón lleva un número de renglón distinto, desde 1.")
			v.Check((l.QuantityMilli == nil || *l.QuantityMilli > 0) && (l.LineTotalCents == nil || *l.LineTotalCents >= 0) &&
				(l.UnitPriceCents == nil || *l.UnitPriceCents >= 0), "lines", "out_of_range", "Cantidades mayores que 0 e importes no negativos.")
			seen[l.LineNo] = true
		}
		if err := v.Err(); err != nil {
			return err
		}
		if d.Lines == nil {
			d.Lines = []api.TicketScanDraftLinesItem{}
		}
		b, err := json.Marshal(d)
		if err != nil {
			return err
		}
		err = dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
			before, err := dominio.RowJSON(ctx, tx, "attachment", id)
			if err != nil {
				return err
			}
			if err := dominio.SaveTicketScanDraft(ctx, tx, id, string(b), dominio.FormatDBTime(reg.Deps().Clock.Now())); err != nil {
				return err
			}
			after, err := dominio.RowJSON(ctx, tx, "attachment", id)
			if err != nil {
				return err
			}
			return dominio.Audit(ctx, tx, api.Actor(r), "update", "attachment", id, before, after, api.RequestID(r))
		})
		if err != nil {
			return err
		}
		return writeTicket(reg, w, r, id)
	}
}
