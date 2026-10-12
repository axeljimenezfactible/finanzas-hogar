# REVISIÓN · frontend de áreas · ronda 1 de 3 (WF-3a, 12-oct-2026)

Revisor independiente: `revisor-frontend-areas-r1`. No escribí nada de lo revisado y no corregí nada.
Agentes revisados: `front-creditos-salud-reportes`, `front-ajustes-a`, `front-ajustes-b` y `front-variantes`.
Contra qué se evalúa: U5-06, U5-07 y U5-09; C3-03; DESIGN.md; U1 (P-17 a P-25, §2.2 y §3.19 a §3.27); U4 (§0 y P-26.01 a P-26.14).
Todas las cifras de este informe salen de comandos que corrí yo; ninguna se toma de los reportes.

## Veredicto: NO APROBADO

Quedan **4 altos, 8 medios y 5 bajos**.

- Lo que sí cumple:
  - Build y tsc.
  - Las 60 pruebas Playwright.
  - C3-03.
  - U5-07, medido en pantalla.
  - Las claves `reports.*` y `alerts.*`.
  - Cero colores hex en el código de las áreas.
  - Contraste del texto en las pantallas.
- Lo que falla:
  - Dos vistas previas que U4 §0.2 exige: la frecuencia de pagos en P-26.05 y el arrastre en el tablero de P-26.14.
  - U5-06 en componentes de la base que usan estas pantallas: objetivos táctiles de menos de 44 px, colores de ECharts fuera de C2 y dinero sin cifras tabulares.
  - Un dato de prueba que no valida contra el esquema de C1.

## Verificación propia (comandos y salidas resumidas)

Los temporales quedaron en `datos/tmp/rev-fa-r1/` y ya se borraron.

| Criterio o comprobación | Estado | Evidencia mía |
|---|---|---|
| tsc | CUMPLE | `npx tsc --noEmit` sale 0, con 0 errores. El error TS6133 que citaban ajustes-a y ajustes-b ya no está. `npm run typecheck:variantes` sale 0. |
| Build | CUMPLE | `npx vite build --outDir datos/tmp/rev-fa-r1/dist` dice «✓ built». Están los chunks P17…P25, P26Ajustes, S01…S07 y P2608…P2614. |
| Build de Storybook | CUMPLE | `npm run build-storybook` dice «Storybook build completed successfully». |
| Playwright de las 3 carpetas | CUMPLE | Usé una configuración propia: preview de mi build en el puerto 4795 y Chromium 1217 local. `playwright test ajustes-a ajustes-b creditos-salud-reportes` da **60 passed**: 20, 26 y 14. |
| **C3-03** | CUMPLE | `npm run verificar:rutas` da 201 rutas invocadas: 0 fuera de openapi.yaml y 0 sin registrar (216 registradas). Además hice un cruce independiente de cada literal `/api/…` contra openapi.yaml y routes.go: creditos-salud-reportes tiene 39 llamadas, ajustes-a 58 y ajustes-b 47. En las tres: 0 fuera de contrato, 0 sin registrar, 0 `fetch` y 0 llamadas con ruta no literal. |
| **U5-07** | CUMPLE (en lo medido) | Escaneé el `innerText` en ejecución: 33 rutas (P-17 a P-25 con sus pestañas y Ajustes 01 a 14), en claro y en oscuro, en 412 táctil y en 1440. Son 132 páginas. Resultado: 0 fechas ISO, 0 palabras en inglés de una lista de 45, 0 identificadores `snake_case`, 0 montos sin centavos o sin coma, y 0 «-$» con guion en lugar del signo menos (U+2212). |
| **U5-06** | **NO CUMPLE** | Lo que pasa: 0 hex, 0 `rgb/hsl/oklch`, 0 `dark:` y 0 valores arbitrarios en las 3 carpetas, las historias y `.storybook/`; solo hay 3 hex, y son texto dentro de `hallazgos.json`. El contraste del texto en las 132 páginas da 0 pares bajo 4.5:1 (o 3:1 en texto grande). Lo que falla: H-03, H-09 y H-10. |
| Variantes con contraste | CUMPLE el script; los componentes **NO** | `node scripts/contraste-variantes.mjs` da: 37 historias, 148 páginas y 5 230 mediciones; 0 fallas nuevas y 0 hosts externos. Pero **5 de 5 hallazgos FV se reproducen**, y la salida es 0 aunque existan. |
| U5-09 | No aplica | Ninguna pantalla revisada registra gastos. P-03 y la cola son de frontend-base. |
| Vista previa «Antes → Después» (U4 §0.2) | **NO CUMPLE** | La regla está bien implementada en `comunes.tsx:145`: no guarda si la vista previa es nula o dio error. Fallan H-01 y H-02. Quedan parciales H-05 y H-07. |
| Claves `reports.*` y `alerts.*` | CUMPLE | `P2608Reportes.tsx:29-34` usa `PATCH /api/v1/settings/reports` con `parameters` y `reset_parameters`. `P2609Notificaciones.tsx:27` hace lo mismo con `/settings/alerts`. Las 11 claves del código son exactamente las 11 sembradas en `backend/db/migrations`. |
| P-17 a P-25 contra U1 | CUMPLE PARCIAL | De 59 textos y elementos de U1 §2.2 y §3.19 a §3.27, faltan 3: «Comparador no disponible» (H-13), «Restricciones alimentarias» en la ficha de P-21 (H-12) y los filtros por tipo y periodo de P-23, que solo existen como el filtro de texto genérico de VistaBD (aceptado). |
| Datos de prueba contra C1 | **NO CUMPLE** | Capturé las respuestas de la API falsa (`page.on('response')`) y las validé con ajv contra openapi.yaml 3.1: 55 respuestas distintas, 54 válidas y 1 inválida (H-04). |
| Credenciales | CUMPLE | Busqué token, secret, contraseñas literales y el patrón del token del bot en las 3 carpetas, sus pruebas, `.storybook` y el script: 0. El `chat_id` se muestra enmascarado. |
| `make verify-all` | Falla por algo ajeno | Sale 2 en `verify-go`: «FALLA gofmt: internal/reportes/…». El paquete `backend/internal/reportes/` está sin rastrear y lo escribe otro agente del backend. No toca nada del frontend. |

