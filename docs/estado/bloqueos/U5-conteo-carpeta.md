# Bloqueo U5-conteo-carpeta

Contexto: WF-3b, agente `u5-capturas` (12-oct-2026). U5-01 dice que `ls /home/deepsec/finanzas-hogar/docs/design/capturas-u5/ | wc -l` debe ser 4 × T «más los archivos de evidencia en subcarpetas, que se cuentan aparte», pero U5-03 y U5-04 fijan las subcarpetas `axe/` y `contraste/` dentro de esa misma carpeta y la tarea del agente pide además `hosts.json`. `ls` cuenta las subcarpetas y los archivos sueltos como entradas.

Qué falta: el valor literal del comando. Con los 4T = 164 PNG, `ls | wc -l` da 164 + `axe` + `contraste` + `evidencia` + `hosts.json` = 168.

Opciones:
- **A (la que apliqué, PROVISIONAL)**: los PNG van en la raíz, la evidencia en `axe/`, `contraste/`, `evidencia/` y `hosts.json`; el conteo de PNG (`ls *.png | wc -l` y `frontend/pruebas-u5/verificar_capturas.py`) da 164 de 164.
- B: mover los PNG a una subcarpeta (p. ej. `png/`) para que el comando literal dé 4T + subcarpetas.
- C: reformular el comando de U5-01 como `ls capturas-u5/*.png | wc -l`.

Consecuencia: el revisor que ejecute el comando literal verá 168. Lo decide el orquestador o Axel.
