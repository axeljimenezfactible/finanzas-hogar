# WF-3a · front-ajustes-b · Ajustes, secciones 8 a 14 (P-26.08 a P-26.14)

Agente `front-ajustes-b` (Opus), 12-oct-2026. Fuentes: U4 parte b (y §0 de la parte a), U1 §2.3, CONVENCIONES-frontend.md, C1 v1.4 (`esquema.gen.ts`), criterios U5, S6, S7 y C3.

## Hecho
- Las 7 secciones como componentes que el marco de **ajustes-a** (`P26Ajustes.tsx`: columna de 14 secciones, «‹ Ajustes» y h2) recoge de `src/pantallas/ajustes-b/secciones-ajustes.ts` (convención de ajustes-a; `indice.ts` no registra pantallas, así nada tapa `P-26`).
- **P-26.08**: 3 reportes (tabla en 1440, tarjetas + hoja en 412), «Activo», día y hora por clave `reports.*` con «Valor inicial», «Restablecer al valor inicial» y vista previa «Próximo envío» (Q2-66 B con ejemplo de mes corto), destinatarios solo con chat confirmado, «Sin destinatarios», último resultado, «Enviar ahora a…», bajo demanda (mensual sin preselección de Telegram, estado de cuenta de Karen, vista web).
- **P-26.09**: 4 alertas con sus claves `alerts.*` y vistas previas (umbrales en centavos HALF_UP sobre el presupuesto del mes, fechas del próximo pago, días de despensa, día de revisión), hora de avisos, chats confirmados (etiqueta, «Recibe avisos»), chats detectados (terminación •••nnnn, confirmar con Idempotency-Key, 409 en palabras, descartar con acción crítica), aviso de prueba.
- **P-26.10**: tema sistema/claro/oscuro (local, sin red, al instante) e interruptor «Atajos de una tecla» (local).
- **P-26.11**: respaldos (encabezado último/próximo/restauración probada, sondeo mientras hay uno «En curso», detalle con restauración en doble confirmación y `RESTAURAR`), hora del respaldo diario (PROVISIONAL, `updateJob` cron), importación CSV con vista previa, semillas, IA opcional con acción crítica al encender, estado del servicio y 3 parámetros K44.
- **P-26.12**: personas (nombre, papel, cuenta, paga directo, qué hace, tope por quincena con vista previa, notas), repartos reutilizables con «el resto», V-16, vista previa con $M editable y diálogo de alcance (`scope` onward/all).
- **P-26.13**: perfiles con IMC y meta sugerida RN-17 con desglose (oráculo: 1,648.75 × 1.55 = 2,555.5625 → 2,556), ajustes por persona y de la casa, «Usar la sugerida», «Usar el de la casa»; metas en tabla/tablero/calendario/galería con vistas previas RN-18 (avance, aportación necesaria, estado, total reservado), «Nueva meta» (Alt+N).
- **P-26.14**: inventario tabla/tablero por modo, orden agotado → reponer → suficiente, vistas previas RN-15/K7/Q2-70 (a), «Agregar producto al inventario» (Alt+N).
- Estados de U4 §0.3: esqueleto, «No se pudo cargar {qué} · Reintentar», «Guardando…», «Guardado» (aria-live), error al guardar con Reintentar, «Valor restablecido · Deshacer», escrituras deshabilitadas sin conexión con la razón.

## Archivos
`frontend/src/pantallas/ajustes-b/`: `indice.ts`, `secciones-ajustes.ts`, `textos.ts`, `calculos.ts`, `comun.tsx`, `P2608Reportes.tsx`, `P2609Notificaciones.tsx`, `P2610Apariencia.tsx`, `P2611Datos.tsx`, `P2612Personas.tsx`, `P2613Salud.tsx`, `P2614Inventario.tsx` · `frontend/pruebas/ajustes-b/datos.ts`, `ajustes-b.escritorio.spec.ts` · `docs/estado/bloqueos/WF3-ajustes-registro.md`, `WF3-ajustes-muestra-tema.md` · este reporte.

