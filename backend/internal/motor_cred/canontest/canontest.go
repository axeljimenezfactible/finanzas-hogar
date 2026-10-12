// Package canontest arma, solo para las pruebas, las entradas canon de S3 leyendo al correr
// oraculo/model.json y datos/semilla/data.json (S3-interfaz.md §10.0: seedFromModel y
// sensitivityFromModel). Ninguna cifra canon está escrita aquí: capital, intereses, partes,
// fechas, límite y contado salen de las hojas; lo que no es hoja (los 15 días, la lectura "B",
// los decimales 2 y 4 de K12, el plazo de 2 quincenas, los meses de devengo y los códigos de
// persona) sale de la interfaz.
package canontest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	mc "finanzas-hogar/internal/motor_cred"
)

// Root devuelve la raíz del repositorio (la carpeta que contiene oraculo/model.json).
func Root() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "oraculo", "model.json")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("no encontré oraculo/model.json")
		}
		dir = parent
	}
}

// Decode lee JSON con números exactos (json.Number).
func Decode(b []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var v any
	err := dec.Decode(&v)
	return v, err
}

// Load lee un archivo del repositorio (ruta relativa a la raíz).
func Load(t testing.TB, rel string) any {
	t.Helper()
	root, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, rel))
	if err != nil {
		t.Fatal(err)
	}
	v, err := Decode(b)
	if err != nil {
		t.Fatalf("%s: %v", rel, err)
	}
	return v
}

// JSON pasa un valor de Go por json.Marshal y lo vuelve a leer como árbol genérico.
func JSON(t testing.TB, v any) any {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	out, err := Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Get sigue una ruta como "despensas[0].pagos[1].monto".
func Get(v any, path string) (any, error) {
	cur := v
	for _, seg := range strings.Split(path, ".") {
		name, idx := seg, []int(nil)
		if k := strings.IndexByte(seg, '['); k >= 0 {
			name = seg[:k]
			for _, s := range strings.Split(strings.TrimSuffix(seg[k+1:], "]"), "][") {
				i, err := strconv.Atoi(s)
				if err != nil {
					return nil, fmt.Errorf("%s: índice %q", path, s)
				}
				idx = append(idx, i)
			}
		}
		if name != "" {
			m, ok := cur.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("%s: %q no está en un objeto", path, name)
			}
			if cur, ok = m[name]; !ok {
				return nil, fmt.Errorf("%s: falta %q", path, name)
			}
		}
		for _, i := range idx {
			l, ok := cur.([]any)
			if !ok || i < 0 || i >= len(l) {
				return nil, fmt.Errorf("%s: índice %d fuera de la lista", path, i)
			}
			cur = l[i]
		}
	}
	return cur, nil
}

func get(t testing.TB, v any, path string) any {
	t.Helper()
	x, err := Get(v, path)
	if err != nil {
		t.Fatal(err)
	}
	return x
}

// CentsOf convierte un número de model.json a centavos: valor × 100 exacto (FORMATO.md §3).
func CentsOf(v any) (mc.Cents, error) {
	n, ok := v.(json.Number)
	if !ok {
		return 0, fmt.Errorf("%v no es número", v)
	}
	r, ok := new(big.Rat).SetString(n.String())
	if !ok {
		return 0, fmt.Errorf("%v no es decimal", v)
	}
	r.Mul(r, big.NewRat(100, 1))
	if !r.IsInt() {
		return 0, fmt.Errorf("%v × 100 no es entero", v)
	}
	return mc.Cents(r.Num().Int64()), nil
}

// C lee dinero en centavos.
func C(t testing.TB, v any, path string) mc.Cents {
	t.Helper()
	c, err := CentsOf(get(t, v, path))
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return c
}

// S lee un texto.
func S(t testing.TB, v any, path string) string {
	t.Helper()
	s, ok := get(t, v, path).(string)
	if !ok {
		t.Fatalf("%s no es texto", path)
	}
	return s
}

// D lee una fecha ISO.
func D(t testing.TB, v any, path string) mc.Date {
	t.Helper()
	d, err := mc.ParseDate(S(t, v, path))
	if err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return d
}

// Len es el largo de una lista.
func Len(t testing.TB, v any, path string) int {
	t.Helper()
	l, ok := get(t, v, path).([]any)
	if !ok {
		t.Fatalf("%s no es lista", path)
	}
	return len(l)
}

