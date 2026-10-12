# P3 · Vuelta 3 · Área «inicio» · Informe de revisión (P3-07)

Revisor: `revisor-inicio-r3` (independiente, vuelta 3 de 3, no corrige). Fecha: 12-oct-2026.

**Objeto.** Los 8 artboards de `docs/design/maquetas/lienzo/project/inicio/` (P-01 Inicio: Resumen, Detalle, Tendencias y Calendario, cada una en móvil 412 y escritorio 1440) y el fragmento `docs/design/maquetas/lienzo/fragmentos/inicio.json`.

**Criterios.** P3 1.3 (P3-06, P3-07, P3-08, P3-12); U5-04, U5-06 y U5-07 en lo que se ve en el markup; DESIGN.md (§7.6, §7.7, §8.1, §9.3, §10.5 y `contraste.csv`); U1 §2.1 y §2.2 (filas de P-01); U3 §1 a §8; `referencia-tipo/format.md` (leído entero, tomado como dato); `datos-sinteticos.json` y `LEEME.md`; contrato `docs/contratos/motores/S2-interfaz.md` 1.3 §3.11 para el hallazgo nuevo.

**Método.** Scripts propios en `datos/tmp/rev-inicio-r3/` (borrados al terminar), todos los de Python con `python3 -I`:
- `rv.js`: corre en `node` el `renderVals()` de cada artboard con un `DCLogic` de prueba (8/8 sin error).
- `rev.py`: forma del tipo, `html.parser` con pila de etiquetas, hex y funciones de color fuera de `data-props`, existencia de cada `var(--x)` en `tokens.json`, grep P3-12, enlaces contra la lista del lienzo y el disco, raíz móvil/escritorio, fragmento, nombres. Expande `sc-for` y `sc-if` con los valores reales de `renderVals()` y resuelve una cascada simple (clases del helmet, especificidad, estilo en línea, herencia de `color` y fondo efectivo) para el contraste, el dinero, el inglés, los emoji, los nombres accesibles y las alturas táctiles.
- `c10.py`: recálculo independiente de 13 pares con los hex de `tokens.json` en claro y oscuro.
- `txt.py`: volcado del texto visible para la tabla de fidelidad.
- `diff` contra la copia de la vuelta 2 (`docs/design/salida/iter-2/inicio/`).

No usé navegador, red ni la herramienta Artifact.

**Veredicto: NO APROBADO.** 1 hallazgo alto (A-3, nuevo), 0 medios y 10 bajos.
- El alto de la vuelta 2 (A-2) quedó **cerrado**.
- A-3 es de la misma familia que A-1 y A-2: el artboard sigue al archivo de datos y el archivo de datos contradice una regla de U3 y del contrato S2. Las vueltas 1 y 2 no lo detectaron.

## 1. Cierre de los hallazgos de la vuelta 2

`diff` contra `iter-2`: solo cambiaron `inicio-p01-detalle-m` (l. 151-152 y 157-158) y `inicio-p01-detalle-e` (l. 190-191 y 196-197). Los otros 6 artboards y `inicio.json` son idénticos.

