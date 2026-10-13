# frontend-cierre · WF-3a cierre (12-oct-2026)

Agente: `frontend-cierre` (Opus). Base: HEAD `0197dad` + árbol de trabajo con los cambios de `backend-cierre` (C1 1.5.1 regenerado, deuda viva «hoy» en el servidor). Sin `git add` ni `commit`. **Cerrado por tope de acciones** (≈ 65 de 60): todo lo de mi alcance quedó hecho y verificado salvo lo que dice §Pendientes.

## Hecho
- **P-26.02 «Parámetros del tablero»** (decisión 4 = opción A de `C3-ajustes-tablero`, PROVISIONAL hasta la puerta 4): bloque nuevo en `S02Presupuesto.tsx` con las 5 claves de K54 (`debt.capacity_warning_pct`, `debt.capacity_limit_pct`, `dashboard.payment_reminder_count`, `budget.signal_months`, `budget.signal_margin_pct`). Lee `getDashboardSettings` y escribe `updateDashboardSettings` (`parameters` y `reset_parameters`), con el componente común `Parametro` (U4 §0.5: «Valor inicial: …», unidad y «Restablecer al valor inicial» solo si `is_initial = false`). Vista previa «Antes → Después» (U4 §0.2) en las 5 filas: tramos del barómetro (regla de U3 §B-1), escala (múltiplo de 10 ≥ 1.5 × tope, U3 §B-1), filas de «Pagos próximos», meses y margen de la sugerencia (RN-31). Validación V-23 con los rangos PROVISIONALES de C1 1.5 (fin del verde de 1 a tope − 1; tope de fin del verde + 1 a 100; 1–20; 2–12; 1–50). Si la respuesta no trae las 5 claves, solo el bloque dice que no cargó (la página sigue viva). Textos en `ajustes-a/textos.ts`.
- **Barómetro** (`Barometro.tsx`): `fill="none"` en el `<svg>` raíz (decisión 5); eje de la aguja como punto HTML de `size.viz.marker` con anillo de 2 px `surface-card` (N-4): mide **10 × 10 px** a 412 y 1440; la pista lleva `data-decorativo` (DESIGN §10.5 C-172).
- **`contraste-dom.js`** (decisión 5): un `<svg>` raíz con trazo propio se mide igual que antes; uno **sin trazo propio** se mide por lo que pintan sus hijos (trazo y relleno). No mide lo marcado `data-decorativo` ni la pintura idéntica al fondo (separaciones de 2 px, DESIGN §8.1); ambos se cuentan en `excluidos` (`svg_decorativos`, `svg_separaciones_color_fondo`).
- **Insignia** (`inicio/comun.tsx:25`): `w-max` → `w-fit max-w-full min-h-6` con texto que se parte; P-01 Detalle en 412 ya no desborda.
- **G-6** (`Tablero.tsx`): escala = `scale_max_pct` de `DashboardMonthlyOverview` (decisión 3, M-R3-1 opción B; C1 1.5.1 ya lo trae requerido). Tipos (`npm run tipos`) y mocks (`mocks:gen`, «221 operaciones de C1 1.5.1») regenerados.
- **K-3 deuda viva de hoy**: nada que cambiar en la interfaz; con el servidor de `backend-cierre`, K-3, P-17 total y P-17 línea dan **$3,895.00** = oráculo a 2026-10-12 (decisión 1).
- **Pruebas de flujos a U3 2.0** (decisión 6): F-1 parte y comprueba en `tarjeta-K-1`; F-4 (`flujos.spec.ts` y `pendientes.spec.ts` F-4b) toca la fila del 30-oct en `tarjeta-K-4`; F-3 ruta A según U1 §4 aclarado: con la lista de F-2 comprueba la condición (c) (vacío «Aún no hay precios para comparar» y «Registrar precios» → `/compras/productos`) y la ruta A completa queda en F-3A3 de `pendientes.spec.ts`; la ruta B sin cambios.
- **C3-09** (`c3-09.spec.ts`) reescrito a U3 2.0: K-4 (pago, hogar, Karen), K-3 (deuda viva de hoy, pagos del mes de todas las líneas, hogar, Karen, disponible B), B-1 (pagos del hogar = oráculo; % y margen con la fórmula de U3 §B-1 en ROUND_HALF_UP) y P-17 (deuda viva total y de la línea): **13 cifras del tablero y P-17**; exige ≥ 10. Lee la fila ya asentada (espera «Karen» y «Disponible hoy»).

