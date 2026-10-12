# artboards-shell · Pista 3, fase 2 (lienzo) · 12-oct-2026

## Hecho
8 artboards en `docs/design/maquetas/lienzo/project/shell/`, uno por escritura, más el fragmento `docs/design/maquetas/lienzo/fragmentos/shell.json`.

| Archivo | Tamaño de preview | Qué muestra |
|---|---|---|
| `shell-p00-entrar-m.dc.html` | 412 × 915 | P-00: Usuario («axel»), Contraseña con «Mostrar contraseña», franja de bloqueo «Demasiados intentos. Vuelve a intentar a las 08:10.» (5 intentos en 15 min), «Entrar» deshabilitado con el motivo y la sesión larga (90 días en este teléfono) |
| `shell-p00-entrar-e.dc.html` | 1440 × 1000 | P-00: formulario centrado de ancho fijo con «Entrar» (enlaza a Inicio) y la tarjeta «Estados»: Vacío (no aplica), Carga (botón en espera, campos bloqueados), Error («Usuario o contraseña incorrectos.»), Bloqueo, Sin servidor (franja G-07 y «Requiere conexión con el servidor de casa») y Sesión vencida |
| `shell-barra-lateral-m.dc.html` | 412 × 1854 | Dos pantallas apiladas: (1) panel ☰ al 85 % sobre Kueski Pay con fondo atenuado: encabezado y usuario, Buscar, Registrar gasto, las 9 secciones en orden fijo, Créditos desplegado con sus líneas en orden alfabético (Kueski Pay activa), Comparador e interruptor de tema Sistema · Claro · Oscuro; (2) pestaña «Buscar» con «/» tecleado: comandos agrupados (Registrar, Ir a, Herramientas) con su alias y la barra inferior G-05 con Buscar activo |
| `shell-barra-lateral-e.dc.html` (Main) | 1440 × 1160 | Barra lateral fija de 240 px: usuario, «Ocultar barra lateral (Ctrl+\)», Buscar (Ctrl+K), Registrar gasto (Alt+G), 9 secciones plegadas con su ▸ y el tema al pie. El contenido lleva el título «Finanzas Hogar», 9 tarjetas que enlazan a la primera página de cada sección y el menú «/» abierto (8 comandos, /gasto resaltado, teclas dibujadas) |
| `shell-p03-registrar-gasto-m.dc.html` | 412 × 915 | P-03 como hoja a pantalla completa: selector Gasto · Ingreso · Peso · Ticket, monto enfocado «134.50», 6 categorías sugeridas + «Más» (Restaurantes y café elegida), «Más opciones» plegado con el resumen de valores por defecto y «Guardar» de 48 px abajo. Es el camino de 3 toques de F-1 |
| `shell-p03-registrar-gasto-e.dc.html` | 1440 × 1200 | P-03 como diálogo de 640 px sobre Inicio atenuado, con «Más opciones» desplegado: descripción o comercio, fecha, persona, reparto, forma de pago (De contado y las 4 líneas), tienda, nota y adjunto. Lleva «Cancelar» y «Guardar» (Ctrl+Enter) y la tarjeta «Estados»: Vacío no aplica, Carga, Error junto al campo, Error al guardar, «Gasto guardado · Deshacer» y su variante sin conexión |
| `shell-p22-registrar-peso-m.dc.html` | 412 × 915 | P-22 en la misma hoja con el tipo «Peso» (F-6): persona Axel, campo de peso enfocado con «kg», «Antes: 84.6 kg el 19-oct-2026», fecha de hoy y «Guardar» |
| `shell-p22-registrar-peso-e.dc.html` | 1440 × 900 | P-22 como diálogo de 480 px sobre Salud: persona Ximena, peso, «Antes: 62.8 kg el 19-oct-2026», fecha y la tarjeta «Estados»: Vacío no aplica, Carga, Error junto al campo, «Peso guardado · Deshacer» y su variante sin conexión |

Todos los artboards siguen la misma forma. Llevan las líneas obligatorias, `lang="es-MX"`, `support.js` y `../ds/finanzas-hogar/tokens.css`. El `<helmet>` trae las reglas `body`, `a` y `a:hover` con tokens. Los colores, espaciados, radios, sombras, medidas y la tipografía salen de tokens (clases `.title-1/2/3 .body .small .caption`), sin `font-size` literal. Las cifras usan `tabular-nums`. Los iconos son Lucide en SVG inline, con `currentColor`, el trazo `--icon-stroke-user-units-*` y `aria-hidden`. Los controles son `button`, `a`, `input`, `select`, `textarea` y `label` reales. Los objetivos miden 44 px en móvil (48 px el primario de hoja). No hay barra de estado ni teclado dibujados. En móvil la raíz es fija de 412 px con border-box. En escritorio la raíz es fluida: fila con `flex-wrap`, barra `flex: 1 1 var(--layout-sidebar-width)` y contenido `flex: 999 1 560px`. Los enlaces del prototipo van a los artboards del mismo ancho de la lista del lienzo, o a «#».

