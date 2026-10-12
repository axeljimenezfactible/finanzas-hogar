package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"finanzas-hogar/internal/precios"
	"finanzas-hogar/internal/precios/ocr"
	"finanzas-hogar/internal/precios/profeco"
)

// datosS5 es la única carpeta de donde profeco_import y ocr_parse_tsv leen archivos (§15.1).
const datosS5 = "oraculo/fixtures/s5/_datos"

// Registro del motor de precios y tiendas (S5): el núcleo responde con precios.Eval; los
// cargadores de PROFECO y del OCR se atienden aquí con sus subpaquetes (S5-interfaz §15.2).
func init() {
	register("s5", func(input json.RawMessage, q Query) (any, string) {
		var out any
		var err error
		switch q.Func {
		case "profeco_import":
			out, err = s5Profeco(input, q.Args)
		case "ocr_parse_tsv":
			out, err = s5OCR(q.Args)
		default:
			out, err = precios.Eval(input, q.Func, q.Args)
		}
		if err != nil {
			return nil, s5Code(err)
		}
		return out, ""
	})
}

func s5Code(err error) string {
	switch {
	case errors.Is(err, profeco.ErrProfecoHeader):
		return "profeco_header"
	case errors.Is(err, profeco.ErrFileTooLarge):
		return "file_too_large"
	}
	return precios.Code(err)
}

func s5Args(raw json.RawMessage, v any) error {
	if len(raw) == 0 || string(raw) == "null" {
		raw = json.RawMessage("{}")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("%w: %v", precios.ErrUnknownField, err)
	}
	return nil
}

// s5OCR parsea un TSV de Tesseract de la fixture (no corre Tesseract).
func s5OCR(args json.RawMessage) (any, error) {
	var a struct {
		TSVFile *string `json:"tsv_file"`
	}
	if err := s5Args(args, &a); err != nil {
		return nil, err
	}
	if a.TSVFile == nil {
		return nil, precios.ErrMissingField
	}
	p, err := precios.RepoPath(*a.TSVFile, datosS5)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", precios.ErrInvalidArgument, err)
	}
	defer f.Close()
	lines, err := ocr.ParseTSV(f, ocr.ProfileGeneric)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", precios.ErrInvalidArgument, err)
	}
	type line struct {
		Description   string `json:"description"`
		QuantityMilli *int64 `json:"quantity_milli"`
		AmountCents   *int64 `json:"amount_cents"`
		Confidence    int    `json:"confidence"`
	}
	out := struct {
		Lines []line `json:"lines"`
	}{Lines: []line{}}
	for _, l := range lines {
		out.Lines = append(out.Lines, line{l.Descripcion, l.CantidadMilli, l.ImporteCents, l.Confianza})
	}
	return out, nil
}

