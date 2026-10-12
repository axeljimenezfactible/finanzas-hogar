// Command ocrmedir mide S5-06: compara el OCR contra la transcripción a mano.
//
//	go run ./internal/precios/ocr/cmd/ocrmedir \
//	    -tickets ../datos/semilla/tickets -esperado ../datos/semilla/tickets-esperado
//
// Con -sintetico mide imágenes de prueba (no cuenta para S5-06).
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"finanzas-hogar/internal/precios/ocr"
)

func main() {
	tickets := flag.String("tickets", "../datos/semilla/tickets", "carpeta con las fotos 01..05")
	esperado := flag.String("esperado", "../datos/semilla/tickets-esperado", "carpeta con las transcripciones NN.csv")
	synth := flag.Bool("sintetico", false, "imágenes de prueba generadas (no cuenta para S5-06)")
	pre := flag.Bool("preprocesar", true, "preprocesar con ImageMagick")
	perfil := flag.String("perfil", "generico", "perfil de análisis: generico o con-codigo")
	flag.Parse()

	eng := ocr.New()
	eng.Preprocess = *pre
	switch *perfil {
	case "generico":
	case "con-codigo":
		eng.Profile = ocr.ProfileWithCode
	default:
		fmt.Fprintln(os.Stderr, "error: perfil desconocido:", *perfil)
		os.Exit(2)
	}
	eng.Warn = func(m string) { fmt.Fprintln(os.Stderr, "aviso:", m) }
	rep, err := ocr.Measure(context.Background(), eng, ocr.MeasureOptions{
		TicketsDir: *tickets, ExpectedDir: *esperado, Synthetic: *synth,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Print(rep.Format())
}
