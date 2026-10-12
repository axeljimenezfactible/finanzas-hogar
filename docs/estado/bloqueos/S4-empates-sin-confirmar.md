# Bloqueo S4-empates-sin-confirmar

Abierto el 11-oct-2026, noche, en la ronda P2c por el agente `correccion-contratos-r1` (hallazgos M-1 y M-4 de `docs/estado/p2c/REVISION-contratos-r1.md`), al quitar de S4 y S5 las marcas que citaban `S4-confirmacion-empates`, que ya está resuelto.

Contexto: Q2-43 = A (bloqueo `S4-confirmacion-empates`, RESUELTO; K46) fija quién confirma y cuándo: Axel, en la puerta 2, sobre la lista de renglones `propuesto`. La confirmación en sí no ocurrió. La lista de nutrición es la Q2-90 («confirmar o corregir los 40 empates de `datos/semilla/nutricion/revisar.tsv`»), que no trae recomendación y quedó sin respuesta. `docs/estado/puerta2/RESPUESTAS.md` lo dice así: «Q2-85, Q2-87 a Q2-91 y Q2-94 sin recomendación → siguen como están, PROVISIONAL».

HECHO (11-oct-2026):
- `datos/semilla/nutricion/tabla.csv` tiene 150 renglones, los 150 con `estado` = `propuesto`.
- `datos/semilla/profeco/vinculos.csv` no existe en el árbol de trabajo, así que todavía no hay lista de ligas PROFECO que confirmar.
- En modo `production`, el cargador no carga ningún renglón `propuesto` (FIRME: `S4-interfaz.md` §9.1; para las ligas, `S5-interfaz.md` §11).
- Las fixtures `s4-kcal-recetas-autor` y `s4-kcal-sin-dato` usan los valores de `tabla.csv` en modo `review`. Su `provisional` actual dice que, si al confirmar se corrige un valor, se regeneran el recetario y la fixture.

Qué falta: que Axel confirme o corrija los empates de nutrición (Q2-90) y, cuando exista, la lista de ligas PROFECO.

Opciones (las letras de Q2-90): A) confirmar tal cual · B) corregir; Axel dice cuáles. Sin recomendación: hacen falta datos que solo Axel tiene.

Impacto: mientras no se confirmen, la app en producción arranca sin kcal (y sin ligas PROFECO), y S4-08 se reporta con dos conteos, confirmados y propuestos. Ninguna cifra canon depende de esto.

Cómo seguí: `S4-interfaz.md` §9.1 da por FIRME quién confirma y cuándo (Q2-43 = A). Los valores de los renglones `propuesto` quedan PROVISIONALES, con cita de este archivo. Las dos fixtures `s4-kcal-*` pasan a citar este archivo en `provisional` (cambio para WF-3).

- **Estado:** ABIERTO
