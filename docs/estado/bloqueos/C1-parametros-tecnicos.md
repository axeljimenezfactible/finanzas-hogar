# Bloqueo C1-parametros-tecnicos
Contexto: Al fijar las convenciones únicas de C1 (paginación, idempotencia, cookie, contraseña, cargas de archivo). Actualizado en la corrección r2 de C1 (11-oct-2026): el esquema 1.1 de S1 ya trae la tabla `idempotency_key`, así que el punto 2 solo pregunta la retención. Estos valores no están en el plan, las decisiones ni el oráculo; los criterios de C1 (C1-15) piden marcarlos PROVISIONAL.
Qué falta: aprobar o cambiar cinco valores técnicos.
1. Máximo de registros por página (`limit`): propuesto 200 (el valor por defecto de 50 sí lo fija U1).
2. `Idempotency-Key`: retención propuesta 24 horas. Las claves se guardan en la tabla `idempotency_key` del esquema 1.1 de S1 (cuenta, clave, método, ruta, huella SHA-256 del cuerpo, respuesta guardada, fecha; S1 C-15) y un job de S6 las purga por antigüedad.
3. Cookie de sesión `fh_session` con `SameSite=Strict` (la decisión A y C1-12 piden SameSite sin fijar el modo; Strict es el más cerrado y la PWA es del mismo sitio).
4. Longitud mínima de la contraseña nueva (N-273): propuesta 12 caracteres.
5. Tamaño máximo del archivo de PROFECO (RAR o CSV) en la importación (N-293): propuesto 1 GB; los adjuntos de gastos van aparte, con 10 MB, bajo el bloqueo existente `criterios-umbrales-ausentes` punto 5.
Opciones:
- A. Aprobar los cinco valores propuestos. Consecuencia: C1 v1.0 queda sin estas marcas PROVISIONAL al pasar la puerta 1.
- B. Cambiar uno o varios valores. Consecuencia: C1 v1.1 con los valores nuevos (no cambian rutas ni esquemas, solo límites y descripciones).
Impacto: Parámetros `Limit` e `IdempotencyKey`, `securitySchemes.cookieAuth`, `PasswordChangeRequest` e `importProfeco`. Ninguna cifra canon depende de esto.
Cómo seguí: C1 v1.0 usa los valores de la opción A marcados `x-provisional` y PROVISIONAL en `C1-datos.md` §1, §2 y §5.
