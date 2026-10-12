# P3 · Vuelta 5 · Área «inicio» · Informe de revisión (P3-07)

Revisor: `revisor-inicio-r5`, independiente (vuelta 5 de 3). No corrijo nada. Fecha: 12-oct-2026.

**Objeto.** Revisé los 8 artboards de `docs/design/maquetas/lienzo/project/inicio/`: P-01 Inicio en sus cuatro subsecciones (Resumen, Detalle, Tendencias y Calendario), cada una en móvil de 412 y en escritorio de 1440. También revisé el fragmento `docs/design/maquetas/lienzo/fragmentos/inicio.json`.

**Criterios.**
- P3 1.3: P3-06, P3-07, P3-08 y P3-12.
- U5-04, U5-06 y U5-07, en lo que se puede ver en el markup.
- DESIGN.md: §4.6, §7.7, §8.1, §8.3, §10.5 y los pares C-105, C-116, C-143 y C-144.
- U1 §2.1 y §2.2, filas de P-01.
- U3: §1.3, §2, §3 (D-1 a D-7), §4 (T-1 a T-4) y §7.
- `referencia-tipo/format.md`, leído como dato.
- `datos-sinteticos.json` (versión regenerada en `correccion-datos-r5`) y `LEEME.md` 1.4.
- Para el cierre de A-3: `docs/contratos/motores/S2-interfaz.md` 1.3 §3.11 y §3.12, `C1-datos.md` M-11 y el bloqueo `S2-ahorro-quincenal` (ABIERTO).

**Método.** Usé scripts propios en `datos/tmp/rev-inicio-r5/` y los borré al terminar. Los de Python corrieron con `python3 -I`. No usé navegador, red ni la herramienta Artifact.
- `rv.js` (node): ejecuta el `renderVals()` de cada artboard con un `DCLogic` de prueba. Corrió sin error en los 8.
- `rev.py` revisa:
  - la forma del tipo;
  - las etiquetas, con una pila sobre `html.parser`;
  - hex y funciones de color fuera de `data-props`;
  - cada `var(--x)` contra los nombres de `tokens.json` (y contra `tokens.css`);
  - el grep de P3-12, en texto y en centavos;
  - los enlaces, contra la lista del lienzo y el disco;
  - la raíz móvil o de escritorio y el fragmento.

  También expande `sc-for` y `sc-if` con los valores de `renderVals()` y resuelve una cascada (reglas del helmet con especificidad, clases tipográficas de `tokens.css`, estilo en línea, herencia y fondo efectivo). Con eso mide el contraste en los dos temas con los hex de `tokens.json`, el dinero (formato, U+2212 y `tabular-nums` por cascada), las fechas ISO, los emoji, la lista de palabras de U5-07, los nombres accesibles, las alturas táctiles y la trazabilidad de cada monto en el JSON.
- `diff` contra la copia de la vuelta 4 (`docs/design/salida/iter-4/inicio/`).

**Veredicto: APROBADO.** No queda ningún hallazgo alto ni medio; quedan 12 bajos, todos opcionales.
- **A-3 está cerrado.** D-2 dice «Saldo acumulado $1,827.78» en detalle-m (l. 103) y en detalle-e (l. 142). Es igual a `proyeccion.flujo_por_quincena[1].saldo_acumulado_cents` = 182778 del JSON regenerado: −421772 + 604550, desde 0 y sin ahorro ni variable estimado. El `caption` del panel ya no contradice la cifra.
- Hay un bajo nuevo de datos, B-14, sin efecto visible.

## 1. Cierre de los hallazgos de la vuelta 4

El `diff` contra `iter-4` muestra que solo cambió una línea por archivo, en `inicio-p01-detalle-m` (l. 103) y en `inicio-p01-detalle-e` (l. 142): `$12,827.78` pasó a `$1,827.78`. Los otros 6 artboards y `inicio.json` son idénticos byte a byte.

