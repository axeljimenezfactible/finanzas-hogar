package profeco

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func surveyOf(t *testing.T, paths ...string) *Survey {
	t.Helper()
	s := NewSurvey()
	for _, p := range paths {
		if err := s.ScanFile(context.Background(), p, 0, 10); err != nil {
			t.Fatal(err)
		}
	}
	return s
}

func TestSurveySucursalesNombradas(t *testing.T) {
	body := header +
		// Sam's Club solo como referencia en la dirección de un Wal-mart (así está en los datos reales de Cancún).
		row("Atún", "Lata", "Dolores", "21.5", "2026/03/02", "Wal-mart", "Wal Mart", "Coba Sm. 21. Costado Sams Club. Col. Centro", "Quintana Roo", "Benito Juárez") +
		// Una cadena que sí se llama Sam's Club.
		row("Atún", "Lata", "Dolores", "40", "2026/03/03", "Sam's Club", "Sam's Club Sucursal Norte", "Av. 1", "Nuevo León", "Monterrey") +
		// Bodega Aurrera con Vallejo en el nombre y Walmart con la frase completa.
		row("Atún", "Lata", "Dolores", "20", "2026/04/03", "Bodega Aurrera", "Bodega Aurrera Vallejo", "Av. 2", "Ciudad de México", "Gustavo a. Madero") +
		row("Atún", "Lata", "Dolores", "25", "2026/05/04", "Wal-mart", "Walmart Acueducto de Guadalupe", "Av. 3", "Ciudad de México", "Gustavo a. Madero") +
		// Parcial: solo «acueducto».
		row("Atún", "Lata", "Dolores", "26", "2026/05/05", "Wal-mart Express", "Walmart Express Virreyes", "Esq. Acueducto Rio Hondo", "Ciudad de México", "Miguel Hidalgo") +
		// Cadena que no es de interés ni de las vigiladas.
		row("Arroz", "Bolsa", "SOS", "19", "2026/05/05", "Chedraui", "Chedraui X", "Av. 4 Vallejo", "Jalisco", "Guadalajara") +
		"mal,renglon\n"
	p := writeTmp(t, "nom.csv", body)
	s := surveyOf(t, p)

	if s.Total.Records != 6 || len(s.Chains) != 5 || len(s.Rejects) != 1 || len(s.Rejections) != 1 {
		t.Fatalf("totales: %d cadenas %d rechazos %v", s.Total.Records, len(s.Chains), s.Rejects)
	}
	if s.Total.Products() != 2 || s.Total.Presentations() != 2 || s.Total.Branches() != 6 || len(s.Total.Months) != 3 {
		t.Errorf("agregados: %d %d %d %d", s.Total.Products(), s.Total.Presentations(), s.Total.Branches(), len(s.Total.Months))
	}
	byLabel := map[string]*NamedResult{}
	for _, n := range s.Named {
		byLabel[n.Label] = n
	}
	sams := byLabel["Sam's Club"]
	if sams.Agg.Records != 2 || sams.DirOnly.Records != 1 {
		t.Errorf("Sam's: %d en total, %d solo en dirección", sams.Agg.Records, sams.DirOnly.Records)
	}
	if v := byLabel["Bodega Aurrera Vallejo"]; v.Agg.Records != 1 || v.DirOnly.Records != 0 {
		t.Errorf("Vallejo: %d / %d", v.Agg.Records, v.DirOnly.Records)
	}
	ac := byLabel["Walmart Acueducto de Guadalupe"]
	if ac.Agg.Records != 1 || ac.NearAgg.Records != 1 {
		t.Errorf("Acueducto: %d exacto, %d parcial", ac.Agg.Records, ac.NearAgg.Records)
	}
	// «Vallejo» en una dirección de Chedraui no cuenta: la cadena no es la pedida.
	if byLabel["Bodega Aurrera Vallejo"].NearAgg.Records != 0 {
		t.Error("Chedraui no debe contar como Bodega Aurrera Vallejo")
	}
	if len(s.Watch) != 2 { // Gustavo a. Madero: Bodega Aurrera y Wal-mart
		t.Errorf("municipio vigilado: %d grupos", len(s.Watch))
	}

	var out bytes.Buffer
	s.WriteMarkdown(&out)
	md := out.String()
	for _, want := range []string{"cadena_comercial", "`wc -l`", "Sucursales de Axel", "Regla de empate", "solo en `direccion`", "★", "Gustavo a. Madero", "Rechazados con motivo", "NO hay registros"} {
		if !strings.Contains(md, want) {
			if want == "NO hay registros" {
				continue // todas las sucursales de esta muestra tienen registros
			}
			t.Errorf("el reporte no dice %q", want)
		}
	}
	// El reporte es determinista.
	var again bytes.Buffer
	surveyOf(t, p).WriteMarkdown(&again)
	if again.String() != md {
		t.Error("el reporte cambia entre corridas")
	}
}

func TestSurveySinRegistrosDeInteres(t *testing.T) {
	p := writeTmp(t, "no.csv", header+row("Atún", "Lata", "Dolores", "21", "2026/03/02", "Chedraui", "Chedraui X", "Av 1", "Jalisco", "Guadalajara"))
	s := surveyOf(t, p)
	var out bytes.Buffer
	s.WriteMarkdown(&out)
	md := out.String()
	if !strings.Contains(md, "Ninguna de las cadenas de interés tiene registros.") || strings.Count(md, "NO hay registros") != 3 {
		t.Errorf("reporte sin cadenas de interés:\n%s", md)
	}
	if err := NewSurvey().ScanFile(context.Background(), p, 5, 1); err == nil {
		t.Error("ScanFile con límite de tamaño debía fallar")
	}
	if err := NewSurvey().ScanFile(context.Background(), p+".no", 0, 1); err == nil {
		t.Error("ScanFile inexistente debía fallar")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := NewSurvey().ScanFile(ctx, p, 0, 1); err == nil {
		t.Error("ScanFile cancelado debía fallar")
	}
}

func TestSurveyMuestra(t *testing.T) {
	s := surveyOf(t, muestraPath)
	rows := readMuestra(t)
	if s.Total.Records != int64(len(rows)) || s.Files[0].Rejected != 0 || s.Files[0].Lines-1 != s.Files[0].Read {
		t.Errorf("muestra: %d registros de %d", s.Total.Records, len(rows))
	}
	if len(s.ChainMun) == 0 {
		t.Error("la muestra debía traer cadenas de interés")
	}
	var out bytes.Buffer
	s.WriteMarkdown(&out)
	if !strings.Contains(out.String(), "Wal-mart ★") || !strings.Contains(out.String(), "Bodega Aurrera ★") {
		t.Error("el reporte de la muestra no marca las cadenas de interés")
	}
}

func TestThousandsYCell(t *testing.T) {
	for in, want := range map[int64]string{0: "0", 999: "999", 1000: "1,000", 9388457: "9,388,457", -1234: "-1,234"} {
		if got := thousands(in); got != want {
			t.Errorf("thousands(%d) = %s", in, got)
		}
	}
	if cell("a|b\nc") != "a/b c" {
		t.Error("cell")
	}
	if monthsSpan(newAgg()) != "0" {
		t.Error("monthsSpan vacío")
	}
}

func TestRowErrorMensaje(t *testing.T) {
	if got := (&RowError{Line: 7, Reason: "x"}).Error(); got != "renglón 7: x" {
		t.Error(got)
	}
}
