# REVISIÓN · tablero (backend) · ronda 2 de 3 · WF-3c · 12-oct-2026

Revisor independiente `revisor-tablero-backend-r2`. Evalúa el trabajo de `api-tablero` y de su corrección (`docs/estado/wf3/api-tablero.md` y `correccion-tablero-backend-r1.md`) **solo** contra `docs/criterios/S6.md` (S6-01 a S6-04, S6-14 y S6-15), S2.md y S3.md (cifras del motor) y C1.md. No corrige nada. Reejecuté todas las verificaciones con comandos propios. Los temporales (scripts de curl, overlays, bitácoras) quedaron en `datos/tmp/revisor-tablero-backend-r2/` y se borraron al terminar.

**Dictamen: NO APROBADO.** Quedan 0 altos, 1 medio y 5 bajos:
- **M-1 de la ronda 1: CERRADO.**
- **M-2 de la ronda 1: SIGUE ABIERTO.** Ya tiene bloqueo y prueba, pero el código no aplica la opción que el bloqueo declara vigente.
- Todas las cifras coinciden con el oráculo.
- No hay verificaciones rotas.

## 1. Hallazgos de la ronda 1

| Hallazgo | Estado | Evidencia propia |
|---|---|---|
| **M-1** · `getBudgetSignals` contaba meses no cerrados | **CERRADO** | Detalle abajo. |
| **M-2** · los pagos de crédito del tablero dejan fuera las líneas archivadas | **ABIERTO** (sigue como M-2, media) | Ver §3. |

**Evidencia de M-1:**
- **Código.** `tablero.go`, líneas 285-306: `closed := monthIndex(mm) < current`. Sin el corte, la racha no se completa (`len(cells) == months`).
- **Servidor de desarrollo** (binario de 15:52:13, posterior a todos los fuentes; proceso arrancado a las 15:52:32; HTTP/2):
  - sin `period`: jul–sep y `none`;
  - `2026-11`: solo 2026-08 y 2026-09 (octubre ya no entra) y `none`;
  - `2027-01` y `2027-02`: `months` vacío y `none` en las 3 filas. En la ronda 1, `2027-01` daba `cut`.
- **Prueba.** `TestOp_getBudgetSignals_PeriodoFuturo` pasa. Con un overlay propio que quita el corte (`closed` siempre true), **falla** con 12 líneas de error. La prueba sí detecta la regresión.
- **Bloqueo.** `docs/estado/bloqueos/S6-tablero-senal-periodo.md` trae qué falta, 3 opciones, impacto y la pregunta para Axel. La opción (a) es la lectura literal de RN-31.

## 2. Verificación por criterio

| Criterio | Resultado | Evidencia propia |
|---|---|---|
| S6-01 | CUMPLE | Detalle 1 abajo. |
| S6-02 | CUMPLE | Detalle 2 abajo. |
| S6-03 | CUMPLE | Detalle 3 abajo. |
| S6-04 | CUMPLE con observaciones (B-1, B-2 y B-4, siguen) | Detalle 4 abajo. |
| S6-14 | CUMPLE | Detalle 5 abajo. |
| S6-15 | CUMPLE con observación (B-2) | Detalle 6 abajo. |
| S2 y S3 (cifras) | CUMPLE en el caso sembrado; con una línea archivada, ver M-2 | §2.1. |
| C1 | CUMPLE con M-2, B-1 y B-5 | Esquemas y rutas iguales al contrato. RN-31 se cumple. RN-30 «de todas las líneas» no se cumple con líneas archivadas (M-2). |
| `make verify-all GOOSE=…` | EXIT 0 | «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». |

**Detalle por criterio:**
1. **S6-01.** `make verify-wf3-go` EXIT 0:
   - go1.27.2 y go.mod 1.27;
   - build y vet sin errores; `gofmt -l` vacío; `go mod tidy -diff` sin cambios;
   - `go test ./...`: 33 paquetes OK y 5 sin pruebas.
2. **S6-02.** `make verify-wf3-cobertura` EXIT 0:
   - con umbral K25: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8;
   - sin umbral: api 82.0, auth 83.6, jobs 92.3, reportes 77.3 y telegram 80.0;
   - con `go test -cover`: api/tablero **87.8 %**, internal/api 88.8 % y api/creditos 80.4 %.
