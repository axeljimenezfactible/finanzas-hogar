// Package jobs es el planificador (S6-08): time.Ticker + tabla job, sin dependencias externas. Cada job se registra por
// código con su función; el renglón de job guarda horario (cron en America/Mexico_City), estado y última ejecución.
// Un disparo se toma con una actualización condicional (next_run_at esperado), así que dos disparos simultáneos o
// repetidos ejecutan una sola vez; tras un reinicio con el job vencido corre una vez y salta a la siguiente fecha.
package jobs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"sort"
	"sync"
	"time"

	"finanzas-hogar/internal/dominio"
)

// Run es el contexto de una corrida.
type Run struct {
	Code      string
	Scheduled time.Time // el disparo que se está atendiendo (next_run_at que se tomó)
	Now       time.Time
	RequestID string // va a audit_log.request_id de los cambios del job
}

// Func es el trabajo de un job; un error queda en job.last_error (resumido, sin secretos).
type Func func(ctx context.Context, run Run) error

// Def define un job: código, nombre visible, horario inicial (solo si el renglón no existe) y función.
type Def struct {
	Code     string
	Name     string
	Schedule string
	Run      Func
}

// Scheduler corre los jobs registrados.
type Scheduler struct {
	db       *sql.DB
	clock    dominio.Clock
	log      *slog.Logger
	mu       sync.Mutex
	defs     map[string]Def
	Interval time.Duration // cada cuánto revisa (30 s)
}

var reCode = regexp.MustCompile(`^[a-z0-9_]{1,40}$`)

// New crea el planificador.
func New(db *sql.DB, clock dominio.Clock, log *slog.Logger) *Scheduler {
	if log == nil {
		log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	return &Scheduler{db: db, clock: clock, log: log, defs: map[string]Def{}, Interval: 30 * time.Second}
}

// Register agrega un job por código (a-z, 0-9 y _; único).
func (s *Scheduler) Register(d Def) error {
	if !reCode.MatchString(d.Code) || d.Name == "" || d.Run == nil {
		return fmt.Errorf("jobs: definición inválida %q", d.Code)
	}
	if _, err := ParseSchedule(d.Schedule); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, dup := s.defs[d.Code]; dup {
		return fmt.Errorf("jobs: %s ya está registrado", d.Code)
	}
	s.defs[d.Code] = d
	return nil
}

// Codes lista los códigos registrados.
func (s *Scheduler) Codes() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, 0, len(s.defs))
	for k := range s.defs {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Prepare crea los renglones que falten y pasa a error los que quedaron en running por un reinicio.
func (s *Scheduler) Prepare(ctx context.Context) error {
	for _, c := range s.Codes() {
		d := s.defs[c]
		if err := dominio.EnsureJob(ctx, s.db, d.Code, d.Name, d.Schedule); err != nil {
			return err
		}
	}
	return dominio.ResetInterruptedJobs(ctx, s.db)
}

// Start prepara y revisa cada Interval hasta que ctx termine (en su propia gorrutina).
func (s *Scheduler) Start(ctx context.Context) error {
	if err := s.Prepare(ctx); err != nil {
		return err
	}
	go func() {
		t := time.NewTicker(s.Interval)
		defer t.Stop()
		for {
			if _, err := s.Tick(ctx); err != nil && !errors.Is(err, context.Canceled) {
				s.log.Error("jobs: revisión fallida", "error", err.Error())
			}
			select {
			case <-ctx.Done():
				return
			case <-t.C:
			}
		}
	}()
	return nil
}

// Tick revisa una vez con la hora del reloj y corre los jobs vencidos. Devuelve cuántos corrió.
func (s *Scheduler) Tick(ctx context.Context) (int, error) {
	now := s.clock.Now()
	ran := 0
	for _, code := range s.Codes() {
		j, err := dominio.GetJobByCode(ctx, s.db, code)
		if errors.Is(err, dominio.ErrNotFound) {
			continue
		}
		if err != nil {
			return ran, err
		}
		if !j.IsEnabled {
			continue
		}
		sched, err := ParseSchedule(j.Schedule)
		if err != nil {
			_ = dominio.FinishJob(ctx, s.db, j.ID, "error", "Horario inválido: "+j.Schedule)
			continue
		}
		if !j.NextRunAt.Valid {
			if err := dominio.SetJobNextRun(ctx, s.db, j.ID, sched.Next(now)); err != nil {
				return ran, err
			}
			continue
		}
		due, err := dominio.ParseDBTime(j.NextRunAt.String)
		if err != nil || now.Before(due) {
			continue
		}
		ok, err := dominio.ClaimJob(ctx, s.db, j.ID, j.NextRunAt.String, now, sched.Next(now))
		if err != nil {
			return ran, err
		}
		if !ok {
			continue
		}
		ran++
		s.execute(ctx, j.ID, Run{Code: code, Scheduled: due.In(dominio.MX), Now: now, RequestID: fmt.Sprintf("job:%s:%d", code, due.Unix())})
	}
	return ran, nil
}

func (s *Scheduler) execute(ctx context.Context, id int64, run Run) {
	status, msg := "ok", ""
	func() {
		defer func() {
			if v := recover(); v != nil {
				status, msg = "error", fmt.Sprint("pánico: ", v)
			}
		}()
		if err := s.defs[run.Code].Run(ctx, run); err != nil {
			status, msg = "error", err.Error()
		}
	}()
	if len(msg) > 500 {
		msg = msg[:500]
	}
	if err := dominio.FinishJob(context.WithoutCancel(ctx), s.db, id, status, msg); err != nil {
		s.log.Error("jobs: no se pudo anotar el resultado", "job", run.Code, "error", err.Error())
	}
	s.log.Info("jobs: corrida", "job", run.Code, "status", status)
}

// RunNow encola una corrida: pone next_run_at = ahora y el próximo Tick la toma (p. ej. «correr ahora» desde Ajustes).
func (s *Scheduler) RunNow(ctx context.Context, code string) error {
	j, err := dominio.GetJobByCode(ctx, s.db, code)
	if err != nil {
		return err
	}
	return dominio.SetJobNextRun(ctx, s.db, j.ID, s.clock.Now())
}
