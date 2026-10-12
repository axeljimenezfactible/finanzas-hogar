package gastos

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"finanzas-hogar/internal/adjuntos"
	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/jobs"
	"finanzas-hogar/internal/precios/ocr"
)

// ocrEngine lee las fotos con el OCR local (Tesseract, spa; motor S5). Las pruebas lo cambian por uno falso.
var ocrEngine ocr.Recognizer = ocr.New()

// ticketJob es el job que vacía la cola de tickets por leer. Cada 5 minutos es PROVISIONAL (bloqueo §17);
// createTicketScan además lo pide «ahora» para que la lectura empiece en el siguiente tick.
const ticketJob = "ticket_ocr"

func registerTicketJob(reg *api.Registry) {
	if reg.Deps().Jobs == nil {
		return
	}
	if err := reg.Deps().Jobs.Register(jobs.Def{Code: ticketJob, Name: "Lectura de tickets", Schedule: "*/5 * * * *", Run: runTicketOCR(reg)}); err != nil {
		panic(err)
	}
}

func uploadName(p *multipart.Part) string {
	if _, params, err := mime.ParseMediaType(p.Header.Get("Content-Disposition")); err == nil {
		return params["filename"]
	}
	return p.FileName()
}

func ticketTooLarge(maxMB int64) error {
	return &api.Problem{Status: http.StatusRequestEntityTooLarge, Code: "payload_too_large",
		Message: "La foto pesa más de " + strconv.FormatInt(maxMB, 10) + " MB. Toma otra con menos resolución."}
}

// createTicketScan (N-075, K3): guarda la foto y la deja en la cola del OCR; responde 202 con el ticket pendiente.
func createTicketScan(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx, d := r.Context(), reg.Deps()
		maxMB, err := dominio.SettingInt(ctx, d.DB, "attachment.max_size_mb")
		if err != nil {
			return err
		}
		maxBytes := maxMB * adjuntos.BytesPerMB
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes+adjuntos.BytesPerMB)
		mr, err := r.MultipartReader()
		if err != nil {
			return &api.Problem{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "Manda la foto como formulario multipart/form-data."}
		}
		st := adjuntos.Store{DataDir: d.Config.DataDir}
		var saved *adjuntos.Saved
		var original string
		fail := func(e error) error {
			if saved != nil {
				_ = st.Remove(saved.RelPath)
			}
			return e
		}
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			var mbe *http.MaxBytesError
			if errors.As(err, &mbe) {
				return fail(ticketTooLarge(maxMB))
			}
			if err != nil {
				return fail(api.BadRequest("El formulario de la foto no es válido."))
			}
			if part.FormName() != "file" {
				return fail(api.BadRequest("El formulario trae un campo que no se conoce: " + part.FormName() + "."))
			}
			if saved != nil {
				return fail(api.BadRequest("Manda una sola foto por petición."))
			}
			if original = uploadName(part); original != "" && adjuntos.CheckName(original) != nil {
				return dominio.Invalid(dominio.FieldIssue{Field: "file", Code: "invalid_format", Message: "El nombre del archivo no es válido: no puede llevar «/», «\\» ni «..»."})
			}
			s, err := st.Save(part, maxBytes, d.Clock.Now())
			switch {
			case errors.Is(err, adjuntos.ErrTooLarge) || errors.As(err, &mbe):
				return ticketTooLarge(maxMB)
			case errors.Is(err, adjuntos.ErrType):
				return &api.Problem{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "Sube la foto del ticket en JPG, PNG o WebP."}
			case err != nil:
				return err
			}
			saved = &s
			if !strings.HasPrefix(s.MimeType, "image/") { // C1: el ticket solo admite fotos (un PDF va como adjunto del gasto)
				return fail(&api.Problem{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "Sube la foto del ticket en JPG, PNG o WebP; un PDF se adjunta al gasto."})
			}
		}
		var v api.Validator
		v.Required(saved != nil, "file")
		if err := v.Err(); err != nil {
			return err
		}
		var id int64
		err = dominio.WithTx(ctx, d.DB, func(tx *sql.Tx) error {
			var orig *string
			if original != "" {
				orig = &original
			}
			if id, err = dominio.GastosInsertTicket(ctx, tx, saved.RelPath, orig, saved.MimeType, saved.Size, saved.SHA256, api.Actor(r)); err != nil {
				return err
			}
			after, err := dominio.RowJSON(ctx, tx, "attachment", id)
			if err != nil {
				return err
			}
			return dominio.Audit(ctx, tx, api.Actor(r), "insert", "attachment", id, nil, after, api.RequestID(r))
		})
		if err != nil {
			return fail(err)
		}
		if d.Jobs != nil {
			if err := d.Jobs.RunNow(ctx, ticketJob); err != nil { // sin el renglón del job, la lectura espera al siguiente arranque
				d.Log.Warn("ticket_ocr: no se pudo pedir la lectura ahora", "attachment_id", id, "error", err.Error())
			}
		}
		t, err := dominio.GetTicketScan(ctx, d.DB, id)
		if err != nil {
			return err
		}
		ts, err := toTicketScan(t)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusAccepted, ts)
	}
}

// runTicketOCR lee cada foto pendiente con el motor de OCR (solo la llamada) y guarda el texto crudo: done o failed.
func runTicketOCR(reg *api.Registry) jobs.Func {
	return func(ctx context.Context, run jobs.Run) error {
		d := reg.Deps()
		pend, err := dominio.GastosPendingTickets(ctx, d.DB)
		if err != nil {
			return err
		}
		st := adjuntos.Store{DataDir: d.Config.DataDir}
		for _, p := range pend {
			status, engine, text := "failed", (*string)(nil), (*string)(nil)
			abs, err := st.Resolve(p.FilePath)
			if err == nil {
				var lines []ocr.Line
				if lines, err = ocrEngine.Recognize(ctx, abs); err == nil {
					raw := make([]string, 0, len(lines))
					for _, l := range lines {
						raw = append(raw, l.Texto)
					}
					status, engine, text = "done", api.Ptr("tesseract "+ocr.RequiredVersion+" spa"), api.Ptr(strings.Join(raw, "\n"))
				}
			}
			if err != nil {
				d.Log.Warn("ticket_ocr: lectura fallida", "attachment_id", p.ID, "error", err.Error())
			}
			err = dominio.WithTx(ctx, d.DB, func(tx *sql.Tx) error {
				before, err := dominio.RowJSON(ctx, tx, "attachment", p.ID)
				if err != nil {
					return err
				}
				if err := dominio.GastosSetTicketOCR(ctx, tx, p.ID, status, engine, text, dominio.FormatDBTime(run.Now)); err != nil {
					return err
				}
				after, err := dominio.RowJSON(ctx, tx, "attachment", p.ID)
				if err != nil {
					return err
				}
				return dominio.Audit(ctx, tx, nil, "update", "attachment", p.ID, before, after, run.RequestID)
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
}
