# Bloqueo S1-ajuste-objetivo-kcal
Contexto: al aplicar K23 (meta de calorías sugerida con Mifflin-St Jeor × factor de actividad, con ajuste por objetivo) escribí RN-17 en `docs/contratos/reglas-negocio.md` 1.2 y la columna `health_profile.goal` (`lose` / `maintain` / `gain`) en la migración `00002_puerta1.sql`. La fórmula y los 5 factores (1.2, 1.375, 1.55, 1.725 y 1.9) salen de la opción B del bloqueo `S1-meta-kcal`, incluida en la D que Axel eligió. No está decidido cuánto cambia la meta con «bajar» o «subir»: la opción D solo daba «por ejemplo −500 kcal para bajar».
Qué falta: ¿Cuánto se ajusta la meta sugerida cuando el objetivo es bajar o subir de peso?
Opciones:
A. Ajuste fijo simétrico: −500 kcal para bajar y +500 kcal para subir (el ejemplo de la opción D). Consecuencia: simple de explicar; en personas con gasto bajo la meta puede quedar muy abajo (el esquema solo acepta 500 kcal o más).
B. Ajuste fijo asimétrico: −500 kcal para bajar y +300 kcal para subir. Consecuencia: se sube de peso más despacio; son dos números que explicar.
C. Ajuste en porcentaje: −15 % para bajar y +10 % para subir. Consecuencia: escala con el gasto de cada persona; el redondeo a kcal enteras va al final (HALF_UP).
D. Cada persona escribe su ajuste en kcal en la pantalla Salud. Consecuencia: hace falta otra columna (migración nueva) y la sugerencia deja de ser automática.
Impacto: S4 (meta sugerida) y la pantalla Salud. «Mantener» no ajusta. Ninguna cifra del oráculo depende de esto.
Cómo seguí: RN-17 queda FIRME en la fórmula, los factores, el redondeo y «mantener»; la magnitud del ajuste de «bajar» y «subir» va marcada PROVISIONAL citando este bloqueo. El esquema 1.2 no cambia con las opciones A, B o C.
