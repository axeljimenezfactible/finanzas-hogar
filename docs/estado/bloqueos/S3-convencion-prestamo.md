# Bloqueo S3-convencion-prestamo
Contexto: agente contrato-S3 (WF-1) al fijar la tabla francesa del préstamo personal (`S3-interfaz.md` §5.3). `credit` guarda la tasa anual (`annual_rate_num/den`), la frecuencia (`payment_frequency`) y el IVA (`vat_on_interest_bp`), pero ninguna regla dice cómo pasar la tasa anual a la tasa de cada periodo ni dónde entra el IVA. S3-09 solo exige que la tabla cuadre al centavo.
Qué falta: ¿Cómo se calcula la tasa por periodo de un préstamo a plazo fijo y cómo entra el IVA?
Opciones:
  (a) Tasa nominal: anual / 12 (mensual), / 24 (quincenal), / 52 (semanal); el IVA multiplica la tasa del periodo (r × 1.16) y el pago fijo se calcula con esa tasa; intereses del periodo = HALF_UP(saldo × r) con IVA incluido. Consecuencia: convención usual de las tablas francesas en México; cifras reproducibles.
  (b) Base diaria: días reales entre pagos × anual / 360, IVA aparte en cada pago. Consecuencia: pagos casi iguales pero no idénticos; el «pago fijo» deja de serlo.
  (c) Axel captura del contrato el pago fijo y el motor solo reparte capital e intereses. Consecuencia: sin fórmula de tasa; la tabla cuadra con el pago dado.
Impacto: `Loan.Amortize` (interfaz §5.3) y la fixture `s3-prestamo-frances`. Las propiedades de S3-09 (Σ amortización = capital, saldo final 0.00, Σ pagos = capital + Σ intereses, último pago ajusta) se cumplen con cualquier opción. Ninguna cifra canon depende de esto.
Cómo seguí: la interfaz usa la opción (a) como PROVISIONAL citando este bloqueo; la fixture usa datos de fixture (no del préstamo real de Axel) y sus valores esperados los calcula el oráculo con la misma convención.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-28 · ¿Cómo se calcula la tasa de cada periodo de un préstamo a plazo fijo y cómo entra el IVA?
- **Respuesta:** **(a)** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Tasa nominal: anual / 12 (mensual), / 24 (quincenal), / 52 (semanal); el IVA multiplica la tasa del periodo (r × 1.16) y el pago fijo se calcula con esa tasa; intereses del periodo = HALF_UP(saldo × r) con IVA incluido.»
- **Consecuencia:** convención usual de las tablas francesas en México; cifras reproducibles.
- **Qué cambia:** `Loan.Amortize` (`S3-interfaz.md` §5.3) y la fixture `s3-prestamo-frances` pierden la marca PROVISIONAL (la fixture sigue con datos de fixture, no del préstamo real de Axel).
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