3. **S6-03.** `make verify-wf3-contrato` EXIT 0:
   - «los 22 archivos coinciden byte a byte con el generador»;
   - «221 operaciones en 176 rutas, 0 diferencias».
   - `TestContrato_Operaciones` da «221 de 221 operaciones cumplen (2xx declarada, 4xx válidas contra openapi.yaml)». Pasan las 5 operaciones del tablero y getDashboardMonthVsPlan.
4. **S6-04.** Las respuestas reales siguen la convención de C1:
   - mes `AAAA-MM` y centavos enteros;
   - `ErrorBody` con `code`, `message`, `fields` y `request_id`;
   - 400 de formato y 422 de regla.
   - `months: []` cumple el esquema `BudgetSignals`.
5. **S6-14.** `ss -ltnp` muestra solo `127.0.0.1:8443` y `192.168.50.56:8443`. Por la IP de Tailscale (100.115.28.35), curl sale con código 7, conexión rechazada. Pasan `TestEscucha_NegociaHTTP2`, `TestEscucha_SoloLANYLocal` y `TestEscucha_HTTPSConCertificado`. El HTTP/2 está en el punto (c) de §2.1.
6. **S6-15.**
   - Los identificadores nuevos están en inglés (`closed`, `current`, `alpnProtos`).
   - Los mensajes están en español de México: «El umbral amarillo (45 %) tiene que ser menor que el tope (2 %). Baja el amarillo o sube el tope.» y «Falta «from»: un mes con el formato AAAA-MM, por ejemplo 2026-10.».
   - `datos/dev/run/finanzas.log` tiene 0 coincidencias de password, contraseña, cookie o token, incluidas mi sesión y mi logout.

### 2.1 Cifras contra el oráculo y los motores (servidor de desarrollo, HTTP/2)

Usé curl con la CA de mkcert y una sesión propia de axel. La contraseña de desarrollo se leyó del archivo y se pasó por la entrada estándar, sin imprimirse. Al final cerré la sesión: logout 204 y luego 401. Recalculé las cifras desde `oraculo/model.json`, agrupando `despensas[].pagos` por mes calendario. Los porcentajes van con ROUND_HALF_UP sobre el ingreso de 2,450,000.

| # | Valor | API | Oráculo (recalculado) | ¿Coincide? |
|---|---|---|---|---|
| 1 | Deuda viva hoy (`getCreditTotals.live_debt_cents`, motor S3) | 1,899,922 | Σ de los pagos de las 4 despensas con fecha ≥ 12-oct = 18,999.22 | Sí |
| 2 | Pagos del mes del hogar y con Karen, oct-2026 (`debt-capacity`) | 59,882 · 97,375 (1 pago) | 598.82 · 973.75 (pago del 30-oct). Además es igual a `getCreditTotals.month_household_cents` y `month_total_cents` | Sí |
| 3 | Pagos del hogar y con Karen, nov y dic (`debt-capacity`) | 119,764 · 194,750 · 148,257 · 223,243 (2 pagos cada mes) | 1,197.64 · 1,947.50 · 1,482.57 · 2,232.43 | Sí |
| 4 | Serie de 12 meses (`monthly-overview` jul-2026 a jun-2027): `credit_payments_cents`, `credit_payments_total_cents` y `debt_ratio_pct` | jul–sep 0; oct a jun: 59,882 · 119,764 · 148,257 · 176,750 · 176,752 · 176,750 · 176,752 · 176,750 · 88,377 | Los 12 meses iguales al oráculo, también la parte con Karen | Sí (12 de 12) |
| 5 | Ratio y margen, oct, nov y dic | `ratio_pct` 2 · 5 · 6; `ratio_total_pct` 4 · 8 · 9; `headroom_cents` 920,118 · 860,236 · 831,743 | 59,882 ÷ 2,450,000 = 2.44 → 2; 4.89 → 5; 6.05 → 6. Margen 980,000 − pagos | Sí |
| 6 | Devengado, efectivo y plan de oct (M-21, S2) | 119,764 · 59,882 · 200,000; `month-vs-plan` 60 % | `periodos[0].devengado.axel` 1,197.64 y `efectivo.axel` 598.82; 119,764 ÷ 200,000 = 59.88 → 60 | Sí |

