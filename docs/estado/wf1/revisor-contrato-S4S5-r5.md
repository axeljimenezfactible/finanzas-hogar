# revisor-contrato-S4S5-r5 · reporte de agente (WF-1, 11-oct-2026)

## Hecho
- Revisé, de forma independiente, `S4-interfaz.md` 1.4, `S5-interfaz.md` 1.4 y `productos-formato.md` 1.1 contra `docs/criterios/S4.md` 1.2, `S5.md` 1.2 y `decisiones.md` §D, §E y §K. Leí todo por partes con `sed -n` y `grep`.
- Volví a calcular todo con mi propio script: el canon S5-01 desde los renglones de `data.json`, S5-02, K40, S4 §8.3, RN-15, S4-07 e IMC, el ejemplo de S4 §6.6, `QuoteFor` y los dos casos de `compare_basket` de S5 §7.3, y la cobertura del catálogo (39 renglones, 61 piezas, 3,094.50).
- Comprobé los dos medios de la ronda 4: los dos están **cerrados**.
- Encontré 3 medios nuevos (M-r5-1 a M-r5-3) y 1 bajo nuevo (B-r5-1). Los 3 bajos de la r4 siguen abiertos.
- Veredicto: **NO APROBADO** (0 altos, 3 medios).

## Archivos
- `docs/estado/wf1/REVISION-contrato-S4S5-r5.md` (informe: veredicto, tabla de criterios, hallazgos, estado de los anteriores).
- `docs/estado/wf1/revisor-contrato-S4S5-r5.md` (este reporte).
- No modifiqué ningún artefacto revisado.

## Verificación
- `python3 -I datos/tmp/rev-s4s5-r5/verif.py /home/deepsec/finanzas-hogar` → **46 OK, 0 FAIL**:
  - **Canon:** idénticos 277350 / 287400 (+10050 = **+$100.50**); totales 309450 / 323300 (+13850 = **+$138.50** = 10050 + 5900 − 2100), iguales a `totales`; piezas 61 / 60 / 55; renglones 33 / 3; ganadores 22 / 8 / 3.
  - **S5-02:** 2933 / 2800, 12692 / 16406, 3375 / 2833.
  - **K40:** 1, 60/61 (0.9836), 30/61 (0.4918), 0, y «no compite» desde el día 91; la recta baja 1/61 por día.
  - **S4:** §8.3 4425 / 2213 / 8850; RN-15 1486571; S4-07 2099 / 1596 / 1799 / 1846 (exactos 2098.50 y 1596.30); RN-17 2556 / 2056; IMC 185 / 250 / 300 / 237.
  - **S4 §6.6:** 8400 + 3000 + 7980 = 19380, con servilletas `null`.
  - **S5 §7.3:** 2200 / 1500 → Bodega; 9430 / 9460 → Walmart; con la regla 1.3, 10780.
  - **Catálogo:** 36 productos y 39 presentaciones; §5.2 = `cuentas` renglón por renglón.
- `ls docs/estado/bloqueos/`: existen los 7 bloqueos citados.
- `grep` en `reglas-negocio.md`: no tiene regla de desempate. Ya está pedida en la solicitud 6 de `contrato-S4S5.md`, y el contrato la usa como PROVISIONAL citando `S5-reglas-comparativa`.
- `grep` en `FORMATO.md` §4 y §5: las listas se comparan por índice (base de M-r5-3).
- Extracto C1 `compras-tiendas-a.yaml`: `reasons` admite `manual` y `needed_milli` es opcional.
- **Criterios cubiertos:** S4-01 a S4-15 y S5-01 a S5-13 tienen sección (tabla en el informe). M-r5-1 toca S4-03 y S4-06; M-r5-2 toca S4-03; M-r5-3 toca S4-03, S5-03 y los oráculos.

## Supuestos
Ninguno que cambie cifras. Los ejemplos de M-r5-1 y M-r5-3 son ILUSTRATIVOS, solo para mostrar el defecto.

## Bloqueos
Ninguno nuevo.

## Pendientes
Para el corrector de la ronda siguiente (todo en `docs/contratos/motores/`):
- **M-r5-1:** `S4-interfaz.md` §8.1 paso 1 y §12.1 `quotes`.
- **M-r5-2:** `S4-interfaz.md` §6.2 paso 5 y §12.1 `catalog`.
- **M-r5-3:** `S4-interfaz.md` §12.2; `S5-interfaz.md` §7.2, §7.4 y §15.2.
- Bajos opcionales: B-r5-1, B-r4-1, B-r4-2 y B-r4-3.

## Solicitudes fuera de alcance
Ninguna.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna (no hay valores operativos de K31 en esta revisión).

## Consumo aproximado
~32 acciones de herramienta; ~215 KB leídos; ~24 KB escritos (informe, este reporte y el script temporal, ya borrado).
