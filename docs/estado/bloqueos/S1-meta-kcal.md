# Bloqueo S1-meta-kcal
Contexto: Al modelar el perfil de salud (reglas-negocio.md RN-17). La decisión D pide un perfil con altura, peso con historial, edad, sexo y actividad, IMC y una «meta de calorías por día». Es la misma pregunta del bloqueo `criterios-meta-kcal` (escrito en paralelo): una sola respuesta resuelve los dos.
Qué falta: ¿La meta de calorías la captura cada persona, o el sistema la sugiere con una fórmula? Si la sugiere: ¿qué fórmula, qué escala de actividad y si se ajusta según el objetivo (bajar, mantener o subir de peso)?
Opciones:
A. Solo captura manual de la meta. Consecuencia: la actividad y el sexo quedan como datos informativos.
B. Sugerencia con Mifflin-St Jeor × factor de actividad (1.2, 1.375, 1.55, 1.725 y 1.9), y la persona confirma o cambia la meta. Consecuencia: la escala de 5 niveles del esquema queda firme.
C. Sugerencia con Harris-Benedict revisada × los mismos factores. Consecuencia: igual que B, con otra fórmula.
D. B más un ajuste según el objetivo (por ejemplo −500 kcal para bajar). Consecuencia: hace falta una columna de objetivo (migración nueva).
Impacto: S4 (meta y comparación de kcal por día) y la pantalla Salud. `health_profile.activity_level` queda PROVISIONAL.
Cómo seguí: El esquema guarda los insumos (estatura, fecha de nacimiento, sexo, actividad, historial de peso) y la meta capturada (`daily_kcal_target`). La escala de actividad de 5 niveles va marcada PROVISIONAL en la migración y en el diccionario.


## Respuesta de Axel (puerta 1, 11-oct-2026)
- **Decisión:** Opción D: meta sugerida con Mifflin-St Jeor × factor de actividad (5 niveles, 1.2 a 1.9) con ajuste por objetivo (bajar, mantener, subir de peso; columna nueva en la migración 00002); la persona confirma o escribe otra.
- **Estado:** RESUELTO. Se aplica en: S1 1.2 (migración 00002), S4, U4.
