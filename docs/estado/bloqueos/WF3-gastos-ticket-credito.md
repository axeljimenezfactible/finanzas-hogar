# WF3-gastos-ticket-credito · Pagar con crédito desde P-05

Abierto por `front-gastos-presupuesto` (WF-3a, 12-oct-2026). Estado: PREGUNTA para Axel.

**Qué falta.** `confirmTicketScan` acepta `credit_purchase: CreditPurchaseCreate`, que exige `credit_id`, `kind` (`installment_plan`, `single_charge` o `loan_disbursement`), `principal_cents` y `term_count`, y opcionalmente intereses, frecuencia y primer pago. U1 (flujo del ticket, paso 6) solo dice «elige la línea y teclea intereses cotizados y plazo». No está decidido qué `kind` usa la pantalla, si se piden la frecuencia y la fecha del primer pago, ni si el capital es el total del ticket.

**Opciones.** (a) `kind=installment_plan`, capital = total, `interest_source=quoted`, y se piden la línea, los intereses y el plazo (2 toques más, como en U1); (b) las mismas, más la frecuencia y el primer pago; (c) mandar a P-19 después de confirmar el ticket de contado.

**Impacto hoy.** En P-05 la forma de pago ofrece Efectivo, Débito y Transferencia; Crédito no aparece, y un aviso lo explica (PROVISIONAL).