| Id | Sev. v2 | Estado | Evidencia (vuelta 3) |
|---|---|---|---|
| A-2 | alta | **Cerrado** | D-6 «Metas» dice «42.4 % · $6,150.00 de $14,500.00» y «78.3 % · $14,100.00 de $18,000.00», con `aria-valuenow="42.4"`/`"78.3"`, `aria-valuetext="42.4 % de la meta"`/`"78.3 % de la meta"` y barra de 42.4 % y 78.3 %. Recalculado: 615000 ÷ 1450000 = 42.41 → 42.4; 1410000 ÷ 1800000 = 78.33 → 78.3 (ROUND_HALF_UP a décimas). La tarjeta 6 sigue con «42 %» entero (resumen-m l. 107, resumen-e l. 144), como pide Q2-17 C |
| B-1 | baja | Abierto | `.ic{…stroke-width:1.5…}` en los 8 helmets; a 20 px DESIGN pide `--icon-stroke-user-units-sm` (1.8) para que el trazo se vea de 1.5 px. El encargo del lienzo dice «stroke 1.5», así que queda en baja |
| B-3 | baja | Abierto | «viejo» aparece 0 veces en D-5 (detalle-m y detalle-e). Falta la decisión |
| B-4 | baja | Abierto | tendencias-m l. 97-98 y tendencias-e l. 133-134: «Meta 2,364», «Meta 1,900» y «Meta (mismo color, punteada)»; U3 T-4b pide «Meta de [nombre]» |
| B-5 | baja | Abierto | Ejemplos de error sin el enlace de la tarjeta o panel (resumen-e l. 169-175, detalle-e l. 240-246); U3 §1.3: «La tarjeta conserva su título y su enlace» |
| B-6 | baja | Abierto | `aria-label` en un `<span>` sin rol para «Pendiente confirmar» (resumen-m l. 83, resumen-e l. 120) |
| B-7 | baja | Abierto (orquestador) | Altos `h` estimados sin navegador; se miden al publicar |
| B-9 | baja | Abierto | `.bd .ic{width:var(--space-4);height:var(--space-4)}` en los 8 (l. 28): icono de insignia a 16 px con un token de espaciado |
| B-10 | baja | Abierto | Etiqueta «Error» en los ejemplos de estado (resumen-e l. 170, detalle-e l. 241, tendencias-e l. 143, calendario-e l. 146); el script de U5-07 la cuenta como inglés |
| B-11 | baja (datos) | Abierto (orquestador) | `datos-sinteticos.json` sigue con `ahorro.metas[].estado` = `atrasada`/`en_ruta` y `inicio.tarjetas.ahorro.next_goal.status` = `behind` (etiqueta «Atrasada»). Inicio ya dice «En curso» |

## 2. Hallazgos de esta vuelta

| Id | Sev. | Regla | Artboards (líneas) | Hallazgo | Corrección |
|---|---|---|---|---|---|
| **A-3** | **alta** | U3 §3, fila D-2 (l. 186); `S2-interfaz.md` 1.3 §3.11 (l. 258-262: `variable_estimated_cents` y `savings_planned_cents` = `null`, «no entra a ningún saldo»; `balance_cents` = ingresos − `spent_cents` − `committed_cents`); bloqueo `docs/estado/bloqueos/S2-ahorro-quincenal.md` **ABIERTO** (pregunta 3: «¿resta del saldo?»); U1 §2.2 P-01 («nunca $0.00 provisional»); P3-12 («para que Inicio, Detalle y Reportes cuadren entre sí») | detalle-m l. 100-103; detalle-e l. 139-142 | D-2 «Quincena en curso» muestra «Gasto variable estimado $0.00», «Ahorro planificado $1,000.00», «Saldo de la quincena $5,045.50» y «Saldo acumulado $11,827.78». U3 D-2 dice que el ahorro planificado llega en `null`, que no se resta de `balance_cents` ni de `cumulative_balance_cents` y que «mientras llegue `null`, la fila muestra «—», nunca un monto en cero». El gasto variable estimado también llega en `null` en S2 1.3. El artboard muestra un monto, lo resta del saldo y así decide la pregunta 3 del bloqueo, que sigue abierta. Además, Inicio no cuadra consigo misma: para la misma quincena (16 al 31 oct), la tarjeta 2 del Resumen dice «$6,045.50 al cierre» (= $15,520.00 − $5,589.55 − $3,884.95), y D-2 dice $5,045.50. La diferencia es justo el ahorro planificado. **Origen**: `datos-sinteticos.json` `proyeccion.flujo_por_quincena[1]` trae `variable_estimado_cents` = 0, `ahorro_planificado_cents` = 100000 y `saldo_cents` = 504550 (con el ahorro restado), mientras `inicio.tarjetas.flujo_quincenal.closing_balance_cents` = 604550 | **Orquestador** (primero): en `proyeccion.flujo_por_quincena`, `variable_estimado_cents` y `ahorro_planificado_cents` en `null` (o un marcador de «pendiente»), y `saldo_cents`/`saldo_acumulado_cents` sin restar el ahorro. Con los datos actuales salen $6,045.50 (= tarjeta 2) y $13,827.78 (= $12,000.00 − $4,217.72 + $6,045.50). Avisar a presupuesto (P-07), que usa la misma serie. **Autor**: en los 2 artboards, «Gasto variable estimado —» y «Ahorro planificado —», con el motivo en `caption` (p. ej. «pendiente de definir»), «Saldo de la quincena $6,045.50» y el saldo acumulado del dato corregido |
| B-12 | baja (datos) | P3-12 (un solo archivo de datos) y Q2-17 C | `datos-sinteticos.json` `ahorro.metas[]`; `LEEME.md` l. 12 (regla 4) | Las décimas de D-6 (42.4 y 78.3) son correctas, pero no están en el JSON (solo `avance_pct` = 42 y 78). `LEEME.md` regla 4 sigue diciendo «Los porcentajes de avance son enteros». Si otra área (P-09, P-26.13) muestra Metas, seguirá la regla vieja | Orquestador: agregar `progress_tenths_pct` (424 y 783) y corregir la regla 4: «entero en la tarjeta de Inicio y un decimal en Metas» |

