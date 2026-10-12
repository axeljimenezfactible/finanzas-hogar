# correccion-especificacion-r1 · Corrección de la ronda 1 de la especificación (WF-3c, K54) · 12-oct-2026

Agente `correccion-especificacion-r1` (Opus). Corrige los hallazgos altos y medios de `docs/estado/wf3/REVISION-especificacion-r1.md`: H-1, H-2 y M-1 a M-4. Los bajos (B-1 a B-5) no entran en esta ronda. Sin `git add` ni `git commit`, sin Telegram, sin sudo y sin paquetes nuevos. No cambia ninguna cifra, ningún id de criterio, ningún token ni ninguna operación de C1.

## Hecho
- **H-1 (alta) · HECHO en mi parte; los manejadores siguen PENDIENTES de `api-tablero`.**
  - En `docs/estado/wf3/s1-c1-v15.md`, §Verificación ya lista las **4 pruebas de los 3 paquetes** que fallan:
    - `cmd/finanzas`: `TestAreas_SinOperacionesEn501`.
    - `internal/api`: `TestContrato_Operaciones` (5 subpruebas) y `TestContrato_Rutas`.
    - `internal/auth`: `TestPermisos_MismasOperacionesAmbasCuentas`.
  - Ya no la da por «esperada»: queda como falla que bloquea el cierre de WF-3c.
  - Pendientes de ese reporte: `api-tablero` registra también las 2 operaciones de Ajustes.
  - No escribí manejadores. `backend/internal/api/tablero/`, `cmd/finanzas/areas.go` y los paquetes de área son del agente `api-tablero` (fase Construcción de `wf3c-tablero.js`); ver solicitudes.
- **H-2 (alta) · HECHO.**
  - `docs/criterios/C1.md`:
    - C1-09 dice «los bloques del tablero de U3 2.0 (K54: K-1 a K-4, B-1, G-1 a G-6 y Más del hogar)», y «Cubre» suma M-18 a M-22. Su Verifica habla de bloques del tablero y de las tarjetas de Más del hogar.
    - §4: «bloques del tablero de U3 2.0 (K-1 a K-4, B-1, G-1 a G-6 y Más del hogar; K54) (C1-09)».
    - §5: línea nueva de trazabilidad con K54 → C1-09. La línea de §G queda como cita histórica.
    - Encabezado y «Cambios de la versión 1.5» con la corrección.
    - Sigue en 1.5, con los mismos ids y severidades.
  - `backend/internal/api/c1gen/c1_plantilla.md` §9: «Los bloques del tablero de Inicio › Resumen (U3 2.0, K54; antes 6 tarjetas, U1 N-020 a N-025)».
  - Regeneré con `build.py`. Solo cambia `docs/contratos/C1-datos.md`, en la línea 400. `openapi.yaml`, `routes.go`, `types_gen.go` y los extractos quedan idénticos.
  - En `gen2_schemas.py` solo cambié el comentario «Inicio: 6 tarjetas». No se emite.
- **M-1 (media) · HECHO** en `docs/design/U3-dashboard.md` 2.0:
  - §2.6 «publicadas en C1 v1.5»: la ventana inválida responde 422 `validation_failed` y el PATCH lleva `parameters` y `reset_parameters`.
  - Las métricas M-18 a M-22 están publicadas en C1-datos §9, con RN-30 y RN-31.
  - B-1, G-1, G-3 y G-6 citan M-20, M-21 y M-22 «(C1 v1.5)», no como propuesta.
  - Regla de §2: los campos nuevos están publicados. La caída solo aplica a un servidor anterior.
  - K-4: N = `parameters["dashboard.payment_reminder_count"].value` de `getDashboardSettings`, leído una sola vez como `getAlertSettings`. «La interfaz no fija N».
  - §1.2: viñeta nueva sobre esa lectura. Si falla, K-4 muestra su error de §7. Es PROVISIONAL, bloqueo `U3-tablero-recordatorio`.
  - A-02 nombra las dos lecturas únicas.
  - «Fuentes que rigen»: C1 v1.5.
  - §13, §14 y §15 al día: C1 v1.5 y el esquema 1.5 están HECHOS. Falta el backend (H-1).
  - Línea de cambios «Corrección r1».
- **M-2 (media) · HECHO**, con PREGUNTA abierta, en `docs/design/U1-arquitectura.md`, sin cambio de versión: `verificar_u1.py` exige 1.2.
  - Regla de subsecciones (§2): en Inicio, el Resumen es el tablero de U3 2.0 (K54, sin tope). En las demás páginas sigue §G; ver PREGUNTA 1.
  - P-01 en §2.1 (propósito y componentes), §2.2 (vacío, carga y error por bloque), §6.1 (fila de tarjetas del Resumen), §6 (P-01 en 412 y 1440) y la tabla de dependencias.
  - §3.3:
    - Encabezado «Resumen = tablero de control de U3 2.0» con el destino de cada bloque.
    - Tabla nueva «Necesidades del tablero (K54)»: 11 filas, una por bloque. Cada fila trae la N-xxx vigente que lo lleva y su operación de C1 v1.5. Las N-xxx salen de las `x-necesidades` de los extractos.
  - F-4: paso 1 «Toca el pago en K-4 («Pagos próximos» del tablero, K54)» y su resumen «Pago en K-4 → Calendario».
  - Agente anotado en el encabezado.
  - Las N-xxx no son nuevas. Usé la opción (a) PROVISIONAL del punto 3 de `C1-tablero`. Ver PREGUNTA 2.
