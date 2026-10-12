// Comando motores: binario común de los motores de Finanzas Hogar (FORMATO.md §4 y §4.1).
//
//	motores <motor> dump --fixtures <carpeta> --out datos/tmp/<archivo> [--caso <nombre>]
//
// Cada motor se registra desde su propio archivo (<motor>.go) con
// func init() { register("<motor>", evaluador) }. Códigos de salida: 0 escribió,
// 1 no pudo leer o escribir, 2 uso incorrecto o motor sin registrar.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Query es una consulta de una fixture.
type Query struct {
	ID   string          `json:"id"`
	Func string          `json:"funcion"`
	Args json.RawMessage `json:"argumentos"`
}

// Evaluator calcula una consulta. input es entradas.datos ya resuelto ($ref y parche aplicados).
// Devuelve la salida (un valor que encoding/json serializa, con dinero en enteros de centavos)
// o, si la función falla, code != "" con el código de error que declara la interfaz del motor.
type Evaluator func(input json.RawMessage, q Query) (out any, code string)

var registry = map[string]Evaluator{}

// register da de alta un motor; un nombre repetido es pánico al arrancar.
func register(motor string, ev Evaluator) {
	if _, dup := registry[motor]; dup {
		panic("motor registrado dos veces: " + motor)
	}
	registry[motor] = ev
}

func main() { os.Exit(run(os.Args[1:], os.Stderr)) }

func run(args []string, stderr io.Writer) int {
	use := func(msg string) int {
		fmt.Fprintln(stderr, "uso: motores <motor> dump --fixtures <carpeta> --out datos/tmp/<archivo> [--caso <nombre>]:", msg)
		return 2
	}
	if len(args) < 2 || args[1] != "dump" {
		return use("falta el subcomando dump")
	}
	ev, ok := registry[args[0]]
	if !ok {
		return use("motor sin registrar " + args[0])
	}
	fs := flag.NewFlagSet("dump", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fix := fs.String("fixtures", "", "carpeta de fixtures")
	out := fs.String("out", "", "archivo de salida bajo datos/tmp/")
	caso := fs.String("caso", "", "solo este caso")
	if err := fs.Parse(args[2:]); err != nil || *fix == "" || *out == "" || fs.NArg() > 0 {
		return use("--fixtures y --out son obligatorios")
	}
	abs, err := filepath.Abs(*out)
	if err != nil || !strings.Contains(filepath.ToSlash(abs), "/datos/tmp/") {
		return use("--out debe quedar bajo datos/tmp/")
	}
	doc, err := dump(args[0], ev, *fix, *caso)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	if err := os.WriteFile(*out, doc, 0o644); err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	return 0
}

type kv struct {
	K string
	V any
}

// ordered es un objeto JSON que conserva el orden de sus claves.
type ordered []kv

func (o ordered) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, e := range o {
		if i > 0 {
			b.WriteByte(',')
		}
		k, _ := json.Marshal(e.K)
		v, err := json.Marshal(e.V)
		if err != nil {
			return nil, err
		}
		b.Write(k)
		b.WriteByte(':')
		b.Write(v)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func (o ordered) get(k string) (any, bool) {
	for _, e := range o {
		if e.K == k {
			return e.V, true
		}
	}
	return nil, false
}

// parseOrdered lee un valor JSON conservando el orden de las claves y los números exactos.
func parseOrdered(dec *json.Decoder) (any, error) {
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	d, ok := t.(json.Delim)
	if !ok {
		return t, nil
	}
	switch d {
	case '{':
		o := ordered{}
		for dec.More() {
			kt, err := dec.Token()
			if err != nil {
				return nil, err
			}
			v, err := parseOrdered(dec)
			if err != nil {
				return nil, err
			}
			o = append(o, kv{kt.(string), v})
		}
		_, err = dec.Token()
		return o, err
	case '[':
		l := []any{}
		for dec.More() {
			v, err := parseOrdered(dec)
			if err != nil {
				return nil, err
			}
			l = append(l, v)
		}
		_, err = dec.Token()
		return l, err
	}
	return nil, errors.New("JSON inválido")
}

func generic(b []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var v any
	err := dec.Decode(&v)
	return v, err
}

// mergePatch aplica un JSON Merge Patch (RFC 7396).
func mergePatch(target, patch any) any {
	pm, ok := patch.(map[string]any)
	if !ok {
		return patch
	}
	tm, ok := target.(map[string]any)
	if !ok {
		tm = map[string]any{}
	}
	for k, v := range pm {
		if v == nil {
			delete(tm, k)
			continue
		}
		tm[k] = mergePatch(tm[k], v)
	}
	return tm
}

// resolve aplica $ref (relativo a la carpeta del caso) y parche a entradas.datos.
func resolve(dir string, raw json.RawMessage) (json.RawMessage, error) {
	v, err := generic(raw)
	if err != nil {
		return nil, err
	}
	m, ok := v.(map[string]any)
	if !ok {
		return raw, nil
	}
	ref, ok := m["$ref"].(string)
	if !ok {
		return raw, nil
	}
	b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(ref)))
	if err != nil {
		return nil, err
	}
	base, err := generic(b)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", ref, err)
	}
	if p, ok := m["parche"]; ok {
		base = mergePatch(base, p)
	}
	return json.Marshal(base)
}

