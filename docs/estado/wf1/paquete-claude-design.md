# Reporte · paquete-claude-design (WF-1, Sonnet) · 11-oct-2026

## Hecho
- `docs/design/tools/paquete_entrada.sh`: regenera todo el paquete desde la raíz del repositorio (PDF y DOCX de DESIGN, PDF de U3, U4 y C1-resumen; copia tal cual de `C2-tokens.json`, `tokens.css` y `paleta.png`; `reporte-p01..p11.png` con `pdftoppm -r 110`; `MANIFIESTO.md`). Falla si el reporte no tiene 11 páginas, si la carpeta no queda con 22 archivos o si hay código en ella. Usa un caché de sellos (`docs/design/tools/paquete/.cache/`) para que `MANIFIESTO.md` solo cambie si cambian las fuentes; `--forzar` reconstruye todo.
- Camino de los PDF: **HTML autónomo con pandoc y `chromium --headless=new --print-to-pdf`** (Inter embebida desde `docs/design/fuentes/inter/`, solo hex neutros y de marca de C2). El script intenta primero `pandoc --pdf-engine=lualatex`; en este equipo falla siempre con `! Font \TU/lmr/m/n/10=[lmroman10-regular]:+tlig; at 10pt not loadable: metric data not found or bad` (aun con `mainfont` Noto Sans; el archivo `lmroman10-regular.otf` existe en texmf). Se registró como hallazgo del entorno; el camino usado por PDF queda en `docs/design/tools/paquete/.cache/<nombre>.camino`.
- U4: `U4-ajustes.md` aún no existe, así que `U4-ajustes.pdf` sale de la concatenación de `U4-ajustes-a.md` y `U4-ajustes-b.md`. El integrador debe volver a correr el script cuando las una (el script ya prefiere `U4-ajustes.md` si existe).
- `docs/design/entrada/prompt.md` (11 958 bytes): 20 pantallas (4 de Inicio, 14 de Ajustes, barra lateral expandida/colapsada, menú «/» con búsqueda) × 4 combinaciones = 80 artboards; 90 hex únicos, todos generados desde `C2-colores.json` por un script temporal (marca y neutros 50–950, 8 roles con insignia, acción crítica, 4 personas, 4 tiendas, 8 series, blanco).
- `docs/design/tools/paquete/C1-resumen.md` (fuente) y `docs/design/entrada/C1-resumen.pdf` (4 páginas): módulos y pantallas, 14 secciones de Ajustes, 31 entidades por tema, tarjetas y métricas por nombre, formatos; sin rutas ni JSON de la API.
- `docs/design/entrada/LEEME.md` para Axel (tabla con orden de subida, 8 pasos, nombres de archivos, vueltas con «iter-N lista», sección «Si algo falla», qué no se hace, comprobación con `sha256sum -c`).
- `docs/design/salida/LEEME.md` y carpetas `iter-1`, `iter-2`, `iter-3` y `final` con `.gitkeep`.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/design/tools/paquete_entrada.sh`
- `/home/deepsec/finanzas-hogar/docs/design/tools/paquete/C1-resumen.md` y `.cache/` (sellos y camino de cada PDF)
- `/home/deepsec/finanzas-hogar/docs/design/entrada/` (22 archivos): `LEEME.md`, `prompt.md`, `DESIGN.pdf`, `DESIGN.docx`, `U3-dashboard.pdf`, `U4-ajustes.pdf`, `C1-resumen.pdf`, `C2-tokens.json`, `tokens.css`, `paleta.png`, `reporte-p01.png` a `reporte-p11.png`, `MANIFIESTO.md`
- `/home/deepsec/finanzas-hogar/docs/design/salida/LEEME.md`, `salida/iter-1|iter-2|iter-3|final/.gitkeep`
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/paquete-claude-design.md` (este reporte)

