# correccion-ds-r2 · Pista 3 fase 1 (sistema de diseño) · 11-oct-2026, noche

Agente `correccion-ds-r2`. Corrige el hallazgo alto y el medio de `REVISION-ds-r2.md` (A-01 y M-01). Los 9 bajos no estaban en el encargo y no se tocaron.

- **Alcance:** `docs/design/sistema/project/`, `docs/design/tools/ds/` (sin cambios), `docs/estado/bloqueos/P3-huecos-componentes.md` y este reporte.
- **Restricciones:** sin git, sin red, sin navegador, sin Artifact.
- **Respaldo previo:** `/tmp/correccion-ds-r2-respaldo/` (project/, ds/ y los dos bloqueos P3).

## Hecho
1. **A-01 · iconos de 24 px en la barra superior.**
   - `components/BarraSuperior/preview.html`: los 3 botones de solo icono pasan de `i20` a `i24`, con `var(--size-icon-md)` y `stroke-width: var(--icon-stroke-user-units-md)`. Son «Más acciones» en móvil, «Más acciones» en escritorio y «Buscar». Las migas y el icono de «+ Registrar» se quedan en `i20` (sm), junto a texto `small`.
   - `components/BarraSuperior/README.md`: agrega la regla (botones de solo icono en `var(--size-icon-md)`; migas y «+ Registrar» en `var(--size-icon-sm)`, DESIGN §5.2 y §7.2) y suma a sus Tokens `var(--size-icon-md)`, `var(--icon-stroke-user-units-md)`, `var(--size-icon-sm)` y `var(--icon-stroke-user-units-sm)`.
   - `README.md`, Iconografía, dice ahora: «`var(--size-icon-md)` … junto a títulos de 20 px o más y en los botones de solo icono de la barra superior y la barra inferior».
2. **M-01 · inventario de §7.5.** Dos componentes nuevos, cada uno con README y preview. Los dos llevan marcador en la línea 1, `lang="es-MX"` y solo `var(--token)` existentes, con iconos copiados del sprite de Lucide 1.55.0.
   - `components/BarraVista` (§7.5.1, grupo «Vistas de base de datos»):
     - en escritorio, selector Tabla · Tablero · Calendario · Galería con `table-2`, `kanban`, `calendar` y `layout-grid`, como alternado con `aria-pressed` (relleno `ui-surface-selected` e icono `ui-icon-selected`);
     - «Filtrar» (`list-filter`), «Ordenar» (`arrow-up-down`) y «Agrupar» como botones de texto de §7.2;
     - búsqueda por texto (los gastos tienen el filtro «texto», U1 §2.5.b);
     - filtros activos como chips quitables de §7.4 («Quitar filtro Ximena», «Quitar filtro Walmart»);
     - en móvil, un solo «Filtrar» con `aria-haspopup="dialog"`, que abre la hoja de `HojaDialogo`, y los chips bajo el título «Gastos», con desplazamiento horizontal propio (U1 §1.3).
   - `components/Galeria` (§7.5.5):
     - tarjetas `ui-surface-card`, `radius-md`, `space-3` y `shadow-0`, con portada 4:3 en `ui-surface-hover` y el icono de la entidad en `ui-icon-default`;
     - hover con solo borde `ui-border-control-hover`; elegida con `ui-surface-selected` y borde `ui-border-selected`;
     - error de imagen con icono y «Foto no disponible»;
     - cifra de gasto con «−», `arrow-up-right` y `tabular-nums`, y tienda como chip;
     - ejemplo en 3 columnas (como a 768) y tickets en 2 (como a 360);
     - el README trae la tabla de columnas por ancho de U1 §6.2 y el texto alternativo con comercio, fecha y total.
   - `components/CalendarioPagos/README.md`: nueva línea «Relación con el calendario genérico (DESIGN §7.5.4)». Dice que representa también la base de §7.5.4 (DESIGN §7.8: «Parte del calendario de §7.5.4»), lista lo que sí dibuja y lo que no (navegación «‹ Hoy ›», «+N más», días de otro mes en `ui-text-disabled`, «+» en hover, renglón genérico y agenda de 360), y remite al punto 12 del bloqueo.
   - `P3-huecos-componentes.md`: nota de actualización y puntos 10 (BarraVista), 11 (Galeria) y 12 (calendario genérico pendiente), cada uno con su pregunta a Axel.
   - **Corrección de lo que dijo `correccion-ds-r1.md`:** su frase «quedan cubiertas todas las filas de la matriz de DESIGN §7.16» no era exacta. La fila «Tabla, tablero, calendario, galería» queda ahora cubierta por `TablaBD`, `Tablero`, `CalendarioPagos` (con lo pendiente del genérico en el punto 12), `Galeria` y `BarraVista`.
