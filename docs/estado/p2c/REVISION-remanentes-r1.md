# P2d · REVISION-remanentes-r1 · revisión del trabajo de `p2c-remanentes` (12-oct-2026, madrugada)

Revisor `revisor-remanentes-r1` (Opus), flujo «P2d · K53». Revisión independiente: no corregí nada. Objeto: `docs/estado/p2c/p2c-remanentes.md` contra M-1 de `REVISION-contratos-r3.md`, M-1 (y B-2, B-5, B-6) de `REVISION-motor-S5-r3.md`, I-1, I-3, I-4 e I-5 de `integrador-p2c.md` (INTEGRACION §5) y los criterios S5 y S1. HEAD = `5936414` [P3d] (Pista 3); comprobé que [P3d] no toca ninguna ruta de P2d (`git diff --quiet 476f4a9 5936414 -- backend oraculo docs/contratos docs/estado/bloqueos docs/plan docs/design/DESIGN.md docs/design/U4-ajustes*.md docs/design/entrada`: sin cambios; 0 archivos de [P3d] fuera de la Pista 3).

## Veredicto

**APROBADO: 0 altos, 0 medios y 2 bajos.**
- **M-1 del motor S5 (`created_at`): cerrado.** Probé los 6 textos de M-1, uno más y 25 sondas extra (32 fixtures mías): motor y oráculo dan exactamente la misma salida en las 32 (`json` idéntico). Las 7 fixtures nuevas fallan contra el lado viejo que corresponde (lo reconstruí desde `476f4a9`).
- **M-1 de contratos (registro): cerrado.** INDICE: 85 archivos, 64 RESUELTOS y 21 ABIERTOS, igual a mi conteo del disco. La suma de filas cuadra. Los 13 de P2c están en el INDICE y en K46 (84/63/21 al cierre de P2c). Las 3 remisiones están en su lugar.
- **I-1, I-3 e I-4: cerrados. I-5: cerrado en los sitios que lista.** Quedan marcas del mismo tipo en `oraculo/rec_*.py` y `oraculo/README.md` (B-1). Están declaradas como Pendiente y en la lista para WF-3.
- `make verify-all`: EXIT=0. `git status`: nada en `docs/design/sistema`, `docs/design/maquetas` ni `docs/estado/wfp3`.

## Verificaciones (evidencia propia, desde `/home/deepsec/finanzas-hogar`)

