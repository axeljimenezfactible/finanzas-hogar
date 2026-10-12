package ocr

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// Prueba de punta a punta con tesseract y magick reales sobre los tickets
// SINTÉTICOS. Se omite si no están instalados. No mide S5-06 (esas no son
// fotos reales): solo comprueba que la tubería completa funciona.
func TestIntegracionSinteticos(t *testing.T) {
	for _, p := range []string{"tesseract", "magick"} {
		if _, err := exec.LookPath(p); err != nil {
			t.Skipf("%s no instalado", p)
		}
	}
	dir := filepath.Join("..", "..", "..", "..", "datos", "semilla", "tickets-sinteticos")
	eng := New()
	eng.Timeout = 60 * time.Second
	eng.TempDir = t.TempDir()
	rep, err := Measure(context.Background(), eng, MeasureOptions{
		TicketsDir: dir, ExpectedDir: filepath.Join(dir, "esperado"), Synthetic: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !rep.Measured {
		t.Skipf("sin sintéticos: %s", rep.Reason)
	}
	t.Logf("\n%s", rep.Format())
	if rep.Correct == 0 {
		t.Errorf("la tubería real no acertó ni un renglón sintético:\n%s", rep.Format())
	}
}
