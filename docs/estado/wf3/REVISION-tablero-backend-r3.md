# REVISIÓN · tablero (backend) · ronda 3 de 3 · WF-3c · 12-oct-2026

Revisor independiente `revisor-tablero-backend-r3`. Evalúa el trabajo de `api-tablero` y sus correcciones (`docs/estado/wf3/api-tablero.md`, `correccion-tablero-backend-r1.md` y `correccion-tablero-backend-r2.md`) **solo** contra `docs/criterios/S6.md` (S6-01 a S6-04, S6-14 y S6-15), S2.md y S3.md (cifras del motor) y C1.md. No corrige nada. Todas las verificaciones las reejecuté con comandos propios. Los temporales (script de curl, overlays, bitácoras de make y el archivo de cookies) quedaron en `datos/tmp/revisor-tablero-backend-r3/` y se borraron al terminar.

**Dictamen: NO APROBADO.** Quedan 0 altos, 1 medio y 6 bajos.
- **M-2 (ronda 1 y 2): SIGUE ABIERTO.** La corrección de la ronda 2 dejó el bloqueo honesto y una parte de la prueba corriendo, pero el código no cambió. Hoy el servidor cuenta solo las líneas activas, y RN-30 y C1 dicen «de todas las líneas».
- Las cifras coinciden con el oráculo en el caso sembrado: 6 grupos de valores y la serie de 12 meses.
- No hay verificaciones rotas: `verify-wf3-go`, `verify-wf3-contrato`, `verify-wf3-cobertura` y `verify-all` salen con EXIT 0.
- Es la ronda 3 de 3, así que **M-2 se escala**. El orquestador decide el alcance (salida i o ii) y Axel decide la regla en la puerta 4.

## 1. Hallazgos de la ronda anterior

| Hallazgo | Estado | Evidencia propia |
|---|---|---|
| **M-2** · los pagos de crédito del tablero dejan fuera las líneas archivadas | **ABIERTO** (media) | Ver §3. |
| B-1 · 400 en la API y 422 en el texto de C1 | Sigue | `"4"` → 400 `bad_request`; `monthly-overview` sin `from` → 400. |
| B-2 · mensaje 400 genérico | Sigue | «Tipo incorrecto: los importes son enteros en centavos y las fechas, texto AAAA-MM-DD.» para un porcentaje. |
| B-3 · `bloqueos/INDICE.md` sin filas del tablero | Sigue | `grep -c tablero` = 0 (mtime 13:25, anterior a los bloqueos nuevos). |
| B-4 · pruebas de 4xx con `c.Raw` | Sigue | `tablero_test.go` 92, 162 y 171; `tablero_datos_test.go` 220, 223 y 318. |
| B-5 · texto de `BudgetSignals` sin la opción (a) de S6-tablero-senal-periodo | Sigue | `inicio.yaml` 1576 sin cambio; 0 menciones de `senal-periodo` en el extracto. |

## 2. Verificación por criterio

| Criterio | Resultado | Evidencia propia |
|---|---|---|
| S6-01 | CUMPLE | `make verify-wf3-go` EXIT 0: go1.27.2 y go.mod 1.27; build y vet sin errores; `gofmt -l .` vacío; `go mod tidy -diff` sin cambios; `go test ./...` con 33 paquetes OK y 5 sin pruebas. |
| S6-02 | CUMPLE | `make verify-wf3-cobertura` EXIT 0. Con umbral K25: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8. Sin umbral: api 82.0, auth 83.6, jobs 93.7, reportes 77.3 y telegram 80.0. Con `go test -cover`: api/tablero **87.8 %**, internal/api 88.8 %, api/creditos 80.4 % y dominio 87.0 %. |
| S6-03 | CUMPLE | `make verify-wf3-contrato` EXIT 0: «los 22 archivos coinciden byte a byte con el generador» y «221 operaciones en 176 rutas, 0 diferencias». `TestContrato_Operaciones`: «221 de 221 operaciones cumplen (2xx declarada, 4xx válidas contra openapi.yaml)», incluidas `getDashboardDebtCapacity`, `getDashboardMonthlyOverview`, `getBudgetSignals`, `getDashboardSettings` y `updateDashboardSettings`. En api/tablero pasan 13 pruebas y queda 1 SKIP, el subtest `opcion_a_PROVISIONAL` (M-2). |
| S6-04 | CUMPLE con observaciones (B-1, B-2, B-4 y B-6) | Mes `AAAA-MM`, centavos enteros, `ErrorBody` con `code`, `message`, `fields` y `request_id`; 400 de formato y 422 de regla. |
| S6-14 | CUMPLE | `ss -ltn` muestra solo `127.0.0.1:8443` y `192.168.50.56:8443`. Por la IP de Tailscale (100.115.28.35), curl sale con código 7. Pasan `TestEscucha_NegociaHTTP2`, `TestEscucha_SoloLANYLocal` y `TestEscucha_HTTPSConCertificado`. |
| S6-15 | CUMPLE con observación (B-2) | Identificadores en inglés. Mensajes en español de México («Falta «from»: un mes con el formato AAAA-MM, por ejemplo 2026-10.», «El umbral amarillo (45 %) tiene que ser menor que el tope (2 %). Baja el amarillo o sube el tope.»). `datos/dev/run/finanzas.log` tiene 0 coincidencias de password, contraseña, cookie o token, incluidas mi sesión y mi logout. |
| S2 y S3 (cifras) | CUMPLE en el caso sembrado; con una línea archivada, ver M-2 | §2.1. |
| C1 | CUMPLE salvo M-2, B-1 y B-5 | Rutas y esquemas iguales al contrato. RN-30 «de todas las líneas» no se cumple con líneas archivadas (M-2). |
| `make verify-all GOOSE=…` | EXIT 0 | «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». |

