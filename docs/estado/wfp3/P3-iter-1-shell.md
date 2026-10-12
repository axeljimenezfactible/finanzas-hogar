# P3-iter-1 · área shell · revisión independiente (vuelta 1 de 3)

Revisor: `revisor-shell-r1` · 12-oct-2026 · Criterio P3-07 (`docs/criterios/P3-claude-design.md` 1.3). No corrige: solo evalúa.

**Veredicto: NO APROBADO** · 0 altos · **2 medios** · 7 bajos. Aprobar exige 0 altos y 0 medios.

**Qué se revisó:** los 8 artboards de `docs/design/maquetas/lienzo/project/shell/` y el fragmento `docs/design/maquetas/lienzo/fragmentos/shell.json`. Se compararon contra P3-08 y P3-12, U5-04, U5-06 y U5-07 (lo que se puede ver en el marcado), DESIGN.md §4.6, §5.4, §6.3, §7.0, §7.2, §7.3 y §7.9 a §7.15, y §10.5. También contra U1 §1.1 a §1.4, §2.1, §2.2, §2.4, §3.1, §3.2, §3.5 y §3.24, F-1, F-6 y §5. Las reglas de forma salen de `referencia-tipo/format.md`. U3 y U4 no se aplican a esta área. Solo se cruzaron dos datos: la tarjeta «Inicio» dice «6 tarjetas» y la de «Ajustes» dice «14 secciones».

## Hallazgos

| Id | Sev. | Criterio | Artboards | Qué falla | Corrección para la vuelta 2 |
|---|---|---|---|---|---|
| S1-01 | media | U5-06 · DESIGN §4.6 y §9.3 (objetivo táctil ≥ 44 px) | `shell-p03-registrar-gasto-m`, `shell-p22-registrar-peso-m`, `shell-p00-entrar-m`, `shell-barra-lateral-m` | En móvil, los `<input>` no declaran alto y su caja propia mide menos de 44 px aunque el contenedor dibujado mida 44 a 64 px. `#monto` y `#peso` miden 32 px (`.title-2` tiene interlineado de 32 px y relleno 0) dentro de una caja de 64 px. `#usuario` y `#contrasena` miden 40 px (24 + 2 × 8) dentro de `.fld`, que mide 44 px. `#q-m` mide unos 35 px: el `font:inherit` en línea anula `.body` y deja el interlineado en `normal`, más 2 × 8 de relleno, dentro de una fila de 56 px. Al tocar la franja que sobra arriba o abajo, el campo no se enfoca. Todos los demás controles móviles (enlaces, botones, chips y segmentos) miden ≥ 44 px, comprobado por script. | Dar a cada campo `align-self:stretch;min-height:var(--size-control-coarse)` (o `height:100%` dentro del contenedor), o convertir el contenedor en `<label for>`. En `#q-m`, quitar el `font:inherit` en línea o declarar `min-height:var(--size-touch-min)`. |
| S1-02 | media | U1 §1.2 punto 3 · enlaces del prototipo | `shell-barra-lateral-e` (enlace «Ajustes» de la barra y tarjeta «Ajustes»), barras de fondo de `shell-p03-registrar-gasto-e` y `shell-p22-registrar-peso-e` | En escritorio, «Ajustes» abre `../ajustes/ajustes-p26-indice-e.dc.html`. Según U1 §1.2.3, en escritorio Ajustes abre «01 Cuenta y roles», y el artboard ya existe: `ajustes/ajustes-p26-s01-cuenta-roles-e.dc.html`. El autor lo dejó así porque no conocía el slug (su Supuesto 8). En móvil, el enlace al índice sí es correcto. | Cambiar los 4 `href` de escritorio a `../ajustes/ajustes-p26-s01-cuenta-roles-e.dc.html`. |
| S1-03 | baja | U1 §5.1 punto 3 y §5.2 · DESIGN §7.11 | `shell-barra-lateral-e`, `shell-barra-lateral-m` | El menú «/» con solo «/» tecleado muestra un subconjunto sin aviso de que hay más. En escritorio salen 8 comandos fuera del orden de §5.2: /peso antes que /credito, y faltan /ingreso, /quincena, /pago, /ahorro, /lista y /salir. En móvil salen 12 y faltan /precio … /respaldo y /atajos. Los dos anchos tampoco muestran el mismo conjunto. | Seguir el orden de §5.2. Mostrar «n comandos más» o una barra de desplazamiento visible, y usar en m y e el mismo conjunto, salvo /barra, que va solo en escritorio. |
| S1-04 | baja | DESIGN §7.15 y §5.4 | `shell-p03-registrar-gasto-e`, `shell-p22-registrar-peso-e` | A los avisos «… guardado · Deshacer» les falta el icono `undo-2` junto a «Deshacer» y el botón `x` para cerrar. | Agregar `undo-2` al enlace y un botón de icono `x` con `aria-label="Cerrar aviso"`. |
| S1-05 | baja | DESIGN §7.12 · §7.3 | `shell-p22-registrar-peso-e`, `shell-p03-registrar-gasto-e` | El diálogo de P-22 es un formulario, pero usa `--layout-dialog-width-short`, que DESIGN reserva para confirmaciones; para formularios fija `-long`. En «Estados», los mensajes de error van en `small` y DESIGN pide `caption`. En P-22 e, la regla `.fld input{font:inherit}` anula `title-2 money` de `#peso-e`, así que el campo queda en 16 px y sin cifras tabulares. | Usar `--layout-dialog-width-long` o justificar el ancho corto. Poner los mensajes de error en `caption`. En `#peso-e`, quitar `font:inherit` de la regla o usar `.fld input:not(.money)`. |
| S1-06 | baja | DESIGN §7.2 «el atajo va en la descripción emergente, nunca dentro del botón» | `shell-barra-lateral-e` | «Buscar» y «Registrar gasto» de la barra llevan dentro las teclas dibujadas «Ctrl+K» y «Alt+G», y ya las tienen también en `title`. | Quitar los `<span class="kbd">` de la barra, porque el `title` ya dice el atajo. En el menú «/» las teclas dibujadas sí son correctas (§7.11). |
| S1-07 | baja | Datos sintéticos (regla «nunca inventes cifras») | `shell-p00-entrar-m`, `shell-p00-entrar-e`, `shell-p03-registrar-gasto-e` | «08:10» no aparece tal cual en `datos-sinteticos.json`: es 07:55 (`ajustes.usuarios[axel].ultimo_acceso`) más 15 (`parametros.ventana_min`). El autor lo declara como Supuesto 2. Además, el campo «Tienda» lista 3 de las 5 tiendas del JSON y faltan «Mercado Portales» y «Tianguis de la Narvarte». | El orquestador decide si agrega la hora de desbloqueo al JSON (por ejemplo en `estados_alternos`) o acepta el supuesto. Listar las 5 tiendas. |
| S1-08 | baja | U1 §1.1 y §1.2 punto 6 | `shell-barra-lateral-e` (Main) | El contenido de Main es una portada «Finanzas Hogar» con 9 tarjetas de sección, y esa pantalla no está en U1 (la entrada es P-01 Inicio). Ningún elemento de la barra está activo. Si F1 la toma de referencia, puede construir una página que no existe en U1. | Marcarla en el artboard con una nota en `caption`: «Lámina del shell: la portada no es una pantalla de U1». Otra opción es poner detrás el encabezado de P-01 con «Inicio» activo. |
| S1-09 | baja | DESIGN §6.3 | `shell-p00-entrar-e`, `shell-p03-registrar-gasto-e`, `shell-p22-registrar-peso-e` | Con `prefers-reduced-motion`, `.spin{animation:none}` detiene `loader-circle`. DESIGN lo deja girar a mitad de velocidad, porque es información. | Cambiar a `animation-duration:2s`. |

