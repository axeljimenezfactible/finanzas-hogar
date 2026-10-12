# Bloqueo S2-ventanas

Abierto el 11-oct-2026 por el agente `correccion-S2-r1` (WF-1), hallazgo M1 de `docs/estado/wf1/REVISION-S2-r1.md`. Se decide en la puerta 2.

Contexto: `periodos` con `{windows: [...]}` (S2-interfaz 1.2, §2.12 y §3.2). Con la semilla y una sola ventana P2 (2026-11-10 a 2026-12-09):
- El motor Go daba devengado 0, efectivo 194750 y diferencia acumulada −292125. Contaba el efectivo anterior a la ventana pero no el devengo anterior: era su supuesto 2, que cambiaba cifras sin bloqueo.
- El oráculo daba `invalid_argument`, pero solo por un costo anterior a la primera ventana; un pago anterior no lo revisaba.
- El canon de esa ventana (R-46) es devengado 194750, efectivo 194750 y diferencia acumulada 97375.
- Además, un costo con fecha en un hueco entre dos ventanas declaradas se descartaba en las dos implementaciones, sin regla escrita.

Qué falta (pregunta para Axel y para el dueño de la interfaz):
1. Con ventanas, ¿qué pasa con un costo o un pago anterior a la primera ventana pedida? §3.2 dice que `cumulative_diff_cents` va «desde el primer evento de las entradas hasta `end` (incluye periodos anteriores no pedidos)». Pero el devengo con ventanas solo está definido como «n ventanas consecutivas a partir de la que contiene su fecha», y un costo anterior no está en ninguna.
2. ¿Qué pasa con un costo o un pago fechado en un hueco entre dos ventanas declaradas?

Opciones:
- A. Rechazar con `invalid_argument` todo costo o pago anterior a la primera ventana o en un hueco entre ventanas. Es la regla que queda aplicada como PROVISIONAL. Consecuencia: quien llama tiene que declarar ventanas contiguas desde el primer evento. P2 sola no se puede pedir, y la diferencia acumulada siempre es completa y cuadra con el canon. Lo posterior a la última ventana se sigue admitiendo, porque no entra en ningún periodo ni en ninguna diferencia.
- B. Admitir ventanas que empiezan después del primer evento y devengar lo anterior con ventanas implícitas hacia atrás (del mismo largo) o con meses calendario, solo para la diferencia acumulada. Con P2 sola daría 194750, 194750 y 97375 (= R-46). Hay que escribir en §2.12 cómo se construyen esas ventanas implícitas.
- C. Admitir esas ventanas y contar devengo y efectivo solo desde el inicio de la primera ventana pedida. Hay que cambiar el texto de §3.2 («incluye periodos anteriores no pedidos»), y P2 sola ya no reproduce R-46.
- Para los huecos, cualquiera de estas tres: (H1) rechazar, que es lo aplicado; (H2) devengar el costo del hueco en n ventanas a partir de la siguiente; (H3) que el costo del hueco cuente en la diferencia acumulada pero en ningún periodo.

Impacto: con la opción A, `periodos` con ventanas que no cubren desde el primer evento, o que dejan un hueco con eventos, devuelve `invalid_argument` en el motor y en el oráculo. No cambia ninguna cifra canon. Las 27 fixtures anteriores tampoco cambian: la única con ventanas (`s2-ventanas-dia-5`) declara P1–P6, que son contiguas y empiezan en el primer evento.

Cómo seguí: apliqué A + H1 en las dos implementaciones, con la marca PROVISIONAL que cita este bloqueo:
- `backend/internal/motor_fin/ledger.go`: `uncovered` y `checkCovered`.
- `oraculo/fin_devengo.py`: `sin_cubrir`, que se llama al inicio de `resumen`.
- La fixture `oraculo/fixtures/s2/s2-ventanas-sin-cubrir.json`, con el campo `provisional` y 4 consultas: control P1+P2 = R-45 y R-46; solo P2; hueco con solo pagos (P1 y P3); hueco con solo un costo (10-dic-2026).

Además, pedí al dueño de `docs/contratos/motores/S2-interfaz.md` que, cuando se decida, escriba la regla en §2.12 y §3.2. Hasta entonces, la interfaz debería citar este bloqueo.

