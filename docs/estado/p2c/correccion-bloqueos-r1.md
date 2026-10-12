# P2c · correccion-bloqueos-r1 (11-oct-2026, noche)

Agente: correccion-bloqueos-r1. Alcance: los 27 bloqueos de `docs/estado/bloqueos/`, `INDICE.md`, `docs/plan/decisiones.md` §K y este reporte. Entrada: `docs/estado/p2c/REVISION-bloqueos-r1.md` y `docs/estado/p2c/bloqueos-respuestas.md`. No decidí ninguna pregunta «sin recomendación», no toqué `bloqueos-respuestas.md` (reporte de otro agente) ni nada fuera del alcance.

## Hecho
- HECHO (H-r1-01, alta): corrí `make verify-c1` y `make verify-s1` y anoto abajo la falla que causa la carpeta de bloqueos. La causa es `S1-parametros-avisos.md`, que ya trae «Respuesta de Axel» (Q2-11 = B) mientras otros archivos siguen citándolo como PROVISIONAL. Eso solo se arregla quitando las marcas fuera de mi alcance: queda como solicitud 1. No quité «Respuesta de Axel» ni cambié el archivo para que pase: sería ocultar el estado real.
- HECHO (H-r1-02, media): `INDICE.md` l. 26 ahora dice que es la regla de las puertas 1 y 1b y que no es la que cuenta desde la puerta 2; l. 28 agrega que las herramientas siguen la regla anterior (67 resueltos y 2 abiertos) hasta que se corrijan los scripts, y que `S4-presentacion-respaldo` está por crear (69 pasa a 70, abiertos a 7). La fila «Flujo WF-1» dice «24 con la sección «Respuesta de Axel (puerta 2)», 20 RESUELTOS», y nombra los 4 abiertos que traen la sección. Comprobé: 25 archivos tienen la sección, 1 es `S1-parametros-avisos` (P1c), 24 son de WF-1; de esos 24, 20 terminan en RESUELTO y 4 en ABIERTO.
- HECHO (H-r1-03, media): `S2-flujo-quincenal.md`, entrada Q2-12: la fuente del valor en uso ahora es `docs/contratos/motores/S2-interfaz.md` l. 250 y l. 447. Comprobé que la l. 250 es la fila de `variable_estimated_cents`/`savings_planned_cents` en `null` y la l. 447 liga al bloqueo.
- HECHO (H-r1-04, media): (i) `S2-ventanas.md`, caso 2 de Q2-01: agregué la línea «Qué falta» con la cita `docs/estado/wf1/REVISION-S2-r3.md` l. 38-48 y `PUERTA-2.md` §2.1, sin inventar opciones. (ii) K46 de `decisiones.md` cita `S4-presentacion-respaldo` como bloqueo por crear. (iii) Este reporte trae «Archivos» y «Solicitudes fuera de alcance».
- NO HECHO a propósito: H-r1-05 (baja, «2,921.25 y 3,895.00» sin signo de pesos en `U3-t2-acumulados.md`) y H-r1-06 (baja, RESPUESTAS.md vs. Q2-12). La orden era alta y media; el primero es un cambio trivial que queda como pendiente, el segundo es del orquestador.

## Archivos
- `docs/estado/bloqueos/INDICE.md` (l. 26, l. 28 y fila «Flujo WF-1»)
- `docs/estado/bloqueos/S2-flujo-quincenal.md` (cita de Q2-12)
- `docs/estado/bloqueos/S2-ventanas.md` (línea «Qué falta» del caso 2)
- `docs/plan/decisiones.md` (K46, Q2-07)
- `docs/estado/p2c/correccion-bloqueos-r1.md` (este reporte)

## Verificación
Árbol de trabajo actual (`make verify-c1 GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`): `parse_schema: clave @param repetida` y `Error 1`; `make verify-s1` termina con `RESUMEN FALLAS: ['S1-16', 'NEG12', 'K44', 'PROG13']`. Las fallas NEG12 y PROG13 y el error de `parse_schema` vienen del trabajo en curso de otros agentes (`backend/db/migrations/00004_puerta2.sql` sin seguimiento y cambios en `gen_dict.py`/`diccionario-datos.md`), no de los bloqueos. S1-16 y K44 sí vienen de los bloqueos (aparecen también aislados, abajo).

