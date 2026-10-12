# revisor-tablero-frontend-r2 · WF-3c · 12-oct-2026

Revisor independiente, ronda 2 de 3 (no corrige). Informe completo: `docs/estado/wf3/REVISION-tablero-frontend-r2.md`.

## Hecho
- HECHO: H-1 y H-2 de la ronda 1 (medios) están cerrados, con pruebas propias en las 4 combinaciones (412 táctil y 1440, claro y oscuro).
- HECHO: hallazgo nuevo **N-1 (media)**. B-1 sin dato (la semilla no trae ingresos) deja `role="meter"` sin `aria-valuenow` y axe da `aria-required-attr` (crítico) en las 4 combinaciones (`frontend/src/componentes/Barometro.tsx:54`).
- HECHO: hallazgo nuevo **N-2 (baja)**. G-4 y G-5 no muestran su título en carga ni en error (`Tablero.tsx:67-68`, `Tendencias.tsx:57,93`).
- HECHO: los bajos H-3 a H-7 de la ronda 1 siguen abiertos, sin cambios en el código.
- Dictamen: **NO APROBADO** (0 altos, 1 medio, 6 bajos).

## Archivos
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/REVISION-tablero-frontend-r2.md` (informe).
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/revisor-tablero-frontend-r2.md` (este reporte).
- Temporales en `datos/tmp/rev-tablero-r2/` (build, configuración y spec de Playwright, observaciones JSON, bitácora de verify-all): borrados al terminar.

## Verificación
- U5-06 / U5-07 · `npx tsc --noEmit`: 0 errores. `npx vite build --outDir ../datos/tmp/rev-tablero-r2/dist`: listo, 80 entradas precacheadas.
- U3-03 / U5-06 · Playwright `pruebas/inicio` con mi configuración (puerto 4187, 412 táctil y 1440): 32 de 32.
- U5-06 · spec propia: 18 de 18. Barómetro 25/35/45 % con el color del rol correcto en valor, tramos e insignia; dígito `tabular-nums` 700; `aria-valuetext` en español. Gráficas: un eje Y, barras de 6.4/5.4 px (SVG) y 24 px (HTML), 2 patrones en G-1. Regiones de tabla con nombre, `tabindex` 0 y anillo `--ui-border-focus`. axe 0 serios con valor, con tablas, con 500 y con 501/404; **1 crítico en el vacío (N-1)**.
- U5-07 · 0 textos en inglés, 0 fechas ISO, 0 montos sin tabulares, 0 controles < 44 px con puntero grueso, sin desbordar, en todos los estados medidos.
- U3 2.0 §2 · con 501 y con 404, B-1, G-1 y G-6 = «Disponible cuando el servidor se actualice» (`text-caption`, rol info, icono de 20 px), sin «Reintentar» ni cifra.
- C3-03 · `npm run verificar:rutas`: 212 rutas invocadas, 0 fuera de `openapi.yaml`, 0 sin registrar (221).
- DESIGN / C2 · grep de hex, `rgb(`/`oklch(`/`hsl(`, `dark:` y `style` con color en 8 archivos: 0.
- Global · `make verify-u5`: EXIT 0. `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: EXIT 0.

## Supuestos
- Las respuestas de C1 con valor del barómetro y de las gráficas son SINTÉTICAS: ingreso de prueba, no del oráculo. Solo se usaron para medir formas y colores.
- Califiqué N-1 como media, igual que H-2 de la ronda 1: es una regla seria o crítica de axe y contradice el artefacto `docs/design/capturas-u5/axe/P-01-*.json`.

## Bloqueos
Ninguno nuevo.

## Pendientes
- Corregir N-1 (media). Los bajos N-2 y H-3 a H-7 quedan a criterio del orquestador.
- `frontend/dist/` del servidor de desarrollo no se tocó. La próxima captura de P-01 (U5-03) mostrará N-1 si no se corrige antes.

## Propuestas para Axel
- PROPUESTA: en DESIGN §8.9 (Accesibilidad), indicar que sin dato se usa `role="img"` con el texto accesible y sin `aria-valuenow`.
- PREGUNTA: con un servidor anterior a C1 v1.5, K-4 muestra la falla de red porque `getDashboardSettings` responde 501. Así lo pide U3 (PROVISIONAL). ¿Se mantiene?

## Paquetes instalados
Ninguno.

## Consumo aproximado
30 acciones · ~130 KB leídos · ~40 KB escritos (spec y configuración temporales, informe y reporte).
