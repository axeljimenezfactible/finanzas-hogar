# correccion-datos-r1 · Corrección de las maquetas de datos (hallazgos M-1 a M-4)

Fecha: 11-oct-2026, noche. Alcance: `docs/design/maquetas/`, `docs/design/tools/ds/generar_sinteticos.py` y `verificar_sinteticos.py` (el generador era necesario: es la fuente del JSON). No se tocó nada más. Sin git add/commit.

## Hecho
- **M-1 · Tickets (P-04, P-05).** Llave nueva `tickets` en `datos-sinteticos.json`:
  - `confirmados[0]` (`tk1`, gasto `g13`, Walmart, 12-oct): 11 productos del catálogo y el renglón «Ajuste del ticket» (−$10.60). Cada renglón trae producto, presentación, cantidad, precio unitario (el de `productos`), importe, cuenta y partes hogar/Karen. Suma $2,918.40 = monto de g13; hogar $1,459.20 y Karen $1,459.20, iguales a g13. Totales `por_cuenta`.
  - `lectura_pendiente` (`ocr1`, Walmart, 21-oct): 10 renglones con `texto_leido`, producto y presentación sugeridos, cuenta sugerida y su origen (la regla de N-071), `confianza` (renglones 9 y 10 de confianza baja, con motivo y alternativas; 6 altas, 2 medias, 2 bajas), total leído $535.99 contra suma $551.89, diferencia −$15.90.
  - `por_confirmar` (2, por fecha descendente) y `borradores[0]` (`ocr2`, Bodega Aurrera, 19-oct, con 2 correcciones; total propuesto $373.60, total del borrador $392.20).
  - `indice_pantallas`: `P-04/P-05` apunta a `tickets` (y a `gastos`, `reparto`, `categorias`, `productos`, `tiendas`); `P-02/P-03` quedó aparte.
- **M-2 · Notificaciones.** `ajustes.notificaciones.alertas`: las 4 alertas de la decisión F (`payment_due` [1, 0]; `budget_threshold` 80/100 %; `stale_prices` con `advance_days` = 1, K50, apagada para dibujar ese estado; `meal_plan_missing` domingo), cada una con activa, destinatarios, anticipación y hora 09:00 (K50). `ajustes.datos_semilla` (3 renglones, estado de carga, sin cifras ni la fecha del 10-oct). `ajustes.parametros.comparador_max_s` = 5 y `comparador_max_s_limites` [1, 30] (K50). El puntero de `P-26` ya no está colgante.
  - Efecto lateral para coherencia: el aviso de ejemplo «plan sin armar» del 18-oct pasó de 19:05 a 09:00 (hora de avisos de K50) y se reordenó.
- **M-3 · Nombres de C1.** `proximo_pago.next_payment` (`due_date`, `credit_id`, `amount_cents`, `household_share_cents`, `external_share_cents`, `days_left`, `payment_count`, `is_overdue`; se conserva `dia_semana` y los `_fmt`), `ahorro.next_goal` (`name`, `target_cents`, `target_date`, `saved_cents`, `progress_pct`, `planned_monthly_cents`, `required_monthly_cents`, `status` `behind`; más `id`, `kind`, `etiqueta_estado`) y `mes_vs_plan.month`. `estados_alternos.pago_vencido` también quedó con esos nombres. `ahorro.metas` (P-09) y `segundo_nivel` conservan los nombres en español y el LEEME lo dice.
- **M-4 · Calendario.** `calendario_hogar` va del 2026-10-01 al 2027-04-17 con **74 eventos** (antes 58): renta 1-oct, telefonías 6-oct, internet 12-oct, streaming 14-oct (realizados), despensa del 12, visitas a mercado realizadas (3, 6, 13, 17 y 20; tomadas de `gastos`, nunca el 10), pagos de crédito ya pagados de octubre y el reporte del 1. LEEME con el rango y la cifra reales.
- **LEEME 1.1:** fila `tickets`, filas de `calendario_hogar`, `inicio` y `ajustes` corregidas, regla de qué usa nombres de C1, nota de T-1 con 6 de 12 meses (B-1), «P3-12» (B-3), 51 comprobaciones.
- **Verificador:** 12 comprobaciones nuevas (de 39 a 51): tarjeta Próximo pago con nombres C1, pago vencido, `next_goal`, ticket confirmado (2), lectura OCR, por confirmar, notificaciones, datos semilla/parámetro, índice sin llaves colgantes, calendario (2). «P-3-12» → «P3-12».

## Verificación
```
python3 -I docs/design/tools/ds/generar_sinteticos.py docs/design/maquetas/datos-sinteticos.json   # determinista
cmp (segunda generación) datos-sinteticos.json                      -> REGENERACION IDENTICA
python3 -I docs/design/tools/ds/verificar_sinteticos.py             -> 51 OK, 0 FALLA, RESULTADO: OK
python3 -I docs/design/tools/ds/verificar_ds.py                     -> RESULTADO: OK · 445 tokens C2 ... 0 fallas
grep de cifras prohibidas (JSON y LEEME)                            -> 0 coincidencias; sin correos ni passwords
```

## Supuestos
- `sucursal: "Sucursal de ejemplo"`: no hay sucursales en las fuentes; no se inventó una dirección.
- `credit_id` conserva el id de línea de este archivo (`kueski`), no el entero de C1.
- `stale_prices` apagada y destinatarios (Axel; Ximena para el plan) son de ejemplo para el estado; K50 fija solo 1 día y 09:00.
- Cuenta de cada renglón de g13: todo «Conjunto 50 %» (el gasto es 50/50 y las partes deben cuadrar con g13).
- Los precios de los renglones son los de `productos` para esa tienda, aunque la fecha del precio sea otra.

## Bloqueos
Ninguno.

## Pendientes
- B-2 (`precio_unidad_base_cents` flotante sin `_fmt`) y B-4 ($9,000.00 de disponible B) siguen abiertos: eran de severidad baja y fuera del alcance pedido.
- Los archivos de `docs/design/entrada/` y `salida/` no se revisaron por quedar fuera del alcance; no citan las llaves renombradas (grep sin coincidencias).

## Consumo aproximado
~38 acciones de herramienta; ~110 KB leídos; ~245 KB escritos contando el JSON regenerado (205 KB, salida del generador, que ya pesaba 191 KB); el resto es unos 40 KB de código, LEEME y este reporte.
