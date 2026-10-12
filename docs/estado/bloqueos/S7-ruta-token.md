# Bloqueo S7-ruta-token · ruta exacta de la clave del token en openclaw.json

Abierto por s7-avisos (WF-3a) el 12-oct-2026. Criterio: S7-06, K30. Estado: **ABIERTO (PREGUNTA a Axel)**.

## HECHO
- K30 y S7-06 nombran la clave `bots.multix.botToken` de `~/.openclaw/openclaw.json`.
- La nota de operación de Axel para el mismo bot (memoria `reference_telegram_inbound`) lee `channels.telegram.bots.multix.botToken`. No se abrió el archivo real para comprobarlo (K30: el token no se lee fuera del servicio).

## PROVISIONAL
`telegram.TokenPaths` busca primero la ruta literal de K30 y, si no existe, `channels.telegram.bots.multix.botToken`. Si ninguna existe, el envío falla con «no trae la clave bots.multix.botToken», sin imprimir el contenido.

## Falta
Que Axel confirme cuál es la ruta real. Opciones: A) dejar las dos (sin costo; recomendada hasta WF4-01); B) dejar solo la que confirme. Impacto: si la ruta real no estuviera en la lista, S7-07 fallaría con un mensaje claro y sin riesgo para el token.
