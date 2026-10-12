# Reporte · oraculo-S2 (WF-1, 11-oct-2026)

## Hecho
- Oráculo financiero S2 en Python 3.14 con biblioteca estándar (`Fraction`, enteros de centavos, HALF_UP solo donde lo dice S2-interfaz §1.3), con todas las funciones de la interfaz 1.2 (§3.1–§3.14) y sus errores (§3.15). Método propio (§0.4): libro de costos y movimientos recorrido día por día; no lee nada de `backend/`.
- `fin_genera_fixtures.py` (único autor, §5.0) escribe las 27 fixtures de §5 y los 3 archivos de `_datos/` (semilla §5.3, insumo-s3 §5.1 con 55 filas, interes-por-monto §5.2). Ids ruta: esperado convertido de `model.json`. Ids libres: cifras fijadas por la interfaz o los criterios, citadas en `fuente`, y el resto calculado por el oráculo. **0 hallazgos**: ninguna cifra fijada difiere del oráculo. El generador es determinista (dos corridas, mismo sha256).
- `fin_run_fixtures.py` → `oraculo/salidas/s2.json` (`salida/1`, determinista; rechaza `--out` fuera del repositorio).
- `compara.py` genérico (FORMATO §5: ramas s2/s3/s4/s5, cruces, data.json, razones con `decimales`, provisionales en su propia línea, DIF, SIN-DUEÑO, códigos 0/1/2) y `fin_compara.py` (= `compara.py s2`).
- `fin_verify.py`: segunda vía en float (agrupa por mes, pesos, redondeo a 2 decimales).
- Higiene S2-04: `model.py` y `verify3.py` ya no escriben al importarse; ejecutados, escriben en `datos/tmp/` (rutas relativas a la raíz). `model.json` y `verify3.json` sin cambio.
- `oraculo/README.md`.

## Archivos
Nuevos: `oraculo/fin_comun.py`, `fin_calendario.py`, `fin_devengo.py`, `fin_vistas.py`, `fin_presupuesto.py`, `fin_flujo.py`, `fin_proyeccion.py`, `fin_metas.py`, `fin_run_fixtures.py`, `fin_genera_fixtures.py`, `fin_verify.py`, `fin_compara.py`, `compara.py`, `README.md`; `oraculo/fixtures/s2/` (27 casos + `_datos/semilla-oraculo.json`, `_datos/insumo-s3.json`, `_datos/interes-por-monto.json`); `oraculo/salidas/s2.json`; este reporte. Modificados: `oraculo/model.py` (escritura bajo `__main__`, destino `datos/tmp/model.json`), `oraculo/verify3.py` (línea 89 → bajo `__main__`, destino `datos/tmp/verify3.json`).

