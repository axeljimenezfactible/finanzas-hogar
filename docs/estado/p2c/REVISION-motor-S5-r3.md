# P2c · REVISION-motor-S5-r3 · 11-oct-2026, noche

Revisor independiente `revisor-motor-S5-r3`, ronda 3 de 3 del flujo P2c «Aplicación de la puerta 2».
- No escribí nada de lo que reviso y no corregí nada.
- Evalué solo contra `docs/criterios/S5.md` 1.2 (completo) y `docs/contratos/motores/S5-interfaz.md` 1.5, con las respuestas de Axel (`docs/estado/puerta2/RESPUESTAS.md` y las recomendadas de `PUERTA-2.md`).
- Objeto: el estado actual de `backend/internal/precios/` (con `profeco/`), `backend/cmd/motores/s5.go`, `oraculo/pre_*.py`, las fixtures s5 nuevas, el bloqueo nuevo `docs/estado/bloqueos/S5-precio-vigente-captura.md` y los reportes `motor-S5-puerta2.md` y `correccion-motor-S5-r2.md`.
- Reejecuté todo con mis propios comandos. Mis sondas las armé en una raíz aislada, `datos/tmp/rev-s5-r3/root/` (copia de `oraculo/*.py`, de las fixtures, de `data.json` y de `productos.json`), para no tocar `oraculo/fixtures/s5/`. Ya la borré.

## Veredicto: NO APROBADO (0 altos, 1 medio, 6 bajos)

**HECHO:**
- M-1 de la ronda 2 quedó cerrado en lo que pedía: con `-06:00` el motor y el oráculo dan 2100 `ticket_manual`, y con `+00:00` dan 2000 `manual`. La fixture `s5-precio-vigente-zona` hace fallar al oráculo viejo, y `s5-ingesta-zona-hora` y `s5-ingesta-zona-minuto` hacen fallar al motor viejo (los dos los reconstruí desde `HEAD`).
- Cobertura, gofmt, vet, regeneración de fixtures, `run_fixtures`, `compara.py s5`, granel, `profeco_import`, `difference_vs_cheapest_cents`, canon, `data.json` y `make verify-all` están bien.

**Lo que impide aprobar (M-1 de esta ronda):** quedan cinco textos de `created_at` en los que un lado acepta y el otro responde `invalid_date`.
- El bloqueo nuevo dice lo contrario: «los dos lados … rechazan con `invalid_date` lo que RFC 3339 §5.6 no admite» y, en su pregunta 3, «aquí no hay divergencia».
- Con uno de esos textos, la zona decide S-8 en el motor (2100) mientras el oráculo rechaza toda la fixture. Es el mismo patrón que M-1 de la ronda 2.

## Verificaciones (evidencia propia)