Fuera del alcance de esta área (lo anoto para el dueño de los datos): `datos-sinteticos.json` tiene `personas[axel].rol = "Titular"` y `ajustes.usuarios[axel].rol = "Co-titular"`. Los artboards usan «Co-titular», igual que U1 N-001 y `ajustes.usuarios`.

## Tabla artboard × regla

✓ = cumple · ✗ = falla (con el id del hallazgo) · n/a = no aplica.

| Artboard | Forma (a) | Tokens (b) | Contraste (c) | P3-12 (d) | Fidelidad U1 (e) | Accesibilidad (f) | Ancho (g) | Enlaces (h) | P3-08 |
|---|---|---|---|---|---|---|---|---|---|
| `shell-p00-entrar-m` (412 × 915) | ✓ | ✓ | ✓ | ✓ (S1-07 baja) | ✓ | ✗ S1-01 | ✓ | ✓ (sin enlaces: «Entrar» está bloqueado) | ✓ |
| `shell-p00-entrar-e` (1440 × 1000, fill) | ✓ | ✓ | ✓ | ✓ (S1-07) | ✓ con los 6 estados | ✓ (S1-09 baja) | ✓ | ✓ | ✓ |
| `shell-barra-lateral-m` (412 × 1854) | ✓ | ✓ | ✓ | ✓ | ✓ (S1-03 baja) | ✗ S1-01 | ✓ | ✓ | ✓ |
| `shell-barra-lateral-e` (Main, 1440 × 1160, fill) | ✓ | ✓ | ✓ | ✓ | ✓ (S1-03, S1-06 y S1-08 bajas) | ✓ | ✓ | ✗ S1-02 | ✓ |
| `shell-p03-registrar-gasto-m` (412 × 915) | ✓ | ✓ | ✓ | ✓ | ✓ | ✗ S1-01 | ✓ | ✓ | ✓ |
| `shell-p03-registrar-gasto-e` (1440 × 1200, fill) | ✓ | ✓ | ✓ | ✓ (S1-07) | ✓ con los estados de §2.2 (S1-04 y S1-05 bajas) | ✓ | ✓ | ✗ S1-02 (barra de fondo) | ✓ |
| `shell-p22-registrar-peso-m` (412 × 915) | ✓ | ✓ | ✓ | ✓ | ✓ | ✗ S1-01 | ✓ | ✓ | ✓ |
| `shell-p22-registrar-peso-e` (1440 × 900, fill) | ✓ | ✓ | ✓ | ✓ | ✓ con los estados de §2.2 (S1-04 y S1-05 bajas) | ✓ | ✓ | ✗ S1-02 (barra de fondo) | ✓ |
| Fragmento `shell.json` | ✓ | n/a | n/a | ✓ | n/a | n/a | ✓ (4 en fill, 4 sin expand) | n/a | n/a |

