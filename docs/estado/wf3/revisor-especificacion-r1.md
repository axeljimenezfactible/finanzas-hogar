# revisor-especificacion-r1 · reporte · 12-oct-2026

Agente revisor independiente, ronda 1 de 3 (WF-3c, K54). No corrige. Informe completo: `docs/estado/wf3/REVISION-especificacion-r1.md`. No hice `git add` ni `git commit`, no envié Telegram, no usé sudo y no instalé paquetes.

## Hecho
- **HECHO** Revisé `u3-v2` (U3 2.0, DESIGN 1.1 §8.9, K54, criterios U3, U5 y P3, U1 F-3) y `s1-c1-v15` (00005, diccionario, reglas, C1 v1.5, criterios C1 y S1) solo contra K54, U3.md, U2.md, C1.md, S1.md y §A a §K.
- **HECHO** Veredicto: **NO APROBADO**, con 2 altos, 4 medios y 5 bajos.
  - **Altos.** H-1: `make verify-all` da EXIT 2 porque 5 operaciones de C1 v1.5 no tienen manejador; fallan 4 pruebas de 3 paquetes. H-2: `C1.md` y `C1-datos.md` conservan «≤ 6 tarjetas» como regla vigente.
  - **Medios.** M-1: U3 §2.6 no está alineado con C1 1.5 (400 contra 422, `reset` contra `reset_parameters`, PROPUESTA ya publicada, N de K-4). M-2: U1 sigue con 6 tarjetas. M-3: los pares del barómetro no están en §10.5, la tabla del script. M-4: el comando de C1-13 busca 1.4.0.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/REVISION-especificacion-r1.md` (informe)
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/revisor-especificacion-r1.md` (este reporte)

## Verificación
- **(a) K54.** Citas de Axel idénticas (658 y 293 caracteres) y 4 decisiones fieles: CUMPLE.
- **(b) U3 2.0.** Cubre ingreso con extra, créditos globales, recordatorio N, barómetro 30/40 con dígito y 6 gráficas. 22 métricas, 34 operaciones y 41 campos existen en C1 1.5: CUMPLE, con M-1.
- **(c) Cifras.** 18 montos del diff de U3: 16 en `model.json` y 2 de S5 K-03, que ya estaban. Lo sintético va marcado y lo no decidido va PROVISIONAL con bloqueo: CUMPLE.
- **(d) DESIGN §8.9.** 0 hex; 38 tokens existen; 11 de 11 pares coinciden con mi cálculo (±0.00) y hay 0 bajo umbral: CUMPLE, con M-3 y B-2.
- **(e) 00005.** goose up / down / `up-to 4` (esquema idéntico) / `down-to 0` / up en rc 0; 12 de 12 pruebas de rango; `@param` completos; diccionario regenerado idéntico: CUMPLE, con B-3.
- **(f) verify-all.** `make verify-all` da **EXIT 2** en `verify-wf1-cobertura`. Las etapas posteriores, corridas aparte, dan EXIT 0. C1: 221 operaciones, 176 rutas y 267 esquemas; `verify_c1 --db` TODO OK: **NO CUMPLE** (H-1).
- **(g) Criterios.** 261 = 201 / 60, con ids y severidades iguales a HEAD; `verify-criterios` OK. Pero queda «≤ 6 tarjetas» en `C1.md` §4 y en C1-09: **NO CUMPLE** (H-2).

## Supuestos
1. «DESIGN.md §8.5» de la tarea se evaluó en §8.9, donde `u3-v2` la puso con marca PROVISIONAL y bloqueo (B-2).
2. Lo que está fuera del alcance de los dos agentes y quedó registrado como solicitud (U1, el INDICE de bloqueos) se califica como inconsistencia de artefactos (M-2, B-4), no como falta de ellos.
3. Para no chocar con otros agentes, corrí `make` con `TMP=datos/tmp/rev-esp-r1/…`.

## Bloqueos
Ninguno nuevo.

## Pendientes
- Para la corrección: H-1 (manejadores de `api-tablero` y el reporte de `s1-c1-v15`), H-2, y M-1 a M-4. Los bajos son opcionales.

## Propuestas para Axel
- Ninguna nueva. Los bloqueos `U3-tablero-*` y `C1-tablero` ya llevan sus preguntas a la puerta 4.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 55 acciones, ~190 KB leídos (por secciones con grep y sed) y ~18 KB escritos. Temporales en `datos/tmp/rev-esp-r1/`, borrados al terminar.
