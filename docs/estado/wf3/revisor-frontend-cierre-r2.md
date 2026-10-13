# revisor-frontend-cierre-r2 · WF-3a cierre (12-oct-2026)

Agente: `revisor-frontend-cierre-r2` (independiente; no corrige). Informe completo: `docs/estado/wf3/REVISION-frontend-cierre-r2.md`. No hice `git add` ni `commit` y no edité archivos del producto.

## Hecho
- **Veredicto: NO APROBADO.** 0 altos, 1 medio y 5 bajos.
- M-1 y M-2 de la ronda 1 están corregidos; lo comprobé contra el servidor real.
- El medio nuevo (M-1) es el mismo defecto de M-2 en D-3 de P-01 Detalle (`Detalle.tsx:92-98`): «Hogar $973.75» mientras `listPersons` no responde. Es un defecto previo; la corrección lo había enviado como solicitud.
- Los bajos: B-1 (F-3A0 pasa aunque su paso dice FALLA), B-2 (U4 P-26.02 sin el bloque y el bloqueo `C3-ajustes-tablero` sin la decisión 4), B-3 (`tableroNoCargo` fuera del estado común) y B-4 (evidencia versionada de U5-04 y matriz de C3 sin regenerar) siguen abiertos desde la ronda 1. B-5 es nuevo: un comentario viejo en `c3-09.spec.ts:55-56`.

## Archivos
- `docs/estado/wf3/REVISION-frontend-cierre-r2.md` (nuevo) y este reporte.

## Verificación (mis comandos; temporales en `datos/tmp/`, ya borrados)
- **(a)** `npx tsc --noEmit` da 0. `npx vite build --outDir datos/tmp/rfc2/dist` da 80 entradas, y `diff -rq` contra `frontend/dist` da 0. `npx playwright test pruebas/ajustes-a/ pruebas/inicio/` da **113 passed, 1 skipped**, rc 0.
- **U4 P-26.02 / §0.2** contra el servidor real (binario compilado de la fuente actual, copia limpia en 8445):
  - Las 5 claves aparecen con su unidad corta, y ninguna de las descripciones largas de `unit` (123, 113, 57, 88 y 112 caracteres).
  - Vista previa de tramos y escala «De 0 a 60 % → De 0 a 70 %». Esc no escribe. Enter guarda 45. «Restablecer» vuelve a 40 con `is_initial` = true. V-23 muestra «del 1 al 39».
- **U5-04** (recaptura de P-01 en 4 combinaciones): Resumen 579, 579, 610 y 610 pares, y Detalle 134, 134, 161 y 161, todos con **0 bajo umbral**. 20 de 20 pares recalculados con mi propia fórmula coinciden (±0.01).
- **U5-06**: el eje mide 10 × 10. P-01 Detalle Android mide 412 de 412. P-26.02 mide 412 de 412, con 0 controles chicos.
- **U5-07**: inglés [] y 0 fechas ISO en P-01 Resumen, P-01 Detalle y P-26.02.
- **(e)**: 0 hex y 0 `rgb(`, `hsl(` u `oklch(` en las 88 líneas agregadas de `frontend/src`.
- **C3-03**: «213 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0». `updateDashboardSettings` sí se llama, y `c3-integracion.py` deja 0 sin cobertura. Su salida 1 se debe solo a la matriz sin regenerar (B-4).
- **C3-08**, cada archivo en una copia limpia: `flujos.spec.ts` da 7 passed (F-1 con 3 toques). `pendientes.spec.ts` da 4 passed (F-3A3, la ruta A completa, con 3 toques).
- **C3-09**: `c3-09.spec.ts` da «120 cifras … diferencias != 0: 0», 13 de ellas del tablero y P-17. Con `/api/v1/persons` retenida, K-4 nunca dice «Hogar $973.75», pero **D-3 sí** (M-1).
- **Global**: `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` da rc 0. `make verify-u5` da «bajo umbral = 2», por la evidencia vieja (B-4).

## Supuestos
- Probé contra un binario compilado de la fuente actual y no contra `backend/bin/finanzas`. Ese binario es de las 17:54 y la corrección del backend cambió `.go` a las 18:54. Con la semilla, que no tiene líneas archivadas, las cifras no cambian.
- M-1 lo califiqué como **media** y no como alta (aunque es una cifra distinta del oráculo) ni como baja (aunque es previo y está fuera de los archivos tocados). Así lo calificó la ronda 1 para el mismo defecto en K-4.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Ninguno de mi alcance.

## Propuestas para Axel
- Que el orquestador decida si D-3 y el Calendario entran en el alcance de la próxima corrección (es la misma corrección de 3 renglones que ya se hizo en K-4) o si quedan como decisión registrada.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- ≈ 58 acciones, ≈ 120 KB leídos, ≈ 16 KB escritos.
