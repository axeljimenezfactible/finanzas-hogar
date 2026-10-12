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

## Respuesta de Axel (puerta 2, 11-oct-2026)
Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo de respuesta elegido por Axel: «lo recomendado + excepciones» (se aprueba en bloque la opción recomendada; hasta ahora no hay excepciones por id). Las letras son las de este archivo, salvo donde se indica.

### Q2-25 · (1) ¿Qué pasa con un costo o pago anterior a la primera ventana pedida?
- **Respuesta:** **A** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Rechazar con `invalid_argument` todo costo o pago anterior a la primera ventana o en un hueco entre ventanas. Es la regla que queda aplicada como PROVISIONAL.»
- **Consecuencia:** quien llama tiene que declarar ventanas contiguas desde el primer evento. P2 sola no se puede pedir, y la diferencia acumulada siempre es completa y cuadra con el canon. Lo posterior a la última ventana se sigue admitiendo, porque no entra en ningún periodo ni en ninguna diferencia.
- **Qué cambia:** `ledger.go` (`uncovered`, `checkCovered`) y `oraculo/fin_devengo.py` (`sin_cubrir`) quitan la marca PROVISIONAL; fixture `s2-ventanas-sin-cubrir`; `S2-interfaz.md` §2.12 y §3.2 escriben la regla (una vez aplicada Q2-01, que redacta la pregunta 4 completa).
- **Estado:** RESUELTO

### Q2-26 · (2) ¿Y un costo o pago fechado en un hueco entre ventanas?
- **Respuesta:** **H1** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «(H1) rechazar, que es lo aplicado»
- **Consecuencia** (el archivo no la separa por opción; es su «Impacto»)**:** Impacto: con la opción A, `periodos` con ventanas que no cubren desde el primer evento, o que dejan un hueco con eventos, devuelve `invalid_argument` en el motor y en el oráculo. No cambia ninguna cifra canon. Las 27 fixtures anteriores tampoco cambian: la única con ventanas (`s2-ventanas-dia-5`) declara P1–P6, que son contiguas y empiezan en el primer evento.
- **Qué cambia:** Igual que Q2-25: `S2-interfaz.md` §2.12 y §3.2; `ledger.go` y `fin_devengo.py` conservan el rechazo sin marca PROVISIONAL.
- **Estado:** RESUELTO

### Q2-27 · (3) ¿Y ventanas fuera de orden, traslapadas o una dentro de otra?
- **Respuesta:** **O1** (la recomendada, aprobada en bloque)
- **Texto de la opción (literal del archivo):** «Rechazar con `invalid_argument` toda lista en la que una ventana no empiece después del fin de la anterior (`start[i] ≤ end[i−1]`). Es la regla que queda aplicada como PROVISIONAL.»
- **Consecuencia:** quien llama ordena las ventanas y no las traslapa. «n ventanas consecutivas» y la diferencia acumulada quedan bien definidas, y ninguna cifra canon cambia.
- **Qué cambia:** `ledger.go` (`periods`), `fin_devengo.py` (`desordenadas`) y la prueba `TestVentanasOrden`; `S2-interfaz.md` §2.12 dice «ventanas en orden ascendente y disjuntas» y cita este bloqueo.
- **Estado:** RESUELTO

## Puerta 2 (11-oct-2026): sin respuesta
Fuente: `docs/estado/puerta2/RESPUESTAS.md` («Las preguntas sin recomendación quedan sin respuesta en esta puerta: PROVISIONAL con la opción que ya está en uso, y pasan a WF-3») y hoja de respuestas de `PUERTA-2.md` §8. Axel puede mandar una excepción por id cuando lea el documento.

### Q2-01 (caso 2, pregunta 4 del bloqueo) · Con ventanas contiguas que no duran un mes, la diferencia acumulada al mismo día de corte depende de cómo se corten las ventanas ($973.75 contra $2,921.25 al 2026-11-09; `docs/estado/wf1/REVISION-S2-r3.md` l. 45)
- PROVISIONAL: sigue en uso el comportamiento actual de motor y oráculo (`ledger.go`, `fin_devengo.py`), que no está escrito (PUERTA-2 §2.1: «Caso 2: sin recomendación»). Q2-01 = A autoriza escribir la pregunta 4 completa y la regla PROVISIONAL, pero no cierra el caso 2 (lo dice PUERTA-2 §2.1). La pregunta 4 quedó redactada abajo en P2c (agente `contratos-puerta2`). Pasa a WF-3.
- **Qué falta:** la pregunta 4 del bloqueo con sus opciones y su impacto. Hay dos fuentes: el hallazgo M1 y la pregunta 4 que el revisor propone (`docs/estado/wf1/REVISION-S2-r3.md` l. 38-48) y `docs/estado/puerta2/PUERTA-2.md` §2.1 («Caso 2: sin recomendación»). No se escriben opciones aquí: las redacta quien aplique Q2-01 y Axel las decide; cuando exista, esta línea se sustituye por la pregunta completa.
- **Estado:** ABIERTO

**Estado del bloqueo (puerta 2): ABIERTO** — Q2-25, Q2-26 y Q2-27 están resueltas; queda abierto el caso 2 (ventanas contiguas que no duran un mes), que es la pregunta 4 del bloqueo, redactada abajo en P2c; el caso 1 también queda abierto, con la regla 1A como PROVISIONAL.

## Pregunta 4 · devengo con ventanas: hueco sin eventos y ventanas que no duran un mes (agregada el 11-oct-2026 en P2c por `contratos-puerta2`; Q2-01 = A)

