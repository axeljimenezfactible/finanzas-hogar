package jobs_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/jobs"
	"finanzas-hogar/internal/testutil"
)

var ctx = context.Background()

// S6-08: día 30 y febrero de 2027 (28), doble disparo = 1 ejecución, reinicio con job vencido = 1 ejecución, todo en America/Mexico_City.
func TestJobsDia30Febrero2027EIdempotencia(t *testing.T) {
	db, _ := testutil.NewDB(t)
	clk := testutil.NewClock(testutil.MX(2026, time.October, 30, 8, 59))
	var runs atomic.Int32
	var scheduled []time.Time
	var mu sync.Mutex
	fn := func(ctx context.Context, r jobs.Run) error {
		runs.Add(1)
		mu.Lock()
		scheduled = append(scheduled, r.Scheduled)
		mu.Unlock()
		return nil
	}
	s := jobs.New(db, clk, nil)
	if err := s.Register(jobs.Def{Code: "prueba_diaria", Name: "Prueba diaria", Schedule: "0 9 * * *", Run: fn}); err != nil {
		t.Fatal(err)
	}
	if err := s.Register(jobs.Def{Code: "prueba_diaria", Name: "x", Schedule: "0 9 * * *", Run: fn}); err == nil {
		t.Error("código duplicado")
	}
	for _, bad := range []jobs.Def{{Code: "Mal", Name: "x", Schedule: "* * * * *", Run: fn}, {Code: "x", Name: "x", Schedule: "61 * * * *", Run: fn}} {
		if s.Register(bad) == nil {
			t.Errorf("definición inválida aceptada: %+v", bad.Code)
		}
	}
	if err := s.Prepare(ctx); err != nil {
		t.Fatal(err)
	}
	if n, _ := s.Tick(ctx); n != 0 {
		t.Fatal("el primer tick solo calcula next_run_at")
	}
	j, _ := dominio.GetJobByCode(ctx, db, "prueba_diaria")
	if j.NextRunAt.String != "2026-10-30T15:00:00.000Z" {
		t.Fatalf("next_run_at = %s; se esperaba 09:00 de México (15:00Z)", j.NextRunAt.String)
	}
	clk.Set(testutil.MX(2026, time.October, 30, 9, 0))
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ { // disparos simultáneos
		wg.Add(1)
		go func() { defer wg.Done(); _, _ = s.Tick(ctx) }()
	}
	wg.Wait()
	if n, _ := s.Tick(ctx); n != 0 || runs.Load() != 1 {
		t.Fatalf("doble disparo: %d ejecuciones", runs.Load())
	}
	j, _ = dominio.GetJobByCode(ctx, db, "prueba_diaria")
	if j.LastStatus.String != "ok" || j.NextRunAt.String != "2026-10-31T15:00:00.000Z" || !scheduled[0].Equal(testutil.MX(2026, time.October, 30, 9, 0)) {
		t.Errorf("tras correr: %+v", j)
	}
	// Reinicio con el job vencido 3 días: otro planificador sobre la misma base corre una vez y salta al futuro.
	clk.Set(testutil.MX(2026, time.November, 3, 10, 0))
	s2 := jobs.New(db, clk, nil)
	_ = s2.Register(jobs.Def{Code: "prueba_diaria", Name: "Prueba diaria", Schedule: "0 9 * * *", Run: fn})
	_, _ = db.Exec(`UPDATE job SET last_status = 'running' WHERE code = 'prueba_diaria'`) // quedó corriendo al apagarse
	if err := s2.Prepare(ctx); err != nil {
		t.Fatal(err)
	}
	n1, _ := s2.Tick(ctx)
	n2, _ := s2.Tick(ctx)
	if n1 != 1 || n2 != 0 || runs.Load() != 2 {
		t.Fatalf("reinicio: %d y %d", n1, n2)
	}
	if j, _ = dominio.GetJobByCode(ctx, db, "prueba_diaria"); j.NextRunAt.String != "2026-11-04T15:00:00.000Z" {
		t.Errorf("tras reinicio next_run_at = %s", j.NextRunAt.String)
	}
	// Febrero de 2027 (28 días): 21:30 del 28 en México es 03:30Z del 1 de marzo.
	clk.Set(testutil.MX(2027, time.February, 28, 21, 0))
	_ = s2.Register(jobs.Def{Code: "fin_de_febrero", Name: "Fin de febrero", Schedule: "30 21 28 2 *", Run: fn})
	_ = s2.Prepare(ctx)
	_, _ = s2.Tick(ctx)
	if j, _ = dominio.GetJobByCode(ctx, db, "fin_de_febrero"); j.NextRunAt.String != "2027-03-01T03:30:00.000Z" {
		t.Fatalf("febrero 2027: %s", j.NextRunAt.String)
	}
	clk.Set(testutil.MX(2027, time.February, 28, 21, 30))
	before := runs.Load()
	_, _ = s2.Tick(ctx)
	_, _ = s2.Tick(ctx)
	if runs.Load()-before != 1 { // solo fin_de_febrero: prueba_diaria vencida ya corrió una vez en el tick de las 21:00
		t.Errorf("febrero: %d corridas", runs.Load()-before)
	}
	if j, _ = dominio.GetJobByCode(ctx, db, "fin_de_febrero"); j.NextRunAt.String != "2028-02-29T03:30:00.000Z" {
		t.Errorf("siguiente 28 de febrero: %s", j.NextRunAt.String)
	}
}

