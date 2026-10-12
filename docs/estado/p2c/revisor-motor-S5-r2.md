# P2c · revisor-motor-S5-r2 · 11-oct-2026, noche

Revisión independiente, ronda 2 de 3, del motor S5 después de `correccion-motor-S5-r1`.
- El informe completo está en `docs/estado/p2c/REVISION-motor-S5-r2.md`.
- No corregí nada, no hice `git add` ni `git commit` y no toqué `docs/design/` ni `docs/estado/wfp3/`.

## Hecho
- **Veredicto: NO APROBADO.** Hay 0 altos, 1 medio y 5 bajos.
- **Los tres medios de r1 están cerrados:**
  - M-1 y M-2 quedan como PROVISIONAL con el bloqueo `S5-motor-oraculo-reglas` (8 preguntas), con marcas en los dos lados.
  - M-3: la nota de S5-09 tiene marca y prueba de texto.
- **M-1 nuevo [media]:**
  - `oraculo/pre_comun.py:39,51` solo acepta `created_at` con `Z`.
  - El motor (`eval.go:259`) acepta RFC 3339 con `±hh:mm`, como dice §15.1.
  - Con `-06:00` o `+00:00`, el oráculo responde `invalid_date` a toda la fixture y el motor da un precio, porque la zona decide S-8.
  - La corrección r1 lo dejó fuera de alcance porque supuso que el archivo es compartido de S2 a S5. Solo lo importan los `pre_*.py` de S5.
- **Bajos:**
  - B-1 a B-3 siguen de r1.
  - B-4: el bloqueo nuevo no está en `INDICE.md`.
  - B-5: la pregunta 5 del bloqueo estaba sin medir. La medí: el motor da `unknown_presentation` y el oráculo `skipped` 6.

## Verificación
Todo desde `/home/deepsec/finanzas-hogar`, con temporales en `datos/tmp/rev-s5-r2/` (ya borrados):
```
$ (cd backend && gofmt -l internal/precios cmd/motores; go vet ./internal/precios/... ./cmd/motores && echo vet-ok)
vet-ok
$ (cd backend && go test -count=1 -coverprofile=../datos/tmp/rev-s5-r2/s5.cov ./internal/precios/...)
precios 96.0% · asistida 100.0% · ocr 99.3% · profeco 95.3% · total 88.8%
$ python3 -I oraculo/pre_genera_fixtures.py && diff -r <copia previa> oraculo/fixtures/s5      → sin salida
$ python3 -I oraculo/pre_run_fixtures.py --out datos/tmp/rev-s5-r2/oraculo-s5.json
32 fixtures, 32 coinciden, 1336 hojas de esperado, 0 diferencias   (cmp contra salidas/s5.json: idéntico)
$ go -C backend build -o ../datos/tmp/rev-s5-r2/motores ./cmd/motores && … s5 dump … && python3 -I oraculo/compara.py s5 --motor …
oráculo contra motor: 935 de 935 · provisionales: 315 de 315 · diferencias: 0
$ (raíz aislada) motores s5 dump --fixtures fx  +  pre_run_fixtures.py --fixtures fx  +  dif.py propio
604 hojas, 13 diferencias: todas declaradas en S5-motor-oraculo-reglas (p1: 2, p2: 8, p3: 3)
rv2-profeco prod: read 21, rejected 12, filtered_out 1, skipped 2, loaded 5, loaded_proposed 0, duplicates 1 (motor = oráculo = mi esperado)
rv2-instante-0 (-06:00): motor price 2100 ticket_manual · oráculo {"error": "invalid_date"}   ← M-1
canon: 309450 / 323300 (+13850) · idénticos 277350 / 287400 (+10050) · 61/60/55 · 33/3 · 22/3, ties 8
$ git diff --quiet HEAD -- datos/semilla/data.json && echo sin-cambios
sin-cambios
$ make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose
EXIT=0 · WF-1: verificación completa OK · Verificación ampliada OK
```

## Supuestos
- **S-r2-1.** Una divergencia entre el motor y el oráculo cuenta como cerrada si queda PROVISIONAL «en uso: el motor», con un bloqueo ABIERTO que trae cifras y opciones y con marcas en los dos lados. En r1, la fixture se pidió «con la regla fijada».
- **S-r2-2.** M-1 lo califiqué medio y no alto. No cambia ninguna cifra de fixtures ni del canon, pero rompe §15.3 (motor = oráculo) con una entrada válida según §15.1.
- **S-r2-3.** Para no tocar `oraculo/fixtures/s5/_datos/`, corrí el oráculo desde una copia de `oraculo/*.py` en una raíz aislada. Su `RAIZ` apunta a esa copia, y el código es idéntico.

## Bloqueos
- No abro ninguno. M-1 no es una regla sin decidir: §15.1 dice RFC 3339.
- Siguen ABIERTOS: `S5-motor-oraculo-reglas` (8 preguntas), `S5-carritos-granel`, `S5-carritos-detalle` y `S5-profeco-filtro`.

## Pendientes
- Para la ronda 3:
  - Comprobar que `instante` del oráculo acepta `±hh:mm`.
  - Comprobar que existe una fixture donde la zona horaria decide S-8.
- Los bajos B-1 a B-5 son para el integrador o para contratos.

## Propuestas para Axel
- Nada que decidir en esta ronda.
- Las preguntas 1 a 3 de `S5-motor-oraculo-reglas` (piezas a granel, precio de lo emparejado y referencias de PROFECO) siguen sin recomendación y no mueven dinero del canon.

## Solicitudes fuera de alcance
1. **Integrador:**
   - Agregar `S5-motor-oraculo-reglas.md` a `docs/estado/bloqueos/INDICE.md` (B-4).
   - Sumar a la pregunta 5 la medición de B-5.
2. **Integrador o contratos:** B-3 (tabla de §15 y «Cómo seguí» de `S5-carritos-granel`).

## Consumo aproximado
Unas 30 acciones de herramienta · unos 120 KB leídos · unos 30 KB escritos:
- Informe y reporte: unos 19 KB.
- Generador, comparador y fixtures temporales: unos 11 KB, ya borrados.
