# Reporte · revisor-paquete-r2 (WF-1, ronda 2 de 3) · 11-oct-2026

## Hecho
- Lecturas por partes:
  - `protocolo-agentes.md` completo;
  - decisiones §K completa;
  - `criterios/P3-claude-design.md` completo;
  - plan §3.6 y §1.5;
  - del paquete: `entrada/LEEME.md`, `prompt.md`, `MANIFIESTO.md`, `salida/LEEME.md` y `paquete_entrada.sh`;
  - reportes `paquete-claude-design.md`, `correccion-paquete-r1.md` y `paquete-regenerar.md`;
  - `extractos/INDICE.md`;
  - texto de los 4 PDF, solo por `grep`, salvo C1-resumen, que leí entero porque tiene 4 páginas.
- Reejecuté las verificaciones (a) a (g) con comandos y scripts propios y escribí el informe `docs/estado/wf1/REVISION-paquete-r2.md`.
- Veredicto: **NO APROBADO**. 0 altos, 1 medio nuevo (H-1: el paso 7 del LEEME pide exportar imagen por pantalla, que plan §1.5 no documenta, y no dice qué hacer si solo hay PDF o HTML) y 7 bajos.
- Los 6 medios de la ronda 1 están cerrados. M-6 deja un residuo bajo: `logout` y `login` en U4-ajustes.pdf.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/REVISION-paquete-r2.md` (informe de revisión)
- `/home/deepsec/finanzas-hogar/docs/estado/wf1/revisor-paquete-r2.md` (este reporte)
- Temporales en `/home/deepsec/finanzas-hogar/datos/tmp/rev-paq-r2/`, borrados al terminar: texto de los PDF, `bin/chk_prompt.py` y dos cajas de arena de regeneración.

## Verificación
- (a) `grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md | sha256sum -c`: 21 OK, 0 fallos. Tamaños de MANIFIESTO = `stat`.
- (a) Regeneración en caja de arena:
  - Sin `--forzar`: `diff -r` contra `entrada/` = 0 diferencias.
  - Con `--forzar`: difieren en bytes solo los 4 PDF y MANIFIESTO. El texto (`pdftotext`) y `pdfinfo` son idénticos salvo las fechas. DOCX, PNG, JSON y CSS son idénticos. Son las diferencias esperadas.
- (b) `pdfinfo`: DESIGN 36, U3 19, U4 37 y C1-resumen 4 páginas. DESIGN.pdf trae los 10 encabezados, 2 veces cada uno. Hay 11 PNG de 910 × 1287 y `paleta.png` de 1400 × 3000.
- (c) Resultado de `chk_prompt.py`:
  - 90 hex, 0 ausentes de `C2-tokens.json` y 0 de `C2-colores.json`;
  - 83 tokens, 0 ausentes;
  - 52 pares token‑hex, 0 discrepancias;
  - escalas, insignias y medidas iguales a C2.
- (c) U3 (4 pestañas) y U4 (14 secciones, mismos nombres que P-26.01 a P-26.14) presentes. Hay reglas de dinero y de accesibilidad, 0 bloques de código y el texto está en es-MX. Pesa 12 284 bytes, contra 12 288 de 12 KiB (B-1).
- (d) Credenciales en LEEME = 0. Las acciones prohibidas solo aparecen en las advertencias. Todas las rutas citadas existen, salvo las plantillas y los archivos que se crean después. Dudas de una persona sin contexto: H-1 (medio) y B-4 a B-6 (bajos).
- (e) No hay archivos de código. `/api/v1` = 0 en los 4 PDF. Nombres de operación (210 de los extractos):
  - U3: 0;
  - C1-resumen: 0;
  - prompt: 0;
  - U4: 3 (B-2);
  - DESIGN: 1, en el texto PROVISIONAL de S1; es una excepción declarada y la acepto.
- (f) C1-resumen tiene 4 páginas y es coherente con `extractos/INDICE.md`: módulos, 6 tarjetas, Detalle con Karen, Calendario, búsqueda, login y 14 secciones. Las 31 entidades coinciden con C1-datos §7.
- (g) `salida/iter-1`, `iter-2`, `iter-3` y `final` existen, vacías salvo `.gitkeep`, y `salida/LEEME.md` las describe.
- P3-12: el texto de `reporte.pdf`, origen de los PNG, tiene 0 correos, 0 credenciales y 0 números de tarjeta.
- Criterios: P3-01, P3-02, P3-05, P3-11, P3-12 y P3-13 cumplen. P3-03 cumple su texto, pero la verificación (d) falla por H-1. P3-04 y P3-06 a P3-10 no aplican al paquete o quedan pendientes de Axel.

## Supuestos
1. «12 KB» se lee como 12 KiB (12 288 bytes), la unidad de MANIFIESTO. Con 12 000 bytes, `prompt.md` no cumpliría (B-1, bajo).
2. «crear proyecto» (paso 2 del LEEME) se acepta porque plan §3.6 paso 6 lo autoriza, aunque P3-03 no lo nombra.
3. `tokens.css` y los bloques de valores de DESIGN.pdf no violan la decisión 42 porque K35 los manda tal cual (B-3, bajo).
4. «con su LEEME» en (g) se lee como `salida/LEEME.md`, porque P3-11 pide las carpetas vacías.
5. P3-04 se evalúa por vuelta: no busqué el diario de herramientas del orquestador.

## Bloqueos
Ninguno.

## Pendientes
- Corregir H-1 (`entrada/LEEME.md` l. 35 y `salida/LEEME.md` l. 22) y volver a pasar el paquete por la ronda 3. Los bajos B-1 a B-7 son opcionales.

## Solicitudes fuera de alcance
- Ninguna. Las correcciones van al agente de corrección del paquete; las rutas están en el informe.

## Paquetes instalados
Ninguno.

## Propuestas para Axel
Ninguna (no hay valores de K31 en este artefacto).

## Consumo aproximado
- Acciones de herramienta: unas 43.
- KB leídos: unos 160 (texto de PDF solo por `grep`).
- KB escritos: unos 22 (informe de 17 KB, este reporte y el script temporal).
