# Bloqueo S3-formula-tarjeta
Contexto: agente contrato-S3 (WF-1, 11-oct-2026) al fijar la interfaz del motor de créditos (`docs/contratos/motores/S3-interfaz.md` §5.2 y §5.4). K18 y RN-26 deciden QUÉ datos captura Axel por producto (tasa anual, base diaria o mensual, si incluye IVA, % de pago mínimo y mínimo fijo), pero la «fórmula general» que los usa no está escrita en `reglas-negocio.md` ni tiene fuente en el repositorio (S3-08 exige que esté escrita con ejemplo numérico y con fuente: estado de cuenta o documentación del producto). La red de los agentes solo admite PROFECO, USDA, Open Food Facts y documentación de herramientas, así que el agente no puede citar Banxico ni CONDUSEF.
Qué falta: (1) Intereses de tarjeta cuando no se paga el total: ¿sobre qué saldo (saldo promedio diario del periodo, saldo al corte, saldo insoluto después del pago) y con qué días (días reales del periodo / 360, / 365, o tasa anual / 12 en base mensual)? (2) Pago mínimo: ¿el mayor entre % × saldo y el mínimo fijo, o % × saldo + intereses + IVA con piso en el fijo? (3) Crédito de tienda con abonos que cobra intereses o comisión: ¿la misma fórmula que la tarjeta? Lo ideal es un estado de cuenta real de Axel (una tarjeta y una tienda) para escribir la fórmula con ejemplo.
Opciones:
  (a) Axel entrega un estado de cuenta de una tarjeta y de un crédito de tienda (PDF en `datos/semilla/estados-cuenta/`, fuera de git): S3 escribe la fórmula de ese producto con ejemplo numérico verificable. Consecuencia: fórmula con fuente; se cierra S3-08 y S3-10 completos.
  (b) Se adopta como fórmula general «saldo promedio diario × tasa anual / 360 × días del periodo» (base diaria) o «saldo al corte × tasa anual / 12» (base mensual), IVA sumado si `vat_on_interest_bp` no es NULL, y mínimo = máx(% × saldo, mínimo fijo), marcada PROVISIONAL hasta que un estado de cuenta la confirme. Consecuencia: el motor calcula, pero toda cifra de intereses y mínimo de tarjeta y tienda sale con la marca PROVISIONAL.
  (c) Mientras no haya fuente, el motor solo calcula lo que ya está decidido (intereses 0.00 si se paga el total a tiempo, meses sin intereses y abonos sin interés) y para lo demás responde «falta capturar». Consecuencia: S3-08 y S3-10 quedan parciales; nada inventado.
Impacto: `Card.InterestForPeriod`, `Card.MinimumPayment` y `StoreCredit` con intereses (interfaz §5.2 y §5.4); fixtures `s3-tarjeta-minimo` y la variante con interés de `s3-tienda-abonos`; criterios S3-08 (un día después de la fecha de pago) y S3-10 (con datos capturados). Ninguna cifra canon depende de esto.
Cómo seguí: la interfaz fija las firmas, los datos de entrada (K18) y el error `ErrMissingCapture` («falta capturar»); la fórmula queda descrita como candidata (opción b) con la marca PROVISIONAL citando este bloqueo. Pago total a tiempo (0.00), meses sin intereses y abonos sin interés no dependen de esto y quedan FIRMES.

Pregunta adicional (4), calendario de la tarjeta: ninguna regla escrita dice si una compra hecha el mismo día de corte entra en ese corte ni cómo se recorta el día de pago en meses cortos. Propuesta del agente (PROVISIONAL, mismo bloqueo): corte = primera fecha ≥ compra cuyo día es `statement_day` (recortado al último día del mes); fecha de pago = primera fecha > corte cuyo día es `payment_due_day` (recortado igual); los meses sin intereses pagan en las fechas de pago de los cortes sucesivos. Alternativa: la compra del día de corte entra en el corte siguiente.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-29 · (1) Intereses cuando no pagas el total: ¿sobre qué saldo y con qué días?
- **Respuesta:** **(a)** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Axel entrega un estado de cuenta de una tarjeta y de un crédito de tienda (PDF en `datos/semilla/estados-cuenta/`, fuera de git): S3 escribe la fórmula de ese producto con ejemplo numérico verificable.»
- Subpregunta literal del archivo: «(1) Intereses de tarjeta cuando no se paga el total: ¿sobre qué saldo (saldo promedio diario del periodo, saldo al corte, saldo insoluto después del pago) y con qué días (días reales del periodo / 360, / 365, o tasa anual / 12 en base mensual)?»
- **Consecuencia:** fórmula con fuente; se cierra S3-08 y S3-10 completos.
- **Qué cambia:** Con (a), S3 escribe la fórmula del producto con ejemplo numérico verificable cuando Axel entregue el PDF de un estado de cuenta de tarjeta y uno de crédito de tienda en `datos/semilla/estados-cuenta/` (fuera de git; acción pendiente de Axel, PUERTA-2 §7); mientras no lleguen, `Card.InterestForPeriod`, `Card.MinimumPayment` y `StoreCredit` con intereses (`S3-interfaz.md` §5.2 y §5.4) y las fixtures `s3-tarjeta-minimo` y `s3-tienda-abonos` (variante con interés) siguen como están (en uso (b)); S3-08 y S3-10 se cierran al llegar el PDF.
- **Estado:** RESUELTO

