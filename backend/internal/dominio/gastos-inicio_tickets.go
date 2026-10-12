package dominio

import (
	"context"
	"database/sql"
	"errors"
)

// Repositorio de lecturas de tickets (attachment con kind = ticket_photo) del agente api-gastos-inicio (K3, RN-23).

// TicketScanRow es una lectura de ticket guardada.
type TicketScanRow struct {
	ID                      int64
	ExpenseID               sql.NullInt64
	Kind                    string
	OriginalName            sql.NullString
	MimeType                string
	SizeBytes               int64
	SHA256, OCRStatus       string
	OCREngine, OCRText      sql.NullString
	UploadedByUserID        sql.NullInt64
	CreatedAt, UpdatedAt    string
	DraftJSON, DraftSavedAt sql.NullString
}

const ticketCols = `a.id, a.expense_id, a.kind, a.original_name, a.mime_type, a.size_bytes, a.sha256, a.ocr_status, a.ocr_engine, a.ocr_text,
 a.uploaded_by_user_id, a.created_at, a.updated_at, a.ocr_draft_json, a.ocr_draft_saved_at`

func scanTicket(sc interface{ Scan(...any) error }) (TicketScanRow, error) {
	var t TicketScanRow
	err := sc.Scan(&t.ID, &t.ExpenseID, &t.Kind, &t.OriginalName, &t.MimeType, &t.SizeBytes, &t.SHA256, &t.OCRStatus, &t.OCREngine, &t.OCRText,
		&t.UploadedByUserID, &t.CreatedAt, &t.UpdatedAt, &t.DraftJSON, &t.DraftSavedAt)
	return t, err
}

