# Corrección del contrato S3 · ronda 2 (agente correccion-contrato-S3-r2)

11-oct-2026 · WF-1 · Objeto: `docs/contratos/motores/S3-interfaz.md` v1.1 → **v1.2** (90,865 → 115,607 bytes) · Informe atendido: `docs/estado/wf1/REVISION-contrato-S3-r2.md` (4 hallazgos medios M1–M4; 0 altos). Vara: `docs/criterios/S3.md` 1.2, decisiones §C y §K, `oraculo/fixtures/FORMATO.md` 1.1.

## Hecho
1. **M1 (S3-01, S3-02, S3-09), contrato alineado a FORMATO 1.1:**
   - (a) Vista nueva `vista_dos_quincenas` (§9.3) con las claves literales de `data.json:financiamiento.dos_quincenas`: `total`, `intereses`, `tasa_pct` = `{num, den, decimales}` y `<persona>.{intereses, total, por_quincena}`. En `s3-kueski-despensa-1` (§10), R-05 y R-06 van con `id` `data.json:financiamiento.dos_quincenas` y `fuente` `data.json`. La compra es d1 a 2 q **sin** cotización: la tasa sale de la tabla, con sus `decimales` 2 capturados, porque `QuotedRate` no tiene de dónde sacar los decimales. `dos-quincenas` queda como cruce libre.
   - (b) La clave JSON de la tasa pasa a `decimales` en §3, en §9.1, en el ejemplo 1 y en la semilla. En Go, el campo `Decimals` lleva la etiqueta `json:"decimales"`.
   - (c) §2 y §9.5: S3 agrega solo `backend/cmd/motores/s3.go` con `func init() { register("s3", evalS3) }`, con `evalS3` de tipo `Evaluator`, y no edita `main.go`. El despacho vive en un paquete nuevo, `backend/internal/motor_cred/evaluar` (`Evaluate`); así no hay ciclo con `comparador` y se prueba sin el binario. Una `funcion` desconocida da `invalid`. Se cita FORMATO 1.1.
2. **M2 (S3-05, S3-03, S3-12), abonos registrados en la línea de tiempo (§5.2):**
   - Hay un evento `payment` por abono. `PrincipalRestored` e `InterestNotRestored` valen 0, salvo en el abono que completa el pago, que lleva K e I completos.
   - Si el pago está parcial y su vencimiento es ≥ `Today`, el resto se proyecta como un solo evento en el vencimiento, con `Amount` = M − S, K e I. Si ya venció sin completarse, no repone.
   - Quedan fijados la fecha de reposición, el orden entre eventos del mismo pago, `PurchaseRef` "" y `Seq` 0 en el inicio y el contado, y las invariantes.
   - El ejemplo 1 de §10.3 trae ahora su `esperado` completo: 4 eventos, series, mínimos y bandera, con las cuentas a mano.
3. **M3 (S3-10, S3-12, S3-06, S3-11), tabla de despacho de `QuotePurchase` (§5.8 nuevo):** por `Kind`, fija entradas obligatorias, intereses en orden, serie, `Source`, `PurchaseKind` (campo nuevo, igual a `credit_purchase.kind`), `RateBasis`, `Rate` y la precedencia de errores.
   - (a) **Tienda sin interés:** se captura como `TermRate` con tasa 0 (el esquema ya lo admite). «Sin datos» da `missing_capture`, que gana sobre `no_rate` (la tienda nunca da `no_rate`). Una tienda que solo tiene datos K18 da `ErrUnsupported`, PROVISIONAL por `S3-formula-tarjeta` (3). Sus intereses y su mínimo por periodo salen con `interes_tarjeta`/`minimo_tarjeta`, que es el algoritmo común con el oráculo.
   - (b) §8.2: el cargo normal de la tarjeta (plazo 1, 0.00) **siempre se genera**.
   - (c) **Préstamo:** se calcula con `Amortize` más las partes de §4. Sin `AnnualRate` da `ErrMissingCapture`. La consulta `amortizar` devuelve las partes vacías. Se fijan `TotalExact` y `ReferencePayment`.
   - (d) **`Source`:** los abonos sin interés llevan `term_table` y el MSI `promotion`. El cargo normal lleva `quoted` como PROVISIONAL, con el bloqueo nuevo `S3-origen-cargo-normal`, porque el CHECK de `interest_source` no tiene un valor que le corresponda.
   - Se ajustaron §5.3, §5.4 y §5.5.
