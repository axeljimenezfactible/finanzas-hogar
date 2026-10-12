package profeco

import (
	"strings"
	"testing"
)

func TestParsePriceCents(t *testing.T) {
	ok := map[string]int64{
		"28": 2800, "34.5": 3450, "17.50": 1750, "$1,234.50": 123450, " 0.01 ": 1,
		"1,000": 100000, "99999999.99": 9999999999, "5.": -1,
	}
	for in, want := range ok {
		got, err := ParsePriceCents(in)
		if want == -1 {
			if err == nil {
				t.Errorf("%q debía rechazarse", in)
			}
			continue
		}
		if err != nil || got != want {
			t.Errorf("%q: got %d, %v; want %d", in, got, err, want)
		}
	}
	for _, bad := range []string{"", "  ", "abc", "-5", "0", "0.00", "1.234", "1,23", "1,,000", ",500", "$", "12a", "1e3", "9999999999999999999", "1.2.3"} {
		if v, err := ParsePriceCents(bad); err == nil {
			t.Errorf("%q debía rechazarse y dio %d", bad, v)
		}
	}
}

func TestParseDate(t *testing.T) {
	ok := map[string][2]string{
		"2026/01/02":          {"2026-01-02", "2006/01/02"},
		"2026-08-31":          {"2026-08-31", "2006-01-02"},
		"31/08/2026":          {"2026-08-31", "02/01/2006"},
		"2026-08-31 10:11:12": {"2026-08-31", "2006-01-02 15:04:05"},
		"20260831":            {"2026-08-31", "20060102"},
		"2026/1/2":            {"2026-01-02", "2006/1/2"},
	}
	for in, want := range ok {
		iso, layout, err := ParseDate(in)
		if err != nil || iso != want[0] || layout != want[1] {
			t.Errorf("%q: got %q %q %v; want %v", in, iso, layout, err, want)
		}
	}
	for _, bad := range []string{"", "ayer", "2026/13/40", "2026/02/30", "13/13/2026"} {
		if _, _, err := ParseDate(bad); err == nil {
			t.Errorf("%q debía rechazarse", bad)
		}
	}
}

func TestCleanCoord(t *testing.T) {
	for in, want := range map[string]string{"21.854349": "21.854349", "-102.294259": "-102.294259", "": "", "abc": "", "1.": "", "1234.5": "", "19": "19"} {
		if got := cleanCoord(in); got != want {
			t.Errorf("cleanCoord(%q) = %q; want %q", in, got, want)
		}
	}
}

func TestFold(t *testing.T) {
	same := [][]string{
		{"Wal-mart", "WALMART", "walmart"},
		{"Sam's Club", "SAMS CLUB", "Sam’s  Club"},
		{"Álvaro Obregón", "alvaro obregon"},
		{"Gustavo a. Madero", "GUSTAVO A MADERO"},
	}
	for _, g := range same {
		for _, s := range g[1:] {
			if Fold(g[0]) != Fold(s) {
				t.Errorf("Fold(%q)=%q distinto de Fold(%q)=%q", g[0], Fold(g[0]), s, Fold(s))
			}
		}
	}
	if Fold("Wal-mart Express") == Fold("Wal-mart") {
		t.Error("Express y la cadena base deben distinguirse")
	}
	// Bytes Windows-1252 «Qui\xe9n» se leen como «Quién».
	if got := Fold("Qui\xe9n"); got != "quien" {
		t.Errorf("Fold cp1252 = %q", got)
	}
}

func TestDecodeTextWindows1252(t *testing.T) {
	// 0x80 = €, 0x93/0x94 = comillas, 0xF1 = ñ, 0x81 queda como punto C1.
	got := decodeText("\x80 \x93a\x94 \xf1 \x81")
	if !strings.HasPrefix(got, "€ “a” ñ ") {
		t.Errorf("decodeText = %q", got)
	}
	if decodeText("ñandú") != "ñandú" {
		t.Error("UTF-8 válido no debe cambiar")
	}
}

func TestIsFocusChain(t *testing.T) {
	for in, want := range map[string]bool{"Wal-mart": true, "Wal-mart Express": true, "Bodega Aurrera": true, "Bodega Aurrera Express": true, "Sam's Club": true, "Chedraui": false, "S Mart": false, "Soriana": false} {
		if IsFocusChain(Fold(in)) != want {
			t.Errorf("IsFocusChain(%q) != %v", in, want)
		}
	}
}

// Lista cerrada de motivos del precio (§11 paso 3, H-r3-03, Q2-04 = A).
func TestParsePriceMotivos(t *testing.T) {
	for in, want := range map[string]string{
		"": "precio vacío", " $ ": "precio vacío", "1,02.00": "precio con comas mal puestas", ",500": "precio con comas mal puestas",
		"1.0,0": "precio con comas mal puestas", "1,,000": "precio con comas mal puestas", "24.505": "precio con más de dos decimales",
		".123": "precio con más de dos decimales", "abc": "precio ilegible", "12.": "precio ilegible", ".5": "precio ilegible",
		"-5": "precio ilegible", "1 021.00": "precio ilegible", "9999999999999999": "precio ilegible", "0.00": "precio no positivo", "$0": "precio no positivo",
	} {
		if _, err := ParsePriceCents(in); err == nil || err.Error() != want {
			t.Errorf("%q: %v; se esperaba %q", in, err, want)
		}
	}
	if c, err := ParsePriceCents("$ 1,021.00"); err != nil || c != 102100 {
		t.Errorf("1,021.00 = %d, %v", c, err)
	}
}
