# P3-iter-2 · área shell · revisión independiente (vuelta 2 de 3)

Revisor: `revisor-shell-r2` · 12-oct-2026 · Criterio P3-07 (`docs/criterios/P3-claude-design.md` 1.3). No corrige: solo evalúa.

**Veredicto: APROBADO** · **0 altos · 0 medios** · 8 bajos (7 siguen abiertos desde la vuelta 1 y 1 es nuevo). Para aprobar hacen falta 0 altos y 0 medios.

**Qué se revisó:** los 8 artboards de `docs/design/maquetas/lienzo/project/shell/` y el fragmento `docs/design/maquetas/lienzo/fragmentos/shell.json`. Se compararon con:

- P3-08 y P3-12.
- U5-04, U5-06 y U5-07, en lo que se puede ver en el marcado.
- DESIGN.md: §4.6, §7.0, §7.2, §7.3, §7.9, §7.11, §7.12, §7.15, §9.3 y §10.5.
- U1: §1.2, §2.1, §2.2, §2.4, §3.5, F-1 y F-6.
- La forma que pide `referencia-tipo/format.md`.

U3 y U4 no se aplican a esta área. Solo se comprobó que las tarjetas digan «6 tarjetas» y «14 secciones».

Contra la vuelta 1, `diff` encontró **8 líneas cambiadas en 7 archivos**. `shell-p00-entrar-e` no cambió y el fragmento tampoco. Además de leer esos cambios, volví a hacer todas las comprobaciones (a) a (h) con scripts propios.

## Hallazgos de la vuelta 1: ¿se cerraron?

| Id | Sev. | Estado | Comprobación |
|---|---|---|---|
| S1-01 | media | **CERRADO** | Medí las alturas declaradas con un script, sin navegador (K45). `#monto` y `#peso` llevan `min-height:var(--size-control-coarse)` (44 px), `align-self:stretch` y relleno 0, y ocupan la caja de 64 px. `#q-m` es `type="search"`, así que el navegador usa border-box: con `min-height:var(--size-touch-min)` mide 44 px. `#usuario` y `#contrasena` (regla `.fld input`) miden 44 + 2 × 8 = **60 px**, porque usan content-box (ver S2-01). En los 4 móviles, los 76 controles miden 44 px o más; los campos ya no tienen franjas muertas. |
| S1-02 | media | **CERRADO** | En escritorio quedan 0 `href` a `ajustes-p26-indice-e`. Los 4 enlaces (barra y tarjeta de Main en `shell-barra-lateral-e`, barras de fondo de `shell-p03-registrar-gasto-e` y `shell-p22-registrar-peso-e`) apuntan a `../ajustes/ajustes-p26-s01-cuenta-roles-e.dc.html`, que existe y es del mismo ancho. En móvil se mantiene `ajustes-p26-indice-m`, que es correcto según U1 §1.2.3. |
| S1-03 | baja | sigue abierto | El menú «/» sigue mostrando conjuntos distintos en m (12 comandos) y en e (8). En e, el orden no es el de §5.2 y no dice cuántos comandos faltan. |
| S1-04 | baja | sigue abierto | Los avisos «… guardado · Deshacer» de P-03 e y P-22 e siguen sin el icono `undo-2` y sin el botón `x` «Cerrar aviso». |
| S1-05 | baja | sigue abierto | El diálogo de P-22 e usa `--layout-dialog-width-short` (línea 72). `.fld input{font:inherit}` (especificidad 0,1,1) le sigue ganando a `.title-2` y `.money` de `#peso-e`: el campo queda en 16 px y sin cifras tabulares. |
| S1-06 | baja | sigue abierto | «Buscar» y «Registrar gasto» de la barra siguen dibujando `Ctrl+K` y `Alt+G` dentro del botón. |
| S1-07 | baja | sigue abierto | La hora «08:10» todavía no está en `datos-sinteticos.json`: hay 0 coincidencias, y sale de sumar 07:55 + 15 min. El campo Tienda de P-03 e lista 3 de las 5 tiendas. |
| S1-08 | baja | sigue abierto | Main sigue siendo una portada que no está en U1, sin nota que lo aclare. |
| S1-09 | baja | sigue abierto | `@media (prefers-reduced-motion: reduce){.spin{animation:none}}` sigue en P-00 e, P-03 e y P-22 e. |

