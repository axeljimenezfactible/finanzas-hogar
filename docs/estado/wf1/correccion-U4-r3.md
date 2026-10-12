# correccion-U4-r3 · WF-1 · 11-oct-2026

Corrige el único hallazgo medio de `docs/estado/wf1/REVISION-U4-r3.md` (R3-1, U4-05). No había altos. Los 5 bajos (R3-2 a R3-6) no se pidieron y no se tocaron. No se reabrió ninguna decisión de §A–§K.

## Hecho
- **R3-1 (media, U4-05), cerrado.** En `U4-ajustes-a.md`, P-26.02:
  - **«Edición en el lugar»** ahora nombra los campos del gasto fijo que guardan solo después de su vista previa en línea (§0.2): monto por vez, periodicidad, cada cuánto, primera fecha, fecha final, persona, reparto, y la categoría cuando la partida no tiene reparto propio (vista previa 3). De la despensa siguen el monto base, el día, la línea y el plazo (vista previa 4).
  - En los campos de lista (periodicidad, persona, reparto y categoría), elegir una opción muestra la vista previa sin guardar. `Enter` o salir del campo guarda, y `Esc` vuelve a la opción guardada.
  - **Vista previa 3** reescrita. Muestra «Antes → Después», calculado en el navegador sin llamar al servidor y en centavos enteros, con lo que la sección ya cargó (`listPlannedItems` `kind = fixed_expense`, `listCategories`, `listCostShareRules` con `parts`):
    - (a) Las próximas 3 fechas, con la regla de fechas del esquema (comentarios de `planned_item.frequency`, `interval_count`, `start_date` y `end_date` en 00001; días de pago de RN-06).
    - (b) «En {mes en curso}: {n} veces × {monto por vez} = {total de la partida}». En quincenal con monto distinto en la quincena del 30, se suman los dos montos. Si ni antes ni después cae una fecha en el mes en curso, se usa el primer mes que tenga alguna.
    - (c) «Gastos fijos de {mes}: {antes} → {después}». Cada fecha cuenta completa en su mes (`accrual_months` = 1). Si una partida tiene otro devengo, cuenta la porción de RN-02.
    - (d) Partes por persona en el mes. Cada fecha se reparte como un gasto (RN-03), con la regla de la partida; si no tiene, la de su categoría, y si tampoco, la de por defecto. Parte por porcentaje = ROUND_HALF_UP(monto × `pct_bp` ÷ 10,000). Parte fija = `fixed_cents`. La persona del resto (el titular) toma lo que falta y absorbe los centavos (K14). Ejemplo: $199.00 con 50/50 da Karen $99.50 y Hogar $99.50, igual que P-26.04. Si una parte fija pasa del monto, aparece un mensaje en lugar de las partes.
    - (e) Persona: «Se atribuye a {antes} → {después}», con los gastos fijos del mes de cada una; el total no cambia.
    - (f) Aviso: «Flujo y proyección (P-07) se recalcula». Los gastos fijos entran en los 12 meses con el interruptor de los meses 7 a 12 encendido o apagado (K21, RN-28).
    - Con datos incompletos, el bloque muestra lo que se puede calcular y dice «Completa la periodicidad y la primera fecha para ver las cifras».
  - **§0.2**: una frase de excepción para que la regla común («las listas guardan al elegir») no contradiga lo anterior. Una lista con vista previa en línea guarda con `Enter` o al salir del campo.
  - **SUPUESTOS 8**: cómo se lee «cada cuánto» mayor que 1 en una partida quincenal, y que el reparto se hace por fecha. Ninguno de los dos toca cifras canon. La vista previa sigue la misma regla de fechas que S2.
  - **Versión 1.3** en el encabezado de la parte a. La parte b sigue en 1.2, sin cambios.

## Archivos
- `docs/design/U4-ajustes-a.md` (modificado: encabezado, §0.2, P-26.02 «Edición en el lugar» y «Vista previa del efecto» punto 3, SUPUESTOS 8).
- `docs/estado/wf1/correccion-U4-r3.md` (este reporte).
- `docs/design/U4-ajustes-b.md`: sin cambios; ningún hallazgo pedido cae ahí.
- Scripts fuera del repositorio, en el scratchpad de la sesión: `u4r3/edita_a.py` (los reemplazos, cada uno con su comprobación) y `u4r3/v3.py` (la verificación completa: la de la ronda 2 más las pruebas de R3-1).

