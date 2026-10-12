# revisor-U4-r4 · WF-1 · 11-oct-2026

Revisor independiente de U4 (ronda 4). No escribí nada de lo revisado ni corregí nada. El informe completo está en `docs/estado/wf1/REVISION-U4-r4.md`.

## Hecho
- Volví a correr, con scripts propios, todas las verificaciones de U4-01 a U4-13 sobre `docs/design/U4-ajustes-a.md` 1.3 y `docs/design/U4-ajustes-b.md` 1.2.
- Comprobé el hallazgo R3-1 (media): **cerrado**. Los bajos R3-2 a R3-6 siguen abiertos.
- Veredicto: **NO APROBADO**. Hay 1 hallazgo medio nuevo, R4-1: el interruptor «activo» y la edición múltiple `Ctrl+/` del gasto fijo, y el interruptor «Activa» de las metas, cambian dinero sin vista previa. Hay además 2 bajos nuevos (R4-2 y R4-3).

## Archivos
- `docs/estado/wf1/REVISION-U4-r4.md` (creado).
- `docs/estado/wf1/revisor-U4-r4.md` (este reporte).
- Temporales en `datos/tmp/revU4r4/` (`ops.py`, `params.py`, `misc.py`, `oraculo.py`, `defaults.py`, `secs.py`): borrados al terminar.

## Verificación
- U4-01: 14 secciones en orden, nombres iguales a los de la decisión G, 7 bloques por sección y 0 celdas vacías (`secs.py`).
- U4-02: 210 `operationId` en los extractos y 210 en `openapi.yaml`; 0 citas inexistentes (98 en la parte a, 45 en la b); 68 de 68 operaciones de Ajustes cubiertas (`ops.py`).
- U4-03 y U4-04: lista de ≈ 53 elementos completa. Muestra al azar de 15 valores por defecto, con 0 sin origen en §A–§K (`defaults.py`).
- U4-05: A/B literal con `calculado_A` y `calculado_B`, mínimos de 8,094.50 y 5,000.00. Recalculé 800.50, 973.75, K40 (0.9836 y 0.4918), RN-17 (2,556, 2,056 y 1,799), IMC 23.7, RN-15 (1,486.571) y el reparto de $199.00 (9,950). Hueco: R4-1.
- U4-06: atajos y comandos «/» dentro de U1 §5.2 y §5.3, sin duplicados. Vistas coherentes con U1 §2.5.a, salvo R3-4. Edición en el lugar coherente con U1 §2.6.
- U4-07 a U4-11: catálogo V-xx; acciones críticas con `color.critical.*`; 0 hex; 0 ISO; 21 tokens existentes; 8 bloqueos existentes.
- U4-12: 15 de 15 `@param` en su sección, con valor inicial y «Vuelve a…». Los 4 de `params_json` llevan PROVISIONAL con `S1-parametros-avisos`. El `params_json` de 00001 da `weekday` 7 y `day_of_month` 1 (`params.py`).
- U4-13: claves `health.*` con 0 y 0, valor por persona con vacío = el de la casa, «mantener» sin campo y vista previa con 1,799.
- No aplican pruebas Go, `go vet` ni `gofmt`: U4 es un documento sin código.

## Supuestos
- Cuento como «cumple» las citas `projection.repeat_cadence_enabled = false` y `price.assisted_reading_enabled = false`: la clave va dentro de un mismo bloque de código junto con su valor.
- `color.role.warning` se cita como grupo; existe en `entrada/C2-tokens.json` con sus variantes `light` y `dark`.

## Bloqueos
- Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md`, que no se decide en WF-1.

## Pendientes
- Corregir R4-1 (media) en `docs/design/U4-ajustes-a.md` P-26.02 y `docs/design/U4-ajustes-b.md` P-26.13.
- Opcionales (bajos): R4-2, R4-3 y R3-2 a R3-6.

## Solicitudes fuera de alcance
- Ninguna.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna. La revisión no toca valores de K31.

## Consumo aproximado
- Acciones de herramienta: 31.
- Leído: ≈ 265 KB (protocolo, decisiones §A–§K, criterios U4, hallazgos de r3, `correccion-U4-r3`, las dos partes de U4 completas, U1 §2.3, §2.5, §2.6 y §5.3, RN-03 a RN-05, `planned_item` del esquema, `@param` de 00003 y salidas de los scripts).
- Escrito: ≈ 30 KB (6 scripts temporales, borrados, y 2 reportes).
