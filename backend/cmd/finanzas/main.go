// Comando finanzas: servidor HTTP local de Finanzas Hogar (esqueleto de S6, WF-0).
//
// Registra todas las operaciones del contrato C1 (docs/contratos/openapi.yaml): GET /health y
// GET /openapi.yaml ya responden; las demás devuelven 501 con el esquema único de error.
// Pendiente de S6 (WF-3): sesión por cookie, TLS con mkcert en finanzas.local:8443, base SQLite,
// jobs y respaldos. Mientras tanto escucha solo en 127.0.0.1 (bloqueo criterios-bind).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	_ "time/tzdata" // zona America/Mexico_City aunque el sistema no tenga tzdata

	"finanzas-hogar/internal/api"
)

func main() {
	addr := flag.String("addr", envOr("FINANZAS_ADDR", "127.0.0.1:8080"), "dirección de escucha (host:puerto)")
	spec := flag.String("openapi", envOr("FINANZAS_OPENAPI", "../docs/contratos/openapi.yaml"), "ruta del contrato OpenAPI que se sirve en /openapi.yaml")
	flag.Parse()

	loc, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		log.Fatalf("no se pudo cargar la zona America/Mexico_City: %v", err)
	}

	mux := http.NewServeMux()
	api.Register(mux, map[string]http.Handler{
		"getHealth":          healthHandler(loc),
		"getOpenApiDocument": openAPIHandler(*spec),
	})

	srv := &http.Server{
		Addr:              *addr,
		Handler:           withContractErrors(mux),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	log.Printf("finanzas: escuchando en http://%s (%d operaciones del contrato C1)", *addr, len(api.Routes))
	log.Fatal(srv.ListenAndServe())
}

// healthHandler responde el esquema Health del contrato: {"status":"ok","time":...}.
func healthHandler(loc *time.Location) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"time":   time.Now().In(loc).Format(time.RFC3339),
		})
	})
}

// openAPIHandler sirve el contrato C1 leyéndolo del disco en cada petición.
func openAPIHandler(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile(path)
		if err != nil {
			log.Printf("finanzas: no se pudo leer el contrato %q: %v", path, err)
			api.WriteError(w, r, http.StatusServiceUnavailable, "unavailable", "El contrato de la API no está disponible en este momento.")
			return
		}
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		_, _ = w.Write(b)
	})
}

// withContractErrors hace que las rutas que no son del contrato (404) y los métodos no
// admitidos (405) también respondan con el esquema único de error de C1.
func withContractErrors(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, pattern := mux.Handler(r); pattern == "" {
			if r.Header.Get("X-Request-Id") == "" {
				r.Header.Set("X-Request-Id", fmt.Sprintf("%x", time.Now().UnixNano()))
			}
			w.Header().Set("X-Request-Id", r.Header.Get("X-Request-Id"))
			w = &contractErrorWriter{ResponseWriter: w, r: r}
		}
		mux.ServeHTTP(w, r)
	})
}

type contractErrorWriter struct {
	http.ResponseWriter
	r        *http.Request
	replaced bool
}

func (c *contractErrorWriter) WriteHeader(code int) {
	switch code {
	case http.StatusNotFound:
		c.replaced = true
		api.WriteError(c.ResponseWriter, c.r, code, "not_found", "La ruta no existe.")
	case http.StatusMethodNotAllowed:
		c.replaced = true
		api.WriteError(c.ResponseWriter, c.r, code, "method_not_allowed", "Método no permitido en esta ruta.")
	default:
		c.ResponseWriter.WriteHeader(code)
	}
}

func (c *contractErrorWriter) Write(b []byte) (int, error) {
	if c.replaced {
		return len(b), nil
	}
	return c.ResponseWriter.Write(b)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
