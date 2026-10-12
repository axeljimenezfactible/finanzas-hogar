# front-ajustes-b-fix (WF-3a bis) · `Ctrl+/` en Ajustes 09, 13 y 14 (H-08) y pruebas de `reports.*` y `alerts.*` (R2-06)

Agente: `front-ajustes-b-fix` (Opus). Fecha: 12-oct-2026. Alcance: `frontend/src/pantallas/ajustes-b/` y sus pruebas (`frontend/pruebas/ajustes-b/`).
Hallazgos que cierra: **H-08** (media) y **R2-06** (baja) de `docs/estado/wf3/REVISION-frontend-areas-r3.md` (INTEGRACION-wf3a §7).
Sin cambios de cifras, motores ni decisiones; 0 operaciones nuevas (las tres escrituras ya estaban en la pantalla y en los extractos).

## Hecho

**HECHO.** `Ctrl+/` conectado con la pieza de la base (`EditarVarias`, `EdicionVarias`, CONVENCIONES-frontend §7), con vista previa sumada «Antes → Después» antes de aplicar; `Enter` aplica a todas y `Esc` cancela sin cambiar ninguna (U4 §0.2 punto 2 y §0.7).

| Sección | Qué se edita a la vez (U4) | Vista previa sumada | Operación (ya en C1) |
|---|---|---|---|
| P-26.09 Notificaciones | «Activa» de las alertas seleccionadas (U4 P-26.09, Atajos) | «Alertas activas: 4 de 4 → 2 de 4» y una fila por alerta que cambia («Pago próximo: Activa → Inactiva») | `updateNotification` (PATCH `/api/v1/notifications/{notification_id}`), una por alerta |
| P-26.13 Metas | «Activa» de varias metas (U4 P-26.13, vista previa (3) y Atajos) | «Total reservado para metas en {mes}: $1,500.00 → $0.00»; si las que cambian no tienen aportación, «… · Sin aportación planificada: el total reservado no cambia» | `updateSavingsGoal`, una por meta |
| P-26.14 Inventario | Modo o nivel de reposición (U4 P-26.14, vista previa y Atajos) | «{k} pasan a Reponer, {j} a Agotado y {s} a Suficiente» con el conteo de estados de las seleccionadas antes → después (regla K7 y RN-15 en el navegador, igual que la vista previa de una fila) y la lista de las filas que cambian («Frijol: Suficiente → Reponer»); sin ritmo, «Falta el ritmo» | `updateInventoryItem`, una por fila |

- 13 y 14 usan `VistaBD` con `def.editarVarias` (casillas, «Seleccionar todas», barra «n seleccionadas · Editar seleccionadas (Ctrl+/) · Quitar selección» de la base).
- 09 no es una `VistaBD` (las alertas son tarjetas con destinatarios y parámetros): `ajustes-b/varias.tsx` (nuevo) da `GrupoVarias`, que replica la selección y el atajo de `VistaBD` (casilla por tarjeta de 36 px fino / 44 px grueso, «Seleccionar todas», la misma barra, `Ctrl+/` con el foco dentro, foco de regreso) y monta el mismo `EdicionVarias` de la base. No toca `src/vistas-bd/`.
- `guardar` manda una escritura por fila y siempre recarga (si una falla, el bloque sigue abierto con el mensaje y la lista vuelve a los valores del servidor, U4 §0.3).
- **R2-06**: dos pruebas nuevas que editan una clave y comprueban el cuerpo exacto: `PATCH /api/v1/settings/reports` `{ parameters: { 'reports.daily_send_time': '22:15' } }` (P-26.08) y `PATCH /api/v1/settings/alerts` `{ parameters: { 'alerts.send_time': '08:30' } }` (P-26.09), con respuesta falsa del esquema de C1 (`conCambios`, ya existente) y validación de todo el tráfico con `erroresC1` (0 errores).
- Datos de prueba: un tercer renglón de inventario (`Frijol`, existencia por unidad 1,500 g, nivel 500 g) para que la prueba de 14 no dependa de la fecha del día (el aceite va por ritmo y su estado cambia con los días); respuestas falsas de C1 para PATCH de `notifications/{id}` y `savings-goals/{id}`, y la de `inventory-items/{id}` ahora refleja `reorder_level_milli`.

## Archivos

- `frontend/src/pantallas/ajustes-b/varias.tsx` (nuevo, 4.6 KB): `Casilla` y `GrupoVarias<T>`.
- `frontend/src/pantallas/ajustes-b/P2609Notificaciones.tsx`: `grupoAlertas()` y la casilla en cada tarjeta.
- `frontend/src/pantallas/ajustes-b/P2613Salud.tsx`: `editarVarias` de metas.
- `frontend/src/pantallas/ajustes-b/P2614Inventario.tsx`: `previaVarias`, `cambioVarias` y `editarVarias`.
- `frontend/src/pantallas/ajustes-b/textos.ts`: `n.inactiva`, `n.alertasActivas`, `n.deTotal`, `i.pasan`, `i.conteo` (es-MX).
- `frontend/pruebas/ajustes-b/ajustes-b.escritorio.spec.ts`: 5 pruebas nuevas (2 de R2-06 y 3 de H-08).
- `frontend/pruebas/ajustes-b/datos.ts`: renglón de inventario 3 y respuestas falsas de 3 PATCH.
- `docs/estado/wf3/front-ajustes-b-fix.md` (este reporte).