func ptr[T any](v T) *T { return &v }

// Seed es seedFromModel: el Input de la semilla Kueski (S3-interfaz.md §10.0 «Datos compartidos»).
func Seed(t testing.TB, m, d any) mc.Input {
	t.Helper()
	cap0 := C(t, m, "despensas[0].capital")
	today := D(t, m, "despensas[0].fecha")
	n4 := Len(t, m, "despensas[0].pagos")
	credit := mc.Credit{ID: 1, Name: "Kueski", Kind: mc.KindRevolving, Holder: "axel",
		Limit: ptr(C(t, m, "parametros.limite_credito")), Frequency: "fortnightly",
		FirstPaymentMinDays: ptr(15), Reading: "B", ReadingAAnchor: ptr(today),
		TermRates: []mc.TermRate{
			{TermCount: 2, Rate: mc.Rate{Num: int64(C(t, d, "financiamiento.dos_quincenas.intereses")), Den: int64(cap0), Decimals: 2}, Source: "manual"},
			{TermCount: n4, Rate: mc.Rate{Num: int64(C(t, m, "despensas[0].intereses")), Den: int64(cap0), Decimals: 4}, Source: "manual"},
		},
		Promotions: []mc.Promotion{}}
	st := mc.CreditState{Credit: credit, Payments: []mc.PaymentRecord{}, Statements: []mc.StatementInput{}}
	for k := 0; k < Len(t, m, "despensas"); k++ {
		p := fmt.Sprintf("despensas[%d]", k)
		share := func(code string, ext bool) mc.ShareInput {
			return mc.ShareInput{Code: code, External: ext, Principal: C(t, m, p+"."+code+".capital"),
				Joint: ptr(C(t, m, p+"."+code+".conjunto")), PaperIn: ptr(C(t, m, p+"."+code+".papel_dentro"))}
		}
		in := mc.PurchaseInput{Date: D(t, m, p+".fecha"), Principal: C(t, m, p+".capital"),
			TermCount: Len(t, m, p+".pagos"), Shares: []mc.ShareInput{share("axel", false), share("karen", true)}}
		if k == 0 {
			in.QuotedInterest, in.RateDecimals = ptr(C(t, m, p+".intereses")), ptr(4)
		}
		st.Purchases = append(st.Purchases, mc.PurchaseState{Ref: fmt.Sprintf("d%d", k+1), Status: "active", Input: in})
	}
	cash := []mc.CashEvent{}
	for k := 0; k < Len(t, m, "papel"); k++ {
		p := fmt.Sprintf("papel[%d]", k)
		if S(t, m, p+".forma") == "de contado" {
			cash = append(cash, mc.CashEvent{Date: D(t, m, p+".compra"), Amount: C(t, m, p+".monto"),
				Person: strings.ToLower(S(t, m, p+".quien")), Description: "Papel de contado"})
		}
	}
	return mc.Input{Today: today, Credits: []mc.CreditState{st}, Cash: cash}
}

// Sensitivity es sensitivityFromModel: el insumo de vista_sensibilidad y los montos de
// vista_interes_por_monto (S3-interfaz.md §10.0, corrección r4-P1).
func Sensitivity(t testing.TB, m any) (mc.SensitivityInput, []mc.Cents) {
	t.Helper()
	si := mc.SensitivityInput{
		Rate:                   mc.Rate{Num: int64(C(t, m, "despensas[0].intereses")), Den: int64(C(t, m, "despensas[0].capital")), Decimals: 4},
		TermCount:              Len(t, m, "despensas[0].pagos"),
		External:               mc.ShareInput{Code: "karen", External: true, Principal: C(t, m, "parametros.karen_fija.capital")},
		HolderCurrentPrincipal: C(t, m, "despensas[0].axel.capital"),
		AccrualMonths:          2,
		CashMonthly:            mc.RoundHalfUp(big.NewInt(int64(C(t, m, "parametros.papel"))), big.NewInt(2)),
	}
	for k := 0; k < Len(t, m, "sensibilidad"); k++ {
		si.Amounts = append(si.Amounts, C(t, m, fmt.Sprintf("sensibilidad[%d].N", k)))
	}
	var ipm []mc.Cents
	for k := 0; k < Len(t, m, "interes_por_monto"); k++ {
		ipm = append(ipm, C(t, m, fmt.Sprintf("interes_por_monto[%d].monto", k)))
	}
	return si, ipm
}
