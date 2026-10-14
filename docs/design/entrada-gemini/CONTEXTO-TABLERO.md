# CONTEXTO-TABLERO · Finanzas Hogar · todo lo que Gemini necesita leer en un solo archivo

Generado el 14-oct-2026 desde el commit HEAD del repositorio (rama main). **Empieza por `BRIEF-WEB.md`** (el primer archivo de abajo). Cada archivo empieza con una línea `## ARCHIVO: <ruta>` y termina con `## FIN: <ruta>`; si la última línea de este documento no es `## FIN-DEL-CONTEXTO`, el archivo llegó cortado: dilo y no rellenes lo que falte.

Contenido (71 archivos, 649 KB):

- `docs/design/entrada-gemini/BRIEF-WEB.md` (5 KB)
- `docs/design/entrada-gemini/BRIEF-TABLERO.md` (12 KB)
- `docs/design/U3-dashboard.md` (125 KB)
- `docs/design/DESIGN.md` (31 KB)
- `docs/contratos/tokens.css` (26 KB)
- `docs/contratos/extractos/inicio.yaml` (66 KB)
- `docs/contratos/extractos/INDICE.md` (8 KB)
- `frontend/package.json` (3 KB)
- `frontend/src/textos/es-MX.ts` (7 KB)
- `frontend/src/estilos/app.css` (1 KB)
- `frontend/src/mocks/manejadores.ts` (3 KB)
- `frontend/src/app/App.tsx` (2 KB)
- `frontend/src/app/IniciarSesion.tsx` (3 KB)
- `frontend/src/app/enrutador.tsx` (2 KB)
- `frontend/src/app/registro.ts` (1 KB)
- `frontend/src/app/rutas.ts` (5 KB)
- `frontend/src/app/sesion.tsx` (3 KB)
- `frontend/src/app/tema.ts` (1 KB)
- `frontend/src/api/cliente.ts` (4 KB)
- `frontend/src/api/conexion.ts` (1 KB)
- `frontend/src/api/consulta.ts` (3 KB)
- `frontend/src/api/errores.ts` (1 KB)
- `frontend/src/api/formatos.ts` (5 KB)
- `frontend/src/api/lecturas.ts` (1 KB)
- `frontend/src/api/tipos.ts` (2 KB)
- `frontend/src/componentes/Avisos.tsx` (3 KB)
- `frontend/src/componentes/Barometro.tsx` (8 KB)
- `frontend/src/componentes/BarreraErrores.tsx` (1 KB)
- `frontend/src/componentes/Boton.tsx` (2 KB)
- `frontend/src/componentes/Campo.tsx` (2 KB)
- `frontend/src/componentes/Chip.tsx` (3 KB)
- `frontend/src/componentes/Estados.tsx` (4 KB)
- `frontend/src/componentes/Grafica.tsx` (9 KB)
- `frontend/src/componentes/Hoja.tsx` (2 KB)
- `frontend/src/componentes/Icono.tsx` (1 KB)
- `frontend/src/componentes/SelectorArchivo.tsx` (2 KB)
- `frontend/src/componentes/TarjetaMetrica.tsx` (2 KB)
- `frontend/src/componentes/utils.ts` (1 KB)
- `frontend/src/pantallas/inicio/Calendario.tsx` (11 KB)
- `frontend/src/pantallas/inicio/Detalle.tsx` (21 KB)
- `frontend/src/pantallas/inicio/GraficaInicio.tsx` (9 KB)
- `frontend/src/pantallas/inicio/Karen.tsx` (3 KB)
- `frontend/src/pantallas/inicio/P01Inicio.tsx` (2 KB)
- `frontend/src/pantallas/inicio/Resumen.tsx` (12 KB)
- `frontend/src/pantallas/inicio/Tablero.tsx` (35 KB)
- `frontend/src/pantallas/inicio/Tendencias.tsx` (15 KB)
- `frontend/src/pantallas/inicio/comun.tsx` (15 KB)
- `frontend/src/pantallas/inicio/fechas.ts` (3 KB)
- `frontend/src/pantallas/inicio/indice.ts` (1 KB)
- `frontend/src/pantallas/inicio/mocks.ts` (3 KB)
- `frontend/src/pantallas/inicio/textos.ts` (17 KB)
- `frontend/src/pantallas/gastos-presupuesto/P04Gasto.tsx` (14 KB)
- `frontend/src/pantallas/gastos-presupuesto/comun.tsx` (11 KB)
- `frontend/src/pantallas/comidas-compras/P15Productos.tsx` (26 KB)
- `frontend/src/pantallas/comidas-compras/comun.tsx` (11 KB)
- `frontend/pruebas/inicio/datos.ts` (9 KB)
- `frontend/pruebas/inicio/inicio.escritorio.spec.ts` (13 KB)
- `frontend/pruebas/inicio/tablero.escritorio.spec.ts` (23 KB)
- `frontend/pruebas/apoyo/api-falsa.ts` (5 KB)
- `frontend/pruebas/apoyo/banco.ts` (1 KB)
- `frontend/pruebas/apoyo/medidas.ts` (8 KB)
- `frontend/pruebas-integracion/humo-inicio.spec.ts` (5 KB)
- `frontend/playwright.config.ts` (2 KB)
- `docs/estado/bloqueos/WF3c-tablero-interfaz.md` (2 KB)
- `docs/estado/bloqueos/WF3c-barometro-sin-dato.md` (2 KB)
- `docs/estado/bloqueos/U3-tablero-k4-personas.md` (3 KB)
- `docs/estado/bloqueos/U3-tablero-barometro.md` (3 KB)
- `docs/estado/bloqueos/U3-tablero-horizonte.md` (2 KB)
- `docs/estado/bloqueos/U3-tablero-recordatorio.md` (2 KB)
- `docs/estado/bloqueos/U3-tablero-senal-presupuesto.md` (2 KB)
- `docs/estado/bloqueos/C3-ajustes-tablero.md` (3 KB)

## ARCHIVO: docs/design/entrada-gemini/BRIEF-WEB.md


# BRIEF-WEB · Gemini desde el navegador (sin terminal) · Tablero de Inicio · Finanzas Hogar

Versión 1.0 · 14-oct-2026. **Este archivo manda sobre `BRIEF-TABLERO.md` en todo lo que dice cómo trabajar.** De `BRIEF-TABLERO.md` sigues vigentes las secciones 1 a 5: contexto, lecturas, **reglas duras**, mapa de rutas y API y la **lista de lo que falta**. Ignora su sección 0 (clonar, rama, push) y los comandos de su sección 6: en la web no puedes correrlos; los corre Axel.

## 1. Por qué existe este archivo
La primera vez, el repositorio completo llegó a la web como un ZIP que se truncó entre las capturas PNG de `docs/design/capturas-u5/`, antes de llegar a `frontend/` y a este brief. Dos arreglos ya hechos:
- **Rama ligera `gemini-ligero`** del mismo repositorio: solo texto (sin PNG, PDF, fuentes ni bases de datos). URL: `https://github.com/axeljimenezfactible/finanzas-hogar/tree/gemini-ligero`. El repositorio completo sigue en la rama `main`.
- **Archivo único `docs/design/entrada-gemini/CONTEXTO-TABLERO.md`** (también en la rama ligera): trae, uno tras otro y separados por `## ARCHIVO: <ruta>`, todos los archivos que necesitas leer. Si la importación del repo vuelve a truncarse, adjunta solo este archivo.
- Las capturas actuales del Inicio (`P-01-android-claro.png`, `P-01-android-oscuro.png`, `P-01-escritorio-claro.png`, `P-01-escritorio-oscuro.png`) **no** están en la rama ligera; Axel puede adjuntarlas aparte si las pides.

## 2. Lo que no puedes hacer y qué haces en su lugar
No puedes ejecutar `git`, `npm`, Playwright ni compilar. Por eso:
1. **Tu entrega son archivos o parches**, que Axel o Claude aplican y verifican. No la calificas de «terminada» hasta que ellos corran las pruebas.
2. **Nunca declares que algo pasa si no lo corriste.** En tu reporte, cada verificación va como «NO EJECUTADO (revisión estática)» con lo que sí revisaste a ojo: tipos, imports, nombres de props y funciones.
3. **No inventes símbolos.** Usa solo funciones, props, rutas, claves de textos y tokens que aparezcan en los archivos que tienes. Si necesitas leer un archivo que no está en el contexto, no adivines su contenido: pídelo en un bloqueo (`docs/estado/bloqueos/GEMINI-<tema>.md`, formato de `BRIEF-TABLERO.md` regla 9) y sigue con lo que no depende de eso.
4. Si el contexto llega incompleto otra vez (archivo cortado, ZIP truncado), **detente y dilo en la primera línea de tu respuesta**: qué archivo es el último que alcanzaste a leer. No rellenes lo que falta.

## 3. Formato de entrega (obligatorio)
Entrega **por tandas** y espera «continúa» entre una y otra. Cada mensaje trae como máximo 5 archivos y termina con una línea «Tanda k de N · Pendiente: …».
- **Archivo pequeño o nuevo (hasta unas 300 líneas):** un encabezado `### ARCHIVO: ruta/relativa/desde/la/raíz` y debajo el archivo **completo** en un bloque de código con su lenguaje.
- **Archivo grande (por ejemplo `Tablero.tsx`, 422 líneas):** un diff unificado (`diff --git a/ruta b/ruta`, con 3 líneas de contexto) que se pueda aplicar con `git apply`. Dilo en el encabezado: `### PARCHE: ruta`.
- Rutas siempre exactas y relativas a la raíz del repositorio. No cambies nombres de archivos existentes.
- Tandas sugeridas: **T1** punto 1 de la lista (cadena P-04 → P-15 › Tendencias) con su prueba; **T2** puntos 2 a 5; **T3** puntos 6 a 8 y sus pruebas; **T4** la tabla de criterios A-01 a A-22 (punto 9), el reporte `docs/estado/gemini/tablero.md`, `docs/estado/gemini/solicitudes.md` y los bloqueos; el punto 10 (texto) solo si sobra espacio.
- El reporte final (`docs/estado/gemini/tablero.md`) lleva: Hecho, Archivos (rutas), Verificación (todo «NO EJECUTADO» salvo lo que sí revisaste, y cómo), Supuestos, Bloqueos y Pendientes, más el bloque de comandos de abajo.

## 4. Lo que corre Axel después (ponlo tal cual al final de tu reporte)
```
git checkout main && git pull
git checkout -b gemini/tablero
# aplicar cada ARCHIVO (copiar el contenido completo a su ruta) y cada PARCHE:
git apply <archivo.patch>
cd frontend && npm ci && npm run typecheck && npm run build && npm run test:e2e
```
Si algo falla, Axel te pega la salida y tú corriges solo lo que falló, en una tanda nueva.

## 5. Recordatorios que más se olvidan
- La interfaz **no calcula** porcentajes, diferencias ni estados; todo llega de la API.
- Texto de interfaz solo en `frontend/src/pantallas/inicio/textos.ts` o en `frontend/src/textos/es-MX.ts`; español de México, tuteo, sin emojis ni signos de exclamación.
- Colores y tamaños solo con tokens (`docs/contratos/tokens.css`); el color nunca va solo.
- Una tarjeta es un solo enlace: nunca un botón dentro de una tarjeta-enlace.
- No edites `docs/contratos/`, `docs/design/capturas-u5/` ni el backend. Lo que encuentres ahí va a `docs/estado/gemini/solicitudes.md` (ruta, cambio, motivo).

## FIN: docs/design/entrada-gemini/BRIEF-WEB.md

## ARCHIVO: docs/design/entrada-gemini/BRIEF-TABLERO.md


# BRIEF para Gemini · Terminar el tablero de Inicio (P-01) · Finanzas Hogar

Versión 1.1 · 13-oct-2026 · lo escribió Claude para Axel. Lo lee Gemini **con el repositorio completo clonado**: aquí solo están el contexto, las reglas y la lista de lo que falta; el detalle vive en los archivos que se citan. Si algo no está en los archivos, **no lo asumas** (regla 9).

> **Si trabajas desde la web (sin terminal ni git), lee primero `BRIEF-WEB.md`: manda sobre este archivo en todo lo que dice cómo trabajar y entregar.**

## 0. Cómo arrancar (hazlo tú, en orden)
1. `git clone https://github.com/axeljimenezfactible/finanzas-hogar.git && cd finanzas-hogar` (repositorio privado; Axel te da el acceso. Nunca escribas credenciales en ningún archivo).
2. `git checkout -b gemini/tablero`. Trabaja y haz commits **solo en esa rama**; súbela con `git push -u origin gemini/tablero`. **No toques `main`** y no hagas merge: Axel y Claude lo revisan después.
3. `cd frontend && npm ci && npx playwright install chromium` (Node reciente; el backend pide Go 1.27, solo si quieres compilarlo).
4. Para ver el tablero sin backend: `npm run dev:mocks` (API simulada con datos sintéticos de `frontend/src/mocks/`).

## 1. Qué es el proyecto
- «Finanzas Hogar»: sistema web de control financiero del hogar de Axel y Ximena (Karen, madre de Axel, participa sin cuenta en la despensa común y en sus compras con el crédito Kueski). Corre solo en la red de casa, como PWA para Android y escritorio.
- Stack: backend Go + SQLite (`backend/`); frontend React 19 + Vite 8 + TypeScript + Tailwind 4 + Radix + ECharts 6 + PWA (`frontend/`). Contrato único con el backend: `docs/contratos/openapi.yaml` (C1 v1.5.1, 221 operaciones). Tokens de diseño: `docs/contratos/C2-tokens.json` y `tokens.css`.
- Propósito del tablero, palabras de Axel (K54, `docs/plan/decisiones.md` §K): «saber cuánto se está gastando realmente mes con mes y ver los desgloses a detalle, incluso de cuánto costó el kg de cebollas, para ver si aumentamos presupuesto o lo recortamos para irme previniendo». Pidió además ingreso mensual y extra, créditos globales, recordatorio de pagos, un barómetro de capacidad de endeudamiento (verde, amarillo, rojo, con el dígito) y gráficas de análisis.
- Estado: el tablero **ya está construido** (`frontend/src/pantallas/inicio/Tablero.tsx`, 422 líneas: K-1 a K-4, B-1, G-1 a G-6 y «Más del hogar», en el orden de U3 §2.4). Tu trabajo es **cerrar lo que falta** (sección 5), no rehacerlo.

## 2. Lecturas obligatorias (en este orden; todo está en el repositorio)
1. `docs/design/U3-dashboard.md`: §1, §2 completo, §7, §11.2 (criterios A-01 a A-22) y §15. Es la especificación del tablero.
2. `docs/design/DESIGN.md`: §1.2–1.3 (principios y tono), §2.4 (roles de color), §3.2–3.3 (escala y dinero), §4.1–4.3 (espaciado), §7.6 (tarjetas), §7.7 (gráficas), §7.14 (estados), §8.9 (barómetro) y §9 (accesibilidad). `docs/design/referencia-visual.md` (qué adoptar del lienzo de Pista 3).
3. `docs/design/U1-arquitectura.md`: §0.2–0.4, §2.4 y §6. Es grande (160 KB): léelo por partes.
4. `docs/estado/bloqueos/`: `WF3c-tablero-interfaz.md`, `WF3c-barometro-sin-dato.md`, `U3-tablero-k4-personas.md`, `U3-tablero-barometro.md`, `U3-tablero-horizonte.md`, `U3-tablero-recordatorio.md`, `U3-tablero-senal-presupuesto.md`, `C3-ajustes-tablero.md`. Traen las decisiones PROVISIONALES vigentes: respétalas.
5. Código: `frontend/src/pantallas/inicio/` (Tablero, Resumen, Detalle, Calendario, comun.tsx, textos.ts), `frontend/src/componentes/` (TarjetaMetrica, Barometro, Grafica, Estados, Boton), `frontend/src/app/rutas.ts` y `enrutador.tsx`.
6. Pruebas existentes: `frontend/pruebas/inicio/tablero.escritorio.spec.ts` y `inicio.escritorio.spec.ts`, más `frontend/pruebas-integracion/humo-inicio.spec.ts`.
7. Evidencia visual actual: `docs/design/capturas-u5/P-01-android-claro.png`, `P-01-android-oscuro.png`, `P-01-escritorio-claro.png`, `P-01-escritorio-oscuro.png`.
8. Reglas de negocio del tablero: `docs/contratos/reglas-negocio.md` (RN-30 barómetro, RN-31 sugerencia de presupuesto) y `docs/contratos/C1-datos.md` §9, métricas M-18 a M-22 (usa `grep`; el archivo es de 180 KB).

## 3. Reglas duras
1. **La interfaz no calcula** porcentajes, diferencias ni estados: llegan del servidor (U3 §2, U1 §0.2). Ningún umbral es constante en el código (30 y 40 llegan en `warning_pct` y `limit_pct`).
2. Dinero en centavos enteros; se muestra `$1,258.68` con cifras tabulares y menos tipográfico U+2212. Fechas `10-oct-2026`. Zona America/Mexico_City.
3. Texto de interfaz en español de México, tuteo, sin jerga, sin signos de exclamación ni emojis. Todo texto nuevo va en `frontend/src/pantallas/inicio/textos.ts` o en `frontend/src/textos/es-MX`, nunca incrustado.
4. Colores, tamaños y espaciados **solo con tokens** (`ui.*`, `color.role.*`, `font.*`, `space.*`). El color nunca va solo: lleva icono, signo o texto. Claro y oscuro. Contraste AA calculado.
5. Una sola acción primaria por pantalla («Registrar gasto»). Cada tarjeta es un solo enlace; un botón dentro de una tarjeta-enlace es HTML inválido y rompe A-08. Objetivos táctiles ≥ 44 px. Móvil 412 px primero; revisa 360, 412, 768, 1024 y 1440 sin desplazamiento horizontal.
6. Cada bloque hace su propia petición y tiene sus estados: carga, vacío, error con «Reintentar» y sin conexión. Si un bloque falla, los demás se ven normales.
7. Sin CDN ni fuentes de terceros; sin dependencias nuevas salvo que lo justifiques por escrito.
8. **No edites** (son generados o están firmados): `docs/contratos/` (C1, C2, C3, `openapi.yaml`, `extractos/`), `frontend/src/api/esquema.gen.ts`, `docs/design/capturas-u5/`, `docs/design/sistema/`, `docs/design/maquetas/`, `docs/estado/wfp3/`. **Backend: solo lectura.** Si encuentras un hueco en el backend o en el contrato, escríbelo en `docs/estado/gemini/solicitudes.md` (ruta, cambio, motivo) y sigue.
9. **No asumas.** Si te falta un dato o un criterio admite dos lecturas, crea `docs/estado/bloqueos/GEMINI-<tema>.md` con: Contexto, Qué falta (la pregunta exacta), Opciones (2 a 4, con su consecuencia), Impacto y Cómo seguí. Continúa con lo que no depende de eso y marca lo hecho como PROVISIONAL.

## 4. Mapa de pantallas y de API que usa el tablero
Rutas del frontend (siempre con `rutaDe('P-xx', params)` de `frontend/src/app/rutas.ts`, nunca a mano):
- Inicio P-01 `/inicio` (pestañas Resumen, Detalle, Tendencias, Calendario).
- Gastos: P-02 `/gastos`, P-04 `/gastos/:id`, P-05 `/gastos/ticket`.
- Presupuesto: P-06 `/presupuesto`, P-07 `/presupuesto/flujo`, P-08 `/presupuesto/ingresos`, P-09 `/presupuesto/ahorro`.
- Compras: P-13 `/compras`, P-14 `/compras/comparar`, P-15 `/compras/productos`, P-16 `/compras/inventario`.
- Créditos: P-17 `/creditos`, P-18 `/creditos/:id`, P-20 `/creditos/comparador`. Ajustes P-26 `/ajustes/:seccion` (parámetros del tablero en Ajustes › Presupuesto, P-26.02).

Operaciones de C1 por bloque (rutas bajo `/api/v1/`; tipos en `frontend/src/api/esquema.gen.ts`; extractos por módulo en `docs/contratos/extractos/`, con `INDICE.md` que dice cuál archivo trae cada operación):
- K-1 `dashboard/month-vs-plan` · K-2 `incomes/month` · K-3 `credits/totals` y `dashboard/next-payment` · K-4 `scheduled-payments`.
- B-1 `dashboard/debt-capacity` · G-1 y G-6 `dashboard/monthly-overview?from=&to=` (una sola petición de 17 meses).
- G-2 `debt/live` · G-3 `budget/categories` y `budget/signals` · G-4 `cash-flow/accrual-vs-cash` · G-5 `credits/debt-trend`.
- Parámetros: `GET` y `PATCH settings/dashboard`; umbrales de presupuesto `GET settings/alerts` (`getAlertSettings`).

## 5. Lo que falta, en orden de valor
1. **Desglose en cadena (U3 §2.2, criterio A-21).** Hoy no existe el enlace de P-04 a P-15. Falta: (a) cada renglón del ticket en `frontend/src/pantallas/gastos-presupuesto/P04Gasto.tsx` enlaza a P-15 con su producto; (b) `frontend/src/pantallas/comidas-compras/P15Productos.tsx` abre en «Tendencias» con ese producto elegido (hoy el selector `p15-historial` es estado local: pásalo a un parámetro de la URL, siguiendo las convenciones de `enrutador.tsx`); (c) la cadena de 4 toques desde Inicio: barra de G-3 → P-02 con `category_id` y `period` → P-04 → P-15 › Tendencias. Agrega su prueba A-21.
2. **Calendario.tsx:** que no pinte «Hogar $973.75» sin la parte de Karen mientras `listPersons` no responde: esqueleto y, si falla, error con «Reintentar». Mismo patrón ya aplicado en K-4 y en D-3 de `Detalle.tsx` (bloqueo `U3-tablero-k4-personas`). Con prueba.
3. **`tableroNoCargo`** debe usar el estado común «Error al cargar · Reintentar» (`componentes/Estados.tsx`) y **`Estados.tsx`** debe revisar su `role="status"` (hallazgos B-3 y N-3).
4. **K-3 sin «Reintentar»** cuando solo falla `getDashboardNextPayment` (bloqueo `WF3c-tablero-interfaz`, punto 2). Hoy muestra «Disponible hoy —» sin botón. Evalúa la opción B (botón fuera de la tarjeta-enlace) sin romper A-08. Si la adoptas, documéntalo como PROVISIONAL.
5. **K-1:** la insignia «[n] categorías rebasadas» es solo texto (punto 3 del mismo bloqueo). Evalúa un enlace aparte a G-3 fuera de la tarjeta. **K-4:** «Hay más pagos en Calendario» sale sin número porque la API no lo da; **déjalo así** (cambiar el contrato no es tuyo; ya hay solicitud).
6. **`humo-inicio.spec.ts`** (`frontend/pruebas-integracion/`) espera las 6 tarjetas de U3 1.x y no entra en ningún `verify-*`. Actualízalo a K-1 a K-4, B-1 y G-1 a G-6 y compáralo con la API.
7. **«Más del hogar»** (tarjetas 2, 4, 5 y 6 de U3 1.x): confirma que siguen sin cambio y que la semana de comidas no muestra costo ni kcal en 0 como si fueran cifras.
8. **Vacíos con sentido.** Sin ingresos capturados, K-2, B-1, G-1 y G-6 deben decir qué falta y ofrecer «Capturar ingreso» (P-08); nunca un 0 inventado. Ojo: en la base de desarrollo de Axel no hay ingresos, así que esos cuatro bloques salen vacíos **por datos, no por error de código**. Para verlos llenos usa los mocks y las respuestas sintéticas de `tablero.escritorio.spec.ts` (ayudante `sinteticas()`).
9. **Criterios A-01 a A-22 de U3 §11.2:** corre cada uno y entrega una tabla con su estado (cumple, parcial, no cumple) y la evidencia.
10. **Al final y solo si sobra tiempo (texto, sin código):** `docs/design/DESIGN.md` P1 (línea 17) todavía dice «≤ 6 tarjetas del Resumen», pero K54 retiró el tope; y las 3 líneas con +$138.50 (`docs/design/U1-arquitectura.md:652` y `:978`, `docs/design/U3-dashboard.md:317`) deben decir +$264.50 (K46). Si dudas, pregunta con un bloqueo.

Fuera de tu alcance (ya lo tiene Claude): alertas de presupuesto 80 %/100 %, precios viejos y día de despensa (`backend/internal/reportes/avisos/fuente.go`), `deriveCredits` (`backend/internal/api/gastos-inicio/vistas.go:235`), la fecha por omisión de `GET /debt/live` y el paso a producción (`ops/`).

## 6. Verificación y definición de terminado
- Desde `frontend/`: `npm run typecheck`, `npm run build` y `npm run test:e2e` (Playwright con API simulada; tarda unos 6 minutos). `npm run verify` corre todo junto (build, rutas contra el contrato, sin hosts externos, manifiesto y pruebas). Debe quedar en verde, con 0 violaciones serias de axe y 0 pares de contraste bajo umbral.
- **No regeneres** las capturas de `docs/design/capturas-u5/` ni la matriz C3 (las regenera el orquestador). Pero deja anotado en tu reporte qué pantallas cambiaste, porque sus capturas quedan viejas.
- Reporte en `docs/estado/gemini/tablero.md` con: Hecho, Archivos (rutas), Verificación (comando y salida resumida), Supuestos, Bloqueos (rutas) y Pendientes.
- Commits pequeños en `gemini/tablero`, con mensajes que empiecen con «[Gemini · Tablero]», y rama subida a `origin`. Sin merge a `main`.

## FIN: docs/design/entrada-gemini/BRIEF-TABLERO.md

## ARCHIVO: docs/design/U3-dashboard.md


# U3 · Dashboard: tablero de control (P-01 Inicio) · Finanzas Hogar

Versión 2.0 (12-oct-2026) · Agente `u3-v2` (WF-3c, decisión K54 de Axel, 12-oct-2026): **«Tablero de control»**. La pestaña «Resumen» deja las 6 tarjetas y pasa a ser un tablero: el gasto real del mes como indicador principal (K-1), ingreso del mes con su extra (K-2), créditos globales (K-3), pagos próximos (K-4), el barómetro de capacidad de endeudamiento (B-1, DESIGN §8.9) y seis gráficas de análisis (G-1 a G-6), con el desglose en cadena hasta el historial del precio por unidad base (§2.2) y la tabla «Qué miro para prevenir» (§2.3). Detalle (§3), Tendencias (§4) y Calendario (§5) se conservan; G-4 y G-5 son T-2 y T-3. Las tarjetas 2, 4, 5 y 6 siguen sin cambio en «Más del hogar» (§2.5, PROVISIONAL); la 1 pasa a K-1 y la 3 a K-3 y K-4.

**Cambios de la versión 2.0:**
- §2 nuevo (sustituye «Resumen: 6 tarjetas»): K-1 a K-4, B-1, §2.1 y G-1 a G-6, §2.2 desglose en cadena, §2.3 «Qué miro para prevenir», §2.4 orden en 412 y rejilla en 1440, §2.5 «Más del hogar», §2.6 operaciones de C1 (3 nuevas, un par de Ajustes, 14 campos nuevos en 3 esquemas y 5 claves de `app_setting`) para C1 v1.5.
- PROVISIONAL con bloqueo para la puerta 4: ingreso neto, pagos que cuentan, redondeo y escala del barómetro (`docs/estado/bloqueos/U3-tablero-barometro.md`); N y rango de K-4 (`U3-tablero-recordatorio.md`); horizonte de las gráficas (`U3-tablero-horizonte.md`); regla «subir o recortar presupuesto» (`U3-tablero-senal-presupuesto.md`); tarjetas que se conservan (`U3-tablero-mas-del-hogar.md`).
- §1, §3, §6 a §15: las referencias a las tarjetas 1 y 3 pasan a K-1, K-3 y K-4; filtros, estados, tamaños, Karen, accesibilidad, criterios A-xx, capturas, trazabilidad, supuestos y dependencias al día. Ninguna cifra canon cambia; las cifras nuevas de los ejemplos salen de `oraculo/model.json` (`periodos[0..2]`) o son marcadores.

- Corrección r1 (12-oct-2026, agente `correccion-especificacion-r1`, hallazgo M-1 de `docs/estado/wf3/REVISION-especificacion-r1.md`): §2.6 queda alineada con C1 v1.5 tal como se publicó (ventana inválida de `getDashboardMonthlyOverview`: 422 `validation_failed`; el PATCH de `updateDashboardSettings` lleva `parameters` y `reset_parameters`); las operaciones, los 14 campos y M-18 a M-22 pasan de PROPUESTA a publicados en C1 v1.5; K-4 lee N de `getDashboardSettings` una sola vez (§1.2) y la interfaz no fija N; «Fuentes que rigen» cita C1 v1.5. Ninguna cifra cambia.
- Corrección r2 (12-oct-2026, agente `correccion-especificacion-r2`, hallazgo M-N1 de `docs/estado/wf3/REVISION-especificacion-r2.md`): K-1 cita la fila M-01 de C1-datos §9, que ya define los 8 campos nuevos de `DashboardMonthVsPlan` con su fórmula, periodo, tablas y ejemplo del oráculo (598.82 de efectivo en oct-2026; en nov-2026, 0.00 de diferencia en lo gastado y 598.82 más en efectivo). Las cifras y la fórmula de K-1 no cambian.

Versión 1.5 (11-oct-2026) · Agente `correccion-U4-r2` (P2c, hallazgos M-1 y M-2 de `docs/estado/p2c/REVISION-U4-r2.md`): (M-1) en D-2 el ahorro planificado deja de ser FIRME: el reparto mitad y mitad sí lo es (Q2-13 B), pero el valor es **PROVISIONAL** por el bloqueo abierto `docs/estado/bloqueos/S2-ahorro-quincenal.md`; mientras llegue `null`, la fila muestra «—»; §15 lo pasa a «Abiertos» y lo quita de los resueltos. (M-2) D-6 muestra el avance con un decimal (`progress_tenths_pct`), porque Q2-17 C deja el entero solo en la tarjeta de Inicio y D-6 es la métrica M-16 «Metas»; §15 pide a C1-datos alinear M-16. Ninguna cifra de canon cambia. Versión 1.4 (11-oct-2026) · Agente `correccion-U4-r1` (P2c, hallazgo M-4 de `docs/estado/p2c/REVISION-U4-r1.md`): se quitan las marcas PROVISIONAL de bloqueos que la puerta 2 dejó RESUELTOS, `S2-presupuesto-total` (Q2-19 B y Q2-20 A), `S2-metas` (Q2-16 A, Q2-17 C y Q2-18 A), `S2-proyeccion-interruptor` (Q2-21 a Q2-24 C) y Q2-13 B y Q2-14 A de `S2-flujo-quincenal`; siguen PROVISIONAL, citando sus bloqueos abiertos, el gasto variable estimado (`S2-flujo-quincenal`, Q2-12) y la aportación de una meta vencida sin cumplir (`S2-metas-vencida`). La tarjeta 6 y D-6 usan los tres estados de Q2-18 A. Versión 1.3 (11-oct-2026) · Agente `U4-puerta2` (P2c, puerta 2): Q2-65 A (T-2 con los acumulados de C1, una sola variante; bloqueo `docs/estado/bloqueos/U3-t2-acumulados.md` RESUELTO) y Q2-86 A (la tarjeta 5 nombra la segunda tienda con `second_store_id`, solicitud a C1 aprobada); ninguna cifra cambia. Versión 1.2 (11-oct-2026) · Agente U3-dashboard (WF-1); ronda 1 de corrección: agente `correccion-U3-r1`; ronda 2: agente `correccion-U3-r2` · Evalúa: `docs/criterios/U3.md` v1.2.

**Cambios de la versión 1.2** (hallazgo alto H-1 y medios H-2 y H-3 de `docs/estado/wf1/REVISION-U3-r2.md`):
- H-1: en T-2 las dos capas tienen el mismo origen que `cumulative_difference_cents` (C1 M-11): capa 1 = efectivo acumulado, cima = devengado acumulado. C1 v1.2 no trae esos dos acumulados; U3 los pide (§15) y deja las capas, sus etiquetas finales y su resumen **PROVISIONAL** con dos variantes (`docs/estado/bloqueos/U3-t2-acumulados.md`). Queda prohibido sumar `cash_cents` desde el primer mes visible y apilar encima `cumulative_difference_cents`. A-11 suma una corrida con reloj en nov-2026.
- H-2: los ejemplos de oct-2026 a mar-2027 son los periodos P1 a P6 del oráculo (C1-datos M-01; `is_canonical`; `S2.md` S2-09, nota de equivalencia). Se retira la frase de T-2 que lo negaba; R-51 y R-52 comprueban las capas en mar-2027.
- H-3: cada gráfica lleva su resumen de una frase y su `<figcaption>` (DESIGN §7.7), con marcadores y solo cifras de C1; sin variaciones calculadas y, en T-2, sin rotular la deuda viva.

**Cambios de la versión 1.1** (hallazgos medios H-1 a H-7 de `docs/estado/wf1/REVISION-U3-r1.md`):
- H-1: la tarjeta sigue DESIGN §7.6 (etiqueta `font.size.small` en `ui.text.secondary`, relleno `space.6` desde 768 px); las gráficas cargan sin ejes falsos (§7.7); nuevo §1.4 con los estados de pago de DESIGN §7.8, incluidos «Por vencer» y «Parcial» (K15); el destino de cada tarjeta se mantiene y se pide a DESIGN (§15).
- H-2: §4 adopta DESIGN §8.1, §8.3, §8.5 y §8.7 (sin contorno en T-1; en T-2, capa 2 en cero con «Pagado por adelantado» y línea «Hoy»; etiquetado selectivo; Ximena con cuadrado y trazo 8 · 4); se retira el supuesto 5.
- H-3 y H-4: D-3 hace dos peticiones explícitas (vencidos, y de hoy a hoy + 60) y su ejemplo trae los 3 pagos que caben en el rango.
- H-5: la tabla de T-2 entra en §9 y en A-13.
- H-6: la fila «Total» de T-2 suma las 12 filas visibles, sin cifra de canon.
- H-7: D-4 y D-5 escriben su fórmula, periodo y filtros mientras C1-datos §9 no tenga su métrica (solicitud en §15).

**Fuentes que rigen** (no se reabren): decisiones §A, §B, §C, §G y §K (K5, K8, K38, K41, K43, K44 y, desde 2.0, K54); U1 1.2 §0.2 (vocabulario, semana por defecto K38), §0.4 (formatos), §0.5 (personas), §1.3 (móvil), §2.4 (estados), §3.3 (P-01: N-020 a N-029), §3.9 (P-07), §6.1 y §6.2 (puntos de quiebre); contrato C1 v1.5 por sus extractos `inicio.yaml`, `presupuesto-a.yaml`, `presupuesto-b.yaml` (y, para los paneles que reutilizan necesidades de otras pantallas, `comidas.yaml`, `compras-tiendas-a.yaml`, `creditos-a.yaml`, `creditos-b.yaml`, `ajustes-b.yaml`, `reportes.yaml`); métricas M-01 a M-17 de `C1-datos.md` §9; DESIGN.md §1, §2.4, §2.6 a §2.10, §3, §4, §5.2 a §5.4, §7.6 (tarjetas de métrica), §7.7 (gráficas), §7.8 (estados de pago) y §8 (patrones de visualización). Donde U3 se aparta de DESIGN lo dice en el lugar, con su razón, y lo pide a DESIGN en §15.

**Reglas que este documento no repite en cada bloque:**
- La interfaz **no recalcula reglas de negocio**: muestra, filtra, ordena y suma lo que está a la vista (U1 §0.2, «Motor Sx»). Toda cifra viene de la operación C1 citada.
- **Sin datos no es cero**: si falta una entrada se muestra «—» con el motivo, nunca un monto en cero ni `0 %` (U1 §0.2).
- **Toda cifra de presupuesto, flujo y ahorro es la parte del hogar**; la parte de Karen es informativa y solo aparece donde ella participa (U1 §0.2, §0.5; decisiones §A).
- Dinero con el formato de U1 §0.4 y DESIGN §3.3: MXN, símbolo `$`, coma de miles, dos decimales, cifras tabulares (`tnum`), negativo con el signo menos tipográfico, color de rol **e** icono (nunca solo color). Fechas en pantalla `10-oct-2026`; con día de la semana donde importa (`Vie 30 oct 2026`). En la API las fechas van en ISO; en este documento solo aparecen en formato de pantalla.
- Colores, tipografía, espacios, radios y sombras **solo por nombre de token** de C2 / DESIGN (0 valores sueltos).
- Ejemplos con cifras: solo las del oráculo o la semilla con su referencia (R-nn de `S2.md`/`S3.md`, K-nn de `S5.md`); lo demás va como marcador entre corchetes, sin número: `[plan]`, `[n]`. Los meses oct-2026 a mar-2027 de los ejemplos son los periodos P1 a P6 del oráculo: C1-datos M-01 dice «P1 = oct-2026», `is_canonical` de C1 marca «oct-2026 a mar-2027» y `docs/criterios/S2.md` S2-09 (nota de equivalencia) explica que con la semilla los dos dan las mismas cifras, porque ningún evento cae entre los días 1 y 9 de un mes. Por eso «octubre de 2026» cita `periodos[0]` en K-1, D-7, T-2, A-06 y A-11 (y `periodos[1]` y `periodos[2]`, noviembre y diciembre, en K-1, B-1, G-1 y G-6). Una prueba con eventos entre los días 1 y 9 toma su esperado de las ventanas que declara, nunca de R-45 a R-50 (S2-09).

---

## 1. Estructura de Inicio

### 1.1 Pestañas y página
- **Título de página**: «Inicio» (`font.size.title-1`, `font.weight.bold`), icono `house`.
- **Pestañas** (orden fijo, decisiones §G): «Resumen» · «Detalle» · «Tendencias» · «Calendario». Abre en «Resumen», que desde U3 2.0 es el tablero de control de §2 (K54; la etiqueta de la pestaña no cambia, decisiones §G). La pestaña elegida se recuerda en el dispositivo (frecuencia «Local», U1 §0.3). En < 768 px la fila de pestañas se desplaza en horizontal y la activa queda siempre visible (U1 §1.3); en ≥ 768 px es una fila completa.
- **Acción primaria de la página**: «Registrar gasto» (U1 §6.4, fila P-01). Es el **único** botón con relleno de marca de la página (DESIGN P2, §2.9): en < 768 px es el «+» de la barra inferior (icono `plus`); en ≥ 768 px, el botón «+ Registrar gasto» de la barra superior. Abre P-03 (flujo F-1, 3 toques, K10).
- Las tarjetas **no** llevan botón con relleno: su acción primaria única es un **enlace** (texto subrayado con `color.brand.700` / `color.brand.300` e icono `chevron-right`) y toda la tarjeta es zona de toque que lleva al mismo destino.

### 1.2 Anatomía común de una tarjeta de métrica
Componente «Tarjeta de métrica» (DESIGN §7.6). Partes, de arriba abajo:

| Parte | Contenido | Tokens |
|---|---|---|
| Etiqueta (título) | Icono de 20 px (`size.icon.sm`, DESIGN §5.2; el mapa de §5.4 prevé iconos en tarjetas) + título en español, sin dos puntos | etiqueta `font.size.small` en `ui.text.secondary` (DESIGN §7.6); el icono toma el color del texto que acompaña (DESIGN §5.3) |
| Valor principal | Una cifra o texto corto; nunca se trunca ni se parte (U1 §6.1 regla 5) | `font.size.title-2`, `font.weight.semibold`, `tnum`; neutro `color.neutral.900` / `color.neutral.100` salvo cifra de rol (DESIGN §2.9) |
| Valor secundario | Una línea a 360 px; una línea más desde 412 px (U1 §6.2) | `font.size.small`, `color.neutral.700` / `color.neutral.300` |
| Insignias de estado | Icono + texto (nunca solo color) | insignia de rol: texto `color.role.<rol>.light` sobre `color.role.<rol>.surface-light` (oscuro: `.dark` sobre `.surface-dark`) |
| Pie | Enlace de la acción primaria | `font.size.small`, `font.weight.medium`, alto de zona `size.touch.min` |

- Superficie: `color.neutral.100` sobre fondo `color.neutral.50` (oscuro: `color.neutral.900` sobre `color.neutral.950`), `shadow.0`, borde `border.width.default`, radio `radius.md`, relleno interior `space.4` (< 768 px) y `space.6` (≥ 768 px), como DESIGN §7.6. Separación entre tarjetas: el medianil de DESIGN §4.2, `space.4` (< 1024 px) y `space.6` (≥ 1024 px).
- **Destino**: la tarjeta entera es un enlace a la pantalla destino de su acción primaria (K-1 → P-06, K-2 → P-08, K-3 y B-1 → P-17, K-4 → pestaña Calendario con cada fila a P-18; en «Más del hogar», P-07, P-10, P-14 y P-09; §2), no a la pestaña Detalle; Detalle se abre con su pestaña. U3 se aparta aquí de DESIGN §7.6 («la tarjeta entera es enlace al Detalle») porque U3-01 pide una acción primaria única por tarjeta y exige que la de la tarjeta 4 lleve a P-10; solicitud 1 a DESIGN (§15).
- Toda la tarjeta mide al menos `size.touch.min` de alto y es un solo objetivo táctil con foco visible (`border.width.focus`, anillo `color.brand.600` / `color.brand.400`). Nombre accesible = título + valor principal + estado.
- Cada tarjeta hace **su propia petición** a C1: si una falla, las otras se ven normales.
- Ninguna tarjeta usa los umbrales 80 % y 100 % como constantes: el estado ya viene calculado en la respuesta (`status` de M-01). Si una implementación necesitara el umbral, lo lee de `getAlertSettings` (C1-22, K44).
- Los días antes del pago (`alerts.payment_due_days_before`, inicial 1 y 0, decisión F y K44) sí se usan, para el estado «Por vencer» de K-4, D-3 y el Calendario (§1.4). Inicio los lee **una sola vez** con `getAlertSettings` (GET `/api/v1/settings/alerts`, campo `parameters.alerts.payment_due_days_before.value`) y nunca como constante; si esa lectura falla, ningún pago se muestra «Por vencer» (queda «Pendiente») y nada más cambia.
- N de «Pagos próximos» (K-4, K54) tampoco es constante: Inicio lo lee **una sola vez** con `getDashboardSettings` (GET `/api/v1/settings/dashboard`, C1 v1.5, campo `parameters["dashboard.payment_reminder_count"].value`, valor inicial 5), igual que `getAlertSettings`. Si esa lectura falla, K-4 muestra su estado de error (§7) y los demás bloques siguen; el valor 5 y esta conducta son PROVISIONALES (bloqueo `docs/estado/bloqueos/U3-tablero-recordatorio.md`).

### 1.3 Estados comunes de las tarjetas (U1 §2.4)
- **Carga**: esqueleto estático con la forma final (etiqueta, bloque de cifra y una o dos líneas) en `ui.skeleton.on-card` (DESIGN §7.6; movimiento de DESIGN §6.2); nunca cifras provisionales. Si tarda más de 10 s: «Esto está tardando más de lo normal · Reintentar».
- **Error**: mensaje de la tarjeta (texto exacto en cada una) + botón secundario «Reintentar» (borde `color.neutral.600` / `color.neutral.400`), icono `circle-alert` con `color.role.error.light` / `color.role.error.dark`. La tarjeta conserva su título y su enlace.
- **Sin conexión** (G-07, K9): franja global «Sin conexión con el servidor de casa» (icono `wifi-off`; «· [n] en cola» si hay capturas guardadas en el teléfono). Cada tarjeta muestra lo último cargado con la línea «Cargado a las [HH:MM]» (`font.size.caption`) y su enlace sigue abriendo la pantalla destino en solo lectura. Si la tarjeta nunca se cargó en el teléfono: «Sin conexión · Esta página aún no se ha cargado en este teléfono». «Registrar gasto» sigue disponible y deja el gasto en la cola (U1 §2.4).

### 1.4 Estados de un pago (DESIGN §7.8)
K-4, D-3 y el Calendario usan los estados de pago de DESIGN §7.8, siempre con icono y texto (nunca solo color). Un pago muestra un solo estado: el primero de esta tabla que se cumpla.

| Estado | Cuándo (dato de C1) | Icono y color | Texto |
|---|---|---|---|
| Pagado | `status` = `paid` (todas las partes pagadas) | `circle-check`, `color.role.success.*` | «Pagado» |
| Vencido | `status` = `overdue` (en `getDashboardNextPayment`, `is_overdue` = true) | `calendar-x`, `color.role.warning.*` | «Vencido» y la línea «Sin marcar como pagado» (U1 N-022, C1 M-03) |
| Por vencer | pendiente, y los días que faltan para el pago son uno de los de `alerts.payment_due_days_before` (§1.2) | `calendar-clock`, `color.role.warning.*` | «Vence mañana» (1), «Vence hoy» (0); otro día configurado en Ajustes: «Vence en [d] días» |
| Parcial | pendiente, con una parte de `shares[]` pagada por completo (`paid_cents` = `amount_cents`) y otra no; el capital se repone solo cuando el pago se completa (K15) | `circle-dot`, `ui.icon.default` | «Falta la parte de Karen» (o «Falta la parte del hogar» si la que falta es la del hogar) |
| Pendiente | cualquier otro pago con alguna parte sin pagar | `circle-dashed`, `ui.icon.default` | «Pendiente» |

- Un pago «Vencido» o «Por vencer» que ya tiene una parte pagada agrega en su línea secundaria «Falta la parte de Karen» (o «del hogar»).
- Qué estados puede mostrar cada bloque: K-4 y D-3 traen `shares[]` y solo pagos con alguna parte sin pagar, así que nunca muestran «Pagado» (sí «Parcial»); el Calendario muestra los cinco (§5).
- «Sin marcar como pagado» acompaña a «Vencido» porque U1 N-022 y C1 M-03 piden esa marca y DESIGN §7.8 solo pone «Vencido». El orden de prioridad y los textos «Falta la parte del hogar» y «Vence en [d] días» tampoco están en DESIGN (supuestos 6 y 7; solicitud 2 a DESIGN, §15).

---

## 2. Resumen: tablero de control (K54)

**Por qué cambia** (decisiones K54, 12-oct-2026): Axel pidió «un dashboard de un sistema de control completo» con ingreso mensual y extra, créditos globales, recordatorio de pagos, un barómetro de capacidad de endeudamiento «con el digito de cuanto se lleva» y «graficas de analisis». Y fijó el propósito del sistema: «saber cuanto se esta gastando realmente mes con mes y ver los desgloces a detalle incluso de cuando costo el kg de cebollas, para ver si aumentamos presupuesto o lo recortamos». Por eso la pestaña «Resumen» deja las 6 tarjetas de U3 1.x y pasa a ser un **tablero de control**. Se retira la regla «≤ 6 tarjetas» (K54 decisión 3). Las pestañas, su orden y la acción primaria de página no cambian (§1.1).

**Indicador principal: el gasto real del mes** (K-1), devengado y en efectivo, contra el plan y contra el mes anterior, con desglose en cadena hasta el precio por unidad base de un producto (§2.2) y la sugerencia «subir o recortar presupuesto» por categoría (G-3).

**Bloques del tablero** (ids fijos; el orden de apilado y la rejilla están en §2.4):

| Id | Bloque | Pregunta que contesta | Operación C1 |
|---|---|---|---|
| K-1 | Gasto del mes (principal) | ¿Cuánto gastamos realmente y cómo vamos contra el plan y contra el mes pasado? | `getDashboardMonthVsPlan` (M-01) + 8 campos nuevos de C1 v1.5 |
| K-2 | Ingreso del mes | ¿Cuánto entró este mes y hubo extra? | `getIncomeMonth` (N-100) + 4 campos nuevos de C1 v1.5 |
| K-3 | Créditos | ¿A cuánto ascienden los créditos y cuánto toca pagar este mes? | `getCreditTotals` (M-12 a hoy) + 2 campos nuevos; `getDashboardNextPayment` (M-04) |
| K-4 | Pagos próximos | ¿Qué pagos vienen y cuál está por vencer o vencido? | `listScheduledPayments` (M-13) |
| B-1 | Barómetro de capacidad de endeudamiento | ¿Qué parte del ingreso se va en pagos de crédito? | **NUEVA** `getDashboardDebtCapacity` |
| G-1 | Ingreso, gasto y pagos por mes | ¿Entra más de lo que gastamos y pagamos? | **NUEVA** `getDashboardMonthlyOverview` |
| G-2 | Deuda viva por línea | ¿En qué línea está la deuda? | `getLiveDebt` (M-12) |
| G-3 | Presupuesto por categoría | ¿Qué categoría se pasa y cuál conviene subir o recortar? | `getBudgetByCategory` (M-09) + **NUEVA** `getBudgetSignals` |
| G-4 | Devengado contra efectivo | ¿Lo que se consume se está pagando? | `getAccrualVsCash` (M-11) = T-2 de §4 |
| G-5 | Deuda viva proyectada | ¿Cuándo baja la deuda? | `getDebtTrend` (M-17) = T-3 de §4 |
| G-6 | Tendencia del barómetro | ¿La capacidad va a empeorar en los próximos meses? | **NUEVA** `getDashboardMonthlyOverview` (la misma petición de G-1) |
| 2.5 | Más del hogar | Quincena, comidas, canasta y ahorro (tarjetas 2, 4, 5 y 6 de U3 1.x sin cambio) | las de cada tarjeta; **PROVISIONAL** |

**Reglas del tablero** (además de las del inicio del documento y de §1):
- Cada bloque hace **su propia petición** y tiene sus estados (§1.3, §7): si uno falla, los demás se ven normales. G-1 y G-6 comparten una petición; si falla, los dos muestran su error.
- La interfaz **no calcula** porcentajes, diferencias ni estados: llegan de C1 (motores S2 y S3). Solo cuenta, ordena y suma lo que está a la vista (U1 §0.2).
- Umbrales: los del presupuesto (80 % y 100 %) se leen una vez de `getAlertSettings` (§1.2). Los del barómetro (30 % y 40 %) llegan en la respuesta de `getDashboardDebtCapacity` y de `getDashboardMonthlyOverview` (`warning_pct`, `limit_pct`, de las claves de `app_setting` `debt.capacity_warning_pct` y `debt.capacity_limit_pct`, K54 y K44). Ningún umbral es constante en el código.
- K-1 a K-4 usan la «Tarjeta de métrica» (§1.2; K-4 con una lista compacta). B-1 usa el barómetro de DESIGN §8.9. Las gráficas G-1 a G-6 siguen todas las reglas de §4 (orden fijo de series, forma y patrón, etiquetas de dinero en texto, `<figure>` con `<figcaption>`, resumen sin variaciones calculadas, «Ver como tabla»).
- Los campos y operaciones **NUEVOS** de §2.6 están **publicados en C1 v1.5** (agente `s1-c1-v15`). Si un servidor anterior a C1 v1.5 responde sin un campo nuevo, la línea que lo usa no se pinta y el bloque muestra lo demás; un bloque que depende solo de una operación nueva muestra «Disponible cuando el servidor se actualice» (`font.size.caption`, `info`), nunca una cifra inventada.

### K-1 · Gasto del mes (indicador principal)
- **Métrica C1**: M-01 (N-020), ampliada en C1 v1.5 para K-1: la fila M-01 de C1-datos §9 trae la fórmula, el periodo (mes y mes calendario anterior), las tablas y el ejemplo del oráculo de los 8 campos nuevos (C1-09; corrección r2) · `getDashboardMonthVsPlan` (GET `/api/v1/dashboard/month-vs-plan`) · esquema `DashboardMonthVsPlan`: `month`, `actual_cents`, `planned_cents`, `used_pct`, `remaining_cents`, `status` (`normal`, `warning`, `exceeded`), `days_left` y, **nuevos en C1 v1.5** (§2.6): `cash_cents`, `previous_month`, `previous_actual_cents`, `previous_cash_cents`, `actual_difference_cents`, `cash_difference_cents`, `categories_warning_count`, `categories_exceeded_count`.
- **Fórmula** (motor S2): gastado (devengado) = M-01 real = Σ parte del hogar devengada en el mes; pagado en efectivo = efectivo del hogar del mes, la misma cifra que `cash_household_cents` de M-11 para ese mes (contado + partes del hogar de pagos de crédito con fecha en el mes); plan, % usado, restante y estado como M-01 (FIRME desde la puerta 2, Q2-19 B y Q2-20 A); mes anterior = las mismas dos cifras del mes calendario anterior, ausentes si el sistema no tiene datos de ese mes (sin datos no es cero); diferencias = mes − mes anterior, con signo, calculadas por el motor; conteos = categorías del mes con `status` `warning` o `exceeded` en M-09.
- **Periodo**: mes en curso (`period` = `2026-10`).
- **Contenido** (icono `wallet`):
  - Título: «Gasto del mes».
  - Valor principal: `[gastado]` con la etiqueta «gastado en [mes]» (neutro).
  - Valor secundario (360): «[n] % de [plan] · te quedan [restante]»; si `remaining_cents` < 0: «−[monto] rebasado» con `color.role.warning.*` e icono `triangle-alert` (DESIGN §3.3). Barra de avance de DESIGN §7.6 con `role="meter"` y `aria-valuetext` «[n] % del plan».
  - Desde 412 suma: «Pagado en efectivo [efectivo]» y la comparación «[diferencia] más que [mes anterior]» (icono `arrow-up-right`, `color.role.expense.*`), «[diferencia] menos que [mes anterior]» (icono `arrow-down-left`, `ui.text.secondary`) o «Igual que [mes anterior]»; sin mes anterior: «Sin datos de [mes anterior]». Las palabras «más» y «menos» van siempre en el texto (DESIGN §7.6).
  - Insignias: `warning` → «Vas en [n] %»; `exceeded` → «Rebasaste el plan» (rol advertencia, `triangle-alert`, como la tarjeta 1 de U3 1.x); si `categories_exceeded_count` > 0, «[n] categorías rebasadas» (rol advertencia), que lleva a G-3.
- **Acción primaria única**: enlace «Ver presupuesto» → P-06 Presupuesto del mes (primer paso del desglose de §2.2).
- **Vacío**: sin `planned_cents`: «—» en el % y «Aún no hay presupuesto para [mes]» con «Llevas [gastado] gastado»; la acción cambia a «Capturar presupuesto» → P-26.02.
- **Error**: «No se pudo cargar el gasto del mes. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable** (oct-2026 = P1, nov-2026 = P2; regla de ejemplos al inicio): reloj en octubre de 2026: gastado **$1,197.64** (`periodos[0].devengado.axel`, M-01) y pagado en efectivo **$598.82** (`periodos[0].efectivo.axel`: la parte del hogar del pago del 30 de octubre, R-13); «Sin datos de septiembre». Reloj en noviembre de 2026: gastado **$1,197.64** y efectivo **$1,197.64** (`periodos[1].devengado.axel` y `periodos[1].efectivo.axel`); contra octubre: «Igual que octubre» en lo gastado y «**$598.82** más que octubre» en efectivo. Plan: marcador `[plan]` (lo captura Axel).

### K-2 · Ingreso del mes
- **Métrica C1**: N-100 (P-08) · `getIncomeMonth` (GET `/api/v1/incomes/month`) · esquema `IncomeMonth`: `month`, `fortnights[]` (`pay_date`, `status`: `expected` o `captured`, `expected_cents`, `amount_cents`), `extras[]`, `captured_cents`, `expected_pending_cents`, `extras_cents`, `total_cents` y, **nuevos en C1 v1.5**: `previous_month`, `previous_fortnights_cents`, `previous_extras_cents`, `fortnights_difference_cents`. Métrica M-18 «Ingreso del mes» de C1-datos §9 (C1 v1.5, §2.6).
- **Fórmula** (motor S2): quincenas del mes = `captured_cents` + `expected_pending_cents` (capturadas, y las esperadas de P-26.03 que aún no se capturan); extra = `extras_cents`; total = `total_cents`; mes anterior y diferencia de quincenas, del motor.
- **Periodo**: mes en curso (`period` = `2026-10`).
- **Contenido** (icono `banknote`):
  - Título: «Ingreso del mes».
  - Valor principal: «+[quincenas]» con `color.role.income.*` e icono `arrow-down-left` (DESIGN §3.3), etiqueta «de quincenas».
  - Valor secundario (360): «Extra este mes +[extra]» o, con `extras_cents` = 0, «Sin extra este mes».
  - Desde 412 suma: «[n] de [m] quincenas capturadas» (conteo de `fortnights[]` a la vista) y «[diferencia] más que [mes anterior]» o «menos que», o «Sin datos de [mes anterior]».
  - Insignia si alguna quincena llega con `status` = `expected`: «Ingreso esperado» (rol información, `info`).
- **Acción primaria única**: enlace «Ver ingresos» → P-08 Ingresos.
- **Vacío** (sin quincena capturada ni esperada en P-26.03 y sin extras): «—» y «Falta tu ingreso quincenal»; la acción cambia a «Capturar ingreso» → P-26.03.
- **Error**: «No se pudo cargar el ingreso del mes. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo**: marcadores `[quincenas]` y `[extra]`; el oráculo no trae ingresos.

### K-3 · Créditos
- **Métrica C1**: M-12 a hoy (deuda viva) y M-04 (disponible) · `getCreditTotals` (GET `/api/v1/credits/totals`) · esquema `CreditTotals`: `live_debt_cents`, `month_total_cents`, `month_household_cents`, `month_external_cents` y, **nuevos en C1 v1.5**, `month_payment_count` y `month_pending_household_cents`; y `getDashboardNextPayment` (GET `/api/v1/dashboard/next-payment`) para `availability[]` (`credit_id`, `availability_reading`, `available_cents`, `availability_confirmed`). Nombres de línea: `listCredits`.
- **Fórmula** (motor S3): deuda viva total = capital + intereses por pagar de todas las líneas (M-12, a hoy); pagos del mes = Σ pagos programados con fecha en el mes en curso de todas las líneas (calendario unificado, M-13), total y partes; cuántos = número de esos pagos; falta pagar = parte del hogar aún sin pagar de esos pagos; disponible = M-04 con la **lectura activa**.
- **Qué lectura muestra** (U3-07, decisiones §C): como la tarjeta 3 de U3 1.x, solo la lectura activa del interruptor (por defecto **B**) con su etiqueta «Lectura B»; la otra se ve lado a lado en Detalle › D-3 y en P-18.
- **Contenido** (icono `credit-card`):
  - Título: «Créditos».
  - Valor principal: `[deuda viva]` con `color.role.debt.*` y la etiqueta «deuda viva total» (rol con icono, DESIGN §2.4).
  - Valor secundario (360): «Pagos de [mes]: [total del mes] · [n] pagos».
  - Desde 412 suma: «Hogar [parte del hogar] · Karen [parte de Karen]» (Karen solo si `month_external_cents` > 0) y «Falta pagar [falta pagar]».
  - Bloque «Disponible hoy», una fila por línea de `availability[]`: «Disponible hoy [monto] · Lectura [A/B]»; si `availability_confirmed` = false, insignia **«Pendiente confirmar»** (rol información, `info`) con nombre accesible «Reposición de capital pendiente de confirmar con Kueski». Con más de una línea revolvente: la primera y «+[n] líneas».
- **Acción primaria única**: enlace «Ver créditos» → P-17 Todas las líneas.
- **Vacío**: sin créditos: «Aún no hay créditos» y la acción «Agregar crédito» → P-26.05; sin pagos en el mes: «Sin pagos en [mes]» y el resto igual.
- **Error**: «No se pudieron cargar los créditos. Revisa que estés en la red de casa.» · «Reintentar». Si solo falla `getDashboardNextPayment`, el bloque «Disponible hoy» dice «—» con «Reintentar» y el resto se ve.
- **Ejemplo trazable** (U3-07): hoy 10-oct-2026: «Pagos de octubre: **$973.75** · 1 pago», «Hogar **$598.82** · Karen **$374.93**» (R-13; `periodos[0].efectivo.n_pagos` = 1); «Disponible hoy **$6,905.50** · Lectura B · Pendiente confirmar» y con el interruptor en A **$10,000.00** (R-64). Deuda viva de hoy: marcador `[deuda viva]`; la cifra de control del motor es la deuda viva al 10-abr-2027, **$6,293.44** (R-53).

### K-4 · Pagos próximos
- **Métrica C1**: M-13 (N-192) · `listScheduledPayments` (GET `/api/v1/scheduled-payments`: `id`, `due_date`, `amount_cents`, `credit_id`, `shares[]` con `person_id`, `amount_cents`, `paid_cents`, `status`; `status`) en **dos peticiones**, como D-3 (§3): (1) vencidos, solo si `getDashboardNextPayment` trae `is_overdue` = true (la respuesta de K-3 se reutiliza); (2) próximos: `status` = `pending` y `period` = `[hoy]..[hoy + 30]`. Las dos con orden `due_date` y `limit` = **N**. Días de «Por vencer»: `getAlertSettings` (§1.2).
- **N y rango** (**PROVISIONAL**, bloqueo `docs/estado/bloqueos/U3-tablero-recordatorio.md`): N = 5 filas en total (primero los vencidos), rango de hoy a hoy + 30 días («los pagos del mes» de Axel, también cerca del fin de mes). N = `parameters["dashboard.payment_reminder_count"].value` de `getDashboardSettings` (C1 v1.5; clave de `app_setting` con valor inicial 5, K44), leído una sola vez, como `getAlertSettings` en §1.2; la interfaz no fija N.
- **Contenido** (icono `calendar-clock`): título «Pagos próximos»; una fila por pago: fecha con día de la semana («Vie 30 oct 2026»), línea, monto total y estado de §1.4 («Vencido» con «Sin marcar como pagado», «Por vencer», «Parcial» con «Falta la parte de Karen», «Pendiente»; «Pagado» no llega). Desde 412, debajo: «Hogar [parte] · Karen [parte]». Cada fila es un enlace a P-18 de su línea (F-4). Si hay más pagos que N: «y [n] más en Calendario».
- **Acción primaria única**: enlace «Ver calendario» → pestaña Calendario de Inicio (§5).
- **Vacío**: «Sin pagos en los próximos 30 días»; sin créditos: «Aún no hay créditos» · «Agregar crédito» → P-26.05.
- **Error**: «No se pudieron cargar los pagos. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable**: hoy 10-oct-2026, una fila: «Vie 30 oct 2026 · Kueski · **$973.75** · Pendiente», «Hogar **$598.82** · Karen **$374.93**» (R-13); el pago del 15 de noviembre (R-14) cae fuera de los 30 días. Jue 29 oct 2026: «Vence mañana» (§1.4, A-17).

### B-1 · Barómetro de capacidad de endeudamiento
- **Métrica C1**: M-20 «Capacidad de endeudamiento» (C1 v1.5) · **NUEVA** `getDashboardDebtCapacity` (GET `/api/v1/dashboard/debt-capacity`, `period` opcional, por omisión el mes en curso) · esquema nuevo `DashboardDebtCapacity` (campos exactos en §2.6): `month`, `income_cents`, `income_is_expected`, `extras_cents`, `payments_cents`, `payments_total_cents`, `payment_count`, `ratio_pct`, `ratio_total_pct`, `warning_pct`, `limit_pct`, `scale_max_pct`, `status`, `headroom_cents`, `has_credits`.
- **Fórmula** (motores S2 y S3; K54 decisión 1): capacidad usada = pagos del mes de todos los créditos ÷ ingreso mensual neto × 100. `payments_cents` = Σ parte del hogar de los pagos programados con fecha en el mes, de todas las líneas, pagados o no (la misma cifra que `month_household_cents` de K-3). `income_cents` = quincenas del mes, capturadas más esperadas sin capturar, **sin extras** (**PROVISIONAL**). `ratio_pct` = entero con ROUND_HALF_UP (protocolo §6). Estado: `normal` (verde) si `ratio_pct` ≤ `warning_pct`; `warning` (amarillo) si `ratio_pct` ≤ `limit_pct`; `exceeded` (rojo) si `ratio_pct` > `limit_pct`. Así se lee K54 decisión 2, «verde hasta 30 %, amarillo de 30 a 40 %, rojo arriba de 40 %»: 30 exacto es verde y 40 exacto es amarillo. `headroom_cents` = `limit_pct` × `income_cents` ÷ 100 (ROUND_HALF_UP) − `payments_cents`: lo que todavía cabe en pagos al mes antes del rojo; es negativo si ya se pasó. `scale_max_pct` = el múltiplo de 10 igual o mayor que 1.5 × `limit_pct` (60 con los valores iniciales).
- **PROVISIONAL** (bloqueo `docs/estado/bloqueos/U3-tablero-barometro.md`, puerta 4): qué es «ingreso neto» (sin extras), qué pagos cuentan (parte del hogar; `ratio_total_pct` con la parte de Karen solo se informa), el redondeo a entero y la escala.
- **Umbrales** (K54 decisión 2, K44): `debt.capacity_warning_pct` = 30 y `debt.capacity_limit_pct` = 40, configurables en Ajustes; llegan en la respuesta.
- **Contenido** (icono `credit-card`; componente «Barómetro», DESIGN §8.9):
  - Título: «Capacidad de endeudamiento».
  - Arco de 0 % a `scale_max_pct` con tres tramos: verde hasta `warning_pct` (`ui.role.success`), amarillo hasta `limit_pct` (`ui.role.warning`) y rojo hasta el final (`ui.role.error`). Marcas en texto «[warning] %» y «[limit] %» en los límites de tramo. Aguja en `ratio_pct`; si `ratio_pct` > `scale_max_pct`, la aguja queda al final y el dígito dice la cifra exacta.
  - Dígito: «[n] %» (`font.size.title-1`, `font.weight.bold`, `tnum`, `ui.text.default`) y debajo «de [limit] % de tope».
  - Insignia de estado, siempre con icono y texto: `normal` → «Verde · dentro del rango» (`circle-check`, rol éxito); `warning` → «Amarillo · cerca del tope» (`triangle-alert`, rol advertencia); `exceeded` → «Rojo · arriba del tope» (`circle-alert`, rol error).
  - Valor secundario (360): «Pagos [pagos] de un ingreso de [ingreso]».
  - Desde 412 suma: «Te caben [headroom] más al mes antes del rojo» o, con `headroom_cents` < 0, «−[monto] arriba del tope» (rol advertencia, `triangle-alert`, DESIGN §3.3); si `income_is_expected`, insignia «Ingreso esperado» (`info`). Desde 768: «Con la parte de Karen: [n] %» (`ratio_total_pct`, informativo, §9).
  - Texto accesible: `role="meter"`, `aria-valuemin` = 0, `aria-valuemax` = `scale_max_pct`, `aria-valuenow` = `ratio_pct` y `aria-valuetext` «Capacidad de endeudamiento [n] % de [limit] %: [verde, amarillo o rojo]». Ejemplo **SINTÉTICO** (no sale del oráculo, que no trae ingresos): «Capacidad de endeudamiento 27 % de 40 %: verde».
- **Acción primaria única**: enlace «Ver créditos» → P-17 Todas las líneas.
- **Vacío**: sin `income_cents`: arco sin aguja, dígito «—», «Falta tu ingreso del mes para calcular la capacidad» y la acción «Capturar ingreso» → P-26.03. Sin créditos (`has_credits` = false): el dígito muestra el `ratio_pct` que llega (0 %, cifra real: no hay pagos) con la línea «Sin créditos activos».
- **Error**: «No se pudo calcular la capacidad de endeudamiento. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable** (oct-2026 a dic-2026 = P1 a P3): pagos del hogar **$598.82** en octubre (1 pago), **$1,197.64** en noviembre (2 pagos) y **$1,482.57** en diciembre (2 pagos) (`periodos[i].efectivo.axel` y `n_pagos`; en P1 a P3 el efectivo del hogar son solo pagos de crédito); con la parte de Karen, **$973.75**, **$1,947.50** y **$2,232.43** (`periodos[i].efectivo.total`). Ingreso: marcador `[ingreso]`; por eso el porcentaje del ejemplo es `[n]`.

### 2.1 Gráficas de análisis (G-1 a G-6)
Reglas comunes: las de §4 (DESIGN §7.7 y §8): orden fijo de series con forma, patrón y etiqueta; etiquetado selectivo (DESIGN §8.5); un solo eje Y; `<figure>` con `<figcaption>` y resumen de una frase sin variaciones calculadas; «Ver como tabla» con las mismas cifras; esqueleto sin ejes falsos; alto `size.chart.height-compact` en < 768 y `size.chart.height` desde 768. **0 donas** en Inicio (G-2 usa barras; §4). G-4 y G-5 son T-2 y T-3 de §4, que **se conservan** en Tendencias y se citan aquí sin cambio. **Horizonte de las gráficas** (**PROVISIONAL**, bloqueo `docs/estado/bloqueos/U3-tablero-horizonte.md`): G-1, los 11 meses anteriores y el mes en curso; G-6, los 6 anteriores, el en curso y los 5 siguientes (la proyección es lo que previene); G-4 y G-5, las ventanas de T-2 y T-3. Para G-1 y G-6 la PWA hace **una** petición a `getDashboardMonthlyOverview` con `from` = mes en curso − 11 y `to` = mes en curso + 5 (17 meses) y cada gráfica muestra su tramo.

### G-1 · Ingreso, gasto y pagos por mes
- **Métrica C1**: M-21 «Resumen mensual del tablero» (C1 v1.5) · **NUEVA** `getDashboardMonthlyOverview` (GET `/api/v1/dashboard/monthly-overview`, `from` y `to` obligatorios) · `months[]`: `month`, `is_projection`, `income_total_cents`, `income_fortnights_cents`, `income_extras_cents`, `income_is_expected`, `accrued_cents`, `cash_cents`, `planned_cents`, `credit_payments_cents`, `credit_payments_total_cents`, `debt_ratio_pct`, `debt_status` (§2.6).
- **Fórmula y periodo**: por mes: ingreso = quincenas + extras (K-2); gasto = devengado del hogar (M-01 real); pagos de crédito = parte del hogar de los pagos con fecha en el mes (B-1). Los 11 meses anteriores y el en curso.
- **Patrón**: **columnas agrupadas**, 3 series (DESIGN §8.2, máximo 3): 1 «Ingreso» = `color.chart.1`, liso; 2 «Gasto» = `color.chart.2`, rayas a 45°; 3 «Pagos de crédito» = `color.chart.3`, rayas a 135°; 2 px de separación, sin contorno. Valores en texto sobre las columnas del mes en curso; los demás, en la tabla. Un mes sin ingreso o sin gasto no lleva esa columna (sin datos no es cero) y la tabla dice «—».
- **Nota fija** bajo la leyenda: «Los pagos de crédito pagan compras que ya cuentan en el gasto de su mes; no se suman al gasto».
- **Resumen** (DESIGN §7.7): «En [mes en curso]: ingreso [monto], gasto [monto] y pagos de crédito [monto]».
- **Vacío**: «Aún no hay meses con ingreso ni gasto» · «Capturar ingreso» → P-26.03.
- **Ejemplo trazable**: gasto **$1,197.64** en octubre y en noviembre de 2026 y **$1,767.51** en diciembre (`periodos[0..2].devengado.axel`); pagos de crédito **$598.82**, **$1,197.64** y **$1,482.57** (`periodos[0..2].efectivo.axel`); ingreso: marcador `[ingreso]`.

### G-2 · Deuda viva por línea
- **Métrica C1**: M-12 (N-092) · `getLiveDebt` (GET `/api/v1/debt/live`, `date` = hoy) · `LiveDebt`: `date`, `total_cents`, `credits[]` (`credit_id`, `total_cents`); nombres con `listCredits`.
- **Patrón**: **barras horizontales** (DESIGN §8.2, «¿cuánto en cada…?»), orden descendente por monto, una serie `color.chart.1` sin patrón, valor en texto al final de cada barra. No es dona: con una sola línea (hoy solo Kueski) un anillo completo no informa y con más de 5 líneas DESIGN §8.4 pide barras; Inicio sigue con 0 donas.
- **Resumen**: «Deuda viva hoy: [total]; la mayor, [línea] con [monto]».
- **Vacío**: «No hay deudas vivas» (con `circle-check`).
- **Ejemplo**: a la fecha de control del motor, 10-abr-2027, Kueski **$6,293.44** (R-53 = **$1,258.70** R-54 + **$5,034.74** R-55, M-12); hoy, marcador `[deuda viva]`.

### G-3 · Presupuesto por categoría (con la sugerencia subir o recortar)
- **Métrica C1**: M-09 (N-081) · `getBudgetByCategory` (`period` = mes en curso; `rows[]`: `category_id`, `planned_cents`, `actual_cents`, `difference_cents`, `used_pct`, `status`, `expense_count`) · umbrales `alerts.budget_warning_pct` y `alerts.budget_limit_pct` de `getAlertSettings` (§1.2) · sugerencia: M-22 (C1 v1.5) · **NUEVA** `getBudgetSignals` (GET `/api/v1/budget/signals`, `period`) · `BudgetSignals`: `month`, `months_considered`, `margin_pct`, `rows[]` (`category_id`, `signal`: `raise`, `cut` o `none`, `months[]` con `month`, `planned_cents`, `actual_cents`, `used_pct`). Nombres con `listCategories`.
- **Regla de la sugerencia** (**PROVISIONAL**, bloqueo `docs/estado/bloqueos/U3-tablero-senal-presupuesto.md`; la calcula el motor S2, no la interfaz): con los **3** meses cerrados anteriores al mes pedido (`budget.signal_months`, valor inicial 3), una categoría lleva `raise` («subir presupuesto») si en los 3 su `used_pct` pasó de 100 + **10** (`budget.signal_margin_pct`, valor inicial 10), y `cut` («recortar presupuesto») si en los 3 quedó abajo de 100 − 10. Un mes sin plan de esa categoría corta la racha. El mes en curso no cuenta porque no ha cerrado.
- **Patrón**: **barras horizontales en % del plan** (DESIGN §8.2, fila de 1.1): una barra por categoría con plan, largo = `used_pct`, serie `color.chart.1`; con `status` `warning` o `exceeded`, la barra va en `ui.role.warning` con `triangle-alert` y el texto «Vas en [n] %» o «Rebasada». Marcas verticales en `alerts.budget_warning_pct` (1 px, trazo 4 · 4, `ui.viz.baseline`) y en `alerts.budget_limit_pct` (2 px, `ui.viz.baseline`), cada una con su «[n] %» en texto. El eje va de 0 % al mayor entre 100 % y el mayor `used_pct` a la vista. Orden descendente por `used_pct` (la interfaz ordena lo visible). Etiqueta de cada barra: «[categoría] · [real] de [plan] · [n] %». Con `signal` ≠ `none`, insignia de rol información (`info`): «Sugerencia: subir presupuesto» o «Sugerencia: recortar presupuesto», con nombre accesible «[categoría]: 3 meses seguidos arriba del plan por más de 10 %» (o «abajo»).
- **Tamaño**: a 360 y 412, las 8 categorías de mayor % y el enlace «Ver todas» → P-06; desde 768, todas. Las categorías con gasto y sin plan no llevan barra: línea «[n] categorías con gasto y sin plan» · «Capturar presupuesto» → P-26.02.
- **Desglose**: tocar una barra abre P-02 con `category_id` y `period` = mes (§2.2).
- **Resumen**: «En [mes]: [n] categorías arriba de [warning] % y [m] rebasadas» (conteo de lo que está a la vista).
- **Vacío**: «Aún no hay presupuesto por categoría» · «Capturar presupuesto» → P-26.02. Sin `getBudgetSignals` (error o C1 sin publicarla), las barras se ven sin insignias de sugerencia.
- **Ejemplo**: marcadores (`[categoría]`, `[real]`, `[plan]`, `[n]`): el plan lo captura Axel (M-09).

### G-4 · Devengado contra efectivo
Es **T-2** de §4 sin cambio (M-11, `getAccrualVsCash`, capas `cumulative_cash_cents` y `cumulative_accrued_cents`, Q2-65 A), con su resumen, su nota de proyección y sus ejemplos (R-45 a R-52). En el tablero no lleva filtros. Usa la misma petición que la pestaña Tendencias.

### G-5 · Deuda viva proyectada
Es **T-3** de §4 sin filtro de crédito (M-17, `getDebtTrend`), con su resumen y su ejemplo (R-53 a R-55). En el tablero no lleva filtros.

### G-6 · Tendencia del barómetro
- **Métrica C1**: M-21 (C1 v1.5) · `getDashboardMonthlyOverview` (la petición de G-1) · `months[]`.`debt_ratio_pct`, `debt_status`, `is_projection`; `warning_pct` y `limit_pct` de la respuesta.
- **Fórmula y periodo**: por mes, la misma capacidad de B-1; meses posteriores al en curso (`is_projection` = true) con los pagos programados (S3) y las quincenas esperadas de P-26.03. Los 6 meses anteriores, el en curso y los 5 siguientes (**PROVISIONAL**, horizonte).
- **Patrón**: **línea** (DESIGN §8.2, tiempo): una serie «Capacidad de endeudamiento», `color.chart.1`, 2 px, marcador círculo; meses de proyección con trazo 8 · 4 y la nota «Proyección con pagos programados e ingreso esperado». Líneas de referencia horizontales de 1 px en `ui.viz.baseline` en `warning_pct` y `limit_pct`, con su «[n] %» en texto al final; línea vertical «Hoy» (DESIGN §8.3). Bajo cada punto, el estado no va solo en color: lo dicen la tabla y la descripción emergente («Diciembre de 2026, [n] %: amarillo»). Eje Y de 0 % al mayor entre `scale_max_pct` y el mayor `debt_ratio_pct` a la vista. Un mes sin ingreso no lleva punto y corta la línea; la tabla dice «—» con «Falta ingreso».
- **Resumen**: «[mes en curso]: [n] % ([estado]); el mayor de la proyección es [mes] con [n] %» (el mayor de lo visible).
- **Vacío**: «Falta tu ingreso para calcular la capacidad» · «Capturar ingreso» → P-26.03.
- **Ejemplo**: pagos del hogar por mes de B-1 (**$598.82**, **$1,197.64**, **$1,482.57**); con ingreso `[ingreso]`, porcentajes `[n]`.

### 2.2 Desglose en cadena (del tablero al precio del kilo)
Responde «cuánto costó el kg de cebollas» (K54) sin pantallas nuevas: cada paso **enlaza** una pantalla que ya existe (U1 §2.1) con sus filtros.

| Paso | Dónde | Qué toca la persona | Qué ve | Necesidad U1 |
|---|---|---|---|---|
| 1 | Inicio › Resumen (K-1 o G-3) | Una barra de G-3, o «Ver presupuesto» en K-1 | Categorías del mes con su semáforo contra el plan (P-06, o P-02 con la categoría) | N-020, N-081 |
| 2 | P-02 Todos los gastos | Llega con `category_id` y `period` = mes | Gastos del mes de esa categoría y sus totales («Resumen del periodo» con el mes anterior) | N-040, N-041, N-042 |
| 3 | P-04 Gasto | Toca un gasto | Renglones del ticket: producto, presentación, cantidad, precio unitario e importe | N-061 |
| 4 | P-15 Productos y precios | Toca el producto de un renglón | Precios del producto por tienda y presentación, con precio por unidad base, fecha y fuente | N-171 |
| 5 | P-15 › Tendencias | Toca «Tendencias» | Historial del precio por unidad base por fecha y tienda (p. ej. «Cebolla · [precio] por kg» en cada tienda) | N-171 |

Toques desde Inicio hasta el historial del precio: 4 (barra, gasto, producto y pestaña). El paso 4 pide que el renglón de P-04 enlace a P-15 con el producto elegido; si U1 o la implementación no lo tienen, va como solicitud (§15). Ejemplo: marcador `[precio]`; la semilla no trae cebolla.

### 2.3 Qué miro para prevenir
La tabla que pidió Axel («a mi me sirven bastante para irme previniendo»): pregunta del mes, dónde se contesta y qué señal mirar. Las acciones son sugerencias de lectura, no reglas del sistema.

| Pregunta | Dónde | Señal que previene | Qué hacer |
|---|---|---|---|
| ¿Cuánto gastamos realmente este mes? | K-1 | % del plan, «Vas en [n] %», «Rebasaste el plan» | Bajar a G-3 y a la categoría |
| ¿Vamos peor que el mes pasado? | K-1 (desde 412) | «[diferencia] más que [mes anterior]» | Ver qué categoría creció (G-3, §2.2) |
| ¿Qué categoría se pasa? | G-3 | Barra amarilla, «Rebasada», marcas de 80 % y 100 % | Desglose §2.2 hasta el gasto y el producto |
| ¿Subimos o recortamos el presupuesto? | G-3 | «Sugerencia: subir presupuesto» o «recortar» (3 meses seguidos) | Ajustar el plan en P-26.02 |
| ¿Subió el precio de algo (el kg de cebolla)? | §2.2 → P-15 › Tendencias | Historial por unidad base y por tienda | Comparar tienda en P-14 |
| ¿Cuánto del ingreso se va en créditos? | B-1 | Color y dígito del barómetro; «Te caben [monto] más» | Con amarillo o rojo, no tomar otra compra a crédito; simularla antes en P-20 |
| ¿Va a empeorar en los próximos meses? | G-6 | Puntos de proyección arriba de la línea de 30 % o 40 % | Planear antes del mes que se pasa |
| ¿Qué pagos vienen? | K-4 | «Vence mañana», «Vence hoy», «Vencido» | Pagar y marcar en P-18 |
| ¿Cuánto debemos en total y cuándo baja? | K-3 y G-5 | Deuda viva total y su curva | Ver la línea en P-18 |
| ¿Ya entró la quincena? ¿Hubo extra? | K-2 | «Ingreso esperado», «Extra este mes» | Capturar la quincena en P-08 |
| ¿Gastamos más de lo que entra? | G-1 | Columna de gasto más alta que la de ingreso | Revisar G-3 y el plan |
| ¿Lo que se consume se está pagando? | G-4 | Capa «diferencia» que crece mes con mes | Ver pagos (K-4) y la deuda (G-5) |

### 2.4 Orden y rejilla (móvil 412 primero, escritorio 1440)
- **360 y 412** (1 columna, orden de apilado): 1 K-1 · 2 B-1 · 3 K-4 · 4 K-2 · 5 K-3 · 6 G-3 · 7 G-1 · 8 G-6 · 9 G-2 · 10 G-5 · 11 G-4 · 12 «Más del hogar» (tarjetas 2, 4, 5 y 6). Primero el propósito (gasto real), luego el «panel máximo» de Axel y lo que pide acción (pagos). A 360, cada bloque con título, valor principal y una línea secundaria; desde 412, una línea más (U1 §6.2). La cifra principal y el dígito del barómetro nunca se truncan (U1 §6.1 regla 5).
- **768** (8 columnas): K-1 y B-1 en una fila (4 + 4); K-4 a todo el ancho; K-2 y K-3 (4 + 4); gráficas, una por fila; «Más del hogar», 2 × 2.
- **1024 y 1440** (12 columnas; en 1024 con la barra lateral fija): fila 1, K-1 (6) y B-1 (6); fila 2, K-2 (4), K-3 (4) y K-4 (4); filas 3 a 5, gráficas de dos en dos (6 + 6): G-3 con G-1, G-6 con G-2, G-5 con G-4; fila 6, «Más del hogar» (1440: 4 tarjetas de 3 columnas; 1024: 2 × 2). Medianil `space.6`, relleno de tarjeta `space.6`; todas las tarjetas de una fila con el mismo alto (DESIGN §7.6).
- **Teclado y lector**: el orden del foco es el de 360 en todos los anchos (el DOM sigue ese orden; la rejilla no lo cambia).

### 2.5 Más del hogar (tarjetas que se conservan)
**PROVISIONAL** (bloqueo `docs/estado/bloqueos/U3-tablero-mas-del-hogar.md`): K54 retira la regla de 6 tarjetas pero no pide quitar de Inicio la quincena, las comidas, la canasta ni el ahorro. Por eso las tarjetas 2, 4, 5 y 6 de U3 1.x siguen **sin cambio de contenido** al final del tablero, bajo el título de sección «Más del hogar» (`font.size.title-3`). Así siguen en pie los flujos que entran por ellas (F-2 y F-3 ruta A de U1 §4) y sus pruebas. La tarjeta 1 pasa a K-1. La tarjeta 3 se reparte entre K-3 (disponible A/B) y K-4 (próximo pago con estado). Cada tarjeta conserva su texto, su operación, sus estados y sus ejemplos:

### Tarjeta 2 · Flujo de la quincena
- **Métrica C1**: M-02 (N-021) · `getDashboardFortnightFlow` (GET `/api/v1/dashboard/fortnight-flow`) · esquema `DashboardFortnightFlow`, campos `start_date`, `end_date`, `income_cents`, `income_is_expected`, `spent_cents`, `committed_cents`, `closing_balance_cents`.
- **Fórmula** (motor S2, efectivo, parte del hogar): ingresos = quincena capturada o esperada + extras de la quincena; egresos hechos = salidas de efectivo hasta hoy (contado + partes del hogar pagadas); comprometido = partes del hogar pendientes de la quincena + fijos y venideros sin pagar; saldo al cierre = ingresos − egresos hechos − comprometido.
- **Periodo**: quincena en curso con los límites de K8 (del 1 al 15 y del 16 al último día; el 31 cae en la segunda). Petición: `period` = quincena en curso (`2026-10-q1` o `2026-10-q2`).
- **Contenido** (icono `banknote`):
  - Título: «Flujo de la quincena»; subtítulo «Del [día] al [día] [mes]» con las fechas de `start_date` y `end_date` (p. ej. «Del 16 al 31 oct»).
  - Valor principal: `[saldo al cierre]` con la etiqueta «al cierre». Negativo: «−[monto] faltante» con `color.role.warning.*` e icono `triangle-alert`.
  - Valor secundario (360): «Entra +[ingresos] · Sale −[egresos + comprometido]» (ingreso con `arrow-down-left` y `color.role.income.*`; salida con `arrow-up-right` y `color.role.expense.*`). Desde 412 separa: «Ya salió −[egresos] · Por pagar −[comprometido]».
  - Insignia si `income_is_expected` = true: «Ingreso esperado» (rol información, icono `info`).
- **Acción primaria única**: enlace «Ver flujo» → P-07 Flujo y proyección.
- **Vacío** (sin quincena capturada ni esperada en P-26.03): ingresos «—» con «Falta tu ingreso quincenal»; la acción cambia a «Capturar ingreso» → P-26.03 (Ajustes › Ingresos). Egresos y comprometido se siguen mostrando.
- **Error**: «No se pudo cargar el flujo de la quincena. Revisa que estés en la red de casa.» · «Reintentar».
- **FIRME** desde la puerta 2 (Q2-14 A; bloqueo `docs/estado/bloqueos/S2-flujo-quincenal.md`, que sigue abierto solo por Q2-12): una partida planificada se da por cumplida con un gasto o ingreso con el mismo `planned_item_id` en la misma quincena (o en el mismo mes si su frecuencia es mensual o mayor), así un fijo ya pagado no cuenta dos veces en «Por pagar»; la tarjeta muestra `committed_cents` tal como llegue. La parte del hogar del pago del 30 de octubre (R-13) no depende de esto.
- **Ejemplo trazable**: hoy 16-oct-2026, quincena «Del 16 al 31 oct»; «Por pagar» incluye la parte del hogar del pago del 30 de octubre, **$598.82** (R-13, M-02); ingresos: marcador `[ingresos]`.

### Tarjeta 4 · Comidas de la semana
- **Métrica C1**: M-05 (N-023) · `getDashboardWeekMeals` (GET `/api/v1/dashboard/week-meals`) · esquema `DashboardWeekMeals`: `week_start`, `week_end`, `cost_cents`, `persons[]` (`person_id`, `avg_kcal_per_day`, `days_with_plan`, `daily_kcal_target`), `unplanned_slots`. Nombres: `display_name` de `listPersons`.
- **Fórmula** (motor S4): costo = Σ costo de las casillas con receta («fuera» no suma); kcal promedio por día de cada persona = Σ kcal asignadas ÷ días con al menos un tiempo asignado; casillas sin armar = pares día × tiempo obligatorio (desayuno, comida, cena) sin receta ni «fuera» para alguna persona, máximo 21 (K43, FIRME).
- **Periodo**: **semana en curso** de lunes a domingo (K8), como dice U1 N-023; K38 **no** cambia el contenido de la tarjeta. La PWA manda siempre `period` explícito con la semana ISO en curso (p. ej. `2026-W41`) para no depender del valor por omisión de la operación.
- **Contenido** (icono `utensils`):
  - Título: «Comidas de la semana»; subtítulo «Lun [d] al Dom [d] [mes]».
  - Valor principal: `[costo]` con la etiqueta «costo de la semana».
  - Valor secundario (360): «[n] casillas sin armar» (si es 0: «Semana completa» con `circle-check`). Desde 412, una línea por persona en orden fijo Axel, Ximena (DESIGN §2.6), con punto `color.person.axel.*` / `color.person.ximena.*` y el nombre en texto: «Axel · [kcal] kcal al día de [meta] · [n] de 7 días». Sin meta: «sin meta». Sin días con plan: «—».
- **Acción primaria única**: enlace «Ver plan» → **P-10 Plan de la semana**, que abre por defecto la semana según **K38** (U1 §0.2): los domingos, la semana que empieza el lunes siguiente; el resto de los días, la semana en curso; «Cambiar semana» siempre disponible en P-10.
  - Caso domingo: hoy **Dom 11 oct 2026** → la tarjeta muestra «Lun 5 al Dom 11 oct» (semana `2026-W41`, N-023) y «Ver plan» abre P-10 en «Lun 12 al Dom 18 oct» (`2026-W42`, K38).
  - Caso entre semana: hoy **Mié 14 oct 2026** → la tarjeta muestra «Lun 12 al Dom 18 oct» (`2026-W42`) y P-10 abre la misma semana.
- **Vacío** (`unplanned_slots` = 21 y sin casillas con receta): valor principal «—»; secundario «Aún no hay plan para esta semana»; la acción cambia a «Armar plan» → P-10 (misma regla K38).
- **Error**: «No se pudo cargar el plan de comidas. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo**: marcadores (`[costo]`, `[kcal]`, `[meta]`, `[n]`); el oráculo no trae comidas (M-05).

### Tarjeta 5 · Mejor canasta
- **Métrica C1**: M-06 (N-024) · `getDashboardBestBasket` (GET `/api/v1/dashboard/best-basket`) · esquema `DashboardBestBasket`: `basket` (`shopping_list` o `expense`), `expense_id`, `store_id`, `total_cents`, `difference_vs_second_cents`, `covered_count`, `item_count`, `stale_price_count`, `assumes_membership` y, desde C1 v1.3, `second_store_id` (Q2-86 A). Nombre de tienda: `listStores`.
- **Fórmula** (motor S5): lista de compras vigente (o la última compra con renglones si está vacía) comparada por tienda; ganadora = menor total; diferencia contra la segunda; cobertura = artículos con precio ÷ artículos; precios viejos usados; Sam's Club como si tuvieran membresía.
- **Periodo**: la **lista de compras vigente** (U1 §0.2): la de la semana por defecto según K38 (los domingos, la siguiente). El parámetro `week` es obligatorio (R3-01) y la PWA manda esa semana (Dom 11 oct 2026 → `2026-W42`; Sáb 10 oct 2026 → `2026-W41`).
- **Contenido** (icono `shopping-cart`):
  - Título: «Mejor canasta».
  - Valor principal: «[tienda] [total]» con punto del color de la tienda (`color.store.<tienda>.*`) y su nombre en texto (DESIGN §2.6 regla 1).
  - Valor secundario (360): «[diferencia] menos que [segunda tienda]», con el nombre de `second_store_id` (Q2-86 A); si C1 aún no lo entrega, «[diferencia] menos que la siguiente tienda». Desde 412 suma: «[covered] de [items] artículos con precio» y, si `stale_price_count` > 0, insignia «[n] precios viejos» (rol advertencia, `triangle-alert`).
  - Si `assumes_membership` = true: insignia «Como si tuvieran membresía» (rol información, `info`).
  - Si `basket` = `expense`: línea «Tu lista está vacía; se compara la compra del [fecha]».
- **Acción primaria única**: enlace «Comparar canasta» → P-14 (abre la semana según K38).
- **Vacío** (sin `store_id`: ningún artículo con precio, o sin lista ni compras con renglones): «—» y «Aún no hay precios para comparar»; la acción cambia a «Registrar precios» → P-15.
- **Error**: «No se pudo comparar la canasta. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable**: compra semilla del 10-oct-2026 (`basket` = `expense`): «Bodega Aurrera **$3,094.50**» (R-01), «**$138.50** menos que la siguiente tienda» (K-03 de `S5.md`: Walmart **$3,233.00**, `data.json` `totales.walmart`). Con `second_store_id` (C1 v1.3, Q2-86 A): «$138.50 menos que Walmart».

### Tarjeta 6 · Ahorro y metas
- **Métrica C1**: M-07 (N-025) · `getDashboardSavings` (GET `/api/v1/dashboard/savings`) · esquema `DashboardSavings`: `saved_total_cents`, `contributed_month_cents`, `planned_month_cents`, `next_goal` (`SavingsGoal`: `name`, `target_cents`, `target_date`, `saved_cents`, `progress_pct`, `status`: tres estados por Q2-18 A, `in_progress`, `achieved` y `overdue`, como en `docs/contratos/motores/S2-interfaz.md` 1.3 §3.13; C1 v1.3 todavía declara `on_track`, `behind` y `achieved`, solicitud a C1 en §15).
- **Fórmula** (motor S2): ahorrado = Σ aportaciones − Σ retiros (fondo y metas); aportado del mes contra planificado del mes (Σ aportación mensual planificada activa); meta más próxima por fecha con % = ahorrado ÷ monto meta × 100, entero.
- **Periodo**: mes en curso (sin parámetros).
- **Contenido** (icono `piggy-bank`):
  - Título: «Ahorro y metas».
  - Valor principal: `[ahorrado]` con `color.role.savings.*` y la etiqueta «ahorrado» (rol con icono, DESIGN §2.4).
  - Valor secundario (360): «Este mes [aportado] de [planificado]». Desde 412 suma la meta más próxima: «[nombre] · [n] % · [fecha meta]» e insignia de estado: «En curso» (`in_progress`; rol información, `info`), «Vencida» (`overdue`; rol advertencia, `triangle-alert`), «Cumplida» (`achieved`; rol éxito, `circle-check`).
- **Acción primaria única**: enlace «Ver ahorro» → P-09 Ahorro y metas.
- **Vacío** (sin fondo ni metas): «—» y «Aún no tienes fondo ni metas de ahorro»; la acción cambia a «Crear meta» → P-26.13 (Ajustes › Salud y metas de ahorro). Sin `next_goal` pero con fondo: se omite la línea de la meta.
- **Error**: «No se pudo cargar el ahorro. Revisa que estés en la red de casa.» · «Reintentar».
- **FIRME** desde la puerta 2 (bloqueo `docs/estado/bloqueos/S2-metas.md` RESUELTO): la tarjeta usa `progress_pct`, entero (Q2-17 C: entero en la tarjeta de Inicio y un decimal en Metas), y los tres estados de Q2-18 A (en curso, cumplida y vencida; no hay «atrasada»). La tarjeta no muestra la aportación necesaria, así que no depende del bloqueo abierto `S2-metas-vencida`.
- **Ejemplo**: marcadores; el oráculo no trae ahorro (M-07).

### 2.6 Operaciones de C1 del tablero (publicadas en C1 v1.5 por el agente `s1-c1-v15`)
Todo dinero en centavos enteros (`MoneyCents`), porcentajes enteros con ROUND_HALF_UP (protocolo §6), meses `Month` (`AAAA-MM`). «Ausente» = el campo no llega cuando falta el dato (sin datos no es cero). Ninguna cifra la calcula la interfaz.

**Operaciones nuevas**

| operationId | Ruta | Parámetros | Esquema de respuesta (200) y campos exactos | Motor |
|---|---|---|---|---|
| `getDashboardDebtCapacity` | GET `/api/v1/dashboard/debt-capacity` | `period` (Month, opcional; por omisión el mes en curso) | `DashboardDebtCapacity`: `month` (Month); `income_cents` (MoneyCents, ausente sin quincena capturada ni esperada); `income_is_expected` (boolean); `extras_cents` (MoneyCents, informativo, no entra al cociente); `payments_cents` (MoneyCents, parte del hogar de los pagos del mes); `payments_total_cents` (MoneyCents, con la parte de Karen); `payment_count` (integer); `ratio_pct` (integer, ausente sin `income_cents` o con 0); `ratio_total_pct` (integer, con `payments_total_cents`; ausente igual); `warning_pct` (integer, `x-app-setting` `debt.capacity_warning_pct`); `limit_pct` (integer, `x-app-setting` `debt.capacity_limit_pct`); `scale_max_pct` (integer); `status` (`normal`, `warning`, `exceeded`; ausente sin `ratio_pct`); `headroom_cents` (MoneyCents con signo; ausente sin `income_cents`); `has_credits` (boolean) | S2 (ingreso) y S3 (pagos) |
| `getDashboardMonthlyOverview` | GET `/api/v1/dashboard/monthly-overview` | `from` y `to` (Month, obligatorios; `from` ≤ `to`, máximo 24 meses; si no, 422 `validation_failed`) | `DashboardMonthlyOverview`: `warning_pct`, `limit_pct` (integer, mismas claves); `months[]`: `month` (Month); `is_projection` (boolean: mes posterior al en curso); `income_total_cents`, `income_fortnights_cents`, `income_extras_cents` (MoneyCents; ausentes sin ingreso capturado ni esperado); `income_is_expected` (boolean); `accrued_cents` (MoneyCents, devengado del hogar, M-01 real; ausente en proyección y en meses sin datos); `cash_cents` (MoneyCents, efectivo del hogar, M-11; ausente igual); `planned_cents` (MoneyCents, plan del mes; ausente sin plan); `credit_payments_cents` y `credit_payments_total_cents` (MoneyCents); `debt_ratio_pct` (integer, ausente sin ingreso); `debt_status` (`normal`, `warning`, `exceeded`; ausente sin `debt_ratio_pct`) | S2 y S3 |
| `getBudgetSignals` | GET `/api/v1/budget/signals` | `period` (Month, opcional; por omisión el mes en curso) | `BudgetSignals`: `month` (Month); `months_considered` (integer, `x-app-setting` `budget.signal_months`); `margin_pct` (integer, `x-app-setting` `budget.signal_margin_pct`); `rows[]`: `category_id` (integer); `signal` (`raise`, `cut`, `none`); `months[]`: `month`, `planned_cents` (ausente sin plan), `actual_cents`, `used_pct` (ausente sin plan) | S2 |
| `getDashboardSettings` y `updateDashboardSettings` | GET y PATCH `/api/v1/settings/dashboard` | PATCH: `parameters` y `reset_parameters`, como `updateAlertSettings` (K44, C1 §11.7) | `DashboardSettings`: `parameters` con `debt.capacity_warning_pct` (30), `debt.capacity_limit_pct` (40), `dashboard.payment_reminder_count` (5), `budget.signal_months` (3) y `budget.signal_margin_pct` (10), cada uno `SettingParameterInteger`; validación: 1 ≤ warning < limit ≤ 100, 1 ≤ N ≤ 20, 2 ≤ meses ≤ 12, 1 ≤ margen ≤ 50 (rangos **PROVISIONALES**) | — |

**Campos nuevos en operaciones existentes** (todos opcionales en el esquema, para no romper P-06, P-08 ni P-17):

| Operación · esquema | Campos nuevos |
|---|---|
| `getDashboardMonthVsPlan` · `DashboardMonthVsPlan` | `cash_cents` (MoneyCents, efectivo del hogar del mes = M-11 `cash_household_cents`); `previous_month` (Month); `previous_actual_cents`, `previous_cash_cents` (MoneyCents, ausentes sin datos de ese mes); `actual_difference_cents`, `cash_difference_cents` (MoneyCents con signo, mes − mes anterior; ausentes si falta el mes anterior); `categories_warning_count`, `categories_exceeded_count` (integer) |
| `getIncomeMonth` · `IncomeMonth` | `previous_month` (Month); `previous_fortnights_cents`, `previous_extras_cents` (MoneyCents, ausentes sin datos); `fortnights_difference_cents` (MoneyCents con signo) |
| `getCreditTotals` · `CreditTotals` | `month_payment_count` (integer); `month_pending_household_cents` (MoneyCents, parte del hogar sin pagar de los pagos del mes) |

**Claves nuevas de `app_setting`** (esquema 1.5, migración 00005 con `@param`; K54 y K44): `debt.capacity_warning_pct` = 30 y `debt.capacity_limit_pct` = 40 (FIRMES, K54 decisión 2); `dashboard.payment_reminder_count` = 5, `budget.signal_months` = 3 y `budget.signal_margin_pct` = 10 (**PROVISIONALES**, bloqueos de §2). Dónde se editan en Ajustes lo decide U4 (§15).

**Métricas de C1-datos §9 (publicadas en C1 v1.5)**: M-18 «Ingreso del mes» (K-2), M-19 «Créditos del mes» (K-3: pagos del mes, cuántos y falta pagar), M-20 «Capacidad de endeudamiento» (B-1), M-21 «Resumen mensual del tablero» (G-1 y G-6) y M-22 «Sugerencia de presupuesto» (G-3), con las fórmulas de sus bloques; las reglas del barómetro y de la sugerencia son RN-30 y RN-31 de `reglas-negocio.md`.

---

## 3. Detalle

**Propósito**: el segundo nivel de los bloques del tablero (D-1 a D-6 siguen el orden de las 6 tarjetas de U3 1.x), sin salir de Inicio (U1 N-027), más el bloque «Participación de Karen» (U1 N-026). Responde «¿de dónde sale la cifra de la tarjeta?» con lo mínimo; lo completo vive en la pantalla destino.
**Contenido**: 7 paneles en el orden de las tarjetas (D-1 a D-6) y al final D-7. Cada panel es un bloque con título (`font.size.title-3`), contenido y un enlace «Ver todo» a su pantalla destino (enlace, no botón con relleno). Componentes: «Tabla compacta» y «Lista de propiedades» (DESIGN §7).
**Filtros**: barra de vista de Detalle (en < 768 px, botón «Filtrar» con icono `list-filter` que abre una hoja inferior y chips con los filtros activos bajo el título, U1 §1.3). Cada panel aplica solo los filtros que su operación admite (§6); un filtro que un panel no admite **no** lo cambia y el panel lo dice con la línea «Este bloque no usa el filtro [nombre]» (`font.size.caption`).
**Periodo por defecto**: mes en curso; la quincena es subperiodo (selector «Mes · 1.ª quincena · 2.ª quincena», valores `period` = `2026-10`, `2026-10-q1`, `2026-10-q2`; K8).
**Estados**: los de §1.3 por panel (cada panel hace su petición), con los textos de §7.

| Panel | Contenido | Métrica · operación C1 (campos) | Vacío |
|---|---|---|---|
| D-1 Presupuesto por categoría | Las 5 categorías con mayor % usado: nombre, barra horizontal real contra plan (patrón «barras para comparar categorías»), «[real] de [plan] · [n] %» e insignia de estado como en K-1 | M-09 (N-081) · `getBudgetByCategory` (`month`, `rows[]`: `category_id`, `planned_cents`, `actual_cents`, `difference_cents`, `used_pct`, `status`, `expense_count`); nombres con `listCategories` | «Aún no hay presupuesto por categoría» · «Capturar presupuesto» → P-26.02 |
| D-2 Quincena en curso | Fila de la quincena en curso como lista de propiedades: ingresos (marca «esperado» si `income_is_expected`), pagos de crédito (parte del hogar), fijos y venideros, gasto variable estimado, ahorro planificado, saldo de la quincena y saldo acumulado | M-10 (N-090) · `getCashFlowFortnights` (`rows[]`: `start_date`, `end_date`, `pay_date`, `income_cents`, `income_is_expected`, `credit_payments_cents`, `fixed_and_upcoming_cents`, `variable_estimate_cents`, `planned_savings_cents`, `balance_cents`, `cumulative_balance_cents`, `is_actual`). Gasto variable estimado: **PROVISIONAL** (`docs/estado/bloqueos/S2-flujo-quincenal.md`, Q2-12, abierta); ahorro planificado: el reparto mitad y mitad es FIRME (Q2-13 B: HALF_UP en la quincena 1 y el resto en la 2); el valor es **PROVISIONAL** (bloqueo abierto `docs/estado/bloqueos/S2-ahorro-quincenal.md`: falta la base, meta por meta o la suma del mes, qué metas cuentan y si resta del saldo). El motor S2 lo entrega en `null` y no lo resta de `balance_cents` ni de `cumulative_balance_cents` (`docs/contratos/motores/S2-interfaz.md` 1.3 §3.11); mientras llegue `null`, la fila muestra «—», nunca un monto en cero (regla «Sin datos no es cero»). C1 v1.3 declara `planned_savings_cents` entero obligatorio; solicitud a C1 en §15 | «Falta tu ingreso quincenal» · «Capturar ingreso» → P-26.03 |
| D-3 Pagos y disponible | (a) Pagos de todas las líneas en dos grupos, cada uno en orden de fecha: primero **los vencidos** (del más antiguo al más reciente) y luego **los próximos 5 pendientes** de hoy a hoy + 60. Cada fila: fecha con día de la semana («Vie 30 oct 2026»), línea, monto total, «Hogar [parte] · Karen [parte]» y estado de §1.4 («Vencido» con «Sin marcar como pagado», «Por vencer», «Parcial» con «Falta la parte de Karen», o «Pendiente»; «Pagado» no llega a este panel). Si hay más de 5 vencidos, la línea «Hay más pagos vencidos» y el enlace «Ver todo». (b) Disponible de cada línea revolvente con **las dos lecturas lado a lado** («Lectura A [monto]» · «Lectura B [monto]», la activa con la insignia «Activa») y el mínimo del horizonte de cada lectura; aviso «Reposición de capital pendiente de confirmar con Kueski» mientras no se confirme | M-13 (N-192) · `listScheduledPayments` (`id`, `due_date`, `amount_cents`, `credit_id`, `shares[]`: `person_id`, `amount_cents`, `paid_cents`, `status`; `status`) en **dos peticiones**, las dos con orden `due_date`, `limit` = 5 y los filtros de persona, familia y crédito de la barra (§6): (1) **vencidos**: `status` = `overdue` y `period` = `[fecha del vencido más antiguo]..[ayer]` (rango inclusivo de C1), donde esa fecha es `next_payment.due_date` de `getDashboardNextPayment` cuando `is_overdue` = true (M-03 pone primero el vencido más antiguo; se reutiliza la respuesta de K-3 o se pide); con `is_overdue` = false no hay vencidos y esta petición no se hace; (2) **próximos**: `status` = `pending` y `period` = `[hoy]..[hoy + 60]` (M-13 «hoy..+60»). Días de «Por vencer»: `getAlertSettings` (§1.2) · M-04 (N-285, N-201) · `previewAvailability` por línea revolvente | «Sin pagos pendientes» (b se sigue mostrando); sin créditos: «Aún no hay créditos» · «Agregar crédito» → P-26.05 |
| D-4 Comidas de hoy | Casillas de hoy por tiempo (desayuno, comida, cena, colación) y persona (Axel, Ximena): receta o «Fuera», kcal y costo | N-121 · `getMealPlanWeek` (`cells[]`: `plan_date`, `meal_slot_id`, `person_id`, `recipe_id`, `is_eating_out`, `cost_cents`, `kcal`), con `period` = semana en curso; la PWA muestra las casillas cuyo `plan_date` es hoy. **Métrica pendiente en C1-datos §9** (solicitud a C1, §15); mientras tanto rige esta definición: «Comidas de hoy» (N-121 vía N-027) = casillas de `cells[]` con `plan_date` = hoy, por tiempo y persona, con `recipe_id` o `is_eating_out`, `kcal` y `cost_cents` tal como llegan (sin sumas); periodo: la semana ISO de lunes a domingo que contiene hoy (K8; también en domingo, porque el panel muestra hoy y K38 solo cambia la semana que abre P-10); filtro: `person_id` (S) | «Hoy no hay comidas en el plan» · «Armar plan» → P-10 (regla K38) |
| D-5 Diferencias de precio | Los 5 artículos con mayor diferencia de precio entre tiendas de la canasta de la tarjeta 5: artículo, tienda ganadora (punto + nombre), diferencia, marca «viejo» si el precio es viejo; y el total por tienda (orden fijo de tiendas, DESIGN §2.6) | N-162 · `getBasketComparison` (`stores[]`: `store_id`, `total_cents`, `difference_vs_cheapest_cents`, `stale_count`, `assumes_membership`; renglones con `product_id`, `winner_store_id`, `difference_cents`, `prices[]`.`is_stale`), con el mismo `week` (y `basket`, `expense_id`) de la tarjeta 5. **Métrica pendiente en C1-datos §9** (solicitud a C1, §15); mientras tanto rige esta definición: «Diferencias de precio» (N-162 vía N-027) = los 5 renglones con mayor `difference_cents`, en orden descendente (la interfaz ordena lo que está a la vista, U1 §0.2), y `stores[]` con `total_cents` tal como llega; periodo: la lista vigente (`week` según K38, como la tarjeta 5); filtro: tienda en la vista (V) | «Aún no hay precios para comparar» · «Registrar precios» → P-15 |
| D-6 Metas | Cada meta activa: nombre, barra de avance (`color.role.savings.*` con el % en texto, con un decimal), «[ahorrado] de [meta]», aportación mensual necesaria, fecha meta y estado (insignias de la tarjeta 6) | M-16 (N-111) · `listSavingsGoals` (`name`, `target_cents`, `target_date`, `saved_cents`, `progress_tenths_pct`, `required_monthly_cents`, `status`), con `status` activo. Avance con un decimal: `progress_tenths_pct` (décimas de punto porcentual con ROUND_HALF_UP, C1 v1.3), por Q2-17 C («un decimal en Metas y entero en la tarjeta de Inicio»): D-6 es la métrica M-16 «Metas» y el entero (`progress_pct`, M-07) queda solo en la tarjeta 6; sin `target_cents` el campo llega ausente y se muestra «—». C1-datos M-16 todavía dice «% de avance entero»; solicitud a C1 en §15. Aportación mensual necesaria con los meses del mes de hoy al de la fecha meta, los dos incluidos (Q2-16 A), y estado con los tres valores de Q2-18 A (`in_progress`, `achieved` y `overdue`; insignias de la tarjeta 6); bloqueo `docs/estado/bloqueos/S2-metas.md` RESUELTO. La aportación de una meta vencida sin cumplir es **PROVISIONAL** (`docs/estado/bloqueos/S2-metas-vencida.md`, abierto) | «Aún no tienes metas de ahorro» · «Crear meta» → P-26.13 |
| D-7 Participación de Karen | Ver §3.1 | M-08 (N-026) · `getPersonParticipation` | Ver §3.1 |

Ejemplos trazables de Detalle: D-3 (a) con hoy 10-oct-2026, rango de hoy a hoy + 60 (hasta el 9-dic-2026, inclusivo): «Vie 30 oct 2026», «Dom 15 nov 2026» y «Lun 30 nov 2026», cada uno por **$973.75** con «Hogar $598.82 · Karen $374.93» y estado «Pendiente» (R-13 a R-15, M-13); no hay vencidos (`getDashboardNextPayment` trae `is_overdue` = false) y el pago del 15-dic-2026 (R-16) cae fuera del rango. D-3 (b) el 10-oct-2026: «Lectura A **$10,000.00**» · «Lectura B **$6,905.50**» (R-64); mínimo del horizonte **$8,094.50** y **$5,000.00** (R-61, R-62). D-5 con la compra semilla: Bodega Aurrera **$3,094.50** (R-01) y Walmart **$3,233.00** (`totales.walmart`), diferencia **$138.50** (K-03).

Al pie de Detalle, enlace «Ver todos los gastos con estos filtros» → P-02, que recibe persona, familia, categoría, tienda, periodo y crédito (`getExpenseSummary` y `listExpenses` los admiten todos).

### 3.1 Participación de Karen (D-7)
- **Métrica C1**: M-08 (N-026) · `getPersonParticipation` (GET `/api/v1/persons/{person_id}/participation`) · esquema `PersonParticipation`: `month`, `accrued_groceries_cents`, `purchases[]` (`description`, `purchase_date`, `total_cents`, `payments_made`, `payments_total`, `balance_cents`), `payment_shares[]` (`due_date`, `amount_cents`, `status`), `live_debt_cents`. El `person_id` es el de la persona con `participation` externa en `listPersons`; con otra persona C1 responde 404.
- **Periodo**: mes en curso.
- **Contenido**, título «Participación de Karen» con punto `color.person.karen.*` y el nombre en texto:
  1. «Su parte de la despensa este mes: [monto]».
  2. «Sus compras a crédito»: tabla compacta con comercio, fecha, total, «Pagos [hechos] de [total]» y saldo.
  3. «Sus pagos de este mes»: fecha con día de la semana, monto de su parte y estado («Pagado» con `circle-check`; «Pendiente»).
  4. «Lo que debe en la línea: [monto]».
  5. Nota fija (texto exacto de U1 N-026): «Karen paga directo a Kueski su parte; su parte es informativa y repone la línea cuando paga».
- **Acción del panel**: enlace «Estado de cuenta de Karen» → P-25, donde se genera el PDF (`generateExternalStatement`) que **Axel reenvía a mano** (decisiones §A, F). Karen no tiene cuenta ni ve el dashboard.
- **Vacío**: «Karen no tiene compras ni pagos este mes»; la nota fija se mantiene.
- **Error**: «No se pudo cargar la participación de Karen. Revisa que estés en la red de casa.» · «Reintentar».
- **Ejemplo trazable**: octubre de 2026, «Su parte de la despensa este mes: **$749.86**» (`periodos[0].devengado.karen`, M-08; P1 = oct-2026); su parte del pago del 30 de octubre, **$374.93** (R-11, R-13).

---

## 4. Tendencias

**Propósito**: ver cómo van las cifras de las tarjetas en el tiempo, con las cuatro series de U1 N-028.
**Contenido**: 4 gráficas en orden fijo T-1 a T-4, cada una con título, subtítulo de periodo, resumen de una frase debajo del título (DESIGN §7.7; su texto está en cada bloque), leyenda, etiquetas de dinero en texto y la alternativa «Ver como tabla» (enlace que cambia la gráfica por una tabla con las mismas cifras; la tabla es también el contenido para lectores de pantalla). Componentes: «Gráfica» y «Tabla compacta» (DESIGN §7) con los patrones de DESIGN §8. Biblioteca: ECharts 6 (decisiones §I).
**Filtros y periodos**: cada gráfica tiene su ventana fija de C1 (12 meses, 12 semanas); los filtros que admite están en su bloque y en §6.
**Estados**: §7.
**Reglas de todas las gráficas** (DESIGN §2.7, §2.10; K27, K41):
- Series en **orden fijo**; el color de una serie nunca va solo: lleva **forma de marcador**, **patrón de trazo o relleno** y **etiqueta directa** en texto.
- Paleta de series `color.chart.1` a `color.chart.8` (`.light` / `.dark`) en ese orden; las gráficas por persona usan la categórica `color.person.*` (DESIGN §2.6 regla 2); personas y tiendas nunca comparten gráfica ni leyenda.
- Etiquetas de dinero siempre en texto, nunca solo al tocar, con el **etiquetado selectivo de DESIGN §8.5** (nunca un número en cada punto): columnas con el valor sobre la columna si cabe con aire; líneas con el último valor y los extremos; área con los valores finales. Formato completo en etiquetas y leyendas; el compacto solo en ejes. Una etiqueta que no cabe no se recorta: se mueve fuera de la marca o pasa a «Ver como tabla». El valor al tocar (< 768) o al pasar el puntero (≥ 768) es un extra (U1 §6.2).
- Marcas de DESIGN §8.1: barras y capas separadas por 2 px del color de la superficie, **nunca un borde dibujado alrededor de una marca**; un solo eje Y; textos en `ui.text.default` o `ui.text.secondary`, nunca en el color de la serie; sin retícula vertical salvo «Hoy».
- Ejes y rejilla `color.neutral.700` / `color.neutral.300` para texto (`font.size.caption`) y divisores decorativos `color.neutral.200` / `color.neutral.800`; cifras con `tnum`.
- Sin donas en Dashboard. Si una implementación futura la necesitara, solo para composición de ≤ 5 partes (U2-13); hoy: 0 donas.
- Los meses sin cifra del oráculo no llevan marca de canon; los meses 7 a 12 se rotulan como proyección (N-093).
- **Resumen y `<figcaption>`** (DESIGN §7.7): cada gráfica va en `<figure>` con `<figcaption>`, que trae el título y el resumen de una frase de su bloque. El resumen va debajo del título, en `ui.text.secondary`. Solo enuncia cifras que llegan de C1 (o sumas de lo que está a la vista), con su mes o su semana. **No calcula variaciones ni porcentajes** («subió 12 %»), porque la interfaz no recalcula (U1 §0.2). En eso U3 se aparta del ejemplo de DESIGN §7.7 (solicitud 5 a DESIGN, §15). Si falta una de sus cifras, esa parte dice «—» con el motivo. En carga, vacío o error no hay resumen: manda el estado de §7. Con «Ver como tabla», el resumen se queda.

### T-1 · Plan contra real por mes
- **Métrica C1**: M-14 (N-085) · `getBudgetMonthlyTrend` (GET `/api/v1/budget/monthly-trend`) · `months[]`: `month`, `planned_cents`, `actual_cents`; filtro `category_id`.
- **Fórmula y periodo**: por mes (máx. 12): plan total y real total (parte del hogar devengada), o de una categoría.
- **Patrón**: **barras agrupadas** por mes (comparar dos cantidades en cada mes). Orden fijo de series: 1 «Real» = `color.chart.1`, relleno sólido; 2 «Plan» = `color.chart.2` con el patrón de la serie 2 (rayas a 45° cada 6 px, DESIGN §8.6) y **sin contorno** (DESIGN §8.1); las dos barras de cada mes se separan 2 px con el color de la superficie. Leyenda: «Real», «Plan».
- **Subtítulo**: «Últimos 12 meses · [categoría o «Todas las categorías»]».
- **Resumen** (DESIGN §7.7): «En [último mes]: real [monto] de un plan de [monto]», con `actual_cents` y `planned_cents` del último mes de `months[]`. Sin plan en ese mes: «En [último mes]: real [monto]; aún sin plan».
- **Vacío**: «Aún no hay presupuesto capturado» · «Capturar presupuesto» → P-26.02. **Ejemplo**: marcador (el plan lo captura Axel; M-14).

### T-2 · Devengado contra efectivo
- **Métrica C1**: M-11 (N-091) · `getAccrualVsCash` (GET `/api/v1/cash-flow/accrual-vs-cash`) · `projection_mode`, `months[]`: `month`, `accrued_cents`, `cash_cents`, `accrued_household_cents`, `accrued_external_cents`, `cash_household_cents`, `cash_external_cents`, `cumulative_difference_cents`, `is_canonical` y, desde C1 v1.3, `cumulative_accrued_cents` y `cumulative_cash_cents` (Q2-65 A).
- **Fórmula y periodo**: mes en curso a 11 meses después; diferencia acumulada = devengado acumulado − efectivo acumulado (motor S2). **No se rotula «deuda viva»** (K19, N-091).
- **Patrón**: **área apilada devengado contra efectivo**, acumulada. Orden fijo de capas, de abajo arriba: 1 «Efectivo acumulado» = `color.chart.1`, relleno sólido, borde superior continuo con marcador círculo; 2 «Diferencia acumulada» (la capa «devengado sin pagar» de DESIGN §8.3; el rótulo sigue a C1 M-11) = `color.chart.2`, relleno con trama diagonal, borde superior discontinuo con marcador cuadrado. La orilla superior de la pila es el devengado acumulado y lleva la etiqueta «Devengado acumulado». **Las dos capas tienen el mismo origen** que `cumulative_difference_cents`, que C1 define como `periodos[i].diferencia_acumulada` de `model.json`: se acumula desde P1 aunque `months[]` empiece en el mes en curso. Capa 1 = efectivo acumulado desde ese origen; cima = devengado acumulado desde ese origen; capa 2 = cima − capa 1 = `cumulative_difference_cents` tal como llega. **Prohibido** sumar `cash_cents` desde el primer mes visible y apilar encima `cumulative_difference_cents`: mezcla dos orígenes y, desde el segundo mes, la cima deja de ser el devengado acumulado.
- **Fuente de las capas** (Q2-65 A; bloqueo `docs/estado/bloqueos/U3-t2-acumulados.md` RESUELTO): C1 v1.3 agrega en `months[]` `cumulative_accrued_cents` y `cumulative_cash_cents`, con el mismo origen que `cumulative_difference_cents` (así `cumulative_difference_cents` = `cumulative_accrued_cents` − `cumulative_cash_cents`). Capa 1 = `cumulative_cash_cents`; cima = `cumulative_accrued_cents`; capa 2 = `cumulative_difference_cents`. Etiquetas «Devengado acumulado» y «Efectivo acumulado». Se retira la variante (b), relativa a la ventana.
- **Pago adelantado** (DESIGN §8.3): si un mes trae `cumulative_difference_cents` negativo (se pagó antes de devengar), la capa 2 **vale cero en ese tramo** y el texto «Pagado por adelantado» lo dice junto al tramo; nunca se dibuja una capa negativa y la gráfica sigue siendo un área apilada (en ese tramo la cima de la pila es el efectivo acumulado). La tabla muestra la diferencia tal como llega, con su signo menos.
- Línea vertical «Hoy» de 1 px en `ui.viz.baseline` con su etiqueta (DESIGN §8.3; única retícula vertical, §8.1).
- **Resumen** (DESIGN §7.7), con las cifras de las etiquetas finales: «Hasta [último mes], con proyección desde [séptimo mes]: devengado acumulado [monto] y efectivo acumulado [monto]». No nombra la diferencia como deuda ni como «debes» (K19, N-091, M-11).
- Etiquetas en texto al final de la pila (DESIGN §8.3 y §8.5): «Devengado acumulado» y «Efectivo acumulado», con su valor final (`cumulative_accrued_cents` y `cumulative_cash_cents` del último mes). El subtítulo da solo el periodo y **no** rotula la deuda viva (K19, N-091, M-11): U3 se aparta aquí del ejemplo de subtítulo de DESIGN §8.3; solicitud 3 a DESIGN (§15).
- Meses 7 a 12: trazo discontinuo y la nota de lo que incluye la proyección según `projection_mode` (K21, N-093). Con el interruptor encendido, esos meses repiten las despensas con el patrón del oráculo como valor inicial sembrado, que cambia con partidas planificadas (Q2-21 a Q2-24 C; bloqueo `docs/estado/bloqueos/S2-proyeccion-interruptor.md` RESUELTO), y la nota lo dice: «Proyección con despensas repetidas»; con el interruptor apagado (por omisión) no se repiten. La tabla trae también las partes del hogar y de Karen de cada mes (`accrued_household_cents`, `accrued_external_cents`, `cash_household_cents`, `cash_external_cents`); la de Karen es informativa, de la despensa en común (§9).
- **Vacío**: «Aún no hay gastos ni pagos para comparar» · «Registrar gasto» (acción de página). **Ejemplo trazable** (oct-2026 = P1 y mar-2027 = P6, regla de ejemplos al inicio; C1-datos M-01, `is_canonical`): fila oct-2026, devengado **$1,947.50**, efectivo **$973.75** y diferencia acumulada **$973.75** (R-45, M-11); fila nov-2026, **$1,947.50**, **$1,947.50** y **$973.75** (R-46). Capas, que no cambian con el mes en curso porque el origen no se mueve: en oct-2026, efectivo acumulado **$973.75** y devengado acumulado **$1,947.50** (R-45); en nov-2026, efectivo acumulado **$2,921.25** (= 973.75 + 1,947.50, efectivo de R-45 y R-46) y devengado acumulado **$3,895.00** (= 1,947.50 + 1,947.50, devengado de R-45 y R-46); en mar-2027, **$13,103.78** (R-52) y **$14,362.48** (R-51), con diferencia acumulada **$1,258.70** (R-50). **Fila «Total» de la tabla**: suma de las 12 filas visibles (la interfaz suma lo que está a la vista, U1 §0.2), sin cifra de canon: R-51 y R-52 son los totales de P1 a P6 (seis filas), no de doce.

### T-3 · Deuda viva proyectada
- **Métrica C1**: M-17 (N-194) · `getDebtTrend` (GET `/api/v1/credits/debt-trend`) · `months[]`: `month`, `live_debt_cents`, `credits[]` (`credit_id`, `live_debt_cents`), `interest_paid_cents`, `interest_pending_cents`; filtro `credit_id`.
- **Fórmula y periodo**: por mes (12), deuda viva al cierre por línea; intereses pagados y por pagar (motor S3).
- **Patrón**: **líneas** (tendencia). Sin filtro: una línea «Deuda viva total» = `color.chart.1`, continua, marcador círculo. Con filtro de crédito: una línea por línea elegida en el orden de `listCredits`, series `color.chart.1`, `color.chart.2`, … con marcadores círculo, cuadrado, triángulo, rombo y trazos continuo, discontinuo, punteado, en ese orden. Intereses pagados y por pagar solo en la tabla.
- **Resumen** (DESIGN §7.7): sin filtro, «Deuda viva al cierre de [primer mes]: [monto]; al cierre de [último mes]: [monto]», con `live_debt_cents` del primer y del último mes de `months[]`. Con filtro de crédito: «Al cierre de [último mes]: [línea] [monto] · [línea] [monto]», con `credits[]` del último mes en el orden de la leyenda. Aquí sí se rotula la deuda viva, porque es la métrica de la gráfica (M-17).
- **Vacío**: «No hay deudas vivas» (con `circle-check`). **Ejemplo trazable**: la cifra de control del motor es la deuda viva al 10-abr-2027, **$6,293.44** (R-53 = **$1,258.70** R-54 + **$5,034.74** R-55, M-17); los puntos mensuales son «al cierre» de cada mes y los entrega C1 tal cual.

### T-4 · Costo semanal y kcal de comidas
- **Métrica C1**: M-15 (N-129) · `getMealPlanTrends` (GET `/api/v1/meal-plan/trends`) · `weeks[]`: `week_start`, `cost_cents`, `persons[]` (`person_id`, `avg_kcal_per_day`, `daily_kcal_target`); filtro `person_id`.
- **Fórmula y periodo**: por semana (12, lunes a domingo, K8): costo de casillas con receta; kcal promedio por día por persona contra su meta (motor S4).
- **Patrón**: dos gráficas de **líneas** apiladas en vertical, mismo eje de semanas:
  - T-4a «Costo por semana»: una línea `color.chart.1`, continua, marcador círculo.
  - T-4b «kcal al día por persona»: orden fijo Axel, Ximena; «Axel» = `color.person.axel.*`, continua, círculo; «Ximena» = `color.person.ximena.*`, trazo 8 · 4 y marcador cuadrado (la segunda persona usa el marcador y el trazo de la serie 2, DESIGN §8.7); la meta de cada persona, línea del mismo color con trazo de puntos 2 · 4, sin marcadores y con la etiqueta «Meta de [nombre]» (DESIGN §8 no fija la línea de meta: supuesto 8 y solicitud 4 a DESIGN, §15). Etiqueta directa con el nombre al final de cada línea (≤ 4 elementos, DESIGN §2.6 regla 2).
- **Resumen** (DESIGN §7.7), uno por gráfica, con la última semana de `weeks[]`: T-4a «Semana Lun [d] al Dom [d mmm]: costo [monto]» (`cost_cents`); T-4b «Semana Lun [d] al Dom [d mmm]: Axel [n] kcal al día con meta de [n]; Ximena [n] kcal al día con meta de [n]» (`avg_kcal_per_day`, `daily_kcal_target`).
- **Vacío**: «Aún no hay semanas con plan de comidas» · «Armar plan» → P-10 (regla K38). **Ejemplo**: marcadores (M-15).

---

## 5. Calendario

**Propósito**: ver en un solo lugar qué entra y qué sale por día: pagos del calendario unificado de créditos y eventos de despensa, con las quincenas marcadas.
- **Métrica C1**: M-13 (N-029) · `listCalendarEvents` (GET `/api/v1/calendar/events`) · `CalendarEventPage` de `CalendarEvent`: `event_type`, `date`, `end_date` (solo promociones), `title`, `amount_cents`, `household_share_cents`, `external_share_cents`, `status` (`pending`, `paid`, `overdue`, `expected`, `captured`, `complete`, `incomplete`, `empty`), `credit_id`, `ref_table`, `ref_id`. Filtros: `period`, `credit_id`, `family`, `event_type`.
- **Tipos que muestra Inicio** (U1 N-029, orden dentro del día): ingresos (`income`), pagos de crédito (`credit_payment`), gastos fijos (`fixed_expense`) y venideros (`upcoming_expense`), fechas de corte (`card_statement`), despensa programada (`planned_grocery`), días de visita a mercados (`market_visit`), promociones (`promotion`, franja de inicio a fin), fechas meta de ahorro (`savings_goal`) y envío del reporte mensual (`monthly_report`). `meal_plan` y `savings_movement` son de otros calendarios y no se piden aquí.
- **Iconos por tipo** (DESIGN §5.4): ingreso `arrow-down-left` con «+» y `color.role.income.*`; pago de crédito `credit-card`; fijo y venidero `arrow-up-right` con «−» y `color.role.expense.*`; corte `calendar-clock`; despensa programada `shopping-cart`; mercado `store` con punto `color.store.market.*`; promoción `calendar` con franja; meta `piggy-bank`; reporte `file-chart-column`. Estados de los pagos de crédito: los cinco de §1.4 (DESIGN §7.8): «Pendiente», «Por vencer», «Vencido» con «Sin marcar como pagado», «Parcial» con «Falta la parte de Karen» y «Pagado». Ingresos esperados: «Esperado» (`info`, rol información). «Parcial» no viene en `CalendarEvent.status`: el Calendario pide además `listScheduledPayments` con el mismo `period`, `credit_id` y `family` y `status` = `pending`, y une cada evento `credit_payment` cuyo `ref_table` es `scheduled_payment` con el pago cuyo `id` es `ref_id` (M-13 cita las dos operaciones). Si esa petición falla, cada pago conserva el estado de `CalendarEvent` y nada más cambia.
- **Pago de crédito**: «[línea] · [monto total]» y, debajo, «Hogar [parte] · Karen [parte]» (Karen solo en pagos de la línea donde ella compra).
- **Quincenas (K8)**: la vista marca cada mes con dos franjas: «1.ª quincena · del 1 al 15» y «2.ª quincena · del 16 al [último día]»; el 31 cae en la segunda. El **día de pago** (15 y 30, o el último día del mes) es un dato aparte: lo marcan los eventos de ingreso y de pago con su fecha, no la franja. Febrero: segunda quincena del 16 al 28 (o 29) y pago el último día, p. ej. «Dom 28 feb 2027». La semana va de lunes a domingo (K8): la rejilla del mes empieza en lunes.
- **Vistas**:

| Ancho | Vista | Rango por defecto (`period`) |
|---|---|---|
| 360 y 412 | Agenda por día: encabezado de día «Vie 30 oct 2026», eventos en el orden de tipos; separador al cambiar de quincena; en 412, mes compacto opcional arriba (puntos por día) | de hoy a hoy + 30 días (rango inclusivo); «Ver 30 días más» amplía a hoy + 60 |
| 768 y 1024 | Mes: rejilla de lunes a domingo con las franjas de quincena; hasta 3 eventos por día y «+[n] más» | mes en curso (`2026-10`) |
| 1440 | Mes con el detalle del día seleccionado en panel derecho (lista completa del día con montos y partes) | mes en curso |

- **Filtros**: tipo de evento (selector múltiple «Tipos»), crédito, familia (K5: «Hogar» incluye lo conjunto; «Karen» incluye lo conjunto; «Conjunto» solo lo conjunto; sin valor: todos), periodo (navegación «‹ mes anterior · mes siguiente ›» y «Hoy»).
- **Acción del panel**: tocar un evento abre su registro (vista previa G-06: página completa en < 768 px, panel en ≥ 1024 px); un pago abre P-19 de su compra o P-18 de su línea.
- **Vacío**: «No hay eventos en este periodo»; si no hay ningún crédito, ingreso ni despensa programada: «Aún no hay pagos ni ingresos programados» · «Capturar ingreso» → P-26.03.
- **Ejemplo trazable**: octubre a diciembre de 2026, pagos de la despensa 1 el «Vie 30 oct 2026», «Dom 15 nov 2026», «Lun 30 nov 2026» y «Mar 15 dic 2026», cada uno por **$973.75** (R-13 a R-16), con «Hogar $598.82 · Karen $374.93».

---

## 6. Filtros y periodos por tarjeta y panel

Filtros de C1 (`componentes-comunes.yaml`): persona `person_id`, familia `family` (K5), categoría `category_id`, tienda `store_id`, periodo `period`, crédito `credit_id`. Claves: **S** = filtro en el servidor (la operación lo admite); **V** = filtra lo que está a la vista (U1 §0.2: la interfaz filtra y ordena lo que muestra, sin recalcular); **—** = no aplica (la operación no lo admite y la cifra no se puede filtrar sin recalcular; el bloque lo dice, §3). Periodo por defecto: **mes calendario**, con la quincena como subperiodo (decisiones §B).

| Tarjeta o panel | Persona | Familia | Categoría | Tienda | Periodo (por defecto) | Crédito |
|---|---|---|---|---|---|---|
| K-1 Gasto del mes | — | — | — | — | S · mes en curso | — |
| K-2 Ingreso del mes | — | — | — | — | S · mes en curso | — |
| K-3 Créditos | — | — | — | — | — · hoy y mes en curso | — |
| K-4 Pagos próximos | — | — | — | — | S · vencidos hasta ayer y de hoy a hoy + 30 (PROVISIONAL) | — |
| B-1 Barómetro | — | — | — | — | S · mes en curso | — |
| G-1 Ingreso, gasto y pagos | — | — | — | — | S · 11 meses anteriores y el en curso (PROVISIONAL) | — |
| G-2 Deuda viva por línea | — | — | — | — | S · hoy | — |
| G-3 Presupuesto por categoría | — | — | V | — | S · mes en curso (sugerencia: 3 meses cerrados) | — |
| G-4 = T-2 · G-5 = T-3 | — | — | — | — | — · las de T-2 y T-3 | — |
| G-6 Tendencia del barómetro | — | — | — | — | S · 6 anteriores, el en curso y 5 siguientes (PROVISIONAL) | — |
| T2 Flujo de la quincena (Más del hogar) | — | — | — | — | S · quincena en curso (K8) | — |
| T4 Comidas de la semana (Más del hogar) | — | — | — | — | S · semana en curso (K8) | — |
| T5 Mejor canasta (Más del hogar) | — | — | — | — | — · lista vigente (`week` según K38) | — |
| T6 Ahorro y metas (Más del hogar) | — | — | — | — | — · mes en curso | — |
| D-1 Presupuesto por categoría | — | — | V | — | S · mes (o quincena) | — |
| D-2 Quincena en curso | — | — | — | — | — · quincena en curso | — |
| D-3 Pagos y disponible | S | S | — | — | S · vencidos hasta ayer y de hoy a hoy + 60 | S |
| D-4 Comidas de hoy | S | — | — | — | S · semana en curso | — |
| D-5 Diferencias de precio | — | — | — | V | — · lista vigente (`week`) | — |
| D-6 Metas | — | — | — | — | — · hoy | — |
| D-7 Participación de Karen | — | — | — | — | — · mes en curso | — |
| T-1 Plan contra real | — | — | S | — | — · 12 meses | — |
| T-2 Devengado contra efectivo | — | — | — | — | — · mes en curso a + 11 | — |
| T-3 Deuda viva proyectada | — | — | — | — | — · 12 meses | S |
| T-4 Costo semanal y kcal | S | — | — | — | — · 12 semanas | — |
| C-1 Calendario | — | S | — | — | S · mes (agenda: hoy a + 30) | S |

Los bloques del tablero no llevan filtros de la barra (son la vista de la casa; G-3 solo ordena y recorta lo visible); el filtro vive en Detalle, Tendencias y Calendario. El enlace «Ver todos los gastos con estos filtros» de Detalle lleva los seis filtros a P-02.

---

## 7. Estados por tarjeta y panel

Reglas comunes en §1.3 (carga con esqueleto y aviso a los 10 s; error con «Reintentar»; sin conexión con lo último cargado y su hora). Celdas = texto exacto o regla.

| Bloque | Vacío (qué se ve · acción) | Carga | Error (texto · acción) | Sin conexión |
|---|---|---|---|---|
| K-1 | «Aún no hay presupuesto para [mes]» · «Capturar presupuesto»; sin mes anterior, «Sin datos de [mes]» | esqueleto de tarjeta con barra | «No se pudo cargar el gasto del mes. Revisa que estés en la red de casa.» · «Reintentar» | último dato + «Cargado a las [HH:MM]» |
| K-2 | «Falta tu ingreso quincenal» · «Capturar ingreso» | esqueleto de tarjeta | «No se pudo cargar el ingreso del mes. …» · «Reintentar» | ídem |
| K-3 | «Aún no hay créditos» · «Agregar crédito»; «Sin pagos en [mes]» | esqueleto de tarjeta + bloque de disponible | «No se pudieron cargar los créditos. …» · «Reintentar» | ídem |
| K-4 | «Sin pagos en los próximos 30 días» / «Aún no hay créditos» · «Agregar crédito» | 3 filas esqueleto | «No se pudieron cargar los pagos. …» · «Reintentar» | ídem |
| B-1 | «Falta tu ingreso del mes para calcular la capacidad» · «Capturar ingreso»; sin créditos, 0 % con «Sin créditos activos» | arco esqueleto en `ui.skeleton.on-card` sin aguja ni marcas falsas | «No se pudo calcular la capacidad de endeudamiento. …» · «Reintentar» | ídem |
| G-1 | «Aún no hay meses con ingreso ni gasto» · «Capturar ingreso» | ídem T-1 | ídem T-1 | ídem |
| G-2 | «No hay deudas vivas» (sin acción) | ídem T-1 | ídem T-1 | ídem |
| G-3 | «Aún no hay presupuesto por categoría» · «Capturar presupuesto» | ídem T-1 | ídem T-1; sin sugerencias si falla `getBudgetSignals` | ídem |
| G-4, G-5 | los de T-2 y T-3 | ídem T-1 | ídem T-1 | ídem |
| G-6 | «Falta tu ingreso para calcular la capacidad» · «Capturar ingreso» | ídem T-1 | ídem T-1 (comparte petición con G-1) | ídem |
| T2 | «Falta tu ingreso quincenal» · «Capturar ingreso» | esqueleto de tarjeta | «No se pudo cargar el flujo de la quincena. …» · «Reintentar» | ídem |
| T4 | «Aún no hay plan para esta semana» · «Armar plan» | esqueleto de tarjeta | «No se pudo cargar el plan de comidas. …» · «Reintentar» | ídem |
| T5 | «Aún no hay precios para comparar» · «Registrar precios» | esqueleto de tarjeta | «No se pudo comparar la canasta. …» · «Reintentar» | ídem |
| T6 | «Aún no tienes fondo ni metas de ahorro» · «Crear meta» | esqueleto de tarjeta | «No se pudo cargar el ahorro. …» · «Reintentar» | ídem |
| D-1 | «Aún no hay presupuesto por categoría» · «Capturar presupuesto» | 5 filas esqueleto | «No se pudo cargar el presupuesto por categoría. …» · «Reintentar» | ídem |
| D-2 | «Falta tu ingreso quincenal» · «Capturar ingreso» | lista esqueleto de 7 propiedades | «No se pudo cargar la quincena. …» · «Reintentar» | ídem |
| D-3 | «Sin pagos pendientes» · «Agregar crédito» si no hay créditos | 5 filas esqueleto + bloque de disponible | «No se pudieron cargar los pagos. …» · «Reintentar» | ídem |
| D-4 | «Hoy no hay comidas en el plan» · «Armar plan» | 4 filas esqueleto | «No se pudieron cargar las comidas de hoy. …» · «Reintentar» | ídem |
| D-5 | «Aún no hay precios para comparar» · «Registrar precios» | 5 filas esqueleto | «No se pudo comparar la canasta. …» · «Reintentar» | ídem |
| D-6 | «Aún no tienes metas de ahorro» · «Crear meta» | 3 filas esqueleto | «No se pudieron cargar las metas. …» · «Reintentar» | ídem |
| D-7 | «Karen no tiene compras ni pagos este mes» (nota fija visible) | bloque esqueleto | «No se pudo cargar la participación de Karen. …» · «Reintentar» | ídem |
| T-1 | «Aún no hay presupuesto capturado» · «Capturar presupuesto» | esqueleto del alto final en `ui.skeleton.on-card`, sin ejes falsos (DESIGN §7.7) | «No se pudo cargar la gráfica. …» · «Reintentar» | ídem |
| T-2 | «Aún no hay gastos ni pagos para comparar» · «Registrar gasto» | ídem T-1 | ídem T-1 | ídem |
| T-3 | «No hay deudas vivas» (sin acción) | ídem T-1 | ídem T-1 | ídem |
| T-4 | «Aún no hay semanas con plan de comidas» · «Armar plan» | ídem T-1 | ídem T-1 | ídem |
| C-1 | «No hay eventos en este periodo» / «Aún no hay pagos ni ingresos programados» · «Capturar ingreso» | rejilla o agenda esqueleto | «No se pudo cargar el calendario. …» · «Reintentar» | ídem; meses nunca cargados: «Sin conexión · Esta página aún no se ha cargado en este teléfono» |

«…» abrevia en esta tabla la segunda frase común: «Revisa que estés en la red de casa.» Los mensajes de error nunca muestran códigos ni texto en inglés; con 401 se va a P-00 con aviso (sesión vencida, U1 §2.4) y con 404 de un registro, «Ya no existe · Volver».

---

## 8. Diseño por tamaño

Móvil primero: se diseña a 360 px y se amplía; referencias de U5: Android 412 × 915 y escritorio 1440 × 900 (U1 §6.1). Rejilla de DESIGN §4.2. Ningún ancho tiene desplazamiento horizontal de la página; solo la fila de pestañas se desplaza. Todo objetivo táctil ≥ `size.touch.min` (44 × 44 px) con `space.2` entre objetivos contiguos.

| Elemento | 360 y 412 (4 columnas, margen `space.4`) | 768 y 1024 | 1440 (12 columnas, margen `space.12`) |
|---|---|---|---|
| Tablero (Resumen) | 1 columna en el orden de apilado de §2.4 (K-1, B-1, K-4, K-2, K-3, G-3, G-1, G-6, G-2, G-5, G-4, «Más del hogar»), relleno de tarjeta `space.4`. A 360: título, valor principal y 1 línea secundaria; a 412 suma 1 dato (U1 §6.2). La cifra principal y el dígito del barómetro nunca se truncan (U1 §6.1 regla 5) | 768: K-1 con B-1 y K-2 con K-3 (4 + 4), K-4 y cada gráfica a todo el ancho, «Más del hogar» 2 × 2 (relleno `space.6`); 1024: rejilla de 1440 con la barra lateral fija de `layout.sidebar.width` y «Más del hogar» 2 × 2 | K-1 (6) y B-1 (6); K-2, K-3 y K-4 (4 cada uno); gráficas de dos en dos (6 + 6); «Más del hogar» en 4 tarjetas de 3 columnas; medianil `space.6`; todas las líneas secundarias visibles |
| Acción primaria de página | «+» de la barra inferior (`size.bottombar.item`) | 768: «+ Registrar gasto» en la barra superior (`layout.topbar.height`); igual en 1024 | igual que 1024 |
| Pestañas | fila con desplazamiento horizontal, activa visible | fila completa | fila completa |
| Detalle | paneles a 1 columna; tablas compactas de 2 líneas (título y cifra; 1 propiedad debajo; a 412, 1 más); filtros en hoja inferior con chips | 768: 1 columna con tablas de hasta 5 columnas; 1024: 2 columnas de paneles | 2 columnas de paneles; D-7 a todo el ancho al final |
| Tendencias y gráficas del tablero | 1 gráfica por fila, alto reducido, leyenda debajo, máximo 6 series, valor al tocar; etiqueta del periodo actual en texto | 1 gráfica por fila con leyenda a la derecha y valor al pasar el puntero | 2 gráficas por fila (T-1 con T-2; T-3 con T-4) |
| Calendario | agenda de hoy a hoy + 30 (412: mes compacto opcional) | mes | mes con panel del día seleccionado |
| Tipografía | `font.size.title-2` cifra; `font.size.small` en `ui.text.secondary` para la etiqueta de tarjeta y el título de gráfica (DESIGN §7.6, §7.7); `font.size.title-3` título de panel de Detalle; `font.size.small` secundario; `font.size.caption` metadatos | igual | igual |
| Alto de controles | `size.control.coarse` (puntero grueso) | `size.control.coarse` si es táctil; `size.control.fine` con puntero fino | `size.control.fine` (puntero fino) |

Teléfono en horizontal: rige el ancho resultante (915 px usa la regla de 768, U1 §6.1 regla 7).

---

## 9. Participación de Karen (reglas)

- Karen **no tiene cuenta y no ve el dashboard** (decisiones §A). Solo Axel y Ximena, co-titulares, ven Inicio, y ven lo mismo (U1 §0.5).
- Aparece **solo en lo que comparte**: (a) la despensa en común y (b) sus compras con el crédito Kueski y sus pagos. En Inicio eso es: la parte de Karen en K-3 (pagos del mes), en K-4 (pagos próximos), en la línea informativa «Con la parte de Karen: [n] %» de B-1 (desde 768), en la tarjeta 2 cuando aplica, en D-3 y en los pagos del Calendario (incluido el estado «Parcial · Falta la parte de Karen», §1.4); el panel D-7; las columnas de su parte en la tabla de T-2 («Ver como tabla»: devengado y efectivo de su parte de la despensa en común, `accrued_external_cents` y `cash_external_cents`, informativas); el valor «Karen» del filtro de familia (K5), que incluye lo conjunto.
- Su parte es **informativa** para el flujo de Axel: no suma a K-1, K-2, al dígito de B-1, a G-1, G-3, G-6 ni a las tarjetas 2 y 6, que usan la parte del hogar (U1 §0.2).
- **No hay indicadores de las finanzas de su familia**: ni ingresos, ni presupuesto, ni ahorro, ni gráficas propias; las gráficas por persona (T-4) solo tienen a Axel y Ximena (plan de comidas y salud: decisiones §A).
- Su estado de cuenta PDF es una acción del hogar desde P-25 (`generateExternalStatement`) y **Axel lo reenvía a mano** (decisiones §A, F).

---

## 10. Color, tipografía y accesibilidad

- Solo tokens de C2 / DESIGN por nombre: fondo y tarjeta `color.neutral.50` / `color.neutral.100` (oscuro `color.neutral.950` / `color.neutral.900`); texto `color.neutral.900` / `color.neutral.700` (oscuro `color.neutral.100` / `color.neutral.300`); enlaces `color.brand.700` / `color.brand.300`; roles `color.role.{income,expense,savings,debt,warning,error,success,info}.{light,dark,surface-light,surface-dark}`; personas `color.person.{axel,ximena,karen,shared}.*`; tiendas `color.store.{walmart,bodega-aurrera,sams-club,market}.*`; series `color.chart.1` a `color.chart.8` y «Otros» `color.neutral.600` / `color.neutral.400`. Ningún color nuevo (DESIGN §2.10 regla 4).
- El color nunca es el único portador de significado: cifras con signo e icono (DESIGN §3.3), estados con icono y texto, series con forma, patrón y etiqueta, personas y tiendas con nombre (K41). Ingreso y éxito no comparten componente: ingreso solo en cifras con «+», éxito solo en insignias con palomita (DESIGN §2.4).
- Cada gráfica tiene «Ver como tabla» con las mismas cifras, orden fijo de series y etiquetas de dinero en texto (§4).
- Contraste: los pares de U3 son los de DESIGN §2.8 y §2.9 (202 de 202 cumplen) y, desde 2.0, los del barómetro de DESIGN §8.9 (`ui.role.success`, `ui.role.warning` y `ui.role.error` sobre `ui.surface.card` y `ui.viz.track`, y el texto de las insignias sobre su superficie: todos ≥ 3:1 de interfaz o ≥ 4.5:1 de texto en claro y en oscuro). El rojo del barómetro es la excepción de K54 a DESIGN §2.4 («error, solo en validación y avisos»).
- Movimiento: las cifras no se animan (DESIGN §6.4); con movimiento reducido las gráficas aparecen sin transición.
- Teclado (escritorio): las pestañas se recorren con flechas; cada bloque del tablero es enfocable en el orden de apilado de §2.4 en todos los anchos.

---

## 11. Paso a implementación (F2, WF-3)

### 11.1 Componentes de DESIGN que usa U3
«Tarjeta de métrica» (§1.2), «Barómetro» (DESIGN §8.9, B-1), «Pestañas», «Insignia de rol», «Enlace», «Botón primario» (solo «Registrar gasto»), «Botón secundario» («Reintentar»), «Esqueleto», «Franja sin conexión», «Barra de vista con filtros» y «Hoja inferior de filtros», «Chip de filtro», «Tabla compacta», «Lista de propiedades», «Barra de avance», «Gráfica» (barras agrupadas, columnas agrupadas, barras horizontales, barras horizontales en % del plan, área apilada, líneas con líneas de referencia) con «Ver como tabla», «Calendario» (mes y agenda) y «Vista previa de registro». Los nombres siguen DESIGN §7 y §8; donde U3 se aparta de su contenido lo dice en el lugar (destino de la tarjeta en §1.2, «Sin marcar como pagado» y prioridad de estados en §1.4, subtítulo de T-2, línea de meta de T-4b, resumen de las gráficas sin variaciones calculadas en §4) y lo pide a DESIGN en §15.

### 11.2 Criterios de aceptación para F2
Pruebas con Playwright 1.64 + axe-core, reloj simulado en `America/Mexico_City` y datos de la semilla y el oráculo.

| Id | Criterio | Cómo se prueba |
|---|---|---|
| A-01 | El Resumen es el tablero de §2: K-1 a K-4, B-1, G-1 a G-6 y «Más del hogar» con las tarjetas 2, 4, 5 y 6, en el orden de apilado de §2.4 a 360 y 412 y en la rejilla de §2.4 a 768, 1024 y 1440; el orden del foco es el de 360 en todos los anchos | contar bloques y leer su orden en el DOM en los 5 anchos |
| A-02 | Cada tarjeta llama solo a su operación de §2 con los parámetros indicados; Inicio lee `getAlertSettings` y `getDashboardSettings` una sola vez cada una (§1.2); D-3 hace sus dos peticiones de `listScheduledPayments` (§3) y el Calendario la suya junto a `listCalendarEvents` (§5) | registro de peticiones de red por tarjeta y panel |
| A-03 | Tarjeta 4 y K38: con reloj en Dom 11 oct 2026 la tarjeta pide `period` = `2026-W41` y «Ver plan» abre P-10 en `2026-W42`; con Mié 14 oct 2026, ambas en `2026-W42` | dos corridas con fecha simulada |
| A-04 | Tarjeta 5 pide `week` de la lista vigente (Dom 11 oct 2026 → `2026-W42`) | petición de red |
| A-05 | K-3 y K-4 con la semilla el 10-oct-2026: K-4 con la fila «Vie 30 oct 2026 · Kueski · $973.75 · Pendiente» y «Hogar $598.82 · Karen $374.93»; K-3 con «Pagos de octubre: $973.75 · 1 pago», «Disponible hoy $6,905.50 · Lectura B» e insignia «Pendiente confirmar»; con lectura A, $10,000.00 | captura y texto |
| A-06 | K-1 en octubre de 2026 (= P1, regla de ejemplos al inicio) muestra gastado $1,197.64, pagado en efectivo $598.82 y «Sin datos de septiembre»; en noviembre de 2026, «Igual que octubre» y «$598.82 más que octubre» en efectivo; sin presupuesto, «—» y «Capturar presupuesto» | texto con dos relojes simulados |
| A-07 | Cada bloque muestra vacío, carga, error y sin conexión con los textos de §7; el aviso de 10 s aparece en la carga lenta | respuestas simuladas (vacía, lenta, 500, sin red) |
| A-08 | Una sola acción con relleno de marca en la página («Registrar gasto»); cada tarjeta, un solo enlace | conteo de botones primarios y enlaces por tarjeta |
| A-09 | Objetivos táctiles ≥ 44 × 44 px y sin desplazamiento horizontal de página en los 5 anchos | medición de cajas y del ancho de documento |
| A-10 | La cifra principal más larga y el dígito del barómetro caben completos a 360 px | prueba de ancho de U1 §6.1 regla 5 |
| A-11 | Gráficas con orden fijo de series, marcador o patrón distinto por serie, etiquetas de dinero en texto y «Ver como tabla» con las mismas cifras; T-2 con reloj en oct-2026: fila oct-2026 con $1,947.50, $973.75 y $973.75 (R-45); T-2 con reloj en nov-2026 (la ventana empieza en nov-2026): fila nov-2026 con $1,947.50, $1,947.50 y $973.75 (R-46); capas (Q2-65 A) en nov-2026, efectivo acumulado $2,921.25 y devengado acumulado $3,895.00, y en mar-2027, $13,103.78 y $14,362.48 (R-52, R-51), las mismas con los dos relojes; la fila «Total» de T-2 suma las 12 filas visibles; con un mes de diferencia negativa la capa 2 vale cero con «Pagado por adelantado»; ninguna marca con contorno y ninguna gráfica con un número en cada punto; cada gráfica en `<figure>` con `<figcaption>` (título y resumen de §4), resumen sin porcentajes ni variaciones calculadas y, en T-2, sin la palabra «deuda» | inspección del DOM y de la tabla; T-2 con dos corridas de reloj simulado (oct-2026 y nov-2026) |
| A-12 | Calendario: franjas de quincena 1 al 15 y 16 al último día; febrero de 2027 con pago «Dom 28 feb 2027»; semana de lunes a domingo | fecha simulada en octubre de 2026 y febrero de 2027 |
| A-13 | Karen: solo en K-3, K-4, la línea informativa de B-1 (desde 768), D-3, D-7, Calendario, las columnas de su parte en la tabla de T-2 y el filtro de familia; ningún indicador de su familia | revisión de textos con «Karen» |
| A-14 | Formato: dinero con `tnum` y signo menos tipográfico; fechas `10-oct-2026` o con día de la semana; 0 textos en inglés | axe-core y búsqueda de textos |
| A-15 | Tema claro y oscuro con los tokens de §10; axe-core sin violaciones AA | axe-core en los dos temas |
| A-16 | D-3 el 10-oct-2026 con la semilla: 3 pagos («Vie 30 oct 2026», «Dom 15 nov 2026», «Lun 30 nov 2026»), ninguno vencido y sin petición de vencidos | petición de red y texto |
| A-17 | Estados de §1.4 con el pago del 30 de octubre: el Jue 29 oct 2026 «Vence mañana»; el Vie 30 oct 2026 «Vence hoy»; el Sáb 31 oct 2026 sin marcar, «Vencido» con «Sin marcar como pagado»; con solo la parte de Karen pagada, «Falta la parte del hogar»; sin respuesta de `getAlertSettings`, ningún «Por vencer» | reloj simulado y respuestas simuladas |
| A-18 | B-1 con respuestas simuladas: `ratio_pct` 30 → «Verde · dentro del rango» (`circle-check`); 31 y 40 → «Amarillo · cerca del tope» (`triangle-alert`); 41 → «Rojo · arriba del tope» (`circle-alert`); sin `income_cents`, dígito «—» y «Capturar ingreso»; `aria-valuetext` «Capacidad de endeudamiento [n] % de 40 %: [estado]»; con `warning_pct` = 25 en la respuesta, la marca del tramo se mueve a 25 % (ningún umbral constante) | respuestas simuladas y lectura del árbol de accesibilidad |
| A-19 | G-1 y G-6 salen de una sola petición a `getDashboardMonthlyOverview` (`from` = mes en curso − 11, `to` = mes en curso + 5); G-1 con 3 series en columnas, la nota fija de pagos y sin columna en meses sin datos; G-6 con las líneas de referencia en `warning_pct` y `limit_pct`, trazo 8 · 4 en proyección y la línea «Hoy» | registro de red e inspección del DOM |
| A-20 | G-3: barras en % del plan, marcas en `alerts.budget_warning_pct` y `alerts.budget_limit_pct` leídas de `getAlertSettings`, orden descendente por `used_pct`, insignia «Sugerencia: subir presupuesto» o «recortar» solo con `signal` = `raise` o `cut`; sin `getBudgetSignals`, barras sin insignias | respuestas simuladas |
| A-21 | Desglose de §2.2: barra de G-3 → P-02 con `category_id` y `period` → P-04 → P-15 con el producto → pestaña «Tendencias», en 4 toques | recorrido con la semilla |
| A-22 | Ningún bloque calcula porcentajes, diferencias ni estados: todos llegan de C1; con un campo nuevo de C1 v1.5 ausente, su línea no se pinta y el resto del bloque se ve | respuestas sin los campos nuevos |

---

## 12. Capturas para U5

Cada captura en claro y en oscuro, a **412 × 915** y **1440 × 900** (U1 §6.1), con reloj simulado:
1. Tablero con datos de la semilla el 10-oct-2026: K-1, K-3 y K-4 con cifras (A-05, A-06); K-2, B-1, G-1 y G-6 con su vacío por falta de ingreso; G-3 con su vacío si no hay plan.
2. Tablero con un ingreso **sintético** de prueba (marcado así en la captura) para ver el barómetro en verde, amarillo y rojo (A-18) y G-6 con proyección.
3. Tablero en carga (esqueletos), con error en un bloque y sin conexión con «· [n] en cola».
4. Tablero a 360 × 800 con la prueba de ancho de la cifra principal y del dígito del barómetro.
4b. «Más del hogar» el Dom 11 oct 2026 y el Mié 14 oct 2026 (tarjeta 4 y destino P-10, K38).
5. Detalle completo, incluido D-7 «Participación de Karen» con la nota fija.
6. Tendencias con T-1 a T-4 y una gráfica en «Ver como tabla».
7. Calendario: agenda (412) y mes con panel del día (1440) de octubre de 2026; mes de febrero de 2027.

---

## 13. Trazabilidad con los criterios U3

| Criterio | Dónde se cumple |
|---|---|
| U3-01 · tablero de control (K54) | §2 (K-1 a K-4, B-1, §2.1 y G-1 a G-6, §2.2, §2.3, §2.4), §2.5 con las tarjetas 2, 4, 5 y 6 (Tarjeta 4 con K38, U1 §0.2, casos domingo y entre semana; Tarjeta 5 con `second_store_id`), §11.2 A-01, A-03 y A-18 a A-21 |
| U3-02 · métricas C1 | §2 a §5 (M-xx, `operationId`, campos, fórmula y periodo en cada bloque); §2.6 con las operaciones y campos NUEVOS publicados en C1 v1.5 y las métricas M-18 a M-22 de C1-datos §9; verificación en el reporte |
| U3-03 · estados | §1.3, §7 (25 filas de bloques × vacío, carga, error, sin conexión), cada bloque de §2 con su vacío y su error |
| U3-04 · móvil primero | §2.4 (orden de apilado en 412 y rejilla en 1440), §8 (360/412, 768/1024, 1440), §1.2 (táctil 44 px) |
| U3-05 · Detalle, Tendencias, Calendario | §3, §4, §5 (propósito, contenido, métricas, filtros, estados); quincenas K8 y febrero en §5; 0 donas (§4); T-2 con las dos capas del mismo origen (Q2-65 A, acumulados de C1 v1.3) |
| U3-06 · Karen | §3.1, §9 |
| U3-07 · próximo pago y A/B | K-3 (lectura activa y «Pendiente confirmar»), K-4 (próximo pago del calendario unificado), D-3 |
| U3-08 · cifras trazables | ejemplos con R-nn o K-nn (oct-2026 a mar-2027 = P1 a P6, regla de ejemplos al inicio); marcadores sin número; verificación en el reporte |
| U3-09 · formatos | reglas al inicio, textos de §2 a §7 |
| U3-10 · filtros | §6 |
| U3-11 · tokens y gráficas | §1.2, §2.1 y G-1 a G-6, B-1 (DESIGN §8.9, tokens de rol existentes), §4 (con el resumen y el `<figcaption>` de DESIGN §7.7), §10 |
| U3-12 · implementación y supuestos | §11, §14, §15, el bloqueo `docs/estado/bloqueos/U3-t2-acumulados.md` (RESUELTO, Q2-65 A) y los cinco `docs/estado/bloqueos/U3-tablero-*.md` (PROVISIONAL de U3 2.0) |

---

## 14. Supuestos

1. Textos de rangos cortos en tarjetas («Del 16 al 31 oct», «Lun 12 al Dom 18 oct») y de meses en gráficas («oct-2026»): derivados del formato de U1 §0.4; no cambian cifras, esquema, operaciones ni tokens.
2. Hora de carga sin conexión como «[HH:MM]» en 24 horas (la misma forma que `send_time`, K43); solo texto.
3. Las tarjetas del Resumen no llevan filtros (vista de la casa); los filtros viven en Detalle, Tendencias y Calendario (§6).
4. K-3 (antes la tarjeta 3), con varias líneas revolventes, muestra la primera de `availability[]` y «+[n] líneas» (las demás en D-3).
5. (Retirado en la versión 1.1: T-2 sigue DESIGN §8.3, con la capa 2 en cero en los tramos de pago adelantado, y ya no cambia a dos líneas.)
6. Prioridad entre estados de pago (Pagado, Vencido, Por vencer, Parcial, Pendiente) y la línea «Falta la parte de …» en un pago vencido o por vencer (§1.4): DESIGN §7.8 da los estados pero no su prioridad; solo orden de presentación y texto.
7. Textos «Falta la parte del hogar» y «Vence en [d] días» (§1.4), para los casos que DESIGN §7.8 no nombra; solo texto.
8. Línea de meta de T-4b con el color de la persona y trazo de puntos 2 · 4 sin marcadores, distinta de los trazos de Axel (continuo) y Ximena (8 · 4); DESIGN §8 no fija la línea de meta; solo dibujo.
9. Redacción de los resúmenes de una frase de T-1 a T-4 (§4): DESIGN §7.7 pide el resumen pero no su texto; U3 lo arma solo con cifras de C1, sin comparaciones; solo texto.
10. (2.0) La pestaña sigue llamándose «Resumen» (decisiones §G); K54 cambia su contenido, no su nombre.
11. (2.0) La comparación contra el mes anterior se dice en dinero con «más que», «menos que» o «Igual que» y los iconos `arrow-up-right` y `arrow-down-left` del mapa de DESIGN §5.4; sin porcentaje de variación (la interfaz no calcula, U1 §0.2).
12. (2.0) La sugerencia «subir o recortar presupuesto» va con rol información (`info`) y la palabra «Sugerencia»: es una lectura del motor, no una alerta.
13. (2.0) Orden de apilado de §2.4 (gasto, barómetro, pagos, ingreso, créditos, gráficas, «Más del hogar»): sigue el propósito de K54; solo orden.
14. (2.0) Nota fija de G-1 («Los pagos de crédito pagan compras que ya cuentan en el gasto de su mes») para que las tres columnas no se lean como sumables; solo texto.

Los supuestos 1 a 9 no cambian cifras, esquema, operaciones ni tokens. Los 10 a 14 tampoco; lo que sí cambia operaciones o claves (campos nuevos de §2.6) se publicó en C1 v1.5 y en el esquema 1.5 y lo que Axel no decidió es PROVISIONAL con bloqueo.

## 15. Dependencias y abiertos
- **U3 2.0 (K54)**, dueños y pedidos:
  - C1 v1.5 (`s1-c1-v15`), HECHO: las 3 operaciones nuevas, el par `getDashboardSettings` / `updateDashboardSettings`, los 14 campos nuevos de `DashboardMonthVsPlan`, `IncomeMonth` y `CreditTotals` (§2.6) y las métricas M-18 a M-22 en C1-datos §9. Backend de las 5 operaciones HECHO en `backend/internal/api/tablero/` (corrección r2, hallazgo A-1: motores S2 y S3, sin recalcular cifras; `make verify-all` da EXIT 0). Falta que los manejadores de área llenen los 14 campos opcionales nuevos de `DashboardMonthVsPlan`, `IncomeMonth` y `CreditTotals` (hoy se omiten y la interfaz cae a su regla de §2) (H-1 de `docs/estado/wf3/REVISION-especificacion-r1.md`). Con un 501 o un campo ausente, los bloques muestran lo que ya existe (regla de §2).
  - Esquema 1.5, HECHO (migración 00005 con `@param` y bloqueo): `debt.capacity_warning_pct` = 30, `debt.capacity_limit_pct` = 40 (FIRMES, K54), `dashboard.payment_reminder_count` = 5, `budget.signal_months` = 3 y `budget.signal_margin_pct` = 10 (PROVISIONALES).
  - U4: en qué sección de Ajustes se editan las 5 claves (propuesta: las dos de `debt.*` en P-26.05 «Créditos y condiciones», las dos de `budget.signal_*` en P-26.02 «Presupuesto» y la de `dashboard.*` en «Notificaciones»), con vista previa de su efecto.
  - U1 (fuera del alcance de `u3-v2`): §2.1 y §2.2 fila P-01 («6 tarjetas»), §3.3 (N-020 a N-029 más las necesidades nuevas del tablero), §6 fila P-01 y F-4 paso 1 («Toca la tarjeta «Próximo pago»» pasa a «Toca el pago en K-4»). F-3 ruta A ya se aclaró (C3-flujos punto 2, opción c).
  - P-04 → P-15: el renglón del ticket enlaza a P-15 con su producto (paso 4 de §2.2).
  - DESIGN 1.1 §8.9 «Barómetro» y la fila nueva de §8.2 ya existen (agente `u3-v2`).
  - Frontend P-01 nuevo, capturas U5 de §12 y C3 regenerada: WF-3c.
- DESIGN §7.6 a §7.8 y §8 ya existen y U3 los sigue desde la versión 1.1. Solicitudes a DESIGN (U2b), con su razón (detalle en `docs/estado/wf1/correccion-U3-r1.md`):
  1. §7.6: la tarjeta entera lleva a la pantalla destino de su acción primaria, no a la pestaña Detalle (U3-01: acción primaria única y la tarjeta 4 a P-10); y la anatomía suma el icono de 20 px junto a la etiqueta (§5.2; §5.4 ya prevé iconos en tarjetas).
  2. §7.8: la línea «Sin marcar como pagado» junto a «Vencido» (U1 N-022, C1 M-03), la prioridad entre estados y los textos «Falta la parte del hogar» y «Vence en [d] días» (U3 §1.4, supuestos 6 y 7).
  3. §8.3: quitar del subtítulo el rótulo «Deuda viva …» o justificarlo frente a K19, N-091 y M-11 («no se rotula deuda viva»).
  4. §8: fijar la línea de meta en una gráfica de líneas (U3 T-4b usa el color de la persona y trazo de puntos 2 · 4 sin marcadores, supuesto 8).
  5. §7.7: el ejemplo de resumen («El gasto subió 12 % respecto a septiembre») pide calcular una variación que la interfaz no hace (U1 §0.2). U3 usa resúmenes que solo enuncian cifras de C1 (§4). Pedido: cambiar el ejemplo o aclarar que la variación tiene que llegar calculada desde C1.
- C1-datos §9: dos métricas para D-4 y D-5, «Comidas de hoy» (N-121 vía N-027) y «Diferencias de precio» (N-162 vía N-027), con la fórmula, el periodo y los filtros que U3 §3 escribe mientras tanto (hallazgo H-7).
- C1 (`getAccrualVsCash`, esquema `AccrualVsCash.months[]` en `openapi.yaml`, y C1-datos M-11): agregar `cumulative_accrued_cents` y `cumulative_cash_cents`, con el mismo origen que `cumulative_difference_cents` (Q2-65 A, aprobada en la puerta 2; bloqueo `docs/estado/bloqueos/U3-t2-acumulados.md` RESUELTO; C1 v1.3).
- C1: `second_store_id` en `DashboardBestBasket` para que la tarjeta 5 nombre la segunda tienda (Q2-86 A: solicitud aprobada; C1 v1.3). Mientras no llegue, «menos que la siguiente tienda» y la tienda se ve en D-5; no bloquea.
- C1 (`SavingsGoal.status`, que sale de `backend/internal/api/c1gen/gen2_schemas.py`): C1 v1.3 todavía declara `on_track`, `behind` y `achieved` («en ruta, atrasada o cumplida»). Q2-18 A dejó tres estados, `in_progress`, `achieved` y `overdue` (`docs/contratos/motores/S2-interfaz.md` 1.3 §3.13), que son los que usan la tarjeta 6 y D-6.
- C1-datos §9, M-16 «Metas» (N-111, N-027): dice «% de avance entero»; con Q2-17 C la métrica de Metas va con un decimal (`progress_tenths_pct`, ya en C1 v1.3) y el entero queda en M-07 (tarjeta 6). El bloqueo `docs/estado/bloqueos/S2-metas.md` lo pide en «Qué cambia» («alinear RN-18 con M-07 y M-16»). D-6 ya usa `progress_tenths_pct`.
- C1 (`getCashFlowFortnights`, `rows[].planned_savings_cents`): C1 v1.3 lo declara `MoneyCents` obligatorio y describe `balance_cents` como «income − egresos − ahorro». Con el bloqueo abierto `docs/estado/bloqueos/S2-ahorro-quincenal.md`, el motor S2 entrega el ahorro en `null` y no lo resta del saldo (S2-interfaz 1.3 §3.11). Pedido: que C1 marque el campo como PROVISIONAL con ese bloqueo, admita `null` y no dé por hecha la pregunta 3 en la descripción de `balance_cents`. U3 D-2 muestra «—» mientras llegue `null`.
- Bloqueos de S2 que tocan cifras del dashboard (U3 solo muestra lo que llegue y marca PROVISIONAL lo que depende de un bloqueo abierto):
  - Abiertos: `docs/estado/bloqueos/S2-flujo-quincenal.md`, solo Q2-12 (gasto variable estimado de D-2); `docs/estado/bloqueos/S2-ahorro-quincenal.md` (valor del ahorro planificado de D-2: base, metas que cuentan y si resta del saldo; el reparto de Q2-13 B ya es FIRME); y `docs/estado/bloqueos/S2-metas-vencida.md` (aportación de una meta vencida sin cumplir en D-6).
  - RESUELTOS en la puerta 2, sin marca: `docs/estado/bloqueos/S2-presupuesto-total.md` (K-1, antes tarjeta 1; Q2-19 B y Q2-20 A), Q2-13 B y Q2-14 A de `S2-flujo-quincenal` (tarjeta 2; Q2-13 B fija solo el reparto del ahorro planificado, cuyo valor en D-2 sigue abierto, arriba), `docs/estado/bloqueos/S2-metas.md` (tarjeta 6 con el avance entero y D-6 con un decimal; Q2-16 A, Q2-17 C y Q2-18 A) y `docs/estado/bloqueos/S2-proyeccion-interruptor.md` (T-2, meses 7 a 12 con el interruptor encendido; Q2-21 a Q2-24 C).
  - Ninguno cambia la estructura, los textos fijos ni las operaciones de U3.
- Bloqueo propio de U3: `docs/estado/bloqueos/U3-t2-acumulados.md` (capas de T-2, sus etiquetas finales y su resumen): RESUELTO en la puerta 2 con Q2-65 A, una sola variante. El bloqueo `S1-parametros-avisos` (RESUELTO, Q2-11 B) no toca el dashboard: U3 no lee ni escribe `send_time`, `advance_days`, `weekday` ni `day_of_month`; los días de «Por vencer» son `alerts.payment_due_days_before`, ya decididos (decisión F, K44).

## FIN: docs/design/U3-dashboard.md

## ARCHIVO: docs/design/DESIGN.md


(EXTRACTO: §1.2, 1.3, 2.4, 3.2, 3.3, 4.1, 4.3, 7.6, 7.7, 7.14, 8.9 y 9. El archivo completo está en la rama main.)

### 1.2 Principios

| # | Principio | Qué significa en la interfaz | Cómo se comprueba |
|---|---|---|---|
| P1 | **Completo pero minimalista** | Todo dato existe, pero cada pantalla muestra solo lo que pide su decisión; lo demás queda a un paso (*progressive disclosure*): «Mostrar N más», menú «⋯», vista previa del registro (U1 G-06). El Resumen del dashboard lleva ≤ 6 tarjetas (decisiones §G). | Revisión por pantalla en U5: ningún dato secundario compite con la cifra principal. |
| P2 | **Una acción primaria por pantalla** | Un solo botón con relleno de marca visible a la vez (§2.9); las demás acciones son secundarias, de texto o van en «⋯». En móvil la primaria nunca se esconde (U1 §1.3). | Conteo de botones primarios por pantalla = 1 (las hojas y diálogos cuentan como pantalla propia). |
| P3 | **Móvil primero** | Se diseña a 360 px y se amplía (U1 §6.1). Registrar un gasto = 3 toques (K10). Todo control táctil mide ≥ 44 × 44 px (§4.6). PWA instalable en Android (decisiones §G). | U5 a 360 y 412 px; conteo de toques de F-1. |
| P4 | **El color significa, no decora** | Regla 60-30-10 y saturación apagada (§2.2); el color de un rol, persona, tienda o serie **nunca va solo**: lo acompaña signo, icono, forma o texto (§2.10). | Vista en escala de grises sin pérdida de significado. |
| P5 | **Cifras exactas y legibles** | Dinero `$1,258.68` (U1 §0.4) con cifras tabulares (§3.4); nunca se trunca ni se parte (U1 §6.1 regla 5); nunca se anima (§6.4). | `$99,999.99` completo en la tarjeta más estrecha a 360 px. |
| P6 | **Calma** | Densidad cómoda, sombras leves solo en claro (§4.5), movimiento corto (§6), sin rojos ni destellos salvo error y acción crítica. | Revisión de U5 contra §2 y §6. |
| P7 | **Claro y oscuro de primera** | Tema por defecto según el sistema, con interruptor (decisiones §G); cada token de color tiene variante clara y oscura. | Toda captura de U5 existe en los dos temas. |
| P8 | **Accesible por cálculo** | WCAG 2.2 AA mínimo, calculado: texto ≥ 4.5:1; texto grande, interfaz y gráficas ≥ 3:1 (decisiones §G). | Dos tablas oficiales, 0 bajo umbral: `contraste.csv` (paleta: 202 pares, §2.8) y `contraste-componentes.csv` (componentes: 344 combinaciones, §10.5). |

### 1.3 Tono de voz

Texto de interfaz en español de México, claro y breve, para una familia (Axel, Ximena y Karen), no para un banco.

- **Trato de tú** y frases cortas. Botones en infinitivo: «Guardar», «Registrar gasto», «Cambiar semana».
- **Sin jerga** financiera ni técnica: «Lo que debes hoy», no «saldo insoluto»; los términos del dominio son los del diccionario de datos.
- **Cifras con su formato**: `$1,258.68`, `10-oct-2026`, `80 %`, tasas como se capturaron (`17.27 %`; K12).
- **Mensajes que dicen qué pasó y qué hacer**: «Falta capturar la tasa de esta tarjeta», no «Error de validación».
- **Sin signos de exclamación ni emojis** en mensajes del sistema; el éxito es discreto: «Guardado · Deshacer» (U1 G-07).
- Textos ya fijados en U1 y decisiones que se respetan tal cual: «Sin conexión con el servidor de casa», «faltan restricciones» (K16), «falta capturar» (K18).

| Bien | Evitar | Por qué |
|---|---|---|
| «Guardado · Deshacer» | «¡Gasto guardado exitosamente!» | Breve, sin exclamación, con salida. |
| «Te quedan $1,320.50 de despensa este mes» | «Disponible: 1320.5» | Dice qué es y usa el formato de dinero. |
| «Este pago vence mañana» | «Pago próximo a vencer (D-1)» | Sin códigos internos. |

### 2.4 Roles semánticos

Texto e icono del rol: paso 700 en claro y 400 en oscuro (tokens `light` y `dark`); insignia: paso 100 en claro y 900 en oscuro (`surface-light` y `surface-dark`). El paso 600 del plan no llega a 4.5:1 en claro con ningún matiz, por eso el texto usa 700 (colorimetría §6.2). La escala completa 50–950 de cada rol está en `color.role.<rol>.*` para los estados vecinos (hover: 800 en claro, 200 en oscuro).

| Rol | Matiz | Claro (texto) | Oscuro (texto) | Refuerzo obligatorio (nunca solo color) |
|---|---|---|---|---|
| Ingreso | 135° | `color.role.income.light` · `#457723` | `color.role.income.dark` · `#98CA7B` | Signo «+» y flecha de entrada (`arrow-down-left`) en cifras. |
| Gasto | 45° (terracota) | `color.role.expense.light` · `#A04B20` | `color.role.expense.dark` · `#F8A078` | Signo «−» y flecha de salida (`arrow-up-right`) en cifras. |
| Ahorro | 200° | `color.role.savings.light` · `#00767B` | `color.role.savings.dark` · `#3FD0D7` | Icono `piggy-bank` y etiqueta. |
| Deuda / crédito | 310° | `color.role.debt.light` · `#7E509D` | `color.role.debt.dark` · `#D0A3F2` | Icono `credit-card` y etiqueta. |
| Advertencia | 75° | `color.role.warning.light` · `#8B5D00` | `color.role.warning.dark` · `#E5AE5B` | Icono `triangle-alert` y texto. |
| Error | 20° (rojo) | `color.role.error.light` · `#A34449` | `color.role.error.dark` · `#FC9999` | Icono `circle-alert` y mensaje; solo en validación y avisos. |
| Éxito | 150° | `color.role.success.light` · `#207B3E` | `color.role.success.dark` · `#7ECF90` | Palomita (`circle-check`) y texto («Guardado», «Pagado»); nunca en cifras. |
| Información | 240° | `color.role.info.light` · `#006FA4` | `color.role.info.dark` · `#6BC3FD` | Icono `info` y texto. |

Decisiones de matiz que este resumen respeta:

- **Gasto en terracota 45°** (K26: rango 40–45°, error se queda en 20°). De los seis grados del rango, 45° es el más separado de error (24.9°) y de advertencia (30.7°) medidos en el hex (colorimetría «Cambios 1.1», punto 1).
- **Ingreso**: decisiones §I lo deja en **150°** «o U2 propone 135°». Valor vigente en los tokens: **135°** (propuesta de U2a ya prevista por decisiones §I; ratificada por Axel en la puerta 2, Q2-81 A, K51). Motivo: separa ingreso de éxito 15° en vez de 0° y forma con advertencia el par de acentos de la armonía. Aun así, lado a lado ingreso y éxito apenas se distinguen (ΔE ≈ 3), por eso ingreso vive en cifras con «+» y éxito solo en avisos con palomita (colorimetría §6.4).
- **Gasto frente a error**: ya no colisionan (≥ 24.9° de matiz), pero con daltonismo rojo-verde siguen cerca: el gasto solo aparece como cifra con «−», el error solo como mensaje con icono; nunca comparten componente (colorimetría §6.5).

### 3.2 Escala

Seis tamaños (12/14/16/20/24/32), la misma escala en todos los anchos; base 16 px = 1 rem. Interlíneas múltiplos de 4 para caer en la retícula de 4 pt (§4.1).

| Token | px | rem | Interlínea | Peso | Uso |
|---|---|---|---|---|---|
| `font.size.caption` | 12 | 0.75 | 16 px | 400 (500 en etiquetas) | metadatos, etiquetas de eje, pies de tabla, insignias, fechas secundarias |
| `font.size.small` | 14 | 0.875 | 20 px | 400 | interfaz densa: barra lateral, tablas, propiedades, menús, botones en escritorio |
| `font.size.body` | 16 | 1 | 24 px | 400 | texto de páginas y notas, campos de formulario, botones en móvil |
| `font.size.title-3` | 20 | 1.25 | 28 px | 600 | títulos de sección y de tarjeta, cifra secundaria |
| `font.size.title-2` | 24 | 1.5 | 32 px | 600 | cifra principal de una tarjeta de métrica, título de hoja o diálogo |
| `font.size.title-1` | 32 | 2 | 40 px | 700 | título de página |

| Token | Valor | Uso |
|---|---|---|
| `font.weight.regular` | 400 | texto corrido y celdas |
| `font.weight.medium` | 500 | etiquetas, elemento activo, botones |
| `font.weight.semibold` | 600 | títulos 20 y 24, cifras principales |
| `font.weight.bold` | 700 | solo título de página |
| `font.letter-spacing` | 0 | el tamaño óptico ya corrige el espaciado |

Decisiones:

- **Nombres por papel, no por tamaño de Tailwind**: `caption … title-1` evita confundirlos con la escala por defecto de Tailwind (donde `text-lg` es 18 px) y deja la escala cerrada en seis valores.
- **Campos de formulario a 16 px** también en móvil: el navegador no amplía la página al enfocar y el monto de F-1 se lee de un vistazo.
- **Tablas y barra lateral a 14 px**: es la densidad de Notion en vistas de base de datos y cabe en 360 px con la cifra completa.
- **Máximo tres tamaños por pantalla además de `caption`**: la jerarquía la dan el peso y el color neutro (900 frente a 700), no muchos tamaños.

### 3.3 Cifras y dinero

| Regla | Valor | Justificación |
|---|---|---|
| Cifras tabulares en **todo** dinero (también porcentajes, cantidades y fechas en tablas) | `font-variant-numeric: tabular-nums` (rasgo `tnum`, presente en los dos archivos de Inter) | Columnas alineadas y cifras que no «bailan» al cambiar. |
| Formato | `$1,258.68`: MXN, símbolo `$`, coma de miles, punto y dos decimales (U1 §0.4) | Formato ya decidido. |
| Alineación en tablas | a la derecha; con dos decimales fijos y `tnum` los puntos decimales quedan en columna | Comparar montos de un vistazo. |
| Ingreso | `+$1,258.68`, `color.role.income.*` e icono `arrow-down-left` | Signo, color e icono: nunca solo color. |
| Gasto o negativo | `−$1,258.68`: signo menos tipográfico U+2212 (presente en Inter) antes de `$`, sin espacio, color del rol e icono | El guion corto (U+002D) se confunde con un guion de texto. |
| Negativo sin rol propio (p. ej., restante rebasado) | `−$230.00` con `color.role.warning.*`, icono `triangle-alert` y la palabra «rebasado» | Advertencia es el estado de un límite superado; error queda para validación (§2.4). |
| Montos neutros (totales, saldos sin rol) | sin signo, `color.neutral.900` (claro) o `color.neutral.100` (oscuro) | El color se reserva para significado (P4). |
| Sin truncar | la cifra nunca lleva puntos suspensivos ni se parte en dos líneas; se quita primero lo secundario (U1 §6.1 regla 5) | `$99,999.99` completo en la tarjeta más estrecha a 360 px; lo verifica U5. |
| Sin animación | las cifras no cuentan hacia arriba ni se deslizan (§6.4) | Un monto siempre es exacto en pantalla. |

Otros formatos con la misma tipografía: fechas `10-oct-2026` y `Vie 30 oct 2026` (U1 §0.4); porcentajes de avance enteros `80 %`; tasas como se capturaron, `17.27 %` y `25.8685 %` (K12).

### 4.1 Espaciado (retícula de 4 pt)

Siete valores y el cero; ningún margen, relleno o separación usa otro valor (los grosores de borde de §4.5 y las medidas de §4.3 y §4.6 son medidas de componente, no espaciado). En Tailwind 4, `--spacing: 0.25rem` hace que `p-1 … p-12` coincidan con estos tokens.

| Token | px | rem | Uso principal |
|---|---|---|---|
| `space.0` | 0 | 0 | sin separación |
| `space.1` | 4 | 0.25 | icono y texto dentro de un control; entre insignias |
| `space.2` | 8 | 0.5 | entre controles relacionados; etiqueta y campo; separación mínima entre objetivos táctiles |
| `space.3` | 12 | 0.75 | relleno horizontal de botones y campos; entre elementos de una lista |
| `space.4` | 16 | 1 | margen lateral y medianil en móvil; relleno de tarjeta en móvil |
| `space.6` | 24 | 1.5 | relleno de tarjeta desde 768; medianil desde 1024; entre bloques de una página |
| `space.8` | 32 | 2 | margen lateral a 1024; entre secciones de una página |
| `space.12` | 48 | 3 | margen lateral a 1440; aire sobre el título de página |

### 4.3 Medidas de estructura

| Token | Valor | Uso y justificación |
|---|---|---|
| `layout.sidebar.width` | 240 px | barra lateral fija desde 1024 (U1 G-01): cabe el nombre de cada sección a 14 px con icono de 20 px; colapsada se oculta por completo (Ctrl+\, U1 §6.2) |
| `layout.sidebar.overlay-width` | 280 px | barra lateral superpuesta a 768 (U1 §6.2 deja «ancho fijo (U2)»): un poco más ancha porque se usa con el dedo |
| `layout.sidebar.mobile-width` | 85 % | panel desde ☰ a 360 y 412 (fijado en U1 §6.2) |
| `layout.topbar.height.mobile` · `layout.topbar.height.desktop` | 56 px (< 768) · 48 px (≥ 768) | barra superior de página (U1 G-02): en móvil aloja botones de 44 px con 6 px de aire |
| `layout.bottombar.height` | 56 px + `env(safe-area-inset-bottom)` | barra inferior de 5 destinos (U1 G-05); respeta la barra de gestos de Android |
| `layout.reading.max-width` | 720 px | páginas de texto y notas a partir de 1440 (U1 §6.2) |
| `layout.preview.width.desktop` · `layout.preview.width.wide` | 400 px (1024) · 480 px (1440) | vista previa lateral de registro (U1 G-06, «más ancho» a 1440) |
| `layout.dialog.width.short` · `layout.dialog.width.long` | 480 px · 640 px | diálogo corto (confirmación) y largo (formularios), desde 768; en < 768 son hojas inferiores (U1 §1.3) |
| `layout.command.width` | 640 px | panel centrado del menú `/` y búsqueda desde 768 (U1 G-03, G-04) |

### 7.6 Tarjetas de métrica

**Anatomía.** Etiqueta (`font.size.small`, `ui.text.secondary`, sin dos puntos); cifra principal (`font.size.title-2`, `font.weight.semibold`, `tabular-nums`; neutra en `ui.text.default` o de rol con signo e icono, §3.3); dato secundario opcional (uno a 360, uno más a 412; comparación contra un periodo con nombre, «12 % más que septiembre», con signo, icono y color de rol); y opcionalmente una **barra de avance** o una mini serie (12 puntos, serie 1, sin ejes, desde 768). Desde K54 el Resumen de Inicio es un tablero de control (U3 2.0) sin tope de tarjetas y con el orden de U3 §2.4; antes llevaba ≤ 6 tarjetas en orden fijo (decisiones §G).

**Barra de avance** (presupuesto, metas): pista `ui.viz.track` de 8 px con `radius.full`; llenado `ui.viz.meter` bajo el primer umbral; desde el primer umbral de alerta (`alerts.*` de `app_setting`, 80 % inicial y configurable, K44) el llenado pasa a `color.role.warning.light` / `.dark` con `triangle-alert` y el texto «80 % usado»; desde el segundo (100 % inicial) se agrega «rebasado» y el excedente `−$230.00` (§3.3). El llenado nunca pasa del ancho de la pista: el excedente se dice en texto. Los umbrales se leen de Ajustes, nunca se fijan en el componente.

| Estado | Tokens y comportamiento |
|---|---|
| Reposo | fondo `ui.surface.card`, `shadow.0`, `radius.md` |
| Hover | borde de 1 px `ui.border.control-hover` (la tarjeta entera es enlace al Detalle) |
| Foco | anillo común en la tarjeta |
| Activo | presionado: borde `ui.border.selected`, sin cambiar el fondo |
| Deshabilitado | no aplica |
| Error | «No se pudo cargar · Reintentar» dentro de la tarjeta, sin cifra (U1 §2.4) |
| Carga · vacío · falta capturar | esqueleto de etiqueta y cifra en `ui.skeleton.on-card`; «Sin datos este mes» con su acción; «Falta capturar · Capturar» (K18) |

**Tamaños.** Relleno `space.4` en móvil y `space.6` desde 768; todas las tarjetas de una fila con el mismo alto; `$99,999.99` cabe completo a 360 px (P5). Columnas: 1 (360, 412), 2 (768), 3 (1024, 1440).

**Accesibilidad.** La tarjeta es un solo enlace con nombre completo («Gasto del mes, $8,240.10, 12 % más que septiembre»); la comparación usa palabras («más», «menos») además de flecha y color; la barra de avance es `role="meter"` con `aria-valuetext` («82 % del plan»). Pares: C-095 a C-105.

### 7.7 Gráficas

**Anatomía.** Título con el estilo de la etiqueta de tarjeta y periodo en `ui.text.secondary`; resumen de una frase debajo del título («El gasto subió 12 % respecto a septiembre»); área de trazado con retícula horizontal de 1 px sólida en `ui.viz.grid` (decorativa) y línea base en `ui.viz.baseline`; ejes en `font.size.caption`, `ui.text.secondary` y `tabular-nums` (dinero del eje compacto: `$0`, `$5 mil`, `$10 mil`); marcas según §8; leyenda (con 2 o más series) con marca de forma y nombre en `ui.text.default`; etiquetas directas en `ui.text.default`; descripción emergente con fondo `ui.surface.raised`, `radius.md` y `shadow.1` (fecha o categoría y una fila por serie: marca, nombre y valor exacto); botón de texto «Ver como tabla» (`table-2`) que alterna una tabla con los mismos datos.

**Variantes.** Barras horizontales, columnas, líneas, área apilada, dona y mini serie (reglas en §8); barra de avance (§7.6).

| Estado | Comportamiento |
|---|---|
| Reposo | dibujada completa, sin animación de entrada (§6.4) |
| Hover | puntero fino: guía vertical de 1 px en `ui.viz.baseline` y descripción emergente del punto o barra más cercano; las demás series no se atenúan |
| Foco | el contenedor recibe el foco; ← → recorren los puntos, ↑ ↓ cambian de serie; el punto enfocado lleva anillo `ui.border.focus` y muestra su descripción |
| Activo | tocar en móvil fija la descripción; tocar fuera la cierra |
| Deshabilitado | no aplica |
| Error | «No se pudo cargar la gráfica · Reintentar» en el alto reservado |
| Carga · vacío | esqueleto del alto final en `ui.skeleton.on-card`, sin ejes falsos; «Aún no hay datos para este periodo» con su acción |

**Tamaños.** Alto 200 px en < 768 y 280 px desde 768; marcas según §8.1; área de toque de cada punto o barra ≥ 44 px aunque la marca sea menor.

**Móvil y escritorio** (U1 §6.2). < 768: máximo 6 series, leyenda debajo, valor al tocar. ≥ 768: hasta 8 series, leyenda a la derecha, valor al pasar el puntero.

**Accesibilidad.** `<figure>` con `<figcaption>` (título y resumen); cada punto enfocable con nombre («Octubre, $8,240.10»); la tabla alterna siempre disponible; con `forced-colors: active` se conservan forma, trazo y patrón. Pares: C-106 a C-120.

### 7.14 Estados vacío, carga, error y estados globales (G-07)

| Estado | Anatomía y texto | Tokens |
|---|---|---|
| Carga | esqueleto con la forma final del contenido (U1 §2.4), estático (§6.2); `loader-circle` solo en acciones de más de 400 ms; a los 10 s: «Esto está tardando más de lo normal · Reintentar»; nunca cifras provisionales | `ui.skeleton.on-page`, `ui.skeleton.on-card` |
| Vacío | icono de la entidad en `size.icon.md`, una frase de qué falta y el botón que lo llena («Aún no hay gastos este mes · Registrar gasto»); nunca una pantalla en blanco | icono `ui.icon.default`, texto `ui.text.secondary`, botón §7.2 |
| Error | franja en la página con `circle-alert`, qué pasó y qué hacer, y «Reintentar» (`refresh-cw`); lo capturado no se pierde | fondo `color.role.error.surface-light` / `.surface-dark`, icono del rol, texto `ui.text.default` |
| Falta capturar (K18) | franja con `triangle-alert`: «Falta capturar · Capturar», abre P-26.05 en esa línea | superficie e icono de advertencia, texto `ui.text.default` |
| Sin conexión (G-07, K9) | franja fija bajo la barra superior: `wifi-off`, «Sin conexión con el servidor de casa» y «· n en cola» si hay capturas guardadas; lo cargado sigue visible con su hora de carga | superficie e icono de advertencia, texto `ui.text.default` |
| Sesión vencida | P-00 con aviso; al entrar vuelve a la página donde estaba | franja informativa |
| Registro inexistente o cambiado | «Ya no existe · Volver» (`file-x`); «Este registro cambió · Recargar» | franja informativa |

Franjas en página (avisos tipo *callout* de Notion): `radius.md`, relleno `space.3`, icono `size.icon.sm` del rol y texto en `ui.text.default` sobre la superficie del rol (pares C-145 a C-152: información, advertencia, error y éxito). Las franjas de error y sin conexión usan `role="alert"` la primera vez y luego quedan como texto; las demás, `role="status"`.

### 8.9 Barómetro (1.1, K54)

Para una sola razón contra dos umbrales; hoy, solo la capacidad de endeudamiento de Inicio (U3 B-1): pagos de crédito del mes ÷ ingreso del mes, con verde, amarillo y rojo (K54). Es una tarjeta de métrica (§7.6) con un arco en lugar de la barra de avance. **Numeración**: el texto de K54 la llama «§8.5», pero ese número ya es «Etiquetas de dinero siempre en texto» y lo citan muchos documentos; para no renumerar, el barómetro entra como §8.9 (PROVISIONAL, bloqueo `docs/estado/bloqueos/U3-tablero-barometro.md` punto 6).

| Parte | Especificación | Tokens |
|---|---|---|
| Arco de escala | semicírculo de 180°, de 0 % (izquierda) a la escala máxima que llega de C1 (derecha); tres tramos separados por 2 px del color de la superficie (§8.1): verde de 0 al primer umbral, amarillo hasta el segundo, rojo hasta el final | `ui.role.success`, `ui.role.warning`, `ui.role.error`; grosor `size.meter.height` |
| Pista del valor | arco concéntrico interior de 0 a la escala máxima | `ui.viz.track`, grosor `size.meter.height` |
| Valor | sobre la pista, de 0 al valor, en el color del tramo donde cae el valor | `ui.role.success` · `ui.role.warning` · `ui.role.error` |
| Aguja | del centro al arco de escala, extremo redondeado, con eje circular en el centro | `ui.text.default`, grosor `size.viz.line`; eje de `size.viz.marker` con anillo de 2 px de `ui.surface.card` |
| Marcas de umbral | marca radial en cada límite de tramo y su «[n] %» fuera del arco | marca de 1 px en `ui.viz.baseline`; texto `font.size.caption`, `ui.text.secondary`, `tabular-nums` |
| Dígito | debajo del centro, «[n] %»; en la línea siguiente, «de [tope] % de tope» | `font.size.title-1`, `font.weight.bold`, `tabular-nums`, `ui.text.default` (neutro: el estado lo dicen el arco y la insignia); segunda línea `font.size.small`, `ui.text.secondary` |
| Insignia de estado | icono + palabra del color + frase: «Verde · dentro del rango» (`circle-check`), «Amarillo · cerca del tope» (`triangle-alert`), «Rojo · arriba del tope» (`circle-alert`) | insignia de rol (§7.4): texto `ui.role.<rol>` sobre `ui.role.<rol>-surface` |

**Reglas.**
- El valor, los dos umbrales, la escala máxima y el estado llegan de C1 (U3 §2.6). El componente no fija 30 ni 40 ni calcula el estado: con otros umbrales en Ajustes (K44), los tramos y las marcas se mueven solos.
- Un valor mayor que la escala deja la aguja al final y el dígito dice la cifra exacta; el dígito nunca se recorta.
- Sin dato (por ejemplo, sin ingreso): arco con tramos y marcas, sin valor ni aguja; dígito «—» y el motivo con su acción (§7.14).
- Sin animación: la aguja se dibuja en su lugar (§6.4), también con movimiento reducido.
- El color nunca va solo (§2.10, K27): el estado lo dicen la insignia con icono y palabra y las marcas con su texto. Con `forced-colors: active`, el arco usa `CanvasText` y los tramos se reconocen por sus marcas y su texto.
- **Excepción de K54 a §2.4**: aquí el rol error marca el tramo rojo («arriba del tope»), aunque §2.4 lo reserva para validación y avisos, porque Axel pidió verde, amarillo y rojo. Los montos de la línea secundaria siguen §3.3: un negativo va en advertencia, nunca en error.

**Tamaños.** Alto del bloque `size.chart.height-compact` en < 768 y `size.chart.height` desde 768 (§7.7). El arco ocupa el ancho de la tarjeta menos su relleno y deja lugar al dígito y a la insignia; «100 %» cabe completo a 360 px (P5).

**Accesibilidad.** `role="meter"`, `aria-valuemin` = 0, `aria-valuemax` = escala máxima, `aria-valuenow` = valor y `aria-valuetext` con título, valor, tope y estado, por ejemplo «Capacidad de endeudamiento 27 % de 40 %: verde» (cifra sintética). El arco es decorativo (`aria-hidden`), porque el medidor ya da el valor; la insignia es texto.

**Pares de contraste** (numerados en §10.5 como C-159 a C-172; los calcula `docs/design/tools/contraste_componentes.py` desde `docs/contratos/C2-tokens.json`, claro y oscuro, truncados a 2 decimales como `contraste.csv` (§2.1), y los comprueba `make verify-wf1-tokens`. El primer cálculo de `u3-v2` redondeaba en lugar de truncar; por eso algunas cifras bajan 0.01. Ningún par queda bajo su mínimo):

| Id | Par | Mínimo | Claro | Oscuro |
|---|---|---|---|---|
| C-159 | `ui.role.success` sobre `ui.surface.card` (tramo verde) | interfaz 3 | 4.73 | 6.42 |
| C-160 | `ui.role.warning` sobre `ui.surface.card` (tramo amarillo) | interfaz 3 | 5.12 | 6.03 |
| C-161 | `ui.role.error` sobre `ui.surface.card` (tramo rojo) | interfaz 3 | 5.39 | 5.81 |
| C-162 | `ui.role.success` sobre `ui.viz.track` (valor sobre su pista) | interfaz 3 | 4.24 | 4.43 |
| C-163 | `ui.role.warning` sobre `ui.viz.track` | interfaz 3 | 4.59 | 4.16 |
| C-164 | `ui.role.error` sobre `ui.viz.track` | interfaz 3 | 4.83 | 4.01 |
| C-165 | `ui.role.success` sobre `ui.role.success-surface` (insignia) | texto 4.5 | 4.81 | 6.19 |
| C-166 | `ui.role.warning` sobre `ui.role.warning-surface` | texto 4.5 | 5.11 | 6.09 |
| C-167 | `ui.role.error` sobre `ui.role.error-surface` | texto 4.5 | 5.37 | 6.05 |
| C-168 | `ui.text.default` (dígito) sobre `ui.surface.card` | texto grande 3 | 10.75 | 10.75 |
| C-169 | `ui.text.secondary` (texto de las marcas y «de [tope] % de tope») sobre `ui.surface.card` | texto 4.5 | 5.03 | 7.97 |
| C-170 | `ui.viz.baseline` (marca de umbral) sobre `ui.surface.card` | interfaz 3 | 3.46 | 6.08 |
| C-171 | `ui.text.default` (aguja) sobre `ui.viz.track` | interfaz 3 | 9.63 | 7.42 |
| C-172 | `ui.viz.track` (pista del valor) sobre `ui.surface.card` | decorativo | 1.11 | 1.44 |

Las superficies de rol sobre la tarjeta dan de 1.00 a 1.04: por eso los tramos usan el color pleno de rol y no su superficie, que no se vería. **Sin tokens nuevos**: todo sale de C2 (`ui.role.*`, `ui.viz.*`, `ui.text.*`, `size.*`, `font.*`).

## 9. Accesibilidad y puntos de quiebre

### 9.1 Objetivo y cómo se cumple

**WCAG 2.2 nivel AA como mínimo, calculado** (decisiones §G). Idioma del documento `lang="es-MX"`.

| Criterio WCAG 2.2 | Cómo lo cumple este sistema |
|---|---|
| 1.3.1 Información y relaciones · 4.1.2 Nombre, función y valor | HTML semántico y patrones ARIA declarados en cada componente (§7): *grid*, *tabs*, *combobox*, *dialog*, *meter* |
| 1.3.4 Orientación | ninguna pantalla se bloquea en vertical u horizontal; en horizontal rige el ancho resultante (U1 §6.1 regla 7) |
| 1.4.1 Uso del color | signo, icono, forma, patrón y texto acompañan siempre al color (§2.10, §7.4, §8.6) |
| 1.4.3 Contraste (mínimo) · 1.4.11 Contraste no textual | texto ≥ 4.5:1; texto grande, interfaz, foco y gráficas ≥ 3:1; lista de pares de §10.5 calculada por U2d, 0 bajo umbral |
| 1.4.4 Cambio de tamaño · 1.4.10 Reajuste | a 200 % de texto y a 320 px de ancho CSS nada se pierde ni pide desplazamiento horizontal de la página; solo tabla, tablero y fila de pestañas se desplazan en su propio contenedor (U1 §6.1 regla 3) |
| 1.4.12 Espaciado del texto | ningún contenedor con alto fijo recorta texto si se amplía interlínea o espaciado; las cifras no se truncan (§3.3) |
| 1.4.13 Contenido al pasar el puntero o al enfocar | las descripciones emergentes se cierran con Esc, se pueden recorrer con el puntero y no desaparecen solas |
| 2.1.1 Teclado · 2.1.2 Sin trampas | todo se hace con teclado (U1 §5.3); solo diálogos y hojas atrapan el foco, y Esc los cierra |
| 2.1.4 Atajos de una tecla | `/`, `?` y las flechas solo actúan fuera de campos de texto (U1 §5.4); interruptor «Atajos de una tecla» en Apariencia para desactivarlos (§9.6; ratificado en la puerta 2, Q2-82 A, K51) |
| 2.3.1 Tres destellos | nada parpadea (§6.4) |
| 2.4.3 Orden del foco · 2.4.7 Foco visible · 2.4.11 Foco no oculto | el orden sigue la lectura; anillo de §9.2; `scroll-padding` igual a las barras fijas para que nada tape el elemento enfocado |
| 2.5.7 Movimientos de arrastre | arrastrar siempre tiene alternativa de un toque («Mover a…», §7.5.3) |
| 2.5.8 Tamaño del objetivo | todo control táctil ≥ 44 × 44 px (§9.3), por encima del piso de 24 px |
| 3.3.1 Identificación de errores · 3.3.3 Sugerencia | el mensaje dice qué pasó y qué hacer (§1.3, §7.0) |
| 3.3.7 Entrada redundante · 3.3.8 Autenticación accesible | los formularios proponen lo ya capturado (últimas categorías, tienda); el acceso admite gestores de contraseñas y pegar, sin pruebas cognitivas |

### 9.2 Foco visible

Anillo de 2 px (`border.width.focus`) en `ui.border.focus` (`color.brand.600` en claro, `color.brand.400` en oscuro) con separación de 2 px (`focus.offset`), por `:focus-visible`; en filas, celdas y elementos de lista va por dentro (`outline-offset: -2px`). Contraste medido contra página, tarjeta y fila seleccionada (C-006 a C-008). Nunca `outline: none` sin reemplazo; el anillo aparece sin transición. Con `forced-colors: active` el anillo usa el color del sistema `Highlight`.

### 9.3 Objetivo táctil ≥ 44 px

- Todo control táctil ocupa al menos `size.touch.min` (**44 × 44 px**) aunque su dibujo sea menor: casilla de 20 px, icono de 20 px, chip de 24 px, punto de una gráfica.
- Con puntero grueso (`any-pointer: coarse`) los controles miden 44 px de alto (`size.control.coarse`); el primario de hoja, 48 px; los destinos de la barra inferior, 56 px de alto (`size.bottombar.item.height`, igual a `layout.bottombar.height`) y al menos 64 px de ancho (`size.bottombar.item.min-width`) (§4.6).
- Entre objetivos contiguos hay `space.2` (8 px) como mínimo.
- Los 36 px de `size.control.fine` solo aplican con puntero fino y sin pantalla táctil.

### 9.4 Puntos de quiebre

Se diseña a 360 px y se amplía (móvil primero, U1 §6.1). Cinco anchos de referencia; las consultas son `min-width`.

| Token | Valor | Rango | Qué cambia (detalle en U1 §6.2 y §4.2) |
|---|---|---|---|
| `breakpoint.base` | 360 px | hasta 411 px | base sin consulta: 4 columnas, barra inferior, hojas a pantalla completa, tabla de 2 líneas, agenda |
| `breakpoint.phone` | 412 px | 412–767 px | un dato secundario más por tarjeta y fila; galería de 2 columnas |
| `breakpoint.tablet` | 768 px | 768–1023 px | 8 columnas; desaparece la barra inferior; tarjetas en 2 columnas; diálogos centrados; calendario mensual; barra lateral superpuesta |
| `breakpoint.desktop` | 1024 px | 1024–1439 px | 12 columnas; barra lateral fija y colapsable; tabla completa; vista previa lateral |
| `breakpoint.wide` | 1440 px | desde 1440 px | más aire; páginas de texto a 720 px; vista previa más ancha; propiedades en 2 columnas |

Además: a 320 px (200 % de zoom en 640 px o 400 % en 1280 px) todo se reajusta a una columna sin desplazamiento horizontal de la página (1.4.10). Referencias de prueba de U5: Android 412 × 915 y escritorio 1440 × 900 (U1 §6.1).

### 9.5 Pruebas que pide este sistema

U5 y WF-3 comprueban: navegación completa solo con teclado en 1440; TalkBack con Chrome en el Android de referencia; un lector de pantalla de escritorio en Linux; zoom de texto al 200 % y ancho de 320 px; vista en escala de grises (P4) y con las simulaciones de daltonismo de colorimetría §8.4; `forced-colors`; y la tabla de contraste de U2d con 0 pares bajo umbral.

### 9.6 Preferencias del sistema

| Preferencia | Comportamiento |
|---|---|
| `prefers-color-scheme` | tema por defecto según el sistema; el interruptor (Sistema · Claro · Oscuro, en Apariencia, en el pie de la barra lateral, `Ctrl+Shift+L` y /tema) fija `data-theme` en la raíz y se recuerda por dispositivo |
| `prefers-reduced-motion` | se anulan las transiciones no esenciales (§6.3) |
| `forced-colors` | colores del sistema; se conservan anillo de foco, bordes de controles, forma, trazo y patrón de las series |
| Atajos de una tecla | Para WCAG 2.1.4: interruptor local en Apariencia, encendido por omisión (ratificado por Axel en la puerta 2, Q2-82 A, K51); lo agregan U1 y U4 |

## FIN: docs/design/DESIGN.md

## ARCHIVO: docs/contratos/tokens.css


/* tokens.css
 * C2 v1.0 · Finanzas Hogar · GENERADO por Style Dictionary 5.6.1 desde docs/contratos/C2-tokens.json.
 * No se edita a mano: cambia C2-tokens.json y compila (docs/design/tools/tokens/README.md).
 */

/* Fuentes autoalojadas (Inter 4.1, SIL OFL 1.1: ../design/fuentes/inter/LICENSE.txt). Rutas relativas a este archivo; la PWA las copia a /fonts/ (WF-3). */
@font-face {
  font-family: "Inter";
  src: url("../design/fuentes/inter/web/InterVariable.woff2") format("woff2");
  font-weight: 100 900;
  font-style: normal;
  font-display: swap;
}
@font-face {
  font-family: "Inter";
  src: url("../design/fuentes/inter/web/InterVariable-Italic.woff2") format("woff2");
  font-weight: 100 900;
  font-style: italic;
  font-display: swap;
}

:root {
  color-scheme: light;
  font-optical-sizing: auto;

  /* Color · primitivos (iguales en los dos temas; C2-colores.json) */
  --color-base-white: #FFFFFF;
  --color-brand-50: #F9FAFF;
  --color-brand-100: #F0F1FF;
  --color-brand-200: #E1E4FF;
  --color-brand-300: #C9CFFF;
  --color-brand-400: #ADB3F4;
  --color-brand-500: #9095E0;
  --color-brand-600: #7579C2;
  --color-brand-700: #5D609F;
  --color-brand-800: #46487B;
  --color-brand-900: #303258;
  --color-brand-950: #1E203A;
  --color-neutral-50: #FAFAFB;
  --color-neutral-100: #F2F2F4;
  --color-neutral-200: #E5E6E9;
  --color-neutral-300: #D1D2D8;
  --color-neutral-400: #B7B8C0;
  --color-neutral-500: #9B9CA4;
  --color-neutral-600: #808189;
  --color-neutral-700: #66676E;
  --color-neutral-800: #4D4E53;
  --color-neutral-900: #36363A;
  --color-neutral-950: #222325;
  --color-role-income-50: #F6FCF3;
  --color-role-income-100: #E9F8E1;
  --color-role-income-200: #D7EFCA;
  --color-role-income-300: #B9E1A4;
  --color-role-income-400: #98CA7B;
  --color-role-income-500: #75B052;
  --color-role-income-600: #5B9335;
  --color-role-income-700: #457723;
  --color-role-income-800: #335A18;
  --color-role-income-900: #22400D;
  --color-role-income-950: #142907;
  --color-role-income-light: var(--color-role-income-700);
  --color-role-income-dark: var(--color-role-income-400);
  --color-role-income-surface-light: var(--color-role-income-100);
  --color-role-income-surface-dark: var(--color-role-income-900);
  --color-role-expense-50: #FFF9F6;
  --color-role-expense-100: #FFEFE8;
  --color-role-expense-200: #FFDED0;
  --color-role-expense-300: #FFC3A9;
  --color-role-expense-400: #F8A078;
  --color-role-expense-500: #E37E4E;
  --color-role-expense-600: #C46231;
  --color-role-expense-700: #A04B20;
  --color-role-expense-800: #7C3815;
  --color-role-expense-900: #59260B;
  --color-role-expense-950: #3B1705;
  --color-role-expense-light: var(--color-role-expense-700);
  --color-role-expense-dark: var(--color-role-expense-400);
  --color-role-expense-surface-light: var(--color-role-expense-100);
  --color-role-expense-surface-dark: var(--color-role-expense-900);
  --color-role-savings-50: #F0FDFE;
  --color-role-savings-100: #D8FAFC;
  --color-role-savings-200: #BBF2F5;
  --color-role-savings-300: #85E5EA;
  --color-role-savings-400: #3FD0D7;
  --color-role-savings-500: #00B2BA;
  --color-role-savings-600: #00949A;
  --color-role-savings-700: #00767B;
  --color-role-savings-800: #005A5E;
  --color-role-savings-900: #003F42;
  --color-role-savings-950: #00292B;
  --color-role-savings-light: var(--color-role-savings-700);
  --color-role-savings-dark: var(--color-role-savings-400);
  --color-role-savings-surface-light: var(--color-role-savings-100);
  --color-role-savings-surface-dark: var(--color-role-savings-900);
  --color-role-debt-50: #FCF8FF;
  --color-role-debt-100: #F7EEFF;
  --color-role-debt-200: #F0DEFF;
  --color-role-debt-300: #E5C3FF;
  --color-role-debt-400: #D0A3F2;
  --color-role-debt-500: #B782DD;
  --color-role-debt-600: #9B67C0;
  --color-role-debt-700: #7E509D;
  --color-role-debt-800: #603C79;
  --color-role-debt-900: #442857;
  --color-role-debt-950: #2C1939;
  --color-role-debt-light: var(--color-role-debt-700);
  --color-role-debt-dark: var(--color-role-debt-400);
  --color-role-debt-surface-light: var(--color-role-debt-100);
  --color-role-debt-surface-dark: var(--color-role-debt-900);
  --color-role-warning-50: #FFF9F1;
  --color-role-warning-100: #FFF0DC;
  --color-role-warning-200: #FDE2BD;
  --color-role-warning-300: #F6CB8F;
  --color-role-warning-400: #E5AE5B;
  --color-role-warning-500: #CF8E1B;
  --color-role-warning-600: #AD7500;
  --color-role-warning-700: #8B5D00;
  --color-role-warning-800: #6B4600;
  --color-role-warning-900: #4C3000;
  --color-role-warning-950: #321E00;
  --color-role-warning-light: var(--color-role-warning-700);
  --color-role-warning-dark: var(--color-role-warning-400);
  --color-role-warning-surface-light: var(--color-role-warning-100);
  --color-role-warning-surface-dark: var(--color-role-warning-900);
  --color-role-error-50: #FFF8F8;
  --color-role-error-100: #FFEEED;
  --color-role-error-200: #FFDDDC;
  --color-role-error-300: #FFC0BF;
  --color-role-error-400: #FC9999;
  --color-role-error-500: #E77678;
  --color-role-error-600: #C75A5E;
  --color-role-error-700: #A34449;
  --color-role-error-800: #7E3236;
  --color-role-error-900: #5A2124;
  --color-role-error-950: #3C1416;
  --color-role-error-light: var(--color-role-error-700);
  --color-role-error-dark: var(--color-role-error-400);
  --color-role-error-surface-light: var(--color-role-error-100);
  --color-role-error-surface-dark: var(--color-role-error-900);
  --color-role-success-50: #F4FDF5;
  --color-role-success-100: #E3FAE6;
  --color-role-success-200: #CDF1D3;
  --color-role-success-300: #A8E4B3;
  --color-role-success-400: #7ECF90;
  --color-role-success-500: #53B56D;
  --color-role-success-600: #339852;
  --color-role-success-700: #207B3E;
  --color-role-success-800: #155E2D;
  --color-role-success-900: #0A421E;
  --color-role-success-950: #052B11;
  --color-role-success-light: var(--color-role-success-700);
  --color-role-success-dark: var(--color-role-success-400);
  --color-role-success-surface-light: var(--color-role-success-100);
  --color-role-success-surface-dark: var(--color-role-success-900);
  --color-role-info-50: #F6FBFF;
  --color-role-info-100: #E7F5FF;
  --color-role-info-200: #CFEBFF;
  --color-role-info-300: #A6DAFF;
  --color-role-info-400: #6BC3FD;
  --color-role-info-500: #35A7EA;
  --color-role-info-600: #008ACB;
  --color-role-info-700: #006FA4;
  --color-role-info-800: #00547E;
  --color-role-info-900: #003B5A;
  --color-role-info-950: #00263C;
  --color-role-info-light: var(--color-role-info-700);
  --color-role-info-dark: var(--color-role-info-400);
  --color-role-info-surface-light: var(--color-role-info-100);
  --color-role-info-surface-dark: var(--color-role-info-900);
  --color-critical-50: #FCFBF0;
  --color-critical-100: #F8F4D9;
  --color-critical-200: #EEE8BD;
  --color-critical-300: #E0D58E;
  --color-critical-400: #CBBB5A;
  --color-critical-500: #B29E15;
  --color-critical-600: #948200;
  --color-critical-700: #766800;
  --color-critical-800: #5A4F00;
  --color-critical-900: #3F3700;
  --color-critical-950: #292300;
  --color-critical-light: var(--color-critical-700);
  --color-critical-dark: var(--color-critical-400);
  --color-critical-surface-light: var(--color-critical-100);
  --color-critical-surface-dark: var(--color-critical-900);
  --color-person-axel-light: #002E64;
  --color-person-axel-dark: #4F83C8;
  --color-person-ximena-light: #A14E78;
  --color-person-ximena-dark: #FFACD3;
  --color-person-karen-light: #B870AC;
  --color-person-karen-dark: #D389C6;
  --color-person-shared-light: #484B66;
  --color-person-shared-dark: #E7EAFF;
  --color-store-walmart-light: #4065AC;
  --color-store-walmart-dark: #7AA3EE;
  --color-store-bodega-aurrera-light: #762C4C;
  --color-store-bodega-aurrera-dark: #FFAAC9;
  --color-store-sams-club-light: #009D82;
  --color-store-sams-club-dark: #99FFE4;
  --color-store-market-light: #5C0204;
  --color-store-market-dark: #C2665B;
  --color-chart-1-light: #474698;
  --color-chart-1-dark: #797DCD;
  --color-chart-2-light: #009A87;
  --color-chart-2-dark: #61E4CD;
  --color-chart-3-light: #883500;
  --color-chart-3-dark: #C16D45;
  --color-chart-4-light: #BF6D2C;
  --color-chart-4-dark: #FFB989;
  --color-chart-5-light: #812E63;
  --color-chart-5-dark: #B86797;
  --color-chart-6-light: #0092B9;
  --color-chart-6-dark: #76D9FF;
  --color-chart-7-light: #683987;
  --color-chart-7-dark: #9C70BC;
  --color-chart-8-light: #579749;
  --color-chart-8-dark: #A0DE93;

  /* Tipografía */
  --font-sans: "Inter", ui-sans-serif, system-ui, "Noto Sans", Roboto, sans-serif, "Noto Color Emoji";
  --font-numeric-tabular: tabular-nums;
  --font-size-caption: 0.75rem;
  --font-size-small: 0.875rem;
  --font-size-body: 1rem;
  --font-size-title-3: 1.25rem;
  --font-size-title-2: 1.5rem;
  --font-size-title-1: 2rem;
  --line-height-caption: 1rem;
  --line-height-small: 1.25rem;
  --line-height-body: 1.5rem;
  --line-height-title-3: 1.75rem;
  --line-height-title-2: 2rem;
  --line-height-title-1: 2.5rem;
  --font-weight-regular: 400;
  --font-weight-medium: 500;
  --font-weight-semibold: 600;
  --font-weight-bold: 700;
  --font-letter-spacing: 0;

  /* Espaciado (retícula de 4 pt) */
  --space-0: 0;
  --space-1: 0.25rem;
  --space-2: 0.5rem;
  --space-3: 0.75rem;
  --space-4: 1rem;
  --space-6: 1.5rem;
  --space-8: 2rem;
  --space-12: 3rem;

  /* Radios */
  --radius-sm: 6px;
  --radius-md: 10px;
  --radius-lg: 14px;
  --radius-full: 9999px;

  /* Bordes */
  --border-width-default: 1px;
  --border-width-focus: 2px;

  /* Foco */
  --focus-offset: 2px;

  /* Sombras (valor del tema claro; en oscuro, --ui-shadow-* vale none y se usa borde) */
  --shadow-0: none;
  --shadow-1: 0 1px 2px 0 color-mix(in srgb, var(--color-neutral-950) 8%, transparent),
    0 1px 3px 0 color-mix(in srgb, var(--color-neutral-950) 6%, transparent);
  --shadow-2: 0 4px 12px 0 color-mix(in srgb, var(--color-neutral-950) 10%, transparent),
    0 12px 32px 0 color-mix(in srgb, var(--color-neutral-950) 8%, transparent);

  /* Medidas de control, icono y componente */
  --size-touch-min: 44px;
  --size-control-coarse: var(--size-touch-min);
  --size-control-primary-sheet: 48px;
  --size-control-fine: 36px;
  --size-icon-sm: 20px;
  --size-icon-md: 24px;
  --size-bottombar-item-height: var(--layout-bottombar-height);
  --size-bottombar-item-min-width: 64px;
  --size-bottombar-add: 40px;
  --size-bottombar-indicator-width: 56px;
  --size-bottombar-indicator-height: 32px;
  --size-sidebar-item: 32px;
  --size-chip-height: 24px;
  --size-badge-height: 20px;
  --size-avatar: 24px;
  --size-mark-category: 10px;
  --size-checkbox: 20px;
  --size-switch-width: 36px;
  --size-switch-height: 20px;
  --size-meter-height: 8px;
  --size-sheet-handle-width: 32px;
  --size-sheet-handle-height: 4px;
  --size-chart-height-compact: 200px;
  --size-chart-height: 280px;
  --size-calendar-cell: 96px;
  --size-calendar-cell-wide: 120px;
  --size-agenda-row: 56px;
  --size-viz-bar-max: 24px;
  --size-viz-bar-radius: 4px;
  --size-viz-line: 2px;
  --size-viz-marker: 10px;
  --size-viz-gap: 2px;

  /* Medidas de estructura */
  --layout-sidebar-width: 240px;
  --layout-sidebar-overlay-width: 280px;
  --layout-sidebar-mobile-width: 85%;
  --layout-topbar-height-mobile: 56px;
  --layout-topbar-height-desktop: 48px;
  --layout-bottombar-height: 56px;
  --layout-reading-max-width: 720px;
  --layout-preview-width-desktop: 400px;
  --layout-preview-width-wide: 480px;
  --layout-dialog-width-short: 480px;
  --layout-dialog-width-long: 640px;
  --layout-command-width: 640px;
  --layout-board-column-width: 272px;
  --layout-toast-max-width: var(--layout-dialog-width-short);

  /* Iconos */
  --icon-stroke: 1.5px;
  --icon-stroke-user-units-md: 1.5;
  --icon-stroke-user-units-sm: 1.8;

  /* Movimiento */
  --duration-fast: 150ms;
  --duration-base: 200ms;
  --duration-slow: 250ms;
  --ease-standard: cubic-bezier(0.2, 0, 0, 1);
  --ease-emphasized: cubic-bezier(0.05, 0.7, 0.1, 1);

  /* Puntos de quiebre (referencia; no sirven dentro de @media) */
  --breakpoint-base: 360px;
  --breakpoint-phone: 412px;
  --breakpoint-tablet: 768px;
  --breakpoint-desktop: 1024px;
  --breakpoint-wide: 1440px;

  /* Series de gráficas: trazo (stroke-dasharray; none = continuo) */
  --viz-series-1-dash: none;
  --viz-series-2-dash: 8 4;
  --viz-series-3-dash: 2 4;
  --viz-series-4-dash: 8 4 2 4;
  --viz-series-5-dash: 4 4;
  --viz-series-6-dash: 12 4;
  --viz-series-7-dash: 8 3 2 3 2 3;
  --viz-series-8-dash: 2 2;

  /* Variables temáticas --ui-* · tema claro (predeterminado) */
  --ui-surface-page: var(--color-neutral-50);
  --ui-surface-card: var(--color-neutral-100);
  --ui-surface-raised: var(--color-neutral-50);
  --ui-surface-field: var(--color-neutral-50);
  --ui-surface-hover: var(--color-neutral-200);
  --ui-surface-hover-subtle: var(--ui-surface-card);
  --ui-surface-pressed: var(--color-neutral-300);
  --ui-surface-selected: var(--color-brand-100);
  --ui-surface-disabled: var(--color-neutral-100);
  --ui-skeleton-on-page: var(--color-neutral-100);
  --ui-skeleton-on-card: var(--color-neutral-200);
  --ui-text-default: var(--color-neutral-900);
  --ui-text-secondary: var(--color-neutral-700);
  --ui-text-disabled: var(--color-neutral-700);
  --ui-text-link: var(--color-brand-700);
  --ui-text-link-hover: var(--color-brand-800);
  --ui-text-on-primary: var(--color-base-white);
  --ui-text-on-destructive: var(--color-base-white);
  --ui-text-on-critical: var(--color-base-white);
  --ui-icon-default: var(--ui-text-secondary);
  --ui-icon-selected: var(--ui-action-primary-bg);
  --ui-border-divider: var(--color-neutral-200);
  --ui-border-control: var(--color-neutral-600);
  --ui-border-control-hover: var(--ui-text-secondary);
  --ui-border-focus: var(--color-brand-600);
  --ui-border-selected: var(--ui-action-primary-bg);
  --ui-border-error: var(--color-role-error-light);
  --ui-action-primary-bg: var(--color-brand-700);
  --ui-action-primary-bg-hover: var(--color-brand-800);
  --ui-action-primary-bg-pressed: var(--color-brand-900);
  --ui-action-destructive-bg: var(--color-role-error-700);
  --ui-action-destructive-bg-hover: var(--color-role-error-800);
  --ui-action-destructive-bg-pressed: var(--color-role-error-900);
  --ui-action-critical-bg: var(--color-critical-700);
  --ui-action-critical-bg-hover: var(--color-critical-800);
  --ui-action-critical-bg-pressed: var(--color-critical-900);
  --ui-viz-grid: var(--ui-border-divider);
  --ui-viz-baseline: var(--ui-border-control);
  --ui-viz-other: var(--ui-border-control);
  --ui-viz-track: var(--ui-border-divider);
  --ui-viz-meter: var(--ui-text-secondary);
  --ui-role-income: var(--color-role-income-light);
  --ui-role-income-surface: var(--color-role-income-surface-light);
  --ui-role-expense: var(--color-role-expense-light);
  --ui-role-expense-surface: var(--color-role-expense-surface-light);
  --ui-role-savings: var(--color-role-savings-light);
  --ui-role-savings-surface: var(--color-role-savings-surface-light);
  --ui-role-debt: var(--color-role-debt-light);
  --ui-role-debt-surface: var(--color-role-debt-surface-light);
  --ui-role-warning: var(--color-role-warning-light);
  --ui-role-warning-surface: var(--color-role-warning-surface-light);
  --ui-role-error: var(--color-role-error-light);
  --ui-role-error-surface: var(--color-role-error-surface-light);
  --ui-role-success: var(--color-role-success-light);
  --ui-role-success-surface: var(--color-role-success-surface-light);
  --ui-role-info: var(--color-role-info-light);
  --ui-role-info-surface: var(--color-role-info-surface-light);
  --ui-critical: var(--color-critical-light);
  --ui-critical-surface: var(--color-critical-surface-light);
  --ui-person-axel: var(--color-person-axel-light);
  --ui-person-ximena: var(--color-person-ximena-light);
  --ui-person-karen: var(--color-person-karen-light);
  --ui-person-shared: var(--color-person-shared-light);
  --ui-store-walmart: var(--color-store-walmart-light);
  --ui-store-bodega-aurrera: var(--color-store-bodega-aurrera-light);
  --ui-store-sams-club: var(--color-store-sams-club-light);
  --ui-store-market: var(--color-store-market-light);
  --ui-chart-1: var(--color-chart-1-light);
  --ui-chart-2: var(--color-chart-2-light);
  --ui-chart-3: var(--color-chart-3-light);
  --ui-chart-4: var(--color-chart-4-light);
  --ui-chart-5: var(--color-chart-5-light);
  --ui-chart-6: var(--color-chart-6-light);
  --ui-chart-7: var(--color-chart-7-light);
  --ui-chart-8: var(--color-chart-8-light);
  --ui-overlay-scrim: color-mix(in srgb, var(--color-neutral-950) 40%, transparent);
  --ui-shadow-1: 0 1px 2px 0 color-mix(in srgb, var(--color-neutral-950) 8%, transparent),
    0 1px 3px 0 color-mix(in srgb, var(--color-neutral-950) 6%, transparent);
  --ui-shadow-2: 0 4px 12px 0 color-mix(in srgb, var(--color-neutral-950) 10%, transparent),
    0 12px 32px 0 color-mix(in srgb, var(--color-neutral-950) 8%, transparent);
}

/* Tema oscuro según el sistema, salvo que el interruptor fije «claro» */
@media (prefers-color-scheme: dark) {
  :root:not([data-theme="light"]) {
    color-scheme: dark;
    --ui-surface-page: var(--color-neutral-950);
    --ui-surface-card: var(--color-neutral-900);
    --ui-surface-raised: var(--color-neutral-900);
    --ui-surface-field: var(--color-neutral-950);
    --ui-surface-hover: var(--color-neutral-800);
    --ui-surface-hover-subtle: var(--ui-surface-card);
    --ui-surface-pressed: var(--color-neutral-700);
    --ui-surface-selected: var(--color-brand-900);
    --ui-surface-disabled: var(--color-neutral-900);
    --ui-skeleton-on-page: var(--color-neutral-900);
    --ui-skeleton-on-card: var(--color-neutral-800);
    --ui-text-default: var(--color-neutral-100);
    --ui-text-secondary: var(--color-neutral-300);
    --ui-text-disabled: var(--color-neutral-400);
    --ui-text-link: var(--color-brand-300);
    --ui-text-link-hover: var(--color-brand-200);
    --ui-text-on-primary: var(--color-brand-950);
    --ui-text-on-destructive: var(--color-role-error-950);
    --ui-text-on-critical: var(--color-critical-950);
    --ui-icon-default: var(--ui-text-secondary);
    --ui-icon-selected: var(--ui-action-primary-bg);
    --ui-border-divider: var(--color-neutral-800);
    --ui-border-control: var(--color-neutral-400);
    --ui-border-control-hover: var(--ui-text-secondary);
    --ui-border-focus: var(--color-brand-400);
    --ui-border-selected: var(--ui-action-primary-bg);
    --ui-border-error: var(--color-role-error-dark);
    --ui-action-primary-bg: var(--color-brand-300);
    --ui-action-primary-bg-hover: var(--color-brand-200);
    --ui-action-primary-bg-pressed: var(--color-brand-100);
    --ui-action-destructive-bg: var(--color-role-error-300);
    --ui-action-destructive-bg-hover: var(--color-role-error-200);
    --ui-action-destructive-bg-pressed: var(--color-role-error-100);
    --ui-action-critical-bg: var(--color-critical-300);
    --ui-action-critical-bg-hover: var(--color-critical-200);
    --ui-action-critical-bg-pressed: var(--color-critical-100);
    --ui-viz-grid: var(--ui-border-divider);
    --ui-viz-baseline: var(--ui-border-control);
    --ui-viz-other: var(--ui-border-control);
    --ui-viz-track: var(--ui-border-divider);
    --ui-viz-meter: var(--ui-text-secondary);
    --ui-role-income: var(--color-role-income-dark);
    --ui-role-income-surface: var(--color-role-income-surface-dark);
    --ui-role-expense: var(--color-role-expense-dark);
    --ui-role-expense-surface: var(--color-role-expense-surface-dark);
    --ui-role-savings: var(--color-role-savings-dark);
    --ui-role-savings-surface: var(--color-role-savings-surface-dark);
    --ui-role-debt: var(--color-role-debt-dark);
    --ui-role-debt-surface: var(--color-role-debt-surface-dark);
    --ui-role-warning: var(--color-role-warning-dark);
    --ui-role-warning-surface: var(--color-role-warning-surface-dark);
    --ui-role-error: var(--color-role-error-dark);
    --ui-role-error-surface: var(--color-role-error-surface-dark);
    --ui-role-success: var(--color-role-success-dark);
    --ui-role-success-surface: var(--color-role-success-surface-dark);
    --ui-role-info: var(--color-role-info-dark);
    --ui-role-info-surface: var(--color-role-info-surface-dark);
    --ui-critical: var(--color-critical-dark);
    --ui-critical-surface: var(--color-critical-surface-dark);
    --ui-person-axel: var(--color-person-axel-dark);
    --ui-person-ximena: var(--color-person-ximena-dark);
    --ui-person-karen: var(--color-person-karen-dark);
    --ui-person-shared: var(--color-person-shared-dark);
    --ui-store-walmart: var(--color-store-walmart-dark);
    --ui-store-bodega-aurrera: var(--color-store-bodega-aurrera-dark);
    --ui-store-sams-club: var(--color-store-sams-club-dark);
    --ui-store-market: var(--color-store-market-dark);
    --ui-chart-1: var(--color-chart-1-dark);
    --ui-chart-2: var(--color-chart-2-dark);
    --ui-chart-3: var(--color-chart-3-dark);
    --ui-chart-4: var(--color-chart-4-dark);
    --ui-chart-5: var(--color-chart-5-dark);
    --ui-chart-6: var(--color-chart-6-dark);
    --ui-chart-7: var(--color-chart-7-dark);
    --ui-chart-8: var(--color-chart-8-dark);
    --ui-overlay-scrim: color-mix(in srgb, var(--color-neutral-950) 60%, transparent);
    --ui-shadow-1: none;
    --ui-shadow-2: none;
  }
}

/* Tema oscuro fijado por el interruptor (Sistema · Claro · Oscuro) */
:root[data-theme="dark"] {
  color-scheme: dark;
  --ui-surface-page: var(--color-neutral-950);
  --ui-surface-card: var(--color-neutral-900);
  --ui-surface-raised: var(--color-neutral-900);
  --ui-surface-field: var(--color-neutral-950);
  --ui-surface-hover: var(--color-neutral-800);
  --ui-surface-hover-subtle: var(--ui-surface-card);
  --ui-surface-pressed: var(--color-neutral-700);
  --ui-surface-selected: var(--color-brand-900);
  --ui-surface-disabled: var(--color-neutral-900);
  --ui-skeleton-on-page: var(--color-neutral-900);
  --ui-skeleton-on-card: var(--color-neutral-800);
  --ui-text-default: var(--color-neutral-100);
  --ui-text-secondary: var(--color-neutral-300);
  --ui-text-disabled: var(--color-neutral-400);
  --ui-text-link: var(--color-brand-300);
  --ui-text-link-hover: var(--color-brand-200);
  --ui-text-on-primary: var(--color-brand-950);
  --ui-text-on-destructive: var(--color-role-error-950);
  --ui-text-on-critical: var(--color-critical-950);
  --ui-icon-default: var(--ui-text-secondary);
  --ui-icon-selected: var(--ui-action-primary-bg);
  --ui-border-divider: var(--color-neutral-800);
  --ui-border-control: var(--color-neutral-400);
  --ui-border-control-hover: var(--ui-text-secondary);
  --ui-border-focus: var(--color-brand-400);
  --ui-border-selected: var(--ui-action-primary-bg);
  --ui-border-error: var(--color-role-error-dark);
  --ui-action-primary-bg: var(--color-brand-300);
  --ui-action-primary-bg-hover: var(--color-brand-200);
  --ui-action-primary-bg-pressed: var(--color-brand-100);
  --ui-action-destructive-bg: var(--color-role-error-300);
  --ui-action-destructive-bg-hover: var(--color-role-error-200);
  --ui-action-destructive-bg-pressed: var(--color-role-error-100);
  --ui-action-critical-bg: var(--color-critical-300);
  --ui-action-critical-bg-hover: var(--color-critical-200);
  --ui-action-critical-bg-pressed: var(--color-critical-100);
  --ui-viz-grid: var(--ui-border-divider);
  --ui-viz-baseline: var(--ui-border-control);
  --ui-viz-other: var(--ui-border-control);
  --ui-viz-track: var(--ui-border-divider);
  --ui-viz-meter: var(--ui-text-secondary);
  --ui-role-income: var(--color-role-income-dark);
  --ui-role-income-surface: var(--color-role-income-surface-dark);
  --ui-role-expense: var(--color-role-expense-dark);
  --ui-role-expense-surface: var(--color-role-expense-surface-dark);
  --ui-role-savings: var(--color-role-savings-dark);
  --ui-role-savings-surface: var(--color-role-savings-surface-dark);
  --ui-role-debt: var(--color-role-debt-dark);
  --ui-role-debt-surface: var(--color-role-debt-surface-dark);
  --ui-role-warning: var(--color-role-warning-dark);
  --ui-role-warning-surface: var(--color-role-warning-surface-dark);
  --ui-role-error: var(--color-role-error-dark);
  --ui-role-error-surface: var(--color-role-error-surface-dark);
  --ui-role-success: var(--color-role-success-dark);
  --ui-role-success-surface: var(--color-role-success-surface-dark);
  --ui-role-info: var(--color-role-info-dark);
  --ui-role-info-surface: var(--color-role-info-surface-dark);
  --ui-critical: var(--color-critical-dark);
  --ui-critical-surface: var(--color-critical-surface-dark);
  --ui-person-axel: var(--color-person-axel-dark);
  --ui-person-ximena: var(--color-person-ximena-dark);
  --ui-person-karen: var(--color-person-karen-dark);
  --ui-person-shared: var(--color-person-shared-dark);
  --ui-store-walmart: var(--color-store-walmart-dark);
  --ui-store-bodega-aurrera: var(--color-store-bodega-aurrera-dark);
  --ui-store-sams-club: var(--color-store-sams-club-dark);
  --ui-store-market: var(--color-store-market-dark);
  --ui-chart-1: var(--color-chart-1-dark);
  --ui-chart-2: var(--color-chart-2-dark);
  --ui-chart-3: var(--color-chart-3-dark);
  --ui-chart-4: var(--color-chart-4-dark);
  --ui-chart-5: var(--color-chart-5-dark);
  --ui-chart-6: var(--color-chart-6-dark);
  --ui-chart-7: var(--color-chart-7-dark);
  --ui-chart-8: var(--color-chart-8-dark);
  --ui-overlay-scrim: color-mix(in srgb, var(--color-neutral-950) 60%, transparent);
  --ui-shadow-1: none;
  --ui-shadow-2: none;
}

/* Dinero, porcentajes y cantidades: cifras tabulares (DESIGN §3.3) */
.money {
  font-variant-numeric: var(--font-numeric-tabular);
}

/* Movimiento reducido: se anulan las transiciones no esenciales (DESIGN §6.3) */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    scroll-behavior: auto !important;
  }
  .icon-spinner { animation: spin 2s cubic-bezier(0, 0, 1, 1) infinite !important; }
}

## FIN: docs/contratos/tokens.css

## ARCHIVO: docs/contratos/extractos/inicio.yaml


# Finanzas Hogar · Contrato C1 v1.5.1 · extracto «inicio». Archivo generado por backend/internal/api/c1gen/build.py: no se edita a mano.
# Fuente única: docs/contratos/openapi.yaml. Índice y regla de uso: docs/contratos/extractos/INDICE.md.
openapi: "3.1.1"
info:
  title: "Finanzas Hogar · API (C1) · extracto inicio"
  version: "1.5.1"
  description: "Extracto de la etiqueta «Inicio» del contrato C1 (docs/contratos/openapi.yaml): sus rutas y los esquemas que usan. Parámetros, encabezados, respuestas y esquemas base viven en componentes-comunes.yaml."
security:
  - cookieAuth: []
tags:
  - name: "Inicio"
    description: "Inicio: tarjetas del resumen, calendario del hogar, búsqueda y vistas de base de datos."
paths:
  /health:
    get:
      operationId: "getHealth"
      tags:
        - "Inicio"
      summary: "Estado del servidor"
      description: "Público. Responde si el servidor y la base están en marcha (franja «Sin conexión», G-07). No expone datos."
      x-necesidades:
        - "N-004"
      x-tablas: []
      security: []
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Health"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/auth/me:
    get:
      operationId: "getCurrentSession"
      tags:
        - "Inicio"
      summary: "Sesión actual"
      description: "Sesión actual."
      x-necesidades:
        - "N-001"
        - "N-052"
      x-tablas:
        - "device_session"
        - "person"
        - "user_account"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/CurrentSession"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/search:
    get:
      operationId: "search"
      tags:
        - "Inicio"
      summary: "Búsqueda global"
      description: "Sin distinguir mayúsculas ni acentos. Máximo 5 por grupo; «Ver todos» usa la lista de cada recurso con `q`. Texto vacío: grupos vacíos (los recientes viven en el dispositivo, N-003)."
      x-necesidades:
        - "N-002"
        - "N-072"
      x-tablas:
        - "category"
        - "credit"
        - "credit_purchase"
        - "expense"
        - "income"
        - "presentation"
        - "product"
        - "recipe"
        - "recipe_ingredient"
        - "savings_goal"
        - "store"
      parameters:
        - $ref: "componentes-comunes.yaml#/components/parameters/Q"
        - name: "type"
          in: "query"
          required: false
          description: "Restringe a uno o varios tipos (product para N-072)."
          style: "form"
          explode: false
          schema:
            type: "array"
            items:
              type: "string"
              enum:
                - "expense"
                - "income"
                - "recipe"
                - "product"
                - "store"
                - "credit"
                - "credit_purchase"
                - "savings_goal"
                - "settings_section"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/SearchResult"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        "422":
          $ref: "componentes-comunes.yaml#/components/responses/UnprocessableEntity"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/calendar/events:
    get:
      operationId: "listCalendarEvents"
      tags:
        - "Inicio"
      summary: "Calendario del hogar"
      description: "Eventos del rango `period` (por defecto el mes en curso; agenda móvil: rango `hoy..hoy+30` o `hoy..hoy+60`). Orden fijo: fecha y, dentro del día, el orden de tipos de N-029. N-086 = tipos de dinero; N-104 = income; N-115 = savings_goal y savings_movement; N-192 y N-208 = credit_payment, promotion y card_statement (N-208 con credit_id)."
      x-necesidades:
        - "N-029"
        - "N-086"
        - "N-104"
        - "N-115"
        - "N-192"
        - "N-208"
      x-tablas:
        - "credit"
        - "credit_payment"
        - "credit_purchase"
        - "credit_statement"
        - "expense"
        - "income"
        - "job"
        - "planned_item"
        - "promotion"
        - "savings_goal"
        - "savings_movement"
        - "scheduled_payment"
        - "scheduled_payment_share"
        - "store"
      parameters:
        - $ref: "componentes-comunes.yaml#/components/parameters/Period"
        - $ref: "componentes-comunes.yaml#/components/parameters/CreditId"
        - $ref: "componentes-comunes.yaml#/components/parameters/Family"
        - name: "event_type"
          in: "query"
          required: false
          description: "Tipos de evento a incluir (todos por defecto)."
          style: "form"
          explode: false
          schema:
            type: "array"
            items:
              type: "string"
              enum:
                - "income"
                - "credit_payment"
                - "fixed_expense"
                - "upcoming_expense"
                - "card_statement"
                - "planned_grocery"
                - "market_visit"
                - "promotion"
                - "savings_goal"
                - "savings_movement"
                - "meal_plan"
                - "monthly_report"
        - $ref: "componentes-comunes.yaml#/components/parameters/Cursor"
        - $ref: "componentes-comunes.yaml#/components/parameters/Limit"
        - $ref: "componentes-comunes.yaml#/components/parameters/Sort"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/CalendarEventPage"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        "422":
          $ref: "componentes-comunes.yaml#/components/responses/UnprocessableEntity"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/views/{entity}:
    get:
      operationId: "queryView"
      tags:
        - "Inicio"
      summary: "Vista de base de datos"
      description: "Endpoint genérico de vistas (decisión G): tabla, tablero, calendario o galería por entidad, con filtros, orden, agrupación y paginación únicos. La respuesta añade `mode` y `groups` (columnas del tablero o días del calendario, con conteo y total). Entidades, agrupaciones y fechas permitidas: C1 §8."
      x-necesidades:
        - "N-043"
        - "N-045"
        - "N-046"
        - "N-111"
        - "N-128"
        - "N-135"
        - "N-170"
        - "N-171"
        - "N-180"
        - "N-202"
        - "N-247"
        - "N-281"
        - "N-287"
        - "N-290"
      x-tablas:
        - "attachment"
        - "backup"
        - "category"
        - "credit"
        - "credit_purchase"
        - "credit_term_rate"
        - "device_session"
        - "expense"
        - "expense_share"
        - "income"
        - "inventory_item"
        - "job"
        - "meal_plan_entry"
        - "notification_delivery"
        - "person"
        - "price"
        - "product"
        - "promotion"
        - "recipe"
        - "recipe_meal_slot"
        - "savings_goal"
        - "savings_movement"
        - "scheduled_payment"
        - "store"
        - "user_account"
        - "weight_entry"
      parameters:
        - name: "entity"
          in: "path"
          required: true
          description: "Entidad (nombre_codigo de la tabla) de la vista."
          schema:
            type: "string"
            enum:
              - "person"
              - "user_account"
              - "device_session"
              - "income"
              - "category"
              - "expense"
              - "savings_goal"
              - "savings_movement"
              - "credit"
              - "credit_term_rate"
              - "credit_purchase"
              - "scheduled_payment"
              - "promotion"
              - "store"
              - "product"
              - "price"
              - "recipe"
              - "meal_plan_entry"
              - "inventory_item"
              - "weight_entry"
              - "notification_delivery"
              - "job"
              - "backup"
        - name: "mode"
          in: "query"
          required: true
          description: "Modo de la vista."
          schema:
            type: "string"
            enum:
              - "table"
              - "board"
              - "calendar"
              - "gallery"
        - name: "group_by"
          in: "query"
          required: false
          description: "Propiedad de agrupación del tablero (nombre_codigo o calculada); valores permitidos por entidad en C1 §8."
          schema:
            type: "string"
            pattern: "^[a-z][a-z0-9_]*$"
        - name: "date_field"
          in: "query"
          required: false
          description: "Propiedad de fecha del calendario (C1 §8)."
          schema:
            type: "string"
            pattern: "^[a-z][a-z0-9_]*$"
        - $ref: "componentes-comunes.yaml#/components/parameters/Q"
        - name: "status"
          in: "query"
          required: false
          description: "Estado propio de la entidad (C1 §8)."
          schema:
            type: "string"
            pattern: "^[a-z][a-z_]*$"
        - $ref: "componentes-comunes.yaml#/components/parameters/PersonId"
        - $ref: "componentes-comunes.yaml#/components/parameters/Family"
        - $ref: "componentes-comunes.yaml#/components/parameters/CategoryId"
        - $ref: "componentes-comunes.yaml#/components/parameters/StoreId"
        - $ref: "componentes-comunes.yaml#/components/parameters/Period"
        - $ref: "componentes-comunes.yaml#/components/parameters/CreditId"
        - $ref: "componentes-comunes.yaml#/components/parameters/Cursor"
        - $ref: "componentes-comunes.yaml#/components/parameters/Limit"
        - $ref: "componentes-comunes.yaml#/components/parameters/Sort"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/ViewPage"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        "404":
          $ref: "componentes-comunes.yaml#/components/responses/NotFound"
        "422":
          $ref: "componentes-comunes.yaml#/components/responses/UnprocessableEntity"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/dashboard/month-vs-plan:
    get:
      operationId: "getDashboardMonthVsPlan"
      tags:
        - "Inicio"
      summary: "Tarjeta «Mes vs plan»"
      description: "Tarjeta «Mes vs plan»."
      x-necesidades:
        - "N-020"
      x-tablas:
        - "budget"
        - "category"
        - "credit_purchase"
        - "expense"
        - "expense_share"
        - "period"
      parameters:
        - $ref: "componentes-comunes.yaml#/components/parameters/Period"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DashboardMonthVsPlan"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        "422":
          $ref: "componentes-comunes.yaml#/components/responses/UnprocessableEntity"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/dashboard/fortnight-flow:
    get:
      operationId: "getDashboardFortnightFlow"
      tags:
        - "Inicio"
      summary: "Tarjeta «Flujo quincenal»"
      description: "Tarjeta «Flujo quincenal»."
      x-necesidades:
        - "N-021"
      x-tablas:
        - "credit"
        - "credit_payment"
        - "credit_purchase"
        - "credit_purchase_share"
        - "expense"
        - "expense_share"
        - "income"
        - "person"
        - "planned_item"
        - "scheduled_payment"
        - "scheduled_payment_share"
      parameters:
        - $ref: "componentes-comunes.yaml#/components/parameters/Period"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DashboardFortnightFlow"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        "422":
          $ref: "componentes-comunes.yaml#/components/responses/UnprocessableEntity"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/dashboard/next-payment:
    get:
      operationId: "getDashboardNextPayment"
      tags:
        - "Inicio"
      summary: "Tarjeta «Próximo pago y disponible»"
      description: "Tarjeta «Próximo pago y disponible»."
      x-necesidades:
        - "N-022"
      x-tablas:
        - "credit"
        - "credit_payment"
        - "credit_purchase"
        - "credit_purchase_share"
        - "person"
        - "scheduled_payment"
        - "scheduled_payment_share"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DashboardNextPayment"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/dashboard/week-meals:
    get:
      operationId: "getDashboardWeekMeals"
      tags:
        - "Inicio"
      summary: "Tarjeta «Comidas de la semana»"
      description: "Tarjeta «Comidas de la semana»."
      x-necesidades:
        - "N-023"
      x-tablas:
        - "health_profile"
        - "meal_plan_entry"
        - "meal_slot"
        - "nutrition_fact"
        - "person"
        - "presentation"
        - "price"
        - "recipe"
        - "recipe_ingredient"
      parameters:
        - $ref: "componentes-comunes.yaml#/components/parameters/Period"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DashboardWeekMeals"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        "422":
          $ref: "componentes-comunes.yaml#/components/responses/UnprocessableEntity"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/dashboard/best-basket:
    get:
      operationId: "getDashboardBestBasket"
      tags:
        - "Inicio"
      summary: "Tarjeta «Mejor canasta»"
      description: "Mejor canasta de la lista de la semana `week`, la que muestra la pantalla (R3-01)."
      x-necesidades:
        - "N-024"
      x-tablas:
        - "app_setting"
        - "expense"
        - "expense_item"
        - "inventory_item"
        - "meal_plan_entry"
        - "presentation"
        - "price"
        - "product"
        - "recipe_ingredient"
        - "shopping_list"
        - "shopping_list_item"
        - "store"
      parameters:
        - $ref: "componentes-comunes.yaml#/components/parameters/ListWeek"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DashboardBestBasket"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        "422":
          $ref: "componentes-comunes.yaml#/components/responses/UnprocessableEntity"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/dashboard/debt-capacity:
    get:
      operationId: "getDashboardDebtCapacity"
      tags:
        - "Inicio"
      summary: "Barómetro de capacidad de endeudamiento"
      description: "Tablero de control (K54, U3 2.0 B-1, RN-30): capacidad de endeudamiento del mes con los umbrales vigentes. period: solo mes calendario AAAA-MM; por omisión el mes en curso; una quincena o una semana responde 422 validation_failed (PROVISIONAL, bloqueo C1-tablero, punto 5). Necesidad PROVISIONAL (bloqueo C1-tablero, punto 3): U1 no tiene aún su necesidad; cuelga de N-022."
      x-provisional: true
      x-necesidades:
        - "N-022"
      x-tablas:
        - "app_setting"
        - "credit"
        - "credit_purchase"
        - "income"
        - "period"
        - "planned_item"
        - "scheduled_payment"
        - "scheduled_payment_share"
      parameters:
        - $ref: "componentes-comunes.yaml#/components/parameters/Period"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DashboardDebtCapacity"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        "422":
          $ref: "componentes-comunes.yaml#/components/responses/UnprocessableEntity"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/dashboard/monthly-overview:
    get:
      operationId: "getDashboardMonthlyOverview"
      tags:
        - "Inicio"
      summary: "Resumen mensual del tablero"
      description: "Tablero de control (K54, U3 2.0 G-1 y G-6): una sola petición para las dos gráficas. `from` y `to` obligatorios (AAAA-MM), `from` ≤ `to` y como máximo 24 meses; si no, 422 validation_failed. Necesidad PROVISIONAL (bloqueo C1-tablero, punto 3): U1 no tiene aún su necesidad; cuelga de N-028."
      x-provisional: true
      x-necesidades:
        - "N-028"
      x-tablas:
        - "app_setting"
        - "budget"
        - "category"
        - "credit"
        - "credit_purchase"
        - "expense"
        - "expense_share"
        - "income"
        - "period"
        - "planned_item"
        - "scheduled_payment"
        - "scheduled_payment_share"
      parameters:
        - $ref: "componentes-comunes.yaml#/components/parameters/OverviewFrom"
        - $ref: "componentes-comunes.yaml#/components/parameters/OverviewTo"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DashboardMonthlyOverview"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        "422":
          $ref: "componentes-comunes.yaml#/components/responses/UnprocessableEntity"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/budget/signals:
    get:
      operationId: "getBudgetSignals"
      tags:
        - "Inicio"
      summary: "Sugerencia subir o recortar presupuesto"
      description: "Tablero de control (K54, U3 2.0 G-3, RN-31): sugerencia por categoría. period: solo mes calendario AAAA-MM; por omisión el mes en curso (PROVISIONAL, bloqueo C1-tablero, punto 5). Necesidad PROVISIONAL (bloqueo C1-tablero, punto 3): U1 no tiene aún su necesidad; cuelga de N-027."
      x-provisional: true
      x-necesidades:
        - "N-027"
      x-tablas:
        - "app_setting"
        - "budget"
        - "category"
        - "credit_purchase"
        - "expense"
        - "expense_share"
        - "period"
      parameters:
        - $ref: "componentes-comunes.yaml#/components/parameters/Period"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/BudgetSignals"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        "422":
          $ref: "componentes-comunes.yaml#/components/responses/UnprocessableEntity"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/dashboard/savings:
    get:
      operationId: "getDashboardSavings"
      tags:
        - "Inicio"
      summary: "Tarjeta «Ahorro»"
      description: "Tarjeta «Ahorro»."
      x-necesidades:
        - "N-025"
      x-tablas:
        - "savings_goal"
        - "savings_movement"
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/DashboardSavings"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
  /api/v1/persons/{person_id}/participation:
    get:
      operationId: "getPersonParticipation"
      tags:
        - "Inicio"
      summary: "Participación de Karen"
      description: "Solo personas externas (participation = external); con otra persona responde 404."
      x-necesidades:
        - "N-026"
      x-tablas:
        - "credit"
        - "credit_payment"
        - "credit_purchase"
        - "credit_purchase_share"
        - "expense"
        - "expense_share"
        - "person"
        - "scheduled_payment"
        - "scheduled_payment_share"
      parameters:
        - name: "person_id"
          in: "path"
          required: true
          description: "Id del registro (person_id)."
          schema:
            type: "integer"
            minimum: 1
      responses:
        "200":
          description: "Correcto."
          headers:
            X-Request-Id:
              $ref: "componentes-comunes.yaml#/components/headers/RequestId"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/PersonParticipation"
        "400":
          $ref: "componentes-comunes.yaml#/components/responses/BadRequest"
        "401":
          $ref: "componentes-comunes.yaml#/components/responses/Unauthorized"
        "404":
          $ref: "componentes-comunes.yaml#/components/responses/NotFound"
        default:
          $ref: "componentes-comunes.yaml#/components/responses/Error"
components:
  securitySchemes:
    cookieAuth:
      type: "apiKey"
      in: "cookie"
      name: "fh_session"
      description: "Sesión larga por dispositivo. Cookie HttpOnly, Secure y SameSite=Strict; el token en claro solo viaja en la cookie (la base guarda su SHA-256 en device_session.token_hash). Duración 90 días renovable con el uso (K11, app_setting session.max_age_days)."
    serviceToken:
      type: "http"
      scheme: "bearer"
      description: "Credencial de servicio del agente externo (K53): encabezado Authorization: Bearer con un token propio, distinto de la cookie de sesión. Solo lo aceptan las operaciones con x-agente-externo (/api/v1/external-agent/*), que no aceptan la cookie; ninguna otra operación lo acepta. El token en claro nunca se guarda ni se registra. PROVISIONAL (bloqueo S5-agente-externo, pregunta e): el esquema 1.4 no tiene dónde guardar la credencial; dónde vive y a nombre de quién quedan audit_log e idempotency_key se decide en la puerta 4."
      x-provisional: true
  schemas:
    BasisPoints:
      type: "integer"
      minimum: 0
      maximum: 10000
      description: "Porcentaje en puntos base: 10000 = 100 %, 5000 = 50 %."
      examples:
        - 5000
    Health:
      type: "object"
      description: "Estado mínimo del servidor (público, sin datos)."
      properties:
        status:
          type: "string"
          enum:
            - "ok"
          description: "Siempre ok si el servidor responde y la base abre."
        time:
          $ref: "componentes-comunes.yaml#/components/schemas/DateTime"
          description: "Hora del servidor."
      required:
        - "status"
        - "time"
      x-calculado: true
    CurrentSession:
      type: "object"
      description: "Sesión de quien hace la petición (N-001). Axel y Ximena tienen los mismos permisos."
      properties:
        user_account_id:
          type: "integer"
          description: "Cuenta (user_account.id)."
        username:
          type: "string"
          description: "Usuario (user_account.username)."
        person_id:
          type: "integer"
          description: "Persona dueña de la cuenta (user_account.person_id)."
        display_name:
          type: "string"
          description: "Nombre visible (person.display_name)."
        participation:
          type: "string"
          enum:
            - "holder"
            - "co_holder"
          description: "Papel (person.participation): titular o co-titular, mismos permisos."
        device_session_id:
          type: "integer"
          description: "Sesión de este dispositivo (device_session.id)."
        expires_at:
          $ref: "componentes-comunes.yaml#/components/schemas/DateTime"
          description: "Vencimiento de la sesión (device_session.expires_at)."
        last_seen_at:
          $ref: "componentes-comunes.yaml#/components/schemas/DateTime"
          description: "Último uso (device_session.last_seen_at)."
      required:
        - "user_account_id"
        - "username"
        - "person_id"
        - "display_name"
        - "participation"
        - "device_session_id"
        - "expires_at"
        - "last_seen_at"
      x-calculado: true
    SavingsGoal:
      type: "object"
      description: "Meta de ahorro (tabla `savings_goal`). Fondo de ahorro o meta con monto y fecha, con su aportación mensual planificada. El avance se calcula con sus movimientos. Campos ausentes = sin valor (NULL)."
      x-tabla: "savings_goal"
      required:
        - "id"
        - "name"
        - "kind"
        - "start_date"
        - "is_active"
        - "created_at"
        - "updated_at"
      properties:
        id:
          type: "integer"
          description: "id"
        name:
          type: "string"
          description: "Nombre (Fondo de emergencia)"
        kind:
          type: "string"
          enum:
            - "fund"
            - "goal"
          description: "fund = fondo de ahorro, goal = meta con monto y fecha"
        target_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Monto objetivo en centavos"
        target_date:
          $ref: "componentes-comunes.yaml#/components/schemas/Date"
          description: "Fecha objetivo"
        planned_monthly_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Aportación mensual planificada en centavos"
        start_date:
          $ref: "componentes-comunes.yaml#/components/schemas/Date"
          description: "Inicio del fondo o meta"
        is_active:
          type: "boolean"
          description: "is_active"
        notes:
          type: "string"
          description: "notes"
        created_at:
          $ref: "componentes-comunes.yaml#/components/schemas/DateTime"
          description: "created_at"
        updated_at:
          $ref: "componentes-comunes.yaml#/components/schemas/DateTime"
          description: "updated_at"
        saved_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Ahorrado = Σ aportaciones − Σ retiros."
          x-calculado: true
        progress_tenths_pct:
          type: "integer"
          description: "saved ÷ target × 1000, entero ROUND_HALF_UP: décimas de punto porcentual (500 = 50.0 %; RN-18). Lo muestra la pantalla de Metas con un decimal; la tarjeta de Inicio usa progress_pct, entero (Q2-17 C, motor S2). Sin tope: una meta cumplida puede pasar de 1000. Ausente sin target_cents, como progress_pct."
          x-calculado: true
        progress_pct:
          type: "integer"
          description: "saved ÷ target × 100, entero ROUND_HALF_UP (solo metas)."
          x-calculado: true
        required_monthly_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Aportación mensual necesaria para llegar a target_date (motor S2)."
          x-calculado: true
        status:
          type: "string"
          enum:
            - "in_progress"
            - "achieved"
            - "overdue"
          description: "Estado de la meta (Q2-18 A; motor S2 1.3 §3.13): achieved si target_cents no es null y saved_cents ≥ target_cents; si no, overdue si target_date no es null y es anterior a hoy; si no, in_progress (también un fondo sin monto o sin fecha)."
          x-calculado: true
    DashboardMonthVsPlan:
      type: "object"
      description: "Tarjeta «Mes vs plan» (N-020, métrica M-01). Fórmulas: real = Σ parte del hogar devengada en el mes (motor S2); plan = Σ monto del renglón de budget vigente en el mes (RN-22) de las categorías activas; % usado = real ÷ plan × 100 entero ROUND_HALF_UP; restante = plan − real; estado normal < 80, warning ≥ 80, exceeded ≥ 100; días que faltan contando hoy. Sin plan: used_pct y status ausentes («—»)."
      properties:
        month:
          $ref: "componentes-comunes.yaml#/components/schemas/Month"
          description: "Mes en curso."
        actual_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Real."
        planned_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Plan."
        used_pct:
          type: "integer"
          description: "% usado."
        remaining_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Restante."
        status:
          type: "string"
          enum:
            - "normal"
            - "warning"
            - "exceeded"
          description: "normal si % < 80; warning si ≥ 80; exceeded si ≥ 100 (decisión F)."
        days_left:
          type: "integer"
          description: "Días que faltan contando hoy."
        cash_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Efectivo del hogar del mes (cash_household_cents de M-11). Nuevo en C1 v1.5 (K54)."
        previous_month:
          $ref: "componentes-comunes.yaml#/components/schemas/Month"
          description: "Mes anterior al medido. Nuevo en C1 v1.5."
        previous_actual_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Real del mes anterior. Ausente sin datos de ese mes."
        previous_cash_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Efectivo del hogar del mes anterior. Ausente sin datos de ese mes."
        actual_difference_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "actual_cents − previous_actual_cents, con signo. Ausente si falta el mes anterior."
        cash_difference_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "cash_cents − previous_cash_cents, con signo. Ausente si falta el mes anterior."
        categories_warning_count:
          type: "integer"
          description: "Categorías activas con estado warning en el mes (umbrales de getAlertSettings)."
        categories_exceeded_count:
          type: "integer"
          description: "Categorías activas con estado exceeded en el mes."
      required:
        - "month"
        - "actual_cents"
        - "days_left"
      x-calculado: true
    DashboardFortnightFlow:
      type: "object"
      description: "Tarjeta «Flujo quincenal» (N-021, métrica M-02), efectivo y parte del hogar (motor S2). Fórmulas: ingresos = quincena capturada (o esperada, con income_is_expected) + extras con fecha en la quincena; egresos hechos = salidas de efectivo hasta hoy (contado + partes del hogar pagadas); comprometido = partes del hogar de pagos pendientes de la quincena + fijos y venideros sin pagar; saldo al cierre = ingresos − egresos hechos − comprometido. Quincena del 1 al 15 o del 16 al último día (K8)."
      properties:
        start_date:
          $ref: "componentes-comunes.yaml#/components/schemas/Date"
          description: "Inicio de la quincena."
        end_date:
          $ref: "componentes-comunes.yaml#/components/schemas/Date"
          description: "Fin."
        income_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Ingresos."
        income_is_expected:
          type: "boolean"
          description: "true si la quincena aún no se captura."
        spent_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Egresos hechos."
        committed_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Comprometido."
        closing_balance_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Saldo al cierre."
      required:
        - "start_date"
        - "end_date"
        - "income_cents"
        - "income_is_expected"
        - "spent_cents"
        - "committed_cents"
        - "closing_balance_cents"
      x-calculado: true
    DashboardNextPayment:
      type: "object"
      description: "Tarjeta «Próximo pago y disponible» (N-022, métricas M-03 y M-04). Próximo pago = entre todas las líneas, el pago con alguna parte pendiente y fecha más cercana ≥ hoy; si hay pendientes con fecha < hoy, primero el más antiguo. Disponible: lectura activa A o B."
      properties:
        next_payment:
          type: "object"
          description: "Próximo pago (N-022 a)."
          properties:
            due_date:
              $ref: "componentes-comunes.yaml#/components/schemas/Date"
              description: "Fecha."
            credit_id:
              type: "integer"
              description: "Línea."
            amount_cents:
              $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
              description: "Monto total (si caen varios el mismo día, la suma)."
            household_share_cents:
              $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
              description: "Parte del hogar."
            external_share_cents:
              $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
              description: "Parte de Karen."
            days_left:
              type: "integer"
              description: "Días que faltan (negativo si está vencido)."
            payment_count:
              type: "integer"
              description: "Pagos que caen ese día."
            is_overdue:
              type: "boolean"
              description: "true si es un pago pendiente con fecha anterior a hoy («sin marcar como pagado»)."
          required:
            - "due_date"
            - "credit_id"
            - "amount_cents"
            - "household_share_cents"
            - "external_share_cents"
            - "days_left"
            - "payment_count"
            - "is_overdue"
          x-calculado: true
        availability:
          type: "array"
          items:
            type: "object"
            description: "Disponible de una línea revolvente con interruptor A/B (N-022 b, motor S3)."
            properties:
              credit_id:
                type: "integer"
                description: "Línea."
              availability_reading:
                type: "string"
                enum:
                  - "A"
                  - "B"
                description: "Lectura activa (credit.availability_reading; B por defecto)."
              available_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Disponible de hoy con la lectura activa."
              availability_confirmed:
                type: "boolean"
                description: "false = mostrar «Reposición de capital pendiente de confirmar con Kueski»."
            required:
              - "credit_id"
              - "availability_reading"
              - "available_cents"
              - "availability_confirmed"
            x-calculado: true
          description: "Una fila por línea revolvente activa."
      required:
        - "availability"
      x-calculado: true
    DashboardWeekMeals:
      type: "object"
      description: "Tarjeta «Comidas de la semana» (N-023, métrica M-05, motor S4). Costo = Σ costo de casillas con receta («fuera» no suma); kcal promedio por día de cada persona = Σ kcal asignadas ÷ días con al menos un tiempo asignado; casillas sin armar = pares día × tiempo obligatorio (desayuno, comida, cena) en que alguna persona no tiene receta ni «fuera» (máx. 21; K43). Semana de lunes a domingo (K8)."
      properties:
        week_start:
          $ref: "componentes-comunes.yaml#/components/schemas/Date"
          description: "Lunes de la semana."
        week_end:
          $ref: "componentes-comunes.yaml#/components/schemas/Date"
          description: "Domingo."
        cost_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Costo de la semana."
        persons:
          type: "array"
          items:
            type: "object"
            description: "Persona."
            properties:
              person_id:
                type: "integer"
                description: "Axel o Ximena."
              avg_kcal_per_day:
                type: "integer"
                description: "kcal promedio por día."
              days_with_plan:
                type: "integer"
                description: "Días con al menos un tiempo asignado (n de 7)."
              daily_kcal_target:
                type: "integer"
                description: "Meta diaria (health_profile.daily_kcal_target)."
            required:
              - "person_id"
              - "days_with_plan"
            x-calculado: true
          description: "Axel y Ximena."
        unplanned_slots:
          type: "integer"
          description: "Casillas sin armar (0 a 21)."
      required:
        - "week_start"
        - "week_end"
        - "cost_cents"
        - "persons"
        - "unplanned_slots"
      x-calculado: true
    DashboardBestBasket:
      type: "object"
      description: "Tarjeta «Mejor canasta» (N-024, métrica M-06, motor S5): lista de compras vigente (o la última compra con renglones si está vacía) comparada por tienda; Sam's Club como si tuvieran membresía."
      properties:
        basket:
          type: "string"
          enum:
            - "shopping_list"
            - "expense"
          description: "Canasta usada."
        expense_id:
          type: "integer"
          description: "Compra usada si la lista está vacía."
        store_id:
          type: "integer"
          description: "Tienda con menor total."
        total_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Total en esa tienda."
        difference_vs_second_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Diferencia contra la segunda tienda."
        second_store_id:
          type: "integer"
          description: "Segunda tienda: contra la que se mide difference_vs_second_cents (Q2-86 A). La tarjeta de Inicio la nombra («$138.50 menos que …»). Ausente si no hay segunda tienda."
        covered_count:
          type: "integer"
          description: "Artículos con precio en esa tienda."
        item_count:
          type: "integer"
          description: "Artículos de la canasta."
        stale_price_count:
          type: "integer"
          description: "Precios viejos usados."
        assumes_membership:
          type: "boolean"
          description: "true si la ganadora exige membresía y se comparó como si la tuvieran."
      required:
        - "basket"
        - "covered_count"
        - "item_count"
        - "stale_price_count"
        - "assumes_membership"
      x-calculado: true
    DashboardSavings:
      type: "object"
      description: "Tarjeta «Ahorro» (N-025, métrica M-07). Ahorrado = Σ aportaciones − Σ retiros; aportado del mes contra planificado del mes; meta más próxima por fecha con % = ahorrado ÷ monto meta × 100 entero."
      properties:
        saved_total_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Ahorrado."
        contributed_month_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Aportado este mes."
        planned_month_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Planificado del mes."
        next_goal:
          $ref: "#/components/schemas/SavingsGoal"
          description: "Meta activa con la fecha meta más próxima."
      required:
        - "saved_total_cents"
        - "contributed_month_cents"
        - "planned_month_cents"
      x-calculado: true
    PersonParticipation:
      type: "object"
      description: "Participación de una persona externa (Karen) (N-026): solo despensa común y su línea revolvente. Nota fija: «Karen paga directo a Kueski su parte; su parte es informativa y repone la línea cuando paga»."
      properties:
        person_id:
          type: "integer"
          description: "Persona externa."
        month:
          $ref: "componentes-comunes.yaml#/components/schemas/Month"
          description: "Mes en curso."
        accrued_groceries_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Su parte devengada de la despensa en el mes (motor S2)."
        purchases:
          type: "array"
          items:
            type: "object"
            description: "Compra particular viva."
            properties:
              credit_purchase_id:
                type: "integer"
                description: "Compra."
              description:
                type: "string"
                description: "Comercio o descripción."
              purchase_date:
                $ref: "componentes-comunes.yaml#/components/schemas/Date"
                description: "Fecha."
              total_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Total."
              payments_made:
                type: "integer"
                description: "Pagos hechos."
              payments_total:
                type: "integer"
                description: "Pagos de la serie."
              balance_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Saldo."
            required:
              - "credit_purchase_id"
              - "purchase_date"
              - "total_cents"
              - "payments_made"
              - "payments_total"
              - "balance_cents"
            x-calculado: true
          description: "Sus compras a crédito vivas."
        payment_shares:
          type: "array"
          items:
            type: "object"
            description: "Su parte de un pago."
            properties:
              scheduled_payment_id:
                type: "integer"
                description: "Pago."
              due_date:
                $ref: "componentes-comunes.yaml#/components/schemas/Date"
                description: "Fecha."
              amount_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Monto de su parte."
              status:
                type: "string"
                enum:
                  - "paid"
                  - "pending"
                description: "Pagada o pendiente."
            required:
              - "scheduled_payment_id"
              - "due_date"
              - "amount_cents"
              - "status"
            x-calculado: true
          description: "Con fecha en el mes."
        live_debt_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Σ de sus partes pendientes."
      required:
        - "person_id"
        - "month"
        - "accrued_groceries_cents"
        - "purchases"
        - "payment_shares"
        - "live_debt_cents"
      x-calculado: true
    CalendarEvent:
      type: "object"
      description: "Evento del calendario del hogar (N-029 y calendarios filtrados)."
      properties:
        event_type:
          type: "string"
          enum:
            - "income"
            - "credit_payment"
            - "fixed_expense"
            - "upcoming_expense"
            - "card_statement"
            - "planned_grocery"
            - "market_visit"
            - "promotion"
            - "savings_goal"
            - "savings_movement"
            - "meal_plan"
            - "monthly_report"
          description: "Tipo de evento."
        date:
          $ref: "componentes-comunes.yaml#/components/schemas/Date"
          description: "Fecha (inicio en promociones)."
        end_date:
          $ref: "componentes-comunes.yaml#/components/schemas/Date"
          description: "Fin (solo promociones)."
        title:
          type: "string"
          description: "Texto del evento en español."
        amount_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Monto total."
        household_share_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Parte del hogar."
        external_share_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Parte de Karen."
        status:
          type: "string"
          enum:
            - "pending"
            - "paid"
            - "overdue"
            - "expected"
            - "captured"
            - "complete"
            - "incomplete"
            - "empty"
          description: "Estado cuando aplica."
        credit_id:
          type: "integer"
          description: "Línea relacionada."
        ref_table:
          type: "string"
          description: "Tabla de origen (nombre_codigo)."
        ref_id:
          type: "integer"
          description: "Id del registro de origen."
      required:
        - "event_type"
        - "date"
        - "title"
        - "ref_table"
        - "ref_id"
      x-calculado: true
    ViewGroup:
      type: "object"
      description: "Grupo de una vista (columna de tablero o día de calendario)."
      properties:
        key:
          type: "string"
          description: "Valor de agrupación (id, enumerado o fecha ISO)."
        label:
          type: "string"
          description: "Etiqueta en español."
        count:
          type: "integer"
          description: "Registros del grupo con los filtros."
        total_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Σ del importe principal de la entidad (C1 §8), si tiene."
        share_bp:
          $ref: "#/components/schemas/BasisPoints"
          description: "total del grupo ÷ total de la vista."
      required:
        - "key"
        - "label"
        - "count"
      x-calculado: true
    ViewItem:
      type: "object"
      description: "Registro proyectado en una vista."
      properties:
        id:
          type: "integer"
          description: "Id del registro."
        title:
          type: "string"
          description: "Título (nombre o descripción)."
        group_key:
          type: "string"
          description: "Grupo al que pertenece (tablero, calendario)."
        date:
          $ref: "componentes-comunes.yaml#/components/schemas/Date"
          description: "Fecha del campo de calendario."
        amount_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Importe principal."
        image_attachment_id:
          type: "integer"
          description: "Adjunto con foto (galería de tickets)."
        record:
          type: "object"
          additionalProperties: true
          description: "Registro completo con las propiedades del esquema de su entidad (components.schemas.<Entidad>)."
      required:
        - "id"
        - "title"
        - "record"
      x-calculado: true
    SearchResult:
      type: "object"
      description: "Búsqueda global (N-002), agrupada por tipo; máximo 5 por grupo."
      properties:
        groups:
          type: "array"
          items:
            type: "object"
            description: "Grupo."
            properties:
              type:
                type: "string"
                enum:
                  - "expense"
                  - "income"
                  - "recipe"
                  - "product"
                  - "store"
                  - "credit"
                  - "credit_purchase"
                  - "savings_goal"
                  - "settings_section"
                description: "Tipo."
              total:
                type: "integer"
                description: "Coincidencias totales."
              items:
                type: "array"
                items:
                  type: "object"
                  description: "Coincidencia."
                  properties:
                    id:
                      type: "string"
                      description: "Id del registro o clave de la sección."
                    title:
                      type: "string"
                      description: "Texto principal."
                    subtitle:
                      type: "string"
                      description: "Texto secundario."
                    date:
                      $ref: "componentes-comunes.yaml#/components/schemas/Date"
                      description: "Fecha."
                  required:
                    - "id"
                    - "title"
                  x-calculado: true
                description: "Máximo 5."
            required:
              - "type"
              - "total"
              - "items"
            x-calculado: true
          description: "Grupos con coincidencias."
      required:
        - "groups"
      x-calculado: true
    DashboardDebtCapacity:
      type: "object"
      description: "Barómetro de capacidad de endeudamiento del tablero (U3 2.0 B-1, métrica M-20, RN-30; K54 decisiones 1 y 2). Capacidad usada = pagos del mes ÷ ingreso mensual neto × 100, entero ROUND_HALF_UP; estado con ese entero: normal (verde) si ratio_pct ≤ warning_pct, warning (amarillo) si ratio_pct ≤ limit_pct, exceeded (rojo) si ratio_pct > limit_pct. Ingreso de S2 y pagos de S3; la interfaz no calcula nada. PROVISIONAL (bloqueos U3-tablero-barometro y C1-tablero): ingreso sin extras, pagos con la parte del hogar, mes calendario y escala."
      properties:
        month:
          $ref: "componentes-comunes.yaml#/components/schemas/Month"
          description: "Mes medido (period; por omisión el mes en curso)."
        income_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          x-provisional: true
          description: "Ingreso mensual neto: quincenas de Axel del mes, capturadas más las esperadas sin capturar, sin extras (PROVISIONAL, bloqueo U3-tablero-barometro). Ausente sin quincena capturada ni esperada (sin datos no es cero)."
        income_is_expected:
          type: "boolean"
          description: "true si alguna quincena de income_cents es esperada (sin capturar)."
        extras_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Extras del mes. Informativo: no entra al cociente."
        payments_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          x-provisional: true
          description: "Parte del hogar de los pagos programados con fecha en el mes, de todas las líneas, pagados o no; la misma cifra que CreditTotals.month_household_cents del mes en curso (PROVISIONAL, bloqueo U3-tablero-barometro)."
        payments_total_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Los mismos pagos con la parte de Karen (informativo)."
        payment_count:
          type: "integer"
          description: "Pagos programados con fecha en el mes."
        ratio_pct:
          type: "integer"
          description: "payments_cents ÷ income_cents × 100, entero ROUND_HALF_UP (RN-01). Ausente sin income_cents o con income_cents = 0."
        ratio_total_pct:
          type: "integer"
          description: "payments_total_cents ÷ income_cents × 100, entero ROUND_HALF_UP (informativo). Ausente igual que ratio_pct."
        warning_pct:
          type: "integer"
          description: "Frontera entre verde y amarillo, en % (30 al sembrar, K54 decisión 2). Valor vigente de app_setting debt.capacity_warning_pct (K44; se cambia en updateDashboardSettings)."
          x-app-setting: "debt.capacity_warning_pct"
        limit_pct:
          type: "integer"
          description: "Tope, en %: arriba es rojo (40 al sembrar, K54 decisión 2). Valor vigente de app_setting debt.capacity_limit_pct (K44; se cambia en updateDashboardSettings)."
          x-app-setting: "debt.capacity_limit_pct"
        scale_max_pct:
          type: "integer"
          x-provisional: true
          description: "Fin de la escala del dibujo, en %: el múltiplo de 10 igual o mayor que 1.5 × limit_pct (60 con 40). Solo dibujo (PROVISIONAL, bloqueo U3-tablero-barometro)."
        status:
          type: "string"
          enum:
            - "normal"
            - "warning"
            - "exceeded"
          description: "Estado del barómetro con ratio_pct (fronteras inclusivas: 30 es normal y 40 es warning con los valores iniciales). Ausente sin ratio_pct."
        headroom_cents:
          $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
          description: "Margen antes del rojo, con signo: limit_pct × income_cents ÷ 100 (ROUND_HALF_UP) − payments_cents; negativo = arriba del tope. Ausente sin income_cents."
        has_credits:
          type: "boolean"
          x-provisional: true
          description: "true si hay al menos una línea de crédito activa o algún pago de una línea archivada con fecha en el mes (PROVISIONAL, bloqueo S6-tablero-lineas-archivadas); con false, payments_cents = 0 y ratio_pct = 0 (cifra real, no falta de dato)."
      required:
        - "month"
        - "income_is_expected"
        - "extras_cents"
        - "payments_cents"
        - "payments_total_cents"
        - "payment_count"
        - "warning_pct"
        - "limit_pct"
        - "scale_max_pct"
        - "has_credits"
      x-calculado: true
      x-provisional: true
    DashboardMonthlyOverview:
      type: "object"
      description: "Resumen mensual del tablero (U3 2.0 G-1 y G-6, métrica M-21; K54): ingreso, gasto devengado y efectivo, plan, pagos de crédito y capacidad de endeudamiento por mes, de `from` a `to` (máximo 24 meses). Motores S2 y S3. PROVISIONAL (bloqueos U3-tablero-barometro y C1-tablero): la capacidad sigue RN-30."
      properties:
        warning_pct:
          type: "integer"
          description: "Frontera entre verde y amarillo, en %. Valor vigente de app_setting debt.capacity_warning_pct (K44; se cambia en updateDashboardSettings)."
          x-app-setting: "debt.capacity_warning_pct"
        limit_pct:
          type: "integer"
          description: "Tope, en %. Valor vigente de app_setting debt.capacity_limit_pct (K44; se cambia en updateDashboardSettings)."
          x-app-setting: "debt.capacity_limit_pct"
        scale_max_pct:
          type: "integer"
          x-provisional: true
          description: "Fin de la escala del eje Y de G-6, en %: la misma de getDashboardDebtCapacity (M-20), el múltiplo de 10 igual o mayor que 1.5 × limit_pct (60 con 40); G-6 dibuja de 0 al mayor entre este valor y el mayor debt_ratio_pct a la vista. Solo dibujo. Nuevo en C1 v1.5.1 (hallazgo M-R3-1; PROVISIONAL, bloqueo U3-tablero-barometro)."
        months:
          type: "array"
          items:
            type: "object"
            description: "Mes del resumen."
            properties:
              month:
                $ref: "componentes-comunes.yaml#/components/schemas/Month"
                description: "Mes."
              is_projection:
                type: "boolean"
                description: "true si el mes es posterior al mes en curso (sin devengado ni efectivo)."
              income_total_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Ingreso del mes = income_fortnights_cents + income_extras_cents. Ausente sin ingreso capturado ni esperado."
              income_fortnights_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Quincenas del mes, capturadas más esperadas sin capturar. Ausente igual."
              income_extras_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Extras del mes. Ausente igual."
              income_is_expected:
                type: "boolean"
                description: "true si alguna quincena del mes es esperada (sin capturar)."
              accrued_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Devengado del hogar del mes (real de M-01, motor S2). Ausente en proyección y en meses sin datos."
              cash_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Efectivo del hogar del mes (cash_household_cents de M-11, motor S2). Ausente igual."
              planned_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Plan del mes (Σ renglón de budget vigente de las categorías activas, RN-22). Ausente sin plan."
              credit_payments_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Parte del hogar de los pagos de crédito con fecha en el mes (la de RN-30)."
              credit_payments_total_cents:
                $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                description: "Los mismos pagos con la parte de Karen."
              debt_ratio_pct:
                type: "integer"
                description: "Capacidad usada del mes (RN-30): credit_payments_cents ÷ income_fortnights_cents × 100, entero ROUND_HALF_UP. Ausente sin ingreso."
              debt_status:
                type: "string"
                enum:
                  - "normal"
                  - "warning"
                  - "exceeded"
                description: "Estado de debt_ratio_pct con warning_pct y limit_pct. Ausente sin debt_ratio_pct."
            required:
              - "month"
              - "is_projection"
              - "income_is_expected"
              - "credit_payments_cents"
              - "credit_payments_total_cents"
            x-calculado: true
          description: "Un elemento por mes de `from` a `to`, en orden."
      required:
        - "warning_pct"
        - "limit_pct"
        - "scale_max_pct"
        - "months"
      x-calculado: true
      x-provisional: true
    BudgetSignals:
      type: "object"
      description: "Sugerencia «subir o recortar presupuesto» por categoría (U3 2.0 G-3, métrica M-22, RN-31; K54, motor S2). Mira los months_considered meses cerrados anteriores al mes pedido: raise si en todos el % usado pasó de 100 + margin_pct; cut si en todos quedó abajo de 100 − margin_pct; none en otro caso. Un mes sin plan corta la racha; con plan y sin gasto cuenta 0 %. Solo sugiere: no cambia el presupuesto ni da un monto. Con un period posterior al mes en curso solo cuentan los meses ya cerrados: months puede venir más corto o vacío y entonces signal es none (PROVISIONAL, bloqueo S6-tablero-senal-periodo, opción a). PROVISIONAL (bloqueos U3-tablero-senal-presupuesto y C1-tablero): toda la regla."
      properties:
        month:
          $ref: "componentes-comunes.yaml#/components/schemas/Month"
          description: "Mes pedido (period; por omisión el mes en curso). No entra en la racha."
        months_considered:
          type: "integer"
          description: "Meses cerrados que mira (3 al sembrar; PROVISIONAL, bloqueo U3-tablero-senal-presupuesto). Valor vigente de app_setting budget.signal_months (K44; se cambia en updateDashboardSettings)."
          x-app-setting: "budget.signal_months"
        margin_pct:
          type: "integer"
          description: "Margen en puntos de % usado (10 al sembrar; PROVISIONAL, bloqueo U3-tablero-senal-presupuesto). Valor vigente de app_setting budget.signal_margin_pct (K44; se cambia en updateDashboardSettings)."
          x-app-setting: "budget.signal_margin_pct"
        rows:
          type: "array"
          items:
            type: "object"
            description: "Categoría."
            properties:
              category_id:
                type: "integer"
                description: "category.id (categoría activa)."
              signal:
                type: "string"
                enum:
                  - "raise"
                  - "cut"
                  - "none"
                description: "raise = subir presupuesto; cut = recortar; none = sin sugerencia."
              months:
                type: "array"
                items:
                  type: "object"
                  description: "Mes mirado."
                  properties:
                    month:
                      $ref: "componentes-comunes.yaml#/components/schemas/Month"
                      description: "Mes."
                    planned_cents:
                      $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                      description: "Plan de la categoría en el mes. Ausente sin plan."
                    actual_cents:
                      $ref: "componentes-comunes.yaml#/components/schemas/MoneyCents"
                      description: "Real (parte del hogar devengada, M-09)."
                    used_pct:
                      type: "integer"
                      description: "% usado = actual ÷ planned × 100, entero ROUND_HALF_UP. Ausente sin plan."
                  required:
                    - "month"
                    - "actual_cents"
                  x-calculado: true
                description: "Los meses mirados, del más antiguo al más reciente."
            required:
              - "category_id"
              - "signal"
              - "months"
            x-calculado: true
          description: "Una fila por categoría activa."
      required:
        - "month"
        - "months_considered"
        - "margin_pct"
        - "rows"
      x-calculado: true
      x-provisional: true
    CalendarEventPage:
      allOf:
        - $ref: "componentes-comunes.yaml#/components/schemas/Page"
        - type: "object"
          properties:
            items:
              type: "array"
              items:
                $ref: "#/components/schemas/CalendarEvent"
      description: "Página de CalendarEvent con la envoltura única (Page)."
    ViewPage:
      allOf:
        - $ref: "componentes-comunes.yaml#/components/schemas/Page"
        - type: "object"
          properties:
            items:
              type: "array"
              items:
                $ref: "#/components/schemas/ViewItem"
            mode:
              type: "string"
              enum:
                - "table"
                - "board"
                - "calendar"
                - "gallery"
              description: "Modo pedido."
            groups:
              type: "array"
              items:
                $ref: "#/components/schemas/ViewGroup"
              description: "Columnas del tablero o días del calendario (vacío en tabla y galería)."
          required:
            - "mode"
            - "groups"
      description: "Página de ViewItem con la envoltura única (Page)."

## FIN: docs/contratos/extractos/inicio.yaml

## ARCHIVO: docs/contratos/extractos/INDICE.md


# Extractos del contrato C1 · índice

Generado por `backend/internal/api/c1gen/build.py` (sin opciones o con `--extractos`); `build.py --check` los compara byte a byte con el generador. No se editan a mano.

**Regla (protocolo-agentes §4, K32):** estos archivos son los **únicos** archivos del contrato que leen los agentes de WF-1 y WF-3. `docs/contratos/openapi.yaml` es la fuente única y no se lee entero; si falta algo que no está en un extracto, se pide al integrador.

- Contrato C1 v1.5.1: 221 operaciones en 176 rutas y 267 esquemas.
- Cada operación aparece en exactamente un extracto. Cada extracto trae sus rutas y solo los esquemas que esas rutas usan, resueltos de forma transitiva y copiados en su `components`.
- Parámetros, encabezados, respuestas, seguridad y los esquemas `Error`, `Page`, `MoneyCents`, `Date`, `DateTime` (y los que ellos referencian) están en `componentes-comunes.yaml`; los extractos los citan con `$ref: componentes-comunes.yaml#/components/...`.
- Tope por archivo: 80 KB (80 000 bytes). Etiquetas divididas en partes consecutivas por orden de rutas: «Presupuesto» en `presupuesto-a.yaml`, `presupuesto-b.yaml`; «Compras y tiendas» en `compras-tiendas-a.yaml`, `compras-tiendas-b.yaml`; «Créditos» en `creditos-a.yaml`, `creditos-b.yaml`; «Ajustes» en `ajustes-a.yaml`, `ajustes-b.yaml`, `ajustes-c.yaml`, `ajustes-d.yaml`, `ajustes-e.yaml`.

| Archivo | Etiqueta | Operaciones | Rutas | Esquemas propios | Tamaño |
|---|---|---|---|---|---|
| `inicio.yaml` | Inicio | 15 | 15 | 20 | 66.4 KB (68038 bytes) |
| `gastos.yaml` | Gastos | 19 | 15 | 23 | 75.1 KB (76857 bytes) |
| `presupuesto-a.yaml` | Presupuesto | 13 | 11 | 14 | 45.5 KB (46580 bytes) |
| `presupuesto-b.yaml` | Presupuesto | 15 | 11 | 17 | 44.7 KB (45782 bytes) |
| `comidas.yaml` | Comidas | 16 | 12 | 19 | 48.2 KB (49399 bytes) |
| `compras-tiendas-a.yaml` | Compras y tiendas | 15 | 13 | 24 | 69.0 KB (70665 bytes) |
| `compras-tiendas-b.yaml` | Compras y tiendas | 16 | 13 | 25 | 62.6 KB (64081 bytes) |
| `creditos-a.yaml` | Créditos | 8 | 8 | 10 | 36.0 KB (36826 bytes) |
| `creditos-b.yaml` | Créditos | 11 | 8 | 23 | 62.9 KB (64376 bytes) |
| `salud.yaml` | Salud | 6 | 4 | 7 | 20.4 KB (20892 bytes) |
| `reportes.yaml` | Reportes | 10 | 10 | 21 | 58.4 KB (59768 bytes) |
| `ajustes-a.yaml` | Ajustes | 14 | 12 | 16 | 42.4 KB (43411 bytes) |
| `ajustes-b.yaml` | Ajustes | 14 | 13 | 21 | 52.6 KB (53910 bytes) |
| `ajustes-c.yaml` | Ajustes | 15 | 12 | 22 | 76.5 KB (78329 bytes) |
| `ajustes-d.yaml` | Ajustes | 14 | 13 | 21 | 56.7 KB (58063 bytes) |
| `ajustes-e.yaml` | Ajustes | 20 | 13 | 24 | 72.6 KB (74379 bytes) |
| `componentes-comunes.yaml` | (comunes) | — | — | 7 esquemas, 16 parámetros, 3 encabezados, 11 respuestas | 15.7 KB (16104 bytes) |
| **Total** | | **221** | | | |

## Operaciones por extracto

- `inicio.yaml`: `getHealth`, `getCurrentSession`, `search`, `listCalendarEvents`, `queryView`, `getDashboardMonthVsPlan`, `getDashboardFortnightFlow`, `getDashboardNextPayment`, `getDashboardWeekMeals`, `getDashboardBestBasket`, `getDashboardDebtCapacity`, `getDashboardMonthlyOverview`, `getBudgetSignals`, `getDashboardSavings`, `getPersonParticipation`
- `gastos.yaml`: `listExpenses`, `createExpense`, `getExpenseTotals`, `getExpenseSummary`, `getExpenseMonthlyTrend`, `getExpense`, `updateExpense`, `deleteExpense`, `listExpenseItems`, `duplicateExpense`, `markExpensePaid`, `uploadExpenseAttachment`, `downloadAttachment`, `deleteAttachment`, `createTicketScan`, `listTicketScans`, `getTicketScan`, `saveTicketScanDraft`, `confirmTicketScan`
- `presupuesto-a.yaml`: `listPeriods`, `listBudgets`, `createBudget`, `updateBudget`, `deleteBudget`, `getBudgetMonth`, `getBudgetByCategory`, `getBudgetFortnights`, `getAccruedGroceries`, `getBudgetAlerts`, `getBudgetMonthlyTrend`, `getCashFlowFortnights`, `getAccrualVsCash`
- `presupuesto-b.yaml`: `getCashFlowAssumptions`, `getLiveDebt`, `listIncomes`, `createIncome`, `getIncomeMonth`, `getIncomeMonthlyTrend`, `updateIncome`, `deleteIncome`, `getSavingsSummary`, `getSavingsMonthlyTrend`, `listSavingsGoals`, `listSavingsMovements`, `createSavingsMovement`, `updateSavingsMovement`, `deleteSavingsMovement`
- `comidas.yaml`: `listMealSlots`, `getMealPlanWeek`, `getMealPlanMonth`, `getMealPlanTrends`, `createMealPlanEntries`, `updateMealPlanEntry`, `deleteMealPlanEntry`, `copyPreviousWeek`, `listRecipes`, `createRecipe`, `suggestRecipeWithAi`, `getRecipe`, `updateRecipe`, `deleteRecipe`, `getRecipeUsage`, `duplicateRecipe`
- `compras-tiendas-a.yaml`: `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `listBaskets`, `getBasketComparison`, `listStores`, `listProductGroups`, `listProducts`, `createProduct`, `getProduct`, `updateProduct`, `createPresentation`, `updatePresentation`
- `compras-tiendas-b.yaml`: `listProductPrices`, `listNutritionFacts`, `createNutritionFact`, `getPriceSummary`, `createPrice`, `updatePrice`, `deletePrice`, `getExternalAgentReference`, `importExternalCatalog`, `importExternalPrices`, `importExternalRecipes`, `listInventoryItems`, `createInventoryItem`, `getInventorySummary`, `listInventoryMovements`, `adjustInventoryItem`
- `creditos-a.yaml`: `listCredits`, `getCreditTotals`, `getDebtTrend`, `getCredit`, `previewAvailability`, `getAvailabilitySeries`, `getCardStatus`, `listCreditTermRates`
- `creditos-b.yaml`: `listCreditPurchases`, `createCreditPurchase`, `previewCreditPurchase`, `getCreditPurchase`, `updateCreditPurchase`, `deleteCreditPurchase`, `listScheduledPayments`, `markScheduledPaymentPaid`, `deleteCreditPayment`, `listPromotions`, `compareScenarios`
- `salud.yaml`: `listHealthProfiles`, `getTodayKcal`, `listWeightEntries`, `createWeightEntry`, `updateWeightEntry`, `deleteWeightEntry`
- `reportes.yaml`: `listReports`, `generateMonthlyReport`, `generateExternalStatement`, `downloadReportFile`, `sendReport`, `listMonthlyReportMonths`, `getMonthlyReport`, `getExternalStatement`, `listNotificationDeliveries`, `listJobs`
- `ajustes-a.yaml`: `getOpenApiDocument`, `login`, `logout`, `changeOwnPassword`, `listUsers`, `listSessions`, `revokeSession`, `revokeSessions`, `getProjectionSettings`, `updateProjectionSettings`, `listPlannedItems`, `createPlannedItem`, `updatePlannedItem`, `previewPlannedItemsImport`
- `ajustes-b.yaml`: `importPlannedItems`, `createSavingsGoal`, `updateSavingsGoal`, `listCategories`, `createCategory`, `updateCategory`, `reorderCategories`, `listPersons`, `updatePerson`, `listCostShareRules`, `updateCostShareRule`, `createStore`, `updateStore`, `setStoreProductGroups`
- `ajustes-c.yaml`: `getPriceSources`, `getPriceSettings`, `updatePriceSettings`, `importProfeco`, `getRecipeSettings`, `reloadNutritionBase`, `updateInventoryItem`, `createCredit`, `updateCredit`, `createCreditTermRate`, `updateCreditTermRate`, `deleteCreditTermRate`, `createPromotion`, `updatePromotion`, `deletePromotion`
- `ajustes-d.yaml`: `getOptimizerSettings`, `updateOptimizerSettings`, `createHealthProfile`, `updateHealthProfile`, `updateJob`, `listNotifications`, `updateNotification`, `setNotificationRecipients`, `sendTestNotification`, `listNotificationRecipients`, `listTelegramChatCandidates`, `confirmTelegramChatCandidate`, `dismissTelegramChatCandidate`, `updateNotificationRecipient`
- `ajustes-e.yaml`: `listBackups`, `createBackup`, `restoreBackup`, `getSeedStatus`, `getSecuritySettings`, `updateSecuritySettings`, `getAlertSettings`, `updateAlertSettings`, `getDashboardSettings`, `updateDashboardSettings`, `getReportSettings`, `updateReportSettings`, `getDataSettings`, `updateDataSettings`, `getHealthSettings`, `updateHealthSettings`, `listAiProviders`, `updateAiProvider`, `getSystemStatus`, `getCaCertificate`

## FIN: docs/contratos/extractos/INDICE.md

## ARCHIVO: frontend/package.json


{
  "name": "finanzas-hogar-frontend",
  "private": true,
  "version": "0.1.0",
  "type": "module",
  "description": "PWA de Finanzas Hogar (WF-3). Contrato C1 v1.4, tokens C2.",
  "scripts": {
    "tokens": "node scripts/tokens.mjs",
    "iconos": "node scripts/iconos.mjs",
    "tipos": "openapi-typescript ../docs/contratos/openapi.yaml -o src/api/esquema.gen.ts",
    "mocks:gen": "node scripts/generar-mocks.mjs",
    "preparar": "npm run tokens && npm run iconos && npm run tipos && npm run mocks:gen",
    "dev": "npm run preparar && vite",
    "dev:mocks": "npm run preparar && vite --mode mocks",
    "typecheck": "tsc --noEmit",
    "build": "npm run preparar && tsc --noEmit && vite build",
    "preview": "vite preview --port 4173 --strictPort",
    "test:e2e": "playwright test",
    "verificar:rutas": "node scripts/rutas-api.mjs",
    "verificar:hosts": "node scripts/hosts-externos.mjs dist",
    "verificar:c3": "python3 -I scripts/c3-integracion.py",
    "test:integracion": "playwright test -c playwright.integracion.config.ts",
    "verify": "npm run verificar",
    "verificar": "npm run build && npm run verificar:rutas && npm run verificar:hosts && python3 -m json.tool dist/manifest.webmanifest > /dev/null && npm run test:e2e",
    "storybook": "npm run tokens && STORYBOOK_DISABLE_TELEMETRY=1 storybook dev -p 6006 --host 127.0.0.1 --no-open",
    "build-storybook": "npm run tokens && STORYBOOK_DISABLE_TELEMETRY=1 storybook build -o storybook-static --quiet",
    "typecheck:variantes": "tsc --noEmit -p .storybook/tsconfig.json",
    "verificar:variantes": "npm run build-storybook && node scripts/contraste-variantes.mjs"
  },
  "dependencies": {
    "@dnd-kit/core": "6.3.1",
    "@dnd-kit/sortable": "10.0.0",
    "@dnd-kit/utilities": "3.2.2",
    "@tanstack/react-table": "9.2.8",
    "class-variance-authority": "0.7.1",
    "clsx": "2.1.1",
    "echarts": "6.1.0",
    "lucide-react": "1.55.0",
    "radix-ui": "1.8.0",
    "react": "19.3.0",
    "react-dom": "19.3.0",
    "tailwind-merge": "3.7.0"
  },
  "devDependencies": {
    "@axe-core/playwright": "4.13.0",
    "@playwright/test": "1.60.0",
    "@storybook/react-vite": "10.6.1",
    "@tailwindcss/vite": "4.3.3",
    "@types/node": "24.10.1",
    "@types/react": "19.3.0",
    "@types/react-dom": "19.3.0",
    "@vitejs/plugin-react": "6.1.2",
    "msw": "3.0.3",
    "openapi-typescript": "7.13.0",
    "storybook": "10.6.1",
    "tailwindcss": "4.3.3",
    "typescript": "5.9.3",
    "vite": "8.3.4",
    "vite-plugin-pwa": "2.0.0",
    "workbox-window": "7.4.1",
    "yaml": "2.9.1"
  }
}

## FIN: frontend/package.json

## ARCHIVO: frontend/src/textos/es-MX.ts


// Todos los textos visibles de la base del frontend, en español de México (U5-07, protocolo §8).
// Las áreas agregan sus textos en src/pantallas/<area>/textos.ts siguiendo este mismo estilo.
export const textos = {
  app: { nombre: 'Finanzas Hogar', cargando: 'Cargando…', tardando: 'Esto está tardando más de lo normal', reintentar: 'Reintentar', volver: 'Volver', cerrar: 'Cerrar', cancelar: 'Cancelar', guardar: 'Guardar', confirmar: 'Confirmar', mas: 'Más', masOpciones: 'Más opciones', menosOpciones: 'Menos opciones', deshacer: 'Deshacer', seDeshizo: 'Se deshizo', saltarContenido: 'Saltar al contenido' },
  secciones: { inicio: 'Inicio', gastos: 'Gastos', presupuesto: 'Presupuesto', comidas: 'Comidas', compras: 'Compras y tiendas', creditos: 'Créditos', salud: 'Salud', reportes: 'Reportes', ajustes: 'Ajustes' },
  subsecciones: { resumen: 'Resumen', detalle: 'Detalle', tendencias: 'Tendencias', calendario: 'Calendario' },
  vistas: { grupo: 'Vistas', columnaAcciones: 'Acciones', tabla: 'Tabla', tablero: 'Tablero', calendario: 'Calendario', galeria: 'Galería', filtrar: 'Filtrar', ordenar: 'Ordenar', agrupar: 'Agrupar', sinGrupo: 'Sin agrupar', quitarFiltro: (n: string) => `Quitar filtro ${n}`, ascendente: 'ascendente', descendente: 'descendente', sinValor: 'Sin valor', registros: (n: number) => (n === 1 ? '1 registro' : `${n} registros`), mesAnterior: 'Mes anterior', mesSiguiente: 'Mes siguiente', sinRegistrosMes: (m: string) => `Sin registros en ${m}`, antes: 'Antes', despues: 'Después', guardadoEnLugar: 'Guardado', cambiaA: 'cambia a',
    seleccionar: (n: string) => `Seleccionar ${n}`, seleccionarTodas: 'Seleccionar todas', seleccionadas: (n: number) => (n === 1 ? '1 seleccionada' : `${n} seleccionadas`), quitarSeleccion: 'Quitar selección',
    editarSeleccionadas: 'Editar seleccionadas', atajoEditar: 'Ctrl+/', eligeFilas: 'Selecciona filas para editarlas a la vez', edicionGrupo: (n: number) => (n === 1 ? 'Editar 1 registro' : `Editar ${n} registros a la vez`),
    propiedad: 'Propiedad', valorNuevo: 'Valor nuevo', eligeOpcion: 'Elige un valor', eligeValor: 'Elige un valor nuevo para ver el efecto', calculando: 'Calculando el efecto…', noSeCalculo: 'No se pudo calcular el efecto',
    sinEfecto: 'Sin cambios en las cifras', guardando: 'Guardando…', ayudaGrupo: 'Enter aplica a todas · Esc cancela sin cambiar ninguna', aplicarA: (n: number) => (n === 1 ? 'Aplicar a 1' : `Aplicar a ${n}`), guardadoVarias: (n: number) => (n === 1 ? 'Guardado en 1 registro' : `Guardado en ${n} registros`) },
  barra: { buscar: 'Buscar', compras: 'Compras', registrarGasto: 'Registrar gasto', plegar: 'Ocultar barra lateral', desplegar: 'Mostrar barra lateral', abrirMenu: 'Abrir menú', desplegarSeccion: (s: string) => `Desplegar ${s}`, plegarSeccion: (s: string) => `Plegar ${s}`, cuenta: 'Cuenta', apariencia: 'Apariencia', cerrarSesion: 'Cerrar sesión', navegacion: 'Navegación principal', inferior: 'Navegación inferior', agregar: 'Registrar' },
  tema: { titulo: 'Tema', sistema: 'Según el sistema', claro: 'Claro', oscuro: 'Oscuro', cambiar: 'Cambiar tema' },
  menu: { titulo: 'Buscar o escribir «/» para comandos', placeholder: 'Buscar en Finanzas Hogar o escribe «/»', sinResultados: 'Sin resultados', comandos: 'Comandos', paginas: 'Páginas', resultados: 'Resultados' },
  sesion: { titulo: 'Iniciar sesión', usuario: 'Usuario', contrasena: 'Contraseña', entrar: 'Entrar', entrando: 'Entrando…', vencida: 'Tu sesión terminó. Vuelve a entrar para seguir donde estabas.', bloqueo: (hora: string) => `Demasiados intentos. Puedes volver a intentar a las ${hora}.`, requiereConexion: 'Requiere conexión con el servidor de casa' },
  conexion: { sinConexion: 'Sin conexión con el servidor de casa', enCola: (n: number) => `${n} en cola`, noSeEnvio: 'No se pudo enviar', revisar: 'Revisar', requiere: 'Requiere conexión con el servidor de casa', paginaNoCargada: 'Sin conexión · Esta página aún no se ha cargado en este teléfono', cargadoA: (h: string) => `Cargado a las ${h}`, enviadas: (n: number) => (n === 1 ? '1 captura enviada' : `${n} capturas enviadas`) },
  cola: { titulo: 'Capturas en el teléfono', vacia: 'No hay capturas pendientes.', pendiente: 'Pendiente de envío', fallo: 'No se pudo enviar', reintentar: 'Reintentar envío', descartar: 'Descartar', confirmarDescartar: '¿Descartar esta captura? Se borrará del teléfono y no se enviará.', gasto: 'Gasto', peso: 'Peso', capturadoEl: (f: string) => `Capturado el ${f}` },
  registro: { titulo: 'Registrar', tipos: { gasto: 'Gasto', ingreso: 'Ingreso', peso: 'Peso', ticket: 'Ticket' }, monto: 'Monto', montoInvalido: 'Escribe un monto mayor que cero.', categoria: 'Categoría', elegirCategoria: 'Elige una categoría.', todasCategorias: 'Todas las categorías', fecha: 'Fecha', persona: 'Persona', formaPago: 'Forma de pago', descripcion: 'Descripción o comercio', nota: 'Nota', formas: { cash: 'De contado', debit: 'Débito', transfer: 'Transferencia', credit: 'A crédito' }, creditoAviso: 'La compra a crédito se registra en Créditos con sus pagos.', irCredito: 'Ir a compra a crédito', gastoGuardado: 'Gasto guardado', gastoEnCola: 'Gasto guardado en el teléfono · Se enviará al volver a la red de casa', pesoGuardado: 'Peso guardado', pesoEnCola: 'Peso guardado en el teléfono · Se enviará al volver a la red de casa', ingresoIr: 'El ingreso se captura en Presupuesto › Ingresos.', irIngresos: 'Ir a Ingresos', ticketIr: 'El ticket se captura con la cámara en Gastos › Capturar ticket.', irTicket: 'Capturar ticket', sinCategorias: 'Aún no hay categorías. Créalas en Ajustes › Categorías.', pesoKg: 'Peso (kg)', pesoInvalido: 'Escribe el peso en kilos, por ejemplo 72.5.', ultimoPeso: (p: string, f: string) => `Antes: ${p} el ${f}`, sinPesoPrevio: 'Sin peso anterior' },
  errores: { generico: 'Algo salió mal. Intenta de nuevo.', pantalla: 'No se pudo mostrar esta página', red: 'No hay conexión con el servidor de casa.', noExiste: 'Ya no existe', cambio: 'Este registro cambió', recargar: 'Recargar', grafica: 'No se pudo cargar la gráfica', revisaCampos: 'Revisa los campos marcados.' },
  archivo: { elegir: 'Elegir archivo', elegirCsv: 'Elegir archivo CSV', ninguno: 'Ningún archivo elegido' },
  vacio: { generico: 'Aún no hay datos', grafica: 'Aún no hay datos para este periodo' },
  pendiente: { titulo: 'Pantalla en construcción', cuerpo: (id: string) => `La pantalla ${id} la construye su área en WF-3.` },
  grafica: { tablaAlterna: 'Ver como tabla', otros: 'Otros' },
  /** Barómetro (DESIGN §8.9, U3 B-1): palabra del color y frase; el estado llega de C1, el componente no lo calcula. */
  barometro: {
    estados: { normal: 'Verde · dentro del rango', warning: 'Amarillo · cerca del tope', exceeded: 'Rojo · arriba del tope' },
    colores: { normal: 'verde', warning: 'amarillo', exceeded: 'rojo' },
    deTope: (n: number) => `de ${n} % de tope`,
    sinDato: '—',
  },
} as const

## FIN: frontend/src/textos/es-MX.ts

## ARCHIVO: frontend/src/estilos/app.css


/* Hoja de entrada. tokens.css y tailwind.theme.css son copias exactas de docs/contratos («npm run tokens»). */
@import "tailwindcss";
@import "./tailwind.theme.css";

/* Inter autoalojada (DESIGN §3.1); «npm run tokens» copia los woff2 a public/fonts. */
@font-face { font-family: "Inter"; src: url("/fonts/InterVariable.woff2") format("woff2"); font-weight: 100 900; font-style: normal; font-display: swap; }
@font-face { font-family: "Inter"; src: url("/fonts/InterVariable-Italic.woff2") format("woff2"); font-weight: 100 900; font-style: italic; font-display: swap; }

@layer base {
  html { font-family: var(--font-sans); font-optical-sizing: auto; -webkit-tap-highlight-color: transparent; }
  body { background-color: var(--ui-surface-page); color: var(--ui-text-default); min-height: 100dvh; }
  :focus-visible { outline: 2px solid var(--ui-border-focus); outline-offset: 2px; }
  /* Iconos Lucide (DESIGN §5.2): nunca se encogen en una fila flexible (U5-06: iconos deformados de 7×16 y 19×20). */
  svg.lucide { flex-shrink: 0; }
}

## FIN: frontend/src/estilos/app.css

## ARCHIVO: frontend/src/mocks/manejadores.ts


// Manejadores de msw para TODAS las operaciones de C1 (solo «npm run dev:mocks»; nunca en el build).
// Genéricos: la respuesta de éxito de cada operación con los ejemplos del contrato (operaciones.gen.json).
// Con estado: sesión, personas, categorías, gastos y pesos (idempotentes por Idempotency-Key, K9).
import { http, HttpResponse, type HttpHandler } from 'msw'
import operaciones from './operaciones.gen.json'
import sinteticos from './sinteticos.gen.json'
type Op = { id: string; metodo: string; ruta: string; estado: number; ejemplo: unknown; binario: boolean }
const ahora = () => '2026-10-21T09:00:00-06:00'
const sesion = { user_account_id: 1, username: 'axel', person_id: 1, display_name: 'Axel', participation: 'holder', device_session_id: 1, expires_at: '2027-01-19T09:00:00-06:00', last_seen_at: ahora() }
const pagina = (items: unknown[]) => HttpResponse.json({ items, next_cursor: null })
const gastos: Record<string, unknown>[] = []
const pesos: Record<string, unknown>[] = []
const porClave = new Map<string, Record<string, unknown>>()
function captura(lista: Record<string, unknown>[]) {
  return async ({ request }: { request: Request }) => {
    const clave = request.headers.get('Idempotency-Key')
    if (!clave || clave.length < 16 || clave.length > 128) return HttpResponse.json({ code: 'validation_failed', message: 'Falta la clave de idempotencia.', request_id: 'mock' }, { status: 422 })
    const previo = porClave.get(clave)
    if (previo) return HttpResponse.json(previo, { status: 200, headers: { 'Idempotent-Replayed': 'true' } })
    const nuevo = { id: lista.length + 1, ...(await request.json() as object), created_at: ahora(), updated_at: ahora() }
    lista.unshift(nuevo)
    porClave.set(clave, nuevo)
    return HttpResponse.json(nuevo, { status: 201 })
  }
}
const conEstado: HttpHandler[] = [
  http.post('/api/v1/auth/login', () => HttpResponse.json(sesion)),
  http.get('/api/v1/auth/me', () => HttpResponse.json(sesion)),
  http.post('/api/v1/auth/logout', () => new HttpResponse(null, { status: 204 })),
  http.get('/api/v1/persons', () => pagina(sinteticos.personas)),
  http.get('/api/v1/categories', () => pagina(sinteticos.categorias)),
  http.get('/api/v1/expenses', () => pagina(gastos)),
  http.post('/api/v1/expenses', captura(gastos)),
  http.get('/api/v1/weight-entries', () => pagina(pesos)),
  http.post('/api/v1/weight-entries', captura(pesos)),
  http.get('/health', () => HttpResponse.json({ status: 'ok' })),
]
const genericos = (operaciones as Op[]).map((o) =>
  http[o.metodo as 'get' | 'post' | 'put' | 'patch' | 'delete'](o.ruta.replace(/\{(\w+)\}/g, ':$1'), () =>
    o.estado === 204 ? new HttpResponse(null, { status: 204 }) : o.binario ? new HttpResponse('%PDF-1.7 ejemplo', { status: o.estado, headers: { 'Content-Type': 'application/pdf' } }) : HttpResponse.json(o.ejemplo as object, { status: o.estado }),
  ),
)
// Cada área puede sumar manejadores con estado en src/pantallas/<area>/mocks.ts (export const manejadores); van primero.
const deAreas = Object.values(import.meta.glob<{ manejadores: HttpHandler[] }>('../pantallas/*/mocks.ts', { eager: true })).flatMap((m) => m.manejadores)
export const manejadores = [...deAreas, ...conEstado, ...genericos]

## FIN: frontend/src/mocks/manejadores.ts

## ARCHIVO: frontend/src/app/App.tsx


import { useEffect } from 'react'
import { ProveedorAvisos, useAvisos } from '../componentes/Avisos'
import { ProveedorIconos } from '../componentes/Icono'
import { Cargando } from '../componentes/Estados'
import { iniciarCola, type ResultadoEnvio } from '../cola/cola'
import { textos } from '../textos/es-MX'
import { ProveedorSesion, useSesion } from './sesion'
import { IniciarSesion } from './IniciarSesion'
import { navegar, useUbicacion } from './enrutador'
import { Layout } from '../shell/Layout'
function Raiz() {
  const { estado } = useSesion()
  const { ruta } = useUbicacion()
  const avisar = useAvisos()
  useEffect(() => {
    iniciarCola()
    const alEnviar = (e: Event) => {
      const r = (e as CustomEvent<ResultadoEnvio>).detail
      if (r.enviadas) avisar({ texto: textos.conexion.enviadas(r.enviadas), rol: 'exito' })
      // U5-13 paso 7: tanto el rechazo (4xx) como la falla del servidor o la sesión vencida (5xx, 401) avisan; la captura sigue en la cola.
      if (r.fallidas || r.retenidas) avisar({ texto: `${textos.conexion.noSeEnvio} · ${textos.conexion.revisar}`, rol: 'error' })
    }
    window.addEventListener('fh:cola-enviada', alEnviar)
    return () => window.removeEventListener('fh:cola-enviada', alEnviar)
  }, [avisar])
  useEffect(() => { if (estado === 'activa' && (ruta === '/' || ruta === '/entrar')) navegar('/inicio', { reemplazar: true }) }, [estado, ruta])
  if (estado === 'cargando') return <div className="p-6"><Cargando /></div>
  if (estado === 'sin-sesion') return <IniciarSesion />
  return <Layout />
}
export function App() {
  return <ProveedorIconos><ProveedorAvisos><ProveedorSesion><Raiz /></ProveedorSesion></ProveedorAvisos></ProveedorIconos>
}

## FIN: frontend/src/app/App.tsx

## ARCHIVO: frontend/src/app/IniciarSesion.tsx


// P-00 Iniciar sesión (U1 §3.2): usuario y contraseña; mensaje de bloqueo (429 con retry_at); aviso de sesión vencida.
import { useState, type FormEvent } from 'react'
import { textos } from '../textos/es-MX'
import { Campo } from '../componentes/Campo'
import { Boton } from '../componentes/Boton'
import { Franja } from '../componentes/Estados'
import { ErrorApi, esSinConexion } from '../api/errores'
import { horaCorta } from '../api/formatos'
import { useConexion } from '../api/conexion'
import { useSesion } from './sesion'
export function IniciarSesion() {
  const { entrar, vencida } = useSesion()
  const conectado = useConexion()
  const [usuario, setUsuario] = useState('')
  const [contrasena, setContrasena] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [enviando, setEnviando] = useState(false)
  const alEnviar = async (e: FormEvent) => {
    e.preventDefault()
    setEnviando(true); setError(null)
    try { await entrar(usuario.trim(), contrasena) } catch (err) {
      if (err instanceof ErrorApi && err.estado === 429 && err.reintentarEn) setError(textos.sesion.bloqueo(horaCorta(new Date(err.reintentarEn))))
      else setError(esSinConexion(err) ? textos.sesion.requiereConexion : err instanceof ErrorApi ? err.message : textos.errores.generico)
    } finally { setEnviando(false); setContrasena('') }
  }
  return (
    <main id="contenido" className="mx-auto flex min-h-dvh w-full max-w-sm flex-col justify-center gap-6 px-4">
      <h1 className="text-title-1 font-bold text-text-default">{textos.app.nombre}</h1>
      {vencida && <Franja rol="info">{textos.sesion.vencida}</Franja>}
      {!conectado && <Franja rol="sinConexion">{textos.sesion.requiereConexion}</Franja>}
      <form onSubmit={alEnviar} className="flex flex-col gap-4" aria-label={textos.sesion.titulo}>
        <h2 className="text-title-3 font-semibold">{textos.sesion.titulo}</h2>
        <Campo etiqueta={textos.sesion.usuario} name="username" autoComplete="username" autoCapitalize="none" required minLength={3} maxLength={64} value={usuario} onChange={(e) => setUsuario(e.target.value)} />
        <Campo etiqueta={textos.sesion.contrasena} name="password" type="password" autoComplete="current-password" required value={contrasena} onChange={(e) => setContrasena(e.target.value)} />
        {error && <Franja rol="error" alerta>{error}</Franja>}
        <Boton type="submit" variante="primario" tam="hoja" disabled={enviando}>{enviando ? textos.sesion.entrando : textos.sesion.entrar}</Boton>
      </form>
    </main>
  )
}

## FIN: frontend/src/app/IniciarSesion.tsx

## ARCHIVO: frontend/src/app/enrutador.tsx


// Enrutador mínimo sobre History API (sin dependencias). navegar() y <Enlace> mantienen la SPA.
import { useSyncExternalStore, type ComponentProps } from 'react'
const EVENTO = 'fh:navegacion'
const suscribir = (f: () => void) => { window.addEventListener('popstate', f); window.addEventListener(EVENTO, f); return () => { window.removeEventListener('popstate', f); window.removeEventListener(EVENTO, f) } }
const instantanea = () => location.pathname + location.search
export function useUbicacion() {
  const actual = useSyncExternalStore(suscribir, instantanea)
  const u = new URL(actual, location.origin)
  return { ruta: u.pathname, consulta: u.searchParams, completa: actual }
}
export function navegar(a: string, opciones: { reemplazar?: boolean } = {}) {
  if (a === instantanea()) return
  history[opciones.reemplazar ? 'replaceState' : 'pushState'](null, '', a)
  window.dispatchEvent(new Event(EVENTO))
  window.scrollTo(0, 0)
}
/** Cambia un parámetro de la consulta (sub=, vista=, hoja=) sin perder los demás. */
export function conConsulta(cambios: Record<string, string | null>) {
  const u = new URL(location.href)
  for (const [k, v] of Object.entries(cambios)) v === null ? u.searchParams.delete(k) : u.searchParams.set(k, v)
  return u.pathname + (u.search ? u.search : '')
}
export function Enlace({ a, onClick, ...resto }: ComponentProps<'a'> & { a: string }) {
  return <a href={a} onClick={(e) => { onClick?.(e); if (e.defaultPrevented || e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey) return; e.preventDefault(); navegar(a) }} {...resto} />
}

## FIN: frontend/src/app/enrutador.tsx

## ARCHIVO: frontend/src/app/registro.ts


// Registro de pantallas por área: cada área exporta `pantallas` desde src/pantallas/<area>/indice.ts
// (componentes con lazy()). Nadie edita este archivo para agregar una pantalla.
import type { ComponentType } from 'react'
import type { Subseccion } from './rutas'
export interface PropsPantalla { params: Record<string, string>; sub: Subseccion | null; consulta: URLSearchParams }
type Indice = { pantallas: Record<string, ComponentType<PropsPantalla>> }
const modulos = import.meta.glob<Indice>('../pantallas/*/indice.ts', { eager: true })
export const PANTALLAS_REGISTRADAS: Record<string, ComponentType<PropsPantalla>> = Object.assign({}, ...Object.values(modulos).map((m) => m.pantallas))

## FIN: frontend/src/app/registro.ts

## ARCHIVO: frontend/src/app/rutas.ts


// Pantallas P-00..P-26 (U1 §2.1) y su ruta en el navegador. Las rutas son PROVISIONALES
// (bloqueo docs/estado/bloqueos/WF3-frontend-rutas.md): U1 define pantallas, no direcciones.
import { textos } from '../textos/es-MX'
export type Subseccion = 'resumen' | 'detalle' | 'tendencias' | 'calendario'
export type IdSeccion = keyof typeof textos.secciones
export interface Pantalla { id: string; nombre: string; seccion: IdSeccion | null; ruta: string; subs?: Subseccion[]; sub?: Subseccion; padre?: string }
const R = 'resumen', D = 'detalle', T = 'tendencias', C = 'calendario'
export const PANTALLAS: Pantalla[] = [
  { id: 'P-00', nombre: 'Iniciar sesión', seccion: null, ruta: '/entrar' },
  { id: 'P-01', nombre: 'Inicio', seccion: 'inicio', ruta: '/inicio', subs: [R, D, T, C], sub: R },
  { id: 'P-05', nombre: 'Capturar ticket', seccion: 'gastos', ruta: '/gastos/ticket' },
  { id: 'P-04', nombre: 'Gasto', seccion: 'gastos', ruta: '/gastos/:id', padre: 'P-02' },
  { id: 'P-02', nombre: 'Todos los gastos', seccion: 'gastos', ruta: '/gastos', subs: [R, D, T, C], sub: D },
  { id: 'P-07', nombre: 'Flujo y proyección', seccion: 'presupuesto', ruta: '/presupuesto/flujo', subs: [R, D, T], sub: R },
  { id: 'P-08', nombre: 'Ingresos', seccion: 'presupuesto', ruta: '/presupuesto/ingresos', subs: [R, D, T, C], sub: D },
  { id: 'P-09', nombre: 'Ahorro y metas', seccion: 'presupuesto', ruta: '/presupuesto/ahorro', subs: [R, D, T, C], sub: R },
  { id: 'P-06', nombre: 'Presupuesto del mes', seccion: 'presupuesto', ruta: '/presupuesto', subs: [R, D, T, C], sub: R },
  { id: 'P-12', nombre: 'Receta', seccion: 'comidas', ruta: '/comidas/recetas/:id', padre: 'P-11' },
  { id: 'P-11', nombre: 'Recetario', seccion: 'comidas', ruta: '/comidas/recetas' },
  { id: 'P-10', nombre: 'Plan de la semana', seccion: 'comidas', ruta: '/comidas', subs: [R, D, T, C], sub: D },
  { id: 'P-14', nombre: 'Comparar canasta', seccion: 'compras', ruta: '/compras/comparar', subs: [R, D], sub: R },
  { id: 'P-15', nombre: 'Productos y precios', seccion: 'compras', ruta: '/compras/productos', subs: [R, D, T], sub: D },
  { id: 'P-16', nombre: 'Inventario', seccion: 'compras', ruta: '/compras/inventario', subs: [R, D], sub: D },
  { id: 'P-13', nombre: 'Lista de compras', seccion: 'compras', ruta: '/compras' },
  { id: 'P-20', nombre: 'Comparador de escenarios', seccion: 'creditos', ruta: '/creditos/comparador' },
  { id: 'P-19', nombre: 'Compra a crédito', seccion: 'creditos', ruta: '/creditos/:id/compras/:compra', padre: 'P-18' },
  { id: 'P-18', nombre: 'Crédito', seccion: 'creditos', ruta: '/creditos/:id', subs: [R, D, T, C], sub: R, padre: 'P-17' },
  { id: 'P-17', nombre: 'Todas las líneas', seccion: 'creditos', ruta: '/creditos', subs: [R, D, T, C], sub: R },
  { id: 'P-21', nombre: 'Salud', seccion: 'salud', ruta: '/salud', subs: [R, D, T], sub: R },
  { id: 'P-24', nombre: 'Reporte mensual', seccion: 'reportes', ruta: '/reportes/mensual' },
  { id: 'P-25', nombre: 'Estado de cuenta de Karen', seccion: 'reportes', ruta: '/reportes/karen' },
  { id: 'P-23', nombre: 'Centro de reportes', seccion: 'reportes', ruta: '/reportes' },
  { id: 'P-26', nombre: 'Ajustes', seccion: 'ajustes', ruta: '/ajustes/:seccion' },
  { id: 'P-26', nombre: 'Ajustes', seccion: 'ajustes', ruta: '/ajustes' },
]
/** Las 14 secciones de Ajustes (U1 §2.3), en orden. */
export const SECCIONES_AJUSTES = ['Cuenta y roles', 'Presupuesto', 'Ingresos', 'Categorías', 'Créditos y condiciones', 'Tiendas y precios', 'Recetas', 'Reportes y frecuencias', 'Notificaciones', 'Apariencia', 'Datos y respaldos', 'Personas y reparto', 'Salud y metas de ahorro', 'Inventario'].map((nombre, i) => ({ id: String(i + 1).padStart(2, '0'), nombre }))
/** Barra lateral (U1 §1.1–§1.2): orden fijo de las 9 secciones y sus subpáginas. */
export const SECCIONES: { id: IdSeccion; hijos: string[] }[] = [
  { id: 'inicio', hijos: [] },
  { id: 'gastos', hijos: ['P-02', 'P-05'] },
  { id: 'presupuesto', hijos: ['P-06', 'P-07', 'P-08', 'P-09'] },
  { id: 'comidas', hijos: ['P-10', 'P-11'] },
  { id: 'compras', hijos: ['P-13', 'P-14', 'P-15', 'P-16'] },
  { id: 'creditos', hijos: ['P-17', 'P-20'] },
  { id: 'salud', hijos: [] },
  { id: 'reportes', hijos: ['P-23', 'P-24', 'P-25'] },
  { id: 'ajustes', hijos: [] },
]
export const PRIMERA: Record<IdSeccion, string> = { inicio: 'P-01', gastos: 'P-02', presupuesto: 'P-06', comidas: 'P-10', compras: 'P-13', creditos: 'P-17', salud: 'P-21', reportes: 'P-23', ajustes: 'P-26' }
const compiladas = PANTALLAS.map((p) => ({ p, re: new RegExp('^' + p.ruta.replace(/:(\w+)/g, '(?<$1>[^/]+)') + '/?$') }))
export function resolver(ruta: string): { pantalla: Pantalla; params: Record<string, string> } | null {
  for (const { p, re } of compiladas) { const m = re.exec(ruta); if (m) return { pantalla: p, params: { ...m.groups } } }
  return null
}
export const pantalla = (id: string) => PANTALLAS.find((p) => p.id === id)!
export function rutaDe(id: string, params: Record<string, string | number> = {}): string {
  return pantalla(id).ruta.replace(/:(\w+)/g, (_, k: string) => encodeURIComponent(String(params[k] ?? '')))
}

## FIN: frontend/src/app/rutas.ts

## ARCHIVO: frontend/src/app/sesion.tsx


// Carga de sesión (K11: sesión larga por dispositivo con cookie HttpOnly fh_session; el frontend nunca ve el token).
// Sin conexión, usa la última sesión cargada en el teléfono para poder capturar gastos y peso (K9).
import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from 'react'
import { pedir } from '../api/cliente'
import { esSinConexion } from '../api/errores'
import { borrarLecturas, guardarLectura, leerLectura } from '../api/lecturas'
import type { Esquemas } from '../api/tipos'
import { navegar } from './enrutador'
export type Sesion = Esquemas['CurrentSession']
type Estado = 'cargando' | 'sin-sesion' | 'activa'
interface Ctx { sesion: Sesion | null; estado: Estado; vencida: boolean; entrar: (usuario: string, contrasena: string) => Promise<void>; salir: () => Promise<void> }
const Contexto = createContext<Ctx | null>(null)
export const useSesion = () => useContext(Contexto)!
const CLAVE_VOLVER = 'fh-volver-a'
export function ProveedorSesion({ children }: { children: ReactNode }) {
  const [sesion, setSesion] = useState<Sesion | null>(null)
  const [estado, setEstado] = useState<Estado>('cargando')
  const [vencida, setVencida] = useState(false)
  useEffect(() => {
    pedir('get', '/api/v1/auth/me').then(
      (s) => { setSesion(s); setEstado('activa') },
      async (e) => {
        const copia = esSinConexion(e) ? await leerLectura('/api/v1/auth/me') : undefined
        if (copia) { setSesion(copia.datos as Sesion); setEstado('activa') } else setEstado('sin-sesion')
      },
    )
    const alVencer = () => {
      try { sessionStorage.setItem(CLAVE_VOLVER, location.pathname + location.search) } catch { /* nada */ }
      setVencida(true); setSesion(null); setEstado('sin-sesion')
    }
    window.addEventListener('fh:sesion-vencida', alVencer)
    return () => window.removeEventListener('fh:sesion-vencida', alVencer)
  }, [])
  const entrar = useCallback(async (username: string, password: string) => {
    const s = await pedir('post', '/api/v1/auth/login', { cuerpo: { username, password } })
    // U5-13 paso 2 (H-5): la sesión recién abierta queda como la «última carga» de /auth/me; así, si se recarga sin red
    // antes de cualquier otra navegación completa, la app arranca con los datos guardados y la franja «Sin conexión».
    await guardarLectura('/api/v1/auth/me', s)
    setSesion(s); setEstado('activa'); setVencida(false)
    let volver = '/inicio'
    try { volver = sessionStorage.getItem(CLAVE_VOLVER) ?? volver; sessionStorage.removeItem(CLAVE_VOLVER) } catch { /* nada */ }
    navegar(volver.startsWith('/entrar') ? '/inicio' : volver, { reemplazar: true })
  }, [])
  const salir = useCallback(async () => {
    try { await pedir('post', '/api/v1/auth/logout') } finally {
      await borrarLecturas() // la cola NO se borra: las capturas pendientes nunca se pierden (U1 §2.4)
      setSesion(null); setEstado('sin-sesion'); navegar('/entrar', { reemplazar: true })
    }
  }, [])
  return <Contexto.Provider value={{ sesion, estado, vencida, entrar, salir }}>{children}</Contexto.Provider>
}

## FIN: frontend/src/app/sesion.tsx

## ARCHIVO: frontend/src/app/tema.ts


// Tema: por defecto según el sistema, con interruptor (decisión G, U5-10). Se guarda por dispositivo.
// Claro y oscuro los resuelven las variables --ui-* de tokens.css con [data-theme]; nunca inversión de color.
import { useSyncExternalStore } from 'react'
export type Tema = 'sistema' | 'claro' | 'oscuro'
const CLAVE = 'fh-tema'
const oyentes = new Set<() => void>()
export function leerTema(): Tema { try { const t = localStorage.getItem(CLAVE); return t === 'claro' || t === 'oscuro' ? t : 'sistema' } catch { return 'sistema' } }
export function ponerTema(t: Tema) {
  try { if (t === 'sistema') localStorage.removeItem(CLAVE); else localStorage.setItem(CLAVE, t) } catch { /* sin almacenamiento */ }
  if (t === 'sistema') document.documentElement.removeAttribute('data-theme')
  else document.documentElement.setAttribute('data-theme', t === 'claro' ? 'light' : 'dark')
  oyentes.forEach((f) => f())
}
/** /tema alterna sistema → claro → oscuro → sistema (U1 §5.2). */
export const siguienteTema = (t: Tema): Tema => (t === 'sistema' ? 'claro' : t === 'claro' ? 'oscuro' : 'sistema')
export const useTema = () => useSyncExternalStore((f) => { oyentes.add(f); return () => oyentes.delete(f) }, leerTema)

## FIN: frontend/src/app/tema.ts

## ARCHIVO: frontend/src/api/cliente.ts


// Cliente del API C1. Única puerta de red del frontend: pedir(metodo, ruta-del-contrato, opciones).
// Convenciones: errores C1 §3 → ErrorApi; sin respuesta → ErrorSinConexion; paginación C1 §5 (cursor);
// Idempotency-Key obligatoria por tipo en los POST de captura; dinero en centavos; fechas ISO.
import { ErrorApi, ErrorSinConexion } from './errores'
import { marcarConServidor, marcarSinServidor } from './conexion'
import { guardarLectura } from './lecturas'
import type { Consulta, Cuerpo, ItemsDe, Metodo, ParamsRuta, RequiereClave, Respuesta, RutasCon } from './tipos'

export type Opciones<P extends RutasCon<M>, M extends Metodo> = {
  ruta?: ParamsRuta<P, M>
  consulta?: Consulta<P, M>
  cuerpo?: Cuerpo<P, M>
  senal?: AbortSignal
} & (RequiereClave<P, M> extends true ? { idempotencia: string } : { idempotencia?: string })

/** Clave de idempotencia (16–128, [A-Za-z0-9_-]): UUID v4, o 32 hex si el navegador no tiene randomUUID. */
export function nuevaClave(): string {
  if (typeof crypto.randomUUID === 'function') return crypto.randomUUID()
  return [...crypto.getRandomValues(new Uint8Array(16))].map((b) => b.toString(16).padStart(2, '0')).join('')
}

export function armarUrl(plantilla: string, ruta?: object, consulta?: object): string {
  let url = plantilla.replace(/\{(\w+)\}/g, (_, k: string) => encodeURIComponent(String((ruta as Record<string, unknown>)?.[k] ?? '')))
  const q = new URLSearchParams()
  for (const [k, v] of Object.entries(consulta ?? {})) {
    if (v === undefined || v === null || v === '') continue
    q.set(k, Array.isArray(v) ? v.join(',') : String(v)) // style=form, explode=false (C1)
  }
  const s = q.toString()
  return s ? `${url}?${s}` : url
}

export async function pedir<M extends Metodo, P extends RutasCon<M>>(metodo: M, ruta: P, opciones?: Opciones<P, M>): Promise<Respuesta<P, M>> {
  const o = (opciones ?? {}) as { ruta?: object; consulta?: object; cuerpo?: unknown; senal?: AbortSignal; idempotencia?: string }
  const url = armarUrl(ruta as string, o.ruta, o.consulta)
  const encabezados: Record<string, string> = { Accept: 'application/json' }
  let cuerpo: BodyInit | undefined
  if (o.cuerpo instanceof FormData) cuerpo = o.cuerpo
  else if (o.cuerpo !== undefined) { cuerpo = JSON.stringify(o.cuerpo); encabezados['Content-Type'] = 'application/json' }
  if (o.idempotencia) encabezados['Idempotency-Key'] = o.idempotencia
  let res: Response
  try {
    res = await fetch(url, { method: metodo.toUpperCase(), headers: encabezados, body: cuerpo, credentials: 'same-origin', signal: o.senal, cache: 'no-store' })
  } catch (e) {
    if ((e as Error).name === 'AbortError') throw e
    marcarSinServidor()
    throw new ErrorSinConexion()
  }
  marcarConServidor()
  const tipo = res.headers.get('Content-Type') ?? ''
  const datos = res.status === 204 ? undefined : tipo.includes('json') ? await res.json() : await res.blob()
  if (!res.ok) {
    const e = (datos ?? {}) as { code?: string; message?: string; fields?: []; request_id?: string; retry_at?: string }
    if (res.status === 401 && ruta !== '/api/v1/auth/login') window.dispatchEvent(new Event('fh:sesion-vencida'))
    throw new ErrorApi(res.status, e.code ?? 'internal_error', e.message ?? 'Algo salió mal. Intenta de nuevo.', e.fields ?? [], e.request_id ?? res.headers.get('X-Request-Id'), e.retry_at ?? null)
  }
  if (metodo === 'get' && tipo.includes('json')) void guardarLectura(url, datos)
  return datos as Respuesta<P, M>
}

/** Una página (C1 §5): `limit` 50 por omisión, máximo 200 (K6); `next_cursor` null = no hay más. */
export function pedirPagina<P extends RutasCon<'get'>>(metodo: 'get', ruta: P, opciones?: Opciones<P, 'get'>) {
  return pedir(metodo, ruta, opciones)
}

/** Todas las páginas de una colección, en orden (para listas cortas: categorías, personas, tiendas). */
export async function todas<P extends RutasCon<'get'>>(metodo: 'get', ruta: P, opciones?: Opciones<P, 'get'>): Promise<ItemsDe<Respuesta<P, 'get'>>[]> {
  const items: ItemsDe<Respuesta<P, 'get'>>[] = []
  let cursor: string | null = null
  do {
    const consulta = { ...((opciones as { consulta?: object })?.consulta ?? {}), limit: 200, ...(cursor ? { cursor } : {}) }
    const pagina = (await pedir(metodo, ruta, { ...(opciones as object), consulta } as unknown as Opciones<P, 'get'>)) as unknown as { items: ItemsDe<Respuesta<P, 'get'>>[]; next_cursor: string | null }
    items.push(...pagina.items)
    cursor = pagina.next_cursor
  } while (cursor)
  return items
}

## FIN: frontend/src/api/cliente.ts

## ARCHIVO: frontend/src/api/conexion.ts


// Estado de conexión con el servidor de casa (G-07, N-004): navigator.onLine + último resultado de red.
import { useSyncExternalStore } from 'react'
let conServidor = typeof navigator === 'undefined' ? true : navigator.onLine
const oyentes = new Set<() => void>()
let sondeo: ReturnType<typeof setInterval> | null = null
const avisar = () => oyentes.forEach((f) => f())
export const hayConexion = () => conServidor && navigator.onLine
export function marcarConServidor() {
  if (conServidor) return
  conServidor = true
  if (sondeo) clearInterval(sondeo)
  sondeo = null
  avisar()
  window.dispatchEvent(new Event('fh:conexion-recuperada'))
}
export function marcarSinServidor() {
  if (!conServidor) return
  conServidor = false
  avisar()
  // Sondeo de /health cada 15 s mientras no haya servidor (PROPUESTA K31: valor configurable).
  sondeo ??= setInterval(() => {
    if (!navigator.onLine) return
    fetch('/health', { cache: 'no-store' }).then((r) => r.ok && marcarConServidor(), () => undefined)
  }, 15_000)
}
if (typeof window !== 'undefined') {
  window.addEventListener('online', () => { conServidor = false; marcarConServidor() })
  window.addEventListener('offline', () => { conServidor = true; marcarSinServidor() })
}
const suscribir = (f: () => void) => { oyentes.add(f); return () => oyentes.delete(f) }
export const useConexion = () => useSyncExternalStore(suscribir, hayConexion)

## FIN: frontend/src/api/conexion.ts

## ARCHIVO: frontend/src/api/consulta.ts


// Hook de lectura con estados (carga, error, sin conexión con copia o sin copia) y recarga.
import { useCallback, useEffect, useRef, useState } from 'react'
import { armarUrl, pedir, type Opciones } from './cliente'
import { esSinConexion } from './errores'
import { leerLectura } from './lecturas'
import type { Respuesta, RutasCon } from './tipos'

export interface EstadoConsulta<T> {
  datos: T | undefined
  cargando: boolean
  error: unknown
  /** Si los datos vienen de la última carga del teléfono (sin conexión): cuándo se cargaron. */
  desdeCopia: number | null
  /** true cuando no hay conexión y esta lectura nunca se cargó en el teléfono. */
  nuncaCargada: boolean
  recargar: () => void
}

export function usarConsulta<P extends RutasCon<'get'>>(metodo: 'get', ruta: P, opciones?: Opciones<P, 'get'>, activa = true): EstadoConsulta<Respuesta<P, 'get'>> {
  const clave = JSON.stringify([ruta, opciones ?? {}])
  const [estado, setEstado] = useState<Omit<EstadoConsulta<Respuesta<P, 'get'>>, 'recargar'>>({ datos: undefined, cargando: activa, error: null, desdeCopia: null, nuncaCargada: false })
  const [vuelta, setVuelta] = useState(0)
  const opRef = useRef(opciones)
  opRef.current = opciones
  useEffect(() => {
    if (!activa) return
    const ctl = new AbortController()
    setEstado((e) => ({ ...e, cargando: true, error: null }))
    pedir(metodo, ruta, { ...(opRef.current as object), senal: ctl.signal } as Opciones<P, 'get'>).then(
      (datos) => setEstado({ datos, cargando: false, error: null, desdeCopia: null, nuncaCargada: false }),
      async (error) => {
        if (ctl.signal.aborted) return
        if (esSinConexion(error)) {
          const o = opRef.current as { ruta?: object; consulta?: object } | undefined
          const copia = await leerLectura(armarUrl(ruta as string, o?.ruta, o?.consulta))
          if (copia) return setEstado({ datos: copia.datos as Respuesta<P, 'get'>, cargando: false, error: null, desdeCopia: copia.cargadoEn, nuncaCargada: false })
          return setEstado({ datos: undefined, cargando: false, error, desdeCopia: null, nuncaCargada: true })
        }
        setEstado({ datos: undefined, cargando: false, error, desdeCopia: null, nuncaCargada: false })
      },
    )
    const recuperada = () => setVuelta((v) => v + 1)
    window.addEventListener('fh:conexion-recuperada', recuperada)
    return () => { ctl.abort(); window.removeEventListener('fh:conexion-recuperada', recuperada) }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [clave, vuelta, activa])
  const recargar = useCallback(() => setVuelta((v) => v + 1), [])
  return { ...estado, recargar }
}

## FIN: frontend/src/api/consulta.ts

## ARCHIVO: frontend/src/api/errores.ts


// Errores del cliente (C1 §3: esquema único Error con code, message, fields, request_id, retry_at).
import type { Esquemas } from './tipos'
export type CampoConError = Esquemas['FieldError']
export class ErrorApi extends Error {
  constructor(
    readonly estado: number,
    readonly codigo: string,
    mensaje: string,
    readonly campos: CampoConError[] = [],
    readonly idPeticion: string | null = null,
    readonly reintentarEn: string | null = null,
  ) {
    super(mensaje)
    this.name = 'ErrorApi'
  }
  /** Mensaje de un campo (`amount_cents`, …) para pintarlo junto al control. */
  campo(nombre: string): string | undefined {
    return this.campos.find((c) => c.field === nombre)?.message
  }
}
/** Sin respuesta del servidor de casa (red cortada, servidor apagado o fuera de la LAN). */
export class ErrorSinConexion extends Error {
  constructor() {
    super('sin_conexion')
    this.name = 'ErrorSinConexion'
  }
}
export const esSinConexion = (e: unknown): e is ErrorSinConexion => e instanceof ErrorSinConexion

## FIN: frontend/src/api/errores.ts

## ARCHIVO: frontend/src/api/formatos.ts


// Formatos de pantalla (U1 §0.4, DESIGN §3.3). Dinero en centavos enteros; nunca punto flotante.
const MESES = ['ene', 'feb', 'mar', 'abr', 'may', 'jun', 'jul', 'ago', 'sep', 'oct', 'nov', 'dic']
const DIAS = ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb']
export const ZONA = 'America/Mexico_City'
export const MENOS = '−'

/** 123456 → «$1,234.56»; con signo: 'gasto' → «−$1,234.56», 'ingreso' → «+$1,234.56»; negativos siempre con U+2212. */
export function dinero(centavos: number, signo: 'ninguno' | 'gasto' | 'ingreso' = 'ninguno'): string {
  if (!Number.isSafeInteger(centavos)) throw new Error('centavos no enteros')
  const neg = centavos < 0
  const abs = Math.abs(centavos)
  const enteros = String(Math.floor(abs / 100)).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
  const texto = `$${enteros}.${String(abs % 100).padStart(2, '0')}`
  if (neg || (signo === 'gasto' && abs > 0)) return MENOS + texto
  return signo === 'ingreso' && abs > 0 ? '+' + texto : texto
}

/** Dinero compacto, SOLO en ejes de gráficas (DESIGN §7.7 y §8.5): 0 → «$0», 50000 → «$500», 500000 → «$5 mil»,
 *  1234567 → «$12.3 mil»; negativos con U+2212. Aritmética entera con ROUND_HALF_UP; en etiquetas, tablas y
 *  descripciones va siempre `dinero()` (formato completo). */
export function dineroCompacto(centavos: number): string {
  const pesos = Math.round(Math.abs(centavos) / 100), signo = centavos < 0 && pesos > 0 ? MENOS : ''
  if (pesos < 1000) return `${signo}$${pesos}`
  const decimas = Math.round(pesos / 100)
  return `${signo}$${Math.floor(decimas / 10)}${decimas % 10 ? `.${decimas % 10}` : ''} mil`
}

/** Texto tecleado → centavos con ROUND_HALF_UP («185.5» → 18550; «1,000.005» → 100001). null si no es válido. */
export function aCentavos(texto: string): number | null {
  const limpio = texto.replace(/[$,\s]/g, '')
  const m = /^(\d{1,12})(?:\.(\d*))?$/.exec(limpio)
  if (!m) return null
  const dec = (m[2] ?? '').padEnd(3, '0')
  let c = Number(m[1]) * 100 + Number(dec.slice(0, 2))
  if (Number(dec[2]) >= 5) c += 1
  return c
}

/** Kilos tecleados → gramos enteros («72.5» → 72500; ROUND_HALF_UP al gramo). */
export function aGramos(texto: string): number | null {
  const m = /^(\d{1,3})(?:\.(\d*))?$/.exec(texto.replace(/[\s,]/g, '.').replace(/\.(?=.*\.)/g, ''))
  if (!m) return null
  const dec = (m[2] ?? '').padEnd(4, '0')
  let g = Number(m[1]) * 1000 + Number(dec.slice(0, 3))
  if (Number(dec[3]) >= 5) g += 1
  return g > 0 ? g : null
}
/** 72500 → «72.5 kg» (sin ceros de sobra; precisión exacta del registro). */
export const kilos = (g: number) => `${(g / 1000).toFixed(3).replace(/\.?0+$/, '')} kg`

const partes = (fecha: string) => {
  const [a, m, d] = fecha.slice(0, 10).split('-').map(Number)
  return { a, m, d, dia: new Date(Date.UTC(a, m - 1, d)).getUTCDay() }
}
/** «2026-10-10» → «10-oct-2026». */
export function fecha(iso: string): string {
  const { a, m, d } = partes(iso)
  return `${d}-${MESES[m - 1]}-${a}`
}
/** «2026-10-30» → «Vie 30 oct 2026». */
export function fechaConDia(iso: string): string {
  const { a, m, d, dia } = partes(iso)
  return `${DIAS[dia]} ${d} ${MESES[m - 1]} ${a}`
}
/** «2026-10» → «oct-2026». */
export const mes = (aaaamm: string) => `${MESES[Number(aaaamm.slice(5, 7)) - 1]}-${aaaamm.slice(0, 4)}`
/** Porcentaje de avance entero con ROUND_HALF_UP: 79.5 → «80 %». */
export const porcentaje = (n: number) => `${Math.floor(n + 0.5)} %`

function enZona(instante = new Date()) {
  const f = new Intl.DateTimeFormat('en-CA', { timeZone: ZONA, year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit', hourCycle: 'h23' })
  const p = Object.fromEntries(f.formatToParts(instante).map((x) => [x.type, x.value]))
  return { fecha: `${p.year}-${p.month}-${p.day}`, hora: `${p.hour}:${p.minute}:${p.second}`, p }
}
/** Hoy en America/Mexico_City (AAAA-MM-DD), para valores por defecto. */
export const hoyISO = (instante = new Date()) => enZona(instante).fecha
/** Instante ISO 8601 con desplazamiento de America/Mexico_City (C1 DateTime; `captured_at`). */
export function ahoraConDesplazamiento(instante = new Date()): string {
  const { fecha: f, hora, p } = enZona(instante)
  const local = Date.UTC(+p.year, +p.month - 1, +p.day, +p.hour, +p.minute, +p.second)
  const min = Math.round((local - Math.floor(instante.getTime() / 1000) * 1000) / 60000)
  const s = min < 0 ? '-' : '+'
  const a = Math.abs(min)
  return `${f}T${hora}${s}${String(Math.floor(a / 60)).padStart(2, '0')}:${String(a % 60).padStart(2, '0')}`
}
/** Hora local «21:30» de un instante (para «Cargado a las 21:30»). */
export const horaCorta = (instante: Date | number) => enZona(new Date(instante)).hora.slice(0, 5)
/** Fecha-hora de API («2026-10-10T21:30:00-06:00») → «10-oct-2026 21:30». */
export const fechaHora = (iso: string) => `${fecha(enZona(new Date(iso)).fecha)} ${horaCorta(new Date(iso))}`

## FIN: frontend/src/api/formatos.ts

## ARCHIVO: frontend/src/api/lecturas.ts


// Última carga de cada lectura (GET) en el teléfono: «lo ya cargado sigue visible con su hora de carga» (U1 §2.4).
// Se borra al cerrar sesión. Las escrituras nunca se guardan aquí (solo gastos y peso van a la cola, K9).
import { leer, poner, quitar } from '../cola/idb'
export interface Lectura { url: string; datos: unknown; cargadoEn: number }
export async function guardarLectura(url: string, datos: unknown) {
  try { await poner<Lectura>('lecturas', { url, datos, cargadoEn: Date.now() }) } catch { /* sin IndexedDB: se sigue sin copia */ }
}
export async function leerLectura(url: string): Promise<Lectura | undefined> {
  try { return await leer<Lectura>('lecturas', url) } catch { return undefined }
}
export const borrarLecturas = () => quitar('lecturas').catch(() => undefined)

## FIN: frontend/src/api/lecturas.ts

## ARCHIVO: frontend/src/api/tipos.ts


// Tipos derivados del contrato C1 (src/api/esquema.gen.ts lo genera openapi-typescript; no se edita).
import type { paths, components } from './esquema.gen'
export type { paths, components }
export type Esquemas = components['schemas']
export type Metodo = 'get' | 'post' | 'put' | 'patch' | 'delete'
type Op<P extends keyof paths, M extends Metodo> = NonNullable<paths[P][M]>
export type RutasCon<M extends Metodo> = { [P in keyof paths]: [Op<P, M>] extends [never] ? never : P }[keyof paths]
type Json<R> = R extends { content: { 'application/json': infer J } } ? J : undefined
export type Respuesta<P extends keyof paths, M extends Metodo> = Op<P, M> extends { responses: infer R }
  ? R extends { 200: infer A } ? Json<A> : R extends { 201: infer A } ? Json<A> : R extends { 202: infer A } ? Json<A> : undefined
  : never
export type Cuerpo<P extends keyof paths, M extends Metodo> = Op<P, M> extends { requestBody?: infer B }
  ? NonNullable<B> extends { content: { 'application/json': infer J } } ? J : NonNullable<B> extends { content: { 'multipart/form-data': unknown } } ? FormData : never
  : never
export type Consulta<P extends keyof paths, M extends Metodo> = Op<P, M> extends { parameters: { query?: infer Q } } ? NonNullable<Q> : never
export type ParamsRuta<P extends keyof paths, M extends Metodo> = Op<P, M> extends { parameters: { path?: infer Q } } ? NonNullable<Q> : never
type Encabezados<P extends keyof paths, M extends Metodo> = Op<P, M> extends { parameters: { header?: infer H } } ? NonNullable<H> : never
export type RequiereClave<P extends keyof paths, M extends Metodo> = [Encabezados<P, M>] extends [never] ? false : Encabezados<P, M> extends { 'Idempotency-Key': string } ? true : false
export type ItemsDe<R> = R extends { items: (infer I)[] } ? I : never

## FIN: frontend/src/api/tipos.ts

## ARCHIVO: frontend/src/componentes/Avisos.tsx


// Avisos breves (DESIGN §7.15, G-07): uno visible a la vez; éxito 5 s (pausa con hover o foco); error sin cierre automático.
import { createContext, useCallback, useContext, useEffect, useRef, useState, type ReactNode } from 'react'
import { CircleAlert, CircleCheck, Info, TriangleAlert, Undo2, X } from 'lucide-react'
import { ICONO_MD } from './Icono'
import { textos } from '../textos/es-MX'
export type RolAviso = 'exito' | 'error' | 'advertencia' | 'info'
export interface Aviso { texto: string; rol?: RolAviso; accion?: { texto: string; alHacer: () => void | Promise<void> } }
const Ctx = createContext<(a: Aviso) => void>(() => undefined)
export const useAvisos = () => useContext(Ctx)
const ICONO = { exito: CircleCheck, error: CircleAlert, advertencia: TriangleAlert, info: Info }
const COLOR = { exito: 'text-role-success', error: 'text-role-error', advertencia: 'text-role-warning', info: 'text-role-info' }
export function ProveedorAvisos({ children }: { children: ReactNode }) {
  const [cola, setCola] = useState<Aviso[]>([])
  const [pausa, setPausa] = useState(false)
  const actual = cola[0]
  const cerrar = useCallback(() => setCola((c) => c.slice(1)), [])
  const mostrar = useCallback((a: Aviso) => setCola((c) => [...c, a]), [])
  const t = useRef<ReturnType<typeof setTimeout>>(undefined)
  useEffect(() => {
    if (!actual || pausa || actual.rol === 'error') return
    t.current = setTimeout(cerrar, 5000)
    return () => clearTimeout(t.current)
  }, [actual, pausa, cerrar])
  const rol = actual?.rol ?? 'exito'
  const Icono = ICONO[rol]
  return (
    <Ctx.Provider value={mostrar}>
      {children}
      <div className="pointer-events-none fixed inset-x-2 bottom-20 z-50 flex justify-center tablet:inset-x-auto tablet:right-6 tablet:bottom-6">
        {actual && (
          <div
            role={rol === 'error' ? 'alert' : 'status'}
            aria-live={rol === 'error' ? 'assertive' : 'polite'}
            onMouseEnter={() => setPausa(true)} onMouseLeave={() => setPausa(false)} onFocus={() => setPausa(true)} onBlur={() => setPausa(false)}
            className="pointer-events-auto flex w-full max-w-120 items-center gap-2 rounded-md border border-border-divider bg-surface-raised px-3 py-2 text-small text-text-default shadow-2"
          >
            <Icono aria-hidden className={`size-5 shrink-0 ${COLOR[rol]}`} />
            <span className="flex-1">{actual.texto}</span>
            {actual.accion && (
              <button type="button" className="inline-flex h-9 items-center gap-1 px-2 font-medium text-text-link coarse:h-11" onClick={async () => { const a = actual.accion!; cerrar(); await a.alHacer() }}>
                {actual.accion.texto === textos.app.deshacer && <Undo2 aria-hidden className="size-5" />}{actual.accion.texto}
              </button>
            )}
            <button type="button" aria-label={textos.app.cerrar} className="inline-flex size-9 items-center justify-center coarse:size-11" onClick={cerrar}><X aria-hidden className="size-6" size={ICONO_MD} /></button>
          </div>
        )}
      </div>
    </Ctx.Provider>
  )
}

## FIN: frontend/src/componentes/Avisos.tsx

## ARCHIVO: frontend/src/componentes/Barometro.tsx


// Barómetro (DESIGN §8.9, U3 B-1, K54): una razón contra dos umbrales en un semicírculo con tramos verde, amarillo y rojo.
// Todo llega de C1 (valor, umbrales, escala y estado): el componente no fija 30 ni 40 ni calcula el estado. SVG en línea con
// colores solo de tokens (clases stroke-/fill- de C2, que cambian solas con el tema); trazos con `non-scaling-stroke` para que
// los grosores sean los de los tokens en píxeles (`size.meter.height`, `size.viz.line`) a cualquier ancho. Las marcas y el
// dígito son HTML (no texto SVG) para que su tamaño sea el de la escala tipográfica y no crezca con el dibujo. Sin animación.
import type { ReactNode } from 'react'
import { CircleAlert, CircleCheck, TriangleAlert } from 'lucide-react'
import { textos } from '../textos/es-MX'
import { cn } from './utils'

export type EstadoBarometro = 'normal' | 'warning' | 'exceeded'
export interface PropsBarometro {
  /** Valor entero en % (C1 `ratio_pct`); ausente = sin dato (sin aguja ni valor, dígito «—»). */
  valor?: number
  /** Límite del tramo verde y del amarillo (C1 `warning_pct`, `limit_pct`). */
  aviso: number
  tope: number
  /** Fin de la escala (C1 `scale_max_pct`). */
  escala: number
  /** Estado que manda C1; ausente sin valor. */
  estado?: EstadoBarometro
  /** Nombre accesible del medidor (título del bloque). */
  titulo: string
  /** `aria-valuetext` completo (título, valor, tope y estado), armado por la pantalla con sus textos; sin dato es el nombre de la imagen. */
  textoAccesible: string
  /** Contenido debajo de la insignia (línea secundaria, avisos). */
  children?: ReactNode
  className?: string
}
const ANCHO = 200, CX = 100, CY = 100, R_ESCALA = 86, R_PISTA = 72
const fraccion = (v: number, escala: number) => Math.min(1, Math.max(0, escala > 0 ? v / escala : 0))
/** Punto del semicírculo para la fracción f (0 = izquierda, 1 = derecha), pasando por arriba. */
const punto = (f: number, r: number) => [CX - r * Math.cos(Math.PI * f), CY - r * Math.sin(Math.PI * f)] as const
const arco = (a: number, b: number, r: number) => {
  const [x1, y1] = punto(a, r), [x2, y2] = punto(b, r)
  return `M ${x1.toFixed(2)} ${y1.toFixed(2)} A ${r} ${r} 0 0 1 ${x2.toFixed(2)} ${y2.toFixed(2)}`
}
const TRAMO: Record<EstadoBarometro, string> = { normal: 'stroke-role-success', warning: 'stroke-role-warning', exceeded: 'stroke-role-error' }
const INSIGNIA: Record<EstadoBarometro, string> = { normal: 'bg-role-success-surface text-role-success', warning: 'bg-role-warning-surface text-role-warning', exceeded: 'bg-role-error-surface text-role-error' }
const ICONO: Record<EstadoBarometro, ReactNode> = { normal: <CircleCheck aria-hidden />, warning: <TriangleAlert aria-hidden />, exceeded: <CircleAlert aria-hidden /> }
const GRUESO = '[vector-effect:non-scaling-stroke] [stroke-width:var(--size-meter-height)]'
const FINO = '[vector-effect:non-scaling-stroke]'

export function Barometro({ titulo, valor, aviso, tope, escala, estado, textoAccesible, children, className }: PropsBarometro) {
  const t = textos.barometro
  const fa = fraccion(aviso, escala), ft = fraccion(tope, escala)
  const conValor = valor !== undefined && estado !== undefined
  const fv = conValor ? fraccion(valor, escala) : 0
  const [ax, ay] = punto(fv, R_ESCALA - 4)
  // Accesibilidad (DESIGN §8.9): con valor, `role="meter"` con `aria-valuenow`; sin dato no hay valor que dar y un `meter` sin
  // `aria-valuenow` es inválido (axe `aria-required-attr`), así que el dibujo es una imagen cuyo nombre es el texto accesible
  // («… sin dato, tope 40 %»). PROVISIONAL: bloqueo docs/estado/bloqueos/WF3c-barometro-sin-dato.md (DESIGN §8.9 no dice el rol sin dato).
  /** Marca de umbral: posición del texto «[n] %» fuera del arco, en % del lienzo (HTML encima del SVG). */
  const marca = (f: number) => { const [x, y] = punto(f, R_ESCALA + 12); return { left: `${(x / ANCHO) * 100}%`, top: `${(y / 110) * 100}%` } }
  return (
    <div className={cn('flex flex-col items-center gap-2', className)}>
      <div {...(conValor
        ? { role: 'meter', 'aria-label': titulo, 'aria-valuemin': 0, 'aria-valuemax': escala, 'aria-valuenow': valor, 'aria-valuetext': textoAccesible }
        : { role: 'img', 'aria-label': textoAccesible })}
        data-prueba="barometro" className="relative w-full max-w-80">
        {/* `fill="none"` en la raíz: ningún hijo hereda el relleno negro inicial; cada trazo pinta con su token (U5-04). */}
        <svg aria-hidden fill="none" viewBox={`0 0 ${ANCHO} 110`} className="block h-auto w-full overflow-visible forced-colors:[&_path]:stroke-[CanvasText]">
          {/* Arco de escala: tres tramos; separación de 2 px del color de la superficie en cada límite (§8.1). */}
          <path d={arco(0, fa, R_ESCALA)} fill="none" className={cn('stroke-role-success', GRUESO)} data-tramo="normal" />
          <path d={arco(fa, ft, R_ESCALA)} fill="none" className={cn('stroke-role-warning', GRUESO)} data-tramo="warning" />
          <path d={arco(ft, 1, R_ESCALA)} fill="none" className={cn('stroke-role-error', GRUESO)} data-tramo="exceeded" />
          {[fa, ft].map((f) => { const [x1, y1] = punto(f, R_ESCALA - 7), [x2, y2] = punto(f, R_ESCALA + 7); return <line key={`s${f}`} x1={x1} y1={y1} x2={x2} y2={y2} className={cn('stroke-surface-card [stroke-width:var(--size-viz-gap)]', FINO)} /> })}
          {/* Marcas de umbral: 1 px en ui.viz.baseline, hacia afuera del arco. */}
          {[fa, ft].map((f) => { const [x1, y1] = punto(f, R_ESCALA + 5), [x2, y2] = punto(f, R_ESCALA + 9); return <line key={`m${f}`} data-marca={f} x1={x1} y1={y1} x2={x2} y2={y2} className={cn('stroke-viz-baseline [stroke-width:1px]', FINO)} /> })}
          {/* Pista del valor (decorativa, DESIGN §10.5 C-172: la mide U5-04 como excluida) y valor en el color del tramo donde cae. */}
          <path d={arco(0, 1, R_PISTA)} fill="none" data-decorativo className={cn('stroke-viz-track', GRUESO)} />
          {conValor && fv > 0 && <path d={arco(0, fv, R_PISTA)} fill="none" className={cn(TRAMO[estado], GRUESO)} data-prueba="barometro-valor" />}
          {/* Aguja con extremo redondeado; su eje va encima en HTML (abajo) para medir `size.viz.marker` a cualquier ancho. */}
          {conValor && <line x1={CX} y1={CY} x2={ax} y2={ay} strokeLinecap="round" data-prueba="barometro-aguja" className={cn('stroke-text-default [stroke-width:var(--size-viz-line)]', FINO)} />}
        </svg>
        {/* Eje de la aguja (DESIGN §8.9; N-4 de la revisión r3 del tablero): punto de `size.viz.marker` (10 px) con anillo de 2 px de
            `ui.surface.card`, en HTML como las marcas: un círculo dentro del SVG crece con el dibujo (medía 16 px a 320 px de ancho). */}
        {conValor && <span aria-hidden data-prueba="barometro-eje" style={{ left: `${(CX / ANCHO) * 100}%`, top: `${(CY / 110) * 100}%` }}
          className="absolute size-[var(--size-viz-marker)] -translate-x-1/2 -translate-y-1/2 rounded-full bg-text-default outline-2 outline-solid outline-surface-card forced-colors:bg-[CanvasText]" />}
        {[aviso, tope].map((n, i) => (
          <span key={`${i}-${n}`} data-prueba="barometro-marca" style={marca(i ? ft : fa)} className="absolute -translate-x-1/2 -translate-y-1/2 whitespace-nowrap text-caption text-text-secondary tabular-nums">{`${n} %`}</span>
        ))}
      </div>
      <div aria-hidden className="flex flex-col items-center">
        <span data-prueba="barometro-digito" className="whitespace-nowrap text-title-1 font-bold text-text-default tabular-nums">{conValor ? `${valor} %` : t.sinDato}</span>
        <span className="text-small text-text-secondary tabular-nums">{t.deTope(tope)}</span>
      </div>
      {conValor && <span className={cn('inline-flex h-6 w-max items-center gap-1 rounded-sm px-2 text-small font-medium [&_svg]:size-5', INSIGNIA[estado])}>{ICONO[estado]}{t.estados[estado]}</span>}
      {children}
    </div>
  )
}

## FIN: frontend/src/componentes/Barometro.tsx

## ARCHIVO: frontend/src/componentes/BarreraErrores.tsx


// Barrera de errores por pantalla (U1 §2.4: «nunca una pantalla en blanco»). Si una pantalla lanza al pintarse, el
// cascarón (barras, franja «Sin conexión», cola y hojas) sigue vivo y se muestra «No se pudo mostrar esta página ·
// Reintentar». Se reinicia sola al cambiar de pantalla (key en Layout) o de subsección o consulta (`clave`).
import { Component, type ErrorInfo, type ReactNode } from 'react'
import { RefreshCw } from 'lucide-react'
import { textos } from '../textos/es-MX'
import { Franja } from './Estados'
import { Boton } from './Boton'
interface Estado { error: unknown }
export class BarreraErrores extends Component<{ children: ReactNode; clave?: string }, Estado> {
  state: Estado = { error: null }
  componentDidUpdate(previas: { clave?: string }) { if (this.state.error && previas.clave !== this.props.clave) this.setState({ error: null }) }
  static getDerivedStateFromError(error: unknown): Estado { return { error: error ?? new Error('sin detalle') } }
  componentDidCatch(error: unknown, info: ErrorInfo) { console.error('Pantalla con error', error, info.componentStack) }
  render() {
    if (!this.state.error) return this.props.children
    return (
      <div data-prueba="barrera-errores">
        <Franja rol="error" alerta accion={<Boton variante="enlace" onClick={() => this.setState({ error: null })}><RefreshCw aria-hidden className="size-5" />{textos.app.reintentar}</Boton>}>{textos.errores.pantalla}</Franja>
      </div>
    )
  }
}

## FIN: frontend/src/componentes/BarreraErrores.tsx

## ARCHIVO: frontend/src/componentes/Boton.tsx


// Botón (DESIGN §7.2) con el patrón de shadcn/ui (cva + cn) y solo tokens C2. El anillo de foco aparece sin transición
// (§7.0, §4.5): la transición cubre color, fondo y borde, nunca outline-color.
import type { ComponentProps } from 'react'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from './utils'
export const estilosBoton = cva(
  'inline-flex select-none items-center justify-center gap-2 rounded-md text-small font-medium transition-[color,background-color,border-color] duration-150 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus aria-disabled:cursor-not-allowed aria-disabled:border aria-disabled:border-dashed aria-disabled:border-border-control aria-disabled:bg-surface-disabled aria-disabled:text-text-disabled',
  {
    variants: {
      variante: {
        primario: 'bg-action-primary-bg text-text-on-primary hover:bg-action-primary-bg-hover active:bg-action-primary-bg-pressed',
        secundario: 'border border-border-control bg-surface-card text-text-default hover:border-border-control-hover hover:bg-surface-hover active:bg-surface-pressed',
        fantasma: 'text-text-default hover:bg-surface-hover active:bg-surface-pressed',
        destructivo: 'bg-action-destructive-bg text-text-on-destructive hover:bg-action-destructive-bg-hover active:bg-action-destructive-bg-pressed',
        enlace: 'px-0 text-text-link underline-offset-2 hover:text-text-link-hover hover:underline',
      },
      tam: { normal: 'h-9 min-w-11 px-3 coarse:h-11', hoja: 'h-12 w-full px-4 text-body', icono: 'size-9 coarse:size-11' },
    },
    defaultVariants: { variante: 'secundario', tam: 'normal' },
  },
)
type Props = ComponentProps<'button'> & VariantProps<typeof estilosBoton> & { /** Motivo en texto: deja el botón visible pero inactivo (§7.0). */ deshabilitadoPor?: string }
export function Boton({ variante, tam, className, deshabilitadoPor, onClick, children, type = 'button', ...resto }: Props) {
  return (
    <button
      type={type}
      className={cn(estilosBoton({ variante, tam }), className)}
      aria-disabled={deshabilitadoPor ? true : undefined}
      title={deshabilitadoPor}
      onClick={deshabilitadoPor ? (e) => e.preventDefault() : onClick}
      {...resto}
    >
      {children}
      {deshabilitadoPor && <span className="sr-only"> · {deshabilitadoPor}</span>}
    </button>
  )
}

## FIN: frontend/src/componentes/Boton.tsx

## ARCHIVO: frontend/src/componentes/Campo.tsx


// Campo de formulario (DESIGN §7.3): etiqueta, control, error con aria-invalid y aria-describedby. Deshabilitado según
// §7.0: relleno ui.surface.disabled, texto ui.text.disabled y borde punteado, sin hover.
import { useId, type ComponentProps } from 'react'
import { CircleAlert } from 'lucide-react'
import { cn } from './utils'
export const estilosControl = 'h-9 w-full rounded-md border border-border-control bg-surface-field px-3 text-body text-text-default placeholder:text-text-secondary enabled:hover:border-border-control-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus aria-invalid:border-border-error coarse:h-11 disabled:cursor-not-allowed disabled:border-dashed disabled:border-border-control disabled:bg-surface-disabled disabled:text-text-disabled'
type Props = ComponentProps<'input'> & { etiqueta: string; error?: string; ayuda?: string; ocultarEtiqueta?: boolean }
export function Campo({ etiqueta, error, ayuda, ocultarEtiqueta, className, id, ...resto }: Props) {
  const auto = useId()
  const idc = id ?? auto
  const desc = error ? `${idc}-error` : ayuda ? `${idc}-ayuda` : undefined
  return (
    <div className={cn('flex flex-col gap-1', className)}>
      <label htmlFor={idc} className={cn('text-small font-medium text-text-secondary', ocultarEtiqueta && 'sr-only')}>{etiqueta}</label>
      <input id={idc} aria-invalid={error ? true : undefined} aria-describedby={desc} className={estilosControl} {...resto} />
      {error ? (
        <p id={desc} className="flex items-center gap-1 text-small text-role-error"><CircleAlert aria-hidden className="size-5" />{error}</p>
      ) : ayuda ? <p id={desc} className="text-small text-text-secondary">{ayuda}</p> : null}
    </div>
  )
}

## FIN: frontend/src/componentes/Campo.tsx

## ARCHIVO: frontend/src/componentes/Chip.tsx


// Chips de persona y tienda (DESIGN §7.4, K41): forma + nombre + posición, nunca solo color.
import type { ComponentProps, ReactNode } from 'react'
import { Check } from 'lucide-react'
import { cn } from './utils'
const PERSONA: Record<string, string> = { axel: 'bg-person-axel', ximena: 'bg-person-ximena', karen: 'bg-person-karen' }
const TIENDA: Record<string, string> = { walmart: 'bg-store-walmart', 'bodega-aurrera': 'bg-store-bodega-aurrera', 'sams-club': 'bg-store-sams-club', market: 'bg-store-market' }
/** Orden fijo de personas y tiendas (DESIGN §7.4). */
export const ORDEN_PERSONAS = ['axel', 'ximena', 'karen', 'shared']
export const ORDEN_TIENDAS = ['walmart', 'bodega-aurrera', 'sams-club', 'market']
export function MarcaPersona({ codigo }: { codigo: string }) {
  return codigo in PERSONA
    ? <span aria-hidden className={cn('inline-block size-2.5 shrink-0 rounded-full', PERSONA[codigo])} />
    : <span aria-hidden className="inline-block size-2.5 shrink-0 rounded-full border-2 border-person-shared" />
}
export function MarcaTienda({ codigo }: { codigo: string }) {
  return <span aria-hidden className={cn('inline-block size-2.5 shrink-0', TIENDA[codigo] ?? 'bg-store-market')} />
}
const base = 'inline-flex h-6 items-center gap-1 border border-border-divider bg-surface-raised px-2 text-small text-text-default'
export const ChipPersona = ({ codigo, nombre }: { codigo: string; nombre: string }) => (
  <span className={cn(base, 'rounded-full')}><MarcaPersona codigo={codigo} />{nombre}</span>
)
export const ChipTienda = ({ codigo, nombre }: { codigo: string; nombre: string }) => (
  <span className={cn(base, 'rounded-sm')}><MarcaTienda codigo={codigo} />{nombre}</span>
)
type Sel = ComponentProps<'button'> & { seleccionado: boolean; marca?: ReactNode; forma?: 'pildora' | 'rectangulo' }
/** Chip seleccionable (filtros y categorías de F-1): <button aria-pressed>, 44 px con puntero grueso. */
export function ChipSeleccionable({ seleccionado, marca, forma = 'pildora', className, children, ...resto }: Sel) {
  return (
    <button
      type="button"
      aria-pressed={seleccionado}
      className={cn(base, 'border-border-control h-8 cursor-pointer px-3 enabled:hover:bg-surface-hover disabled:cursor-not-allowed disabled:border-dashed disabled:border-border-control disabled:bg-surface-disabled disabled:text-text-disabled focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus coarse:h-11', forma === 'pildora' ? 'rounded-full' : 'rounded-sm', seleccionado && 'border-border-selected bg-surface-selected', className)}
      {...resto}
    >
      {seleccionado && <Check aria-hidden className="size-5 text-icon-selected" />}
      {marca}
      {children}
    </button>
  )
}

## FIN: frontend/src/componentes/Chip.tsx

## ARCHIVO: frontend/src/componentes/Estados.tsx


// Estados comunes (U1 §2.4, DESIGN §7.14): carga con esqueleto, vacío, error, sin conexión y franjas.
import { useEffect, useState, type ReactNode } from 'react'
import { CircleAlert, Info, RefreshCw, TriangleAlert, WifiOff } from 'lucide-react'
import { textos } from '../textos/es-MX'
import { ErrorApi, esSinConexion } from '../api/errores'
import { horaCorta } from '../api/formatos'
import type { EstadoConsulta } from '../api/consulta'
import { Boton } from './Boton'
import { cn } from './utils'
export const Esqueleto = ({ className }: { className?: string }) => <div aria-hidden className={cn('h-6 rounded-md bg-skeleton-on-page', className)} />
export function Cargando({ alReintentar, children }: { alReintentar?: () => void; children?: ReactNode }) {
  const [tarde, setTarde] = useState(false)
  useEffect(() => { const t = setTimeout(() => setTarde(true), 10_000); return () => clearTimeout(t) }, [])
  return (
    <div aria-busy="true" aria-label={textos.app.cargando} className="flex flex-col gap-3">
      {children ?? <><Esqueleto className="w-2/3" /><Esqueleto /><Esqueleto className="w-1/2" /></>}
      {tarde && <Franja rol="info" accion={alReintentar && <Boton variante="enlace" onClick={alReintentar}>{textos.app.reintentar}</Boton>}>{textos.app.tardando}</Franja>}
    </div>
  )
}
const ROL = {
  info: { icono: Info, clase: 'bg-role-info-surface', color: 'text-role-info' },
  advertencia: { icono: TriangleAlert, clase: 'bg-role-warning-surface', color: 'text-role-warning' },
  error: { icono: CircleAlert, clase: 'bg-role-error-surface', color: 'text-role-error' },
  sinConexion: { icono: WifiOff, clase: 'bg-role-warning-surface', color: 'text-role-warning' },
}
/** Franja tipo «callout» de Notion. */
export function Franja({ rol, children, accion, alerta }: { rol: keyof typeof ROL; children: ReactNode; accion?: ReactNode; alerta?: boolean }) {
  const { icono: Icono, clase, color } = ROL[rol]
  return (
    <div role={alerta ? 'alert' : 'status'} className={cn('flex flex-col gap-2 rounded-md p-3 text-small text-text-default tablet:flex-row tablet:items-center', clase)}>
      {/* Por debajo de 768 px el texto y la acción se apilan a todo el ancho (densidad cómoda, DESIGN §7.14–§7.15). */}
      <span className="flex min-w-0 flex-1 items-start gap-2 tablet:items-center">
        <Icono aria-hidden className={cn('size-5 shrink-0', color)} />
        <span className="min-w-0 flex-1">{children}</span>
      </span>
      {accion}
    </div>
  )
}
export function Vacio({ icono, texto, accion }: { icono?: ReactNode; texto: string; accion?: ReactNode }) {
  return (
    <div className="flex flex-col items-center gap-3 py-12 text-center text-text-secondary">
      {icono && <span aria-hidden className="text-icon-default [&_svg]:size-6 [&_svg]:stroke-[1.5]">{icono}</span>}
      <p className="text-body">{texto}</p>
      {accion}
    </div>
  )
}
export function ErrorCarga({ error, alReintentar }: { error: unknown; alReintentar?: () => void }) {
  const texto = esSinConexion(error) ? textos.errores.red : error instanceof ErrorApi ? (error.estado === 404 ? textos.errores.noExiste : error.message) : textos.errores.generico
  return <Franja rol="error" alerta accion={alReintentar && <Boton variante="enlace" onClick={alReintentar}><RefreshCw aria-hidden className="size-5" />{textos.app.reintentar}</Boton>}>{texto}</Franja>
}
/** Pinta una consulta con todos sus estados; `children` recibe los datos ya cargados. */
export function ConDatos<T>({ consulta, children, esqueleto }: { consulta: EstadoConsulta<T>; children: (d: T) => ReactNode; esqueleto?: ReactNode }) {
  if (consulta.cargando && consulta.datos === undefined) return <Cargando alReintentar={consulta.recargar}>{esqueleto}</Cargando>
  if (consulta.nuncaCargada) return <Franja rol="sinConexion">{textos.conexion.paginaNoCargada}</Franja>
  if (consulta.error || consulta.datos === undefined) return <ErrorCarga error={consulta.error} alReintentar={consulta.recargar} />
  return (
    <>
      {consulta.desdeCopia && <p className="text-caption text-text-secondary">{textos.conexion.cargadoA(horaCorta(consulta.desdeCopia))}</p>}
      {children(consulta.datos)}
    </>
  )
}

## FIN: frontend/src/componentes/Estados.tsx

## ARCHIVO: frontend/src/componentes/Grafica.tsx


// Gráfica base (DESIGN §7.7 y §8.6, K27) con Apache ECharts 6 (renderizado SVG, módulos sueltos).
// Paleta categórica de C2 (--ui-chart-1..8) leída del tema activo; cada serie con marcador, trazo y patrón (decal).
// Accesibilidad (§7.7, U5-07): la descripción accesible es la propia «título. resumen» en español (aria.label.description),
// nunca la automática de ECharts (en inglés y con centavos crudos). La descripción emergente y la guía vertical usan
// tokens de C2 releídos al cambiar el tema; las series no se recolorean ni se atenúan al pasar el puntero (emphasis off).
// Dinero (§7.7 y §8.5): eje Y compacto («$5 mil») solo cuando la gráfica es de dinero; etiquetas directas con el formato
// completo y selectivas (líneas: último valor y extremos; área: valor final; columnas: el valor con ≤ 4 series, y si no
// cabe pasa a «Ver como tabla»); todo el texto del SVG hereda tabular-nums del contenedor.
import { useEffect, useRef, useState } from 'react'
import { ChevronRight } from 'lucide-react'
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { AriaComponent, GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import { SVGRenderer } from 'echarts/renderers'
import { textos } from '../textos/es-MX'
import { dinero, dineroCompacto } from '../api/formatos'
import { Vacio } from './Estados'
echarts.use([BarChart, LineChart, AriaComponent, GridComponent, LegendComponent, TooltipComponent, SVGRenderer])

export const MARCADORES = ['circle', 'rect', 'triangle', 'diamond', 'path://M0,0L10,0L5,10Z', 'path://M4,0H6V4H10V6H6V10H4V6H0V4H4Z', 'path://M1,0L5,4L9,0L10,1L6,5L10,9L9,10L5,6L1,10L0,9L4,5L0,1Z', 'emptyCircle']
export const TRAZOS: (number[] | 'solid')[] = ['solid', [8, 4], [2, 4], [8, 4, 2, 4], [4, 4], [12, 4], [8, 3, 2, 3, 2, 3], [2, 2]]
export function patron(i: number, superficie: string) {
  const rayas = (rot: number, paso: number) => ({ symbol: 'rect', symbolSize: 1, dashArrayX: [1, 0], dashArrayY: [1.5, paso - 1.5], rotation: rot, color: superficie })
  const puntos = (tam: number, paso: number) => ({ symbol: 'circle', symbolSize: tam / paso, dashArrayX: [[paso, 0]], dashArrayY: [paso, 0], color: superficie })
  return [undefined, rayas(Math.PI / 4, 6), rayas(-Math.PI / 4, 6), puntos(2, 6), { symbol: 'rect', symbolSize: 1, dashArrayX: [[1, 7]], dashArrayY: [1, 7], rotation: Math.PI / 4, color: superficie }, rayas(Math.PI / 4, 10), rayas(-Math.PI / 4, 10), puntos(3, 8)][i]
}
const css = (v: string) => getComputedStyle(document.documentElement).getPropertyValue(v).trim()
function usarTemaVivo() {
  const [n, setN] = useState(0)
  useEffect(() => {
    const mq = matchMedia('(prefers-color-scheme: dark)')
    const f = () => setN((x) => x + 1)
    const obs = new MutationObserver(f)
    obs.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] })
    mq.addEventListener('change', f)
    return () => { obs.disconnect(); mq.removeEventListener('change', f) }
  }, [])
  return n
}
/** Índices con etiqueta directa en una línea (DESIGN §8.5): el último valor y los extremos, nunca uno en cada punto. */
export function selectivos(v: (number | null)[]) {
  const idx = v.map((x, i) => [x, i] as const).filter(([x]) => x !== null) as [number, number][]
  if (!idx.length) return new Set<number>()
  const max = idx.reduce((a, b) => (b[0] > a[0] ? b : a)), min = idx.reduce((a, b) => (b[0] < a[0] ? b : a))
  return new Set([idx[idx.length - 1][1], max[1], min[1]])
}
/** Columnas con etiqueta directa solo hasta 4 series (§8.5: «si no cabe, en la tabla»). */
const MAX_SERIES_CON_ETIQUETA = 4
const DINERO = (v: number) => dinero(v)
export interface Serie { nombre: string; valores: number[] }
interface Props {
  titulo: string; resumen: string; tipo: 'barras' | 'lineas' | 'area'; categorias: string[]; series: Serie[]
  /** Formato completo de etiquetas, descripción emergente y tabla. Sin él, la gráfica es de dinero (`dinero`). */
  formato?: (v: number) => string
  /** Formato del eje Y. Por omisión: `dineroCompacto` si la gráfica es de dinero; si no, el mismo `formato` (kilos, kcal). */
  formatoEje?: (v: number) => string
}
export function Grafica({ titulo, resumen, tipo, categorias, series, formato: formatoDado, formatoEje: ejeDado }: Props) {
  const ref = useRef<HTMLDivElement>(null)
  const tema = usarTemaVivo()
  const visibles = series.slice(0, 8) // 8 series máximo; la novena va en «Otros» (la arma quien llama)
  const formato = formatoDado ?? DINERO
  const formatoEje = ejeDado ?? formatoDado ?? dineroCompacto
  // Sin ningún valor distinto de 0 no se dibujan ejes falsos (DESIGN §7.7): estado vacío «Aún no hay datos para este periodo».
  const sinDatos = !visibles.some((s) => s.valores.some((v) => v))
  useEffect(() => {
    if (!ref.current || sinDatos) return
    const graf = echarts.init(ref.current, null, { renderer: 'svg' })
    const sup = css('--ui-surface-card'), sec = css('--ui-text-secondary')
    const columnasConEtiqueta = tipo === 'barras' && visibles.length <= MAX_SERIES_CON_ETIQUETA
    const conEtiqueta = (valores: number[]): Set<number> =>
      tipo === 'lineas' ? selectivos(valores) : tipo === 'area' ? new Set(valores.length ? [valores.length - 1] : []) : new Set()
    graf.setOption({
      aria: { enabled: true, label: { enabled: true, description: `${titulo}. ${resumen}` }, decal: { show: true } },
      textStyle: { fontFamily: 'Inter, ui-sans-serif, system-ui', color: sec },
      // Con leyenda abajo, la rejilla le deja su renglón para que no se encime con las etiquetas del eje X.
      grid: { left: 8, right: 8, top: 24, bottom: visibles.length > 1 ? 32 : 8, containLabel: true },
      legend: visibles.length > 1 ? { type: 'scroll', bottom: 0, textStyle: { color: css('--ui-text-default') }, inactiveColor: css('--ui-text-disabled') } : undefined,
      tooltip: {
        trigger: 'axis',
        className: 'fh-descripcion-grafica',
        valueFormatter: (v: unknown) => formato(Number(v)),
        backgroundColor: css('--ui-surface-raised'),
        borderColor: css('--ui-border-divider'),
        borderWidth: 1,
        padding: 8,
        textStyle: { color: css('--ui-text-default'), fontFamily: 'Inter, ui-sans-serif, system-ui', fontSize: 14 },
        extraCssText: `border-radius: ${css('--radius-md')}; box-shadow: var(--ui-shadow-1); font-variant-numeric: tabular-nums;`,
        axisPointer: { type: 'line', lineStyle: { color: css('--ui-viz-baseline'), width: 1, type: 'solid' } },
      },
      xAxis: { type: 'category', data: categorias, axisLine: { lineStyle: { color: css('--ui-viz-baseline') } }, axisLabel: { color: sec } },
      yAxis: { type: 'value', splitLine: { lineStyle: { color: css('--ui-viz-grid') } }, axisLabel: { color: sec, formatter: (v: number) => formatoEje(v) } },
      series: visibles.map((s, i) => {
        const marcados = conEtiqueta(s.valores)
        return {
          name: s.nombre,
          type: tipo === 'barras' ? 'bar' : 'line',
          data: s.valores.map((v, j) => (marcados.has(j) ? { value: v, label: { show: true } } : v)),
          label: { show: columnasConEtiqueta, position: 'top', color: css('--ui-text-default'), fontSize: 12, formatter: (p: { value: unknown }) => formato(Number(p.value)) },
          labelLayout: { hideOverlap: true, moveOverlap: 'shiftY' },
          animation: false,
          color: css(`--ui-chart-${i + 1}`),
          symbol: MARCADORES[i],
          symbolSize: 8,
          lineStyle: { type: TRAZOS[i], width: 2 },
          areaStyle: tipo === 'area' ? { opacity: 0.2 } : undefined,
          stack: tipo === 'area' ? 'total' : undefined,
          itemStyle: { decal: patron(i, sup) },
          emphasis: { disabled: true },
        }
      }),
    })
    const ro = new ResizeObserver(() => graf.resize())
    ro.observe(ref.current)
    return () => { ro.disconnect(); graf.dispose() }
  }, [tema, tipo, categorias, visibles, formato, formatoEje, titulo, resumen, sinDatos])
  if (sinDatos) return (
    <figure className="flex flex-col gap-2 rounded-lg border border-border-divider bg-surface-card p-4">
      <figcaption className="text-small font-medium text-text-secondary">{titulo}</figcaption>
      <Vacio texto={textos.vacio.grafica} />
    </figure>
  )
  return (
    <figure className="flex flex-col gap-2 rounded-lg border border-border-divider bg-surface-card p-4">
      <figcaption>
        <span className="block text-small font-medium text-text-secondary">{titulo}</span>
        <span className="block text-small text-text-default">{resumen}</span>
      </figcaption>
      <div ref={ref} tabIndex={0} role="img" aria-label={`${titulo}. ${resumen}`} className="h-50 w-full tabular-nums tablet:h-70" />
      <details className="group text-small">
        <summary className="inline-flex min-h-9 cursor-pointer list-none items-center gap-1 text-text-link coarse:min-h-11 [&::-webkit-details-marker]:hidden">
          <ChevronRight aria-hidden className="size-5 group-open:rotate-90" />{textos.grafica.tablaAlterna}
        </summary>
        <table className="mt-2 w-full tabular-nums">
          <thead><tr><th className="text-left" scope="col">{titulo}</th>{visibles.map((s) => <th key={s.nombre} scope="col" className="text-right">{s.nombre}</th>)}</tr></thead>
          <tbody>{categorias.map((c, j) => <tr key={c}><th scope="row" className="text-left font-normal">{c}</th>{visibles.map((s) => <td key={s.nombre} className="text-right">{formato(s.valores[j] ?? 0)}</td>)}</tr>)}</tbody>
        </table>
      </details>
    </figure>
  )
}

## FIN: frontend/src/componentes/Grafica.tsx

## ARCHIVO: frontend/src/componentes/Hoja.tsx


// Hojas y diálogos (DESIGN §7.12, U1 §1.3): hoja a pantalla completa en móvil, ventana modal en escritorio. Radix Dialog.
import type { ReactNode } from 'react'
import { Dialog } from 'radix-ui'
import { X } from 'lucide-react'
import { ICONO_MD } from './Icono'
import { textos } from '../textos/es-MX'
interface Props { abierta: boolean; alCambiar: (a: boolean) => void; titulo: string; descripcion?: string; children: ReactNode; pie?: ReactNode; alAbrirFoco?: (e: Event) => void }
export function Hoja({ abierta, alCambiar, titulo, descripcion, children, pie, alAbrirFoco }: Props) {
  return (
    <Dialog.Root open={abierta} onOpenChange={alCambiar}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-40 bg-overlay-scrim" />
        <Dialog.Content
          onOpenAutoFocus={alAbrirFoco}
          aria-describedby={descripcion ? undefined : undefined}
          className="fixed inset-0 z-50 flex flex-col bg-surface-raised text-text-default tablet:inset-auto tablet:top-1/2 tablet:left-1/2 tablet:max-h-[85dvh] tablet:w-full tablet:max-w-lg tablet:-translate-x-1/2 tablet:-translate-y-1/2 tablet:rounded-lg tablet:border tablet:border-border-divider tablet:shadow-2"
        >
          <div className="flex items-center gap-2 border-b border-border-divider px-4 py-2">
            <Dialog.Title className="flex-1 text-title-3 font-semibold">{titulo}</Dialog.Title>
            <Dialog.Close aria-label={textos.app.cerrar} className="inline-flex size-9 items-center justify-center rounded-md hover:bg-surface-hover coarse:size-11"><X aria-hidden className="size-6" size={ICONO_MD} /></Dialog.Close>
          </div>
          {descripcion ? <Dialog.Description className="px-4 pt-2 text-small text-text-secondary">{descripcion}</Dialog.Description> : <Dialog.Description className="sr-only">{titulo}</Dialog.Description>}
          <div className="flex-1 overflow-y-auto px-4 py-3">{children}</div>
          {pie && <div className="border-t border-border-divider px-4 py-3 pb-[max(env(safe-area-inset-bottom),12px)]">{pie}</div>}
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}

## FIN: frontend/src/componentes/Hoja.tsx

## ARCHIVO: frontend/src/componentes/Icono.tsx


// Iconos Lucide (DESIGN §5.2): solo dos tamaños, 20 px (`size.icon.sm`) y 24 px (`size.icon.md`), y trazo constante
// de 1.5 px dibujados (`icon.stroke`) con el ancho de trazo absoluto del componente: stroke-width 1.8 a 20 px y 1.5 a
// 24 px en unidades del viewBox de 24 (`icon.stroke.user-units.sm` · `.md`).
// ProveedorIconos fija esos valores por omisión para toda la app (20 px). Un icono de 24 px (solo icono, barras,
// estados vacíos) lleva className="size-6" y size={ICONO_MD}, para que el trazo se calcule sobre 24 px.
import type { ReactNode } from 'react'
import { LucideProvider } from 'lucide-react'
export const ICONO_SM = 20
export const ICONO_MD = 24
export const TRAZO_ICONO = 1.5
export function ProveedorIconos({ children }: { children: ReactNode }) {
  return <LucideProvider size={ICONO_SM} strokeWidth={TRAZO_ICONO} absoluteStrokeWidth>{children}</LucideProvider>
}

## FIN: frontend/src/componentes/Icono.tsx

## ARCHIVO: frontend/src/componentes/SelectorArchivo.tsx


// Selector de archivo con botón propio (U5-07): el <input type="file"> nativo pinta textos del navegador en inglés
// («Choose File», «No file chosen»), así que va oculto y lo abre un Boton con texto en español; el nombre del archivo
// elegido (o «Ningún archivo elegido») se muestra al lado. El input conserva su nombre accesible para las pruebas.
import { useId, useRef, type ReactNode } from 'react'
import { FileUp } from 'lucide-react'
import { textos } from '../textos/es-MX'
import { Boton } from './Boton'
interface Props {
  etiqueta: string
  /** Texto del botón; por omisión «Elegir archivo». */
  boton?: string
  accept: string
  archivo: File | null
  alElegir: (f: File | null) => void
  /** Motivo visible para no dejar elegir (sin conexión, importando…). */
  deshabilitadoPor?: string
  /** Lo que se muestra con un archivo elegido; por omisión, su nombre. */
  detalle?: ReactNode
}
export function SelectorArchivo({ etiqueta, boton = textos.archivo.elegir, accept, archivo, alElegir, deshabilitadoPor, detalle }: Props) {
  const ref = useRef<HTMLInputElement>(null)
  const id = useId()
  return (
    <div role="group" aria-labelledby={id} className="flex flex-col gap-1">
      <span id={id} className="text-small font-medium text-text-default">{etiqueta}</span>
      <div className="flex flex-wrap items-center gap-2">
        <Boton deshabilitadoPor={deshabilitadoPor} onClick={() => ref.current?.click()}><FileUp aria-hidden />{boton}</Boton>
        <span className="min-w-0 break-words text-small text-text-secondary tabular-nums">{archivo ? (detalle ?? archivo.name) : textos.archivo.ninguno}</span>
      </div>
      <input ref={ref} type="file" hidden accept={accept} aria-label={etiqueta} onChange={(e) => { alElegir(e.target.files?.[0] ?? null); e.target.value = '' }} />
    </div>
  )
}

## FIN: frontend/src/componentes/SelectorArchivo.tsx

## ARCHIVO: frontend/src/componentes/TarjetaMetrica.tsx


// Tarjeta de métrica (DESIGN §7.6): etiqueta, cifra exacta con cifras tabulares, nota (también tabular: suele llevar dinero,
// DESIGN §3.4); tocar abre el detalle.
import type { ReactNode } from 'react'
import { cn } from './utils'
type Rol = 'income' | 'expense' | 'savings' | 'debt' | 'warning'
const COLOR: Record<Rol, string> = { income: 'text-role-income', expense: 'text-role-expense', savings: 'text-role-savings', debt: 'text-role-debt', warning: 'text-role-warning' }
interface Props { titulo: string; valor: string; nota?: ReactNode; rol?: Rol; icono?: ReactNode; alAbrir?: () => void; children?: ReactNode }
export function TarjetaMetrica({ titulo, valor, nota, rol, icono, alAbrir, children }: Props) {
  const Contenedor = alAbrir ? 'button' : 'div'
  return (
    <Contenedor
      {...(alAbrir ? { type: 'button' as const, onClick: alAbrir } : {})}
      className={cn('flex w-full flex-col gap-1 rounded-lg border border-border-divider bg-surface-card p-4 text-left', alAbrir && 'border-border-control hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus')}
    >
      <span className="flex items-center gap-2 text-small font-medium text-text-secondary">{icono}{titulo}</span>
      <span className={cn('whitespace-nowrap text-title-2 font-semibold tabular-nums', rol ? COLOR[rol] : 'text-text-default')}>{valor}</span>
      {nota && <span className="text-small text-text-secondary tabular-nums">{nota}</span>}
      {children}
    </Contenedor>
  )
}

## FIN: frontend/src/componentes/TarjetaMetrica.tsx

## ARCHIVO: frontend/src/componentes/utils.ts


import { clsx, type ClassValue } from 'clsx'
import { extendTailwindMerge } from 'tailwind-merge'
// Los tamaños de texto de C2 (text-body, text-small…) no deben chocar con los colores (text-text-default…).
const unir = extendTailwindMerge({ extend: { classGroups: { 'font-size': [{ text: ['caption', 'small', 'body', 'title-1', 'title-2', 'title-3'] }] } } })
export const cn = (...c: ClassValue[]) => unir(clsx(c))

## FIN: frontend/src/componentes/utils.ts

## ARCHIVO: frontend/src/pantallas/inicio/Calendario.tsx


// Calendario de P-01 (U3 §5): pagos del calendario unificado, ingresos y eventos de despensa; quincenas K8; semana de lunes a domingo.
// < 768 px agenda de hoy a hoy + 30 (o + 60); ≥ 768 px mes con franjas de quincena; ≥ 1440 px con el panel del día seleccionado.
import { useState } from 'react'
import { ArrowDownLeft, ArrowUpRight, Calendar, CalendarClock, ChevronLeft, ChevronRight, CreditCard, FileChartColumn, PiggyBank, ShoppingCart, Store } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import { Enlace } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { usarConsulta } from '../../api/consulta'
import { dinero, fecha, fechaConDia } from '../../api/formatos'
import { Boton } from '../../componentes/Boton'
import { cn } from '../../componentes/utils'
import type { Esquemas } from '../../api/tipos'
import { AMPLIO, EnlaceAccion, EsqueletoFilas, Estados, Filtros, IconoInfo, Insignia, InsigniaPago, TABLET, diasAviso, estadoPago, karen, listaDe, nombreCredito, usarMedia, useInicio } from './comun'
import { diaDe, mesDe, quincenaDe, rejillaMes, sumarDias, sumarMeses, ultimoDia } from './fechas'
import { MESES_LARGOS, textosInicio as t } from './textos'
type Evento = Esquemas['CalendarEvent']
type Tipo = Exclude<Evento['event_type'], 'meal_plan' | 'savings_movement'>
const K = t.calendario
const TIPOS: Tipo[] = ['income', 'credit_payment', 'fixed_expense', 'upcoming_expense', 'card_statement', 'planned_grocery', 'market_visit', 'promotion', 'savings_goal', 'monthly_report']
const ICONO: Record<Tipo, typeof Calendar> = { income: ArrowDownLeft, credit_payment: CreditCard, fixed_expense: ArrowUpRight, upcoming_expense: ArrowUpRight, card_statement: CalendarClock, planned_grocery: ShoppingCart, market_visit: Store, promotion: Calendar, savings_goal: PiggyBank, monthly_report: FileChartColumn }
const COLOR: Partial<Record<Tipo, string>> = { income: 'text-role-income', fixed_expense: 'text-role-expense', upcoming_expense: 'text-role-expense', market_visit: 'text-store-market' }
const nums = (v: string[]) => (v.length ? v.map(Number) : undefined)

export default function Calendario() {
  const c = useInicio(), hoy = c.hoy
  const ancho = usarMedia(TABLET), amplio = usarMedia(AMPLIO)
  const [mesV, setMesV] = useState(mesDe(hoy)), [dias, setDias] = useState(30), [sel, setSel] = useState(hoy)
  const [f, setF] = useState<Record<string, string[]>>({ tipos: [], credito: [], familia: [] })
  const period = ancho ? mesV : `${hoy}..${sumarDias(hoy, dias)}`
  const family = f.familia[0] as 'household' | 'external' | 'joint' | undefined
  const ev = usarConsulta('get', '/api/v1/calendar/events', { consulta: { period, event_type: f.tipos.length ? (f.tipos as Tipo[]) : TIPOS, credit_id: nums(f.credito), family, limit: 200 } })
  const pagos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { period, credit_id: nums(f.credito), family, status: 'pending', limit: 200 } })
  const aviso = diasAviso(c), idK = karen(c)?.id
  const porId = new Map(((pagos.datos?.items ?? []) as Esquemas['ScheduledPayment'][]).map((p) => [p.id, p]))
  const grupos = [
    { id: 'tipos', nombre: K.tipos, multiple: true, opciones: TIPOS.map((x) => ({ valor: x, texto: K.tiposEvento[x] })) },
    { id: 'credito', nombre: t.detalle.credito, opciones: listaDe(c.creditos).map((x) => ({ valor: String(x.id), texto: x.name })) },
    { id: 'familia', nombre: t.detalle.familia, opciones: (['household', 'external', 'joint'] as const).map((v) => ({ valor: v, texto: t.detalle.familias[v] })) },
  ]
  function Fila({ e, corto }: { e: Evento; corto?: boolean }) {
    const tipo = e.event_type as Tipo, Icono = ICONO[tipo] ?? Calendar
    const pago = tipo === 'credit_payment' ? estadoPago({ status: e.status === 'paid' || e.status === 'overdue' ? e.status : 'pending', due_date: e.date, shares: e.ref_table === 'scheduled_payment' ? porId.get(e.ref_id)?.shares : undefined }, hoy, aviso, idK) : undefined
    const monto = e.amount_cents === undefined ? '' : tipo === 'income' ? dinero(e.amount_cents, 'ingreso') : tipo === 'fixed_expense' || tipo === 'upcoming_expense' ? dinero(e.amount_cents, 'gasto') : dinero(e.amount_cents)
    const titulo = tipo === 'credit_payment' ? `${nombreCredito(c, e.credit_id)} · ${monto}` : [e.title, monto].filter(Boolean).join(' · ')
    const cuerpo = <>
      <span className={cn('flex items-center gap-2', corto && 'min-w-0')}><Icono aria-hidden className={cn('size-5 shrink-0', COLOR[tipo] ?? 'text-icon-default')} /><span className={cn('tabular-nums', corto ? 'truncate' : 'font-medium', COLOR[tipo])}>{titulo}</span></span>
      {!corto && tipo === 'credit_payment' && e.household_share_cents !== undefined && <span className="text-text-secondary tabular-nums">{t.comun.hogarYKaren(dinero(e.household_share_cents), e.external_share_cents ? dinero(e.external_share_cents) : undefined)}</span>}
      {!corto && tipo === 'promotion' && e.end_date && <span className="text-text-secondary">{K.delAl(fecha(e.date), fecha(e.end_date))}</span>}
      {!corto && pago && <InsigniaPago e={pago} />}
      {!corto && tipo === 'income' && e.status === 'expected' && <Insignia rol="info" icono={<IconoInfo />}>{t.pago.esperado}</Insignia>}
    </>
    const clase = cn('flex flex-col gap-1 text-small', !corto && 'border-t border-border-divider py-2')
    return (tipo === 'credit_payment' || tipo === 'card_statement') && e.credit_id !== undefined && !corto
      ? <li><Enlace a={rutaDe('P-18', { id: e.credit_id })} className={cn(clase, 'min-h-11 rounded-sm text-text-default no-underline hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-border-focus')}>{cuerpo}</Enlace></li>
      : <li className={clase}>{cuerpo}</li>
  }
  const filtros = <Filtros grupos={grupos} valores={f} alCambiar={(id, v) => setF((x) => ({ ...x, [id]: v }))} />
  return (
    <div className="flex flex-col gap-4">
      {filtros}
      <Estados consulta={ev} falla={K.falla} esqueleto={<EsqueletoFilas n={6} />}>{(d) => {
        const items = ([...d.items] as Evento[]).filter((e) => TIPOS.includes(e.event_type as Tipo)).sort((a, b) => a.date.localeCompare(b.date) || TIPOS.indexOf(a.event_type as Tipo) - TIPOS.indexOf(b.event_type as Tipo))
        const delDia = (dia: string) => items.filter((e) => e.date === dia)
        const vacio = !items.length && (listaDe(c.creditos).length === 0
          ? <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{K.sinProgramados}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '03' })} texto={K.accionVacio} /></div>
          : <p className="text-body text-text-secondary">{K.vacio}</p>)
        if (!ancho) {
          const fechas = [...new Set(items.map((e) => e.date))]
          return <div data-prueba="agenda" className="flex flex-col gap-3">{vacio}
            {fechas.map((dia, i) => <section key={dia} aria-label={fechaConDia(dia)} className="flex flex-col">
              {(i === 0 || quincenaDe(fechas[i - 1]) !== quincenaDe(dia)) && <p className="mb-2 rounded-sm bg-surface-hover-subtle px-2 py-1 text-caption font-medium text-text-secondary">{diaDe(dia) <= 15 ? K.q1 : K.q2(ultimoDia(mesDe(dia)))}</p>}
              <h3 className="text-small font-semibold">{fechaConDia(dia)}</h3><ul>{delDia(dia).map((e, j) => <Fila key={j} e={e} />)}</ul></section>)}
            {dias < 60 && <Boton variante="secundario" className="w-max" onClick={() => setDias(60)}>{K.mas30}</Boton>}
          </div>
        }
        const u = ultimoDia(mesV), [a, m] = [Number(mesV.slice(0, 4)), Number(mesV.slice(5, 7))]
        return <div className="grid grid-cols-1 gap-4 wide:grid-cols-3">
          <div data-prueba="mes" className="flex flex-col gap-3 wide:col-span-2">
            <div className="flex flex-wrap items-center gap-2">
              <Boton tam="icono" variante="fantasma" aria-label={K.anterior} onClick={() => setMesV(sumarMeses(mesV, -1))}><ChevronLeft aria-hidden className="size-6" size={ICONO_MD} /></Boton>
              <h2 className="min-w-40 text-title-3 font-semibold">{`${MESES_LARGOS[m - 1]} ${a}`}</h2>
              <Boton tam="icono" variante="fantasma" aria-label={K.siguiente} onClick={() => setMesV(sumarMeses(mesV, 1))}><ChevronRight aria-hidden className="size-6" size={ICONO_MD} /></Boton>
              <Boton variante="secundario" onClick={() => { setMesV(mesDe(hoy)); setSel(hoy) }}>{K.hoy}</Boton>
            </div>
            <div className="flex flex-wrap gap-2 text-caption text-text-secondary"><span className="rounded-sm border border-border-divider bg-surface-card px-2 py-1">{K.q1}</span><span className="rounded-sm bg-surface-hover-subtle px-2 py-1">{K.q2(u)}</span></div>
            {vacio}
            <table className="w-full table-fixed border-collapse text-small">
              <thead><tr>{K.semana.map((s) => <th key={s} scope="col" className="py-1 text-caption font-medium text-text-secondary">{s}</th>)}</tr></thead>
              <tbody>{rejillaMes(mesV).map((sem, i) => <tr key={i}>{sem.map((dia, j) => {
                if (!dia) return <td key={j} className="border border-border-divider" />
                const evs = delDia(dia), q2 = diaDe(dia) > 15
                return <td key={j} className={cn('h-24 border border-border-divider p-0 align-top', q2 && 'bg-surface-hover-subtle')}>
                  <button type="button" aria-pressed={sel === dia} aria-current={dia === hoy ? 'date' : undefined} aria-label={`${fechaConDia(dia)} · ${evs.length}`} onClick={() => setSel(dia)} className={cn('flex h-full min-h-11 w-full flex-col gap-1 p-1 text-left hover:bg-surface-hover focus-visible:outline-2 focus-visible:outline-border-focus', sel === dia && 'outline-2 outline-border-selected')}>
                    <span className={cn('tabular-nums', dia === hoy && 'font-bold text-text-link')}>{diaDe(dia)}</span>
                    {(diaDe(dia) === 1 || diaDe(dia) === 16) && <span className="truncate text-caption text-text-secondary">{diaDe(dia) === 1 ? t.detalle.q1 : t.detalle.q2}</span>}
                    <ul className="flex min-w-0 flex-col gap-1">{evs.slice(0, 3).map((e, k) => <Fila key={k} e={e} corto />)}</ul>
                    {evs.length > 3 && <span className="text-caption text-text-secondary">{K.masN(evs.length - 3)}</span>}
                  </button></td>
              })}</tr>)}</tbody>
            </table>
          </div>
          <section aria-label={K.diaSeleccionado} className={cn('flex flex-col gap-2 rounded-md border border-border-divider bg-surface-card p-4', !amplio && 'tablet:p-6')}>
            <h3 className="text-small font-semibold">{fechaConDia(sel)}</h3>
            {delDia(sel).length ? <ul>{delDia(sel).map((e, j) => <Fila key={j} e={e} />)}</ul> : <p className="text-small text-text-secondary">{K.sinEventosDia}</p>}
          </section>
        </div>
      }}</Estados>
    </div>
  )
}

## FIN: frontend/src/pantallas/inicio/Calendario.tsx

## ARCHIVO: frontend/src/pantallas/inicio/Detalle.tsx


// Detalle de P-01 (U3 §3): D-1 a D-7, cada panel con su petición, sus estados y los filtros que admite (§6).
import { CircleCheck } from 'lucide-react'
import { Enlace, conConsulta, navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { usarConsulta } from '../../api/consulta'
import { dinero, fecha, fechaConDia, porcentaje } from '../../api/formatos'
import { MarcaPersona, MarcaTienda } from '../../componentes/Chip'
import { ChipSeleccionable } from '../../componentes/Chip'
import { EnlaceAccion, EsqueletoFilas, Estados, FallaBloque, Filtros, IconoAviso, IconoInfo, IconoOk, Insignia, InsigniaPago, Panel, codigoTienda, decimas, diasAviso, entero, estadoPago, karen, listaDe, nombreCredito, titulares, useInicio, type GrupoFiltro } from './comun'
import { mesDe, semanaISO, semanaVigente, sumarDias } from './fechas'
import { textosInicio as t } from './textos'
import type { Esquemas } from '../../api/tipos'
type Pago = Esquemas['ScheduledPayment']
import Karen from './Karen'
import type { PropsPantalla } from '../../app/registro'

const F = t.detalle
type Clave = 'periodo' | 'persona' | 'familia' | 'credito' | 'categoria' | 'tienda'
const NOMBRES: Record<Clave, string> = { periodo: F.periodo, persona: F.persona, familia: F.familia, credito: F.credito, categoria: F.categoria, tienda: F.tienda }
const lista = (v: string | null) => (v ? v.split(',').filter(Boolean) : [])
const nums = (v: string[]) => (v.length ? v.map(Number) : undefined)
const progreso = 'h-2 w-full appearance-none overflow-hidden rounded-full bg-viz-track [&::-webkit-progress-bar]:bg-viz-track [&::-webkit-progress-value]:bg-viz-meter [&::-moz-progress-bar]:bg-viz-meter'

export default function Detalle({ consulta }: Pick<PropsPantalla, 'consulta'>) {
  const c = useInicio(), hoy = c.hoy
  const f = { periodo: consulta.get('periodo') ?? 'mes', persona: lista(consulta.get('persona')), familia: lista(consulta.get('familia')), credito: lista(consulta.get('credito')), categoria: lista(consulta.get('categoria')), tienda: lista(consulta.get('tienda')) }
  const period = f.periodo === 'mes' ? mesDe(hoy) : `${mesDe(hoy)}-${f.periodo}`
  const activos = (Object.keys(NOMBRES) as Clave[]).filter((k) => (k === 'periodo' ? f.periodo !== 'mes' : f[k].length > 0))
  const noUsa = (usa: Clave[]) => activos.filter((k) => !usa.includes(k)).map((k) => F.noUsa(NOMBRES[k]))
  const poner = (k: string, v: string[]) => navegar(conConsulta({ [k]: v.length ? v.join(',') : null }), { reemplazar: true })
  const grupos: GrupoFiltro[] = [
    { id: 'persona', nombre: F.persona, opciones: listaDe(c.personas).map((p) => ({ valor: String(p.id), texto: p.display_name, marca: <MarcaPersona codigo={p.code} /> })) },
    { id: 'familia', nombre: F.familia, opciones: (['household', 'external', 'joint'] as const).map((v) => ({ valor: v, texto: F.familias[v] })) },
    { id: 'credito', nombre: F.credito, opciones: listaDe(c.creditos).map((x) => ({ valor: String(x.id), texto: x.name })) },
    { id: 'categoria', nombre: F.categoria, opciones: listaDe(c.categorias).map((x) => ({ valor: String(x.id), texto: x.name })) },
    { id: 'tienda', nombre: F.tienda, opciones: listaDe(c.tiendas).map((x) => ({ valor: String(x.id), texto: x.name, marca: <MarcaTienda codigo={codigoTienda(x)} /> })) },
  ]
  // PROVISIONAL (bloqueo WF3-inicio-vacios): P-02 recibe los filtros con los nombres de C1.
  const aGastos = Object.entries({ period, person_id: f.persona.join(','), family: f.familia.join(','), category_id: f.categoria.join(','), store_id: f.tienda.join(','), credit_id: f.credito.join(',') }).filter(([, v]) => v).map(([k, v]) => `${k}=${encodeURIComponent(v)}`).join('&')
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-2" role="group" aria-label={F.periodo}>
        {(['mes', 'q1', 'q2'] as const).map((p) => <ChipSeleccionable key={p} seleccionado={f.periodo === p} onClick={() => navegar(conConsulta({ periodo: p === 'mes' ? null : p }), { reemplazar: true })}>{F[p]}</ChipSeleccionable>)}
        <Filtros grupos={grupos} valores={{ persona: f.persona, familia: f.familia, credito: f.credito, categoria: f.categoria, tienda: f.tienda }} alCambiar={poner} />
      </div>
      <div className="grid grid-cols-1 gap-4 desktop:grid-cols-2 desktop:gap-6">
        <D1 period={period} categorias={f.categoria} aviso={noUsa(['periodo', 'categoria'])} />
        <D2 aviso={noUsa([])} />
        <D3 persona={f.persona} familia={f.familia} credito={f.credito} aviso={noUsa(['persona', 'familia', 'credito'])} />
        <D4 persona={f.persona} aviso={noUsa(['persona'])} />
        <D5 tiendas={f.tienda} aviso={noUsa(['tienda'])} />
        <D6 aviso={noUsa([])} />
        <Karen aviso={noUsa([])} />
      </div>
      <Enlace a={`${rutaDe('P-02')}?${aGastos}`} className="inline-flex min-h-11 items-center text-small font-medium text-text-link underline underline-offset-2">{F.verGastos}</Enlace>
    </div>
  )
}
function Avisos({ a }: { a: string[] }) { return <>{a.map((x) => <p key={x} className="text-caption text-text-secondary">{x}</p>)}</> }
function D1({ period, categorias, aviso }: { period: string; categorias: string[]; aviso: string[] }) {
  const c = useInicio()
    const q = usarConsulta('get', '/api/v1/budget/categories', { consulta: { period } })
    const nombre = (id: number) => listaDe(c.categorias).find((x) => x.id === id)?.name ?? t.comun.sinDato
    return <Panel etiqueta={t.d1.titulo} titulo={t.d1.titulo}><Avisos a={aviso} />
      <Estados consulta={q} falla={t.d1.falla} esqueleto={<EsqueletoFilas n={5} />}>{(d) => {
        const filas = d.rows.filter((r) => r.planned_cents !== undefined && (!categorias.length || categorias.includes(String(r.category_id)))).sort((a, b) => (b.used_pct ?? -1) - (a.used_pct ?? -1)).slice(0, 5)
        if (!filas.length) return <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d1.vacio}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '02' })} texto={t.d1.accionVacio} /></div>
        return <ul className="flex flex-col gap-3">{filas.map((r) => <li key={r.category_id} className="flex flex-col gap-1 text-small">
          <span className="flex flex-wrap items-center justify-between gap-2"><span className="font-medium">{nombre(r.category_id)}</span>{r.status === 'warning' ? <Insignia rol="advertencia" icono={<IconoAviso />}>{t.t1.vasEn(r.used_pct ?? 0)}</Insignia> : r.status === 'exceeded' ? <Insignia rol="advertencia" icono={<IconoAviso />}>{t.t1.rebasaste}</Insignia> : null}</span>
          <progress className={progreso} max={Math.max(r.planned_cents ?? 0, r.actual_cents, 1)} value={r.actual_cents} aria-label={nombre(r.category_id)} />
          <span className="tabular-nums text-text-secondary">{t.d1.fila(dinero(r.actual_cents), dinero(r.planned_cents ?? 0), r.used_pct !== undefined ? porcentaje(r.used_pct) : t.comun.sinDato)}</span>
        </li>)}</ul>
      }}</Estados><EnlaceAccion a={rutaDe('P-06')} texto={t.acciones.verTodo} /></Panel>
  }
function D2({ aviso }: { aviso: string[] }) {
  const hoy = useInicio().hoy
    const q = usarConsulta('get', '/api/v1/cash-flow/fortnights')
    return <Panel etiqueta={t.d2.titulo} titulo={t.d2.titulo}><Avisos a={aviso} />
      <Estados consulta={q} falla={t.d2.falla} esqueleto={<EsqueletoFilas n={7} />}>{(d) => {
        const r = d.rows.find((x) => x.start_date <= hoy && hoy <= x.end_date)
        if (!r || r.income_cents === 0) return <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d2.vacio}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '03' })} texto={t.d2.accionVacio} /></div>
        const filas: [string, string][] = [[t.d2.ingresos, dinero(r.income_cents, 'ingreso') + (r.income_is_expected ? ` · ${t.d2.esperado}` : '')], [t.d2.pagosCredito, dinero(r.credit_payments_cents, 'gasto')], [t.d2.fijos, dinero(r.fixed_and_upcoming_cents, 'gasto')], [t.d2.variable, dinero(r.variable_estimate_cents, 'gasto')], [t.d2.ahorro, typeof r.planned_savings_cents === 'number' ? dinero(r.planned_savings_cents) : t.comun.sinDato], [t.d2.saldo, dinero(r.balance_cents)], [t.d2.acumulado, dinero(r.cumulative_balance_cents)]]
        return <dl className="grid grid-cols-2 gap-x-3 gap-y-2 text-small">{filas.map(([k, v]) => <div key={k} className="contents"><dt className="text-text-secondary">{k}</dt><dd className="text-right tabular-nums">{v}</dd></div>)}</dl>
      }}</Estados><EnlaceAccion a={rutaDe('P-07')} texto={t.acciones.verTodo} /></Panel>
  }
function D3({ persona, familia, credito, aviso }: { persona: string[]; familia: string[]; credito: string[]; aviso: string[] }) {
  const c = useInicio(), hoy = c.hoy
    const np = c.proximo.datos?.next_payment
    const filtros = { person_id: nums(persona), family: familia[0] as 'household' | 'external' | 'joint' | undefined, credit_id: nums(credito), sort: 'due_date', limit: 5 }
    const vencidos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { ...filtros, status: 'overdue', period: `${np?.due_date ?? hoy}..${sumarDias(hoy, -1)}` } }, np?.is_overdue === true)
    const proximos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { ...filtros, status: 'pending', period: `${hoy}..${sumarDias(hoy, 60)}` } })
    const aviso2 = diasAviso(c), idK = karen(c)?.id
    const av = c.proximo.datos?.availability ?? []
    // «Hogar · Karen» y la insignia del pago dependen de quién es Karen (listPersons). Igual que K-4 (Tablero.tsx): mientras esa
    // lectura no llega no se pinta ninguna fila (esqueleto); sin ella todo caería en el hogar («Hogar $973.75»), una cifra distinta
    // del oráculo (R-13; M-1 de REVISION-frontend-cierre-r2). Si falla, el estado de error de D-3 con «Reintentar»; esta conducta
    // es PROVISIONAL (bloqueo docs/estado/bloqueos/U3-tablero-k4-personas.md). Una parte de una persona que no está en listPersons
    // no se suma al hogar: esa fila se queda sin la línea «Hogar · Karen».
    const sinPersonas = c.personas.datos === undefined && (Boolean(c.personas.error) || c.personas.nuncaCargada)
    const conocidas = new Set(listaDe(c.personas).map((x) => x.id))
    const Fila = ({ p }: { p: Pago }) => {
      const partes = p.shares ?? [], suma = (xs: typeof partes) => xs.reduce((a, s) => a + s.amount_cents, 0)
      const resueltas = partes.every((s) => conocidas.has(s.person_id))
      const k = suma(partes.filter((s) => s.person_id === idK)), h = suma(partes.filter((s) => s.person_id !== idK && conocidas.has(s.person_id)))
      return <li className="flex flex-col gap-1 border-t border-border-divider py-2 text-small tabular-nums">
        <span className="flex flex-wrap items-center justify-between gap-2"><span className="font-medium">{fechaConDia(p.due_date)} · {nombreCredito(c, p.credit_id)}</span><span>{dinero(p.amount_cents)}</span></span>
        <span className="flex flex-wrap items-center justify-between gap-2 text-text-secondary">{resueltas && <span>{t.comun.hogarYKaren(dinero(h), k > 0 ? dinero(k) : undefined)}</span>}<InsigniaPago e={estadoPago(p, hoy, aviso2, idK)} /></span>
      </li>
    }
    return <Panel etiqueta={t.d3.titulo} titulo={t.d3.titulo}><Avisos a={aviso} />
      {sinPersonas ? <FallaBloque texto={t.d3.falla} alReintentar={c.personas.recargar} />
        : c.personas.datos === undefined ? <EsqueletoFilas n={5} />
        : <>
      {np?.is_overdue && <Estados consulta={vencidos} falla={t.d3.falla} esqueleto={<EsqueletoFilas n={2} />}>{(d) => <div><h3 className="text-small font-medium">{t.d3.vencidos}</h3><ul>{(d.items as Pago[]).map((p) => <Fila key={p.id} p={p} />)}</ul>{d.next_cursor && <p className="text-small">{t.d3.masVencidos}</p>}</div>}</Estados>}
      <Estados consulta={proximos} falla={t.d3.falla} esqueleto={<EsqueletoFilas n={5} />}>{(d) => {
        if (!d.items.length && !np?.is_overdue) return listaDe(c.creditos).length || av.length ? <p className="text-body text-text-secondary">{t.d3.sinPagos}</p> : <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d3.sinCreditos}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '05' })} texto={t.d3.accionVacio} /></div>
        return <div><h3 className="text-small font-medium">{t.d3.proximos}</h3><ul>{(d.items as Pago[]).map((p) => <Fila key={p.id} p={p} />)}</ul></div>
      }}</Estados></>}
      {av.length > 0 && <div className="flex flex-col gap-2"><h3 className="text-small font-medium">{t.d3.disponible}</h3>{av.map((a) => <Disponible key={a.credit_id} id={a.credit_id} />)}</div>}
      <EnlaceAccion a={rutaDe('P-17')} texto={t.acciones.verTodo} /></Panel>
  }
function Disponible({ id }: { id: number }) {
  const c = useInicio()
    const q = usarConsulta('get', '/api/v1/credits/{credit_id}/availability-preview', { ruta: { credit_id: id } })
    return <Estados consulta={q} falla={t.d3.fallaDisponible} esqueleto={<EsqueletoFilas n={1} />}>{(d) => <div className="flex flex-col gap-1 text-small tabular-nums">
      <span className="font-medium">{nombreCredito(c, id)}</span>
      <span className="grid grid-cols-1 gap-2 phone:grid-cols-2">{(['A', 'B'] as const).map((l) => { const r = l === 'A' ? d.reading_a : d.reading_b
        return <span key={l} className="flex flex-col gap-1 rounded-sm border border-border-divider p-2"><span className="flex flex-wrap items-center gap-2">{t.comun.lectura(l)} {dinero(r.available_today_cents)}{d.availability_reading === l && <Insignia rol="info" icono={<CircleCheck aria-hidden />}>{t.comun.activa}</Insignia>}</span><span className="text-text-secondary">{t.d3.minimo(dinero(r.horizon_min_cents))}</span></span> })}</span>
      {!d.availability_confirmed && <Insignia rol="info" icono={<IconoInfo />}>{t.t3.reposicion(nombreCredito(c, id))}</Insignia>}
    </div>}</Estados>
  }
function D4({ persona, aviso }: { persona: string[]; aviso: string[] }) {
  const c = useInicio(), hoy = c.hoy
    const q = usarConsulta('get', '/api/v1/meal-plan/week', { consulta: { period: semanaISO(hoy), person_id: nums(persona) } })
    const tiempos = usarConsulta('get', '/api/v1/meal-slots', { consulta: { limit: 200 } })
    const deHoy = (q.datos?.cells ?? []).filter((x) => x.plan_date === hoy && (x.recipe_id !== undefined || x.is_eating_out))
    const recetas = usarConsulta('get', '/api/v1/recipes', { consulta: { limit: 200 } }, deHoy.some((x) => x.recipe_id !== undefined))
    const pers = titulares(c)
    return <Panel etiqueta={t.d4.titulo} titulo={t.d4.titulo}><Avisos a={aviso} />
      <Estados consulta={q} falla={t.d4.falla} esqueleto={<EsqueletoFilas n={4} />}>{() => {
        if (!deHoy.length) return <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d4.vacio}</p><EnlaceAccion a={rutaDe('P-10')} texto={t.d4.accionVacio} /></div>
        const orden = listaDe(tiempos) as Esquemas['MealSlot'][]
        const filas = [...deHoy].sort((a, b) => (orden.find((s) => s.id === a.meal_slot_id)?.sort_order ?? 9) - (orden.find((s) => s.id === b.meal_slot_id)?.sort_order ?? 9) || pers.findIndex((p) => p.id === a.person_id) - pers.findIndex((p) => p.id === b.person_id))
        return <ul>{filas.map((x) => { const p = pers.find((y) => y.id === x.person_id)
          return <li key={`${x.meal_slot_id}-${x.person_id}`} className="flex flex-wrap items-center justify-between gap-2 border-t border-border-divider py-2 text-small tabular-nums">
            <span className="flex items-center gap-2"><span className="font-medium">{orden.find((s) => s.id === x.meal_slot_id)?.name ?? t.comun.sinDato}</span>{p && <><MarcaPersona codigo={p.code} />{p.display_name}</>}</span>
            <span>{x.is_eating_out ? t.d4.fuera : (listaDe(recetas) as { id: number; name: string }[]).find((r) => r.id === x.recipe_id)?.name ?? t.comun.sinDato}</span>
            <span className="text-text-secondary">{x.kcal !== undefined ? t.d4.kcal(entero(x.kcal)) : t.comun.sinDato} · {x.cost_cents !== undefined ? dinero(x.cost_cents) : t.comun.sinDato}</span></li> })}</ul>
      }}</Estados><EnlaceAccion a={rutaDe('P-10')} texto={t.acciones.verTodo} /></Panel>
  }
function D5({ tiendas, aviso }: { tiendas: string[]; aviso: string[] }) {
  const c = useInicio(), hoy = c.hoy
    const week = semanaVigente(hoy)
    const canasta = usarConsulta('get', '/api/v1/dashboard/best-basket', { consulta: { week } })
    const b = canasta.datos
    const q = usarConsulta('get', '/api/v1/basket-comparison', { consulta: { week, basket: b?.basket, expense_id: b?.expense_id } }, b !== undefined && b.store_id !== undefined)
    const ts = listaDe(c.tiendas)
    const vacio = <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d5.vacio}</p><EnlaceAccion a={rutaDe('P-15')} texto={t.d5.accionVacio} /></div>
    return <Panel etiqueta={t.d5.titulo} titulo={t.d5.titulo}><Avisos a={aviso} />
      <Estados consulta={canasta} falla={t.d5.falla} esqueleto={<EsqueletoFilas n={5} />}>{(bb) => bb.store_id === undefined ? vacio :
        <Estados consulta={q} falla={t.d5.falla} esqueleto={<EsqueletoFilas n={5} />}>{(d) => {
          const ver = (id?: number) => !tiendas.length || tiendas.includes(String(id))
          const filas = d.items.filter((x) => x.difference_cents !== undefined && ver(x.winner_store_id)).sort((a, b2) => (b2.difference_cents ?? 0) - (a.difference_cents ?? 0)).slice(0, 5)
          const tienda = (id?: number) => ts.find((s) => s.id === id)
          return <div className="flex flex-col gap-3">
            <ul>{filas.map((x) => { const g = tienda(x.winner_store_id); const viejo = x.prices.some((p) => p.store_id === x.winner_store_id && p.is_stale)
              return <li key={x.product_id} className="flex flex-wrap items-center justify-between gap-2 border-t border-border-divider py-2 text-small tabular-nums">
                <Producto id={x.product_id} /><span className="flex items-center gap-1"><MarcaTienda codigo={codigoTienda(g)} />{g?.name ?? t.comun.sinDato}</span><span>{dinero(x.difference_cents ?? 0)}</span>
                {viejo && <Insignia rol="advertencia" icono={<IconoAviso />}>{t.d5.viejo}</Insignia>}</li> })}</ul>
            <div><h3 className="text-small font-medium">{t.d5.totales}</h3><ul>{[...d.stores].filter((s) => ver(s.store_id)).sort((a, b2) => orden(tienda(a.store_id)) - orden(tienda(b2.store_id))).map((s) => <li key={s.store_id} className="flex items-center justify-between gap-2 py-1 text-small tabular-nums"><span className="flex items-center gap-1"><MarcaTienda codigo={codigoTienda(tienda(s.store_id))} />{tienda(s.store_id)?.name ?? t.comun.sinDato}</span><span>{dinero(s.total_cents)}</span></li>)}</ul></div>
          </div>
        }}</Estados>}</Estados><EnlaceAccion a={rutaDe('P-14')} texto={t.acciones.verTodo} /></Panel>
  }
function D6({ aviso }: { aviso: string[] }) {
    const q = usarConsulta('get', '/api/v1/savings-goals', { consulta: { status: 'active', limit: 50 } })
    return <Panel etiqueta={t.d6.titulo} titulo={t.d6.titulo}><Avisos a={aviso} />
      <Estados consulta={q} falla={t.d6.falla} esqueleto={<EsqueletoFilas n={3} />}>{(d) => {
        if (!d.items.length) return <div className="flex flex-col gap-2"><p className="text-body text-text-secondary">{t.d6.vacio}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '13' })} texto={t.d6.accionVacio} /></div>
        return <ul className="flex flex-col gap-3">{(d.items as Esquemas['SavingsGoal'][]).map((g) => { const est = g.status ? t.t6.estados[g.status] : undefined
          return <li key={g.id} className="flex flex-col gap-1 text-small tabular-nums">
            <span className="flex flex-wrap items-center justify-between gap-2"><span className="font-medium">{g.name}</span>{est && <Insignia rol={g.status === 'achieved' ? 'exito' : g.status === 'overdue' ? 'advertencia' : 'info'} icono={g.status === 'achieved' ? <IconoOk /> : g.status === 'overdue' ? <IconoAviso /> : <IconoInfo />}>{est}</Insignia>}</span>
            {g.target_cents !== undefined && <progress className={progreso + ' [&::-webkit-progress-value]:bg-role-savings! [&::-moz-progress-bar]:bg-role-savings!'} max={g.target_cents} value={Math.min(g.saved_cents ?? 0, g.target_cents)} aria-label={g.name} />}
            <span className="text-text-secondary">{t.d6.avance(g.progress_tenths_pct !== undefined ? decimas(g.progress_tenths_pct) : t.comun.sinDato)} · {t.d6.deDe(dinero(g.saved_cents ?? 0), g.target_cents !== undefined ? dinero(g.target_cents) : t.comun.sinDato)}</span>
            <span className="text-text-secondary">{g.required_monthly_cents !== undefined ? t.d6.aportacion(dinero(g.required_monthly_cents)) : null}{g.target_date ? ` · ${t.d6.fechaMeta(fecha(g.target_date))}` : ''}</span>
          </li> })}</ul>
      }}</Estados><EnlaceAccion a={rutaDe('P-09')} texto={t.acciones.verTodo} /></Panel>
  }
const ORDEN_T = ['walmart', 'bodega-aurrera', 'sams-club', 'market']
const orden = (s?: { name: string; kind: string }) => ORDEN_T.indexOf(codigoTienda(s))
function Producto({ id }: { id: number }) {
  const q = usarConsulta('get', '/api/v1/products/{product_id}', { ruta: { product_id: id } })
  return <span className="font-medium">{q.datos?.name ?? t.comun.sinDato}</span>
}

## FIN: frontend/src/pantallas/inicio/Detalle.tsx

## ARCHIVO: frontend/src/pantallas/inicio/GraficaInicio.tsx


// Gráficas de Tendencias (U3 §4, DESIGN §7.7 y §8): <figure> con <figcaption> (título + resumen), «Ver como tabla» que cambia
// la gráfica por la tabla (la tabla queda siempre para lectores de pantalla), ECharts 6 en SVG con colores leídos de los tokens.
import { useEffect, useRef, useState, type ReactNode } from 'react'
import * as echarts from 'echarts/core'
import { BarChart, LineChart } from 'echarts/charts'
import { AriaComponent, GridComponent, LegendComponent, MarkAreaComponent, MarkLineComponent, TooltipComponent } from 'echarts/components'
import { SVGRenderer } from 'echarts/renderers'
import { Boton } from '../../componentes/Boton'
import { cn } from '../../componentes/utils'
import { Vacio } from '../../componentes/Estados'
import { textos } from '../../textos/es-MX'
import { textosInicio as t } from './textos'
echarts.use([BarChart, LineChart, AriaComponent, GridComponent, LegendComponent, MarkAreaComponent, MarkLineComponent, TooltipComponent, SVGRenderer])

export const token = (v: string) => getComputedStyle(document.documentElement).getPropertyValue(v).trim()
/** «$12.3 mil»: formato compacto, solo en ejes (U3 §4, DESIGN §8.5). Vive en api/formatos.ts (`dineroCompacto`) para que
 *  Grafica y Tendencias usen el mismo. */
export { dineroCompacto as compacto } from '../../api/formatos'
/** Índices con etiqueta (DESIGN §8.5): el último valor y los extremos. Es el mismo de componentes/Grafica. */
export { selectivos } from '../../componentes/Grafica'
function usarTema() {
  const [n, setN] = useState(0)
  useEffect(() => {
    const mq = matchMedia('(prefers-color-scheme: dark)'), f = () => setN((x) => x + 1)
    const obs = new MutationObserver(f)
    obs.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme', 'class'] })
    mq.addEventListener('change', f)
    return () => { obs.disconnect(); mq.removeEventListener('change', f) }
  }, [])
  return n
}
type SerieOpcion = { type?: string; data?: unknown[]; label?: Record<string, unknown>; endLabel?: unknown; labelLayout?: unknown }
const seriesDe = (op: echarts.EChartsCoreOption) => ([] as SerieOpcion[]).concat((op.series ?? []) as SerieOpcion[])
const valorDe = (d: unknown): number | null => (typeof d === 'number' ? d : d && typeof d === 'object' && 'value' in d ? valorDe((d as { value: unknown }).value) : null)
/** Todas las series en 0 o sin valor: estado vacío en lugar de ejes falsos con «$0» (DESIGN §7.7). */
export const sinDatos = (op: echarts.EChartsCoreOption) => seriesDe(op).every((s) => (s.data ?? []).every((d) => !valorDe(d)))
/** Ancho mínimo de cada columna para su etiqueta directa («$1,197.64» a 12 px); si no cabe, el valor queda en la
 *  descripción emergente y en «Ver como tabla» (DESIGN §8.5). */
const ANCHO_ETIQUETA_COLUMNA = 72
/** Etiquetas sin encimarse (DESIGN §7.7): se recorren en vertical o se ocultan si chocan; con leyenda, los nombres al final
 *  de la línea pasan a la leyenda; las de columnas angostas se quitan. En las líneas, la etiqueta del primer punto se alinea
 *  a la izquierda y la del último a la derecha: así no invaden las etiquetas del eje Y ni salen por el borde (§8.5). */
function sinEncimar(op: echarts.EChartsCoreOption, ancho: number): echarts.EChartsCoreOption {
  const series = seriesDe(op)
  const eje = (Array.isArray(op.xAxis) ? op.xAxis[0] : op.xAxis) as { data?: unknown[] } | undefined
  const columnas = Math.max(1, series.filter((s) => s.type === 'bar').length)
  const caben = ancho / Math.max(1, eje?.data?.length ?? 1) / columnas >= ANCHO_ETIQUETA_COLUMNA
  const sinEtiqueta = (d: unknown) => (d && typeof d === 'object' && 'label' in d ? { ...(d as object), label: { ...(d as { label: object }).label, show: false } } : d)
  const conLabel = (d: unknown): d is { label: object } => !!d && typeof d === 'object' && 'label' in d
  const alineada = (d: unknown, i: number, n: number) => (n > 1 && (i === 0 || i === n - 1) && conLabel(d) ? { ...d, label: { ...d.label, align: i === 0 ? 'left' : 'right' } } : d)
  return {
    ...op,
    series: series.map((s) => ({
      ...s,
      labelLayout: s.labelLayout ?? { hideOverlap: true, moveOverlap: 'shiftY' },
      ...(s.endLabel && op.legend ? { endLabel: { show: false } } : {}),
      ...(s.type === 'bar' && !caben ? { label: { ...s.label, show: false }, data: (s.data ?? []).map(sinEtiqueta) } : {}),
      ...(s.type === 'line' ? { data: (s.data ?? []).map((d, i, todos) => alineada(d, i, todos.length)) } : {}),
    })),
  }
}
export interface Tabla { columnas: string[]; filas: string[][]; total?: string[] }
/** Envoltorio de la tabla alterna (U5-06; axe `scrollable-region-focusable`). Oculta: solo `sr-only`, sin desplazamiento
 *  (con `overflow-x-auto` quedaba una región desplazable sin foco). Visible: región desplazable con foco de teclado, anillo
 *  de foco del token y nombre = título de la gráfica. La usan FiguraGrafica y FiguraBarras (Tablero). */
export function ContenedorTabla({ comoTabla, titulo, children }: { comoTabla: boolean; titulo: string; children: ReactNode }) {
  if (!comoTabla) return <div className="sr-only">{children}</div>
  return <div role="region" aria-label={titulo} tabIndex={0} className="overflow-x-auto rounded-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus">{children}</div>
}
interface Props { id: string; titulo: string; subtitulo: string; resumen: string; notas?: string[]; filtros?: ReactNode; opcion: () => echarts.EChartsCoreOption; tabla: Tabla; dependencias: unknown[] }
export function FiguraGrafica({ id, titulo, subtitulo, resumen, notas, filtros, opcion, tabla, dependencias }: Props) {
  const [comoTabla, setComoTabla] = useState(false)
  const ref = useRef<HTMLDivElement>(null)
  const tema = usarTema()
  useEffect(() => {
    if (comoTabla || !ref.current) return
    const g = echarts.init(ref.current, null, { renderer: 'svg' })
    const reducido = matchMedia('(prefers-reduced-motion: reduce)').matches
    g.setOption({ animation: !reducido, aria: { enabled: true, label: { enabled: true, description: `${titulo}. ${resumen}` }, decal: { show: false } }, textStyle: { fontFamily: 'Inter, ui-sans-serif, system-ui', color: token('--ui-text-secondary') }, ...sinEncimar(opcion(), ref.current.clientWidth) })
    const ro = new ResizeObserver(() => g.resize())
    ro.observe(ref.current)
    return () => { ro.disconnect(); g.dispose() }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [comoTabla, tema, ...dependencias])
  const vacia = sinDatos(opcion())
  const tablaHtml = (
    <ContenedorTabla comoTabla={comoTabla} titulo={titulo}>
      <table className="w-full text-small tabular-nums">
        <caption className="sr-only">{titulo}</caption>
        <thead><tr>{tabla.columnas.map((c, i) => <th key={c} scope="col" className={cn('whitespace-nowrap px-2 py-1 font-medium text-text-secondary', i ? 'text-right' : 'text-left')}>{c}</th>)}</tr></thead>
        <tbody>{tabla.filas.map((f) => <tr key={f[0]} className="border-t border-border-divider">{f.map((v, i) => (i ? <td key={i} className="whitespace-nowrap px-2 py-1 text-right">{v}</td> : <th key={i} scope="row" className="whitespace-nowrap px-2 py-1 text-left font-normal">{v}</th>))}</tr>)}</tbody>
        {tabla.total && <tfoot><tr className="border-t border-border-control font-medium">{tabla.total.map((v, i) => (i ? <td key={i} className="px-2 py-1 text-right">{v}</td> : <th key={i} scope="row" className="px-2 py-1 text-left">{v}</th>))}</tr></tfoot>}
      </table>
    </ContenedorTabla>
  )
  return (
    <figure data-prueba={`grafica-${id}`} className="flex min-w-0 flex-col gap-3 rounded-md border border-border-divider bg-surface-card p-4 tablet:p-6">
      <figcaption className="flex flex-col gap-1">
        <span className="text-small font-medium text-text-secondary">{titulo}</span>
        <span className="text-caption text-text-secondary">{subtitulo}</span>
        <span className="text-small text-text-default tabular-nums">{resumen}</span>
      </figcaption>
      {filtros}
      {vacia && <Vacio texto={textos.vacio.grafica} />}
      {!vacia && !comoTabla && <div ref={ref} aria-hidden className="h-60 w-full tabular-nums tablet:h-72" />}
      {!vacia && tablaHtml}
      {notas?.map((n) => <p key={n} className="text-caption text-text-secondary">{n}</p>)}
      {!vacia && <Boton variante="enlace" className="w-max" aria-pressed={comoTabla} onClick={() => setComoTabla((x) => !x)}>{comoTabla ? t.acciones.verComoGrafica : t.acciones.verComoTabla}</Boton>}
    </figure>
  )
}
/** Ejes y rejilla comunes (DESIGN §8.1): un solo eje Y, sin retícula vertical, textos en los tokens de texto. */
export function ejes(categorias: string[], formatoY: (v: number) => string, ancho: boolean) {
  const sec = token('--ui-text-secondary')
  return {
    grid: { left: 8, right: ancho ? 120 : 16, top: 32, bottom: ancho ? 8 : 40, containLabel: true },
    legend: { [ancho ? 'right' : 'bottom']: 0, ...(ancho ? { top: 'middle', orient: 'vertical' } : {}), textStyle: { color: sec } },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: categorias, axisLine: { lineStyle: { color: token('--ui-viz-baseline') } }, axisTick: { show: false }, axisLabel: { color: sec, fontSize: 12 }, splitLine: { show: false } },
    yAxis: { type: 'value', splitLine: { lineStyle: { color: token('--ui-viz-grid') } }, axisLabel: { color: sec, fontSize: 12, formatter: formatoY } },
  }
}

## FIN: frontend/src/pantallas/inicio/GraficaInicio.tsx

## ARCHIVO: frontend/src/pantallas/inicio/Karen.tsx


// D-7 Participación de Karen (U3 §3.1, §9): solo lo que comparte; Karen no tiene cuenta (decisiones §A).
import { usarConsulta } from '../../api/consulta'
import { dinero, fecha, fechaConDia } from '../../api/formatos'
import { rutaDe } from '../../app/rutas'
import { MarcaPersona } from '../../componentes/Chip'
import { EnlaceAccion, EsqueletoFilas, Estados, IconoOk, Insignia, Panel, karen, useInicio } from './comun'
import { textosInicio as t } from './textos'
import { CircleDashed } from 'lucide-react'

export default function Karen({ aviso }: { aviso?: string[] }) {
  const c = useInicio()
  const k = karen(c)
  const q = usarConsulta('get', '/api/v1/persons/{person_id}/participation', { ruta: { person_id: k?.id ?? 0 } }, k !== undefined)
  const th = 'py-2 pr-3 text-left text-caption font-medium text-text-secondary', td = 'py-2 pr-3 tabular-nums'
  return (
    <Panel etiqueta={t.d7.titulo} className="desktop:col-span-2" titulo={<><MarcaPersona codigo="karen" />{t.d7.titulo}</>}>
      {aviso?.map((a) => <p key={a} className="text-caption text-text-secondary">{a}</p>)}
      <Estados consulta={q} falla={t.d7.falla} esqueleto={<EsqueletoFilas n={4} />}>{(d) => {
        const sinNada = d.purchases.length === 0 && d.payment_shares.length === 0 && d.accrued_groceries_cents === 0
        if (sinNada) return <p className="text-body text-text-secondary">{t.d7.vacio}</p>
        return (
          <div className="flex flex-col gap-4">
            <p className="text-body tabular-nums">{t.d7.despensa(dinero(d.accrued_groceries_cents))}</p>
            {d.purchases.length > 0 && <div className="overflow-x-auto"><table className="w-full text-small"><caption className="mb-1 text-left text-small font-medium text-text-default">{t.d7.compras}</caption>
              <thead><tr><th scope="col" className={th}>{t.d7.comercio}</th><th scope="col" className={th}>{t.d7.fecha}</th><th scope="col" className={th + ' text-right'}>{t.d7.total}</th><th scope="col" className={th}>{t.d7.pagosDe(0, 0).split(' ')[0]}</th><th scope="col" className={th + ' text-right'}>{t.d7.saldo}</th></tr></thead>
              <tbody>{d.purchases.map((p) => <tr key={p.credit_purchase_id} className="border-t border-border-divider"><td className={td}>{p.description ?? t.comun.sinDato}</td><td className={td}>{fecha(p.purchase_date)}</td><td className={td + ' text-right'}>{dinero(p.total_cents)}</td><td className={td}>{t.d7.pagosDe(p.payments_made, p.payments_total)}</td><td className={td + ' text-right'}>{dinero(p.balance_cents)}</td></tr>)}</tbody>
            </table></div>}
            {d.payment_shares.length > 0 && <div className="flex flex-col gap-2"><h3 className="text-small font-medium">{t.d7.pagos}</h3>
              <ul className="flex flex-col">{d.payment_shares.map((p) => <li key={p.scheduled_payment_id} className="flex flex-wrap items-center justify-between gap-2 border-t border-border-divider py-2 text-small tabular-nums">
                <span>{fechaConDia(p.due_date)}</span><span>{dinero(p.amount_cents)}</span>
                {p.status === 'paid' ? <Insignia rol="exito" icono={<IconoOk />}>{t.pago.pagado}</Insignia> : <Insignia rol="neutro" icono={<CircleDashed aria-hidden />}>{t.pago.pendiente}</Insignia>}
              </li>)}</ul></div>}
            <p className="text-body tabular-nums">{t.d7.debe(dinero(d.live_debt_cents))}</p>
          </div>
        )
      }}</Estados>
      <p className="text-small text-text-secondary">{t.d7.nota}</p>
      <EnlaceAccion a={rutaDe('P-25')} texto={t.d7.accion} />
    </Panel>
  )
}

## FIN: frontend/src/pantallas/inicio/Karen.tsx

## ARCHIVO: frontend/src/pantallas/inicio/P01Inicio.tsx


// P-01 Inicio (U1 §2.1–§2.2, §3.3; U3 2.0): Resumen (= tablero de control, K54) · Detalle · Tendencias · Calendario. El cascarón pinta título, pestañas y «Registrar gasto».
import { useEffect } from 'react'
import { Camera, Banknote } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { conConsulta, navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { Boton } from '../../componentes/Boton'
import { ProveedorInicio, useInicio } from './comun'
import { esDiaDePago } from './fechas'
import { textosInicio as t } from './textos'
import Tablero from './Tablero'
import Detalle from './Detalle'
import Tendencias from './Tendencias'
import Calendario from './Calendario'
const CLAVE = 'fh:inicio:pestana'
const SUBS = ['resumen', 'detalle', 'tendencias', 'calendario']

export default function P01Inicio({ sub, consulta }: PropsPantalla) {
  // La pestaña elegida se recuerda en el dispositivo (U3 §1.1, frecuencia «Local»).
  useEffect(() => {
    try {
      if (!consulta.has('sub')) { const g = localStorage.getItem(CLAVE); if (g && g !== 'resumen' && SUBS.includes(g)) navegar(conConsulta({ sub: g }), { reemplazar: true }) }
      else if (sub) localStorage.setItem(CLAVE, sub)
    } catch { /* sin almacenamiento local: abre en Resumen */ }
  }, [sub, consulta])
  return (
    <ProveedorInicio>
      <div className="flex flex-col gap-4">
        <Acciones />
        {sub === 'detalle' ? <Detalle consulta={consulta} /> : sub === 'tendencias' ? <Tendencias /> : sub === 'calendario' ? <Calendario /> : <Tablero />}
      </div>
    </ProveedorInicio>
  )
}
/** Acciones secundarias de P-01 (U1 §2.2): «Foto de ticket»; «Capturar quincena» solo el día de pago si falta (ingreso aún esperado). */
function Acciones() {
  const c = useInicio()
  const falta = esDiaDePago(c.hoy) && c.flujo.datos?.income_is_expected === true
  return (
    <div className="flex flex-wrap justify-end gap-2">
      <Boton variante="secundario" onClick={() => navegar(rutaDe('P-05'))}><Camera aria-hidden className="size-5" />{t.acciones.fotoTicket}</Boton>
      {falta && <Boton variante="secundario" onClick={() => navegar(`${rutaDe('P-08')}?hoja=quincena`)}><Banknote aria-hidden className="size-5" />{t.acciones.capturarQuincena}</Boton>}
    </div>
  )
}

## FIN: frontend/src/pantallas/inicio/P01Inicio.tsx

## ARCHIVO: frontend/src/pantallas/inicio/Resumen.tsx


// Tarjeta de métrica de P-01 (U3 §1.2) y «Más del hogar» (U3 2.0 §2.5, PROVISIONAL, bloqueo U3-tablero-mas-del-hogar): las tarjetas
// 2, 4, 5 y 6 de U3 1.x sin cambio de contenido. La 1 pasó a K-1 y la 3 a K-3 y K-4 (Tablero.tsx). Cada tarjeta es un solo enlace.
import type { ReactNode } from 'react'
import { ArrowDownLeft, ArrowUpRight, Banknote, PiggyBank, ShoppingCart, Utensils } from 'lucide-react'
import { Enlace } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { usarConsulta, type EstadoConsulta } from '../../api/consulta'
import { dinero, fecha, porcentaje } from '../../api/formatos'
import { MarcaPersona, MarcaTienda } from '../../componentes/Chip'
import { cn } from '../../componentes/utils'
import { Aislado, CargadoA, EnlaceAccion, EsqueletoTarjeta, Estados, IconoAviso, IconoInfo, IconoOk, Insignia, PieEnlace, codigoTienda, entero, listaDe, titulares, useInicio } from './comun'
import { diaDe, mesAbrev, rangoSemana, semanaISO, semanaVigente } from './fechas'
import { textosInicio as t } from './textos'

export interface Vista { valor?: string; cuerpo?: ReactNode; rolValor?: 'advertencia' | 'ahorro' | 'ingreso' | 'deuda'; iconoValor?: ReactNode; etiquetaValor?: string; subtitulo?: string; secundario?: ReactNode; extra?: ReactNode; insignias?: ReactNode; estado?: string; accion: { a: string; texto: string } }
const COLOR_VALOR = { advertencia: 'text-role-warning', ahorro: 'text-role-savings', ingreso: 'text-role-income', deuda: 'text-role-debt' }
interface PropsTarjeta<T> { n: number | string; className?: string; titulo: string; icono: ReactNode; consulta: EstadoConsulta<T>; falla: string; accionBase: { a: string; texto: string }; vista: (d: T) => Vista }
export const CAJA = 'flex h-full min-h-11 flex-col gap-2 rounded-md border border-border-divider bg-surface-card p-4 shadow-0 tablet:p-6'
export const Titulo = ({ icono, titulo }: { icono: ReactNode; titulo: string }) => <span className="flex items-center gap-2 text-small font-medium text-text-secondary [&_svg]:size-5">{icono}{titulo}</span>

export function Tarjeta<T>({ n, className, titulo, icono, consulta, falla, accionBase, vista }: PropsTarjeta<T>) {
  const listo = !consulta.cargando || consulta.datos !== undefined
  if (listo && !consulta.error && !consulta.nuncaCargada && consulta.datos !== undefined) {
    return <Aislado texto={falla} alReintentar={consulta.recargar}><Cargada n={n} className={className} titulo={titulo} icono={icono} desde={consulta.desdeCopia} v={() => vista(consulta.datos as T)} /></Aislado>
  }
  return (
    <section data-prueba={`tarjeta-${n}`} aria-label={titulo} className={cn(CAJA, className)}>
      <Titulo icono={icono} titulo={titulo} />
      <Estados consulta={consulta} falla={falla} esqueleto={<EsqueletoTarjeta />}>{() => null}</Estados>
      <EnlaceAccion a={accionBase.a} texto={accionBase.texto} />
    </section>
  )
}
function Cargada({ n, className, titulo, icono, desde, v: armar }: { n: number | string; className?: string; titulo: string; icono: ReactNode; desde: number | null; v: () => Vista }) {
  const v = armar()
  const nombre = [titulo, v.valor, v.etiquetaValor, v.estado, v.accion.texto].filter(Boolean).join(' · ')
  return (
    <Enlace a={v.accion.a} data-prueba={`tarjeta-${n}`} aria-label={nombre} className={cn(CAJA, className, 'text-text-default no-underline hover:border-border-control-hover hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus')}>
      <Titulo icono={icono} titulo={titulo} />
      {v.subtitulo && <span className="text-caption text-text-secondary">{v.subtitulo}</span>}
      {v.valor !== undefined && <span className="flex flex-wrap items-baseline gap-x-2">
        <span className={cn('flex items-center gap-1 whitespace-nowrap text-title-2 font-semibold tabular-nums [&_svg]:size-5', v.rolValor ? COLOR_VALOR[v.rolValor] : 'text-text-default')}>{v.iconoValor}{v.valor}</span>
        {v.etiquetaValor && <span className="text-small text-text-secondary">{v.etiquetaValor}</span>}
      </span>}
      {v.cuerpo}
      {v.secundario && <span className="flex flex-col gap-1 text-small text-text-secondary tabular-nums">{v.secundario}</span>}
      {v.extra && <span className="hidden flex-col gap-1 text-small text-text-secondary tabular-nums phone:flex">{v.extra}</span>}
      {v.insignias && <span className="flex flex-wrap gap-2">{v.insignias}</span>}
      <CargadoA desde={desde} />
      <PieEnlace texto={v.accion.texto} />
    </Enlace>
  )
}

/** «Más del hogar» (U3 2.0 §2.5): título de sección y las tarjetas 2, 4, 5 y 6; 1 columna, 2 × 2 desde 768 y 4 en 1440. */
export function MasDelHogar({ className }: { className?: string }) {
  const c = useInicio()
  const hoy = c.hoy
  const comidas = usarConsulta('get', '/api/v1/dashboard/week-meals', { consulta: { period: semanaISO(hoy) } })
  const canasta = usarConsulta('get', '/api/v1/dashboard/best-basket', { consulta: { week: semanaVigente(hoy) } })
  const ahorro = usarConsulta('get', '/api/v1/dashboard/savings')
  const idGasto = canasta.datos?.basket === 'expense' ? canasta.datos.expense_id : undefined
  const gasto = usarConsulta('get', '/api/v1/expenses/{expense_id}', { ruta: { expense_id: idGasto ?? 0 } }, idGasto !== undefined)
  const tiendas = listaDe(c.tiendas)
  const personas = titulares(c)
  return (
    <section data-prueba="mas-del-hogar" aria-labelledby="mas-del-hogar" className={cn('flex flex-col gap-4', className)}>
      <h2 id="mas-del-hogar" className="text-title-3 font-semibold text-text-default">{t.tablero.masDelHogar}</h2>
      <div className="grid grid-cols-1 gap-4 tablet:grid-cols-2 desktop:gap-6 wide:grid-cols-4">
      <Tarjeta n={2} titulo={t.t2.titulo} icono={<Banknote aria-hidden />} consulta={c.flujo} falla={t.t2.falla} accionBase={{ a: rutaDe('P-07'), texto: t.t2.accion }} vista={(d) => {
        // PROVISIONAL (bloqueo WF3-inicio-vacios): «sin quincena capturada ni esperada» = income_cents en 0.
        const sinIngreso = d.income_cents === 0
        const negativo = d.closing_balance_cents < 0
        const entra = <span className="inline-flex items-center gap-1"><ArrowDownLeft aria-hidden className="size-5 text-role-income" /><span className={sinIngreso ? '' : 'text-role-income'}>{sinIngreso ? t.comun.sinDato : dinero(d.income_cents, 'ingreso')}</span></span>
        const sale = (m: number) => <span className="inline-flex items-center gap-1"><ArrowUpRight aria-hidden className="size-5 text-role-expense" /><span className="text-role-expense">{dinero(m, 'gasto')}</span></span>
        return {
          subtitulo: t.t2.delAl(diaDe(d.start_date), diaDe(d.end_date), mesAbrev(d.end_date)),
          valor: negativo ? `${dinero(d.closing_balance_cents)} ${t.t2.faltante}` : dinero(d.closing_balance_cents), etiquetaValor: negativo ? undefined : t.t2.alCierre,
          rolValor: negativo ? 'advertencia' : undefined, iconoValor: negativo ? <IconoAviso /> : undefined,
          secundario: <><span className="flex flex-wrap items-center gap-1">{t.t2.entra} {entra} · {t.t2.sale} {sale(d.spent_cents + d.committed_cents)}</span>{sinIngreso && <span>{t.t2.faltaIngreso}</span>}</>,
          extra: <span className="flex flex-wrap items-center gap-1">{t.t2.yaSalio} {sale(d.spent_cents)} · {t.t2.porPagar} {sale(d.committed_cents)}</span>,
          estado: d.income_is_expected ? t.t2.esperado : undefined,
          insignias: d.income_is_expected && !sinIngreso ? <Insignia rol="info" icono={<IconoInfo />}>{t.t2.esperado}</Insignia> : undefined,
          accion: sinIngreso ? { a: rutaDe('P-26', { seccion: '03' }), texto: t.t2.accionVacio } : { a: rutaDe('P-07'), texto: t.t2.accion },
        }
      }} />
      <Tarjeta n={4} titulo={t.t4.titulo} icono={<Utensils aria-hidden />} consulta={comidas} falla={t.t4.falla} accionBase={{ a: rutaDe('P-10'), texto: t.t4.accion }} vista={(d) => {
        const vacia = d.unplanned_slots >= 21 && d.cost_cents === 0
        return {
          subtitulo: rangoSemana(d.week_start, d.week_end),
          valor: vacia ? t.comun.sinDato : dinero(d.cost_cents), etiquetaValor: vacia ? undefined : t.t4.costo,
          secundario: vacia ? <span>{t.t4.sinPlan}</span> : d.unplanned_slots === 0 ? <span className="inline-flex items-center gap-1"><IconoOk />{t.t4.completa}</span> : <span>{t.t4.sinArmar(d.unplanned_slots)}</span>,
          extra: vacia ? undefined : <>{personas.map((p) => { const x = d.persons.find((q) => q.person_id === p.id); if (!x) return null
            return <span key={p.id} className="inline-flex items-center gap-2"><MarcaPersona codigo={p.code} />{!x.days_with_plan || x.avg_kcal_per_day === undefined ? t.t4.personaSinDias(p.display_name) : t.t4.persona(p.display_name, entero(x.avg_kcal_per_day), x.daily_kcal_target !== undefined ? entero(x.daily_kcal_target) : null, x.days_with_plan)}</span> })}</>,
          accion: { a: rutaDe('P-10'), texto: vacia ? t.t4.accionVacio : t.t4.accion },
        }
      }} />
      <Tarjeta n={5} titulo={t.t5.titulo} icono={<ShoppingCart aria-hidden />} consulta={canasta} falla={t.t5.falla} accionBase={{ a: rutaDe('P-14'), texto: t.t5.accion }} vista={(d) => {
        if (d.store_id === undefined || d.total_cents === undefined) return { valor: t.comun.sinDato, secundario: <span>{t.t5.sinPrecios}</span>, accion: { a: rutaDe('P-15'), texto: t.t5.accionVacio } }
        const tienda = tiendas.find((s) => s.id === d.store_id), segunda = tiendas.find((s) => s.id === d.second_store_id)
        const dif = d.difference_vs_second_cents
        return {
          valor: `${tienda?.name ?? t.comun.sinDato} ${dinero(d.total_cents)}`, iconoValor: <MarcaTienda codigo={codigoTienda(tienda)} />,
          secundario: <>{dif !== undefined && <span>{segunda ? t.t5.menosQue(dinero(dif), segunda.name) : t.t5.menosQueSiguiente(dinero(dif))}</span>}
            {d.basket === 'expense' && <span>{t.t5.listaVacia(gasto.datos?.expense_date ? fecha(gasto.datos.expense_date) : t.comun.sinDato)}</span>}</>,
          extra: <span>{t.t5.cobertura(d.covered_count, d.item_count)}</span>,
          insignias: (d.stale_price_count > 0 || d.assumes_membership) ? <>{d.stale_price_count > 0 && <span className="hidden phone:inline-flex"><Insignia rol="advertencia" icono={<IconoAviso />}>{t.t5.viejos(d.stale_price_count)}</Insignia></span>}{d.assumes_membership && <Insignia rol="info" icono={<IconoInfo />}>{t.t5.membresia}</Insignia>}</> : undefined,
          accion: { a: rutaDe('P-14'), texto: t.t5.accion },
        }
      }} />
      <Tarjeta n={6} titulo={t.t6.titulo} icono={<PiggyBank aria-hidden />} consulta={ahorro} falla={t.t6.falla} accionBase={{ a: rutaDe('P-09'), texto: t.t6.accion }} vista={(d) => {
        // PROVISIONAL (bloqueo WF3-inicio-vacios): «sin fondo ni metas» = sin meta, sin ahorro y sin plan del mes.
        if (!d.next_goal && d.saved_total_cents === 0 && d.planned_month_cents === 0) return { valor: t.comun.sinDato, secundario: <span>{t.t6.sinAhorro}</span>, accion: { a: rutaDe('P-26', { seccion: '13' }), texto: t.t6.accionVacio } }
        const g = d.next_goal
        const est = g?.status ? t.t6.estados[g.status] : undefined
        return {
          valor: dinero(d.saved_total_cents), rolValor: 'ahorro', iconoValor: <PiggyBank aria-hidden />, etiquetaValor: t.t6.ahorrado,
          secundario: <span>{t.t6.esteMes(dinero(d.contributed_month_cents), dinero(d.planned_month_cents))}</span>,
          extra: g ? <span>{t.t6.meta(g.name, g.progress_pct !== undefined ? porcentaje(g.progress_pct) : t.comun.sinDato, g.target_date ? fecha(g.target_date) : t.comun.sinDato)}</span> : undefined,
          estado: est,
          insignias: g?.status ? <span className="hidden phone:inline-flex"><Insignia rol={g.status === 'achieved' ? 'exito' : g.status === 'overdue' ? 'advertencia' : 'info'} icono={g.status === 'achieved' ? <IconoOk /> : g.status === 'overdue' ? <IconoAviso /> : <IconoInfo />}>{est}</Insignia></span> : undefined,
          accion: { a: rutaDe('P-09'), texto: t.t6.accion },
        }
      }} />
      </div>
    </section>
  )
}

## FIN: frontend/src/pantallas/inicio/Resumen.tsx

## ARCHIVO: frontend/src/pantallas/inicio/Tablero.tsx


// Resumen de P-01 = tablero de control (U3 2.0 §2, K54): K-1 gasto del mes (indicador principal), B-1 barómetro, K-4 pagos
// próximos, K-2 ingreso, K-3 créditos, G-3, G-1, G-6, G-2, G-5 (= T-3), G-4 (= T-2) y «Más del hogar», en ese orden en el DOM
// (orden del foco de 360 en todos los anchos, §2.4). Cada bloque hace su propia petición y tiene sus estados (§1.3, §7); G-1 y
// G-6 comparten una. La interfaz no calcula porcentajes, diferencias ni estados: llegan de C1 v1.5 (motores S2 y S3).
import type { ReactNode } from 'react'
import { useState } from 'react'
import { ArrowDownLeft, ArrowUpRight, Banknote, CalendarClock, CircleCheck, CreditCard, Info, TriangleAlert, Wallet } from 'lucide-react'
import { Enlace } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { usarConsulta, type EstadoConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import type { Respuesta } from '../../api/tipos'
import { dinero, fecha, fechaConDia, mes, porcentaje } from '../../api/formatos'
import { Barometro } from '../../componentes/Barometro'
import { Boton } from '../../componentes/Boton'
import { patron } from '../../componentes/Grafica'
import { cn } from '../../componentes/utils'
import { textos } from '../../textos/es-MX'
import { CAJA, MasDelHogar, Tarjeta, Titulo } from './Resumen'
import { T2, T3 } from './Tendencias'
import { compacto, ContenedorTabla, ejes, FiguraGrafica, token } from './GraficaInicio'
import { EnlaceAccion, EsqueletoFilas, Estados, FallaBloque, IconoAviso, IconoInfo, Insignia, InsigniaPago, TABLET, diasAviso, estadoPago, karen, listaDe, nombreCredito, usarMedia, usarRecargaAlCambiar, useInicio } from './comun'
import { mesDe, sumarDias, sumarMeses } from './fechas'
import { MESES_LARGOS, textosInicio as t } from './textos'

const T = t.tablero, G = t.graficas, B = textos.barometro
const nombreMes = (m: string) => MESES_LARGOS[Number(m.slice(5, 7)) - 1]
const mesLargo = (m: string) => { const n = nombreMes(m); return `${n.charAt(0).toUpperCase()}${n.slice(1)} de ${m.slice(0, 4)}` }
const nada = t.comun.sinDato
const val = (v?: number) => (v === undefined ? nada : dinero(v))
const ESQUELETO = <div aria-hidden className="h-60 w-full rounded-md bg-skeleton-on-card tablet:h-72" />
const CAJA_GRAFICA = 'flex min-w-0 flex-col gap-3 rounded-md border border-border-divider bg-surface-card p-4 tablet:p-6'
const LIGA_FILA = 'flex min-h-11 flex-col justify-center gap-1 rounded-sm py-1 text-text-default no-underline hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus'
type Mensual = Respuesta<'/api/v1/dashboard/monthly-overview', 'get'>

/** Rejilla de §2.4: 1 columna; 768 con 8 columnas (K-1 + B-1, K-4 entera, K-2 + K-3, una gráfica por fila); desde 1024 con 12
 *  (K-1 + B-1; K-2, K-3 y K-4 de 4; gráficas 6 + 6). K-4 va antes que K-2 en el DOM y se coloca a la derecha en 1024. */
const POS = {
  k1: 'tablet:col-span-4 desktop:col-span-6',
  b1: 'tablet:col-span-4 desktop:col-span-6',
  k4: 'tablet:col-span-8 desktop:col-span-4 desktop:col-start-9 desktop:row-start-2',
  k2: 'tablet:col-span-4 desktop:col-span-4 desktop:col-start-1 desktop:row-start-2',
  k3: 'tablet:col-span-4 desktop:col-span-4 desktop:col-start-5 desktop:row-start-2',
  g: 'tablet:col-span-8 desktop:col-span-6',
  mas: 'tablet:col-span-8 desktop:col-span-12',
}
const Bloque = ({ id, className, children }: { id: string; className: string; children: ReactNode }) => <div data-bloque={id} className={cn('min-w-0', className)}>{children}</div>

export default function Tablero() {
  const c = useInicio()
  const actual = mesDe(c.hoy)
  // G-1 y G-6: UNA petición de 17 meses (U3 §2.1; horizonte PROVISIONAL, bloqueo U3-tablero-horizonte).
  const mensual = usarConsulta('get', '/api/v1/dashboard/monthly-overview', { consulta: { from: sumarMeses(actual, -11), to: sumarMeses(actual, 5) } })
  const capacidad = usarConsulta('get', '/api/v1/dashboard/debt-capacity', { consulta: { period: actual } })
  usarRecargaAlCambiar(mensual.recargar)
  return (
    <section data-prueba="tablero" aria-label={T.region} className="grid grid-cols-1 gap-4 tablet:grid-cols-8 desktop:grid-cols-12 desktop:gap-6">
      <Bloque id="K-1" className={POS.k1}><K1 /></Bloque>
      <Bloque id="B-1" className={POS.b1}><B1 q={capacidad} /></Bloque>
      <Bloque id="K-4" className={POS.k4}><K4 /></Bloque>
      <Bloque id="K-2" className={POS.k2}><K2 /></Bloque>
      <Bloque id="K-3" className={POS.k3}><K3 /></Bloque>
      <Bloque id="G-3" className={POS.g}><G3 /></Bloque>
      <Bloque id="G-1" className={POS.g}><G1 q={mensual} /></Bloque>
      <Bloque id="G-6" className={POS.g}><G6 q={mensual} escala={capacidad.datos?.scale_max_pct} /></Bloque>
      <Bloque id="G-2" className={POS.g}><G2 /></Bloque>
      <Bloque id="G-5" className={POS.g}><T3 id="g5" sinFiltros /></Bloque>
      <Bloque id="G-4" className={POS.g}><T2 id="g4" /></Bloque>
      <Bloque id="mas" className={POS.mas}><MasDelHogar /></Bloque>
    </section>
  )
}

/** «[diferencia] más que [mes]», «menos que» o «Igual que», con icono de rol (DESIGN §3.3 y §7.6: la palabra va siempre). */
function Comparacion({ etiqueta, dif, mesAnterior }: { etiqueta?: string; dif: number; mesAnterior: string }) {
  const icono = dif > 0 ? <ArrowUpRight aria-hidden className="size-5 shrink-0 text-role-expense" /> : dif < 0 ? <ArrowDownLeft aria-hidden className="size-5 shrink-0 text-text-secondary" /> : null
  return <span className="inline-flex items-center gap-1">{icono}{etiqueta ? `${etiqueta} · ` : ''}{T.comparar(dinero(Math.abs(dif)), dif, mesAnterior)}</span>
}

// ---- K-1 · Gasto del mes (M-01 + 8 campos de C1 v1.5) ----
function K1() {
  const c = useInicio()
  return <Tarjeta n="K-1" titulo={T.k1.titulo} icono={<Wallet aria-hidden />} consulta={c.plan} falla={T.k1.falla} accionBase={{ a: rutaDe('P-06'), texto: T.k1.accion }} vista={(d) => {
    const conPlan = d.planned_cents !== undefined
    const usado = d.used_pct, resta = d.remaining_cents
    const rebasado = resta !== undefined && resta < 0
    const ant = d.previous_month ? nombreMes(d.previous_month) : undefined
    const alerta = d.status === 'warning' || d.status === 'exceeded'
    const sinPrevio = d.previous_actual_cents === undefined && d.previous_cash_cents === undefined
    const estado = d.status === 'warning' ? T.k1.vasEn(usado ?? 0) : d.status === 'exceeded' ? T.k1.rebasaste : undefined
    return {
      valor: dinero(d.actual_cents), etiquetaValor: T.k1.gastadoEn(nombreMes(d.month)),
      cuerpo: conPlan && usado !== undefined ? (
        <div role="meter" aria-label={T.k1.barra} aria-valuemin={0} aria-valuemax={Math.max(100, usado)} aria-valuenow={usado} aria-valuetext={T.k1.delPlan(porcentaje(usado))} className="h-[var(--size-meter-height)] w-full overflow-hidden rounded-full bg-viz-track">
          <div className={cn('h-full rounded-full', alerta ? 'bg-role-warning' : 'bg-viz-meter')} style={{ width: `${Math.min(100, Math.max(0, usado))}%` }} />
        </div>) : undefined,
      secundario: conPlan ? <>
        <span>{T.k1.dePlan(usado !== undefined ? porcentaje(usado) : nada, dinero(d.planned_cents!), rebasado ? nada : val(resta))}</span>
        {rebasado && <span className="inline-flex items-center gap-1 text-role-warning"><IconoAviso />{T.k1.rebasado(dinero(resta!))}</span>}
      </> : <><span>{`${nada} · ${T.k1.sinPlan(nombreMes(d.month))}`}</span><span>{T.k1.llevas(dinero(d.actual_cents))}</span></>,
      extra: (d.cash_cents !== undefined || ant) ? <>
        {d.cash_cents !== undefined && <span>{T.k1.efectivo(dinero(d.cash_cents))}</span>}
        {ant && sinPrevio && <span>{T.sinDatosDe(ant)}</span>}
        {ant && !sinPrevio && d.actual_difference_cents !== undefined && <Comparacion etiqueta={T.k1.gastado} dif={d.actual_difference_cents} mesAnterior={ant} />}
        {ant && !sinPrevio && d.cash_difference_cents !== undefined && <Comparacion etiqueta={T.k1.enEfectivo} dif={d.cash_difference_cents} mesAnterior={ant} />}
      </> : undefined,
      estado,
      insignias: (estado || (d.categories_exceeded_count ?? 0) > 0) ? <>
        {estado && <Insignia rol="advertencia" icono={<IconoAviso />}>{estado}</Insignia>}
        {(d.categories_exceeded_count ?? 0) > 0 && <Insignia rol="advertencia" icono={<IconoAviso />}>{T.k1.rebasadas(d.categories_exceeded_count!)}</Insignia>}
      </> : undefined,
      accion: conPlan ? { a: rutaDe('P-06'), texto: T.k1.accion } : { a: rutaDe('P-26', { seccion: '02' }), texto: T.k1.accionVacio },
    }
  }} />
}

// ---- K-2 · Ingreso del mes (N-100, M-18) ----
function K2() {
  const c = useInicio()
  const q = usarConsulta('get', '/api/v1/incomes/month', { consulta: { period: mesDe(c.hoy) } })
  return <Tarjeta n="K-2" titulo={T.k2.titulo} icono={<Banknote aria-hidden />} consulta={q} falla={T.k2.falla} accionBase={{ a: rutaDe('P-08'), texto: T.k2.accion }} vista={(d) => {
    const qs = d.fortnights ?? []
    if (!qs.length && d.extras_cents === 0) return { valor: nada, secundario: <span>{T.k2.faltaIngreso}</span>, accion: { a: rutaDe('P-26', { seccion: '03' }), texto: T.k2.accionVacio } }
    const esperado = qs.some((x) => x.status === 'expected')
    const ant = d.previous_month ? nombreMes(d.previous_month) : undefined
    return {
      // Quincenas del mes = capturadas + esperadas sin capturar (U3 K-2: suma de lo que está a la vista, del motor S2).
      valor: dinero(d.captured_cents + d.expected_pending_cents, 'ingreso'), rolValor: 'ingreso', iconoValor: <ArrowDownLeft aria-hidden />, etiquetaValor: T.k2.deQuincenas,
      secundario: <span>{d.extras_cents > 0 ? T.k2.extra(dinero(d.extras_cents, 'ingreso')) : T.k2.sinExtra}</span>,
      extra: <>
        <span>{T.k2.capturadas(qs.filter((x) => x.status === 'captured').length, qs.length)}</span>
        {ant && (d.previous_fortnights_cents === undefined ? <span>{T.sinDatosDe(ant)}</span> : d.fortnights_difference_cents !== undefined && <Comparacion dif={d.fortnights_difference_cents} mesAnterior={ant} />)}
      </>,
      estado: esperado ? T.k2.esperado : undefined,
      insignias: esperado ? <Insignia rol="info" icono={<IconoInfo />}>{T.k2.esperado}</Insignia> : undefined,
      accion: { a: rutaDe('P-08'), texto: T.k2.accion },
    }
  }} />
}

// ---- K-3 · Créditos (M-12 a hoy, M-19 y disponible M-04 con la lectura activa) ----
function K3() {
  const c = useInicio()
  const q = usarConsulta('get', '/api/v1/credits/totals')
  return <Tarjeta n="K-3" titulo={T.k3.titulo} icono={<CreditCard aria-hidden />} consulta={q} falla={T.k3.falla} accionBase={{ a: rutaDe('P-17'), texto: T.k3.accion }} vista={(d) => {
    if (c.creditos.datos && !listaDe(c.creditos).length) return { valor: T.k3.sinCreditos, accion: { a: rutaDe('P-26', { seccion: '05' }), texto: T.k3.accionVacio } }
    const m = nombreMes(mesDe(c.hoy))
    const sinPagos = d.month_total_cents === 0 && !d.month_payment_count
    const av = c.proximo.datos?.availability ?? [], linea = av[0]
    const disponible = c.proximo.error ? <span className="border-t border-border-divider pt-2">{T.k3.disponibleSinDato}</span> : linea && (
      <span className="flex flex-col gap-1 border-t border-border-divider pt-2">
        <span className="tabular-nums text-text-default">{t.t3.disponibleHoy(dinero(linea.available_cents), linea.availability_reading)}</span>
        {!linea.availability_confirmed && <Insignia rol="info" icono={<IconoInfo />} etiqueta={t.t3.reposicion(nombreCredito(c, linea.credit_id))}>{t.t3.pendienteConfirmar}</Insignia>}
        {av.length > 1 && <span>{t.t3.masLineas(av.length - 1)}</span>}
      </span>
    )
    return {
      valor: dinero(d.live_debt_cents), rolValor: 'deuda', iconoValor: <CreditCard aria-hidden />, etiquetaValor: T.k3.deudaViva,
      secundario: <>{sinPagos ? <span>{T.k3.sinPagos(m)}</span> : <span>{T.k3.pagosDe(m, dinero(d.month_total_cents), d.month_payment_count)}</span>}{disponible}</>,
      extra: sinPagos ? undefined : <>
        <span>{t.comun.hogarYKaren(dinero(d.month_household_cents), d.month_external_cents > 0 ? dinero(d.month_external_cents) : undefined)}</span>
        {d.month_pending_household_cents !== undefined && <span>{T.k3.faltaPagar(dinero(d.month_pending_household_cents))}</span>}
      </>,
      accion: { a: rutaDe('P-17'), texto: T.k3.accion },
    }
  }} />
}

// ---- K-4 · Pagos próximos (M-13): vencidos (si los hay) y de hoy a hoy + 30, N de getDashboardSettings (PROVISIONAL) ----
function K4() {
  const c = useInicio()
  const n = c.ajustesTablero.datos?.parameters?.['dashboard.payment_reminder_count']?.value
  const N = typeof n === 'number' && n > 0 ? n : undefined
  const np = c.proximo.datos?.next_payment
  const conVencidos = N !== undefined && np?.is_overdue === true
  const vencidos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { status: 'overdue', sort: 'due_date', limit: N } }, conVencidos)
  const proximos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { status: 'pending', sort: 'due_date', limit: N, period: `${c.hoy}..${sumarDias(c.hoy, 30)}` } }, N !== undefined)
  const idK = karen(c)?.id, aviso = diasAviso(c)
  const sinN = c.ajustesTablero.error || c.ajustesTablero.nuncaCargada || (!c.ajustesTablero.cargando && N === undefined)
  // «Hogar · Karen» y la insignia («Falta la parte de Karen», «Falta la parte del hogar») dependen de quién es Karen (listPersons).
  // Mientras esa lectura no llega no se pinta ninguna fila (esqueleto): sin ella todo caería en el hogar («Hogar $973.75»), una
  // cifra distinta del oráculo (M-2 de REVISION-frontend-cierre-r1). Si falla, el estado de error del bloque con «Reintentar»,
  // igual que cuando falla la lectura de N (U3 §1.2); esta conducta es PROVISIONAL (bloqueo docs/estado/bloqueos/U3-tablero-k4-personas.md).
  const sinPersonas = c.personas.datos === undefined && (Boolean(c.personas.error) || c.personas.nuncaCargada)
  const reintentar = () => { if (sinN) c.ajustesTablero.recargar(); if (sinPersonas) c.personas.recargar() }
  return (
    <section data-prueba="tarjeta-K-4" aria-label={T.k4.titulo} className={CAJA}>
      <Titulo icono={<CalendarClock aria-hidden />} titulo={T.k4.titulo} />
      {sinN || sinPersonas ? <FallaBloque texto={T.k4.falla} alReintentar={reintentar} />
        : N === undefined || c.personas.datos === undefined ? <EsqueletoFilas n={3} />
        : <Estados consulta={proximos} falla={T.k4.falla} esqueleto={<EsqueletoFilas n={3} />}>{(d) => {
          const filas = [...(conVencidos ? vencidos.datos?.items ?? [] : []), ...d.items].slice(0, N)
          if (!filas.length) {
            if (c.creditos.datos && !listaDe(c.creditos).length) return <><p className="text-small text-text-secondary">{T.k4.sinCreditos}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '05' })} texto={T.k4.accionVacio} /></>
            return <p className="text-small text-text-secondary">{T.k4.sinPagos}</p>
          }
          return <>
            <ul className="flex flex-col divide-y divide-border-divider">
              {filas.map((p) => {
                const partes = p.shares ?? [], suma = (xs: typeof partes) => xs.reduce((a, s) => a + s.amount_cents, 0)
                const k = suma(partes.filter((s) => s.person_id === idK)), h = suma(partes.filter((s) => s.person_id !== idK))
                return <li key={p.id}><Enlace a={p.credit_id !== undefined ? rutaDe('P-18', { id: p.credit_id }) : rutaDe('P-17')} className={LIGA_FILA}>
                  <span className="flex flex-wrap items-center gap-x-1 gap-y-1 text-small tabular-nums">{T.k4.fila(fechaConDia(p.due_date), nombreCredito(c, p.credit_id), dinero(p.amount_cents))}{' · '}<InsigniaPago e={estadoPago(p, c.hoy, aviso, idK)} /></span>
                  {partes.length > 0 && <span className="hidden text-small text-text-secondary tabular-nums phone:block">{t.comun.hogarYKaren(dinero(h), k > 0 ? dinero(k) : undefined)}</span>}
                </Enlace></li>
              })}
            </ul>
            {d.next_cursor && <p className="text-small text-text-secondary">{T.k4.hayMas}</p>}
          </>
        }}</Estados>}
      <EnlaceAccion a={`${rutaDe('P-01')}?sub=calendario`} texto={T.k4.accion} />
    </section>
  )
}

// ---- B-1 · Barómetro de capacidad de endeudamiento (M-20, getDashboardDebtCapacity) ----
function B1({ q }: { q: EstadoConsulta<Respuesta<'/api/v1/dashboard/debt-capacity', 'get'>> }) {
  if (esServidorViejo(q.error)) return (
    <section data-prueba="tarjeta-B-1" aria-label={T.b1.titulo} className={CAJA}>
      <Titulo icono={<CreditCard aria-hidden />} titulo={T.b1.titulo} />
      <ServidorViejo />
      <EnlaceAccion a={rutaDe('P-17')} texto={T.b1.accion} />
    </section>
  )
  return <Tarjeta n="B-1" titulo={T.b1.titulo} icono={<CreditCard aria-hidden />} consulta={q} falla={T.b1.falla} accionBase={{ a: rutaDe('P-17'), texto: T.b1.accion }} vista={(d) => {
    const sinIngreso = d.income_cents === undefined
    const conValor = d.ratio_pct !== undefined && d.status !== undefined
    const texto = conValor ? T.b1.valorTexto(`${d.ratio_pct} %`, d.limit_pct, B.colores[d.status!]) : T.b1.sinValorTexto(d.limit_pct)
    return {
      cuerpo: <Barometro titulo={T.b1.titulo} valor={conValor ? d.ratio_pct : undefined} estado={conValor ? d.status : undefined} aviso={d.warning_pct} tope={d.limit_pct} escala={d.scale_max_pct} textoAccesible={texto} />,
      secundario: sinIngreso ? <span>{T.b1.faltaIngreso}</span> : <>
        <span>{T.b1.pagosDe(dinero(d.payments_cents), dinero(d.income_cents!))}</span>
        {!d.has_credits && <span>{T.b1.sinCreditos}</span>}
      </>,
      extra: sinIngreso ? undefined : <>
        {d.headroom_cents !== undefined && (d.headroom_cents >= 0 ? <span>{T.b1.teCaben(dinero(d.headroom_cents))}</span>
          : <span className="inline-flex items-center gap-1 text-role-warning"><IconoAviso />{T.b1.arriba(dinero(d.headroom_cents))}</span>)}
        {d.ratio_total_pct !== undefined && <span className="hidden tablet:inline">{T.b1.conKaren(d.ratio_total_pct)}</span>}
      </>,
      insignias: d.income_is_expected && !sinIngreso ? <Insignia rol="info" icono={<IconoInfo />}>{T.b1.esperado}</Insignia> : undefined,
      estado: conValor ? `${d.ratio_pct} % · ${B.estados[d.status!]}` : nada,
      accion: sinIngreso ? { a: rutaDe('P-26', { seccion: '03' }), texto: T.b1.accionVacio } : { a: rutaDe('P-17'), texto: T.b1.accion },
    }
  }} />
}

// ---- Piezas de las gráficas del tablero ----
/** Servidor anterior a C1 v1.5 (U3 2.0 §2, reglas del tablero): la operación nueva responde 501 `not_implemented` o 404
 *  «La ruta no existe» (C1 no declara 404 en getDashboardDebtCapacity ni en getDashboardMonthlyOverview, así que un 404 ahí
 *  solo es ruta inexistente). Solo se consulta en los bloques que dependen solo de una operación nueva: B-1, G-1 y G-6. */
const esServidorViejo = (e: unknown) => e instanceof ErrorApi && (e.estado === 501 || e.estado === 404)
/** «Disponible cuando el servidor se actualice» en `font.size.caption` y rol información (par C-073 de contraste-componentes):
 *  sin cifra y sin «Reintentar», porque reintentar no lo arregla. */
const ServidorViejo = () => (
  <p role="status" className="flex w-fit max-w-full items-center gap-1 rounded-sm bg-role-info-surface px-2 py-1 text-caption text-role-info [&_svg]:size-5 [&_svg]:shrink-0">
    <Info aria-hidden />{T.servidorViejo}
  </p>
)
/** Estados de una gráfica con su título visible mientras carga o falla (§7: «ídem T-1»). `soloNueva`: el bloque depende solo
 *  de una operación de C1 v1.5 y, con un servidor anterior, muestra ServidorViejo en lugar de la falla. */
function EstadosGrafica<D>({ titulo, consulta, soloNueva = false, children }: { titulo: string; consulta: EstadoConsulta<D>; soloNueva?: boolean; children: (d: D) => ReactNode }) {
  if (soloNueva && esServidorViejo(consulta.error)) return <section aria-label={titulo} className={CAJA_GRAFICA}><span className="text-small font-medium text-text-secondary">{titulo}</span><ServidorViejo /></section>
  if (consulta.datos !== undefined && !consulta.error && !consulta.nuncaCargada) return <Estados consulta={consulta} falla={G.falla} esqueleto={ESQUELETO}>{children}</Estados>
  return <section aria-label={titulo} className={CAJA_GRAFICA}><span className="text-small font-medium text-text-secondary">{titulo}</span><Estados consulta={consulta} falla={G.falla} esqueleto={ESQUELETO}>{children}</Estados></section>
}
function VacioBloque({ id, titulo, texto, accion }: { id: string; titulo: string; texto: ReactNode; accion?: ReactNode }) {
  return <section data-prueba={`grafica-${id}`} aria-label={titulo} className={CAJA_GRAFICA}><span className="text-small font-medium text-text-secondary">{titulo}</span><p className="flex items-center gap-2 text-body text-text-secondary">{texto}</p>{accion}</section>
}
/** Figura con barras horizontales en HTML (G-2, G-3): mismas reglas que FiguraGrafica (figcaption, resumen, «Ver como tabla»
 *  con las mismas cifras, tabla siempre disponible para lectores de pantalla) y filas que pueden ser enlaces de 44 px. */
function FiguraBarras({ id, titulo, subtitulo, resumen, tabla, notas, children }: { id: string; titulo: string; subtitulo: string; resumen: string; tabla: { columnas: string[]; filas: string[][] }; notas?: ReactNode; children: ReactNode }) {
  const [comoTabla, setComoTabla] = useState(false)
  return (
    <figure data-prueba={`grafica-${id}`} className={CAJA_GRAFICA}>
      <figcaption className="flex flex-col gap-1">
        <span className="text-small font-medium text-text-secondary">{titulo}</span>
        <span className="text-caption text-text-secondary">{subtitulo}</span>
        <span className="text-small text-text-default tabular-nums">{resumen}</span>
      </figcaption>
      {!comoTabla && children}
      <ContenedorTabla comoTabla={comoTabla} titulo={titulo}>
        <table className="w-full text-small tabular-nums">
          <caption className="sr-only">{titulo}</caption>
          <thead><tr>{tabla.columnas.map((c, i) => <th key={c} scope="col" className={cn('whitespace-nowrap px-2 py-1 font-medium text-text-secondary', i ? 'text-right' : 'text-left')}>{c}</th>)}</tr></thead>
          <tbody>{tabla.filas.map((f) => <tr key={f[0]} className="border-t border-border-divider">{f.map((v, i) => (i ? <td key={i} className="whitespace-nowrap px-2 py-1 text-right">{v}</td> : <th key={i} scope="row" className="whitespace-nowrap px-2 py-1 text-left font-normal">{v}</th>))}</tr>)}</tbody>
        </table>
      </ContenedorTabla>
      {notas}
      <Boton variante="enlace" className="w-max" aria-pressed={comoTabla} onClick={() => setComoTabla((x) => !x)}>{comoTabla ? t.acciones.verComoGrafica : t.acciones.verComoTabla}</Boton>
    </figure>
  )
}
const tooltipDinero = { trigger: 'axis', valueFormatter: (v: unknown) => (typeof v === 'number' ? dinero(v) : nada) }

// ---- G-1 · Ingreso, gasto y pagos por mes (M-21): columnas agrupadas, 3 series en orden fijo con patrón ----
function G1({ q }: { q: EstadoConsulta<Mensual> }) {
  const c = useInicio(), ancho = usarMedia(TABLET), actual = mesDe(c.hoy)
  return <EstadosGrafica titulo={T.g1.titulo} consulta={q} soloNueva>{(d) => {
    const m = d.months.filter((x) => !x.is_projection && x.month <= actual).slice(-12)
    if (!m.length || m.every((x) => x.income_total_cents === undefined && x.accrued_cents === undefined))
      return <VacioBloque id="g1" titulo={T.g1.titulo} texto={T.g1.vacio} accion={<EnlaceAccion a={rutaDe('P-26', { seccion: '03' })} texto={T.g1.accionVacio} />} />
    const iu = Math.max(0, m.findIndex((x) => x.month === actual)), u = m[iu] ?? m[m.length - 1]
    const cs = m.map((x) => mes(x.month)), sup = token('--ui-surface-card')
    const serie = (nombre: string, i: number, v: (number | undefined)[]) => ({
      name: nombre, type: 'bar', barGap: '8%', itemStyle: { color: token(`--ui-chart-${i + 1}`), decal: patron(i, sup) }, emphasis: { disabled: true },
      data: v.map((x, j) => (x === undefined ? null : j === iu ? { value: x, label: { show: true, position: 'top', color: token('--ui-text-default'), fontSize: 12, formatter: () => dinero(x) } } : x)),
    })
    return <FiguraGrafica id="g1" titulo={T.g1.titulo} subtitulo={T.g1.subtitulo(mes(m[0].month), mes(m[m.length - 1].month))} dependencias={[d, ancho]}
      resumen={T.g1.resumen(mes(u.month), val(u.income_total_cents), val(u.accrued_cents), dinero(u.credit_payments_cents))} notas={[T.g1.nota]}
      tabla={{ columnas: [G.mes, T.g1.ingreso, T.g1.quincenas, T.g1.extras, T.g1.gasto, T.g1.pagos], filas: m.map((x) => [mes(x.month), val(x.income_total_cents), val(x.income_fortnights_cents), val(x.income_extras_cents), val(x.accrued_cents), dinero(x.credit_payments_cents)]) }}
      opcion={() => ({ ...ejes(cs, compacto, ancho), tooltip: tooltipDinero, series: [serie(T.g1.ingreso, 0, m.map((x) => x.income_total_cents)), serie(T.g1.gasto, 1, m.map((x) => x.accrued_cents)), serie(T.g1.pagos, 2, m.map((x) => x.credit_payments_cents))] })} />
  }}</EstadosGrafica>
}

// ---- G-6 · Tendencia del barómetro (M-21): línea con proyección punteada, referencias en los umbrales y «Hoy» ----
function G6({ q, escala }: { q: EstadoConsulta<Mensual>; escala?: number }) {
  const c = useInicio(), ancho = usarMedia(TABLET), actual = mesDe(c.hoy)
  return <EstadosGrafica titulo={T.g6.titulo} consulta={q} soloNueva>{(d) => {
    const m = d.months.filter((x) => x.month >= sumarMeses(actual, -6) && x.month <= sumarMeses(actual, 5))
    if (!m.length || m.every((x) => x.debt_ratio_pct === undefined))
      return <VacioBloque id="g6" titulo={T.g6.titulo} texto={T.g6.vacio} accion={<EnlaceAccion a={rutaDe('P-26', { seccion: '03' })} texto={T.g6.accionVacio} />} />
    const cs = m.map((x) => mes(x.month)), hoyI = m.findIndex((x) => x.month === actual)
    const color = (e?: 'normal' | 'warning' | 'exceeded') => (e ? B.colores[e] : nada)
    const pct = (v?: number) => (v === undefined ? nada : `${v} %`)
    const real = m.map((x) => (!x.is_projection ? x.debt_ratio_pct ?? null : null))
    const proy = m.map((x, i) => (x.is_projection || i === hoyI ? x.debt_ratio_pct ?? null : null))
    const proyectados = m.filter((x) => x.is_projection && x.debt_ratio_pct !== undefined)
    const mayor = proyectados.reduce<typeof proyectados[number] | undefined>((a, b) => (!a || b.debt_ratio_pct! > a.debt_ratio_pct! ? b : a), undefined)
    const u = m[hoyI] ?? m[m.length - 1]
    // Escala de G-6 = la del barómetro (M-R3-1, decisión 3 del orquestador = opción B: `scale_max_pct` entero requerido en
    // DashboardMonthlyOverview, C1 1.5.1). La de getDashboardDebtCapacity y `limit_pct` quedan solo si un servidor viejo no lo manda.
    const tope = Math.max(d.scale_max_pct ?? escala ?? d.limit_pct, ...m.map((x) => x.debt_ratio_pct ?? 0))
    const c1 = token('--ui-chart-1'), base = token('--ui-viz-baseline'), sec = token('--ui-text-secondary')
    const referencia = (v: number) => ({ yAxis: v, lineStyle: { color: base, width: 1, type: 'solid' }, label: { formatter: `${v} %`, position: 'end', color: sec } })
    return <FiguraGrafica id="g6" titulo={T.g6.titulo} subtitulo={T.g1.subtitulo(cs[0], cs[cs.length - 1])} dependencias={[d, ancho, escala]}
      resumen={mayor ? T.g6.resumen(mes(u.month), pct(u.debt_ratio_pct), color(u.debt_status), mes(mayor.month), pct(mayor.debt_ratio_pct)) : T.g6.resumenSinProyeccion(mes(u.month), pct(u.debt_ratio_pct), color(u.debt_status))}
      notas={[T.g6.nota]}
      tabla={{ columnas: [G.mes, T.g6.capacidad, T.g6.estado, T.g1.pagos, T.g1.ingreso, T.g6.proyeccion], filas: m.map((x) => [mes(x.month), x.debt_ratio_pct === undefined ? T.g6.faltaIngreso : pct(x.debt_ratio_pct), color(x.debt_status), dinero(x.credit_payments_cents), val(x.income_total_cents), x.is_projection ? T.g6.si : T.g6.no]) }}
      opcion={() => {
        const e = ejes(cs, (v: number) => `${v} %`, ancho)
        return { ...e, yAxis: { ...e.yAxis, min: 0, max: tope },
          tooltip: { trigger: 'axis', formatter: (ps: { dataIndex: number }[]) => { const x = m[ps[0]?.dataIndex ?? 0]; return x.debt_ratio_pct === undefined ? `${mesLargo(x.month)}: ${T.g6.faltaIngreso}` : T.g6.punto(mesLargo(x.month), x.debt_ratio_pct, color(x.debt_status)) } },
          series: [
            { name: T.g6.serie, type: 'line', symbol: 'circle', symbolSize: 8, connectNulls: false, emphasis: { disabled: true }, lineStyle: { color: c1, width: 2 }, itemStyle: { color: c1 },
              data: real.map((v, i) => (v !== null && i === hoyI ? { value: v, label: { show: true, position: 'top', color: token('--ui-text-default'), fontSize: 12, formatter: `${v} %` } } : v)),
              markLine: { symbol: 'none', silent: true, data: [referencia(d.warning_pct), referencia(d.limit_pct), ...(hoyI >= 0 ? [{ xAxis: cs[hoyI], lineStyle: { color: base, width: 1, type: 'solid' }, label: { formatter: G.hoy, color: sec } }] : [])] } },
            { name: T.g6.serie, type: 'line', symbol: 'circle', symbolSize: 8, connectNulls: false, emphasis: { disabled: true }, lineStyle: { color: c1, width: 2, type: [8, 4] }, itemStyle: { color: c1 }, data: proy },
          ] }
      }} />
  }}</EstadosGrafica>
}

// ---- G-2 · Deuda viva por línea (M-12): barras horizontales en orden descendente con el valor al final ----
function G2() {
  const c = useInicio()
  const q = usarConsulta('get', '/api/v1/debt/live', { consulta: { date: c.hoy } })
  return <EstadosGrafica titulo={T.g2.titulo} consulta={q}>{(d) => {
    const filas = d.credits.filter((x) => x.total_cents > 0).sort((a, b) => b.total_cents - a.total_cents)
    if (!filas.length || d.total_cents === 0) return <VacioBloque id="g2" titulo={T.g2.titulo} texto={<><CircleCheck aria-hidden className="size-5 text-role-success" />{T.g2.vacio}</>} />
    const max = filas[0].total_cents
    return <FiguraBarras id="g2" titulo={T.g2.titulo} subtitulo={T.g2.subtitulo(fecha(d.date))} resumen={T.g2.resumen(dinero(d.total_cents), nombreCredito(c, filas[0].credit_id), dinero(filas[0].total_cents))}
      tabla={{ columnas: [T.g2.linea, T.g2.monto], filas: filas.map((x) => [nombreCredito(c, x.credit_id), dinero(x.total_cents)]) }}>
      <ul aria-hidden className="flex flex-col gap-3">
        {filas.map((x) => <li key={x.credit_id} className="flex flex-col gap-1">
          <span className="text-small text-text-default">{nombreCredito(c, x.credit_id)}</span>
          <span className="flex items-center gap-2"><span className="block h-[var(--size-viz-bar-max)] shrink-0 rounded-r-[var(--size-viz-bar-radius)] bg-chart-1" style={{ width: `calc((100% - 7rem) * ${x.total_cents / max})` }} /><span className="whitespace-nowrap text-small text-text-default tabular-nums">{dinero(x.total_cents)}</span></span>
        </li>)}
      </ul>
    </FiguraBarras>
  }}</EstadosGrafica>
}

// ---- G-3 · Presupuesto por categoría (M-09) con la sugerencia subir o recortar (M-22, PROVISIONAL) ----
function G3() {
  const c = useInicio(), amplio = usarMedia(TABLET), periodo = mesDe(c.hoy)
  const q = usarConsulta('get', '/api/v1/budget/categories', { consulta: { period: periodo } })
  const senales = usarConsulta('get', '/api/v1/budget/signals', { consulta: { period: periodo } })
  const umbral = (k: 'alerts.budget_warning_pct' | 'alerts.budget_limit_pct') => { const v = c.alertas.datos?.parameters?.[k]?.value; return typeof v === 'number' ? v : undefined }
  const w = umbral('alerts.budget_warning_pct'), l = umbral('alerts.budget_limit_pct')
  return <EstadosGrafica titulo={T.g3.titulo} consulta={q}>{(d) => {
    const cats = listaDe(c.categorias), nombre = (id: number) => cats.find((x) => x.id === id)?.name ?? nada
    const conPlan = d.rows.filter((r) => r.planned_cents !== undefined && r.used_pct !== undefined).sort((a, b) => b.used_pct! - a.used_pct!)
    const sinPlan = d.rows.filter((r) => r.planned_cents === undefined && r.actual_cents > 0).length
    const capturar = <EnlaceAccion a={rutaDe('P-26', { seccion: '02' })} texto={T.g3.accionVacio} />
    if (!conPlan.length) return <VacioBloque id="g3" titulo={T.g3.titulo} texto={T.g3.vacio} accion={capturar} />
    const visibles = amplio ? conPlan : conPlan.slice(0, 8)
    const eje = Math.max(100, ...visibles.map((r) => r.used_pct!))
    const pos = (n: number) => `${(Math.min(n, eje) / eje) * 100}%`
    const s = senales.error ? undefined : senales.datos
    const senal = (id: number) => s?.rows.find((r) => r.category_id === id)?.signal
    const alerta = (r: { status?: string }) => r.status === 'warning' || r.status === 'exceeded'
    const textoSenal = (x?: string) => (x === 'raise' ? T.g3.subir : x === 'cut' ? T.g3.recortar : T.g3.ninguna)
    return <FiguraBarras id="g3" titulo={T.g3.titulo} subtitulo={T.g3.subtitulo(mes(d.month ?? periodo))}
      resumen={T.g3.resumen(nombreMes(periodo), visibles.filter(alerta).length, w ?? 0, visibles.filter((r) => r.status === 'exceeded').length)}
      tabla={{ columnas: [T.g3.categoria, T.g3.real, T.g3.plan, T.g3.usado, T.g3.sugerencia], filas: visibles.map((r) => [nombre(r.category_id), dinero(r.actual_cents), dinero(r.planned_cents!), porcentaje(r.used_pct!), textoSenal(senal(r.category_id))]) }}
      notas={<>
        {!amplio && conPlan.length > visibles.length && <EnlaceAccion a={rutaDe('P-06')} texto={T.g3.verTodas} />}
        {sinPlan > 0 && <p className="flex flex-wrap items-center gap-2 text-small text-text-secondary">{T.g3.sinPlan(sinPlan)}{capturar}</p>}
      </>}>
      <div aria-hidden className="relative h-5 text-caption text-text-secondary tabular-nums">
        {w !== undefined && <span className="absolute -translate-x-1/2" style={{ left: pos(w) }}>{T.g3.marca(w)}</span>}
        {l !== undefined && <span className="absolute -translate-x-1/2" style={{ left: pos(l) }}>{T.g3.marca(l)}</span>}
      </div>
      <ul className="flex flex-col gap-1">
        {visibles.map((r) => {
          const x = senal(r.category_id), n = nombre(r.category_id)
          return <li key={r.category_id}><Enlace a={`${rutaDe('P-02')}?categoria=${r.category_id}&periodo=${periodo}`} aria-label={T.g3.abrir(n)} className={LIGA_FILA}>
            <span className="flex flex-wrap items-center gap-x-2 text-small tabular-nums">{T.g3.etiqueta(n, dinero(r.actual_cents), dinero(r.planned_cents!), porcentaje(r.used_pct!))}
              {alerta(r) && <span className="inline-flex items-center gap-1 text-role-warning"><TriangleAlert aria-hidden className="size-5" />{r.status === 'exceeded' ? T.g3.rebasada : T.g3.vasEn(r.used_pct!)}</span>}
            </span>
            <span className="relative block h-[var(--size-viz-bar-max)] w-full">
              <span className={cn('absolute inset-y-0 left-0 rounded-r-[var(--size-viz-bar-radius)]', alerta(r) ? 'bg-role-warning' : 'bg-chart-1')} style={{ width: pos(r.used_pct!) }} />
              {w !== undefined && <span aria-hidden data-marca="aviso" className="absolute -inset-y-1 border-l border-dashed border-viz-baseline" style={{ left: pos(w) }} />}
              {l !== undefined && <span aria-hidden data-marca="tope" className="absolute -inset-y-1 border-l-2 border-viz-baseline" style={{ left: pos(l) }} />}
            </span>
            {(x === 'raise' || x === 'cut') && s && <Insignia rol="info" icono={<IconoInfo />} etiqueta={T.g3.motivo(n, s.months_considered, x === 'raise', s.margin_pct)}>{textoSenal(x)}</Insignia>}
          </Enlace></li>
        })}
      </ul>
    </FiguraBarras>
  }}</EstadosGrafica>
}

## FIN: frontend/src/pantallas/inicio/Tablero.tsx

## ARCHIVO: frontend/src/pantallas/inicio/Tendencias.tsx


// Tendencias de P-01 (U3 §4): T-1 a T-4 en orden fijo, series en orden fijo con marcador, trazo o patrón y etiqueta directa.
import { useState, type ReactNode } from 'react'
import { CircleCheck } from 'lucide-react'
import { rutaDe } from '../../app/rutas'
import { usarConsulta } from '../../api/consulta'
import { dinero, mes } from '../../api/formatos'
import { ChipSeleccionable, MarcaPersona } from '../../componentes/Chip'
import { patron } from '../../componentes/Grafica'
import { abrirRegistro } from '../../shell/estado'
import { Boton } from '../../componentes/Boton'
import { EnlaceAccion, Estados, TABLET, entero, listaDe, titulares, usarMedia, useInicio } from './comun'
import { compacto, ejes, FiguraGrafica, selectivos, token } from './GraficaInicio'
import { diaDe, mesAbrev, mesDe, rangoSemana, sumarDias } from './fechas'
import { textosInicio as t } from './textos'
const G = t.graficas
const esqueleto = <div aria-hidden className="h-60 w-full rounded-md bg-skeleton-on-card tablet:h-72" />
const etiqueta = (formato: (v: number) => string) => ({ show: true, color: token('--ui-text-default'), fontSize: 12, formatter: (p: { value: number }) => formato(p.value) })
const conEtiquetas = (v: (number | null)[], formato: (v: number) => string, solo?: Set<number>) => { const s = solo ?? selectivos(v); return v.map((x, i) => (x === null ? null : s.has(i) ? { value: x, label: etiqueta(formato) } : x)) }
function VacioGrafica({ titulo, texto, accion }: { titulo: string; texto: ReactNode; accion?: ReactNode }) {
  return <section aria-label={titulo} className="flex flex-col gap-2 rounded-md border border-border-divider bg-surface-card p-4 tablet:p-6"><span className="text-small font-medium text-text-secondary">{titulo}</span><p className="flex items-center gap-2 text-body text-text-secondary">{texto}</p>{accion}</section>
}
function Chips<T extends string | number>({ etiqueta: e, opciones, valor, alCambiar, todas }: { etiqueta: string; opciones: { valor: T; texto: string; marca?: ReactNode }[]; valor: T[]; alCambiar: (v: T[]) => void; todas: string }) {
  return <div role="group" aria-label={e} className="flex flex-wrap gap-2">
    <ChipSeleccionable seleccionado={!valor.length} onClick={() => alCambiar([])}>{todas}</ChipSeleccionable>
    {opciones.map((o) => <ChipSeleccionable key={o.valor} seleccionado={valor.includes(o.valor)} marca={o.marca} onClick={() => alCambiar(valor.includes(o.valor) ? valor.filter((x) => x !== o.valor) : [...valor, o.valor])}>{o.texto}</ChipSeleccionable>)}
  </div>
}

export default function Tendencias() {
  return <div className="grid grid-cols-1 gap-4 wide:grid-cols-2 wide:gap-6"><T1 /><T2 /><T3 /><T4 /></div>
}
function T1() {
  const c = useInicio(), ancho = usarMedia(TABLET)
  const [cat, setCat] = useState<number[]>([])
  const q = usarConsulta('get', '/api/v1/budget/monthly-trend', { consulta: { category_id: cat.length ? cat : undefined } })
  const cats = listaDe(c.categorias)
  const nombre = cat.length ? cats.filter((x) => cat.includes(x.id)).map((x) => x.name).join(', ') : G.todas
  const filtros = <Chips etiqueta={t.detalle.categoria} todas={G.todas} valor={cat} alCambiar={setCat} opciones={cats.map((x) => ({ valor: x.id, texto: x.name }))} />
  return <Estados consulta={q} falla={G.falla} esqueleto={esqueleto}>{(d) => {
    const m = d.months
    if (!m.some((x) => x.planned_cents !== undefined)) return <VacioGrafica titulo={G.t1.titulo} texto={G.t1.vacio} accion={<><EnlaceAccion a={rutaDe('P-26', { seccion: '02' })} texto={G.t1.accionVacio} />{filtros}</>} />
    const u = m[m.length - 1], cs = m.map((x) => mes(x.month)), ult = new Set([m.length - 1])
    const sup = token('--ui-surface-card')
    return <FiguraGrafica id="t1" titulo={G.t1.titulo} subtitulo={G.ultimos12(nombre)} filtros={filtros} dependencias={[d, ancho]}
      resumen={u.planned_cents !== undefined ? G.t1.resumen(mes(u.month), dinero(u.actual_cents), dinero(u.planned_cents)) : G.t1.resumenSinPlan(mes(u.month), dinero(u.actual_cents))}
      tabla={{ columnas: [G.mes, G.t1.real, G.t1.plan], filas: m.map((x) => [mes(x.month), dinero(x.actual_cents), x.planned_cents !== undefined ? dinero(x.planned_cents) : t.comun.sinDato]) }}
      opcion={() => ({ ...ejes(cs, compacto, ancho), labelLayout: { hideOverlap: true }, series: [
        { name: G.t1.real, type: 'bar', barGap: '0%', data: conEtiquetas(m.map((x) => x.actual_cents), (v) => dinero(v), ult), label: { position: 'top' }, itemStyle: { color: token('--ui-chart-1'), borderColor: sup, borderWidth: 1 } },
        { name: G.t1.plan, type: 'bar', data: m.map((x) => x.planned_cents ?? null), itemStyle: { color: token('--ui-chart-2'), decal: patron(1, sup), borderColor: sup, borderWidth: 1 } },
      ] })} />
  }}</Estados>
}
/** T-2; en el tablero es G-4 sin cambio (U3 2.0 §2.1), con la misma petición. */
export function T2({ id = 't2' }: { id?: string }) {
  const c = useInicio(), ancho = usarMedia(TABLET)
  const q = usarConsulta('get', '/api/v1/cash-flow/accrual-vs-cash')
  return <Estados consulta={q} falla={G.falla} esqueleto={esqueleto}>{(d) => {
    const m = d.months
    if (!m.length || m.every((x) => x.accrued_cents === 0 && x.cash_cents === 0)) return <VacioGrafica titulo={G.t2.titulo} texto={G.t2.vacio} accion={<Boton variante="enlace" className="w-max" onClick={() => abrirRegistro('gasto')}>{G.t2.accionVacio}</Boton>} />
    const u = m[m.length - 1], cs = m.map((x) => mes(x.month)), sup = token('--ui-surface-card')
    const hoyI = m.findIndex((x) => x.month === mesDe(c.hoy))
    const adelantado = m.filter((x) => x.cumulative_difference_cents < 0).map((x) => mes(x.month))
    const suma = (k: 'accrued_cents' | 'cash_cents' | 'accrued_household_cents' | 'accrued_external_cents' | 'cash_household_cents' | 'cash_external_cents') => dinero(m.reduce((a, x) => a + x[k], 0))
    const C = G.t2.columnas
    return <FiguraGrafica id={id} titulo={G.t2.titulo} subtitulo={G.t2.subtitulo(mes(m[0].month), mes(u.month))} dependencias={[d, ancho]}
      resumen={G.t2.resumen(mes(u.month), m[6] ? mes(m[6].month) : t.comun.sinDato, dinero(u.cumulative_accrued_cents), dinero(u.cumulative_cash_cents))}
      notas={[d.projection_mode === 'with_repeat_cadence' ? G.t2.conRepetidas : G.t2.sinRepetidas, ...adelantado.map(G.t2.adelantadoEn)]}
      tabla={{ columnas: [G.mes, C.devengadoMes, C.efectivoMes, G.t2.diferencia, G.t2.devengado, G.t2.efectivo, C.devHogar, C.devKaren, C.efeHogar, C.efeKaren],
        filas: m.map((x) => [mes(x.month), dinero(x.accrued_cents), dinero(x.cash_cents), dinero(x.cumulative_difference_cents), dinero(x.cumulative_accrued_cents), dinero(x.cumulative_cash_cents), dinero(x.accrued_household_cents), dinero(x.accrued_external_cents), dinero(x.cash_household_cents), dinero(x.cash_external_cents)]),
        total: [G.total, suma('accrued_cents'), suma('cash_cents'), t.comun.sinDato, t.comun.sinDato, t.comun.sinDato, suma('accrued_household_cents'), suma('accrued_external_cents'), suma('cash_household_cents'), suma('cash_external_cents')] }}
      opcion={() => {
        // Áreas apiladas (DESIGN §8.3 y §8.5): los valores finales van en la leyenda, abajo y fuera del trazo (una etiqueta sobre
        // el relleno se encima con el área y, a la derecha, la leyenda tapaba su marca). La tabla sigue con todos los valores.
        const e = ejes(cs, compacto, ancho)
        const finales: Record<string, string> = { [G.t2.efectivo]: dinero(u.cumulative_cash_cents), [G.t2.diferencia]: dinero(Math.max(0, u.cumulative_difference_cents)) }
        return { ...e, grid: { ...e.grid, right: 16, bottom: ancho ? 40 : 64 }, legend: { bottom: 0, textStyle: e.legend.textStyle, formatter: (n: string) => `${n} ${finales[n] ?? ''}`.trim() }, series: [
        { name: G.t2.efectivo, type: 'line', stack: 'pila', symbol: 'circle', symbolSize: 8, lineStyle: { color: token('--ui-chart-1'), width: 2 }, itemStyle: { color: token('--ui-chart-1') }, areaStyle: { color: token('--ui-chart-1'), opacity: 1 },
          data: m.map((x) => x.cumulative_cash_cents),
          markLine: hoyI >= 0 ? { symbol: 'none', silent: true, lineStyle: { color: token('--ui-viz-baseline'), width: 1, type: 'solid' }, label: { formatter: G.hoy, color: token('--ui-text-secondary') }, data: [{ xAxis: cs[hoyI] }] } : undefined,
          markArea: m.length > 6 ? { silent: true, itemStyle: { color: token('--ui-viz-track'), opacity: 0.35 }, label: { formatter: G.proyeccion, color: token('--ui-text-secondary') }, data: [[{ xAxis: cs[6] }, { xAxis: cs[m.length - 1] }]] } : undefined },
        { name: G.t2.diferencia, type: 'line', stack: 'pila', symbol: 'rect', symbolSize: 8, lineStyle: { color: token('--ui-chart-2'), width: 2, type: [8, 4] }, itemStyle: { color: token('--ui-chart-2') }, areaStyle: { color: token('--ui-chart-2'), opacity: 1, decal: patron(2, sup) },
          data: m.map((x) => Math.max(0, x.cumulative_difference_cents)) },
      ] } }} />
  }}</Estados>
}
/** T-3; en el tablero es G-5 sin filtro de crédito (U3 2.0 §2.1). */
export function T3({ id = 't3', sinFiltros = false }: { id?: string; sinFiltros?: boolean }) {
  const c = useInicio(), ancho = usarMedia(TABLET)
  const [sel, setSel] = useState<number[]>([])
  const q = usarConsulta('get', '/api/v1/credits/debt-trend', { consulta: { credit_id: sel.length ? sel : undefined } })
  const creditos = listaDe(c.creditos), elegidos = creditos.filter((x) => sel.includes(x.id))
  const filtros = sinFiltros ? undefined : <Chips etiqueta={t.detalle.credito} todas={G.sinFiltro} valor={sel} alCambiar={setSel} opciones={creditos.map((x) => ({ valor: x.id, texto: x.name }))} />
  return <Estados consulta={q} falla={G.falla} esqueleto={esqueleto}>{(d) => {
    const m = d.months
    if (!m.length || m.every((x) => x.live_debt_cents === 0)) return <VacioGrafica titulo={G.t3.titulo} texto={<><CircleCheck aria-hidden className="size-5 text-role-success" />{G.t3.vacio}</>} accion={filtros} />
    const u = m[m.length - 1], cs = m.map((x) => mes(x.month))
    const series = elegidos.length ? elegidos.map((cr) => ({ nombre: cr.name, v: m.map((x) => x.credits.find((y) => y.credit_id === cr.id)?.live_debt_cents ?? null) })) : [{ nombre: G.t3.total, v: m.map((x) => x.live_debt_cents) }]
    const MARC = ['circle', 'rect', 'triangle', 'diamond'], TRAZO: (string | number[])[] = ['solid', [8, 4], [2, 4], [8, 4, 2, 4]]
    return <FiguraGrafica id={id} titulo={G.t3.titulo} subtitulo={G.ultimos12(elegidos.length ? elegidos.map((x) => x.name).join(', ') : G.t3.total)} filtros={filtros} dependencias={[d, ancho, sel.join()]}
      resumen={elegidos.length ? G.t3.resumenFiltro(mes(u.month), series.map((s) => `${s.nombre} ${s.v[s.v.length - 1] !== null ? dinero(s.v[s.v.length - 1]!) : t.comun.sinDato}`).join(' · ')) : G.t3.resumen(mes(m[0].month), dinero(m[0].live_debt_cents), mes(u.month), dinero(u.live_debt_cents))}
      tabla={{ columnas: [G.mes, ...series.map((s) => s.nombre), G.t3.interesesPagados, G.t3.interesesPorPagar], filas: m.map((x, i) => [mes(x.month), ...series.map((s) => (s.v[i] !== null ? dinero(s.v[i]!) : t.comun.sinDato)), dinero(x.interest_paid_cents), dinero(x.interest_pending_cents)]) }}
      opcion={() => ({ ...ejes(cs, compacto, ancho), series: series.slice(0, 8).map((s, i) => ({ name: s.nombre, type: 'line', symbol: MARC[i % 4], symbolSize: 8, lineStyle: { color: token(`--ui-chart-${i + 1}`), width: 2, type: TRAZO[i % 4] }, itemStyle: { color: token(`--ui-chart-${i + 1}`) }, label: { position: 'top' }, data: conEtiquetas(s.v, (v) => dinero(v)) })) })} />
  }}</Estados>
}
function T4() {
  const c = useInicio(), ancho = usarMedia(TABLET)
  const [sel, setSel] = useState<number[]>([])
  const q = usarConsulta('get', '/api/v1/meal-plan/trends', { consulta: { person_id: sel.length ? sel : undefined } })
  const pers = titulares(c), vistas = sel.length ? pers.filter((p) => sel.includes(p.id)) : pers
  const filtros = <Chips etiqueta={t.detalle.persona} todas={G.sinFiltro} valor={sel} alCambiar={setSel} opciones={pers.map((p) => ({ valor: p.id, texto: p.display_name, marca: <MarcaPersona codigo={p.code} /> }))} />
  return <Estados consulta={q} falla={G.falla} esqueleto={esqueleto}>{(d) => {
    const w = d.weeks
    if (!w.length) return <VacioGrafica titulo={G.t4.tituloA} texto={G.t4.vacio} accion={<EnlaceAccion a={rutaDe('P-10')} texto={G.t4.accionVacio} />} />
    const u = w[w.length - 1], cs = w.map((x) => `${diaDe(x.week_start)} ${mesAbrev(x.week_start)}`), rango = rangoSemana(u.week_start, sumarDias(u.week_start, 6))
    const kcal = (p: number, k: 'avg_kcal_per_day' | 'daily_kcal_target') => w.map((x) => x.persons.find((y) => y.person_id === p)?.[k] ?? null)
    const fmtK = (v: number) => entero(v)
    const ultimo = (p: number, k: 'avg_kcal_per_day' | 'daily_kcal_target') => { const v = u.persons.find((y) => y.person_id === p)?.[k]; return v !== undefined ? entero(v) : t.comun.sinDato }
    return <div className="flex flex-col gap-4">
      <FiguraGrafica id="t4a" titulo={G.t4.tituloA} subtitulo={G.t4.subtitulo} filtros={filtros} dependencias={[d, ancho]} resumen={G.t4.resumenA(rango, dinero(u.cost_cents))}
        tabla={{ columnas: [G.semana, G.t4.costo], filas: w.map((x, i) => [cs[i], dinero(x.cost_cents)]) }}
        opcion={() => ({ ...ejes(cs, compacto, ancho), series: [{ name: G.t4.costo, type: 'line', symbol: 'circle', symbolSize: 8, lineStyle: { color: token('--ui-chart-1'), width: 2 }, itemStyle: { color: token('--ui-chart-1') }, label: { position: 'top' }, data: conEtiquetas(w.map((x) => x.cost_cents), (v) => dinero(v)) }] })} />
      <FiguraGrafica id="t4b" titulo={G.t4.tituloB} subtitulo={G.t4.subtitulo} dependencias={[d, ancho, sel.join()]} resumen={G.t4.resumenB(rango, vistas.map((p) => G.t4.resumenPersona(p.display_name, ultimo(p.id, 'avg_kcal_per_day'), ultimo(p.id, 'daily_kcal_target'))).join('; '))}
        tabla={{ columnas: [G.semana, ...vistas.flatMap((p) => [p.display_name, G.t4.metaDe(p.display_name)])], filas: w.map((x, i) => [cs[i], ...vistas.flatMap((p) => { const y = x.persons.find((z) => z.person_id === p.id); return [y?.avg_kcal_per_day !== undefined ? entero(y.avg_kcal_per_day) : t.comun.sinDato, y?.daily_kcal_target !== undefined ? entero(y.daily_kcal_target) : t.comun.sinDato] })]) }}
        opcion={() => ({ ...ejes(cs, fmtK, ancho), series: vistas.flatMap((p, i) => { const color = token(`--ui-person-${p.code}`) || token(`--ui-chart-${i + 1}`)
          return [
            { name: p.display_name, type: 'line', symbol: i ? 'rect' : 'circle', symbolSize: 8, lineStyle: { color, width: 2, type: i ? [8, 4] : 'solid' }, itemStyle: { color }, endLabel: { show: true, formatter: p.display_name, color: token('--ui-text-default') }, label: { position: 'top' }, data: conEtiquetas(kcal(p.id, 'avg_kcal_per_day'), fmtK) },
            { name: G.t4.metaDe(p.display_name), type: 'line', symbol: 'none', lineStyle: { color, width: 1, type: [2, 4] }, itemStyle: { color }, endLabel: { show: true, formatter: G.t4.metaDe(p.display_name), color: token('--ui-text-secondary') }, data: kcal(p.id, 'daily_kcal_target') },
          ] }) })} />
    </div>
  }}</Estados>
}

## FIN: frontend/src/pantallas/inicio/Tendencias.tsx

## ARCHIVO: frontend/src/pantallas/inicio/comun.tsx


// Piezas comunes de P-01 (U3 §1): tarjeta-enlace, insignias, estados por bloque, estado de un pago (§1.4) y datos compartidos.
import { Component, createContext, useContext, useEffect, useState, type ReactNode } from 'react'
import { ListFilter, X } from 'lucide-react'
import { Hoja } from '../../componentes/Hoja'
import { ChipSeleccionable } from '../../componentes/Chip'
import { textos } from '../../textos/es-MX'
import { ChevronRight, CircleAlert, CircleCheck, CircleDashed, CircleDot, CalendarClock, CalendarX, Info, RefreshCw, TriangleAlert } from 'lucide-react'
import { Enlace } from '../../app/enrutador'
import { usarConsulta, type EstadoConsulta } from '../../api/consulta'
import { horaCorta, hoyISO } from '../../api/formatos'
import { Boton } from '../../componentes/Boton'
import { Cargando, Esqueleto } from '../../componentes/Estados'
import { cn } from '../../componentes/utils'
import type { Esquemas, Respuesta } from '../../api/tipos'
import { difDias, quincenaDe, mesDe } from './fechas'
import { textosInicio as t } from './textos'

export const entero = (n: number) => String(Math.round(n)).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
/** Décimas de punto porcentual → «47.5 %» (C1 `progress_tenths_pct`). */
export const decimas = (n: number) => `${Math.trunc(n / 10)}.${Math.abs(n % 10)} %`

type RolInsignia = 'advertencia' | 'info' | 'exito' | 'neutro'
const INSIGNIA: Record<RolInsignia, string> = { advertencia: 'bg-role-warning-surface text-role-warning', info: 'bg-role-info-surface text-role-info', exito: 'bg-role-success-surface text-role-success', neutro: 'bg-surface-raised text-text-secondary border border-border-divider' }
export function Insignia({ rol, icono, children, etiqueta }: { rol: RolInsignia; icono: ReactNode; children: ReactNode; etiqueta?: string }) {
  // Ancho al contenido sin pasar del contenedor: el texto largo («Reposición de capital pendiente de confirmar con Kueski») se
  // parte en renglones en vez de desbordar la página a 475 px en Android (u5-tablero, U5-06). `min-h-6` conserva los 24 px.
  return <span aria-label={etiqueta} className={cn('inline-flex min-h-6 w-fit max-w-full items-center gap-1 rounded-sm px-2 text-small font-medium [&_svg]:size-5 [&_svg]:shrink-0', INSIGNIA[rol])}>{icono}<span className="min-w-0">{children}</span></span>
}
export const IconoAviso = () => <TriangleAlert aria-hidden />
export const IconoInfo = () => <Info aria-hidden />
export const IconoOk = () => <CircleCheck aria-hidden />

/** Pie de tarjeta o panel: texto de enlace con chevron (la zona entera es el enlace). */
export const PieEnlace = ({ texto }: { texto: string }) => (
  <span className="mt-auto inline-flex min-h-11 items-center gap-1 text-small font-medium text-text-link underline underline-offset-2">{texto}<ChevronRight aria-hidden className="size-5" /></span>
)
export function EnlaceAccion({ a, texto }: { a: string; texto: string }) {
  return <Enlace a={a} className="inline-flex min-h-11 items-center gap-1 rounded-sm text-small font-medium text-text-link underline underline-offset-2 hover:text-text-link-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus">{texto}<ChevronRight aria-hidden className="size-5" /></Enlace>
}

export function FallaBloque({ texto, alReintentar }: { texto: string; alReintentar?: () => void }) {
  return (
    <div role="alert" className="flex flex-col gap-2 rounded-md bg-role-error-surface p-3 text-small text-text-default">
      <span className="flex items-start gap-2"><CircleAlert aria-hidden className="size-5 shrink-0 text-role-error" />{texto}</span>
      {alReintentar && <Boton variante="secundario" className="w-max" onClick={alReintentar}><RefreshCw aria-hidden className="size-5" />{t.acciones.reintentar}</Boton>}
    </div>
  )
}
export const CargadoA = ({ desde }: { desde: number | null }) => (desde ? <p className="text-caption text-text-secondary">{t.comun.cargadoA(horaCorta(desde))}</p> : null)

/** Aísla un bloque: si sus datos no tienen la forma de C1, muestra su falla y las demás tarjetas siguen (U3 §1.2). */
export class Aislado extends Component<{ texto: string; alReintentar?: () => void; children: ReactNode }, { fallo: boolean }> {
  state = { fallo: false }
  componentDidCatch() { this.setState({ fallo: true }) }
  render() { return this.state.fallo ? <FallaBloque texto={this.props.texto} alReintentar={() => { this.setState({ fallo: false }); this.props.alReintentar?.() }} /> : this.props.children }
}

/** Estados de una lectura (U3 §1.3): esqueleto con la forma final, aviso a los 10 s, falla con su texto, sin conexión. */
export function Estados<T>({ consulta, falla, esqueleto, children }: { consulta: EstadoConsulta<T>; falla: string; esqueleto: ReactNode; children: (d: T) => ReactNode }) {
  if (consulta.cargando && consulta.datos === undefined) return <EsperaLarga alReintentar={consulta.recargar}>{esqueleto}</EsperaLarga>
  if (consulta.nuncaCargada) return <p role="status" className="rounded-md bg-role-warning-surface p-3 text-small text-text-default">{t.comun.nuncaCargada}</p>
  if (consulta.error || consulta.datos === undefined) return <FallaBloque texto={falla} alReintentar={consulta.recargar} />
  return <Aislado texto={falla} alReintentar={consulta.recargar}><CargadoA desde={consulta.desdeCopia} /><Pintar f={children} d={consulta.datos} /></Aislado>
}
const EsperaLarga = Cargando
/** Llama a la función de pintado dentro del límite de Aislado (un dato con otra forma no tumba la página). */
function Pintar<T>({ f, d }: { f: (d: T) => ReactNode; d: T }) { return <>{f(d)}</> }
export const EsqueletoTarjeta = () => <div className="flex flex-col gap-2"><Esqueleto className="w-1/2 bg-skeleton-on-card" /><Esqueleto className="h-8 w-2/3 bg-skeleton-on-card" /><Esqueleto className="w-5/6 bg-skeleton-on-card" /></div>
export const EsqueletoFilas = ({ n }: { n: number }) => <div className="flex flex-col gap-2">{Array.from({ length: n }, (_, i) => <Esqueleto key={i} className="h-10 bg-skeleton-on-card" />)}</div>

/** Recarga al guardar un gasto o un peso (fh:datos-cambiaron, CONVENCIONES §4). */
export function usarRecargaAlCambiar(...recargas: (() => void)[]) {
  useEffect(() => {
    const f = () => recargas.forEach((r) => r())
    window.addEventListener('fh:datos-cambiaron', f)
    return () => window.removeEventListener('fh:datos-cambiaron', f)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, recargas)
}

// ---- Datos compartidos de Inicio: se leen una sola vez por visita (A-02: las alertas se leen una vez) ----
export type Persona = Esquemas['Person']
type Lista<T> = { items: T[]; next_cursor?: string | null }
export interface Compartido {
  hoy: string
  alertas: EstadoConsulta<{ parameters: { 'alerts.payment_due_days_before': { value: number[] | null }; 'alerts.budget_warning_pct'?: { value: number | null }; 'alerts.budget_limit_pct'?: { value: number | null } } }>
  /** getDashboardSettings (C1 v1.5, K54): N de «Pagos próximos»; se lee una sola vez, como las alertas (U3 §1.2). */
  ajustesTablero: EstadoConsulta<Respuesta<'/api/v1/settings/dashboard', 'get'>>
  personas: EstadoConsulta<Lista<Persona>>
  creditos: EstadoConsulta<Lista<Esquemas['Credit']>>
  tiendas: EstadoConsulta<Lista<Esquemas['Store']>>
  categorias: EstadoConsulta<Lista<Esquemas['Category']>>
  proximo: EstadoConsulta<Respuesta<'/api/v1/dashboard/next-payment', 'get'>>
  flujo: EstadoConsulta<Respuesta<'/api/v1/dashboard/fortnight-flow', 'get'>>
  plan: EstadoConsulta<Respuesta<'/api/v1/dashboard/month-vs-plan', 'get'>>
}
const Ctx = createContext<Compartido | null>(null)
export const useInicio = () => useContext(Ctx)!
export function ProveedorInicio({ children }: { children: ReactNode }) {
  const hoy = hoyISO()
  const alertas = usarConsulta('get', '/api/v1/settings/alerts')
  const ajustesTablero = usarConsulta('get', '/api/v1/settings/dashboard')
  const personas = usarConsulta('get', '/api/v1/persons', { consulta: { limit: 200 } })
  const creditos = usarConsulta('get', '/api/v1/credits', { consulta: { limit: 200 } })
  const tiendas = usarConsulta('get', '/api/v1/stores', { consulta: { limit: 200 } })
  const categorias = usarConsulta('get', '/api/v1/categories', { consulta: { limit: 200 } })
  const proximo = usarConsulta('get', '/api/v1/dashboard/next-payment')
  const flujo = usarConsulta('get', '/api/v1/dashboard/fortnight-flow', { consulta: { period: quincenaDe(hoy) } })
  const plan = usarConsulta('get', '/api/v1/dashboard/month-vs-plan', { consulta: { period: mesDe(hoy) } })
  usarRecargaAlCambiar(plan.recargar, flujo.recargar)
  const valor = { hoy, alertas, ajustesTablero, personas, creditos, tiendas, categorias, proximo, flujo, plan } as Compartido
  return <Ctx.Provider value={valor}>{children}</Ctx.Provider>
}
export const listaDe = <T,>(c: EstadoConsulta<Lista<T>>): T[] => (Array.isArray(c.datos?.items) ? c.datos.items : [])
export const nombreCredito = (c: Compartido, id?: number) => listaDe(c.creditos).find((x) => x.id === id)?.name ?? t.comun.sinDato
export const karen = (c: Compartido) => listaDe(c.personas).find((p) => p.participation === 'external')
export const titulares = (c: Compartido) => listaDe(c.personas).filter((p) => p.participation !== 'external').sort((a, b) => ORDEN.indexOf(a.code) - ORDEN.indexOf(b.code))
const ORDEN = ['axel', 'ximena', 'karen', 'shared']
/** Días de aviso «Por vencer» (K44): solo de la lectura de alertas (C1-22); si falla, ninguno (U3 §1.2). */
export const diasAviso = (c: Compartido): number[] => { const v = c.alertas.datos?.parameters?.['alerts.payment_due_days_before']?.value; return Array.isArray(v) ? v : [] }
/** Color de tienda de C2 por nombre de la cadena. PROVISIONAL: C1 no trae un código de tienda (bloqueo WF3-inicio-tiendas). */
export function codigoTienda(s?: { name: string; kind: string }) {
  if (!s) return 'market'
  if (s.kind === 'market') return 'market'
  const n = s.name.normalize('NFD').replace(/[̀-ͯ']/g, '').toLowerCase().replace(/\s+/g, '-')
  return ['walmart', 'bodega-aurrera', 'sams-club'].find((c) => n.startsWith(c)) ?? 'market'
}

// ---- Estado de un pago (U3 §1.4, DESIGN §7.8): el primero que se cumpla ----
export interface EstadoPago { texto: string; icono: ReactNode; rol: RolInsignia; linea?: string }
type Parte = { person_id: number; amount_cents: number; paid_cents: number }
export function estadoPago(p: { status?: string; due_date: string; shares?: Parte[] }, hoy: string, aviso: number[], idKaren?: number): EstadoPago {
  const partes = p.shares ?? []
  const pagada = (s: Parte) => s.paid_cents >= s.amount_cents
  const k = partes.filter((s) => s.person_id === idKaren), h = partes.filter((s) => s.person_id !== idKaren)
  const falta = partes.some(pagada) && partes.some((s) => !pagada(s)) ? (k.length && k.every(pagada) ? t.pago.faltaHogar : h.length && h.every(pagada) ? t.pago.faltaKaren : undefined) : undefined
  if (p.status === 'paid') return { texto: t.pago.pagado, icono: <CircleCheck aria-hidden />, rol: 'exito' }
  if (p.status === 'overdue') return { texto: t.pago.vencido, icono: <CalendarX aria-hidden />, rol: 'advertencia', linea: [t.pago.sinMarcar, falta].filter(Boolean).join(' · ') }
  const d = difDias(hoy, p.due_date)
  if (aviso.includes(d)) return { texto: d === 0 ? t.pago.venceHoy : d === 1 ? t.pago.venceManana : t.pago.venceEn(d), icono: <CalendarClock aria-hidden />, rol: 'advertencia', linea: falta }
  if (falta) return { texto: falta, icono: <CircleDot aria-hidden />, rol: 'neutro' }
  return { texto: t.pago.pendiente, icono: <CircleDashed aria-hidden />, rol: 'neutro' }
}
export const InsigniaPago = ({ e }: { e: EstadoPago }) => <span className="flex flex-col gap-1"><Insignia rol={e.rol} icono={e.icono}>{e.texto}</Insignia>{e.linea && <span className="text-caption text-text-secondary">{e.linea}</span>}</span>

/** Título de panel de Detalle o bloque (font.size.title-3). */
export const TituloPanel = ({ children, id }: { children: ReactNode; id?: string }) => <h2 id={id} className="flex items-center gap-2 text-title-3 font-semibold text-text-default">{children}</h2>
export const Panel = ({ titulo, children, className, etiqueta }: { titulo: ReactNode; children: ReactNode; className?: string; etiqueta: string }) => (
  <section aria-label={etiqueta} className={cn('flex flex-col gap-3 rounded-md border border-border-divider bg-surface-card p-4 tablet:p-6', className)}><TituloPanel>{titulo}</TituloPanel>{children}</section>
)

// ---- Filtros (U3 §6): botón «Filtrar» que abre una hoja con chips; chips activos con «Quitar filtro» ----
export interface GrupoFiltro { id: string; nombre: string; multiple?: boolean; opciones: { valor: string; texto: string; marca?: ReactNode }[] }
export function Filtros({ grupos, valores, alCambiar }: { grupos: GrupoFiltro[]; valores: Record<string, string[]>; alCambiar: (id: string, v: string[]) => void }) {
  const [abierta, setAbierta] = useState(false)
  const activos = grupos.flatMap((g) => (valores[g.id] ?? []).map((v) => ({ g, v, texto: g.opciones.find((o) => o.valor === v)?.texto ?? v })))
  const alternar = (g: GrupoFiltro, v: string) => { const a = valores[g.id] ?? []; alCambiar(g.id, a.includes(v) ? a.filter((x) => x !== v) : g.multiple ? [...a, v] : [v]) }
  return (
    <div className="flex flex-wrap items-center gap-2">
      <Boton variante="secundario" onClick={() => setAbierta(true)}><ListFilter aria-hidden className="size-5" />{t.detalle.filtrar}</Boton>
      {activos.map(({ g, v, texto }) => (
        <ChipSeleccionable key={g.id + v} seleccionado aria-label={textos.vistas.quitarFiltro(`${g.nombre} ${texto}`)} onClick={() => alternar(g, v)}>{g.nombre}: {texto}<X aria-hidden className="size-5" /></ChipSeleccionable>
      ))}
      <Hoja abierta={abierta} alCambiar={setAbierta} titulo={t.detalle.filtros} pie={<Boton variante="secundario" tam="hoja" onClick={() => setAbierta(false)}>{t.detalle.listo}</Boton>}>
        <div className="flex flex-col gap-4">
          {grupos.map((g) => (
            <fieldset key={g.id} className="flex flex-col gap-2">
              <legend className="mb-2 text-small font-medium text-text-secondary">{g.nombre}</legend>
              <div className="flex flex-wrap gap-2">{g.opciones.map((o) => <ChipSeleccionable key={o.valor} seleccionado={(valores[g.id] ?? []).includes(o.valor)} marca={o.marca} onClick={() => alternar(g, o.valor)}>{o.texto}</ChipSeleccionable>)}</div>
            </fieldset>
          ))}
        </div>
      </Hoja>
    </div>
  )
}
/** Ancho de la ventana por puntos de quiebre de C2 (tablet 768, escritorio 1024, amplio 1440). */
export function usarMedia(consulta: string) {
  const [si, setSi] = useState(() => typeof matchMedia === 'function' && matchMedia(consulta).matches)
  useEffect(() => { const m = matchMedia(consulta); const f = () => setSi(m.matches); f(); m.addEventListener('change', f); return () => m.removeEventListener('change', f) }, [consulta])
  return si
}
export const TABLET = '(min-width: 48rem)'
export const AMPLIO = '(min-width: 90rem)'

## FIN: frontend/src/pantallas/inicio/comun.tsx

## ARCHIVO: frontend/src/pantallas/inicio/fechas.ts


// Fechas de Inicio (K8: quincenas 1–15 y 16–último día; semana ISO de lunes a domingo; K38: los domingos la lista vigente es la de la semana siguiente).
// Todo en fechas de calendario AAAA-MM-DD (America/Mexico_City vía hoyISO), sin horas ni flotantes.
const MESES = ['ene', 'feb', 'mar', 'abr', 'may', 'jun', 'jul', 'ago', 'sep', 'oct', 'nov', 'dic']
const utc = (iso: string) => { const [a, m, d] = iso.slice(0, 10).split('-').map(Number); return Date.UTC(a, m - 1, d) }
const aIso = (ms: number) => new Date(ms).toISOString().slice(0, 10)
export const sumarDias = (iso: string, n: number) => aIso(utc(iso) + n * 86_400_000)
/** Días de `a` a `b` (b − a). */
export const difDias = (a: string, b: string) => Math.round((utc(b) - utc(a)) / 86_400_000)
/** 0 = lunes … 6 = domingo. */
export const diaSemanaLunes = (iso: string) => (new Date(utc(iso)).getUTCDay() + 6) % 7
export const lunesDe = (iso: string) => sumarDias(iso, -diaSemanaLunes(iso))
export const ultimoDia = (aaaamm: string) => new Date(Date.UTC(Number(aaaamm.slice(0, 4)), Number(aaaamm.slice(5, 7)), 0)).getUTCDate()
/** Semana ISO 8601 «2026-W41» de una fecha. */
export function semanaISO(iso: string): string {
  const jueves = sumarDias(lunesDe(iso), 3)
  const anio = Number(jueves.slice(0, 4))
  const n = Math.floor(difDias(`${anio}-01-01`, jueves) / 7) + 1
  return `${anio}-W${String(n).padStart(2, '0')}`
}
/** Lista vigente (K38, U1 §0.2): los domingos, la semana que empieza el lunes siguiente. */
export const semanaVigente = (iso: string) => semanaISO(diaSemanaLunes(iso) === 6 ? sumarDias(iso, 1) : iso)
export const mesDe = (iso: string) => iso.slice(0, 7)
export const quincenaDe = (iso: string) => `${mesDe(iso)}-${Number(iso.slice(8, 10)) <= 15 ? 'q1' : 'q2'}`
export const diaDe = (iso: string) => Number(iso.slice(8, 10))
export const mesAbrev = (iso: string) => MESES[Number(iso.slice(5, 7)) - 1]
/** «2026-10» ± n meses. */
export function sumarMeses(aaaamm: string, n: number) {
  const t = Number(aaaamm.slice(0, 4)) * 12 + Number(aaaamm.slice(5, 7)) - 1 + n
  return `${Math.floor(t / 12)}-${String((t % 12) + 1).padStart(2, '0')}`
}
/** «Lun 5 al Dom 11 oct» (o «Lun 28 sep al Dom 4 oct» si cambia el mes). */
export function rangoSemana(inicio: string, fin: string) {
  const mismo = inicio.slice(5, 7) === fin.slice(5, 7)
  return `Lun ${diaDe(inicio)}${mismo ? '' : ` ${mesAbrev(inicio)}`} al Dom ${diaDe(fin)} ${mesAbrev(fin)}`
}
/** Día de pago de quincena (15 y 30, o el último día del mes si es menor; K8). */
export const esDiaDePago = (iso: string) => { const d = diaDe(iso); return d === 15 || d === Math.min(30, ultimoDia(mesDe(iso))) }
/** Rejilla de un mes de lunes a domingo (K8): semanas con 7 fechas; las de otros meses en null. */
export function rejillaMes(aaaamm: string): (string | null)[][] {
  const primero = `${aaaamm}-01`
  const u = ultimoDia(aaaamm)
  const celdas: (string | null)[] = Array.from({ length: diaSemanaLunes(primero) }, () => null)
  for (let d = 1; d <= u; d++) celdas.push(`${aaaamm}-${String(d).padStart(2, '0')}`)
  while (celdas.length % 7) celdas.push(null)
  return Array.from({ length: celdas.length / 7 }, (_, i) => celdas.slice(i * 7, i * 7 + 7))
}

## FIN: frontend/src/pantallas/inicio/fechas.ts

## ARCHIVO: frontend/src/pantallas/inicio/indice.ts


// Pantallas del área Inicio (CONVENCIONES-frontend §3): el cascarón las descubre solo.
import { lazy } from 'react'
export const pantallas = { 'P-01': lazy(() => import('./P01Inicio')) }

## FIN: frontend/src/pantallas/inicio/indice.ts

## ARCHIVO: frontend/src/pantallas/inicio/mocks.ts


// Manejadores de msw del tablero (solo «npm run dev:mocks»; nunca en el build): las operaciones nuevas de C1 v1.5 (K54) con
// datos SINTÉTICOS coherentes entre sí (no son del oráculo): ingreso de $20,000.00 al mes, pagos del hogar de la semilla
// (R-13: $598.82 en octubre; $1,197.64 y $1,482.57 en noviembre y diciembre) y los umbrales iniciales 30 % y 40 % (K54).
// El estado y el porcentaje van escritos como los mandaría el motor: la interfaz no los calcula.
import { http, HttpResponse, type HttpHandler } from 'msw'
const INGRESO = 2000000
const PAGOS: Record<string, [number, number, number, 'normal' | 'warning' | 'exceeded']> = {
  '2026-10': [59882, 97375, 3, 'normal'], '2026-11': [119764, 194750, 6, 'normal'], '2026-12': [148257, 223243, 7, 'normal'],
}
const parametro = (key: string, value: number, unit: string) => ({ key, value_type: 'integer', unit, value })
const ajustes = { parameters: {
  'debt.capacity_warning_pct': parametro('debt.capacity_warning_pct', 30, '%'), 'debt.capacity_limit_pct': parametro('debt.capacity_limit_pct', 40, '%'),
  'dashboard.payment_reminder_count': parametro('dashboard.payment_reminder_count', 5, 'pagos'), 'budget.signal_months': parametro('budget.signal_months', 3, 'meses'),
  'budget.signal_margin_pct': parametro('budget.signal_margin_pct', 10, '%'),
} }
const mesMas = (m: string, n: number) => { const t = Number(m.slice(0, 4)) * 12 + Number(m.slice(5, 7)) - 1 + n; return `${Math.floor(t / 12)}-${String((t % 12) + 1).padStart(2, '0')}` }
export const manejadores: HttpHandler[] = [
  http.get('/api/v1/settings/dashboard', () => HttpResponse.json(ajustes)),
  http.get('/api/v1/dashboard/debt-capacity', ({ request }) => {
    const month = new URL(request.url).searchParams.get('period') ?? '2026-10'
    const [h, tot, n, status] = PAGOS[month] ?? [0, 0, 0, 'normal']
    return HttpResponse.json({ month, income_cents: INGRESO, income_is_expected: false, extras_cents: 0, payments_cents: h, payments_total_cents: tot, payment_count: h ? 1 : 0, ratio_pct: n, ratio_total_pct: n + 2, warning_pct: 30, limit_pct: 40, scale_max_pct: 60, status, headroom_cents: 800000 - h, has_credits: true })
  }),
  http.get('/api/v1/dashboard/monthly-overview', ({ request }) => {
    const u = new URL(request.url), desde = u.searchParams.get('from') ?? '2025-11', hasta = u.searchParams.get('to') ?? '2027-03'
    const months = []
    for (let m = desde; m <= hasta && months.length < 24; m = mesMas(m, 1)) {
      const [h, tot, n, status] = PAGOS[m] ?? [0, 0, 0, 'normal'], proy = m > '2026-10'
      months.push({ month: m, is_projection: proy, income_total_cents: INGRESO, income_fortnights_cents: INGRESO, income_extras_cents: 0, income_is_expected: proy, ...(m === '2026-10' ? { accrued_cents: 119764, cash_cents: 59882 } : {}), credit_payments_cents: h, credit_payments_total_cents: tot, debt_ratio_pct: n, debt_status: status })
    }
    return HttpResponse.json({ warning_pct: 30, limit_pct: 40, months })
  }),
  http.get('/api/v1/budget/signals', ({ request }) => HttpResponse.json({ month: new URL(request.url).searchParams.get('period') ?? '2026-10', months_considered: 3, margin_pct: 10, rows: [] })),
]

## FIN: frontend/src/pantallas/inicio/mocks.ts

## ARCHIVO: frontend/src/pantallas/inicio/textos.ts


// Textos visibles de P-01 Inicio (U3 v1.5), en español de México (U5-07). Mismo estilo que src/textos/es-MX.ts.
const red = 'Revisa que estés en la red de casa.'
const dias = (d: number) => (d === 1 ? '1 día' : `${d} días`)
export const MESES_LARGOS = ['enero', 'febrero', 'marzo', 'abril', 'mayo', 'junio', 'julio', 'agosto', 'septiembre', 'octubre', 'noviembre', 'diciembre']
export const textosInicio = {
  acciones: { fotoTicket: 'Foto de ticket', capturarQuincena: 'Capturar quincena', reintentar: 'Reintentar', verTodo: 'Ver todo', verComoTabla: 'Ver como tabla', verComoGrafica: 'Ver como gráfica' },
  comun: { cargadoA: (h: string) => `Cargado a las ${h}`, nuncaCargada: 'Sin conexión · Esta página aún no se ha cargado en este teléfono', sinDato: '—', hogar: 'Hogar', karen: 'Karen', hogarYKaren: (h: string, k?: string) => (k ? `Hogar ${h} · Karen ${k}` : `Hogar ${h}`), activa: 'Activa', lectura: (l: string) => `Lectura ${l}`, provisional: 'Provisional' },
  pago: { pagado: 'Pagado', vencido: 'Vencido', sinMarcar: 'Sin marcar como pagado', venceHoy: 'Vence hoy', venceManana: 'Vence mañana', venceEn: (d: number) => `Vence en ${dias(d)}`, faltaKaren: 'Falta la parte de Karen', faltaHogar: 'Falta la parte del hogar', pendiente: 'Pendiente', esperado: 'Esperado' },
  t1: {
    titulo: 'Mes contra plan', teQuedan: 'te quedan', rebasado: 'rebasado', deDe: (r: string, p: string, n: string) => `${r} de ${p} · ${n}`, faltan: (d: number) => (d === 1 ? 'Falta 1 día del mes' : `Faltan ${d} días del mes`),
    vasEn: (n: number) => `Vas en ${n} %`, rebasaste: 'Rebasaste el plan', accion: 'Ver presupuesto', sinPlan: (m: string) => `Aún no hay presupuesto para ${m}`, llevas: (r: string) => `Llevas ${r} gastado`, accionVacio: 'Capturar presupuesto',
    falla: `No se pudo cargar el presupuesto del mes. ${red}`,
  },
  t2: {
    titulo: 'Flujo de la quincena', delAl: (a: number, b: number, m: string) => `Del ${a} al ${b} ${m}`, alCierre: 'al cierre', faltante: 'faltante', entra: 'Entra', sale: 'Sale', yaSalio: 'Ya salió', porPagar: 'Por pagar',
    esperado: 'Ingreso esperado', accion: 'Ver flujo', faltaIngreso: 'Falta tu ingreso quincenal', accionVacio: 'Capturar ingreso', falla: `No se pudo cargar el flujo de la quincena. ${red}`,
  },
  t3: {
    titulo: 'Próximo pago', variosPagos: (n: number) => `${n} pagos este día`, faltan: (d: number) => (d === 0 ? 'Hoy' : d === 1 ? 'Mañana' : `Faltan ${dias(d)}`), vencioHace: (d: number) => `Sin marcar como pagado · venció hace ${dias(d)}`,
    disponibleHoy: (m: string, l: string) => `Disponible hoy ${m} · Lectura ${l}`, pendienteConfirmar: 'Pendiente confirmar', reposicion: (n: string) => `Reposición de capital pendiente de confirmar con ${n}`, masLineas: (n: number) => `+${n} líneas`,
    accion: 'Ver crédito', sinPagos: 'Sin pagos pendientes', sinCreditos: 'Aún no hay créditos', accionVacio: 'Agregar crédito', falla: `No se pudieron cargar los pagos. ${red}`,
  },
  t4: {
    titulo: 'Comidas de la semana', costo: 'costo de la semana', sinArmar: (n: number) => (n === 1 ? '1 casilla sin armar' : `${n} casillas sin armar`), completa: 'Semana completa',
    persona: (nombre: string, kcal: string, meta: string | null, d: number) => `${nombre} · ${kcal} kcal al día ${meta ? `de ${meta}` : '· sin meta'} · ${d} de 7 días`, personaSinDias: (nombre: string) => `${nombre} · —`,
    accion: 'Ver plan', sinPlan: 'Aún no hay plan para esta semana', accionVacio: 'Armar plan', falla: `No se pudo cargar el plan de comidas. ${red}`,
  },
  t5: {
    titulo: 'Mejor canasta', menosQue: (d: string, t: string) => `${d} menos que ${t}`, menosQueSiguiente: (d: string) => `${d} menos que la siguiente tienda`, cobertura: (c: number, i: number) => `${c} de ${i} artículos con precio`,
    viejos: (n: number) => (n === 1 ? '1 precio viejo' : `${n} precios viejos`), membresia: 'Como si tuvieran membresía', listaVacia: (f: string) => `Tu lista está vacía; se compara la compra del ${f}`,
    accion: 'Comparar canasta', sinPrecios: 'Aún no hay precios para comparar', accionVacio: 'Registrar precios', falla: `No se pudo comparar la canasta. ${red}`,
  },
  t6: {
    titulo: 'Ahorro y metas', ahorrado: 'ahorrado', esteMes: (a: string, p: string) => `Este mes ${a} de ${p}`, meta: (n: string, p: string, f: string) => `${n} · ${p} · ${f}`,
    estados: { in_progress: 'En curso', overdue: 'Vencida', achieved: 'Cumplida' }, accion: 'Ver ahorro', sinAhorro: 'Aún no tienes fondo ni metas de ahorro', accionVacio: 'Crear meta', falla: `No se pudo cargar el ahorro. ${red}`,
  },
  /** Tablero de control (U3 2.0 §2, K54). */
  tablero: {
    region: 'Tablero de control', masDelHogar: 'Más del hogar',
    servidorViejo: 'Disponible cuando el servidor se actualice',
    comparar: (cifra: string, dif: number, m: string) => (dif > 0 ? `${cifra} más que ${m}` : dif < 0 ? `${cifra} menos que ${m}` : `Igual que ${m}`),
    sinDatosDe: (m: string) => `Sin datos de ${m}`,
    k1: {
      titulo: 'Gasto del mes', gastadoEn: (m: string) => `gastado en ${m}`, dePlan: (n: string, p: string, r: string) => `${n} de ${p} · te quedan ${r}`, rebasado: (m: string) => `${m} rebasado`,
      delPlan: (n: string) => `${n} del plan`, efectivo: (m: string) => `Pagado en efectivo ${m}`, gastado: 'Gastado', enEfectivo: 'Efectivo', barra: 'Avance del plan',
      vasEn: (n: number) => `Vas en ${n} %`, rebasaste: 'Rebasaste el plan', rebasadas: (n: number) => (n === 1 ? '1 categoría rebasada' : `${n} categorías rebasadas`),
      sinPlan: (m: string) => `Aún no hay presupuesto para ${m}`, llevas: (r: string) => `Llevas ${r} gastado`,
      accion: 'Ver presupuesto', accionVacio: 'Capturar presupuesto', falla: `No se pudo cargar el gasto del mes. ${red}`,
    },
    k2: {
      titulo: 'Ingreso del mes', deQuincenas: 'de quincenas', extra: (m: string) => `Extra este mes ${m}`, sinExtra: 'Sin extra este mes',
      capturadas: (n: number, m: number) => `${n} de ${m} quincenas capturadas`, esperado: 'Ingreso esperado', faltaIngreso: 'Falta tu ingreso quincenal',
      accion: 'Ver ingresos', accionVacio: 'Capturar ingreso', falla: `No se pudo cargar el ingreso del mes. ${red}`,
    },
    k3: {
      titulo: 'Créditos', deudaViva: 'deuda viva total', pagosDe: (m: string, total: string, n?: number) => `Pagos de ${m}: ${total}${n === undefined ? '' : n === 1 ? ' · 1 pago' : ` · ${n} pagos`}`,
      faltaPagar: (m: string) => `Falta pagar ${m}`, sinPagos: (m: string) => `Sin pagos en ${m}`, sinCreditos: 'Aún no hay créditos', disponibleSinDato: 'Disponible hoy —',
      accion: 'Ver créditos', accionVacio: 'Agregar crédito', falla: `No se pudieron cargar los créditos. ${red}`,
    },
    k4: {
      titulo: 'Pagos próximos', hayMas: 'Hay más pagos en Calendario', sinPagos: 'Sin pagos en los próximos 30 días', sinCreditos: 'Aún no hay créditos',
      fila: (f: string, l: string, m: string) => `${f} · ${l} · ${m}`, accion: 'Ver calendario', accionVacio: 'Agregar crédito', falla: `No se pudieron cargar los pagos. ${red}`,
    },
    b1: {
      titulo: 'Capacidad de endeudamiento', valorTexto: (n: string, tope: number, color: string) => `Capacidad de endeudamiento ${n} de ${tope} %: ${color}`,
      sinValorTexto: (tope: number) => `Capacidad de endeudamiento sin dato, tope ${tope} %`,
      pagosDe: (p: string, i: string) => `Pagos ${p} de un ingreso de ${i}`, teCaben: (m: string) => `Te caben ${m} más al mes antes del rojo`, arriba: (m: string) => `${m} arriba del tope`,
      conKaren: (n: number) => `Con la parte de Karen: ${n} %`, sinCreditos: 'Sin créditos activos', esperado: 'Ingreso esperado', faltaIngreso: 'Falta tu ingreso del mes para calcular la capacidad',
      accion: 'Ver créditos', accionVacio: 'Capturar ingreso', falla: `No se pudo calcular la capacidad de endeudamiento. ${red}`,
    },
    g1: {
      titulo: 'Ingreso, gasto y pagos por mes', subtitulo: (a: string, b: string) => `De ${a} a ${b}`, ingreso: 'Ingreso', gasto: 'Gasto', pagos: 'Pagos de crédito', quincenas: 'Quincenas', extras: 'Extras',
      nota: 'Los pagos de crédito pagan compras que ya cuentan en el gasto de su mes; no se suman al gasto',
      resumen: (m: string, i: string, g: string, p: string) => `En ${m}: ingreso ${i}, gasto ${g} y pagos de crédito ${p}`,
      vacio: 'Aún no hay meses con ingreso ni gasto', accionVacio: 'Capturar ingreso',
    },
    g2: {
      titulo: 'Deuda viva por línea', subtitulo: (f: string) => `Al ${f}`, linea: 'Línea', monto: 'Deuda viva', resumen: (t: string, l: string, m: string) => `Deuda viva hoy: ${t}; la mayor, ${l} con ${m}`, vacio: 'No hay deudas vivas',
    },
    g3: {
      titulo: 'Presupuesto por categoría', subtitulo: (m: string) => `Mes de ${m}`, etiqueta: (c: string, r: string, p: string, n: string) => `${c} · ${r} de ${p} · ${n}`, vasEn: (n: number) => `Vas en ${n} %`, rebasada: 'Rebasada',
      subir: 'Sugerencia: subir presupuesto', recortar: 'Sugerencia: recortar presupuesto',
      motivo: (c: string, n: number, arriba: boolean, m: number) => `${c}: ${n} meses seguidos ${arriba ? 'arriba' : 'abajo'} del plan por más de ${m} %`,
      marca: (n: number) => `${n} %`, verTodas: 'Ver todas', sinPlan: (n: number) => (n === 1 ? '1 categoría con gasto y sin plan' : `${n} categorías con gasto y sin plan`),
      resumen: (m: string, n: number, w: number, x: number) => `En ${m}: ${n} ${n === 1 ? 'categoría' : 'categorías'} arriba de ${w} % y ${x} ${x === 1 ? 'rebasada' : 'rebasadas'}`,
      categoria: 'Categoría', real: 'Real', plan: 'Plan', usado: '% del plan', sugerencia: 'Sugerencia', ninguna: 'Sin sugerencia', abrir: (c: string) => `Ver gastos de ${c}`,
      vacio: 'Aún no hay presupuesto por categoría', accionVacio: 'Capturar presupuesto',
    },
    g6: {
      titulo: 'Tendencia del barómetro', serie: 'Capacidad de endeudamiento', nota: 'Proyección con pagos programados e ingreso esperado', faltaIngreso: '— · Falta ingreso',
      estado: 'Estado', capacidad: 'Capacidad', proyeccion: 'Proyección', si: 'Sí', no: 'No', punto: (m: string, n: number, e: string) => `${m}, ${n} %: ${e}`,
      resumen: (m: string, n: string, e: string, mm: string, nn: string) => `${m}: ${n} (${e}); el mayor de la proyección es ${mm} con ${nn}`,
      resumenSinProyeccion: (m: string, n: string, e: string) => `${m}: ${n} (${e})`,
      vacio: 'Falta tu ingreso para calcular la capacidad', accionVacio: 'Capturar ingreso',
    },
  },
  detalle: {
    filtrar: 'Filtrar', filtros: 'Filtros', periodo: 'Periodo', mes: 'Mes', q1: '1.ª quincena', q2: '2.ª quincena', persona: 'Persona', familia: 'Familia', credito: 'Crédito', categoria: 'Categoría', tienda: 'Tienda',
    familias: { household: 'Hogar', external: 'Karen', joint: 'Conjunto' }, noUsa: (f: string) => `Este bloque no usa el filtro ${f}`, verGastos: 'Ver todos los gastos con estos filtros', listo: 'Listo', limpiar: 'Quitar filtros',
  },
  d1: { titulo: 'Presupuesto por categoría', fila: (r: string, p: string, n: string) => `${r} de ${p} · ${n}`, vacio: 'Aún no hay presupuesto por categoría', accionVacio: 'Capturar presupuesto', falla: `No se pudo cargar el presupuesto por categoría. ${red}` },
  d2: {
    titulo: 'Quincena en curso', ingresos: 'Ingresos', esperado: 'esperado', pagosCredito: 'Pagos de crédito (parte del hogar)', fijos: 'Fijos y venideros', variable: 'Gasto variable estimado', ahorro: 'Ahorro planificado', saldo: 'Saldo de la quincena', acumulado: 'Saldo acumulado',
    vacio: 'Falta tu ingreso quincenal', accionVacio: 'Capturar ingreso', falla: `No se pudo cargar la quincena. ${red}`,
  },
  d3: {
    titulo: 'Pagos y disponible', vencidos: 'Vencidos', proximos: 'Próximos pagos', masVencidos: 'Hay más pagos vencidos', disponible: 'Disponible por línea', minimo: (m: string) => `Mínimo del horizonte ${m}`,
    sinPagos: 'Sin pagos pendientes', sinCreditos: 'Aún no hay créditos', accionVacio: 'Agregar crédito', falla: `No se pudieron cargar los pagos. ${red}`, fallaDisponible: `No se pudo cargar el disponible. ${red}`,
  },
  d4: { titulo: 'Comidas de hoy', fuera: 'Fuera', kcal: (n: string) => `${n} kcal`, vacio: 'Hoy no hay comidas en el plan', accionVacio: 'Armar plan', falla: `No se pudieron cargar las comidas de hoy. ${red}` },
  d5: { titulo: 'Diferencias de precio', viejo: 'viejo', totales: 'Total por tienda', articulo: 'Artículo', gana: 'Gana', diferencia: 'Diferencia', vacio: 'Aún no hay precios para comparar', accionVacio: 'Registrar precios', falla: `No se pudo comparar la canasta. ${red}` },
  d6: {
    titulo: 'Metas', deDe: (a: string, m: string) => `${a} de ${m}`, aportacion: (m: string) => `Aportación necesaria ${m} al mes`, fechaMeta: (f: string) => `Fecha meta ${f}`, avance: (p: string) => `Avance ${p}`,
    vacio: 'Aún no tienes metas de ahorro', accionVacio: 'Crear meta', falla: `No se pudieron cargar las metas. ${red}`,
  },
  d7: {
    titulo: 'Participación de Karen', despensa: (m: string) => `Su parte de la despensa este mes: ${m}`, compras: 'Sus compras a crédito', pagos: 'Sus pagos de este mes', debe: (m: string) => `Lo que debe en la línea: ${m}`,
    nota: 'Karen paga directo a Kueski su parte; su parte es informativa y repone la línea cuando paga', comercio: 'Comercio', fecha: 'Fecha', total: 'Total', pagosDe: (h: number, t: number) => `Pagos ${h} de ${t}`, saldo: 'Saldo', monto: 'Su parte', estado: 'Estado',
    accion: 'Estado de cuenta de Karen', vacio: 'Karen no tiene compras ni pagos este mes', falla: `No se pudo cargar la participación de Karen. ${red}`,
  },
  graficas: {
    falla: `No se pudo cargar la gráfica. ${red}`, hoy: 'Hoy', proyeccion: 'Proyección', mes: 'Mes', semana: 'Semana', total: 'Total', todas: 'Todas las categorías', ultimos12: (c: string) => `Últimos 12 meses · ${c}`, sinFiltro: 'Todas',
    t1: { titulo: 'Plan contra real por mes', real: 'Real', plan: 'Plan', resumen: (m: string, r: string, p: string) => `En ${m}: real ${r} de un plan de ${p}`, resumenSinPlan: (m: string, r: string) => `En ${m}: real ${r}; aún sin plan`, vacio: 'Aún no hay presupuesto capturado', accionVacio: 'Capturar presupuesto' },
    t2: {
      titulo: 'Devengado contra efectivo', subtitulo: (a: string, b: string) => `De ${a} a ${b}`, efectivo: 'Efectivo acumulado', diferencia: 'Diferencia acumulada', devengado: 'Devengado acumulado', adelantado: 'Pagado por adelantado', adelantadoEn: (m: string) => `Pagado por adelantado en ${m}`,
      resumen: (u: string, s: string, d: string, e: string) => `Hasta ${u}, con proyección desde ${s}: devengado acumulado ${d} y efectivo acumulado ${e}`,
      conRepetidas: 'Proyección con despensas repetidas', sinRepetidas: 'Proyección sin despensas repetidas',
      columnas: { devengadoMes: 'Devengado', efectivoMes: 'Efectivo', devHogar: 'Devengado del hogar', devKaren: 'Devengado de Karen', efeHogar: 'Efectivo del hogar', efeKaren: 'Efectivo de Karen' },
      vacio: 'Aún no hay gastos ni pagos para comparar', accionVacio: 'Registrar gasto',
    },
    t3: { titulo: 'Deuda viva proyectada', total: 'Deuda viva total', resumen: (a: string, ma: string, b: string, mb: string) => `Deuda viva al cierre de ${a}: ${ma}; al cierre de ${b}: ${mb}`, resumenFiltro: (m: string, partes: string) => `Al cierre de ${m}: ${partes}`, interesesPagados: 'Intereses pagados', interesesPorPagar: 'Intereses por pagar', vacio: 'No hay deudas vivas' },
    t4: {
      tituloA: 'Costo por semana', tituloB: 'kcal al día por persona', subtitulo: 'Últimas 12 semanas', costo: 'Costo', metaDe: (n: string) => `Meta de ${n}`, semana: (a: string) => `Semana ${a}`,
      resumenA: (s: string, m: string) => `Semana ${s}: costo ${m}`, resumenPersona: (n: string, k: string, m: string) => `${n} ${k} kcal al día con meta de ${m}`, resumenB: (s: string, partes: string) => `Semana ${s}: ${partes}`,
      vacio: 'Aún no hay semanas con plan de comidas', accionVacio: 'Armar plan',
    },
  },
  calendario: {
    titulo: 'Calendario', anterior: 'Mes anterior', siguiente: 'Mes siguiente', hoy: 'Hoy', tipos: 'Tipos', mas30: 'Ver 30 días más', masN: (n: number) => `+${n} más`, q1: '1.ª quincena · del 1 al 15', q2: (u: number) => `2.ª quincena · del 16 al ${u}`,
    delAl: (a: string, b: string) => `del ${a} al ${b}`, diaSeleccionado: 'Día seleccionado', sinEventosDia: 'Sin eventos este día', semana: ['Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb', 'Dom'],
    tiposEvento: { income: 'Ingresos', credit_payment: 'Pagos de crédito', fixed_expense: 'Gastos fijos', upcoming_expense: 'Gastos venideros', card_statement: 'Fechas de corte', planned_grocery: 'Despensa programada', market_visit: 'Visita a mercados', promotion: 'Promociones', savings_goal: 'Metas de ahorro', monthly_report: 'Reporte mensual' },
    vacio: 'No hay eventos en este periodo', sinProgramados: 'Aún no hay pagos ni ingresos programados', accionVacio: 'Capturar ingreso', falla: `No se pudo cargar el calendario. ${red}`,
  },
} as const

## FIN: frontend/src/pantallas/inicio/textos.ts

## ARCHIVO: frontend/src/pantallas/gastos-presupuesto/P04Gasto.tsx


// P-04 Gasto (U1 §2.1, §2.2, §3.6: N-060 a N-066). Propiedades editables, reparto, devengo, efectivo, renglones y adjunto.
// WF-3a bis: foto del ticket (downloadAttachment) y «Quitar adjunto» (deleteAttachment, N-064). Expense solo trae
// has_attachment: el id del adjunto sale de queryView en galería (image_attachment_id) o de la respuesta al adjuntar
// (PROVISIONAL, bloqueo WF3a-bis-gastos-presupuesto: un PDF adjuntado antes no se puede ubicar).
import { useRef, useState } from 'react'
import { Copy, CreditCard, Paperclip, Trash2, CircleCheck, Unlink } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { navegar } from '../../app/enrutador'
import { armarUrl, nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import { aCentavos, dinero, fecha, hoyISO, mes } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { useAvisos } from '../../componentes/Avisos'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Esqueleto, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { textos } from '../../textos/es-MX'
import { CeldaEditable } from '../../vistas-bd/CeldaEditable'
import { VistaBD } from '../../vistas-bd/VistaBD'
import type { Propiedad } from '../../vistas-bd/tipos'
import { Bloque, Confirmar, Fila, centavosATexto, mensajeDe, useCatalogos, useSinRed } from './comun'
import { t } from './textos'

type Gasto = Esquemas['Expense']
const milesimas = (m: number) => `${Math.floor(m / 1000)}${m % 1000 ? `.${String(m % 1000).padStart(3, '0').replace(/0+$/, '')}` : ''}`

export default function P04Gasto({ params }: PropsPantalla) {
  const id = Number(params.id)
  const g = usarConsulta('get', '/api/v1/expenses/{expense_id}', { ruta: { expense_id: id } })
  if (g.error instanceof ErrorApi && g.error.estado === 404)
    return <Vacio texto={t.p04.noExiste} accion={<Boton onClick={() => navegar(rutaDe('P-02'))}>{textos.app.volver}</Boton>} />
  return (
    <ConDatos consulta={g} esqueleto={<div className="flex flex-col gap-3">{Array.from({ length: 8 }, (_, i) => <Esqueleto key={i} />)}</div>}>
      {(gasto) => <Detalle gasto={gasto} recargar={g.recargar} />}
    </ConDatos>
  )
}

function Detalle({ gasto, recargar }: { gasto: Gasto; recargar: () => void }) {
  const cats = useCatalogos()
  const avisar = useAvisos()
  const sinRed = useSinRed()
  const archivo = useRef<HTMLInputElement>(null)
  const [hoja, setHoja] = useState<'pagado' | 'eliminar' | 'quitar' | null>(null)
  const [subido, setSubido] = useState<{ id: number; mime: string } | null>(null)
  const galeria = usarConsulta('get', '/api/v1/views/{entity}', { ruta: { entity: 'expense' }, consulta: { mode: 'gallery', period: gasto.expense_date.slice(0, 7), limit: 200 } }, gasto.has_attachment)
  const fotoId = (galeria.datos?.items as Esquemas['ViewItem'][] | undefined)?.find((i) => i.id === gasto.id)?.image_attachment_id
  const adjunto = fotoId ? { id: fotoId, mime: 'image' } : subido
  const urlAdjunto = adjunto ? armarUrl('/api/v1/attachments/{attachment_id}/file', { attachment_id: adjunto.id }) : ''
  const [pago, setPago] = useState({ monto: centavosATexto(gasto.amount_cents), fecha: hoyISO(), error: '' })
  const ruta = { expense_id: gasto.id }
  const parchar = async (_: Gasto, cuerpo: Esquemas['ExpenseUpdate']) => { await pedir('patch', '/api/v1/expenses/{expense_id}', { ruta, cuerpo }); recargar() }
  const opc = <T extends { id: number }>(l: T[], n: (x: T) => string) => l.map((x) => ({ valor: String(x.id), texto: n(x) }))
  const props: Propiedad<Gasto>[] = [
    { id: 'expense_date', nombre: t.comun.fecha, valor: (x) => x.expense_date, pintar: (x) => fecha(x.expense_date), editar: { tipo: 'fecha', textoInicial: (x) => x.expense_date, guardar: (x, v) => parchar(x, { expense_date: v }) } },
    { id: 'amount_cents', nombre: t.comun.monto, valor: (x) => x.amount_cents, pintar: (x) => dinero(x.amount_cents), editar: { tipo: 'dinero', textoInicial: (x) => centavosATexto(x.amount_cents),
      vistaPrevia: async (x, v) => { const c = aCentavos(v); if (c == null || c <= 0) throw new Error(textos.registro.montoInvalido); return [{ etiqueta: t.comun.monto, antes: dinero(x.amount_cents), despues: dinero(c) }] },
      guardar: (x, v) => parchar(x, { amount_cents: aCentavos(v)! }) } },
    { id: 'category_id', nombre: t.comun.categoria, valor: (x) => cats.categoria(x.category_id), editar: { tipo: 'opcion', opciones: opc(cats.categorias, (c) => c.name), textoInicial: (x) => String(x.category_id), guardar: (x, v) => parchar(x, { category_id: Number(v) }) } },
    { id: 'person_id', nombre: t.comun.persona, valor: (x) => cats.persona(x.person_id), editar: { tipo: 'opcion', opciones: opc(cats.personas, (p) => p.display_name), textoInicial: (x) => String(x.person_id ?? ''), guardar: (x, v) => parchar(x, { person_id: v ? Number(v) : null }) } },
    { id: 'description', nombre: t.comun.descripcion, valor: (x) => x.description ?? '—', editar: { tipo: 'texto', textoInicial: (x) => x.description ?? '', guardar: (x, v) => parchar(x, { description: v || null }) } },
    { id: 'store_id', nombre: t.comun.tienda, valor: (x) => cats.tienda(x.store_id), editar: { tipo: 'opcion', opciones: opc(cats.tiendas, (s) => s.name), textoInicial: (x) => String(x.store_id ?? ''), guardar: (x, v) => parchar(x, { store_id: v ? Number(v) : null }) } },
    { id: 'payment_method', nombre: t.comun.formaPago, valor: (x) => textos.registro.formas[x.payment_method] },
    { id: 'kind', nombre: t.comun.tipo, valor: (x) => t.tipos[x.kind], editar: { tipo: 'opcion', opciones: Object.entries(t.tipos).map(([valor, texto]) => ({ valor, texto })), textoInicial: (x) => x.kind, guardar: (x, v) => parchar(x, { kind: v as Gasto['kind'] }) } },
    { id: 'cost_share_rule_id', nombre: t.comun.reparto, valor: (x) => cats.regla(x.cost_share_rule_id) },
    { id: 'accrual_months', nombre: t.p04.mesesDevengo, valor: (x) => x.accrual_months, editar: { tipo: 'numero', textoInicial: (x) => String(x.accrual_months),
      vistaPrevia: async (x, v) => { const n = Number(v); if (!Number.isInteger(n) || n < 1) throw new Error(textos.errores.revisaCampos); return [{ etiqueta: t.p04.mesesDevengo, antes: String(x.accrual_months), despues: String(n) }] },
      guardar: (x, v) => parchar(x, { accrual_months: Number(v) }) } },
    { id: 'notes', nombre: t.comun.nota, valor: (x) => x.notes ?? '—', editar: { tipo: 'texto', textoInicial: (x) => x.notes ?? '', guardar: (x, v) => parchar(x, { notes: v || null }) } },
  ]
  const accion = async (f: () => Promise<unknown>, ok: string, despues?: (r: unknown) => void) => {
    try { const r = await f(); avisar({ texto: ok, rol: 'exito' }); despues ? despues(r) : recargar() } catch (e) { avisar({ texto: mensajeDe(e), rol: 'error' }) }
  }
  const adjuntar = (f: File | undefined) => {
    if (!f) return
    const datos = new FormData(); datos.append('file', f); datos.append('kind', 'ticket_photo')
    accion(() => pedir('post', '/api/v1/expenses/{expense_id}/attachments', { ruta, cuerpo: datos, idempotencia: nuevaClave() }), t.p04.adjuntado,
      (r) => { const a = r as Esquemas['Attachment']; setSubido({ id: a.id, mime: a.mime_type }); recargar(); galeria.recargar() })
  }
  const irCompra = async () => {
    try { const c = await pedir('get', '/api/v1/credit-purchases/{credit_purchase_id}', { ruta: { credit_purchase_id: gasto.credit_purchase_id! } }); navegar(rutaDe('P-19', { id: c.credit_id, compra: c.id })) }
    catch (e) { avisar({ texto: mensajeDe(e), rol: 'error' }) }
  }
  const marcar = () => {
    const c = aCentavos(pago.monto)
    if (c == null || c <= 0) return setPago({ ...pago, error: textos.registro.montoInvalido })
    accion(() => pedir('post', '/api/v1/expenses/{expense_id}/mark-paid', { ruta, cuerpo: { amount_cents: c, expense_date: pago.fecha }, idempotencia: nuevaClave() }), t.p04.marcado, () => { setHoja(null); recargar() })
  }
  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-wrap gap-2">
        {!gasto.has_attachment && <Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => archivo.current?.click()}><Paperclip aria-hidden className="size-5" />{t.p04.adjuntar}</Boton>}
        <input ref={archivo} type="file" hidden accept="image/jpeg,image/png,image/webp,application/pdf" onChange={(e) => adjuntar(e.target.files?.[0])} />
        <Boton deshabilitadoPor={sinRed} onClick={() => accion(() => pedir('post', '/api/v1/expenses/{expense_id}/duplicate', { ruta, idempotencia: nuevaClave() }), t.p04.duplicado, (r) => navegar(rutaDe('P-04', { id: (r as Gasto).id })))}><Copy aria-hidden className="size-5" />{t.p04.duplicar}</Boton>
        {gasto.kind === 'upcoming' && <Boton deshabilitadoPor={sinRed} onClick={() => setHoja('pagado')}><CircleCheck aria-hidden className="size-5" />{t.p04.marcarPagado}</Boton>}
        {gasto.credit_purchase_id != null && <Boton onClick={irCompra}><CreditCard aria-hidden className="size-5" />{t.p04.compraCredito}</Boton>}
        <Boton variante="destructivo" deshabilitadoPor={sinRed} onClick={() => setHoja('eliminar')}><Trash2 aria-hidden className="size-5" />{t.comun.eliminar}</Boton>
      </div>
      {gasto.has_attachment && (
        <Bloque titulo={t.p04.ticketTitulo} accion={<Boton deshabilitadoPor={sinRed ?? (adjunto ? undefined : t.p04.sinUbicar)} onClick={() => setHoja('quitar')}><Unlink aria-hidden className="size-5" />{t.p04.quitarAdjunto}</Boton>}>
          {adjunto ? (
            adjunto.mime.startsWith('image')
              ? <a href={urlAdjunto} target="_blank" rel="noopener" className="self-start rounded-lg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus" aria-label={t.p04.abrirArchivo}>
                  <img src={urlAdjunto} alt={t.p04.fotoTicket} className="max-h-96 w-full rounded-lg border border-border-divider object-contain tablet:w-96" />
                </a>
              : <a href={urlAdjunto} target="_blank" rel="noopener" className="self-start text-body text-text-link underline">{t.p04.abrirArchivo}</a>
          ) : galeria.cargando ? <Esqueleto /> : <p className="text-small text-text-secondary">{t.p04.sinUbicar}</p>}
        </Bloque>
      )}
      <Bloque titulo={t.p04.propiedades}>
        <div>{props.map((p) => <Fila key={p.id} etiqueta={p.nombre}>{p.editar && !sinRed ? <CeldaEditable fila={gasto} prop={p}>{p.pintar ? p.pintar(gasto) : String(p.valor(gasto))}</CeldaEditable> : (p.pintar ? p.pintar(gasto) : String(p.valor(gasto)))}</Fila>)}</div>
      </Bloque>
      <Bloque titulo={t.p04.repartoTitulo}>
        <div>
          {(gasto.shares ?? []).map((s) => <Fila key={s.person_id} etiqueta={cats.persona(s.person_id)}>{dinero(s.amount_cents)}</Fila>)}
          <Fila etiqueta={t.p04.parteHogar}><span className="font-semibold">{dinero(gasto.household_share_cents ?? 0)}</span></Fila>
          <Fila etiqueta={t.p04.parteKaren}><span className="font-semibold">{dinero(gasto.external_share_cents ?? 0)}</span></Fila>
        </div>
      </Bloque>
      <div className="grid gap-6 desktop:grid-cols-2">
        <Bloque titulo={t.p04.devengo}><div>{(gasto.accruals ?? []).map((a) => <Fila key={a.month} etiqueta={mes(a.month)}>{dinero(a.amount_cents)}</Fila>)}</div></Bloque>
        <Bloque titulo={t.p04.efectivo}><div>{(gasto.cash_flows ?? []).map((c) => <Fila key={c.date} etiqueta={fecha(c.date)}>{dinero(c.amount_cents)}</Fila>)}</div></Bloque>
      </div>
      <Renglones id={gasto.id} />
      <Hoja abierta={hoja === 'pagado'} alCambiar={(a) => setHoja(a ? 'pagado' : null)} titulo={t.p04.marcarPagado}
        pie={<div className="flex justify-end gap-2"><Boton onClick={() => setHoja(null)}>{textos.app.cancelar}</Boton><Boton variante="primario" onClick={marcar}>{textos.app.guardar}</Boton></div>}>
        <div className="flex flex-col gap-3">
          <Campo etiqueta={t.p04.montoReal} inputMode="decimal" value={pago.monto} error={pago.error || undefined} onChange={(e) => setPago({ ...pago, monto: e.target.value, error: '' })} />
          <Campo etiqueta={t.p04.fechaReal} type="date" value={pago.fecha} onChange={(e) => setPago({ ...pago, fecha: e.target.value })} />
        </div>
      </Hoja>
      <Confirmar abierta={hoja === 'quitar'} alCambiar={(a) => setHoja(a ? 'quitar' : null)} titulo={t.p04.quitarAdjunto} texto={t.p04.confirmarQuitar} accion={t.p04.quitarAdjunto}
        alConfirmar={() => accion(() => pedir('delete', '/api/v1/attachments/{attachment_id}', { ruta: { attachment_id: adjunto!.id } }), t.p04.quitado, () => { setHoja(null); setSubido(null); recargar() })} />
      <Confirmar abierta={hoja === 'eliminar'} alCambiar={(a) => setHoja(a ? 'eliminar' : null)} titulo={t.comun.eliminar} texto={t.p04.confirmarEliminar} accion={t.comun.eliminar}
        alConfirmar={() => accion(() => pedir('delete', '/api/v1/expenses/{expense_id}', { ruta }), t.p04.eliminado, () => navegar(rutaDe('P-02')))} />
    </div>
  )
}

function Renglones({ id }: { id: number }) {
  const r = usarConsulta('get', '/api/v1/expenses/{expense_id}/items', { ruta: { expense_id: id }, consulta: { limit: 200 } })
  type Renglon = Esquemas['ExpenseItem']
  return (
    <Bloque titulo={t.p04.renglones}>
      <ConDatos consulta={r}>{(d) => (
        <VistaBD def={{
          clave: 'P-04-renglones', filas: d.items, idFila: (x: Renglon) => x.id, vistas: ['tabla'], titulo: (x: Renglon) => x.description, vacio: <Vacio texto={t.p04.sinRenglones} />,
          propiedades: [
            { id: 'line_no', nombre: '#', valor: (x: Renglon) => x.line_no },
            { id: 'description', nombre: t.p04.producto, valor: (x: Renglon) => x.description },
            { id: 'quantity_milli', nombre: t.p04.cantidad, valor: (x: Renglon) => x.quantity_milli, pintar: (x: Renglon) => milesimas(x.quantity_milli), alinear: 'derecha' },
            { id: 'unit_price_cents', nombre: t.p04.precioUnitario, valor: (x: Renglon) => x.unit_price_cents ?? null, pintar: (x: Renglon) => (x.unit_price_cents == null ? '—' : dinero(x.unit_price_cents)), alinear: 'derecha' },
            { id: 'line_total_cents', nombre: t.p04.importe, valor: (x: Renglon) => x.line_total_cents, pintar: (x: Renglon) => dinero(x.line_total_cents), alinear: 'derecha' },
            { id: 'household_share_cents', nombre: t.p04.parteHogar, valor: (x: Renglon) => x.household_share_cents ?? 0, pintar: (x: Renglon) => dinero(x.household_share_cents ?? 0), alinear: 'derecha' },
          ],
        }} />
      )}</ConDatos>
    </Bloque>
  )
}

## FIN: frontend/src/pantallas/gastos-presupuesto/P04Gasto.tsx

## ARCHIVO: frontend/src/pantallas/gastos-presupuesto/comun.tsx


// Piezas compartidas del área Gastos y Presupuesto: catálogos, selector de mes, barra de avance, selector nativo,
// confirmación y utilidades de periodo. Solo tokens C2; textos en ./textos.
import { useCallback, useEffect, useState, type ReactNode } from 'react'
import { ChevronLeft, ChevronRight, TriangleAlert } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import { todas } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { ErrorApi } from '../../api/errores'
import { aCentavos as aCentavosLocal, hoyISO, mes, porcentaje } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { conConsulta, navegar } from '../../app/enrutador'
import { Boton } from '../../componentes/Boton'
import { Campo as CampoTexto, estilosControl } from '../../componentes/Campo'
import { Hoja } from '../../componentes/Hoja'
import { cn } from '../../componentes/utils'
import { textos } from '../../textos/es-MX'
import { t } from './textos'

// ── Periodo (filtro global `period` de C1: aquí solo el mes calendario AAAA-MM) ──
export const mesActual = () => hoyISO().slice(0, 7)
export function periodoDe(consulta: URLSearchParams): string {
  const p = consulta.get('periodo')
  return p && /^\d{4}-(0[1-9]|1[0-2])$/.test(p) ? p : mesActual()
}
export function sumarMes(aaaamm: string, n: number): string {
  const total = Number(aaaamm.slice(0, 4)) * 12 + Number(aaaamm.slice(5, 7)) - 1 + n
  return `${Math.floor(total / 12)}-${String((total % 12) + 1).padStart(2, '0')}`
}
export const primerDia = (aaaamm: string) => `${aaaamm}-01`
export function ultimoDia(aaaamm: string): string {
  const d = new Date(Date.UTC(Number(aaaamm.slice(0, 4)), Number(aaaamm.slice(5, 7)), 0))
  return `${aaaamm}-${String(d.getUTCDate()).padStart(2, '0')}`
}
export function SelectorMes({ valor }: { valor: string }) {
  const ir = (n: number) => navegar(conConsulta({ periodo: sumarMes(valor, n) }), { reemplazar: true })
  return (
    <div className="flex items-center gap-1" role="group" aria-label={t.comun.periodo}>
      <Boton tam="icono" variante="fantasma" aria-label={t.comun.mesAnterior} onClick={() => ir(-1)}><ChevronLeft aria-hidden className="size-6" size={ICONO_MD} /></Boton>
      <span className="min-w-24 text-center text-body font-medium" aria-live="polite">{mes(valor)}</span>
      <Boton tam="icono" variante="fantasma" aria-label={t.comun.mesSiguiente} onClick={() => ir(1)}><ChevronRight aria-hidden className="size-6" size={ICONO_MD} /></Boton>
    </div>
  )
}

// ── Catálogos (categorías, personas, tiendas, reglas de reparto): una carga por sesión de la página ──
type Catalogos = { categorias: Esquemas['Category'][]; personas: Esquemas['Person'][]; tiendas: Esquemas['Store'][]; reglas: Esquemas['CostShareRule'][] }
let cache: Promise<Catalogos> | null = null
const cargar = () => (cache ??= Promise.all([
  todas('get', '/api/v1/categories').catch(() => []),
  todas('get', '/api/v1/persons').catch(() => []),
  todas('get', '/api/v1/stores').catch(() => []),
  todas('get', '/api/v1/cost-share-rules').catch(() => []),
]).then(([categorias, personas, tiendas, reglas]) => {
  const c = { categorias, personas, tiendas, reglas } as Catalogos
  if (!categorias.length) cache = null // reintenta en la siguiente pantalla
  return c
}))
export function useCatalogos() {
  const [c, setC] = useState<Catalogos>({ categorias: [], personas: [], tiendas: [], reglas: [] })
  useEffect(() => { let vivo = true; cargar().then((x) => vivo && setC(x)); return () => { vivo = false } }, [])
  const nombre = <T extends { id: number }>(lista: T[], id: number | null | undefined, campo: (x: T) => string) => (id == null ? '—' : (lista.find((x) => x.id === id) ? campo(lista.find((x) => x.id === id)!) : '—'))
  return {
    ...c,
    categoria: (id?: number | null) => nombre(c.categorias, id, (x) => x.name),
    persona: (id?: number | null) => nombre(c.personas, id, (x) => x.display_name),
    tienda: (id?: number | null) => nombre(c.tiendas, id, (x) => x.name),
    regla: (id?: number | null) => nombre(c.reglas, id, (x) => x.name),
  }
}

// ── Recarga cuando otra hoja guardó un gasto (fh:datos-cambiaron) ──
export function useAlCambiarDatos(...recargas: (() => void)[]) {
  const f = useCallback(() => recargas.forEach((r) => r()), recargas) // eslint-disable-line react-hooks/exhaustive-deps
  useEffect(() => { window.addEventListener('fh:datos-cambiaron', f); return () => window.removeEventListener('fh:datos-cambiaron', f) }, [f])
}

/** Motivo para desactivar una escritura sin conexión (U1 §2.4: solo lectura). */
export const useSinRed = () => (useConexion() ? undefined : textos.conexion.requiere)
export const mensajeDe = (e: unknown) => (e instanceof ErrorApi ? e.message : textos.errores.generico)

// ── Barra de avance (DESIGN §7.6): pista 8 px, llenado neutro; aviso y rebasado según el estado que da el API ──
export function BarraAvance({ pct, estado, excedente, nombre }: { pct: number; estado?: 'normal' | 'warning' | 'exceeded' | null; excedente?: string; /** Nombre accesible del medidor: categoría, quincena o meta (U5-03, aria-meter-name). */ nombre: string }) {
  const alerta = estado === 'warning' || estado === 'exceeded'
  return (
    <div className="flex flex-col gap-1">
      <div role="meter" aria-label={nombre} aria-valuemin={0} aria-valuemax={100} aria-valuenow={Math.min(pct, 100)} aria-valuetext={t.comun.delPlan(porcentaje(pct))} className="h-2 w-full overflow-hidden rounded-full bg-viz-track">
        <div className={cn('h-full rounded-full', alerta ? 'bg-role-warning' : 'bg-viz-meter')} style={{ width: `${Math.max(0, Math.min(pct, 100))}%` }} />
      </div>
      {alerta && (
        <span className="flex items-center gap-1 text-caption text-role-warning">
          <TriangleAlert aria-hidden className="size-5 shrink-0" />{t.comun.usado(porcentaje(pct))}{estado === 'exceeded' && <> · {t.comun.rebasado}{excedente && <span className="tabular-nums"> {excedente}</span>}</>}
        </span>
      )}
    </div>
  )
}

// ── Selector nativo con etiqueta (DESIGN §7.3) ──
export function Selector({ etiqueta, valor, opciones, alCambiar, className, vacio }: { etiqueta: string; valor: string; opciones: { valor: string; texto: string }[]; alCambiar: (v: string) => void; className?: string; vacio?: string }) {
  return (
    <label className={cn('flex flex-col gap-1 text-small text-text-secondary', className)}>
      {etiqueta}
      <select value={valor} onChange={(e) => alCambiar(e.target.value)} className={cn(estilosControl, 'w-full')}>
        {vacio !== undefined && <option value="">{vacio}</option>}
        {opciones.map((o) => <option key={o.valor} value={o.valor}>{o.texto}</option>)}
      </select>
    </label>
  )
}

/** Título de bloque dentro de la pantalla (h2), con acción opcional a la derecha. */
export function Bloque({ titulo, accion, children }: { titulo: string; accion?: ReactNode; children: ReactNode }) {
  return (
    <section className="flex flex-col gap-3">
      <div className="flex flex-wrap items-center justify-between gap-2"><h2 className="text-title-3 font-semibold text-text-default">{titulo}</h2>{accion}</div>
      {children}
    </section>
  )
}

/** Confirmación de una acción destructiva (U1 §2.2: «Eliminar (con confirmación)»). */
export function Confirmar({ abierta, alCambiar, titulo, texto, accion, alConfirmar }: { abierta: boolean; alCambiar: (a: boolean) => void; titulo: string; texto: string; accion: string; alConfirmar: () => Promise<void> }) {
  const [enCurso, setEnCurso] = useState(false)
  return (
    <Hoja abierta={abierta} alCambiar={alCambiar} titulo={titulo} descripcion={texto}
      pie={<div className="flex justify-end gap-2"><Boton onClick={() => alCambiar(false)}>{textos.app.cancelar}</Boton><Boton variante="destructivo" aria-busy={enCurso} onClick={async () => { setEnCurso(true); try { await alConfirmar() } finally { setEnCurso(false) } }}>{accion}</Boton></div>}>
      <p className="text-body text-text-default">{texto}</p>
    </Hoja>
  )
}

/** Fila de propiedad: etiqueta a la izquierda y valor a la derecha (móvil en dos renglones). */
export function Fila({ etiqueta, children }: { etiqueta: string; children: ReactNode }) {
  return (
    <div className="flex flex-col gap-1 border-b border-border-divider py-2 tablet:flex-row tablet:items-center tablet:gap-4">
      <span className="text-small text-text-secondary tablet:w-48 tablet:shrink-0">{etiqueta}</span>
      <div className="min-h-11 flex-1 text-body text-text-default tabular-nums tablet:min-h-9 coarse:min-h-11 flex items-center">{children}</div>
    </div>
  )
}
export const rejillaTarjetas = 'grid grid-cols-1 gap-3 tablet:grid-cols-2 desktop:grid-cols-3'
/** Centavos → texto editable «1234.56» sin flotantes (protocolo §6). */
export const centavosATexto = (c: number) => `${c < 0 ? '-' : ''}${Math.floor(Math.abs(c) / 100)}.${String(Math.abs(c) % 100).padStart(2, '0')}`

/** Hoja de captura con monto, fecha y campos extra (P-08 quincena/extra, P-09 movimiento). Valida monto > 0 (C1 MoneyCents). */
export function HojaCaptura({ titulo, abierta, alCerrar, montoInicial = '', fechaInicial, extra, alGuardar }: { titulo: string; abierta: boolean; alCerrar: () => void; montoInicial?: string; fechaInicial: string; extra?: ReactNode; alGuardar: (c: { monto: number; fecha: string; nota: string }) => Promise<void> }) {
  const [v, setV] = useState({ monto: montoInicial, fecha: fechaInicial, nota: '', error: '', general: '' })
  const [enCurso, setEnCurso] = useState(false)
  useEffect(() => { if (abierta) setV({ monto: montoInicial, fecha: fechaInicial, nota: '', error: '', general: '' }) }, [abierta, montoInicial, fechaInicial])
  const sinRed = useSinRed()
  const guardar = async () => {
    const m = aCentavosLocal(v.monto)
    if (m == null || m <= 0) return setV({ ...v, error: textos.registro.montoInvalido })
    setEnCurso(true)
    try { await alGuardar({ monto: m, fecha: v.fecha, nota: v.nota }); alCerrar() } catch (e) { setV({ ...v, general: mensajeDe(e) }) } finally { setEnCurso(false) }
  }
  return (
    <Hoja abierta={abierta} alCambiar={(a) => !a && alCerrar()} titulo={titulo}
      pie={<div className="flex justify-end gap-2"><Boton onClick={alCerrar}>{textos.app.cancelar}</Boton><Boton variante="primario" aria-busy={enCurso} deshabilitadoPor={sinRed} onClick={guardar}>{textos.app.guardar}</Boton></div>}>
      <div className="flex flex-col gap-3">
        {extra}
        <CampoTexto etiqueta={t.comun.monto} inputMode="decimal" value={v.monto} error={v.error || undefined} onChange={(e) => setV({ ...v, monto: e.target.value, error: '' })} />
        <CampoTexto etiqueta={t.comun.fecha} type="date" value={v.fecha} onChange={(e) => setV({ ...v, fecha: e.target.value })} />
        <CampoTexto etiqueta={t.comun.nota} value={v.nota} onChange={(e) => setV({ ...v, nota: e.target.value })} />
        {v.general && <p role="alert" className="text-small text-role-error">{v.general}</p>}
      </div>
    </Hoja>
  )
}

## FIN: frontend/src/pantallas/gastos-presupuesto/comun.tsx

## ARCHIVO: frontend/src/pantallas/comidas-compras/P15Productos.tsx


// P-15 Productos y precios (U1 §2.1–§2.2, §3.17 N-170 a N-176). Los precios de PROFECO no se editan (N-176).
// WF-3a bis: N-173 con listNutritionFacts y N-175 con edición en el lugar (updateProduct, updatePresentation y `Ctrl+/`).
import { useDeferredValue, useState } from 'react'
import { Package } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import type { PropsPantalla } from '../../app/registro'
import type { Esquemas } from '../../api/tipos'
import { nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { useConexion } from '../../api/conexion'
import { aCentavos, dinero, fecha, hoyISO } from '../../api/formatos'
import { conConsulta, navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Franja, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { Grafica } from '../../componentes/Grafica'
import { useAvisos } from '../../componentes/Avisos'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { CeldaEditable } from '../../vistas-bd/CeldaEditable'
import type { Propiedad } from '../../vistas-bd/tipos'
import { BarraAcciones, MenuMas, Seccion, Selector, Tiendita, aMilli, cantidad, filasDe, invalidar, mensajeError, numero, porId, usarCatalogo, nombreTienda } from './comun'
import { textosCC as t } from './textos'

type Producto = Esquemas['Product']
type Precio = Esquemas['Price']
type Presentacion = Esquemas['Presentation']
/** C1: base_unit_price_cents es el precio por 1 kg, 1 l o 1 pieza (K7); el contenido viene en milésimas de g, ml o pieza. */
const MILLI_POR_UNIDAD_PRECIO: Record<string, number> = { g: 1_000_000, ml: 1_000_000, piece: 1_000 }
/** a × b ÷ c con HALF_UP en enteros (§6: centavos enteros, nunca flotantes en el resultado). */
const mulDivHalfUp = (a: number, b: number, c: number) => Math.floor((2 * a * b + c) / (2 * c))
const EDITABLES = ['manual', 'ticket_manual', 'ticket_ocr']
const opc = (r: Record<string, string>) => Object.entries(r).map(([valor, texto]) => ({ valor, texto }))

export default function P15Productos({ sub, consulta }: PropsPantalla) {
  const conexion = useConexion()
  const tiendas = porId(usarCatalogo('tiendas'))
  const gruposL = usarCatalogo('grupos')
  const grupos = porId(gruposL)
  const [f, setF] = useState({ q: '', tienda: '', grupo: '', frescura: consulta.get('frescura') ?? '', fuente: '' })
  const q = useDeferredValue(f.q.trim())
  const activa = sub === 'detalle'
  const productos = usarConsulta('get', '/api/v1/products', { consulta: { q: q || undefined, store_id: f.tienda ? [Number(f.tienda)] : undefined, product_group_id: f.grupo ? Number(f.grupo) : undefined, freshness: (f.frescura || undefined) as Producto['freshness'], source: (f.fuente || undefined) as Precio['source'] | undefined, limit: 200 } }, activa)
  const resumen = usarConsulta('get', '/api/v1/prices/summary', {}, sub === 'resumen')
  const sinRed = conexion ? undefined : textos.conexion.requiere
  const hoja = consulta.get('hoja')
  const cerrar = () => navegar(conConsulta({ hoja: null, id: null }), { reemplazar: true })
  const mejor = (p: Producto) => p.best_price
  const opcionesGrupo = [{ valor: '', texto: t.sinCategoria }, ...(gruposL ?? []).map((g) => ({ valor: String(g.id), texto: g.name }))]
  const editarProducto = async (p: Producto, cuerpo: Esquemas['ProductUpdate']) => {
    await pedir('patch', '/api/v1/products/{product_id}', { ruta: { product_id: p.id }, cuerpo })
    invalidar('productos'); productos.recargar()
  }
  const grupoDe = (p: Producto) => (p.product_group_id ? String(p.product_group_id) : '')
  return (
    <div className="flex flex-col gap-4">
      <BarraAcciones acciones={<>
        <Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => navegar(conConsulta({ hoja: 'precio' }))}>{t.capturarPrecio}</Boton>
        <MenuMas acciones={[{ texto: t.fotoTicket, alHacer: () => navegar(rutaDe('P-05')) }, { texto: t.nuevoProducto, alHacer: () => navegar(conConsulta({ hoja: 'nuevo' })), deshabilitadoPor: sinRed }]} />
      </>} />
      {sub === 'resumen' && <ConDatos consulta={resumen}>{(r) => (
        <div className="flex flex-col gap-4">
          <div className="grid gap-3 tablet:grid-cols-3">
            {([['fresh', t.conPrecioFresco, r.fresh_count], ['stale', t.conPrecioViejo, r.stale_count], ['no_price', t.sinPrecioCuenta, r.no_price_count]] as const).map(([k, titulo, n]) =>
              <TarjetaMetrica key={k} titulo={titulo} valor={numero(n)} nota={t.frescuras[k]} rol={k === 'fresh' ? undefined : 'warning'} alAbrir={() => navegar(conConsulta({ sub: null, frescura: k }))} />)}
          </div>
          <Seccion titulo={t.ultimaFecha}>
            <ul className="flex flex-col gap-1">{r.stores.map((s) => <li key={s.store_id} className="flex items-center justify-between text-small"><Tiendita id={s.store_id} mapa={tiendas} /><span>{s.latest_observed_date ? fecha(s.latest_observed_date) : t.sinPrecio}</span></li>)}</ul>
          </Seccion>
        </div>
      )}</ConDatos>}
      {sub === 'tendencias' && <Historial />}
      {activa && <>
        <div className="grid gap-2 tablet:grid-cols-3 desktop:grid-cols-5">
          <Campo etiqueta={t.buscar} type="search" value={f.q} onChange={(e) => setF({ ...f, q: e.target.value })} />
          <Selector id="p15-tienda" etiqueta={t.tienda} valor={f.tienda} vacio={t.todas} alCambiar={(v) => setF({ ...f, tienda: v })} opciones={[...tiendas.values()].map((s) => ({ valor: String(s.id), texto: s.name }))} />
          <Selector id="p15-grupo" etiqueta={t.categoria} valor={f.grupo} vacio={t.todas} alCambiar={(v) => setF({ ...f, grupo: v })} opciones={(gruposL ?? []).map((g) => ({ valor: String(g.id), texto: g.name }))} />
          <Selector id="p15-frescura" etiqueta={t.frescura} valor={f.frescura} vacio={t.todas} alCambiar={(v) => setF({ ...f, frescura: v })} opciones={opc(t.frescuras)} />
          <Selector id="p15-fuente" etiqueta={t.fuente} valor={f.fuente} vacio={t.todas} alCambiar={(v) => setF({ ...f, fuente: v })} opciones={opc(t.fuentes)} />
        </div>
        <ConDatos consulta={productos}>{(d) => {
          const filas = filasDe<Producto>(d)
          if (!filas.length) return <Vacio icono={<Package aria-hidden className="size-6" size={ICONO_MD} />} texto={t.sinProductos} />
          return <VistaBD def={{
            clave: 'P-15-productos', filas, idFila: (p) => p.id, vistas: ['tabla', 'tablero', 'galeria'], titulo: (p) => p.name, alAbrir: (p) => navegar(conConsulta({ hoja: 'producto', id: String(p.id) })),
            columnasTablero: { propiedad: 'frescura', columnas: Object.values(t.frescuras) },
            propiedades: [
              { id: 'nombre', nombre: t.nombre, valor: (p) => p.name }, // abre la ficha; el nombre se edita en el lugar dentro de ella
              { id: 'categoria', nombre: t.categoria, valor: (p) => grupos.get(p.product_group_id ?? 0)?.name ?? '', pintar: (p) => grupos.get(p.product_group_id ?? 0)?.name ?? t.sinCategoria,
                editar: { tipo: 'opcion', opciones: opcionesGrupo, textoInicial: grupoDe, guardar: (p, v) => editarProducto(p, { product_group_id: v ? Number(v) : null }) } },
              { id: 'unidad', nombre: t.unidadBase, valor: (p) => t.unidadesNombre[p.base_unit] },
              { id: 'mejor', nombre: t.mejorPrecio, valor: (p) => mejor(p)?.base_unit_price_cents ?? null, alinear: 'derecha',
                pintar: (p) => { const b = mejor(p); return b ? <span className="flex flex-col items-end"><span className="tabular-nums">{t.porUnidad(dinero(b.base_unit_price_cents), t.unidadPrecio[p.base_unit])}</span><Tiendita id={b.store_id} mapa={tiendas} /></span> : '—' } },
              { id: 'ultimo', nombre: t.ultimoPrecio, valor: (p) => p.last_observed_date ?? null, pintar: (p) => (p.last_observed_date ? fecha(p.last_observed_date) : '—') },
              { id: 'frescura', nombre: t.frescura, valor: (p) => t.frescuras[p.freshness ?? 'no_price'], grupo: (p) => t.frescuras[p.freshness ?? 'no_price'] },
            ],
            // Edición en grupo (U1 §5.3 `Ctrl+/`): categoría de varios productos, con vista previa sumada antes de aplicar.
            editarVarias: {
              propiedades: [{ id: 'product_group_id', nombre: t.categoria, tipo: 'opcion', opciones: opcionesGrupo }],
              vistaPrevia: async (sel, _prop, v) => [{ etiqueta: t.cambianCategoria, antes: numero(0), despues: numero(sel.filter((p) => grupoDe(p) !== v).length) }],
              guardar: async (sel, _prop, v) => {
                for (const p of sel.filter((x) => grupoDe(x) !== v)) await pedir('patch', '/api/v1/products/{product_id}', { ruta: { product_id: p.id }, cuerpo: { product_group_id: v ? Number(v) : null } })
                invalidar('productos'); productos.recargar()
              },
            },
          }} />
        }}</ConDatos>
      </>}
      {hoja === 'producto' && consulta.get('id') && <HojaProducto id={Number(consulta.get('id'))} cerrar={cerrar} alCambiar={productos.recargar} />}
      {hoja === 'precio' && <HojaPrecio cerrar={cerrar} alGuardar={productos.recargar} />}
      {hoja === 'nuevo' && <HojaNuevo cerrar={cerrar} alGuardar={productos.recargar} />}
    </div>
  )
}

function Historial() {
  const productos = usarCatalogo('productos') ?? []
  const tiendas = porId(usarCatalogo('tiendas'))
  const [id, setId] = useState('')
  const precios = usarConsulta('get', '/api/v1/products/{product_id}/prices', { ruta: { product_id: Number(id) || 0 }, consulta: { sort: 'observed_date', limit: 200 } }, !!id)
  const p = productos.find((x) => String(x.id) === id)
  return (
    <div className="flex flex-col gap-3">
      <Selector id="p15-historial" etiqueta={t.producto} valor={id} alCambiar={setId} vacio={t.elegirProducto} opciones={productos.map((x) => ({ valor: String(x.id), texto: x.name }))} />
      {!id ? <Vacio texto={t.elegirParaHistorial} /> : <ConDatos consulta={precios}>{(d) => {
        const l = filasDe<Precio>(d).filter((x) => x.base_unit_price_cents != null)
        if (!l.length) return <Vacio texto={textos.vacio.grafica} />
        const fechas = [...new Set(l.map((x) => x.observed_date))].sort()
        const ids = [...new Set(l.map((x) => x.store_id))]
        // Precio vigente por tienda: el último capturado hasta cada fecha (se mantiene hasta que cambia)
        const series = ids.map((s) => { let ult = 0; return { nombre: nombreTienda(tiendas.get(s), tiendas), valores: fechas.map((fe) => { const x = l.filter((y) => y.store_id === s && y.observed_date === fe).pop(); if (x) ult = x.base_unit_price_cents ?? ult; return ult }) } })
        return <Grafica titulo={t.historial} resumen={t.resumenHistorial(p?.name ?? '')} tipo="lineas" categorias={fechas.map(fecha)} series={series} />
      }}</ConDatos>}
    </div>
  )
}

/** Esc dentro de una edición en el lugar descarta solo esa edición (U4 §0.2), no cierra la hoja: Radix escucha Esc antes
 *  que la celda, así que la hoja ignora el cierre mientras el foco está en un control de `[data-edicion-en-lugar]`. */
const editandoEnLugar = () => { const el = document.activeElement; return el instanceof HTMLElement && /^(INPUT|SELECT)$/.test(el.tagName) && !!el.closest('[data-edicion-en-lugar]') }

function HojaProducto({ id, cerrar, alCambiar }: { id: number; cerrar: () => void; alCambiar: () => void }) {
  const avisar = useAvisos()
  const conexion = useConexion()
  const tiendas = porId(usarCatalogo('tiendas'))
  const det = usarConsulta('get', '/api/v1/products/{product_id}', { ruta: { product_id: id } })
  const precios = usarConsulta('get', '/api/v1/products/{product_id}/prices', { ruta: { product_id: id }, consulta: { sort: '-observed_date', limit: 100 } })
  const nutricion = usarConsulta('get', '/api/v1/products/{product_id}/nutrition-facts', { ruta: { product_id: id }, consulta: { limit: 50 } })
  const [corrige, setCorrige] = useState<{ id: number; precio: string; confirmarBorrar?: boolean } | null>(null)
  const [nut, setNut] = useState({ kcal: '', cita: '' })
  const [pres, setPres] = useState({ nombre: '', tam: '' })
  const sinRed = conexion ? undefined : textos.conexion.requiere
  const hacer = async (f: () => Promise<unknown>) => { try { await f(); avisar({ texto: t.guardado, rol: 'exito' }); setCorrige(null); det.recargar(); precios.recargar(); nutricion.recargar() } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) } }
  // N-175: edición en el lugar de la presentación (CeldaEditable guarda con Enter o al salir; Esc descarta)
  const editarPres = async (x: Presentacion, cuerpo: Esquemas['PresentationUpdate']) => {
    await pedir('patch', '/api/v1/presentations/{presentation_id}', { ruta: { presentation_id: x.id }, cuerpo })
    det.recargar(); precios.recargar(); invalidar('productos'); alCambiar()
  }
  const propNombre: Propiedad<Esquemas['ProductDetail']> = { id: 'nombre', nombre: t.nombre, valor: (x) => x.name,
    editar: { tipo: 'texto', textoInicial: (x) => x.name, guardar: async (x, v) => {
      await pedir('patch', '/api/v1/products/{product_id}', { ruta: { product_id: x.id }, cuerpo: { name: v.trim() } })
      det.recargar(); invalidar('productos'); alCambiar()
    } } }
  return (
    <Hoja abierta alCambiar={(a) => !a && !editandoEnLugar() && cerrar()} titulo={det.datos?.name ?? t.producto}>
      <ConDatos consulta={det}>{(p) => {
        const nombreP = (pid: number) => p.presentations.find((x) => x.id === pid)?.name ?? '—'
        return (
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-1" data-edicion-en-lugar><span className="text-small font-medium text-text-secondary">{t.nombre}</span>
              <CeldaEditable fila={p} prop={propNombre}><span className="font-medium">{p.name}</span></CeldaEditable></div>
            <Seccion titulo={t.presentaciones}>
              {p.presentations.length > 0 && <div data-edicion-en-lugar><VistaBD def={{
                clave: 'P-15-presentaciones', filas: p.presentations, idFila: (x) => x.id, vistas: ['tabla'], titulo: (x) => x.name,
                propiedades: [
                  { id: 'nombre', nombre: t.presentacion, valor: (x) => x.name, editar: { tipo: 'texto', textoInicial: (x) => x.name, guardar: (x, v) => editarPres(x, { name: v.trim() }) } },
                  { id: 'marca', nombre: t.marca, valor: (x) => x.brand ?? '', pintar: (x) => x.brand || '—', editar: { tipo: 'texto', textoInicial: (x) => x.brand ?? '', guardar: (x, v) => editarPres(x, { brand: v.trim() || null }) } },
                  { id: 'tamano', nombre: `${t.tamano} (${t.unidades[p.base_unit]})`, valor: (x) => x.size_milli, alinear: 'derecha',
                    pintar: (x) => <span className="tabular-nums">{cantidad(x.size_milli, p.base_unit)}{x.is_bulk ? ` · ${t.granel}` : ''}</span>,
                    editar: { tipo: 'numero', textoInicial: (x) => String(x.size_milli / 1000),
                      // El tamaño cambia el precio por kg, l o pieza: «Antes → Después» con el último precio de la presentación
                      vistaPrevia: async (x, v) => {
                        const nuevo = aMilli(v)
                        if (nuevo === null) throw new Error(t.cantidadInvalida)
                        const suyos = filasDe<Precio>(precios.datos ?? { items: [] }).filter((y) => y.presentation_id === x.id)
                        const cambios = [{ etiqueta: t.tamano, antes: cantidad(x.size_milli, p.base_unit), despues: cantidad(nuevo, p.base_unit) }]
                        const u = suyos[0]
                        if (u && u.base_unit_price_cents != null) cambios.push({ etiqueta: t.ultimoPrecioUnidad(t.unidadPrecio[p.base_unit], nombreTienda(tiendas.get(u.store_id), tiendas), fecha(u.observed_date)), antes: dinero(u.base_unit_price_cents), despues: dinero(mulDivHalfUp(u.price_cents, MILLI_POR_UNIDAD_PRECIO[p.base_unit], nuevo)) })
                        if (suyos.length) cambios.push({ etiqueta: t.preciosSeRecalculan, antes: numero(0), despues: numero(suyos.length) })
                        return cambios
                      },
                      guardar: (x, v) => editarPres(x, { size_milli: aMilli(v)! }) } },
                ],
              }} /></div>}
              <div className="grid grid-cols-2 items-end gap-2">
                <Campo etiqueta={t.nuevaPresentacion} value={pres.nombre} onChange={(e) => setPres({ ...pres, nombre: e.target.value })} />
                <Campo etiqueta={`${t.tamano} (${t.unidades[p.base_unit]})`} inputMode="decimal" value={pres.tam} onChange={(e) => setPres({ ...pres, tam: e.target.value })} />
              </div>
              <Boton deshabilitadoPor={sinRed} disabled={!pres.nombre.trim() || aMilli(pres.tam) === null} onClick={() => hacer(() => pedir('post', '/api/v1/products/{product_id}/presentations', { ruta: { product_id: id }, cuerpo: { name: pres.nombre.trim(), size_milli: aMilli(pres.tam)! }, idempotencia: nuevaClave() }).then(() => setPres({ nombre: '', tam: '' })))}>{t.nuevaPresentacion}</Boton>
            </Seccion>
            <Seccion titulo={t.precios}>
              <ConDatos consulta={precios}>{(d) => (
                <ul className="flex flex-col divide-y divide-border-divider">{filasDe<Precio>(d).map((x) => (
                  <li key={x.id} className="flex flex-col gap-1 py-2">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <Tiendita id={x.store_id} mapa={tiendas} />
                      <span className="tabular-nums font-medium">{dinero(x.price_cents)}{x.is_stale && <span className="ml-1 rounded-sm bg-role-warning-surface px-1 text-caption text-role-warning">{t.viejo}</span>}</span>
                    </div>
                    <span className="text-caption tabular-nums text-text-secondary">{[nombreP(x.presentation_id), x.base_unit_price_cents != null ? t.porUnidad(dinero(x.base_unit_price_cents), t.unidadPrecio[p.base_unit]) : '', fecha(x.observed_date), t.fuentes[x.source], x.age_days != null ? t.diasN(x.age_days) : ''].filter(Boolean).join(' · ')}</span>
                    {!EDITABLES.includes(x.source) ? <span className="text-caption text-text-secondary">{t.profecoNoSeEdita}</span> : corrige?.id === x.id ? (
                      corrige.confirmarBorrar
                        ? <div className="flex flex-wrap items-center gap-2"><span className="text-small">{t.confirmarEliminarPrecio}</span><Boton onClick={() => setCorrige(null)}>{textos.app.cancelar}</Boton><Boton variante="destructivo" deshabilitadoPor={sinRed} onClick={() => hacer(() => pedir('delete', '/api/v1/prices/{price_id}', { ruta: { price_id: x.id } }))}>{t.eliminar}</Boton></div>
                        : <div className="flex flex-wrap items-end gap-2"><Campo etiqueta={t.precio} inputMode="decimal" value={corrige.precio} onChange={(e) => setCorrige({ ...corrige, precio: e.target.value })} ayuda={aCentavos(corrige.precio) ? `${dinero(x.price_cents)} → ${dinero(aCentavos(corrige.precio)!)}` : undefined} />
                          <Boton onClick={() => setCorrige(null)}>{textos.app.cancelar}</Boton>
                          <Boton variante="primario" deshabilitadoPor={sinRed} disabled={!aCentavos(corrige.precio)} onClick={() => hacer(() => pedir('patch', '/api/v1/prices/{price_id}', { ruta: { price_id: x.id }, cuerpo: { price_cents: aCentavos(corrige.precio)! } }))}>{textos.app.guardar}</Boton></div>
                    ) : <div className="flex gap-2"><Boton variante="enlace" onClick={() => setCorrige({ id: x.id, precio: String(x.price_cents / 100) })}>{t.corregir}</Boton><Boton variante="enlace" onClick={() => setCorrige({ id: x.id, precio: '', confirmarBorrar: true })}>{t.eliminar}</Boton></div>}
                  </li>
                ))}</ul>
              )}</ConDatos>
            </Seccion>
            <Seccion titulo={t.nutricion}>
              <ConDatos consulta={nutricion}>{(dn) => {
                const datos = filasDe<Esquemas['NutritionFact']>(dn)
                return datos.length ? <ul className="flex flex-col gap-1 text-small" data-prueba="nutricion">{datos.map((n) => <li key={n.id}><span className="tabular-nums">{t.kcal100(numero(n.kcal_milli_per_100g / 1000), p.base_unit === 'ml' ? 'ml' : 'g')}</span> · {t.fuentesNutricion[n.source]}{n.source_ref ? ` · ${n.source_ref}` : ''}{n.is_preferred ? ` · ${t.preferida}` : ''}</li>)}</ul> : <p className="text-small text-text-secondary">{t.sinNutricion}</p>
              }}</ConDatos>
              <div className="grid grid-cols-2 items-end gap-2">
                <Campo etiqueta={t.kcalPor100} inputMode="decimal" value={nut.kcal} onChange={(e) => setNut({ ...nut, kcal: e.target.value })} />
                <Campo etiqueta={t.cita} value={nut.cita} onChange={(e) => setNut({ ...nut, cita: e.target.value })} />
              </div>
              <Boton deshabilitadoPor={sinRed} disabled={aMilli(nut.kcal) === null || !nut.cita.trim()} onClick={() => hacer(() => pedir('post', '/api/v1/products/{product_id}/nutrition-facts', { ruta: { product_id: id }, cuerpo: { kcal_milli_per_100g: aMilli(nut.kcal)!, source: 'manual', source_ref: nut.cita.trim() }, idempotencia: nuevaClave() }).then(() => setNut({ kcal: '', cita: '' })))}>{t.agregarNutricion}</Boton>
            </Seccion>
          </div>
        )
      }}</ConDatos>
    </Hoja>
  )
}

function HojaPrecio({ cerrar, alGuardar }: { cerrar: () => void; alGuardar: () => void }) {
  const avisar = useAvisos()
  const conexion = useConexion()
  const productos = usarCatalogo('productos') ?? []
  const tiendas = usarCatalogo('tiendas') ?? []
  const [v, setV] = useState({ producto: '', presentacion: '', tienda: '', precio: '', fecha: hoyISO() })
  const det = usarConsulta('get', '/api/v1/products/{product_id}', { ruta: { product_id: Number(v.producto) || 0 } }, !!v.producto)
  const [error, setError] = useState<string | null>(null)
  async function guardar() {
    const c = aCentavos(v.precio)
    if (!v.presentacion || !v.tienda) return setError(t.elegir)
    if (!c) return setError(t.montoInvalido)
    try { await pedir('post', '/api/v1/prices', { cuerpo: { presentation_id: Number(v.presentacion), store_id: Number(v.tienda), price_cents: c, observed_date: v.fecha }, idempotencia: nuevaClave() }); avisar({ texto: t.precioGuardado, rol: 'exito' }); invalidar('productos'); alGuardar(); cerrar() } catch (e) { setError(mensajeError(e)) }
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={t.capturarPrecio} pie={<Boton tam="hoja" variante="primario" deshabilitadoPor={conexion ? undefined : textos.conexion.requiere} onClick={guardar}>{textos.app.guardar}</Boton>}>
      <div className="flex flex-col gap-3">
        <Selector id="p15-p-producto" etiqueta={t.producto} valor={v.producto} vacio={t.elegirProducto} alCambiar={(x) => setV({ ...v, producto: x, presentacion: '' })} opciones={productos.map((p) => ({ valor: String(p.id), texto: p.name }))} />
        <Selector id="p15-p-presentacion" etiqueta={t.presentacion} valor={v.presentacion} vacio={t.elegir} alCambiar={(x) => setV({ ...v, presentacion: x })} opciones={(det.datos?.presentations ?? []).map((p) => ({ valor: String(p.id), texto: `${p.name} · ${cantidad(p.size_milli, det.datos?.base_unit)}` }))} />
        <Selector id="p15-p-tienda" etiqueta={t.tienda} valor={v.tienda} vacio={t.elegir} alCambiar={(x) => setV({ ...v, tienda: x })} opciones={tiendas.filter((s) => s.is_active).map((s) => ({ valor: String(s.id), texto: s.name }))} />
        <Campo etiqueta={t.precio} inputMode="decimal" value={v.precio} onChange={(e) => setV({ ...v, precio: e.target.value })} />
        <Campo etiqueta={t.fecha} type="date" value={v.fecha} onChange={(e) => setV({ ...v, fecha: e.target.value })} ayuda={fecha(v.fecha)} />
        {error && <Franja rol="error" alerta>{error}</Franja>}
      </div>
    </Hoja>
  )
}

function HojaNuevo({ cerrar, alGuardar }: { cerrar: () => void; alGuardar: () => void }) {
  const avisar = useAvisos()
  const conexion = useConexion()
  const grupos = usarCatalogo('grupos') ?? []
  const [v, setV] = useState({ nombre: '', grupo: '', unidad: 'g', pres: '', tam: '' })
  const [error, setError] = useState<string | null>(null)
  async function guardar() {
    if (!v.nombre.trim()) return setError(t.faltaNombre)
    const tam = aMilli(v.tam)
    try {
      await pedir('post', '/api/v1/products', { cuerpo: { name: v.nombre.trim(), base_unit: v.unidad as Producto['base_unit'], product_group_id: v.grupo ? Number(v.grupo) : undefined, presentations: v.pres.trim() && tam ? [{ name: v.pres.trim(), size_milli: tam }] : undefined }, idempotencia: nuevaClave() })
      invalidar('productos'); avisar({ texto: t.productoCreado, rol: 'exito' }); alGuardar(); cerrar()
    } catch (e) { setError(mensajeError(e)) }
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={t.nuevoProducto} pie={<Boton tam="hoja" variante="primario" deshabilitadoPor={conexion ? undefined : textos.conexion.requiere} onClick={guardar}>{textos.app.guardar}</Boton>}>
      <div className="flex flex-col gap-3">
        <Campo etiqueta={t.nombre} value={v.nombre} onChange={(e) => setV({ ...v, nombre: e.target.value })} />
        <Selector id="p15-n-grupo" etiqueta={t.categoria} valor={v.grupo} vacio={t.elegir} alCambiar={(x) => setV({ ...v, grupo: x })} opciones={grupos.map((g) => ({ valor: String(g.id), texto: g.name }))} />
        <Selector id="p15-n-unidad" etiqueta={t.unidadBase} valor={v.unidad} alCambiar={(x) => setV({ ...v, unidad: x })} opciones={opc(t.unidadesNombre)} />
        <div className="grid grid-cols-2 gap-2">
          <Campo etiqueta={t.presentacionInicial} value={v.pres} onChange={(e) => setV({ ...v, pres: e.target.value })} />
          <Campo etiqueta={`${t.tamano} (${t.unidades[v.unidad]})`} inputMode="decimal" value={v.tam} onChange={(e) => setV({ ...v, tam: e.target.value })} />
        </div>
        {error && <Franja rol="error" alerta>{error}</Franja>}
      </div>
    </Hoja>
  )
}

## FIN: frontend/src/pantallas/comidas-compras/P15Productos.tsx

## ARCHIVO: frontend/src/pantallas/comidas-compras/comun.tsx


// Piezas compartidas de Comidas y Compras (P-10 a P-16): semanas ISO (K8, K38), catálogos del contrato,
// cantidades en milésimas (Milli de C1), marca de tienda y el menú «⋯». Solo tokens de C2.
import { useEffect, useState, type ReactNode } from 'react'
import { DropdownMenu } from 'radix-ui'
import { ChevronLeft, ChevronRight, Ellipsis } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import { todas } from '../../api/cliente'
import type { Esquemas } from '../../api/tipos'
import { fecha, hoyISO } from '../../api/formatos'
import { ErrorApi } from '../../api/errores'
import { conConsulta, navegar } from '../../app/enrutador'
import { Boton } from '../../componentes/Boton'
import { estilosControl } from '../../componentes/Campo'
import { ChipTienda } from '../../componentes/Chip'
import { textos } from '../../textos/es-MX'
import { textosCC as t } from './textos'

export type Tienda = Esquemas['Store']
export type Producto = Esquemas['Product']

// ── Semanas: lunes a domingo (K8); los domingos la semana por defecto es la siguiente (K38, U1 §0.2) ──
const DIA = 86_400_000
const utc = (iso: string) => Date.UTC(Number(iso.slice(0, 4)), Number(iso.slice(5, 7)) - 1, Number(iso.slice(8, 10)))
const aISO = (ms: number) => new Date(ms).toISOString().slice(0, 10)
export const sumarDias = (iso: string, n: number) => aISO(utc(iso) + n * DIA)
export function semanaDe(iso: string): string {
  const d = new Date(utc(iso))
  d.setUTCDate(d.getUTCDate() + 4 - (d.getUTCDay() || 7))
  const anio = d.getUTCFullYear()
  const n = Math.ceil(((+d - Date.UTC(anio, 0, 1)) / DIA + 1) / 7)
  return `${anio}-W${String(n).padStart(2, '0')}`
}
export function lunesDe(semana: string): string {
  const anio = Number(semana.slice(0, 4))
  const ene4 = Date.UTC(anio, 0, 4)
  return aISO(ene4 - ((new Date(ene4).getUTCDay() || 7) - 1) * DIA + (Number(semana.slice(6)) - 1) * 7 * DIA)
}
export const diasDe = (semana: string) => Array.from({ length: 7 }, (_, i) => sumarDias(lunesDe(semana), i))
export const moverSemana = (semana: string, n: number) => semanaDe(sumarDias(lunesDe(semana), 7 * n))
export const semanaPorDefecto = (hoy = hoyISO()) => semanaDe(new Date(utc(hoy)).getUTCDay() === 0 ? sumarDias(hoy, 1) : hoy)
export const textoSemana = (s: string) => t.semana(Number(s.slice(6)), fecha(lunesDe(s)), fecha(sumarDias(lunesDe(s), 6)))
const DIAS_CORTOS = ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb']
export const diaCorto = (iso: string) => `${DIAS_CORTOS[new Date(utc(iso)).getUTCDay()]} ${Number(iso.slice(8, 10))}`

/** Semana que muestra la pantalla: la de ?semana= o la por defecto, fijada al abrir (U1 §0.2). */
export function useSemana(consulta: URLSearchParams) {
  const [porDefecto] = useState(() => semanaPorDefecto())
  const q = consulta.get('semana')
  const semana = q && /^\d{4}-W\d{2}$/.test(q) ? q : porDefecto
  const cambiar = (s: string) => navegar(conConsulta({ semana: s === porDefecto ? null : s }), { reemplazar: true })
  return { semana, porDefecto, cambiar }
}
export const conSemana = (ruta: string, semana: string, porDefecto: string) => (semana === porDefecto ? ruta : `${ruta}?semana=${semana}`)

export function SelectorSemana({ semana, porDefecto, cambiar }: ReturnType<typeof useSemana>) {
  return (
    <div role="group" aria-label={t.cambiarSemana} className="flex flex-wrap items-center gap-1">
      <Boton tam="icono" variante="fantasma" aria-label={t.semanaAnterior} onClick={() => cambiar(moverSemana(semana, -1))}><ChevronLeft aria-hidden className="size-6" size={ICONO_MD} /></Boton>
      <span className="text-body font-medium" data-prueba="semana">{textoSemana(semana)}</span>
      <Boton tam="icono" variante="fantasma" aria-label={t.semanaSiguiente} onClick={() => cambiar(moverSemana(semana, 1))}><ChevronRight aria-hidden className="size-6" size={ICONO_MD} /></Boton>
      {semana === porDefecto
        ? <span className="rounded-full bg-surface-selected px-2 text-caption text-text-secondary">{t.semanaVigente}</span>
        : <Boton variante="enlace" onClick={() => cambiar(porDefecto)}>{t.irSemanaVigente}</Boton>}
    </div>
  )
}

// ── Catálogos (una carga por sesión de la página; se invalidan al crear) ──
type Catalogos = { tiendas: Tienda; productos: Producto; personas: Esquemas['Person']; tiempos: Esquemas['MealSlot']; recetas: Esquemas['Recipe']; grupos: Esquemas['ProductGroup'] }
type Clave = keyof Catalogos
const cargadores: { [K in Clave]: () => Promise<unknown> } = {
  tiendas: () => todas('get', '/api/v1/stores', { consulta: { limit: 200 } }),
  productos: () => todas('get', '/api/v1/products', { consulta: { limit: 200 } }),
  personas: () => todas('get', '/api/v1/persons'),
  tiempos: () => todas('get', '/api/v1/meal-slots', { consulta: { limit: 200 } }),
  recetas: () => todas('get', '/api/v1/recipes', { consulta: { limit: 200 } }),
  grupos: () => todas('get', '/api/v1/product-groups', { consulta: { limit: 200 } }),
}
const cache: Partial<Record<Clave, Promise<unknown>>> = {}
export const invalidar = (k: Clave) => { delete cache[k] }
export function usarCatalogo<K extends Clave>(k: K): Catalogos[K][] | undefined {
  const [datos, setDatos] = useState<Catalogos[K][]>()
  useEffect(() => {
    let vivo = true
    const p = (cache[k] ??= cargadores[k]().catch((e) => { delete cache[k]; throw e })) as Promise<Catalogos[K][]>
    p.then((v) => vivo && setDatos(v), () => vivo && setDatos([]))
    return () => { vivo = false }
  }, [k])
  return datos
}
export const porId = <T extends { id: number }>(l?: T[]) => new Map((l ?? []).map((x) => [x.id, x]))
/** Personas con plan de comidas: titular y co-titular (Karen es externa, decisión A). */
export const personasPlan = (l?: Esquemas['Person'][]) => (l ?? []).filter((p) => p.participation !== 'external' && p.is_active)

// ── Tiendas: la marca de color sale del nombre de la cadena (C1 no trae código de tienda: PROVISIONAL) ──
export function codigoTienda(s: Tienda | undefined, mapa?: Map<number, Tienda>): string {
  if (!s) return 'market'
  if (s.kind === 'market') return 'market'
  if (s.kind === 'branch' && s.parent_id && mapa?.get(s.parent_id)) return codigoTienda(mapa.get(s.parent_id), mapa)
  return s.name.normalize('NFD').replace(/[̀-ͯ']/g, '').toLowerCase().trim().replace(/\s+/g, '-')
}
/** Nombre visible: una sucursal con cadena se muestra «Cadena · Sucursal» (PROVISIONAL, bloqueo WF3a-bis-comidas-compras punto 3). */
export function nombreTienda(s: Tienda | undefined, mapa?: Map<number, Tienda>): string {
  if (!s) return '—'
  const cadena = s.kind === 'branch' && s.parent_id ? mapa?.get(s.parent_id) : undefined
  return cadena ? `${cadena.name} · ${s.name}` : s.name
}
export function Tiendita({ id, mapa }: { id?: number | null; mapa: Map<number, Tienda> }) {
  const s = id ? mapa.get(id) : undefined
  return s ? <ChipTienda codigo={codigoTienda(s, mapa)} nombre={nombreTienda(s, mapa)} /> : <span className="text-text-secondary">{t.sinTienda}</span>
}

// ── Cantidades (Milli: milésimas de la unidad base) ──
const num = new Intl.NumberFormat('es-MX', { maximumFractionDigits: 2 })
export const numero = (n: number) => num.format(n)
export function cantidad(milli: number | null | undefined, unidad?: string | null): string {
  if (milli == null) return '—'
  const v = milli / 1000
  if (unidad === 'g') return v >= 1000 ? `${num.format(v / 1000)} kg` : `${num.format(v)} g`
  if (unidad === 'ml') return v >= 1000 ? `${num.format(v / 1000)} l` : `${num.format(v)} ml`
  return `${num.format(v)} ${v === 1 ? 'pza' : 'pzas'}`
}
export const kcal = (n: number | null | undefined) => (n == null ? '—' : `${num.format(n)} kcal`)
export const porciones = (milli: number) => num.format(milli / 1000)
/** «1.5» → 1500 (enteros, sin flotantes); null si no es mayor que cero. */
export function aMilli(texto: string): number | null {
  const s = texto.trim().replace(',', '.')
  if (!/^\d+(\.\d{1,3})?$/.test(s)) return null
  const [e, d = ''] = s.split('.')
  const v = Number(e) * 1000 + Number(d.padEnd(3, '0'))
  return v > 0 ? v : null
}
export const mensajeError = (e: unknown) => (e instanceof ErrorApi ? e.message : t.noGuardado)

// ── Controles pequeños ──
export function Selector({ etiqueta, valor, alCambiar, opciones, vacio, id }: { etiqueta: string; valor: string; alCambiar: (v: string) => void; opciones: { valor: string; texto: string }[]; vacio?: string; id: string }) {
  return (
    <label htmlFor={id} className="flex flex-col gap-1 text-small font-medium">
      {etiqueta}
      <select id={id} value={valor} onChange={(e) => alCambiar(e.target.value)} className={`${estilosControl} coarse:h-11`}>
        {vacio !== undefined && <option value="">{vacio}</option>}
        {opciones.map((o) => <option key={o.valor} value={o.valor}>{o.texto}</option>)}
      </select>
    </label>
  )
}
export interface AccionMas { texto: string; alHacer: () => void; deshabilitadoPor?: string }
export function MenuMas({ acciones }: { acciones: AccionMas[] }) {
  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger asChild>
        <Boton tam="icono" variante="secundario" aria-label={textos.app.masOpciones}><Ellipsis aria-hidden className="size-6" size={ICONO_MD} /></Boton>
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content align="end" sideOffset={4} className="z-50 min-w-56 rounded-md border border-border-divider bg-surface-raised p-1 text-text-default shadow-1">
          {acciones.map((a) => (
            <DropdownMenu.Item key={a.texto} disabled={!!a.deshabilitadoPor} onSelect={a.alHacer} className="flex h-9 cursor-pointer items-center rounded-sm px-3 text-small outline-none data-[highlighted]:bg-surface-hover data-[disabled]:cursor-not-allowed data-[disabled]:text-text-disabled coarse:h-11">
              {a.texto}{a.deshabilitadoPor && <span className="sr-only"> · {a.deshabilitadoPor}</span>}
            </DropdownMenu.Item>
          ))}
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  )
}
/** Barra de acciones de la pantalla: a la izquierda el contexto, a la derecha la primaria y «⋯». */
export const BarraAcciones = ({ children, acciones }: { children?: ReactNode; acciones: ReactNode }) => (
  <div className="flex flex-col gap-2 tablet:flex-row tablet:items-center tablet:justify-between">
    <div className="min-w-0">{children}</div>
    <div className="flex flex-wrap items-center gap-2">{acciones}</div>
  </div>
)
export const Seccion = ({ titulo, children, accion }: { titulo: string; children: ReactNode; accion?: ReactNode }) => (
  <section className="flex flex-col gap-2">
    <div className="flex items-center justify-between gap-2"><h2 className="text-title-3 font-semibold">{titulo}</h2>{accion}</div>
    {children}
  </section>
)
/** Filas de una página de C1 (Page.items llega como unknown[] en el tipo generado por la composición allOf). */
export const filasDe = <T,>(d: { items?: unknown[] }) => (d.items ?? []) as T[]

## FIN: frontend/src/pantallas/comidas-compras/comun.tsx

## ARCHIVO: frontend/pruebas/inicio/datos.ts


// Datos de prueba de P-01 con las cifras del oráculo y la semilla (U3: R-01, R-11, R-13 a R-15, R-45, R-46, R-61, R-62, R-64, M-01, M-08).
// Formas de C1 v1.5 (tablero K54: M-18 a M-22); ingresos y plan son marcadores (U3 los deja como marcador): el oráculo no trae
// ingresos, así que con la semilla K-2, B-1, G-1 y G-6 muestran su vacío por falta de ingreso. Sin credenciales.
import type { BrowserContext } from '@playwright/test'
const t = '2026-10-01T09:00:00-06:00'
const pago = (id: number, due_date: string) => ({ id, credit_purchase_id: 1, seq: id, due_date, amount_cents: 97375, principal_cents: 90000, interest_cents: 7375, created_at: t, updated_at: t, credit_id: 1, status: 'pending', shares: [{ person_id: 1, amount_cents: 59882, paid_cents: 0, status: 'pending' }, { person_id: 3, amount_cents: 37493, paid_cents: 0, status: 'pending' }] })
/** Meses «AAAA-MM» de `a` a `b`, inclusive. */
export function rangoMeses(a: string, b: string) {
  const r: string[] = []
  for (let n = Number(a.slice(0, 4)) * 12 + Number(a.slice(5, 7)) - 1; ; n++) { const m = `${Math.floor(n / 12)}-${String((n % 12) + 1).padStart(2, '0')}`; r.push(m); if (m >= b || r.length > 36) return r }
}
const meses = Array.from({ length: 12 }, (_, i) => { const n = 9 + i; return `${2026 + Math.floor(n / 12)}-${String((n % 12) + 1).padStart(2, '0')}` })
export function respuestas(): Record<string, unknown | ((u: URL) => unknown)> {
  return {
    '/api/v1/settings/alerts': { parameters: { 'alerts.payment_due_days_before': { key: 'alerts.payment_due_days_before', value_type: 'integer_list', unit: 'días', value: [1, 0] } } },
    '/api/v1/credits': { items: [{ id: 1, name: 'Kueski', kind: 'revolving_line', holder_person_id: 1, payment_frequency: 'fortnightly', availability_reading: 'B', availability_confirmed: false, is_active: true, created_at: t, updated_at: t }], next_cursor: null },
    '/api/v1/stores': { items: [['Walmart', 'chain'], ['Bodega Aurrera', 'chain'], ["Sam's Club", 'chain']].map(([name, kind], i) => ({ id: i + 1, name, kind, assisted_reading_enabled: false, requires_membership: i === 2, has_membership: false, is_active: true, created_at: t, updated_at: t })), next_cursor: null },
    '/api/v1/dashboard/month-vs-plan': { month: '2026-10', actual_cents: 119764, days_left: 22, cash_cents: 59882, previous_month: '2026-09', categories_warning_count: 0, categories_exceeded_count: 0 },
    '/api/v1/settings/dashboard': { parameters: Object.fromEntries(([['debt.capacity_warning_pct', 30, '%'], ['debt.capacity_limit_pct', 40, '%'], ['dashboard.payment_reminder_count', 5, 'pagos'], ['budget.signal_months', 3, 'meses'], ['budget.signal_margin_pct', 10, '%']] as const).map(([key, value, unit]) => [key, { key, value_type: 'integer', unit, value }])) },
    '/api/v1/incomes/month': { month: '2026-10', fortnights: [], extras: [], captured_cents: 0, expected_pending_cents: 0, extras_cents: 0, total_cents: 0, previous_month: '2026-09' },
    '/api/v1/credits/totals': { live_debt_cents: 389500, fortnight_total_cents: 0, fortnight_household_cents: 0, fortnight_external_cents: 0, month_total_cents: 97375, month_household_cents: 59882, month_external_cents: 37493, interest_paid_month_cents: 0, interest_paid_year_cents: 0, month_payment_count: 1, month_pending_household_cents: 59882 },
    '/api/v1/dashboard/debt-capacity': { month: '2026-10', income_is_expected: false, extras_cents: 0, payments_cents: 59882, payments_total_cents: 97375, payment_count: 1, warning_pct: 30, limit_pct: 40, scale_max_pct: 60, has_credits: true },
    '/api/v1/dashboard/monthly-overview': (u: URL) => ({ warning_pct: 30, limit_pct: 40, months: rangoMeses(u.searchParams.get('from')!, u.searchParams.get('to')!).map((month) => {
      const pagos: Record<string, [number, number]> = { '2026-10': [59882, 97375], '2026-11': [119764, 194750], '2026-12': [148257, 223243] }
      const [h, tot] = pagos[month] ?? [0, 0]
      return { month, is_projection: month > '2026-10', income_is_expected: false, ...(month === '2026-10' ? { accrued_cents: 119764, cash_cents: 59882 } : {}), credit_payments_cents: h, credit_payments_total_cents: tot }
    }) }),
    '/api/v1/budget/categories': { month: '2026-10', rows: [] },
    '/api/v1/budget/signals': { month: '2026-10', months_considered: 3, margin_pct: 10, rows: [] },
    '/api/v1/debt/live': { date: '2026-10-10', total_cents: 389500, credits: [{ credit_id: 1, total_cents: 389500 }], purchases: [] },
    '/api/v1/dashboard/fortnight-flow': { start_date: '2026-10-01', end_date: '2026-10-15', income_cents: 1500000, income_is_expected: true, spent_cents: 119764, committed_cents: 0, closing_balance_cents: 1380236 },
    '/api/v1/dashboard/next-payment': { next_payment: { due_date: '2026-10-30', credit_id: 1, amount_cents: 97375, household_share_cents: 59882, external_share_cents: 37493, days_left: 20, payment_count: 1, is_overdue: false }, availability: [{ credit_id: 1, availability_reading: 'B', available_cents: 690550, availability_confirmed: false }] },
    '/api/v1/dashboard/week-meals': { week_start: '2026-10-05', week_end: '2026-10-11', cost_cents: 0, persons: [], unplanned_slots: 21 },
    '/api/v1/dashboard/best-basket': { basket: 'expense', expense_id: 10, store_id: 2, total_cents: 309450, difference_vs_second_cents: 13850, second_store_id: 1, covered_count: 33, item_count: 36, stale_price_count: 0, assumes_membership: false },
    '/api/v1/expenses/10': { id: 10, expense_date: '2026-10-10', amount_cents: 309450, category_id: 1, person_id: 1, created_at: t, updated_at: t },
    '/api/v1/dashboard/savings': { saved_total_cents: 0, contributed_month_cents: 0, planned_month_cents: 0 },
    '/api/v1/persons/3/participation': { person_id: 3, month: '2026-10', accrued_groceries_cents: 74986, purchases: [{ credit_purchase_id: 1, description: 'Despensa 1', purchase_date: '2026-10-10', total_cents: 389500, payments_made: 0, payments_total: 4, balance_cents: 149972 }], payment_shares: [{ scheduled_payment_id: 1, due_date: '2026-10-30', amount_cents: 37493, status: 'pending' }], live_debt_cents: 149972 },
    // Filtra por `status` y por `period` = «desde..hasta» como el servidor (K-4 pide de hoy a hoy + 30; D-3, a hoy + 60).
    '/api/v1/scheduled-payments': (u: URL) => {
      const [a, b] = (u.searchParams.get('period') ?? '..').split('..'), st = u.searchParams.get('status')
      const items = [pago(1, '2026-10-30'), pago(2, '2026-11-15'), pago(3, '2026-11-30')].filter((p) => (!st || p.status === st) && (!a || p.due_date >= a) && (!b || p.due_date <= b))
      return { items, next_cursor: null }
    },
    '/api/v1/credits/1/availability-preview': { credit_id: 1, availability_reading: 'B', availability_confirmed: false, reading_a: { available_today_cents: 1000000, horizon_min_cents: 809450 }, reading_b: { available_today_cents: 690550, horizon_min_cents: 500000 } },
    '/api/v1/cash-flow/fortnights': { projection_mode: 'contracted_only', rows: [{ start_date: '2026-10-01', end_date: '2026-10-15', pay_date: '2026-10-15', income_cents: 1500000, income_is_expected: true, credit_payments_cents: 0, fixed_and_upcoming_cents: 0, variable_estimate_cents: 0, balance_cents: 1380236, cumulative_balance_cents: 1380236, is_actual: true }] },
    '/api/v1/cash-flow/accrual-vs-cash': { projection_mode: 'contracted_only', months: meses.map((month, i) => {
      const acc = i < 2 ? 194750 : 0, cash = i === 0 ? 97375 : i === 1 ? 194750 : 0
      const cumAcc = i === 0 ? 194750 : 389500, cumCash = i === 0 ? 97375 : 292125
      return { month, accrued_cents: acc, cash_cents: cash, accrued_household_cents: acc ? 119764 : 0, accrued_external_cents: acc ? 74986 : 0, cash_household_cents: cash ? 59882 * (i + 1) : 0, cash_external_cents: cash ? 37493 * (i + 1) : 0, cumulative_accrued_cents: cumAcc, cumulative_cash_cents: cumCash, cumulative_difference_cents: cumAcc - cumCash, is_canonical: i < 6 }
    }) },
    '/api/v1/budget/monthly-trend': { months: [] },
    '/api/v1/credits/debt-trend': { months: meses.map((month, i) => ({ month, live_debt_cents: Math.max(0, 389500 - i * 97375), credits: [{ credit_id: 1, live_debt_cents: Math.max(0, 389500 - i * 97375) }], interest_paid_cents: 0, interest_pending_cents: 0 })) },
    '/api/v1/meal-plan/trends': { weeks: [] },
    '/api/v1/calendar/events': { items: [
      { event_type: 'income', date: '2026-10-15', title: 'Quincena de Axel', amount_cents: 1500000, status: 'expected', ref_table: 'income', ref_id: 1 },
      { event_type: 'credit_payment', date: '2026-10-30', title: 'Pago', amount_cents: 97375, household_share_cents: 59882, external_share_cents: 37493, status: 'pending', credit_id: 1, ref_table: 'scheduled_payment', ref_id: 1 },
    ], next_cursor: null },
  }
}
/** Instala las respuestas de Inicio (después de ApiFalsa); `cambios` reemplaza o falla rutas: número = estado HTTP de error. */
export async function instalarInicio(ctx: BrowserContext, cambios: Record<string, unknown> = {}, registro: string[] = []) {
  const r = { ...respuestas(), ...cambios }
  await ctx.route(/\/api\/v1\//, async (route) => {
    const u = new URL(route.request().url())
    if (route.request().method() !== 'GET' || !(u.pathname in r)) return route.fallback()
    registro.push(u.pathname + u.search)
    const crudo = r[u.pathname], v = typeof crudo === 'function' ? crudo(u) : crudo
    if (typeof v === 'number') return route.fulfill({ status: v, contentType: 'application/json', body: JSON.stringify({ code: 'internal', message: 'Falla de prueba', request_id: 'prueba' }) })
    return route.fulfill({ status: 200, contentType: 'application/json', headers: { 'X-Request-Id': 'prueba' }, body: JSON.stringify(v) })
  })
}

## FIN: frontend/pruebas/inicio/datos.ts

## ARCHIVO: frontend/pruebas/inicio/inicio.escritorio.spec.ts


// P-01 Inicio (U3, U5-06/07/09, C3-03): corre a 412×915 táctil y a 1440×900. API simulada con context.route.
import { expect, test, type Locator, type Page } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { instalarInicio } from './datos'
const reloj = (page: Page, iso = '2026-10-10T12:00:00-06:00') => page.clock.setFixedTime(new Date(iso))
const tarjeta = (page: Page, n: number | string) => page.locator(`[data-prueba="tarjeta-${n}"]`)
// Desde U3 2.0 (K54) el Resumen es el tablero (pruebas en tablero.escritorio.spec.ts); las tarjetas 2, 4, 5 y 6 siguen en «Más del hogar».

test('P-01 Más del hogar: tarjetas 2, 4, 5 y 6 en orden, cifras de la semilla y una petición por tarjeta', async ({ page, context }) => {
  const reg: string[] = []
  await new ApiFalsa().instalar(context); await instalarInicio(context, {}, reg); await reloj(page)
  await page.goto('/inicio')
  await expect(page.getByRole('heading', { level: 1, name: 'Inicio' })).toBeVisible()
  const tarjetas = page.getByTestId('mas-del-hogar').locator('[data-prueba^="tarjeta-"]')
  await expect(tarjetas).toHaveCount(4)
  expect(await tarjetas.evaluateAll((t) => t.map((x) => x.getAttribute('data-prueba')))).toEqual([2, 4, 5, 6].map((n) => `tarjeta-${n}`))
  await expect(tarjeta(page, 2)).toContainText('Del 1 al 15 oct')
  await expect(tarjeta(page, 2)).toContainText('Ingreso esperado')
  await expect(tarjeta(page, 4)).toContainText('Aún no hay plan para esta semana')
  for (const txt of ['Bodega Aurrera $3,094.50', '$138.50 menos que Walmart', 'Tu lista está vacía; se compara la compra del 10-oct-2026']) await expect(tarjeta(page, 5)).toContainText(txt)
  await expect(tarjeta(page, 6)).toContainText('Aún no tienes fondo ni metas de ahorro')
  // A-02, A-03, A-04: semana en curso para comidas y lista vigente (K38) para la canasta; getAlertSettings una sola vez.
  expect(reg).toContain('/api/v1/dashboard/week-meals?period=2026-W41')
  expect(reg).toContain('/api/v1/dashboard/best-basket?week=2026-W41')
  expect(reg.filter((r) => r.startsWith('/api/v1/settings/alerts'))).toHaveLength(1)
  // A-08 y A-09: cada tarjeta es un solo enlace ≥ 44 px; sin desplazamiento horizontal.
  for (const n of [2, 4, 5, 6]) {
    const t = tarjeta(page, n)
    expect(await t.evaluate((x) => x.tagName)).toBe('A')
    await expect(t.locator('a, button')).toHaveCount(0)
    expect((await t.boundingBox())!.height).toBeGreaterThanOrEqual(44)
  }
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true)
  const texto = await page.locator('main').innerText()
  expect(texto).not.toMatch(/\b(Save|Cancel|Settings|Dashboard|Loading|Error|Home|Delete|Search|Add)\b/)
  expect(texto).not.toMatch(/\d{4}-\d{2}-\d{2}/)
})

test('P-01 estados: una tarjeta con falla no tumba las demás y «Reintentar» la recupera', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await instalarInicio(context, { '/api/v1/dashboard/month-vs-plan': 500 }); await reloj(page)
  await page.goto('/inicio')
  await expect(tarjeta(page, 'K-1')).toContainText('No se pudo cargar el gasto del mes. Revisa que estés en la red de casa.')
  await expect(tarjeta(page, 'K-1').getByRole('link', { name: 'Ver presupuesto' })).toBeVisible()
  await expect(tarjeta(page, 'K-3')).toContainText('Pagos de octubre: $973.75 · 1 pago')
  await context.unrouteAll({ behavior: 'ignoreErrors' })
  await new ApiFalsa().instalar(context); await instalarInicio(context)
  await tarjeta(page, 'K-1').getByRole('button', { name: 'Reintentar' }).click()
  await expect(tarjeta(page, 'K-1')).toContainText('Llevas $1,197.64 gastado')
})

test('P-01 A-17: el Jue 29 oct 2026 el pago de K-4 dice «Vence mañana»', async ({ page, context }) => {
  await new ApiFalsa().instalar(context)
  await instalarInicio(context, { '/api/v1/dashboard/next-payment': { next_payment: { due_date: '2026-10-30', credit_id: 1, amount_cents: 97375, household_share_cents: 59882, external_share_cents: 37493, days_left: 1, payment_count: 1, is_overdue: false }, availability: [] } })
  await reloj(page, '2026-10-29T12:00:00-06:00')
  await page.goto('/inicio')
  await expect(tarjeta(page, 'K-4')).toContainText('Vie 30 oct 2026 · Kueski · $973.75 · Vence mañana')
})

test('P-01 Detalle: D-3 con tres pagos y sin pedir vencidos (A-16), lecturas A y B, y participación de Karen', async ({ page, context }) => {
  const reg: string[] = []
  await new ApiFalsa().instalar(context); await instalarInicio(context, {}, reg); await reloj(page)
  await page.goto('/inicio?sub=detalle')
  const d3 = page.getByRole('region', { name: 'Pagos y disponible' })
  for (const f of ['Vie 30 oct 2026 · Kueski', 'Dom 15 nov 2026 · Kueski', 'Lun 30 nov 2026 · Kueski', 'Lectura A $10,000.00', 'Lectura B $6,905.50', 'Mínimo del horizonte $8,094.50', 'Mínimo del horizonte $5,000.00']) await expect(d3).toContainText(f)
  await expect(d3.getByText('Pendiente', { exact: true })).toHaveCount(3)
  expect(reg.some((r) => r.includes('status=overdue'))).toBe(false)
  expect(reg.some((r) => r.startsWith('/api/v1/scheduled-payments') && decodeURIComponent(r).includes('period=2026-10-10..2026-12-09'))).toBe(true)
  const d7 = page.getByRole('region', { name: 'Participación de Karen' })
  await expect(d7).toContainText('Su parte de la despensa este mes: $749.86')
  await expect(d7).toContainText('Karen paga directo a Kueski su parte; su parte es informativa y repone la línea cuando paga')
  await expect(d7.getByRole('link', { name: 'Estado de cuenta de Karen' })).toBeVisible()
})

const d3 = (page: Page) => page.getByRole('region', { name: 'Pagos y disponible' })
const PARTES = 'Hogar $598.82 · Karen $374.93'

test('P-01 Detalle: D-3 no pinta ninguna fila mientras listPersons no responde; nunca «Hogar $973.75» sin la parte de Karen (M-1 de REVISION-frontend-cierre-r2, C3-09, R-13)', async ({ page, context }) => {
  const reg: string[] = []
  // Registra si en algún instante la página dijo «Hogar $973.75» (el pago entero cargado al hogar porque aún no se sabe quién es Karen).
  await page.addInitScript(() => {
    const w = window as unknown as { __hogarTodo: boolean }; w.__hogarTodo = false
    new MutationObserver(() => { if (document.body?.textContent?.includes('Hogar $973.75')) w.__hogarTodo = true }).observe(document, { subtree: true, childList: true, characterData: true })
  })
  await new ApiFalsa().instalar(context); await instalarInicio(context, {}, reg); await reloj(page)
  let soltar = () => {}
  const personas = new Promise<void>((r) => { soltar = r })
  await context.route(/\/api\/v1\/persons(\?|$)/, async (route) => { await personas; await route.fallback() })
  await page.goto('/inicio?sub=detalle')
  // Los pagos próximos y el disponible ya respondieron y listPersons sigue retenida: esqueleto, sin filas ni cifras repartidas.
  await expect.poll(() => reg.some((x) => x.startsWith('/api/v1/scheduled-payments') && x.includes('status=pending'))).toBe(true)
  await expect(d3(page)).toContainText('Lectura B $6,905.50')
  await page.waitForTimeout(400)
  await expect(d3(page).getByRole('listitem')).toHaveCount(0)
  await expect(d3(page)).not.toContainText('$973.75')
  await expect(d3(page)).not.toContainText('Hogar')
  soltar()
  // Al soltarla, las tres filas con las dos partes separadas (oráculo R-13).
  for (const f of ['Vie 30 oct 2026 · Kueski', 'Dom 15 nov 2026 · Kueski', 'Lun 30 nov 2026 · Kueski']) await expect(d3(page)).toContainText(f)
  await expect(d3(page).getByText(PARTES, { exact: true })).toHaveCount(3)
  expect(await page.evaluate(() => (window as unknown as { __hogarTodo: boolean }).__hogarTodo)).toBe(false)
})

test('P-01 Detalle: D-3 con listPersons en falla muestra su estado de error con «Reintentar» y sin cifras repartidas; «Reintentar» lo recupera (PROVISIONAL, bloqueo U3-tablero-k4-personas)', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await instalarInicio(context, { '/api/v1/persons': 500 }); await reloj(page)
  await page.goto('/inicio?sub=detalle')
  await expect(d3(page)).toContainText('No se pudieron cargar los pagos. Revisa que estés en la red de casa.')
  await expect(d3(page)).toContainText('Lectura B $6,905.50') // el disponible no depende de listPersons
  await expect(d3(page)).not.toContainText('$973.75')
  await expect(d3(page)).not.toContainText('Hogar')
  await context.unrouteAll({ behavior: 'ignoreErrors' })
  await new ApiFalsa().instalar(context); await instalarInicio(context)
  await d3(page).getByRole('button', { name: 'Reintentar' }).click()
  await expect(d3(page)).toContainText('Vie 30 oct 2026 · Kueski')
  await expect(d3(page).getByText(PARTES, { exact: true })).toHaveCount(3)
})

test('P-01 Detalle: en D-3 una parte de una persona que no está en listPersons no se suma al hogar (M-1 de REVISION-frontend-cierre-r2)', async ({ page, context }) => {
  const parte = (person_id: number, amount_cents: number) => ({ person_id, amount_cents, paid_cents: 0, status: 'pending' })
  const pago = { id: 1, credit_purchase_id: 1, seq: 1, due_date: '2026-10-30', amount_cents: 97375, principal_cents: 90000, interest_cents: 7375, created_at: '2026-10-01T09:00:00-06:00', updated_at: '2026-10-01T09:00:00-06:00', credit_id: 1, status: 'pending', shares: [parte(1, 59882), parte(9, 37493)] }
  await new ApiFalsa().instalar(context); await instalarInicio(context, { '/api/v1/scheduled-payments': { items: [pago], next_cursor: null } }); await reloj(page)
  await page.goto('/inicio?sub=detalle')
  await expect(d3(page)).toContainText('Vie 30 oct 2026 · Kueski')
  await expect(d3(page).getByRole('listitem')).toHaveCount(1)
  await expect(d3(page)).toContainText('Pendiente')
  await expect(d3(page)).not.toContainText('Hogar')
})

test('P-01 Tendencias: T-2 con R-45 y R-46 en la tabla, figuras con figcaption y la pestaña recordada', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await instalarInicio(context); await reloj(page)
  await page.goto('/inicio?sub=tendencias')
  const t2 = page.getByTestId('grafica-t2')
  await expect(t2.locator('figcaption')).toContainText('Devengado contra efectivo')
  await expect(t2.locator('figcaption')).toContainText('devengado acumulado $3,895.00 y efectivo acumulado $2,921.25')
  await t2.getByRole('button', { name: 'Ver como tabla' }).click()
  await expect(t2.getByRole('row', { name: /oct-2026/ })).toContainText('$1,947.50$973.75$973.75')
  await expect(t2.getByRole('row', { name: /nov-2026/ })).toContainText('$1,947.50$1,947.50$973.75')
  expect((await t2.innerText()).toLowerCase()).not.toContain('deuda')
  await expect(page.getByTestId('grafica-t3').locator('figcaption')).toContainText('Deuda viva al cierre de oct-2026: $3,895.00')
  await expect(page.getByText('Aún no hay presupuesto capturado')).toBeVisible()
  await expect(page.getByText('Aún no hay semanas con plan de comidas')).toBeVisible()
  await page.goto('/inicio')
  await expect(page).toHaveURL(/sub=tendencias/)
})

test('P-01 Calendario: pagos con su parte, quincenas K8 y vista según el ancho', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await instalarInicio(context); await reloj(page)
  await page.goto('/inicio?sub=calendario')
  if (page.viewportSize()!.width < 768) {
    const ag = page.getByTestId('agenda')
    await expect(ag).toContainText('Vie 30 oct 2026')
    await expect(ag).toContainText('Kueski · $973.75')
    await expect(ag).toContainText('Hogar $598.82 · Karen $374.93')
    await expect(ag).toContainText('2.ª quincena · del 16 al 31')
  } else {
    const mes = page.getByTestId('mes')
    await expect(mes).toContainText('octubre 2026')
    await expect(mes).toContainText('1.ª quincena · del 1 al 15')
    await mes.getByRole('button', { name: /Vie 30 oct 2026/ }).click()
    await expect(page.getByRole('region', { name: 'Día seleccionado' })).toContainText('Hogar $598.82 · Karen $374.93')
  }
})

test('F-1 desde Inicio: 3 toques y Inicio vuelve a pedir el mes y la quincena; F-3 y F-4 paso 1', async ({ page, context }) => {
  const reg: string[] = []
  await new ApiFalsa().instalar(context); await instalarInicio(context, {}, reg); await reloj(page)
  await page.goto('/inicio')
  await expect(tarjeta(page, 'K-4')).toContainText('$973.75')
  const antes = reg.filter((r) => r.startsWith('/api/v1/dashboard/month-vs-plan')).length
  let toques = 0
  const tocar = async (l: Locator) => { if (page.viewportSize()!.width < 768) await l.tap(); else await l.click(); toques++ }
  await tocar(page.getByRole('button', { name: 'Registrar gasto' }).last())
  await page.getByLabel('Monto').fill('185.50')
  await tocar(page.getByRole('button', { name: 'Despensa' }))
  await tocar(page.getByRole('button', { name: 'Guardar' }))
  console.log(`U5-09 toques desde Inicio (${page.viewportSize()!.width} px): ${toques}`)
  expect(toques).toBeLessThanOrEqual(3)
  await expect.poll(() => reg.filter((r) => r.startsWith('/api/v1/dashboard/month-vs-plan')).length).toBeGreaterThan(antes)
  await expect.poll(() => reg.filter((r) => r.startsWith('/api/v1/dashboard/fortnight-flow')).length).toBeGreaterThan(1)
  await tarjeta(page, 5).click()
  await expect(page).toHaveURL(/\/compras\/comparar$/)
  await page.goBack()
  await tarjeta(page, 'K-4').getByRole('link', { name: /Vie 30 oct 2026/ }).click()
  await expect(page).toHaveURL(/\/creditos\/1$/)
})

## FIN: frontend/pruebas/inicio/inicio.escritorio.spec.ts

## ARCHIVO: frontend/pruebas/inicio/tablero.escritorio.spec.ts


// P-01 Resumen = tablero de control (U3 2.0 §2, K54): corre a 412×915 táctil y a 1440×900, en claro y en oscuro. API simulada.
// Cifras de la semilla y del oráculo (R-13, R-64, M-01; A-05, A-06); el ingreso de las pruebas de barómetro es SINTÉTICO.
import { expect, test, type Page } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'
import { ApiFalsa } from '../apoyo/api-falsa'
import { controlesChicos, dineroSinTabulares, iconosFueraDeNorma, textosEnIngles } from '../apoyo/medidas'
import { instalarInicio, rangoMeses } from './datos'
const reloj = (page: Page, iso = '2026-10-10T12:00:00-06:00') => page.clock.setFixedTime(new Date(iso))
const bloque = (page: Page, id: string) => page.locator(`[data-bloque="${id}"]`)
const ORDEN = ['K-1', 'B-1', 'K-4', 'K-2', 'K-3', 'G-3', 'G-1', 'G-6', 'G-2', 'G-5', 'G-4', 'mas']
/** Capacidad SINTÉTICA (ingreso de prueba, no del oráculo) con el porcentaje y el estado que mandaría el motor. */
const capacidad = (ratio: number, status: string, extra: Record<string, unknown> = {}) => ({ month: '2026-10', income_cents: 2000000, income_is_expected: false, extras_cents: 0, payments_cents: 59882, payments_total_cents: 97375, payment_count: 1, ratio_pct: ratio, ratio_total_pct: ratio + 2, warning_pct: 30, limit_pct: 40, scale_max_pct: 60, status, headroom_cents: 740118, has_credits: true, ...extra })

for (const tema of ['light', 'dark'] as const) {
  test(`P-01 tablero (${tema}): bloques en orden, cifras de la semilla y una petición por bloque (A-01, A-02, A-05, A-06, A-09, A-14)`, async ({ page, context }) => {
    const reg: string[] = []
    await page.emulateMedia({ colorScheme: tema })
    await new ApiFalsa().instalar(context); await instalarInicio(context, {}, reg); await reloj(page)
    await page.goto('/inicio')
    await expect(page.getByRole('heading', { level: 1, name: 'Inicio' })).toBeVisible()
    expect(await page.evaluate(() => matchMedia('(prefers-color-scheme: dark)').matches)).toBe(tema === 'dark')
    // A-01: orden del DOM = orden de apilado de 360/412 en todos los anchos.
    await expect(page.locator('[data-bloque]')).toHaveCount(ORDEN.length)
    expect(await page.locator('[data-bloque]').evaluateAll((xs) => xs.map((x) => x.getAttribute('data-bloque')))).toEqual(ORDEN)
    // A-06: K-1 en octubre de 2026 (P1): gastado $1,197.64, efectivo $598.82, «Sin datos de septiembre»; sin plan, «Capturar presupuesto».
    const k1 = bloque(page, 'K-1')
    for (const x of ['Gasto del mes', '$1,197.64', 'gastado en octubre', '— · Aún no hay presupuesto para octubre', 'Llevas $1,197.64 gastado', 'Pagado en efectivo $598.82', 'Sin datos de septiembre', 'Capturar presupuesto']) await expect(k1).toContainText(x)
    // A-05: K-3 y K-4 con la semilla el 10-oct-2026.
    const k3 = bloque(page, 'K-3')
    for (const x of ['$3,895.00', 'deuda viva total', 'Pagos de octubre: $973.75 · 1 pago', 'Hogar $598.82 · Karen $374.93', 'Falta pagar $598.82', 'Disponible hoy $6,905.50 · Lectura B', 'Pendiente confirmar']) await expect(k3).toContainText(x)
    const k4 = bloque(page, 'K-4')
    await expect(k4.getByRole('listitem')).toHaveCount(1)
    await expect(k4).toContainText('Vie 30 oct 2026 · Kueski · $973.75 · Pendiente')
    await expect(k4).toContainText('Hogar $598.82 · Karen $374.93')
    await expect(k4.getByRole('link', { name: 'Ver calendario' })).toHaveAttribute('href', /\/inicio\?sub=calendario$/)
    // Sin ingreso en la semilla: K-2, B-1 y G-6 con su vacío; G-3 sin plan.
    await expect(bloque(page, 'K-2')).toContainText('Falta tu ingreso quincenal')
    await expect(bloque(page, 'K-2')).toContainText('Capturar ingreso')
    const b1 = bloque(page, 'B-1')
    await expect(b1.getByTestId('barometro-digito')).toHaveText('—')
    await expect(b1).toContainText('Falta tu ingreso del mes para calcular la capacidad')
    // Sin dato no hay `meter` (sin `aria-valuenow` sería inválido, axe `aria-required-attr`): imagen con el texto accesible (DESIGN §8.9).
    await expect(b1.getByRole('meter')).toHaveCount(0)
    const sinDato = b1.getByRole('img', { name: 'Capacidad de endeudamiento sin dato, tope 40 %', exact: true })
    await expect(sinDato).toHaveAttribute('data-prueba', 'barometro')
    for (const a of ['aria-valuenow', 'aria-valuemin', 'aria-valuemax', 'aria-valuetext']) expect(await sinDato.getAttribute(a)).toBeNull()
    await expect(b1.getByTestId('barometro-aguja')).toHaveCount(0)
    await expect(bloque(page, 'G-6')).toContainText('Falta tu ingreso para calcular la capacidad')
    await expect(bloque(page, 'G-3')).toContainText('Aún no hay presupuesto por categoría')
    // G-1 con el gasto de octubre (periodos[0].devengado.axel) y la nota fija; G-2 con la deuda viva por línea.
    await expect(page.getByTestId('grafica-g1').locator('figcaption')).toContainText('En oct-2026: ingreso —, gasto $1,197.64 y pagos de crédito $598.82')
    await expect(page.getByTestId('grafica-g1')).toContainText('Los pagos de crédito pagan compras que ya cuentan en el gasto de su mes; no se suman al gasto')
    await expect(page.getByTestId('grafica-g2').locator('figcaption')).toContainText('Deuda viva hoy: $3,895.00; la mayor, Kueski con $3,895.00')
    await expect(page.getByTestId('grafica-g4').locator('figcaption')).toContainText('Devengado contra efectivo')
    await expect(page.getByTestId('grafica-g5').locator('figcaption')).toContainText('Deuda viva proyectada')
    // «Más del hogar» conserva las tarjetas 2, 4, 5 y 6 (§2.5).
    const mas = bloque(page, 'mas')
    await expect(mas.getByRole('heading', { level: 2, name: 'Más del hogar' })).toBeVisible()
    expect(await mas.locator('[data-prueba^="tarjeta-"]').evaluateAll((xs) => xs.map((x) => x.getAttribute('data-prueba')))).toEqual(['tarjeta-2', 'tarjeta-4', 'tarjeta-5', 'tarjeta-6'])
    await expect(mas).toContainText('Bodega Aurrera $3,094.50')
    // A-02: una petición por operación; G-1 y G-6 de una sola (from = mes − 11, to = mes + 5); ajustes leídos una vez.
    const de = (p: string) => reg.filter((r) => r.startsWith(p))
    expect(de('/api/v1/dashboard/monthly-overview')).toEqual(['/api/v1/dashboard/monthly-overview?from=2025-11&to=2027-03'])
    expect(de('/api/v1/settings/alerts')).toHaveLength(1)
    expect(de('/api/v1/settings/dashboard')).toHaveLength(1)
    expect(de('/api/v1/dashboard/debt-capacity')).toEqual(['/api/v1/dashboard/debt-capacity?period=2026-10'])
    expect(de('/api/v1/incomes/month')).toEqual(['/api/v1/incomes/month?period=2026-10'])
    expect(de('/api/v1/budget/signals')).toEqual(['/api/v1/budget/signals?period=2026-10'])
    expect(de('/api/v1/scheduled-payments').map(decodeURIComponent)).toEqual(['/api/v1/scheduled-payments?status=pending&sort=due_date&limit=5&period=2026-10-10..2026-11-09'])
    // A-08: K-1, K-2, K-3 y B-1 son un solo enlace cada uno (la tarjeta entera).
    for (const id of ['K-1', 'K-2', 'K-3', 'B-1']) {
      const t = bloque(page, id).locator(`[data-prueba="tarjeta-${id}"]`)
      expect(await t.evaluate((x) => x.tagName)).toBe('A')
      await expect(t.locator('a, button')).toHaveCount(0)
    }
    // A-09, A-10, A-14, U5-06/07: ≥ 44 px, sin desplazamiento horizontal, dígito sin recortar, tabulares, sin inglés ni fechas ISO.
    // 44 × 44 con puntero grueso (412 táctil); con puntero fino la barra lateral usa 32 px (DESIGN §7.9).
    if (await page.evaluate(() => matchMedia('(any-pointer: coarse)').matches)) expect((await controlesChicos(page)).chicos).toEqual([])
    expect((await iconosFueraDeNorma(page)).mal).toEqual([])
    expect(await dineroSinTabulares(page)).toEqual([])
    expect(await textosEnIngles(page)).toEqual([])
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true)
    const texto = await page.locator('main').innerText()
    expect(texto).not.toMatch(/\d{4}-\d{2}-\d{2}/)
    // U5-03/U5-06: axe WCAG 2.0–2.2 A/AA sobre el vacío de la semilla (B-1, K-2 y G-6 sin ingreso), en cada ancho y tema.
    const graves = (await new AxeBuilder({ page: page as never }).withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'wcag22aa']).analyze()).violations
      .filter((v) => v.impact === 'serious' || v.impact === 'critical').map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(' ')).join(' | ')}`)
    expect(graves).toEqual([])
  })
}

test('P-01 B-1: el barómetro pinta lo que manda C1 y ningún umbral es constante (A-18)', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await reloj(page)
  const casos: [number, string, string, string][] = [[30, 'normal', 'Verde · dentro del rango', 'verde'], [31, 'warning', 'Amarillo · cerca del tope', 'amarillo'], [40, 'warning', 'Amarillo · cerca del tope', 'amarillo'], [41, 'exceeded', 'Rojo · arriba del tope', 'rojo']]
  for (const [n, st, insignia, color] of casos) {
    await context.unrouteAll({ behavior: 'ignoreErrors' }); await new ApiFalsa().instalar(context)
    await instalarInicio(context, { '/api/v1/dashboard/debt-capacity': capacidad(n, st) })
    await page.goto('/inicio')
    const b1 = bloque(page, 'B-1')
    await expect(b1.getByTestId('barometro-digito')).toHaveText(`${n} %`)
    await expect(b1).toContainText(insignia)
    const medidor = b1.getByRole('meter', { name: 'Capacidad de endeudamiento' })
    await expect(medidor).toHaveAttribute('aria-valuetext', `Capacidad de endeudamiento ${n} % de 40 %: ${color}`)
    await expect(medidor).toHaveAttribute('aria-valuenow', String(n))
    await expect(b1.getByTestId('barometro-aguja')).toHaveCount(1)
    await expect(b1).toContainText('Pagos $598.82 de un ingreso de $20,000.00')
    await expect(b1).toContainText('Te caben $7,401.18 más al mes antes del rojo')
  }
  // Con warning_pct = 25 en la respuesta, la marca del tramo se mueve (los umbrales llegan de C1, K44).
  await context.unrouteAll({ behavior: 'ignoreErrors' }); await new ApiFalsa().instalar(context)
  await instalarInicio(context, { '/api/v1/dashboard/debt-capacity': capacidad(27, 'warning', { warning_pct: 25, headroom_cents: -1200 }) })
  await page.goto('/inicio')
  await expect(bloque(page, 'B-1').getByTestId('barometro-marca')).toHaveText(['25 %', '40 %'])
  await expect(bloque(page, 'B-1')).toContainText('−$12.00 arriba del tope')
  // Los tramos usan los tokens de rol (sin colores sueltos): éxito, advertencia y error distintos entre sí.
  const trazos = await bloque(page, 'B-1').locator('path[data-tramo]').evaluateAll((ps) => ps.map((p) => getComputedStyle(p).stroke))
  expect(new Set(trazos).size).toBe(3)
})

/** Respuestas SINTÉTICAS (ingreso de prueba) con las 6 gráficas del tablero pintadas: G-1, G-6, G-2, G-3, G-4 y G-5. */
const sinteticas = () => {
  const ratio: Record<string, [number, string]> = { '2026-10': [27, 'normal'], '2026-11': [35, 'warning'], '2026-12': [44, 'exceeded'] }
  return {
    '/api/v1/dashboard/debt-capacity': capacidad(27, 'normal'),
    '/api/v1/dashboard/monthly-overview': (u: URL) => ({ warning_pct: 30, limit_pct: 40, months: rangoMeses(u.searchParams.get('from')!, u.searchParams.get('to')!).map((month) => {
      const r = ratio[month], proy = month > '2026-10'
      return { month, is_projection: proy, income_is_expected: proy, ...(month >= '2026-09' ? { income_total_cents: 2000000, income_fortnights_cents: 2000000, income_extras_cents: 0 } : {}), ...(month === '2026-10' ? { accrued_cents: 119764, cash_cents: 59882 } : {}), credit_payments_cents: month === '2026-10' ? 59882 : 0, credit_payments_total_cents: 0, ...(r ? { debt_ratio_pct: r[0], debt_status: r[1] } : {}) }
    }) }),
    '/api/v1/settings/alerts': { parameters: { 'alerts.payment_due_days_before': { value: [1, 0] }, 'alerts.budget_warning_pct': { value: 80 }, 'alerts.budget_limit_pct': { value: 100 } } },
    '/api/v1/budget/categories': { month: '2026-10', rows: [
      { category_id: 1, planned_cents: 300000, actual_cents: 330000, difference_cents: -30000, used_pct: 110, status: 'exceeded', expense_count: 4 },
      { category_id: 2, planned_cents: 100000, actual_cents: 50000, difference_cents: 50000, used_pct: 50, status: 'normal', expense_count: 1 },
      { category_id: 3, planned_cents: 100000, actual_cents: 85000, difference_cents: 15000, used_pct: 85, status: 'warning', expense_count: 2 },
      { category_id: 4, actual_cents: 12000, expense_count: 1 },
    ] },
    '/api/v1/budget/signals': { month: '2026-10', months_considered: 3, margin_pct: 10, rows: [{ category_id: 1, signal: 'raise', months: [] }, { category_id: 2, signal: 'cut', months: [] }, { category_id: 3, signal: 'none', months: [] }] },
  }
}

test('P-01 tablero con ingreso SINTÉTICO: G-1, G-6 y G-3 con sugerencias y desglose a P-02 (A-19, A-20, A-21 paso 1)', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await reloj(page)
  await instalarInicio(context, sinteticas())
  await page.goto('/inicio')
  const g3 = page.getByTestId('grafica-g3')
  await expect(g3.locator('figcaption')).toContainText('En octubre: 2 categorías arriba de 80 % y 1 rebasada')
  const filas = g3.getByRole('listitem')
  await expect(filas).toHaveCount(3)
  await expect(filas.nth(0)).toContainText('Rebasada')
  await expect(filas.nth(0)).toContainText('Sugerencia: subir presupuesto')
  await expect(filas.nth(1)).toContainText('Vas en 85 %')
  await expect(filas.nth(2)).toContainText('Sugerencia: recortar presupuesto')
  await expect(filas.nth(0).locator('[aria-label$="3 meses seguidos arriba del plan por más de 10 %"]')).toHaveCount(1)
  await expect(g3).toContainText('1 categoría con gasto y sin plan')
  await expect(g3.locator('[data-marca="aviso"]')).toHaveCount(3)
  const g1 = page.getByTestId('grafica-g1')
  await expect(g1.locator('figcaption')).toContainText('En oct-2026: ingreso $20,000.00, gasto $1,197.64 y pagos de crédito $598.82')
  await g1.getByRole('button', { name: 'Ver como tabla' }).click()
  await expect(g1.getByRole('row', { name: /oct-2026/ })).toContainText('$20,000.00')
  await expect(g1.getByRole('row', { name: /ago-2026/ })).toContainText('—')
  const g6 = page.getByTestId('grafica-g6')
  await expect(g6.locator('figcaption')).toContainText('oct-2026: 27 % (verde); el mayor de la proyección es dic-2026 con 44 %')
  await expect(g6).toContainText('Proyección con pagos programados e ingreso esperado')
  for (const x of ['30 %', '40 %', 'Hoy']) await expect(g6.locator('svg').getByText(x, { exact: true }).first()).toBeAttached()
  await g6.getByRole('button', { name: 'Ver como tabla' }).click()
  await expect(g6.getByRole('row', { name: /nov-2026/ })).toContainText('35 %amarillo')
  await expect(g6.getByRole('row', { name: /jul-2026/ })).toContainText('— · Falta ingreso')
  expect(await textosEnIngles(page)).toEqual([])
  // Desglose §2.2 paso 1: la barra de la categoría abre P-02 con la categoría y el mes.
  await filas.nth(0).getByRole('link').click()
  await expect(page).toHaveURL(/\/gastos\?categoria=1&periodo=2026-10$/)
})

test('P-01 estados del tablero: un bloque con falla no tumba los demás y «Reintentar» lo recupera (A-07)', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await instalarInicio(context, { '/api/v1/dashboard/debt-capacity': 500, '/api/v1/settings/dashboard': 500, '/api/v1/dashboard/monthly-overview': 500 }); await reloj(page)
  await page.goto('/inicio')
  await expect(bloque(page, 'B-1')).toContainText('No se pudo calcular la capacidad de endeudamiento. Revisa que estés en la red de casa.')
  await expect(bloque(page, 'B-1').getByRole('link', { name: 'Ver créditos' })).toBeVisible()
  await expect(bloque(page, 'K-4')).toContainText('No se pudieron cargar los pagos. Revisa que estés en la red de casa.')
  await expect(bloque(page, 'G-1')).toContainText('No se pudo cargar la gráfica.')
  await expect(bloque(page, 'G-6')).toContainText('No se pudo cargar la gráfica.')
  await expect(bloque(page, 'K-1')).toContainText('$1,197.64')
  await expect(bloque(page, 'K-3')).toContainText('Pagos de octubre: $973.75 · 1 pago')
  await context.unrouteAll({ behavior: 'ignoreErrors' })
  await new ApiFalsa().instalar(context); await instalarInicio(context)
  await bloque(page, 'B-1').getByRole('button', { name: 'Reintentar' }).click()
  await expect(bloque(page, 'B-1')).toContainText('Falta tu ingreso del mes para calcular la capacidad')
  await bloque(page, 'K-4').getByRole('button', { name: 'Reintentar' }).click()
  await expect(bloque(page, 'K-4')).toContainText('Vie 30 oct 2026 · Kueski · $973.75 · Pendiente')
})

test('P-01 K-4 no pinta ninguna fila mientras listPersons no responde: nunca «Hogar $973.75» sin la parte de Karen (M-2 de REVISION-frontend-cierre-r1, C3-09, R-13)', async ({ page, context }) => {
  const reg: string[] = []
  // Registra si en algún instante la página dijo «Hogar $973.75» (el pago entero cargado al hogar porque aún no se sabe quién es Karen).
  await page.addInitScript(() => {
    const w = window as unknown as { __hogarTodo: boolean }; w.__hogarTodo = false
    new MutationObserver(() => { if (document.body?.textContent?.includes('Hogar $973.75')) w.__hogarTodo = true }).observe(document, { subtree: true, childList: true, characterData: true })
  })
  await new ApiFalsa().instalar(context); await instalarInicio(context, {}, reg); await reloj(page)
  let soltar = () => {}
  const personas = new Promise<void>((r) => { soltar = r })
  await context.route(/\/api\/v1\/persons(\?|$)/, async (route) => { await personas; await route.fallback() })
  await page.goto('/inicio')
  const k4 = bloque(page, 'K-4')
  // Los pagos próximos ya respondieron y listPersons sigue retenida: esqueleto, sin filas ni cifras repartidas.
  await expect.poll(() => reg.some((x) => x.startsWith('/api/v1/scheduled-payments') && x.includes('status=pending'))).toBe(true)
  await page.waitForTimeout(400)
  await expect(k4.getByRole('listitem')).toHaveCount(0)
  await expect(k4).not.toContainText('$973.75')
  await expect(k4).not.toContainText('Hogar')
  soltar()
  await expect(k4).toContainText('Vie 30 oct 2026 · Kueski · $973.75 · Pendiente')
  await expect(k4).toContainText('Hogar $598.82 · Karen $374.93')
  expect(await page.evaluate(() => (window as unknown as { __hogarTodo: boolean }).__hogarTodo)).toBe(false)
})

test('P-01 K-4 con listPersons en falla: estado de error con «Reintentar» y sin cifras repartidas; «Reintentar» lo recupera (PROVISIONAL, bloqueo U3-tablero-k4-personas)', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await instalarInicio(context, { '/api/v1/persons': 500 }); await reloj(page)
  await page.goto('/inicio')
  const k4 = bloque(page, 'K-4')
  await expect(k4).toContainText('No se pudieron cargar los pagos. Revisa que estés en la red de casa.')
  await expect(k4).not.toContainText('$973.75')
  await expect(bloque(page, 'K-1')).toContainText('$1,197.64') // los demás bloques siguen
  await context.unrouteAll({ behavior: 'ignoreErrors' })
  await new ApiFalsa().instalar(context); await instalarInicio(context)
  await k4.getByRole('button', { name: 'Reintentar' }).click()
  await expect(k4).toContainText('Vie 30 oct 2026 · Kueski · $973.75 · Pendiente')
  await expect(k4).toContainText('Hogar $598.82 · Karen $374.93')
})

test('P-01 servidor anterior a C1 v1.5 (501 o 404 en las operaciones nuevas): B-1, G-1 y G-6 dicen «Disponible cuando el servidor se actualice», sin cifra ni «Reintentar» (U3 2.0 §2, U3-03)', async ({ page, context }) => {
  await reloj(page)
  for (const estado of [501, 404]) {
    await context.unrouteAll({ behavior: 'ignoreErrors' })
    await new ApiFalsa().instalar(context); await instalarInicio(context, { '/api/v1/dashboard/debt-capacity': estado, '/api/v1/dashboard/monthly-overview': estado })
    await page.goto('/inicio')
    for (const id of ['B-1', 'G-1', 'G-6']) {
      const b = bloque(page, id)
      await expect(b.getByRole('status')).toHaveText('Disponible cuando el servidor se actualice')
      await expect(b.getByRole('alert')).toHaveCount(0)
      await expect(b).not.toContainText('Revisa que estés en la red de casa')
      await expect(b.getByRole('button', { name: 'Reintentar' })).toHaveCount(0)
      expect(await b.innerText()).not.toMatch(/\$|\d\s?%/)
    }
    await expect(bloque(page, 'B-1').getByRole('link', { name: 'Ver créditos' })).toBeVisible()
    // Los bloques de operaciones que ya existen siguen con sus cifras de la semilla.
    await expect(bloque(page, 'K-1')).toContainText('$1,197.64')
    await expect(bloque(page, 'K-3')).toContainText('Pagos de octubre: $973.75 · 1 pago')
  }
})

for (const tema of ['light', 'dark'] as const) {
  test(`P-01 tablero (${tema}): axe WCAG 2.0–2.2 A/AA sin infracciones serias ni críticas con las 6 gráficas, ocultas y como tabla (U5-06, U5-03)`, async ({ page, context }) => {
    await page.emulateMedia({ colorScheme: tema })
    await new ApiFalsa().instalar(context); await reloj(page)
    await instalarInicio(context, sinteticas())
    await page.goto('/inicio')
    for (const g of ['g3', 'g1', 'g6', 'g2', 'g5', 'g4']) await expect(page.getByTestId(`grafica-${g}`)).toBeVisible()
    const graves = async () => (await new AxeBuilder({ page: page as never }).withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'wcag22aa']).analyze()).violations
      .filter((v) => v.impact === 'serious' || v.impact === 'critical').map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(' ')).join(' | ')}`)
    expect(await graves()).toEqual([])
    // La tabla oculta no es región desplazable; ya visible, es región con foco y nombre = título (FiguraBarras y FiguraGrafica).
    for (const [g, titulo] of [['g2', 'Deuda viva por línea'], ['g1', 'Ingreso, gasto y pagos por mes']] as const) {
      const f = page.getByTestId(`grafica-${g}`)
      await expect(f.getByRole('region')).toHaveCount(0)
      await f.getByRole('button', { name: 'Ver como tabla' }).click()
      const region = f.getByRole('region', { name: titulo })
      await expect(region).toHaveAttribute('tabindex', '0')
      await region.focus(); await expect(region).toBeFocused()
    }
    expect(await graves()).toEqual([])
  })
}

test('P-01 K-1 en noviembre de 2026 (P2): «Igual que octubre» en lo gastado y «$598.82 más que octubre» en efectivo (A-06)', async ({ page, context }) => {
  await new ApiFalsa().instalar(context)
  await instalarInicio(context, { '/api/v1/dashboard/month-vs-plan': { month: '2026-11', actual_cents: 119764, days_left: 20, cash_cents: 119764, previous_month: '2026-10', previous_actual_cents: 119764, previous_cash_cents: 59882, actual_difference_cents: 0, cash_difference_cents: 59882, categories_warning_count: 0, categories_exceeded_count: 0 } })
  await reloj(page, '2026-11-10T12:00:00-06:00')
  await page.goto('/inicio')
  const k1 = bloque(page, 'K-1')
  await expect(k1).toContainText('gastado en noviembre')
  await expect(k1).toContainText('Gastado · Igual que octubre')
  await expect(k1).toContainText('Efectivo · $598.82 más que octubre')
})

## FIN: frontend/pruebas/inicio/tablero.escritorio.spec.ts

## ARCHIVO: frontend/pruebas/apoyo/api-falsa.ts


// API falsa para las pruebas base: context.route sobre /api/v1 y /health, con estado en memoria.
// Respeta C1: Page {items,next_cursor}, Error {code,message,request_id}, Idempotency-Key (200 + Idempotent-Replayed).
import type { BrowserContext, Request } from '@playwright/test'
export const SESION = { user_account_id: 1, username: 'axel', person_id: 1, display_name: 'Axel', participation: 'holder', device_session_id: 7, expires_at: '2027-01-10T09:00:00-06:00', last_seen_at: '2026-10-12T09:00:00-06:00' }
const t = '2026-10-01T09:00:00-06:00'
export const CATEGORIAS = ['Despensa', 'Comida fuera', 'Transporte', 'Servicios', 'Salud', 'Hogar', 'Ropa'].map((name, i) => ({ id: i + 1, name, sort_order: i + 1, is_active: true, recent_usage_count: 7 - i, created_at: t, updated_at: t }))
export const PERSONAS = [['axel', 'Axel', 'holder'], ['ximena', 'Ximena', 'co_holder'], ['karen', 'Karen', 'external']].map(([code, display_name, participation], i) => ({ id: i + 1, code, display_name, participation, pays_share_directly: false, is_active: true, created_at: t, updated_at: t }))
export interface Captura { clave: string; cuerpo: Record<string, unknown> }
export class ApiFalsa {
  conSesion = true
  gastos: Captura[] = []
  pesos: Captura[] = []
  peticiones: { metodo: string; url: string; encabezados: Record<string, string>; cuerpo: string | null }[] = []
  fallarSiguienteGasto = false
  /** Mientras sea true, cada POST de gasto responde 500 (servidor con falla), sin guardar nada. */
  fallaServidorGastos = false
  limitarLogin = false
  async instalar(ctx: BrowserContext) {
    await ctx.route(/\/(api\/v1|health)(\/|\?|$)/, async (route) => {
      const pet = route.request()
      const u = new URL(pet.url())
      this.peticiones.push({ metodo: pet.method(), url: u.pathname + u.search, encabezados: pet.headers(), cuerpo: pet.postData() })
      const json = (estado: number, cuerpo: unknown, enc: Record<string, string> = {}) => route.fulfill({ status: estado, contentType: 'application/json', headers: { 'X-Request-Id': 'prueba', ...enc }, body: JSON.stringify(cuerpo) })
      const error = (estado: number, code: string, message: string, extra = {}) => json(estado, { code, message, request_id: 'prueba', ...extra })
      const ruta = `${pet.method()} ${u.pathname}`
      if (ruta === 'GET /health') return json(200, { status: 'ok' })
      if (ruta === 'POST /api/v1/auth/login') {
        if (this.limitarLogin) return error(429, 'too_many_attempts', 'Demasiados intentos.', { retry_at: '2026-10-12T21:45:00-06:00' })
        const b = pet.postDataJSON()
        if (b.username === 'axel' && b.password === 'contrasena-de-prueba') { this.conSesion = true; return json(200, SESION) }
        return error(422, 'invalid_credentials', 'Usuario o contraseña incorrectos.')
      }
      if (!this.conSesion) return error(401, 'unauthorized', 'Inicia sesión.')
      if (ruta === 'GET /api/v1/auth/me') return json(200, SESION)
      if (ruta === 'POST /api/v1/auth/logout') { this.conSesion = false; return route.fulfill({ status: 204 }) }
      if (ruta === 'GET /api/v1/categories') return json(200, { items: CATEGORIAS, next_cursor: null })
      if (ruta === 'GET /api/v1/persons') return json(200, { items: PERSONAS, next_cursor: null })
      if (ruta === 'POST /api/v1/expenses' && this.fallaServidorGastos) return error(500, 'internal_error', 'Ocurrió un error en el servidor de casa.')
      if (ruta === 'POST /api/v1/expenses') return this.capturar(pet, this.gastos, json, error, () => this.fallarSiguienteGasto && !(this.fallarSiguienteGasto = false))
      if (ruta === 'POST /api/v1/weight-entries') return this.capturar(pet, this.pesos, json, error, () => false)
      if (pet.method() === 'GET') return json(200, { items: [], next_cursor: null })
      return error(501, 'not_implemented', 'No implementado en la API falsa.')
    })
  }
  private capturar(pet: Request, lista: Captura[], json: (e: number, c: unknown, h?: Record<string, string>) => Promise<void>, error: (e: number, c: string, m: string) => Promise<void>, fallar: () => boolean) {
    const clave = pet.headers()['idempotency-key'] ?? ''
    if (clave.length < 16 || clave.length > 128 || !/^[A-Za-z0-9_-]+$/.test(clave)) return error(400, 'bad_request', 'Falta la clave de idempotencia.')
    if (fallar()) return error(422, 'validation_failed', 'La categoría ya no existe.')
    const cuerpo = pet.postDataJSON()
    const previo = lista.findIndex((c) => c.clave === clave)
    if (previo >= 0) return json(200, { id: previo + 1, ...lista[previo].cuerpo, created_at: t, updated_at: t }, { 'Idempotent-Replayed': 'true' })
    lista.push({ clave, cuerpo })
    return json(201, { id: lista.length, ...cuerpo, created_at: t, updated_at: t })
  }
}

## FIN: frontend/pruebas/apoyo/api-falsa.ts

## ARCHIVO: frontend/pruebas/apoyo/banco.ts


// Vite de desarrollo para el banco de pruebas de la base (pruebas/banco/): monta las piezas de src/ sin pasar por
// dist/ ni por una pantalla de área. Sin red: solo 127.0.0.1 y los paquetes ya instalados.
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { createServer, type ViteDevServer } from 'vite'
const RAIZ = fileURLToPath(new URL('../..', import.meta.url))
export interface Banco { url: (caso: string) => string; cerrar: () => Promise<void> }
export async function iniciarBanco(): Promise<Banco> {
  const srv: ViteDevServer = await createServer({ root: RAIZ, configFile: path.join(RAIZ, 'vite.config.ts'), logLevel: 'error', clearScreen: false, server: { host: '127.0.0.1', port: 5181, strictPort: false, hmr: false } })
  await srv.listen()
  const base = srv.resolvedUrls!.local[0].replace(/\/$/, '')
  return { url: (caso) => `${base}/pruebas/banco/index.html?caso=${caso}`, cerrar: () => srv.close() }
}

## FIN: frontend/pruebas/apoyo/banco.ts

## ARCHIVO: frontend/pruebas/apoyo/medidas.ts


// Mediciones de DESIGN en el DOM real (U5-06): área de toque ≥ 44 × 44 px (§4.6 size.touch.min), iconos Lucide de
// 20 o 24 px con trazo de 1.5 px dibujados, 24 px en los botones de solo icono (§5.2) y cifras tabulares en todo
// dinero visible (§3.4).
import type { Page } from '@playwright/test'
export interface Control { control: string; ancho: number; alto: number }
export interface IconoMal { clase: string; ancho: number; alto: number; trazo: number; control: string; regla: string }
/** Controles visibles cuya área de toque es menor que 44 × 44 px. La etiqueta que envuelve una casilla cuenta como su
 *  área de toque (tocarla marca la casilla). Se omiten los ocultos (sr-only, fuera de pantalla, aria-hidden, inert). */
export function controlesChicos(page: Page): Promise<{ total: number; chicos: Control[] }> {
  return page.evaluate(() => {
    const SEL = 'a[href], button, input:not([type=hidden]), select, textarea, summary, [role=button], [role=tab], [role=checkbox], [role=switch], [role=radio], [role=menuitem], [role=option], [role=slider]'
    const chicos: { control: string; ancho: number; alto: number }[] = []
    let total = 0
    for (const el of Array.from(document.querySelectorAll<HTMLElement>(SEL))) {
      const r = el.getBoundingClientRect(), cs = getComputedStyle(el)
      if (r.width <= 1 || r.height <= 1 || cs.visibility === 'hidden' || el.closest('[aria-hidden="true"], [inert]')) continue
      if (r.right <= 0 || r.left >= window.innerWidth) continue
      if (cs.clipPath === 'inset(50%)' || cs.clip === 'rect(0px, 0px, 0px, 0px)') continue
      const etiqueta = (el as HTMLInputElement).labels?.[0]
      const caja = etiqueta && etiqueta.contains(el) ? etiqueta.getBoundingClientRect() : r
      total++
      if (caja.width + 0.5 < 44 || caja.height + 0.5 < 44) {
        const nombre = el.getAttribute('aria-label') || el.textContent?.trim().slice(0, 40) || el.getAttribute('name') || ''
        chicos.push({ control: `${el.tagName.toLowerCase()}${el.getAttribute('role') ? `[${el.getAttribute('role')}]` : ''} «${nombre}»`, ancho: Math.round(caja.width), alto: Math.round(caja.height) })
      }
    }
    return { total, chicos }
  })
}
/** Iconos Lucide visibles fuera de §5.2:
 *  - tamaño distinto de 20 o 24 px, o trazo dibujado distinto de 1.5 px (stroke-width en unidades del viewBox de 24 ×
 *    ancho / 24);
 *  - regla de uso: el icono de un botón o enlace sin texto visible (solo icono, con su nombre en aria-label o en texto
 *    sr-only) mide 24 px (`size.icon.md`: «solos en barras… botones de solo icono»), y todo icono de la barra inferior
 *    G-05 («Navegación inferior») mide 24 px aunque lleve nombre debajo (§7.10 `size.icon.md`). Queda fuera la barra lateral
 *    («Navegación principal»), a la que §5.2 asigna `size.icon.sm` en todas sus filas, incluidas las flechas
 *    Desplegar/Plegar de sección. */
export function iconosFueraDeNorma(page: Page): Promise<{ total: number; mal: IconoMal[] }> {
  return page.evaluate(() => {
    const CONTROL = 'button, a[href], [role=button], [role=link], [role=tab], [role=menuitem], [role=option]'
    const oculto = (el: Element) => {
      const cs = getComputedStyle(el), r = el.getBoundingClientRect()
      return cs.display === 'none' || cs.visibility === 'hidden' || cs.clipPath === 'inset(50%)' || cs.clip === 'rect(0px, 0px, 0px, 0px)' || r.width <= 1 || r.height <= 1
    }
    /** ¿El control muestra texto? Cuenta solo el texto que se ve (no el sr-only, ni el oculto, ni el de los svg). */
    const conTexto = (c: Element) => {
      const w = document.createTreeWalker(c, NodeFilter.SHOW_TEXT)
      for (let n = w.nextNode(); n; n = w.nextNode()) {
        if (!n.textContent?.trim()) continue
        let visible = true
        for (let p = n.parentElement; p && visible; p = p === c ? null : p.parentElement) if (p.closest('svg') || oculto(p)) visible = false
        if (visible) return true
      }
      return false
    }
    const nombre = (c: Element) => c.getAttribute('aria-label') || c.getAttribute('title') || c.textContent?.trim().slice(0, 40) || ''
    const mal: { clase: string; ancho: number; alto: number; trazo: number; control: string; regla: string }[] = []
    let total = 0
    for (const svg of Array.from(document.querySelectorAll<SVGSVGElement>('svg.lucide'))) {
      const r = svg.getBoundingClientRect()
      if (r.width <= 1 || r.height <= 1 || svg.closest('[inert]')) continue
      total++
      const vb = svg.viewBox.baseVal?.width || 24
      const dibujado = (Number(svg.getAttribute('stroke-width')) * r.width) / vb
      const ancho = Math.round(r.width * 10) / 10, alto = Math.round(r.height * 10) / 10
      const control = svg.closest(CONTROL)
      const fila = { clase: svg.getAttribute('class') ?? '', ancho, alto, trazo: Math.round(dibujado * 100) / 100, control: control ? `${control.tagName.toLowerCase()} «${nombre(control)}»` : '' }
      if (!((ancho === 20 || ancho === 24) && alto === ancho && Math.abs(dibujado - 1.5) < 0.01)) mal.push({ ...fila, regla: 'tamaño o trazo' })
      else if (control && ancho !== 24 && !control.closest('nav[aria-label="Navegación principal"]') && !conTexto(control)) mal.push({ ...fila, regla: 'solo icono → 24 px' })
      else if (control && ancho !== 24 && control.closest('nav[aria-label="Navegación inferior"]')) mal.push({ ...fila, regla: 'barra inferior → 24 px' })
    }
    return { total, mal }
  })
}
/** Textos visibles con dinero ($1,234.56, y el compacto $5 mil de los ejes) cuyo elemento no tiene
 *  font-variant-numeric: tabular-nums. Incluye el texto SVG de las gráficas de ECharts (ejes, etiquetas directas y
 *  leyenda: DESIGN §7.7 y U5-06), aunque su contenedor sea aria-hidden porque la tabla accesible va aparte; los
 *  demás SVG (iconos) y lo aria-hidden fuera de las gráficas se omiten. Las del SVG se marcan «gráfica: …». */
export function dineroSinTabulares(page: Page): Promise<string[]> {
  return page.evaluate(() => {
    const fuera: string[] = []
    const w = document.createTreeWalker(document.body, NodeFilter.SHOW_TEXT)
    for (let n = w.nextNode(); n; n = w.nextNode()) {
      const txt = n.textContent ?? ''
      if (!/\$\d/.test(txt)) continue
      const el = n.parentElement
      if (!el) continue
      const enGrafica = !!el.closest('svg') && !!el.closest('[_echarts_instance_]')
      if (!enGrafica && el.closest('svg, canvas, option, script, style, [aria-hidden="true"]')) continue
      if (el.closest('option, script, style')) continue
      const r = el.getBoundingClientRect()
      if (r.width <= 1 || r.height <= 1) continue
      if (!getComputedStyle(el).fontVariantNumeric.includes('tabular-nums')) fuera.push(`${enGrafica ? 'gráfica: ' : ''}${txt.trim().slice(0, 80)}`)
    }
    return fuera
  })
}
/** U5-07: inglés visible o anunciado (texto, aria-label, title, alt). Incluye la descripción automática de ECharts
 *  («This is a chart… series… the data…») y palabras de interfaz en inglés. Devuelve los textos que lo contienen. */
export const INGLES = /\b(This is a chart|It consists of|series count|The data is as follows|the data for|representing|Save|Cancel|Settings|Dashboard|Loading|Home|Delete|Search|Add|Submit|Close)\b/
export function textosEnIngles(page: Page): Promise<string[]> {
  return page.evaluate((patron) => {
    const re = new RegExp(patron)
    const fuera: string[] = []
    const w = document.createTreeWalker(document.body, NodeFilter.SHOW_TEXT)
    for (let n = w.nextNode(); n; n = w.nextNode()) {
      const el = n.parentElement
      if (!el || el.closest('script, style, [aria-hidden="true"]')) continue
      const t = n.textContent ?? ''
      if (re.test(t)) fuera.push(`texto «${t.trim().slice(0, 80)}»`)
    }
    for (const el of Array.from(document.querySelectorAll('[aria-label], [title], [alt], [aria-description]'))) {
      for (const at of ['aria-label', 'title', 'alt', 'aria-description']) {
        const v = el.getAttribute(at)
        if (v && re.test(v)) fuera.push(`${at} «${v.slice(0, 120)}»`)
      }
    }
    return fuera
  }, INGLES.source)
}

## FIN: frontend/pruebas/apoyo/medidas.ts

## ARCHIVO: frontend/pruebas-integracion/humo-inicio.spec.ts


// Humo real (C3-09 parcial, S6-12): entra con la cuenta de desarrollo y abre Inicio contra el servidor de casa
// en marcha (ops/dev-run.sh) con la semilla. Sin context.route: todas las cifras vienen de la API real.
import { readFileSync } from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import { expect, test } from '@playwright/test'
import { dinero } from '../src/api/formatos'

// La contraseña se lee del archivo de semilla de desarrollo (modo 600, fuera del repositorio); nunca se imprime.
function contrasenaDesarrollo(usuario: string): string {
  const archivo = process.env.FH_DEV_SEMILLA ?? path.join(os.homedir(), '.config', 'finanzas-hogar', 'dev-semilla.txt')
  for (const renglon of readFileSync(archivo, 'utf8').split('\n')) {
    const i = renglon.indexOf('=')
    if (!renglon.startsWith('#') && i > 0 && renglon.slice(0, i).trim() === usuario) return renglon.slice(i + 1).trim()
  }
  throw new Error(`No hay contraseña de desarrollo para «${usuario}» en el archivo de semilla`)
}

// Las 6 tarjetas de U3 en orden fijo y la operación de C1 de cada una.
const TARJETAS = [
  { n: 1, ruta: '/api/v1/dashboard/month-vs-plan', op: 'getDashboardMonthVsPlan' },
  { n: 2, ruta: '/api/v1/dashboard/fortnight-flow', op: 'getDashboardFortnightFlow' },
  { n: 3, ruta: '/api/v1/dashboard/next-payment', op: 'getDashboardNextPayment' },
  { n: 4, ruta: '/api/v1/dashboard/week-meals', op: 'getDashboardWeekMeals' },
  { n: 5, ruta: '/api/v1/dashboard/best-basket', op: 'getDashboardBestBasket' },
  { n: 6, ruta: '/api/v1/dashboard/savings', op: 'getDashboardSavings' },
]

// Todas las cifras en centavos (*_cents) de una respuesta, ya con el formato de pantalla (con y sin signo de gasto).
function cifras(o: unknown, salida = new Set<string>(), clave = ''): Set<string> {
  if (Array.isArray(o)) o.forEach((v) => cifras(v, salida, clave))
  else if (o && typeof o === 'object') for (const [k, v] of Object.entries(o)) cifras(v, salida, k)
  else if (typeof o === 'number' && Number.isSafeInteger(o) && clave.endsWith('_cents') && o !== 0) {
    salida.add(dinero(o)); salida.add(dinero(Math.abs(o)))
  }
  return salida
}

test('humo real · entrar como axel, Inicio con las 6 tarjetas desde la API y el disponible A $10,000.00', async ({ page }) => {
  const base = new URL(test.info().project.use.baseURL ?? '')
  const respuestas = new Map<string, unknown>()
  const ajenas = new Set<string>()
  page.on('request', (r) => { const u = new URL(r.url()); if (/^https?:$/.test(u.protocol) && u.host !== base.host) ajenas.add(u.host) })
  page.on('response', async (r) => {
    const u = new URL(r.url())
    const t = TARJETAS.find((x) => x.ruta === u.pathname)
    if (t && r.status() === 200 && u.host === base.host) respuestas.set(t.ruta, await r.json().catch(() => undefined))
  })

  await page.goto('/inicio')
  await expect(page.getByRole('heading', { name: 'Iniciar sesión' })).toBeVisible()
  await page.getByLabel('Usuario').fill('axel')
  await page.getByLabel('Contraseña').fill(contrasenaDesarrollo('axel'))
  await page.getByRole('button', { name: 'Entrar' }).click()
  await expect(page).toHaveURL(/\/inicio/)
  await expect(page.getByRole('heading', { level: 1, name: 'Inicio' })).toBeVisible()

  for (const t of TARJETAS) {
    const tarjeta = page.getByTestId(`tarjeta-${t.n}`)
    await expect(tarjeta, `tarjeta ${t.n} visible`).toBeVisible()
    await expect.poll(() => respuestas.has(t.ruta), { message: `${t.op} respondió 200 desde ${base.host}` }).toBe(true)
    // Cargada = la tarjeta es el enlace al detalle (no la caja de carga, vacío o error).
    await expect(tarjeta, `tarjeta ${t.n} cargada`).toHaveJSProperty('tagName', 'A')
    const texto = (await tarjeta.innerText()).replace(/\s+/g, ' ')
    const esperadas = [...cifras(respuestas.get(t.ruta))]
    const halladas = esperadas.filter((c) => texto.includes(c))
    test.info().annotations.push({ type: `tarjeta-${t.n} ${t.op}`, description: `cifras de la API: ${esperadas.length}; en pantalla: ${halladas.join(' | ') || '—'} · texto: ${texto.slice(0, 160)}` })
    if (esperadas.length) expect(halladas.length, `tarjeta ${t.n}: alguna cifra de ${t.op} en pantalla`).toBeGreaterThan(0)
  }

  // C3-09 (parcial), cifras de oraculo/model.json: próximo pago 973.75 y disponible de la lectura activa (B 6,905.50).
  await expect(page.getByTestId('tarjeta-3')).toContainText('$973.75')
  await expect(page.getByTestId('tarjeta-3')).toContainText('$6,905.50')

  // Disponible inicial A 10,000.00 (credito.eventos[0].disponible_A): Inicio → Detalle → «Disponible por línea».
  const vista = page.waitForResponse((r) => /\/api\/v1\/credits\/\d+\/availability-preview$/.test(new URL(r.url()).pathname) && r.status() === 200)
  await page.getByRole('tab', { name: 'Detalle' }).click()
  const dispo = (await (await vista).json()) as { reading_a: { available_today_cents: number }; reading_b: { available_today_cents: number } }
  expect(dispo.reading_a.available_today_cents, 'API: disponible lectura A').toBe(1_000_000)
  expect(dispo.reading_b.available_today_cents, 'API: disponible lectura B').toBe(690_550)
  await expect(page.getByText('$10,000.00').first()).toBeVisible()
  await expect(page.getByText(dinero(dispo.reading_b.available_today_cents)).first()).toBeVisible()
  expect([...ajenas], 'peticiones a hosts externos').toEqual([])
})

## FIN: frontend/pruebas-integracion/humo-inicio.spec.ts

## ARCHIVO: frontend/playwright.config.ts


import fs from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import { defineConfig } from '@playwright/test'

// Chromium local sin red: FH_CHROMIUM si viene; si no, el chrome-headless-shell más nuevo que ya esté en la caché de
// Playwright (así `make verify-wf3-frontend` → `npm run verify` corre sin variables extra). Sin ninguno, el de Playwright.
function chromiumLocal(): string | undefined {
  if (process.env.FH_CHROMIUM) return process.env.FH_CHROMIUM
  const base = path.join(os.homedir(), '.cache', 'ms-playwright')
  try {
    const revs = fs.readdirSync(base).filter((n) => /^chromium_headless_shell-\d+$/.test(n)).sort((a, b) => Number(b.split('-')[1]) - Number(a.split('-')[1]))
    for (const r of revs) { const exe = path.join(base, r, 'chrome-headless-shell-linux64', 'chrome-headless-shell'); if (fs.existsSync(exe)) return exe }
  } catch { /* sin caché: el de Playwright */ }
  return undefined
}

// Pruebas base (frontend-base). Corren contra «vite preview» del build real (con service worker);
// FH_CHROMIUM: ruta a un Chromium local si el de Playwright no está descargado (sin red en ejecución).
// el API se simula con context.route (pruebas/apoyo/api-falsa.ts). Sin red externa.
export default defineConfig({
  testDir: 'pruebas',
  outputDir: '../datos/tmp/playwright-frontend',
  workers: 1,
  reporter: [['list']],
  timeout: 45_000,
  use: { baseURL: 'http://localhost:4173', locale: 'es-MX', timezoneId: 'America/Mexico_City', trace: 'off', testIdAttribute: 'data-prueba', launchOptions: { executablePath: chromiumLocal() } },
  projects: [
    { name: 'movil', use: { browserName: 'chromium', viewport: { width: 412, height: 915 }, hasTouch: true, isMobile: true } },
    { name: 'escritorio', use: { browserName: 'chromium', viewport: { width: 1440, height: 900 } }, testMatch: /(login|shell|\.escritorio)\.spec\.ts/ },
  ],
  webServer: { command: 'npm run preview', url: 'http://localhost:4173', reuseExistingServer: true, timeout: 60_000 },
})

## FIN: frontend/playwright.config.ts

## ARCHIVO: docs/estado/bloqueos/WF3c-tablero-interfaz.md


# Bloqueo WF3c-tablero-interfaz · Puntos del tablero de Inicio que U3 2.0 no cierra para la interfaz

Abierto por `front-tablero` (WF-3c, 12-oct-2026). Cada punto está construido como dice la columna «Hecho (PROVISIONAL)» y sigue lo demás.

| # | Qué falta | Hecho (PROVISIONAL) | Opciones | Impacto |
|---|---|---|---|---|
| 1 | K-4: «y [n] más en Calendario» pide un conteo que `listScheduledPayments` no da (solo `next_cursor`) | Con `next_cursor` ≠ null: «Hay más pagos en Calendario», sin número | A) Así. B) C1 agrega un conteo (p. ej. `total_count` en la lista o `upcoming_count` en `CreditTotals`). C) Una tercera petición con `limit` 200 solo para contar | Texto de una línea de K-4 |
| 2 | K-3: «Disponible hoy —» con «Reintentar» si solo falla `getDashboardNextPayment`; un botón dentro de la tarjeta-enlace es HTML inválido y rompe A-08 (un solo enlace por tarjeta) | «Disponible hoy —» sin botón; el bloque se reintenta al recargar Inicio o al volver la conexión | A) Así. B) Poner «Reintentar» fuera del enlace, bajo la tarjeta | Un botón en un caso de falla parcial |
| 3 | K-1: la insignia «[n] categorías rebasadas» «lleva a G-3», pero la tarjeta entera ya es el enlace a P-06 (A-08) | Insignia de texto; G-3 es el bloque 6 del mismo tablero | A) Así. B) Enlace aparte fuera de la tarjeta | Un toque menos hacia G-3 |
| 4 | B-1: DESIGN §8.9 fija el alto del bloque en `size.chart.height-compact` / `size.chart.height`; el arco escala con el ancho | Arco de ancho completo hasta 320 px (alto ≈ 176 px) más dígito, insignia y líneas; la tarjeta crece con su contenido | A) Así. B) Fijar el alto y recortar líneas secundarias | Alto de la tarjeta B-1 |
| 5 | G-1: «2 px de separación» entre columnas; ECharts 6 da `barGap` en % | `barGap` 8 % | A) Así. B) Calcular el % desde el ancho en cada cambio de tamaño | Visual |
| 6 | G-2 y G-3 como barras HTML (no ECharts) para tener filas-enlace de 44 px con icono `triangle-alert` e insignias de sugerencia (A-20, A-21) | `<figure>` con `<figcaption>`, resumen, «Ver como tabla» y la tabla siempre para lectores de pantalla, como las demás | A) Así. B) ECharts con la lista de enlaces aparte | Ninguno en cifras |
| 7 | K-4: N que no sea entero > 0 en `getDashboardSettings` | Se trata como falla de lectura: K-4 muestra su error con «Reintentar» (U3 §1.2) | A) Así. B) Otra conducta | Solo con datos inválidos |

## FIN: docs/estado/bloqueos/WF3c-tablero-interfaz.md

## ARCHIVO: docs/estado/bloqueos/WF3c-barometro-sin-dato.md


# Bloqueo WF3c-barometro-sin-dato · Rol accesible del barómetro cuando no hay dato

Abierto por `correccion-tablero-frontend-r2` (WF-3c, 12-oct-2026), a partir del hallazgo N-1 de `docs/estado/wf3/REVISION-tablero-frontend-r2.md` y de su PROPUESTA de §8. No cambia tokens, cifras ni textos visibles.

**Qué falta.** DESIGN §8.9 (Accesibilidad) pide `role="meter"` con `aria-valuenow` = valor, pero no dice qué hacer en el estado «Sin dato» que la misma §8.9 define (por ejemplo, sin ingreso; es el estado de la semilla de `datos/semilla/`). Un `meter` sin `aria-valuenow` es inválido: axe 4.13 da `aria-required-attr` (impacto crítico) y U5-03 falla.

**Hecho (PROVISIONAL).** `frontend/src/componentes/Barometro.tsx`: con valor, el `meter` de siempre (`aria-valuemin` 0, `aria-valuemax` = escala, `aria-valuenow` = valor, `aria-valuetext` con título, valor, tope y estado). Sin dato, el dibujo es `role="img"` con `aria-label` = el texto accesible que ya armaba la pantalla («Capacidad de endeudamiento sin dato, tope 40 %»), sin atributos `aria-value*`. El dígito «—» y el motivo con su acción siguen igual.

| Opción | Qué es | Impacto |
|---|---|---|
| A (hecha) | Sin dato: `role="img"` con el texto accesible | axe 0 serios en las 4 combinaciones; el lector dice el título, «sin dato» y el tope |
| B | Sin dato: contenedor sin rol y el mismo texto en `sr-only` | Mismo resultado; el texto se lee como párrafo y no como imagen |
| C | Mantener `meter` con `aria-valuenow` = 0 y `aria-valuetext` «sin dato» | Se descarta: anuncia un 0 % que no existe (cifra inventada) |

**Pregunta para quien edita DESIGN (o Axel).** ¿Se agrega a DESIGN §8.9, Accesibilidad, la línea «Sin dato: `role="img"` con el texto accesible, sin `aria-valuenow`» (opción A), o se prefiere B? Mientras no se decida, queda A.

## FIN: docs/estado/bloqueos/WF3c-barometro-sin-dato.md

## ARCHIVO: docs/estado/bloqueos/U3-tablero-k4-personas.md


# Bloqueo U3-tablero-k4-personas · K-4 «Pagos próximos» cuando `listPersons` aún no responde o falla (U3 §2 K-4, C3-09)

Fecha: 2026-10-12 · Agente: `correccion-frontend-cierre-r1` (Opus, WF-3a cierre) · Origen: hallazgo M-2 de `docs/estado/wf3/REVISION-frontend-cierre-r1.md`. Ninguna cifra cambia.

## Qué falta (HECHO)
- K-4 (U3 §2) pinta por fila «Hogar [parte] · Karen [parte]» y el estado del pago de U3 §1.4 («Falta la parte de Karen» / «Falta la parte del hogar»). Las dos cosas necesitan saber quién es Karen: la persona con `participation` externa de `listPersons`.
- Antes, mientras `listPersons` no respondía, la fila decía un instante «Hogar $973.75» (todo al hogar), una cifra distinta del oráculo (R-13: Hogar $598.82 · Karen $374.93).
- U3 dice qué hace K-4 si falla la lectura de N (`getDashboardSettings`, §1.2: estado de error del bloque) pero no dice qué hace si falla `listPersons`.

## Qué se hizo (PROVISIONAL)
- Mientras `listPersons` no llega: esqueleto de filas (no se pinta ninguna fila).
- Si `listPersons` falla o nunca se cargó sin conexión: el estado de error de K-4 («No se pudieron cargar los pagos. Revisa que estés en la red de casa.» · «Reintentar»), igual que cuando falla N; «Reintentar» vuelve a pedir lo que falló.
- Código: `frontend/src/pantallas/inicio/Tablero.tsx` (`K4`); prueba: `frontend/pruebas/inicio/tablero.escritorio.spec.ts`.

## PREGUNTA
¿Qué muestra K-4 si falla `listPersons`?

| Opción | Resultado | Costo | Riesgo |
|---|---|---|---|
| A. PROPUESTA (aplicada): estado de error del bloque con «Reintentar» | Nunca una cifra distinta del oráculo; misma conducta que con N | Ninguno | Bajo: se pierde la lista mientras falla una lectura secundaria |
| B. Las filas sin la línea «Hogar · Karen» y con el estado de pago sin las variantes «Falta la parte de…» | Se ven fecha, línea y monto | Bajo | Medio: un pago con una parte ya pagada diría «Pendiente» en vez de «Falta la parte de Karen» |

Impacto: solo el caso de falla de `listPersons`; con la lectura bien, K-4 no cambia. La misma conducta faltaba en D-3 (`Detalle.tsx`) y en el Calendario (`Calendario.tsx`), que usan `karen(c)` igual. Decide Axel en la puerta 4 o el orquestador.

## Ampliación a D-3 (12-oct-2026, M-1 de `docs/estado/wf3/REVISION-frontend-cierre-r2.md`)
- D-3 «Pagos y disponible» de P-01 Detalle aplica ya la misma opción A: esqueleto mientras `listPersons` no llega; si falla o nunca se cargó, el estado de error de D-3 («No se pudieron cargar los pagos. Revisa que estés en la red de casa.» · «Reintentar», que vuelve a pedir `listPersons`). «Disponible por línea» no depende de `listPersons` y se sigue viendo.
- Además, una parte de una persona que no aparece en `listPersons` ya no se suma a «Hogar»: esa fila se queda sin la línea «Hogar · Karen».
- Código: `frontend/src/pantallas/inicio/Detalle.tsx` (`D3`); pruebas: `frontend/pruebas/inicio/inicio.escritorio.spec.ts`; reporte: `docs/estado/wf3/correccion-m1-detalle.md`.
- Sigue pendiente la insignia de `Calendario.tsx` (los montos vienen del servidor; solo cambiaría con un pago a medias).

## FIN: docs/estado/bloqueos/U3-tablero-k4-personas.md

## ARCHIVO: docs/estado/bloqueos/U3-tablero-barometro.md


# Bloqueo U3-tablero-barometro · qué entra en la capacidad de endeudamiento (U3 2.0, B-1 y G-6)

Abierto el 12-oct-2026 por el agente `u3-v2` (WF-3c). Lo contesta Axel en la **puerta 4**. Fuente: decisiones K54 (decisión 1: «pagos mensuales de todos los créditos ÷ ingreso mensual neto»; decisión 2: verde hasta 30 %, amarillo de 30 a 40 %, rojo arriba de 40 %). Lo FIRME de K54 no se reabre: la fórmula general, los umbrales 30 y 40 y que sean configurables (`debt.capacity_warning_pct`, `debt.capacity_limit_pct`). Lo que sigue es lo que Axel no definió. Mientras tanto rige la opción marcada PROVISIONAL, que ya usa `docs/design/U3-dashboard.md` 2.0 (B-1, G-6, §2.6). Ninguna opción cambia cifras canon.

| # | Qué falta | Opciones (la PROVISIONAL va primero) | Impacto |
|---|---|---|---|
| 1 | **Qué es «ingreso mensual neto».** El sistema solo conoce lo que se captura en P-08 y P-26.03 (no calcula impuestos). | (a) PROVISIONAL: las quincenas del mes, capturadas más las esperadas que aún no se capturan, **sin extras**. Un extra de un solo mes pintaría el barómetro de verde ese mes y escondería el riesgo de los meses siguientes; el extra se ve aparte en K-2. (b) Quincenas más extras del mes. (c) Promedio de las quincenas de los últimos 3 meses. | Cambia `income_cents` de `getDashboardDebtCapacity` y `income_total_cents` frente a `income_fortnights_cents` en G-6. |
| 2 | **Qué pagos cuentan.** En Kueski, Karen paga directo su parte (decisiones §A). | (a) PROVISIONAL: solo la **parte del hogar** de los pagos con fecha en el mes, pagados o no; el porcentaje con la parte de Karen (`ratio_total_pct`) solo se informa desde 768 px. (b) El monto total de cada pago, con la parte de Karen. (c) Solo lo que falta pagar del mes. | Cambia `payments_cents` y el estado del barómetro. |
| 3 | **Redondeo y fronteras.** | (a) PROVISIONAL: `ratio_pct` entero con ROUND_HALF_UP (protocolo §6); el estado se decide con ese entero: ≤ 30 verde, 31 a 40 amarillo, ≥ 41 rojo (lectura literal de «verde hasta 30 %» y «rojo arriba de 40 %»). (b) Un decimal (décimas, como `progress_tenths_pct`) y estado con el valor exacto. | Cambia el dígito y el estado en los bordes. |
| 4 | **Escala del arco.** | (a) PROVISIONAL: de 0 % al múltiplo de 10 igual o mayor que 1.5 × el tope (60 % con 40), que da a la zona roja el mismo ancho que la verde y la amarilla juntas hasta 40; un valor mayor deja la aguja al final con la cifra exacta. (b) De 0 a 100 %. (c) Clave nueva de `app_setting`. | Solo dibujo; no cambia cifras. |
| 5 | **Mes que se mide.** | (a) PROVISIONAL: el mes calendario en curso (decisiones §B, periodo base); G-6 da los demás meses. (b) El promedio de los próximos 3 meses. | Cambia el dígito en meses con pagos dobles. |
| 6 | **Número de la sección de DESIGN.** El texto de K54 dice «DESIGN.md §8.5 barómetro», pero §8.5 ya es «Etiquetas de dinero siempre en texto» y la citan muchos documentos. | (a) PROVISIONAL: el barómetro entra como **§8.9** sin renumerar. (b) Renumerar §8.5 a §8.8 y corregir todas las citas. | Solo referencias. |

**Qué cambia al contestar:** `docs/design/U3-dashboard.md` (B-1, G-6, §2.6 y supuestos), `docs/design/DESIGN.md` §8.9 (solo con 4 o 6), C1 v1.5 (`DashboardDebtCapacity` y `DashboardMonthlyOverview`) y el motor que calcule la razón.

## FIN: docs/estado/bloqueos/U3-tablero-barometro.md

## ARCHIVO: docs/estado/bloqueos/U3-tablero-horizonte.md


# Bloqueo U3-tablero-horizonte · cuántos meses ve cada gráfica del tablero (U3 2.0, §2.1)

Abierto el 12-oct-2026 por el agente `u3-v2` (WF-3c). Lo contesta Axel en la **puerta 4**. Fuente: decisiones K54 («graficas de analisis… para irme previniendo»). Axel no fijó el horizonte. Mientras tanto rige la opción PROVISIONAL, que usa `docs/design/U3-dashboard.md` 2.0 (§2.1, G-1, G-6).

| # | Qué falta | Opciones (la PROVISIONAL va primero) | Impacto |
|---|---|---|---|
| 1 | **G-1 · ingreso, gasto y pagos por mes.** | (a) PROVISIONAL: los 11 meses anteriores y el mes en curso (lo real; el gasto futuro no existe). (b) 6 anteriores, el en curso y 5 siguientes, con el plan en lugar del gasto en los meses futuros. | Ventana de `getDashboardMonthlyOverview` y meses con datos. |
| 2 | **G-6 · tendencia del barómetro.** | (a) PROVISIONAL: los 6 meses anteriores, el en curso y los 5 siguientes. Los pagos futuros ya se conocen por el calendario de pagos y las quincenas esperadas por P-26.03, así que se ve antes un mes que se pasa del 40 %. (b) Los 12 meses anteriores. (c) El en curso y los 11 siguientes. | Ventana y número de puntos de proyección. |
| 3 | **Una o varias peticiones.** | (a) PROVISIONAL: una sola petición de 17 meses (`from` = en curso − 11, `to` = en curso + 5) y cada gráfica muestra su tramo. (b) Una petición por gráfica. | Rendimiento; no cambia cifras. |

G-4 (T-2: el mes en curso y los 11 siguientes) y G-5 (T-3: 12 meses) conservan las ventanas que ya aprobó la puerta 2.

**Qué cambia al contestar:** U3 2.0 §2.1, G-1, G-6, §6 y A-19; C1 v1.5 (límite de `from` y `to`).

## FIN: docs/estado/bloqueos/U3-tablero-horizonte.md

## ARCHIVO: docs/estado/bloqueos/U3-tablero-recordatorio.md


# Bloqueo U3-tablero-recordatorio · N y rango de «Pagos próximos» (U3 2.0, K-4)

Abierto el 12-oct-2026 por el agente `u3-v2` (WF-3c). Lo contesta Axel en la **puerta 4**. Fuente: decisiones K54, palabras de Axel: «una vista rapida de recordatorio de pagos del mes, puede ser de pagos N o de creditos proximos». Axel no fijó N ni el rango. Mientras tanto rige la opción PROVISIONAL, que usa `docs/design/U3-dashboard.md` 2.0, K-4.

| # | Qué falta | Opciones (la PROVISIONAL va primero) | Impacto |
|---|---|---|---|
| 1 | **N, cuántos pagos se ven.** | (a) PROVISIONAL: **5**, como D-3 (U1 N-027 «próximos 5 pagos»), en la clave nueva de `app_setting` `dashboard.payment_reminder_count` (valor inicial 5, de 1 a 20, K44). Mientras C1 v1.5 no publique la clave, la interfaz usa 5. (b) 3, para que quepa en un vistazo a 412 px. (c) Todos los del rango. | Alto de K-4 y una petición con `limit` = N. |
| 2 | **Rango.** | (a) PROVISIONAL: primero los vencidos y luego los pendientes de hoy a hoy + 30 días. Cubre «los pagos del mes» y no se queda vacío al final del mes. (b) Solo el mes calendario en curso. (c) De hoy a hoy + 60, como D-3. | Qué pagos se ven; con la semilla el 10-oct-2026, (a) da 1 pago y (c) da 3. |
| 3 | **Pagos o créditos.** «pagos N o créditos próximos». | (a) PROVISIONAL: una fila por pago (fecha, línea, monto y estado) y varios pagos del mismo día en filas separadas. (b) Una fila por línea con su próximo pago. | Forma de la lista. |

**Qué cambia al contestar:** U3 2.0 K-4 y §6, el esquema 1.5 (valor inicial de la clave) y C1 v1.5 (`DashboardSettings`).

## FIN: docs/estado/bloqueos/U3-tablero-recordatorio.md

## ARCHIVO: docs/estado/bloqueos/U3-tablero-senal-presupuesto.md


# Bloqueo U3-tablero-senal-presupuesto · regla «subir o recortar presupuesto» por categoría (U3 2.0, G-3)

Abierto el 12-oct-2026 por el agente `u3-v2` (WF-3c). Lo contesta Axel en la **puerta 4**. Fuente: decisiones K54, propósito del sistema según Axel: «para ver si aumentamos presupuesto o lo recortamos para irme previniendo al respecto». El orquestador dio un ejemplo («tres meses seguidos arriba o abajo del plan por más de un umbral configurable»), pero Axel no fijó la regla. Mientras tanto rige la opción PROVISIONAL. Ya la usan `docs/design/U3-dashboard.md` 2.0 (G-3, §2.6) y la operación propuesta `getBudgetSignals`. La calcula el motor S2 y la interfaz solo la muestra.

| # | Qué falta | Opciones (la PROVISIONAL va primero) | Impacto |
|---|---|---|---|
| 1 | **Cuántos meses.** | (a) PROVISIONAL: **3 meses cerrados** antes del mes pedido (clave nueva `budget.signal_months`, valor inicial 3, de 2 a 12). El mes en curso no cuenta porque no ha cerrado. (b) 2 meses. (c) 6 meses. | Qué tan pronto aparece la sugerencia. |
| 2 | **Umbral.** | (a) PROVISIONAL: un solo margen de **10 %** para los dos lados (clave nueva `budget.signal_margin_pct`, valor inicial 10): «subir» si en los 3 meses el % usado pasó de 110; «recortar» si quedó abajo de 90. (b) Márgenes distintos para subir y recortar (p. ej. 10 y 20). (c) Comparar el promedio de los meses contra el plan, sin racha. | Cuántas categorías llevan sugerencia. |
| 3 | **Meses sin plan o sin gasto.** | (a) PROVISIONAL: un mes sin plan de esa categoría corta la racha (sin datos no es cero); un mes con plan y sin gasto cuenta como 0 % usado. (b) Se saltan sin cortar la racha. | Categorías nuevas o de temporada. |
| 4 | **Qué dice y dónde.** | (a) PROVISIONAL: insignia de rol información «Sugerencia: subir presupuesto» o «Sugerencia: recortar presupuesto» en G-3, sin monto sugerido. (b) Además un monto sugerido (p. ej. el promedio real de los 3 meses). (c) Además un aviso de S7 el día 1. | Texto, y en (b) y (c) una cifra o un aviso nuevos. |

**Qué cambia al contestar:** U3 2.0 G-3 y §2.6; C1 v1.5 (`BudgetSignals`); esquema 1.5 (las dos claves); motor S2 (función nueva de sugerencia, con su prueba).

## FIN: docs/estado/bloqueos/U3-tablero-senal-presupuesto.md

## ARCHIVO: docs/estado/bloqueos/C3-ajustes-tablero.md


# Bloqueo C3-ajustes-tablero · `updateDashboardSettings` no tiene pantalla (C3-02, C3-03)

Fecha: 2026-10-12 · Agente: `c3-tablero` (Sonnet, WF-3c) · Contexto: `docs/contratos/C3-integracion.md` v1.3 (sin firma). Hay una pregunta abierta y nada de esto cambia cifras ni motores.

## Qué falta (HECHO)
- C1 v1.5 trae la pareja `getDashboardSettings` / `updateDashboardSettings` (GET y PATCH `/api/v1/settings/dashboard`; `x-necesidades`: N-278 la de escritura).
- U1 §3.3 (fila «Parámetros del tablero»: «las cinco claves de K54 en Ajustes: leer y cambiar», N-275 y N-278) y U4 la ponen en Ajustes › Presupuesto (P-26.02).
- En `frontend/src` solo Inicio llama `getDashboardSettings` (`pantallas/inicio/comun.tsx`, para N de «Pagos próximos»). **Ninguna pantalla llama `updateDashboardSettings`**; `S02Presupuesto.tsx` no tiene el apartado.
- `python3 -I frontend/scripts/c3-integracion.py` lo informa: «Operaciones con pantalla que el frontend no llama: 2 [listPeriods, updateDashboardSettings]; sin excepción justificada o sin cobertura: 1 [updateDashboardSettings]» y sale con 1.
- La fila P-26.02 sigue en «sí» en implementado porque la regla de C3 es por necesidad (N-278 la cubren `createBudget`, `updateBudget`, `createPlannedItem`, `updatePlannedItem`, `deleteBudget`); pero el apartado «Parámetros del tablero» no existe.

## PREGUNTA
¿Qué se hace con `updateDashboardSettings`?

| Opción | Resultado | Costo | Riesgo |
|---|---|---|---|
| A. PROPUESTA (la recomendada): el agente de Ajustes (frontend, P-26.02) agrega el apartado «Parámetros del tablero» con las cinco claves de K54 (`debt.capacity_warning_pct`, `debt.capacity_limit_pct`, `dashboard.payment_reminder_count`, `budget.signal_months`, `budget.signal_margin_pct`), con «Restablecer», que llama PATCH con `parameters` y `reset_parameters` | C3-02 y C3-03 en verde sin excepción; cierra el ciclo K44 | Medio: una sección más con su prueba y sus capturas de P-26.02 | Bajo |
| B. Declararla excepción técnica PROVISIONAL (como `listPeriods`) | El script sale con 0 | Ninguno | Alto: queda el valor del barómetro (30 % y 40 %) sin forma de cambiarse desde la app, contra K44 y U1 |
| C. Quitar la operación de C1 o dejarla solo para el agente externo | Contrato más chico | Medio: cambia C1 y el generador | Contradice U1 y U3 §2.6 |

Impacto: sin decidir, C3 no se firma con 0 filas abiertas en C3-03 (el script falla por esta operación). **No se tomó ninguna decisión**: el script no la admite como excepción; el pendiente queda a la vista en la tabla al pie de la matriz. Decide Axel o el orquestador.

## FIN: docs/estado/bloqueos/C3-ajustes-tablero.md

## FIN-DEL-CONTEXTO
