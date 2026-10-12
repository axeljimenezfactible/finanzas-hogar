# integrador-backend · WF-3a «Construcción»

Agente integrador-backend (Opus) · 12-oct-2026 · repositorio en HEAD «[Puerta 2 bis] K53» (sin git add/commit: lo hace el integrador final). Criterios: `docs/criterios/S6.md` 1.2. Estado de revisiones recibido: backend-core aprobado (2 rondas), backend-areas aprobado (3 rondas), frontend-core y frontend-areas **no aprobados** (3 rondas cada uno).

## Hecho

- **Go completo en verde** (paso 1): `gofmt -l .` vacío, `go vet ./...` y `go build ./...` sin errores, `go test ./... -count=1 -cover` con 32 paquetes OK y 5 sin pruebas; cobertura de los 10 paquetes de S6-02 abajo.
- **Contrato** (paso 2): `build.py --check` (23 archivos idénticos), `verify_c1.py --db` sobre base migrada con goose (TODO OK) y paridad de rutas: 216 en `routes.go` = 216 en `openapi.yaml`; el servidor registra 216 de 216 (`implementadas=216` en el log de arranque).
- **mkcert** (paso 3): instalado mkcert 1.4.4 con la única orden sudo permitida; CA local del usuario deepsec creada con `TRUST_STORES=nss mkcert -install` (el `-install` completo usa `sudo trust anchor` para el almacén del sistema, y ese sudo no está permitido); certificado `datos/certs/finanzas.local.pem` con SAN `finanzas.local`, `localhost`, 192.168.50.56 y 127.0.0.1 (vence el 12-ene-2029), llave en modo 600; `datos/certs/` en `.gitignore`. Nombre `finanzas.local` → 192.168.50.56 con `avahi-publish -a -R` en segundo plano (sin `/etc/hosts`, sin sudo).
- **Servidor de desarrollo CORRIENDO** (paso 4), para WF-3b:
  - `ops/dev-run.sh` (nuevo): compila `backend/bin/finanzas`, migra y siembra `datos/finanzas-dev.db` si no existe (`finanzas migrate` y `finanzas semilla`), publica `finanzas.local` y arranca `finanzas serve` en segundo plano en **192.168.50.56:8443 y 127.0.0.1:8443**.
  - **PID 1844644** (`datos/dev/run/finanzas.pid`); avahi-publish PID 1844253 (`datos/dev/run/avahi-publish.pid`); log `datos/dev/run/finanzas.log` (0 renglones ERROR).
  - **Detener:** `bash ops/dev-stop.sh` (SIGTERM; a los 10 s, SIGKILL; también detiene avahi-publish; `--conservar-nombre` deja el nombre publicado).
  - **Telegram apagado en desarrollo:** el proceso corre con `HOME=datos/dev/home`, sin `.openclaw/openclaw.json`. Los jobs `daily_summary`, `weekly_summary` y `alerts` ya corrieron (`last_status = ok`) con 0 destinatarios y 0 entregas (`notification_delivery` = 0). Ningún envío.
  - Semilla: 3 personas, 2 cuentas (`$argon2id$`), 160 productos, 186 presentaciones, 60 recetas, 16 pagos programados; 6 «cotejo OK» contra el oráculo (reparto $1,903.00 / $1,191.50; despensas 1–4; C3-09 $973.75 y $6,293.44). Contraseñas iniciales de **desarrollo** aleatorias en `~/.config/finanzas-hogar/dev-semilla.txt` (fuera del repositorio, modo 600, nunca impresas).
- **Makefile** (paso 5): objetivos finales `verify-wf3-secretos` (S6-13), `verify-wf3-respaldo` (S6-09 sobre una copia de la base de desarrollo), `verify-wf3-servidor` (S6-11, S6-12, S6-14 contra el servidor corriendo) y `verify-wf3-final` (verify-wf3 + los tres); `.PHONY`, encabezado y `make help` actualizados. `verify-all` sin cambios.
- `make verify-wf3-final` → OK (4 min 54 s; la interfaz son ~3.5 min: build + 149 pruebas de Playwright, 1 omitida). `make verify-all GOOSE=…` → OK (29 s).
- `ops/mkcert.md` v1.1 (lo hecho, `datos/certs/`, IP, `TRUST_STORES=nss`) y `ops/runbook.md` 0.1 (borrador para WF-4: dirección, nombre, cómo correr, cómo respaldar).

