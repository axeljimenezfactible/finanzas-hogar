# Bloqueo U5-contraste-ui-controles

Contexto: WF-3b, agente `u5-capturas` (12-oct-2026). U5-04 pide «texto ≥ 4.5:1, texto grande y componentes de UI ≥ 3:1 (WCAG 2.x)» y **0 pares bajo umbral**, sin decir qué cuenta como «componente de UI» ni contra qué se mide su límite.

Qué falta: la regla exacta para los controles que llevan texto. El script mide, por cada control (`button`, `[role=button|tab|checkbox|radio|switch|combobox]`, `input`, `select`, `textarea`) con límite visible (borde o relleno distinto de lo que lo rodea), la mejor razón entre su borde o su relleno y el fondo de al lado. Con esa lectura literal salen pares bajo 3:1 en tres familias: (1) **campo**: su límite es lo que lo identifica y WCAG 1.4.11 sí lo exige; (2) **seleccion**: pestañas de vista (Tabla, Galería) y chips cuyo estado se marca solo con relleno claro (`#F0F1FF` sobre `#FAFAFB`, 1.075:1); (3) **boton**: tarjetas-botón y chips con borde de separación (`#E5E6E9` sobre `#FAFAFB`, 1.196:1). En (2) y (3) el texto del propio control ya lo identifica (WCAG 1.4.11 no exige el límite cuando otra pista lo identifica), pero el estado seleccionado de (2) solo se distingue por el relleno.

Opciones:
- **A (la que apliqué, PROVISIONAL)**: contar los tres subtipos, informar el total literal y el subtotal por subtipo en el reporte y en cada JSON (`subtipo`). No oculto nada.
- B: contar solo `campo` (y `seleccion` por el estado) y dejar los botones con texto fuera del 3:1; el texto y los íconos siguen en 4.5:1 y 3:1.
- C: contar solo `campo`.

Consecuencia: el veredicto de U5-04 («0 bajo umbral») cambia según la opción. Lo decide Axel con la revisión; el agente de corrección de DESIGN.md/C2 decide si sube el borde de los chips a `border-control` o lo deja.
