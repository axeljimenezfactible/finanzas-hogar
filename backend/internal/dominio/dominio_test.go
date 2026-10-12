package dominio_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/testutil"
)

var ctx = context.Background()

func TestMigracionesYApertura(t *testing.T) {
	d, fx := testutil.NewDB(t)
	if fx.AxelUserID == 0 || fx.XimenaUserID == 0 || fx.KarenPersonID == 0 {
		t.Fatalf("fixtures incompletas: %+v", fx)
	}
	v, err := dominio.Migrate(ctx, d) // segunda vez: nada pendiente
	if err != nil || v < 4 {
		t.Fatalf("Migrate otra vez: v=%d err=%v", v, err)
	}
	if _, err := dominio.Open(filepath.Join(testutil.TempDir(t), "no", "existe", "x.db")); err == nil {
		t.Error("Open en carpeta inexistente debió fallar")
	}
	if _, err := dominio.OpenReadOnly(filepath.Join(testutil.TempDir(t), "x.db")); err == nil {
		t.Error("OpenReadOnly de un archivo inexistente debió fallar")
	}
	// Karen no puede tener cuenta (disparador de S1) y la violación se reconoce como restricción.
	_, err = d.Exec(`INSERT INTO user_account (person_id, username, password_hash) VALUES (?, 'karen', '$argon2id$x')`, fx.KarenPersonID)
	if !dominio.IsConstraint(err) {
		t.Errorf("cuenta para Karen: se esperaba restricción, salió %v", err)
	}
	_, err = d.Exec(`INSERT INTO person (code, display_name, participation) VALUES ('axel', 'Otro', 'external')`)
	if !dominio.IsUnique(err) || dominio.IsUnique(errors.New("x")) {
		t.Errorf("IsUnique: %v", err)
	}
	err = dominio.WithTx(ctx, d, func(tx *sql.Tx) error {
		_, _ = tx.Exec(`UPDATE person SET notes = 'x' WHERE id = ?`, fx.AxelPersonID)
		return errors.New("revertir")
	})
	var notes sql.NullString
	_ = d.QueryRow(`SELECT notes FROM person WHERE id = ?`, fx.AxelPersonID).Scan(&notes)
	if err == nil || notes.Valid {
		t.Errorf("WithTx no revirtió: %v %v", err, notes)
	}
}

func TestErrores(t *testing.T) {
	cases := []struct {
		e    *dominio.Error
		kind dominio.Kind
		code string
	}{
		{dominio.NotFound("x"), dominio.KindNotFound, "not_found"}, {dominio.Conflict("in_use", "x"), dominio.KindConflict, "in_use"},
		{dominio.Forbidden("read_only_source", "x"), dominio.KindForbidden, "read_only_source"}, {dominio.Unavailable("x"), dominio.KindUnavailable, "unavailable"},
		{dominio.BadRequest("x"), dominio.KindBadRequest, "bad_request"}, {dominio.Invalid(dominio.FieldIssue{Field: "a"}), dominio.KindValidation, "validation_failed"},
	}
	for _, c := range cases {
		if c.e.Kind != c.kind || c.e.Code != c.code || c.e.Error() == "" {
			t.Errorf("%+v", c.e)
		}
	}
	wrapped := &dominio.Error{Kind: dominio.KindConflict, Code: "conflict", Message: "m", Err: sql.ErrConnDone}
	if !errors.Is(wrapped, sql.ErrConnDone) || wrapped.Error() == "" {
		t.Error("Unwrap")
	}
	if e, ok := dominio.AsError(sql.ErrNoRows); !ok || e.Kind != dominio.KindNotFound {
		t.Error("sql.ErrNoRows debe ser 404")
	}
	if _, ok := dominio.AsError(errors.New("x")); ok {
		t.Error("error común no es de negocio")
	}
}

