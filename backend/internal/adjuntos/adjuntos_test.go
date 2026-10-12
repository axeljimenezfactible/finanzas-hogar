package adjuntos_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"finanzas-hogar/internal/adjuntos"
	"finanzas-hogar/internal/testutil"
)

func TestAdjunto_NombresYRutas(t *testing.T) {
	for _, bad := range []string{"../x.png", `..\x.png`, "a/b.png", "..", ".", "", "x\x00.png", strings.Repeat("a", 256)} {
		if adjuntos.CheckName(bad) == nil {
			t.Errorf("nombre %q aceptado", bad)
		}
	}
	if adjuntos.CheckName("Ticket del súper (1).jpg") != nil {
		t.Error("nombre normal rechazado")
	}
	data := testutil.TempDir(t)
	st := adjuntos.Store{DataDir: data}
	png := append([]byte("\x89PNG\r\n\x1a\n"), bytes.Repeat([]byte{0}, 100)...)
	s, err := st.Save(bytes.NewReader(png), 1024, testutil.MX(2026, time.October, 12, 12, 0))
	if err != nil || s.MimeType != "image/png" || s.Size != int64(len(png)) || !strings.HasPrefix(s.RelPath, "adjuntos/2026/10/") {
		t.Fatalf("%+v %v", s, err)
	}
	if p, err := st.Resolve(s.RelPath); err != nil || !strings.HasPrefix(p, data) {
		t.Errorf("resolver: %v", err)
	}
	os.WriteFile(filepath.Join(data, "fuera.png"), png, 0o600)
	os.Symlink(filepath.Join(data, "fuera.png"), filepath.Join(data, "adjuntos", "enlace.png"))
	for _, bad := range []string{"adjuntos/../fuera.png", "/etc/passwd", "fuera.png", "adjuntos/enlace.png", "adjuntos", `adjuntos\..\fuera.png`} {
		if _, err := st.Resolve(bad); err == nil {
			t.Errorf("ruta %q aceptada", bad)
		}
	}
	if _, err := st.Save(bytes.NewReader(png), 50, time.Now()); !errors.Is(err, adjuntos.ErrTooLarge) {
		t.Errorf("tamaño: %v", err)
	}
	if _, err := st.Save(strings.NewReader("texto simple"), 50, time.Now()); !errors.Is(err, adjuntos.ErrType) {
		t.Errorf("tipo: %v", err)
	}
	if err := st.Remove(s.RelPath); err != nil {
		t.Error(err)
	}
	if err := st.Remove(s.RelPath); err != nil {
		t.Errorf("borrar dos veces: %v", err)
	}
}
