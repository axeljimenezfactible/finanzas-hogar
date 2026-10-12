# front-comidas-compras-fix · WF-3a bis (Opus) · 12-oct-2026

Ronda acotada de corrección de Comidas y Compras (decisiones del orquestador sobre `docs/estado/bloqueos/C3-matriz.md`: punto 2 y punto 3 (b)). No cambia cifras, motores ni decisiones de Axel. Sin git add ni commit. No se envió Telegram. No hay credenciales en el código, las pruebas ni este reporte.

## Hecho
- **HECHO · P-12 N-123** (U1 §3.14: «"Agregar al plan" usa N-123»; §3.12 N-123). El botón «Agregar al plan» ya no manda a P-10 `?hoja=asignar`. Ahora abre en la misma P-12 la hoja «Agregar al plan · {receta}» con estos campos:
  - semana: la semana por defecto de U1 §0.2, con selector de semana;
  - días: chips con selección múltiple; empieza en hoy si cae en la semana y, si no, en el lunes;
  - tiempo: solo los tiempos aptos de la receta, o todos si la receta no tiene;
  - persona: Ambos por defecto, Axel o Ximena (las de `personasPlan`);
  - porciones por persona: 1 por defecto, con la ayuda «Porciones por día en total: 2» (decisión D).

  Guarda con **una** llamada a `createMealPlanEntries` (`POST /api/v1/meal-plan-entries`). La `Idempotency-Key` se repite si se reintenta la misma captura y cambia si cambia cualquier dato. El aviso dice «Receta agregada al plan en n casillas» y trae la acción «Ver plan», que abre P-10 en esa semana. Los errores 409 del servidor (`already_exists`, `slot_conflict`) se muestran dentro de la hoja, en español y tal como llegan, y la hoja no se cierra. Sin conexión, el botón queda deshabilitado con «Requiere conexión…».
- **HECHO · P-15 N-173.** La sección «Nutrición» de la ficha ya lee `listNutritionFacts` (`GET /api/v1/products/{product_id}/nutrition-facts`). Antes usaba `nutrition_facts` de `getProduct`. Muestra kcal por 100 g o 100 ml, la fuente, la cita y «preferida». Pasa por `ConDatos` (carga, error y Reintentar) y se recarga al agregar un dato manual (`createNutritionFact`).
- **HECHO · P-15 N-175**, edición en el lugar:
  - (1) Catálogo: «Categoría» se edita en la celda con `updateProduct` (`product_group_id`, con «Sin categoría» = null). El nombre se queda como el enlace que abre la ficha, porque VistaTabla pierde el botón de abrir si la primera columna es editable.
  - (2) Ficha: el **nombre** se edita en el lugar con `updateProduct`.
  - (3) Ficha: las **presentaciones** pasan a una VistaBD tabla. Nombre, marca y tamaño se editan en el lugar con `updatePresentation`.
  - (4) El **tamaño** tiene vista previa «Antes → Después» antes de guardar (U4 §0.2): el tamaño, el último precio por kg, l o pieza de esa presentación (recalculado con enteros y HALF_UP: precio × 1,000,000 ÷ contenido en milésimas para g y ml; × 1,000 para pieza; C1 «Precio por 1 kg, 1 l o 1 pieza», RN §precio por unidad base) y «Precios de esta presentación que se recalculan 0 → n». Ejemplo de la prueba: «Tamaño 1 kg → 900 g · Último precio por kg (Bodega Aurrera, 10-oct-2026) $28.90 → $32.11 · Precios … 0 → 1». Sin vista previa no guarda.
  - (5) **`Ctrl+/`** (U1 §5.3, pieza `editarVarias` de la base): cambia la categoría de varios productos a la vez, con la vista previa sumada «Productos que cambian de categoría 0 → n». `Esc` cancela sin escribir; `Enter` aplica a todas (un PATCH por producto que cambia).
- **HECHO · Esc dentro de la ficha.** Editar en el lugar y pulsar `Esc` cerraba la hoja entera, porque Radix escucha antes que la celda. La ficha ya no se cierra mientras el foco está en un control de `[data-edicion-en-lugar]` (PROVISIONAL, bloqueo punto 4; el arreglo de fondo es de frontend-base).
- **HECHO · precio por unidad (corrección encontrada al revisar P-14 y P-15).** C1 define `base_unit_price_cents` como el «Precio por 1 kg, 1 l o 1 pieza (K7)», pero P-14 y P-15 escribían «$28.90 por g» o «por ml». Nueva clave `unidadPrecio` = kg, l, pieza. El renglón de precios de la ficha ya lleva `tabular-nums` (`dineroSinTabulares` lo marcaba).
- **HECHO · P-13 y P-14 contra la semilla de desarrollo** (solo lecturas con curl, CA de mkcert, sesión de desarrollo de Axel leída del archivo fuera del repositorio, sin imprimirla):
  - P-14: el selector de canastas ofrece la «compra del 10-oct-2026» de Bodega $3,094.50 (36 artículos). La comparación responde Bodega $3,094.50; Walmart $3,359.00 (+$264.50); idénticos +$100.50; 33 idénticos y 3 sustitutos; 22, 8 y 3. **El total de Walmart y el reparto idénticos/sustitutos no son los canon de U1 §3.16** ($3,233.00, +$138.50): PREGUNTA en el bloqueo, punto 1. La pantalla pinta lo que manda la API.
  - P-13: la semilla no trae lista generada para 2026-W42 ni W41, así que P-13 sale en su estado vacío (bloqueo, punto 2).
  - Las tiendas de la semilla son sucursales («Vallejo», «Acueducto de Guadalupe»). P-13 a P-15 ahora dicen «Bodega Aurrera · Vallejo» (PROVISIONAL, bloqueo punto 3).
