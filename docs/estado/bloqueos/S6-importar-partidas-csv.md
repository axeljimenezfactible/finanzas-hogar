# Bloqueo S6-importar-partidas-csv · fase WF-3a (agente api-ajustes, 12-oct-2026)

Contexto: `previewPlannedItemsImport` e `importPlannedItems` leen `gastos_fijos_ingresos.csv` con las columnas de plan §8.2 (`tipo, concepto, monto, periodicidad, dia_o_fecha, categoria, persona, porcentaje_conjunto, notas`). C1, S1, U4 y el diccionario no dicen cómo pasa cada columna a `planned_item`.

Qué falta: confirmar la correspondencia (PROVISIONAL, ya implementada):
1. `tipo`: ingreso → income, gasto_fijo → fixed_expense, gasto_venidero → upcoming_expense. `periodicidad`: quincenal → fortnightly, mensual → monthly, anual → yearly, unica → once (sin semanal: la plantilla no la trae).
2. `monto`: pesos con punto decimal y hasta 2 decimales; acepta `$` y comas de miles («15,000.00»); se convierte a centavos sin punto flotante; debe ser > 0.
3. `dia_o_fecha`: AAAA-MM-DD se usa tal cual como `start_date`; un día 1–31 en mensual = la próxima fecha con ese día desde hoy (recortada al último día); en quincenal, vacío o número = el día de pago de la quincena en curso (15 o último); en anual y única se exige la fecha completa.
4. `categoria`: nombre (sin mayúsculas ni acentos) o `code` de una categoría activa; vacía en ingresos.
5. `persona`: axel/ximena → `person_id`; hogar → la regla por defecto (en ingresos, `person_id` vacío = hogar); karen → la regla reutilizable activa llamada «Karen»; conjunto → la que empieza con «Conjunto». La regla se busca por nombre porque `cost_share_rule` no tiene `code`.
6. `porcentaje_conjunto`: solo con conjunto; debe ser igual a la parte de Karen en esa regla (si no, error del renglón; no se crean reglas de un solo uso).
7. `notas` → `notes`. `interval_count` = 1, `accrual_months` = 1, `is_active` = true.
8. Archivo: UTF-8 (se quita el BOM), separador coma o punto y coma (se detecta en el encabezado), encabezados en cualquier orden; `porcentaje_conjunto` y `notas` pueden faltar. Tope: 8 MiB (el de un cuerpo JSON). Solo la parte `file` con `text/csv`; si no, 415.

Opciones:
- A) Aprobar la correspondencia anterior.
- B) Agregar una columna `code` a `cost_share_rule` (migración nueva) y buscar karen/conjunto por clave. Consecuencia: S1 1.5 y C1 cambian.
- C) Que `porcentaje_conjunto` distinto cree una regla de un solo uso (`is_template = 0`). Consecuencia: el CSV crea reglas además de partidas.

PROPUESTA del agente: A, con B si Axel prevé renombrar reglas.

Impacto: solo la importación (N-304, P-26.11). Ninguna cifra canon depende de esto; el CSV real aún no existe (plan §9 punto 9).

- **Estado:** ABIERTO
