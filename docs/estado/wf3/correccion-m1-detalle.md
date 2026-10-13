# Corrección M-1 · D-3 de P-01 Detalle con `listPersons` pendiente (12-oct-2026)

Origen: hallazgo M-1 (media, C3-09 y U5-06) de `docs/estado/wf3/REVISION-frontend-cierre-r2.md`. Sin `git add` ni commit. No se tocó `backend/` ni `docs/design/`.

## Qué estaba mal
`D3` en `frontend/src/pantallas/inicio/Detalle.tsx` pintaba las filas antes de que respondiera `listPersons`. Mientras tanto `idK` (Karen) era `undefined`, así que el pago entero se sumaba al hogar. Las 3 filas decían «Kueski $973.75 · Hogar $973.75 · Pendiente» en lugar de «Hogar $598.82 · Karen $374.93» (oráculo R-13).

## Qué se cambió
| Archivo | Cambio |
|---|---|
| `frontend/src/pantallas/inicio/Detalle.tsx` (`D3`) | Mismo patrón que K-4 (`Tablero.tsx`). (1) Mientras `c.personas.datos === undefined`, se muestra el esqueleto de 5 filas: ni filas, ni montos, ni «Hogar». (2) Si `listPersons` falla o nunca se cargó sin conexión, se muestra el estado de error de D-3 (`FallaBloque` con `t.d3.falla` y «Reintentar» → `c.personas.recargar()`). (3) «Hogar» suma solo las partes de personas que están en `listPersons` y no son Karen. Si una parte es de una persona desconocida, esa fila se queda sin la línea «Hogar · Karen». «Disponible por línea» no depende de `listPersons` y se sigue viendo. No hay textos nuevos (se reutiliza `t.d3.falla`). |
| `frontend/pruebas/inicio/inicio.escritorio.spec.ts` | 3 pruebas nuevas en D-3 (ver abajo). Corren en `movil` y `escritorio`. |
| `docs/estado/bloqueos/U3-tablero-k4-personas.md` | Se amplía a D-3 (opción A, PROVISIONAL), como pide la corrección que sugiere la revisión. |

## Pruebas nuevas
Usan la API simulada de la convención §8 (`ApiFalsa` + `instalarInicio` con `context.route`). Es la misma técnica que la prueba de K-4.
1. **Retención**: la ruta retiene `/api/v1/persons` en `/inicio?sub=detalle`. Ya respondieron los pagos (`status=pending`) y el disponible («Lectura B $6,905.50»), pero D-3 tiene 0 filas, sin «$973.75» y sin «Hogar». Al soltar la petición aparecen las 3 fechas y «Hogar $598.82 · Karen $374.93» exactamente 3 veces. Un `MutationObserver` sobre todo el documento confirma que «Hogar $973.75» **nunca** apareció.
2. **Falla**: con `/api/v1/persons` → 500, D-3 muestra «No se pudieron cargar los pagos. Revisa que estés en la red de casa.», sin «$973.75» ni «Hogar», y el disponible sigue visible. «Reintentar» lo recupera y aparecen las 3 filas con las dos partes.
3. **Persona desconocida**: un pago con una parte de `person_id` 9, que no está en `listPersons`. La fila sale con su fecha, su monto y «Pendiente», sin la línea «Hogar».

**Control positivo**: con el `Detalle.tsx` anterior (build temporal), las 3 pruebas fallan. En la de retención aparecen 3 filas. En la de falla, D-3 dice «Kueski $973.75 · Hogar $973.75 · Pendiente» ×3, el mismo texto que reprodujo el revisor. Después se restauró la versión corregida y se volvió a compilar.

## Verificación
- `npx tsc --noEmit`: rc 0.
- `npm run build`: rc 0 («precache 80 entries»). `dist/` sale de la fuente corregida.
- `npx playwright test pruebas/inicio`: **42 passed**, rc 0. Son las 36 que había y las 3 nuevas × 2 proyectos.
- Revisé las 70 líneas agregadas (fuente y prueba). Hex o `rgb(`/`hsl(`/`oklch(`: **0**. Textos visibles en inglés: **0**. Solo coinciden identificadores del API en la prueba (`/api/v1/scheduled-payments`, `status=pending`), y la fuente no agrega cadenas visibles.

## Pendiente (fuera de este encargo)
- La insignia de `Calendario.tsx:33-42` («Falta la parte de Karen / del hogar») sigue usando `karen(c)` sin esperar a `listPersons`. Los montos vienen del servidor; solo se notaría con un pago a medias. La revisión sugiere el mismo arreglo.
- B-5 de la revisión: el comentario viejo en `frontend/pruebas-u5/flujos/c3-09.spec.ts:55-56` no se tocó.
