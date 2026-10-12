# Índice de bloqueos · Finanzas Hogar

Consolidado por el integrador de WF-0 (`integrador-wf0`) el 11-oct-2026 para la **puerta 1**. Reúne los 32 archivos de esta carpeta (31 abiertos y 1 resuelto) y 2 preguntas que dejó la revisión independiente sin archivo propio. En total son **31 preguntas para Axel**, agrupadas por la fase que necesita la respuesta. Cada archivo conserva el detalle completo; aquí van la pregunta exacta, las opciones con su consecuencia, el impacto y lo que quedó PROVISIONAL mientras tanto.

- **Ningún bloqueo esconde una cifra inventada.** La revisión r3 lo comprobó (`docs/estado/wf0/REVISION-r3.md` §4): lo provisional lleva la marca PROVISIONAL y cita su bloqueo, y las cifras de C1 salen de `oraculo/model.json` o de `datos/semilla/data.json`.
- **Una respuesta cierra dos archivos** en 4 casos: `C1-lista-compras` = `S1-lista-compras`; `C1-huecos-u1-s1` = `S1-huecos-negocio` (3 preguntas); `S1-frescura-peso` = `criterios-frescura-peso`; `S1-meta-kcal` = `criterios-meta-kcal`. Son complementarios `S1-restricciones-optimizador` y `criterios-optimizador`, y `criterios-login-bloqueo` con `C1-parametros-tecnicos`.
- **Cuándo responder.** La columna «Fase que la usa» sale del «Impacto» de cada archivo. Las de la tabla A fijan C1 y U1, que se aprueban en la puerta 1. Las de la tabla B las necesitan los motores de WF-1, que arranca después de la puerta 1. Las de la tabla C son de diseño (WF-1, se aprueban en la puerta 2). Las de la tabla D son de construcción y operación (WF-3 y WF-4).
- Cuando Axel responda, el orquestador anota la respuesta en cada archivo y regenera lo afectado (C1 v1.1 con `backend/internal/api/c1gen/build.py`; migración 00002 si cambia el esquema).

> **Estado al 11-oct-2026 (puerta 1): las 31 preguntas quedaron RESUELTAS.** Axel respondió en la sesión; la respuesta está al final de cada archivo (sección «Respuesta de Axel») y consolidada en `docs/plan/decisiones.md` §K. Las dos preguntas sin archivo ya lo tienen: `C1-decimales-tasa.md` (A12) y `S1-tope-quincena-karen.md` (A13). Lo que sigue: aplicar las respuestas (C1 v1.1, migración 00002, U1, criterios y colorimetría) en el flujo de la puerta 1.

> **Estado al 11-oct-2026 (puerta 2): 21 de los 27 bloqueos abiertos quedaron RESUELTOS y 6 siguen ABIERTOS.** Axel contestó en modo «lo recomendado + excepciones» (`docs/estado/puerta2/RESPUESTAS.md`): la opción recomendada de `PUERTA-2.md` es la respuesta, más tres respuestas explícitas (Q2-11 = B, Q2-64 = A, Q2-96 = A). Las preguntas «sin recomendación» (Q2-48, Q2-50 a Q2-52, Q2-56 a Q2-58) quedan **sin respuesta, PROVISIONAL con la opción en uso, y pasan a WF-3**; Q2-12 (la recomendada es «C o A», dos opciones) y la pregunta 4 de `S2-ventanas` (Q2-01, casos 1 y 2: el caso 1 con la regla 1A como PROVISIONAL y el caso 2 sin recomendación) también quedan abiertos. Cada archivo trae al final «Respuesta de Axel (puerta 2, 11-oct-2026)» y/o «Puerta 2 (11-oct-2026): sin respuesta», y la línea «Estado del bloqueo (puerta 2)». Detalle en la sección «Respuestas de la puerta 2 (11-oct-2026)» al final y en `docs/plan/decisiones.md` §K46 a §K52. Después de las respuestas, P2c abrió 4 bloqueos nuevos (`S2-ahorro-quincenal`, `S2-metas-vencida`, `S4-presentacion-respaldo` y `S5-carritos-granel`, los 4 ABIERTOS; ver «Bloqueos nuevos de P2c»; al cierre de P2c son 13, todos ABIERTOS) y la Pista 3 abrió 2 (`P3-contraste-chip-hover` y `P3-huecos-componentes`).

> **Estado al 12-oct-2026 (cierre de WF-3a): 40 bloqueos nuevos, todos ABIERTOS** (3 de la semilla y 37 de los demás agentes de WF-3a; tablas al final). El código corre con la opción PROVISIONAL de cada uno. Los de la columna K31 se deciden en la puerta 4; los demás, en WF-3b o WF-4. Detalle: `docs/estado/wf3/INTEGRACION-wf3a.md`.

## Resumen

**Conteo al 11-oct-2026, cierre de P1c (integrador `integrador-p1c`):** esta carpeta tiene **43 archivos de bloqueo** (sin contar este índice): **42 resueltos y 1 abierto**. **Conteo al cierre de WF-1 (integrador `integrador-wf1`, 11-oct-2026):** 69 archivos, **42 resueltos y 27 abiertos** (`S1-parametros-avisos` de P1c y 26 nuevos de WF-1). **Conteo al cierre de P2c (puerta 2, 11-oct-2026, agente `bloqueos-respuestas`):** 69 archivos, **63 resueltos y 6 abiertos** (`S2-flujo-quincenal`, `S2-ventanas`, `S5-carrito-comparativa`, `S5-carritos-detalle`, `S5-profeco-filtro` y `S5-terminos`). **Conteo al 11-oct-2026, noche (corrección `correccion-bloqueos-r2`, con la carpeta en disco):** 75 archivos, **63 resueltos y 12 abiertos**: los 69 de arriba (63 y 6), más 4 que abrió P2c después de las respuestas (`S2-ahorro-quincenal`, `S2-metas-vencida`, `S4-presentacion-respaldo`, `S5-carritos-granel`) y 2 de la Pista 3 (`P3-contraste-chip-hover`, `P3-huecos-componentes`). Las 12 abiertas: las 6 de WF-1, las 4 de P2c y las 2 de la Pista 3. Las cifras de esta carpeta cambian si otro agente abre o cierra un bloqueo; la fuente es la carpeta, no este párrafo. **Conteo al cierre de P2c (integrador `integrador-p2c`, 11-oct-2026, noche, con la carpeta en disco):** 84 archivos, **63 resueltos y 21 abiertos**: las 6 de WF-1, **13 de P2c** (los 4 de arriba más 9 que abrieron `C1-puerta2` y las correcciones de la ronda; ver «Bloqueos nuevos de P2c») y las 2 de la Pista 3. `S5-terminos` sigue ABIERTO en su archivo: la lectura del 11-oct-2026 quedó «no verificado» por desafío anti-robot y Axel contestó en la puerta 2 bis («Ignóralo», `docs/estado/puerta2/RESPUESTAS.md`), que se registra como K53 en P2d; hasta entonces no se cierra aquí. **Conteo al 12-oct-2026, madrugada (P2d; `p2c-remanentes`, con la carpeta en disco):** 85 archivos, **64 resueltos y 21 abiertos**: `k53-aplicar` abrió `S5-agente-externo` (ABIERTO, se pregunta en la puerta 4) y cerró `S5-terminos` con K53 (línea «Estado del bloqueo (puerta 2 bis): RESUELTO»). Las 21 abiertas: 5 de WF-1, 13 de P2c, 1 de P2d y 2 de la Pista 3. Si P2d cambia esos dos archivos después de esta medición, manda la carpeta.

| Origen | Archivos | Resueltos | Abiertos | Cómo se resolvieron |
|---|---|---|---|---|
| WF-0 y puerta 1 (32 de WF-0 más los 2 que se crearon en la puerta 1 para A12 y A13) | 34 | 34 | 0 | 33 con «Respuesta de Axel» de la puerta 1 (K1 a K31); `criterios-nxx` con su «Resolución» de WF-0 |
| Flujo P1 | 8 | 8 | 0 | «Respuesta de Axel (puerta 1b)», K36 a K43 |
| Flujo P1c | 1 | 1 | 0 | `S1-parametros-avisos.md`, «Respuesta de Axel (puerta 2)», Q2-11 = B, K47 (antes abierto; ver «Bloqueos nuevos de P1c» al final) |
| Flujo WF-1 | 26 | 21 | 5 | 24 con la sección «Respuesta de Axel (puerta 2)», 20 RESUELTOS (Q2-13 a Q2-71, K46); 6 abiertos (4 de ellos, `S2-flujo-quincenal`, `S2-ventanas`, `S5-carrito-comparativa` y `S5-terminos`, traen la sección porque tienen preguntas contestadas): `S2-flujo-quincenal`, `S2-ventanas`, `S5-carrito-comparativa`, `S5-carritos-detalle`, `S5-profeco-filtro` y `S5-terminos` (ver «Respuestas de la puerta 2» al final). Desde P2d (12-oct-2026), `S5-terminos` cuenta RESUELTO por K53 (su línea «Estado del bloqueo (puerta 2 bis)»): quedan 5 abiertos |
| Flujo P2c (bloqueos nuevos de la ronda, ver «Bloqueos nuevos de P2c») | 13 | 0 | 13 | 4 de `contratos-puerta2` al aplicar Q2-03, Q2-07, Q2-13, Q2-16 y Q2-18; 2 de `C1-puerta2`; 7 de las correcciones (`correccion-S1-r1` 1, `correccion-contratos-r1` 3, `correccion-motor-S5-r1` 1, `correccion-motor-S5-r2` 1 y `correccion-U4-r1` 1) |
| Flujo P2d (ver «Bloqueos nuevos de P2d») | 1 | 0 | 1 | `S5-agente-externo`, abierto por `k53-aplicar` al registrar K53; se pregunta en la puerta 4 |
| Pista 3 (ver «Bloqueos nuevos de la Pista 3») | 2 | 0 | 2 | abiertos por el agente `ds-brandbook`; no detienen la Pista 3 |
| **Total** | **85** | **64** | **21** | |

Regla de conteo de las puertas 1 y 1b (no es la que cuenta desde la puerta 2; ver el párrafo siguiente): un bloqueo está resuelto si su archivo trae la sección «Respuesta de Axel» (o, en el caso de `criterios-nxx`, su «Resolución»); sin ella sigue abierto y la marca PROVISIONAL que lo cita se conserva. Las marcas PROVISIONAL activas de los artefactos citan solo el bloqueo abierto (detalle por archivo en `docs/estado/p1c/INTEGRACION.md` §4). La tabla que sigue es la de la puerta 1 (sus 31 preguntas, todas resueltas); la columna «Lo que quedó mientras tanto» describe el estado anterior a las respuestas.

Regla de conteo desde la puerta 2: un archivo con respuesta a medias, o con una decisión tomada pero sin resultado, ya no cuenta como resuelto solo por traer «Respuesta de Axel». Cada archivo tocado en la puerta 2 termina con la línea **«Estado del bloqueo (puerta 2): RESUELTO»** o **«ABIERTO»**, y el conteo de arriba sale de esa línea. Un bloqueo con varias preguntas es RESUELTO solo si todas tienen respuesta. Hasta que `verify_c1.py` (l. 322-323) y `s1_verifica.py` (l. 242 y 575) lean esa línea (solicitud fuera de alcance en `docs/estado/p2c/correccion-bloqueos-r1.md`), las herramientas siguen la regla anterior y cuentan 67 resueltos y 18 abiertos en la carpeta de hoy (medido el 12-oct-2026 por `p2c-remanentes`): los 2 de WF-1 sin «Respuesta de Axel» (`S5-carritos-detalle` y `S5-profeco-filtro`), los 13 de P2c, `S5-agente-externo` de P2d y los 2 de la Pista 3. Con esa regla, `S2-flujo-quincenal`, `S2-ventanas`, `S5-carrito-comparativa` y `S5-terminos` cuentan como resueltos, aunque tienen preguntas abiertas. El bloqueo `S4-presentacion-respaldo` (variante de Q2-07, K46) ya existe y está en el conteo (75 al cerrar las respuestas; 84 al cierre de P2c).

