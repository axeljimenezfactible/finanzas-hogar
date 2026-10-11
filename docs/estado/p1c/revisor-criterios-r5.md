# revisor-criterios-r5 · Revisión independiente de los criterios 1.2, ronda 5 (flujo P1e) · 11-oct-2026

## Hecho
- Revisé la ronda 5 de los criterios 1.2 solo contra `docs/criterios/INDICE.md` §1 y `docs/plan/decisiones.md` §K (K35 a K44). No corregí nada.
- Comprobé los 6 hallazgos de `REVISION-criterios-r4.md`: los 6 están **cerrados** (H-1 alta; H-2 y H-3 medias; H-4, H-5 y H-6 bajas).
- Reejecuté las comprobaciones (a) a (j) con scripts propios. Base temporal con goose 00001→00003, Verifica de C1-22 y parte ejecutable de U4-12 contra la base y `ajustes-*.yaml`: 0 fallas. Verifica literal de C1-12 sobre `openapi.yaml`: 0. `verif_final.sh` y `make verify-criterios`: rc = 0. 0 marcas PROVISIONAL activas que citen un bloqueo resuelto. Cifras presentes y recalculadas. El diff de la ronda 5 queda dentro de lo listado.
- **Veredicto: APROBADO** (0 altos, 0 medios, 2 bajos: B-1 en S7-10 y B-2 en U4 §4/§5). El informe completo está en `docs/estado/p1c/REVISION-criterios-r5.md`.

## Archivos
- Creado: `docs/estado/p1c/REVISION-criterios-r5.md` (informe).
- Creado: `docs/estado/p1c/revisor-criterios-r5.md` (este reporte).
- Temporales: `datos/tmp/rev5c/` (`t.db`, `chk.py`, `c122.py`, `c112.py`, `prov.py`, muestras de regex y salidas). Borrados al terminar.

## Verificación
- `bash docs/criterios/tools/verif_final.sh`, rc = 0: 261 (201 / 60 / 0); citas 306 / 0; refs 261/261; canon 178 / 0; S1-09 SQL OK (54 tablas) y equivalencia 12 / 0; matriz 210 operaciones, 0 N-xxx sin operación; rutas `/tmp` 0; INDICE contra conteo 0.
- `make verify-criterios GOOSE=<scratchpad>/gobin/goose`, rc = 0: «las 11 comprobaciones dan lo esperado».
- `goose -dir backend/db/migrations sqlite3 datos/tmp/rev5c/t.db up`: versión 3.
- `chk.py`: 18 archivos, 261 criterios (201 / 60 / 0), 0 repetidos y 0 problemas de estructura o §5; §2 y §10 coinciden; cifras 219.88 ×7, 629.35 ×7, 0.9836 ×6, 0.4918 ×6.
- `c122.py`: 15 `@param` sembrados con valor y cita; lectura `value`/`initial_value` y escritura tipada; 0 `x-provisional`; solo `minimum: 0`; alertas en `getAlertSettings`/`updateAlertSettings`; los 4 parámetros de avisos como `@param_sin_valor`; `weekday` 7/7 y `day_of_month` 1 en `params_json`; 0 `thresholds_pct`/`days_before`; U4-12, 12 claves con valor y sección: 0 diferencias. Total: 0 fallas.
- `uv run --offline --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/rev5c/t.db`: 29 OK, TODO OK (con `python3` sin `jsonschema`, solo C1-14 falla por la dependencia).
- `c112.py`: operaciones públicas `getHealth` y `login`; 0 propiedades `password|hash|token|secret` en respuestas sin `x-app-setting`; excluida solo `auth.password_min_length`.
- `prov.py`: 43 bloqueos, 41 con «Respuesta de Axel»; 89 líneas con PROVISIONAL; las 10 que nombran un bloqueo resuelto son avisos de retiro o Verifica de ausencia. 0 activas.
- Regex de P3-03: 6 advertencias dan 0; 6 credenciales dan 6; «Paso clave: exporta…» da 1, que el criterio acepta con justificación.
- `git diff HEAD -U0 -- docs/criterios/` y fechas de modificación: la ronda 5 tocó solo C1, P3, S7, U4, INDICE y U3, en las partes listadas.

## Supuestos
- En la Verifica de C1-12, «busca `password`, `hash`, `token`, `secret` en esquemas de respuesta» lo leo sobre **nombres de propiedad**, con las refs resueltas, igual que la ronda 4. Las descripciones en prosa, como «sin token», no cuentan.
- La ronda 5 la identifico por fecha de modificación: archivos escritos después de `REVISION-criterios-r4.md` (16:56:47). Los bloques de la ronda 4 ya los revisó la ronda 4.

## Bloqueos
Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md` (sin «Respuesta de Axel»); no lo toqué.

## Pendientes
- Bajos opcionales: B-1 (S7-10, probar el cambio de `day_of_month` en la Verifica) y B-2 (U4 §4 y §5, valores de la decisión F de U4-12). No bloquean la aprobación.

## Solicitudes fuera de alcance
- Para el integrador de P1e: confirmar qué agente modificó a las 16:59:04-05 `00003_puerta1b.sql`, `S1-parametros-avisos.md`, `openapi.yaml`, los extractos, `C1-datos.md`, `diccionario-datos.md` y `routes.go`. Fue después de `REVISION-C1-r3` y de `REVISION-criterios-r4`, y antes de la ronda 5, que no lo declara. Mis comprobaciones sobre el estado actual pasan; la pregunta es solo de trazabilidad antes del commit.

## Propuestas para Axel
Ninguna.

## Consumo aproximado
Unas 41 acciones de herramienta. Leído ≈ 150 KB: protocolo, bloqueo, sección K, informe de la ronda 4, reporte de la ronda 5, tramos de los criterios, salidas de scripts y del diff. Escrito ≈ 19 KB: informe, reporte y scripts temporales.