## Archivos
- `frontend/src/pantallas/ajustes-a/S02Presupuesto.tsx`, `frontend/src/pantallas/ajustes-a/textos.ts`
- `frontend/src/componentes/Barometro.tsx`, `frontend/src/pantallas/inicio/comun.tsx`, `frontend/src/pantallas/inicio/Tablero.tsx` (solo G-6)
- `frontend/pruebas-u5/contraste-dom.js`
- `frontend/pruebas-u5/flujos/flujos.spec.ts`, `pendientes.spec.ts`, `c3-09.spec.ts` (+ capturas de flujos que reescriben las corridas: `capturas/F1-inicio-con-el-gasto.jpg`, `F3A-tarjeta-sin-precios.jpg` nueva, `F4b-resumen-despues.jpg`, `F5b-*.jpg`)
- `frontend/pruebas/ajustes-a/parametros-tablero.escritorio.spec.ts` (nueva), `frontend/pruebas/ajustes-a/ajustes.escritorio.spec.ts` (dato falso de `/api/v1/settings/dashboard`)
- Generados (no se editan): `frontend/src/api/esquema.gen.ts`, `frontend/src/mocks/*.gen.json`

## Verificación
- `npx tsc --noEmit`: sin errores. `npm run build`: OK (`✓ built`, PWA 80 entradas).
- **S6/U4 · P-26.02**: `npx playwright test pruebas/ajustes-a/` → **77 passed, 1 skipped** (movil y escritorio; incluye las 4 nuevas: 5 claves con valor inicial, vista previa de tramos y escala, Esc sin escritura, Enter = un PATCH `{parameters:{'debt.capacity_limit_pct':45}}`, V-23 «del 1 al 39», «Restablecer» = `{reset_parameters:['budget.signal_months']}`; tráfico validado contra C1 con `erroresC1` = []).
- **Inicio**: `pruebas/inicio/` en verde en la corrida conjunta (las 9 fallas de esa corrida eran de ajustes-a por el dato falso que faltaba; corregido arriba).
- **C3-03**: `npm run verificar:rutas` → «213 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 221)»; `c3-integracion.py` → «Operaciones con pantalla que el frontend no llama: 1 ['listPeriods']; sin excepción justificada o sin cobertura: 0» (`updateDashboardSettings` ya se llama). Ver Pendientes: el script sale con 1 por otra causa no revisada.
- **U5-08**: `npm run verificar:hosts` → salida 0.
- **C3-08** contra el servidor real, cada archivo en copia nueva (`servidor-prueba.sh arrancar`, 8444): `flujos.spec.ts` **7 de 7** (F-1, F-3A0, F-2, F-3, F-4, F-5, F-6); `pendientes.spec.ts` **4 de 4** (F-5, F-4b, F-2 copiar, F-3A3); `f7-f8.spec.ts` **2 de 2**.
- **C3-09** contra `https://finanzas.local:8443` (solo lectura), 2 corridas: «120 cifras leídas en pantalla · diferencias != 0: 0». Tablero: K-4 973.75/598.82/374.93; K-3 3,895.00, 973.75, 598.82, 374.93, 6,905.50; B-1 598.82, 2 %, 9,201.18; P-17 3,895.00 y 3,895.00.
- **U5-04 / U5-06** (prueba temporal ya borrada, P-01 Resumen y Detalle, 412/1440 × claro/oscuro): bajo umbral **0** en las 4; hijos del barómetro medidos = DESIGN C-159…C-171 (claro 3.466, 4.738, 5.127, 5.397, 10.76; oscuro 5.818, 6.036, 6.09, 6.425, 10.76); `svg_separaciones_color_fondo` 2 y `svg_decorativos` 1; eje 10×10; Detalle `scrollWidth` 412 de 412 y 1440 de 1440.
- Grep de las líneas agregadas en `frontend/src`: 0 hex/rgb/hsl; 0 inglés visible en los textos nuevos.

## Supuestos
- Lugar del bloque: al final de P-26.02 (U4 P-26.02 aún no lo ubica; la decisión 4 solo dice «sección en Ajustes › Presupuesto»).
- Vista previa de la escala con la regla de U3 §B-1 (no la manda el servidor en Ajustes). Unidades: las que traiga `unit` de la API.
- La pista del barómetro y las separaciones del color de la superficie no cuentan para U5-04 (DESIGN §10.5 las declara decorativas / §8.1).