| # | Tema | Archivo(s) | Fase que la usa | Lo que quedó mientras tanto |
|---|---|---|---|---|
| A1 | Lista de compras compartida, local o calculada | `C1-lista-compras` = `S1-lista-compras` | C1, S1, U1 (puerta 1) | opción A (compartida en el servidor), PROVISIONAL |
| A2 | Vigencia de un cambio de presupuesto | `C1-huecos-u1-s1` 1 = `S1-huecos-negocio` 1 | C1, S1 | cualquier vigencia sin traslape, PROVISIONAL |
| A3 | Correcciones a medias del ticket con OCR | `C1-huecos-u1-s1` 2 = `S1-huecos-negocio` 2 | C1, S1 | solo la lectura cruda, PROVISIONAL |
| A4 | Registro del chat de Telegram de Ximena | `C1-huecos-u1-s1` 3 = `S1-huecos-negocio` 3 | C1, S7 | captura manual, PROVISIONAL |
| A5 | Valores del filtro «familia» | `C1-filtro-familia` | C1, U1 | tres valores incluyentes, PROVISIONAL |
| A6 | Cinco valores técnicos de C1 (y cuatro límites de R3-02) | `C1-parametros-tecnicos` | C1, S6 | valores propuestos, PROVISIONAL |
| A7 | Definición de cinco campos calculados | `C1-campos-calculados` | C1; S3, S4, S5 | definiciones propuestas, `x-provisional` |
| A8 | Días de la semana y de la quincena | `U1-periodos-semana-quincena` | U1, S1, S2 | lunes a domingo; 1–15 y 16–fin, PROVISIONAL |
| A9 | Captura de gastos fuera de casa | `U1-captura-fuera-de-casa` | U1, C1, S6 | sin captura sin conexión, PROVISIONAL |
| A10 | Qué cuenta como «toque» | `criterios-toques` | U1, U5, C3 | regla (a) escrita en U1 §4; falta que Axel la confirme |
| A11 | Bloqueo del login y duración de la sesión | `criterios-login-bloqueo` | C1, S6 | 5 en 15 min (canon); 90 días renovables, PROVISIONAL |
| A12 | Decimales de las tasas en pantalla | sin archivo (revisión R2-06) | C1, U4 | `rate_pct` con 4 decimales |
| A13 | Tope por quincena de Karen | sin archivo (revisión R2-03) | S3, S8 | columna vacía (sin tope) |
| B1 | Redondeo de series nuevas | `S1-redondeo-series` | S2, S3 | lo del oráculo FIRME; generalización A PROVISIONAL |
| B2 | Capital que repone la parte de Karen | `S1-reposicion-parte-externa` | S3, S7 | pago completo FIRME; reparto PROVISIONAL |
| B3 | Valores y alcance de las restricciones duras | `S1-restricciones-optimizador` | S3, U4 | definiciones FIRMES; valores y alcance PROVISIONALES |
| B4 | Normalización y medidas del puntaje | `criterios-optimizador` | S3 | recomendación ponderada PROVISIONAL |
| B5 | Fórmulas de tarjeta y crédito de tienda | `criterios-credito-formulas` | S3 | casos con intereses PROVISIONALES |
| B6 | Qué es «deuda viva» en el invariante | `criterios-invariante-deuda` | S2 | identidad (b) PROVISIONAL |
| B7 | Promedio diario y semanal | `criterios-promedios` | S2 | sin cifra canon |
| B8 | Proyección de los meses 7 a 12 | `criterios-proyeccion-12m` | S2 | meses 7–12 PROVISIONALES |
| B9 | Peso de un precio viejo | `S1-frescura-peso` = `criterios-frescura-peso` | S5, S7 | 30 días FIRME; peso sin valor |
| B10 | Meta de calorías | `S1-meta-kcal` = `criterios-meta-kcal` | S4 | meta capturada; escala de actividad PROVISIONAL |
| B11 | Tickets de prueba del OCR | `criterios-ocr-tickets` | S5 | sin medir hasta definirlo |
| B12 | Paquetes con 85 % de cobertura | `criterios-cobertura` | S2 a S6 | 85 % seguro en 3 paquetes; 2 por confirmar |
| C1 | Gasto y error con el mismo rojo | `U2-gasto-error` | U2 (puerta 2) | valores decididos (25° y 20°) con regla de uso |
| C2 | ΔL entre series de gráficas | `criterios-deltal` | U2 | consecutivas como mínimo; todos los pares informativo |
| D1 | Cómo obtener Go 1.27 | `S6-version-go` | S6 (WF-3) | `go.mod` en `go 1.26` |
| D2 | Escuchar solo en la LAN o también por Tailscale | `criterios-bind` | S6, WF-4 | sin abrir Tailscale |
| D3 | Cómo llega el token del bot al servicio | `criterios-token-servicio` | S7, WF-4 | pruebas con archivo falso |
| D4 | Ocho valores operativos sin fijar | `criterios-umbrales-ausentes` | S3, S6, S7, S8, WF-4 | cada valor PROVISIONAL en su criterio |
| — | Qué significa «N-xx» | `criterios-nxx` | criterios | **resuelto** con (b): necesidad de datos `N-xxx`; pantalla `P-xx` |

## A. Fijan C1, U1 y el esquema (puerta 1)

### A1 · Lista de compras compartida, local o calculada
- **Archivos:** `C1-lista-compras.md` = `S1-lista-compras.md`.
- **Pregunta:** «¿Las marcas de comprado, los artículos a mano y la tienda o cantidad elegida por artículo se comparten entre los teléfonos de Axel y Ximena (se guardan en el servidor), viven solo en el dispositivo, o la lista es solo calculada?»
- **Opciones:** A, compartidas en el servidor: los dos ven las mismas marcas y «Terminar compra» funciona desde cualquier teléfono; ya está en el esquema 1.1 y en C1 (5 operaciones). B, solo en el dispositivo: salen 2 tablas y 3 operaciones, y si Axel marca en su teléfono Ximena no lo ve. C, lista solo calculada: sin marcas ni artículos a mano; cambia el paso 3 del flujo F-3.
- **Impacto:** `shopping_list`, `shopping_list_item`, RN-15 y las operaciones de N-150 a N-153, N-165, N-127 y N-073. Ninguna cifra canon.
- **Mientras tanto:** opción A, PROVISIONAL.

### A2 · Vigencia de un cambio de presupuesto
- **Archivos:** `C1-huecos-u1-s1.md` punto 1 = `S1-huecos-negocio.md` punto 1.
- **Pregunta:** «Cuando Axel o Ximena cambian el monto mensual de una categoría en Ajustes o en Presupuesto, ¿el nuevo monto rige solo para ese mes o para ese mes y los siguientes hasta otro cambio?» (N-087, N-275, N-278).
- **Opciones:** A, rige para ese mes y los siguientes (renglón sin fin que cierra el anterior). B, rige solo para ese mes. El esquema y C1 admiten las dos; solo cambia lo que manda la pantalla.
- **Impacto:** comportamiento de N-087, N-275 y N-278. Ninguna cifra canon.
- **Mientras tanto:** la API acepta cualquier vigencia que no se traslape (409 `budget_overlap`), PROVISIONAL en `budget.valid_to_date`, RN-22 y C1.

### A3 · Correcciones a medias del ticket con OCR
- **Archivos:** `C1-huecos-u1-s1.md` punto 2 = `S1-huecos-negocio.md` punto 2.
- **Pregunta:** «Ticket por confirmar (N-075 «Guardar para después»): ¿se guardan las correcciones que la persona ya hizo a la lectura del OCR, o solo la lectura cruda?»
- **Opciones:** A, solo la lectura cruda (`attachment.ocr_text`). B, se guardan como borrador: pide una columna nueva en S1 y una operación en C1.
- **Mientras tanto:** opción A, PROVISIONAL.

### A4 · Registro del chat de Telegram de Ximena
- **Archivos:** `C1-huecos-u1-s1.md` punto 3 = `S1-huecos-negocio.md` punto 3.
- **Pregunta:** «Registro del chat de Telegram (N-299): ¿Ximena captura su chat_id a mano en Ajustes, o el sistema lo detecta cuando ella escribe al bot MultiX?»
- **Opciones:** A, captura manual (el dato se guarda y nunca se devuelve). B, el bot lo detecta y ella confirma: trabajo de S7 y una confirmación en C1. El esquema no cambia en ninguna.
- **Mientras tanto:** opción A, PROVISIONAL.

### A5 · Valores del filtro «familia»
- **Archivo:** `C1-filtro-familia.md`.
- **Pregunta:** «(1) ¿El filtro «familia» tiene dos valores (Hogar, Karen) o tres (Hogar, Karen, Conjunto)? (2) ¿«Hogar» incluye los registros conjuntos (todo lo que tiene parte del hogar) o solo los que son 100 % del hogar?»
- **Opciones:** A, tres valores incluyentes: «todo lo de Karen» es un solo filtro (sirve para su participación y su estado de cuenta). B, tres valores exclusivos: ver todo lo de Karen pide dos valores a la vez. C, dos valores incluyentes sin «Conjunto»: U1 se corrige en N-040.
- **Impacto:** significado del parámetro `family` en 8 operaciones. Ninguna cifra canon.
- **Mientras tanto:** opción A, PROVISIONAL. U1 §0.2 dice «Hogar o Karen» y N-040 trae tres valores; U1 se alinea con la respuesta.

### A6 · Cinco valores técnicos de C1
- **Archivo:** `C1-parametros-tecnicos.md`.
- **Pregunta:** «Aprobar o cambiar cinco valores técnicos»: (1) máximo de registros por página, 200; (2) retención de la clave `Idempotency-Key`, 24 horas; (3) cookie de sesión `fh_session` con `SameSite=Strict`; (4) contraseña nueva de al menos 12 caracteres; (5) archivo de PROFECO de hasta 1 GB (los adjuntos de gastos van aparte, con 10 MB, en D4).
- **Sugerencia de la revisión r3 (R3-02):** aprobar en la misma respuesta 4 límites que hoy no llevan marca: cursor de hasta 512 caracteres, texto de búsqueda de hasta 200, usuario de 3 a 64 caracteres y clave de idempotencia de 16 a 128.
- **Opciones:** A, aprobar los valores: C1 pierde esas marcas PROVISIONAL. B, cambiar uno o varios: C1 v1.1 cambia límites y descripciones, no rutas ni esquemas.
- **Impacto:** parámetros `Limit` e `IdempotencyKey`, `cookieAuth`, `PasswordChangeRequest` e `importProfeco`. Ninguna cifra canon.
- **Mientras tanto:** opción A, PROVISIONAL.

### A7 · Definición de cinco campos calculados
- **Archivo:** `C1-campos-calculados.md`.
- **Pregunta:** «confirmar cinco definiciones»: (1) clasificación del IMC de la OMS (bajo peso < 18.5, normal 18.5 a 24.9, sobrepeso 25 a 29.9, obesidad ≥ 30); (2) precio por unidad base en centavos por 1 kg, 1 l o 1 pieza, con ROUND_HALF_UP; (3) membresía de Sam's: con qué frecuencia de compra se anualiza el ahorro y cómo se expresa el punto de equilibrio (compras al año que cubren el costo anual); (4) rendimiento de compra por peso y puntaje ponderado (dependen de B4); (5) estado de inventario: reponer = existencia estimada ≤ nivel de reposición y agotado = 0.
- **Opciones:** A, las cinco como están escritas. B, cada motor propone su definición en `reglas-negocio.md` y Axel la aprueba en la puerta 2. C, quitar de C1 los campos que no se aprueben: U1 pierde esas cifras.
- **Impacto:** `bmi_class`, `base_unit_price_cents`, `membership.*`, `value_per_peso_milli`, `score_bp` y el estado de inventario. Ninguna cifra canon.
- **Mientras tanto:** opción A con `x-provisional`; las fórmulas firmes (IMC = peso ÷ estatura², sumas de dinero) no llevan marca.

### A8 · Días de la semana y de la quincena
- **Archivo:** `U1-periodos-semana-quincena.md`.
- **Pregunta:** «1. ¿La semana va de lunes a domingo o de domingo a sábado? 2. ¿Qué días abarca cada quincena: del 1 al 15 y del 16 al último día del mes, o de fecha de pago a fecha de pago (del 15 al 29 y del 30 al 14 del mes siguiente)? En particular, ¿a qué quincena pertenece el día 31?»
- **Opciones:** semana A, lunes a domingo (ISO 8601; el resumen del domingo cierra la semana) o semana B, domingo a sábado. Quincena C, del 1 al 15 y del 16 al fin de mes (cada quincena es medio mes; el 31 va en la segunda) o quincena D, de pago a pago (una quincena cruza dos meses; pide la migración 00002 con otro CHECK de `period`).
- **Impacto:** promedios semanales, costo semanal de comidas, tendencias por semana, flujo por quincena y la tarjeta «Flujo quincenal». No cambian pantallas ni flujos.
- **Mientras tanto:** semana A y quincena C, PROVISIONAL en U1 §0.2, en el CHECK de `period` y en RN-06.

### A9 · Captura de gastos fuera de casa
- **Archivo:** `U1-captura-fuera-de-casa.md`.
- **Pregunta:** «¿qué debe hacer la app instalada en el teléfono cuando Axel o Ximena quieren registrar un gasto (o un peso) fuera de casa, sin conexión con el servidor?»
- **Opciones:** 1, sin captura fuera de casa: la app muestra «Sin conexión», deja ver lo último y deshabilita las capturas; riesgo de olvidos. 2, cola en el teléfono para gastos y peso, que se envía al volver a la red de casa. 3, cola solo para gastos. Con 2 o 3, C1 y S6 deben recibir capturas atrasadas sin duplicarlas (S1 lo resolvería con la migración 00002 y una clave de idempotencia en `expense` y `weight_entry`).
- **Impacto:** el estado «Sin conexión» y variantes de F-1 y F-6. El camino de 3 toques no cambia.
- **Mientras tanto:** opción 1, PROVISIONAL.

### A10 · Qué cuenta como «toque»
- **Archivo:** `criterios-toques.md` (criterios U1-04, U5-09 y C3-08).
- **Pregunta:** «¿Qué cuenta como «toque»? ¿Solo toques en controles (abrir captura, elegir, guardar) sin contar teclear monto y concepto? ¿Los campos con valor por defecto (fecha de hoy, categoría reciente, persona, reparto 50/50) cuentan como pasos?»
- **Opciones:** (a) toques en controles desde la pantalla de inicio hasta guardar; teclear y los valores por defecto no cuentan. (b) cada dígito cuenta (inviable en 3 toques). (c) atajo «+» para abrir y 2 toques más para guardar.
- **Mientras tanto:** U1 §4 escribió la regla (a) y con ella registrar un gasto suma 3 toques («+», categoría y «Guardar»); el revisor lo dio por demostrado. U1 y la corrección de criterios piden cerrar el bloqueo con (a): falta solo que Axel lo confirme.