Aislado, para separar mi carpeta del trabajo ajeno: copié `git archive HEAD` a `/tmp/p2c-cbr1/a` («HEAD») y a `/tmp/p2c-cbr1/b` («HEAD + los 28 archivos de bloqueos del árbol + decisiones.md»):
| Copia | verify-c1 | verify-s1 |
|---|---|---|
| a · HEAD | `RESULTADO: TODO OK` | `S1-16 OK`, `K44 OK`; `RESUMEN FALLAS: ['DDL']` |
| b · HEAD + bloqueos | `FALLA C1-15 toda marca PROVISIONAL remite a un bloqueo abierto (sin «Respuesta de Axel»)` · `marcas PROVISIONAL en openapi.yaml: 6 · … sin bloqueo abierto: [Notification/params_json, send_time, advance_days, NotificationUpdate/params_json, …]` · `RESULTADO: 1 falla(s)` | `S1-16 FALLA … líneas con PROVISIONAL: 19 (migración 6, diccionario 7, reglas 6) · bloqueos abiertos citados: [] · líneas sin bloqueo en su bloque=17`; `K44 FALLA … marcas PROVISIONAL que citan un bloqueo resuelto=[('migraciones','S1-parametros-avisos') ×5, ('reglas','S1-parametros-avisos') ×5]`; `RESUMEN FALLAS: ['S1-16','K44','DDL']` |
`DDL` falla en las dos copias porque `git archive` no trae `.git` y esa comprobación compara contra git HEAD: no es del cambio. La diferencia entre a y b es exactamente C1-15, S1-16 y K44, lo que confirma H-r1-01. Tras mis correcciones el resultado de b es el mismo: mis cambios (texto de INDICE, una cita, una línea «Qué falta» y K46) no pueden arreglar ni empeorar esas tres, y las corregirá quien quite las marcas. `make verify-all` global no se puede dar por bueno hasta entonces.
Conteo de H-r1-02 por script propio (archivos con «Respuesta de Axel (puerta 2» y su última línea «Estado del bloqueo (puerta 2)»): 25 con la sección, 1 de P1c, 24 de WF-1 (20 RESUELTO, 4 ABIERTO).

## Supuestos
- La copia de `INDICE.md` que dice «67 y 2» describe a las herramientas, no a la regla vigente: dejé las dos reglas visibles y dije cuál manda para cada cosa.
- La línea «Qué falta» del caso 2 no lista opciones porque PUERTA-2 §2.1 la declara «sin recomendación» y el revisor solo dio la redacción de la pregunta, no las opciones.

## Bloqueos
- Ninguno nuevo. Por crear (lo pide la solicitud 3): `docs/estado/bloqueos/S4-presentacion-respaldo.md`.

## Solicitudes fuera de alcance
1. **Quitar las marcas PROVISIONAL de Q2-11 = B** (H-r1-01; antes del commit de los bloqueos). Responsable: el agente que aplica Q2-11 = B (ola de P2c de esquema/contratos). Rutas: `backend/db/migrations/00003_*.sql` o su continuación (6 líneas), `docs/contratos/diccionario-datos.md` (7), RN-29 y reglas (6), `openapi.yaml` por `c1gen` (6 marcas: `Notification.params_json`, `send_time`, `advance_days`, `NotificationUpdate.params_json`…), `docs/contratos/C1-datos.md` (P-26.08, P-26.09 y las líneas de `params_json`) y los criterios C1-22, U4-12 y S7. Motivo: `verify_c1.py` l. 322-323 y `s1_verifica.py` l. 242 y 575 dan por resuelto todo bloqueo con «Respuesta de Axel» y por eso C1-15, S1-16 y K44 fallan mientras haya marcas. El integrador no hace commit de los bloqueos sin este cambio; se cierra cuando `make verify-all` pasa con ambos juntos.
2. **Hacer que los scripts lean la línea «Estado del bloqueo (puerta 2)»** (H-r1-02). Responsable: quien mantiene `backend/internal/api/c1gen/verify_c1.py` (l. 322-323) y `backend/db/tools/s1_verifica.py` (l. 242 y 575). Cambio: si el archivo trae esa línea, usar su última aparición (RESUELTO/ABIERTO); si no, la regla de «Respuesta de Axel». Motivo: 4 archivos ABIERTOS (`S2-flujo-quincenal`, `S2-ventanas`, `S5-carrito-comparativa`, `S5-terminos`) traen la sección y se cuentan como resueltos; la primera marca PROVISIONAL de WF-3 que los cite (Q2-12, Q2-48, caso 2 de Q2-01) fallaría C1-15 o S1-16 con el bloqueo abierto.
3. **Crear `docs/estado/bloqueos/S4-presentacion-respaldo.md`** (H-r1-04 ii). Responsable: el agente que aplique Q2-07 en el contrato S4/S5 (M-r5-2). Cambio: documentar como PROVISIONAL la variante sin recomendación de Q2-07 (la regla de respaldo que el motor ya aplica) con Qué falta, Opciones e Impacto (protocolo §2), y actualizar el conteo de `INDICE.md` (69 a 70, abiertos 6 a 7). Fuente: `RESPUESTAS.md` l. 17 y `PUERTA-2.md` Q2-07. Lo pedí a quien aplique Q2-07 porque yo no decido la variante.
4. **Redactar la pregunta 4 completa de `S2-ventanas`** (caso 2 de Q2-01; H-r1-04 i). Responsable: quien aplique Q2-01 (motor S2, `ledger.go`, `fin_devengo.py`, `S2-interfaz.md` §2.12 y §3.2). Cambio: escribir la pregunta con sus opciones y cifras ($973.75 contra $2,921.25 al 2026-11-09, `REVISION-S2-r3.md` l. 45) y la regla PROVISIONAL; el bloqueo sigue ABIERTO y pasa a WF-3 hasta que Axel decida.
5. **H-r1-05** (baja): en `docs/estado/bloqueos/U3-t2-acumulados.md`, entrada Q2-65, «Qué cambia», escribir «$2,921.25 y $3,895.00». Responsable: quien cierre la ronda de bloqueos.
6. **H-r1-06** (baja): el orquestador confirma Q2-12 con Axel o anota en `RESPUESTAS.md` l. 18 que queda abierta (hoy esa lista no incluye Q2-12 y los bloqueos, INDICE y K46 sí la dejan abierta).

## Pendientes
- Solicitudes 1 a 6; hasta que la 1 se aplique, `make verify-c1` y `make verify-s1` no pasan con los bloqueos del árbol.
- Re-revisión de la ronda 2 de bloqueos.

## Propuestas para Axel
- PROPUESTA: ninguna nueva. Siguen abiertas Q2-12 (C o A), el caso 2 de Q2-01, la variante de Q2-07 y la decisión de `S5-terminos` (B, C o D) del reporte anterior.

## Consumo aproximado
Acciones de herramienta: 14 · KB leídos: ~45 · KB escritos: ~9.
