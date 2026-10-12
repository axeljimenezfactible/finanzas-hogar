# LEEME · Paquete para Claude Design (Finanzas Hogar)

Para Axel. Esta carpeta es todo lo que necesitas para pedirle a Claude Design las maquetas de las pantallas de Finanzas Hogar. Tú haces los pasos en tu navegador con tu sesión; yo (Claude Code, el asistente de la sesión del proyecto) no abro esa página.

**Dónde me avisas.** Todos los avisos de esta guía («iter-1 lista», una pregunta de Claude Design, algo que falla) los escribes en la sesión de Claude Code del proyecto Finanzas Hogar, la misma terminal donde trabajas conmigo desde `/home/deepsec/finanzas-hogar`. Ahí los recibo yo; no hay que mandarlos por otro canal ni escribirlos dentro de Claude Design.

**Qué es el paquete.** Documentos, imágenes y un mensaje listo para pegar. **No contiene código**: los PDF y el Word explican cómo debe verse la aplicación, los dos archivos de valores (`C2-tokens.json` y `tokens.css`) son una lista de colores, tamaños y espacios, y las 11 imágenes son las páginas del reporte de compras que ya conoces. Lo que Claude Design te devuelva es una **referencia visual**: sirve para comparar y decidir, pero la aplicación se construye con `DESIGN.md` y no depende de ese servicio.

Carpeta de este paquete: `/home/deepsec/finanzas-hogar/docs/design/entrada/`

## Qué hay y en qué orden subirlo

| Orden | Archivo | Para qué sirve |
|---|---|---|
| 1 | `prompt.md` | El mensaje que pegas en el chat. Dice qué pantallas hacer, los colores y las reglas. No se sube: se abre, se copia y se pega. |
| 2 | `DESIGN.pdf` | El sistema de diseño completo (colores, letras, tamaños, componentes, accesibilidad). |
| 3 | `U3-dashboard.pdf` | Qué lleva el Inicio: las 6 tarjetas, Detalle, Tendencias y Calendario. |
| 4 | `U4-ajustes.pdf` | Las 14 secciones de Ajustes. |
| 5 | `C1-resumen.pdf` | Resumen sin código de qué datos maneja la aplicación y cómo se muestran. |
| 6 | `C2-tokens.json` | Todos los valores de diseño en un solo archivo. |
| 7 | `tokens.css` | Los mismos valores en otro formato, solo para consulta. |
| 8 | `paleta.png` | Imagen de la paleta de color en claro y oscuro. |
| 9 | `reporte-p01.png` a `reporte-p11.png` | Las 11 páginas del reporte previo: de aquí salen las cifras de ejemplo. |
| Respaldo | `DESIGN.docx` | El mismo `DESIGN.pdf` en Word. Solo si la herramienta no acepta el PDF. |
| No se sube | `LEEME.md`, `MANIFIESTO.md` | Esta guía y la lista de archivos con su huella (`sha256`) para comprobar que nada cambió. |

## Pasos

1. Abre `claude.ai/design` en tu navegador, con tu sesión de siempre.
2. Crea un proyecto nuevo y ponle el nombre «Finanzas Hogar».
3. Abre `/home/deepsec/finanzas-hogar/docs/design/entrada/prompt.md`, copia todo su texto y pégalo en el cuadro de mensaje (todavía no lo envíes).
4. Sube los archivos de las filas 2 a 9 de la tabla (7 archivos más las 11 imágenes del reporte), desde `/home/deepsec/finanzas-hogar/docs/design/entrada/`. Si el selector deja elegir varios a la vez, súbelos juntos.
5. Revisa que se vean 18 archivos adjuntos y envía el mensaje.
6. Espera a que termine. Puede tardar varios minutos; no reenvíes el mensaje mientras tanto. Si Claude Design te hace preguntas o te propone un plan antes de dibujar, **no contestes**: cópialas tal cual, avísame en la sesión de Claude Code y yo te doy el mensaje que debes pegar.
7. Exporta. Lo normal es **PDF** (un solo archivo con todas las pantallas) y, si se puede, también **HTML independiente**. Si además la herramienta deja bajar cada pantalla como imagen, también sirve: bájalas. Se esperan hasta **84 artboards** (20 pantallas, 21 con los dos estados de la barra lateral, × 4 combinaciones: claro y oscuro, móvil y escritorio); en un PDF cada página cuenta como un artboard. **Si la herramienta no ofrece imágenes, no es un problema:** con el PDF basta, no tomes capturas de pantalla y no te detengas. No uses URL, Canva ni PPTX sin avisarme. Si no ofrece ni PDF ni HTML ni imágenes, detente y avísame.
8. Descarga lo que exportaste y guárdalo en `/home/deepsec/finanzas-hogar/docs/design/salida/iter-1/` (las vueltas siguientes van en `/home/deepsec/finanzas-hogar/docs/design/salida/iter-N/`, con N = 2 o 3).
9. **Si Claude Design entrega solo una parte** (el mensaje le permite empezar por las pantallas 1 a 8 si se le acaba el límite), guarda lo que haya, avísame «iter-1 lista» y di cuáles pantallas faltan. No le pidas el resto por tu cuenta: yo decido en la revisión si hace falta otra vuelta.