### Ajustes: tabla sección → campos (U4 contra el código)

Un script mío leyó cada tabla de campos de U4 y buscó cada etiqueta en el código del área. Para las que no aparecían tal cual, busqué a mano; cuando dice «verificadas a mano», todas aparecieron.

| Sección | Etiquetas de U4 | Encontradas tal cual | Verificadas a mano | Faltan |
|---|---|---|---|---|
| 01 Cuenta y roles | 17 | 14 | 3 (cerrar sesión, cerrar por alcance, parámetros) | 0 |
| 02 Presupuesto | 20 | 15 | 5 (monto por categoría, alcance, «Solo este mes», quitar, V-17) | 0 campos; falta `Ctrl+/` (H-07) |
| 03 Ingresos | 7 | 7 | — | 0 |
| 04 Categorías | 8 | 8 | — | 0 campos; faltan el arrastre con «Aplicar» y `Ctrl+/` (H-07) |
| 05 Créditos y condiciones | 52 | 48 | 4 (eliminar tasa, eliminar promoción, pesos, parámetros) | Alta de promoción (`createPromotion`, H-06); vista previa 6 (H-01); vista previa 7 incompleta (H-05) |
| 06 Tiendas y precios | 24 | 23 | 1 | Calendario por día de visita, conteo de fresco a viejo y avance de la importación (H-07) |
| 07 Recetas | 10 | 8 | 2 (eliminar receta, restricciones) | «Nueva receta» y «Duplicar receta» (H-07) |
| 08 Reportes | 14 | 9 | 5 (hora y día de cada reporte, con etiquetas genéricas «Hora», «Día de la semana» y «Día del mes») | 0 |
| 09 Notificaciones | 24 | 22 | 2 | Vista previa de umbral por categoría y `Ctrl+/` (H-08) |
| 10 Apariencia | 2 | 2 | — | 0 |
| 11 Datos y respaldos | 20 | 18 | 2 («Próximo» en el encabezado, `allowed_uses`) | 0 |
| 12 Personas y reparto | 12 | 12 | — | 0 |
| 13 Salud y metas | 24 | 24 | — | `Ctrl+/` (H-08) |
| 14 Inventario | 9 | 9 | — | El arrastre guarda sin vista previa (H-02); `Ctrl+/` (H-08) |

## Hallazgos

