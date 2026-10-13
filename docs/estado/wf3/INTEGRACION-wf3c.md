# INTEGRACIÓN de WF-3c «Inicio como tablero de control (K54)»

> **Actualización WF-3c bis (12-oct-2026):** `make verify-wf3` y `make verify-all` en verde, C3 v1.4 **firmada**; backend (1 alto) y frontend (1 medio) sin aprobar tras 2 rondas. Ver §15.

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

---

## 15. WF-3c bis · cierre del tablero (12-oct-2026)

Integrador: `integrador-wf3c-bis` (Opus) · Base: HEAD `0197dad` («[WF-3c]») más el árbol de trabajo de WF-3c bis. Las §1 a §14 siguen siendo el cierre de WF-3c y no se reescriben; lo que esta ronda cerró de ellas se dice en §15.3.

Estado del orquestador al integrar:
- Backend: **no aprobado**, 2 rondas. Queda 1 alto (A-1).
- Frontend: **no aprobado**, 2 rondas. Queda 1 medio (M-1) y 5 bajos.
- C3: **aprobada y FIRMADA**, 1 ronda (C3 v1.4).

Commit: «[WF-3c bis] Cierre del tablero: deuda viva de hoy, pagos de todas las líneas, C1 1.5.x, Parámetros del tablero en Ajustes, barómetro, flujos y C3 1.4». Su hash va en `integrador-wf3c-bis.md` y en la salida del orquestador; no va aquí, porque este archivo es parte del commit.

### 15.1 En corto

**HECHO**
- Se aplicaron las 6 decisiones del orquestador sobre `WF3c-ronda3-escalamiento`:
  1. **Deuda viva de hoy (H-4).** `creditos.LiveDebtAt` deja fuera las compras con fecha posterior al corte, que es la regla de `getLiveDebt`. K-3 y P-17 dicen **$3,895.00**, igual que el oráculo al 12-oct-2026 (antes $18,999.22). La «Deuda viva proyectada» (G-5) arranca en **$2,921.25** al cierre de oct-2026, y sus 12 cierres son iguales a `getLiveDebt` del último día de cada mes. PROVISIONAL: `S6-deuda-viva-hoy` (a).
  2. **Pagos de todas las líneas (M-2).** El barómetro, «Créditos del mes» de K-3 y G-5 cuentan también las líneas archivadas. PROVISIONAL: `S6-tablero-lineas-archivadas` (a) y su punto 5. Se quitó el SKIP de `opcion_a_PROVISIONAL`.
  3. **Escala de G-6 (M-R3-1).** `scale_max_pct` entra a `DashboardMonthlyOverview` y el contrato pasa a **C1 v1.5.1** (221 operaciones y 176 rutas, como antes). G-6 y el barómetro usan la misma escala: 0 a 60 % con los umbrales iniciales. En la tabla de `WF3c-ronda3-escalamiento` es la opción (A); los reportes de los agentes la llaman «decisión 3, opción B».
  4. **Parámetros del tablero en Ajustes (H-3).** P-26.02 tiene el bloque «Parámetros del tablero» con las 5 claves de K54, la vista previa «Antes → Después», la validación V-23 y «Restablecer al valor inicial». `updateDashboardSettings` ya tiene pantalla. PROVISIONAL: `C3-ajustes-tablero` (A).
  5. **Barómetro (U5-04).** El `<svg>` raíz lleva `fill="none"`. `contraste-dom.js` mide lo que pintan los hijos y cuenta aparte, en `excluidos`, la pista marcada `data-decorativo` y las separaciones del color del fondo. El eje de la aguja mide 10 × 10 px.
  6. **Flujos a U3 2.0 (H-1 y H-2).** F-1, F-4 y C3-09 leen K-1 y K-4. La ruta A de F-3 sigue la condición (c) de U1 §4 (PROVISIONAL, `C3-flujos` punto 2).
- **Además:**
  - K-2 trae el mes anterior (`getIncomeMonth`, campos `previous_*`).
  - K-3 trae `month_payment_count` y `month_pending_household_cents`, que es «falta pagar» (PROVISIONAL, `S6-deuda-viva-hoy` punto 3).
  - `monthParam` da su propio 400.
  - P-01 Detalle ya no desborda en Android: mide 412 de 412.
  - K-4 ya no pinta «Hogar $973.75» mientras `listPersons` no responde (PROVISIONAL, `U3-tablero-k4-personas`).
- **Capturas U5:** `u5-cierre` regeneró las 164 con el `dist` de esta ronda. `verify-u5` da **0 pares bajo umbral** (antes daba 2).
- **C3 v1.4 FIRMADA** por `revisor-c3-cierre-r1` en la ronda 1: se cumplen C3-01 a C3-11, con 0 altos, 0 medios y 3 bajos (§15.4).
- **Verificación global:** `make verify-wf3` da **EXIT 0**, por primera vez desde WF-3b, y `make verify-all` da **EXIT 0** (§15.2).
- **Servidor de desarrollo:** este integrador lo reinició con el binario y el `dist` actuales. Responde por **HTTP/2** (§15.2).

