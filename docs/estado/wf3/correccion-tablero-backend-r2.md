# correccion-tablero-backend-r2 · M-2 de la revisión del tablero (backend), ronda 2 · WF-3c · 12-oct-2026

Agente `correccion-tablero-backend-r2` (Opus). El encargo trae un solo hallazgo: el M-2 (media) de `docs/estado/wf3/REVISION-tablero-backend-r2.md`, sin altos. Los bajos B-1 a B-5 no eran parte del encargo y no se tocaron.

Alcance: los archivos de `api-tablero` que nombra el hallazgo. **`backend/internal/api/creditos/tablero.go` no es de ese alcance**:
- el paquete `internal/api/creditos` es de `api-creditos` (CONVENCIONES-backend §2);
- el ALCANCE de `api-tablero` en `runs/wf3c-tablero.js` no lo incluye;
- el archivo lo escribió `correccion-especificacion-r2`.

El revisor lo dice en el mismo hallazgo («ATENCIÓN ORQUESTADOR») y deja la decisión al orquestador. El encargo no amplió el alcance. Por eso **no apliqué** el cambio en `api/creditos` ni elegí la salida (ii), que cambia la opción PROVISIONAL. Corregí lo que sí es del alcance y dejé la decisión escrita como PREGUNTA, con el cambio exacto y probado.

## Hecho
- **HECHO · El bloqueo ya no dice algo falso.** En `docs/estado/bloqueos/S6-tablero-lineas-archivadas.md`, la frase «Mientras tanto rige la opción PROVISIONAL (a)» cambió por:
  - «La propuesta PROVISIONAL es la (a) con la (a1). **Todavía no rige en el código**: hoy el servidor hace lo de la (b). Eso no es una elección: falta un cambio fuera del alcance de `api-tablero`.»
  - Con esto, el bloqueo, el reporte de la ronda 1 y el código cuentan lo mismo sobre lo que pasa hoy.
- **HECHO · Punto 4 nuevo en el bloqueo**, «Decisión pendiente del orquestador: el alcance (no la regla)». Tiene una tabla con las dos salidas del revisor:
  - **(i)** ampliar el alcance o que lo aplique el dueño de `api/creditos`. Con esta, rige la (a) como se propone;
  - **(ii)** igualar los artefactos a la (b). Cambia la opción propuesta.
  - Para cada una dice quién la aplica y qué efecto tiene. PROPUESTA: la (i).
- **HECHO · La prueba ya no queda entera en SKIP.** `TestOp_Dashboard_LineaArchivada` (`backend/internal/api/tablero/tablero_datos_test.go`, líneas 365-427; el `t.Skip` está en la 408) ahora tiene dos partes.
  - **Parte que corre siempre.** Comprueba lo que vale con cualquier opción del bloqueo:
    - el efectivo de S2 de octubre sigue en 700,000 después de archivar;
    - G-1/G-6 (`monthly-overview`) y el barómetro (`debt-capacity`) dan la misma cifra de octubre;
    - en octubre, noviembre y diciembre, con `has_credits` false, los pagos, el conteo y `ratio_pct` son 0 (texto de C1). Esto descarta la (a2).
  - **Subtest `opcion_a_PROVISIONAL`.** Tiene las aserciones de la (a) y la (a1), con las mismas cifras que la ronda 1. Queda en SKIP con la referencia al punto 4 del bloqueo. Al aplicar la solicitud 1, basta con quitar ese `t.Skip`.
  - `go test -v` lo muestra explícito: `--- PASS: TestOp_Dashboard_LineaArchivada` y `--- SKIP: …/opcion_a_PROVISIONAL`.
- **HECHO · Solicitud 1 comprobada otra vez**, con `go test -overlay` en `datos/tmp/`, sin tocar el árbol:

| Código | Prueba | Resultado |
|---|---|---|
| actual, que hace lo de la (b) | completa, sin SKIP | **FALLA**: «octubre con la línea archivada: pagos 0 (antes 650000; parte de crédito del efectivo S2 650000), total 0, 0 %, normal»; `debt-capacity` de octubre y noviembre en 0 |
| con la solicitud 1, (a) con (a1) | completa, sin SKIP | **PASA**. `api/tablero`, `api/creditos` e `internal/api` OK; `go vet` sin errores |
| (a) con (a2): todas las líneas, `has_credits` solo por activas | completa, sin SKIP | **FALLA** en la parte que corre siempre: «has_credits false con pagos o porcentaje» en octubre y noviembre |

  Así se ve que la parte que corre siempre sí detecta un error y no frena la opción propuesta.
