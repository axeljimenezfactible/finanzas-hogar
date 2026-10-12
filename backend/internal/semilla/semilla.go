// Package semilla carga de forma idempotente los datos de arranque de datos/semilla/ a una base migrada
// (WF-3a, agente semilla): personas y cuentas, reglas de reparto, categorías, tiendas, catálogo de productos,
// tabla nutricional confirmada, recetario, la compra del 10-oct-2026 de data.json, el crédito Kueski con sus
// despensas (todas las cifras salen del motor motor_cred y se cotejan con oraculo/model.json) y la muestra de PROFECO.
//
// Reglas:
//   - Idempotente: cada entidad se busca por su llave natural (código, nombre, fecha e importe) y solo se crea si
//     falta. Sembrar dos veces no duplica ni cambia lo ya cargado (tampoco las contraseñas).
//   - Todo ocurre en una sola transacción: si algo falla (o una cifra del motor difiere del oráculo) no queda nada.
//   - Dinero en centavos enteros; nada de cifras copiadas a mano de las despensas.
//   - Cada alta escribe audit_log (actor nulo = proceso del sistema, request_id «semilla»).
//   - Las contraseñas iniciales NUNCA están en el repositorio: vienen de variables de entorno o de un archivo
//     fuera del repositorio (ver PasswordsFromEnv). Nada se imprime.
package semilla

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/recetas"
)

// RequestID es el request_id con que la semilla firma audit_log.
const RequestID = "semilla"

// Options son las entradas de una carga.
type Options struct {
	// Dir es datos/semilla (data.json, productos.json, recetario.json, nutricion/, profeco/).
	Dir string
	// ModelPath es oraculo/model.json; vacío = Dir/model.json y, si no existe, Dir/../../oraculo/model.json.
	ModelPath string
	// Passwords son las contraseñas iniciales por usuario (axel, ximena). Solo se usan para cuentas nuevas.
	Passwords map[string]string
	// Argon2 reemplaza los parámetros del hash (pruebas); nil = dominio.DefaultArgon2.
	Argon2 *dominio.Argon2Params
	// Clock da «hoy» (rechazo de fechas futuras de PROFECO); nil = reloj del sistema.
	Clock dominio.Clock
	// ImportPlannedCSV importa gastos_fijos_ingresos.csv cuando el archivo existe. Lo instala quien tenga el
	// importador de Ajustes (PROVISIONAL: solicitud fuera de alcance); nil = solo avisa.
	ImportPlannedCSV func(ctx context.Context, db dominio.Querier, path string) (int, error)
}

// Entry cuenta lo creado y lo que ya existía de una entidad.
type Entry struct {
	Name              string
	Created, Existing int
}

// Report es el resultado de una carga.
type Report struct {
	Entries []Entry
	Notices []string // avisos para Axel (archivos que faltan, cosas PROVISIONALES)
	Checks  []string // cotejos contra el oráculo que pasaron
	index   map[string]int
}

func (r *Report) count(name string, created bool) {
	if r.index == nil {
		r.index = map[string]int{}
	}
	i, ok := r.index[name]
	if !ok {
		i = len(r.Entries)
		r.index[name] = i
		r.Entries = append(r.Entries, Entry{Name: name})
	}
	if created {
		r.Entries[i].Created++
	} else {
		r.Entries[i].Existing++
	}
}

// Get devuelve lo creado de una entidad (para pruebas y para el reporte).
func (r *Report) Get(name string) Entry {
	if i, ok := r.index[name]; ok {
		return r.Entries[i]
	}
	return Entry{Name: name}
}

// TotalCreated suma todo lo creado.
func (r *Report) TotalCreated() int {
	n := 0
	for _, e := range r.Entries {
		n += e.Created
	}
	return n
}

// String es el reporte en español de México, sin ningún secreto.
func (r *Report) String() string {
	var b strings.Builder
	b.WriteString("Semilla de Finanzas Hogar\n")
	for _, e := range r.Entries {
		fmt.Fprintf(&b, "  %-34s nuevos %5d · ya estaban %5d\n", e.Name, e.Created, e.Existing)
	}
	for _, c := range r.Checks {
		fmt.Fprintf(&b, "  cotejo OK: %s\n", c)
	}
	for _, n := range r.Notices {
		fmt.Fprintf(&b, "  AVISO: %s\n", n)
	}
	return b.String()
}

