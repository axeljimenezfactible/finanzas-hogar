# REVISIÓN motor S5 · P2c · ronda 1 de 3 (revisor-motor-S5-r1) · 11-oct-2026, noche

Revisor independiente `revisor-motor-S5-r1`, flujo P2c «Aplicación de la puerta 2». No escribí nada de lo que reviso y no corregí nada. Evalué solo contra `docs/criterios/S5.md` 1.2 (completo) y `docs/contratos/motores/S5-interfaz.md` 1.5, con las respuestas de Axel (`docs/estado/puerta2/RESPUESTAS.md`; recomendadas de `PUERTA-2.md`). Objeto: el trabajo de `motor-S5-puerta2`, es decir el diff de `backend/internal/precios/` (con `profeco/`), `backend/cmd/motores/s5.go`, `oraculo/pre_calculo.py`, `oraculo/pre_profeco.py`, `oraculo/pre_genera_fixtures.py`, las fixtures s5 nuevas y regeneradas y el reporte `docs/estado/p2c/motor-S5-puerta2.md`. Volví a correr cada verificación con mis comandos. Los casos propios están en `datos/tmp/rev-s5-p2c-r1/` (scripts, raíz aislada con copia del oráculo y fixtures propias con esperados hechos a mano con `Fraction`); los borré al terminar.

## Veredicto
**NO APROBADO: 0 altos, 3 medios y 3 bajos.**
- Lo que pedía la ronda sí cumple:
  - Un solo significado de la mejor tienda en los dos modos (Q2-05 y Q2-08 = A).
  - Importe a granel con un solo redondeo, PROVISIONAL `S5-carritos-granel`.
  - Lectura PROFECO con la gramática y la lista cerrada de §11 (Q2-04 = A).
  - Orden de las listas de §15.2, y las marcas FIRME con su Q2 (Q2-47, Q2-53, Q2-55 y Q2-59 a Q2-63, todas recomendadas).
  - Canon intacto (+$138.50 y +$100.50), `data.json` intacto, cobertura ≥ 85 % y `make verify-all` OK.
- Ninguna decisión «sin recomendación» (Q2-48, Q2-50 a Q2-52, Q2-56 a Q2-58) se tomó: siguen PROVISIONAL con la opción en uso.
- Quedan abiertos tres puntos medios:
  - Con renglones a granel partidos, el motor y el oráculo dan `piece_count` distinto. Nadie lo había documentado (M-1).
  - Hay divergencias entre el motor y el oráculo, ya conocidas o copiadas de un lado al otro, que no tienen bloqueo ni solicitud (M-2).
  - La marca PROVISIONAL de «sin costo capturado» se quitó y el criterio S5-09 la pide (M-3).

## Verificaciones (evidencia mía)

