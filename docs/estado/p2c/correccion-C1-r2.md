# Corrección de C1 v1.3 (puerta 2) · ronda 2

Agente: `correccion-C1-r2`, flujo P2c, 11-oct-2026, noche (el reloj del equipo ya marcaba 12-oct). Corrige el único hallazgo medio de `docs/estado/p2c/REVISION-C1-r2.md`, **M-1**. No había hallazgos altos. Los bajos B-1 a B-5 no venían en el encargo y no los toqué. No decidí ninguna pregunta «sin recomendación».

Alcance tocado: solo `backend/internal/api/c1gen/` (`gen2_schemas.py`, `verify_c1.py` y `README.md`) y las salidas que regenera `build.py`. Cambiaron 3 salidas: `docs/contratos/openapi.yaml`, `extractos/ajustes-c.yaml` y `extractos/INDICE.md`. `C1-datos.md`, `routes.go` y los otros 16 extractos se regeneraron idénticos. `gen3_ops.py` no cambió: usa la constante `CREDIT_VAT_RULE_OP`, que vive en `gen2_schemas.py`. Sin `git add` ni `git commit`.

## Hecho

### M-1 (media) · el CHECK también rechaza un porcentaje con `vat_included` null o ausente · CORREGIDO

**HECHO: la regla de la base.** La 00004 (l. 99) dice `CHECK (vat_included IN (0, 1) AND (vat_on_interest_bp IS NULL OR vat_included IS 0))`, y la columna no tiene valor por omisión. Lo comprobé con la base (v-1):

| Alta de `store_credit` con `vat_on_interest_bp` | `vat_included` | Base |
|---|---|---|
| 1600 | true | rechaza |
| 1600 | null | rechaza |
| 1600 | ausente | rechaza |
| 1600 | false | acepta |
| sin porcentaje | true, null o false | acepta |

| Merge Patch sobre un crédito con (false, 1600) | Base |
|---|---|
| `{"vat_included": true}` | rechaza |
| `{"vat_included": null}` | rechaza |
| `{"vat_included": true, "vat_on_interest_bp": null}` | acepta |

**Cambio en los esquemas.** `gen2_schemas.py`, constante `CREDIT_VAT_RULE`, que va en las descripciones de `CreditCreate` y `CreditUpdate`:
- Ahora escribe la regla completa: «vat_on_interest_bp con valor solo se admite con vat_included = false».
- Si la combinación que quedaría guardada trae un porcentaje y `vat_included` = true, null o ausente, la API responde 422 `validation_failed` con `field = vat_on_interest_bp` y no guarda nada. En el alta cuenta el cuerpo. En el Merge Patch cuenta el cuerpo sobre lo ya guardado.
- Agrega los 2 ejemplos de la revisión:
  - el alta `{"kind": "store_credit", …, "vat_on_interest_bp": 1600}` sin `vat_included`;
  - `{"vat_included": null}` sobre un crédito guardado con `vat_included` = false y `vat_on_interest_bp` = 1600.
- Conserva los ejemplos de r1 y la frase «La API no cambia el otro campo por su cuenta: no vacía vat_on_interest_bp ni cambia vat_included». Le añade «(tampoco deduce false a partir del porcentaje)». No es una regla nueva: solo hace explícito lo que ya decía «ni cambia vat_included» para el caso del alta sin ese campo.

**Cambio en las operaciones.** Constante `CREDIT_VAT_RULE_OP`, que usan `createCredit` y `updateCredit`. Es la versión corta de la misma regla: «Si lo que quedaría guardado (el cuerpo o, en un cambio, el cuerpo sobre lo ya guardado) trae vat_on_interest_bp con valor y vat_included = true, null o ausente, responde 422 validation_failed con fields[].field = vat_on_interest_bp; la API no cambia el otro campo por su cuenta (…)».

**Dónde quedó.**
- `openapi.yaml`: l. 6168 (`createCredit`), l. 6330 (`updateCredit`), l. 15487 (`CreditCreate`) y l. 15589 (`CreditUpdate`).
- `extractos/ajustes-c.yaml`: l. 343, l. 387, l. 1344 y l. 1446.
- Las dos operaciones conservan su respuesta 422.

**Tamaño.** `ajustes-c.yaml` pasa de 78,546 a 79,094 bytes. Sigue bajo el tope de 80,000. `INDICE.md` cambia solo esa fila.

**Comprobación nueva C1-04 en `verify_c1.py`** (l. 457 en adelante, después de EXT):
- Copia la base a memoria con `db.backup()`. Así funciona también con `--db`, que se abre solo lectura (así la corre `make verify-c1`).
- Apaga la llave a `person` en la copia: solo se prueba el CHECK de `credit`.
- Cada combinación va dentro de un `SAVEPOINT` que se deshace:
  - alta con 1600 y `vat_included` true, null, false o ausente;
  - las mismas sin porcentaje, como control;
  - los 3 Merge Patch sobre (false, 1600).
