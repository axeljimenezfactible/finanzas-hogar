# correccion-tablero-frontend-r1 · WF-3c · 12-oct-2026

Agente `correccion-tablero-frontend-r1` (Opus). Corrige los 2 hallazgos medios de `docs/estado/wf3/REVISION-tablero-frontend-r1.md` (H-1 y H-2). Los 5 bajos (H-3 a H-7) no estaban en el encargo y no se tocaron. Sin cifras nuevas: ningún monto ni porcentaje cambió.

## Hecho

**H-1 (media) · Estado «servidor anterior a C1 v1.5» (U3 2.0 §2, reglas del tablero; U3-03).**
- `Tablero.tsx`: `esServidorViejo(error)` = `ErrorApi` con estado 501 (`not_implemented`, routes.go) o 404 («La ruta no existe», server.go). HECHO verificado en el contrato: `getDashboardDebtCapacity` declara solo 200, 400, 401 y 422, así que en estas operaciones un 404 solo puede ser una ruta inexistente.
- `ServidorViejo`: `<p role="status">` con `T.servidorViejo` («Disponible cuando el servidor se actualice», `textos.ts:40`, que antes nadie usaba), `text-caption` (`font.size.caption`), rol información (`text-role-info` sobre `bg-role-info-surface`, el par C-073 de `contraste-componentes.csv`: 4.96 en claro y 6.11 en oscuro) e icono `Info` de 20 px. Sin cifra y sin «Reintentar».
- B-1: con servidor viejo pinta la caja de la tarjeta (`data-prueba="tarjeta-B-1"`, título con icono, `ServidorViejo` y el enlace «Ver créditos» → P-17, igual que en los demás estados de la tarjeta de U3 §1.3). Con 500 sigue igual que antes (falla + «Reintentar», A-07).
- G-1 y G-6: `EstadosGrafica` recibe `soloNueva`. Solo esos dos bloques lo activan, porque son los únicos de U3 2.0 §2 que dependen solo de una operación nueva (`getDashboardMonthlyOverview`). G-2 usa `getDebtLive` y G-3 combina `getBudgetByCategory` con `getBudgetSignals`, así que siguen con su falla normal.

**H-2 (media) · axe serio `scrollable-region-focusable` (U5-06; U5-03).**
- `GraficaInicio.tsx`: componente nuevo `ContenedorTabla`, que se exporta. Con la tabla oculta usa solo `sr-only`, sin `overflow-x-auto`, y ya no queda una región desplazable. Con la tabla visible («Ver como tabla») usa `role="region"`, `aria-label` = título de la gráfica, `tabIndex={0}`, `overflow-x-auto` y el anillo de foco del token (`focus-visible:outline-border-focus`).
- Lo usan `FiguraGrafica` (G-1, G-6, G-4, G-5 y las T-1 a T-4 de Tendencias) y `FiguraBarras` de `Tablero.tsx` (G-2 y G-3). Las dos dejaron el envoltorio `cn('overflow-x-auto', !comoTabla && 'sr-only')`.

**Pruebas nuevas** en `frontend/pruebas/inicio/tablero.escritorio.spec.ts` (corren en 412×915 táctil y en 1440×900):
- «servidor anterior a C1 v1.5»: con 501 y luego con 404 en las dos operaciones nuevas, B-1, G-1 y G-6 tienen `role="status"` = «Disponible cuando el servidor se actualice», 0 `role="alert"`, 0 «Reintentar», sin «Revisa que estés en la red de casa» y sin `$` ni `n %`. B-1 conserva «Ver créditos». K-1 y K-3 siguen con las cifras de la semilla ($1,197.64; «Pagos de octubre: $973.75 · 1 pago»).
- «axe WCAG 2.0–2.2 A/AA» en claro y en oscuro, es decir 4 combinaciones con los dos anchos. Usa las 6 gráficas pintadas con las respuestas SINTÉTICAS de A-19 a A-21, que pasaron a un ayudante `sinteticas()` con los mismos valores. Pide 0 violaciones serias o críticas con las tablas ocultas y otra vez con G-2 y G-1 en «Ver como tabla». Comprueba que la tabla oculta no es región y que la visible es región con nombre, `tabindex=0` y que recibe el foco.

## Archivos

