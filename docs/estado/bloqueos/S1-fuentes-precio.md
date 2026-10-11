# Bloqueo S1-fuentes-precio
Contexto: al aplicar K22 en RN-12 (`docs/contratos/reglas-negocio.md` 1.2). `price.source` admite cinco fuentes: `ticket_manual`, `ticket_ocr` (OCR ya confirmado), `site_assisted` (lectura asistida del sitio, hoy apagada), `profeco` y `manual` (otra captura a mano, sin ticket). K22 decide tres: PROFECO nunca compite (solo referencia) y el ticket a mano y el OCR confirmado compiten con el mismo peso (propuesta del orquestador). No dice nada de `site_assisted` ni de `manual`.
Qué falta: ¿Compiten en la comparativa los precios de la lectura asistida de sitios y los capturados a mano sin ticket? Si compiten, ¿con el mismo peso que el ticket?
Opciones:
A. Los dos compiten con el mismo peso que el ticket y la misma regla de frescura. Consecuencia: más precios en la comparativa; un precio tecleado sin ticket pesa igual que uno con ticket.
B. Los dos solo como referencia, como PROFECO. Consecuencia: la comparativa usa solo tickets; los demás se ven marcados.
C. `manual` compite y `site_assisted` queda como referencia hasta que cada cadena tenga veredicto de términos (S5). Consecuencia: va con la lectura asistida apagada por defecto.
Impacto: S5 (mejor tienda por ingrediente y por canasta). No cambia el esquema ni cifras del oráculo.
Cómo seguí: RN-12 queda FIRME para las tres fuentes que decidió K22; para `site_assisted` y `manual` lleva la marca PROVISIONAL citando este bloqueo. Hoy `site_assisted` no produce precios (interruptor apagado).

## Respuesta de Axel (puerta 1b, 11-oct-2026)
- **Decisión:** opción **C**: `manual` compite con el mismo peso que el ticket y la misma regla de frescura; `site_assisted` queda solo como referencia (como PROFECO) hasta que cada cadena tenga veredicto de términos (S5).
- **Estado:** RESUELTO (K37). Se aplica en: RN-12 1.3, criterios S5 1.2, S5 (WF-1).