| # | Verificación | Resultado | Evidencia |
|---|---|---|---|
| a1 | INDICE = disco | CUMPLE | Mi script `/tmp/rev-rem-r1/estados.py` toma la última línea «Estado del bloqueo …» o «- **Estado:**» de cada archivo. Resultado: 85 archivos, `RESUELTO 64`, `ABIERTO 19` y 2 sin línea (`P3-contraste-chip-hover` y `P3-huecos-componentes`, abiertos por su texto). Total: 64/21. INDICE «**Total** \| **85** \| **64** \| **21**». Filas: 34+8+1+26+13+1+2 = 85; resueltos 34+8+1+21 = 64; abiertos 5+13+1+2 = 21 |
| a2 | Los 13 de P2c registrados | CUMPLE | Cada uno de los 13 nombres aparece al menos una vez en el INDICE y una vez en `decisiones.md` l. 120 (K46: «84 archivos, 63 resueltos y 21 abiertos, con los 13 bloqueos nuevos de P2c»). Las secciones «Bloqueos nuevos de P2c» (4) y «Dados de alta al cierre de P2c» (9) suman 13. La fila P2c dice 13/0/13 |
| a3 | Regla de las herramientas («67 y 18») | CUMPLE | Archivos con «Respuesta de Axel» (o «Resolución» en `criterios-nxx`): 67 de 85; sin ella: 18 (los 13 de P2c, `S5-agente-externo`, `S5-carritos-detalle`, `S5-profeco-filtro` y los 2 de P3) |
| a4 | Remisiones y línea de estado (I-4) | CUMPLE | `S3-formula-tarjeta.md` l. 47, `S4-confirmacion-empates.md` l. 22 y `U4-iva-credito.md` l. 21 traen «Remisión (P2c …)» al bloqueo abierto. `C1-hora-reportes-job.md` y `C1-reparto-alcance.md` terminan (l. 18 de 18) con `- **Estado:** ABIERTO` |
| a5 | P2d en el INDICE | CUMPLE | Fila «Flujo P2d» 1/0/1. Sección «Bloqueos nuevos de P2d» con `S5-agente-externo`. `S5-terminos` dice «RESUELTO por K53» en l. 381, 497 y 527, y en su archivo («puerta 2 bis»). Las menciones de «abierto» en l. 12 y 356 son párrafos fechados de la puerta 2 |
| b1 | Los 6 textos de M-1 en motor y oráculo | CUMPLE | Sondas `datos/tmp/rev-rem-r1/fx/rv-m1-*`: `manual` 2000 a las 12:30:00Z contra `ticket_manual` 2100 con el texto. `T7:00:00Z`, `T7:00:00-06:00`, `12:30:00,5Z`, `…Z\n`, todo en dígitos U+0660 a U+0669 (y solo la hora con uno), `0000-01-01T00:00:00Z`: **`invalid_date` en los dos** |
| b2 | Controles válidos y sondas extra | CUMPLE | Iguales en los dos lados. Válidos: `07:00:00-06:00` → 2100 `ticket_manual`; `12:00:00+00:00` → 2000 `manual`; `-00:00`, 10 dígitos de fracción y `9999-12-31T23:00:00-02:00` → 2100; `0001-01-01T00:30:00+01:00` → 2000. `invalid_date` en los dos: `.Z`, `t…z`, espacio o `\r` final, salto inicial, `+0600`, dígitos de ancho completo, `0000-12-31T23:30:00-01:00` y 29-feb-2026. `observed_date`: `0000-01-01`, `\n` final, `20261010`, `2026-W41-6`, `2026-283`, arábigo-índicos y 29-feb-2026 → `invalid_date` en los dos; `0001-01-01` → 2000 en los dos. `today` `0000-01-01` y con `\n` → `invalid_date` en los dos. **32 de 32 iguales; `or.json == mo.json`** |
| b3 | Las 7 fixtures detectan el código viejo | CUMPLE | `git archive 476f4a9` en `/tmp/rev-rem-r1/head`. El oráculo viejo falla en `salto-linea` y `digitos-no-ascii`. El motor viejo da `__falta__` en `anio-cero`, `coma-decimal`, `fecha-anio-cero`, `hora-un-digito` y `-zona` (`compara.py s5 --fixtures …/nuevas`: 7 diferencias) |
| b4 | Las fixtures se regeneran idénticas | CUMPLE | `pre_genera_fixtures.py` en una raíz aislada (`datos/tmp/rev-rem-r1/raiz`, sin `fixtures/s5`): «43 fixtures». `diff -r` contra `oraculo/fixtures/s5`: vacío (también `_datos/`) |
| b5 | `run_fixtures` reproduce `salidas/s5.json` | CUMPLE | «43 fixtures, 43 coinciden, 1366 hojas de esperado, 0 diferencias». `cmp` contra `oraculo/salidas/s5.json`: idéntico |
| b6 | `compara.py s5` | CUMPLE | Con mi binario (`go build -o datos/tmp/rev-rem-r1/motores ./cmd/motores`, `s5 dump`): oráculo contra motor 963/963, provisionales 317/317, **diferencias: 0** |
| b7 | gofmt, vet, pruebas y cobertura (S5-13, K25 ≥ 85 %) | CUMPLE | `gofmt -l internal/precios internal/recetas cmd/motores`: vacío. `go vet`: limpio. Cobertura: precios 96.0 %, asistida 100.0 %, ocr 99.3 %, profeco 95.3 %, recetas 93.2 %; total 90.4 %. `TestEvalCreatedAtGramatica`, `TestEvalErrores` y `TestEvalCreatedAtZona`: PASS |
| b8 | Bajos del motor S5 | CUMPLE (B-2, B-5, B-6) | `compare.go` l. 165–166 cita `S5-carritos-granel` (opción A). `profeco/filter.go` l. 45 cita `S5-profeco-filtro` (preguntas 1 a 3). `profeco/links.go` l. 50–51 cita la pregunta 2. `S5-motor-oraculo-reglas.md` l. 77 trae la cifra de r2. `types.go` l. 4 dice «1.5». B-3 sigue abierto; es de `k53-aplicar` |
| c1 | Paquete: `sha256sum -c` | CUMPLE | 21 renglones del MANIFIESTO, **21 OK** |
| c2 | `/api/v1` en los PDF | CUMPLE | `pdftotext`: `/api/v1` = 0 y nombres de operación = 0 en DESIGN, U3-dashboard, U4-ajustes y C1-resumen |
| c3 | El paquete corresponde a las fuentes actuales | CUMPLE | Recalculé los 5 sellos: hash del script más las fuentes, con el mismo filtro `perl` de `filtrar_api`. DESIGN, DESIGN-docx, U3-dashboard, U4-ajustes y C1-resumen son **iguales** a los de `.cache`. U4 filtrado: 0 `/api/v1`. U4 l. 490 conserva «(N-296 y N-297)» |
| d1 | DESIGN.md l. 933 (I-3) | CUMPLE | «**FIRME (Q2-11 = B, K47; bloqueo … S1-parametros-avisos.md RESUELTO)**». Solo la parte del `schedule` es PROVISIONAL, y cita el bloqueo abierto `C1-hora-reportes-job` |
| d2 | Ninguna marca PROVISIONAL cita un bloqueo RESUELTO (alcance de `p2c-remanentes`) | CUMPLE | `/tmp/rev-rem-r1/ventana.py` (cita en el renglón y los 2 siguientes) sobre todos los `.go` de `precios` y `recetas`, `pre_comun.py`, `pre_genera_fixtures.py`, `pre_calculo.py`, `pre_profeco.py`, `productos-formato.md`, `S2-interfaz.md` y `DESIGN.md`. Revisó 56 renglones PROVISIONAL. Única coincidencia: DESIGN l. 933, que es d1. Los 10 sitios de I-5 que tocó dicen FIRME con su Q2 (los leí) |
| d3 | Lo mismo en todo el repositorio | OBSERVACIÓN (B-1) | `/tmp/rev-rem-r1/marcas.py`, un renglón a la vez, sobre contratos, diseño, `backend`, `oraculo` y criterios: 116 renglones. Casi todos son bitácoras, criterios que dicen «RESUELTO», la migración 00003 o renglones que también citan el bloqueo abierto. Marcas activas viejas fuera del alcance de `p2c-remanentes`: `motor_cred` y `oraculo/cred_*` (lista para WF-3 de `correccion-contratos-r2`), `oraculo/rec_*.py` y 3 fixtures s4 (B-1), `oraculo/README.md` l. 37–38 (B-1) y `c1gen/gen6_c1.py` l. 46 (C1) |
| e | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | CUMPLE | **EXIT=0**, 21.7 s. «RESUMEN TODO OK»; «TOTAL comprobadas: 178 \| con diferencias: 0»; «s5 · 43 fixtures, 43 coinciden … salidas/s5.json idéntica»; «sha256sum -c: 21 archivos sin diferencias»; «Verificación ampliada OK». `git status` igual antes y después |
| f | `git status` | CUMPLE | 0 rutas en `docs/design/sistema`, `docs/design/maquetas` y `docs/estado/wfp3`. Todo lo demás está en uno de los dos alcances. Hay 5 rutas que `p2c-remanentes` dice no haber tocado (`README.md`, `docs/criterios/C1.md`, `INDICE.md`, `S5.md` y `docs/plan/plan.md`). Su mtime (02:01:09 a 02:01:28) es posterior a su reporte (02:00:29). Corresponden a `k53-aplicar` (K53 en criterios y plan) |

