# Corrección del contrato S2 · ronda 2

Agente: `correccion-contrato-S2-r2` · 11-oct-2026 · Alcance: `docs/contratos/motores/S2-interfaz.md`, `oraculo/fixtures/FORMATO.md`, este reporte (más un bloqueo nuevo de S2, como permite el protocolo §2). Vara: `docs/estado/wf1/REVISION-contrato-S2-r2.md` (6 hallazgos medios M1–M6), `docs/criterios/S2.md` 1.2 y decisiones §B y §K. Ninguna cifra canon cambia y no se reabrió ninguna decisión.

## Hecho
- **M1 · argumentos (S2-03).** Agregué en §3 una regla general: claves en inglés, `{}` si no lleva, `unknown_field` y `invalid_argument`, «exactamente uno», y los datos salen siempre de `entradas.datos`. Cada función de §3.1–§3.14 tiene ahora su línea «Argumentos `{…}` → salida» con obligatoriedad: `periodos {months}|{windows}`, `horizonte {start, end, months}`, `deuda_viva {cut}`, `promedios {total_cents, start, end}|{month, basis}` (con el total de cada `basis`), `promedio_mensual {total_cents, month_count}`, `valor_estable {items?}`, `papel {}`, `presupuesto_mes {month}`, `aplicar_cambio_presupuesto {row_id, month, amount_cents, scope?}` (los renglones salen de `entradas.datos.budgets`), `plan_contra_real {from, count, category?}`, `flujo_quincenal {from, count}`, `proyeccion {from?}`, `metas {}` y las 6 vistas. `plan_contra_real` ya tiene función Python: `fin_presupuesto.py: plan_contra_real(datos, desde, n, categoria=None)`. Para que los mapas de M2 tengan sus personas, la firma de `valor_estable` pasa a recibir `persons`.
- **M2 · mapas `{code: cents}`.** Regla única en §1.1: los `by_person` y los `shares` de las salidas llevan siempre al titular y a cada persona externa de `persons`, con 0 si no tiene parte, y nunca a la co-titular. Los `by_category` llevan solo las categorías con importe distinto de 0, y quien lee una categoría ausente toma 0. Van ejemplos: `papel_contado` → `{axel: 0, karen: 9950}`, y una compra solo del titular → `karen: 0`. En §1.4 agregué la validación de personas por C-09 y en FORMATO §4 la regla 5 del recorte.
- **M3 · campos sin fórmula (S2-11, S2-12).** En §3.8, cada campo tiene fórmula: categoría con y sin renglón (`no_plan`: `plan_cents`, `remaining_cents` y `used_pct` en null, banderas en `false`), `total.actual_cents` (todas las filas, como dice M-01), `total.remaining_cents`, y `accrued_household_cents` y `cash_household_cents` (`by_person[holder]` de `periodos`, todas las categorías; con la semilla en oct-2026, 119764 = ejemplo de M-01 y 59882). En §3.13: `month_summary` (`saved_cents`, `contributed_this_month_cents` = solo depósitos del mes de hoy con fecha ≤ hoy, `planned_this_month_cents`) y las metas sin `target_cents` o sin `target_date` (null y `plan` vacío). Queda una sola duda real, si el resumen cuenta solo las metas activas, y va PROVISIONAL con el bloqueo nuevo `S2-metas-resumen` (§8).
- **M4 · `projected` en `Expense`.** La agregué a la tabla de §2.4: bool, por omisión `false`, solo `true` dentro de `cadence`, como en `CreditPurchase`.
- **M5 · autoría de las fixtures.** Encabezado y §5.0 nuevo: el agente Python de S2 es el autor único de las 27 fixtures y de los 3 archivos de `_datos/`, con el script `oraculo/fin_genera_fixtures.py`. §5.0 también fija de dónde sale `esperado` y qué hace el agente Go mientras falten esos archivos: el ayudante `seedFromModel` arma el mismo `Input` con la receta; las pruebas que leen archivos fallan con nombre y autor, nunca con `t.Skip`; y cuando el archivo existe, `TestCanon` exige semilla = `seedFromModel`. §5.3 nuevo trae la receta completa de `semilla-oraculo.json`: las 12 claves, incluidos `stable_items`, los códigos de `groceries_category` y `paper_category`, y el `kind` y los `id` de los 6 gastos. FORMATO §1 exige que cada interfaz declare sus autores. §7 y §9 quedaron al día.
- **M6 · `verify3.py:89`.** El cambio queda asignado al agente Python de S2 (encabezado y §0.5): solo la ruta de salida pasa a `datos/tmp/verify3.json`, sin tocar el cálculo ni `oraculo/verify3.json`. La verificación queda escrita: `grep` en 0 líneas, `cmp` contra `oraculo/verify3.json` (sha256 `d86087f3…5041`) y sha256 de `model.json`. Comprobé hoy con una copia redirigida que la salida es igual byte a byte. §7 cita la verificación de S2-04.
- Versiones: `S2-interfaz.md` 1.2 y `FORMATO.md` 1.2, cada una con su nota «Cambios de la 1.2».

