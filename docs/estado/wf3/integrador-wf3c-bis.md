# integrador-wf3c-bis · WF-3c bis «Cierre del tablero» · 12-oct-2026

Agente: `integrador-wf3c-bis` (Opus). Base: HEAD `0197dad` («[WF-3c]») más el árbol de trabajo de WF-3c bis.

Alcance:
- `docs/estado/wf3/INTEGRACION-wf3c.md`: la sección «WF-3c bis», que es la §15, y una línea de aviso arriba;
- `README.md`;
- `docs/estado/bloqueos/INDICE.md`;
- este reporte;
- el commit.

No toqué código, pruebas, contratos, criterios ni las carpetas de la Pista 3.

## Hecho

- **HECHO · Verificación global.**
  - `make verify-wf3` dio **EXIT 0**, incluido `verify-u5` con 0 pares bajo umbral.
  - `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` dio **EXIT 0**.
  - Además: `verify-wf3-secretos` y `verify-wf3-servidor` dieron EXIT 0.
- **HECHO · Servidor de desarrollo con el `dist` nuevo.**
  - Corrí `bash ops/dev-stop.sh --conservar-nombre` y `bash ops/dev-run.sh`, que recompila `backend/bin/finanzas`.
  - El PID nuevo arrancó a las 20:20:47 y el binario es de las 20:20:42; no hay ningún `.go` más nuevo. El `dist` es de las 20:13 y salió de la fuente actual.
  - `curl --http2` a `https://finanzas.local:8443/` da `http_version=2` con 200.
  - El `index.html` servido es igual, byte por byte, a `frontend/dist/index.html`.
- **HECHO · `INTEGRACION-wf3c.md` §15 «WF-3c bis».** Trae:
  - qué cerró, contra la lista de §10;
  - C3 v1.4 **firmada**, con sus 3 bajos;
  - lo que queda sin aprobar: A-1 (alto, backend) y M-1 (medio, frontend), con el cambio exacto y una tabla de opciones;
  - lo consolidado para la puerta 4:
    - valores K31 y K54 PROVISIONALES;
    - bloqueos abiertos por fase: 154 archivos, 64 resueltos y 90 abiertos;
    - las preguntas `U3-tablero-*`, `WF3a-bis-semilla-dev`, `C3-flujos` y `S5-agente-externo`;
    - +$138.50 contra +$264.50;
    - la tabla para Axel;
  - el incidente de `c3-cierre`, ya corregido;
  - el consumo por agente y el commit.
- **HECHO · `INDICE.md`.**
  - Una línea de estado al cierre de WF-3c bis.
  - Las filas de `C1-unidad-parametros-tablero` y `U3-tablero-k4-personas`, que no tenían fila. Era la solicitud de `correccion-frontend-cierre-r1` y la nota del revisor de frontend r2.
  - El conteo de WF-3c bis: 3 nuevos y 2 actualizados.
  - Una nota de que `WF3c-ronda3-escalamiento` y `C3-ajustes-tablero` ya se atendieron, pero siguen abiertos hasta que Axel ratifique.
- **HECHO · `README.md`.**
  - Párrafo «Estado al 12-oct-2026, WF-3c bis».
  - Versiones: C1 v1.5.1 y C3 v1.4 firmada. Con esto queda atendida la solicitud 4 de `backend-cierre`, que pedía cambiar «C1 v1.5.0» por «C1 v1.5.1».
  - Sección «Qué entregó WF-3c bis».
  - Punto 9 de «Cómo sigue».
- **HECHO · Commit**, con la ruta y el mensaje del orquestador (ver Verificación).

## Archivos

- `docs/estado/wf3/INTEGRACION-wf3c.md`: §15 nueva y una línea de aviso bajo el título.
- `README.md`
- `docs/estado/bloqueos/INDICE.md`
- `docs/estado/wf3/integrador-wf3c-bis.md`: este reporte.
- Temporales en `datos/tmp/int-wf3cbis/` (registros de `make` y el script de conteo de bloqueos). Ya están borrados.