## Verificación
- `ls docs/design/entrada | wc -l` = 22; archivos de 0 bytes = 0 (P3-02).
- Dos corridas seguidas sin cambios: `diff` de `MANIFIESTO.md` = 0 líneas. Cambiando una fuente (línea añadida a `C1-resumen.md`): solo se reconstruyó `C1-resumen.pdf`, y el `diff` mostró solo su fila y su línea de sha256; fuente restaurada, corrida final: `diff` = 0. `--forzar` regenera los PDF con fecha nueva (hash distinto), por eso el caché.
- `grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md | sha256sum -c` desde `docs/design/entrada`: 21 OK, 0 fallos (P3-01, P3-02). Corrida final del script: «sha256sum -c: 0 diferencias».
- `pdfinfo`: DESIGN.pdf 36 páginas, U3-dashboard.pdf 19, U4-ajustes.pdf 30, C1-resumen.pdf 4 (todas > 0; C1-resumen ≤ 10). `DESIGN.docx`: 16 partes en el zip (válido).
- `pdftotext DESIGN.pdf`: los 10 encabezados de sección encontrados (Principios y tono, Colorimetría, Tipografía, Espaciado…, Iconografía, Movimiento, Componentes, Patrones de visualización, Accesibilidad…, Tokens); «Colorimetr» aparece 2 veces (≥ 1). Revisión visual de las páginas 6 y 20 de DESIGN.pdf y 8 de U4: tablas legibles, sin texto cortado (se corrigió `overflow-wrap:anywhere` a `break-word` tras ver columnas partidas).
- 11 PNG del reporte: 910 × 1287 px cada uno (110 ppp); `paleta.png` 1400 × 3000.
- Hex de `prompt.md` contra `C2-colores.json` (sin distinguir mayúsculas): 90 únicos, 0 ausentes. `grep -c "color" prompt.md` = 40. `prompt.md` pesa 11 958 bytes (≤ 12 KB), sin código.
- `find docs/design/entrada` con `*.go`, `*.ts`, `*.tsx`, `*.py`, `*.sql`: 0.
- Patrones de credencial (P3-01, P3-03, P3-12) sobre `LEEME.md` y `prompt.md`: 0 coincidencias; sobre el texto de los 4 PDF: 0. Las palabras compartir, publicar, invitar y suscripción en `LEEME.md` aparecen solo en la advertencia de lo que no se hace (líneas «Si se agota el límite de tu suscripción» y «No compartas el proyecto…»).
- Tamaño total del paquete: 9 548 KB (9.4 MB) en 22 archivos.
- `C1-resumen.md`: 0 coincidencias de rutas de API, `openapi` o llaves.
- Criterios P3 aplicables al paquete:
  - P3-01 (paquete listo): 22 elementos y `sha256sum -c` en 0 (cumplido); **falta el revisor propio del paquete** (lo asigna el orquestador, ver Pendientes).
  - P3-02: línea de evidencia `P3-02`: `ls -l` = 11 archivos + 11 PNG = 22; 4 PDF con `Pages:` ≥ 1; «Colorimetr» ≥ 1; hex de C2 en `prompt.md`; `MANIFIESTO.md` con 21 sha256; `sha256sum -c` 0 diferencias. Cumple.
  - P3-03: `LEEME.md` con 8 pasos numerados (abrir, subir, pegar, enviar, exportar, descargar, guardar en `docs/design/salida/iter-N/` y `final/`); sin credenciales; ejecución en seco: todas las rutas citadas existen. Cumple.
  - P3-04: no aplica al orquestador; el `LEEME.md` y el `prompt.md` dicen que lo descargado es dato. Este agente hizo 0 acciones de navegador (herramientas usadas: Bash, Read, Edit, Write).
  - P3-05: `LEEME.md` trae «Si algo falla» con 2 o 3 intentos y la frase «el proyecto sigue con DESIGN.md»; `grep -c "DESIGN.md" LEEME.md` = 3. Cumple.
  - P3-06: carpetas `iter-1..3` creadas y explicadas (aviso «iter-N lista», `siguiente-prompt.md`); la evidencia por vuelta es posterior. Pendiente de Axel.
  - P3-11: `grep -c "docs/design/salida/iter-N" LEEME.md` = 1; `ls -d` de `salida/`, `salida/final/`, `entrada/` e `iter-1..3` existe. Cumple.
  - P3-12: 0 patrones de credencial en `prompt.md`; las capturas son páginas del reporte ya aprobado por el plan (no revisadas una por una por este agente).
  - P3-07 a P3-10: no aplican al paquete (se evalúan por vuelta). P3-13: ver Supuestos y Bloqueos.

