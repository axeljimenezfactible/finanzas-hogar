# P2d · p2c-remanentes · reporte (12-oct-2026, madrugada)

Agente `p2c-remanentes` (Opus), flujo «P2d · K53», en paralelo con `k53-aplicar` (alcances disjuntos). Base: HEAD `476f4a9` «[Puerta 2]». Sin `git add` ni commit. No toqué `docs/design/sistema/`, `docs/design/maquetas/`, `docs/estado/wfp3/`, S5-interfaz, C1/c1gen, criterios, plan ni README.

## Hecho

**M-1 de `REVISION-motor-S5-r3` (`created_at`): corregido en los dos lados.**
- Motor, `backend/internal/precios/eval.go`. Antes de `time.Parse` exige la gramática de RFC 3339 §5.6 en ASCII (`rfc3339Text`, `^[0-9]{4}-…(Z|[+-][0-9]{2}:[0-9]{2})$`). Los rangos siguen en `time.Parse` y `rfc3339Offset`. Rechaza el año 0000 (PROVISIONAL).
- Oráculo, `oraculo/pre_comun.py`. `_RFC` y `_ISO` usan `[0-9]` en lugar de `\d` y se aplican con `fullmatch` en lugar de `match` con `$`.
- Siete fixtures nuevas, generadas por `pre_genera_fixtures.py`:
  - `s5-ingesta-hora-un-digito`
  - `-hora-un-digito-zona`
  - `-coma-decimal`
  - `-salto-linea`
  - `-digitos-no-ascii`
  - `-anio-cero`
  - `-fecha-anio-cero`
- Todas esperan `invalid_date`. Las dos del año 0000 llevan `provisional` con cita de `S5-precio-vigente-captura`, pregunta 3. `oraculo/salidas/s5.json` regenerada.
- Prueba Go nueva, `TestEvalCreatedAtGramatica` (8 textos), y un renglón nuevo en `TestEvalErrores` (`observed_date` `0000-01-01`).
- **Hallazgo propio, mismo tipo que M-1:**
  - `time.Parse(time.DateOnly)` acepta `observed_date` `0000-01-01` y `date.fromisoformat` no. Lo medí con un programa aparte.
  - `ParseDate` (`types.go`) ahora lo rechaza, PROVISIONAL, con la misma pregunta 3.
- El año 0000 **no lo decidí**. Lo sumé a la pregunta 3 del bloqueo con las opciones A1 y A2. Los dos lados lo rechazan porque el oráculo ya lo hacía, y así coinciden (§15.3).

**Bajos del motor S5:**
- B-2: marcas PROVISIONAL en `compare.go` (total exacto, `S5-carritos-granel`, opción A), `profeco/filter.go` (`S5-profeco-filtro`, preguntas 1 a 3) y `profeco/links.go` (pregunta 2).
- B-5: la cifra medida por `revisor-motor-S5-r2` va en la pregunta 5 de `S5-motor-oraculo-reglas.md`.
- B-6: `types.go` dice «S5-interfaz.md 1.5».
- B-3 no lo hice: es la tabla de §15 de S5-interfaz, que está en el alcance de `k53-aplicar`.

**`S5-precio-vigente-captura.md` al día:**
- Corregí el punto 3 de «Contexto» y la pregunta 3, que decía «aquí no hay divergencia».
- Agregué la tabla de los 6 textos, con su estado antes y después, y las fixtures nuevas.
- Agregué «Cómo seguí» con fecha y la línea `- **Estado:** ABIERTO`.

**I-4.** Agregué la línea `- **Estado:** ABIERTO` (la misma forma que los otros bloqueos de P2c) a `C1-hora-reportes-job.md` y `C1-reparto-alcance.md`.

