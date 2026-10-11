# Corrección S1 · esquema 1.2 · ronda 1 (flujo P1)

Agente `correccion-S1-r1` · 11-oct-2026 · Corrige los hallazgos de severidad alta y media de `docs/estado/p1/REVISION-S1-r1.md` (H-01, H-02 y H-03). No reabre §A–§K y **no cambia ninguna cifra del oráculo**. Los hallazgos bajos (H-04 a H-08) no estaban en esta asignación; ver Pendientes.

## Hecho
1. **H-01 (alta, S1-03): ejemplos sin rótulo con cifras que no están en el oráculo.**
   - **RN-12** y **RN-27** empiezan con «**Ejemplo ilustrativo (sin dato en el oráculo).**» y una frase dice qué parte sí es del oráculo. Conservan sus citas: `data.json: fecha = 2026-10-10` y `model.json: periodos[0].devengado.despensa_total = 1,947.50`. Las cuentas no cambian: pesos 1 y 0.5 con la lectura (a), 60/61 y 30/61 con la (b); 1,947.50 / 31 = 62.8225… → 62.82; 1,947.50 × 7 / 31 = 439.758… → 439.76; 62.82 × 7 = 439.74.
   - **RN-11**: el encabezado ya no dice que todas las cifras son del oráculo. Los valores de las restricciones se rotulan como ilustrativos porque los captura Axel (K16), y el colchón de 5,500.00 lleva «ilustrativo (sin dato en el oráculo)». 1,814.50 cita su cuenta: `data.json: financiamiento.dos_quincenas.total / 2 = 1,814.50`. La división 3,629.00 / 2 es exacta y coincide con `axel.por_quincena` 1,115.85 + `karen.por_quincena` 698.65.
   - **Verificador** (`backend/db/tools/s1_verifica.py`, bloque S1-03). (a) Junta las cifras del oráculo: las hojas numéricas de `model.json` y `data.json` y las cifras escritas dentro de sus textos, por ejemplo `sustitutos[0].nota`. (b) En cada ejemplo sin rótulo marca toda cifra de dos decimales que no esté en el oráculo. Solo la deja pasar si es el resultado de una cuenta citada «ruta / n = valor», que el verificador comprueba con fracciones y ROUND_HALF_UP. (c) El rótulo de una frase cubre solo esa frase (corta en «.» y «;»). Un ejemplo que empieza con «**Ejemplo ilustrativo» es ilustrativo completo. (d) Cuenta aparte los ilustrativos, las cifras ilustrativas y los rótulos con la leyenda literal. Falla con una sola cifra fuera del oráculo sin rótulo.
   - **Por el verificador estricto también se tocaron RN-07 y RN-08.** Tenían cifras de contraste (lo que saldría con otro redondeo) que no están en el oráculo y no tenían rótulo: RN-07 1,258.67, 1,573.35 y 1,573.37; RN-08 258.69 y 1,293.43. A esas frases solo se les antepuso «ilustrativo (sin dato en el oráculo)»; no cambió ninguna cifra ni otra palabra. La frase de RN-07 «la serie real sería…» es la segunda mitad de H-08.
   - **Prueba del verificador:** sobre las reglas de antes (respaldo) da **FALLA** con 5 ejemplos: RN-07, RN-08, RN-11 (5,500.00 y 1,814.50), RN-12 (0.98 y 0.49) y RN-27 (62.82, 439.76 y 439.74). Sobre las reglas de ahora da **OK**.
2. **H-02 (media, S1-09 y K8): `shopping_list.start_date`.**
   - El comentario de 00001 L1111 dice ahora «Lunes que abre la semana del plan (K8: semana de lunes a domingo)», con la regla «fecha válida y en lunes».
   - La base lo exige con la opción que propuso el revisor: en la 00002, sección K8, van los disparadores `trg_shopping_list_monday_ins` (BEFORE INSERT) y `trg_shopping_list_monday_upd` (BEFORE UPDATE OF start_date). Los dos rechazan con «shopping_list: la semana de la lista empieza en lunes (K8)».
   - El Down los borra en orden inverso, antes de `DROP TABLE telegram_chat_candidate`.
   - NEG12 suma 2 pruebas negativas (alta en martes, cambio a miércoles) y 2 positivas (alta en lunes, cambio a otro lunes).
   - Es compatible con C1 v1.1, que pide la semana de la lista como `period` ISO (`2026-W42`, de lunes a domingo; `openapi.yaml` L4393).
   - La línea «Esquema» de RN-15 nombra los disparadores. El diccionario se regeneró.
3. **H-03 (media, S1-04): consumidores en U1.**
   - `@used_new` del borrador de `attachment` (00002 L31): «U1 P-05 (N-071, N-075 y N-076)». Fuente: U1 L28 y L484-485; P-05 en L774.
   - `@used` de `telegram_chat_candidate`: «U1 P-26.09: N-299 (chats detectados) y N-316 (Confirmar chat detectado)». Fuente: U1 L745 y L803.
   - El diccionario se regeneró.
