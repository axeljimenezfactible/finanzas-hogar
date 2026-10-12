// Package sistema: operaciones de sistema de auth-sistema (salud, contrato, respaldos, semilla, ajustes de
// seguridad y de datos, IA opcional, estado del servidor, certificado de la CA local, adjuntos y jobs). Lo
// registra auth.Register.
package sistema

import (
	"context"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"os"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/ia"
	"finanzas-hogar/internal/jobs"
)

var (
	// Version se fija al compilar: -ldflags "-X finanzas-hogar/internal/api/sistema.Version=…".
	Version = "dev"
	// CACertPath es el rootCA.pem de mkcert que sirve getCaCertificate (serve --ca / FINANZAS_CA). Vacío = 503.
	// PROVISIONAL: va aquí porque api.Config no tiene el campo (solicitud fuera de alcance a backend-base).
	CACertPath string
)

var (
	securityKeys = []string{"auth.max_failed_attempts", "auth.lockout_window_minutes", "auth.password_min_length", "session.max_age_days"}
	dataKeys     = []string{"attachment.max_size_mb", "idempotency.retention_hours", "ocr.timeout_seconds"}
	// aiUses: solo OCR de tickets y sugerencias de recetas (decisión H).
	aiUses = []string{"ticket_ocr", "recipe_suggestions"}
	// seedCodes: las 3 semillas de C1 (SeedStatus).
	seedCodes = []string{"purchase_2026_10_10", "kueski_credit", "fixed_items_csv"}
)

// Register registra las 19 operaciones de sistema.
func Register(reg *api.Registry) {
	d := reg.Deps()
	if d.IA == nil {
		d.IA = ia.NewRegistry(d.DB) // proveedores apagados (S6-13)
	}
	started := d.Clock.Now()
	reg.Handle("getHealth", getHealth(reg))
	reg.Handle("getOpenApiDocument", getOpenAPI(reg))
	reg.Handle("getSecuritySettings", getSettings(reg, securityKeys, securityOut))
	reg.Handle("updateSecuritySettings", updateSettings(reg, securityKeys, securityOut))
	reg.Handle("getDataSettings", getSettings(reg, dataKeys, dataOut))
	reg.Handle("updateDataSettings", updateSettings(reg, dataKeys, dataOut))
	reg.Handle("listAiProviders", listAiProviders(reg))
	reg.Handle("updateAiProvider", updateAiProvider(reg))
	reg.Handle("getSystemStatus", getSystemStatus(reg, started))
	reg.Handle("getCaCertificate", getCaCertificate)
	reg.Handle("getSeedStatus", getSeedStatus(reg))
	reg.Handle("listJobs", listJobs(reg))
	reg.Handle("updateJob", updateJob(reg))
	registerBackups(reg)
	registerAttachments(reg)
}

func getHealth(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Cache-Control", "no-store")
		var one int
		if err := reg.Deps().DB.QueryRowContext(r.Context(), `SELECT 1`).Scan(&one); err != nil {
			return dominio.Unavailable("La base de datos no responde. Revisa que el servidor de la casa esté encendido.")
		}
		return api.WriteJSON(w, http.StatusOK, api.Health{Status: "ok", Time: dominio.FormatDateTime(reg.Deps().Clock.Now())})
	}
}

func getOpenAPI(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		b, err := os.ReadFile(reg.Deps().Config.OpenAPIPath)
		if err != nil {
			return dominio.Unavailable("El contrato de la API no está disponible en este momento.")
		}
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		_, err = w.Write(b)
		return err
	}
}

// --- Ajustes de seguridad y de datos (K44, C1 §11.7) ---

func intParam(ctx context.Context, q dominio.Querier, key string) (api.SettingParameterInteger, error) {
	p, err := dominio.ReadParameter(ctx, q, key)
	if err != nil {
		return api.SettingParameterInteger{}, err
	}
	out := api.SettingParameterInteger{Key: p.Key, ValueType: p.ValueType, Unit: p.Unit, IsInitial: p.IsInitial}
	if out.Value, err = nullInt(p.Value); err != nil {
		return out, err
	}
	out.InitialValue, err = nullInt(p.InitialValue)
	return out, err
}

func nullInt(raw json.RawMessage) (api.Nullable[int64], error) {
	if len(raw) == 0 || string(raw) == "null" {
		return api.Null[int64](), nil
	}
	var v int64
	if err := json.Unmarshal(raw, &v); err != nil {
		return api.Nullable[int64]{}, err
	}
	return api.Some(v), nil
}

func readParams(ctx context.Context, q dominio.Querier, keys []string) ([]api.SettingParameterInteger, error) {
	out := make([]api.SettingParameterInteger, len(keys))
	for i, k := range keys {
		p, err := intParam(ctx, q, k)
		if err != nil {
			return nil, err
		}
		out[i] = p
	}
	return out, nil
}

