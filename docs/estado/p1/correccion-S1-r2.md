# Corrección S1 · ronda 2 · agente correccion-S1-r2 (flujo P1)

Fecha: 2026-10-11. Encargo: corregir el hallazgo **R2-01** (alta, S1-03, verificación j) de `docs/estado/p1/REVISION-S1-r2.md`. La revisión r2 no tiene hallazgos medios. Los 8 bajos (R2-02 a R2-09) no se tocaron. Ninguna decisión de §A–§K se reabrió.

## Hecho
- **RN-06** (`docs/contratos/reglas-negocio.md` L124, antes L122): la frase de semana y quincena lleva el rótulo «Ejemplo ilustrativo (sin dato en el oráculo) de los periodos (K8)». Sus fechas 5-oct-2026, 11-oct-2026 y 31-oct-2026 no están en el oráculo. El verificador extendido encontró en la misma línea el «15-oct» del cálculo del primer pago (2026-10-15 tampoco está en el oráculo). Por eso el resultado queda citado (`model.json: despensas[0].pagos[0..3].fecha`) y la razón del primer pago va en una frase aparte con el rótulo «Ejemplo ilustrativo (sin dato en el oráculo) del primer pago». Ninguna fecha cambia: al «15-oct» solo se le agregó el año.
- **RN-11** (L177, antes L175): el encabezado del ejemplo dice ahora «las fechas de compras y pagos y las demás cifras son del oráculo, con su ruta o con su cuenta, salvo lo rotulado como ilustrativo». El escenario del 16-dic-2026 lleva rótulo: «como ejemplo ilustrativo (sin dato en el oráculo), moverla al 16-dic-2026, el día siguiente a ese pago, la cumple». La fecha de la restricción «deuda en cero antes del 1-may-2027» es un valor de restricción (K16) y no está en el oráculo. Lleva «como valor ilustrativo (sin dato en el oráculo)» en su propia frase, porque el rótulo de una frase solo cubre esa frase.
- Cabecera de `reglas-negocio.md`: se agregó el párrafo «Correcciones de la revisión S1-r2» con estos cambios. Ninguna cifra ni fecha cambia.
- **Verificador** (`backend/db/tools/s1_verifica.py`, bloque S1-03, y cabecera): ahora revisa también las fechas de los ejemplos, en tres formatos:
  - ISO (`2026-10-10`);
  - dd-mmm-aaaa (`10-oct-2026`);
  - dd-mmm sin año (`15-oct`). El año sale de la fecha completa o del mes-año (`oct-2026`) más cercano dentro de la misma frase. Se busca primero a la derecha (un año menos si el mes de la fecha es mayor). Si no hay, a la izquierda o en el último del ejemplo (un año más si el mes de la fecha es menor). Si el ejemplo no da ningún año, la fecha queda como «????» y cuenta como fuera del oráculo.

  Las fechas del oráculo son las que aparecen en valores, llaves o textos de `model.json` y `data.json` (29 distintas). S1-03 falla si una fecha fuera del oráculo aparece en una frase sin rótulo, o si una fecha es imposible (31-nov) aunque lleve rótulo. Una autoprueba de 5 casos comprueba que la marca funciona: las frases que señaló la revisión (RN-06; RN-11 con el 16-dic y el 1-may), el 15-oct, un cruce de diciembre a enero y una fecha imposible. Si la autoprueba falla, S1-03 falla.
- Comprobé que, con la versión de r2 de `reglas-negocio.md`, el verificador extendido da FALLA (Verificación 1).

## Archivos
- `docs/contratos/reglas-negocio.md`: modificado (cabecera, RN-06 L124, RN-11 L177).
- `backend/db/tools/s1_verifica.py`: modificado (cabecera y bloque S1-03).
- `docs/estado/p1/correccion-S1-r2.md`: este reporte.
- Sin cambios: `backend/db/migrations/00001_esquema.sql`, `backend/db/migrations/00002_puerta1.sql`, `docs/contratos/diccionario-datos.md`, `backend/db/tools/gen_dict.py`.

## Verificación
1. **Verificador extendido contra la versión de r2 de las reglas.** Comando: `GOOSE=<scratchpad>/gobin/goose python3 -I backend/db/tools/s1_verifica.py`.
   - Resultado: rc 1, «RESUMEN FALLAS: ['S1-03']». Las otras 23 comprobaciones dan OK.
   - S1-03: «ejemplos con fechas que no están en model.json ni en data.json y sin rótulo=2 [('RN-06', ['2026-10-15', '2026-10-05', '2026-10-11', '2026-10-31']), ('RN-11', ['2026-12-16', '2027-05-01'])]». La autoprueba pasa (5 casos).
   - En una primera corrida, el verificador todavía no leía el mes-año y marcaba RN-02 («????-10-30», «????-11-15», «????-11-30»), porque su tabla da el año como «oct-2026». Agregué el mes-año como fuente del año y esas fechas quedan en el oráculo (2026-10-30, 2026-11-15, 2026-11-30). RN-02 no se tocó.
