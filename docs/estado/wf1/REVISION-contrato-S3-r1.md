# Revisión del contrato S3 · ronda 1 de 3

Revisor: `revisor-contrato-S3-r1` (independiente; no escribió nada de lo revisado y no corrige) · 11-oct-2026 · WF-1.
Objeto: `docs/contratos/motores/S3-interfaz.md` v1.0 (46,597 bytes) y el reporte `docs/estado/wf1/contrato-S3.md`.
Vara: **solo** `docs/criterios/S3.md` v1.2, decisiones §C y §K. Todo lo numérico lo recalculé con scripts propios en `datos/tmp/rev-s3r1/` (ya borrados).

## Veredicto
**NO APROBADO**: 0 hallazgos altos, **8 medios** y 4 bajos. Las cifras están bien: con las reglas escritas en la interfaz reproduje sin una sola diferencia las hojas S3 de `model.json`, y el conteo de las 437 hojas cuadra. Lo que falta es que el contrato alcance para implementar sin preguntar: choca con el formato común de fixtures (`oraculo/fixtures/FORMATO.md`), deja tipos sin declarar, deja sin definir las 30 filas S2 de `TestCanon85` y deja ambiguas cuatro reglas del comparador y del calendario.

## Verificaciones propias
| # | Qué | Comando o método | Resultado |
|---|---|---|---|
| a | Conteo de hojas por rama y cruce con §12 | `python3 -I datos/tmp/rev-s3r1/hojas.py oraculo/model.json`: recorre `model.json`, excluye booleanos, textos y null, y casa cada hoja con los 42 renglones de §12 (transcritos a mano) | `total 596` · S3 **437** (`parametros` 9, `despensas` 144, `credito` 143, `sensibilidad` 45, `disposicion_por_despensa` 68, `interes_por_monto` 28) · S2 159 · otras 0 · suma declarada 437 · **sin dueño 0 · dos dueños 0 · conteo por renglón distinto []** |
| a' | Fórmulas de la interfaz contra `model.json` | `python3 -I datos/tmp/rev-s3r1/formulas.py oraculo/model.json datos/semilla/data.json` (`Fraction`, centavos enteros, HALF_UP propio, calendario propio con la regla de §5.1, línea de tiempo propia con las reglas de §5.2) | **339 comprobaciones, 0 diferencias**: tasa y tasa_pct (binary64), R-05 «17.27» y R-06 3,629.00 desde `data.json`; 4 despensas (capital, intereses, total, 16 fechas, montos, Axel, Karen, capital repuesto); composición e intereses por persona en `despensas` y `disposicion_por_despensa`; `karen_fija`; 17 eventos × 7 campos, series A/B esperado y calculado, mínimos; 5 filas × 9 hojas de sensibilidad (con el booleano); 7 × 3 de `interes_por_monto`; hallazgo 1 de §14 confirmado (desde el total redondeado difiere solo en `[0]` y `[5]`); cruces S2: deuda viva 6,293.44, 1,258.70, 3,776.04 y efectivo P1…P6 (pagos + papel de contado, ventanas del 10 al 9); fixtures: 29-feb-2028, 0 % → 77362×3 + 77364, MSI y abonos 33333×2 + 33334, 400000 → 103474, normalización 100/200/400 → 1, 2/3, 0; estimación 6×7×15×6×2 = 7,560 |
| — | `model.json` intacto | `sha256sum oraculo/model.json` | `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` (igual al de S3-16) |
| b | Criterio → sección | Lectura de §13 contra §2 de S3.md | Las 20 tienen sección (tabla de abajo) |
| c | Suficiencia | Lectura completa por partes y cruce con C1 (`creditos-b.yaml`), 00001/00003, S2-interfaz §2.5 y FORMATO.md | Quedan 8 preguntas (P1–P8, abajo) |
| d | Coherencia con K12–K18 y K31 | Lectura | K12, K13, K14, K15, K16 y K18 coherentes; K31: 7 días como propuesta, no como valor; K43 y K44: `comparator.timeout_seconds` = 5 (00003, líneas 45 y 86); **K17: normaliza sobre «válidos» y no sobre «generados» (H4)** |
| e | 3 casos del comparador | §10.1 | Requisitos e intención de los 4 óptimos explícitos; las cifras (créditos, compras, canastas, los 4 valores de restricción y los óptimos esperados) se delegan a `S3-escenarios.md`, como permite S3-13 (H10, bajo) |

