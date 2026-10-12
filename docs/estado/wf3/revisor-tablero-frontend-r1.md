# revisor-tablero-frontend-r1 · reporte

Agente `revisor-tablero-frontend-r1` (WF-3c, ronda 1 de 3, 12-oct-2026). Revisor independiente del trabajo de `front-tablero`; no corrige. Informe completo: `docs/estado/wf3/REVISION-tablero-frontend-r1.md`. Sin `git add` ni `git commit`; sin Telegram; sin red.

## Hecho
- HECHO: reejecuté tsc, build (en un directorio temporal, sin tocar `frontend/dist`), C3-03 (`verificar:rutas`), hosts externos, las 26 pruebas de `pruebas/inicio` (412 táctil y 1440, claro y oscuro) y la suite completa de `frontend/pruebas` contra el build temporal.
- HECHO: pruebas propias (temporales) del barómetro con 25 %, 35 % y 45 %, medidas de gráficas (SVG y HTML), idioma, fechas ISO, objetivos ≥ 44 px, iconos, cifras tabulares, axe WCAG 2.0–2.2 A/AA, rejilla de 1440 y estados carga / 500 / 501 / vacío.
- HECHO: tabla de fidelidad a U3 2.0 (elemento → presente) y dictamen: **NO APROBADO**, 0 altos, 2 medios (H-1 estado «servidor anterior» ausente; H-2 axe serio `scrollable-region-focusable`), 5 bajos.

## Archivos
- `docs/estado/wf3/REVISION-tablero-frontend-r1.md` (informe) y este reporte. Temporales en `datos/tmp/rev-tablero-r1/` (borrados al terminar).

## Verificación
- U5-06: barómetro de §8.9 con tramos de rol (valor pintado = `--ui-role-success/warning/error` para 25/35/45 en claro y oscuro), dígito `tabular-nums` 32 px 700, insignia con icono y palabra; 0 controles < 44 px (táctil); 0 iconos fuera de norma; rejilla de 1440 según U3 §2.4. axe: 0 `color-contrast`, **6 nodos `scrollable-region-focusable` (serio) por combinación** (H-2).
- U5-07: 0 textos en inglés (incluido el texto de los SVG), dinero `$1,234.56` tabular, fechas `10-oct-2026` / `Vie 30 oct 2026` / `oct-2026`, 0 fechas ISO, en las 4 combinaciones.
- U5-09: no aplica.
- C3-03: `npm run verificar:rutas` → 212 rutas invocadas, 0 fuera de `openapi.yaml`, 0 sin registrar (221).
- U3 2.0: 16 de 17 elementos presentes; falta «Disponible cuando el servidor se actualice» (H-1). Estados: carga, error y vacío en los 11 bloques.
- DESIGN §8 / §8.5: un eje Y; barras 6.4 / 5.4 px (G-1) y 24 px (G-2, G-3); leyenda de G-1 con patrón; dinero en texto o en la tabla; marcador de G-6 de 8 px sin anillo (H-5).
- tsc 0 errores; build listo; `pruebas/inicio` 26 de 26; suite completa: **253 passed, 4 skipped, 0 failed** (5.5 min).

## Supuestos
- «U3 2.0» = `docs/design/U3-dashboard.md` versión 2.0, evaluado con los criterios de `docs/criterios/U3.md` 1.4.
- §2.3 «Qué miro para prevenir» no es un bloque de la interfaz (no figura en la tabla de bloques ni en el orden de §2.4).
- Las «4 combinaciones» son 412×915 táctil y 1440×900, en claro y oscuro (U5 §1).

## Bloqueos
- Ninguno nuevo.

## Pendientes
- Ronda 2 tras la corrección de H-1 y H-2 (y, si se quiere, de los bajos).

## Propuestas para Axel
- PREGUNTA: ¿la tabla «Qué miro para prevenir» debe verse en Inicio? Hoy solo vive en U3 §2.3.

## Paquetes instalados
Ninguno.

## Consumo aproximado
≈ 40 acciones; ≈ 230 KB leídos; ≈ 30 KB escritos (temporales incluidos).
