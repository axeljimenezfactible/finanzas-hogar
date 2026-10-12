package reportes

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/jobs"
)

// MaxTelegramDocumentBytes: la Bot API de Telegram acepta documentos de hasta 50 MB en sendDocument
// («Bots can currently send files of any type of up to 50 MB in size», core.telegram.org/bots/api#senddocument).
const MaxTelegramDocumentBytes = 50 * 1000 * 1000

// MonthlyJobCode es el job del reporte mensual (S7-10).
const MonthlyJobCode = "monthly_report"

// ReportKeys son las claves de la sección «Reportes y frecuencias» (K47, K50).
var ReportKeys = []string{"reports.daily_send_time", "reports.weekly_weekday", "reports.weekly_send_time",
	"reports.monthly_day_of_month", "reports.monthly_send_time"}

// DocumentSender manda un PDF a un chat (sendDocument, texto plano sin parse_mode). Lo implementa internal/telegram
// (agente s7-avisos); devuelve el código HTTP de Telegram. El error nunca debe traer el token.
type DocumentSender interface {
	SendDocument(ctx context.Context, chatID, fileName string, pdf []byte, caption string) (int, error)
}

var (
	senderMu sync.RWMutex
	sender   DocumentSender
)

// SetSender instala el cliente de Telegram (lo hace cmd/finanzas al arrancar; las pruebas ponen uno falso).
func SetSender(s DocumentSender) { senderMu.Lock(); sender = s; senderMu.Unlock() }

func currentSender() DocumentSender { senderMu.RLock(); defer senderMu.RUnlock(); return sender }

// Service genera, guarda, lista y envía reportes.
type Service struct {
	DB       *sql.DB
	Clock    dominio.Clock
	Log      *slog.Logger
	DataDir  string
	Compiler *Compiler
}

// New arma el servicio con las dependencias de la API.
func New(d *api.Deps) *Service {
	return &Service{DB: d.DB, Clock: d.Clock, Log: d.Log, DataDir: d.Config.DataDir,
		Compiler: &Compiler{Engine: "lualatex", TmpDir: d.Config.TmpDir}}
}

func (s *Service) today() dominio.Date { return dominio.Today(s.Clock) }

// LastClosedMonth es el día 1 del mes anterior a hoy.
func (s *Service) LastClosedMonth() dominio.Date {
	t := s.today()
	return dominio.MonthStart(t.Year, t.Month-1)
}

// Monthly arma datos y documento del mes.
func (s *Service) Monthly(ctx context.Context, month dominio.Date) (*Monthly, error) {
	return BuildMonthly(ctx, s.DB, s.today(), month)
}

// GenerateMonthly compila y guarda el PDF del mes (versionado). Si LaTeX falla no guarda nada.
func (s *Service) GenerateMonthly(ctx context.Context, month dominio.Date, actor *int64, reqID string) (api.GeneratedReport, error) {
	m, err := s.Monthly(ctx, month)
	if err != nil {
		return api.GeneratedReport{}, err
	}
	pdf, err := s.Compiler.Compile(ctx, "mensual", MonthlyDoc(m))
	if err != nil {
		return api.GeneratedReport{}, err
	}
	start := dominio.MonthStart(month.Year, month.Month)
	mk := start.MonthKey()
	id, err := s.store(ctx, "monthly_pdf", start, dominio.MonthEnd(start), nil, filepath.Join("mensual", mk), "reporte-mensual-"+mk, pdf, actor, reqID)
	if err != nil {
		return api.GeneratedReport{}, err
	}
	return s.Get(ctx, id)
}

// GenerateExternal compila y guarda el estado de cuenta (nunca se envía: Axel lo reenvía a mano).
func (s *Service) GenerateExternal(ctx context.Context, pid int64, start, end dominio.Date, actor *int64, reqID string) (api.GeneratedReport, error) {
	x, err := BuildExternal(ctx, s.DB, s.today(), pid, start, end)
	if err != nil {
		return api.GeneratedReport{}, err
	}
	pdf, err := s.Compiler.Compile(ctx, "estado-cuenta", ExternalDoc(x))
	if err != nil {
		return api.GeneratedReport{}, err
	}
	var code string
	if err := s.DB.QueryRowContext(ctx, `SELECT code FROM person WHERE id = ?`, pid).Scan(&code); err != nil {
		return api.GeneratedReport{}, err
	}
	id, err := s.store(ctx, "external_statement_pdf", start, end, &pid, filepath.Join("estados-de-cuenta", code),
		"estado-de-cuenta-"+code+"-"+start.String()+"-al-"+end.String(), pdf, actor, reqID)
	if err != nil {
		return api.GeneratedReport{}, err
	}
	return s.Get(ctx, id)
}

