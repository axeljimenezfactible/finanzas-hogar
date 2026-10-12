package ocr

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeDescription(t *testing.T) {
	cases := map[string]string{
		"  Leche   Entera, ":      "LECHE ENTERA",
		"Jabón de Tocador.":       "JABON DE TOCADOR",
		"PIÑA  ÁÉÍÓÚ ÜÑ -":        "PINA AEIOU UN",
		"aceite 900 ml":           "ACEITE 900 ML",
		"SAL; ":                   "SAL",
		"":                        "",
		"arroz 1 kg.":             "ARROZ 1 KG",
		"café con leche: grande:": "CAFE CON LECHE: GRANDE",
	}
	for in, want := range cases {
		if got := NormalizeDescription(in); got != want {
			t.Errorf("%q -> %q, quiero %q", in, got, want)
		}
	}
}

func TestReadExpectedCSV(t *testing.T) {
	exp, err := ReadExpectedCSV(strings.NewReader("descripcion,cantidad,importe\nLECHE,,28.50\nJITOMATE,0.750,22.13\n\"ATUN, AGUA\",3,66.00\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(exp) != 3 || exp[0].CantidadMilli != nil || exp[0].ImporteCents != 2850 ||
		*exp[1].CantidadMilli != 750 || exp[2].Descripcion != "ATUN, AGUA" || *exp[2].CantidadMilli != 3000 {
		t.Errorf("%+v", exp)
	}
	bad := []string{
		"",
		"x,y,z\nA,,1.00\n",
		"descripcion,cantidad,importe\nA,,\n",
		"descripcion,cantidad,importe\nA,,abc\n",
		"descripcion,cantidad,importe\nA,x,1.00\n",
		"descripcion,cantidad\nA,1\n",
	}
	for _, b := range bad {
		if _, err := ReadExpectedCSV(strings.NewReader(b)); err == nil {
			t.Errorf("%q debía fallar", b)
		}
	}
}

func TestScore(t *testing.T) {
	got := []Line{
		{Descripcion: "Leche Entera 1 L.", ImporteCents: i64(2850)},
		{Descripcion: "ARROZ BLANCO", CantidadMilli: i64(2000), ImporteCents: i64(5800)},
		{Descripcion: "JITOMATE", CantidadMilli: i64(750), ImporteCents: i64(2200)}, // importe mal
		{Descripcion: "HUEVO", ImporteCents: i64(4590)},                             // cantidad: esperada 12000
		{Descripcion: "TOTAL", ImporteCents: i64(100), EsResumen: true},
		{Descripcion: "ENCABEZADO"},                                // sin importe
		{Descripcion: "SOBRANTE", ImporteCents: i64(99)},           // extra
		{Descripcion: "LECHE ENTERA 1 L", ImporteCents: i64(2850)}, // duplicado: no cuenta dos veces
	}
	exp := []Expected{
		{"LECHE ENTERA 1 L", nil, 2850},
		{"ARROZ BLANCO", i64(2000), 5800},
		{"JITOMATE", i64(750), 2213},
		{"HUEVO", i64(12000), 4590},
		{"ATUN", nil, 6600}, // no salió
	}
	r := Score("t", got, exp)
	if r.Expected != 5 || r.Correct != 2 || len(r.Misses) != 3 {
		t.Errorf("%+v", r)
	}
	if r.Extra != 4 { // JITOMATE mal, HUEVO mal, SOBRANTE y el duplicado de LECHE
		t.Errorf("extra = %d", r.Extra)
	}
}

type fakeRec struct {
	lines map[string][]Line
	err   error
}

func (f fakeRec) Recognize(_ context.Context, img string) ([]Line, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.lines[filepath.Base(img)], nil
}

func writeTickets(t *testing.T, n int, csvs int) (string, string) {
	t.Helper()
	td, ed := t.TempDir(), t.TempDir()
	for i := 1; i <= n; i++ {
		name := "0" + string(rune('0'+i))
		ext := ".jpg"
		if i%2 == 0 {
			ext = ".PNG"
		}
		must(t, os.WriteFile(filepath.Join(td, name+ext), []byte("x"), 0o600))
		if i <= csvs {
			must(t, os.WriteFile(filepath.Join(ed, name+".csv"), []byte("descripcion,cantidad,importe\nLECHE,,10.00\nARROZ,2,20.00\n"), 0o600))
		}
	}
	must(t, os.WriteFile(filepath.Join(td, "LEEME.md"), []byte("no cuenta"), 0o600))
	must(t, os.Mkdir(filepath.Join(td, "sub"), 0o700))
	return td, ed
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func perfect() map[string][]Line {
	m := map[string][]Line{}
	for _, n := range []string{"01.jpg", "02.PNG", "03.jpg", "04.PNG", "05.jpg"} {
		m[n] = []Line{
			{Descripcion: "leche", ImporteCents: i64(1000)},
			{Descripcion: "ARROZ", CantidadMilli: i64(2000), ImporteCents: i64(2000)},
		}
	}
	return m
}

func TestMeasureSinMedir(t *testing.T) {
	ctx := context.Background()
	// carpeta inexistente
	rep, err := Measure(ctx, fakeRec{}, MeasureOptions{TicketsDir: filepath.Join(t.TempDir(), "no"), ExpectedDir: t.TempDir()})
	if err != nil || rep.Measured || rep.Passes || !strings.Contains(rep.Reason, "faltan las fotos de Axel") {
		t.Errorf("sin carpeta: %+v %v", rep, err)
	}
	if !strings.Contains(rep.Format(), "S5-06 sin medir: hay 0 imágenes") {
		t.Errorf("Format: %q", rep.Format())
	}
	// 4 de 5
	td, ed := writeTickets(t, 4, 4)
	rep, _ = Measure(ctx, fakeRec{}, MeasureOptions{TicketsDir: td, ExpectedDir: ed})
	if rep.Measured || !strings.Contains(rep.Reason, "hay 4 imágenes") {
		t.Errorf("4 fotos: %+v", rep)
	}
	// 5 fotos pero falta una transcripción
	td, ed = writeTickets(t, 5, 4)
	rep, _ = Measure(ctx, fakeRec{}, MeasureOptions{TicketsDir: td, ExpectedDir: ed})
	if rep.Measured || !strings.Contains(rep.Reason, "falta la transcripción") {
		t.Errorf("sin csv: %+v", rep)
	}
	// sintético sin imágenes
	rep, _ = Measure(ctx, fakeRec{}, MeasureOptions{TicketsDir: t.TempDir(), ExpectedDir: ed, Synthetic: true})
	if rep.Measured || !rep.Synthetic || !strings.Contains(rep.Format(), "NO cuenta para S5-06") {
		t.Errorf("sintético vacío: %+v", rep)
	}
	// transcripciones sin renglones
	td, ed = writeTickets(t, 5, 0)
	for _, n := range []string{"01", "02", "03", "04", "05"} {
		must(t, os.WriteFile(filepath.Join(ed, n+".csv"), []byte("descripcion,cantidad,importe\n"), 0o600))
	}
	rep, _ = Measure(ctx, fakeRec{lines: perfect()}, MeasureOptions{TicketsDir: td, ExpectedDir: ed})
	if rep.Measured || !strings.Contains(rep.Reason, "no tienen renglones") {
		t.Errorf("csv vacío: %+v", rep)
	}
}

func TestMeasureErrores(t *testing.T) {
	ctx := context.Background()
	td, ed := writeTickets(t, 5, 5)
	if _, err := Measure(ctx, fakeRec{err: errors.New("boom")}, MeasureOptions{TicketsDir: td, ExpectedDir: ed}); err == nil {
		t.Error("error del OCR debía propagarse")
	}
	must(t, os.WriteFile(filepath.Join(ed, "03.csv"), []byte("basura"), 0o600))
	if _, err := Measure(ctx, fakeRec{}, MeasureOptions{TicketsDir: td, ExpectedDir: ed}); err == nil {
		t.Error("csv inválido debía fallar")
	}
	f := filepath.Join(t.TempDir(), "archivo")
	must(t, os.WriteFile(f, []byte("x"), 0o600))
	if _, err := Measure(ctx, fakeRec{}, MeasureOptions{TicketsDir: filepath.Join(f, "x"), ExpectedDir: ed}); err == nil {
		t.Error("ReadDir con error distinto de «no existe» debía fallar")
	}
}

func TestMeasurePorcentajes(t *testing.T) {
	ctx := context.Background()
	td, ed := writeTickets(t, 5, 5)
	// 10 esperados, todos correctos
	rep, err := Measure(ctx, fakeRec{lines: perfect()}, MeasureOptions{TicketsDir: td, ExpectedDir: ed})
	if err != nil || !rep.Measured || !rep.Passes || rep.PctTenths != 1000 || rep.Correct != 10 {
		t.Fatalf("perfecto: %+v %v", rep, err)
	}
	if !strings.Contains(rep.Format(), "TOTAL 10/10 = 100.0 %") || !strings.Contains(rep.Format(), "cumple (>= 80 %)") {
		t.Errorf("Format: %s", rep.Format())
	}
	// 8 de 10 = 80 % cumple justo
	m := perfect()
	m["01.jpg"] = nil
	rep, _ = Measure(ctx, fakeRec{lines: m}, MeasureOptions{TicketsDir: td, ExpectedDir: ed})
	if rep.Correct != 8 || !rep.Passes || rep.PctTenths != 800 {
		t.Errorf("80 %%: %+v", rep)
	}
	// 7 de 10 no cumple
	m["02.PNG"] = m["02.PNG"][:1]
	rep, _ = Measure(ctx, fakeRec{lines: m}, MeasureOptions{TicketsDir: td, ExpectedDir: ed})
	if rep.Correct != 7 || rep.Passes {
		t.Errorf("70 %%: %+v", rep)
	}
	out := rep.Format()
	if !strings.Contains(out, "falla: LECHE") || !strings.Contains(out, "NO cumple (< 80 %)") {
		t.Errorf("Format: %s", out)
	}
	// sintético mide aunque no sean 5 y lo declara
	td2, ed2 := writeTickets(t, 2, 2)
	rep, _ = Measure(ctx, fakeRec{lines: perfect()}, MeasureOptions{TicketsDir: td2, ExpectedDir: ed2, Synthetic: true})
	if !rep.Measured || !strings.HasPrefix(rep.Format(), "SINTETICO") {
		t.Errorf("sintético: %+v", rep)
	}
	// redondeo HALF_UP de la décima: 2/3 = 66.7
	rep = Report{Measured: true, Expected: 3, Correct: 2}
	rep.PctTenths = (rep.Correct*2000 + rep.Expected) / (2 * rep.Expected)
	if rep.PctTenths != 667 {
		t.Errorf("décimas = %d", rep.PctTenths)
	}
}
