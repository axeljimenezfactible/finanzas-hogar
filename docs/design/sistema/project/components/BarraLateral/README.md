Barra lateral tipo Notion con el espacio «Finanzas Hogar», dos acciones fijas, las 9 secciones en orden fijo, subpáginas plegables y el tema al pie.

- **Cuándo:** navegación principal desde 1024 (fija y colapsable), superpuesta a 768 y como panel desde el botón `menu` en móvil.
- **Provee quien lo usa:** la sección y subpágina abiertas (`aria-current="page"`), el estado plegado de cada sección, las líneas de Créditos por su nombre y el usuario de la sesión.
- **Orden fijo:** Inicio, Gastos, Presupuesto, Comidas, Compras y tiendas, Créditos, Salud, Reportes, Ajustes. No se reordenan ni se ocultan.
- **Hacer:** fondo `var(--ui-surface-card)`, ancho `var(--layout-sidebar-width)`; elementos de `var(--size-sidebar-item)` con puntero fino y `var(--size-control-coarse)` con grueso; reposo `var(--ui-text-secondary)`; hover `var(--ui-surface-hover)`; activo `var(--ui-surface-selected)` con texto `var(--ui-text-default)` a 500 e icono `var(--ui-icon-selected)`; el chevron es un botón propio con `aria-expanded` que gira 90° en `var(--duration-base)`; subpáginas con sangría `var(--space-4)`.
- **No hacer:** abrirla deslizando desde el borde en móvil (choca con «atrás»); deshabilitar secciones.
- **Supuesto provisional (P3-huecos-componentes, punto 1):** iconos de Gastos (`receipt`), Presupuesto (`wallet`), Salud (`heart-pulse`) y Reportes (`file-chart-column`); DESIGN §5.4 solo fija Inicio (`house`), Comidas (`utensils`), Compras (`shopping-cart`), Créditos (`credit-card`) y Ajustes (`settings`).
- **Tokens:** `var(--ui-surface-card)`, `var(--ui-surface-hover)`, `var(--ui-surface-selected)`, `var(--ui-text-secondary)`, `var(--ui-text-default)`, `var(--ui-icon-default)`, `var(--ui-icon-selected)`, `var(--ui-border-selected)`, `var(--ui-border-divider)`, `var(--ui-person-axel)`.
