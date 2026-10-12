# Bloqueo WF3-inicio-vacios · Reglas de «vacío» y destinos de P-01 que C1 y U3 no fijan

Fase: WF-3a · Construcción · agente `front-inicio` · 12-oct-2026 · Estado: **ABIERTO** (todo lo de abajo está PROVISIONAL en el código, con comentario que cita este archivo).

## Qué falta
U3 define el texto de cada estado vacío, pero en cinco casos C1 v1.4 no trae un campo que diga «está vacío» y la interfaz no debe recalcular (U1 §0.2). Además faltan dos destinos y un texto.

| # | Dónde | Qué falta | Regla PROVISIONAL usada | Opciones |
|---|---|---|---|---|
| 1 | Tarjeta 2 y D-2 | Cómo saber «sin quincena capturada ni esperada» | `income_cents` = 0 | A) esa regla; B) C1 agrega `income_source` (`captured`/`expected`/`none`) |
| 2 | Tarjeta 6 | Cómo saber «sin fondo ni metas» | sin `next_goal`, `saved_total_cents` = 0 y `planned_month_cents` = 0 | A) esa regla; B) C1 agrega `has_goals` |
| 3 | Tarjeta 4 | Vacío = «`unplanned_slots` = 21 y sin casillas con receta»; C1 no trae cuántas casillas tienen receta | `unplanned_slots` ≥ 21 y `cost_cents` = 0 | A) esa regla; B) C1 agrega `planned_slots` |
| 4 | D-6 | Si «cada meta activa» incluye el fondo (`kind` = `fund`) | se muestra todo lo que devuelve `listSavingsGoals?status=active` | A) todo; B) solo `kind` = `goal` |
| 5 | «Ver todos los gastos con estos filtros» (Detalle) | Nombres de los parámetros que recibe P-02 | los nombres de los filtros de C1 (`period`, `person_id`, `family`, `category_id`, `store_id`, `credit_id`) | A) esos; B) los que defina `front-gastos` |
| 6 | Calendario | Destino al tocar un evento que no es pago (vista previa G-06 no existe en la base) | pagos y cortes abren P-18 de su línea; los demás eventos no son enlace | A) así; B) la base agrega G-06 y todos abren su vista previa |
| 7 | T-2, nota de proyección con el interruptor apagado | U3 solo da el texto con el interruptor encendido («Proyección con despensas repetidas») | «Proyección sin despensas repetidas» | A) ese texto; B) otro que dé Axel |
| 8 | Acción «Capturar quincena» (U1 §2.2: «solo los días 15 y 30 si falta») | Qué es «si falta» | el día 15 o 30 (último día si el mes es más corto) con `income_is_expected` = true en `getDashboardFortnightFlow` | A) así; B) otra señal de C1 |

## Impacto
Solo cambia qué texto ve la persona en esos casos; ninguna cifra cambia. Si Axel elige B en 1–3, hace falta una versión menor de C1 y ajustar `Resumen.tsx` (una línea por caso).