// trim recorta la salida calculada con esperado como plantilla (FORMATO §4, reglas 1–3 y 5).
func trim(val, tmpl any) any {
	switch t := tmpl.(type) {
	case ordered:
		m, ok := val.(map[string]any)
		if !ok {
			return "__tipo__"
		}
		o := ordered{}
		for _, e := range t {
			if v, ok := m[e.K]; ok {
				o = append(o, kv{e.K, trim(v, e.V)})
			} else {
				o = append(o, kv{e.K, "__falta__"})
			}
		}
		return o
	case []any:
		l, ok := val.([]any)
		if !ok {
			return "__tipo__"
		}
		out := []any{}
		for i, v := range l {
			if i < len(t) {
				out = append(out, trim(v, t[i]))
			} else {
				out = append(out, "__sobra__")
			}
		}
		return out
	}
	switch val.(type) {
	case map[string]any, []any:
		return "__tipo__"
	}
	return val
}

type fixture struct {
	Formato  string `json:"formato"`
	Motor    string `json:"motor"`
	Caso     string `json:"caso"`
	Entradas struct {
		Datos     json.RawMessage `json:"datos"`
		Consultas []Query         `json:"consultas"`
	} `json:"entradas"`
	Esperado json.RawMessage `json:"esperado"`
}

// evalCase evalúa las consultas de una fixture y recorta cada salida.
func evalCase(path, motor string, ev Evaluator) (ordered, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var fx fixture
	if err := json.Unmarshal(raw, &fx); err != nil {
		return nil, fmt.Errorf("%s: %v", path, err)
	}
	name := strings.TrimSuffix(filepath.Base(path), ".json")
	if fx.Formato != "fixture/1" || fx.Motor != motor || fx.Caso != name {
		return nil, fmt.Errorf("%s: formato, motor o caso inválido", path)
	}
	datos, err := resolve(filepath.Dir(path), fx.Entradas.Datos)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", path, err)
	}
	var esp ordered
	if len(fx.Esperado) > 0 {
		e, err := parseOrdered(json.NewDecoder(bytes.NewReader(fx.Esperado)))
		if err != nil {
			return nil, fmt.Errorf("%s: esperado: %v", path, err)
		}
		esp, _ = e.(ordered)
	}
	res := ordered{}
	for _, q := range fx.Entradas.Consultas {
		out, code := ev(datos, q)
		if code != "" {
			res = append(res, kv{q.ID, ordered{{"error", code}}})
			continue
		}
		b, err := json.Marshal(out)
		if err != nil {
			return nil, fmt.Errorf("%s/%s: %v", name, q.ID, err)
		}
		v, _ := generic(b)
		if t, ok := esp.get(q.ID); ok {
			v = trim(v, t)
		}
		res = append(res, kv{q.ID, v})
	}
	return res, nil
}

// dump recorre las fixtures (*.json que no empiezan con «_», por nombre) y arma salida/1.
func dump(motor string, ev Evaluator, dir, only string) ([]byte, error) {
	names, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}
	sort.Strings(names)
	casos := ordered{}
	for _, p := range names {
		name := strings.TrimSuffix(filepath.Base(p), ".json")
		if strings.HasPrefix(name, "_") || (only != "" && name != only) {
			continue
		}
		res, err := evalCase(p, motor, ev)
		if err != nil {
			return nil, err
		}
		casos = append(casos, kv{name, res})
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	err = enc.Encode(ordered{{"formato", "salida/1"}, {"motor", motor}, {"casos", casos}})
	return buf.Bytes(), err
}
