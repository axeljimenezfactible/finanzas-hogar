# integrador-wf3a · cierre de WF-3a «Construcción»

Agente `integrador-wf3a` (Opus) · 12-oct-2026 · HEAD de partida «[Puerta 2 bis] K53». Revisiones recibidas: backend-core APROBADO (2 rondas) y backend-areas APROBADO (3 rondas); frontend-core NO APROBADO y frontend-areas NO APROBADO (3 rondas cada uno). El documento de integración es [`INTEGRACION-wf3a.md`](INTEGRACION-wf3a.md).

## Hecho

- **HECHO** · Verificación global: `make verify-wf3` y `make verify-all GOOSE=…` quedaron en rc 0. También pasaron `verify-wf3-secretos`, `verify-wf3-respaldo` y `verify-wf3-servidor`. La primera corrida falló por una prueba intermitente de Playwright; la detallo en Verificación, y no la corregí porque está fuera de mi alcance.
- **HECHO** · `INTEGRACION-wf3a.md` reúne:
  - lo hecho por agente y paquete;
  - los 16 criterios S6 y los 14 S7 con su estado;
  - la cobertura por paquete y la tabla operación → prueba por extracto (216 de 216);
  - las pantallas implementadas, probadas y capturadas;
  - los hallazgos abiertos de los dos frontends no aprobados;
  - los 40 bloqueos nuevos;
  - las propuestas PROVISIONALES de K31;
  - los pendientes de WF-3b y WF-4, los paquetes, las descargas y el consumo de 42 corridas.
- **HECHO** · `README.md`: estado «12-oct-2026, WF-3a», árbol (`frontend/`, `ops/`, `cmd/finanzas/`), sección «Qué entregó WF-3a», `verify-wf3` y `verify-wf3-final` en «Cómo correr la verificación», y el punto 6 de «Cómo sigue».
- **HECHO** · `docs/estado/bloqueos/INDICE.md`: nota de estado al principio y la sección «Bloqueos nuevos de WF-3a» con los 37 archivos que faltaban (los 3 de la semilla ya estaban). Esto cubre las solicitudes de INDICE de `integrador-backend` y de `integrador-frontend-c3`.
- **HECHO** · `.gitignore`: agregué de forma explícita `datos/*.db`, `datos/*.db-*`, `frontend/node_modules/`, `frontend/dist/`, `playwright-report/` y `test-results/`. `*.db`, `datos/certs/`, `datos/dev/` y `frontend/.gitignore` ya los cubrían.
- **HECHO** · Un solo commit con `git add -A -- . ':!docs/design/sistema' ':!docs/design/maquetas' ':!docs/estado/wfp3'`. Antes revisé lo que entra (ver Verificación).

## Archivos

| Archivo | Cambio |
|---|---|
| `docs/estado/wf3/INTEGRACION-wf3a.md` | nuevo |
| `docs/estado/wf3/integrador-wf3a.md` | nuevo (este reporte) |
| `README.md` | estado, árbol, «Qué entregó WF-3a», verificación y «Cómo sigue» |
| `docs/estado/bloqueos/INDICE.md` | nota de estado y 37 filas nuevas |
| `.gitignore` | bloque WF-3a |

## Verificación

- `make -k verify-wf3 verify-wf3-secretos verify-wf3-respaldo verify-wf3-servidor` (primera corrida): rc 2. El único fallo fue `pruebas/cola.spec.ts:85` [movil]: pasó el tiempo de 45 s y la salida dice «element was detached from the DOM». Quedaron 148 passed, 1 skipped y 1 failed. Los otros tres objetivos pasaron:
  - **S6-13**: 0 coincidencias.
  - **S6-09**: 56 tablas, 2,482 filas, 0 diferencias e `integrity_check` = ok.
  - **S6-11**: curl 200 y SAN con `finanzas.local`. En Android, NO VERIFICADO.
  - **S6-12**: el manifiesto y los íconos responden 200 y hay 0 cargas de terceros.
  - **S6-14**: escucha solo en 127.0.0.1 y 192.168.50.56; Tailscale no responde.
