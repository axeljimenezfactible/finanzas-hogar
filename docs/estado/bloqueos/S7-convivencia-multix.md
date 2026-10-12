# Bloqueo S7-convivencia-multix · cómo detectar el chat de Ximena sin chocar con MultiX

Abierto por s7-avisos (WF-3a) el 12-oct-2026. Criterio: S7-09 y K4. Estado: **ABIERTO (PREGUNTA a Axel)**.

## HECHO
- K4: el bot detecta el chat cuando Ximena le escribe y ella lo confirma en Ajustes. El servicio ya tiene la mitad de la app: `telegram.RegistrarActualizacion` recibe un objeto `Update` oficial de Telegram y deja el chat como candidato pendiente (solo `chat_id` y nombre); `listTelegramChatCandidates`, `confirmTelegramChatCandidate` y `dismissTelegramChatCandidate` funcionan (pruebas `TestDestinatarios` y `TestOp_*`).
- La API de bots entrega cada actualización a **un solo** consumidor: con `getUpdates` el segundo lector recibe 409 «terminated by other getUpdates request» y las actualizaciones confirmadas (`offset`) desaparecen para el otro; con un webhook activo `getUpdates` deja de funcionar (core.telegram.org/bots/api#getupdates). MultiX ya consume las actualizaciones del bot.
- No se llamó `getUpdates` contra el bot real ni se leyó `openclaw.json`.

## Falta (qué decide Axel)
Quién le pasa al servicio el mensaje de Ximena. Opciones:
- **A. MultiX reenvía (recomendada).** Cuando MultiX recibe un mensaje de un chat que no conoce (o que empieza con `/finanzas`), escribe el `Update` crudo en una carpeta de entrada local (p. ej. `datos/entrada-telegram/*.json`, permisos 0600) o lo manda a un endpoint local del servicio; el servicio lo pasa por `RegistrarActualizacion`. Costo: un cambio pequeño en MultiX. Riesgo: bajo; nadie más lee el bot.
- **B. Código de enlace.** Ajustes muestra un código de 6 dígitos con vigencia corta; Ximena lo manda al bot y MultiX lo reenvía como en A. Igual que A más una pantalla; evita confirmar un chat equivocado. Requiere C1 nuevo (no hay operación para el código).
- **C. Bitácora de MultiX.** El servicio lee el registro de MultiX. Depende del formato interno de MultiX (no documentado aquí); frágil.
- **D. Segundo bot solo para Finanzas.** Limpio técnicamente pero contradice K30 (bot MultiX) y agrega otra credencial.

## Impacto
Mientras no se decida, nadie entra a la lista de candidatos de forma automática: Ximena queda «pendiente de acción» en S7-07 y el chat de Axel se registra a mano con su OK (ver `ops/telegram-disparo.md`). Todo lo demás de S7-09 (destinatarios configurables, Karen rechazada, omitir y reportar sin chat, confirmación) ya está hecho.
