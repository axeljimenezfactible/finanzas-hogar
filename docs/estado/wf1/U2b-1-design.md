# U2b-1-design · reporte final (WF-1, 11-oct-2026)

## Hecho
- `docs/design/DESIGN.md` nuevo: secciones 1 a 6 completas (principios y tono; colorimetría resumida fiel a la 1.2 con tokens por nombre; tipografía; espaciado, rejilla, radios, sombras y bordes; iconografía; movimiento), encabezados 7 a 10 vacíos marcados «(la escribe U2b-2)» y sección Supuestos. Tablas con valores en px, rem y ms y nombres de token propuestos (`font.*`, `space.*`, `radius.*`, `shadow.*`, `size.*`, `layout.*`, `motion.*`, `icon.*`, `border.*`, `focus.*`, `overlay.*`) para que la sección 10 los lleve a tokens.
- Regla K41 escrita en §2.6 y §2.10 («una persona y una tienda nunca se distinguen solo por color en una misma vista»); K42 (Walmart oscuro `#7AA3EE`) y K26 (gasto 45°) documentados sin marcas provisionales; ingreso con sus dos valores (150° de decisiones §I y 135° vigente) y el motivo.
- Inter 4.1 variable (OFL 1.1) y Lucide 1.55.0 (`lucide-static`, ISC + MIT de Feather) descargados de sus orígenes oficiales (GitHub rsms/inter y registro npm), cada archivo comprimido en una carpeta nueva de `datos/tmp/`, extraídos los archivos útiles a `docs/design/fuentes/` sin modificar, con licencias; temporales borrados (solo los míos; `datos/tmp/s3chk.py` es de otro agente y no se tocó).
- Script de verificación `docs/design/fuentes/scripts/verificar.py` (solo biblioteca estándar): firmas y ejes de las fuentes, rasgo `tnum`, glifo U+2212, licencias, sprite, nombres de iconos de §5.4, hex y rutas de token de DESIGN.md contra `C2-colores.json`, 19 contrastes de §2.8 recalculados, URLs de terceros y encabezados.
- `docs/design/fuentes/FUENTES.md`: procedencia, sha256 de archivos comprimidos y conservados, integridad npm.

## Archivos
- `docs/design/DESIGN.md` (nuevo, 48.7 KB, escrito en 4 partes)
- `docs/design/fuentes/FUENTES.md` (nuevo)
- `docs/design/fuentes/scripts/verificar.py` (nuevo)
- `docs/design/fuentes/inter/web/InterVariable.woff2`, `docs/design/fuentes/inter/web/InterVariable-Italic.woff2`, `docs/design/fuentes/inter/ttf/InterVariable.ttf`, `docs/design/fuentes/inter/ttf/InterVariable-Italic.ttf`, `docs/design/fuentes/inter/LICENSE.txt` (descargados, sin modificar)
- `docs/design/fuentes/lucide/sprite.svg`, `docs/design/fuentes/lucide/LICENSE`, `docs/design/fuentes/lucide/package-1.55.0.json` (descargados, sin modificar)
- `docs/estado/wf1/U2b-1-design.md` (este reporte)

## Verificación
- `python3 -I docs/design/fuentes/scripts/verificar.py .` → código 0, `RESULTADO: OK`:
  - Fuentes: `InterVariable.woff2` 352,240 B y `InterVariable-Italic.woff2` 387,976 B con firma wOF2; TTF 879,708 B y 910,252 B con ejes `opsz` 14–32 y `wght` 100–900, 9 instancias, rasgos GSUB con `tnum`, `case`, `zero`, `cv01`–`cv14`, `ss01`–`ss08`; glifos `+ − $ , .` presentes.
  - Licencias: Inter «SIL Open Font License, Version 1.1» presente; Lucide «ISC License» y MIT de Feather presentes; cabecera `@license lucide-static v1.55.0 - ISC` en el sprite.
  - Iconos: 1,877 símbolos, todos `viewBox 0 0 24 24`, 0 `stroke-width` fijos (heredan el trazo 1.5 del CSS); 55 nombres citados en §5.4, 0 faltan.
  - DESIGN.md: 190 tokens en C2; 80 hex citados (74 distintos), **0 fuera de C2**; 209 rutas `color.*`, 0 inexistentes; 76 pares token → hex, 0 discrepancias; 19 contrastes de §2.8 recalculados, 19 coinciden (truncado a 2 decimales).
