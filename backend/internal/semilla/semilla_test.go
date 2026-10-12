package semilla_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/semilla"
	"finanzas-hogar/internal/testutil"
)

var testPasswords = map[string]string{"axel": testutil.FakeCredential, "ximena": testutil.FakeCredential + "-x"}

// emptyDB es una base migrada y vacía (sin las fixtures de testutil.NewDB) en datos/tmp.
func emptyDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := dominio.Open(filepath.Join(testutil.TempDir(t), "semilla.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err := dominio.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	return db
}

// oracleSeed es oraculo/fixtures/s2/_datos/semilla-oraculo.json (para las cargas desde una carpeta copiada a datos/tmp).
func oracleSeed(t *testing.T) string {
	return filepath.Join(testutil.Root(t), "oraculo", "fixtures", "s2", "_datos", "semilla-oraculo.json")
}

func seedDir(t *testing.T) string { return filepath.Join(testutil.Root(t), "datos", "semilla") }

func opts(t *testing.T) semilla.Options {
	return semilla.Options{Dir: seedDir(t), Passwords: testPasswords, Argon2: &testutil.FastArgon2}
}

func one[T any](t *testing.T, db *sql.DB, q string, args ...any) T {
	t.Helper()
	var v T
	if err := db.QueryRow(q, args...).Scan(&v); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
	return v
}

// miniSeed copia solo los archivos pequeños de la semilla a datos/tmp (profeco/db pesa 16 GB y no se copia).
func miniSeed(t *testing.T) string {
	t.Helper()
	dst := testutil.TempDir(t)
	src := seedDir(t)
	for _, f := range []string{"data.json", "productos.json", "recetario.json", "nutricion/tabla.csv", "profeco/muestra/qqp-muestra.csv"} {
		b, err := os.ReadFile(filepath.Join(src, f))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Dir(filepath.Join(dst, f)), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dst, f), b, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dst
}

func TestSembrarCargaTodoYCruzaConElOraculo(t *testing.T) {
	db := emptyDB(t)
	rep, err := semilla.RunWith(context.Background(), db, opts(t))
	if err != nil {
		t.Fatal(err)
	}
	for table, want := range map[string]int{"person": 3, "user_account": 2, "cost_share_rule": 6, "category": 3, "store": 4, "product_group": 9, "product": 160,
		"presentation": 186, "recipe": 60, "expense": 6, "expense_item": 36, "expense_share": 10, "credit": 1, "credit_term_rate": 2, "credit_purchase": 4,
		"credit_purchase_share": 8, "scheduled_payment": 16, "scheduled_payment_share": 32, "promotion": 0, "nutrition_fact": 0, "price": 36} {
		if got := one[int](t, db, "SELECT count(*) FROM "+table); got != want {
			t.Errorf("%s: %d, se esperaban %d", table, got, want)
		}
	}
	// Personas: Karen es externa y no tiene cuenta (disparador).
	if got := one[string](t, db, `SELECT participation FROM person WHERE code = 'karen'`); got != "external" {
		t.Errorf("Karen: %s", got)
	}
	// La compra del 10-oct-2026: 61 piezas, $3,094.50, Axel $1,903.00 y Karen $1,191.50.
	if got := one[int](t, db, `SELECT sum(quantity_milli) FROM expense_item`); got != 61000 {
		t.Errorf("piezas: %d", got)
	}
	if got := one[int](t, db, `SELECT sum(line_total_cents) FROM expense_item`); got != 309450 {
		t.Errorf("importe: %d", got)
	}
	for p, want := range map[string]int{"axel": 190300, "karen": 119150} {
		if got := one[int](t, db, `SELECT s.amount_cents FROM expense_share s JOIN person p ON p.id = s.person_id WHERE p.code = ?`, p); got != want {
			t.Errorf("parte de %s: %d, se esperaban %d", p, got, want)
		}
	}
	if got := one[int](t, db, `SELECT count(*) FROM cost_share_rule WHERE is_template = 0`); got != 3 {
		t.Errorf("reglas de un solo uso (mixto): %d", got)
	}
	for name, bp := range map[string]int{"Mixto: Frijoles bayos La Costeña refritos 220 g": 5000, "Mixto: Salchicha de pavo Sabori Parma 500 g": 5000, "Mixto: Granos de elote Herdez 220 g": 6000} {
		if got := one[int](t, db, `SELECT pct_bp FROM cost_share_rule_part WHERE rule_id = (SELECT id FROM cost_share_rule WHERE name = ?)`, name); got != bp {
			t.Errorf("%s: %d pb, se esperaban %d", name, got, bp)
		}
	}
	// Despensa 1 como la calcula motor_cred: 4 pagos de 973.75 el 30-oct, 15-nov, 30-nov y 15-dic-2026; el último ajusta centavos.
	rows, err := db.Query(`SELECT due_date, amount_cents FROM scheduled_payment WHERE credit_purchase_id = (SELECT id FROM credit_purchase WHERE purchase_date = '2026-10-10') ORDER BY seq`)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for rows.Next() {
		var d string
		var c int
		rows.Scan(&d, &c)
		got = append(got, d+" "+dominio.FormatPesos(int64(c)))
	}
	rows.Close()
	if strings.Join(got, ", ") != "2026-10-30 $973.75, 2026-11-15 $973.75, 2026-11-30 $973.75, 2026-12-15 $973.75" {
		t.Errorf("calendario de la despensa 1: %v", got)
	}
	// C3-09: A, B, primer pago y deuda viva salen del motor sobre la base sembrada.
	c, err := semilla.Verify(context.Background(), db, seedDir(t), "")
	if err != nil {
		t.Fatal(err)
	}
	if c.DisponibleA != 1000000 || c.DisponibleB != 690550 || c.PrimerPago != 97375 || c.PrimerPagoFecha != "2026-10-30" || c.DeudaViva != 629344 || c.DeudaVivaFecha != "2027-04-10" {
		t.Errorf("cruce C3-09: %+v", c)
	}
	// Integridad.
	if got := one[string](t, db, `PRAGMA integrity_check`); got != "ok" {
		t.Errorf("integrity_check: %s", got)
	}
	if n := one[int](t, db, `SELECT count(*) FROM pragma_foreign_key_check`); n != 0 {
		t.Errorf("foreign_key_check: %d", n)
	}
	if rep.Get("recetas").Created != 60 || len(rep.Checks) < 6 {
		t.Errorf("reporte: %+v", rep.Entries)
	}
	// Bitácora: cada alta quedó en audit_log, firmada «semilla» y sin actor.
	if n := one[int](t, db, `SELECT count(*) FROM audit_log WHERE request_id = 'semilla' AND user_account_id IS NULL AND action = 'insert'`); n < 1000 {
		t.Errorf("audit_log: %d", n)
	}
}

func counts(t *testing.T, db *sql.DB) map[string]int {
	t.Helper()
	out := map[string]int{}
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND name NOT LIKE 'goose_%'`)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for rows.Next() {
		var n string
		rows.Scan(&n)
		names = append(names, n)
	}
	rows.Close()
	for _, n := range names {
		out[n] = one[int](t, db, "SELECT count(*) FROM "+n)
	}
	return out
}

func TestSembrarDosVecesNoDuplica(t *testing.T) {
	db := emptyDB(t)
	if _, err := semilla.RunWith(context.Background(), db, opts(t)); err != nil {
		t.Fatal(err)
	}
	first := counts(t, db)
	hash := one[string](t, db, `SELECT password_hash FROM user_account WHERE username = 'axel'`)
	o := opts(t)
	o.Passwords = map[string]string{"axel": "otra-contrasena-distinta", "ximena": "otra-contrasena-distinta-2"} // no se aplican a cuentas existentes
	rep, err := semilla.RunWith(context.Background(), db, o)
	if err != nil {
		t.Fatal(err)
	}
	if rep.TotalCreated() != 0 {
		t.Errorf("la segunda carga creó %d renglones: %+v", rep.TotalCreated(), rep.Entries)
	}
	second := counts(t, db)
	for table, n := range first {
		if second[table] != n {
			t.Errorf("%s: %d tras la primera carga y %d tras la segunda", table, n, second[table])
		}
	}
	if got := one[string](t, db, `SELECT password_hash FROM user_account WHERE username = 'axel'`); got != hash {
		t.Error("la segunda carga cambió la contraseña de una cuenta existente")
	}
	// Con la base ya sembrada tampoco hacen falta contraseñas.
	o.Passwords = nil
	if _, err := semilla.RunWith(context.Background(), db, o); err != nil {
		t.Errorf("sin contraseñas y con las cuentas ya creadas debe pasar: %v", err)
	}
}

func TestContrasenasInicialesNuncaEnClaro(t *testing.T) {
	db := emptyDB(t)
	if _, err := semilla.RunWith(context.Background(), db, opts(t)); err != nil {
		t.Fatal(err)
	}
	h := one[string](t, db, `SELECT password_hash FROM user_account WHERE username = 'ximena'`)
	if !strings.HasPrefix(h, "$argon2id$") {
		t.Fatalf("hash: %.12s", h)
	}
	if ok, err := dominio.VerifyPassword(testPasswords["ximena"], h); err != nil || !ok {
		t.Errorf("la contraseña inicial no verifica: %v %v", ok, err)
	}
	for _, q := range []string{`SELECT count(*) FROM audit_log WHERE after_json LIKE '%' || ? || '%' OR before_json LIKE '%' || ? || '%'`} {
		for _, pw := range testPasswords {
			if n := one[int](t, db, q, pw, pw); n != 0 {
				t.Errorf("la contraseña inicial aparece en audit_log")
			}
		}
	}
}

func TestSinContrasenasNoSiembraNada(t *testing.T) {
	db := emptyDB(t)
	o := opts(t)
	o.Passwords = nil
	_, err := semilla.RunWith(context.Background(), db, o)
	if err == nil || !strings.Contains(err.Error(), "FINANZAS_SEED_PASSWORD_AXEL") {
		t.Fatalf("debía pedir las contraseñas iniciales, dio: %v", err)
	}
	if n := one[int](t, db, `SELECT count(*) FROM person`); n != 0 {
		t.Errorf("la carga fallida dejó %d personas (debe ser todo o nada)", n)
	}
	o.Passwords = map[string]string{"axel": "corta", "ximena": testutil.FakeCredential}
	if _, err := semilla.RunWith(context.Background(), db, o); err == nil || !strings.Contains(err.Error(), "más corta") {
		t.Errorf("una contraseña de menos de 12 caracteres debe rechazarse: %v", err)
	}
}

func TestArchivoDeContrasenas(t *testing.T) {
	// Fuera del repositorio a propósito (t.TempDir): dentro de él la semilla lo rechaza.
	f := filepath.Join(t.TempDir(), "claves.txt")
	body := "# falsas, solo de prueba\naxel=" + testPasswords["axel"] + "\nximena=" + testPasswords["ximena"] + "\n"
	if err := os.WriteFile(f, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv(semilla.EnvPasswordFile, f)
	if _, err := semilla.PasswordsFromEnv(); err == nil || !strings.Contains(err.Error(), "chmod 600") {
		t.Errorf("un archivo legible por otros debe rechazarse: %v", err)
	}
	if err := os.Chmod(f, 0o600); err != nil {
		t.Fatal(err)
	}
	pw, err := semilla.PasswordsFromEnv()
	if err != nil || pw["axel"] != testPasswords["axel"] || pw["ximena"] != testPasswords["ximena"] {
		t.Errorf("lectura del archivo: %v %v", pw, err)
	}
	t.Setenv(semilla.EnvPasswordAxel, "valor-del-entorno-1")
	if pw, _ := semilla.PasswordsFromEnv(); pw["axel"] != "valor-del-entorno-1" {
		t.Error("la variable de entorno debe ganarle al archivo")
	}
	// Dentro del repositorio (datos/tmp) se rechaza aunque el modo sea 0600.
	in := filepath.Join(testutil.TempDir(t), "claves.txt")
	if err := os.WriteFile(in, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(semilla.EnvPasswordFile, in)
	if _, err := semilla.PasswordsFromEnv(); err == nil || !strings.Contains(err.Error(), "repositorio") {
		t.Errorf("un archivo dentro del repositorio debe rechazarse: %v", err)
	}
}

func TestOraculoDistintoEsHallazgoYNoSiembra(t *testing.T) {
	db := emptyDB(t)
	dir := testutil.TempDir(t)
	b, err := os.ReadFile(filepath.Join(seedDir(t), "..", "..", "oraculo", "model.json"))
	if err != nil {
		t.Fatal(err)
	}
	// Cambia el monto del primer pago de la despensa 1 (973.75 → 973.76): el motor y el oráculo ya no coinciden.
	var m map[string]any
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	if err := dec.Decode(&m); err != nil {
		t.Fatal(err)
	}
	pago := m["despensas"].([]any)[0].(map[string]any)["pagos"].([]any)[0].(map[string]any)
	pago["monto"] = json.Number("973.76")
	bad, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	mp := filepath.Join(dir, "model.json")
	if err := os.WriteFile(mp, bad, 0o600); err != nil {
		t.Fatal(err)
	}
	o := opts(t)
	o.ModelPath = mp
	_, err = semilla.RunWith(context.Background(), db, o)
	if err == nil || !strings.Contains(err.Error(), "HALLAZGO") || !strings.Contains(err.Error(), "97376") {
		t.Fatalf("debía reportar el hallazgo con los dos valores: %v", err)
	}
	if n := one[int](t, db, `SELECT count(*) FROM credit_purchase`) + one[int](t, db, `SELECT count(*) FROM person`); n != 0 {
		t.Errorf("la carga con hallazgo dejó %d renglones", n)
	}
}

func TestNutricionSoloConfirmada(t *testing.T) {
	dir := miniSeed(t)
	p := filepath.Join(dir, "nutricion", "tabla.csv")
	b, _ := os.ReadFile(p)
	lines := strings.Split(string(b), "\n")
	// Confirma el renglón 1 (aceite-oliva): 884 kcal, 100 g de grasa.
	lines[1] = strings.Replace(lines[1], ",propuesto,,", ",confirmado,axel,2026-10-12", 1)
	if err := os.WriteFile(p, []byte(strings.Join(lines, "\n")), 0o600); err != nil {
		t.Fatal(err)
	}
	db := emptyDB(t)
	o := opts(t)
	o.Dir = dir
	o.ModelPath = filepath.Join(seedDir(t), "..", "..", "oraculo", "model.json")
	o.OracleSeedPath = oracleSeed(t)
	rep, err := semilla.RunWith(context.Background(), db, o)
	if err != nil {
		t.Fatal(err)
	}
	if n := one[int](t, db, `SELECT count(*) FROM nutrition_fact`); n != 1 {
		t.Fatalf("nutrition_fact: %d (solo el confirmado debe entrar)", n)
	}
	if k := one[int](t, db, `SELECT kcal_milli_per_100g FROM nutrition_fact`); k != 884000 {
		t.Errorf("kcal: %d", k)
	}
	if k := one[int](t, db, `SELECT fat_mg_per_100g FROM nutrition_fact`); k != 100000 {
		t.Errorf("grasa: %d", k)
	}
	if s := one[string](t, db, `SELECT source_license FROM nutrition_fact`); s != "CC0" {
		t.Errorf("licencia: %s", s)
	}
	found := false
	for _, n := range rep.Notices {
		found = found || strings.Contains(n, "149 renglones")
	}
	if !found {
		t.Errorf("el aviso debe decir que 149 siguen propuestos: %v", rep.Notices)
	}
}

func TestPartidasCSVAvisaOImporta(t *testing.T) {
	dir := miniSeed(t)
	model := filepath.Join(seedDir(t), "..", "..", "oraculo", "model.json")
	o := opts(t)
	o.Dir, o.ModelPath, o.OracleSeedPath = dir, model, oracleSeed(t)
	rep, err := semilla.RunWith(context.Background(), emptyDB(t), o)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(rep.Notices, "\n"), "Falta datos/semilla/gastos_fijos_ingresos.csv") {
		t.Errorf("sin CSV debe avisar: %v", rep.Notices)
	}
	// Con el CSV presente y el importador enganchado, se llama una sola vez con la ruta del archivo.
	csv := filepath.Join(dir, "gastos_fijos_ingresos.csv")
	if err := os.WriteFile(csv, []byte("tipo,concepto,monto,periodicidad,dia_o_fecha,categoria,persona,porcentaje_conjunto,notas\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var called []string
	o.ImportPlannedCSV = func(_ context.Context, _ dominio.Querier, path string) (int, error) {
		called = append(called, path)
		return 3, nil
	}
	rep, err = semilla.RunWith(context.Background(), emptyDB(t), o)
	if err != nil {
		t.Fatal(err)
	}
	if len(called) != 1 || called[0] != csv || rep.Get("partidas de ingresos y gastos fijos").Created != 3 {
		t.Errorf("importador: %v %+v", called, rep.Entries)
	}
	// Con el CSV pero sin importador: avisa, no inventa.
	o.ImportPlannedCSV = nil
	rep, err = semilla.RunWith(context.Background(), emptyDB(t), o)
	if err != nil || !strings.Contains(strings.Join(rep.Notices, "\n"), "no tiene el importador enganchado") {
		t.Errorf("CSV sin importador: %v %v", err, rep.Notices)
	}
}

func TestSinDataJSONFalla(t *testing.T) {
	o := opts(t)
	o.Dir = testutil.TempDir(t)
	if _, err := semilla.RunWith(context.Background(), emptyDB(t), o); err == nil || !strings.Contains(err.Error(), "data.json") {
		t.Errorf("sin data.json debía fallar claro: %v", err)
	}
}
