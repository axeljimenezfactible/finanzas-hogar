Aviso breve flotante con icono de rol, una frase y, si aplica, «Deshacer» y cerrar.

- **Cuándo:** confirmar lo que se guardó solo («Guardado · Deshacer»), avisar una falla al guardar o que algo quedó en el teléfono sin conexión.
- **Provee quien lo usa:** el rol (éxito, error, advertencia, información), el texto, la acción opcional y qué revierte.
- **Hacer:** fondo `var(--ui-surface-raised)`, `var(--radius-md)`, `var(--shadow-2)` en claro y borde `var(--ui-border-divider)`; ancho máximo `var(--layout-toast-max-width)`; texto `small` en `var(--ui-text-default)`; acción como enlace subrayado en `var(--ui-text-link)`; uno a la vez; el de éxito dura 5 s y se pausa con hover o foco; el de error se queda hasta cerrarlo. Entra subiendo 8 px en `var(--duration-base)` con `var(--ease-emphasized)`.
- **Accesibilidad:** éxito e información en `role="status"`, error en `role="alert"`; no roba el foco.
- **No hacer:** exclamaciones ni emoji; avisos que se van solos si son de error.
- **Tokens:** `var(--ui-surface-raised)`, `var(--ui-border-divider)`, `var(--ui-role-success)`, `var(--ui-role-error)`, `var(--ui-role-info)`, `var(--ui-text-link)`, `var(--ui-text-link-hover)`, `var(--ui-icon-default)`.
