# Reporte · revisor-U1-r2 (flujo P1c, 11-oct-2026)

## Hecho
- Revisión independiente de U1 1.2, ronda 2, solo contra `docs/criterios/U1.md` 1.2 y K35, K36, K38, K43 y K44. Veredicto: **NO APROBADO** (0 altos, 1 medio, 7 bajos). El informe está en `docs/estado/p1c/REVISION-U1-r2.md`.
- Comprobé el hallazgo H-01 de r1:
  - **Cerrado** en `U1-arquitectura.md` (N-298, N-300 y filas 318-319 de §2.3) y en `verificar_u1.py`.
  - **Abierto** en «Para C1» de `docs/estado/p1c/U1-1.2.md` (línea 60): la solicitud 1 de `correccion-U1-r1.md` no se aplicó. Queda como H-r2-01 (media).
- H-02 a H-07 (bajas) siguen igual, con las rutas actualizadas. Encontré un bajo nuevo, H-r2-02: el alcance de `send_time` en los reportes no coincide entre U1 y 00003.

## Archivos
- `docs/estado/p1c/REVISION-U1-r2.md` (nuevo).
- `docs/estado/p1c/revisor-U1-r2.md` (este reporte).
- Temporales en `datos/tmp/revisor-u1-r2/` (`chk.py`, `bl.py`, `num.py`, `oa.py` y dos extractos de «Para C1»). Los borré al terminar.

## Verificación
- `python3 -I docs/design/tools/verificar_u1.py`: sale con 0, 71 líneas, «RESULTADO: TODO OK». La única PENDIENTE es N-320 a N-327.
- `chk.py` (mío):
  - 211 N (138 L + 73 E), sin duplicados, la última es N-327; no hay referencias a N inexistentes.
  - La matriz §3.29 tiene 42 filas, todas con las mismas columnas, y no le falta ninguna N.
  - 27 P-xx y 14 P-26.xx.
  - F-1 acumula 1, 2 y 3 toques.
  - PROVISIONAL solo aparece en las listas de cambios (líneas 22, 37, 38 y 51).
  - `advance_days` en las líneas 23, 29, 303, 319, 790, 792 y 1207; `days_before` en las líneas 23, 29, 318, 790 y 792.
- `oa.py` (mío):
  - `openapi.yaml` 1.1.0 tiene 200 N con operación.
  - De U1 quedan sin operación ni Local exactamente N-320 a N-327, que son las mismas 8 «Nuevas» de «Para C1».
  - No hay N de openapi ausentes en U1 ni N Local con operación.
- `bl.py` (mío): bloqueos abiertos: `S1-parametros-avisos` y `criterios-nxx`. U1 describe bien el estado de los 13 bloqueos que cita.
- `num.py` (mío):
  - Las cifras nuevas del diff son 0 kcal, [1, 0], 1 GB, 10 MB, 12 caracteres, 15 min, 24 h, 5 s, 30 y 90 días. Todas salen de K6, K11, K22, K36, K40, K43, K44 o F.
  - No hay topes: «≥ 1», «365» y «no menos» dan 0 coincidencias.
- `diff` de «Para C1» entre `U1-1.2.md` y `correccion-U1-r1.md`: solo difiere la viñeta de N-297, N-298 y N-300 (línea 19 del extracto). `U1-1.2.md` es de las 15:06 y la corrección de las 15:18.
- `grep` en `00001_esquema.sql`, líneas 1240-1247: `days_before` [1, 0], `thresholds_pct` [80, 100], `weekday` 7 y `day_of_month` 1.
- `grep` en el encabezado de `00003_puerta1b.sql`: 12 claves `@param` y la definición de `advance_days` en la línea 46.

## Supuestos
- K43 dice que «los tres valores numéricos (10 MB, 5 s, `advance_days`) nacen con esos valores». Lo leo como que aprueba la clave `advance_days` sin valor inicial. Me baso en la propuesta aprobada (`C1-valores-k31` línea 6: «los valores por defecto los propone S7») y en 00003 línea 46. Así lo tratan también U1 y S1. No cambia cifras ni esquema.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Para cerrar H-r2-01: el integrador cambia `docs/estado/p1c/U1-1.2.md` línea 60 por la viñeta de «Para C1 (corregida)» de `docs/estado/p1c/correccion-U1-r1.md`.

## Solicitudes fuera de alcance
- `docs/estado/p1c/U1-1.2.md` línea 60: aplicar la solicitud 1 de `correccion-U1-r1.md` (H-r2-01).
- `backend/db/migrations/00003_puerta1b.sql` línea 45 y RN-29 (S1): alinear con U1 el alcance de `send_time`, es decir, qué reportes lo llevan (H-r2-02).
- S1: confirmar si `weekday` y `day_of_month` de `notification.params_json` son parámetros de K44, o preguntarlo en un bloqueo (H-07).

## Propuestas para Axel
- Ninguna.

## Consumo aproximado
- Acciones de herramienta: 33.
- Leído: ~175 KB (protocolo 6, §K 15, reportes de r1 y corrección ~30, diff de U1 27, criterios ~8, salida del verificador ~10, código del verificador ~18, migraciones y bloqueos por grep ~15, salidas de mis scripts ~20, «Para C1» y líneas sueltas ~26).
- Escrito: ~22 KB (scripts ~6, REVISION 12.5, este reporte ~4).
