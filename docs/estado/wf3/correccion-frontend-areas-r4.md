# correccion-frontend-areas-r4 · WF-3a (12-oct-2026)

Agente: `correccion-frontend-areas-r4` (Opus). Corrige los dos hallazgos medios de `docs/estado/wf3/REVISION-frontend-areas-r4.md`: R4-01 (dueño front-ajustes-a) y R4-02 (dueño front-ajustes-b). No toqué nada fuera de esos hallazgos. R3-01 es bajo, su dueño es frontend-base y no está en mi encargo. No hice git add ni commit.

## Hecho

**R4-01 · U4 P-26.02 (fila «Si se traslapa con otro monto» → V-17) y U4 §0.3 «Error al guardar»**. Archivo: `frontend/src/pantallas/ajustes-a/S02Presupuesto.tsx`.
- `confirmar` (:73-87) ahora envuelve el `PATCH /api/v1/budgets/{budget_id}` o el `POST /api/v1/budgets` en `try/catch`. Ya no quedan promesas rechazadas sin capturar.
- Si el error es `ErrorApi` con `codigo === 'budget_overlap'` (409 de C1, `presupuesto-a.yaml` createBudget y updateBudget), la hoja sigue abierta y debajo del monto aparece `t.v.v17`: «Esta categoría ya tiene un monto en esos meses. Ajusta el que existe.», con `CircleAlert` y `text-role-error`, en `role="alert"` (:113).
- Con cualquier otro error, la hoja también sigue abierta y muestra el texto de §0.3 `t.comun.noGuardo`: «No se pudo guardar. Revisa la conexión con el servidor de casa · Reintentar» (:114). El botón Reintentar vuelve a llamar a `confirmar`. Es el mismo patrón visual de `Estado` en `comunes.tsx`: 44 px en táctil, sin mínimo en escritorio.
- No se pierde nada de lo capturado: el monto y el alcance se quedan en la hoja. El mensaje se borra al reintentar, al cambiar el alcance, al cerrar o cancelar la hoja (`cerrar`, :87) y al editar otro monto (:98).
- Imports nuevos: `CircleAlert` de lucide-react y `ErrorApi` de `api/errores`. `v17` ya tiene uso.

**R4-02 · U4 §0.3 «Error al guardar» en P-26.13**. Archivo: `frontend/src/pantallas/ajustes-b/P2613Salud.tsx`.
- Ahora `CampoP.extra` acepta también una función que recibe el estado de guardado del campo (`ReturnType<typeof useGuardado>`, :61-62), y `CampoPerfil` la llama con su `g` (:84).
- Los dos botones pasan por `g.ejecutar`, el ayudante de `ajustes-b/comun.tsx`, en lugar de `pedir(...).then(recargar)`: «Usar el de la casa» (:112-113) y «Usar la sugerida» (:125). Si el PATCH falla, debajo del mismo campo, con `IndicadorGuardado`, aparece «No se pudo guardar. Revisa la conexión con el servidor de casa» con Reintentar. Un 422 muestra el mensaje del servidor, que es el comportamiento que ya tenía el ayudante.
- No usé el `g` de `Perfil`: su indicador solo se pinta cuando falta el perfil, así que ahí el error no se vería. Por eso uso el `g` del campo.

**Pruebas nuevas** (corren en los proyectos `movil` y `escritorio`):
- `frontend/pruebas/ajustes-a/ajustes.escritorio.spec.ts:175`, «R4-01 · P-26.02 …». Edita el monto de Despensa a 4500, elige «Solo este mes» y confirma.
  - El servidor responde `409 {code:'budget_overlap', message, request_id}`, validado contra C1 con `erroresC1`. La prueba comprueba V-17 en `role=alert` dentro de la hoja, la hoja abierta, el alcance conservado y «$4,500.00».
  - Después responde 500 `internal_error`. Comprueba «No se pudo guardar…», que V-17 ya no está y que el alcance sigue marcado.
  - Después Reintentar con 200. Comprueba que la hoja se cierra, que se manda `scope=this_month` con el cuerpo `{amount_cents: 450000}` y que hubo 2 intentos fallidos con ese mismo cuerpo.
  - Cierra con `pageerror` = [].
- `frontend/pruebas/ajustes-b/ajustes-b.escritorio.spec.ts:108`, «R4-02 · P-26.13 …». El perfil de Axel trae `goal_adjust_lose_kcal: 300`, con la lectura validada contra C1. Los PATCH a `/api/v1/health-profiles/{id}` responden 500.
  - «Usar el de la casa» muestra «No se pudo guardar…» y Reintentar debajo de «Mi ajuste para bajar (kcal)».
  - «Usar la sugerida» muestra el mismo mensaje debajo de «Meta confirmada (kcal por día)», y hay 2 mensajes en total.
  - Los cuerpos son exactos: `[{goal_adjust_lose_kcal: null}, {daily_kcal_target: 2556}]`.
  - Cierra con `pageerror` = [].
