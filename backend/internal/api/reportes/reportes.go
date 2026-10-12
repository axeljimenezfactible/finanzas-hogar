// Package reportesapi expone las operaciones de Reportes de C1 (agente s7-pdf): reportes en PDF, vista web mensual,
// estado de cuenta de personas externas y la sección de ajustes «Reportes y frecuencias».
package reportesapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/jobs"
	"finanzas-hogar/internal/reportes"
)

// Register registra las 10 operaciones de s7-pdf y el job del reporte mensual.
func Register(reg *api.Registry) {
	s := reportes.New(reg.Deps())
	h := &handlers{reg: reg, s: s}
	reg.Handle("listReports", h.listReports)
	reg.Handle("generateMonthlyReport", h.generateMonthlyReport)
	reg.Handle("generateExternalStatement", h.generateExternalStatement)
	reg.Handle("downloadReportFile", h.downloadReportFile)
	reg.Handle("sendReport", h.sendReport)
	reg.Handle("listMonthlyReportMonths", h.listMonthlyReportMonths)
	reg.Handle("getMonthlyReport", h.getMonthlyReport)
	reg.Handle("getExternalStatement", h.getExternalStatement)
	reg.Handle("getReportSettings", h.getReportSettings)
	reg.Handle("updateReportSettings", h.updateReportSettings)
	if j := reg.Deps().Jobs; j != nil && reg.Deps().DB != nil {
		cron, err := reportes.MonthlySchedule(context.Background(), reg.Deps().DB)
		if err != nil {
			cron = "0 8 1 * *" // solo si la base aún no tiene las claves: decisión F y K50 (día 1, 08:00)
		}
		_ = j.Register(jobs.Def{Code: reportes.MonthlyJobCode, Name: "Reporte mensual en PDF", Schedule: cron, Run: s.RunMonthlyJob})
	}
}

type handlers struct {
	reg *api.Registry
	s   *reportes.Service
}

func invalid(field, code, msg string) error {
	return dominio.Invalid(dominio.FieldIssue{Field: field, Code: code, Message: msg})
}

// latexErr: si LaTeX falla, 503 con mensaje simple (el detalle queda en el registro).
func (h *handlers) latexErr(err error) error {
	if errors.Is(err, reportes.ErrLatex) {
		h.reg.Deps().Log.Error("reportes: LaTeX", "error", err.Error())
		return dominio.Unavailable("No se pudo generar el PDF. Revisa que LaTeX esté instalado y vuelve a intentar.")
	}
	return err
}

func (h *handlers) listReports(w http.ResponseWriter, r *http.Request) error {
	p, err := api.ParsePage(r, []string{"generated_at", "period_start_date"}, "-generated_at")
	if err != nil {
		return err
	}
	per, err := api.QueryPeriod(r)
	if err != nil {
		return err
	}
	pids, err := api.QueryIDs(r, "person_id")
	if err != nil {
		return err
	}
	q := `SELECT id, kind, period_start_date, period_end_date, person_id, file_path, sha256, size_bytes, generated_at, sent_at, created_at, updated_at
		FROM generated_report g WHERE 1 = 1`
	var args []any
	if k := r.URL.Query().Get("kind"); k != "" {
		if k != "monthly_pdf" && k != "external_statement_pdf" {
			return api.BadRequest("kind debe ser monthly_pdf o external_statement_pdf.")
		}
		q += ` AND kind = ?`
		args = append(args, k)
	}
	if per != nil {
		q += ` AND period_end_date >= ? AND period_start_date <= ?`
		args = append(args, per.Start.String(), per.End.String())
	}
	if len(pids) > 0 {
		q += ` AND person_id IN (?` + repeat(",?", len(pids)-1) + `)`
		for _, id := range pids {
			args = append(args, id)
		}
	}
	q += " " + p.OrderBy(map[string]string{"generated_at": "g.generated_at", "period_start_date": "g.period_start_date", "id": "g.id"}) + " " + p.LimitOffset()
	items, err := h.s.List(r.Context(), q, args...)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}

