# INTEGRACIÓN de WF-3b «Verificación visual y C3»

Integrador: `integrador-wf3b` (Opus) · 12-oct-2026 · Base: HEAD `6d691d0` «[WF-3a bis]» más el árbol de trabajo de WF-3b · Manda: `docs/criterios/U5.md` y `docs/criterios/C3.md` (comandos literales), `DESIGN.md`, U1, U3, U4, `docs/contratos/C3-integracion.md` y `docs/plan/decisiones.md` §A–§K. Los revisores evaluaron solo contra sus criterios; este documento consolida sus reportes y las corridas propias del integrador. No cambia cifras de los motores ni decisiones de Axel.

## 1. En corto

- **U5 NO APROBADO tras 3 rondas** (ronda 3 «devuelto», 0 altos y 2 medios). Se escala a Axel: `docs/estado/bloqueos/U5-ronda3-escalamiento.md`. De los 13 criterios, 11 se cumplen con evidencia reejecutada, U5-05 y U5-06 no, y U5-11 ya se cumple en el árbol actual sin revisor que lo confirme (§3).
- **C3 v1.2 SIN FIRMA tras 3 rondas** (ronda 3 «devuelto», 1 alto y 6 bajos). Se escala a Axel: `docs/estado/bloqueos/C3-ronda3-escalamiento.md`. Se cumplen 9 de 11 criterios. C3-08 no: pasan 7 de 8 flujos de U1 y falla F-3 ruta A con la semilla. Por eso C3-07 (firma) no se cumple (§4).
- **Evidencia versionada en este commit.** Son las 164 capturas reales (41 pantallas × Android 412×915 táctil y escritorio 1440×900, claro y oscuro) con su axe-core, contraste en el DOM, hosts y tema (`docs/design/capturas-u5/`, 505 archivos, 24 MB). También entran los flujos de extremo a extremo y la cola sin conexión contra el servidor real (`frontend/pruebas-u5/`) y C3 v1.2 con «capturado» 41 de 41.
- **Correcciones dentro del flujo.** `frontend-min-fix` cerró los 2 medios del frontend que dejó WF-3a bis (core R5-01 fechas de P-10, áreas R5-01 errores sin aviso en Ajustes) y la prueba inestable R5-05. `semilla-c3` aplicó las decisiones PROVISIONALES del orquestador sobre `WF3a-bis-semilla-dev` (puntos 4 y 7). `correccion-u5-r1`/`r2` y `correccion-c3-r1`/`r2` atendieron las rondas.
- **Verificación del integrador:** `make verify-wf3` rc 0: Playwright 241 passed, 4 skipped y 0 failed, y la prueba R5-05 de WF-3a bis ya pasa; incluye el nuevo `verify-u5`. `make verify-all` rc 0, después de restaurar un cambio de solo formato en `docs/design/entrada/prompt.md` (§2).
- **Bloqueos:** 8 nuevos, todos ABIERTOS, con 2 escalamientos de ronda 3. En disco hay **137 archivos de bloqueo** (§6).
- **Servidor de desarrollo** corriendo para la puerta 4 en `https://finanzas.local:8443` (y `https://127.0.0.1:8443`), base `datos/finanzas-dev.db` con la semilla.

## 2. Verificación global (corrida por este integrador, 12-oct-2026)

- **`make verify-wf3` → rc 0** («== WF-3: verificación completa OK (Go, cobertura K25, contrato C1, interfaz y evidencia U5)»):
  - Go: `go build`, `go vet`, `gofmt -l` y `go mod tidy -diff` limpios; `go test ./...` con 32 paquetes OK y 5 sin pruebas.
  - K25: dominio 87.0 %, motor_fin 89.2 %, motor_cred 94.0 %, recetas 93.2 % y precios 88.8 %. Sin umbral: api 82.0, auth 83.6, jobs 92.3, reportes 77.3 y telegram 80.0.
  - Contrato: `build.py --check` dice «los 23 archivos coinciden byte a byte»; S6-03 da 216 operaciones en 172 rutas y 0 diferencias, y pasan `TestRutasIgualesAlContrato` y `TestContratoSinSesionYEsqueleto`.
  - Interfaz: `npm run verify` con build, rutas, 0 hosts externos en `dist` y manifiesto válido. Playwright con la API simulada: **241 passed, 4 skipped, 0 failed** (5.2 min). La prueba R5-05 que dejaba `verify-wf3` en rojo en WF-3a bis ya pasa.
  - `verify-u5`: los 7 renglones de abajo.