## Supuestos
1. `C1-resumen.md` vive en `docs/design/tools/paquete/` y no en `docs/design/entrada/`, porque P3-02 exige «exactamente» 22 elementos y Markdown no es entrada de Claude Design. Si el integrador prefiere el `.md` en la carpeta de entrada, hay que cambiar P3-02 a 23 elementos.
2. «≤ 12 KB» se interpretó como menos de 12 000 bytes (`prompt.md` pesa 11 958).
3. 20 pantallas × 4 combinaciones (80 artboards) es el conteo de artboards de `prompt.md` para P3-06; el prompt pide priorizar las pantallas 1 a 8 si el límite del servicio no alcanza.
4. Claude Design acepta PDF y PNG (plan §1.5); `C2-tokens.json` y `tokens.css` pueden ser rechazados, y `LEEME.md` dice que se siga sin ellos porque los colores y medidas van en el mensaje. Las medidas de tipografía, espaciado y radios del `prompt.md` salen de `tokens.css` (no de `C2-colores.json`), solo los hex salen de `C2-colores.json`.
5. `U4-ajustes.pdf` nombra operaciones de la API en prosa (p. ej. `listCategories`); se consideró contenido de documento, no código (decisión 42).
6. Los PDF de Chromium no son idénticos byte a byte entre corridas; la estabilidad del manifiesto depende del caché de sellos (si se borra el caché, el manifiesto cambia aunque las fuentes no cambien).
7. El mensaje al usuario del `LEEME.md` pide subir 18 archivos (7 + 11 imágenes) y pegar `prompt.md`; `DESIGN.docx` es respaldo y no se sube salvo que el PDF falle.

## Bloqueos
Ninguno nuevo. Hallazgo de entorno (no bloquea): `lualatex` no carga Latin Modern en este equipo (ver «Hecho»); posiblemente ligado al `-Syu` pendiente de Arch.

## Pendientes
- Revisor propio del paquete (P3-01): `docs/estado/wf1/` aún sin su reporte.
- Volver a correr `docs/design/tools/paquete_entrada.sh` cuando el integrador una `U4-ajustes-a.md` y `U4-ajustes-b.md` en `U4-ajustes.md`, y cuando cambie `DESIGN.md`, `U3-dashboard.md`, `C2-tokens.json` o `tokens.css` (los archivos de origen cambiaron durante WF-1: `DESIGN.md` y `U3-dashboard.md` modificados a las 18:27 y 18:37).
- Cuando exista `docs/design/U4-ajustes.md`, retirar las partes `a` y `b` es decisión del integrador; el script no las borra.

## Solicitudes fuera de alcance
- `docs/design/U4-ajustes.md`: unir `U4-ajustes-a.md` y `U4-ajustes-b.md` (integrador de WF-1) y volver a correr el script.
- `docs/criterios/P3-claude-design.md` P3-02: si se quiere conservar `C1-resumen.md` junto al PDF, ajustar «exactamente 22 elementos» (ver Supuesto 1).
- `.gitignore` o commit del integrador: decidir si `docs/design/tools/paquete/.cache/` se versiona (son 10 archivos de texto pequeños).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna (no se usó ningún valor operativo de K31).

## Consumo aproximado
- Acciones de herramienta: 38. Lectura: unos 150 KB (protocolo, §K de decisiones, criterios P3, extractos de plan, DESIGN §1–§2 y encabezados, U1 §1 y §2.1, U3 por extractos, C1-datos por encabezados y tablas de entidades y métricas, tokens). Escritura: unos 80 KB (script 11 KB, `C1-resumen.md` 10 KB, plantilla y generador del prompt 19 KB, `prompt.md` 12 KB, `LEEME.md` 6 KB, `salida/LEEME.md` 3 KB, este reporte 8 KB); los PDF, el DOCX y los PNG los genera el script (9.5 MB, no cuentan como lectura ni escritura de contexto).
