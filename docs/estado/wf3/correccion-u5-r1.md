# WF-3b · correccion-u5-r1 · Reporte

Agente `correccion-u5-r1` (Opus) · 12-oct-2026 · HEAD `6d691d0` + árbol de trabajo. Corrige los 11 hallazgos altos y medios de `REVISION-u5-r1.md` (A-1 a A-5, M-1 a M-6). No toqué los bajos (B-1 a B-6), `docs/design/sistema/`, `docs/design/maquetas/`, `docs/estado/wfp3/` ni el backend.

## Hecho
1. **A-1 · U5-03 (axe).** `BarraAvance` exige `nombre` y lo pone como `aria-label` del `role="meter"`; el valor sigue en `aria-valuetext`. Nombres: «Real», «Quincena n», la categoría de la fila (P-06), el nombre de la meta y «Aportado este mes» (P-09). En `VistaTabla`, una columna sin nombre (las acciones de P-08) ya no lleva botón de orden, solo «Acciones» para lectores de pantalla. Las columnas con nombre conservan su texto visible como nombre accesible. Le había puesto `aria-label="Ordenar por …"`, pero eso hacía que `getByLabel('Monto')` encontrara dos elementos en F-1, así que lo quité. En P-26.05 y P-26.06, la columna Nombre ya no pinta un `<button>` dentro del botón de la fila.
2. **A-2 · U5-01 (P-26.13 a 415 px).** La barra de vista de `VistaBD` pasa a `flex-wrap min-w-0 max-w-full`.
3. **A-3 · U5-04 (componentes de UI).** La barra de vista ahora es un grupo de alternados (`role="group"`, `aria-pressed`, DESIGN §7.5.1) y marca la vista elegida con un indicador inferior de 2 px en `border-border-selected`. Las pestañas Gasto/Ingreso/Peso/Ticket de P-03 siguen siendo `role="tab"` y llevan el mismo indicador. Pasan a `border-border-control`: los chips seleccionables sin elegir, las tarjetas-botón (`TarjetaMetrica` con `alAbrir`), las tarjetas de galería y los botones de P-10 (tiempos de comida y días). Los chips de persona y tienda, que no se pueden pulsar, se quedan con `border-divider`.
4. **A-4 · U5-07 (inglés).** Nuevo `componentes/SelectorArchivo.tsx`: el input nativo queda oculto y lo abre un `Boton` «Elegir archivo CSV» o «Elegir archivo», con el nombre del archivo o «Ningún archivo elegido». Se usa en P-26.11 y P-26.06. En P-26.04, la columna «Clave que usan los cálculos» muestra «Despensa» y «Comida fuera» (equivalencias de U4 §P-26.04).
5. **A-5 · U5-13 (H-5).** Causa: después de entrar con el formulario nunca se guardaba la copia de `/api/v1/auth/me`; solo se guardaba con una navegación completa. Al recargar sin red no había sesión guardada y salía la pantalla de entrada. Ahora `entrar()` guarda la sesión recibida como última lectura de `/auth/me` (`app/sesion.tsx`). Al cerrar sesión se borra con `borrarLecturas`, igual que antes.
6. **M-1.** El `h1` de página usa `text-title-1 font-bold` (32 px/700, `shell/Barras.tsx`).
7. **M-2.** La barra lateral mide `w-60` (240 px, fija y flotante). El encabezado reserva `mr-12` para el botón de plegar, con 16 px de separación.
8. **M-3.** Los iconos Lucide de `size-3` y `size-4` pasan a `size-5` (20 px, trazo 1.5) en 21 archivos. Los `size-6` sin `size={ICONO_MD}` lo reciben. Los iconos de `Vacio` llevan `stroke-[1.5]`. Además, `svg.lucide { flex-shrink: 0 }` en `estilos/app.css` evita que se deformen. Quité los `mt-0.5` de los iconos.
9. **M-4.** Cambié `gap-0.5` por `gap-1` en la lista de la barra lateral, en los destinos de la barra inferior (de ahí salían los 2 px que se veían en Android), en P-10 y en `inicio/comun.tsx` y `Calendario.tsx`. `Vacio` pasa de `py-10` a `py-12`. El relleno inferior del contenido en móvil era 64 px y ahora es `calc(var(--spacing)*14 + env(safe-area-inset-bottom))`: exactamente `layout.bottombar.height` (§4.3), una medida de estructura.
10. **M-5.** En `Franja` (`componentes/Estados.tsx`), por debajo de 768 px el texto y la acción se apilan a todo el ancho. Aplica a P-17, P-18 y a cualquier franja con acción.
11. **M-6.** Si todas las series valen 0 o no tienen valor, `Grafica` y `FiguraGrafica` muestran «Aún no hay datos para este periodo» con el título y sin ejes. En `FiguraGrafica` cambian tres cosas en las etiquetas:
    - Todas las series llevan `labelLayout` con `hideOverlap` y `moveOverlap: 'shiftY'`. Antes estaba en la raíz de la opción, donde ECharts lo ignora.
    - Cuando hay leyenda, las etiquetas de fin de línea pasan a ella.
    - Las etiquetas de columnas de menos de 72 px se quitan; el valor queda en la descripción emergente y en «Ver como tabla» (DESIGN §8.5).

    En `Grafica`, la rejilla deja 32 px para la leyenda (`type: 'scroll'`).