- **HECHO · El cambio no mueve cifras del servidor de desarrollo.** La base `datos/finanzas-dev.db` tiene 1 línea de crédito, activa (consulta de solo lectura `select is_active, count(*) from credit`). Las cifras que el revisor comparó con el oráculo quedan igual con la (a).
- **PREGUNTA (para el orquestador):** ¿se amplía el alcance a `backend/internal/api/creditos/tablero.go` y `c1gen` (salida i), o se igualan los artefactos a la (b) (salida ii)? Mientras no se decida, **M-2 no se puede cerrar en el código desde `api-tablero`**. Hoy el servidor cuenta solo las líneas activas.

## Archivos
- Cambiados:
  - `backend/internal/api/tablero/tablero_datos_test.go`: solo `TestOp_Dashboard_LineaArchivada`.
  - `docs/estado/bloqueos/S6-tablero-lineas-archivadas.md`: encabezado, «Estado del código» y punto 4 nuevo. Los puntos 1 a 3, el impacto y la pregunta para Axel no cambian.
- Nuevo: este reporte.
- Sin cambio:
  - `backend/internal/api/creditos/tablero.go` (solicitud 1);
  - `backend/internal/api/tablero/tablero.go`: el comentario de la línea 148 ya cita el bloqueo;
  - `c1gen`, C1, reglas de negocio y migraciones.

## Verificación
- **S6-01:** `make verify-wf3-go` EXIT 0:
  - build y vet sin errores; `gofmt -l .` vacío; `go mod tidy -diff` sin cambios;
  - `go test ./...`: 33 paquetes OK y 5 sin pruebas.
- **S6-02:**
  - `make verify-wf3-cobertura` EXIT 0. K25 OK: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8. Sin umbral: api 82.0, auth 83.6, jobs 93.0, reportes 77.3 y telegram 80.0;
  - `go test -cover`: api/tablero **87.8 %** (igual que antes), api/creditos 80.4 % e internal/api 88.5 %.
- **S6-03:** `make verify-wf3-contrato` EXIT 0:
  - «los 22 archivos coinciden byte a byte con el generador»;
  - «221 operaciones en 176 rutas, 0 diferencias»;
  - `TestContrato_Operaciones`: «221 de 221 operaciones cumplen (2xx declarada, 4xx válidas contra openapi.yaml)».
  - En api/tablero pasan las 11 pruebas `TestOp_`; el único SKIP es el subtest `opcion_a_PROVISIONAL`.
- **S6-04:** la prueba usa `c.Do`, que valida cada respuesta contra `openapi.yaml` (monthly-overview y 3 llamadas a debt-capacity).
- **S6-14:** no cambió el servidor, así que no hubo que reiniciarlo:
  - `ss -ltnp` muestra solo `127.0.0.1:8443` y `192.168.50.56:8443`;
  - `GET /health` responde 200 por HTTP/2 por loopback y por la IP de la LAN (`curl -w %{http_version}` = 2).
- **S6-15:**
  - identificadores en inglés (`byPeriod`, subtest `opcion_a_PROVISIONAL`, que es el nombre de la opción del bloqueo);
  - mensajes y comentarios en español de México;
  - `datos/dev/run/finanzas.log` tiene 0 coincidencias de password, contraseña, cookie o token.
- **S2 y S3 (cifras):** ninguna cambia. No toqué ni el motor ni el código que corre en el servidor.
- **`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`:** EXIT 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)».

## Supuestos
- «Los archivos de api-tablero» = el ALCANCE de `api-tablero` en `runs/wf3c-tablero.js` más lo que escribió su ronda 1 de corrección (el bloqueo). `api/creditos/` y `c1gen` quedan fuera. Es la misma lectura de la ronda 1 y del revisor.
- Elegir entre la salida (i) y la (ii) no le toca a un agente de corrección: la (ii) cambia la opción PROVISIONAL propuesta, y el revisor la deja al orquestador.

## Bloqueos
- **Actualizado:** `S6-tablero-lineas-archivadas`, con el estado real del código y el punto 4. Pregunta para Axel (puerta 4): ¿la (a) o la (b)? Pregunta para el orquestador: ¿la salida (i) o la (ii)?
- **Siguen vigentes:** `S6-tablero-senal-periodo`, `S6-tablero-sin-datos`, `C1-tablero` y `U3-tablero-*`.

