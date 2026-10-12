# Finanzas Hogar

Sistema web de control financiero del hogar de Axel y Ximena. Corre en la máquina de casa, solo en la red local, y se usa desde el navegador del escritorio y como aplicación instalada (PWA) en sus teléfonos Android.

**Estado al 12-oct-2026, WF-3c (Inicio como tablero de control, K54):** Inicio › Resumen ya es un tablero de control. Arriba va el gasto real del mes, y junto a él el ingreso, los créditos y los pagos próximos. El **barómetro de capacidad de endeudamiento** divide los pagos de crédito del mes entre el ingreso neto: verde hasta 30 %, amarillo hasta 40 % y rojo arriba, con umbrales que se pueden cambiar. Debajo hay 6 gráficas para anticiparse y «Más del hogar». Para eso se hicieron U3 2.0, DESIGN 1.1 (sin tokens nuevos), el esquema 1.5 (migración 00005) y C1 v1.5.0 (221 operaciones, 176 rutas y 267 esquemas). El servidor de desarrollo ya habla **HTTP/2** y las 164 capturas se regeneraron con el tablero. `make verify-all` pasa. `make verify-wf3` falla solo en U5-04, con 2 pares del dibujo del barómetro que son un falso positivo comprobado. **Faltan aprobaciones:** la especificación y el backend quedaron con 1 medio cada uno tras 3 rondas, y C3 v1.3 volvió sin firma con 4 altos. El frontend sí quedó aprobado. Detalle y propuesta de una ronda acotada: [`docs/estado/wf3/INTEGRACION-wf3c.md`](docs/estado/wf3/INTEGRACION-wf3c.md) y `docs/estado/bloqueos/WF3c-ronda3-escalamiento.md`.

**Estado al 12-oct-2026, WF-3b (verificación visual y C3):** las 41 pantallas de U1 ya tienen sus **164 capturas reales**: la aplicación en marcha con la semilla, en Android 412×915 táctil y en escritorio 1440×900, en claro y en oscuro. Cada captura trae su axe-core (0 violaciones serias o críticas), su contraste medido en el DOM (20,919 pares y 0 bajo umbral) y sus hosts (solo `finanzas.local`). Todo está versionado en `docs/design/capturas-u5/`. Los flujos de U1 se probaron de extremo a extremo contra el servidor real, en copias aisladas de la base. También se probó la cola sin conexión (U5-13: sin duplicados y con reenvío idempotente). Las 111 cifras de pantalla cuadran con el oráculo (C3-09). **U5 no quedó aprobado tras 3 rondas:** quedan 2 hallazgos medios. Uno es la barra lateral sin iconos de sección. El otro, la marca «capturado» en C3, ya está resuelto en C3 v1.2. **C3 v1.2 sigue sin firma tras 3 rondas:** pasan 7 de 8 flujos y F-3 «Comparar canasta» por la ruta A falla con la semilla, porque la lista del plan no tiene precios de ingredientes. Los dos casos se escalaron a Axel con opciones; la recomendada es una corrección y una verificación acotadas. Dentro del flujo se cerraron los 2 medios que el frontend arrastraba de WF-3a bis. La semilla aplicó las decisiones PROVISIONALES del orquestador: gastos proyectados y precios de Walmart. `make verify-wf3` (ahora con `verify-u5`, la evidencia U5 sin servidor) y `make verify-all` pasan. Hay 8 bloqueos nuevos (137 archivos en total). Para la puerta 4 quedan pendientes S6-11 en los Android (Axel), S7-07 con envío real (orquestador) y los valores K31. Detalle: [`docs/estado/wf3/INTEGRACION-wf3b.md`](docs/estado/wf3/INTEGRACION-wf3b.md).

