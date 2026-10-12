# Bloqueo S5-terminos

Abierto por `correccion-S5-r1` (WF-1, 11-oct-2026) a partir del hallazgo H-03 de `docs/estado/wf1/REVISION-S5-r1.md`. Se decide en la puerta 2.

Contexto: el criterio S5-07 (`docs/criterios/S5.md`) y el plan §3.4 piden el entregable `docs/estado/S5-terminos.md` con un veredicto escrito por sitio (Walmart, Bodega Aurrera y Sam's Club), cada uno con la URL de sus términos y la fecha de lectura, construido leyendo solo las páginas legales y sin extraer precios. La verificación es `grep -c "Veredicto" docs/estado/S5-terminos.md` = 3. Ese archivo no existe y ningún agente de WF-1 puede escribirlo: las reglas de WF-1 y el protocolo §7 prohíben entrar a sitios de tiendas («nada de sitios de tiendas»; Internet solo para PROFECO, USDA, Open Food Facts y documentación oficial de herramientas). El código ya cumple la otra mitad de S5-07: la lectura asistida está apagada por defecto y `go test ./internal/precios/asistida/ -run LecturaAsistida` comprueba que el interruptor se rechaza sin veredicto.

Qué falta: Axel, ¿cómo se obtiene el veredicto de términos de los 3 sitios (walmart.com.mx, bodegaaurrera.com.mx y sams.com.mx)?

Opciones:
- (A) Axel autoriza por escrito que un agente lea **solo** las páginas legales (términos y condiciones, aviso de uso) de los 3 sitios, sin buscar ni extraer precios, y escriba `docs/estado/S5-terminos.md` con veredicto, URL y fecha por sitio. Consecuencia: S5-07 se puede cerrar en WF-2; es una excepción explícita a la regla «nada de sitios de tiendas».
- (B) Axel lee los términos y dicta el veredicto (permitido / no permitido / dudoso) con la URL y la fecha; un agente lo transcribe. Consecuencia: ningún agente entra a los sitios; depende del tiempo de Axel.
- (C) S5-07 queda pendiente hasta la puerta 2 (o después) con la lectura asistida apagada; el criterio se reporta «NO CUMPLE: pendiente del bloqueo S5-terminos». Consecuencia: `site_assisted` sigue solo como referencia (K37) y nada cambia en cifras ni en esquema.
- (D) Axel descarta la lectura asistida en esta etapa: se ajusta S5-07 para pedir solo el interruptor apagado y el veredicto pasa a un flujo posterior. Consecuencia: hay que editar el criterio y el plan §3.4 (lo hace el integrador).

Impacto: S5-07 **no se marca cumplido**. La lectura asistida sigue desactivada y no activable (no hay veredicto). `site_assisted` no compite (K37). Nada más de S5 depende de esto.

Cómo seguí: no escribí `docs/estado/S5-terminos.md` (un archivo con la palabra «Veredicto» sin veredicto real podría engañar al interruptor o a la verificación) ni entré a ningún sitio de tienda. Dejé el interruptor como estaba (apagado, rechaza sin veredicto). Nada provisional entró en cifras ni en esquema.

## Resultado
11-oct-2026, 22:28 CST. Q2-64 = A (Axel). El agente S5-terminos intentó la portada de los 3 sitios (3 solicitudes de 12); los 3 devolvieron un desafío anti-robot («Verifica tu identidad») y no fue posible leer ninguna página legal. Se escribió `docs/estado/S5-terminos.md` con 3 veredictos «no verificado» y la evidencia del bloqueo. **El bloqueo sigue ABIERTO**: S5-07 no se cumple en el fondo. Nadie reintenta leer esos sitios. PREGUNTA para Axel, qué sigue: **B** (tú lees los términos y dictas el veredicto con URL y fecha), **C** (queda pendiente, lectura asistida apagada y S5-07 en «NO CUMPLE») o **D** (descartas la lectura asistida en esta etapa). El agente además propuso intentar directo las URL legales si Axel las da; eso sería una autorización nueva y no está dada. La lectura asistida sigue apagada y `site_assisted` solo como referencia (K37).

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-64 · ¿Cómo se obtiene el veredicto de términos de walmart.com.mx, bodegaaurrera.com.mx y sams.com.mx?
- **Respuesta:** **A** (respuesta explícita de Axel)
- **Texto de la opción (literal del archivo):** «Axel autoriza por escrito que un agente lea **solo** las páginas legales (términos y condiciones, aviso de uso) de los 3 sitios, sin buscar ni extraer precios, y escriba `docs/estado/S5-terminos.md` con veredicto, URL y fecha por sitio.»
- **Autorización (tal como la describe RESPUESTAS.md):** excepción a «nada de sitios de tiendas», limitada a las páginas legales (términos, condiciones, aviso de privacidad, políticas de uso) y, si hace falta para ubicarlas, la portada; sin buscar ni extraer precios. Resultado en `docs/estado/S5-terminos.md` con URL y fecha.
- **Consecuencia:** S5-07 se puede cerrar en WF-2; es una excepción explícita a la regla «nada de sitios de tiendas».
- **Qué cambia:** `docs/estado/S5-terminos.md` con URL y fecha por sitio (escrito el 11-oct-2026 por el agente S5-terminos de P2c: los 3 sitios devolvieron un desafío anti-robot y quedaron «no verificado»). Autorización de Axel tal como la describe `RESPUESTAS.md`: «excepción a "nada de sitios de tiendas", limitada a las páginas legales (términos, condiciones, aviso de privacidad, políticas de uso) y, si hace falta para ubicarlas, la portada; sin buscar ni extraer precios». Solo el agente S5-terminos usa WebFetch y solo en esas páginas.
- **Estado:** RESUELTO (la decisión de Axel; el bloqueo sigue ABIERTO, ver abajo)

**Estado del bloqueo (puerta 2): ABIERTO** — la decisión (A) está tomada, pero el intento de lectura del 11-oct-2026 no obtuvo veredicto (ver «Resultado» arriba); falta que Axel decida el siguiente paso (B, URL legales, C o D) y S5-07 sigue sin cumplirse.

## Respuesta de Axel (puerta 2 bis, 11-oct-2026, 22:55)
Fuente: `docs/estado/puerta2/RESPUESTAS.md`, sección «Puerta 2 bis» (en el chat, con la herramienta de preguntas). Se registró como K53 en P2d (`docs/plan/decisiones.md` §K, agente `k53-aplicar`, 12-oct-2026).
- **Palabras de Axel (literal):** «Ignoralo, todo el tema de la investigacion de los precios lo haremos aparte, a que me refiero con esto, crearemos el sistema para que tenga yo una sesion de claude code que pueda mediante claude chrome extension sacar los precios y los productos, si bien estos datos los podemos insertar manualmente esto tambien impacta a las recetas tambien pueden ser agregadas de forma autonoma por eso la regla no funciona ya que usaremos a un agente externo»
- **Decisión (K53):** la aplicación NO lee sitios de tiendas. Se retira la lectura asistida, con su interruptor y el veredicto de términos. Un agente externo (sesión de Claude Code de Axel con Claude in Chrome) obtiene precios y productos, puede agregar recetas y los mete por la API con una credencial de servicio propia.
- **Consecuencia:** ya no hace falta veredicto de términos de walmart.com.mx, bodegaaurrera.com.mx ni sams.com.mx. Las opciones B, C y D de este archivo quedan sin objeto. `docs/estado/S5-terminos.md` se conserva como registro histórico («no aplica desde K53»). S5-07 conserva id y severidad, con texto nuevo (`docs/criterios/S5.md`). K37 (en lo de `site_assisted`) y K48 quedan sustituidas en esa parte.
- **Lo que Axel no decidió** (peso de esos precios en la comparativa, qué puede escribir el agente, recetas en borrador o directas) se pregunta en la puerta 4: bloqueo nuevo `S5-agente-externo`.

**Estado del bloqueo (puerta 2 bis): RESUELTO** por K53.
