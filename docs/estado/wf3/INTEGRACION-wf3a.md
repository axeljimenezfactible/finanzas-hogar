# INTEGRACIÓN de WF-3a «Construcción»

Agente `integrador-wf3a` (Opus) · 12-oct-2026 · repositorio en HEAD «[Puerta 2 bis] K53» antes del commit de WF-3a · contrato C1 v1.4.0 (216 operaciones en 172 rutas), esquema 1.4, criterios S6 1.2, S7, U5 y C3 1.1.

Estado de las revisiones al cerrar: **backend-core APROBADO** (2 rondas), **backend-areas APROBADO** (3 rondas), **frontend-core NO APROBADO** (3 rondas), **frontend-areas NO APROBADO** (3 rondas). El servidor de desarrollo queda corriendo para WF-3b en `https://finanzas.local:8443` (192.168.50.56:8443 y 127.0.0.1:8443), PID 1844644; se detiene con `bash ops/dev-stop.sh`.

## 1. En corto

- **HECHO.** El servidor Go atiende las 216 operaciones de C1 v1.4, y las 216 pasan su prueba de contrato contra `openapi.yaml`. Incluye autenticación argon2id, sesiones de 90 días, límite de intentos, jobs, respaldos con restauración probada, adjuntos, IA apagada, semilla real, avisos por Telegram (sin envío real) y PDF con LaTeX. Sirve la PWA por HTTPS con mkcert, solo en la red de casa.
- **HECHO.** El frontend tiene el cascarón, las 41 pantallas de U1 (27 `P-xx` y 14 secciones de Ajustes), la PWA instalable y la cola sin conexión. Las 41 pantallas más la fila Global tienen pruebas de Playwright con la API simulada (149 pasan, 1 omitida, 0 fallan). Una prueba de humo contra el servidor real cuadra con el oráculo: $973.75, $6,905.50 y $10,000.00.
- **HECHO.** C3 (`docs/contratos/C3-integracion.md`) está en v0.9, como **borrador sin firma**: implementado = sí en 35 de 42 filas, probado = sí en 42 de 42 y capturado = no en las 42 (eso lo hace WF-3b).
- **PREGUNTA.** El frontend no quedó aprobado tras 3 rondas. Siguen abiertos 1 hallazgo alto y 5 medios, casi todos en archivos de `frontend-base` (§6). ¿Se encarga una corrección de `frontend-base` con revisión puntual antes de WF-3b, o se acepta con estos pendientes?
- **PREGUNTA.** Hay 40 bloqueos nuevos, todos abiertos y con opción PROVISIONAL en el código (§8). Los de valores operativos (K31) llevan su propuesta en §7.
- **PENDIENTE.** S7-07 (envío real por Telegram) lo hace el orquestador. La prueba de S6-11 en Chrome Android la hace Axel a mano.

## 2. Verificación global (corrida por este integrador)

Todo se corrió el 12-oct-2026 con el servidor de desarrollo en marcha (PID 1844644).

- **`make verify-wf3` → rc 0** («WF-3: verificación completa OK», 4 min 13 s). Detalle:
  - Go (S6-01): go1.27.2, `go build`, `go vet`, `gofmt -l` y `go mod tidy -diff` limpios; `go test ./...` con 32 paquetes OK, 5 sin pruebas y 0 fallas.
  - Cobertura K25 (S6-02): ver §5.
  - Contrato (S6-03): `build.py --check` dice «los 23 archivos coinciden byte a byte con el generador»; `TestRutasIgualesAlContrato` da «216 operaciones en 172 rutas, 0 diferencias».
  - Interfaz: `npm run verify` con tsc y build limpios. La salida dice «C3-03 frontend: 202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)» y «U5-08 dist: hosts externos que se pedirían = 0». Playwright: **149 passed, 1 skipped, 0 failed** (3.5 min).
- **Primera corrida (con `make -k verify-wf3 verify-wf3-secretos verify-wf3-respaldo verify-wf3-servidor`): rc 2 por una prueba intermitente.** Falló `pruebas/cola.spec.ts:85` [movil] («falla del servidor (500) al enviar la cola…»): pasó el tiempo de 45 s en `tap` sobre «Revisar», con el mensaje «element was detached from the DOM, retrying». En esa corrida, 148 pasaron. La causa probable es la carga: corrí `go test` del contrato al mismo tiempo. Después, el spec aislado con `--repeat-each=8` dio **32 de 32 pasan**, y la segunda corrida completa de arriba quedó en verde. **No lo corregí**, porque no es de mi alcance. Queda como hallazgo bajo para `frontend-base`: la franja se vuelve a pintar con el sondeo y la prueba debería esperar a que quede estable.
- **`make verify-wf3-secretos` (S6-13) → OK**: «S6-13 grep: 0 coincidencias».
- **`make verify-wf3-respaldo` (S6-09) → OK.** Hizo el respaldo de una copia de la base de desarrollo con `VACUUM INTO`: 1,810,432 bytes, `integrity_check: ok` y «restauración de prueba: ok». La restaurada en `datos/tmp` da `integrity_check = ok`, **56 tablas, 2,482 filas y 0 diferencias de conteo**.
- **`make verify-wf3-servidor` (S6-11, S6-12 y S6-14) → OK:**
  - S6-11: mkcert v1.4.4; `finanzas.local` → 192.168.50.56; curl con la CA da 200; SAN `DNS:finanzas.local, DNS:localhost, IP Address:192.168.50.56, IP Address:127.0.0.1`. «Chrome Android con la CA instalada: NO VERIFICADO (manual de Axel, ops/mkcert.md §4)».
  - S6-12: manifiesto con nombre, `start_url`, `display standalone` e íconos 192, 512 y maskable con 200 `image/png`; `sw.js` 200; 0 cargas de terceros en `dist`.
  - S6-14: solo `127.0.0.1:8443` y `192.168.50.56:8443`; Tailscale 100.115.28.35:8443 no responde.