4. **M4 (S3-20, S3-12, S3-13), marcas PROVISIONAL (§14.1 nuevo):**
   - Regla de las listas: sin repetidos y en orden alfabético por bytes. Cada salida tiene su contenido: `Purchase`, `CalendarEntry`, interés/mínimo de tarjeta, `Scenario` (siempre `S3-objetivos-medida`, más lo de su compra y lo de la serie del objetivo 2) y `Result` (unión de sus escenarios más `S3-objetivos-medida`; `[]` con `missing_constraints`).
   - Regla de las fixtures: un objeto con el bloqueo principal y «otros bloqueos» dentro de `que`, porque FORMATO 1.1 solo admite un objeto.
   - Llevan marca `s3-comparador-caso-1…-3` y `-carga` (el caso 2 también nombra `S3-formula-tarjeta` y `S3-origen-cargo-normal`), además de `s3-unificado-sin-corte` y `s3-unificado`.
   - §11 tiene pruebas nuevas: `TestReposicionAbonoParcial`, `TestTiendaConInteres`, `TestDespachoCotizacion`, `TestProvisionalListas` y `TestEvaluar`. §13 está actualizado.
5. §14: se agregaron un **Hallazgo 5** nuevo (no corregido, ver Pendientes) y el bloque «Cambios de la versión 1.2». Ninguna cifra canon cambia.

## Archivos
- Modificado: `/home/deepsec/finanzas-hogar/docs/contratos/motores/S3-interfaz.md` (v1.2).
- Creado: `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S3-origen-cargo-normal.md`.
- Creado: `/home/deepsec/finanzas-hogar/docs/estado/wf1/correccion-contrato-S3-r2.md` (este reporte).
- Temporales en `datos/tmp/corr-s3-r2/` (scripts de edición `p1.py`…`p3.py`, `verifica.py` y una copia de la v1.1): borrados al terminar.

## Verificación
- El artefacto es un contrato sin código. `backend/internal/motor_cred`, `oraculo/cred_*.py` y `oraculo/compara.py` todavía no existen, así que no hay pruebas Go, oráculo ni `compara.py` que correr. La verificación es la consistencia del contrato con scripts propios (`Fraction`, `python3 -I`).
- `python3 -I datos/tmp/corr-s3-r2/verifica.py` → `bloques json 2` · `comprobaciones 36 fallas 0`. Cubre:
  - los 2 bloques JSON del contrato son JSON válido;
  - **ejemplo 1 (M2):** recalcula la línea de tiempo con una implementación propia de la regla nueva de §5.2 (serie de d1 `SplitTrunc`/`SplitHalfUp`, abonos ordenados, abono que completa, resto proyectado) y coincide con el `esperado` escrito en `events` (4), `series_a_cents`, `series_b_cents`, `min_a_cents` y `min_b_cents`; d1 da 97375 ×4 y capital 77363 ×3 + 77361, y Karen da 37493 ×4;
  - **R-05/R-06 (M1):** las 9 hojas de `data.json:financiamiento.dos_quincenas` salen con las reglas de §4. Total 362900, intereses 53450, Axel 32870/223170/111585, Karen 20580/139730/69865 y `tasa_pct` = HALF_UP(100 × 53450/309450, 2) = 17.27. Las 8 cifras aparecen en el contrato;
  - el orden alfabético de los 4 nombres de bloqueo y la existencia de sus archivos;
  - ninguna clave `decimals` queda en el JSON; «el integrador une» ya no aparece; solo queda 1 cita de «FORMATO.md v1.0» (histórica, en «Cambios de la versión 1.1»); existen §5.8, §14.1, §9.5 con `register("s3", evalS3)`, el `id` `data.json:` y `invalid`; los `id` libres casan con `^[a-z0-9][a-z0-9_-]*$`;
  - el sha256 de `oraculo/model.json` es `1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4` (intacto).
- `go build ./... && go vet ./...` en `backend` limpio. No toqué código; queda como constancia de que el módulo `finanzas-hogar` compila. La ruta de import de `s3.go` (`finanzas-hogar/internal/motor_cred/evaluar`) usa ese módulo.
- Criterios de S3 tocados y su evidencia en el contrato:

| Criterio | Evidencia |
|---|---|
| S3-01 | §9.3 `vista_dos_quincenas`; §10 R-05/R-06 con `id` de `data.json` |
| S3-02 | La línea `data.json` de `compara.py s3` deja de salir «0 de 0» |
| S3-03, S3-05 | §5.2 y el ejemplo 1 completo, verificado |
| S3-06, S3-08, S3-09, S3-10 | §5.8 |
| S3-11 | Préstamo con partes en el calendario |
| S3-12 | §8.2: el cargo normal siempre se genera; marcas del comparador |
| S3-13 | Fixtures del comparador marcadas |
| S3-20 | Bloqueo nuevo y §14.1 |

