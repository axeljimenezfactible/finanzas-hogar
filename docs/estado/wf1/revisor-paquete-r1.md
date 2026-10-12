# Reporte · revisor-paquete-r1 (WF-1, revisor independiente, ronda 1 de 3) · 11-oct-2026

## Hecho
- Revisión independiente del paquete de entrada para Claude Design contra P3 v1.2, plan §3.6 y §K (K35). Informe: `docs/estado/wf1/REVISION-paquete-r1.md`.
- Veredicto: **NO APROBADO**: 0 altos, 6 medios (M-1 a M-6), 7 bajos (B-1 a B-7).
- Verificaciones (a)–(g) reejecutadas con comandos y scripts propios (`check_prompt.py`, `check_ui.py`, `paths.py`, `greps.sh` en `datos/tmp/rev-paquete-r1/scripts/`, ya borrados). El script del paquete se corrió sobre una copia del árbol (sin `--forzar` y con `--forzar`) para no tocar la entrega.
- No se corrigió nada del paquete.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-paquete-r1.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-paquete-r1.md` (este reporte)

## Verificación
- (a) `cd docs/design/entrada && grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md | sha256sum -c` → 21 OK, rc=0. Script sobre copia sin `--forzar`: todo «sin cambios», `diff -r` = 0 contra la entrada real. Con `--forzar`: cambian solo los 4 PDF y sus filas del MANIFIESTO; texto (`pdftotext -layout`, md5) y páginas idénticos; DOCX y 11 PNG idénticos byte a byte.
- (b) `pdfinfo`: DESIGN 36, U3 19, U4 31, C1-resumen 4 páginas; los 10 encabezados de DESIGN.md en DESIGN.pdf (2 apariciones cada uno); «Colorimetr» = 2; 11 PNG 910 × 1287.
- (c) `prompt.md` 11 958 bytes; 90 hex únicos, 0 ausentes de `C2-colores.json`/`C2-tokens.json`; 81 tokens, 0 ausentes; 52 pares token-hex y 22 pasos de escala, 0 discrepancias; insignias = `color.role.*.surface-*`; semánticos de página, tarjeta, texto, primario, enlace y foco = `ui.*` resuelto; 6 tarjetas, Detalle con Karen, T-1..T-4, 5 estados de pago, 14 secciones de U4 y orden de barra lateral de U1 presentes; patrón P3-12 = 0.
- (d) Patrón de credencial P3-01/P3-03 sobre LEEME = 0; «compartir/publicar/invitar/suscripción» solo en l. 51 y 58; todas las rutas reales citadas existen (`paths.py`). Dudas de una persona sin contexto: M-1 a M-5.
- (e) `grep -c /api/v1` en el texto: U3-dashboard 13, DESIGN 0, U4 0, C1-resumen 0; nombres de operación: U3 25, U4 89 → M-6. DESIGN.md con 6 bloques de código (B-1).
- (f) C1-resumen 4 páginas; 0 rutas y 0 operaciones; 9 módulos = 9 etiquetas de `extractos/INDICE.md`; «31 entidades» = `C1-datos.md` §7.
- (g) `salida/iter-1..3` y `final/` existen con `.gitkeep`; `salida/LEEME.md` presente.
- Criterios cubiertos: P3-01, P3-02, P3-03, P3-04 (parte del paquete), P3-05, P3-06 (preparación), P3-11, P3-12, P3-13; P3-07 a P3-10 no aplican al paquete.

## Supuestos
1. La verificación (c) «0 hex ausentes de C2» es la lectura de la Verifica de P3-02; la otra lectura (todos los hex de C2 en el prompt) se documenta en B-2 y no se usa para el veredicto.
2. «Con su LEEME» en (g) = `docs/design/salida/LEEME.md`; un LEEME en cada `iter-N` rompería P3-11 y el conteo de P3-06.
3. Las rutas de la API no son «código» en sentido estricto (P3-02 alta); por eso M-6 es medio y no alto.

## Bloqueos
Ninguno.

## Pendientes
- Corregir M-1 a M-6 (agente de corrección del paquete) y volver a correr `docs/design/tools/paquete_entrada.sh`; luego la ronda 2 de este revisor.

## Solicitudes fuera de alcance
- `docs/criterios/P3-claude-design.md`, P3-02 Verifica: fijar la lectura de «los hex de `C2-tokens.json` aparecen en `prompt.md`» (subconjunto o totalidad), porque admite dos lecturas (B-2).
- Integrador de WF-1: decidir si se versiona `docs/design/tools/paquete/.cache/` (B-5).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna (no se usó ningún valor operativo de K31).

## Consumo aproximado
- Acciones de herramienta: 37. Leído: ~225 KB (protocolo, §K, criterios P3, plan §3.6, LEEME de entrada y salida, prompt, MANIFIESTO, script, 2 reportes, C1-resumen, INDICE de extractos, extractos de U3/U4/U1/DESIGN por grep). Escrito: ~23 KB (este reporte, el informe y los scripts temporales).