// store escribe el PDF en <DataDir>/reportes/<dir>/<base>-vN.pdf (N = versión del mismo periodo) y su renglón.
func (s *Service) store(ctx context.Context, kind string, start, end dominio.Date, pid *int64, dir, base string, pdf []byte, actor *int64, reqID string) (int64, error) {
	if s.DataDir == "" {
		return 0, dominio.Unavailable("Falta la carpeta de datos para guardar el PDF.")
	}
	var n int64
	var pidv any
	if pid != nil {
		pidv = *pid
	}
	if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM generated_report WHERE kind = ? AND period_start_date = ? AND period_end_date = ?
		AND person_id IS ?`, kind, start.String(), end.String(), pidv).Scan(&n); err != nil {
		return 0, err
	}
	abs := filepath.Join(s.DataDir, "reportes", dir)
	if err := os.MkdirAll(abs, 0o700); err != nil {
		return 0, err
	}
	var rel, full string
	var f *os.File
	for v := n + 1; ; v++ {
		rel = filepath.ToSlash(filepath.Join("reportes", dir, base+"-v"+strconv.FormatInt(v, 10)+".pdf"))
		full = filepath.Join(s.DataDir, filepath.FromSlash(rel))
		var err error
		if f, err = os.OpenFile(full, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600); err == nil {
			break
		} else if !errors.Is(err, os.ErrExist) {
			return 0, err
		}
	}
	if _, err := f.Write(pdf); err != nil {
		f.Close()
		os.Remove(full)
		return 0, err
	}
	if err := f.Close(); err != nil {
		os.Remove(full)
		return 0, err
	}
	sum := sha256.Sum256(pdf)
	var id int64
	err := dominio.WithTx(ctx, s.DB, func(tx *sql.Tx) error {
		if err := tx.QueryRowContext(ctx, `INSERT INTO generated_report (kind, period_start_date, period_end_date, person_id, file_path, sha256, size_bytes, generated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`, kind, start.String(), end.String(), pidv, rel, hex.EncodeToString(sum[:]), len(pdf),
			dominio.FormatDBTime(s.Clock.Now())).Scan(&id); err != nil {
			return err
		}
		after, err := dominio.RowJSON(ctx, tx, "generated_report", id)
		if err != nil {
			return err
		}
		return dominio.Audit(ctx, tx, actor, "insert", "generated_report", id, nil, after, reqID)
	})
	if err != nil {
		os.Remove(full)
		return 0, err
	}
	return id, nil
}

// Get devuelve el reporte como en C1.
func (s *Service) Get(ctx context.Context, id int64) (api.GeneratedReport, error) {
	r, _, err := s.get(ctx, id)
	return r, err
}

const reportCols = `id, kind, period_start_date, period_end_date, person_id, file_path, sha256, size_bytes, generated_at, sent_at, created_at, updated_at`

func (s *Service) get(ctx context.Context, id int64) (api.GeneratedReport, string, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT `+reportCols+` FROM generated_report WHERE id = ?`, id)
	if err != nil {
		return api.GeneratedReport{}, "", err
	}
	list, paths, err := s.scan(ctx, rows)
	if err != nil {
		return api.GeneratedReport{}, "", err
	}
	if len(list) == 0 {
		return api.GeneratedReport{}, "", dominio.NotFound("No existe ese reporte.")
	}
	return list[0], paths[0], nil
}

func (s *Service) scan(ctx context.Context, rows *sql.Rows) ([]api.GeneratedReport, []string, error) {
	var out []api.GeneratedReport
	var paths []string
	func() {
		defer rows.Close()
		for rows.Next() {
			var r api.GeneratedReport
			var pid sql.NullInt64
			var sent sql.NullString
			var path string
			if err := rows.Scan(&r.ID, &r.Kind, &r.PeriodStartDate, &r.PeriodEndDate, &pid, &path, &r.SHA256, &r.SizeBytes, &r.GeneratedAt,
				&sent, &r.CreatedAt, &r.UpdatedAt); err != nil {
				out = nil
				return
			}
			r.GeneratedAt, r.CreatedAt, r.UpdatedAt = api.APITime(r.GeneratedAt), api.APITime(r.CreatedAt), api.APITime(r.UpdatedAt)
			if pid.Valid {
				r.PersonID = api.Ptr(pid.Int64)
			}
			if sent.Valid {
				r.SentAt = api.Ptr(api.APITime(sent.String))
			}
			out = append(out, r)
			paths = append(paths, path)
		}
	}()
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}
	for i := range out {
		st, err := s.deliveryStatus(ctx, out[i].ID)
		if err != nil {
			return nil, nil, err
		}
		out[i].DeliveryStatus = st
	}
	return out, paths, nil
}