- **`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` → rc 0 («== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)»: 55 tablas STRICT; goose v3.28.0 up, down-to 0 y up; C1 con 23 archivos idénticos al generador y `verify_c1` con 33 OK; S1 «RESUMEN TODO OK»; paquete con 21 archivos sin diferencias). Es la segunda corrida: la primera dio rc 2 en `verify-wf1-paquete` por `prompt.md` (ver abajo)**
- **Commit sin bases, certificados, llaves, contraseñas ni `node_modules`:** se probó con `git add -A --dry-run` y el pathspec del commit. Sobre los archivos candidatos: 0 rutas `.db`/`.sqlite`/`.pem`/`.key`/`.crt`/`.env`/`node_modules`/`tsbuildinfo`/`test-results`, 0 llaves privadas o tokens de bot por patrón y 0 archivos que contengan alguna de las contraseñas de desarrollo. Se comparó contra el archivo de fuera del repositorio sin imprimir su valor.
- **Corrección del integrador en el árbol (supuesto 1 de `integrador-wf3b.md`).** `docs/design/entrada/prompt.md` estaba modificado desde las 10:50 de hoy. Fuera de `MANIFIESTO.md` y sin autor en ningún reporte de WF-3b, solo cambiaba el formato: 10 renglones de una tabla Markdown con columnas alineadas, y `git diff --word-diff` no da ninguna palabra distinta. Ese cambio hacía fallar `verify-wf1-paquete` (`prompt.md: FAILED`) y con eso `verify-all`. Se restauró la versión del commit (`git checkout -- docs/design/entrada/prompt.md`) y el paquete vuelve a «21 archivos sin diferencias».

Comandos U5 sin servidor, corridos por el integrador: `make verify-u5` es nuevo y entra en `verify-wf3`.
- `python3 -I frontend/pruebas-u5/verificar_capturas.py` → «T = 41 · 4T = 164 · PNG correctos: 164 de 164 · faltantes: 0». El comando literal `ls docs/design/capturas-u5/ | wc -l` da 168: 164 PNG más `axe/`, `contraste/`, `evidencia/` y `hosts.json` (bloqueo `U5-conteo-carpeta`).
- Comando literal de U5-03 (`python3 -I -c "import json,glob; v=[…'serious','critical'…]; print(len(v))"`) → **0**, sobre 164 JSON de `axe/`.
- U5-04: 164 JSON de `contraste/`, **20,919 pares** evaluados (4,588 · 4,623 · 5,833 · 5,875 por combinación) y **0 bajo umbral**.
- U5-08: `hosts.json` → `["finanzas.local:8443"]` en 164 de 164 capturas.
- U5-11: 41 de 41 filas `| P-` de C3 v1.2 con «capturado» = sí y sus 4 rutas existentes.
- C3, comandos literales: `grep -c "^| P-"` = 41 (T de U1); `grep -c "| no |"` = 1 (fila Global, con bloqueo); `grep -c "^## Ronda" docs/estado/U5-revision.md` = 3, la última con «Veredicto: devuelto»; `ls docs/estado/bloqueos/U5-*` = 5 y `ls docs/estado/bloqueos/C3-*` = 4.

