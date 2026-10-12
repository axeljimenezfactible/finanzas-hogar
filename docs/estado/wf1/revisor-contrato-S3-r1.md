# revisor-contrato-S3-r1 · reporte de agente (WF-1, 11-oct-2026)

## Hecho
- Revisé `docs/contratos/motores/S3-interfaz.md` v1.0 y `docs/estado/wf1/contrato-S3.md` contra `docs/criterios/S3.md` v1.2 y decisiones §C y §K. Informe: `docs/estado/wf1/REVISION-contrato-S3-r1.md`.
- Veredicto: **NO APROBADO** (0 altos, 8 medios, 4 bajos). Las cifras y el conteo de hojas están bien; faltan definiciones para implementar sin preguntar.

## Archivos
- `docs/estado/wf1/REVISION-contrato-S3-r1.md` (nuevo)
- `docs/estado/wf1/revisor-contrato-S3-r1.md` (este reporte)
- Temporales `datos/tmp/rev-s3r1/hojas.py` y `datos/tmp/rev-s3r1/formulas.py`: creados y borrados.

## Verificación
- `python3 -I datos/tmp/rev-s3r1/hojas.py oraculo/model.json` → `total hojas numericas 596`; `S3 437 S2 159 otras {}`; `renglones 42 suma declarada 437`; `hojas S3 437 sin duenio 0 dos duenios 0 conteo distinto []`.
- `python3 -I datos/tmp/rev-s3r1/formulas.py oraculo/model.json datos/semilla/data.json` → `coinciden 339 diferencias 0`. Reproduce con las reglas de la interfaz: tasa y tasa_pct, R-05 y R-06, las 4 despensas, partes y composición, `disposicion_por_despensa`, `karen_fija`, los 17 eventos, las series A/B, los mínimos, la sensibilidad, `interes_por_monto`, los cruces de deuda viva y de efectivo P1…P6, y las cifras de las fixtures.
- `sha256sum oraculo/model.json` → `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` (intacto).
- `grep -n reporte-compra-latex oraculo/*.py` → `verify3.py:89` todavía escribe fuera del repositorio (H11).
- Cruces de lectura: C1 `creditos-b.yaml` (`missing_constraints` con 4 valores, `violated_constraint`, `score_bp` «entre los escenarios generados»), 00003 (`comparator.timeout_seconds` = 5; disparador ≥ 0), 00001 (`optimizer.weights` 1,1,1,1; `vat_on_interest_bp` NULL = la tasa ya trae IVA), RN-11 y RN-21, `S2-interfaz.md` §2.5 y §5, `oraculo/fixtures/FORMATO.md`.
- Criterios: S3-02, 03, 04, 05, 06, 07, 08, 09, 10, 15, 16, 17, 18, 19 y 20 cubiertos; S3-01, 11, 12, 13 y 14 parciales (tabla en el informe).

## Supuestos
- El revisor evalúa el contrato, no una implementación: «Cumple» quiere decir que el contrato especifica lo que pide el criterio y que las cifras que fija cuadran con `model.json`.

## Bloqueos
Ninguno nuevo.

## Pendientes
Lo corrige el autor del contrato: H1–H8 (medios) y H9–H12 (bajos) de `docs/estado/wf1/REVISION-contrato-S3-r1.md`. Después sigue la ronda 2.

## Solicitudes fuera de alcance
1. `docs/contratos/reglas-negocio.md` RN-11: el tiempo máximo ya es 5 s y configurable (K43, K44); hoy dice «lo propone S3 … (K31)». También aclarar «generados» frente a «válidos» en la normalización (H4).
2. `oraculo/fixtures/FORMATO.md` (integrador o contrato-S2): si S3 conserva `$model`, extender el formato común en vez de tener dos (H1).
3. `docs/criterios/S3.md` S3-12: si la lectura de `app_setting` y el 503 solo pueden verificarse en la API, mover esa parte del «Verifica» a WF-2 (H8).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna: no hay valores K31 que proponer en esta revisión.

## Consumo aproximado
26 acciones de herramienta; unos 160 KB leídos; unos 31 KB escritos (dos scripts temporales de unos 9 KB y dos informes de unos 22 KB).