4. **Notas de versión** (las versiones siguen en 1.2 porque es el mismo flujo, sin commit):
   - reglas: «Correcciones de la revisión S1-r1» (L13);
   - diccionario: la misma línea en el encabezado, que es texto fijo que `gen_dict.py` conserva;
   - cabeceras de la 00002 y de `s1_verifica.py`.

## Archivos
- `backend/db/migrations/00001_esquema.sql`: solo el comentario de L1111; el DDL no cambia.
- `backend/db/migrations/00002_puerta1.sql`: cabecera (L9), L31 (`@used_new`), `@used` de `telegram_chat_candidate`, sección K8 con 2 disparadores (L69-88) y 2 `DROP TRIGGER` en el Down (L164-165).
- `backend/db/tools/s1_verifica.py`: cabecera, bloque S1-03 y 2 negativas más 2 positivas en NEG12/POS12.
- `docs/contratos/reglas-negocio.md`: nota en L13; RN-07 L134, RN-08 L144, RN-11 L175, RN-12 L185, RN-15 L207 y RN-27 L272.
- `docs/contratos/diccionario-datos.md`: regenerado. Nota en L9; «Resumen» con 71 disparadores (17 de reglas); L1107, L1774, L1795-1796 y L2033.
- `docs/estado/p1/correccion-S1-r1.md`: este reporte.

## Verificación
1. `GOOSE=<scratchpad>/gobin/goose python3 -I backend/db/tools/s1_verifica.py` → rc 0, **RESUMEN TODO OK**: 23 de 23 OK (S1-01 a S1-17, KUESKI, K12, NEG, NEG12, DOWN y GOOSE).
   - KUESKI: 389500 y 239528.
   - NEG12: 26 negativas rechazadas y 12 positivas aceptadas.
   - DOWN: `sqlite_master` igual al de 1.1 (232 objetos, 0 diferencias) y `app_setting` igual.
   - GOOSE: up / down-to 0 / up con rc [0, 0, 0], tablas [55, 0, 55] y versión 2.
2. **S1-03:** 28 reglas, 0 sin ejemplo, 0 citas que no coinciden y **0 ejemplos con cifras de dos decimales fuera de `model.json` y `data.json` sin rótulo**. 1 cuenta citada coincide.
   - **Ilustrativos, contados aparte:** 12 ejemplos completos (RN-12, RN-13, RN-14, RN-16, RN-17, RN-18, RN-19, RN-22, RN-23, RN-24, RN-26 y RN-27).
   - 11 frases rotuladas en 7 ejemplos con dato: RN-04 (1), RN-07 (3), RN-08 (1), RN-10 (1), RN-11 (3), RN-15 (1) y RN-28 (1).
   - 31 cifras ilustrativas no están en el oráculo.
   - 6 rótulos usan la leyenda literal «ilustrativo (sin dato en el oráculo)».
3. **El verificador nuevo contra las reglas de antes** (respaldo, con el mismo bloque S1-03): FALLA con 5 ejemplos (RN-07, RN-08, RN-11, RN-12 y RN-27). Contra las de ahora: OK.
4. `python3 -I backend/db/tools/gen_dict.py` → «tablas 55 · núcleo 45 · complementarias 10 · índices 152 · disparadores 71 · columnas sin descripción: []».
   - Una segunda corrida a otro archivo, comparada con `cmp`, sale idéntica (174,605 bytes).
   - El diff contra el diccionario anterior trae solo los cambios esperados: nota, «Resumen», L1107, L1774, los 2 disparadores y L2033.
5. goose v3.28.0 a mano: up (v2, 55 tablas, 71 disparadores, 2 de lunes) / version / down-to 0 (v0, 0 tablas) / version / up (v2) / version. rc 0 en los 6 pasos.
6. `diff <(sed 's/--.*//' respaldo/00001) <(sed 's/--.*//' 00001)` sale vacío: el DDL de 00001 no cambió.
7. sqlite3 en memoria:
   - con 1.1, una lista que empieza en martes se acepta;
   - con 1.2, la base rechaza el martes 2026-10-13 y el domingo 2026-10-11 con «shopping_list: la semana de la lista empieza en lunes (K8)»;
   - con 1.2, el lunes 2026-10-12 entra.

**Línea S1-03, regla por regla** (salida de `s1_verifica.py`; «cuentas citadas» son las de la forma «ruta / n = valor»):