## Hallazgos nuevos

| Id | Sev. | Criterio | Artboard | Qué falla | Corrección sugerida |
|---|---|---|---|---|---|
| S2-01 | baja | DESIGN §9.3 («con puntero grueso los controles miden 44 px de alto»), coherencia m/e | `shell-p00-entrar-m` | El arreglo de S1-01 dio a `.fld input` `min-height:var(--size-control-coarse)` y dejó el relleno `var(--space-2) 0`. El Supuesto de `correccion-shell-r1.md` dice que hay un reinicio a border-box, pero no existe: no hay ninguna regla `box-sizing` en `sistema/project/tokens.css` ni en el `<helmet>`. Un `<input type="text">` o `"password"` usa content-box por omisión, así que cada campo mide 44 + 16 = **60 px** y su caja `.fld` crece a **62 px**. En P-00 e, la misma caja mide 44 px. El objetivo táctil se cumple de sobra; lo que falla es el tamaño de control que fija DESIGN. | En `.fld input`, cambiar `padding:var(--space-2) 0` por `padding:0`, o agregar `box-sizing:border-box`. Así el campo mide 44 px dentro de una caja de 46 px. |

## Tabla artboard × regla

✓ = cumple · ✗ = falla (con el id del hallazgo) · n/a = no aplica. Todas las bajas quedan entre paréntesis.

| Artboard | Forma (a) | Tokens (b) | Contraste (c) | P3-12 (d) | Fidelidad (e) | Accesibilidad (f) | Ancho (g) | Enlaces (h) | P3-08 |
|---|---|---|---|---|---|---|---|---|---|
| `shell-p00-entrar-m` (412 × 915) | ✓ | ✓ | ✓ | ✓ (S1-07) | ✓ | ✓ (S2-01) | ✓ | ✓ (sin enlaces: «Entrar» bloqueado) | ✓ |
| `shell-p00-entrar-e` (1440 × 1000, fill) | ✓ | ✓ | ✓ | ✓ (S1-07) | ✓ con 6 estados | ✓ (S1-09) | ✓ | ✓ | ✓ |
| `shell-barra-lateral-m` (412 × 1854) | ✓ | ✓ | ✓ | ✓ | ✓ (S1-03) | ✓ | ✓ | ✓ | ✓ |
| `shell-barra-lateral-e` (Main, 1440 × 1160, fill) | ✓ | ✓ | ✓ | ✓ | ✓ (S1-03, S1-06, S1-08) | ✓ | ✓ | ✓ | ✓ |
| `shell-p03-registrar-gasto-m` (412 × 915) | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `shell-p03-registrar-gasto-e` (1440 × 1200, fill) | ✓ | ✓ | ✓ | ✓ (S1-07) | ✓ con los estados de §2.2 (S1-04, S1-05) | ✓ (S1-09) | ✓ | ✓ | ✓ |
| `shell-p22-registrar-peso-m` (412 × 915) | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| `shell-p22-registrar-peso-e` (1440 × 900, fill) | ✓ | ✓ | ✓ | ✓ | ✓ con los estados de §2.2 (S1-04, S1-05) | ✓ (S1-09) | ✓ | ✓ | ✓ |
| Fragmento `shell.json` | ✓ | n/a | n/a | ✓ | n/a | n/a | ✓ (4 en fill, 4 sin expand) | n/a | n/a |

## (a) Forma · verificación propia (`python3 -I`, scripts en `datos/tmp/rev-shell-r2/`, borrados al terminar)