**Estado al 12-oct-2026, WF-3a bis (ronda acotada de corrección del frontend):** el servidor Go atiende las 216 operaciones de C1 v1.4, y las 216 pasan su prueba de contrato. Corre por HTTPS (mkcert) solo en la red de casa, en `https://finanzas.local:8443`, con la semilla real más una **semilla de desarrollo** (`finanzas semilla --desarrollo`, con montos de las fixtures del oráculo S2): 5 de las 6 tarjetas de Inicio ya muestran cifras. El backend quedó **aprobado**. El frontend corrigió el hallazgo alto y los medios de la ronda 3: objetivos táctiles, gráficas accesibles, barra inferior, edición en grupo con `Ctrl+/` y las necesidades de U1 que no tenían llamada. Aun así, **no quedó aprobado tras 5 rondas**: le quedan 2 hallazgos medios (fechas de P-10 y errores sin aviso en Ajustes 01–07) y varios bajos, ninguno toca cifras. Playwright con la API simulada da 230 pruebas que pasan, 4 omitidas y 1 que falla en cada corrida completa (`correccion-bis.escritorio.spec.ts:79`). Por esa prueba, `make verify-wf3` queda en rojo; `make verify-all` pasa. C3 v0.95 sigue sin firma: 42 de 42 pantallas implementadas y probadas, 0 capturadas (eso es WF-3b). Hay 4 bloqueos nuevos (129 archivos en total). Detalle: [`docs/estado/wf3/INTEGRACION-wf3a.md`](docs/estado/wf3/INTEGRACION-wf3a.md), §15 «WF-3a bis».

**Estado al 12-oct-2026, P2d (K53):** la aplicación no leerá sitios de tiendas: se retiró la lectura asistida, y un agente externo (la sesión de Claude Code de Axel con Claude in Chrome) meterá precios, productos y recetas por la API (C1 v1.4), con el peso de esos datos PROVISIONAL hasta la puerta 4 (`docs/estado/bloqueos/S5-agente-externo.md`).

**Estado al 11-oct-2026, puerta 2:** Etapa B. Axel contestó la puerta 2 (98 preguntas) en modo «lo recomendado + excepciones», sin excepciones. Las preguntas sin recomendación quedan PROVISIONAL con la opción en uso y pasan a WF-3 o a la puerta 4 (`docs/estado/puerta2/RESPUESTAS.md`). El flujo P2c las aplicó: esquema 1.4 (migración 00004, con los parámetros de avisos y reportes en `app_setting`), C1 v1.3, contratos S2 1.3, S3 1.5, S4 1.5 y S5 1.5, motores S2 y S5 con sus oráculos, U3 1.5, U4 1.6 y criterios 1.3. `make verify-all` pasa completo. Tras 3 rondas quedan **sin aprobar** los bloqueos, S1, los contratos de motores y el motor S5: les quedan hallazgos medios y bajos, y el alto que compartían ya pasa. **Los términos de las tiendas quedaron «no verificado»** por un desafío anti-robot y se escalaron a Axel. Axel contestó en la puerta 2 bis que la aplicación no leerá sitios de tiendas: un agente externo traerá precios, productos y recetas. Es K53 y se registra en P2d. El paquete para Claude Design no se regeneró, porque U4 trae 2 rutas `/api/v1` que el filtro del script no quita. Bloqueos: 84 archivos, 63 resueltos y 21 abiertos. Detalle: [`docs/estado/p2c/INTEGRACION.md`](docs/estado/p2c/INTEGRACION.md). La Pista 3 (K45) sigue en curso.

