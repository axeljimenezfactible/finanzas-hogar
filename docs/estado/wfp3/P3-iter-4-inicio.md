# P3 · Vuelta 4 · Área «inicio» · Informe de revisión (P3-07)

Revisor: `revisor-inicio-r4` (independiente, vuelta 4 de 3; no corrige). Fecha: 12-oct-2026.

**Objeto.** Los 8 artboards de `docs/design/maquetas/lienzo/project/inicio/` (P-01 Inicio: Resumen, Detalle, Tendencias y Calendario, cada una en móvil 412 y escritorio 1440) y el fragmento `docs/design/maquetas/lienzo/fragmentos/inicio.json`.

**Criterios.** P3 1.3 (P3-06, P3-07, P3-08, P3-12); U5-04, U5-06 y U5-07, en lo que se ve en el markup; DESIGN.md (§8.1, §10.5, tabla de pares C-105/C-116); U1 §2.1 y §2.2 (filas de P-01); U3 (§2 tarjetas, §3 fila D-2, §4 etiquetas de dinero, §7 estados); `referencia-tipo/format.md` (leído como dato); `datos-sinteticos.json` y `LEEME.md` 1.3; contrato `docs/contratos/motores/S2-interfaz.md` 1.3 §3.11 y bloqueo `docs/estado/bloqueos/S2-ahorro-quincenal.md` (ABIERTO) para el cierre de A-3.

**Método.** Scripts propios en `datos/tmp/rev-inicio-r4/`, borrados al terminar; los de Python corrieron con `python3 -I`:
- `rv.js` (node): ejecuta el `renderVals()` de cada artboard con un `DCLogic` de prueba. Resultado: 8/8 sin error.
- `rev.py`: revisa la forma del tipo, cierra etiquetas con una pila sobre `html.parser`, busca hex y funciones de color fuera de `data-props`, compara cada `var(--x)` con `tokens.json`, corre el grep de P3-12, valida los enlaces contra la lista del lienzo y el disco, la raíz móvil o de escritorio y el fragmento. También expande `sc-for` y `sc-if` con los valores de `renderVals()` y resuelve una cascada simple (reglas del helmet con especificidad, estilo en línea, herencia y fondo efectivo). Con eso mide el contraste en los dos temas, el dinero, la lista de palabras de U5-07, los emoji, los nombres accesibles y las alturas táctiles.
- `traza.py`: busca cada monto de los artboards en `datos-sinteticos.json` y recalcula el saldo de D-2 con la fórmula de S2 1.3 §3.11.
- `diff` contra la copia de la vuelta 3 (`docs/design/salida/iter-3/inicio/`).

No usé navegador, red ni la herramienta Artifact.

**Veredicto: NO APROBADO.** Queda 1 hallazgo alto (A-3, reabierto en parte), con 0 medios y 11 bajos.
- A-3 está corregido en 4 de sus 5 puntos: «Gasto variable estimado —», «Ahorro planificado —», «Saldo de la quincena $6,045.50», que ya es igual a la tarjeta 2, y el `caption` con el motivo.
- Falta el quinto: **«Saldo acumulado $12,827.78» sigue fuera del contrato.** Arrastra el ahorro planificado de la 1.ª quincena de octubre (−$1,000.00) y un saldo inicial de $12,000.00 que el modelo no tiene. Así contradice el `caption` del mismo panel («No se restan de ningún saldo»).
- Todo lo demás pasa: forma, tokens, contraste, P3-12, enlaces, raíces y fragmento.

## 1. Cierre de los hallazgos de la vuelta 3

`diff` contra `iter-3`: solo cambiaron `inicio-p01-detalle-m` (l. 62, 100-103, la 105 nueva y la 196) y `inicio-p01-detalle-e` (l. 139-142, la 144 nueva y la 252), además de los altos de detalle en `inicio.json` (4250 → 4298 y 3000 → 3032). Los otros 6 artboards son idénticos byte a byte.

