package main

import (
	"finanzas-hogar/internal/api"
	"finanzas-hogar/internal/api/ajustes"
	"finanzas-hogar/internal/api/avisos"
	"finanzas-hogar/internal/api/comidas"
	"finanzas-hogar/internal/api/comprasa"
	"finanzas-hogar/internal/api/comprasb"
	"finanzas-hogar/internal/api/creditos"
	gastos "finanzas-hogar/internal/api/gastos-inicio"
	"finanzas-hogar/internal/api/presupuesto"
	reportesapi "finanzas-hogar/internal/api/reportes"
	"finanzas-hogar/internal/api/tablero"
	"finanzas-hogar/internal/auth"
)

// areas junta los Register de cada área de la API (CONVENCIONES-backend.md §2 y §14): una línea por agente.
// Lo que no esté registrado responde 501; TestAreas_SinOperacionesEn501 exige todas las operaciones de C1 (221 en C1 v1.5).
var areas = []func(*api.Registry){
	auth.Register,        // auth-sistema (incluye sistema.Register: salud, contrato, respaldos, adjuntos, jobs)
	gastos.Register,      // api-gastos-inicio (paquete gastos en internal/api/gastos-inicio)
	presupuesto.Register, // api-presupuesto
	comidas.Register,     // api-comidas-salud
	comprasa.Register,    // api-compras-a
	comprasb.Register,    // api-compras-b
	creditos.Register,    // api-creditos
	ajustes.Register,     // api-ajustes
	avisos.Register,      // s7-avisos (cliente de Telegram con el token del archivo de MultiX; nunca se imprime)
	reportesapi.Register, // s7-pdf
	tablero.Register,     // tablero de control de Inicio y sus parámetros (C1 v1.5, K54; correccion-especificacion-r2)
}