| # | Verificación | Resultado | Evidencia |
|---|---|---|---|
| a1 | `go vet ./internal/precios/... ./cmd/motores/` | OK | sin salida |
| a2 | `gofmt -l internal/precios cmd/motores` | vacío | `gofmt:[]` |
| a3 | Cobertura (K25 ≥ 85 %) | CUMPLE | `go test -count=1 -coverprofile`: precios 96.0 %, asistida 100.0 %, ocr 99.3 %, profeco 95.3 %; perfil total 88.8 %. El reporte da 95.9 % en precios (diferencia de una décima, sin efecto) |
| a4 | Pruebas por criterio (`-run`) | 0 FAIL, 0 SKIP | Comparativa 4, Unidad 2, MejorTienda 5, Frescura 3, Fuente 4, Profeco 12, LecturaAsistida 1, Membresia 2, Mercado 2, PreciosViejos 1, EvalCarritos 3 y ParsePriceMotivos 1, todas PASS |
| a5 | `float` en dinero (S5-13) | CUMPLE | `grep float32\|float64` fuera de pruebas: solo `ocr/tsv.go` (confianza del OCR) |
| a6 | Red (S5-08) | CUMPLE | `grep net/http\|http.Get\|http.Client` fuera de pruebas: 0; dominios de tiendas en `.go`: 0 |
| b1 | Fixtures reproducibles | CUMPLE | copié `oraculo/fixtures/s5`, corrí `python3 -I oraculo/pre_genera_fixtures.py` («32 fixtures») y `diff -r` dio 0 diferencias. El generador solo importa la biblioteca estándar, no el oráculo |
| b2 | `run_fixtures` reproduce `salidas/s5.json` | CUMPLE | `pre_run_fixtures.py --out datos/tmp/…`: «32 fixtures, 32 coinciden, 1336 hojas de esperado, 0 diferencias»; `cmp` con `oraculo/salidas/s5.json`: idéntica |
| b3 | `compara.py s5` con motor propio | 0 diferencias | Construí el binario y corrí `motores s5 dump`. Oráculo contra esperado, motor contra esperado y oráculo contra motor: 935 de 935 hojas cada uno; provisionales: 315 de 315 |
| b4 | Motor contra oráculo, salida completa sin recortar (las 32 fixtures con `esperado` vacío) | 1 clave fuera de §15.2 | La única diferencia es `membership.note` («sin costo capturado»), que solo emite el motor, en `s5-membresia-anual` y `s5-membresia-equilibrio`. §15.2 no la lista |
| c1 | Granel del revisor r3 (0.5 a 1.01 y 1.03) | CUMPLE | Mi `s5-rv-granel` q0: totales 101 y 103, líneas mostradas 51 y 52, unidad base 135, 137, 721 y 736. Coinciden el motor, el oráculo y mi esperado |
| c2 | Granel propio | CUMPLE en las cifras de §15.2, salvo `piece_count` (M-1) | q1, excedente a granel: Bodega 0.5 de arroz y 0.5 + 0.25 de atún; total 126.25 → 126, idénticos 101, Walmart 103, Walmart primero y Bodega +23. q2, sustituto a granel: 0.333 × 750 g a 21.00 contra 0.25 × 1 kg a 28.00; 2800 = 2800 exacto, sin ganador, totales 699 y 700. q3, tres canastas: `tercera` incompleta y con el menor total (50) va al final con `null`; el atún queda como sustituto. Las tres coinciden con mi esperado |
| c3 | `profeco_import`, CSV propio de 23 registros: los 17 del tipo del revisor (precio 0, «1,021.00», cadena vacía, ATUN/ATÚN y «Wal-Mart»/«Gustavo A. Madero» normalizados) más 6 de la lista cerrada | CUMPLE | Producción: read 23, rejected 14, filtered_out 2, skipped 2, loaded 4, loaded_proposed 0 y duplicates 1. Revisión: 14, 2, 1, 4, 1 y 1. Se cumple la identidad de §11. `coverage` lleva el texto tal como llega (3 filas) y los 14 motivos con su texto exacto. El motor, el oráculo y mi esperado coinciden |
| c4 | `difference_vs_cheapest_cents` con un solo significado | CUMPLE | `compare_carts` (q3) y `compare_basket`: con una incompleta más barata, `s5-rv-canasta` da s2 0, s1 +50 y s3 `null`, y `cheapest_store` s2. Si ninguna está completa, `s5-rv-canasta-ninguna` da todo `null` y `cheapest_store` `null`. Con canasta vacía, `s5-rv-canasta-vacia` da s1 y s2 con 0 y `cheapest_store` s1 (letra de §7.4, P-1). El tipo `*int64` es el mismo en los dos modos. Mis 5 fixtures: 295 de 295 hojas en los tres sentidos |
| c5 | Canon +138.50 y +100.50 | INTACTO | Motor (`s5-comparativa-canon`): Bodega 309450 y Walmart 323300 (diferencia 13850); idénticos 277350 y 287400 (diferencia 10050); 33/3 renglones; piezas 61/60/55; ganadores 22/3 y empates 8. Desde `data.json` con `Decimal`: 138.5 y 100.5; 33 y 3 |
| c6 | Orden de `stores[]` en `s5-carritos-presentaciones` (HECHO 8 del reporte) | CUMPLE | Con mis conjuntos de productos: Bodega completa; Walmart sin leche (missing 1); tercera sin crema ni leche (missing 2). El orden es Bodega 0, tercera `null` y walmart `null`, como dice el reporte; el resumen de `contratos-puerta2.md` estaba incompleto |
| d1 | `data.json` intacto | CUMPLE | sha256 `2fe7afae…1d31b`; `git diff HEAD` vacío; `model.json` y `verify3.json` = `MODEL_SHA` y `VERIFY3_SHA` |
| d2 | `make verify-all GOOSE=…/goose` | EXIT=0 | «Verificación ampliada OK». S5: «32 fixtures, 32 coinciden… salidas/s5.json idéntica», «935 de 935» y K25 OK en los 4 paquetes de precios |
| e1 | S5-05 sin cambio de cifras | CUMPLE | `sqlite3 -readonly qqp.sqlite`: 37,548,200 registros, 0 precios con coma y 0 con coma en los decimales, `MAX(fecha)` 2026-08-31 |
| e2 | S5-06 y S5-07 (fuera del diff de este motor) | sin cambios del motor | `datos/semilla/tickets/`: 0 fotos, así que S5-06 queda «sin medir». `docs/estado/S5-terminos.md` existe y `grep -c Veredicto` da 3; es obra de S5-terminos, no la revisé. `tesseract 5.5.3` |
| e3 | Marcas FIRME y PROVISIONAL contra la puerta 2 | CUMPLE | Q2-47, Q2-53, Q2-55, Q2-59, Q2-60, Q2-61 y Q2-63 tienen recomendada A y Q2-62 tiene B (`PUERTA-2.md` §7.2: «A en 1 a 3 y 4b, B en 4a»). Las fixtures llevan `provisional` solo con bloqueos ABIERTOS: `S5-carritos-granel`, `S5-carritos-detalle` y `S5-profeco-filtro` (4). Las de `S5-reglas-comparativa` y `S5-profeco-almacen` (RESUELTOS) van con `null` |

