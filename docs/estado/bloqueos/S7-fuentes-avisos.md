# Bloqueo S7-fuentes-avisos · cifras y definiciones que faltan para las alertas

Abierto por s7-avisos (WF-3a) el 12-oct-2026. Criterios: S7-04, S7-08 (la definición de «casillas sin armar», las horas, días, anticipación, umbrales y días antes del pago son FIRMES y no se tocan aquí). Estado: **ABIERTO**.

## 1. Día de despensa (PREGUNTA a Axel)
La alerta (3) sale `alerts.stale_prices_advance_days` días antes del «día de despensa», pero ningún documento ni tabla define ese día. Opciones: A) el lunes que abre cada `shopping_list` (semana K8); B) los días de pago 15 y 30 (fin de mes en febrero); C) una clave nueva `alerts.grocery_weekday` en `app_setting` (migración nueva 00005, K44). Recomendación: C (configurable, como todo parámetro). Mientras tanto `FuenteBD.Despensa` devuelve «sin fuente» y la alerta se omite y queda anotada en el registro; el evaluador está hecho y probado con una fecha fija (matriz de `TestAlerta_Matriz`).

## 2. Funciones de otras áreas (solicitudes fuera de alcance)
- Presupuesto devengado del mes por categoría y total (lo mismo que `getBudgetAlerts`): api-presupuesto exporta una función de dominio y `FuenteBD.Presupuesto` la llama. Hoy devuelve «sin fuente».
- Precios viejos de los productos del plan en la fecha de despensa (`precios.StalePricesForPlan`, S5-12) con el libro cargado desde la base: api-compras-b exporta la carga del `precios.Book`. Hoy devuelve «sin fuente».
Las cifras nunca se inventan: sin fuente, la alerta no sale.

## 3. Qué semana revisa la alerta (4) (PROPUESTA)
El día `alerts.meal_plan_missing_weekday` revisa la semana que empieza el **lunes siguiente** (el domingo avisa de la semana que viene). Alternativa: la semana en curso. Recomendación: la siguiente.

## 4. Cadencia y reintentos (PROPUESTA, K31)
- Los jobs `daily_summary`, `weekly_summary` y `alerts` revisan cada minuto (`* * * * *`) y comparan la hora local con su clave de `app_setting`; así un cambio en Ajustes se respeta al instante y ninguna hora vive en el código. Si el servicio estuvo apagado a la hora, sale al encender ese mismo día. Si Ajustes cambia el horario del job con `updateJob`, el aviso puede retrasarse: PROPUESTA ocultar `schedule` de estos tres jobs en Ajustes o dejarlo de solo lectura.
- Un envío fallido queda `failed` con su código HTTP y no se reintenta solo (evita mandar dos veces si Telegram sí lo recibió). Alternativa: un reintento a los 10 minutos.
- Pagos ya hechos antes de su fecha: la alerta (1) usa `scheduled_payment` sin cruzar `credit_payment`. PREGUNTA: ¿se omite el aviso de un pago ya liquidado?