- **Los 8 artboards cumplen todo esto:**
  - **Líneas obligatorias.** Están exactas: `<!doctype html>`, `<html lang="es-MX">`, `<meta charset="utf-8">`, `<script src="./support.js"></script>`, `<link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css">`, `<x-dc>`, `<helmet>`, `</helmet>` y `</x-dc>`.
  - **Reglas del `<helmet>`.** Trae `body{margin:0;background:var(--…);color:var(--…);font-family:var(--…)}`, más `a{}` y `a:hover{}`.
  - **Script de datos.** Es `<script type="text/x-dc" data-dc-script data-props='…'>`, su JSON parsea y `$preview` coincide con el `w` × `h` del fragmento en los 8. Tiene `class Component extends DCLogic` y `renderVals()`.
  - **Etiquetas.** `html.parser` no encuentra etiquetas sin cerrar ni desparejadas.
  - **Títulos.** Cada artboard tiene su `<title>` y no se repite ninguno.
  - **ids.** No hay ids repetidos. Los `for`, `aria-labelledby`, `aria-describedby`, `aria-controls` y `aria-activedescendant` apuntan a ids que existen.
- **Nombres de archivo.** Los 90 archivos del lienzo cumplen `[a-z0-9.-]+/[a-z0-9.-]+\.dc\.html` y no hay **ningún stem repetido**.
- **Fragmento.** JSON válido. `boards` lista exactamente los 8 archivos y `order` es una permutación sin repeticiones. Todos llevan `page:"shell"` e `is_interactive:false`, y ninguno trae x ni y.
- **Lo que no debe aparecer:** 0 `innerHTML`, 0 `<iframe>` y 0 URL de red. Tampoco hay `<button>`, `<input>` ni `<a>` dentro de un `<a>` (format.md: «a `<button>` or input inside it swallows the click»).
- **El verificador detecta fallas.** Lo corrí sobre copias alteradas y marcó: `$preview` distinto del fragmento, regla `body` alterada, etiqueta desparejada, hex `#fafafb`, `var(--radius-xx)` inexistente, la cifra canon `883.75`, un enlace de escritorio a un artboard móvil y un campo móvil de 40 px.

## (b) Tokens

- **0 hex literales** fuera de `data-props`, y también 0 `rgb()`, `hsl()` y `oklch()`. Los SVG solo usan `fill:none` y `stroke:currentColor`, puestos por clase.
- **Todos los `var(--x)` existen** (entre 37 y 55 distintos por archivo) entre los 365 `name` de `tokens.json`. La excepción es `--font-sans`: está en `tokens.json` como `type.families.sans`, sin entrada `name`, y en `tokens.css` como `--font-sans`. Lo usan los 90 artboards del lienzo, tal como pide el `body` obligatorio. No es falla.

## (c) Contraste (U5-04 en lo estático, P3-07): recálculo propio

**Método.** Un script hace una cascada mínima: lee las reglas del `<helmet>` (selectores de clase, etiqueta, atributo, descendiente y `>`, con su especificidad) y los estilos en línea. Así obtiene, para cada nodo de texto, icono y borde de control, el color efectivo y el fondo del ancestro más cercano. Repite el cálculo con cada regla `:hover` aplicada. Los hex salen de `tokens.json`, resolviendo alias por alias en claro y en oscuro, con la fórmula WCAG 2.x.

**Qué queda fuera del cálculo:**
- El velo `ui-overlay-scrim` y el contenido que tapa, marcados `aria-hidden`.
- Los bordes `ui-border-divider`, que DESIGN §7.0 declara decorativos.
- Las marcas de persona, que siempre van con el nombre.

**Resultado: 65 pares × 2 temas, 0 bajo umbral.** La columna «§10.5» dice si el par aparece con ese nombre en DESIGN §10.5 (con los nombres pasados de `ui.x.y` a `ui-x-y`). Los que no aparecen se calcularon igual y cumplen.

