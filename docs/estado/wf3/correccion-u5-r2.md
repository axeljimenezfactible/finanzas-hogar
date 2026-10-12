# WF-3b · correccion-u5-r2 · Reporte

Agente `correccion-u5-r2` (Opus) · 12-oct-2026 · HEAD `6d691d0` + árbol de trabajo. Corrige los 5 hallazgos altos y medios de `REVISION-u5-r2.md` (N-1, N-2, N-3, M-6b, M-1b). No toqué los bajos (B-1 a B-9) ni reabrí decisiones. No hice `git add` ni commit.

## Hecho
1. **N-1 · U5-07 · nota de la semilla en P-26.13 (alta).** En `backend/internal/semilla/desarrollo.go`, `devNote` pasa a «Dato de desarrollo: no existe en producción.» y es la nota única de las filas del overlay. Antes, 4 filas le sumaban « Fixture <archivo>, <clave>[n].»: presupuesto, partida planificada, meta de ahorro y aportación. Ahora ninguna lo hace. El origen de cada fila sigue en la tabla del encabezado del archivo y en `semilla-dev.md`. El cargador es idempotente y no reescribe las notas de las filas que ya existen, así que **volví a crear la base de desarrollo desde cero**: `dev-stop` → respaldo temporal → `finanzas migrate` → `finanzas semilla` (11 «cotejo OK») → `finanzas semilla --desarrollo` (sin `--con-gasto-fijo`, como la base anterior: 1 presupuesto, 1 partida, 1 meta, 1 aportación, 36 precios y 12 comidas, todos nuevos) → `dev-run`. El respaldo se borró al terminar. Las sesiones anteriores se perdieron (es solo desarrollo).
2. **N-2 · U5-06 #12 · capturas Android sin puntero grueso (media).** En `pruebas-u5/capturas.spec.ts`, cada pantalla abre una **página nueva** en el contexto con sesión y la cierra al terminar. La cookie es del contexto, así que se sigue entrando una sola vez. P-00 sigue en su contexto anónimo. Antes de capturar se lee `matchMedia('(any-pointer: coarse)').matches` y se guarda como `puntero_grueso` en `manifiesto.json`. En Android se exige `true`.
3. **N-3 · U5-06 #7 · sombra en reposo de las secciones de Ajustes (media).** Quité `shadow-1` del `<section>` de `Bloque` en `ajustes-a/comunes.tsx` y `ajustes-b/comun.tsx`, que queda en `shadow.0`. La sección conserva la separación por superficie (`ui.surface.card` sobre `ui.surface.page`) y el borde `ui.border.divider` (§4.5 y §7.13). Las perillas de los interruptores y los diálogos (`shadow-2`) siguen igual porque flotan.
4. **M-6b · U5-06 #14 · etiquetas encimadas en P-01 ▸ Tendencias (media).**
   - `inicio/GraficaInicio.tsx` (`sinEncimar`, usado por todas las `FiguraGrafica`): en las series de línea, la etiqueta del primer punto se alinea a la izquierda y la del último a la derecha. «$18,025.47» ya no invade las etiquetas del eje Y y el valor extremo sigue visible (§8.5).
   - `inicio/Tendencias.tsx` (T2 «Devengado contra efectivo»): quité las etiquetas de fin de serie que se dibujaban sobre el relleno. Los valores finales pasan a la **leyenda** con formato completo: «Efectivo acumulado $19,397.22» y «Diferencia acumulada $0.00» (§8.5, «Formato completo en etiquetas y leyendas»).
   - En T2 la leyenda queda **abajo en los dos tamaños**, y la rejilla le deja 40 px en escritorio y 64 px en Android, donde la leyenda ocupa 2 renglones. En escritorio, la leyenda a la derecha (margen de 120 px) caía sobre el área y su marca del mismo color desaparecía. Esa era la «leyenda sin marca» en Ec.
   - El devengado acumulado sigue en el resumen de la figura y en «Ver como tabla».
5. **M-1b · U5-06 #2 · título de P-00 (media).** En `app/IniciarSesion.tsx`, el h1 pasa a `text-title-1 font-bold` (32/700).
6. **Script de U5-07.** En `capturas.spec.ts`, cada captura revisa su `innerText`:
   - (a) cada importe en formato `$1,234.56`. El compacto «$1.5 mil» solo se permite en ejes (§7.7).
   - (b) 0 fechas `[0-9]{4}-[0-9]{2}-[0-9]{2}`.
   - (c) 0 palabras de la lista del criterio, más «Fixture», «Submit», «Close» y la descripción automática de ECharts.
   - (d) 0 claves internas `\b[a-z]+_[a-z_]+\b`.
   - El resultado queda en `manifiesto.json` (`u5_07`). Los puntos (b), (c) y (d) son aserciones. `resumen_u5.py` imprime las líneas U5-07 y la del puntero grueso.