func securityOut(ps []api.SettingParameterInteger) any {
	return api.SecuritySettings{Parameters: api.SecuritySettingsParameters{AuthMaxFailedAttempts: ps[0], AuthLockoutWindowMinutes: ps[1], AuthPasswordMinLength: ps[2], SessionMaxAgeDays: ps[3]}}
}

func dataOut(ps []api.SettingParameterInteger) any {
	return api.DataSettings{Parameters: api.DataSettingsParameters{AttachmentMaxSizeMB: ps[0], IdempotencyRetentionHours: ps[1], OCRTimeoutSeconds: ps[2]}}
}

func getSettings(reg *api.Registry, keys []string, shape func([]api.SettingParameterInteger) any) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ps, err := readParams(r.Context(), reg.Deps().DB, keys)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, shape(ps))
	}
}

// updateSettings: parameters (valor nuevo por clave) y reset_parameters (vuelven al inicial) con dominio.ApplyParams,
// que valida tipo, clave de la sección, «cambiar y restablecer a la vez» y rango (422) y escribe audit_log.
func updateSettings(reg *api.Registry, keys []string, shape func([]api.SettingParameterInteger) any) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var in struct {
			Parameters      map[string]json.RawMessage `json:"parameters"`
			ResetParameters []string                   `json:"reset_parameters"`
		}
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		if len(in.Parameters) == 0 && len(in.ResetParameters) == 0 {
			return api.BadRequest("No hay cambios que guardar: manda al menos un parámetro.")
		}
		seen := map[string]bool{}
		for _, k := range in.ResetParameters {
			if seen[k] {
				return api.BadRequest("Una clave viene repetida en reset_parameters.")
			}
			seen[k] = true
		}
		ctx, d := r.Context(), reg.Deps()
		err := dominio.WithTx(ctx, d.DB, func(tx *sql.Tx) error {
			return dominio.ApplyParams(ctx, tx, dominio.ParamChange{Set: in.Parameters, Reset: in.ResetParameters}, keys, api.Actor(r), api.RequestID(r))
		})
		if err != nil {
			return err
		}
		return getSettings(reg, keys, shape)(w, r)
	}
}

// --- IA opcional (S6-13) ---

func aiProvider(ctx context.Context, reg *api.Registry, code string) (api.AiProvider, error) {
	reg_ := reg.Deps().IA
	on, err := reg_.Enabled(ctx, code)
	if err != nil {
		return api.AiProvider{}, err
	}
	out := api.AiProvider{Provider: code, IsEnabled: on, AllowedUses: aiUses, IsConfigured: reg_.Available(code)}
	if on { // solo se consulta si está activo
		out.IsAvailable = reg_.Check(ctx, code) == nil
	}
	return out, nil
}

func listAiProviders(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		out := api.AiProviders{Providers: []api.AiProvider{}}
		for _, c := range []string{ia.LocalModel, ia.ClaudeAPI} {
			p, err := aiProvider(r.Context(), reg, c)
			if err != nil {
				return err
			}
			out.Providers = append(out.Providers, p)
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

func updateAiProvider(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		code := r.PathValue("provider")
		if code != ia.LocalModel && code != ia.ClaudeAPI {
			return dominio.NotFound("Ese proveedor de IA no existe.")
		}
		var in struct {
			IsEnabled *bool `json:"is_enabled"`
		}
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		var v api.Validator
		v.Required(in.IsEnabled != nil, "is_enabled")
		if err := v.Err(); err != nil {
			return err
		}
		val, _ := json.Marshal(*in.IsEnabled)
		ctx, d := r.Context(), reg.Deps()
		err := dominio.WithTx(ctx, d.DB, func(tx *sql.Tx) error {
			return dominio.SetSetting(ctx, tx, ia.SettingKey(code), val, api.Actor(r), api.RequestID(r))
		})
		if err != nil {
			return err
		}
		p, err := aiProvider(ctx, reg, code)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, p)
	}
}

// --- Estado del servidor y CA local ---

func getSystemStatus(reg *api.Registry, started time.Time) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var size int64
		var seq int
		var name, file string
		// PRAGMA database_list da la ruta del archivo; no consulta tablas
		if err := reg.Deps().DB.QueryRowContext(r.Context(), `PRAGMA database_list`).Scan(&seq, &name, &file); err != nil {
			return err
		}
		if st, err := os.Stat(file); err == nil {
			size = st.Size()
		}
		return api.WriteJSON(w, http.StatusOK, api.SystemStatus{Version: Version, StartedAt: dominio.FormatDateTime(started), DatabaseSizeBytes: size})
	}
}