### Q2-30 · (2) Pago mínimo
- **Respuesta:** **(a)** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Axel entrega un estado de cuenta de una tarjeta y de un crédito de tienda (PDF en `datos/semilla/estados-cuenta/`, fuera de git): S3 escribe la fórmula de ese producto con ejemplo numérico verificable.»
- Subpregunta literal del archivo: «(2) Pago mínimo: ¿el mayor entre % × saldo y el mínimo fijo, o % × saldo + intereses + IVA con piso en el fijo?»
- **Consecuencia:** fórmula con fuente; se cierra S3-08 y S3-10 completos.
- **Qué cambia:** Con (a), S3 escribe la fórmula del producto con ejemplo numérico verificable cuando Axel entregue el PDF de un estado de cuenta de tarjeta y uno de crédito de tienda en `datos/semilla/estados-cuenta/` (fuera de git; acción pendiente de Axel, PUERTA-2 §7); mientras no lleguen, `Card.InterestForPeriod`, `Card.MinimumPayment` y `StoreCredit` con intereses (`S3-interfaz.md` §5.2 y §5.4) y las fixtures `s3-tarjeta-minimo` y `s3-tienda-abonos` (variante con interés) siguen como están (en uso (b)); S3-08 y S3-10 se cierran al llegar el PDF.
- **Estado:** RESUELTO

### Q2-31 · (3) Crédito de tienda con intereses o comisión
- **Respuesta:** **(a)** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Axel entrega un estado de cuenta de una tarjeta y de un crédito de tienda (PDF en `datos/semilla/estados-cuenta/`, fuera de git): S3 escribe la fórmula de ese producto con ejemplo numérico verificable.»
- Subpregunta literal del archivo: «(3) Crédito de tienda con abonos que cobra intereses o comisión: ¿la misma fórmula que la tarjeta?»
- **Consecuencia:** fórmula con fuente; se cierra S3-08 y S3-10 completos.
- **Qué cambia:** Con (a), S3 escribe la fórmula del producto con ejemplo numérico verificable cuando Axel entregue el PDF de un estado de cuenta de tarjeta y uno de crédito de tienda en `datos/semilla/estados-cuenta/` (fuera de git; acción pendiente de Axel, PUERTA-2 §7); mientras no lleguen, `Card.InterestForPeriod`, `Card.MinimumPayment` y `StoreCredit` con intereses (`S3-interfaz.md` §5.2 y §5.4) y las fixtures `s3-tarjeta-minimo` y `s3-tienda-abonos` (variante con interés) siguen como están (en uso (b)); S3-08 y S3-10 se cierran al llegar el PDF.
- **Estado:** RESUELTO

### Q2-32 · (4) Calendario de la tarjeta (corte y pago)
- **Respuesta:** **A (propuesta del agente, letras de PUERTA-2)** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «corte = primera fecha ≥ compra cuyo día es `statement_day` (recortado al último día del mes); fecha de pago = primera fecha > corte cuyo día es `payment_due_day` (recortado igual); los meses sin intereses pagan en las fechas de pago de los cortes sucesivos»
- **Consecuencia** (el archivo no la da; trae la alternativa que se descarta)**:** Alternativa: la compra del día de corte entra en el corte siguiente.
- **Qué cambia:** Calendario de la tarjeta en `S3-interfaz.md` (corte y fecha de pago, meses cortos, meses sin intereses) sin marca PROVISIONAL; fixtures `s3-unificado` y `s3-unificado-sin-corte`.
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