| # | Verificación | Resultado | Evidencia (mis comandos, desde `/home/deepsec/finanzas-hogar`) |
|---|---|---|---|
| a1 | gofmt y vet (S5-13) | CUMPLE | `gofmt -l internal/precios cmd/motores`: vacío. `go vet ./internal/precios/... ./cmd/motores`: limpio |
| a2 | Cobertura K25 ≥ 85 % (S5-13) | CUMPLE | `go test -count=1 -coverprofile=datos/tmp/rev-s5-r3/s5.cov ./internal/precios/...`: precios 96.0 %, asistida 100.0 %, ocr 99.3 %, profeco 95.3 %; `go tool cover -func … \| tail -1`: total 88.8 % (con los dos `main` sin pruebas) |
| a3 | Sin `float` en dinero ni red (S5-07, S5-08, S5-13) | CUMPLE | `float64` solo en `ocr/tsv.go:13, 64 y 68` (confianza de Tesseract). `net/http`, `http.Get` o `http.Client` fuera de pruebas: 0. URLs de las tiendas en Go fuera de pruebas: 0. `grep -c Veredicto docs/estado/S5-terminos.md` = 3 |
| a4 | Pruebas nombradas y nuevas | CUMPLE | `-run 'Comparativa\|Membresia\|Unidad\|MejorTienda\|Frescura\|Fuente\|LecturaAsistida\|Mercado\|PreciosViejos'`: `ok`. `-run 'TestEvalErrores\|TestEvalCreatedAtZona' -v`: los dos PASS. La línea de Python de S5-04 imprime `[0.9836, 0.4918, 0.0]` |
| b1 | Las fixtures se regeneran idénticas | CUMPLE | `pre_genera_fixtures.py` en mi raíz aislada: «36 fixtures»; `diff -r` contra mi copia previa: vacío |
| b2 | `run_fixtures` reproduce `salidas/s5.json` | CUMPLE | `pre_run_fixtures.py --out datos/tmp/rev-s5-r3/oraculo-s5.json`: «36 fixtures, 36 coinciden, 1359 hojas de esperado, 0 diferencias»; `cmp` contra `oraculo/salidas/s5.json`: idéntico |
| b3 | `compara.py s5` | CUMPLE | Con mi binario (`go build -o datos/tmp/rev-s5-r3/motores ./cmd/motores`, `s5 dump`): oráculo contra esperado 958/958, motor contra esperado 958/958, oráculo contra motor 958/958, provisionales 315/315, **diferencias: 0** |
| b4 | Las fixtures nuevas detectan el error | CUMPLE | Oráculo con `pre_comun.py` de `HEAD`: falla `s5-precio-vigente-zona` y `s5-ingesta-hora-24`. Motor con `eval.go` de `HEAD`: en `s5-ingesta-zona-hora` y `s5-ingesta-zona-minuto` no da error (sale «__falta__» donde se espera `invalid_date`) |
| c1 | Granel 0.5, 1.01 y 1.03 | CUMPLE | `rv3-granel`, mi cálculo escrito antes de correr: Bodega 0.5 × 101 = 50.5, 1.01 × 2201 = 2223.01 y 1.03 × 1999 = 2058.97, total 4332.48 → **4332**; Walmart 0.5 × 103 = 51.5, 1.01 × 2199 = 2220.99 y 1.03 × 2003 = 2063.09, total 4335.58 → **4336**; diferencia +4; líneas 51/52, 2223/2221 y 2059/2063; ganadores Bodega 2 y Walmart 1. El motor y el oráculo dan exactamente eso, y su salida completa es igual |
| c2 | `profeco_import` con 17 registros propios | CUMPLE | Precio 0, «1,021.00», precio vacío, cadena vacía, ATUN, ATÚN y «atún dolores / 140 g / Dolores», liga propuesta, sin liga, municipio fuera, duplicado, «ayer», «12.345», 14 columnas, «-5.00» y Walmart Tepeyac. Producción: read 17, rejected 7, filtered_out 1, skipped 2, loaded 6, loaded_proposed 0 y duplicates 1. Revisión: skipped 1 y loaded_proposed 1. La identidad de §11 da 17 en los dos modos. Los 7 motivos salen de la lista de §11 paso 3 y 4 («-5.00» → «precio ilegible»). Motor y oráculo iguales en la salida completa |
| c3 | `difference_vs_cheapest_cents` con un solo significado | CUMPLE | `rv3-dif`, con `compare_basket` y con `compare_carts`: t1 3000 → 0, t2 3200 → +200, t3 (incompleta, menor total 500) al final con `null`; `cheapest_store` t1. `rv3-dif-ninguna`: las dos `null` y `cheapest_store` `null`. `compare.go:35`: `*int64`, «null en incompletas (§7.4)» |
| c4 | Canon intacto (S5-01) | CUMPLE | `s5-comparativa-canon` en mi dump: 309450 y 323300 (**+13850**), idénticos 277350 y 287400 (**+10050**), piezas 61/60, `identical_pieces` 55, ganadores 22/3 y `ties` 8. `make verify-all`: «pre_verify.py: 15 comprobaciones OK · TODO COINCIDE» |
| c5 | M-1 r2: `created_at` con zona | CUMPLE | 7 sondas válidas, iguales en los dos lados y en mi cálculo: `07:00:00-06:00` → 2100, `12:00:00+00:00` → 2000, `12:00:00-00:00` → 2000, `12:30:00.5Z` → 2100, `08:29:59.999-04:00` → 2000, `2026-10-11T00:15:00+11:44` → 2100 y `2026-10-11T12:29:00+23:58` → 2100 |
| c6 | `created_at` fuera de RFC 3339 (§15.1, §15.3) | **NO CUMPLE** (M-1) | 15 sondas. Iguales (`invalid_date` en los dos): `+24:00`, `+05:60`, `24:00:00Z`, `23:59:60Z`, mes 13, 30 de febrero, `-6:00`, espacio y minúsculas. Distintas: las 6 de M-1 abajo |
| d1 | `data.json` intacto | CUMPLE | `git diff --quiet HEAD -- datos/semilla/data.json`: sin cambios; sha256 `2fe7afae…1031d31b` |
| f1 | Verificación global | CUMPLE | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: EXIT=0; «s5 · 36 fixtures, 36 coinciden … salidas/s5.json idéntica»; «WF-1: verificación completa OK»; «Verificación ampliada OK». Después, los 43 archivos de `oraculo/fixtures/s5/` (36 fixtures y 7 de `_datos/`), `salidas/s5.json` y `data.json` tenían el mismo sha256 que antes |
| — | S5-06 (OCR) | sin medir | `ls datos/semilla/tickets/ \| wc -l` = 0. No es objeto de esta ronda |