func TestJobsErroresApagadoYRunNow(t *testing.T) {
	db, _ := testutil.NewDB(t)
	clk := testutil.NewClock(testutil.MX(2026, time.October, 12, 12, 0))
	s := jobs.New(db, clk, nil)
	s.Interval = 10 * time.Millisecond
	_ = s.Register(jobs.Def{Code: "falla", Name: "Falla", Schedule: "*/5 * * * *", Run: func(context.Context, jobs.Run) error { return errors.New("sin red") }})
	_ = s.Register(jobs.Def{Code: "panico", Name: "Pánico", Schedule: "0 0 1 1 *", Run: func(context.Context, jobs.Run) error { panic("x") }})
	_ = s.Register(jobs.Def{Code: "apagado", Name: "Apagado", Schedule: "* * * * *", Run: func(context.Context, jobs.Run) error { t.Error("no debe correr"); return nil }})
	if got := s.Codes(); len(got) != 3 || got[0] != "apagado" {
		t.Errorf("Codes: %v", got)
	}
	c, cancel := context.WithCancel(ctx)
	if err := s.Start(c); err != nil {
		t.Fatal(err)
	}
	time.Sleep(30 * time.Millisecond)
	cancel()
	_, _ = db.Exec(`UPDATE job SET is_enabled = 0 WHERE code = 'apagado'`)
	for _, code := range []string{"falla", "panico", "apagado"} {
		if err := s.RunNow(ctx, code); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.RunNow(ctx, "no_existe"); err == nil {
		t.Error("RunNow de un job inexistente")
	}
	if n, err := s.Tick(ctx); err != nil || n != 2 {
		t.Fatalf("RunNow: %d %v", n, err)
	}
	for code, want := range map[string]string{"falla": "sin red", "panico": "pánico: x"} {
		if j, _ := dominio.GetJobByCode(ctx, db, code); j.LastStatus.String != "error" || j.LastError.String != want {
			t.Errorf("%s: %s %q", code, j.LastStatus.String, j.LastError.String)
		}
	}
	_, _ = db.Exec(`UPDATE job SET schedule = 'cada rato', is_enabled = 1 WHERE code = 'apagado'`)
	_, _ = s.Tick(ctx)
	if j, _ := dominio.GetJobByCode(ctx, db, "apagado"); j.LastStatus.String != "error" {
		t.Error("horario inválido en la tabla")
	}
}

func TestCron(t *testing.T) {
	base := testutil.MX(2026, time.October, 12, 12, 0) // lunes
	cases := map[string]time.Time{
		"*/15 * * * *":     testutil.MX(2026, time.October, 12, 12, 15),
		"0 9 * * 0":        testutil.MX(2026, time.October, 18, 9, 0),
		"0 9 * * 7":        testutil.MX(2026, time.October, 18, 9, 0),
		"0 8 1 * *":        testutil.MX(2026, time.November, 1, 8, 0),
		"0 8 1 * 3":        testutil.MX(2026, time.October, 14, 8, 0), // día 1 o miércoles
		"30 21 * 1-3 *":    testutil.MX(2027, time.January, 1, 21, 30),
		"0,30 13-14 * * *": testutil.MX(2026, time.October, 12, 13, 0),
	}
	for expr, want := range cases {
		s, err := jobs.ParseSchedule(expr)
		if err != nil {
			t.Fatalf("%s: %v", expr, err)
		}
		if got := s.Next(base); !got.Equal(want) {
			t.Errorf("%s: %s, se esperaba %s", expr, got, want)
		}
	}
	for _, bad := range []string{"* * * *", "a * * * *", "*/0 * * * *", "5-1 * * * *", "0 24 * * *", "0 0 0 * *"} {
		if _, err := jobs.ParseSchedule(bad); err == nil {
			t.Errorf("%q debió fallar", bad)
		}
	}
	if s, _ := jobs.ParseSchedule("0 0 31 2 *"); !s.Next(base).IsZero() {
		t.Error("31 de febrero nunca llega")
	}
}
