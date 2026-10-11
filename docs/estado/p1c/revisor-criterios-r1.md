# revisor-criterios-r1 · Flujo P1c (puerta 1b) · 11-oct-2026

Revisor independiente, ronda 1 de 3, de los criterios 1.2. Informe: `docs/estado/p1c/REVISION-criterios-r1.md`. Veredicto: **NO APROBADO** (3 altos, 2 medios, 8 bajos).

## Hecho
- Leídos: protocolo completo; decisiones §K (K1–K44, con K35–K44 al detalle); INDICE §1–§7 y §10–§11; reporte `criterios-1.2.md`; diffs de C1, U4, U5, C2, S6, S3, U1 y U2; criterios S2-08, S5-04, S5-10, S4-07, U2-02, U2-08, S3-12, U1-16, S6-10, S7-04, S7-08, S7-10 y P3-01 a P3-13; `tools/conteo.py` y el diff de `verif_final.sh`; respuestas de Axel en `S1-fuentes-precio`; bloqueo abierto `S1-parametros-avisos`.
- Reejecutado con scripts propios: conteo por archivo, secciones, Verifica, Evidencia y §5; cifras de K39 y K40 con `Decimal`; los *one-liners* de S2-08 y S5-04 tal cual; búsqueda de PROVISIONAL contra los 41 bloqueos con respuesta; L OKLab de `#7AA3EE`; tamaño y número de operaciones de los extractos; `make verify-criterios`; `conteo.py`.
- Informe con veredicto, tabla (a)–(g), 13 hallazgos con ruta y corrección, y estado de los anteriores (B-01 a B-04 de r3 y pendiente 4 de INTEGRACION §7: corregidos).

## Archivos
- `docs/estado/p1c/REVISION-criterios-r1.md` (nuevo).
- `docs/estado/p1c/revisor-criterios-r1.md` (este reporte).
- Temporales `datos/tmp/revcrit-r1/` (`conteo_rev.py`, `prov.py`, `oneliners.py`, `verify.log`), borrados al terminar.

## Verificación
- `python3 -I datos/tmp/revcrit-r1/conteo_rev.py`: 18 archivos, 261 (201 / 60 / 0), 5 secciones en orden, todos con Verifica y Evidencia. Único caso «sin §5»: WF4-07, cubierto por el rango «WF4-01…WF4-12». 41 bloqueos con respuesta; sin respuesta: `S1-parametros-avisos` (abierto) y `criterios-nxx` (resuelto por su sección «Resolución»).
- `python3 -I docs/criterios/tools/conteo.py`: «INDICE contra conteo propio: 0 diferencias».
- `Decimal`: 219.88 (no 219.87), 629.35 (no 629.37), 0.9836, 0.4918, 0, y peso(29) = 1. Los *one-liners* de S2-08 y S5-04 imprimen lo esperado. 973.75, 1,947.50 y 2,517.38 en `model.json`, con sha256 igual al de INDICE §1.
- `make verify-criterios GOOSE=…/gobin/goose`: **exit 2**. Las comprobaciones de los criterios dan lo esperado; falla `matriz_c102.py` porque N-320…N-327 (U1 1.2) no están en la matriz de C1 §10 (hallazgo B-01, fuera de los criterios).
- `grep -n` de K35–K44 en los criterios: todas aplicadas en su archivo (ver la tabla (b) del informe).
- Extractos: 17, el mayor de 76,202 bytes; 201 operaciones en los extractos y 201 en `openapi.yaml`; `uploadExpenseAttachment` ya documenta 413 en C1 (C1-23 no inventa el código).
- L OKLab de `#7AA3EE` = 0.716525 (B-03).

## Supuestos
- `make verify-criterios` falla por un artefacto de otro agente (U1 1.2 contra C1), no por los criterios; lo registré como baja (B-01) y no como alta.
- S7 no estaba en la lista de K aplicadas, pero lo evalué porque la revisión es contra §K en todo el directorio, y S7-08 contradice K43 (A-03).

## Bloqueos
Ninguno nuevo.

## Pendientes
- Ronda 2: comprobar A-01, A-02, A-03, M-01 y M-02 corregidos, y volver a correr `make verify-criterios` cuando C1 v1.2 mapee N-320…N-327.

## Solicitudes fuera de alcance
- `docs/contratos/C1-datos.md` §10 (agente de C1 v1.2): mapear N-320…N-327 o declararlas Local o PROVISIONAL con bloqueo, para que `make verify-criterios` vuelva a exit 0.
- `Makefile:155-158`: sumar la cadena de `conteo.py` a las esperadas (repite la solicitud de `criterios-1.2`).
- Orquestador: llevar `S1-parametros-avisos` a la puerta 2. De esa respuesta dependen C1-22, U4-12, S7-04 y S7-08.

## Propuestas para Axel
Ninguna propia. La de mínimos de seguridad de `criterios-1.2` está fuera del esquema y de los criterios, como pide K44.

## Consumo aproximado
Acciones de herramienta: ~47 de 60. KB leídos: ~175 (protocolo 6, §K 15, INDICE ~20, reporte 8, diffs y criterios por sed/grep ~95, salidas de scripts ~30). KB escritos: ~19 (informe 12, reporte 4, scripts 3). Sin paquetes instalados.
