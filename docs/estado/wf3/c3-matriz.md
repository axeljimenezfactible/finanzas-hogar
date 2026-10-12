# WF-3b · c3-matriz · Reporte

Agente `c3-matriz` (Opus) · 2026-10-12 · HEAD `6d691d0` «[WF-3a bis]» más el árbol de trabajo con `correccion-u5-r2`. No hice commit: eso lo hace el integrador. Criterios: `docs/criterios/C3.md` 1.1 (C3-01 a C3-06 y C3-08 a C3-11) y U5-11.

## Hecho
- **C3 pasa a v1.0, sin firma** (`docs/contratos/C3-integracion.md`). La firma de C3-07 la agrega el revisor. Estado: **implementado sí 42/42, probado sí 42/42, capturado sí 41/42**. La columna «capturado» ya trae, en cada fila `P-`, las 4 rutas de `docs/design/capturas-u5/`. Con eso queda resuelto el hallazgo M-8 de `REVISION-u5-r3` (U5-11). La fila Global queda en «no», con bloqueo.
- La matriz la sigue generando el script: `python3 -I frontend/scripts/c3-integracion.py … --escribir`. En modo de comprobación sale «Matriz del documento = matriz generada: sí» y termina con código 0.
- **Extendí el script** (+45 y −6 líneas), aunque estaba fuera de mi alcance declarado (ver Supuestos 1). Agregué tres cosas:
  - La columna «capturado» se calcula del disco: 4 PNG por pantalla, no vacíos, de 412 o 1440 px de ancho.
  - Una comprobación de C3-02 por necesidad: `N-xxx` de U1 sin operación ni Local, y `N-xxx` de C1 sin fila.
  - La opción `--capturas`.
- **Corrí yo mismo todas las pruebas** (conteos en §Verificación):
  - `go test ./... -count=1`, `npm run build` y Playwright con la API simulada.
  - Contra el servidor real: los flujos de `frontend/pruebas-u5/flujos/`, en copias aisladas en `127.0.0.1:8444`; `c3-09.spec.ts` y `humo-inicio.spec.ts`, en `finanzas.local:8443`, solo lectura.
  - Una consulta de la API, solo lectura, para el lado «= API» de C3-09.
