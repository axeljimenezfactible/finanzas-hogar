# Reporte · `revisor-criterios-r2` (P2c, ronda 2 de 3)

Fecha: 11-oct-2026, noche. Revisor independiente de los criterios cambiados en P2c (`criterios-puerta2` y `correccion-criterios-r1`). El informe completo está en `docs/estado/p2c/REVISION-criterios-r2.md`. No hice cambios en los criterios ni en otros artefactos, no hice `git add` ni `commit` y no usé Telegram, navegador ni WebFetch.

## Hecho
- **Veredicto: NO APROBADO** (0 altos, 1 medio y 2 bajos).
- Los hallazgos de la ronda 1 **H-1, M-1 y M-2 están cerrados**. Las bajas B-1 a B-4 siguen, porque no se pidió corregirlas.
- **M-1 (medio, nuevo).** C1-15 define «bloqueo abierto» por archivo: «sin "Respuesta de Axel"» y, en la Verifica, además «"Estado: ABIERTO"».
  - Con esa regla, las marcas que exigen S2-09, S2-10, S2-16 y S5-13 pasan a ser defectos. Esas marcas citan `S2-ventanas`, `S2-flujo-quincenal` y `S5-carrito-comparativa`, que tienen respuesta parcial y siguen ABIERTOS.
  - INDICE §6 usa la misma definición por archivo.
  - Dos de los bloqueos que C1-15 pone como ejemplo (`C1-hora-reportes-job` y `C1-reparto-alcance`) no tienen la línea «Estado».
- Lo que pedía la orden se cumple:
  - (a) `make verify-criterios` pasa las 11 comprobaciones.
  - (b) Hay 261 criterios (201 de severidad alta y 60 media) y ningún id ni severidad cambió.
  - (c) Los 10 valores FIRMES nuevos coinciden con la 00004, con PUERTA-2 §4 y con K50.
  - (d) S5-07 mide `docs/estado/S5-terminos.md` y declara «NO CUMPLE en el fondo».
  - (e) Ningún bloqueo resuelto aparece como abierto en el texto vigente.
  - `make verify-all` termina en rc 0.

## Verificación (comandos y salidas)
- `make verify-criterios`: `TOTAL 261 (alta 201 / media 60 / baja 0) | archivos: 18`; `INDICE contra conteo propio: 0 diferencias`; `residuos «N-xx» y rutas /tmp en los criterios: 0 y 0; las 11 comprobaciones dan lo esperado`.
- `python3 -I datos/tmp/rev-crit-r2/cuenta_rev.py` (mío): `archivos 18 total 261 {'alta': 201, 'media': 60, 'baja': 0} ids unicos 261`; `cambios de ids/severidad contra HEAD: []`.
- `goose -dir backend/db/migrations sqlite3 datos/tmp/rev-crit-r2/t.db up`: `successfully migrated database to version: 4`. El `SELECT` de las claves `reports.*`, `alerts.*`, `optimizer.*`, `ocr.*`, `price.*` y `comparator.*` da los valores de la tabla c1 del informe.
- Pruebas de rango con `UPDATE`. Rechazados: comparador 0 y 31, `price.stale_after_days` 0, día de la semana 8, día del mes 32 y las horas "24:00" y "9:00". Aceptados: comparador 30, `ocr.timeout_seconds` 0 y la anticipación 0.
- `firmes.py` (mío) encontró 35 menciones `` `clave` = valor `` en los criterios. 34 coinciden con la semilla y 1 es el escenario de S7-10 (31 en noviembre). Las 9 diferencias aparentes las causaron mis propias pruebas de rango y las revisé una por una.
- `grep -c "Veredicto" docs/estado/S5-terminos.md` da 3, los 3 «no verificado». `go test ./internal/precios/... -run LecturaAsistida -count=1 -v` da `--- PASS: TestLecturaAsistidaSinVeredicto`.
- `bloq.py`, `citas_e.py` y `marcas.py` (míos) revisaron 84 bloqueos y las marcas de `docs/contratos/`, `backend/internal/precios`, `backend/internal/recetas` y `oraculo/pre_*` y `rec_*`. No hay citas vigentes a un bloqueo resuelto como si siguiera abierto. Las marcas citadas en M-1 están en el informe.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` termina en rc 0: `Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)`.
- Borré `datos/tmp/rev-crit-r2/` al terminar.

## Supuestos
1. El estado de un bloqueo con respuesta parcial lo da su última línea «Estado del bloqueo (puerta 2)». Lo mismo hacen los bloqueos y S2-16.
2. Las listas «Cambios de la versión 1.2» y la §11 de INDICE son historial. Si mencionan un bloqueo ya resuelto como abierto, no lo cuento como defecto.
3. Leí `make verify-all` como la verificación global. Hice la parte de criterios con mis scripts y no volví a revisar los motores.

## Bloqueos
No abrí ninguno.

## Pendientes
- Ronda 3: comprobar el cierre de M-1 (una sola regla por pregunta en C1-15 e INDICE §6) y, si se corrigieron, B-1 y B-2.

## Solicitudes fuera de alcance
1. **Dueño de los bloqueos.** Agregar la línea «Estado» a `C1-hora-reportes-job.md`, `C1-reparto-alcance.md` y `S5-precio-vigente-captura.md`, si la corrección de M-1 conserva esa condición.
2. **Integrador, motor S5.** Quitar la marca vieja de `ocr/types.go` l. 13–15, porque `ocr.timeout_seconds` ya está sembrada y es FIRME. Revisar `oraculo/pre_genera_fixtures.py` l. 254, que cita `S5-reglas-comparativa`, ya RESUELTO.
3. **Integrador, motor S4.** No se tocó en P2c y conserva marcas de `S4-asignacion-mercado` y `S4-reglas-lista`, que están RESUELTOS. S4-15 lo va a reprobar en WF-3.
4. **Contrato S2.** En `S2-interfaz.md` l. 354, la marca de `S2-flujo-quincenal` tiene que decir la pregunta que cita.

## Propuestas para Axel
Ninguna nueva.

## Consumo aproximado
- Acciones de herramienta: unas 46.
- KB leídos: unos 150 (criterios por tramos, bloqueos citados, PUERTA-2 §3.1, §4 y §5, K47 y K50, reportes de P2c y salidas de scripts).
- KB escritos: unos 19 (este reporte, el informe y scripts temporales ya borrados).
