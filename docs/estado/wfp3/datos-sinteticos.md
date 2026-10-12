# Reporte · datos-sinteticos (Pista 3, fase 1) · 11-oct-2026, noche

Agente `datos-sinteticos` (Sonnet). Reanudación del flujo wf_6b8245eb-3e3: se empezó de cero, sin parciales anteriores.

## Hecho
- `docs/design/maquetas/datos-sinteticos.json` (191 KB, 34 llaves de primer nivel): personas, ingresos, 13 categorías, presupuesto, 33 gastos, ahorro (fondo y 2 metas), 4 líneas de crédito (Kueski Pay revolvente con 3 compras vivas y 2 programadas, Banorte Oro, préstamo, crédito Coppel), 3 promociones, comparador con 6 escenarios, 5 tiendas, 33 productos con 105 precios (22 viejos), canasta comparada, 12 recetas, plan semanal 7×4×2 con «fuera», lista de compras por tienda, inventario (14), salud (8 semanas de peso), reportes y avisos, proyección de 24 quincenas y 12 meses, deuda viva, calendario (58 eventos), las 6 tarjetas de Inicio con nombres de campo de C1.
- `docs/design/maquetas/LEEME.md`: esquema, reglas para autores, cifras que deben coincidir, fechas inevitables, qué es ilustrativo, cómo regenerar y comprobar.
- `docs/design/tools/ds/verificar_sinteticos.py` (39 comprobaciones) y, como apoyo declarado, `docs/design/tools/ds/generar_sinteticos.py` (generador determinista; las sumas salen de código, no a mano).

## Verificación
- `python3 -I docs/design/tools/ds/verificar_sinteticos.py` → 39 de 39 en OK y `RESULTADO: OK`.
- Grep P3-12 (patrón completo de la Verifica, sobre el JSON) = 0; las 19 cifras canon en centavos como valor entero = 0; credenciales y correos = 0.
- Comando de fechas de P3-12 (model.json ∪ data.json ∩ JSON): 15 fechas, todas de la regla de quincena (día 15, día 30 o fin de febrero): 2026-10-30, 2026-11-15, 2026-11-30, 2026-12-15, 2026-12-30, 2027-01-15, 2027-01-30, 2027-02-15, 2027-02-28, 2027-03-15, 2027-03-30, 2027-04-15, 2027-04-30, 2027-05-15, 2027-06-15. Ninguna con un monto del reporte ni de Kueski (los pagos de la despensa son otros). Ninguna fecha cae en 9 o 10 después del 1-oct-2026.
- Muestra de cuadre: Mes contra plan $20,525.78 de $27,586.00 (74 %); Mercado y Despensa común en aviso (94 % y 83 %); Vivienda y Mascotas excedidas (100 % y 110 %); próximo pago vie 30 oct con 2 pagos; disponible B $6,881.60 y A $11,440.00.

## Supuestos
1. «Hoy» es el miércoles 21-oct-2026, no la fecha real, para que el mes tenga datos y no se cruce con el 10-oct de la semilla.
2. Las series con historia empiezan antes de 2026-10 (peso desde 31-ago, comidas desde 27-jul, ahorro desde may-2026, tendencia mensual may-sep, préstamo jun-2026, Coppel ago-2026, meses sin intereses jun-2026); la instrucción pedía 2026-10 a 2027-09, pero una gráfica necesita pasado. Ninguna de esas fechas está en `model.json`.
3. Lo ilustrativo (no es regla): lecturas A/B, estado de metas, pago mínimo de la tarjeta, gasto variable estimado ($5,800.00 por quincena), orden de tiendas con cobertura de 90 % o más, sustitutos de tamaño comparable, 12 compras al año en Sam's, devengo de préstamo/tienda/MSI al pagarse, saldo inicial de efectivo de $12,000.00. El «real» del mes excluye venideros. Quedó escrito en LEEME para que los autores no lo tomen como regla de S2/S3/S5.
4. Nombres reales públicos con términos inventados: Kueski Pay, Banorte, Santander, Coppel, Liverpool, Walmart, Bodega Aurrera, Sam's Club, Mercado Portales, Tianguis de la Narvarte. Los datos de salud de Axel y Ximena son inventados.
5. Axel, Ximena y Karen aparecen con su nombre (Q2-98 A). Los usuarios son `axel` y `ximena`; no hay correos.
6. La despensa programada y la compra del comparador son de ejemplo; los escenarios se calculan con la misma función que las compras reales.
7. El generador `generar_sinteticos.py` no estaba en el alcance literal (solo el verificador); se dejó junto a él porque sin generador no se pueden cambiar datos manteniendo las sumas.

## Bloqueos
- Ninguno. Excepción de tamaño (ver Consumo).

## Pendientes
- Revisión visual de Axel o del orquestador de que 33 gastos, 12 recetas y 105 precios alcanzan para las ~25 pantallas; si un autor necesita más datos, se agregan en el generador y se regenera.
- Los apuntadores de `indice_pantallas` son guía, no contrato: los autores pueden pedir ajustes.
- Un posible ajuste de gusto: Axel supera 3,000 kcal en dos días del plan (porciones dobles en comida y cena).

## Consumo aproximado
- Acciones: unas 40. Leído: unos 180 KB (U1 §0 a §3, U3 tarjetas, U4 campos, P3-12, extractos por grep).
- Escrito: unos 290 KB, por encima del tope de 150 KB: el JSON (191 KB) es el entregable y el generador (78 KB) lo produce; sin ellos no hay forma de cumplir «todos los datos de ~25 pantallas coherentes entre sí». Se prefirió terminar completo y declarar el exceso.
