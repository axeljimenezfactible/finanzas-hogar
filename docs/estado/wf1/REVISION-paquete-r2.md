# REVISION-paquete-r2 · Paquete de entrada para Claude Design (ronda 2 de 3)

Revisor: `revisor-paquete-r2` (independiente; no escribió nada de lo que revisa y no corrige nada). Fecha: 11-oct-2026.
Contra: `docs/criterios/P3-claude-design.md` v1.2 (lo aplicable al paquete y a las carpetas de salida), `docs/plan/plan.md` §3.6 (y §1.5, que P3 cita en su trazabilidad) y `docs/plan/decisiones.md` §K (K35).
Objeto: `docs/design/entrada/` (22 archivos), `docs/design/salida/` (LEEME y 4 carpetas), `docs/design/tools/paquete_entrada.sh`, reportes `paquete-claude-design.md`, `correccion-paquete-r1.md` y `paquete-regenerar.md`.

## Veredicto

**NO APROBADO**: 0 altos, **1 medio** (H-1, nuevo) y 7 bajos. Los 6 medios de la ronda 1 quedan cerrados; de M-6 queda un residuo bajo (B-2).
Para la puerta 3 (P3-01) no hay hallazgos altos abiertos, pero la regla de esta revisión pide 0 altos y 0 medios para aprobar.

## Tabla de criterios y verificaciones