- Solo cuenta como «rechaza» un `IntegrityError` con «CHECK constraint failed». Cualquier otro error hace fallar la comprobación, así que una prueba mal armada no puede pasar sin que se note.
- Con lo que la base rechaza y acepta, arma las frases que exige:
  - «vat_on_interest_bp con valor y vat_included = true, null o ausente»;
  - «vat_on_interest_bp con valor solo se admite con vat_included = false».

  Las exige en los 8 textos: `CreditCreate`, `CreditUpdate`, `createCredit` y `updateCredit`, cada uno en `openapi.yaml` y en su extracto. En los esquemas exige además los ejemplos de cada combinación que se rechaza (el alta sin `vat_included` y los dos Merge Patch) y la forma de pasar a «incluye IVA».
- Si un día cambia el CHECK (por ejemplo, con un valor por omisión), las frases que se esperan cambian y la comprobación falla hasta que se corrija el texto.

**Las comprobaciones detectan el defecto.** Contra la salida anterior, C1-04 falló con 20 faltantes. Los primeros eran «true, null o ausente» en `CreditCreate`, el ejemplo del alta y `{"vat_included": null}` (v-2).

**Documentación.** `c1gen/README.md` describe C1-04 junto a las comprobaciones de r1.

## Verificación (comandos y salidas)

