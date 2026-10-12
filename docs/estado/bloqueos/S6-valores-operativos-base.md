# S6-valores-operativos-base · Valores operativos PROVISIONALES de backend-base (WF-3a)

- **Fase:** S6 (WF-3a «Construcción»), agente backend-base, 12-oct-2026.
- **Estado:** ABIERTO hasta la puerta 4 (K31). El código ya funciona con la propuesta; cambiarla es una constante o una opción.
- **Qué falta (decisión de Axel):**
  1. **Hora y retención del respaldo diario** (S6-09): propuesta `backup_daily` a las **03:15** (America/Mexico_City), conservar **30** respaldos diarios. Por qué: a esa hora nadie captura; queda después del resumen diario (21:30) y antes del reporte mensual (08:00); 30 copias = un mes para volver atrás y, con la base de ~1 MB de hoy, ocupan poco. Dónde: `backend/cmd/finanzas/main.go` (`backupSchedule`, `backupKeep`); ops/runbook.md (WF-4) debe copiarlo.
  2. **Parámetros de argon2id**: propuesta t=3, m=64 MiB, p=4, sal 16 B, clave 32 B (segunda opción recomendada de RFC 9106 §4; ~0.1–0.2 s por inicio de sesión en la máquina de casa). Dónde: `dominio.DefaultArgon2`.
  3. **Cabeceras** (K31 «cabeceras»): `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`, `X-Frame-Options: DENY`; `Cache-Control: no-store` en la API; CSP en la PWA `default-src 'self'; img-src 'self' data: blob:; style-src 'self' 'unsafe-inline'; connect-src 'self'; font-src 'self'; worker-src 'self'; manifest-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'` ('unsafe-inline' en estilos por los estilos en línea de Radix y ECharts). Sin HSTS mientras el certificado sea de mkcert (un cambio de certificado bloquearía la PWA). Dónde: `api.securityHeaders`.
  4. **Límites de cuerpo**: 8 MiB por cuerpo JSON (413 al pasarse); hasta 32 MiB en memoria por cuerpo con Idempotency-Key, el resto a un temporal (para cargas de PROFECO de hasta 1 GB). Dónde: `api.Config.MaxJSONBytes`, `IdemMemBytes`.
  5. **Job `idempotency_purge`** (C1 dice «un job la purga» sin nombre ni horario): cada hora al minuto 5.
  6. **Rutas por omisión** de `finanzas serve` (`datos/finanzas.db`, `datos/app`, `datos/respaldos`, `127.0.0.1:8443`): las fija WF-4 en `ops/runbook.md` con `FINANZAS_*`; el binario rechaza 0.0.0.0, :: y «:8443» (S6-14).
  7. **Campo desconocido en un cuerpo JSON = 400** `bad_request` (C1 §3 no lo dice; los esquemas tienen `additionalProperties: false` y se trata como formato, igual que un tipo incorrecto).
- **Opciones:** (a) aprobar la propuesta; (b) cambiar valores puntuales; (c) mover 1, 4 y 5 a `app_setting` (requiere migración nueva y C1, por eso no se hizo).
- **Impacto:** ninguno en cifras, esquema ni endpoints; solo operación y seguridad. Se relaciona con `S5-agente-externo` (credencial del agente externo y a nombre de quién quedan sus claves de idempotencia: hoy sin verificador las 4 operaciones dan 401).
