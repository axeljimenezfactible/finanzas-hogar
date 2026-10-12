# correccion-criterios-r2 · P2c · 11-oct-2026, noche

Alcance: `docs/criterios/*.md`, `INDICE.md` y este reporte. Fuente: `docs/estado/p2c/REVISION-criterios-r2.md`, hallazgo M-1 (media). No toqué ids ni severidades. No hice git add ni commit.

## Hecho
- **HECHO.** M-1 corregido con una sola regla, por pregunta, en tres lugares:
  - `docs/criterios/C1.md` C1-15, cuerpo: una marca es vigente si cita un bloqueo sin «Respuesta de Axel», o una pregunta que el bloqueo declara ABIERTA en su «Estado» o en «Estado del bloqueo (puerta 2): ABIERTO». Un bloqueo con la respuesta completa está RESUELTO y su marca se quita.
  - `docs/criterios/C1.md` C1-15, Verifica: un bloqueo cuenta como abierto si `grep -L "Respuesta de Axel"` lo lista o si `grep -n "Estado del bloqueo (puerta 2): ABIERTO"` encuentra la línea; con respuesta parcial se comprueba además que la pregunta citada por la marca sea la abierta (pregunta 4 de `S2-ventanas`, Q2-12 de `S2-flujo-quincenal`, Q2-48 de `S5-carrito-comparativa`).
  - `docs/criterios/INDICE.md` §6 (párrafo «Estado =»): misma regla por pregunta.
- **HECHO.** Quité de C1-15 la afirmación «Estado: ABIERTO» sobre los cuatro bloqueos de C1. Ahora dice que todos carecen de «Respuesta de Axel»; solo `S3-iva-lectura` y `S1-reparto-edicion-gasto` traen «Estado: ABIERTO»; `C1-hora-reportes-job` y `C1-reparto-alcance` no traen línea «Estado» y no hace falta. Con esto las marcas de `S2-interfaz.md` (l. 155, 177, 258) y `S5-interfaz.md` (l. 31, 74) ya no son defecto y S2-09, S2-10, S2-16 y S5-13 dejan de contradecir a C1-15.

## Verificación
- `make verify-criterios`: `criterios definidos: 261`, `IDs distintos citados: 261`, `citados que no existen: {}`; `citas verificadas: 321`, `sin fuente: 0`; `TOTAL comprobadas: 178`, `con diferencias: 0`; `cifras comparadas: 12`, `diferencias: 0`; `TOTAL 261 (alta 201 / media 60 / baja 0)`, 18 archivos; `INDICE contra conteo propio: 0 diferencias`; `residuos «N-xx» y rutas /tmp en los criterios: 0 y 0; las 11 comprobaciones dan lo esperado`.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: `WF-1: verificación completa OK` y `Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)`.
- Regla nueva contra los bloqueos: `grep -c "Estado del bloqueo (puerta 2): ABIERTO"` da 1 en `S2-ventanas`, `S2-flujo-quincenal` y `S5-carrito-comparativa`; `grep -L "Respuesta de Axel"` lista `C1-hora-reportes-job`, `C1-reparto-alcance`, `S3-iva-lectura`, `S1-reparto-edicion-gasto` y `S5-precio-vigente-captura`. Todas las marcas vigentes citadas quedan cubiertas por la regla.

## Supuestos
- Ninguno que cambie cifras, esquema o endpoints. Elegí quitar la afirmación «Estado: ABIERTO» (primera opción de la corrección sugerida) y no pedir que se agregue la línea a los bloqueos.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- B-1 y B-2 (bajas) de la revisión no estaban en mi encargo y no se tocaron. Tampoco K53 ni U2 (135°).

## Propuestas para Axel
- Ninguna.

## Solicitudes fuera de alcance
- Opcional: el dueño de los bloqueos puede agregar la línea «Estado: ABIERTO» a `C1-hora-reportes-job`, `C1-reparto-alcance` y `S5-precio-vigente-captura`; ya no hace falta para C1-15.
- Notas del revisor al integrador (marcas PROVISIONAL viejas en `ocr/types.go`, `pre_genera_fixtures.py`, motor S4) siguen vigentes; no son criterios.

## Consumo aproximado
- Acciones de herramienta: 8 · KB leídos: ~45 · KB escritos: ~5
