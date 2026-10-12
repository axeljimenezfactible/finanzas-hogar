# P2c · correccion-contratos-r1 · 11-oct-2026, noche

Objeto: los hallazgos altos y medios de `docs/estado/p2c/REVISION-contratos-r1.md` (A-1, A-2 y M-1 a M-4) sobre lo que escribió `contratos-puerta2`.
Alcance: `S3-interfaz.md` (1.4 → 1.5), `S4-interfaz.md` y `S5-interfaz.md` (siguen en 1.5; solo marcas y fixtures), tres bloqueos nuevos y este reporte. No toqué `S2-interfaz.md`, `S2-ventanas.md` ni `oraculo/fixtures/FORMATO.md`, porque ningún hallazgo alto o medio los toca; los bajos B-1 a B-5 no entran en esta corrección.
Fuentes: `RESPUESTAS.md`; `PUERTA-2.md` §3.9, §3.11 a §3.15, §3.18 a §3.20, §3.24 y §5 (Q2-90); K46; RN-26 1.4; `00004_puerta2.sql`; `U4-ajustes-a.md` P-26.05; las fixtures de `oraculo/fixtures/s3`, `s4` y `s5`.
No hice `git add` ni `commit` y no toqué la Pista 3.

## Hecho
- **A-1 (alta) · Q2-67 = 1 aplicada en S3 (FIRME).**
  - §5: `Credit` gana la entrada `VATIncluded *bool` (`vat_included`). Sus valores: true = la tasa ya trae IVA; false = se suma `VATOnInterestBP`; null = falta capturar.
  - §5.3, §5.5 y §5.8: la lista de K18 sube a **5 datos** (`AnnualRate`, `RateBase`, `VATIncluded`, `MinPaymentPctBP`, `MinPaymentFixed`). Si falta cualquiera, el motor responde `ErrMissingCapture`, igual que dicen RN-26 1.4 y `Credit.missing_terms` con `vat`. El IVA de la fórmula se suma solo con `VATIncluded` = false.
  - §5.4 y §5.8: `Amortize` recibe `vatBP` = «IVA que se suma» (nil = ninguno). `QuotePurchase` lo toma de `VATIncluded`.
  - §9 y §10: la semilla Kueski lleva `"vat_included": true` (decisión C, mapeo de la 00004 y solicitud 4 de `S1-puerta2.md`). Las fixtures `s3-tarjeta-minimo` y `s3-tienda-abonos-interes` describen `vat_included`.
  - §11: `TestTarjetaSinCaptura` (5 datos), `TestDespachoCotizacion` (préstamo con true, false y null) y `TestProvisionalListas`.
  - §14 y §14.1: `U4-iva-credito` sale de todas las listas y fixtures.
  - Q2-67 no cubre dos casos, que quedan PROVISIONALES con el bloqueo nuevo **`S3-iva-lectura`** (ABIERTO, sin recomendación):
    1. Un préstamo con `VATIncluded` null. Sigue lo que está en uso: no suma IVA y lleva `Provisional` = [`S3-iva-lectura`]. U4 no captura el IVA del préstamo.
    2. `VATIncluded` = false sin porcentaje. Responde `ErrMissingCapture`, sin número.
- **A-2 (alta) · `verify-c1` sigue fallando.** Ese error no es de los contratos: viene de `gen2_schemas.py` y de la 00004, que están fuera de mi ALCANCE, y no los toqué. Coincide con H-01 de `REVISION-S1-r1.md` y con la solicitud 1 de `S1-puerta2.md`. Va en solicitudes fuera de alcance.
- **M-1 (media) · Marcas que citaban bloqueos resueltos, quitadas.**
  - S4 §9.1: quién confirma y cuándo pasa a FIRME (Q2-43 = A). La confirmación misma no ocurrió: Q2-90 no tiene recomendación ni respuesta, y los 150 renglones de `tabla.csv` siguen `propuesto`. Por eso esos valores quedan PROVISIONALES con el bloqueo nuevo **`S4-empates-sin-confirmar`**.
  - S4 §11: vertical `market` FIRME (Q2-39 a Q2-41).
  - S4 §12: `provisional` solo para bloqueos ABIERTOS (`S4-presentacion-respaldo` y `S4-empates-sin-confirmar`). Fila `s4-canal-mercados` FIRME. La fila `s4-kcal-recetas-autor` cita el bloqueo nuevo.
  - S5 §15: solo llevan `provisional` los casos de bloqueos ABIERTOS. Fila `s5-cotizacion-recomendada` FIRME (Q2-60). `s5-profeco-muestra` cita `S5-profeco-filtro` (ABIERTO), porque sus conteos dependen del filtro (PUERTA-2 §3.19).
  - Cada contrato trae su nota en «Cambios de la versión 1.5».