Otras comprobaciones sin hallazgo:
- Las sumas de T-2 ($242,520.65 y $257,740.56) no están en el JSON como valor. Son la suma de los 12 meses de `proyeccion.meses[].devengado_total_cents` y `efectivo_total_cents`, y la diferencia final (−$15,219.91) es `diferencia_acumulada_cents` del último mes. U3 §4 permite «sumas de lo que está a la vista».
- Los ejes de las gráficas usan formato compacto («$10 mil», «$1 mil»), que U3 §4 permite «solo en ejes». Llevan `font-variant-numeric: tabular-nums` por las clases `.ax` y `.lb`.
- El único `color-mix` es `.ar1` (área de T-2, 20 % de `--ui-chart-1` sobre `--ui-surface-card`), que DESIGN §8.1 permite expresamente.

## 3. Tabla artboard × regla

Reglas: F = forma del tipo (líneas obligatorias, `data-props` válido y `$preview` = fragmento, `html.parser` sin etiquetas abiertas, `lang="es-MX"`, `<title>`); T = solo tokens (0 hex, todos los `var()` en `tokens.json`); C = contraste (U5-04); D = P3-12; R = raíz (móvil 412 fija con `min-height` ≥ 915; escritorio fluido con `flex:999 1 560px`, `max-width:76rem` y `expand:"fill"`); E = enlaces; A = accesibilidad (controles con nombre, sin anidar, iconos `aria-hidden`, sin emoji ni barra de estado falsa); X = táctil ≥ 44 px (móvil); $ = dinero `$1,234.56`, tabular, signo U+2212; S = estados de U1 §2.2 (escritorio); U = fidelidad con U1 y U3.

| Artboard | F | T | C | D | R | E | A | X | $ | S | U | Falladas | Corrección |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| resumen-m | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ (B-6) | ✔ | ✔ | n/a | ✔ | B-1, B-6, B-9 | Opcionales |
| resumen-e | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ (B-6) | n/a | ✔ | ✔ | ✔ | B-1, B-5, B-6, B-9, B-10 | Opcionales |
| detalle-m | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | n/a | ✘ A-3 | **A-3**, B-1, B-3, B-9 | D-2 (l. 100-103) |
| detalle-e | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | n/a | ✔ | ✔ | ✘ A-3 | **A-3**, B-1, B-3, B-5, B-9, B-10 | D-2 (l. 139-142) |
| tendencias-m | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | n/a | ✔ | B-1, B-4, B-9 | Opcionales |
| tendencias-e | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | n/a | ✔ | ✔ | ✔ | B-1, B-4, B-9, B-10 | Opcionales |
| calendario-m | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | n/a | ✔ | B-1, B-9 | Opcionales |
| calendario-e | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ | n/a | ✔ | ✔ | ✔ | B-1, B-9, B-10 | Opcionales |