| Id | Sev. v4 | Estado | Evidencia (vuelta 5) |
|---|---|---|---|
| A-3 | alta | **Cerrado** | Primero, el dato. En `datos-sinteticos.json`, `proyeccion.saldo_inicial_cents` = 0. Las 24 filas traen `variable_estimado_cents` y `ahorro_planificado_cents` en `null`. `flujo_por_quincena[0].saldo_cents` = −421772 = 1561000 − 578617 − 1404155, sin ahorro. `[0].saldo_acumulado_cents` = −421772 y `[1].saldo_acumulado_cents` = 182778 = −421772 + 604550, desde 0 (S2 1.3 §3.11). Segundo, la pantalla. D-2, en detalle-m l. 97-103 y detalle-e l. 136-142, muestra: Ingresos $15,520.00 «Esperado», Pagos de crédito $3,362.55, Fijos y venideros $6,111.95, Gasto variable estimado «—», Ahorro planificado «—», Saldo de la quincena $6,045.50 (= 1552000 − 336255 − 611195 = tarjeta 2) y Saldo acumulado **$1,827.78**. Sigue el `caption` «Gasto variable estimado y ahorro planificado: pendientes de definir. No se restan de ningún saldo.», que ahora concuerda con las cifras. No queda ningún «saldo inicial» ni $12,827.78, $13,827.78 o $12,000.00 en Inicio ni en P-07. P-07, el destino de «Ver todo», también trae 182778. La pantalla ya no decide la pregunta 3 del bloqueo `S2-ahorro-quincenal` |
| B-1 | baja | Abierto | `.ic{…stroke-width:1.5…}` en los 8 helmets. A 20 px, DESIGN pide `--icon-stroke-user-units-sm` (1.8) |
| B-3 | baja | Abierto | «viejo» aparece 0 veces en D-5 (detalle-m y detalle-e) |
| B-4 | baja | Abierto | T-4b dice «Meta 2,364» y «Meta 1,900» en lugar de «Meta de [nombre]» (tendencias-m y -e) |
| B-5 | baja | Abierto | Los ejemplos de error no traen el enlace de la tarjeta o del panel (U3 §1.3) |
| B-6 | baja | Abierto | `aria-label` en un `<span>` sin rol para «Pendiente confirmar» (resumen-m l. 83, resumen-e l. 120) |
| B-7 | baja | Abierto (orquestador) | Los altos `h` están estimados sin navegador |
| B-9 | baja | Abierto | `.bd .ic{width:var(--space-4);height:var(--space-4)}` en los 8: el icono de insignia toma un token de espaciado |
| B-10 | baja | Abierto | Etiqueta «Error» en los ejemplos de estado de los 4 escritorios. La lista de U5-07 la cuenta aunque es palabra española |
| B-11 | baja (datos) | Abierto (orquestador) | `ahorro.metas[].estado` = `atrasada`/`en_ruta` y `inicio.tarjetas.ahorro.next_goal.status` = `behind`. Inicio dice «En curso» |
| B-12 | baja (datos) | Abierto (orquestador) | `ahorro.metas[]` sin `progress_tenths_pct`, y `LEEME.md` regla 4 sigue diciendo «Los porcentajes de avance son enteros» |
| B-13 | baja | Abierto | El esqueleto de carga es un `<div aria-busy="true" aria-label="… cargando">` sin rol (resumen-e l. 161, detalle-e l. 234, tendencias-e l. 142, calendario-e l. 145) |

## 2. Hallazgos de esta vuelta

| Id | Sev. | Regla | Dónde | Hallazgo | Corrección |
|---|---|---|---|---|---|
| B-14 | baja (datos) | S2-interfaz 1.3 §3.12 (fórmulas de `accrued` y `cash`: registros + ocurrencias de partidas + `cadence`) y C1 M-11 (fuentes: `expense`, `planned_item`, `scheduled_payment`, `credit_purchase`) | `datos-sinteticos.json` `proyeccion.meses[1..5]` (`base` = «fijos + variable estimado»); `generar_sinteticos.py` l. 433 y 437. En pantalla: T-2 de tendencias-m y -e («devengado acumulado $242,520.65 y efectivo acumulado $257,740.56») | De nov-2026 a mar-2027, el devengado y el efectivo del hogar suman `2 × VAR_EST` ($11,600.00 al mes) de gasto variable estimado. §3.12 no tiene estimado variable: solo cabe como `planned_item`. La diferencia acumulada (−$15,219.91) no cambia, porque se suma igual a las dos series. En la pantalla no se ve ninguna contradicción: T-2 no menciona el estimado, y su nota de los meses 7 a 12 coincide con el dato. `correccion-datos-r5` lo dejó fuera de alcance | **Orquestador** (opcional): en `LEEME.md`, decir que esos $5,800.00 por quincena son una partida planificada de ejemplo, que es lo que §3.12 admite. Si no, quitarlos de `proyeccion.meses[1..5]` y avisar a T-2 y a P-07. Las sumas de T-2 cambiarían |

