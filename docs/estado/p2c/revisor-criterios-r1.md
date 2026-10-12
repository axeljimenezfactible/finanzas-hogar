# P2c · revisor-criterios-r1 · 11-oct-2026, noche

Revisor independiente, ronda 1 de 3, de los criterios que cambió `criterios-puerta2`. No corregí nada; no hice `git add` ni `commit`; no usé navegador, WebFetch ni Telegram. Informe completo: `docs/estado/p2c/REVISION-criterios-r1.md`.

## Hecho
- Veredicto **NO APROBADO**: 1 alto (H-1), 2 medios (M-1, M-2), 4 bajos (B-1 a B-4).
- H-1: U4-05 vuelve firme la preselección «Solo a los gastos nuevos» de P-26.12, que es el bloqueo abierto `U4-reparto-preseleccion` (sin recomendación; U4-ajustes-b 1.5 la marca PROVISIONAL) y se la atribuye a Q2-68 = C.
- M-1: la Verifica de C1-15 enumera las marcas vigentes de C1 y omite `S1-reparto-edicion-gasto` (abierto, `C1-datos.md` l. 1044).
- M-2: S5-13 dice que las marcas citan «solo» cuatro bloqueos; el motor y el oráculo S5 citan además `S5-motor-oraculo-reglas` y `S5-precio-vigente-captura` (abiertos en P2c). INDICE §6 omite los tres de M-1, M-2 y H-1.
- Lo pedido como mínimo: (a) 11 comprobaciones OK; (b) 261, 201 / 60, ids y severidades iguales a HEAD; (c) los diez valores de la puerta 2 y los rangos coinciden con la 00004 y con PUERTA-2 §4; (d) S5-07 mide `docs/estado/S5-terminos.md` (3 «Veredicto», los 3 «no verificado», no se marca cumplido); (e) ninguna cita vigente a un bloqueo resuelto como abierto (solo en historiales de versión).

## Verificación (comandos y salidas)
- `make verify-criterios` → rc 0; `bash docs/criterios/tools/verif_final.sh` → las 11 secciones con lo esperado (261 definidos, 321 citas, 178 cifras canon, 12 cifras S1-09, matriz C1 212 operaciones, INDICE 0 diferencias).
- `python3 -I datos/tmp/rev-crit-r1/micuenta.py` → 18 archivos, TOTAL 261 {alta 201, media 60}, 261 ids únicos; md5 de «id + severidad» igual a `git show HEAD:` en los 8 archivos tocados.
- `goose -dir backend/db/migrations sqlite3 datos/tmp/rev-crit-r1/r.db up` (v3.28.0) y `python3 -I valores.py` → semilla 00004 = PUERTA-2 §4 en 13 claves; 53 citas «clave = valor» en los criterios, 52 iguales y 1 dato de prueba (S7-10, 31).
- `bash rangos.sh r.db` → los rangos de C1-22/U4-12/S5-04 coinciden con los disparadores (detalle en el informe).
- `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db c1.db` → 33 OK, `RESULTADO: TODO OK`, «claves de puerta 2 con rango distinto del aprobado: []».
- `python3 -I citas_bloq.py` → estado de 85 bloqueos y contexto de cada mención; 0 citas vigentes erróneas.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → rc 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)», 21.6 s.
- Los scripts y bases temporales vivieron en `datos/tmp/rev-crit-r1/` y se borraron al final.

## Supuestos
1. El estado de un bloqueo es su última línea «Estado del bloqueo (puerta 2)» o «Estado: …»; sin ella, «Respuesta de Axel» = resuelto (regla de `bloqueos/INDICE.md` l. 30).
2. Las listas «Cambios de la versión 1.2 …» son historial: lo que dicen de `S1-parametros-avisos` describe su momento y no cuenta para (e).
3. B-2 (K53): lo dejo en bajo porque RESPUESTAS.md asigna el registro de K53 a P2d y lo registrado hoy (K48, bloqueo S5-terminos) dice ABIERTO.

## Bloqueos
Ninguno nuevo.

## Pendientes
- Ronda 2: comprobar H-1, M-1, M-2 y los bajos tras la corrección.
- P2d: al registrar K53, retirar o reformular S5-07 y quitar `S5-terminos` de INDICE §6 (B-2).

## Propuestas para Axel
- PREGUNTA ya abierta por otro agente que el criterio no debe contestar: en P-26.12, ¿qué alcance aparece preseleccionado (A «solo los nuevos», B «también los ya capturados», C ninguno)? Bloqueo `U4-reparto-preseleccion`, puerta 4.

## Solicitudes fuera de alcance
- U2-02 y U2-08 todavía llaman «propuesta» al ingreso en 135° aunque Q2-81 = A lo ratificó (B-4): para la ronda siguiente de criterios o el integrador.
- `docs/estado/S5-terminos.md` está sin seguimiento en git: lo tiene que agregar el integrador.

## Consumo aproximado
Acciones de herramienta: unas 50. KB leídos: unos 200 (diff de criterios por archivo y recortado, PUERTA-2 §2–§5 por secciones, RESPUESTAS, K47–K50, bloqueos por partes, extractos con `grep`, salidas de verificación). KB escritos: unos 20 (este reporte, el informe y scripts temporales ya borrados).