## Hallazgos

**B-1 (baja) · Siguen marcas viejas del tipo de I-5 en `oraculo/`.**
- Rutas:
  - `oraculo/rec_plan.py` l. 3, 151 y 213 (`S4-reglas-lista` y `S4-asignacion-mercado`, RESUELTOS);
  - `oraculo/rec_genera_fixtures.py` l. 225, 235, 395 y 644, más las descripciones de `s4-canal-mercados`, `s4-lista-reposicion-ritmo` y `s4-salud-imc-franja`;
  - `oraculo/README.md` l. 37–38 (`S2-metas` y `S2-proyeccion-interruptor`, RESUELTOS; l. 38 dice «PROVISIONAL»).
- Contexto:
  - La nota 4 de `REVISION-criterios-r3` contaba las de `rec_*.py` dentro de las «10 líneas». I-5 decía que quedaban para WF-3.
  - `p2c-remanentes` las declara Pendiente, con dueño (S4 o WF-3).
  - Las de `rec_*.py` están en la receta 1 de la lista para WF-3 (`correccion-contratos-r2.md` l. 43).
- Por qué es bajo: no cambian ninguna cifra ni ninguna salida, `compara.py` da 0, y el pendiente está escrito.
- Lo que no cuadra con su justificación («hay que regenerar s4»): solo la necesitan l. 235, 395 y 644 y las 3 fixtures. Los 4 comentarios (`rec_plan.py` l. 3, 151 y 213, y `rec_genera_fixtures.py` l. 225) y `oraculo/README.md` se pueden pasar a FIRME sin regenerar nada.
- Corrección sugerida: el dueño de S4 o WF-3 pasa los 4 comentarios y `oraculo/README.md` l. 37–38 a FIRME con su Q2. Para las descripciones, regenera s4 y corre `make verify-all`.

