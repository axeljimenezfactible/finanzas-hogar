# integrador-p1 · Reporte final (flujo P1, 11-oct-2026)

## Hecho
- `Makefile`: `verify-s1` corre `backend/db/tools/s1_verifica.py` (se quitó el `define` que extraía el Anexo B); objetivos nuevos `verify-u1` (`python3 -I docs/design/tools/verificar_u1.py`) y `verify-criterios` (`docs/criterios/tools/verif_final.sh`, que siempre sale con 0, así que el Makefile exige 13 cadenas esperadas y dos «0»; variable `CRITERIOS_SH` para la prueba negativa); `verify-all` = verify-wf0 + verify-c1 + verify-s1 + verify-u1 + verify-criterios; ayuda, cabecera y requisito Go 1.27.
- Confirmado: `verify-c1` ya cubre los extractos (`build.py --check` compara 21 archivos, 18 de ellos en `docs/contratos/extractos/`). `verify-sql` aplica 00001 y luego 00002 (wildcard ordenado), y el Down en orden inverso deja 0 objetos.
- `make verify-all` completo con goose v3.28.0: pasa dos veces (antes y después de editar los documentos).
- `README.md`: estado al 11-oct-2026 (puerta 1 aplicada, C1 sin aprobar, sigue WF-1), árbol con `backend/db/tools`, `docs/contratos/extractos`, `docs/criterios/tools`, `docs/design/tools` y `docs/estado/p1`, tabla de versiones vigentes y sección de verificación actualizada.
- `docs/criterios/INDICE.md`: conteo propio por script en los 18 archivos = 257 (197 / 60 / 0); coincide con la tabla de §2, con `traza.py` y con los reportes (la diferencia con el 256 de `criterios-UC-1.1.md` es U5-13, de `correccion-criterios-r1`). Se corrigió «248» en §1 y la fecha de la línea 3, y se agregó §10.
- `docs/estado/bloqueos/INDICE.md`: sección «Bloqueos nuevos de P1» con los 8 abiertos, su fase, la pregunta, lo que queda PROVISIONAL y la fase que lo usa.
- Conteo de PROVISIONAL y bloqueos citados, en `docs/estado/p1/INTEGRACION.md` §4.
- `docs/estado/p1/INTEGRACION.md`: revisiones, hecho por artefacto, verificaciones, PROVISIONAL, hallazgos registrados, bloqueos nuevos, pendientes para WF-1, paquetes y consumo por agente.
- Commit único del flujo P1 (hash en la salida estructurada).

## Archivos
- `/home/deepsec/finanzas-hogar/Makefile`
- `/home/deepsec/finanzas-hogar/README.md`
- `/home/deepsec/finanzas-hogar/docs/criterios/INDICE.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/INDICE.md`
- `/home/deepsec/finanzas-hogar/docs/estado/p1/INTEGRACION.md` (nuevo)
- `/home/deepsec/finanzas-hogar/docs/estado/p1/integrador-p1.md` (nuevo)

## Verificación
- `make -n verify-all GOOSE=/x` → sin errores de sintaxis.
- `make verify-all GOOSE=<scratchpad>/gobin/goose` → código 0, unos 10 s. Detalles: SQL con 00001 y 00002, 55 tablas STRICT, 152 índices y 71 disparadores; Down en orden inverso deja 0 objetos; goose up, down-to 0 y up sin errores. Go 1.27.2 con build, vet y gofmt limpios. Colorimetría idéntica, 202 pares y 0 `cumple=false`. OpenAPI OK. C1: 21 archivos idénticos (extractos incluidos), `verify_c1.py` con 24 OK y «TODO OK». S1: «RESUMEN TODO OK». U1: 53 OK, 0 PENDIENTE y «TODO OK». Criterios: 257 criterios, 178 cifras canon sin diferencias y 12 de 12 cifras del calendario.
- Segunda corrida de `make verify-all` después de editar los documentos → código 0, mismas cifras.
- Prueba negativa: `make verify-criterios CRITERIOS_SH=<falso>` → código 2, con dos FALLA; con el script real → código 0.
- Conteo propio (`python3 -I conteo.py docs/criterios`) → 18 archivos, 257 (197 / 60 / 0), sin IDs repetidos, todos en la versión 1.1.
- `grep -c PROVISIONAL` → 00001: 0, 00002: 0, diccionario: 1, reglas: 11, U1: 7, openapi.yaml: 6, C1-datos.md: 18.
- `make clean` → «Limpio: sin datos/tmp, sin backend/bin y sin __pycache__».

## Supuestos
- `CRITERIOS_SH` es solo una variable del Makefile para la prueba negativa; no cambia lo que se verifica por omisión.
- Lo que el Makefile exige de `verif_final.sh` es la salida aprobada por `revisor-criterios-r3`; «sin fuente: 5» queda como información (hallazgo bajo B-04).

## Bloqueos
Ninguno nuevo de este agente.

## Pendientes
Están en `docs/estado/p1/INTEGRACION.md` §7. Los principales:
1. Decisión del orquestador sobre C1, que no quedó aprobado: H-r3-01 a H-r3-04.
2. Representar en S1 el «falta capturar» del IVA (K18).
3. Fusionar «Para C1» de U1 (B6 y B7).
4. Citas que faltan en los criterios S4-07, S5-04, U1-16 y C1.
5. Cambiar por un enlace los anexos duplicados de `docs/estado/wf0/`.
6. U4 y U2 (WF-1).
7. Los 8 bloqueos para Axel.

## Solicitudes fuera de alcance
- `docs/design/U1-arquitectura.md` línea 66: `period` → `week` (H-r3-03 de C1). Dueño: U1.
- `docs/criterios/U1.md` U1-16: `period` → `week` (solicitud 3 de `correccion-C1-r2.md`, sin aplicar). Dueño: criterios.
- `docs/contratos/diccionario-datos.md`: `this_month`, `onward`, `chat_id_masked` y la partición de `budget.valid_to_date` (H-r3-04). Dueño: S1.
- `backend/internal/api/c1gen/gen3_ops.py`: `week` en `getBasketComparison` (H-r3-01) y la marca y cita de `U1-lista-vigente-domingo` en `getMealPlanWeek` (H-r3-02); después, regenerar con `build.py`. Dueño: C1.
- `docs/estado/p1/U1-1.1.md`: adenda K18 y «N-289: quitar `updatePerson`» (B6 y B7 de U1 r3). Dueño: U1.
- `docs/estado/wf0/S1-modelo-datos.md` (Anexos A y B) y `docs/estado/wf0/U1-arquitectura.md` (anexo `u1_verify.py`): cambiarlos por enlaces a `backend/db/tools/` y a `docs/design/tools/verificar_u1.py`, para que haya una sola copia.

## Paquetes instalados
Ninguno.

## Consumo aproximado
- Unas 36 acciones de herramienta.
- Unos 185 KB leídos: protocolo, §K, Makefile, extractos de README, de los índices, de los reportes y de las revisiones, y salidas de los comandos.
- Unos 55 KB escritos: parches del Makefile, del README y de los índices, INTEGRACION.md, este reporte y los scripts del scratchpad.
- `cerrado_por_tope`: false.
