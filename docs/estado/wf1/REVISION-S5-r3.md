# REVISIÓN S5 · ronda 3 de 3 (WF-1)

Revisor `revisor-S5-r3` · 11-oct-2026 · independiente: no escribí nada de lo revisado y no corregí nada. Evalué solo contra `docs/criterios/S5.md` 1.2, `docs/plan/decisiones.md` §E y §K y `docs/contratos/motores/S5-interfaz.md` 1.4. Cada verificación la volví a correr con mis comandos; los cálculos numéricos los hice con scripts propios en `datos/tmp/rev5r3/` (borrados al terminar).

## Veredicto
**NO APROBADO.** Queda 1 hallazgo alto (S5-07, que solo Axel puede destrabar) y 4 medios: R2-04 sigue abierto y hay 3 nuevos que encontré con fixtures propias. Lo que pide el encargo sí cumple: canon, frescura K40/K37, PROFECO, OCR honesto, Sam's, mercados, cobertura ≥ 85 %, independencia y `data.json` intacto. Los hallazgos medios están en el contrato y en el oráculo; ninguno toca cifras canon ni conteos reales de PROFECO.

## Tabla de criterios (evidencia propia)

| Criterio | Estado | Evidencia (mis comandos) |
|---|---|---|
| S5-01 (alta) | CUMPLE | `canon.py` propio (Decimal) sobre `datos/semilla/data.json`: B 309450 · W 323300 · **+13850**; idénticos 277350 · 287400 · **+10050**; piezas 61/60/55; renglones 33/3; ganadores 22/8/3. `-run Comparativa`: 4 PASS (la prueba lee `data.json`, `precios_test.go:67`). `compara.py s5`: 747/747 en las tres comparaciones, 0 diferencias. La lectura +100.50 (mismos artículos, mismas cantidades) contra +138.50 (carrito contra carrito) queda explicada en §7.2, K-03 y K-06. |
| S5-02 (alta) | CUMPLE (motor) | Sustitutos con mi cálculo exacto: 29.33 vs 28.00 /kg; 126.92 vs 164.06 /kg; 33.75 vs 28.33 /l; mismos ganadores que las notas. `-run Unidad`: 2 PASS. Con cantidades a granel el oráculo redondea a medias (H-r3-02). |
| S5-03 (alta) | CUMPLE | `-run MejorTienda`: 5 PASS; en la canasta del 10-oct gana Bodega. Hay una inconsistencia de `difference_vs_cheapest_cents` entre modos (H-r3-04). |
| S5-04 (alta) | CUMPLE | `fresc.py` propio contra las 64 edades del motor (`s5-frescura-fronteras`): 0 diferencias en peso = (90 − d)/61, `is_stale` (d ≥ 30), `competes` (peso > 0) y pantalla HALF_UP a 4 decimales: 29 → 1, 30 → 0.9836, 60 → 0.4918, 90 → 0, 91 → no compite y marcado. `python3 -I` → `[0.9836, 0.4918, 0.0]`. `manual`, `ticket_manual` y `ticket_ocr` de la misma edad pesan 1 y compiten; `site_assisted` (10.00) y `profeco` (9.00) son más baratos que el ticket (29.00), quedan en `references[]` y `quote_for` da `null` con ellos. `-run Frescura` 3 PASS y `-run Fuente` 4 PASS. |
| S5-05 (alta) | CUMPLE en sustancia (H-r3-01 y H-r3-03 abiertos) | `cat *.csv \| wc -l` por corte: 14,889,243 / 13,270,548 / 9,388,473; `sqlite3 -readonly qqp.sqlite`: 14,889,219 / 13,270,524 / 9,388,457 (= wc − 24/24/16 encabezados), total 37,548,200, `MAX(fecha_registro)` = 2026/08/31. Los 3 RAR existen en `raw/<corte>-rar/` y git los ignora (`raw/.gitignore`). La muestra pesa 601 KB (≤ 2 MB). El cargador lee en flujo (`bufio` 1 MB + `encoding/csv`, sin `ReadAll`). `grep -c cadena_comercial` = 13; el documento cita `QQP_diccionario.csv` y dice con todas sus letras qué pasa con las 3 sucursales. Los conteos se reproducen con los comandos de su §5. |
| S5-06 (alta) | SIN MEDIR (declarado con honestidad) | `tesseract 5.5.3` con `spa`; `ls datos/semilla/tickets/ \| wc -l` = 0 (solo `.gitkeep`); `go run ./internal/precios/ocr/cmd/ocrmedir` → «S5-06 sin medir: hay 0 imágenes y se necesitan 5: faltan las fotos de Axel». `PricesFromOCR` devuelve `ErrNotConfirmed` y `-run OCRSinConfirmar` pasa. Ningún porcentaje se presenta como cumplido. |
| S5-07 (alta) | **NO CUMPLE** | `docs/estado/S5-terminos.md` no existe, así que el `grep -c Veredicto` no puede dar 3. `-run LecturaAsistida` pasa y el grep de dominios de tiendas en `.go` da 0. Bloqueo `S5-terminos` abierto (H-r3-00). |
| S5-08 (alta) | CUMPLE | `grep net/http\|http.Get\|http.Client` en `internal/precios` sin pruebas: 0. PROFECO se lee de archivos locales. |
| S5-09 (alta) | CUMPLE | Mi cálculo: 100000/12500 = 8; 100000/12000 = 8.33, y hacia arriba da 9; 13850 × 6 = 83100 < 100000, así que no conviene; con 80000 equilibra en 6 y conviene. Sin costo: `null` y «sin costo capturado (PROVISIONAL, S5-09)». `-run Membresia`: 2 PASS. |
| S5-10 (media) | CUMPLE | Las 5 fuentes, cada una marcada como «compite» o «referencia» (salida `s5-fuentes-compiten`); promoción vencida fuera; `-run Fuente` 4 PASS. |
| S5-11 (media) | CUMPLE | Tianguis de martes: del sábado 10-oct → 13-oct, del miércoles 14 → 20-oct; gana uva (6000 < 8700) y no compite en avena. `-run Mercado`: 2 PASS. |
| S5-12 (media) | CUMPLE | `-run PreciosViejos`: 1 PASS. |
| S5-13 (alta) | CUMPLE | `go vet` limpio; `gofmt -l` vacío; cobertura: precios 96.1 %, asistida 100 %, ocr 99.3 % (con `exec.go`), profeco 95.3 %, perfil total 88.8 % (89.6 % sin `exec.go` ni `cmd`); 78 PASS, 0 FAIL, 0 SKIP. `float64` solo en la confianza del OCR (`ocr/tsv.go`). |