## Pregunta 3 · orden y traslape de las ventanas (agregada el 11-oct-2026 por `correccion-S2-r2`)

Contexto: hallazgo M1 de `docs/estado/wf1/REVISION-S2-r2.md` (S2-03 y S2-16, reabre en parte el M1 de r1). La regla A + H1 no impedía pedir ventanas fuera de orden, traslapadas o una dentro de otra. Con la semilla, todos los eventos caían en alguna ventana pedida, así que pasaban la cobertura y el motor y el oráculo daban las mismas cifras, que estaban mal:
- `[P2, P1]`: P2 daba devengado 0, efectivo 194750 y diferencia −292125. Según §3.2 y R-46 tendría que dar 194750 / 194750 / 97375.
- `[P3, P1, P2]`: el devengado de P1 daba 446487, porque la segunda mitad de la despensa del 10-dic-2026 se devengaba en la ventana de octubre (la «siguiente» en el orden recibido). Las diferencias daban −263631 y 600849.
- Traslape P1 + X (2026-10-25 a 2026-12-09): el pago del 2026-10-30 se contaba en las dos ventanas.
- A (2026-10-10 a 2027-01-09) con B (2026-11-10 a 2026-12-09) dentro: B daba devengado 446487 y diferencia 600849.

La causa es que las dos implementaciones leen «n ventanas consecutivas» (§2.12) como consecutivas en el orden recibido y calculan la diferencia acumulada como suma corrida en ese orden. Ese supuesto cambia cifras y no estaba escrito. En modo meses no pasa: cada mes tiene un lugar fijo en el calendario.

Qué falta (pregunta para Axel y para el dueño de la interfaz): ¿cómo se tratan unas ventanas que no vienen en orden cronológico, o que comparten días (traslape o una dentro de otra)?

Opciones:
- O1. Rechazar con `invalid_argument` toda lista en la que una ventana no empiece después del fin de la anterior (`start[i] ≤ end[i−1]`). Es la regla que queda aplicada como PROVISIONAL. Consecuencia: quien llama ordena las ventanas y no las traslapa. «n ventanas consecutivas» y la diferencia acumulada quedan bien definidas, y ninguna cifra canon cambia.
- O2. Ordenar las ventanas por fecha de inicio para calcular, y devolver los resultados en el orden recibido. `[P2, P1]` daría lo mismo que `[P1, P2]` (R-45 y R-46, en orden inverso). Los traslapes y las ventanas contenidas se seguirían rechazando, porque un día en dos ventanas cuenta su pago dos veces y no está definido en cuál ventana empieza el devengo. Admitirlos pediría una regla nueva en §2.12 y §3.2.

Impacto: con O1, `periodos` con ventanas desordenadas, traslapadas (aunque sea un día) o contenidas devuelve `invalid_argument` en el motor y en el oráculo. Ninguna cifra canon cambia. Las 30 fixtures siguen siendo 30. Solo cambia `s2-ventanas-sin-cubrir`, que pasa de 4 a 9 consultas.

Cómo seguí: apliqué O1 en las dos implementaciones, con la marca PROVISIONAL que cita este bloqueo. Cada lado lo escribe con su propia forma:
- `backend/internal/motor_fin/ledger.go` (`periods`): compara cada inicio con el fin de la ventana anterior.
- `oraculo/fin_devengo.py` (`desordenadas`, que se llama al inicio de `resumen` antes de `sin_cubrir`): compara cada inicio con el fin más tardío de todas las anteriores. Con ventanas válidas (inicio ≤ fin) las dos formas son equivalentes.
- `oraculo/fin_genera_fixtures.py` agrega 5 consultas con error esperado a `s2-ventanas-sin-cubrir`: `desordenadas-p2-p1`, `desordenadas-p3-p1-p2`, `traslape-p1-x`, `traslape-un-dia` (Y empieza el 2026-11-09, el último día de P1; es la frontera `start = end` anterior) y `contenida-a-b`.
- Prueba Go `TestVentanasOrden` (`cover_test.go`): comprueba que el rechazo viene de la regla de orden y no de la de cobertura. El control P1, P2 y P3 reproduce R-45…R-47.

Pedí al dueño de `docs/contratos/motores/S2-interfaz.md` que §2.12 diga «ventanas en orden ascendente y disjuntas» y que §2.12 y §3.2 citen este bloqueo.
