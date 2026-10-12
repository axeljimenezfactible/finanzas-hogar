# Corrección de C1 v1.3 (puerta 2) · ronda 1

Agente: `correccion-C1-r1`, flujo P2c, 11-oct-2026, noche (el reloj del equipo ya marcaba 12-oct 00:1x). Corrige los 3 hallazgos medios de `docs/estado/p2c/REVISION-C1-r1.md` (H-1, H-2 y H-3). No había hallazgos altos. Los bajos (B-1 a B-5) no venían en el encargo y no los toqué. Tampoco decidí ninguna pregunta «sin recomendación».

Alcance tocado: solo `backend/internal/api/c1gen/` (`gen2_schemas.py`, `gen3_ops.py`, `c1_plantilla.md`, `verify_c1.py` y `README.md`) y las salidas que regenera `build.py` (`docs/contratos/openapi.yaml`, `C1-datos.md` y `extractos/`). `routes.go` se regeneró idéntico. Sin `git add` ni `git commit`.

## Hecho

### H-1 (media) · nombres con un guion bajo de más · CORREGIDO
- **Causa confirmada.** `O()` usa el `_` final como marca de opcional y lo quita. `_ins` no lo quita: copia el nombre tal cual. Por eso `visit_date_=` y `cheapest_store_id_=` se publicaban con el guion bajo.
- **Corrección.** En `gen2_schemas.py`, los argumentos pasan a `visit_date=` y `cheapest_store_id=`. Ninguna de las dos estaba en `required` antes ni lo está ahora: `_ins` no toca `required`, así que siguen opcionales, como dicen sus descripciones («Solo en grupos de mercado con tienda» y «Ausente si ninguna está completa»).
- **Para que no se repita.** `_ins` ahora rechaza con `assert` cualquier nombre que termine en `_`. Prueba negativa: `('_ins: nombre con «_» final', ['visit_date_'])`.
- **Comprobaciones nuevas en `verify_c1.py`.**
  - C1-14: 0 propiedades que terminen en `_` en `openapi.yaml` y en cada extracto.
  - C1-14: cada nombre entre comillas invertidas de «C1 v1.3 (puerta 2)» de C1-datos.md existe en `openapi.yaml` con su tipo: esquema, `operationId`, ruta `Esquema.prop[].sub` o `prop[].sub`, clave de `app_setting`, propiedad, valor de `enum`, tabla o ruta `/api/…`.
  - Dos excepciones, ambas explícitas: `send_time` y `advance_days`, que la sección menciona como retirados (C1-23 comprueba que ya no estén), y `scopes`, la clave propia de `params_json` que se nombra en su descripción.
- **Las comprobaciones detectan el defecto.** Corridas contra la salida anterior, las dos fallaron con exactamente los nombres del hallazgo (ver Verificación).
- **Resultado.**
  - `openapi.yaml` l. 13517 dice `visit_date:` y l. 13718 dice `cheapest_store_id:`.
  - Cambian igual `extractos/compras-tiendas-a.yaml` (2 nombres) y `extractos/reportes.yaml` (1 nombre).

### H-2 (media) · M-11 no nombraba los acumulados · CORREGIDO
- **Fórmula.** La fila M-11 de `c1_plantilla.md` (que sale en C1-datos.md l. 372) dice ahora:
  - al cierre de cada mes, con el mismo origen para los tres (desde el primer evento de las entradas, S2-interfaz §3.2; Q2-25 A y Q2-65 A):
    - `cumulative_accrued_cents` = Σ devengado;
    - `cumulative_cash_cents` = Σ efectivo;
    - `cumulative_difference_cents` = devengado acumulado − efectivo acumulado;
  - el origen no es el primer mes visible.
- **Ejemplo**, con cifras sumadas de `periodos[i]` de `oraculo/model.json` (Decimal), no copiadas a mano:

| Mes | Devengado acumulado | Efectivo acumulado | Diferencia acumulada | Referencia |
|---|---|---|---|---|
| oct-2026 | $1,947.50 | $973.75 | $973.75 | R-45 |
| nov-2026 (primer mes de una ventana que empieza ahí) | $3,895.00 | $2,921.25 | $973.75 | R-46: el mes trae $1,947.50 y $1,947.50 |
| mar-2027 | $14,362.48 | $13,103.78 | $1,258.70 | R-51, R-52 y R-50 |

  La fila conserva el contraste: sumar solo lo visible daría $1,947.50 y $1,947.50.
- **Comprobación nueva C1-09 en `verify_c1.py`.**
  - Recalcula desde `model.json` los acumulados de P1, P2 y P6.
  - Exige que las 9 cifras y los 3 nombres de campo estén en la fila M-11.
  - Comprueba que Σ devengado − Σ efectivo = `diferencia_acumulada` en los 6 periodos.

