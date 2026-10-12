# correccion-datos-r4 · 12-oct-2026

## Hecho
1. **Quincena en curso** (`proyeccion.flujo_por_quincena[1]`, 2026-10-q2, 16 a 31 oct):
   - `variable_estimado_cents`: 0 → `null`.
   - `ahorro_planificado_cents`: 100000 → `null`.
   - `saldo_cents`: 504550 → **604550** ($6,045.50; igual a `inicio.tarjetas.flujo_quincenal.closing_balance_cents`).
   - `saldo_acumulado_cents`: 1182778 → **1282778** ($12,827.78).
   - Efecto en cadena: cada `saldo_acumulado_cents` de las quincenas `[2]` a `[23]` sube 100000 ($1,000.00); sus demás llaves no cambian. La quincena `[0]` no cambia.
2. **Inventario**: `umbral_reposicion` huevo null → 12 (pza), leche null → 2000 (ml), tortilla null → 1000 (g). Estado de los tres sigue «reponer» (existencias 6, 1500 y 700). Regla del generador ahora: agotado si 0; reponer si hay nivel y existencia ≤ nivel; suficiente en el resto. `inventario.resumen` recalculado: 1 / 5 / 8 (sin cambio). Croquetas, jabón, detergente y papel siguen con nivel vacío y «suficiente».
3. Archivos: `docs/design/tools/ds/generar_sinteticos.py`, `docs/design/tools/ds/verificar_sinteticos.py` (+2 comprobaciones de inventario, +2 de quincena en curso; el chequeo del acumulado tolera null), `docs/design/maquetas/datos-sinteticos.json` (regenerado; generador determinista, verificado con cmp), `docs/design/maquetas/LEEME.md` (v1.3, sección «Notas de lectura (r4)»: las dos llaves van en null por el bloqueo S2-ahorro-quincenal abierto y la pantalla muestra «—»).

## Verificación
- `python3 -I tools/ds/verificar_sinteticos.py` → 61 comprobaciones, 0 FALLA, `RESULTADO: OK` (nuevas: 51, 52, 55, 56). Incluye el grep de P3-12 (cifras prohibidas con y sin coma, y en centavos) = 0.
- Diff de llaves entre el JSON anterior y el nuevo: vacío (mismas llaves). Solo cambian los valores listados arriba (comparación recursiva).
- Regeneración determinista (segunda corrida idéntica byte a byte).

## Supuestos
- La cifra «$13,827.78» del encargo y de la revisión A-3 sale de $12,000.00 − $4,217.72 + $6,045.50, pero el saldo real de la quincena 1 es −$5,217.72 (con el ahorro de $1,000.00 restado). Con el alcance pedido (solo la quincena en curso) el acumulado correcto es **$12,827.78** (= $11,827.78 anterior + $1,000.00). Para llegar a $13,827.78 habría que quitar también el ahorro de la quincena 1 (1-15 oct), que queda fuera del alcance; no se hizo. Los autores deben pintar $12,827.78.
- Niveles de reposición elegidos (12 pza, 2000 ml, 1000 g): arbitrarios, solo cumplen existencia ≤ nivel (U4 P-26.14, K7).
- Las quincenas futuras (`[2]` a `[7]`) siguen restando el ahorro de $1,000.00 y el variable de $5,800.00 en su propio saldo; el bloqueo abierto afectaría también a esas filas, pero el encargo solo pidió la quincena en curso.

## Bloqueos
- Ninguno nuevo. S2-ahorro-quincenal sigue abierto (es la causa del null).

## Pendientes (para otros agentes)
- Autores de inicio (detalle, D-2): «Gasto variable estimado —», «Ahorro planificado —», saldo de la quincena $6,045.50, saldo acumulado $12,827.78.
- Autores de ajustes-b (s14) y compras-p16: huevo «Nivel 12 pza», leche «Nivel 2,000 ml», tortilla «Nivel 1,000 g» en lugar de «Sin aviso»; días restantes sin cambio (2,3 / 3,0 / 3,5 → enteros en pantalla); resumen 1/5/8.
- Autor de presupuesto (P-07): si muestra `saldo_acumulado_cents` de quincenas posteriores, suben $1,000.00; y la fila de la quincena en curso va con «—» en variable y ahorro.

## Consumo aproximado
~25 acciones, ~70 KB leídos (JSON solo por script), ~12 KB escritos.
