# P2d · revisor-k53-r1 · reporte (12-oct-2026, madrugada)

Revisor independiente de `k53-aplicar` (Opus). No corrijo. Informe completo en `docs/estado/p2c/REVISION-k53-r1.md`. No hice `git add` ni commit. No toqué `docs/design/sistema/`, `docs/design/maquetas/` ni `docs/estado/wfp3/`, ni ningún archivo fuera de `docs/estado/p2c/`.

## Hecho

- Revisé los puntos (a) a (g) del encargo contra RESPUESTAS.md «Puerta 2 bis», §K (K37, K48, K53), los criterios C1 y S5 y el protocolo §4.
- Veredicto: **NO APROBADO**. 1 alto, 2 medios y 4 bajos.
- **A-1 (alto):** S5-04 y S5-10 de `docs/criterios/S5.md` siguen diciendo FIRME que «`site_assisted` solo es referencia hasta el veredicto de términos». Eso contradice K53: el peso es PREGUNTA de la puerta 4. También contradice la fila K37 del INDICE de criterios y S5-interfaz 1.6 §3.
- **M-1 (medio):** límites de lote de 500 y 100, con 422, sin fuente, sin supuesto declarado y sin pregunta en `S5-agente-externo`.
- **M-2 (medio):** la prosa de `balance_cents` (I-6) omite el gasto variable registrado que S2 1.3 §3.11 sí resta.
- **Bajos:**
  - B-1: restos de la lectura asistida en §15 y §16 de S5-interfaz.
  - B-2: S5-08 sigue hablando del módulo de lectura asistida, y al último comando de S5-07 le falta la salida esperada.
  - B-3: error 409 «sin veredicto» en `c1_plantilla.md` l. 157.
  - B-4: precisión del reporte (18 extractos, no 20; `ajustes-f.yaml` nuevo; idempotencia de recetas).

## Verificación (comandos y salidas)

| Comando | Salida |
|---|---|
| `python3 -I`: literal de RESPUESTAS.md contra `decisiones.md` y los 2 bloqueos | True, True, True |
| `python3 -I`: K37 y K48 contra `git show HEAD:` | prefijo intacto; solo agregan la marca de K53 |
| Conteo de bloqueos en disco | 85 archivos: 64 resueltos y 21 abiertos; igual que el INDICE |
| `grep` de los CHECK en `00001_esquema.sql` | `price.source` sin `external_agent` (l. 711); `recipe.origin` cerrado (l. 970) |
| `make verify-c1` | 22 IDÉNTICO; 33 OK; `RESULTADO: TODO OK` |
| `make verify-openapi` | `openapi.yaml: OK` |
| `pyyaml` sobre `openapi.yaml` | 1.4.0; 216 operaciones, 172 rutas, 262 esquemas; `serviceToken` solo en las 4 `external-agent`, con `x-provisional` |
| `verify_c1.py --db` (base temporal) | sin N-xxx: 5; `con serviceToken sin la marca: []`; C1-15: 15 marcas, `sin bloqueo abierto: []` |
| `make verify-criterios` | 261 (201 / 60); 178 comprobadas, 0 diferencias; INDICE con 0 diferencias |
| Los 5 comandos de S5-07, con `bash` | 0, 1, 0, 4 y `7:  version: "1.4.0"` |
| `make verify-all GOOSE=…/goose` | rc 0; «Verificación ampliada OK» |
| `git status --porcelain` de las 3 carpetas vedadas | 0 |

## Supuestos

1. Leí la regla de «opción en uso» de RESPUESTAS.md como aplicable a lo PROVISIONAL de K53. Como la elección «compite» viene del encargo del orquestador y está rotulada como PROPUESTA, no la conté como hallazgo de `k53-aplicar`. La dejé como PREGUNTA para el orquestador.
2. Atribuí cada archivo modificado a `k53-aplicar` o a `p2c-remanentes` según sus reportes y el diff. K46 de `decisiones.md` es de `p2c-remanentes`, porque su propio texto lo dice.
3. Taché A-1 de alto y no de medio: el criterio dice FIRME, en negritas, lo que §K deja como pregunta de la puerta 4.

## Bloqueos

Ninguno nuevo. Ningún hallazgo requiere una decisión de Axel. M-1 propone una pregunta (g) en `S5-agente-externo`, pero la abre quien corrija.

## Pendientes

- Una ronda de corrección de `k53-aplicar` (o de quien la tome) para A-1, M-1 y M-2. Después, `revisor-k53-r2`.
- Los bajos B-1 a B-4, a criterio del corrector.

## Propuestas para Axel

Ninguna nueva. Las del bloqueo `S5-agente-externo` siguen para la puerta 4.

## Consumo aproximado

- Acciones de herramienta: unas 40.
- KB leídos: unos 120, por partes con `grep`, `sed` y scripts. No leí `openapi.yaml`, `C1-datos.md` ni `decisiones.md` enteros.
- KB escritos: unos 19 (este reporte y el informe).