## Solicitudes fuera de alcance
1. **`backend/internal/api/creditos/tablero.go`** (dueño: `api-creditos`, o el orquestador si amplía el alcance), `PaymentsByMonth`, para aplicar la (a) y la (a1). Este es el diff exacto; con él, `go vet` y las pruebas de `api`, `api/creditos` y `api/tablero` pasan:
   ```diff
   -// cuenta que getCreditTotals (month_household_cents y month_total_cents): calendario unificado del motor S3 con las
   -// líneas activas; parte del hogar = partes de las personas que no son externas. Aquí no se recalcula nada: solo se
   +// cuenta que getCreditTotals (month_household_cents y month_total_cents), pero con todas las líneas, activas o
   +// archivadas (RN-30 «de todas las líneas», igual que el efectivo de S2; PROVISIONAL, bloqueo
   +// S6-tablero-lineas-archivadas, opción a): calendario unificado del motor S3; parte del hogar = partes de las
   +// personas que no son externas. Aquí no se recalcula nada: solo se
   @@ PaymentsByMonth
   -// hasCredits dice si hay al menos una línea de crédito activa.
   +// hasCredits dice si hay al menos una línea de crédito activa o algún pago (de una línea archivada) con fecha entre
   +// from y to; con false no hay pagos (PROVISIONAL, bloqueo S6-tablero-lineas-archivadas, opción a1).
   -	states := s.states(nil)
   +	var states []mc.CreditState
   +	for _, c := range s.credits {
   +		states = append(states, s.state(c))
   +		hasCredits = hasCredits || c.IsActive
   +	}
   @@
   -		return nil, true, motorErr(err, "period")
   +		return nil, hasCredits, motorErr(err, "period")
    	}
   +	hasCredits = hasCredits || len(cal.Entries) > 0
   @@
   -	return out, true, nil
   +	return out, hasCredits, nil
   ```
   Luego, quitar la línea `t.Skip("PENDIENTE: requiere la solicitud fuera de alcance…")` del subtest `opcion_a_PROVISIONAL` de `TestOp_Dashboard_LineaArchivada`, en `backend/internal/api/tablero/tablero_datos_test.go`. `motor.go` (`states`) y `getCreditTotals` no cambian.
2. **`backend/internal/api/c1gen/gen2_schemas.py`, línea 529** (dueño: `s1-c1-v15`), texto de `DashboardDebtCapacity.has_credits`:
   - hoy dice `'b:true si hay al menos una línea de crédito activa; con false, payments_cents = 0 y ratio_pct = 0 (cifra real, no falta de dato).'`;
   - debe decir `'b:true si hay al menos una línea de crédito activa o algún pago de una línea archivada con fecha en el mes (PROVISIONAL, bloqueo S6-tablero-lineas-archivadas); con false, payments_cents = 0 y ratio_pct = 0 (cifra real, no falta de dato).'`;
   - después, regenerar con `build.py` (openapi.yaml, extractos, C1-datos.md y `types_gen.go`) y correr `make verify-wf3-contrato`.
3. **`docs/estado/bloqueos/INDICE.md`** (dueño: `integrador-wf3c`): agregar a la tabla de S6 las dos filas que ya pidió la solicitud 3 de `correccion-tablero-backend-r1.md` (`S6-tablero-senal-periodo` y `S6-tablero-lineas-archivadas`). Completa B-3.

## Pendientes
- **M-2 sigue abierto en el código** hasta que el orquestador elija la salida (i) o la (ii) del punto 4 del bloqueo. Con la (i), aplicar las solicitudes 1 y 2 y quitar el SKIP del subtest. Con la (ii), reescribir el subtest para comprobar la (b) y pedir que se corrijan RN-30 y M-20/M-21.
- B-1 a B-5 de la revisión: no eran parte del encargo.

## Propuestas para Axel
- **PROPUESTA (S6-tablero-lineas-archivadas):** dejar la (a). Archivar una tarjeta que ya terminaste de pagar no debería borrar sus pagos del barómetro ni de las gráficas de meses pasados, porque el efectivo del tablero sí los sigue contando.
- **PROPUESTA (al orquestador):** la salida (i). Son unas 10 líneas, ya probadas con overlay en dos rondas, y no cambian ninguna cifra del motor ni del servidor de desarrollo.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 27 acciones, unos 120 KB leídos y unos 18 KB escritos. Los temporales de `datos/tmp/correccion-tablero-backend-r2/` (overlays, copias y bitácoras de make) se borraron.
