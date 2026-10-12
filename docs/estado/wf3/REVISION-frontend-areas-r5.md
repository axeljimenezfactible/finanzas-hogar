# REVISIÓN frontend de áreas · ronda 5 · WF-3a (12-oct-2026)

Revisor: `revisor-frontend-areas-r5`. Soy independiente: no escribí nada de lo que reviso y no corregí nada. Trabajo revisado: front-creditos-salud-reportes, front-ajustes-a, front-ajustes-b y front-variantes, con las correcciones `correccion-frontend-areas-r4` y anteriores. Criterios: U5-06, U5-07 y U5-09 (`docs/criterios/U5.md`), C3-03 (`C3.md`), DESIGN.md, U1 y U4. Árbol revisado: HEAD `a35451e` más los cambios sin commit del área.

Todas las mediciones son mías y las hice en `datos/tmp/rev-fa-r5/`, que ya borré:
- un build propio (`vite build --outDir datos/tmp/rev-fa-r5/dist`), servido con `vite preview` en el puerto 4799;
- una configuración temporal de Playwright que extiende `frontend/playwright.config.ts`;
- tres sondas propias: errores de guardado en Ajustes, un barrido de 32 rutas × 2 temas × 2 proyectos y un 503 en P-20.

No usé `frontend/dist` ni el puerto 4173.

## Veredicto: NO APROBADO · 0 altos, 1 medio, 2 bajos

| Severidad | Abiertos | Ids |
|---|---|---|
| alta | 0 | — |
| media | 1 | **R5-01** (nuevo, front-ajustes-a): en P-26.01, P-26.02, P-26.03, P-26.05 y P-26.07 hay guardados y borrados que fallan sin avisar (U4 §0.3). Por la misma causa nunca aparece V-18 al eliminar una promoción. |
| baja | 2 | **H-13** (reabierto: viene de r1 y nunca se cerró, front-creditos-salud-reportes) y **R3-01** (sigue abierto, frontend-base) |

**HECHO.** R4-01 y R4-02 quedaron cerrados, con prueba del área y con una sonda mía. El resto de lo que pide el encargo cumple:
- build y tsc;
- 119 pruebas de las 3 áreas, 0 fallidas;
- 0 rutas fuera del contrato y 0 hex fuera de los tokens;
- 0 textos en inglés, 0 fechas ISO y 0 dinero sin formato en 128 páginas;
- variantes con «RESULTADO: OK»;
- los campos de U4 presentes en las 14 secciones;
- las 11 claves `reports.*`/`alerts.*`;
- U5-09 en 3 toques.

Lo que impide aprobar es R5-01. Es el mismo tipo de defecto que R4-01: un `await pedir(…)` sin `catch` en un manejador que se lanza con `void`. Esta vez está en otros seis guardados y en el diálogo de acción crítica común de ajustes-a. La ronda 4 solo había señalado el de P-26.02 y el de P-26.13.

## Tabla criterio → estado → evidencia (mi corrida)

