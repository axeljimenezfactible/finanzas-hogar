# frontend-min-fix · WF-3b · 12-oct-2026

Agente: `frontend-min-fix` (Opus). Base: HEAD `6d691d0` («[WF-3a bis]»), con el árbol de trabajo compartido. Otros agentes de WF-3b tenían cambios sin commit en `backend/internal/semilla/` y `backend/internal/api/comprasa/`; no son míos y no los toqué.

Encargo: los 3 cambios locales de INTEGRACION-wf3a.md §15.8, que salen de §15.4: core R5-01, áreas R5-01 y R5-05. No hice `git add` ni commit.

## Hecho

### 1. core R5-01 · P-10 Tendencias con `fechaConDia()` (U5-07, CONVENCIONES-frontend §5)
- `P10Plan.tsx`, `Tendencias`: las categorías de las dos gráficas, «Costo semanal de comidas» y «kcal promedio por día contra la meta», salen ahora de `fechaConDia(w.week_start)` y muestran «Lun 5 oct 2026». Antes se armaban a mano con `diaCorto(…) + ' ' + textoMes(…)` y daban «Lun 5 oct-2026». Esas mismas etiquetas son los encabezados de fila de «Ver como tabla», así que la tabla también queda corregida.
- Prueba nueva en `pruebas/comidas-compras/pantallas.escritorio.spec.ts`: «P-10 Tendencias (core R5-01)…», en `movil` y en `escritorio`. Revisa 4 cosas:
  - el `<text>` del SVG de cada gráfica contiene «Lun 5 oct 2026»;
  - los `rowheader` de «Ver como tabla» son exactamente `['Lun 5 oct 2026', 'Lun 12 oct 2026']`;
  - en `main` hay 0 formatos híbridos (`Lun 5 oct-2026`) y 0 fechas ISO;
  - la página no se desborda.

### 2. áreas R5-01 · U4 §0.3 «Error al guardar» y V-18 al eliminar una promoción en uso
**Piezas comunes nuevas en `comunes.tsx`:**
- `textoFallo(x, porCodigo)` decide el texto de un guardado que falla:
  - si el código de C1 trae su V-xx (por ejemplo `already_exists` → V-12), se usa ese V-xx;
  - un 4xx de validación usa el mensaje del campo, igual que `Editable` (`:144`);
  - cualquier otro error usa el texto de U4 §0.3: «No se pudo guardar. Revisa la conexión con el servidor de casa».
- `ErrorGuardado` pinta el error: `circle-alert` y `text-role-error` con `role="alert"`, y «· Reintentar» cuando el error es el de §0.3.

**`DialogoCritico`:**
- Ahora tiene `catch`. Si la confirmación falla, el diálogo queda abierto y muestra «No se pudo guardar… · Reintentar». Reintentar repite el mismo paso.
- Prop nueva `porCodigo`. Un código de C1 muestra su V-xx y una acción alterna en lugar del verbo crítico.
- El error se limpia al cerrar el diálogo.
- Esto corrige a la vez los borrados de S01 (cerrar sesiones), S02 (quitar el monto), S05 (eliminar tasa y eliminar promoción) y S07 (eliminar receta).

**Eliminar una promoción (S05):**
- Un 409 `in_use` muestra V-18: «No se puede borrar porque hay registros que la usan. Puedes desactivarla.».
- En lugar de «Eliminar la promoción de …» aparece «Desactivar», que manda `PATCH /api/v1/promotions/{promotion_id}` con `{ is_active: false }` (U4:293; C1, `deletePromotion`).
- Texto nuevo `t.comun.desactivar = 'Desactivar'` en `ajustes-a/textos.ts`. Es la única línea que toqué en ese archivo.

**Los 6 manejadores.** Cada uno tiene ahora `try/catch`. Ante un error muestran el aviso en español dentro de la hoja, que no se cierra, y lo capturado no se pierde:

