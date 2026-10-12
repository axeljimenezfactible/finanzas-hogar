package ocr

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type call struct {
	name string
	args []string
}

type fakeRunner struct {
	mu      sync.Mutex
	calls   []call
	version string
	langs   string
	tsv     string
	failOn  string // nombre del programa que falla
	delay   time.Duration
}

func (f *fakeRunner) Run(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
	f.mu.Lock()
	f.calls = append(f.calls, call{name, args})
	f.mu.Unlock()
	if f.delay > 0 {
		select {
		case <-time.After(f.delay):
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
	if f.failOn == name {
		return nil, []byte("falló a propósito"), errors.New("exit status 1")
	}
	switch {
	case name == "tesseract" && len(args) == 1 && args[0] == "--version":
		return nil, []byte(f.version), nil // la versión sale por stderr, como en tesseract
	case name == "tesseract" && len(args) == 1 && args[0] == "--list-langs":
		return []byte(f.langs), nil, nil
	case name == "tesseract":
		return []byte(f.tsv), nil, nil
	case name == "magick":
		return nil, nil, os.WriteFile(args[len(args)-1], []byte("png"), 0o600)
	}
	return nil, nil, errors.New("programa inesperado " + name)
}

func (f *fakeRunner) count(name string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, c := range f.calls {
		if c.name == name {
			n++
		}
	}
	return n
}

func okRunner() *fakeRunner {
	return &fakeRunner{
		version: "tesseract 5.5.3\n leptonica-1.87.0\n",
		langs:   "List of available languages in \"/usr/share/tessdata/\" (2):\nosd\nspa\n",
		tsv: tsvHead +
			tsvRow("5", "1", "1", "1", "1", "95", "LECHE") +
			tsvRow("5", "1", "1", "1", "2", "95", "28.50"),
	}
}

func tmpImage(t *testing.T, name string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// S5-06: TestOCRParseo
func TestOCRParseo(t *testing.T) {
	r := okRunner()
	e := &Engine{Runner: r, Profile: ProfileGeneric}
	lines, err := e.Recognize(context.Background(), tmpImage(t, "a.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 1 || lines[0].Descripcion != "LECHE" || *lines[0].ImporteCents != 2850 || lines[0].Confianza != 95 {
		t.Fatalf("%+v", lines)
	}
	// la última llamada es tesseract con la configuración documentada
	last := r.calls[len(r.calls)-1]
	got := strings.Join(last.args[1:], " ")
	want := "stdout -l spa --psm 6 --oem 1 -c preserve_interword_spaces=1 tsv"
	if last.name != "tesseract" || got != want || !filepath.IsAbs(last.args[0]) {
		t.Errorf("llamada: %v", last)
	}
	if r.count("magick") != 0 {
		t.Error("sin Preprocess no debe llamar a magick")
	}
}

// S5-06: TestOCRVersion
func TestOCRVersion(t *testing.T) {
	ctx := context.Background()
	cases := map[string]*fakeRunner{
		"version distinta": {version: "tesseract 5.3.0\n", langs: "spa\n"},
		"sin version":      {version: "otra cosa\n", langs: "spa\n"},
		"sin spa":          {version: "tesseract 5.5.3\n", langs: "osd\neng\n"},
		"no instalado":     {failOn: "tesseract"},
	}
	for name, r := range cases {
		t.Run(name, func(t *testing.T) {
			e := &Engine{Runner: r}
			if err := e.CheckAvailable(ctx); !errors.Is(err, ErrOCRUnavailable) {
				t.Errorf("err = %v, quiero ErrOCRUnavailable", err)
			}
			if _, err := e.Recognize(ctx, tmpImage(t, "a.png")); !errors.Is(err, ErrOCRUnavailable) {
				t.Errorf("Recognize err = %v", err)
			}
		})
	}
	// falla en --list-langs (la versión sí responde)
	r := &fakeRunner{version: "tesseract 5.5.3\n", langs: "spa\n"}
	e := &Engine{Runner: &failSecond{r}}
	if err := e.CheckAvailable(ctx); !errors.Is(err, ErrOCRUnavailable) {
		t.Errorf("list-langs: %v", err)
	}
}

type failSecond struct{ *fakeRunner }

func (f *failSecond) Run(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
	if len(args) == 1 && args[0] == "--list-langs" {
		return nil, nil, errors.New("boom")
	}
	return f.fakeRunner.Run(ctx, name, args...)
}

func TestOCRVersionSeVerificaUnaVez(t *testing.T) {
	r := okRunner()
	e := &Engine{Runner: r}
	img := tmpImage(t, "a.jpeg")
	for i := 0; i < 3; i++ {
		if _, err := e.Recognize(context.Background(), img); err != nil {
			t.Fatal(err)
		}
	}
	vers := 0
	for _, c := range r.calls {
		if len(c.args) == 1 && c.args[0] == "--version" {
			vers++
		}
	}
	if vers != 1 {
		t.Errorf("--version llamado %d veces, quiero 1", vers)
	}
}

func TestRecognizeBadImage(t *testing.T) {
	e := &Engine{Runner: okRunner()}
	ctx := context.Background()
	for _, p := range []string{"foto.gif", filepath.Join(t.TempDir(), "no-existe.png"), t.TempDir() + "/x.png/"} {
		if _, err := e.Recognize(ctx, p); !errors.Is(err, ErrBadImage) {
			t.Errorf("%q: err = %v", p, err)
		}
	}
	dir := filepath.Join(t.TempDir(), "carpeta.png")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := e.Recognize(ctx, dir); !errors.Is(err, ErrBadImage) {
		t.Errorf("directorio: %v", err)
	}
}

func TestRecognizeTesseractFails(t *testing.T) {
	r := okRunner()
	r.failOn = "tesseract"
	// la verificación también usa tesseract: la dejamos pasar y fallamos al leer
	e := &Engine{Runner: &failOnRead{r}}
	_, err := e.Recognize(context.Background(), tmpImage(t, "a.png"))
	if err == nil || errors.Is(err, ErrOCRUnavailable) || !strings.Contains(err.Error(), "falló") {
		t.Errorf("err = %v", err)
	}
}

type failOnRead struct{ *fakeRunner }

func (f *failOnRead) Run(ctx context.Context, name string, args ...string) ([]byte, []byte, error) {
	if name == "tesseract" && len(args) > 1 {
		return nil, []byte("falló a propósito"), errors.New("exit status 1")
	}
	f.fakeRunner.failOn = ""
	return f.fakeRunner.Run(ctx, name, args...)
}

func TestRecognizeTimeout(t *testing.T) {
	r := okRunner()
	r.delay = 200 * time.Millisecond
	e := &Engine{Runner: r, Timeout: 20 * time.Millisecond}
	e.checked = true // la verificación no es lo que se mide aquí
	_, err := e.Recognize(context.Background(), tmpImage(t, "a.png"))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, quiero DeadlineExceeded", err)
	}
}

func TestPreprocess(t *testing.T) {
	r := okRunner()
	var warns []string
	e := &Engine{Runner: r, Preprocess: true, TempDir: t.TempDir(), Warn: func(s string) { warns = append(warns, s) }}
	if _, err := e.Recognize(context.Background(), tmpImage(t, "a.jpg")); err != nil {
		t.Fatal(err)
	}
	var mg, tess call
	for _, c := range r.calls {
		if c.name == "magick" {
			mg = c
		}
		if c.name == "tesseract" && len(c.args) > 1 {
			tess = c
		}
	}
	joined := strings.Join(mg.args, " ")
	for _, w := range []string{"-colorspace Gray", "-deskew 40%", "-lat 41x41-6%", "-auto-orient"} {
		if !strings.Contains(joined, w) {
			t.Errorf("magick sin %q: %s", w, joined)
		}
	}
	if !strings.HasSuffix(tess.args[0], "pre.png") {
		t.Errorf("tesseract no leyó el preprocesado: %s", tess.args[0])
	}
	if ents, _ := os.ReadDir(e.TempDir); len(ents) != 0 {
		t.Errorf("el temporal no se borró: %v", ents)
	}
	if len(warns) != 0 {
		t.Errorf("avisos inesperados: %v", warns)
	}
}

func TestPreprocessFallbackAndWarn(t *testing.T) {
	r := okRunner()
	r.failOn = "magick"
	var warns []string
	e := &Engine{Runner: r, Preprocess: true, TempDir: t.TempDir(), Warn: func(s string) { warns = append(warns, s) }}
	img := tmpImage(t, "a.png")
	lines, err := e.Recognize(context.Background(), img)
	if err != nil || len(lines) != 1 {
		t.Fatalf("%v %v", lines, err)
	}
	if len(warns) != 1 || !strings.Contains(warns[0], "preprocesado omitido") {
		t.Errorf("avisos: %v", warns)
	}
	// sin Warn tampoco truena
	e2 := &Engine{Runner: r, Preprocess: true, TempDir: t.TempDir()}
	if _, err := e2.Recognize(context.Background(), img); err != nil {
		t.Fatal(err)
	}
	// TempDir inexistente: el preprocesado falla y se usa la original
	e3 := &Engine{Runner: okRunner(), Preprocess: true, TempDir: filepath.Join(t.TempDir(), "no", "existe")}
	if _, err := e3.Recognize(context.Background(), img); err != nil {
		t.Fatal(err)
	}
}

func TestNew(t *testing.T) {
	e := New()
	if _, ok := e.Runner.(ExecRunner); !ok || !e.Preprocess || e.Profile.Name != "generico" {
		t.Errorf("%+v", e)
	}
}