## Archivos

| Archivo | Cambio |
|---|---|
| `Makefile` | + `verify-wf3-final`, `verify-wf3-secretos`, `verify-wf3-respaldo`, `verify-wf3-servidor` (variables `WF3_HOST`, `WF3_PORT`, `WF3_DEV_DB`, `MKCERT`, `CURL`) |
| `ops/dev-run.sh` | nuevo (ejecutable) |
| `ops/dev-stop.sh` | nuevo (ejecutable) |
| `ops/mkcert.md` | v1.0 → v1.1 |
| `ops/runbook.md` | nuevo, borrador 0.1 |
| `.gitignore` | + `datos/certs/` (pedido en la tarea) y `datos/dev/` (datos del servidor de desarrollo; ver Supuestos) |
| `docs/estado/bloqueos/WF3-integrador-red.md` | nuevo bloqueo: IP fija, nombre permanente, CA en el sistema |
| `docs/estado/wf3/integrador-backend.md` | este reporte |

Fuera de git (ignorados): `datos/certs/*`, `datos/finanzas-dev.db*`, `datos/dev/`, `backend/bin/finanzas`; fuera del repositorio: CA en `~/.local/share/mkcert/`, contraseñas de desarrollo en `~/.config/finanzas-hogar/dev-semilla.txt`.

## Verificación

- **S6-01** · `go version` = go1.27.2; `go.mod`: `go 1.27`; `go build ./...` y `go vet ./...` sin diagnósticos; `gofmt -l .` sin archivos; `go mod tidy -diff` sin cambios (verify-wf3-go). **Cumple.**
- **S6-02** · comando literal (10 paquetes, `./internal/$p/...`), 0 fallas en `go test ./... -count=1`:
  - Con umbral ≥ 85 %: dominio **87.0 %**, motor_fin **89.2 %**, motor_cred **94.0 %**, recetas **93.2 %**, precios **88.8 %**. **Cumple.**
  - Sin umbral (K25): api 82.0 %, auth 83.6 %, jobs 92.3 % (93.0 % en la corrida de verify-wf3-cobertura), reportes 77.3 %, telegram 80.0 %.