| Criterio o verificación | Estado | Evidencia propia (comandos de este revisor) |
|---|---|---|
| (a) `sha256sum -c` sobre MANIFIESTO | Cumple | `grep -E '^[0-9a-f]{64}  ' MANIFIESTO.md \| sha256sum -c`: 21 OK, 0 fuera de OK. Tamaños de la tabla de MANIFIESTO contra `stat -c %s`: 0 diferencias. |
| (a) El script regenera sin diferencias salvo las esperadas | Cumple | Copié las fuentes a una caja de arena (`datos/tmp/rev-paq-r2/root`) y corrí la copia del script, para no tocar el paquete entregado. Sin `--forzar`: los 6 productos salen «sin cambios en las fuentes» y `diff -r` contra `docs/design/entrada` da **0 diferencias**. Con `--forzar` (otra copia): solo difieren en bytes los 4 PDF y MANIFIESTO. El texto de los 4 PDF (`pdftotext`) es idéntico y `pdfinfo` también, salvo las fechas. DOCX, PNG, JSON y CSS salen idénticos byte a byte. Esas son las diferencias esperadas: las fechas que Chromium pone en el PDF. |
| (b) Cada PDF abre | Cumple | `pdfinfo`: DESIGN 36 págs., U3-dashboard 19, U4-ajustes 37, C1-resumen 4, todos A4. |
| (b) DESIGN.pdf con los 10 encabezados | Cumple | `pdftotext -layout`: «1. Principios y tono» … «10. Tokens», cada uno 2 veces (índice y cuerpo); «Colorimetr» = 2. |
| (b) 11 PNG del reporte | Cumple | `file`: `reporte-p01..p11.png` son PNG de 910 × 1287. La regeneración desde `/home/deepsec/reporte-compra-latex/build/reporte.pdf` (11 págs.) da los mismos bytes. `paleta.png` mide 1400 × 3000. |
| (c) prompt.md: 0 hex o tokens ausentes de C2 | Cumple | Script propio `chk_prompt.py` (`python3 -I`): 90 hex distintos, 0 ausentes de `C2-tokens.json` y 0 de `docs/contratos/C2-colores.json`. 83 tokens citados, 0 ausentes. 52 pares token‑hex, 0 discrepancias. Escalas de marca y neutro, 11 + 11 pasos sin discrepancias. Las 8 insignias = `color.role.*.surface-light/dark`. Medidas `space.*`, `radius.*`, `font.size.*` y títulos 20/24/32 iguales a C2. El mapeo de superficies, texto, divisor, botón primario, enlace y foco coincide con `ui.*` resuelto. |
| (c) Pantallas de U3 y 14 secciones de U4 | Cumple | Resumen, Detalle, Tendencias y Calendario (U3 §2–§5). Las 14 secciones con los mismos nombres que `U4-ajustes-a/b.md` (P-26.01 a P-26.14). Barra lateral y menú «/». 20 pantallas, 84 artboards. |
| (c) Reglas de dinero, accesibilidad, es-MX, sin código | Cumple | §«Reglas de dinero, fechas e idioma» (l. 87-91) y §«Accesibilidad (WCAG 2.2 AA)» (l. 93-95). 0 bloques de código. Español de México con trato de tú. |
| (c) ≤ 12 KB | Cumple al límite | 12 284 bytes ≤ 12 288 (12 KiB, la unidad de MANIFIESTO); ver B-1. |
| (d) LEEME para una persona sin contexto técnico | **No cumple** | Una duda que obliga a adivinar: H-1 (medio). Dudas menores: B-4 a B-6 (bajos). |
| (d) Nombres de carpetas y archivos exactos y existentes | Cumple | Script propio sobre los nombres entre comillas invertidas del LEEME: existen todos los archivos del paquete, `salida/`, `iter-1/`, `iter-2/` y `final/`. No existen todavía, y es lo esperado, `iter-N/` (plantilla), `siguiente-prompt.md`, `P3-iter-1.md` y los nombres de ejemplo de artboards. `DESIGN.md` existe en `docs/design/`. |
| (d) Sin contraseñas, credenciales ni acciones fuera de la lista | Cumple | Patrón de credencial de P3-01/P3-03 sobre LEEME = 0. «compartas/publiques/invites/suscripción» solo aparecen en las advertencias (l. 55 y 62). «Crea un proyecto nuevo» (paso 2) es la acción «crear proyecto» que autoriza plan §3.6 paso 6; la acepto con esa justificación. |
| (e) Sin archivos de código ni rutas de la API | Cumple, con residuo bajo | No hay `.go`, `.py`, `.ts`, `.tsx`, `.sql` ni `.sh`. `/api/v1` y «MÉTODO /ruta» = 0 en los 4 PDF, en prompt, en LEEME, en JSON y en CSS. Contra las 210 operaciones de `extractos/*.yaml`: U3 0, C1-resumen 0, prompt 0. U4 conserva `logout` (2) y `login` (1) (B-2). DESIGN.pdf conserva `updateNotification` (1), en el texto PROVISIONAL del bloqueo S1; es una excepción declarada en el script y la acepto. `tokens.css` es CSS por mandato de K35 (B-3). |
| (f) C1-resumen ≤ 10 páginas y coherente con extractos/INDICE.md | Cumple | 4 páginas. Los módulos coinciden con las etiquetas del índice: Inicio, Gastos, Presupuesto, Comidas, Compras y tiendas, Créditos, Salud, Reportes y Ajustes. Las 6 tarjetas corresponden a los 6 `getDashboard*`, Detalle con Karen a `getPersonParticipation`, Calendario a `listCalendarEvents`, Ctrl+K a `search` y P-00 a `login`. Las 14 secciones caben en `ajustes-a..e`. Las «31 entidades» = C1-datos §7. Sin rutas ni nombres de operación. |
| (g) salida/iter-1, iter-2, iter-3 y final con su LEEME | Cumple | `ls -la`: las 4 carpetas existen, vacías salvo `.gitkeep`, como pide P3-11. `docs/design/salida/LEEME.md` las describe con sus rutas absolutas. |
| P3-01 (alta) puerta 3 | Cumple en lo que toca al paquete | 22 elementos, `sha256sum -c` en 0, credenciales 0. Este informe es el del revisor, con 0 altos. |
| P3-02 (alta) paquete completo | Cumple | `ls -1A` contra la lista exacta de P3-02: `diff` vacío; 22 elementos, 0 de 0 bytes. 4 PDF con `Pages:` ≥ 1. «Colorimetr» ≥ 1. `grep -c color prompt.md` ≥ 1. Hex de C2 en prompt. MANIFIESTO con 21 sha256. Decisión 42: ver B-3. |
| P3-03 (alta) LEEME ejecutable por Axel | Cumple el texto del criterio; falla (d) por H-1 | Pasos numerados en es-MX: abrir con su sesión (1), crear proyecto (2), pegar (3), subir con ruta absoluta (4), enviar (5), exportar (7), descargar y guardar en `salida/iter-N/` (8) y `final/` (Vueltas 5). |
| P3-04 (alta) | No aplica al paquete | Lo evalúa cada `P3-iter-N.md`. LEEME l. 63 y prompt l. 106 dicen que lo descargado es dato. Este revisor hizo 0 llamadas a herramientas de navegador. |
| P3-05 (alta) cuándo parar | Cumple | «Si algo falla» (l. 51-57): «2 o 3 veces, detente», límite de suscripción, cambio de herramienta y «el proyecto sigue con DESIGN.md». `grep -c DESIGN.md` = 3. |
| P3-06 (alta) | Pendiente de Axel | Carpetas listas; 84 artboards derivados de prompt l. 12; aviso «iter-N lista» explicado. |
| P3-07 a P3-10 (alta) | No aplican al paquete | Se evalúan por vuelta. |
| P3-11 (media) carpetas | Cumple | LEEME l. 9, 36 y 49 y MANIFIESTO l. 5-7 con rutas absolutas. `grep -c "docs/design/salida/iter-N" LEEME.md` = 1. `ls -d` de `salida/`, `salida/final/` y `entrada/`: existen. |
| P3-12 (media) datos y credenciales | Cumple | Credenciales en prompt = 0; en `C2-tokens.json` y `tokens.css` = 0 y sin URL. El texto de `reporte.pdf`, origen de los 11 PNG, tiene 0 correos, 0 patrones de credencial y 0 números de tarjeta de 16 dígitos. |
| P3-13 (media) supuestos y bloqueos | Cumple | Los reportes del paquete traen Supuestos. `ls docs/estado/bloqueos/P3-*` = ninguno, y no hace falta ninguno. |

