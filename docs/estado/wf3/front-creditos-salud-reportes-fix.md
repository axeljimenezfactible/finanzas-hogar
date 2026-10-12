# front-creditos-salud-reportes-fix · WF-3a bis

Fecha: 12-oct-2026 · Agente: `front-creditos-salud-reportes-fix` (Opus) · Alcance: `frontend/src/pantallas/creditos-salud-reportes/` y `frontend/pruebas/creditos-salud-reportes/`.
Encargo: bloqueo `docs/estado/bloqueos/C3-matriz.md` punto 2, decisión del orquestador «implementar las necesidades sin llamada en sus pantallas»: **P-17 N-193** y **P-18 N-208**; «Ver como tabla» ≥ 44 px táctil (H-03 parte 2) si la base no lo cubre. Nada de esto cambia cifras ni motores.

## Hecho

1. **P-17 · N-193 (aviso A/B) con `getCredit`.** `P17Lineas.tsx`: `FranjaAB` ya no pinta el aviso con lo que trae `listCredits`. Por cada línea revolvente lee `GET /api/v1/credits/{credit_id}` (C1 `getCredit`, `x-necesidades` N-193) y pinta la lectura activa (A o B) y `availability_confirmed` de esa respuesta: «Lectura B activa · Kueski Pay · Reposición de capital pendiente de confirmar con Kueski» con el enlace «Confirmar en Ajustes › Créditos y condiciones». Mientras llega la respuesta, o si falla sin copia, usa lo de `listCredits` para no esconder el aviso. P-18 ya leía `getCredit`. Ahora usa la pieza pura `FranjaLinea` con esos datos.
2. **P-18 · N-208 (calendario de la línea) con `listCalendarEvents`.** `P18Credito.tsx`, pestaña «Calendario»: `GET /api/v1/calendar/events` con `credit_id=<línea>`, `event_type=credit_payment,promotion,card_statement` (los tipos que C1 da a N-192/N-208) y `period=<día 1 del mes en curso, o el pago sin marcar más viejo>..<hoy + 60>` (sintaxis de rango de C1 `Period`). Se pinta así:
   - **Agenda (todas las anchuras):** pagos y cortes de hoy a 60 días, más los pagos sin marcar con fecha pasada (igual que antes con `enAgenda`), en el orden de C1. Fila de pago: «Vie 30 oct 2026 · Pago de Kueski Pay · $973.75» y «Hogar $598.82 · Karen $374.93 · pendiente». El estado de cada parte sale del pago programado ya cargado (`ref_id`). Si ese pago no está cargado, salen las partes y el estado del evento. Al tocar un pago se abre la hoja «Marcar pago» (N-210) con ese pago. La fila de corte dice «Saldo al corte».
   - **Promociones vigentes** como franjas: «Promoción: Walmart a 4 pagos · Vigente del … al …».
   - **De tableta en adelante, «Calendario de la línea»** (nombre de N-208 en U1): `VistaBD` con vistas calendario y tabla, igual que en P-17, porque N-192 pide el mes visible fuera del teléfono. Un pago abre «Marcar pago». Una promoción o un corte abre P-26.05 de esa línea (ver Supuestos).
   - Después de marcar o desmarcar un pago se recargan los eventos. Mientras cargan los pagos, que dan el inicio del rango, se muestra «Cargando» y no un error.
3. **«Ver como tabla» (H-03 parte 2): la base ya lo cubre** (`frontend-base-fix`: 36 px con puntero fino, 44 px táctil). No cambié nada. Agregué una prueba que lo mide en P-17 y P-18 (`?sub=tendencias`): ≥ 44 en `movil` y ≥ 36 en `escritorio`. Pasa.
4. **`Ctrl+/`:** U4 lo pide solo en secciones de Ajustes (P-26.02, .04, .05, .07, .09, .13). Ninguna pantalla de esta área (P-17 a P-25) lo tiene en U4, así que no aplica prueba de `Ctrl+/`.
5. Textos nuevos en `textos.ts`, en español de México: `EVENTO_LINEA` (Pago, Promoción, Corte de tarjeta), `calendarioLinea` («Calendario de la línea») y `vigencia`. `comun.tsx` solo exporta `agruparEstado`, que ya existía.

## Archivos

- `frontend/src/pantallas/creditos-salud-reportes/P17Lineas.tsx`: `FranjaLinea`, `AvisoAB` (getCredit) y `FranjaAB`.
- `frontend/src/pantallas/creditos-salud-reportes/P18Credito.tsx`: `periodoLinea`, `AgendaLinea`, `CalendarioLinea` y la consulta `listCalendarEvents`.
- `frontend/src/pantallas/creditos-salud-reportes/textos.ts` y `comun.tsx`: textos y una exportación.
- `frontend/pruebas/creditos-salud-reportes/datos.ts`: API falsa de `GET /api/v1/calendar/events`, que filtra por `period`, `credit_id` y `event_type` como el servidor. Exporta `KUESKI`.
- `frontend/pruebas/creditos-salud-reportes/pantallas.escritorio.spec.ts`: 3 pruebas nuevas (N-193, N-208 y «Ver como tabla»).
- `docs/estado/wf3/front-creditos-salud-reportes-fix.md` (este reporte).