## Supuestos
- La tienda sin interés se captura como `credit_term_rate` con tasa 0. El esquema lo admite y RN-08 ya dice que una compra planeada usa la tasa del plazo. No hay columna nueva.
- Una compra de tienda con intereses cotizados por el acreedor usa la regla (1) de la línea (RN-08, compra real). El valor sale de una captura, no se inventa.
- La tarjeta da `ErrUnsupported` con una promoción de tasa ≠ 0 y `ErrNoRate` con más de un pago sin promoción, porque no hay regla para meses con intereses.
- Al préstamo sin `AnnualRate` le toca `ErrMissingCapture`, porque la tasa es captura del producto, como en K18.
- Ninguno cambia cifras, esquema ni endpoints. Lo que sí cambiaría el esquema (`Source` del cargo normal) quedó como bloqueo.

## Bloqueos
- Nuevo: `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/S3-origen-cargo-normal.md`. Pregunta qué `interest_source` lleva un cargo normal de tarjeta. El motor emite `quoted` PROVISIONAL.
- Siguen abiertos y citados: `S3-formula-tarjeta` (ahora también cubre la cotización de una tienda con interés por datos K18), `S3-convencion-prestamo` y `S3-objetivos-medida`. `S1-parametros-avisos` no toca a S3.

## Pendientes
1. **Hallazgo 5, nuevo, no corregido.** Está fuera de M1–M4 y choca con K14. En la línea revolvente, una promoción al 0 % con los dos modos de RN-07 (§5.1) da montos 77362 ×3 + 77364 y capital 77363 ×3 + 77361: interés −1 en tres pagos. `scheduled_payment.interest_cents ≥ 0` (00001) lo rechaza. MSI y tienda sin interés ya usan capital = monto. Ruta: `S3-interfaz.md` §5.1 (Serie) y §14 (Hallazgo 5). Lo decide el revisor r3 o Axel.
2. Hallazgos bajos B1–B7 de la r2, sin corregir porque el encargo eran altos y medios. Rutas: B2 en §11.1; B3 en §5.3; B4 en §11 y §10.2; B5 en §5.2; B6 en §8.2; B7 en §8.5 (con solicitud previa a RN-11 y C1). B1 es de los archivos de bloqueo (solicitud abajo).

## Solicitudes fuera de alcance
1. `oraculo/fixtures/FORMATO.md` §2 (integrador). Cambio: que `provisional` admita una lista de objetos `{bloqueo, que}`. Motivo: M4; el caso 2 del comparador y las fixtures del calendario unificado dependen de 2 o 3 bloqueos, y hoy van con «otros bloqueos» dentro de `que`.
2. `docs/estado/bloqueos/INDICE.md` (orquestador). Cambio: registrar `S3-origen-cargo-normal` para la puerta 2. Motivo: bloqueo nuevo.
3. `docs/estado/bloqueos/S3-formula-tarjeta.md` (dueño del bloqueo). Cambio: que la pregunta 3 diga también «cómo se arman los abonos de una compra de tienda con interés por datos K18 (hoy `ErrUnsupported`, `S3-interfaz.md` §5.8)» y que se actualicen las referencias viejas (B1: §5.3/§5.5, `CardPeriodInterest`, `CardMinimumPayment`; en `S3-convencion-prestamo.md`, §5.4 `Amortize`). Motivo: M3 y B1.
4. U4 (pantalla Ajustes › Créditos) y C1 (operaciones de tasas por plazo). Cambio: aceptar tasa 0 por plazo en un crédito de tienda («abonos sin interés»). Motivo: §5.8; el esquema ya lo admite y la interfaz no debe rechazarlo.
5. `docs/contratos/diccionario-datos.md` y C1. Cambio: `Purchase.purchase_kind` ↔ `credit_purchase.kind`; cuando se resuelva `S3-origen-cargo-normal`, el valor de `interest_source` del cargo normal. Motivo: M3 (d).

## Paquetes instalados
Ninguno.

## Propuestas para Axel
1. **`S3-origen-cargo-normal`:** propongo la opción (a), `quoted`. El banco fija los intereses del cargo normal en el estado de cuenta (0.00 si se paga a tiempo). No toca el esquema y no cambia ninguna cifra. Si se prefiere un significado limpio, la (b) agrega un valor al CHECK con una migración nueva.
2. **Hallazgo 5:** cuando los intereses de una compra son 0, el capital de cada pago es igual a su monto en todos los tipos. Ya es así en MSI y en la tienda sin interés, cumple el CHECK del esquema y no cambia ninguna cifra canon ni los montos de S3-07.
3. **FORMATO:** `provisional` como lista (solicitud 1). No es un valor K31.

## Consumo aproximado
32 acciones de herramienta · ~205 KB leídos (protocolo, §K, informe r2, criterios S3 completos, FORMATO 1.1 §1–§6, contrato por partes ~85 KB, 3 bloqueos S3, esquema por grep, RN-07/08/26) · ~57 KB escritos (scripts de edición ~37 KB, verificación 5 KB, bloqueo 3 KB, reporte 10 KB; el contrato creció 24.7 KB).