| Id | Sev. v3 | Estado | Evidencia (vuelta 4) |
|---|---|---|---|
| A-3 | alta | **Reabierto en parte** (sigue alta) | ✔ detalle-m l. 100-101 y detalle-e l. 139-140: «Gasto variable estimado —» y «Ahorro planificado —» (U3 D-2: «—», nunca $0.00). ✔ l. 102 y 141: «Saldo de la quincena $6,045.50» = tarjeta 2 «$6,045.50 al cierre» = `inicio.tarjetas.flujo_quincenal.closing_balance_cents` 604550 = $15,520.00 − $3,362.55 − $6,111.95. ✔ l. 105 y 144: `caption` «pendientes de definir. No se restan de ningún saldo.» ✘ l. 103 y 142: «Saldo acumulado $12,827.78» (§2). La cifra que sugirió la vuelta 3 ($13,827.78) también sumaba el saldo inicial, así que la corrijo en §2 |
| B-1 | baja | Abierto | `.ic{…stroke-width:1.5…}` en los 8 helmets; a 20 px DESIGN pide `--icon-stroke-user-units-sm` |
| B-3 | baja | Abierto | «viejo» aparece 0 veces en D-5 (detalle-m y detalle-e) |
| B-4 | baja | Abierto | T-4b: «Meta 2,364», «Meta 1,900» en lugar de «Meta de [nombre]» (tendencias-m l. 97-98, tendencias-e l. 133-134) |
| B-5 | baja | Abierto | Los ejemplos de error no traen el enlace de la tarjeta o del panel (U3 §1.3) |
| B-6 | baja | Abierto | `aria-label` en un `<span>` sin rol para «Pendiente confirmar» (resumen-m l. 83, resumen-e l. 120) |
| B-7 | baja | Abierto (orquestador) | Altos `h` estimados sin navegador; los de detalle crecieron 48 y 32 px por el renglón nuevo |
| B-9 | baja | Abierto | `.bd .ic{width:var(--space-4);height:var(--space-4)}`: el icono de insignia toma un token de espaciado |
| B-10 | baja | Abierto | Etiqueta «Error» en los ejemplos de estado de los 4 escritorios; la lista de U5-07 la cuenta aunque es palabra española |
| B-11 | baja (datos) | Abierto (orquestador) | `ahorro.metas[].estado` sigue en `atrasada`/`en_ruta` y `inicio.tarjetas.ahorro.next_goal.status` en `behind` («Atrasada»). Inicio dice «En curso» |
| B-12 | baja (datos) | Abierto (orquestador) | `ahorro.metas[]` sin décimas (`progress_tenths_pct` ausente) y `LEEME.md` regla 4 sigue diciendo «Los porcentajes de avance son enteros» |

## 2. Hallazgos de esta vuelta

