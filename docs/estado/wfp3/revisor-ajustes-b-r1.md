# revisor-ajustes-b-r1 · Pista 3, fase 2 (lienzo) · revisión de la vuelta 1

Fecha: 12-oct-2026. Revisor independiente, vuelta 1 de 3; no corrige. Informe completo: `docs/estado/wfp3/P3-iter-1-ajustes-b.md`.

## Hecho
- Revisé los 14 artboards de ajustes-b (`ajustes-p26-s08` a `s14`, `-m` y `-e`) y el fragmento `fragmentos/ajustes-b.json` contra P3-08, P3-12, U5-04/06/07, DESIGN.md, U1 §1.3/§2.1–§2.3, U4 §0 y P-26.08 a P-26.14, y `referencia-tipo/format.md`.
- **Veredicto: NO APROBADO.** 1 alta, 11 medias y 10 bajas.
  - Alta: H-01, la perilla del interruptor en `color-base-white` queda a 1.52:1 y 1.98:1 en oscuro.
  - Medias: H-02, regla `.in` perdida (campos sin token y de menos de 44 px); H-03, claves técnicas visibles; H-04, «Valor inicial» lejos del campo en Notificaciones; H-05, fecha nativa del navegador; H-06, sin barra inferior en móvil; H-07, barra lateral sin colapsar; H-08, muestra clara sin tema; H-09, sin tabla de vista previa del CSV; H-10, campos de U4 ausentes; H-11, días con decimal y «Reponer» sin nivel; H-12, sugerencia de kcal sin fecha de nacimiento.

## Verificación
Hice la verificación con scripts propios (`python3 -I`) en `datos/tmp/rev-ajb-r1/`, que ya borré.
- **Forma, `data-props`, `$preview` frente al fragmento, `html.parser`, nombres y stems:** 14/14 OK.
- **Tokens:** 0 hex y 0 `rgb()`; todos los `var()` existen en `tokens.json`.
- **Contraste en el DOM con los tokens resueltos en claro y oscuro:** 23 combinaciones de texto y unos 1,362 nodos por tema, con 0 bajo 4.5:1. Recalculé 10 pares de texto y 4 de UI; fallan 2 de UI, en oscuro.
- **P3-12:** 0 cifras canon (21 patrones, con y sin coma, y en centavos) y 0 credenciales. La intersección de fechas con el calendario real no deja ninguna fecha junto a un monto del reporte. Los 12 valores de muestra están en `datos-sinteticos.json`, y las cifras derivadas cuadran.
- **Accesibilidad:** etiquetas, nombres de botones y enlaces, y `aria-hidden` correctos. Fallan los objetivos táctiles de s11-m y s12-m. 0 emoji, 0 barra de estado falsa, 0 `innerHTML`, `iframe` o red.
- **Enlaces:** 193 `href`; 0 rotos y 0 a artboards de otro ancho.

## Supuestos
1. Tomé la clave `type.families.sans` de `tokens.json` como fuente de `--font-sans`, porque `tokens.css` la genera con ese nombre.
2. Los defectos que vienen de `datos-sinteticos.json` (H-11, H-12 y la terminación del chat) se los cuento a los artboards, porque se ven en ellos, y propongo escalarlos a correccion-datos.
3. Calculé el contraste con un motor de selectores propio (clases, etiquetas, atributos, descendiente e hijo; sin pseudoclases). Los esqueletos y el estado deshabilitado (`.dis`) son decorativos o exentos.
4. La fila «Estados» se exige solo en escritorio, como dice la tarea. En Apariencia, U4 declara que no aplica.
5. Escribí el informe en `docs/estado/wfp3/`, como pide la tarea; P3-07 nombra `docs/estado/`. La copia a `docs/design/salida/iter-1/` la hace el orquestador.

## Bloqueos
Ninguno.

## Pendientes
- La vuelta 2 la corrige el autor de ajustes-b. H-11 y H-12 piden también un cambio en `datos-sinteticos.json`.
- Revisión visual y nueva medición de las alturas cuando el orquestador publique, porque aquí no hubo navegador (K45).
- H-06 y H-07 afectan también a ajustes-a y a la barra del shell: conviene unificarlas.

## Consumo aproximado
- Leído: unos 230 KB (U4 parte b y §0, U1 §1.3 y §2.3, DESIGN §7.13, §7.14, §9.3 y §9.4, format.md, P3, U5, extractos de texto de los 14 artboards y consultas al JSON).
- Escrito: unos 30 KB (scripts temporales, ya borrados, más los dos informes).
- Unas 45 acciones.
- 0 acciones de navegador y 0 llamadas a Artifact.
