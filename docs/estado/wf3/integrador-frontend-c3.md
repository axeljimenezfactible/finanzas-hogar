# Integrador frontend / C3 · WF-3a

Agente `integrador-frontend-c3` (Opus) · 2026-10-12 · HEAD «[Puerta 2 bis] K53» (sin git add ni commit). Criterios: `docs/criterios/C3.md` 1.1 (y S6-12, U5-08 en lo que toca a `dist`). Estado recibido: backend-core y backend-areas aprobados; frontend-core y frontend-areas **no aprobados** (3 rondas cada uno; hallazgos abiertos en `REVISION-frontend-*-r3.md`, no son de este agente). Servidor de desarrollo del integrador-backend en marcha (PID 1844644); no lo detuve ni lo reinicié.

## Hecho
- **HECHO** · `npm ci` (lockfile, sin paquetes nuevos), `npx tsc --noEmit` y `npm run build` limpios; Playwright completo con mocks: **149 pasan, 1 omitida, 0 fallan** (150 en 3.5 min).
- **HECHO** · `frontend/scripts/c3-integracion.py` (python3 -I, sin red): extrae del frontend las rutas que llama (cliente `pedir/pedirPagina/todas/encolar/usarConsulta`, más `armarUrl` y `fetch` literales; sin mocks ni `*.gen.ts`), las cruza con `openapi.yaml` y `routes.go` (C3-03), arma la matriz desde U1 §2.1/§2.3/§3.29 y las `x-necesidades` de C1 (C3-02), lee los resultados de Go y Playwright (C3-04) y escribe o compara la matriz de `C3-integracion.md`. Complementa a `scripts/rutas-api.mjs` de frontend-base (que solo ve las 5 funciones del cliente).
- **HECHO** · El binario sirve `frontend/dist`: `GET /` por `https://finanzas.local:8443` con la CA de mkcert = 200 y es byte a byte `dist/index.html`; manifiesto válido; íconos, `sw.js` y el JS principal con 200. 0 hosts externos en `dist`.
- **HECHO** · Humo real con Playwright contra `https://finanzas.local:8443` (`frontend/pruebas-integracion/humo-inicio.spec.ts`, `frontend/playwright.integracion.config.ts`): entra como `axel` con la contraseña de desarrollo (leída de `~/.config/finanzas-hogar/dev-semilla.txt`; no aparece en código, salidas ni en este reporte), abre Inicio, las 6 tarjetas cargan de la API real y aparecen $973.75, $6,905.50 (lectura B) y **$10,000.00 (lectura A)** = API = oráculo. Pasa en 6 s.
- **HECHO** · `docs/contratos/C3-integracion.md` **v0.9 BORRADOR sin firma**: definiciones, 41 filas `P-` + Global, excepciones una por una, operaciones no llamadas, faltantes, humo y bloque de firma vacío.
- **PREGUNTA** · Bloqueo nuevo `docs/estado/bloqueos/C3-matriz.md` (5 puntos: fila P-26, 7 filas con una necesidad sin llamada, 7 operaciones no llamadas, regla de «probado», tarjetas vacías con la semilla).
- **PROPUESTA** · Scripts npm `verificar:c3` y `test:integracion` en `frontend/package.json`; objetivos de Makefile en Solicitudes.

## Archivos
| Ruta | Cambio |
|---|---|
| `frontend/scripts/c3-integracion.py` | nuevo: C3-02, C3-03, C3-04 y matriz (`--escribir`, o comparar si no se pasa) |
| `frontend/playwright.integracion.config.ts` | nuevo: Playwright contra el servidor en marcha (sin webServer ni mocks, sin `ignoreHTTPSErrors`) |
| `frontend/pruebas-integracion/humo-inicio.spec.ts` | nuevo: humo real de Inicio (C3-09 parcial) |
| `frontend/package.json` | + `verificar:c3`, `test:integracion` (sin dependencias nuevas) |
| `docs/contratos/C3-integracion.md` | nuevo, v0.9 borrador |
| `docs/estado/bloqueos/C3-matriz.md` | nuevo |
| `docs/estado/wf3/integrador-frontend-c3.md` | este reporte |
`frontend/dist` se reconstruyó con `npm run build` (el servidor lo lee del disco en cada petición).