- **M-2 (media) · S3 §5.8, fila `bank_card`.** Pasan a FIRMES el capital = monto (Q2-38), las fechas (Q2-32; la fila decía «Fechas PROVISIONALES, pregunta 4», que también estaba resuelta) y `quoted` (Q2-37). En §10 (l. 594) se quitó «PROVISIONALES» y la línea remite a §14.1.
- **M-3 (media) · `S3-formula-tarjeta`.** Usé la salida que el hallazgo da como alternativa: el bloqueo nuevo **`S3-formula-tarjeta-pdf`** (ABIERTO). Lo que tiene es la respuesta (a) a Q2-29 a Q2-31, la espera de los PDF y la fórmula (b) en uso. Las 11 marcas de S3 y las listas `Provisional` (interés, mínimo y corte `computed`) ahora lo citan. No reabrí `S3-formula-tarjeta.md` porque no está en mi ALCANCE. La alta en INDICE y en K46 va en solicitudes.
- **M-4 (media) · Lista de cambios para los motores.** Va completa abajo, como complemento de la de `contratos-puerta2.md`: las 14 fixtures, las de S3 y los comentarios PROVISIONAL del código, con archivo y línea.
- **Bloqueos nuevos (ABIERTOS, sin recomendación):**
  - `docs/estado/bloqueos/S3-formula-tarjeta-pdf.md`
  - `docs/estado/bloqueos/S3-iva-lectura.md`
  - `docs/estado/bloqueos/S4-empates-sin-confirmar.md`

## Verificación
- Reemplazos con conteo exacto: `/tmp/p2c-ccr1/edit_s3.py` (40), `edit_s45.py` (S4: 5 + 1 inserción; S5: 3 + 1) y un ajuste de redacción en el encabezado de S3. Cada script aborta sin escribir si un texto no aparece las veces esperadas. Respaldo en `/tmp/p2c-ccr1/bak/`.
- PROVISIONAL contra el estado de los bloqueos, con `/tmp/p2c-ccr1/prov.py` sobre S2 a S5: **0 citas a bloqueos no ABIERTOS** fuera de la bitácora. Las únicas que quedan son renglones de historia de versiones anteriores: S4 l. 26 y 31; S5 l. 22, 28 y 34. Estado de los bloqueos nuevos y citados: `S3-formula-tarjeta-pdf`, `S3-iva-lectura`, `S4-empates-sin-confirmar`, `S4-presentacion-respaldo` y `S5-profeco-filtro`, todos ABIERTO.
- `grep -cF` de los textos señalados, todos con 0 apariciones: «PROVISIONAL `S3-serie-sin-intereses`», «PROVISIONAL (`S3-origen-cargo-normal`)», «Fechas PROVISIONALES», «sin canon y PROVISIONALES», «lista de **4**», «4 datos», la marca de §9.1 de S4, la regla de §12 con `S4-reglas-lista.md`, la de §15 de S5 y la fila `s5-cotizacion-recomendada`.
- `make -k verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP=…/datos/tmp/correccion-contratos-r1/vw` termina con rc=2:
  - WF-0 OK (SQL con la 00004, Go, colorimetría y OpenAPI).
  - **verify-c1 FALLA**: `gen2: tipo de parámetro desconocido: hora HH:MM`, `make: *** [Makefile:130: verify-c1] Error 1`. Es A-2, fuera de alcance, y no cambió.
  - S1: RESUMEN TODO OK.
  - U1: 63 OK, 0 pendientes.
  - Criterios: 261 definidos, 178 comprobadas con 0 diferencias e INDICE con 0 diferencias.
  - K25: motor_cred 93.7 %, motor_fin 89.2 %, recetas 93.2 % y precios 96.1 %.
  - Oráculos: S2 30/30, S3 27 fixtures (906/906 y 437/437 hojas canon), S4 26/26 y S5 29/29, con 0 diferencias.
  - `model.json` sha256 `1a57066689dd…` intacto.
  - WF-1: verificación completa OK.
- No cambié código, fixtures ni salidas. Los cambios de abajo son para WF-3.