**NO APROBADO** (salida propuesta en §15.5):
- **Backend, 1 alto (A-1).** Desde la corrección de M-1, el comparador de escenarios (`compareScenarios`) ofrece y **recomienda** líneas archivadas para una compra nueva.
- **Frontend, 1 medio (M-1) y 5 bajos.** Mientras `listPersons` no responde, D-3 de P-01 Detalle dice un instante «Hogar $973.75» sin la parte de Karen. Es un defecto de WF-3a.

### 15.2 Verificación global (la corrió este integrador el 12-oct-2026, de 20:05 a 20:25)

| Línea | Comando | Resultado |
|---|---|---|
| S6-01 | `make verify-wf3` › `verify-wf3-go` | go1.27.2. `go build`, `go vet`, `gofmt -l` y `go mod tidy -diff` sin cambios. `go test ./...`: 33 paquetes OK y 5 sin pruebas. |
| S6-02 | `verify-wf3-cobertura` | K25 OK: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8. Sin umbral: api 82.3, auth 83.6, jobs 92.3, reportes 77.3 y telegram 80.0. |
| S6-03 y S6-04 | `verify-wf3-contrato` | «los 22 archivos coinciden byte a byte con el generador»; «S6-03: 221 operaciones en 176 rutas, 0 diferencias». `TestRutasIgualesAlContrato` y `TestContratoSinSesionYEsqueleto` dan PASS. |
| Interfaz | `verify-wf3-frontend` (`npm run verify`) | «mocks: 221 operaciones de C1 1.5.1». Playwright con la API simulada: **271 passed y 4 skipped**, 0 fallas, en 6.0 min. Las trazas «Explota» del registro vienen de `pruebas/banco`, la prueba de la barrera de errores, y no son fallas. |
| U5-01, U5-03, U5-04, U5-08 y U5-11 | `verify-u5` | 164 de 164 PNG. axe: 0 violaciones serias o críticas. **U5-04: 25818 pares evaluados, bajo umbral = 0.** Hosts: solo `finanzas.local:8443`. C3: 41 filas «capturado» con sus 4 rutas. |
| Total | `make verify-wf3` | **EXIT 0**: «== WF-3: verificación completa OK (Go, cobertura K25, contrato C1, interfaz y evidencia U5)». Tardó 7 min 10 s. |
| Global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | **EXIT 0**: «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». goose v3.28.0 hace up, down-to 0 y up otra vez. K54 OK con las 5 claves. |
| S6-13 | `make verify-wf3-secretos` | 0 coincidencias. |
| S6-11, S6-12 y S6-14 | `make verify-wf3-servidor` (después del reinicio) | mkcert v1.4.4. `finanzas.local` → 192.168.50.56. curl 200. SAN: `DNS:finanzas.local`, `DNS:localhost` y las 2 IP. Manifest con íconos 192 y 512 y `sw.js` 200. 0 cargas de terceros. Escucha solo en `127.0.0.1:8443` y `192.168.50.56:8443`; por Tailscale no responde. Chrome Android: NO VERIFICADO, es manual de Axel. |
| Reinicio | `bash ops/dev-stop.sh --conservar-nombre` y `bash ops/dev-run.sh` | El proceso arrancó a las 20:20:47 con `backend/bin/finanzas`, que `dev-run.sh` compiló a las 20:20:42. No hay ningún `.go` más nuevo. El `dist` es de las 20:13 (lo generó `verify-wf3-frontend`), y las fuentes de `frontend/src` editadas a mano son de las 18:55 o antes. Se conserva `avahi-publish`. |
| HTTP/2 y `dist` | `curl --http2 --cacert <CA de mkcert> -w %{http_version}` | `https://finanzas.local:8443/` da **2** con 200. `https://127.0.0.1:8443/api/v1/auth/me` da 2 con 401, sin sesión. El `index.html` servido es igual, byte por byte, a `frontend/dist/index.html` (`cmp`). |
| C3-07 | el comando del bloque de firma de C3 | SHA-256 `e341d182…d163`, igual que el bloque. |
| Higiene | `grep -ciE 'password\|contraseña\|token'` en `datos/dev/run/finanzas.log` | 0. |

No corrí `verify-wf3-respaldo`, porque esta ronda no tocó respaldos. La última corrida (continuación de `frontend-cierre`) dio EXIT 0.

### 15.3 Qué cerró, contra la lista de §10

