# Finanzas Hogar

Sistema web de control financiero del hogar de Axel y Ximena. Corre en la máquina de casa, solo en la red local, y se usa desde el navegador del escritorio y como aplicación instalada (PWA) en sus teléfonos Android.

**Estado al 11-oct-2026:** Etapa B. WF-0 «Cimientos» cerrado y **puerta 1 aplicada** (flujo P1): Axel aprobó la puerta y respondió los 31 bloqueos y 3 preguntas de proceso (`docs/plan/decisiones.md` §K), y sus respuestas ya están en el **esquema 1.2**, **C1 v1.1** (con extractos por módulo), **U1 1.1**, los **criterios 1.1** y la **colorimetría 1.1**. Los revisores aprobaron S1, U1, colorimetría y criterios; **C1 v1.1 no quedó aprobado** en la ronda 3 de 3 (2 hallazgos altos y 2 medios, en [`docs/estado/p1/REVISION-C1-r3.md`](docs/estado/p1/REVISION-C1-r3.md)) y pasa a decisión del orquestador antes de WF-1. Quedan 8 bloqueos nuevos de P1 abiertos (ninguno cambia cifras del oráculo). Sigue **WF-1 «Motores y sistema de diseño»**. Todavía no hay motores de cálculo, interfaz ni servicio en operación. Resumen del flujo: [`docs/estado/p1/INTEGRACION.md`](docs/estado/p1/INTEGRACION.md).

## Qué es

- **Para quién.** Axel (titular) y Ximena (co-titular, mismos permisos), cada quien con su usuario. Karen, madre de Axel, participa sin cuenta en la despensa común y en sus compras con el crédito Kueski; recibe un estado de cuenta en PDF que Axel le reenvía.
- **Qué hace.** Presupuesto por mes calendario con quincenas (15 y 30), gastos con reparto configurable (50/50 por defecto), ahorro y metas con fecha, créditos (línea revolvente tipo Kueski, tarjeta, préstamo a plazo, crédito de tienda) con un comparador de escenarios, plan de comidas con recetas, calorías e IMC, inventario y lista de compras, precios por tienda (ticket, OCR local, PROFECO), reportes por Telegram y PDF mensual en LaTeX.
- **Dónde corre.** Máquina de casa (Arch), servicio systemd, base SQLite de un archivo, respaldos locales y HTTPS con certificado propio (dirección propuesta `https://finanzas.local:8443`).
- **Con qué se construye.** Go + `net/http` + SQLite (modernc.org/sqlite) + goose; React 19 + Vite + TypeScript + Tailwind 4 + shadcn/ui. Un oráculo independiente en Python (`oraculo/`) verifica las cifras al centavo.

## Decisiones clave

Las decisiones cerradas de Axel están en [`docs/plan/decisiones.md`](docs/plan/decisiones.md) y ningún agente las reabre. El plan aprobado está en [`docs/plan/plan.md`](docs/plan/plan.md) (70 preguntas respondidas en su §2) y las reglas de trabajo de los agentes en [`docs/plan/protocolo-agentes.md`](docs/plan/protocolo-agentes.md).

| Tema | Decisión | Dónde |
|---|---|---|
| Personas | Axel y Ximena con cuenta (argon2id, 5 intentos en 15 min, sesión larga por dispositivo); Karen externa sin acceso, paga directo a Kueski | decisiones §A |
| Dinero | MXN en centavos exactos, ROUND_HALF_UP, el último pago de una serie ajusta centavos; zona America/Mexico_City; fechas `10-oct-2026` en pantalla e ISO `2026-10-10` en archivos y API | decisiones §B |
| Periodos e ingresos | Mes calendario con quincenas 15 y 30 (o último día); solo los ingresos de Axel; despensa bimestral devengada por mitades | decisiones §B |
| Créditos | 4 tipos; tasa Kueski capturada por compra más tabla editable por plazo (2 q = 17.27 %, 4 q = 25.8685 %); interruptor de lectura A/B (por defecto B); comparador de escenarios con 4 objetivos y 3 restricciones duras, sin solver externo | decisiones §C |
| Comida y salud | Recetario inicial más recetas propias; 4 tiempos; perfil, IMC y meta de calorías; inventario híbrido; mercados y cadenas | decisiones §D |
| Precios | Ticket manual, OCR con confirmación, PROFECO; lectura asistida de sitios desactivada hasta tener veredicto de términos; un precio es viejo a los 30 días | decisiones §E |
| Reportes | Resumen diario y semanal por Telegram, PDF mensual el día 1, estado de cuenta de Karen; 4 alertas | decisiones §F |
| Interfaz | Fiel a Notion, móvil primero (registrar gasto en 3 toques o menos), claro y oscuro, WCAG 2.2 AA, matiz base **C · Índigo 280°** | decisiones §G e §I |
| Infraestructura | Solo red local, Go + SQLite + systemd, respaldos locales, IA opcional desactivada por defecto | decisiones §H |
| Proceso | Etapa B por flujos (WF-0, WF-1, WF-3, WF-4) con 4 puertas de aprobación; un commit por flujo | decisiones §J; plan §3.3 |

