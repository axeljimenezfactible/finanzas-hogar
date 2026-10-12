# correccion-inicio-r4 · 12-oct-2026 (vuelta 5)

## Hecho
1. Respaldo previo en docs/design/salida/iter-4/inicio/ (8 artboards de inicio, 2 de P-07 y inicio.json).
2. Hallazgo 1 [alta] (A-3, saldo acumulado de D-2): "Saldo acumulado" de la quincena en curso pasa de $12,827.78 a **$1,827.78** (182778 = -421772 + 604550) en inicio-p01-detalle-m (l. 103) e inicio-p01-detalle-e (l. 142). La nota del panel no se tocó.
3. P-07 (presupuesto-p07-flujo-m y -e), que pintaba la serie vieja: la constante F se reemplazó por la serie de proyeccion.flujo_por_quincena del JSON actual (24 quincenas; variable y ahorro en null, saldo sin ahorro, acumulado desde 0, saldo_inicial_cents = 0). Textos ajustados: KPI "Saldo inicial" ahora dice "El acumulado parte de cero" (antes "Efectivo al 1-oct-2026"); las dos notas de "—" dicen que variable y ahorro están pendientes en todas las quincenas, no se restan de ningún saldo y el acumulado parte de $0.00; en Supuestos, el variable de $5,800.00 queda "solo de referencia (no entra a ningún saldo)". La lógica ya toleraba null.

## Verificación
Script propio (python3 -I, /tmp/verif_r4.py) sobre los 10 archivos: líneas obligatorias, data-props parsea con $preview (412 / 1440), etiquetas balanceadas, 0 hex literales fuera de data-props y SVG, 0 cifras prohibidas P3-12 (texto con/sin coma y centavos), todos los var(--x) existen en tokens.css, sin innerHTML/iframe/red. RESULTADO: OK. Fragmento inicio.json válido y lista exactamente los 8 archivos de inicio (sin cambios de w/h). grep de 12,827 / 13,827 / 11,827 en inicio y presupuesto: 0.

## Supuestos
- Alturas del fragmento sin cambio: los textos modificados tienen largo similar.
- Etiqueta de base "Fijos y variable estimado" en P-07 se mantiene tal cual viene en el JSON (campo base).

## Bloqueos
Ninguno nuevo. S2-ahorro-quincenal sigue abierto (causa de los "—").

## Pendientes
- Otros artboards de presupuesto (P-06, P-09) no se revisaron por quedar fuera de alcance; si muestran ahorro restado de la quincena, avisar al autor de presupuesto.

## Consumo aproximado
~14 acciones, ~60 KB leídos, ~12 KB escritos (P-07 reescrito por script).