type loader struct {
	ctx  context.Context
	tx   *sql.Tx
	opt  Options
	rep  *Report
	ids  map[string]int64 // llaves de trabajo entre pasos (person:axel, rule:hogar…)
	argo dominio.Argon2Params
	cat  *recetas.Catalog   // catálogo validado de productos.json
	nut  *recetas.Nutrition // tabla nutricional validada (modo revisión: incluye «propuesto»)
	pres map[string]int64   // nombre exacto de la presentación y nombres de ticket → id
}

func (l *loader) id(key string) int64 { return l.ids[key] }

// find devuelve el id del primer renglón que devuelve la consulta.
func (l *loader) find(q string, args ...any) (int64, bool, error) {
	var id int64
	err := l.tx.QueryRowContext(l.ctx, q, args...).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	return id, err == nil, err
}

// insert inserta un renglón y escribe su audit_log (action insert).
func (l *loader) insert(table, q string, args ...any) (int64, error) {
	res, err := l.tx.ExecContext(l.ctx, q, args...)
	if err != nil {
		return 0, fmt.Errorf("semilla: insertar en %s: %w", table, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	after, err := dominio.RowJSON(l.ctx, l.tx, table, id)
	if err != nil {
		return 0, err
	}
	if err := dominio.Audit(l.ctx, l.tx, nil, "insert", table, id, nil, after, RequestID); err != nil {
		return 0, err
	}
	return id, nil
}

// ensure busca por la llave natural y, si falta, inserta. Cuenta lo creado y lo existente bajo la entidad `name`.
func (l *loader) ensure(name, table, findQ string, findArgs []any, insQ string, insArgs ...any) (int64, bool, error) {
	if id, ok, err := l.find(findQ, findArgs...); err != nil {
		return 0, false, err
	} else if ok {
		l.rep.count(name, false)
		return id, false, nil
	}
	id, err := l.insert(table, insQ, insArgs...)
	if err != nil {
		return 0, false, err
	}
	l.rep.count(name, true)
	return id, true, nil
}

func a(args ...any) []any { return args }

// DefaultImportPlannedCSV es el importador de gastos_fijos_ingresos.csv que usa Run. Queda en nil hasta que Ajustes
// exporte su importador (solicitud fuera de alcance en docs/estado/wf3/semilla.md); entonces cmd/finanzas lo instala
// en una línea. Mientras sea nil, si el archivo existe la carga solo avisa.
var DefaultImportPlannedCSV func(ctx context.Context, db dominio.Querier, path string) (int, error)

// Run carga la semilla en la base (una transacción). Es la firma que usa «finanzas semilla».
func Run(ctx context.Context, db *sql.DB, dir string) error {
	opt := Options{Dir: dir, ImportPlannedCSV: DefaultImportPlannedCSV}
	pw, err := PasswordsFromEnv()
	if err != nil {
		return err
	}
	opt.Passwords = pw
	rep, err := RunWith(ctx, db, opt)
	if rep != nil {
		fmt.Print(rep.String())
	}
	return err
}

// RunWith carga la semilla con opciones explícitas y devuelve el reporte.
func RunWith(ctx context.Context, db *sql.DB, opt Options) (*Report, error) {
	if opt.Dir == "" {
		return nil, errors.New("semilla: falta la carpeta de la semilla (datos/semilla)")
	}
	if _, err := os.Stat(filepath.Join(opt.Dir, "data.json")); err != nil {
		return nil, fmt.Errorf("semilla: no encuentro data.json en %s: %w", opt.Dir, err)
	}
	rep := &Report{}
	err := dominio.WithTx(ctx, db, func(tx *sql.Tx) error {
		l := &loader{ctx: ctx, tx: tx, opt: opt, rep: rep, ids: map[string]int64{}, pres: map[string]int64{}, argo: dominio.DefaultArgon2}
		if opt.Argon2 != nil {
			l.argo = *opt.Argon2
		}
		steps := []func() error{l.personasYCuentas, l.reglasYCategorias, l.tiendas, l.catalogo, l.nutricion, l.recetario,
			l.kueski, l.despensas, l.profeco, l.partidasCSV, l.cruce}
		for _, s := range steps {
			if err := s(); err != nil {
				return err
			}
		}
		return l.avisos()
	})
	if err != nil {
		return nil, err
	}
	return rep, nil
}