### A11 · Bloqueo del login y duración de la sesión
- **Archivo:** `criterios-login-bloqueo.md` (criterios S6-06 y S6-07).
- **Pregunta:** «(1) Duración del bloqueo tras 5 intentos fallidos (¿hasta que pasen los 15 min, otro plazo?) y si el contador es por usuario, por IP o ambos. (2) Duración exacta de la sesión larga (¿60, 90, 180 días?).»
- **Opciones:** (a) ventana deslizante de 15 minutos por usuario e IP y sesión de 90 días renovable con el uso. (b) bloqueo fijo de 15 minutos tras el quinto fallo y sesión de 180 días. (c) S6 propone en C1 y Axel aprueba en la puerta 1.
- **Mientras tanto:** el límite de 5 intentos en 15 minutos es decisión de Axel; C1 describe conteo por usuario e IP con ventana deslizante y 90 días renovables, marcados PROVISIONAL.

### A12 · Decimales de las tasas en pantalla (sin archivo; revisión R2-06)
- **Pregunta sugerida por la revisión:** «¿Cuántos decimales de tasa se muestran en pantalla?»
- **Por qué importa:** la decisión C usa 17.27 % (2 quincenas) y 25.8685 % (4 quincenas). Las tasas se guardan como fracción exacta (53450/309450 = 17.2726 % y 80050/309450 = 25.8685 % con 4 decimales y ROUND_HALF_UP). Con 4 decimales, la de 2 quincenas se vería 17.2726 %; con 2 decimales, la de 4 quincenas se vería 25.87 %. Otra opción es mostrarla como se capturó.
- **Impacto:** solo la presentación (`CreditTermRate.rate_pct`); los cálculos usan la fracción. La revisión también pide documentar que el sufijo «Z» de las fechas y horas solo se acepta en entradas.
- **Mientras tanto:** `rate_pct` con 4 decimales.

### A13 · Tope por quincena de Karen (sin archivo; revisión R2-03)
- **Pregunta sugerida por la revisión:** «¿El tope por quincena de Karen del objetivo 2 es 374.93 (nota «cap» del oráculo) o lo fija Axel?»
- **Por qué importa:** `oraculo/model.json` dice «fuera del cap de 374.93» en dos notas de `papel`; S1 trata 374.93 como la parte de Karen por pago, no como un tope. El tope solo alimenta el objetivo 2 del comparador («pago por quincena parejo y bajo tope»); no es restricción dura (decisión C).
- **Mientras tanto:** `person.fortnight_payment_cap_cents` vacía (sin tope); ninguna cifra entró al esquema ni a C1.

## B. Las necesitan los motores de WF-1

### B1 · Redondeo de series nuevas
- **Archivo:** `S1-redondeo-series.md`.
- **Pregunta:** «Para series nuevas que no están en el oráculo, ¿qué base de redondeo se usa en (1) las partes por persona de cada pago cuando el total de una persona no se divide exacto entre el número de pagos, y (2) el devengo mensual cuando el costo tiene centavos impares? ¿Y el «pago de referencia» se muestra en pantallas o solo la serie real?»
- **Opciones:** A (la recomienda S1), conservar los dos modos del oráculo: montos con base truncada, capital con base HALF_UP, partes de personas externas con base HALF_UP y el titular toma el resto, devengo con base HALF_UP y el último ajusta; el pago de referencia solo en cotizaciones. B, todo HALF_UP: rompe 1,258.68/1,258.70 y 883.75/883.77. C, todo truncado: rompe 773.63/773.61 y las series A/B. D, como A pero con las partes también truncadas.
- **Impacto:** S3 (series de compras nuevas) y S2 (devengo). El esquema no cambia.
- **Mientras tanto:** lo que fija el oráculo es FIRME (RN-07); la generalización queda PROVISIONAL con la opción A como propuesta.

### B2 · Capital que repone la parte de Karen
- **Archivo:** `S1-reposicion-parte-externa.md`.
- **Pregunta:** «Cuando Karen paga su parte de un pago (374.93) por separado, antes o después que Axel, ¿cuánto capital repone ese pago en el disponible?»
- **Opciones:** A, proporcional al capital del pago (con Karen, 297.88 de 773.63, y Axel repone 475.75). B, solo al completarse el pago (lo más conservador). C, primero su capital propio (coincide con A en la despensa 1). D, lo que diga Kueski cuando Axel confirme la reposición (va junto con el interruptor A/B).
- **Impacto:** disponible A/B entre pagos parciales y estado de cuenta de Karen. El esquema no cambia.
- **Mientras tanto:** la regla por pago completo es FIRME (sale del oráculo); el reparto entre partes, PROVISIONAL.

### B3 · Valores y alcance de las restricciones duras
- **Archivo:** `S1-restricciones-optimizador.md`.
- **Pregunta:** «(1) El valor de cada restricción y el «± días» de los escenarios. (2) ¿Son del hogar (una sola para todo), por crédito, o se fijan en cada compra planeada? (3) ¿El colchón se mide en la línea que se usa o en la suma del disponible de todas las líneas? (4) ¿Las compras vivas se cuentan en todos los créditos o solo en el que se usa? (5) ¿La deuda en cero es de todos los créditos o del que se usa?»
- **Opciones:** A, globales con valor por defecto en Ajustes que se cambia en cada comparación; colchón en la línea usada y compras vivas y deuda en el crédito usado (sin cambiar el esquema). B, por crédito (migración 00002). C, por compra planeada (migración 00002). D, globales sobre todo el hogar (lo más estricto; sin cambiar el esquema).
- **Impacto:** S3 y Ajustes «Créditos y condiciones». Las claves `optimizer.*` quedan sin valor.
- **Mientras tanto:** definiciones FIRMES en RN-11; valores y alcance PROVISIONALES.

### B4 · Normalización y medidas del puntaje
- **Archivo:** `criterios-optimizador.md` (criterios S3-13 y S3-14).
- **Pregunta:** «(1) ¿Cómo se normaliza cada objetivo antes de ponderar (mín-máx entre los escenarios generados, relativo al mejor, rango fijo)? (2) ¿Qué mide «pago por quincena parejo y bajo tope»: desviación estándar de los pagos, diferencia máx-mín, pico contra el tope? (3) ¿Qué mide «comprar más por peso»: unidades base por peso gastado de la canasta, usando qué precios?»
- **Opciones:** (a) mín-máx entre escenarios; parejo = desviación estándar de pagos; por peso = unidades base por peso. (b) relativo al mejor escenario de cada objetivo. (c) el agente propone las tres definiciones y Axel las aprueba en la puerta 2.
- **Mientras tanto:** los 4 óptimos por objetivo no dependen de esto; la recomendación ponderada queda PROVISIONAL.

### B5 · Fórmulas de tarjeta y crédito de tienda
- **Archivo:** `criterios-credito-formulas.md` (criterios S3-08 y S3-10).
- **Pregunta:** «Fórmulas de: pago mínimo de tarjeta (¿% del saldo, mínimo fijo, ambos?), intereses de tarjeta cuando no se paga el total (tasa anual, base diaria, IVA), y comisión o interés del crédito de tienda con abonos.» Los datos reales de las tarjetas y tiendas de Axel no están en el repositorio.
- **Opciones:** (a) Axel captura por producto la tasa y el método de su estado de cuenta. (b) Axel entrega un estado de cuenta de cada producto para fijar las fórmulas con ejemplo. (c) solo el caso sin intereses y meses sin intereses; el resto, captura manual de pagos.
- **Mientras tanto:** los casos con intereses quedan PROVISIONALES; no se inventa tasa ni mínimo.

### B6 · Qué es «deuda viva» en el invariante
- **Archivo:** `criterios-invariante-deuda.md` (criterio S2-07).
- **Pregunta:** «¿Qué significa «deuda viva» en ese invariante?» (plan §3.4, fila S2: «devengado acumulado − efectivo acumulado = deuda viva»). Con el oráculo, la forma literal solo se cumple en P2, P4 y P6; siempre se cumple deuda viva = (devengado acumulado − efectivo acumulado) + saldo no devengado.
- **Opciones:** (a) dos nombres canon: «deuda viva devengada» (1,258.70 al final) y «deuda viva total» (6,293.44). (b) invariante algebraico con el saldo no devengado (10-abr-2027: 6,293.44 = 1,258.70 + 5,034.74); hay que reescribir el criterio del plan. (c) cambiar el oráculo (no se recomienda: invalida cifras canon).
- **Mientras tanto:** S2-07 exige la identidad (b), PROVISIONAL. Ninguna cifra cambió.

### B7 · Promedio diario y semanal
- **Archivo:** `criterios-promedios.md` (criterio S2-08).
- **Pregunta:** «¿Cómo se calculan el promedio diario y el semanal?» El oráculo solo trae el promedio mensual (efectivo 2,183.96; devengado 2,393.75).
- **Opciones:** (a) total del horizonte ÷ días reales (182) y ÷ (días/7) para el semanal. (b) promedio mensual ÷ 30 y × 12/52. (c) promedio mensual ÷ días del mes calendario.
- **Mientras tanto:** las cifras diarias y semanales no tienen valor canon y quedan PROVISIONALES.

### B8 · Proyección de los meses 7 a 12
- **Archivo:** `criterios-proyeccion-12m.md` (criterio S2-09).
- **Pregunta:** «¿Qué supuestos rigen los meses 7 a 12 (10-abr-2027 a 9-oct-2027)? ¿Se repite la cadencia bimestral de despensa a crédito (compras 10-jun, 10-ago…) y el papel alterno de $199? ¿Solo se proyectan gastos fijos, ingresos y créditos ya contratados?»
- **Opciones:** (a) repetir la cadencia bimestral de la despensa y el papel alterno. (b) solo compromisos contratados y gastos fijos e ingresos del CSV de Axel. (c) configurable en Ajustes, con (b) por defecto.
- **Mientras tanto:** los meses 1 a 6 coinciden con el oráculo; los meses 7 a 12 llevan la marca PROVISIONAL.

### B9 · Peso de un precio viejo
- **Archivos:** `S1-frescura-peso.md` = `criterios-frescura-peso.md` (criterio S5-04).
- **Pregunta:** «(1) ¿Cuánto menos pesa un precio viejo en la comparativa de tiendas y canastas? (2) Confirmar que «a los 30 días» significa una edad de 30 días o más, por ejemplo que un precio del 10-oct-2026 ya es viejo el 9-nov-2026.» Y, del archivo de criterios: «¿Cómo se ponderan entre sí las fuentes (ticket manual, OCR confirmado, PROFECO)?»
- **Opciones:** A, peso fijo de 0.5. B, el precio viejo solo se usa si no hay uno fresco, y se marca. C, decaimiento lineal hasta 0 a los 90 días. D, sin peso, solo la marca (contradice «pesa menos»). Para las fuentes, la opción (a) de criterios: todas con el mismo peso.
- **Mientras tanto:** el umbral de 30 días es FIRME (`price.stale_after_days = 30`, lectura «30 o más», PROVISIONAL para criterios); el peso queda sin valor.

### B10 · Meta de calorías
- **Archivos:** `S1-meta-kcal.md` = `criterios-meta-kcal.md` (criterio S4-07).
- **Pregunta:** «¿La meta de calorías la captura cada persona, o el sistema la sugiere con una fórmula? Si la sugiere: ¿qué fórmula, qué escala de actividad y si se ajusta según el objetivo (bajar, mantener o subir de peso)?»
- **Opciones:** A, solo captura manual. B, sugerencia con Mifflin-St Jeor × factor de actividad (1.2 a 1.9) que la persona confirma. C, Harris-Benedict revisada con los mismos factores. D, B más un ajuste por objetivo (columna nueva).
- **Mientras tanto:** el esquema guarda los insumos y la meta capturada; la escala de actividad de 5 niveles es PROVISIONAL.

### B11 · Tickets de prueba del OCR
- **Archivo:** `criterios-ocr-tickets.md` (criterio S5-06).
- **Pregunta:** «(1) ¿De dónde salen los 5 tickets: reales que entrega Axel (tickets de Bodega Aurrera/Walmart de la compra del 10-oct u otros) o sintéticos generados? (2) ¿Qué es un «renglón correcto» (descripción + cantidad + importe exactos, solo importe, descripción con tolerancia)?»
- **Opciones:** (a) 5 tickets reales fotografiados por Axel. (b) 5 sintéticos (prueba menos fiel). (c) 3 reales y 2 sintéticos.
- **Mientras tanto:** el ≥ 80 % de renglones correctos no se puede declarar cumplido; con sintéticos sería «prueba limitada». Con (a) o (c), Axel fotografía los tickets.

### B12 · Paquetes con 85 % de cobertura de pruebas
- **Archivo:** `criterios-cobertura.md` (criterios S4-15, S5-13 y S6-02).
- **Pregunta:** «¿Qué paquetes son «motores»? … ¿El 85 % aplica a recetas y precios, y a api/auth/jobs/reportes?»
- **Opciones:** (a) dominio, motor_fin, motor_cred, recetas y precios. (b) solo dominio, motor_fin y motor_cred; los demás se miden sin umbral. (c) todo `internal/`.
- **Mientras tanto:** 85 % exigido en dominio, motor_fin y motor_cred; recetas y precios, sujetos a esta respuesta.

## C. Diseño (WF-1, se aprueba en la puerta 2)

