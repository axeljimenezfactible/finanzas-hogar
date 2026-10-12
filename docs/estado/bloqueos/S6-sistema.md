# Bloqueo S6-sistema · decisiones abiertas de auth-sistema (PROVISIONAL)

12-oct-2026 · agente auth-sistema (WF-3a). Cada punto está implementado con el valor indicado y marcado PROVISIONAL en el código.

| # | PREGUNTA | Valor PROVISIONAL | Opciones |
|---|---|---|---|
| 1 | C1 permite `auth.password_min_length` = 0 (y `auth.max_failed_attempts` = 0); K6 pide ≥ 12 | Se respeta C1 (mínimo 0) | Subir el `minimum` de C1 (generador + disparador) a 12 y 1 |
| 2 | ¿Cambiar la contraseña cierra las otras sesiones? (C1 lista `device_session` en x-tablas) | No las cierra | Cerrar todas salvo la actual |
| 3 | ¿El usuario distingue mayúsculas? (los teléfonos ponen «Axel») | Exacto, como está en la base | Comparar sin mayúsculas |
| 4 | Estado de la semilla (`getSeedStatus`): no hay tabla de semillas | Cargada si `audit_log` tiene altas con `request_id = 'seed:<código>'`; `row_count` = esas altas. `internal/semilla` debe seguir esa convención | Tabla nueva (migración) |
| 5 | Ruta de la CA para `getCaCertificate` | `serve --ca` / `FINANZAS_CA` → `sistema.CACertPath`; sin ruta, 503 | Campo `CACertPath` en `api.Config` (backend-base) |
| 6 | IA: modelo de Ollama, instrucciones, dónde vive la llave de Claude | Adaptadores `ia.Ollama` (ping a `/api/tags`) e `ia.Claude` (llave en un archivo fuera del repositorio) existen y están apagados; OCR y recetas responden «no configurado»; `serve` todavía no instala ninguno; `is_configured` = hay adaptador instalado | Decidir modelo y llave antes de encenderla |
| 7 | Orden por defecto de `listJobs` (C1 no lo dice) | `code` | `next_run_at` |
| 8 | Cabeceras de seguridad (K31) | Las de backend-base (bloqueo S6-valores-operativos-base); descargas con `Content-Disposition: attachment` y `Cache-Control: private, no-store` | — |
| 9 | Nombre `finanzas.local` en Android y la IP de la LAN | `ops/mkcert.md` con `192.168.50.X`; `serve` acepta solo 127.0.0.1 y 192.168.50.0/24 | Entrada en el DNS del router; IP fija de WF-4 |
