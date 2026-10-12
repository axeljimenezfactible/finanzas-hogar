# revisor-inicio-r1 · Pista 3, fase 2 (lienzo), vuelta 1 de 3

Fecha: 12-oct-2026 · Agente: `revisor-inicio-r1` (revisor independiente; no corrige) · Área: `inicio`

## Hecho

- Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/inicio/`: P-01 Resumen, Detalle, Tendencias y Calendario, en móvil y escritorio. También revisé el fragmento `docs/design/maquetas/lienzo/fragmentos/inicio.json`. Los comparé con P3 1.3 (P3-04, P3-06, P3-07, P3-08 y P3-12), con U5 (U5-04, U5-06 y U5-07), con DESIGN.md, U1, U3 y con `referencia-tipo/format.md`.
- Escribí el informe P3-07 en `docs/estado/wfp3/P3-iter-1-inicio.md`. Trae la tabla artboard × regla, la fidelidad a U1 y U3, el contraste, la lista de P3-08 y las líneas de evidencia.
- **Veredicto: no aprobado.** Hay 1 alto, 0 medios y 8 bajos.
  - **A-1 (alto).** La insignia de estado de las metas dice «Atrasada» en la tarjeta 6 y en D-6, y «En ruta» en D-6. Eso contradice la decisión Q2-18 A, que fija tres estados (En curso, Vencida y Cumplida) y deja fuera «atrasada», como dice U3. Las dos metas deberían decir «En curso». El origen está en `datos-sinteticos.json` (`status: behind` y `etiqueta_estado: Atrasada`; `ahorro.metas[].estado`). El autor copió los datos tal como venían, así que la corrección empieza por el orquestador, en el archivo de datos.
  - **Bajos.** Son estos:
    - Trazo de icono de 1.5 a 20 px (DESIGN pide 1.8).
    - «Pendiente» sin icono en D-3.
    - La marca «viejo» de D-5 y el supuesto 7 del autor.
    - «Meta de [nombre]» en T-4b.
    - Enlace de la tarjeta ausente en los ejemplos de error.
    - Nombre accesible de «Pendiente confirmar».
    - Altos estimados.
    - Tope de escritura excedido por el autor.

## Verificación

Usé scripts propios con `python3 -I` y `node` en `datos/tmp/rev-inicio/`: el verificador de forma y tokens, un motor de cascada CSS para el contraste, la revisión de accesibilidad, la traza de montos y el evaluador de `renderVals()`. Ya los borré.

- **Forma: 8 de 8 bien.**
  - Las líneas obligatorias están completas y en orden.
  - `data-props` es JSON válido y `$preview` es igual a `w`×`h` del fragmento.
  - `html.parser` no encuentra etiquetas sin cerrar.
  - `lang="es-MX"` y cada archivo tiene su `<title>`.
  - Los nombres son válidos y los stems únicos en las 90 maquetas del lienzo.
  - `renderVals()` se ejecuta y todos los huecos tienen valor.
  - El fragmento lista exactamente los 8 archivos, con `expand: "fill"` solo en escritorio, `is_interactive: false` y sin x ni y.
- **Tokens.**
  - 0 hex y 0 `rgb()`/`hsl()` fuera de `data-props`.
  - Hay un solo `color-mix`, el del área de T-2, y DESIGN §8.1 lo permite.
  - Todos los `var(--x)` existen en `tokens.css`. En `tokens.json` existen todos menos `font-sans`, que es `type.families.sans`.
- **Contraste.**
  - Evalué 62 pares (tema × primer plano × fondo) y ninguno queda bajo el umbral.
  - Los 62 coinciden con `contraste.csv` o `contraste-componentes.csv` al ±0.01.
  - Recalculé a mano 9 pares extra (error, medidor, tiendas, Karen y selección), y todos pasan.
- **P3-12.**
  - El grep de las 21 cifras canon da 0, igual que los grep de centavos y de credenciales.
  - Las fechas que coinciden con el modelo son las 15 quincenales de `LEEME.md`, y en Inicio todas van con montos sintéticos.
  - Todos los montos salen de `datos-sinteticos.json`, salvo 2 acumulados de T-2. Esos 2 son sumas de `proyeccion.meses`, cuadran con `diferencia_acumulada` y el autor los declara.
  - Tomé una muestra de 10 trazas y está en el informe.
- **Accesibilidad.**
  - No hay controles falsos, controles sin nombre ni controles anidados.
  - Los 207 SVG llevan `aria-hidden` o `role="img"` con nombre.
  - Medí los objetivos táctiles en móvil resolviendo los tokens y ninguno queda por debajo de 44 px.
  - No hay emoji ni barra de estado falsa.
  - Todo el dinero lleva `tabular-nums` y los negativos usan U+2212.
- **Enlaces.**
  - Revisé los del markup y los de `renderVals()` (barra lateral, agenda y celdas).
  - Todos existen y son del mismo ancho, o son «#».

## Supuestos

1. Califiqué A-1 como alto porque la escala del encargo pone «decisión contradicha» en alto, aunque el autor no lo inventó: el archivo de datos usa el enum de C1 v1.3 (`on_track`/`behind`), que U3 §2 marca como pendiente de corregir en C1. Si el orquestador decide que los datos mandan sobre U3 en este punto, A-1 baja de severidad, pero entonces habría que registrar la excepción a Q2-18 A.
2. Varias palabras no cuentan como inglés. «Error», en la fila «Estados», es palabra del español y es el nombre del estado en U1 §2.2 y U3 §7, aunque esté en la lista de U5-07. Lo mismo pasa con «Total». «Kueski Pay» es el nombre de la línea en los datos y se trata como nombre propio.
3. `font-sans` cumple «existe en tokens.json» porque es `type.families.sans` con su nombre de CSS (nota `meta.note` de `tokens.json`).
4. «Exactamente 6 tarjetas» se cumple: en Resumen hay 6 enlaces-tarjeta en cada ancho. Las tarjetas de la fila «Estados» son ejemplos de estado de la tarjeta 1, no tarjetas adicionales.
5. `is_interactive: false` sigue la plantilla del encargo, aunque los artboards tienen enlaces de prototipo. Si el orquestador quiere recorrerlos en Play, tiene que ponerlo en `true` (format.md).
6. El informe va en `docs/estado/wfp3/`, como pidió el encargo, y no en `docs/estado/`, como dice el texto de P3-07.

## Bloqueos

Ninguno.

## Pendientes

- **Orquestador:**
  - Corregir los estados de meta en `datos-sinteticos.json` (A-1; afecta también a P-09) y volver a correr `verificar_sinteticos.py`.
  - Medir los altos al publicar (B-7).
  - Copiar este informe a `docs/design/salida/iter-1/` (P3-06).
- **Vuelta 2:** el autor corrige A-1 y los bajos que decida, y otro revisor vuelve a pasar los mismos scripts.

## Consumo aproximado

- **Leído:** unos 270 KB. Incluye P3 entero; U5 (§U5-04 a U5-07); format.md entero; U3 §1 a §10; U1 (§0.4, §0.5, §1.2 a §1.4, §2.2, §2.4, §3.3); DESIGN (§2.10, §5.2, §8). También el reporte del autor y un informe hermano como referencia de formato, el cuerpo de los artboards sin trazos SVG, partes de `datos-sinteticos.json` y de `LEEME.md`, y `tokens.css`.
- **Escrito:** unos 45 KB. Son los scripts temporales, ya borrados, este reporte y `P3-iter-1-inicio.md`.
- **Acciones de herramienta:** unas 50.
