# Índice de bloqueos · Finanzas Hogar

Consolidado por el integrador de WF-0 (`integrador-wf0`) el 11-oct-2026 para la **puerta 1**. Reúne los 32 archivos de esta carpeta (31 abiertos y 1 resuelto) y 2 preguntas que dejó la revisión independiente sin archivo propio. En total son **31 preguntas para Axel**, agrupadas por la fase que necesita la respuesta. Cada archivo conserva el detalle completo; aquí van la pregunta exacta, las opciones con su consecuencia, el impacto y lo que quedó PROVISIONAL mientras tanto.

- **Ningún bloqueo esconde una cifra inventada.** La revisión r3 lo comprobó (`docs/estado/wf0/REVISION-r3.md` §4): lo provisional lleva la marca PROVISIONAL y cita su bloqueo, y las cifras de C1 salen de `oraculo/model.json` o de `datos/semilla/data.json`.
- **Una respuesta cierra dos archivos** en 4 casos: `C1-lista-compras` = `S1-lista-compras`; `C1-huecos-u1-s1` = `S1-huecos-negocio` (3 preguntas); `S1-frescura-peso` = `criterios-frescura-peso`; `S1-meta-kcal` = `criterios-meta-kcal`. Son complementarios `S1-restricciones-optimizador` y `criterios-optimizador`, y `criterios-login-bloqueo` con `C1-parametros-tecnicos`.
- **Cuándo responder.** La columna «Fase que la usa» sale del «Impacto» de cada archivo. Las de la tabla A fijan C1 y U1, que se aprueban en la puerta 1. Las de la tabla B las necesitan los motores de WF-1, que arranca después de la puerta 1. Las de la tabla C son de diseño (WF-1, se aprueban en la puerta 2). Las de la tabla D son de construcción y operación (WF-3 y WF-4).
- Cuando Axel responda, el orquestador anota la respuesta en cada archivo y regenera lo afectado (C1 v1.1 con `backend/internal/api/c1gen/build.py`; migración 00002 si cambia el esquema).

> **Estado al 11-oct-2026 (puerta 1): las 31 preguntas quedaron RESUELTAS.** Axel respondió en la sesión; la respuesta está al final de cada archivo (sección «Respuesta de Axel») y consolidada en `docs/plan/decisiones.md` §K. Las dos preguntas sin archivo ya lo tienen: `C1-decimales-tasa.md` (A12) y `S1-tope-quincena-karen.md` (A13). Lo que sigue: aplicar las respuestas (C1 v1.1, migración 00002, U1, criterios y colorimetría) en el flujo de la puerta 1.

## Resumen

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