## Árbol del repositorio

```
finanzas-hogar/
├── Makefile                         verificación (make verify-wf0, verify-c1, verify-s1, verify-u1, verify-criterios, verify-all) y limpieza
├── README.md
├── backend/                         Go, módulo finanzas-hogar (go.mod declara go 1.27; K28)
│   ├── cmd/finanzas/main.go         servidor mínimo: GET /health y GET /openapi.yaml; las demás rutas responden 501
│   ├── db/migrations/
│   │   ├── 00001_esquema.sql        migración goose: esquema 1.1 de S1 (WF-0)
│   │   └── 00002_puerta1.sql        migración goose de la puerta 1: esquema 1.2, 55 tablas STRICT en total
│   ├── db/tools/                    s1_verifica.py (verificación de S1, make verify-s1) y gen_dict.py
│   └── internal/api/
│       ├── routes.go                201 rutas generadas desde el contrato C1 (no se edita a mano)
│       ├── generate.go              go generate → c1gen/build.py
│       └── c1gen/                   generador y verificador de C1 (openapi.yaml, routes.go, C1-datos.md y extractos)
├── datos/semilla/                   data.json (compra del 10-oct-2026) y profeco/ (diccionario y metadatos)
├── docs/
│   ├── plan/                        plan.md, decisiones.md, protocolo-agentes.md y fuentes/ (Kueski)
│   ├── contratos/                   C1-datos.md, openapi.yaml, diccionario-datos.md, reglas-negocio.md, C2-colores.json
│   │   └── extractos/               openapi.yaml por módulo (17 YAML e INDICE.md); los genera build.py, no se editan
│   ├── criterios/                   INDICE.md y 18 archivos de criterios de aceptación, uno por fase
│   │   └── tools/                   verif_final.sh y sus scripts (conteos, trazas, citas, cifras canon; make verify-criterios)
│   ├── design/                      U1-arquitectura.md, colorimetria.md, contraste.csv, paleta.html y .png
│   │   └── tools/                   colorimetria.py (regenera la colorimetría) y verificar_u1.py (make verify-u1)
│   └── estado/
│       ├── wf0/                     reportes de agentes, revisiones REVISION-r1 a r3, correcciones e INTEGRACION.md
│       ├── p1/                      flujo de la puerta 1: reportes, revisiones por artefacto, correcciones e INTEGRACION.md
│       └── bloqueos/                preguntas para Axel (una por archivo, con su respuesta al final) e INDICE.md
└── oraculo/                         model.py y verify3.py con sus JSON: cifras canon de regresión (no se editan)
```

`frontend/`, `ops/`, `datos/semilla/nutricion/` y `docs/design/{entrada,salida,capturas-u5}/` existen en disco pero están vacías, así que git no las guarda; se llenan en WF-1, WF-3 y WF-4. `datos/tmp/` es solo para bases temporales de verificación y nunca entra al repositorio.

## Versiones vigentes de los artefactos (11-oct-2026)

| Artefacto | Versión | Ruta | Estado en P1 |
|---|---|---|---|
| Modelo de datos (S1) | esquema 1.2 | `backend/db/migrations/00001_esquema.sql` + `00002_puerta1.sql`, `docs/contratos/diccionario-datos.md` 1.2, `docs/contratos/reglas-negocio.md` | aprobado (ronda 3): 55 tablas STRICT, 152 índices, 71 disparadores, 28 reglas de negocio |
| Contrato C1 | v1.1 (`info.version` 1.1.0) | `docs/contratos/openapi.yaml`, `docs/contratos/C1-datos.md`, `docs/contratos/extractos/`, `backend/internal/api/` | **no aprobado** (ronda 3 de 3: 2 altos y 2 medios); 201 operaciones en 162 rutas, matriz de 227 filas |
| Arquitectura de información (U1) | 1.1 | `docs/design/U1-arquitectura.md` | aprobado (ronda 3): 27 pantallas, 14 secciones de Ajustes, 203 necesidades |
| Colorimetría (U2a) | 1.1 | `docs/design/colorimetria.md`, `docs/contratos/C2-colores.json`, `docs/design/contraste.csv`, `docs/design/paleta.html` | aprobado (ronda 2): 190 tokens, 202 pares de contraste sin ninguno bajo umbral; gasto en 45° (K26) |
| Criterios de aceptación | 1.1 | `docs/criterios/` | aprobado (ronda 3): 257 criterios (197 alta y 60 media) en 18 archivos |
| Bloqueos | puerta 1 respondida | `docs/estado/bloqueos/INDICE.md` | 31 preguntas respondidas (§K); 8 bloqueos nuevos de P1 abiertos |