## Verificación
- tsc: `npx tsc --noEmit` → 0 errores en `ajustes-a`/`ajustes-b`; el único error del árbol es ajeno: `pruebas/creditos-salud-reportes/datos.ts(48,79) TS6133`, que hoy rompe `npm run build`. Build aislado `npx vite build --outDir datos/tmp/ajustes-b/dist` → «✓ built».
- Playwright (copia de la configuración base con build y puerto propios, 4791): **14 de 14** (7 pantallas × 412×915 táctil y 1440×900).
- C3-03: `npm run verificar:rutas` → «201 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0».
- U5-06: hex literales en mi carpeta y pruebas = 0; sin `dark:` ni colores en línea; solo clases de C2; controles `coarse:h-11`/`min-h-11` (44 px), `focus-visible` en todo control propio; dinero con `tabular-nums`; estados con icono + texto. Única medida en línea: el ancho % de barras de datos (reparto, avance).
- U5-07: texto en inglés visible (Save, Cancel, Settings, Dashboard, Loading, Error, Home, Delete, Search, Add) = 0 (la única coincidencia de `grep -w` es el identificador `AlertDialog.Cancel` de Radix); fechas con `fecha`/`fechaConDia`/`fechaHora` (la prueba verifica «Lun 12 oct 2026, 21:30» y «12-oct-2026 03:00»).
- U5-09 y U5-13: no aplican (estas secciones no capturan gastos ni peso). Flujos F-1..F-6: ninguno pasa por P-26.08..14 (F-5 va por P-23/P-24).
- S6-09: tabla con integridad y restauración probada; restaurar pide doble confirmación y `RESTAURAR`. S6-10: `attachment.max_size_mb` (inicial 10 MB) editable y restablecible (prueba P-26.11). S6-13: IA apagada por defecto, encender pide acción crítica (prueba).
- S7: la app no envía Telegram; solo llama `sendTestNotification`, `generateMonthlyReport` del API (S7-07); nunca se muestra el `chat_id` completo ni el token.
- `make verify-all GOOSE=…` → «Verificación ampliada OK».

## Supuestos
- Pruebas en `frontend/pruebas/ajustes-b/` (CONVENCIONES §1 y §8: `testDir: 'pruebas'`).
- Ajustes se escriben con `PATCH` (C1), aunque U4 nombra `PUT` en P-26.08; manda el contrato.
- Destinatario válido = destinatario con `verified_at`; Karen nunca aparece (decisión A).
- Estado de cuenta de Karen: periodo por defecto el mes en curso (U1 F-5).

## Bloqueos (PROVISIONAL)
- `WF3-ajustes-registro` (resuelto entre áreas con `secciones-ajustes.ts`; propuesta A para la base).
- `WF3-ajustes-muestra-tema` (muestra claro y oscuro a la vez imposible con tokens en `:root`).
- Se citan los abiertos: `C1-hora-reportes-job`, `U4-reparto-preseleccion`, `C1-reparto-alcance`, `S2-metas-vencida`, `WF3-frontend-k31` (hora del respaldo).

## Pendientes
- `Ctrl+/` (edición múltiple con vista previa sumada) en alertas, metas e inventario: `VistaBD` no lo trae.
- Arrastre en el tablero de inventario con vista previa y «Aplicar»: `VistaBD` guarda al soltar (`moverA`).
- P-26.09: vista previa de umbrales «por categoría al tocar una categoría de ejemplo».
- Prueba axe de estas pantallas.

## Propuestas para Axel
- `WF3-ajustes-muestra-tema`: opción A. `U4-reparto-preseleccion`: queda «Solo a los gastos nuevos». Marcador $M del reparto: $100.00 al abrir (PROVISIONAL, editable).

## Paquetes instalados
Ninguno.

## Consumo aproximado
~59 acciones · ~200 KB leídos · ~175 KB escritos (**pasa el tope de 150 KB** del protocolo §4 por ~25 KB: 7 secciones con sus textos y pruebas).