- Agregué al documento estas secciones nuevas: §4 tabla de los flujos (C3-08), §5 cifras pantalla = API = oráculo (C3-09), §6 ejecuciones, §7 SUPUESTOS y §8 historial.
- Abrí 2 bloqueos nuevos: `C3-capturas-global.md` y `C3-flujos.md`.
- Dejé intactos los archivos de `u5-flujos`: respaldé `resultados/` y `capturas/` antes de correr sus pruebas y los restauré después (`diff -rq` sin diferencias). La base de desarrollo sigue con 6 gastos y 0 pesos. Borré las copias de la base y `datos/tmp/c3m`.
- **Aviso al orquestador.** El resumen de U5 que trae mi encargo no coincide con la ronda 3 de `docs/estado/U5-revision.md`. En esa ronda U5-01 da 164 de 164 (P-26.13 a 412 px; yo medí el ancho de los 164 PNG: 0 fuera de 412 o 1440), U5-03 da 0, U5-04 da 0 bajo umbral y U5-13 pasa. U5 queda devuelto solo por M-7 (barra lateral sin iconos, U5-06 #10) y M-8 (esta matriz).

## Archivos
- `docs/contratos/C3-integracion.md`: de v0.95 a **v1.0**, ≈ 44.6 KB.
- `frontend/scripts/c3-integracion.py`: extendido; está **fuera del alcance declarado**.
- `docs/estado/bloqueos/C3-capturas-global.md`: nuevo.
- `docs/estado/bloqueos/C3-flujos.md`: nuevo.
- `docs/estado/wf3/c3-matriz.md`: este reporte.

## Verificación
- **C3-01** · `grep -c "^| P-" docs/contratos/C3-integracion.md` = **41** = T de U1 (27 `P-xx` y 14 `P-26.xx`), más 1 fila Global. El §1 trae las 3 definiciones (implementado, probado y capturado), con el formato de fila declarado.
- **C3-02** · El script imprime estas líneas:
  - «pantallas U1: 41 (P-xx 27 + P-26.xx 14); filas: 42»
  - «pantallas sin fila: 0; filas sin operación ni Local: 0»
  - «necesidades: U1 §3.29 211 N-xxx (3 Local); C1 x-necesidades 208; de U1 sin operación ni Local: 0; de C1 sin fila en U1: 0»
  - «operaciones: 216; sin pantalla: 5 (excepciones: 5, de ellas K53: 4; sin justificar: 0)»
  - Las 5 excepciones están justificadas una por una al pie de la matriz: 4 de `x-agente-externo` y `getOpenApiDocument`. Más `listPeriods`, la única operación no llamada, con su cobertura comprobada.
- **C3-03** · Salida del script: «210 rutas invocadas; fuera de openapi.yaml: 0; sin registrar en routes.go: 0 (routes.go: 216; openapi: 216 operaciones)». `cd frontend && npm run build`: sin errores. Los archivos que regenera `preparar` quedaron iguales a HEAD (`git diff` vacío).
- **C3-04** · Backend, `cd backend && go test ./... -count=1 -v`: 32 paquetes `ok`, 0 `FAIL`; 1126 pruebas y subpruebas PASS, 0 FAIL y 0 SKIP; 216 de 216 `TestContrato_Operaciones/<operationId>`. Frontend, `cd frontend && npx playwright test`: **245 ejecutadas, 241 pasan, 4 omitidas, 0 fallan** (5.2 min). La prueba inestable `correccion-bis.escritorio.spec.ts:79` pasó. Probado sí en 42/42, con 0 filas que citen una prueba que no haya pasado en esta corrida.
- **C3-05** · Salida del script: «164 esperadas (4 × 41 P-); filas con 4/4: 41; filas en «no»: ['Global']». Medí los PNG: 82 de 412 px y 82 de 1440 px, 0 fuera de ancho. Coincide con la ronda 3 de `U5-revision.md` (U5-01 164 de 164).
- **C3-06** · `grep -c "| no |" docs/contratos/C3-integracion.md` = **1** = filas con bloqueo citado: la fila Global, por `docs/estado/bloqueos/C3-capturas-global.md`.
- **C3-07** · No me toca: el documento no trae bloque de firma.
- **C3-08** · **No se cumple: pasan 5 de 8 flujos.** Tabla en §4.
  - `flujos.spec.ts`: 6 de 7 pasan; la que falla es la de F-3. `pendientes.spec.ts`: 4 de 4 pasan.
  - F-1 pasa en **3 toques** con la regla K10. También pasan F-2, F-4 y F-6; F-5 pasa con reserva (no existe la fila sep-2026).
  - F-3 falla con la semilla: la ruta A no tiene precios de ingredientes y la cifra de la ruta B difiere de U1 (H-4).
  - F-7 y F-8 no tienen prueba contra el servidor real.
  - Bloqueo: `C3-flujos.md`.
- **C3-09** · `c3-09.spec.ts` contra `finanzas.local:8443`: **111 cifras, 0 diferencias** con `model.json`; pasa.
  - Lado API: P1–P6 de `getAccrualVsCash` contra `model.json`, 18 cifras, 0 diferencias. A 1000000 y B 690550 centavos; pago 97375 del 2026-10-30; partes $598.82 y $374.93 (humo de Inicio, 1 de 1).
  - Deuda viva: 629344 centavos en `getDebtTrend`. Total: 24 cifras de la API iguales al oráculo.
- **C3-10** · Leí 10 filas (Global, P-00, P-01, P-06, P-13, P-18, P-26, P-26.05, P-26.10 y P-26.13). Los nombres son los de U1, en español de México. Las fechas de ejecución y de evidencia van en ISO (2026-10-12, 2026-10-30, 2027-04-10). Las fechas de pantalla aparecen solo cuando cito texto literal de la interfaz.
- **C3-11** · `grep -n "v1"` da la línea 3, «Versión v1.0 · 2026-10-12 · sin firma». `ls docs/estado/bloqueos/C3-*` da 3 archivos: `C3-capturas-global.md`, `C3-flujos.md` y `C3-matriz.md`. SUPUESTOS están en §7 del documento y en este reporte.
- **U5-11** · Las 41 pantallas capturadas tienen «capturado» en «sí» y las 4 rutas en la matriz.

## Supuestos
1. **Toqué un archivo fuera de mi alcance**, `frontend/scripts/c3-integracion.py`, por tres razones:
   - La matriz va entre marcadores «generado por el script; no editar a mano».
   - El script escribía «no · pendiente WF-3b» fijo en «capturado».
   - Editar a mano rompería la comprobación de C3-02, C3-03 y C3-05, que sale con 1 si la matriz difiere.

   El cambio solo agrega cosas y es compatible con lo anterior. Si el integrador no lo acepta, hay que revertirlo y C3 vuelve a «capturado no» en las 42 filas.
2. «Capturado» exige que las 4 imágenes existan y además que sean PNG no vacíos con el ancho de U5-01 (412 o 1440 px, factor 1). Uso la convención de nombres de `u5-capturas`.
3. La fila Global queda en «no» en «capturado» (PROVISIONAL): U5-01 solo pide capturas de las `P-xx`, y G-03, G-04 y G-06 no tienen captura propia.
4. «Probado» sigue con la regla PROVISIONAL del punto 4 a de `C3-matriz.md`: título de la prueba y API simulada. La evidencia contra el servidor real está en §4 y §5.
5. Para C3-08 tomo N = 8, porque U1 §4 declara F-7 y F-8 para cubrir los obligatorios de U1-05. Mi encargo hablaba de 6. Si Axel fija 6, C3-08 queda en 5 de 6, porque falta F-3.
6. Los flujos que escriben datos corrieron en copias aisladas de la base. C3-09, el humo y mi consulta de la API corrieron sobre la base de desarrollo, solo lectura. El inicio y cierre de sesión de `axel` deja registros de sesión, como ya pasaba con las pruebas de `u5-flujos`. No imprimí ni escribí ninguna contraseña.
7. La deuda viva de la API se reconoce por el monto (629344), porque la respuesta de `getDebtTrend` no trae la cadena «2027-04-10». El revisor debe confirmar el campo de fecha.
8. Tomé el estado de U5 de `U5-revision.md` ronda 3, que es la evidencia más reciente en disco, y no del resumen del encargo.

## Bloqueos
- `docs/estado/bloqueos/C3-capturas-global.md` (nuevo): «capturado» de la fila Global. Opciones: (a) PROVISIONAL «no»; (b) aceptar las capturas de P-01 y las de U5-13; (c) capturar G-03, G-04 y G-06; (d) que no aplique.
- `docs/estado/bloqueos/C3-flujos.md` (nuevo), 5 puntos:
  - F-7 y F-8 sin prueba.
  - F-3 ruta A con la semilla.
  - La cifra de F-3 (H-4).
  - La fila sep-2026 de F-5.
  - El recálculo de F-4 (H-10).
- `docs/estado/bloqueos/C3-matriz.md`: siguen PROVISIONALES los puntos 1 (P-26 «Local»), 3 c (`listPeriods`) y 4 a («probado»).

## Pendientes
- Firma de C3 (C3-07): la hace el revisor independiente, después de decidir los dos bloqueos nuevos.
- C3-08: pruebas F-7 y F-8 contra el servidor real; precios de ingredientes en la semilla para F-3 ruta A; repetir los flujos.
- Hay que regenerar la matriz con `--escribir` si cambia el frontend, la semilla o las capturas (por ejemplo, tras corregir M-7). La firma va sobre esa versión.
- Observación para el integrador: en mi corrida, la fila del pago seguía en «pagada / pendiente» después de «Deshacer», y aun así `pendientes.spec.ts:57` pasó. Hay que revisar si «Deshacer» revierte la marca o si la prueba no lo comprueba.

## Propuestas para Axel
| Decisión | Qué se obtiene | Costo | Riesgo | Recomendación |
|---|---|---|---|---|
| Probar F-7 (comprar a crédito con el comparador) y F-8 (ticket con OCR) en el servidor real | C3-08 completo con los 4 flujos obligatorios | 2 pruebas nuevas, un agente | Bajo: corren en una copia de la base | Sí, antes de firmar C3 |
| Cargar precios de los ingredientes de las recetas en la semilla | «Mejor canasta» compara la lista del plan y F-3 pasa con la semilla | Datos del agente externo (K53) o de la semilla | Bajo | Sí |
| Fila Global: aceptar como «capturado» las capturas de Inicio y las de «Sin conexión» | Matriz con 42 de 42 en «sí» | Ninguno (cambiar una constante y regenerar) | Bajo: G-03, G-04 y G-06 quedan sin imagen propia | Sí, opción (b) |
| Cifra de F-3 (+$138.50 de U1 contra la del motor) y si marcar el pago del hogar mueve el disponible | Cierra H-4 y H-10 | Una respuesta de Axel | Ninguno | Decidir en la puerta 4 |

## Consumo aproximado
- ≈ 55 acciones de herramienta, dentro del tope de 60.
- ≈ 130 KB leídos: por partes y con `grep`; U1 solo §4 y extractos.
- ≈ 70 KB de contenido nuevo: el documento, 2 bloqueos, este reporte y +45 líneas del script. Con las reescrituras completas que el generador hace de C3 (3 escrituras de ≈ 44 KB), el total escrito en disco ronda 180 KB.
- Sin paquetes instalados y sin red fuera de `finanzas.local` y `127.0.0.1`.
