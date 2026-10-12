# Bloqueo WF3-comidas-compras-interfaz (PROVISIONAL) · front-comidas-compras · 12-oct-2026

Puntos de P-10 a P-16 que el contrato o U1 no fijan. Se implementó la opción marcada «PROVISIONAL»; nada bloquea la entrega.

| # | Qué falta | Opciones | PROVISIONAL hoy | Impacto |
|---|---|---|---|---|
| 1 | Cómo pasa P-13 «Terminar compra» la semana a P-05 (N-073, prellenado) | a) `?semana=2026-W42` en la ruta de P-05; b) estado compartido | a) `rutaDe('P-05')?semana=…` | Lo debe leer el área de gastos (P-05) |
| 2 | Cómo pasa P-14 «Simular pago con crédito» el total a P-20 | a) `?monto=<centavos>`; b) otro nombre | a) `?monto=309450` | Lo debe leer el área de créditos (P-20) |
| 3 | Marca de color de tienda (K41): `Store` de C1 no trae código | a) agregar `code` a `Store` en C1; b) derivarlo del nombre | b) nombre normalizado (`walmart`, `bodega-aurrera`, `sams-club`; mercados = `market`; sucursal = la de su cadena) | Una tienda renombrada pierde su color |
| 4 | «Sugerir con IA» (P-11) solo si la IA está activa en P-26.11: P-11 no tiene lectura de ese estado en sus extractos | a) leer el ajuste de IA de P-26.11; b) mostrarla siempre y respetar el 409 `ai_disabled` | b) en «⋯»; con IA apagada se muestra el mensaje de C1 | La acción se ve aunque la IA esté apagada |
| 5 | Unidad de `quantity_milli` al registrar un sobrante con `recipe_id` (N-184) | a) porciones × 1000 (como `leftovers.servings_milli`); b) gramos | a) | Si es b), el sobrante queda mal medido |
| 6 | Semana visible en la dirección | a) `?semana=AAAA-Www` (se quita si es la vigente); b) otra llave | a), junto con `WF3-frontend-rutas` | Ninguno fuera del área |
| 7 | Historial de precios (P-15 Tendencias) en fechas sin captura de una tienda | a) se mantiene el último precio capturado; b) hueco | a) | Visual; las cifras capturadas no cambian |
