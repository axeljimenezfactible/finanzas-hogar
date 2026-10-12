# ds-brandbook · Pista 3 fase 1 (sistema de diseño) · 11-oct-2026, noche

Agente `ds-brandbook` (Opus). Reanudación tras el corte del flujo wf_6b8245eb-3e3: empecé de cero. No toqué `tokens.json` ni `fonts/`. Sin red, sin navegador, sin Artifact, sin git.

## Hecho
- `docs/design/sistema/project/README.md` (22.9 KB): brand book en es-MX, en imperativo y nombrando `var(--…)` de `tokens.json`. Sin título, sin procedencia. Cubre: 8 principios, tono (de tú, sin emoji ni exclamaciones, textos fijados de U1), color (60-30-10, superficies/texto/bordes con sus fondos, acciones, 8 roles con refuerzo, acción crítica, personas y tiendas con forma y regla K41, series), contraste medido (de DESIGN §2.8), tipografía (6 estilos, pesos), dinero (tabulares, U+2212 antes de `$`, rebasado), espaciado y rejilla, radios/sombras/bordes, iconografía (Lucide 1.55.0, 1.5 px, 20/24), movimiento y `prefers-reduced-motion`, estados, dataviz (§8), accesibilidad, puntos de quiebre 360/412/768/1024/1440 y «Datos sintéticos» (K45.1).
- 13 componentes, cada uno con `README.md` (primera frase = resumen; cuándo, qué provee quien lo usa, hacer/no hacer, tokens) y `preview.html` (marcador en la línea 1, HTML + CSS con `var(--token)`, sin React, sin red, sin iframes, controles reales, `aria-label` en iconos sueltos): `Boton`, `CampoTexto`, `TarjetaMetrica`, `TablaBD`, `Tablero`, `CalendarioPagos`, `BarraLateral`, `MenuBarra`, `ChipPersonaTienda`, `AvisoBreve`, `EstadosVista`, `InterruptorOpciones` (interruptor encendido/apagado y lista con bloque «Antes → Después», U4 §0.2) y `Grafica` (SVG en línea: 3 líneas con marcador, trazo y etiqueta directa; leyenda y las 8 series con color, marcador, trazo y patrón). Grupos: Acciones, Formularios, Datos, Vistas de base de datos, Navegación, Avisos, Gráficas.
- Iconos: copiados del sprite `docs/design/fuentes/lucide/sprite.svg` (60 símbolos, en línea), nada dibujado a mano.
- `components/Cover/preview.html` (carpeta sola): bloques `color-brand-700`, `color-brand-300`, `color-neutral-700` y un `color-critical-700` pequeño; patrón de teselas cortadas por los radios (fila «geométrico, modular»: vistas de BD, radios 6/10/14, retícula de 4 pt), 23 teselas; «Finanzas Hogar» a 88 px en dos líneas y el lema «Las cuentas de la casa, claras y en calma.»; derivación de 4 líneas en comentario; mismo archivo en los dos temas (bloques primitivos visibles sobre `ui-surface-page` claro y oscuro).
- `project/design-system.json`: índice con los campos del encargo (`v` 3, `layout` files, `createdOnFiles` y `lastChange` con la hora UTC de `date -u`).
- `docs/design/sistema/verificar_brandbook.py`: verificador propio (fuera de `project/`, no se publica).
- Bloqueos nuevos: `docs/estado/bloqueos/P3-contraste-chip-hover.md` y `P3-huecos-componentes.md`.

