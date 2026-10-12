# WF-3a · backend-base · Reporte

Agente backend-base (Opus) · 12-oct-2026 · repositorio en HEAD «[Puerta 2 bis] K53» (sin commit: lo hace el integrador). Criterios: `docs/criterios/S6.md` 1.2. Convenciones para los demás agentes: `docs/estado/wf3/CONVENCIONES-backend.md`.

## Hecho
- **Tipos desde el contrato**: `c1gen/gen8_types.py` (nuevo) genera `backend/internal/api/types_gen.go`: 343 tipos de `components/schemas` (structs con etiquetas json alineadas como gofmt; `Nullable[T]` para los tipos con null; alias para los escalares) y `Operations` (216: sesión/pública/agente externo, Idempotency-Key, tipo y esquema del cuerpo, códigos 2xx). Integrado en `build.py` (SOURCES, STAGES y OUTPUTS); los otros generadores y los extractos no cambian; `build.py --check` compara 23 archivos y da 0.
- **`internal/dominio`**: apertura de SQLite con modernc (WAL, llaves foráneas, `_txlock=immediate`), apertura de solo lectura, migraciones embebidas con goose como biblioteca (`goose.NewProvider`), `WithTx`, errores de negocio con clase → HTTP, dinero en centavos ROUND_HALF_UP, fechas en America/Mexico_City, filtro `period` (mes, quincena, semana ISO, año, rango), quincenas con día de pago 15 y 30/último día, repositorios de `user_account` (sin función para crear cuentas), argon2id PHC, `device_session` (90 días renovables), `login_attempt` (ventana deslizante por usuario y por IP), `app_setting` (lectura tipada, `@param` leído de las migraciones con valor inicial, cambiar y restablecer con auditoría), `idempotency_key`, `job`, `audit_log` (`RowJSON` + `Audit`) y respaldos (`VACUUM INTO`, `integrity_check`, restauración de prueba con conteos por tabla, retención).
- **`internal/api`**: `Registry` (cada área exporta `Register(reg)`; operación fuera de C1 o duplicada = pánico), middleware de sesión por cookie (renueva y reemite `fh_session` HttpOnly; Secure; SameSite=Strict), agente externo por bearer (sin verificador = 401, PROVISIONAL), idempotencia automática en las 44 operaciones, error único de C1 (`Problem`, `dominio.Error`, 500 sin detalle), 404/405 de C1, cabeceras PROVISIONALES, registro de accesos sin cookies ni cuerpos, `DecodeJSON` (400/413), `Validator` (422), paginación con cursor opaco (1–200, 50 por defecto), parámetros de ruta y filtros, PWA desde `frontend/dist` con regreso a `index.html`, `RegisterBase` (salud y contrato mientras auth-sistema no los registre). Lo no implementado sigue en 501.
- **`internal/jobs`**: planificador `time.Ticker` + tabla `job`, cron de 5 campos en America/Mexico_City, toma atómica por disparo, reinicio con job vencido = una corrida, `RunNow`.
- **`internal/ia`**: interfaz `Provider` (OCR de tickets y sugerencias de recetas), adaptador `Off`, `Registry` que lee `ai.local_model_enabled` y `ai.claude_api_enabled` (apagados) y no llama a nada si están apagados.
- **`internal/testutil`**: base temporal en `datos/tmp/` migrada con fixtures (Axel, Ximena; Karen sin cuenta; credencial falsa), reloj simulable, servidor TLS con sesión, `Do` que valida cada respuesta contra `openapi.yaml` con `kin-openapi/openapi3filter` (solo pruebas).
- **`cmd/finanzas`**: subcomandos `serve` (rechaza 0.0.0.0, :: y «:8443»; TLS con `--cert/--key` de mkcert; jobs base `backup_daily` e `idempotency_purge`), `migrate`, `semilla` (esqueleto: falla con mensaje claro mientras no exista `internal/semilla`) y `respaldo`; `areas.go` junta los `Register`.
- **`backend/db/embed.go`**: `go:embed migrations/*.sql` (go:embed no sube de carpeta; las migraciones no se tocaron).
- **`go.mod`**: goose, x/crypto (argon2), kin-openapi; `go mod tidy` (cierra I-2: modernc deja de ser indirecta).
- **Makefile**: `verify-wf3-go`, `verify-wf3-cobertura`, `verify-wf3-contrato`, `verify-wf3-frontend` y `verify-wf3`; `verify-all` no los incluye.
- **Reparto**: `getReportSettings` y `updateReportSettings` no tenían dueño; quedaron en s7-pdf (CONVENCIONES §14, PROPUESTA).