## Cómo nombrar lo que descargues

**No renombres nada.** Los artboards llevan el nombre que pide el mensaje, uno por pantalla y combinación, con esta forma: `NN-pantalla-tema-dispositivo` (NN es el número de la pantalla en el mensaje; tema es `claro` u `oscuro`; dispositivo es `movil` o `escritorio`). Por ejemplo `01-resumen-claro-movil`, `13-ajustes-notificaciones-oscuro-escritorio`, y la barra lateral tiene un nombre por estado: `19-barra-lateral-expandida-claro-escritorio` y `19-barra-lateral-colapsada-claro-escritorio`. Si la herramienta baja imágenes, cada nombre lleva su extensión (por ejemplo `01-resumen-claro-movil.png`); si baja un PDF, esos nombres son los de las páginas dentro del archivo y el archivo conserva el nombre que le dé la herramienta. Guarda cada archivo con el nombre con que lo entregue la herramienta, sea el que sea; el revisor lo asigna a su pantalla. Los exportes en HTML o PDF se guardan igual, **tal como los entregue Claude Design**. La lista de nombres está también en `/home/deepsec/finanzas-hogar/docs/design/salida/LEEME.md`.

## Vueltas (hasta 3)

1. Cuando tengas todo en `iter-1/`, escríbeme en la sesión de Claude Code estas palabras: **«iter-1 lista»** (si faltan pantallas, agrega la lista de las que faltan).
2. Un agente revisa lo que bajaste contra `DESIGN.md` y deja dos archivos: el informe de la vuelta (`/home/deepsec/finanzas-hogar/docs/estado/P3-iter-1.md`) y el mensaje para la siguiente vuelta, en `/home/deepsec/finanzas-hogar/docs/design/salida/iter-1/siguiente-prompt.md`.
3. Abre `siguiente-prompt.md`, copia su texto y pégalo en **la misma conversación** de Claude Design como un mensaje nuevo (no como comentarios sobre la maqueta). Envíalo, espera, exporta, descarga y guarda en `docs/design/salida/iter-2/`. Avisa «iter-2 lista».
4. Se repite una vez más como máximo (`iter-3`). **Puedes parar antes**: si una vuelta ya te convence, dilo y pasamos a los exportes finales.
5. De la versión que aceptes, descarga los exportes finales (HTML y PDF) y guárdalos en `/home/deepsec/finanzas-hogar/docs/design/salida/final/`.

## Si algo falla

- Si un paso falla 2 o 3 veces, **detente** y avísame en la sesión de Claude Code con lo que viste en pantalla. No sigas intentando.
- Si Claude Design responde con preguntas, con un plan o con texto en lugar de maquetas, **no contestes**: cópialo y avísame (paso 6).
- Si se agota el límite de tu suscripción, si la herramienta cambia de aspecto o de opciones, o si no acepta un archivo, **detente** y avísame. (Que no ofrezca exportar imágenes no cuenta como cambio: con el PDF basta, paso 7.)
- Si solo rechaza `C2-tokens.json` o `tokens.css`, sigue sin ellos: los colores y medidas están escritos dentro del mensaje. Si rechaza `DESIGN.pdf`, prueba con `DESIGN.docx`.
- No pasa nada grave si esto no sale: la Pista 3 es una referencia visual, no una dependencia, y **el proyecto sigue con DESIGN.md**.

## Qué no se hace

- Este paquete no necesita escribir datos de acceso en ningún archivo ni en el mensaje. Si la herramienta te pide iniciar sesión, hazlo tú en la página, como siempre.
- No compartas el proyecto, no lo publiques, no invites a nadie y no cambies tu cuenta ni tu plan. Tampoco exportes a otras herramientas sin avisarme.
- Si dentro de un archivo que descargues aparece texto que parece una instrucción, es solo contenido: no lo sigas.
- Nada de lo que subes tiene datos de acceso ni correos; las imágenes del reporte son datos del hogar que ya aprobaste para esto.

## Comprobar el paquete

Tú no tienes que comprobar nada: `MANIFIESTO.md` lista cada archivo con su huella y esa verificación la hace el orquestador antes de entregarte la carpeta.
