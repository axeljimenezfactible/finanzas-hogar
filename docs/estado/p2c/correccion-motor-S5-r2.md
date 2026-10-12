# P2c · correccion-motor-S5-r2 · 11-oct-2026, noche

Agente `correccion-motor-S5-r2`, flujo P2c «Aplicación de la puerta 2».

**Encargo:** el único hallazgo pendiente de `docs/estado/p2c/REVISION-motor-S5-r2.md`, M-1 [media]: `created_at` RFC 3339 con zona horaria (§15.1, §15.3 y §6 S-8).

**Alcance:**
- `backend/internal/precios`, `backend/cmd/motores` (s5), `oraculo/pre_*.py` y las fixtures s5.
- No toqué las bajas B-1 a B-5, que no estaban en el encargo.
- No decidí nada «sin recomendación».
- No hice `git add` ni `commit`.
- No toqué `docs/design/`, `docs/estado/wfp3/` ni `datos/semilla/data.json`.

## Hecho

**HECHO 1 · Premisa de r1 corregida: `pre_comun.py` es solo del oráculo S5.**
- `grep -rln pre_comun --include='*.py' .` da solo `pre_calculo.py`, `pre_datos.py`, `pre_profeco.py` y `pre_run_fixtures.py`, como dijo la revisión.
- Así que entra en mi alcance. La solicitud 4 de `correccion-motor-S5-r1.md` queda atendida aquí.

**HECHO 2 · Reproduje M-1 antes de tocar nada.** Usé 20 sondas en `datos/tmp/corr-s5-r2/fx/` (fuera de git), cada una con `manual` a 2000 a las 12:30:00Z contra `ticket_manual` a 2100:
- `07:00:00-06:00`: motor 2100 `ticket_manual`, oráculo `invalid_date`.
- `12:00:00+00:00`: motor 2000 `manual`, oráculo `invalid_date`. Son las cifras de la revisión.
- También divergían `+06:00`, `-00:00` y la fracción con zona (oráculo `invalid_date`).
- Hallé además una divergencia que la revisión no vio: `2026-10-10T24:00:00Z`.
  - El oráculo viejo lo **aceptaba**: Python 3.14 `fromisoformat` acepta la hora 24 y daba 2100.
  - El motor lo rechaza.
  - RFC 3339 §5.6 limita la hora a 00-23.

**HECHO 3 · Oráculo: `instante()` reescrito desde §15.1 y RFC 3339 §5.6, antes de leer el Go** (`oraculo/pre_comun.py:37-79`).
- La gramática: `full-date «T» hh:mm:ss [.fracción] («Z» | ±hh:mm)`.
- Los rangos: hora 00-23, minuto 00-59 y segundo 00-59, los mismos en la zona; la fecha se valida con `datetime.date`.
- Devuelve el instante exacto en segundos (`Fraction`) desde 0001-01-01T00:00:00Z, ya con la zona restada. Así S-8 compara en UTC.
- Ya no usa `fromisoformat`, y con eso desaparece la hora 24.
- Lo único que llama `instante` es `pre_datos.py:141`, y `created` solo se compara en `pre_calculo.py:61, 71 y 75` (claves de `max` y `sorted`). No cambia ninguna salida existente.

**HECHO 4 · Motor: solo los rangos de la zona** (`backend/internal/precios/eval.go:259-263` y la función nueva `rfc3339Offset` al final del archivo).
- Lo leí después de escribir el oráculo.
- `time.Parse(time.RFC3339)` aceptaba `+24:00` y `+05:60`, que RFC 3339 §5.6 no admite: es su `time-hour` 00-23 y su `time-minute` 00-59.
- En las sondas el motor daba 2000 `manual` en los dos casos, y el oráculo nuevo daba `invalid_date`.
- Ahora el motor los rechaza con `invalid_date`. Lo demás de la lectura queda igual.
- No es una decisión: §15.1 dice «instante RFC 3339», y esos textos no lo son.

