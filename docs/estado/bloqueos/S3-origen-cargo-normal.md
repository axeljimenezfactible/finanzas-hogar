# Bloqueo S3-origen-cargo-normal
Contexto: agente correccion-contrato-S3-r2 (WF-1, 11-oct-2026), al fijar la tabla de despacho de `QuotePurchase` por tipo de crédito (`docs/contratos/motores/S3-interfaz.md` §5.8; hallazgo M3 (d) de `docs/estado/wf1/REVISION-contrato-S3-r2.md`). El cargo normal de una tarjeta (`credit_purchase.kind = single_charge`: un pago, intereses 0.00 si se paga el total a más tardar en la fecha de pago, S3-08) se guarda en `credit_purchase`, cuyo `interest_source` es NOT NULL y solo admite `quoted`, `term_table`, `promotion` y `amortization` (`backend/db/migrations/00001_esquema.sql`, tabla `credit_purchase`). Ninguno de los cuatro describe un cargo normal: no hay cotización por compra, ni tabla por plazo, ni promoción, ni amortización.
Qué falta: ¿qué `interest_source` lleva un cargo normal de tarjeta?
Opciones:
  (a) `quoted`: los intereses del cargo normal los fija el banco en su estado de cuenta (0.00 si se paga el total a tiempo; si no, los trae el corte), que se lee como «cotización del acreedor al pagar». Consecuencia: sin cambio de esquema; el significado de `quoted` se amplía y el diccionario lo anota.
  (b) Un valor nuevo (`none` o `statement`) en el CHECK de `credit_purchase.interest_source`. Consecuencia: migración nueva, diccionario y enumerado de C1; el motor emite ese valor.
  (c) `interest_source` admite NULL cuando `kind = single_charge`. Consecuencia: migración y CHECK nuevos; C1 lo marca anulable; el motor emite null.
Impacto: `Purchase.Source` del cargo normal (`S3-interfaz.md` §5.8) y, por arrastre, los escenarios de tarjeta del comparador (`s3-comparador-caso-2`), las fixtures `s3-unificado-sin-corte` y `s3-unificado` y el calendario unificado. Ningún monto ni ninguna cifra canon dependen de esto: los intereses del cargo normal (0.00) son FIRMES.
Cómo seguí: el motor emite `quoted` (opción a, propuesta del agente porque no toca el esquema) con la marca PROVISIONAL: toda salida del cargo normal lleva `S3-origen-cargo-normal` en su lista `provisional` (`S3-interfaz.md` §14.1) y las fixtures cuyo `esperado` trae ese `source` lo nombran en `provisional.que`. Las fixtures que solo prueban los 0.00 no piden `source`.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-37 · ¿Qué `interest_source` lleva un cargo normal de tarjeta?
- **Respuesta:** **(a)** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «`quoted`: los intereses del cargo normal los fija el banco en su estado de cuenta (0.00 si se paga el total a tiempo; si no, los trae el corte), que se lee como «cotización del acreedor al pagar».»
- **Consecuencia:** sin cambio de esquema; el significado de `quoted` se amplía y el diccionario lo anota.
- **Qué cambia:** `Purchase.Source` del cargo normal (`S3-interfaz.md` §5.8 y §14.1: se quita `S3-origen-cargo-normal` de la lista `provisional`); el diccionario de datos anota que `quoted` se amplía a «cotización del acreedor al pagar»; fixtures `s3-unificado` y `s3-unificado-sin-corte` y `s3-comparador-caso-2`. Los intereses $0.00 ya eran FIRMES.
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