- **M-3 (media) · HECHO.**
  - DESIGN §10.5 suma C-159 a C-172, 14 pares del barómetro en claro y oscuro: tramos (3), valor sobre la pista (3), insignias (3), dígito, texto de las marcas, marca de umbral, aguja sobre la pista y la pista (decorativa).
  - Introducción: 172 pares × 2 = 344 combinaciones, con 330 con umbral, 14 decorativas y 58 que no están en `contraste.csv`. P8 dice 344.
  - `contraste_componentes.py`: subí `ESPERADOS` y el docstring, y regeneré `docs/design/contraste-componentes.csv`.
  - La tabla de §8.9 lleva la columna Id y las cifras del script (truncadas a 2 decimales, §2.1). El cálculo de `u3-v2` redondeaba: 8 cifras bajan 0.01, por ejemplo 4.74 → 4.73 y 10.76 → 10.75. Ninguna queda bajo su mínimo.
  - Suma la marca de umbral (3.46 / 6.08) y la aguja (9.63 / 7.42). El revisor redondeó a 3.47 / 6.09 y 9.64 / 7.42.
  - Línea de cambios en la versión 1.1.
- **M-4 (media) · HECHO.** La Verifica de C1-13 busca `grep -n -m1 "version: \"1.5.0\""`. Rc 0, línea 7.
- **Notas en los reportes de origen:**
  - `s1-c1-v15.md` y `u3-v2.md`: sección «Corrección r1».
  - `docs/estado/bloqueos/C1-tablero.md`: nota del punto 3. Sigue abierto.

