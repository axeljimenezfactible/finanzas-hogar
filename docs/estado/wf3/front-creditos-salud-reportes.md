# WF-3a · front-creditos-salud-reportes (12-oct-2026)

## Hecho
- Las 8 pantallas del área en `frontend/src/pantallas/creditos-salud-reportes/`, registradas solas con `indice.ts`, que exporta P-17, P-18, P-19, P-20, P-21, P-23, P-24 y P-25:
  - **P-17**: totales (N-191), aviso A/B (N-193) y tarjetas de línea (N-190) con «Falta capturar · Capturar» (K18). Detalle en galería y tabla. Tendencias con la deuda por línea y los intereses (N-194). Calendario con la agenda de hoy a 60 días, las promociones y el calendario mensual en tableta o más (N-192).
  - **P-18**: variante según el tipo de línea:
    - revolvente: lecturas A/B con la activa destacada, mínimos y tasas por plazo (N-203, N-204);
    - tarjeta: estado del ciclo (N-205);
    - préstamo: amortización (N-206);
    - tienda: abonos (N-207).
    Además: compras (N-202), datos de la línea con los 5 datos de N-283 o «Falta capturar», serie A/B y la agenda de la línea.
  - **P-19**: alta con vista previa del servidor en cada cambio. Muestra «Calculando…» mientras calcula. Si falla, dice «No se pudo calcular el calendario» y no guarda. Se puede prellenar desde P-20. El detalle trae partes por persona, pagos con estado por parte, editar, eliminar y ver el gasto vinculado.
  - **P-20**: datos, restricciones (solo para esta comparación), pesos, recomendación con puntaje, los 4 óptimos lado a lado, válidos y descartados con su motivo. Avisa «Faltan restricciones · Capturar» y «Ningún escenario cumple…» con la restricción que más descarta.
  - **P-21**: fichas de Axel y Ximena con IMC, «Confirmar» la meta sugerida (K23), ajuste por objetivo en el lugar con «Restablecer» (K36, N-327) y kcal de hoy. Historial de peso editable («Antes → Después») y con eliminar. Tendencias de peso y de kcal contra la meta.
  - **P-23**: próximos envíos, reportes generados (tabla y calendario, descargar, «Reintentar envío») y avisos enviados. Generar el mensual con `?hoja=generar`.
  - **P-24**: mes por defecto el último cerrado, «parcial» marcado y las secciones del PDF. «Descargar PDF», «Compartir» o «Generar PDF» si no existe.
  - **P-25**: periodo (mes o rango), compras, despensas y partes de pagos, pagado y saldo, nota fija, «Generar PDF», «Compartir» y estados anteriores.
- La hoja «Marcar pago» (G-06, N-210) marca cada parte por separado, con fecha y monto por omisión y «Deshacer». Se usa en P-17, P-18 y P-19.
- Estados de U1 §2.2/§2.4 con `ConDatos`. Sin conexión, P-19 y P-20 muestran «Requiere conexión» y las escrituras de las demás pantallas quedan con `deshabilitadoPor`. Textos es-MX en `textos.ts`; dinero y fechas con `formatos.ts`; solo tokens C2.

## Archivos
- `frontend/src/pantallas/creditos-salud-reportes/`: `indice.ts`, `textos.ts`, `comun.tsx`, `P17Lineas.tsx`, `P18Credito.tsx`, `P19Compra.tsx`, `P20Comparador.tsx`, `P21Salud.tsx`, `P23Centro.tsx`, `P24Mensual.tsx`, `P25Karen.tsx`
- `frontend/pruebas/creditos-salud-reportes/`: `datos.ts` (API del área con la semilla canónica de U1 §3.21), `pantallas.escritorio.spec.ts`, `flujos.spec.ts`
- `docs/estado/bloqueos/WF3-creditos-salud-reportes.md`

## Verificación
- `npx tsc --noEmit`: 0 errores en todo el frontend.
- `npm run build`: correcto, con la app completa y el service worker.
- **C3-03**: `npm run verificar:rutas` da 201 rutas invocadas, 0 fuera de openapi.yaml y 0 sin registrar en routes.go (216 registradas). El área usa 39 rutas, todas de creditos-a/b, salud, reportes y comidas (`/meal-plan/trends`), más `/persons`, `/categories` y `/stores`.
- **Playwright**: `npx playwright test pruebas/creditos-salud-reportes` da **20/20**. Son 8 pantallas en 412×915 táctil y en 1440×900, más 4 flujos:
  - **F-4**: P-18 → «Calendario» → la 1.ª fila dice «Vie 30 oct 2026 · $973.75 · Hogar $598.82 · Karen $374.93 · pendiente». Al marcar la parte del hogar se envía un POST con Idempotency-Key ≥ 16; «Deshacer» manda el DELETE.
  - **F-5**: P-23 → P-24 sep-2026 → descarga `reporte-mensual-2026-09.pdf`. En oct-2026, «Generar PDF» y luego «Descargar».
  - **F-5 Karen**: P-23 → P-25 → «Generar PDF» → «Compartir».
  - **F-6 desde Salud**: 1 registro y el historial se recarga.

  Cada prueba de pantalla comprueba que no hay desplazamiento horizontal y, en móvil, que la acción principal mide ≥ 44 px.
- **U5-07 y C2**: grep de hex/rgb/hsl/oklch en mi carpeta = 0. Grep de palabras (Save, Cancel, Settings, Dashboard, Loading, Error, Home, Delete, Search, Add) = 0.
- **U5-06**: solo clases de C2, foco visible en lo interactivo y `coarse:` 44 px en los componentes de la base. **U5-09 y U5-13** no pasan por mis pantallas: P-22 y la cola son de la base.
- `make verify-all GOOSE=…/gobin/goose`: salida 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)».

## Supuestos
- La carpeta del área es `creditos-salud-reportes` (como dice el ALCANCE). CONVENCIONES §3 sugería `creditos`, `salud` y `reportes`; el registro funciona igual.
- Las pruebas están en `pruebas/creditos-salud-reportes/` (CONVENCIONES §8: el `testDir` es `pruebas/`) y usan ApiFalsa con `context.route`, como dice §8. Los mocks msw genéricos siguen sirviendo para `dev:mocks`.
- Lo demás está en el bloqueo, del punto 1 al 9.

## Bloqueos
`docs/estado/bloqueos/WF3-creditos-salud-reportes.md`

## Pendientes
- Puntos 2 a 6 del bloqueo: dependen de cambios en C1.
- Revisión visual en las 4 combinaciones (U5-06): la hace el revisor.

## Propuestas para Axel
Decidir los puntos 1, 3, 7 y 8 del bloqueo. Los valores provisionales están en la tabla.

## Paquetes instalados
Ninguno.

## Consumo aproximado
~58 acciones, ~160 KB leídos y ~150 KB escritos (en el tope).