`make verify-u5-servidor` (nuevo, fuera de `verify-wf3` y de `verify-all`) regenera las capturas contra el servidor en marcha, en unos 9 minutos, y reescribe archivos versionados. El integrador **no** lo corrió, para no cambiar la evidencia que revisaron `revisor-u5-r3` y `revisor-c3-r3`: los 164 PNG tienen los mismos SHA-256 que en la revisión C3 r3. Su última corrida completa es la de `revisor-u5-r3`: 174 passed en 8.8 min.

## 3. U5 · criterios (docs/criterios/U5.md)

Fuente principal: `revisor-u5-r3` (ronda 3, con comandos reejecutados). En las filas marcadas «integrador» hay además una corrida propia (§2).

| Criterio | Estado | Evidencia |
|---|---|---|
| U5-01 (alta) | **Cumple** | `verificar_capturas.py`: 164 de 164 (integrador y revisor r3); `scrollWidth = innerWidth` en 164 de 164 (r3). El `ls \| wc -l` literal da 168, con el bloqueo `U5-conteo-carpeta` (A PROVISIONAL) |
| U5-02 (alta) | **Cumple** | r3 reejecutó el comando del capturista desde cero: `npx playwright test -c playwright.u5.config.ts` → 174 passed (8.8 min). Las cifras visibles $10,000.00, $6,905.50, $973.75, $6,293.44, $3,094.50 y $1,191.50 coinciden con `model.json` |
| U5-03 (alta) | **Cumple** | Comando literal = 0 (integrador). Con el axe propio de r3 (wcag22a): 0 serias o críticas en las 4 combinaciones |
| U5-04 (alta) | **Cumple** | 20,919 pares y 0 bajo umbral (integrador). r3 recalculó 20 de 20 pares con Δ máx. 0.0005. La regla de «componente de UI» sigue PROVISIONAL (`U5-contraste-ui-controles`, A) |
| U5-05 (alta) | **No cumple** | 3 rondas en `docs/estado/U5-revision.md`; la última dice «Veredicto: devuelto». No hay ronda 4: escalado (`U5-ronda3-escalamiento`) |
| U5-06 (alta) | **No cumple** | 13 de 14 elementos. Falla el #10: la barra lateral no tiene iconos de sección (M-7, Ec y Eo, DESIGN §7.9 y §5.4) |
| U5-07 (alta) | **Cumple** | r3: 351 importes por combinación y 0 fuera de formato, 0 fechas ISO, 0 palabras de la lista y 0 claves internas. Hay una baja B-10 («Versión dev») y la lista de nombres propios sigue PROVISIONAL (`U5-nombres-propios`) |
| U5-08 (alta) | **Cumple** | 1 host (`finanzas.local:8443`) en 164 cargas (integrador, capturista y r3) |
| U5-09 (alta) | **Cumple** | r3: 3 toques («+», Despensa, Guardar), la tarjeta 1 de Inicio cambia y «$187.25» aparece en P-02. `cola.spec.ts` dentro de `verify-wf3`: «toques para registrar un gasto (sin conexión): 3» |
| U5-10 (media) | **Cumple** | Fondo del body = C2 (`#FAFAFB` / `#222325`) en 164 de 164, sin `data-theme`; `tema.spec` pasa |
| U5-11 (media) | **Cumple en el árbol actual (sin revisor)** | r3 lo marcó como falla (M-8) porque leyó C3 v0.95. Desde v1.0 (`c3-matriz`), C3 marca «sí» con las 4 rutas: 41 de 41 (integrador, `make verify-u5`; también C3-05 en `revisor-c3-r3`). Si Axel elige A, lo confirma la verificación acotada |
| U5-12 (media) | **Cumple** | 5 bloqueos `U5-*` (`conteo-carpeta`, `contraste-ui-controles`, `nombres-propios`, `tipografia-tres-tamanos` y `ronda3-escalamiento`). Hay §Supuestos en `u5-capturas.md` y en los reportes de los revisores. Nada se marcó «aprobado» con violaciones abiertas |
| U5-13 (alta) | **Cumple** | r3, en un servidor aislado (`127.0.0.1:8445`), cumplió los 7 pasos: franja y contador 1 → 2, que se mantienen tras recargar; claves de 36 caracteres y `captured_at` ISO con desplazamiento; gastos 6 → 9 y pesos 0 → 1. El reenvío da 200 con `Idempotent-Replayed: true` y el mismo id. Un 500 deja «No se pudo enviar · 1 en cola» y 0 registros. La base de desarrollo no se tocó |

