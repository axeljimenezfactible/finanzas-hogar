# revisor-tablero-frontend-r3 · WF-3c · 12-oct-2026

Agente `revisor-tablero-frontend-r3` (revisor independiente, ronda 3 de 3; no corrige). Informe completo: `docs/estado/wf3/REVISION-tablero-frontend-r3.md`.

## Hecho

- Revisé `front-tablero` con sus correcciones r1 y r2 contra U5-06, U5-07 (U5-09 no aplica), U3 2.0, DESIGN §7.6, §8, §8.5 y §8.9, y C3-03.
- **Dictamen: APROBADO** (0 altos, 0 medios). N-1 (media, r2) **cerrado**: sin dato, el barómetro es `role="img"` con «Capacidad de endeudamiento sin dato, tope 40 %» y sin `aria-value*`; axe 0 serios en la semilla en las 4 combinaciones.
- Bajos: 2 nuevos (N-3 esqueleto de carga común con `aria-label` sin rol, fuera del alcance; N-4 eje de la aguja de 16 px en lugar de 10 px) y 6 que siguen (N-2, H-3 a H-7).

## Archivos

- `docs/estado/wf3/REVISION-tablero-frontend-r3.md` (informe)
- `docs/estado/wf3/revisor-tablero-frontend-r3.md` (este reporte)
- Temporales en `datos/tmp/rev-tablero-r3/` (build, configuración y specs propias, observaciones JSON): **borrados**.

## Verificación

| Línea | Comando | Salida resumida |
|---|---|---|
| tipos | `npx tsc --noEmit -p .` | EXIT 0 |
| build | `npx vite build --outDir ../datos/tmp/rev-tablero-r3/dist --emptyOutDir` | EXIT 0; 80 entradas precacheadas |
| U5-06 / U5-07 (Inicio, 4 combinaciones) | Playwright con configuración propia (412×915 táctil y 1440×900; claro y oscuro) sobre `frontend/pruebas/inicio` | 32 de 32 |
| U5-06 · barómetro DESIGN §8.9 | spec propia R1 (25, 35, 45 %) | colores de valor, tramos, insignia, dígito y aguja = tokens de rol en claro y oscuro; dígito 700 `tabular-nums`; `aria-valuetext` en español; axe 0 serios |
| U5-06 · sin dato (N-1) | spec propia R2 (semilla) | 0 `meter`; 1 `img` con el nombre; axe 0 serios; axe con todas las reglas sobre B-1: 0 |
| DESIGN §8 / §8.5 | specs propias R3 y R5 | un eje Y; barras de G-1 5.4 y 6.4 px; 2 patrones en la leyenda; resúmenes con dinero; aria en español; axe 0 serios |
| U3 2.0 | specs propias R6 y R7 | 360/412/768 líneas de B-1, valor > escala, sin créditos, ingreso esperado, rejilla 768 y 1440: presentes |
| Estados (f) | spec propia R4 | carga con esqueleto (G-4 y G-5 sin título, N-2); 500 con mensaje, `alert` y «Reintentar»; 501 con «Disponible cuando el servidor se actualice»; axe en carga: `aria-prohibited-attr` del `Cargando` común (N-3) |
| U5-06 objetivos ≥ 44 px, U5-07 | ayudantes `pruebas/apoyo/medidas.ts` en todos los estados | 0 controles chicos en `movil`; 0 inglés; 0 ISO; 0 montos sin `tabular-nums` |
| C3-03 | `npm run verificar:rutas` | 212 rutas; 0 fuera de `openapi.yaml`; 0 sin registrar (221) |
| Hosts | `node scripts/hosts-externos.mjs <dist temporal>` | solo literales que no son peticiones |
| DESIGN / C2 | greps de hex, `rgb(`, `dark:`, `yAxis: [`, color en `style` en 8 archivos | 0 |
| Global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | EXIT 0 |

## Supuestos

- En escritorio (puntero fino) los controles de 36 px son válidos (`size.control.fine`); el mínimo de 44 px se mide en `movil` (táctil), como en la ronda 2.
- Las respuestas con ingreso son SINTÉTICAS (la semilla no trae ingresos); no se comparó ninguna cifra con el oráculo porque la interfaz no calcula: pinta lo que llega de C1.

## Bloqueos

- Ninguno nuevo. Sigue `docs/estado/bloqueos/WF3c-barometro-sin-dato.md` (PROVISIONAL, opción A hecha).

## Pendientes (fuera de mi alcance)

- `frontend/src/componentes/Estados.tsx:15` (frontend-base): `role="status"` en el `div` de `Cargando`, o `sr-only` en lugar de `aria-label` (N-3).
- `docs/design/capturas-u5/` (WF-3b o integrador): regenerar P-01 con `make verify-u5-servidor`; las capturas y el axe de P-01 son de antes del tablero.
- `docs/design/DESIGN.md` §8.9 (quien edita DESIGN): fijar el rol sin dato según el bloqueo.

## Propuestas para Axel

- PREGUNTA: ¿opción A (`role="img"`) u opción B (`sr-only`) para el barómetro sin dato en DESIGN §8.9?
- PROPUESTA: juntar N-3 (una línea en `Estados.tsx`, ayuda a todas las pantallas), N-4 y los bajos de `Tablero.tsx` en una sola pasada.

## Paquetes instalados

Ninguno.

## Consumo aproximado

≈ 48 acciones · ≈ 130 KB leídos · ≈ 40 KB escritos (incluye specs temporales).
