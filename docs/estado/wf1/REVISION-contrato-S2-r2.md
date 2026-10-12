# Revisión del contrato S2 · ronda 2 de 3

Revisor: `revisor-contrato-S2-r2` (independiente; no escribió nada de lo revisado y no corrige) · 11-oct-2026 · Artefactos: `docs/contratos/motores/S2-interfaz.md` 1.1 y `oraculo/fixtures/FORMATO.md` 1.1 · Vara: `docs/criterios/S2.md` 1.2, decisiones §B y §K.

## Veredicto
**NO APROBADO**: 0 hallazgos altos, 6 medios y 5 bajos. Los 10 hallazgos medios de la r1 quedan **cerrados** (los comprobé con cálculos propios). Las cifras están bien: mi cálculo independiente reproduce las 159 hojas de S2, los 43 cruces con S3 y las 7 filas de K19, con **242 OK y 0 FALLA**. Los 6 medios nuevos son preguntas que un implementador tendría que adivinar, más un hueco de alcance que deja sin dueño una corrección exigida por S2-04. Ninguno cambia cifras.

## Verificación propia (scripts en `datos/tmp/rev-s2-r2/`, ya borrados)
- `sha256sum oraculo/model.json` → `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` (sin cambio).
- `python3 -I hojas.py oraculo/model.json`: 596 hojas numéricas. Por rama: parametros 9, periodos 122, despensas 144, credito 143, valor_estable 12, sensibilidad 45, papel 12, disposicion_por_despensa 68, horizonte 13, interes_por_monto 28. **S2 = 159** (periodos[0] 17 y periodos[1..5] 21 cada uno), S3 = 437.
- `python3 -I tabla4.py`: expandí las llaves de §4 con mi propio analizador. Resultado: 20 renglones, 159 rutas, 159 distintas, 0 «n» mal declarado, **0 sin dueño, 0 sobrantes, 0 rutas en ramas de S3** (ninguna duplicada con S3).
- `python3 -I calc.py oraculo/model.json`: cálculo propio con `Fraction`, que solo toma como entrada las hojas S3 (`despensas`, `papel`, `parametros`) y aplica las reglas de la interfaz (mes calendario, devengo §1.3, vistas §3.14). Resultado: **159 de 159 hojas S2 iguales a `model.json`**; `sensibilidad` 15 de 15; `interes_por_monto` 28 de 28, con la serie de §5.2 (Σ devengado = Σ efectivo; capital e interés ≥ 0 por pago). También coinciden:
  - las 7 filas de K19: 292125 = 97375 + 194750 … 629344 = 125870 + 503474, con la identidad cierta en los 7 cierres;
  - los promedios: 3141/21988, 6492/45442, 8991/62935 y 6282/43976, más los mensuales 218396 y 239375;
  - el flujo: 59882 ×4 y Karen 37493 ×16; en los ciclos 2 a 4, 353502 = 2 × 176751 y 149972 = 2 × 74986; las diferencias por quincena salen +1 ×9 y −1 ×3;
  - el calendario: 2026-10-31 cae en q2 (2026-10-16 a 2026-10-31, paga el 30); la q2 de febrero va del 2027-02-16 al 2027-02-28; 2026-10-14 es la semana `2026-W42` (12 al 18); con zona horaria, 05:59:59Z → 2026-10-30 y 06:00:00Z → 2026-10-31; octubre + 50000 = 147375;
  - el presupuesto (88 %, aviso sí, excedido no; 79 999 → 80 con aviso no) y la meta (500 décimas, 50 %).
  - Total: **242 OK, 0 FALLA**.
- `python3 -I insumo.py`: la tabla de §5.1 tiene 55 filas únicas, exactamente R-01…R-36, R-61…R-67 y R-74…R-85, y todas sus llaves resuelven. R-67 tiene 85 llaves con los tipos compra, inicio, pago y papel; R-65 y R-66 tienen 14 valores cada una. En R-04, el binary64 de 8005000/309450 es igual a `tasa_pct`; R-05 da HALF_UP 17.27 = `data.json`. El ejemplo R-65 de §5.1 coincide con la conversión.
- `grep -n reporte-compra-latex oraculo/*.py` → `verify3.py:89` sigue escribiendo fuera del repositorio (hallazgo M6).
- C1 `BudgetScope` (`componentes-comunes.yaml`): `this_month` · `onward`, por omisión `onward` = §2.8/§3.9 (K2). RN-03 (FIRME) y el esquema: regla por defecto «Hogar» y Despensa «Conjunto 50/50» = §2.7 (decisión B).

