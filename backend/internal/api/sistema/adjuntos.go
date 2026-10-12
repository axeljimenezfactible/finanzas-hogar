package sistema

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"

	"finanzas-hogar/internal/adjuntos"
	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/dominio"
)

// attachmentKinds: attachment.kind. Sin kind en la petición se guarda «other» (PROVISIONAL, bloqueo S6-adjuntos).
var attachmentKinds = map[string]bool{"ticket_photo": true, "invoice": true, "bank_statement": true, "other": true}

func registerAttachments(reg *api.Registry) {
	reg.Handle("uploadExpenseAttachment", uploadExpenseAttachment(reg))
	reg.Handle("downloadAttachment", downloadAttachment(reg))
	reg.Handle("deleteAttachment", deleteAttachment(reg))
}

func store(reg *api.Registry) adjuntos.Store {
	return adjuntos.Store{DataDir: reg.Deps().Config.DataDir}
}

func tooLarge(maxMB int64) error {
	return &api.Problem{Status: http.StatusRequestEntityTooLarge, Code: "payload_too_large",
		Message: fmt.Sprintf("El archivo pesa más de %d MB, el máximo vigente. Mándalo más ligero.", maxMB)}
}

const attachmentCols = `id, expense_id, kind, original_name, mime_type, size_bytes, sha256, ocr_status, ocr_engine, ocr_text, uploaded_by_user_id, created_at, updated_at`

func getAttachment(r *http.Request, q dominio.Querier, id int64) (api.Attachment, string, error) {
	var a api.Attachment
	var exp, by sql.NullInt64
	var orig, eng, txt sql.NullString
	var file string
	err := q.QueryRowContext(r.Context(), `SELECT `+attachmentCols+`, file_path FROM attachment WHERE id = ?`, id).
		Scan(&a.ID, &exp, &a.Kind, &orig, &a.MimeType, &a.SizeBytes, &a.SHA256, &a.OCRStatus, &eng, &txt, &by, &a.CreatedAt, &a.UpdatedAt, &file)
	if errors.Is(err, sql.ErrNoRows) {
		return a, "", dominio.NotFound("Ese adjunto no existe.")
	}
	a.CreatedAt, a.UpdatedAt = api.APITime(a.CreatedAt), api.APITime(a.UpdatedAt)
	if exp.Valid {
		a.ExpenseID = api.Ptr(exp.Int64)
	}
	if by.Valid {
		a.UploadedByUserID = api.Ptr(by.Int64)
	}
	if orig.Valid {
		a.OriginalName = api.Ptr(orig.String)
	}
	if eng.Valid {
		a.OCREngine = api.Ptr(eng.String)
	}
	if txt.Valid {
		a.OCRText = api.Ptr(txt.String)
	}
	return a, file, err
}

