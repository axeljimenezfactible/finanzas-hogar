# correccion-frontend-core-r4 · WF-3a · 12-oct-2026

Agente: `correccion-frontend-core-r4` (Opus). Corrige los dos hallazgos medios de `REVISION-frontend-core-r4.md`: **R4-01** y **R4-02**. Los bajos (R4-03, R4-04, R4-05, H-07, H-08 y H-09) no estaban en el encargo y no se tocaron. No se reabrió ninguna decisión ni se tocó una cifra. Los motores, el backend, Pista 3 y los archivos de diseño quedan igual.

## Hecho

### R4-01 · P-03: pestaña «Ingreso» (N-102) y forma de pago «A crédito» (N-053), U1 §3.5
- `src/shell/HojaRegistro.tsx`: la pestaña «Ingreso» monta `<FormIngreso refMonto={ref} alTerminar={cerrar} />` en lugar del desvío «Ir a Ingresos». El foco va al Monto al abrir la hoja. Guarda con `POST /api/v1/incomes` y clave de idempotencia. Sin conexión, la hoja muestra «Requiere conexión».
- Con la forma de pago `credit`, la hoja monta `<LineaCredito monto dia comercio={descripcion} alSalir />` en lugar de la franja con «Ir a compra a crédito». La pieza lista las líneas activas en orden alfabético con su nombre y su tipo (`GET /api/v1/credits`) y sigue en P-19 con el monto, la fecha y el comercio ya llenos. Si no hay líneas, muestra el estado vacío con «Ir a Ajustes › Créditos».
- Las dos piezas ya existían en `pantallas/gastos-presupuesto/hoja-registro.tsx` y no se cambiaron. Esto cierra la «solicitud 1» de `front-gastos-presupuesto-fix.md`.
- `pruebas/gastos-presupuesto/correccion-bis.escritorio.spec.ts:166` y `:179`: los `test.skip(…isVisible())` condicionales ahora son aserciones: el desvío viejo debe tener 0 elementos y el Monto debe estar visible. Si el desvío vuelve, la prueba falla; ya no se salta.
- `pruebas/medidas.spec.ts:95`: esa línea esperaba el enlace «Ir a compra a crédito» que R4-01 quita, así que ahora espera el estado vacío de N-053. Sin este cambio, `make verify-wf3-frontend` fallaba. Además, la medición U5-06 de la hoja ya cubre el contenido nuevo de N-053 y N-102.

### R4-02 · Gráficas: eje compacto, `tabular-nums` en el SVG y etiquetas de dinero (DESIGN §7.7 y §8.5, U5-06)
1. **Eje compacto.** `src/api/formatos.ts:21` tiene ahora `dineroCompacto`, que es el `compacto` de Inicio: aritmética entera con ROUND_HALF_UP, signo con U+2212, por ejemplo `$0`, `$5 mil`, `$12.3 mil`. `Grafica` tiene una prop opcional nueva, `formatoEje`. Por omisión, una gráfica de dinero (la que no pasa `formato`) usa `dineroCompacto` en `yAxis.axisLabel`. Las de kilos y kcal usan su propio `formato`, como antes.
2. **Cifras tabulares en el SVG.** Se agregó `tabular-nums` al contenedor de ECharts en `Grafica.tsx:121` y en `GraficaInicio.tsx:64`. Los `<text>` y la descripción emergente lo heredan. No hizo falta una regla en `app.css`.
3. **Etiquetas directas con el formato completo, selectivas (§8.5):**
   - en líneas, el último valor y los extremos, con `selectivos`, que ahora vive en `Grafica.tsx:41`;
   - en área, el valor final;
   - en columnas, el valor sobre cada columna solo con 4 series o menos (`MAX_SERIES_CON_ETIQUETA`);
   - `labelLayout.hideOverlap`: una etiqueta que no cabe no se recorta; queda en «Ver como tabla»;
   - color `--ui-text-default`, nunca el de la serie (§8.1).
