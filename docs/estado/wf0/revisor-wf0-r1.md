# Reporte WF-0 · revisor-wf0-r1 (revisor independiente, ronda 1 de 3)

10-oct-2026 · Estado: **terminado** (sin cierre por tope) · Veredicto: **no aprobado** por 1 hallazgo alto · Detalle completo: `docs/estado/wf0/REVISION-r1.md`.

## Hecho
- Leí completos el protocolo, decisiones.md, los criterios `INDICE.md`, `S1.md`, `U1.md`, `U2.md` y `C1.md`, los 5 reportes de WF-0 y los 30 bloqueos.
- Del plan leí §0–§3 y §4–§5 con `grep -n` + `sed -n`. Los artefactos grandes (`openapi.yaml`, U1, C1, diccionario, migración) los leí por secciones o con scripts.
- Reejecuté las 9 verificaciones de la tarea con scripts propios (sin reutilizar código de los agentes) y además los verificadores de U1 y C1 extraídos de sus reportes.
- Evalué 60 criterios (17 de S1, 14 de U1, 15 de C1 y 14 de U2, de los que 5 no aplican en WF-0). Revisé además la estructura y la trazabilidad de los 248 criterios.
- Hice la tabla de cobertura de las decisiones A–I: 40 aplicables, 39 cubiertas y 1 parcial porque depende de WF-1.
- Registré 12 hallazgos: 1 alto, 7 medios y 4 bajos.
- No corregí ningún archivo de otros agentes.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf0/REVISION-r1.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf0/revisor-wf0-r1.md` (este reporte)
- Temporales: `/home/deepsec/finanzas-hogar/datos/tmp/` (`rev.db`, `rev-goose.db`), creada y **borrada**; `datos/` vuelve a tener solo `semilla/`.
- Scripts del revisor en el scratchpad de la sesión, carpeta `rev-r1/`: `up.sql` y `down.sql` (secciones de la migración), `regresion_rev.py`, `negativas_rev.py`, `matriz_rev.py`, `c1_rev.py`, `color_rev.py`, `reglas_rev.py` y `criterios_rev.py`. También `mirror/`, donde corrí colorimetria.py sin tocar el repositorio.

## Verificación
1. **SQL**: Up extraído con awk → `sqlite3 -bail datos/tmp/rev.db` → `UP OK`, 49 tablas STRICT, 100 índices con nombre, 62 disparadores. `PRAGMA foreign_key_check` = 0 filas, `integrity_check` = ok. Con goose v3.28.0 (`-dir backend/db/migrations sqlite3 … up/down/up`): OK, versión 1; down deja 0 tablas; up otra vez ok.
2. **Go**: `cd backend && go build ./... && go vet ./...` → BUILD OK, VET OK; `gofmt -l .` sin salida (Go 1.26.5, `go.mod` con `go 1.26`).
3. **Colorimetría**:
   - `colorimetria.py --sin-png` en un árbol espejo regenera md, csv, json y html idénticos (`cmp`) a los del repositorio: 202 filas, 0 `cumple=false`, «VERIFICACIÓN OK».
   - Script propio: recálculo de los 202 pares, 0 diferencias mayores a 0.01 y 0 bajo umbral. Los 190 tokens del JSON = Anexo A, y 171 valores de las tablas del cuerpo coinciden.
   - Matices a ±0.7° de lo declarado; L por paso ±0.0016.
   - Pares < 20° y ΔL < 0.15: los 7 que declara §7.7.
   - Series: ΔL consecutivo mín. 0.179 / 0.223; contraste mín. 3.15 / 3.11.
4. **Matriz** (`matriz_rev.py`):
   - 199 N-xxx en U1; 220 pares en U1 §3.29 = 220 filas en C1 §10; `x-necesidades` = C1 §10.
   - N-xxx sin endpoint: N-003 y N-301 (locales) y **N-151, N-152, N-153, N-165** (bloqueo). Operación sin N-xxx: `GET /openapi.yaml` (justificada).
   - `u1_verify.py` (anexo de U1) da «FALLA U1-03» por las mismas 4. `verify_c1.py` (anexo de S6a) da «TODO OK» porque pone esas 4 en su lista de excepciones.
5. **Decisiones**: tabla en `REVISION-r1.md` §3.
6. **PROVISIONAL**: toda marca remite a un bloqueo existente (U1 10, C1 44 + 63 en YAML, diccionario 6, reglas 22, migración 4). Las semillas de la migración son valores decididos. En el OpenAPI, el único valor por defecto es `limit` 50 (PROVISIONAL), y los límites que no vienen de las decisiones están marcados.
7. **Criterios** (`criterios_rev.py`): 248 = INDICE (189 / 59); 248 con «Verifica» y «Evidencia»; 0 enunciados vagos sin cifra; 35 sin línea en su §5 Trazabilidad.
8. **Toques**: U1 §4 tiene la regla de conteo; F-1 = «+», categoría, «Guardar» = 3.
9. **Regresión** (`regresion_rev.py` sobre la base migrada):
   - Suma de pagos de la despensa 1 = 389500; parte de Axel = 239528.
   - 17 eventos y series A y B de 14 valores; mínimos 809450 / 500000; efectivo y devengado de P1–P6 por mes calendario.
   - **109 valores, 0 diferencias**. `foreign_key_check` vacío e `integrity_check` ok después de los inserts.
10. **Pruebas negativas** (`negativas_rev.py`): 15 rechazos esperados. Tres casos los acepta la base y los resuelve la API:
    - partes 60 % + 50 % (invariante 8);
    - tercera cuenta;
    - participación cambiada a externa con cuenta viva (invariante 12).
11. **OpenAPI**: `uvx --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml` → OK, salida 0. `c1_rev.py`:
    - 874 de 874 respuestas de error → `Error`;
    - 253 propiedades `_cents` → `MoneyCents` y 0 `number`;
    - 906 propiedades de entidad respaldadas (824 columnas, 71 calculadas, 11 anidadas);
    - 6 filtros con definición única;
    - públicas solo `GET /health` y `POST /api/v1/auth/login`.

## Supuestos
- Corrí `colorimetria.py` sobre una copia en el scratchpad y no en su ruta del repositorio, para no reescribir artefactos ajenos (el script escribe en la raíz de su propio árbol). El resultado es el mismo código sobre las mismas entradas y salió idéntico byte a byte.
- Clasifiqué H-01 como alta porque C1-02 y U1-03 (alta) solo admiten excepciones técnicas, y porque la coherencia S1/C1/U1 es el encargo de este revisor (plan §5).
- Clasifiqué H-03 (U2a) como media porque U2 se aprueba en la puerta 2.

## Bloqueos
Ninguno nuevo. Los existentes que pesan en la puerta 1 son `C1-lista-compras`, `C1-huecos-u1-s1` y `criterios-nxx` (ver H-01, H-02 y H-06).

## Pendientes
- Ronda 2: reverificar H-01 (matriz con 0 N-xxx sin endpoint ni declaración Local) y la corrección o justificación escrita de H-02 a H-08.

## Solicitudes fuera de alcance
1. **Integrador**: guardar en el repositorio los generadores y verificadores que hoy viven en el scratchpad efímero, antes de la ronda 2 (H-05):
   - `…/scratchpad/c1/*` en `docs/contratos/tools/c1/`;
   - `u1_verify.py` en `docs/design/verificar_u1.py`;
   - el generador del diccionario de S1 y su verificador en `ops/tools/`;
   - `gen_canon.py` de criterios;
   - mis scripts `…/scratchpad/rev-r1/*.py` en `docs/estado/wf0/herramientas-revision/`. Para la regresión: extraer Up con awk, crear `datos/tmp/rev.db`, correr `regresion_rev.py` y luego `negativas_rev.py`.
2. **S1, C1, U1, criterios**: las correcciones de `REVISION-r1.md` §5.

## Paquetes instalados
- Ninguno con pacman ni pip global.
- Usé entornos aislados de la caché de uv: `openapi-spec-validator` 0.9.0 (`uvx`), y `pyyaml` 6.0.3 y `jsonschema` 4.26.0 (`uv run --with`).
- Usé el binario goose v3.28.0 que S1 compiló en el scratchpad, sin reinstalarlo.
- Sin red fuera de eso, sin credenciales y sin Telegram.

## Consumo aproximado
- Unas 75 acciones de herramienta, sin subagentes.
- Leído, unos 420 KB: protocolo, decisiones, plan por secciones, 5 criterios, 5 reportes, 30 bloqueos, y secciones de U1, C1, diccionario, reglas, colorimetría y migración. El OpenAPI, `model.json` y el CSV solo los leí por script.
- Escrito: `REVISION-r1.md` (unos 33 KB), este reporte (unos 6 KB) y scripts en el scratchpad (unos 30 KB).