## Hallazgos

### H-1 (media) · P3-03 / P3-06 / verificación (d): el paso 7 pide una salida que la herramienta no documenta y no dice qué hacer si no existe
- Ruta: `/home/deepsec/finanzas-hogar/docs/design/entrada/LEEME.md` (l. 35) y `/home/deepsec/finanzas-hogar/docs/design/salida/LEEME.md` (l. 22).
- Qué pasa: el paso 7 dice «Exporta cada pantalla como imagen. Se esperan hasta 84 imágenes… Si la herramienta ofrece un solo exporte con todo, exporta **también** en HTML o en PDF». Hay tres problemas:
  - Las salidas documentadas de Claude Design (plan §1.5) son URL, Canva, PDF, PPTX, HTML independiente y paquete de handoff. Ninguna es una imagen por pantalla.
  - `prompt.md` («Qué entregar», l. 98) pide «artboard… listo para exportar como imagen, **o** un solo exporte HTML o PDF con todos».
  - `salida/LEEME.md` (l. 23) acepta que «un PDF de N páginas cuenta N pantallas».
- La duda que queda: si el menú de exportar solo ofrece PDF, HTML o PPTX, ¿basta con el PDF, hay que detenerse («si la herramienta cambia… de opciones, detente», l. 55) o hay que tomar capturas, que no están entre las acciones autorizadas? Según la opción que elija, la Pista 3 se detiene sin necesidad o se sale de las acciones permitidas. Además, P3-06 cuenta lo descargado contra los artboards, así que la duda también cambia la evidencia.
- Corrección sugerida: reescribir el paso 7 así: «Exporta en PDF (un archivo con todas las pantallas) y, si se puede, también en HTML independiente. Si además deja bajar cada pantalla como imagen, también sirve. No uses URL, Canva ni PPTX sin avisarme».
  - En `salida/LEEME.md` l. 22, cambiar «84 imágenes» por «84 artboards (imágenes o páginas de un PDF)».
  - Ajustar el ejemplo `.png` del LEEME l. 41 para que no suponga solo imágenes.

### B-1 (baja) · verificación (c): prompt.md queda 4 bytes bajo el tope y hay dos lecturas del tope
- Ruta: `/home/deepsec/finanzas-hogar/docs/design/entrada/prompt.md` y `/home/deepsec/finanzas-hogar/docs/estado/wf1/paquete-claude-design.md` (Supuesto 2).
- Qué pasa: pesa 12 284 bytes. Cumple con 12 KiB (12 288), la unidad de MANIFIESTO y del Supuesto 4 de `correccion-paquete-r1.md`. No cumple con la lectura «menos de 12 000 bytes» del Supuesto 2 del reporte original.
- Corrección sugerida: recortar unos 300 bytes, por ejemplo con frases repetidas entre l. 12 y l. 98, para que cumpla con las dos lecturas, o fijar la unidad en un solo lugar.