- **Excepción técnica (punto 3 (c))**: no aplica a mi área. P-12 y P-15 ya llaman operaciones ligadas a cada necesidad. `queryView`, `listPeriods`, `listIncomes` y `deleteAttachment` no son de Comidas y Compras.

## Archivos
- `frontend/src/pantallas/comidas-compras/P12Receta.tsx`: `HojaAgregarPlan` (N-123).
- `frontend/src/pantallas/comidas-compras/P15Productos.tsx`: N-173, N-175, `Ctrl+/`, Esc en la ficha, precio por kg y `tabular-nums`.
- `frontend/src/pantallas/comidas-compras/P14Canasta.tsx` y `P13Lista.tsx`: precio por kg y `nombreTienda`.
- `frontend/src/pantallas/comidas-compras/comun.tsx`: `nombreTienda`.
- `frontend/src/pantallas/comidas-compras/textos.ts`: 13 claves nuevas, es-MX.
- `frontend/pruebas/comidas-compras/correccion-bis.escritorio.spec.ts` (nuevo, 8 pruebas).
- `docs/estado/bloqueos/WF3a-bis-comidas-compras.md` (nuevo) y este reporte.

## Verificación
- `npx tsc --noEmit`: 0 errores en `comidas-compras`. En el momento de la corrida había errores ajenos a mi área en `gastos-presupuesto/P02Gastos.tsx` y `P04Gasto.tsx`, de otro agente en curso. Antes, `npm run build` completo pasó (tsc + vite + PWA, 82 entradas). La última compilación se hizo con `npx vite build` (✓) por esos errores ajenos.
- `npx playwright test pruebas/comidas-compras/` (FH_CHROMIUM headless-shell 1217): **30 passed, 1 skipped**: el `Ctrl+/` en móvil, igual que en la base (U4 §0.7). `--repeat-each=4` de la prueba `Ctrl+/` en escritorio: 4 de 4.
- `npm run verificar:rutas`: «205 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0» (C3-03).
- `python3 -I frontend/scripts/c3-integracion.py` (sin `--escribir`): «implementado sí 39/42». P-12 y P-15 ya no están en «no implementado». «Operaciones con pantalla que el frontend no llama: 4 [deleteAttachment, listIncomes, listPeriods, queryView]»: `listNutritionFacts`, `updateProduct` y `updatePresentation` salieron de la lista.
- grep de hex (`#[0-9a-fA-F]{3,8}`) en pantallas y pruebas del área: **0**. `dark:` y colores en línea: **0**. Inglés: `textosEnIngles` = [] en P-12 (hoja) y en P-15 (ficha y edición). `dineroSinTabulares` = [] en la ficha.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: **exit 0** («Verificación ampliada OK»).
- Pruebas, una por cambio:

| Cambio | Prueba |
|---|---|
| N-123 | escribe una sola vez: días 12 y 13, tiempo cena (3), personas [1,2], 1000, con clave de idempotencia; la URL sigue en P-12. Otra prueba: un 409 se ve en la hoja |
| N-173 | GET `nutrition-facts`; «18 kcal por 100 g · USDA · FDC 170457 · preferida» |
| N-175 · catálogo | PATCH `/products/2 {product_group_id:1}` |
| N-175 · ficha | nombre por PATCH; tamaño con vista previa, Esc sin PATCH y la hoja sigue abierta, Enter con `{size_milli:900000}`; marca por PATCH |
| `Ctrl+/` | vista previa, Esc sin PATCH, Enter con 2 PATCH |
| P-14 con sucursales | «Bodega Aurrera · Vallejo es la más barata: $3,094.50.» |

## Supuestos
- El tiempo «apto» de N-123 sale de `recipe.meal_slot_ids`, igual que P-10. Si la receta no tiene, se ofrecen todos.
- «Agregar al plan» no lee el plan de la semana antes de escribir. Los choques los responde el servidor con 409, que se muestra tal cual.
- `Ctrl+/` en P-15 se apoya en U1 §5.3 (atajo global de las vistas de base de datos). U4 no lo pide para P-15. Solo se ofrece «Categoría», que no mueve dinero.

## Bloqueos
`docs/estado/bloqueos/WF3a-bis-comidas-compras.md`, 4 puntos:
1. Cifras de P-14 ≠ canon (dueño: semilla o backend S5).
2. P-13 sin lista en la semilla.
3. Formato «Cadena · Sucursal».
4. Esc en la hoja (base).

## Pendientes
- Hallazgo para **frontend-base**: tras `Esc` en el bloque de `Ctrl+/` de P-15, el foco regresó a `BODY` en 2 de 4 corridas (en las otras 2 regresó a la casilla). La prueba vuelve a enfocar la casilla antes del segundo `Ctrl+/`. En el banco de la base no pasa; parece que la tabla se vuelve a pintar.
- WF-3b: P-13 con cifras después de F-2 contra el servidor real, y las capturas U5 de P-12 a P-15.

## Propuestas para Axel
- Decidir el punto 1 del bloqueo: ¿el canon de P-14 en desarrollo es el de U1 §3.16 ($3,233.00 y +$138.50) o lo que dé el overlay?
- Confirmar el formato «Cadena · Sucursal» (punto 3).

## Paquetes instalados
Ninguno.

## Consumo aproximado
~57 acciones · ~165 KB leídos · ~35 KB escritos.