**Registro (solicitud 2 de `correccion-contratos-r2`).**
- `aplica-registro.py` en ensayo: exit 1, con 2 anclas que ya no existen («| Flujo P2c … | 4 | 0 | 4 |» y «Total 75/63/12»). El integrador ya había dado de alta los 13 bloqueos de P2c.
- Apliqué con un script propio (`/tmp/p2c-rem/parche.py`, con los mismos controles de conteo: ensayo y después `--aplicar`) solo lo que faltaba:
  - la regla de las herramientas: 67 resueltos y 17 abiertos, después 18 con P2d;
  - «conteo de 75»;
  - la remisión de los 4 a los otros 9;
  - la nota del año 0000 en la fila de `S5-precio-vigente-captura`;
  - en K46 de `decisiones.md` (l. 120, una sola línea), el conteo al cierre de P2c (84, 63, 21) con los 13 nombres;
  - las 3 remisiones de `S3-formula-tarjeta`, `S4-confirmacion-empates` y `U4-iva-credito`.
- La fila de `S1-reparto-edicion-gasto` ya traía sus 6 celdas.

**INDICE contra el disco, con P2d.** `k53-aplicar` abrió `S5-agente-externo` y cerró `S5-terminos` (K53). Para que el conteo coincida con el disco:
- agregué la fila «Flujo P2d» (1, 0, 1);
- en WF-1 puse 26, 21 y 5, y en Total 85, 64 y 21;
- agregué el conteo del 12-oct, la fila Q2-64 y la sección «Bloqueos nuevos de P2d». Copié el archivo; no decidí nada.

**I-1.**
- En `U4-ajustes-b.md` (l. 35 y 393) y en el unido (l. 490 y 848), la ruta va ahora con la forma «(GET `…`) (PUT `…`)», que el filtro quita.
- `paquete_entrada.sh`: exit 0. Se regeneraron DESIGN.pdf, DESIGN.docx, U3-dashboard.pdf (U3 de P2c no se había empaquetado), U4-ajustes.pdf y MANIFIESTO.

**I-3.** `DESIGN.md` l. 933 pasa de PROVISIONAL a FIRME (Q2-11 = B, K47) y no cambia ningún valor. Solo la parte del `schedule` cita el bloqueo abierto `C1-hora-reportes-job`, como en U4-b l. 393.

**I-5.** Marcas que citaban bloqueos RESUELTOS, cambiadas a FIRME con su Q2:
- `productos-formato.md` l. 67 (Q2-39 a Q2-42);
- `S2-interfaz.md` l. 354 (Q2-14 = A; la fixture ya trae `provisional` null);
- `ocr/types.go` l. 13 (Q2-76, K50);
- `pre_genera_fixtures.py` l. 254 (Q2-61);
- `backend/internal/recetas`: `plan.go` l. 66, `shopping.go` l. 113, 243 y 259, y `health.go` l. 31 y 32 (Q2-44, 45, 46, 39 y 40).
- `gen2_schemas.py` l. 364 y 374 ya no traen la marca: su cita está en l. 461 con «(resuelto)».

## Verificación (comandos y salidas)
| Comando | Salida |
|---|---|
| `gofmt -l internal/precios internal/recetas cmd/motores`; `go vet …` | vacío; limpio |
| `go test -count=1 -coverprofile=… ./internal/precios/... ./internal/recetas/...` | precios 96.0 %, asistida 100.0 %, ocr 99.3 %, profeco 95.3 %, recetas 93.2 %; `TestEvalCreatedAtGramatica` PASS |
| `python3 -I oraculo/pre_genera_fixtures.py` y `sha256sum -c` (hashes tomados antes) | «43 fixtures». Las 36 viejas y los 7 `_datos` siguen iguales; solo cambió `salidas/s5.json` |
| `pre_run_fixtures.py --out oraculo/salidas/s5.json` | «43 fixtures, 43 coinciden, 1366 hojas de esperado, 0 diferencias» |
| `motores s5 dump` y `compara.py s5` | 963/963 en los tres cruces; provisionales 317/317; **diferencias: 0** |
| Fixtures nuevas contra el código de HEAD (`git archive` en `/tmp/p2c-rem/head`) | Oráculo viejo: 2 diferencias (`salto-linea`, `digitos-no-ascii`). Motor viejo: `__falta__` en `anio-cero`, `coma-decimal`, `fecha-anio-cero`, `hora-un-digito` y `-zona`. Cada fixture detecta al menos un lado viejo |
| `bash docs/design/tools/paquete_entrada.sh` | exit 0; «pdftotext \| grep -c '/api/v1' = 0 en los 4 PDF»; «sha256sum -c: 0 diferencias»; 22 archivos |
| `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` (final, después de todo) | **EXIT=0**; `RESUMEN TODO OK`; `TOTAL comprobadas: 178 \| con diferencias: 0`; `s5 · 43 fixtures … salidas/s5.json idéntica`; `sha256sum -c: 21 archivos sin diferencias (más MANIFIESTO.md)`; `Verificación ampliada OK` |
| Conteo de estados en `docs/estado/bloqueos` (última línea de estado de cada archivo) | 85 archivos: **64 RESUELTO, 21 ABIERTO**, 0 sin estado. El Total del INDICE da 85, 64 y 21, y la suma de las filas coincide |
| `diff` de `decisiones.md` antes y después | 1 línea (K46). No pisé el §K de `k53-aplicar` |

