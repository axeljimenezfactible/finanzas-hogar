# Plan · Finanzas Hogar (Etapa A → B)

**Estado: LISTO PARA APROBACIÓN.** 70 preguntas hechas con la herramienta de preguntas en 17 rondas, **0 abiertas**. Quedan 2 decisiones que solo Axel toma al aprobar (§8.1): el candidato de matiz y la pila tecnológica recomendada.
**Fecha:** 2026-10-10 · Orquestador Fable 5.1 · Etapa A en modo plan, sin workflows ni código.
**Convención:** cada bloque va marcado **HECHO** (dicho por Axel o leído en archivos), **PROPUESTA** (recomendación mía) o **DECISIÓN AXEL** (ya respondida).

---

## 0. Contexto

Axel y Ximena (esposa) quieren un sistema web de control financiero del hogar que nace del reporte Bodega vs Walmart y del modelo Kueski-Despensa. Karen (madre de Axel) participa solo en la despensa en común y en sus compras particulares con el mismo crédito. El producto es el sistema web (panel + dashboard + configuración estilo Notion); los reportes PDF/Telegram son salidas. Corre en la máquina de casa, en red local, con dos cuentas co-titulares. La Etapa B se ejecuta con Workflow en ultracode, tres pistas en paralelo, contratos en disco y tope de 300K tokens por agente.

---

## 1. HECHOS

### 1.1 Prompt canónico
- El texto pegado en el chat llegó truncado; el canon es `~/prompt_unico_sistema_finanzas_hogar.md` (16 KB, 10-oct 19:41), que absorbe `~/prompt_sistema_finanzas_hogar_por_fases.md` y `~/prompt_ux_ui_dashboard.md`.
- Este prompt **sustituye** la memoria `project_sistema_finanzas_hogar_vision.md` («Axel lo construirá, Claude no»): ahora Claude lo construye en Etapa B. La memoria se actualiza al arrancar la Etapa B.

### 1.2 Base existente (`~/reporte-compra-latex/`)
- `data.json`: compra 10-oct-2026, 61 artículos en Bodega Aurrera Vallejo vs 60 en Walmart Acueducto; cada artículo con `cuenta` ∈ {axel, karen, conjunto, mixto}.
- `model/model.py` (v3, Decimal, ROUND_HALF_UP, asserts) + `model/verify3.py` (segundo cálculo independiente en float); `model.json` y `verify3.json` coinciden. **Patrón modelo + verificador independiente: obligatorio.**
- Cifras canon para regresión (S3): compra $3,094.50; intereses $800.50 (25.8685 % a 4 quincenas; 2 quincenas = 17.27 %); Axel $2,395.28 ($598.82/q), Karen $1,499.72 ($374.93/q); pagos 30-oct, 15-nov, 30-nov, 15-dic-2026; despensa base $4,000 → $5,034.74, pagos $1,258.68 ×3 + $1,258.70; Axel $883.75/q (último $883.77) desde 30-dic; valor estable $2,616.87/mes (Axel $1,767.51, Karen $849.36); deuda viva 10-abr-2027 $6,293.44; lecturas A/B del disponible (mín. $8,094.50 / $5,000.00); papel $199 alterno; en febrero el «30» es el 28; periodos P1–P6 devengado/efectivo con totales $14,362.48 / $13,103.78.
- Reporte PDF 11 páginas, secciones 1–6; paleta previa (descartada por Axel: paleta nueva completa); pipeline synth/ → build/ (LaTeX) → review/ (capturas pdftoppm).

### 1.3 Estructura real del hogar (R5/R6)
- **Axel** titular; **Ximena** esposa y **co-titular** con login propio; **Karen** madre de Axel, **participante externa**: solo entra en la despensa común y en sus compras particulares con el crédito; **no se modelan las finanzas de su familia**; **no tiene acceso**; **paga directo a Kueski** su parte; recibe un estado de cuenta PDF que Axel le reenvía a mano.
- El presupuesto se planifica **solo con el ingreso de Axel** (quincenal 15/30, capturado por él) más extras.

### 1.4 Entorno verificado (2026-10-10)
| | casa (producción, 100.115.28.35 / LAN 192.168.50.56) | servidor 100.77.211.11 (fuera del alcance; solo referencia) |
|---|---|---|
| SO / hardware | Arch, kernel 7.1.5-zen, X11, 16 hilos, 27 GB, 677 GB libres, GPU AMD Vega integrada | Arch, Docker 29.8, RTX 5050 8 GB, 789 GB libres |
| Runtimes | Python 3.14.6, Go 1.26.5, Node 26.5, Bun 1.3.14, Deno 2.9.4, Rust 1.97.1, uv 0.11.32 | Python 3.14.7, Node 26.11.1 |
| Datos | SQLite 3.53.4; cliente psql 18.4 (servicio inactivo); MariaDB escuchando en 3306 | SQLite 3.53.4 |
| Documentos | lualatex, pdflatex, xelatex, latexmk, pdftoppm ✓ · Tesseract 5.5.3 con idioma **spa** ✓ · Chromium ✓ · unrar y bsdtar ✓ · mkcert ✗ (openssl ✓) | sin LaTeX, sin Chromium |
| Red | avahi activo (mDNS); puertos ocupados 22, 3306, 5010, 631; Tailscale activo | `tailscale serve` sin config |
| Telegram | bot MultiX en `~/.openclaw/openclaw.json` (chat_id de Axel; token solo se lee del archivo) | sin integración |

### 1.5 Claude Design (verificado en anthropic.com/news/claude-design-anthropic-labs y claude.com/product/design)
- Anthropic Labs, **research preview**, lanzado 17-abr-2026, URL `claude.ai/design`, incluido en Pro/Max/Team/Enterprise, usa los límites de la suscripción.
- Entradas documentadas: prompts, imágenes, DOCX/PPTX/XLSX, código base, captura web. **Markdown no está listado** → DESIGN.md se sube convertido a PDF (lualatex o pandoc) además del texto en el prompt.
- Salidas: URL, Canva, PDF, PPTX, HTML independiente, paquete de handoff para Claude Code.
- Herramienta `mcp__claude-in-chrome__file_upload`: **solo sube archivos de carpetas compartidas con la sesión**; rechaza otras rutas (riesgo R4).

### 1.6 Kueski Pay (verificado en kueskipay.com, kueski.com y tyc.kueski.com, 2026-10-10)
- Planes «de 2, 4, 6, 8, 10 y 12 quincenas **con intereses**»; promociones por temporada a 2 y 4 quincenas sin intereses; «APR mínima 0 % y máxima de 486 %»; ejemplo oficial: $1,000 → $1,101.87 con IVA, 0.19 % diario, 5 quincenas, CAT promedio 153.31 % (02-mar-2026); «Personalizar tu tasa de interés». **No existe tabla pública por plazo**; la página de comisiones (centro de ayuda, tras Cloudflare, leída con Obscura `--stealth`) dice que «dependerá de la evaluación de crédito».
- Tablero de promociones 2026 (PDF oficial, 1,732 líneas extraídas): **37 promociones «TASA 0 % DE INTERÉS DIARIA»**, casi todas exclusivamente a 4 quincenas, por comercio, con compra mínima ($2,000–$2,500) y vigencia de días (p. ej. Chedraui solo el 30-oct-2026). Walmart, Bodega Aurrera y Sam's **no aparecen**.
- La tasa del reporte (25.87 % a 4 q, saldo decreciente) ≈ 0.69 % diario: tasa personalizada de Axel, dentro del rango oficial.