| # | Manejador | Archivo | Dónde sale el aviso | Códigos con V-xx |
|---|---|---|---|---|
| 1 | «Pesos iguales» (`PATCH /settings/optimizer`) | `S05Creditos.tsx` (Comparador) | bajo el botón | — |
| 2 | «Nueva tasa» (`POST /credits/{id}/term-rates`) | `S05Creditos.tsx` (TasasPlazo) | bajo la fila de alta | `already_exists` → V-12 |
| 3 | «Nuevo crédito» (`POST /credits`) | `S05Creditos.tsx` | en el pie de la hoja, sobre los botones | `already_exists` → V-12 |
| 4 | «Nueva promoción» (`POST /promotions`) | `S05Creditos.tsx` | en el pie de la hoja | `already_exists` → V-12 |
| 5 | «Nuevo gasto fijo» (`POST /planned-items`) | `S02Presupuesto.tsx` | en el pie de la hoja | — |
| 6 | «Capturar ingreso quincenal» (`POST /planned-items`) | `S03Ingresos.tsx` | en el pie de la hoja | — |

Dos ajustes más, del mismo tipo de defecto:
- **Idempotencia.** En las altas 3, 5 y 6 la clave se renueva solo después de guardar bien; así, Reintentar reusa la misma clave. Antes, la clave `useState(nuevaClave)` no cambiaba después de un alta, y una segunda alta con la hoja montada habría repetido la clave. La prueba 2 comprueba que el reintento usa la misma clave.
- **«Restablecer al valor inicial» de `Parametro` (`comunes.tsx`).** Tenía un `await` sin `catch`, el mismo defecto aunque el revisor no lo contó. Ahora un fallo muestra un aviso de error con «Reintentar».

**4 pruebas nuevas**, al final de `pruebas/ajustes-a/ajustes.escritorio.spec.ts`. Usan `page.route` (gana sobre las rutas de contexto de `preparar()`), corren en `movil` y en `escritorio`, y cada una comprueba el mensaje y 0 `pageerror`:
1. «Pesos iguales» con 500 → «No se pudo guardar…» → Reintentar → 2 PATCH con pesos 1/1/1/1. La respuesta 200 se valida contra C1.
2. Alta de promoción con 500 → la hoja queda abierta con «Chedraui» y «4» → Reintentar con la misma `Idempotency-Key` → `POST /promotions`.
3. Eliminar promoción con 500 → el diálogo crítico queda abierto con «No se pudo guardar…» → Reintentar → `DELETE /promotions/1`.
4. Eliminar promoción con 409 `in_use` → V-18 exacto y el verbo crítico desaparece → «Desactivar» → `PATCH /promotions/1 {is_active:false}`. El 409 se valida con `erroresC1`.

### 3. R5-05 · `correccion-bis.escritorio.spec.ts:79`, ahora estable
**Diagnóstico.** Hice 3 corridas completas (`npx playwright test`, 245 pruebas) con una sonda temporal, ya quitada, que registraba antes y después del `Enter`: el `document.activeElement` y el estado del botón «Aplicar».

| Corrida | Resultado | Foco al pulsar `Enter` |
|---|---|---|
| 1 | 241 passed | `SELECT\|Propiedad`, dentro del bloque; después del Enter, «Guardando…» |
| 2 | 241 passed | igual que la corrida 1 |
| 3 | **1 failed**, en `:79` | `INPUT\|Seleccionar todas`, **fuera del bloque**; «Aplicar» habilitado y 0 PATCH |

**Causa.** Al cerrar con `Esc`, `VistaBD.cerrarGrupo` (`src/vistas-bd/VistaBD.tsx:41`) devuelve el foco a «Seleccionar todas» dentro de un `requestAnimationFrame`.
- En una corrida completa con carga, ese cuadro puede llegar tarde: después de que el bloque se volvió a abrir y tomó el foco.
- Entonces el foco salta fuera del bloque.
- El `Enter` va con `page.keyboard.press` al elemento enfocado. `EdicionVarias` escucha `Enter` solo en su `<section onKeyDown>` (`EdicionVarias.tsx:55` y `:67`), así que nadie lo oye.