| Criterio | Estado | Evidencia |
|---|---|---|
| C3-03 build | CUMPLE | `npx tsc --noEmit -p .`: rc 0. `npx tsc --noEmit -p .storybook/tsconfig.json`: rc 0. `npx vite build --outDir <tmp>/dist`: «✓ built», «precache 81 entries». |
| C3-03 rutas | CUMPLE | `npm run verificar:rutas` da «208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)». `python3 -I scripts/c3-integracion.py`, sin `--escribir`, da «210 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (routes.go: 216; openapi: 216 operaciones)» e «implementado sí 42/42». |
| U5-08 (referencia) | CUMPLE | `node scripts/hosts-externos.mjs <tmp>/dist`: «hosts externos que se pedirían = 0». Solo hay literales que no son peticiones: la licencia de Inter y bit.ly en workbox. |
| Pruebas Playwright de las 3 áreas | CUMPLE | `npx playwright test -c <tmp> ajustes-a ajustes-b creditos-salud-reportes`: **119 passed, 1 skipped, 0 failed** (2.2 min). Entre ellas pasan «R4-01 · P-26.02…» y «R4-02 · P-26.13…» en `movil` y en `escritorio`. |
| U5-06: 0 hex y colores solo por token | CUMPLE | `grep` de `#rgb`/`#rrggbb` en `src` (ts, tsx y css): solo aparecen en `src/estilos/tokens.css`, que es generado. `rgb(`, `hsl(`, `oklch(` y `dark:` en ts/tsx: 0. |
| U5-06: variantes con contraste (front-variantes) | CUMPLE | `npm run build-storybook` y `node scripts/contraste-variantes.mjs --salida <tmp>`: 37 historias, 148 páginas, 5 402 mediciones únicas y 0 razones distintas contra U2d. Resultado: «FALLAS nuevas: 0», «0 de 5 se reproducen» (FV-01 a FV-05 CERRADO), «problemas de la corrida: 0» y «RESULTADO: OK». |
| U5-06 y U5-07: barrido visual | CUMPLE | Sonda propia con 18 rutas de P-17 a P-25 (con sus subvistas resumen, detalle, tendencias y calendario, P-19 nueva y existente, y P-20) y /ajustes/01 a /ajustes/14. Cada ruta en claro y oscuro, a 412 táctil y a 1440: **128 páginas**. Resultado: 0 `pageerror`, 0 desbordes horizontales, 0 controles bajo 44 px en táctil, 0 cifras de dinero sin `tabular-nums` y 0 textos en inglés (`textosEnIngles`). La segunda pasada a 1440 da 0 fechas ISO visibles y 0 cifras `$1234` sin formato en 64 páginas. El único «No se pudo» del barrido es «No se pudo enviar» en /reportes: es el estado de envío fallido de un reporte en los datos de prueba, como pide N-245. |
| U5-09 | CUMPLE | `-g "U5-09\|F-1 desde Inicio"`: 4 passed (2 pruebas × 2 proyectos), registrar gasto en 3 toques. |
| U4: campos de las 14 secciones | CUMPLE | Script propio sobre las tablas de U4: 242 filas de campos y acciones, **0 faltantes**. 207 tienen el texto exacto en el código y 35 tienen todas sus palabras. El detalle está en la tabla de abajo. |
| U4 §0.2: «Antes → Después» | CUMPLE | Las 13 secciones con efecto numérico tienen vista previa en línea (P-26.10 Apariencia no tiene efecto numérico). Con `Ctrl+/` y con el tablero, las pruebas de 02, 04, 05, 07, 09, 13 y 14 pasan. También pasan las de K44 en P-26.01, las vistas previas 6, 7 y 10 de P-26.05 y la de frescura de P-26.06. |
| U4 §0.6: catálogo V-xx | CUMPLE PARCIAL | ajustes-a define 23 textos V y ajustes-b define V-16. V-13 llega con el mensaje del servidor, que es el mismo texto (`auth.go:210`). **V-18 se define (`textos.ts:17`) y no se usa**: ver R5-01. |
| U4: claves `reports.*` y `alerts.*` editables | CUMPLE | Las 11 claves de C1 (`extractos/ajustes-f.yaml`) aparecen en P2608 y P2609 y se guardan con `PATCH /api/v1/settings/reports` (`updateReportSettings`) y `PATCH /api/v1/settings/alerts` (`updateAlertSettings`), con `parameters` y `reset_parameters`. |
| U4 §0.3 «Error al guardar» | **NO CUMPLE** | **R5-01**. Sonda propia en los 2 proyectos: los 3 casos medidos dan `pageerror` = 1 y 0 «No se pudo guardar». |
| U1: P-17 a P-25 | CUMPLE, salvo H-13 (baja) | Los títulos de U1 §3.19 a §3.27 están en `creditos-salud-reportes/textos.ts`: Deuda viva total, Calendario, Compras de la línea, Tasas por plazo, Estado de la tarjeta, Amortización, Abonos, Vista previa del cálculo, Promociones vigentes, Faltan restricciones, Usar escenario, Calorías de hoy, Historial de peso, Meta sin confirmar, «El de la casa: n kcal», Próximos envíos, Avisos enviados, Reintentar envío, Generar PDF, «parcial» y Estados anteriores. Las pruebas por pantalla pasan. Falta «Comparador no disponible» (H-13). |
| DESIGN.md | CUMPLE | Lo cubren los controles de U5-06 de arriba: variantes, barrido y tokens. |

### U4 sección por sección (mi script y mi lectura)