### Preguntas que quedarían al implementador (c)
- P1 (H1) ¿Qué formato de fixture, carpeta, CLI y comparador rigen: §9–§10 o `FORMATO.md`?
- P2 (H2) ¿Qué campos y llaves JSON llevan `CreditState`, `PurchaseState`, `Result`, `Scenario`, `SensitivityRow`, `DebtLine`, `PersonPeak` y `Reading`?
- P3 (H3) ¿Qué afirma cada una de las 30 subpruebas S2 de `TestCanon85` (R-37…R-60 y R-68…R-73)?
- P4 (H5) ¿Con qué valor arranca la lectura A cuando `From` ≠ `ReadingAAnchor` (en el tablero y en las ventanas del comparador)?
- P5 (H6) ¿Qué quincenas entran en la desviación estándar del objetivo 2 (con las de total 0, y con qué frontera)? ¿`over_cap` desempata, como dice el bloqueo, o no, como dice §8.6?
- P6 (H7) En el calendario unificado, ¿una tarjeta aporta sus pagos programados, sus cortes o las dos cosas (doble conteo)?
- P7 (H8) ¿Dónde se prueba en S3 que el tiempo se lee de `app_setting` y que responde 503, si el motor es puro y no puede importar `net/http`?
- P8 (H4) ¿La normalización mín-máx es sobre los escenarios válidos (§8.5) o sobre los generados (K17, RN-11, C1 `score_bp`, S3-14)?

