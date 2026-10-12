Estados de una vista cuando aún no hay contenido normal: carga con esqueleto, vacío con su acción y franjas de error, falta capturar y sin conexión.

- **Cuándo:** toda vista o tarjeta que carga datos; ninguna pantalla queda en blanco.
- **Provee quien lo usa:** la forma final del contenido (para el esqueleto), el icono de la entidad, la frase de qué falta, la acción que lo llena y el texto de qué pasó y qué hacer.
- **Hacer:** esqueleto estático en `var(--ui-skeleton-on-card)` o `var(--ui-skeleton-on-page)`; vacío con icono de 24 px en `var(--ui-icon-default)`, frase en `var(--ui-text-secondary)` y el botón («Aún no hay gastos este mes · Registrar gasto»); franjas con `var(--radius-md)`, relleno `var(--space-3)`, icono del rol y texto `var(--ui-text-default)` sobre `var(--ui-role-error-surface)` o `var(--ui-role-warning-surface)`; error y sin conexión con `role="alert"` la primera vez.
- **No hacer:** cifras provisionales; spinner en lugar del esqueleto (solo en acciones de más de 400 ms); perder lo capturado.
- **Tokens:** `var(--ui-skeleton-on-card)`, `var(--ui-skeleton-on-page)`, `var(--ui-icon-default)`, `var(--ui-text-secondary)`, `var(--ui-role-error)`, `var(--ui-role-error-surface)`, `var(--ui-role-warning)`, `var(--ui-role-warning-surface)`, `var(--ui-text-link)`.