| Id | Comando | Salida |
|---|---|---|
| v-1 | Script propio: Up de 00001 a 00004 en memoria, `INSERT`/`UPDATE` dentro de `SAVEPOINT` | (1, 1600), (NULL, 1600) y ausente con 1600: rechazadas. (0, 1600): aceptada. Sin porcentaje: las 3 aceptadas. Merge Patch sobre (0, 1600): `vat_included = NULL` y `= 1` rechazados; `= 1` con porcentaje NULL, aceptado. Es lo mismo que imprime C1-04 en v-4 |
| v-2 | `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py` con la comprobación nueva y la salida anterior | rc 1. «FALLA C1-04 cada combinación … que el CHECK de la 00004 rechaza (alta y Merge Patch) está en la regla 422 …», «lo que falta en los textos: … (20)» y «RESULTADO: 1 falla(s)». Las otras 32 seguían OK |
| v-3 | `python3 -I backend/internal/api/c1gen/build.py`. Después, las salidas anteriores (regeneradas con `--out datos/tmp/cc1r2/antes` usando el respaldo de `gen2_schemas.py`), comparadas con `cmp` | rc 0. Cambian solo `openapi.yaml` (4 líneas: 6168, 6330, 15487 y 15589), `extractos/ajustes-c.yaml` (4 líneas) y `extractos/INDICE.md` (1 línea: «77.2 KB (79094 bytes)» en la versión final). Los otros 18 archivos quedan idénticos |
| v-4 | `verify_c1.py` (sin `--db`, migraciones en memoria) | rc 0. «C1-04 IVA …: alta con 1600 {'true': 'rechaza', 'null': 'rechaza', 'false': 'acepta', 'ausente': 'rechaza'} · sin porcentaje {todas 'acepta'} · Merge Patch sobre (false, 1600) {true: 'rechaza', null: 'rechaza', true + null: 'acepta'} · textos revisados: 8 · operaciones sin respuesta 422: [] · lo que falta en los textos: [] (0)». 33 OK. «RESULTADO: TODO OK» |
| v-5 | `python3 -I backend/internal/api/c1gen/build.py --check` | rc 0. 21 «IDÉNTICO». «RESULTADO: los 21 archivos coinciden byte a byte con el generador» |
| v-6 | `make verify-c1 GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP_ROOT=…/datos/tmp/cc1r2/tr` (base de goose, con `--db` solo lectura) | rc 0. «RESULTADO: los 21 archivos coinciden byte a byte con el generador», «comprobaciones OK: 33» (antes 32, más C1-04) y «RESULTADO: TODO OK» |
| v-7 | `make verify-openapi` | «docs/contratos/openapi.yaml: OK» (openapi-spec-validator 0.9.0) |
| v-8 | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP_ROOT=…/datos/tmp/cc1r2/va` | rc 0. «== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». C1 con 33 OK |
| v-9 | Lecturas de `verify_c1.py` (v-4) y `grep` | C1-01: «rutas = 168 · operaciones = 212». C1-13: «info.version: 1.3.0». EXT: 212/212 y «> 80 KB: []». S6-03: 212/212. `x-provisional: true`: 6 (sin cambio). `grep` de la redacción anterior («vat_included = true junto con», «vat_included = true no admite») fuera de `datos/` y `.git`: 0 archivos |

Usé un `TMP_ROOT` propio en todas las corridas de `make`, por B-4 de la revisión. Ninguna cifra se ajustó a mano.

## Supuestos
- **S-1.** «Ausente» en el alta equivale a NULL en la base, porque `credit.vat_included` no tiene `DEFAULT`. C1-04 lo prueba con un `INSERT` sin la columna.
- **S-2.** «Tampoco deduce false a partir del porcentaje» no es una decisión nueva. Es la frase de r1 «ni cambia vat_included», que la revisión pidió conservar, aplicada al alta.
- **S-3.** En el Merge Patch, la regla se evalúa sobre la combinación que quedaría guardada, como en r1. Por eso también cubre `{"vat_on_interest_bp": 1600}` sobre un crédito con `vat_included` null: es la misma fila del CHECK.

## Bloqueos
Ninguno nuevo. No cambia el estado de `C1-reparto-alcance`, `C1-hora-reportes-job` ni `S3-iva-lectura`, que siguen abiertos.

## Pendientes
- **Bajos de `REVISION-C1-r2.md`, no encargados:**
  - B-1: falta la fila «C1 v1.3 · corrección r1» (y ahora también «r2») en Cambios de `c1_plantilla.md`, y una frase del 422 en el punto de IVA de C1-datos.md l. 36;
  - B-2: el resumen del orden de los grupos;
  - B-3: qué hace la API con claves de v1.2 en `params_json`;
  - B-4: el TMP compartido del `Makefile`;
  - B-5: los bajos de r1, que siguen abiertos.
- **`C1-puerta2.md`** todavía cita «comprobaciones OK: 29». Ahora son 33. No es mi archivo.

## Solicitudes fuera de alcance
1. **S6 (implementación de la API).** Validar antes de escribir y responder 422 `validation_failed` con `fields[].field = vat_on_interest_bp` cuando la combinación que quedaría guardada trae `vat_on_interest_bp` con valor y `vat_included` distinto de false. Eso incluye:
   - (a) `vat_included` = true (lo que ya se pidió en r1);
   - (b) **`vat_included` ausente en el alta**, por ejemplo `POST /api/v1/credits` con `{"name": "Tarjeta X", "kind": "store_credit", "holder_person_id": 1, "payment_frequency": "monthly", "vat_on_interest_bp": 1600}`;
   - (c) **`vat_included` null**, por ejemplo `PATCH /api/v1/credits/{id}` con `{"vat_included": null}` sobre un crédito con (false, 1600), o `{"vat_on_interest_bp": 1600}` sobre uno con `vat_included` null.

   En ningún caso la API cambia el otro campo ni deduce `vat_included` = false. Así la violación del CHECK de la 00004 nunca sale como 5xx. Fuente: C1 v1.3, `CreditCreate`, `CreditUpdate`, `createCredit` y `updateCredit`.
2. **Integrador o agente de criterios.** Al pasar `docs/criterios/C1.md` a 1.3, registrar en C1-04 la comprobación nueva: las combinaciones del CHECK de `credit.vat_included` probadas contra la base y exigidas en la regla 422.
3. **Integrador.** Corregir el conteo de `C1-puerta2.md` (29 → 33), o anotar que lo cambiaron las correcciones r1 y r2.

## Propuestas para Axel
- **PREGUNTA (no bloquea; hoy no se hace).** ¿Quieres que la API complete `vat_included` = false cuando llega un porcentaje sin ese dato? Así funcionaba la captura en v1.2.

  | | Resultado | Costo | Riesgo |
  |---|---|---|---|
  | Como hoy (recomendado) | La API rechaza con un aviso y la pantalla pide el dato | Ninguno | Ninguno: nada se guarda sin que tú lo confirmes |
  | Que lo complete | Menos pasos al capturar | Una regla nueva en C1 y en S6 | Guarda «no incluye IVA» sin que lo hayas dicho |

  Recomiendo dejarlo como está.

## Consumo aproximado
- 22 acciones de herramienta.
- Unos 75 KB leídos: la revisión (15 KB), el reporte de r1 (12 KB), fragmentos del generador, de `verify_c1.py`, de `parse_schema.py` y de la 00004, y las salidas de las verificaciones. `openapi.yaml` solo se leyó con `python3 -I` y `grep` por fragmentos.
- Unos 18 KB escritos a mano:
  - unos 2 KB en `gen2_schemas.py`;
  - unos 5 KB en `verify_c1.py`;
  - 1.4 KB en el README;
  - unos 9 KB en este reporte.

  Las salidas regeneradas las escribe `build.py`.
- Respaldos de las fuentes antes de editar: `/tmp/{gen2_schemas.py,gen3_ops.py,verify_c1.py,README.md}.bak-cc1r2`. La carpeta temporal `datos/tmp/cc1r2/` se borra al terminar.
