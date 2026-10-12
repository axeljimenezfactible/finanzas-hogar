package semilla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
)

// money lee un número JSON de pesos y lo vuelve centavos exactos (valor × 100 debe ser entero; nunca float).
type money int64

func (m *money) UnmarshalJSON(b []byte) error {
	r, ok := new(big.Rat).SetString(string(bytes.TrimSpace(b)))
	if !ok {
		return fmt.Errorf("número inválido %s", b)
	}
	r.Mul(r, big.NewRat(100, 1))
	if !r.IsInt() {
		return fmt.Errorf("%s no es una cantidad en centavos exactos", b)
	}
	*m = money(r.Num().Int64())
	return nil
}

type dataFile struct {
	Fecha          string `json:"fecha"`
	TiendaA        string `json:"tienda_a"`
	TiendaB        string `json:"tienda_b"`
	ItemsIdenticos []struct {
		Producto    string `json:"producto"`
		Cantidad    int64  `json:"cantidad"`
		BodegaTotal money  `json:"bodega_total"`
		Cuenta      string `json:"cuenta"`
	} `json:"items_identicos"`
	Sustitutos []struct {
		Bodega      string `json:"bodega"`
		CantB       int64  `json:"cant_b"`
		BodegaTotal money  `json:"bodega_total"`
	} `json:"sustitutos"`
	AtunSexta money `json:"atun_sexta_pieza_bodega"`
	Totales   struct {
		Bodega          money `json:"bodega"`
		ArticulosBodega int64 `json:"articulos_bodega"`
	} `json:"totales"`
	Cuentas map[string]struct {
		Lineas []struct {
			Producto string `json:"producto"`
			Cant     int64  `json:"cant"`
			Importe  money  `json:"importe"`
		} `json:"lineas"`
	} `json:"cuentas"`
	Financiamiento struct {
		DosQuincenas struct {
			Intereses money `json:"intereses"`
		} `json:"dos_quincenas"`
		CuatroQuincenas struct {
			Intereses money `json:"intereses"`
		} `json:"cuatro_quincenas"`
	} `json:"financiamiento"`
	Promociones json.RawMessage `json:"promociones"`
}

func readData(dir string) (*dataFile, error) {
	b, err := os.ReadFile(filepath.Join(dir, "data.json"))
	if err != nil {
		return nil, err
	}
	var d dataFile
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, fmt.Errorf("semilla: data.json: %w", err)
	}
	return &d, nil
}

type modelPago struct {
	Fecha string `json:"fecha"`
	Monto money  `json:"monto"`
	Axel  money  `json:"axel"`
	Karen money  `json:"karen"`
}

type modelPerson struct {
	Capital   money `json:"capital"`
	Conjunto  money `json:"conjunto"`
	Intereses money `json:"intereses"`
	Total     money `json:"total"`
}

type modelDespensa struct {
	N         int         `json:"n"`
	Fecha     string      `json:"fecha"`
	Capital   money       `json:"capital"`
	Intereses money       `json:"intereses"`
	Total     money       `json:"total"`
	Pagos     []modelPago `json:"pagos"`
	Axel      modelPerson `json:"axel"`
	Karen     modelPerson `json:"karen"`
}

// modelFile es la parte de oraculo/model.json que usa la semilla. Las cifras NO se copian: el motor las calcula y se cotejan.
type modelFile struct {
	Parametros struct {
		Limite money `json:"limite_credito"`
	} `json:"parametros"`
	Despensas []modelDespensa `json:"despensas"`
	Credito   struct {
		Eventos []struct {
			Fecha       string `json:"fecha"`
			Tipo        string `json:"tipo"`
			DisponibleA money  `json:"disponible_A"`
			DisponibleB money  `json:"disponible_B"`
		} `json:"eventos"`
	} `json:"credito"`
	Papel     json.RawMessage `json:"papel"`
	Horizonte struct {
		Fin       string `json:"fin"`
		DeudaViva money  `json:"deuda_viva_10_abr"`
	} `json:"horizonte"`
}

// ModelPath resuelve oraculo/model.json: la ruta dada; Dir/model.json (plan §3.4: copia en la semilla); o, en el
// repositorio, Dir/../../oraculo/model.json.
func ModelPath(dir, explicit string) (string, error) {
	cands := []string{explicit, filepath.Join(dir, "model.json"), filepath.Join(dir, "..", "..", "oraculo", "model.json")}
	for _, c := range cands {
		if c == "" {
			continue
		}
		if _, err := os.Stat(c); err == nil {
			return c, nil
		}
	}
	return "", fmt.Errorf("semilla: no encuentro model.json (probé %s y oraculo/model.json): el crédito Kueski se cota contra el oráculo", filepath.Join(dir, "model.json"))
}

func readModel(dir, explicit string) (*modelFile, error) {
	p, err := ModelPath(dir, explicit)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var m modelFile
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("semilla: %s: %w", p, err)
	}
	if len(m.Despensas) == 0 {
		return nil, fmt.Errorf("semilla: %s no trae despensas", p)
	}
	return &m, nil
}