- Las respuestas 500 no pasan por `erroresC1`. C1 las declara en `default` con el esquema `Error`, y `pruebas/ajustes-a/validar-c1.ts` solo lee códigos explícitos. El cuerpo está tipado con `satisfies Esquemas['Error']`. Esto queda en Propuestas.

## Archivos
- `frontend/src/pantallas/ajustes-a/S02Presupuesto.tsx` (R4-01)
- `frontend/src/pantallas/ajustes-b/P2613Salud.tsx` (R4-02)
- `frontend/pruebas/ajustes-a/ajustes.escritorio.spec.ts`: prueba R4-01; `Esquemas` agregado al import de tipos
- `frontend/pruebas/ajustes-b/ajustes-b.escritorio.spec.ts`: prueba R4-02; imports de `RESPUESTAS`, `erroresC1` y `Esquemas`
- `docs/estado/wf3/correccion-frontend-areas-r4.md` (este reporte)

## Verificación
Corrí todo de forma aislada, como el revisor: un build propio en `datos/tmp/corr-fa-r4/dist`, `vite preview` en el puerto 4798 y una configuración temporal de Playwright que extiende `playwright.config.ts`. No usé `frontend/dist` ni el puerto 4173. Ya borré los temporales.

- **C3-03 build**: `npx tsc --noEmit -p .` da rc 0 (cubre `src` y `pruebas`). `npx vite build --outDir datos/tmp/corr-fa-r4/dist` da rc 0, con «precache 81 entries». `npx tsc --noEmit -p .storybook/tsconfig.json` da rc 0.
- **C3-03 rutas**: `npm run verificar:rutas` da «208 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 216)».
- **U5-08 (referencia)**: `node scripts/hosts-externos.mjs <dist tmp>` da «hosts externos que se pedirían = 0».
- **U4 §0.3 / V-17 (R4-01 y R4-02)**: `npx playwright test -c <tmp> -g "R4-0"` da **4 passed**: 2 pruebas × 2 proyectos.
- **Pruebas Playwright de las 3 áreas** (`ajustes-a ajustes-b creditos-salud-reportes`): **119 passed, 1 skipped, 0 failed** (2.3 min). Antes eran 115 + 1 omitida; las 4 nuevas son las R4. La omitida es la misma medida táctil en escritorio. `erroresTrafico` = [] en todas, por el `afterEach`.
- **U5-06 colores**: `grep -E "#hex|rgb(|hsl(|oklch(|dark:"` sobre los dos archivos tocados da 0. Solo uso tokens: `text-role-error`, `text-text-link`.
- **U5-07**: los textos nuevos salen de `textos.ts` (`t.v.v17`, `t.comun.noGuardo`, `textos.app.reintentar`, `T.comun.noSeGuardo`). No escribí literales en la interfaz.
- **`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`**: rc 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)».
- **Lo que no medí**: no corrí las pruebas nuevas contra el código anterior, porque habría tenido que cambiar archivos del repositorio mientras otros agentes trabajan. Aun así fallarían con ese código: antes `v17` no se usaba en ningún lado y los botones de P-26.13 no tenían `catch`, así que ningún mensaje podía aparecer. La sonda del revisor ya midió el `pageerror`. Tampoco corrí `make verify-wf3-frontend`, porque escribe en `frontend/dist` y usa el 4173, ni la suite completa de todas las áreas.

## Supuestos
- HECHO, con la instrucción literal del revisor: en P-26.02, todo error que no sea `budget_overlap` muestra el texto común de §0.3, también un 4xx como `422 out_of_range`. El `Editable` de ajustes-a muestra el mensaje del servidor en los 4xx. No cambié esa regla para no decidir por mi cuenta. Queda en Propuestas.
- La V-17 no lleva Reintentar, porque el texto pide ajustar el monto que ya existe. La hoja se queda abierta para cambiar el alcance o cancelar.

## Bloqueos
- Ninguno nuevo.

## Pendientes
- R3-01 (bajo): `ajv` sin declarar en `frontend/package.json`. Su dueño es frontend-base y no está en mi encargo.
- Una revisión puntual de R4-01 y R4-02, como pide el revisor en su PREGUNTA al orquestador.

## Propuestas para Axel
- PROPUESTA: en la hoja de alcance de P-26.02, mostrar el `message` del servidor para los 4xx que no sean `budget_overlap`, como `422 out_of_range`, como ya hace el `Editable` de ajustes-a. El texto de §0.3 quedaría solo para la red y los 5xx. Hoy sigue la corrección literal del revisor.
- PROPUESTA para el dueño de `pruebas/ajustes-a/validar-c1.ts` (front-ajustes-a): que `erroresC1` use `responses.default` cuando el código no esté declarado. Así las respuestas 500 de las pruebas también se validarían contra C1.

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- Acciones: ~24.
- KB leídos: ~120.
- KB escritos: ~16.