4. **Inicio reutiliza las mismas piezas.** `GraficaInicio.tsx:16` y `:18` reexportan `compacto` (= `dineroCompacto`) y `selectivos`, así que `Tendencias.tsx` no cambió.
5. **Medición.** `pruebas/apoyo/medidas.ts:77` (`dineroSinTabulares`) ahora incluye los `<text>` del SVG de ECharts, aunque su contenedor sea `aria-hidden` (el de Inicio lo es). Los marca como «gráfica: …». Los demás SVG (iconos) y lo `aria-hidden` fuera de las gráficas se siguen omitiendo.

## Archivos
- `frontend/src/shell/HojaRegistro.tsx` (frontend-base)
- `frontend/src/componentes/Grafica.tsx` (frontend-base)
- `frontend/src/api/formatos.ts` (frontend-base)
- `frontend/src/pantallas/inicio/GraficaInicio.tsx` (front-inicio)
- `frontend/pruebas/apoyo/medidas.ts` (frontend-base)
- `frontend/pruebas/medidas.spec.ts` (frontend-base, solo la línea 95: consecuencia directa de R4-01)
- `frontend/pruebas/gastos-presupuesto/correccion-bis.escritorio.spec.ts` (front-gastos-presupuesto, solo `:166` y `:179`)
- `docs/estado/wf3/correccion-frontend-core-r4.md` (este reporte)

## Verificación
- **Build y tipos.** `npm run build` termina limpio (preparar, `tsc --noEmit` y vite). Los avisos de las fuentes Inter `woff2` ya estaban antes.
- **R4-01 (U1 §3.5).** Corren **y pasan** en las dos combinaciones (movil 412 y escritorio 1440):
  - `P-03 N-102 · la pestaña «Ingreso» registra el ingreso en la misma hoja`: POST `{kind:'extra', description:'Reembolso', amount_cents:120000}` y la hoja se cierra;
  - `P-03 N-053 · forma de pago «Crédito» elige la línea activa…`: opciones `['Elige la línea', 'Línea de la tienda · Crédito de tienda', 'Tarjeta azul · Tarjeta bancaria']` y URL `/creditos/2/compras/nueva?monto=185.50&fecha=…`.
  - Antes eran 4 pruebas saltadas.
- **R4-02 (DESIGN §7.7/§8.5, U5-06).** Prueba propia temporal en `datos/tmp/correccion-fcore-r4/` (ya borrada), 412 y 1440, con los datos falsos del área, medida sobre los `<text>` del SVG:
  - P-02 Tendencias: eje `$0 | $1 mil | $2 mil | $3 mil | $4 mil`; etiquetas `$4,000.00 | $1,500.00 | $0.00 | $1,200.00`; sin tabulares 0.
  - P-06 Tendencias: eje `$0 | $5 mil … $25 mil`; etiquetas `$22,000.00 | $22,000.00 | $21,000.00 | $18,500.00`; sin tabulares 0.
  - P-07 Tendencias: eje `$0 | $3 mil … $18 mil`; etiquetas `$18,000.00 | $17,000.00 | $1,000.00`; sin tabulares 0.
  - P-08 Tendencias: eje `$0 | $5 mil … $30 mil`; etiquetas `$30,000.00 | $2,500.00`; sin tabulares 0.
  - P-09 Tendencias: eje `$0 | $10 mil … $50 mil`; etiquetas `$45,000.00 | $1,500.00 | $3,000.00`; sin tabulares 0.
  - P-01 Tendencias (GraficaInicio): 15 cifras en el SVG, sin tabulares **0** (antes 15).
  - En todas: ningún texto de dinero fuera de los dos formatos y `dineroSinTabulares` = `[]`.
