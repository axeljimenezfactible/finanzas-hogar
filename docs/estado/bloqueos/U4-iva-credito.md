# Bloqueo U4-iva-credito
Contexto: U4a (WF-1) especificaba en Ajustes › Créditos y condiciones los cinco datos de K18 de la tarjeta bancaria y del crédito de tienda (tasa anual, base diaria o mensual, si incluye IVA, porcentaje de pago mínimo y mínimo fijo), con el aviso «falta capturar» mientras falte alguno.
Qué falta: C1 v1.2 y el esquema 1.3 guardan «si incluye IVA» en `credit.vat_on_interest_bp`, donde NULL significa «la tasa ya incluye IVA»; por eso `Credit.missing_terms` solo puede listar 4 datos (`annual_rate`, `rate_base`, `min_payment_pct_bp`, `min_payment_fixed_cents`) y nunca avisa que falta el del IVA. ¿Cómo se distingue «no capturado» de «incluye IVA»?
Opciones:
1. Columna nueva `credit.vat_included` (booleano que admite NULL = falta capturar) y `vat` en `missing_terms`; `vat_on_interest_bp` solo cuando `vat_included = 0`. Consecuencia: migración 00004 y C1 v1.3; la pantalla cumple K18 completo.
2. Usar `vat_on_interest_bp = 0` para «incluye IVA» y NULL para «falta capturar». Consecuencia: cambia el significado de NULL en esquema, C1 y motor S3 (la semilla Kueski pasaría a 0); sin columna nueva.
3. Dejarlo como está: «incluye IVA» por omisión y sin aviso para ese dato. Consecuencia: K18 queda parcial (4 de 5 datos con «falta capturar»).
Impacto: el campo «¿la tasa incluye IVA?» de P-26.05 queda PROVISIONAL en `docs/design/U4-ajustes-a.md`; no cambia ninguna cifra del oráculo (Kueski ya incluye IVA, decisión C).
Cómo seguí: el campo se especifica con la opción sí o no y el porcentaje de IVA, escribiendo `vat_on_interest_bp` como hoy, marcado PROVISIONAL con cita de este archivo; el aviso «falta capturar» cubre los otros 4 datos.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-67 · ¿Cómo se distingue «no capturado» de «la tasa ya incluye IVA»?
- **Respuesta:** **1** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Columna nueva `credit.vat_included` (booleano que admite NULL = falta capturar) y `vat` en `missing_terms`; `vat_on_interest_bp` solo cuando `vat_included = 0`.»
- **Consecuencia:** migración 00004 y C1 v1.3; la pantalla cumple K18 completo.
- **Qué cambia:** Migración 00004 (`credit.vat_included`, admite NULL = falta capturar), C1 v1.3 (`Credit.missing_terms` incluye `vat`), motor S3 y la semilla; campo «¿la tasa incluye IVA?» de P-26.05 en `docs/design/U4-ajustes-a.md` sin marca PROVISIONAL. Ninguna cifra cambia: Kueski ya incluye IVA.
- **Estado:** RESUELTO

Remisión (P2c, 11-oct-2026, noche; la aplicó `p2c-remanentes` el 12-oct-2026): Q2-67 = 1 es FIRME. Los dos casos que no cubre (préstamo con `vat_included` NULL; `vat_included` = 0 sin porcentaje) viven en `docs/estado/bloqueos/S3-iva-lectura.md` (ABIERTO, lo abrió `correccion-contratos-r1`).

**Estado del bloqueo (puerta 2): RESUELTO**
