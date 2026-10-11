# Reporte WF-0 · revisor-wf0-r3

Revisor independiente, ronda 3 de 3 del flujo WF-0 «Cimientos» · 11-oct-2026 · Estado: **terminado**, sin cierre por tope (`cerrado_por_tope = false`).
Informe completo de la revisión: `/home/deepsec/finanzas-hogar/docs/estado/wf0/REVISION-r3.md`.

**Veredicto: APROBADO.** Hay 0 hallazgos de severidad alta, 0 de media y 10 de baja.
- R2-01 y R2-02 (alta, C1) quedan cerrados.
- Los bajos se registran: R3-01 a R3-04 son nuevos y R2-03 a R2-08 vienen de r2, en archivos que no se tocaron esta ronda.

## Hecho
1. Leí el protocolo, las decisiones, las secciones §1 a §4 del plan (con `sed -n`), los criterios INDICE, S1, U1, U2 y C1, los entregables de WF-0, los 32 bloqueos y los reportes de `docs/estado/wf0/` (REVISION-r1 y r2 y las correcciones de r1 y r2).
2. Reejecuté las 9 verificaciones de la tarea con scripts propios y con los verificadores de los agentes (detalle abajo y en REVISION-r3 §1).
3. Escribí la tabla de criterios (S1, U1, C1, U2a y criterios), la cobertura de las decisiones A a I (35 renglones), el estado de los hallazgos de r2 y los hallazgos de r3.
4. No corregí nada ni escribí fuera de las dos rutas asignadas. Creé `datos/tmp/` para las bases temporales y la borré al final (no existía antes).

## Archivos
- Creados: `/home/deepsec/finanzas-hogar/docs/estado/wf0/REVISION-r3.md` y `/home/deepsec/finanzas-hogar/docs/estado/wf0/revisor-wf0-r3.md` (este reporte).
- Temporales, fuera del repositorio: el scratchpad de la sesión, carpeta `rev-r3/`, con `up.sql`, `down.sql`, `mirror/`, `u2a_r3.py`, `u2a_pairs.py`, `chart_r3.py`, `matriz_r3.py`, `c103_r3.py`, `c1conv_r3.py`, `decisiones_r3.py`, `prov_r3.py`, `criterios_r3.py`, `regresion_r3.py` y `u1_verify.py` (extraído del anexo del reporte de U1).
- Temporales en el repositorio, ya borrados: `datos/tmp/rev.db`, `rev-down.db`, `rev-goose.db` y `rev-regr.db`.

## Verificación

**1. SQL en `datos/tmp/rev.db`.** Up con `sqlite3 -bail`: `UP_OK`.
- 54 tablas STRICT, 147 índices, 68 disparadores.
- `PRAGMA foreign_key_check` = 0 filas; `integrity_check` = ok.
- Down deja 0 objetos. goose v3.28.0 up/down/up: OK.
- Semillas: 7 `app_setting`, 4 `meal_slot` y 7 `notification`, todas con valores de las decisiones.
- Negativas rechazadas: cuenta de Karen, hash que no es argon2id, segunda co-titular, `12.5` en centavos, `10-oct-2026`, ventana del 10 al 9, 2027-02-30 y paso a externa con cuenta.

**2. Go.** `cd backend && go build ./... && go vet ./...`: BUILD_OK y VET_OK (go1.26.5).
- `gofmt -l .` sin salida.
- `routes.go` = `openapi.yaml`: 197 = 197.

**3. Colorimetría.** `colorimetria.py --sin-png`, corrido en un árbol espejo: «VERIFICACIÓN OK».
- md, csv, json y html idénticos byte a byte al repositorio.
- `contraste.csv`: 202 filas, 0 `cumple=false`; recálculo propio sin diferencias.
- 190 tokens de `C2-colores.json` = 190 filas del Anexo A de `colorimetria.md`, 0 hex distintos.
- Las tablas del cuerpo usan solo hex de token. Plan §4.3: 22 de 22.

**4. Matriz C1.** `matriz_r3.py`:
- 199 N-xxx: 0 sin operación ni Local, 0 pantallas sin operación.
- 1 operación sin N-xxx: `getOpenApiDocument`, la excepción «estáticos».
- §10: 220 filas = 220 pares de U1; 0 filas BLOQUEO.