## Supuestos
1. Para M-3 elegí el bloqueo nuevo y no reabrir `S3-formula-tarjeta.md`, porque ese archivo no está en mi ALCANCE («bloqueos nuevos»).
2. `VATIncluded` es `*bool` (JSON true, false o null), igual que `ReadingConfirmed bool` ↔ `availability_confirmed` INTEGER.
3. Q2-67 = 1 es FIRME para la tarjeta y la tienda (K18, RN-26 1.4). Para el préstamo no hay fuente sobre NULL, y U4 no captura su IVA. Por eso abrí `S3-iva-lectura` y no apliqué `ErrMissingCapture`, que dejaría sin cálculo a todo préstamo nuevo.
4. `S3-iva-lectura` pregunta 2: escribí `ErrMissingCapture` como lectura PROVISIONAL porque es la que no inventa tasa (RN-26). Nada está en uso para ese caso.
5. La revisión pedía `null` para las 14 fixtures, pero `s4-kcal-recetas-autor` y `s4-kcal-sin-dato` dependen de Q2-90, que no tiene recomendación ni respuesta (K46: quedan PROVISIONALES con bloqueo abierto). Por eso citan `S4-empates-sin-confirmar`. `s5-profeco-muestra` cita `S5-profeco-filtro` y deja de mencionar las ligas: `vinculos.csv` no existe y quién confirma ya es FIRME.
6. S3 sube a 1.5 porque cambia la interfaz (`Credit`). S4 y S5 se quedan en 1.5.

## Bloqueos
- Nuevos, ABIERTOS: `S3-formula-tarjeta-pdf` (espera los PDF de Axel), `S3-iva-lectura` (2 preguntas) y `S4-empates-sin-confirmar` (Q2-90).
- A-2 sin resolver: `make verify-all` sigue roto en `verify-c1` hasta que alguien con alcance en C1 aplique la solicitud 1 de `S1-puerta2.md`.

## Pendientes
- Lo de «Solicitudes fuera de alcance» (integrador o agentes de C1 y de bloqueos).
- Los cambios de motor, oráculo y fixtures de abajo (WF-3).
- Ronda 2 del revisor.

## Cambios que deben implementar los motores (complemento de `contratos-puerta2.md`)
**S3 (`backend/internal/motor_cred`) y su oráculo (`oraculo/cred_*.py`)**
- `Credit.VATIncluded *bool` (`vat_included`).
- `CardPeriodInterest`, `CardMinimumPayment` y la regla (4) de la tienda: `ErrMissingCapture` si falta cualquiera de los 5 datos (con `VATIncluded` nil incluido) o si `VATIncluded` = false llega sin `VATOnInterestBP`. El IVA se suma solo con false.
- `QuotePurchase` de `term_loan`:
  - true → `vatBP` nil;
  - false → `VATOnInterestBP`, y sin él, `ErrMissingCapture`;
  - nil → `vatBP` nil y `Provisional` = [`S3-iva-lectura`].
- Listas `Provisional`:
  - `S3-formula-tarjeta` → `S3-formula-tarjeta-pdf` en interés con fórmula, mínimo y corte `computed`;
  - fuera `U4-iva-credito` en todas;
  - `amortizar` siempre `[]`.
- Fixtures:
  - `s3-tarjeta-minimo` y `s3-tienda-abonos-interes`: `provisional` = `{"bloqueo": "docs/estado/bloqueos/S3-formula-tarjeta-pdf.md", …}`, sin «otros bloqueos».
  - `vat_included: false` en los créditos que traen IVA 1600: `s3-tarjeta-minimo` id 2, `s3-tienda-abonos-interes` id 6, `s3-tarjeta-sin-intereses` id 2 y `s3-unificado`, `-sin-corte` y `-duplicado` ids 2 (tarjeta) y 3 (préstamo). El CHECK de la 00004 no admite porcentaje con NULL ni con 1.
  - `null` en los que no tienen datos.
  - `_datos/kueski-semilla.json`: `vat_included: true`.
  - `s3-prestamo-frances`: `provisional` null.
- Pruebas `TestTarjetaSinCaptura`, `TestDespachoCotizacion` y `TestProvisionalListas` según §11.
- Comentarios PROVISIONAL de bloqueos resueltos, a quitar o pasar a FIRME:
  - `loan.go` l. 6 (`S3-convencion-prestamo`) y l. 35 (`U4-iva-credito`; además pasa a leer `VATIncluded`);
  - `revolving.go` l. 36 (`S3-serie-sin-intereses`);
  - `comparador/request.go` l. 27 y `comparador/generate.go` l. 46 (`S3-objetivos-medida`);
  - `oraculo/cred_genera_fixtures.py` l. 385 y 543 (`S3-serie-sin-intereses`).