**(c) HTTP/2.**
- `%{http_version}` = 2 por loopback y por la IP de la LAN (`--resolve finanzas.local:8443:192.168.50.56`).
- `openssl s_client -alpn h2,http/1.1` da «ALPN protocol: h2».
- Con `--http1.1` responde en HTTP/1.1.
- El cambio en `tls.go` es mínimo: `NextProtos: alpnProtos` con `h2` y `http/1.1`.

**(d) Umbrales desde `app_setting`.**
- Con `PATCH` amarillo 1 y tope 2 → 200:
  - octubre pasa a `warning`, con escala 10 y margen −10,882;
  - noviembre pasa a `exceeded`, con margen −70,764;
  - `monthly-overview` da `warning`, `exceeded` y `exceeded`.
- Amarillo 45 con tope 2 → 422 `out_of_range` en español.
- `"4"` → 400 (B-1).
- `reset_parameters` → 200, con los 5 parámetros en `is_initial` true. Octubre vuelve a `normal`, 30/40, escala 60 y margen 920,118.
- La base de desarrollo quedó con sus valores iniciales.

**(e) Sin cifras a mano.**
- Los pagos salen de `motor_cred.UnifiedCalendar` (por `creditos.PaymentsByMonth`).
- El ingreso, el devengado, el efectivo, el plan y el % usado salen de `motor_fin`.
- El manejador solo divide con `dominio.PercentOf` y `MulDivRound` (ROUND_HALF_UP).
- La única resta es la del margen (`tablero.go`, línea 164): tope × ingreso ÷ 100 − pagos, como define C1.

## 3. Hallazgos

### M-2 (media, sigue de la ronda 1) · El código del tablero no aplica la opción que su bloqueo declara vigente: las líneas archivadas siguen fuera
- **Dónde:**
  - `backend/internal/api/creditos/tablero.go`, líneas 26-50: `PaymentsByMonth` usa `s.states(nil)`, que filtra `c.IsActive` (`motor.go`, líneas 204-217), y devuelve `true` fijo;
  - `backend/internal/api/tablero/tablero_datos_test.go`, línea 369: `t.Skip`;
  - `docs/estado/bloqueos/S6-tablero-lineas-archivadas.md`.
- **Qué se atendió:**
  - el supuesto ya está en un bloqueo con 3 puntos, opciones, impacto y pregunta;
  - hay una prueba que compara contra el efectivo de S2;
  - la solicitud fuera de alcance trae el cambio exacto.
- **Por qué sigue abierto:**
  - El bloqueo dice «Mientras tanto rige la opción PROVISIONAL (a)», pero el servidor corre la (b).
  - La prueba que demuestra la regla está en SKIP, así que `go test` pasa sin comprobarla.
  - RN-30 dice «de todas las líneas». El código de hoy cuenta solo las activas.
  - Es una inconsistencia entre artefactos (bloqueo, RN-30 y código) que cambia cifras visibles en G-1, G-6 y el barómetro.
- **Comprobado con overlays propios** (`go test -overlay`, sin tocar el árbol):
  - **Código actual, sin el Skip: FALLA.** «octubre con la línea archivada: pagos 0 (antes 650000; parte de crédito del efectivo S2 650000), total 0, 0 %, normal», y `debt-capacity` de octubre y noviembre en 0.
  - **Con el cambio de la solicitud 1, sin el Skip: PASA.** `api/tablero`, `api/creditos` e `internal/api` siguen OK.
- **Corrección sugerida:**
  - Aplicar la solicitud 1 de `correccion-tablero-backend-r1.md` tal cual, solo en `PaymentsByMonth` de `backend/internal/api/creditos/tablero.go`:
    - recorrer `s.credits` con `s.state(c)` para todas las líneas;
    - `hasCredits = hasCredits || c.IsActive`;
    - después del calendario, `hasCredits = hasCredits || len(cal.Entries) > 0`;
    - devolver `hasCredits` en lugar de `true`, y actualizar el comentario.
  - Quitar la línea 369 (`t.Skip`) de `tablero_datos_test.go`.
  - Aplicar también la solicitud 2: el texto de `has_credits` por `c1gen`.
  - **Atención, orquestador:** la corrección de la ronda 3 solo puede tocar los archivos de `api-tablero`, y `api/creditos/` no está entre ellos. Hace falta ampliar ese alcance a este único archivo (y al generador para la solicitud 2), o que lo aplique quien tenga `api/creditos`.
  - Si no se amplía, la otra salida es que los artefactos digan lo mismo que el código: el bloqueo declara que hasta la puerta 4 rige (b), la prueba sin Skip comprueba (b) y se pide corregir RN-30. Eso cambia la opción PROVISIONAL propuesta, así que lo decide el orquestador, no el revisor.