| Id | Sev. | Criterio | Ruta | Dueño |
|---|---|---|---|---|
| H-01 | alta | U4 §0.2 y P-26.05 vista previa 6 (Q2-09 A) | `frontend/src/pantallas/ajustes-a/S05Creditos.tsx:139` | front-ajustes-a |
| H-02 | alta | U4 §0.2 punto 2 (P-26.14, Q2-09 A) | `frontend/src/pantallas/ajustes-b/P2614Inventario.tsx:85` | front-ajustes-b (y frontend-base por VistaTablero) |
| H-03 | alta | U5-06: objetivos táctiles y tokens de color | `frontend/src/componentes/Grafica.tsx`, migas de `src/shell/`, `src/vistas-bd/` | frontend-base |
| H-04 | alta | El dato de prueba no valida contra el esquema de C1 | `frontend/pruebas/ajustes-b/datos.ts:21` | front-ajustes-b |
| H-05 | media | U4 P-26.05 punto 7 | `S05Creditos.tsx:140` | front-ajustes-a |
| H-06 | media | U4 P-26.05 bloque 3 (alta de promoción) | `S05Creditos.tsx:246` y `textos.ts:94` | front-ajustes-a |
| H-07 | media | U4 §0.2 punto 2 y §0.7, pendientes de 02, 04, 05, 06 y 07 | `frontend/src/pantallas/ajustes-a/` | front-ajustes-a (y frontend-base) |
| H-08 | media | U4 §0.2 punto 2 y §0.7, pendientes de 09, 13 y 14 | `frontend/src/pantallas/ajustes-b/` | front-ajustes-b (y frontend-base) |
| H-09 | media | U5-06: cifras tabulares en todo el dinero | `frontend/src/pantallas/creditos-salud-reportes/*`, `src/vistas-bd/VistaGaleria.tsx` | front-creditos-salud-reportes y frontend-base |
| H-10 | media | U5-06: FV-01, FV-04 y FV-05 | `src/componentes/Boton.tsx`, `Grafica.tsx`, `Campo.tsx` y `Chip.tsx` | frontend-base |
| H-11 | media | Datos de prueba sin tipar contra C1 | `frontend/pruebas/{ajustes-a,ajustes-b,creditos-salud-reportes}/` | los 3 agentes de área |
| H-12 | media | U1 N-230 y N-231 | `frontend/src/pantallas/creditos-salud-reportes/P21Salud.tsx:92` | front-creditos-salud-reportes |
| H-13 | baja | U1 N-226: el texto «Comparador no disponible» | `P20Comparador.tsx:89` | front-creditos-salud-reportes |
| H-14 | baja | Protocolo §6 y CONVENCIONES §5 (sin flotantes) | `creditos-salud-reportes/comun.tsx:151` | front-creditos-salud-reportes |
| H-15 | baja | C2: la clase `font-normal` no existe y no genera CSS | `ajustes-b/P2608Reportes.tsx:156` | front-ajustes-b |
| H-16 | baja | CONVENCIONES §1 y uso del script como puerta | `frontend/package.json`, `.gitignore` y `scripts/contraste-variantes.mjs` | front-variantes |
| H-17 | baja | Protocolo §4: más de 150 KB escritos | Reportes de ajustes-a (~190 KB) y ajustes-b (~175 KB) | front-ajustes-a y front-ajustes-b |

### Detalle y corrección sugerida

- **H-01 (alta).** U4 P-26.05 punto 6 no es PROVISIONAL. La vista previa de la frecuencia de pagos llama dos veces a `previewCreditPurchase`: una con la frecuencia guardada y otra con la tecleada. Con la semilla da 30-oct, 15-nov, 30-nov y 15-dic-2026, cada uno de $973.75, total $3,895.00.
  - El código devuelve `{ pendiente: 'Vista previa pendiente' }` sin citar ningún bloqueo. Con eso, `confirmar()` deja guardar sin vista previa.
  - **Corrección:** implementar la vista previa 6 con la compra de ejemplo (capital y plazo de la última compra, editables; «Escribe un monto y un plazo de ejemplo» si la línea no tiene compras). Si de verdad falta algo, abrir un bloqueo y no guardar mientras tanto.
- **H-02 (alta).** En el tablero por modo de inventario, `moverA` llama a `actualizar(...tracking_mode)` al soltar la tarjeta. U4 §0.2 punto 2 pide, para P-26.14: vista previa del estado antes y después, botón «Aplicar» o `Enter`, y `Esc` para devolver la tarjeta.
  - **Corrección:** quitar `moverA`, para que el tablero sea solo lectura como el de 04, hasta que `VistaTablero` (frontend-base) admita una vista previa antes de aplicar.
- **H-03 (alta).** Lo medí en 412 táctil sobre las pantallas revisadas:
  - El `summary` «Ver como tabla» mide 346×20 en /creditos, /creditos/1 y /salud en Tendencias. Es FV-02 y queda por debajo de 44 px y del mínimo de 24 px de WCAG 2.2.
  - Las migas «Todas las líneas» (122×20) y «Crédito» (65×20).
  - Celdas que se tocan de 36 a 40 px: «84.6 kg» (60×36), «Despensa Walmart» (71×40) y «Mensual PDF» (61×40).
  - Un evento de calendario de 40×16 en /ajustes/05.
  - Un botón sin texto de 42×36 en /ajustes/12, 13 y 14.
  - La descripción emergente de ECharts usa #FFFFFF, #6D6E73 y #B7B9BE también en oscuro (FV-03).
  - **Corrección (frontend-base):** `coarse:min-h-11` en el `summary`, las migas, las celdas y los eventos, y la descripción emergente con `--ui-surface-raised`, `--ui-text-default` y el borde del token.
