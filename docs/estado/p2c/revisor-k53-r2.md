# P2d · revisor-k53-r2 · reporte (12-oct-2026, madrugada)

Agente `revisor-k53-r2` (Opus), revisor independiente del flujo «P2d · K53». Evalué el trabajo de `k53-aplicar` con la corrección `correccion-k53-r1`. Informe completo: `docs/estado/p2c/REVISION-k53-r2.md`.

No corregí nada. No hice `git add` ni commit. No usé Telegram, navegador ni WebFetch. No toqué `docs/design/sistema/`, `docs/design/maquetas/` ni `docs/estado/wfp3/`.

## Hecho

- Revisé los tres hallazgos de la ronda 1 contra el árbol de trabajo:
  - **A-1:** corregido en S5-04, S5-10, §3, §4 y §5 de `S5.md`.
  - **M-1:** corregido con la pregunta (g) en `S5-agente-externo` y su cita en C1.
  - **M-2:** la prosa de `balance_cents` ya cuadra con S2 1.3 §3.11.
- Corrí las verificaciones (a) a (g) del encargo. Todas pasan; la tabla está en el informe.
- **Veredicto: NO APROBADO, con 0 altos, 1 medio y 7 bajos.**
  - **El medio:** `docs/criterios/S5.md` §1 «Objetivo» (l. 35) todavía pide la lectura asistida «desactivada hasta que exista veredicto escrito por sitio». K53 retiró esa condición.
  - **Los bajos:** B-1 a B-4 vienen de la r1 y no estaban en el encargo de la corrección. B-5 a B-7 son nuevos: la pregunta (g) no llegó a los resúmenes, el INDICE de criterios no refleja la corrección de A-1 y los resúmenes de §E en README e INDICE no llevan la nota de K53.

## Verificación (comandos y salidas)

| Comando | Salida |
|---|---|
| `python3 -I /tmp/revk53r2/a.py` y `a2.py` | literal de 446 caracteres presente en `decisiones.md` y en los 2 bloqueos; K37 y K48 empiezan con su texto completo de HEAD |
| `python3 -I /tmp/revk53r2/b.py` | `archivos 85 resueltos 64 abiertos 21`, igual que el INDICE |
| `make verify-c1` | `RESULTADO: los 22 archivos coinciden byte a byte con el generador`, `comprobaciones OK: 33`, `RESULTADO: TODO OK` |
| `make verify-openapi` | `docs/contratos/openapi.yaml: OK` |
| `python3 -I …/build.py --check` | 22 `IDÉNTICO`, rc 0 |
| `uv run --with pyyaml python3 -I /tmp/revk53r2/d.py` | `version 1.4.0 ops 216 rutas 172 esquemas 262`; `ops con serviceToken: 4` |
| `make verify-criterios` | `TOTAL 261 (alta 201 / media 60 / baja 0)`; `INDICE contra conteo propio: 0 diferencias` |
| Los comandos de S5-07, en `bash -c` | 0 / 1 / 0 / 4 / `7:  version: "1.4.0"` |
| `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | `rc=0`; `== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)` |
| `git status --porcelain -- docs/design/sistema docs/design/maquetas docs/estado/wfp3 \| wc -l` | `0` |

## Supuestos

1. **El orden de los criterios de severidad.** Califiqué M-1 como medio y no alto: no decide ninguna pregunta de la puerta 4, no rompe ninguna verificación y en la práctica pide lo mismo que K53, el interruptor apagado. Pero sí deja la condición que K53 retiró en el texto que manda al revisor de S5. Si el orquestador lee «contradice §K» en sentido estricto, sería alto. El veredicto no cambia: NO APROBADO.
2. **B-2, el shell de los comandos.** Supuse que los «Verifica» de los criterios se corren en bash, porque `make` usa `/bin/sh`. El INDICE de criterios no fija un shell. En zsh, `--include=*.go` sin comillas da un falso «sin resultados».
3. **B-5.** La fila del INDICE de bloqueos la escribió `p2c-remanentes`. La cuento como baja porque el propio INDICE dice que el archivo manda.

## Bloqueos

Ninguno nuevo. Ni el revisor ni yo decidimos nada de `S5-agente-externo`, que sigue ABIERTO con las preguntas (a) a (g).

## Pendientes

- **Para la corrección de la r2:** M-1, es decir, §1 de `S5.md` y su línea en «Cambios de la versión 1.4».
- **Para el integrador:**
  - B-5 y B-7: la fila del INDICE de bloqueos, K53 con la pregunta (g), y las notas de K53 en README l. 29 e INDICE de criterios l. 105.
  - B-1 a B-4, si el orquestador los quiere en esta ronda.
- **La PREGUNTA de la r1 sigue para el orquestador:** cuál PROVISIONAL rige para (a), «compite como `manual`» o la regla en uso, «referencia».

## Propuestas para Axel

Ninguna nueva. Las de la puerta 4 están en `docs/estado/bloqueos/S5-agente-externo.md`, preguntas (a) a (g).

## Consumo aproximado

- **Acciones:** unas 32, de 60.
- **KB leídos:** unos 120, de 400. Fueron reportes y revisión de la r1, RESPUESTAS, el bloqueo, extractos de `decisiones.md`, `S5.md`, `S5-interfaz.md`, `S2-interfaz.md` §3.11, `c1_plantilla.md` y los dos INDICE por `grep` y `sed`. No leí `openapi.yaml` ni `C1-datos.md` enteros: los leí con pyyaml y `grep`.
- **KB escritos:** unos 19, de 150: el informe (unos 14) y este reporte (unos 4), más los scripts de `/tmp/revk53r2` (unos 4).
