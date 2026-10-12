# P2c · correccion-motor-S5-r1 · 11-oct-2026, noche

Corregí los tres hallazgos medios (M-1, M-2 y M-3) de `docs/estado/p2c/REVISION-motor-S5-r1.md`. En esta ronda no había ninguno alto.
- Lo hice dentro de mi alcance: `backend/internal/precios/` (con `profeco/`), `backend/cmd/motores/s5.go`, `oraculo/pre_*.py` y las fixtures s5, más este reporte y un bloqueo nuevo, como permite la regla suprema.
- No toqué las observaciones de baja severidad (B-1 a B-3), ni `docs/design/` ni `docs/estado/wfp3/`.
- No hice `git add` ni `git commit`.
- No decidí ninguna pregunta «sin recomendación».

## Hecho

**HECHO 1 · M-3 (S5-09, alta): volvió la marca PROVISIONAL a la nota de membresía.**
- `backend/internal/precios/extras.go:90`: `Note` vuelve a decir `"sin costo capturado (PROVISIONAL, S5-09)"`, que es el texto de `HEAD`. El comentario explica que la marca se queda hasta que el integrador alinee el criterio S5-09 con §8 del contrato 1.5.
- `backend/internal/precios/precios_test.go:612`: `TestMembresiaSinCosto` ahora comprueba el texto exacto de `Note`, además de `verdict = missing_fee`.
- El dump del motor emite esa nota en `s5-membresia-equilibrio`. La clave `note` no está en §15.2, así que no se compara con el oráculo.

**HECHO 2 · M-1 y M-2: bloqueo nuevo `docs/estado/bloqueos/S5-motor-oraculo-reglas.md`, con 8 preguntas.**
- Cada pregunta trae qué falta, las opciones, las cifras medidas y el impacto. Todas quedan PROVISIONAL, en uso: el motor. No elegí qué lado tiene razón.
- Preguntas:
  - 1 (M-1): qué es «renglón» en S-11 (`piece_count` a granel).
  - 2 (M-2a): qué precio toma la parte emparejada cuando la misma presentación viene a precios distintos.
  - 3 (M-2b): si `references[]` lleva una referencia por fuente o todas.
  - 4 (M-2c): formatos de fecha.
  - 5 (M-2c): liga a una presentación desconocida.
  - 6 (M-2c): tope de `rejections[]`.
  - 7 (M-2d): gramática del precio, S-1 a S-3.
  - 8 (M-2d): liga repetida.
- Las preguntas 1 a 3 son de negocio o de pantalla: las fija contratos si las juzga técnicas y, si no, Axel. Las preguntas 4 a 8 son técnicas, de §11, y van a contratos.

**HECHO 3 · Marcas PROVISIONAL en el código de los dos lados**, cada una con «bloqueo S5-motor-oraculo-reglas pregunta N»:
- `compare.go`: `agg`/`take` (pregunta 2) y `pieces` (pregunta 1).
- `pricing.go`: `References` (pregunta 3).
- `profeco/parse.go`: `ParsePriceCents` (pregunta 7) y `dateLayouts` (pregunta 4).
- `profeco/import.go`: `MaxRejections` (pregunta 6).
- `profeco/links.go`: liga repetida (pregunta 8).
- `cmd/motores/s5.go`: presentación desconocida (pregunta 5).
- `oraculo/pre_calculo.py`: `_piezas`, `_agrupa` y `current_price`.
- `oraculo/pre_profeco.py`: `_precio`, `_FECHAS`, liga repetida, `presentacion_id` desconocido y el tope de rechazos.
- En el oráculo, el comentario «supuesto del agente motor-S5-puerta2» ahora dice que S-1 a S-3 y la liga repetida son una **copia del motor en uso** y que no son evidencia independiente (§15.3).
- Solo cambiaron comentarios: ninguna salida cambia.

**HECHO 4 · Medí las divergencias con mis propios casos.** Fixture borrador sin `esperado` en `datos/tmp/corr-s5-r1/fx/` (fuera de git):
- `piezas-1` (Bodega 0.5 de arroz y 0.5 + 0.25 de atún): `piece_count` de Bodega es 2 en el motor y 3 en el oráculo. Es la única diferencia en 60 hojas.
- `piezas-2` (Walmart 0.5 de atún y 1 bolsa de arroz, sobra 0.5): Walmart es 2 en el motor y 3 en el oráculo. Es la única diferencia en 60 hojas.
- `precio-emparejado` (Bodega 1 atún a $20.00 y otro a $22.00; Walmart 1 a $21.00): 8 diferencias.
  - `identical_total_cents` de Bodega: 2000 contra 2100.
  - Ganador: Bodega contra ninguno.
  - `ties`: 0 contra 1.
  - Línea excedente: 2200 contra 2100.
  - Precio por unidad base: 14286 contra 15000.
  - El total de Bodega es 4200 en los dos.