func (h *handlers) generateMonthlyReport(w http.ResponseWriter, r *http.Request) error {
	var in api.MonthlyReportRequest
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	m, err := dominio.ParseDate(in.Month + "-01")
	if err != nil || len(in.Month) != 7 {
		return invalid("month", "invalid_format", "Usa un mes como 2026-10.")
	}
	today := dominio.Today(h.reg.Deps().Clock)
	if m.After(dominio.MonthStart(today.Year, today.Month)) {
		return invalid("month", "out_of_range", "El reporte es de un mes que ya empezó.")
	}
	rep, err := h.s.GenerateMonthly(r.Context(), m, api.Actor(r), api.RequestID(r))
	if err != nil {
		return h.latexErr(err)
	}
	if in.SendTelegram {
		if rep, err = h.s.Send(r.Context(), rep.ID, api.Actor(r), api.RequestID(r)); err != nil {
			var de *dominio.Error
			if errors.As(err, &de) {
				return err
			}
			h.reg.Deps().Log.Warn("reportes: envío", "error", err.Error())
		}
	}
	return api.WriteJSON(w, http.StatusAccepted, rep)
}

func (h *handlers) generateExternalStatement(w http.ResponseWriter, r *http.Request) error {
	var in api.ExternalStatementRequest
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	var v api.Validator
	start, e1 := dominio.ParseDate(in.PeriodStartDate)
	end, e2 := dominio.ParseDate(in.PeriodEndDate)
	v.Check(e1 == nil, "period_start_date", "invalid_format", "Usa una fecha como 2026-10-01.")
	v.Check(e2 == nil, "period_end_date", "invalid_format", "Usa una fecha como 2026-10-31.")
	v.Check(e1 != nil || e2 != nil || !end.Before(start), "period_end_date", "out_of_range", "El fin del periodo no puede ser antes del inicio.")
	v.Check(in.PersonID >= 1, "person_id", "required", "Elige a la persona.")
	if err := v.Err(); err != nil {
		return err
	}
	rep, err := h.s.GenerateExternal(r.Context(), in.PersonID, start, end, api.Actor(r), api.RequestID(r))
	if errors.Is(err, dominio.ErrNotFound) {
		return invalid("person_id", "not_found", "No existe esa persona.")
	}
	var de *dominio.Error
	if errors.As(err, &de) && de.Kind == dominio.KindNotFound {
		return invalid("person_id", "not_found", "No existe esa persona.")
	}
	if err != nil {
		return h.latexErr(err)
	}
	return api.WriteJSON(w, http.StatusAccepted, rep)
}

func (h *handlers) downloadReportFile(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "report_id")
	if err != nil {
		return err
	}
	path, name, err := h.s.File(r.Context(), id)
	if err != nil {
		return err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return dominio.NotFound("El PDF de ese reporte ya no está en disco.")
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
	return nil
}

func (h *handlers) sendReport(w http.ResponseWriter, r *http.Request) error {
	id, err := api.PathID(r, "report_id")
	if err != nil {
		return err
	}
	rep, err := h.s.Send(r.Context(), id, api.Actor(r), api.RequestID(r))
	if err != nil {
		var de *dominio.Error
		if errors.As(err, &de) {
			return err
		}
		h.reg.Deps().Log.Warn("reportes: envío", "error", err.Error())
	}
	return api.WriteJSON(w, http.StatusAccepted, rep)
}

