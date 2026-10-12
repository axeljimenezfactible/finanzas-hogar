# correccion-frontend-areas-r1 · WF-3a (12-oct-2026)

Corrige los hallazgos de `docs/estado/wf3/REVISION-frontend-areas-r1.md`. Solo se tocaron archivos de front-ajustes-a, front-ajustes-b y front-creditos-salud-reportes. Lo que es de frontend-base quedó como solicitud. No hice git add ni git commit.

## Hecho

| Hallazgo | Estado | Qué cambió |
|---|---|---|
| H-01 (alta) | HECHO | P-26.05, vista previa 6: llama dos veces a `previewCreditPurchase`, con la frecuencia guardada y con la tecleada. La compra de ejemplo es la última de la línea, con capital y plazo editables en el mismo bloque y la fecha de esa compra. Muestra fechas, pago y total «Antes → Después». Si la línea no tiene compras dice «Escribe un monto y un plazo de ejemplo» y no guarda (`Previa.bloquea`). Para líneas no revolventes, los intereses quedan PROVISIONAL (bloqueo nuevo `WF3-ajustes-compra-ejemplo`). |
| H-02 (alta) | HECHO | P-26.14: quité `moverA`, así que el tablero queda de solo lectura hasta que `VistaTablero` admita vista previa, «Aplicar» y Esc. |
| H-04 (alta) | HECHO | `GET /dashboard/next-payment` falso ya trae `availability`. |
| H-05 (media) | HECHO | Vista previa 7: «Primer pago de una compra hecha hoy» y el calendario de la compra de ejemplo, ambos con el valor guardado. El «Después» dice «Vista previa pendiente» · PROVISIONAL (SC-U4-1). |
| H-06 (media) | HECHO | Formulario «Nueva promoción» con `POST /api/v1/promotions` e Idempotency-Key. Tiene comercio, tienda, línea, plazo, tasa, compra mínima, vigencia y fuente; nace activa. Se abre con el botón del bloque, con el vacío o con Alt+N cuando el foco está en el calendario, y se guarda con Guardar o Ctrl+Enter. En el calendario, ←/→ cambian de mes. |
| H-07 (media) | PARCIAL | Hecho: «Duplicar receta» en 07 (`duplicateRecipe` y abre P-12 de la copia). Pendiente: ver Pendientes. |
| H-08 (media) | PARCIAL | Hecho: en P-26.09, el umbral por categoría con una «Categoría de ejemplo» (`getBudgetByCategory`), en la vista previa y fijo. Pendiente: `Ctrl+/`, que depende de VistaBD. |
| H-09 (media) | HECHO en las áreas | Puse `tabular-nums` en la raíz de P-17 a P-25, que se hereda a tarjetas, notas, franjas, galería y calendario. También en la nota de proyección de P-26.02. |
| H-11 (media) | HECHO | Las respuestas falsas se validan con tsc contra `esquema.gen.ts`: `satisfies` con un mapa por ruta en ajustes-a y ajustes-b, y `json(r, cuerpo: CuerpoDe<r>)` en creditos-salud-reportes. Así aparecieron y se corrigieron 41 diferencias de tipo. |
| H-12 (media) | HECHO | P-21 muestra la «Diferencia» del servidor (`difference_kcal`) y ya no el porcentaje calculado en el navegador. La ficha agrega «Restricciones alimentarias» (`dietary_restrictions`, «Sin dato» si está vacía). |
| H-03 y H-10 | FUERA DE ALCANCE | Son de frontend-base. Las pedí en solicitudes. |