// deliveryStatus resume los envíos del reporte (enum de C1): sent (todos), failed (alguno falló) o not_sent.
func (s *Service) deliveryStatus(ctx context.Context, id int64) (*string, error) {
	var total, sent, failed int
	if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(*), COALESCE(SUM(status = 'sent'), 0), COALESCE(SUM(status = 'failed'), 0)
		FROM notification_delivery WHERE dedupe_key = ?`, dedupeKey(id)).Scan(&total, &sent, &failed); err != nil {
		return nil, err
	}
	switch {
	case total > 0 && failed > 0:
		return api.Ptr("failed"), nil
	case total > 0 && sent == total:
		return api.Ptr("sent"), nil
	default:
		return api.Ptr("not_sent"), nil
	}
}

func dedupeKey(id int64) string { return "monthly_report:report:" + strconv.FormatInt(id, 10) }

// File devuelve la ruta absoluta del PDF guardado y un nombre de descarga.
func (s *Service) File(ctx context.Context, id int64) (string, string, error) {
	r, rel, err := s.get(ctx, id)
	if err != nil {
		return "", "", err
	}
	name := filepath.Base(rel)
	if r.Kind == "monthly_pdf" {
		name = "reporte-mensual-" + r.PeriodStartDate[:7] + ".pdf"
	}
	return filepath.Join(s.DataDir, filepath.FromSlash(rel)), name, nil
}

type recipient struct {
	id     int64
	chatID string
}

// Send manda el reporte mensual por Telegram a los destinatarios suscritos a monthly_report (nunca a personas externas;
// S7-09) y anota cada envío en notification_delivery. Un destinatario que ya lo recibió no se repite.
func (s *Service) Send(ctx context.Context, id int64, actor *int64, reqID string) (api.GeneratedReport, error) {
	r, rel, err := s.get(ctx, id)
	if err != nil {
		return r, err
	}
	if r.Kind != "monthly_pdf" {
		return r, dominio.Conflict("conflict", "El estado de cuenta no se envía por Telegram: Axel lo reenvía a mano.")
	}
	snd := currentSender()
	if snd == nil {
		return r, dominio.Unavailable("El envío por Telegram no está configurado todavía.")
	}
	pdf, err := os.ReadFile(filepath.Join(s.DataDir, filepath.FromSlash(rel)))
	if err != nil {
		return r, dominio.Conflict("conflict", "No se encontró el PDF guardado; genera el reporte otra vez.")
	}
	if sum := sha256.Sum256(pdf); hex.EncodeToString(sum[:]) != r.SHA256 {
		return r, dominio.Conflict("conflict", "El PDF guardado no coincide con su huella; genera el reporte otra vez.")
	}
	if len(pdf) > MaxTelegramDocumentBytes {
		return r, dominio.Conflict("conflict", "El PDF pasa de 50 MB, el límite de Telegram.")
	}
	var notifID int64
	if err := s.DB.QueryRowContext(ctx, `SELECT id FROM notification WHERE code = 'monthly_report'`).Scan(&notifID); err != nil {
		return r, fmt.Errorf("reportes: aviso monthly_report: %w", err)
	}
	rows, err := s.DB.QueryContext(ctx, `SELECT nr.id, nr.chat_id FROM notification_subscription ns
		JOIN notification n ON n.id = ns.notification_id JOIN notification_recipient nr ON nr.id = ns.recipient_id
		JOIN person p ON p.id = nr.person_id
		WHERE n.id = ? AND nr.is_active = 1 AND p.participation <> 'external' AND p.is_active = 1
		  AND NOT EXISTS (SELECT 1 FROM notification_delivery d WHERE d.dedupe_key = ? AND d.recipient_id = nr.id AND d.status = 'sent')
		ORDER BY nr.id`, notifID, dedupeKey(id))
	if err != nil {
		return r, err
	}
	var rs []recipient
	for rows.Next() {
		var x recipient
		if err := rows.Scan(&x.id, &x.chatID); err != nil {
			rows.Close()
			return r, err
		}
		rs = append(rs, x)
	}
	rows.Close()
	mo, _ := dominio.ParseDate(r.PeriodStartDate)
	caption := "Reporte mensual de " + MesLargo(mo.Year, mo.Month) + " (PDF)."
	failed := 0
	for _, x := range rs {
		code, err := snd.SendDocument(ctx, x.chatID, "reporte-mensual-"+mo.MonthKey()+".pdf", pdf, caption)
		status, msg := "sent", ""
		if err != nil || code != 200 {
			status, msg, failed = "failed", "Telegram respondió HTTP "+strconv.Itoa(code), failed+1
			if code == 0 {
				msg = "No se pudo conectar con Telegram."
			}
		}
		var sentAt any
		if status == "sent" {
			sentAt = dominio.FormatDBTime(s.Clock.Now())
		}
		var httpStatus any
		if code > 0 {
			httpStatus = code
		}
		if _, err := s.DB.ExecContext(ctx, `INSERT INTO notification_delivery (notification_id, recipient_id, dedupe_key, status, message_text, http_status, error, sent_at)
			VALUES (?, ?, ?, ?, ?, ?, NULLIF(?, ''), ?) ON CONFLICT (dedupe_key, recipient_id) DO UPDATE SET status = excluded.status,
			message_text = excluded.message_text, http_status = excluded.http_status, error = excluded.error, sent_at = excluded.sent_at`,
			notifID, x.id, dedupeKey(id), status, caption, httpStatus, msg, sentAt); err != nil {
			return r, err
		}
	}
	if failed == 0 && len(rs) > 0 {
		err := dominio.WithTx(ctx, s.DB, func(tx *sql.Tx) error {
			before, err := dominio.RowJSON(ctx, tx, "generated_report", id)
			if err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `UPDATE generated_report SET sent_at = ? WHERE id = ?`, dominio.FormatDBTime(s.Clock.Now()), id); err != nil {
				return err
			}
			after, err := dominio.RowJSON(ctx, tx, "generated_report", id)
			if err != nil {
				return err
			}
			return dominio.Audit(ctx, tx, actor, "update", "generated_report", id, before, after, reqID)
		})
		if err != nil {
			return r, err
		}
	}
	out, err := s.Get(ctx, id)
	if err != nil {
		return out, err
	}
	if failed > 0 {
		return out, fmt.Errorf("reportes: %d de %d envíos por Telegram fallaron", failed, len(rs))
	}
	return out, nil
}

// EffectiveDay es el día en que sale el reporte en ese mes: el configurado o el último del mes (Q2-66 B).
func EffectiveDay(y int, m time.Month, day int64) int {
	last := dominio.MonthEnd(dominio.MonthStart(y, m)).Day
	if int(day) > last {
		return last
	}
	return int(day)
}

// MonthlySchedule da el horario cron del job a partir de reports.monthly_day_of_month y reports.monthly_send_time.
// Con día 29–31 corre los días 28 a 31 y el job solo actúa en el día efectivo del mes.
func MonthlySchedule(ctx context.Context, q dominio.Querier) (string, error) {
	day, err := dominio.SettingInt(ctx, q, "reports.monthly_day_of_month")
	if err != nil {
		return "", err
	}
	hhmm, err := dominio.SettingTime(ctx, q, "reports.monthly_send_time")
	if err != nil {
		return "", err
	}
	hh, mm, ok := strings.Cut(hhmm, ":")
	h, err1 := strconv.Atoi(hh)
	mi, err2 := strconv.Atoi(mm)
	if !ok || err1 != nil || err2 != nil {
		return "", fmt.Errorf("reportes: hora inválida %q", hhmm)
	}
	days := strconv.FormatInt(day, 10)
	if day > 28 {
		days = "28-31"
	}
	return fmt.Sprintf("%d %d %s * *", mi, h, days), nil
}

// SyncMonthlyJob deja el renglón del job con el horario de los ajustes y su siguiente corrida.
func (s *Service) SyncMonthlyJob(ctx context.Context, q dominio.Querier) error {
	cron, err := MonthlySchedule(ctx, q)
	if err != nil {
		return err
	}
	sched, err := jobs.ParseSchedule(cron)
	if err != nil {
		return err
	}
	// Si el job ya corrió hoy, la siguiente corrida se cuenta desde el fin del día: cambiar la hora no repite el reporte.
	base := s.Clock.Now()
	var last sql.NullString
	if err := q.QueryRowContext(ctx, `SELECT last_run_at FROM job WHERE code = ?`, MonthlyJobCode).Scan(&last); err != nil && err != sql.ErrNoRows {
		return err
	}
	if last.Valid {
		if t, err := dominio.ParseDBTime(last.String); err == nil && dominio.DateOf(t) == dominio.DateOf(base) {
			base = dominio.DateOf(base).AddDays(1).StartOfDay().Add(-time.Second)
		}
	}
	_, err = q.ExecContext(ctx, `UPDATE job SET schedule = ?, next_run_at = ? WHERE code = ?`, cron,
		dominio.FormatDBTime(sched.Next(base)), MonthlyJobCode)
	return err
}

// RunMonthlyJob es el job del día configurado: genera el PDF del mes anterior, lo guarda y lo envía. Si LaTeX falla,
// el job termina en error y no se envía nada (S7-10).
func (s *Service) RunMonthlyJob(ctx context.Context, run jobs.Run) error {
	day, err := dominio.SettingInt(ctx, s.DB, "reports.monthly_day_of_month")
	if err != nil {
		return err
	}
	d := dominio.DateOf(run.Scheduled)
	if d.Day != EffectiveDay(d.Year, d.Month, day) {
		return nil
	}
	month := dominio.MonthStart(d.Year, d.Month-1)
	rep, err := s.GenerateMonthly(ctx, month, nil, run.RequestID)
	if err != nil {
		return fmt.Errorf("no se generó el reporte mensual de %s: %w", MesLargo(month.Year, month.Month), err)
	}
	if _, err := s.Send(ctx, rep.ID, nil, run.RequestID); err != nil {
		return fmt.Errorf("el reporte se guardó pero no se envió: %w", err)
	}
	return nil
}

// Settings lee la sección «Reportes y frecuencias».
func Settings(ctx context.Context, q dominio.Querier) (api.ReportSettings, error) {
	var out api.ReportSettings
	tp := func(key string) (api.SettingParameterTime, error) {
		p, err := dominio.ReadParameter(ctx, q, key)
		if err != nil {
			return api.SettingParameterTime{}, err
		}
		o := api.SettingParameterTime{Key: p.Key, ValueType: p.ValueType, Unit: p.Unit, IsInitial: p.IsInitial,
			Value: api.Null[string](), InitialValue: api.Null[string]()}
		var v string
		if json.Unmarshal(p.Value, &v) == nil && len(p.Value) > 0 && string(p.Value) != "null" {
			o.Value = api.Some(v)
		}
		if json.Unmarshal(p.InitialValue, &v) == nil && len(p.InitialValue) > 0 && string(p.InitialValue) != "null" {
			o.InitialValue = api.Some(v)
		}
		return o, nil
	}
	ip := func(key string) (api.SettingParameterInteger, error) {
		p, err := dominio.ReadParameter(ctx, q, key)
		if err != nil {
			return api.SettingParameterInteger{}, err
		}
		o := api.SettingParameterInteger{Key: p.Key, ValueType: p.ValueType, Unit: p.Unit, IsInitial: p.IsInitial,
			Value: api.Null[int64](), InitialValue: api.Null[int64]()}
		var v int64
		if json.Unmarshal(p.Value, &v) == nil && len(p.Value) > 0 && string(p.Value) != "null" {
			o.Value = api.Some(v)
		}
		if json.Unmarshal(p.InitialValue, &v) == nil && len(p.InitialValue) > 0 && string(p.InitialValue) != "null" {
			o.InitialValue = api.Some(v)
		}
		return o, nil
	}
	var err error
	pr := &out.Parameters
	if pr.ReportsDailySendTime, err = tp("reports.daily_send_time"); err != nil {
		return out, err
	}
	if pr.ReportsWeeklyWeekday, err = ip("reports.weekly_weekday"); err != nil {
		return out, err
	}
	if pr.ReportsWeeklySendTime, err = tp("reports.weekly_send_time"); err != nil {
		return out, err
	}
	if pr.ReportsMonthlyDayOfMonth, err = ip("reports.monthly_day_of_month"); err != nil {
		return out, err
	}
	if pr.ReportsMonthlySendTime, err = tp("reports.monthly_send_time"); err != nil {
		return out, err
	}
	return out, nil
}

// List corre la consulta de listReports (columnas de reportCols) y arma la página.
func (s *Service) List(ctx context.Context, query string, args ...any) ([]api.GeneratedReport, error) {
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	items, _, err := s.scan(ctx, rows)
	if items == nil {
		items = []api.GeneratedReport{}
	}
	return items, err
}