// ListTicketScans lee la cola de lecturas sin confirmar (pending, done o failed).
func ListTicketScans(ctx context.Context, q Querier, orderBy, limitOffset string) ([]TicketScanRow, error) {
	rows, err := q.QueryContext(ctx, `SELECT `+ticketCols+` FROM attachment a WHERE a.kind = 'ticket_photo' AND a.expense_id IS NULL
 AND a.ocr_status IN ('pending', 'done', 'failed') `+orderBy+" "+limitOffset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []TicketScanRow
	for rows.Next() {
		t, err := scanTicket(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// GetTicketScan lee una lectura de ticket (404 si no existe o no es foto de ticket).
func GetTicketScan(ctx context.Context, q Querier, id int64) (TicketScanRow, error) {
	t, err := scanTicket(q.QueryRowContext(ctx, `SELECT `+ticketCols+` FROM attachment a WHERE a.id = ? AND a.kind = 'ticket_photo'`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return t, NotFound("No encontramos esa lectura de ticket.")
	}
	return t, err
}

// SaveTicketScanDraft guarda el borrador de «Guardar para después» (K3); una lectura confirmada responde 409.
func SaveTicketScanDraft(ctx context.Context, q Querier, id int64, draftJSON, savedAt string) error {
	t, err := GetTicketScan(ctx, q, id)
	if err != nil {
		return err
	}
	if t.OCRStatus == "confirmed" || t.ExpenseID.Valid {
		return Conflict("conflict", "Este ticket ya está confirmado; sus cambios se hacen en el gasto.")
	}
	_, err = q.ExecContext(ctx, `UPDATE attachment SET ocr_draft_json = ?, ocr_draft_saved_at = ? WHERE id = ?`, draftJSON, savedAt, id)
	return err
}

// PendingTicket es una foto de ticket en la cola del OCR (ocr_status = pending, sin gasto).
type PendingTicket struct {
	ID       int64
	FilePath string
}

// GastosInsertTicket guarda la foto subida como ticket por leer (createTicketScan): sin gasto y con OCR pendiente.
func GastosInsertTicket(ctx context.Context, q Querier, filePath string, originalName *string, mimeType string, size int64, sha string, userID *int64) (int64, error) {
	var id int64
	err := q.QueryRowContext(ctx, `INSERT INTO attachment (expense_id, kind, file_path, original_name, mime_type, size_bytes, sha256, ocr_status, uploaded_by_user_id)
 VALUES (NULL, 'ticket_photo', ?, ?, ?, ?, ?, 'pending', ?) RETURNING id`, filePath, originalName, mimeType, size, sha, userID).Scan(&id)
	return id, err
}

// GastosPendingTickets lista la cola del job ticket_ocr en orden de llegada.
func GastosPendingTickets(ctx context.Context, q Querier) ([]PendingTicket, error) {
	rows, err := q.QueryContext(ctx, `SELECT id, file_path FROM attachment WHERE kind = 'ticket_photo' AND expense_id IS NULL AND ocr_status = 'pending' ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []PendingTicket{}
	for rows.Next() {
		var p PendingTicket
		if err := rows.Scan(&p.ID, &p.FilePath); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// GastosSetTicketOCR guarda el resultado de la lectura: done con motor y texto crudo, o failed sin texto.
func GastosSetTicketOCR(ctx context.Context, q Querier, id int64, status string, engine, text *string, now string) error {
	_, err := q.ExecContext(ctx, `UPDATE attachment SET ocr_status = ?, ocr_engine = ?, ocr_text = ?, updated_at = ? WHERE id = ? AND ocr_status = 'pending'`,
		status, engine, text, now, id)
	return err
}

// ---------- confirmación del ticket (confirmTicketScan, N-074) ----------

// TicketPresentation es lo que la confirmación necesita de una presentación del catálogo.
type TicketPresentation struct {
	ID, ProductID, SizeMilli int64
	ProductName              string
}

// GastosTicketPresentations lee las presentaciones de los renglones del ticket; las que no existen no vienen.
func GastosTicketPresentations(ctx context.Context, q Querier, ids []int64) (map[int64]TicketPresentation, error) {
	out := map[int64]TicketPresentation{}
	if len(ids) == 0 {
		return out, nil
	}
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := q.QueryContext(ctx, `SELECT pr.id, pr.product_id, pr.size_milli, pd.name FROM presentation pr JOIN product pd ON pd.id = pr.product_id
 WHERE `+gxInList("pr.id", len(ids)), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p TicketPresentation
		if err := rows.Scan(&p.ID, &p.ProductID, &p.SizeMilli, &p.ProductName); err != nil {
			return nil, err
		}
		out[p.ID] = p
	}
	return out, rows.Err()
}

// GastosExpenseItemIDs devuelve el id de cada renglón de un gasto por su número de renglón.
func GastosExpenseItemIDs(ctx context.Context, q Querier, expenseID int64) (map[int64]int64, error) {
	rows, err := q.QueryContext(ctx, `SELECT line_no, id FROM expense_item WHERE expense_id = ?`, expenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64]int64{}
	for rows.Next() {
		var line, id int64
		if err := rows.Scan(&line, &id); err != nil {
			return nil, err
		}
		out[line] = id
	}
	return out, rows.Err()
}

// GastosStockItem busca el renglón de inventario por unidad (tracking_mode = stock, RN-15) de un producto.
func GastosStockItem(ctx context.Context, q Querier, productID int64) (int64, bool, error) {
	var id int64
	err := q.QueryRowContext(ctx, `SELECT id FROM inventory_item WHERE product_id = ? AND tracking_mode = 'stock'`, productID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	return id, err == nil, err
}

// gxAuditWrite escribe audit_log con la foto «después» del renglón.
func gxAuditWrite(ctx context.Context, q Querier, actor *int64, action, table string, id int64, before []byte, requestID string) error {
	after, err := RowJSON(ctx, q, table, id)
	if err != nil {
		return err
	}
	return Audit(ctx, q, actor, action, table, id, before, after, requestID)
}

// GastosTicketStockIn suma la compra de un renglón del ticket a la existencia y deja su inventory_movement purchase
// ligado al renglón, en la misma transacción (RN-15, invariante 12 de S1), con audit_log de los dos.
func GastosTicketStockIn(ctx context.Context, tx *sql.Tx, itemID int64, date string, delta, expenseItemID int64, now string, actor *int64, requestID string) error {
	before, err := RowJSON(ctx, tx, "inventory_item", itemID)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE inventory_item SET quantity_milli = quantity_milli + ?, updated_at = ? WHERE id = ?`, delta, now, itemID); err != nil {
		return err
	}
	if err := gxAuditWrite(ctx, tx, actor, "update", "inventory_item", itemID, before, requestID); err != nil {
		return err
	}
	var mid int64
	if err := tx.QueryRowContext(ctx, `INSERT INTO inventory_movement (inventory_item_id, movement_date, kind, quantity_delta_milli, expense_item_id)
 VALUES (?, ?, 'purchase', ?, ?) RETURNING id`, itemID, date, delta, expenseItemID).Scan(&mid); err != nil {
		return err
	}
	return gxAuditWrite(ctx, tx, actor, "insert", "inventory_movement", mid, nil, requestID)
}

// TicketPrice es el precio de una presentación leído en un ticket confirmado (fuente ticket_ocr).
type TicketPrice struct {
	PresentationID, StoreID, PriceCents, ExpenseItemID, AttachmentID int64
	ObservedDate                                                     string
}

// GastosSaveTicketPrice guarda el precio del ticket. Si ese día ya hay un precio ticket_ocr regular de esa presentación
// en esa tienda (otro ticket confirmado), la llave UNIQUE de price no deja otro: se actualiza con el de este ticket
// (PROVISIONAL, bloqueo S6-gastos-inicio §18).
func GastosSaveTicketPrice(ctx context.Context, tx *sql.Tx, p TicketPrice, now string, actor *int64, requestID string) (int64, error) {
	var id int64
	err := tx.QueryRowContext(ctx, `SELECT id FROM price WHERE presentation_id = ? AND store_id = ? AND observed_date = ? AND source = 'ticket_ocr'
 AND is_promotional = 0`, p.PresentationID, p.StoreID, p.ObservedDate).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		if err := tx.QueryRowContext(ctx, `INSERT INTO price (presentation_id, store_id, price_cents, observed_date, source, is_promotional, expense_item_id, attachment_id)
 VALUES (?, ?, ?, ?, 'ticket_ocr', 0, ?, ?) RETURNING id`, p.PresentationID, p.StoreID, p.PriceCents, p.ObservedDate, p.ExpenseItemID, p.AttachmentID).Scan(&id); err != nil {
			return 0, err
		}
		return id, gxAuditWrite(ctx, tx, actor, "insert", "price", id, nil, requestID)
	} else if err != nil {
		return 0, err
	}
	before, err := RowJSON(ctx, tx, "price", id)
	if err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE price SET price_cents = ?, expense_item_id = ?, attachment_id = ?, updated_at = ? WHERE id = ?`,
		p.PriceCents, p.ExpenseItemID, p.AttachmentID, now, id); err != nil {
		return 0, err
	}
	return id, gxAuditWrite(ctx, tx, actor, "update", "price", id, before, requestID)
}

// GastosConfirmTicket liga el ticket a su gasto, lo marca confirmed y vacía el borrador (RN-23). false = ya estaba
// confirmado (otra confirmación ganó la carrera).
func GastosConfirmTicket(ctx context.Context, q Querier, id, expenseID int64, now string) (bool, error) {
	res, err := q.ExecContext(ctx, `UPDATE attachment SET expense_id = ?, ocr_status = 'confirmed', ocr_draft_json = NULL, ocr_draft_saved_at = NULL,
 updated_at = ? WHERE id = ? AND kind = 'ticket_photo' AND expense_id IS NULL AND ocr_status <> 'confirmed'`, expenseID, now, id)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n == 1, err
}
