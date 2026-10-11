# Reporte · revisor-wf0-r2 (revisor independiente de WF-0, ronda 2 de 3)

11-oct-2026 · Resultado: **NO APROBADO** (2 altas: R2-01 y R2-02, ambas de C1; 6 bajas). Detalle completo, tabla de criterios, cobertura de decisiones A–I y hallazgos en `docs/estado/wf0/REVISION-r2.md`. Cerrado por tope: no.

## Hecho
- Leí el protocolo, `decisiones.md`, el plan §1 y §3 (con `sed -n`), los criterios INDICE, S1, U1, U2 y C1, la revisión r1 y los cuatro reportes de corrección; revisé por partes la migración, el diccionario, las reglas, el OpenAPI, C1, U1, la colorimetría y el backend.
- Reejecuté las 9 verificaciones de la tarea con scripts propios y, aparte, los verificadores de los agentes (`verify_c1.py`, `build.py --check`, `s1_verifica.py` extraído del Anexo B de S1, `u1_verify.py`, `colorimetria.py` en árbol espejo).
- Escribí la tabla de criterios (S1 17, U1 14, C1 15, U2 14, criterios 4 aspectos), la cobertura de decisiones A–I, el estado de los 12 hallazgos de r1 y 8 hallazgos nuevos.
- No corregí nada de lo revisado.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf0/REVISION-r2.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf0/revisor-wf0-r2.md` (este reporte)
- Temporales: `/home/deepsec/finanzas-hogar/datos/tmp/` (rev.db, rev-down.db, rev-goose.db, dec.db), creada y **borrada al final**; `datos/` queda solo con `semilla/`.
- Scripts del revisor en el scratchpad (efímero): `rev-r2/matriz_r2.py`, `regresion_r2.py`, `dicc_r2.py`, `color_r2.py`, `u2_r2.py`, `decisiones_r2.py`, `provisional_r2.py`, `criterios_r2.py`, `citas_r2.py`, `reglas_r2.py`, `reglas2_r2.py`, `up.sql`, `down.sql`, `mirror/`.

## Verificación
- SQL: Up con `sqlite3 -bail` → 54 tablas STRICT, 147 índices, 68 disparadores, `foreign_key_check` 0 filas, `integrity_check` ok; Down → 0 objetos; goose v3.28.0 up/down/up ok.
- Go: `go build ./...` OK, `go vet ./...` OK, `gofmt -l .` vacío (go1.26.5; `go.mod` = 1.26, bloqueo `S6-version-go`).
- Colorimetría: regenera md, csv, json y html idénticos byte a byte; 202 pares, 0 `cumple=false`; recálculo propio de los 202 sin diferencias; 190 tokens = 190 filas del Anexo A; 0 tokens fuera de `color.critical` en 80°–120°.
- Matriz: 199 N-xxx, 193 operaciones; 4 N-xxx sin operación ni Local (N-151, N-152, N-153, N-165); 1 operación sin N-xxx justificada (`getOpenApiDocument`); 0 pantallas sin operación. `verify_c1.py --db` → 2 fallas (C1-02, C1-03). `build.py --check` → DIFIERE `openapi.yaml` y `C1-datos.md`.
- OpenAPI: `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator` → OK.
- Regresión propia: 389500 y 239528; 142 valores, 0 diferencias (17 eventos, series A/B de 14, mínimos, P1–P6 por mes calendario y por persona). `s1_verifica.py` → TODO OK (164 valores, 46 de 46 negativas).
- Criterios: 248 (189/59) = INDICE; 248 con «Verifica» y «Evidencia»; 248 trazados en §5; 27 de 288 citas no literales.
- Decisiones A–I: 37 comprobaciones de esquema, OpenAPI y U1 en OK; 5 renglones con C1 desfasado.
- F-1: 3 toques demostrados en U1 §4.

## Supuestos
- Evalué C1 tal como está publicado en el repositorio (`openapi.yaml` 316ce18d…, generado a las 23:50:01), no como lo dejaría una regeneración futura.
- Usé el binario de goose v3.28.0 que ya estaba en el scratchpad compartido (`gobin/goose`, `goose --version` = v3.28.0) y el validador desde la caché de uv, sin red.
- Ninguno cambia cifras, esquema, endpoints ni tokens.

## Bloqueos
- Ninguno nuevo. No escribí en `docs/estado/bloqueos/`.

## Pendientes
- Ronda 3: regenerar C1 contra el esquema 1.1 de S1 (R2-01, R2-02) y reejecutar `matriz_r2.py`, `verify_c1.py --db` y `u1_verify.py`.

## Solicitudes fuera de alcance
1. **Orquestador / integrador:** asignar explícitamente a C1 la regeneración de `openapi.yaml`, `routes.go` y `C1-datos.md` con `backend/internal/api/c1gen/` (R2-01, R2-02). En la ronda 1 ningún agente la tenía.
2. **Integrador:** copiar al repositorio los scripts que hoy solo están en anexos o en el scratchpad (R2-08); si se quieren reutilizar los del revisor, están en `rev-r2/` del scratchpad.
3. **Criterios:** correcciones de R2-07. **S1:** R2-03 y R2-04. **U1:** R2-05. **C1:** R2-06.

## Paquetes instalados
- Ninguno. Usé `uvx --offline` (openapi-spec-validator 0.9.0) y `uv run --no-project --with pyyaml --with jsonschema` desde la caché de uv; goose v3.28.0 ya compilado en el scratchpad. Sin internet, sin credenciales, sin Telegram.

## Consumo aproximado
- Unas 85 acciones de herramienta, sin subagentes.
- Leído: unos 330 KB (protocolo, decisiones, plan §1 y §3, criterios, revisión r1, 4 reportes de corrección, secciones de U1, diccionario, reglas, migración, C1 y colorimetría; el OpenAPI y la migración por script).
- Escrito: `REVISION-r2.md` (unos 38 KB), este reporte y unos 40 KB de scripts en el scratchpad.