| # | Tipo | Primer plano / fondo | Claro | Oscuro | Umbral | §10.5 | Estado | Resultado |
|---|---|---|---|---|---|---|---|---|
| 1 | texto | `ui-role-error` / `ui-surface-card` | 5.40 | 5.82 | 4.5 | no | reposo | cumple |
| 2 | texto | `ui-text-default` / `ui-role-error-surface` | 10.72 | 11.20 | 4.5 | no | reposo | cumple |
| 3 | texto | `ui-text-default` / `ui-role-info-surface` | 10.83 | 10.59 | 4.5 | no | reposo | cumple |
| 4 | texto | `ui-text-default` / `ui-role-warning-surface` | 10.74 | 10.86 | 4.5 | no | reposo | cumple |
| 5 | texto | `ui-text-default` / `ui-surface-card` | 10.76 | 10.76 | 4.5 | sí | hover, reposo | cumple |
| 6 | texto | `ui-text-default` / `ui-surface-field` | 11.53 | 14.07 | 4.5 | sí | hover, reposo | cumple |
| 7 | texto | `ui-text-default` / `ui-surface-hover` | 9.64 | 7.42 | 4.5 | sí | hover, reposo | cumple |
| 8 | texto | `ui-text-default` / `ui-surface-page` | 11.53 | 14.07 | 4.5 | sí | reposo | cumple |
| 9 | texto | `ui-text-default` / `ui-surface-raised` | 11.53 | 10.76 | 4.5 | sí | hover, reposo | cumple |
| 10 | texto | `ui-text-default` / `ui-surface-selected` | 10.73 | 10.90 | 4.5 | sí | hover, reposo | cumple |
| 11 | texto | `ui-text-disabled` / `ui-surface-disabled` | 5.04 | 6.09 | 4.5 | sí | hover, reposo | cumple |
| 12 | texto | `ui-text-link` / `ui-role-error-surface` | 5.14 | 8.26 | 4.5 | no | reposo | cumple |
| 13 | texto | `ui-text-link` / `ui-surface-raised` | 5.53 | 7.94 | 4.5 | sí | reposo | cumple |
| 14 | texto | `ui-text-link-hover` / `ui-role-error-surface` | 7.56 | 9.98 | 4.5 | no | hover | cumple |
| 15 | texto | `ui-text-link-hover` / `ui-surface-card` | 7.59 | 9.59 | 4.5 | sí | hover | cumple |
| 16 | texto | `ui-text-link-hover` / `ui-surface-hover` | 6.80 | 6.62 | 4.5 | no | hover | cumple |
| 17 | texto | `ui-text-link-hover` / `ui-surface-raised` | 8.14 | 9.59 | 4.5 | no | hover | cumple |
| 18 | texto | `ui-text-on-primary` / `ui-action-primary-bg` | 5.77 | 10.47 | 4.5 | sí | reposo | cumple |
| 19 | texto | `ui-text-on-primary` / `ui-action-primary-bg-hover` | 8.49 | 12.66 | 4.5 | sí | hover | cumple |
| 20 | texto | `ui-text-secondary` / `ui-surface-card` | 5.04 | 7.97 | 4.5 | sí | hover, reposo | cumple |
| 21 | texto | `ui-text-secondary` / `ui-surface-field` | 5.40 | 10.43 | 4.5 | sí | hover, reposo | cumple |
| 22 | texto | `ui-text-secondary` / `ui-surface-hover` | 4.51 | 5.50 | 4.5 | sí | hover, reposo | cumple |
| 23 | texto | `ui-text-secondary` / `ui-surface-page` | 5.40 | 10.43 | 4.5 | sí | reposo | cumple |
| 24 | texto | `ui-text-secondary` / `ui-surface-raised` | 5.40 | 7.97 | 4.5 | sí | hover, reposo | cumple |
| 25 | UI | `ui-border-control` / `ui-surface-card` | 3.47 | 6.09 | 3.0 | sí | hover, reposo | cumple |
| 26 | UI | `ui-border-control` / `ui-surface-page` | 3.71 | 7.96 | 3.0 | sí | reposo | cumple |
| 27 | UI | `ui-border-control` / `ui-surface-raised` | 3.71 | 6.09 | 3.0 | no | hover, reposo | cumple |
| 28 | UI | `ui-border-control-hover` / `ui-surface-card` | 5.04 | 7.97 | 3.0 | sí | hover | cumple |
| 29 | UI | `ui-border-control-hover` / `ui-surface-page` | 5.40 | 10.43 | 3.0 | sí | hover | cumple |
| 30 | UI | `ui-border-control-hover` / `ui-surface-raised` | 5.40 | 7.97 | 3.0 | no | hover | cumple |
| 31 | UI | `ui-border-divider` / `ui-surface-card` | 1.12 | 1.45 | 3.0 | no | reposo | decorativo (exento) |
| 32 | UI | `ui-border-divider` / `ui-surface-hover` | 1.00 | 1.00 | 3.0 | no | hover, reposo | decorativo (exento) |
| 33 | UI | `ui-border-divider` / `ui-surface-page` | 1.20 | 1.90 | 3.0 | sí | hover, reposo | decorativo (exento) |
| 34 | UI | `ui-border-divider` / `ui-surface-raised` | 1.20 | 1.45 | 3.0 | sí | reposo | decorativo (exento) |
| 35 | UI | `ui-border-error` / `ui-surface-card` | 5.40 | 5.82 | 3.0 | sí | hover, reposo | cumple |
| 36 | UI | `ui-border-focus` / `ui-surface-field` | 3.82 | 7.89 | 3.0 | no | hover, reposo | cumple |
| 37 | UI | `ui-border-selected` / `ui-surface-card` | 5.16 | 7.94 | 3.0 | no | reposo | cumple |
| 38 | UI | `ui-border-selected` / `ui-surface-raised` | 5.53 | 7.94 | 3.0 | no | hover, reposo | cumple |
| 39 | UI | `ui-icon-default` / `ui-surface-card` | 5.04 | 7.97 | 3.0 | sí | hover, reposo | cumple · icono con aria-hidden |
| 40 | UI | `ui-icon-default` / `ui-surface-field` | 5.40 | 10.43 | 3.0 | no | hover, reposo | cumple · icono con aria-hidden |
| 41 | UI | `ui-icon-default` / `ui-surface-hover` | 4.51 | 5.50 | 3.0 | sí | hover, reposo | cumple · icono con aria-hidden |
| 42 | UI | `ui-icon-default` / `ui-surface-page` | 5.40 | 10.43 | 3.0 | sí | reposo | cumple · icono con aria-hidden |
| 43 | UI | `ui-icon-default` / `ui-surface-raised` | 5.40 | 7.97 | 3.0 | no | hover, reposo | cumple · icono con aria-hidden |
| 44 | UI | `ui-icon-selected` / `ui-surface-selected` | 5.15 | 8.04 | 3.0 | sí | hover, reposo | cumple · icono con aria-hidden |
| 45 | UI | `ui-person-axel` / `ui-surface-card` | 11.93 | 3.10 | 3.0 | no | reposo | cumple |
| 46 | UI | `ui-person-shared` / `ui-surface-raised` | 8.14 | 10.09 | 3.0 | no | reposo | cumple |
| 47 | UI | `ui-role-error` / `ui-role-error-surface` | 5.37 | 6.05 | 3.0 | no | reposo | cumple · icono con aria-hidden |
| 48 | UI | `ui-role-error` / `ui-surface-card` | 5.40 | 5.82 | 3.0 | no | reposo | cumple · icono con aria-hidden |
| 49 | UI | `ui-role-info` / `ui-role-info-surface` | 4.96 | 6.12 | 3.0 | no | reposo | cumple · icono con aria-hidden |
| 50 | UI | `ui-role-success` / `ui-surface-raised` | 5.08 | 6.42 | 3.0 | no | reposo | cumple · icono con aria-hidden |
| 51 | UI | `ui-role-warning` / `ui-role-warning-surface` | 5.12 | 6.09 | 3.0 | no | reposo | cumple · icono con aria-hidden |
| 52 | UI | `ui-text-default` / `ui-surface-card` | 10.76 | 10.76 | 3.0 | sí | reposo | cumple · icono con aria-hidden |
| 53 | UI | `ui-text-default` / `ui-surface-hover` | 9.64 | 7.42 | 3.0 | sí | hover | cumple · icono con aria-hidden |
| 54 | UI | `ui-text-default` / `ui-surface-page` | 11.53 | 14.07 | 3.0 | sí | reposo | cumple · icono con aria-hidden |
| 55 | UI | `ui-text-default` / `ui-surface-raised` | 11.53 | 10.76 | 3.0 | sí | reposo | cumple · icono con aria-hidden |
| 56 | UI | `ui-text-default` / `ui-surface-selected` | 10.73 | 10.90 | 3.0 | sí | reposo | cumple · icono con aria-hidden |
| 57 | UI | `ui-text-link-hover` / `ui-surface-card` | 7.59 | 9.59 | 3.0 | sí | hover | cumple · icono con aria-hidden |
| 58 | UI | `ui-text-link-hover` / `ui-surface-hover` | 6.80 | 6.62 | 3.0 | no | hover | cumple · icono con aria-hidden |
| 59 | UI | `ui-text-link-hover` / `ui-surface-raised` | 8.14 | 9.59 | 3.0 | no | hover | cumple · icono con aria-hidden |
| 60 | UI | `ui-text-on-primary` / `ui-action-primary-bg` | 5.77 | 10.47 | 3.0 | sí | hover, reposo | cumple · icono con aria-hidden |
| 61 | UI | `ui-text-on-primary` / `ui-action-primary-bg-hover` | 8.49 | 12.66 | 3.0 | sí | hover | cumple · icono con aria-hidden |
| 62 | UI | `ui-text-secondary` / `ui-surface-card` | 5.04 | 7.97 | 3.0 | sí | reposo | cumple · icono con aria-hidden |
| 63 | UI | `ui-text-secondary` / `ui-surface-hover` | 4.51 | 5.50 | 3.0 | sí | hover | cumple · icono con aria-hidden |
| 64 | UI | `ui-text-secondary` / `ui-surface-page` | 5.40 | 10.43 | 3.0 | sí | reposo | cumple · icono con aria-hidden |
| 65 | UI | `ui-text-secondary` / `ui-surface-raised` | 5.40 | 7.97 | 3.0 | sí | reposo | cumple · icono con aria-hidden |

