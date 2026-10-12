package ocr

import (
	"strings"
	"testing"
)

func i64(v int64) *int64 { return &v }

func TestParseLineCases(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		p       Profile
		desc    string
		qty     *int64
		unit    string
		amount  *int64
		resumen bool
	}{
		{"solo descripcion e importe", "LECHE ENTERA 1 L   28.50", ProfileGeneric, "LECHE ENTERA 1 L", nil, "", i64(2850), false},
		{"cantidad entera al inicio", "2 ARROZ BLANCO 1 KG  58.00", ProfileGeneric, "ARROZ BLANCO 1 KG", i64(2000), "", i64(5800), false},
		{"cantidad decimal con unidad", "0.750 KG JITOMATE 22.13", ProfileGeneric, "JITOMATE", i64(750), "KG", i64(2213), false},
		{"coma decimal en cantidad", "0,500 kg QUESO 80.00", ProfileGeneric, "QUESO", i64(500), "KG", i64(8000), false},
		{"signo peso, miles y letra de iva", "CAFE SOLUBLE 200 G  $1,234.50 A", ProfileGeneric, "CAFE SOLUBLE 200 G", nil, "", i64(123450), false},
		{"importe negativo (descuento)", "DESCUENTO PROMO -5.00", ProfileGeneric, "DESCUENTO PROMO", nil, "", i64(-500), false},
		{"sin importe", "TIENDA SINTETICA DE PRUEBA", ProfileGeneric, "TIENDA SINTETICA DE PRUEBA", nil, "", nil, false},
		{"codigo de barras se conserva en generico", "7501234567890 REFRESCO COLA 2L 32.00 B", ProfileGeneric, "7501234567890 REFRESCO COLA 2L", nil, "", i64(3200), false},
		{"codigo de barras se quita con perfil", "7501234567890 REFRESCO COLA 2L 32.00 B", ProfileWithCode, "REFRESCO COLA 2L", nil, "", i64(3200), false},
		{"o por cero en el importe", "PAN BLANCO 3O.OO", ProfileGeneric, "PAN BLANCO", nil, "", i64(3000), false},
		{"l por uno en el importe", "SAL l2.5O", ProfileGeneric, "SAL", nil, "", i64(1250), false},
		{"numero sin dos decimales no es importe", "HUEVO 12 PZA 45.9", ProfileGeneric, "HUEVO 12 PZA 45.9", nil, "", nil, false},
		{"tres decimales no es importe", "QUESO 0.750", ProfileGeneric, "QUESO 0.750", nil, "", nil, false},
		{"cantidad sin texto despues no es cantidad", "2 15.00", ProfileGeneric, "2", nil, "", i64(1500), false},
		{"total es resumen", "TOTAL   $1,457.05", ProfileGeneric, "TOTAL", nil, "", i64(145705), true},
		{"subtotal es resumen", "Sub total 251.53", ProfileGeneric, "Sub total", nil, "", i64(25153), true},
		{"cambio es resumen", "CAMBIO 42.95", ProfileGeneric, "CAMBIO", nil, "", i64(4295), true},
		{"iva con puntos es resumen", "I.V.A. 16.00", ProfileGeneric, "I.V.A", nil, "", i64(1600), true},
		{"producto con total en medio no es resumen", "JUGO TOTAL 1L 22.00", ProfileGeneric, "JUGO TOTAL 1L", nil, "", i64(2200), false},
		{"linea vacia", "   ", ProfileGeneric, "", nil, "", nil, false},
		{"importe demasiado largo", "X 1234567890123.00", ProfileGeneric, "X 1234567890123.00", nil, "", nil, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := ParseLine(c.in, c.p)
			if l.Descripcion != c.desc {
				t.Errorf("descripcion = %q, quiero %q", l.Descripcion, c.desc)
			}
			if !eqQty(l.CantidadMilli, c.qty) {
				t.Errorf("cantidad = %v, quiero %v", ptr(l.CantidadMilli), ptr(c.qty))
			}
			if l.Unidad != c.unit {
				t.Errorf("unidad = %q, quiero %q", l.Unidad, c.unit)
			}
			if !eqQty(l.ImporteCents, c.amount) {
				t.Errorf("importe = %v, quiero %v", ptr(l.ImporteCents), ptr(c.amount))
			}
			if l.EsResumen != c.resumen {
				t.Errorf("resumen = %v, quiero %v", l.EsResumen, c.resumen)
			}
		})
	}
}