## Hallazgos

**M-1 [media] `piece_count` a granel: el motor y el oráculo dan cifras distintas, y no está documentado ni tiene fixture.**
- Rutas:
  - `backend/internal/precios/compare.go:99-104` (`pieces`, sobre la presentación agrupada del paso 1).
  - `oraculo/pre_calculo.py:244` (`_piezas`, sobre cada línea de `rows[]`: emparejada y excedente).
  - `S5-interfaz.md` §7.2 paso 6 (supuesto S-11: «un renglón con cantidad no entera cuenta como 1 pieza», sin decir qué renglón).
- Evidencia:
  - `s5-rv-granel` q1: Bodega captura 0.5 de arroz y 0.5 + 0.25 de atún. Motor **2**, oráculo **3**.
  - `s5-rv-piezas`: Walmart captura 0.5 de atún y 1 bolsa de arroz; se emparejan 0.5 y sobra 0.5. Motor **2**, oráculo **3**.
  - Si «renglón» es la línea capturada, Bodega en `s5-rv-piezas` (atún 0.25 + 0.25 y arroz 0.5) daría 3, y el motor y el oráculo dan 2.
  - Hay tres lecturas posibles y no hay regla. Es la misma clase de hueco que H-r3-02. La fixture `s5-carritos-granel` no tiene renglones partidos ni excedente a granel.
- Ninguna cifra canon cambia.
- **Corrección:**
  - Solicitud a contratos para §7.2 paso 6: decir a qué renglón se refiere S-11 (la línea capturada, la presentación agrupada del paso 1 o la línea de `rows[]`). Si no es técnica, va a `S5-carritos-granel` como pregunta nueva.
  - Alinear el lado que no cumpla.
  - Agregar a `s5-carritos-granel` una consulta con un renglón partido y un excedente a granel.

**M-2 [media] Hay divergencias entre el motor y el oráculo y supuestos copiados del motor que no tienen bloqueo ni solicitud: el siguiente implementador tendría que adivinar.**
El reporte las lista en «Pendientes» y «Supuestos» (bien separadas), pero no abre bloqueos ni las deja como solicitudes fuera de alcance. Comprobé las dos primeras con mi `s5-rv-pendientes`:
- (a) **Misma presentación con precios distintos en una canasta** (§7.2 pasos 1 y 2 no dicen qué precio toma la parte emparejada). Bodega trae 1 atún a 20.00 y otro a 22.00; Walmart, 1 a 21.00.
  - Motor: `identical_total_cents` 2000, ganador Bodega, `ties` 0 y `winners` Bodega 1.
  - Oráculo: 2100, sin ganador, `ties` 1 y `winners` 0; las líneas valen 2100/2100 contra 2000/2200.
  - Toca cifras que se ven en pantalla, como el desglose de idénticos y los ganadores.
- (b) **`current_price.references[]`** (§6): con 2 de PROFECO y 1 `site_assisted`, el motor da 2 (la última de cada fuente) y el oráculo 3 (todas). El orden de §15.2 sí es el mismo en los dos.
- (c) Documentadas por el implementador y no reproducidas por mí:
  - Formatos de fecha aceptados distintos.
  - Una liga a `presentacion_id` desconocido: `unknown_presentation` en el motor, registro omitido en el oráculo.
  - Tope de 100 `rejections[]` solo en el motor.
- (d) S-1, S-2 y S-3 (espacios dentro del precio, más de 15 dígitos enteros, «1.2.3» como «más de dos decimales») y «liga repetida → `invalid_argument`»:
  - El oráculo copia al motor. `oraculo/pre_profeco.py` lo dice: «supuesto del agente motor-S5-puerta2».
  - Por eso que coincidan no es evidencia independiente (§15.3).
- **Corrección** (el corrector o el integrador):
  - Llevarlas como solicitudes a contratos: §7.2 pasos 1 y 2 para (a), §6 para (b) y §11 para (c) y (d).
  - Las que no sean técnicas abren bloqueo o se suman a `S5-carritos-detalle` y `S5-profeco-filtro`, como PROVISIONAL «en uso: el motor».
  - Con la regla fijada, una fixture por caso. No se elige aquí qué lado tiene razón.