Notas de la tabla:
- F: los 8 tienen las 9 líneas obligatorias exactas, `data-props` parsea, `$preview` coincide con `w`×`h` del fragmento y el parser no deja etiquetas abiertas. Los títulos son el nombre de la pantalla («Inicio · Resumen», etc.); el móvil y el escritorio de una misma pantalla comparten título, como pide el encargo («nombre de la pantalla»). Los stems son únicos y los nombres cumplen `[a-z0-9.-]`.
- T: 0 hex literales fuera de `data-props`; de 54 a 70 `var()` distintos por archivo, todos en `tokens.json` (`font-sans` sale de `type.families.sans`).
- X: en móvil, todos los `a`, `button` y `role="tab"` declaran `min-height` o `height` ≥ 44 px, salvo las 6 tarjetas-enlace del Resumen, que miden lo que su contenido (4 a 7 líneas, muy por encima de 44 px).
- $: 0 montos mal formados (fuera de los ejes compactos permitidos), 0 guiones ASCII antes de `$` y 0 montos sin cifras tabulares. 0 fechas ISO visibles. 0 emoji. Inglés: solo «Error» (B-10).
- E: todos los `href` apuntan a archivos de la lista del lienzo que existen en disco, con el mismo ancho, o a `#`.

## 4. Fidelidad (pantalla → U1 y U3)

| Pantalla | Elemento de U1/U3 | Estado |
|---|---|---|
| Resumen | Exactamente 6 tarjetas, en orden 1 a 6, con icono, título, valor principal, secundario de 412, insignias y enlace único (U3 §2) | ✔ (6 `a.card` en m y e) |
| Resumen | T1: «$7,060.22 te quedan», «$20,525.78 de $27,586.00 · 74 %», «Faltan 11 días del mes» (hoy 21-oct), sin insignia (74 %) | ✔ |
| Resumen | T2: «Del 16 al 31 oct», «$6,045.50 al cierre», «Entra +$15,520.00» ($14,780.00 + extra $740.00), «Ya salió −$5,589.55 · Por pagar −$3,884.95», «Ingreso esperado» | ✔ |
| Resumen | T3: $1,861.58 = $954.60 + $906.98, «Vie 30 oct 2026 · Kueski Pay», «Hogar $1,026.25 · Karen $835.33», «Faltan 9 días · 2 pagos este día», «Disponible hoy $6,881.60 · Lectura B» y «Pendiente confirmar» | ✔ (B-6) |
| Resumen | T4, T5 y T6 con sus métricas; T6 «42 %» entero y «En curso» | ✔ |
| Detalle | D-1 a D-7 en orden, cada uno con «Ver todo» a su destino; barra «Filtrar» y selector «Mes · 1.ª quincena · 2.ª quincena»; enlace final «Ver todos los gastos con estos filtros» → P-02 | ✔ |
| Detalle | D-2 «Quincena en curso»: 7 propiedades | ✘ **A-3** (ahorro planificado y gasto variable con monto en vez de «—»; saldo distinto de la tarjeta 2) |
| Detalle | D-3: «Vencidos: ninguno», 5 pendientes de hoy al 20 dic 2026 (hoy + 60), las dos lecturas lado a lado con «Activa» y mínimo del horizonte, aviso de reposición | ✔ |
| Detalle | D-4: 4 tiempos × Axel y Ximena, «Fuera» con «—» | ✔ |
| Detalle | D-5: 5 artículos y total por tienda | ✔ (B-3: sin marca «viejo») |
| Detalle | D-6: un decimal (A-2 cerrado), aportación necesaria, fecha meta, «En curso» | ✔ |
| Detalle | D-7: las 5 partes, nota fija exacta de N-026, enlace a P-25 | ✔ |
| Tendencias | T-1 barras agrupadas «Real»/«Plan», T-2 área acumulada con «Hoy», proyección desde el mes 7 y «Pagado por adelantado», T-3 líneas de deuda viva, T-4a y T-4b; cada una en `figure` con `figcaption`, resumen de una frase y «Ver como tabla» | ✔ (B-4) |
| Calendario | Móvil: agenda de hoy a hoy + 30 («De hoy al Vie 20 nov 2026»), separador de quincena, «Ver 30 días más», pagos con «Hogar · Karen» y estado. Escritorio: mes con panel del día | ✔ |
| Estados (U1 §2.2) | Vacío, carga y error compactos en los 4 escritorios, con los textos de U3 §7 | ✔ (B-5, B-10) |

