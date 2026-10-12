# Bloqueo S2-proyeccion-interruptor
Contexto: interfaz del motor financiero (WF-1, `docs/contratos/motores/S2-interfaz.md` §2.10 y §3.12). S2-09 pide que, con el interruptor de Ajustes encendido (`projection.repeat_cadence_enabled`, K21), los meses 7–12 sumen «la cadencia bimestral de despensa a crédito y el papel alterno de $199», con monto y fechas que declare `reglas-negocio.md` con su fuente; si no se pueden fijar con los documentos, se abre este bloqueo. RN-28 lo da solo como ejemplo ilustrativo («se suma una despensa a crédito cada dos meses después del 10-abr-2027 y el papel de 199.00 en los meses alternos, como en el oráculo»).
Qué falta: (1) monto de cada despensa proyectada: ¿capital de la despensa base 4,000.00 (despensas 2–4) con la tasa vigente de la tabla por plazo, o el de la última despensa comprada, o una partida planificada que Axel captura? (2) Fecha: ¿día 10 cada dos meses desde la última despensa (10-jun-2027, 10-ago-2027) o el día de `start_date` de una partida planificada? (3) Reparto: ¿la parte fija de Karen del oráculo (1,191.50 de capital y 1,499.72 con intereses) o la regla de la categoría Despensa? (4) Papel alterno: ¿199.00 de contado en los meses sin despensa (may, jul y sep de 2027), pagado por quién (en el oráculo lo paga Karen, `papel[2]` y `papel[4]`)?
Opciones: A) Repetir el patrón del oráculo: despensa de 4,000.00 de capital el día 10 cada dos meses desde la última, mismas partes que las despensas 2–4, papel de contado de Karen el día 10 de los meses alternos. B) Que la cadencia salga de partidas planificadas que Axel captura (`planned_item` con `credit_id`, `interval_count` = 2) y S3 calcule sus series; sin partidas no hay cadencia. C) A como valor inicial sembrado y B para cambiarlo.
Impacto: con el interruptor encendido, los meses 7–12 y la fixture `s2-proyeccion-7-12-cadencia` quedan PROVISIONAL. El modo por omisión (apagado) no depende de esto.
Cómo seguí: la interfaz recibe la cadencia ya armada (`cadence`: compras proyectadas con su serie de S3 y gastos proyectados) y la usa solo en los meses 7–12 con el interruptor encendido; la salida lleva `provisional = "S2-proyeccion-interruptor"`. Ninguna cifra canon depende de esto.

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-21 · (1) Monto de cada despensa proyectada
- **Respuesta:** **C** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «C) A como valor inicial sembrado y B para cambiarlo.»
- Subpregunta literal del archivo: «(1) monto de cada despensa proyectada: ¿capital de la despensa base 4,000.00 (despensas 2–4) con la tasa vigente de la tabla por plazo, o el de la última despensa comprada, o una partida planificada que Axel captura?»
- A (literal): A) Repetir el patrón del oráculo: despensa de 4,000.00 de capital el día 10 cada dos meses desde la última, mismas partes que las despensas 2–4, papel de contado de Karen el día 10 de los meses alternos.
- B (literal): B) Que la cadencia salga de partidas planificadas que Axel captura (`planned_item` con `credit_id`, `interval_count` = 2) y S3 calcule sus series; sin partidas no hay cadencia.
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: con el interruptor encendido, los meses 7–12 y la fixture `s2-proyeccion-7-12-cadencia` quedan PROVISIONAL. El modo por omisión (apagado) no depende de esto.
- **Qué cambia:** Interruptor `projection.repeat_cadence_enabled` (K21): con él encendido, los meses 7 a 12 usan el patrón del oráculo como valor inicial sembrado en `app_setting` (K44) y cambian por partidas planificadas; `S2-interfaz.md` §2.10 y §3.12 y RN-28 (`reglas-negocio.md`) quitan la marca; fixture `s2-proyeccion-7-12-cadencia`. Apagado, que es el valor por omisión, no depende de esto.
- **Estado:** RESUELTO

