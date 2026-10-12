# revisor-inicio-r5 (vuelta 5, Pista 3, área Inicio)

Fecha: 12-oct-2026. Soy revisor independiente: no corregí artboards ni datos.

## Hecho
- Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/inicio/` y `fragmentos/inicio.json`. Los criterios fueron P3 1.3 (P3-06, P3-07, P3-08 y P3-12), U5-04/06/07, DESIGN.md, U1 §2.1/§2.2, U3 y `referencia-tipo/format.md`. Para el cierre de A-3 usé también `S2-interfaz.md` 1.3 §3.11 y §3.12, C1 M-11 y el bloqueo `S2-ahorro-quincenal`.
- Escribí el informe P3-07 en `docs/estado/wfp3/P3-iter-5-inicio.md`.
- **Veredicto: APROBADO.** Quedan 0 hallazgos altos, 0 medios y 12 bajos, todos opcionales.
  - **A-3 (alta) está cerrado.** D-2 dice «Saldo acumulado $1,827.78» en detalle-m l. 103 y detalle-e l. 142. Es igual a `flujo_por_quincena[1].saldo_acumulado_cents` = 182778 del JSON regenerado (−421772 + 604550, desde 0, sin ahorro ni variable). Las otras 6 propiedades de D-2 y el `caption` siguen correctos y ya concuerdan con la cifra. En Inicio no queda ningún «saldo inicial», y P-07 trae el mismo 182778.
  - **B-14 (nuevo, baja, datos, sin efecto visible).** `proyeccion.meses[1..5]` suma $11,600.00 al mes de gasto variable estimado al devengado y al efectivo; su `base` dice «fijos + variable estimado». S2 §3.12 y C1 M-11 solo lo admiten como `planned_item`. La diferencia acumulada no cambia. Las sumas de T-2 ($242,520.65 y $257,740.56) lo incluyen.
  - Siguen abiertos los bajos B-1, B-3, B-4, B-5, B-6, B-7, B-9, B-10, B-11, B-12 y B-13. Nadie los tocó.

## Verificación
- Scripts propios en `datos/tmp/rev-inicio-r5/`, borrados al final:
  - `rv.js` (node) ejecuta `renderVals()`: corrió en 8/8 sin error.
  - `rev.py` (python3 -I) revisa la forma, los tokens, P3-12, los enlaces, la raíz y el fragmento. Expande `sc-for` y `sc-if` y resuelve la cascada para medir contraste, dinero, idioma, accesibilidad, alturas táctiles y la trazabilidad de los montos.
- `diff` contra `salida/iter-4/inicio/`: solo cambió una línea en detalle-m (l. 103) y otra en detalle-e (l. 142). Lo demás, `inicio.json` incluido, es idéntico.
- **Forma (8/8 OK):**
  - Las líneas obligatorias están exactas.
  - `data-props` es válido y `$preview` coincide con el fragmento.
  - `html.parser` no deja etiquetas abiertas.
  - Todos llevan `lang="es-MX"` y `<title>`.
  - Los nombres son válidos y los 90 stems del lienzo no se repiten.
  - No hay `innerHTML`, `iframe` ni red.
- **Tokens:** 0 hex fuera de `data-props`. El único `color-mix` es el área de T-2, que DESIGN §8.1 permite. Los 86 `var()` existen en `tokens.json` y en `tokens.css`.
- **Contraste**, en claro y oscuro con los hex de `tokens.json`:
  - Texto: 22 pares, 0 bajo 4.5. Los mínimos son 4.80 en claro y 5.90 en oscuro.
  - UI: 0 bajo 3, salvo los decorativos que DESIGN exceptúa (C-105, C-116, C-144 y el halo de §8.1) y las superficies que no son límite de control.
- **P3-12:**
  - El grep canon da 0, con y sin coma y en centavos, sobre los artboards, el fragmento, el JSON y `iter-4`. El de credenciales también da 0.
  - De los 82 montos, 80 están en el JSON. Los otros 2 son las sumas de T-2, que recalculé contra `proyeccion.meses`.
  - Revisé las 15 fechas que coinciden: ninguna va con un monto real.
- **Enlaces:** todos son `#` o archivos de la lista del lienzo del mismo ancho que existen en disco.
- **Raíces:** las 4 móviles son fijas de 412 px con `min-height` igual al alto del fragmento (≥ 915). Los 4 escritorios son fluidos, con `flex-wrap`, `flex:999 1 560px`, `max-width` y `expand:"fill"`.
- **Táctil:** en móvil todo control mide ≥ 44 px; las 6 tarjetas-enlace miden lo que su contenido.
- **Dinero:** formato `$1,234.56` tabular con U+2212.
- **Idioma y formato:** sin emoji, sin fechas ISO visibles y sin barra de estado falsa. De la lista de U5-07 solo aparece «Error» (B-10).

## Supuestos
- El móvil y el escritorio de una misma pantalla comparten `<title>`, porque el encargo pide el «nombre de la pantalla». No lo cuento como falla.
- `../ds/finanzas-hogar/tokens.css` y `support.js` todavía no existen bajo `lienzo/project/`. Los publica el orquestador (`ds-tokens-css.md`), así que no los cuento como falla de forma.
- Las superficies tintadas (insignias con texto ≥ 4.5, tarjeta sobre página, pastilla y fila seleccionada con `aria-current`) no son límites de control según WCAG 1.4.11. No les exijo 3:1.
- B-14 va en baja y no en media: §3.12 admite el monto como `planned_item`, la pantalla no lo nombra y no hay ninguna contradicción visible. El arreglo es una aclaración o un ajuste de datos del orquestador.
- «Ticket» («Foto de ticket») es uso de México y está en U1. «Streaming» es el nombre de un gasto sintético. Ninguno de los dos cuenta como inglés.

## Bloqueos
- Ninguno para la revisión.

## Pendientes
- **Orquestador:**
  - B-14 (aclarar en `LEEME.md` o ajustar `proyeccion.meses` y avisar a T-2 y P-07).
  - B-11 y B-12.
  - B-7, al publicar.
  - La copia de los informes de revisión en `docs/design/salida/iter-N/` (P3-06): `iter-3/` e `iter-4/` no la traen.
  - La vuelta 5 rebasa el máximo de 3 de P3-06. Lo decide el orquestador o Axel.
- **Autor:** las bajas opcionales B-13, B-9, B-1, B-4, B-5, B-10, B-6 y B-3.

## Consumo aproximado
~150 KB leídos, ~25 KB escritos (informe, reporte y scripts temporales), ~40 acciones de herramienta. cerrado_por_tope = false.