## 5. Contraste (U5-04 y P3-07)

Pares de texto evaluados: 22 combinaciones distintas de texto y fondo efectivo por tema (44 en total). **0 bajo 4.5:1.** El mínimo es `ui-role-income` sobre `ui-surface-card` en claro, 4.80. Pares de UI: 43 por tema. **0 bajo 3:1**, salvo los decorativos que DESIGN exceptúa: `ui-border-divider`, `ui-viz-grid` y la separación de 2 px en `ui-surface-card` entre barras (DESIGN §8.1).

Recálculo independiente con los hex de `tokens.json` (claro | oscuro):

| Par | Claro | Oscuro | Umbral |
|---|---|---|---|
| `ui-text-default` / `ui-surface-page` | #36363a / #fafafb = 11.53 | #f2f2f4 / #222325 = 14.07 | 4.5 |
| `ui-text-secondary` / `ui-surface-card` | #66676e / #f2f2f4 = 5.04 | #d1d2d8 / #36363a = 7.97 | 4.5 |
| `ui-text-link` / `ui-surface-card` | #5d609f / #f2f2f4 = 5.16 | #c9cfff / #36363a = 7.94 | 4.5 |
| `ui-role-income` / `ui-surface-card` | #457723 / #f2f2f4 = 4.80 | #98ca7b / #36363a = 6.34 | 4.5 |
| `ui-role-expense` / `ui-surface-card` | #a04b20 / #f2f2f4 = 5.33 | #f8a078 / #36363a = 5.90 | 4.5 |
| `ui-role-savings` / `ui-surface-card` | #00767b / #f2f2f4 = 4.84 | #3fd0d7 / #36363a = 6.42 | 4.5 |
| `ui-role-info` / `ui-role-info-surface` | #006fa4 / #e7f5ff = 4.96 | #6bc3fd / #003b5a = 6.12 | 4.5 |
| `ui-role-warning` / `ui-role-warning-surface` | #8b5d00 / #fff0dc = 5.12 | #e5ae5b / #4c3000 = 6.09 | 4.5 |
| `ui-role-success` / `ui-role-success-surface` | #207b3e / #e3fae6 = 4.81 | #7ecf90 / #0a421e = 6.19 | 4.5 |
| `ui-text-on-primary` / `ui-action-primary-bg` | #ffffff / #5d609f = 5.77 | #1e203a / #c9cfff = 10.47 | 4.5 |
| `ui-border-control` / `ui-surface-card` | 3.47 | 6.09 | 3 (UI) |
| `ui-person-axel` / `ui-surface-card` | 11.93 | 3.10 | 3 (UI) |
| `ui-chart-2` / `ui-surface-card` | 3.15 | 7.73 | 3 (UI) |

Cruce con `docs/design/contraste.csv`: «Texto secundario sobre tarjeta» 5.03, «Cifra de ingreso sobre tarjeta» 4.79, «Cifra de ahorro sobre tarjeta» 4.84 e «Insignia de información» 4.96. Las diferencias son de redondeo (≤ 0.01).

## 6. P3-08: decisiones

| Punto | Estado |
|---|---|
| Resumen con ≤ 6 tarjetas | ✔ exactamente 6, en el orden fijo |
| 14 secciones de Ajustes | n/a en Inicio; los enlaces a P-26.02 y P-26.03 existen |
| Barra lateral colapsable | ✔ botón «Colapsar barra lateral» con `aria-expanded="false"` en los 4 escritorios; en móvil, menú ☰ hacia `shell-barra-lateral-m` |
| Español de México | ✔ salvo «Error» (B-10, baja) |
| Dinero `$1,234.56` con cifras tabulares | ✔; `grep -rnoE '\$[0-9]{4,}\.[0-9]{2}'` sobre los 8 = 0 |
| Fechas `10-oct-2026` | ✔ («17-ene-2027», «19-oct-2026», «Vie 30 oct 2026»); 0 fechas ISO visibles |
| Paleta de C2 (índigo 280°) | ✔ solo tokens; un `color-mix` permitido |
| Notion como única referencia | ✔ superficies planas, `shadow-0`, sin referencias ajenas |
| Decisiones de la puerta 2 | ✔ Q2-17 C (A-2 cerrado) y Q2-18 A; ✘ D-2 decide una pregunta abierta de `S2-ahorro-quincenal` (A-3) |

