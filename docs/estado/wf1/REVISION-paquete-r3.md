# REVISION-paquete-r3 · Paquete de entrada para Claude Design (WF-1, ronda 3 de 3)

Fecha: 11-oct-2026 · Revisor: `revisor-paquete-r3` (independiente; no escribí nada de lo revisado y no corregí nada) · Alcance: `docs/criterios/P3-claude-design.md` 1.2 (lo aplicable al paquete y a las carpetas de salida), plan §1.5 y §3.6, decisiones §K (K35). Todas las cifras de este informe salen de mis comandos y de mi script en `datos/tmp/rev-paq-r3/` (borrado al terminar). La única corrida que escribe en el paquete fue `paquete_entrada.sh` sin `--forzar`, con respaldo previo de `entrada/`, `salida/` y `.cache/`; `diff -rq` contra el respaldo salió vacío, así que no cambió nada.

## Veredicto
**NO APROBADO**: 0 altos, **1 medio** (H-1, una duda en el paso 2 del LEEME) y 5 bajos. El hallazgo medio de la ronda 2 (formato del exporte en el paso 7) queda **cerrado**. La corrección de H-1 es una frase en `LEEME.md` más regenerar `MANIFIESTO.md`. Ningún criterio alto falla: P3-01, P3-02, P3-03 y P3-05 cumplen en su texto, y P3-06 queda preparado y «pendiente de Axel».

## Verificaciones pedidas (a)–(g)
| # | Estado | Evidencia propia |
|---|---|---|
| (a) sha256 y regeneración | Cumple | `grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md \| sha256sum -c`: 21 OK y 0 fallos. Tabla del MANIFIESTO igual a las líneas de sha256sum (`diff` vacío); tamaños de la tabla iguales a `stat`; 22 elementos, 0 de 0 bytes; fuera de la lista solo queda `MANIFIESTO.md`. El script sin `--forzar` corrió los 5 pasos: «sin cambios en las fuentes» en 4 PDF, DOCX y PNG; «sha256sum -c: 0 diferencias»; «Listo: 22 archivos». `diff -rq` contra el respaldo (entrada y caché) salió vacío. `cmp`: `C2-tokens.json`, `tokens.css` y `paleta.png` son idénticos a sus fuentes. |
| (b) PDF, encabezados y PNG | Cumple | `pdfinfo`: DESIGN 36, U3 19, U4 37 y C1-resumen 4 páginas; A4; ninguno cifrado. Índice de DESIGN.pdf (pdftotext, l. 2-11): «1. Principios y tono» … «10. Tokens», los mismos 10 `## N.` de DESIGN.md; «Colorimetr» aparece 2 veces. DESIGN.docx: Word 2007+ con 16 partes. 11 `reporte-pNN.png` de 910 × 1287; el reporte de origen tiene 11 páginas. |
| (c) prompt.md | Cumple (bajos B-2 y B-1) | Script propio `prompt_check.py`. Hex: 90 únicos y 0 fuera de `C2-colores.json`. Tokens: 81 nombres citados y 0 fuera de `C2-tokens.json` (los `color.*`: 0 fuera de C2-colores). 52 pares token→hex sin discrepancias. Escalas de marca y neutros: 11 + 11 pasos, 0 discrepancias. Las 16 insignias son `surface-light`/`surface-dark` de su rol. `space.*`, `radius.*` y `font.size.*` coinciden con C2. Pantallas: Resumen con 6 tarjetas en el orden de U3 §2; Detalle con «Participación de Karen»; Tendencias con T-1 a T-4; Calendario; las 14 secciones con los nombres de U4 (P-26.01 a P-26.14), barra lateral y menú «/». Dinero `$1,258.68`, fechas, porcentajes, WCAG 2.2 AA (4.5:1, 3:1, 44 px). es-MX y trato de tú. Pesa 12 284 bytes. 0 bloques de código y 0 llaves; `<paso>` es un marcador de texto. |
| (d) LEEME | Cumple con 1 duda media (H-1) | Patrón de credencial de P3-01/P3-03: 0. «compart/publi/invit/suscripci» solo aparecen en advertencias (l. 55, 62). Rutas: existen todas las de entrada y salida; las que faltan son archivos futuros (`siguiente-prompt.md`, `P3-iter-N.md`) o nombres de ejemplo. Acciones que pide: abrir, crear proyecto (permitido por plan §3.6 paso 6), pegar, subir, enviar, esperar, exportar, descargar y guardar, y avisar en la sesión. Dudas: ver H-1 (media), B-4 y B-5. |
| (e) Código y rutas de la API | Cumple (residuo bajo B-3) | 0 archivos `.go/.ts/.tsx/.js/.jsx/.py/.sql/.sh/.html/.yaml`. En los 4 PDF y el DOCX: `/api/` 0, «MÉTODO /» 0, `operationId` 0, SQL 0. Nombres de operación: U3 0, C1 0, U4 0 con prefijo (con `login`/`logout` sueltos: 3, B-2 de r2) y DESIGN 1 (`updateNotification` en el texto PROVISIONAL del bloqueo S1, previsto en el script). |
| (f) C1-resumen | Cumple (bajo B-1) | 4 páginas (≤ 10). «C1 1.2» coincide con «C1 v1.2.0» de `extractos/INDICE.md`. Sus módulos (§2) son las 9 etiquetas del índice (Inicio, Gastos, Presupuesto, Comidas, Compras y tiendas, Créditos, Salud, Reportes, Ajustes), más Buscar/Registrar y «Fuera de la barra». 14 secciones iguales a U4. 0 rutas y 0 operaciones. |
| (g) Carpetas de salida | Cumple | `salida/iter-1`, `iter-2`, `iter-3` y `final` existen, vacías salvo `.gitkeep`; `ls -d salida/iter-* \| wc -l` = 3. `salida/LEEME.md` describe las 4 (l. 7-14). No hay un LEEME por carpeta: leo «su LEEME» como `salida/LEEME.md`, porque P3-11 pide carpetas vacías. |

