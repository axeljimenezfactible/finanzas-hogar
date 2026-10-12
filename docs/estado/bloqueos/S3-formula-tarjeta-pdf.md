# Bloqueo S3-formula-tarjeta-pdf

Abierto el 11-oct-2026, noche, en la ronda P2c por el agente `correccion-contratos-r1` (hallazgo M-3 de `docs/estado/p2c/REVISION-contratos-r1.md`). Aplica la regla de conteo de `docs/estado/bloqueos/INDICE.md`: una decisión tomada pero sin resultado no cuenta como resuelta.

Contexto: en la puerta 2, Axel contestó (a) a las preguntas 1 a 3 de `S3-formula-tarjeta` (Q2-29 a Q2-31, aprobadas en bloque; `docs/estado/puerta2/RESPUESTAS.md` y K46). Con (a), S3 escribe la fórmula de intereses y de pago mínimo de la tarjeta y del crédito de tienda a partir de un estado de cuenta real de cada producto, con un ejemplo comprobable. `S3-formula-tarjeta` quedó RESUELTO por la decisión, pero el resultado todavía no existe: la fórmula espera los PDF.

HECHO (11-oct-2026):
- La carpeta `datos/semilla/estados-cuenta/` no existe en el árbol de trabajo. No se ha entregado ningún estado de cuenta.
- `docs/estado/puerta2/PUERTA-2.md` §7 deja la acción pendiente para Axel: «con (a) en Q2-29 a Q2-31, entregar el PDF de un estado de cuenta de tarjeta y uno de crédito de tienda».
- Mientras tanto, `docs/contratos/motores/S3-interfaz.md` 1.5 (§5.3, §5.5, §5.8, §6 y §14.1) usa la fórmula general de la opción (b) de `S3-formula-tarjeta`. Con base diaria es Σ saldo diario × tasa anual / 360 y con base mensual, saldo al corte × tasa anual / 12. El IVA se suma según `vat_included` (Q2-67 = 1), y el mínimo es el mayor entre % × saldo y el mínimo fijo. Todo eso va marcado PROVISIONAL y cita este archivo.

Qué falta: los dos PDF (una tarjeta y un crédito de tienda) en `datos/semilla/estados-cuenta/`, fuera de git. Con ellos, S3 escribe la fórmula de cada producto con su ejemplo y quita la marca.

Opciones (la respuesta ya es (a); esto es lo que puede pasar ahora):
1. Axel entrega los dos PDF. S3 escribe la fórmula con fuente, S3-08 y S3-10 se cumplen completos y este bloqueo se cierra.
2. Axel cambia por excepción la respuesta de Q2-29, Q2-30 o Q2-31: a (b), y la fórmula general queda como regla, o a (c), y el motor responde «falta capturar» en vez de calcular.

Sin recomendación: la respuesta (a) ya es de Axel y solo falta su entrega.

Impacto:
- El interés calculado con la fórmula y todo pago mínimo de `CardPeriodInterest` y `CardMinimumPayment` (consultas `interes_tarjeta` y `minimo_tarjeta`).
- El `ErrUnsupported` de una compra de tienda con interés por datos K18.
- Los cortes `computed` del calendario unificado.
- Las fixtures `s3-tarjeta-minimo` y `s3-tienda-abonos-interes`.

Ninguna cifra canon depende de esto. Son FIRMES el pago total a tiempo ($0.00), los meses sin intereses, los abonos sin interés y el calendario de la tarjeta (Q2-32 = A).

Cómo seguí: las marcas PROVISIONAL de S3 que citaban `S3-formula-tarjeta` (preguntas 1 a 3) ahora citan este archivo. Su nombre sustituye a `S3-formula-tarjeta` en las listas `Provisional` de §14.1 y en el `provisional` de las dos fixtures (cambio para WF-3).

- **Estado:** ABIERTO