## Criterios
| Criterio | Estado | Evidencia propia |
|---|---|---|
| S2-01 (alta) | Cubierto | §5.1 verificado (55 filas, llaves, R-04/R-05); 30 filas S2 en `TestCanon` (R-37…R-60, R-68…R-73); 159/159 calculadas por mí |
| S2-02 (alta) | Cubierto | Tabla §4 = 159 sin dueño 0; mi cálculo 159/159; FORMATO §2.1 liga `id` ↔ hoja |
| S2-03 (alta) | Cubierto con hallazgos M1, M2, M5 | Fixtures con periodos, valor estable, horizonte, papel, sensibilidad (15/15), interés por monto (28/28) y fin de mes; método distinto en §0.4. Las claves de argumentos, los mapas y la autoría de las fixtures quedan sin fijar |
| S2-04 (alta) | Cubierto con hallazgo M6 | §0.5 y §7 (sha256); nadie tiene asignada la corrección de `verify3.py:89` |
| S2-05 (alta) | Cubierto | §1.3 (K14), §2.5 (verifica sin corregir), invariantes §3.2 y `TestSumaPartes` |
| S2-06 (alta) | Cubierto | §1.2, §3.1, §5.2; mis casos de febrero, 31-oct, semana y 147375 |
| S2-07 (alta) | Cubierto | §3.4; las 7 filas reproducidas con la lectura de meses cerrados |
| S2-08 (alta) | Cubierto | §1.3, §3.5; los diarios y semanales de K39 reproducidos |
| S2-09 (alta) | Cubierto con hallazgos M4 y B1 | §3.12 con fórmulas, interruptor apagado por omisión y la cadencia PROVISIONAL; meses 7–12 = 2027-04-15 (desp. 3) + desp. 4 |
| S2-10 (alta) | Cubierto | §3.11: 16 quincenas y contraste con el valor estable reproducidos |
| S2-11 (alta) | Cubierto con hallazgo M3 | Fronteras exactas en enteros y umbrales `alerts.*`; campos de salida sin fórmula |
| S2-12 (media) | Cubierto con hallazgo M3 | §3.13: Σ plan = restante; `month_summary` sin fórmula |
| S2-13 (alta) | Cubierto | §0.3, §1.1 |
| S2-14 (media) | Cubierto | §1.2, §3.1; comprobado con `zoneinfo` |
| S2-15 (alta) | Cubierto | §7 (vet, gofmt, cobertura ≥ 85 %, K25) |
| S2-16 (media) | Cubierto | §8; existen los 4 bloqueos `S2-*` |

Decisiones: K2, K8, K14, K19, K20/K39, K21 (apagado por omisión, cadencia PROVISIONAL) y K44 (umbrales en `alerts.*`) se cumplen. §B se cumple: quincena 15/30 o último día, devengo mitad por mes, ingresos solo del titular (`income_not_holder`) y reparto según RN-03. No hay contradicciones.

