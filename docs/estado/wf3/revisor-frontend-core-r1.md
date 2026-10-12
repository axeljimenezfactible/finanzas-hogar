# revisor-frontend-core-r1 · WF-3a · reporte (12-oct-2026)

## Hecho
- Revisión independiente, ronda 1 de 3, del trabajo de frontend-base, front-inicio, front-gastos-presupuesto y front-comidas-compras contra U5 (U5-06..U5-10, U5-13 con API falsa), C3-03, DESIGN.md, U1 §2.1/§2.2/§3 y U3. No corregí nada.
- Informe completo con tabla criterio → estado → evidencia y hallazgos: `docs/estado/wf3/REVISION-frontend-core-r1.md`.
- **Veredicto: NO APROBADO**: 5 altos (H-01 cola sin aviso ante 5xx; H-02 íconos de 16 px y trazo 2; H-03..H-05 objetivos táctiles < 44 px), 2 medios (H-06 dinero sin `tabular-nums`; H-10 `make verify-wf3-frontend` no corre la verificación del frontend), 4 bajos.

## Archivos
- `docs/estado/wf3/REVISION-frontend-core-r1.md` (nuevo).
- `docs/estado/wf3/revisor-frontend-core-r1.md` (este).
- Temporales en `datos/tmp/rev-fc-r1/` (build, configuración de Playwright en el puerto 4193, 3 pruebas propias, script de rutas): **borrados**.

## Verificación (comandos propios; salidas resumidas)
- Build: `npm run preparar` exit 0 · `npx tsc --noEmit` exit 0 (0 errores) · `npx vite build --outDir ../datos/tmp/rev-fc-r1/dist` exit 0 (precache 80).
- Pruebas de las áreas: `npx playwright test -c datos/tmp/rev-fc-r1/pw.config.ts pruebas/inicio pruebas/gastos-presupuesto pruebas/comidas-compras pruebas/cola.spec.ts pruebas/shell.spec.ts pruebas/login.spec.ts` → **57 passed**.
- C3-03: `node scripts/rutas-api.mjs` → 201 rutas, 0 fuera de openapi.yaml, 0 sin registrar (216). Script propio sobre todas las cadenas `/api/v1/…`: 162, 0 fuera; `fetch(` solo en `api/cliente.ts` y `api/conexion.ts` (`/health`). **Cumple.**
- U5-08: `hosts-externos.mjs` → 0; `grep` propio en `dist` → solo literales (xmlns, licencias, mensajes); barrido Playwright: único host `localhost:4193`. **Cumple.**
- U5-07: barrido propio de 55 pantallas × (datos, API vacía) × (412, 1440): inglés 0, ISO 0, dinero mal formado 0 (ejes `$1 mil` permitidos por DESIGN §7.7). **Cumple** (con «NaN» en P-16 solo con API vacía, H-07 baja).
- U5-09: 3 toques en las 5 mediciones (K10: solo `tap/click` en controles). **Cumple.**
- U5-10: fondo oscuro `rgb(34, 35, 37)` = C2 `#222325`, claro `rgb(250, 250, 251)` = `#FAFAFB`, `filter: none`; interruptor real en sus 3 valores con sistema oscuro y claro (6/6). **Cumple.**
- U5-13: pasos 1, 2 (base) y 3, 6 (propios: misma clave, `captured_at` 30 s antes del envío, reenvío → 200 mismo `id`, 1 registro) cumplen; paso 7 cumple con 422 pero **no con 500** (sin aviso ni contador; la captura se conserva y se reenvía a los 30 s). **No cumple (H-01).**
- U5-06: Inicio 6 tarjetas exactas; demás Resúmenes ≤ 6; hex/`rgb(`/`dark:` = 0; íconos 280/280 con trazo 2 y 104 de 16 px (**no cumple**, H-02); 43 controles < 44 px a 412 táctil (**no cumple**, H-03..H-05); 3 importes sin `tabular-nums` (H-06).
- U1 §2.2: textos de vacío/carga/error de P-02..P-16 presentes; P-01 con los de U3 §7. **Cumple.**
- `make verify-all`: no ejecutado (reescribiría `frontend/dist` compartido y ocuparía el puerto 4173); su parte de frontend equivale a `npm run build`, que sí corrí por partes (ver H-10).

## Supuestos
- U5-13 paso 7 «fallo del servidor» se probó también con 500 (lectura literal); la prueba de la base usa 422.
- Objetivos táctiles medidos con `getBoundingClientRect` del control (no cuenta áreas de toque extendidas por seudoelementos; `grep` de `before:/after:` con `inset` en las áreas = 0).
- El enlace «Saltar al contenido» (1×1, oculto hasta recibir foco) no cuenta como objetivo táctil.

## Bloqueos
Ninguno nuevo.

## Pendientes
- Ronda 2 tras las correcciones de H-01..H-06 y H-10 (reejecutar las mismas mediciones de DOM y el caso 500 de la cola).

## Propuestas para Axel
- Ninguna decisión pendiente de Axel en esta revisión: los hallazgos son correcciones contra DESIGN, U1 y U5 ya aprobados.

## Paquetes instalados
Ninguno.

## Consumo aproximado
~44 acciones · ~120 KB leídos · ~30 KB escritos (pruebas temporales ~12 KB, informes ~18 KB).