### C1 · Gasto y error con el mismo rojo
- **Archivo:** `U2-gasto-error.md`.
- **Pregunta:** «¿Aceptas que gasto y error compartan prácticamente el mismo rojo, distinguidos solo por contexto, signo e icono, o prefieres separarlos con alguna de las opciones 2 a 4?»
- **Opciones:** 1, mantener gasto 25° y error 20° (lo entregado; ΔE ≈ 1, el color no ayuda a distinguirlos). 2, mover gasto a 40°–45°: hay que rehacer Karen y Bodega Aurrera en la paleta. 3, gasto sin color de rol (texto neutro con «−»). 4, gasto 25° con croma bajo (rompe la regla de croma uniforme).
- **Impacto:** solo los tokens `color.role.expense.*` (con la opción 2 también `color.person.karen.*` y `color.store.bodega-aurrera.*`).
- **Mientras tanto:** los valores decididos con la regla de uso de la opción 1. Nada PROVISIONAL.

### C2 · ΔL entre series de gráficas
- **Archivo:** `criterios-deltal.md` (criterio U2-09).
- **Pregunta:** «¿«Entre series» significa entre series consecutivas en el orden de uso o entre todos los pares (28 pares con 8 series)?»
- **Opciones:** (a) consecutivas, con forma, etiqueta y orden fijo como refuerzo. (b) todos los pares (obliga a combinar patrones de relleno con luminosidad). (c) pares que aparecen juntos en una gráfica (máximo 5 series).
- **Mientras tanto:** (a) es el mínimo exigido y U2a lo cumple (ΔL mínimo 0.162 en claro y 0.198 en oscuro); (b) se reporta como información.

## D. Construcción y operación (WF-3 y WF-4)

### D1 · Cómo obtener Go 1.27
- **Archivo:** `S6-version-go.md` (criterio S6-01).
- **Pregunta:** «¿Cómo se obtiene Go 1.27 en la máquina de casa?» Hoy hay Go 1.26.5 y pacman solo ofrece 1.26.5 sin actualizar todo el sistema.
- **Opciones:** A, Axel corre `sudo pacman -Syu` (actualización completa, con reinicio si toca el kernel) y `go.mod` pasa a 1.27. B, Go descarga su propia cadena 1.27 (`GOTOOLCHAIN=auto`) en la caché del usuario; la primera compilación necesita internet. C, seguir con Go 1.26 hasta la siguiente actualización normal; S6-01 no se puede declarar cumplido.
- **Mientras tanto:** `backend/go.mod` declara `go 1.26`; `go build` y `go vet` pasan con 1.26.5.

### D2 · Escuchar solo en la LAN o también por Tailscale
- **Archivo:** `criterios-bind.md` (criterio S6-14).
- **Pregunta:** «¿El servicio debe escuchar solo en la LAN (192.168.50.x) o también por Tailscale (acceso remoto de Axel y Ximena fuera de casa)?»
- **Opciones:** (a) solo LAN, lectura literal de «solo red local». (b) LAN y Tailscale: más cómodo, más superficie. (c) LAN ahora y Tailscale como decisión posterior.
- **Impacto:** dirección de escucha en S6 y WF-4 y nombres del certificado de mkcert. Se relaciona con A9: por Tailscale la app alcanzaría el servidor fuera de casa.
- **Mientras tanto:** no se abre Tailscale; el esqueleto escucha solo en 127.0.0.1.

### D3 · Cómo llega el token del bot al servicio
- **Archivo:** `criterios-token-servicio.md` (criterios S7-06, S7-07 y WF4-01).
- **Pregunta:** «¿Cómo llega el token al servicio si el archivo está en el home de otro usuario (`deepsec`) y el servicio corre como `finanzas`?»
- **Opciones:** (a) el servicio corre como `deepsec`, sin usuario propio. (b) permiso de lectura (ACL) del archivo para `finanzas`. (c) archivo de entorno de systemd fuera del repositorio, con permisos 600, creado por Axel.
- **Mientras tanto:** las pruebas usan un archivo falso y el token jamás se imprime ni entra al repositorio.

### D4 · Ocho valores operativos sin fijar
- **Archivo:** `criterios-umbrales-ausentes.md` (criterios S3-12, S6-09, S6-10, S7-04, S7-08, S7-10, S8-02, WF4-04 y WF4-08).
- **Pregunta:** «(1) hora exacta y contenido del resumen diario (noche) y del semanal (domingo) y del PDF mensual (día 1); (2) anticipación de la alerta «precios viejos antes del día de despensa» y definición de «plan de comidas sin armar» (¿cualquier celda vacía de la semana siguiente?); (3) hora y retención (cuántas versiones) del respaldo diario; (4) lista exacta de cabeceras de seguridad; (5) límites de tamaño y tipo de adjuntos; (6) tamaño máximo del PDF (más allá del límite de la API de Telegram); (7) tiempo máximo aceptable de la búsqueda exhaustiva del optimizador; (8) permisos de archivos de la base y respaldos.»
- **Opciones:** (a) cada agente propone el valor con justificación y Axel lo aprueba en la puerta que toque. (b) Axel fija los 8 valores ahora. (c) valores estándar documentados (cabeceras OWASP, 90 días de retención, etc.) marcados PROVISIONAL.
- **Mientras tanto:** cada valor queda PROVISIONAL en su criterio. El punto 7 lo necesita S3 en WF-1.

## Resuelto

### criterios-nxx · Qué significa «N-xx»
- Resuelto el 10-oct-2026 en la ronda 1 de correcciones con la opción (b): «N-xx» es una necesidad de datos (`N-xxx` en U1 §3) y la pantalla es `P-xx` (`P-26.xx` para las secciones de Ajustes). No cambió cifras, esquema ni operaciones. Solo se reabre si Axel prefiere otra lectura (opción c).

## Respuestas de la puerta 1 (11-oct-2026)

| # | Respuesta de Axel (resumen) |
|---|---|
| A1 | A: lista compartida en el servidor |
| A2 | Configurable por cambio («solo este mes» / «de este mes en adelante»); preselección «en adelante» |
| A3 | B: se guardan las correcciones del OCR como borrador |
| A4 | B: el bot detecta el chat y Ximena confirma |
| A5 | A: tres valores incluyentes |
| A6 | Aprobados los 5 valores y los 4 límites de R3-02 |
| A7 | A: las cinco definiciones como están |
| A8 | Lunes a domingo; quincenas 1–15 y 16–fin de mes |
| A9 | 2: cola en el teléfono para gastos y peso |
| A10 | (a) toques en controles; 3 toques |
| A11 | (a) ventana deslizante 15 min por usuario e IP; sesión 90 días renovable |
| A12 | Como se capturó (17.27 % y 25.8685 %) |
| A13 | Sin tope fijo; Axel lo captura por persona |
| B1 | A: los dos modos del oráculo |
| B2 | B: capital repuesto solo al completarse el pago |
| B3 | A: globales en Ajustes, ajustables por comparación; línea y crédito usados |
| B4 | (a) mín-máx; desviación estándar; unidades por peso |
| B5 | (a) Axel captura tasa y método por producto |
| B6 | (b) identidad con el saldo no devengado |
| B7 | (a) total ÷ días reales; semanal = diario × 7 |
| B8 | (c) configurable; por defecto solo compromisos y CSV |
| B9 | C: decaimiento lineal 30→90 días; PROFECO solo referencia |
| B10 | D: Mifflin-St Jeor × actividad con ajuste por objetivo |
| B11 | (a) 5 tickets reales de Axel; renglón = descripción + cantidad + importe |
| B12 | (a) dominio, motores, recetas y precios |
| C1 | 2: gasto a terracota 40–45°; mover Karen y Bodega Aurrera |
| C2 | (a) consecutivas con refuerzo |
| D1 | pacman -Sy (sin u) + go del repositorio: Go 1.27.2 instalado |
| D2 | (a) solo LAN |
| D3 | (b) ACL de lectura de openclaw.json para el servicio |
| D4 | (a) cada agente propone; Axel aprueba en su puerta |

## Bloqueos nuevos de P1 (11-oct-2026)
Los dejaron los agentes del flujo «P1 · Aplicar respuestas de la puerta 1» por ambigüedades que §K no resuelve. Los 8 quedaron **RESUELTOS** el 11-oct-2026 (puerta 1b, respuestas K36–K43 en `docs/plan/decisiones.md`; cada archivo lleva su «Respuesta de Axel (puerta 1b)»); el flujo P1c quita las marcas PROVISIONAL que los citan. Ninguno cambia cifras del oráculo ni el esquema; lo que depende de ellos lleva la marca PROVISIONAL y cita el archivo. Registrado por el integrador de P1 (`integrador-p1`).

| Archivo | Fase que lo dejó | Pregunta para Axel | Qué queda PROVISIONAL | Fase que lo usa |
|---|---|---|---|---|
| `C1-valores-k31.md` | C1 v1.1 | Aprobar o cambiar 4 propuestas de valores operativos (K31) | 4 descripciones de `openapi.yaml` (alta de adjuntos, `compareScenarios`, `DashboardWeekMeals`, `NotificationUpdate`) y el estado de N-023, N-064, N-226, N-297 y N-300 en la matriz de `C1-datos.md` §10 | La puerta que fija K31 para cada valor (optimizador en la puerta 2; los demás en la 4) |
| `S1-ajuste-objetivo-kcal.md` | S1 esquema 1.2 | Cuánto se ajusta la meta sugerida de calorías al bajar o subir de peso (K23 no da la cifra) | El ajuste por objetivo en RN-17 y en C1; «mantener» no ajusta | S4 (WF-1) y la pantalla Salud |
| `S1-fuentes-precio.md` | S1 esquema 1.2 | Si compiten en la comparativa los precios de la lectura asistida y los capturados a mano sin ticket, y con qué peso | Las fuentes que compiten en RN-12 | S5 (WF-1) |
| `U1-lista-vigente-domingo.md` | U1 1.1 | Los domingos, qué semana abren por defecto el plan de comidas y la lista de compras | La semana por omisión de P-10, P-13 y P-14 (U1 §0.2 y §7); C1 debe marcarla igual (hallazgo alto H-r3-02 de `docs/estado/p1/REVISION-C1-r3.md`) | U1 y C1 (antes de WF-3) |
| `U2a-rosas-ximena-bodega.md` | Colorimetría 1.1 | Si se acepta que Karen y Bodega Aurrera queden cerca de Ximena solo con el refuerzo declarado, u otra salida | 4 tokens: `color.person.karen.light/.dark` y `color.store.bodega-aurrera.light/.dark` | U2 `DESIGN.md` (WF-1, puerta 2) |
| `U2a-walmart-oscuro.md` | Colorimetría 1.1 | Si se acepta bajar Walmart en oscuro de L 0.720 a 0.716 (`#7BA4F0` → `#7AA3EE`, ΔE 0.4) | `color.store.walmart.dark` | U2 `DESIGN.md` (WF-1, puerta 2) |
| `criterios-frescura-extremos.md` | Criterios 1.1 | Cuánto pesa un precio de 30 días y cómo cambia día a día entre el 30 y el 90 (K22 en los extremos) | El peso de los precios de 30 a 89 días en la comparativa | S5 (WF-1) |
| `criterios-promedios-redondeo.md` | Criterios 1.1 | Si el promedio semanal multiplica por 7 el diario sin redondear o ya redondeado a centavos (K20) | Solo el promedio semanal (unos centavos por periodo) | S2 (WF-1) |

Detalle de dónde se cita cada uno y cuántas marcas PROVISIONAL hay por archivo: `docs/estado/p1/INTEGRACION.md`.

## Respuestas de la puerta 1b (11-oct-2026)
| Bloqueo | Decisión de Axel | K |
|---|---|---|
| `S1-ajuste-objetivo-kcal.md` | Ajuste configurable; inicia en 0 para bajar y subir; por persona con el valor de la casa como inicial; «mantener» no ajusta | K36 |
| `S1-fuentes-precio.md` | C: `manual` compite igual que el ticket; `site_assisted` solo referencia hasta veredicto de términos | K37 |
| `U1-lista-vigente-domingo.md` | 1: los domingos abren la semana siguiente | K38 |
| `criterios-promedios-redondeo.md` | (a): diario exacto × 7, un solo redondeo al final | K39 |
| `criterios-frescura-extremos.md` | (b): peso = (90 − días) / 61; día 30 = 0.9836 | K40 |
| `U2a-rosas-ximena-bodega.md` | 1: aceptar con refuerzo declarado; DESIGN.md prohíbe persona y tienda solo por color | K41 |
| `U2a-walmart-oscuro.md` | 1: Walmart oscuro en L 0.716 | K42 |
| `C1-valores-k31.md` | A: aprobadas las cuatro propuestas | K43 |
| (regla general) | Todo parámetro numérico operativo nace con su valor decidido y es configurable en Ajustes | K44 |
| (proceso) | La Pista 3 (Claude Design) la ejecuta Axel; el orquestador entrega `docs/design/entrada/` y revisa `docs/design/salida/` | K35 |

## Bloqueos nuevos de P1c (11-oct-2026)
> **Actualización de la puerta 2 (11-oct-2026):** `S1-parametros-avisos.md` quedó RESUELTO con Q2-11 = B (K47): los cuatro parámetros viven en `app_setting`.
Los dejaron los agentes del flujo «P1c · Aplicar respuestas de la puerta 1b» por un choque que §K no resuelve. Registrado por el integrador de P1c (`integrador-p1c`). No cambia cifras del oráculo; lo que depende de él lleva la marca PROVISIONAL y cita el archivo.