### 2.1 Cifras contra el oráculo y los motores (servidor de desarrollo, HTTP/2)

**Servidor.**
- Binario de las 15:52:13; proceso arrancado a las 15:52:32.
- Ningún fuente de Go fuera de pruebas es posterior al binario (`find -newer` vacío).
- La base de desarrollo tiene 1 línea de crédito, activa, y 0 archivadas (consulta de solo lectura).

**Sesión.** curl con la CA de mkcert y una sesión propia de axel. La contraseña de desarrollo se leyó del archivo y se pasó por la entrada estándar, sin imprimirse. Al final: logout 204 y luego `auth/me` 401.

**Recálculo.** Hice el recálculo desde `oraculo/model.json`, agrupando `despensas[].pagos` por mes calendario (la parte de axel y el monto). Los porcentajes van con ROUND_HALF_UP sobre 2,450,000.

| # | Valor | API | Oráculo (recalculado) | ¿Coincide? |
|---|---|---|---|---|
| 1 | Deuda viva (`getCreditTotals.live_debt_cents`, S3) | 1,899,922 | Σ de los pagos con fecha ≥ 12-oct = 1,899,922 | Sí |
| 2 | Pagos de octubre, hogar · total · conteo (`debt-capacity`) | 59,882 · 97,375 · 1 | 59,882 · 97,375 · 1 (pago del 30-oct). Además es igual a `getCreditTotals.month_household_cents` y `month_total_cents` | Sí |
| 3 | Ratio y ratio total, oct · nov · dic | 2/4 · 5/8 · 6/9 | 2.44 → 2 y 3.97 → 4; 4.89 → 5 y 7.95 → 8; 6.05 → 6 y 9.11 → 9 | Sí |
| 4 | Margen (`headroom_cents`), oct · nov · dic | 920,118 · 860,236 · 831,743 | 980,000 − pagos | Sí |
| 5 | Serie de 12 meses (`monthly-overview` jul-2026 a jun-2027): `credit_payments_cents`, `credit_payments_total_cents` y `debt_ratio_pct` | jul–sep 0; de oct a jun: 59,882 · 119,764 · 148,257 · 176,750 · 176,752 · 176,750 · 176,752 · 176,750 · 88,377; ratio 2·5·6·7·7·7·7·7·4 | Los 12 meses iguales, también la parte con Karen (97,375 … 125,870) | Sí (12 de 12) |
| 6 | Devengado, efectivo y plan de oct (S2) | 119,764 · 59,882 · 200,000; `month-vs-plan` 60 % | `periodos[0].devengado.axel` 1,197.64 y `efectivo.axel` 598.82; 119,764 ÷ 200,000 = 59.88 → 60 | Sí |

**(c) HTTP/2.**
- `%{http_version}` = 2 por loopback y por la IP de la LAN (`--resolve finanzas.local:8443:192.168.50.56`).
- `openssl s_client -alpn h2,http/1.1` da «ALPN protocol: h2».
- Con `--http1.1` responde 200 en HTTP/1.1.
- El cambio en `tls.go` es mínimo: `NextProtos: alpnProtos` con `h2` y `http/1.1`.

**(d) Umbrales desde `app_setting`.**
- Al inicio, los 5 parámetros estaban en `is_initial` true (30, 40, 5, 3 y 10).
- `PATCH` con amarillo 1 y tope 2 → 200:
  - octubre pasa a `warning`, con escala 10 y margen −10,882;
  - noviembre pasa a `exceeded`, con margen −70,764;
  - `monthly-overview` de oct a dic da `warning`, `exceeded` y `exceeded`.