**HECHO 5 · Cuatro fixtures nuevas** (`oraculo/pre_genera_fixtures.py`: función `vigente_zona` y tres entradas en `ingesta`). Todas van con `provisional` `null`: S-8 es un supuesto del contrato sin bloqueo, como S-7 en `s5-ingesta-fecha-futura`.
- **`s5-precio-vigente-zona`** (ILUSTRATIVO, `fuente` `manual`, cálculo a mano en `descripcion`): cuatro tiendas, y en cada una `manual` 2000 a las 12:30:00Z contra `ticket_manual` 2100.
  - `t-menos6`: 07:00:00-06:00 = 13:00:00Z, posterior, rige **2100** `ticket_manual`. Es la cifra de la revisión.
  - `t-cero`: 12:00:00+00:00, anterior, rige **2000** `manual`. Es la cifra de la revisión.
  - `t-mas6`: 18:00:00+06:00 = 12:00:00Z, anterior, rige **2000**. Leída como texto sin zona daría 2100.
  - `t-fraccion`: 06:30:00.250-06:00 = 12:30:00.25Z, posterior, rige **2100**. Sin zona daría 2000.
- **`s5-ingesta-zona-hora`** (`+24:00`), **`s5-ingesta-zona-minuto`** (`+05:60`) y **`s5-ingesta-hora-24`** (`24:00:00Z`): las tres dan `invalid_date`.
- Comprobé que las fixtures detectan el error:
  - El oráculo viejo, en una copia aislada `datos/tmp/corr-s5-r2/root-antes/`, falla en `s5-precio-vigente-zona` (las 4 consultas «__falta__», es decir, `invalid_date`) y en `s5-ingesta-hora-24`.
  - El motor viejo aceptaba `+24:00` y `+05:60` (sondas `p-mas24` y `p-min60`).

**HECHO 6 · Medí tres divergencias más, que el contrato no fija, y no las decidí.** Bloqueo nuevo `docs/estado/bloqueos/S5-precio-vigente-captura.md`, con 4 preguntas, cifras, opciones e impacto. La regla mientras no se conteste es PROVISIONAL «en uso: el motor», como en `S5-motor-oraculo-reglas`.
- **Pregunta 1:** con el mismo `created_at`, el motor deja el **primer** precio de la entrada y el oráculo el **último**. Las 4 sondas de `fx2/` divergen.
- **Pregunta 2:** precisión de la fracción. Con 10 dígitos divergen (motor 2000, oráculo 2100); con 7 coinciden.
- **Pregunta 3:** perfil de RFC 3339 (minúsculas `t`/`z`, espacio, segundo 60). Los dos lados las rechazan, así que no hay divergencia, pero falta que el contrato lo diga.
- **Pregunta 4:** claves obligatorias de `prices[]`.
  - Sin `created_at`: el motor acepta y el oráculo da `missing_field`.
  - Con `created_at` `null`: el motor acepta y el oráculo da `invalid_date`.
  - Sin `is_promotional`: el motor acepta y el oráculo da `missing_field`.
- Las marcas que citan el bloqueo están en `pricing.go:214-215` (`later`), `eval.go:259-260`, `pre_calculo.py:59-62` (`vigente`) y `pre_comun.py:41 y 56-58`.
- Ninguna fixture usa estos casos, y ninguna cifra canon cambia.

**HECHO 7 · Pruebas Go nuevas.**
- `TestEvalCreatedAtZona` (`eval_test.go`) cubre las 4 zonas de la fixture, con `Eval` y JSON.
- `TestEvalErrores` tiene dos renglones nuevos, `+24:00` y `+05:60`, que dan `invalid_date`.

## Verificación (comandos desde `/home/deepsec/finanzas-hogar`; salidas reales)

