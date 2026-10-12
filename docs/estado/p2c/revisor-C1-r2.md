# P2c · revisor-C1-r2 · reporte

Agente: `revisor-C1-r2`, revisor independiente (ronda 2 de 3). Flujo P2c, 11-oct-2026, noche. No corregí nada y no hice `git add` ni `git commit`. Solo escribí este reporte y `docs/estado/p2c/REVISION-C1-r2.md`.

## Hecho
- **Veredicto: NO APROBADO.** 0 altos, 1 medio (M-1) y 5 bajos (B-1 a B-5). El detalle está en `REVISION-C1-r2.md`.
- **Hallazgos de la ronda 1:**
  - H-1: CERRADO.
  - H-2: CERRADO.
  - H-3: CERRADO en lo que pedía, el caso `vat_included = true`.
- **M-1 (media).** El CHECK de la 00004 también rechaza `vat_on_interest_bp` con valor cuando `vat_included` es NULL o ausente. Pasa en dos casos:
  - un alta con solo `{"vat_on_interest_bp": 1600}`;
  - el Merge Patch `{"vat_included": null}` sobre un crédito con (false, 1600).

  C1 no dice qué responde la API en ninguno de los dos. Lo comprobé con `sqlite3` sobre una base migrada por mí.
- **Lo que pedía el encargo, con evidencia propia:**
  - (a) `verify-c1`, `verify-openapi`, `build.py --check` y `verify-all` pasan, con un TMP propio.
  - (b) Las 10 claves nuevas de S1 y las 3 de `@param_rango` están en sus operaciones de settings, con los rangos de los disparadores. `NotificationUpdate` no tiene `send_time` ni `advance_days`.
  - (c) `vat_included`, `scope`, `assigned_store_id` y los acumulados de `AccrualVsCash` están y están documentados.
  - (d) `info.version` y `C1_VERSION` son 1.3.0.
  - (e) Hay 168 rutas, 212 operaciones y 255 esquemas, iguales en todos los artefactos.
- **Extractos.** Comparé los 17 con `openapi.yaml` normalizando los `$ref`: son idénticos, con 212 operaciones y 304 copias de esquemas.

## Verificación (comandos y salidas)

| Comando | Salida |
|---|---|
| `make verify-c1` (TMP compartido) | Falla por el entorno: «grep: …/datos/tmp/verify-wf0/verify_c1.log: No such file or directory», `Error 2`. Lo borró otra corrida en paralelo (B-4) |
| `make verify-c1 TMP=…/datos/tmp/rev-c1-r2-mk` | «los 21 archivos coinciden byte a byte», «comprobaciones OK: 32», «RESULTADO: TODO OK» |
| `make verify-openapi` | «docs/contratos/openapi.yaml: OK» |
| `python3 -I …/c1gen/build.py --check` | 21 archivos IDÉNTICO. `openapi.yaml` de 702,557 bytes y `routes.go` con 212 rutas |
| `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose TMP_ROOT=…/datos/tmp/rev-c1-r2/va` | rc=0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)» |
| Scripts propios con `yaml.safe_load` (`rev.py` a `rev8.py`) | Conteos; claves por operación; rangos; propiedades con «_» final (0 de 610); esquemas de crédito, reparto, lista y comparativa; extractos contra `openapi.yaml` (0 diferencias); `x-provisional` (6, todos con bloqueo abierto) |
| Recálculo de M-11 desde `oraculo/model.json` con `Decimal` | P1 1947.50 / 973.75 / 973.75; P2 3895.00 / 2921.25 / 973.75; P6 14362.48 / 13103.78 / 1258.70. Coincide con la fila de `C1-datos.md` l. 372 |
| Base propia: Up de 00001 a 00004 con `sqlite3`, y `INSERT` o `UPDATE` de `credit` | (NULL, 1600) rechazado; (1, 1600) rechazado; (0, 1600) aceptado; `vat_included` → NULL sobre (0, 1600) rechazado |

## Supuestos
1. **La falla de V-1 viene de otra corrida en paralelo.** El log desapareció dentro de la misma receta, entre la escritura y la lectura, y la misma verificación pasó con un TMP propio. No revisé la lista de procesos.
2. **M-1 es medio.** Uso el mismo criterio que dio «media» a H-3 en la ronda 1: quien implemente S6 tendría que adivinar, y la violación puede salir como 5xx. Que la API deduzca `vat_included = false` sería una regla nueva, y no la propongo como corrección técnica.
3. **Criterios.** Evalué contra `docs/criterios/C1.md` 1.2, la versión vigente. Lo específico de la puerta 2 lo evalué contra `PUERTA-2.md` §2 y §3 y contra `RESPUESTAS.md`.

## Bloqueos
No abro ninguno nuevo: M-1 es una corrección técnica dentro de una respuesta de Axel (Q2-67 opción 1). Siguen abiertos y bien citados `C1-reparto-alcance`, `C1-hora-reportes-job`, `S3-iva-lectura` y `S1-reparto-edicion-gasto`.

## Pendientes
- **Para la ronda 3:** corregir M-1 y, si se encargan, B-1 a B-3.
- **Siguen abiertos desde r1 (B-5):** `docs/criterios/C1.md` a 1.3; `bloqueos/INDICE.md` con los 2 bloqueos nuevos de C1; el conteo 29 → 32 de `C1-puerta2.md`; y los textos B-1 y B-2 de r1.
- **K53** (importación del agente externo, retiro de la lectura asistida): queda para P2d, no para esta ronda.

## Solicitudes fuera de alcance
1. **Integrador (`Makefile`):** un TMP por corrida o `TMP_ROOT` obligatorio en P2c (B-4). Con agentes en paralelo, `make verify-c1` falla por el entorno.
2. **S6 (WF-3):** además del caso `true` que ya pidió `correccion-C1-r1`, validar antes de escribir el caso con `vat_included` NULL o ausente y un porcentaje, para que el CHECK de la 00004 nunca salga como 5xx.

## Propuestas para Axel
- Ninguna decisión nueva hace falta para cerrar M-1: la corrección documenta la respuesta 422 sin cambiar valores.
- PREGUNTA, opcional y para la puerta 4: ¿quieres que la API deduzca «no incluye IVA» cuando se captura solo el porcentaje? Hoy C1 dice que no cambia un campo por su cuenta.

## Consumo aproximado
- Unas 50 acciones de herramienta.
- Unos 150 KB leídos, por partes: reportes `C1-puerta2` y `correccion-C1-r1`, `PUERTA-2.md` §2, §3.1 y §3.13 a §3.25, fragmentos de la 00004, de `C1-datos.md` y de `S4-interfaz.md`, y salidas de mis scripts. No leí `openapi.yaml` entero: solo con `yaml.safe_load` y `grep`.
- Unos 21 KB escritos: `REVISION-C1-r2.md` (15 KB) y este reporte. Además, scripts y una base temporal en `datos/tmp/rev-c1-r2/`, ya borrados.