No era el estado de la vista previa ni la conexión: en la corrida que falló, «Aplicar» estaba habilitado.

**Corrección**, solo con esperas estables y sin reintentos:
- después de `Esc`, `await expect(todas).toBeFocused()`. Esa aserción solo se cumple cuando ya corrió el `requestAnimationFrame` pendiente, y además comprueba el regreso de foco de U4 §0.7;
- al reabrir, `toBeFocused()` sobre el combo «Propiedad» del bloque, que es el foco inicial;
- `Enter` con `bloque.getByRole('combobox', { name: 'Valor nuevo' }).press('Enter')`, que enfoca el control del bloque justo antes de la tecla;
- `expect(bloque).toHaveCount(0)` antes de contar los 2 PATCH.

La prueba no se debilita: sigue probando que `Enter` aplica a todas, con los mismos 2 PATCH y los mismos cuerpos, y suma 2 aserciones de foco.

**Sonda de confirmación.** La corrí y la borré. Con una copia de la prueba en la que el `requestAnimationFrame` llega hasta 2 s tarde, la versión corregida pasa 3 de 3: espera ese regreso de foco, unos 3.3 s. Además, `--repeat-each=5`: 5 passed.

## Archivos

- `frontend/src/pantallas/comidas-compras/P10Plan.tsx` (1 línea)
- `frontend/src/pantallas/ajustes-a/comunes.tsx` (`textoFallo`, `ErrorGuardado`, `DialogoCritico` con `catch` y `porCodigo`, `Parametro` restablecer)
- `frontend/src/pantallas/ajustes-a/S05Creditos.tsx` (manejadores 1 a 4 y V-18 con «Desactivar»)
- `frontend/src/pantallas/ajustes-a/S02Presupuesto.tsx` (manejador 5)
- `frontend/src/pantallas/ajustes-a/S03Ingresos.tsx` (manejador 6)
- `frontend/src/pantallas/ajustes-a/textos.ts` (`desactivar`)
- `frontend/pruebas/comidas-compras/pantallas.escritorio.spec.ts` (prueba de P-10 Tendencias)
- `frontend/pruebas/ajustes-a/ajustes.escritorio.spec.ts` (4 pruebas áreas R5-01)
- `frontend/pruebas/gastos-presupuesto/correccion-bis.escritorio.spec.ts` (R5-05)
- `docs/estado/wf3/frontend-min-fix.md` (este reporte)

## Verificación

- **C3-03 build**: `npx tsc --noEmit -p .` → rc 0, sin salida. `npm run build` → «precache 81 entries», `dist/sw.js`, rc 0.
- **U5-07 (P-10)**: la prueba «P-10 Tendencias (core R5-01)» pasa en `movil` y `escritorio`, con los encabezados de tabla `Lun 5 oct 2026` y `Lun 12 oct 2026`, el texto del SVG y 0 híbridos y 0 ISO en `main`.
- **U4 §0.3 y §0.6 (áreas R5-01)**: `npx playwright test pruebas/ajustes-a -g "R5-01"` → 8 passed (4 pruebas × 2 proyectos), 0 `pageerror`, `erroresTrafico` = [].
- **R5-05**: 3 corridas completas de diagnóstico → 241/241/240 passed (la corrida 3 falló en `:79`, con foco fuera del bloque). Después de la corrección: `--repeat-each=5` → 5 passed, más la sonda de cuadro tardío, 3 de 3.
- **Corrida completa 1, dentro de `make verify-wf3`**: rc 0. Llega a «== WF-3: verificación completa OK (Go, cobertura K25, contrato C1 e interfaz)», con Playwright en **241 passed, 4 skipped, 0 failed** (5.2 min).
- **Corrida completa 2** (`npx playwright test`): rc 0, con **241 passed, 4 skipped, 0 failed** (5.2 min). Los 4 saltados son los atajos de teclado en `movil`, que ya se saltaban antes. Con las corridas de diagnóstico 1 y 2, son 4 corridas completas con 0 fallas.
- **U5-06 hex**: `git diff -U0 -- frontend/src | grep '^+' | grep -E '#[0-9a-fA-F]{3}([0-9a-fA-F]{3})?\b'` → 0 líneas (rc 1). Solo hay clases de token (`text-role-error`, `text-text-link`).
- **U5-07 inglés**: en las líneas nuevas de `frontend/src`, los únicos literales en inglés son identificadores de código (`'delete'` como método HTTP y `'error'` como rol), no texto visible. Los textos nuevos son «Desactivar», V-18 y el §0.3, todos de U4.
- **U5-08**: no hay llamadas nuevas fuera de `openapi.yaml`; `PATCH /api/v1/promotions/{promotion_id}` ya está en C1 (`updatePromotion`). La red usada fue solo `localhost:4173`.