| Sección | Campos U4 (exactos + aprox.) | Faltan | Vista previa | Error al guardar (§0.3) |
|---|---|---|---|---|
| P-26.01 Cuenta y roles | 17 (14 + 3) | 0 | K44 (prueba) | **R5-01**: «Cerrar sesiones» falla sin aviso (medido). Además, descargar el certificado tampoco avisa si falla (`S01Cuenta.tsx:152`). |
| P-26.02 Presupuesto | 20 (16 + 4) | 0 | Vistas previas 1, 3 (g, h) y 4 | Monto: CUMPLE, R4-01 cerrado. **R5-01**: «Quitar el monto» y «Nuevo gasto fijo». |
| P-26.03 Ingresos | 7 (7 + 0) | 0 | En línea | **R5-01**: «Capturar ingreso quincenal». |
| P-26.04 Categorías | 8 (8 + 0) | 0 | Tablero y `Ctrl+/` (prueba) | CUMPLE (`S04Categorias.tsx:102,146`). |
| P-26.05 Créditos y condiciones | 51 (39 + 12) | 0 | 6, 7 y 10, A/B y `Ctrl+/` | **R5-01**: «Pesos iguales» y «Nueva promoción» medidos; «Nueva tasa», «Nuevo crédito» y eliminar tasa o promoción, por código. V-18 no aparece. |
| P-26.06 Tiendas y precios | 24 (21 + 3) | 0 | Frescura PROVISIONAL | CUMPLE (`S06Tiendas.tsx:153,225`). |
| P-26.07 Recetas | 10 (8 + 2) | 0 | `Ctrl+/` | **R5-01**: eliminar receta, que usa el diálogo crítico común. Nueva receta CUMPLE (`:145`). |
| P-26.08 Reportes | 14 (9 + 5) | 0 | Horas y días | CUMPLE (`g.ejecutar`). |
| P-26.09 Notificaciones | 24 (22 + 2) | 0 | `Ctrl+/` sumado | CUMPLE. |
| P-26.10 Apariencia | 2 (2 + 0) | 0 | No aplica | No aplica: es local. |
| P-26.11 Datos y respaldos | 20 (18 + 2) | 0 | Parámetro | CUMPLE: restaurar tiene `catch` y activar la IA usa `g.ejecutar`. |
| P-26.12 Personas y reparto | 12 (12 + 0) | 0 | Alcance | CUMPLE (`g.ejecutar`). |
| P-26.13 Salud y metas | 24 (23 + 1) | 0 | Perfil y `Ctrl+/` | CUMPLE: R4-02 cerrado. |
| P-26.14 Inventario | 9 (8 + 1) | 0 | `Ctrl+/` | CUMPLE. |

## Hallazgos abiertos

| Id | Sev. | Criterio | Ruta | Dueño | Evidencia mía | Corrección sugerida |
|---|---|---|---|---|---|---|
| R5-01 | media | U4 §0.3 «Error al guardar» (rige las 14 secciones; P-26.01, 02, 03, 05 y 07 dicen «Carga, error y guardado según §0.3»). U4 P-26.05, fila «Eliminar promoción»: «si hay compras que la usan → V-18 y ofrece desactivar» (U4:293; C1 `deletePromotion` responde 409 `in_use`). | (1) `frontend/src/pantallas/ajustes-a/comunes.tsx:228`, `DialogoCritico`: `try { await alConfirmar(); alCambiar(false) } finally {…}` sin `catch`. Lo usan `S01Cuenta.tsx:137`, `S02Presupuesto.tsx:119`, `S05Creditos.tsx:297` y `:365`, y `S07Recetas.tsx:84`. (2) `S05Creditos.tsx:371` y `:392`: «Pesos iguales» con `void patch(…)`. (3) `S05Creditos.tsx:268-272`: «Nueva tasa». (4) `S05Creditos.tsx:405-411`: «Nuevo crédito». (5) `S05Creditos.tsx:436-445`: «Nueva promoción». (6) `S02Presupuesto.tsx:226-231`: «Nuevo gasto fijo». (7) `S03Ingresos.tsx:60-65`: «Capturar ingreso quincenal». Ninguno tiene `catch`. `textos.ts:17` define `v18` y nadie lo usa. | front-ajustes-a | Sonda propia en `movil` y `escritorio`, con el resto de la API igual que en `preparar()` de las pruebas del área. **S-C**: «Pesos iguales» con `PATCH /settings/optimizer` en 500 → `pageerror` 1 («Error interno del servidor.») y 0 «No se pudo guardar». **S-D**: «Nueva promoción» llena, `Ctrl+Enter`, con `POST /promotions` en 500 → `pageerror` 1, 0 mensajes y la hoja abierta sin aviso. **S-F**: «Cerrar sesiones» con `POST /sessions/revoke` en 500 → `pageerror` 1, 0 mensajes y el diálogo crítico abierto sin aviso. Los casos (3), (4), (6) y (7) y los borrados por el diálogo los vi en el código, con el mismo patrón y la misma consecuencia. En cambio S04, S06 y S07 (nueva receta) sí capturan el error, y ajustes-b pasa todo por `g.ejecutar`. | (a) En `DialogoCritico`, capturar el error. Mostrar dentro del diálogo el texto de §0.3 «No se pudo guardar. Revisa la conexión con el servidor de casa · Reintentar» (con `circle-alert` y `text-role-error`) y dejarlo abierto. Para eliminar una promoción con `ErrorApi.codigo === 'in_use'`, mostrar `t.v.v18` y la acción «Desactivar» (`PATCH /api/v1/promotions/{promotion_id}` con `is_active: false`). (b) Envolver los handlers (2) a (7) en `try/catch`, con el mismo trato que `S06Tiendas.tsx:225`: 409 → su V-xx, cualquier otro error → `t.comun.noGuardo`, sin perder lo capturado en la hoja. (c) Agregar pruebas: 500 en «Pesos iguales», en un alta y en un borrado crítico, y 409 `in_use` al eliminar una promoción. Cada una debe comprobar el mensaje y 0 `pageerror`. |
| H-13 | baja | U1 §3.22 N-226: «si se pasa [el tiempo máximo], P-20 avisa “Comparador no disponible” sin resultados parciales» | `frontend/src/pantallas/creditos-salud-reportes/P20Comparador.tsx:67` y `:113` (`ErrorCarga` con el mensaje del servidor) | front-creditos-salud-reportes | **S-G**, sonda propia: `POST /scenario-comparisons` responde 503 `{code:'unavailable'}`. La pantalla muestra la franja «tiempo máximo agotado · Reintentar». Lo esperado era «Comparador no disponible», que aparece 0 veces; `pageerror` 0. El texto real del servidor es «tiempo máximo agotado» (`motor_cred/errors.go:19`). Este hallazgo ya estaba en r1 como baja y no se siguió en r2, r3 ni r4. | Con `ErrorApi.codigo === 'unavailable'`, mostrar un texto propio de `textos.ts`: «Comparador no disponible» (y Reintentar), sin resultados. |
| R3-01 | baja | Higiene de pruebas | `frontend/package.json` | frontend-base | `ajv` no está ni en `dependencies` ni en `devDependencies`. `pruebas/ajustes-a/validar-c1.ts:7` lo importa (`ajv/dist/2020.js`), así que solo funciona como dependencia transitiva. | Agregar `"ajv": "8.20.0"` a `devDependencies`. |