**S4 (`backend/internal/recetas`) y su oráculo**
- `provisional` null en `s4-canal-mercados`, `s4-lista-reposicion-ritmo` y `s4-salud-imc-franja`.
- `s4-kcal-recetas-autor` y `s4-kcal-sin-dato`: `bloqueo` = `docs/estado/bloqueos/S4-empates-sin-confirmar.md`, con el mismo `que`.
- Comentarios:
  - `plan.go` l. 66 y `shopping.go` l. 113 (`S4-reglas-lista`);
  - `shopping.go` l. 243 (`S4-asignacion-mercado`);
  - `oraculo/rec_plan.py` l. 3, 151 y 213;
  - `oraculo/rec_genera_fixtures.py` l. 225, 235, 395 y 644 (textos de `descripcion` de las 3 fixtures anteriores).

**S5 (`backend/internal/precios`) y su oráculo**
- `provisional` null en `s5-canasta-regla-unica`, `s5-cotizacion-recomendada`, `s5-mejor-tienda-empate-orden`, `s5-mejor-tienda-empate-recomendados`, `s5-mejor-tienda-empate-viejos`, `s5-mejor-tienda-incompleta`, `s5-membresia-anual` y `s5-promocion-vigencia`.
- `s5-profeco-muestra`: `bloqueo` = `docs/estado/bloqueos/S5-profeco-filtro.md`, `que` = «conteos con el filtro y la liga normalizada en uso (A, Q2-56 a Q2-58 sin respuesta)».
- Comentarios:
  - `pricing.go` l. 227 y `basket.go` l. 21 y 59 (`S5-reglas-comparativa`);
  - `profeco/sqlite.go` l. 109 (`S5-profeco-almacen`);
  - `oraculo/pre_genera_fixtures.py` l. 254.
- Para no omitir comentarios de varias líneas: `grep -rn PROVISIONAL backend/internal oraculo/*.py | grep -E "S3-(convencion-prestamo|objetivos-medida|origen-cargo-normal|serie-sin-intereses)|U4-iva-credito|S4-(asignacion-mercado|confirmacion-empates|reglas-lista)|S5-(reglas-comparativa|profeco-almacen)"` y revisar a mano los bloques de comentario vecinos.
- Al cambiar los `provisional`, se regeneran `oraculo/salidas/s3.json`, `s4.json` y `s5.json`.

## Solicitudes a C1
- `Credit.vat_included` (booleano que admite null) y `vat` en el enum de `Credit.missing_terms`, en el orden de RN-26: `annual_rate`, `rate_base`, `vat`, `min_payment_pct_bp`, `min_payment_fixed_cents`. Ya está en la solicitud 1 de `S1-puerta2.md`.
- Si Axel contesta A en la pregunta 2 de `S3-iva-lectura`, `vat` también va con `vat_included` = 0 sin porcentaje.

## Solicitudes fuera de alcance
1. **A-2 (agente de C1 o integrador):** enseñar a `backend/internal/api/c1gen/gen2_schemas.py` (`_pv`) el tipo «hora HH:MM» (string con patrón `^([01][0-9]|2[0-3]):[0-5][0-9]$`) y aplicar el resto de la solicitud 1 de `S1-puerta2.md`. Después, `build.py --check` y `make verify-all`.
2. **Integrador o agente de bloqueos:**
   - Dar de alta en `docs/estado/bloqueos/INDICE.md` los 3 bloqueos nuevos ABIERTOS (`S3-formula-tarjeta-pdf`, `S3-iva-lectura` y `S4-empates-sin-confirmar`) y ajustar el conteo de K46 en `docs/plan/decisiones.md`: P2c abrió 4 + 3 = 7 bloqueos nuevos.
   - Poner una línea de remisión en `S3-formula-tarjeta.md` (preguntas 1 a 3 → `S3-formula-tarjeta-pdf`), en `S4-confirmacion-empates.md` (→ `S4-empates-sin-confirmar`) y en `U4-iva-credito.md` (casos no cubiertos → `S3-iva-lectura`).
3. **Integrador:** unir la lista «Cambios que deben implementar los motores» de este reporte a la de `contratos-puerta2.md` (o a la de INTEGRACION de P2c).
4. **U4 o integrador:** `docs/design/U4-ajustes-a.md` P-26.05 sigue con la marca PROVISIONAL (`U4-iva-credito`) y con «vacío = incluye». Con Q2-67 = 1, el campo se guarda en `vat_included`.

## Propuestas para Axel
- Contestar por id `S3-iva-lectura` (préstamo sin el dato del IVA; «no incluye IVA» sin porcentaje) y Q2-90 (`S4-empates-sin-confirmar`).
- Entregar los PDF de un estado de cuenta de tarjeta y uno de crédito de tienda (`S3-formula-tarjeta-pdf`).

## Consumo aproximado
40 acciones de herramienta · unos 190 KB leídos · unos 50 KB escritos (28 KB de renglones cambiados en los contratos, 8 KB en bloqueos y 13 KB en este reporte).
