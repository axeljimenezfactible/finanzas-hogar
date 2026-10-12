# Revisión del contrato S2 · ronda 3 de 3

Revisor: `revisor-contrato-S2-r3` (independiente; no escribió nada de lo revisado y no corrige) · 11-oct-2026 · Vara: `docs/criterios/S2.md` 1.2, `docs/plan/decisiones.md` §B y §K. Artefactos: `docs/contratos/motores/S2-interfaz.md` 1.2 (465 líneas, 79.8 KB), `oraculo/fixtures/FORMATO.md` 1.2, reporte `docs/estado/wf1/contrato-S2.md` y corrección `correccion-contrato-S2-r2.md`. Todas las cifras de abajo salen de mis scripts en `datos/tmp/rev-s2-r3/` (ya borrados), no de los reportes.

## Veredicto
**APROBADO.** 0 hallazgos altos, 0 medios, 13 bajos. Los 6 hallazgos medios de la r2 quedaron cerrados. Ninguna cifra contradice `model.json`, §B ni §K. La interfaz alcanza para que los agentes Go y Python implementen sin leer `model.py`. Los bajos son mejoras: el integrador decide si se aplican antes de implementar.

## Verificaciones propias
- **(a) Hojas de `model.json`** (`python3 -I hojas.py`): parametros 9, periodos 122, despensas 144, credito 143, valor_estable 12, sensibilidad 45, papel 12, disposicion_por_despensa 68, horizonte 13, interes_por_monto 28. Total **596**: **S2 159** y S3 437. Expandí las llaves de la tabla §4 con mi propio analizador: 20 renglones → 159 rutas, 159 distintas, 0 duplicadas, 0 con `n` mal, 0 que no sean hoja, **0 hojas de S2 sin dueño** y 0 fuera de las ramas de S2. Duplicadas con S3: la tabla §12 de `S3-interfaz.md` solo usa ramas de S3 (credito 11 renglones, despensas 16, disposicion 2, interes_por_monto 3, parametros 5, sensibilidad 5) → **0**.
- **Receta de la semilla (§5.3), implementación propia** (`recalc.py`, `Fraction`, mes calendario): las validaciones de §2.4 y §2.5 pasan. Con §1.3, §3.2, §3.3, §3.4, §3.6 y §3.7, las vistas dan **159 de 159** hojas S2 iguales a `model.json` (0 diferencias).
- **K19, 7 cierres** (deuda = diferencia + saldo): 2026-10-31 292125 = 97375 + 194750; 11-30 97375 = 97375 + 0; 12-31 377606 = 125869 + 251737; 2027-01-31 125870 = 125870 + 0; 02-28 377606 = 125869 + 251737; 03-31 125870 = 125870 + 0; 04-10 629344 = 125870 + 503474. Los 7 coinciden con S2-07. La forma sin saldo vale solo en P2, P4 y P6.
- **K20/K39, promedios** (diario y semanal con un solo redondeo): oct 3141/21988; nov 6492/45442; feb 8991/62935; devengado de la despensa de oct 6282/43976. Los semanales con el diario ya redondeado dan 21987, 45444, 62937 y 43974, y se rechazan. Mensuales: 218396 (R-56) y 239375 (R-57).
- **Cruces:** `vista_sensibilidad` da 5 de 5 filas iguales (198753, 230220, 261687, 293154, 324621, con `axel_mensual` y `karen_mensual`), y 149972 + `axel_total` = `total` en las 5. Las series de `_datos/interes-por-monto.json` según §5.2 cumplen Σ montos = total y Σ capital = capital, con intereses ≥ 0 en los 28 pagos. La de 3,094.50 reproduce 77363 ×3 + 77361, el `capital_repuesto` de la despensa 1.
- **Otros ejemplos de la interfaz:** presupuesto 176751/200000 → `used_pct` 88, `warning` true y `limit` false. Metas: 314672/629344 → 500 décimas. `s2-mes-31-dias`: 97375 + 50000 = 147375. Flujo de 16 quincenas: una por quincena de 2026-10-q2 a 2027-06-q1; contra 88376 da +1 en 9 quincenas y −1 en 3; ciclos 353502 = 2 × 176751 y 149972 = 2 × 74986.
- **S2-04:** el sha256 de `model.json` sigue en `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` antes y después; el de `verify3.json` es `d86087f3…5041`. Corrí una copia de `verify3.py` con la salida redirigida a `datos/tmp`: `cmp` sale IGUALES. `grep -n reporte-compra-latex oraculo/*.py` sigue dando 1 línea (`verify3.py:89`). Es lo esperado: el contrato asigna la corrección al agente Python de S2 en la implementación.
- **(b)** En §9 están los 16 criterios. Cada sección citada existe (§0.3–§0.5 son incisos de §0) y cada criterio se cita en el cuerpo. Están los 5 bloqueos `S2-*` citados en §8.
- **(e)** §5 trae las 27 fixtures con su columna «Fuente» y su criterio; §5.0 nombra al autor de las 27 y de los 3 archivos `_datos/`.