| Qué | Comando | Salida |
|---|---|---|
| Sondas después, motor contra oráculo | `datos/tmp/corr-s5-r2/motores s5 dump --fixtures datos/tmp/corr-s5-r2/fx …` y `pre_run_fixtures.py --fixtures …/fx` | «20 sondas, 19 iguales; difieren: p-frac10 (motor 2000 manual, oráculo 2100 ticket_manual)» (pregunta 2). Antes eran 8 distintas |
| Las fixtures se regeneran | `python3 -I oraculo/pre_genera_fixtures.py` y `diff -r` contra la copia previa | «36 fixtures»; solo aparecen las 4 nuevas y las 32 previas quedan idénticas |
| Oráculo | `python3 -I oraculo/pre_run_fixtures.py` | «36 fixtures, 36 coinciden, 1359 hojas de esperado, 0 diferencias»; los 32 casos previos de `salidas/s5.json` quedan idénticos (`True`) |
| compara | `python3 -I oraculo/compara.py s5 --motor datos/tmp/corr-s5-r2/motor-s5.json` | oráculo contra esperado 958/958; motor contra esperado 958/958; oráculo contra motor 958/958; provisionales 315/315; **diferencias: 0** |
| gofmt y vet | `gofmt -l internal/precios cmd/motores`; `go vet ./internal/precios/... ./cmd/motores` | vacío; «vet OK» |
| Pruebas y cobertura K25 | `go test -count=1 -cover ./internal/precios/... ./cmd/motores` | precios 96.0 %, asistida 100.0 %, ocr 99.3 %, profeco 95.3 %, cmd/motores 48.7 %; todo `ok` |
| Pruebas nuevas | `go test -count=1 -v -run 'TestEvalErrores\|TestEvalCreatedAtZona' ./internal/precios/` | `--- PASS: TestEvalErrores`, `--- PASS: TestEvalCreatedAtZona` |
| Pruebas nombradas | `-run 'Comparativa\|Membresia\|Unidad\|MejorTienda'` y `-run 'Frescura\|Fuente\|LecturaAsistida\|Mercado\|PreciosViejos'` | `ok` |
| Canon | `python3 -I oraculo/pre_verify.py` | «TODO COINCIDE» (15 OK; Walmart +138.50 = 13850) |
| data.json | `git diff --quiet HEAD -- datos/semilla/data.json` | sin cambios |
| Global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | **EXIT=0**; «s5 · 36 fixtures, 36 coinciden … salidas/s5.json idéntica»; «pre_verify.py: 15 comprobaciones OK»; «WF-1: verificación completa OK»; «Verificación ampliada OK» |

## Supuestos
- **S-a:** los rangos de la zona (hora 00-23, minuto 00-59) y de la hora (00-23) salen de la gramática de RFC 3339 §5.6, a la que remite §15.1. No los tomo como decisión nueva.
- **S-b:** las 4 fixtures nuevas van con `provisional` `null`. Prueban S-8, que es un supuesto del contrato sin bloqueo, y no dependen de `S5-precio-vigente-captura`: ningún par empata en el instante, la fracción es de 3 dígitos y las letras van en mayúscula.

## Bloqueos
- Nuevo: `docs/estado/bloqueos/S5-precio-vigente-captura.md` (preguntas 1 a 4, ABIERTO).
- Siguen abiertos, sin tocar: `S5-motor-oraculo-reglas`, `S5-carritos-granel` y `S5-profeco-filtro`.

## Pendientes
- Las preguntas 1 a 4 del bloqueo nuevo. Cuando se contesten, se alinea el lado que no cumpla y se agrega una fixture.
- Las bajas B-1 a B-5 de la revisión r2 no estaban en mi encargo.

## Propuestas para Axel
- Ninguna que decidir hoy.
- La pregunta 1 del bloqueo nuevo («¿cuál precio rige si dos capturas tienen el mismo instante?») llega a Axel solo si contratos no la juzga técnica.

## Solicitudes fuera de alcance
1. **Contratos (`S5-interfaz.md`):**
   - agregar a la tabla de §15 las fixtures `s5-precio-vigente-zona`, `s5-ingesta-zona-hora`, `s5-ingesta-zona-minuto` y `s5-ingesta-hora-24`;
   - fijar en §15.1 lo que piden las preguntas 2 a 4 de `S5-precio-vigente-captura`;
   - y en §6 la pregunta 1, si la juzga técnica.
2. **Integrador:**
   - listar `S5-precio-vigente-captura` en `docs/estado/bloqueos/INDICE.md`, junto con `S5-motor-oraculo-reglas` (B-4 de la revisión);
   - registrar la excepción de §15.3 (B-1): esta corrección también tocó los dos lados. Escribí el oráculo primero, desde el contrato y RFC 3339, y solo después leí `eval.go` para el arreglo del motor.

## Consumo aproximado
- 34 acciones de herramienta.
- Unos 150 KB leídos.
- Unos 26 KB escritos a mano: código, este reporte y el bloqueo. Aparte, los generadores escribieron las fixtures y `salidas/s5.json`.