| Id | Sev. | Regla | Artboards (líneas) | Hallazgo | Corrección |
|---|---|---|---|---|---|
| **A-3** (reabierto, parcial) | **alta** | `S2-interfaz.md` 1.3 §3.11 (`savings_planned_cents` y `variable_estimated_cents`: «no entra a ningún saldo»; `cumulative_balance_cents`: «Σ `balance_cents` desde la primera quincena de la consulta (`from`) hasta esta, inclusive; **parte de 0** (el modelo no tiene saldo inicial de cuentas; M-10 empieza en la Q1 del mes en curso)»); U3 §3 fila D-2 («no lo resta de `balance_cents` ni de `cumulative_balance_cents`»); bloqueo `S2-ahorro-quincenal` ABIERTO (pregunta 3); P3-12 (cifras que cuadran entre sí) | detalle-m l. 103 (y `caption` l. 105); detalle-e l. 142 (y `caption` l. 144) | D-2 muestra «Saldo acumulado $12,827.78» (`proyeccion.flujo_por_quincena[1].saldo_acumulado_cents` = 1282778). Esa cifra es `saldo_inicial_cents` 1200000 + `flujo_por_quincena[0].saldo_cents` −521772 + 604550. El −521772 de la 1.ª quincena (1 al 15 oct) resta `ahorro_planificado_cents` 100000; con §3.11 ese saldo es 1561000 − 578617 − 1404155 = −421772. Además, §3.11 dice que el acumulado parte de 0, sin saldo inicial. Con el contrato y los mismos datos, D-2 debe decir **$1,827.78** (−421772 + 604550 = 182778). Hoy la pantalla decide, para la 1.ª quincena, la pregunta 3 de un bloqueo abierto («¿resta del saldo?»). También inventa un saldo de cuentas que el modelo no tiene, y el `caption` del mismo panel afirma lo contrario. `correccion-datos-r4` lo dejó así a propósito («fuera del alcance»), y el orquestador fijó $12,827.78. Ninguna de las dos cosas cambia el contrato | **Orquestador** (primero): en `datos-sinteticos.json` `proyeccion.flujo_por_quincena`, que ninguna fila reste el ahorro planificado ni el gasto variable estimado (S2 §3.11: los dos en `null`), y `saldo_acumulado_cents` desde 0 en 2026-10-q1, sin `saldo_inicial_cents`. Para las dos primeras filas: −421772 y 182778. Hay que avisar a presupuesto (P-07), que pinta la misma serie, el «saldo inicial» y los KPI. Si Axel quiere un saldo inicial, se abre un cambio de S2 y C1 por bloqueo; el artboard no lo decide. **Autor**: en detalle-m l. 103 y detalle-e l. 142, «Saldo acumulado $1,827.78» (o el valor corregido del JSON), sin tocar el `caption` |
| B-13 | baja | U5-06/DESIGN §9 (nombres accesibles), misma clase que B-6 | resumen-e l. 161, detalle-e l. 234, tendencias-e l. 142, calendario-e l. 145 | El esqueleto de los ejemplos de carga es un `<div aria-busy="true" aria-label="… cargando">` sin rol. Un `aria-label` en un elemento genérico no se anuncia | `role="status"` en el contenedor, o el texto «Cargando…» visible para lectores en lugar del `aria-label` |

Comprobaciones sin hallazgo:
- Las sumas de T-2 ($242,520.65 y $257,740.56) siguen siendo las 2 únicas cifras de 82 que no están en el JSON como valor. Son sumas de lo que está a la vista (U3 §4), igual que en la vuelta 3.
- «$0» y «$500» (tendencias-m y -e) son etiquetas de eje `class="ax te"`, como «$10 mil». U3 §4 permite el formato compacto «solo en ejes».
- Por debajo de 3:1 solo aparecen trazos decorativos: el halo de los marcadores en `ui-surface-card` (DESIGN §8.1), la pista `ui-viz-track` (C-105, «decorativo») y la retícula `ui-viz-grid` (C-116, «decorativo»).
- «08:00 Lun 2 nov» (calendario-m) es la hora de un evento («Reporte mensual»), no una barra de estado falsa.

## 3. Tabla artboard × regla

Reglas:
- F, forma del tipo: líneas obligatorias, `data-props` válido y `$preview` = fragmento, `html.parser` sin etiquetas abiertas, `lang="es-MX"`, `<title>`.
- T, solo tokens: 0 hex y todos los `var()` en `tokens.json` y en `sistema/project/tokens.css`.
- C, contraste (U5-04).
- D, P3-12.
- R, raíz: móvil de 412 fija con `min-height` ≥ 915; escritorio fluido con `flex:999 1 560px`, `max-width` y `expand:"fill"`.
- E, enlaces.
- A, accesibilidad.
- X, táctil ≥ 44 px (móvil).
- $, dinero tabular con U+2212.
- S, estados de U1 §2.2 (escritorio).
- U, fidelidad con U1 y U3.