## Criterios
| Criterio | Sev. | Sección | Estado | Evidencia propia |
|---|---|---|---|---|
| S3-01 | alta | §4, §5.1, §7, §10, §11 | Parcial | Las 55 filas S3 se reproducen con las reglas (a'); las 30 filas S2 no tienen aserción definida (H3); la fuente de R-05 y R-06 no cabe en `$model` (H1) |
| S3-02 | alta | §12, §10.2 | Cumple | (a): 437 hojas, 0 sin dueño, 0 con dos dueños, los 42 conteos iguales; 596 en total |
| S3-03 | alta | §5.2 | Cumple | (a'): 17 eventos, 28 valores y mínimos 8,094.50 y 5,000.00; omitida = B y `PendingConfirmation`; ancla de A sin definir fuera del canon (H5) |
| S3-04 | alta | §5.1 | Cumple | Regla ≥ 15 días y mín(30, último día): 16 fechas, 28-feb-2027; 29-feb-2028 en fixture |
| S3-05 | alta | §5.2 | Cumple | K15 en la fecha en que el pago queda completo; Σ capital repuesto = capital (a'); fixture `s3-reposicion-separada` |
| S3-06 | alta | §5.1 | Cumple | Orden cotizado → promoción → tabla → `ErrNoRate`; 400000 → 103474; «17.27 %» y «25.8685 %» con los decimales capturados (a') |
| S3-07 | alta | §5.1 | Cumple | 4 fronteras en `s3-promocion-0`; nada precargado; 0 % → 77362×3 + 77364 (a') |
| S3-08 | alta | §5.3 | Cumple (PROVISIONAL declarado) | 0.00 FIRME; «falta capturar» con `ErrMissingCapture`; fórmula y calendario PROVISIONALES con el bloqueo `S3-formula-tarjeta`; MSI verificado |
| S3-09 | alta | §5.4 | Cumple (PROVISIONAL declarado) | Francés con último pago = saldo + interés; tasa del periodo con el bloqueo `S3-convencion-prestamo` |
| S3-10 | media | §5.5 | Cumple | Abonos con `SplitTrunc` verificado; interés PROVISIONAL; «falta capturar» |
| S3-11 | alta | §6 | Parcial | Orden, llave y Σ por quincena definidos; el doble conteo de la tarjeta (corte + pagos programados) queda sin regla (H7) |
| S3-12 | alta | §8.1–§8.3, §8.7, §10.1 | Parcial | K16 completo (sin valores por defecto, «faltan restricciones», Override campo por campo, colchón en la línea usada, vivas y deuda en el crédito usado, tope solo marca); `comparator.timeout_seconds` = 5 en 00003; la prueba de «se lee de la clave» y del 503 no tiene dónde vivir en un motor puro (H8) |
| S3-13 | alta | §8.4, §8.6, §10.1 | Parcial | Objetivos 2 y 4 PROVISIONALES con bloqueo (válido); la serie del objetivo 2 no queda unívoca para que motor y oráculo coincidan (H6); casos sin cifras, conforme a S3-13 (H10) |
| S3-14 | alta | §8.5 | Parcial | Empate declarado (vale 1); pesos 1,1,1,1 = `optimizer.weights` de 00001 línea 1396; 100/200/400 → 1, 2/3, 0 verificado; conjunto de normalización distinto del texto de K17 (H4) |
| S3-15 | alta | §1.5 | Cumple | Sin `net/http`, `database/sql`, HiGHS ni OR-Tools en `motor_cred` |
| S3-16 | alta | §10.2 | Cumple en el contrato | Oráculo con `fractions`, no lee Go; sha256 intacto; `oraculo/verify3.py:89` sigue escribiendo en `/home/deepsec/reporte-compra-latex/model/` y no tiene dueño en §2 (H11) |
| S3-17 | alta | §1.1–§1.3 | Cumple | `Cents` int64, `big.Rat`; float solo para comparar `tasa` y `tasa_pct`; desviación con raíz entera |
| S3-18 | alta | §11 | Cumple | ≥ 85 % por paquete, `gofmt` y `go vet` |
| S3-19 | media | §4 | Cumple | Intereses de Karen HALF_UP(I × 1,191.50 / capital) = 308.22 en las 4; titular con el resto; pagos verificados (a') |
| S3-20 | media | §14 y reporte | Cumple | 3 bloqueos con las 5 partes del protocolo §2; sección Supuestos en el reporte |

## Hallazgos
### H1 · media · El contrato choca con el formato común de fixtures
Ruta: `docs/contratos/motores/S3-interfaz.md` §2, §9, §10 y §11, contra `oraculo/fixtures/FORMATO.md` v1.0 («Rige para todos los motores (S2, S3 …)») y `S2-interfaz.md` §5.
Choques: carpeta plana `oraculo/fixtures/s3-*.json` contra `oraculo/fixtures/s3/<caso>.json`; llaves `{caso, fuente (texto), provisional (texto), entradas, esperado}` contra `{formato, motor, caso, descripcion, fuente [{tipo, ruta}], provisional {bloqueo, que}, entradas {datos, consultas}, esperado {id: salida}}`; `esperado` con `{"$model": puntero JSON}` contra valores escritos más la ruta `periodos[0]…` comparada por `compara.py`; CLI `go run ./cmd/finanzas s3 dump --fixture|--canon`, salida `{caso, salida}` y código 1 por error de cálculo, contra `backend/cmd/motores s3 dump --fixtures <carpeta> --out`, salida `salida/1` y errores dentro con código 0; `cred_compara.py` y `cred_run_fixtures.py --escribir` contra `oraculo/salidas/s3.json` y `python3 -I oraculo/compara.py s3`. Además, R-05 y R-06 no están en `model.json` y el contrato solo dice «contra plan §1.2»: `FORMATO.md` ya tiene la fuente `data.json` (`datos/semilla/data.json`, `financiamiento.dos_quincenas`).
Corrección: alinear §2 y §9–§11 a `FORMATO.md`, o pedir al integrador que extienda `FORMATO.md` (por ejemplo con `$model`) y citar la versión. Fijar la fuente de R-05 y R-06 en `data.json` y retirar el SUPUESTO de §10.

### H2 · media · Tipos usados en firmas pero no declarados
Ruta: `S3-interfaz.md` §5.2, §6, §7 y §8. `CreditState` (en `Request.Credits` y `UnifiedCalendar`) no tiene ni una línea de definición. `PurchaseState` solo tiene un comentario. `Result`, `Scenario`, `SensitivityRow`, `DebtLine`, `PersonPeak`, `StatementInput` (definido en prosa) y `Reading` tampoco se declaran. Como §9 rechaza los campos desconocidos (`unknown_field`) y las entradas de las fixtures salen de estos tipos, el núcleo, el optimizador y la regresión (tres agentes) tendrían que inventar las mismas llaves por separado.
Corrección: declarar cada tipo con sus campos y su llave JSON, sobre todo `CreditState` (crédito, compras con `Ref` y `Status`, pagos registrados, cortes, contado), y dar un ejemplo JSON de `entradas` para una fixture de línea de tiempo y otra del comparador.

### H3 · media · `TestCanon85` no dice qué prueba en las 30 filas S2
Ruta: `S3-interfaz.md` §11 (renglón `TestCanon85`) y §10. S3-01 (alta) exige 85 de 85 PASS y 0 SKIP en `motor_cred`, pero R-37…R-44 (valor estable), R-45…R-52 (devengado y efectivo), R-56…R-58, R-60 y R-68…R-73 (papel) solo tienen «el insumo que S3 produce (calendario, partes, efectivo por ventana, deuda viva)». No hay función de efectivo por ventana (el contado vive en `AvailabilityTimeline`, no en `UnifiedCalendar`) ni un insumo S3 para el valor estable o el devengado.
Corrección: una tabla R-nn → función S3 → aserción exacta para las 30 filas (por ejemplo, R-45…R-50 efectivo = Σ `UnifiedCalendar` + contado en la ventana del 10 al 9, que yo verifiqué que cuadra; R-37…R-44 y el devengado: qué insumo y qué igualdad). Si alguna fila no tiene insumo S3, decir cómo se cuenta sin `t.Skip`, o abrir un bloqueo.

### H4 · media · La normalización mín-máx usa los escenarios válidos, no los «generados»
Ruta: `S3-interfaz.md` §8.5 («Sobre los escenarios **válidos**»). K17, RN-11 («entre los escenarios generados»), C1 `Scenario.score_bp` y S3-14 dicen «generados». Cambiar el conjunto cambia el mínimo y el máximo de cada objetivo, el `score_bp` y, en general, la recomendación. La decisión C («genera escenarios … que respetan las restricciones duras») permite leer «generados = válidos», pero el contrato no lo declara ni lo justifica, y el oráculo y el revisor de motores leerán S3-14 al pie de la letra.
Corrección: declarar la lectura con su justificación (cita de la decisión C) en §8.5 y en el reporte, y pedir a RN-11 y C1 que la aclaren; o abrir un bloqueo `S3-normalizacion-conjunto`; o normalizar sobre válidos + descartados.

### H5 · media · Lectura A sin regla cuando `From` no es el ancla
Ruta: `S3-interfaz.md` §5.2 («Inicio en `From`: A = límite … canon: ancla = `From`»). Si `From` es posterior a `ReadingAAnchor` (el tablero de hoy, o la ventana [fecha, último pago] del comparador en §8.3–§8.4), A volvería a 10,000.00 e ignoraría las compras y los pagos desde el ancla, cuando la lectura A es «10,000 libres hoy (ancla) y los pagos suman encima». Afecta al colchón y al objetivo 3 con la lectura A. El canon no lo detecta.
Corrección: definir A en `From` = límite − capital de las compras posteriores al ancla y anteriores o iguales a `From` + capital repuesto desde el ancla, o exigir `From` = ancla y que el comparador filtre la ventana después. Añadir una prueba con `From` ≠ ancla.

### H6 · media · Objetivo 2 sin definición unívoca, y el bloqueo y §8.6 no coinciden
Ruta: `S3-interfaz.md` §8.4 (2) y §8.6; `docs/estado/bloqueos/S3-objetivos-medida.md` (2). Falta decir qué quincenas son «las quincenas de la ventana»: ¿las de `FortnightOf` de cada día de [fecha, último pago] o solo las que tienen el día de pago dentro? ¿Entran las de total 0? Como es PROVISIONAL se acepta, pero S3-13 exige que motor y oráculo (dos agentes) coincidan exactamente en los 3 casos. Además, el bloqueo propone que `over_cap` «solo marca y desempata», y §8.6 no lo usa al desempatar.
Corrección: fijar el conjunto de quincenas (propuesta: todas las `FortnightOf` de la fecha de compra a la del último pago, inclusive, con las de 0) y alinear el desempate entre §8.6 y el bloqueo.

### H7 · media · Calendario unificado: la tarjeta puede contarse dos veces
Ruta: `S3-interfaz.md` §6 («pagos programados de todos los créditos y los cortes de tarjeta»; llave con `Source`). El cargo normal y los MSI de la tarjeta ya generan `ScheduledPayment` (§5.3), y el corte (`Source = statement`, «pago para no generar intereses») incluye esos mismos cargos. Con los dos, Σ por quincena (S3-11, que S2 y S7 consumen) cuenta dos veces el mismo dinero, y la llave con `Source` no lo marca como duplicado.
Corrección: regla explícita (por ejemplo: con corte capturado o calculado, el corte sustituye a los pagos programados de los cargos que incluye; sin corte, van los pagos programados), marcada PROVISIONAL con el bloqueo `S3-formula-tarjeta` mientras dure el calendario de la tarjeta. Prueba en `TestUnificado`.

### H8 · media · S3-12: «se lee de `app_setting`» y «503» no tienen dónde probarse
Ruta: `S3-interfaz.md` §8.7 y §11 (`TestRestriccionesTiempo`: «503 con 0»). El motor es puro y la API crea el `context`. Con eso, en `motor_cred/comparador` no se puede probar que el valor salga de `comparator.timeout_seconds` ni que la respuesta sea 503, y S3-15 prohíbe `net/http` en el código del motor.
Corrección: exportar una función pura (por ejemplo `TimeoutFromSetting(seconds int) time.Duration` y la tabla de errores a la respuesta de C1, `ErrTimeout` → 503 `unavailable`, como dato y no con `net/http`) y probarla en S3; o declarar que esa parte de S3-12 se verifica en WF-2 y pedir que se ajuste el criterio. Decir qué prueba escribe el tiempo medido de cada caso.

### H9 · baja · Detalles sin fijar
Ruta: `S3-interfaz.md` §3, §5.2, §10 y §12. (i) El evento de inicio lleva `monto`, `capital_repuesto` e `interes_no_repuesto` = 0 en `model.json`, y §5.2 no lo dice. (ii) No se dice cómo se pasa `Kind` (start, purchase, cash, payment) al `tipo` del canon (inicio, compra, pago, papel). (iii) No se dice qué muestra `Rate.Display()` con `Decimals = -1`. (iv) La fixture `s3-prestamo-frances` «100000 a 7» no da tasa, frecuencia ni primera fecha. (v) El tipo `Reading` no está declarado.
Corrección: una línea por punto.

### H10 · baja · Casos del comparador sin cifras y caso 2 frágil
Ruta: `S3-interfaz.md` §10.1. Cumple S3-13, porque se fijan en `S3-escenarios.md` antes de correr, pero el contrato no da los 4 valores de restricción ni los óptimos esperados por llave. En el caso 2, el cargo normal de 0.00 empata en intereses con la promoción en todas las tiendas, y en la tienda de mejor U/C sería también el óptimo 4. Que el óptimo 1 sea distinto queda en manos del puntaje.
Corrección: exigir que `S3-escenarios.md` traiga, por caso, los 4 valores de restricción, la llave canónica de cada óptimo esperado y la cuenta a mano que los separa, y prever cómo desempatar el caso 2.

### H11 · baja · `verify3.py` sin dueño
Ruta: `oraculo/verify3.py:89` (escribe en `/home/deepsec/reporte-compra-latex/model/verify3.json`); `S3-interfaz.md` §2 no se lo asigna a nadie, solo aparece como solicitud en el reporte. S3-16 exige corregirlo en WF-1.
Corrección: asignarlo en §2 al agente de regresión S3, con la salida en `datos/tmp/` u `oraculo/salidas/`.

### H12 · baja · Solicitudes a `reglas-negocio.md` incompletas
Ruta: `docs/estado/wf1/contrato-S3.md` (solicitud 2) y `docs/contratos/reglas-negocio.md` RN-11. S3-13 pide los objetivos 2 y 4 escritos en reglas-negocio, y la solicitud no los incluye (serie y fórmula U/C). Además, RN-11 sigue diciendo que el tiempo máximo «lo propone S3 … (K31)», cuando ya es 5 s por K43 y K44.
Corrección: ampliar la solicitud y añadir la actualización de RN-11 por K43.

## Estado de los anteriores
Ronda 1: no hay hallazgos previos.