## Archivos
- `frontend/src/pantallas/ajustes-a/`: `comunes.tsx` (`Previa.bloquea`, `Editable.extra` y `dependencias`), `S05Creditos.tsx`, `S07Recetas.tsx`, `S02Presupuesto.tsx` y `textos.ts`.
- `frontend/src/pantallas/ajustes-b/`: `P2614Inventario.tsx`, `P2609Notificaciones.tsx` y `textos.ts`.
- `frontend/src/pantallas/creditos-salud-reportes/`: `P17Lineas.tsx`, `P18Credito.tsx`, `P19Compra.tsx`, `P20Comparador.tsx`, `P21Salud.tsx`, `P23Centro.tsx`, `P24Mensual.tsx`, `P25Karen.tsx` y `textos.ts`.
- `frontend/pruebas/ajustes-a/ajustes.escritorio.spec.ts` (5 pruebas nuevas), `frontend/pruebas/ajustes-b/datos.ts` y `frontend/pruebas/creditos-salud-reportes/datos.ts`.
- `docs/estado/bloqueos/WF3-ajustes-compra-ejemplo.md` (nuevo) y este reporte.

## Verificación
- `npx tsc --noEmit`: 0 errores, ahora con los datos de prueba tipados.
- `npx vite build --outDir datos/tmp/corr-fa-r1/dist`: «✓ built».
- Playwright con preview propio en el puerto 4796 y Chromium 1217: `ajustes-a ajustes-b creditos-salud-reportes` da **70 passed**. Antes eran 60; las 5 nuevas corren en móvil y en escritorio.
  - H-01: «P-26.05 vista previa 6». Muestra 30-oct-2026 · 15-nov-2026 · 30-nov-2026 · 15-dic-2026, «$973.75 cada uno (4 pagos)» y $3,895.00, y guarda `{ payment_frequency: 'monthly' }` con Enter.
  - H-01: «sin compras pide la compra de ejemplo y no guarda». Con Enter hay 0 PATCH; al teclear el ejemplo aparece $3,895.00.
  - H-05: «vista previa 7». Muestra el primer pago 30-oct-2026, «Vista previa pendiente» y PROVISIONAL (SC-U4-1).
  - H-06: «alta de promoción con Ctrl+Enter». Manda POST /api/v1/promotions con `credit_id` 1, plazo 4, `rate_num` 0, vigencia 30-oct-2026 e `is_active` true.
  - H-07: «duplicar receta abre P-12 de la copia». Manda POST /recipes/1/duplicate y lleva a /comidas/recetas/2.
- `npm run verificar:rutas` (C3-03): 202 rutas invocadas, 0 fuera de openapi.yaml y 0 sin registrar.
- U5-06 y U5-07 no se volvieron a medir en pantalla: no hice el escaneo de 132 páginas. Lo que pide H-03 y H-10 sigue abierto en frontend-base.
- `make verify-all`: no lo corrí por el tope de acciones. El revisor lo reportó fallando por gofmt en `backend/internal/reportes/`, ajeno a esto.

## Supuestos
- La «última compra» es la de `purchase_date` más reciente (y, si empatan, la de id mayor), porque `listCreditPurchases` no tiene `sort`.
- En las líneas revolventes la compra de ejemplo usa `term_table` (RN-08, compra planeada, como la vista previa 4 de P-26.02).
- H-09 se resolvió por herencia CSS de `tabular-nums` en la raíz de cada pantalla del área.

## Bloqueos
- `docs/estado/bloqueos/WF3-ajustes-compra-ejemplo.md` (nuevo): intereses de la compra de ejemplo en líneas no revolventes, PROVISIONAL.
- Siguen citados: SC-U4-1 (vista previa 7) y `S3-formula-tarjeta-pdf`.

## Pendientes
- H-07, P-26.06: el calendario por día de visita, el conteo de precios que pasan de fresco a viejo y el avance de la importación con `listJobs`. Están en el alcance de ajustes-a y no se hicieron por el tope.
- H-07, P-26.07: «Nueva receta». Necesita la ruta de alta de P-12 en `rutas.ts` (frontend-base).
- H-07 y H-08: `Ctrl+/` en 02, 04, 05, 09, 13 y 14, y el arrastre con «Aplicar» en el tablero de 04. Necesitan soporte en VistaBD y VistaTablero (frontend-base).
- H-08: falta una prueba Playwright del umbral por categoría.
- Bajos que no se pidieron: H-13 a H-17.