| Artboard | F | T | C | D | R | E | A | X | $ | S | U | Falladas | Corrección |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| resumen-m | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ (B-6) | ✔ | ✔ | n/a | ✔ | B-1, B-6, B-9 | Opcionales |
| resumen-e | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ (B-6, B-13) | n/a | ✔ | ✔ | ✔ | B-1, B-5, B-6, B-9, B-10, B-13 | Opcionales |
| detalle-m | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | n/a | ✘ A-3 | **A-3**, B-1, B-3, B-9 | D-2 l. 103 |
| detalle-e | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ (B-13) | n/a | ✔ | ✔ | ✘ A-3 | **A-3**, B-1, B-3, B-5, B-9, B-10, B-13 | D-2 l. 142 |
| tendencias-m | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | n/a | ✔ | B-1, B-4, B-9 | Opcionales |
| tendencias-e | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ (B-13) | n/a | ✔ | ✔ | ✔ | B-1, B-4, B-9, B-10, B-13 | Opcionales |
| calendario-m | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | n/a | ✔ | B-1, B-9 | Opcionales |
| calendario-e | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ (B-13) | n/a | ✔ | ✔ | ✔ | B-1, B-9, B-10, B-13 | Opcionales |

Notas de la tabla:
- **F.** Las líneas obligatorias están en los 8 artboards: `<!doctype html>` en la l. 1, `<html lang="es-MX">`, `<meta charset="utf-8">`, `<script src="./support.js"></script>`, `<link rel="stylesheet" href="../ds/finanzas-hogar/tokens.css">`, `<x-dc>` con `<helmet><style>` primero, la regla `body{margin:0;background:var(--…);color:var(--…);font-family:var(--…)}` y `a{}`/`a:hover{}`. `data-dc-script` viene después de `</x-dc>`. El `data-props` parsea y `$preview` es igual a `w`×`h` del fragmento (detalle-m 412×4298 y detalle-e 1440×3032 incluidos). `html.parser` no deja etiquetas abiertas ni cierres cruzados. Los 8 nombres cumplen `[a-z0-9.-]+\.dc\.html` y sus stems no se repiten en otras áreas del lienzo. El móvil y el escritorio de cada pantalla comparten `<title>` («Inicio · Detalle», etc.); ver Supuestos en el reporte.
- **T.** 0 hex, 0 `rgb()/hsl()/oklch()` y 0 nombres de color fuera de `data-props`. Hay entre 54 y 70 `var()` distintos por archivo, 86 en total, todos definidos en `tokens.json` y en `tokens.css`.
- **X.** En móvil, cada `a`, `button` y `role="tab"` declara `min-height` o `height` ≥ 44 px (`--size-touch-min`, `--size-agenda-row` y `--size-bottombar-item-height`). Las únicas excepciones son las 6 tarjetas-enlace del Resumen, que miden lo que su contenido, de 4 a 7 líneas.
- **$.** 82 montos distintos con el formato `$1,234.56`, todos con `tabular-nums` por cascada (`.n`, `.ax`). No hay guion ASCII antes de `$`, ni fechas ISO visibles, ni emoji. De la lista de U5-07 solo aparece «Error» (B-10).
- **E.** Cada `href` apunta a `#` o a un archivo de la lista del lienzo con el mismo ancho: shell, presupuesto, gastos, comidas, compras, créditos, reportes y `ajustes-p26-s01/-s02/-s03`, que existen en disco.

## 4. Fidelidad (pantalla → U1 y U3)

| Pantalla | Elemento de U1/U3 | Estado |
|---|---|---|
| Resumen | Exactamente 6 tarjetas en el orden fijo de U1 §2.1 (Mes vs plan, Flujo quincenal, Próximo pago y disponible, Comidas de la semana, Mejor canasta, Ahorro) | ✔ 6 `a.card` en m y en e |
| Resumen | Métricas de las 6 tarjetas (T1 $7,060.22; T2 $6,045.50 al cierre; T3 $1,861.58 y disponible $6,881.60 · Lectura B; T5 «Bodega Aurrera $1,551.70» y «$86.60 menos que Walmart»; T6 $41,730.00 y 42 % entero) | ✔ trazadas al JSON (§7) |
| Resumen | Estados de U1 §2.2: vacío por tarjeta con su acción, esqueleto de las 6 sin `$0.00`, error por tarjeta con «Reintentar» | ✔ fila «Estados» en el escritorio (B-5, B-10, B-13) |
| Detalle | D-1 a D-7 en orden, con «Ver todo» hacia su destino | ✔ |
| Detalle | D-2 «Quincena en curso»: 7 propiedades (ingresos con «Esperado», pagos de crédito, fijos y venideros, variable estimado, ahorro planificado, saldo de la quincena, saldo acumulado) | ✘ **A-3**: 6 de 7 correctas; el saldo acumulado no sigue S2 §3.11 |
| Detalle | D-3 a D-7 (pagos y disponible, comidas de hoy, diferencias de precio, metas con un decimal, participación de Karen con la nota fija) | ✔ (B-3) |
| Tendencias | T-1, T-2, T-3, T-4a y T-4b en `figure`/`svg role="img"` con `aria-label`, resumen y «Ver como tabla» | ✔ (B-4) |
| Calendario | Móvil: agenda de hoy a hoy + 30 con separador de quincena y «Ver 30 días más». Escritorio: mes con panel del día | ✔ |
| Estados (U1 §2.2) | Vacío, carga y error compactos en los 4 escritorios | ✔ (B-5, B-10, B-13) |

