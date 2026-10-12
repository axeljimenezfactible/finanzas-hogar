#!/usr/bin/env bash
# Regenera el brief de referencia de la Pista 3 (K45, sustituye a K35; criterios P3-02). No se sube a ningún sitio.
# Uso: docs/design/tools/paquete_entrada.sh [--forzar]
#   Sin opciones solo reconstruye lo que cambió en las fuentes: así MANIFIESTO.md
#   solo cambia si cambian las fuentes (los PDF de Chromium llevan fecha de creación
#   y no son idénticos byte a byte entre corridas; el caché de sellos lo evita).
#   --forzar reconstruye todo (los PDF cambian de hash).
# Requiere: pandoc, chromium (o lualatex), pdftoppm, pdfinfo, pdftotext, perl, sha256sum.
# Los archivos escritos a mano (LEEME.md, prompt.md) no se tocan; solo se hashean.
# La fuente de C1-resumen.pdf vive en docs/design/tools/paquete/C1-resumen.md
# (en la carpeta de entrada solo van los 22 archivos de P3-02).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$ROOT"
DESIGN="docs/design"
OUT="$DESIGN/entrada"
FUENTES_PAQ="$DESIGN/tools/paquete"
CACHE="$FUENTES_PAQ/.cache"
REPORTE_PDF="${REPORTE_PDF:-/home/deepsec/reporte-compra-latex/build/reporte.pdf}"
TMP="$ROOT/datos/tmp/paquete"
FORZAR=0
[[ "${1:-}" == "--forzar" ]] && FORZAR=1

# Fecha fija para los formatos que la respetan (DOCX, lualatex): salidas reproducibles.
export SOURCE_DATE_EPOCH=1791676800 FORCE_SOURCE_DATE=1 TZ=America/Mexico_City LC_ALL=C.UTF-8

mkdir -p "$OUT" "$CACHE" "$TMP"
limpiar() { if [[ "$TMP" == "$ROOT"/datos/tmp/paquete && -d "$TMP" ]]; then rm -rf -- "$TMP"; fi; }
trap limpiar EXIT
SCRIPT_HASH="$(sha256sum "${BASH_SOURCE[0]}" | cut -c1-64)"

log() { printf '%s\n' "$*" >&2; }
die() { log "ERROR: $*"; exit 1; }
for c in pandoc pdftoppm pdfinfo sha256sum; do command -v "$c" >/dev/null || die "falta $c"; done

# --- caché de sellos -------------------------------------------------------
# necesita_construir <nombre> <archivo_salida> <fuentes...>
necesita_construir() {
  local nombre="$1" salida="$2"; shift 2
  local sello; sello="$( { printf '%s\n' "$SCRIPT_HASH"; cat -- "$@"; } | sha256sum | cut -c1-64)"
  NUEVO_SELLO="$sello"
  [[ $FORZAR -eq 1 ]] && return 0
  [[ -s "$salida" && -f "$CACHE/$nombre.sello" && "$(cat "$CACHE/$nombre.sello")" == "$sello" ]] && return 1
  return 0
}
guardar_sello() { printf '%s' "$NUEVO_SELLO" > "$CACHE/$1.sello"; }

