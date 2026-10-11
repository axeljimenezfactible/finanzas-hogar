# revisor-U1-r3 · Reporte de agente (flujo P1, 11-oct-2026)

## Hecho
- Revisé por mi cuenta U1 1.1 en la ronda 3 de 3. Evalué contra `docs/criterios/U1.md` 1.1, las decisiones §K (K2 a K5, K8 a K10, K12, K13, K16, K21, K23 y K18) y los hallazgos R2-05 y R3-01.
- **Veredicto: APROBADO.** No hay hallazgos de severidad alta ni media; hay 7 de severidad baja (B1 a B5 siguen abiertos y B6 y B7 son nuevos). El hallazgo alto A1 de la ronda 2 (K18) queda cerrado.
- Volví a correr `verificar_u1.py` y comprobé lo mismo con dos scripts propios que usan otro método: pantalla de cada N-xxx según su sección, matriz §3.29 con sus columnas L y E, y `openapi.yaml` leído con PyYAML.
- Informe: `docs/estado/p1/REVISION-U1-r3.md`. Trae el veredicto, las verificaciones a a h, los 15 orígenes con su línea, la tabla de criterios, el recorrido de 29 capacidades, los hallazgos y el estado de los anteriores.

## Archivos
- `docs/estado/p1/REVISION-U1-r3.md` (nuevo)
- `docs/estado/p1/revisor-U1-r3.md` (nuevo, este reporte)
- Temporales en `datos/tmp/revisor-u1-r3/` (`rev.py`, `rev2.py`, `diff_plus.txt`, `diff_minus.txt`, `vu1_out.txt`): borrados al terminar.

## Verificación
- **`python3 -I docs/design/tools/verificar_u1.py`:** sale con 0 y dice «RESULTADO: TODO OK». Da 52 OK, 0 FALLA, 9 INFO y 2 PENDIENTE: N-316, N-318 y N-319, y los 4 campos de K18 en `Credit`, `CreditCreate` y `CreditUpdate`. Coincide con `correccion-U1-r2.md`.
- **`python3 -I datos/tmp/revisor-u1-r3/rev.py`:**
  - pantallas: 27 `P-xx` y 14 `P-26.xx` (nombres y orden exactos), así que T = 41;
  - necesidades: 203 N-xxx (138 L y 65 E), sin duplicadas ni inexistentes; lo declarado en las líneas 18 y 766 coincide;
  - matriz: 41 filas, 0 N-xxx fuera, 0 en la columna equivocada y 0 pantallas sin fila;
  - Local: N-003, N-301 y N-317;
  - contra `openapi.yaml` 1.0.0 (197 operaciones): sin operación quedan N-316, N-318 y N-319, y la única operación sin N-xxx es `getOpenApiDocument`;
  - los campos de K18 faltan en los 3 esquemas;
  - «Para C1 · Nuevas» menos las Local da lo mismo que las N-xxx sin operación;
  - cada K tiene sus líneas;
  - hay una sola definición de «lista de compras vigente» (66);
  - no hay medidas táctiles menores de 44 px, salvo el piso de 24 × 24 px;
  - PROVISIONAL: 66 y 1152 citan un bloqueo abierto; 421 y 1151 citan un bloqueo resuelto más K31;
  - §6.4 cubre las 27 pantallas y §2.4 declara qué hace cada una sin conexión;
  - F-1 da 3 toques;
  - las líneas cambiadas no traen cifras en $.
- **`python3 -I datos/tmp/revisor-u1-r3/rev2.py`:**
  - las 29 capacidades de U1-09 tienen pantalla y N-xxx;
  - §2.5 tiene 31 entidades;
  - la regla del domingo solo está en 23, 66 y 1152;
  - N-242 solo aparece en la fila de P-21;
  - `updatePerson` lleva [N-289, N-310].
- **`grep`:** `00002_puerta1.sql:112-115` («NULL = falta capturar»), `00001_esquema.sql:448` (`vat_on_interest_bp` NULL = «ya trae IVA»), RN-26, `S7.md:10` y `:46` (K31 en S7-08), y 0 apariciones de K18 en `C1.md`.

## Supuestos
- Ninguno cambia cifras, esquema, endpoints ni tokens. Hay tres criterios de evaluación propios:
  1. **U1-03 y U1-16 quedan en «Cumple en U1; C1 v1.1 pendiente».** En el orden del flujo, C1 v1.1 viene después de U1 1.1. Lo pendiente está en «Para C1» y el script lo reporta como PENDIENTE.
  2. **Una marca PROVISIONAL cuenta como justificada si cita un bloqueo resuelto por K31.** Es el caso de 421 y 1151: K31 manda «PROVISIONAL con la propuesta escrita» hasta la aprobación.
  3. **La representación del IVA no es hallazgo de U1.** Vive en S1 y RN-26; U1 sigue a K18.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Nada de esta revisión. Para otros agentes, ver «Solicitudes fuera de alcance» y §7.5 del informe (C1 v1.1).

## Solicitudes fuera de alcance
1. **`docs/estado/p1/U1-1.1.md`, sección «Para C1»:**
   - Cambio: sumar la adenda «Para C1 (adenda K18 · corrección r2)» de `docs/estado/p1/correccion-U1-r2.md`.
   - Motivo: B6. `C1.md` 1.1 no menciona K18 y `verificar_u1.py` solo lo marca como PENDIENTE.
2. **`docs/estado/p1/U1-1.1.md`, «Para C1», entrada de N-289:**
   - Cambio: agregar «quitar `updatePerson`; el tope por persona lo edita N-310 (K13)».
   - Motivo: B7.
3. **`docs/estado/bloqueos/INDICE.md`:**
   - Cambio: registrar `U1-lista-vigente-domingo` como abierto.
   - Motivo: lo pidió `U1-1.1.md` y sigue sin estar.
4. **`backend/db/migrations/00001_esquema.sql:448` y `00002_puerta1.sql` (S1), y `docs/contratos/reglas-negocio.md` RN-26:**
   - Cambio: que «si incluye IVA» pueda quedar en «falta capturar» en `bank_card` y `store_credit`.
   - Motivo: K18. Ya lo pidió la solicitud 2 de `correccion-U1-r2.md`; conviene resolverlo antes de C1 v1.1.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- **Acciones de herramienta:** 27 de 60.
- **Leído:** unos 195 KB. Nunca leí enteros U1, `openapi.yaml`, `C1-datos.md`, el diccionario ni las migraciones; solo `grep`, rangos y salidas de scripts.
- **Escrito:** unos 98 KB:
  - temporales, ya borrados: unos 72 KB (diff 51 KB, salida del script 6 KB y 2 scripts 15 KB);
  - `REVISION-U1-r3.md`: 19.7 KB;
  - este reporte: unos 5 KB.
