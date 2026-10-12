# revisor-ds-r2 · Pista 3 fase 1 (sistema de diseño) · revisión ronda 2 de 3 · 11-oct-2026, noche

Agente `revisor-ds-r2`, revisor independiente. No corrige nada. El informe completo está en `docs/estado/wfp3/REVISION-ds-r2.md`.

## Hecho

- **Dictamen: NO APROBADO.** Quedan 1 hallazgo alto, 1 medio y 9 bajos.
- Los 4 hallazgos de la r1 (A-01, M-01, M-02 y M-03) quedan **cerrados**. El A-01 se cierra porque la excepción quedó declarada; lo que falta depende del bloqueo `P3-contraste-chip-hover.md`, pendiente de Axel.
- **A-01 (alto, nuevo):** en `BarraSuperior`, los 3 botones de solo icono («⋯» en móvil y en escritorio, y Buscar) usan `size-icon-sm` (20 px). DESIGN §5.2 y §7.2 piden `size-icon-md` (24 px) en la barra superior. Además, el README del sistema redacta esa regla incompleta.
- **M-01 (medio, nuevo):** DESIGN §7.5.1 (barra de la vista) y §7.5.5 (galería) no tienen componente, y ningún reporte ni bloqueo lo dice. `correccion-ds-r1.md` afirma que toda la matriz §7.16 está cubierta.
- **Bajos nuevos:**
  - B-01: inferencia sin marcar en VistaPrevia;
  - B-02: tope de familias contado distinto en cada archivo, con 12 en el límite;
  - B-03: foco interior en las pestañas;
  - B-04: un token de diálogo usado como ancho del marco de muestra.
- Siguen abiertos los 5 bajos de la r1.

## Verificación (comandos y salidas)

Usé scripts propios en `datos/tmp/rev-ds-r2/`, ya borrados: `v_tokens.py`, `v_prev.py`, `v_struct.py`, `v_icons.py` y `v_misc.py`.

```
$ python3 -I docs/design/tools/ds/verificar_ds.py
RESULTADO: OK · 445 tokens C2 -> 355 en listas de tokens.json + 6 estilos de tipo (18 font.*) + 2 fuentes; 0 fallas
$ python3 -I v_tokens.py 20261012
C2: 444 hojas + $root de icon.stroke = 445 · color 190 · ui 150 · font 18 · size 32 · layout 14 · space 8 · viz 8 · motion 5 · breakpoint 5 · radius 4 · shadow 3 · icon 3 · border 2 · overlay 2 · focus 1
tokens.json: color 266 · spacing 8 · radius 4 · shadow 3 · fontWeight 4 · border 2 · focus 1 · size 32 · layout 14 · icon 3 · motion 5 · breakpoint 5 · viz 8 · + type (6 estilos, 1 pila, 2 fuentes)
names 355 dup [] bad [] · var( 0 color-mix( 0 url( 0 · themes ['light','dark'] · sin usage []
alias refs 222 rotos [] valores no-hex []
muestra 20 (claro y oscuro) discrepancias [] · comparacion completa 532 discrepancias []
hex C2 154 tokens.json 154 · C2-ds set() · ds-C2 set()
font InterVariable.woff2 wOF2 352240 sha256 igual (693b77d4f32e) · InterVariable-Italic.woff2 wOF2 387976 sha256 igual (e564f652916d)
$ python3 -I v_struct.py      # 19 previews: 0 problemas (marcador, HTML, prohibidas, URL, hex, var inexistente, botón sin nombre, input sin etiqueta, lang, emoji)
claro ui-surface-hover: person-karen 2.83 · store-sams-club 2.74 | oscuro ui-surface-hover: person-axel 2.14 · store-market 2.11
raised/selected: mínimo 3.05 (claro) y 3.06 (oscuro)
$ python3 -I v_prev.py        # cascada propia, 4 estados, 2 temas
FALLA? RELLENO ChipPersonaTienda light hover/active ui-person-karen sobre ui-surface-hover 2.83   (excepción declarada)
pares totales 2466 · texto mín 4.51 (MenuBarra claro) / 5.04 oscuro · gráfico mín 3.28 claro / 3.06 oscuro (sin la excepción)
$ python3 -I v_icons.py       # 101 iconos, 0 sin coincidencia en sprite.svg (Lucide 1.55.0)
BarraSuperior: chevron-left ellipsis chevron-right search plus · PanelAjustes: user-cog wallet banknote tags credit-card store chef-hat file-chart-column bell palette database users heart-pulse package circle-check rotate-ccw
$ python3 -I v_misc.py
README: 122 vars, 0 inexistentes · no empieza con # · P3-12 cifras canon: 0 · credenciales: 0
fechas que coinciden con oraculo/semilla ['2026-10-30'] · 41 montos, ninguno con centavos coincide
size-bottombar-add 40px · layout-topbar-height 56/48 · size-bottombar-indicator 56×32 · size-sheet-handle 32×4 · size-control-primary-sheet 48 · layout-dialog-width 480/640 · layout-preview-width 400/480 · ui-critical {color-critical-light|dark} · ui-role-success {color-role-success-light|dark}
previews: 143 var(--…) distintos, todos existen
$ diff (respaldo /tmp/correccion-ds-r1-respaldo vs actual, pesos normalizados): 14 previews con 0 líneas · tokens.json sin fontWeight: 0 líneas
$ grep 'btn btn-icono' components/BarraSuperior/preview.html → 3 (Más acciones, Buscar, Más acciones), todos con svg .i20
```