## Archivos
- `docs/contratos/motores/S2-interfaz.md` (modificado; 63.5 → 77.9 KB; 67 líneas nuevas o cambiadas)
- `oraculo/fixtures/FORMATO.md` (modificado; 14.8 → 15.8 KB)
- `docs/estado/bloqueos/S2-metas-resumen.md` (nuevo)
- `docs/estado/wf1/correccion-contrato-S2-r2.md` (este reporte)

## Verificación
Scripts en `datos/tmp/corr-s2-r2/` (ya borrados).
- `sha256sum oraculo/model.json` → `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`, sin cambio antes y después. `oraculo/verify3.json` → `d86087f3…5041`, sin cambio.
- `python3 -I semilla_check.py oraculo/model.json`: arma la semilla con la receta de §5.3 → `validaciones §2.4/§2.5: 28 de 28`; con las reglas de §1.3, §3.2, §3.3, §3.4, §3.6 y §3.7 → **`hojas S2 desde la receta: 159 de 159 coinciden`** (122 de `periodos`, 12 de `valor_estable`, 13 de `horizonte` y 12 de `papel`). La receta alcanza para las fixtures canon (M5).
- Copia de `verify3.py` con la salida redirigida a `datos/tmp` → `cmp` contra `oraculo/verify3.json`: iguales byte a byte (M6). `grep -n reporte-compra-latex oraculo/*.py` sigue dando 1 línea (`verify3.py:89`): la corrección le toca al agente Python de S2 en la implementación, no a este contrato.
- `python3 -I verifica.py`:
  - M1: 0 subsecciones de §3.1–§3.14 sin «Argumentos», y la función Python de `plan_contra_real` está presente.
  - §4: 159 rutas, 159 distintas, todas son hojas de `model.json`, 0 renglones con `n` mal, 0 fuera de las ramas de S2.
  - Presencia de M2 (§1.1 y §1.4), M3 (§3.8 y §3.13), M4, M5 (§5.0 y §5.3), M6 (§0.5 con el sha256), FORMATO 1.2 y el bloqueo nuevo: 10 de 10 OK.
  - Tablas Markdown: 0 renglones con columnas distintas en los dos archivos.
- Las cifras de los ejemplos nuevos salen de `model.json`: 119764 = `periodos[0].devengado.axel`, 59882 = `periodos[0].efectivo.axel`, 194750 = `periodos[0].devengado.despensa_total`, 9950 = `valor_estable.componentes.papel_contado_mes`, y 503474 / 149972 / 353502 / 19900 de `despensas[1]` y `parametros.papel`.
- Criterios de la fase, según la revisión r2:
  - S2-03: M1, M2 y M5 cerrados (§1.1, §3, §5.0, §5.3).
  - S2-04: M6 asignado, con su verificación (§0.5, §7).
  - S2-09: M4 cerrado (§2.4).
  - S2-11 y S2-12: M3 cerrado (§3.8 y §3.13; lo no decidido va PROVISIONAL con `S2-metas-resumen`).
  - Las demás filas de criterios no cambiaron. No hay pruebas Go ni oráculo que correr: el motor y `fin_*.py` aún no existen; este artefacto es el contrato.

