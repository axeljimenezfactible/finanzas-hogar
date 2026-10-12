# correccion-datos-r5 · 12-oct-2026

## Hecho
- `docs/design/tools/ds/generar_sinteticos.py`: el flujo por quincena sigue S2-interfaz 1.3 §3.11. `variable_estimado_cents` y `ahorro_planificado_cents` van en `null` en las 24 quincenas (decisión: `null`, que es lo que dice §3.11 campo por campo; el supuesto de $5,800.00 queda solo como texto informativo en `proyeccion.supuestos` y `presupuesto.metodo_variable`). `saldo_cents` = ingresos − pagos de crédito − fijos y venideros. `SALDO_INI` = 0. Se quitó la corrección aparte de la quincena en curso (ya no hace falta).
- `docs/design/maquetas/datos-sinteticos.json`: regenerado. Único cambio de llave-valor fuera de la serie: `proyeccion.saldo_inicial_cents` 1200000 → 0 (la llave se conserva por la regla «diff de llaves idéntico»). Los otros 81 cambios están en `proyeccion.flujo_por_quincena` (los dos campos en null, `saldo_cents` de las filas que restaban, `saldo_acumulado_cents` de toda la serie).
- `docs/design/tools/ds/verificar_sinteticos.py`: la comprobación de la serie encadena desde 0 sin restar variable ni ahorro; comprobaciones nuevas: sin saldo inicial, null en las 24 filas, y quincena en curso con acumulado 182778 y 1.ª quincena −421772.
- `docs/design/maquetas/LEEME.md` 1.4: nota r5 (regla, valores nuevos, qué no cambió); la nota r4 de la quincena en curso queda marcada como reemplazada.

## Valores nuevos (centavos, saldo / acumulado)
| Quincena | saldo | acumulado |
|---|---|---|
| 2026-10-q1 | −421772 | −421772 (−$4,217.72) |
| 2026-10-q2 (en curso) | 604550 | 182778 ($6,045.50 / $1,827.78) |
| 2026-11-q1 | −106434 | 76344 |
| 2026-11-q2 | 998151 | 1074495 |
| 2026-12-q1 | 282659 | 1357154 |

La serie reproduce los $1,827.78 del revisor: −421772 + 604550 = 182778. La 1.ª quincena era −521772 porque restaba $1,000.00 de ahorro; ahora 1561000 − 578617 − 1404155 = −421772.

## Verificación
- `python3 -I verificar_sinteticos.py`: RESULTADO: OK (63 comprobaciones, 3 nuevas o reescritas para la regla).
- Grep P3-12 (con/sin coma) sobre el JSON: 0 coincidencias (comprobación 1 del verificador).
- Diff de llaves entre el JSON anterior y el nuevo: idéntico (mismas llaves y misma longitud de listas); 82 diferencias de valor, todas en `proyeccion`.
- `inicio.tarjetas.flujo_quincenal.closing_balance_cents` sigue en 604550 = saldo de 2026-10-q2.

## Supuestos
- `saldo_inicial_cents` queda en 0 en vez de borrarse, por la regla de llaves idénticas. Las pantallas no deben mostrarlo.
- Las etiquetas `base` («fijos + variable estimado») conservan su texto: P-07 las traduce por clave y cambiarlas rompe el mapa. Ya no describen una resta.
- `proyeccion.meses` (12 meses devengado/efectivo, §3.12) sigue usando el estimado variable en sus fórmulas; es otra pieza del contrato y la tarea no la incluye.

## Bloqueos
- Ninguno.

## Pendientes (para otros agentes)
- Autores de Inicio (detalle-m l. 103 y detalle-e l. 142): «Saldo acumulado $1,827.78».
- Autor de P-07 (`presupuesto-p07-flujo-m/e.dc.html`): tiene la serie copiada a mano en `const F` con `saldo_inicial_cents` 1200000 y los valores viejos; hay que reemplazarla por la del JSON nuevo y quitar cualquier «saldo inicial» visible. Con `variable_estimado_cents`/`ahorro_planificado_cents` en null, sus columnas ya pintan «—» (el código maneja null; hoy además esconde las filas con `mes_n > 6`, que ahora también son null).

## Consumo aproximado
~15 acciones, ~60 KB leídos, ~8 KB escritos.
