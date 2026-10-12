package asistida

import (
	"errors"
	"testing"
)

func TestLecturaAsistidaSinVeredicto(t *testing.T) {
	full := []Term{{Site: SiteWalmart, Verdict: "permitido", URL: "https://ejemplo.invalid/t", ReadDate: "2026-10-11"}}
	// Apagada por defecto: sin interruptor global ni de tienda no se activa aunque haya veredicto.
	for _, c := range []struct {
		site          string
		global, store bool
		terms         []Term
	}{
		{SiteWalmart, false, true, full},
		{SiteWalmart, true, false, full},
		{"", true, true, full},
		{SiteWalmart, true, true, nil},
		{SiteBodegaAurrera, true, true, full},
		{SiteWalmart, true, true, []Term{{Site: SiteWalmart, Verdict: "permitido", URL: " ", ReadDate: "2026-10-11"}}},
		{SiteWalmart, true, true, []Term{{Site: SiteWalmart, URL: "u", ReadDate: "2026-10-11"}}},
		{SiteWalmart, true, true, []Term{{Site: SiteWalmart, Verdict: "v", URL: "u"}}},
	} {
		if err := CanEnable(c.site, c.global, c.store, c.terms); !errors.Is(err, ErrNoTermsVerdict) {
			t.Errorf("%+v: %v, se esperaba ErrNoTermsVerdict", c, err)
		}
	}
	if err := CanEnable(SiteWalmart, true, true, full); err != nil {
		t.Errorf("veredicto completo rechazado: %v", err)
	}
	// El archivo de términos de hoy no trae veredictos completos: ninguna cadena se activa.
	md := "# Términos\n\n## Walmart\n- **Veredicto:** permitido con baja frecuencia\n- URL: https://ejemplo.invalid/terminos\n- Fecha de lectura: 2026-10-11\n\n## Bodega Aurrera\nVeredicto: \n\n## Sam's Club\n- Veredicto: pendiente\n- URL: https://ejemplo.invalid/sams\n"
	terms := ParseTerms(md)
	if len(terms) != 3 || terms[0].Site != SiteWalmart || terms[1].Site != SiteBodegaAurrera || terms[2].Site != SiteSamsClub {
		t.Fatalf("ParseTerms: %+v", terms)
	}
	if terms[0].Verdict != "permitido con baja frecuencia" || terms[0].ReadDate != "2026-10-11" {
		t.Errorf("Walmart: %+v", terms[0])
	}
	if CanEnable(SiteWalmart, true, true, terms) != nil {
		t.Error("Walmart con veredicto completo rechazado")
	}
	for _, s := range []string{SiteBodegaAurrera, SiteSamsClub} {
		if CanEnable(s, true, true, terms) == nil {
			t.Errorf("%s activado sin veredicto completo", s)
		}
	}
	for _, c := range [][3]string{{"walmart", "", SiteWalmart}, {"x", "Bodega Aurrera Vallejo", SiteBodegaAurrera},
		{"x", "Sam's Club", SiteSamsClub}, {"x", "WALMART Acueducto", SiteWalmart}, {"x", "Tianguis", ""}} {
		if got := SiteFor(c[0], c[1]); got != c[2] {
			t.Errorf("SiteFor(%q, %q) = %q", c[0], c[1], got)
		}
	}
}