### H-3 (media) · qué responde la API a «incluye IVA» con porcentaje · CORREGIDO
- **En los esquemas.** `gen2_schemas.py`, bloque de `vat_included` (Q2-67 opción 1), con la constante `CREDIT_VAT_RULE`. Va en las descripciones de `CreditCreate` y `CreditUpdate`:
  - si la combinación que quedaría guardada trae `vat_included` = true y `vat_on_interest_bp` con valor, la API responde 422 `validation_failed` con detalle en `fields[]` (`field = vat_on_interest_bp`) y no guarda nada;
  - la combinación que quedaría guardada es, en el alta, la del cuerpo, y en el Merge Patch, la del cuerpo sobre lo ya guardado. Ejemplos: `{"vat_included": true, "vat_on_interest_bp": 1600}`, o `{"vat_included": true}` sobre un crédito que ya tiene 1600;
  - la API no cambia el otro campo por su cuenta;
  - para pasar a «la tasa ya incluye IVA» se manda `{"vat_included": true, "vat_on_interest_bp": null}`.
- **En las operaciones.** `gen3_ops.py`: `createCredit` y `updateCredit` llevan `CREDIT_VAT_RULE_OP`, una versión corta de la misma regla (422, `fields[].field`, sin cambiar el otro campo) que remite a los esquemas.
- **Por qué es corta.** Con el texto completo repetido en las dos operaciones, `extractos/ajustes-c.yaml` quedaba en 79,426 bytes, muy cerca del tope de 80,000. Con la versión corta queda en 78,546.
- **Lo que no toqué.** El texto de `vat_included` y `vat_on_interest_bp` (viene de la nota de columna de S1) y las marcas `x-provisional`, que siguen siendo 6. Son B-1 y B-2, bajos y no encargados.

### Documentación del generador
`c1gen/README.md` describe las 3 comprobaciones nuevas y la guarda de `_ins`.

## Verificación (comandos y salidas)