| Archivo | Fase que lo dejó | Pregunta para Axel | Qué queda PROVISIONAL | Fase que lo usa |
|---|---|---|---|---|
| `S1-parametros-avisos.md` | S1 esquema 1.3 (reducido por `correccion-S1-r1` y ampliado por `correccion-S1-r2`) | Si `send_time`, `advance_days`, `weekday` (7) y `day_of_month` (1) viven en `notification.params_json` (opción A, K43 al pie de la letra), en claves de `app_setting` (opción B, K44, como ya quedaron los umbrales 80 y 100 % y los días antes del pago [1, 0]) o mixto (opción C) | Solo el **lugar** de esos cuatro parámetros (el valor de `weekday` y `day_of_month` es FIRME por la decisión F; `send_time` y `advance_days` no tienen valor: lo propone S7, K31). Marcas en `00003_puerta1b.sql`, `diccionario-datos.md`, `reglas-negocio.md` (RN-29), `openapi.yaml`, `C1-datos.md` y en los criterios C1-22, U4-12 y S7 | Puerta 2; C1, U4 y S7 |

Recomendaciones de los agentes, que no coinciden: S1 recomienda la **opción B** (`docs/estado/p1c/correccion-S1-r2.md`, una sola fuente por valor, como los umbrales); la corrección de criterios recomendó la **opción A** (`docs/estado/p1c/correccion-criterios-r1.md`, sin cambio de esquema ni contrato). Ojo: los criterios 1.2 todavía describen los umbrales y los días antes del pago como si vivieran en `params_json` (hallazgo A, alto, de `docs/estado/p1c/REVISION-criterios-r3.md`); el esquema 1.3 ya los tiene en `app_setting`.

## Bloqueos nuevos de WF-1 (11-oct-2026, para la puerta 2)
> **Actualización de la puerta 2 (11-oct-2026):** 20 de los 26 quedaron RESUELTOS; siguen abiertos `S2-flujo-quincenal`, `S2-ventanas`, `S5-carrito-comparativa`, `S5-carritos-detalle`, `S5-profeco-filtro` y `S5-terminos`. El texto de abajo es el de la consolidación de WF-1 y no se reescribe.

Consolidado por `integrador-wf1`. WF-1 dejó **26 archivos nuevos, todos abiertos**; `S1-parametros-avisos` (P1c) también sigue abierto: nadie lo decidió en WF-1 y lo que depende de él (U4 «Notificaciones» y «Reportes y frecuencias») lleva la marca PROVISIONAL con sus dos variantes. La columna «Fase» es la fase que dejó el archivo y la que lo usa. Las recomendaciones de los agentes para cada bloqueo están en `docs/estado/wf1/INTEGRACION.md` §7.2, y los artefactos que se escalan a Axel con sus hallazgos abiertos, en §4.