## Criterios P3
| Criterio | Estado | Evidencia |
|---|---|---|
| P3-01 (alta) paquete listo y revisado | Cumple para el paquete | 22 elementos; sha256 con 0 diferencias; credenciales 0 en LEEME; este informe no tiene altos (el medio H-1 impide «aprobado», no la regla de altos de P3-01). |
| P3-02 (alta) paquete completo | Cumple | Lista exacta de P3-02 = `ls -1A`. 4 PDF con Pages ≥ 1. «Colorimetr» ≥ 1. `grep -c color prompt.md` ≥ 1. Hex de C2 en el prompt. 21 sha256. |
| P3-03 (alta) LEEME ejecutable | Cumple el texto; duda H-1 | Pasos numerados en es-MX: abrir con su sesión, subir con ruta absoluta (l. 32), pegar, enviar, exportar, descargar y guardar en `salida/iter-N/` (l. 36) y en `final/` (l. 49). Sin credenciales. Ensayo en seco: todas las rutas existen. |
| P3-04 (alta) | No aplica al paquete | LEEME l. 63 y prompt l. 106: lo descargado es dato. Este revisor hizo 0 acciones de navegador. El conteo en el registro del orquestador se hace en cada `P3-iter-N.md`. |
| P3-05 (alta) cuándo parar | Cumple | l. 51-57: «2 o 3 veces», «detente», «el proyecto sigue con DESIGN.md»; `grep -c DESIGN.md` = 3. Sin bucles: «no reenvíes» (l. 34) y «No sigas intentando» (l. 53). |
| P3-06 (alta) | Preparado; pendiente de Axel | 3 carpetas `iter`; 0 `P3-iter-*.md`. «Un PDF de N páginas cuenta N» (LEEME l. 35, salida l. 22-23) coincide con el criterio. |
| P3-07 a P3-10 (alta) | No aplican al paquete | Se evalúan en cada vuelta. |
| P3-11 (media) | Cumple | Rutas absolutas en LEEME (l. 9, 36, 49) y MANIFIESTO (l. 5-7). `grep -c "docs/design/salida/iter-N"` = 1. `ls -d` de entrada, salida y final correcto. 0 menciones de `file_upload` o de carpeta compartida. |
| P3-12 (media) | Cumple | Credenciales 0 en el prompt. 0 correos en los 4 PDF, DOCX, JSON, CSS, prompt y LEEME. El texto del reporte de origen no tiene correos, números de 12 o más dígitos, RFC, CURP ni CLABE. La única coincidencia del patrón está en `U4-ajustes.pdf`: «tras cambiar la contraseña: «Contraseña cambiada…»». Es un mensaje de la interfaz, sin valor secreto, y así lo justifico. |
| P3-13 (media) | Aún no aplica | 0 bloqueos `P3-*`. |

## Hallazgos
### H-1 (media) · P3-03 / verificación (d): el paso 2 no dice qué elegir si al crear el proyecto la herramienta pide tipo de proyecto o sistema de diseño
- Ruta: `/home/deepsec/finanzas-hogar/docs/design/entrada/LEEME.md` (l. 30; relacionada con l. 55).
- Qué pasa: el paso 2 solo dice «Crea un proyecto nuevo y ponle el nombre "Finanzas Hogar"». Si al crearlo aparece una elección (tipo de proyecto, nivel de detalle o un sistema de diseño guardado en la cuenta o en la organización), el LEEME no dice qué elegir. Tampoco dice si eso cuenta como «cambia de opciones», que según la l. 55 obliga a detenerse; para quien nunca vio la herramienta no hay referencia de qué es un «cambio». Además, elegir un sistema de diseño ajeno contradice «No cambies la paleta» (prompt l. 102), la decisión 43 y P3-08 (paleta de C2; Notion como única referencia). No pude ver esa pantalla (K35: el orquestador no abre claude.ai/design), así que el hallazgo es la duda que deja el texto, no un dato de la herramienta.
- Corrección sugerida: añadir al paso 2: «Si te pide elegir un tipo de proyecto, elige el de prototipo o maqueta. Si te ofrece un sistema de diseño, no elijas ninguno: los colores y medidas van en el mensaje. Si no puedes seguir sin elegir algo que no está aquí, detente y avísame.» Después, regenerar `MANIFIESTO.md` con `paquete_entrada.sh` sin `--forzar`.