## Verificación

| Criterio | Comando | Resultado |
|---|---|---|
| Tipos | `npx tsc --noEmit` (proyecto entero) | 0 errores (en una corrida intermedia aparecieron 3 de `gastos-presupuesto/P02Gastos.tsx`, de otro agente en curso; en la corrida final, 0) |
| Build | `npx vite build --outDir ../datos/tmp/ajb-fix/dist2` (misma configuración; outDir propio para no pisar `dist/` ni el preview de otros agentes) | `✓ built`, precache de 81 entradas |
| Pruebas del área | `npx playwright test --config datos/tmp/ajb-fix/pw.config.ts ajustes-b --repeat-each 2` (preview propio en 4193) | **60 passed**, 0 failed, 0 flaky (15 pruebas × 2 proyectos × 2 vueltas); antes, una vuelta: 30 passed |
| U4 §0.2/§0.7 · H-08 09 | prueba «H-08 · P-26.09 Ctrl+/…» | vista previa: «Alertas activas 4 de 4 → 2 de 4 · Pago próximo Activa → Inactiva · Precios viejos antes del día de despensa Activa → Inactiva»; `Esc`: 0 PATCH, la selección sigue y el foco regresa a la última casilla; `Enter`: PATCH `notifications/4` y `notifications/6` con `{ is_enabled: false }`, aviso «Guardado en 2 registros» y selección limpia |
| U4 · H-08 13 | prueba «H-08 · P-26.13 Ctrl+/…» | «Total reservado para metas en oct-2026 $1,500.00 → $0.00»; `Esc`: 0 PATCH; `Enter`: PATCH `savings-goals/1` y `/2` con `{ is_active: false }` |
| U4 · H-08 14 | prueba «H-08 · P-26.14 Ctrl+/…» | «1 pasa a Reponer, 0 a Agotado y 0 a Suficiente · Suficiente 1 · Reponer 0 · Agotado 1 → Suficiente 0 · Reponer 1 · Agotado 1 · Frijol Suficiente → Reponer» (Arroz no aparece: no cambia); modo «Ritmo de consumo» → «Falta el ritmo»; `Esc`: 0 PATCH; `Enter`: PATCH `inventory-items/2` y `/3` con `{ reorder_level_milli: 2000000 }` |
| Táctil (U4 §0.7: en Android no hay atajos) | las tres pruebas en `movil` 412×915 | el mismo bloque se abre con «Editar seleccionadas», se cancela con «Cancelar» y aplica con «Aplicar a 2»; el foco regresa a «Editar seleccionadas» |
| U5-06 (09, nuevo control) | `controlesChicos` en `movil` | 0 controles de selección menores de 44 px |
| U5-07 | `textosEnIngles(page)` con el bloque abierto en 09, 13 y 14 | `[]` en los dos proyectos |
| R2-06 | pruebas «R2-06 · P-26.08…» y «R2-06 · P-26.09…» | cuerpos exactos (arriba); `erroresTrafico` = `[]` en todas las pruebas (cuerpos y respuestas contra C1) |
| C3-03 (0 rutas fuera del contrato) | `npm run -s verificar:rutas` | «205 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0» |
| Hex (U5) | `grep -E '#[0-9a-fA-F]{3}\b|#[0-9a-fA-F]{6}\b'` en los archivos tocados | 0 |
| Global | `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` | salida 0, «Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)» |

## Supuestos

- PROPUESTA (no PREGUNTA, porque no toca cifras ni decisiones): U4 P-26.09 pide `Ctrl+/` para activar o desactivar alertas pero no dice el contenido de la vista previa; como la base no aplica sin vista previa calculada, se muestra lo único que cambia (cuántas alertas quedan activas y cuáles cambian). Si Axel prefiere otro texto, se cambia en `textos.ts`.
- En 14 la línea de U4 «{k} pasan a Reponer, {j} a Agotado y {s} a Suficiente» va como etiqueta del primer renglón del bloque (`Cambio` de la base solo tiene etiqueta, antes y después) y el antes → después es el conteo de estados; con k = 1 dice «pasa». Cuenta solo las filas que cambian de estado.
- En 14, el nivel tecleado se interpreta en la unidad base de cada producto (g, ml o pieza), igual que la edición de una sola fila; un valor que no es cantidad no calcula vista previa («No se pudo calcular el efecto · Reintentar» de la base) y por lo tanto no aplica.
- El grupo escribe una vez por fila seleccionada (también en las que ya tenían ese valor), como el ejemplo de CONVENCIONES §7 y U4 P-26.02 (h).

## Bloqueos

Ninguno nuevo.

## Pendientes

- Ninguno del alcance. R2-05 y R2-07 son de otras áreas; R3-01 (`ajv` sin declarar y `validar-c1.ts`/`tablero.ts` en `ajustes-a/`) sigue en la base (ver solicitud).

## Propuestas para Axel

- Revisar el texto de la vista previa de 09 («Alertas activas: a de 4 → b de 4» más la lista) y del bloque de 14 (frase de U4 como etiqueta con el conteo de estados antes → después).

## Paquetes instalados

Ninguno.

## Consumo aproximado

~50 acciones; ~150 KB leídos; ~30 KB escritos. Temporales en `datos/tmp/ajb-fix/` (config de Playwright, builds y salidas): borrados al terminar.
