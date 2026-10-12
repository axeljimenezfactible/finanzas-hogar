# front-gastos-presupuesto · WF-3a (12-oct-2026)

## Hecho
- Área `frontend/src/pantallas/gastos-presupuesto/`, registrada sola por `indice.ts` (convenciones §3). Pantallas:
  - **P-02**: resumen (6 tarjetas), detalle y calendario, con los filtros en la consulta (`periodo`, `familia`, `tipo`, `pago`, `ticket`, `categoria`), fila de totales y «Mostrar más» por cursor. Vista de BD con tabla, tablero por categoría (arrastrar cambia la categoría), calendario y galería. Se edita en el lugar fecha, descripción, categoría, monto («Antes → Después») y tipo. Tendencias: 6 categorías + «Otras». Muestra el vacío de U1 y el vacío con filtros.
  - **P-04**: propiedades editables, reparto, devengo, efectivo y renglones. Acciones: adjuntar ticket, duplicar, marcar pagado (solo venideros), abrir la compra a crédito y eliminar con confirmación. Si el gasto ya no existe: «Este gasto ya no existe» + Volver.
  - **P-05**: 4 pasos. Cámara o galería; «Leyendo el ticket…» con la foto y consulta cada 2 s; si falla: «No se pudo leer» con Repetir y Capturar a mano. Los renglones se pueden editar y muestran su confianza (baja resaltada). Incluye Agregar renglón, revisa si la suma cuadra con el total, Guardar para después (K3), Confirmar (abre P-04), Cancelar y la lista de tickets por confirmar. Sin conexión muestra «Requiere conexión».
  - **P-06**: 6 tarjetas, barra de avance `role=meter` y alertas. Los estados normal, aviso y excedido vienen del API (no hay umbrales fijos en el código). Muestra quincenas, tabla y tablero por estado, «Ver gastos» y la despensa devengada. **Ajustar presupuesto** en una hoja: alcance K2 con «De este mes en adelante» preseleccionado, «Antes → Después» del plan y del restante, `updateBudget` (con scope y month) o `createBudget`. Tiene tendencias por categoría y calendario de dinero.
  - **P-07**: interruptor A/B (`role=switch`, `updateProjectionSettings`) y modo usado. Muestra el saldo de la quincena, el acumulado, la diferencia acumulada y la deuda viva a la fecha elegida, más los supuestos. Detalle con quincenas y meses; tendencias devengado contra efectivo. «Simular compra» abre P-20.
  - **P-08**: tarjetas; quincenas esperadas o capturadas con «Capturar quincena» (prellenada con lo esperado); extras editables en el lugar y eliminar. Abre las hojas con `?hoja=quincena` y `?hoja=ingreso`; tiene tendencias y calendario.
  - **P-09**: totales, metas con avance y lo necesario por mes; galería, tablero por estado y calendario de metas. Movimientos editables y eliminar; aportación o retiro con `?hoja=movimiento`; «Nueva meta» abre P-26.13; tendencias y calendario.
- Pruebas: `frontend/pruebas/gastos-presupuesto/`, 8 pruebas que corren en 412×915 (táctil) y en 1440×900. Cubren las 7 pantallas, F-1 (U5-09) y el flujo complementario del ticket (P-05 → P-04).

## Archivos
`frontend/src/pantallas/gastos-presupuesto/{indice.ts,textos.ts,comun.tsx,P02Gastos.tsx,P04Gasto.tsx,P05Ticket.tsx,P06Presupuesto.tsx,P07Flujo.tsx,P08Ingresos.tsx,P09Ahorro.tsx}` · `frontend/pruebas/gastos-presupuesto/{datos.ts,gastos-presupuesto.escritorio.spec.ts}` · `docs/estado/bloqueos/WF3-gastos-contrato.md` · `docs/estado/bloqueos/WF3-gastos-ticket-credito.md`

