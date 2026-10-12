package ocr

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Runner ejecuta un programa SIN shell. Es la única costura con el sistema:
// en las pruebas se sustituye por un ejecutor simulado.
type Runner interface {
	Run(ctx context.Context, name string, args ...string) (stdout, stderr []byte, err error)
}

// Engine corre el OCR. Sirve para varias imágenes y es seguro entre hilos.
type Engine struct {
	Runner     Runner
	Profile    Profile
	Timeout    time.Duration // 0 = DefaultTimeoutSeconds (PROVISIONAL)
	Preprocess bool          // escala de grises, enderezado y umbral con ImageMagick
	TempDir    string        // dónde dejar el preprocesado; "" = os.TempDir()
	// Warn recibe avisos no fatales (p. ej. el preprocesado falló y se usó la
	// imagen original). Puede ser nil.
	Warn func(string)

	mu      sync.Mutex
	checked bool
}

// New devuelve un motor que usa los programas locales.
func New() *Engine {
	return &Engine{Runner: ExecRunner{}, Profile: ProfileGeneric, Preprocess: true}
}

var versionRe = regexp.MustCompile(`(?m)^tesseract\s+v?(\d+\.\d+\.\d+)`)

// CheckAvailable verifica una sola vez tesseract 5.5.3 y el idioma spa.
// Si algo falla devuelve ErrOCRUnavailable (envuelto con el motivo) y vuelve a
// intentarlo en la siguiente llamada.
func (e *Engine) CheckAvailable(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.checked {
		return nil
	}
	out, errOut, err := e.Runner.Run(ctx, "tesseract", "--version")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrOCRUnavailable, err)
	}
	m := versionRe.FindSubmatch(append(append([]byte{}, out...), errOut...))
	if m == nil || string(m[1]) != RequiredVersion {
		got := "desconocida"
		if m != nil {
			got = string(m[1])
		}
		return fmt.Errorf("%w: versión %s", ErrOCRUnavailable, got)
	}
	out, errOut, err = e.Runner.Run(ctx, "tesseract", "--list-langs")
	if err != nil {
		return fmt.Errorf("%w: %v", ErrOCRUnavailable, err)
	}
	if !hasLang(append(append([]byte{}, out...), errOut...), Lang) {
		return fmt.Errorf("%w: falta el idioma %s", ErrOCRUnavailable, Lang)
	}
	e.checked = true
	return nil
}

func hasLang(out []byte, lang string) bool {
	for _, ln := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(ln) == lang {
			return true
		}
	}
	return false
}

func (e *Engine) warn(msg string) {
	if e.Warn != nil {
		e.Warn(msg)
	}
}

// tesseractArgs es la configuración documentada (ver doc.go).
func tesseractArgs(img string) []string {
	return []string{img, "stdout", "-l", Lang, "--psm", "6", "--oem", "1",
		"-c", "preserve_interword_spaces=1", "tsv"}
}

// Recognize lee una imagen y devuelve el BORRADOR de renglones. No guarda
// nada: la persona confirma o corrige y solo entonces PricesFromOCR entrega
// renglones.
func (e *Engine) Recognize(ctx context.Context, imagen string) ([]Line, error) {
	img, err := validImage(imagen)
	if err != nil {
		return nil, err
	}
	if err := e.CheckAvailable(ctx); err != nil {
		return nil, err
	}
	t := e.Timeout
	if t <= 0 {
		t = DefaultTimeoutSeconds * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, t)
	defer cancel()

	in := img
	if e.Preprocess {
		out, cleanup, perr := e.preprocess(ctx, img)
		if perr != nil {
			e.warn("preprocesado omitido: " + perr.Error())
		} else {
			defer cleanup()
			in = out
		}
	}
	stdout, stderr, err := e.Runner.Run(ctx, "tesseract", tesseractArgs(in)...)
	if err != nil {
		return nil, fmt.Errorf("ocr: tesseract falló: %w: %s", err, truncate(strings.TrimSpace(string(stderr)), 200))
	}
	return ParseTSV(bytes.NewReader(stdout), e.Profile)
}

func validImage(path string) (string, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg", ".png":
	default:
		return "", fmt.Errorf("%w: %q (solo .jpg, .jpeg o .png)", ErrBadImage, filepath.Base(path))
	}
	abs, err := filepath.Abs(path) // ruta absoluta: nunca empieza con "-"
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrBadImage, err)
	}
	st, err := os.Stat(abs)
	if err != nil || st.IsDir() {
		return "", fmt.Errorf("%w: %q no existe", ErrBadImage, filepath.Base(path))
	}
	return abs, nil
}

// preprocessArgs: orienta según EXIF, escala de grises, amplía a 1600 px de
// ancho si es más chica (Tesseract rinde mejor con letra de unos 30 px),
// endereza (deskew 40 %), normaliza el contraste y aplica un umbral LOCAL
// (-lat) que sobrevive a sombras de la foto. Salida en PNG sin pérdida.
func preprocessArgs(in, out string) []string {
	return []string{in, "-auto-orient", "-colorspace", "Gray", "-resize", "1600x<",
		"-deskew", "40%", "+repage", "-normalize", "-lat", "41x41-6%", "-type", "Bilevel", out}
}

func (e *Engine) preprocess(ctx context.Context, img string) (string, func(), error) {
	dir, err := os.MkdirTemp(e.TempDir, "ocr-")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	out := filepath.Join(dir, "pre.png")
	if _, stderr, err := e.Runner.Run(ctx, "magick", preprocessArgs(img, out)...); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("magick: %w: %s", err, truncate(strings.TrimSpace(string(stderr)), 200))
	}
	return out, cleanup, nil
}