### Q2-22 · (2) Fecha de cada despensa proyectada
- **Respuesta:** **C** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «C) A como valor inicial sembrado y B para cambiarlo.»
- Subpregunta literal del archivo: «(2) Fecha: ¿día 10 cada dos meses desde la última despensa (10-jun-2027, 10-ago-2027) o el día de `start_date` de una partida planificada?»
- A (literal): A) Repetir el patrón del oráculo: despensa de 4,000.00 de capital el día 10 cada dos meses desde la última, mismas partes que las despensas 2–4, papel de contado de Karen el día 10 de los meses alternos.
- B (literal): B) Que la cadencia salga de partidas planificadas que Axel captura (`planned_item` con `credit_id`, `interval_count` = 2) y S3 calcule sus series; sin partidas no hay cadencia.
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: con el interruptor encendido, los meses 7–12 y la fixture `s2-proyeccion-7-12-cadencia` quedan PROVISIONAL. El modo por omisión (apagado) no depende de esto.
- **Qué cambia:** Interruptor `projection.repeat_cadence_enabled` (K21): con él encendido, los meses 7 a 12 usan el patrón del oráculo como valor inicial sembrado en `app_setting` (K44) y cambian por partidas planificadas; `S2-interfaz.md` §2.10 y §3.12 y RN-28 (`reglas-negocio.md`) quitan la marca; fixture `s2-proyeccion-7-12-cadencia`. Apagado, que es el valor por omisión, no depende de esto.
- **Estado:** RESUELTO

### Q2-23 · (3) Reparto de cada despensa proyectada
- **Respuesta:** **C** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «C) A como valor inicial sembrado y B para cambiarlo.»
- Subpregunta literal del archivo: «(3) Reparto: ¿la parte fija de Karen del oráculo (1,191.50 de capital y 1,499.72 con intereses) o la regla de la categoría Despensa?»
- A (literal): A) Repetir el patrón del oráculo: despensa de 4,000.00 de capital el día 10 cada dos meses desde la última, mismas partes que las despensas 2–4, papel de contado de Karen el día 10 de los meses alternos.
- B (literal): B) Que la cadencia salga de partidas planificadas que Axel captura (`planned_item` con `credit_id`, `interval_count` = 2) y S3 calcule sus series; sin partidas no hay cadencia.
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: con el interruptor encendido, los meses 7–12 y la fixture `s2-proyeccion-7-12-cadencia` quedan PROVISIONAL. El modo por omisión (apagado) no depende de esto.
- **Qué cambia:** Interruptor `projection.repeat_cadence_enabled` (K21): con él encendido, los meses 7 a 12 usan el patrón del oráculo como valor inicial sembrado en `app_setting` (K44) y cambian por partidas planificadas; `S2-interfaz.md` §2.10 y §3.12 y RN-28 (`reglas-negocio.md`) quitan la marca; fixture `s2-proyeccion-7-12-cadencia`. Apagado, que es el valor por omisión, no depende de esto.
- **Estado:** RESUELTO

### Q2-24 · (4) Papel alterno
- **Respuesta:** **C** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «C) A como valor inicial sembrado y B para cambiarlo.»
- Subpregunta literal del archivo: «(4) Papel alterno: ¿199.00 de contado en los meses sin despensa (may, jul y sep de 2027), pagado por quién (en el oráculo lo paga Karen, `papel[2]` y `papel[4]`)?»
- A (literal): A) Repetir el patrón del oráculo: despensa de 4,000.00 de capital el día 10 cada dos meses desde la última, mismas partes que las despensas 2–4, papel de contado de Karen el día 10 de los meses alternos.
- B (literal): B) Que la cadencia salga de partidas planificadas que Axel captura (`planned_item` con `credit_id`, `interval_count` = 2) y S3 calcule sus series; sin partidas no hay cadencia.
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: con el interruptor encendido, los meses 7–12 y la fixture `s2-proyeccion-7-12-cadencia` quedan PROVISIONAL. El modo por omisión (apagado) no depende de esto.
- **Qué cambia:** Interruptor `projection.repeat_cadence_enabled` (K21): con él encendido, los meses 7 a 12 usan el patrón del oráculo como valor inicial sembrado en `app_setting` (K44) y cambian por partidas planificadas; `S2-interfaz.md` §2.10 y §3.12 y RN-28 (`reglas-negocio.md`) quitan la marca; fixture `s2-proyeccion-7-12-cadencia`. Apagado, que es el valor por omisión, no depende de esto.
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
