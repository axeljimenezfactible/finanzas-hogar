# revisor-C1-r1 · Reporte (flujo P1c, ronda 1 de 3)

Revisor independiente de C1 v1.2, 11-oct-2026. No escribió nada de lo revisado ni corrigió nada. El informe completo está en `docs/estado/p1c/REVISION-C1-r1.md`.

## Hecho
- Leí, por partes, el protocolo, §K de decisiones (K1–K44), C1.md 1.2 (cambios y C1-12 a C1-15 y C1-21 a C1-23), el bloqueo abierto `S1-parametros-avisos`, el diff de `c1gen/`, `extractos/INDICE.md`, el reporte `C1-v1.2.md` y C1-datos.md (Cambios, §11.7 y §13). openapi.yaml solo lo consulté con scripts; no lo leí entero.
- Reejecuté las verificaciones (a)–(k) con comandos y scripts propios.
- **Veredicto: NO APROBADO.** 1 alto (A-01, C1-23: `send_time` sin patrón y `advance_days` sin `minimum: 0`), 0 medios y 5 bajos (B-01 a B-05).

## Archivos
- Creados: `docs/estado/p1c/REVISION-C1-r1.md` y `docs/estado/p1c/revisor-C1-r1.md`.
- Temporales: `datos/tmp/rc1/` (base c1.db, check.py, check2.py, check3.py y salidas), borrados al terminar. No toqué `datos/tmp/verify-wf0/` porque no es mío.

## Verificación
- `build.py --check`: IDÉNTICO en los 21 archivos.
- goose v3.28.0 `up` con 00001, 00002 y 00003 (versión 3), luego `uv run --no-project --with pyyaml --with jsonschema python3 -I verify_c1.py --db`: TODO OK (27 OK). Con `python3 -I` del sistema falla C1-14 porque falta jsonschema.
- openapi-spec-validator 0.9.0: openapi.yaml y los 17 extractos OK.
- verificar_u1.py: TODO OK, 0 PENDIENTE.
- go build, go vet y gofmt -l: sin errores; gofmt vacío.
- Extractos: 210 = 210, 0 repetidas, el mayor tiene 76 220 B, INDICE cuadra, sha256 de la carpeta `e7fb8de8…69bfa`.
- K44: las 10 claves de 00003 (y todas las de 00001 a 00003) están en el contrato. Las 15 `@param` tienen 1 lectura con `value` e `initial_value`, 1 cambio y 1 «restablecer», sin `maximum` ni `x-provisional`, y sus comentarios citan su decisión.
- K38 y K43: 0 marcas que citen bloqueos resueltos. Las 3 `x-provisional` citan `S1-parametros-avisos` (abierto).
- C1-23: 413, 503 `unavailable` y M-05 con 21 están bien. No hay patrón de `send_time` ni mínimo de `advance_days` (A-01).

## Supuestos
- Donde C1.md 1.2 (C1-22, excepción de avisos) contradice a K44, mandan K44 y el bloqueo abierto en su versión reducida: los umbrales y los días antes del pago son claves de `app_setting`. Esa diferencia se reporta como B-01 contra los criterios y no se cuenta contra C1.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Corregir A-01 (agente de corrección de C1) y revisarlo en la ronda 2.

## Solicitudes fuera de alcance
1. `docs/criterios/C1.md`: actualizar C1-22 (excepción solo para `send_time`, `advance_days`, `weekday` y `day_of_month`; umbrales y días antes en `app_setting` por K44), C1-13 (C1 v1.2 y 1.2.0), C1-21 (210 operaciones) y C1-12 (excepción para las claves `x-app-setting`). Motivo: hoy contradicen a K44 y al bloqueo reducido (B-01).

## Propuestas para Axel
- Ninguna nueva. Respaldo la del autor: que `alerts.payment_due_days_before` no acepte la lista vacía. Se decide en la puerta 2.

## Consumo aproximado
- Acciones de herramienta: 35 de 60.
- Leído: ~205 KB (nunca enteros openapi.yaml, C1-datos.md, U1, migraciones ni diccionario).
- Escrito: ~17 KB (este reporte y el informe de revisión) más los scripts temporales (~10 KB, borrados).