## 5. Contraste (U5-04 y P3-07)

Pares de texto medidos en el markup expandido con los hex de `tokens.json`: 18 por tema. **0 bajo 4.5:1.** Los mínimos son `ui-role-income` sobre `ui-surface-card` en claro (4.80) y `ui-role-expense` sobre `ui-surface-card` en oscuro (5.90). En UI, 53 pares por tema, y **0 bajo 3:1 fuera de los decorativos** (§2).

| Par | Claro | Oscuro | Umbral |
|---|---|---|---|
| `ui-text-default` / `ui-surface-page` | #36363a / #fafafb = 11.53 | #f2f2f4 / #222325 = 14.07 | 4.5 |
| `ui-text-default` / `ui-surface-card` | 10.76 | 10.76 | 4.5 |
| `ui-text-default` / `ui-surface-selected` | 10.73 | 10.90 | 4.5 |
| `ui-text-secondary` / `ui-surface-card` (incluye el `caption` nuevo de D-2) | #66676e / #f2f2f4 = 5.04 | #d1d2d8 / #36363a = 7.97 | 4.5 |
| `ui-text-secondary` / `ui-surface-page` | 5.40 | 10.43 | 4.5 |
| `ui-text-link` / `ui-surface-card` | #5d609f / #f2f2f4 = 5.16 | #c9cfff / #36363a = 7.94 | 4.5 |
| `ui-role-income` / `ui-surface-card` | #457723 / #f2f2f4 = 4.80 | #98ca7b / #36363a = 6.34 | 4.5 |
| `ui-role-expense` / `ui-surface-card` | #a04b20 / #f2f2f4 = 5.33 | #f8a078 / #36363a = 5.90 | 4.5 |
| `ui-role-savings` / `ui-surface-card` | #00767b / #f2f2f4 = 4.84 | #3fd0d7 / #36363a = 6.42 | 4.5 |
| `ui-role-info` / `ui-role-info-surface` | 4.96 | 6.12 | 4.5 |
| `ui-role-warning` / `ui-role-warning-surface` | 5.12 | 6.09 | 4.5 |
| `ui-role-success` / `ui-role-success-surface` | 4.81 | 6.19 | 4.5 |
| `ui-text-on-primary` / `ui-action-primary-bg` | #ffffff / #5d609f = 5.77 | #1e203a / #c9cfff = 10.47 | 4.5 |

Las cifras coinciden con la vuelta 3 y con `contraste.csv` (redondeo ≤ 0.01). El cambio de esta vuelta, «—» y el `caption` en `ui-text-secondary` sobre tarjeta, no agrega pares nuevos.

## 6. P3-08: decisiones

| Punto | Estado |
|---|---|
| Resumen con ≤ 6 tarjetas | ✔ exactamente 6, en el orden fijo |
| 14 secciones de Ajustes | n/a en Inicio; los enlaces a P-26.01, P-26.02 y P-26.03 existen |
| Barra lateral colapsable | ✔ sin cambio desde la vuelta 3 |
| Español de México | ✔, salvo «Error» (B-10, baja) |
| Dinero `$1,234.56` con cifras tabulares | ✔ |
| Fechas `10-oct-2026` / `Vie 30 oct 2026` | ✔, 0 ISO visibles |
| Paleta de C2 y Notion como referencia | ✔ solo tokens |
| Decisiones y bloqueos | ✔ D-2 ya no decide el ahorro de la quincena en curso; ✘ el saldo acumulado sí decide la pregunta 3 de `S2-ahorro-quincenal` para la 1.ª quincena y agrega un saldo inicial que §3.11 excluye (A-3) |