- **`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → rc 0**: «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». Los oráculos S2 a S5 dan 0 diferencias y K25 se cumple en los 11 paquetes de motores.
- **Prueba de contrato aislada**: `go test ./internal/api -run TestContrato_Operaciones -count=1 -v` → rc 0, «S6-03: 216 de 216 operaciones cumplen», 216 subpruebas PASS y 0 FAIL.
- **Secretos en lo que entra al commit** (script propio que solo imprime conteos): revisé 429 archivos y no encontré llaves privadas, tokens de bot, llaves de API, hashes argon2 completos ni las contraseñas de desarrollo de `~/.config/finanzas-hogar/dev-semilla.txt`. `git check-ignore` confirma que `datos/finanzas-dev.db`, `datos/certs/`, `datos/dev/`, `frontend/node_modules` y `frontend/dist` quedan fuera.

## 3. Hecho por agente y paquete

| Agente | Paquete o ruta | Qué dejó |
|---|---|---|
| `backend-base` | `backend/internal/api/` (`http.go`, `server.go`, `nullable.go`, `types_gen.go`), `c1gen/gen8_types.py`, `backend/db/embed.go`, `cmd/finanzas/areas.go`, `Makefile` (`verify-wf3-*`), `CONVENCIONES-backend.md` | 343 tipos generados desde C1, registro de handlers por área, errores, paginación, idempotencia de 24 h, `goose` embebido y convenciones |
| `auth-sistema` | `internal/auth/`, `internal/api/sistema/`, `internal/adjuntos/`, `internal/respaldo/`, `internal/jobs/`, `internal/ia/` | login con argon2id, 5 intentos en 15 min, sesión de 90 días por dispositivo, jobs, respaldo con `VACUUM INTO`, adjuntos de 10 MB e IA opcional apagada |
| `semilla` | `internal/semilla/`, `cmd/finanzas/semilla.go` | carga idempotente de `datos/semilla/` en una sola transacción: 3 personas, 2 cuentas, 160 productos, 186 presentaciones, 60 recetas y 16 pagos programados; 6 cotejos con el oráculo |
| `api-gastos-inicio` | `internal/api/gastos-inicio/` | Gastos e Inicio: 12 tarjetas y lecturas A y B (incluida la mejor canasta) |
| `api-presupuesto` | `internal/api/presupuesto/` | 28 operaciones de presupuesto, flujo, ingresos y ahorro (motor S2) |
| `api-comidas-salud` | `internal/api/comidas/` | 26 operaciones de comidas y salud (motor S4) |
| `api-compras-a` | `internal/api/comprasa/` | 18 operaciones: lista, canasta y comparativa (S4 y S5) |
| `api-compras-b` | `internal/api/comprasb/` | 23 operaciones: precios, productos, inventario y las 4 del agente externo (K53) |
| `api-creditos` | `internal/api/creditos/` | 29 operaciones de créditos y del comparador (motor S3) |
| `api-ajustes` | `internal/api/ajustes/` | 19 operaciones de Ajustes: categorías, personas, partidas e importación de CSV |
| `s7-avisos` | `internal/telegram/`, `internal/reportes/avisos/`, `internal/api/avisos/`, `ops/telegram-disparo.md` | resúmenes diario y semanal, 4 alertas, destinatarios y chat, texto plano y token leído sin copiarlo |
| `s7-pdf` | `internal/reportes/` (plantillas LaTeX embebidas), `internal/api/reportes/` | reporte mensual, estado de cuenta de Karen, vista web mensual y cruce PDF = API |
| `pruebas-contrato` | `internal/api/contrato_test.go`, `contrato_datos_test.go`, `contrato_tiempo_test.go` | una prueba por operación validada con kin-openapi, idempotencia de 24 h, página máxima y capturas atrasadas |
| `frontend-base` | `frontend/` (Vite, React, Tailwind, Radix, TanStack, dnd-kit, ECharts, PWA, msw, Playwright), `src/api`, `src/cola`, `src/shell`, `src/componentes`, `src/vistas-bd`, `CONVENCIONES-frontend.md` | cascarón, cliente de API tipado, cola sin conexión en IndexedDB, tema claro, oscuro y del sistema, vistas de BD y fuentes e íconos autoalojados |
| `front-inicio` | `src/pantallas/inicio/` | P-01 Inicio y su Detalle |
| `front-gastos-presupuesto` | `src/pantallas/gastos-presupuesto/` | P-02 a P-09 y F-1 |
| `front-comidas-compras` | `src/pantallas/comidas-compras/` | P-10 a P-16 y F-2 |
| `front-creditos-salud-reportes` | `src/pantallas/creditos-salud-reportes/` | P-17 a P-21, P-23 a P-25 y F-4 a F-6 |
| `front-ajustes-a` | `src/pantallas/ajustes-a/` | P-26 y secciones 01 a 07 |
| `front-ajustes-b` | `src/pantallas/ajustes-b/` | secciones 08 a 14 |
| `front-variantes` | `frontend/.storybook/`, `scripts/contraste-variantes.*` | Storybook 10.6.1 sin red y contraste de variantes |
| `integrador-backend` | `ops/dev-run.sh`, `ops/dev-stop.sh`, `ops/mkcert.md` 1.1, `ops/runbook.md` 0.1, `Makefile` (`verify-wf3-final`, `-secretos`, `-respaldo`, `-servidor`) | mkcert, servidor de desarrollo corriendo y cierre de S6 |
| `integrador-frontend-c3` | `frontend/scripts/c3-integracion.py`, `playwright.integracion.config.ts`, `pruebas-integracion/humo-inicio.spec.ts`, `docs/contratos/C3-integracion.md` v0.9 | matriz C3 generada y comparada por script, y humo real de Inicio |
| correcciones y revisores | `correccion-*-rN.md`, `revisor-*-rN.md`, `REVISION-*-rN.md` | 7 correcciones y 10 revisiones (backend-core 2 rondas, backend-areas 3, frontend-core 3 y frontend-areas 3) |

Los reportes completos de cada agente están en `docs/estado/wf3/<agente>.md`.

## 4. Criterios S6 y S7

### S6 (servidor, 16 criterios): 15 cumplen y 1 cumple en la máquina con una parte manual pendiente

| Id | Estado | Evidencia (resumida) |
|---|---|---|
| S6-01 | Cumple | go 1.27.2; `go build`, `go vet`, `gofmt -l` y `go mod tidy -diff` limpios (verify-wf3-go) |
| S6-02 | Cumple | 32 paquetes OK y 0 fallas; cobertura K25 en §5 |
| S6-03 | Cumple | `build.py --check`: 23 archivos idénticos; 216 rutas en `routes.go` = 216 en `openapi.yaml`; 216 de 216 operaciones cumplen (§5) |
| S6-04 | Cumple | `TestContrato_Idempotencia24h`, `TestContrato_CapturasAtrasadas` y `TestContrato_PaginaMaxima` pasan |
| S6-05 | Cumple | 15 pruebas de `internal/auth` pasan; `password_hash` = `$argon2id$` en las 2 cuentas |
| S6-06 | Cumple | `-run Intentos`: 3 pasan |
| S6-07 | Cumple | `-run Sesion`: 2 pasan |
| S6-08 | Cumple | `internal/jobs`: 3 pasan; 7 jobs con `last_status` en el servidor vivo |
| S6-09 | Cumple (hora y retención PROVISIONAL) | `make verify-wf3-respaldo`, ver §2 |
| S6-10 | Cumple (tipos PROVISIONAL) | `-run Adjunto`: 3 pasan; 10 MB entran y 10 MB + 1 byte dan 413 |
| S6-11 | **Cumple en la máquina; Chrome Android NO VERIFICADO (lo prueba Axel)** | mkcert 1.4.4; curl con la CA da 200; SAN con `finanzas.local`; `ops/mkcert.md` §4 |
| S6-12 | Cumple | manifiesto válido, íconos 192 y 512, `sw.js` 200 y 0 cargas de terceros en `dist`. Hallazgo menor: el manifiesto sale como `text/plain` (§11) |
| S6-13 | Cumple | `-run IA`: 2 pasan; grep de llaves y tokens en `backend/` = 0 |
| S6-14 | Cumple | escucha solo en 192.168.50.56:8443 y 127.0.0.1:8443; nada en `0.0.0.0`, `::` ni Tailscale |
| S6-15 | Cumple | `-run Log`: 2 pasan; el log no trae contraseñas, cookies ni tokens |
| S6-16 | Cumple | supuestos y bloqueos `S6-*` documentados (§8) |

### S7 (avisos y reportes, 14 criterios): 13 cumplen y 1 queda pendiente del orquestador

| Id | Estado | Evidencia (de `s7-avisos`, `s7-pdf` y de la revisión `REVISION-backend-areas-r3`) |
|---|---|---|
| S7-01 | Cumple | `latexmk -lualatex`: mensual con rc 0 y 2 páginas, estado de cuenta con rc 0 y 1 página; 0 errores y 0 referencias indefinidas |
| S7-02 | Cumple | `-run Cruce`: 36 de 36 cifras y fechas en el agente; 256 en 6 PDF en la revisión; 0 diferencias |
| S7-03 | Cumple | `-run Karen`: 20 de 20; sin ingreso de Axel; Karen no recibe Telegram (422) |
| S7-04 | Cumple | diario a las 21:30 a 2 destinatarios; semanal el domingo; respeta los cambios de Ajustes |
| S7-05 | Cumple | sin `parse_mode`; un mensaje de 12,296 bytes sale en 3 partes de ≤ 4096 |
| S7-06 | Cumple (la ACL real es de WF-4) | grep del patrón de token en todo el repositorio = 0; `TestToken` con un token falso |
| **S7-07** | **PENDIENTE: lo hace el orquestador** | contra el servidor falso, `sendTestNotification` → 202 y `http_status` 200; el servidor de desarrollo corre sin token (HOME aislado) |
| S7-08 | Cumple | matriz de 11 fechas con 0 diferencias; umbrales de 80 % y 100 % |
| S7-09 | Cumple | Axel con chat, Ximena sin chat se omite, Karen rechazada |
| S7-10 | Cumple | PDF v1 el día 1; un día 31 cae al último día del mes; plantilla rota → 0 envíos |
| S7-11 | Cumple | 34 de 34 cifras de `getMonthlyReport` = PDF |
| S7-12 | Cumple | 39 cifras con el patrón de dinero, 0 fechas ISO y 0 palabras en inglés |
| S7-13 | Cumple | 7 colores = tokens de C2; contraste mínimo 5.04:1 |
| S7-14 | Cumple | `go vet` limpio; solo `https://api.telegram.org`; LaTeX con `-no-shell-escape` |

