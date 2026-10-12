# WF-3a · front-inicio · P-01 Inicio

Agente `front-inicio` (Opus) · 12-oct-2026 · Alcance: `frontend/src/pantallas/inicio/`, `frontend/pruebas/inicio/` (CONVENCIONES-frontend §1 y §8) y este reporte. Evalúan: U5 (U5-06, 07, 08, 09, 13), C3-03 y U3 §11.2 como guía. Sin `git add` ni `commit`.

## HECHO
- **P-01 Inicio** registrada por `src/pantallas/inicio/indice.ts` (el cascarón la descubre solo; no se tocó la base). Cuatro subsecciones con estados vacío, carga (esqueleto con la forma final y aviso a los 10 s de la base), falla con su texto exacto de U3 §7 y «Reintentar», y sin conexión («Cargado a las HH:MM» o «Sin conexión · Esta página aún no se ha cargado en este teléfono»).
- **Resumen (U3 §2)**: 6 tarjetas en orden fijo. Cada tarjeta es **un solo enlace** a su destino (P-06, P-07, P-18 de la línea, P-10, P-14, P-09; en vacío, P-26.02/.03/.05/.13 o P-15), ≥ 44 px, con foco visible. Cada una hace su propia petición. Un dato con otra forma lo atrapa un límite por bloque (`Aislado`) y las demás tarjetas siguen. Tarjeta 3: lectura activa, «Pendiente confirmar» y los estados de pago de §1.4 con los días de aviso leídos **una vez** de las alertas (K44, nunca constantes). Tarjeta 4 pide la semana ISO en curso y la tarjeta 5 la lista vigente según K38. Desde 412 px aparece el dato secundario extra (U1 §6.2).
- **Detalle (U3 §3)**: periodo Mes · 1.ª · 2.ª quincena y «Filtrar», que abre una hoja con persona, familia (K5), crédito, categoría y tienda, con chips activos y «Quitar filtro». Los paneles D-1 a D-7 aplican solo lo que admiten (§6) y dicen «Este bloque no usa el filtro …». D-3 hace dos peticiones (los vencidos solo si `is_overdue`; los próximos en `hoy..hoy+60`) y muestra las lecturas A y B lado a lado con su mínimo. D-7 (§3.1) trae la nota fija y el enlace a P-25. Al pie está «Ver todos los gastos con estos filtros».
- **Tendencias (U3 §4)**: T-1, T-2, T-3, T-4a y T-4b con ECharts 6 en SVG. Los colores se leen de los tokens (`--ui-chart-N`, `--ui-person-*`) y se repintan al cambiar el tema. Cada gráfica tiene orden fijo, marcador y trazo o patrón por serie y etiquetado selectivo (último valor y extremos). Hay un `<figure>` + `<figcaption>` con título y resumen sin variaciones calculadas, y «Ver como tabla» cambia la gráfica por la tabla (la tabla también queda para lectores de pantalla). T-2 se pinta como área apilada con los acumulados de C1, «Hoy», zona «Proyección», «Pagado por adelantado» y fila «Total» de las filas visibles, sin la palabra «deuda». En móvil la leyenda va debajo y en tableta o más, a la derecha. A 1440 px hay 2 gráficas por fila.
- **Calendario (U3 §5)**: debajo de 768 px es una agenda de hoy a hoy + 30 («Ver 30 días más» la lleva a + 60), con separador de quincena. Desde 768 px es un mes de lunes a domingo con franjas de 1.ª y 2.ª quincena, hasta 3 eventos y «+n más», con el panel del día seleccionado (a la derecha en 1440). Los pagos se cruzan con `listScheduledPayments` para el estado «Parcial». Los filtros son Tipos (varios), Crédito y Familia.
- **Karen** solo aparece en la tarjeta 3, D-3, D-7, el Calendario, las columnas de su parte en la tabla de T-2 y el filtro de familia (U3 §9).
- **Acciones de página** (U1 §2.2): «Foto de ticket» (P-05) y «Capturar quincena» (P-08 `?hoja=quincena`), esta solo el día de pago con el ingreso aún esperado. «Registrar gasto» lo pinta el cascarón. La pestaña elegida se recuerda en el dispositivo (U3 §1.1).
- Textos es-MX en `textos.ts`, dinero y fechas con `src/api/formatos.ts` y solo clases de tokens C2.

## Archivos
- `frontend/src/pantallas/inicio/`: `indice.ts`, `P01Inicio.tsx`, `Resumen.tsx`, `Detalle.tsx`, `Karen.tsx`, `Tendencias.tsx`, `GraficaInicio.tsx`, `Calendario.tsx`, `comun.tsx`, `fechas.ts`, `textos.ts`
- `frontend/pruebas/inicio/`: `inicio.escritorio.spec.ts` (corre en 412×915 táctil y 1440×900), `datos.ts`
- `docs/estado/bloqueos/WF3-inicio-vacios.md`, `docs/estado/bloqueos/WF3-inicio-tiendas.md` (nuevos)