| Hallazgo de WF-3c | Estado | Evidencia |
|---|---|---|
| H-1 (alta, C3): F-1, F-3 y F-4 | **cerrado** | Contra el servidor real, con una copia aislada por archivo: `flujos.spec.ts` 7 de 7, `pendientes.spec.ts` 4 de 4 y `f7-f8.spec.ts` 2 de 2. Lo corrieron `c3-cierre` y el revisor de C3. |
| H-2 (alta, C3): `c3-09.spec.ts` | **cerrado** | «120 cifras leídas en pantalla · diferencias != 0: 0», con 13 del tablero y P-17. |
| H-3 (alta, C3): `updateDashboardSettings` sin pantalla | **cerrado, PROVISIONAL** | Bloque de P-26.02 (`S02Presupuesto.tsx:82`). Lo decide `C3-ajustes-tablero` (A). |
| H-4 (alta, C3) y la diferencia 1 de `u5-tablero` | **cerrado, PROVISIONAL** | K-3 y P-17 dan $3,895.00 = `getLiveDebt` = oráculo. Lo decide `S6-deuda-viva-hoy` (a). |
| U5-06 (alta): P-01 Detalle a 475 px | **cerrado** | `scrollWidth` 412 de 412 (revisor de frontend r2). |
| M-R3-1 (media, especificación) | **cerrado** en C1 y en el código | C1 v1.5.1 y `scale_max_pct` 60 en vivo. Falta el texto de U3 G-6 y §2.6 (solicitud en §15.5). |
| M-2 (media, backend) | **cerrado, PROVISIONAL** | Opción (a). Su corrección abrió A-1 (§15.5). |
| U5-04 (media): 2 pares del barómetro | **cerrado** | `verify-u5`: bajo umbral = 0. |
| H-5 (media, C3): F-3A3 | **cerrado** | C3 v1.4 §4. |
| Diferencia 2 de `u5-tablero`: G-5 desde $18,025.47 | **cerrado** | G-5 arranca en $2,921.25. |
| B-R3-2 (baja): campos sin llenar | **cerrado** | `getIncomeMonth` y `getCreditTotals` llenan los campos nuevos. |
| B-N4 (baja): «esquema 1.4» en el `Makefile` (líneas 8, 66 y 149) | **abierto** | Fuera del alcance de esta ronda. |
| N-3 (baja): `Estados.tsx`, `role="status"` | **abierto** | Lo pidió `frontend-cierre` como solicitud fuera de alcance. |
| Los demás bajos de §10 | **abiertos** | Siguen en sus informes `REVISION-*`. |

### 15.4 C3 v1.4: FIRMADA

- **Revisor:** `revisor-c3-cierre-r1`, independiente. Ronda 1. **Veredicto: APROBADO.** El bloque de firma está en `docs/contratos/C3-integracion.md` §9. SHA-256 del contenido firmado: `e341d182e08567d797ae503c3d3970b8d3542d650a632806559fc2df5507d163`, que este integrador recalculó (§15.2).
- **Matriz:** implementado 42 de 42, probado 42 de 42 y capturado 41 de 42. La fila Global sigue en «no» con `C3-capturas-global`. Hay 221 operaciones; 5 no tienen fila y las 5 tienen su excepción (las 4 de K53 y `getOpenApiDocument`).
- **Evidencia propia del revisor, en una copia limpia:**
  - Go: 1151 PASS, 0 FAIL y 0 SKIP.
  - Playwright: 271 pasan y 4 omitidas.
  - Flujos: 8 de 8 contra el servidor real.
  - C3-09: 120 cifras y 0 diferencias. K-3 da $3,895.00 y la deuda viva al 10-abr-2027 da $6,293.44.
- **3 hallazgos bajos (no abren ningún criterio):**
  - **B-1:** F-3A0 (`flujos.spec.ts:71`) sale en verde aunque su paso diga FALLA, porque le falta `expect(reg.ok)`.
  - **B-2:** U1 §4 F-3, U1 §3.16 y U3 tarjeta 5 siguen con +$138.50 (§15.6.7).
  - **B-3:** la línea 3 de C3 todavía dice «sin firma». Se corrige en la próxima versión, con firma nueva, porque es parte del contenido firmado.
- **Con esto se cumple** lo que pedían `C3-ronda3-escalamiento` (WF-3b) y el punto 3 de `WF3c-ronda3-escalamiento`: C3 queda firmada. Las dos fichas siguen abiertas en el índice hasta que Axel ratifique en la puerta 4 lo PROVISIONAL que citan.

### 15.5 Lo que queda sin aprobar y salida propuesta

