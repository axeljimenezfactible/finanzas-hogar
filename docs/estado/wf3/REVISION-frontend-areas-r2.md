# REVISIÓN frontend de áreas · ronda 2 de 3 · WF-3a (12-oct-2026)

Revisor: `revisor-frontend-areas-r2` (independiente; no escribí lo revisado ni corregí nada). Trabajo revisado: front-creditos-salud-reportes, front-ajustes-a, front-ajustes-b y front-variantes, con la corrección `correccion-frontend-areas-r1` y su continuación 1. Criterios: U5-06, U5-07, U5-09 (docs/criterios/U5.md), C3-03 (C3.md), DESIGN.md, U1 y U4.

## Veredicto: NO APROBADO · 3 altos y 6 medios

| Severidad | Abiertos | Ids |
|---|---|---|
| alta | 3 | R2-01 (nuevo), R2-02 (nuevo), H-03 (reabierto en parte) |
| media | 6 | H-07 y H-08 (reabiertos en parte), H-09 (resto), H-10, R2-03 (nuevo), R2-04 (nuevo) |
| baja | 3 | R2-05, R2-06, R2-07 |

Todo lo de abajo lo medí yo, en `datos/tmp/rev-fa-r2/` (ya borrado): build propio, un preview aislado en el puerto 4798 y una configuración temporal de Playwright. Había otro `playwright test` con `vite preview` en el 4173 (otro agente), así que no usé `frontend/dist` ni ese puerto.

## Tabla criterio → estado → evidencia