- Amarillo 45 con tope 2 → 422 `out_of_range` en español.
- `"4"` → 400 (B-1, B-2).
- `reset_parameters` → 200, con los 5 en `is_initial` true. Octubre vuelve a `normal`, 30/40, escala 60 y margen 920,118.
- La base de desarrollo quedó con sus valores iniciales.

**(e) Sin cifras a mano.**
- Los pagos salen de `motor_cred.UnifiedCalendar` (por `creditos.PaymentsByMonth`).
- El ingreso, el devengado, el efectivo y el plan salen de `motor_fin`.
- El manejador solo suma las partes de ingreso del motor del mes (`tablero.go`, líneas 122-123) y divide con `dominio.PercentOf` y `MulDivRound` (ROUND_HALF_UP).
- La única resta es la del margen (línea 164): tope × ingreso ÷ 100 − pagos, como define C1.

## 3. Hallazgos

### M-2 (media, sigue de las rondas 1 y 2) · El servidor cuenta solo las líneas activas; RN-30 y C1 dicen «de todas las líneas», y la opción PROVISIONAL del bloqueo no está en el código

**Dónde:**
- `backend/internal/api/creditos/tablero.go`, líneas 26-50 (sin cambios desde las 14:45:49). `PaymentsByMonth` usa `s.states(nil)`, que filtra `c.IsActive` (`motor.go`, líneas 205-217), y devuelve `true` fijo. El comentario de las líneas 11-15 dice «con las líneas activas» sin citar el bloqueo.
- `backend/internal/api/tablero/tablero_datos_test.go`, línea 408: `t.Skip` del subtest `opcion_a_PROVISIONAL`.
- Contra qué choca:
  - `docs/contratos/reglas-negocio.md`, RN-30 (línea 345: «de todas las líneas»);
  - `docs/contratos/extractos/inicio.yaml`, línea 1448: `payments_cents` «de todas las líneas»;
  - `docs/estado/bloqueos/S6-tablero-lineas-archivadas.md`, con la propuesta PROVISIONAL (a) y la (a1).

**HECHO · qué cambió en la ronda 2:**
- El bloqueo ya no dice algo falso. Ahora dice «Todavía no rige en el código: hoy el servidor hace lo de la opción (b)» y agrega el punto 4, con las salidas (i) y (ii) para el orquestador.
- `TestOp_Dashboard_LineaArchivada` corre siempre la parte común a todas las opciones y pasa:
  - el efectivo de S2 se queda en 700,000;
  - G-1/G-6 y el barómetro dan la misma cifra;
  - con `has_credits` false, los pagos son 0.

**HECHO · por qué sigue abierto:**
- No se aplicó la salida (i) ni la (ii) que proponía la ronda 2. El código sigue haciendo lo de la (b).
- Contradice RN-30, la descripción de `payments_cents` en C1 y la opción PROVISIONAL que propone el bloqueo.
- Lo que demuestra la regla sigue en SKIP, así que `go test` pasa sin comprobarla.
- Es una inconsistencia entre artefactos (contrato, bloqueo y código) que cambia cifras visibles en G-1, G-6 y el barómetro en cuanto haya una línea archivada.
- La base de desarrollo no tiene líneas archivadas, así que hoy las cifras del servidor no cambian.

**HECHO · comprobado con overlays propios** (`go test -overlay`, sin tocar el árbol):

| Código | Prueba | Resultado |
|---|---|---|
| actual | sin el `t.Skip` | **FALLA**: «octubre con la línea archivada: pagos 0 (antes 650000; parte de crédito del efectivo S2 650000), total 0, 0 %, normal»; `debt-capacity` de octubre y noviembre con `PaymentsCents:0` |
| con el diff de la solicitud 1 de `correccion-tablero-backend-r2.md` | sin el `t.Skip` | **PASA**: `api/tablero`, `api/creditos` e `internal/api` OK; `go vet ./internal/api/creditos/` sin errores |

**Corrección sugerida.** Es la misma de la ronda 2. La decisión ya no es de un agente de corrección.
- **Salida (i)**, la que propone el bloqueo. Hay que ampliar el alcance a `backend/internal/api/creditos/tablero.go` y a `c1gen`, o que lo aplique el dueño de `api/creditos`.
  - Aplicar tal cual la solicitud 1 de `correccion-tablero-backend-r2.md`:
    - recorrer `s.credits` con `s.state(c)`;
    - `hasCredits = hasCredits || c.IsActive`;
    - después del calendario, `hasCredits = hasCredits || len(cal.Entries) > 0`;
    - devolver `hasCredits` y actualizar el comentario citando el bloqueo.
  - Quitar el `t.Skip` de la línea 408 de `tablero_datos_test.go`.
  - Aplicar la solicitud 2: el texto de `has_credits` en `gen2_schemas.py`, línea 529, y regenerar con `build.py`.