## Verificación

| Línea | Comando | Salida resumida |
|---|---|---|
| S6-01 | `make verify-wf3` › `verify-wf3-go` | go1.27.2. `go build`, `go vet`, `gofmt -l` y `go mod tidy -diff` sin cambios. `go test ./...`: 33 paquetes OK y 5 sin pruebas. |
| S6-02 | `verify-wf3-cobertura` | K25 OK: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8. Sin umbral: api 82.3. |
| S6-03 y S6-04 | `verify-wf3-contrato` | «los 22 archivos coinciden byte a byte con el generador»; «S6-03: 221 operaciones en 176 rutas, 0 diferencias». |
| Interfaz y C3-03 | `verify-wf3-frontend` | «mocks: 221 operaciones de C1 1.5.1». Playwright: 271 passed y 4 skipped, 0 fallas. |
| U5-01, U5-03, U5-04, U5-08 y U5-11 | `verify-u5` | 164 de 164 PNG. axe: 0 violaciones serias. **25818 pares y 0 bajo umbral.** Hosts: solo `finanzas.local:8443`. 41 filas C3. |
| Total | `make verify-wf3` | rc=0: «== WF-3: verificación completa OK (Go, cobertura K25, contrato C1, interfaz y evidencia U5)», en 7 min 10 s. |
| Global | `make verify-all GOOSE=…/goose` | rc=0: «== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». |
| S6-13 | `make verify-wf3-secretos` | «S6-13 grep: 0 coincidencias», rc=0. |
| S6-11, S6-12 y S6-14 | `make verify-wf3-servidor` | rc=0. curl 200 y SAN con `DNS:finanzas.local`. Manifest e íconos autoalojados y `sw.js` 200. 0 cargas de terceros. Escucha solo en 127.0.0.1 y 192.168.50.56. Tailscale no responde. Chrome Android: NO VERIFICADO (manual). |
| HTTP/2 | `curl --http2 --cacert <CA de mkcert> -w %{http_version}` | `http_version=2 code=200`. `cmp` del `index.html` servido contra `frontend/dist/index.html`: iguales. |
| C3-07 | el comando del bloque de firma | `e341d182e08567d797ae503c3d3970b8d3542d650a632806559fc2df5507d163`, igual que el bloque. |
| Conteo de bloqueos | script propio sobre `docs/estado/bloqueos/*.md` | «TOTAL 154 archivos · 64 resueltos · 90 abiertos». Después de mis filas, 0 abiertos sin fila en `INDICE.md`. |
| Higiene | `grep -ciE 'password\|contraseña\|token'` en `datos/dev/run/finanzas.log` | 0. |
| Commit | `git add -A --dry-run` con las exclusiones del orquestador | 171 rutas. 0 `.db`, `.pem`, `.key`, `node_modules`, `dist` o `datos/tmp`. 0 en las carpetas de la Pista 3. |

## Supuestos

- Para el conteo de bloqueos usé la regla de la puerta 2: la última línea «Estado del bloqueo» manda; si no la hay, cuentan «Respuesta de Axel» o «Resolución». La fase de origen es la sección de `INDICE.md` donde está la fila. Da 21 abiertos antes de WF-3, igual que la tabla de la puerta 2 del índice.
- La decisión 3 la describo por su contenido: `scale_max_pct` en C1 1.5.1. En la tabla de `WF3c-ronda3-escalamiento` es la opción (A); los reportes la llaman «opción B». No es una decisión mía; solo anoto la diferencia de nombre.
- El aviso de una línea bajo el título de `INTEGRACION-wf3c.md` lo tomé como parte de la sección «WF-3c bis», porque solo remite a ella.
- No corrí `verify-wf3-respaldo`: esta ronda no tocó respaldos, y la continuación de `frontend-cierre` lo dio en EXIT 0.

## Bloqueos

Ninguno nuevo. Di de alta en el índice las 2 fichas nuevas del frontend.

## Pendientes