## 5. Cobertura y contrato

**Cobertura por paquete** (`make verify-wf3-cobertura`, K25 ≥ 85 % en los 5 primeros):

| Paquete | Cobertura | Umbral |
|---|---|---|
| `internal/dominio` | 87.0 % | ≥ 85 % (cumple) |
| `internal/motor_fin` (S2) | 89.2 % | ≥ 85 % (cumple) |
| `internal/motor_cred` (S3) | 94.0 % | ≥ 85 % (cumple) |
| `internal/recetas` (S4) | 93.2 % | ≥ 85 % (cumple) |
| `internal/precios` (S5) | 88.8 % | ≥ 85 % (cumple) |
| `internal/api` | 82.0 % | sin umbral (K25) |
| `internal/auth` | 83.6 % | sin umbral |
| `internal/jobs` | 92.3 % | sin umbral |
| `internal/reportes` | 77.3 % | sin umbral |
| `internal/telegram` | 80.0 % | sin umbral |

**Operación → prueba** (`go test ./internal/api -run TestContrato_Operaciones -count=1 -v`, una subprueba `TestContrato_Operaciones/<operationId>` por operación, contada por extracto de C1). Esto sustituye las cifras viejas de `pruebas-contrato.md` (191 de 216, de una ronda anterior):

| Extracto | Cumplen | Extracto | Cumplen |
|---|---|---|---|
| `inicio.yaml` | 12 de 12 | `creditos-a.yaml` | 8 de 8 |
| `gastos.yaml` | 19 de 19 | `creditos-b.yaml` | 11 de 11 |
| `presupuesto-a.yaml` | 13 de 13 | `salud.yaml` | 6 de 6 |
| `presupuesto-b.yaml` | 15 de 15 | `reportes.yaml` | 10 de 10 |
| `comidas.yaml` | 16 de 16 | `ajustes-a.yaml` | 12 de 12 |
| `compras-tiendas-a.yaml` | 15 de 15 | `ajustes-b.yaml` | 11 de 11 |
| `compras-tiendas-b.yaml` | 16 de 16 | `ajustes-c.yaml` | 12 de 12 |
| | | `ajustes-d.yaml` | 13 de 13 |
| | | `ajustes-e.yaml` | 11 de 11 |
| | | `ajustes-f.yaml` | 16 de 16 |
| **Total** | **216 de 216** (0 fallas; «S6-03: 216 de 216 operaciones cumplen (2xx declarada, 4xx válidas contra openapi.yaml)») | | |

## 6. Pantallas: implementadas, probadas y capturadas (C3 v0.9)

Fuente: `docs/contratos/C3-integracion.md` v0.9, generada y comparada con `python3 -I frontend/scripts/c3-integracion.py` («matriz del documento = matriz generada: sí»). Son 41 pantallas de U1 (27 `P-xx` y 14 `P-26.xx`) más la fila Global.

| Columna | Sí | No | Nota |
|---|---|---|---|
| Implementado (la pantalla llama al menos una operación de cada necesidad y la ruta está registrada) | 35 | 7 | en «no»: P-02, P-03, P-05, P-12, P-17, P-18 y P-26.05 (bloqueo `C3-matriz` punto 2) |
| Probado (prueba Go de cada operación más Playwright con la API simulada) | 42 | 0 | la regla de «probado» por título es PROVISIONAL (`C3-matriz` punto 4) |
| Capturado (U5) | 0 | 42 | pendiente de WF-3b (`docs/design/capturas-u5/` todavía no existe) |

- Rutas: el frontend llama a 204 rutas; 0 fuera de `openapi.yaml` y 0 sin registrar en `routes.go` (C3-03). Hay 7 operaciones con pantalla en U1 que el frontend no llama en ningún lado, entre ellas `queryView` (`C3-matriz` punto 3).
- Humo real (C3-09, parcial): Inicio contra `https://finanzas.local:8443` con la CA en NSS, sin ignorar el certificado. Las 6 tarjetas cargan y las cifras $973.75, $6,905.50 (lectura B) y $10,000.00 (lectura A) son iguales a la API y al oráculo. Con la semilla, 5 de las 6 tarjetas salen vacías porque no hay datos de octubre (`C3-matriz` punto 5).
- Flujos con la API simulada: F-1 en ≤ 3 toques, F-2, F-4, F-5 y F-6 pasan. Contra el servidor real (C3-08) se hacen en WF-3b.

## 7. Artefactos no aprobados tras 3 rondas (hallazgos abiertos)

**frontend-core** (`REVISION-frontend-core-r3.md`): NO APROBADO, con 0 altos, 2 medios y varios bajos.

| Id | Sev. | Ruta | Qué pasa |
|---|---|---|---|
| N-04 | media | `frontend/src/componentes/Grafica.tsx:44` y `:76` | la descripción accesible automática de ECharts sale en inglés y con dinero sin formato en toda pantalla con gráfica (P-02, P-06 a P-10, P-15) |
| N-05 | media | `frontend/src/shell/Barras.tsx:42-56` (`BarraInferior`) | la barra inferior no sigue DESIGN §7.10: iconos de 24 px, fondo `ui.surface.raised` y píldora del destino activo |
| N-06 | baja | `src/app/App.tsx`, `src/shell/Layout.tsx` | falta una barrera de errores global («nunca una pantalla en blanco») |
| I-01 (nuevo, de este integrador) | baja | `frontend/pruebas/cola.spec.ts:104` | prueba intermitente: falló 1 vez de 3 corridas completas porque el botón «Revisar» se vuelve a pintar («element was detached from the DOM»); aislada pasó 32 de 32. Debe esperar a que la franja quede estable |
| N-03, H-07, H-08, H-09, H-11 | bajas | `Barras.tsx:17,26`, `Layout.tsx:47`, Inicio y API falsa | botones de 32 px con puntero fino, ancho del panel e Inicio sin `TarjetaMetrica` |

