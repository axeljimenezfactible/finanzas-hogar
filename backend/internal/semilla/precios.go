package semilla

// Precios de Walmart de la semilla REAL (WF-3b, agente semilla-c3, 12-oct-2026).
//
// DECISIÓN PROVISIONAL DEL ORQUESTADOR hasta la puerta 4 (bloqueo WF3a-bis-semilla-dev, punto 7 = sí): los precios de
// Walmart del carrito comparado del 10-oct-2026 de data.json (fuente «manual», dato real de Axel) entran a la semilla
// real, porque la comparativa de canasta Bodega contra Walmart (S5) se reproduce con ellos: +138.50 en la canasta
// completa y +100.50 en los artículos idénticos. Antes solo los cargaba el overlay de desarrollo.
//
// Los precios de Bodega Aurrera NO se cargan aquí: la compra base de la comparativa sale de los renglones del propio
// ticket (expense_item), y los precios unitarios de Bodega siguen siendo del overlay de desarrollo.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// walmartFile son los precios de Walmart del carrito comparado del 10-oct (datos/semilla/data.json).
type walmartFile struct {
	Fecha          string `json:"fecha"`
	ItemsIdenticos []struct {
		Producto     string `json:"producto"`
		Cantidad     int64  `json:"cantidad"`
		WalmartTotal money  `json:"walmart_total"`
	} `json:"items_identicos"`
	Sustitutos []struct {
		Walmart      string `json:"walmart"`
		CantW        int64  `json:"cant_w"`
		WalmartTotal money  `json:"walmart_total"`
	} `json:"sustitutos"`
}

// precios carga los 36 precios de Walmart del carrito comparado de data.json. Sigue el mismo criterio de siempre:
// un renglón cuyo importe no se divide exacto entre las piezas no tiene precio unitario y no se inventa (se avisa).
func (l *loader) precios() error {
	b, err := os.ReadFile(filepath.Join(l.opt.Dir, "data.json"))
	if err != nil {
		return err
	}
	var w walmartFile
	if err := json.Unmarshal(b, &w); err != nil {
		return fmt.Errorf("semilla: data.json: %w", err)
	}
	wstore, ok, err := l.find(`SELECT s.id FROM store s JOIN store c ON c.id = s.parent_id WHERE s.kind = 'branch' AND c.name = 'Walmart' ORDER BY s.id LIMIT 1`)
	if err != nil {
		return err
	}
	if !ok {
		l.rep.Notices = append(l.rep.Notices, "No hay sucursal de Walmart sembrada: no se cargaron precios de Walmart.")
		return nil
	}
	omitted := 0
	put := func(name string, qty int64, total money) error {
		id, ok := l.pres[name]
		if !ok || qty <= 0 || int64(total)%qty != 0 {
			omitted++
			return nil
		}
		return l.price("precios Walmart (carrito 10-oct)", id, wstore, int64(total)/qty, w.Fecha, "manual", "data.json carrito comparado 10-oct-2026", 0)
	}
	for _, it := range w.ItemsIdenticos {
		if err := put(it.Producto, it.Cantidad, it.WalmartTotal); err != nil {
			return err
		}
	}
	for _, s := range w.Sustitutos {
		if err := put(s.Walmart, s.CantW, s.WalmartTotal); err != nil {
			return err
		}
	}
	if omitted > 0 {
		l.rep.Notices = append(l.rep.Notices, fmt.Sprintf("Walmart: %d renglón(es) del carrito no se cargaron (sin presentación en el catálogo o sin precio unitario exacto).", omitted))
	}
	return nil
}