### B-1 (baja) · coherencia prompt/C1-resumen ↔ U3: el Calendario omite tipos de evento
- Rutas: `/home/deepsec/finanzas-hogar/docs/design/entrada/prompt.md` l. 18 y `/home/deepsec/finanzas-hogar/docs/design/tools/paquete/C1-resumen.md` §5.
- Qué pasa: U3 §5 lista también «fechas meta de ahorro (`savings_goal`)» y «envío del reporte mensual (`monthly_report`)»; el prompt omite además «venideros». `U3-dashboard.pdf`, que también se sube, sí los trae. Lo dejo en baja porque es una omisión del resumen, no una contradicción, y la maqueta solo es referencia (decisión 43). Ya se había señalado como B-3 en r1.
- Corrección sugerida: agregar los dos tipos al prompt cuando se pueda (ver B-2 por el tamaño) y a C1-resumen §5.

### B-2 (baja) · tamaño de prompt.md
- Qué pasa: pesa 12 284 bytes. Cumple con 12 KiB (12 288) y no con 12 000 bytes, la lectura del Supuesto 2 del reporte original. Le quedan 4 bytes: cualquier cambio al prompt (B-1) lo pasa. Sigue B-1 de r2.

### B-3 (baja) · «No contiene código» no es exacto (B-3 de r2, sigue abierto)
- Rutas: `LEEME.md` l. 7 y `MANIFIESTO.md` l. 8 (la escribe el script).
- Qué pasa: DESIGN.pdf trae bloques CSS (`@font-face`, `@media (prefers-reduced-motion…)`) y un bloque JSON de DESIGN.md, y `tokens.css` es CSS. Es bajo porque la decisión 42 incluye «DESIGN.md + tokens»: son valores de estilo, no código de la aplicación.
- Corrección sugerida: «sin código de la aplicación».

### B-4 (baja) · exportes finales
- Ruta: `LEEME.md` l. 49.
- Qué pasa: pide «(HTML y PDF)» sin repetir el «con el PDF basta» del paso 7, y no dice qué hacer si la versión aceptada es una vuelta anterior a la última: copiar de `iter-N/` o volver a exportar.
- Corrección sugerida: «si solo hay PDF, basta; si aceptas una vuelta anterior, avísame».

### B-5 (baja) · bajos de r2 sin cambio (fuera del encargo de `correccion-paquete-r2`)
- B-2 r2: `login` y `logout` siguen en U4-ajustes.pdf (3 apariciones).
- B-4 r2: no dice que el aviso de `siguiente-prompt.md` llega en la sesión.
- B-5 r2: a la l. 55 le falta «salvo los casos del punto siguiente».
- B-6 r2: la l. 41 sigue diciendo «La lista de nombres está también en…».
- B-7 r2: `paquete-claude-design.md` sigue con «80 artboards». La evidencia vigente está en las correcciones r1 y r2 y en este informe.

## Estado de los hallazgos de la ronda 2
| Hallazgo r2 | Estado | Evidencia |
|---|---|---|
| H-1 (media) · paso 7: salida no documentada y sin indicar qué hacer | **Cerrado** | l. 35: «Lo normal es PDF… y, si se puede, también HTML independiente»; las imágenes son opcionales; «con el PDF basta, no tomes capturas de pantalla y no te detengas»; «No uses URL, Canva ni PPTX sin avisarme»; detenerse solo «si no ofrece ni PDF ni HTML ni imágenes». La l. 55 aclara que no ofrecer imágenes no cuenta como cambio. `salida/LEEME.md` l. 22-23: «84 artboards (imágenes o páginas de un PDF)», «Un PDF de N páginas cuenta N». El prompt (l. 98) acepta «un solo exporte HTML o PDF». Coincide con plan §1.5 y P3-06. |
| B-1 a B-7 (bajos) | Abiertos | Ver B-2, B-3 y B-5 de este informe. |

## Comandos y scripts
`sha256sum -c`; `comm`/`diff` de la tabla del MANIFIESTO contra disco; respaldo y `bash docs/design/tools/paquete_entrada.sh` + `diff -rq`; `cmp` con las fuentes; `pdfinfo`, `pdftotext -layout` y `unzip -p` del DOCX; `datos/tmp/rev-paq-r3/scripts/prompt_check.py` (python3 -I, solo biblioteca estándar); `grep` de patrones de credencial, correos, rutas de la API, operaciones y SQL; `ls -d`/`ls -A` de salida. Todo en `datos/tmp/rev-paq-r3/`, borrado al terminar.
