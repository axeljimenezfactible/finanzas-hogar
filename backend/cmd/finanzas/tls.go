package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"time"
)

// lanNet es la red de la casa (decisión H y K29: solo LAN 192.168.50.0/24; Tailscale no se abre).
var (
	lanNet   = netip.MustParsePrefix("192.168.50.0/24")
	loopback = netip.MustParseAddr("127.0.0.1")
)

// parseAddrs parte «IP:puerto,IP:puerto» (FINANZAS_ADDR) y valida cada dirección (S6-14): solo 127.0.0.1 o una IP de
// 192.168.50.0/24; nunca comodín (0.0.0.0, ::, «:8443») ni la red de Tailscale (100.64.0.0/10) ni otra red.
func parseAddrs(list string) ([]string, error) {
	var out []string
	for _, a := range strings.Split(list, ",") {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		if err := checkLAN(a); err != nil {
			return nil, err
		}
		host, _, _ := net.SplitHostPort(a)
		ip, err := netip.ParseAddr(host)
		if err != nil || !(ip == loopback || lanNet.Contains(ip)) {
			return nil, fmt.Errorf("dirección %q: solo 127.0.0.1 o una IP de la LAN 192.168.50.0/24 (S6-14; Tailscale no se abre, K29)", a)
		}
		out = append(out, a)
	}
	if len(out) == 0 {
		return nil, errors.New("falta la dirección de escucha (FINANZAS_ADDR, p. ej. 192.168.50.10:8443,127.0.0.1:8443)")
	}
	return out, nil
}

// alpnProtos son los protocolos que el servidor ofrece por ALPN: HTTP/2 primero y HTTP/1.1 para los clientes que no
// lo hablan (WF-3c, api-tablero). Sin NextProtos el listener TLS no anuncia h2 y todo cae a HTTP/1.1; con h2 el tablero
// de Inicio hace sus peticiones en paralelo sobre una sola conexión. http.Server ya trae el manejador de h2
// (TLSNextProto) porque su TLSConfig es nil (newHTTPServer); el TLS lo pone el listener de serveAll.
var alpnProtos = []string{"h2", "http/1.1"}

// loadTLS carga el certificado de mkcert (finanzas.local); el certificado real lo emite el integrador (ops/mkcert.md).
func loadTLS(certFile, keyFile string) (*tls.Config, error) {
	if certFile == "" || keyFile == "" {
		return nil, errors.New("falta el certificado de mkcert: usa --cert y --key (ops/mkcert.md) o --dev-http solo en 127.0.0.1")
	}
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("certificado de mkcert: %w", err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS12, NextProtos: alpnProtos}, nil
}

func newHTTPServer(h http.Handler) *http.Server {
	return &http.Server{Handler: h, ReadHeaderTimeout: 10 * time.Second, ReadTimeout: 10 * time.Minute, WriteTimeout: 10 * time.Minute, IdleTimeout: 2 * time.Minute}
}

// serveAll escucha en cada dirección con el mismo manejador (HTTPS si tc no es nil); si una falla o ctx termina, se
// apagan todas.
func serveAll(ctx context.Context, addrs []string, h http.Handler, tc *tls.Config, log *slog.Logger) error {
	var lns []net.Listener
	for _, a := range addrs {
		ln, err := net.Listen("tcp", a)
		if err != nil {
			for _, l := range lns {
				l.Close()
			}
			return err
		}
		if tc != nil {
			ln = tls.NewListener(ln, tc)
		}
		lns = append(lns, ln)
	}
	errc := make(chan error, len(lns))
	srvs := make([]*http.Server, 0, len(lns))
	for _, ln := range lns {
		srv := newHTTPServer(h)
		srvs = append(srvs, srv)
		go func(ln net.Listener) { errc <- srv.Serve(ln) }(ln)
		log.Info("finanzas: escuchando", "addr", ln.Addr().String(), "tls", tc != nil)
	}
	var err error
	select {
	case <-ctx.Done():
	case err = <-errc:
	}
	shut, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	for _, s := range srvs {
		_ = s.Shutdown(shut)
	}
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