## Archivos
- `/home/deepsec/finanzas-hogar/docs/criterios/C1.md`
- `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/c1_plantilla.md`, `/home/deepsec/finanzas-hogar/backend/internal/api/c1gen/gen2_schemas.py` (solo un comentario)
- `/home/deepsec/finanzas-hogar/docs/contratos/C1-datos.md` (regenerado)
- `/home/deepsec/finanzas-hogar/docs/design/U3-dashboard.md`
- `/home/deepsec/finanzas-hogar/docs/design/U1-arquitectura.md`
- `/home/deepsec/finanzas-hogar/docs/design/DESIGN.md`
- `/home/deepsec/finanzas-hogar/docs/design/tools/contraste_componentes.py` y `/home/deepsec/finanzas-hogar/docs/design/contraste-componentes.csv`: fuera de la lista de archivos de u3-v2. Hacen falta para que M-3 quede comprobada por `--check`.
- `/home/deepsec/finanzas-hogar/docs/estado/bloqueos/C1-tablero.md` (nota)
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/s1-c1-v15.md`, `/home/deepsec/finanzas-hogar/docs/estado/wf3/u3-v2.md` (secciones «Corrección r1»)
- `/home/deepsec/finanzas-hogar/docs/estado/wf3/correccion-especificacion-r1.md` (este reporte)

## Verificación
- **H-1 · `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose`: EXIT 2.**
  - Pasan WF-0, C1 (22 archivos idénticos, `verify_c1 --db` TODO OK), S1 (RESUMEN TODO OK), U1 (TODO OK) y criterios.
  - Falla `verify-wf1-cobertura`, solo en las 4 pruebas de los 3 paquetes:
    - `areas_test.go:22`: «5 de 221 operaciones sin registrar (responden 501)».
    - `contrato_test.go:873` y `:906`: «S6-03: contrato 221 operaciones, servidor 221, sobrantes 0, faltantes 0, sin implementar 5».
    - `auth_test.go:331`: «Axel 501, Ximena 501».
  - Las etapas siguientes, corridas aparte (`make verify-wf1-oraculos verify-wf1-recetario verify-wf1-tokens verify-wf1-u4 verify-wf1-paquete`), dan **EXIT 0**.
  - **C1-02 / S6-03:** las 5 operaciones siguen en 501 hasta `api-tablero`.
- **Repaso del grupo después de las ediciones:** `make verify-c1 verify-s1 verify-u1 verify-criterios verify-wf1-tokens`, **EXIT 0**.
  - C1: 22 archivos idénticos y TODO OK.
  - S1: RESUMEN TODO OK.
  - U1: TODO OK. Las líneas PROVISIONAL nuevas, 472 y 473, citan bloqueos abiertos (`U3-tablero-mas-del-hogar` y `C1-tablero`).
  - Criterios: 261 (201 / 60), 326 citas, 0 sin fuente.
  - Tokens: `contraste_componentes.py --check` IDENTICO, 0 bajo umbral.
- **C1-09 y C1-13:**
  - `grep -n "6 tarjetas" docs/criterios/C1.md`: solo la lista de cambios, la cita histórica de §G en §5 y la línea de K54. En §4 ya no hay umbral de 6 tarjetas.
  - `C1-datos.md:400`: «Los bloques del tablero … (U3 2.0, K54; antes 6 tarjetas…)».
  - `grep -n -m1 "version: \"1.5.0\"" docs/contratos/openapi.yaml`: «7:  version: "1.5.0"», rc 0.
- **Criterios (ids):** comparé los ids y severidades `**XX-nn (sev)**` de los 18 archivos entre HEAD y la copia de trabajo. 0 diferencias.
- **C1 regenerado:** `build.py` y después `build.py --check`: «los 22 archivos coinciden byte a byte». `cmp` de `openapi.yaml` antes y después: idéntico.
- **U3-02 y U3 §2.6:** en U3, `grep -c` de «400 `validation_failed`», «`reset`», «C1 v1.2 por sus» y «mientras C1 no la publique» da 0. En la ronda r1 comprobé en los extractos `inicio.yaml` y `ajustes-e.yaml` que existen el 422 de `getDashboardMonthlyOverview`, `reset_parameters` en `updateDashboardSettings` y `DashboardSettings.parameters` con `SettingParameterInteger.value`.
- **U2-04 (M-3):**
  - `python3 -I docs/design/tools/contraste_componentes.py --check`: «pares: 172 · combinaciones: 344 · con umbral: 330 · decorativas exentas: 14 · bajo umbral: 0 · cruce con contraste.csv: igual 286 · nuevo 58 · difiere 0 · cruces secundarios: 0», rc 0.
  - `--independiente`: «10 razones comprobadas, todas coinciden», rc 0.
- **U1 (M-2):** `make verify-u1` TODO OK. En U1 `grep "6 tarjetas"` deja solo las menciones históricas («antes 6 tarjetas», «decisión G: como máximo 6»).

## Supuestos
1. H-1 se corrige así: el reporte se pone al día y los manejadores van a `api-tablero`. Así lo dice la corrección del revisor y así lo ordena `wf3c-tablero.js`, donde `api-tablero` corre justo después de esta revisión. Escribir manejadores aquí invadiría su alcance.
2. M-2 sin N-xxx nuevas. Las propias del tablero (opción b de `C1-tablero`, punto 3) piden U1 1.3, y `verificar_u1.py` exige la versión 1.2 y que las N-xxx nuevas sean N-320 a N-327. Con N-328 o más, `verify-u1` y `verify-all` fallarían. Por eso §3.3 documenta las necesidades del tablero sobre las N-xxx vigentes (opción a, PROVISIONAL).
3. Las cifras de §8.9 son las del script oficial, truncadas (DESIGN §2.1). No son las redondeadas del primer cálculo.
4. Que K-4 muestre su estado de error de §7 si falla la lectura de `getDashboardSettings` es PROVISIONAL (bloqueo `U3-tablero-recordatorio`). Así no queda N fijo en el código.

## Bloqueos
No abrí ninguno nuevo. Siguen `C1-tablero` (punto 3 con la nota) y los cinco `U3-tablero-*`.

## Pendientes
- **`api-tablero`** registra los 5 manejadores con los motores S2 y S3, sin recalcular cifras:
  - Inicio: `getDashboardDebtCapacity`, `getDashboardMonthlyOverview` y `getBudgetSignals`.
  - Ajustes: `getDashboardSettings` y `updateDashboardSettings` (K44, con `reset_parameters` y 422 si el umbral amarillo no es menor que el tope).
  - Con eso pasan las 4 pruebas y `make verify-all` da EXIT 0. Sin eso, el integrador no cierra WF-3c con C1 v1.5.
- **Dueño de U1 o de sus herramientas**, si Axel elige la opción (b): U1 1.3 con N-xxx propias del tablero, `verificar_u1.py` aceptando la 1.3 y C1 citándolas.
- **B-1 a B-5** del revisor (bajos): fuera de esta ronda. B-4 (`docs/estado/bloqueos/INDICE.md`) es del integrador.

## Preguntas y propuestas para Axel (puerta 4)
1. **PREGUNTA.** K54 retira el tope de 6 tarjetas en el Resumen de Inicio. ¿Sigue el tope de §G en el Resumen de las demás páginas (P-02, P-18…)? U1 lo deja así porque K54 dice «en esa parte sustituye a §G».
2. **PROPUESTA** (`C1-tablero`, punto 3): mantener la opción (a), el tablero sobre las N-xxx vigentes, hasta la puerta 4. La opción (b), con N-xxx propias, cuesta U1 1.3 y cambiar el verificador de U1 y C1, y no cambia ninguna pantalla.

## Paquetes instalados
Ninguno.

## Consumo aproximado
46 acciones; unos 170 KB leídos (por secciones con grep y sed; `openapi.yaml` solo por extractos); unos 40 KB escritos a mano. Las salidas del generador y el CSV se regeneran por script.
