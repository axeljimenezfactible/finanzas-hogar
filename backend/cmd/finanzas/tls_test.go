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

// S6-14: solo 127.0.0.1 y la LAN 192.168.50.0/24; nunca comodín ni Tailscale.
func TestEscucha_SoloLANYLocal(t *testing.T) {
	for _, bad := range []string{"0.0.0.0:8443", "[::]:8443", ":8443", "100.101.102.103:8443", "192.168.1.5:8443", "10.0.0.2:8443", "finanzas.local:8443", "127.0.0.1:8443,0.0.0.0:8443", ""} {
		if _, err := parseAddrs(bad); err == nil {
			t.Errorf("%q se aceptó", bad)
		}
	}
	got, err := parseAddrs("192.168.50.10:8443, 127.0.0.1:8443")
	if err != nil || len(got) != 2 {
		t.Fatalf("%v %v", got, err)
	}
	if _, err := loadTLS("", ""); err == nil {
		t.Error("sin certificado debe fallar (salvo --dev-http)")
	}
}

// TLS: el binario sirve HTTPS con el certificado que se le da y en todas las direcciones pedidas.
func TestEscucha_HTTPSConCertificado(t *testing.T) {
	dir := testutil.TempDir(t)
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	tpl := &x509.Certificate{SerialNumber: big.NewInt(7), Subject: pkix.Name{CommonName: "finanzas.local"}, DNSNames: []string{"finanzas.local"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	der, _ := x509.CreateCertificate(rand.Reader, tpl, tpl, &key.PublicKey, key)
	kder, _ := x509.MarshalPKCS8PrivateKey(key)
	certFile, keyFile := filepath.Join(dir, "cert.pem"), filepath.Join(dir, "key.pem")
	os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600)
	os.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: kder}), 0o600)
	tc, err := loadTLS(certFile, keyFile)
	if err != nil {
		t.Fatal(err)
	}
	// dos puertos libres en 127.0.0.1
	var addrs []string
	for i := 0; i < 2; i++ {
		l, _ := net.Listen("tcp", "127.0.0.1:0")
		addrs = append(addrs, l.Addr().String())
		l.Close()
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- serveAll(ctx, addrs, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "ok") }), tc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	}()
	pool := x509.NewCertPool()
	cert, _ := x509.ParseCertificate(der)
	pool.AddCert(cert)
	hc := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool, ServerName: "finanzas.local"}}, Timeout: 5 * time.Second}
	for _, a := range addrs {
		var res *http.Response
		for i := 0; i < 50; i++ {
			if res, err = hc.Get("https://" + a + "/"); err == nil {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
		if err != nil || res.StatusCode != 200 || res.TLS == nil {
			t.Fatalf("%s: %v", a, err)
		}
		res.Body.Close()
		if r, err := http.Get("http://" + a + "/"); err == nil && r.StatusCode == 200 {
			t.Errorf("%s responde sin TLS", a)
		}
	}
	cancel()
	if err := <-done; err != nil {
		t.Errorf("apagado: %v", err)
	}
}