func TestFechasYPeriodos(t *testing.T) {
	d, err := dominio.ParseDate("2026-10-31")
	if err != nil || d.String() != "2026-10-31" || dominio.FortnightOf(d).String() != "2026-10-q2" {
		t.Fatalf("31 de octubre en q2: %v %v", d, err)
	}
	for _, bad := range []string{"2026-02-30", "10-oct-2026", "2026-1-01"} {
		if _, err := dominio.ParseDate(bad); err == nil {
			t.Errorf("ParseDate(%q) debió fallar", bad)
		}
	}
	feb := dominio.Fortnight{Year: 2027, Month: time.February, Half: 2}
	if feb.PayDay().String() != "2027-02-28" || feb.End().String() != "2027-02-28" || feb.Start().Day != 16 {
		t.Errorf("febrero 2027: pago %s fin %s", feb.PayDay(), feb.End())
	}
	q1 := dominio.Fortnight{Year: 2026, Month: time.October, Half: 1}
	if q1.PayDay().Day != 15 || q1.Next().String() != "2026-10-q2" || q1.Prev().String() != "2026-09-q2" || q1.Start().Day != 1 || q1.End().Day != 15 {
		t.Error("quincena 1")
	}
	if (dominio.Fortnight{Year: 2026, Month: time.October, Half: 2}).PayDay().Day != 30 || dominio.DaysIn(2028, time.February) != 29 {
		t.Error("pago del 30 / bisiesto")
	}
	if !d.After(d.AddDays(-1)) || !d.AddDays(-1).Before(d) || d.Compare(d) != 0 || d.DaysUntil(d.AddDays(10)) != 10 || d.IsZero() {
		t.Error("comparaciones")
	}
	if d.ISOWeek() != "2026-W44" || d.MonthKey() != "2026-10" || d.Weekday() != time.Saturday || dominio.NewDate(2026, 2, 31).String() != "2026-03-03" {
		t.Errorf("semana %s", d.ISOWeek())
	}
	if d.StartOfDay().Format(dominio.DateTimeLayout) != "2026-10-31T00:00:00-06:00" {
		t.Error("StartOfDay")
	}
	clk := testutil.NewClock(time.Date(2026, 10, 11, 3, 30, 0, 0, time.UTC))
	if dominio.Today(clk).String() != "2026-10-10" {
		t.Errorf("hoy en MX: %s", dominio.Today(clk))
	}
	ts, err := dominio.ParseDateTime("2026-10-11T03:30:00Z")
	if err != nil || dominio.FormatDateTime(ts) != "2026-10-10T21:30:00-06:00" {
		t.Errorf("Z → -06:00: %v", err)
	}
	for _, bad := range []string{"2026-10-10 21:30", "2026-13-10T21:30:00Z"} {
		if _, err := dominio.ParseDateTime(bad); err == nil {
			t.Errorf("ParseDateTime(%q)", bad)
		}
	}
	db := dominio.FormatDBTime(ts)
	if db != "2026-10-11T03:30:00.000Z" || dominio.DBTimeToAPI(db) != "2026-10-10T21:30:00-06:00" || dominio.DBTimeToAPI("x") != "x" {
		t.Errorf("DBTime %s", db)
	}
	if _, err := dominio.ParseDBTime("2026-10-11T03:30:00Z"); err != nil {
		t.Error(err)
	}
	periods := map[string][2]string{"2026-10": {"2026-10-01", "2026-10-31"}, "2027-02-q2": {"2027-02-16", "2027-02-28"}, "2026-W42": {"2026-10-12", "2026-10-18"},
		"2026": {"2026-01-01", "2026-12-31"}, "2026-10-10..2026-11-09": {"2026-10-10", "2026-11-09"}, "2027-W53": {}}
	for s, want := range periods {
		p, err := dominio.ParsePeriod(s)
		if want[0] == "" {
			if err == nil {
				t.Errorf("%s debió fallar (2027 tiene 52 semanas): %v", s, p)
			}
			continue
		}
		if err != nil || p.Start.String() != want[0] || p.End.String() != want[1] || !p.Contains(p.End) || p.Days() < 1 {
			t.Errorf("%s: %+v %v", s, p, err)
		}
	}
	if _, err := dominio.ParsePeriod("2026-11-09..2026-10-10"); err == nil || err.(*dominio.Error).Kind != dominio.KindValidation {
		t.Errorf("rango invertido = 422: %v", err)
	}
	for _, bad := range []string{"oct-2026", "2026-02-30..2026-03-01"} {
		if _, err := dominio.ParsePeriod(bad); err == nil || err.(*dominio.Error).Kind != dominio.KindBadRequest {
			t.Errorf("%s = 400: %v", bad, err)
		}
	}
	if h, m, err := dominio.ParseHHMM("21:30"); err != nil || h != 21 || m != 30 {
		t.Error("HH:MM")
	}
	for _, bad := range []string{"24:00", "9:00", "ab:cd"} {
		if _, _, err := dominio.ParseHHMM(bad); err == nil {
			t.Errorf("ParseHHMM(%q)", bad)
		}
	}
}