| Sev. | Id · área | Qué pasa | Cambio exacto | Dueño propuesto |
|---|---|---|---|---|
| **alta** | A-1 · backend (`REVISION-backend-cierre-r2.md`) | La corrección de M-1 quitó el filtro `c.IsActive` de `states` (`backend/internal/api/creditos/motor.go:210`). `compareScenarios` arma `comparador.Request.Credits` con `s.states(nil)` (`condiciones.go:676`), así que genera, puntúa y **recomienda** escenarios en líneas archivadas. Con una línea archivada, el revisor midió 6 escenarios evaluados y la recomendación en la archivada; con el filtro anterior, 3 escenarios y la recomendación en la activa. «Usar escenario» abriría `createCreditPurchase` sobre la línea archivada. Esto contradice `diccionario-datos.md:776` (`credit.is_active`: «deja de ofrecerse en capturas nuevas») y U1 N-053. Con la base de desarrollo de hoy **no se nota**, porque Kueski es la única línea y está activa. | En `condiciones.go:676`, pasar solo las líneas activas (por ejemplo, un `s.activeStates()` que filtre `c.IsActive`) sin tocar `states`. Agregar `TestOp_compareScenarios_LineaArchivada`, que exija la misma recomendación y el mismo `evaluated_count` que con solo la línea activa. Corregir `S6-tablero-lineas-archivadas.md:55` y decir qué operaciones usan `states`. | backend (`api-creditos`) |
| **media** | M-1 · frontend (`REVISION-frontend-cierre-r2.md`) | `frontend/src/pantallas/inicio/Detalle.tsx:92-98` (D-3): mientras `listPersons` no responde, cada fila dice «Hogar $973.75». Al llegar la respuesta cambia a «Hogar $598.82 · Karen $374.93». Viene de WF-3a. `Calendario.tsx:33-42` tiene el mismo patrón, pero solo en la insignia. | Lo mismo que se hizo en K-4: pintar las filas solo con `c.personas.datos !== undefined`, con esqueleto mientras tanto, y el estado de error con «Reintentar» si falla (PROVISIONAL `U3-tablero-k4-personas`, que habría que ampliar a D-3). Agregar una prueba con `/api/v1/persons` retenido en `?sub=detalle` y un `MutationObserver`. | frontend |
| baja | B-1 a B-5 · frontend | B-1: F-3A0 sin `expect(reg.ok)` (es el mismo B-1 de C3). B-2: U4 P-26.02 no describe el bloque, y `C3-ajustes-tablero.md:21` sigue diciendo «No se tomó ninguna decisión». B-3: `tableroNoCargo` no usa el estado común «Error al cargar · Reintentar». B-4: evidencia vieja, **ya cerrado de hecho** (`verify-u5` en 0 y C3 firmada con la matriz igual a la generada). B-5: comentario viejo en `c3-09.spec.ts:55-56`. | Ver el informe. | frontend, U4 y U5 |
| sin calificar | solicitud de `backend-cierre` | `backend/internal/api/gastos-inicio/vistas.go:235` (`deriveCredits`) todavía usa `mc.LiveDebt`. La vista de BD de créditos sigue contando las compras futuras. | Cambiar a `creditos.LiveDebtAt(...)` con los mismos argumentos. | backend |
| sin calificar | observación de `u5-cierre` | `GET /api/v1/debt/live` sin `date` respondió `date` 2027-09-30 y `total_cents` 0. | Confirmar la fecha por omisión contra C1 y corregirla si no es «hoy». | backend |
| sin calificar | solicitud de `c3-cierre` | `frontend/pruebas-integracion/humo-inicio.spec.ts` espera las 6 tarjetas de U3 1.x. No entra en ningún `verify-*`. | Leer K-1 a K-4 y B-1 y compararlos con la API. | frontend o U5 |
| sin calificar | solicitud de `c3-cierre` | Hoy se puede correr `flujos.spec.ts` contra la base de desarrollo por omisión (§15.9). | Que `servidor-prueba.sh arrancar` exporte o imprima `FH_BASE`, o que las pruebas que escriben se nieguen a correr contra 8443. | U5 |
| sin calificar | solicitudes de documentos | `docs/criterios/C1.md` (líneas 6, 9, 106 y 107) sigue en 1.5.0, y la Verifica literal de C1-13 busca `version: "1.5.0"`. U3 G-6 y §2.6 no citan `DashboardMonthlyOverview.scale_max_pct`. | Cambiar a 1.5.1 con el nombre «C1 v1.5.1 (escala de G-6, WF-3c bis)», y citar el campo en U3. | criterios y especificación |

**PROPUESTA del integrador (no decidido).** Backend y frontend ya hicieron sus 2 rondas.

