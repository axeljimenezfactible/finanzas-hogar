package profeco

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const muestraPath = "../../../../datos/semilla/profeco/muestra/qqp-muestra.csv"

const header = "producto,presentacion,marca,categoria,catalogo,precio,fecha_registro,cadena_comercial,giro,nombre_comercial,direccion,estado,municipio,latitud,longitud\n"

// writeTmp escribe un CSV de prueba y devuelve su ruta.
func writeTmp(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func row(prod, pres, marca, precio, fecha, cadena, nombre, dir, edo, mun string) string {
	return fmt.Sprintf("%s,%s,%s,Cat,PACIC,%s,%s,%s,Supermercado,%s,%s,%s,%s,19.5,-99.1\n", prod, pres, marca, precio, fecha, cadena, nombre, dir, edo, mun)
}

// catalogo de prueba: Wal-mart y Bodega Aurrera en Gustavo a. Madero (CDMX) y
// Wal-mart en Benito Juárez (solo en Quintana Roo, con estado).
func testStores() StoreSet {
	return StoreSet{
		Chains: []Chain{
			{ID: 1, Name: "Walmart", ProfecoName: "Wal-mart"},
			{ID: 2, Name: "Bodega Aurrera", ProfecoName: "Bodega Aurrera"},
			{ID: 3, Name: "Sam's Club", ProfecoName: "Sam's Club"},
		},
		Branches: []Branch{
			{ID: 11, ChainID: 1, Name: "Tepeyac", ProfecoName: "Walmart Sucursal Tepeyac", State: "Ciudad de México", Municipality: "Gustavo a. Madero"},
			{ID: 12, ChainID: 2, Name: "Vallejo", ProfecoName: "Bodega Aurrera Express Sucursal Poniente 140", Municipality: "Gustavo A. Madero"},
			{ID: 13, ChainID: 1, Name: "Cancún", ProfecoName: "Wal Mart", State: "Quintana Roo", Municipality: "Benito Juárez"},
		},
	}
}

func mustLinks(t *testing.T, csvBody string) *Links {
	t.Helper()
	l, err := ReadLinks(strings.NewReader("producto,presentacion,marca,presentacion_id,estado,confirmo,fecha_confirmacion\n" + csvBody))
	if err != nil {
		t.Fatal(err)
	}
	return l
}

func TestProfecoCodificacion(t *testing.T) {
	// UTF-8 con BOM, columnas en otro orden y encabezado en mayúsculas.
	utf := "\xEF\xBB\xBFPRODUCTO,Presentacion,MARCA,categoria,catalogo,precio,fecha_registro,cadena_comercial,giro,nombre_comercial,direccion,estado,MUNICIPIO,latitud,longitud\n" +
		row("Atún", "Lata 140 Gr.", "Dolores", "21.5", "2026/03/02", "Wal-mart", "Walmart Sucursal Tepeyac", "Calz. de Guadalupe 431", "Ciudad de México", "Gustavo a. Madero")
	rd, err := NewReader(strings.NewReader(utf))
	if err != nil {
		t.Fatal(err)
	}
	rec, err := rd.Next()
	if err != nil || rec.Producto != "Atún" || rec.PrecioCents != 2150 || rec.Fecha != "2026-03-02" || rec.Line != 2 {
		t.Fatalf("UTF-8: %+v %v", rec, err)
	}
	if _, err := rd.Next(); !errors.Is(err, io.EOF) {
		t.Fatalf("esperaba EOF, %v", err)
	}

	// Windows-1252: «Atún» = 0xFA, «Ciudad de México» = 0xE9; así viene QQP_metadatos.csv.
	win := header + "At\xfan,Lata 140 Gr.,Dolores,Cat,PACIC,21.5,2026/03/02,Wal-mart,Super,Walmart Sucursal Tepeyac,Calz. de Guadalupe 431,Ciudad de M\xe9xico,Gustavo a. Madero,19.5,-99.1\n"
	rd, err = NewReader(strings.NewReader(win))
	if err != nil {
		t.Fatal(err)
	}
	rec, err = rd.Next()
	if err != nil || rec.Producto != "Atún" || rec.Estado != "Ciudad de México" {
		t.Fatalf("Windows-1252: %+v %v", rec, err)
	}

	// Saltos de línea dentro de comillas se cuentan como un solo registro y se limpian.
	multi := header + "\"Atún\nlata\",P,M,C,K,10,2026-03-02,Wal-mart,G,N,\"Calle 1\n  Col. X\",E,M,,\n"
	rd, _ = NewReader(strings.NewReader(multi))
	rec, err = rd.Next()
	if err != nil || rec.Producto != "Atún lata" || rec.Direccion != "Calle 1 Col. X" || rec.Latitud != "" {
		t.Fatalf("multilínea: %+v %v", rec, err)
	}
	if rd.Read != 1 {
		t.Errorf("Read = %d; want 1", rd.Read)
	}
}

func TestProfecoEncabezado(t *testing.T) {
	for name, body := range map[string]string{
		"vacío":    "",
		"faltan":   "producto,precio\nA,1\n",
		"repetida": strings.Replace(header, "marca", "producto", 1),
	} {
		if _, err := NewReader(strings.NewReader(body)); !errors.Is(err, ErrProfecoHeader) {
			t.Errorf("%s: esperaba ErrProfecoHeader, got %v", name, err)
		}
	}
	p := writeTmp(t, "x.csv", "producto,precio\n")
	if _, err := OpenFile(p, 0); !errors.Is(err, ErrProfecoHeader) {
		t.Errorf("OpenFile: %v", err)
	}
	if _, err := OpenFile(filepath.Join(t.TempDir(), "no-existe.csv"), 0); err == nil {
		t.Error("archivo inexistente debía fallar")
	}
}

func TestProfecoArchivoGrande(t *testing.T) {
	p := writeTmp(t, "g.csv", header+row("A", "P", "M", "1", "2026/01/02", "Wal-mart", "N", "D", "E", "M"))
	if _, err := OpenFile(p, 10); !errors.Is(err, ErrFileTooLarge) {
		t.Errorf("esperaba ErrFileTooLarge, got %v", err)
	}
	im, _ := NewImporter(Options{MaxBytes: 10})
	if _, err := im.ImportFile(context.Background(), p); !errors.Is(err, ErrFileTooLarge) {
		t.Errorf("ImportFile: %v", err)
	}
}

func TestProfecoConteo(t *testing.T) {
	body := header +
		row("Atún", "Lata", "Dolores", "21.5", "2026/03/02", "Wal-mart", "Walmart Sucursal Tepeyac", "Calz", "Ciudad de México", "Gustavo a. Madero") + // bueno
		row("Atún", "Lata", "Dolores", "abc", "2026/03/02", "Wal-mart", "N", "D", "E", "Gustavo a. Madero") + // precio ilegible
		row("Atún", "Lata", "Dolores", "21.5", "ayer", "Wal-mart", "N", "D", "E", "Gustavo a. Madero") + // fecha ilegible
		row("", "Lata", "Dolores", "21.5", "2026/03/02", "Wal-mart", "N", "D", "E", "M") + // sin producto
		row("Atún", "Lata", "Dolores", "21.5", "2026/03/02", "", "N", "D", "E", "M") + // sin cadena
		row("Atún", "Lata", "Dolores", "21.5", "2026/12/31", "Wal-mart", "N", "D", "E", "Gustavo a. Madero") + // futura
		"solo,tres,campos\n" + // columnas de menos
		"\"abierta,x,y\n" // comilla sin cerrar al final
	p := writeTmp(t, "conteo.csv", body)
	links := mustLinks(t, "Atún,Lata,Dolores,7,confirmado,axel,2026-10-11\n")
	im, err := NewImporter(Options{Corte: "2026", Stores: testStores(), Links: links, Today: "2026-10-11", Prices: &MemoryPrices{}})
	if err != nil {
		t.Fatal(err)
	}
	res, err := im.ImportFile(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
	if res.Read != 8 || res.Rejected != 7 || res.Loaded != 1 {
		t.Fatalf("conteos: %+v", res)
	}
	if !res.Balanced() {
		t.Errorf("no cuadra: %+v", res)
	}
	for _, want := range []string{"precio ilegible", "fecha ilegible", "producto vacío", "cadena comercial vacía", "fecha futura"} {
		if res.RejectReasons[want] != 1 {
			t.Errorf("falta el motivo %q: %v", want, res.RejectReasons)
		}
	}
	if len(res.Rejections) != 7 || res.Rejections[0].Row != 2 {
		t.Errorf("detalle de rechazos: %+v", res.Rejections)
	}
	if res.MaxDate != "2026-03-02" || res.MinDate != "2026-03-02" || res.Lines != 9 {
		t.Errorf("fechas o líneas: %+v", res)
	}
}

// expected recorre la muestra por su cuenta (sin usar el cargador) y devuelve
// los renglones que el filtro debe dejar pasar.
type muestraRow struct {
	line                              int
	rec                               []string
	chainFold, munFold, stateFold     string
	productKey                        string
	cadena, nombre, municipio, estado string
}

func readMuestra(t *testing.T) []muestraRow {
	t.Helper()
	f, err := os.Open(muestraPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	cr := csv.NewReader(f)
	cr.LazyQuotes = true
	head, err := cr.Read()
	if err != nil {
		t.Fatal(err)
	}
	if len(head) != 15 || head[0] != "producto" {
		t.Fatalf("encabezado de la muestra: %v", head)
	}
	var out []muestraRow
	for {
		r, err := cr.Read()
		if errors.Is(err, io.EOF) {
			return out
		}
		if err != nil {
			t.Fatal(err)
		}
		line, _ := cr.FieldPos(0)
		out = append(out, muestraRow{line: line, rec: r, chainFold: Fold(r[7]), munFold: Fold(r[12]), stateFold: Fold(r[11]),
			productKey: LinkKey(r[0], r[1], r[2]), cadena: r[7], nombre: r[9], municipio: r[12], estado: r[11]})
	}
}

func TestProfecoMuestra(t *testing.T) {
	rows := readMuestra(t)
	if len(rows) < 1500 || len(rows) > 2500 {
		t.Fatalf("la muestra debía traer unas 2,000 filas y trae %d", len(rows))
	}
	st, err := os.Stat(muestraPath)
	if err != nil || st.Size() > 2<<20 {
		t.Fatalf("la muestra pesa más de 2 MB: %v %v", st, err)
	}
	im, _ := NewImporter(Options{Corte: "2026", Stores: testStores(), Mode: ModeReview})
	res, err := im.ImportFile(context.Background(), muestraPath)
	if err != nil {
		t.Fatal(err)
	}
	if res.Read != int64(len(rows)) || res.Rejected != 0 || !res.Balanced() {
		t.Fatalf("conteo de la muestra: leídos %d de %d, rechazados %d", res.Read, len(rows), res.Rejected)
	}
	if res.Lines-1 != res.Read {
		t.Errorf("wc -l (%d) − 1 distinto de registros (%d)", res.Lines, res.Read)
	}
	if res.MinDate < "2026-01-01" || res.MaxDate > "2026-08-31" {
		t.Errorf("fechas fuera de la cobertura: %s a %s", res.MinDate, res.MaxDate)
	}
}

func TestProfecoFiltroYLiga(t *testing.T) {
	rows := readMuestra(t)
	// Esperado, calculado aparte: cadena de las tres y municipio de las sucursales
	// (Benito Juárez solo con estado Quintana Roo).
	wantChain := map[string]bool{Fold("Wal-mart"): true, Fold("Bodega Aurrera"): true, Fold("Sam's Club"): true}
	passes := func(r muestraRow) bool {
		if !wantChain[r.chainFold] {
			return false
		}
		return r.munFold == "gustavo a madero" || (r.munFold == "benito juarez" && r.stateFold == "quintana roo")
	}
	var matching []muestraRow
	for _, r := range rows {
		if passes(r) {
			matching = append(matching, r)
		}
	}
	if len(matching) == 0 {
		t.Fatal("la muestra no trae registros de las cadenas en los municipios de prueba")
	}
	// Ligas: la mitad de las presentaciones filtradas, confirmadas; un tercio, propuestas.
	var sb strings.Builder
	seen := map[string]int{}
	confirmed, proposed := map[string]bool{}, map[string]bool{}
	for _, r := range matching {
		if _, ok := seen[r.productKey]; ok {
			continue
		}
		n := len(seen) + 1
		seen[r.productKey] = n
		switch n % 3 {
		case 0:
			fmt.Fprintf(&sb, "%s,%s,%s,%d,propuesto,,\n", r.rec[0], r.rec[1], r.rec[2], n)
			proposed[r.productKey] = true
		case 1:
			fmt.Fprintf(&sb, "%s,%s,%s,%d,confirmado,ximena,2026-10-11\n", r.rec[0], r.rec[1], r.rec[2], n)
			confirmed[r.productKey] = true
		}
	}
	links := mustLinks(t, sb.String())
	// Conteos esperados con «primero gana» por (presentación, tienda, día).
	storeOf := func(r muestraRow) int64 {
		switch {
		case r.nombre == "Walmart Sucursal Tepeyac":
			return 11
		case r.nombre == "Bodega Aurrera Express Sucursal Poniente 140":
			return 12
		case r.nombre == "Wal Mart" && r.munFold == "benito juarez":
			return 13
		}
		return map[string]int64{Fold("Wal-mart"): 1, Fold("Bodega Aurrera"): 2, Fold("Sam's Club"): 3}[r.chainFold]
	}
	type k3 struct {
		pres, store int64
		day         string
	}
	var wantSkipNoLink, wantLoadedConf, wantLoadedProp, wantDup int64
	seenPrice := map[k3]bool{}
	for _, r := range matching {
		switch {
		case !confirmed[r.productKey] && !proposed[r.productKey]:
			wantSkipNoLink++
			continue
		}
		key := k3{int64(seen[r.productKey]), storeOf(r), strings.ReplaceAll(r.rec[6], "/", "-")}
		if seenPrice[key] {
			wantDup++
			continue
		}
		seenPrice[key] = true
		if confirmed[r.productKey] {
			wantLoadedConf++
		} else {
			wantLoadedProp++
		}
	}
	// Modo revisión: cuentan las dos.
	mem := &MemoryPrices{}
	im, _ := NewImporter(Options{Corte: "2026", Stores: testStores(), Links: links, Mode: ModeReview, Prices: mem})
	res, err := im.ImportFile(context.Background(), muestraPath)
	if err != nil {
		t.Fatal(err)
	}
	if res.FilteredOut != int64(len(rows)-len(matching)) {
		t.Errorf("filtrados %d; esperaba %d", res.FilteredOut, len(rows)-len(matching))
	}
	if res.Skipped != wantSkipNoLink || res.SkippedReasons[SkipNoLink] != wantSkipNoLink {
		t.Errorf("sin liga %d; esperaba %d", res.Skipped, wantSkipNoLink)
	}
	if res.Loaded != wantLoadedConf || res.LoadedProposed != wantLoadedProp || res.Duplicates != wantDup {
		t.Errorf("cargados %d/%d dup %d; esperaba %d/%d dup %d", res.Loaded, res.LoadedProposed, res.Duplicates, wantLoadedConf, wantLoadedProp, wantDup)
	}
	if !res.Balanced() {
		t.Errorf("no cuadra: %+v", res)
	}
	if int64(len(mem.Rows)) != res.Loaded+res.LoadedProposed {
		t.Errorf("precios en memoria %d", len(mem.Rows))
	}
	var covTotal int64
	for _, c := range res.Coverage {
		covTotal += c.PriceCount
	}
	if covTotal != res.Loaded+res.LoadedProposed {
		t.Errorf("cobertura suma %d", covTotal)
	}
	for _, p := range mem.Rows {
		if !strings.HasPrefix(p.SourceRef, "QQP 2026 qqp-muestra.csv:") || p.PriceCents <= 0 || len(p.ObservedDate) != 10 {
			t.Fatalf("precio mal formado: %+v", p)
		}
	}
}

func TestProfecoLigaSinConfirmar(t *testing.T) {
	body := header +
		row("Atún", "Lata", "Dolores", "21.5", "2026/03/02", "Wal-mart", "Walmart Sucursal Tepeyac", "Calz", "Ciudad de México", "Gustavo a. Madero") +
		row("Arroz", "Bolsa", "SOS", "19", "2026/03/02", "Wal-mart", "Otra", "Calz", "Ciudad de México", "Gustavo a. Madero") +
		row("Frijol", "Bolsa", "Verde Valle", "30", "2026/03/02", "Wal-mart", "Otra", "Calz", "Ciudad de México", "Gustavo a. Madero")
	p := writeTmp(t, "pend.csv", body)
	links := mustLinks(t, "Atún,Lata,Dolores,1,confirmado,axel,2026-10-11\nArroz,Bolsa,SOS,2,propuesto,,\n")

	prod := &MemoryPrices{}
	im, _ := NewImporter(Options{Corte: "2026", Stores: testStores(), Links: links, Prices: prod}) // modo por omisión = producción
	r, err := im.ImportFile(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
	if r.Loaded != 1 || r.LoadedProposed != 0 || r.Skipped != 2 || r.SkippedReasons[SkipLinkPending] != 1 || r.SkippedReasons[SkipNoLink] != 1 {
		t.Errorf("producción: %+v", r)
	}
	if len(prod.Rows) != 1 || prod.Rows[0].PresentationID != 1 || prod.Rows[0].StoreID != 11 || prod.Rows[0].PriceCents != 2150 {
		t.Errorf("precio de producción: %+v", prod.Rows)
	}

	rev := &MemoryPrices{}
	im, _ = NewImporter(Options{Corte: "2026", Stores: testStores(), Links: links, Mode: ModeReview, Prices: rev})
	r, _ = im.ImportFile(context.Background(), p)
	if r.Loaded != 1 || r.LoadedProposed != 1 || r.Skipped != 1 || r.SkippedReasons[SkipNoLink] != 1 || !r.Balanced() {
		t.Errorf("revisión: %+v", r)
	}
	// La otra sucursal de Wal-mart no tiene liga a sucursal: el precio va a la cadena (id 1).
	if len(rev.Rows) != 2 || rev.Rows[1].StoreID != 1 {
		t.Errorf("precio de revisión: %+v", rev.Rows)
	}

	if _, err := NewImporter(Options{Mode: "otro"}); err == nil {
		t.Error("modo desconocido debía fallar")
	}
}

func TestProfecoDuplicadosYEstado(t *testing.T) {
	// Mismo producto, tienda y día: el primero gana. Benito Juárez de Quintana Roo
	// pasa; el de la Ciudad de México no (el estado distingue).
	body := header +
		row("Atún", "Lata", "Dolores", "21.5", "2026/03/02", "Wal-mart", "Wal Mart", "Costado Sams Club", "Quintana Roo", "Benito Juárez") +
		row("Atún", "Lata", "Dolores", "99.9", "2026/03/02", "Wal-mart", "Wal Mart", "Costado Sams Club", "Quintana Roo", "Benito Juárez") +
		row("Atún", "Lata", "Dolores", "22", "2026/03/03", "Wal-mart", "Wal Mart", "Costado Sams Club", "Quintana Roo", "Benito Juárez") +
		row("Atún", "Lata", "Dolores", "23", "2026/03/03", "Wal-mart", "Wal Mart", "X", "Ciudad de México", "Benito Juárez")
	p := writeTmp(t, "dup.csv", body)
	links := mustLinks(t, "Atún,Lata,Dolores,5,confirmado,axel,2026-10-11\n")
	mem := &MemoryPrices{}
	im, _ := NewImporter(Options{Corte: "2026", Stores: testStores(), Links: links, Prices: mem})
	total, files, err := im.Import(context.Background(), p, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 || total.Read != 8 {
		t.Fatalf("total: %+v", total)
	}
	// Primera pasada: 2 cargados (días 02 y 03), 1 duplicado, 1 filtrado. Segunda: todo duplicado.
	if files[0].Loaded != 2 || files[0].Duplicates != 1 || files[0].FilteredOut != 1 || files[1].Loaded != 0 || files[1].Duplicates != 3 {
		t.Errorf("por archivo: %+v / %+v", files[0], files[1])
	}
	if len(mem.Rows) != 2 || mem.Rows[0].PriceCents != 2150 || mem.Rows[0].StoreID != 13 {
		t.Errorf("filas: %+v", mem.Rows)
	}
	if total.Loaded != 2 || total.Duplicates != 4 || !total.Balanced() || len(total.Coverage) != 1 || total.Coverage[0].PriceCount != 2 {
		t.Errorf("acumulado: %+v", total)
	}
	if _, _, err := im.Import(context.Background(), filepath.Join(t.TempDir(), "nada.csv")); err == nil {
		t.Error("archivo inexistente debía fallar")
	}
}

func TestProfecoCancelacion(t *testing.T) {
	p := writeTmp(t, "c.csv", header+row("A", "P", "M", "1", "2026/01/02", "Wal-mart", "N", "D", "E", "Gustavo a. Madero"))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	im, _ := NewImporter(Options{Stores: testStores()})
	if _, err := im.ImportFile(ctx, p); !errors.Is(err, context.Canceled) {
		t.Errorf("esperaba cancelación, got %v", err)
	}
}

func TestReadLinksValidaciones(t *testing.T) {
	const h = "producto,presentacion,marca,presentacion_id,estado,confirmo,fecha_confirmacion\n"
	bad := map[string]string{
		"sin encabezado":      "",
		"falta columna":       "producto,presentacion\nA,B\n",
		"id inválido":         h + "A,B,C,x,propuesto,,\n",
		"id cero":             h + "A,B,C,0,propuesto,,\n",
		"producto vacío":      h + ",B,C,1,propuesto,,\n",
		"propuesta con firma": h + "A,B,C,1,propuesto,axel,\n",
		"confirma sin quién":  h + "A,B,C,1,confirmado,,2026-10-11\n",
		"confirma otro":       h + "A,B,C,1,confirmado,luis,2026-10-11\n",
		"fecha mala":          h + "A,B,C,1,confirmado,axel,11/10/2026\n",
		"estado raro":         h + "A,B,C,1,quizá,,\n",
		"repetida":            h + "A,B,C,1,propuesto,,\nA,B,C,2,propuesto,,\n",
		"columnas":            h + "A,B,C\n",
	}
	for name, body := range bad {
		if _, err := ReadLinks(strings.NewReader(body)); !errors.Is(err, ErrLinks) {
			t.Errorf("%s: esperaba ErrLinks, got %v", name, err)
		}
	}
	l, err := ReadLinks(strings.NewReader("\xEF\xBB\xBF" + h + "Atún,Lata,Dolores,1,CONFIRMADO,Axel,2026-10-11\n"))
	if err != nil || l.Len() != 1 {
		t.Fatalf("liga válida: %v", err)
	}
	// La liga se encuentra aunque cambien mayúsculas, acentos y guiones.
	if _, ok := l.Lookup(&Record{Producto: "ATUN", Presentacion: "lata", Marca: "dolores"}); !ok {
		t.Error("Lookup con Fold")
	}
	var nilLinks *Links
	if _, ok := nilLinks.Lookup(&Record{}); ok || nilLinks.Len() != 0 {
		t.Error("Links nulo")
	}
	p := writeTmp(t, "v.csv", h+"A,B,C,1,propuesto,,\n")
	if l, err := ReadLinksFile(p); err != nil || l.Len() != 1 {
		t.Errorf("ReadLinksFile: %v", err)
	}
	if _, err := ReadLinksFile(filepath.Join(t.TempDir(), "no.csv")); err == nil {
		t.Error("ReadLinksFile inexistente")
	}
}