func TestDinero(t *testing.T) {
	cases := [][4]int64{{5, 2, 0, 3}, {-5, 2, 0, -3}, {5, -2, 0, -3}, {7, 3, 0, 2}, {125868, 1, 0, 125868}}
	for _, c := range cases {
		if got := dominio.RoundHalfUp(c[0], c[1]); got != c[3] {
			t.Errorf("RoundHalfUp(%d,%d)=%d, se esperaba %d", c[0], c[1], got, c[3])
		}
	}
	if dominio.MulDivRound(503474, 1, 4) != 125869 || dominio.ApplyBasisPoints(19900, 1600) != 3184 || dominio.PercentOf(80, 0) != 0 || dominio.PercentOf(805, 1000) != 81 {
		t.Error("MulDivRound/ApplyBasisPoints/PercentOf")
	}
	for c, want := range map[int64]string{97375: "$973.75", 125870: "$1,258.70", -5: "-$0.05", 0: "$0.00", 123456789: "$1,234,567.89"} {
		if got := dominio.FormatPesos(c); got != want {
			t.Errorf("FormatPesos(%d)=%s", c, got)
		}
	}
	for _, f := range []func(){func() { dominio.RoundHalfUp(1, 0) }, func() { dominio.MulDivRound(1, 1, 0) }} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("denominador 0 debe entrar en pánico")
				}
			}()
			f()
		}()
	}
}