## Hallazgos

**M-1 [media] `created_at`: quedan textos que no son RFC 3339 (o que sí lo son) en los que el motor y el oráculo no coinciden, y el bloqueo afirma que no hay divergencia.**
- Rutas:
  - `backend/internal/precios/eval.go:261`: `time.Parse(time.RFC3339, …)` cae a su lector general cuando falla el rápido, y ese lector acepta la hora de un dígito y la coma decimal. `rfc3339Offset` (línea 575) solo revisa la zona.
  - `oraculo/pre_comun.py:42`: `_RFC` usa `\d`, que en Python acepta dígitos Unicode, y `$`, que deja pasar un salto de línea final.
  - `oraculo/pre_comun.py:64`: `dt.date(a, me, d)` rechaza el año 0000.
  - `docs/estado/bloqueos/S5-precio-vigente-captura.md`, «Contexto» (punto 3) y pregunta 3.
  - Contrato: `S5-interfaz.md` §15.1 (`created_at` «instante RFC 3339»), §15.3 (motor = oráculo) y §6 (S-8).
- Evidencia: mis fixtures `rv3-x-*`. En cada una compiten `manual` a 2000 con `created_at` 12:30:00Z y `ticket_manual` a 2100 con el texto de la tabla.

| Texto del `ticket_manual` | ¿RFC 3339? | Motor | Oráculo |
|---|---|---|---|
| `2026-10-10T7:00:00Z` | no (hora de 1 dígito) | 2000 `manual` (lo lee como 07:00Z) | `invalid_date` |
| `2026-10-10T7:00:00-06:00` | no | **2100 `ticket_manual`** (13:00Z, decide S-8) | `invalid_date` |
| `2026-10-10T12:30:00,5Z` | no (`time-secfrac` es «.» y dígitos) | 2100 `ticket_manual` | `invalid_date` |
| `2026-10-10T12:30:00Z` y un salto de línea | no | `invalid_date` | 2100 `ticket_manual` |
| `2026-10-10T13:00:00Z` con dígitos arábigo-índicos (U+0660 a U+0669) | no (`DIGIT` es %x30-39) | `invalid_date` | 2100 `ticket_manual` |
| `0000-01-01T00:00:00Z` | sí (`date-fullyear` = 4DIGIT) | 2000 `manual` | `invalid_date` |

- Por qué es medio y no bajo:
  - Es la misma clase y el mismo efecto que M-1 de la ronda 2: un lado da un precio y el otro rechaza toda la fixture.
  - Además, el bloqueo, que es lo que lee el siguiente implementador, dice que no hay divergencia.
  - No mueve ninguna cifra canon, y ninguna fixture usa estos textos, así que `compara.py s5` da 0 diferencias.
- No es una regla sin decidir, salvo el año 0000:
  - §15.1 dice RFC 3339, y los cinco primeros textos no lo son. La corrección r2 usó ese mismo argumento para `+24:00` y `+05:60` (su HECHO 4).
  - El año 0000 sí es RFC 3339 válido. Que se acepte o que el perfil lo excluya es materia de la pregunta 3 del bloqueo.
