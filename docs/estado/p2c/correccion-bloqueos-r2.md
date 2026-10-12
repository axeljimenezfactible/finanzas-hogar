# Corrección de bloqueos · ronda 2 (11-oct-2026, noche)

Agente: `correccion-bloqueos-r2`. Entrada: `docs/estado/p2c/REVISION-bloqueos-r2.md` (hallazgos alta y media: H-r2-01, H-r2-02, H-r2-03). Alcance: los bloqueos, `INDICE.md`, `decisiones.md` §K46 y este reporte. No hice `git add` ni `git commit`. Las bajas (H-r2-04, H-r2-05, H-r2-06 y H-r2-07) no estaban en la orden y no las toqué.

## Hecho

### H-r2-02 (media) · INDICE y K46 contra el disco: corregido
HECHO. La carpeta tiene hoy **75 archivos de bloqueo, 63 RESUELTOS y 12 ABIERTOS**, no los 72, 63 y 9 de la revisión. Es 69 de HEAD (63 y 6) más 4 de P2c (la revisión vio 3: `S5-carritos-granel` se abrió después) y 2 de la Pista 3 (`P3-contraste-chip-hover`, `P3-huecos-componentes`, dejados por `ds-brandbook`). Conté con la regla de la puerta 2 (última línea «Estado del bloqueo (puerta 2)»; si falta, «Respuesta de Axel»; si no, «Estado»).
Cambios en `docs/estado/bloqueos/INDICE.md`:
- l. 12 (estado de la puerta 2): la pregunta 4 de `S2-ventanas` queda abierta con los casos 1 y 2; menciona los 4 bloqueos de P2c y los 2 de la Pista 3.
- Resumen: párrafo con el conteo actual (75, 63 y 12); filas «Flujo P2c | 4 | 0 | 4» y «Pista 3 | 2 | 0 | 2»; total 75, 63 y 12.
- l. 28: se quitó «Está por crear…». Las herramientas cuentan **67 resueltos y 8 abiertos** (comprobado con la regla de `verify_c1.py` l. 322-323): los 2 de WF-1 sin «Respuesta de Axel» (`S5-carritos-detalle`, `S5-profeco-filtro`), los 4 de P2c y los 2 de la Pista 3. La revisión pedía «67 y 5»; con los 3 bloqueos que ella contó y sin los de la Pista 3 serían 67 y 5, y hoy son 67 y 8.
- Sección nueva «Bloqueos nuevos de P2c» con la pregunta, las opciones y lo PROVISIONAL de `S2-ahorro-quincenal`, `S2-metas-vencida`, `S4-presentacion-respaldo` y `S5-carritos-granel`, tomados de sus archivos, sin recomendación. Sección breve «Bloqueos nuevos de la Pista 3» que solo los cuenta y los remite a su dueño.
- «Respuestas de la puerta 2»: la fila Q2-01 pasa a «casos 1 y 2» (1A y 2A PROVISIONALES) y la introducción remite la variante de Q2-07 a `S4-presentacion-respaldo`.
Cambios en `docs/plan/decisiones.md`:
- K46: se quitó «por crear … que aún no existe»; ahora dice que el bloqueo existe, está ABIERTO y entra en el conteo de 75. La pregunta 4 de `S2-ventanas` queda abierta con los casos 1 y 2 (antes solo el 2).
- Párrafo final (l. 150): añade los 4 bloqueos nuevos de P2c.

### H-r2-03 (media) · cita de Q2-12: corregida
HECHO. `docs/estado/bloqueos/S2-flujo-quincenal.md`, entrada Q2-12: ya no cita líneas, cita «`S2-interfaz.md` 1.3 §3.11, fila `variable_estimated_cents`, y §8, fila «Variable estimado de M-10»». Comprobado en el archivo actual: esas dos filas existen (l. 258 y l. 460 hoy).

### H-r2-01 (alta) · verificación global rota: no se corrige en los bloqueos
Sigue **sin resolverse, y no es de mi alcance**. La corrección sugerida dice expresamente que los bloqueos no se tocan (quitar «Respuesta de Axel» de `S1-parametros-avisos.md` escondería una decisión de Axel, Q2-11 = B). El arreglo es del agente que aplica Q2-11 = B. Estado verificado hoy en el árbol vivo (`grep PROVISIONAL … S1-parametros-avisos`):
- `docs/contratos/openapi.yaml`: 6 marcas (l. 17911 a 17960); `docs/contratos/C1-datos.md`: 6 (l. 572, 576, 772 a 775); `backend/db/migrations/00003_puerta1b.sql`: 5 (l. 56 a 59 y 62).
- `diccionario-datos.md` y `reglas-negocio.md` ya no tienen marcas (alguien las quitó en el árbol); la migración `00004_puerta2.sql` ya trae `@resuelve S1-parametros-avisos`.
Ver solicitud 1.

## Verificación