## Solicitudes fuera de alcance (frontend-base)
1. **H-03.** Agregar `coarse:min-h-11` en:
   - el `summary` «Ver como tabla» de `src/componentes/Grafica.tsx`;
   - las migas de `src/shell/`, con `inline-flex items-center`;
   - las celdas editables de `src/vistas-bd/VistaTabla.tsx` y `CeldaEditable.tsx`;
   - los eventos de `src/vistas-bd/VistaCalendario.tsx`;
   - el botón de icono de 42×36 de /ajustes/12 a 14.

   La descripción emergente de ECharts debe usar `var(--ui-surface-raised)`, `var(--ui-text-default)` y el borde del token, en claro y en oscuro.
2. **H-10.** Aplicar FV-01, FV-04 y FV-05 tal como están en `frontend/.storybook/hallazgos.json` en `Boton.tsx`, `Grafica.tsx`, `Campo.tsx` y `Chip.tsx`.
3. **H-09.** Agregar `tabular-nums` a los valores de `src/vistas-bd/VistaGaleria.tsx`, para pantallas que no lo heredan.
4. **H-07 y H-08.**
   - `VistaBD`: selección de filas con `Ctrl+/` y un `editarVarias.vistaPrevia` sumado; Enter aplica y Esc cancela.
   - `VistaTablero`: `columnasTablero.vistaPrevia` más «Aplicar»/Enter y Esc para devolver la tarjeta.
   - `VistaCalendario`: ←/→ nativo. Hoy P-26.05 lo resuelve con los botones del mes.
   - `src/app/rutas.ts`: ruta de alta de P-12.

## Propuestas para Axel
- Elegir la opción del bloqueo `WF3-ajustes-compra-ejemplo` (A, B o C).

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 60 acciones, unos 240 KB leídos y unos 45 KB escritos (código, pruebas, bloqueo y reporte). Los temporales de `datos/tmp/corr-fa-r1/` se borraron.

## Continuación 1 (12-oct-2026)

Se retomó el trabajo desde los Pendientes de arriba, sin rehacer lo que ya estaba. Solo se tocaron archivos de front-ajustes-a y front-ajustes-b, más un bloqueo nuevo. No hice git add ni git commit.

### Hecho
| Hallazgo | Estado | Qué cambió |
|---|---|---|
| H-07, P-26.06: calendario por día de visita | HECHO | `CalendarioVisitas` en el bloque de tiendas. Pinta una semana de lunes a domingo (K8) con los mercados y tianguis activos en su `visit_weekday`, de forma recurrente. Tiene botones «Semana anterior» y «Semana siguiente», y ←/→ cuando el foco está dentro (fuera de los campos). Al tocar un mercado se abre su detalle. En móvil es una lista de una columna y en escritorio, 7 columnas. Los botones miden 44 px en táctil. Lo resolví dentro del área porque `VistaCalendario` de la base solo pinta fechas fijas por mes (ver solicitud 1). |
| H-07, P-26.06: conteo de fresco a viejo | HECHO, PROVISIONAL | En la vista previa de «Días para que un precio sea viejo», una nota dice «Con este valor, n productos pasan hoy de precio fresco a viejo · PROVISIONAL (bloqueo WF3-ajustes-precios-conteo)». C1 no da ese conteo, así que abrí el bloqueo nuevo. Mientras tanto se usa su opción A: los productos de `listProducts` cuyo `last_observed_date` tiene hoy entre el valor nuevo y el vigente menos 1 días. También quité el `pendiente: 'Registros'` que salía sin sentido en las dos vistas previas de frescura. |
| H-07, P-26.06: avance de la importación | HECHO en parte, PROVISIONAL | Mientras se importa, el botón queda inactivo con su motivo y un estado `role="status"` dice «Importando el corte; puede tardar unos minutos». El avance como tarea de `listJobs` queda PROVISIONAL: no hay job de importación PROFECO y `importProfeco` es síncrono (bloqueo existente `S6-compras-b-fuentes`, punto 2). No inventé un código de job. |
| H-07, P-26.07: «Nueva receta» | HECHO | Se abre con el botón del bloque, con el estado vacío («Sin recetas · Nueva receta») o con Alt+N. La hoja pide el nombre (V-01, V-09, V-12 si da 409) y al menos un tiempo apto («Elige al menos un tiempo de comida.»). Se guarda con Guardar o Ctrl+Enter. Hace `POST /api/v1/recipes` con Idempotency-Key y abre P-12 de la receta nueva (`rutaDe('P-12', { id })`) para capturar ingredientes y pasos. No hizo falta una ruta de alta nueva en `rutas.ts`. |
| H-08: prueba del umbral por categoría | HECHO | Prueba nueva «P-26.09 umbral por categoría». La respuesta falsa de `GET /api/v1/budget/categories` está tipada contra C1. Comprueba que el selector solo ofrece categorías con monto y que, con Despensa ($8,000.00), la nota dice $6,400.00 y $8,000.00. |

