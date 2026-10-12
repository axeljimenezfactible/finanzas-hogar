# revisor-criterios-r3 · P2c · 11-oct-2026, noche

Rol: revisor independiente de los criterios, ronda 3 de 3. No escribí nada de lo que reviso y no corregí nada. No hice git add ni commit, no envié Telegram, no usé navegador ni WebFetch. Base de evaluación: `docs/criterios/tools/verif_final.sh`, PUERTA-2 §4 y §5 y los reportes de P2c (`criterios-puerta2.md`, `correccion-criterios-r1.md`, `correccion-criterios-r2.md` y `REVISION-criterios-r2.md`).

## Hecho
- **Veredicto: APROBADO, con 0 altos, 0 medios y 4 bajos.** El informe está en `docs/estado/p2c/REVISION-criterios-r3.md`.
- **El M-1 de la ronda 2 está CERRADO.** C1-15 (cuerpo y Verifica), INDICE §6, S2-16 y S5-13 usan ahora una sola regla, pregunta por pregunta. La afirmación «Estado: ABIERTO» sobre `C1-hora-reportes-job` y `C1-reparto-alcance` se quitó.
- **Hallazgos bajos:**
  - B-1, que sigue: C1-22 no copia al pie de la letra el texto de `verify_c1.py`.
  - B-2, que sigue: «Q2-01 caso 1 = A».
  - B-3, nuevo: al final del cuerpo de C1-15 queda una oración con la definición vieja de «abierto».
  - B-4, nuevo: S3-08 manda abrir `S3-formula-tarjeta`, que ya está RESUELTO. El bloqueo abierto es `S3-formula-tarjeta-pdf`. S3 quedó fuera del alcance de P2c.

## Verificación (comandos propios; scripts en `datos/tmp/rev-crit-r3/`, ya borrados)
- `make verify-criterios`, rc 0:
  - `criterios definidos: 261 | IDs distintos citados: 261 | citados que no existen: {}`
  - `citas verificadas: 321 | sin fuente: 0`
  - `TOTAL comprobadas: 178 | con diferencias: 0`
  - `cifras comparadas: 12 | diferencias: 0`
  - `TOTAL 261 (alta 201 / media 60 / baja 0) | archivos: 18`
  - `INDICE contra conteo propio: 0 diferencias`
  - `las 11 comprobaciones dan lo esperado`
- `python3 -I cuenta.py` (regex propia y comparación con `git show HEAD:`): `criterios 261 unicos 261 {'alta': 201, 'media': 60}`, `cambios vs HEAD []`.
- `goose … up` en una base temporal da `successfully migrated database to version: 4`. Con un `SELECT` de `app_setting` salen las 10 claves de 00004 con 21:30, 7, 19:00, 1, 08:00, 09:00, 1, 7, 7 y 30. Coinciden con PUERTA-2 §4 (Q2-72 a Q2-76) y con K50.
- Con `UPDATE` de prueba, los rangos coinciden con C1-22: comparador de 1 a 30, S ≥ 1, día de la semana de 1 a 7, día del mes de 1 a 31 y hora HH:MM.
- `python3 -I firmes.py` encuentra 30 menciones de `clave = valor` en los criterios. La única diferencia con la semilla es el escenario de S7-10 (31 en noviembre).
- `go test ./internal/precios/... -run LecturaAsistida -count=1 -v`: `TestLecturaAsistidaSinVeredicto` PASS. `grep -c "Veredicto" docs/estado/S5-terminos.md` = 3. Hay 0 dominios de tiendas en `.go` fuera de las pruebas.
- `bloq.py`, `citas_e.py` y `marcas.py`, que son míos:
  - Los 84 bloqueos quedan en 62 resueltos, 18 sin respuesta y 4 parciales.
  - 0 citas vigentes de un bloqueo resuelto como si siguiera abierto, y 0 al revés.
  - Las marcas de C1 v1.3 citan solo los 4 bloqueos que nombra C1-15.
  - Cada marca vigente que cita un bloqueo parcial nombra la pregunta abierta.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: rc 0 y `Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)`.

## Supuestos
- Tomé como «marca vigente» toda línea con PROVISIONAL fuera de las listas «Cambios de la versión …» y de la tabla de versiones de `C1-datos.md`. Las líneas históricas no las conté como marcas, igual que las rondas 1 y 2.
- Clasifiqué cada bloqueo por la sección «Respuesta de Axel» y por las líneas «Estado del bloqueo (puerta 2): ABIERTO» y «Estado: …» de su propio archivo. No usé `docs/estado/bloqueos/INDICE.md`.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- B-1 a B-4 (bajos), para la próxima pasada de criterios o para P2d. K53 y S5-07 son de P2d.

## Propuestas para Axel
- Ninguna nueva. Las de `criterios-puerta2.md` siguen en pie: S5-07 con K53, el criterio para el agente externo de precios y los criterios propios de Q2-67, Q2-69 y Q2-82.

## Solicitudes fuera de alcance (para el integrador)
1. `docs/contratos/motores/productos-formato.md` l. 67 tiene una marca PROVISIONAL vieja con `S4-asignacion-mercado`, que está RESUELTO.
2. `S2-interfaz.md` l. 354 cita `S2-flujo-quincenal` sin la pregunta.
3. Siguen las marcas viejas de código que anotó r2: `ocr/types.go` l. 13, `pre_genera_fixtures.py` l. 254 y 10 líneas en el motor y el oráculo de S4.

## Consumo aproximado
- Acciones de herramienta: unas 38. KB leídos: unos 120 (salidas de scripts y greps, reportes de P2c, C1-15, INDICE §5–§8 y PUERTA-2 §4 y §5). KB escritos: unos 20 (dos informes y scripts temporales).
