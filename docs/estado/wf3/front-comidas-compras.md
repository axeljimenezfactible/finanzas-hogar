# front-comidas-compras · WF-3a · 12-oct-2026

## Hecho
- P-10 a P-16 en `frontend/src/pantallas/comidas-compras/` (el cascarón las registra solo con `indice.ts`). Cada pantalla tiene sus subsecciones de U1, sus vistas de BD, los estados vacío, carga y error de U1 §2.2 (mediante `ConDatos`), sus textos es-MX en `textos.ts` y sus hojas. Usa solo los tokens de C2 y el cliente tipado (`pedir` y `usarConsulta`), con `Idempotency-Key` en los POST que la piden.
  - P-10: semana por defecto de §0.2 (domingo abre la siguiente, K38) y selector de semana. Cuadrícula días × tiempos con casillas de 44 px. Hoja «Elegir receta» con fila de días, tiempo, Ambos, Axel o Ximena, porciones, búsqueda y orden por uso en 8 semanas; ahí también están «Fuera de casa» y «Vaciar casilla». «⋯» tiene «Copiar semana anterior» y «Asignar receta». La primaria «Generar lista de compras» abre P-13. Resumen en 5 tarjetas; Tendencias en 2 gráficas; Calendario mensual; tabla, tablero y galería de asignaciones.
  - P-11: galería*, tabla y tablero. Filtros por texto, tiempo, origen y costo o kcal máximos. Edición del nombre en el lugar. «Nueva receta» y «Sugerir con IA» (bloqueo #4).
  - P-12: porciones con − y +, que recalculan con `servings_milli` del motor S4. 4 totales. Tabla de ingredientes con «—» y «Sin precio», o «sin fuente». Pasos en bloques, uso en 8 semanas, «Agregar al plan», y Editar, Duplicar o Eliminar (este con aviso si hay plan futuro y `confirm_planned`). La ruta `nueva` crea la receta.
  - P-13: `week` en toda lectura y escritura. Total, totales por tienda o mercado con día de visita, tabla* y tablero por tienda (mover = cambiar tienda). Casilla «comprado». Hojas «Agregar artículo» (`?hoja=articulo`) y «Cambiar tienda o cantidad». Vacíos «Nada que comprar…» y «Sin plan esta semana · Planear».
  - P-14: selector de canasta (lista de la semana o compras registradas). En Resumen: frase de la más barata, tarjetas por tienda (total, +diferencia, +en idénticos, cobertura, idénticos, sustitutos, sin precio y viejos), ganadores con empates, membresía de Sam's, «Usar <ganadora>» y «Simular pago con crédito». En Detalle: artículos × tiendas.
  - P-15: Resumen (fresco, viejo, sin precio y fecha por tienda). Detalle con 5 filtros, tabla*, tablero por frescura y galería. Hoja del producto con presentaciones, precios (corregir o eliminar; PROFECO no se edita) y nutrición manual con cita. Hojas «Capturar precio» y «Nuevo producto». Tendencias con el historial por tienda.
  - P-16: Resumen por estado. Detalle con filtros y orden agotado, reponer, suficiente y nombre. Existencia editable en el lugar con «Antes → Después» (conteo). Hojas de movimientos, «Ajustar existencia» con motivo y «Registrar sobrante».
- Pruebas en `frontend/pruebas/comidas-compras/` (CONVENCIONES §1 y §8, con `ApiFalsa` más rutas del área):
  - `pantallas.escritorio.spec.ts`: 7 pruebas, una por pantalla, en 412×915 y 1440×900, que además comprueban que no haya desborde horizontal.
  - `flujos.spec.ts`: F-2, F-3 ruta A y F-3 ruta B, en 412×915 táctil.

## Archivos
`frontend/src/pantallas/comidas-compras/{indice.ts,textos.ts,comun.tsx,P10Plan.tsx,P11Recetario.tsx,P12Receta.tsx,P13Lista.tsx,P14Canasta.tsx,P15Productos.tsx,P16Inventario.tsx}`, `frontend/pruebas/comidas-compras/{datos.ts,pantallas.escritorio.spec.ts,flujos.spec.ts}`, `docs/estado/bloqueos/WF3-comidas-compras-interfaz.md`, este reporte.

## Verificación
- `npx tsc --noEmit`: **0 errores en comidas-compras**. El total da 48 errores, todos en `src/pantallas/ajustes-a/` (otro agente, en curso). Por eso `npm run build` (que corre tsc completo) no pasa hoy por causa ajena.
- `npx vite build --outDir ../datos/tmp/fcc/dist`: ✓ built, PWA generada.
- `npx playwright test` (configuración temporal en `datos/tmp`, puerto 4183, mismo Chromium y proyectos que la base): **17/17 pasan** (7 pantallas × 2 tamaños + 3 flujos).
- U5-06: 0 hex literales en la carpeta (`grep -E '#[0-9a-fA-F]{3,8}'` = 0). Tampoco hay `style=`, `dark:`, `[npx]`, `rgb(` ni `hsl(`: solo clases de C2. ~~Todos los controles miden 44 px con puntero grueso~~ (corregido por `correccion-frontend-core-r1`, H-05: la afirmación no se había medido; la casilla de P-13 dibuja 20 × 20 dentro de una etiqueta de 44 × 44 que la marca, y el enlace de regreso «Recetario» de P-12 medía 20 px de alto (ahora 44, `shell/Barras.tsx`); hoy lo mide `pruebas/medidas.spec.ts` en las 10 vistas del área: 0 controles menores de 44 × 44) y llevan `focus-visible` con el anillo de `border-focus`.
- U5-07: el grep `Save|Cancel|Settings|Dashboard|Loading|Error|Home|Delete|Search|Add` da 0 en la carpeta. El dinero sale de `dinero()` (canon `$3,094.50`, `+$138.50`) y las fechas de `fecha()` y `fechaConDia()` (`12-oct-2026`). Ninguna fecha ISO es visible.
- U5-08: `hosts-externos.mjs` sobre el build solo encuentra literales que no son peticiones (licencia Inter, workbox). Hay 0 `fetch` directos en la carpeta.
- U5-09 y U5-13: no aplican a estas pantallas (gastos y peso son de la base). Las escrituras sin red quedan «Requiere conexión con el servidor de casa» (K9).
- C3-03: `npm run verificar:rutas` da «201 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0». De ellas, 34 rutas literales son de esta carpeta: comidas.yaml, compras-tiendas-a.yaml y compras-tiendas-b.yaml más `/api/v1/persons`.
- C3-04 (frontend): F-2 hace el POST `meal-plan-entries` con `{plan_dates:['2026-10-12'], meal_slot_id:3, person_ids:[1,2], recipe_id:1, servings_milli:1000}` y la clave de 16 caracteres o más, y luego `generate?week=2026-W42`. F-3 A manda `assign-store?week=2026-W42` con `{store_id:2}`. F-3 B pide `basket-comparison?basket=expense&expense_id=77`.
- Canon S5 en P-14: Bodega $3,094.50 contra Walmart $3,233.00, +$138.50 y +$100.50 en idénticos, 33 idénticos, 3 sustitutos y empates 8.
- `make verify-all GOOSE=…`: **falla en `verify-wf1-cobertura`**: `internal/dominio` da 72.9 % (< 85 %). Es backend; no toqué Go.

## Supuestos
Los PROVISIONALES están en el bloqueo `WF3-comidas-compras-interfaz` (#1–#7). Además:
- «Casillas sin armar» se cuenta sobre 7 × tiempos obligatorios = 21, como en el ejemplo F-2 de U1.
- Las personas del plan son las que no son externas (Karen queda fuera).
- Las pruebas van en `pruebas/<área>/` (CONVENCIONES §8) y no dentro de `src/pantallas`.
- Los mocks msw de área son opcionales (§3.8) y no los hice: las pruebas usan `context.route`.

## Bloqueos
`docs/estado/bloqueos/WF3-comidas-compras-interfaz.md`.

## Pendientes
- Correr `npm run build` y `npx playwright test` completos cuando `ajustes-a` compile.
- Capturas U5-01 (las toma el revisor).

## Propuestas para Axel
- Agregar `code` a `Store` en C1 (bloqueo #3) para no depender del nombre.
- Que P-05 lea `?semana=` y P-20 lea `?monto=` (#1, #2).
- Que el Playwright de la base acepte puerto y carpeta de build por variable, para que varios agentes prueben a la vez sin pisarse.

## Paquetes instalados
Ninguno.

## Consumo aproximado
~52 acciones · ~170 KB leídos · ~148 KB escritos.