| Opción | Resultado | Costo | Riesgo | Recomendación |
|---|---|---|---|---|
| (A) Corrección acotada antes de la puerta 4 | Backend y frontend aprobables; el comparador nunca ofrece una tarjeta archivada y Detalle nunca muestra una cifra repartida mal | 2 agentes cortos (unas 15 líneas de Go con su prueba; D-3, Calendario y su prueba) y 1 revisión de cada uno | Bajo: no cambia ninguna cifra de los motores ni del oráculo, y C3 no se toca | **Sí** |
| (B) Ir a la puerta 4 tal cual, con A-1 y M-1 anotados | La puerta 4 no espera | Ninguno ahora | Medio: si Axel archiva una tarjeta, el comparador puede recomendarla para una compra nueva | No |
| (C) Volver a filtrar `states` por `c.IsActive` | Arregla el comparador | Una línea | Alto: K-3 y G-5 dejarían de coincidir con `getLiveDebt` cuando haya una línea archivada (vuelve M-1 de la ronda 1) | No |

### 15.6 Pendientes para la puerta 4 (consolidado; PROPUESTA, no decidido)

#### 15.6.1 Valores K31 PROVISIONALES

WF-3c y WF-3c bis **no agregaron valores K31**: `grep K31` en los 6 reportes de autor de WF-3c bis y en sus 3 bloqueos nuevos da 0, y los 17 bloqueos de WF-3c y WF-3c bis tienen «—» en la columna K31. La lista sigue siendo la de `INTEGRACION-wf3b.md` §9, que viene de `INTEGRACION-wf3a.md` §9:

| Tema | Valor PROVISIONAL en el código | Bloqueo |
|---|---|---|
| Respaldos | Diario a las 03:15 (America/Mexico_City); se guardan 30 diarios. La restauración reemplaza todo salvo `backup` y `device_session` | `S6-respaldos`, `S6-valores-operativos-base` |
| Adjuntos | JPG, PNG, WebP y PDF, por contenido (415 si no); 10 MB | `S6-adjuntos` |
| Cabeceras y CSP | `nosniff`, `no-referrer`, `DENY` y `no-store` en la API; CSP `default-src 'self'` … `frame-ancestors 'none'` | `S6-valores-operativos-base`, `WF3-frontend-servidor` |
| Resúmenes y reporte | Diario: fecha, gastos de hoy y del mes, y pagos de 7 días. Semanal: semana ISO, gastos de la semana y del mes, pagos de 7 días y casillas sin armar. Mensual: resumen, devengado contra efectivo, categorías, quincenas, créditos y parte de Karen | `S7-fuentes-avisos`, `S7-reporte-mensual-secciones` |
| Avisos | Revisión cada minuto; un envío fallido queda `failed` y no se reintenta solo | `S7-fuentes-avisos` |
| Contraseñas, límites e idempotencia | argon2id t=3, m=64 MiB, p=4; 8 MiB por JSON; purga cada hora, al minuto 5 | `S6-valores-operativos-base` |
| PWA sin conexión | `/health` cada 15 s; la cola se reintenta cada 30 s | `WF3-frontend-k31` |

#### 15.6.2 Valores PROVISIONALES de K54 (todos sembrados por la 00005 o fijos en U3 2.0)

| Valor | PROVISIONAL hoy | Bloqueo |
|---|---|---|
| Fin del verde del barómetro (`debt.capacity_warning_pct`) | 30 %, de 1 a (tope − 1) | `U3-tablero-barometro`, `C1-tablero` punto 2 |
| Tope del amarillo (`debt.capacity_limit_pct`) | 40 %, de (fin del verde + 1) a 100 | los mismos |
| Escala del barómetro y de G-6 | múltiplo de 10 ≥ 1.5 × tope: **0 a 60 %** | `U3-tablero-barometro` punto 4 |
| Ingreso mensual neto | quincenas capturadas más esperadas, sin extras | `U3-tablero-barometro` punto 1 |
| Pagos que cuentan | solo la parte del hogar | `U3-tablero-barometro` punto 2 |
| Redondeo y colores | entero con ROUND_HALF_UP: ≤ 30 verde, 31 a 40 amarillo y ≥ 41 rojo | `U3-tablero-barometro` punto 3 |
| Filas de «Pagos próximos» (`dashboard.payment_reminder_count`) | 5, de 1 a 20; vencidos y luego de hoy a hoy + 30 días | `U3-tablero-recordatorio` |
| Sugerencia de presupuesto (`budget.signal_months` y `budget.signal_margin_pct`) | 3 meses cerrados (de 2 a 12) y margen de 10 puntos (de 1 a 50); un mes sin plan corta la racha | `U3-tablero-senal-presupuesto` |
| Horizontes | G-1: 11 meses anteriores más el actual. G-6: 6 anteriores, el actual y 5 siguientes. Una sola petición de 17 meses | `U3-tablero-horizonte` |
| Unidad que se ve en P-26.02 | «%», «%», «pagos», «meses» y «puntos», puestas en `textos.ts` (puente B) | `C1-unidad-parametros-tablero` (se recomienda la A) |
| Deuda viva a una fecha y «falta pagar» | solo las compras hasta el corte; falta pagar = parte del hogar menos los abonos del hogar | `S6-deuda-viva-hoy` (a) |
| Líneas archivadas | cuentan en el barómetro, K-3 y G-5 | `S6-tablero-lineas-archivadas` (a) y punto 5 |

