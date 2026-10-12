# Bloqueo S3-iva-lectura

Abierto el 11-oct-2026, noche, en la ronda P2c por el agente `correccion-contratos-r1`, al aplicar Q2-67 = 1 en `docs/contratos/motores/S3-interfaz.md` (hallazgo A-1 de `docs/estado/p2c/REVISION-contratos-r1.md`).

Contexto: Q2-67 = 1 (bloqueo `U4-iva-credito`, RESUELTO; K46) agrega la columna `credit.vat_included` en la migración `00004_puerta2.sql`. Sus valores: 1 = la tasa ya incluye IVA; 0 = se suma el de `vat_on_interest_bp`; NULL = falta capturar. RN-26 1.4 y la nota de la columna dicen qué hace S3 con la tarjeta y el crédito de tienda: mientras `vat_included` sea NULL, `vat` va en `missing_terms` y el motor responde «falta capturar». Quedan dos casos que ninguna fuente resuelve.

HECHO (11-oct-2026):
- `docs/design/U4-ajustes-a.md` (P-26.05) solo captura «¿la tasa incluye IVA?» para la tarjeta y la tienda. Del préstamo captura la tasa anual y nada sobre el IVA.
- La 00004 migra los créditos que ya existían con el significado de 1.3: con porcentaje → `vat_included` = 0; con NULL → 1.
- El CHECK de la 00004 acepta `vat_included` = 0 con `vat_on_interest_bp` NULL; solo rechaza un porcentaje junto con 1 o con NULL. RN-26 1.4 pone `vat` en `missing_terms` solo mientras `vat_included` sea NULL.
- Hoy el motor (`backend/internal/motor_cred/loan.go` l. 35) lee `vatBP` nil como «la tasa ya trae IVA» y no suma IVA. Todavía no conoce `vat_included`.

Qué falta (preguntas para Axel):
1. Préstamo a plazo fijo (`term_loan`) con `vat_included` NULL: ¿qué hace S3?
   - A) Lo mismo que en 1.3: no suma IVA, porque lee «la tasa ya incluye IVA». Es lo que el motor hace hoy.
   - B) Responde «falta capturar» (`ErrMissingCapture`), igual que con un préstamo sin tasa anual. Para eso U4 tiene que capturar el IVA del préstamo y C1 tiene que avisarlo.
   - C) La API no deja guardar un préstamo sin ese dato (se valida al guardar), así que el motor nunca recibe NULL.
2. Tarjeta, tienda o préstamo con `vat_included` = 0 y sin `vat_on_interest_bp`: ¿qué hace S3?
   - A) Responde «falta capturar» (`ErrMissingCapture`), y C1 también lista `vat` en `missing_terms` en ese caso (cambia RN-26 1.4).
   - B) La API no deja guardar 0 sin porcentaje (con el CHECK o con una validación; cambia el esquema o C1), y si el motor lo recibe, responde `ErrInvalid`.

Sin recomendación.

Impacto: solo afecta a préstamos sin el dato del IVA (pregunta 1) y a créditos marcados «no incluye IVA» sin porcentaje (pregunta 2). Ninguna cifra canon cambia. Kueski es línea revolvente y no usa estos datos, y los préstamos de las fixtures de S3 traen IVA 1600.

Cómo seguí: `S3-interfaz.md` 1.5 (§5, §5.3, §5.4, §5.8 y §14.1) marca PROVISIONAL, con cita de este archivo:
- Pregunta 1: la opción A, que es la que está en uso. Ese préstamo lleva `Purchase.Provisional` = [`S3-iva-lectura`].
- Pregunta 2: `ErrMissingCapture`. No da número, porque RN-26 dice que el motor no inventa tasa, y deja sin decidir si C1 lo avisa.

- **Estado:** ABIERTO