| Qué | Comando | Salida |
|---|---|---|
| Conteo de la carpeta | script propio con la regla de la puerta 2 | `75 {'RESUELTO': 63, 'ABIERTO': 12}`; nuevos respecto de HEAD: los 6 listados; HEAD: 63 y 6 |
| Conteo con la regla de las herramientas | script con `Respuesta de Axel\|Estado: RESUELTO` | 67 resueltos y 8 abiertos |
| Criterios | `make verify-criterios GOOSE=…` | rc=0; «INDICE contra conteo propio: 0 diferencias» |
| U1 | `make verify-u1 GOOSE=…` | rc=0; «RESULTADO: TODO OK» |
| S1 en el árbol vivo | `make verify-s1 GOOSE=…` | «RESUMEN TODO OK» |
| C1 en el árbol vivo | `make verify-c1 GOOSE=…` | **FALLA**: `gen2: tipo de parámetro desconocido: hora HH:MM`, Error 1. Viene de `00004_puerta2.sql` y del generador, que otro agente está aplicando; no es de mis archivos |
| Global | `make verify-all GOOSE=…` | pasa WF-0 (SQL, Go, colorimetría, OpenAPI) y se detiene en `verify-c1` por lo de arriba. **No pasa** |
| Efecto aislado de mis cambios | clon de HEAD con solo los bloqueos, INDICE y `decisiones.md` del árbol (ya borrado) | verify-c1: «21 archivos coinciden byte a byte» y `FALLA C1-15`; verify-s1: `RESUMEN FALLAS: ['S1-16', 'K44']`; verify-criterios y verify-u1: rc=0. Es la misma falla de H-r2-01 y se debe a las 17 marcas que siguen en `openapi.yaml`, `C1-datos.md` y la migración 00003, no a lo que cambié hoy |

El criterio de cierre «`make verify-all` pasa con los dos cambios juntos» **no se cumple todavía** y no puede cumplirse desde esta carpeta.

## Supuestos
- Conté los 2 bloqueos de la Pista 3 porque están en `docs/estado/bloqueos/` y el INDICE debe coincidir con el disco; no los edité ni redacté sus opciones (solo las resumí con sus palabras). Si la Pista 3 abre otro, el conteo cambia.
- «Fase que lo usa» de los 4 bloqueos de P2c sale del «Impacto» de cada archivo (S2 M-02 y M-10, pantalla de Metas, S4, S5), más WF-3, donde dicen que se deciden.
- Q2-01 = A autoriza escribir la pregunta 4 pero no la cierra (PUERTA-2 §2.1); por eso la fila dice «sin respuesta».

## Bloqueos
Ninguno nuevo. Siguen vigentes los de los 4 bloqueos de P2c y los de WF-1 que ya estaban en el INDICE.

## Pendientes
- H-r2-01 completo (ver solicitud 1).
- Bajas fuera de la orden: H-r2-04 (`S2-ventanas.md` l. 85, «Qué falta» del caso 2 obsoleto), H-r2-05 (`$` en `U3-t2-acumulados.md` l. 24-25), H-r2-06 (RESPUESTAS.md no lista Q2-12 como abierta; es del orquestador) y H-r2-07 (columnas «Decisión» y «K» en la tabla de la puerta 2 del INDICE).
- Si `contratos-puerta2` o la Pista 3 abren otro bloqueo, recontar el INDICE.

## Propuestas para Axel
Ninguna nueva. Quedan 4 preguntas de P2c sin recomendación (ver «Bloqueos nuevos de P2c» en el INDICE); se deciden en WF-3 o por id. Pregunta para el orquestador: las dos de la Pista 3 (`P3-contraste-chip-hover` recomienda A) no tienen cabida en la tabla de la puerta 2; decidir si se llevan a Axel junto con las demás.

## Solicitudes fuera de alcance
1. (Alta, H-r2-01) Al agente que aplica Q2-11 = B: quitar las marcas PROVISIONAL que citan `S1-parametros-avisos` en `docs/contratos/openapi.yaml` (6, por `c1gen`), `docs/contratos/C1-datos.md` (P-26.08, P-26.09 y las filas de params_json, 6), `backend/db/migrations/00003_puerta1b.sql` (5; es una migración ya aplicada, valorar si se corrige el comentario o se deja la nota en 00004) y los criterios C1-22, U4-12 y S7. Arreglar además `gen2: tipo de parámetro desconocido: hora HH:MM` para que `make verify-c1` vuelva a correr. El integrador no hace commit de la carpeta de bloqueos sin ese cambio. Cierra cuando `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` pasa con todo junto.
2. Al integrador: confirmar el conteo del INDICE (75, 63 y 12) contra la carpeta justo antes del commit.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/INDICE.md`
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S2-flujo-quincenal.md`
- `/home/deepsec/finanzas-hogar/docs/plan/decisiones.md`
- `/home/deepsec/finanzas-hogar/docs/estado/p2c/correccion-bloqueos-r2.md`

## Consumo aproximado
27 acciones de herramienta · unos 130 KB leídos · unos 17 KB escritos.