3. `design-system.json`: `lastChange.at` 2026-10-12T05:40:25Z, con la nota de esta corrección. Las demás claves no cambian.

## Verificación (comandos y salidas)
```
$ python3 -I docs/design/tools/ds/verificar_ds.py | tail -1
RESULTADO: OK · 445 tokens C2 -> 355 en listas de tokens.json + 6 estilos de tipo (18 font.*) + 2 fuentes; 0 fallas
$ python3 -I docs/design/sistema/verificar_brandbook.py   (líneas relevantes)
components/BarraSuperior/preview.html: pares texto 10, gráficos 9 · grafico dark mín 5.50 · grafico light mín 3.71 · texto dark mín 7.42 · texto light mín 5.40
components/BarraVista/preview.html: pares texto 10, gráficos 16 · grafico dark mín 3.28 · grafico light mín 3.73 · texto dark mín 5.04 · texto light mín 5.02
components/Galeria/preview.html: pares texto 9, gráficos 8 · grafico dark mín 3.06 · grafico light mín 3.28 · texto dark mín 5.50 · texto light mín 4.51
previews 21 · variables distintas usadas 161 (todas en tokens.json: True) · cifras canon 0
RESULTADO: OK · 0 fallas · 1 avisos de excepción declarada      (el chip de Karen, ya conocido)
$ python3 -I docs/design/tools/ds/verificar_sinteticos.py | tail -1
RESULTADO: OK
```
- **Iconos contra `fuentes/lucide/sprite.svg`** (contenido del `<symbol>`, sin espacios):
  - BarraSuperior: 6 de 6 (chevron-left, ellipsis × 2, chevron-right, search, plus);
  - BarraVista: 12 de 12 (table-2, kanban, calendar, layout-grid, list-filter × 2, arrow-up-down, search, x × 4);
  - Galeria: 10 de 10 (scan-line × 5, arrow-up-right × 5);
  - 0 sin coincidencia.
- **BarraSuperior:** si se normaliza `i24` → `i20`, el preview queda idéntico al respaldo; solo cambiaron los 3 iconos.
- **Contraste propio con los hex de `tokens.json`** (claro / oscuro), para los pares que la cascada del verificador no recorre:
  - marca de Ximena sobre `ui-surface-hover` (chip quitable en hover): 4.34 / 4.78;
  - marca de Walmart sobre `ui-surface-hover`: 4.58 / 3.28;
  - `ui-text-secondary` sobre `ui-surface-selected` (vista elegida): 5.02 / 8.08;
  - `ui-icon-selected` sobre `ui-surface-selected`: 5.15 / 8.04;
  - «Foto no disponible» (`ui-text-secondary`) e icono sobre `ui-surface-hover`: 4.51 / 5.50;
  - `ui-role-expense` sobre `ui-surface-selected`: 5.32 / 5.97;
  - Sam's Club y Mercado sobre `ui-surface-raised`: 3.28 / 10.18 y 13.74 / 3.06.
  
  Todos cumplen. Los chips quitables usan Ximena y Walmart, que no están entre los 4 pares de la excepción del chip.
- **P3-12:** el grep de las 21 cifras canon del criterio sobre `docs/design/sistema/project` da **0**, y el patrón de credenciales da **0**.
  - Montos con centavos en los archivos tocados: $842.37, $1,163.90 y $275.48 (nuevos), más los que ya estaban. Ninguno aparece en `oraculo/model.json` ni en `datos/semilla/data.json`.
  - Fechas nuevas (13, 20 y 27-nov-2026): 0 coincidencias con `oraculo/model.json`.
  - El ejemplo de texto alternativo de DESIGN («$1,258.68») es cifra canon y **no** se copió: se usó «Ticket de Walmart, 13-nov-2026, $842.37».
