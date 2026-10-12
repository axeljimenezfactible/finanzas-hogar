# revisor-frontend-core-r2 · WF-3a · 12-oct-2026

Agente: `revisor-frontend-core-r2`, revisor independiente en la ronda 2 de 3.
- No escribió código del frontend y no corrigió nada.
- Informe completo: `docs/estado/wf3/REVISION-frontend-core-r2.md`.

## Hecho

- **Veredicto: NO APROBADO.** Quedan 1 alto (N-01), 1 medio (N-02) y 1 bajo (N-03).
- Los 7 hallazgos de la ronda 1 que tocaba comprobar (H-01..H-06 y H-10) quedan **cerrados**, verificados con mediciones propias.
- Se reejecutaron con comandos propios: build y tsc, e2e del grupo, `make verify-wf3-frontend` y `make verify-all`.
- Se escribieron dos pruebas propias en `datos/tmp`, independientes de `pruebas/apoyo/medidas.ts`:
  - recorrido de 30 vistas a 412 px táctil en claro y a 1440 px en oscuro, más la hoja F-1;
  - comprobación de tema, de las tarjetas de Inicio, de los estados abiertos y de la cola con 2 capturas y falla 500.
- Se cruzaron las rutas del frontend con los extractos de C1 mediante un script propio.

## Archivos

- `docs/estado/wf3/REVISION-frontend-core-r2.md` (nuevo).
- `docs/estado/wf3/revisor-frontend-core-r2.md` (nuevo, este archivo).
- Temporales en `datos/tmp/rev-fcore-r2/`, ya borrados:
  - `spec/rev.spec.ts`, `spec/rev2.spec.ts` y `spec/pw.config.ts`;
  - `rutas.py`;
  - los registros de build, e2e y make.

## Verificación

- **C3-03**:
  - `npm run build`: exit 0, 0 `error TS`.
  - `rutas.py`: 162 rutas literales, 0 fuera de los 170 paths de los extractos.
  - make: «202 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0».
- **U5-06**:
  - Táctil: 697 controles a 412 px y 0 menores de 44 en reposo, pero la hoja F-1 tiene «Más opciones» de 93×20 (**N-01, alta**).
  - Iconos: 915 medidos, 0 fuera de 20/24 px o de trazo 1.5. Quedan 14 controles de solo icono con 20 px (**N-02, media**).
  - Dinero: 0 importes sin `tabular-nums`.
  - Colores y medidas: 0 hex y 0 `rgb()`/`hsl()`/`oklch()` en el alcance. Queda un valor arbitrario, `max-w-[85vw]` (H-08, bajo, pendiente).
  - Inicio: 6 tarjetas exactas.
  - Los componentes de DESIGN §7 se usan en las áreas.
- **U5-07**: 0 términos en inglés (lista de 37) y 0 fechas ISO, en 61 estados medidos. Sin red se ven «Sin conexión con el servidor de casa» y «2 en cola».
- **U5-08**: hosts pedidos = `localhost:4173`. En `dist`, solo literales que no se piden.
- **U5-09**: 3 toques en las 5 mediciones del agente, que reejecuté (sin conexión; P-02 a 412 y a 1440; Inicio a 412 y a 1440).
- **U5-10**:
  - Oscuro: `body` = `#222325` (`ui.surface.page.dark`) en las 30 vistas.
  - Claro: `#FAFAFB`.
  - El interruptor recorre Claro, Oscuro y Según el sistema.
- **U5-13 pasos 1–3 y 6–7** (variante propia):
  - Con 500: «No se pudo enviar · 2 en cola» y el aviso. El reintento manual con falla vuelve a avisar.
  - Al recuperar el servidor: 2 registros con 0 duplicados (4 POST, 2 claves de 36 caracteres) y `captured_at` con `-06:00`.
- **e2e**:
  - Grupo: 63 passed.
  - `make verify-wf3-frontend`: 133 passed, exit 0.
  - Pruebas propias: 6 passed.
- **`make verify-all`**: exit 2 por gofmt en `backend/internal/reportes/reportes_test.go`. Es backend de S7, fuera del alcance revisado; queda para el orquestador.

## Supuestos

- HECHO: las etiquetas `label[for]` de 20 px de alto no son objetivos táctiles, porque el campo asociado mide ≥ 44. No se cuentan como hallazgo.
- HECHO: las flechas «Plegar» y «Desplegar» de la barra lateral con 20 px cumplen §5.2, que asigna la barra lateral a `size.icon.sm`.
- PREGUNTA (no bloquea): U5-07 incluye «Error» en la lista de palabras en inglés, pero también es español. Hoy el conteo es 0, así que no hubo que decidir.

## Bloqueos

Ninguno nuevo.

## Pendientes

- Para los dueños (ronda 3):
  - N-01: `frontend-base`.
  - N-02: `frontend-base`, `front-inicio`, `front-gastos-presupuesto` y `front-comidas-compras`.
  - N-03 (bajo): `frontend-base`.
- Bajos de la ronda 1 que siguen abiertos: H-07, H-08, H-09 y H-11.
- `make verify-all` roto por gofmt en backend de S7, fuera de este alcance.

## Propuestas para Axel

- PROPUESTA: que `pruebas/medidas.spec.ts` mida también los estados abiertos (hoja F-1 con «Más opciones», panel de cola, menú «/» y vista previa) y la regla «solo icono → 24 px». Así no se repite lo que pasó con N-01 y N-02, que la prueba no podía detectar.

## Paquetes instalados

Ninguno.

## Consumo aproximado

~46 acciones · ~120 KB leídos · ~30 KB escritos.