## Qué entregó WF-0

Tabla histórica de WF-0 (cifras de entonces; las vigentes están en la tabla anterior).


| Entregable | Ruta | Resumen |
|---|---|---|
| Modelo de datos (S1) | `backend/db/migrations/00001_esquema.sql`, `docs/contratos/diccionario-datos.md`, `docs/contratos/reglas-negocio.md` | 54 tablas STRICT, 147 índices, 68 disparadores; dinero en centavos, fechas ISO con CHECK, tasas como fracción exacta; 22 reglas de negocio con ejemplo del oráculo |
| Contrato C1 v1.0 | `docs/contratos/openapi.yaml`, `docs/contratos/C1-datos.md`, `backend/internal/api/` | OpenAPI 3.1.1 con 158 rutas y 197 operaciones; matriz pantalla → necesidad → operación de 220 filas sin huecos; esqueleto Go que compila |
| Arquitectura de información (U1) | `docs/design/U1-arquitectura.md` | barra lateral de 9 entradas, 27 pantallas y 14 secciones de Ajustes, 199 necesidades de datos, 8 flujos (registrar gasto en 3 toques) |
| Colorimetría del matiz C (U2a) | `docs/design/colorimetria.md`, `docs/contratos/C2-colores.json`, `docs/design/contraste.csv`, `docs/design/tools/colorimetria.py` | 190 tokens de color, 202 pares de contraste sin ninguno bajo umbral; se regenera con un script determinista |
| Criterios de aceptación | `docs/criterios/` | 248 criterios medibles (189 de severidad alta) para las 18 fases, con 85 cifras canon de regresión |
| Bloqueos | `docs/estado/bloqueos/INDICE.md` | 32 archivos (31 abiertos) y 2 preguntas de la revisión, agrupados en 31 preguntas para Axel |
| Integración | `docs/estado/wf0/INTEGRACION.md` | qué se hizo por fase, verificación, supuestos, pendientes para WF-1 y rondas de revisión |

## Cómo correr la verificación

Desde la raíz del repositorio:

```sh
make verify-wf0      # SQL en base temporal, go build y go vet, colorimetría y validador OpenAPI
make verify-c1       # build.py --check (C1 y extractos idénticos al generador) y verify_c1.py contra la base migrada
make verify-s1       # backend/db/tools/s1_verifica.py: esquema 1.2, oráculo, negativas, Down y goose
make verify-u1       # docs/design/tools/verificar_u1.py: U1 1.1
make verify-criterios  # docs/criterios/tools/verif_final.sh: los 18 archivos de criterios
make verify-all      # todo lo anterior
make clean           # borra los temporales si una verificación falló a medias
```

Requisitos: GNU Make y bash, `sqlite3` 3.37 o posterior, Go 1.27 o posterior, Python 3.11 o posterior, y `uv`/`uvx` con `openapi-spec-validator` 0.9.0 en su caché (para `verify-c1` también `pyyaml` y `jsonschema`; para `verify-criterios`, `pyyaml`). Si la caché está vacía, la primera corrida necesita internet: `make verify-wf0 UVX_FLAGS= UV_FLAGS=`. Si `goose` v3.28.0 está en el PATH (o con `GOOSE=/ruta/a/goose`), también se prueba su ciclo up, down y up.

Qué comprueba `make verify-wf0`, sin escribir en ningún archivo del repositorio:

1. **SQL.** Aplica la sección Up de cada migración con `sqlite3 -bail` en `datos/tmp/verify-wf0/`. Exige que todas las tablas sean STRICT, `PRAGMA foreign_key_check` sin filas e `integrity_check` = ok. Luego aplica Down en orden inverso (la base debe quedar vacía) y Up otra vez.
2. **Go.** `go build ./...`, `go vet ./...` y `gofmt -l .` en `backend/`.
3. **Colorimetría.** Corre `docs/design/tools/colorimetria.py --sin-png` en un árbol espejo dentro de `datos/tmp`, exige «VERIFICACIÓN OK» y que `colorimetria.md`, `contraste.csv`, `C2-colores.json` y `paleta.html` salgan idénticos byte a byte a los del repositorio, con 0 filas `cumple=false`.
4. **OpenAPI.** `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml`.