**Independencia y canon (g):** `git status` muestra `datos/semilla/data.json` y `oraculo/model.json` sin cambios (sha256 de data.json `2fe7afae…1d31b`). Ni el motor lee salidas del oráculo ni el oráculo lee Go. `pre_verify.py` → «TODO COINCIDE». `pre_run_fixtures.py` → 29/29 y deja `oraculo/salidas/s5.json` idéntico byte a byte. Los cambios de `oraculo/model.py` y `verify3.py` son de S2-04, no de S5.

## Estado de los hallazgos de la ronda 2

| Hallazgo r2 | Estado | Evidencia |
|---|---|---|
| R2-01 [alta] S5-07 sin `S5-terminos.md` | **SIGUE ABIERTO** (H-r3-00) | El archivo no existe. El bloqueo `S5-terminos` tiene las 5 secciones del protocolo. Ningún agente puede cerrarlo sin leer términos de tiendas. |
| R2-02 [media] filtro de `profeco_import` | **CERRADO como PROVISIONAL** | El bloqueo `S5-profeco-filtro` (opción A) tiene formato correcto, y el motor y el oráculo coinciden en `s5-profeco-filtro` y en mis variantes: modo revisión, mismo CSV → 0 diferencias; CSV en Windows-1252 con «Juárez», «Peñón» y «Aurrerá» → 0 diferencias (en los dos casos solo sobra la clave `loaded_proposed` del motor). El caso H-01 (cadena con municipio sin sucursal) y «Wal-mart» = «WALMART» ya coinciden. Las reglas de lectura y de liga siguen divergiendo: es un hallazgo nuevo (H-r3-03). |
| R2-03 [media] `compare_carts` varias presentaciones y `missing_count` | **CERRADO como PROVISIONAL** | El bloqueo `S5-carritos-detalle` (opción A) existe, junto con la fixture `s5-carritos-presentaciones` y `TestEvalCarritosVariasPresentaciones` (PASS). Mi variante «dos de tres» da 0 diferencias. Quedan dos huecos nuevos: el granel (H-r3-02) y la diferencia contra una tienda incompleta (H-r3-04). |
| R2-04 [media] contrato 1.4 desfasado | **SIGUE ABIERTO** (H-r3-01) | `S5-interfaz.md` sigue en 1.4 sin las solicitudes 1–6 de `correccion-S5-r2`. |