## 7. Líneas de evidencia

- **P3-06.** Los 8 artboards están en disco y en el fragmento: `boards` = `order` = disco, sin duplicados, con `title`, `page:"inicio"` e `is_interactive:false` y sin `x`/`y`. `expand:"fill"` va solo en los 4 de escritorio, con `w` 412/1440. La copia de la vuelta 3 está en `docs/design/salida/iter-3/inicio/` (9 archivos). `docs/design/salida/iter-4/` no existe: la copia de esta vuelta y de este informe le toca al orquestador.
- **P3-07.** Este informe: tabla artboard × regla (§3) y recálculo de 13 pares (§5).
- **P3-08.** §6.
- **P3-12.**
  - El grep canon de P3-12 (21 cifras, con y sin coma) sobre los 8 artboards, `inicio.json`, `datos-sinteticos.json` y `salida/iter-3/inicio` da **0**. Las mismas cifras en centavos sobre los artboards dan 0, y el grep de credenciales da 0.
  - Fechas: el comando de P3-12 lista 15 coincidencias (días 15, 28 y 30, de oct-2026 a jun-2027). En Inicio aparecen el 30-oct y el 15-nov, con $954.60, $906.98 y $803.27 y partes $572.76, $453.49, $381.84 y $803.27. Ninguna es del reporte ni de Kueski.
  - Revisé el texto visible de resumen-m, detalle-m y calendario-m: sin correos, cuentas ni datos de facturación; solo los nombres Axel, Ximena y Karen.
  - 80 de 82 montos están en el JSON; los otros 2 son las sumas de T-2. Muestra de 10:
    - $6,045.50: `inicio.tarjetas.flujo_quincenal.closing_balance_cents` = `flujo_por_quincena[1].saldo_cents`.
    - $12,827.78: `flujo_por_quincena[1].saldo_acumulado_cents` (A-3).
    - $7,060.22: `presupuesto.resumen.restante_cents`.
    - $1,861.58: `creditos.lineas[0].proximo_pago.monto_cents`.
    - $6,881.60: `creditos.lineas[0].disponible_b_hoy_cents`.
    - $1,062.20: `plan_semana.totales.costo_semana_cents`.
    - $1,551.70: `canasta.por_tienda[1].total_cents`.
    - $86.60: `inicio.tarjetas.mejor_canasta.difference_vs_second_cents`.
    - $41,730.00: `ahorro.ahorrado_total_cents`.
    - $1,402.98: `reporte_mensual.karen.despensa_devengada_mes_cents`.
  - Recálculo de A-3 con §3.11: saldo de la 1.ª quincena −421772 (el JSON trae −521772); acumulado de la quincena en curso **182778**. Ese valor no existe en el JSON.

## 8. Correcciones para la siguiente vuelta (en orden)

1. **A-3 (orquestador, luego autor).**
   - Orquestador: en `proyeccion.flujo_por_quincena`, que ningún `saldo_cents` reste el ahorro planificado ni el gasto variable estimado, y que `saldo_acumulado_cents` parta de 0 en 2026-10-q1, sin `saldo_inicial_cents`. Con eso salen −421772 y 182778 en las dos primeras filas. Hay que avisar a presupuesto (P-07). Si se quiere un saldo inicial, se abre un cambio de contrato por bloqueo.
   - Autor: detalle-m l. 103 y detalle-e l. 142 → «Saldo acumulado $1,827.78».
2. **B-11 y B-12 (orquestador, datos).**
3. Bajas opcionales del autor: B-13, B-9, B-1, B-4, B-5, B-10, B-6 y B-3 (decisión).
4. B-7: el orquestador mide los altos al publicar.