## Hallazgos de esta ronda
### Medios
- **M1 · S2-03 · claves de `argumentos` sin declarar.** Solo se declaran para `calendario` (`{date}`/`{now}`), `valor_estable` (`{items}`), `promedios`, `vista_periodo` (`{month}`, en el ejemplo de FORMATO), `vista_sensibilidad` y `vista_interes_por_monto`. Para `periodos`, `horizonte`, `deuda_viva`, `promedio_mensual`, `presupuesto_mes`, `aplicar_cambio_presupuesto`, `plan_contra_real`, `flujo_quincenal`, `proyeccion`, `vista_horizonte` y `vista_valor_estable` no hay nombres JSON. Quedan preguntas abiertas: ¿`months` o `meses`?, ¿`from`/`count` o `desde`/`n`?, ¿`cut`?, ¿`aplicar_cambio_presupuesto` toma los renglones de `entradas.datos.budgets` o de los argumentos? Además, `plan_contra_real` no tiene función Python. Quien escriba las fixtures y el evaluador Go tendrían que adivinarlo. Ruta: `S2-interfaz.md` §3.2–§3.14. **Corrección:** agregar a cada función una línea «Argumentos `{…}` → salida», con nombres en inglés y obligatoriedad, más el nombre Python de §3.10.
- **M2 · S2-03 · mapas `{code: cents}` sin regla de contenido.** No se dice qué claves llevan `by_person`, `by_category` y `shares`, ni si una persona o categoría en 0 aparece. Ejemplo: en el componente `papel_contado`, ¿lleva `axel: 0`? ¿Y la co-titular? El recorte de FORMATO §4 marca `__falta__` cuando `esperado` trae una clave que la otra implementación omite. Ruta: `S2-interfaz.md` §1.1/§3.2/§3.6/§3.11. **Corrección:** fijar una regla única, por ejemplo: `by_person` lleva siempre al titular y a cada persona externa de `persons` (0 si no tiene parte) y nunca a la co-titular; `by_category`, solo las categorías con importe distinto de 0.
- **M3 · S2-11, S2-12 · campos de salida sin fórmula.** En `presupuesto_mes`: `accrued_household_cents` y `cash_household_cents` (¿todas las categorías o solo las activas? ¿incluye `cash` lo pagado por Karen?); `remaining_cents`, `used_pct`, `warning` y `limit` de una categoría sin renglón (`no_plan`: solo está descrito el plan 0); `total.actual_cents`. En `metas`: `month_summary.saved_cents` (¿activas? ¿fondos?), `contributed_this_month_cents` (¿depósitos o neto?, ¿en qué mes?) y `progress_*` de un `fund` sin `target_cents`. Ruta: `S2-interfaz.md` §3.8, §3.13. **Corrección:** una fórmula por campo, o la marca PROVISIONAL con su bloqueo.
- **M4 · S2-09 · `Expense` sin campo `projected`.** §2.10 define `cadence.expenses` como «Expense con projected = true», pero la tabla de §2.4 no tiene `projected` y §2 declara que un campo desconocido es `unknown_field`. Con eso, la fixture `s2-proyeccion-7-12-cadencia` la rechaza un lector estricto y otro no. Ruta: `S2-interfaz.md` §2.4, §2.10. **Corrección:** agregar `projected` (bool, por omisión `false`, solo `true` dentro de `cadence`) a la tabla de §2.4.
- **M5 · S2-03 (alcance disjunto, protocolo §3) · fixtures sin autor.** Solo `_datos/insumo-s3.json` tiene dueño (el agente Python). Nadie tiene asignadas las 27 fixtures, `_datos/semilla-oraculo.json` ni `_datos/interes-por-monto.json`, aunque `TestCanon` y `TestOraculoFixtures` (Go) las necesitan. Al contenido de la semilla también le falta: `stable_items` (los usan `s2-valor-estable` y la consulta `valor_estable` de `s2-flujo-quincenal`), los códigos de `groceries_category` y `paper_category`, el `kind` y los `id` de los 4 gastos de despensa y de los 2 de papel. Ruta: `S2-interfaz.md` encabezado y §5. **Corrección:** declarar el autor de cada archivo de `oraculo/fixtures/s2/`, qué hace el otro agente mientras no existan, y el contenido completo de `semilla-oraculo.json`.
- **M6 · S2-04 · corrección de `verify3.py:89` sin dueño.** S2-04 exige que en WF-1 deje de haber escrituras a `~/reporte-compra-latex/`; hoy `oraculo/verify3.py:89` todavía escribe ahí. La interfaz cita la regla (§0.5) pero no asigna a nadie el cambio, y el alcance del agente Python (`fin_*.py`, `compara.py`, fixtures) no incluye `verify3.py`. Así, la verificación de S2-04 fallaría en la revisión de la implementación. Ruta: `S2-interfaz.md` §0.5/§7. **Corrección:** asignar el cambio al agente Python de S2 (redirigir a `datos/tmp/` o a `oraculo/verify3.json` sin alterarlo) o dejarlo como solicitud explícita con su dueño.