## Bloqueos
- Ninguno nuevo. Siguen PROVISIONALES: `C3-ajustes-tablero` (opción A aplicada), `C1-tablero` (rangos), `C3-flujos` (condición (c) de la ruta A).

## Pendientes
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` y `verify-wf3-*`: **no corridos** (tope de acciones).
- `python3 -I frontend/scripts/c3-integracion.py` sale con **1** aunque la línea de operaciones da 0 sin cubrir; causa no revisada (probable matriz C3 vieja frente al C1 1.5.1 regenerado). Para el agente de C3.
- Capturas U5 de P-01 y P-26.02 (`docs/design/capturas-u5/`) no regeneradas: fuera de mi alcance; con este cambio `verify-u5` debería quedar en 0 bajo umbral al regenerarlas.

## Propuestas para Axel
- Ratificar en la puerta 4 los rangos de las 5 claves y el lugar del bloque en P-26.02.

## Solicitudes fuera de alcance
- `frontend/src/pantallas/inicio/Tablero.tsx` K-4 (fila de «Pagos próximos», línea `hogarYKaren`): mientras `listPersons` no responde, la fila dice un instante «Hogar $973.75» sin Karen (todo al hogar). Cambio: pintar la línea de partes solo con `c.personas.datos !== undefined` (o esqueleto).
- `frontend/src/componentes/Estados.tsx:15` (N-3): `role="status"` en el `div` de `Cargando`.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- ≈ 65 acciones, ≈ 230 KB leídos, ≈ 30 KB escritos.

## Continuación 1 (12-oct-2026, frontend-cierre · Opus)

Hice solo los Pendientes de arriba; no rehíce nada de lo Hecho ni toqué código: esta continuación no editó ningún archivo del producto, solo este reporte.

### Hecho
- **HECHO** · `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: salida **0** («== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)»; goose v3.28.0 up, down-to 0 y up sin errores; K54 OK con las 5 claves y 11 negativas rechazadas).
- **HECHO** · Objetivos `verify-wf3-*`, uno por uno (salida de cada uno):
  - `verify-wf3-go` **0** (S6-01): build, vet, gofmt y `go mod tidy -diff` limpios; `go test ./...`: 33 paquetes OK, 5 sin pruebas.
  - `verify-wf3-cobertura` **0** (S6-02, K25): dominio 87.0 %, motor_fin 89.2 %, motor_cred 94.0 %, recetas 93.2 %, precios 88.8 %.
  - `verify-wf3-contrato` **0** (S6-03, S6-04): «los 22 archivos coinciden byte a byte con el generador»; `TestRutasIgualesAlContrato` y `TestContratoSinSesionYEsqueleto` PASS.
  - `verify-wf3-frontend` **0**: `npm run verify` = build (`✓ built`), C3-03 «213 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 221)», U5-08 «hosts externos que se pedirían = 0», manifiesto válido y la suite completa de Playwright: **267 passed, 4 skipped, 0 fallas** (5.9 min).
  - `verify-wf3-secretos` **0** (S6-13): 0 coincidencias.
  - `verify-wf3-respaldo` **0** (S6-09): copia de la base de desarrollo, `integrity_check` ok, 56 tablas, 2927 filas, 0 diferencias de conteo.
  - `verify-wf3-servidor` **0** (S6-11, S6-12, S6-14): curl con la CA de mkcert 200; SAN con `DNS:finanzas.local`; manifiesto e íconos 192/512 autoalojados; `sw.js` 200; 0 cargas de terceros en `dist`; escucha solo en 127.0.0.1:8443 y 192.168.50.56:8443; Tailscale no responde. Chrome Android: NO VERIFICADO (manual de Axel).
  - `verify-u5` **2 (falla)**: «U5-04 contraste: 164 JSON, 24466 pares evaluados, bajo umbral = 2». Los 2 pares son el mismo nodo, el `<svg>` raíz del barómetro (`tarjeta-B-1 > barometro > svg`, fg #000000 sobre #36363A, razón 1.746) en `contraste/P-01-android-oscuro.json` y `P-01-escritorio-oscuro.json`. Es **evidencia vieja**: `docs/design/capturas-u5/` no cambió desde HEAD (`git status` 0 entradas) y en HEAD ya trae esos 2 pares. La fuente ya está corregida (`fill="none"` en la raíz, `Barometro.tsx:62`, decisión 5) y el agente anterior midió 0 bajo umbral en P-01 en las 4 combinaciones. Queda en 0 cuando se regeneren las capturas (solicitud abajo).
- **HECHO** · Causa de «`c3-integracion.py` sale con 1»: el script se corrió **sin** `--pw-json` ni `--go-log`. Sin esas entradas «probado» da 0/42 y la matriz generada no puede ser igual a la del documento; su docstring lo pide con ambas. Con las entradas reales (`npx playwright test --reporter=json`: 271 ejecutadas, 267 pasan, 4 omitidas, 0 fallas; `go test ./internal/api -run 'TestContrato_Operaciones|TestContrato_Rutas' -v`: 221 subpruebas PASS) da: «implementado sí 42/42; probado sí 42/42; capturado sí 41/42», «sin pantalla: 5 (sin justificar: 0)», «no llama: 1 ['listPeriods']; sin excepción o sin cobertura: 0», C3-03 «215 rutas invocadas; fuera: 0; sin registrar: 0». Sigue saliendo con **1** por una sola causa: «Matriz del documento = matriz generada: NO». Lo comparé con difflib (sin `--escribir`) y son 3 bloques, todos esperados por esta ronda:
  1. P-12 y P-26.0x: los números de línea de `ajustes-a/ajustes.escritorio.spec.ts` se corren +2 (dato falso de `/api/v1/settings/dashboard`).
  2. P-26 «PW: 43 → 47» y P-26.02 «PW: 4 → 8»: las pruebas nuevas de `parametros-tablero.escritorio.spec.ts`.
  3. Desaparece de la tabla de «operaciones no llamadas» la fila `updateDashboardSettings | P-26.02 | … PANTALLA PENDIENTE`: ya se llama (H-3 de REVISION-c3-tablero cerrado).
  El script no necesita cambios; hay que regenerar el documento (solicitud abajo).

### Verificación (Continuación 1)
- S6-01, S6-02, S6-03/S6-04, S6-09, S6-11, S6-12, S6-13 y S6-14: los objetivos de arriba, cada uno con salida 0.
- C3-02, C3-03 y C3-04: `c3-integracion.py --pw-json … --go-log …` con los conteos de arriba; C3-07/matriz: difiere solo en los 3 bloques descritos.
- U5-01 «PNG correctos: 164 de 164»; U5-03 «violaciones serias o críticas = 0»; U5-08 «fuera de finanzas.local/localhost = 0»; U5-11 «sin marca o sin ruta = 0»; U5-04 «bajo umbral = 2» (evidencia de HEAD; ver arriba).
- Temporales: `datos/tmp/fc1`, `datos/tmp/verify-wf3` y `datos/tmp/playwright-frontend` borrados. `git status` igual que antes de la continuación (59 entradas); no se tocaron archivos versionados.

### Supuestos (Continuación 1)
- Corrí los objetivos de `verify-wf3` uno por uno, no el agregado `verify-wf3`, para ver la salida de cada uno; el agregado falla hoy solo por `verify-u5`.
- No regeneré `docs/design/capturas-u5/` ni `docs/contratos/C3-integracion.md`: no están en mi ALCANCE.

### Bloqueos (Continuación 1)
- Ninguno nuevo.

### Pendientes (Continuación 1)
- Ninguno de mi alcance. Lo que falta para dejar `verify-u5` y la matriz de C3 en verde es de otros agentes (solicitudes de abajo).

### Solicitudes fuera de alcance (Continuación 1; se suman a las 2 de arriba)
- `docs/design/capturas-u5/` (agente de U5 o integrador): con el servidor de desarrollo ya reiniciado y el `dist` de esta ronda, correr `make verify-u5-servidor` (≈ 9 min). Regenera las 164 capturas y sus JSON con `fill="none"` y la regla nueva de `contraste-dom.js`, y deja `verify-u5` en «bajo umbral = 0».
- `docs/contratos/C3-integracion.md` (agente de C3): `npx playwright test --reporter=json` → R.json; `go test ./internal/api -run 'TestContrato_Operaciones|TestContrato_Rutas' -v` → G.log; luego `python3 -I frontend/scripts/c3-integracion.py --pw-json R.json --go-log G.log --escribir`. Cambia solo los 3 bloques descritos. Conviene hacerlo **después** de regenerar las capturas de U5, porque la columna «capturado» lee esa carpeta.

### Consumo aproximado (Continuación 1)
- ≈ 24 acciones, ≈ 45 KB leídos, ≈ 7 KB escritos (este apartado).
