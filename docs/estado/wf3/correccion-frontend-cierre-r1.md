# correccion-frontend-cierre-r1 · WF-3a cierre (12-oct-2026)

Agente: `correccion-frontend-cierre-r1` (Opus). Base: HEAD `0197dad` + árbol de trabajo (cambios de `backend-cierre`, `frontend-cierre` y sus correcciones). Corrige los 2 hallazgos medios de `docs/estado/wf3/REVISION-frontend-cierre-r1.md` (M-1 y M-2; no hay altos). Los bajos (B-3 a B-6) no estaban en el encargo y no se tocaron. Sin `git add` ni `commit`.

## Hecho
- **M-1 · unidad de las 5 claves de K54 en P-26.02** (U4 §0.5, U5-06):
  - HECHO: el servidor manda en `unit` la tercera columna de `@param` de `00005_tablero.sql:25-29` (`parseParamDefs`, `Unit: f[2]`), una descripción de 57 a 123 caracteres; C1 (`SettingParameterInteger.unit`) promete «Unidad en español (intentos, minutos, días, MB…)».
  - Bloqueo nuevo `docs/estado/bloqueos/C1-unidad-parametros-tablero.md` con las opciones (A) la API da unidad corta sin editar la 00005, (B) unidades propias de P-26.02 en `textos.ts`, (C) `Parametro` muestra la unidad solo si es corta; impacto y PROPUESTA (A definitiva, B como puente).
  - PROVISIONAL mientras se decide: opción B. `Parametro` (`comunes.tsx`) acepta una prop opcional `unidad` que reemplaza `param.unit` (sin ella, S01, S05 y S06 siguen igual); P-26.02 le pasa `t.presupuesto.unidades[clave]` (`%`, `%`, `pagos`, `meses`, `puntos`: la primera palabra de la descripción del servidor). Ya no se repite la descripción ni aparece «de 1 a 100» en las filas.
  - Pruebas: los datos falsos ya no inventan unidades. `pruebas/ajustes-a/unidades-servidor.ts` (nuevo) lee `unit` de las mismas anotaciones `@param` de las migraciones que lee el servidor; lo usan `parametros-tablero.escritorio.spec.ts` y `ajustes.escritorio.spec.ts`. La prueba de las 5 claves comprueba que cada `unit` del servidor es larga (> 20 caracteres), que no aparece en el bloque, que el único «de 1 a 100» es el de la nota de rangos, y «Valor inicial: 30 % / 40 % / 5 pagos / 3 meses / 10 puntos» y «4 meses».
  - Comprobación contra el servidor real (copia limpia en 8444, prueba temporal ya borrada): ver Verificación.
- **M-2 · K-4 «Hogar $973.75» mientras `listPersons` no responde** (C3-09, U5-06):
  - `Tablero.tsx` (`K4`): no se pinta ninguna fila (esqueleto) hasta que `c.personas.datos !== undefined`; la línea «Hogar · Karen» y la insignia (que usa `idK` para «Falta la parte de Karen / del hogar») ya no se calculan con `idK` indefinido.
  - Si `listPersons` falla o nunca se cargó sin conexión: estado de error del bloque («No se pudieron cargar los pagos…» · «Reintentar»), igual que cuando falla la lectura de N (U3 §1.2); «Reintentar» vuelve a pedir lo que falló. Esta conducta es PROVISIONAL: bloqueo nuevo `docs/estado/bloqueos/U3-tablero-k4-personas.md` (U3 no dice qué hace K-4 si falla `listPersons`).
  - Dos pruebas nuevas en `pruebas/inicio/tablero.escritorio.spec.ts`: (1) retiene `/api/v1/persons` hasta que los pagos ya respondieron: 0 filas, sin `$973.75` ni «Hogar» en K-4; al soltarla, «Hogar $598.82 · Karen $374.93»; un `MutationObserver` comprueba que «Hogar $973.75» **nunca** estuvo en la página. (2) `/api/v1/persons` = 500: error de K-4 sin cifras, K-1 sigue, «Reintentar» lo recupera con las cifras del oráculo.
  - Comprobé que las dos pruebas **fallan sin la corrección** (4 fallas en movil y escritorio con el `Tablero.tsx` anterior) y pasan con ella.

## Archivos
- `frontend/src/pantallas/ajustes-a/comunes.tsx` (prop `unidad` de `Parametro`), `frontend/src/pantallas/ajustes-a/S02Presupuesto.tsx` (`ParametrosTablero` pasa `unidad`), `frontend/src/pantallas/ajustes-a/textos.ts` (`presupuesto.unidades`)
- `frontend/src/pantallas/inicio/Tablero.tsx` (`K4`)
- `frontend/pruebas/ajustes-a/unidades-servidor.ts` (nuevo), `frontend/pruebas/ajustes-a/parametros-tablero.escritorio.spec.ts`, `frontend/pruebas/ajustes-a/ajustes.escritorio.spec.ts` (dato falso de `/api/v1/settings/dashboard`; +1 renglón por el `import`)
- `frontend/pruebas/inicio/tablero.escritorio.spec.ts` (2 pruebas nuevas)
- `docs/estado/bloqueos/C1-unidad-parametros-tablero.md` (nuevo), `docs/estado/bloqueos/U3-tablero-k4-personas.md` (nuevo)
- `frontend/dist/` recompilado (lo sirve el binario desde disco; el servidor de 8443 ya sirve esta versión). Capturas de flujos reescritas por las corridas: las mismas 4 que ya estaban modificadas y `F3A-tarjeta-sin-precios.jpg`.

