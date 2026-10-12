# revisor-motor-S2-r2 · P2c · ronda 2 de 3 de la revisión del motor S2

Agente: `revisor-motor-S2-r2` (revisor independiente) · 11-oct-2026, noche · Alcance: revisar, no corregir. Evalué contra `docs/criterios/S2.md` 1.2 completo y `docs/contratos/motores/S2-interfaz.md` 1.3. Informe: `docs/estado/p2c/REVISION-motor-S2-r2.md`.

## Hecho
- Veredicto: **APROBADO** (0 altos y 0 medios).
- M1 de r1 (media; `cover_test.go:239` pegaba la cifra canon 97375) quedó **CERRADO**. Lo comprobé con una mutación propia en una copia del paquete: la aserción falla con `esperado=125869` cuando se le pide `periodos[2]`. Mi script `canon_lits.py` busca las hojas de `model.json` entre los literales de las pruebas y ya no encuentra ningún esperado canon pegado.
- Reejecuté todas las verificaciones (a) a (e) y `make verify-all` con mis comandos. Cobertura de `motor_fin`: 89.2 %. 32 fixtures, 3185 hojas, 159 de 159 de `model.json` y 0 diferencias.
- La regla 4 de §2.12 (PROVISIONAL) la probé en los dos lados con 2 fixtures mías y cifras calculadas a mano (gasto de 90001 a 3 meses: 30000, 30000 y 30001):
  - 5 huecos que el motor y el oráculo rechazan con `invalid_argument` **por la regla 4** (mensaje leído en los dos lados);
  - 2 rechazos por la regla 1;
  - 4 listas contiguas aceptadas (del 15 al 14, ventanas cortas del caso 2, cambio de año y 29 de febrero de 2028);
  - meses con hueco aceptados;
  - además, A + C de la interfaz.
- Estado de los anteriores: siguen abiertos B5, B6 y B7 de r1 y B2 y parte de B3 de r3, todos bajos y sin asignar. B4 de r3 está cerrado. No encontré hallazgos nuevos.
- Lo único que cambió después de r1 es `cover_test.go` (según `stat`), y coincide con la corrección.

## Verificación (comandos y salidas)
Están en la tabla del informe, en las filas a1 a g. Las principales:
- `cd backend && go test ./internal/motor_fin/ -count=1 -cover` → `ok … coverage: 89.2% of statements`; 20 pruebas de primer nivel y 294 subpruebas PASS, 0 FAIL y 0 SKIP.
- `gofmt -l .` vacío; `go vet ./...` limpio.
- `python3 -I oraculo/fin_run_fixtures.py --fixtures oraculo/fixtures/s2 --out datos/tmp/rev-s2-r2/oraculo-s2.json` → «32 fixtures, 32 coinciden, 0 diferencias»; `cmp` contra `salidas/s2.json` idéntico.
- `motores s2 dump` desde la raíz y `compara.py s2` → 3185 de 3185 en las tres comparaciones, 159 de 159 y «diferencias: 0».
- Fixtures mías: `compara.py s2 --fixtures datos/tmp/rev-s2-r2/fx` → 93 de 94 hojas. La que falta es `hora-24` (B5): el motor da `invalid_date` y el oráculo lo acepta.
- `sha256sum -c` de `model.json`, `verify3.json` y `data.json` → OK; `git diff --quiet` sin cambios.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → exit 0, «Verificación ampliada OK».

Borré `datos/tmp/rev-s2-r2/` al terminar. No toqué carpetas temporales de otros agentes, ni archivos de código, del oráculo o de las fixtures.

## Supuestos
1. Mantengo como bajos los abiertos de r1 (B5, B6 y B7). Los artefactos que tendrían que cambiar para B6 y B7 (la interfaz y el bloqueo) no son del implementador del motor. Además, ninguno cambia una cifra ni bloquea a S6 sin aviso: el compilador y `go doc` muestran `Instant`. Es la misma severidad que les dio r1.
2. `292125` no es una hoja de `model.json` (lo comprobé con el script). Sale de la interfaz §2.12 y §3.2 y de REVISION-S2-r3 M1, así que tenerlo pegado en la prueba no es «No aceptable».

## Bloqueos
Ninguno nuevo. Siguen ABIERTOS, sin decidir y con su marca en los dos lados: `S2-ventanas` pregunta 4 (casos 1 y 2), `S2-metas-vencida`, `S2-ahorro-quincenal` y `S2-flujo-quincenal` pregunta 1 (Q2-12).

## Pendientes
- B5: hora 24 en `oraculo/fin_comun.py` `instante`. Le toca al agente del oráculo S2.
- B6: §1.1 y §2.1 de la interfaz con `Instant`. Le toca al dueño del contrato.
- B7: la línea de HECHO de `S2-metas-vencida.md` y de `bloqueos/INDICE.md:405`. Le toca al dueño de los bloqueos.
- B2 de r3: PREGUNTA para §3.15.
- B3 de r3: la tabla de §5 de la interfaz (27 fixtures contra 32 reales).

## Propuestas para Axel
- Ninguna nueva. Siguen en pie las del reporte del implementador: contestar `S2-ventanas` pregunta 4 y `S2-metas-vencida`.

## Consumo aproximado
- Unas 31 acciones de herramienta.
- Unos 150 KB leídos: criterios S2 (34 KB), fragmentos de la interfaz, informes r1 y de la corrección, reporte del implementador, diffs de Go y Python y bitácoras.
- Unos 22 KB escritos: los 2 informes y los scripts y fixtures temporales, ya borrados.
