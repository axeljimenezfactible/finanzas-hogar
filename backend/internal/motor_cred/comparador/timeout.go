package comparador

import (
	"math"
	"time"

	mc "finanzas-hogar/internal/motor_cred"
)

// TimeoutFromSetting convierte comparator.timeout_seconds (K43, K44; inicial 5) en duración:
// < 0 → ErrInvalid; 0 vence de inmediato (S3-interfaz.md §8.7).
func TimeoutFromSetting(seconds int64) (time.Duration, error) {
	if seconds < 0 || seconds > math.MaxInt64/int64(time.Second) {
		return 0, invalid("comparator.timeout_seconds = %d", seconds)
	}
	return time.Duration(seconds) * time.Second, nil
}

// checkCtx traduce un contexto vencido o cancelado a ErrTimeout (503 unavailable por ResponseFor).
func checkCtx(ctx interface{ Err() error }) error {
	if ctx.Err() != nil {
		return mc.ErrTimeout
	}
	return nil
}