## Verificación
Comando: `python3 -I <scratchpad>/u4r3/v3.py /home/deepsec/finanzas-hogar`. Además, un conteo de celdas vacías en las dos partes. Salida resumida:
- **U4-01:** 14 secciones, P-26.01 a P-26.14 en orden, sin bloques faltantes. 321 filas de tabla (con encabezados), 0 con celdas vacías.
- **U4-02:** 210 `operationId` en los extractos, igual que en `openapi.yaml`. Se citan 100 en la parte a y 45 en la b, 0 inexistentes. Cubren las 73 de 73 operaciones de Ajustes. `listPlannedItems`, `listCategories` y `listCostShareRules` existen (extractos ajustes-a y ajustes-b).
- **U4-04 y U4-12:** las 15 claves `@param` de 00003 y las 19 claves sembradas de `app_setting` están citadas. `send_time`, `advance_days`, `weekday` y `day_of_month` llevan `S1-parametros-avisos`. «Mínimo de seguridad» fuera de SUPUESTOS: 0.
- **U4-05 (caso obligatorio A/B):** las series A y B de 14 valores aparecen literales y son iguales a `calculado_*`. Están los mínimos 8,094.50 y 5,000.00 y el disponible de hoy. A − B es constante. La despensa 1 trae $3,094.50, $800.50, $3,895.00 y $973.75. La cotización da 800.50.
- **U4-05 (R3-1):**
  - Los 8 campos están nombrados en «Edición en el lugar», y está la regla de las listas.
  - La vista previa 3 tiene los 16 elementos que exige el script: (a) a (f), la fórmula ROUND_HALF_UP con `pct_bp`, K14, P-07, K21/RN-28, centavos enteros, las 3 operaciones y «sin llamar al servidor». La frase vieja («solo fechas») aparece 0 veces.
  - Reparto: 19,900 × 5,000 ÷ 10,000 = 9,950 → Karen $99.50 y Hogar $99.50. `parametros.papel` = 199.0.
  - Regla de fechas (a) probada contra el oráculo:
    - Mensual cada 2 desde el 10-oct-2026 da 10-oct-2026, 10-dic-2026, 10-feb-2027 y 10-abr-2027, igual a las compras de `despensas[*]`.
    - Los días de pago desde el 10-oct-2026 + 15 dan `despensas[0].pagos` (30-oct, 15-nov, 30-nov y 15-dic-2026).
    - Febrero paga el 28 (`despensas[2].pagos[0]`).
    - Recorte a fin de mes: mensual desde el 31-ene-2027 da 31-ene, 28-feb y 31-mar-2027.
    - «Una vez» con fin igual al inicio da solo esa fecha.
- **Regresión de la ronda 2 (U4-05 y U4-13):** IMC 23.7; meta 2,555.5625 → 2,556 y 2,056; S4-07 1,799; RN-18 50 % (50.0 %). P-26.05 conserva las vistas previas 1 a 10, y R2-1, R2-2 y R2-3 siguen cerrados.
- **U4-09 y U4-11:** 0 hex, 0 `rgb(` y 0 fechas ISO en las dos partes. 21 tokens citados, 0 inexistentes. 0 bloqueos citados que no existan. PROVISIONAL en la parte a: 9 líneas, todas con cita. En la b: 27 líneas, y las 5 sin cita en la misma línea (26, 104, 209, 361 y 362) son las de rondas anteriores, bajas y fuera de esta ronda.
- **U4-06:** no se agregaron atajos ni comandos «/». `Enter` y `Esc` ya estaban en §0.2.
- **No aplican** pruebas Go, `compara.py` ni la compilación de tokens: U4 es un documento de diseño sin código. No creé temporales en `datos/tmp`; lo que hay ahí es de otros agentes y no lo toqué.

## Supuestos
- Quincenal con «cada cuánto» mayor que 1: se cuenta uno de cada {cada cuánto} días de pago, desde el primero que no es anterior a la primera fecha. Es una lectura de los comentarios de `planned_item`. No toca cifras canon (el oráculo no trae gastos fijos) y la vista previa queda atada a la regla del motor S2. Va como SUPUESTOS 8 en la parte a.
- El reparto por fecha aplica RN-03 a cada gasto que generará la partida («cada gasto dice de quién es»). Así, los centavos de la vista previa son los mismos que se congelarán en `expense_share`.

## Bloqueos
- Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md` (no lo toca esta corrección).

## Pendientes
- Ninguno de R3-1. Los bajos R3-2 a R3-6 de `REVISION-U4-r3.md` no se pidieron en esta ronda y siguen abiertos: los decide el integrador o el revisor de la ronda siguiente.

## Solicitudes fuera de alcance
- Ninguna.

## Paquetes instalados
- Ninguno.

## Propuestas para Axel
- Ninguna. R3-1 no pide valores de K31.

## Consumo aproximado
- Acciones de herramienta: 29.
- Leído: ≈ 125 KB (protocolo, §K y §B de decisiones, criterios U4, REVISION-U4-r3, partes de U4-a y U4-b, RN-01/02/03/05/06/22/28, `CREATE TABLE planned_item`, fragmentos de los extractos, el script de verificación de r2 y las salidas).
- Escrito: ≈ 32 KB (el cambio neto en U4-a es de ≈ 5 KB, aunque el script reescribe el archivo completo; más los scripts `edita_a.py` y `v3.py` y este reporte).