Pares evaluados: 65 por tema. Bajo umbral y no exentos: **0**. Los 4 pares `ui-border-divider` bajo 3:1 son decorativos (DESIGN §7.0), y el borde con información en esos componentes es `ui-border-selected` o `ui-border-control`, que cumplen.

## (d) P3-12 · cifras sintéticas y credenciales

- **Cifras canon.** Corrí el `grep -rnoE` completo de P3-12, con y sin separador de miles, sobre `project/shell`, `fragmentos/shell.json` y la copia `docs/design/salida/iter-1/shell/`: **0** coincidencias.
- **Credenciales.** Patrón de P3-12 sobre las mismas rutas: **0**.
- **Fechas.** El comando de P3-12 sobre `oraculo/model.json`, `datos/semilla/data.json` y el JSON sintético da 15 fechas (30-oct-2026 a 15-jun-2027). Las únicas fechas del shell son 21-oct-2026 (hoy) y 19-oct-2026, más sus formas ISO dentro de `renderVals()`, y **ninguna está en esa lista**.
- **Revisión visual** de P-00 m, P-03 m, P-22 m y barra-m (marcado completo) y del texto visible de los 4 de escritorio: no hay correo, nombre de cuenta ni datos de facturación. «axel» es el `usuario` sintético.
- **Muestra de 10 valores**, comparada con `datos-sinteticos.json`:

| # | Valor en el artboard | Fuente en el JSON | Coincide |
|---|---|---|---|
| 1 | «Disponible hoy $6,881.60» (barra-m) | `creditos.kueski.disponible_hoy.B_fmt` | sí |
| 2 | Monto «134.50» · Café y pan · Ximena · De contado · Hogar (P-03 m) | `gastos[g28]` (`monto_cents` 13450) | sí |
| 3 | Monto «411.35» · Bodega Aurrera (P-03 e) | `gastos[g29]` (`monto_cents` 41135, `tienda_id` bodega-aurrera) | sí |
| 4 | «Antes: 84.6 kg el 19-oct-2026» (P-22 m) | `salud.perfiles[axel].peso_actual_kg` y `peso_fecha` | sí |
| 5 | «Antes: 62.8 kg el 19-oct-2026» (P-22 e) | `salud.perfiles[ximena]` | sí |
| 6 | «90 días» de sesión | `ajustes.parametros.sesion_dias` | sí |
| 7 | «5 intentos fallidos» y «15 minutos» | `parametros.intentos_fallidos` y `ventana_min` | sí |
| 8 | «Foto o PDF de hasta 10 MB» | `parametros.adjunto_max_mb` | sí |
| 9 | «Mié 21 oct 2026» | `meta.hoy` (2026-10-21) y `meta.dia_semana` | sí |
| 10 | Líneas en orden alfabético: Banorte Oro, Crédito Coppel, Kueski Pay, Préstamo personal | `creditos.lineas` (N-053) | sí |
| extra | «08:10» | no está tal cual (07:55 + 15 min) | derivado: S1-07 |

