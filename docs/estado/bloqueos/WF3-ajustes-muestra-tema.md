# Bloqueo WF3-ajustes-muestra-tema · Muestra de tarjeta en claro y oscuro (P-26.10)

- Fase: WF-3a, agente `front-ajustes-b` (12-oct-2026). Estado: **PROVISIONAL** (no detiene el trabajo).
- **Qué falta (PREGUNTA):** U4 P-26.10 pide en 1440 px «una muestra de tarjeta en claro y oscuro» a la vez. Los tokens de C2 (`tokens.css`) solo cambian de tema en `:root` (`:root[data-theme="dark"]` y `:root:not([data-theme="light"])`), así que una tarjeta no puede pintarse en el otro tema sin colores fuera de tokens.
- **Lo que quedó:** una sola muestra de tarjeta en el tema vigente, con su rótulo («Así se ve en claro / en oscuro / Sistema»); el cambio de tema es su propia vista previa (U4: «la página cambia de tema al instante»).
- **Opciones:** (A) quedarse así; (B) que C2 publique los mismos valores también bajo `[data-theme]` en cualquier elemento (cambio en `docs/contratos/tokens.css`, dueño C2) para pintar las dos muestras; (C) quitar la muestra doble de U4.
- **Recomendación:** A (no toca C2 y cumple «el cambio es su propia vista previa»).
- **Impacto:** solo P-26.10 en escritorio.