## 7. Líneas de evidencia

- **P3-06**: los 8 artboards están en disco y en el fragmento: `boards` = `order` = disco, sin duplicados, con `title`, `page:"inicio"` e `is_interactive:false`, sin `x`/`y`, y `expand:"fill"` solo en los 4 de escritorio. Faltantes: 0. La copia de la vuelta 2 está en `docs/design/salida/iter-2/inicio/` (9 archivos). `docs/design/salida/iter-3/` está vacía: la copia de esta vuelta y de este informe le toca al orquestador.
- **P3-07**: este informe, con la tabla artboard × regla (§3) y el recálculo de 13 pares (§5).
- **P3-08**: §6; grep de montos sin separador de miles = 0.
- **P3-12**: grep de las 21 cifras canon (con y sin coma) sobre los 8 artboards, `inicio.json` y `datos-sinteticos.json` = **0**. Grep de credenciales = 0. Revisión visual de resumen-m, detalle-m y calendario-m: sin correos, cuentas ni datos de facturación; solo los nombres Axel, Ximena y Karen. Muestra de 12 montos trazados al JSON: $7,060.22 (`presupuesto.resumen.restante_cents`), $6,045.50 (`inicio.tarjetas.flujo_quincenal.closing_balance_cents`), $1,861.58 (`creditos.lineas[0].proximo_pago.monto_cents`), $6,881.60 (`creditos.lineas[0].disponible_b_hoy_cents`), $1,062.20 (`plan_semana.totales.costo_semana_cents`), $1,551.70 (`canasta.por_tienda[1].total_cents`), $86.60 (`inicio.tarjetas.mejor_canasta.difference_vs_second_cents`), $41,730.00 (`ahorro.ahorrado_total_cents`), $1,402.98 (`reporte_mensual.karen.despensa_devengada_mes_cents`), $3,802.33 (`reporte_mensual.karen.deuda_viva_cents`), $732.90 (`salud.tendencia_semanal[11].costo_cents`) y $2,783.34 (`ahorro.metas[0].aportacion_necesaria_cents`). Todos los montos de los artboards están en el JSON, salvo los ejes compactos y las dos sumas de T-2 (§2). Fechas: el comando de P3-12 lista 15 fechas coincidentes (días 15, 28 y 30, de oct-2026 a jun-2027). Revisé las de Inicio (30-oct y 15-nov): sus montos son $954.60, $906.98 y $803.27, con partes $572.76, $453.49, $381.84 y $803.27. Ninguno es del reporte ni de Kueski.

## 8. Correcciones para la siguiente vuelta (en orden)

1. **A-3 (orquestador, luego autor).**
   - Orquestador: en `datos-sinteticos.json` `proyeccion.flujo_por_quincena`, `variable_estimado_cents` y `ahorro_planificado_cents` en `null` mientras sigan abiertos `S2-flujo-quincenal` (Q2-12) y `S2-ahorro-quincenal`, y `saldo_cents`/`saldo_acumulado_cents` sin restar el ahorro, como en S2 1.3 §3.11. Avisar al área presupuesto (P-07), que usa la misma serie.
   - Autor: en `inicio-p01-detalle-m` l. 100-103 y `inicio-p01-detalle-e` l. 139-142, «—» con motivo en las dos filas provisionales y el saldo de la quincena igual al de la tarjeta 2 ($6,045.50).
2. **B-12 y B-11 (orquestador).** Décimas de avance y estados de meta en el JSON; regla 4 de `LEEME.md`.
3. Bajas opcionales del autor: B-9, B-1, B-4, B-5, B-10, B-6 y B-3 (decisión).
4. B-7: el orquestador mide los altos al publicar.