- `referencias` (2 PROFECO y 1 `site_assisted`): el motor da 2 y el oráculo 3. Sobra la de PROFECO del 01-oct a 1900.
- Fechas, con una copia literal de `dateLayouts` corrida con `go run` y `_fecha` del oráculo, sobre 14 textos:
  - Solo el motor acepta `2026/8/3`, `2026-8-3`, `3/8/2026`, `20260831` y `…T10:20:30Z`.
  - Solo el oráculo acepta la hora sin segundos y `31/08/2026 10:20:30`.
  - Los dos aceptan la fracción de segundo. El reporte del motor decía que solo el oráculo la aceptaba; no es así.

**HECHO 5 · Hallazgo nuevo, fuera de los que me tocan (no lo corregí).** El `created_at` de un precio con zona horaria distinta de `Z` da resultados distintos en los dos lados:
- Con `2026-10-10T12:00:01-06:00`, el oráculo responde toda la fixture con `invalid_date`.
- El motor (`time.RFC3339`, `precios/eval.go:259`) lo acepta.
- La causa está en `oraculo/pre_comun.py:51` (`instante`): solo admite el sufijo `Z`. §15.1 dice «instante RFC 3339», que admite zona horaria.
- `pre_comun.py` lo comparten los oráculos S2 a S5. Va como solicitud.

## Verificación
Todo desde `/home/deepsec/finanzas-hogar`:
```
$ (cd backend && gofmt -l ./internal/precios ./cmd/motores && go vet ./internal/precios/... ./cmd/motores && echo vet-ok)
vet-ok
$ (cd backend && go test ./internal/precios/... -run Membresia -count=1 -v)
--- PASS: TestMembresiaEquilibrio (0.00s)
--- PASS: TestMembresiaSinCosto (0.00s)
ok  finanzas-hogar/internal/precios
$ (cd backend && go test ./internal/precios/... -count=1 -cover)
ok  finanzas-hogar/internal/precios          coverage: 96.0% of statements
ok  finanzas-hogar/internal/precios/asistida coverage: 100.0% of statements
ok  finanzas-hogar/internal/precios/ocr      coverage: 99.3% of statements
ok  finanzas-hogar/internal/precios/profeco  coverage: 95.3% of statements
$ cp -r oraculo/fixtures/s5 datos/tmp/corr-s5-r1/s5-antes && python3 -I oraculo/pre_genera_fixtures.py && diff -r datos/tmp/corr-s5-r1/s5-antes oraculo/fixtures/s5
32 fixtures en /home/deepsec/finanzas-hogar/oraculo/fixtures/s5
(diff -r sin salida: 0 diferencias)
$ python3 -I oraculo/pre_run_fixtures.py --out datos/tmp/corr-s5-r1/oraculo-s5.json && cmp datos/tmp/corr-s5-r1/oraculo-s5.json oraculo/salidas/s5.json
32 fixtures, 32 coinciden, 1336 hojas de esperado, 0 diferencias
(cmp sin salida: salidas/s5.json idéntica)
$ go -C backend build -o ../datos/tmp/motores ./cmd/motores && datos/tmp/motores s5 dump --fixtures oraculo/fixtures/s5 --out datos/tmp/motor-s5.json && python3 -I oraculo/compara.py s5 --motor datos/tmp/motor-s5.json
oráculo contra esperado: 935 de 935 hojas coinciden
motor contra esperado: 935 de 935 hojas coinciden
oráculo contra motor: 935 de 935 hojas coinciden
provisionales: 315 de 315 hojas coinciden
diferencias: 0
$ grep -o '"note": *"[^"]*"' datos/tmp/motor-s5.json | sort | uniq -c
      1 "note": "sin costo capturado (PROVISIONAL, S5-09)"
$ python3 -I oraculo/pre_verify.py | tail -1
TODO COINCIDE
$ sha256sum oraculo/model.json oraculo/verify3.json datos/semilla/data.json
1a57066689ddeed31e0001259f341c7a3f782a46cd4a05128ff453aad2f1bfd4  oraculo/model.json     (= MODEL_SHA)
d86087f3797e08ceabd2c721bf423115345b318ae7fe5fa13b3400cd87b95041  oraculo/verify3.json   (= VERIFY3_SHA)
2fe7afae5b1be77fafeed787d4c7bd4d68fe469b78a86898b370c1ba1031d31b  datos/semilla/data.json (git diff HEAD: sin cambios)
$ make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose   # salida en datos/tmp/corr-s5-r1/verify-all.log
EXIT=0
   K25 OK  finanzas-hogar/internal/precios 96.0 % · asistida 100.0 % · ocr 99.3 % · profeco 95.3 %
   s5 · 32 fixtures, 32 coinciden, 1336 hojas de esperado, 0 diferen · salidas/s5.json idéntica
   oráculo contra motor: 935 de 935 hojas coinciden · provisionales: 315 de 315 hojas coinciden
== WF-1: verificación completa OK (Go y cobertura K25, oráculos S2 a S5, recetario, tokens y contraste, U4 y paquete)
== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)
```
Canon intacto: ninguna fixture ni salida cambió, y `s5-comparativa-canon` sigue en 935 de 935 hojas.

