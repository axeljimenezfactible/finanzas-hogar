# revisor-tablero-backend-r2 · Revisión independiente del tablero (backend), ronda 2 de 3 · WF-3c · 12-oct-2026

Agente revisor; no corrige. El informe completo está en `docs/estado/wf3/REVISION-tablero-backend-r2.md`.

## Hecho
- **HECHO · dictamen: NO APROBADO.** Hay 0 altos, 1 medio y 5 bajos.
- **HECHO · M-1 de la ronda 1 CERRADO.**
  - El corte por mes cerrado funciona en el servidor de desarrollo: `2026-11` mira ago y sep; `2027-01` y `2027-02` dan `months` vacío y `none`.
  - La prueba nueva falla si se quita el corte.
- **HECHO · M-2 de la ronda 1 SIGUE ABIERTO (media).**
  - El bloqueo y la prueba existen.
  - El código sigue contando solo las líneas activas, contra la opción (a) que el bloqueo declara vigente y contra RN-30.
  - La prueba está en SKIP.
  - Lo comprobé con overlays: el código actual falla; con la solicitud 1 pasa.
- **HECHO · cifras.** Recalculé desde `oraculo/model.json`: deuda viva, pagos del mes (oct, nov y dic, del hogar y con Karen), la serie de 12 meses, los ratios con ROUND_HALF_UP, el margen y el devengado y efectivo de octubre. Todo coincide.
- **HECHO · bajas.** B-1 a B-4 siguen abiertas. B-5 es nueva: el texto de `BudgetSignals` en C1 no refleja la opción (a).

## Archivos
- Nuevos:
  - `docs/estado/wf3/REVISION-tablero-backend-r2.md`;
  - este reporte.
- Temporales en `datos/tmp/revisor-tablero-backend-r2/` (scripts de curl, overlays y bitácoras): **borrados**.

## Verificación
- **S6-01:** `make verify-wf3-go` EXIT 0:
  - go1.27.2 y go.mod 1.27;
  - build, vet y gofmt limpios; `go mod tidy -diff` sin cambios;
  - 33 paquetes OK y 5 sin pruebas.
- **S6-02:** `make verify-wf3-cobertura` EXIT 0.
  - K25: dominio 87.0, motor_fin 89.2, motor_cred 94.0, recetas 93.2 y precios 88.8.
  - Sin umbral: api 82.0, auth 83.6, jobs 92.3, reportes 77.3 y telegram 80.0.
  - api/tablero 87.8 %.
- **S6-03:**
  - `make verify-wf3-contrato` EXIT 0: 22 archivos byte a byte, y «221 operaciones en 176 rutas, 0 diferencias»;
  - `TestContrato_Operaciones`: «221 de 221 operaciones cumplen».
- **S6-04:**
  - las respuestas siguen la convención de C1 y `months: []` cumple el esquema;
  - B-1, B-2 y B-4 siguen (bajas).
- **S6-14:**
  - `ss -ltnp` muestra solo `127.0.0.1:8443` y `192.168.50.56:8443`;
  - por Tailscale, curl sale con código 7;
  - curl da `%{http_version}` = 2 por loopback y por la LAN; ALPN h2;
  - con `--http1.1` responde en 1.1;
  - las 3 pruebas `TestEscucha_*` pasan.
- **S6-15:**
  - identificadores en inglés y mensajes en español de México;
  - `finanzas.log` tiene 0 coincidencias de password, contraseña, cookie o token.
- **S2 y S3:**
  - 12 de 12 meses de la serie iguales al oráculo;
  - deuda viva 1,899,922 = Σ de los pagos pendientes del oráculo;
  - ratios 2 · 5 · 6 y margen 920,118 · 860,236 · 831,743.
- **(d) umbrales:**
  - con 1/2: `warning`, `exceeded` y `exceeded`;
  - 45 con tope 2 da 422;
  - el reset vuelve a 30/40 y `normal`.
- **M-1:**
  - en vivo, `2026-11`, `2027-01` y `2027-02` dan `none` sin meses abiertos;
  - con un overlay sin el corte, `TestOp_getBudgetSignals_PeriodoFuturo` FALLA (12 líneas de error).
- **M-2** (overlays):
  - sin el Skip y con el código actual, FALLA: «pagos 0 (antes 650000…)»;
  - con la solicitud 1, PASA, junto con `api/tablero`, `api/creditos` e `internal/api`.
- **`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`:** EXIT 0, «Verificación ampliada OK».
- Sesión de desarrollo cerrada: logout 204 y luego 401.

## Supuestos
- «Deuda viva» se verificó con `getCreditTotals.live_debt_cents`, porque ninguna de las 5 operaciones del tablero la expone. Usa el mismo motor S3 y la misma base.
- El oráculo por mes calendario se obtuvo agrupando `despensas[].pagos` de `model.json`. Coincide con los `periodos` P1 a P3, que usó la ronda 1.

## Bloqueos
Ninguno nuevo. Revisé `S6-tablero-senal-periodo` y `S6-tablero-lineas-archivadas`.

## Pendientes
- **M-2.** Aplicar la solicitud 1 de `correccion-tablero-backend-r1.md` en `backend/internal/api/creditos/tablero.go` y quitar `t.Skip` (`tablero_datos_test.go`, línea 369). Queda fuera del alcance de la corrección de la ronda 3: hace falta que el orquestador lo amplíe.
- B-1 a B-5 (bajas, opcionales).

## Propuestas para Axel
- **PROPUESTA:** para S6-tablero-lineas-archivadas, aceptar la (a). Una tarjeta liquidada y archivada no debería borrar sus pagos pasados del barómetro, porque el efectivo del tablero sí los cuenta.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 42 acciones, unos 110 KB leídos y unos 20 KB escritos.