// getCaCertificate sirve solo bloques CERTIFICATE: si la ruta apunta a otra cosa (p. ej. la llave de la CA), 503.
func getCaCertificate(w http.ResponseWriter, r *http.Request) error {
	unavailable := dominio.Unavailable("El certificado de la casa no está configurado en el servidor. Pídele a Axel que lo revise.")
	if CACertPath == "" {
		return unavailable
	}
	b, err := os.ReadFile(CACertPath)
	if err != nil {
		return unavailable
	}
	rest, n := b, 0
	for {
		var blk *pem.Block
		blk, rest = pem.Decode(rest)
		if blk == nil {
			break
		}
		if blk.Type != "CERTIFICATE" {
			return unavailable
		}
		n++
	}
	if n == 0 {
		return unavailable
	}
	w.Header().Set("Content-Type", "application/x-pem-file")
	w.Header().Set("Content-Disposition", `attachment; filename="finanzas-hogar-ca.pem"`)
	_, err = w.Write(b)
	return err
}

// getSeedStatus: PROVISIONAL (bloqueo S6-sistema): una semilla está cargada si audit_log tiene altas con
// request_id = «seed:<código>» (convención que debe seguir internal/semilla); row_count = esas altas.
func getSeedStatus(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		out := api.SeedStatus{Items: []api.SeedStatusItemsItem{}}
		for _, c := range seedCodes {
			var n int64
			if err := reg.Deps().DB.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM audit_log WHERE request_id = ? AND action = 'insert'`, "seed:"+c).Scan(&n); err != nil {
				return err
			}
			out.Items = append(out.Items, api.SeedStatusItemsItem{Code: c, IsLoaded: n > 0, RowCount: n})
		}
		return api.WriteJSON(w, http.StatusOK, out)
	}
}

// --- Jobs ---

func jobOut(j dominio.Job) api.Job {
	out := api.Job{ID: j.ID, Code: j.Code, Name: j.Name, Schedule: j.Schedule, IsEnabled: j.IsEnabled, CreatedAt: api.APITime(j.CreatedAt), UpdatedAt: api.APITime(j.UpdatedAt)}
	if j.NotificationID.Valid {
		out.NotificationID = api.Ptr(j.NotificationID.Int64)
	}
	if j.LastRunAt.Valid {
		out.LastRunAt = api.Ptr(api.APITime(j.LastRunAt.String))
	}
	if j.LastStatus.Valid {
		out.LastStatus = api.Ptr(j.LastStatus.String)
	}
	if j.LastError.Valid {
		out.LastError = api.Ptr(j.LastError.String)
	}
	if j.NextRunAt.Valid {
		out.NextRunAt = api.Ptr(api.APITime(j.NextRunAt.String))
	}
	return out
}

func listJobs(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"code", "name", "next_run_at", "last_run_at"}, "code")
		if err != nil {
			return err
		}
		rows, err := reg.Deps().DB.QueryContext(r.Context(), `SELECT id FROM job `+
			p.OrderBy(map[string]string{"code": "code", "name": "name", "next_run_at": "next_run_at", "last_run_at": "last_run_at", "id": "id"})+" "+p.LimitOffset())
		if err != nil {
			return err
		}
		var ids []int64
		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return err
			}
			ids = append(ids, id)
		}
		rows.Close()
		items := []api.Job{}
		for _, id := range ids {
			j, err := dominio.GetJob(r.Context(), reg.Deps().DB, id)
			if err != nil {
				return err
			}
			items = append(items, jobOut(j))
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

func updateJob(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "job_id")
		if err != nil {
			return err
		}
		var in api.JobUpdate
		if err := reg.DecodeJSON(r, &in); err != nil {
			return err
		}
		if in.Schedule == nil && in.IsEnabled == nil {
			return api.BadRequest("No hay cambios que guardar: manda la hora o si el job está activo.")
		}
		ctx, d := r.Context(), reg.Deps()
		now := d.Clock.Now()
		var next *string
		if in.Schedule != nil {
			s, err := jobs.ParseSchedule(*in.Schedule)
			if err != nil {
				return dominio.Invalid(dominio.FieldIssue{Field: "schedule", Code: "invalid_format",
					Message: "La hora no es válida: usa 5 campos (minuto hora día mes día-de-la-semana), por ejemplo «30 21 * * *»."})
			}
			n := dominio.FormatDBTime(s.Next(now))
			next = &n
		}
		err = dominio.WithTx(ctx, d.DB, func(tx *sql.Tx) error {
			before, err := dominio.RowJSON(ctx, tx, "job", id)
			if err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `UPDATE job SET schedule = coalesce(?, schedule), is_enabled = coalesce(?, is_enabled), next_run_at = coalesce(?, next_run_at), updated_at = ? WHERE id = ?`,
				in.Schedule, in.IsEnabled, next, dominio.FormatDBTime(now), id); err != nil {
				return err
			}
			after, err := dominio.RowJSON(ctx, tx, "job", id)
			if err != nil {
				return err
			}
			return dominio.Audit(ctx, tx, api.Actor(r), "update", "job", id, before, after, api.RequestID(r))
		})
		if errors.Is(err, sql.ErrNoRows) {
			return dominio.NotFound("Ese job no existe.")
		}
		if err != nil {
			return err
		}
		j, err := dominio.GetJob(ctx, d.DB, id)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusOK, jobOut(j))
	}
}