Contexto: hallazgo M1 de `docs/estado/wf1/REVISION-S2-r3.md` (S2-16 y S2-03). Con ventanas, un costo de n meses se devenga en «n ventanas consecutivas» a partir de la que contiene su fecha (`S2-interfaz.md` §3.2), no en n meses. Las reglas A, H1 y O1 (preguntas 1 a 3, FIRMES desde la puerta 2) no cubren dos casos que tienen esa misma causa. En los dos, el motor (`ledger.go`, `windowed`) y el oráculo (`fin_devengo.py`, `k0 + j`) dan las mismas cifras; la regla no estaba escrita.

Qué falta (pregunta para Axel): con ventanas, ¿el devengo «mitad por mes» sigue a la ventana o al mes? En concreto:
- **Caso 1 · hueco sin eventos que cruza un devengo.** Datos del revisor: la semilla sin compras a crédito más un gasto de contado de $1,000.00 (100000) el 2026-10-10 con `accrual_months` 2. Con las ventanas A (2026-10-10 a 2026-11-09) y C (2026-12-10 a 2027-01-09), y el hueco B (2026-11-10 a 2026-12-09) sin eventos: A da devengado 50000, efectivo 100000 y diferencia −50000; C da devengado **50000**, efectivo 0 y diferencia 0. Con A, B y C contiguas, C da devengado **0** (B lleva los 50000). En modo meses `[2026-10, 2026-12]`, diciembre da **0**: la mitad va a noviembre, que no se pidió.
- **Caso 2 · ventanas contiguas que no duran un mes.** Con la semilla, `[P1]` da diferencia acumulada **97375** ($973.75) al 2026-11-09, y `[d = 2026-10-10, r = 2026-10-11 a 2026-11-09]` da **292125** ($2,921.25) al mismo 2026-11-09, porque la segunda mitad de la despensa 1 cae en r. La interfaz define esa diferencia «desde el primer evento … hasta `end`», y con las mismas entradas y el mismo `end` cambia según el corte. Las ventanas d y r son contiguas, ascendentes y disjuntas, así que la regla del caso 1 no lo cierra (PUERTA-2 §2.1).

Opciones del caso 1 (hueco sin eventos):
- **1A.** Rechazar con `invalid_argument` toda lista con hueco (`start[i] ≠ end[i−1] + 1 día`), tenga o no eventos. A + C da `invalid_argument`. Es la regla que queda aplicada como PROVISIONAL desde la puerta 2 (`S2-interfaz.md` 1.3, §2.12 regla 4).
- **1B.** Rechazar solo el hueco que cruza un devengo: el que recibiría una parte de algún costo si las ventanas fueran contiguas. A + C da `invalid_argument`; un hueco que ningún devengo cruza se admite.
- **1C.** Conservar el comportamiento anterior: la parte del hueco pasa a la siguiente ventana declarada. A + C da en C devengado 50000, distinto del modo meses (0).
- **1D.** La parte del hueco no entra en ningún periodo. A + C da en C devengado 0, como el modo meses. Es lo que resulta de la opción 2B.

Opciones del caso 2 (ventanas que no duran un mes):
- **2A.** El devengo sigue a la ventana: la parte k va a la ventana k posiciones después de la que contiene la fecha. `[P1]` da 97375 y `[d, r]` da 292125 al 2026-11-09: la diferencia acumulada depende del corte. Es el comportamiento actual de motor y oráculo, escrito como PROVISIONAL en `S2-interfaz.md` 1.3 §3.2.
- **2B.** El devengo sigue al mes: la parte k se fecha en la fecha del costo + k meses (mismo día, recortado al último día del mes, como las ocurrencias `monthly` de §2.7) y va a la ventana que contiene esa fecha; si cae en un hueco o después de la última ventana, no entra en ningún periodo. La diferencia acumulada al mismo `end` deja de depender del corte: `[d, r]` daría 97375 al 2026-11-09, igual que `[P1]`. Esta cifra sale de la regla (la segunda mitad de la despensa 1 quedaría fechada el 2026-11-10, fuera de r); ningún motor la implementa todavía. Con las ventanas del oráculo (P1–P6, del día 10 al 9) da lo mismo que hoy (R-45…R-50).
- **2C.** Solo ventanas de un mes: cada ventana termina el día anterior al mismo día del mes siguiente (recortado al fin de mes). Así «n ventanas consecutivas» son n meses. `[d, r]` da `invalid_argument` y `[P1]` da 97375.

Sin recomendación: PUERTA-2 §2.1 deja el caso 2 sin recomendación y ningún reporte propone opciones; la 1A es la regla provisional que autorizó Q2-01, no una recomendación.

Impacto: con 1A, `periodos` con una lista de ventanas con hueco devuelve `invalid_argument` en el motor y en el oráculo. Ninguna cifra canon cambia: la única fixture con ventanas válidas, `s2-ventanas-dia-5`, declara P1–P6, que son contiguas. Con 2A no cambia nada; con 2B o 2C solo cambian consultas con ventanas que no duran un mes, y ninguna es canon.

Cómo seguí (P2c): 1A y 2A quedan escritas como PROVISIONAL en `S2-interfaz.md` 1.3 (§2.12 regla 4 y §3.2), con estas cifras. Falta que los agentes de motor y oráculo (WF-3) implementen 1A, porque hoy `uncovered` y `sin_cubrir` solo rechazan eventos fechados en un hueco, y que agreguen a `s2-ventanas-sin-cubrir` la consulta A + C (`invalid_argument`) y la consulta `[d, r]` (292125 al 2026-11-09, PROVISIONAL 2A).

- **Estado de la pregunta 4:** ABIERTO (pasa a WF-3; Axel puede contestar por id).
