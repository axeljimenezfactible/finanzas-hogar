package profeco

import "strconv"

// branchKey es la llave de una sucursal dentro de su cadena.
func branchKey(chainID int64, folded string) string {
	return strconv.FormatInt(chainID, 10) + "|" + folded
}

// Chain es una cadena del catálogo (store.kind = 'chain') con su nombre en PROFECO.
type Chain struct {
	ID          int64
	Name        string
	ProfecoName string // cadena_comercial en PROFECO (store.profeco_name)
}

// Branch es una sucursal del catálogo (store.kind = 'branch').
type Branch struct {
	ID           int64
	ChainID      int64
	Name         string
	ProfecoName  string // nombre_comercial en PROFECO
	State        string // opcional; vacío = cualquier estado
	Municipality string
}

// StoreSet son las tiendas del catálogo contra las que se filtra.
type StoreSet struct {
	Chains   []Chain
	Branches []Branch
}

// Filter decide qué registros de PROFECO interesan: la cadena comercial ha de
// ser una de las cadenas y el municipio uno de los de las sucursales (S5-interfaz
// §11). Compara con Fold, así que no distingue acentos ni mayúsculas.
type Filter struct {
	chainByName  map[string]Chain
	branchByName map[string]Branch // clave: cadena|nombre_comercial
	munAny       map[string]bool   // municipio en cualquier estado
	munState     map[string]bool   // estado|municipio
}

// NewFilter arma el filtro. Sin sucursales (con municipio) no pasa ningún registro.
func NewFilter(s StoreSet) *Filter {
	f := &Filter{chainByName: map[string]Chain{}, branchByName: map[string]Branch{}, munAny: map[string]bool{}, munState: map[string]bool{}}
	for _, c := range s.Chains {
		if k := Fold(c.ProfecoName); k != "" {
			f.chainByName[k] = c
		}
	}
	for _, b := range s.Branches {
		if k := Fold(b.Municipality); k != "" {
			if st := Fold(b.State); st != "" {
				f.munState[st+"|"+k] = true
			} else {
				f.munAny[k] = true
			}
		}
		if k := Fold(b.ProfecoName); k != "" {
			f.branchByName[branchKey(b.ChainID, k)] = b
		}
	}
	return f
}

// Match devuelve la cadena a la que pertenece el registro y la tienda del
// precio: la sucursal cuyo profeco_name = nombre_comercial; si no hay, la
// cadena. ok es false si el registro no es de las cadenas o los municipios.
func (f *Filter) Match(rec *Record) (chain Chain, storeID int64, ok bool) {
	chain, ok = f.chainByName[Fold(rec.CadenaComercial)]
	if !ok {
		return Chain{}, 0, false
	}
	mun := Fold(rec.Municipio)
	if !f.munAny[mun] && !f.munState[Fold(rec.Estado)+"|"+mun] {
		return Chain{}, 0, false
	}
	if b, found := f.branchByName[branchKey(chain.ID, Fold(rec.NombreComercial))]; found {
		return chain, b.ID, true
	}
	return chain, chain.ID, true
}
