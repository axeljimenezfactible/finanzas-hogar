# Bloqueo S1-reposicion-parte-externa
Contexto: Al formalizar la regla de Karen (reglas-negocio.md RN-04: «paga directo a Kueski su parte y repone la línea cuando ella paga»). El oráculo solo repone capital por pago completo: el pago del 30-oct-2026 de 973.75 repone 773.63 (`model.json: despensas[0].pagos[0]`).
Qué falta: Cuando Karen paga su parte de un pago (374.93) por separado, antes o después que Axel, ¿cuánto capital repone ese pago en el disponible?
Opciones:
A. Proporcional al capital del pago: capital × parte pagada / monto del pago. Con Karen: 773.63 × 374.93 / 973.75 = 297.88, y Axel repone 475.75. Consecuencia: el disponible sube en dos pasos y al completarse el pago suma exactamente 773.63.
B. Solo al completarse el pago programado (las dos partes). Consecuencia: si Karen paga antes, la línea no sube hasta que paga Axel; es lo más conservador.
C. Primero su capital propio: la parte de Karen repone su capital por pago (1,191.50 / 4 = 297.875 → 297.88) y el titular el resto. Consecuencia: en la despensa 1 coincide con A; difiere cuando las partes no son proporcionales al capital.
D. Lo que diga Kueski cuando Axel confirme cómo repone el capital (va junto con el interruptor A/B). Consecuencia: queda en espera de esa confirmación.
Impacto: El disponible A/B entre pagos parciales y el estado de cuenta de Karen (S3, S7). El esquema no cambia: `credit_payment` guarda quién pagó, cuánto y cuándo.
Cómo seguí: Documenté la regla a nivel de pago completo (FIRME, sale del oráculo) y marqué PROVISIONAL cómo se divide el capital entre las partes.
