package profeco

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// execer es lo que comparten *sql.DB y *sql.Tx.
type execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// DBPrices escribe en la tabla price del esquema con source = 'profeco'. Nunca
// es promocional. INSERT OR IGNORE sobre la restricción única (presentación,
// tienda, día, fuente, promoción): un precio repetido se cuenta como duplicado.
type DBPrices struct {
	X execer
}

// AddPrice implementa PriceSink.
func (d DBPrices) AddPrice(ctx context.Context, p Price) (bool, error) {
	res, err := d.X.ExecContext(ctx,
		`INSERT OR IGNORE INTO price (presentation_id, store_id, price_cents, observed_date, source, is_promotional, source_ref)
		 VALUES (?, ?, ?, ?, 'profeco', 0, ?)`,
		p.PresentationID, p.StoreID, p.PriceCents, p.ObservedDate, p.SourceRef)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}

// LoadStores lee las cadenas y sucursales activas con nombre PROFECO.
func LoadStores(ctx context.Context, db execer) (StoreSet, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, kind, coalesce(parent_id, 0), name, coalesce(profeco_name, ''), coalesce(state, ''), coalesce(municipality, '')
		FROM store WHERE is_active = 1 AND kind IN ('chain', 'branch') AND profeco_name IS NOT NULL ORDER BY id`)
	if err != nil {
		return StoreSet{}, err
	}
	defer rows.Close()
	var s StoreSet
	for rows.Next() {
		var id, parent int64
		var kind, name, pname, state, mun string
		if err := rows.Scan(&id, &kind, &parent, &name, &pname, &state, &mun); err != nil {
			return StoreSet{}, err
		}
		if kind == "chain" {
			s.Chains = append(s.Chains, Chain{ID: id, Name: name, ProfecoName: pname})
		} else {
			s.Branches = append(s.Branches, Branch{ID: id, ChainID: parent, Name: name, ProfecoName: pname, State: state, Municipality: mun})
		}
	}
	return s, rows.Err()
}

// MaxBytesFromSettings lee price.profeco_import_max_gb de app_setting (K6, K44)
// y lo devuelve en bytes. Nunca es una constante en código.
func MaxBytesFromSettings(ctx context.Context, db execer) (int64, error) {
	var v string
	err := db.QueryRowContext(ctx, `SELECT value_json FROM app_setting WHERE setting_key = 'price.profeco_import_max_gb' AND user_account_id IS NULL`).Scan(&v)
	if err != nil {
		return 0, fmt.Errorf("profeco: falta price.profeco_import_max_gb: %w", err)
	}
	gb, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil || gb <= 0 {
		return 0, fmt.Errorf("profeco: price.profeco_import_max_gb inválido: %q", v)
	}
	return gb << 30, nil
}

// LoadIntoDB importa los archivos en una sola transacción sobre la base de la
// app: o entra todo o no entra nada. Lee tiendas y límite de tamaño de la propia
// base; las ligas, el modo y el corte los pone el llamador en o.
func LoadIntoDB(ctx context.Context, db *sql.DB, o Options, paths ...string) (*ImportResult, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if o.Stores.Chains == nil && o.Stores.Branches == nil {
		if o.Stores, err = LoadStores(ctx, tx); err != nil {
			return nil, err
		}
	}
	if o.MaxBytes == 0 {
		if o.MaxBytes, err = MaxBytesFromSettings(ctx, tx); err != nil {
			return nil, err
		}
	}
	o.Prices = DBPrices{X: tx}
	im, err := NewImporter(o)
	if err != nil {
		return nil, err
	}
	total, _, err := im.Import(ctx, paths...)
	if err != nil {
		return total, err
	}
	return total, tx.Commit()
}

// QQP es la base local aparte datos/profeco/qqp.sqlite con todos los registros
// bien formados (FIRME, puerta 2, Q2-53 = A).
// Va fuera del esquema de la app. Escribe por lotes dentro de una transacción.
type QQP struct {
	db    *sql.DB
	tx    *sql.Tx
	stmt  *sql.Stmt
	batch int
	Rows  int64
}

const qqpBatch = 50000

const qqpSchema = `CREATE TABLE IF NOT EXISTS profeco_record (
	id INTEGER PRIMARY KEY,
	archivo TEXT NOT NULL,
	renglon INTEGER NOT NULL,
	producto TEXT NOT NULL, presentacion TEXT NOT NULL, marca TEXT NOT NULL,
	categoria TEXT NOT NULL, catalogo TEXT NOT NULL,
	precio TEXT NOT NULL, fecha_registro TEXT NOT NULL,
	cadena_comercial TEXT NOT NULL, giro TEXT NOT NULL, nombre_comercial TEXT NOT NULL,
	direccion TEXT NOT NULL, estado TEXT NOT NULL, municipio TEXT NOT NULL,
	latitud TEXT NOT NULL, longitud TEXT NOT NULL,
	precio_cents INTEGER NOT NULL, fecha TEXT NOT NULL,
	UNIQUE (archivo, renglon)
) STRICT;
CREATE INDEX IF NOT EXISTS ix_profeco_cadena_mun ON profeco_record (cadena_comercial, municipio, fecha)`

// NewQQP prepara la tabla en db (una base aparte, nunca la de la app).
func NewQQP(ctx context.Context, db *sql.DB) (*QQP, error) {
	if _, err := db.ExecContext(ctx, qqpSchema); err != nil {
		return nil, err
	}
	return &QQP{db: db}, nil
}

func (q *QQP) begin(ctx context.Context) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	st, err := tx.PrepareContext(ctx, `INSERT INTO profeco_record
		(archivo, renglon, producto, presentacion, marca, categoria, catalogo, precio, fecha_registro, cadena_comercial, giro, nombre_comercial, direccion, estado, municipio, latitud, longitud, precio_cents, fecha)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	q.tx, q.stmt, q.batch = tx, st, 0
	return nil
}

// AddRecord implementa RecordSink.
func (q *QQP) AddRecord(ctx context.Context, file string, r *Record) error {
	if q.tx == nil {
		if err := q.begin(ctx); err != nil {
			return err
		}
	}
	if _, err := q.stmt.ExecContext(ctx, file, r.Row, r.Producto, r.Presentacion, r.Marca, r.Categoria, r.Catalogo,
		r.PrecioTexto, r.FechaTexto, r.CadenaComercial, r.Giro, r.NombreComercial, r.Direccion, r.Estado, r.Municipio,
		r.Latitud, r.Longitud, r.PrecioCents, r.Fecha); err != nil {
		return err
	}
	q.Rows++
	q.batch++
	if q.batch >= qqpBatch {
		return q.Flush()
	}
	return nil
}

// Flush confirma el lote en curso.
func (q *QQP) Flush() error {
	if q.tx == nil {
		return nil
	}
	err := q.stmt.Close()
	if e := q.tx.Commit(); err == nil {
		err = e
	}
	q.tx, q.stmt = nil, nil
	return err
}