## Verificación
`python3 -I /tmp/verif-shell/verif_shell.py /home/deepsec/finanzas-hogar` da **RESULTADO: OK** en 8 de 8 archivos. El script queda fuera del repositorio. Revisa lo siguiente:
- las líneas obligatorias (doctype, lang, charset, support.js, tokens.css, x-dc, helmet), el `<title>`, las reglas `body` y `a`/`a:hover`, y que exista `class Component extends DCLogic` con `renderVals()`;
- que `data-props` parsee como JSON y que `$preview` coincida con `w`/`h` del fragmento; en móvil 412 de ancho y al menos 915 de alto, con `width:412px` y border-box; en escritorio 1440 con `expand: fill` y sin ancho en px en la raíz;
- con `html.parser`, que no queden etiquetas sin cerrar ni cierres sin pareja;
- que no haya hex literales fuera de `data-props` ni `rgb()`/`hsl()`;
- que aparezcan 0 cifras canon: es la expresión `CANON` de `tools/ds/verificar_ds.py` (P3-12), con y sin coma;
- que cada `var(--x)` exista en `tokens.json` (los `name` y `font-<familia>`);
- que no haya `innerHTML`, iframes, URLs de red, emoji ni negativos con guion; que el dinero tenga el formato `$1,234.56` y `tabular-nums`; que no aparezca lorem ni placeholder;
- que todo `href` apunte a un artboard de la lista con el mismo ancho, o a «#»;
- que el fragmento liste exactamente los 8 archivos y que `order` coincida con `boards`.

Antes de correrlo probé el verificador con casos malos: detecta una etiqueta sin cerrar, `#fff`, `#22232566` y las cifras canon.

## Supuestos
1. **Categorías sugeridas (N-050).** El JSON no trae el conteo de 90 días, así que usé `categorias[].num_gastos` del mes y desempaté por el orden de P-26.04. Quedan: Frutas y verduras (mercado) 5, Servicios y suscripciones 4, Súper de la semana 3, Transporte 3, Restaurantes y café 3 y Despensa común 2.
2. **Hora de desbloqueo «08:10» (P-00).** No está en el JSON. La derivé de `ajustes.usuarios[axel].ultimo_acceso` (07:55) más `parametros.ventana_min` (15). Los 5 intentos y los 90 días de sesión salen de `ajustes.parametros`.
3. **Capturas de P-03.**
   - En móvil, Ximena captura `g28` (Café y pan, $134.50, Restaurantes y café, De contado, Hogar, 21-oct-2026).
   - En escritorio, Ximena captura `g29` (Bodega Aurrera · lácteos y reposición, $411.35, Súper de la semana, tienda Bodega Aurrera).
   - El adjunto se ve vacío aunque `g29.tiene_ticket` es verdadero: supuse que el ticket se liga después desde P-05.
4. **P-22.** El campo de peso está vacío y enfocado, porque no hay un peso nuevo en el JSON. El «Antes» sale de `salud.perfiles` (Axel 84.6 kg y Ximena 62.8 kg, los dos del 19-oct-2026).
5. **P-00 de escritorio** va sin barra lateral (U1 §1.1: queda fuera de la barra y aún no hay sesión).
6. **`shell-barra-lateral-m`** apila dos pantallas (412 × 1854). En móvil el menú «/» es la pestaña Buscar y no puede convivir con el panel ☰. La página de fondo es Kueski Pay para que Créditos se vea desplegado (U1 §1.2 punto 4: solo se despliega la sección abierta).
7. **Iconos de sección.** Usé `receipt` para Gastos, porque DESIGN §5.4 no fija uno para esa sección. Los demás: `wallet` en Presupuesto, `credit-card` en Créditos, `heart-pulse` en Salud, `file-chart-column` en Reportes y `settings` en Ajustes. Las rutas Lucide están escritas de memoria y no las comparé con el sprite 1.55.0.
8. **Enlaces.**
   - Ajustes enlaza a `ajustes-p26-indice-*` y no a «01 Cuenta y roles», porque no conozco el slug de s01.
   - Kueski Pay enlaza a `creditos-p18-credito-*`: supuse que el artboard de P-18 es el de Kueski. Las otras tres líneas van a «#».
   - El tipo «Ingreso» de la hoja enlaza a P-08.
9. **Fondo bajo el velo** en P-03 y P-22 de escritorio: barra lateral simplificada sin iconos y encabezado de Inicio o Salud, marcados con `inert` y `aria-hidden`.
10. **Menú «/» de escritorio.** Lo dibujé en el flujo, bajo las tarjetas de Main, para no tapar los enlaces de entrada. En la app real flota al 20 % del alto.
11. **Valores fijos.** El tema queda en «Sistema» (`ajustes.tema`). Ningún artboard es interactivo (`is_interactive: false`): solo funcionan los enlaces.
12. **Alturas.** Las estimé sin renderizar, porque K45 no me deja usar navegador.

## Bloqueos
- Ninguno.

## Pendientes
- Orquestador: renderizar los 8 artboards para confirmar alturas (sobre todo 1854, 1160 y 1200) y el tema oscuro. Publicar `support.js`, `tokens.css` y `fonts/` en `ds/finanzas-hogar/`.
- Fijar el slug de `ajustes-p26-s01-*` y, si existe, enlazar Ajustes de la barra de escritorio a «01 Cuenta y roles».
- Revisar los trazos Lucide contra el sprite 1.55.0.

## Consumo aproximado
- Leídos: unos 130 KB (format.md completo, secciones de craft, U1, DESIGN §5.4/§7/§10.5, nombres de tokens, subconjuntos del JSON sintético, LEEME y README).
- Escritos: unos 115 KB (8 artboards con 105 KB, fragmento con 2.6 KB, este reporte y el verificador en /tmp).
- Acciones de herramienta: unas 30.