- **El orquestador decide** la corrección acotada de A-1 (backend) y M-1 (frontend), que es la recomendada en §15.5. Ningún área sin aprobar tiene ronda pendiente dentro de este flujo.
- Las solicitudes de abajo, que no entran en mi alcance.

## Propuestas para Axel

- **Antes de la puerta 4:** corregir A-1 y M-1, que son unas 15 líneas de Go y un ajuste de D-3, cada uno con su prueba. Si no se corrige, el comparador puede recomendar una tarjeta archivada para una compra nueva.
- **En la puerta 4:**
  - confirmar `S6-deuda-viva-hoy` (a), `S6-tablero-lineas-archivadas` (a) y su punto 5, `C3-ajustes-tablero` (A), `C1-unidad-parametros-tablero` (A), `U3-tablero-k4-personas` (A), `C3-flujos` punto 2 (c) y los valores de K54 tal cual;
  - resolver Q2-49 contra la semilla con precios de Walmart.

  Todo está en la tabla de §15.6.10.

## Solicitudes fuera de alcance

1. `backend/internal/api/creditos/condiciones.go:676` (A-1): pasar al comparador solo las líneas activas (por ejemplo, `s.activeStates()`) sin tocar `states`. Agregar `TestOp_compareScenarios_LineaArchivada`. Corregir `docs/estado/bloqueos/S6-tablero-lineas-archivadas.md:55`.
2. `frontend/src/pantallas/inicio/Detalle.tsx:92-98` y `Calendario.tsx:33-42` (M-1): no pintar filas hasta tener `listPersons`, y mostrar error con «Reintentar» si falla. Agregar una prueba con `/api/v1/persons` retenido y un `MutationObserver`.
3. `docs/criterios/C1.md`, líneas 6, 9, 106 y 107: «1.5.0» → «1.5.1» (en la 107, `version: "1.5.1"`), con el nombre «C1 v1.5.1 (escala de G-6, WF-3c bis)». Es la solicitud 1 de `correccion-backend-cierre-r1`.
4. `docs/design/U3-dashboard.md` G-6 (línea 229) y §2.6: citar `DashboardMonthlyOverview.scale_max_pct` (C1 v1.5.1).
5. `docs/design/U1-arquitectura.md:652` y `:978`, y `docs/design/U3-dashboard.md:317`: cambiar Walmart $3,233.00 (+$138.50) por $3,359.00 (+$264.50), la cifra de K46 (Q2-47 = A). En idénticos se queda +$100.50.
6. `backend/internal/api/gastos-inicio/vistas.go:235` (`deriveCredits`): cambiar `mc.LiveDebt(...)` por `creditos.LiveDebtAt(...)`.
7. `backend/internal/api/creditos/`: confirmar la fecha por omisión de `GET /api/v1/debt/live`, que sin `date` respondió 2027-09-30 con `total_cents` 0 (observación de `u5-cierre`).
8. `frontend/pruebas-integracion/humo-inicio.spec.ts`: leer K-1 a K-4 y B-1 del tablero de U3 2.0.
9. `frontend/pruebas-u5/flujos/servidor-prueba.sh` o `playwright.config.ts`: que las pruebas que escriben no puedan correr contra 8443 sin `FH_BASE` (incidente de §15.9).
10. `docs/contratos/C3-integracion.md` línea 3: dice «sin firma». Va en la próxima versión y necesita firma nueva.
11. `Makefile` (líneas 8, 66 y 149): «esquema 1.4» → «esquema 1.5» (B-N4 de WF-3c).
12. `frontend/src/componentes/Estados.tsx:15`: `role="status"` (N-3).
13. Frontend B-1 a B-5 de `REVISION-frontend-cierre-r2.md`:
    - `flujos.spec.ts:71`: agregar `expect(reg.ok)`;
    - U4 P-26.02 y `C3-ajustes-tablero.md:21`;
    - `tableroNoCargo`;
    - el comentario de `c3-09.spec.ts:55-56`.

## Paquetes instalados

Ninguno.

## Consumo aproximado

Unas 45 acciones, unos 200 KB leídos y unos 50 KB escritos.