## (a) Forma · verificación propia (`python3 -I`, script del revisor, ya borrado)
- **8 de 8** cumplen todo esto:
  - tienen las líneas obligatorias exactas: `<!doctype html>`, `<html lang="es-MX">`, `<meta charset="utf-8">`, `<script src="./support.js"></script>`, `<link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css">`, `<x-dc>`, `<helmet>`, `</helmet>` y `</x-dc>`;
  - el `<helmet>` trae las reglas `body{margin:0;background;color;font-family}`, `a{}` y `a:hover{}`, todas con `var(--…)`;
  - el `data-dc-script` usa comilla simple, su JSON parsea y `$preview` coincide con `w` × `h` del fragmento;
  - la lógica tiene `class Component extends DCLogic` con `renderVals()`;
  - `html.parser` no encuentra etiquetas sin cerrar ni cierres sin pareja;
  - cada artboard tiene su propio `<title>`, sin duplicados.
- Nombres: solo letras minúsculas, dígitos, «-» y «.». **0 stems duplicados en todo el lienzo** (se revisaron todas las áreas).
- Fragmento: `boards` lista exactamente los 8 archivos y `order` es una permutación de `boards`. Todos llevan `page:"shell"` e `is_interactive:false`, y ninguno trae x ni y.
- 0 `innerHTML`, 0 iframes, 0 URLs de red y 0 controles de icono sin nombre accesible. Los 121 `<svg>` tienen `aria-hidden` efectivo.
- Para confirmar que el verificador detecta errores, lo corrí sobre copias alteradas. Detectó la línea de tokens.css cambiada, una etiqueta sin cerrar, `#fff`, un `var` inexistente, la cifra canon `883.75`, una fecha ISO visible, el monto `$1234.5`, un botón de icono sin nombre, un `$preview` distinto del fragmento y un enlace de escritorio que apuntaba a un artboard móvil.

## (b) Tokens
- **0 hex literales** fuera de `data-props`, y también 0 `rgb()` y 0 `hsl()`. Los SVG solo usan `fill="none"` y `stroke` con `currentColor`, puesto por la clase.
- **Todos los `var(--x)` existen**: entre 37 y 55 distintos por archivo, comprobados contra los 365 `name` de `tokens.json` más `font-sans` (de `type.families`). También están definidos en `tokens.css`.

## (c) Contraste (U5-04 en lo estático, P3-07): recálculo propio
Usé la fórmula WCAG 2.x y los hex de `tokens.json` resueltos alias por alias, en claro y en oscuro. Evalué 46 pares de texto o UI sobre su fondo, tomados de lo que usa el marcado: **0 bajo umbral**. Los pares de color de persona se marcan como decorativos porque siempre van con su nombre, y aun así cumplen 3:1. El fondo atenuado bajo `--ui-overlay-scrim` está marcado `aria-hidden` o `inert` y no se evalúa. La columna «En §10.5» dice si el par aparece con ese nombre en DESIGN §10.5. Los roles allí se escriben `color.role.*`, así que «no» solo quiere decir que el nombre no coincide: igual se calcularon y cumplen.