- **`diff -rq` contra el respaldo:**
  - cambian solo `README.md`, `BarraSuperior/README.md`, `BarraSuperior/preview.html`, `CalendarioPagos/README.md` y `design-system.json`;
  - hay 2 carpetas nuevas, `BarraVista` y `Galeria`;
  - `tools/ds/` no tiene diferencias.

## Supuestos
Todos los que DESIGN y U1 no fijan van en sus README («Supuesto provisional») y en `P3-huecos-componentes.md`:
- **Punto 10 (BarraVista):**
  - «Agrupar» sin icono;
  - texto de la vista elegida en `ui-text-secondary`;
  - búsqueda con el estilo de `CampoTexto`;
  - orden y disposición de la fila, con los chips en una línea propia;
  - «Filtrar» a la izquierda de los chips en móvil.
- **Punto 11 (Galeria):**
  - `scan-line` en `size-icon-md` como icono de ticket;
  - portada con `radius-sm`;
  - «Foto no disponible» en `caption` `ui-text-secondary`;
  - `space-4` entre tarjetas;
  - fecha como segunda propiedad;
  - selección como en `Tablero`;
  - a 360, título, cifra y comercio.
- La vista previa de la galería no lleva foto: no hay una imagen sintética de ticket en el repositorio.
- Títulos, tiendas, montos y fechas de los ejemplos son sintéticos.

## Bloqueos
- `P3-huecos-componentes.md` (abierto, provisional; ahora con 12 puntos). No detiene la Pista 3.
- `P3-contraste-chip-hover.md` sigue abierto, sin cambios.

## Pendientes
- **Calendario genérico (§7.5.4), punto 12:** no tiene vista previa propia de la navegación «‹ Hoy ›», «+N más», los días de otro mes en `ui-text-disabled`, el «+» en hover, el renglón genérico ni la agenda de 360. Si Axel pide un componente `Calendario`, `ui-text-disabled` necesita una excepción declarada en `verificar_brandbook.py`.
- **BarraVista, sin dibujar:** qué abren «Filtrar», «Ordenar» y «Agrupar» en escritorio, el contenido de la hoja de filtros de móvil, el selector de vista en móvil y la barra de 768 a 1023. Son preguntas del punto 10.
- **Fuera de mi alcance:** la fila de `P3-huecos-componentes.md` en `docs/estado/bloqueos/INDICE.md` todavía dice «Seis huecos». Conviene que el integrador la actualice a 12 puntos.
- **Orquestador, al publicar:** revisar en los dos temas, con la comprobación de render, los previews de `BarraVista` y `Galeria` y la barra superior con iconos de 24 px. Aquí no hubo navegador. Comprobar que `height=600` de la galería y `height=330` de la barra de la vista bastan.
- **Bajos de la r2 sin atender** (no estaban en el encargo): B-01 (inferencia en VistaPrevia), B-02 (conteo del tope de familias), B-03 (foco de pestañas), B-04 (`layout-dialog-width-long` como ancho de marco) y los 5 bajos de la r1.

## Consumo aproximado
- **Leído:** ≈ 110 KB:
  - revisión y reporte r1, 29 KB;
  - DESIGN §5.2, §5.4, §7.2, §7.4, §7.5 y §7.8, ≈ 28 KB;
  - U1 §1.3 y §2.5, ≈ 8 KB;
  - previews y README de BarraSuperior, Tablero, ChipPersonaTienda, CampoTexto, TablaBD y CalendarioPagos, ≈ 25 KB;
  - salidas de verificadores, ≈ 20 KB.
- **Escrito:** ≈ 37 KB:
  - BarraVista, 12.3 KB;
  - Galeria, 10.3 KB;
  - ediciones, ≈ 3 KB;
  - bloqueo, +3.1 KB;
  - este reporte, ≈ 8 KB;
  - ayudante de iconos en `/tmp`, 0.4 KB.
- **Acciones de herramienta:** ≈ 40.
