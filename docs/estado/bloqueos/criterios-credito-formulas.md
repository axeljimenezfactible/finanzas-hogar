# Bloqueo criterios-credito-formulas
Contexto: S3-08 y S3-10: decisiones C pide modelar tarjeta bancaria (corte, pago, mínimo, meses sin intereses) y crédito de tienda con abonos. Solo hay criterio numérico para «intereses 0 si se paga total antes de fecha de pago» y para el préstamo (tabla francesa).
Qué falta: Fórmulas de: pago mínimo de tarjeta (¿% del saldo, mínimo fijo, ambos?), intereses de tarjeta cuando no se paga el total (tasa anual, base diaria, IVA), y comisión o interés del crédito de tienda con abonos. Datos reales de las tarjetas y tiendas de Axel (tasas, fechas) no están en el repositorio.
Opciones:
- (a) Axel captura por producto la tasa y el método de su estado de cuenta; el motor solo implementa el cálculo con esos parámetros.
- (b) Axel entrega 1 estado de cuenta de cada producto para fijar las fórmulas con ejemplo numérico.
- (c) Implementar solo el caso 0 intereses y MSI sin interés, y dejar el resto como captura manual de pagos.
Impacto: Los casos con intereses de tarjeta y de crédito de tienda quedan PROVISIONAL; no se inventa tasa ni mínimo.
Cómo seguí: S3-08 y S3-10 exigen los casos no ambiguos (0 intereses, MSI, abonos sin interés, francesa) y marcar PROVISIONAL el resto.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Opción (a): Axel captura por producto, desde su estado de cuenta, la tasa anual, la base (diaria o mensual), si incluye IVA, el porcentaje de pago mínimo y el mínimo fijo; el motor aplica la fórmula general; sin datos muestra «falta capturar». No se inventa tasa ni mínimo.
- **Estado:** RESUELTO. Se aplica en: criterios S3-08 y S3-10, S1 1.2 (columnas de tarjeta y crédito de tienda si faltan), S3, U4.
