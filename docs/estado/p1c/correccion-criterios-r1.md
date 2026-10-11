# correccion-criterios-r1 (P1c, 11-oct-2026)

Corrige los hallazgos alta y media de `docs/estado/p1c/REVISION-criterios-r1.md` (A-01, A-02, A-03, M-01, M-02). Los bajos (B-01 a B-08) no se tocaron (fuera del encargo).

## Hecho
- **A-01 · C1-22** (`docs/criterios/C1.md`): los parámetros de avisos (umbrales 80/100 %, días antes del pago, `send_time`, `advance_days`) se sacan de la lista de claves de `app_setting`; viven hoy en `notification.params_json` (K43) y su lugar final es PROVISIONAL citando el bloqueo abierto `S1-parametros-avisos`. La Verifica los excluye de la comprobación contra 00003 y los comprueba aparte (existen, valor inicial, se cambian con `updateNotification`, «restablecer», marca). La marca cubre el lugar, no los valores. Evidencia pide tabla aparte. Entrada en «Cambios de la versión 1.2».
- **A-02 · U4-12** (`docs/criterios/U4.md`): la columna «clave» admite `notification.params_json` para los avisos con marca PROVISIONAL citando el bloqueo; el cruce con 00003 los excluye. Entrada en «Cambios».
- **A-03 · S7** (`docs/criterios/S7.md`): S7 sube a 1.2. S7-08 (4) lleva la regla FIRME de K43 sin marca (casos de prueba: 21 armadas, una sin receta ni «fuera», «fuera» no cuenta); solo el valor inicial de la anticipación de (3) queda PROVISIONAL (K31 y el bloqueo). S7-04, S7-08 y S7-10 exigen leer `send_time`, `advance_days`, umbrales y días antes del pago de su parámetro configurable (K44) con prueba de que un cambio se respeta y de que «restablecer» vuelve a la matriz original. §3 (defectos) y §5 (citas K43 y K44) actualizados.
- **M-01 · INDICE §6**: fila nueva de `S1-parametros-avisos.md` con estado ABIERTO, pregunta y criterios afectados (C1-22, U4-12, S7-04, S7-08, S7-10); la frase de estado dice 25 filas: 24 resueltas y 1 abierta. §10: S7 en 1.2. §11: cuatro filas de la revisión r1. Cabecera actualizada.
- **M-02 · P3** (`docs/criterios/P3-claude-design.md`): P3-01 queda en «paquete listo y revisado» (paquete completo, revisor propio sin altos, `sha256sum -c` en 0); el aviso de Axel «iter-N lista» pasa a la Verifica y la Evidencia de P3-06, con la Pista 3 «pendiente de Axel» hasta que avise y hasta la puerta 4.

## Archivos
`docs/criterios/C1.md`, `U4.md`, `S7.md`, `P3-claude-design.md`, `INDICE.md` y este reporte. No se tocó `docs/criterios/tools/`.

## Verificación
- `bash docs/criterios/tools/verif_final.sh` (verificación completa de criterios): conteo 261 (201 / 60 / 0) = INDICE; traza: 18 archivos, 0 criterios sin línea en §5; refs: 261 definidos y 261 citados, 0 inexistentes; citas: 305 verificadas, 0 sin fuente; N-xx residuos 0; canon 178 comprobadas, 0 diferencias; S1-09 SQL y equivalencia con el oráculo: 12 de 12 cifras; rutas /tmp: 0; `conteo.py`: «INDICE contra conteo propio: 0 diferencias» (S7 1.2 con su lista «Cambios de la versión 1.2»).
- Paso 9 (`matriz_c102.py`) reporta N-320…N-327 sin operación en la matriz de C1: es el hallazgo bajo B-01 de la revisión (C1-datos.md, fuera de alcance), no cambió con estas correcciones.
- `grep -c S1-parametros-avisos`: C1 3, U4 3, S7 5, INDICE 4 (antes 0 en los cuatro).

## Supuestos
- S7-10: la hora del job mensual se lee de «un parámetro configurable» (K44); si es `send_time` u otro lo documenta el agente de S7. No se fijó valor.
- La marca PROVISIONAL sobre los avisos cubre el lugar de los parámetros, no sus valores iniciales (80/100, [1, 0]), porque son los decididos y no cambian en ninguna opción del bloqueo.
- `send_time` y `advance_days` figuran en 00003 como `@param_sin_valor` (claves de K43 sin valor); su valor inicial lo propone S7 (K31).

## Bloqueos
Ninguno nuevo. Sigue abierto `docs/estado/bloqueos/S1-parametros-avisos.md` (sin «Respuesta de Axel»).

## Pendientes
- Al responder Axel el bloqueo: quitar la excepción de C1-22 y U4-12, las marcas de S7-04/S7-08 sobre el lugar y actualizar INDICE §6.
- Hallazgos bajos B-01 a B-08 de `REVISION-criterios-r1.md` (no pedidos).

## Solicitudes fuera de alcance
- Ninguna nueva (B-01: `docs/contratos/C1-datos.md` §10 y `docs/design/U1-arquitectura.md` §3.29 ya están en la revisión).

## Propuestas para Axel
- Responder `S1-parametros-avisos` en la puerta 2: la opción A (todo en `notification.params_json`) no cambia esquema ni contrato y es la que ya usa el esquema 1.3.

## Consumo aproximado
Acciones de herramienta: 22; leído ≈ 130 KB (por partes); escrito ≈ 18 KB.
