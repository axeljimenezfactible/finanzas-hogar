# Revisor ajustes-a · vuelta 3 (agente `revisor-ajustes-a-r3`)

- Fecha: 12-oct-2026. Soy revisor independiente: no corrijo. El informe completo está en `docs/estado/wfp3/P3-iter-3-ajustes-a.md` (criterio P3-07).
- Alcance: los 16 artboards de la parte a en `docs/design/maquetas/lienzo/project/ajustes/` (índice y s01 a s07, `-m` y `-e`) y `docs/design/maquetas/lienzo/fragmentos/ajustes-a.json`.

## Hecho

- Comprobé los 2 medios de la vuelta 2 y **los dos quedan cerrados**:
  - **M-01**: los diálogos críticos de s02 y s05 ya llevan relleno `ui-action-critical-bg`, `shield-alert` y título en `ui-critical`, y superficie con radio y sombra.
  - **M-02**: los 5 `-m` tienen borde `ui-border-error` en el campo inválido.
- Un diff contra `docs/design/salida/iter-2/ajustes-a/` confirma que el corrector solo tocó el `<helmet>` y el sprite.
- **Veredicto: APROBADO. 0 altos, 0 medios, 8 bajos.** Seis bajos vienen de antes: B-01 a B-06. Dos son nuevos:
  - B-07: `.ro` sin regla en s06.
  - B-08: la anatomía de los diálogos de ejemplo. En móvil no son hoja inferior; en escritorio no tienen ancho de 480 px; el cuerpo no va en texto secundario; y s02 tiene 2 `autofocus`.

## Verificación

Hice scripts propios con `python3 -I` en `datos/tmp/rev-aj-a-r3/` y los borré al final. Usé un árbol `html.parser` propio y una cascada CSS propia: selectores, especificidad, herencia y fondo efectivo. Comprueban:

- la forma y `data-props`/`$preview` contra el fragmento: 16/16;
- el CSS bien formado; 0 etiquetas sin cerrar; 0 ids repetidos; 0 referencias ARIA sin destino;
- 0 hex; 0 `var()` inexistentes;
- las clases sin regla que case: solo s04, s06 y s07-m, todas bajas;
- `<use>` sin símbolo: 0;
- el contraste: 27 pares de texto, 38 de interfaz y 14 reglas de estado, en claro y oscuro, con 0 bajo umbral. Recalculé 10 pares y coinciden con `contraste-componentes.csv` (±0.01);
- P3-12: 0 cifras canon, 0 centavos canon y 0 credenciales. 60 de 66 montos y 24 de 25 fechas están en el JSON; el resto son derivados que comprobé a mano. En pantalla, las fechas que coinciden con el oráculo van con montos sintéticos;
- los objetivos táctiles en `-m`: todos ≥ 44 px;
- la raíz: 412 fija en `-m`; fluida con `flex-wrap` y `expand: "fill"` en `-e`;
- los enlaces: 0 rotos;
- la columna de las 14 secciones: en orden y con `aria-current`;
- el idioma, U+2212, y que no hay ni emoji ni ISO.

## Supuestos

1. Sin navegador (K45), el render se infiere del CSS del `<helmet>`, de los tokens y de `tokens.css` (que solo define `body`, `caption`, `small` y `title-1..3`). Las reglas con `@media` se evalúan al ancho del artboard.
2. «Error», «No», «Total» y «Plan» se leen como español, porque son textos literales de U4. «Liverpool» y «Banorte» son nombres propios mexicanos que vienen del JSON.
3. Los diálogos se dibujan en línea como ejemplo («Así se ve el diálogo al confirmar»). Por eso su anatomía (hoja o caja) cuenta como estilo (B-08) y no como elemento ausente.
4. El choque `groceries`/`eating_out` entre U4 y U5-07 sigue en manos del orquestador (lo heredo de la vuelta 2).

## Bloqueos

Ninguno.

## Pendientes

- Bajos opcionales B-01 a B-08 (correcciones en el informe). No hay vuelta 4 (P3-06).
- B-05: medir las alturas al publicar.

## Consumo aproximado

Unos 190 KB leídos, unos 20 KB escritos (los dos informes) y unas 80 acciones. Pasé el tope de 60 acciones por partir la verificación en pasos chicos. La revisión está completa (`cerrado_por_tope = false`).