7. **Comprobación nueva de Tendencias.** `pruebas-u5/tendencias.spec.ts` (`--grep Tendencias`) captura `/inicio?sub=tendencias` en las 4 combinaciones dentro de `evidencia/tendencias/`, una subcarpeta que no cuenta para U5-01. Compara las cajas de todos los `<text>` del SVG de ECharts de cada figura (sin los iconos de los chips): 0 encimadas y 0 fuera del lienzo. Cada caja se reduce 20 % arriba y abajo para descontar el interlineado: el «$0» del eje Y y el primer mes del eje X se tocaban 3 px sin que las letras se tocaran (lo revisé en la imagen).
8. Recompilé (`npm run build`), reinicié el servidor y **regeneré las 164 capturas** con axe y contraste, porque el cambio de método afecta las 82 de Android y la nueva base afecta a todas.

## Archivos
- Aplicación: `backend/internal/semilla/desarrollo.go` · `frontend/src/app/IniciarSesion.tsx` · `frontend/src/pantallas/ajustes-a/comunes.tsx` · `frontend/src/pantallas/ajustes-b/comun.tsx` · `frontend/src/pantallas/inicio/GraficaInicio.tsx` · `frontend/src/pantallas/inicio/Tendencias.tsx` · `frontend/dist/` (recompilado).
- Pruebas: `frontend/pruebas-u5/capturas.spec.ts` · `frontend/pruebas-u5/tendencias.spec.ts` (nuevo) · `frontend/pruebas-u5/resumen_u5.py` · `frontend/playwright.u5.config.ts` (agrega `tendencias` al `testMatch` de los 4 proyectos y la línea de uso).
- Evidencia regenerada: `docs/design/capturas-u5/` (164 PNG, `axe/` 164, `contraste/` 164, `hosts.json`, `evidencia/manifiesto.json`, `evidencia/tema-*.json`, `evidencia/cifras-oraculo.json`) y, nueva, `evidencia/tendencias/` (4 PNG y 4 JSON).
- Base `datos/finanzas-dev.db`, recreada (ignorada por git). Este reporte.

## Verificación
- **tsc**: `npx tsc --noEmit` da 0 errores. Los dos specs de `pruebas-u5` también pasan con `--strict`.
- **Build**: `npm run build` sale bien (PWA con 82 entradas). Backend: `go vet ./internal/semilla` y `go test ./...` dan OK, sin ningún FAIL.
- **Pruebas Playwright afectadas** (`pruebas/ajustes-a`, `pruebas/ajustes-b`, `pruebas/inicio`, `login.spec.ts`, `medidas.spec.ts`): 125 pasan y 1 se omite. La omitida es `ajustes.escritorio.spec.ts:372`, que se salta sola en ese proyecto.
- **Corrida U5**: `npx playwright test -c playwright.u5.config.ts`, partida en Android y escritorio, da **87 + 87 pasan**: 164 capturas, 4 comprobaciones de T, 2 del interruptor y 4 de Tendencias.
- **U5-01**: `python3 -I frontend/pruebas-u5/verificar_capturas.py` da `T = 41 · 4T = 164 · PNG correctos: 164 de 164 · faltantes: 0 · mal: 0`. Los 164 PNG son posteriores a la nueva base. Sigue abierto el bloqueo `U5-conteo-carpeta`: `ls | wc -l` = 168.
- **U5-02**: las cifras de la semilla se ven en $1,191.50 (P-02, P-04, P-19, P-25), $10,000.00 (P-17, P-18), $3,094.50 (P-01, P-02, P-03, P-04, P-14), $6,293.44 (P-09, P-26.13), $6,905.50 (P-01, P-03, P-17, P-18, P-26.05) y $973.75 (P-01, P-03, P-17, P-18, P-24, P-26.05). Ninguna captura quedó cargando ni en la pantalla de error.
- **U5-03**: `resumen_u5.py` da **0 serias o críticas y 0 violaciones de otro impacto** en las 4 combinaciones (incompletos: 15 · 15 · 6 · 6). Incluye `target-size` (WCAG 2.2) ya con el puntero grueso.
- **U5-04**: **0 pares bajo umbral** en las 4 combinaciones (4,490 · 4,511 · 5,714 · 5,735 pares evaluados) y 0 indeterminados.
- **U5-06 #12 (N-2)**: en `manifiesto.json`, `puntero_grueso = true` en **82 de 82** capturas Android, con aserción en cada una. Medida directa en P-02 Android: 1 control de menos de 44 px fuera de una etiqueta (antes, la página reutilizada daba controles de 36 px).
- **U5-06 #7 (N-3)**: `box-shadow` computado de los `section` de P-26.01, .05, .10, .13 y .14 da `none` en las 4 combinaciones.
- **U5-06 #14 (M-6b)**: `--grep Tendencias` pasa 4 de 4, con 0 textos encimados y 0 fuera del lienzo en T1, T2 y T3. T4 muestra sus estados vacíos. Lo revisé a ojo en Ec y Ao:
  - «$18,025.47» empieza en el primer punto, a la derecha de «$18 mil» y «$21 mil».
  - T2 no tiene texto sobre el área, y la leyenda de abajo muestra las 2 marcas (círculo y cuadro) con sus valores finales.