## (e) Fidelidad con U1 (§1.2, §2.1, §2.2, §2.4 y §3.5)

| Pantalla | Elemento de la especificación | m | e |
|---|---|---|---|
| P-00 | Usuario y contraseña con `<label for>`, «Mostrar contraseña» | ✓ | ✓ |
| P-00 | Bloqueo como mensaje «Demasiados intentos. Vuelve a intentar a las HH:MM» (5 intentos en 15 min) | ✓ (estado principal, `role="alert"`) | ✓ (tarjeta Estados) |
| P-00 | Sesión larga por dispositivo | ✓ «90 días en este teléfono» | ✓ «en este dispositivo» |
| P-00 | Formulario a todo el ancho en móvil; centrado y de ancho fijo en escritorio | ✓ | ✓ (`max-width:25rem`) |
| P-00 | Estados de §2.2 y §2.4: vacío, carga, error, bloqueo, sin servidor y sesión vencida | n/a | ✓ los 6 |
| P-00 | Sin recuperación de contraseña en pantalla | ✓ | ✓ |
| G-01 | Encabezado «Finanzas Hogar», usuario con menú, Buscar y Registrar gasto, las 9 secciones en orden fijo, pie con el tema | ✓ | ✓ |
| G-01 | Hijos solo en Créditos, en orden alfabético; Ajustes plegado | ✓ | ✓ |
| G-01 | Tocar el nombre abre la primera subpágina (Ajustes: 01 Cuenta y roles en e, índice en m, §1.2.3) | ✓ | ✓ (S1-02 cerrado) |
| G-01 | Colapsable en escritorio / panel ☰ con velo en móvil | ✓ | ✓ |
| G-03/G-04 | Menú «/» con grupos, alias y teclas | ✓ (S1-03) | ✓ (S1-03) |
| G-05 | Barra inferior de 5 destinos con «+» | ✓ | n/a |
| P-03 | Selector de tipo, monto `inputmode="decimal"` tabular, 6 chips + «Más», «Más opciones», Guardar | ✓ | ✓ |
| P-03 | «Más opciones»: fecha, persona, reparto, forma de pago, tienda, nota y adjunto, más «Descripción o comercio» (N-054) | ✓ plegado con resumen (N-052) | ✓ desplegado (S1-07: tiendas) |
| P-03 | Estados de §2.2 y §2.4: carga, error junto al campo, «No se guardó · Reintentar», aviso con Deshacer y su variante sin conexión | n/a | ✓ (S1-04) |
| P-22 | Persona, campo peso con kg, último peso, Guardar | ✓ | ✓ |
| P-22 | Estados de §2.2 y §2.4: carga, error junto al campo, aviso y su variante sin conexión | n/a | ✓ (S1-04) |
| F-1 / F-6 | P-01 → «+» → P-03 → Guardar → P-01; P-03 → Peso → P-22 → Guardar | ✓ por enlaces | ✓ por enlaces |

## (f) Accesibilidad (U5-06 y U5-07, en lo que se ve en el marcado)

- **Controles reales:** `button`, `a`, `input`, `select`, `textarea`, `label for`, `fieldset`/`legend` y grupos con `aria-labelledby`. Todo campo tiene nombre (`label for`, `aria-label` o `aria-labelledby`). Los controles de icono tienen nombre y los SVG tienen `aria-hidden` efectivo (el del fondo de barra-m lo hereda de su contenedor).
- **Objetivos táctiles en móvil:** 76 controles, **0 bajo 44 px**. Así se reparten:

