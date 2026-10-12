// Paquete asistida: interruptor de la lectura asistida de sitios de tiendas (S5-07, §13).
// En WF-1 no implementa ningún cliente HTTP ni contiene URLs de tiendas: solo decide si el
// interruptor puede activarse a partir del veredicto escrito en docs/estado/S5-terminos.md.
package asistida

import (
	"errors"
	"strings"
)

// ErrNoTermsVerdict: falta el veredicto completo del sitio o algún interruptor está apagado.
var ErrNoTermsVerdict = errors.New("lectura asistida sin veredicto de términos")

// Sitios con veredicto obligatorio.
const (
	SiteWalmart       = "walmart"
	SiteBodegaAurrera = "bodega_aurrera"
	SiteSamsClub      = "sams_club"
)

// Term es el veredicto leído de un sitio; un campo vacío cuenta como faltante.
type Term struct {
	Site, Verdict, URL, ReadDate string
}

// CanEnable exige price.assisted_reading_enabled = true, store.assisted_reading_enabled = 1 y
// un veredicto completo (Veredicto, URL y Fecha de lectura) del sitio; si falta cualquiera,
// ErrNoTermsVerdict. Aun activada, site_assisted sigue siendo solo referencia (K37).
func CanEnable(site string, globalEnabled, storeEnabled bool, terms []Term) error {
	if site == "" || !globalEnabled || !storeEnabled {
		return ErrNoTermsVerdict
	}
	for _, t := range terms {
		if t.Site == site && strings.TrimSpace(t.Verdict) != "" && strings.TrimSpace(t.URL) != "" &&
			strings.TrimSpace(t.ReadDate) != "" {
			return nil
		}
	}
	return ErrNoTermsVerdict
}

// SiteFor deduce el sitio de una cadena: por código igual al del sitio o por su nombre
// (sin acentos ni mayúsculas): «bodega» → bodega_aurrera, «sam» → sams_club, «walmart» → walmart.
func SiteFor(code, name string) string {
	switch code {
	case SiteWalmart, SiteBodegaAurrera, SiteSamsClub:
		return code
	}
	n := strings.ToLower(name)
	switch {
	case strings.Contains(n, "bodega"):
		return SiteBodegaAurrera
	case strings.Contains(n, "sam"):
		return SiteSamsClub
	case strings.Contains(n, "walmart"):
		return SiteWalmart
	}
	return ""
}

// ParseTerms lee las secciones «## <sitio>» de S5-terminos.md con sus líneas «Veredicto:»,
// «URL:» y «Fecha de lectura:». El sitio sale de SiteFor sobre el título de la sección.
func ParseTerms(md string) []Term {
	var out []Term
	cur := -1
	for _, raw := range strings.Split(md, "\n") {
		l := strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(raw), "-*"))
		if strings.HasPrefix(raw, "## ") {
			out = append(out, Term{Site: SiteFor("", strings.TrimPrefix(raw, "## "))})
			cur = len(out) - 1
			continue
		}
		if cur < 0 {
			continue
		}
		l = strings.ReplaceAll(l, "**", "")
		for _, f := range []struct {
			key string
			dst *string
		}{{"Veredicto:", &out[cur].Verdict}, {"URL:", &out[cur].URL}, {"Fecha de lectura:", &out[cur].ReadDate}} {
			if v, ok := strings.CutPrefix(l, f.key); ok {
				*f.dst = strings.TrimSpace(v)
			}
		}
	}
	return out
}