## Hallazgos de esta ronda

**H-r3-00 [alta] S5-07 · `docs/estado/S5-terminos.md` no existe.** El criterio es alto y no tiene salida «sin medir», así que no cumple aunque el interruptor rechace la activación y no haya URLs de tiendas. Ruta: `docs/estado/S5-terminos.md` (falta) y `docs/estado/bloqueos/S5-terminos.md`. **Corrección:** el orquestador lleva el bloqueo a la puerta 2. Con la opción A o B se escriben los 3 veredictos leyendo solo páginas legales, con URL y fecha. Con la opción D el integrador reescribe S5-07 y el plan §3.4 para que «lectura asistida apagada sin veredicto» sea el estado aceptado en WF-1.

**H-r3-01 [media] Contrato S5 1.4 desfasado (R2-04 sin cerrar).** Hechos:
- **§11** sigue diciendo «2024 y 2025 en CSV», `raw/<corte>/`, extracción a `datos/tmp/` y base `datos/profeco/qqp.sqlite`; además no recoge la opción A de `S5-profeco-filtro` ni la convención de `rejections[].file/row` y `source_ref`.
- **§7.2 y §15.2** no definen `missing_count`, `stale_count`, el ganador con varias presentaciones, el orden de `rows[]`, `stores[]` y `winners[]`, ni qué es un sustituto. Con 3 carritos, la **misma** presentación presente en 2 de las 3 canastas sale como `kind: "substitute"` en el motor y en el oráculo (mi fixture «dos de tres»: leche `santa-clara-4x1l-taza` en Bodega y Walmart), y eso contradice §7.2 paso 3 («presentaciones distintas del mismo producto»).
- **§7.4** no fija la convención de `cheapest_store`.
- **§12** no menciona `tickets-README.md` ni `tickets-sinteticos/`.
- El criterio S5-05 dice «este último en RAR».

Rutas: `docs/contratos/motores/S5-interfaz.md` (§7.2, §7.4, §11, §12, §15.2), `docs/criterios/S5.md` (S5-05), `docs/plan/plan.md` §3.4. **Corrección:** el integrador aplica las solicitudes 1 a 6 de `docs/estado/wf1/correccion-S5-r2.md` (contrato 1.5). Agrega la regla para 3 o más canastas: «idéntico = la misma presentación en todas las canastas; lo demás es sustituto o solo en una tienda». Si se quiere otra regla, abre la pregunta en `S5-carritos-detalle`.

**H-r3-02 [media] `compare_carts` con cantidades a granel: el motor y el oráculo dan cifras distintas y ninguna fixture lo prueba.** Mis fixtures:
- `rv-carritos-granel`: dos renglones idénticos de 0.5 a 1.01 y 1.03. Total Bodega: motor **101**, oráculo **102**. Total Walmart: **103** contra **104**.
- Mismo caso, `base_unit_price_cents`: 135, 137, 721 y 736 en el motor; 136, 139, 729 y 743 en el oráculo.
- `rv-carritos-granel-sust`, sustituto con 0.333 y 0.667: 2933 y 2800 en el motor; 2935 y 2801 en el oráculo.

El oráculo (`_agrupa` en `oraculo/pre_calculo.py`) redondea cada importe de línea y calcula el precio por unidad base desde el importe ya redondeado. Eso es redondeo intermedio, y va contra S5-02 («aritmética exacta») y §7.2 paso 3 («precio por unidad base exacto»). El motor (`backend/internal/precios/compare.go`) redondea una sola vez. El §7.2 no dice si el importe de un renglón se redondea por línea, como lo imprime un ticket. Hoy no cambia el canon, porque todas sus cantidades son enteras, pero el ganador de un sustituto podría voltearse. `correccion-S5-r2` lo dejó como observación sin bloqueo ni fixture. **Corrección:** fijar en §7.2 la regla del importe de línea. Si no es técnica, abrir el bloqueo `S5-carritos-granel` y marcarla PROVISIONAL. El precio por unidad base siempre sale de la fracción exacta. Alinear el oráculo y agregar una fixture a granel.

