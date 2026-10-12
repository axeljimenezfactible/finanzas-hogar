# Bloqueo S7-reporte-mensual-secciones · secciones del reporte mensual que dependen de otras áreas

Estado: ABIERTO · 12-oct-2026 · agente s7-pdf (WF-3a) · contenido del reporte PROVISIONAL hasta la puerta 4 (K31).

## Qué falta (PREGUNTAS)
1. **Devengado contra efectivo en C1**: `MonthlyReport` (getMonthlyReport) no trae devengado, efectivo ni diferencia acumulada (R-45…R-50); viven en `getAccrualVsCash`, que solo lista **del mes en curso en adelante**. El PDF los muestra calculados con las mismas llamadas al motor S2 y los cita como `accrual:AAAA-MM`. Opciones: A) C1 1.5 agrega `accrual_vs_cash` (renglón `AccrualVsCashMonthsItem`) a `MonthlyReport` (recomendada: la vista web y el PDF quedarían en una sola respuesta); B) la vista web llama además a `getAccrualVsCash` (no sirve para meses cerrados).
2. **Cálculo compartido**: presupuesto, quincenas, créditos y participación externa se arman en `internal/reportes/datos.go` con las mismas llamadas que `getBudgetMonth`, `getBudgetByCategory`, `getBudgetFortnights` (motor S2 vía `dominio.FinEngine`) y con SQL equivalente a `getCreditTotals`; no hay forma de llamar a otra operación sin sesión (el job del día 1 no tiene). Opciones: A) cada área exporta su constructor (`presupuesto.BudgetMonthFor(ctx, db, mes)`…); B) `api.Registry.Invoke(ctx, operationID, params)` interno de backend-base. Mientras: **plan de comidas y canasta no se muestran** (`meals.weeks` vacío, sin `basket`): el reporte no inventa cifras.
3. **Valores PROVISIONAL de esta versión**: ingreso del mes = Σ `income.amount_cents` del mes; la «quincena» de `credits` en un mes cerrado = su segunda quincena; intereses pagados = intereses de pagos programados ya cubiertos; `accrued_groceries_cents` = devengado total de la persona en el mes (motor S2).
4. **Formato de negativos**: `dominio.FormatPesos` da `-$1,197.64` y la expresión de S7-12 pide `$-1,197.64`. El PDF usa el formato compartido (`-$`); la prueba de S7-12 acepta el signo antes de `$`. ¿Cuál vale?

## Impacto
Las cifras visibles cruzan con la API con 0 diferencias (S7-02, S7-11). Si C1 cambia (pregunta 1) o un área cambia su cálculo (pregunta 2), el reporte se ajusta en `datos.go` y `vistas.go`.
