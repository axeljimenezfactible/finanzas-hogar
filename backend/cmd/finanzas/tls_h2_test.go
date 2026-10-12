package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"finanzas-hogar/internal/testutil"
)

// certPrueba crea un certificado autofirmado de prueba para finanzas.local y 127.0.0.1 (archivos falsos en
// datos/tmp/, nunca el de mkcert) y devuelve sus rutas y el almacén con el que el cliente lo acepta.
func certPrueba(t *testing.T) (certFile, keyFile string, pool *x509.CertPool) {
	t.Helper()
	dir := testutil.TempDir(t)
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tpl := &x509.Certificate{SerialNumber: big.NewInt(11), Subject: pkix.Name{CommonName: "finanzas.local"}, DNSNames: []string{"finanzas.local"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	kder, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	certFile, keyFile = filepath.Join(dir, "cert.pem"), filepath.Join(dir, "key.pem")
	if err := os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: kder}), 0o600); err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	pool = x509.NewCertPool()
	pool.AddCert(cert)
	return certFile, keyFile, pool
}

// HTTP/2 por ALPN (WF-3c): el mismo serveAll de producción negocia h2 con un cliente que lo ofrece y sigue sirviendo
// HTTP/1.1 al que solo habla HTTP/1.1. S6-14 no cambia: la dirección sigue siendo 127.0.0.1 (TestEscucha_*).
func TestEscucha_NegociaHTTP2(t *testing.T) {
	certFile, keyFile, pool := certPrueba(t)
	tc, err := loadTLS(certFile, keyFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(tc.NextProtos) != 2 || tc.NextProtos[0] != "h2" || tc.NextProtos[1] != "http/1.1" {
		t.Fatalf("NextProtos = %v, se esperaba [h2 http/1.1]", tc.NextProtos)
	}
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- serveAll(ctx, []string{addr}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, r.Proto) }),
			tc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	}()
	defer func() {
		cancel()
		if err := <-done; err != nil {
			t.Errorf("apagado: %v", err)
		}
	}()

	// 1) Cliente TLS crudo que ofrece h2 y http/1.1: el servidor elige h2.
	var conn *tls.Conn
	for i := 0; i < 50; i++ {
		if conn, err = tls.Dial("tcp", addr, &tls.Config{RootCAs: pool, ServerName: "finanzas.local", NextProtos: []string{"h2", "http/1.1"}}); err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("conexión TLS: %v", err)
	}
	if got := conn.ConnectionState().NegotiatedProtocol; got != "h2" {
		t.Errorf("ALPN negoció %q, se esperaba h2", got)
	}
	conn.Close()

	// 2) Cliente HTTP con HTTP/2: la petición viaja como HTTP/2.0.
	h2 := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool, ServerName: "finanzas.local"}, ForceAttemptHTTP2: true}}
	res, err := h2.Get("https://" + addr + "/")
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	if res.ProtoMajor != 2 || string(body) != "HTTP/2.0" {
		t.Errorf("cliente h2: protocolo %s, el manejador vio %q", res.Proto, body)
	}

	// 3) Cliente que solo habla HTTP/1.1: sigue funcionando (respaldo de ALPN).
	h1 := &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool, ServerName: "finanzas.local", NextProtos: []string{"http/1.1"}},
		TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{}}}
	res, err = h1.Get("https://" + addr + "/")
	if err != nil {
		t.Fatal(err)
	}
	body, _ = io.ReadAll(res.Body)
	res.Body.Close()
	if res.ProtoMajor != 1 || string(body) != "HTTP/1.1" {
		t.Errorf("cliente HTTP/1.1: protocolo %s, el manejador vio %q", res.Proto, body)
	}
}