| Criterio | Estado | Evidencia (mi corrida) |
|---|---|---|
| C3-03: build | CUMPLE | `npx tsc --noEmit` da 0 errores. `npx vite build --outDir datos/tmp/rev-fa-r2/dist` da «✓ built» y la PWA con 80 entradas en precache. No corrí `npm run build` porque escribe `frontend/dist`, que estaba usando otra corrida. |
| C3-03: rutas | CUMPLE | `npm run verificar:rutas`: «202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)». Por mi cuenta crucé con openapi las operaciones que cita U4 en las 14 secciones: todas tienen su ruta en el código de ajustes-a o ajustes-b, salvo `queryView`, que va por VistaBD, y las de otras pantallas (`createIncome` en P-08 y `suggestRecipeWithAi` en P-11 y P-12). |
| U5-08 (referencia) | CUMPLE | `node scripts/hosts-externos.mjs` sobre mi dist: solo literales que no son peticiones (react.dev, tailwindcss.com, licencia de Inter y bit.ly en workbox). |
| Pruebas Playwright de las 3 áreas | CUMPLE | Con `ajustes-a ajustes-b creditos-salud-reportes` da **78 passed**, en móvil 412 táctil y en escritorio 1440. Las copias instrumentadas también dan 78 passed en claro y 78 passed en oscuro. |
| Datos de prueba contra el esquema C1 | **NO CUMPLE** | Grabé las 775 respuestas `/api/` de las 78 pruebas y las validé con ajv 2020 contra `openapi.yaml` 3.1. Pasan 795 validaciones de respuesta y de cuerpo de petición, y todos los GET validan (next-payment ya trae `availability`). Fallan las respuestas falsas de escritura de ajustes-a y ajustes-b (R2-01). |
| U5-06: objetivos táctiles de 44 px | **NO CUMPLE** | Con `controlesChicos` al final de cada prueba en móvil medí 832 controles y 3 están bajo 44: «Semana anterior» y «Semana siguiente» miden 40×44 en /ajustes/06 (R2-02), y el evento «Walmart · Sin intereses» mide 40×44 en /ajustes/05 (H-03, VistaCalendario). Además, «Ver como tabla» sigue en 20 px con puntero fino (FV-02, 12 mediciones). |
| U5-06: cifras tabulares | NO CUMPLE (medio) | Con `dineroSinTabulares` en claro y en oscuro quedan 0 casos en P-17 a P-25. Queda 1 en /ajustes/02: el aviso del alcance «De este mes en adelante: oct-2026 y los siguientes quedan en $…» (S02Presupuesto.tsx:99). |
| U5-06: colores por token, foco y deshabilitado | NO CUMPLE (alto y medio, de frontend-base) | `contraste-variantes.mjs` con Storybook construido en tmp: 37 historias, 148 páginas y 5 230 mediciones. «FALLAS nuevas: 0», pero «hallazgos conocidos: 5 de 5 se reproducen»: FV-01 tiene 128 mediciones, FV-02 12, FV-03 384, FV-04 52 y FV-05 16. En el código, `Grafica.tsx:48` deja `tooltip` sin colores y `Boton.tsx:6` conserva `transition-colors` junto al anillo de foco. |
| U5-06: hex en el código de las áreas | CUMPLE | `grep -E "#[0-9a-fA-F]{3,8}\b\|rgb(\|hsl(\|oklch(\|dark:"` sobre las 3 carpetas de pantallas da 0. |
| U5-07: idioma, fechas y dinero | CUMPLE | Revisé el `innerText` del estado final de 156 pruebas (78 en claro y 78 en oscuro): 0 fechas ISO, 0 fechas con barras y 0 importes fuera de `$1,234.56`. En inglés solo aparece «Mar», que es la abreviatura en español de martes y marzo (falso positivo). |
| U5-09 (referencia, no es de estas áreas) | CUMPLE | «F-1 (U5-09)» da 3 toques en móvil y 3 en escritorio. Desde Inicio, también 3. |
| DESIGN §7, accesibilidad del Botón (l. 655: «nombre visible o `aria-label`») | NO CUMPLE (medio) | En /ajustes/12 hay 6 botones sin nombre, en /13 hay 4 y en /14 otros 6: son celdas editables vacías (R2-03). |
| U4: las 14 secciones (campos, validaciones y vista previa) | NO CUMPLE (medio) | Ver la tabla por sección. Faltan `Ctrl+/` y el arrastre con vista previa en los tableros 04 y 14 (H-07 y H-08). Las vistas previas 6 y 7 de P-26.05 y el alta de promoción ya están. |
| U4: claves `reports.*` y `alerts.*` | CUMPLE (con nota baja) | Las 5 claves `reports.*` y las 6 `alerts.*` de U4 son las mismas que las del extracto C1. P2608Reportes.tsx:34 las guarda con `PATCH /api/v1/settings/reports` y P2609Notificaciones.tsx:27 con `PATCH /api/v1/settings/alerts`. Ninguna prueba edita una de esas claves (R2-06). |
| U1: P-17 a P-25 | NO CUMPLE (medio) | Encontré los elementos N-190 a N-262 por su texto y su ruta, salvo N-225 en P-20 (R2-04). P-22 es de la base. P-21 ya cumple N-230 y N-231. |
| front-variantes (su alcance) | CUMPLE | Las 37 historias y la cobertura de 4 combinaciones dan 0 fallas nuevas y 0 problemas de corrida. Los 5 hallazgos FV-xx son de la base, no de este agente. |

## Tabla de las 14 secciones de Ajustes (U4)

Método: en cada sección crucé con un script tres cosas. (a) Las filas de las tablas de campos de U4 contra los textos del código, incluidas sus partes. (b) Las operaciones que cita U4 contra las rutas que llama el código. (c) Lo que mostraron las pruebas. Las filas que el script no encontró las revisé a mano, y casi todas son etiquetas compuestas o vienen de los metadatos del parámetro que da el servidor.