**H-r3-03 [media] `profeco_import`: las reglas de lectura y de liga difieren entre el motor y el oráculo, y ni §11 ni §15.2 las fijan.** Mi CSV `rv-bordes` tiene 17 registros. En producción, el motor da rejected 4, filtered_out 1, skipped 1 y loaded 10; el oráculo da 3, 2, 2 y 9. Las diferencias:
- Precio `0.00`: el motor lo rechaza («precio no positivo»); el oráculo lo **carga** y metería un precio de referencia de $0.00.
- Precio `"1,021.00"`: el motor carga 102100; el oráculo lo rechaza («precio ilegible»).
- `cadena_comercial` vacía: el motor la rechaza; el oráculo la cuenta como filtrada.
- Producto `ATUN DOLORES` contra la liga `ATÚN DOLORES`: el motor liga, porque compara con `Fold` (`profeco/links.go:52`); el oráculo lo salta, porque usa clave exacta (`oraculo/pre_profeco.py:125`).
- Los textos de motivo difieren: «precio con más de dos decimales» contra «precio ilegible»; «columnas: 14 en vez de 15» contra «columnas».
- El motor emite `loaded_proposed`, una clave que §15.2 no lista.

En modo revisión pasa lo mismo. Los cortes reales no se afectan (0 rechazos en 37,548,200). Rutas: `backend/internal/precios/profeco/parse.go`, `links.go`, `oraculo/pre_profeco.py` y `S5-interfaz.md` §11 y §15.2. **Corrección:** escribir en §11 la gramática del precio (positivo, ¿separador de miles?), la lista cerrada de motivos con su texto, el trato de la cadena vacía y si la liga se compara normalizada (con PROVISIONAL y bloqueo si no es técnico). Después alinear el oráculo o el motor y agregar una fixture con estos casos.

**H-r3-04 [media] `difference_vs_cheapest_cents` tiene dos significados según el modo.** En `compare_carts` (`rankCarts`, `compare.go:356`; oráculo `pre_calculo.py:251`) la diferencia se mide contra el menor total aunque esa tienda esté incompleta. En `s5-carritos-presentaciones`, `tercera` tiene `missing_count` 2 y sale con diferencia 0, Walmart con +131.50 y Bodega con +261.00. En `compare_basket` (`rankStores`) se mide contra la primera tienda completa. El §7.4 dice que una tienda con faltantes «no puede ser la más barata», y el §15.2 dice que las claves de `stores[]` son «las mismas» en los dos modos. La opción A de `S5-carritos-detalle` no responde esta pregunta. **Corrección:** agregar la pregunta 4 al bloqueo `S5-carritos-detalle` (o decidirla en §7.2), usar la misma referencia en los dos modos o documentar la diferencia, y fijarla en la fixture.

**H-r3-05 [baja]** El comando de Tesseract de `ocr/doc.go` agrega `--oem 1 -c preserve_interword_spaces=1` al de §12. Está documentado en `docs/estado/wf1/ocr.md`, pero el contrato no lo recoge. **Corrección:** sumarlo a la solicitud de §12.

**H-r3-06 [baja]** `docs/estado/S5-cobertura-profeco.md` §1 y §4.2 dan cifras del corte 2026: Bodega Vallejo 4,461, Sam's 24,788 y Virreyes 21,333. La §6 da las de los tres cortes: 23,809, 94,293 y 77,782. El encabezado lo advierte, pero Axel puede leer la tabla del resumen como total. **Corrección:** poner las cifras de los tres cortes en §1 y §4.2, o marcar cada celda «(corte 2026)».

## Fixtures propias (borradas al terminar)
Armé una raíz sombra en `datos/tmp/rev5r3/root` con copias de `oraculo/pre_*.py`, `_datos/`, `productos.json` y `data.json`, para no tocar `oraculo/fixtures/s5/`. Ahí corrí `../motores s5 dump --fixtures fx` y `python3 -I oraculo/pre_run_fixtures.py --fixtures fx` y comparé hoja por hoja con un script propio. Resultado: `rv-filtro-revision` 0 diferencias, `rv-filtro-cp1252` 0, `rv-carritos-dos-de-tres` 0, `rv-carritos-granel` 8, `rv-carritos-granel-sust` 2, `rv-filtro-bordes` 13 y `rv-filtro-bordes-rev` 13. En las dos de filtro con 0, la única hoja distinta es la clave extra `loaded_proposed` del motor, que no cuento como diferencia.