## Supuestos

- S-1. Para «Nueva tasa», «Nuevo crédito» y «Nueva promoción», el 409 `already_exists` se muestra con V-12 («Ya existe uno con ese nombre.»), igual que `S06Tiendas.tsx:225` y el validador de plazo duplicado de `TasasPlazo`. U4 no nombra el V-xx de esos 409.
- S-2. Un 4xx de validación sin código propio muestra el mensaje del servidor, como ya hace `Editable`. Solo los errores de conexión y los 5xx usan el texto §0.3 con Reintentar.
- S-3. V-18 deja «Desactivar» como acción principal y quita el verbo crítico, porque repetir el borrado volvería a dar 409. «Cancelar» sigue con el foco inicial (§0.4).
- S-4. La posición del aviso en las hojas de alta es el pie, sobre «Guardar», para que se vea en móvil sin hacer scroll. U4 §0.3 no fija la posición dentro de una hoja.

## Bloqueos

Ninguno nuevo.

## Pendientes

- P-1. **El servidor Go no trae estas correcciones.** El que corre en `:8443` sirve el frontend embebido de su compilación anterior. No lo reinicié porque otros agentes de WF-3b lo usan. Para que las capturas U5 y C3-08 las incluyan, el integrador debe recompilar el binario con el `frontend/dist` nuevo y reiniciarlo (`ops/dev-stop.sh` y `ops/dev-run.sh`).
- P-2. Siguen abiertos fuera de mi alcance: R5-02, R5-03, R5-04, H-13, R3-01, R4-03 a R4-05 y H-07 a H-09 (§15.4).
- P-3. `frontend/tsconfig.tsbuildinfo` sigue sin rastrear (solicitud 4 de §15.8).

## Propuestas para Axel

- PR-1 (fuera de mi alcance, dueño `frontend-base`). Robustecer `VistaBD.cerrarGrupo`: cancelar el `requestAnimationFrame` pendiente en `abrirGrupo`, o no mover el foco si el bloque ya se volvió a abrir. En uso real (Esc y Ctrl+/ dentro del mismo cuadro) es casi imposible, pero es el mismo origen de R5-05.
- PR-2. Llevar a `DialogoCritico` de ajustes-b (`ajustes-b/comun`) el mismo trato de `porCodigo`, si en algún momento lo necesita. Hoy usa `g.ejecutar` y la revisión lo da por bueno.

## Consumo aproximado

- Leído: ~135 KB, por partes y con grep: revisiones r5, §15.4/§15.8, extractos de U4, archivos tocados y logs.
- Escrito: ~28 KB, en código, pruebas y este reporte.
- Acciones: ~60.
- Corridas: 3 completas de diagnóstico (unos 16 min), 1 dentro de `make verify-wf3` y 1 completa final.