**frontend-areas** (`REVISION-frontend-areas-r3.md`): NO APROBADO, con 1 alto, 3 medios y 5 bajos. Lo que bloquea está en archivos de `frontend-base`.

| Id | Sev. | Ruta | Qué pasa |
|---|---|---|---|
| H-03 | **alta** | `src/vistas-bd/VistaCalendario.tsx:30`; `src/componentes/Grafica.tsx:48` y `:78` | el evento del calendario mide 40 × 44 en 412 táctil; «Ver como tabla» mide 20 px con puntero fino; la descripción emergente usa colores fijos en oscuro |
| H-07 | media | `src/vistas-bd/VistaBD.tsx`, `VistaTabla.tsx`; Ajustes 02, 04, 05 y 07 | falta la edición en grupo con `Ctrl+/` y vista previa (U4 §0.2 y §0.7); VistaBD no tiene selección de filas |
| H-08 | media | lo mismo en Ajustes 09, 13 y 14 | igual que H-07 |
| H-10 | media | `src/componentes/Boton.tsx:6`, `Grafica.tsx`, `Campo.tsx`, `Chip.tsx` | foco con transición, series de C2 recoloreadas al pasar el cursor, estilo de deshabilitado |
| R2-05, R2-06, R2-07, R3-01, R3-02 | bajas | `S05Creditos.tsx:158`, pruebas de ajustes-b, `S06Tiendas.tsx:128`, `validar-c1.ts` (`ajv` sin declarar), `CeldaEditable.tsx:39` | detalle en la revisión |

**backend-core** (APROBADO, r2) deja 4 bajos y **backend-areas** (APROBADO, r3) deja 8 bajos (B-1 a B-8), todos en sus revisiones; ninguno bloquea.

**PREGUNTA para Axel o el orquestador** (la hace también la revisión de frontend-areas): ¿se encarga a `frontend-base` una corrección de H-03, H-10, la base de `Ctrl+/`, N-04 y N-05 (más R3-01 y R3-02), con una revisión puntual, antes de WF-3b? ¿O se acepta con estos pendientes? Las correcciones son locales y no tocan cifras.

## 8. Bloqueos nuevos de WF-3a

Son **40 archivos, todos ABIERTOS**, dados de alta en `docs/estado/bloqueos/INDICE.md` (secciones «Bloqueos de la semilla de WF-3a» y «Bloqueos nuevos de WF-3a»):

- 19 `S6-*`: 3 de la semilla, más adjuntos, ajustes, comidas, compras (3), créditos, gastos e Inicio, importación de CSV, presupuesto (3), respaldos, sistema y valores operativos.
- 6 `S7-*`: convivencia con MultiX, fechas del estado de cuenta, fuentes de las alertas, lualatex, secciones del reporte mensual y ruta del token.
- 14 `WF3-*`: interfaz de las áreas, rutas, servidor y red.
- 1 `C3-matriz`.

Los que esperan a Axel con más urgencia son estos tres:

1. `C3-matriz`: las 7 filas en «no» y las 7 operaciones sin llamada; decide si C3 se puede firmar en WF-3b.
2. `S6-gastos-inicio` §15.4 y §15.5.
3. `S7-fuentes-avisos` §1: el día de despensa. Sin él, la alerta de precios viejos no sale.

## 9. Propuestas PROVISIONALES para la puerta 4 (K31)

Todos estos valores ya están en el código como PROVISIONAL, con su bloqueo citado. Cambiar cualquiera es una constante o un valor de `app_setting`, sin migración salvo donde se indica.

| Tema | Propuesta | Por qué | Bloqueo |
|---|---|---|---|
| Respaldos: hora | Diario a las **03:15** (America/Mexico_City) | A esa hora nadie captura; queda después del resumen de las 21:30 y antes del reporte mensual de las 08:00 | `S6-respaldos`, `S6-valores-operativos-base` |
| Respaldos: retención | Conservar **30** diarios | Un mes para volver atrás; cada copia pesa unos 2 MB (1,810,432 bytes medidos hoy) | `S6-respaldos` |
| Respaldos: restauración | Reemplaza toda la base salvo `backup` y `device_session`, que se reponen del estado previo; trabajo síncrono (202 ya resuelto) mientras la base pese poco | No perder el catálogo de respaldos ni revivir sesiones cerradas | `S6-respaldos` (PREGUNTA) |
| Adjuntos: tipos de archivo | **JPG, PNG, WebP y PDF**, detectados por el contenido y no por el nombre; otro tipo da 415. HEIC queda como opción b | Son los que ya declara C1 (U1 N-064) | `S6-adjuntos` |
| Adjuntos: tamaño | 10 MB con 1 MB = 1,048,576 bytes | Es la propuesta de S1 que cita S6-10 | `S6-adjuntos` |
| Cabeceras | `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`, `X-Frame-Options: DENY`; `Cache-Control: no-store` en la API; CSP `default-src 'self'; img-src 'self' data: blob:; style-src 'self' 'unsafe-inline'; connect-src 'self'; font-src 'self'; worker-src 'self'; frame-ancestors 'none'` | El frontend no usa scripts en línea ni recursos de terceros. Ya se envían; falta `Cache-Control` en los estáticos y el tipo del manifiesto (§11) | `S6-valores-operativos-base`, `WF3-frontend-servidor` |
| Contenido del resumen diario | Fecha, gastos de hoy, gastos del mes y pagos de los próximos 7 días | Contesta las tres preguntas de cada noche: cuánto gasté, cómo va el mes y qué viene | `S7-fuentes-avisos` (propuesta en `s7-avisos.md`) |
| Contenido del resumen semanal | Semana ISO, gastos de la semana, gastos del mes, pagos de los próximos 7 días y casillas sin armar de la semana siguiente | Lo mismo, más el plan de comidas que hay que armar | ídem |
| Contenido del reporte mensual | Resumen del mes, devengado contra efectivo, por categoría, por quincena, créditos y parte de Karen; comidas y canasta cuando haya cálculo compartido | Lo que ya cruza el PDF con la API (S7-02) | `S7-reporte-mensual-secciones` |
| Avisos: cadencia y reintentos | Los jobs revisan cada minuto contra la hora de `app_setting`; un envío fallido queda `failed` y no se reintenta solo | Un cambio en Ajustes se respeta al instante y no se manda dos veces | `S7-fuentes-avisos` §4 |
| Contraseñas | argon2id t=3, m=64 MiB, p=4 (RFC 9106 §4) | Tarda 0.1 a 0.2 s por inicio de sesión en la máquina de casa | `S6-valores-operativos-base` |
| Límites del cuerpo | 8 MiB por JSON (413 al pasarse); 32 MiB en memoria con Idempotency-Key | Las cargas grandes, como PROFECO, pasan por un temporal | `S6-valores-operativos-base` |
| Purga de idempotencia | Cada hora, al minuto 5 | C1 dice «un job la purga» sin horario | `S6-valores-operativos-base` |
| PWA sin conexión | Sondeo de `/health` cada 15 s; reintento de la cola cada 30 s; la copia de lecturas se guarda hasta cerrar sesión; la cola no se borra al cerrar sesión | U1 §2.4: nunca se pierde una captura | `WF3-frontend-k31` |