**M-3 [media] S5-09: se quitó la marca PROVISIONAL de «sin costo capturado».**
- `backend/internal/precios/extras.go:88` decía «sin costo capturado (PROVISIONAL, S5-09)» y ahora dice «sin costo capturado».
- El criterio S5-09 (alta) pide, al verificar, «con costo ausente, salida «sin costo capturado» PROVISIONAL».
- §8 del contrato 1.5 no trae la marca. Es una inconsistencia entre el criterio y el contrato, y el motor la resolvió por su cuenta.
- `TestMembresiaSinCosto` solo comprueba `verdict = missing_fee`.
- **Corrección:** devolver la marca a `Note` y hacer que la prueba verifique el texto. Si se quiere quitar, primero el integrador alinea S5-09 y §8, y después se cambia el motor.

**B-1 [baja] Independencia según la letra de §15.3.**
- §15.3 dice que el oráculo «lo escribe un agente distinto del que implementa `internal/precios`».
- En P2c, `motor-S5-puerta2` cambió el motor, el oráculo (`pre_calculo.py`, `pre_profeco.py`) y los esperados (`pre_genera_fixtures.py`). Así lo asignó el reparto de la ronda, y WF-1 (`correccion-S5-r2`) hizo lo mismo.
- Lo que sigue intacto: el oráculo no importa Go, el motor no lee salidas del oráculo y el generador solo usa la biblioteca estándar.
- Lo que cambió el contrato (granel, §7.4 y §11) lo recalculé a mano con casos propios: 295 de 295 hojas.
- Lo copiado del motor está en M-2(d).
- **Corrección:** el integrador registra la excepción del reparto de P2c, o pide que §15.3 diga quién alinea el oráculo cuando cambia el contrato.

**B-2 [baja] Faltan marcas PROVISIONAL en el código.**
- La regla del importe a granel (`compare.go`, alrededor de las líneas 160 y 309; §7.2 paso 8, `S5-carritos-granel`) no lleva marca.
- La liga y el filtro normalizados (`profeco/links.go:45-52` y `profeco/filter.go:35`; `S5-profeco-filtro`, Q2-56 a Q2-58) tampoco.
- El contrato y las fixtures sí las llevan.
- **Corrección:** poner el comentario «PROVISIONAL (bloqueo …, en uso A)» en esos tres puntos.

**B-3 [baja] Hay artefactos fuera del alcance del motor que no se actualizaron.**
- La tabla de casos de §15 del contrato no lista `s5-carritos-granel`, `s5-profeco-bordes` ni `s5-profeco-motivos`.
- En `docs/estado/bloqueos/S5-carritos-granel.md`, «Cómo seguí» sigue diciendo que falta alinear el oráculo y agregar la fixture a granel, y las dos cosas ya están hechas.
- **Corrección:** solicitudes al integrador para contratos y bloqueos.

**Observación fuera del objeto (sin severidad aquí).** El extracto C1 `compras-tiendas-a.yaml:998` (y `reportes.yaml:802`) declara la clave `cheapest_store_id_`, con guion bajo al final. Sale de `gen2_schemas.py:420` (`_ins(…, cheapest_store_id_=…)`). Ya lo registra `REVISION-C1-r1.md`. Lo cito porque es la salida de S5 en la API.

## Estado de los hallazgos anteriores (WF-1, `REVISION-S5-r3.md`, y M-r5-3)

| Hallazgo | Estado | Evidencia |
|---|---|---|
| H-r3-01 [media] contrato desfasado | CERRADO (contratos, no el motor) | `S5-interfaz.md` 1.5 trae §7.2 pasos 2, 3 y 8, §7.4, §11 y §15.2 |
| H-r3-02 [media] granel | CERRADO en totales, idénticos, líneas y unidad base | c1 y c2: el motor, el oráculo y mi esperado coinciden. Queda el hueco nuevo de `piece_count` (M-1) |
| H-r3-03 [media] lectura PROFECO | CERRADO | c3: 23 registros, mismos conteos y motivos en producción y revisión. Lo copiado del motor está en M-2(d) |
| H-r3-04 [media] `difference_vs_cheapest_cents` | CERRADO | c4 en los dos modos, con `null` en las incompletas |
| H-r3-05 [baja] comando de Tesseract | Fuera del objeto | `ocr/` no se tocó |
| H-r3-06 [baja] etiquetas del corte en `S5-cobertura-profeco.md` | ABIERTO, fuera del alcance del motor | el reporte lo deja así |
| M-r5-3 (Q2-08) orden de listas | CERRADO | `references[]` por fecha descendente y fuente (mismo orden en los dos lados, M-2(b)); `stale_prices_for_plan` sin repetidos; b3 y b4 sin diferencias de orden |
