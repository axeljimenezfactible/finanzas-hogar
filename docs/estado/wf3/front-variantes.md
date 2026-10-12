# front-variantes · WF-3a (12-oct-2026)

Agente `front-variantes` (Sonnet). Alcance: historias de los componentes de DESIGN.md §7 que ya existen en `frontend/src/componentes/`, configuración de Storybook 10 y el script de contraste. No se tocó ninguna pantalla ni ningún componente de frontend-base.

## Hecho

- **Storybook 10.6.1** (`@storybook/react-vite`, compatible con Vite 8) en `frontend/.storybook/`. Sin addons, sin telemetría, sin red; se quitan los plugins de la PWA; Inter y los estilos salen de la app (`tokens.css` + Tailwind solo con escalas de C2). Barra de herramientas «Tema» (según el sistema, claro, oscuro: pone `data-theme` igual que `src/app/tema.ts`) y vistas de 360, 412, 768, 1024 y 1440 px.
- **37 historias** en 8 archivos `*.stories.tsx` (Botón, Campo, Chip, Tarjeta de métrica, Avisos, Estados, Hoja, Gráfica): variantes, reposo, deshabilitado, ocupado, error, con ayuda, sobre página, tarjeta, superficie flotante y fila elegida, tamaños (normal, hoja, icono), de 1 a 8 series, importe largo a 360 px, etc. Hover, foco y presionado se ven con el puntero y el teclado; el script los provoca.
- **Script de contraste** `frontend/scripts/contraste-variantes.mjs` + `.py`: sirve `storybook-static` en local, abre cada historia en claro y oscuro, en escritorio (1280 px, puntero fino) y móvil (360 px, táctil), provoca hover, presionado y foco (Tab) sobre cada elemento interactivo, mide texto, valor y marcador de campos, iconos, marcas de forma, marcas de series, bordes (afuera y adentro) y anillo de foco contra el fondo efectivo (composición de alfas por la cadena de ancestros) y calcula la razón **importando `contraste()` de `docs/design/tools/contraste_componentes.py`** (Decimal, truncada a 2 decimales; umbral 4.5 texto, 3 texto grande e interfaz). Además comprueba: color dentro de los tokens `--ui-*` (nombra cada color con su token), estilo de deshabilitado de §7.0, anillo de 2 px separado 2 px y **sin transición**, áreas táctiles (44 px grueso, 24 px fino), 0 errores de consola, 0 hosts externos y 0 desborde horizontal a 360 px.
- **Hallazgos** (`.storybook/hallazgos.json`): fallas reales de componentes que no son de mi alcance. Se informan como HALLAZGO con la corrección exacta y no cambian el código de salida; cualquier falla distinta sí lo vuelve 1, y un hallazgo corregido se informa CERRADO.

## Archivos

- `frontend/.storybook/main.ts`, `preview.tsx`, `tsconfig.json`, `apoyo.tsx`, `hallazgos.json`
- `frontend/src/componentes/{Avisos,Boton,Campo,Chip,Estados,Grafica,Hoja,TarjetaMetrica}.stories.tsx`
- `frontend/scripts/contraste-variantes.mjs`, `frontend/scripts/contraste-variantes.py`
- `frontend/package.json` (scripts y devDependencies, ver abajo) y `frontend/.gitignore` (`storybook-static/`): archivos de frontend-base, tocados solo para esto (ver Supuestos).
- `docs/estado/wf3/front-variantes.md` (este reporte).

## Comandos

```bash
cd /home/deepsec/finanzas-hogar/frontend
npm run storybook             # 127.0.0.1:6006
npm run build-storybook       # → storybook-static/ (ignorado por git)
npm run typecheck:variantes   # tsc -p .storybook/tsconfig.json
npm run verificar:variantes   # build + render + contraste (FH_CHROMIUM solo si falta el Chromium de Playwright; el script busca el chrome-headless-shell de ~/.cache/ms-playwright)
# opciones del script: --solo <texto del id de historia> (corrida parcial), --salida <carpeta> (por omisión ../datos/tmp/variantes: CSV con todas las mediciones)
```