### B-2 (baja) · M-6 residual / verificación (e): U4-ajustes.pdf conserva 3 nombres de operación y nombres de campo de la API
- Ruta: `/home/deepsec/finanzas-hogar/docs/design/tools/paquete_entrada.sh` (función `filtrar_api` y `sin_operaciones`).
- Qué pasa:
  - El filtro solo sustituye nombres con prefijo en camelCase, así que se escapan `logout` (U4-ajustes-a l. 94 y l. 372) y `login` (l. 374). En el PDF se leen «— / logout», «…operación de C1, logout , consulta de C1…» y «login (P-00)».
  - La comprobación `sin_operaciones` tampoco los detecta, y `correccion-paquete-r1.md` informa «U4 0/0».
  - Siguen nombres de campo de la API (`reset_parameters`, `is_initial`, `all_except_current`, `user_account_id`).
  - Las listas «Cobertura de operaciones» quedan como «consulta de C1, consulta de C1, operación de C1…», que no aportan nada para maquetar.
- Corrección sugerida:
  - Que el filtro tome la lista exacta de nombres de `docs/contratos/extractos/INDICE.md` (210) en lugar del prefijo, y que `sin_operaciones` use esa misma lista.
  - Opcional: quitar del PDF las secciones «Matriz de operaciones» y «Cobertura de operaciones», que son para quien implementa la API.

### B-3 (baja) · decisión 42 / P3-02: «Ningún archivo del paquete contiene código» no es exacto (B-1 de la ronda 1, sigue abierto)
- Ruta: `/home/deepsec/finanzas-hogar/docs/design/entrada/MANIFIESTO.md` (l. 8, la genera `paquete_entrada.sh`) y `/home/deepsec/finanzas-hogar/docs/design/entrada/LEEME.md` (l. 7).
- Qué pasa: `tokens.css` es CSS, por mandato de K35. Además, DESIGN.pdf trae:
  - bloques CSS (`@font-face`, `:root {`, sombras con `color-mix(… var(--color-neutral-950) …)`, `@media (prefers-reduced-motion: reduce)`);
  - un bloque JSON de tokens;
  - instrucciones de Tailwind (`@import "tailwindcss"`, `@theme inline`, `@custom-variant coarse`);
  - 21 rutas del repositorio y nombres de scripts (`colorimetria.py`, `generar_c2_tokens.py`, `verificar_c2.py`…).
- Es bajo porque K35 manda DESIGN.pdf como «DESIGN.md convertido» y `tokens.css` tal cual: son valores de estilo, no código de la aplicación.
- Corrección sugerida: cambiar la frase a «sin código de la aplicación; `tokens.css` y los bloques de valores de DESIGN son solo referencia de estilo». Opcionalmente, quitar de la copia que va a pandoc los bloques de código y las instrucciones de Tailwind, igual que `filtrar_api`.

### B-4 (baja) · verificación (d): no dice cómo se entera Axel de que `siguiente-prompt.md` ya está listo
- Ruta: `/home/deepsec/finanzas-hogar/docs/design/entrada/LEEME.md` (Vueltas, pasos 2-3).
- Qué pasa: dice que un agente «deja dos archivos», pero no que el aviso de que ya están llegará en la misma sesión de Claude Code. No bloquea, porque el canal ya está definido en la l. 5.
- Corrección sugerida: agregar «te aviso en la sesión cuando estén listos».

### B-5 (baja) · verificación (d): la regla general de «Si algo falla» va antes de sus excepciones
- Ruta: `/home/deepsec/finanzas-hogar/docs/design/entrada/LEEME.md` (l. 55-56).
- Qué pasa: la l. 55 dice «si no acepta un archivo, detente». La l. 56 dice que se siga sin `C2-tokens.json` o `tokens.css` y que se pruebe `DESIGN.docx`.
- Corrección sugerida: «si no acepta un archivo, detente (salvo los casos del punto siguiente)».