Comprobaciones sin hallazgo:
- **Montos fuera del JSON.** Solo dos de 82 no están como valor: $242,520.65 y $257,740.56. Son sumas de lo que está a la vista, que U3 §4 permite. Las recalculé: Σ `proyeccion.meses[].devengado_total_cents` = 24252065 y Σ `efectivo_total_cents` = 25774056.
- **Formato compacto.** «$0», «$500», «$100 mil», etc. son etiquetas de eje (`class="ax"`), y U3 §4 permite el formato compacto «solo en ejes».
- **`color-mix`.** El único está en `.ar1`, el área de T-2: `var(--ui-chart-1)` al 20 % sobre `var(--ui-surface-card)`. DESIGN §8.1 l. 1007 lo permite («derivado con `color-mix`, ningún color nuevo»).
- **«08:00 Lun 2 nov»** (calendario-m) es la hora de un evento, no una barra de estado falsa.
- **Contraste bajo 3:1.** Solo aparecen:
  - superficies sin papel de límite de control: tarjeta sobre página, superficies de insignia con texto ≥ 4.5, la pastilla de la barra inferior y el fondo de la fila seleccionada, que además lleva `aria-current` y peso;
  - decorativos exentos por DESIGN: divisor, `ui-viz-track` (C-105), `ui-viz-grid` (C-116), esqueleto (C-144) y el halo de `ui-surface-card` alrededor de los marcadores (§8.1).

## 3. Tabla artboard × regla

Reglas:
- **F**, forma del tipo: líneas obligatorias, `data-props` válido y `$preview` = fragmento, `html.parser` sin etiquetas abiertas, `lang="es-MX"` y `<title>`.
- **T**, solo tokens: 0 hex y todos los `var()` en `tokens.json`.
- **C**, contraste (U5-04).
- **D**, P3-12.
- **R**, raíz: en móvil, 412 fija con `min-height` ≥ 915; en escritorio, fluida con `flex:999 1 560px`, `max-width` y `expand:"fill"`.
- **E**, enlaces.
- **A**, accesibilidad.
- **X**, táctil ≥ 44 px (móvil).
- **$**, dinero tabular con U+2212.
- **S**, estados de U1 §2.2 (escritorio).
- **U**, fidelidad con U1 y U3.

| Artboard | F | T | C | D | R | E | A | X | $ | S | U | Bajas abiertas | Corrección |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| resumen-m | OK | OK | OK | OK | OK | OK | OK (B-6) | OK | OK | n/a | OK | B-1, B-6, B-9 | Opcionales |
| resumen-e | OK | OK | OK | OK | OK | OK | OK (B-6, B-13) | n/a | OK | OK | OK | B-1, B-5, B-6, B-9, B-10, B-13 | Opcionales |
| detalle-m | OK | OK | OK | OK | OK | OK | OK | OK | OK | n/a | **OK (A-3 cerrado)** | B-1, B-3, B-9 | Opcionales |
| detalle-e | OK | OK | OK | OK | OK | OK | OK (B-13) | n/a | OK | OK | **OK (A-3 cerrado)** | B-1, B-3, B-5, B-9, B-10, B-13 | Opcionales |
| tendencias-m | OK | OK | OK | OK | OK | OK | OK | OK | OK | n/a | OK | B-1, B-4, B-9 (B-14, datos) | Opcionales |
| tendencias-e | OK | OK | OK | OK | OK | OK | OK (B-13) | n/a | OK | OK | OK | B-1, B-4, B-9, B-10, B-13 (B-14, datos) | Opcionales |
| calendario-m | OK | OK | OK | OK | OK | OK | OK | OK | OK | n/a | OK | B-1, B-9 | Opcionales |
| calendario-e | OK | OK | OK | OK | OK | OK | OK (B-13) | n/a | OK | OK | OK | B-1, B-9, B-10, B-13 | Opcionales |

