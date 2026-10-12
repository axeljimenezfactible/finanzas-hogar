# Bloqueo WF3c-barometro-sin-dato · Rol accesible del barómetro cuando no hay dato

Abierto por `correccion-tablero-frontend-r2` (WF-3c, 12-oct-2026), a partir del hallazgo N-1 de `docs/estado/wf3/REVISION-tablero-frontend-r2.md` y de su PROPUESTA de §8. No cambia tokens, cifras ni textos visibles.

**Qué falta.** DESIGN §8.9 (Accesibilidad) pide `role="meter"` con `aria-valuenow` = valor, pero no dice qué hacer en el estado «Sin dato» que la misma §8.9 define (por ejemplo, sin ingreso; es el estado de la semilla de `datos/semilla/`). Un `meter` sin `aria-valuenow` es inválido: axe 4.13 da `aria-required-attr` (impacto crítico) y U5-03 falla.

**Hecho (PROVISIONAL).** `frontend/src/componentes/Barometro.tsx`: con valor, el `meter` de siempre (`aria-valuemin` 0, `aria-valuemax` = escala, `aria-valuenow` = valor, `aria-valuetext` con título, valor, tope y estado). Sin dato, el dibujo es `role="img"` con `aria-label` = el texto accesible que ya armaba la pantalla («Capacidad de endeudamiento sin dato, tope 40 %»), sin atributos `aria-value*`. El dígito «—» y el motivo con su acción siguen igual.

| Opción | Qué es | Impacto |
|---|---|---|
| A (hecha) | Sin dato: `role="img"` con el texto accesible | axe 0 serios en las 4 combinaciones; el lector dice el título, «sin dato» y el tope |
| B | Sin dato: contenedor sin rol y el mismo texto en `sr-only` | Mismo resultado; el texto se lee como párrafo y no como imagen |
| C | Mantener `meter` con `aria-valuenow` = 0 y `aria-valuetext` «sin dato» | Se descarta: anuncia un 0 % que no existe (cifra inventada) |

**Pregunta para quien edita DESIGN (o Axel).** ¿Se agrega a DESIGN §8.9, Accesibilidad, la línea «Sin dato: `role="img"` con el texto accesible, sin `aria-valuenow`» (opción A), o se prefiere B? Mientras no se decida, queda A.