### Bajos
- **B1 · S2-09 (sigue abierto el B5 de la r1).** `TestProyeccionDefecto` no afirma «0 eventos de la semilla entre los días 1 y 9», que el Verifica de S2-09 pide. **Corrección:** agregar esa afirmación en §7.
- **B2 (siguen abiertos B1 y B4 de la r1).** `plan_contra_real` y `promedios {month, basis}` no tienen fixture. No se dice qué mes usa `s2-mes-30-dias` ni qué error da un evento fuera de las ventanas declaradas.
- **B3 · §2.7.** No hay regla para una partida `fortnightly` cuyo `start_date` no cae en día de pago: ¿se valida o se mueve al siguiente día de pago? Tampoco se dice qué pasa si llegan `today` y `now` a la vez (§2.1).
- **B4 · §3.11.** Por la fórmula, un gasto `variable` registrado con fecha posterior a hoy no aparece en ningún campo del flujo. Si es intencional, conviene decirlo explícitamente.
- **B5 · FORMATO §5.** Falta precisar que el contraste con `model.json` y el conteo de `SIN-DUEÑO` van por las hojas de `esperado`, no por todas las hojas del nodo. En un cruce, `sensibilidad[k]` tiene 9 hojas numéricas y S2 emite 3. `promedios` con `basis` `accrued_groceries` solo está definido por un ejemplo.

Observación fuera de alcance, no cuenta: `S3-interfaz.md` 1.1 cita FORMATO v1.0 y dice que «el integrador une los dos subcomandos en un solo main». Con FORMATO 1.1 §4.1, S3 solo agrega `cmd/motores/s3.go` con `register`. La solicitud 1 de `correccion-contrato-S2-r1.md` ya lo cubre.

## Estado de los hallazgos de la r1
| Hallazgo r1 | Estado | Comprobación |
|---|---|---|
| FORMATO: ligar `id` ↔ hoja; sensibilidad; recorte y `.#` | Cerrado | FORMATO §2.1 (regla única, ruta = `id` + ruta relativa), §4 regla 2 (listas completas, `__sobra__`, el largo calculado cuenta); `vista_sensibilidad` con claves literales; 15/15 propio |
| Piezas compartidas (`cmd/motores`, `compara.py`) y S3 fuera de FORMATO | Cerrado | FORMATO §4.1 (main.go del agente Go de S2, `register`), §5 (compara.py del agente Python de S2), §6; S3 1.1 ya usa `fixture/1`, `s3 dump` y `compara.py s3` |
| `_datos/insumo-s3.json` sin esquema | Cerrado | §5.1 `insumo-s3/1`: llaves = rutas, 55 filas, R-05 con `decimales`, R-65/R-66 listas, R-67 85 llaves (verificado) |
| Partidas planificadas en accrued/cash y `provisional` | Cerrado | §3.2 (fuera de `periodos`), §3.12 (fórmulas, `planned_expense_cents` total, `net_cash_cents`, marca por liga independiente del interruptor) |
| Campos del flujo (`variable_spent`, acumulado, base, 31-oct) | Cerrado | §3.11 tabla con fórmulas; acumulado desde `from` en 0; `today` 2026-10-31 en `s2-mes-31-dias` (queda B4 bajo) |
| Flujo solo R-13…R-21 y contraste con el valor estable | Cerrado | 16 quincenas y contraste en `TestFlujoQuincenal`; mis cifras coinciden |
| `share_rule` null = todo al titular | Cerrado | §2.7: regla efectiva con la prioridad de RN-03; Despensa → 50/50 |
| Listas sin contenido ni orden | Cerrado | §3.2–§3.13 declaran contenido y orden (queda M2 sobre los mapas) |
| Falta la fixture de interés por monto | Cerrado | `s2-interes-por-monto` (§5, §5.2); 28/28 propio |
| Tipos con huecos (persona en ingresos y partidas, importe a crédito) | Cerrado | §2.6 `person` obligatorio, §2.7 `person`, §2.4 `amount_cents` = capital o `sum_mismatch` |
