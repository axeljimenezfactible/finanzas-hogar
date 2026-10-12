# Corrección k53 · ronda 2 (P2d · K53)

Agente `correccion-k53-r2` (Opus), 12-oct-2026, madrugada. Corrijo el único hallazgo que el orquestador me encargó de `docs/estado/p2c/REVISION-k53-r2.md`: M-1 (medio). Tengo el mismo alcance que `k53-aplicar`. Reglas:
- No decidí ninguna pregunta de la puerta 4.
- No toqué los bajos B-1 a B-7, que no venían en el encargo.
- No hice `git add` ni commit. No usé Telegram, navegador ni WebFetch.
- No toqué `docs/design/sistema/`, `docs/design/maquetas/` ni `docs/estado/wfp3/`, ni el alcance de `p2c-remanentes`.

## Hecho

### M-1 (medio) · `docs/criterios/S5.md` §1 «Objetivo»

Un solo archivo modificado: `docs/criterios/S5.md`. Sigue en la versión 1.4, con los mismos 13 criterios (10 alta y 3 media).

**§1 «Objetivo» (l. 36; antes l. 35)**
- Antes: «… de modo que reproduzca la comparativa del reporte y mantenga la lectura asistida de sitios **desactivada** hasta que exista veredicto escrito por sitio.»
- Ahora: «… de modo que reproduzca la comparativa del reporte y no lea sitios de tiendas (K53: la lectura asistida no aplica, el interruptor queda apagado y `docs/estado/S5-terminos.md` es registro histórico).»
- Es el texto de la corrección sugerida, sin cambios.
- Con esto §1 ya no deja abierta la vía «veredicto → activar». Queda igual que K53 (1), S5-07 («No aplica desde K53»), S5-04 («ya no aplica la condición…») y §3 («Cualquier lectura de sitios de tiendas desde la aplicación (K53)… no aceptable»).

**«Cambios de la versión 1.4» (l. 8):** agregué una viñeta nueva, «§1 «Objetivo» (K53; corrección `correccion-k53-r2`, hallazgo M-1 de `docs/estado/p2c/REVISION-k53-r2.md`)». Dice qué cambió y que no cambia ningún criterio, ninguna severidad ni ningún comando de «Verifica».

**Encabezado (l. 3):** «y corrección `correccion-k53-r1`» pasa a «y correcciones `correccion-k53-r1` y `correccion-k53-r2`», como hizo la r1 (ver Supuesto 1).

**Diff completo contra el respaldo** (`/tmp/correccion-k53-r2/S5.md.bak`, la copia previa a mi edición): 3 bloques, que son l. 3, la viñeta nueva de l. 8 y §1. No cambió nada más del archivo.

## Verificación (comandos y salidas)