## Criterios
| Criterio | Estado | Evidencia propia |
|---|---|---|
| S2-01 (alta) | Cubierto | §4, §5.1, §7: `TestCanon` con 85 subpruebas (30 de S2 = R-37…R-60 y R-68…R-73; 55 de S3 = R-01…R-36, R-61…R-67 y R-74…R-85), leídas de `model.json` al correr; la receta reproduce las 159 |
| S2-02 (alta) | Cubierto | Conteo 596/159/437; tabla §4 sin huecos ni duplicadas; 0 cruces con S3 §12 |
| S2-03 (alta) | Cubierto | Argumentos por función (§3.1–§3.14), regla de mapas (§1.1, FORMATO §4 regla 5), autoría (§5.0), semilla completa (§5.3, 159/159 propias), cruces 15 + 28 = 43 verificados |
| S2-04 (alta) | Cubierto en el contrato | §0.5 y §7 asignan y verifican `verify3.py:89`; copia redirigida igual byte a byte; el sha no cambia (bajo B10) |
| S2-05 (alta) | Cubierto | §1.3 (K14), §2.5 (verifica sin corregir), §3.2 (invariantes); validaciones de la semilla en verde |
| S2-06 (alta) | Cubierto | §1.2 (K8), ejemplos de §3.1, §5.2 (31-oct, 147375); bajo B2 (mes de 30 días) |
| S2-07 (alta) | Cubierto | §3.4; 7 de 7 cierres reproducidos |
| S2-08 (alta) | Cubierto | §1.3, §3.5; diarios, semanales y mensuales reproducidos |
| S2-09 (alta) | Cubierto | §3.12, §2.4 (`projected`), §2.10 PROVISIONAL, interruptor apagado por omisión (K21, RN-28); bajo B1 |
| S2-10 (alta) | Cubierto | §3.11; 16 quincenas y contraste reproducidos |
| S2-11 (alta) | Cubierto | §3.8 con fórmula por campo, banderas en enteros y umbrales de `alerts.*` (K44); §3.9 (K2, `onward` por omisión); bajos B6, B7 y B13 |
| S2-12 (media) | Cubierto | §3.13: fórmulas firmes, `month_summary`, metas sin monto o sin fecha; lo demás PROVISIONAL (`S2-metas`, `S2-metas-resumen`); bajo B13 |
| S2-13 (alta) | Cubierto | §0.3, §1.1, `grep` de §7 |
| S2-14 (media) | Cubierto | §1.2, §3.1 (05:59:59Z → 2026-10-30; 06:00:00Z → 2026-10-31), zona por base del sistema |
| S2-15 (alta) | Cubierto | §7 (`vet`, `gofmt`, cobertura ≥ 85 %, K25) |
| S2-16 (media) | Cubierto | §8 con 5 bloqueos S2 que existen y la nota de `S1-parametros-avisos` |

§B y §K: no encontré contradicciones. Revisé K2 (alcance con `onward` por omisión), K8 (límites de quincena y semana), K14, K19, K20/K39, K21 (interruptor apagado y cadencia PROVISIONAL) y K44 (umbrales 80/100 de `app_setting`), además de §B: ingresos solo del titular, quincena 15 y 30 o último día, devengo bimestral mitad por mes y reparto por defecto de RN-03.

## Hallazgos de la ronda anterior
| # | Hallazgo r2 | Estado | Evidencia |
|---|---|---|---|
| M1 | Claves de `argumentos` y función Python de `plan_contra_real` | **Cerrado** | Regla general en §3 y línea «Argumentos» en §3.1–§3.14: `months`/`windows`, `start`/`end`/`months`, `cut`, `total_cents`/`month_count`, `month`, `row_id`/`month`/`amount_cents`/`scope` (renglones de `entradas.datos.budgets`), `from`/`count`/`category` y `from`. Agrega `fin_presupuesto.py: plan_contra_real` |
| M2 | Claves de los mapas `{code: cents}` | **Cerrado** | §1.1: titular y externas siempre, con 0, nunca la co-titular; categorías solo con importe ≠ 0. Ejemplo `papel_contado` → `{axel: 0, karen: 9950}`, que mi cálculo confirma. FORMATO §4, regla 5 |
| M3 | Campos sin fórmula en presupuesto y metas | **Cerrado** | §3.8: categoría con y sin renglón, `total.*`, `accrued_household_cents` y `cash_household_cents` (119764 y 59882, que reproduje). §3.13: `month_summary` y metas sin monto o sin fecha; lo no decidido va PROVISIONAL con `S2-metas-resumen` |
| M4 | `projected` en `Expense` | **Cerrado** | §2.4: bool, `false` por omisión |
| M5 | Autoría de las fixtures y semilla incompleta | **Cerrado** | Encabezado y §5.0 (agente Python, `fin_genera_fixtures.py`, qué hace Go mientras falten, sin `t.Skip`); §5.3 con las 12 claves, `stable_items`, códigos, `kind` e ids; 159/159 propias |
| M6 | `verify3.py:89` sin dueño | **Cerrado** | Encabezado, §0.5 y §7: lo corrige el agente Python de S2 y la verificación queda escrita; la copia redirigida sale igual |