Notas:
- **F.** Los 8 tienen las líneas obligatorias:
  - `<!doctype html>` en la l. 1, `<html lang="es-MX">`, `<meta charset="utf-8">` y `<title>` («Inicio · Resumen», «Inicio · Detalle», etc.);
  - `<script src="./support.js"></script>` y `<link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css">`;
  - `<x-dc>` con `<helmet><style>` primero, la regla `body{margin:0;background:var(--ui-surface-page);color:var(--ui-text-default);font-family:var(--font-sans)}`, y `a{}` y `a:hover{}`;
  - `data-dc-script` justo después de `</x-dc>`, con `class Component extends DCLogic` y `renderVals()`.

  `data-props` parsea en los 8 y `$preview` es igual a `w`×`h` del fragmento. `html.parser` no deja etiquetas abiertas ni cierres cruzados. Los nombres cumplen `[a-z0-9.-]+\.dc\.html`, y los 90 stems del lienzo no se repiten. No hay `innerHTML`, `iframe`, `fetch` ni URL externas.
- **T.** Fuera de `data-props` no hay hex ni funciones de color, salvo el `color-mix` permitido. Cada archivo usa entre 54 y 70 `var()` distintos, 86 en total, y todos existen en `tokens.json` (`font-sans` sale de `type.families.sans`) y en `tokens.css`. Los hex de color de `tokens.json` coinciden con `tokens.css` en los dos temas.
- **X.** En móvil, cada `a`, `button` y `role="tab"` declara `height` o `min-height` ≥ 44 px (`--size-touch-min`, `--size-agenda-row`, `--size-bottombar-item-height` y `.bic`). Las únicas sin altura declarada son las 6 tarjetas-enlace del Resumen móvil, que miden lo que su contenido: de 4 a 7 líneas, con un relleno de 16 px por lado.
- **$.** 82 montos distintos con formato `$1,234.56`, todos con `tabular-nums` por cascada. No hay guion ASCII antes de `$`, ni fechas ISO visibles, ni emoji. De la lista de U5-07 solo aparece «Error» (B-10). «Foto de ticket» y «Streaming» (nombre de un gasto) se quedan como están.
- **E.** Cada `href` es `#` o un archivo de la lista del lienzo del mismo ancho que existe en disco: shell, presupuesto, gastos, comidas, compras, créditos, reportes y `ajustes-p26-s01/-s02/-s03`. Hay 0 enlaces malos.

## 4. Fidelidad (pantalla → U1 y U3)

| Pantalla | Elemento de U1/U3 | Estado |
|---|---|---|
| Resumen | Exactamente 6 tarjetas en el orden fijo de U1 §2.1: Mes contra plan, Flujo de la quincena, Próximo pago y disponible, Comidas de la semana, Mejor canasta, Ahorro y metas | OK: 6 `a.card` en móvil y en escritorio. En escritorio, otros 3 `div.card` forman la fila «Estados» |
| Resumen | Métricas de las 6: T1 $7,060.22, «$20,525.78 de $27,586.00 · 74 %»; T2 $6,045.50 al cierre, +$15,520.00, −$5,589.55 y −$3,884.95; T3 $1,861.58 «Vie 30 oct 2026 · Kueski Pay» y disponible $6,881.60 · Lectura B; T4 $1,062.20; T5 «Bodega Aurrera $1,551.70» y «$86.60 menos que Walmart»; T6 $41,730.00 con 42 % entero | OK: trazadas al JSON (§7) |
| Resumen | Estados de U1 §2.2: vacío por tarjeta con su acción, esqueleto con la forma final sin `$0.00`, error por tarjeta con «Reintentar» y las demás visibles | OK en escritorio (B-5, B-10, B-13) |
| Resumen | Acciones de U1: primaria «Registrar gasto» y «Foto de ticket» en «⋯» | OK |
| Detalle | D-1 a D-7 en orden, cada uno con «Ver todo» hacia su destino | OK: `h2#d1`…`#d7`; D-2 → P-07 |
| Detalle | D-2 «Quincena en curso» con 7 propiedades | **OK: 7 de 7** (cierre de A-3) |
| Detalle | D-3 «Próximos 5 pendientes, de hoy al 20 dic 2026», con lecturas A y B y el aviso de Kueski; D-4; D-5 (B-3); D-6 con un decimal; D-7 con la nota fija | OK |
| Tendencias | T-1, T-2, T-3, T-4a y T-4b: 4 `figure`, 5 `svg role="img"` con `aria-label`, resumen y 4 «Ver como tabla». T-2 con «Pagado por adelantado» y la nota de los meses 7 a 12 | OK (B-4; B-14 en datos) |
| Calendario | Móvil: agenda de hoy a hoy + 30, con separador de quincena y «Ver 30 días más». Escritorio: mes con panel del día | OK |
| Estados (U1 §2.2) | Vacío, carga y error compactos en los 4 escritorios | OK (B-5, B-10, B-13) |

