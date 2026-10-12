# revisor-backend-core-r2 · reporte del agente

12-oct-2026 · WF-3a · revisor independiente, ronda 2 de 3, de backend-core: backend-base, auth-sistema, semilla, api-gastos-inicio, api-presupuesto y api-comidas-salud. Solo leí, ejecuté y escribí este reporte y el informe. No corregí código, no hice git add ni commit, no envié nada por Telegram y no usé ni escribí credenciales reales: las contraseñas de la semilla fueron valores falsos de prueba y nunca salieron de `datos/tmp`.

## Hecho
- HECHO: informe `docs/estado/wf3/REVISION-backend-core-r2.md`. Veredicto: **APROBADO**, con 0 altos, 0 medios y 4 bajos (B-1 a B-4).
- HECHO: volví a probar uno por uno los 9 hallazgos de la ronda 1 y los 9 están cerrados.
  - H2 (`getDashboardBestBasket`) se cerró **mientras yo revisaba**, por un cambio concurrente de otro agente entre las 04:53 y las 04:59. Repetí las verificaciones sobre el estado final.
- HECHO: para S6-01 a S6-10, S6-13, S6-15, S6-16, C3-03, C3-04 y C3-09 corrí mis propios comandos, con el binario real, una base sembrada en `datos/tmp` y un servidor `--dev-http` en 127.0.0.1:18943.
- HECHO: probé 5 operaciones al azar por área (semilla 20261012), más algunas escrituras, y validé cada respuesta contra `openapi.yaml` con jsonschema 2020-12.

## Archivos
- `docs/estado/wf3/REVISION-backend-core-r2.md` (nuevo).
- `docs/estado/wf3/revisor-backend-core-r2.md` (este).
- Temporales en `datos/tmp/rbc2/`: binario, bases, logs y scripts de prueba en Python. Los borré al terminar.

## Verificación (resumen; el detalle está en el informe)
- **S6-01**: go1.27.2, `go 1.27`; build, vet y gofmt sin salida.
- **S6-02**: 0 FAIL en el estado final. Con umbral: dominio 87.0 %, motor_fin 89.2 %, motor_cred 94.0 %, recetas 93.2 % y precios 88.8 %. Sin umbral: api 82.0 %, auth 83.6 %, jobs 92.3 %, reportes 77.3 % y telegram 80.0 %.
- **S6-03**: `build.py --check` rc=0 (23 archivos idénticos). `verify_c1.py --db` «TODO OK» (216 = 216). `TestContrato_Rutas` da 0 sin implementar y `TestContrato_Operaciones` 216/216. Hay TestOp para 216 de 216 operaciones.
- **S6-04**: gasto y peso atrasados con la misma clave → 1 fila cada uno. `limit` 201 → 400. De 20 frases al azar, 18 son mensajes y están en español; 0 en inglés.
- **S6-05**: 2 cuentas `$argon2id$`, Karen 401, cookie `HttpOnly;Secure;SameSite=Strict`, contraseña de 11 → 422 y de 12 → 204. Prueba de permisos con las 10 áreas: PASS.
- **S6-06**: casos (a), (b) y (c) PASS con reloj simulado. Leí las pruebas.
- **S6-07**: PASS (90 días; el uso del día 60 lleva el vencimiento al 150; sin uso, el día 91 da 401). En vivo, al cerrar un dispositivo su cookie da 401 y el otro sigue con 200.
- **S6-08**: PASS (día 30, febrero de 2027, doble disparo y reinicio).
- **S6-09**: `finanzas respaldo` pasa `integrity_check` ok; restauración con 56 tablas y 0 diferencias.
- **S6-10**: PASS (10 MB → 201; 10 MB + 1 → 413; el límite cambiado se respeta; factor 1 << 20 PROVISIONAL).
- **S6-13**: PASS; grep = 0; IA `is_enabled: false`.
- **S6-15**: PASS; en los logs no aparecen contraseñas ni cookies.
- **S6-16**: 19 `S6-*`, incluidos los 3 `S6-semilla-*`.
- **C3-03**: 216/216 registradas en el servidor.
- **C3-04**: pruebas de backend en verde.
- **C3-09**: A 1000000 / B 690550, pago 2026-10-30 de 97375 y deuda al 2027-04-10 de 629344, igual que `model.json`.
- **S6-11**: NO VERIFICADO (no hay mkcert).
- **S6-12**: fuera de esta revisión.
- **S6-14**: parcial (WF-4).
- **Global**: `make verify-all GOOSE=…` RC=0, dos veces; la última sobre el árbol estable. `make verify-wf3-go verify-wf3-cobertura verify-wf3-contrato` RC=0.

## Supuestos
- El veredicto vale para el árbol a las 04:59 del 12-oct-2026. Los archivos del cambio concurrente no aparecían en ningún reporte de `docs/estado/wf3/` cuando revisé; lo que escriba su autor después no está evaluado aquí.
- En el lote de la semilla, la tarjeta «Mejor canasta» sale con `covered_count` 0 porque no se cargó ningún precio (bloqueo `S6-semilla-sucursales-municipio`). No lo tomé como hallazgo: C3-09 no incluye esa tarjeta.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Los 4 bajos (B-1 a B-4) quedan para cuando sus dueños los atiendan; no bloquean.
- S6-11 (mkcert y Android), S6-12 (PWA) y S6-14 (IP de la LAN en `ops/runbook.md`) quedan para el integrador, WF-4 y Axel.

## Propuestas para Axel
- PROPUESTA: cuando exista `ops/runbook.md` en WF-4, copiar ahí la hora (03:15) y la retención (30) de respaldos, que siguen PROVISIONAL, para aprobarlas en la puerta 4.

## Paquetes instalados
- Ninguno. Usé `uv run --no-project --with pyyaml --with jsonschema` solo en ejecución, igual que en `verify_c1.py`.

## Consumo aproximado
- Unas 58 acciones; alrededor de 150 KB leídos (criterios S6, reporte de corrección, fragmentos de código y pruebas, y salidas de comandos); alrededor de 32 KB escritos (los dos reportes y unos 14 KB de scripts temporales).
