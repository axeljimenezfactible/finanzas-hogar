package main

import "testing"

// «finanzas semilla» ya no es un esqueleto: internal/semilla instala el cargador.
func TestSemillaInstalada(t *testing.T) {
	if seedRunner == nil {
		t.Fatal("seedRunner no está instalado: «finanzas semilla» respondería que falta internal/semilla")
	}
}
