# Reporte · revisor-S1-r3 (flujo P1c, ronda 3 de 3, 11-oct-2026)

## Hecho
- Revisé el esquema 1.3 (`backend/db/migrations/00003_puerta1b.sql` entera; las diferencias contra HEAD de las reglas, del diccionario y de las herramientas) contra `docs/criterios/S1.md` 1.1 y §K (K36, K37, K39, K40, K43 y K44).
- Reejecuté las verificaciones (a)–(j) con mis propios scripts en `datos/tmp/rev3/` y probé los disparadores y el CHECK de `health_profile` con datos reales.
- Comprobé que M-01 de la ronda 2 quedó cerrado.
- Veredicto: **APROBADO**. Quedan 0 hallazgos de severidad alta o media y 6 de baja: B-01 a B-04 siguen de la ronda 2, y B-05 y B-06 son nuevos.
- Informe completo: `docs/estado/p1c/REVISION-S1-r3.md`.

## Archivos
- `docs/estado/p1c/REVISION-S1-r3.md` (nuevo)
- `docs/estado/p1c/revisor-S1-r3.md` (nuevo, este reporte)
- Temporales: `datos/tmp/rev3/` (rev.py, rev2.py, rev3.py, bases .db, dic.md y s1v.txt), borrados al terminar.

## Verificación
- `goose up` en `datos/tmp/rev3/rev.db`: rc 0, versión 3. `integrity_check` = ok, `foreign_key_check` = vacío, 55 de 55 tablas STRICT (más `goose_db_version`, que es de goose).
- Down de la 00003 sobre una copia de la 1.2: `sqlite_master` idéntico (279 objetos) y datos de las 55 tablas iguales. Solo difiere el `sqlite_sequence` de `goose_db_version`, que lleva goose. El `updated_at` de `notification` se conserva.
- `goose up` / `down-to 0` / `up`: códigos [0, 0, 0], tablas 55 → 0 → 55, versión 3, mismo esquema.
- `git show HEAD:` de 00001 y 00002: idénticos byte por byte y también sin comentarios.
- `GOOSE=… python3 -I backend/db/tools/s1_verifica.py`: exit 0, 29 líneas, RESUMEN TODO OK, 389500 y 239528.
- `python3 -I backend/db/tools/gen_dict.py datos/tmp/rev3/dic.md` y `cmp`: idéntico.
- Valores en la base: 80, 100, 30, 90, 90, 5, 15, 10, 5, 0, 0 y 24, más 12, 1 y [1, 0], todos de §A–§K. Ninguno inventado.
- PROVISIONAL: 19 líneas. Solo citan `S1-parametros-avisos`, que está abierto. Ninguna cita un bloqueo resuelto.
- Decimal HALF_UP desde model.json: 219.88 (219.87 con el diario redondeado), 629.35 (629.37), 0.9836, 0.4918, 439.76, 0.6593, 2,556 y 2,056.
- Mínimos de seguridad: 0 apariciones en la 00003, las reglas y el diccionario. El Up solo compara `< 0` y `>= 0`.

## Supuestos
- El `sqlite_sequence` de `goose_db_version` (3 contra 4) es contabilidad de goose y no cuenta como diferencia de datos de la 1.2.
- En (c) excluí `created_at` y `updated_at` al comparar datos de dos bases creadas en momentos distintos (dependen de la hora en que se corre). El esquema lo comparé completo.

## Bloqueos
- Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md`, que no es mío.

## Pendientes
- Ninguno de esta revisión. Los hallazgos bajos B-01 a B-06 quedan a criterio del integrador.

## Solicitudes fuera de alcance
- Ninguna. No corregí nada.

## Propuestas para Axel
- Ninguna propia.

## Consumo aproximado
- Acciones de herramienta: 25.
- Leído: unos 140 KB:
  - protocolo 5.6; §K 14.7; 00003 20; bloqueos ~9; S1.md ~22; corrección r2 8.5;
  - diferencias de las reglas ~14; secciones del diccionario ~12; reporte S1 ~5;
  - salidas de verificación ~25.
- Escrito: unos 27 KB (scripts temporales ~10, REVISION-S1-r3.md 13.5 y este reporte ~3.5).
- cerrado_por_tope: false.
