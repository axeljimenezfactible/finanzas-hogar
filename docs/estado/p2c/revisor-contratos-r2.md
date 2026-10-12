# P2c · revisor-contratos-r2 · 11-oct-2026, noche

Revisor independiente, ronda 2 de 3, de los contratos de motores de la puerta 2. El informe completo está en `docs/estado/p2c/REVISION-contratos-r2.md`.

## Hecho
- Veredicto **NO APROBADO**: 1 alto, 2 medios y 4 bajos.
  - **A-1:** `make verify-all` sigue roto en `verify-c1` («hora HH:MM»), fuera del alcance de los contratos.
  - **M-1:** 3 bloqueos ABIERTOS que citan los contratos no están en INDICE ni en K46.
  - **M-2:** a la lista para WF-3 le faltan 7 comentarios PROVISIONAL del código.
- Hallazgos de la ronda 1:
  - Cerrados: A-1, M-1 y M-2. M-3 también, en el contrato; su registro pasa a M-1.
  - Sigue abierto: A-2, que ahora es A-1.
  - Cerrado en parte: M-4, que ahora es M-2.
- Lo que cumple: (a) Q2-01 a Q2-08 aplicadas; (b) ninguna «sin recomendación» decidida; (c) 0 marcas PROVISIONAL activas que citen un bloqueo resuelto; (d) ninguna cifra canon cambió.

## Verificación (comandos propios en `datos/tmp/revisor-contratos-r2/`, ya borrados)
- `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP=…/datos/tmp/revisor-contratos-r2/vw`: rc=2.
  - `verify-c1`: «gen2: tipo de parámetro desconocido: hora HH:MM», `Makefile:130`.
  - Lo demás OK: WF-0; S1 TODO OK; U1 63 OK; criterios 178 con 0 diferencias; K25 ≥ 89.2 %; oráculos S2 30/30, S3 906/906 (canon 437/437), S4 26/26 y S5 29/29; `model.json` intacto; WF-1 OK.
- `prov.py`: cada renglón con «provisional» de S2 a S5 contra el estado de su bloqueo. Hay 0 marcas activas que citen un bloqueo no ABIERTO; las demás coincidencias son bitácora o negaciones.
- `fix.py`: 25 fixtures cuyo `provisional` cita un bloqueo RESUELTO, todas cubiertas por las listas o por §14.1. `s3prov.py` y `s3where.py`: el origen de cada lista en el `esperado` de S3.
- `code.py`: 62 sitios de código con PROVISIONAL y un bloqueo resuelto; 7 no están en la lista (M-2). La receta `grep` del reporte no los encuentra.
- `nums.py`: cifras de `HEAD` contra el árbol de trabajo en S2 a S5; no desaparece ninguna.
- `frag.py`: lo que cambió la corrección r1, por fragmentos, contra su respaldo `/tmp/p2c-ccr1/bak/`.
- Conteos: `tabla.csv` tiene 150 de 150 renglones `propuesto`. `grep` en INDICE y en decisiones: 0 menciones de `S3-formula-tarjeta-pdf`, `S3-iva-lectura`, `S4-empates-sin-confirmar` y `S1-reparto-edicion-gasto`.

## Supuestos
1. El estado de un bloqueo sale de su última línea «Estado del bloqueo (puerta 2)» o, si no la tiene, de sus líneas «**Estado:**» (INDICE, regla desde la puerta 2).
2. Cuento los renglones de bitácora («Cambios de la versión…», M-r…) como historia y no como marca activa.
3. A-1 y M-1 están fuera del alcance de los contratos, pero la regla de severidad los cuenta igual (verificación global rota; inconsistencia entre artefactos).

## Bloqueos
Ninguno nuevo.

## Pendientes
- Integrador o agente de C1: A-1.
- Integrador o agente de bloqueos: M-1.
- Corrector de contratos: M-2 y los bajos.
- Después, la ronda 3.

## Propuestas para Axel
Ninguna nueva. Siguen las de la corrección r1: contestar `S3-iva-lectura` y Q2-90 por id, y entregar los PDF de `S3-formula-tarjeta-pdf`.

## Consumo aproximado
55 acciones de herramienta · unos 170 KB leídos · unos 16 KB escritos.