**Estado al 11-oct-2026, cierre de WF-1:** Etapa B. WF-0, la puerta 1 (P1), la puerta 1b (P1c) y **WF-1 «Motores y sistema de diseño»** están integrados. WF-1 dejó los 4 motores en Go (S2 financiero, S3 créditos y comparador, S4 recetas y salud, S5 precios, PROFECO y OCR), cada uno con su oráculo independiente en Python y sus fixtures, que reproducen `oraculo/model.json` y `datos/semilla/data.json` con 0 diferencias; el recetario de 60 recetas con kcal y fuente; el catálogo de 160 productos; `DESIGN.md` y los tokens C2; U3 (Inicio) y U4 (Ajustes, unido en `docs/design/U4-ajustes.md`); y el paquete para Claude Design en `docs/design/entrada/`, que ya no se sube a ningún sitio: desde K45 (11-oct-2026) la Pista 3 la ejecuta el orquestador con Artifacts y el paquete queda como brief de referencia (antes K35). `make verify-all` pasa completo. **Siguiente: puerta 2.** Axel decide 6 artefactos sin aprobar (motor S2, motor S5, contrato S4/S5, U4 y el paquete), 27 bloqueos abiertos y los valores operativos propuestos: todo en [`docs/estado/wf1/INTEGRACION.md`](docs/estado/wf1/INTEGRACION.md) y en `docs/estado/bloqueos/INDICE.md`.

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
├── Makefile                         verificación (verify-wf0, verify-c1, verify-s1, verify-u1, verify-criterios, verify-wf1, verify-all, verify-wf3 y verify-wf3-final) y limpieza
├── README.md
├── backend/                         Go, módulo finanzas-hogar (go 1.27; K28); solo biblioteca estándar y modernc.org/sqlite
│   ├── cmd/finanzas/                servidor (serve, migrate, semilla, respaldo): 216 operaciones de C1 v1.4, TLS y la PWA de frontend/dist
│   ├── cmd/motores/                 binario común de los motores: motores <s2|s3|s4|s5> dump --fixtures … --out … (correr desde la raíz)
│   ├── db/migrations/               00001 (esquema 1.1), 00002 (1.2) y 00003 (1.3): 55 tablas STRICT
│   ├── db/tools/                    s1_verifica.py (make verify-s1) y gen_dict.py
│   └── internal/
│       ├── api/                     rutas generadas desde C1 v1.2 y c1gen/ (generador y verificador de C1)
│       ├── motor_fin/               S2: devengo, efectivo, presupuesto, flujo, metas, proyección
│       ├── motor_cred/              S3: créditos y series (más comparador/, evaluar/ y canontest/)
│       ├── recetas/                 S4: recetario, escalado, lista, inventario, plan semanal, kcal y salud
│       └── precios/                 S5: comparativa y frescura (más ocr/, profeco/ y asistida/, apagada)
├── datos/semilla/                   data.json, productos.json (160), recetario.json (60 = recetario-a + recetario-b) y su esquema,
│                                    nutricion/ (tabla.csv y procedencia), profeco/ (muestra; raw/ y db/ ignorados), tickets-sinteticos/, tickets/
├── docs/
│   ├── plan/                        plan.md, decisiones.md, protocolo-agentes.md y fuentes/
│   ├── contratos/                   C1 (openapi.yaml, C1-datos.md, extractos/), diccionario, reglas de negocio, C2-colores.json,
│   │                                C2-tokens.json, tokens.css, tailwind.theme.css
│   │   └── motores/                 interfaces S2, S3, S4 y S5 y formato del catálogo de productos
│   ├── criterios/                   INDICE.md y 18 archivos de criterios (tools/: verif_final.sh)
│   ├── design/                      U1, DESIGN.md, U3-dashboard.md, U4-ajustes{,-a,-b}.md, colorimetría y contraste
│   │   ├── entrada/                 brief de referencia de la Pista 3 (LEEME.md y MANIFIESTO.md); no se sube (K45)
│   │   ├── fuentes/                 Inter y Lucide autoalojados
│   │   └── tools/                   colorimetria.py, verificar_u1.py, contraste_componentes.py, tokens/ (Style Dictionary) y paquete_entrada.sh
│   └── estado/
│       ├── wf0/, p1/, p1c/          reportes, revisiones, correcciones e INTEGRACION.md de cada flujo
│       ├── wf1/                     reportes de WF-1, revisiones REVISION-*-rN, correcciones e INTEGRACION.md
│       └── bloqueos/                preguntas para Axel (una por archivo) e INDICE.md
├── frontend/                        WF-3a: Vite 8, React 19, TypeScript, Tailwind 4, Radix, TanStack Table 9, dnd-kit, ECharts 6 y PWA;
│                                    src/ (shell, api, cola, componentes, vistas-bd, pantallas/<área>), pruebas/ (Playwright) y scripts/
├── ops/                             WF-3a: dev-run.sh y dev-stop.sh (servidor de desarrollo), mkcert.md, runbook.md (borrador) y telegram-disparo.md
└── oraculo/                         model.json y verify3.json (canon, no se editan), oráculos fin_/cred_/rec_/pre_*.py,
                                     compara.py, fixtures/s2..s5 y salidas/