| Sección | Filas U4 | Operaciones de U4 / con ruta en el código | Presentes | Faltan o fallan |
|---|---|---|---|---|
| P-26.01 Cuenta y roles | 17 | 11 / 11. `logout` va por `useSesion().salir` y `queryView` por VistaBD | Sesiones con los 3 alcances (`user_except_current` y `all_except_current`), certificado y K44 | La respuesta falsa de `PATCH settings/security` no valida (R2-01) |
| P-26.02 Presupuesto | 20 | 9 / 9 | Monto con alcance «Confirmar», gastos fijos, despensa programada y proyección. V-17 sale con el mensaje 409 del servidor | `Ctrl+/` (H-07). Cifras no tabulares en el aviso del alcance (H-09). La respuesta falsa de `PATCH budgets` no valida (R2-01) |
| P-26.03 Ingresos | 7 | 4 / 3. `createIncome` se hace en P-08 y el enlace «Registrar ingreso extra» existe | Todos | — |
| P-26.04 Categorías | 8 | 6 / 6 | Reparto con vista previa en la celda | El tablero no tiene arrastre con vista previa, «Aplicar» y Esc (U4 l. 235). `Ctrl+/` (H-07) |
| P-26.05 Créditos y condiciones | 53 | 18 / 18 | Vista previa 6 con dos `previewCreditPurchase` (S05Creditos.tsx:127). Vista previa 7 con el valor guardado y el «Después» PROVISIONAL (SC-U4-1). Alta de promoción con `POST /promotions` y Ctrl+Enter. ←/→ en el calendario (l. 302) | `Ctrl+/` (H-07). El evento del calendario mide 40 px de ancho (H-03). La respuesta falsa de `PATCH credits` y `POST promotions` no valida (R2-01) |
| P-26.06 Tiendas y precios | 24 | 14 / 13 | Calendario por día de visita. Conteo de fresco a viejo PROVISIONAL (WF3-ajustes-precios-conteo). Avance de la importación PROVISIONAL (S6-compras-b-fuentes) | Las flechas de semana miden 40×44 (R2-02). El punto de equilibrio no llama a `getBasketComparison` y queda PROVISIONAL por WF3-ajustes-a §2, pero sin decir en pantalla cuál bloqueo (R2-07) |
| P-26.07 Recetas | 10 | 14 / 12. Los tiempos van por `getRecipeSettings.meal_slots`, que U4 admite, y la IA se pide en P-11 y P-12 | «Nueva receta» con Alt+N y «Duplicar receta», con sus pruebas | `Ctrl+/` (H-07) |
| P-26.08 Reportes y frecuencias | 13 | 7 / 7 | Las 5 claves `reports.*` con `PATCH settings/reports` | Ninguna prueba edita una clave (R2-06) |
| P-26.09 Notificaciones | 23 | 11 / 11 | Las 6 claves `alerts.*` con `PATCH settings/alerts`. Umbral por categoría con «Categoría de ejemplo» y su prueba | `Ctrl+/` (H-08). La respuesta falsa de `DELETE telegram-chat-candidates` es 200 y C1 pide 204 (R2-01) |
| P-26.10 Apariencia | 2 | 0 (es del dispositivo) | Tema al instante | — |
| P-26.11 Datos y respaldos | 20 | 9 / 9 | Respaldos, parámetro con vista previa y activación crítica de la IA | La respuesta falsa de `PATCH ai/providers` no valida (R2-01) |
| P-26.12 Personas y reparto | 12 | 4 / 4 | Papeles, reglas y V-16 | 6 celdas editables sin nombre (R2-03) |
| P-26.13 Salud y metas | 24 | 8 / 8 | Perfil con meta sugerida y metas | `Ctrl+/` (H-08). 4 celdas sin nombre (R2-03) |
| P-26.14 Inventario | 9 | 3 / 3 | Modos, estado con texto y días restantes. El tablero ya no guarda sin vista previa | El tablero es de solo lectura: falta el arrastre con vista previa, «Aplicar» y Esc (U4 P-26.14). `Ctrl+/` (H-08). 6 celdas sin nombre (R2-03) |

## Hallazgos de esta ronda

