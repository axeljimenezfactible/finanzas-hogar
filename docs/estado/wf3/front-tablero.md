# front-tablero · P-01 Resumen = tablero de control (U3 2.0, K54) y U5-06 #10

Agente `front-tablero` (WF-3c, 12-oct-2026). Contrato C1 v1.5.0 (tipos regenerados con `npm run tipos`; mocks con `npm run mocks:gen`: 221 operaciones). Sin `git add` ni `git commit`.

## Hecho
- **Tablero** (`src/pantallas/inicio/Tablero.tsx`): K-1 «Gasto del mes» (indicador principal: gastado, % del plan con barra `role="meter"`, efectivo y comparación con el mes anterior, insignias), B-1 barómetro, K-4 «Pagos próximos» (N de `getDashboardSettings`, vencidos + de hoy a hoy + 30, filas-enlace a P-18), K-2 «Ingreso del mes», K-3 «Créditos» (deuda viva, pagos del mes, falta pagar, disponible con la lectura activa), G-3 (barras en % del plan, marcas de `alerts.budget_*`, sugerencia subir/recortar, fila-enlace a P-02 con `categoria` y `periodo`), G-1 (columnas agrupadas de 3 series con patrón), G-6 (línea con proyección 8 · 4, referencias en `warning_pct` y `limit_pct`, «Hoy»), G-2 (deuda viva por línea), G-5 = T-3 sin filtros, G-4 = T-2, y «Más del hogar» con las tarjetas 2, 4, 5 y 6. Orden del DOM = orden de 360/412; rejilla 8 columnas en 768 y 12 desde 1024 (K-4 a la derecha de K-2 y K-3 con colocación explícita, sin cambiar el foco). G-1 y G-6 de **una** petición `from` = mes − 11, `to` = mes + 5. Ningún umbral, porcentaje, diferencia ni estado se calcula: llegan de C1.
- **Barómetro** (`src/componentes/Barometro.tsx`, DESIGN §8.9): SVG en línea, tres tramos con `stroke-role-success/warning/error`, pista `viz-track`, valor en el color del tramo, aguja `text-default` con eje y anillo `surface-card`, marcas de 1 px `viz-baseline` con «[n] %» en HTML (tamaño de la escala tipográfica), dígito `title-1` bold tabular, «de [tope] % de tope», insignia con icono y palabra; `role="meter"` con `aria-valuetext`; sin dato: arco sin aguja y «—»; sin animación; `forced-colors`.
- **Resumen.tsx** pasa a exportar la tarjeta de métrica (`Tarjeta`, `Vista` con `cuerpo`) y `MasDelHogar`; se retiraron las tarjetas 1 y 3 (pasaron a K-1, K-3 y K-4). `Tendencias.tsx` exporta `T2` y `T3` (`id`, `sinFiltros`). `GraficaInicio.tsx`: descripción accesible «título. resumen» en español (CONVENCIONES §6). `comun.tsx`: `getDashboardSettings` leído una vez en el proveedor de Inicio.
- **Barra lateral** (U5-06 #10): icono de sección de 20 px (trazo 1.8 del proveedor), `text-icon-default` y `text-icon-selected` en la activa, entre el chevron y el nombre; ancho de 240 px intacto. 4 iconos PROVISIONALES (bloqueo `WF3c-barra-iconos-seccion`).
- **Mocks** (`src/pantallas/inicio/mocks.ts`): las 4 lecturas nuevas con datos SINTÉTICOS coherentes para `npm run dev:mocks`.
- **Textos** es-MX en `src/pantallas/inicio/textos.ts` (`tablero`) y `src/textos/es-MX.ts` (`barometro`).

## Archivos
Nuevos: `frontend/src/pantallas/inicio/Tablero.tsx`, `frontend/src/componentes/Barometro.tsx`, `frontend/src/pantallas/inicio/mocks.ts`, `frontend/pruebas/inicio/tablero.escritorio.spec.ts`, `docs/estado/bloqueos/WF3c-barra-iconos-seccion.md`, `docs/estado/bloqueos/WF3c-tablero-interfaz.md`, este reporte.
Cambiados: `frontend/src/shell/BarraLateral.tsx`, `frontend/src/pantallas/inicio/{P01Inicio,Resumen,Tendencias,GraficaInicio,comun,textos}.ts(x)`, `frontend/src/textos/es-MX.ts`, `frontend/pruebas/inicio/{datos.ts,inicio.escritorio.spec.ts}`. Regenerados por script: `frontend/src/api/esquema.gen.ts`, `frontend/src/mocks/*.gen.json`.

## Verificación
- `npx tsc --noEmit`: 0 errores. `npm run build`: listo. `npm run verificar:rutas` (C3-03): 212 rutas invocadas, 0 fuera de openapi.yaml, 0 sin registrar en routes.go (221). `npm run verificar:hosts` (U5-08): sin peticiones externas (2 literales que no son peticiones, ya existentes).
- `npx playwright test pruebas/inicio`: 26 de 26 (412×915 táctil y 1440×900). Suite completa `make verify-wf3-frontend`: **253 passed, 4 skipped, 0 failed**.
- U3 A-01: orden de los 12 bloques en el DOM en los dos anchos. A-02: una petición por operación; `monthly-overview?from=2025-11&to=2027-03` una sola vez; alertas y ajustes del tablero una vez. A-05: K-3 «Pagos de octubre: $973.75 · 1 pago», «Disponible hoy $6,905.50 · Lectura B», «Pendiente confirmar»; K-4 «Vie 30 oct 2026 · Kueski · $973.75 · Pendiente» y «Hogar $598.82 · Karen $374.93». A-06: octubre $1,197.64, $598.82, «Sin datos de septiembre»; noviembre «Gastado · Igual que octubre», «Efectivo · $598.82 más que octubre». A-07: falla de B-1, K-4, G-1 y G-6 sin tumbar K-1 ni K-3, y «Reintentar». A-08: K-1, K-2, K-3 y B-1 un solo enlace. A-09/U5-06: 0 controles < 44 × 44 con puntero grueso, 0 iconos fuera de norma (incluida la barra lateral a 1440), sin desplazamiento horizontal. A-14/U5-07: 0 textos en inglés (incluida la descripción de ECharts), dinero tabular, sin fechas ISO; claro y oscuro. A-17: K-4 «Vence mañana» el 29-oct. A-18: 30 verde, 31 y 40 amarillo, 41 rojo, `aria-valuetext` «Capacidad de endeudamiento n % de 40 %: …», marca en 25 % con `warning_pct` = 25, tres trazos de token distintos. A-19: G-1 3 series y nota fija; G-6 «30 %», «40 %» y «Hoy», proyección y «— · Falta ingreso». A-20: orden descendente, «Rebasada», «Vas en 85 %», sugerencias con nombre accesible. A-21 paso 1: barra → `/gastos?categoria=1&periodo=2026-10`.
- Greps: 0 hex, 0 `dark:` y 0 colores en `style` en los archivos nuevos; 0 palabras inglesas en los textos.

## Supuestos
- K-2 «quincenas del mes» = `captured_cents` + `expected_pending_cents` (suma de lo visible que U3 K-2 escribe); vacío = sin `fortnights[]` y sin extras.
- K-1 muestra las dos comparaciones (gastado y efectivo) como dice A-06; sin `previous_*` del mes anterior: «Sin datos de [mes]». Un campo de C1 v1.5 ausente no pinta su línea (A-22).
- K-4 pide vencidos sin `period` (solo `status=overdue`) y solo si `next_payment.is_overdue`; el «Ver calendario» lleva a `/inicio?sub=calendario`.
- G-3: «[n] categorías arriba de [warning] %» cuenta `warning` + `exceeded` a la vista.

## Bloqueos
- `docs/estado/bloqueos/WF3c-barra-iconos-seccion.md` (nuevo): iconos de Gastos, Presupuesto, Salud y Reportes.
- `docs/estado/bloqueos/WF3c-tablero-interfaz.md` (nuevo): 7 puntos (conteo «y [n] más», «Reintentar» dentro de K-3, insignia de K-1 hacia G-3, alto de B-1, separación de G-1, G-2/G-3 en HTML, N inválido).
- Siguen los de U3: `U3-tablero-barometro`, `U3-tablero-recordatorio`, `U3-tablero-horizonte`, `U3-tablero-senal-presupuesto`, `U3-tablero-mas-del-hogar`, y `C1-tablero`.

## Pendientes
- Los manejadores Go de las 5 operaciones nuevas (agente `api-tablero`): hasta entonces el servidor responde 501 y el tablero muestra la falla de B-1, G-1, G-6, K-4 (por `settings/dashboard`) y sin sugerencias en G-3.
- Pruebas fuera de mi alcance que buscan las tarjetas 1 y 3 retiradas (ver solicitudes).
- Pasos 2 a 5 del desglose §2.2 (P-02 → P-04 → P-15 › Tendencias) son de las áreas gastos y compras.

## Propuestas para Axel
- Confirmar los 4 iconos de sección y los 7 puntos de `WF3c-tablero-interfaz` en la puerta 4, junto con los bloqueos U3-tablero-*.

## Paquetes instalados
Ninguno.

## Consumo aproximado
≈ 57 acciones; ≈ 240 KB leídos; ≈ 80 KB escritos (sin contar los generados por script).