- **S6-03** · (a) `build.py --check`: «los 23 archivos coinciden byte a byte con el generador», rc 0. (b) goose up (versión 4) + `verify_c1.py --db`: «RESULTADO: TODO OK», incluida «S6-03 rutas en routes.go: 216 · en openapi.yaml: 216 · sobrantes: [] · faltantes: []». (c) `TestRutasIgualesAlContrato`: «216 operaciones en 172 rutas, 0 diferencias»; `TestContrato_Rutas`: «contrato 216, servidor 216, sobrantes 0, faltantes 0, sin implementar 0». (d) `TestContrato_Operaciones`: «**216 de 216** operaciones cumplen (2xx declarada, 4xx válidas contra openapi.yaml)». La tabla operación → prueba es la de `docs/estado/wf3/pruebas-contrato.md` (sus cifras de 191/216 eran de una ronda anterior; hoy 216/216). **Cumple.**
- **S6-04** · `TestContrato_Idempotencia24h` (clave retenida 23 h 59 min, a las 24 h 1 min ya no), `TestContrato_CapturasAtrasadas` (3 gastos y 2 pesos con 48 h de atraso reenviados no se duplican), `TestContrato_PaginaMaxima` (200 máximo; 201/0/-1/abc = 400; 50 por omisión): PASS. La revisión de 20 mensajes la hicieron los revisores de backend (aprobado). **Cumple** según pruebas.
- **S6-05** · `go test ./internal/auth/... -count=1 -v`: 15 PASS, 0 FAIL; `SELECT substr(password_hash,1,10) FROM user_account` en la base de desarrollo = `$argon2id$` × 2.
- **S6-06** · `-run Intentos`: 3 PASS, 0 FAIL.
- **S6-07** · `-run Sesion`: 2 PASS, 0 FAIL.
- **S6-08** · `go test ./internal/jobs/... -count=1 -v`: 3 PASS, 0 FAIL; en el servidor vivo, tabla `job` con 7 jobs y `last_status` visible.
- **S6-09** · `make verify-wf3-respaldo` sobre copia de `datos/finanzas-dev.db`: `finanzas respaldo` (VACUUM INTO) → 1,810,432 bytes, `integrity_check: ok`, `restauración de prueba: ok`; restaurada en `datos/tmp/…/restaura.db`: `PRAGMA integrity_check` = **ok**; **56 tablas, 2,476 filas, 0 diferencias** de conteo con el origen. Hora 03:15 y retención 30: **PROVISIONAL** con la propuesta escrita en `ops/runbook.md` §4 (K31, puerta 4; bloqueo `S6-respaldos`). `ops/respaldo.sh` y la unidad systemd: WF-4.
- **S6-10** · `-run Adjunto`: 3 PASS (`TestAdjunto_SubirDescargarConYSinSesion`, `TestAdjunto_RutaConPuntosRechazada`, `TestAdjunto_Limite10MBYLimiteCambiado`), 0 FAIL. Tipos permitidos: PROVISIONAL (K31, bloqueo `S6-adjuntos`).
- **S6-11** · `mkcert -version` = **v1.4.4** (instalado hoy); `curl --cacert "$(mkcert -CAROOT)/rootCA.pem" … https://finanzas.local:8443/` = **200**; `openssl s_client … | openssl x509 -noout -ext subjectAltName` = `DNS:finanzas.local, DNS:localhost, IP Address:192.168.50.56, IP Address:127.0.0.1`; sin `--cacert` curl rechaza (código 60); HTTP plano en :8443 = 400 (no sirve contenido sin TLS). **Cumple en la máquina; Chrome Android: NO VERIFICADO, pendiente de Axel** (`ops/mkcert.md` §4).
- **S6-12** · `curl … /manifest.webmanifest | python3 -I -m json.tool`: válido; `name` «Finanzas Hogar», `display: standalone`, `start_url: /inicio`, íconos `/icons/icon-192.png` y `/icon-512.png` (+ maskable) servidos con 200 `image/png`; `sw.js` 200. Grep literal en `frontend/dist` (*.html, *.css): 1 renglón, el comentario de licencia `/*! tailwindcss v4.3.3 | MIT License | https://tailwindcss.com */` (no carga nada); 0 cargas de terceros. En JS solo hay literales que no se piden (w3.org de SVG, react.dev y bit.ly en mensajes de error; U5-08 del frontend: «hosts externos que se pedirían = 0»). Respaldo SPA y cabeceras CSP, `nosniff`, `X-Frame-Options: DENY` y `Referrer-Policy: no-referrer` presentes. **Cumple**, con un hallazgo menor: el manifiesto sale como `text/plain` (ver Solicitudes).
- **S6-13** · `-run IA`: 2 PASS, 0 FAIL; grep literal de llaves y tokens en `backend` (*.go, *.yaml, *.sql, *.env) = **0 coincidencias** (`make verify-wf3-secretos`). **Cumple.**
- **S6-14** · `ss -ltnp | grep 8443`: solo `127.0.0.1:8443` y `192.168.50.56:8443` (pid 1844644), sin `0.0.0.0` ni `::`; `tailscale0` = 100.115.28.35: `curl https://100.115.28.35:8443/` → «Could not connect» (no escucha). Dirección declarada en `ops/runbook.md` §1 (PROVISIONAL, bloqueo `WF3-integrador-red`). **Cumple.** Nota: 127.0.0.1 es local (no sale de la máquina) y lo pide la tarea para WF-3b.
- **S6-15** · `-run Log`: 2 PASS, 0 FAIL; el log del servidor vivo no tiene contraseñas, cookies ni tokens (solo arranque y direcciones).
- **S6-16** · Supuestos y bloqueos abajo; bloqueo nuevo `WF3-integrador-red.md`; los `S6-*` existentes son de otros agentes.
- **Global** · `make verify-wf3-final` → «WF-3a: cierre OK» (rc 0, 4:54). `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)» (rc 0, 29 s). verify-all no incluye el frontend, así que no tarda por él.

## Supuestos

- HECHO/PROPUESTA: `datos/dev/` agregado a `.gitignore` para adjuntos, reportes, respaldos, HOME aislado, PID y logs del servidor de desarrollo (la tarea solo nombró `datos/certs/`; sin esto, `datos/app` y `datos/respaldos`, que son las rutas por omisión de `serve`, quedarían sin ignorar). Se revierte quitando una línea.
- PROPUESTA: `HOME` aislado para que el servidor de desarrollo no pueda leer el token de MultiX. Las pruebas de Telegram de WF-3b o de la orquestación (S7-07) no deben esperar envíos reales de este servidor.
- PROPUESTA: contraseñas iniciales de desarrollo aleatorias en `~/.config/finanzas-hogar/dev-semilla.txt` (no son las de la casa). WF-3b las lee de ahí para entrar como `axel` o `ximena`.
- La IP 192.168.50.56 la dio la tarea y coincide con `wlp4s0`; queda PROVISIONAL (bloqueo `WF3-integrador-red`).
- `verify-wf3-frontend` vuelve a compilar `frontend/dist`; el servidor lee la carpeta del disco en cada petición, así que sirve siempre la última compilación sin reiniciarse.

## Bloqueos

- Nuevo: `docs/estado/bloqueos/WF3-integrador-red.md` (IP fija, nombre `finanzas.local` permanente, CA en el sistema). No detiene WF-3b.
- Citados (de otros agentes): `S6-respaldos` (hora, retención, qué conserva una restauración), `S6-sistema` (Android y `.local`), `S6-adjuntos` (tipos), `WF3-frontend-servidor` (cabeceras de caché y tipo del manifiesto: ver Solicitudes).

## Pendientes

- Chrome Android con la CA instalada (S6-11): manual de Axel.
- WF-4: unidad systemd, rutas de producción, ACL del token, IP y nombre permanentes, segunda copia de respaldos (`ops/runbook.md` §6).
- Frontend-core y frontend-areas siguen sin aprobar (no es de este agente); `verify-wf3-frontend` sí pasa hoy.

## Solicitudes fuera de alcance

1. `backend/internal/api/http.go`, `staticHandler`: servir `*.webmanifest` con `Content-Type: application/manifest+json` (hoy `text/plain; charset=utf-8` con `X-Content-Type-Options: nosniff`; por ejemplo `mime.AddExtensionType(".webmanifest", "application/manifest+json")` en un `init`). Agregar `Cache-Control: no-cache` a `index.html`, `sw.js`, `manifest.webmanifest` y a la respuesta del respaldo SPA, y `Cache-Control: public, max-age=31536000, immutable` a `/assets/*` (punto 3 del bloqueo `WF3-frontend-servidor`; hoy no hay `Cache-Control` en los estáticos). Dueño: backend-base o auth-sistema.
2. `docs/estado/bloqueos/WF3-frontend-servidor.md`: anotar que los puntos 1 (mismo origen) y 2 (respaldo SPA solo fuera de `/api/` y `/health`) están comprobados en el servidor vivo, que el 4 (CSP) ya se envía, y que el 3 sigue pendiente (solicitud 1).
3. `docs/estado/bloqueos/S6-respaldos.md` punto 1: «Debe copiarse a `ops/runbook.md` (no existe…)» → ya está en `ops/runbook.md` §4 como PROVISIONAL.
4. `docs/estado/bloqueos/INDICE.md`: agregar la fila de `WF3-integrador-red.md`.

## Propuestas para Axel

- Reservar 192.168.50.56 en el módem para esta máquina y, si el módem lo permite, una entrada DNS `finanzas.local → 192.168.50.56` (resuelve también los dos Android).
- Aprobar en la puerta 4 la hora (03:15) y la retención (30) de los respaldos, y considerar una segunda copia en otro disco local.
- Cuando frontend-core y frontend-areas se aprueben, incluir `verify-wf3` en `verify-all` (hoy agrega ~4.5 min, casi todo Playwright).

## Paquetes instalados

- `mkcert 1.4.4-3` (pacman, repositorio extra) con `sudo pacman -S --noconfirm mkcert`. Ningún módulo de Go ni paquete de npm.

## Consumo aproximado

- ~46 acciones de herramienta; ~170 KB leídos; ~50 KB escritos.