## 10. Pendientes

**Para WF-3b:**

1. **U5:** capturas de las 42 filas en las 4 combinaciones (`docs/design/capturas-u5/`) y la columna «capturado» de C3.
2. **Flujos:** C3-08 con F-1 a F-6 contra el servidor real.
3. **C3-09 completo:** 10 o más cifras, la deuda viva de $6,293.44 y P1 a P6.
4. **Firma de C3 (C3-07):** la hace un revisor independiente, después de decidir los 5 puntos de `C3-matriz` y regenerar la matriz con `c3-integracion.py --escribir`.
5. **Semilla:** presupuesto e ingreso de octubre, para que Inicio tenga cifras en las 6 tarjetas al capturar (`C3-matriz` punto 5).
6. **Corrección del frontend:** la de §7, si Axel la pide.
7. **Makefile:** los objetivos `verify-wf3-c3` y `verify-wf3-humo` que pidió `integrador-frontend-c3`, y meter `verify-wf3` en `verify-all` cuando el frontend se apruebe.
8. **Servidor:** las solicitudes de §11.

**Para WF-4:**

1. Unidad systemd, rutas de producción (`FINANZAS_*`) y binario con `dist` embebido (`WF3-frontend-servidor` opción A).
2. ACL de lectura del token de MultiX (S7-06, K30) y ruta exacta de la clave (`S7-ruta-token`).
3. IP fija y nombre `finanzas.local` permanentes; CA en el almacén del sistema y en los dos Android (`WF3-integrador-red`, `S6-sistema`).
4. `ops/respaldo.sh`, el temporizador y una segunda copia de respaldos (`ops/runbook.md` §6).
5. `lualatex` con `luaotfload` en la máquina de casa (`S7-lualatex`).
6. S7-07 real (orquestador) y S6-11 en Android (Axel).

## 11. Solicitudes fuera de alcance (consolidadas; ningún agente de WF-3a las aplicó)

1. `backend/internal/api/http.go` (`staticHandler`), dueño `backend-base` o `auth-sistema`: servir `*.webmanifest` como `application/manifest+json` (hoy sale `text/plain; charset=utf-8` con `nosniff`). Agregar `Cache-Control: no-cache` a `index.html`, `sw.js`, `manifest.webmanifest` y al respaldo SPA, y `Cache-Control: public, max-age=31536000, immutable` a `/assets/*`. Lo piden `integrador-backend` e `integrador-frontend-c3`.
2. `docs/estado/bloqueos/WF3-frontend-servidor.md`: anotar que los puntos 1 y 2 ya se comprobaron en el servidor vivo, que el 4 (CSP) ya se envía y que el 3 sigue pendiente.
3. `docs/estado/bloqueos/S6-respaldos.md` punto 1: la propuesta ya está en `ops/runbook.md` §4 como PROVISIONAL.
4. `Makefile` (`backend-base`): `verify-wf3-c3` y `verify-wf3-humo` (texto exacto en `integrador-frontend-c3.md`, Solicitudes 1).
5. `frontend/package.json`: declarar `"ajv": "8.20.0"` en `devDependencies` y mover `validar-c1.ts` y `tablero.ts` a `pruebas/apoyo/` (R3-01).

## 12. Paquetes instalados y descargas

| Origen | Paquete | Quién | Nota |
|---|---|---|---|
| pacman | `mkcert 1.4.4-3` (repositorio extra) | `integrador-backend` | Con la única orden sudo permitida (`sudo pacman -S --noconfirm mkcert`). La CA se instaló solo en NSS (`TRUST_STORES=nss`), porque el `-install` completo pide otro sudo |
| Go | `github.com/pressly/goose/v3 v3.28.0`, `golang.org/x/crypto v0.58.0`, `github.com/getkin/kin-openapi v0.149.0` (solo pruebas), más sus indirectos; `golang.org/x/sys` de 0.48 a 0.49 | `backend-base` | `go.mod` y `go.sum` |
| npm | react 19.3.0, react-dom 19.3.0, radix-ui 1.8.0, class-variance-authority 0.7.1, clsx 2.1.1, tailwind-merge 3.7.0, @tanstack/react-table 9.2.8, @dnd-kit/core 6.3.1, @dnd-kit/sortable 10.0.0, @dnd-kit/utilities 3.2.2, echarts 6.1.0, lucide-react 1.55.0 | `frontend-base` | dependencias |
| npm (dev) | vite 8.3.4, @vitejs/plugin-react 6.1.2, typescript 5.9.3, tailwindcss 4.3.3, @tailwindcss/vite 4.3.3, vite-plugin-pwa 2.0.0, workbox-window 7.4.1, openapi-typescript 7.13.0, msw 3.0.3, @playwright/test 1.60.0, @axe-core/playwright 4.13.0, @types/react 19.3.0, @types/react-dom 19.3.0, @types/node 24.10.1, yaml 2.9.1 | `frontend-base` | 519 paquetes con los transitivos; `package-lock.json` versionado |
| npm (dev) | storybook 10.6.1, @storybook/react-vite 10.6.1 (y 84 transitivos) | `front-variantes` | `npm audit` sin vulnerabilidades |
| Descargas | Ninguna en ejecución. Inter (woff2) sale de `docs/design/fuentes/` (ya versionado). El Chromium de Playwright 1.60 ya estaba en la máquina (`FH_CHROMIUM`) | — | La app no carga nada de la red: 0 hosts externos (U5-08) |

Los demás agentes no instalaron nada. Los revisores usaron `uv run --with pyyaml --with jsonschema` solo para ejecutar, y el `ajv` que ya traía `node_modules`.

## 13. Consumo por agente (aproximado, de cada reporte)

| Agente | Acciones | KB leídos | KB escritos |
|---|---|---|---|
| backend-base | 55 | 235 | 205 (por encima del tope de 150; más 183 generados) |
| auth-sistema | 51 | 165 | 140 |
| semilla | 56 | 330 | 75 |
| api-gastos-inicio | 57 | 180 | 100 |
| api-presupuesto | 56 | 185 | 125 |
| api-comidas-salud | 45 | 330 | 125 |
| api-compras-a | 54 | 200 | 125 |
| api-compras-b | 52 | 260 | 125 |
| api-creditos | 52 | 200 | 140 |
| api-ajustes | 48 | 200 | 125 |
| s7-avisos | 56 | 270 | 115 |
| s7-pdf | 60 | 240 | 110 |
| pruebas-contrato | 48 | 230 | 85 |
| frontend-base | 51 | 175 | 189 (por encima del tope; sin contar los generados) |
| front-inicio | 48 | 220 | 130 |
| front-gastos-presupuesto | 54 | 165 | 145 |
| front-comidas-compras | 52 | 170 | 148 |
| front-creditos-salud-reportes | 58 | 160 | 150 |
| front-ajustes-a | 58 | 215 | 190 (por encima del tope) |
| front-ajustes-b | 59 | 200 | 175 (por encima del tope) |
| front-variantes | 40 | 150 | 65 |
| correccion-backend-core-r1 | 60 | 150 | 30 |
| correccion-backend-areas-r1 / r2 | 38 / 36 | 170 / 115 | 30 / 30 |
| correccion-frontend-core-r1 / r2 | 45 / 35 | 170 / 110 | 40 / 22 |
| correccion-frontend-areas-r1 / r2 | 60 / 52 | 240 / 150 | 45 / 40 |
| revisor-backend-core-r1 / r2 | 42 / 58 | 200 / 150 | 32 / 32 |
| revisor-backend-areas-r1 / r2 / r3 | 60 / 58 / 52 | 150 / 170 / 190 | 30 / 30 / 60 |
| revisor-frontend-core-r1 / r2 / r3 | 44 / 46 / 50 | 120 / 120 / 150 | 30 / 30 / 45 |
| revisor-frontend-areas-r1 / r2 / r3 | 55 / 57 / 40 | 230 / 150 / 120 | 25 / 35 / 45 |
| integrador-backend | 46 | 170 | 50 |
| integrador-frontend-c3 | 45 | 95 | 65 |
| integrador-wf3a | 50 | 160 | 60 |
| **Total (42 corridas)** | **2,139** | **7,760 (~7.8 MB)** | **3,593 (~3.6 MB)** |