func TestAjustes(t *testing.T) {
	d, fx := testutil.NewDB(t)
	defs, err := dominio.ParamDefs()
	if err != nil || len(defs) < 20 {
		t.Fatalf("ParamDefs: %d %v", len(defs), err)
	}
	if defs["attachment.max_size_mb"].ValueType() != "integer" || defs["alerts.payment_due_days_before"].ValueType() != "integer_list" ||
		defs["alerts.send_time"].ValueType() != "time" || defs["comparator.timeout_seconds"].Range == "" {
		t.Error("tipos o rangos de @param")
	}
	for k, want := range map[string]int64{"auth.max_failed_attempts": 5, "auth.lockout_window_minutes": 15, "session.max_age_days": 90, "attachment.max_size_mb": 10, "idempotency.retention_hours": 24, "auth.password_min_length": 12} {
		if n, err := dominio.SettingInt(ctx, d, k); err != nil || n != want {
			t.Errorf("%s = %d (%v), se esperaba %d", k, n, err, want)
		}
	}
	if l, err := dominio.SettingIntList(ctx, d, "alerts.payment_due_days_before"); err != nil || len(l) != 2 {
		t.Error("lista", l, err)
	}
	if s, err := dominio.SettingTime(ctx, d, "alerts.send_time"); err != nil || s != "09:00" {
		t.Error("hora", s, err)
	}
	if b, err := dominio.SettingBool(ctx, d, "ai.claude_api_enabled"); err != nil || b {
		t.Error("IA apagada por defecto", b, err)
	}
	if _, err := dominio.SettingInt(ctx, d, "no.existe"); err == nil {
		t.Error("clave inexistente")
	}
	if _, err := dominio.SettingTime(ctx, d, "auth.max_failed_attempts"); err == nil {
		t.Error("5 no es HH:MM")
	}
	if _, err := dominio.SettingInt(ctx, d, "optimizer.weights"); err == nil {
		t.Error("objeto no es entero")
	}
	actor := fx.AxelUserID
	sec := []string{"auth.max_failed_attempts", "auth.password_min_length", "comparator.timeout_seconds", "alerts.send_time", "alerts.payment_due_days_before"}
	apply := func(ch dominio.ParamChange) error {
		return dominio.WithTx(ctx, d, func(tx *sql.Tx) error { return dominio.ApplyParams(ctx, tx, ch, sec, &actor, "req-1") })
	}
	if err := apply(dominio.ParamChange{Set: map[string]json.RawMessage{"auth.password_min_length": json.RawMessage(`16`), "alerts.send_time": json.RawMessage(`"08:15"`)}}); err != nil {
		t.Fatal(err)
	}
	p, err := dominio.ReadParameter(ctx, d, "auth.password_min_length")
	if err != nil || string(p.Value) != "16" || string(p.InitialValue) != "12" || p.IsInitial || p.Unit != "caracteres" {
		t.Errorf("ReadParameter tras cambio: %+v %v", p, err)
	}
	if err := apply(dominio.ParamChange{Reset: []string{"auth.password_min_length"}}); err != nil {
		t.Fatal(err)
	}
	if p, _ := dominio.ReadParameter(ctx, d, "auth.password_min_length"); !p.IsInitial {
		t.Error("restablecer")
	}
	if _, err := dominio.ReadParameter(ctx, d, "ai.claude_api_enabled"); err == nil {
		t.Error("ai.* no tiene @param")
	}
	var n int
	_ = d.QueryRow(`SELECT COUNT(*) FROM audit_log WHERE table_name = 'app_setting' AND request_id = 'req-1' AND user_account_id = ?`, actor).Scan(&n)
	if n != 3 {
		t.Errorf("audit_log de ajustes = %d, se esperaban 3", n)
	}
	bad := []dominio.ParamChange{
		{Set: map[string]json.RawMessage{"price.stale_after_days": json.RawMessage(`10`)}},                                     // fuera de la sección
		{Set: map[string]json.RawMessage{"auth.max_failed_attempts": json.RawMessage(`"5"`)}},                                  // tipo
		{Set: map[string]json.RawMessage{"auth.max_failed_attempts": json.RawMessage(`null`)}},                                 // null
		{Set: map[string]json.RawMessage{"alerts.payment_due_days_before": json.RawMessage(`null`)}},                           // lista null
		{Set: map[string]json.RawMessage{"alerts.send_time": json.RawMessage(`"25:00"`)}},                                      // hora
		{Set: map[string]json.RawMessage{"comparator.timeout_seconds": json.RawMessage(`31`)}},                                 // rango (disparador)
		{Set: map[string]json.RawMessage{"alerts.send_time": json.RawMessage(`"08:00"`)}, Reset: []string{"alerts.send_time"}}, // ambos
		{Reset: []string{"no.existe"}},
	}
	for i, ch := range bad {
		err := apply(ch)
		if e, ok := dominio.AsError(err); !ok || e.Kind != dominio.KindValidation || len(e.Fields) == 0 {
			t.Errorf("cambio malo %d: %v", i, err)
		}
	}
	err = dominio.WithTx(ctx, d, func(tx *sql.Tx) error {
		if err := dominio.SetSetting(ctx, tx, "ai.local_model_enabled", json.RawMessage(`true`), &actor, ""); err != nil {
			return err
		}
		if err := dominio.SetSetting(ctx, tx, "optimizer.min_available_buffer_cents", json.RawMessage(`100000`), nil, ""); err != nil {
			return err
		}
		return dominio.SetSetting(ctx, tx, "x.y", json.RawMessage(`{`), nil, "")
	})
	if err == nil {
		t.Error("JSON inválido en SetSetting")
	}
	if _, err := dominio.RowJSON(ctx, d, "app_setting; DROP", 1); err == nil {
		t.Error("RowJSON con tabla inválida")
	}
	if _, err := dominio.RowJSON(ctx, d, "app_setting", 99999); !errors.Is(err, dominio.ErrNotFound) {
		t.Error("RowJSON inexistente")
	}
}

