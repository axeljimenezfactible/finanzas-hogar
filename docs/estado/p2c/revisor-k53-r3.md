# Reporte · revisor-k53-r3 (P2d · K53)

Agente `revisor-k53-r3` (Opus), 12-oct-2026, madrugada. Soy revisor independiente y no corrijo nada. El informe completo está en `docs/estado/p2c/REVISION-k53-r3.md`.

## Hecho

- Revisé el trabajo de `k53-aplicar` con las correcciones r1 y r2. Lo contrasté con:
  - RESPUESTAS.md, «Puerta 2 bis»;
  - decisiones §K (K37, K48 y K53);
  - los criterios C1 y S5;
  - protocolo §4.
- **Veredicto: NO APROBADO, con 0 altos, 1 medio y 7 bajos.**
  - El medio de la ronda 2 (`S5.md` §1) está CORREGIDO.
  - **M-1 nuevo:** C1 1.4 sigue describiendo con la regla anterior a K53 tres campos:
    - `Store.assisted_reading_enabled` y `StoreUpdate.assisted_reading_enabled` dicen «Queda en 0 hasta que docs/estado/S5-terminos.md tenga veredicto para ese sitio» (`openapi.yaml` l. 14076 y 14212).
    - `Price.source` dice que `site_assisted` es «lectura asistida del sitio» (l. 14747).

    Las tres descripciones salen de los comentarios de columna del esquema 1.4 y el generador las copia. Chocan con K53, con `updateStore` (l. 4855), con `Store.assisted_reading_verdict` (l. 14122) y con `importExternalPrices`.
  - Los bajos B-1 a B-7 de la ronda 2 siguen sin cambios. No venían en el encargo de la r2.
- Escribí 2 archivos: `docs/estado/p2c/REVISION-k53-r3.md` y este reporte. Los scripts de verificación están en `/tmp/revk53r3/`, fuera del repositorio.

## Verificación (comandos y salidas)

| Comando | Salida |
|---|---|
| `python3 -I /tmp/revk53r3/a.py` y `a2.py` | literal de 446 caracteres, idéntico en `decisiones.md` y en los dos bloqueos. K37 y K48 empiezan con su texto de HEAD completo y solo agregan la marca de K53 |
| `python3 -I /tmp/revk53r3/b2.py` | `archivos 85 {'RESUELTO': 64, 'ABIERTO': 21, '??': 0}`, igual que el Total del INDICE |
| `grep -rn "CHECK" backend/db/migrations/ \| grep -i "source\|origin"` | l. 711, `price.source`, sin `external_agent`; l. 970, `recipe.origin` = `claude_seed` y `household` |
| `make verify-c1`, `make verify-openapi` y `build.py --check` | rc 0 en los tres. 22 `IDÉNTICO`, `comprobaciones OK: 33`, `TODO OK` y `openapi.yaml: OK` |
| `uv run … /tmp/revk53r3/d.py`, `d2.py`, `d3.py` y `d4.py` | `version 1.4.0 ops 216 rutas 172 esquemas 262`. Solo 4 operaciones con `serviceToken`, cada una con `x-provisional` y con cita a `S5-agente-externo`. Las descripciones viejas de `Store`, `StoreUpdate` y `Price.source` (M-1) |
| `uv run … verify_c1.py` | rc 0, 33 `OK`; C1-02 tiene 5 operaciones sin N-xxx (las 4 y `GET /openapi.yaml`); C1-15 da `sin bloqueo abierto: []` |
| Los 5 comandos de S5-07, con `bash -c` | `0 lineas`; `1`; `0 archivos`; `4`; `7:  version: "1.4.0"` |
| `make verify-criterios GOOSE=…` | rc 0. `TOTAL 261 (alta 201 / media 60 / baja 0)`, `citas verificadas: 321 \| sin fuente: 0` y `INDICE contra conteo propio: 0 diferencias` |
| `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | rc 0. s5 `43 fixtures, 43 coinciden`, y al final `== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)` |
| `git status --porcelain -- docs/design/sistema docs/design/maquetas docs/estado/wfp3 \| wc -l` | `0` |

## Supuestos

1. **M-1 es medio y no bajo.** Es el mismo tipo de residuo que el M-1 de la ronda 2: una condición «hasta que exista veredicto» que K53 retiró, ahora en C1, que es lo que lee WF-3. El encargo de `k53-aplicar` pedía que `Store` dijera «no aplica desde K53», y su reporte lo da por hecho.
2. **El arreglo va en el generador, no en el esquema de S1.** K53 dice que P2d no cambia el esquema de S1. Sustituir la descripción en `gen2_schemas.py` es del alcance de `k53-aplicar`.
3. **K46 no es de `k53-aplicar`.** El cambio en K46 de `decisiones.md` es de `p2c-remanentes`, según su reporte (l. 45 y 80).

## Bloqueos

- No abrí ninguno.
- `S5-agente-externo` sigue ABIERTO, con las preguntas (a) a (g), para la puerta 4.

## Pendientes

- **Para la corrección:** M-1, más B-3, que va junto. El detalle está en el informe.
- **Para el integrador:** los bajos B-1, B-2 y B-4 a B-7.

## Propuestas para Axel

- Ninguna nueva. La pregunta más visible de la puerta 4 sigue siendo la (a) de `S5-agente-externo`: ¿los precios del agente externo compiten como `manual` o son solo referencia, como PROFECO?

## Consumo aproximado

- **Lectura:** unos 130 KB, de 400 permitidos. Incluye:
  - los reportes de `k53-aplicar`, `correccion-k53-r2` y `REVISION-k53-r2`;
  - RESPUESTAS.md;
  - K37, K48 y K53;
  - `S5-agente-externo.md`;
  - extractos de `S5-interfaz.md`, `S5.md`, de los dos INDICE y de `gen1_base.py`;
  - las salidas de los comandos.

  No leí `openapi.yaml` entero: lo leí con pyyaml y con `grep`.
- **Escritura:** unos 24 KB, de 150 permitidos: el informe (16 KB), este reporte y los scripts de `/tmp`.
- **Acciones:** unas 40, de 60 permitidas.