### B-6 (baja) · verificación (d): «La lista de nombres está también en salida/LEEME.md» no es exacto
- Ruta: `/home/deepsec/finanzas-hogar/docs/design/entrada/LEEME.md` (l. 41).
- Qué pasa: `salida/LEEME.md` trae la convención y ejemplos, no la lista. La lista numerada de las 20 pantallas, que sirve para decir «cuáles faltan» (paso 9), está en `prompt.md` §«Pantallas a maquetar».
- Corrección sugerida: «La convención está en `salida/LEEME.md` y la lista numerada de pantallas en `prompt.md`».

### B-7 (baja) · trazabilidad: el reporte original del paquete quedó desactualizado
- Ruta: `/home/deepsec/finanzas-hogar/docs/estado/wf1/paquete-claude-design.md`.
- Qué pasa: todavía dice:
  - 80 artboards;
  - `prompt.md` de 11 958 bytes;
  - U4 de 30 páginas;
  - LEEME «con `sha256sum -c`»;
  - Supuesto 5: U4 nombra operaciones.
- La evidencia vigente está en `correccion-paquete-r1.md` y en este informe.
- Corrección sugerida: el integrador cita como evidencia de P3-01/P3-02 la corrección r1 y esta revisión, o agrega una nota «sustituido por» al reporte original.

## Estado de los hallazgos de la ronda 1

| Hallazgo r1 | Estado | Evidencia |
|---|---|---|
| M-1 · dónde y a quién avisar «iter-N lista» | **Cerrado** | LEEME l. 5 «Dónde me avisas»: en la sesión de Claude Code del proyecto; lo recibe Claude Code. Lo repiten l. 34, 45 y 53. |
| M-2 · cuántos archivos y entrega parcial | **Cerrado** | l. 35 (hasta 84) y paso 9 (l. 37): guardar lo que haya, avisar con la lista de faltantes y no pedir el resto. `salida/LEEME.md` l. 22 dice lo mismo. H-1 es otra cosa: el formato de exporte, no el conteo. |
| M-3 · preguntas o plan en lugar de maquetas | **Cerrado** | Paso 6 (l. 34) y l. 54: no contestar, copiar y avisar; el orquestador da el mensaje que hay que pegar. |
| M-4 · dos convenciones de nombre, la 19 y «menu-barra» | **Cerrado** | Una sola convención `NN-pantalla-tema-dispositivo` en prompt l. 12, LEEME l. 41 y `salida/LEEME.md` l. 18. «No renombres nada». La 19 tiene dos nombres, expandida y colapsada. «menu-barra» pasó a `20-menu-comandos` con la explicación de la pantalla 20. |
| M-5 · comando de terminal en «Comprobar» | **Cerrado** | l. 66-68: «Tú no tienes que comprobar nada»; ya no hay comandos. |
| M-6 · rutas `/api/v1` y nombres de operación en U3/U4 | **Cerrado**; residuo bajo en B-2 | `/api/v1` = 0 en los 4 PDF. Nombres de operación (210 de los extractos): U3 0 (antes 25), U4 3 (antes 89: `logout` ×2 y `login` ×1). |

## Comandos y scripts (todos en `datos/tmp/rev-paq-r2/`, borrados al terminar)
- `sha256sum -c`, `ls -1A` + `diff` contra la lista de P3-02, `stat` contra los tamaños de MANIFIESTO, `pdfinfo` y `pdftotext -layout` de los 4 PDF, `file` de los PNG.
- `bin/chk_prompt.py`: hex, tokens, pares token‑hex, escalas, insignias y medidas de `prompt.md` contra `C2-tokens.json` y `C2-colores.json`.
- Lista de 210 `operationId` sacada de `docs/contratos/extractos/*.yaml` y buscada con `grep -w -F -f` en el texto de los 4 PDF, del DOCX, de prompt, de LEEME, de JSON y de CSS.
- Regeneración en dos cajas de arena con copia de las fuentes y del caché (`root` sin `--forzar` y `root2` con `--forzar`), más `diff -r`, `cmp` y comparación de `pdftotext`.
- Comprobación de existencia de cada ruta o archivo citado entre comillas invertidas en los dos LEEME.
- `pdftotext` de `/home/deepsec/reporte-compra-latex/build/reporte.pdf`, para buscar correos, credenciales y números de tarjeta.
