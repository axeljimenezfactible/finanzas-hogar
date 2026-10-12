package precios

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ParseCents convierte un decimal en texto («3094.5», «21.0», «150») a centavos exactos.
func ParseCents(s string) (int64, error) {
	ip, fp, _ := strings.Cut(s, ".")
	if ip == "" || strings.HasPrefix(ip, "-") || strings.HasPrefix(ip, "+") {
		return 0, fmt.Errorf("%w: importe %q", ErrInvalidArgument, s)
	}
	fp = strings.TrimRight(fp, "0")
	if len(fp) > 2 {
		return 0, fmt.Errorf("%w: importe %q con más de 2 decimales", ErrInvalidArgument, s)
	}
	fp += strings.Repeat("0", 2-len(fp))
	v, err := strconv.ParseInt(ip+fp, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: importe %q", ErrInvalidArgument, s)
	}
	return v, nil
}

// LoadCatalogFile lee datos/semilla/productos.json (productos-formato.md §2): la categoría es
// el grupo; unidad_base pieza → piece; canal mercado → market, cadena → chain.
func LoadCatalogFile(path string) (*Catalog, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var f struct {
		Grupos []struct {
			ID, Nombre, Canal string
		} `json:"grupos"`
		Productos []struct {
			ID, Nombre, Categoria string
			UnidadBase            string `json:"unidad_base"`
			Presentaciones        []struct {
				ID, Nombre     string
				ContenidoMilli int64 `json:"contenido_milli"`
			} `json:"presentaciones"`
		} `json:"productos"`
	}
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidArgument, err)
	}
	var gs []Group
	for _, g := range f.Grupos {
		ch := map[string]string{"mercado": "market", "cadena": "chain"}[g.Canal]
		gs = append(gs, Group{g.ID, g.Nombre, ch})
	}
	var ps []Product
	var pr []Presentation
	for _, p := range f.Productos {
		u := p.UnidadBase
		if u == "pieza" {
			u = "piece"
		}
		ps = append(ps, Product{p.ID, p.Nombre, u, p.Categoria})
		for _, x := range p.Presentaciones {
			pr = append(pr, Presentation{x.ID, p.ID, x.Nombre, x.ContenidoMilli})
		}
	}
	return NewCatalog(gs, ps, pr)
}

// Report es la semilla del 10-oct-2026 (data.json) modelada como compra contra carrito (§3).
type Report struct {
	Bodega, Walmart Cart
	Totals          map[string]int64 // totales.* en centavos (importes) o piezas (articulos_*)
	Winners         map[string]int64 // resultados_por_producto
	Identical       int              // len(items_identicos)
	Substitutes     int              // len(sustitutos)
}

// LoadReport arma las dos canastas del 2026-10-10 desde data.json: «bodega» (compra registrada:
// 33 idénticos con bodega_unit, la sexta lata de atún y 3 sustitutos) y «walmart» (carrito: 33
// idénticos con walmart_unit y 3 sustitutos). Cada texto empata con el nombre exacto de una
// presentación del catálogo. Fuente manual: PROVISIONAL (bloqueo S5-carrito-comparativa p2).
func LoadReport(path string, c *Catalog) (*Report, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var f struct {
		Fecha string `json:"fecha"`
		Items []struct {
			Producto string
			Cantidad json.Number
			BU       json.Number `json:"bodega_unit"`
			WU       json.Number `json:"walmart_unit"`
		} `json:"items_identicos"`
		Sust []struct {
			Bodega, Walmart string
			CantB           json.Number `json:"cant_b"`
			CantW           json.Number `json:"cant_w"`
			TotB            json.Number `json:"bodega_total"`
			TotW            json.Number `json:"walmart_total"`
		} `json:"sustitutos"`
		Atun    json.Number            `json:"atun_sexta_pieza_bodega"`
		Totales map[string]json.Number `json:"totales"`
		Res     map[string]int64       `json:"resultados_por_producto"`
	}
	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.UseNumber()
	if err := dec.Decode(&f); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidArgument, err)
	}
	date, err := ParseDate(f.Fecha)
	if err != nil {
		return nil, err
	}
	byName := map[string]string{}
	for _, p := range c.Presentations {
		byName[p.Name] = p.Code
	}
	r := &Report{Totals: map[string]int64{}, Winners: f.Res, Identical: len(f.Items), Substitutes: len(f.Sust),
		Bodega:  Cart{Code: "bodega", Store: "bodega", Kind: "expense", ObservedDate: date, Source: SourceManual},
		Walmart: Cart{Code: "walmart", Store: "walmart", Kind: "cart", ObservedDate: date, Source: SourceManual}}
	line := func(name string, qty json.Number, unit, total string) (CartLine, error) {
		code, ok := byName[name]
		if !ok {
			return CartLine{}, fmt.Errorf("%w: %q", ErrUnknownPresentation, name)
		}
		q, err := strconv.ParseInt(qty.String(), 10, 64)
		if err != nil || q <= 0 {
			return CartLine{}, fmt.Errorf("%w: cantidad de %q", ErrInvalidArgument, name)
		}
		var u int64
		if unit != "" {
			u, err = ParseCents(unit)
		} else if t, err2 := ParseCents(total); err2 != nil {
			err = err2
		} else if t%q != 0 {
			err = fmt.Errorf("%w: %q no da centavos exactos", ErrInvalidArgument, name)
		} else {
			u = t / q
		}
		return CartLine{code, q * 1000, u}, err
	}
	atun := ""
	for _, it := range f.Items {
		lb, err := line(it.Producto, it.Cantidad, it.BU.String(), "")
		if err != nil {
			return nil, err
		}
		lw, err := line(it.Producto, it.Cantidad, it.WU.String(), "")
		if err != nil {
			return nil, err
		}
		r.Bodega.Lines = append(r.Bodega.Lines, lb)
		r.Walmart.Lines = append(r.Walmart.Lines, lw)
		if strings.HasPrefix(it.Producto, "Atún") {
			atun = lb.Presentation
		}
	}
	if f.Atun != "" {
		if atun == "" {
			return nil, fmt.Errorf("%w: sin renglón de atún para la sexta lata", ErrInvalidArgument)
		}
		u, err := ParseCents(f.Atun.String())
		if err != nil {
			return nil, err
		}
		r.Bodega.Lines = append(r.Bodega.Lines, CartLine{atun, 1000, u})
	}
	for _, s := range f.Sust {
		lb, err := line(s.Bodega, s.CantB, "", s.TotB.String())
		if err != nil {
			return nil, err
		}
		lw, err := line(s.Walmart, s.CantW, "", s.TotW.String())
		if err != nil {
			return nil, err
		}
		r.Bodega.Lines = append(r.Bodega.Lines, lb)
		r.Walmart.Lines = append(r.Walmart.Lines, lw)
	}
	for k, v := range f.Totales {
		if strings.HasPrefix(k, "articulos_") {
			n, err := strconv.ParseInt(v.String(), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%w: %s", ErrInvalidArgument, k)
			}
			r.Totals[k] = n
			continue
		}
		if r.Totals[k], err = ParseCents(v.String()); err != nil {
			return nil, err
		}
	}
	return r, nil
}