## Verificación
- `npx tsc --noEmit`: **0 errores en `src/pantallas/inicio` y `pruebas/inicio`**. El total del proyecto da 1 error ajeno: `pruebas/creditos-salud-reportes/datos.ts(48,79)` TS6133.
- `npx vite build`: «built in 612ms», PWA con 82 entradas. `npm run build` y `make verify-wf3-frontend` fallan **solo** por ese error ajeno de tsc.
- `npx playwright test pruebas/inicio pruebas/shell.spec.ts pruebas/cola.spec.ts pruebas/login.spec.ts` (con FH_CHROMIUM): **24 passed**, de ellas 14 de Inicio (7 × 2 tamaños). Las pruebas base siguen verdes con P-01 registrada.
- **U5-06**: 6 tarjetas (≤ 6), orden 1→6 comprobado, cada tarjeta = 1 enlace y 0 controles dentro, alto ≥ 44 px, sin desplazamiento horizontal a 412 y 1440. `grep` de hex en mi carpeta = **0**. Medidas arbitrarias `[..px]`/`style=` = **0**.
- **U5-07**: `grep -rcE "Save|Cancel|Settings|Dashboard|Loading|Error|Home|Delete|Search|Add" src/pantallas/inicio` = **0** (sensible a mayúsculas). En la prueba, el texto de `main` no tiene esas palabras ni fechas ISO. Las fechas se ven como `10-oct-2026` y `Vie 30 oct 2026`, el dinero como `$1,234.56`.
- **U5-08**: `npm run verificar:hosts` da «U5-08 dist: hosts externos que se pedirían = 0».
- **U5-09**: «U5-09 toques desde Inicio (412 px): 3» y «(1440 px): 3». Después de guardar, Inicio vuelve a pedir `month-vs-plan` y `fortnight-flow` (F-1 paso 4).
- **U5-13**: sin cambios a la cola; `cola.spec.ts` pasa (3 toques sin conexión y sin duplicar).
- **C3-03**: `npm run verificar:rutas` da «201 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0». Todas las llamadas van por `usarConsulta` con la ruta literal del contrato.
- **U3 §11.2** cubiertos por prueba: A-01, A-02 (alertas 1 vez; `week-meals?period=2026-W41`, `best-basket?week=2026-W41`), A-03/A-04 (Sáb 10 oct), A-05 (lectura B), A-06, A-07 (falla y Reintentar en la tarjeta 1, con las demás visibles), A-08, A-09 (412/1440), A-11 parcial (tabla T-2: R-45 «$1,947.50 $973.75 $973.75» y R-46; capas $3,895.00 y $2,921.25; figcaption; sin «deuda»), A-12 parcial (franjas de octubre), A-16 (3 pagos, sin petición de vencidos, `period=2026-10-10..2026-12-09`), A-17 (Jue 29 oct «Vence mañana»). F-3 y F-4 paso 1: tocar las tarjetas 5 y 3 lleva a `/compras/comparar` y `/creditos/1`.
- **S6 / S7**: no aplican a esta pantalla (jobs, respaldos y Telegram son del servidor). No se envió nada por Telegram.
- `make verify-all GOOSE=…` da exit 2 en `verify-wf1-cobertura`, por `internal/dominio` 75.2 % (< 85 %): es del backend y ajeno a este alcance.

## SUPUESTOS
- PROVISIONAL (bloqueo `WF3-inicio-vacios`): las reglas de vacío de las tarjetas 2, 4 y 6, D-6 con todo lo activo, los nombres de filtros hacia P-02, los destinos de eventos del Calendario, el texto de proyección apagada y la condición de «Capturar quincena».
- PROVISIONAL (bloqueo `WF3-inicio-tiendas`): el color de tienda se saca del nombre de la cadena.
- Las pruebas usan respuestas simuladas con las cifras del oráculo de U3. Ingresos, plan y ahorro son marcadores (U3 los deja así).
- En escritorio los filtros de Detalle también se abren en la hoja (no hay una barra en línea).

## BLOQUEOS
- `docs/estado/bloqueos/WF3-inicio-vacios.md` (8 preguntas, todas con opción PROVISIONAL aplicada).
- `docs/estado/bloqueos/WF3-inicio-tiendas.md` (código de tienda en C1).

## PENDIENTES
- Faltan pruebas de A-10 y A-09 a 360, 768 y 1024 px. También de A-12 en febrero de 2027, de A-15 (axe en claro y oscuro), de T-1 con datos, de «Pagado por adelantado», de D-4 y D-5 con datos y del aviso de 10 s (el aviso es de la base).
- En T-2 los meses 7 a 12 se marcan con una zona «Proyección» y la nota, pero el trazo no es discontinuo por tramo. Falta partir la serie.
- No está el mes compacto opcional en 412 ni los `mocks.ts` de msw del área (los dos son opcionales).

## PROPUESTAS PARA AXEL
- Responder las 8 preguntas de `WF3-inicio-vacios` (lo que más importa es la 1: cómo se sabe que falta el ingreso quincenal).
- Aprobar `Store.code` en C1 (opción B de `WF3-inicio-tiendas`), que también lo necesita Compras.

## Paquetes instalados
Ninguno (se usaron los de frontend-base: ECharts 6, lucide, Radix, Playwright).

## Consumo aproximado
~48 acciones · ~220 KB leídos · ~130 KB escritos.
