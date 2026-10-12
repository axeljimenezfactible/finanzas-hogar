# front-ajustes-a · P-26 Ajustes, secciones 01 a 07 (WF-3a, 12-oct-2026)

## Hecho
- **P-26** (`P26Ajustes.tsx`):
  - En escritorio, la columna interna con las 14 secciones y la sección a la derecha.
  - En móvil, la lista de secciones, y cada sección abre como página con «‹ Ajustes».
  - Las secciones 08 a 14 las aporta otra área por `secciones-ajustes.ts` (bloqueo WF3-ajustes-a §3).
- **Piezas comunes de U4 §0** (`comunes.tsx`):
  - Edición en el lugar: `Enter` o salir guarda, `Esc` descarta. Lleva los mensajes V-01 a V-25 y la vista previa «Antes → Después».
  - Una vista previa que falla no deja guardar. Una pendiente se rotula «Vista previa pendiente» y sí deja guardar.
  - Estados del guardado: «Guardando…», «Guardado» (`circle-check`, aria-live, 2 s) y error con Reintentar.
  - Sin conexión, toda escritura queda deshabilitada con la razón en texto.
  - Interruptores, con o sin vista previa.
  - Diálogo de acción crítica: `shield-alert`, «Cancelar» con el foco inicial y hoja inferior en móvil.
  - Parámetro K44: «Valor inicial: …», «Restablecer al valor inicial» y el aviso «Valor restablecido · Deshacer».
- **Cálculos en el navegador** (`calculos.ts`), en centavos enteros con ROUND_HALF_UP y BigInt:
  - Fechas de una partida planificada: quincenas 15 y 30 o último día (RN-06), recorte al último día del mes.
  - Reparto según RN-03 y K14.
  - Tasa como fracción exacta, con los decimales tecleados (K12).
  - Corte y pago de la tarjeta (Q2-32 A).
  - Peso de un precio según K22 y K40.
  - Tabla francesa con fracciones exactas (Q2-28 (a)).
  - Pesos normalizados.
- **01 Cuenta y roles**:
  - Usuarios y Karen sin cuenta.
  - Contraseña con V-01, V-14, V-15 y `Ctrl+Enter`.
  - Sesiones en tabla, con cierre de una (confirmación simple) y cierre por alcance como acción crítica con conteo.
  - Salir del dispositivo, indicador de app instalada y descarga del certificado.
  - 4 parámetros K44 con vista previa.
- **02 Presupuesto**:
  - Proyección de los meses 7 a 12 (K21).
  - Montos por categoría y total, con el diálogo de alcance («De este mes en adelante» preseleccionado, K2) y el total antes → después.
  - Quitar un monto como acción crítica.
  - Gastos fijos en tabla, editables en el lugar, con la vista previa 3 (a) a (g): próximas 3 fechas, n × monto, gastos fijos del mes, partes por persona, «Se atribuye a», aviso de P-07 y el interruptor «activo».
  - Nuevo gasto fijo (`Alt+N`).
  - Despensa programada con la vista previa 4: fechas y `previewCreditPurchase` antes → después, más «El disponible baja».
- **03 Ingresos**: quincena del 15 y del 30 y «Rige desde», con la vista previa del ingreso del mes. Muestra a Ximena, el estado vacío con alta y «Registrar ingreso extra» (hoja de P-03, pestaña Ingreso).
- **04 Categorías**:
  - Tabla y tablero por reparto.
  - Nombre, padre, tipo, reparto con vista previa (ejemplo de $199.00 y gastos fijos sin reparto propio).
  - «Activa» con aviso de archivo y aviso de Despensa y Comida fuera.
  - Reordenar con `ListaOrdenable` y `PUT /categories/order`. Nueva categoría (`Alt+N`, V-12 ante un 409).
- **05 Créditos y condiciones**:
  - Líneas en tabla, galería y tablero, con «Falta capturar».
  - Ficha de la línea en hoja con todos los campos por tipo:
    - Límite: vista previa 5 con Δ sobre las dos lecturas.
    - Condiciones de tarjeta y tienda: vista previa 4 (a) y (b); (c) PROVISIONAL por S3-formula-tarjeta-pdf.
    - Días de corte y pago: vista previa 9.
    - Tasa del préstamo: vista previa 10 con amortización exacta.
    - Comisiones: vista previa 8.
    - Días mínimos al primer pago: vista previa 7 PROVISIONAL por SC-U4-1.
  - A/B con las dos lecturas lado a lado y «Usar la lectura X». Aviso «Pendiente confirmar» y casilla «Ya lo confirmé».
  - Tasas por plazo: vista previa 2 con $3,094.50; plazo repetido da V-12; eliminar es acción crítica.
  - Promociones: calendario, tabla y tablero, nota 11, eliminar como acción crítica.
  - Comparador: «Faltan restricciones», colchón, máximo de compras, deuda en cero, 2 parámetros K44 (tiempo máximo de 1 a 30 con V-23), 4 pesos con % normalizado, el aviso «Al menos un peso…», «Pesos iguales» y el enlace al tope en P-26.12.
