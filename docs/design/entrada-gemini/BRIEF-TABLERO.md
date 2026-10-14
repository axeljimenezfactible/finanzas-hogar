# BRIEF para Gemini · Terminar el tablero de Inicio (P-01) · Finanzas Hogar

Versión 1.1 · 13-oct-2026 · lo escribió Claude para Axel. Lo lee Gemini **con el repositorio completo clonado**: aquí solo están el contexto, las reglas y la lista de lo que falta; el detalle vive en los archivos que se citan. Si algo no está en los archivos, **no lo asumas** (regla 9).

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