## Verificación
- `npx tsc --noEmit` (incluye `pruebas/`): salida 0. `npm run build`: OK, PWA 80 entradas.
- **U4 P-26.02 §0.5 / U5-06 / U5-07** · `npx playwright test pruebas/ajustes-a/ pruebas/inicio/`: **113 passed, 1 skipped**, rc 0 (antes 109 + 1; +4 = 2 pruebas nuevas × 2 proyectos). Incluye las 4 de `parametros-tablero` (5 claves con la `unit` real, vista previa, Esc 0 escrituras, Enter 1 PATCH, V-23 «del 1 al 39», `reset_parameters`; `erroresC1` = []; `textosEnIngles` = []; `controlesChicos` = [] en movil).
- **U4 P-26.02 contra el servidor real** (copia limpia 8444, prueba temporal borrada): longitudes de `unit` que manda la API = 123, 113, 57, 88 y 112; ninguna aparece en el bloque. Texto: «Valor inicial: 30 % 30 % … Valor inicial: 40 % 40 % … Valor inicial: 5 pagos 5 pagos … Valor inicial: 3 meses 3 meses … Valor inicial: 10 puntos 10 puntos · Rangos: … PROVISIONAL (bloqueo C1-tablero)». `scrollWidth` 412 de 412.
- **C3-09 / M-2 contra el servidor real** (misma copia, `listPersons` retenida con `route`): K-4 retenida = «Pagos próximos · Ver calendario», 0 filas; suelta = «Vie 30 oct 2026 · Kueski · $973.75 · Pendiente · Hogar $598.82 · Karen $374.93»; «Hogar $973.75» nunca apareció.
- **C3-09** (`c3-09.spec.ts`, copia limpia): «C3-09: **120 cifras** leídas en pantalla · diferencias != 0: 0»; K-4 973.75 / 598.82 / 374.93 con dif 0.
- **C3-08** (cada archivo en copia nueva de `servidor-prueba.sh`): `flujos.spec.ts` **7 passed** (F-1, F-3A0, F-2, F-3, F-4, F-5, F-6); `pendientes.spec.ts` **4 passed** (F-5, F-4b, F-2 copiar, F-3A3).
- **C3-03** `npm run verificar:rutas`: «213 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (registradas: 221)». **U5-08** `npm run verificar:hosts`: «hosts externos que se pedirían = 0», rc 0.
- `make verify-wf3-frontend`: rc **0**; suite completa **271 passed, 4 skipped, 0 fallas** (6.0 min; antes 267). Los «forma inesperada» del registro son las pruebas de `BarreraErrores` del banco, esperados.
- `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: rc **0** («== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)»).
- 0 hex: en las líneas agregadas de `frontend/src` hay 0 `#hex`, 0 `rgb(`, 0 `hsl(`, 0 `oklch(`.
- Temporales: `datos/tmp/cfc-r1`, `datos/tmp/u5-flujos`, `datos/tmp/playwright-u5-flujos` y `datos/tmp/playwright-frontend` borrados; servidor de 8444 detenido; el de 8443 no se tocó (solo sirve el `dist` nuevo).

## Supuestos
- Unidades cortas de B (`%`, `pagos`, `meses`, `puntos`): la primera palabra de la descripción que hoy manda la API para cada clave; no son decisión nueva y quedan PROVISIONALES en el bloqueo.
- K-4 sin `listPersons`: estado de error por analogía con la falla de N (U3 §1.2); PROVISIONAL en `U3-tablero-k4-personas`.

## Bloqueos
- Nuevos: `C1-unidad-parametros-tablero` (de dónde sale la unidad corta) y `U3-tablero-k4-personas` (K-4 si falla `listPersons`). Siguen PROVISIONALES: `C3-ajustes-tablero` (opción A aplicada), `C1-tablero` (rangos), `C3-flujos`.

## Pendientes
- Ninguno de mi alcance.
- Regenerar `docs/contratos/C3-integracion.md` (ya pedido por `frontend-cierre`, B-6): ahora cambian también los números de renglón de `ajustes.escritorio.spec.ts` (+1) y suben las cuentas PW de P-01 (2 pruebas nuevas) y de P-26.02.

## Propuestas para Axel
- Puerta 4: decidir `C1-unidad-parametros-tablero` (recomendada A: el servidor da la unidad corta que C1 ya promete; se retira el puente B) y `U3-tablero-k4-personas` (recomendada A: error de K-4 con «Reintentar»).

## Solicitudes fuera de alcance
- `docs/estado/bloqueos/INDICE.md`: agregar las filas de `C1-unidad-parametros-tablero.md` y `U3-tablero-k4-personas.md` (agente: correccion-frontend-cierre-r1; fase: puerta 4).
- `frontend/src/pantallas/inicio/Detalle.tsx:92-98` (D-3) y `frontend/src/pantallas/inicio/Calendario.tsx:33-42`: mismo defecto que M-2 (con `idK` indefinido D-3 dice un instante «Hogar $973.75» y el estado del pago ignora la parte de Karen). Cambio: pintar las filas solo con `c.personas.datos !== undefined` (esqueleto mientras tanto) y el mismo estado de error si falla.
- Opción A de `C1-unidad-parametros-tablero` (backend), solo si se decide: migración nueva que vuelve a anotar las 5 claves con `@param` y unidad corta (la última anotación gana en `parseParamDefs`), sin editar la 00005; después quitar `presupuesto.unidades` y la prop `unidad` en P-26.02.
- `docs/contratos/C3-integracion.md`: regenerar con `--pw-json` y `--go-log` y `--escribir` después de las capturas de U5 (como pidió `frontend-cierre`).

## Paquetes instalados
- Ninguno.

## Consumo aproximado
- ≈ 47 acciones, ≈ 115 KB leídos, ≈ 22 KB escritos.
