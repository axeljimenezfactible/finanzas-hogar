# Revisor ajustes-a · vuelta 2 (agente `revisor-ajustes-a-r2`)

- Fecha: 12-oct-2026. Revisor independiente: no corrige. Informe completo en `docs/estado/wfp3/P3-iter-2-ajustes-a.md` (criterio P3-07).
- Alcance: los 16 artboards de la parte a de `docs/design/maquetas/lienzo/project/ajustes/` (índice y s01 a s07, `-m` y `-e`) y `docs/design/maquetas/lienzo/fragmentos/ajustes-a.json`.

## Hecho

- Comprobé los 8 medios de la vuelta 1. Se cierran 6: enlaces a s08, `#instrucciones`, objetivos táctiles, estados, fechas sin cero y barra lateral §7.9. Quedan 2 a medias: las validaciones (M-05) y los elementos de U4 (M-08), por reglas CSS y un símbolo que faltan en los archivos tocados.
- **Veredicto: NO aprobado. 0 altos, 2 medios, 6 bajos.**
  - M-01: los diálogos críticos de s02 y s05 (`-m` y `-e`) no tienen las reglas `.btn.c` ni `.dlg.cr h3`, y les falta `<symbol id="ic-shield-alert">`. En s05 falta también `.dlg`. Por eso no se ven como acción crítica: sin relleno crítico, sin icono y, en s05, sin caja de diálogo.
  - M-02: en s01-m, s03-m, s04-m, s05-m y s07-m falta `.inw.er`, así que el campo inválido no tiene borde de error.
  - Bajos: `.btn.t`, `.row`, `.tx` y `.cn` sin regla en s04 y s07-m; concordancia en s05-e; leyenda «Enter…»; atajos incompletos; alturas estimadas; primaria no única.

## Verificación

Hice scripts propios con `python3 -I` en `datos/tmp/rev-aj-a-r2/` y los borré al final. Comprueban:

- la forma y `data-props`/`$preview` contra el fragmento: 16/16;
- `html.parser`: 0 errores;
- 0 hex;
- todos los `var(--x)` existen;
- **las clases sin regla y los `<use>` sin símbolo**, que es la prueba nueva de esta vuelta y la que encontró M-01, M-02 y B-01;
- el contraste en matriz claro y oscuro: 450 combinaciones y 0 pares reales bajo umbral;
- P3-12: 0 cifras canon y 0 credenciales; 55 de 60 montos y 24 de 25 fechas están en el JSON, y el resto son sumas declaradas;
- las alturas táctiles calculadas en `-m`: todas ≥ 44 px;
- los enlaces: 0 rotos;
- la columna interna: 14 secciones en orden;
- el idioma, U+2212 y que no hay ni emoji ni fechas ISO.

## Supuestos

1. Sin navegador (K45), el render se infiere del CSS del `<helmet>`, de los tokens y de `tokens.css`. Que una clase no tenga regla se interpreta como «no se dibuja», porque `tokens.css` solo define `body`, `caption`, `small` y `title-1..3`.
2. El contraste se evalúa como matriz de tokens usados × superficies usadas. Las combinaciones bajo umbral se descartaron a mano tras comprobar que no aparecen juntas en el marcado.
3. El anillo `ui-person-axel` del avatar sobre `ui-surface-hover` es la anatomía de DESIGN §7.4 y no cuenta como hallazgo.
4. Las claves `groceries` y `eating_out` siguen como choque entre U4 y U5-07 y quedan para el orquestador.

## Bloqueos

Ninguno.

## Pendientes

- Vuelta 3: corregir M-01 y M-02 copiando las reglas y el símbolo de s01-e y s06-m (ver la tabla por artboard del informe), y los bajos que se puedan.
- B-05: medir las alturas al publicar.

## Consumo aproximado

Unos 210 KB leídos, unos 22 KB escritos (los dos informes), 40 acciones.