| Id | Sev. | Criterio | Ruta | Dueño | Qué pasa | Corrección sugerida |
|---|---|---|---|---|---|---|
| R2-01 | alta | Prueba con datos que no validan contra C1 | `frontend/pruebas/ajustes-a/ajustes.escritorio.spec.ts:61`, `frontend/pruebas/ajustes-b/datos.ts:64` | front-ajustes-a, front-ajustes-b | Las respuestas falsas de escritura no son de C1. ajustes-a contesta toda escritura con 200 `{}`, así que fallan contra ajv `PATCH settings/security` (falta `parameters`), `PATCH budgets/{id}` (Budget), `PATCH credits/{id}` (Credit) y `POST promotions`, donde C1 pide 201 con un Promotion. ajustes-b contesta con 202 o 200 y `{ id: 99, ...cuerpo }`: `PATCH ai/providers/{provider}` falla porque faltan `provider` y `allowed_uses`, y `DELETE telegram-chat-candidates/{id}` da 200 donde C1 pide 204. Si se ejercita, `PATCH settings/reports` también falla, porque faltan las claves obligatorias. tsc no lo detecta porque estas respuestas no están tipadas: H-11 solo tipó los GET. | Una respuesta por operación con `satisfies Respuesta<ruta, método>` y el código de C1 (201 en los POST de alta, 204 en DELETE). Agregar a la prueba una validación ajv del tráfico grabado. |
| R2-02 | alta | U5-06: 44 px | `frontend/src/pantallas/ajustes-a/S06Tiendas.tsx:199` y `:201` | front-ajustes-a | Con 412 táctil, «Semana anterior» y «Semana siguiente» miden 40×44. El `Boton tam="icono"` (`coarse:size-11`) se encoge dentro del `span.flex` porque el rango de la semana es largo. | Poner `shrink-0` en los dos botones, o dejar que el rango baje de línea. Pedir a frontend-base `shrink-0` en `tam: 'icono'`. |
| R2-03 | media | DESIGN §7 Botón, accesibilidad («nombre visible o `aria-label`»), y U4 §0.3 | `frontend/src/vistas-bd/VistaTabla.tsx:11`, `frontend/src/vistas-bd/CeldaEditable.tsx:39` | frontend-base (lo afectado es de ajustes-b) | Con un valor nulo y sin `pintar`, la celda pinta `''` y el `<button>` editable queda sin nombre ni texto. Son 6 en /ajustes/12, 4 en /13 y 6 en /14. | En `Celda`, poner `textos.comun.sinValor` cuando el valor está vacío y `aria-label` con «{prop.nombre}: {valor}». En ajustes-b, `pintar` con «Sin dato» en las propiedades opcionales. |
| R2-04 | media | U1 N-225 (P-20 «Valores por defecto») | `frontend/src/pantallas/creditos-salud-reportes/P20Comparador.tsx:71` | front-creditos-salud-reportes | P-20 no lee `getOptimizerSettings` (`GET /api/v1/settings/optimizer`) ni muestra los valores globales de las restricciones duras, el tope por quincena de cada persona (N-309) ni los pesos. Los campos salen vacíos con «Vacío: se usan las de Ajustes». | Precargar los campos con los valores globales, con la etiqueta «De Ajustes», y mostrar el tope por persona y los pesos normalizados. Si cambian, que sea solo para esta comparación. |
| R2-05 | baja | U4 P-26.05, vista previa 10 | `frontend/src/pantallas/ajustes-a/S05Creditos.tsx:158-161` | front-ajustes-a | La vista previa de la tasa del préstamo toma `credit-purchases?limit=1` sin orden y la frecuencia de la línea (`c.payment_frequency`). U4 pide la frecuencia y la primera fecha de la compra del préstamo. | Usar la compra del préstamo con su `payment_frequency`, y ordenar o filtrar. |
| R2-06 | baja | Cobertura de pruebas de las claves `reports.*` y `alerts.*` | `frontend/pruebas/ajustes-b/ajustes-b.escritorio.spec.ts` | front-ajustes-b | Ninguna prueba edita una clave ni comprueba el cuerpo de `PATCH settings/reports` o `PATCH settings/alerts`. | Agregar una prueba por operación, con su respuesta falsa válida en C1. |
| R2-07 | baja | Regla «PROVISIONAL citando bloqueo» | `frontend/src/pantallas/ajustes-a/S06Tiendas.tsx:128` | front-ajustes-a | El punto de equilibrio solo dice «Vista previa pendiente · PROVISIONAL», sin el bloqueo (WF3-ajustes-a §2). | Que diga «PROVISIONAL (bloqueo WF3-ajustes-a)», como en las demás vistas previas. |

## Estado de los hallazgos de la ronda 1

