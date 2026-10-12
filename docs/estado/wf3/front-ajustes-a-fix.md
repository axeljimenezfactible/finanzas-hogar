# front-ajustes-a-fix · WF-3a bis (12-oct-2026, Opus)

Alcance: `frontend/src/pantallas/ajustes-a/` y sus pruebas (`frontend/pruebas/ajustes-a/`). Corrige lo que quedaba de ajustes-a en `REVISION-frontend-areas-r3.md` y en `INTEGRACION-wf3a.md` §7 (H-07, la parte de H-03 en 05, R2-05 y R2-07) y la N-284 de `docs/estado/bloqueos/C3-matriz.md` (decisión del orquestador, punto 2). Ninguna cifra ni motor cambia; 0 operaciones nuevas fuera de C1 v1.4.

## HECHO

1. **H-07 · `Ctrl+/` con vista previa sumada (U4 §0.2 punto 2 y §0.7)**, sobre `def.editarVarias` de la base, en las 4 secciones que pide la revisión:
   - **02 Presupuesto, gastos fijos**: categoría o activo. Vista previa 3 (h) sumada sobre las partidas elegidas en el mes en curso, calculada en el navegador en centavos: (c) «Gastos fijos de {mes}: antes → después», (d) las partes por persona con todas las fechas de las n partidas en ese mes (ROUND_HALF_UP, el resto al titular, K14), la parte fija que pasa del monto (RN-03), «{k} de {n} partidas toman el reparto de su categoría…» al cambiar la categoría y (f) el aviso de P-07. Un `updatePlannedItem` por partida (`previaFijosGrupo` en `S02Presupuesto.tsx`).
   - **04 Categorías**: tipo, reparto o activa. Con reparto: «{n} categorías pasan a {reparto}», el ejemplo de $199.00 por cada reparto actual, los gastos fijos activos sin reparto propio que cambian sus partes (R5-2) y «Aplica a los gastos que registres desde ahora…». Con activa = No: el aviso de archivar con los gastos sumados y, si toca Despensa o Comida fuera, el de los cálculos. Un `updateCategory` por categoría.
   - **05 Promociones** (vista tabla): activa, línea, plazo, compra mínima y fin de vigencia (PROVISIONAL, bloqueo `WF3a-bis-ajustes-a` §2). Valor común o «n valores distintos» → valor nuevo, «Promociones vigentes hoy: a de n → b de n» y el aviso de R5-3. V-08, V-21 y V-03 se validan antes de aplicar. Un `updatePromotion` por promoción.
   - **07 Recetario**: tiempos aptos con «Agregar {tiempo}» o «Quitar {tiempo}» (PROVISIONAL, §3). «Recetas aptas para {tiempo}: a de n → b de n». Si una receta se quedaría sin tiempo: «Elige al menos un tiempo de comida.» y no se aplica. Un `updateRecipe` por receta.
   - Lo común está en `grupo.ts`:
     - `guardarVarias` hace un PATCH por fila. La que falla conserva su valor, se recarga la lista y se muestra el primer error con el bloque abierto.
     - `nota` pinta los avisos sin cifra como «— → —» (PROVISIONAL, §1).
     - `resumen` y `siNo` sirven a las cuatro secciones.
   - `Enter` aplica, `Esc` cancela sin escribir y el foco regresa. Eso lo pone la base.
2. **H-03, parte 1 (evento del calendario de 05 en 412 táctil)**: ya lo cubre la base (agenda por día debajo de 768 px). Medido en la prueba U5-06 de esta área: **344 × 44**. No hizo falta cambiar nada en ajustes-a.
3. **N-284 (P-26.05)**: `LecturaAB` lee la línea con `getCredit` (`GET /api/v1/credits/{credit_id}`). El valor del interruptor A/B, el aviso «Pendiente confirmar con el acreedor» y la casilla «Ya confirmé» salen de ahí. Después de confirmar o de cambiar la lectura se vuelve a leer. `verificar:c3` ya no lista P-26.05 como «no implementado».
4. **R2-05 (vista previa 10, tasa del préstamo)**:
   - Toma la compra vigente más reciente del préstamo con `credit-purchases?credit_id=…&status=active&sort=-purchase_date,-id&limit=1`, con **su** `payment_frequency` (antes usaba la de la línea).
   - Agrega la nota «Compra de ejemplo: {capital} a {n} pagos, del {primer pago}».