- **U5-06 #2 (M-1b)**: h1 de P-00 «Finanzas Hogar» a **32px/700** en las 4 combinaciones.
- **U5-07 (N-1)**: en las 4 combinaciones hay **351 importes, 0 fuera de formato, 0 fechas ISO, 0 palabras en inglés (con «Fixture») y 0 claves internas**. P-26.13 muestra «Dato de desarrollo: no existe en producción.» En la base no queda ninguna nota con «Fixture» ni con `_`.
- **U5-08**: 1 host único, `finanzas.local:8443`, en las 164 capturas.
- **U5-10**: 0 capturas con el fondo distinto del token de C2 o con `data-theme`. El interruptor pasa en Android y escritorio (`tema.spec.ts`).
- **C3**: no se tocó ningún contrato. La semilla real se volvió a cargar con sus 11 cotejos «OK» del oráculo.

## Supuestos
1. Dejé la misma nota simple en las 4 filas del overlay que llevaban «Fixture», no solo en la meta, porque es el mismo defecto: las notas del presupuesto y de la partida también pueden verse en Presupuesto y en Ajustes.
2. Para volver a sembrar había que recrear la base: el cargador hace «ensure» y no actualiza notas. Usé el procedimiento de `semilla-dev.md` y de `semilla-c3.md`, sin `--con-gasto-fijo`.
3. Alineé la primera y la última etiqueta en lugar de quitarlas, para que el valor extremo siga a la vista (§8.5: «líneas con el último valor y los extremos»).
4. En las áreas apiladas, los valores finales van en la leyenda, como sugirió el revisor. «Devengado acumulado $…», que era la etiqueta de la cima, ya no se dibuja: el mismo dato sigue en el resumen de la figura y en la tabla.
5. En U5-07, los importes se registran pero no son aserción. Así lo pide el criterio: el revisor acepta el compacto de ejes de §7.7.
6. La comprobación de Tendencias reduce cada caja de texto un 20 % arriba y abajo (ver Hecho 7). Las cajas completas siguen en el JSON.

## Bloqueos
Ninguno nuevo. Siguen abiertos `U5-conteo-carpeta`, `U5-nombres-propios` y `U5-tipografia-tres-tamanos`, que no son de esta corrección.

## Pendientes
- El integrador hace el commit. Después viene la ronda 3 del revisor U5.
- Siguen los bajos B-1 a B-9 de la ronda 2, que quedaron fuera de alcance.
- En Android, la captura de página completa sigue dejando la barra inferior fija a media imagen. En Tendencias tapa «Hoy» y «Proyección» de T2. Es un artefacto de `fullPage`, no de la app.
- `manifiesto.json` marca `h1_coincide: false` en P-00 porque compara «Finanzas Hogar» con el nombre de la ruta. Ya pasaba antes y no es un hallazgo de la revisión.

## Propuestas para Axel (PROPUESTA, no decidido)
| Qué | Resultado | Costo | Riesgo |
|---|---|---|---|
| Llevar la leyenda de todas las gráficas de Tendencias abajo, como T2, o calcular el margen derecho según el ancho de la leyenda | Ninguna leyenda cae sobre el trazo en escritorio (con «Deuda viva total» casi se toca) | Bajo | Las gráficas de escritorio quedan un poco menos altas |
| `ops/dev-run.sh --sembrar` corre también `semilla --desarrollo` y avisa que, si hay que rehacer notas, se borra la base | Volver a sembrar es un solo comando | Bajo | Ninguno (solo desarrollo) |
| Agregar «Tendencias» de P-01 y de P-10 como evidencia fija de U5, como ya hace `evidencia/tendencias/` | Las gráficas quedan revisadas en cada corrida | Bajo | Ninguno |
Recomendación: la primera, en la siguiente corrección de frontend.

## Consumo aproximado
≈ 52 acciones de herramienta · ≈ 120 KB leídos · ≈ 22 KB escritos a mano (código ≈ 9 KB y este reporte ≈ 11 KB). Además, la corrida generó de nuevo ≈ 24 MB de PNG y JSON.