2. **Después de la corrección**, mismo comando: rc 0, «RESUMEN TODO OK» (23 comprobaciones OK).
   - S1-03: 28 reglas; 0 sin ejemplo; 0 citas que no coinciden; 0 cifras fuera del oráculo sin rótulo.
   - Fechas: 63 en los ejemplos y 29 distintas en el oráculo; 0 fuera del oráculo sin rótulo; 0 imposibles; 14 ilustrativas fuera del oráculo. La autoprueba pasa (5 casos).
   - S1-12: 164 valores, 0 diferencias. KUESKI: 389500 y 239528. NEG: 46/46 rechazadas. NEG12: 26/26 rechazadas y 12/12 positivas aceptadas.
   - DOWN: igual al esquema 1.1 (232 objetos). GOOSE: códigos [0, 0, 0], tablas [55, 0, 55], versión final 2.
3. **Diccionario.** `python3 -I backend/db/tools/gen_dict.py datos/tmp/dic-regen.md` da «tablas 55 núcleo 45 complementarias 10 índices 152 disparadores 71» y «columnas sin descripción: []». `cmp` contra `docs/contratos/diccionario-datos.md`: idéntico (174,605 bytes).
4. **goose v3.28.0 a mano** en `datos/tmp/g.db`:
   - up: rc 0 (00001 y 00002, versión 2);
   - down-to 0: rc 0 (00002 y 00001, versión 0);
   - up: rc 0 (versión 2);
   - status: las dos migraciones aplicadas. `integrity_check` ok, `foreign_key_check` 0 filas, 55 tablas.
   - Un primer intento de down-to dio rc 1 porque zsh pasó «down-to 0» como un solo argumento. Con los argumentos separados, pasa.
5. **Fechas de los textos rotulados**, comprobadas con `datetime`:
   - el 10-oct-2026 es sábado, y su semana va del lunes 2026-10-05 al domingo 2026-10-11;
   - el 31-oct-2026 cae en la quincena 2 (días 16 a 31);
   - el 15-oct-2026 queda a 5 días de la compra y el 30-oct-2026 a 20 días;
   - el 16-dic-2026 es el día siguiente al 15-dic-2026.
6. **Conteo aparte de ilustrativos** (S1-03):
   - 12 ejemplos completos: RN-12, RN-13, RN-14, RN-16, RN-17, RN-18, RN-19, RN-22, RN-23, RN-24, RN-26 y RN-27;
   - 15 frases rotuladas en 8 ejemplos con dato: RN-04 (1), RN-06 (2), RN-07 (3), RN-08 (1), RN-10 (1), RN-11 (5), RN-15 (1) y RN-28 (1);
   - 31 cifras y 14 fechas ilustrativas fuera del oráculo;
   - 10 rótulos con la leyenda literal «ilustrativo (sin dato en el oráculo)».

## Supuestos
- El «15-oct» de RN-06 entra en R2-01: es la misma regla, la misma línea (L122 en r2) y el mismo defecto (fecha fuera del oráculo sin rótulo). Lo encontró el verificador al revisar también las fechas sin año. No cambia cifras, esquema, endpoints ni tokens.
- La fecha de la restricción «deuda en cero» (1-may-2027) ya estaba cubierta por el encabezado de RN-11 («los valores de las restricciones son ilustrativos»). Se rotuló también en su frase para que el verificador, que revisa frase por frase, la acepte. No cambia cifras.

## Bloqueos
Ninguno nuevo.

## Pendientes
Ninguno de R2-01. Los bajos R2-02 a R2-09 quedan como los dejó la revisión, porque están fuera del encargo.

## Solicitudes fuera de alcance
1. `docs/estado/p1/S1-esquema-1.2.md`, línea S1-03: enlazar este reporte para el conteo vigente de ilustrativos y fechas. Complementa R2-04, que ya pide enlazar `correccion-S1-r1.md`. Motivo: el reporte S1 todavía trae cifras de antes de las correcciones.

## Paquetes instalados
Ninguno.

## Consumo aproximado
- 22 acciones de herramienta.
- ~80 KB leídos.
- ~18 KB escritos: parches de ~6.5 KB en el verificador y de ~3.5 KB en las reglas, más este reporte de ~7.6 KB.
