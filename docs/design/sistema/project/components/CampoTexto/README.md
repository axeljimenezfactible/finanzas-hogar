Campo con etiqueta arriba, valor a 16 px, prefijo o sufijo opcional y ayuda o error debajo.

- **Cuándo:** toda captura de texto, dinero, porcentaje, número con unidad, fecha u hora. El monto de «Registrar gasto» usa `title-2`.
- **Provee quien lo usa:** `<label for>`, el `<input>` con `inputmode` (dinero: `decimal`), el prefijo (`$`) o sufijo (`%`, `kg`, `kcal`), la ayuda y el mensaje de error enlazados con `aria-describedby`, y `aria-invalid="true"` cuando falla.
- **Hacer:** fondo `var(--ui-surface-field)`, borde `var(--border-width-default)` en `var(--ui-border-control)`, `var(--radius-sm)`, alto de control por puntero, `var(--space-2)` entre etiqueta y control y `var(--space-4)` entre campos; cifras tabulares en dinero; el anillo de foco va fuera del borde.
- **Error:** borde `var(--ui-border-error)`, `circle-alert` y mensaje en `var(--ui-role-error)` que dice qué pasó y qué hacer («Falta capturar la tasa de esta tarjeta»); el valor capturado no se borra.
- **No hacer:** usar el marcador de posición como etiqueta; campos de menos de 16 px en móvil; borrar lo capturado al fallar.
- **Tokens:** `var(--ui-surface-field)`, `var(--ui-border-control)`, `var(--ui-border-control-hover)`, `var(--ui-border-focus)`, `var(--ui-border-error)`, `var(--ui-text-default)`, `var(--ui-text-secondary)`, `var(--ui-role-error)`.
