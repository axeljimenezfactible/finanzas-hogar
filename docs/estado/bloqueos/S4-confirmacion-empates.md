# Bloqueo S4-confirmacion-empates
Contexto: corrección del contrato S4/S5 (WF-1, agente `correccion-contrato-S4S5-r1`, 11-oct-2026), hallazgo M5 de `docs/estado/wf1/REVISION-contrato-S4S5-r1.md`. El riesgo R9 del plan («empate asistido con confirmación») y S4-11 («los empates con USDA/OFF son asistidos y con confirmación humana; ninguno entra automático») exigen que una persona confirme cada empate producto ↔ fuente nutricional; S5 §11 pide lo mismo para las ligas de PROFECO (`datos/semilla/profeco/vinculos.csv`). La tabla `nutrition_fact` del esquema 1.3 no tiene columna de confirmación, así que la confirmación se registra en los CSV antes de cargar (columnas nuevas `estado`, `confirmo` y `fecha_confirmacion`, S4 §9.1 y S5 §11 de las interfaces 1.1). Ni §A–§K, ni el plan ni los criterios dicen **quién** confirma ni **cuándo** en WF-1.
Qué falta: ¿quién confirma los empates de `nutricion/tabla.csv` y las ligas de `profeco/vinculos.csv`, y en qué momento?
Opciones:
- **A)** Axel en la puerta 2, sobre una lista que entrega el agente de nutrición (y la del agente PROFECO) con los renglones `propuesto`: los que apruebe pasan a `confirmado` con `confirmo = axel` y la fecha. Consecuencia: WF-1 mide S4-08 con los renglones propuestos y confirmados (dos conteos) y nada propuesto entra a la base de la app.
- **B)** Axel o Ximena en la app después del despliegue (pantalla de Salud o de Tiendas y precios). Consecuencia: hace falta una columna de estado en `nutrition_fact` (esquema) y operaciones en C1; la app arranca sin kcal de lo no confirmado.
- **C)** Axel antes de que cierre WF-1, en una sesión aparte. Consecuencia: WF-1 depende de un paso humano a mitad del flujo.
- Un agente (por ejemplo el revisor) no puede confirmar: R9 y S4-11 piden confirmación humana.
Impacto: queda PROVISIONAL quién confirma y cuándo (S4 §9.1 y S5 §11). No cambia ninguna cifra canon; las fórmulas de kcal no dependen de la confirmación. El cargador ya queda fijado: en modo producción no carga un renglón `propuesto` (S4-11, «ninguno entra automático»); en modo revisión (pruebas, fixtures y medición de WF-1) sí lo usa y lo cuenta aparte.
Cómo seguí: opción A, marcada PROVISIONAL citando este bloqueo en `docs/contratos/motores/S4-interfaz.md` §9.1 y `docs/contratos/motores/S5-interfaz.md` §11 (versiones 1.1).

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-43 · ¿Quién confirma los empates producto ↔ dato nutricional y las ligas de PROFECO, y cuándo?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Axel en la puerta 2, sobre una lista que entrega el agente de nutrición (y la del agente PROFECO) con los renglones `propuesto`: los que apruebe pasan a `confirmado` con `confirmo = axel` y la fecha.»
- **Consecuencia:** WF-1 mide S4-08 con los renglones propuestos y confirmados (dos conteos) y nada propuesto entra a la base de la app.
- **Qué cambia:** `S4-interfaz.md` §9.1 y `S5-interfaz.md` §11 quitan la marca PROVISIONAL de quién y cuándo; Axel confirma la lista de renglones «propuesto» (la de nutrición es la Q2-90, sin respuesta); en producción no se carga ningún renglón «propuesto»; no toca esquema ni C1. Acción pendiente de Axel en PUERTA-2 §7.
- **Estado:** RESUELTO

**Estado del bloqueo (puerta 2): RESUELTO**
