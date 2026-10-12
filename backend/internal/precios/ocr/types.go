package ocr

import (
	"errors"
	"strings"
)

// Constantes del contrato (S5-interfaz §12).
const (
	RequiredVersion = "5.5.3"
	Lang            = "spa"

	// DefaultTimeoutSeconds es el tiempo máximo por imagen: 30 s, FIRME (Q2-76,
	// K50). La clave ocr.timeout_seconds de app_setting existe desde la
	// migración 00004 y este es su valor inicial.
	DefaultTimeoutSeconds = 30
)

// Errores del paquete.
var (
	// ErrOCRUnavailable: falta tesseract, su versión no es la requerida o no
	// tiene el idioma spa.
	ErrOCRUnavailable = errors.New("ocr: tesseract no disponible (se requiere 5.5.3 con idioma spa)")
	// ErrNotConfirmed: un OCR sin confirmar no crea ningún precio (S5-06).
	ErrNotConfirmed = errors.New("ocr: el borrador no está confirmado por una persona")
	// ErrBadImage: la ruta no es una imagen legible (.jpg, .jpeg o .png).
	ErrBadImage = errors.New("ocr: imagen no válida")
)

// Line es un renglón leído. Las claves JSON están en español como las del
// contrato (§12); la fixture s5-ocr-parseo las traduce a inglés.
type Line struct {
	Descripcion   string `json:"descripcion"`
	CantidadMilli *int64 `json:"cantidad_milli"` // milésimas de la unidad impresa; nil si no se imprimió
	Unidad        string `json:"unidad,omitempty"`
	ImporteCents  *int64 `json:"importe_cents"` // nil: la persona debe descartarlo o completarlo
	Confianza     int    `json:"confianza"`     // 0 a 100, promedio HALF_UP de las palabras
	Texto         string `json:"texto"`         // texto leído sin interpretar
	EsResumen     bool   `json:"es_resumen"`    // TOTAL, SUBTOTAL, IVA, CAMBIO...: no es un producto
}

// ConfirmedItem es un renglón ya confirmado por una persona.
type ConfirmedItem struct {
	Descripcion   string
	CantidadMilli *int64
	ImporteCents  int64
}

// PricesFromOCR convierte renglones en renglones de gasto/precios ticket_ocr
// solo si una persona los confirmó. Omite los renglones sin importe y los de
// resumen (total, iva...). Con confirmado = false devuelve ErrNotConfirmed y
// ningún renglón (S5-06).
func PricesFromOCR(lineas []Line, confirmado bool) ([]ConfirmedItem, error) {
	if !confirmado {
		return nil, ErrNotConfirmed
	}
	items := make([]ConfirmedItem, 0, len(lineas))
	for _, l := range lineas {
		if l.ImporteCents == nil || l.EsResumen || strings.TrimSpace(l.Descripcion) == "" {
			continue
		}
		items = append(items, ConfirmedItem{
			Descripcion:   l.Descripcion,
			CantidadMilli: l.CantidadMilli,
			ImporteCents:  *l.ImporteCents,
		})
	}
	return items, nil
}
