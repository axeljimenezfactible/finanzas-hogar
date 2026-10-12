# revisor-ds-r1 · Pista 3 fase 1 (sistema de diseño) · ronda 1 de 3 · 11-oct-2026, noche

Agente `revisor-ds-r1` (Opus), revisor independiente. No corregí nada de lo revisado. Escribí solo este reporte y `docs/estado/wfp3/REVISION-ds-r1.md`. Sin red, sin navegador, sin Artifact, sin git ni Telegram.

## Hecho
- Revisé `docs/design/sistema/project/` completo contra DESIGN.md 1.0, C2-tokens.json 1.0, P3-claude-design.md 1.3, format.md, cover.md y craft.md. Dictamen: **NO APROBADO**, con 1 alto (A-01, contraste del hover del chip, heredado de DESIGN §7.4), 3 medios (M-01 faltan §7.10 y §7.12; M-02 peso 500 sin valor y el verificador dice 18 → 18; M-03 supuestos provisionales sin marca en los README) y 5 bajos. El informe está en `docs/estado/wfp3/REVISION-ds-r1.md`.
- Hice mis propios scripts en `datos/tmp/rev-ds-r1/` (`tok.py`: conteo de C2, alias, hex, nombres, fuentes; `prev.py`: marcador, HTML, red, variables, accesibilidad y contraste por cascada CSS mínima en los dos temas; `misc.py`: nombres del README, emoji, formato de dinero, hover del chip, foco, portada). Los borré al final.

## Verificación (comandos y salidas)
```
$ python3 -I docs/design/tools/ds/verificar_ds.py
RESULTADO: OK · 445 tokens C2 -> 351 en listas de tokens.json + 6 estilos de tipo (18 font.*) + 2 fuentes; 0 fallas
$ python3 -I datos/tmp/rev-ds-r1/tok.py
C2 total 445 (color 190, ui 150, size 32, font 18, layout 14, space 8, viz 8, motion 5, breakpoint 5, radius 4, shadow 3, icon 3, border 2, overlay 2, focus 1)
themes order ['light', 'dark'] · aliases 222 bad/nonhex 0
color compares 532 · missing in css [] · mismatches 2 (ui-overlay-scrim claro/oscuro: color-mix -> #22232566/#22232599, conversión declarada)
20 muestras al azar (semilla 20261011): claro y oscuro iguales a tokens.css
hex C2 154 · hex tokens.json 154 · diferencias [] []
names 351 dups [] bad [] · var( 0 color-mix( 0 url( 0 · familias de lista <= 60
css vars not in tokens.json: las 19 font-* / line-height-* (ver M-02)
font InterVariable.woff2 wOF2, sha256 = docs/design/fuentes/inter/web/ (True); Italic igual (True)
$ python3 -I datos/tmp/rev-ds-r1/prev.py   (14 previews)
marcador válido 14/14 · HTML mal formado 0 · etiquetas prohibidas 0 · URL externas 0 · variables inexistentes 0 · hex literales 0
botones de solo icono sin aria-label 0 · labels con for roto 0 · pares calculados 652
texto mínimo: claro 4.51 (MenuBarra), oscuro 5.50 · UI mínimo 3.06
fallos: Boton :hover deshabilitado (falso positivo, la especificidad lo anula); Grafica anillo y retícula (decorativos); EstadosVista esqueleto (decorativo)
$ python3 -I datos/tmp/rev-ds-r1/misc.py
README var() no existentes [] · emoji 0 · dinero sin formato: solo ejes compactos ($6, $5 mil)
chip-hover light ui-person-karen 2.83 · light ui-store-sams-club 2.74 · dark ui-person-axel 2.14 · dark ui-store-market 2.11   <- A-01
foco: claro 3.19–3.82, oscuro 4.16–7.89 (>= 3 en las 6 superficies)
portada, bloque contra fondo: claro 1.45–5.53, oscuro 2.73–10.38
$ grep -rnoE "<21 cifras canon de P3-12>" docs/design/sistema/project | wc -l
0
$ grep -rn -i -E "<patrón de credenciales de P3-12>" docs/design/sistema/project | wc -l
0
fechas de project/ ∩ oraculo/model.json: {2026-10-30: tokens.json}   <- B-01 (sin monto)
```
- Contrasté 21 afirmaciones del README con DESIGN (lista en el informe): 0 contradicciones; la tabla «Contraste medido» la recalculé con los hex de tokens.json y da lo mismo.

## Supuestos
- Clasifiqué como decorativos, y fuera del umbral de 3:1, el anillo de `ui-surface-card` de los marcadores (DESIGN §8.1), la retícula `ui-viz-grid` y el esqueleto `ui-skeleton-on-*`, porque sus notas de token dicen «decorativo». También los divisores `ui-border-divider` (README: «decorativo»).
- Texto grande (umbral 3:1) solo con `title-1` y `title-2` (24 px o más); `title-3` (20 px, 600) lo medí a 4.5:1.
- A-01 va como **alto** por la regla de severidad del encargo («contraste bajo umbral»), aunque la vista previa sigue a DESIGN §7.4 al pie de la letra y el agente lo escaló bien. La parte de color la resuelve Axel; la parte de la nota en el README se puede corregir ya.
- No validé los campos de `design-system.json` más allá de los que pide el encargo: su gramática no está en `referencia-tipo/`.
- No rendericé nada (sin navegador): la altura de cada tarjeta y el aspecto real de la portada quedan para cuando el orquestador lea el artifact publicado.

## Bloqueos
- Ninguno nuevo. Siguen abiertos `docs/estado/bloqueos/P3-contraste-chip-hover.md` (base de A-01) y `P3-huecos-componentes.md` (base de M-03).

## Pendientes (para la ronda 2)
- A-01: nota en el README de ChipPersonaTienda ya; el cambio de color, cuando Axel responda el bloqueo.
- M-01: componentes de §7.10 y §7.12, o declararlos pendientes.
- M-02: familia de pesos (o declarar que el 500 solo va en `usage`) y corregir la línea `font 18 -> 18` del verificador.
- M-03: marcar los supuestos provisionales en cada README.
- B-01 a B-05: opcionales (fecha 30-oct, ubicación de `design-system.json` al publicar, `usage` con fondos, `rx` literal de la portada, ruta del script y licencia OFL).

## Consumo aproximado
- Leído ≈ 175 KB: format.md, cover.md y craft.md ≈ 31 KB; criterios P3 por grep ≈ 12 KB; reportes ds-tokens y ds-brandbook ≈ 15 KB; README 23 KB; previews de Boton y Cover ≈ 8 KB (el resto por script); bloqueos ≈ 3 KB; DESIGN.md por secciones y grep ≈ 35 KB; salidas de scripts ≈ 45 KB.
- Escrito ≈ 40 KB: scripts temporales ≈ 22 KB (borrados), informe 13 KB, este reporte ≈ 5 KB.
- ≈ 32 acciones de herramienta. No cerré por tope (`cerrado_por_tope=false`).