func TestCuentasSesionesEIntentos(t *testing.T) {
	d, fx := testutil.NewDB(t)
	u, err := dominio.UserByUsername(ctx, d, "AXEL")
	if err != nil || u.ID != fx.AxelUserID || u.DisplayName != "Axel" || u.PasswordHash[:10] != "$argon2id$" {
		t.Fatalf("UserByUsername: %+v %v", u, err)
	}
	if _, err := dominio.UserByUsername(ctx, d, "karen"); !errors.Is(err, dominio.ErrNotFound) {
		t.Error("Karen no tiene cuenta")
	}
	us, err := dominio.ListUsers(ctx, d)
	if err != nil || len(us) != 2 {
		t.Fatalf("2 cuentas: %d %v", len(us), err)
	}
	if ok, err := dominio.VerifyPassword(testutil.FakeCredential, u.PasswordHash); !ok || err != nil {
		t.Error("VerifyPassword correcta", err)
	}
	if ok, _ := dominio.VerifyPassword("otra", u.PasswordHash); ok {
		t.Error("VerifyPassword incorrecta")
	}
	for _, bad := range []string{"x", "$argon2i$v=19$m=1,t=1,p=1$a$b", "$argon2id$v=1$m=1,t=1,p=1$a$b", "$argon2id$v=19$m=x$a$b", "$argon2id$v=19$m=1,t=1,p=1$!$b"} {
		if _, err := dominio.VerifyPassword("x", bad); err == nil {
			t.Errorf("hash malo %q", bad)
		}
	}
	h, err := dominio.HashPassword("nueva-contrasena-falsa", testutil.FastArgon2)
	now := testutil.MX(2026, time.October, 12, 12, 0)
	if err != nil || dominio.SetPasswordHash(ctx, d, u.ID, h, now) != nil || dominio.TouchLogin(ctx, d, u.ID, now) != nil {
		t.Fatal("cambio de contraseña")
	}
	if u2, _ := dominio.UserByID(ctx, d, u.ID); u2.PasswordHash != h || !u2.LastLoginAt.Valid {
		t.Error("SetPasswordHash/TouchLogin")
	}
	tok, hash, err := dominio.NewSessionToken()
	if err != nil || dominio.HashToken(tok) != hash || len(hash) != 64 {
		t.Fatal("token")
	}
	day := 24 * time.Hour
	id, err := dominio.CreateSession(ctx, d, u.ID, hash, "Pixel", "Chrome", "192.168.50.20", now, 90*day)
	if err != nil {
		t.Fatal(err)
	}
	_, hash2, _ := dominio.NewSessionToken()
	id2, _ := dominio.CreateSession(ctx, d, fx.XimenaUserID, hash2, "", "", "", now, 90*day)
	if s, err := dominio.ActiveSession(ctx, d, hash, now.Add(60*day)); err != nil || s.ID != id {
		t.Fatal("sesión vigente el día 60", err)
	}
	if err := dominio.RenewSession(ctx, d, id, "192.168.50.21", now.Add(60*day), 90*day); err != nil {
		t.Fatal(err)
	}
	s, _ := dominio.GetSession(ctx, d, id)
	if exp, _ := dominio.ParseDBTime(s.ExpiresAt); !exp.Equal(now.Add(150*day)) || s.IPAddress.String != "192.168.50.21" {
		t.Errorf("renovación al día 150: %s", s.ExpiresAt)
	}
	if _, err := dominio.ActiveSession(ctx, d, hash2, now.Add(91*day)); !errors.Is(err, dominio.ErrNotFound) {
		t.Error("sin uso, el día 91 vence")
	}
	if l, err := dominio.ListSessions(ctx, d, now); err != nil || len(l) != 2 {
		t.Errorf("ListSessions: %d %v", len(l), err)
	}
	if err := dominio.RevokeSession(ctx, d, id2, now); err != nil || dominio.RevokeSession(ctx, d, id2, now) == nil {
		t.Error("RevokeSession doble")
	}
	if n, err := dominio.RevokeOtherSessions(ctx, d, 0, now); err != nil || n != 1 {
		t.Errorf("RevokeOtherSessions: %d %v", n, err)
	}
	if _, err := dominio.ActiveSession(ctx, d, hash, now); !errors.Is(err, dominio.ErrNotFound) {
		t.Error("sesión cerrada")
	}
	if _, err := dominio.GetSession(ctx, d, 9999); !errors.Is(err, dominio.ErrNotFound) {
		t.Error("GetSession inexistente")
	}
	// S6-06 (b): fallos en los minutos 0, 1, 2, 3 y 14; rechazado en 14.5 y aceptado en 15.5.
	t0 := testutil.MX(2026, time.October, 12, 21, 0)
	for _, m := range []float64{0, 1, 2, 3, 14} {
		if err := dominio.RecordLoginAttempt(ctx, d, "axel", &u.ID, "192.168.50.20", false, t0.Add(time.Duration(m*float64(time.Minute)))); err != nil {
			t.Fatal(err)
		}
	}
	lk, err := dominio.CheckLockout(ctx, d, "Axel", "192.168.50.99", t0.Add(14*time.Minute+30*time.Second), 5, 15*time.Minute)
	if err != nil || !lk.Locked || !lk.RetryAt.Equal(t0.Add(15*time.Minute)) {
		t.Errorf("minuto 14.5: %+v %v", lk, err)
	}
	if lk, _ := dominio.CheckLockout(ctx, d, "axel", "", t0.Add(15*time.Minute+30*time.Second), 5, 15*time.Minute); lk.Locked {
		t.Error("minuto 15.5 debe aceptar")
	}
	// Por IP: otra cuenta desde la misma IP también queda bloqueada.
	if lk, _ := dominio.CheckLockout(ctx, d, "ximena", "192.168.50.20", t0.Add(14*time.Minute+30*time.Second), 5, 15*time.Minute); !lk.Locked {
		t.Error("bloqueo por IP")
	}
	if err := dominio.RecordLoginAttempt(ctx, d, "nadie", nil, "", true, t0); err != nil {
		t.Error(err)
	}
}

