Vista de tablero: columnas por grupo con nombre, conteo y suma, y tarjetas planas con título, propiedades y cifra.

- **Cuándo:** cuando el registro avanza por grupos o se reparte por categoría.
- **Provee quien lo usa:** los grupos en su orden, las tarjetas (título, hasta 3 propiedades, cifra con rol), la acción «Nuevo» por columna y la alternativa sin arrastre «Mover a…» en el menú «⋯».
- **Hacer:** columnas de `var(--layout-board-column-width)`; tarjeta en `var(--ui-surface-card)`, `var(--radius-md)`, relleno `var(--space-3)`, plana; hover solo con borde `var(--ui-border-control-hover)`; elegida en `var(--ui-surface-selected)`; en arrastre `var(--shadow-1)` y borde `var(--ui-border-divider)`; columna destino con borde punteado de 2 px en `var(--ui-border-selected)`.
- **Móvil:** una columna visible a 360 y 412; el tablero se desplaza en su propio contenedor.
- **No hacer:** cambiar el fondo en hover (baja el contraste de la cifra); hacer del arrastre la única forma de mover.
- **Tokens:** `var(--ui-surface-card)`, `var(--ui-surface-selected)`, `var(--ui-surface-raised)`, `var(--ui-border-control-hover)`, `var(--ui-border-divider)`, `var(--ui-role-expense)`, `var(--ui-store-walmart)`, `var(--ui-store-sams-club)`, `var(--ui-store-market)`.