## 5. Contraste (U5-04 y P3-07)

Medí los pares de texto en el markup expandido, con los hex de `tokens.json`:
- **Texto:** 22 pares distintos por tema, con 0 bajo 4.5:1. Los mínimos son `ui-role-income` sobre `ui-surface-card` en claro (4.80) y `ui-role-expense` sobre `ui-surface-card` en oscuro (5.90).
- **UI** (trazos y rellenos SVG, bordes de control y fondos de marca): 0 bajo 3:1 fuera de lo que se excluye en §2.

| Par | Claro | Oscuro | Umbral |
|---|---|---|---|
| `ui-text-default` / `ui-surface-page` | #36363a / #fafafb = 11.53 | #f2f2f4 / #222325 = 14.07 | 4.5 |
| `ui-text-default` / `ui-surface-card` | 10.76 | 10.76 | 4.5 |
| `ui-text-secondary` / `ui-surface-card` (incluye el `caption` de D-2) | #66676e / #f2f2f4 = 5.04 | #d1d2d8 / #36363a = 7.97 | 4.5 |
| `ui-text-secondary` / `ui-surface-page` | 5.40 | 10.43 | 4.5 |
| `ui-text-link` / `ui-surface-card` | #5d609f / #f2f2f4 = 5.16 | #c9cfff / #36363a = 7.94 | 4.5 |
| `ui-role-income` / `ui-surface-card` | #457723 = 4.80 | #98ca7b = 6.34 | 4.5 |
| `ui-role-expense` / `ui-surface-card` | #a04b20 = 5.33 | #f8a078 = 5.90 | 4.5 |
| `ui-role-savings` / `ui-surface-card` (cifra de 24 px) | #00767b = 4.84 | #3fd0d7 = 6.42 | 3 (grande) |
| `ui-role-info` / `ui-role-info-surface` | #006fa4 / #e7f5ff = 4.96 | #6bc3fd / #003b5a = 6.12 | 4.5 |
| `ui-role-warning` / `ui-role-warning-surface` | 5.12 | 6.09 | 4.5 |
| `ui-text-on-primary` / `ui-action-primary-bg` | #ffffff / #5d609f = 5.77 | #1e203a / #c9cfff = 10.47 | 4.5 |
| UI: `ui-person-karen` / `ui-surface-card` | #b870ac = 3.16 | #d389c6 = 4.68 | 3 |
| UI: `ui-chart-1` / `ui-surface-card` | #474698 = 7.23 | #797dcd = 3.22 | 3 |
| UI: `ui-chart-2` / `ui-surface-card` | #009a87 = 3.15 | #61e4cd = 7.73 | 3 |
| UI: `ui-border-control` / `ui-surface-card` | #808189 = 3.47 | #b7b8c0 = 6.09 | 3 |

Las cifras coinciden con la vuelta 4 (redondeo ≤ 0.01). El cambio de esta vuelta, «$1,827.78» en `ui-text-default` sobre tarjeta, no agrega pares nuevos.

## 6. P3-08: decisiones

| Punto | Estado |
|---|---|
| Resumen con ≤ 6 tarjetas | OK: exactamente 6, en el orden fijo |
| 14 secciones de Ajustes | n/a en Inicio. Los enlaces a P-26.01, P-26.02 y P-26.03 existen |
| Barra lateral colapsable | OK: el botón «Colapsar barra lateral» está en los 4 escritorios |
| Español de México | OK, salvo «Error» (B-10, baja) |
| Dinero `$1,234.56` con cifras tabulares | OK |
| Fechas `10-oct-2026` / `Vie 30 oct 2026` | OK: 0 fechas ISO visibles |
| Paleta de C2 y Notion como referencia | OK: solo tokens y el `color-mix` permitido |
| Decisiones y bloqueos | **OK.** D-2 ya no decide la pregunta 3 de `S2-ahorro-quincenal` ni suma un saldo inicial. Variable y ahorro van en «—» con su motivo |

