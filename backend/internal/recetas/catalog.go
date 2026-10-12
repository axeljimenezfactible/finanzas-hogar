package recetas

import (
	"encoding/json"
	"os"
)

// Catalog es el catálogo de grupos, productos y presentaciones (productos-formato.md).
type Catalog struct {
	Groups        map[string]*Group
	Products      map[string]*Product
	Presentations map[string]*Presentation
	presOrder     []string
}

// NewCatalog valida y arma el catálogo: códigos únicos, canal y unidad válidos, grupo existente,
// gramos por unidad solo en ml o pieza (CHECK del esquema) y contenido > 0.
func NewCatalog(groups []Group, products []Product, pres []Presentation) (*Catalog, error) {
	c := &Catalog{Groups: map[string]*Group{}, Products: map[string]*Product{}, Presentations: map[string]*Presentation{}}
	for i := range groups {
		g := groups[i]
		if g.Code == "" || c.Groups[g.Code] != nil || (g.Channel != ChannelMarket && g.Channel != ChannelChain) {
			return nil, errf(ErrInvalidArgument, "grupo inválido o repetido %q", g.Code)
		}
		c.Groups[g.Code] = &g
	}
	for i := range products {
		p := products[i]
		if p.Code == "" || c.Products[p.Code] != nil {
			return nil, errf(ErrInvalidArgument, "producto inválido o repetido %q", p.Code)
		}
		if p.BaseUnit != UnitG && p.BaseUnit != UnitML && p.BaseUnit != UnitPiece {
			return nil, errf(ErrInvalidArgument, "unidad base inválida de %q", p.Code)
		}
		if p.Group != "" && c.Groups[p.Group] == nil {
			return nil, errf(ErrInvalidArgument, "grupo desconocido %q en %q", p.Group, p.Code)
		}
		if p.GramsPerUnitMilli != nil && (p.BaseUnit == UnitG || *p.GramsPerUnitMilli <= 0) {
			return nil, errf(ErrInvalidArgument, "gramos por unidad inválidos en %q", p.Code)
		}
		c.Products[p.Code] = &p
	}
	for i := range pres {
		x := pres[i]
		if x.Code == "" || c.Presentations[x.Code] != nil || x.SizeMilli <= 0 {
			return nil, errf(ErrInvalidArgument, "presentación inválida o repetida %q", x.Code)
		}
		if c.Products[x.Product] == nil {
			return nil, errf(ErrUnknownProduct, "presentación %q de producto desconocido %q", x.Code, x.Product)
		}
		c.Presentations[x.Code] = &x
		c.presOrder = append(c.presOrder, x.Code)
	}
	return c, nil
}

// presentationOf devuelve la presentación si existe y es del producto.
func (c *Catalog) presentationOf(product, code string) (*Presentation, error) {
	x := c.Presentations[code]
	if x == nil || x.Product != product {
		return nil, errf(ErrUnknownPresentation, "presentación %q del producto %q", code, product)
	}
	return x, nil
}

// fallbackPresentation es la presentación típica del producto o, si no hay, la primera del catálogo.
func (c *Catalog) fallbackPresentation(product string) string {
	if p := c.Products[product]; p != nil && p.TypicalPresentation != "" {
		return p.TypicalPresentation
	}
	for _, code := range c.presOrder {
		if c.Presentations[code].Product == product {
			return code
		}
	}
	return ""
}

// channelOf es la vertical del producto por su grupo; sin grupo es ErrProductWithoutChannel (§6.5, S4-12).
func (c *Catalog) channelOf(product string) (Channel, string, error) {
	p := c.Products[product]
	if p == nil {
		return "", "", errf(ErrUnknownProduct, "%q", product)
	}
	g := c.Groups[p.Group]
	if g == nil {
		return "", "", errf(ErrProductWithoutChannel, "%q sin grupo de compra", product)
	}
	return g.Channel, g.Code, nil
}

type catFile struct {
	Grupos []struct {
		ID     string `json:"id"`
		Nombre string `json:"nombre"`
		Canal  string `json:"canal"`
		Orden  int    `json:"orden"`
	} `json:"grupos"`
	Productos []struct {
		ID                 string  `json:"id"`
		Nombre             string  `json:"nombre"`
		Categoria          string  `json:"categoria"`
		UnidadBase         string  `json:"unidad_base"`
		GramosPorUnidad    *string `json:"gramos_por_unidad"`
		PresentacionTipica string  `json:"presentacion_tipica"`
		Presentaciones     []struct {
			ID             string `json:"id"`
			Nombre         string `json:"nombre"`
			ContenidoMilli int64  `json:"contenido_milli"`
		} `json:"presentaciones"`
	} `json:"productos"`
}

var channelES = map[string]Channel{"mercado": ChannelMarket, "cadena": ChannelChain}
var unitES = map[string]Unit{"g": UnitG, "ml": UnitML, "pieza": UnitPiece}

// LoadCatalogFile lee datos/semilla/productos.json con las reglas de productos-formato.md §2.
func LoadCatalogFile(path string) (*Catalog, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, errf(ErrInvalidArgument, "no se pudo leer el catálogo: %v", err)
	}
	return ParseCatalog(b)
}

// ParseCatalog arma el catálogo desde el JSON de productos.json.
func ParseCatalog(b []byte) (*Catalog, error) {
	var f catFile
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, errf(ErrInvalidArgument, "catálogo: %v", err)
	}
	var gs []Group
	var ps []Product
	var xs []Presentation
	for _, g := range f.Grupos {
		gs = append(gs, Group{Code: g.ID, Name: g.Nombre, Channel: channelES[g.Canal], SortOrder: g.Orden})
	}
	for _, p := range f.Productos {
		u, ok := unitES[p.UnidadBase]
		if !ok {
			return nil, errf(ErrInvalidArgument, "unidad_base %q de %q", p.UnidadBase, p.ID)
		}
		prod := Product{Code: p.ID, Name: p.Nombre, BaseUnit: u, Group: p.Categoria, TypicalPresentation: p.PresentacionTipica}
		if p.GramosPorUnidad != nil {
			m, err := parseMilli(*p.GramosPorUnidad)
			if err != nil {
				return nil, errf(ErrInvalidArgument, "gramos_por_unidad de %q: %v", p.ID, err)
			}
			prod.GramsPerUnitMilli = &m
		}
		ps = append(ps, prod)
		for _, x := range p.Presentaciones {
			xs = append(xs, Presentation{Code: x.ID, Product: p.ID, Name: x.Nombre, SizeMilli: x.ContenidoMilli})
		}
	}
	return NewCatalog(gs, ps, xs)
}
