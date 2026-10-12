# P2c · revisor-S1-r2 · reporte del agente

Revisor independiente de S1, ronda 2 de 3, 11-oct-2026, noche. El informe completo está en `docs/estado/p2c/REVISION-S1-r2.md`.

## Hecho
- Revisé contra S1.md, K44, K46 a K50, PUERTA-2 §3 y §4 y RESPUESTAS.md:
  - `backend/db/migrations/00004_puerta2.sql`, completa;
  - las partes nuevas de `diccionario-datos.md` 1.4 y de `reglas-negocio.md` 1.4 (RN-03, RN-07, RN-11, RN-12, RN-26 y RN-29);
  - `S1-puerta2.md`, `correccion-S1-r1.md` y el bloqueo nuevo `S1-reparto-edicion-gasto.md`.
- Veredicto: **NO APROBADO** (1 alto, 2 medios y 5 bajos).
  - H-01 sigue abierto: verify-c1 y verify-all fallan por «hora HH:MM».
  - Medios nuevos: N-01 (`vat_included` no está en el mapa de Ajustes, S1-13) y N-02 (sentido no escrito de (0, NULL) en el IVA del crédito).
  - Cerrados: M-01 y M-02.
- No corregí nada. No hice git add ni commit, no mandé nada por Telegram y no usé navegador.

## Verificación
Comandos y salidas en la tabla de `REVISION-S1-r2.md`. En resumen:
- `make verify-sql GOOSE=…` rc=0: 55 tablas, 152 índices, 79 disparadores, integrity ok y fk 0.
- `make verify-s1 GOOSE=…` rc=0: `RESUMEN TODO OK`.
- `make verify-c1 GOOSE=…` rc=2: `gen2: tipo de parámetro desconocido: hora HH:MM`.
- `make -k verify-all GOOSE=…` rc=2: falla solo verify-c1. WF-0, S1, U1, criterios y WF-1 pasan.
- `python3 -I datos/tmp/rev-s1-r2/chk.py`, mi script:
  - 10 `@param` completos, con valor = sembrado = PUERTA-2 §4 y F;
  - 120 casos de rango y formato, 0 discrepancias;
  - `params_json` sin las 4 claves; el Down es byte por byte igual a 1.3 y repone valores cambiados;
  - matriz del CHECK de IVA: 13 casos × INSERT/UPDATE;
  - traspaso de créditos de 1.3 con `updated_at` conservado.
- goose con datos: up-to 3, up, down-to 3 y up, con rc 0 en los cuatro pasos.
- `git diff --stat HEAD` de 00001 a 00003: vacío.
- `gen_dict.py` a un archivo temporal y `cmp`: idéntico.

## Supuestos
1. Leí «Q2-74 mantener 5 s» como valor 5 de `comparator.timeout_seconds`, que la 00003 ya sembraba. La 00004 solo agrega su rango de 1 a 30.
2. No cuento como defecto el caso ('1', NULL) que acepta la tabla STRICT, porque SQLite lo guarda como entero 1.

## Bloqueos
- Ninguno nuevo de mi parte. N-02 puede pedir un bloqueo si el corrector lo considera una decisión: está en la corrección sugerida.

## Pendientes
- Ronda 3: volver a correr `make verify-all GOOSE=…` cuando C1 tome la solicitud 1 (H-01) y revisar N-01, N-02, B-05 y B-06.

## Propuestas para Axel
- Ninguna nueva. La pregunta de `S1-reparto-edicion-gasto` sigue abierta, sin recomendación.

## Solicitudes fuera de alcance
1. Orquestador → C1: H-01, más FALLA K44 y C1-15 de `verify_c1.py`; `vat` en `missing_terms` también con (0, NULL), o 422 (N-02).
2. Orquestador → S3: `S3-interfaz.md` todavía cita como PROVISIONAL `U4-iva-credito`, que ya está resuelto.
3. Integrador: agregar `S1-reparto-edicion-gasto` al INDICE (B-06) y ratificar las ediciones de herramientas (B-04).

## Consumo aproximado
- Acciones de herramienta: ~30.
- KB leídos: ~150 (migración 23, reportes y revisión r1 ~30, PUERTA-2, RESPUESTAS y decisiones ~25, criterios ~8, RN y diccionario por partes ~25, salidas ~40).
- KB escritos: ~32 (informe 14, este reporte 3.5 y script temporal ~12, ya borrado; los .db y el diccionario temporales también se borraron).