- **U5-06.** `npx playwright test pruebas/medidas.spec.ts`, ahora con el SVG medido: «dinero sin tabulares 0» en Inicio (4 vistas), Gastos y presupuesto (15) y Comidas y compras (10). Estados abiertos: 11 estados, 299 controles, 0 menores de 44 px, 0 iconos fuera de norma, dinero sin tabulares 0. Incluye «F-1 Gasto · Más opciones» con N-053 y «F-1 Ingreso» con N-102.
- **N-04 y H-03.** La prueba del banco de la gráfica pasa en claro y oscuro: `aria-label` en español, descripción emergente con tokens, series sin recolorear.
- **`make verify-wf3-frontend`: exit 0.**
  - «C3-03 frontend: 208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)».
  - «U5-08 dist: hosts externos que se pedirían = 0».
  - e2e **231 passed, 4 skipped** (5.0 min), antes 223 y 8. Los 4 saltados son los legítimos por dispositivo: 3 atajos `Ctrl+/` en movil y 1 medida táctil en escritorio. Ninguno es de P-03.
  - La primera corrida dio 1 fallo en `medidas.spec.ts:95` por el enlace retirado; se corrigió como se explica arriba.
- **`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: exit 0.** Termina con «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)». `git status` es igual antes y después: los archivos generados no cambiaron.

## Supuestos
- **SUPUESTO 1.** Las etiquetas directas selectivas se aplican a todas las gráficas de `Grafica`, también las de kilos y kcal (P-10, P-21), con su `formato` completo. §7.7 pide etiquetas directas para cualquier gráfica y §8.5 pide no poner un número en cada punto. El eje compacto, en cambio, es solo para dinero.
- **SUPUESTO 2.** El tope de 4 series para etiquetar columnas sale de la corrección sugerida en R4-02. §8.5 dice «si cabe con aire». `hideOverlap` cumple el «si no cabe, en la tabla».
- **SUPUESTO 3.** `dineroCompacto` ya no escribe «−$0» para un negativo menor de 50 centavos, como hacía `compacto`; ahora escribe `$0`. En lo demás es idéntico.
- **SUPUESTO 4.** Toqué `medidas.spec.ts:95`, que no se nombra literalmente en R4-01. Esa línea afirmaba justo lo que R4-01 retira, y el archivo es de frontend-base, dueño del hallazgo.

## Bloqueos
Ninguno nuevo.

## Pendientes
- **Integrador.** Regenerar la matriz C3 (`npm run verificar:c3` con `--escribir`, según `front-gastos-presupuesto-fix.md`, «solicitud 3») para que P-03 cuente `GET /api/v1/credits` y `POST /api/v1/incomes`. También actualizar `docs/contratos/C3-integracion.md` y `docs/estado/bloqueos/C3-matriz.md` si corresponde.
- **Integrador.** El servidor de `https://finanzas.local:8443` sirve el frontend compilado dentro del binario. Para ver estos cambios en la LAN hay que recompilarlo o relanzar `ops/dev-run.sh`. No lo hice porque está fuera de mi alcance.
- **Bajos de la ronda 4, no encargados:**
  - R4-03, `GraficaInicio.tsx:49`: `aria.label` en inglés. Sigue impreso como «PENDIENTE U5-07» en `medidas.spec.ts`.
  - R4-04, `Barras.tsx:52`: alto de 55 px.
  - R4-05, la excepción `PENDIENTES_INGLES`.
  - H-07, H-08 y H-09.

## Propuestas para Axel
- **PROPUESTA.** Con «A crédito» elegido, el botón «Guardar» de la hoja no hace nada, sin aviso, porque la compra se termina en P-19 con «Continuar en Compra a crédito». Eso ya pasaba antes y no estaba en ningún hallazgo. Propongo deshabilitarlo con el motivo visible («La compra a crédito se termina en Créditos») o esconderlo mientras esa forma esté elegida. Es decisión de diseño: no lo cambié.
- **PROPUESTA.** Agregar a la prueba del banco de la gráfica (`pruebas/base.escritorio.spec.ts`) dos aserciones permanentes: el eje compacto (`$4 mil`) y la etiqueta completa (`$4,000.00`). Así R4-02 queda cubierto en la suite y no solo en mi medición temporal. Ese archivo no lo nombraba el hallazgo.
- **PROPUESTA.** En `src/textos/es-MX.ts`, `registro.irIngresos`, `registro.irCredito` y `registro.creditoAviso` ya no se usan en la aplicación. `registro.ingresoIr` solo sigue en las historias de Storybook. Se pueden quitar en una limpieza de frontend-base.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 32 acciones, unos 110 KB leídos y unos 25 KB escritos (código, pruebas y este reporte).