| Regla | Ejemplo | Tipo | Citas con ruta que coinciden | Cuentas citadas | Cifras fuera del oráculo sin rótulo |
|---|---|---|---|---|---|
| RN-01 | sí | con dato | 1 de 1 | 0 | 0 |
| RN-02 | sí | con dato | 1 de 1 | 0 | 0 |
| RN-03 | sí | con dato | 5 de 5 | 0 | 0 |
| RN-04 | sí | con dato; 1 frases ilustrativas | 1 de 1 | 0 | 0 |
| RN-05 | sí | con dato | 2 de 2 | 0 | 0 |
| RN-06 | sí | con dato | 1 de 1 | 0 | 0 |
| RN-07 | sí | con dato; 3 frases ilustrativas | 1 de 1 | 0 | 0 |
| RN-08 | sí | con dato; 1 frases ilustrativas | 2 de 2 | 0 | 0 |
| RN-09 | sí | con dato | 2 de 2 | 0 | 0 |
| RN-10 | sí | con dato; 1 frases ilustrativas | 1 de 1 | 0 | 0 |
| RN-11 | sí | con dato; 3 frases ilustrativas | 4 de 4 | 1 | 0 |
| RN-12 | sí | ilustrativo completo | 1 de 1 | 0 | 0 |
| RN-13 | sí | ilustrativo completo | 1 de 1 | 0 | 0 |
| RN-14 | sí | ilustrativo completo | 0 de 0 | 0 | 0 |
| RN-15 | sí | con dato; 1 frases ilustrativas | 1 de 1 | 0 | 0 |
| RN-16 | sí | ilustrativo completo | 0 de 0 | 0 | 0 |
| RN-17 | sí | ilustrativo completo | 0 de 0 | 0 | 0 |
| RN-18 | sí | ilustrativo completo | 1 de 1 | 0 | 0 |
| RN-19 | sí | ilustrativo completo | 0 de 0 | 0 | 0 |
| RN-20 | sí | con dato | 1 de 1 | 0 | 0 |
| RN-21 | sí | con dato | 6 de 6 | 0 | 0 |
| RN-22 | sí | ilustrativo completo | 1 de 1 | 0 | 0 |
| RN-23 | sí | ilustrativo completo | 2 de 2 | 0 | 0 |
| RN-24 | sí | ilustrativo completo | 0 de 0 | 0 | 0 |
| RN-25 | sí | con dato | 2 de 2 | 0 | 0 |
| RN-26 | sí | ilustrativo completo | 0 de 0 | 0 | 0 |
| RN-27 | sí | ilustrativo completo | 1 de 1 | 0 | 0 |
| RN-28 | sí | con dato; 1 frases ilustrativas | 3 de 3 | 0 | 0 |

## Supuestos
- Ninguno que cambie cifras, esquema, endpoints o tokens. Los disparadores aplican K8 (FIRME) a la columna que el propio esquema 1.1 ataba al bloqueo `U1-periodos-semana-quincena`, hoy RESUELTO, y son la opción que propuso el revisor en H-02.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Estos hallazgos bajos de `REVISION-S1-r1.md` quedan abiertos porque no estaban en esta asignación (solo alta y media):
  - **H-04:** que la base exija todos o ninguno de los datos de K18 (`00002_puerta1.sql`).
  - **H-05:** la recomendación de la lectura (b) en RN-12.
  - **H-06:** `rate_base` contra `rate_basis`.
  - **H-07:** la fecha de la cabecera de la 00002.
  - **H-08, primera mitad:** los rótulos anteriores no usan la leyenda literal. Son 10 ejemplos completos (RN-13, RN-14, RN-16, RN-17, RN-18, RN-19, RN-22, RN-23, RN-24 y RN-26) y las frases de RN-04, RN-07 («Ilustrativo para una parte externa…»), RN-10, RN-11 («tercer escenario ilustrativo»), RN-15 y RN-28. El verificador ya los cuenta: 6 rótulos literales.
  - Rutas: `docs/contratos/reglas-negocio.md` y `backend/db/migrations/00002_puerta1.sql`.

## Solicitudes fuera de alcance
1. `docs/estado/p1/S1-esquema-1.2.md`, línea S1-03.
   - Cambio: poner la salida nueva del verificador con los ilustrativos contados aparte (12 completos, 11 frases en 7 ejemplos, 31 cifras ilustrativas y 0 cifras fuera del oráculo sin rótulo), o enlazar este reporte.
   - Motivo: H-01 pide ese conteo en el reporte, y ese archivo es del agente S1-esquema-1.2.
2. C1 y S6 (informativa, sin cambio en C1).
   - La API debe guardar `shopping_list.start_date` como el lunes de la semana ISO de `period`; la base ya rechaza cualquier otro día.
   - C1 v1.1 ya usa la semana ISO (`2026-W42`, `openapi.yaml` L4393).

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Acciones de herramienta: 33.
- Leído: unos 100 KB, siempre por partes (grep y sed).
- Escrito de contenido propio: unos 23 KB (bloque S1-03 4.4 KB, script de parches 8.6 KB y este reporte unos 10 KB).
- Las ediciones se aplicaron con `str.replace` (neto +6 KB en los archivos).
- El diccionario lo regeneró `gen_dict.py` (174,605 bytes, 7 bloques de cambio).
- Los temporales de `datos/tmp/corr-s1r1/` se borraron al terminar.
