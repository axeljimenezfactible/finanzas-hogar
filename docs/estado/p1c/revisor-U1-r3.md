# revisor-U1-r3 · Reporte del agente (flujo P1c, 11-oct-2026)

## Hecho
- Revisé U1 1.2 de forma independiente (ronda 3 de 3), solo contra `docs/criterios/U1.md` 1.2 y K35, K36, K38, K43 y K44. Informe: `docs/estado/p1c/REVISION-U1-r3.md`.
- **Veredicto: APROBADO.** Hay 0 hallazgos altos, 0 medios y 10 bajos (3 nuevos y 7 que siguen de la ronda 2).
- H-r2-01 (media) queda **cerrado** por la «otra salida» que admitió la ronda 2. Lo que queda es bajo (H-r3-01): `U1-1.2.md` conserva la «Para C1» vieja.

## Archivos
- `docs/estado/p1c/REVISION-U1-r3.md` (nuevo)
- `docs/estado/p1c/revisor-U1-r3.md` (este reporte)
- Temporales en `datos/tmp/rev-u1-r3/` (script de recuento, script de pruebas negativas y 6 raíces de prueba): borrados al terminar.

## Verificación
- `python3 -I docs/design/tools/verificar_u1.py` sale con 0, da 63 OK, 0 FALLA y «TODO OK». Hay una sola línea PENDIENTE, `C1 v1.2`, con N-320 a N-327, y apunta a «Para C1 (corregida r2)».
- Recuento propio con `check.py`:
  - 211 N-xxx (138 L y 73 E), sin duplicados;
  - la matriz §3.29 coincide en los dos sentidos y la columna L/E es coherente;
  - 27 pantallas y 14 secciones, sin P-26.15;
  - PROVISIONAL solo en el historial (líneas 22, 38, 39 y 52);
  - F-1 = 3 toques, sin cambios.
- Cruce propio con `grep` de las `x-necesidades` de openapi.yaml: las N-xxx de U1 sin operación son las 3 Local más N-320 a N-327. No hay ninguna N-xxx citada por una operación que falte en U1.
- Cruce con 00003: las 15 líneas `@param` (35 a 49) coinciden con la tabla de §2.3 y con N-320 a N-327 en clave, valor y sección. RN-29 (líneas 291 a 305) también coincide.
- Bloqueos: el único abierto es `S1-parametros-avisos`, y U1 lo cita como abierto.
- Pruebas negativas sobre copias:
  - V0, copia sin cambios: sale con 0.
  - V1 a V5 salen con 1: «Para C1» con la viñeta vieja, PROVISIONAL con un bloqueo resuelto, ajuste de 100 kcal, sección P-26.15 y umbral del 85 %.

## Supuestos
- Doy por cerrado H-r2-01 porque la corrección r2 aplicó la «otra salida» que la propia ronda 2 escribió: el verificador apunta a la sección corregida. Además, la fila 9 de U1 nombra esa sección y la prueba V1 demuestra que la viñeta vieja da FALLA. La copia vieja que sigue en `U1-1.2.md` queda como hallazgo bajo, no medio, porque tanto el artefacto como su verificador remiten a la sección correcta.
- Las 8 N-xxx sin operación son PENDIENTE esperado de U1-03 hasta C1 v1.2, como pide el encargo (a).

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Ninguno del revisor.

## Solicitudes fuera de alcance
1. `docs/estado/p1c/U1-1.2.md`, sección «Para C1»: aplicar la solicitud 1 de `correccion-U1-r2.md` o poner una nota de que esa sección la sustituye (H-r3-01). Motivo: que no queden dos «Para C1» distintas.
2. `docs/design/U1-arquitectura.md`: los bajos H-r3-02, H-r3-03, H-r2-02, H-02, H-03, H-04 y H-05, para una versión siguiente o para WF-1. Motivo: redacción y trazabilidad; no cambian cifras.
3. S1 (`00003_puerta1b.sql` / RN-29): H-07 (`weekday` y `day_of_month` como parámetros de K44) y H-r2-02 (alcance de `send_time`). Motivo: coherencia con K44 y con U1.

## Propuestas para Axel
- Ninguna.

## Consumo aproximado
- Acciones de herramienta: ~36 de 60.
- Leído: ~190 KB. Desglose:
  - protocolo 6 y decisiones §K 15;
  - reporte U1-1.2 15, correccion-U1-r2 15 y REVISION-U1-r2 por partes ~8;
  - diff de U1 37;
  - criterios U1 ~12;
  - verificador por partes ~6 y su salida ~5;
  - 00003 por `grep` ~8, RN-29 ~2 y bloqueo S1 ~3;
  - INTEGRACION §7 ~3;
  - salidas de mis scripts ~15.
- Escrito: ~25 KB. Desglose: informe ~12, este reporte ~4, scripts temporales ~8.
