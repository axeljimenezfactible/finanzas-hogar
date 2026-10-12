# correccion-tablero-frontend-r2 · WF-3c · 12-oct-2026

Agente `correccion-tablero-frontend-r2` (Opus). Encargo: el hallazgo N-1 (media) de `docs/estado/wf3/REVISION-tablero-frontend-r2.md`. Los bajos (N-2, H-3 a H-7) no estaban en el encargo y no se tocaron.

## Hecho

- **HECHO · N-1 (U5-06, U5-03; DESIGN §8.9 Accesibilidad).** `frontend/src/componentes/Barometro.tsx`: el contenedor del dibujo ya no es un `meter` sin `aria-valuenow`.
  - Con valor (`ratio_pct` y `status` de C1): sin cambios, `role="meter"`, `aria-label` = título, `aria-valuemin` 0, `aria-valuemax` = escala de C1, `aria-valuenow` = valor, `aria-valuetext` = «Capacidad de endeudamiento [n] % de [tope] %: [color]».
  - Sin dato (la semilla, sin ingreso): `role="img"` con `aria-label` = el texto accesible que ya armaba `Tablero.tsx` («Capacidad de endeudamiento sin dato, tope 40 %»), sin ningún `aria-value*`. Lo visible no cambia (arco con tramos y marcas, sin aguja, dígito «—», motivo y acción).
  - Lo que DESIGN §8.9 no dice (el rol sin dato) queda PROVISIONAL con bloqueo nuevo `docs/estado/bloqueos/WF3c-barometro-sin-dato.md` (opciones A hecha, B, C descartada por inventar un 0 %).
- **HECHO · pruebas** (`frontend/pruebas/inicio/tablero.escritorio.spec.ts`):
  - Prueba de la semilla (bloques en orden, A-01…A-14): B-1 sin dato pide 0 `meter`, una imagen con nombre exacto «Capacidad de endeudamiento sin dato, tope 40 %» (`data-prueba="barometro"`) y ningún `aria-valuenow`/`-min`/`-max`/`-valuetext`; al final corre **axe WCAG 2.0–2.2 A/AA (0 serios ni críticos) sobre el vacío de la semilla**. Corre en `movil` 412×915 táctil y `escritorio` 1440×900, en claro y oscuro = 4 combinaciones.
  - Prueba A-18 (con valor): además del `aria-valuetext`, comprueba `aria-valuenow` = el valor de C1 en los 4 casos.

## Archivos

- `frontend/src/componentes/Barometro.tsx` (atributos ARIA según haya valor; comentario con el bloqueo)
- `frontend/pruebas/inicio/tablero.escritorio.spec.ts` (B-1 sin dato por rol imagen, axe en la semilla, `aria-valuenow` en A-18)
- `docs/estado/bloqueos/WF3c-barometro-sin-dato.md` (nuevo)
- `docs/estado/wf3/correccion-tablero-frontend-r2.md` (este reporte)

## Verificación

Temporales en `datos/tmp/corr-tablero-r2/` (borrados al terminar).

| Criterio | Comando | Resultado |
|---|---|---|
| tipos | `cd frontend && npx tsc --noEmit -p .` | 0 errores |
| build | `npm run build` (reescribe `frontend/dist/`, que sirve el servidor de desarrollo) | listo; 80 entradas precacheadas |
| **U5-06 / U5-03 · N-1** | `npx playwright test pruebas/inicio` | **32 de 32**; la prueba de la semilla (B-1 sin dato = imagen, 0 `meter`, axe 0 serios ni críticos) y la de axe con las 6 gráficas pasan en `movil` y `escritorio`, claro y oscuro (4 combinaciones) |
| Control negativo | build del `Barometro.tsx` anterior en `datos/tmp/…/dist-viejo`, `vite preview` en 4187 y la prueba de la semilla | **4 de 4 fallan** en las 4 combinaciones (`getByRole('meter')` esperado 0, recibido 1): la prueba ya no da por bueno el `meter` sin valor |
| A-18 con valor | dentro de `pruebas/inicio` | 30, 31, 40 y 41 %: `meter` con `aria-valuenow` = valor y `aria-valuetext` de siempre |
| DESIGN / C2 | grep de hex, `rgb(`, `oklch(`, `hsl(`, `dark:` en `Barometro.tsx` | 0 |
| S6-01, S6-02, S6-03/04, C3-03, U5-06/07/08/09, U5 sin servidor | `make verify-wf3 GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | **EXIT 0**: Go OK (33 paquetes), cobertura K25 OK, 221 operaciones en 176 rutas con 0 diferencias, C3-03 212 rutas (0 fuera de `openapi.yaml`, 0 sin registrar), hosts externos 0, e2e **259 pasan, 4 omitidas, 0 fallan**, U5-06 Inicio 0 controles < 44 px, U5 sin servidor (U5-01 164/164, U5-03 0, U5-04 0, U5-08 0, U5-11 0) |
| Global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | **EXIT 0** («Verificación ampliada OK») |

No corrí `make verify-u5-servidor` ni `playwright.u5.config.ts`: escriben en `docs/design/capturas-u5/`, fuera de mi alcance (ver Pendientes).

## Supuestos

- HECHO que se apoya en la PROPUESTA del revisor (REVISION r2 §8): sin dato, `role="img"` con el texto accesible. DESIGN §8.9 no lo dice; queda PROVISIONAL con el bloqueo `WF3c-barometro-sin-dato.md`. No se tocó DESIGN (no es de mi alcance).
- El texto accesible sin dato es el mismo que ya armaba `Tablero.tsx` (`T.b1.sinValorTexto`); no se agregaron textos ni cifras.

## Bloqueos

- Nuevo: `docs/estado/bloqueos/WF3c-barometro-sin-dato.md` (rol del barómetro sin dato; opción A hecha, B posible, C descartada).

## Pendientes

- Regenerar la evidencia U5 contra el servidor de desarrollo (`make verify-u5-servidor`, ≈ 9 min) para que `docs/design/capturas-u5/axe/P-01-*.json` refleje B-1 sin dato con este cambio; lo hace quien es dueño de `docs/design/capturas-u5/` (WF-3b o el integrador). `frontend/dist/` ya está reconstruido con la corrección.
- Bajos abiertos de la revisión, fuera de este encargo: N-2 (G-4 y G-5 sin título en carga y error), H-3 a H-7.

## Propuestas para Axel

- PREGUNTA (bloqueo `WF3c-barometro-sin-dato.md`): ¿se fija en DESIGN §8.9, Accesibilidad, «Sin dato: `role="img"` con el texto accesible, sin `aria-valuenow`» (opción A, la que quedó), o se prefiere el texto en `sr-only` sin rol (B)?

## Paquetes instalados

Ninguno.

## Consumo aproximado

≈ 30 acciones · ≈ 65 KB leídos · ≈ 35 KB escritos.
