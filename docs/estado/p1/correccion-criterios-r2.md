# Reporte · correccion-criterios-r2 (11-oct-2026)

## Hecho
- R2-M01 (media, K8 FIRME), ambos puntos:
  - `docs/criterios/S2.md` S2-06: la quincena es subperiodo con los límites de K8 (día 1 al 15; 16 al último día, con el 31 en la segunda; semana de lunes a domingo). El día de pago (15 y 30, o el último día; el 28 en febrero de 2027) queda como otro dato. Se quitó «la quincena de fin de mes es el 30, no el 31».
  - La Verifica de S2-06 agrega el caso de 31 días: un gasto de fixture del 2026-10-31 cae en la quincena 2 de octubre (`2026-10-16` a `2026-10-31`) y no en otra; el pago del 2026-10-30 también cae en la quincena 2; un motor con la quincena 2 cerrada el 30 falla. Incluye también febrero 2027 (16 a 28), el mes de 30 días y la semana del miércoles 14-oct-2026 (`2026-10-12` a `2026-10-18`). La Evidencia pide los 4 casos.
  - S2: viñeta nueva en «Cambios de la versión 1.1» y traza de K8 (cita textual) en la §5 hacia S2-06.
  - `docs/criterios/S1.md` S1-03: «quincena 15 y 30» pasa a «día de pago 15 y 30»; se agrega la regla de periodos de K8 con ejemplo ilustrativo rotulado (sábado 31-oct-2026 en la quincena 2; semana del 14-oct-2026 = 12 al 18 de octubre). Viñeta en la lista de cambios y S1-03 agregado a la traza de K8 en la §5.
- Los 6 hallazgos bajos (R2-B01 a R2-B06) no estaban en el alcance de esta corrección y no se tocaron.

## Archivos
- /home/deepsec/finanzas-hogar/docs/criterios/S2.md (parches con str.replace)
- /home/deepsec/finanzas-hogar/docs/criterios/S1.md (parches con str.replace)
- /home/deepsec/finanzas-hogar/docs/estado/p1/correccion-criterios-r2.md

## Verificación
- `bash docs/criterios/tools/verif_final.sh`:
  - cuenta.py: 134 (105 alta / 29 media / 0 baja) en S y WF4; incompletos ninguno.
  - traza.py: 18 archivos, 257 criterios; 0 sin línea en su §5; sin secciones faltantes.
  - refs.py: 257 definidos, 0 citados inexistentes.
  - citas.py: 276 citas (275 antes: la cita nueva de K8 en S2 tiene fuente), 5 sin fuente (los mismos de H-07, fuera de alcance).
  - canon.py: 178 cifras, 0 diferencias; sha256 de model.json `1a5706…bfd4` sin cambio.
  - s109_sql.sh: UP OK, 18 filas, q2 oct `2026-10-16`→`2026-10-31`, 0 quincenas fuera de límites, semana 12→18 oct, integrity_check ok; s109_equiv.py: 12 de 12, 0 diferencias.
  - matriz_c102.py: igual que antes (no tocado). Rutas /tmp en S y WF4: 0.
- Los demás comandos del prompt (build.py, verify_c1, goose, make verify-color, etc.) no aplican: solo cambió texto de criterios.

## Supuestos
- Ninguno que cambie cifras. El valor esperado del gasto del 31-oct-2026 lo declara la fixture (el oráculo no tiene eventos el día 31); no se inventó ninguna cifra.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- R2-B01 a R2-B06 (bajos) y H-07 (5 citas sin fuente), fuera de este alcance.

## Solicitudes fuera de alcance
- Ninguna.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Acciones de herramienta: 8. Leído: ≈ 60 KB. Escrito: ≈ 6 KB.