## Verificación
- `npx tsc --noEmit`: 0 errores en esta área (las pruebas incluidas). El tsc y el `npm run build` globales fallan hoy por **otras áreas aún en obra** (ajustes-a, ajustes-b, comidas-compras, creditos-salud-reportes). Por eso hice el build en una copia aislada en `datos/tmp/` con solo esta área: tsc 0 y `vite build` OK (PWA, 33 entradas).
- `npx playwright test pruebas/gastos-presupuesto`: **16/16** (movil y escritorio). Las pruebas base cola, login y shell: 10/10 en la misma copia.
- **U5-09**: 3 toques en movil y escritorio, y el gasto aparece en la tabla de P-02.
- **U5-07**: 0 palabras de la lista en inglés (Save, Cancel…), 0 fechas ISO en los textos; dinero con `dinero()` y fechas con `fecha()`/`mes()`.
- **U5-06**: 0 hex, 0 rgb/hsl/oklch, 0 `dark:` y 0 valores arbitrarios `[...]`. Objetivos de ≥ 44 px medidos en movil (Registrar gasto, Tomar foto, interruptor). (Nota de `correccion-frontend-core-r1`, H-04: solo se habían medido 3 controles; los botones de categoría de P-06 Detalle medían 71 × 20 y el enlace «Todos los gastos» de P-04 130 × 20. Corregidos en `vistas-bd/VistaTabla.tsx` y `shell/Barras.tsx`; `pruebas/medidas.spec.ts` mide hoy todos los controles visibles de las 15 vistas del área: 0 menores de 44 × 44.) Sin desplazamiento horizontal a 412 px en las 7 pantallas; foco visible con las clases de la base.
- **U5-13**: P-02, P-04, P-06, P-08 y P-09 desactivan las escrituras sin red con «Requiere conexión con el servidor de casa»; P-05 pone el aviso en lugar del contenido; no toqué la cola (K9).
- **C3-03** (`npm run verificar:rutas`): 200 rutas, 0 fuera de openapi.yaml, 0 sin registrar en routes.go. Sin `fetch` directo.
- `make verify-all`: falla en `verify-go` por `internal/api/comprasa` (handlers de otro agente aún sin escribir). Las migraciones pasan (55 tablas, goose up/down/up). No viene de esta área.

## Supuestos
- Nombre del área `gastos-presupuesto` y pruebas en `pruebas/gastos-presupuesto/` (convenciones §1, §3 y §8). Uso ApiFalsa + `context.route` según §8 y no los mocks de msw.
- El periodo de P-02, P-06 y P-08 es el mes calendario en `?periodo=AAAA-MM` (por omisión, el mes en curso). Los demás periodos de C1 (quincena, semana, rango) quedan sin selector.
- Concepto de la quincena capturada: «Quincena» (`description` es obligatorio en `IncomeCreate`).
- Ajustes › Ingresos = `/ajustes/03` y Metas = `/ajustes/13` (rutas PROVISIONALES, bloqueo WF3-frontend-rutas).
- El ancho del llenado de la barra de avance va en `style` como porcentaje del dato; no es un color ni una medida de diseño.

## Bloqueos
- `WF3-gastos-contrato`: adjunto del gasto, vista previa del reparto, filtros de los totales y búsqueda de catálogo en P-05.
- `WF3-gastos-ticket-credito`: datos de la compra a crédito al confirmar un ticket.

## Pendientes
- Lo que depende de los dos bloqueos. Revisión visual contra los artboards de la Pista 3 (`lienzo/project/gastos|presupuesto/*.dc.html`): no los leí por el tope de lectura. Mocks msw propios (`mocks.ts`, opcional §3.8). Volver a correr el build y la verificación globales cuando las demás áreas compilen.

## Propuestas para Axel
- Bloqueo 1: opción (a), `attachments` dentro de `Expense`, con un solo cambio en C1.
- Bloqueo 2: opción (a), vista previa en el servidor (el motor S2 es la única fuente de las partes).
- Bloqueo 3: opción (a).
- Ticket a crédito: opción (a), que respeta los +2 toques de U1.

## Paquetes instalados
Ninguno.

## Consumo aproximado
~54 acciones · ~165 KB leídos · ~145 KB escritos.