Salida esperada al 11-oct-2026 (esquema 1.2): 00001 y 00002 aplicadas en ese orden, 55 tablas (55 STRICT), 152 índices, 71 disparadores, 0 filas en `foreign_key_check`; Down en orden inverso deja 0 objetos; goose v3.28.0 up, down-to 0 y up sin errores; Go sin errores; 202 pares de contraste con 0 `cumple=false`; `openapi.yaml: OK`. `make verify-all` agrega `build.py --check` (los 21 archivos de C1, extractos incluidos, idénticos al generador), `verify_c1.py --db` (24 comprobaciones, «TODO OK»), `s1_verifica.py` («RESUMEN TODO OK»; regresión Kueski: suma de pagos de la despensa 1 = 389500 centavos y parte de Axel = 239528 centavos, igual que `oraculo/model.json`), `verificar_u1.py` (53 comprobaciones OK, «TODO OK») y `verif_final.sh` (257 criterios en 18 archivos, 178 cifras canon sin diferencias, 12 de 12 cifras del calendario contra el oráculo). Corrida completa: unos 10 segundos.

Comandos equivalentes sin make:

```sh
mkdir -p datos/tmp
for f in backend/db/migrations/0*.sql; do awk '/^-- \+goose Up/{s=1;next} /^-- \+goose Down/{s=0} s' "$f" | sqlite3 -bail datos/tmp/wf0.db; done
sqlite3 datos/tmp/wf0.db "SELECT count(*) FROM pragma_foreign_key_check; PRAGMA integrity_check;"
(cd backend && go build ./... && go vet ./...)
uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml
rm -rf datos/tmp
```

Para la colorimetría sin make, copia el script a un árbol espejo (`<espejo>/docs/design/tools/`) y córrelo ahí con `--sin-png`: el script escribe en la raíz de su propio árbol, y así no reescribe el repositorio.

## Cómo sigue

1. **Puerta 1 (aplicada el 11-oct-2026).** Axel aprobó C1 v1.0, U1 y los criterios y respondió todos los bloqueos (`docs/plan/decisiones.md` §K). El flujo P1 aplicó las respuestas (esquema 1.2 con la migración 00002, C1 v1.1 regenerado con `backend/internal/api/c1gen/build.py`, U1 1.1, criterios 1.1 y colorimetría 1.1). Antes de arrancar WF-1, el orquestador decide qué hacer con los 4 hallazgos de C1 que siguen abiertos y Axel puede responder los 8 bloqueos nuevos (ver [`docs/estado/p1/INTEGRACION.md`](docs/estado/p1/INTEGRACION.md)).
2. **WF-1 «Motores y sistema de diseño»** (plan §3.3 y §5). Pista 1: S2 motor financiero con su oráculo y fixtures; S3 motor de créditos, comparador de escenarios y regresión Kueski; S4 recetas, comidas y salud, recetario de 60 recetas y nutrición; S5 precios y tiendas, PROFECO y OCR. Pista 2: U2 `DESIGN.md` y tokens C2 a partir de `C2-colores.json`; U3 especificación del dashboard; U4 especificación de Ajustes. Cada motor se verifica contra `oraculo/` al centavo y contra los criterios de `docs/criterios/`. Revisores independientes y un integrador cierran el flujo.
3. **Puerta 2:** `DESIGN.md`, tokens C2 y motores verificados. Después vienen la Pista 3 (Claude Design, puerta 3), WF-3 (backend, frontend y verificación visual; puerta 4) y WF-4 (operación, semilla, seguridad e informe final).

Lo que WF-1 debe tomar de WF-0 (insumos, pendientes y solicitudes abiertas) está en [`docs/estado/wf0/INTEGRACION.md`](docs/estado/wf0/INTEGRACION.md) §5 y §7, y lo de la puerta 1 en [`docs/estado/p1/INTEGRACION.md`](docs/estado/p1/INTEGRACION.md).

## Reglas de trabajo

- Todo agente lee [`docs/plan/protocolo-agentes.md`](docs/plan/protocolo-agentes.md): no asume (lo que falta va como bloqueo), escribe solo en su alcance, no hace commits y deja su reporte en `docs/estado/<flujo>/<agente>.md`.
- Toda cifra que exista en `oraculo/model.json` se reproduce exactamente; si no coincide, es un hallazgo, no un ajuste.
- Documentos e interfaz en español de México; identificadores de código en inglés.
- El integrador de cada flujo hace un solo commit. Git es local, sin remoto.