- Integridad npm del paquete de Lucide recalculada = la publicada por el registro (`sha512-h18lun…m/A==`).
- `grep -rn "fonts.googleapis\|fonts.gstatic\|cdn\.\|unpkg\|jsdelivr" docs/design docs/contratos --include='*.md' --include='*.css' --include='*.json' | wc -l` → 0.
- `grep -n -i 'linear\|stripe\|material' docs/design/DESIGN.md` → 0 líneas.
- **U2-01** (cubierto en lo que toca a U2b-1): `grep -n "^## " docs/design/DESIGN.md` → 11 encabezados; los 10 del índice de plan §4.1 en orden (líneas 7, 43, 236, 320, 414, 499, 554, 556, 558, 560) más «Supuestos» (562). Secciones 1–6 con párrafos y tablas; 7–10 vacías por diseño del reparto, las llena U2b-2.
- **U2-10**: §3.1 (Inter 4.1 variable OFL autoalojada, `@font-face` local, `font-display: swap`), §3.2 (escala 12/14/16/20/24/32), §3.3 (`tabular-nums` en todo dinero; negativos con U+2212, color de rol e icono), §5.1–5.2 (Lucide ISC autoalojado, trazo 1.5 px, 20/24 px); licencias en `docs/design/fuentes/inter/LICENSE.txt` y `docs/design/fuentes/lucide/LICENSE`; 0 referencias a terceros.
- **U2-11**: §4.1 espaciado 4/8/12/16/24/32/48; §4.2 rejilla 4/8/12 columnas con los rangos 360/412/768/1024/1440; §4.4 radios 6/10/14; §4.5 dos niveles de sombra en claro y bordes en oscuro, foco visible 2 px con separación 2 px; §4.6 objetivo táctil 44 px (`grep -n "44 px"` → 5 líneas, 0 mínimos táctiles menores; 36 px solo con puntero fino y declarado no táctil); §6.1 150/200/250 ms con curvas estándar y enfática; §6.3 `prefers-reduced-motion` (2 apariciones).
- Apoyo a otros criterios: U2-02 (texto de §2.2 y §2.4 con 280°, 150°/135°, 100°, 310°, gasto 45°), U2-08 (`grep -n -i 'solo por color\|persona y.*tienda'` → 6 líneas; tabla de pares cercanos con refuerzo en §2.6), U2-13 (Notion única referencia, §1.1), U2-14 (español de México, tokens en inglés, sección Supuestos).

## Supuestos
1. La ruta estática en la PWA (`/fonts/`, `/icons/`) es propuesta; la fija WF-3.
2. En la PWA los iconos se usan con el componente de Lucide para React que trae shadcn/ui, misma versión 1.55.0, trazo absoluto 1.5; la dependencia exacta la decide WF-3.
3. Los nombres de token no cromáticos son propuestas de U2b-1 para la sección 10; ningún valor de color cambia.
Ninguno cambia cifras, esquema, endpoints ni tokens de color.

## Bloqueos
Ninguno nuevo. El único abierto del flujo (`docs/estado/bloqueos/S1-parametros-avisos.md`) no toca las secciones 1–6.

## Pendientes
- Secciones 7, 8, 9 y 10 de `docs/design/DESIGN.md` (agente U2b-2), incluidos `C2-tokens.json`, `tokens.css` y el tema de Tailwind con los tokens propuestos en §3–§6.
- Agregar a la sección Supuestos de `docs/design/DESIGN.md` los supuestos de U2b-2.

## Solicitudes fuera de alcance
1. `docs/criterios/U2.md` (U2-11) y la trazabilidad citan `C2.md`, que no existe en `docs/contratos/`: crear `docs/contratos/C2.md` (quien escriba la sección 10) o corregir la cita a `C2-tokens.json`. Motivo: el revisor no podrá comparar «con los tokens de C2 (C2.md)».
2. `docs/design/colorimetria.md` §3 dice «tipo Notion/Linear»: quitar «Linear». Motivo: decisiones §G y U2-13, Notion es la única referencia visual (DESIGN.md no lo repite).
3. `docs/design/U1-arquitectura.md` §6.2 (línea 1135) dice «ancho fijo (U2)» para la barra lateral a 768: citar `layout.sidebar.overlay-width` = 280 px de DESIGN.md §4.3. Motivo: trazabilidad.
4. WF-3: copiar `docs/design/fuentes/inter/web/*.woff2` y `docs/design/fuentes/lucide/sprite.svg` a los estáticos del frontend y fijar la versión 1.55.0 de Lucide. Motivo: autoalojado (protocolo §8).

## Paquetes instalados
Ninguno. Descargas (no paquetes del sistema): `Inter-4.1.zip` (32 MB) y `lucide-static-1.55.0.tgz` (6.6 MB), borradas tras extraer; se conservan 3.1 MB en `docs/design/fuentes/`.

## Propuestas para Axel
Ningún valor K31 cae en este alcance. Decisiones de diseño para ratificar en la puerta 2:
1. Negativos `−$1,258.68` con signo menos tipográfico U+2212 antes de `$` (DESIGN.md §3.3).
2. Negativo sin rol propio (restante rebasado) en advertencia con icono y la palabra «rebasado» (§3.3).
3. Trato de tú en la interfaz (§1.3).
4. Ingreso en 135° (vigente en los tokens desde U2a; decisiones §I dice 150° «o U2 propone 135°»): confirmar explícitamente.
5. Lucide 1.55.0 se publicó el 10-oct-2026, un día antes de la descarga; son SVG estáticos sin código, integridad verificada. Alternativa si se prefiere una versión con más antigüedad: fijar una de hace ≥ 2 semanas al armar el frontend.

## Consumo aproximado
35 acciones de herramienta · ~165 KB leídos · ~68 KB escritos (DESIGN.md 48.7 KB, script 8.2 KB, FUENTES.md 2.6 KB, reportes ~8 KB). Sin cierre por tope.