### Archivos (continuación 1)
- `frontend/src/pantallas/ajustes-a/S06Tiendas.tsx`, `S07Recetas.tsx` y `textos.ts`.
- `frontend/pruebas/ajustes-a/ajustes.escritorio.spec.ts`: 3 pruebas nuevas, el catálogo `/api/v1/products` con fechas relativas a hoy y la respuesta de `POST /api/v1/recipes` con `satisfies`.
- `frontend/pruebas/ajustes-b/ajustes-b.escritorio.spec.ts` (1 prueba nueva) y `frontend/pruebas/ajustes-b/datos.ts` (`GET /api/v1/budget/categories`).
- `docs/estado/bloqueos/WF3-ajustes-precios-conteo.md` (nuevo) y esta sección.

### Verificación (continuación 1)
- `npx tsc --noEmit`: 0 errores, con los datos de prueba nuevos tipados.
- `npx vite build --outDir datos/tmp/corr-fa-r1c1/dist`: correcto (PWA, 80 entradas en precache).
- Playwright con un preview propio en el puerto 4797 (configuración temporal en `datos/tmp/`): `ajustes-a ajustes-b creditos-salud-reportes` da **78 passed**. Antes eran 70; las 4 pruebas nuevas corren en móvil y en escritorio.
  - «P-26.06 calendario por día de visita»: «Tianguis del martes» aparece en el martes y el lunes queda vacío. ArrowRight cambia la semana y ArrowLeft la regresa. Al tocar el tianguis se abre su hoja.
  - «P-26.06 vista previa de frescura»: con 20 días da «Con este valor, 1 producto pasa hoy de precio fresco a viejo · PROVISIONAL (bloqueo WF3-ajustes-precios-conteo)». Los datos tienen precios de hace 25, 10 y 40 días y un producto sin precio.
  - «P-26.07 nueva receta con Alt+N»: sin tiempo apto muestra el mensaje y hace 0 POST. Con «Comida» manda `POST /api/v1/recipes` con `{ name: 'Caldo de verduras', meal_slot_ids: [2] }` y lleva a `/comidas/recetas/3`.
  - «P-26.09 umbral por categoría (H-08)»: en el selector aparecen «Elige una categoría» y «Despensa». La nota dice «Con el presupuesto de Despensa de este mes ($8,000.00), el primer aviso llega al gastar $6,400.00 y el de tope al gastar $8,000.00».
  - Repetición (`--repeat-each 2`) de ajustes-a más `medidas.spec.ts:58`: 86 passed.
