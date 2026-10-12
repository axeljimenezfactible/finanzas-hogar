# Bloqueo S2-metas
Contexto: interfaz del motor financiero, metas de ahorro (`S2-interfaz.md` §3.13; S2-12; RN-18; M-07 y M-16 de `C1-datos.md` §9).
Qué falta: (1) cómo se cuentan los «meses que faltan» de M-16 («aportación mensual necesaria = (monto meta − ahorrado) ÷ meses que faltan»); (2) qué porcentaje de avance se muestra: RN-18 dice «con un decimal» (50.0 %) y M-07/M-16 dicen «× 100 entero»; (3) qué valores tiene el «estado» de una meta en M-16.
Opciones: (1) A) del mes de hoy al mes de `target_date`, los dos incluidos; B) del mes siguiente al de hoy al de `target_date`; C) meses completos entre hoy y `target_date`. (2) A) un decimal en todas partes (RN-18); B) entero en todas partes (C1); C) un decimal en Metas y entero en la tarjeta de Inicio. (3) A) `en_curso`, `cumplida`, `vencida`; B) además `atrasada` cuando lo ahorrado va por debajo del plan acumulado.
Impacto: `months_left`, `needed_monthly_cents`, `plan` y `state` quedan PROVISIONAL, y con ellos la fixture `s2-metas`. El avance, lo ahorrado y el ajuste del último centavo (Σ plan = restante, S2-12) no dependen de esto.
Cómo seguí: la interfaz emite los dos porcentajes (`progress_tenths_pct` de RN-18 y `progress_pct` de C1), los dos calculados con HALF_UP desde la fracción exacta, y marca PROVISIONAL lo demás. Solicitud al integrador: alinear RN-18 con M-07/M-16 cuando Axel decida.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-16 · (1) ¿Cómo se cuentan los «meses que faltan» de una meta?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «A) del mes de hoy al mes de `target_date`, los dos incluidos»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: `months_left`, `needed_monthly_cents`, `plan` y `state` quedan PROVISIONAL, y con ellos la fixture `s2-metas`. El avance, lo ahorrado y el ajuste del último centavo (Σ plan = restante, S2-12) no dependen de esto.
- **Qué cambia:** `months_left` y `needed_monthly_cents` (M-16, `S2-interfaz.md` §3.13); fixture `s2-metas`; P-26.13 de U4; contesta también la antigua pregunta 2 de `U4-vista-previa-formulas`.
- **Estado:** RESUELTO

### Q2-17 · (2) ¿Qué formato tiene el porcentaje de avance?
- **Respuesta:** **C** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «C) un decimal en Metas y entero en la tarjeta de Inicio»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: `months_left`, `needed_monthly_cents`, `plan` y `state` quedan PROVISIONAL, y con ellos la fixture `s2-metas`. El avance, lo ahorrado y el ajuste del último centavo (Σ plan = restante, S2-12) no dependen de esto.
- **Qué cambia:** Alinear RN-18 (`reglas-negocio.md`, un decimal) con M-07 y M-16 (`C1-datos.md` §9, entero): un decimal en Metas (`progress_tenths_pct`) y entero en la tarjeta de Inicio (`progress_pct`); solicitud del archivo al integrador.
- **Estado:** RESUELTO

### Q2-18 · (3) ¿Qué estados puede tener una meta?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «A) `en_curso`, `cumplida`, `vencida`»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: `months_left`, `needed_monthly_cents`, `plan` y `state` quedan PROVISIONAL, y con ellos la fixture `s2-metas`. El avance, lo ahorrado y el ajuste del último centavo (Σ plan = restante, S2-12) no dependen de esto.
- **Qué cambia:** Campo `state` de M-16 y la fixture `s2-metas`: tres estados (`en_curso`, `cumplida`, `vencida`); P-26.13 de U4.
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