**B-2 (baja) · Frase de alcance inexacta en el reporte.**
- `p2c-remanentes.md` («No toqué … criterios, plan ni README») y el mismo reporte («en K46 de `decisiones.md` (l. 120, una sola línea)»).
- La edición de K46 era pedida (INTEGRACION §11, punto 3; M-1 de contratos r3) y es una sola línea. Pero `docs/plan/decisiones.md` está en `docs/plan/`.
- Corrección sugerida: decir «no toqué `plan.md` ni el §K de `decisiones.md`; solo l. 120 (K46)».

## Estado de lo que se pedía

| Origen | Estado | Evidencia |
|---|---|---|
| M-1 de `REVISION-motor-S5-r3` | **CERRADO** | b1 a b7. Bloqueo al día: pregunta 3 con la tabla de los 6 textos y las opciones A1 y A2 (l. 76–88), «Cómo seguí» y `- **Estado:** ABIERTO` (l. 121). Ya no dice «aquí no hay divergencia» (`grep`: 0) |
| B-2, B-5 y B-6 del motor S5 | CERRADOS | b8 |
| B-3 del motor S5 (§15 de S5-interfaz) | SIGUE, fuera de su alcance | Es de `k53-aplicar` o de contratos. Lo declara |
| M-1 de `REVISION-contratos-r3` | **CERRADO** | a1 a a4 |
| I-1 (paquete) | **CERRADO** | c1 a c3. U4-b l. 35 y 393 y el unido l. 490 y 848 llevan «(GET `…`) (PUT `…`)»; no cambia nada más |
| I-3 (DESIGN l. 933) | **CERRADO** | d1 |
| I-4 (línea de estado de C1) | **CERRADO** | a4 |
| I-5 (marcas viejas) | **CERRADO en los sitios que lista**; queda B-1 | d2. `gen2_schemas.py` l. 364 y 374 ya no tienen la marca. `gen6_c1.py` l. 46 es de C1 (declarado) |
| Año 0000 (hallazgo propio) | PROVISIONAL bien hecho | No hay decisión inventada. Se toma lo que el oráculo ya hacía, con cita a la pregunta 3 (ABIERTA), opciones A1 y A2, y fixtures con `provisional` que cita el bloqueo. `ParseDate` y `ca.Year() >= 1` con el mismo PROVISIONAL |

## Notas para el integrador (no son hallazgos de `p2c-remanentes`)
1. **Lectura de fechas de PROFECO.** `profeco/parse.go` acepta más formatos que `_FECHAS` de `pre_profeco.py`, que además usa `\d`. Ya está como PROVISIONAL en `S5-motor-oraculo-reglas`, pregunta 4. Cuando se conteste, conviene el mismo cuidado que en M-1: dígitos ASCII y año 0000.
2. **`S5-agente-externo`.** `p2c-remanentes` lo registró con el estado que tenía en disco. Vuelvan a medir el INDICE cuando cierre `k53-aplicar`, como pide su reporte.
3. **Referencias por renglón.** La lista para WF-3 de `correccion-contratos-r2` cita `pre_genera_fixtures.py` l. 596 (`s5-profeco-muestra`). Hoy está en l. 637. Búsquese por nombre.

## Separación
- **HECHO:** todo lo de las tablas, medido por mí con los comandos citados. Scripts en `/tmp/rev-rem-r1/`. Sondas y salidas en `datos/tmp/rev-rem-r1/`, que git ignora.
- **PREGUNTA:** ninguna para Axel.
- **PROPUESTA:** B-1 y B-2 (técnicas). Ninguna de negocio.
