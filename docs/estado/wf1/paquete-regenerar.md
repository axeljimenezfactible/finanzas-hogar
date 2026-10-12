# Reporte · paquete-regenerar (WF-1, Sonnet) · 11-oct-2026

## Hecho
- Se ejecutó `docs/design/tools/paquete_entrada.sh` (sin `--forzar`). Solo se reconstruyó `U4-ajustes.pdf` (fuentes `U4-ajustes-a.md` y `-b.md`, modificadas 19:30, después del paquete previo de 18:57); pasó de 30 a 31 páginas. DESIGN (18:27), U3 (18:37) y C1-resumen eran anteriores a la generación previa: sus sellos coinciden y no se reconstruyeron.
- `MANIFIESTO.md` regenerado: cambian solo la fila y la línea de sha256 de `U4-ajustes.pdf`.

## Archivos (tamaño en bytes, sha256 corto)
- `C1-resumen.pdf` 280159 74de369d
- `C2-tokens.json` 100073 570ff553
- `DESIGN.docx` 82222 a48a7693
- `DESIGN.pdf` 2859638 0ca9177b
- `LEEME.md` 5970 179f875f
- `MANIFIESTO.md` 6956 ba064f6e
- `U3-dashboard.pdf` 1186388 ae9752c6
- `U4-ajustes.pdf` 2026629 79596557
- `paleta.png` 865657 8b6bb79f
- `prompt.md` 11958 fd60b9c9
- `reporte-p01.png` 193383 2d8b5215
- `reporte-p02.png` 126390 f23e15e7
- `reporte-p03.png` 217623 a16f3234
- `reporte-p04.png` 224544 493a14a3
- `reporte-p05.png` 213182 639cbe5b
- `reporte-p06.png` 174292 60f15487
- `reporte-p07.png` 247150 c553dcc0
- `reporte-p08.png` 189586 2743264a
- `reporte-p09.png` 194116 2bd286c7
- `reporte-p10.png` 229803 60ae8cb3
- `reporte-p11.png` 310313 80126d7b
- `tokens.css` 26121 2b9429bc

## Verificación
- Sellos de `.cache` recalculados por fuera con la regla del script (sha256 del script + contenido de fuentes): DESIGN, DESIGN-docx, U3-dashboard, U4-ajustes y C1-resumen coinciden (5/5).
- `grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md | sha256sum -c`: 0 fallos (el script lo reporta también: «0 diferencias»). 22 archivos en la carpeta.
- Archivos .go, .py, .ts, .js, .sh en `docs/design/entrada/`: 0.
- `LEEME.md` solo cita nombres (U4-ajustes.pdf, reporte-p01..p11.png, 18 adjuntos = 7 + 11); no cita tamaños ni páginas, y todos los nombres existen. Sin cambios necesarios.
- Comparación antes/después por sha256: únicos cambios U4-ajustes.pdf y MANIFIESTO.md.

## Supuestos
- Los PDF no son reproducibles byte a byte (fecha de Chromium); el caché evita regenerar si las fuentes no cambian.
- Camino de PDF sigue siendo html+chromium (lualatex falla por métrica de fuente lmroman10, hallazgo previo).

## Bloqueos
- Ninguno.

## Pendientes
- Si U3, DESIGN o U4 reciben más cambios, volver a correr el script. U4-ajustes.md unificado aún no existe.

## Solicitudes fuera de alcance
- Ninguna.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna.

## Consumo aproximado
- Acciones 4, KB leídos ~35, KB escritos ~3 (más artefactos regenerados por el script).
