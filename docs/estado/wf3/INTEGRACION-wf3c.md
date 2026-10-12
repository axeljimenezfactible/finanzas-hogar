# INTEGRACIÓN de WF-3c «Inicio como tablero de control (K54)»

12-oct-2026 · Integrador: `integrador-wf3c` (Opus) · Base: HEAD `5a089da` («[WF-3b]») más el árbol de trabajo de WF-3c.

Estado del orquestador al integrar:
- Especificación: no aprobada, 3 rondas.
- Backend: no aprobado, 3 rondas.
- Frontend: **aprobado**, 3 rondas.
- C3: no aprobada.

Commit: «[WF-3c] Inicio como tablero de control (K54): …». Su hash va en el reporte `integrador-wf3c.md` y en la salida del orquestador. No va en este archivo, porque el archivo es parte del commit.

## 1. En corto

**HECHO**
- **K54 está aplicada.** Axel pidió que Inicio › Resumen deje las 6 tarjetas y sea un tablero de control. Ahora tiene:
  - el gasto real del mes como indicador principal (K-1);
  - ingreso (K-2), créditos (K-3) y pagos próximos (K-4);
  - el **barómetro de capacidad de endeudamiento** (B-1): pagos de crédito del mes ÷ ingreso mensual neto × 100. Verde hasta 30 %, amarillo hasta 40 % y rojo arriba de 40 %, con umbrales configurables;
  - 6 gráficas para anticiparse (G-1 a G-6);
  - «Más del hogar» al final.
- **Versiones:**
  - Documentos: U3 2.0, DESIGN 1.1 (§8.9 «Barómetro», sin tokens nuevos), criterios U3 1.4, U5 1.3, P3 1.4, C1 1.5 y S1 1.3.
  - Esquema y contrato: esquema 1.5 (migración `00005_tablero.sql`) y C1 v1.5.0, con 221 operaciones, 176 rutas y 267 esquemas.
  - API y pantalla: el servidor atiende las 5 operaciones nuevas y P-01 tiene su tablero.
- **HTTP/2:** el servidor de desarrollo corre con el binario y el `dist` actuales. Responde con `curl -w %{http_version}` = **2**, por loopback y por la LAN.
- **Capturas:** las 164 se regeneraron, las 4 de P-01 ya muestran el tablero, y hay 16 de las subsecciones de P-01.

**NO APROBADO** (el detalle y las salidas propuestas están en `docs/estado/bloqueos/WF3c-ronda3-escalamiento.md`):
- Especificación: queda 1 medio, M-R3-1.
- Backend: queda 1 medio, M-2.
- C3 v1.3: devuelta y **sin firma**, con 4 altos.

**Verificación global:**
- `make verify-all` da **EXIT 0**.
- `make verify-wf3` da **EXIT 2**. La única causa es U5-04: 2 pares bajo umbral, del `<svg>` raíz del barómetro en P-01 oscuro. Es un falso positivo de la medición, comprobado (§2).

## 2. Verificación global (la corrió este integrador el 12-oct-2026)