- **06 Tiendas y precios**:
  - Tiendas en tabla, galería y tablero. La ficha en hoja tiene: nombre, cadena, PROFECO, sitio, dirección, grupos con `setStoreProductGroups`, membresía, lectura asistida bloqueada sin veredicto (V-19) y acción crítica para encenderla, activa y notas.
  - Fuentes de precio: manual, OCR, estado de la IA, lectura asistida general como acción crítica y el corte PROFECO, que se importa con validación V-20 contra el máximo vigente.
  - 3 parámetros K44 con la vista previa de frescura (30 días = 9-nov-2026, 0.9836; 0.4918; 8-ene-2027) y V-24 y V-25.
- **07 Recetas**:
  - Porciones base, tiempos y conteo por origen.
  - Recetario en tabla, tablero y galería: nombre y tiempos aptos («Elige al menos un tiempo de comida.»), enlace a P-12 y eliminar como acción crítica con el uso del plan.
  - Bases nutricionales con «Recargar» y su vista previa.
  - Restricciones por persona y estado de la IA.
- **Pruebas** (`frontend/pruebas/ajustes-a/ajustes.escritorio.spec.ts`; corre en 412×915 táctil y en 1440×900): índice y 7 secciones (elementos principales, sin desborde horizontal, botones ≥ 44 px en móvil), más 5 flujos: K44 con «Antes → Después» y PATCH, acción crítica con «Cancelar» enfocado, alcance del presupuesto, A/B y V-25.

## Archivos
`frontend/src/pantallas/ajustes-a/`: `indice.ts`, `secciones-ajustes.ts`, `P26Ajustes.tsx`, `S01Cuenta.tsx` a `S07Recetas.tsx`, `comunes.tsx`, `calculos.ts`, `textos.ts`.
También: `frontend/pruebas/ajustes-a/ajustes.escritorio.spec.ts`, `docs/estado/bloqueos/WF3-ajustes-a.md` y este reporte.

## Verificación
- `npx tsc --noEmit`: 0 errores en `ajustes-a`. El proyecto tiene 1 error ajeno, `pruebas/creditos-salud-reportes/datos.ts(48,79)` TS6133.
- `npm run build`: **falla por ese error ajeno**. `npx vite build` sale sin errores (chunks `P26Ajustes`, `S01`…`S07`).
- `npx playwright test pruebas/ajustes-a`: **26 passed** (13 en movil y 13 en escritorio).
- **U5-06**: 0 hex, 0 `style=`, 0 valores arbitrarios `[…]`; solo clases de C2. Foco visible en todo control y objetivos de 44 px en móvil, comprobados en la prueba.
- **U5-07**: grep de Save, Cancel, Settings, Dashboard, Loading, Error, Home, Delete, Search y Add en texto visible = 0. El único resultado es el componente `AlertDialog.Cancel`, que muestra «Cancelar». Dinero con `dinero()` y fechas con `fecha()`.
- **U5-09 y U5-13**: no aplican a Ajustes; P-26 no toca P-03 ni la cola.
- **C3-03**: `npm run verificar:rutas` da 201 rutas invocadas, 0 fuera de `openapi.yaml` y 0 sin registrar.
- **Flujos F-1 a F-6**: ninguno pasa por P-26.01 a P-26.07. F-4 solo enlaza «Falta capturar», que se comprueba en la prueba de 05.
- `make verify-all GOOSE=…`: OK.

## Supuestos
- Las pruebas van en `frontend/pruebas/ajustes-a/`, como dice CONVENCIONES §1 y §8, y usan ApiFalsa con `context.route` (§8), no msw.
- La ficha de una línea y la de una tienda se abren en `Hoja`.

## Bloqueos
`docs/estado/bloqueos/WF3-ajustes-a.md`:
1. Tipo de partida de la despensa programada.
2. Semana para el punto de equilibrio de la membresía.
3. Contrato de P-26 compartida con la parte b.

## Pendientes
- 02: vista previa 3 (h), es decir, `Ctrl+/` en varias partidas.
- 04: arrastre entre columnas del tablero con vista previa y botón «Aplicar»; hoy el tablero es solo lectura. Falta también `Ctrl+/`.
- 05: vista previa 6 (frecuencia de pagos con `previewCreditPurchase` doble); hoy dice «Vista previa pendiente». Faltan `←`/`→` y `Ctrl+/` en promociones, y altas de promoción.
- 06: calendario por día de visita; conteo «cuántos precios cambian de fresco a viejo hoy»; avance de la importación con `listJobs`.
- 07: «Nueva receta» y «Duplicar receta», porque `rutas.ts` no tiene ruta de alta de P-12; hoy cada receta enlaza a P-12.
- `mocks.ts` con msw para `npm run dev:mocks`; hoy se usan los mocks genéricos.

## Propuestas para Axel
1. Despensa programada como `fixed_expense` con `accrual_months = 2` (bloqueo §1, opción a): reutiliza el esquema.
2. Punto de equilibrio con la semana en curso de la lista (bloqueo §2, opción a).

## Paquetes instalados
Ninguno.

## Consumo aproximado
- Unas 58 acciones.
- Unos 215 KB leídos.
- Unos 190 KB escritos, **por encima del tope de 150 KB**: el código ocupa 178 KB.
