# revisor-S1-r2 · Reporte final (flujo P1c, 11-oct-2026)
Revisor independiente, ronda 2 de 3 del esquema 1.3. No corrigí nada. Mi informe completo es `docs/estado/p1c/REVISION-S1-r2.md`.

## Hecho
- Releí el protocolo, la sección §K de decisiones (K36 a K44), la 00003 completa, los cambios de las reglas, del diccionario y de `gen_dict.py` contra HEAD, el criterio S1, el bloqueo `S1-parametros-avisos` y los reportes de S1 y de correccion-S1-r1.
- Reejecuté por mi cuenta las verificaciones (a) a (j). Las 10 pasan.
- Los dos hallazgos de la ronda 1 quedaron cerrados: H-01 (alta) y H-02 (media).
- Hallazgos nuevos: 1 medio (M-01) y 4 bajos (B-01 a B-04). El veredicto es **NO APROBADO**.

## Archivos
- Creados: `docs/estado/p1c/REVISION-S1-r2.md` y `docs/estado/p1c/revisor-S1-r2.md`.
- Temporales en `datos/tmp/rev2/` (los scripts v.py, b.py y h.py, las bases y los registros), ya borrados.

## Verificación
- Script `v.py` (sqlite3 de Python) sobre 00001 + 00002 + 00003:
  - `integrity_check` ok, `foreign_key_check` [], 55 tablas y 55 STRICT, 152 índices y 75 disparadores.
  - K44: 15 de 15 valores iguales a los que esperaba según §K. 1.3 tiene 20 claves (10 de la 1.2 y 10 nuevas) y ninguna inesperada.
  - Los valores extremos se aceptan (no hay mínimos en el esquema) y el tipo o el signo equivocado se rechaza (9 de 9).
  - `notification` en 1.3: `weekly_summary` y `meal_plan_missing` con `{"weekday": 7}`, `monthly_report` con `{"day_of_month": 1}`. Esto sustenta M-01.
- Script `b.py`: una base 1.2 copiada, con Up y luego Down de la 00003, deja `sqlite_master` (278 filas) y los datos (55 tablas) idénticos.
- goose v3.28.0 up / down-to 0 / up: códigos 0, 0 y 0; tablas 55, 0 y 55; versión 3.
- `git diff --quiet HEAD` de 00001 y 00002: exit 0.
- `GOOSE=… python3 -I backend/db/tools/s1_verifica.py`: exit 0 y RESUMEN TODO OK. KUESKI 389500 / 239528, 0 diferencias.
- `gen_dict.py` hacia `datos/tmp/rev2/dic.md`, comparado con `cmp` contra el diccionario: idéntico.
- Script `h.py`: todas las marcas PROVISIONAL citan solo `S1-parametros-avisos`, que sigue abierto. Las 15 claves `@param` están en la fila de su sección.
- Cifras con Decimal y HALF_UP desde `model.json`: 973.75 × 7 / 31 = 219.88 y 2517.38 × 7 / 28 = 629.35 (con el diario redondeado darían 219.87 y 629.37). Pesos: 60/61 = 0.9836 y 30/61 = 0.4918.

## Supuestos
- Que `weekday` 7 y `day_of_month` 1 son «valores numéricos operativos decididos en §F» lo saco del texto de la decisión F y de la línea 2257 del diccionario. Por eso lo califiqué de inconsistencia (media) y no de valor inventado.

## Bloqueos
- Ninguno nuevo. `docs/estado/bloqueos/S1-parametros-avisos.md` sigue abierto y es legítimo.

## Pendientes
- Para el corrector de la ronda 3, que ajuste lo siguiente:
  - M-01, en `backend/db/migrations/00003_puerta1b.sql:35-57`, RN-29, la sección K44 del diccionario (se regenera) y `docs/estado/bloqueos/S1-parametros-avisos.md`.
  - Opcionalmente, B-01 a B-04.

## Solicitudes fuera de alcance
- Ninguna.

## Propuestas para Axel
- Ninguna nueva. Para la puerta 2, conviene llevar en una sola pregunta del bloqueo dónde viven `send_time`, `advance_days`, `weekday` y `day_of_month`.

## Consumo aproximado
- 28 acciones de herramienta, unos 205 KB leídos y unos 17 KB escritos.