- `/home/deepsec/finanzas-hogar/frontend/src/pantallas/inicio/Tablero.tsx`: imports (`Info`, `ErrorApi`, `ContenedorTabla`); `esServidorViejo`, `ServidorViejo`, B1, `EstadosGrafica` (`soloNueva`), `FiguraBarras`; G1 y G6 con `soloNueva`.
- `/home/deepsec/finanzas-hogar/frontend/src/pantallas/inicio/GraficaInicio.tsx`: `ContenedorTabla` y su uso en `FiguraGrafica`.
- `/home/deepsec/finanzas-hogar/frontend/pruebas/inicio/tablero.escritorio.spec.ts`: import de `AxeBuilder`, ayudante `sinteticas()` (es solo una extracción: A-19 a A-21 usan los mismos datos), las 2 pruebas nuevas (3 por proyecto).
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/correccion-tablero-frontend-r1.md` (este reporte).
- `textos.ts` no cambió: `servidorViejo` ya estaba en la línea 40.

## Verificación

Los temporales quedaron en `datos/tmp/correccion-tablero-r1/` y se borraron al terminar. El build se hizo en un directorio temporal y se sirvió con `vite preview` en el puerto 4185, sin tocar `frontend/dist/`, que es lo que sirve el binario de desarrollo.

| Criterio | Comando | Resultado |
|---|---|---|
| tipos | `cd frontend && npx tsc --noEmit` (incluye `pruebas/`) | 0 errores |
| build | `npx vite build --outDir ../datos/tmp/correccion-tablero-r1/dist --emptyOutDir` | listo; 80 entradas precacheadas |
| U3-03 / U3 2.0 §2 (H-1) | `npx playwright test -c <config temporal> pruebas/inicio` | **32 de 32** (antes 26; +3 por proyecto) |
| U5-06 / U5-03 (H-2) | la misma corrida; prueba axe en las 4 combinaciones (412 y 1440 × claro y oscuro) | 0 serios o críticos con las tablas ocultas y con «Ver como tabla» |
| Las pruebas nuevas detectan el defecto | las 3 pruebas nuevas contra el `frontend/dist/` anterior (solo lectura, puerto 4186) | **3 fallan**, como se esperaba: axe da `scrollable-region-focusable` en `grafica-g3, g1, g6, g2, g5, g4 > .overflow-x-auto.sr-only` (6 nodos, igual que la revisión) y B-1 no muestra «Disponible…» |
| Suite completa | `npx playwright test -c <config temporal>` (todo `frontend/pruebas`) | **259 passed, 4 skipped, 0 failed** (5.8 min; la revisión tenía 253: +6 nuevas) |
| C3-03 | `npm run verificar:rutas` | 212 rutas invocadas; 0 fuera de `openapi.yaml`; 0 sin registrar en `routes.go` (221) |
| U5-08 estático | `node scripts/hosts-externos.mjs <dist temporal>` | hosts externos que se pedirían = 0 |
| Tokens (DESIGN, C2) | grep de hex, `rgb(`, `oklch(`, `dark:` y colores en `style` en `Tablero.tsx` y `GraficaInicio.tsx` | 0 y 0 |
| U5 sin servidor | `make verify-u5` | EXIT 0 (U5-01 164/164, U5-03 serios = 0, U5-04 bajo umbral = 0, U5-08 0, U5-11 0) |
| Global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | **EXIT 0** («Verificación ampliada OK») |

## Supuestos

- **HECHO**: en B-1, G-1 y G-6, 501 y 404 cuentan como «servidor anterior a C1 v1.5». En el contrato estas dos operaciones no declaran 404, y el servidor responde 404 `not_found` «La ruta no existe» cuando la ruta no existe y 501 `not_implemented` cuando la ruta está registrada sin manejador.
- **PROPUESTA (sin decisión nueva)**: B-1 conserva el enlace «Ver créditos» en ese estado. U3 §2 solo pide el texto y «nunca una cifra inventada»; el enlace lleva a P-17, que ya existe, y es el mismo que tiene la tarjeta en carga y en falla. Si Axel lo quiere fuera, basta con quitar una línea de B1.
- La región de la tabla visible lleva como nombre el título de la gráfica, que es lo que pide la revisión. La tabla repite ese título en su `caption`. Lo dejé así, como pide la corrección sugerida.

## Bloqueos

Ninguno nuevo.

## Pendientes

- `frontend/dist/`, que sirve el binario de desarrollo en https://finanzas.local:8443, sigue con el build anterior. Lo regenera el integrador con `npm run build` o con `make verify-wf3-frontend` (`npm run verificar` reconstruye `dist/` y corre la suite en el 4173). Yo no lo toqué para no cambiar lo que sirve el servidor de desarrollo a media revisión.
- La próxima captura de P-01 (`make capturas-u5` o el objetivo equivalente) debe reescribir `docs/design/capturas-u5/axe/P-01-*.json`. Con esta corrección deberían seguir en 0 serios.
- Los hallazgos bajos H-3 a H-7 de la revisión quedan fuera de este encargo.

## Propuestas para Axel

- Ninguna que necesite decisión. La de «Ver créditos» en B-1 con servidor viejo está en Supuestos.

## Paquetes instalados

Ninguno. `@axe-core/playwright` 4.13.0 ya estaba en `frontend/package.json`.

## Consumo aproximado

41 acciones · ~95 KB leídos · ~18 KB escritos.