### B-1 (baja, sigue) · 400 en la API, 422 en el texto de C1
- `"debt.capacity_warning_pct":"4"` responde 400 `bad_request`, pero `ajustes-e.yaml` dice «de otro tipo responde 422».
- Nuevo caso de la misma familia: `GET /api/v1/dashboard/monthly-overview` sin `from` responde 400 («Falta «from»…»). La descripción de la operación (`inicio.yaml`, línea 586) dice «`from` y `to` obligatorios…; si no, 422 validation_failed». En cambio, el parámetro `OverviewFrom` (`componentes-comunes.yaml`, línea 149) solo pide 422 para el orden o para más de 24 meses.
- Hay que alinear el texto en `c1gen` (`gen3_ops.py`).

### B-2 (baja, sigue) · Mensaje 400 genérico
El 400 de tipo incorrecto en `updateDashboardSettings` sigue diciendo «Tipo incorrecto: los importes son enteros en centavos y las fechas, texto AAAA-MM-DD.», cuando el campo es un porcentaje.

### B-3 (baja, sigue y crece) · `docs/estado/bloqueos/INDICE.md` no tiene ninguna fila del tablero
Faltan 10 bloqueos: `S6-tablero-sin-datos`, `S6-tablero-senal-periodo`, `S6-tablero-lineas-archivadas`, `C1-tablero`, los 5 `U3-tablero-*` y `WF3c-tablero-interfaz`. La solicitud 3 de la corrección lo pide y el alcance del integrador incluye INDICE.md.

### B-4 (baja, sigue) · Pruebas de 4xx con `c.Raw`, que no valida el esquema
`tablero_test.go`, líneas 92, 162 y 171, y `tablero_datos_test.go`, líneas 220, 223 y 318. La cobertura de contrato no se pierde, porque `TestContrato_Operaciones` sí valida los 4xx.

### B-5 (baja, nueva) · El texto de `BudgetSignals` en C1 no refleja la opción (a) de S6-tablero-senal-periodo
- `inicio.yaml`, línea 1576: «Mira los months_considered meses cerrados anteriores al mes pedido», con cita solo a `U3-tablero-senal-presupuesto` y `C1-tablero`.
- No dice que, con un `period` posterior al mes en curso, `months` puede venir con menos meses que `months_considered` o vacío, ni cita el bloqueo nuevo.
- No rompe el esquema, porque `months: []` es válido.
- Corrección sugerida: solicitud a `c1gen` para agregar «con un period posterior al mes en curso solo cuentan los meses cerrados; months puede venir más corto o vacío (PROVISIONAL, bloqueo S6-tablero-senal-periodo)».

## 4. Lo que está bien
- M-1 quedó resuelto en un solo lugar, con una prueba que falla sin el corte y un bloqueo completo. Con el `period` por omisión no cambia nada.
- M-2 está bien diagnosticado: el bloqueo explica la contradicción interna de RN-30 («de todas las líneas» frente a «la misma cifra que `getCreditTotals`»). La solicitud fuera de alcance es exacta y la verifiqué: con ella todo pasa.
- Las cifras del motor coinciden con el oráculo en 12 de 12 meses de la serie, en la deuda viva, en los ratios y en el margen. Los umbrales se leen de `app_setting` en cada petición.
- El HTTP/2 se negocia por ALPN y la escucha sigue solo en LAN y loopback.

## 5. Para la ronda 3
- Cerrar M-2: aplicar la solicitud 1 y quitar el Skip. Requiere ampliar el alcance (ver arriba).
- B-1 a B-5 son opcionales.