## Estado de los hallazgos de la ronda 4

| Hallazgo | Sev. | Estado | Evidencia mía |
|---|---|---|---|
| R4-01 (P-26.02: 409 `budget_overlap` sin V-17, y sin «No se pudo guardar») | media | **CERRADO** | Código: `S02Presupuesto.tsx:77-84`, `try/catch` con `ErrorApi.codigo === 'budget_overlap'` → `t.v.v17` (`:113`); cualquier otro error → `t.comun.noGuardo` con Reintentar (`:114`). La prueba «R4-01 · P-26.02…» pasa en los 2 proyectos (camino `PATCH`, 409 → 500 → Reintentar 200, `pageerror` = []). Mi sonda **S-A** probó el camino `POST` (monto nuevo de una categoría sin monto, 409): V-17 visible 1 vez, hoja abierta, `pageerror` 0, en los 2 proyectos. |
| R4-02 (P-26.13: «Usar el de la casa» y «Usar la sugerida» sin aviso) | media | **CERRADO** | Código: `P2613Salud.tsx:113` y `:125` pasan por `g.ejecutar` (que captura y pone `t.comun.noSeGuardo` con Reintentar, `comun.tsx:79-91`), y el error se pinta con `IndicadorGuardado` del campo (`:85`). La prueba «R4-02 · P-26.13…» pasa en los 2 proyectos (500 → 2 mensajes, cuerpos exactos, `pageerror` = []). |
| R3-01 (`ajv`) | baja | **ABIERTO** | Ver la tabla de hallazgos. |

## PREGUNTA para el orquestador o Axel

R5-01 es del mismo tipo que R4-01, pero está en otros puntos de ajustes-a. Se corrige dentro del alcance de front-ajustes-a: un `catch` en `DialogoCritico` y en seis handlers, más V-18. ¿Se encarga a esa área con una revisión puntual después, o se acepta como pendiente? H-13 y R3-01 son bajos y no impiden aprobar. No decido yo.

## Lo que no corrí

- `make verify-wf3-frontend`, porque escribe en `frontend/dist` y usa el puerto 4173. En su lugar corrí lo mismo de forma aislada: tsc, build, rutas, hosts, Playwright de las áreas, U5-09 y variantes.
- La suite completa de Playwright de todas las áreas: solo corrí las 3 áreas que me tocan y las 2 pruebas de U5-09.
- `make verify-all`: no toqué nada del repositorio fuera de mis 2 reportes.
- No medí con sonda «Nueva tasa», «Nuevo crédito», «Nuevo gasto fijo», «Capturar ingreso» ni el 409 `in_use` de promociones. Tienen el mismo código sin `catch` que los 3 casos medidos.
- Storybook se reconstruyó en `frontend/storybook-static`, porque `contraste-variantes.mjs` lo lee de ahí fijo. Es una salida ignorada por git (`frontend/.gitignore:16`) y se generó con las mismas fuentes.
