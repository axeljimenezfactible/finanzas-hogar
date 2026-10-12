# Bloqueo S2-flujo-quincenal
Contexto: interfaz del motor financiero, flujo por quincena (`S2-interfaz.md` §3.11; métricas M-02 y M-10 de `C1-datos.md` §9). M-10 incluye «variable estimado (método de S2)» y «ahorro planificado» por quincena; ningún documento (decisiones §A–§K, plan, reglas) fija esos métodos ni cómo se reconoce que una ocurrencia planificada ya se capturó.
Qué falta: (1) método del gasto variable estimado de una quincena futura; (2) cómo se reparte el ahorro mensual planificado (`planned_monthly_cents`) entre las dos quincenas; (3) cuándo una ocurrencia planificada (sueldo, renta) se da por cumplida para no contarla dos veces.
Opciones: (1) A) presupuesto de las categorías variables del mes en proporción a los días de la quincena; B) promedio diario del gasto variable real de los últimos 3 meses × días de la quincena; C) 0 hasta que Axel lo decida. (2) A) todo en la quincena 1; B) mitad y mitad con HALF_UP en la 1 y el resto en la 2. (3) A) un gasto o ingreso con el mismo `planned_item_id` en la misma quincena (o el mismo mes si la frecuencia es mensual o mayor); B) solo por `planned_item_id` y fecha exacta.
Impacto: `variable_estimated_cents`, `savings_planned_cents` y la regla de cumplimiento quedan PROVISIONAL; el saldo de quincenas futuras depende de ellos. Las cifras canon (pagos de despensa por quincena, R-13…R-31) no.
Cómo seguí: la interfaz declara los campos con la marca PROVISIONAL y propone (3A) como lectura de trabajo; las fixtures canon no traen gasto variable, ahorro ni partidas.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-13 · (2) ¿Cómo se reparte el ahorro mensual planificado entre las dos quincenas?
- **Respuesta:** **B** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «B) mitad y mitad con HALF_UP en la 1 y el resto en la 2.»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: `variable_estimated_cents`, `savings_planned_cents` y la regla de cumplimiento quedan PROVISIONAL; el saldo de quincenas futuras depende de ellos. Las cifras canon (pagos de despensa por quincena, R-13…R-31) no.
- **Qué cambia:** `savings_planned_cents` por quincena (M-10 de `C1-datos.md` §9); `docs/contratos/motores/S2-interfaz.md` §3.11 escribe la regla sin marca PROVISIONAL; el saldo de las quincenas futuras (M-02, M-10) se construye en WF-3. Nota: PUERTA-2 §3.2 describe B como «mitad y mitad, con el centavo de redondeo en la 1»; el texto literal del archivo es el de arriba.
- **Estado:** RESUELTO

### Q2-14 · (3) ¿Cuándo una partida planificada (sueldo, renta) se da por cumplida?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «A) un gasto o ingreso con el mismo `planned_item_id` en la misma quincena (o el mismo mes si la frecuencia es mensual o mayor)»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: `variable_estimated_cents`, `savings_planned_cents` y la regla de cumplimiento quedan PROVISIONAL; el saldo de quincenas futuras depende de ellos. Las cifras canon (pagos de despensa por quincena, R-13…R-31) no.
- **Qué cambia:** La regla de cumplimiento de §3.11 de `S2-interfaz.md` deja de ser PROVISIONAL (en uso era 3A como lectura de trabajo); fixture `s2-liga-provisional` (quita la marca).
- **Estado:** RESUELTO

## Puerta 2 (11-oct-2026): sin respuesta
Fuente: `docs/estado/puerta2/RESPUESTAS.md` («Las preguntas sin recomendación quedan sin respuesta en esta puerta: PROVISIONAL con la opción que ya está en uso, y pasan a WF-3») y hoja de respuestas de `PUERTA-2.md` §8. Axel puede mandar una excepción por id cuando lea el documento.

### Q2-12 · (1) ¿Cómo se estima el gasto variable de una quincena futura?
- PROVISIONAL: sigue en uso `variable_estimated_cents` = `null` (no entra a ningún saldo), que es lo más cercano a la opción C «0 hasta que Axel decida» (`docs/contratos/motores/S2-interfaz.md` 1.3 §3.11, fila `variable_estimated_cents`, y §8, fila «Variable estimado de M-10», que lo liga a este bloqueo; se cita por sección porque el archivo sigue cambiando). La hoja de respuestas de PUERTA-2 §8 recomienda «C o A», dos opciones y no una: el modo «lo recomendado + excepciones» no permite escoger entre ellas sin inventar la respuesta de Axel. Axel puede contestarla por id (C o A); pasa a WF-3 con marca PROVISIONAL.
- **Estado:** ABIERTO

**Estado del bloqueo (puerta 2): ABIERTO** — queda abierta Q2-12 (la recomendada nombra dos opciones, C o A); Q2-13 y Q2-14 están resueltas.