| Comando | Salida | Resultado |
|---|---|---|
| `python3 -I /tmp/correccion-k53-r2/edit.py` (cada reemplazo exige 1 sola aparición) | `OK /home/deepsec/finanzas-hogar/docs/criterios/S5.md 3 reemplazos` | OK |
| `diff /tmp/correccion-k53-r2/S5.md.bak docs/criterios/S5.md` | `3c3` (encabezado), `7a8` (viñeta nueva) y `35c36` (§1) | OK |
| `grep -n -o "…veredicto escrito…\|…hasta que exista veredicto…\|…hasta veredicto…" docs/criterios/S5.md` | Quedan 3 menciones: l. 8 (la viñeta nueva, que cita el texto retirado), l. 64 (S5-07, «Texto anterior, como registro: «Lectura asistida desactivada hasta veredicto»») y l. 113 (§5, cita literal de §E: «… lectura asistida… desactivada hasta tener veredicto escrito… PROFECO», que comprueba `citas.py`). Ninguna pide hoy esperar un veredicto | OK |
| `make verify-criterios GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | `archivos: 18 \| criterios definidos: 261`; `citas verificadas: 321 \| sin fuente: 0`; `TOTAL comprobadas: 178 \| con diferencias: 0`; `cifras comparadas: 12 \| diferencias: 0`; `TOTAL 261 (alta 201 / media 60 / baja 0) \| archivos: 18`; `INDICE contra conteo propio: 0 diferencias`; «las 11 comprobaciones dan lo esperado»; rc 0 | OK |
| `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` (log en `/tmp/correccion-k53-r2/verify-all.log`, 147 renglones) | `rc=0`. goose v3.28.0: «up, down-to 0 y up otra vez sin errores»; `go build` y `go vet` sin errores. C1: 22 `IDÉNTICO`, «los 22 archivos coinciden byte a byte con el generador», `comprobaciones OK: 33`, `RESULTADO: TODO OK`. S1: `RESUMEN TODO OK`. U1: `comprobaciones OK: 63 · INFO: 5 · PENDIENTE: 0`, `RESULTADO: TODO OK`. Criterios: las mismas cifras de arriba. Oráculos: s2 32/32, s3 27 fixtures, s4 26/26 y s5 `43 fixtures, 43 coinciden`, con 0 diferencias. `pre_verify.py: 15 comprobaciones OK`. Al final: `== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)` | OK |
| `git status --porcelain -- docs/design/sistema docs/design/maquetas docs/estado/wfp3 \| wc -l` | `0` | OK |
| `find . -newer /tmp/correccion-k53-r2/S5.md.bak` (sin `.git` ni `datos/tmp`) | `docs/criterios/S5.md` y `docs/contratos/C2-tokens.json`. El segundo lo reescribe `verify-all` (`generar_c2_tokens.py`), y el log dice «C2-tokens.json idéntico byte a byte». `git status --porcelain` no lo marca | OK |

## Supuestos

1. **Encabezado de la versión.** Agregué `correccion-k53-r2` al renglón de versión (l. 3), como hizo `correccion-k53-r1` con su nombre. El encargo solo pedía anotarlo en «Cambios de la versión 1.4». Es un cambio de registro, no de contenido.
2. **«Entregar … `docs/estado/S5-terminos.md`» en §1.** Lo dejé, porque la corrección sugerida solo cambia la cláusula final. Por K53 el documento existe y se conserva como registro histórico, y el paréntesis nuevo lo dice. Quitarlo de los entregables iría más allá del hallazgo.
3. **l. 64 y l. 113.** No las toqué. La l. 64 es el «Texto anterior, como registro» de S5-07. La l. 113 es una cita literal de §E que comprueba `citas.py`, y K53 decidió no editar §E. El resumen de §E que sí pide una nota (B-7) está en `README.md` y en `docs/criterios/INDICE.md`, no en `S5.md`, y no venía en el encargo.

## Bloqueos

- No abrí bloqueos nuevos.
- `S5-agente-externo` sigue **ABIERTO**, con las preguntas (a) a (g). No lo toqué.

## Pendientes

- **Bajos B-1 a B-7 de `REVISION-k53-r2.md`.** No venían en este encargo:
  - B-1: `S5-interfaz.md` 1.6 §15 y §16.
  - B-2: S5-08, la salida esperada del último comando de S5-07 y `--include='*.go'` con comillas.
  - B-3: `c1_plantilla.md` l. 159.
  - B-4: la precisión del reporte de `k53-aplicar`.
  - B-5: la pregunta (g) en la fila del INDICE de bloqueos y en K53. Es para el integrador.
  - B-6: `docs/criterios/INDICE.md`, la fila K53 y la l. 3.
  - B-7: `README.md` l. 29 y `docs/criterios/INDICE.md` l. 105.
- **PREGUNTA del revisor, para el orquestador (sin cambios desde la r1).** Hay dos candidatos para lo PROVISIONAL de (a):
  - «compite como `manual`», la PROPUESTA del encargo de P2d, que usan `S5.md`, K53 y `S5-interfaz` 1.6;
  - «referencia», la opción en uso del motor y de la fixture.

  RESPUESTAS.md no fija cuál rige para K53.

## Propuestas para Axel (puerta 4)

- No tengo propuestas nuevas. Siguen abiertas las preguntas (a) a (g) de `S5-agente-externo`. La más visible es la (a): ¿los precios del agente externo compiten como `manual` o son solo referencia, como PROFECO?

## Consumo aproximado

- **Lectura:** unos 80 KB, de 400 permitidos: la revisión r2 (14 KB), `S5.md` (28 KB), el reporte de la r1 (11 KB), los renglones K37, K48 y K53 de `decisiones.md`, la «Puerta 2 bis» de RESPUESTAS.md, el inicio de `citas.py` y extractos de los logs.
- **Escritura:** menos de 1 KB en `S5.md` y unos 6 KB en este reporte, de 150 permitidos.
- **Acciones:** unas 14, de 60 permitidas.