5. **R2-07 (06, punto de equilibrio de la membresía)**: la fila y la vista previa del costo ahora dicen «Vista previa pendiente · PROVISIONAL (bloqueo WF3-ajustes-a §2)», que es el bloqueo que ya lo documentaba.
6. **Pruebas Playwright (una por cambio, con API falsa validada contra C1)** en `ajustes.escritorio.spec.ts`:
   - Cuatro de `Ctrl+/` (02, 04, 05 y 07). Cada una cubre: casillas, `Control+/`, propiedad y valor, el bloque `vista-previa` con cifras, `Escape` sin ningún PATCH y `Enter` con un PATCH por fila y su cuerpo exacto. La de 07 cubre además la validación que no aplica.
   - Una de N-284: la lista dice B y sin confirmar, `getCredit` dice A y confirmada, y se muestra lo de `getCredit`.
   - Una de R2-05 (consulta con `status`, `sort` y `limit`, y nota con la primera fecha) y una de R2-07.
   - En los bloques abiertos también se miden `textosEnIngles`, `dineroSinTabulares` y, en 412 táctil, `controlesChicos`.
   - La API falsa suma `PATCH planned-items/{id}`, `promotions/{id}` y `recipes/{id}` con respuesta de C1, y `GET credits/1`.

## Archivos

- `frontend/src/pantallas/ajustes-a/grupo.ts` (nuevo)
- `frontend/src/pantallas/ajustes-a/S02Presupuesto.tsx`, `S04Categorias.tsx`, `S05Creditos.tsx`, `S06Tiendas.tsx`, `S07Recetas.tsx`, `textos.ts` (claves `grupo.*` y `tiendas.bloqueoEquilibrio`)
- `frontend/pruebas/ajustes-a/ajustes.escritorio.spec.ts`
- `docs/estado/bloqueos/WF3a-bis-ajustes-a.md` (nuevo)
- `docs/estado/wf3/front-ajustes-a-fix.md` (este)

## Verificación

- `npx tsc --noEmit -p frontend`: 0 errores. En una corrida intermedia salió un error en `pruebas/ajustes-b/datos.ts`; ese archivo es de otro agente, que lo estaba editando al mismo tiempo, y en la corrida final ya no salió.
- `npm run build`: OK (PWA, 82 entradas en precaché).
- `npm run verificar:rutas`: «C3-03 frontend: 208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0».
- `npm run verificar:hosts`: solo los literales permitidos de siempre (licencia de Inter y workbox).
- `python3 -I scripts/c3-integracion.py`: «implementado sí 41/42». El único «no» que queda es P-03 (N-053, N-102), que no es de esta área.
- Playwright `pruebas/ajustes-a` (412 táctil y 1440 escritorio), con una configuración temporal en el puerto 4191 porque otros agentes usaban el 4173 a la vez (ya borrada): **59 passed, 1 skipped** (la medida táctil en escritorio), 0 failed.
- `grep` de colores hex en `src/pantallas/ajustes-a` y `pruebas/ajustes-a`: 0. Inglés visible: `textosEnIngles` = [] en los cuatro bloques de `Ctrl+/`.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: rc = 0 («Verificación ampliada OK»).
- Por criterio (los revisores evalúan contra U5, C3 y U4):
  - **U5-06**: 44 px en 412 táctil (`controlesChicos` = [] con el bloque abierto; evento de 05 a 344 × 44); colores solo por token.
  - **U5-07**: textos en es-MX, 0 inglés.
  - **U5-09**: pruebas con API falsa validada contra C1 (`erroresTrafico` = [] en todas).
  - **C3-03**: 0 rutas fuera del contrato.
  - **C3-02**: P-26.05 queda «implementado = sí» con N-284.
  - **U4 §0.2 y §0.7**: vista previa sumada antes de aplicar; `Enter` aplica a todas y `Esc` no cambia ninguna (comprobado por pantalla).

## Supuestos (PROVISIONALES, con bloqueo `docs/estado/bloqueos/WF3a-bis-ajustes-a.md`)

- §1: los avisos de U4 en el bloque de `Ctrl+/` van como renglón «— → —» porque la base solo pinta `Cambio[]`.
- §2: los campos que se editan a la vez en promociones.
- §3: «Agregar» y «Quitar» un tiempo apto en el recetario.
- 04: en «Ejemplo con un gasto de $199.00» las partes salen en el orden de `partesTexto` («Karen $99.50 · Hogar $99.50»). Es la función que ya usa la edición de una sola categoría; no la cambié.

## Bloqueos

- Nuevo: `WF3a-bis-ajustes-a` (3 puntos).
- Citados: `WF3-ajustes-a` §2 (punto de equilibrio).

## Pendientes

- Lo que la base pide a los demás no es de esta área (ajustes-b: H-08 y R2-06; P-03: N-053 y N-102).
- «Probado» de C3 cuenta en WF-3b.
- No corrí la suite completa de Playwright de todas las áreas en mi puerto. Solo cambié `ajustes-a`, y `make verify-all` pasó.

## Propuestas para Axel

- Decidir los 3 puntos de `WF3a-bis-ajustes-a`. La recomendación es la opción (a) en los tres, más el (b) del punto 1 a `frontend-base` para mostrar las notas como notas.

## Paquetes instalados

Ninguno.

## Consumo aproximado

~55 acciones · ~200 KB leídos · ~50 KB escritos.
