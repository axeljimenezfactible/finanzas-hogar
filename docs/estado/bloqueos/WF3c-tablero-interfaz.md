# Bloqueo WF3c-tablero-interfaz · Puntos del tablero de Inicio que U3 2.0 no cierra para la interfaz

Abierto por `front-tablero` (WF-3c, 12-oct-2026). Cada punto está construido como dice la columna «Hecho (PROVISIONAL)» y sigue lo demás.

| # | Qué falta | Hecho (PROVISIONAL) | Opciones | Impacto |
|---|---|---|---|---|
| 1 | K-4: «y [n] más en Calendario» pide un conteo que `listScheduledPayments` no da (solo `next_cursor`) | Con `next_cursor` ≠ null: «Hay más pagos en Calendario», sin número | A) Así. B) C1 agrega un conteo (p. ej. `total_count` en la lista o `upcoming_count` en `CreditTotals`). C) Una tercera petición con `limit` 200 solo para contar | Texto de una línea de K-4 |
| 2 | K-3: «Disponible hoy —» con «Reintentar» si solo falla `getDashboardNextPayment`; un botón dentro de la tarjeta-enlace es HTML inválido y rompe A-08 (un solo enlace por tarjeta) | «Disponible hoy —» sin botón; el bloque se reintenta al recargar Inicio o al volver la conexión | A) Así. B) Poner «Reintentar» fuera del enlace, bajo la tarjeta | Un botón en un caso de falla parcial |
| 3 | K-1: la insignia «[n] categorías rebasadas» «lleva a G-3», pero la tarjeta entera ya es el enlace a P-06 (A-08) | Insignia de texto; G-3 es el bloque 6 del mismo tablero | A) Así. B) Enlace aparte fuera de la tarjeta | Un toque menos hacia G-3 |
| 4 | B-1: DESIGN §8.9 fija el alto del bloque en `size.chart.height-compact` / `size.chart.height`; el arco escala con el ancho | Arco de ancho completo hasta 320 px (alto ≈ 176 px) más dígito, insignia y líneas; la tarjeta crece con su contenido | A) Así. B) Fijar el alto y recortar líneas secundarias | Alto de la tarjeta B-1 |
| 5 | G-1: «2 px de separación» entre columnas; ECharts 6 da `barGap` en % | `barGap` 8 % | A) Así. B) Calcular el % desde el ancho en cada cambio de tamaño | Visual |
| 6 | G-2 y G-3 como barras HTML (no ECharts) para tener filas-enlace de 44 px con icono `triangle-alert` e insignias de sugerencia (A-20, A-21) | `<figure>` con `<figcaption>`, resumen, «Ver como tabla» y la tabla siempre para lectores de pantalla, como las demás | A) Así. B) ECharts con la lista de enlaces aparte | Ninguno en cifras |
| 7 | K-4: N que no sea entero > 0 en `getDashboardSettings` | Se trata como falla de lectura: K-4 muestra su error con «Reintentar» (U3 §1.2) | A) Así. B) Otra conducta | Solo con datos inválidos |