## Verificación
- **Paso 1** · `npm ci` rc 0 (4 s; aviso: el `postinstall` de esbuild 0.28.2 queda bloqueado por `allowScripts` de npm 12 y no hace falta: el build pasa). `npx tsc --noEmit` rc 0. `npm run build` rc 0 (`built in 404ms`, precache 80 entradas). `npx playwright test` (con mocks): `expected 149, skipped 1, unexpected 0, flaky 0` (omitida: `ajustes-a/ajustes.escritorio.spec.ts:317` en el proyecto escritorio, que es de 412 táctil).
- **C3-01** · `grep -c "^| P-" docs/contratos/C3-integracion.md` = **41** = T (27 `P-xx` de U1 §2.1 + 14 `P-26.xx` de U1 §2.3) + 1 fila Global. Encabezado con las 3 definiciones (§1). Cada estado en dos columnas: «sí/no» y evidencia.
- **C3-02** · `python3 -I frontend/scripts/c3-integracion.py --pw-json … --go-log …` → «pantallas U1: 41 (P-xx 27 + P-26.xx 14); filas: 42», «pantallas sin fila: 0; filas sin operación ni Local: 0», «operaciones: 216; sin pantalla: 5 (excepciones: 5, de ellas K53: 4; sin justificar: 0)». Las 4 de K53 llevan `x-agente-externo`; `getOpenApiDocument` es la «Excepción técnica» que declara C1. «Matriz del documento = matriz generada: sí», rc 0. P-26 como «Local (contenedor)» es PROVISIONAL (bloqueo C3-matriz, punto 1).
- **C3-03** · mismo script: «204 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (routes.go: 216; openapi: 216 operaciones)». `npm run build` sin errores. Implementado = sí en **35/42** filas; las 7 en «no» (P-02, P-03, P-05, P-12, P-17, P-18, P-26.05) tienen una necesidad cuya operación de C1 no llama la pantalla (detalle en la matriz y en el bloqueo, punto 2). Ninguna falta por ruta sin registrar.
- **C3-04** · `go test ./internal/api -run 'TestContrato_Operaciones|TestContrato_Rutas' -count=1 -v` → `ok` 1.77 s; **216 PASS / 0 FAIL** en subpruebas `TestContrato_Operaciones/<operationId>`; «216 de 216 operaciones cumplen»; «servidor 216, sobrantes 0, faltantes 0, sin implementar 0». Playwright: 150 ejecutadas, 149 pasan. Probado = sí en **42/42** filas (regla de títulos PROVISIONAL, bloqueo punto 4).
- **C3-05** · Capturado = no en 42/42: «pendiente WF-3b». `docs/design/capturas-u5/` todavía no existe.
- **C3-06** · `grep -c "| no |"` = 42 (todas por «capturado», 7 también por «implementado»); todas citan WF-3b o `C3-matriz`. No se firma en este estado.
- **C3-07** · Sin firma (borrador v0.9); el bloque de firma (§5) queda vacío para el revisor independiente.
- **C3-08** · No se ejecutó contra el servicio (es de WF-3b). Con mocks pasan los flujos F-1 (≤ 3 toques: `gastos-presupuesto.escritorio.spec.ts:133`, `inicio.escritorio.spec.ts:112`), F-4, F-5 y F-6 (`creditos-salud-reportes/flujos.spec.ts`).
- **C3-09 (parcial)** · `npx playwright test -c playwright.integracion.config.ts` → **1 passed (6.5 s)**, sin `ignoreHTTPSErrors`: Chromium confía en la CA porque mkcert la puso en NSS (`~/.pki/nssdb`; `certutil -L` muestra «mkcert development CA»). Tarjetas 6/6 cargadas desde finanzas.local; T1 $1,197.64 (1/1 cifras de la API); T3 $973.75, $598.82, $374.93 y $6,905.50 (4/4); T2, T4, T5 y T6 sin cifras distintas de 0 en la API (estados vacíos). Inicio → Detalle: API `reading_a.available_today_cents` = 1000000 y `reading_b` = 690550; en pantalla «$10,000.00». Hosts externos = 0.
- **C3-10** · Nombres de pantallas tal como U1 (es-MX), notas en español y fecha ISO en el encabezado (2026-10-12).
- **C3-11** · Versión «0.9» y fecha en el encabezado (sin «v1» todavía: borrador); Supuestos abajo; bloqueo `docs/estado/bloqueos/C3-matriz.md`.
- **S6-12** · `curl --cacert "$(mkcert -CAROOT)/rootCA.pem" https://finanzas.local:8443/` → `200 text/html`, `ssl_verify_result=0`, igual a `dist/index.html`, `<title>Finanzas Hogar</title>` y `rel="manifest"`. `/manifest.webmanifest` → 200, JSON válido: `name` «Finanzas Hogar», `start_url` `/inicio`, `display` `standalone`, `scope` `/`, `lang` `es-MX`, íconos 192, 512 y maskable 512 con 200 `image/png`; `/sw.js` 200. El manifiesto sale como `text/plain` (ya pedido por integrador-backend, solicitud 1).
- **U5-08 (dist)** · `npm run verificar:hosts` → «hosts externos que se pedirían = 0» (solo hay literales que no se piden: licencia de Inter y workbox).
- **Global** · `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → rc 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)» (29 s).

## Supuestos
- PROVISIONAL: la fila P-26 se marca como contenedor «Local» (bloqueo, punto 1).
- PROVISIONAL: «implementado» se mide **por necesidad**: basta con una de las operaciones de C1 de cada `N-xxx` de la fila, llamada por el cierre de importaciones de la pantalla. Exigir todas las operaciones daría falsos «no», porque C1 liga varias alternativas a la misma necesidad.
- PROVISIONAL: «probado» en Playwright se cuenta por el título de la prueba más 3 reglas escritas (bloqueo, punto 4). Las pruebas de esa columna usan la API simulada.
- Entradas de pantalla en la base: P-00 = `src/app/IniciarSesion.tsx`; P-03 y P-22 = `src/shell/HojaRegistro.tsx` (su comentario lo declara); Global = cierre de `src/main.tsx`.
- La contraseña de desarrollo se lee del archivo que documentaron semilla e integrador-backend; `FH_DEV_SEMILLA` permite otra ruta.

## Bloqueos
- Nuevo: `docs/estado/bloqueos/C3-matriz.md` (5 puntos).
- Citados: `WF3-frontend-servidor` (tipo del manifiesto y cabeceras de caché), `WF3-frontend-rutas` (rutas del navegador PROVISIONALES).

## Pendientes
1. Decidir los puntos 1 a 5 de `C3-matriz.md` y luego regenerar la matriz (`--escribir`).
2. WF-3b: capturas U5 (columna «capturado»), C3-08 con los 6 flujos contra el servidor, C3-09 completo (≥ 10 cifras, deuda viva $6,293.44, P1–P6) y la firma del revisor (C3-07).
3. La aprobación de frontend-core y frontend-areas (hallazgos abiertos en sus r3) sigue pendiente. Cuando cambien pantallas, hay que volver a correr el script: la comparación falla si la matriz quedó vieja.

## Propuestas para Axel
| Qué decidir | Resultado | Costo | Riesgo | Recomendación |
|---|---|---|---|---|
| 7 pantallas que usan otra consulta o mandan a otra pantalla (punto 2) | Si se acepta: la matriz queda en «sí» con solo ajustar el contrato | Bajo (documentos) | Ninguno en cifras | Aceptar las de P-03, P-12, P-17, P-18 y P-26.05; pedir de verdad la foto en la galería de P-02 y el emparejar productos en P-05 |
| Quitar adjunto, nutrición y edición del catálogo sin botón (punto 3) | Funciones que U1 prometió y no se ven | Medio (frontend) | Bajo | Pedirlas en la corrección del frontend |
| Semilla con presupuesto e ingreso de octubre (punto 5) | Inicio muestra cifras en las 6 tarjetas para las capturas | Bajo | Ninguno | Sí, antes de WF-3b |

## Paquetes instalados
Ninguno nuevo. `npm ci` reinstaló exactamente `package-lock.json`.

## Solicitudes fuera de alcance
1. `Makefile` (backend-base): objetivo `verify-wf3-c3`: `cd backend && go test ./internal/api -run TestContrato_Operaciones -count=1 -v > ../datos/tmp/c3/go.log`, después `cd frontend && PLAYWRIGHT_JSON_OUTPUT_NAME=../datos/tmp/c3/pw.json npx playwright test --reporter=json`, y al final `python3 -I frontend/scripts/c3-integracion.py --pw-json datos/tmp/c3/pw.json --go-log datos/tmp/c3/go.log`. Y `verify-wf3-humo`: `cd frontend && npm run test:integracion` (requiere el servidor de `ops/dev-run.sh`). Agregar los dos a `verify-wf3-servidor` o `verify-wf3-final`.
2. `backend/internal/api/http.go` (backend-base o auth-sistema): servir `*.webmanifest` como `application/manifest+json` (repite la solicitud 1 de integrador-backend).
3. `docs/estado/bloqueos/INDICE.md`: agregar el renglón «C3-matriz · integrador-frontend-c3 · 2026-10-12 · 5 preguntas de la matriz C3 v0.9».

## Consumo aproximado
~45 acciones; ~95 KB leídos (C3.md, convenciones, secciones de U1 con sed, extractos de reportes, salidas); ~65 KB escritos. Temporales `datos/tmp/c3/` y `datos/tmp/playwright-integracion/` borrados al terminar.