- **Corrección sugerida:**
  - Motor: antes de `time.Parse`, exigir la gramática de §5.6 en ASCII, por ejemplo `^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(\.[0-9]+)?(Z|[+-][0-9]{2}:[0-9]{2})$`. `rfc3339Offset` se queda para los rangos.
  - Oráculo: `[0-9]` (o `re.ASCII`) en lugar de `\d`, y `fullmatch` (o `\Z`) en lugar de `$`.
  - Año 0000: no decidirlo. Sumarlo a la pregunta 3 del bloqueo con estas cifras y marcarlo PROVISIONAL.
  - Agregar fixtures `s5-ingesta-*` con la hora de un dígito, la coma decimal, el salto de línea y los dígitos no ASCII, todas con `invalid_date`.
  - Corregir en el bloqueo la frase del punto 3 de «Contexto» y la de la pregunta 3.

**B-1 [baja, sigue de r1] Independencia según la letra de §15.3.** La corrección r2 volvió a tocar los dos lados. Declara que escribió el oráculo primero, desde el contrato y RFC 3339, y que leyó `eval.go` después. Corrección: el integrador registra la excepción del reparto de P2c (solicitud 2 de `correccion-motor-S5-r2.md`).

**B-2 [baja, sigue de r1] Faltan marcas PROVISIONAL.** `grep S5-carritos-granel backend/internal/precios/compare.go` y `grep S5-profeco-filtro profeco/links.go profeco/filter.go` siguen sin resultados.

**B-3 [baja, sigue de r1, crece] Tabla de §15.** `S5-interfaz.md` §15 no lista `s5-carritos-granel`, `s5-profeco-bordes`, `s5-profeco-motivos`, `s5-precio-vigente-zona`, `s5-ingesta-zona-hora`, `s5-ingesta-zona-minuto` ni `s5-ingesta-hora-24` (`grep -c`: 0). Es de contratos (solicitud 1 de la corrección r2).

**B-4 [baja, sigue de r2] Índice de bloqueos.** `docs/estado/bloqueos/INDICE.md` no lista `S5-motor-oraculo-reglas` ni `S5-precio-vigente-captura` (`grep -c`: 0). Es del integrador.

**B-5 [baja, sigue de r2] Pregunta 5 de `S5-motor-oraculo-reglas` sin la cifra medida.** La línea 77 sigue diciendo «No medí qué pasa con un id que no está en el catálogo». La cifra está en `REVISION-motor-S5-r2.md` B-5.

**B-6 [baja, nueva] Versión del contrato en el código.** `backend/internal/precios/types.go:4` (el comentario del paquete que muestra `go doc`) dice «S5-interfaz.md 1.4», pero el contrato vigente es 1.5. Corrección: actualizar el número.

## Estado de los hallazgos de la ronda 2

| Hallazgo r2 | Estado | Evidencia |
|---|---|---|
| M-1 [media] `created_at` con zona (§15.1, §15.3, §6 S-8) | **CERRADO** en lo que pedía. Su residuo es M-1 de esta ronda | `pre_comun.py:54-77` acepta `Z` y `±hh:mm`, con o sin fracción, y compara en UTC. `eval.go:261` y `:575` rechazan `+24:00` y `+05:60`. Fixture `s5-precio-vigente-zona` con las cifras de r2 (2100 con `-06:00` y 2000 con `+00:00`) y `provisional` `null`. Mis 7 sondas válidas coinciden. Las fixtures detectan el oráculo viejo y el motor viejo (b4). La premisa falsa de r1 quedó corregida en el HECHO 1 de la corrección r2 |
| B-1 a B-5 [baja] | SIGUEN ABIERTOS | ver arriba; no estaban en el encargo de la corrección r2 |

## Separación
- **HECHO:** todo lo de la tabla de verificaciones y la tabla de M-1, medido por mí.
- **PREGUNTA:** ninguna nueva para Axel. El año 0000 va a la pregunta 3 del bloqueo, que es técnica, para contratos. Las preguntas 1, 2 y 4 del bloqueo y las de `S5-motor-oraculo-reglas` siguen sin respuesta, y no las decido.
- **PROPUESTA:** corregir M-1 con las dos gramáticas estrictas en ASCII y cuatro fixtures `s5-ingesta-*`. Es poco código en cada lado y no cambia ninguna salida actual.