- **H-04 (alta).** La respuesta falsa de `GET /api/v1/dashboard/next-payment` no trae la propiedad obligatoria `availability`. ajv contra openapi.yaml dice «must have required property 'availability'».
  - **Corrección:** completar el objeto según el esquema y tiparlo con `satisfies Respuesta<'/api/v1/dashboard/next-payment','get'>`.
- **H-05 (media).** Mientras SC-U4-1 siga abierta, U4 pide mostrar el calendario de la compra de ejemplo con el valor guardado (`previewCreditPurchase`) y rotular como pendiente solo el «Después». Hoy solo dice «PROVISIONAL (SC-U4-1)».
- **H-06 (media).** Hay un texto «Nueva promoción» que nadie usa y ninguna llamada a `POST /api/v1/promotions`. Desde Ajustes no se pueden crear promociones, y U4 bloque 3 lo pide. Faltan también `←` y `→`.
- **H-07 (media).** Pendientes que declaró ajustes-a y que siguen sin hacerse:
  - `Ctrl+/` (edición múltiple con vista previa sumada) en 02, 04 y 05.
  - El arrastre con vista previa y «Aplicar» en el tablero de 04.
  - En 06: el calendario por día de visita, el conteo de precios que pasan de fresco a viejo y el avance de la importación con `listJobs`.
  - En 07: «Nueva receta» y «Duplicar receta», que necesitan la ruta de alta de P-12 en `rutas.ts` (frontend-base).
- **H-08 (media).** Pendientes de ajustes-b: `Ctrl+/` en 09, 13 y 14, y la vista previa de umbral «por categoría al tocar una categoría de ejemplo» en P-26.09.
- **H-09 (media).** Montos sin `tabular-nums`, medidos en pantalla:
  - P-17: «Hogar $598.82 · Karen $374.93», «Hogar $1,197.64 · Karen $749.86», la galería del detalle ($3,895.00, $8,105.00, $1,200.00 y $28,800.00), «Kueski Pay · $973.75» y «compra mínima $1,500.00».
  - P-18: «Mínimo A $8,094.50 · mínimo B $5,000.00» y «Hogar $2,395.28 · Karen $1,499.72».
  - P-19: las partes de cada pago.
  - P-23: el texto del aviso.
  - P-26.02: la nota de proyección.
  - **Corrección:** envolver el dinero en `tabular-nums` (`comun.tsx` ya lo hace en 5 lugares) y agregar `tabular-nums` a `VistaGaleria` (frontend-base).
- **H-10 (media).** Los tres hallazgos se reproducen en la corrida, aunque el script sale con 0:
  - FV-01: el anillo de foco del Botón aparece con transición.
  - FV-04: ECharts recolorea las series con colores fuera de C2.
  - FV-05: Campo y Chip deshabilitados no traen el estilo de §7.0.
  - La corrección exacta está en `frontend/.storybook/hallazgos.json`.
- **H-11 (media).** Los datos falsos de las 3 carpetas de pruebas son JSON sin tipo (`Record<string, unknown>`), así que tsc no detecta cuando se apartan de C1; así se coló H-04.
  - **Corrección:** `satisfies Respuesta<ruta,'get'>` en cada respuesta.
- **H-12 (media).** Dos problemas en P-21:
  - Calcula en el navegador el porcentaje de kcal planeadas contra la meta (`planned_kcal*100/daily_kcal_target`). Esa cifra no está en N-231, que pide kcal planeadas, meta y diferencia, y no la devuelve el motor S4. Mostrar `difference_kcal` o pedir la cifra a C1.
  - La ficha de la persona no muestra «Restricciones alimentarias» (N-230: vacías hoy, decisión D).
- **H-13 (baja).** Con 503 por tiempo, `ErrorCarga` muestra el mensaje del servidor; U1 pide «Comparador no disponible», sin resultados parciales.
- **H-14 (baja).** El monto por omisión se arma con `((amount_cents - paid_cents) / 100).toFixed(2)`. Usar un formateo de centavos enteros, sin flotantes.
- **H-15 (baja).** `font-normal` no está en las escalas de C2: usar la clase de peso de C2.
- **H-16 (baja).** Dos cosas de front-variantes:
  - Editó `frontend/package.json` y `.gitignore`, que son de frontend-base. Está documentado y el integrador debe conservarlo.
  - `verificar:variantes` sale 0 aunque FV-01 a FV-05 sigan abiertos. Si entra a `verify-wf3-*`, esos 5 deben cerrarse antes o contar como falla.
- **H-17 (baja).** Para la siguiente ronda: partir el trabajo para no pasar el tope de 150 KB escritos.

## Estado de hallazgos anteriores

Esta es la ronda 1 y no hay hallazgos de revisiones anteriores. Los FV-01 a FV-05 de front-variantes siguen ABIERTOS: se reproducen 5 de 5 y quedan incluidos en H-03 y H-10.
