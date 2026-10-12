# front-gastos-presupuesto-fix · WF-3a bis (12-oct-2026)

Agente: front-gastos-presupuesto-fix (Opus). Alcance: `frontend/src/pantallas/gastos-presupuesto/`, `frontend/pruebas/gastos-presupuesto/` y este reporte. Sin `git add` ni `commit`. No cambia cifras, motores ni decisiones de Axel.

## Hecho
- **P-02 N-043**: tablero por categoría sobre `queryView` (`/api/v1/views/{entity}`, `expense`, `mode=board`, `group_by=category_id`): cada cabecera dice «Categoría · total (parte del hogar) (%)» y las columnas van por monto descendente.
- **P-02 N-045**: subsección Calendario sobre `queryView` (`mode=calendar`, `date_field=expense_date`): cada día empieza con «Total del día · $X · n gastos» y sigue la lista; los venideros llevan su marca.
- **P-02 N-046**: galería con la foto del ticket (`image_attachment_id` de `queryView` en galería y la descarga con `downloadAttachment` vía `armarUrl`, con carga diferida). En la tabla, la columna «Ticket» dice «Sí» o «No» y no descarga nada. Ver bloqueo §1 y §3.
- **P-02 `Ctrl+/`** (U1 §5.3): edición en grupo de categoría o tipo, con vista previa sumada («Gastos que cambian 0 → n» y «Parte del hogar que cambia $0.00 → $X»). `Esc` no guarda y `Enter` aplica.
- **P-04 N-064**: bloque «Ticket adjunto» con la foto (`downloadAttachment`, abre el archivo) y «Quitar adjunto» con confirmación y `deleteAttachment`. El id del adjunto se obtiene como dice el bloqueo §2. El botón deshabilitado «Ver ticket aún no está disponible» se quitó.
- **P-05 N-072**: cada renglón tiene «Producto del catálogo · Buscar en catálogo». La búsqueda usa `listProducts` con `q` y espera 300 ms; se elige el producto y su presentación (`getProduct`) o «Usar solo el producto». La confirmación manda `presentation_id`.
- **P-05 N-073**: con `?semana=2026-W42` (lo que manda P-13 «Terminar compra»), `getShoppingList` (`is_bought=true`) propone «Agregar n renglones» con nombre, presentación, piezas y tienda (bloqueo §4).
- **P-03 N-053 y N-102**: las piezas `LineaCredito` (`listCredits`: líneas activas en orden alfabético, nombre · tipo, y P-19 prellenada con `monto` en pesos, `fecha` y `comercio`) y `FormIngreso` (`createIncome` extra con clave de idempotencia; sin conexión muestra «Requiere conexión») están en `hoja-registro.tsx`. Montarlas en `src/shell/HojaRegistro.tsx` está fuera de mi alcance: es la solicitud 1, probada en una copia temporal.
- **P-08**: `listIncomes` (`kind=fixed`) da el registro de cada quincena capturada para corregir en el lugar su fecha, monto («Antes → Después») y nota (N-105). «Registrar ingreso extra» usa `FormIngreso`, la misma pieza de P-03. La pantalla escucha `fh:datos-cambiaron`. `listPeriods` queda como excepción técnica PROVISIONAL (bloqueo §6).
- **R3-01 y R3-02**: no están en mi carpeta (`pruebas/ajustes-a/validar-c1.ts`, `src/vistas-bd/CeldaEditable.tsx`), así que no hice nada.

## Archivos
- `frontend/src/pantallas/gastos-presupuesto/P02Gastos.tsx`, `P04Gasto.tsx`, `P05Ticket.tsx`, `P08Ingresos.tsx`, `textos.ts` (modificados)
- `frontend/src/pantallas/gastos-presupuesto/hoja-registro.tsx` (nuevo; no importa `./comun`, para no sumar sus lecturas a P-03 en C3)
- `frontend/pruebas/gastos-presupuesto/correccion-bis.escritorio.spec.ts` (nuevo, 11 pruebas × 2 proyectos)
- `docs/estado/bloqueos/WF3a-bis-gastos-presupuesto.md` (nuevo) y este reporte

