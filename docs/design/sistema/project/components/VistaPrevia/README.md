Vista previa de un registro: encabezado con «‹» o `x`, título editable en el lugar, propiedades y cuerpo, sin salir de la página.

- **Cuándo:** al tocar un registro de una tabla, tablero o calendario (G-06): página completa en menos de 768, panel central a 768 y panel lateral derecho desde 1024 (más ancho a 1440).
- **Provee quien lo usa:** el título, las propiedades en orden (4 y «Mostrar N más» en móvil), el cuerpo y las acciones del registro («Marcar pago» en un pago).
- **Hacer:** fondo `var(--ui-surface-raised)`, borde `var(--ui-border-divider)`, `var(--radius-md)`, `var(--shadow-1)`; ancho `var(--layout-preview-width-desktop)` a 1024 y `var(--layout-preview-width-wide)` a 1440; edición en el lugar con hover de solo borde `var(--ui-border-control-hover)` (sobre lo que flota, como sobre una tarjeta, DESIGN §7.3); el gasto en `var(--ui-role-expense)` con `−` y `arrow-up-right`; `x` con `aria-label`.
- **No hacer:** mostrarla a página completa desde 768 (a 768 es panel central y desde 1024 panel lateral).
- **Supuesto provisional (P3-huecos-componentes, punto 8):** título en `title-2`, propiedades en `small` y rellenos `var(--space-4)`; registro y nota de ejemplo.
- **Tokens:** `var(--ui-surface-raised)`, `var(--ui-border-divider)`, `var(--radius-md)`, `var(--shadow-1)`, `var(--layout-preview-width-desktop)`, `var(--layout-preview-width-wide)`, `var(--ui-border-control-hover)`, `var(--ui-text-secondary)`, `var(--ui-role-expense)`, `var(--ui-icon-default)`.