| Línea | Comando | Resultado |
|---|---|---|
| S6-01 | `make verify-wf3` › `verify-wf3-go` | go1.27.2. `go build`, `go vet`, `gofmt -l` y `go mod tidy -diff` sin cambios. `go test ./...`: 33 paquetes OK y 5 sin pruebas. |
| S6-02 | `verify-wf3-cobertura` | K25 OK: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8. Sin umbral: api 82.0, auth 83.6, jobs 93.0, reportes 77.3 y telegram 80.0. |
| S6-03 | `verify-wf3-contrato` | «los 22 archivos coinciden byte a byte con el generador»; «S6-03: 221 operaciones en 176 rutas, 0 diferencias». `TestRutasIgualesAlContrato` y `TestContratoSinSesionYEsqueleto` dan PASS. |
| Interfaz y C3-03 | `verify-wf3-frontend` (`npm run verify`) | `tsc`, `vite build` y 80 entradas precacheadas. «C3-03 frontend: 212 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 221)». Playwright con la API simulada: **259 passed y 4 skipped** en 5.8 min, con las 9 pruebas de `tablero.escritorio.spec.ts`. |
| U5-01, U5-03, U5-08 y U5-11 | `verify-u5` | 164 de 164 PNG. axe: 164 JSON y 0 violaciones serias. Hosts: solo `finanzas.local:8443`. C3: 41 filas «capturado» con sus 4 rutas. |
| **U5-04** | `verify-u5` | **FALLA: «24466 pares evaluados, bajo umbral = 2»**, con EXIT 2. Comprobado en `contraste/P-01-android-oscuro.json` y `P-01-escritorio-oscuro.json`, 1 par cada uno; las 162 restantes dan 0. Es el `<svg>` raíz `aria-hidden` del barómetro (`Barometro.tsx:61`), que no tiene `stroke` ni `fill`. `contraste-dom.js:91-96` toma su `fill` inicial negro (1.746 sobre #36363A). Los trazos reales pasan (de 4.74 a 10.76). La salida está en `WF3c-ronda3-escalamiento`, punto 4. |
| Global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | **EXIT 0**: «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». En S1: K54 OK, DOWN14 OK, GOOSE OK y RESUMEN TODO OK. U1: 63 OK y 0 pendientes. Oráculos: 15 OK, todo coincide. Tokens C2 idénticos. |
| S6-13 | `make verify-wf3-secretos` | 0 coincidencias. |
| S6-09 | `make verify-wf3-respaldo` | Copia de la base de desarrollo. `integrity_check` ok; 56 tablas, 2907 filas y 0 diferencias de conteo. |
| S6-11, S6-12 y S6-14 | `make verify-wf3-servidor` | mkcert v1.4.4. `finanzas.local` → 192.168.50.56. curl 200, con el SAN `DNS:finanzas.local`. Manifest con íconos 192 y 512 y `sw.js` 200; 0 cargas de terceros. Escucha solo en `127.0.0.1:8443` y `192.168.50.56:8443`; por Tailscale no responde. Chrome Android: NO VERIFICADO, es manual de Axel. |
| HTTP/2 | `curl --cacert <CA de mkcert> -w %{http_version}` | `https://finanzas.local:8443/` da **2** (200). Por la LAN (`--resolve …:192.168.50.56`) da **2**. Con `--http1.1` da 1.1. La prueba `TestEscucha_NegociaHTTP2` (`backend/cmd/finanzas/tls_h2_test.go`) pasa y falla si se quita `NextProtos`, según `api-tablero`. |
| `dist` servido | `index.html` del servidor contra `frontend/dist/assets` | El servidor sirve `assets/index-CBAvKt4u.js`, que es el `dist` que generó `verify-wf3-frontend` a las 17:25. El binario es de las 16:36 y su proceso arrancó a las 16:37; no hay ningún `.go` más nuevo que el binario. |

## 3. K54 · decisión aplicada

- `docs/plan/decisiones.md`, línea 158: K54, del 12-oct-2026. Trae las palabras literales de Axel y 4 decisiones: (1) la fórmula del barómetro, (2) los umbrales 30 y 40 configurables, (3) se retira el tope de «≤ 6 tarjetas», (4) el propósito preventivo. La registró `u3-v2`. El revisor de especificación comparó las citas con el texto del orquestador y son idénticas.
- Bajo abierto B-r1-1: no se registró el «segundo problema» (la lentitud) ni la cita de Axel sobre ese punto.

## 4. Especificación: U3 2.0, DESIGN 1.1 y criterios

- **U3 2.0** (`docs/design/U3-dashboard.md`): §2 nueva con K-1 a K-4, B-1 y G-1 a G-6.
  - §2.2: desglose en cadena hacia P-02 y P-06.
  - §2.3: «Qué miro para prevenir».
  - §2.4: orden en 412 y rejilla en 1440.
  - §2.5: «Más del hogar».
  - §2.6: operaciones de C1.
  - Las cifras de ejemplo salen de `oraculo/model.json`.
- **DESIGN 1.1:** §8.9 «Barómetro» (arco de tres tramos con los tokens de rol, aguja, dígito `title-1` tabular, insignia con icono y palabra, `role="meter"`) y dos filas nuevas en §8.2. §7.6 queda sin el tope de 6 tarjetas. **Ningún token nuevo:** `C2-tokens.json` y `tokens.css` no cambian, y `sincronia.mjs` sale idéntico. Hay 14 pares de contraste nuevos (C-159 a C-172), 28 valores y ninguno bajo su mínimo.
- **Criterios:** U3 1.4, U5 1.3 (U5-06 sin el tope de 6), P3 1.4, C1 1.5 y S1 1.3. Siguen siendo 261 criterios (201 alta y 60 media), sin cambiar ningún id ni severidad. `make verify-criterios` pasa dentro de `verify-all`.
- **U1 1.2:** se editó, pero no cambió de versión. §3.3 reutiliza P-01 y F-3 aclara la ruta A con la condición (c), PROVISIONAL del orquestador (`C3-flujos`, punto 2).
- **Dictamen:** la especificación tuvo 3 rondas y quedó **NO APROBADA**, con 0 altos y 1 medio (M-R3-1). Ver §9.

## 5. Esquema 1.5 y C1 v1.5

- **Migración `backend/db/migrations/00005_tablero.sql`** (esquema 1.4 → 1.5): siembra 5 claves de `app_setting` del hogar.
  - `debt.capacity_warning_pct` = 30 y `debt.capacity_limit_pct` = 40 son FIRMES (K54).
  - `dashboard.payment_reminder_count` = 5, `budget.signal_months` = 3 y `budget.signal_margin_pct` = 10 son PROVISIONALES.
  - Agrega los disparadores de rango `trg_app_setting_param_k54_ins` y `_upd`.
  - No crea tablas. El Down deja el esquema 1.4 idéntico.
  - Ninguna migración aplicada se editó: `DDL OK` compara de 00001 a 00004 con HEAD.
- **Diccionario y reglas:** `diccionario-datos.md` 1.5 se regeneró y tiene 30 parámetros. `reglas-negocio.md` 1.5 trae RN-29 (las 5 claves), RN-30 (capacidad de endeudamiento, entero con ROUND_HALF_UP) y RN-31 (sugerencia de presupuesto, PROVISIONAL).
- **C1 v1.5.0** (`openapi.yaml`, `info.version` 1.5.0, generado con `c1gen`):

| Conteo | v1.4 (WF-3b) | v1.5 (WF-3c) |
|---|---|---|
| Operaciones | 216 | **221** (+5) |
| Rutas | no se recontó | **176** |
| Esquemas | no se recontó | **267** |
| Extractos | no se recontó | **22**: se regeneraron, y `ajustes-f.yaml` desaparece al repartirse |

- **Operaciones nuevas:**
  - `getDashboardDebtCapacity` (GET `/api/v1/dashboard/debt-capacity?period`);
  - `getDashboardMonthlyOverview` (GET `/api/v1/dashboard/monthly-overview?from&to`, con un máximo de 24 meses);
  - `getBudgetSignals` (GET `/api/v1/budget/signals?period`);
  - la pareja `getDashboardSettings` / `updateDashboardSettings` (GET y PATCH `/api/v1/settings/dashboard`, con el patrón K44 y `reset_parameters`).
- **Campos nuevos:** 14 campos opcionales: `DashboardMonthVsPlan` 8, `IncomeMonth` 4 y `CreditTotals` 2.
- **Métricas:** M-18 a M-22 en C1-datos §9, y la fila M-01 con los 8 campos nuevos.
- `verify_c1 --db`: «comprobaciones OK: 33, TODO OK», según el revisor de especificación r3.

## 6. API (servidor Go)

- **Paquete nuevo `backend/internal/api/tablero/`:** barómetro, resumen mensual, señales de presupuesto y parámetros. Además:
  - `internal/api/creditos/tablero.go` tiene los pagos por mes;
  - `internal/dominio/tablero_datos.go` tiene `MonthHasData`, la regla «sin datos», PROVISIONAL según `S6-tablero-sin-datos`.
- **K-1:** `addMonthComparison` (`api/gastos-inicio/inicio.go`) llena los 8 campos opcionales de `DashboardMonthVsPlan` con el motor S2.
- **Pruebas contra el motor:**
  - `payments_cents` es igual a `getCreditTotals` (S3);
  - el devengado y el efectivo de M-21 son iguales a M-01 (S2);
  - RN-30 cuadra con una cuenta a mano usando ROUND_HALF_UP.
- **Cobertura:** `api/tablero` 87.8 % y `api/creditos` 80.4 %, según el revisor backend r3.
- **Cifras en el servidor de desarrollo** (revisor backend r3, con HTTP/2): son iguales al oráculo en 6 grupos y en la serie de 12 meses. Por ejemplo, octubre da 2 % y 4 %, con margen de $9,201.18, y noviembre y diciembre dan 5 % y 6 %.
- **Umbrales:** con 45 contra un tope de 40 la respuesta es 422, en español. `reset_parameters` vuelve a 30 y 40.
- **HTTP/2:** `cmd/finanzas/tls.go` anuncia `h2` por ALPN (`NextProtos: ["h2", "http/1.1"]`), y S6-14 no cambia.
- **Dictamen:** el backend tuvo 3 rondas y quedó **NO APROBADO**, con 0 altos y 1 medio (M-2). Ver §9.

## 7. Pantalla P-01 (frontend) y sus capturas

- **Archivos nuevos:**
  - `src/pantallas/inicio/Tablero.tsx` (K-1 a K-4, B-1, G-1 a G-6 y «Más del hogar»);
  - `src/componentes/Barometro.tsx` (DESIGN §8.9; sin dato usa `role="img"`, PROVISIONAL según `WF3c-barometro-sin-dato`);
  - `src/pantallas/inicio/mocks.ts`.
- **Barra lateral:** muestra los iconos de sección (M-7 de U5 ronda 3). 4 son PROVISIONALES (`WF3c-barra-iconos-seccion`).
- **Dictamen:** el frontend tuvo 3 rondas y quedó **APROBADO**, con 0 altos, 0 medios y 8 bajos.
- **Capturas de P-01** (regeneradas por `u5-tablero` entre las 16:37 y las 16:44, después de la corrección r2 del frontend; desde entonces no cambió ningún archivo de `frontend/src`, salvo generados con el mismo contenido):
  - `docs/design/capturas-u5/P-01-escritorio-claro.png` (1440 × 3029), `P-01-escritorio-oscuro.png`, `P-01-android-claro.png` y `P-01-android-oscuro.png` (412 × 5299);
  - axe da 0 serios en las 4;
  - contraste: 0 bajo umbral en claro y 1 par en cada oscuro (el falso positivo de §2).
- **Subsecciones:** `docs/design/capturas-u5/evidencia/inicio-subsecciones/` tiene 16 PNG y 16 JSON (Resumen, Detalle, Tendencias y Calendario × 4 combinaciones).
- **Cifras:** `evidencia/tablero-cifras/` guarda el texto visible y las respuestas crudas de `/api/v1/dashboard/*`.
- **Lo que se ve en `P-01-escritorio-claro.png`** (lo revisó el integrador):

| Bloque | Qué muestra |
|---|---|
| K-1 | «$1,197.64 gastado en octubre», «60 % de $2,000.00», «Pagado en efectivo $598.82» y «Sin datos de septiembre». |
| B-1 | Arco con marcas en 30 % y 40 %, «2 %», «de 40 % de tope», «Verde · dentro del rango», «Te caben $9,201.18», «Con la parte de Karen: 4 %» e «Ingreso esperado». |
| K-2 | «+$24,500.00». |
| K-3 | «**$18,999.22** deuda viva total». Es el hallazgo H-4. |
| K-4 | «Vie 30 oct 2026 · Kueski · $973.75 · Hogar $598.82 · Karen $374.93». |
| Gráficas | G-3, G-1, G-6, deuda viva por línea con $3,895.00, deuda viva proyectada, y devengado contra efectivo. |
| Más del hogar | 4 tarjetas. |

  Estas cifras son las mismas que `u5-tablero` comparó con `model.json`: próximo pago, barómetro, disponible B $6,905.50, devengado y efectivo acumulados y deuda viva al 10-abr-2027 de $6,293.44.
- **Observaciones visuales del integrador** (bajas, sin criterio que falle):
  - En G-1 la leyenda queda sobre el extremo derecho del área de trazo, junto a la barra de oct-2026.
  - En G-6, la etiqueta «Capacidad de endeudamiento» pisa «30 %». También lo vio `u5-tablero`.
  - La tarjeta K-1 se estira a la altura de B-1 y deja un hueco en blanco.

## 8. C3 v1.3 (criterios `docs/criterios/C3.md` 1.1): **NO firmada**

`c3-tablero` regeneró la matriz con C1 v1.5. Quedó así: implementado 42 de 42, probado 42 de 42 y capturado 41 de 42 (la fila Global sigue en «no», según `C3-capturas-global`). La versión es v1.3, sin firma.

`revisor-c3-tablero` (ronda 1) la **DEVOLVIÓ** con 4 altos, 1 medio y 4 bajos:

| Criterio | Dictamen | Nota |
|---|---|---|
| C3-01 y C3-02 | cumple | 41 pantallas y 221 operaciones. Hay 5 operaciones sin pantalla con excepción. |
| **C3-03** | **no cumple** | H-3: `updateDashboardSettings` sin pantalla (`C3-ajustes-tablero`). |
| C3-04 | cumple | Go: 1147 PASS. Playwright: 259 pasan, 0 fallan. |
| C3-05 y C3-06 | cumple | 164 PNG. «no» solo en la fila Global. |
| **C3-07** | **sin firma** | — |
| **C3-08** | **no cumple** | H-1: F-1, F-3 y F-4 fallan por los selectores `tarjeta-1` y `tarjeta-3`, que ya no existen. |
| **C3-09** | **no cumple** | H-2: `c3-09.spec.ts` falla. H-4: la deuda viva de K-3 y P-17 es $18,999.22 y el oráculo, `getLiveDebt` y P-07 dan $3,895.00. |
| C3-10 y C3-11 | cumple | — |

Además:
- 55 de 57 cifras de pantalla son iguales a `model.json`.
- El medio H-5: la matriz dice que F-3A3 falla y hoy pasa.

## 9. Rondas y dictámenes

| Área | Autor | Rondas | Dictamen final | Abiertos |
|---|---|---|---|---|
| Especificación (K54, U3 2.0, DESIGN 1.1, esquema 1.5 y C1 1.5) | `u3-v2` y `s1-c1-v15` | 3 (más 2 correcciones) | **NO APROBADO** | 1 medio (M-R3-1) y 9 bajos (B-R3-1, B-R3-2, B-N1 a B-N5, B-r1-1 y B-r1-3; B-r1-4 lo cerró este integrador) |
| Backend del tablero (API y HTTP/2) | `api-tablero` | 3 (más 2 correcciones) | **NO APROBADO** | 1 medio (M-2) y 5 bajos (B-1, B-2, B-4, B-5 y B-6; B-3 lo cerró este integrador) |
| Frontend del tablero (P-01 y barra lateral) | `front-tablero` | 3 (más 2 correcciones) | **APROBADO** | 8 bajos (N-2, N-3, N-4 y H-3 a H-7) |
| U5 del tablero (capturas y cifras) | `u5-tablero` | sin ronda de revisión | — | U5-04 con 2 pares, la página de 475 px y 4 observaciones (§10) |
| C3 v1.3 | `c3-tablero` | 1 | **DEVUELTO, sin firma** | 4 altos, 1 medio y 4 bajos |

## 10. Hallazgos abiertos (consolidados)

| Sev. | Id · origen | Qué | Archivo y cambio exacto | Dueño propuesto |
|---|---|---|---|---|
| alta | H-1 · C3 | F-1, F-3 y F-4 fallan contra el servidor real. | `frontend/pruebas-u5/flujos/flujos.spec.ts:18-19,41-42,226` y `pendientes.spec.ts:62`: `tarjeta-1` pasa a la tarjeta K-1 y `tarjeta-3` a `tarjeta-K-4`. A F-3A3 le falta preparar los precios de Vallejo (`c3-tablero.md`). | `u5-flujos` u otro que designe el orquestador |
| alta | H-2 · C3 | `c3-09.spec.ts` falla en la línea 31 (`tarjeta-3`). | `frontend/pruebas-u5/flujos/c3-09.spec.ts:31-37` | el mismo |
| alta | H-3 · C3 | Ninguna pantalla llama a `updateDashboardSettings`. | Apartado «Parámetros del tablero» en `S02Presupuesto.tsx`, según la respuesta a `C3-ajustes-tablero`. | frontend, después de la puerta 4 |
| alta | H-4 · C3 (también es la diferencia 1 de `u5-tablero`) | K-3 y P-17 muestran una deuda viva total de $18,999.22 y el oráculo a hoy da $3,895.00. | `getCreditTotals.live_debt_cents`: la misma regla de fecha que `getLiveDebt`, si Axel lo confirma. | backend (`api-creditos`) |
| alta | U5-06 · `u5-tablero` | P-01 Detalle en Android mide 475 px de ancho en lugar de 412. | `frontend/src/pantallas/inicio/comun.tsx:25`: cambiar `w-max` por `max-w-full`, con texto que envuelve y alto automático. | frontend |
| media | M-R3-1 · especificación | G-6 usa `scale_max_pct`, que no está en `DashboardMonthlyOverview`. | Opción A: agregar el campo en C1 (c1gen, `build.py`, `tablero.go`, U3 §2.6 y M-21). Opción B: cambiar el texto de U3 G-6. | especificación |
| media | M-2 · backend | Los pagos del tablero solo cuentan las líneas activas, contra RN-30 y C1. | Salida (i): opción (a) en `backend/internal/api/creditos/tablero.go` (`PaymentsByMonth`) y quitar el SKIP de `opcion_a_PROVISIONAL`. Salida (ii): texto de RN-30 y C1. | backend |
| media | U5-04 · `make verify-wf3` | 2 pares del `<svg>` raíz del barómetro. Es un falso positivo comprobado. | `frontend/src/componentes/Barometro.tsx:61`: `fill="none"` en el `<svg>` raíz, y luego regenerar las capturas de P-01. La otra opción es `frontend/pruebas-u5/contraste-dom.js:91-96`. | frontend o U5 |
| media | H-5 · C3 | La matriz dice que F-3A3 falla y hoy pasa en 3 toques. | `docs/contratos/C3-integracion.md` §4, al regenerarla. | C3 |
| media | Diferencia 2 · `u5-tablero` | La «deuda viva proyectada» empieza en $18,025.47 en oct-2026 y el oráculo da $2,921.25. | Tiene el mismo origen que H-4 (`S6-semilla-despensas-futuras`). | backend |
| baja | B-R3-2 · especificación (parcial) | Los 8 campos de `DashboardMonthVsPlan` ya se llenan. `grep` de `FortnightsDifferenceCents` y `MonthPaymentCount` fuera de `types_gen.go` da 0: los de `IncomeMonth` y `CreditTotals` no los llena ningún manejador. | `api/gastos-inicio` y `api/creditos` | backend |
| baja | B-N4 · especificación | El `Makefile` (líneas 8, 66 y 149) y `s1_verifica.py:190` todavía dicen «esquema 1.4». | Cambiar el texto a «esquema 1.5». | integrador o S1 |
| baja | otros | Especificación: B-R3-1, B-N1, B-N2, B-N3, B-N5, B-r1-1 y B-r1-3. Backend: B-1 (400 en la API y 422 en C1), B-2, B-4, B-5 y B-6. Frontend: N-2, N-3 (`Estados.tsx:15`, `role="status"`), N-4 y H-3 a H-7. C3: H-6 a H-9 (H-8: U1 §3.3 y §3.29 sin N-191 ni N-275 en P-01). `u5-tablero`: el subtítulo «Últimos 12 meses» debería decir «Próximos 12 meses»; a `playwright.u5.config.ts` le faltan `subsecciones.spec.ts` y `tablero-cifras.spec.ts` en `testMatch`. | Ver cada informe `REVISION-*`. | según el informe |

## 11. Bloqueos nuevos (14, todos ABIERTOS) y actualizados

Los registra `docs/estado/bloqueos/INDICE.md` en «Bloqueos nuevos de WF-3c». Los archivos son:
- `U3-tablero-barometro`, `U3-tablero-horizonte`, `U3-tablero-mas-del-hogar`, `U3-tablero-recordatorio` y `U3-tablero-senal-presupuesto`, de `u3-v2`;
- `C1-tablero`, de `s1-c1-v15`;
- `S6-tablero-sin-datos`, de `api-tablero`;
- `S6-tablero-senal-periodo` y `S6-tablero-lineas-archivadas`, de `correccion-tablero-backend-r1`;
- `C3-ajustes-tablero`, de `c3-tablero`;
- `WF3c-barra-iconos-seccion` y `WF3c-tablero-interfaz`, de `front-tablero`;
- `WF3c-barometro-sin-dato`, de `correccion-tablero-frontend-r2`;
- `WF3c-ronda3-escalamiento`, de este integrador.

Se actualizó `C3-flujos.md`. Ninguno toca un valor K31. Con el índice al día se cierran B-3 (backend) y B-r1-4 (especificación).

## 12. Para la puerta 4 (PROPUESTA, no decidido)

| Qué decide | Resultado | Costo | Riesgo | Recomendación |
|---|---|---|---|---|
| Ronda acotada «WF-3c bis» (`WF3c-ronda3-escalamiento`) | `verify-wf3` en verde y C3 firmable | 3 agentes cortos y una corrida de U5 de P-01 | Bajo: ningún motor ni cifra del oráculo cambia | **Sí**, antes de la puerta 4 |
| «Deuda viva total» de K-3 y P-17 (H-4) | Una sola cifra de deuda viva: $3,895.00 hoy, igual que el oráculo | Una regla de fecha en `getCreditTotals` | Bajo | Solo las compras hechas a la fecha |
| Dónde se editan los 5 parámetros del tablero (`C3-ajustes-tablero` y `C1-tablero` punto 1) | Los umbrales 30 y 40 se pueden cambiar desde la pantalla | Un apartado en P-26.02 | Bajo | Ajustes › Presupuesto |
| Líneas archivadas (`S6-tablero-lineas-archivadas`) | Los pagos que ya hiciste no desaparecen de las gráficas al archivar una tarjeta | Unas 10 líneas | Bajo | Opción (a), como propone el revisor |
| Escala de G-6 (M-R3-1) | El medidor y la tendencia se leen en la misma escala (0 a 60 %) | C1 1.5.1 | Bajo | Opción A, como propone el revisor |
| Valores PROVISIONALES de K54 | N = 5 pagos, 3 meses y margen de 10 % para la sugerencia, y horizontes de 12 y 17 meses | Ninguno: ya están sembrados | Bajo | Confirmarlos tal cual |

## 13. Consumo por agente (aproximado, tomado de cada reporte)

| Agente | Acciones | KB leídos | KB escritos |
|---|---|---|---|
| `u3-v2` | 46 | 195 | 110 |
| `s1-c1-v15` | ~60 | ~235 | ~60 |
| `revisor-especificacion-r1` | ~55 | ~190 | ~18 |
| `correccion-especificacion-r1` | 46 | ~170 | ~40 |
| `revisor-especificacion-r2` | 55 | ~190 | ~17 |
| `correccion-especificacion-r2` | 58 | ~175 | ~38 |
| `revisor-especificacion-r3` | 46 | ~150 | ~17 |
| `api-tablero` | ~45 | ~165 | ~40 |
| `revisor-tablero-backend-r1` | ~47 | ~140 | ~19 |
| `correccion-tablero-backend-r1` | ~33 | ~150 | ~25 |
| `revisor-tablero-backend-r2` | ~42 | ~110 | ~20 |
| `correccion-tablero-backend-r2` | ~27 | ~120 | ~18 |
| `revisor-tablero-backend-r3` | ~36 | ~105 | ~22 |
| `front-tablero` | ≈57 | ≈240 | ≈80 |
| `revisor-tablero-frontend-r1` | ≈40 | ≈230 | ≈30 |
| `correccion-tablero-frontend-r1` | 41 | ~95 | ~18 |
| `revisor-tablero-frontend-r2` | 30 | ~130 | ~40 |
| `correccion-tablero-frontend-r2` | ≈30 | ≈65 | ≈35 |
| `revisor-tablero-frontend-r3` | ≈48 | ≈130 | ≈40 |
| `u5-tablero` | ≈40 | ≈190 | ≈12, más ≈25 MB de evidencia por script |
| `c3-tablero` | ~40 | ~190 | ~30 |
| `revisor-c3-tablero` | 45 | ~120 | ~22 |
| `integrador-wf3c` | ~32 | ~85, más 1 captura | ~45 |
| **Total aproximado** (23 agentes) | **~1,000** | **~3,570** | **~800** (más la evidencia que genera el script) |

Todos quedaron dentro de los topes del protocolo (§4).

## 14. Commit

- **Instrucción:** `git add -A -- . ':!docs/design/sistema' ':!docs/design/maquetas' ':!docs/estado/wfp3'`, y luego el commit con el mensaje del orquestador y sus 4 líneas finales.
- **Comprobación previa** (`git add -A --dry-run`): 352 rutas. 0 `.db`, 0 `.pem` o `.key`, 0 `node_modules`, 0 `dist` y 0 `datos/tmp`; los cubre `.gitignore`. En las 3 carpetas de la Pista 3 no había cambios.
- **Qué entra:**
  - K54, U3 2.0, DESIGN 1.1, los criterios, C1 v1.5 y sus extractos, la migración 00005, el diccionario y las reglas;
  - el servidor: tablero, HTTP/2 y pruebas;
  - el frontend: tablero, barómetro, barra lateral y pruebas;
  - las 164 capturas regeneradas y su evidencia;
  - C3 v1.3 sin firma;
  - los 14 bloqueos, el índice, los reportes de WF-3c, esta integración y el README.
