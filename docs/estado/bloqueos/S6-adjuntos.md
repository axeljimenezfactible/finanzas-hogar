# Bloqueo S6-adjuntos · valores de adjuntos sin decidir (PROVISIONAL)

12-oct-2026 · agente auth-sistema (WF-3a) · K31 (puerta 4), K43, S6-10.

| # | Qué falta | Opciones | Valor PROVISIONAL y por qué | Impacto |
|---|---|---|---|---|
| 1 | Factor de bytes de «10 MB»: C1 no lo declara (S6-10 dice que lo fija C1) | a) 1 MB = 1,048,576 bytes; b) 1 MB = 1,000,000 bytes | a) `adjuntos.BytesPerMB = 1 << 20`: es la propuesta de S1 que cita S6-10 y lo que muestra un teléfono como «10 MB» | Con b) el límite baja 4.6 %; cambia una constante y la prueba `TestAdjunto_Limite10MBYLimiteCambiado`. Falta que C1 (generador) lo escriba en `uploadExpenseAttachment` |
| 2 | Tipos admitidos (K31) | a) JPG, PNG, WebP y PDF; b) agregar HEIC | a) los que ya declara C1 (U1 N-064); el tipo se detecta por el contenido, no por el nombre ni la cabecera | Otro tipo responde 415 |
| 3 | `kind` sin valor en la subida (opcional en C1, obligatorio en la base) | a) `other`; b) `ticket_photo`; c) obligatorio (422) | a) `other` (neutro) | La pantalla P-xx debería mandarlo siempre |
| 4 | Casos de 409 en `deleteAttachment` (C1 lo declara sin decir cuándo) | a) borrar y dejar `price`/`credit_statement.attachment_id` en NULL (ON DELETE SET NULL del esquema); b) 409 `in_use` si lo usan | a) lo que dice el esquema 1.4 | Con b) hay que decidir el mensaje |
