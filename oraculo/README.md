# Oráculo · Finanzas Hogar

Verificadores independientes en Python 3.14 (solo biblioteca estándar, siempre con `python3 -I`). Nada aquí lee código ni salidas de `backend/` salvo `compara.py`, que contrasta la salida del motor Go cuando existe.

## Canon (no se modifica)
- `model.json` (sha256 `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4`) y `verify3.json` (sha256 `d86087f3797e08ceabd2c721bf423115345b318ae7fe5fa13b3400cd87b95041`).
- `model.py` (Decimal) y `verify3.py` (float): importarlos ya no escribe nada (S2-04). Ejecutados, escriben en `datos/tmp/model.json` y `datos/tmp/verify3.json` (rutas relativas a la raíz del repositorio); `model.py` además dice si su salida es igual byte a byte a `model.json`.

## Oráculo financiero S2 (`fin_*.py`)
Contrato: `docs/contratos/motores/S2-interfaz.md` 1.2 y `oraculo/fixtures/FORMATO.md` 1.2. Método propio (S2-03): un libro de costos y movimientos que se recorre día por día; dinero en centavos enteros, intermedios con `Fraction`, HALF_UP solo donde lo dice la interfaz (§1.3).

| Archivo | Contenido |
|---|---|
| `fin_comun.py` | `ErrorMotor`, HALF_UP exacto, fechas, quincenas (K8), validación de `Input` (§2) |
| `fin_calendario.py` | `calendario`, `hoy`, `limites_mes`, `quincena_de`, `limites_quincena`, `semana_de`, `limites_semana` |
| `fin_devengo.py` | `periodos`, `horizonte`, `deuda_viva`, `promedios`, `promedio_mensual`, `valor_estable`, `papel` |
| `fin_presupuesto.py` | `presupuesto_mes`, `aplicar_cambio_presupuesto`, `plan_contra_real` |
| `fin_flujo.py` | `flujo_quincenal`, ocurrencias de partidas planificadas y la liga PROVISIONAL |
| `fin_proyeccion.py` | `proyeccion` (12 meses, interruptor K21) |
| `fin_metas.py` | `metas` |
| `fin_vistas.py` | vistas con las claves literales de `model.json` (§3.14) |
| `fin_genera_fixtures.py` | escribe las 27 fixtures y `_datos/` de `fixtures/s2/` (único autor, §5.0) |
| `fin_run_fixtures.py` | corre las fixtures y escribe `salidas/s2.json` (`salida/1`) |
| `fin_verify.py` | segunda vía en float de las 159 hojas de S2 |
| `compara.py` | comparación genérica de todos los motores (FORMATO §5); `fin_compara.py` = `compara.py s2` |

## Comandos (desde la raíz del repositorio)
```
python3 -I oraculo/fin_genera_fixtures.py     # solo si cambia la interfaz; es determinista
python3 -I oraculo/fin_run_fixtures.py        # 27 fixtures, 27 coinciden, 0 diferencias
python3 -I oraculo/compara.py s2              # usa datos/tmp/motor-s2.json si existe; --sin-motor para omitirlo
python3 -I oraculo/fin_verify.py              # 159 de 159 a 2 decimales
```

## Lecturas de trabajo y supuestos del oráculo (ninguno toca cifras canon)
- Liga ocurrencia ↔ captura: opción 3A de `docs/estado/bloqueos/S2-flujo-quincenal.md` (misma quincena para semanal y quincenal; mismo mes para mensual, anual y `once`). PROVISIONAL.
- Metas: `months_left` con la opción (1)A de `S2-metas.md` (del mes de hoy al de la fecha meta, ambos incluidos); `state` se emite (`in_progress`, `achieved`, `overdue`) pero ninguna fixture lo compara porque la interfaz no fija sus valores. Resumen solo con metas activas (`S2-metas-resumen`). PROVISIONAL.
- Cadencia de prueba de `s2-proyeccion-7-12-cadencia`: opción A de `S2-proyeccion-interruptor.md`. PROVISIONAL.
- Ventanas (§2.12): un evento anterior a la primera ventana declarada es `invalid_argument`; las porciones que caen después de la última se omiten.
- Proyección: las ocurrencias planificadas se generan solo dentro de los 12 meses proyectados.
- `compara.py`: si falta `datos/tmp/motor-s2.json` y no se pasó `--sin-motor`, avisa y compara sin motor (el flujo WF-1 lo pide así); las hojas de fixtures provisionales van en su propia línea.