- **Salida (ii)**, si no se amplía el alcance. Que el bloqueo declare que hasta la puerta 4 rige la (b), que el subtest compruebe la (b) sin SKIP, y pedir que se corrijan RN-30 y la descripción de `payments_cents` (M-20/M-21) por `c1gen`.
- **PREGUNTA (al orquestador):** ¿la salida (i) o la (ii)? **PREGUNTA (a Axel, puerta 4):** la del bloqueo, ¿los pagos de una tarjeta archivada siguen en el tablero (a) o desaparecen (b)?

### B-1 (baja, sigue) · 400 en la API, 422 en el texto de C1
- `"debt.capacity_warning_pct":"4"` → 400 `bad_request`, pero `ajustes-e.yaml` dice «de otro tipo responde 422».
- `monthly-overview` sin `from` → 400, pero `inicio.yaml`, línea 586, dice 422.
- Hay que alinear el texto en `c1gen` (`gen3_ops.py`) o el código.

### B-2 (baja, sigue) · Mensaje 400 genérico en `updateDashboardSettings`
El campo es un porcentaje y el mensaje dice «Tipo incorrecto: los importes son enteros en centavos y las fechas, texto AAAA-MM-DD.».

### B-3 (baja, sigue) · `docs/estado/bloqueos/INDICE.md` no tiene ninguna fila del tablero
- Faltan `S6-tablero-sin-datos`, `S6-tablero-senal-periodo`, `S6-tablero-lineas-archivadas`, `C1-tablero`, los 5 `U3-tablero-*` y `WF3c-tablero-interfaz`.
- Le toca al integrador (`integrador-wf3c` tiene INDICE.md en su alcance).

### B-4 (baja, sigue) · Pruebas de 4xx con `c.Raw`, que no valida el esquema
- Están en `tablero_test.go` (líneas 92, 162 y 171) y en `tablero_datos_test.go` (220, 223 y 318).
- No se pierde cobertura de contrato, porque `TestContrato_Operaciones` sí valida los 4xx.

### B-5 (baja, sigue) · El texto de `BudgetSignals` en C1 no refleja la opción (a) de S6-tablero-senal-periodo
- `inicio.yaml`, línea 1576, no dice que, con un `period` posterior al mes en curso, `months` puede venir más corto o vacío.
- Tampoco cita el bloqueo.
- Corrección: una solicitud a `c1gen`.

### B-6 (baja, nueva; fuera del alcance de `api-tablero`) · El 400 de un `period` mal formado propone una quincena que la operación rechaza
- `GET /api/v1/dashboard/debt-capacity?period=2026-13` → 400: «El periodo «2026-13» no tiene un formato válido. Usa, por ejemplo, 2026-10 o 2026-10-q1.».
- El ejemplo de quincena lleva a un 422, porque la operación solo acepta un mes (`inicio.yaml`, línea 547; `tablero.go`, líneas 64-79).
- El mensaje sale del ayudante común `api.QueryPeriod` (backend-base).
- Corrección sugerida: que `monthParam` traduzca el 400 de `QueryPeriod` a un mensaje propio («…por ejemplo 2026-10»), o una solicitud a backend-base para que el ejemplo dependa del parámetro.

## 4. Lo que está bien
- Las 5 operaciones del tablero responden por HTTP/2 y validan contra `openapi.yaml`. Sus cifras salen del motor y coinciden con el oráculo, incluida la serie de 12 meses.
- Los umbrales se leen de `app_setting`: al cambiarlos cambia el color, y al restablecerlos vuelve el original. El 422 de regla está en español.
- La escucha es solo loopback y LAN; Tailscale queda cerrado; ALPN con h2 y http/1.1.
- `go test`, cobertura, contrato y `verify-all` en verde.
- El bloqueo `S6-tablero-lineas-archivadas` ya dice la verdad sobre el código, y la parte común de la prueba detecta la (a2).

## 5. Propuestas para Axel
- **PROPUESTA:** decidir en la puerta 4 la pregunta de `S6-tablero-lineas-archivadas`. El revisor coincide con la (a): archivar una tarjeta liquidada no debería borrar sus pagos pasados del barómetro y de las gráficas mientras el efectivo de S2 los sigue contando.
- **PROPUESTA (al orquestador):** la salida (i). Son unas 10 líneas en un archivo, ya probadas con overlay en tres rondas, y no mueven ninguna cifra del servidor de desarrollo, que no tiene líneas archivadas.