## 4. C3 · criterios (docs/criterios/C3.md 1.1) sobre C3 v1.2

Fuente principal: `revisor-c3-r3` (ronda 3, independiente, con ejecuciones propias). Hay además los greps del integrador (§2).

| Criterio | Estado | Evidencia |
|---|---|---|
| C3-01 (alta) | **Cumple** | 41 filas `P-`, igual a T de U1 (27 + 14), más la fila Global; 9 columnas; §1 define implementado, probado y capturado |
| C3-02 (alta) | **Cumple** | Script propio de r3: 0 pantallas sin fila; 211 `N-xxx` de U1 y las mismas 211 en §3.29; 216 tripletas iguales en `openapi.yaml` y `routes.go`; 5 operaciones sin pantalla, todas en la lista de excepciones |
| C3-03 (alta) | **Cumple** | 0 rutas del frontend fuera de `openapi.yaml` y 0 sin registrar en `routes.go`. `npm run build` rc 0. `listPeriods` es la única operación con pantalla y sin llamada (excepción PROVISIONAL, R3-7) |
| C3-04 (alta) | **Cumple** | `go test ./... -v`: 1126 PASS y 216 de 216 `TestContrato_Operaciones`. Playwright: 241 pasan, 4 omitidas, 0 fallan. Los 60 `archivo:línea` citados existen y pasan |
| C3-05 (alta) | **Cumple** | 41 filas con «sí» y 164 PNG válidos, con las 4 rutas exactas por fila (r3 e integrador) |
| C3-06 (alta) | **Cumple** | `grep -c "\| no \|"` = 1: la fila Global, que cita el bloqueo `C3-capturas-global` |
| C3-07 (alta) | **No cumple: sin firma** | r3 no agrega el bloque de firma porque C3-08 falla. Escalado (`C3-ronda3-escalamiento`) |
| C3-08 (alta) | **No cumple** | 7 de 8 flujos de U1 contra el servidor real, en copias aisladas: F-1 (3 toques), F-2, F-3 ruta B (cifra de K46), F-4, F-5 con reserva, F-6, F-7 y F-8 pasan. **F-3 ruta A falla** con la semilla (`flujos.spec.ts:122`), porque la lista del plan no tiene precios de ingredientes. Los 4 obligatorios de U1-05 (F-1, F-2, F-7 y F-8) pasan |
| C3-09 (alta) | **Cumple** | `c3-09.spec.ts` contra `finanzas.local:8443`: 111 cifras en pantalla, 0 diferencias con `model.json` |
| C3-10 (media) | **Cumple** | es-MX y fechas de evidencia en ISO (r3 leyó 11 filas y §3 a §8) |
| C3-11 (media) | **Cumple** | Versión v1.2 con historial (líneas 3 y 163 a 165); 4 bloqueos `C3-*`; §Supuestos y §Bloqueos en los reportes de `c3-matriz` y de las correcciones |

## 5. Rondas de revisión