#### 15.6.3 Bloqueos abiertos por fase de origen (conteo de este integrador sobre la carpeta, 154 archivos)

Regla de la puerta 2: un archivo cuenta como resuelto si su última línea «Estado del bloqueo» dice RESUELTO o, si no la tiene, si trae «Respuesta de Axel» o «Resolución».

| Origen | Archivos | Resueltos | Abiertos | Quién decide los abiertos |
|---|---|---|---|---|
| WF-0, puerta 1, P1 y P1c | 43 | 43 | 0 | — |
| WF-1 | 26 | 21 | 5 | puerta 4: `S2-flujo-quincenal`, `S2-ventanas`, `S5-carrito-comparativa`, `S5-carritos-detalle` y `S5-profeco-filtro` |
| P2c | 13 | 0 | 13 | puerta 4 |
| P2d | 1 | 0 | 1 | puerta 4: `S5-agente-externo` |
| Pista 3 | 2 | 0 | 2 | Pista 3 |
| WF-3a (3 de la semilla y 37 más) | 40 | 0 | 40 | puerta 4 y WF-4; 8 con valores K31 |
| WF-3a bis | 4 | 0 | 4 | puerta 4 |
| WF-3b | 8 | 0 | 8 | puerta 4 |
| WF-3c | 14 | 0 | 14 | puerta 4 |
| WF-3c bis | 3 | 0 | 3 | puerta 4 |
| **Total** | **154** | **64** | **90** | |

#### 15.6.4 Preguntas `U3-tablero-*`

- `U3-tablero-barometro`: (1) qué es el ingreso mensual neto; (2) qué pagos cuentan; (3) redondeo y colores; (4) escala del arco. La opción PROVISIONAL es la (a) de cada una (§15.6.2). Además está B-R3-1: la justificación del punto 4 trae una cuenta errónea.
- `U3-tablero-horizonte`: los meses que ve cada gráfica (a).
- `U3-tablero-mas-del-hogar`: las tarjetas 2, 4, 5 y 6 siguen sin cambio bajo «Más del hogar» (a).
- `U3-tablero-recordatorio`: N, el rango y la forma de K-4 (a).
- `U3-tablero-senal-presupuesto`: los meses, el margen y la racha de G-3 (a).
- `U3-tablero-k4-personas` (nuevo): qué hace K-4 si `listPersons` falla. La PROVISIONAL es el estado de error del bloque con «Reintentar». Si se confirma, D-3 tiene que hacer lo mismo (M-1 de §15.5).

#### 15.6.5 `WF3a-bis-semilla-dev` (semilla de desarrollo, 7 puntos)

Son preguntas sobre lo que siembra `finanzas semilla --desarrollo`, que nunca va a producción:
1. La renta no cabe con P1–P6 (PREGUNTA).
2. El candado «nunca en producción» (PROVISIONAL).
3. El presupuesto solo trae Despensa (PREGUNTA).
4. P3–P6 de `getAccrualVsCash` contra `model.json` (resuelto de forma provisional en WF-3b).
5. El costo y las kcal de la semana quedan en 0 (PREGUNTA).
6. No hay ingresos capturados.
7. Los precios de Walmart en la semilla, que chocan con Q2-49 (punto siguiente).

#### 15.6.6 `C3-flujos`

- **Punto 2 (F-3, ruta A):** rige la (c) PROVISIONAL. La ruta A exige una lista vigente con todos sus precios en al menos una tienda; si falta alguno, la tarjeta lleva a «Registrar precios». Axel la ratifica. La (b), que la tarjeta abra P-14 con la compra elegida, queda como PROPUESTA.
- **Punto 3:** lo contesta K46 (+$264.50; ver §15.6.7).
- **Punto 4:** F-5 con oct-2026 y 6 toques.
- **Observación sobre Q2-49:** Q2-49 = A dice «la app arranca sin precios de Walmart», pero la semilla trae los 36 precios de Walmart de `data.json` (`WF3a-bis-semilla-dev` punto 7). Sin ellos, la ruta B no tiene con qué comparar. Que Axel diga cuál de las dos vale.

#### 15.6.7 +$138.50 contra +$264.50 (compra semilla del 10-oct-2026)