## Verificación

| Criterio | Comando | Resultado |
|---|---|---|
| Tipos | `cd frontend && npx tsc --noEmit` | 0 errores en todo el frontend en la última corrida. En una corrida intermedia hubo errores de otros agentes en `ajustes-a/S05Creditos.tsx` y `pruebas/ajustes-b/datos.ts`, que ya se corrigieron solos |
| Build | `npx vite build --outDir ../datos/tmp/csr-fix/dist` | OK (PWA generada). Lo construí en una carpeta propia para no pisar `dist/` ni el preview de otros agentes. `npm run build` = preparar + `tsc` + `vite build`. `tsc` pasa |
| Pruebas del área | `npx playwright test --config <copia temporal con puerto 4183> creditos-salud-reportes` | **26 passed, 0 failed**: 13 en `movil` 412×915 táctil y 13 en `escritorio` 1440, entre ellas F-4, que marca la parte del hogar desde la agenda nueva |
| C3-02 / C3-06 (implementado) | `python3 frontend/scripts/c3-integracion.py`, sin `--escribir` | P-17 y P-18 **ya no salen en NO-IMPLEMENTADO** (antes: «P-17 N-193 (getCredit)» y «P-18 N-208 (listCalendarEvents)»). Total: implementado sí 39/42. Lo que queda es de otras áreas |
| C3-03 (rutas) | `npm run verificar:rutas` | fuera de openapi.yaml: 0; sin registrar en routes.go: 0 |
| U5-06 (44 px) | prueba N-208: `controlesChicos` en `movil`; `revisar()` en las nuevas | 0 controles bajo 44 × 44 en la pestaña Calendario de P-18. Enlace del aviso, fila de la agenda y «Ver como tabla» ≥ 44 en táctil |
| U5-07 (inglés y cifras) | `textosEnIngles` y `dineroSinTabulares` en las pruebas nuevas; grep de textos en inglés en `src/pantallas/creditos-salud-reportes/` | 0 y 0 (las dos con tema claro); grep = 0 |
| Colores (DESIGN, C2) | `grep -rnE "#[0-9a-fA-F]{3,8}\b"` en el área y sus pruebas | 0. Solo clases de token: `bg-surface-card`, `border-border-divider`, `text-text-secondary`, `bg-surface-hover-subtle` y `outline-border-focus` |
| U5-08 (hosts) | `node scripts/hosts-externos.mjs ../datos/tmp/csr-fix/dist` | hosts externos que se pedirían = 0 |
| Global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | salida 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)» |

## Supuestos (HECHO con marca PROVISIONAL, reversibles)

- **Inicio del rango de P-18:** día 1 del mes en curso, o el pago sin marcar más viejo, hasta hoy + 60. C1 da para la agenda móvil `hoy..hoy+60`. Lo extendí hacia atrás para conservar los pagos «sin marcar con fecha pasada», que U1 N-192 pide y que la agenda ya mostraba, y para llenar el mes visible de tableta. La agenda filtra a hoy en adelante más los vencidos.
- **Promoción o corte en el calendario de tableta:** al tocarlos se abre P-26.05 de esa línea, donde U4 ubica promociones y datos del estado de cuenta (N-283). Si se prefiere que no abran nada, es un cambio de una línea (`abrirEvento`).
- **Promociones en P-18:** salen del evento (título «comercio a n pagos» y vigencia). La compra mínima, que P-17 sí muestra con `listPromotions`, no viene en `CalendarEvent`. No agregué otra llamada.

## Bloqueos

Ninguno nuevo.

## Pendientes (fuera de mi alcance)

- Regenerar la matriz C3 con `python3 frontend/scripts/c3-integracion.py --escribir` (integrador) y marcar en `docs/estado/bloqueos/C3-matriz.md` punto 2 que P-17 N-193 y P-18 N-208 quedaron resueltos por (c).

## Propuestas para Axel

- PROPUESTA: si se quiere la compra mínima en las franjas de promoción de P-18, agregar `min_purchase_cents` al `CalendarEvent` de tipo `promotion` en una versión futura de C1, o leer `listPromotions` filtrado por línea. Hoy no hace falta para N-208.

## Paquetes instalados

Ninguno.

## Consumo aproximado

Unas 45 acciones; unos 125 KB leídos; unos 22 KB escritos (código, pruebas y este reporte). Los temporales de `datos/tmp/csr-fix/` (build propio, configuración de Playwright y salidas) se borraron al terminar.
