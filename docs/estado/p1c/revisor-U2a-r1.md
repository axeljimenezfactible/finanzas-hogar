# revisor-U2a-r1 · Reporte (flujo P1c, puerta 1b, 11-oct-2026)

## Hecho
- Revisé la colorimetría 1.2 contra `docs/criterios/U2.md` 1.2 (U2-02 a U2-09) y §K (K26, K27, K41 y K42). Informe completo: `docs/estado/p1c/REVISION-U2a-r1.md`.
- Veredicto: **NO APROBADO**, con 0 hallazgos altos, 1 medio (M-01: §12 todavía presenta a Karen y Bodega Aurrera como «propuestas sujetas a la puerta 2», contra K41) y 2 bajos (B-01 y B-02).
- Volví a ejecutar las verificaciones (a) a (f) con comandos y scripts propios. Las seis pasan. Corrí además 4 pruebas negativas sobre las comprobaciones del script y revisé el estado de los 6 hallazgos de la ronda anterior.

## Archivos
- `docs/estado/p1c/REVISION-U2a-r1.md` (nuevo)
- `docs/estado/p1c/revisor-U2a-r1.md` (nuevo, este reporte)

## Verificación
- `make verify-color` con `pipefail`: salida 0, «VERIFICACIÓN OK», 202 filas con 0 en cumple=false, 4 archivos idénticos byte a byte.
- Regeneración propia en un árbol espejo (`python3 -I <espejo>/colorimetria.py --sin-png`): salida 0 y los 4 archivos idénticos al repositorio.
- Script propio contra `git show HEAD:docs/contratos/C2-colores.json`: 190 tokens con los mismos nombres; 0 diferencias en `$value`, `$type` y `$extensions`; `$description` cambia solo en los 5 tokens de K41 y K42.
- `grep -c PROVISIONAL`: 0 en los 5 archivos. En el script solo queda el mecanismo. En `docs/` fuera de `estado/`, ningún artefacto vigente cita como PROVISIONAL a los bloqueos resueltos.
- OKLab propio (matrices de Ottosson, ida y vuelta) de los 5 tokens y de Ximena; pares frente a Ximena; escalonamiento de las tiendas en oscuro (ΔL mínimo 0.1008). Todo coincide con el documento.
- Barrido propio de U2-08 (roles, acción crítica, categóricas y marca): 6 pares en claro y 5 en oscuro, los mismos 11 de §7.7. Ningún par queda entre una categórica o la marca y un rol.
- U2-04: recalculé 32 pares de contraste, todos dentro de ±0.01 (el CSV trunca, como está documentado). U2-05, U2-06 y U2-07: pares recalculados, sin cambios. U2-09: ΔL consecutivo mínimo 0.174 y 0.218 en OKLab, contraste mínimo 3.15 y 3.11.
- Pruebas negativas (hex de `FIRMES_12`, marca PROVISIONAL, L de Walmart, reasignar a K42): las 4 salen con código 1.

## Supuestos
- U2-03 no se reevaluó a fondo: los 190 hex y `$extensions.oklch` son idénticos a HEAD, así que su estado es el de `docs/estado/p1/REVISION-U2a-r2.md`. Por el mismo motivo, en U2-09 no recalculé la simulación de daltonismo con culori.
- Las partes de U2-02 a U2-09 que piden evidencia en `DESIGN.md` (líneas «U2-0x», tabla rol → matiz) son de WF-1. Las evalué sobre los artefactos de U2a.
- Califiqué M-01 como media: es una incoherencia entre artefactos (`colorimetria.md` §12 frente a K41 y al propio §7.1–§7.2), no una marca PROVISIONAL ni un valor cambiado.

## Bloqueos
- Ninguno.

## Pendientes
- Ninguno en mi alcance. Las correcciones de M-01, B-01 y B-02 le tocan al agente de U2a (`docs/design/tools/colorimetria.py` y regeneración).

## Solicitudes fuera de alcance
- Ninguna nueva. Se mantiene la de U2a: `docs/design/DESIGN.md` (WF-1) debe recoger las dos reglas de «Cambios 1.2» y el hallazgo H-04 de la ronda anterior.

## Propuestas para Axel
- Ninguna. Esta revisión no deja valores K31 ni K44 a propuesta.

## Consumo aproximado
- Acciones de herramienta: unas 33.
- Leído: unos 140 KB, siempre por partes (protocolo, §K, colas de los bloqueos, criterios U2, reporte U2a-1.2, diff del script, secciones de `colorimetria.md` y salidas de los scripts).
- Escrito: unos 30 KB (scripts temporales en `datos/tmp/rev-u2a-r1/`, ya borrados, y los dos informes).
- Temporales: borré `datos/tmp/rev-u2a-r1/`. `datos/tmp/verify-wf0/` es la carpeta de trabajo que el propio `Makefile` recrea en cada `make verify-*`. No la borré porque otros agentes del flujo pueden estar usándola al mismo tiempo.