- `npx playwright test pruebas/cola.spec.ts --repeat-each=8`: **32 passed**. La falla es intermitente.
- `make verify-wf3` (segunda corrida): **rc 0**, «WF-3: verificación completa OK», 149 passed, 1 skipped y 0 failed, en 4 min 13 s.
  - **S6-01**: limpio; 32 paquetes OK.
  - **S6-02**: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8.
  - **S6-03**: 23 archivos idénticos; 216 operaciones en 172 rutas con 0 diferencias.
  - **C3-03** (script del frontend): 202 rutas, 0 fuera de `openapi.yaml` y 0 sin registrar. `c3-integracion.py` cuenta 204, porque también ve `armarUrl` y `fetch`.
  - **U5-08**: 0 hosts externos.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: **rc 0**, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)».
- `go test ./internal/api -run TestContrato_Operaciones -count=1 -v`: rc 0, «S6-03: 216 de 216 operaciones cumplen», 216 PASS. Contadas por extracto, todas están completas (tabla en INTEGRACION §5).
- **S7**: 13 de 14 cumplen según `s7-avisos`, `s7-pdf` y `REVISION-backend-areas-r3`. **S7-07 está PENDIENTE del orquestador.**
- **S6-11 en el móvil**: PENDIENTE, lo hace Axel.
- Antes del commit:
  - `git check-ignore` confirma que quedan fuera `datos/finanzas-dev.db`, `datos/certs/finanzas.local.pem`, `datos/dev`, `frontend/node_modules`, `frontend/dist` y `frontend/src/api/esquema.gen.ts`.
  - Un script propio que solo imprime conteos revisó los 429 archivos nuevos o cambiados: 0 llaves privadas, 0 tokens de bot, 0 llaves de API, 0 hashes argon2 completos y 0 apariciones de las contraseñas de desarrollo.
  - En la lista no hay `.db`, `.pem`, `.key`, `node_modules`, `playwright-report` ni `test-results`.

## Supuestos

- La prueba intermitente no impide cerrar WF-3a: con la segunda corrida verde, la tarea pedía reportarla y no arreglarla. La dejo como hallazgo I-01 (bajo) para `frontend-base`.
- Las cifras de consumo de cada agente salen de su propio reporte y son aproximadas. Para las corridas que dieron un rango tomé el valor central.
- No agregué a `Makefile` los objetivos `verify-wf3-c3` ni `verify-wf3-humo` que pidió `integrador-frontend-c3`: no arreglan una falla, así que quedan como pendiente de WF-3b.

## Bloqueos

No abrí bloqueos nuevos. Di de alta en el índice los 37 que faltaban (40 nuevos de WF-3a, todos ABIERTOS).

## Pendientes

1. Decidir si `frontend-base` corrige H-03, H-10, la base de `Ctrl+/`, N-04, N-05 e I-01 antes de WF-3b (INTEGRACION §7).
2. S7-07 real (orquestador) y S6-11 en Android (Axel).
3. WF-3b: capturas U5, flujos F-1 a F-6 contra el servidor, C3-09 completo y firma de C3. Antes hay que decidir los 5 puntos de `C3-matriz`.
4. Solicitudes fuera de alcance (INTEGRACION §11): el tipo del manifiesto y `Cache-Control` en `http.go`, las notas en dos bloqueos, el Makefile y `ajv` en `package.json`.
5. WF-4: systemd, ACL del token, IP y nombre permanentes, CA en el sistema y en Android, respaldo programado y `lualatex` con `luaotfload`.

## Propuestas para Axel

Están en la tabla de INTEGRACION §14 (resultado, costo, riesgo y recomendación). En corto:

- corregir el frontend antes de capturar;
- aprobar los valores de K31 de §9;
- reservar la IP 192.168.50.56;
- sembrar octubre;
- probar la CA en Android.

## Paquetes instalados

Ninguno. Los de WF-3a (mkcert por pacman, 3 módulos de Go y los de npm) están listados en INTEGRACION §12.

## Consumo aproximado

Unas 50 acciones de herramienta, unos 160 KB leídos (reportes de integradores y revisiones r3 por secciones, cabeceras de los bloqueos y salidas de make) y unos 60 KB escritos. Los temporales en `datos/tmp/integrador-wf3a/` y `datos/tmp/playwright-frontend/` se borran al terminar.
