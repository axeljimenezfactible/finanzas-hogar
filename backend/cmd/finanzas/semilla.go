package main

import seed "finanzas-hogar/internal/semilla"

// finanzas semilla carga datos/semilla/ con internal/semilla (idempotente, una transacción). Las contraseñas iniciales
// salen del entorno (FINANZAS_SEED_PASSWORD_AXEL y _XIMENA) o de un archivo fuera del repositorio
// (FINANZAS_SEED_PASSWORD_FILE); nunca de un valor por omisión.
func init() { seedRunner = seed.Run }