Cuatro agentes de construcción se pasaron del tope de 150 KB escritos del protocolo §4 (backend-base, frontend-base, front-ajustes-a y front-ajustes-b). Lo declararon en su reporte y no cortaron trabajo. **PROPUESTA:** que el tope de escritura de los agentes de construcción de WF-3b sea de 200 KB, sin contar el código generado.

## 14. Propuestas para Axel

| Qué decidir | Resultado | Costo | Riesgo | Recomendación |
|---|---|---|---|---|
| Corregir el frontend antes de WF-3b (§7) | Las pantallas cumplen los 44 px, los colores en oscuro y la edición en grupo que promete U4 | Una corrección de `frontend-base` más una revisión puntual | Bajo: no toca cifras | Sí, antes de las capturas de U5, para no capturar dos veces |
| Aprobar los valores de §9 en la puerta 4 | Quedan firmes respaldos, adjuntos, cabeceras y el contenido de los avisos | Solo revisar la tabla | Ninguno en cifras | Aprobar como están, y decidir aparte si se admite HEIC |
| Reservar la IP 192.168.50.56 en el módem | `finanzas.local` siempre apunta a la misma máquina | Unos minutos en el módem | Sin reservarla, el teléfono puede perder la app tras un reinicio | Sí |
| Semilla con presupuesto e ingreso de octubre | Inicio muestra cifras en las 6 tarjetas para las capturas | Bajo | Ninguno | Sí, antes de WF-3b |
| Probar en Android la CA de mkcert (S6-11) | Confirma que la PWA se instala sin aviso de certificado | 10 minutos con `ops/mkcert.md` §4 | Si falla, la app no se instala en el teléfono | Hacerlo antes de WF-3b |

## 15. WF-3a bis (12-oct-2026): ronda acotada de corrección del frontend

Agente `integrador-wf3a-bis` (Opus). El orquestador encargó esta ronda antes de WF-3b, con el precedente de P1d y P1e. Nada de esto cambia cifras de los motores ni decisiones de Axel. Las decisiones del orquestador sobre `docs/estado/bloqueos/C3-matriz.md` fueron: punto 1 (a), fila P-26 «Local (contenedor)»; punto 2, implementar las necesidades sin llamada; punto 3 (b), N-064, N-173 y N-175, y (c), excepción técnica PROVISIONAL solo donde la pantalla ya cubre la necesidad; punto 4 (a) ahora y (b) en WF-3b; punto 5 (b), semilla de desarrollo. Todas siguen PROVISIONALES hasta la puerta 4.

Trabajaron: `frontend-base-fix`, `front-ajustes-a-fix`, `front-ajustes-b-fix`, `front-comidas-compras-fix`, `front-creditos-salud-reportes-fix`, `front-gastos-presupuesto-fix`, `semilla-dev`, `correccion-frontend-core-r4`, `correccion-frontend-areas-r4` y `c3-regen`. Revisaron (rondas 4 y 5): `revisor-frontend-core-r4/r5` y `revisor-frontend-areas-r4/r5`. Sus reportes están en `docs/estado/wf3/`.

### 15.1 En corto

- **HECHO.** Quedaron cerrados el alto y los 5 medios abiertos al cierre de WF-3a (§7): objetivos táctiles y calendario (H-03), gráficas accesibles en español (N-04), barra inferior (N-05), edición en grupo con `Ctrl+/` y vista previa (H-07, H-08), foco y estados (H-10). También se cerraron los 4 medios que salieron en la ronda 4 y la prueba intermitente de la cola (I-01).
- **HECHO.** C3 v0.95 (borrador sin firma): implementado **42 de 42**, probado 42 de 42 y capturado 0 de 42 (WF-3b). `listPeriods` es la única operación sin llamada y queda como excepción técnica PROVISIONAL, con su cobertura comprobada por el script.
- **HECHO.** Semilla de desarrollo (`finanzas semilla --desarrollo`). El servidor de desarrollo quedó re-sembrado y **corriendo** (PID 2419156), y las 6 tarjetas de Inicio responden con cifras (§15.6).
- **PREGUNTA.** El frontend **no quedó aprobado tras 5 rondas**. Quedan 2 medios y 7 bajos nuevos o reabiertos, más bajos viejos (§15.4). Ninguno toca cifras.
- **PREGUNTA.** `make verify-wf3` queda en **rojo** (rc 2) por una sola prueba: `pruebas/gastos-presupuesto/correccion-bis.escritorio.spec.ts:79` [escritorio]. Falla en **cada** corrida completa: 2 de 2 de `c3-regen` y 2 de 2 de este integrador. Aislada pasa. Es R5-05 de `REVISION-frontend-core-r5`, que el revisor clasificó como baja porque la vio pasar sola. `make verify-all` pasa (rc 0).

### 15.2 Verificación (corrida por este integrador, 12-oct-2026, 08:13–08:40)