- **HECHO:**
  - Dicen Walmart **$3,233.00 (+$138.50)**: U1 §3.16 (`U1-arquitectura.md:652`), U1 §4 F-3 (`:978`) y el «Ejemplo trazable» de la tarjeta 5 de U3 (`U3-dashboard.md:317`). Es el canon carrito contra carrito.
  - K46 aprobó Q2-47 = A (`basket = expense`): cada tienda se cotiza con los productos y las cantidades de la compra base. Con esa regla, la API (`getBasketComparison`: Walmart 335900 contra Bodega Aurrera 309450) y la pantalla dan **$3,359.00 (+$264.50)**. La diferencia es +$21.00 por la sexta lata de atún y +$105.00 por los dos Coffee Mate de 640 g que cubren 780 g.
  - En idénticos, las dos versiones dan +$100.50.
- **Pendiente:** el dueño de U1 y U3 actualiza esas 3 líneas a la cifra de K46. Según `C3-flujos` punto 3, no hace falta una respuesta nueva de Axel. Es el bajo B-2 de la firma de C3.

#### 15.6.8 `S5-agente-externo` (K53)

Preguntas (a) a (f). Las PROPUESTAS son:
- (a) los precios del agente compiten como `manual`;
- (b) el agente solo escribe productos, presentaciones, precios y recetas;
- (c) sus recetas entran como borrador;
- (d) `price.source` = `site_assisted` (la d2);
- (e) solo el transporte `serviceToken` de C1 1.4;
- (f) borrador y origen de una receta, sin propuesta.

#### 15.6.9 Lo que ya venía de WF-3b

- `U5-ronda3-escalamiento` y `C3-ronda3-escalamiento`. C3 ya quedó firmada, así que el segundo solo espera que Axel ratifique.
- S6-11 en los Android de Axel (manual, `ops/mkcert.md` §4).
- S7-07, el envío real por Telegram. Lo hace solo el orquestador.

#### 15.6.10 Tabla para Axel (PROPUESTA)

| Qué decide | Resultado | Costo | Riesgo | Recomendación |
|---|---|---|---|---|
| Corrección acotada de A-1 y M-1 (§15.5) | El comparador no ofrece tarjetas archivadas y Detalle no muestra repartos mal ni un instante | 2 agentes cortos y 2 revisiones | Bajo | Sí, antes de la puerta 4 |
| Deuda viva de hoy (`S6-deuda-viva-hoy`) | Una sola cifra: $3,895.00 hoy, igual que el oráculo; una despensa programada suma el día de su compra | Ninguno: ya está | Bajo | Confirmar la (a) |
| Líneas archivadas (`S6-tablero-lineas-archivadas`) | Archivar una tarjeta con pagos no la borra del barómetro ni de la deuda hasta que se termina de pagar | Ninguno: ya está | Bajo, con A-1 corregido | Confirmar la (a) y el punto 5 |
| Parámetros del tablero en Ajustes › Presupuesto | Los umbrales y los demás valores se cambian desde la pantalla | Ninguno: ya está | Bajo | Confirmar `C3-ajustes-tablero` (A) |
| Unidad de los parámetros (`C1-unidad-parametros-tablero`) | Unidades cortas que vienen del servidor | Una migración nueva de anotaciones, sin editar la 00005 | Bajo | La (A); se retira el puente (B) |
| K-4 y D-3 si falla `listPersons` | Nunca se ve una cifra mal repartida | Ninguno en K-4; D-3 en la corrección acotada | Bajo | La (A): error con «Reintentar» |
| F-3 ruta A (`C3-flujos` punto 2) | La tarjeta «Mejor canasta» solo recomienda tienda si todo tiene precio | Ninguno: ya está | Bajo | Confirmar la (c) |
| Valores PROVISIONALES de K54 (§15.6.2) | El tablero se queda como está | Ninguno | Bajo | Confirmarlos tal cual |
| Ejemplos de U1 y U3 (§15.6.7) | Los documentos dicen lo mismo que la pantalla | 3 líneas de texto | Ninguno | +$264.50 (K46) |

### 15.7 Rondas y dictámenes de WF-3c bis

| Área | Autor | Rondas | Dictamen final | Abiertos |
|---|---|---|---|---|
| Backend (H-4, M-2, M-R3-1, K-2 y K-3) | `backend-cierre` y `correccion-backend-cierre-r1` | 2 (r1: 3 medios y 4 bajos, todos cerrados en r2) | **NO APROBADO** | 1 alto nuevo (A-1) |
| Frontend (P-26.02, barómetro, Detalle, G-6, flujos y C3-09) | `frontend-cierre` (más una continuación) y `correccion-frontend-cierre-r1` | 2 (r1: 2 medios, cerrados en r2) | **NO APROBADO** | 1 medio (M-1, de WF-3a) y 5 bajos |
| U5 (164 capturas) | `u5-cierre` | sin revisión propia; lo comprueban `verify-u5` y C3-05 | — | 2 observaciones (§15.5) |
| C3 v1.4 | `c3-cierre` | 1 | **APROBADO y FIRMADO** | 3 bajos |

