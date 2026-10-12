# revisor-frontend-areas-r1 · WF-3a (12-oct-2026)

## Hecho

Revisé de forma independiente a `front-creditos-salud-reportes`, `front-ajustes-a`, `front-ajustes-b` y `front-variantes` contra U5-06, U5-07, U5-09, C3-03, DESIGN.md, U1 y U4. No corregí nada. El informe completo está en `docs/estado/wf3/REVISION-frontend-areas-r1.md`.

**Veredicto: NO APROBADO.** Quedan 4 altos, 8 medios y 5 bajos.

## Archivos

- `docs/estado/wf3/REVISION-frontend-areas-r1.md`
- `docs/estado/wf3/revisor-frontend-areas-r1.md` (este reporte)

## Verificación

Los comandos son míos; los temporales estuvieron en `datos/tmp/rev-fa-r1/` y ya se borraron.

- **tsc**: `npx tsc --noEmit` sale 0. `npm run typecheck:variantes` sale 0.
- **Build**: `npx vite build --outDir datos/tmp/rev-fa-r1/dist` dice «✓ built». `npm run build-storybook` termina bien.
- **Playwright**: con una configuración propia (preview en el puerto 4795), `ajustes-a ajustes-b creditos-salud-reportes` da 60 passed.
- **C3-03**: CUMPLE.
  - `npm run verificar:rutas` da 201 rutas, 0 fuera de contrato y 0 sin registrar.
  - Cruce independiente de los literales `/api/…` de las 3 carpetas (39, 58 y 47 llamadas): 0 fuera y 0 sin registrar.
- **U5-07**: CUMPLE en lo medido. Escaneé en ejecución 132 páginas (33 rutas × 2 temas × 2 tamaños): 0 fechas ISO, 0 textos en inglés, 0 identificadores crudos y 0 formatos de dinero incorrectos.
- **U5-06**: NO CUMPLE.
  - Lo que pasa: 0 hex en el código y 0 pares de contraste bajo umbral en las pantallas.
  - Lo que falla:
    - Objetivos táctiles de menos de 44 px y colores de ECharts fuera de C2 en componentes de la base (H-03).
    - Dinero sin `tabular-nums` (H-09).
    - FV-01, FV-04 y FV-05 (H-10).
- **U5-09**: no aplica a estas áreas.
- **U4**: la tabla de sección contra campos de las 14 secciones está en el informe.
  - Las claves `reports.*` y `alerts.*` (11) se editan con `PATCH /settings/reports` y `/settings/alerts`; coinciden 11 de 11 con la migración.
  - Faltan dos vistas previas que U4 §0.2 exige: H-01 y H-02.
- **U1, de P-17 a P-25**: de 59 elementos revisados faltan 3 (H-12 y H-13).
- **Datos de prueba contra C1**: validé con ajv 55 respuestas distintas de la API falsa. Una no valida (H-04).
- **Variantes**: `node scripts/contraste-variantes.mjs` da 5 230 mediciones y 0 fallas nuevas, pero los 5 hallazgos FV conocidos se reproducen.
- **`make verify-all`**: sale 2 por gofmt en `backend/internal/reportes/`, un paquete sin rastrear de otro agente del backend. No toca este alcance.

## Supuestos

- En lugar de `npm run build` usé `tsc --noEmit` más `vite build` con salida temporal. Así no reescribí `dist/` ni los generados de `npm run preparar`, que son de frontend-base.
- Ubiqué el «Antes → Después» según U4 §0.2: una vista previa PROVISIONAL solo vale si U4 o un bloqueo la declaran así.
- Clasifiqué H-12 como media y no alta: es un porcentaje de presentación, no dinero ni una cifra del oráculo.

## Bloqueos

Ninguno nuevo.

## Pendientes

La ronda 2 debe comprobar H-01 a H-17. Las correcciones van a cargo de cada agente dueño, según la tabla del informe.

## Propuestas para Axel

- Que `npm run verificar:variantes` cuente como falla cualquier hallazgo FV abierto antes de agregarlo a `verify-wf3-*`.
- Que los datos de prueba se tipen con `satisfies Respuesta<…>`.

## Paquetes instalados

Ninguno. Usé el ajv que ya estaba en `frontend/node_modules`.

## Consumo aproximado

Unas 55 acciones, unos 230 KB leídos y unos 25 KB escritos.
