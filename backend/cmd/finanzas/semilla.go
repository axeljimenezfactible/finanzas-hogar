package main

import (
	"context"
	"database/sql"

	seed "finanzas-hogar/internal/semilla"
)

// devSeedRunner es la semilla de DESARROLLO (finanzas semilla --desarrollo): un overlay aparte, nunca en producción.
// Lo instala internal/semilla (desarrollo.go); la bandera la lee semilla() en main.go.
var devSeedRunner func(ctx context.Context, db *sql.DB, dir string, conGastoFijo bool) error

// finanzas semilla carga datos/semilla/ con internal/semilla (idempotente, una transacción). Las contraseñas iniciales
// salen del entorno (FINANZAS_SEED_PASSWORD_AXEL y _XIMENA) o de un archivo fuera del repositorio
// (FINANZAS_SEED_PASSWORD_FILE); nunca de un valor por omisión.
func init() { seedRunner, devSeedRunner = seed.Run, seed.RunDev }