## Verificación
- `npx tsc --noEmit`: 0 errores. `npm run build`: código de salida 0.
- C3-03 `npm run verificar:rutas`: 208 rutas invocadas, 0 fuera de openapi.yaml y 0 sin registrar en routes.go.
- Pruebas nuevas, `npx playwright test pruebas/gastos-presupuesto/correccion-bis.escritorio.spec.ts`: **17 pasan y 5 se saltan**. Las que se saltan son `Ctrl+/` en móvil (atajo de teclado) y las 4 de P-03, que se saltan solas hasta que se monte la solicitud 1.
- Área, shell y cola (`pruebas/gastos-presupuesto/ pruebas/shell.spec.ts pruebas/cola.spec.ts`): **41 pasan, 5 se saltan, 0 fallan**. Incluye F-1/U5-09 (≤ 3 toques).
- Copia temporal con la solicitud 1 aplicada (`datos/tmp/fe-p03`, ya borrada): build con salida 0; P-03 N-102, P-03 N-053, P-02 y F-1 (U5-09) dan **8 pasan de 8** en móvil y escritorio. Antes, en esa copia, toda la carpeta del área más shell y cola dio 41 pasan, 1 se salta y 4 fallan. Las 4 fallas eran de P-02: la API falsa responde sin `groups`. Se corrigió en P-02 (usa sumas locales) y luego todo pasó.
- `python3 -I scripts/c3-integracion.py` (sin `--escribir`): implementado **41/42**. Solo queda «NO-IMPLEMENTADO P-03 N-053 (listCredits), N-102 (createIncome)», que espera la solicitud 1. Operaciones no llamadas: 1 (`listPeriods`, bloqueo §6). `queryView`, `listIncomes`, `deleteAttachment`, `getShoppingList` y `listProducts` ya se llaman.
- Hex en mi carpeta y en mi prueba: 0. Inglés visible: 0. El grep de palabras inglesas da 15 coincidencias, todas identificadores de código: `rol: 'error'`, `pedir('delete', …)`, `Esquemas['Attachment']` y `Esquemas['Product']`. Sin `fetch`, `dark:` ni colores en línea.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: salida 0, con «WF-1: verificación completa OK» y «Verificación ampliada OK».
- Criterios: C3-02 y C3-03 como arriba. C3-04 no lo cuento: el script no recibió `--pw-json` (le toca al integrador). U5-07 (textos es-MX) y U5-08 (sin hosts externos, sin CDN) se cumplen en mis archivos. S6 y S7 no aplican a este alcance (solo frontend).

## Supuestos
Todos son PROVISIONALES y están en el bloqueo `WF3a-bis-gastos-presupuesto`, §1 a §7:
- Sumas locales cuando `queryView` no cubre el filtro.
- Id del adjunto por galería.
- Detección tabla/galería.
- Reglas de N-073.
- Toda línea activa va a P-19.
- `listPeriods` como excepción.
- «Total del día» abre el Detalle.
- PATCH en serie en `Ctrl+/`.

## Bloqueos
- Nuevo: `docs/estado/bloqueos/WF3a-bis-gastos-presupuesto.md` (7 puntos).

## Pendientes
- Solicitud 1, para frontend-base: montar `FormIngreso` y `LineaCredito` en `src/shell/HojaRegistro.tsx`. Cierra P-03 en C3; las pruebas de P-03 se activan solas.
- Solicitud 2, para frontend-base: portada y filtro de galería en `DefVista`. Después, quitar la detección por DOM en `P02Gastos.tsx` (`Ticket`).
- Solicitud 3, para el integrador: regenerar C3 con `--escribir` y actualizar `C3-matriz.md`.
- WF-3b: probar contra el servidor real con la semilla de desarrollo.

## Propuestas para Axel
- **PREGUNTA (bloqueo §2)**: ¿el contrato agrega el id del adjunto a cada gasto? Sin eso, un PDF adjuntado antes no se puede ver ni quitar desde la pantalla del gasto.
- **PREGUNTA (bloqueo §4)**: al terminar la compra, ¿los artículos de la lista se agregan solos al ticket o con un toque (así está hoy)?
- **PROPUESTA (bloqueo §1)**: agregar los filtros tipo, forma de pago y «con ticket» a la vista genérica del servidor, para que los totales por categoría y por día siempre sean del filtro completo.

## Paquetes instalados
Ninguno.

## Consumo aproximado
Unas 54 acciones, unos 270 KB leídos y unos 80 KB escritos.