### 15.8 Bloqueos de WF-3c bis

- **Nuevos (3, ABIERTOS):**
  - `S6-deuda-viva-hoy`, de `backend-cierre`;
  - `C1-unidad-parametros-tablero` y `U3-tablero-k4-personas`, de `correccion-frontend-cierre-r1`.

  Este integrador dio de alta en `INDICE.md` las 2 fichas del frontend, que no tenían fila.
- **Actualizados:** `S6-tablero-lineas-archivadas`, con el punto 5, y `S6-tablero-sin-datos`.
- **Atendidos en WF-3c bis, pero abiertos hasta la ratificación de la puerta 4:** `WF3c-ronda3-escalamiento` y `C3-ajustes-tablero`.
- Ninguno toca un valor K31.

### 15.9 Incidente de la ronda (ya corregido)

`c3-cierre` lanzó una primera pasada de los flujos sin `FH_BASE=https://127.0.0.1:8444`, y escribió en la base de desarrollo (`datos/finanzas-dev.db`, el servidor de 8443). Lo detectó porque C3-09 dio 19 diferencias. Lo corrigió así:
- detuvo el servidor;
- apartó la base y los 2 archivos que había creado;
- reconstruyó la base con `finanzas migrate`, `finanzas semilla` y `finanzas semilla --desarrollo`;
- comparó tabla por tabla contra el estado documentado.

No se cita ninguna cifra de esa pasada. Se perdieron las sesiones de desarrollo, así que hay que volver a entrar. La prevención está en §15.5 (`servidor-prueba.sh`).

### 15.10 Consumo por agente (aproximado, tomado de cada reporte)

| Agente | Acciones | KB leídos | KB escritos |
|---|---|---|---|
| `backend-cierre` | 57 | ~150 | ~30 (sin lo que regeneró `build.py`) |
| `revisor-backend-cierre-r1` | ~50 | ~190 | ~25 |
| `correccion-backend-cierre-r1` | 29 | ~95 | ~22 |
| `revisor-backend-cierre-r2` | 42 | ~110 | ~14 |
| `frontend-cierre` | ≈65 (cerró por tope, por encima de las 60 del protocolo §4) | ≈230 | ≈30 |
| `frontend-cierre`, continuación 1 | ≈24 | ≈45 | ≈7 |
| `revisor-frontend-cierre-r1` | ≈52 | ≈150 | ≈16 |
| `correccion-frontend-cierre-r1` | ≈47 | ≈115 | ≈22 |
| `revisor-frontend-cierre-r2` | ≈58 | ≈120 | ≈16 |
| `u5-cierre` | ≈22 | ≈60 | ≈6, más la evidencia que genera el script |
| `c3-cierre` | ≈50 | ≈210 | ≈28 |
| `revisor-c3-cierre-r1` | 55 | ~125 | ~16 |
| `integrador-wf3c-bis` | ~45 | ~200 | ~50 |
| **Total aproximado** (13 entradas) | **~596** | **~1,800** | **~283** |

Solo `frontend-cierre` pasó un tope (las acciones); su continuación terminó los pendientes. Los demás quedaron dentro de los topes del protocolo (§4).

### 15.11 Commit de WF-3c bis

- **Instrucción:** `git add -A -- . ':!docs/design/sistema' ':!docs/design/maquetas' ':!docs/estado/wfp3'`, y luego el commit con el mensaje del orquestador y sus 4 líneas finales.
- **Comprobación previa** (`git add -A --dry-run` con las mismas exclusiones): 171 rutas. 0 `.db`, 0 `.pem` o `.key`, 0 `node_modules`, 0 `dist` y 0 `datos/tmp`. En las 3 carpetas de la Pista 3 no había cambios.
- **Qué entra:**
  - el servidor: `LiveDebtAt`, los pagos de todas las líneas, `previous_*` y los campos de K-3, `scale_max_pct` y el 400 de `monthParam`, con sus pruebas;
  - C1 v1.5.1: c1gen, `openapi.yaml`, `C1-datos.md`, los extractos, `routes.go` y `types_gen.go`;
  - `reglas-negocio.md` (RN-21 y RN-30);
  - el frontend: P-26.02 «Parámetros del tablero», el barómetro, K-4, la insignia, G-6, `contraste-dom.js` y las pruebas;
  - las 164 capturas regeneradas y su evidencia;
  - C3 v1.4 firmada;
  - los 3 bloqueos nuevos y los 2 actualizados, el índice, los reportes de WF-3c bis, esta sección y el README.
