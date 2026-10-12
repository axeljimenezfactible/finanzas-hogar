Tarjeta que muestra una sola cifra con su etiqueta, un dato de comparación y, si hay plan, una barra de avance.

- **Cuándo:** en el Resumen (6 tarjetas como máximo, orden fijo) y donde una sola cifra contesta la pregunta; nunca una gráfica para una sola cifra.
- **Provee quien lo usa:** etiqueta sin dos puntos, cifra ya formateada (`$6,418.25`, `+$12,305.00`), dato secundario con periodo nombrado («9 % más que septiembre»), el valor y los umbrales de Ajustes para la barra, y el nombre accesible completo de la tarjeta-enlace.
- **Hacer:** fondo `var(--ui-surface-card)`, `var(--radius-md)`, plana; relleno `var(--space-4)` en móvil y `var(--space-6)` desde 768; cifra en `title-2` con `.money`; rol con signo, icono y color; barra con `role="meter"` y `aria-valuetext`. Desde el primer umbral (80 % al inicio) la barra pasa a `var(--ui-role-warning)` con `triangle-alert` y «82 % usado».
- **No hacer:** truncar la cifra (`$99,999.99` cabe completo a 360 px); dibujar el llenado más allá de la pista; fijar umbrales en el componente.
- **Estados:** carga con esqueleto `var(--ui-skeleton-on-card)`; «Sin datos este mes»; «Falta capturar · Capturar»; «No se pudo cargar · Reintentar» sin cifra.
- **Tokens:** `var(--ui-surface-card)`, `var(--ui-text-secondary)`, `var(--ui-text-default)`, `var(--ui-role-expense)`, `var(--ui-role-income)`, `var(--ui-role-warning)`, `var(--ui-viz-track)`, `var(--ui-viz-meter)`, `var(--size-meter-height)`, `var(--ui-border-control-hover)`, `var(--ui-border-selected)`.
