# revisor-U1-r1 · Reporte del agente (flujo P1c, 11-oct-2026)

## Hecho
- Revisé U1 1.2 (`docs/design/U1-arquitectura.md`, solo el diff contra HEAD y búsquedas puntuales) contra `docs/criterios/U1.md` 1.2 y contra K35, K36, K38, K43 y K44.
- Volví a correr `docs/design/tools/verificar_u1.py` y además escribí dos scripts propios: recuento de N-xxx, P-xx y P-26.xx; cobertura de la matriz §3.29; cruce con `x-necesidades` de `openapi.yaml`; marcas PROVISIONAL; integridad de tablas; cifras nuevas del diff.
- Informe con veredicto, tabla de criterios y hallazgos: `docs/estado/p1c/REVISION-U1-r1.md`. **Veredicto: NO APROBADO**, con 0 hallazgos altos, 1 medio (H-01) y 6 bajos.

## Archivos
- `docs/estado/p1c/REVISION-U1-r1.md` (nuevo)
- `docs/estado/p1c/revisor-U1-r1.md` (este reporte)
- Temporales `datos/tmp/rev-u1-r1/` (check.py, check2.py, u1.diff), ya borrados.

## Verificación
- `python3 -I docs/design/tools/verificar_u1.py`: sale con 0 y da «RESULTADO: TODO OK». Sale 1 línea PENDIENTE C1 v1.2 con N-320 a N-327.
- `check.py` (propio): 211 N-xxx definidas (L 138, E 73), sin duplicados, la última es N-327. Las 211 están en la matriz §3.29. La matriz tiene 26 P-xx más 14 P-26.xx, y P-26 queda cubierta por sus secciones. §2.1 tiene 27 pantallas y §2.3 tiene 14 secciones con nombres exactos. N-xxx sin operación ni Local: exactamente N-320 a N-327, igual a «Para C1». N-xxx en openapi que no existen en U1: ninguna. PROVISIONAL: solo en la historia de cambios (líneas 22, 36, 37 y 50), más 2 usos comunes de la palabra en las líneas 248 y 327.
- `check2.py` (propio): ninguna fila de tabla tiene columnas de más o de menos. Las cifras nuevas del diff salen todas de §A, §F o §K. «≥ 1», «365» y «no menos» dan 0 coincidencias.
- F-1 (líneas 899-922): ningún bloque del diff la toca y da 3 toques. Las claves de «Para C1» coinciden con las que siembra `00003_puerta1b.sql`.

## Supuestos
- Tomé las filas PROVISIONAL de las listas de cambios de 1.1 como historia, no como marcas vigentes. Lo dejo como hallazgo bajo H-05.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Corregir H-01 (media) y los bajos que se decidan, y luego la ronda 2.

## Solicitudes fuera de alcance
1. `docs/contratos/reglas-negocio.md` RN-29 (S1): confirmar si el día 1 del reporte mensual y el domingo del semanal y del aviso «plan sin armar» son parámetros de K44, o dejarlo en un bloqueo (H-07).
2. C1 v1.2 (`docs/contratos/`): aplicar «Para C1» de U1, incluida la corrección de H-01 (valor inicial de `days_before` [1, 0] para «Restablecer»).

## Propuestas para Axel
- Ninguna propia. Las validaciones de coherencia que propuso U1-1.2 no entran al esquema ni a U1 como valores. Están bien dejadas para la puerta 2.

## Consumo aproximado
- Acciones de herramienta: 27.
- Leído: ~155 KB (protocolo 6, §K 15, criterios U1 20, reporte U1-1.2 15, diff de U1 ~45 en partes, diff del verificador 8, salidas del verificador y de los scripts ~15, bloqueos 6, INTEGRACION §7 2, RN-29 5, búsquedas ~15).
- Escrito: ~22 KB (scripts temporales 6, informe 11, este reporte 4).