## Supuestos

- **Familias:** el «≤12 families» de format.md lo leo como un tope de 12 claves «otras». Por eso B-02 queda como bajo y no como alto. Si el tipo contara las 14 familias, sería un problema de forma; se confirma al publicar.
- **Iconos de solo icono fuera de una barra** (el «x» de VistaPrevia y de AvisoBreve): DESIGN §5.2 dice «botones de solo icono» en `md`, y §7.2 dice `sm` salvo en barras. Es ambiguo, así que no lo levanto como hallazgo. La barra superior sí está nombrada en los dos lugares: de ahí sale A-01.
- **El A-01 de la r1** se da por cerrado con la corrección en dos pasos que fijó la r1. El cambio de la regla espera a Axel.

## Bloqueos

- Ninguno nuevo.
- Siguen abiertos `P3-contraste-chip-hover.md` y `P3-huecos-componentes.md`, este con los puntos 1 a 9.

## Pendientes

- **Corrección r2:**
  - A-01: `i20` → `i24` en los 3 `.btn-icono` de BarraSuperior, y ajustar su README y el del sistema;
  - M-01: componentes `Galeria` y barra de la vista, o listarlos como pendientes en el bloqueo y en el reporte;
  - opcional: B-01 a B-04 y los bajos de la r1.
- **Orquestador, al publicar:**
  - leer el `tokens.css` compilado y confirmar `--font-weight-*`, `--viz-*` y `--breakpoint-*` (B-02);
  - publicar `design-system.json` en la raíz;
  - hacer la comprobación de render de los 5 componentes nuevos en los dos temas (aquí no hubo navegador).

## Consumo aproximado

- **Leído:** ≈ 200 KB.
  - format.md, craft.md, cover.md y from-code.md (parte): 31 KB.
  - REVISION-ds-r1 y correccion-ds-r1: 24 KB.
  - README del sistema: 23 KB.
  - README de los componentes: 17 KB.
  - 5 previews nuevos y Cover, sin los trazos de los SVG: 30 KB.
  - DESIGN §4.2–4.5, §5.2, §5.4, §6.1–6.2, §7.2, §7.5, §7.9–7.13 y §7.16: 35 KB.
  - Criterios P3: 6 KB.
  - Salidas de los scripts: 30 KB.
- **Escrito:** ≈ 52 KB.
  - Scripts temporales: 28 KB, ya borrados.
  - REVISION-ds-r2.md: 19 KB.
  - Este reporte: 5 KB.
- **Acciones de herramienta:** ≈ 33.