| Archivo | Fase | Pregunta (primera línea de «Qué falta») |
|---|---|---|
| `S2-flujo-quincenal.md` | S2 | (1) método del gasto variable estimado de una quincena futura; (2) cómo se reparte el ahorro mensual planificado (`planned_monthly_cents`) entre las dos quincenas; (3) cuándo una ocurrencia planificada (sueldo, renta) se da… |
| `S2-metas-resumen.md` | S2 | ¿el resumen del mes cuenta solo las metas y fondos activos o también los inactivos que aún tienen saldo? M-07 dice «ahorrado = Σ aportaciones − Σ retiros (fondo y metas)» sin decir «activos», y en la misma fila dice… |
| `S2-metas.md` | S2 | (1) cómo se cuentan los «meses que faltan» de M-16 («aportación mensual necesaria = (monto meta − ahorrado) ÷ meses que faltan»); (2) qué porcentaje de avance se muestra: RN-18 dice «con un decimal» (50.0 %) y M-07/M-16 dicen… |
| `S2-presupuesto-total.md` | S2 | (1) si existe un renglón del total, ¿el plan total es ese monto o la suma de las categorías? (2) Con un plan de 0.00 en una categoría, ¿qué porcentaje y qué banderas se muestran? |
| `S2-proyeccion-interruptor.md` | S2 | (1) monto de cada despensa proyectada: ¿capital de la despensa base 4,000.00 (despensas 2–4) con la tasa vigente de la tabla por plazo, o el de la última despensa comprada, o una partida planificada que Axel captura? (2)… |
| `S2-ventanas.md` | S2 | (pregunta para Axel y para el dueño de la interfaz): 1. Con ventanas, ¿qué pasa con un costo o un pago anterior a la primera ventana pedida? §3.2 dice que `cumulative_diff_cents` va «desde el primer evento de las entradas… |
| `S3-convencion-prestamo.md` | S3 | ¿Cómo se calcula la tasa por periodo de un préstamo a plazo fijo y cómo entra el IVA? |
| `S3-formula-tarjeta.md` | S3 | (1) Intereses de tarjeta cuando no se paga el total: ¿sobre qué saldo (saldo promedio diario del periodo, saldo al corte, saldo insoluto después del pago) y con qué días (días reales del periodo / 360, / 365, o tasa anual / 12… |
| `S3-objetivos-medida.md` | S3 | (1) Canasta y tienda (objetivo 4): ¿la tienda es una dimensión más del escenario (cada tienda con precios vigentes) y qué es «la canasta» de una compra planeada? (a) la lista de compras de la semana (P-13) con los precios… |
| `S3-origen-cargo-normal.md` | S3 | ¿qué `interest_source` lleva un cargo normal de tarjeta? |
| `S3-serie-sin-intereses.md` | S3 | con intereses 0, ¿qué capital lleva cada pago? K14 y RN-07 fijan montos con base truncada y capital con base HALF_UP, y el último pago ajusta. Para 3,094.50 a 4 quincenas eso da montos 773.62 ×3 + 773.64 y capital 773.63 ×3 +… |
| `S4-asignacion-mercado.md` | S4 | (cuatro preguntas para Axel): 1. Si varios mercados surten el grupo del producto, ¿a cuál va el artículo de la lista? |
| `S4-confirmacion-empates.md` | S4 | ¿quién confirma los empates de `nutricion/tabla.csv` y las ligas de `profeco/vinculos.csv`, y en qué momento? |
| `S4-reglas-lista.md` | S4 | (tres preguntas para Axel): 1. **Sobrante previsto** (`MealPlanWeek.leftovers`, RN-15): ¿el plan cocina exactamente las porciones planeadas o la receta completa (múltiplos de sus 2 porciones base)? |
| `S5-carrito-comparativa.md` | S5 | (preguntas para Axel): 1. Con `basket = expense`, ¿contra qué se compara la compra registrada en las otras tiendas? |
| `S5-carritos-detalle.md` | S5 | tres preguntas para Axel (o para el dueño del contrato): 1. Ganador de un sustituto cuando una tienda tiene varias presentaciones del producto: ¿precio por unidad base de **lo que esa tienda cobró** por el producto (importe… |
| `S5-profeco-almacen.md` | S5 | (tres preguntas para Axel): 1. ¿Dónde viven los registros crudos de PROFECO (millones de renglones, la mayoría sin producto del catálogo)? |
| `S5-profeco-filtro.md` | S5 | tres preguntas para Axel (o para el dueño del contrato): 1. ¿El `municipality` de una tienda que **no** es sucursal (una cadena o un mercado) entra al filtro de municipios? |
| `S5-reglas-comparativa.md` | S5 | (cuatro preguntas para Axel): 1. ¿Cómo «pesa menos» un precio viejo al elegir la tienda más barata? |
| `S5-terminos.md` | S5 | Axel, ¿cómo se obtiene el veredicto de términos de los 3 sitios (walmart.com.mx, bodegaaurrera.com.mx y sams.com.mx)? **RESUELTO por K53** (puerta 2 bis, P2d): la aplicación no lee sitios de tiendas. |
| `U3-t2-acumulados.md` | U3 | ¿C1 agrega a `AccrualVsCash.months[]` el devengado acumulado y el efectivo acumulado con el mismo origen que `cumulative_difference_cents` (opción A), o T-2 se dibuja relativo a la ventana con las sumas de lo visible (opción B)? |
| `U4-dia-mes-reporte.md` | U4 | si Axel cambia el día del mes del reporte mensual a 29, 30 o 31, ¿qué pasa en los meses que no tienen ese día (febrero, abril, junio, septiembre, noviembre)? |
| `U4-iva-credito.md` | U4 | C1 v1.2 y el esquema 1.3 guardan «si incluye IVA» en `credit.vat_on_interest_bp`, donde NULL significa «la tasa ya incluye IVA»; por eso `Credit.missing_terms` solo puede listar 4 datos (`annual_rate`, `rate_base`,… |
| `U4-reparto-retroactivo.md` | U4 | cuando Axel cambia el porcentaje de una regla de reparto reutilizable (p. ej. «Conjunto 50/50»), ¿el cambio vale para los gastos ya capturados o solo para los nuevos? |
| `U4-vista-previa-formulas.md` | U4 | (dos preguntas vigentes para Axel; la 2 se retiró): 1. **Comisiones de un crédito** (anual, por pago tardío y por apertura; `credit.annual_fee_cents`, `late_fee_cents`, `opening_fee_cents`): ¿deben entrar al costo de una… |
| `WF1-catalogo-huevo-por-kilo.md` | WF1 | ¿Cómo debe registrarse un precio por kilo de un producto cuya unidad base es `pieza` (huevo)? Es decisión de modelo (S5/Ajustes), no del catálogo. |

### Propuestas para Axel de valores operativos (K31 y K44), recogidas de todos los reportes de WF-1
| Valor | Propuesta | Por qué (resumen) | Agente |
|---|---|---|---|
| Hora de los reportes (K31 punto 1) | diario 21:30; semanal domingo 19:00; mensual día 1 a las 08:00 | de noche ya se capturó el día; el domingo se planea la semana; el día 1 el mes ya cerró. Dónde se guarda depende de `S1-parametros-avisos` | `U4b-ajustes` |
| Avisos (K31 punto 2) | `advance_days` de precios viejos = 1; hora de avisos 09:00 (detalle en el reporte) | un día antes de la despensa da tiempo de actualizar precios | `U4b-ajustes` |
| Tiempo máximo del comparador (K31 punto 7) | mantener 5 s (K43); límites de seguridad 1 s a 30 s | peor caso medido 111–123 ms (3,870 candidatos), 2.5 % del tope; casos reales ≤ 1 ms | `motor-S3-comparador-go`, `oraculo-S3` |
| `optimizer.date_window_days` | 7 días | cubre una semana y mueve el primer pago a la quincena vecina; peor caso ≈ 7,600 escenarios, muy por debajo de 5 s | `contrato-S3` |
| `ocr.timeout_seconds` | 30 s | 2 imágenes sintéticas tardan ≈ 2 s; una foto de 12 MP en un equipo lento, 5–10 s; margen de 3× | `ocr`, `contrato-S4S5` |
| Rango de frescura (K44) | `price.stale_after_days` ≥ 1 y `price.zero_weight_after_days` ≥ `price.stale_after_days`; fórmula general del peso (Z − días) ÷ (Z − S + 1), que con 30 y 90 da K40 | con otros valores la recta no baja de 1 a 0; el motor los rechaza con `invalid_argument` | `motor-S5-go`, `U4a-ajustes` |

## Bloqueos nuevos de P2c (11-oct-2026, noche)

Los abrió `contratos-puerta2` al aplicar las respuestas de la puerta 2 en los contratos de los motores (`docs/estado/p2c/contratos-puerta2.md`). Los 4 están **ABIERTOS**, **sin recomendación** (los otros 9 de P2c van en «Dados de alta al cierre de P2c», abajo: 13 en total), y se deciden en WF-3 o cuando Axel conteste por id. Ninguno cambia una cifra canon. Dado de alta en este índice por `correccion-bloqueos-r2`; las opciones y el impacto completos están en cada archivo.

| Archivo | Nació de | Pregunta para Axel | Opciones | Qué queda PROVISIONAL | Fase que lo usa |
|---|---|---|---|---|---|
| `S2-ahorro-quincenal.md` | Q2-13 = B (`S2-flujo-quincenal`, pregunta 2) | El reparto del ahorro mensual planificado (mitad y mitad, HALF_UP en la quincena 1) ya es FIRME; falta (1) si se reparte meta por meta o sobre la suma del mes, (2) qué metas cuentan en una quincena futura y (3) si el ahorro planificado resta del saldo o solo se informa | P1: (A) sobre la suma del mes, (B) meta por meta. P2: (A) todas las metas activas con `planned_monthly_cents`, (B) además solo del mes de `start_date` al de `target_date` y sin las cumplidas. P3: (A) solo se informa, (B) resta del saldo como egreso comprometido | `savings_planned_cents` sigue en `null` y no entra a ningún saldo (`S2-interfaz.md` 1.3 §3.11) | S2 (M-02 y M-10), WF-3 |
| `S2-metas-vencida.md` | Q2-16 = A y Q2-18 = A | ¿Qué aportación muestra la pantalla de Metas para una meta vencida que no se cumplió (`months_left` 0 y restante mayor que 0)? El motor y el oráculo difieren | (A) todo el restante en el mes de hoy (lo que hace el motor), (B) sin aportación calculada: `needed_monthly_cents` null y `plan` vacío (lo que hace el oráculo), (C) `needed_monthly_cents` = restante y `plan` vacío | La opción A, y el oráculo se alinea a ella (`S2-interfaz.md` 1.3 §3.13) | S2 (pantalla de Metas), WF-3 |
| `S4-presentacion-respaldo.md` | Q2-07 = A (hallazgo M-r5-2); la variante quedó sin recomendación (PUERTA-2 §2.3 y §8) | Sin precio en S5, ¿qué presentación rige un artículo de la lista cuando (a) viene de varias recetas con presentación distinta, (b) es solo de reposición por ritmo o (c) es `calculated` con necesidad 0 y sin ingrediente esa semana? | (i) la presentación típica del catálogo (`typical_presentation` en `catalog.products[]`), (ii) para (a) la de la primera receta en orden de fecha y tiempo, y para (b) y (c) además i o iii, (iii) `null` en `presentation`, `packages_to_buy` y `leftover_milli`, con el producto en `missing_prices[]` | La regla del motor: ii para (a), luego i, luego la primera presentación del catálogo (`S4-interfaz.md` 1.5 §6.2 paso 5). El oráculo debe alinearse y falta una fixture con los tres casos | S4, WF-3 |
| `S5-carritos-granel.md` | Q2-03 = A (hallazgo H-r3-02) | Con cantidades a granel, ¿el importe de cada renglón se redondea como en el ticket o se suma exacto? | (A) exacto: HALF_UP de la suma exacta, un solo redondeo (lo que hace el motor), (B) como el ticket: cada línea HALF_UP y el total es la suma de las redondeadas (lo que hace el oráculo), (C) otra regla por línea, por ejemplo la de cada tienda, documentada con un ticket real | La opción A (`S5-interfaz.md` 1.5 §7.2 paso 8). Falta alinear el oráculo y agregar una fixture a granel | S5, WF-3 |

### Dados de alta al cierre de P2c (integrador `integrador-p2c`, 11-oct-2026, noche)

`C1-puerta2` y las correcciones de la ronda abrieron 9 bloqueos más. Los 9 están **ABIERTOS** y ninguno cambia una cifra canon. Ninguno trae recomendación de la puerta; donde el agente dejó una PROPUESTA se anota como tal (no es respuesta). Mientras no se contesten rige la opción en uso, marcada PROVISIONAL con cita del archivo. Las opciones y el impacto completos están en cada archivo.

| Archivo | Lo abrió | Pregunta para Axel | Opciones | Qué queda PROVISIONAL | Fase que lo usa |
|---|---|---|---|---|---|
| `C1-hora-reportes-job.md` | `C1-puerta2`, al aplicar Q2-11 = B (K47) | ¿Qué manda la hora y el día de un reporte programado: las claves `reports.*` de `app_setting` o `job.schedule`? | (A) mandan las claves y la API reescribe `schedule`; (B) manda `job.schedule` y las claves solo informan; (C) quitar `schedule` de los jobs de reportes (migración nueva). PROPUESTA del agente: A | `JobUpdate.schedule` y la fila de C1-datos §13: `updateJob` acepta `schedule` en todos los jobs (en uso) | C1, S7 (WF-3), U4 P-26.08 |
| `C1-reparto-alcance.md` | `C1-puerta2`, al aplicar Q2-68 = C | ¿Qué gastos recalcula `CostShareRuleUpdate.scope = all`? | (A) solo los que nombran la regla; (B) todos los que la usan, también por categoría o por omisión; (C) B limitado a un rango de fechas. PROPUESTA del agente: B | El alcance de `all` (C1-datos §7.2 y §13) | C1, motor de reparto (WF-3), U4 P-26.12 |
| `S1-reparto-edicion-gasto.md` | `correccion-S1-r1` (M-02 de `REVISION-S1-r1.md`) | Tras cambiar una regla con «solo los nuevos», si se edita un gasto anterior, ¿con qué partes se recalcula? | (A) proporción de las partes congeladas; (B) guardar cómo se calculó cada parte (migración); (C) con la regla vigente | La C, que es la que hace hoy `updateExpense` (RN-03 e invariante 18 del diccionario) | S1, C1, S7 |
| `S3-formula-tarjeta-pdf.md` | `correccion-contratos-r1` (M-3 de `REVISION-contratos-r1.md`) | Falta la entrega de los 2 PDF de estado de cuenta (una tarjeta y un crédito de tienda) para escribir la fórmula de cada producto; la respuesta (a) ya es de Axel | Sin recomendación: solo falta la entrega | La fórmula general de la opción (b) de `S3-formula-tarjeta` en `S3-interfaz.md` 1.5 (§5.3, §5.5, §5.8, §6 y §14.1) | S3, WF-3 |
| `S3-iva-lectura.md` | `correccion-contratos-r1`, al aplicar Q2-67 = 1 | (1) Préstamo sin el dato de IVA: ¿no suma IVA, «falta capturar» o la API no deja guardarlo? (2) Crédito «no incluye IVA» sin porcentaje: ¿«falta capturar» con aviso en C1 o la API lo impide? | P1: A, B o C; P2: A o B. Sin recomendación | P1: la A (en uso); P2: `ErrMissingCapture` sin decidir si C1 lo avisa | S3, S1, C1, U4 |
| `S4-empates-sin-confirmar.md` | `correccion-contratos-r1` (M-1 y M-4 de `REVISION-contratos-r1.md`) | Que Axel confirme o corrija los empates de nutrición (Q2-90) y, cuando exista, la lista de ligas PROFECO | (A) confirmar tal cual; (B) corregir. Sin recomendación: faltan datos que solo Axel tiene | Los valores de los renglones `propuesto` (`S4-interfaz.md` §9.1) y las fixtures `s4-kcal-*` | S4, WF-3 |
| `S5-motor-oraculo-reglas.md` | `correccion-motor-S5-r1` (M-1 y M-2 de `REVISION-motor-S5-r1.md`) | Cuatro reglas en que el motor y el oráculo difieren (conteo de piezas, desglose de idénticos, precios por fuente y lectura de fechas) | En cada una: (A) lo que hace el motor, (B) lo que hace el oráculo, (C) otra. Sin recomendación | En uso: el motor; el código de los dos lados cita la pregunta | S5, WF-3 |
| `S5-precio-vigente-captura.md` | `correccion-motor-S5-r2` (M-1 de `REVISION-motor-S5-r2.md`) | Cuatro puntos de `S5-interfaz.md` 1.5 que el contrato no fija: (1) desempate con el mismo `created_at`, (2) precisión de la fracción de segundo, (3) perfil de RFC 3339 de `created_at` (con el año 0000, también en `observed_date`; corrección del 12-oct-2026) y (4) claves obligatorias de `prices[]` | Por pregunta, en el archivo | En uso: el motor, igual que en `S5-motor-oraculo-reglas`. La revisión r3 (M-1, medio) pide sumar el año 0000 a la pregunta 3 y corregir el texto «sin divergencia» | S5, WF-3 |
| `U4-reparto-preseleccion.md` | `correccion-U4-r1` (M-1 de `REVISION-U4-r1.md`) | En el diálogo de alcance de P-26.12, ¿qué opción aparece preseleccionada? | (A) «Solo a los gastos nuevos»; (B) «También a los ya capturados»; (C) ninguna. PROPUESTA de `U4-puerta2`: A | La A, que es la que describe U4 1.5 (P-26.12 y supuesto 6) | U4, puerta 4 |

La pregunta 4 de `S2-ventanas` (Q2-01) no es un bloqueo nuevo: es una pregunta que `contratos-puerta2` redactó dentro del bloqueo existente. **Queda abierta con sus dos casos**: el caso 1 (hueco sin eventos que cruza un devengo) con la regla 1A como PROVISIONAL, y el caso 2 (ventanas contiguas que no duran un mes) con 2A, el comportamiento actual, también PROVISIONAL y sin recomendación.

## Bloqueos nuevos de la Pista 3 (11-oct-2026)

Los dejó el agente `ds-brandbook` de la Pista 3 (sistema de diseño). No pasan por la puerta 2 ni los tocó este índice más que para contarlos; su dueño es la Pista 3 y los archivos no usan el formato «Estado del bloqueo (puerta 2)» (traen «Estado: ABIERTO» en el encabezado). Ninguno detiene la Pista 3.

| Archivo | Pregunta para Axel | Qué queda PROVISIONAL |
|---|---|---|
| `P3-contraste-chip-hover.md` | Contraste de la marca categórica sobre el hover del chip seleccionable: 4 marcas bajo 3:1 sobre `ui-surface-hover`. Opciones: (A) el hover cambia solo el borde, (B) se deja y se declara decorativa la marca, (C) otra. El agente recomienda A | La vista previa respeta DESIGN §7.4 tal cual y el verificador lo marca como excepción declarada |
| `P3-huecos-componentes.md` | Seis huecos de DESIGN.md resueltos como supuesto con tokens existentes (iconos de 4 secciones de la barra lateral, bloque «Antes → Después», perilla del interruptor, círculo de «hoy», comandos de ejemplo del menú «/», pesos y 2 px del anillo). Pregunta explícita: ¿se confirman los iconos? | Todos los supuestos del archivo |

## Respuestas de la puerta 2 (11-oct-2026)

Fuente: `docs/estado/puerta2/RESPUESTAS.md` y hoja de respuestas de `docs/estado/puerta2/PUERTA-2.md` §8. Modo «lo recomendado + excepciones». «Origen» dice si la respuesta fue explícita de Axel o la recomendada aprobada en bloque. Consolidado en `docs/plan/decisiones.md` §K46 a §K52. Q2-01 a Q2-10 (artefactos de WF-1) y Q2-78 a Q2-98 no son bloqueos y se tratan en otras partes de P2c y en §K46, §K49, §K51 y §K52; aquí solo aparece Q2-01 (casos 1 y 2, pregunta 4 de `S2-ventanas`) porque queda abierta. La variante de Q2-07 sin recomendación vive en el bloqueo nuevo `S4-presentacion-respaldo` (ver «Bloqueos nuevos de P2c»).

| Q2 | Bloqueo | Respuesta (letra) | Estado | Origen |
|---|---|---|---|---|
| Q2-01 (casos 1 y 2, pregunta 4 del bloqueo) | `S2-ventanas` | — (sin respuesta; Q2-01 = A autoriza escribir la pregunta, no la cierra) | ABIERTO | caso 1: 1A PROVISIONAL; caso 2: 2A PROVISIONAL, sin recomendación; pasa a WF-3 |
| Q2-11 | `S1-parametros-avisos` | B | RESUELTO | explícita |
| Q2-12 | `S2-flujo-quincenal` | — (sin respuesta) | ABIERTO | PROVISIONAL en uso; pasa a WF-3 |
| Q2-13 | `S2-flujo-quincenal` | B | RESUELTO | recomendada |
| Q2-14 | `S2-flujo-quincenal` | A | RESUELTO | recomendada |
| Q2-15 | `S2-metas-resumen` | A | RESUELTO | recomendada |
| Q2-16 | `S2-metas` | A | RESUELTO | recomendada |
| Q2-17 | `S2-metas` | C | RESUELTO | recomendada |
| Q2-18 | `S2-metas` | A | RESUELTO | recomendada |
| Q2-19 | `S2-presupuesto-total` | B | RESUELTO | recomendada |
| Q2-20 | `S2-presupuesto-total` | A | RESUELTO | recomendada |
| Q2-21 | `S2-proyeccion-interruptor` | C | RESUELTO | recomendada |
| Q2-22 | `S2-proyeccion-interruptor` | C | RESUELTO | recomendada |
| Q2-23 | `S2-proyeccion-interruptor` | C | RESUELTO | recomendada |
| Q2-24 | `S2-proyeccion-interruptor` | C | RESUELTO | recomendada |
| Q2-25 | `S2-ventanas` | A | RESUELTO | recomendada |
| Q2-26 | `S2-ventanas` | H1 | RESUELTO | recomendada |
| Q2-27 | `S2-ventanas` | O1 | RESUELTO | recomendada |
| Q2-28 | `S3-convencion-prestamo` | (a) | RESUELTO | recomendada |
| Q2-29 | `S3-formula-tarjeta` | (a) | RESUELTO | recomendada |
| Q2-30 | `S3-formula-tarjeta` | (a) | RESUELTO | recomendada |
| Q2-31 | `S3-formula-tarjeta` | (a) | RESUELTO | recomendada |
| Q2-32 | `S3-formula-tarjeta` | A (propuesta del agente, letras de PUERTA-2) | RESUELTO | recomendada |
| Q2-33 | `S3-objetivos-medida` | (b) | RESUELTO | recomendada |
| Q2-34 | `S3-objetivos-medida` | A | RESUELTO | recomendada |
| Q2-35 | `S3-objetivos-medida` | A | RESUELTO | recomendada |
| Q2-36 | `S3-objetivos-medida` | A | RESUELTO | recomendada |
| Q2-37 | `S3-origen-cargo-normal` | (a) | RESUELTO | recomendada |
| Q2-38 | `S3-serie-sin-intereses` | (a) | RESUELTO | recomendada |
| Q2-39 | `S4-asignacion-mercado` | A | RESUELTO | recomendada |
| Q2-40 | `S4-asignacion-mercado` | A | RESUELTO | recomendada |
| Q2-41 | `S4-asignacion-mercado` | A | RESUELTO | recomendada |
| Q2-42 | `S4-asignacion-mercado` | A | RESUELTO | recomendada |
| Q2-43 | `S4-confirmacion-empates` | A | RESUELTO | recomendada |
| Q2-44 | `S4-reglas-lista` | A | RESUELTO | recomendada |
| Q2-45 | `S4-reglas-lista` | A | RESUELTO | recomendada |
| Q2-46 | `S4-reglas-lista` | A | RESUELTO | recomendada |
| Q2-47 | `S5-carrito-comparativa` | A | RESUELTO | recomendada |
| Q2-48 | `S5-carrito-comparativa` | — (sin respuesta) | ABIERTO | PROVISIONAL en uso; pasa a WF-3 |
| Q2-49 | `S5-carrito-comparativa` | A | RESUELTO | recomendada |
| Q2-50 | `S5-carritos-detalle` | — (sin respuesta) | ABIERTO | PROVISIONAL en uso; pasa a WF-3 |
| Q2-51 | `S5-carritos-detalle` | — (sin respuesta) | ABIERTO | PROVISIONAL en uso; pasa a WF-3 |
| Q2-52 | `S5-carritos-detalle` | — (sin respuesta) | ABIERTO | PROVISIONAL en uso; pasa a WF-3 |
| Q2-53 | `S5-profeco-almacen` | A | RESUELTO | recomendada |
| Q2-54 | `S5-profeco-almacen` | A | RESUELTO | recomendada |
| Q2-55 | `S5-profeco-almacen` | A | RESUELTO | recomendada |
| Q2-56 | `S5-profeco-filtro` | — (sin respuesta) | ABIERTO | PROVISIONAL en uso; pasa a WF-3 |
| Q2-57 | `S5-profeco-filtro` | — (sin respuesta) | ABIERTO | PROVISIONAL en uso; pasa a WF-3 |
| Q2-58 | `S5-profeco-filtro` | — (sin respuesta) | ABIERTO | PROVISIONAL en uso; pasa a WF-3 |
| Q2-59 | `S5-reglas-comparativa` | A | RESUELTO | recomendada |
| Q2-60 | `S5-reglas-comparativa` | A | RESUELTO | recomendada |
| Q2-61 | `S5-reglas-comparativa` | A | RESUELTO | recomendada |
| Q2-62 | `S5-reglas-comparativa` | B | RESUELTO | recomendada |
| Q2-63 | `S5-reglas-comparativa` | A | RESUELTO | recomendada |
| Q2-64 | `S5-terminos` | A | RESUELTO (decisión); bloqueo ABIERTO al cierre de P2c y RESUELTO por K53 en P2d (puerta 2 bis; lo que sigue abierto, en `S5-agente-externo`) | explícita |
| Q2-65 | `U3-t2-acumulados` | A | RESUELTO | recomendada |
| Q2-66 | `U4-dia-mes-reporte` | B | RESUELTO | recomendada |
| Q2-67 | `U4-iva-credito` | 1 | RESUELTO | recomendada |
| Q2-68 | `U4-reparto-retroactivo` | C | RESUELTO | recomendada |
| Q2-69 | `U4-vista-previa-formulas` | (a) | RESUELTO | recomendada |
| Q2-70 | `U4-vista-previa-formulas` | (a) | RESUELTO | recomendada |
| Q2-71 | `WF1-catalogo-huevo-por-kilo` | 1 | RESUELTO | recomendada |

### Valores operativos K31 y K44 (Q2-72 a Q2-77), todos A: aprobados

| Q2 | Valor | Respuesta | Estado |
|---|---|---|---|
| Q2-72 | Hora de los reportes: diario 21:30; semanal domingo 19:00; mensual día 1 a las 08:00 | A | FIRME (K50) |
| Q2-73 | Avisos: anticipación de precios viejos 1 día; hora de avisos 09:00 | A | FIRME (K50) |
| Q2-74 | Tiempo máximo del comparador 5 s; límites de seguridad de 1 s a 30 s | A | FIRME (K50) |
| Q2-75 | `optimizer.date_window_days` = 7 días | A | FIRME (K50) |
| Q2-76 | `ocr.timeout_seconds` = 30 s | A | FIRME (K50) |
| Q2-77 | Rango de frescura: `price.stale_after_days` ≥ 1 y `price.zero_weight_after_days` ≥ el anterior; peso (Z − días) ÷ (Z − S + 1) | A | FIRME (K50) |

Dónde se guardan las claves de Q2-72 y Q2-73 lo fija Q2-11 = B (`app_setting`).

## Bloqueos nuevos de P2d (12-oct-2026, madrugada)

Lo dio de alta `p2c-remanentes` con la carpeta en disco. El archivo lo escribió `k53-aplicar` al registrar K53 y puede cambiar mientras P2d corre; las opciones, lo PROVISIONAL y el impacto completos están en el archivo. No decidí nada.

| Archivo | Lo abrió | Pregunta para Axel | Opciones | Qué queda PROVISIONAL | Fase que lo usa |
|---|---|---|---|---|---|
| `S5-agente-externo.md` | `k53-aplicar`, al registrar K53 (Q2-64, puerta 2 bis) | (a) Peso de los precios del agente externo en la comparativa; (b) qué puede escribir el agente; (c) si sus recetas entran como borrador. Y las que salieron contra el esquema 1.4: (d) valor de `price.source`, (e) credencial de servicio y (f) borrador y origen de una receta | Por pregunta, en el archivo. PROPUESTA del orquestador para (a) a (d) y de `k53-aplicar` para (e); (f) sin propuesta | (a) compiten como `manual`; (b) solo productos, presentaciones, precios y recetas; (c) borrador; (d) la (d2), `site_assisted`; (e) solo el transporte `serviceToken` de C1 1.4 | C1, S1, S5, WF-3, puerta 4 |

`S5-terminos` (WF-1) quedó **RESUELTO por K53** en su línea «Estado del bloqueo (puerta 2 bis)». Lo que K53 dejó abierto vive en `S5-agente-externo`.

## Bloqueos de la semilla de WF-3a (12-oct-2026)

Los abrió el agente `semilla` con el nombre `semilla-*.md`; `correccion-backend-core-r1` los renombró a `S6-semilla-*.md` para que los encuentre `ls docs/estado/bloqueos/S6-*` (criterio S6-16, hallazgo M3 de REVISION-backend-core-r1) y los dio de alta aquí. El contenido no cambió; las opciones, lo PROVISIONAL y el impacto completos están en cada archivo. No decidí nada.

| Archivo | Lo abrió | Pregunta para Axel | Opciones | Qué queda PROVISIONAL | Fase que lo usa |
|---|---|---|---|---|---|
| `S6-semilla-despensas-futuras.md` | `semilla` | ¿Las despensas 2, 3 y 4 (10-dic-2026, 10-feb-2027, 10-abr-2027) se siembran como compras a crédito activas sin gasto, o las crea la proyección? | A: 4 `credit_purchase` con sus 16 pagos de `motor_cred`; B: solo la despensa 1; C: A con `status = 'cancelled'` | A (cambia cifras sembradas: deuda viva al 10-abr-2027 = $6,293.44, igual al oráculo) | S6 (WF-3), C3-09 |
| `S6-semilla-sucursales-municipio.md` | `semilla` | ¿En qué alcaldía o municipio y estado están las sucursales Bodega Aurrera Vallejo y Walmart Acueducto de Guadalupe? | A: Axel las captura en Ajustes; B: Axel las dicta y van a la semilla; C: suponerlas (descartada) | Sucursales sin municipio; 0 precios de la muestra PROFECO cargados | S5/S6 (PROFECO) |
| `S6-semilla-catalogo-alias.md` | `semilla` | ¿Dónde viven los sinónimos y los nombres de ticket del catálogo una vez sembrado? | A: se leen de `productos.json`; B: migración nueva con tablas de alias; C: no usarlos | No se guardan en la base (el esquema 1.4 no tiene columna) | S5 (OCR y búsqueda) |

## Bloqueos nuevos de WF-3a (12-oct-2026, para WF-3b, la puerta 4 y WF-4)

Los dio de alta `integrador-wf3a` al cerrar WF-3a. Son 37 archivos, más los 3 de la semilla de la sección anterior: **40 bloqueos nuevos de WF-3a, todos ABIERTOS**. El código ya funciona con la opción marcada PROVISIONAL en cada archivo, que cita su bloqueo; ninguno inventa una cifra (los motores y sus oráculos no se tocaron). La pregunta exacta, las opciones y el impacto están en cada archivo. Los valores operativos de la puerta 4 (K31) con su propuesta están resumidos en `docs/estado/wf3/INTEGRACION-wf3a.md` §7.

| Archivo | Lo abrió | Tema (pregunta para Axel o para el siguiente flujo) | K31 | Fase que lo usa |
|---|---|---|---|---|
| `C3-matriz.md` | `integrador-frontend-c3` | 5 preguntas de la matriz C3 v0.9: fila P-26, 7 filas con «implementado = no», 7 operaciones que el frontend no llama, regla de «probado» y semilla sin datos de octubre | — | WF-3b (firma de C3) |
| `S6-adjuntos.md` | `auth-sistema` | Factor de «10 MB», tipos admitidos (JPG, PNG, WebP y PDF; ¿HEIC?), `kind` por omisión y 409 al borrar | sí | puerta 4 |
| `S6-ajustes-api.md` | `api-ajustes` | 6 puntos de Ajustes que C1, S1, RN y U4 no fijan (p. ej. `PlannedItem.next_date`) | — | puerta 4 |
| `S6-comidas-calculos.md` | `api-comidas-salud` | Cálculos de comidas y salud que C1 y S4 no fijan | — | puerta 4 |
| `S6-compras-a-motores.md` | `api-compras-a` | Cómo la API de compras (parte a) alimenta y lee los motores S4 y S5 | — | puerta 4 |
| `S6-compras-b-fuentes.md` | `api-compras-b` | Veredicto de cada cadena con sitio y fuentes de precio tras K53 | — | puerta 4 |
| `S6-compras-b-nutricion.md` | `api-compras-b` | Recarga de la base de nutrición, recetario e inventario | — | puerta 4 |
| `S6-creditos.md` | `api-creditos` | Decisiones PROVISIONALES de la API de créditos | — | puerta 4 |
| `S6-gastos-inicio.md` | `api-gastos-inicio` | Reparto RN-03 fuera del motor y huecos de gastos e Inicio (§15.4 y §15.5 esperan a Axel) | — | puerta 4 |
| `S6-importar-partidas-csv.md` | `api-ajustes` | Correspondencia de columnas de `gastos_fijos_ingresos.csv` (plan §8.2) | — | puerta 4 |
| `S6-presupuesto-motor.md` | `api-presupuesto` | Entradas del motor S2 desde la base | — | puerta 4 |
| `S6-presupuesto-periodos.md` | `api-presupuesto` | Cuándo se llena la tabla `period` | — | puerta 4 |
| `S6-presupuesto-quincenas.md` | `api-presupuesto` | Plan de la quincena en `getBudgetFortnights` | — | puerta 4 |
| `S6-respaldos.md` | `auth-sistema` | Hora (03:15) y retención (30) del respaldo; qué conserva una restauración | sí | puerta 4 y WF-4 |
| `S6-sistema.md` | `auth-sistema` | Decisiones abiertas de autenticación y sistema (entre ellas Android y `.local`) | sí | puerta 4 y WF-4 |
| `S6-valores-operativos-base.md` | `backend-base` | Respaldo, argon2id, cabeceras, límites de cuerpo, job de purga de idempotencia, rutas por omisión y campo desconocido = 400 | sí | puerta 4 |
| `S7-convivencia-multix.md` | `s7-avisos` | Cómo detectar el chat de Ximena sin chocar con MultiX | — | puerta 4 |
| `S7-estado-cuenta-fechas.md` | `s7-pdf` | Fechas de pago del estado de cuenta de Karen en C1 | — | puerta 4 |
| `S7-fuentes-avisos.md` | `s7-avisos` | Día de despensa, fuentes de otras áreas para las alertas, semana que revisa la alerta 4, cadencia y reintentos | sí | puerta 4 |
| `S7-lualatex.md` | `s7-pdf` | `lualatex` sin `luaotfload` (fuentes OpenType) en la máquina de casa | — | WF-4 |
| `S7-reporte-mensual-secciones.md` | `s7-pdf` | Secciones del reporte mensual que dependen de otras áreas y signo negativo | sí | puerta 4 |
| `S7-ruta-token.md` | `s7-avisos` | Ruta exacta de la clave del token en `openclaw.json` (K30) | — | WF-4 y S7-07 |
| `WF3-ajustes-a.md` | `front-ajustes-a` | Ajustes, secciones 01 a 07 | — | WF-3b |
| `WF3-ajustes-compra-ejemplo.md` | `correccion-frontend-areas-r1` | Intereses de la compra de ejemplo en P-26.05 (vistas previas 6 y 7) | — | WF-3b |
| `WF3-ajustes-muestra-tema.md` | `front-ajustes-b` | Muestra de tarjeta en claro y oscuro (P-26.10) | — | WF-3b |
| `WF3-ajustes-precios-conteo.md` | `correccion-frontend-areas-r1` | Conteo de precios que pasan de fresco a viejo en P-26.06 | — | WF-3b |
| `WF3-ajustes-registro.md` | `front-ajustes-b` | Registro de P-26 entre dos agentes | — | WF-3b |
| `WF3-comidas-compras-interfaz.md` | `front-comidas-compras` | Interfaz PROVISIONAL de P-10 a P-16 | — | WF-3b |
| `WF3-creditos-salud-reportes.md` | `front-creditos-salud-reportes` | P-17 a P-21 y P-23 a P-25 | — | WF-3b |
| `WF3-frontend-k31.md` | `frontend-base` | Sondeo de `/health` (15 s), reintento de la cola (30 s), copia de lecturas y cola al cerrar sesión | sí | puerta 4 |
| `WF3-frontend-rutas.md` | `frontend-base` | Direcciones (URL) de las pantallas | — | WF-3b |
| `WF3-frontend-servidor.md` | `frontend-base` | Cómo sirve el servidor la PWA: mismo origen, respaldo SPA, `Cache-Control`, tipo del manifiesto y CSP | sí | WF-3b y WF-4 |
| `WF3-gastos-contrato.md` | `front-gastos-presupuesto` | Huecos de C1 v1.4 para P-02, P-04 y P-05 | — | WF-3b |
| `WF3-gastos-ticket-credito.md` | `front-gastos-presupuesto` | Pagar con crédito desde P-05 | — | WF-3b |
| `WF3-inicio-tiendas.md` | `front-inicio` | Color de tienda sin código en C1 | — | WF-3b |
| `WF3-inicio-vacios.md` | `front-inicio` | Reglas de «vacío» y destinos de P-01 que C1 y U3 no fijan | — | WF-3b |
| `WF3-integrador-red.md` | `integrador-backend` | IP fija 192.168.50.56, nombre `finanzas.local` permanente y CA en el sistema | — | WF-4 |

## Bloqueos nuevos de WF-3a bis (12-oct-2026, para WF-3b y la puerta 4)

Los dio de alta `integrador-wf3a-bis` al cerrar la ronda acotada de corrección del frontend (WF-3a bis). Son **4 archivos nuevos, todos ABIERTOS**, más `C3-matriz.md` (ya contado en WF-3a), que ahora trae la sección «Decisiones del orquestador (2026-10-12, WF-3a bis)»: sus decisiones 1 a 5 están aplicadas, pero siguen PROVISIONALES hasta la puerta 4, así que el archivo sigue ABIERTO. **Conteo con la carpeta en disco: 129 archivos de bloqueo** (85 hasta P2d y la Pista 3, 40 de WF-3a y 4 de WF-3a bis). El código funciona con la opción PROVISIONAL de cada punto, que cita su archivo; ninguno cambia cifras de los motores ni sus oráculos.

| Archivo | Lo abrió | Tema (pregunta para Axel o para el siguiente flujo) | K31 | Fase que lo usa |
|---|---|---|---|---|
| `WF3a-bis-ajustes-a.md` | `front-ajustes-a-fix` | 3 puntos de la edición en grupo con `Ctrl+/` en Ajustes 02, 04, 05 y 07: avisos sin cifra en el bloque sumado, qué campos de una promoción se editan a la vez y cómo se cambia una selección múltiple (tiempos aptos) | — | WF-3b y puerta 4 |
| `WF3a-bis-comidas-compras.md` | `front-comidas-compras-fix` | 4 puntos: cifras de P-14 con la semilla de desarrollo distintas de las canon de U1 §3.16 (Walmart $3,359.00 contra $3,233.00), P-13 sin lista sembrada, nombre visible de una sucursal y `Esc` dentro de una hoja | — | WF-3b (semilla y C3-09) |
| `WF3a-bis-gastos-presupuesto.md` | `front-gastos-presupuesto-fix` | 7 puntos: filtros de `queryView` que faltan para N-040, id del adjunto en `Expense` (N-064), galería N-046, N-073, N-053 con préstamos a plazo y tarjetas (R5-03 de `REVISION-frontend-core-r5`), `listPeriods` como excepción técnica y detalles menores | — | WF-3b y puerta 4 (dueño de C1) |
| `WF3a-bis-semilla-dev.md` | `semilla-dev` | 7 puntos de la semilla de desarrollo (`finanzas semilla --desarrollo`): renta contra P1–P6, candado «nunca en producción», presupuesto solo de Despensa, P3–P6 de `getAccrualVsCash` contra model.json, costo y kcal de la semana en 0, ingresos y precios de Walmart | — | WF-3b (C3-09) y puerta 4 |

## Bloqueos nuevos de WF-3b (12-oct-2026, para la puerta 4)

Los dio de alta el flujo WF-3b «Verificación visual y C3»; los registra `integrador-wf3b`. Son **8 archivos nuevos, todos ABIERTOS**, entre ellos los 2 escalamientos de ronda 3 (U5 y C3 quedaron devueltos en la ronda 3 y no hay ronda 4). Además, `WF3a-bis-semilla-dev.md` (ya contado en WF-3a bis) trae ahora las **decisiones del orquestador** sobre sus puntos 4 (opción a: gastos proyectados de la fixture S2) y 7 (sí: 36 precios de Walmart en la semilla real), aplicadas por `semilla-c3` y PROVISIONALES hasta la puerta 4; sus puntos 1, 3 y 5 siguen como pregunta a Axel. **Conteo con la carpeta en disco: 137 archivos de bloqueo** (129 hasta WF-3a bis y 8 de WF-3b). Ninguno cambia cifras de los motores ni sus oráculos. Detalle: `docs/estado/wf3/INTEGRACION-wf3b.md` §6.

| Archivo | Lo abrió | Tema (pregunta para Axel o para el siguiente flujo) | K31 | Fase que lo usa |
|---|---|---|---|---|
| `U5-ronda3-escalamiento.md` | `revisor-u5-r3` | U5 devuelto en la ronda 3 (0 altos, 2 medios: M-7 iconos de la barra lateral; M-8 marca «capturado» en C3, ya presente en C3 v1.2). Opciones A (corrección acotada y verificación puntual, recomendada), B (aceptar con pendientes con fecha) o C (otra ronda completa) | — | puerta 4 |
| `C3-ronda3-escalamiento.md` | `revisor-c3-r3` | C3 v1.2 devuelto y sin firma en la ronda 3: 1 alto (R3-1, F-3 ruta A falla con la semilla; pasan 7 de 8 flujos). Opciones A (U1 aclara la ruta A + verificación acotada de C3-08 y firma, recomendada), B (precios de ingredientes en la semilla) o C (firmar con excepción) | — | puerta 4 |
| `C3-flujos.md` | `c3-matriz` | Flujos de U1 contra el servidor real (C3-08): punto 2 (F-3 ruta A, caras 2a y 2b) abierto; punto 3 lo contesta K46 (+$264.50, corregir los ejemplos de U1 y U3); punto 4 (F-5 con oct-2026 y 6 toques); observación sobre Q2-49 (precios de Walmart en la semilla) | — | puerta 4 |
| `C3-capturas-global.md` | `c3-matriz` | «Capturado» de la fila Global de C3 (G-01 a G-07): (a) PROVISIONAL «no», (b) aceptar P-01 y U5-13 como evidencia, (c) capturas propias de G-03, G-04 y G-06, (d) sacarla del conteo | — | puerta 4 |
| `U5-conteo-carpeta.md` | `u5-capturas` | El comando literal de U5-01 da 168 (164 PNG + `axe/`, `contraste/`, `evidencia/`, `hosts.json`); A PROVISIONAL: contar los PNG (164 de 164) | — | puerta 4 (criterios) |
| `U5-contraste-ui-controles.md` | `u5-capturas` | Qué cuenta como «componente de UI» en el 3:1 de U5-04 (campo, selección, botón); A PROVISIONAL: contar los tres subtipos | — | puerta 4 (DESIGN/C2) |
| `U5-nombres-propios.md` | `revisor-u5-r1` | Si se aceptan en pantalla Claude, Ollama, Tesseract y las siglas API, CSV, PDF y OCR (U5-07) | — | puerta 4 |
| `U5-tipografia-tres-tamanos.md` | `revisor-u5-r1` | «Máximo tres tamaños por pantalla» de DESIGN §3.2 contra los 4 a 5 tamaños que piden sus propios componentes; A PROVISIONAL: solo el área de contenido | — | puerta 4 (DESIGN) |

## Bloqueos nuevos de WF-3c (12-oct-2026, para la puerta 4)

Los abrió el flujo WF-3c, «Inicio como tablero de control» (K54), y los registra `integrador-wf3c`. Son **14 archivos nuevos y todos están ABIERTOS**. Con esta tabla se cierran los bajos B-3 (`REVISION-tablero-backend-r3.md`) y B-r1-4 (`REVISION-especificacion-r3.md`), porque el índice no tenía filas del tablero.

Mientras no haya respuesta, el código y los documentos corren con la opción PROVISIONAL de cada archivo. Ninguno toca valores K31. Las cifras de los motores no cambian, y las del tablero salen de RN-30 y RN-31 sobre los motores S2 y S3.

| Archivo | Lo abrió | Tema (pregunta para Axel o para el siguiente flujo) | K31 | Fase que lo usa |
|---|---|---|---|---|
| `U3-tablero-barometro.md` | `u3-v2` | Qué entra en la capacidad de endeudamiento, que es lo FIRME de K54. Hay 4 preguntas: (1) qué es el ingreso mensual neto: (a) las quincenas capturadas más las esperadas, sin extras; (2) qué pagos cuentan: (a) solo la parte del hogar; (3) redondeo: (a) entero con ROUND_HALF_UP, ≤ 30 verde, 31 a 40 amarillo, ≥ 41 rojo; (4) escala del arco: (a) 0 a 60 %. B-R3-1: la justificación del punto 4 trae una cuenta errónea. | — | puerta 4 |
| `U3-tablero-horizonte.md` | `u3-v2` | Cuántos meses ve cada gráfica. G-1: (a) 11 meses anteriores más el actual. G-6: (a) 6 anteriores, el actual y 5 siguientes. Petición: (a) una sola, de 17 meses. | — | puerta 4 |
| `U3-tablero-mas-del-hogar.md` | `u3-v2` | Qué pasa con las tarjetas 2, 4, 5 y 6 de U3 1.x: (a) siguen sin cambio bajo «Más del hogar». Las tarjetas 1 y 3 pasan a K-1, K-3 y K-4. | — | puerta 4 |
| `U3-tablero-recordatorio.md` | `u3-v2` | «Pagos próximos» (K-4). N: (a) 5 (`dashboard.payment_reminder_count`, de 1 a 20). Rango: (a) los vencidos y luego de hoy a hoy + 30 días. Forma: (a) una fila por pago. | — | puerta 4 |
| `U3-tablero-senal-presupuesto.md` | `u3-v2` | Sugerencia de subir o recortar el presupuesto (G-3, RN-31): (a) 3 meses cerrados (`budget.signal_months`); (a) margen de 10 % (`budget.signal_margin_pct`); (a) un mes sin plan corta la racha. | — | puerta 4 |
| `C1-tablero.md` | `s1-c1-v15` | Lo que C1 v1.5 y el esquema 1.5 fijan sin decisión. (1) Sección de Ajustes de las 5 claves: (a) Presupuesto, P-26.02. (2) Rangos de los disparadores de la 00005. (3) Necesidades de U1 de las 5 operaciones nuevas. (4) Nombre de la clave del recordatorio: (a) `dashboard.payment_reminder_count`. (5) Formato de `period` en `getDashboardDebtCapacity` y `getBudgetSignals`: (a) solo `AAAA-MM`, y una quincena responde 422. | — | puerta 4 (U4 y U1 pueden cerrar los puntos 1 y 3) |
| `S6-tablero-sin-datos.md` | `api-tablero` | Cuándo un mes «no tiene datos»: (a) devengado y efectivo del hogar = 0, en `dominio.MonthHasData`; (b) antes del primer gasto o pago registrado. | — | puerta 4 |
| `S6-tablero-senal-periodo.md` | `correccion-tablero-backend-r1` | Qué meses mira `getBudgetSignals` cuando `period` es posterior al mes en curso: (a) solo los cerrados, y entonces la señal es `none`; (b) 422 `out_of_range`. B-5 pide que C1 lo diga. | — | puerta 4 |
| `S6-tablero-lineas-archivadas.md` | `correccion-tablero-backend-r1` | Si los pagos de una línea archivada cuentan en el barómetro y en las gráficas: (a) sí, en los meses en que hubo pagos, y es la PROVISIONAL **sin aplicar en el código**; (b) solo las líneas activas, que es lo que hace hoy el código. Es el hallazgo M-2, abierto tras 3 rondas (ver `WF3c-ronda3-escalamiento`, punto 2). | — | orquestador (salida i o ii) y puerta 4 |
| `C3-ajustes-tablero.md` | `c3-tablero` | `updateDashboardSettings` (N-278) no tiene pantalla. Opciones: un apartado «Parámetros del tablero» en P-26.02, una excepción justificada, u otra sección. Es el H-3 de C3. | — | puerta 4 (U4) |
| `WF3c-barra-iconos-seccion.md` | `front-tablero` | 4 iconos de sección de la barra lateral que DESIGN §5.4 no nombra: Presupuesto `wallet`, Salud `heart-pulse`, Reportes y Gastos. A) confirmarlos; B) Gastos con `arrow-up-right`; C) otros. | — | puerta 4 (DESIGN) |
| `WF3c-tablero-interfaz.md` | `front-tablero` | 3 puntos de interfaz que U3 2.0 no cierra. (1) K-4 dice «hay más pagos» sin número. (2) K-3 sin «Reintentar» dentro de la tarjeta-enlace. (3) La insignia de K-1 es solo texto. | — | puerta 4 |
| `WF3c-barometro-sin-dato.md` | `correccion-tablero-frontend-r2` | Rol accesible del barómetro sin dato: A) `role="img"` con su texto, que es lo que está hecho; B) `sr-only`. DESIGN §8.9 debe fijarlo. | — | puerta 4 (DESIGN) |
| `WF3c-ronda3-escalamiento.md` | `integrador-wf3c` | Lo que WF-3c no aprobó. Especificación: M-R3-1, la escala de G-6. Backend: M-2. C3 v1.3 devuelto: H-1 a H-5. `make verify-wf3` sale con EXIT 2 por U5-04: 2 pares del `<svg>` raíz del barómetro, un falso positivo comprobado. P-01 Detalle mide 475 px en Android. La PROPUESTA es una ronda acotada, «WF-3c bis». | — | orquestador, antes de la puerta 4 |

**Actualizado en WF-3c:** `C3-flujos.md`. En el punto 2, el orquestador decidió de forma PROVISIONAL la opción (c): la ruta A de F-3 exige una lista vigente con todos sus precios en al menos una tienda, y U1 §4 F-3 ya lo dice. Además, `c3-tablero` dejó la nota de las 4 pruebas de flujos que el Inicio de U3 2.0 dejó viejas.

**Conteo de WF-3c:** 14 archivos nuevos (14 abiertos) y 1 actualizado (`C3-flujos.md`, sigue abierto).
