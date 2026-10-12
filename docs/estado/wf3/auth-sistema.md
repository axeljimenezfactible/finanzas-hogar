# auth-sistema · WF-3a «Construcción»

12-oct-2026 · agente **auth-sistema** (Opus). Alcance: autenticación, 19 operaciones de sistema, respaldos, adjuntos, IA apagada, TLS y escucha. Las 27 operaciones del reparto (CONVENCIONES-backend §14) están implementadas y cada una tiene su prueba `TestOp_<operationId>` que valida contra `openapi.yaml`.

## HECHO

- **`internal/auth`** (`auth.Register`, la única línea que va en `areas.go`; registra también `sistema.Register`): `login` (argon2id, verificación falsa cuando la cuenta no existe para no revelar usuarios, 401 `invalid_credentials` genérico, 429 `too_many_attempts` con `retry_at` y «Vuelve a intentar a las HH:MM» aunque la contraseña sea correcta; el intento bloqueado no se anota), `logout`, `changeOwnPassword` (mínimo vigente `auth.password_min_length`, 422 `invalid_current_password`, audit sin hash), `listUsers` (nunca `password_hash`), `listSessions` (`is_current`, filtro `person_id`, orden `-last_seen_at`), `revokeSession` y `revokeSessions` (los tres alcances de C1, audit por sesión sin `token_hash`), `getCurrentSession`. Cookie `fh_session` HttpOnly; Secure; SameSite=Strict; 90 días renovados con cada uso (middleware de backend-base).
- **`internal/api/sistema`**: `getHealth` (consulta la base), `getOpenApiDocument`, ajustes de seguridad y de datos (`dominio.ApplyParams`: tipo, sección, rango, audit), IA (`listAiProviders`, `updateAiProvider`), `getSystemStatus` (versión por `-ldflags`, inicio, tamaño del archivo), `getCaCertificate` (solo bloques CERTIFICATE; nunca la llave), `getSeedStatus` (PROVISIONAL), `listJobs`, `updateJob` (cron validado con `jobs.ParseSchedule`, recalcula `next_run_at`), respaldos y adjuntos.
- **`internal/respaldo`**: `Run` (VACUUM INTO versionado `finanzas-AAAAMMDD-HHMMSS-<tipo>-<id>.db` en la carpeta configurable `BackupDir`, integrity_check, restauración de prueba y comparación de conteos con el origen) y `Restore` (valida estado, huella e integridad; respaldo `pre_migration` previo; API de respaldo de SQLite; migraciones; repone `backup` y `device_session`; integrity_check).
- **`internal/adjuntos`**: disco local `DataDir/adjuntos/AAAA/MM/<32 hex>.<ext>`, tipo detectado por contenido, tamaño ≤ `attachment.max_size_mb` × 1,048,576 (PROVISIONAL), rutas que no salen de la carpeta (sin `..`, sin absolutas, sin enlaces simbólicos), nombre original con `/`, `\` o `..` rechazado (422).
- **`internal/ia/adaptadores.go`**: `ia.Ollama` e `ia.Claude` con cliente HTTP inyectable, atados al interruptor `ai.*_enabled` (apagado por defecto): apagados devuelven `ErrDisabled` sin tocar la red; `Registry.Check` solo consulta si está encendido.
- **`cmd/finanzas`**: `tls.go` (`parseAddrs`: solo 127.0.0.1 y 192.168.50.0/24; `loadTLS` TLS ≥ 1.2; `serveAll`: mismo manejador en cada dirección). `serve` acepta varias direcciones (`FINANZAS_ADDR=192.168.50.X:8443,127.0.0.1:8443`), exige `--cert/--key` (o `--dev-http` solo en 127.0.0.1) y toma `--ca`/`FINANZAS_CA`. La PWA ya la sirve el binario desde `frontend/dist` (backend-base); se comprobó.
- **`ops/mkcert.md`**: instalación, emisión (integrador), arranque, CA en los dos Android, verificación.

## Archivos

Nuevos: `backend/internal/auth/{auth.go,auth_test.go}`, `backend/internal/api/sistema/{sistema.go,respaldos.go,adjuntos.go,sistema_test.go}`, `backend/internal/respaldo/{respaldo.go,respaldo_test.go}`, `backend/internal/adjuntos/{adjuntos.go,adjuntos_test.go}`, `backend/internal/ia/adaptadores.go`, `backend/cmd/finanzas/{tls.go,tls_test.go}`, `ops/mkcert.md`, `docs/estado/bloqueos/{S6-adjuntos.md,S6-respaldos.md,S6-sistema.md}`, este reporte. Cambiado: `backend/cmd/finanzas/main.go` (solo `serve`: varias direcciones, TLS obligatorio, `--ca`, `--dev-http`; import `sistema`).

## Verificación

- **S6-01**: `go version` = go1.27.2, `go.mod` = `go 1.27`, `go build ./...` y `go vet ./...` sin diagnósticos; `gofmt -l .` lista solo `internal/reportes/vistas.go` (de s7-pdf, no mío). CUMPLE en lo mío.
- **S6-02**: `go test ./... -count=1`: 0 fallas. Con umbral: dominio 87.6 %, motor_fin 89.2 %, motor_cred 94.0 %, recetas 93.2 %, precios 88.8 %. Sin umbral: api 79.9 % (incluye subpaquetes; sistema 82.2 %), **auth 83.6 %**, jobs 93.0 %, reportes 22.3 %, telegram 80.0 %. Extra: respaldo 75.0 %, adjuntos 83.3 %.
- **S6-03**: `build.py --check` = 0 (23 archivos idénticos); `make verify-wf3-contrato` OK. Tabla operación → prueba (27/27): login→TestOp_login, logout→TestOp_logout, changeOwnPassword→TestOp_changeOwnPassword, listUsers→TestOp_listUsers, listSessions→TestOp_listSessions, revokeSession→TestOp_revokeSession, revokeSessions→TestOp_revokeSessions, getCurrentSession→TestOp_getCurrentSession (internal/auth); getHealth, getOpenApiDocument, listBackups, createBackup, restoreBackup, getSeedStatus, getSecuritySettings, updateSecuritySettings, getDataSettings, updateDataSettings, listAiProviders, updateAiProvider, getSystemStatus, getCaCertificate, uploadExpenseAttachment, downloadAttachment, deleteAttachment, listJobs, updateJob → `TestOp_<id>` (internal/api/sistema). Todas cubren 2xx y sus 4xx con `c.Do` (valida contra C1).
- **S6-04**: errores con el esquema único, mensajes en español de México, página ≤ 200 (`limit=201` → 400), Idempotency-Key en `createBackup` y `uploadExpenseAttachment` (repetición con `Idempotent-Replayed`, sin clave → 400). CUMPLE en lo mío.
- **S6-05**: `go test ./internal/auth/... -count=1 -v`: 16 PASS. `TestOp_login` (Set-Cookie con HttpOnly, Secure, SameSite=Strict, Max-Age=7776000; `substr(password_hash,1,10)` = `$argon2id$` en las 2 cuentas), `TestOp_changeOwnPassword` (11 caracteres → 422, 12 → 204), `TestCuentas_DosYKarenSinCuenta` (2 cuentas; Karen 401; `POST /api/v1/users` 405; la base rechaza su cuenta), `TestPermisos_MismasOperacionesAmbasCuentas` («212 operaciones con el mismo resultado para las dos cuentas», ninguna 403; las 4 que cierran la sesión o restauran se prueban aparte). CUMPLE.
- **S6-06**: `go test ./internal/auth/... -run Intentos -count=1 -v`: (a) 5 fallos en 40 s → minuto 14 = 429 «a las 12:15», minuto 16 = 200; (b) fallos 0, 1, 2, 3, 14 → 14.5 = 429, 15.5 = 200; (c) por usuario y por IP (C1-datos §2): mismo usuario desde otra IP = 429, otra cuenta desde la misma IP = 429, otra cuenta desde otra IP = 200. 3 PASS. CUMPLE.
- **S6-07**: `go test ./internal/auth/... -run Sesion -count=1 -v`: 2 dispositivos, cerrar uno → su cookie 401, el otro 200; Max-Age y `expires_at` = día 90; uso el día 60 → día 150; sin uso, día 91 → 401. 2 PASS. CUMPLE.
- **S6-08**: de backend-base (no cambié `internal/jobs`).
- **S6-09**: `TestRespaldo_VacuumIntoIntegridadYRestauracion` (56 tablas iguales al origen) y `TestRespaldo_RestaurarConservaCatalogoYSesiones`; con el binario: `finanzas respaldo` → «integrity_check: ok, restauración de prueba: ok»; `sqlite3 datos/tmp/restaura.db "PRAGMA integrity_check;"` = `ok`; conteo por tabla origen contra restauración: 56 tablas, 0 diferencias (no vacías: goose_db_version 5, app_setting 30, meal_slot 4, notification 7, job 2, backup 2). Hora y retención PROVISIONAL (bloqueo S6-respaldos); falta copiarlas a `ops/runbook.md` (WF-4). CUMPLE salvo el runbook.
- **S6-10**: `go test ./internal/api/... -run Adjunto -count=1 -v`: 3 PASS (subir, descargar con sesión, sin sesión 401, nombre generado `adjuntos/2026/10/<32 hex>.png`, `../` y `..\` → 422, ruta alterada en la base → 404, exactamente 10 MB → 201, 10 MB + 1 byte → 413 `payload_too_large` «10 MB», límite cambiado a 1 → 1 MB 201 y 1 MB + 1 413). Factor 1 MB = 1,048,576 bytes PROVISIONAL (C1 no lo declara). CUMPLE con esa reserva.
- **S6-11**: `mkcert` **no está instalado** (no lo instalé: el certificado lo emite el integrador). `TestEscucha_HTTPSConCertificado` sirve HTTPS en 2 direcciones con un certificado `finanzas.local` y no responde en HTTP. Comandos listos en `ops/mkcert.md` §5. Android: pendiente de Axel. NO VERIFICADO.
- **S6-12**: binario en `127.0.0.1:18443 --dev-http`: `/manifest.webmanifest` válido (name «Finanzas Hogar», display standalone, start_url /inicio, íconos 192 y 512 autoalojados), `/sw.js` 200 `text/javascript`. `grep` de `https?://` en `frontend/dist` *.html/*.css: 1 línea, `https://tailwindcss.com` en el comentario de licencia del CSS (no es una carga). CUMPLE en el servidor (el resto es de frontend-base).
- **S6-13**: `go test ./internal/api/... -run IA -count=1 -v`: `TestIA_ApagadaNoLlamaALaRed` PASS (apagada: 0 llamadas del cliente que falla; encendida: 1). `grep -rn -i "api[_-]\?key\|sk-ant\|bot[0-9]\{6,\}\|password *= *\"" backend --include=*.go --include=*.yaml --include=*.sql --include=*.env` = 0. CUMPLE.
- **S6-14**: `serve -addr 0.0.0.0:18443` → rechazado; `100.100.1.2:18443` (Tailscale) → rechazado; sin certificado → rechazado; `ss -ltnp | grep 18443` = solo `127.0.0.1:18443`. `TestEscucha_SoloLANYLocal` PASS. La IP de la LAN va en `ops/runbook.md` (WF-4). CUMPLE en el código.
- **S6-15**: `go test ./internal/api/... -run Log -count=1`: `TestLog_SinContrasenasNiCookies` PASS (ni la contraseña mala, ni la nueva, ni la falsa, ni el token, ni su hash, ni `fh_session=` en el registro). CUMPLE.
- **S6-16**: §Supuestos y 3 bloqueos nuevos.
- **Globales**: `make verify-wf3-contrato` y `make verify-wf3-cobertura` OK; `make verify-wf3` y `make verify-all` fallan **solo** por `gofmt` de `internal/reportes/vistas.go` (s7-pdf, en curso).