## Verificación

- Build de Storybook: `npm run build-storybook` → «Storybook build completed successfully» (37 historias; `npm run storybook` arranca y sirve `iframe.html` e `index.json` con 200).
- Tipos: `npx tsc --noEmit -p .storybook/tsconfig.json` → 0 errores en `src/componentes/**` y `.storybook/`; `npx tsc --noEmit` (raíz) → 0 errores fuera de `src/pantallas/` (los 43 errores de esa carpeta son de agentes de área en curso: módulos aún inexistentes y tipos de `ajustes-a`; ver Bloqueos).
- `npm run verificar:variantes`: 37 historias, 148 páginas (historia x tema x contexto), **5 230 mediciones únicas** (claro 2 606, oscuro 2 624), 4 726 con razón de contraste, 232 pares de color distintos; por estado: reposo 2 434, hover 654, presionado 654, foco 1 488.
  - Cruce con U2d: 7 406 mediciones coinciden en par de hex con `docs/design/contraste-componentes.csv` (C-001..C-158) y **0** con razón distinta.
  - **FALLAS nuevas: 0**; problemas de la corrida (consola, hosts externos, desborde, render): **0**; hosts externos: 0 en escritorio y 0 en móvil; `RESULTADO: OK`.
  - Prueba negativa: una historia temporal con `#999999` fuera de tokens, `text-text-disabled` sobre `surface-hover` (4.20) y un botón de 24 px produjo FALLA en las tres comprobaciones y salida 1; se borró.
- Criterios que toca (evidencia de componentes; la evaluación por pantalla es de U5): **U5-04** (contraste medido en el DOM sin pares bajo umbral) cumplido en los 8 componentes salvo los hallazgos FV-03/FV-04 (colores fuera de tokens), **U5-08** (0 hosts externos) 0 en Storybook, **U5-06** (estados y tamaños) cubierto salvo FV-01, FV-02 y FV-05.

### Hallazgos que se reproducen (corrección pedida en solicitudes_fuera_de_alcance)

| Id | Componente | Qué pasa (DESIGN.md) | Mediciones |
|---|---|---|---|
| FV-01 | Botón | El anillo de foco aparece con transición: `transition-colors` incluye `outline-color` (§7.0 y §4.5: sin transición) | 128 |
| FV-02 | Gráfica | «Ver como tabla» mide 20 px de alto (§4.6/§9.3: 44 px con puntero grueso; WCAG 2.2: 24 px) | 24 |
| FV-03 | Gráfica | La descripción emergente de ECharts usa fondo #FFFFFF, texto #6D6E73 y borde #B7B9BE también en oscuro (§7.7: `ui.surface.raised`, `ui.text.default`) | 384 |
| FV-04 | Gráfica | Al pasar el puntero, ECharts recolorea símbolos y leyenda con variantes que no están en C2 (§7.7: las demás series no se atenúan) | 52 |
| FV-05 | Campo, Chip | Deshabilitados sin el estilo de §7.0 (relleno `ui.surface.disabled`, texto `ui.text.disabled`, borde punteado) | 16 |

## Supuestos (marcados PROVISIONAL cuando se decidió algo)