func uploadExpenseAttachment(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		expenseID, err := api.PathID(r, "expense_id")
		if err != nil {
			return err
		}
		ctx, d := r.Context(), reg.Deps()
		var one int
		if err := d.DB.QueryRowContext(ctx, `SELECT 1 FROM expense WHERE id = ?`, expenseID).Scan(&one); errors.Is(err, sql.ErrNoRows) {
			return dominio.NotFound("Ese gasto no existe.")
		} else if err != nil {
			return err
		}
		maxMB, err := dominio.SettingInt(ctx, d.DB, "attachment.max_size_mb")
		if err != nil {
			return err
		}
		maxBytes := maxMB * adjuntos.BytesPerMB
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes+adjuntos.BytesPerMB) // margen para las cabeceras del multipart
		mr, err := r.MultipartReader()
		if err != nil {
			return &api.Problem{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "Manda el archivo como formulario multipart/form-data."}
		}
		now := d.Clock.Now()
		st := store(reg)
		var saved *adjuntos.Saved
		var original, kind string
		cleanup := func() {
			if saved != nil {
				_ = st.Remove(saved.RelPath)
			}
		}
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			var mbe *http.MaxBytesError
			if errors.As(err, &mbe) {
				cleanup()
				return tooLarge(maxMB)
			}
			if err != nil {
				cleanup()
				return api.BadRequest("El formulario del archivo no es válido.")
			}
			switch part.FormName() {
			case "kind":
				b, _ := io.ReadAll(io.LimitReader(part, 64))
				kind = string(b)
			case "file":
				if saved != nil {
					cleanup()
					return api.BadRequest("Manda un solo archivo por petición.")
				}
				original = rawFileName(part)
				if original != "" && adjuntos.CheckName(original) != nil {
					return dominio.Invalid(dominio.FieldIssue{Field: "file", Code: "invalid_format", Message: "El nombre del archivo no es válido: no puede llevar «/», «\\» ni «..»."})
				}
				s, err := st.Save(part, maxBytes, now)
				switch {
				case errors.Is(err, adjuntos.ErrTooLarge) || errors.As(err, &mbe):
					return tooLarge(maxMB)
				case errors.Is(err, adjuntos.ErrType):
					return &api.Problem{Status: http.StatusUnsupportedMediaType, Code: "unsupported_media_type", Message: "Ese tipo de archivo no se admite: sube una foto JPG, PNG o WebP, o un PDF."}
				case err != nil:
					return err
				}
				saved = &s
			default:
				cleanup()
				return api.BadRequest("El formulario trae un campo que no se conoce: " + part.FormName() + ".")
			}
		}
		if kind == "" {
			kind = "other"
		}
		var v api.Validator
		v.Required(saved != nil, "file")
		v.Check(attachmentKinds[kind], "kind", "invalid_format", "El tipo de adjunto no es válido: ticket_photo, invoice, bank_statement u other.")
		if err := v.Err(); err != nil {
			cleanup()
			return err
		}
		var id int64
		err = dominio.WithTx(ctx, d.DB, func(tx *sql.Tx) error {
			var orig any
			if original != "" {
				orig = original
			}
			if err := tx.QueryRowContext(ctx, `INSERT INTO attachment (expense_id, kind, file_path, original_name, mime_type, size_bytes, sha256, uploaded_by_user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`,
				expenseID, kind, saved.RelPath, orig, saved.MimeType, saved.Size, saved.SHA256, api.Actor(r)).Scan(&id); err != nil {
				return err
			}
			after, err := dominio.RowJSON(ctx, tx, "attachment", id)
			if err != nil {
				return err
			}
			return dominio.Audit(ctx, tx, api.Actor(r), "insert", "attachment", id, nil, after, api.RequestID(r))
		})
		if err != nil {
			cleanup()
			return err
		}
		a, _, err := getAttachment(r, d.DB, id)
		if err != nil {
			return err
		}
		return api.WriteJSON(w, http.StatusCreated, a)
	}
}

func downloadAttachment(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "attachment_id")
		if err != nil {
			return err
		}
		_, rel, err := getAttachment(r, reg.Deps().DB, id)
		if err != nil {
			return err
		}
		full, err := store(reg).Resolve(rel)
		if err != nil {
			if errors.Is(err, adjuntos.ErrPath) {
				reg.Deps().Log.Warn("adjunto con ruta fuera de la carpeta", "attachment_id", id, "request_id", api.RequestID(r))
			}
			return dominio.NotFound("El archivo de ese adjunto no está en el servidor.")
		}
		f, err := os.Open(full)
		if err != nil {
			return dominio.NotFound("El archivo de ese adjunto no está en el servidor.")
		}
		defer f.Close()
		st, err := f.Stat()
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="adjunto-`+strconv.FormatInt(id, 10)+path.Ext(rel)+`"`)
		w.Header().Set("Content-Length", strconv.FormatInt(st.Size(), 10))
		w.Header().Set("Cache-Control", "private, no-store")
		w.WriteHeader(http.StatusOK)
		_, err = io.Copy(w, f)
		return err
	}
}

func deleteAttachment(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := api.PathID(r, "attachment_id")
		if err != nil {
			return err
		}
		ctx, d := r.Context(), reg.Deps()
		_, rel, err := getAttachment(r, d.DB, id)
		if err != nil {
			return err
		}
		// price.attachment_id y credit_statement.attachment_id quedan en NULL (ON DELETE SET NULL del esquema)
		err = dominio.WithTx(ctx, d.DB, func(tx *sql.Tx) error {
			before, err := dominio.RowJSON(ctx, tx, "attachment", id)
			if err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `DELETE FROM attachment WHERE id = ?`, id); err != nil {
				return err
			}
			return dominio.Audit(ctx, tx, api.Actor(r), "delete", "attachment", id, before, nil, api.RequestID(r))
		})
		if err != nil {
			return err
		}
		if err := store(reg).Remove(rel); err != nil {
			d.Log.Warn("no se pudo borrar el archivo del adjunto", "attachment_id", id, "request_id", api.RequestID(r), "error", err.Error())
		}
		return api.NoContent(w)
	}
}

// rawFileName es el nombre tal como lo mandó el cliente (Part.FileName ya lo recorta con filepath.Base y escondería
// un «../»): se valida para rechazarlo y solo se guarda como dato en original_name.
func rawFileName(p *multipart.Part) string {
	_, params, err := mime.ParseMediaType(p.Header.Get("Content-Disposition"))
	if err != nil {
		return p.FileName()
	}
	return params["filename"]
}
