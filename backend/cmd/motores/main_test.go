package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDumpRecorteYParche(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "_datos"), 0o755)
	os.WriteFile(filepath.Join(dir, "_datos", "d.json"), []byte(`{"a":1,"b":{"c":2}}`), 0o644)
	os.WriteFile(filepath.Join(dir, "s9-x.json"), []byte(`{"formato":"fixture/1","motor":"s9","caso":"s9-x",
"entradas":{"datos":{"$ref":"_datos/d.json","parche":{"b":null,"z":3}},"consultas":[{"id":"q","funcion":"f","argumentos":{}},{"id":"e","funcion":"err","argumentos":{}}]},
"esperado":{"q":{"z":0,"falta":0,"l":[0]},"e":{"error":"x"}}}`), 0o644)
	ev := func(in json.RawMessage, q Query) (any, string) {
		if q.Func == "err" {
			return nil, "invalid_argument"
		}
		var m map[string]any
		json.Unmarshal(in, &m)
		m["l"] = []int{1, 2}
		return m, ""
	}
	b, err := dump("s9", ev, dir, "")
	want := `"q": {` + "\n" + `        "z": 3,` + "\n" + `        "falta": "__falta__",` + "\n" + `        "l": [` + "\n" + `          1,` + "\n" + `          "__sobra__"`
	if err != nil || !strings.Contains(string(b), want) || !strings.Contains(string(b), `"error": "invalid_argument"`) || strings.Contains(string(b), `"b"`) {
		t.Fatalf("dump: %v\n%s", err, b)
	}
	if run([]string{"s2"}, io.Discard) != 2 || run([]string{"nada", "dump"}, io.Discard) != 2 ||
		run([]string{"s2", "dump", "--fixtures", dir, "--out", "/tmp/x.json"}, io.Discard) != 2 {
		t.Error("uso incorrecto debe salir con 2")
	}
}
