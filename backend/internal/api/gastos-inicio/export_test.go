package gastos

import "finanzas-hogar/internal/precios/ocr"

// SetOCR cambia el motor de OCR por uno falso en las pruebas y devuelve cómo restaurarlo.
func SetOCR(r ocr.Recognizer) func() {
	old := ocrEngine
	ocrEngine = r
	return func() { ocrEngine = old }
}
