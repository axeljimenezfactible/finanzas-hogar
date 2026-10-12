# P2c · revisor-motor-S5-r1 (ronda 1 de 3) · 11-oct-2026, noche

Rol: revisor independiente del trabajo de `motor-S5-puerta2`. No corregí nada. No hice `git add` ni `git commit`, no mandé Telegram y no usé navegador ni WebFetch. No toqué `docs/design/` ni `docs/estado/wfp3/`. Informe: `docs/estado/p2c/REVISION-motor-S5-r1.md`.

## Hecho
- HECHO: el veredicto es **NO APROBADO**, con 0 altos, 3 medios (M-1 `piece_count` a granel, M-2 divergencias y supuestos sin bloqueo, M-3 marca PROVISIONAL de S5-09) y 3 bajos (B-1 independencia según la letra de §15.3, B-2 marcas PROVISIONAL en el código, B-3 artefactos fuera de alcance).
- HECHO: volví a correr todas las verificaciones de (a) a (d) con mis comandos. Además corrí casos propios en una raíz aislada (`datos/tmp/rev-s5-p2c-r1/raiz`, con una copia del oráculo, el binario del motor construido por mí y 5 fixtures propias cuyos esperados escribí a mano con `Fraction`) y una comparación completa, sin recortar, del motor contra el oráculo en las 32 fixtures y en las mías.
- HECHO: lo que pedía la ronda (H-r3-02, H-r3-03, H-r3-04 y M-r5-3) quedó cerrado para las reglas que fija el contrato 1.5. No se decidió ninguna pregunta «sin recomendación».

## Verificación (comandos y salidas resumidas; detalle en la tabla del informe)
```
$ cd backend && go vet ./internal/precios/... ./cmd/motores/          → sin salida
$ gofmt -l internal/precios cmd/motores                               → vacío
$ go test -count=1 -coverprofile=… ./internal/precios/...             → precios 96.0 %, asistida 100 %, ocr 99.3 %, profeco 95.3 %; total 88.8 %
$ go test ./internal/precios/... -run <criterio> -v                   → 40 PASS, 0 FAIL, 0 SKIP (12 patrones)
$ cp -a oraculo/fixtures/s5 …; python3 -I oraculo/pre_genera_fixtures.py; diff -r   → 0 diferencias
$ python3 -I oraculo/pre_run_fixtures.py --out datos/tmp/…/oraculo-s5.json         → 32/32, 1336 hojas, 0 diferencias; cmp con salidas/s5.json: idéntica
$ go -C backend build -o …/motores ./cmd/motores; motores s5 dump …; compara.py s5   → 935/935 ×3, provisionales 315/315, diferencias: 0
$ (fixtures sin esperado) difcompleto.py oráculo motor                → solo `membership.note` (fuera de §15.2)
$ (raíz aislada) 5 fixtures propias: compara.py s5                    → 295/295 ×3, diferencias: 0
$ (raíz aislada) s5-rv-piezas y s5-rv-pendientes sin esperado         → piece_count 3 contra 2; identical_total 2100 contra 2000; references 3 contra 2
$ sha256sum datos/semilla/data.json; git diff HEAD -- datos/semilla/data.json      → 2fe7afae…1d31b; sin cambios
$ sqlite3 -readonly qqp.sqlite "SELECT COUNT(*), SUM(precio LIKE '%,%'), …"        → 37548200|0|0|2026-08-31
$ make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose               → EXIT=0, «Verificación ampliada OK»
```

## Supuestos
- S-r1: tomé como «los 17 registros del revisor» los tipos de caso que describe H-r3-03 (precio 0, «1,021.00», cadena vacía, ATUN/ATÚN). El CSV original ya no existe, así que armé uno propio de 23 registros que los incluye.
- S-r2: clasifiqué las divergencias entre el motor y el oráculo en casos que el contrato no fija como **media**, no como alta. Seguí el criterio de WF-1 (`REVISION-S5-r3.md`, H-r3-02 y H-r3-03); ninguna toca cifras canon ni fixtures.
- S-r3: B-1 (independencia) es baja. El reparto de P2c le dio a un solo agente el motor y el oráculo, y el oráculo sigue sin leer Go. La parte copiada del motor la cuento en M-2(d).

## Bloqueos
No abrí bloqueos: `docs/estado/bloqueos/` no está en mi alcance. Los que propongo abrir o ampliar están en «Solicitudes fuera de alcance».

## Pendientes
- Ronda 2: verificar M-1, M-2 y M-3, y los bajos, con las mismas fixtures propias (están descritas en el informe para rehacerlas).
- Solicitudes fuera de alcance, para el integrador:
  - Contratos: §7.2 paso 6 (renglón de S-11), §7.2 pasos 1 y 2 (precio de la parte emparejada), §6 (referencias), §11 (fechas, liga a presentación desconocida, tope de `rejections[]`, S-1 a S-3) y la tabla de §15 con las 3 fixtures nuevas.
  - Bloqueos: actualizar `S5-carritos-granel.md` («Cómo seguí») y, si no son técnicas, las preguntas nuevas de M-1 y M-2.
  - Criterios: alinear S5-09 con §8 si se quiere quitar la marca PROVISIONAL.
  - C1: la clave `cheapest_store_id_`, que ya registra `REVISION-C1-r1.md`.

## Propuestas para Axel
- Ninguna nueva. Las preguntas de fondo de M-2(a) (qué precio cuenta cuando compraste el mismo producto a dos precios en una compra) y M-1 (cómo se cuentan las piezas a granel) solo llegarán a Axel si contratos no las puede fijar como regla técnica.

## Consumo aproximado
Unas 47 acciones de herramienta · unos 180 KB leídos · unos 45 KB escritos (los 2 informes, unos 22 KB, más scripts y fixtures temporales en `datos/tmp/rev-s5-p2c-r1/`, ya borrados).