## Archivos
Nuevos: `backend/internal/api/c1gen/gen8_types.py`, `backend/internal/api/types_gen.go` (generado), `backend/internal/api/{server.go,http.go,nullable.go,api_test.go}`, `backend/internal/dominio/{dominio.go,errores.go,tiempo.go,dinero.go,ajustes.go,auditoria.go,cuentas.go,idempotencia.go,trabajos.go,respaldo.go,dominio_test.go}`, `backend/internal/jobs/{cron.go,jobs.go,jobs_test.go}`, `backend/internal/ia/ia.go`, `backend/internal/testutil/{testutil.go,servidor.go}`, `backend/cmd/finanzas/areas.go`, `backend/db/embed.go`, `docs/estado/wf3/CONVENCIONES-backend.md`, `docs/estado/wf3/backend-base.md`, `docs/estado/bloqueos/S6-valores-operativos-base.md`. Cambiados: `backend/internal/api/c1gen/build.py` (3 listas y el docstring), `backend/cmd/finanzas/main.go` (reescrito), `backend/go.mod`, `backend/go.sum`, `Makefile`.

## Verificación
- **S6-01**: `go version` = go1.27.2; `go.mod`: `go 1.27`; `go build ./...`, `go vet ./...` sin diagnósticos; `gofmt -l .` vacío. CUMPLE.
- **S6-02**: `go test ./... -count=1`: 14 paquetes ok, 0 fallas. Con umbral: dominio **87.3 %**, motor_fin 89.2 %, motor_cred 94.0 %, recetas 93.2 %, precios 88.8 %. Sin umbral: api 81.4 %, jobs 92.3 %; auth, reportes y telegram todavía no existen (los crean auth-sistema y s7-*). CUMPLE para lo de este agente.
- **S6-03**: (a) `build.py --check` → «los 23 archivos coinciden byte a byte» (incluye `types_gen.go`); (b) `verify-all` corre `verify_c1.py --db` → «TODO OK»; (c) `TestRutasIgualesAlContrato`: «216 operaciones en 172 rutas, 0 diferencias» entre `paths` de `openapi.yaml`, `api.Routes` y `api.Operations`; (d) `TestContratoSinSesionYEsqueleto` llama las 216 operaciones sin sesión (401 o pública) y con sesión (501, o 400 por falta de Idempotency-Key) y valida el 100 % de esas respuestas contra el esquema. La prueba de comportamiento por operación y su tabla son de cada área (CONVENCIONES §12, `TestOp_<operationId>`). PARCIAL (base lista).
- **S6-04**: `TestIdempotencia`: misma clave y cuerpo = 1 ejecución, 201 y luego 200 con `Idempotent-Replayed: true`; otro cuerpo = 409 `idempotency_key_reused`; la clave es por cuenta; un 422 libera la clave; pasadas 25 h la tabla ya no la tiene; claves vacía, corta o con espacios = 400; `973.75`, `"973.75"`, campo desconocido, JSON roto y datos de sobra = 400. `TestPaginacion`: `limit=201` = 400, 50 por defecto, cursor ligado a la consulta. El ciclo K9 por renglón (gasto y peso) es de api-gastos-inicio y api-comidas-salud. PARCIAL (base lista).
- **S6-05**: `TestSesionCookieYRenovacion`: `Set-Cookie` con `HttpOnly`, `Secure`, `SameSite=Strict`, `Path=/`, `Max-Age=7776000`; fixtures con hash `$argon2id$`; `TestMigracionesYApertura`: una cuenta para Karen la rechaza la base. Login y contraseña nueva de 11/12 caracteres son de auth-sistema. PARCIAL.
- **S6-06**: `TestCuentasSesionesEIntentos` (repositorio): fallos en 0, 1, 2, 3 y 14 → bloqueado en 14.5 con `retry_at` = minuto 15, libre en 15.5; otra cuenta desde la misma IP también queda bloqueada. La prueba `-run Intentos` en `internal/auth` es de auth-sistema. PARCIAL.
- **S6-07**: uso el día 60 → vence el día 150; sin uso, el día 91 = 401; sesión cerrada = 401 (api y dominio). La prueba `-run Sesion` en `internal/auth` es de auth-sistema. PARCIAL.
- **S6-08**: `go test ./internal/jobs/... -v`: día 30 (next_run_at 09:00 de México = 15:00Z), 4 disparos simultáneos + 1 repetido = 1 ejecución, reinicio con el job vencido 3 días y en `running` = 1 ejecución y salto al día siguiente, febrero de 2027 (21:30 del 28 = 03:30Z del 1 de marzo; siguiente 2028-02-29T03:30Z), errores y pánicos anotados, job apagado no corre. CUMPLE.
- **S6-09**: `TestRespaldoYRestauracion`: 3 respaldos con `VACUUM INTO`, `integrity_check` = ok, restauración de prueba ok con filas por tabla iguales a las del origen, renglón `backup` en ok con sha256, retención deja 2 de 3. Manual: `finanzas respaldo` → «integrity_check: ok», «restauración de prueba: ok». Hora 03:15 y 30 copias PROVISIONAL (bloqueo `S6-valores-operativos-base`); `ops/respaldo.sh` y `ops/runbook.md` son de WF-4. PARCIAL (falta aprobación).
- **S6-10**: de auth-sistema (adjuntos); la base da 413 con `Problem` y el valor `attachment.max_size_mb` = 10. NO APLICA a este agente.
- **S6-11**: `serve` acepta el certificado de mkcert (`--cert`, `--key`); no se instaló mkcert ni se generó certificado (WF-4/auth-sistema). Android: pendiente de Axel. NO VERIFICADO.
- **S6-12**: `TestFrontendYCabeceras` sirve `index.html`, `manifest.webmanifest` y regresa `index.html` en rutas de la PWA, con 404 JSON en `/api/…`. El manifiesto real es de frontend-base. PARCIAL.
- **S6-13**: `TestIAApagadaSinLlamadas`: interruptores apagados al sembrar; con un adaptador que saldría a la red (transporte que falla la prueba), `Active` devuelve `ErrDisabled` sin llamarlo. `grep -rn -i "api[_-]\?key\|sk-ant\|bot[0-9]\{6,\}\|password *= *\"" backend --include=*.go --include=*.yaml --include=*.sql --include=*.env` = **0**. CUMPLE para la base.
- **S6-14**: `finanzas serve --addr 0.0.0.0:8443` → «usa la IP de la LAN (192.168.50.x) o 127.0.0.1, nunca un comodín». Prueba de humo en 127.0.0.1:18443: `/health` 200, `/api/v1/users` 401. La IP LAN la fija WF-4. PARCIAL.
- **S6-15**: `TestLogSinSecretosYErrores`: el registro no contiene el token de sesión; error interno con la contraseña falsa en su texto → la respuesta no la muestra; códigos en inglés y mensajes en español. CUMPLE para la base.
- **S6-16**: §Supuestos y bloqueo `docs/estado/bloqueos/S6-valores-operativos-base.md`. CUMPLE.
- `make verify-wf3` → «== WF-3: verificación completa OK» (la parte de interfaz corrió `npm run build` de `frontend/package.json`, que ya existe). `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)»; su cobertura K25 ya mide `internal/dominio` (87.3 %).
- El ejemplo de CONVENCIONES §13 se compiló y se probó en un paquete temporal (pasa la validación del contrato) y se borró.

## Supuestos
1. `backend/db/embed.go` fuera de la lista de alcance: es la única forma de embeber las migraciones sin copiarlas (go:embed no sube de carpeta).
2. Un campo JSON desconocido es 400 (como un tipo incorrecto), no 422.
3. Repetición idempotente: solo se guardan respuestas 2xx; un 201 guardado se repite como 200 cuando la operación declara 200 (regla de K9 generalizada).
4. Una clave con más de `idempotency.retention_hours` cuenta como purgada aunque el job no haya corrido.
5. Un intento rechazado por bloqueo no se anota como fallo; `retry_at` = momento en que la ventana baja de 5 fallos (con 5 fallos es «el más antiguo + 15 min» de C1).
6. Cursor opaco por desplazamiento con huella de la consulta (no por llave).
7. La sesión se renueva y la cookie se reemite en cada petición con sesión.
8. El planificador crea el renglón `job` que falte con un horario inicial del código; después manda la tabla.
9. Nombres de paquete por agente (CONVENCIONES §2).

## Bloqueos
- Nuevo: `docs/estado/bloqueos/S6-valores-operativos-base.md` (hora y retención de respaldos, argon2id, cabeceras, límites de cuerpo, job de purga, rutas por omisión, campo desconocido = 400). Relacionado: `S5-agente-externo` (credencial del agente externo).

## Pendientes
- Pruebas por operación y tabla de S6-03: cada área. Líneas en `backend/cmd/finanzas/areas.go`: integrador.
- `internal/semilla`: no tiene dueño en el reparto (PREGUNTA).
- I-1 (`cmd/motores` desde `backend/`) sigue abierto: `backend/cmd/motores/main.go` no es de este agente.
- `ops/runbook.md` y `ops/respaldo.sh` (WF-4): copiar la propuesta del bloqueo.
- Constantes de los `enum` en `types_gen.go`: no se generan (se puede agregar a gen8 si las áreas lo piden).

## Propuestas para Axel
- Aprobar en la puerta 4 los 7 valores de `S6-valores-operativos-base` (respaldo diario 03:15 con 30 copias; argon2id RFC 9106; cabeceras; límites de cuerpo).
- Dejar `getReportSettings` y `updateReportSettings` con s7-pdf.

## Paquetes instalados
Módulos de Go: `github.com/pressly/goose/v3 v3.28.0`, `golang.org/x/crypto v0.58.0`, `github.com/getkin/kin-openapi v0.149.0` (solo pruebas), con sus indirectos en `go.mod`/`go.sum` (`mfridman/interpolate`, `sethvargo/go-retry`, `go.uber.org/multierr`, dependencias de kin-openapi) y `golang.org/x/sys` 0.48 → 0.49. Sin npm, sin sudo.

## Consumo aproximado
~55 acciones; ~235 KB leídos; ~205 KB escritos a mano (código Go ~150 KB incluidas pruebas, generador 8 KB, CONVENCIONES 17 KB, reporte y bloqueo 13 KB) más 183 KB generados por `build.py` (`types_gen.go`). **Se pasó el tope de 150 KB escritos** del protocolo §4 por las pruebas necesarias para la cobertura ≥ 85 % y el contrato; no se cortó trabajo.