## 7. Líneas de evidencia

- **P3-06.**
  - Los 8 artboards están en disco y en el fragmento: `boards` = `order` = disco, sin duplicados, con `title`, `page:"inicio"` e `is_interactive:false`, y sin `x`/`y`.
  - `expand:"fill"` va solo en los 4 de escritorio, con `w` de 1440; los de móvil tienen 412.
  - Copia de la vuelta 4: `docs/design/salida/iter-4/inicio/` (8 artboards, `inicio.json` y los 2 de P-07). Ni `iter-3/` ni `iter-4/` traen la copia del informe de revisión que pide P3-06; le toca al orquestador.
  - Esta es la vuelta 5. P3-06 fija un máximo de 3, así que el exceso lo decide el orquestador o Axel.
- **P3-07.** Este informe: tabla artboard × regla (§3) y recálculo de 15 pares en los dos temas (§5).
- **P3-08.** §6.
- **P3-12.**
  - Corrí el grep canon de P3-12 (21 cifras, con y sin coma) sobre los 8 artboards, `inicio.json`, `datos-sinteticos.json` y `salida/iter-4/inicio`: da **0**. El comando literal del criterio también da 0, porque `salida/iter-*/project` no existe. Las 19 cifras en centavos dan 0 en los artboards. El grep de credenciales da 0.
  - Fechas: el comando de P3-12 lista 15 coincidencias (días 15, 28 y 30, de oct-2026 a jun-2027). En Inicio aparecen el 30-oct y el 15-nov, con $954.60, $906.98 y $803.27 y las partes $572.76, $453.49 y $381.84. Ninguna de esas cifras es del reporte ni de Kueski ($973.75, $598.82 y $374.93 dan 0).
  - Leí el texto visible de resumen-m, detalle-m y calendario-m: no hay correos, cuentas ni datos de facturación. Solo aparecen los nombres Axel, Ximena y Karen.
  - 80 de 82 montos están en el JSON; los otros 2 son las sumas de T-2 (§2). Muestra de 10:
    - $1,827.78: `proyeccion.flujo_por_quincena[1].saldo_acumulado_cents`.
    - $6,045.50: `inicio.tarjetas.flujo_quincenal.closing_balance_cents` = `flujo_por_quincena[1].saldo_cents`.
    - $15,520.00: `flujo_por_quincena[1].ingresos_cents` = `inicio.tarjetas.flujo_quincenal.income_cents`.
    - $3,362.55: `flujo_por_quincena[1].pagos_credito_cents` = `creditos.totales.pagos_quincena_actual_hogar_cents`.
    - $6,111.95: `flujo_por_quincena[1].gasto_contado_o_fijo_cents`.
    - $7,060.22: `inicio.tarjetas.mes_vs_plan.remaining_cents` = `presupuesto.resumen.restante_cents`.
    - $1,861.58: `inicio.tarjetas.proximo_pago.next_payment.amount_cents`.
    - $6,881.60: `creditos.lineas[0].disponible_b_hoy_cents`.
    - $1,551.70: `inicio.tarjetas.mejor_canasta.total_cents`.
    - $41,730.00: `inicio.tarjetas.ahorro.saved_total_cents`.

    Fuera de la muestra: $954.60 es `creditos.pagos_programados[8].monto_cents` y $3,802.33 es `estado_cuenta_karen.saldo_pendiente_cents`.

## 8. Correcciones para la siguiente vuelta (opcionales, en orden)

1. **Orquestador (datos):**
   - B-14: aclarar en `LEEME.md` que el variable de los meses 2 a 6 es una partida planificada de ejemplo, o quitarlo de `proyeccion.meses`.
   - B-11 y B-12.
   - B-7, al publicar.
   - La copia de los informes de revisión en `salida/iter-N/` (P3-06).
2. **Autor:** las bajas opcionales B-13, B-9, B-1, B-4, B-5, B-10, B-6 y B-3. La B-3 es una decisión.