## Verificación
```
$ python3 -I docs/design/sistema/verificar_brandbook.py
(14 previews; pares texto 2–14 y gráficos 0–31 por archivo; mínimos: texto claro 4.51, texto oscuro 5.50, gráfico claro 3.15, gráfico oscuro 3.06)
previews 14 · variables distintas usadas 145 (todas en tokens.json: True) · cifras canon 0
components/ChipPersonaTienda/preview.html: grafico ui-person-karen sobre ui-surface-hover light 2.83:1 < 3.0 · excepción declarada …
RESULTADO: OK · 0 fallas · 1 avisos de excepción declarada
```
- Qué comprueba: marcador `@dsCard` y `height` en rango (portada 240–360); etiquetas cerradas (`html.parser`); sin `iframe/object/embed/noscript` ni `http(s)://`; sin «Inter» escrito como pila; todo `var(--x)` de README, READMEs y previews existe en `tokens.json`; contraste WCAG con los hex de `tokens.json` en claro y oscuro de cada texto (≥ 4.5) y de cada marca, icono, borde de control o relleno con significado (≥ 3) según una cascada CSS mínima, incluidos `:hover`, `:active`, `:focus-within` y `::placeholder`; sin emoji; sin `$` de 4 cifras sin coma ni guion corto antes de `$`; cifras canon P3-12; carpetas completas y `Cover` sola; README sin título; claves del índice.
- Prueba negativa en una copia (`/tmp`): inyecté un color de bajo contraste, un `<div>` sin cerrar, un `var(--ui-no-existe)`, `-$1234.50` y quité un marcador → detectó las 9 fallas.
- Dos defectos reales que encontró y corregí: el `plus` del botón de `EstadosVista` heredaba `ui-icon-default` sobre el primario (1.03:1) → `.vacio > .icon`; `label{font-weight:500}` perdía contra la clase `.small` → `.grupo label`.
- Comando del criterio P3-12 sobre `docs/design/sistema/project`: 0 cifras canon; patrón de credenciales: 0; URLs en previews: 0; «Kueski»: 0.

## Supuestos
- Nombres de componente en español sin acentos (identificadores ASCII del formato). `design-system.json` sigue los campos del encargo: el formato del tipo lista ese nombre como ruta de tabla y no trae su gramática en `referencia-tipo/`.
- Las clases tipográficas (`.title-2`, `.small`, `.caption`…) y `--font-sans` son las que el tipo compila desde `tokens.json`; los pesos 500/600/700 van como número porque no hay variable de peso.
- Todo lo de `P3-huecos-componentes.md` (iconos de Gastos, Presupuesto, Salud y Reportes; aspecto del bloque «Antes → Después»; perilla de 16 px; círculo de hoy de 24 px; comandos de ejemplo del menú «/»).
- Las vistas previas muestran la variante de un ancho (móvil primero con puntero grueso y la medida fina por `@media (pointer: fine)`); los cambios por punto de quiebre van en los README.
- Cifras, fechas, nombres de líneas de crédito y montos son inventados (noviembre de 2026, días que no caen en 9, 10, 15, 28 ni 30).

## Bloqueos
- `P3-contraste-chip-hover.md`: DESIGN §7.4 deja la marca de Karen y Sam's Club (claro) y de Axel y Mercado (oscuro) bajo 3:1 sobre `ui-surface-hover`; se respetó DESIGN y se declaró la excepción. Pregunta a Axel/U2 con recomendación (hover solo con borde).
- `P3-huecos-componentes.md`: supuestos PROVISIONALES de lo que DESIGN no fija.

## Pendientes
- Orquestador: publicar `project/` y revisar en la página que las clases tipográficas y `--duration-*`/`--ease-*` se compilan (si no, las vistas previas pierden tamaño o transición, no color).
- Grupo de iconos en `assets/` y la licencia OFL de Inter (pendiente de ds-tokens): no eran de este alcance.
- Si Axel resuelve el bloqueo del chip, cambiar `button.chip:hover` en `ChipPersonaTienda` y quitar la excepción del verificador.

## Consumo aproximado
- Leído ≈ 190 KB (DESIGN §1–§9 y Supuestos ≈ 110 KB, format/craft/cover ≈ 31 KB, reporte ds-tokens 9 KB, demo.json por script 7 KB, U1/U4/plan/criterios ≈ 18 KB, salidas de scripts ≈ 15 KB).
- Escrito ≈ 143 KB (README 23 KB, previews ≈ 76 KB con ≈ 12 KB de iconos expandidos por script, READMEs 18 KB, verificador 13 KB, bloqueos e índice 4 KB, este reporte 6 KB).
- ≈ 35 acciones de herramienta.