- `npm run verificar:rutas` (C3-03): 202 rutas invocadas, 0 fuera de openapi.yaml y 0 sin registrar (216 registradas).
- `node scripts/hosts-externos.mjs` sobre mi build: solo literales que no son peticiones (`scripts.sil.org` en la licencia de Inter y `bit.ly` en workbox). Igual que antes.
- `make verify-wf3-frontend` (`npm run verify`): **139 passed y 3 failed**. Fallaron «P-26.02 monta» (móvil), «P-26.06 calendario» (escritorio) y `medidas.spec.ts:58` (frontend-base). Durante la corrida, otro agente usaba el mismo puerto 4173 y `frontend/dist`: había otro `playwright test` y otro `vite preview --port 4173` vivos. Con el preview aislado, esas 3 pruebas pasan 2 de 2 veces. Hay que repetirlo sin corridas en paralelo; no lo atribuyo a estos cambios sin esa prueba.
- `make verify-all GOOSE=…/goose`: **falla en `verify-go` por gofmt** en `backend/internal/reportes/reportes_test.go` y `backend/internal/reportes/vistas.go`. Antes de eso pasan las migraciones (Up, Down y Up) y `go build` y `go vet`. Es ajeno a este alcance; ver la solicitud 3.
- U5-06: los controles nuevos llevan `min-h-11` en táctil y los textos con cifras llevan `tabular-nums`. No volví a hacer el escaneo completo de U5-06/U5-07.

### Supuestos (continuación 1)
- La semana del calendario de visitas arranca en el lunes de la semana de hoy (K8). Solo entran los mercados y tianguis activos con día de visita.
- «Nueva receta» captura nombre y tiempos aptos, lo mismo que se edita en la tabla, y deja ingredientes y pasos para P-12 («Ingredientes y pasos se editan en P-12»). No manda `base_servings` para que el servidor ponga su valor por defecto (2, decisión D).

### Bloqueos (continuación 1)
- Nuevo: `docs/estado/bloqueos/WF3-ajustes-precios-conteo.md`. Opciones A (en uso, PROVISIONAL), B y C.
- Citado: `S6-compras-b-fuentes` punto 2 (no hay job de importación PROFECO), para el avance como tarea.

### Pendientes (continuación 1)
- `Ctrl+/` con vista previa sumada en 02, 04, 05, 07, 09, 13 y 14, y el arrastre con «Aplicar» y Esc en los tableros de 04 y 14. Siguen dependiendo de `VistaBD` y `VistaTablero` (frontend-base, solicitud 4 de arriba). Comprobé que hoy `VistaBD` todavía no los admite.
- H-03 (descripción emergente de ECharts y lo que falte de 44 px) y H-10 (FV-01, FV-04 y FV-05) siguen en frontend-base. Por ejemplo, `Boton.tsx` todavía tiene `transition-colors` junto al anillo de foco.
- El avance de la importación PROFECO como tarea, cuando se resuelva `S6-compras-b-fuentes` punto 2.

### Solicitudes fuera de alcance (continuación 1)
1. frontend-base, `frontend/src/vistas-bd/VistaCalendario.tsx` y `tipos.ts`: admitir un calendario recurrente por día de la semana (por ejemplo `diaSemana?: (fila) => 1..7` y una vista de semana con ←/→). Así P-26.06 puede pasar `CalendarioVisitas` a `VistaBD` con `vistas: [..., 'calendario']`.
2. C1 o backend (bloqueo `WF3-ajustes-precios-conteo`, opción B): un `stale_after_days` opcional en `getPriceSummary` para que el servidor dé el conteo.
3. Dueño de `backend/internal/reportes/`: correr `gofmt -w backend/internal/reportes/reportes_test.go backend/internal/reportes/vistas.go` para que `make verify-all` pase `verify-go`.

### Propuestas para Axel (continuación 1)
- Elegir la opción del bloqueo `WF3-ajustes-precios-conteo` (A, B o C). Recomiendo la B: el conteo lo daría el servidor con la misma regla RN-12 y sin el tope de 200 productos.

### Paquetes instalados (continuación 1)
Ninguno.

### Consumo aproximado (continuación 1)
Unas 47 acciones, unos 120 KB leídos y unos 25 KB escritos (código, pruebas, bloqueo y esta sección). Los temporales de `datos/tmp/corr-fa-r1c1/` se borran al terminar.