| # | Tipo | Primer plano / fondo | Uso | Claro | Oscuro | Umbral | En §10.5 | Resultado |
|---|---|---|---|---|---|---|---|---|
| 1 | texto | `ui-text-default` / `ui-surface-page` | P-00 m cuerpo; Main; fondo de Estados | 11.53 (#36363a/#fafafb) | 14.07 (#f2f2f4/#222325) | 4.5 | sí | cumple |
| 2 | texto | `ui-text-secondary` / `ui-surface-page` | subtítulos, migas, «Requiere conexión» | 5.40 (#66676e/#fafafb) | 10.43 (#d1d2d8/#222325) | 4.5 | sí | cumple |
| 3 | texto | `ui-text-default` / `ui-surface-card` | tarjeta P-00 e, encabezado barra, tiles, .st | 10.76 (#36363a/#f2f2f4) | 10.76 (#f2f2f4/#36363a) | 4.5 | sí | cumple |
| 4 | texto | `ui-text-secondary` / `ui-surface-card` | elementos de barra en reposo, kbd, descripciones | 5.04 (#66676e/#f2f2f4) | 7.97 (#d1d2d8/#36363a) | 4.5 | sí | cumple |
| 5 | texto | `ui-text-default` / `ui-surface-raised` | hojas, diálogos, menú «/», avisos | 11.53 (#36363a/#fafafb) | 10.76 (#f2f2f4/#36363a) | 4.5 | sí | cumple |
| 6 | texto | `ui-text-secondary` / `ui-surface-raised` | barra inferior, alias, «Más opciones» | 5.40 (#66676e/#fafafb) | 7.97 (#d1d2d8/#36363a) | 4.5 | sí | cumple |
| 7 | texto | `ui-text-default` / `ui-surface-field` | valores de campos | 11.53 (#36363a/#fafafb) | 14.07 (#f2f2f4/#222325) | 4.5 | sí | cumple |
| 8 | texto | `ui-text-secondary` / `ui-surface-field` | prefijo $ y sufijo kg | 5.40 (#66676e/#fafafb) | 10.43 (#d1d2d8/#222325) | 4.5 | sí | cumple |
| 9 | texto | `ui-text-default` / `ui-surface-selected` | elemento activo, chip y segmento elegidos | 10.73 (#36363a/#f0f1ff) | 10.90 (#f2f2f4/#303258) | 4.5 | sí | cumple |
| 10 | texto | `ui-text-default` / `ui-surface-hover` | fila resaltada del menú, letra del avatar | 9.64 (#36363a/#e5e6e9) | 7.42 (#f2f2f4/#4d4e53) | 4.5 | sí | cumple |
| 11 | texto | `ui-text-secondary` / `ui-surface-hover` | alias en la fila resaltada | 4.51 (#66676e/#e5e6e9) | 5.50 (#d1d2d8/#4d4e53) | 4.5 | sí | cumple |
| 12 | texto | `ui-text-on-primary` / `ui-action-primary-bg` | botón primario y «+» | 5.77 (#ffffff/#5d609f) | 10.47 (#1e203a/#c9cfff) | 4.5 | sí | cumple |
| 13 | texto | `ui-text-on-primary` / `ui-action-primary-bg-hover` | primario en hover | 8.49 (#ffffff/#46487b) | 12.66 (#1e203a/#e1e4ff) | 4.5 | sí | cumple |
| 14 | texto | `ui-text-default` / `ui-role-error-surface` | franja de bloqueo y «No se guardó» | 10.72 (#36363a/#ffeeed) | 11.20 (#f2f2f4/#5a2124) | 4.5 | no | cumple |
| 15 | texto | `ui-text-default` / `ui-role-warning-surface` | franja «Sin conexión» | 10.74 (#36363a/#fff0dc) | 10.86 (#f2f2f4/#4c3000) | 4.5 | no | cumple |
| 16 | texto | `ui-text-default` / `ui-role-info-surface` | franja «Sesión vencida» | 10.83 (#36363a/#e7f5ff) | 10.59 (#f2f2f4/#003b5a) | 4.5 | no | cumple |
| 17 | texto | `ui-role-error` / `ui-surface-card` | mensajes de error junto al campo (Estados) | 5.40 (#a34449/#f2f2f4) | 5.82 (#fc9999/#36363a) | 4.5 | no | cumple |
| 18 | texto | `ui-text-link` / `ui-role-error-surface` | enlace «Reintentar» en la franja de error | 5.14 (#5d609f/#ffeeed) | 8.26 (#c9cfff/#5a2124) | 4.5 | no | cumple |
| 19 | texto | `ui-text-link` / `ui-surface-raised` | enlace «Deshacer» en los avisos | 5.53 (#5d609f/#fafafb) | 7.94 (#c9cfff/#36363a) | 4.5 | sí | cumple |
| 20 | texto | `ui-text-link-hover` / `ui-surface-raised` | «Deshacer» en hover | 8.14 (#46487b/#fafafb) | 9.59 (#e1e4ff/#36363a) | 4.5 | no | cumple |
| 21 | texto | `ui-text-disabled` / `ui-surface-disabled` | «Entrar» deshabilitado (exento WCAG 1.4.3) | 5.04 (#66676e/#f2f2f4) | 6.09 (#b7b8c0/#36363a) | 4.5 | sí | cumple |
| 22 | UI | `ui-icon-default` / `ui-surface-card` | iconos de la barra | 5.04 (#66676e/#f2f2f4) | 7.97 (#d1d2d8/#36363a) | 3.0 | sí | cumple |
| 23 | UI | `ui-icon-default` / `ui-surface-raised` | iconos del menú, cerrar | 5.40 (#66676e/#fafafb) | 7.97 (#d1d2d8/#36363a) | 3.0 | no | cumple |
| 24 | UI | `ui-icon-default` / `ui-surface-field` | botón «Mostrar contraseña», icono de fecha | 5.40 (#66676e/#fafafb) | 10.43 (#d1d2d8/#222325) | 3.0 | no | cumple |
| 25 | UI | `ui-icon-selected` / `ui-surface-selected` | logo, píldora activa, icono activo | 5.15 (#5d609f/#f0f1ff) | 8.04 (#c9cfff/#303258) | 3.0 | sí | cumple |
| 26 | UI | `ui-role-error` / `ui-role-error-surface` | icono de franja de error | 5.37 (#a34449/#ffeeed) | 6.05 (#fc9999/#5a2124) | 3.0 | no | cumple |
| 27 | UI | `ui-role-warning` / `ui-role-warning-surface` | icono wifi-off | 5.12 (#8b5d00/#fff0dc) | 6.09 (#e5ae5b/#4c3000) | 3.0 | no | cumple |
| 28 | UI | `ui-role-info` / `ui-role-info-surface` | icono info | 4.96 (#006fa4/#e7f5ff) | 6.12 (#6bc3fd/#003b5a) | 3.0 | no | cumple |
| 29 | UI | `ui-role-success` / `ui-surface-raised` | icono de aviso guardado | 5.08 (#207b3e/#fafafb) | 6.42 (#7ecf90/#36363a) | 3.0 | no | cumple |
| 30 | UI | `ui-border-control` / `ui-surface-page` | borde de campo P-00 m | 3.71 (#808189/#fafafb) | 7.96 (#b7b8c0/#222325) | 3.0 | sí | cumple |
| 31 | UI | `ui-border-control` / `ui-surface-card` | borde de campo P-00 e, chips | 3.47 (#808189/#f2f2f4) | 6.09 (#b7b8c0/#36363a) | 3.0 | sí | cumple |
| 32 | UI | `ui-border-control` / `ui-surface-raised` | bordes de campos y chips en hoja/diálogo | 3.71 (#808189/#fafafb) | 6.09 (#b7b8c0/#36363a) | 3.0 | no | cumple |
| 33 | UI | `ui-border-error` / `ui-surface-card` | campo con error (Estados) | 5.40 (#a34449/#f2f2f4) | 5.82 (#fc9999/#36363a) | 3.0 | sí | cumple |
| 34 | UI | `ui-border-focus` / `ui-surface-raised` | anillo de foco del monto | 3.82 (#7579c2/#fafafb) | 6.03 (#adb3f4/#36363a) | 3.0 | no | cumple |
| 35 | UI | `ui-border-focus` / `ui-surface-page` | anillo de foco en página | 3.82 (#7579c2/#fafafb) | 7.89 (#adb3f4/#222325) | 3.0 | sí | cumple |
| 36 | UI | `ui-border-focus` / `ui-surface-card` | anillo de foco en tarjeta/barra | 3.56 (#7579c2/#f2f2f4) | 6.03 (#adb3f4/#36363a) | 3.0 | sí | cumple |
| 37 | UI | `ui-border-selected` / `ui-surface-raised` | borde de chip/segmento elegido | 5.53 (#5d609f/#fafafb) | 7.94 (#c9cfff/#36363a) | 3.0 | no | cumple |
| 38 | UI | `ui-border-selected` / `ui-surface-card` | segmento de tema elegido | 5.16 (#5d609f/#f2f2f4) | 7.94 (#c9cfff/#36363a) | 3.0 | no | cumple |
| 39 | UI | `ui-action-primary-bg` / `ui-surface-raised` | relleno de «Guardar» en hoja | 5.53 (#5d609f/#fafafb) | 7.94 (#c9cfff/#36363a) | 3.0 | sí | cumple |
| 40 | UI | `ui-action-primary-bg` / `ui-surface-page` | relleno primario en página | 5.53 (#5d609f/#fafafb) | 10.38 (#c9cfff/#222325) | 3.0 | sí | cumple |
| 41 | UI | `ui-action-primary-bg` / `ui-surface-card` | relleno primario en tarjeta | 5.16 (#5d609f/#f2f2f4) | 7.94 (#c9cfff/#36363a) | 3.0 | sí | cumple |
| 42 | UI | `ui-person-axel` / `ui-surface-card` | anillo del avatar (decorativo, con nombre) | 11.93 (#002e64/#f2f2f4) | 3.10 (#4f83c8/#36363a) | 3.0 | no | cumple |
| 43 | UI | `ui-person-axel` / `ui-surface-raised` | marca de persona (decorativa, con nombre) | 12.78 (#002e64/#fafafb) | 3.10 (#4f83c8/#36363a) | 3.0 | no | cumple |
| 44 | UI | `ui-person-ximena` / `ui-surface-selected` | marca de Ximena elegida (decorativa) | 4.83 (#a14e78/#f0f1ff) | 7.02 (#ffacd3/#303258) | 3.0 | no | cumple |
| 45 | UI | `ui-person-karen` / `ui-surface-raised` | marca de Karen (decorativa) | 3.38 (#b870ac/#fafafb) | 4.68 (#d389c6/#36363a) | 3.0 | no | cumple |
| 46 | UI | `ui-person-shared` / `ui-surface-raised` | marca Conjunto (decorativa) | 8.14 (#484b66/#fafafb) | 10.09 (#e7eaff/#36363a) | 3.0 | no | cumple |

Pares evaluados: 46 · bajo umbral (no exentos): 0

## (d) P3-12 · cifras sintéticas y credenciales
- Busqué con `grep -rnoE` el patrón completo de P3-12 en `project/shell` y `fragmentos/shell.json`, con y sin separador de miles: **0**.
- Patrón de credenciales de P3-12 sobre las mismas rutas: **0**.
- `\$[0-9]{4,}\.[0-9]{2}` (P3-08): **0**.
- Fechas: las únicas son 21-oct-2026 (hoy) y 19-oct-2026, más sus ISO, que solo están en `renderVals()`. El comando de P3-12 sobre `oraculo/model.json`, `datos/semilla/data.json` y el JSON sintético da 15 coincidencias en total, y **ninguna de las dos fechas del shell está entre ellas**. 0 apariciones de 2026-10-19 o 2026-10-21 en el modelo y en la semilla.
- Revisé a mano los 8 artboards: no hay correo, nombre de cuenta ni datos de facturación. «axel» es el `usuario` sintético de `ajustes.usuarios`.
- Muestra de 10 valores comparados con `datos-sinteticos.json`:

| # | Valor en el artboard | Fuente en el JSON | Coincide |
|---|---|---|---|
| 1 | «Disponible hoy $6,881.60» (barra-m) | `creditos.kueski.disponible_hoy.B_fmt` (lectura activa B) | sí |
| 2 | Monto «134.50» · Café y pan · Ximena · De contado · Hogar (P-03 m) | `gastos[g28]` (13450 centavos) | sí |
| 3 | Monto «411.35» · Bodega Aurrera · lácteos y reposición · Súper de la semana (P-03 e) | `gastos[g29]` (41135 centavos) | sí |
| 4 | «Antes: 84.6 kg el 19-oct-2026» (P-22 m) | `salud.perfiles[axel].peso_actual_kg` / `peso_fecha` | sí |
| 5 | «Antes: 62.8 kg el 19-oct-2026» (P-22 e) | `salud.perfiles[ximena]` | sí |
| 6 | «90 días» de sesión | `ajustes.parametros.sesion_dias` | sí |
| 7 | «5 intentos fallidos» y «15 minutos» | `parametros.intentos_fallidos` / `ventana_min` | sí |
| 8 | «Foto o PDF de hasta 10 MB» | `parametros.adjunto_max_mb` | sí |
| 9 | «Mié 21 oct 2026» | `meta.hoy` y `meta.dia_semana` | sí |
| 10 | Las 6 categorías sugeridas y su orden | `categorias[].num_gastos` con desempate por `orden` (N-050) | sí: la 6.ª queda en Despensa común (orden 3), antes que Salud, Ropa y Mascotas |
| extra | Líneas en orden alfabético (Banorte Oro, Crédito Coppel, Kueski Pay, Préstamo personal) | `creditos.lineas` y `cuentas_y_formas_de_pago` | sí |
| extra | «08:10» | no está literal (07:55 + 15 min) | derivado: S1-07 |

## (e) Fidelidad con U1 (§2.1, §2.2, §2.4 y §3)

| Pantalla | Elemento de la especificación | m | e |
|---|---|---|---|
| P-00 | Usuario y contraseña con `<label for>`, «Mostrar contraseña» | ✓ | ✓ |
| P-00 | Bloqueo como mensaje «Demasiados intentos. Vuelve a intentar a las HH:MM» (5 intentos en 15 min, N-010 y N-011) | ✓ (estado principal) | ✓ (tarjeta Estados) |
| P-00 | Sesión larga por dispositivo (decisión A) | ✓ «90 días en este teléfono» | ✓ «en este dispositivo» |
| P-00 | Forma: formulario a todo el ancho en móvil; centrado y de ancho fijo en escritorio (§6.4) | ✓ | ✓ (`max-width:25rem`) |
| P-00 | Estados de §2.2: vacío «—», carga (botón en espera y campos bloqueados), error «Usuario o contraseña incorrectos», bloqueo, sin servidor (franja G-07 y «Requiere conexión», §2.4), sesión vencida | n/a (van en escritorio) | ✓ los 6 |
| P-00 | Sin recuperación de contraseña en pantalla (§2.2) | ✓ | ✓ |
| G-01 | Encabezado «Finanzas Hogar» y usuario con su menú (Cuenta, Apariencia, Cerrar sesión) | ✓ | ✓ |
| G-01 | Dos acciones fijas: Buscar y Registrar gasto | ✓ | ✓ |
| G-01 | 9 secciones en orden fijo, con ▸ solo donde hay subpáginas (Inicio y Salud no llevan) | ✓ | ✓ |
| G-01 | Hijos dinámicos solo en Créditos, en orden alfabético (N-005); solo la sección abierta va desplegada | ✓ (Kueski Pay activa) | ✓ (Main sin página abierta) |
| G-01 | Pie con el tema Sistema · Claro · Oscuro | ✓ | ✓ |
| G-01 | Colapsable en escritorio (`panel-left`, Ctrl+\, /barra) | n/a | ✓ |
| G-01 | Panel ☰ al 85 % con fondo atenuado en móvil (§6.2) | ✓ | n/a |
| G-01 | «Tocar el nombre abre la primera subpágina»; en escritorio, Ajustes abre «01 Cuenta y roles» (§1.2.3) | ✓ | ✗ S1-02 |
| G-03/G-04 | Un solo panel; con «/» muestra comandos agrupados, alias y teclas dibujadas en escritorio | ✓ (S1-03) | ✓ (S1-03) |
| G-05 | Barra inferior: Inicio · Comidas · «+» · Compras · Buscar, con iconos de §5.4, píldora activa y «+» de 40 px | ✓ | n/a |
| P-03 | Selector de tipo (Gasto · Ingreso · Peso · Ticket) | ✓ | ✓ |
| P-03 | Monto grande enfocado, `inputmode="decimal"`, prefijo $, tabular | ✓ | ✓ |
| P-03 | 6 chips sugeridos (N-050) + «Más» | ✓ (2 columnas, 44 px) | ✓ (3 columnas) |
| P-03 | «Más opciones»: fecha, persona, reparto (Hogar, Karen, Conjunto), forma de pago (De contado y las líneas en orden alfabético, N-053), tienda, nota, adjunto | ✓ plegado con resumen de valores por defecto (N-052) | ✓ desplegado (S1-07: tiendas) |
| P-03 | Guardar (48 px abajo en la hoja; Ctrl+Enter en escritorio) y Cerrar | ✓ | ✓ |
| P-03 | Estados de §2.2: vacío no aplica, carga «la hoja no se cierra», error junto al campo (monto ≤ 0, falta categoría), «No se guardó · Reintentar», aviso «Gasto guardado · Deshacer» y su variante sin conexión (§2.4) | n/a | ✓ |
| P-22 | Persona, campo de peso enfocado con kg, último peso (N-241), Guardar | ✓ | ✓ (con fecha) |
| P-22 | Estados de §2.2: carga, error junto al campo, aviso y su variante sin conexión | n/a | ✓ |
| F-1 | P-01 → «+» → P-03 → categoría → Guardar → P-01 (3 toques) | ✓ por los enlaces | ✓ |
| F-6 | P-03 → pestaña Peso → P-22 → Guardar | ✓ por los enlaces | ✓ (desde Salud) |

U3 y U4 no se aplican: el shell no tiene tarjetas de Inicio ni campos de Ajustes. Lo que sí cruzan es coherente: «Resumen del hogar en 6 tarjetas», «14 secciones» y Ajustes plegado por defecto.

## (f) Accesibilidad (U5-06 y U5-07, en lo que se ve en el marcado)
- Los controles son `button`, `a`, `input`, `select`, `textarea`, `label for`, `fieldset/legend` y grupos con `aria-labelledby`, todos reales.
- Los controles de icono llevan `aria-label` y los SVG decorativos, `aria-hidden`.
- Diálogos y hojas: `role="dialog"`, `aria-modal` y `aria-labelledby`. Menú: combobox, listbox y `aria-activedescendant`. Segmentos: `aria-pressed`. Barra: `aria-current`. Cada artboard trae su regla de `:focus-visible`.
- Objetivos táctiles en móvil: lo medí con un script sobre alturas en línea y de clase. Todos los enlaces, botones, chips y segmentos miden ≥ 44 px, la barra inferior 56 px y el primario de hoja 48 px. Los campos, no: S1-01.
- Sin emoji, sin barra de estado ni teclado falsos. Los textos están en español de México y tratan de tú.
- La lista de palabras en inglés de U5-07 coincide con «Error» en los títulos de las tarjetas «Estados» de P-00 e, P-03 e y P-22 e. Es palabra española (U1 §2.2 usa esa misma columna) y es una nota del lienzo, no texto de la app. Lo anoto como falso positivo para U5.
- Dinero: `$6,881.60` lleva `.money` (tabular). Los montos capturados van con prefijo `$` y `.money`. No hay negativos, así que no aplica el signo U+2212. Fechas `21-oct-2026` y «Mié 21 oct 2026», sin ISO visible.

## (g) Anchos
- **Móvil**: raíz `width:412px; box-sizing:border-box`. Altos: 915 en P-00, P-03 y P-22. La barra mide 1854 = 915 + 24 (gap) + 915, coherente con `$preview`.
- **Escritorio**: raíz sin ancho en px (solo `min-height`); barra `flex:1 1 var(--layout-sidebar-width)` y contenido `flex:999 1 560px` dentro de una fila con `flex-wrap`, como pide format.md. Los contenedores usan `max-width` en rem o en token. En el fragmento, los 4 llevan `expand:"fill"` y `w:1440`.
- Las alturas de escritorio son estimadas: K45 no permite renderizar. P-00 e podría pasar un poco de 1000 (cuento unos 1080 con la tarjeta «Estados» de 4 + 2). No es falla: la página es fluida con `fill`.

## (h) Enlaces
Las 116 referencias existen en disco dentro del lienzo y apuntan a un artboard del mismo ancho, o a «#». La excepción es S1-02, que existe y tiene el mismo ancho pero no es el destino que fija U1.

## P3-08 · lista de 8 puntos (área shell, vuelta 1)
1. Resumen con ≤ 6 tarjetas: no aplica; la tarjeta de Main dice «6 tarjetas» ✓.
2. 14 secciones de Ajustes: ✓ («14 secciones»; Ajustes plegado).
3. Barra lateral colapsable: ✓ (`panel-left`, Ctrl+\, /barra).
4. Español de México: ✓.
5. Dinero `$1,234.56` con cifras tabulares: ✓.
6. Fechas `10-oct-2026`: ✓.
7. Paleta de C2, índigo 280°: ✓ (solo tokens `ui-*`).
8. Notion como única referencia: ✓ (barra y menú «/» al estilo Notion; ninguna otra referencia).

Ninguna decisión queda contradicha.

## Líneas de evidencia
- **P3-04**: este revisor usó 0 herramientas de navegador y 0 llamadas a Artifact. Journal del flujo consultado: `/home/deepsec/.claude/projects/-home-deepsec/503eb21a-a9cc-4b4f-b2cc-7c818f38ec4e/subagents/workflows/wf_73bdc7c4-bb8/journal.jsonl`. Tiene 0 menciones de `mcp__claude-in-chrome__` y 0 coincidencias del patrón de credenciales. Ese journal guarda entradas `started` y `result`, no cada llamada de herramienta, así que la verificación de llamadas es parcial.
- **P3-06 (área shell)**: 4 pantallas (P-00, G-01 con su menú «/», P-03 y P-22) × 2 anchos = 8 artboards, y no falta ninguno. Los flujos F-1 y F-6 se pueden seguir por los enlaces.
- **P3-07**: este archivo. 46 pares recalculados (el criterio pide 10).
- **P3-08**: 8 de 8, y 0 montos ≥ 1,000 sin separador.
- **P3-12**: 0 cifras canon, 0 credenciales, fechas sin coincidencias, 3 artboards revisados a mano (los 8, de hecho).
- **P3-13**: el reporte del autor (`docs/estado/wfp3/artboards-shell.md`) trae Supuestos (12) y Bloqueos (ninguno).

## Texto con aspecto de instrucción
- `referencia-tipo/format.md` incluye frases imperativas del publicador del tipo (por ejemplo, «read it first», «never lift published design source…»). Las tomé como datos que fijan la forma. Ninguna es una orden para este revisor y no actué por ellas.
- En los 8 artboards, el fragmento y el reporte del autor no hay texto dirigido a quien revisa.