## Supuestos (PROVISIONAL, con bloqueo)

1 MB = 1,048,576 bytes; tipos JPG/PNG/WebP/PDF; `kind` ausente = `other`; borrar un adjunto deja en NULL sus referencias (S6-adjuntos). Respaldo 03:15 y 30 diarios; restaurar conserva `backup` y `device_session`; respaldo y restauración síncronos (S6-respaldos). Semilla por `request_id = 'seed:<código>'`; CA por `--ca`; usuario exacto; cambiar contraseña no cierra otras sesiones; `listJobs` por `code`; IA sin llamada real decidida (S6-sistema).

## Bloqueos

`docs/estado/bloqueos/S6-adjuntos.md`, `S6-respaldos.md`, `S6-sistema.md`. Ajeno: `gofmt` de `internal/reportes/vistas.go` rompe `verify-wf3`/`verify-all`.

## Pendientes

- Integrador: línea `auth.Register,` en `backend/cmd/finanzas/areas.go`; instalar mkcert y emitir el certificado (`ops/mkcert.md` §1–2); correr S6-11.
- WF-4: `ops/runbook.md` con la IP de la LAN y la hora/retención de respaldos.
- Axel: prueba en los dos Android (S6-11).

## Propuestas para Axel

Ver las tablas de los 3 bloqueos; en corto: confirmar 1 MB = 1,048,576 bytes; dejar los respaldos a las 03:15 con 30 copias; que restaurar conserve sesiones y catálogo; subir el mínimo de `auth.password_min_length` en C1 a 12 (hoy permite 0); cerrar las otras sesiones al cambiar la contraseña; entrada `finanzas.local` en el DNS del router.

## Paquetes instalados

Ninguno (se usan `modernc.org/sqlite` y la biblioteca estándar; no se instaló mkcert).

## Consumo aproximado

~51 acciones · ~165 KB leídos · ~140 KB escritos.
