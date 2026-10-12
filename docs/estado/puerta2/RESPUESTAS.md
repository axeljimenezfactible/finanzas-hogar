# Puerta 2 · Respuestas de Axel (11-oct-2026, 22:40–22:55, en el chat de la sesión; herramienta de preguntas)

**Documento base:** `docs/estado/puerta2/PUERTA-2.md` (commit 1ccb468), 98 preguntas Q2-01 a Q2-98.

## Modo de respuesta (decisión de Axel)
«Lo recomendado + excepciones»: **se aprueban en bloque todas las preguntas que traen una opción recomendada** en la hoja de respuestas (§8 de PUERTA-2.md). Las preguntas **sin recomendación** quedan **sin respuesta en esta puerta: PROVISIONAL con la opción que ya está en uso**, y pasan a WF-3 (regla de PUERTA-2.md §0). Axel leerá el documento cuando pueda y mandará **excepciones por id**; hasta entonces no hay ninguna excepción.

## Respuestas explícitas (preguntadas una por una)
| Q2 | Pregunta | Respuesta de Axel | Nota |
|---|---|---|---|
| Q2-11 | S1-parametros-avisos · dónde viven `send_time`, `advance_days`, `weekday` (7) y `day_of_month` (1) | **B · todo en `app_setting`** | Primero eligió «C · Mixta» con una descripción del orquestador que no coincidía con la C literal del bloqueo; al aclararlo eligió B. Una migración siembra las claves con los valores de la decisión F y de Q2-72/Q2-73, las quita de `params_json` y su Down las repone. |
| Q2-64 | S5-terminos · veredicto de términos de walmart.com.mx, bodegaaurrera.com.mx y sams.com.mx | **A · un agente lee solo las páginas legales** | Autorización explícita de Axel: excepción a «nada de sitios de tiendas», limitada a las páginas legales (términos, condiciones, aviso de privacidad, políticas de uso) y, si hace falta para ubicarlas, la portada; sin buscar ni extraer precios. Resultado en `docs/estado/S5-terminos.md` con URL y fecha. |
| Q2-96 | DESIGN.md 1.0 y C2-tokens.json 1.0 como base del sistema de diseño de la Pista 3 | **A · aprobar tal cual** | La Pista 3 por Artifacts arranca de inmediato (K45). |
| (modo) | Cómo resolver las 98 | **Lo recomendado + excepciones** | Ver arriba. |

## Cómo se aplica (P2c)
- Q2-01 a Q2-10: autorizadas las correcciones recomendadas (A; Q2-10 B). Q2-07: variante sin recomendación → la regla de respaldo que el motor ya aplica se documenta como PROVISIONAL (S4-presentacion-respaldo).
- Q2-12 a Q2-71 (bloqueos): la opción recomendada es la respuesta; en cada archivo de bloqueo se escribe «Respuesta de Axel (puerta 2)». Los «sin recomendación» (Q2-48, Q2-50 a Q2-52, Q2-56 a Q2-58) quedan abiertos con nota de PROVISIONAL «en uso A».
- Q2-72 a Q2-77 (valores K31/K44): aprobados → pasan de PROVISIONAL a FIRME en criterios, bloqueos y esquema (claves de `app_setting` con sus valores iniciales).
- Q2-78 a Q2-94 (ratificaciones): las que tienen recomendación quedan ratificadas; Q2-85, Q2-87 a Q2-91 y Q2-94 sin recomendación → siguen como están, PROVISIONAL.
- Q2-95 A (WF-3 no espera a la Pista 3), Q2-97 A (ronda P2c y commit «[Puerta 2]»), Q2-98 A (nombres reales en las maquetas; cifras sintéticas).

## Puerta 2 bis · Respuestas de Axel (11-oct-2026, 22:55, en el chat; herramienta de preguntas)
| Tema | Respuesta de Axel (literal) | Cómo se aplica |
|---|---|---|
| Q2-64 · términos de las tiendas (las tres portadas mostraron desafío anti-robot; el orquestador no lo salta) | «Ignoralo, todo el tema de la investigacion de los precios lo haremos aparte, a que me refiero con esto, crearemos el sistema para que tenga yo una sesion de claude code que pueda mediante claude chrome extension sacar los precios y los productos, si bien estos datos los podemos insertar manualmente esto tambien impacta a las recetas tambien pueden ser agregadas de forma autonoma por eso la regla no funciona ya que usaremos a un agente externo» | **K53 (nueva decisión, se registra en P2d):** la aplicación NO lee sitios de tiendas (se retira la «lectura asistida» del plan: interruptor, veredicto de términos y S5-07 dejan de aplicar). Un **agente externo** (una sesión de Claude Code de Axel con la extensión Claude in Chrome) obtiene precios y productos, y también puede agregar recetas, y los mete en la aplicación por su API. Consecuencias que se escriben como PROVISIONAL hasta la puerta 4: (1) C1 expone importación para el agente externo (productos, presentaciones, precios con origen y fecha; recetas) con una credencial de servicio propia; (2) el peso de esos precios en la comparativa (¿compiten como `manual` o solo referencia como PROFECO?) es PREGUNTA de la puerta 4; (3) `docs/estado/S5-terminos.md` queda como registro histórico («no aplica desde K53»); el bloqueo S5-terminos se cierra con esta respuesta. |
| Excepciones a las 70 recomendadas | «Ninguna» | Se aplican tal cual. |
| Las 28 sin recomendación | «Dejarlas PROVISIONAL» | Siguen con la opción en uso; se preguntan en la puerta 4. |
| Encadenar de fondo esta noche | «Sí, encadena hasta WF-3» | P2c → P2d (K53) → Pista 3 fases 1 y 2 → WF-3 sin detenerse en la puerta 3 (Axel revisa la URL del lienzo mañana; cambios como correcciones del frontend); el orquestador se detiene solo en la puerta 4. |