func TestIdempotenciaYJobs(t *testing.T) {
	d, fx := testutil.NewDB(t)
	now := testutil.MX(2026, time.October, 12, 12, 0)
	nb := now.Add(-24 * time.Hour)
	key := "3f6c2a8e-7b1d-4c9a-9e2f-5a7d1b0c4e81"
	id, err := dominio.ReserveIdempotency(ctx, d, fx.AxelUserID, key, "POST", "/api/v1/expenses", dominio.HashToken("a"), now, nb)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dominio.ReserveIdempotency(ctx, d, fx.AxelUserID, key, "POST", "/api/v1/expenses", dominio.HashToken("a"), now, nb); err == nil {
		t.Error("reserva doble")
	}
	if err := dominio.CompleteIdempotency(ctx, d, id, 201, []byte(`{"id":812}`)); err != nil {
		t.Fatal(err)
	}
	r, err := dominio.FindIdempotency(ctx, d, fx.AxelUserID, key, nb)
	if err != nil || r.Status.Int64 != 201 || r.Body.String != `{"id":812}` || r.Method != "POST" {
		t.Fatalf("FindIdempotency: %+v %v", r, err)
	}
	if _, err := dominio.FindIdempotency(ctx, d, fx.AxelUserID, key, now.Add(time.Hour)); !errors.Is(err, dominio.ErrNotFound) {
		t.Error("vencida = no existe")
	}
	if n, err := dominio.PurgeIdempotency(ctx, d, now.Add(time.Hour)); err != nil || n != 1 {
		t.Errorf("purga %d %v", n, err)
	}
	id, _ = dominio.ReserveIdempotency(ctx, d, fx.XimenaUserID, key, "POST", "/x", dominio.HashToken("b"), now, nb)
	if err := dominio.ReleaseIdempotency(ctx, d, id); err != nil {
		t.Fatal(err)
	}
	if err := dominio.CompleteIdempotency(ctx, d, id, 200, nil); err != nil {
		t.Error(err)
	}
	if err := dominio.EnsureJob(ctx, d, "backup_daily", "Respaldo diario", "15 3 * * *"); err != nil {
		t.Fatal(err)
	}
	_ = dominio.EnsureJob(ctx, d, "backup_daily", "Otro nombre", "0 0 * * *")
	j, err := dominio.GetJobByCode(ctx, d, "backup_daily")
	if err != nil || j.Name != "Respaldo diario" || !j.IsEnabled {
		t.Fatalf("EnsureJob no cambia uno existente: %+v %v", j, err)
	}
	next := now.Add(time.Hour)
	if err := dominio.SetJobNextRun(ctx, d, j.ID, next); err != nil {
		t.Fatal(err)
	}
	j, _ = dominio.GetJob(ctx, d, j.ID)
	ok1, err1 := dominio.ClaimJob(ctx, d, j.ID, j.NextRunAt.String, next, next.Add(24*time.Hour))
	ok2, _ := dominio.ClaimJob(ctx, d, j.ID, j.NextRunAt.String, next, next.Add(24*time.Hour))
	if !ok1 || err1 != nil || ok2 {
		t.Errorf("ClaimJob doble: %v %v %v", ok1, ok2, err1)
	}
	if err := dominio.ResetInterruptedJobs(ctx, d); err != nil {
		t.Fatal(err)
	}
	if j, _ = dominio.GetJob(ctx, d, j.ID); j.LastStatus.String != "error" {
		t.Error("ResetInterruptedJobs")
	}
	if err := dominio.FinishJob(ctx, d, j.ID, "ok", ""); err != nil {
		t.Fatal(err)
	}
	if l, err := dominio.ListJobs(ctx, d); err != nil || len(l) != 1 || l[0].LastStatus.String != "ok" {
		t.Error("ListJobs", err)
	}
	if _, err := dominio.GetJob(ctx, d, 999); !errors.Is(err, dominio.ErrNotFound) {
		t.Error("GetJob inexistente")
	}
}

