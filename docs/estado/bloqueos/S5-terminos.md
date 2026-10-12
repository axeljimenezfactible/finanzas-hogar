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