| Artboard | Controles | < 44 px |
|---|---|---|
| barra-m | 48 | 0 |
| P-00 m | 4 | 0 |
| P-03 m | 15 | 0 |
| P-22 m | 9 | 0 |

  Los campos miden 44 px (`#monto`, `#peso` y `#q-m`) y 60 px (`#usuario` y `#contrasena`, ver S2-01).
- **Otros puntos de U5-06:** sin emoji, sin barra de estado ni teclado falsos (P-03 m usa `inputmode="decimal"`), y foco visible con `:focus-visible` en los 8.
- **Idioma.** De la lista de palabras en inglés de U5-07, solo aparece «Error», en las tarjetas «Estados» de P-00 e, P-03 e y P-22 e. Es palabra del español, la columna de U1 §2.2 se llama igual y es una nota del lienzo. Lo cuento como falso positivo, igual que en la vuelta 1.
- **Dinero y fechas.** Todo importe visible cumple `$1,234.56` y lleva `.money` (tabular). No hay negativos, así que no aplica U+2212. Ninguna fecha ISO es visible.

## (g) Anchos

- **Móvil:** raíz `width:412px; box-sizing:border-box`. Altos: 915 en P-00, P-03 y P-22, y 1854 en barra-m (915 + 24 + 915). Según la suma estimada, el contenido de P-03 m mide unos 690 px y cabe en 915.
- **Escritorio:** ninguna raíz tiene ancho en px (solo `min-height`). La barra y el contenido van en una fila con `flex-wrap:wrap` y `flex:999 1 560px`, y no hay `sticky` ni `100vh`. P-00 e no tiene barra lateral: es un formulario centrado con `max-width` en rem. En el fragmento, los 4 llevan `expand:"fill"` y `w:1440`.

## (h) Enlaces

Hay 116 `href` en los 8 archivos: 24 son «#» y los otros 92 apuntan a archivos que existen en el lienzo y son del mismo ancho (m → m, e → e). Hay 0 rotos y 0 de otro ancho. Los 4 de Ajustes en escritorio van a `ajustes-p26-s01-cuenta-roles-e`.

## P3-08 · decisiones (área shell, vuelta 2)

1. Resumen con ≤ 6 tarjetas: no aplica, pero la tarjeta de Main dice «6 tarjetas» ✓.
2. 14 secciones de Ajustes: ✓ («14 secciones»; Ajustes plegado).
3. Barra lateral colapsable: ✓.
4. Español de México y trato de tú: ✓.
5. Dinero `$1,234.56` con cifras tabulares: ✓ (en P-22 e el campo de peso no es dinero: S1-05).
6. Fechas `10-oct-2026`: ✓.
7. Paleta de C2, índigo 280°: ✓ (solo tokens `ui-*`).
8. Notion como única referencia: ✓.

Ninguna decisión queda contradicha.

## Líneas de evidencia

- **P3-04:** este revisor usó 0 herramientas de navegador, 0 llamadas a Artifact y 0 accesos a red.
- **P3-07:** este archivo. Se recalcularon 65 pares por tema; el criterio pide 10.
- **P3-08:** la lista de arriba, sin contradicciones.
- **P3-12:** 0 cifras canon, 0 credenciales y 0 fechas del shell entre las coincidencias de fechas.
- **U5-04:** 65 pares × 2 temas, 0 bajo umbral (estático, sin DOM renderizado).
- **U5-06:** 0 objetivos táctiles < 44 px; S2-01 es una baja de tamaño de control.
- **U5-07:** 0 palabras en inglés reales («Error» es falso positivo), 0 ISO visibles y dinero bien formado.

## Fuera del alcance de esta área (para el orquestador)

- **Copias en `salida/`.** Las copias de la vuelta 1 están en `docs/design/salida/iter-1/<area>/`, en todas las áreas, pero el grep de P3-12 busca en `docs/design/salida/iter-*/project`, así que no las recorre. Además, `salida/iter-1/shell/` no tiene copia del informe de revisión que pide P3-06. Lo corrí a mano sobre `salida/iter-1/shell` y dio 0.
- **Rol de Axel.** Sigue la diferencia en `datos-sinteticos.json`: `personas[axel].rol = "Titular"` contra `ajustes.usuarios[axel].rol = "Co-titular"`. El shell usa «Co-titular».
