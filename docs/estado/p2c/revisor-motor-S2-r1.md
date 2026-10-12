# P2c · revisor-motor-S2-r1 · 11-oct-2026, noche

Informe completo: `docs/estado/p2c/REVISION-motor-S2-r1.md`. **Veredicto: NO APROBADO.** 0 altos, 1 medio (M1) y 5 bajos: B2 y B3 (que se arrastran de r3, B3 en parte) y los nuevos B5, B6 y B7.

## Hecho
- Revisé el diff de `backend/internal/motor_fin/` (más `go doc` de `Instant` y `Params`), de `oraculo/fin_devengo.py`, `fin_metas.py`, `fin_flujo.py`, `fin_proyeccion.py`, `fin_presupuesto.py` y `fin_genera_fixtures.py`, las fixtures nuevas `s2-ventanas-a-mas-c` y `s2-metas-vencida`, las 5 cambiadas y el reporte `motor-S2-puerta2.md`. Los evalué contra `S2.md` 1.2 completo y `S2-interfaz.md` 1.3, con `RESPUESTAS.md` y `PUERTA-2.md` (Q2-01 y Q2-12 a Q2-27) como mandato.
- Reejecuté por mi cuenta todas las verificaciones (a) a (e) y la global, y escribí 4 fixtures propias en `datos/tmp/rev-s2-r1/`, con cifras a mano para las ventanas, las metas, los instantes y B2.
- No corregí nada. No hice `git add` ni `commit`.

## Verificación (comandos y salidas resumidas; detalle en el informe)
- `go test ./internal/motor_fin/ -count=1 -cover` → ok, 89.2 %. Con `-v`: 20 + 294 PASS, 0 FAIL y 0 SKIP. `gofmt -l .` vacío, `go vet ./...` limpio y 0 `float` fuera de las pruebas.
- `fin_genera_fixtures.py` en una copia → 32 fixtures, sin HALLAZGO; `diff -r` contra el repositorio sin diferencias.
- `fin_run_fixtures.py --out <mío>` → 32 de 32; `cmp` contra `salidas/s2.json` idéntico.
- `motores s2 dump` (exit 0) y luego `compara.py s2` → 3185 de 3185 en las tres comparaciones, 159 de 159, 43 de 43, 256 de 256 y «diferencias: 0».
- Fixtures propias (15 consultas, 78 hojas) → motor 78 de 78 y oráculo 77 de 78. La regla 4 rechaza todo hueco en los dos lados, también un hueco sin devengo que lo cruce. Las listas contiguas se aceptan en el cambio de mes, en el de año y en febrero. Los meses con hueco se siguen aceptando. La única diferencia es `T24:00:00Z` (B5). `rev-b2` confirma que B2 sigue abierto.
- `sha256sum -c` antes y después: `model.json`, `verify3.json` y `data.json` OK; `git diff --quiet` sin cambios.
- Independencia: el oráculo no lee `backend/` y el generador solo lee `model.json` y `data.json`.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → exit 0, «Verificación ampliada OK».

## Supuestos
1. HECHO con lectura propia: la regla de §3 «valores esperados pegados… en vez de leídos de model.json» vale para cualquier prueba Go, no solo para `TestCanon`. Por eso M1 es medio. Las cifras de la interfaz y de M1 que no son hojas de `model.json` sí pueden ir en el código.
2. HECHO con lectura propia: B5 y B6 son bajos, con la misma escala que B1 y B3 de r3: ninguna cifra cambia y ninguna fixture falla.

## Bloqueos
Ninguno nuevo.

## Pendientes (para la corrección r1)
- M1: `cover_test.go:239` debe leer `periodos[0].diferencia_acumulada` de `model.json`.
- B5: el oráculo debe rechazar la hora 24 en `instante`, y `s2-zona-medianoche` debe llevar la consulta `hora-24`.

## Solicitudes fuera de alcance (al integrador)
- Dueño de `S2-interfaz.md`: §5 con 32 fixtures (`s2-ventanas-a-mas-c`, `s2-metas-vencida`) y la consulta A + C en su fixture propia (§5 y §8), lo que cierra B3; §1.1 y §2.1 con el tipo `Instant` (B6); §3.15 con la regla de `invalid_date` (B2, PREGUNTA).
- Dueño de los bloqueos: actualizar el párrafo HECHO de `S2-metas-vencida.md` (B7).

## Propuestas para Axel
- Ninguna nueva. Siguen pendientes de su respuesta `S2-ventanas` pregunta 4 (casos 1 y 2), `S2-metas-vencida`, `S2-ahorro-quincenal` y Q2-12.

## Consumo aproximado
42 acciones de herramienta · unos 175 KB leídos · unos 22 KB escritos (los 2 informes y los temporales borrados).