func (h *handlers) listMonthlyReportMonths(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	db := h.reg.Deps().DB
	today := dominio.Today(h.reg.Deps().Clock)
	cur := dominio.MonthStart(today.Year, today.Month)
	first := cur
	var minDate sql.NullString
	if err := db.QueryRowContext(ctx, `SELECT MIN(expense_date) FROM expense`).Scan(&minDate); err != nil {
		return err
	}
	if minDate.Valid {
		if d, err := dominio.ParseDate(minDate.String); err == nil && d.Before(cur) {
			first = dominio.MonthStart(d.Year, d.Month)
		}
	}
	out := api.ReportMonths{Months: []api.ReportMonthsMonthsItem{}}
	for m := first; !m.After(cur); m = dominio.MonthStart(m.Year, m.Month+1) {
		it := api.ReportMonthsMonthsItem{Month: m.MonthKey(), IsPartial: m == cur}
		var id sql.NullInt64
		if err := db.QueryRowContext(ctx, `SELECT MAX(id) FROM generated_report WHERE kind = 'monthly_pdf' AND period_start_date = ?`, m.String()).Scan(&id); err != nil {
			return err
		}
		if id.Valid {
			it.GeneratedReportID = api.Ptr(id.Int64)
		}
		out.Months = append(out.Months, it)
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) getMonthlyReport(w http.ResponseWriter, r *http.Request) error {
	per, err := api.QueryPeriod(r)
	if err != nil {
		return err
	}
	m := h.s.LastClosedMonth()
	if per != nil {
		if per.Kind != dominio.PeriodMonth {
			return invalid("period", "not_allowed", "Esta vista es por mes: usa un periodo como 2026-10.")
		}
		m = per.Start
	}
	mr, err := h.s.Monthly(r.Context(), m)
	if err != nil {
		return err
	}
	var id sql.NullInt64
	if err := h.reg.Deps().DB.QueryRowContext(r.Context(), `SELECT MAX(id) FROM generated_report WHERE kind = 'monthly_pdf' AND period_start_date = ?`,
		dominio.MonthStart(m.Year, m.Month).String()).Scan(&id); err != nil {
		return err
	}
	if id.Valid {
		rep, err := h.s.Get(r.Context(), id.Int64)
		if err != nil {
			return err
		}
		mr.Report.PDF = &rep
	}
	return api.WriteJSON(w, http.StatusOK, mr.Report)
}

func (h *handlers) getExternalStatement(w http.ResponseWriter, r *http.Request) error {
	pid, err := api.PathID(r, "person_id")
	if err != nil {
		return err
	}
	per, err := api.QueryPeriod(r)
	if err != nil {
		return err
	}
	today := dominio.Today(h.reg.Deps().Clock)
	start := dominio.MonthStart(today.Year, today.Month)
	end := dominio.MonthEnd(start)
	if per != nil {
		start, end = per.Start, per.End
	}
	x, err := reportes.BuildExternal(r.Context(), h.reg.Deps().DB, today, pid, start, end)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, x.Statement)
}

func (h *handlers) getReportSettings(w http.ResponseWriter, r *http.Request) error {
	out, err := reportes.Settings(r.Context(), h.reg.Deps().DB)
	if err != nil {
		return err
	}
	return api.WriteJSON(w, http.StatusOK, out)
}

func (h *handlers) updateReportSettings(w http.ResponseWriter, r *http.Request) error {
	var in api.ReportSettingsUpdate
	if err := h.reg.DecodeJSON(r, &in); err != nil {
		return err
	}
	ch := dominio.ParamChange{Set: map[string]json.RawMessage{}, Reset: in.ResetParameters}
	if p := in.Parameters; p != nil {
		add := func(k string, v any) {
			b, _ := json.Marshal(v)
			ch.Set[k] = b
		}
		if p.ReportsDailySendTime != nil {
			add("reports.daily_send_time", *p.ReportsDailySendTime)
		}
		if p.ReportsWeeklyWeekday != nil {
			add("reports.weekly_weekday", *p.ReportsWeeklyWeekday)
		}
		if p.ReportsWeeklySendTime != nil {
			add("reports.weekly_send_time", *p.ReportsWeeklySendTime)
		}
		if p.ReportsMonthlyDayOfMonth != nil {
			add("reports.monthly_day_of_month", *p.ReportsMonthlyDayOfMonth)
		}
		if p.ReportsMonthlySendTime != nil {
			add("reports.monthly_send_time", *p.ReportsMonthlySendTime)
		}
	}
	ctx := r.Context()
	err := dominio.WithTx(ctx, h.reg.Deps().DB, func(tx *sql.Tx) error {
		if err := dominio.ApplyParams(ctx, tx, ch, reportes.ReportKeys, api.Actor(r), api.RequestID(r)); err != nil {
			return err
		}
		return h.s.SyncMonthlyJob(ctx, tx)
	})
	if err != nil {
		return err
	}
	return h.getReportSettings(w, r)
}