| Artefacto | Ronda | Revisor | Veredicto | Lo atendió |
|---|---|---|---|---|
| U5 | 1 | `revisor-u5-r1` | devuelto: 5 altos, 6 medios, 6 bajos | `correccion-u5-r1` (axe de `BarraAvance`, P-26.13 a 415 px, alternados de vista, `SelectorArchivo`, sesión sin conexión, h1 de 32 px y barra de 240 px) |
| U5 | 2 | `revisor-u5-r2` | devuelto: 1 alto, 4 medios | `correccion-u5-r2` (nota de la semilla en P-26.13, puntero grueso en Android, sombra de las secciones de Ajustes y etiquetas de Tendencias) |
| U5 | 3 | `revisor-u5-r3` | **devuelto: 0 altos, 2 medios (M-7, M-8) y bajas** | Escalado a Axel |
| C3 | 1 | `revisor-c3-r1` | devuelto: 3 altos, 5 bajos (5 de 8 flujos) | `correccion-c3-r1` (pruebas de F-7 y F-8 contra el servidor real) |
| C3 | 2 | `revisor-c3-r2` | devuelto: 1 alto, 6 bajos (7 de 8) | `correccion-c3-r2` (F-3 ruta B contra la cifra de K46) |
| C3 | 3 | `revisor-c3-r3` | **devuelto y sin firma: 1 alto (R3-1), 6 bajos** | Escalado a Axel |

## 6. Artefactos no aprobados y hallazgos abiertos

