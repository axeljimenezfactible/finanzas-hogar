package ocr

import (
	"errors"
	"testing"
)

// S5-06: TestOCRSinConfirmar
func TestOCRSinConfirmar(t *testing.T) {
	lines := []Line{
		ParseLine("LECHE 28.50", ProfileGeneric),
		ParseLine("2 ARROZ 58.00", ProfileGeneric),
	}
	items, err := PricesFromOCR(lines, false)
	if !errors.Is(err, ErrNotConfirmed) || items != nil {
		t.Fatalf("sin confirmar: %v %v", items, err)
	}
	items, err = PricesFromOCR(nil, false)
	if !errors.Is(err, ErrNotConfirmed) || items != nil {
		t.Fatalf("nil sin confirmar: %v %v", items, err)
	}
}

func TestPricesFromOCRConfirmado(t *testing.T) {
	lines := []Line{
		ParseLine("LECHE 28.50", ProfileGeneric),
		ParseLine("2 ARROZ 58.00", ProfileGeneric),
		ParseLine("TIENDA SIN IMPORTE", ProfileGeneric),
		ParseLine("TOTAL 86.50", ProfileGeneric),
		{Descripcion: " ", ImporteCents: i64(100)},
	}
	items, err := PricesFromOCR(lines, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("items = %+v", items)
	}
	if items[0].Descripcion != "LECHE" || items[0].CantidadMilli != nil || items[0].ImporteCents != 2850 {
		t.Errorf("0: %+v", items[0])
	}
	if items[1].Descripcion != "ARROZ" || *items[1].CantidadMilli != 2000 || items[1].ImporteCents != 5800 {
		t.Errorf("1: %+v", items[1])
	}
}
