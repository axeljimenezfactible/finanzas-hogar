# correccion-paquete-r2 · Corrección del paquete de entrada para Claude Design (ronda 2)

Fecha: 11-oct-2026. Fuente: `docs/estado/wf1/REVISION-paquete-r2.md` (solo H-1, el único medio; los bajos B-1 a B-7 quedan fuera de esta corrección).

## Hecho
- H-1 (media, P3-03 / P3-06 / verificación d): paso 7 de `docs/design/entrada/LEEME.md` reescrito. Ahora dice: exportar en PDF (un archivo con todas las pantallas) y, si se puede, HTML independiente; si además deja bajar imágenes, también sirve; con solo PDF basta, no se toman capturas y no hay que detenerse; URL, Canva y PPTX no se usan sin avisar; si no hay ni PDF, ni HTML, ni imágenes, detenerse y avisar. Se cuenta «84 artboards» y cada página de un PDF vale un artboard.
- Ejemplos de nombres (sección «Cómo nombrar lo que descargues») sin extensión fija, con la regla de cuándo llevan `.png` y que un PDF conserva el nombre de la herramienta.
- «Si algo falla»: aclarado que no ofrecer imágenes no cuenta como «cambio de opciones» (evita detenerse sin necesidad).
- `docs/design/salida/LEEME.md`: «84 imágenes» pasa a «84 artboards (imágenes o páginas de un PDF)»; «Un PDF de N páginas cuenta N artboards. Con PDF basta: no se toman capturas…».
- `MANIFIESTO.md` regenerado con `paquete_entrada.sh` (sin `--forzar`): solo cambió el hash/tamaño de LEEME.md (7539 a 8212 bytes).

## Archivos
- Modificados: `/home/deepsec/finanzas-hogar/docs/design/entrada/LEEME.md`, `/home/deepsec/finanzas-hogar/docs/design/salida/LEEME.md`, `/home/deepsec/finanzas-hogar/docs/design/entrada/MANIFIESTO.md` (regenerado).
- Sin tocar: `prompt.md` (12284 bytes, su «Qué entregar» ya es coherente con el nuevo paso 7), PDF, DOCX, PNG, JSON, CSS, `paquete_entrada.sh`.
- Respaldo previo en el scratchpad de la sesión (fuera del repositorio).

## Verificación
- `docs/design/tools/paquete_entrada.sh` (sin `--forzar`): 5/5 pasos; 4 PDF y DOCX «sin cambios en las fuentes»; 11 PNG sin cambios; `/api/v1` = 0 en los 4 PDF; «sha256sum -c: 0 diferencias»; «Listo: 22 archivos».
- `diff -rq` contra copia previa: solo difiere MANIFIESTO.md. Los otros 20 hashes del MANIFIESTO son idénticos a los anteriores.
- `grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md | sha256sum -c`: 0 líneas fuera de OK. `ls -1A entrada | wc -l` = 22.
- Credenciales/contraseñas en LEEME: 0. `prompt.md`: sin cambios (12284 bytes, tope 12288).
- Criterios: P3-02 (22 elementos, huellas en 0), P3-03 (pasos numerados, es-MX, ahora sin ambigüedad de formato de exporte), P3-05 (cuándo parar: detenerse solo si no hay ningún formato utilizable), P3-06 (el conteo de artboards admite PDF e imágenes), P3-11 (carpetas `iter-1/`, `iter-2/`, `iter-3/` y `final/` intactas con `.gitkeep`).
- No hay pruebas Go, oráculo ni `compara.py` aplicables: el artefacto es documental y el cambio no toca cifras.

## Supuestos
- Ninguno que cambie cifras, esquema, endpoints o tokens. Interpretación: el criterio de parada del paso 7 («ni PDF ni HTML ni imágenes») es mi redacción para cerrar la ambigüedad, en línea con la corrección sugerida por el revisor.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Bajos B-1 a B-7 de la revisión (fuera del encargo): prompt.md a 4 bytes del tope, nombres de operación en U4 (`logout`, `login`), «ningún archivo contiene código» frente a `tokens.css`, aviso de que `siguiente-prompt.md` ya existe, orden de la regla general de «Si algo falla», «la lista de nombres está también en salida/LEEME.md» y trazabilidad del reporte original.
- El LEEME pasó de 7.4 a 8.0 KB; no tiene tope.

## Solicitudes fuera de alcance
- Ninguna.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna.

## Consumo aproximado
- Acciones de herramienta: 8. Leído: ~45 KB. Escrito: ~12 KB (LEEME x2, MANIFIESTO regenerado, este reporte).