**U5 (no aprobado tras 3 rondas)**, según `REVISION-u5-r3.md`:
- M-7 (media, U5-06 #10): la barra lateral de escritorio no tiene iconos de sección (`frontend/src/shell/BarraLateral.tsx`; el mapa de iconos está en `shell/Barras.tsx`).
- M-8 (media, U5-11): ya resuelto en C3 v1.2 (§3). Falta solo que lo confirme un revisor.
- Bajas: B-1 (destinos de la barra inferior), B-2 (`U5-nombres-propios`), B-3 (`U5-tipografia-tres-tamanos`), B-5 a B-8 (B-8: axe da `aria-prohibited-attr` incompleto en P-26.10), B-10 («Versión dev» en P-26.11, `sistema.Version` sin `-ldflags`), B-11 (sombra en la pestaña elegida de P-03 y P-22), B-12 (pista del interruptor de 40 × 24 y sombra de la perilla), B-13 (P-00 muestra «Tu sesión terminó» sin sesión previa), B-14 («Finanzas H…» truncado en la barra) y B-15 (hoja inferior sin `shadow.2`).
- Observación sin severidad: P-06 muestra en la quincena 2 (16 a 31 de octubre de 2026) el mismo real de la quincena 1, $1,197.64, aunque el único gasto es del 10 de octubre. Que lo revise el dueño de Presupuesto.

**C3 v1.2 (sin firma tras 3 rondas)**, según `REVISION-c3-r3.md`:
- R3-1 (alta, C3-08): F-3 ruta A falla con la semilla y depende del punto 2 de `C3-flujos`.
- Bajas: R3-2 (F-8 cita `:139` y la prueba está en `:141`), R3-3 (§5 cita `getDebtTrend` en lugar de `getLiveDebt`), R3-4 (lista de pruebas truncada con «…»), R3-5 (`flujos.spec.ts:70` sin `expect(reg.ok)`), R3-6 (F-5 con reserva) y R3-7 (`listPeriods` sin llamada).
- Observación de r3: después de F-3A3, P-13 muestra **dos** renglones «Bodega Aurrera · Vallejo» ($450.00 y $859.00).

**Frontend (no aprobado tras 5 rondas en WF-3a).** `frontend-min-fix` cerró en WF-3b los 2 medios que seguían abiertos (core R5-01 y áreas R5-01) y estabilizó R5-05. Lo comprobó con sus propias corridas (241 passed, 4 skipped, 0 failed, dos veces) y el integrador lo confirma con `make verify-wf3` (§2). Ningún revisor independiente de frontend volvió a revisar esos 3 cambios; solo los cubren de forma indirecta las rondas de U5. Los bajos de `INTEGRACION-wf3a.md` §15.4 siguen igual.

## 7. Bloqueos nuevos de WF-3b (8, todos ABIERTOS) y actualizados

| Archivo | Lo abrió | Tema |
|---|---|---|
| `U5-ronda3-escalamiento.md` | `revisor-u5-r3` | U5 devuelto en la ronda 3; opciones A (recomendada), B y C |
| `C3-ronda3-escalamiento.md` | `revisor-c3-r3` | C3 sin firma en la ronda 3; opciones A (recomendada), B y C |
| `C3-flujos.md` | `c3-matriz` | C3-08: punto 2 (F-3 ruta A) abierto, punto 3 (lo contesta K46), punto 4 (F-5) y observación sobre Q2-49 |
| `C3-capturas-global.md` | `c3-matriz` | «Capturado» de la fila Global (a PROVISIONAL «no») |
| `U5-conteo-carpeta.md` | `u5-capturas` | El `ls \| wc -l` literal da 168 (A PROVISIONAL: contar los PNG) |
| `U5-contraste-ui-controles.md` | `u5-capturas` | Qué cuenta como componente de UI en el 3:1 (A PROVISIONAL) |
| `U5-nombres-propios.md` | `revisor-u5-r1` | Claude, Ollama, Tesseract y las siglas API, CSV, PDF y OCR en pantalla |
| `U5-tipografia-tres-tamanos.md` | `revisor-u5-r1` | «Máximo tres tamaños» de DESIGN §3.2 (A PROVISIONAL: solo el área de contenido) |

Actualizado: `WF3a-bis-semilla-dev.md`, con las decisiones del orquestador (PROVISIONALES hasta la puerta 4) aplicadas por `semilla-c3`:
- Punto 4 = (a): gastos proyectados de la fixture S2. Con ellos, `getAccrualVsCash` P1 a P6 = `model.json` y la deuda viva al 10-abr-2027 = $6,293.44.
- Punto 7 = sí: 36 precios de Walmart en la semilla real. Con ellos, la comparativa de la compra del 10-oct da +$264.50 y no los +$138.50 de `data.json`; la diferencia se explica por el método.
Los puntos 1, 3 y 5 siguen como pregunta a Axel. Índice al día: `docs/estado/bloqueos/INDICE.md` §«Bloqueos nuevos de WF-3b», con 137 archivos en disco.

## 8. Pendientes para la puerta 4

1. **S6-11 en Chrome Android** con la CA de mkcert instalada: lo prueba **Axel** a mano en sus teléfonos (`ops/mkcert.md` §4). En la máquina se cumple (WF-3a).
2. **S7-07, envío real por Telegram**: lo hace **el orquestador**. El servidor de desarrollo corre sin token (HOME aislado) y ningún agente de WF-3b envió nada.
3. **Escalamientos.** `U5-ronda3-escalamiento` (A, B o C) y `C3-ronda3-escalamiento` (A, B o C), que incluye el punto 2 de `C3-flujos`. Con la A en los dos, basta una corrección acotada (iconos de la barra lateral) y una verificación acotada de U5-06 #10, U5-11 y C3-08 (F-3), y luego la firma.
4. **Preguntas abiertas:**
   - el punto 4 de `C3-flujos`: F-5 con oct-2026 y 6 toques;
   - la observación sobre Q2-49: la app arranca sin precios de Walmart, pero la semilla ahora los trae;
   - los puntos 1, 3 y 5 de `WF3a-bis-semilla-dev`;
   - `C3-capturas-global`;
   - los 4 bloqueos `U5-*` de reglas;
   - `C3-matriz` 1, 3c y 4a, que siguen PROVISIONALES desde WF-3a bis.
5. **Documentos de otro dueño:** U1 1.2 (§4 F-3 y §3.16) y U3 (tarjeta 5, «Ejemplo trazable») siguen diciendo Walmart $3,233.00 (+$138.50). Con K46 (Q2-47 = A), la pantalla y la API dan +$264.50. Los corrige su dueño por encargo del orquestador (`C3-flujos` punto 3).
6. **Valores K31 PROVISIONALES**: la lista consolidada está en §9.
7. **Del frontend:** las observaciones de §6 (P-06, P-13 duplicado y «Versión dev»). También queda la que imprime `pruebas/medidas.spec.ts` en `verify-wf3`: «PENDIENTE U5-07 (solicitud a front-inicio)», porque las gráficas de `/inicio?sub=tendencias` llevan un `aria-label` en inglés que genera la biblioteca de gráficas («This is a chart…»). No es texto visible y por eso no lo cuentan U5-07 ni r3, pero sí lo lee un lector de pantalla. Siguen además las solicitudes fuera de alcance que siguen de WF-3a bis (`INTEGRACION-wf3a.md` §15.8). La 4 (`frontend/tsconfig.tsbuildinfo` en `.gitignore`) queda aplicada en este commit.

## 9. Valores K31 PROVISIONALES (lista consolidada de WF-3a y WF-3b)

WF-3b **no agregó valores K31 nuevos**: `grep K31` en los 15 reportes de agentes de WF-3b da 0 coincidencias. Los 8 bloqueos nuevos tienen «—» en la columna K31. La lista sigue siendo la de `INTEGRACION-wf3a.md` §9, y WF-3b puso a prueba una parte contra el servidor real (U5-13: cola, reenvío idempotente y sondeo).

| Tema | Valor PROVISIONAL en el código | Bloqueo | Ejercido en WF-3b |
|---|---|---|---|
| Respaldos: hora | Diario a las 03:15 (America/Mexico_City) | `S6-respaldos`, `S6-valores-operativos-base` | — |
| Respaldos: retención | 30 diarios | `S6-respaldos` | — |
| Respaldos: restauración | Reemplaza todo salvo `backup` y `device_session`; trabajo síncrono | `S6-respaldos` (PREGUNTA) | — |
| Adjuntos: tipos | JPG, PNG, WebP y PDF por contenido (415 si no) | `S6-adjuntos` | F-8 adjunta el ticket (PNG) |
| Adjuntos: tamaño | 10 MB (1 MB = 1,048,576 bytes) | `S6-adjuntos` | — |
| Cabeceras y CSP | `nosniff`, `no-referrer`, `DENY`, `no-store` en la API; CSP `default-src 'self'` … `frame-ancestors 'none'` | `S6-valores-operativos-base`, `WF3-frontend-servidor` | U5-08: 0 hosts externos en 164 cargas |
| Resumen diario | Fecha, gastos de hoy, gastos del mes y pagos de los próximos 7 días | `S7-fuentes-avisos` | — |
| Resumen semanal | Semana ISO, gastos de la semana y del mes, pagos de 7 días y casillas sin armar | `S7-fuentes-avisos` | — |
| Reporte mensual | Resumen, devengado contra efectivo, categorías, quincenas, créditos y parte de Karen | `S7-reporte-mensual-secciones` | F-5 genera y descarga el PDF |
| Avisos | Revisión cada minuto; un envío fallido queda `failed` y no se reintenta solo | `S7-fuentes-avisos` §4 | — |
| Contraseñas | argon2id t=3, m=64 MiB, p=4 | `S6-valores-operativos-base` | — |
| Límites del cuerpo | 8 MiB por JSON (413); 32 MiB en memoria con Idempotency-Key | `S6-valores-operativos-base` | — |
| Purga de idempotencia | Cada hora, al minuto 5 | `S6-valores-operativos-base` | U5-13: el reenvío dentro de la ventana da `Idempotent-Replayed: true` |
| PWA sin conexión | Sondeo de `/health` cada 15 s; reintento de la cola cada 30 s; copia de lecturas hasta cerrar sesión; la cola no se borra al cerrar sesión | `WF3-frontend-k31` | U5-13 y `cola.spec.ts` |

Otras decisiones PROVISIONALES de WF-3b que no son K31 y también se confirman en la puerta 4:
- semilla: puntos 4 (a) y 7 (sí);
- `C3-capturas-global` (a);
- `U5-conteo-carpeta` A;
- `U5-contraste-ui-controles` A;
- `U5-tipografia-tres-tamanos` A.

## 10. Consumo por agente (aproximado, de cada reporte)

| Agente | Modelo o rol | Acciones | KB leídos | KB escritos |
|---|---|---|---|---|
| frontend-min-fix | corrección del frontend | ≈ 60 | ≈ 135 | ≈ 28 |
| semilla-c3 | Sonnet | ≈ 60 | ≈ 200 | ≈ 45 |
| u5-capturas | capturista | ≈ 38 | ≈ 120 | ≈ 45 (más ≈ 24 MB de PNG y JSON generados) |
| u5-flujos | flujos | ≈ 70 (**pasó el tope de 60**) | ≈ 150 | ≈ 75 (más 22 JPEG) |
| c3-matriz | Opus | ≈ 55 | ≈ 130 | ≈ 70 nuevos; ≈ 180 en disco contando las 3 reescrituras del generador (**pasó el tope de 150**) |
| revisor-u5-r1 | revisor | ≈ 47 | ≈ 115 (más 13 láminas JPEG) | ≈ 32 (más ≈ 38 temporales, borrados) |
| correccion-u5-r1 | corrección | ≈ 58 | ≈ 120 | ≈ 30 |
| revisor-u5-r2 | revisor | ≈ 60 | ≈ 150 (más 8 hojas) | ≈ 30 |
| correccion-u5-r2 | corrección | ≈ 52 | ≈ 120 | ≈ 22 |
| revisor-u5-r3 | revisor | ≈ 59 | ≈ 150 (más 7 hojas JPEG, ≈ 0.84 MB) | ≈ 30 |
| revisor-c3-r1 | revisor | 37 | ≈ 125 | ≈ 13 |
| correccion-c3-r1 | corrección | ≈ 52 | ≈ 150 | ≈ 40 |
| revisor-c3-r2 | revisor | ≈ 36 | ≈ 110 | ≈ 24 |
| correccion-c3-r2 | corrección | ≈ 47 | ≈ 150 | ≈ 22 |
| revisor-c3-r3 | revisor | ≈ 40 | ≈ 110 | ≈ 17 |
| integrador-wf3b | Opus | ≈ 58 | ≈ 130 | ≈ 60 |

Paquetes instalados en WF-3b: ninguno. Los reportes que lo dicen afirman que Playwright, `@axe-core/playwright`, ImageMagick y numpy ya estaban.

## 11. Propuestas para Axel (PROPUESTA, no decidido)

| Qué | Resultado | Costo | Riesgo |
|---|---|---|---|
| A en los dos escalamientos: iconos de la barra lateral; U1 aclara que la ruta A de F-3 pide una lista con todos los precios en al menos una tienda; verificación acotada de U5-06 #10, U5-11 y C3-08 (F-3) y firma de C3 | U5 aprobado y C3 firmado sin abrir otra ronda completa | Bajo: un componente, una aclaración en U1 y unos 15 min de verificación | Bajo. La tarjeta «Mejor canasta» sigue anunciando la compra del 10-oct y abriendo la semana, salvo que también se elija la (b) de la cara 2b |
| B: aceptar U5 y C3 con los pendientes con fecha | Se pasa a WF-4 ya | Ninguno ahora | La barra lateral queda distinta de DESIGN y la firma de C3 lleva una excepción explícita en C3-08 |
| C: otra ronda completa de U5 y de C3 | Revisión completa | Medio: unos 9 min de capturas y 2 revisiones | Rompe el límite de 3 rondas (U5-05; plan §3.1, principio 5) |

Recomendación de los dos revisores y de este integrador: **A**.

## 12. Commit

Un solo commit, con `git add -A -- . ':!docs/design/sistema' ':!docs/design/maquetas' ':!docs/estado/wfp3'`. Antes de hacerlo se comprobó que no entran bases `.db`, certificados, llaves, contraseñas ni `node_modules` (§2). El título que pidió el orquestador decía «… y C3 firmado». **No se usó tal cual porque sería falso: C3 queda sin firma** (§4, C3-07). Se usó «… y C3 v1.2 sin firma (U5 y C3 escalados a Axel tras la ronda 3)». Las líneas de atribución siguen la indicación vigente de la sesión.
