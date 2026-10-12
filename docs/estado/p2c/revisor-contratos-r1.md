# P2c · revisor-contratos-r1 · reporte · 11-oct-2026, noche

Informe completo: `docs/estado/p2c/REVISION-contratos-r1.md`. **Veredicto: NO APROBADO (2 altos, 4 medios, 5 bajos).**

## Hecho
- Revisé el diff completo de `S2-interfaz.md`, `S4-interfaz.md`, `S5-interfaz.md` y `S3-interfaz.md`, la pregunta 4 de `S2-ventanas.md`, los 4 bloqueos nuevos y el reporte `contratos-puerta2.md`. Lo comparé con PUERTA-2 §0, §2 y §3, `RESPUESTAS.md` y K46.
- Saqué las cifras S2 con el oráculo: 97375, 292125, 50000 y 0. Confirmé en el código del motor tres reglas PROVISIONAL «en uso»: respaldo de presentación, metas vencidas y redondeo a granel. También los motivos de rechazo de PROFECO.
- Conté por segmento las marcas PROVISIONAL de los 4 contratos y revisé el estado de cada bloqueo citado. Revisé el `provisional` de las fixtures de S2, S4 y S5.

## Verificación (comandos y salidas)
- `python3 -I datos/tmp/revisor-contratos-r1/s2cifras.py .`: `P1 … 97375`; `d,r … 194750, 292125`.
- `python3 -I datos/tmp/revisor-contratos-r1/s2ac.py .`: A + C deja C con devengado 50000 y diferencia 0; A + B + C deja C con devengado 0; en meses, diciembre da 0.
- `python3 -I datos/tmp/revisor-contratos-r1/prov.py . <4 contratos>`: S2 solo cita bloqueos ABIERTOS. S3, S4 y S5 citan bloqueos RESUELTOS (tabla en el informe, fila c).
- `make verify-all GOOSE=… TMP=datos/tmp/revisor-contratos-r1/vw`: rc=2. WF-0 OK y luego `verify-c1`: `gen2: tipo de parámetro desconocido: hora HH:MM`.
- `make -k verify-criterios verify-wf1-cobertura verify-wf1-oraculos …`: rc=0. Criterios 261/261 con 0 diferencias; K25 OK (≥ 89.2 %); S2 30/30, S3 27, S4 26/26 y S5 29/29, con 0 diferencias; `model.json` y `verify3.json` intactos.

## Supuestos
1. En (c) no cuento como marca vigente la bitácora de versiones anteriores («Cambios de la 1.x»).
2. A-1 y A-2 están fuera del rango que se dio a `contratos-puerta2`. Los califico por su efecto en el artefacto y en la verificación global, sin culpar al agente.

## Bloqueos
- Ninguno nuevo. Propongo reabrir las preguntas 1 a 3 de `S3-formula-tarjeta` (M-3).

## Pendientes
- Ronda 2: A-1, A-2 y M-1 a M-4. Los bajos son opcionales.

## Propuestas para Axel
- Ninguna nueva. B-4 sugiere que la elección del supuesto 3 de `contratos-puerta2` (diferencia `null` en las tiendas incompletas, sobre «puede ser negativa») entre en la lista de posibles excepciones.

## Solicitudes fuera de alcance
- A-2 es para S1 o el integrador: el generador C1 no conoce el tipo «hora HH:MM».
- M-3 es para el agente de bloqueos o el integrador.

## Consumo aproximado
53 acciones · unos 190 KB leídos · unos 16 KB escritos (informe, reporte y scripts temporales ya borrados).
