// Comando finanzas: servidor local de Finanzas Hogar (S6) y sus tareas de operación.
//
//	finanzas serve    servidor HTTPS (mkcert) con la API de C1, la PWA (frontend/dist) y los jobs
//	finanzas migrate  aplica las migraciones embebidas (goose como biblioteca)
//	finanzas semilla  carga datos/semilla/ (esqueleto: llama a internal/semilla cuando exista)
//	finanzas respaldo respaldo con VACUUM INTO + integrity_check + restauración de prueba (S6-09)
//
// Solo red local (S6-14, K29): serve rechaza las direcciones comodín (0.0.0.0, ::, «:8443»); la dirección LAN
// (192.168.50.x:8443) la fija WF-4 en ops/runbook.md con FINANZAS_ADDR. Rutas por omisión: PROVISIONAL (WF-4).
package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/sistema"
	"finanzas-hogar/internal/dominio"
	"finanzas-hogar/internal/jobs"
)

// seedRunner lo instala internal/semilla cuando exista (fuera del alcance de backend-base).
var seedRunner func(ctx context.Context, db *sql.DB, dir string) error

// Valores PROVISIONALES (K31, puerta 4) de los jobs base; la justificación está en docs/estado/wf3/backend-base.md.
const (
	backupSchedule = "15 3 * * *" // 03:15 de México: casa dormida, después del resumen diario (21:30) y antes del reporte mensual (08:00)
	backupKeep     = 30           // 30 respaldos diarios: un mes para volver atrás
	purgeSchedule  = "5 * * * *"  // purga de Idempotency-Key cada hora (retención: idempotency.retention_hours)
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	var err error
	switch os.Args[1] {
	case "serve":
		err = serve(log, os.Args[2:])
	case "migrate":
		err = migrate(os.Args[2:])
	case "semilla":
		err = semilla(os.Args[2:])
	case "respaldo":
		err = respaldo(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "finanzas:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "uso: finanzas serve|migrate|semilla|respaldo [opciones]  (finanzas <subcomando> -h para ver las opciones)")
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func openDB(path string) (*sql.DB, error) {
	d, err := dominio.Open(path)
	if err != nil {
		return nil, err
	}
	if _, err := dominio.Migrate(context.Background(), d); err != nil {
		d.Close()
		return nil, err
	}
	return d, nil
}

// checkLAN rechaza las direcciones comodín: el servicio escucha solo en la IP de la LAN que se le dé (S6-14).
func checkLAN(addr string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("dirección %q inválida: %w", addr, err)
	}
	ip := net.ParseIP(host)
	if host == "" || ip == nil || ip.IsUnspecified() {
		return fmt.Errorf("dirección %q: usa la IP de la LAN (192.168.50.x) o 127.0.0.1, nunca un comodín", addr)
	}
	return nil
}

func serve(log *slog.Logger, args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	dbPath := fs.String("db", envOr("FINANZAS_DB", "datos/finanzas.db"), "base SQLite")
	addr := fs.String("addr", envOr("FINANZAS_ADDR", "127.0.0.1:8443"), "IP:puerto de escucha, separadas por coma: 127.0.0.1 y la IP de la LAN 192.168.50.x (sin comodín)")
	cert := fs.String("cert", envOr("FINANZAS_CERT", ""), "certificado de mkcert (finanzas.local); obligatorio salvo --dev-http")
	ca := fs.String("ca", envOr("FINANZAS_CA", ""), "rootCA.pem de mkcert que se descarga desde Ajustes (getCaCertificate)")
	devHTTP := fs.Bool("dev-http", false, "sin TLS y solo en 127.0.0.1 (desarrollo)")
	key := fs.String("key", envOr("FINANZAS_KEY", ""), "llave del certificado de mkcert")
	front := fs.String("frontend", envOr("FINANZAS_FRONTEND", "frontend/dist"), "carpeta de la PWA compilada")
	spec := fs.String("openapi", envOr("FINANZAS_OPENAPI", "docs/contratos/openapi.yaml"), "contrato C1 que se sirve en /openapi.yaml")
	data := fs.String("datos", envOr("FINANZAS_DATOS", "datos/app"), "carpeta de adjuntos y reportes")
	backups := fs.String("respaldos", envOr("FINANZAS_RESPALDOS", "datos/respaldos"), "carpeta de respaldos (solo local)")
	_ = fs.Parse(args)
	addrs, err := parseAddrs(*addr)
	if err != nil {
		return err
	}
	var tc *tls.Config
	if *devHTTP {
		for _, a := range addrs {
			if !strings.HasPrefix(a, "127.0.0.1:") {
				return fmt.Errorf("--dev-http solo escucha en 127.0.0.1, no en %s", a)
			}
		}
	} else if tc, err = loadTLS(*cert, *key); err != nil {
		return err
	}
	sistema.CACertPath = *ca
	d, err := openDB(*dbPath)
	if err != nil {
		return err
	}
	defer d.Close()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	clock := dominio.SystemClock{}
	sched := jobs.New(d, clock, log)
	tmp := *data + "/tmp"
	if err := baseJobs(sched, d, *backups, tmp); err != nil {
		return err
	}
	deps := &api.Deps{DB: d, Clock: clock, Log: log, Jobs: sched, Config: api.Config{DataDir: *data, BackupDir: *backups, TmpDir: tmp,
		OpenAPIPath: *spec, FrontendDir: existing(*front), SecureHeaders: true}}
	if err := os.MkdirAll(tmp, 0o700); err != nil {
		return err
	}
	reg := api.NewRegistry(deps)
	for _, area := range areas {
		area(reg)
	}
	api.RegisterBase(reg)
	if err := sched.Start(ctx); err != nil {
		return err
	}
	log.Info("finanzas: arranque", "direcciones", strings.Join(addrs, ","), "tls", tc != nil, "operaciones", len(api.Routes), "implementadas", len(reg.Implemented()))
	return serveAll(ctx, addrs, reg.Handler(), tc, log)
}

func existing(dir string) string {
	if st, err := os.Stat(dir); err == nil && st.IsDir() {
		return dir
	}
	return ""
}

// baseJobs registra los jobs de backend-base: respaldo diario y purga de Idempotency-Key.
func baseJobs(s *jobs.Scheduler, d *sql.DB, backups, tmp string) error {
	err := s.Register(jobs.Def{Code: "backup_daily", Name: "Respaldo diario", Schedule: backupSchedule, Run: func(ctx context.Context, r jobs.Run) error {
		if _, err := dominio.Backup(ctx, d, backups, tmp, "scheduled", r.Now); err != nil {
			return err
		}
		_, err := dominio.PruneBackups(ctx, d, "scheduled", backupKeep, r.Now)
		return err
	}})
	if err != nil {
		return err
	}
	return s.Register(jobs.Def{Code: "idempotency_purge", Name: "Purga de claves de idempotencia", Schedule: purgeSchedule, Run: func(ctx context.Context, r jobs.Run) error {
		h, err := dominio.SettingInt(ctx, d, "idempotency.retention_hours")
		if err != nil {
			return err
		}
		_, err = dominio.PurgeIdempotency(ctx, d, r.Now.Add(-time.Duration(h)*time.Hour))
		return err
	}})
}

func migrate(args []string) error {
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)
	dbPath := fs.String("db", envOr("FINANZAS_DB", "datos/finanzas.db"), "base SQLite")
	_ = fs.Parse(args)
	d, err := dominio.Open(*dbPath)
	if err != nil {
		return err
	}
	defer d.Close()
	v, err := dominio.Migrate(context.Background(), d)
	if err == nil {
		fmt.Printf("migraciones aplicadas: versión %d\n", v)
	}
	return err
}

func semilla(args []string) error {
	fs := flag.NewFlagSet("semilla", flag.ExitOnError)
	dbPath := fs.String("db", envOr("FINANZAS_DB", "datos/finanzas.db"), "base SQLite")
	dir := fs.String("dir", "datos/semilla", "carpeta de la semilla real")
	dev := fs.Bool("desarrollo", false, "carga SOLO el overlay de desarrollo (presupuesto, sueldo, plan de comidas…) sobre una base ya sembrada; nunca en producción")
	conGasto := fs.Bool("con-gasto-fijo", false, "con --desarrollo: agrega la renta mensual (rompe las cifras P1–P6 de model.json)")
	_ = fs.Parse(args)
	if *conGasto && !*dev {
		return errors.New("--con-gasto-fijo solo va con --desarrollo")
	}
	if *dev {
		if devSeedRunner == nil {
			return errors.New("la semilla de desarrollo no está instalada (falta internal/semilla)")
		}
		d, err := openDB(*dbPath)
		if err != nil {
			return err
		}
		defer d.Close()
		return devSeedRunner(context.Background(), d, *dir, *conGasto)
	}
	if seedRunner == nil {
		return errors.New("la carga de la semilla todavía no está implementada (falta internal/semilla)")
	}
	d, err := openDB(*dbPath)
	if err != nil {
		return err
	}
	defer d.Close()
	return seedRunner(context.Background(), d, *dir)
}

func respaldo(args []string) error {
	fs := flag.NewFlagSet("respaldo", flag.ExitOnError)
	dbPath := fs.String("db", envOr("FINANZAS_DB", "datos/finanzas.db"), "base SQLite")
	dir := fs.String("respaldos", envOr("FINANZAS_RESPALDOS", "datos/respaldos"), "carpeta de respaldos (solo local)")
	tmp := fs.String("tmp", "datos/tmp", "carpeta temporal de la restauración de prueba")
	kind := fs.String("tipo", "manual", "manual, scheduled o pre_migration")
	_ = fs.Parse(args)
	d, err := openDB(*dbPath)
	if err != nil {
		return err
	}
	defer d.Close()
	res, err := dominio.Backup(context.Background(), d, *dir, *tmp, *kind, time.Now())
	if err != nil {
		return err
	}
	fmt.Printf("respaldo %d: %s (%d bytes, sha256 %s)\nintegrity_check: %s\nrestauración de prueba: %s\n", res.ID, res.Path, res.SizeBytes, res.SHA256, res.Integrity, res.RestoreResult)
	for t, n := range res.Counts {
		fmt.Printf("  %s\t%d\n", t, n)
	}
	return nil
}
