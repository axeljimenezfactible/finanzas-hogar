# P3 · Contraste de la marca categórica sobre el hover del chip seleccionable

Fecha: 11-oct-2026 (agente ds-brandbook, Pista 3 fase 1). Estado: ABIERTO, no detiene la Pista 3.

**Qué pasa.** DESIGN.md §7.4 pide, en el chip seleccionable o de filtro, hover con fondo `ui.surface.hover`, y mide la marca de forma (10 px) sobre `ui.surface.raised`. Calculado con los hex de `tokens.json` (fórmula WCAG 2.x), cuatro marcas quedan bajo 3:1 sobre `ui-surface-hover`:

| Marca | Claro | Oscuro |
|---|---|---|
| `ui-person-karen` | **2.83** | 3.23 |
| `ui-store-sams-club` | **2.74** | 7.03 |
| `ui-person-axel` | 10.69 | **2.14** |
| `ui-store-market` | 11.48 | **2.11** |

Sobre `ui-surface-raised` y `ui-surface-selected` las 8 marcas cumplen 3:1 en los dos temas (mínimo 3.05). El par no está en la lista de §10.5.

**Qué se hizo.** La vista previa `components/ChipPersonaTienda/preview.html` respeta DESIGN tal cual (hover en `ui-surface-hover`); el verificador `docs/design/sistema/verificar_brandbook.py` lo marca como excepción declarada. La marca va con `aria-hidden="true"` y el nombre en texto, así que no es la única portadora del significado, pero 1.4.11 sigue pidiendo 3:1 a la marca si se considera necesaria para entender.

**Pregunta para Axel / U2.** Opciones: (A) el hover del chip cambia solo el borde (`ui.border.control-hover`), como las tarjetas, y el fondo sigue en `ui.surface.raised`; (B) se deja como está y se declara la marca decorativa en §7.4; (C) otra. Recomendación del agente: A (no cambia ningún color, sigue el patrón de §7.5.3 y §7.6).

**Actualización 11-oct-2026, noche (agente correccion-ds-r1, revisión ds-r1 A-01).** El README publicado de `components/ChipPersonaTienda` declara ahora la excepción con los 4 valores y el nombre de este bloqueo. El color y la regla `button.chip:hover` no se tocaron: se cambian solo cuando Axel responda (con la opción A, hover solo con borde `ui-border-control-hover` sobre `ui-surface-raised`). Nota: `docs/design/sistema/verificar_brandbook.py` solo declara como excepción el par de Karen; no se tocó porque está fuera del alcance de esta corrección.