## Supuestos
1. **Año 0000 rechazado en los dos lados (PROVISIONAL).** Elegí lo que ya hacía el oráculo para cumplir §15.3, sin decidir el perfil. Cita la pregunta 3 de `S5-precio-vigente-captura`, que es técnica, para contratos.
2. **Forma de la línea de estado.** En I-4 usé la forma `- **Estado:** ABIERTO` de los otros bloqueos de P2c, no «Estado: ABIERTO» al pie de la letra.
3. **`S5-agente-externo` en el INDICE.** Es de P2d, no de P2c. Lo registré para que «INDICE = disco» se cumpla, con el estado que tenía en disco al medir. Si `k53-aplicar` lo cambia después, manda la carpeta.
4. **`types.go` dice 1.5 aunque `k53-aplicar` lleve S5-interfaz a 1.6.** Su bloqueo dice que el motor no cambia hasta WF-3, así que el código implementa 1.5.

## Bloqueos
- Ninguno nuevo. Amplié la pregunta 3 de `S5-precio-vigente-captura` (año 0000 en `created_at` y en `observed_date`).

## Pendientes (fuera de mi alcance)
- B-3 del motor S5: §15 de `S5-interfaz.md` no lista 14 fixtures, las 7 de r2 y las 7 de hoy. Es de contratos o de `k53-aplicar`.
- `backend/internal/api/c1gen/gen6_c1.py` l. 46 (comentario) dice «bloqueo abierto S1-parametros-avisos». Es de C1/c1gen.
- Siguen marcas de `S4-reglas-lista` y `S4-asignacion-mercado` (RESUELTOS) en `oraculo/rec_plan.py` l. 3, 151 y 213 y en `rec_genera_fixtures.py` l. 225, 235, 395 y 644. Estas tres últimas son descripciones de fixtures s4, así que hay que regenerar s4. Es de S4 o de WF-3.
- I-2, I-6 e I-7 (C1 y orquestador), y los bajos de contratos r3 (por ejemplo, `S2-ventanas.md` l. 85).
- Que el integrador vuelva a medir el INDICE cuando cierre `k53-aplicar`.

## Propuestas para Axel
- Ninguna de negocio. Una PROPUESTA técnica para contratos: A1 en la pregunta 3, es decir, año 0001 o posterior, que es lo que ya hacen los dos lados. La API sella `created_at` al capturar y nunca genera el año 0.

## Consumo aproximado
- Unas 36 acciones de herramienta.
- Unos 130 KB leídos, por partes con `sed`, `grep` y `awk`. No leí `openapi.yaml`, `DESIGN.md` ni INDICE completos.
- Unos 60 KB escritos: el parche, unos 11 KB netos; el INDICE de P2d, 2 KB; las 7 fixtures, 16 KB; `salidas/s5.json`, unos 3 KB; los scripts en `/tmp/p2c-rem`, unos 20 KB; y este reporte, unos 8 KB.