| Paso | Resultado |
|---|---|
| `cd frontend && npm run build` | rc 0 (preparar, `tsc --noEmit` y vite; precache 81 entradas) |
| `npx playwright test` completo (API simulada, JSON) | 235 ejecutadas: **230 pasan, 4 omitidas, 1 falla** (5.1 min). Omitidas: 3 de `Ctrl+/` en el proyecto móvil (atajo de teclado) y `ajustes-a/ajustes.escritorio.spec.ts:372` en escritorio (es de 412 táctil). La falla es `gastos-presupuesto/correccion-bis.escritorio.spec.ts:79` [escritorio], «P-02 Ctrl+/ · editar el tipo de varias filas…»: «Expected: 2 / Received: 0» en `expect.poll(() => parches.length)` tras `Enter` |
| `npx playwright test pruebas/cola.spec.ts --repeat-each=3` | **12 passed** (3 corridas × 4 pruebas, 1.0 min). I-01 ya no se reproduce |
| `node scripts/rutas-api.mjs` | «C3-03 frontend: 208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)» |
| `node scripts/hosts-externos.mjs dist` | «U5-08 dist: hosts externos que se pedirían = 0» (solo literales que no son peticiones: w3.org, react.dev, tailwindcss.com, licencia de Inter y bit.ly en workbox) |
| `cd backend && go test ./... -count=1` | rc 0: 32 paquetes OK, 5 sin pruebas, 0 fallas (entre ellos `internal/semilla`, con las 8 pruebas nuevas de la semilla de desarrollo) |
| `make -k verify-wf3` | **rc 2**. Pasan: Go (S6-01: build, vet, gofmt, `go mod tidy -diff` y 32 paquetes OK), K25 (S6-02: dominio 87.0 %, motor_fin 89.2 %, motor_cred 94.0 %, recetas 93.2 % y precios 88.8 %) y contrato (S6-03: «los 23 archivos coinciden byte a byte», «216 operaciones en 172 rutas, 0 diferencias»). Falla `verify-wf3-frontend` (`npm run verify`): build, rutas, hosts y manifiesto OK; Playwright 230 / 4 / **1 falla**, la misma prueba |
| `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | **rc 0**: «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)» |
| `go test ./internal/api -run 'TestContrato_Operaciones\|TestContrato_Rutas'` + `python3 -I frontend/scripts/c3-integracion.py --pw-json … --go-log …` (sin `--escribir`) | 216 subpruebas PASS. «implementado sí 42/42; probado sí 42/42; capturado sí 0/42», «sin pantalla: 5 (sin justificar: 0)», «no llama: 1 ['listPeriods']; sin excepción o sin cobertura: 0», **«Matriz del documento = matriz generada: sí»** |
| Servidor: `bash ops/dev-stop.sh`; `finanzas migrate`; `finanzas semilla`; `finanzas semilla --desarrollo`; `bash ops/dev-run.sh` | Todo rc 0: «migraciones aplicadas: versión 4»; semilla real con todo «ya estaba» (idempotente); semilla de desarrollo con 0 nuevos y «ya estaban» 1 presupuesto, 1 partida, 1 meta, 1 aportación, 36 + 36 precios y 12 comidas; servidor en `127.0.0.1:8443` y `192.168.50.56:8443`, PID 2419156 |
| Humo con la CA de mkcert (script propio, sin imprimir credenciales) | `GET /` 200 e **igual byte a byte** a `frontend/dist/index.html`; 3 recursos de `/assets/` y `/sw.js` iguales a los del disco. Inicio de sesión como Axel: 200. Tarjetas en §15.6 |
| Secretos en lo que entra al commit (script propio, solo conteos) | 0 llaves privadas, 0 tokens de bot, 0 llaves de API, 0 hashes argon2 y 0 contraseñas de desarrollo. 0 archivos `.db`, certificados, `node_modules`, `dist` o `tsbuildinfo` |

### 15.3 Hallazgos cerrados (con la ronda en que se cerraron)

| Revisión | Id | Sev. | Qué era | Cerrado en |
|---|---|---|---|---|
| frontend-areas | H-03 | alta | evento del calendario de 40 × 44 en 412, «Ver como tabla» de 20 px y descripción emergente con colores fijos | r4 (`frontend-base-fix`) |
| frontend-areas | H-07, H-08 | media | edición en grupo con `Ctrl+/` y vista previa en Ajustes 02, 04, 05, 07, 09, 13 y 14 | r4 (`frontend-base-fix`, `front-ajustes-a-fix`, `front-ajustes-b-fix`) |
| frontend-areas | H-10 | media | foco con transición, series recoloreadas y estilo de deshabilitado | r4 |
| frontend-areas | R2-05, R2-06, R2-07, R3-02 | bajas | detalles de Ajustes y celda sin nombre | r4 |
| frontend-areas | R4-01, R4-02 | media | P-26.02 sin V-17 ni «No se pudo guardar»; P-26.13 sin aviso | r5 (`correccion-frontend-areas-r4`) |
| frontend-core | N-04 | media | descripción de ECharts en inglés | r4 |
| frontend-core | N-05 | media | barra inferior fuera de DESIGN §7.10 | r4 |
| frontend-core | N-06, N-03 | bajas | barrera de errores global; botones de 32 px | r4 |
| frontend-core | I-01 | baja | `cola.spec.ts` intermitente | `frontend-base-fix`; aquí, 12 de 12 en 3 corridas |
| frontend-core | R4-01, R4-02 | media | P-03 sin pestaña Ingreso (N-102) ni líneas de crédito (N-053); eje compacto y etiquetas de dinero en las gráficas | r5 (`correccion-frontend-core-r4`) |

### 15.4 Hallazgos abiertos (frontend-core r5 y frontend-areas r5)

Los dos revisores usan el id «R5-01» para hallazgos distintos; aquí van como «core R5-01» y «áreas R5-01».

| Id | Sev. | Ruta | Dueño | Qué pasa |
|---|---|---|---|---|
| core R5-01 | **media** | `frontend/src/pantallas/comidas-compras/P10Plan.tsx:245` | front-comidas-compras | P-10 Tendencias rotula las semanas «Lun 5 oct-2026»: formato híbrido fuera de U5-07. La corrección es `fechaConDia()` o `fecha()` (CONVENCIONES-frontend §5) |
| áreas R5-01 | **media** | `frontend/src/pantallas/ajustes-a/comunes.tsx:228` (`DialogoCritico`) y 6 manejadores en S02, S03 y S05 | front-ajustes-a | guardados y borrados de P-26.01, 02, 03, 05 y 07 que fallan sin avisar (U4 §0.3); V-18 nunca aparece al eliminar una promoción en uso |
| R5-05 (core) | baja según el revisor; aquí **bloquea `make verify-wf3`** | `frontend/pruebas/gastos-presupuesto/correccion-bis.escritorio.spec.ts:79-97` | front-gastos-presupuesto | falla en 4 de 4 corridas completas y pasa aislada. El `Enter` llega con el foco en el `<select>` del bloque y no produce los 2 PATCH. Sugerencia del revisor: esperar el foco del bloque `edicion-varias` o tocar «Aplicar», y revisar dónde escucha `Enter` `EdicionVarias.tsx:55` |
| R5-02 (core) | baja / PREGUNTA | `frontend/src/shell/HojaRegistro.tsx:89` | frontend-base | con «A crédito», «Guardar» no hace nada ni avisa |
| R5-03 (core) | baja / PREGUNTA | `gastos-presupuesto/hoja-registro.tsx:36`, `:54` | front-gastos-presupuesto | N-053 también ofrece préstamos a plazo y tarjetas; ya está en `WF3a-bis-gastos-presupuesto.md` punto 5 (PROVISIONAL) |
| R5-04 (core) | baja | `frontend/src/pantallas/inicio/GraficaInicio.tsx:76` | front-inicio | leyenda en `ui.text.secondary` en lugar de `ui.text.default` |
| H-13 (áreas, reabierto) | baja | `creditos-salud-reportes/P20Comparador.tsx:67`, `:113` | front-creditos-salud-reportes | con 503 `unavailable`, P-20 dice «tiempo máximo agotado» en lugar de «Comparador no disponible» (N-226) |
| R3-01 (áreas) | baja | `frontend/package.json` | frontend-base | `ajv` sin declarar en `devDependencies` |
| R4-03, R4-04, R4-05, H-07, H-08, H-09 (core) | bajas | `GraficaInicio.tsx` (descripción automática en inglés dentro de `aria-hidden`), alto de la barra inferior, `api-falsa.ts:41`, `Layout.tsx:48 max-w-[85vw]`, Inicio sin `TarjetaMetrica` | frontend-base, front-inicio | sin cambios |

### 15.5 C3 v0.95

`docs/contratos/C3-integracion.md` v0.95 la regeneró `c3-regen` con `python3 -I frontend/scripts/c3-integracion.py --escribir`. Este integrador la volvió a generar sin `--escribir` con sus propias corridas, y salió igual: «Matriz del documento = matriz generada: sí».

- Implementado: **42 de 42** (antes 35). Las 7 filas en «no» (P-02, P-03, P-05, P-12, P-17, P-18 y P-26.05) ya llaman sus operaciones: `deleteAttachment` en P-04 («Quitar adjunto»), `listNutritionFacts`, `updateProduct` y `updatePresentation` en P-15, `queryView` en P-02 (tablero, calendario y galería) y P-04, y `listIncomes` en P-08, entre otras.
- Probado: 42 de 42, con la regla PROVISIONAL de `C3-matriz` punto 4 (a). Capturado: 0 de 42 (WF-3b).
- La fila P-26 dice «Local (contenedor de P-26.01…P-26.14)» (punto 1 (a), PROVISIONAL).
- `listPeriods` es la única operación con pantalla que no se llama. Es excepción técnica PROVISIONAL (punto 3 (c)): N-080 lo cubre `getBudgetMonth`, N-082 `getBudgetFortnights` y N-100 `getIncomeMonth`/`listIncomes`. El script falla (rc 1) si aparece otra sin excepción escrita o sin cobertura.
- Sigue **sin firma** (C3-07). La firma va en WF-3b, con las capturas U5, C3-08 (flujos contra el servidor real) y C3-09 completo.

### 15.6 Semilla de desarrollo (`finanzas semilla --desarrollo`, agente `semilla-dev`)

Es una capa idempotente sobre la semilla real (`backend/internal/semilla/desarrollo.go`, banderas en `backend/cmd/finanzas`). Exige la semilla real ya cargada, marca todo con «(desarrollo)» y se niega si `FINANZAS_ENTORNO=produccion` (candado PROVISIONAL). Los montos **no están en el código**: los lee en ejecución de las fixtures del oráculo S2 (`oraculo/fixtures/s2`).

| Qué siembra | Valor | De qué fixture |
|---|---|---|
| Sueldo quincenal de Axel (partida planificada desde el 15-oct) | $12,000.00 el 15 y $12,500.00 el 30 | `s2-proyeccion-csv`, `planned_items[1]` |
| Presupuesto de Despensa desde el 1-oct-2026 | $2,000.00 | `s2-presupuesto-scope`, `budgets[1]` |
| Meta de ahorro (desarrollo) y aportación del 5-oct | $6,293.44 al 10-abr-2027, plan $449.53 al mes; aportación $3,146.72 (50.0 %) | `s2-metas` |
| Precios de la compra del 10-oct | 36 de Bodega Aurrera (renglones del ticket, `ticket_manual`) y 36 de Walmart (carrito comparado de `datos/semilla/data.json`, `manual`) | ticket y carrito de la semilla |
| Plan de comidas del 12 al 17-oct | 12 entradas, 6 recetas del recetario | sin fixture: dato de desarrollo |
| Renta $8,000.00 mensual | solo con `--con-gasto-fijo`; por omisión **no** se siembra, porque movería P1–P6 | `s2-proyeccion-csv`, `planned_items[2]` |

Inicio con la base re-sembrada (API con sesión de Axel; centavos):
- **Mes contra plan:** real 119,764 de 200,000, `normal`, quedan 80,236.
- **Quincena 1–15 oct:** ingreso 1,200,000 (esperado), gastado 0, saldo 1,200,000.
- **Comidas de la semana 12–18 oct:** responde, pero con costo 0 (bloqueo `WF3a-bis-semilla-dev` punto 5).
- **Mejor canasta:** 309,450, con 26,450 menos que la segunda tienda.
- **Ahorro:** 314,672 ahorrado y plan 44,953 al mes, de una meta de 629,344.
- **Próximo pago:** 30-oct, 97,375 (hogar 59,882 y Karen 37,493); disponible B 690,550.

Resultado: **5 de 6 tarjetas con montos distintos de 0**; antes de esta ronda eran 1 de 6. `getAccrualVsCash` queda igual que con la semilla real sola.

### 15.7 Bloqueos nuevos de WF-3a bis

Hay **4 archivos nuevos, todos ABIERTOS**, dados de alta en `docs/estado/bloqueos/INDICE.md`, sección «Bloqueos nuevos de WF-3a bis». La carpeta tiene 129 archivos.
- `WF3a-bis-ajustes-a.md`: 3 puntos de `Ctrl+/` en Ajustes.
- `WF3a-bis-comidas-compras.md`: 4 puntos. El principal: las cifras de P-14 con la semilla de desarrollo (Walmart $3,359.00) no son las canon de U1 §3.16 ($3,233.00).
- `WF3a-bis-gastos-presupuesto.md`: 7 puntos. Entre ellos, los filtros que le faltan a `queryView`, el id del adjunto en `Expense`, N-053 y `listPeriods`.
- `WF3a-bis-semilla-dev.md`: 7 puntos. Entre ellos, la renta contra P1–P6, P3–P6 de la semilla real contra model.json, y el costo y las kcal de la semana en 0.

`C3-matriz.md` sigue ABIERTO, con las decisiones del orquestador anotadas y aplicadas, pero PROVISIONALES.

### 15.8 Pendientes y preguntas

- **PREGUNTA (orquestador).** ¿Se hace otra ronda mínima antes de WF-3b, o se pasa a WF-3b con estos pendientes anotados? Esa ronda tendría 3 cambios locales, sin cifras:
  1. core R5-01: `fechaConDia()` en `P10Plan.tsx:245`, más una aserción en la prueba.
  2. áreas R5-01: `catch` en `DialogoCritico` y en los 6 manejadores, V-18 con «Desactivar», y 4 pruebas.
  3. R5-05: la prueba `correccion-bis.escritorio.spec.ts:79`. Es la única que deja `make verify-wf3` en rojo.
- **PREGUNTA (Axel, puerta 4).** R5-02: con «A crédito», ¿«Guardar» se deshabilita con su motivo o se esconde? R5-03 / N-053: ¿se ofrecen también préstamos a plazo y tarjetas?
- **PENDIENTE (WF-3b).** Capturas U5, C3-08 con F-1 a F-6 contra el servidor real, C3-09 completo (antes, decidir `WF3a-bis-semilla-dev` punto 4) y firma de C3.
- **Solicitudes fuera de alcance, sin aplicar:**
  1. `ops/dev-run.sh`, rama `--sembrar`: agregar `"$BIN" semilla --desarrollo -db "$DB" -dir "$ROOT/datos/semilla"` después de `finanzas semilla`. Hoy se aplica a mano (`semilla-dev`).
  2. `ops/runbook.md`: «`finanzas semilla --desarrollo` solo en desarrollo; producción lleva `FINANZAS_ENTORNO=produccion`».
  3. `frontend/package.json`: agregar `"ajv": "8.20.0"` a `devDependencies` (R3-01).
  4. `.gitignore`: agregar `frontend/tsconfig.tsbuildinfo`. Es un archivo generado (06:46) que no entró al commit.

### 15.9 Paquetes instalados y consumo

Paquetes: ninguno en toda la ronda (cada reporte dice «Ninguno»). Consumo de este integrador: unas 58 acciones, unos 180 KB leídos y unos 30 KB escritos.