## Archivos
- Componentes y shell: `frontend/src/componentes/{SelectorArchivo.tsx (nuevo), Chip.tsx, TarjetaMetrica.tsx, Estados.tsx, Grafica.tsx, Boton.stories.tsx, TarjetaMetrica.stories.tsx}`, `frontend/src/shell/{Barras.tsx, BarraLateral.tsx, HojaRegistro.tsx, Layout.tsx}`.
- Vistas, app, textos y estilos: `frontend/src/vistas-bd/{VistaBD.tsx, VistaTabla.tsx, VistaGaleria.tsx}`, `frontend/src/app/sesion.tsx`, `frontend/src/textos/es-MX.ts`, `frontend/src/estilos/app.css`.
- Pantallas, gastos e inicio: `gastos-presupuesto/{comun.tsx, P06Presupuesto.tsx, P09Ahorro.tsx}`, `inicio/{GraficaInicio.tsx, comun.tsx, Calendario.tsx}`, `comidas-compras/P10Plan.tsx`.
- Pantallas de ajustes: `ajustes-a/{S04Categorias.tsx, S05Creditos.tsx, S06Tiendas.tsx, textos.ts, P26Ajustes.tsx, S01Cuenta.tsx, S02Presupuesto.tsx, S07Recetas.tsx, comunes.tsx}` y `ajustes-b/{P2611Datos.tsx, P2608Reportes.tsx, P2610Apariencia.tsx, P2612Personas.tsx, P2614Inventario.tsx, comun.tsx}`.
- Pantallas de créditos, salud y reportes, solo iconos: `creditos-salud-reportes/{P17Lineas.tsx, P18Credito.tsx, P19Compra.tsx, P21Salud.tsx, P23Centro.tsx}`.
- Pruebas: en `frontend/pruebas/{ajustes-a/ajustes.escritorio.spec.ts, ajustes-b/ajustes-b.escritorio.spec.ts, gastos-presupuesto/correccion-bis.escritorio.spec.ts}`, el localizador de vista cambia de `getByRole('tab', { name: 'Tablero'|'Galería' })` a `getByRole('button', { …, exact: true })`.
- Evidencia regenerada por script: `docs/design/capturas-u5/` (164 PNG, `axe/`, `contraste/`, `hosts.json`, `evidencia/`) y `frontend/pruebas-u5/flujos/{resultados,capturas}` de U5-13.

## Verificación
- tsc y build: `npm run build` → `tsc --noEmit` sin errores y `vite build` correcto (PWA, 82 entradas). Servidor reiniciado con `ops/dev-stop.sh` + `ops/dev-run.sh` (`/health` 200).
- **U5-01:** `python3 -I frontend/pruebas-u5/verificar_capturas.py` → `PNG correctos: 164 de 164 · faltantes: 0 · mal: 0`. En el DOM, P-26.13 Android mide `innerWidth/scrollWidth` 412/412, y lo mismo en las 23 pantallas medidas.
- **U5-03:** `npx playwright test -c playwright.u5.config.ts` (corrida completa, 170 passed, 8.0 min; luego otra de P-02/06/08/09/26.05/26.06/26.13 tras el último cambio, 28 passed) y `resumen_u5.py` → **0 serias/críticas** en las 4 combinaciones (incompletos 15/15/6/6).
- **U5-04:** `resumen_u5.py` → **0 pares bajo umbral** en Ac, Ao, Ec y Eo (4,742 / 4,749 / 5,952 / 5,966 evaluados; 0 indeterminados), incluidos los subtipos `seleccion` y `boton`.
- **U5-06:** medición temporal en el DOM sobre 23 pantallas, en Android claro y escritorio oscuro:
  - `h1` a 32px/700.
  - Barra lateral de 240 px; encabezado en x 8–183 y botón de plegar en x 199–231 (16 px de separación).
  - 0 iconos fuera de 20×20 o 24×24 con trazo dibujado de 1.50.
  - Franja de P-17 en Android con el texto a 328 px en 2 líneas.
  - Tendencias: P-01 tiene 2 estados vacíos y 3 gráficas con datos (cada una con un solo «$0», el arranque del eje); P-10 tiene 2 estados vacíos.
  - Espaciado fuera de la retícula: solo quedaron los 2 px de la barra inferior (ya corregidos y en la corrida final de capturas, sin volver a medir en el DOM) y los 56 px de estructura.