- PROVISIONAL: se eligió Storybook 10 y no la página `/_componentes` (pesa poco: 86 paquetes, build en ~1 s). Si Axel prefiere la página interna, las historias son CSF3 y se reutilizan.
- Para dejar corriendo Storybook y el script fue necesario tocar `frontend/package.json` (scripts `storybook`, `build-storybook`, `typecheck:variantes`, `verificar:variantes`; devDependencies exactas `storybook` y `@storybook/react-vite` 10.6.1) y `frontend/.gitignore` (`storybook-static/`). Es lo mínimo; el integrador debe conservarlo.
- `.storybook/tsconfig.json` es aparte porque `frontend/tsconfig.json` no incluye `.storybook/`; las historias sí entran en `tsc --noEmit` de la raíz porque están en `src/`.
- Los estados hover, foco y presionado no se simulan con clases: se provocan con el puntero y el teclado reales de Chromium (los estilos `hover:` de Tailwind solo valen con `@media (hover: hover)`, que el contexto móvil no tiene, igual que Android).
- Umbral de áreas táctiles con puntero fino: 24 px (mínimo de WCAG 2.2) porque §4.6 baja solo el «dibujo» a 36 px y §7.4 fija el chip en 24 px visuales; con puntero grueso 44 px, como pide §9.3.
- Se considera decorativo (exento) el color igual a `ui.border.divider`, `ui.viz.grid`, `ui.viz.track` y los esqueletos (§7.1, §10.5 «decorativo»).
- Los datos de las historias son literales sintéticos (centavos enteros formateados con `dinero()`); no se usó `datos-sinteticos.json`.
- No se pudo comprobar `npm run build` de la app completa: falla por archivos de otros agentes en `src/pantallas/` (ver Bloqueos). El `dist/` que existe es el de frontend-base; mis archivos no entran al build de la app (Vite solo empaqueta lo importado desde `main.tsx`).

## Bloqueos

- Ninguno nuevo (no se creó archivo en `docs/estado/bloqueos/`).
- Externo: `tsc` de la raíz y `npm run build` fallan hoy por `src/pantallas/gastos-presupuesto/` (módulos `P02Gastos`, `P04Gasto`, … aún sin escribir), `ajustes-a/` y `creditos-salud-reportes/`. Son de los agentes de área en curso, no de este alcance.

## Pendientes

- Componentes de DESIGN.md §7 que **no existen** en `frontend/src/componentes/` y por eso no tienen historias (frontend-base debe crearlos o las áreas decidir; ver solicitudes): botón de **acción crítica** y estado **ocupado** (§7.2), prefijo/sufijo y **monto grande** de Campo, **casilla, interruptor y opción segmentada** (§7.3), **avatar**, **insignia de rol** y **chip de filtro quitable** (§7.4), **barra de avance** con umbrales (§7.6), estados de la tarjeta de métrica (carga, error, «Falta capturar») y el **presionado** con borde `ui.border.selected`, gráficas de **dona** y **mini serie** (§7.7, §8), **calendario de pagos** (§7.8), **diálogo de acción crítica** (§7.12) y **panel de Ajustes** (§7.13).
- Historias de `src/vistas-bd/` (tabla, tablero, calendario, galería) y de `src/shell/` (barras, menú «/»): no estaban en mi alcance (`src/componentes/**`); se pueden agregar con el mismo `apoyo.tsx` y el script las mide solas.
- Mientras FV-01..FV-05 sigan abiertos el script los informa como HALLAZGO; al corregirse, `npm run verificar:variantes` los marca CERRADO y se pueden borrar de `hallazgos.json`.
- Campo pinta la etiqueta con `text-text-secondary` y §7.3 pide `ui.text.default` (observación visual, no es falla de contraste: pasa 4.5).

## Propuestas para Axel (PROPUESTA, no decidido)

- Dejar `npm run verificar:variantes` como paso de `make verify-wf3-*` (tarda ~70 s y no necesita red) para que un color fuera de C2 o un par bajo umbral rompa la verificación.
- Corregir FV-01..FV-05 en una sola pasada de frontend-base (son 4 archivos y pocas líneas).

## Paquetes instalados

- npm (devDependencies, versión exacta): `storybook@10.6.1`, `@storybook/react-vite@10.6.1` (+ 84 dependencias transitivas; `npm audit` sin vulnerabilidades). Sin `sudo` ni mkcert.

## Consumo aproximado

- Acciones: ~40 · KB leídos: ~150 (CONVENCIONES-frontend, DESIGN §7, contraste_componentes.py, componentes base, tokens) · KB escritos: ~65 (historias 24, configuración 9, script 30, reporte 7). Temporales de `datos/tmp/variantes/` y `storybook-static/` borrados al terminar.