Herramientas de los agentes:
- `python3 -I backend/internal/api/c1gen/build.py --check`: IDÉNTICO ×3.
- `uv run --no-project --offline --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db datos/tmp/rev-goose.db`: «RESULTADO: TODO OK».
- `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml`: OK.
- `u1_verify.py`: U1-03 en OK; su única FALLA es la excepción «estáticos» (R3-03).

**5. Decisiones A–I.** `decisiones_r3.py`: 35 de 35 renglones con respaldo en S1, U1 o C1 (tabla en REVISION-r3 §3).

**6. PROVISIONAL.** `prov_r3.py`:
- 108 marcas en texto y 54 nodos `x-provisional`; todas las de contenido citan un bloqueo existente.
- Las 18 cifras de C1 §9 salen de `model.json` (16) o de `data.json` (2).
- 0 valores inventados en cifras, esquema o contrato.
- 4 límites técnicos sin marca (R3-02).

**7. Criterios.** `criterios_r3.py`:
- 248 criterios (189 alta y 59 media); los 248 con «Verifica» y «Evidencia» y en la §5 de su archivo.
- Siguen las incoherencias menores de R2-07.

**8. ≤ 3 toques.** U1 §4 y F-1: «+», chip y «Guardar» = 3. `ExpenseCreate` exige solo campos cubiertos por valores por defecto o por esos 3 toques.

**9. Regresión de S1.** `regresion_r3.py`:
- Suma de pagos de la despensa 1 = **389500** y parte de Axel = **239528**.
- 98 valores comparados contra `model.json`, 0 diferencias: despensas, pagos y partes, 17 eventos A/B, series de 14 valores, mínimos 809450/500000, efectivo y devengado de oct-2026 a mar-2027.
- `fk_check` 0.

## Supuestos
1. Califiqué como «baja» los hallazgos que no cambian cifras, esquema, rutas ni contratos y que tienen su bloqueo o su corrección obvia. Seguí la definición de severidad del INDICE §1 y el mismo criterio que la ronda 2 para los heredados.
2. Las 11 propiedades anidadas de los esquemas de entrada (`items`, `ingredients`, `parts`, `shares`, `presentations`, `expense`, `cost_share_rule_parts`) cumplen C1-03. Son listas u objetos de esquemas hijos cuyas propiedades sí son columnas de su tabla; no son campos sin respaldo.

## Bloqueos
Ninguno nuevo. Los 32 existentes en `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/` se revisaron y ninguno esconde un valor inventado.

## Pendientes
Ninguno de esta revisión. Para el integrador y la puerta 1 quedan los 10 hallazgos de baja (REVISION-r3 §6) y las preguntas abiertas de §7.

## Solicitudes fuera de alcance
1. **Integrador**, sobre `docs/estado/wf0/U1-arquitectura.md`: anotar el cierre de U1-03 y ajustar o incorporar `u1_verify.py` para que acepte la excepción de C1 §11.2 (R3-03).
2. **C1**, en la próxima regeneración tras la puerta 1 (`backend/internal/api/c1gen/gen3_ops.py`): alinear la semana por omisión de la lista de compras con U1 §0.2 (R3-01), declarar los límites técnicos de R3-02 y atender R2-06.
3. **U2 en WF-1** (`docs/design/tools/colorimetria.py`, `colorimetria.md` §7.7): agregar la marca al barrido de pares cercanos (R3-04).
4. **S1** (`00001_esquema.sql`, diccionario, `reglas-negocio.md`) y **criterios**: los heredados R2-03, R2-04, R2-07 y R2-08.

## Paquetes instalados
Ninguno. Usé lo que ya estaba en las cachés del usuario:
- goose v3.28.0 en `scratchpad/gobin/`.
- `uvx --offline openapi-spec-validator==0.9.0`.
- `uv run --offline --with pyyaml --with jsonschema`.

## Consumo aproximado
- Unas 75 acciones de herramienta, sin subagentes ni internet.
- Leí unos 380 KB en lecturas dirigidas (plan por secciones, criterios, reportes, bloqueos, extractos de U1, C1 y colorimetría) y escribí unos 45 KB (dos reportes y los scripts del scratchpad).