## Supuestos
- S-a. Para M-1 y M-2, «corregir» quiere decir documentar: un bloqueo, solicitudes y marcas PROVISIONAL «en uso: el motor». No alineé ningún lado ni agregué fixtures porque la revisión pide hacerlo «con la regla fijada» y dice: «No se elige aquí qué lado tiene razón».
- S-b. Abrí un bloqueo nuevo en vez de sumar preguntas a `S5-carritos-granel`, `S5-carritos-detalle` o `S5-profeco-filtro`. Esos archivos no son de mi alcance, y la regla suprema permite abrir uno nuevo.
- S-c. El texto de la nota es el que tenía `HEAD`: «sin costo capturado (PROVISIONAL, S5-09)».

## Bloqueos
- **Nuevo:** `docs/estado/bloqueos/S5-motor-oraculo-reglas.md` (ABIERTO, 8 preguntas). No detiene nada: las 8 reglas quedan PROVISIONAL con lo que hace el motor.
- **Siguen abiertos** (no los toqué): `S5-carritos-granel`, `S5-carritos-detalle` (Q2-50 a Q2-52), `S5-profeco-filtro` (Q2-56 a Q2-58) y `S5-carrito-comparativa` pregunta 2 (Q2-48).

## Pendientes
- Cuando contratos (o Axel) fije cada una de las 8 reglas: alinear el lado que no la cumpla y agregar una fixture por caso.
  - Las cifras base están en el bloqueo.
  - La pregunta 1 va como consulta nueva en `s5-carritos-granel`, con un renglón partido y un excedente a granel (`piezas-1` y `piezas-2`).
- B-1, B-2 y B-3 (bajos) no estaban en mi encargo.
- El `created_at` con zona horaria (HECHO 5) sigue sin corregir.

## Propuestas para Axel
- Si contratos te pasa las preguntas 1 a 3 de `S5-motor-oraculo-reglas`:
  - Cómo contar las piezas de un renglón a granel.
  - Qué precio toma lo emparejado cuando el mismo artículo viene a dos precios.
  - Si se muestran todas las referencias de PROFECO o solo la última.
- No tienen recomendación y ninguna mueve dinero del canon.

## Solicitudes fuera de alcance
1. **Contratos** (`S5-interfaz.md`): fijar las 8 preguntas de `S5-motor-oraculo-reglas` o pasarlas a Axel.
   - §7.2 paso 6: qué es «renglón» en S-11.
   - §7.2 pasos 1 y 2: precio de la parte emparejada.
   - §6: `references[]`.
   - §11: fechas, presentación desconocida, tope de `rejections[]`, S-1 a S-3 y liga repetida.
2. **Integrador:** alinear el criterio S5-09 («sin costo capturado» PROVISIONAL) con §8 del contrato 1.5, que no lleva la marca. Mientras tanto, el motor la conserva.
3. **Integrador:**
   - Agregar `S5-motor-oraculo-reglas.md` a `docs/estado/bloqueos/INDICE.md`, en «Bloqueos nuevos de P2c» y en el conteo.
   - Decidir si las preguntas 1 y 2 se suman a `S5-carritos-granel` y a `S5-carritos-detalle`.
4. **Siguiente revisión o integrador:** `oraculo/pre_comun.py:51` (`instante`) rechaza un `created_at` RFC 3339 con zona horaria que no sea `Z`, y el motor lo acepta (HECHO 5). Es compartido por S2 a S5.

## Consumo aproximado
Unas 33 acciones de herramienta · unos 140 KB leídos · unos 30 KB escritos: el bloqueo (unos 8 KB), este reporte (unos 9 KB), los comentarios y la prueba (unos 5 KB) y lo temporal en `datos/tmp/corr-s5-r1/` (unos 6 KB).