// s5Profeco carga CSV de muestra con el importador de PROFECO, filtrando con las tiendas del
// Input. Los ids del importador son posiciones (tienda en stores, presentación en el catálogo,
// desde 1); presentacion_id del CSV de ligas puede traer el id o el código de la presentación.
func s5Profeco(input, args json.RawMessage) (any, error) {
	var a struct {
		Files     []string `json:"files"`
		LinksFile *string  `json:"links_file"`
		Mode      *string  `json:"mode"`
	}
	if err := s5Args(args, &a); err != nil {
		return nil, err
	}
	if a.LinksFile == nil || a.Mode == nil {
		return nil, precios.ErrMissingField
	}
	in, err := precios.ParseInput(input)
	if err != nil {
		return nil, err
	}
	b := in.Book
	presID := map[string]int64{}
	for i, p := range b.Catalog.Presentations {
		presID[p.Code] = int64(i + 1)
	}
	var set profeco.StoreSet
	storeID := map[string]int64{}
	for i, st := range b.Stores {
		storeID[st.Code] = int64(i + 1)
	}
	for i, st := range b.Stores {
		switch st.Kind {
		case precios.KindChain:
			set.Chains = append(set.Chains, profeco.Chain{ID: int64(i + 1), Name: st.Name, ProfecoName: st.ProfecoName})
		case precios.KindBranch:
			set.Branches = append(set.Branches, profeco.Branch{ID: int64(i + 1), ChainID: storeID[st.Parent],
				Name: st.Name, ProfecoName: st.ProfecoName, Municipality: st.Municipality})
		}
	}
	lp, err := precios.RepoPath(*a.LinksFile, datosS5)
	if err != nil {
		return nil, err
	}
	links, err := s5Links(lp, presID)
	if err != nil {
		return nil, err
	}
	var files []string
	asked := map[string]string{} // ruta resuelta → ruta tal como llega en files[] (rejections[].file)
	for _, f := range a.Files {
		p, err := precios.RepoPath(f, datosS5)
		if err != nil {
			return nil, err
		}
		files = append(files, p)
		asked[p] = f
	}
	im, err := profeco.NewImporter(profeco.Options{Corte: "fixture", Stores: set, Links: links,
		Mode: profeco.Mode(*a.Mode), Prices: &profeco.MemoryPrices{}, Today: b.Today.Format("2006-01-02")})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", precios.ErrInvalidArgument, err)
	}
	r, _, err := im.Import(context.Background(), files...)
	if err != nil {
		return nil, err
	}
	type cov struct {
		Chain        string `json:"chain"`
		Municipality string `json:"municipality"`
		PriceCount   int64  `json:"price_count"`
	}
	type rej struct {
		File   string `json:"file"`
		Row    int    `json:"row"`
		Reason string `json:"reason"`
	}
	out := struct {
		Read           int64 `json:"read"`
		Rejected       int64 `json:"rejected"`
		FilteredOut    int64 `json:"filtered_out"`
		Skipped        int64 `json:"skipped"`
		Loaded         int64 `json:"loaded"`
		LoadedProposed int64 `json:"loaded_proposed"`
		Duplicates     int64 `json:"duplicates"`
		Coverage       []cov `json:"coverage"`
		Rejections     []rej `json:"rejections"`
	}{r.Read, r.Rejected, r.FilteredOut, r.Skipped, r.Loaded, r.LoadedProposed, // §11: Read = Rejected + FilteredOut + Skipped + Loaded + LoadedProposed + Duplicates

		r.Duplicates, []cov{}, []rej{}}
	for _, c := range r.Coverage {
		out.Coverage = append(out.Coverage, cov{c.Chain, c.Municipality, c.PriceCount})
	}
	for _, x := range r.Rejections {
		out.Rejections = append(out.Rejections, rej{asked[x.File], x.Row, x.Reason})
	}
	return out, nil
}

// s5Links lee el CSV de ligas y traduce presentacion_id escrito como código a su posición.
func s5Links(path string, presID map[string]int64) (*profeco.Links, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", precios.ErrInvalidArgument, err)
	}
	rows, err := csv.NewReader(bytes.NewReader(raw)).ReadAll()
	if err != nil || len(rows) == 0 {
		return nil, fmt.Errorf("%w: ligas %v", precios.ErrInvalidArgument, err)
	}
	col := -1
	for i, h := range rows[0] {
		if h == "presentacion_id" {
			col = i
		}
	}
	for _, r := range rows[1:] {
		if col < 0 || col >= len(r) {
			break
		}
		if _, err := strconv.ParseInt(r[col], 10, 64); err != nil {
			id, ok := presID[r[col]]
			if !ok { // PROVISIONAL, bloqueo S5-motor-oraculo-reglas pregunta 5: el oráculo omite el registro (skipped)
				return nil, fmt.Errorf("%w: %s", precios.ErrUnknownPresentation, r[col])
			}
			r[col] = strconv.FormatInt(id, 10)
		}
	}
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.WriteAll(rows); err != nil {
		return nil, err
	}
	l, err := profeco.ReadLinks(&buf)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", precios.ErrInvalidArgument, err)
	}
	return l, nil
}