- **U5-07:** P-26.04 ya no muestra «groceries» ni «eating_out». P-26.11 tiene 0 inputs de archivo visibles, el botón «Elegir archivo CSV» y el texto «Ningún archivo elegido».
- **U5-08:** `hosts.json` registra 1 host, `finanzas.local:8443`, en las 164 capturas. **U5-10:** 0 capturas con fondo distinto de C2 o con `data-theme`.
- **U5-13:** con `servidor-prueba.sh` en :8444, `captura-sin-conexion` da **5 passed**. La sesión empieza desde la primera carga, con el formulario y sin navegación completa (la secuencia de H-5). Paso 2b tras recargar sin red: «Sin conexión con el servidor de casa · 2 en cola». Al volver la red, llegan 1 gasto y 1 peso, y el reintento con la misma clave da `Idempotent-Replayed: true`.
- Pruebas Playwright del frontend: `npx playwright test --workers=2` → **241 passed, 4 skipped, 0 failed**. En la corrida anterior fallaban 2 (F-1, `getByLabel('Monto')` contra el `aria-label` de orden, ya quitado).
- C3: no reejecuté C3-09 (no cambié cifras ni API).

## Supuestos
1. La barra de vista deja de ser `tablist` y pasa a alternados con `aria-pressed`, como pide DESIGN §7.5.1. Por eso cambié los localizadores de 3 pruebas de otras áreas. Las pestañas de P-03 y de subsección siguen como `tab`.
2. Las tarjetas-botón y los chips seleccionables llevan borde `border-control`, la corrección que sugiere el revisor. Si Axel resuelve el bloqueo `U5-contraste-ui-controles` excluyendo los botones con texto, basta con revertir la clase en `Chip.tsx`, `TarjetaMetrica.tsx`, `VistaGaleria.tsx` y `P10Plan.tsx`.
3. En las claves de P-26.04, un código que no sea `groceries` ni `eating_out` se muestra tal cual. C1 hoy solo tiene esos dos.
4. Estado vacío con `py-12` (48 px), no con `py-8`. Relleno inferior en móvil igual a `layout.bottombar.height` (56 px + área segura) en lugar de 48 px, para que la barra fija no tape contenido.
5. Ancho mínimo de 72 px por columna para su etiqueta directa: «$1,197.64» a 12 px más holgura.

## Bloqueos
Ninguno nuevo. Siguen abiertos para Axel `U5-contraste-ui-controles` (hoy da 0 pares bajo 3:1 con el borde subido) y `U5-nombres-propios` (claves internas ya en español; siglas y nombres de B-2 pendientes).

## Pendientes
- Que el revisor vuelva a correr la ronda. No revisé a ojo los PNG de Tendencias de P-01 en Ao y Ec: no comprobé que las etiquetas ya no se encimen.
- Bajos B-1 a B-6 sin tocar (fuera de alcance).
- `git add` y commit los hace el integrador. Temporales borrados (`datos/tmp/corr-u5`, `u5-flujos`, `u5-partes`, `playwright-*` de mis corridas, configuración y spec temporales).

## Propuestas para Axel (PROPUESTA, no decidido)
| Opción | Resultado | Costo | Riesgo |
|---|---|---|---|
| Dejar el borde `border-control` en chips y tarjetas-botón | U5-04 en 0 sin depender del bloqueo | Ninguno, ya aplicado | Las tarjetas del tablero se ven con borde más marcado |
| Volver al borde suave y excluir los botones con texto (lectura B del bloqueo) | Diseño más ligero | Revertir 4 clases | Depende de que el revisor acepte la lectura B |

Recomendación: dejar el borde `border-control`.

## Consumo aproximado
≈ 58 acciones · ≈ 120 KB leídos · ≈ 30 KB escritos a mano (código ≈ 18 KB, reporte ≈ 9 KB, scripts temporales ≈ 7 KB ya borrados). Las corridas generaron además ≈ 24 MB de PNG y JSON.