# --- Markdown -> PDF -------------------------------------------------------
escribir_css() {
cat > "$TMP/estilo.css" <<CSS
@font-face { font-family: "Inter"; src: url("file://$ROOT/$DESIGN/fuentes/inter/ttf/InterVariable.ttf"); font-weight: 100 900; }
@page { size: A4; margin: 14mm 12mm 16mm 12mm; @bottom-center { content: counter(page) " / " counter(pages); font: 8pt "Inter", sans-serif; color: #66676E; } }
html { font-size: 9.5pt !important; }
body { font-family: "Inter", "Noto Sans", sans-serif !important; line-height: 1.4 !important; color: #222325 !important; background: #FFFFFF !important; max-width: none !important; margin: 0 !important; padding: 0 !important; }
h1 { font-size: 18pt; margin: 0 0 6pt; } h2 { font-size: 13.5pt; margin: 14pt 0 5pt; border-bottom: 1px solid #D1D2D8; padding-bottom: 2pt; break-after: avoid; }
h3 { font-size: 11pt; margin: 10pt 0 4pt; break-after: avoid; } h4 { font-size: 10pt; break-after: avoid; }
table { border-collapse: collapse; width: 100%; margin: 6pt 0; font-size: 8pt; table-layout: auto; display: table !important; }
th, td { border: 1px solid #D1D2D8; padding: 2.5pt 4pt; vertical-align: top; text-align: left; overflow-wrap: break-word; }
th { background: #F2F2F4; } thead { display: table-header-group; } tr { break-inside: avoid; }
code { font-family: "Noto Sans Mono", monospace; font-size: .86em; overflow-wrap: break-word; background: #F2F2F4; padding: 0 1.5pt; }
pre { white-space: pre-wrap; font-size: 7.6pt; background: #F2F2F4; padding: 5pt; overflow-wrap: break-word; } pre code { background: none; }
blockquote { margin: 6pt 0; padding-left: 8pt; border-left: 3px solid #C9CFFF; color: #4D4E53; }
a { color: #5D609F; } img { max-width: 100%; }
nav#TOC { font-size: 9pt; break-after: page; } nav#TOC ul { list-style: none; padding-left: 10pt; }
CSS
}

# md_a_pdf <nombre> <titulo> <salida.pdf> <fuente.md...>
md_a_pdf() {
  local nombre="$1" titulo="$2" salida="$3"; shift 3
  local fuente="$TMP/$nombre.md"
  : > "$fuente"
  local f; for f in "$@"; do cat -- "$f" >> "$fuente"; printf '\n\n' >> "$fuente"; done
  # Camino 1: lualatex (el pedido por el plan).
  if timeout 300 pandoc "$fuente" -f markdown -o "$TMP/$nombre.lualatex.pdf" --pdf-engine=lualatex \
       --metadata pagetitle="$titulo" -V geometry:margin=2cm -V mainfont="Noto Sans" -V monofont="Noto Sans Mono" -V fontsize=10pt \
       >"$TMP/$nombre.lualatex.log" 2>&1 && [[ -s "$TMP/$nombre.lualatex.pdf" ]]; then
    cp "$TMP/$nombre.lualatex.pdf" "$salida"; echo "lualatex" > "$CACHE/$nombre.camino"
    log "  $nombre: camino lualatex"; return 0
  fi
  # Camino 2: HTML autónomo + Chromium sin interfaz.
  command -v chromium >/dev/null || die "lualatex falló y no hay chromium"
  escribir_css
  pandoc "$fuente" -f markdown --columns=100000 --standalone --embed-resources --css "$TMP/estilo.css" \
    --metadata pagetitle="$titulo" --toc --toc-depth=2 -o "$TMP/$nombre.html" 2>"$TMP/$nombre.pandoc.log"
  chromium --headless=new --no-sandbox --disable-gpu --allow-file-access-from-files --no-pdf-header-footer \
    --user-data-dir="$TMP/chromium-$nombre" --virtual-time-budget=20000 \
    --print-to-pdf="$TMP/$nombre.pdf" "file://$TMP/$nombre.html" >"$TMP/$nombre.chromium.log" 2>&1 || true
  [[ -s "$TMP/$nombre.pdf" ]] || die "no se generó $salida (ver $TMP/$nombre.chromium.log)"
  cp "$TMP/$nombre.pdf" "$salida"
  echo "html+chromium (lualatex falló: $(grep -m1 -E '^!' "$TMP/$nombre.lualatex.log" 2>/dev/null | tr -d '\n' || true))" > "$CACHE/$nombre.camino"
  log "  $nombre: camino html+chromium"
}

# --- Filtro de rutas y nombres de operación (decisión 42, M-6 de la ronda 1) ---
# U3 y U4 se escribieron para quien programa la API: traen rutas "/api/v1/..." y nombres
# de operación (getX, listX, updateX...). El paquete para Claude Design no lleva ni lo uno
# ni lo otro (igual que C1-resumen): se filtra la copia que va a pandoc; las fuentes no cambian.
# filtrar_api <salida.md> <fuente.md...>
filtrar_api() {
  local salida="$1"; shift
  command -v perl >/dev/null || die "falta perl (filtro de rutas de la API)"
  : > "$salida"
  local f; for f in "$@"; do cat -- "$f" >> "$salida"; printf '\n\n' >> "$salida"; done
  perl -0pi -e '
    s{\((?:GET|POST|PUT|PATCH|DELETE) `/api/v1/[^`]*`, (campo `[^`]*`)\)}{($1)}g;   # ruta + campo: queda el campo
    s{\s*\((?:GET|POST|PUT|PATCH|DELETE) `/api/v1/[^`]*`\)}{}g;                    # (GET `/api/v1/...`)
    s{`operationId`}{el identificador de la operación}g;
    s{`(?:get|list)[A-Z][A-Za-z0-9]*`}{consulta de C1}g;                             # nombres de operación de lectura
    s{`(?:update|create|delete|set|change|confirm|dismiss|duplicate|generate|import|preview|query|reload|reorder|restore|revoke|send|suggest|post|put|patch|run|export|test|reset|cancel|apply|archive|upload|download|move|merge|retry|dry)[A-Z][A-Za-z0-9]*`}{operación de C1}g;
    s{(consulta de C1|operación de C1)( / \1)+}{$1}g;                               # "X / X" -> "X"
    s{\b(de|con|por|mediante) (consulta|operación) de C1}{$1 la $2 de C1}g;           # concordancia en prosa
  ' "$salida"
}
# sin_api <pdf...>: falla si el texto de algún PDF aún trae rutas "/api/v1"
sin_api() {
  local pdf n
  for pdf in "$@"; do
    n="$(pdftotext "$pdf" - 2>/dev/null | grep -c '/api/v1' || true)"
    [[ "$n" == "0" ]] || die "$(basename "$pdf") trae $n líneas con /api/v1"
  done
}
# sin_operaciones <pdf...>: falla si el texto trae nombres de operación (getX, listX, updateX...)
# Se aplica a U3, U4 y C1-resumen; DESIGN.pdf nombra una operación solo en el texto PROVISIONAL del bloqueo S1.
sin_operaciones() {
  local pdf m
  for pdf in "$@"; do
    m="$(pdftotext "$pdf" - 2>/dev/null | grep -cE '\b(get|list|update|create|delete)[A-Z][A-Za-z]+' || true)"
    [[ "$m" == "0" ]] || die "$(basename "$pdf") trae $m líneas con nombres de operación"
  done
}

construir_pdf() { # nombre titulo salida fuentes...
  local nombre="$1" titulo="$2" salida="$3"; shift 3
  if necesita_construir "$nombre" "$salida" "$@"; then
    md_a_pdf "$nombre" "$titulo" "$salida" "$@"; guardar_sello "$nombre"
  else log "  $nombre: sin cambios en las fuentes"; fi
}

log "1/5 PDF y DOCX"
construir_pdf DESIGN "DESIGN · Finanzas Hogar" "$OUT/DESIGN.pdf" "$DESIGN/DESIGN.md"
filtrar_api "$TMP/U3-filtrado.md" "$DESIGN/U3-dashboard.md"
construir_pdf U3-dashboard "U3 · Dashboard · Finanzas Hogar" "$OUT/U3-dashboard.pdf" "$TMP/U3-filtrado.md"
if [[ -f "$DESIGN/U4-ajustes.md" ]]; then U4_FUENTES=("$DESIGN/U4-ajustes.md")
else U4_FUENTES=("$DESIGN/U4-ajustes-a.md" "$DESIGN/U4-ajustes-b.md"); fi
filtrar_api "$TMP/U4-filtrado.md" "${U4_FUENTES[@]}"
construir_pdf U4-ajustes "U4 · Ajustes · Finanzas Hogar" "$OUT/U4-ajustes.pdf" "$TMP/U4-filtrado.md"
construir_pdf C1-resumen "C1 · Resumen del contrato de datos · Finanzas Hogar" "$OUT/C1-resumen.pdf" "$FUENTES_PAQ/C1-resumen.md"

if necesita_construir DESIGN-docx "$OUT/DESIGN.docx" "$DESIGN/DESIGN.md"; then
  pandoc "$DESIGN/DESIGN.md" -f markdown --columns=100000 --toc --toc-depth=2 --metadata title="DESIGN · Finanzas Hogar" -o "$OUT/DESIGN.docx"
  guardar_sello DESIGN-docx; log "  DESIGN.docx generado"
else log "  DESIGN.docx: sin cambios en las fuentes"; fi

log "2/5 Copias tal cual"
cp -p docs/contratos/C2-tokens.json "$OUT/C2-tokens.json"
cp -p docs/contratos/tokens.css "$OUT/tokens.css"
cp -p "$DESIGN/paleta.png" "$OUT/paleta.png"

log "3/5 Páginas del reporte previo"
[[ -s "$REPORTE_PDF" ]] || die "no existe $REPORTE_PDF"
PAGS="$(pdfinfo "$REPORTE_PDF" | awk '/^Pages:/{print $2}')"
[[ "$PAGS" == "11" ]] || die "el reporte tiene $PAGS páginas y se esperaban 11"
if necesita_construir reporte-png "$OUT/reporte-p01.png" "$REPORTE_PDF"; then
  pdftoppm -r 110 -png "$REPORTE_PDF" "$TMP/rp"
  for i in 01 02 03 04 05 06 07 08 09 10 11; do mv "$TMP/rp-$i.png" "$OUT/reporte-p$i.png"; done
  guardar_sello reporte-png; log "  11 PNG generados"
else log "  reporte-p01..p11: sin cambios en las fuentes"; fi

log "4/5 Revisión de contenido (sin código)"
if find "$OUT" \( -name '*.go' -o -name '*.ts' -o -name '*.tsx' -o -name '*.py' -o -name '*.sql' \) | grep -q .; then die "hay código en $OUT"; fi
for f in LEEME.md prompt.md; do [[ -s "$OUT/$f" ]] || die "falta $OUT/$f (se escribe a mano)"; done
command -v pdftotext >/dev/null || die "falta pdftotext (comprobación de rutas de la API)"
sin_api "$OUT/U3-dashboard.pdf" "$OUT/U4-ajustes.pdf" "$OUT/C1-resumen.pdf" "$OUT/DESIGN.pdf"
sin_operaciones "$OUT/U3-dashboard.pdf" "$OUT/U4-ajustes.pdf" "$OUT/C1-resumen.pdf"
log "  pdftotext | grep -c '/api/v1' = 0 en los 4 PDF; sin nombres de operación en U3, U4 y C1-resumen"

log "5/5 MANIFIESTO.md"
declare -A SIRVE=(
  [LEEME.md]="Qué pasa en la Pista 3 por Artifacts (K45): quién la ejecuta, dónde quedan las salidas y cómo aprueba Axel; ya no se sube nada."
  [prompt.md]="Brief de los agentes de la Pista 3 (K45): pantallas de U1 a maquetar, cifras sintéticas, tokens de color en texto, reglas de dinero y accesibilidad."
  [DESIGN.pdf]="Sistema de diseño completo (principios, color, tipografía, componentes, accesibilidad)."
  [DESIGN.docx]="El mismo sistema de diseño en Word, por si el PDF no se puede leer."
  [U3-dashboard.pdf]="Qué lleva el Inicio: Resumen (6 tarjetas), Detalle, Tendencias y Calendario."
  [U4-ajustes.pdf]="Las 14 secciones de Ajustes con sus campos y estados."
  [C1-resumen.pdf]="Resumen sin código del contrato de datos: módulos, pantallas, entidades, métricas y formatos."
  [C2-tokens.json]="Todos los valores de diseño (colores, tipografía, espaciado, radios) en un solo archivo."
  [tokens.css]="Los mismos valores como variables de hojas de estilo, para consulta."
  [paleta.png]="Imagen de la paleta de color en claro y oscuro."
)
ORDEN=(LEEME.md prompt.md DESIGN.pdf DESIGN.docx U3-dashboard.pdf U4-ajustes.pdf C1-resumen.pdf C2-tokens.json tokens.css paleta.png)
for i in 01 02 03 04 05 06 07 08 09 10 11; do
  SIRVE[reporte-p$i.png]="Página $((10#$i)) de 11 del reporte previo de compras: referencia de forma y orden; sus cifras son reales y no se copian a las maquetas (K45.1)."
  ORDEN+=("reporte-p$i.png")
done
M="$OUT/MANIFIESTO.md"
{
  echo "# MANIFIESTO · Brief de referencia de la Pista 3 (K45)"
  echo
  echo "Lo genera \`docs/design/tools/paquete_entrada.sh\`; no se edita a mano. Cambia solo si cambian los archivos de origen."
  echo
  echo "- Carpeta del brief de referencia (no se sube; K45): \`/home/deepsec/finanzas-hogar/docs/design/entrada/\`"
  echo "- Copia en disco de cada vuelta (la hace el orquestador): \`/home/deepsec/finanzas-hogar/docs/design/salida/iter-N/\`, con N = 1, 2 o 3"
  echo "- Exportación final del lienzo: \`/home/deepsec/finanzas-hogar/docs/design/salida/final/\`"
  echo "- Ningún archivo del paquete contiene código ni credenciales."
  echo
  echo "## Archivos (${#ORDEN[@]}; este MANIFIESTO no se lista a sí mismo)"
  echo
  echo "| # | Archivo | Tamaño | sha256 | Para qué sirve |"
  echo "|---|---|---|---|---|"
  n=0
  for f in "${ORDEN[@]}"; do
    [[ -s "$OUT/$f" ]] || die "falta $OUT/$f"
    n=$((n+1)); b="$(stat -c %s "$OUT/$f")"; h="$(sha256sum "$OUT/$f" | cut -c1-64)"
    printf '| %s | `%s` | %s bytes (%s KB) | `%s` | %s |\n' "$n" "$f" "$b" "$(awk -v b="$b" 'BEGIN{printf "%.1f", b/1024}')" "$h" "${SIRVE[$f]}"
  done
  echo
  echo "## Cómo verificarlo"
  echo
  echo "Desde la carpeta de entrada, este comando no debe mostrar ninguna diferencia (todas las líneas terminan en OK):"
  echo
  echo '```'
  echo "cd /home/deepsec/finanzas-hogar/docs/design/entrada && grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md | sha256sum -c"
  echo '```'
  echo
  echo "## Líneas de sha256sum"
  echo
  echo '```'
  for f in "${ORDEN[@]}"; do printf '%s  %s\n' "$(sha256sum "$OUT/$f" | cut -c1-64)" "$f"; done
  echo '```'
} > "$M"

# Cuenta final: 22 elementos (21 + MANIFIESTO).
TOTAL="$(find "$OUT" -maxdepth 1 -type f | wc -l)"
[[ "$TOTAL" == "22" ]] || die "se esperaban 22 archivos en $OUT y hay $TOTAL"
( cd "$OUT" && grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md | sha256sum -c --quiet ) && log "sha256sum -c: 0 diferencias"
for p in DESIGN U3-dashboard U4-ajustes C1-resumen; do log "  $p.pdf: $(pdfinfo "$OUT/$p.pdf" | awk '/^Pages:/{print $2}') páginas; camino: $(cat "$CACHE/$p.camino" 2>/dev/null || echo sin registro)"; done
log "Listo: $TOTAL archivos, $(du -sk "$OUT" | cut -f1) KB en $OUT"
