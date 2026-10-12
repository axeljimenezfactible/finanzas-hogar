Vista de tabla de base de datos: encabezado con icono de propiedad, filas de 14 px, dinero a la derecha y fila seleccionada en marca suave.

- **Cuándo:** listas de registros que se comparan o se editan en el lugar (gastos, productos, pagos).
- **Provee quien lo usa:** `<table>` con `<th scope="col">` y `aria-sort`, las filas, la casilla de selección con `aria-label`, el estado `aria-selected`, y cada cifra ya formateada con signo e icono de rol.
- **Hacer:** fondo `var(--ui-surface-page)`, divisores `var(--ui-border-divider)`, hover `var(--ui-surface-hover-subtle)`, seleccionada `var(--ui-surface-selected)`; alto de fila por puntero; relleno de celda `var(--space-2)`; persona y tienda como chip sobre `var(--ui-surface-raised)` y en columnas distintas; patrón ARIA *grid* cuando se edita.
- **Móvil:** a 360 filas de 2 líneas (título y cifra) y el resto en la vista previa; solo la tabla se desplaza en horizontal, nunca la página.
- **No hacer:** usar `var(--ui-surface-hover)` en filas (baja la cifra de rol de 4.5:1); poner color sin signo en cifras; mezclar personas y tiendas en una columna.
- **Tokens:** `var(--ui-surface-page)`, `var(--ui-border-divider)`, `var(--ui-surface-hover-subtle)`, `var(--ui-surface-selected)`, `var(--ui-text-secondary)`, `var(--ui-role-expense)`, `var(--ui-action-primary-bg)`, `var(--size-checkbox)`.