// S6-09: VACUUM INTO, integrity_check = ok y restauración con las mismas filas por tabla que el origen.
func TestRespaldoYRestauracion(t *testing.T) {
	d, _ := testutil.NewDB(t)
	dir := testutil.TempDir(t)
	now := testutil.MX(2026, time.October, 12, 3, 15)
	before, err := dominio.TableCounts(ctx, d)
	if err != nil {
		t.Fatal(err)
	}
	var res dominio.BackupResult
	for i := 0; i < 3; i++ {
		res, err = dominio.Backup(ctx, d, filepath.Join(dir, "respaldos"), filepath.Join(dir, "tmp"), "scheduled", now.Add(time.Duration(i)*24*time.Hour))
		if err != nil {
			t.Fatal(err)
		}
	}
	if res.Integrity != "ok" || res.RestoreResult != "ok" || res.SizeBytes == 0 || len(res.SHA256) != 64 {
		t.Fatalf("respaldo: %+v", res)
	}
	for tbl, n := range before {
		want := n
		if tbl == "backup" {
			want = n + 3 // el renglón de cada respaldo vive en la base
		}
		if tbl != "backup" && res.Counts[tbl] != want {
			t.Errorf("tabla %s: restaurada %d, origen %d", tbl, res.Counts[tbl], want)
		}
	}
	if ic, err := dominio.IntegrityCheck(ctx, res.Path); err != nil || ic != "ok" {
		t.Error("integrity_check", ic, err)
	}
	var status, sha string
	_ = d.QueryRow(`SELECT status, sha256 FROM backup WHERE id = ?`, res.ID).Scan(&status, &sha)
	if status != "ok" || sha != res.SHA256 {
		t.Errorf("renglón backup: %s %s", status, sha)
	}
	n, err := dominio.PruneBackups(ctx, d, "scheduled", 2, now)
	if err != nil || n != 1 {
		t.Errorf("retención: %d %v", n, err)
	}
	files, _ := filepath.Glob(filepath.Join(dir, "respaldos", "*.db"))
	if len(files) != 2 {
		t.Errorf("quedan %d archivos", len(files))
	}
	if _, err := dominio.Backup(ctx, d, filepath.Join(files[0], "x"), dir, "manual", now); err == nil {
		t.Error("carpeta inválida")
	}
	os.WriteFile(filepath.Join(dir, "roto.db"), []byte("no es sqlite"), 0o600)
	if _, _, err := dominio.RestoreTest(ctx, filepath.Join(dir, "roto.db"), dir); err == nil {
		t.Error("archivo roto")
	}
}