| Id | Comando | Salida |
|---|---|---|
| v-0 | `goose -dir backend/db/migrations sqlite3 datos/tmp/cc1r1/c1.db up` | «successfully migrated database to version: 4» |
| v-1 | Comprobaciones nuevas contra la salida anterior: `uv run … verify_c1.py --db datos/tmp/cc1r1/c1.db` | rc 1. «que terminan en «_»: openapi ['cheapest_store_id_', 'visit_date_'] · extractos ['compras-tiendas-a.yaml:cheapest_store_id_', 'compras-tiendas-a.yaml:visit_date_', 'reportes.yaml:cheapest_store_id_']» y «que no existen en openapi.yaml: ['BasketComparison.cheapest_store_id', 'groups[].visit_date']». «RESULTADO: 2 falla(s)» |
| v-2 | `python3 -I backend/internal/api/c1gen/build.py` | rc 0 |
| v-3 | Diferencias entre las salidas anteriores (regeneradas con las fuentes de antes en `--out datos/tmp/cc1r1/antes`) y las nuevas | `openapi.yaml`: 6 líneas cambiadas. Son l. 6168 (`createCredit`), l. 6330 (`updateCredit`), l. 13517 (`visit_date`), l. 13718 (`cheapest_store_id`), l. 15487 (`CreditCreate`) y l. 15589 (`CreditUpdate`). `C1-datos.md`: solo l. 372 (M-11). `routes.go`: 0. Extractos: `ajustes-c.yaml` (4 descripciones), `compras-tiendas-a.yaml` (2 nombres), `reportes.yaml` (1 nombre) e `INDICE.md` (los tamaños de esos tres) |
| v-4 | `make verify-c1 GOOSE=… TMP_ROOT=datos/tmp/cc1r1/tr` | rc 0. «RESULTADO: los 21 archivos coinciden byte a byte con el generador», «comprobaciones OK: 32» (antes 29, más 3 nuevas) y «RESULTADO: TODO OK» |
| v-5 | `verify_c1.py --db datos/tmp/cc1r1/c1.db` (salida completa) | rc 0, 32 OK y 0 FALLA. C1-01: «rutas = 168 · operaciones = 212». C1-13: «info.version: 1.3.0». C1-14: «propiedades distintas: 610 · que terminan en «_»: openapi [] · extractos []» y «nombres … 77 · comprobados: 66 · que no existen en openapi.yaml: []». C1-09: «faltan en la fila: [] · campos que no nombra: [] · periodos donde Σ devengado − Σ efectivo ≠ diferencia_acumulada: []». EXT: 212/212, «> 80 KB: []». S6-03: 212/212 |
| v-6 | Lectura del YAML con `python3 -I` | `visit_date` y `cheapest_store_id` existen, no son obligatorias y quedan en su lugar (después de `visit_weekday` y de `stores`). `CreditCreate`, `CreditUpdate`, `createCredit` y `updateCredit` traen «422 validation_failed» y `vat_on_interest_bp` como campo del detalle, y las dos operaciones conservan la respuesta 422. `x-provisional`: 6. `grep` de `visit_date_` y `cheapest_store_id_` en `openapi.yaml` y en los extractos: 0 |
| v-7 | Prueba negativa de `_ins` (la función aislada) | «rechazado: ('_ins: nombre con «_» final', ['visit_date_'])» |
| v-8 | `make verify-openapi` | «docs/contratos/openapi.yaml: OK» |
| v-9 | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP_ROOT=datos/tmp/cc1r1/va` | rc 0. «== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)» |
| v-10 | `build.py --check` después de editar el README | «RESULTADO: los 21 archivos coinciden byte a byte con el generador» |

Ninguna cifra se ajustó a mano. Las de M-11 salen de `model.json` y las recalcula C1-09 en cada corrida.

## Supuestos
- **S-1 (H-3, ambos sentidos).** La regla cubre cualquier combinación que quedaría guardada. También cubre `{"vat_on_interest_bp": 1600}` sobre un crédito con `vat_included` = true, porque es el mismo CHECK de la 00004. Es la regla que sugiere la revisión («en el cuerpo o ya guardado») más la convención de C1 de no sustituir valores (`SETW`, `scope` sin `parts`). No es una decisión nueva.
- **S-2 (H-3).** El ejemplo con `vat_on_interest_bp: null` no es una regla nueva: aplica el Merge Patch ya documentado en `CreditUpdate` («null borra el valor»).
- **S-3 (H-2).** Uso «desde el primer evento de las entradas» con las palabras de S2-interfaz §3.2 (Q2-25 A, FIRME). Equivale al origen P1 de `model.json` para el canon.
- **S-4 (H-1).** Las dos propiedades siguen siendo opcionales. Ese era el sentido del `_`, y antes tampoco eran obligatorias.

## Bloqueos
Ninguno nuevo. No cambió el estado de los que ya estaban abiertos (`C1-reparto-alcance`, `C1-hora-reportes-job` y `S3-iva-lectura`).

## Pendientes
- **B-1 a B-5 (bajos) de `REVISION-C1-r1.md`.** No venían en este encargo:
  - B-1: marca de más en `vat_included`;
  - B-2: vocabulario de la base en su descripción, que viene de la nota de columna de S1;
  - B-3: anotar en `S3-iva-lectura` lo que tomó C1;
  - B-4: `C1.md` a 1.3;
  - B-5: agregar los bloqueos nuevos a `bloqueos/INDICE.md`.
- **El reporte del implementador.** `docs/estado/p2c/C1-puerta2.md` l. 84 cita «comprobaciones OK: 29». Ahora son 32. No es mi archivo.
- **Lo que queda en P2d.** K53: lectura asistida e importación del agente externo, como dice la revisión.

## Solicitudes fuera de alcance
1. **S6 (implementación de la API).** Validar `vat_included` = true junto con `vat_on_interest_bp` con valor antes de escribir: en el alta, con el cuerpo; en el PATCH, con el cuerpo aplicado sobre lo guardado. Responder 422 `validation_failed` con `fields[].field = vat_on_interest_bp`, para que la violación del CHECK de la 00004 nunca salga como 5xx (C1 v1.3, `CreditCreate`, `CreditUpdate`, `createCredit` y `updateCredit`).
2. **Integrador o agente de criterios.** Al pasar `docs/criterios/C1.md` a 1.3 (B-4), conviene registrar las 3 comprobaciones nuevas en C1-09 y C1-14: el nombre de propiedad sin `_` final, los nombres de la sección v1.3 que existen en el contrato, y M-11 recalculada desde el oráculo.
3. **Integrador.** Corregir el conteo de `C1-puerta2.md` (29 → 32), o dejar la nota de que lo cambió esta corrección.

## Propuestas para Axel
- Ninguna decisión nueva hace falta para estos tres hallazgos.
- Propuesta técnica para una ronda futura, sin urgencia: extender la comprobación de nombres de C1-14 a todo C1-datos.md, no solo a la sección «C1 v1.3 (puerta 2)». Atraparía el mismo tipo de error en cualquier versión.

## Consumo aproximado
- 38 acciones de herramienta.
- Unos 95 KB leídos: revisión, fragmentos del generador y de la plantilla, bloqueo, salidas de las verificaciones, y diferencias acotadas. `openapi.yaml` solo se leyó con `python3 -I` por fragmentos.
- Unos 22 KB escritos a mano: unos 6 KB en `verify_c1.py`, 2.5 KB en `gen2_schemas.py`, 0.3 KB en `gen3_ops.py`, 1.6 KB en `c1_plantilla.md`, 0.7 KB en el README y unos 10 KB de este reporte. Las salidas regeneradas las escribe `build.py`.
- Respaldos de las fuentes antes de editar: `/tmp/{gen2_schemas.py,gen3_ops.py,c1_plantilla.md,verify_c1.py}.bak-cc1r1`. Carpeta temporal `datos/tmp/cc1r1/`, borrada al terminar.
