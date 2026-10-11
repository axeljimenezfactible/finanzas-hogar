# Bloqueo criterios-nxx
**Estado: RESUELTO el 10-oct-2026 con la opción (b); ver «Resolución» al final. Lo de abajo es el texto original.**

Contexto: U1-02, C1-02, C3-02: la tarea del flujo habla de «toda N-xx de U1 mapeada a un endpoint».
Qué falta: plan.md y decisiones.md no definen el identificador «N-xx». Se asumió que es el identificador de pantalla o nodo del mapa de U1 (formato fijado por U1). ¿Es correcto o N-xx significa otra cosa (p. ej. necesidades)?
Opciones:
- (a) N-xx = identificador de pantalla o nodo del inventario de U1 (lo asumido).
- (b) N-xx = necesidad o requisito de usuario (entonces se necesita su lista fuente).
- (c) Otro significado indicado por Axel.
Impacto: Si (b), la matriz de U1 y C1 cambia de unidad; los criterios U1-02/03, C1-02, C3-02 se reescriben.
Cómo seguí: Los criterios dicen «identificador N-xx (formato que fija U1)» y miden unicidad y cobertura sin depender del significado.

## Resolución (10-oct-2026, ronda 1 de correcciones de WF-0, agente `correccion-criterios-r1`)
Opción elegida: **(b)**. «N-xx» es una necesidad de datos. U1 (`docs/design/U1-arquitectura.md` §0.1) numera las necesidades `N-xxx` (199 en su §3, tipo L o E) y las pantallas `P-xx` (`P-26.xx` para las 14 secciones de Ajustes); C1 mapea cada `N-xxx` a operaciones (`C1-datos.md` §10) y C3 y U5 nombran las pantallas con `P-xx`.
Base de la lectura: la tarea de U1 definía N-xx como necesidad de datos («tabla de NECESIDADES DE DATOS con id N-xx … el contrato C1 mapeará cada N-xx a endpoints», citada en `docs/estado/wf0/U1-arquitectura.md`, solicitud 1(a)); la pidieron también el reporte de C1 (`docs/estado/wf0/S6a-contrato-c1.md`, solicitud 4) y el revisor independiente (`docs/estado/wf0/REVISION-r1.md`, hallazgo H-06).
Efecto: no cambia cifras, esquema ni endpoints; solo el identificador con el que los criterios nombran lo que U1 y C1 ya entregaron. Si Axel prefiriera otra lectura (opción c), se reabre.
Aplicado en `docs/criterios/`: U1-02, U1-03, U1-04, U1-05, U1-07, U1-08, U1-09, C1-02, C1-10, C3-01, C3-02, C3-08, U5-01, U5-03, U4-01, P3-06 y `INDICE.md` (ver su §9).