## Verificación
- `sha256sum oraculo/model.json oraculo/verify3.json` antes y después: `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` y `d86087f3797e08ceabd2c721bf423115345b318ae7fe5fa13b3400cd87b95041` (iguales).
- `python3 -I oraculo/model.py` → `datos/tmp/model.json igual a oraculo/model.json: True`, `OK asserts`. `python3 -I oraculo/verify3.py; cmp datos/tmp/verify3.json oraculo/verify3.json` → iguales. Importar `model` y `verify3` no escribe nada. Temporales borrados.
- `grep -n reporte-compra-latex oraculo/*.py` → 0 líneas (**S2-04**).
- `python3 -I oraculo/fin_run_fixtures.py` (y con `--fixtures oraculo/fixtures/s2 --out oraculo/salidas/s2.json`) → `27 fixtures, 27 coinciden, 0 diferencias` (**S2-03**). Fixtures: s2-periodos-p1…p6, s2-valor-estable, s2-horizonte, s2-papel, s2-sensibilidad, s2-interes-por-monto, s2-deuda-viva-identidad, s2-fin-de-mes-feb-2027, s2-mes-30-dias, s2-mes-31-dias, s2-semana, s2-zona-medianoche, s2-promedios, s2-ventanas-dia-5, s2-proyeccion-7-12-defecto, s2-proyeccion-7-12-cadencia (PROVISIONAL), s2-proyeccion-csv, s2-liga-provisional (PROVISIONAL), s2-flujo-quincenal, s2-presupuesto-scope, s2-alertas-80-100, s2-metas (PROVISIONAL).
- `python3 -I oraculo/compara.py s2` → `oráculo contra esperado: 2318 de 2318`, `model.json (ramas de s2): 159 de 159 hojas coinciden` (**S2-02**), `cruces con otras ramas: 43 de 43`, `data.json: 0 de 0`, `provisionales: 867 de 867`, `diferencias: 0`, código 0. A las 18:54 existía `datos/tmp/motor-s2.json` del agente Go (generado con estas fixtures): `motor contra esperado: 2318 de 2318` y `oráculo contra motor: 2318 de 2318`; ese agente lo borró después, así que la última corrida dice «sin motor (no existe datos/tmp/motor-s2.json)».
- Pruebas negativas de `compara.py` en copia temporal: una hoja alterada → `DIF s2-periodos-p1 periodos[0].devengado.total esperado=194751 oraculo=194750` y `158 de 159`; un id ruta inexistente → código 2.
- `python3 -I oraculo/fin_verify.py` → `159 de 159 hojas de S2 coinciden a 2 decimales (float contra la vía exacta)` y `159 de 159 (float contra model.json)`.
- Errores de §3.15 comprobados: unknown_field, shares_mismatch, income_not_holder, budget_overlap, invalid_zone, sum_mismatch, cadence_before_month_7, invalid_argument (función desconocida), oracle_view_shape; out_of_range e invalid_argument viven en fixtures.
- S2-05 (lado oráculo): 267 invariantes (Σ por persona = total, pagos + gastos = efectivo, Σ partes = monto) en `salidas/s2.json`, 0 con diferencia. S2-06: los 4 casos (feb-2027, 30 días, 31 días con el gasto del 31-oct en q2: 50000 en `variable_spent_cents`, `household_pending_cents` 59882; semana 2026-W42 12–18 oct). S2-07: 7 cierres con la identidad, sin saldo solo en P2, P4 y P6. S2-08: 3141/21988, 6492/45442, 8991/62935, 6282/43976, 218396 y 239375. S2-11: banderas en 80 y 100 (79 999 de 100 000 → 80 % `normal`). S2-12: Σ plan = restante (44953 ×6 + 44954). S2-14: 05:59:59Z → 30-oct, 06:00:00Z → 31-oct, 2027-07-01T05:30Z → 30-jun.

## Supuestos
Ninguno cambia cifras canon ni esquema; los que tocan salidas no canon están marcados:
- Liga ocurrencia ↔ captura con la opción 3A de `S2-flujo-quincenal` (lectura de trabajo de la interfaz); `once` usa la ventana del mes. PROVISIONAL.
- `months_left` con la opción (1)A de `S2-metas`; `state` (in_progress, achieved, overdue) no se compara en ninguna fixture porque la interfaz no fija sus valores. PROVISIONAL.
- Cadencia de prueba de `s2-proyeccion-7-12-cadencia` con la opción A de `S2-proyeccion-interruptor`. PROVISIONAL.
- Ventanas: un evento anterior a la primera ventana es `invalid_argument`; la proyección genera ocurrencias solo dentro de sus 12 meses.

## Bloqueos
Ninguno nuevo. Se usan los abiertos: `docs/estado/bloqueos/S2-flujo-quincenal.md`, `S2-metas.md`, `S2-metas-resumen.md`, `S2-proyeccion-interruptor.md`, `S2-presupuesto-total.md`.

## Pendientes
- El integrador repite §6 con el binario del motor (`datos/tmp/motor-s2.json`) al juntar y corre `python3 -I oraculo/compara.py s2` sin `--sin-motor`.

## Solicitudes fuera de alcance
1. `oraculo/fixtures/FORMATO.md` §5: decir que si falta el archivo del motor y no se pasa `--sin-motor`, `compara.py` avisa y compara sin motor (así lo pide la tarea de WF-1), en lugar del código 2.
2. `docs/contratos/motores/S2-interfaz.md` §3.13: fijar `months_left` y los valores de `state` cuando Axel decida `S2-metas`; entonces se agrega `state` a `s2-metas`.
3. S2-interfaz §2.12 y §3.12: declarar qué pasa con un evento anterior a la primera ventana y si las ocurrencias anteriores a `from` entran a `cumulative_diff_cents`. S2-interfaz §3.11: decir si `once` cuenta como «mensual o mayor» para la ventana de la liga.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: el oráculo no tiene valores operativos de K31.

## Consumo aproximado
34 acciones de herramienta; ~232 KB leídos; ~141 KB escritos con herramientas (código 134 KB + reportes). Los scripts generaron además 243 KB de fixtures y 145 KB de `salidas/s2.json`.
