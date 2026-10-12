# P2c · revisor-S1-r1 · reporte del agente

Revisor independiente (ronda 1 de 3) de S1-puerta2 (esquema 1.4, migración 00004), 11-oct-2026, noche. El informe con veredicto, verificaciones y hallazgos está en `docs/estado/p2c/REVISION-S1-r1.md`.

## Hecho
- Leí por partes:
  - `00004_puerta2.sql` completa;
  - `RESPUESTAS.md`;
  - PUERTA-2 §0, §3.1, §3.8 a §3.12, §3.22 a §3.27 y §4;
  - decisiones F y K44 a K52;
  - `docs/criterios/S1.md`;
  - el reporte `S1-puerta2.md`;
  - el bloqueo `U4-iva-credito`;
  - el diff de `s1_verifica.py` y de `reglas-negocio.md`;
  - las filas nuevas del diccionario.
- Volví a correr `make verify-s1`, `make verify-sql` y `make verify-c1` con GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose.
- Escribí y corrí 2 scripts propios (`python3 -I`, en `datos/tmp/rev-s1r1/`). Cubren goose up/down/up, valores contra PUERTA-2 §4 y F, 34 casos de rango y hora, `params_json` en Up y Down con valores iniciales y cambiados, `updated_at`, el DDL de `credit`, la migración de `vat_included` y la matriz de su CHECK.
- Regeneré el diccionario a un archivo temporal y lo comparé con `cmp`.
- Veredicto: **NO APROBADO**, con 1 alto (H-01, `verify-c1` y `verify-all` rotos por el tipo «hora HH:MM»), 2 medios (M-01, CHECK de `vat_included` con NULL; M-02, pregunta de Q2-68 C sin bloqueo) y 4 bajos.
- No corregí nada, no hice `git add` ni `commit` y no toqué la Pista 3.

## Verificación
Detalle completo en la tabla de `REVISION-S1-r1.md`. Salidas clave:
- `make verify-s1 …`: `RESUMEN TODO OK`, rc=0.
- `make verify-sql …`: 55 tablas, 152 índices, 79 disparadores, integrity_check=ok, foreign_key_check=0 filas y goose sin errores; rc=0.
- `make verify-c1 …`: `gen2: tipo de parámetro desconocido: hora HH:MM` y `Error 1`, rc=2.
- `chk.py`: OK en A-GOOSE, B-VALORES (0 diferencias en 13 claves), B-RANGOS (21 de 21 negativas rechazadas y 13 de 13 positivas aceptadas), C-PARAMS, E-DOWN, E-DOWNVAL, C-UPVAL y D-BACKFILL. D-CHECK: (NULL, 1600) se acepta.
- `chk2.py`: `updated_at` conservado en notification y credit; el DDL de `credit` es el de 1.3 más 106 caracteres de `vat_included`, sin nada borrado.
- `gen_dict.py` a un temporal: `CMP IDENTICO`, sha256 bab3528a…e1b5.
- 00001 a 00003: `git diff --stat HEAD` vacío y sha256 iguales a HEAD.
- Los temporales de `datos/tmp/rev-s1r1/` y los logs `datos/tmp/rev-s1r1-*.log` ya están borrados.

## Supuestos
1. La severidad de H-01 es alta porque así lo pide la regla de la tarea («verificación global rota»), aunque la causa se corrige fuera del ALCANCE de S1, en C1 (c1gen).
2. Para la matriz del CHECK de `vat_included` leí la opción 1 al pie de la letra («`vat_on_interest_bp` solo cuando `vat_included = 0`»).

## Bloqueos
Ninguno nuevo. M-02 propone que el corrector abra `S1-reparto-edicion-gasto.md`; no lo abrí yo porque no corrijo.

## Pendientes
- Ronda 2: comprobar H-01 (requiere la solicitud 1 de C1 asignada por el orquestador), M-01 y M-02, y volver a correr `make verify-all`.

## Propuestas para Axel
- Ninguna nueva. Quedan las del implementador: la sección de `ocr.timeout_seconds` (B-02) y el supuesto 3, créditos de 1.3 sin porcentaje que pasan como «incluye IVA».

## Solicitudes fuera de alcance
1. Orquestador: asignar en P2c la solicitud 1 de `S1-puerta2.md` (c1gen con «hora HH:MM», 10 claves, `vat_included` y `missing_terms` con `vat`), que sostiene H-01.
2. Criterios: que S1-16 contemple `@resuelve` sobre migraciones cerradas (B-03).

## Consumo aproximado
- Acciones de herramienta: 31.
- KB leídos: ~165.
- KB escritos: ~22 (este reporte, el informe y los scripts temporales ya borrados).
