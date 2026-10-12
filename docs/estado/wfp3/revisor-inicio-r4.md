# revisor-inicio-r4 (vuelta 4, Pista 3, área Inicio)

Fecha: 12-oct-2026. Revisor independiente; no corregí artboards ni datos.

## Hecho
- Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/inicio/` y `fragmentos/inicio.json` contra P3 1.3 (P3-06, P3-07, P3-08 y P3-12), U5-04/06/07, DESIGN.md, U1 §2.1/§2.2, U3 y `referencia-tipo/format.md`. Para el cierre de A-3 usé también `S2-interfaz.md` 1.3 §3.11 y el bloqueo `S2-ahorro-quincenal`.
- Informe P3-07: `docs/estado/wfp3/P3-iter-4-inicio.md`.
- Veredicto: **NO APROBADO**. Queda 1 alto, con 0 medios y 11 bajos.
  - **A-3 (alta), reabierto en parte.** Ya está bien: «—» en gasto variable estimado y en ahorro planificado, el saldo de la quincena en $6,045.50 (igual a la tarjeta 2) y el `caption`. Falta «Saldo acumulado $12,827.78» (detalle-m l. 103, detalle-e l. 142). Esa cifra suma un saldo inicial de $12,000.00 que S2 1.3 §3.11 excluye («parte de 0»). También resta el ahorro planificado de la 1.ª quincena de octubre (−$1,000.00), con lo que decide la pregunta 3 de un bloqueo abierto y contradice el `caption` del mismo panel. Con el contrato y los mismos datos, la cifra es **$1,827.78** (−421772 + 604550). La sugerencia de la vuelta 3 ($13,827.78) también sumaba el saldo inicial; la corrijo.
  - B-13 (nuevo, baja): `aria-label` sin rol en los esqueletos de carga de los 4 escritorios.
  - Siguen abiertos B-1, B-3, B-4, B-5, B-6, B-7, B-9, B-10, B-11 y B-12. Nadie los tocó en esta vuelta.

## Verificación
- Scripts propios en `datos/tmp/rev-inicio-r4/`, borrados al final: `rv.js` (node) y `rev.py` y `traza.py` (python3 -I).
- `diff` contra `salida/iter-3/inicio/`: solo cambiaron las líneas de D-2, el `caption` nuevo y los altos de detalle. Los otros 6 artboards son idénticos.
- Forma, 8/8 OK: líneas obligatorias exactas, `data-props` válido con `$preview` = fragmento, `html.parser` sin etiquetas abiertas, `lang="es-MX"`, `<title>`, nombres válidos y stems únicos en todo el lienzo. `renderVals()` corre sin error en los 8.
- Tokens: 0 hex y 0 funciones de color fuera de `data-props`. Los 86 `var()` existen en `tokens.json` y en `tokens.css`.
- Contraste, con el markup expandido y la cascada resuelta en claro y oscuro:
  - Texto: 18 pares por tema, 0 bajo 4.5. Los mínimos son 4.80 en claro y 5.90 en oscuro.
  - UI: 53 pares por tema, 0 bajo 3, salvo los decorativos que DESIGN exceptúa (C-105, C-116 y el halo de §8.1).
- P3-12: el grep canon da 0 (con y sin coma, y en centavos) y el de credenciales también da 0. De los 82 montos, 80 están en el JSON; los otros 2 son las sumas de T-2, que U3 §4 permite. Fechas coincidentes revisadas: ninguna va con un monto real.
- Enlaces: todos apuntan a archivos de la lista del lienzo con el mismo ancho, o a `#`.
- Raíces: las 4 móviles son fijas de 412 px con `min-height` ≥ 915. Los 4 escritorios son fluidos y llevan `expand:"fill"`.
- Táctil ≥ 44 px en móvil. Dinero tabular con U+2212. Sin emoji, sin fechas ISO visibles y sin barra de estado falsa.

## Supuestos
- El móvil y el escritorio de una misma pantalla comparten `<title>`, porque el encargo pide «nombre de la pantalla». No lo cuento como falla.
- `../ds/finanzas-hogar/tokens.css` aún no existe bajo `lienzo/project/`. Lo publica el orquestador junto con `fonts/` (`ds-tokens-css.md`), y no lo cuento como falla de forma.
- A-3 sigue en alta por la misma regla que en la vuelta 3: contradice S2 1.3 §3.11 y U3 D-2 en la pantalla principal y decide una pregunta de un bloqueo abierto. Lo sostengo aunque `correccion-datos-r4` haya dejado la 1.ª quincena fuera de alcance y el orquestador haya fijado $12,827.78, porque ninguna de las dos cosas cambia el contrato.
- «Error» es palabra española. B-10 sigue en baja porque la lista automática de U5-07 la cuenta.

## Bloqueos
- Ninguno para la revisión. Para cerrar A-3, primero el orquestador tiene que corregir `datos-sinteticos.json`; después el autor cambia 2 líneas.

## Pendientes
- **Orquestador**:
  - A-3 en `proyeccion.flujo_por_quincena`: ningún saldo con ahorro ni variable estimado, y el acumulado desde 0 en 2026-10-q1, sin `saldo_inicial_cents`. Avisar a presupuesto P-07, que pinta la serie, el «saldo inicial» y los KPI.
  - B-11 y B-12.
  - B-7, al publicar.
  - Copia de esta vuelta y de este informe en `docs/design/salida/iter-4/` (P3-06); esa carpeta todavía no existe.
- **Autor**:
  - A-3: «Saldo acumulado $1,827.78» (o el valor corregido) en detalle-m l. 103 y detalle-e l. 142.
  - Bajas opcionales.
- Fuera de mi área, como nota: las quincenas futuras (`[2]` a `[23]`) también restan $5,800.00 de variable estimado y $1,000.00 de ahorro en su saldo, contra S2 §3.11. Le toca a P-07.

## Consumo aproximado
~120 KB leídos, ~30 KB escritos (informe, reporte y scripts temporales), ~38 acciones de herramienta. cerrado_por_tope = false.
