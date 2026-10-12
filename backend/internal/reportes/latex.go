package reportes

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"
)

//go:embed plantillas/*.tex
var plantillasFS embed.FS

// Plantillas son las de backend/internal/reportes/plantillas (embebidas en el binario).
func Plantillas() fs.FS { return plantillasFS }

var (
	tplMu       sync.RWMutex
	tplOverride fs.FS
)

// SetTemplates cambia las plantillas de todos los compiladores sin Templates propio (solo pruebas: plantilla rota
// de S7-10); devuelve la función que deja las embebidas.
func SetTemplates(f fs.FS) (restore func()) {
	tplMu.Lock()
	tplOverride = f
	tplMu.Unlock()
	return func() { tplMu.Lock(); tplOverride = nil; tplMu.Unlock() }
}

// LatexTimeout es el tiempo límite de una compilación. PROVISIONAL (K31, puerta 4): 120 s; una compilación normal
// tarda 1–3 s en la máquina de casa, así que 120 s solo corta un LaTeX colgado sin cortar uno lento.
const LatexTimeout = 120 * time.Second

// Compiler pinta una plantilla y la compila con lualatex en un directorio temporal propio.
type Compiler struct {
	Engine    string        // "lualatex" (decisión del plan; S7-01)
	Timeout   time.Duration // LatexTimeout si es 0
	TmpDir    string        // directorio de temporales (datos/tmp)
	Templates fs.FS         // nil = las embebidas; las pruebas ponen una rota para S7-10
	KeepDir   string        // si no está vacío, deja ahí el .tex, el .log y el .pdf (evidencia de S7-01)
}

// ErrLatex marca que LaTeX falló: no hay PDF y nada se guarda ni se envía (S7-10).
var ErrLatex = errors.New("reportes: LaTeX falló")

// Render pinta la plantilla (mensual o estado-cuenta) con el documento.
func (c *Compiler) Render(name string, d *Doc) ([]byte, error) {
	fsys := c.Templates
	if fsys == nil {
		tplMu.RLock()
		fsys = tplOverride
		tplMu.RUnlock()
	}
	if fsys == nil {
		fsys = plantillasFS
	}
	t, err := template.New(name).Delims("<<", ">>").Funcs(template.FuncMap{"tex": Tex}).
		ParseFS(fsys, "plantillas/comun.tex", "plantillas/"+name+".tex")
	if err != nil {
		return nil, fmt.Errorf("%w: plantilla %s: %v", ErrLatex, name, err)
	}
	var b bytes.Buffer
	if err := t.ExecuteTemplate(&b, name+".tex", d); err != nil {
		return nil, fmt.Errorf("%w: plantilla %s: %v", ErrLatex, name, err)
	}
	return b.Bytes(), nil
}

// Compile devuelve el PDF. Falla si lualatex sale con error, si el registro trae líneas «!» o «undefined», o si el
// PDF no empieza con %PDF-; el texto del error es el de LaTeX (sin datos personales: solo la línea del error).
func (c *Compiler) Compile(ctx context.Context, name string, d *Doc) ([]byte, error) {
	src, err := c.Render(name, d)
	if err != nil {
		return nil, err
	}
	if c.TmpDir != "" {
		if err := os.MkdirAll(c.TmpDir, 0o700); err != nil {
			return nil, err
		}
	}
	dir, err := os.MkdirTemp(c.TmpDir, "latex-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, name+".tex"), src, 0o600); err != nil {
		return nil, err
	}
	to := c.Timeout
	if to <= 0 {
		to = LatexTimeout
	}
	engine := c.Engine
	if engine == "" {
		engine = "lualatex"
	}
	cctx, cancel := context.WithTimeout(ctx, to)
	defer cancel()
	cmd := exec.CommandContext(cctx, engine, "-interaction=nonstopmode", "-halt-on-error", "-no-shell-escape", name+".tex")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "TEXMFOUTPUT="+dir, "openout_any=p", "openin_any=p")
	out, runErr := cmd.CombinedOutput()
	logb, _ := os.ReadFile(filepath.Join(dir, name+".log"))
	pdf, _ := os.ReadFile(filepath.Join(dir, name+".pdf"))
	if c.KeepDir != "" {
		_ = os.MkdirAll(c.KeepDir, 0o700)
		_ = os.WriteFile(filepath.Join(c.KeepDir, name+".tex"), src, 0o600)
		_ = os.WriteFile(filepath.Join(c.KeepDir, name+".log"), logb, 0o600)
		if len(pdf) > 0 {
			_ = os.WriteFile(filepath.Join(c.KeepDir, name+".pdf"), pdf, 0o600)
		}
	}
	if cctx.Err() != nil {
		return nil, fmt.Errorf("%w: se pasó del tiempo límite (%s)", ErrLatex, to)
	}
	if runErr != nil {
		return nil, fmt.Errorf("%w: %s", ErrLatex, firstError(logb, out))
	}
	if bad := firstError(logb, nil); bad != "" {
		return nil, fmt.Errorf("%w: %s", ErrLatex, bad)
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		return nil, fmt.Errorf("%w: no se produjo el PDF", ErrLatex)
	}
	return pdf, nil
}

// firstError da la primera línea de error del registro («!» o «undefined»); "" si no hay.
func firstError(logb, out []byte) string {
	for _, l := range strings.Split(string(logb), "\n") {
		if strings.HasPrefix(l, "!") || strings.Contains(l, "undefined") {
			return strings.TrimSpace(l)
		}
	}
	if len(out) > 0 {
		s := strings.TrimSpace(string(out))
		if i := strings.LastIndex(s, "\n"); i >= 0 {
			s = s[i+1:]
		}
		if len(s) > 200 {
			s = s[:200]
		}
		return s
	}
	return ""
}
