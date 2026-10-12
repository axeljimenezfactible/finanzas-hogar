// Package ocr lee tickets térmicos con Tesseract en local (S5, K24, S5-06).
//
// Flujo: preprocesado opcional con ImageMagick (escala de grises, enderezado
// y umbral local) -> tesseract 5.5.3 con idioma spa en formato TSV -> análisis
// de renglones {descripcion, cantidad, importe, confianza}.
//
// Regla suprema: la salida es siempre un BORRADOR que una persona confirma.
// Nada de este paquete escribe precios; PricesFromOCR devuelve ErrNotConfirmed
// mientras la persona no confirme (S5-06).
//
// Configuración de Tesseract (documentada en docs/estado/wf1/ocr.md):
//
//	tesseract <imagen> stdout -l spa --psm 6 --oem 1 -c preserve_interword_spaces=1 tsv
//
// psm 6 = un solo bloque uniforme de texto (un ticket es una columna de
// renglones); oem 1 = solo la red LSTM; preserve_interword_spaces conserva los
// huecos entre descripción e importe.
//
// Dependencias: solo biblioteca estándar; tesseract y magick por os/exec sin shell.
package ocr