### 1.7 Fuentes de datos externas verificadas
| Fuente | Estado verificado | Uso previsto |
|---|---|---|
| PROFECO Quién es Quién en los Precios | datos.profeco.gob.mx: datasets 2024/2025/**2026 (RAR)**, metadatos y diccionario CSV descargados. Cobertura ene-2024 a 31-ago-2026, **actualización mensual**. Campos: producto, presentacion, marca, categoria, catalogo, precio, fecha_registro, **cadena_comercial**, giro, nombre_comercial, direccion, estado, municipio, latitud, longitud | precios de referencia por cadena y municipio; cobertura de las sucursales de Axel por confirmar en S5 |
| USDA FoodData Central | API gratuita con llave data.gov (1,000 req/h), **dominio público CC0**; descargas: Foundation Foods 04/2026 (CSV 3.7 MB zip), Branded 04/2026 (CSV 428 MB zip), SR Legacy congelado 2018 | kcal por ingrediente (carga local del CSV Foundation, sin API en runtime) |
| Open Food Facts | ODbL + DbCL; volcados CSV (~0.9 GB gz), JSONL, Parquet; API «1 llamada = 1 escaneo real», scraping bloqueado | productos empacados por código de barras (volcado local filtrado) |
| Tablas mexicanas (SMAE) | **no verificada licencia abierta**; se trata como captura manual citando fuente | complemento manual |

### 1.8 Tecnologías verificadas (versión estable al 2026-10-10, fuente oficial)
| Tecnología | Versión / fecha | Nota |
|---|---|---|
| SQLite | 3.54.0 (2026-10-09) | archivo único; casa tiene 3.53.4 |
| DuckDB | 1.5.6 (2026-09-28) | solo opcional para analítica |
| PostgreSQL | 18.6 (18 soportado a 2030) | descartado por Axel (SQLite) |
| Go | 1.27.2 estable (1.26.9 aún estable) | casa tiene 1.26.5 → actualizar |
| modernc.org/sqlite | ≥ v1.60.1 (2026-09-29), BSD-3 | SQLite puro Go, sin cgo |
| pressly/goose | v3.28.0 (2026-09-02), MIT | migraciones |
| Bun | 1.4.3 (2026-10-10) | descartado (TS no preferido) |
| Python / FastAPI | 3.14 / 0.143.0 | Python queda como **oráculo verificador** |
| React | 19.3.0 (2026-09-09) | pila recomendada |
| Vite / vite-plugin-pwa | 8.3.4 / 2.0.0 (soporta Vite 8) | build + PWA |
| Tailwind CSS | 4.3.3 (2026-07-16) | |
| shadcn CLI / radix-ui | 4.21.4 (2026-10-07) / 1.7.0 (~2026-10-05) | componentes accesibles |
| TanStack Table | 9.2.8 (React y Svelte) | vistas de base de datos |
| dnd-kit | 6.3.1 | arrastrar y soltar |
| BlockNote | 0.55.0, MPL-2.0, **pre-1.0** | editor de bloques estilo Notion; solo para notas |
| Apache ECharts | 6.1.0 (2026-05-19) | gráficas recomendadas |
| Observable Plot / visx / Tremor | 0.6.17 (2026-02) / 4.0.0 (2026-06-11) / 3.18.7 | alternativas evaluadas |
| Svelte / SvelteKit / shadcn-svelte | 5.57.2 / **3.0.1 (2026-10-06, 4 días)** / 1.7.0 | pila alternativa |
| htmx | 4.0.0 (2026-08-28; 2.x mantenido) | descartado por Axel (app en navegador) |
| Style Dictionary | 5.6.1 (2026-10-10), DTCG | tokens → CSS |
| culori | 4.0.2 | cálculos OKLCH en scripts |
| Storybook | 10.6.1 estable (11 alpha) | variantes de componentes |
| Playwright | 1.64.0 (2026-10-07) | capturas y E2E |
| mkcert | 1.4.4 (paquete Arch) | CA local para HTTPS en LAN |
| Tesseract | 5.5.3 + spa (instalado) | OCR de tickets |
| HiGHS / OR-Tools | 1.15.1 (C/cgo) / 9.15 (**sin Go**) | **no se usan**: espacio de escenarios discreto, búsqueda exhaustiva en Go |

---

## 2. Registro de preguntas y respuestas (70, todas cerradas)

| # | Tema | DECISIÓN AXEL |
|---|---|---|
| **Infraestructura** | | |
| 1 | Dónde corre | Esta máquina (casa); el servidor queda fuera |
| 2 | Acceso móvil | Solo red local de casa |
| 4 | Bot Telegram | Bot MultiX existente |
| 39 | PWA | Sí, instalable, certificado propio (mkcert) |
| 49 | Despliegue | Servicio systemd nativo |
| 50 | Login | Usuario y contraseña, sesión larga por dispositivo |
| 54 | Repositorio | `~/finanzas-hogar`, git local sin remoto; nombre provisional «Finanzas Hogar» |
| 57 | Respaldos | Solo local (otra carpeta o disco) |
| 64 | Instalar paquetes | Sí, sin lista previa; se reporta al final |
| **Personas y roles** | | |
| 3 | Karen acceso | Sin acceso; Axel captura lo suyo |
| 17 | Quién es Ximena | Esposa de Axel; Karen es la madre de Axel |
| 21 | Ximena acceso | Co-titular con login propio, mismos permisos |
| 22 | Ingresos Ximena | No entran; solo el ingreso de Axel |
| 23 | Familias | Sistema solo para Axel y Ximena; Karen participante en despensa común y sus compras |
| 25 | Pago de Karen | Paga directo a Kueski su parte |
| 32 | Estado de cuenta Karen | El sistema lo genera; Axel lo reenvía a mano |
| 69 | Avisos Ximena | Sí; registra su chat en el bot MultiX |
| **Reglas financieras** | | |
| 5 | Periodo | Mes calendario; quincena 15/30 subperiodo; despensa bimestral devengada aparte |
| 6 | Ingreso | Quincenal 15/30; Axel captura montos |
| 7 | Reparto conjunto | Porcentaje configurable por gasto/categoría, 50/50 por defecto |
| 8 | Ahorro | Sí: fondo y metas con fecha |
| 24 | Qué se comparte con Karen | Configurable por gasto |
| 56 | Lectura A/B | Interruptor; por defecto B (conservadora) |
| 59 | Fecha y moneda | 10-oct-2026; MXN, $, 2 decimales |
| 66 | Idioma | Solo es-MX |
| **Créditos** | | |
| 9 | Tipos | Revolvente Kueski, tarjeta bancaria, préstamo fijo, crédito de tienda (tiene varios) |
| 10 | Objetivo optimizador | Los 4 (intereses, pago parejo bajo tope, disponible, comprar más) |
| 11 | Restricciones duras | Colchón mínimo, máximo de compras vivas, deuda cero antes de fecha |
| 12 | Tasas Kueski | Buscar en internet → no hay tabla pública; tasa por compra + tabla editable (ver 1.6) |
| 13 | Combinación de objetivos | Comparador de escenarios + puntaje ponderado por defecto |
| **Comida, salud, inventario** | | |
| 14 | Personas | 2 adultos (Axel y Ximena) + datos de salud |
| 15 | Recetas | Recetario inicial de Claude + las de Axel |
| 16 | Inventario | Híbrido; verduras y carne no especializada en mercados; resto en cadenas |
| 18 | Salud | Perfil, IMC, historial de peso, calorías por receta y por día |
| 19 | Nutrición | USDA, Open Food Facts, tablas mexicanas, etiquetas de tiendas cruzadas |
| 20/29 | Mercados | Varios mercados o tianguis; Axel los captura en Ajustes |
| 67 | Dieta | Ninguna restricción |
| 68 | Tiempos | Desayuno, comida, cena y colaciones |
| 70 | Fuera de casa | Gasto normal en su categoría; fuera del plan de comidas |
| **Precios y tiendas** | | |
| 26 | Fuentes | Captura manual, OCR local, lectura asistida de sitios, PROFECO |
| 27 | Sam's | No la tienen; evaluarla (punto de equilibrio) |
| 28 | Frescura | 30 días |
| **Reportes y avisos** | | |
| 30 | Reportes | Diario y semanal Telegram, mensual PDF LaTeX, vista web mensual |
| 31 | Alertas | Pago próximo, presupuesto 80/100 %, precios viejos, plan sin armar |
| 60 | Prefijo Telegram | [Finanzas-Hogar] |
| **UX/UI** | | |
| 33 | Dispositivos | Android los dos + escritorio Linux |
| 34 | Tema | Según el sistema, con interruptor |
| 35/37 | Paleta | Nueva completa; 3 candidatos de matiz a elegir (§4) |
| 36 | Saturación/densidad | Apagada y cómoda, tipo Notion |
| 38 | Ajustes | 11 del prompt + Personas y reparto, Salud y metas, Inventario (14) |
| 40 | Referencias | Solo Notion |
| **Claude Design** | | |
| 41 | Qué es | claude.ai/design; Axel inicia sesión |
| 42 | Materiales | DESIGN.md + tokens, C1 resumido, capturas del reporte (sin código) |
| 43 | Implementar | Referencia visual; se reimplementa con tokens C2 |
| 44 | Iteraciones | 3 |
| 45 | Acciones autorizadas | Pegar prompt/enviar/crear proyecto; subir archivos acordados; descargar exportes y capturas. **No** comentarios en línea ni edición directa |
| 63 | Carpeta de subida | Intentar desde carpeta local; si la herramienta rechaza, pedir a Axel compartirla en esa puerta |
| **Tecnología y proceso** | | |
| 46 | Backend | Go |
| 47 | Interfaz | Fiel a Notion, aplicación en el navegador |
| 48 | Base de datos | SQLite |
| 51 | IA opcional | Modelo local y API de Claude, desactivables; núcleo determinista |
| 52 | Semilla | Compra 10-oct, crédito Kueski, gastos fijos e ingresos en archivo de Axel |
| 53 | Puertas | C1, DESIGN.md+C2, antes de Claude Design, C3 |
| 55 | Tope total | Sin tope total; solo 300K por agente |
| 58 | Históricos | Ninguno |
| 65 | Fecha límite | Sin fecha; calidad primero |

---

## 3. Plan de la Etapa B (PROPUESTA)

### 3.1 Principios operativos
1. **Tope 300K por agente, sin comprimir.** Cada agente se planifica a ≤ 240K estimados (§5). Si un agente se acerca al tope, escribe `docs/estado/<fase>/<agente>.md` con lo hecho y lo pendiente y termina; el orquestador reparte el resto en un agente nuevo.
2. **Comunicación por disco.** Contratos en `docs/contratos/`, criterios en `docs/criterios/`, reportes de agente en `docs/estado/`, bloqueos en `docs/estado/bloqueos/`. Nunca SendMessage a agentes de flujo (memoria: levanta copias concurrentes).
3. **Protocolo de bloqueo.** Un agente que no sabe algo **no asume**: escribe `docs/estado/bloqueos/<fase>-<tema>.md` (qué falta, opciones, impacto), sigue con lo que no depende de eso y termina. El orquestador recoge los bloqueos al cerrar el flujo y los pregunta a Axel en la puerta siguiente (o antes si bloquean la fase).
4. **Criterios antes del trabajo.** WF-0 escribe `docs/criterios/<fase>.md` para todas las fases antes de que arranque ningún agente de trabajo; los revisores evalúan solo contra esos criterios.
5. **Puerta de calidad por fase:** revisor independiente (Opus) aprueba o devuelve, máximo 3 rondas; si no aprueba, se escala a Axel.
6. **Modelos.** Opus 5.5: arquitectura, motores, optimizador, sistema de diseño, especificaciones, revisión crítica, integración. Sonnet 5.5: recetario, cargadores de datos, fixtures y pruebas repetitivas, tablas de contraste, variantes de componentes, capturas. Fable 5.1 orquesta y ejecuta personalmente la Pista 3 (navegador) por ser secuencial y estar acotada por la lista de acciones autorizadas.
7. **Telegram.** Al cerrar cada flujo y cada puerta: mensaje al chat de Axel con prefijo `[Finanzas-Hogar]`, texto plano (sin parse_mode: los nombres de archivo con guiones bajos rompen Markdown), formato Qué / Archivos / Verificación / Estado / Pendientes.
8. **Verificación independiente obligatoria** en todo cálculo: motor Go + oráculo Python (model.py/verify3.py extendidos) con igualdad al centavo.
9. **Semilla.** `datos/semilla/gastos_fijos_ingresos.csv` lo entrega Axel (plantilla en §8.2); si no existe al llegar WF-4, la carga se hace sin esa parte y se deja el importador listo.

### 3.2 Repositorio `~/finanzas-hogar/` (git local)
```
backend/            Go: cmd/finanzas, internal/{dominio,motor_fin,motor_cred,recetas,precios,api,auth,jobs,reportes,telegram,respaldo}
frontend/           Vite + React + TypeScript + Tailwind 4 + shadcn/ui; src/{app,shell,dashboard,ajustes,comidas,creditos,compras,salud,db-views}
oraculo/            Python: copia de model.py/verify3.py + extensiones; fixtures JSON; se ejecuta en CI local
datos/semilla/      data.json, model.json (copias), gastos_fijos_ingresos.csv (Axel), recetario.json, nutricion/, profeco/
docs/plan/          este plan + registro Q&A
docs/criterios/     criterios de aceptación por fase (escritos en WF-0)
docs/contratos/     C1-datos.md + openapi.yaml · C2-tokens.json + tokens.css · C3-integracion.md
docs/design/        DESIGN.md, U1-arquitectura.md, U3-dashboard.md, U4-ajustes.md, entrada/ (paquete Claude Design), salida/ (capturas, HTML), capturas-u5/
docs/estado/        reportes por agente y bloqueos/
ops/                finanzas-hogar.service, respaldo.sh, mkcert.md, runbook.md
```

### 3.3 Secuencia de flujos y puertas
```
WF-0 Cimientos (7 agentes) ──► PUERTA 1: C1 + U1 + criterios
WF-1 Motores + sistema de diseño (18 agentes) ──► PUERTA 2: DESIGN.md + C2 + motores verificados
PUERTA 3: Axel inicia sesión en claude.ai/design y comparte carpeta ──► Pista 3 (orquestador + 2 agentes)
WF-3 Construcción backend + frontend + U5 (18 agentes) ──► PUERTA 4: C3 verificado + U5 aprobado
WF-4 Operación, semilla, seguridad, informe (5 agentes) ──► Informe final + Telegram
```
Las pistas 1 y 2 corren en paralelo dentro de WF-0, WF-1 y WF-3; la Pista 3 va entre las puertas 2 y 4 y alimenta la implementación del frontend como referencia visual.

### 3.4 Pista 1 — Sistema

| Fase | Objetivo | Entregables | Criterios de aceptación (medibles) | Depende de | Riesgos | Agentes (modelo, flujo) | Puerta |
|---|---|---|---|---|---|---|---|
| **S1 Dominio y datos** | Esquema SQLite + diccionario + reglas formalizadas | `backend/db/migrations/0001…`, `docs/contratos/diccionario-datos.md`, `docs/contratos/reglas-negocio.md` | 100 % de entidades de §3.4.1 presentes; migraciones aplican en SQLite limpio; cada regla con ejemplo numérico del reporte; cada tabla referenciada por ≥1 pantalla U1 o ≥1 motor | Q&A | sobre-modelado | 1 Opus (WF-0) | Revisor WF-0 |
| **S2 Motor financiero** | Presupuesto mensual, promedios diario/semanal/mensual, proyección 12 meses devengado y efectivo, flujo por quincena, metas de ahorro, real vs plan | `internal/motor_fin`, `oraculo/fin_*.py`, fixtures | Reproduce las ~40 cifras canon de 1.2 al centavo; oráculo Python coincide en 100 % de fixtures; invariantes: Σ partes = total, devengado acumulado − efectivo acumulado = deuda viva | S1 | redondeo acumulado | motor Opus; oráculo+fixtures Sonnet (WF-1) | Revisor motores |
| **S3 Motor de créditos** | 4 tipos de crédito, calendario unificado, promociones manuales, comparador de escenarios con 4 objetivos y 3 restricciones duras | `internal/motor_cred`, `oraculo/cred_*.py`, `docs/estado/S3-escenarios.md` | Regresión Kueski: 800.50, 973.75, 598.82/374.93, 5,034.74, 1,258.68/.70, 883.75/.77, series A y B (14 valores c/u), 6,293.44, fechas incl. 28-feb; tarjeta: intereses 0 si se paga total antes de fecha de pago; préstamo: tabla francesa cuadra al centavo; optimizador: para 3 casos de prueba produce 4 óptimos distintos y respeta restricciones en 100 % de escenarios (verificado por oráculo) | S1, S2 | tasa personalizada; sin solver en Go | núcleo Opus; optimizador Opus; regresión Sonnet; integrador Opus (WF-1) | Revisor motores |
| **S4 Recetas y comidas** | Recetas con ingredientes exactos y presentación comprable, escalado, plan semanal por persona (4 tiempos, «fuera»), lista de compras neteada con inventario híbrido, sobrantes, costo por receta/día/semana, perfil de salud, IMC, meta kcal, kcal por receta y día | `internal/recetas`, `datos/semilla/recetario.json` (60 recetas), `datos/semilla/nutricion/` | 100 % recetas válidas contra esquema; escalado exacto (2→4 porciones); lista = Σ ingredientes − stock; kcal/día con fuente citada por ingrediente (USDA/OFF/manual); 5 recetas verificadas a mano por el revisor | S1, S5 (precios) | nombres en inglés (USDA) | motor Opus; recetario Sonnet; cargador nutrición Sonnet (WF-1) | Revisor recetas+precios |
| **S5 Precios y tiendas** | Ingesta multi-fuente, equivalencia por unidad base, mejor tienda por ingrediente y canasta, frescura 30 d, membresía Sam's (punto de equilibrio), promociones, mercados con día de visita | `internal/precios`, cargador PROFECO, OCR Tesseract, `docs/estado/S5-terminos.md`, `docs/estado/S5-cobertura-profeco.md` | Reproduce comparativa del reporte (33 idénticos, 3 sustitutos, +$138.50 / +$100.50); PROFECO 2026 cargado y cobertura por cadena/municipio reportada; OCR ≥ 80 % renglones correctos en 5 tickets de prueba; veredicto de términos por sitio **antes** de activar lectura asistida (desactivada por defecto) | S1 | ToS de tiendas; RAR grande; cobertura PROFECO | comparador Opus; PROFECO+OCR Sonnet (WF-1) | Revisor recetas+precios |
| **S6 Backend local** | API JSON (OpenAPI = C1), auth 2 cuentas, sesiones largas, jobs, respaldos, adjuntos, HTTPS mkcert, PWA servida por el binario, adaptadores IA opcionales | `backend/` completo, `docs/contratos/openapi.yaml` | `go vet` limpio; `go test` ≥ 85 % en dominio y motores; cada endpoint con prueba; login bloquea tras 5 intentos/15 min; respaldo diario `VACUUM INTO` + `integrity_check` + restauración probada; HTTPS válido en Chrome Android con CA instalada | S1–S5 | tope 300K → 3 agentes disjuntos | API núcleo Opus; API motores Opus; auth+jobs+respaldos+TLS Opus; pruebas Sonnet (WF-3) | Revisores backend ×2 |
| **S7 Reportes** | Diario y semanal Telegram (Axel y Ximena), mensual PDF LaTeX día 1, estado de cuenta Karen PDF, vista web mensual, 4 alertas | `internal/reportes`, plantillas `.tex`, `internal/telegram` | PDF compila sin errores; cifras del PDF = API (cruce automático); envío real al chat de Axel HTTP 200; alertas disparan en fixtures de fecha | S2, S3, S6 | Markdown Telegram; tamaño PDF | reportes+Telegram Opus (WF-3) | Revisor backend |
| **S8 Calidad y operación** | Pruebas, seguridad, privacidad, despliegue systemd, runbook, informe final | `ops/`, `docs/estado/informe-final.md` | `govulncheck` y `npm audit` sin críticos; cabeceras de seguridad; servicio `enabled` y vivo tras reinicio; semilla cargada y cruzada con model.json; informe con evidencia y lista de lo no verificado | todo | — | ops Opus; semilla Sonnet; seguridad Opus; informe Opus; revisor Opus (WF-4) | Revisor final |

**3.4.1 Entidades de S1 (HECHO + PROPUESTA de nombres):** persona (Axel, Ximena; Karen marcada externa), cuenta_usuario (2), sesion_dispositivo, periodo (mes calendario, quincena 15/30 con ajuste a fin de mes), ingreso (fijo quincenal, extra), categoria, gasto (fijo/variable/venidero; devengado vs efectivo; adjunto), regla_reparto (hogar / Karen / conjunto con %), meta_ahorro + movimiento_ahorro, credito (tipo, límite, lectura A/B, fechas de corte y pago, comisiones) + tasa_por_plazo, compra_credito (capital, intereses cotizados con IVA, plazo, partes por persona), pago_programado, promocion (comercio, plazo, mínimo, vigencia), tienda (cadena / sucursal / mercado con día de visita; membresía con costo anual), producto, presentacion (unidad base y factor), precio (tienda, fecha, fuente, vigencia 30 d), receta, ingrediente_receta, tiempo_comida (desayuno, comida, cena, colación), plan_comida (día × tiempo × persona × receta × porciones | «fuera»), inventario_item (stock por unidad o ritmo de consumo), perfil_salud + peso_historico, nutricion_ingrediente (kcal/100 g, fuente), notificacion + destinatario (chat_id), job, respaldo. Reglas: Decimal a centavos ROUND_HALF_UP, último pago ajusta centavos, devengado = total de la despensa bimestral / 2 por mes, moneda MXN, zona America/Mexico_City, fechas `10-oct-2026`.

### 3.5 Pista 2 — UX/UI

| Fase | Objetivo | Entregables | Criterios | Depende de | Riesgos | Agentes | Puerta |
|---|---|---|---|---|---|---|---|
| **U1 Arquitectura de información** | Mapa del sitio, barra lateral (Inicio, Gastos, Presupuesto, Comidas, Compras y tiendas, Créditos, Salud, Reportes, Ajustes), inventario de ~25 pantallas con subsecciones Resumen/Detalle/Tendencias/Calendario, 6 flujos clave, vistas de BD por entidad, menú `/`, atajos | `docs/design/U1-arquitectura.md` | Matriz pantalla ↔ endpoint C1 completa; registrar gasto desde móvil ≤ 3 toques; cada flujo con pasos numerados y pantalla por paso | Q&A, S1 borrador | — | 1 Opus (WF-0) | Revisor WF-0 |
| **U2 Sistema de diseño y DESIGN.md** | DESIGN.md completo (§4) con el candidato elegido; tokens DTCG + CSS + tema Tailwind; tabla de contraste calculada para todos los pares usados | `docs/design/DESIGN.md`, `docs/contratos/C2-tokens.json`, `tokens.css` | 0 pares bajo umbral (texto 4.5:1, UI 3:1) en claro y oscuro; tokens compilan con Style Dictionary; paleta de gráficas distinguible en gris (ΔL ≥ 0.10 entre series) y simulación de daltonismo documentada; tipografía e iconos autoalojados | elección de matiz (Axel) | colisiones de matiz | colorimetría Opus (WF-0); DESIGN.md Opus + tokens/contraste Sonnet (WF-1) | Revisor diseño |
| **U3 Dashboard** | Especificación de Resumen (≤ 6 tarjetas: mes vs plan, flujo quincenal, próximo pago y disponible A/B, comidas de la semana con costo y kcal, mejor canasta, ahorro) + Detalle/Tendencias/Calendario + participación de Karen | `docs/design/U3-dashboard.md` | Cada tarjeta mapeada a métricas C1; estados vacío/carga/error definidos; móvil primero (412 px) | U1, C1 | sobrecarga | 1 Opus (WF-1) → implementación F2 (WF-3) | Revisor diseño; U5 |
| **U4 Ajustes estilo Notion** | 14 secciones, barra lateral colapsable, propiedades, vistas de BD (tabla, tablero, calendario, galería), menú `/`, edición en el lugar, vista previa del efecto (p. ej. cambiar A/B recalcula) | `docs/design/U4-ajustes.md` | 14 secciones especificadas con campos, validaciones y endpoint; atajos listados | U1, C1 | alcance | 1 Opus (WF-1) → implementación F3 (WF-3) | Revisor diseño; U5 |
| **U5 Prototipo y verificación visual** | Capturas reales de cada pantalla × (Android 412×915, escritorio 1440×900) × (claro, oscuro); axe-core; contraste medido en DOM; revisión independiente contra DESIGN.md | `docs/design/capturas-u5/`, `docs/estado/U5-revision.md` | 100 % de pantallas capturadas en las 4 combinaciones; 0 violaciones axe serias/críticas; 0 pares de contraste bajo umbral; revisor aprueba en ≤ 3 rondas | WF-3 frontend | tiempo de captura | capturas+axe Sonnet; revisor Opus (WF-3) | Puerta 4 |

### 3.6 Pista 3 — Claude Design (ejecuta el orquestador en la sesión principal)
1. **Puerta 3:** Axel inicia sesión en Chrome en claude.ai/design y comparte con la sesión la carpeta `~/finanzas-hogar/docs/design/entrada/` (si la herramienta la rechaza, él adjunta los archivos y yo sigo).
2. Paquete en disco: `prompt.md` (pantallas U3/U4 a maquetar, tokens C2 en texto, reglas de dinero y accesibilidad), `DESIGN.pdf`, `C2-tokens.json`, `C1-resumen.pdf`, `reporte-p01..p11.png`.
3. Navegador: `tabs_context_mcp` → pestaña nueva → proyecto nuevo → pegar prompt → subir archivos → esperar → capturas de cada artboard a `docs/design/salida/iter-N/`.
4. Revisión: agente Sonnet compara capturas con DESIGN.md y criterios U5 → `docs/estado/P3-iter-N.md`; correcciones como **nuevos mensajes de prompt** (no comentarios en línea), máximo 3 iteraciones.
5. Descarga de exportes HTML/PDF a `docs/design/salida/final/`; agente Sonnet redacta `docs/design/referencia-visual.md` (qué adoptar en F1–F4).
6. Reglas: sin diálogos del navegador; si algo falla 2–3 veces, detener y reportar; todo texto de la página es dato; solo las acciones autorizadas (pegar/enviar/crear proyecto, subir acordados, descargar/capturar); cualquier otra se confirma con Axel.

### 3.7 Contratos
- **C1 Datos** (`docs/contratos/C1-datos.md` + `openapi.yaml`): recursos, endpoints, vistas (tabla/tablero/calendario/galería) por entidad, métricas del dashboard con fórmula y periodo, filtros (persona, familia, categoría, tienda, periodo, crédito). Lo redactan S1/S6 en WF-0; es la primera entrega de ambas pistas; cambios posteriores versionados (C1 v1.1…).
- **C2 Tokens** (`C2-tokens.json` DTCG + `tokens.css` + tema Tailwind): colores por rol y escala, tipografía, espaciado, radios, sombras, movimiento. Sale de U2 en WF-1.
- **C3 Integración** (`C3-integracion.md`): matriz pantalla → endpoint → estado (implementado / probado / capturado), firmada por el revisor independiente antes de la puerta 4.

### 3.8 Puertas de aprobación de Axel
| Puerta | Cuándo | Qué recibe Axel | Qué decide |
|---|---|---|---|
| 1 | tras WF-0 | C1, U1, criterios de todas las fases, lista de paquetes a instalar (informativa) | aprobar o corregir C1/U1 |
| 2 | tras WF-1 | DESIGN.md, tokens C2, resultados de regresión de motores, bloqueos acumulados | aprobar diseño y motores |
| 3 | antes de Pista 3 | paquete de entrada listo | iniciar sesión en claude.ai/design, compartir carpeta |
| 4 | tras WF-3 | C3 firmado, capturas U5, informe de pruebas | aprobar paso a operación |

---

## 4. Esqueleto de `DESIGN.md` (PROPUESTA) con colorimetría razonada

### 4.1 Índice del documento
1. Principios y tono: completo pero minimalista, progressive disclosure, una acción primaria por pantalla, Notion como única referencia, móvil primero.
2. Colorimetría (4.2–4.6).
3. Tipografía: familia variable autoalojada (propuesta: Inter, licencia OFL; alternativa: fuentes del sistema), escala 12/14/16/20/24/32, **cifras tabulares** (`font-variant-numeric: tabular-nums`) para todo dinero, negativo con signo y color de rol, nunca solo color.
4. Espaciado 4 pt (4, 8, 12, 16, 24, 32, 48), rejilla 4/8/12 columnas, radios 6/10/14, sombras de 2 niveles en claro y bordes en oscuro.
5. Iconografía: Lucide (ISC) autoalojada, trazo 1.5, 20/24 px.
6. Movimiento: 150–250 ms, curvas estándar/enfática, `prefers-reduced-motion` desactiva transiciones no esenciales.
7. Componentes: botones, campos, tablas y vistas de BD (tabla, tablero, calendario, galería), tarjetas de métrica, gráficas, calendario de pagos, barra lateral colapsable, menú `/`, panel de Ajustes, estados vacío/carga/error, toasts.
8. Patrones de visualización: barras para comparar categorías, líneas para tendencias, área apilada devengado vs efectivo, dona solo para composición ≤ 5 partes, etiquetas de dinero siempre en texto, no solo tooltip; reglas de la skill dataviz.
9. Accesibilidad y puntos de quiebre: 360, 412, 768, 1024, 1440; foco visible; objetivo táctil ≥ 44 px; WCAG 2.2 AA.
10. Tokens: `C2-tokens.json` (DTCG) + `tokens.css` (variables `--color-*`, `--space-*`, …) + tema Tailwind 4.

### 4.2 Método (fiel al círculo cromático)
- Espacio **OKLCH**; escalas de **11 pasos** (50–950) con luminosidad fija por paso (L = 0.985, 0.962, 0.925, 0.865, 0.785, 0.695, 0.605, 0.515, 0.425, 0.335, 0.255) y croma por regla (perfil 0.10→1.0→0.45 del croma máximo), con recorte al gamut sRGB por bisección. Mismo L en todos los matices → luminosidad percibida uniforme.
- **Armonía complementaria dividida**: matiz base (marca, superficies tintadas, acento principal) + eje opuesto (base + 180°) reservado para error y acciones críticas + dos acentos a ±30° del opuesto para advertencia y gasto. Justificación: el opuesto garantiza el máximo contraste de matiz para lo que debe interrumpir; los acentos divididos evitan el choque directo rojo/verde.
- **60-30-10**: 60 % neutros tintados con el matiz base (croma 0.012), 30 % escala de marca en pasos claros (50–200 claro / 900–950 oscuro), 10 % acentos semánticos y categóricos.
- **Saturación apagada** (croma máx. 0.11 marca, 0.14 roles a L 0.56 en claro, 0.12 a L 0.74 en oscuro), tipo Notion.
- **Roles semánticos** (8): ingreso, gasto, ahorro, deuda/crédito, advertencia, error, éxito, información. **Personas** (Axel, Ximena, Karen, conjunto) y **tiendas** (Walmart, Bodega Aurrera, Sam's, mercado) van en una **paleta categórica aparte** con luminosidad escalonada (distinguible en gris) y separación ≥ 20° de los roles; en gráficas se refuerzan con forma, orden fijo y etiqueta.
- **Contraste calculado** (WCAG 2.x) para cada par usado; umbrales: texto ≥ 4.5:1, texto grande y UI/gráficas ≥ 3:1, en ambos temas.

### 4.3 Candidatos de matiz (valores calculados hoy con el método 4.2; DECISIÓN AXEL al aprobar)

**A · Esmeralda** — base **155°**, opuesto **335°** (magenta) para error/crítico. Lectura: crecimiento, ahorro, comida fresca.
- Marca: 50 `#F5FCF7` · 100 `#E8F7EC` · 200 `#D0EFDA` · 300 `#AEE1BE` · 400 `#87CB9F` · 500 `#60B17F` · 600 `#439564` · 700 `#30784E` · 800 `#225B3A` · 900 `#154027` · 950 `#0B2A18`
- Neutro tintado: 50 `#F9FAFA` · 100 `#F1F3F2` · 200 `#E4E7E5` · 300 `#CFD4D0` · 400 `#B4BBB6` · 500 `#979F9A` · 600 `#7C847F` · 700 `#636965` · 800 `#4B504C` · 900 `#343835` · 950 `#212422`
- Roles (claro 600 sobre fondo 50 → contraste · oscuro 400 sobre 950 → contraste): éxito 155° `#008B4F` 4.18 · `#66C189` 7.13 — error 335° `#A55095` 4.80 · `#D98DC9` 6.42 — advertencia 75° `#9C6900` 4.53 · `#D79F4C` 6.68 — info 240° `#007CB7` 4.40 · `#5CB4EF` 6.89 — ingreso 125° `#628202` 4.26 · `#99B860` 7.00 — gasto 25° `#B84D49` 4.79 · `#ED8C84` 6.48 — ahorro 200° `#00858A` 4.25 · `#24C1C9` 7.12 — deuda 290° `#7563C0` 4.67 · `#AA9DF1` 6.57.
- Observación: éxito e ingreso comparten familia verde con la marca → en A, ingreso se separa a 125° (lima) y se refuerza con icono.

**B · Teal** — base **195°**, opuesto **15°** (rojo) para error/crítico. Lectura: calma, control, sobriedad.
- Marca: 50 `#F2FDFC` · 100 `#E2F7F7` · 200 `#C5F0EF` · 300 `#9AE2E1` · 400 `#68CCCB` · 500 `#2AB2B2` · 600 `#009495` · 700 `#007777` · 800 `#005A5B` · 900 `#004040` · 950 `#002929`
- Neutro: 50 `#F9FAFA` · 100 `#F1F3F3` · 200 `#E3E7E7` · 300 `#CDD4D4` · 400 `#B2BBBB` · 500 `#959F9F` · 600 `#7A8484` · 700 `#616A69` · 800 `#495050` · 900 `#333838` · 950 `#202424`
- Roles: éxito 195° `#008686` 4.22 · `#25C2C2` 7.15 — error 15° `#B84C59` 4.76 · `#EC8A92` 6.43 — advertencia 75° `#9C6900` 4.53 · `#D79F4C` 6.69 — info 240° `#007CB7` 4.40 · `#5CB4EF` 6.89 — ingreso 150° `#218A45` 4.20 · `#6FC082` 7.12 — gasto 10° `#B74B61` 4.79 · `#EB8A99` 6.44 — ahorro 180° `#008777` 4.24 · `#39C3AE` 7.16 — deuda 290° `#7563C0` 4.67 · `#AA9DF1` 6.57.
- Observación: ahorro (180°) queda a 15° de la marca → en B se mueve a 210° o se expresa con la marca misma; error y gasto están a 5° → gasto se separa a 35° (terracota).

**C · Índigo** — base **280°**, opuesto **100°** (oliva-amarillo) para énfasis crítico. Lectura: Notion/Linear, sobrio y moderno.
- Marca: 50 `#F9FAFF` · 100 `#F0F1FF` · 200 `#E1E4FF` · 300 `#C9CFFF` · 400 `#ADB3F4` · 500 `#9095E0` · 600 `#7579C2` · 700 `#5D609F` · 800 `#46487B` · 900 `#303258` · 950 `#1E203A`
- Neutro: 50 `#FAFAFB` · 100 `#F2F2F4` · 200 `#E5E6E9` · 300 `#D1D2D8` · 400 `#B7B8C0` · 500 `#9B9CA4` · 600 `#808189` · 700 `#66676E` · 800 `#4D4E53` · 900 `#36363A` · 950 `#222325`
- Roles: éxito 280° (marca) → **se reasigna a 150°** porque éxito y marca no deben confundirse; error 100° `#857500` 4.43 · `#BDAC4A` 6.86 (un error amarillo-oliva es poco convencional: en C se propone error en 20° rojo y reservar el opuesto 100° para «acción crítica» de confirmación) — advertencia 75° `#9C6900` 4.54 · `#D79F4C` 6.71 — info 240° `#007CB7` 4.41 · `#5CB4EF` 6.91 — ingreso 150° `#218A45` 4.21 · `#6FC082` 7.14 — gasto 25° `#B84D49` 4.80 · `#ED8C84` 6.51 — ahorro 200° `#00858A` 4.26 · `#24C1C9` 7.15 — deuda 310° `#8D5AB1` 4.76 · `#C295E3` 6.52.

**Contrastes estructurales calculados (iguales en los tres candidatos, ±0.1):** texto 900 / fondo 50 = **11.4:1**; texto secundario 700 / fondo 50 = **5.3:1**; texto 600 / fondo 100 = 3.45:1 (**solo texto grande o iconos**); blanco / marca 600 = 3.7–4.0:1 (**insuficiente para texto de botón** → el botón primario usa **marca 700** de relleno: marca 700 / fondo 50 = 5.1–5.5:1 y el blanco contrasta aún más); tema oscuro: texto 100 / fondo 950 = **14.1:1**, texto 300 / fondo 900 = **7.9:1**, marca 300 / fondo 950 = **10.4–10.7:1**; bordes 200/50 = 1.2:1 y 800/950 = 1.9:1 (**decorativos**; los bordes de campos usan neutro 600 en claro ≥ 4.2:1 y neutro 400 en oscuro).

**Paleta categórica propuesta (independiente del candidato, L escalonada para gris):** Axel 250° L 0.55 `#0F74C5` · Ximena 350° L 0.68 `#D375A4` · Karen 80° L 0.78 `#E0AE57` · Conjunto 250° L 0.45 croma 0.03 `#495766` · Walmart 220° L 0.72 `#2FB5D8` · Bodega Aurrera 60° L 0.80 `#FFA658` · Sam's 175° L 0.60 `#11957C` · Mercado 30° L 0.50 `#9C4438`. Colisiones a resolver en U2: Karen (80°) vs advertencia (75°) → Karen a 95° o advertencia a 60°; Bodega (60°) vs advertencia; Sam's (175°) vs ahorro/éxito en B. Regla: separación ≥ 20° o diferencia de L ≥ 0.15 y refuerzo con forma.

### 4.4 Qué entrega U2 sobre este esqueleto
Escalas finales del candidato elegido en claro y oscuro, roles con colisiones resueltas, paleta de gráficas de 8 series (orden de uso, simulación de daltonismo con culori), **tabla completa de contraste** por script (todos los pares componente/estado/tema), tokens DTCG, CSS y tema Tailwind.

---

## 5. Presupuesto de tokens y agentes (PROPUESTA)

**Método de estimación:** tokens por agente ≈ 25K base (sistema, herramientas, prompt y criterios) + lectura de entradas (bytes de archivos × 0.3 × 1.3 por relecturas) + salida (líneas escritas × 12 × 1.6 por iteraciones) + 25 % de uso de herramientas. **Margen de error ±35 %** (sin calibración previa en este dominio; se calibra con WF-0 y se recalcula el resto en la puerta 1). Ningún agente se planifica arriba de 240K. Las rondas de corrección de las puertas de calidad (≤ 3) pueden añadir hasta 30 % por fase.

| Flujo / fase | Tarea | Modelo | Tokens est. | Agentes trabajo | Revisor | Integrador | Total fase |
|---|---|---|---|---|---|---|---|
| WF-0 S1 | Modelo de datos, diccionario, reglas | Opus | 130K | 1 | compartido | compartido | 130K |
| WF-0 S6a | Contrato C1 + esqueleto API + OpenAPI | Opus | 150K | 1 | | | 150K |
| WF-0 U1 | Arquitectura de información | Opus | 110K | 1 | | | 110K |
| WF-0 U2a | Colorimetría final del candidato + escalas + contraste | Opus | 120K | 1 | | | 120K |
| WF-0 crit. | Criterios de aceptación de todas las fases | Sonnet | 90K | 1 | | | 90K |
| WF-0 rev/int | Revisor (coherencia S1/C1/U1) · Integrador (C1 v1, árbol, scripts) | Opus | 140K + 110K | | 1 | 1 | 250K |
| **WF-0** | | | | **5** | **1** | **1** | **850K** |
| WF-1 S2 | Motor financiero · oráculo + fixtures | Opus · Sonnet | 200K · 130K | 2 | | | 330K |
| WF-1 S3 | Motor de créditos · optimizador de escenarios · regresión Kueski | Opus · Opus · Sonnet | 220K · 200K · 140K | 3 | | | 560K |
| WF-1 S4 | Motor recetas/comidas/salud · recetario 60 recetas · cargador nutrición | Opus · Sonnet · Sonnet | 190K · 160K · 130K | 3 | | | 480K |
| WF-1 S5 | Comparador precios/tiendas · PROFECO + OCR | Opus · Sonnet | 200K · 150K | 2 | | | 350K |
| WF-1 U2b/c | DESIGN.md completo · tokens + tabla de contraste | Opus · Sonnet | 210K · 100K | 2 | | | 310K |
| WF-1 U3/U4 | Especificación dashboard · especificación Ajustes | Opus · Opus | 150K · 150K | 2 | | | 300K |
| WF-1 rev/int | Revisor motores · revisor recetas+precios · revisor diseño · integrador | Opus ×4 | 180K · 170K · 170K · 150K | | 3 | 1 | 670K |
| **WF-1** | | | | **14** | **3** | **1** | **3,000K** |
| Pista 3 | Revisor visual de iteraciones · redactor de referencia visual (orquestador ejecuta el navegador) | Sonnet · Sonnet | 120K · 80K | 2 | | | 200K (+~300K sesión principal) |
| WF-3 S6 | API núcleo · API motores · auth+jobs+respaldos+TLS · pruebas backend | Opus ×3 · Sonnet | 230K · 230K · 200K · 160K | 4 | | | 820K |
| WF-3 S7 | Reportes LaTeX + Telegram + estado Karen | Opus | 190K | 1 | | | 190K |
| WF-3 F1–F4 | Shell/sidebar/menú `/`/vistas BD · dashboard+gráficas · Ajustes 14 secciones · comidas/créditos/compras/salud | Opus ×4 | 230K ×4 | 4 | | | 920K |
| WF-3 F5 | Variantes de componentes + Storybook | Sonnet | 150K | 1 | | | 150K |
| WF-3 U5 | Capturas Playwright + axe · revisor visual independiente | Sonnet · Opus | 150K · 200K | 1 | 1 | | 350K |
| WF-3 rev/int | Revisores backend ×2 · revisores frontend ×2 · integrador backend · integrador frontend/C3 | Opus ×6 | 180K ×4 · 150K · 180K | | 4 | 2 | 1,050K |
| **WF-3** | | | | **11** | **5** | **2** | **3,480K** |
| WF-4 | Ops systemd/mkcert/respaldos/runbook · semilla · seguridad · informe final · revisor final | Opus · Sonnet · Opus · Opus · Opus | 140K · 120K · 150K · 140K · 150K | 4 | 1 | | 700K |
| **TOTAL** | | | | **36** | **10** | **4** | **≈ 8.2M (rango 5.4M–11.1M con ±35 %; hasta ≈ 10.7M con rondas de corrección)** |

**Justificación de modelos.** Opus en todo lo que decide estructura o cifras (S1–S3 motores, S6 API, F1–F4, DESIGN.md, especificaciones, revisores e integradores). Sonnet en volumen mecánico con criterios cerrados (recetario, cargadores PROFECO/USDA/OFF, fixtures, pruebas repetitivas, tabla de contraste por script, capturas, Storybook).

---

## 6. Comparativa de tecnologías (PROPUESTA; DECISIÓN AXEL al aprobar)

Las tres pilas cumplen las decisiones ya tomadas (casa, LAN, SQLite, systemd, PWA, app en navegador fiel a Notion); difieren en lenguaje y en el marco de la interfaz. Explicado en resultado visible / costo / riesgo.

| Criterio | **Pila A (recomendada): Go + React** | Pila B: Go + Svelte | Pila C: Python + React |
|---|---|---|---|
| Qué se ve | Igual en las tres: la misma app estilo Notion | igual | igual |
| Velocidad | Binario Go arranca en ms; interfaz compilada con Vite 8 | igual en servidor; Svelte genera menos JS (carga algo más rápida en móvil) | servidor más lento al arrancar; sobra para 2 usuarios |
| Simplicidad | 1 binario + 1 archivo SQLite + carpeta estática; interfaz con la caja de herramientas más completa para Notion (editor de bloques, arrastrar, tablas) | igual en servidor; menos piezas listas para el estilo Notion | 1 proceso Python + dependencias; el modelo v3 ya está en Python |
| Madurez (verificado 1.8) | Go 1.27.2, React 19.3 (sep-2026), Tailwind 4.3.3, shadcn 4.21, TanStack 9.2, dnd-kit 6.3, ECharts 6.1: todo estable y con años de uso | Svelte 5.57 estable; **SvelteKit 3.0.1 salió hace 4 días**; shadcn-svelte 1.7 | FastAPI 0.143 estable; React igual que A |
| Costo de operación | Mínimo: sin runtime externo, sin servicio de BD | mínimo | Python + entorno virtual que mantener |
| Curva para Axel | Lee Go (Matriz de Proyectos); la interfaz no la editaría a mano en ninguna pila | igual | lee Python mejor; la interfaz igual |
| Riesgo | Bajo. Único punto pre-1.0: BlockNote 0.55 (solo para notas; sustituible por campos simples) | Medio: marco mayor recién publicado; menos componentes Notion listos | Bajo-medio: dos lenguajes igual que A, pero sin la velocidad ni el binario único que Axel pidió |
| Verificación cruzada | Motor Go vs oráculo Python (dos lenguajes = independencia real) | igual | motor y oráculo en el mismo lenguaje: menos independencia |

**Recomendación: Pila A.** Componentes por capa (todos verificados en 1.8):
- Datos: SQLite 3.54 vía modernc.org/sqlite (puro Go) + goose; `VACUUM INTO` para respaldos; DuckDB solo si la analítica lo exige (no previsto).
- Backend: Go 1.27, `net/http` estándar (enrutador 1.22+), OpenAPI 3.1 generado, argon2id (`golang.org/x/crypto`), cookies HttpOnly/Secure/SameSite, jobs con `time.Ticker` + tabla `job` (sin dependencia externa), Telegram por HTTP directo, LaTeX y Tesseract por `exec`.
- Motor de optimización: búsqueda exhaustiva sobre escenarios discretos (crédito × plazo × fecha × promoción) con restricciones duras; HiGHS y OR-Tools descartados (sin Go nativo).
- Frontend: Vite 8 + React 19 + TypeScript + Tailwind 4 + shadcn/ui 4 (Radix 1.7) + TanStack Table 9 (vistas de BD) + dnd-kit 6 (tablero y orden) + BlockNote 0.55 (solo notas) + ECharts 6 (gráficas, con decals para daltonismo) + vite-plugin-pwa 2.
- Diseño y pruebas: Style Dictionary 5.6 (tokens DTCG → CSS/Tailwind), culori 4 (cálculo OKLCH), Storybook 10.6, Playwright 1.64 (capturas y E2E), axe-core.
- Operación: systemd (unidad propia, usuario `finanzas`), mkcert 1.4.4 (CA local, certificado para `finanzas.local`), avahi (mDNS ya activo), respaldo local versionado.

Descartadas por decisión de Axel: Postgres 18.6, Bun 1.4.3, htmx 4.0 (páginas de servidor). Descartadas por verificación: OR-Tools (sin Go), HiGHS (cgo), Tremor (dashboard React opinado, sin releases públicos; ECharts cubre más), Observable Plot (SVG estático, menos interacción), visx (bajo nivel, más código).

---

## 7. Riesgos

| # | Riesgo | Impacto | Mitigación |
|---|---|---|---|
| R1 | Casa es único punto de falla; respaldo solo local (decisión de Axel) | pérdida total de historial si falla el disco | respaldo diario versionado + `integrity_check` + restauración probada; si existe segundo disco, copia ahí; recomendación futura: copia al servidor |
| R2 | IP LAN dinámica (192.168.50.56 por DHCP) | la PWA deja de abrir | nombre mDNS `finanzas.local` (avahi activo) + sugerir reserva DHCP en el router |
| R3 | HTTPS en LAN: la CA de mkcert debe instalarse en los dos Android (acción manual) | sin HTTPS no hay PWA instalable | `ops/mkcert.md` con pasos; verificación en puerta 4 desde el teléfono |
| R4 | Claude Design: research preview cambia sin aviso; `file_upload` exige carpeta compartida; Markdown no documentado; límites del plan | Pista 3 bloqueada o degradada | PDF del DESIGN.md + texto en el prompt; carpeta compartida en puerta 3; si falla 2–3 veces, detener y reportar; la Pista 3 es referencia visual, no dependencia |
| R5 | BlockNote pre-1.0 | cambios de API | solo para notas; aislado tras una interfaz propia; sustituible |
| R6 | Kueski: tasa personalizada y reposición de capital sin confirmar | cifras de disponible erróneas | tasa por compra desde la cotización; interruptor A/B por defecto B; aviso visible «pendiente confirmar» |
| R7 | PROFECO: cobertura desconocida de las sucursales de Axel; RAR mensual | fuente de precios poco útil | criterio S5 exige reporte de cobertura; precios manuales siempre válidos |
| R8 | Lectura asistida de sitios de tiendas: términos de uso y cambios de HTML | bloqueo legal o técnico | revisión documentada por sitio antes de activar; desactivada por defecto; baja frecuencia; nunca masiva |
| R9 | Datos nutricionales en inglés (USDA) y calidad variable (OFF) | empates erróneos | empate asistido con confirmación; fuente citada por ingrediente; captura manual como respaldo |
| R10 | Tope 300K en frontend y API | agentes cortados | ya dividido en agentes disjuntos por carpeta; protocolo de detención y reparto |
| R11 | Agentes no pueden preguntar | suposiciones ocultas | protocolo de bloqueo en disco; revisores verifican que no haya valores inventados (lista de «supuestos» obligatoria en cada reporte de agente) |
| R12 | 16 matices (8 roles + 4 personas + 4 tiendas) | colisiones y daltonismo | paleta categórica con L escalonada, separación ≥ 20°, forma y etiqueta; simulación de daltonismo en U2 |
| R13 | Botón primario con blanco sobre marca 600 no llega a 4.5:1 (calculado 3.7–4.0) | fallo WCAG | regla: relleno marca 700 en claro, marca 300 con texto 950 en oscuro |
| R14 | Telegram con Markdown rompe por guiones bajos | avisos no entregados | texto plano; encabezado en mayúsculas |
| R15 | Ximena debe escribir al bot para registrar su chat | sin avisos para ella | paso en runbook; destinatarios configurables |
| R16 | OCR en tickets térmicos | renglones erróneos | siempre confirmación humana; umbral ≥ 80 % en pruebas |
| R17 | Semilla de gastos fijos depende del archivo de Axel | arranque sin presupuesto completo | importador listo; se carga cuando exista |

---

## 8. Decisiones y acciones de Axel

### 8.1 Decisiones al aprobar este plan
1. **Candidato de matiz:** A Esmeralda (155°), B Teal (195°) o C Índigo (280°) (§4.3).
2. **Pila:** aprobar la Pila A (Go + React) o indicar otra (§6).
3. **Nombre y dirección:** «Finanzas Hogar» y `https://finanzas.local:8443` (propuesta; el puerto 8443 está libre en casa).

### 8.2 Acciones de Axel durante la Etapa B
- Antes de la puerta 1 (o cuando pueda): dejar `~/finanzas-hogar/datos/semilla/gastos_fijos_ingresos.csv` con columnas `tipo (ingreso|gasto_fijo|gasto_venidero), concepto, monto, periodicidad (quincenal|mensual|anual|unica), dia_o_fecha, categoria, persona (axel|ximena|hogar|karen|conjunto), porcentaje_conjunto (opcional), notas`. El plan arranca sin él.
- Puerta 3: iniciar sesión en claude.ai/design en Chrome y compartir la carpeta de entrada con la sesión.
- Operación: instalar la CA de mkcert en los dos Android (pasos en `ops/mkcert.md`); Ximena escribe una vez al bot MultiX para registrar su chat.
- Pendiente externo: confirmar con Kueski si cada pago repone solo capital (cambia el interruptor A/B).

---

## 9. Verificación de extremo a extremo (cómo se comprueba al cerrar la Etapa B)
1. `go test ./...` con cobertura ≥ 85 % en dominio y motores; `oraculo/` ejecuta y coincide al centavo con el motor en 100 % de fixtures, incluida la regresión Kueski-Despensa completa (lista de 1.2).
2. `npm run build` sin errores; Playwright genera las capturas de todas las pantallas en 4 combinaciones; axe-core sin violaciones serias; script de contraste sin pares bajo umbral.
3. Desde un Android en la LAN: abrir `https://finanzas.local:8443` sin advertencia de certificado, instalar la PWA, iniciar sesión, registrar un gasto en ≤ 3 toques, verlo en el dashboard y en la vista de tabla.
4. Disparar manualmente el resumen diario y una alerta de pago: llegan al chat de Axel por el bot MultiX (HTTP 200 y mensaje visible).
5. Generar el reporte mensual PDF y el estado de cuenta de Karen: compilan, abren, y sus cifras coinciden con la API (cruce automático).
6. `systemctl status finanzas-hogar` activo; reiniciar la máquina y comprobar que vuelve solo; ejecutar `ops/respaldo.sh`, restaurar en una ruta temporal y pasar `integrity_check`.
7. Informe final en `docs/estado/informe-final.md` y por Telegram: hecho, verificado con evidencia (rutas, salidas, capturas), consumo real de tokens por fase vs presupuesto, paquetes instalados, y **lo que no se pudo verificar**.

## 10. Primer paso al aprobar
Actualizar las memorias (visión del sistema, estructura del hogar, decisiones), crear `~/finanzas-hogar` con el árbol de §3.2, copiar este plan a `docs/plan/`, y lanzar **WF-0 Cimientos** con Workflow en ultracode (7 agentes, ≈ 850K). Al terminar: Telegram `[Finanzas-Hogar]` y puerta 1.