## Hallazgos de esta ronda (todos bajos)
- **B1 · S2-09 · §7** (sigue abierto el B1 de la r2). `TestProyeccionDefecto` no afirma «0 eventos de la semilla entre los días 1 y 9», que pide el Verifica de S2-09. **Corrección:** agregar esa afirmación a la fila de §7.
- **B2 · S2-03 y S2-06 · §5 y §5.2** (sigue abierto en parte el B2 de la r2). Faltan: fixture de `plan_contra_real` y de `promedios {month, basis}`; el mes de `s2-mes-30-dias`; y qué pasa en modo `windows` con un evento antes de la primera ventana o entre ventanas (¿error?, ¿se ignora?). **Corrección:** una consulta de cada función en fixtures que ya existen, el mes (p. ej. 2026-11) y la regla de §3.2 para eventos fuera de las ventanas.
- **B3 · §2.1 y §2.7** (sigue abierto el B3 de la r2). No se dice qué pasa con una partida `fortnightly` cuyo `start_date` no cae en día de pago, ni si llegan `today` y `now` a la vez. **Corrección:** `invalid_argument` en los dos casos, o la regla de ajuste.
- **B4 · §3.11** (sigue abierto el B4 de la r2). Un gasto `variable` registrado con fecha posterior a hoy no aparece en ningún campo del flujo. **Corrección:** decir que es intencional o mandarlo a `committed_cents`.
- **B5 · FORMATO §5** (sigue abierto en parte el B5 de la r2; la parte de `accrued_groceries` quedó cerrada en §3.5). Falta decir que el contraste con `model.json` y el `SIN-DUEÑO` van por las hojas de `esperado`; hoy solo se deduce de «43 de 43». **Corrección:** una frase en §5.
- **B6 · S2-11 · §3.9.** Con `this_month` en un renglón que empezó antes del mes salen dos renglones nuevos (el del mes y el de la continuación), pero no se dice cuál toma `máximo + 1`. Es justo el ejemplo de RN-22 que usará `s2-presupuesto-scope`. **Corrección:** «los ids nuevos se asignan en orden de `valid_from`».
- **B7 · S2-11 · §3.8.** La regla del plan 0 (`used_pct` null, banderas en `false`, `no_plan`) está escrita en el párrafo del total. Para una categoría con renglón de 0, la fórmula `actual·100 ≥ pct·plan` daría `true`. **Corrección:** decir que la regla PROVISIONAL del plan 0 aplica también a las categorías, como dice el bloqueo `S2-presupuesto-total`.
- **B8 · §5.0 y FORMATO §2.1.** «El valor convertido de `model.json`» no excluye los textos (`periodo`, `etiqueta`, `inicio`, `fin`, `cubre`, `nota`, `formula`). `vista_periodo` emite la etiqueta del mes calendario, distinta de la de `model.json`. **Corrección:** el `esperado` de una consulta con `id` ruta lleva solo las hojas numéricas convertidas (más los textos que la vista declara con su valor propio).
- **B9 · S2-03 · §6.** El Verifica de S2-03 corre `python3 -I oraculo/fin_run_fixtures.py` sin argumentos, pero §6 lo llama con `--fixtures` y `--out` y no declara valores por omisión. **Corrección:** por omisión `oraculo/fixtures/s2` y `oraculo/salidas/s2.json`.
- **B10 · S2-04 · §0.5.** `oraculo/model.py` (líneas 208–209) reescribe `model.json` al importarse, no solo al ejecutarse. Hoy deja los mismos bytes, pero un `fin_*.py` que lo importe escribiría en `model.json` en cada corrida. **Corrección:** prohibir importar `model.py` desde `fin_*.py`.
- **B11 · §3.6 y §3.12.** Falta el signo de `cent_diff_cents` (¿total − suma o suma − total?; con la semilla es 0) y la base de `index` en la proyección (¿1–12?). **Corrección:** fijar los dos.
- **B12 · §2.** No se dice si una clave de `Input` ausente equivale a una lista vacía. `_datos/interes-por-monto.json` se describe solo con «personas, `params` y 7 compras». **Corrección:** «una clave ausente = lista vacía» o exigir las 12.
- **B13 · S2-11 y S2-12 · §5.** Los criterios hablan de «fixtures» en 79/80/99/100/101 % y de «3 fixtures» de metas; la interfaz pone cada grupo en un solo archivo con varios casos. **Corrección:** declarar la equivalencia (un caso = una meta o una categoría) o separar los archivos.
