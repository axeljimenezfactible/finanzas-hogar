Menú «/» y búsqueda en un solo panel flotante: campo arriba, comandos agrupados, una fila resaltada a la vez, alias y atajo.

- **Cuándo:** «/» fuera de campos, `Ctrl+K` o `Ctrl+P` en escritorio; pestaña «Buscar» en móvil (hoja a pantalla completa).
- **Provee quien lo usa:** los grupos («En esta pantalla», «Registrar», «Ir a», «Vista», «Herramientas»), cada comando con icono, nombre, alias y atajo, y cuáles requieren conexión.
- **Hacer:** fondo `var(--ui-surface-raised)`, borde `var(--ui-border-divider)`, `var(--radius-lg)`, `var(--shadow-2)`, ancho `var(--layout-command-width)`; filas de alto de control; resaltada en `var(--ui-surface-hover)` con `aria-activedescendant` (el foco se queda en el campo); atajo en tecla dibujada con `var(--radius-sm)`; patrón ARIA *combobox* con *listbox*; búsqueda sin acentos ni mayúsculas.
- **No hacer:** guardar nada desde el menú; ocultar comandos sin conexión (muéstralos en `var(--ui-text-disabled)` con «Requiere conexión»).
- **Vacío:** «Sin resultados para "texto"» con la sugerencia de /ir.
- **Supuesto provisional (P3-huecos-componentes, punto 5):** los comandos y alias de la vista previa (`/gasto`, `/barra`, `/tema`, «Respaldar ahora») son ejemplos, no la lista definitiva; `/peso`, `Ctrl+\` y `Ctrl+Shift+L` sí salen de U1 y DESIGN §9.6.
- **Tokens:** `var(--ui-surface-raised)`, `var(--ui-surface-hover)`, `var(--ui-border-divider)`, `var(--ui-text-default)`, `var(--ui-text-secondary)`, `var(--ui-text-disabled)`, `var(--ui-icon-default)`.
