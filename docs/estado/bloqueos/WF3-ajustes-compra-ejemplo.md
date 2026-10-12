# Bloqueo WF3-ajustes-compra-ejemplo · intereses de la compra de ejemplo en P-26.05 (vistas previas 6 y 7)

Abierto por `correccion-frontend-areas-r1` (WF-3a, 12-oct-2026). Estado: ABIERTO.

## Qué falta
U4 P-26.05 puntos 6 y 7 fijan la compra de ejemplo (capital y plazo de la última compra de la línea, editables; fecha de esa compra) y que se calcule con `previewCreditPurchase`, pero no dicen con qué intereses se arma el cuerpo cuando la línea **no es revolvente** (tarjeta bancaria, crédito de tienda, préstamo a plazo).
- Línea revolvente (Kueski): HECHO y firme. Es una compra planeada, así que usa `interest_source = term_table` (RN-08), igual que la vista previa 4 de P-26.02. Con la semilla da 30-oct, 15-nov, 30-nov y 15-dic-2026, $973.75 cada uno, total $3,895.00.
- Otras líneas: en uso, la opción A, marcada PROVISIONAL en la pantalla.

## Opciones
- **A (en uso, PROVISIONAL):** copiar los términos de la última compra de la línea: `interest_source`, `rate_num`, `rate_den`, `rate_basis`, `interest_includes_vat` y `rate_decimals`. `interest_cents` solo se copia si el capital y el plazo no cambiaron.
- **B:** usar siempre `term_table` y, si la línea no tiene tabla, mostrar «Sin tasa para calcular» y no guardar.
- **C:** en líneas no revolventes, mostrar solo las fechas (con la frecuencia) sin montos.

## Impacto
Solo afecta los montos de la vista previa en líneas no revolventes. No guarda nada y no toca cifras de los motores. Si el motor S3 rechaza el cuerpo (422), la vista previa sale con error y el campo no se guarda (U4 §0.2).