```

`docs/design/capturas-u5/` guarda las 164 capturas de U5 con su evidencia (`axe/`, `contraste/`, `evidencia/`, `hosts.json`; WF-3b) y sí se versiona; las salidas de Playwright (`test-results/`) no. `datos/certs/` (certificado de mkcert), `datos/dev/` y `datos/*.db` (servidor y base de desarrollo) nunca entran al repositorio. `datos/tmp/` es solo para temporales de verificación y nunca entra al repositorio; las descargas crudas (`datos/semilla/*/raw/`, ≈ 13 GB) y la base de PROFECO (`profeco/db/`, ≈ 17 GB) tampoco.

## Versiones vigentes de los artefactos (11-oct-2026, puerta 2)

> **Cambios del 12-oct-2026 (WF-3c, K54)**, la tabla todavía no los refleja:
>
> - Esquema 1.5 (`00005_tablero.sql`), diccionario y reglas 1.5.
> - C1 v1.5.0: 221 operaciones, 176 rutas y 267 esquemas.
> - U3 2.0 y DESIGN 1.1.
> - Criterios: U3 1.4, U5 1.3, P3 1.4, C1 1.5 y S1 1.3.
> - C3 v1.3, sin firma.
>
> Ver `docs/estado/wf3/INTEGRACION-wf3c.md` §4 y §5.

| Artefacto | Versión | Ruta | Estado (puerta 2, P2c) |
|---|---|---|---|
| Modelo de datos (S1) | esquema 1.4 | `backend/db/migrations/00001_esquema.sql` a `00004_puerta2.sql`, `docs/contratos/diccionario-datos.md` 1.4, `docs/contratos/reglas-negocio.md` 1.4 | **no aprobado** tras 3 rondas en P2c (el alto, `verify-c1` roto, ya pasa; quedan 7 bajos): 55 tablas STRICT, 152 índices, 79 disparadores, 29 reglas; 10 claves nuevas de `app_setting` (Q2-11 B, K47 y K50) y `credit.vat_included` (Q2-67) |
| Contrato C1 | v1.3 (`info.version` 1.3.0) | `docs/contratos/openapi.yaml`, `docs/contratos/C1-datos.md`, `docs/contratos/extractos/`, `backend/internal/api/` | aprobado (P2c, ronda 3: 6 bajos): 212 operaciones en 168 rutas, 255 esquemas; `getReportSettings` y `updateReportSettings` |
| Arquitectura de información (U1) | 1.2 | `docs/design/U1-arquitectura.md` | aprobado (ronda 3): 27 pantallas, 14 secciones de Ajustes, 211 necesidades (138 L + 73 E) |
| Colorimetría (U2a) | 1.2 | `docs/design/colorimetria.md`, `docs/contratos/C2-colores.json`, `docs/design/contraste.csv`, `docs/design/paleta.html` | aprobado (ronda 2): 190 tokens sin cambio de valor (K41 y K42 quitan 5 marcas), 202 pares de contraste sin ninguno bajo umbral |
| Criterios de aceptación | 1.3 | `docs/criterios/` | aprobado (P2c, ronda 3: 4 bajos); 261 criterios (201 alta y 60 media) en 18 archivos; K31 y K44 FIRMES con clave y valor; S5-07 sigue NO CUMPLE |
| Motor S2 (`motor_fin`) | contrato 1.3 | `backend/internal/motor_fin/`, `docs/contratos/motores/S2-interfaz.md`, `oraculo/fin_*.py` | aprobado (P2c, ronda 2); 0 diferencias en 3185 hojas y 159 de `model.json` |
| Motor S3 (`motor_cred`) | contrato 1.5 (solo marcas) | `backend/internal/motor_cred/`, `S3-interfaz.md`, `oraculo/cred_*.py` | motor sin cambios en P2c; el contrato va en la cadena de contratos, **no aprobada** (ronda 3: 1 medio y 5 bajos; el alto ya pasa); 906 de 906 hojas y 437 de 437 de `model.json` |
| Motor S4 (`recetas`) y recetario | contrato 1.5 | `backend/internal/recetas/`, `S4-interfaz.md`, `datos/semilla/recetario.json` | motor aprobado (WF-1, ronda 2); contrato **no aprobado** (cadena de contratos de P2c, ronda 3); 8451 de 8451 hojas |
| Motor S5 (`precios`) | contrato 1.5 | `backend/internal/precios/`, `S5-interfaz.md`, `oraculo/pre_*.py` | **no aprobado** (P2c, ronda 3: 1 medio y 6 bajos); 958 de 958 hojas; términos de las tiendas «no verificado» (K53 en P2d) |
| U2: DESIGN.md y tokens C2 | 1.0 | `docs/design/DESIGN.md`, `docs/contratos/C2-tokens.json`, `tokens.css` | aprobados tal cual por Axel (Q2-96 A); 4 líneas de ratificación (Q2-81, Q2-82 y Q2-84) |
| U3 Inicio · U4 Ajustes | 1.5 · 1.6 (a 1.6 y b 1.5) | `docs/design/U3-dashboard.md`, `docs/design/U4-ajustes.md` | aprobados (P2c, ronda 3: 5 bajos) |
| Paquete para Claude Design | — | `docs/design/entrada/` | **sin regenerar**: el script se detiene porque U4 1.6 trae 2 rutas `/api/v1` que su filtro no quita; sigue el de WF-1 (U3 1.2 y U4 1.4) y `sha256sum -c` no da diferencias |
| Bloqueos | cierre de P2c | `docs/estado/bloqueos/INDICE.md` | 84 archivos: 63 resueltos y 21 abiertos (6 de WF-1, 13 de P2c y 2 de la Pista 3) |

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

## Qué entregó WF-1

Detalle, verificaciones, rondas extra, interrupción y continuación, hallazgos, bloqueos, propuestas para Axel y consumo por agente: [`docs/estado/wf1/INTEGRACION.md`](docs/estado/wf1/INTEGRACION.md). En corto:

- 4 motores con oráculo independiente: S2 30 fixtures, S3 27, S4 26 y S5 29; `make verify-wf1` corre el oráculo y el motor de cada uno y `oraculo/compara.py` da 0 diferencias. Cobertura ≥ 85 % (K25) en `motor_fin`, `motor_cred`, `recetas` y `precios`.
- Recetario de 60 recetas validado (JSON Schema, generador del oráculo y pruebas Go), catálogo de 160 productos, tabla nutricional con fuente por renglón, muestra de PROFECO y OCR listo para medirse con las 5 fotos de Axel (K24).
- `DESIGN.md`, tokens C2 compilados con Style Dictionary (reproducibles), contraste de componentes sin pares bajo umbral, U3 y U4.
- Paquete para Claude Design en `docs/design/entrada/` con `LEEME.md` paso a paso para Axel.

## Qué entregó WF-3a

Detalle por agente, criterios S6 y S7, cobertura, operación → prueba, pantallas, hallazgos abiertos, bloqueos, propuestas para la puerta 4 y consumo: [`docs/estado/wf3/INTEGRACION-wf3a.md`](docs/estado/wf3/INTEGRACION-wf3a.md). En corto:

- **Servidor Go** (`backend/cmd/finanzas`, `backend/internal/api/<área>`, `auth`, `jobs`, `respaldo`, `adjuntos`, `semilla`, `telegram`, `reportes`). Tiene 216 de 216 operaciones con prueba de contrato contra `openapi.yaml`. La cobertura K25 se cumple en dominio, motor_fin, motor_cred, recetas y precios. Los motores y sus cifras no se tocaron.
- **Frontend** (`frontend/`). Tiene las 41 pantallas de U1, la PWA instalable con la cola sin conexión y el tema claro, oscuro y del sistema; las fuentes y los íconos están autoalojados (0 hosts externos). Las 149 pruebas de Playwright con la API simulada pasan, y la prueba de humo de Inicio contra el servidor real cuadra con el oráculo.
- **C3 v0.9** (`docs/contratos/C3-integracion.md`): borrador sin firma generado por `frontend/scripts/c3-integracion.py`. La firma va en WF-3b.
- **Operación** (`ops/`): servidor de desarrollo con `bash ops/dev-run.sh` y `bash ops/dev-stop.sh`, mkcert y el borrador del runbook para WF-4.

## Qué entregó WF-3b

Criterios U5 y C3 uno por uno, rondas, hallazgos abiertos, bloqueos, pendientes de la puerta 4, valores K31 y consumo: [`docs/estado/wf3/INTEGRACION-wf3b.md`](docs/estado/wf3/INTEGRACION-wf3b.md). En corto:

- **Capturas U5** (`docs/design/capturas-u5/`): 164 PNG, más `axe/`, `contraste/`, `evidencia/` y `hosts.json`, generados por `frontend/playwright.u5.config.ts` (`frontend/pruebas-u5/`) contra el servidor real. Registro de la revisión: `docs/estado/U5-revision.md`, 3 rondas.
- **Flujos de extremo a extremo** (`frontend/pruebas-u5/flujos/`): F-1 a F-8, la cola sin conexión (U5-13) y C3-09. Corren contra una copia aislada de la base de desarrollo (`servidor-prueba.sh`, `127.0.0.1:8444`) y guardan sus capturas y resultados en JSON.
- **C3 v1.2** (`docs/contratos/C3-integracion.md`): 41 pantallas implementadas, probadas y capturadas; 7 de 8 flujos; sin firma.
- **Semilla**: gastos proyectados de la fixture S2 (devengado contra efectivo P1 a P6 = `model.json`) y 36 precios de Walmart, PROVISIONAL (`WF3a-bis-semilla-dev`).

## Qué entregó WF-3c

Inicio como tablero de control (K54). Detalle, rondas, hallazgos abiertos, bloqueos y consumo: [`docs/estado/wf3/INTEGRACION-wf3c.md`](docs/estado/wf3/INTEGRACION-wf3c.md). En corto:

- **Especificación:** K54 en `docs/plan/decisiones.md`; U3 2.0, con K-1 a K-4, el barómetro B-1 y las gráficas G-1 a G-6; DESIGN 1.1, §8.9 «Barómetro»; esquema 1.5 (`backend/db/migrations/00005_tablero.sql`, 5 claves de `app_setting`); C1 v1.5.0, con 5 operaciones nuevas.
- **Servidor:**
  - `backend/internal/api/tablero/`, los pagos por mes en `api/creditos/tablero.go` y K-1 con la comparación contra el mes anterior;
  - HTTP/2 por ALPN (`cmd/finanzas/tls.go`, prueba `tls_h2_test.go`).
- **Pantalla:** `frontend/src/pantallas/inicio/Tablero.tsx`, `frontend/src/componentes/Barometro.tsx` y la barra lateral con los iconos de sección.
- **Capturas:** las 164 se regeneraron, más 16 de las subsecciones de P-01 en `docs/design/capturas-u5/evidencia/inicio-subsecciones/`.
- **Abierto:**
  - la especificación (M-R3-1) y el backend (M-2) no se aprobaron;
  - C3 v1.3 sin firma: 4 pruebas de flujos quedaron viejas con el nuevo Inicio, a `updateDashboardSettings` le falta su pantalla, y la «deuda viva total» de K-3 ($18,999.22) no es la del oráculo ($3,895.00);
  - P-01 Detalle mide 475 px en Android;
  - `verify-wf3` falla en U5-04 por un falso positivo.
  
  La salida está en `docs/estado/bloqueos/WF3c-ronda3-escalamiento.md`.

## Cómo correr la verificación

Desde la raíz del repositorio:

```sh
make verify-wf0      # SQL en base temporal, go build y go vet, colorimetría y validador OpenAPI
make verify-c1       # build.py --check (C1 y extractos idénticos al generador) y verify_c1.py contra la base migrada
make verify-s1       # backend/db/tools/s1_verifica.py: esquema 1.3, oráculo, negativas, Down y goose
make verify-u1       # docs/design/tools/verificar_u1.py: U1 1.2
make verify-criterios  # docs/criterios/tools/verif_final.sh: los 18 archivos de criterios
make verify-all      # todo lo anterior
make verify-wf3      # WF-3: Go, cobertura K25, contrato C1, interfaz (build, rutas, hosts y Playwright; ~5 min) y verify-u5; fuera de verify-all
make verify-u5       # U5 sin servidor sobre la evidencia versionada: U5-01, U5-03, U5-04, U5-08 y U5-11 (segundos)
make verify-u5-servidor  # U5 con bash ops/dev-run.sh corriendo: regenera las 164 capturas, axe y contraste (~9 min) y luego verify-u5
make verify-wf3-final  # verify-wf3 + secretos (S6-13), respaldo (S6-09) y servidor (S6-11, S6-12, S6-14); necesita bash ops/dev-run.sh
make clean           # borra los temporales si una verificación falló a medias
```

Requisitos: GNU Make y bash, `sqlite3` 3.37 o posterior, Go 1.27 o posterior, Python 3.11 o posterior, y `uv`/`uvx` con `openapi-spec-validator` 0.9.0 en su caché (para `verify-c1` también `pyyaml` y `jsonschema`; para `verify-criterios`, `pyyaml`). Si la caché está vacía, la primera corrida necesita internet: `make verify-wf0 UVX_FLAGS= UV_FLAGS=`. Si `goose` v3.28.0 está en el PATH (o con `GOOSE=/ruta/a/goose`), también se prueba su ciclo up, down y up.

Qué comprueba `make verify-wf0`, sin escribir en ningún archivo del repositorio:

1. **SQL.** Aplica la sección Up de cada migración con `sqlite3 -bail` en `datos/tmp/verify-wf0/`. Exige que todas las tablas sean STRICT, `PRAGMA foreign_key_check` sin filas e `integrity_check` = ok. Luego aplica Down en orden inverso (la base debe quedar vacía) y Up otra vez.
2. **Go.** `go build ./...`, `go vet ./...` y `gofmt -l .` en `backend/`.
3. **Colorimetría.** Corre `docs/design/tools/colorimetria.py --sin-png` en un árbol espejo dentro de `datos/tmp`, exige «VERIFICACIÓN OK» y que `colorimetria.md`, `contraste.csv`, `C2-colores.json` y `paleta.html` salgan idénticos byte a byte a los del repositorio, con 0 filas `cumple=false`.
4. **OpenAPI.** `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml`.

Salida esperada al 11-oct-2026 (esquema 1.3): 00001, 00002 y 00003 aplicadas en ese orden, 55 tablas (55 STRICT), 152 índices, 75 disparadores, 0 filas en `foreign_key_check`; Down en orden inverso (00003, 00002, 00001) deja 0 objetos; goose v3.28.0 up, down-to 0 y up sin errores; Go sin errores; 202 pares de contraste con 0 `cumple=false`; `openapi.yaml: OK`. `make verify-all` agrega `build.py --check` (los 21 archivos de C1, extractos incluidos, idénticos al generador), `verify_c1.py --db` (29 comprobaciones, «TODO OK»), `s1_verifica.py` («RESUMEN TODO OK», con las 3 migraciones del comodín; regresión Kueski: suma de pagos de la despensa 1 = 389500 centavos y parte de Axel = 239528 centavos, igual que `oraculo/model.json`), `verificar_u1.py` (63 comprobaciones OK, «TODO OK») y `verif_final.sh` (261 criterios en 18 archivos, 178 cifras canon sin diferencias, 12 de 12 cifras del calendario contra el oráculo). Corrida completa: unos 10 segundos.

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

1. **Puerta 1 (aplicada el 11-oct-2026).** Axel aprobó C1 v1.0, U1 y los criterios y respondió todos los bloqueos (`docs/plan/decisiones.md` §K, K1 a K34). El flujo P1 aplicó las respuestas (esquema 1.2, C1 v1.1, U1 1.1, criterios 1.1 y colorimetría 1.1): [`docs/estado/p1/INTEGRACION.md`](docs/estado/p1/INTEGRACION.md).
2. **Puerta 1b (aplicada el 11-oct-2026).** Axel respondió los 8 bloqueos de P1 y fijó dos reglas (K35 a K44). El flujo P1c las aplicó: esquema 1.3 con la migración 00003, C1 v1.2 regenerado con `backend/internal/api/c1gen/build.py`, U1 1.2, colorimetría 1.2 y criterios 1.2. Antes de arrancar WF-1, el orquestador decide qué hacer con los hallazgos de los criterios que siguen abiertos; las propuestas para la puerta 2 (mínimos de seguridad y el bloqueo `S1-parametros-avisos`) están en [`docs/estado/p1c/INTEGRACION.md`](docs/estado/p1c/INTEGRACION.md).
3. **WF-1 «Motores y sistema de diseño» (integrado el 11-oct-2026).** Motores S2 a S5 con oráculos, recetario, catálogo, PROFECO y OCR, `DESIGN.md`, tokens C2, U3, U4 y el paquete para Claude Design: [`docs/estado/wf1/INTEGRACION.md`](docs/estado/wf1/INTEGRACION.md). `make verify-wf1` lo comprueba.
4. **Puerta 2 (aplicada el 11-oct-2026, flujo P2c).** Axel contestó las 98 preguntas en modo «lo recomendado + excepciones», sin excepciones (`docs/estado/puerta2/RESPUESTAS.md`; `docs/plan/decisiones.md` §K46 a §K52). P2c aplicó: esquema 1.4, C1 v1.3, contratos S2 1.3, S3 1.5, S4 1.5 y S5 1.5, motores S2 y S5, U3 1.5, U4 1.6 y criterios 1.3. Lo que no se aprobó, los bloqueos y los pendientes para WF-3 y la Pista 3 están en [`docs/estado/p2c/INTEGRACION.md`](docs/estado/p2c/INTEGRACION.md). Sigue P2d (registrar K53, la aplicación no lee sitios de tiendas) y después WF-3; el orquestador se detiene en la puerta 4.
5. **Pista 3 (Claude Design por Artifacts), la ejecuta el orquestador (K45, sustituye a K35).** Sin navegador y sin subir archivos: los agentes escriben en disco el sistema de diseño y los artboards de todas las pantallas de U1 con cifras sintéticas, y solo el orquestador los publica con la herramienta Artifact; hay hasta 3 vueltas revisadas por agentes, cada una copiada en `docs/design/salida/iter-N/`, y la exportación final va en `docs/design/salida/final/`. `docs/design/entrada/` queda como brief de referencia. La puerta 3 es al final: Axel abre la URL privada del lienzo y aprueba la referencia visual. Plan técnico: [`docs/design/pista3-plan.md`](docs/design/pista3-plan.md). Después vienen WF-3 (backend, frontend y verificación visual; puerta 4) y WF-4 (operación, semilla, seguridad e informe final).
6. **WF-3a «Construcción» (integrado el 12-oct-2026).** Servidor Go, frontend y PWA, pruebas de contrato y C3 v0.9 en borrador: [`docs/estado/wf3/INTEGRACION-wf3a.md`](docs/estado/wf3/INTEGRACION-wf3a.md). **Siguiente: WF-3b**, con las capturas de U5, los flujos F-1 a F-6 contra el servidor y la firma de C3. WF-3a bis (12-oct-2026) ya corrigió el frontend y aplicó las decisiones del orquestador sobre `docs/estado/bloqueos/C3-matriz.md` (§15 de la integración). Antes de WF-3b falta decidir si se cierran los 2 medios que quedan y la prueba que falla en la corrida completa (§15.8). Después vienen la puerta 4 (valores K31, §9) y WF-4 (instalación con systemd).
7. **WF-3b «Verificación visual y C3» (integrado el 12-oct-2026).** Capturas U5, flujos contra el servidor real y C3 v1.2: [`docs/estado/wf3/INTEGRACION-wf3b.md`](docs/estado/wf3/INTEGRACION-wf3b.md). **Siguiente: puerta 4.** Axel decide los dos escalamientos de ronda 3, `docs/estado/bloqueos/U5-ronda3-escalamiento.md` y `C3-ronda3-escalamiento.md` (se recomienda la A en los dos), prueba S6-11 en sus Android y confirma los valores K31. El orquestador hace S7-07 con el envío real por Telegram. Después viene WF-4, la instalación con systemd.
8. **WF-3c «Inicio como tablero de control» (integrado el 12-oct-2026).** K54 con U3 2.0, DESIGN 1.1, esquema 1.5, C1 v1.5, API, pantalla, HTTP/2 y capturas: [`docs/estado/wf3/INTEGRACION-wf3c.md`](docs/estado/wf3/INTEGRACION-wf3c.md). **Siguiente:** el orquestador decide `WF3c-ronda3-escalamiento`; se propone una ronda acotada «WF-3c bis» antes de la puerta 4. En la puerta 4, Axel contesta:
   - qué es la «deuda viva total»;
   - dónde se editan los parámetros del tablero;
   - qué pasa con las líneas archivadas;
   - los valores PROVISIONALES de K54, junto con lo que ya estaba pendiente de WF-3b.

Lo que WF-1 debe tomar de los flujos anteriores (insumos, pendientes y solicitudes abiertas) está en [`docs/estado/wf0/INTEGRACION.md`](docs/estado/wf0/INTEGRACION.md) §5 y §7, [`docs/estado/p1/INTEGRACION.md`](docs/estado/p1/INTEGRACION.md) §7 y [`docs/estado/p1c/INTEGRACION.md`](docs/estado/p1c/INTEGRACION.md) §7.

## Reglas de trabajo

- Todo agente lee [`docs/plan/protocolo-agentes.md`](docs/plan/protocolo-agentes.md): no asume (lo que falta va como bloqueo), escribe solo en su alcance, no hace commits y deja su reporte en `docs/estado/<flujo>/<agente>.md`.
- Toda cifra que exista en `oraculo/model.json` se reproduce exactamente; si no coincide, es un hallazgo, no un ajuste.
- Documentos e interfaz en español de México; identificadores de código en inglés.
- El integrador de cada flujo hace un solo commit. Git es local, sin remoto.
