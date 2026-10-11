# Reporte de agente · revisor-C1-r3 (flujo P1, 11-oct-2026)

Revisor independiente, ronda 3 de 3, de C1 v1.1. No corrige nada. Informe: `docs/estado/p1/REVISION-C1-r3.md`.

## Hecho
- Leí el protocolo completo, §K de `decisiones.md`, `docs/criterios/C1.md` v1.1 y `correccion-C1-r2.md`. De `C1-datos.md` leí solo Cambios, §4, §6 y §7.1, y una línea de §9. De `openapi.yaml` solo por grep y scripts. De `extractos/INDICE.md`, la cabecera.
- Volví a correr las verificaciones (a) a (j) con comandos y scripts propios en `datos/tmp/rev-c1-r3/`.
- Comprobé los 3 hallazgos de r2 y los bajos H6, H7 y H8.
- Veredicto: **NO APROBADO**, con 2 hallazgos altos (H-r3-01 y H-r3-02), 2 medios (H-r3-03 y H-r3-04) y 5 bajos.

## Archivos
- Creado: `docs/estado/p1/REVISION-C1-r3.md`
- Creado: `docs/estado/p1/revisor-C1-r3.md`
- Temporales: `datos/tmp/rev-c1-r3/` (base c1.db y los scripts chk1 a chk5), borrados al terminar.

## Verificación
- `python3 -I backend/internal/api/c1gen/build.py --check`: IDÉNTICO en los 21 archivos (openapi.yaml, routes.go, C1-datos.md, INDICE.md y 17 extractos).
- `goose v3.28.0 … sqlite3 datos/tmp/rev-c1-r3/c1.db up`: 00001 y 00002 OK. Después `verify_c1.py --db`: 24 OK, 0 FALLA, «TODO OK».
- `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator`: openapi.yaml OK y los 17 extractos OK.
- `python3 -I docs/design/tools/verificar_u1.py`: TODO OK, 0 N-xxx sin operación ni Local.
- `go build ./... && go vet ./... && gofmt -l .` con go1.27.2: OK, OK, vacío.
- Extractos: 201 operationId = `grep -c` de openapi.yaml, 0 duplicados, el mayor de 76,202 B, `componentes-comunes.yaml` presente, todos parsean.
- §K: scope, borrador, candidatos, Idempotency-Key en gasto y peso, rate_decimals, tope por persona, restricciones, campos de Credit, proyección, goal, family y `week`. Están todos, salvo la semana en `getBasketComparison`.
- PROVISIONAL: 4 `x-provisional`, que citan 3 bloqueos abiertos. Sin marca: la regla del domingo en `getMealPlanWeek`.
- Cifras de ejemplo: 0 sin origen en el oráculo ni en la semilla.

## Supuestos
- Acepto que una marca PROVISIONAL cite un bloqueo de otra fase (U1-, S1-, criterios-) y no solo `C1-<tema>`, como ya hace el contrato.

## Bloqueos
Ninguno nuevo.

## Pendientes
Las correcciones de H-r3-01 a H-r3-04 (rutas y propuesta en el informe). Es la ronda 3 de 3, así que el orquestador decide si C1 v1.1 va a Axel con estos hallazgos o si hay otra corrección.

## Solicitudes fuera de alcance
1. `docs/design/U1-arquitectura.md` línea 66 y `docs/criterios/U1.md` U1-16: cambiar `period` por `week` (H-r3-03).
2. `docs/contratos/diccionario-datos.md`: aplicar las solicitudes 1 y 2 de `correccion-C1-r2.md` (H-r3-04).
3. S1, `credit.vat_on_interest_bp`: estado «sin dato» para el IVA (H-r3-06, baja).

## Paquetes instalados
Ninguno.

## Consumo aproximado
- Acciones de herramienta: unas 38 de 60.
- Leído: unos 140 KB, entre documentos por partes y salidas de comandos y scripts. Ningún archivo grande leído entero.
- Escrito: unos 33 KB (5 scripts temporales de unos 13 KB, el informe de 16 KB y este reporte de 4 KB).