| Id | Sev. | Estado | Evidencia mía |
|---|---|---|---|
| H-01 (vista previa 6, frecuencia) | alta | **CERRADO** | S05Creditos.tsx:127 llama dos veces a `previewCreditPurchase`. La prueba «P-26.05 vista previa 6» pasa en las dos vistas (30-oct-2026 … $3,895.00), igual que «sin compras … no guarda». |
| H-02 (tablero 14 guarda sin vista previa) | alta | **CERRADO como alto** | Se quitó `moverA` (P2614Inventario.tsx:84, sin `columnasTablero.moverA`), así que ya no guarda sin vista previa. El arrastre con vista previa que pide U4 sigue faltando y queda dentro de H-08 (medio). |
| H-03 (44 px y descripción emergente) | alta | **REABIERTO EN PARTE** (frontend-base) | Ya cumplen: las celdas (`coarse:min-h-11 coarse:min-w-11`), «Ver como tabla» con puntero grueso (`coarse:py-3`), las migas (ocultas en 412) y el botón de 42×36. Siguen abiertos: el evento de VistaCalendario de 40×44 en /ajustes/05, «Ver como tabla» en 20 px con puntero fino (FV-02, 12 mediciones) y la descripción emergente con #FFFFFF, #6D6E73 y #B7B9BE también en oscuro (FV-03, 384 mediciones; Grafica.tsx:48). |
| H-04 (next-payment sin `availability`) | alta | **CERRADO** | datos.ts:26 ya trae `availability`. Grabado y validado con ajv: válido. |
| H-05 (vista previa 7) | media | **CERRADO** | S05Creditos.tsx:133-134 muestra el primer pago y el calendario con el valor guardado, y el «Después» como pendiente. Su prueba pasa. |
| H-06 (alta de promoción, ←/→) | media | **CERRADO** | `POST /api/v1/promotions` con Idempotency-Key (l. 398). La prueba «alta de promoción con Ctrl+Enter» pasa y ←/→ funciona (l. 302). |
| H-07 (pendientes de ajustes-a) | media | **REABIERTO EN PARTE** | Ya está: en 06, el calendario por día de visita, el conteo PROVISIONAL y el avance PROVISIONAL; en 07, «Nueva receta» y «Duplicar receta»; todo con pruebas que pasan. Falta `Ctrl+/` con vista previa sumada en 02, 04, 05 y 07, y el arrastre con vista previa y «Aplicar» en el tablero de 04 (U4 l. 235). En vistas-bd no hay manejo de `Ctrl+/`. Dueños: front-ajustes-a y frontend-base (VistaBD y VistaTablero). |
| H-08 (pendientes de ajustes-b) | media | **REABIERTO EN PARTE** | Ya está el umbral por categoría con su prueba. Faltan `Ctrl+/` en 09, 13 y 14, y el arrastre con vista previa, «Aplicar» y Esc en el tablero de 14, que hoy es de solo lectura. Dueños: front-ajustes-b y frontend-base. |
| H-09 (cifras tabulares) | media | **REABIERTO (resto)** | P-17 a P-25 tienen 0 casos en claro y en oscuro, y VistaGaleria ya tiene `tabular-nums`. Queda el aviso del alcance de /ajustes/02 (S02Presupuesto.tsx:99, `Aviso` sin `tabular-nums`). Dueño: front-ajustes-a. |
| H-10 (FV-01, FV-04 y FV-05) | media | **ABIERTO** (frontend-base) | Se reproducen los tres: 128, 52 y 16 mediciones. `Boton.tsx:6` conserva `transition-colors`. |
| H-11 (datos de prueba sin tipar) | media | **CERRADO para los GET, sustituido por R2-01 (alto)** | Todos los GET validan con ajv, pero las respuestas falsas de escritura no están tipadas y no validan. |
| H-12 (P-21: % en el navegador y restricciones) | media | **CERRADO** | P21Salud.tsx usa `difference_kcal`, ya no calcula `* 100 /` y la ficha muestra `dietary_restrictions` (l. 77). |

## Lo que no corrí
- `make verify-all` y `make verify-wf3-frontend`: usan `frontend/dist` y el puerto 4173, que en ese momento ocupaba la corrida de otro agente. En su lugar corrí lo mismo de forma aislada: tsc, build, rutas, hosts, Playwright y variantes.
