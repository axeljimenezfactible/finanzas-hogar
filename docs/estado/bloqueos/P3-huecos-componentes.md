# P3 · Huecos de DESIGN.md resueltos como supuesto en los componentes del sistema de diseño

Fecha: 11-oct-2026 (agente ds-brandbook). Estado: ABIERTO, PROVISIONAL; no detiene la Pista 3. Cada punto usa solo tokens existentes; ningún color nuevo.

1. **Iconos de 4 secciones de la barra lateral.** DESIGN §5.4 fija `house` (Inicio), `utensils` (Comidas), `shopping-cart` (Compras), `settings` (Ajustes) y `credit-card` (crédito), pero no Gastos, Presupuesto, Salud ni Reportes. Supuesto: `receipt` (Gastos, el del ejemplo de §5.2), `wallet` (Presupuesto), `heart-pulse` (Salud) y `file-chart-column` (Reportes), tomados de las secciones homónimas de Ajustes en §5.4. Pregunta: ¿se confirman?
2. **Bloque «Antes → Después»** (U4 §0.2, U4-05). DESIGN no fija su aspecto. Supuesto: fondo `ui-surface-card`, `radius-sm`, relleno `space-2` × `space-3`, texto `small`, «antes» en `ui-text-secondary`, flecha `arrow-right` en `ui-icon-default`, `role="status"`.
3. **Perilla del interruptor.** DESIGN fija la pista (36 × 20) pero no la perilla. Supuesto: 16 px con 2 px de aire (20 − 2 × 2) y desplazamiento de 16 px.
4. **Círculo de «hoy» en el calendario y relleno de chips.** Supuesto: círculo de 24 px (`space-6`) y relleno horizontal del chip `space-2`.
5. **Comandos y alias del menú «/»** de la vista previa (`/gasto`, `/barra`, `/tema`, «Respaldar ahora»): ejemplos; `/peso`, `Ctrl+\` y `Ctrl+Shift+L` salen de U1 y DESIGN §9.6.
6. **Pesos 500/600/700 y el 2 px del anillo de Conjunto, del avatar y de los marcadores** van como número en el CSS: el formato del tipo no compila `font.weight.*` como variable y DESIGN no tiene token para esos 2 px.

**Actualización 11-oct-2026, noche (agente correccion-ds-r1, revisión ds-r1 M-01, M-02 y M-03).**
- Punto 6, corregido en parte: el formato del tipo sí compila una familia de valores simples. `tokens.json` trae ahora la familia `fontWeight` (`font-weight-regular` 400, `font-weight-medium` 500, `font-weight-semibold` 600, `font-weight-bold` 700, mismos nombres que `tokens.css`) y los previews usan `var(--font-weight-*)`. Sigue como número el 2 px del anillo de Conjunto, del avatar, de los marcadores y del indicador de pestaña (DESIGN no tiene token para él). Tamaños e interlíneas siguen solo dentro de los estilos de Tipografía (clases `.caption` … `.title-1`); el `tokens.css` que compile el tipo no tendrá `--font-size-*`, `--line-height-*`, `--font-letter-spacing` ni `--font-numeric-tabular`.
- Los README de BarraLateral (punto 1), InterruptorOpciones (puntos 2 y 3), CalendarioPagos y ChipPersonaTienda (punto 4) y MenuBarra (punto 5) llevan ahora la línea «Supuesto provisional (P3-huecos-componentes, punto N)».
7. **Barra superior, pestañas y barra inferior** (DESIGN §7.10, componentes `BarraSuperior` y `BarraInferior`). DESIGN no fija: la variante de «+ Registrar» (supuesto: secundario, para dejar una sola primaria), el alto y relleno de la pestaña (supuesto: alto de control y `space-3`), el token del indicador de 2 px (va como número), el dibujo de «‹» (supuesto: `chevron-left`; `chevron-right` sí está en §5.4 para migas), si las migas van subrayadas (supuesto: sí, regla de enlace de §7.2) y el nombre del «+» de la barra inferior (supuesto: sin nombre visible y nombre accesible «Registrar gasto», P-03). Pregunta: ¿se confirman?
8. **Hoja inferior, diálogo y vista previa** (DESIGN §7.12, componentes `HojaDialogo` y `VistaPrevia`). DESIGN no fija los rellenos (supuesto: `space-4` en hoja y vista previa, `space-6` en diálogo, `space-2` entre acciones) ni el estilo del título de la vista previa (supuesto: `title-2`, como hojas y diálogos). Pregunta: ¿se confirman?
9. **Panel de Ajustes** (DESIGN §7.13, componente `PanelAjustes`). DESIGN no fija el estilo del título de sección (supuesto: `title-3`) ni el color del texto «Guardado» (supuesto: `ui.role.success`, como su icono; los roles pintan texto e icono). Pregunta: ¿se confirman?


**Actualización 11-oct-2026, noche (agente correccion-ds-r2, revisión ds-r2 A-01 y M-01).**
- A-01 de la r2 corregido sin supuestos nuevos: los botones de solo icono de `BarraSuperior` («Buscar» y «⋯») usan `size-icon-md` (DESIGN §5.2 y §7.2).
- M-01 de la r2: se agregan `BarraVista` (§7.5.1) y `Galeria` (§7.5.5). Sus huecos van en los puntos 10 y 11. Lo que falta del calendario genérico (§7.5.4) va en el punto 12.
10. **Barra de la vista** (DESIGN §7.5.1, componente `BarraVista`). DESIGN no fija:
    - el icono de «Agrupar» (supuesto: sin icono; §5.4 no le asigna uno);
    - el color del texto de la vista elegida (supuesto: se queda en `ui.text.secondary`; §7.2 solo cambia el relleno a `ui.surface.selected` y el icono a `ui.icon.selected`);
    - el aspecto de la búsqueda por texto (supuesto: campo de §7.3 con `search` a la izquierda);
    - el orden y la disposición en escritorio (supuesto: selector a la izquierda, después Filtrar, Ordenar, Agrupar y la búsqueda; los chips de filtro activo, en una línea propia dentro de la barra);
    - el lugar de «Filtrar» en móvil (supuesto: a la izquierda de los chips, bajo el título).
    Quedan sin dibujar, porque DESIGN y U1 no los fijan: qué abren «Filtrar», «Ordenar» y «Agrupar» en escritorio (¿menú o panel?); el contenido de la hoja de filtros de móvil; si el selector de vista aparece en móvil y cómo; y qué barra lleva el tramo de 768 a 1023 (U1 §1.3 compara escritorio ≥ 1024 con móvil < 768). Pregunta: ¿se confirman los supuestos? ¿Cómo se resuelven estos cuatro puntos?
11. **Galería** (DESIGN §7.5.5, componente `Galeria`). DESIGN no fija:
    - el icono de la entidad de un ticket (supuesto: `scan-line`, el de «Ticket con OCR» en §5.4, en `size.icon.md` como el estado vacío de §7.14);
    - el radio de la portada (supuesto: `radius.sm`);
    - el estilo de «Foto no disponible» (supuesto: `caption` en `ui.text.secondary`, bajo el icono);
    - la separación entre tarjetas (supuesto: `space.4`);
    - cuál es la segunda propiedad de un ticket (supuesto: la fecha; U1 §2.5.a da «foto, comercio, total»);
    - el patrón de selección (supuesto: el de `Tablero`, `listbox` con `option` y `aria-selected`);
    - qué lleva la tarjeta a 360 en 2 columnas (supuesto: título, cifra y comercio).
    La vista previa no lleva foto: no hay una imagen sintética de ticket en el repositorio. Pregunta: ¿se confirman? ¿Se agrega una foto sintética de ticket para la maqueta?
12. **Calendario genérico** (DESIGN §7.5.4). `CalendarioPagos` (§7.8) lo representa en su base, como dice ahora su README. **Pendiente:** no hay vista previa propia de lo que solo tiene el genérico:
    - la navegación «‹ Hoy ›»;
    - «+N más»;
    - los días de otro mes en `ui.text.disabled`;
    - el «+» que aparece en hover;
    - el renglón genérico de título y cifra;
    - la agenda de 360 con encabezado fijo.
    Si se dibuja, el contraste de `ui.text.disabled` necesita una excepción declarada en el verificador. Pregunta: ¿basta con `CalendarioPagos` o se agrega un componente `Calendario` propio?