## Supuestos
- Los códigos `despensa` y `papel` de la semilla son identificadores de la fixture: S1 y S6 aún no siembran categorías, y ningún canon depende del código. El `kind` `variable` de las 2 categorías y de los 6 gastos no cambia ninguna cifra: los 4 gastos a crédito no cuentan como gasto en el flujo y los 2 de papel los paga Karen. Así lo comprueban las 159 de 159.
- Con base en C-09 (el dinero del hogar Axel–Ximena va a nombre del titular), una parte, un `paid_by`, un abono, una persona de `share_rule` o de `stable_items` a nombre de la co-titular, o un `code` que no está en `persons`, da `invalid_argument`. El `person` de una partida es atribución y no se valida así. Es validación, no cifra.
- Elegí códigos de error entre los que ya existen: `row_id` inexistente, `purchase_id` inexistente en `vista_interes_por_monto`, y categoría inactiva o inexistente en `plan_contra_real` → `invalid_argument`.
- `total.actual_cents` suma también las categorías activas sin renglón, porque M-01 dice «real = Σ parte del hogar devengada». `accrued_household_cents` incluye las categorías inactivas: es el «real» literal de M-01, y su ejemplo 1,197.64 es `periodos[0].devengado.axel`. Se conserva la regla anterior de §3.8: una categoría inactiva no suma a `total`.
- «Aportado del mes» = solo depósitos: M-07 distingue «aportaciones» de «retiros».
- La firma `StableValue(items, persons)` y `valor_estable(partidas, personas)` reemplaza a la versión con `holder`. Nadie la implementó aún: `backend/internal/motor_fin` y `oraculo/fin_*.py` no existen.

## Bloqueos
- `docs/estado/bloqueos/S2-metas-resumen.md` (nuevo): ¿el `month_summary` cuenta solo las metas activas o también las inactivas con saldo? Lectura de trabajo: A (activas), marcada PROVISIONAL.
- Siguen abiertos sin cambio: `S2-proyeccion-interruptor`, `S2-flujo-quincenal`, `S2-metas` y `S2-presupuesto-total`. `S1-parametros-avisos` no afecta a S2.

## Pendientes
- Hallazgos bajos B1–B5 de la r2: no estaban en mi encargo y no los toqué. B1: falta afirmar en `TestProyeccionDefecto` «0 eventos entre los días 1 y 9». B2: falta fixture de `plan_contra_real` y de `promedios {month, basis}`, falta decir el mes de `s2-mes-30-dias` y el error de un evento fuera de las ventanas. B3: `fortnightly` con `start_date` que no cae en día de pago, y `today` y `now` a la vez. B4: un `variable` con fecha futura en el flujo. B5: FORMATO §5, contraste por las hojas de `esperado`.
- La implementación aplica esta interfaz: el agente Python con `oraculo/fin_*.py`, `oraculo/fin_genera_fixtures.py`, `oraculo/fixtures/s2/**`, `oraculo/compara.py` y `oraculo/verify3.py:89`; el agente Go con `backend/internal/motor_fin`, `backend/cmd/motores/main.go`, `s2.go` y `seedFromModel`.

## Solicitudes fuera de alcance
1. Orquestador, prompts de implementación de S2. Asignar al agente Python, además de `fin_*.py` y `compara.py`: `oraculo/fin_genera_fixtures.py`, todo `oraculo/fixtures/s2/` (27 casos y 3 `_datos/`) y `oraculo/verify3.py` (solo la ruta de la línea 89). El agente Go no escribe en `oraculo/`. Motivo: S2-interfaz 1.2, encabezado, §0.5 y §5.0 (M5, M6).
2. `docs/contratos/motores/S3-interfaz.md`, `S4-interfaz.md` y `S5-interfaz.md`. Citar FORMATO 1.2 y declarar el autor de cada archivo de su `oraculo/fixtures/<motor>/` (FORMATO §1), las claves de `argumentos` de cada función (§2) y la regla de sus mapas por persona o categoría (§4, regla 5). Motivo: la 1.2 agrega esos tres requisitos comunes.
3. `docs/contratos/C1-datos.md` M-07. Cuando Axel decida `S2-metas-resumen`, escribir si «ahorrado» y «aportado del mes» cuentan solo las metas activas. Motivo: M3.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
- `S2-metas-resumen` (puerta 2): propongo la opción A, solo metas activas, el mismo conjunto que la lista de metas. Así la tarjeta de Inicio y la pantalla de Metas suman lo mismo, y archivar una meta la saca de los dos lados. Si se quiere conservar a la vista el ahorro de una meta archivada, la opción B lo hace.
- Sin valores K31 en este encargo.

## Consumo aproximado
- 31 acciones de herramienta.
- Lectura ≈ 195 KB: protocolo, revisión r2, §K, `S2.md` completo, `S2-interfaz.md` por partes, FORMATO, extractos de reportes, C1 y `model.json`.
- Escritura ≈ 62 KB de contenido nuevo: scripts temporales ≈ 36 KB, ediciones ≈ 15.5 KB, bloqueo 1.7 KB y reporte ≈ 9 KB. Los dos documentos editados se reescribieron completos por script (78 KB + 16 KB).