func ptr(p *int64) any {
	if p == nil {
		return nil
	}
	return *p
}

func TestParseFixed(t *testing.T) {
	good := map[string]int64{"1,234.50": 123450, "$45.00": 4500, "-5.00": -500, "7": 700, "0.5": 50}
	for in, want := range good {
		if got, ok := parseFixed(in, 2); !ok || got != want {
			t.Errorf("parseFixed(%q) = %d,%v quiero %d", in, got, ok, want)
		}
	}
	for _, in := range []string{"", "abc", "1.234", ".50", "1.2x", "9999999999999.00"} {
		if _, ok := parseFixed(in, 2); ok {
			t.Errorf("parseFixed(%q) debía fallar", in)
		}
	}
	if v, ok := parseFixed("0.750", 3); !ok || v != 750 {
		t.Errorf("milésimas: %d %v", v, ok)
	}
}

func TestParseQty(t *testing.T) {
	if _, ok := parseQty("7501234567890"); ok {
		t.Error("un código de barras no es cantidad")
	}
	if v, ok := parseQty("12"); !ok || v != 12000 {
		t.Errorf("12 -> %d %v", v, ok)
	}
}

const tsvHead = "level\tpage_num\tblock_num\tpar_num\tline_num\tword_num\tleft\ttop\twidth\theight\tconf\ttext\n"

func tsvRow(level, block, par, line, word, conf, text string) string {
	return strings.Join([]string{level, "1", block, par, line, word, "0", "0", "10", "10", conf, text}, "\t") + "\n"
}

func TestParseTSV(t *testing.T) {
	tsv := tsvHead +
		tsvRow("1", "0", "0", "0", "0", "-1", "") +
		tsvRow("4", "1", "1", "1", "0", "-1", "") +
		tsvRow("5", "1", "1", "1", "1", "96.5", "2") +
		tsvRow("5", "1", "1", "1", "2", "90.4", "ARROZ") +
		tsvRow("5", "1", "1", "1", "3", "93.0", "58.00") +
		tsvRow("5", "1", "1", "2", "1", "80", "TOTAL") +
		tsvRow("5", "1", "1", "2", "2", "-1", "ignorada") + // conf -1
		tsvRow("5", "1", "1", "2", "3", "81", " ") + // texto vacío
		tsvRow("5", "1", "1", "2", "4", "70", "58.00") +
		"basura corta\n\n"
	lines, err := ParseTSV(strings.NewReader(tsv), ProfileGeneric)
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 {
		t.Fatalf("renglones = %d, quiero 2: %+v", len(lines), lines)
	}
	a := lines[0]
	if a.Descripcion != "ARROZ" || a.CantidadMilli == nil || *a.CantidadMilli != 2000 || *a.ImporteCents != 5800 {
		t.Errorf("renglón 1: %+v", a)
	}
	if a.Confianza != 93 { // (96.5+90.4+93.0)/3 = 93.3
		t.Errorf("confianza = %d, quiero 93", a.Confianza)
	}
	b := lines[1]
	if !b.EsResumen || b.Confianza != 75 { // (80+70)/2 = 75
		t.Errorf("renglón 2: %+v", b)
	}
}

func TestParseTSVErrors(t *testing.T) {
	if _, err := ParseTSV(strings.NewReader("no es tsv\n"), ProfileGeneric); err == nil {
		t.Error("encabezado inválido debía fallar")
	}
	lines, err := ParseTSV(strings.NewReader(""), ProfileGeneric)
	if err != nil || len(lines) != 0 {
		t.Errorf("vacío: %v %v", lines, err)
	}
	long := tsvHead + tsvRow("5", "1", "1", "1", "1", "90", strings.Repeat("A", 5*1024*1024))
	if _, err := ParseTSV(strings.NewReader(long), ProfileGeneric); err == nil {
		t.Error("línea gigante debía fallar")
	}
}

func TestTruncate(t *testing.T) {
	if truncate("abc", 5) != "abc" || truncate("abcdef", 3) != "abc" {
		t.Error("truncate")
	}
}
