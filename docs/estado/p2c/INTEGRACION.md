# P2c · Integración · Aplicación de la puerta 2 (11-oct-2026, noche)

Integrador: `integrador-p2c` (Opus). Fuentes que mandan: `docs/estado/puerta2/RESPUESTAS.md` (modo «lo recomendado + excepciones»; Q2-11 = B, Q2-64 = A, Q2-96 = A; puerta 2 bis: «ninguna» excepción, las «sin recomendación» quedan PROVISIONAL, K53) y `docs/estado/puerta2/PUERTA-2.md` (Q2-01 a Q2-98). Cada agente dejó su reporte en esta carpeta; este documento los consolida. HECHO = comprobado en disco o con un comando; PREGUNTA = falta una decisión; PROPUESTA = sugerencia sin decidir.

## 0. Resumen

- **HECHO.** `make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` termina en 0 con todo el árbol: WF-0, C1 1.3 (21 archivos idénticos al generador, 33 comprobaciones), S1 (esquema 1.4, `RESUMEN TODO OK`), U1, criterios (178 comprobadas, 0 diferencias), cobertura K25 y oráculos S2 a S5 con 0 diferencias. El alto que compartían las revisiones r3 de bloqueos, S1 y contratos («gen2: tipo de parámetro desconocido: hora HH:MM») ya no se reproduce: lo cerraron las correcciones de C1 (r1 y r2), posteriores a esas revisiones.
- **HECHO.** El paquete para Claude Design **no se pudo regenerar**: `bash docs/design/tools/paquete_entrada.sh` sale con `ERROR: U4-ajustes.pdf trae 2 líneas con /api/v1` (§2.2). Restauré `docs/design/entrada/` y los sellos a la versión versionada, así que el paquete sigue siendo el de WF-1 y `sha256sum -c` da 21 OK. Queda desfasado respecto de U3 1.5, U4 1.6 y las 4 líneas nuevas de DESIGN.md.
- **Revisiones (3 rondas como máximo).** Aprobados: C1 (r3), U4/U3 (r3), motor S2 (r2) y criterios (r3). **No aprobados tras 3 rondas:** bloqueos, S1, contratos de motores y motor S5 (§4). **S5-terminos: no verificado** (las 3 portadas devolvieron un desafío anti-robot) y **se escala a Axel**; Axel ya contestó en la puerta 2 bis («Ignóralo»), y esa respuesta se registra como K53 en P2d (§4.5).
- **HECHO.** Bloqueos al cierre: 84 archivos, **63 resueltos y 21 abiertos** (6 de WF-1, 13 nuevos de P2c y 2 de la Pista 3). `INDICE.md` ya los lista todos (§6).
- **Base del commit.** HEAD no es 1ccb468: la Pista 3 hizo 4 commits encima, 7c59d82 y 4079cee [P3a], 873f348 [P3b] y 0271f3b [P3c]. El último lo hizo a las 01:40, mientras yo integraba. El commit de la puerta 2 va sobre 0271f3b y **no incluye** archivos de la Pista 3 (§10).

## 1. Hecho por artefacto

| Artefacto | Versión (HEAD → ahora) | Autor (reporte) | Revisión | Qué cambió |
|---|---|---|---|---|
| Esquema S1 | 1.3 → **1.4** (`backend/db/migrations/00004_puerta2.sql`; diccionario 1.4; reglas de negocio 1.4) | `S1-puerta2`, `correccion-S1-r1`, `correccion-S1-r2` | **no aprobado** (r3: 1 alto, hoy cerrado de hecho; 7 bajos) | Q2-11 B (K47): 10 claves del hogar en `app_setting` (5 `reports.*`, 3 `alerts.*`, `optimizer.date_window_days` = 7, `ocr.timeout_seconds` = 30) con los valores de K50 y la decisión F; `params_json` de los 3 reportes queda `{}` y el Down lo repone; `credit.vat_included` (Q2-67, opción 1). 55 tablas STRICT, 152 índices, 79 disparadores; 00001 a 00003 sin cambios en su DDL |
| Contrato C1 | v1.2 → **v1.3** (`info.version` 1.3.0) | `C1-puerta2`, `correccion-C1-r1`, `correccion-C1-r2` | **aprobado** (r3: 6 bajos) | 168 rutas, 212 operaciones y 255 esquemas (antes 167, 210 y 252). Nuevas `getReportSettings` y `updateReportSettings` (`/api/v1/settings/reports`, N-296 y N-297); `AlertSettings` suma 3 claves; tipo «hora HH:MM» en el generador; `params_json` sin `send_time` ni `advance_days`. Todo sale de `backend/internal/api/c1gen/` (`build.py --check`: idéntico) |
| Contratos de motores | S2 1.2 → **1.3**; S3 1.4 → **1.5** (solo marcas); S4 1.4 → **1.5**; S5 1.4 → **1.5** | `contratos-puerta2`, `correccion-contratos-r1`, `correccion-contratos-r2` | **no aprobado** (r3: 1 alto, hoy cerrado de hecho; 1 medio; 5 bajos) | Q2-01 a Q2-08 aplicadas; S2 §2.12 ventanas contiguas (regla 4 PROVISIONAL, `S2-ventanas` pregunta 4); S4 cotización por elección (`choice_store`, `choice_presentation`); S5 una sola «mejor tienda», granel PROVISIONAL, PROFECO como regla técnica; marcas de 20 bloqueos resueltos pasan a FIRME |
| Motor S2 y oráculo | contrato 1.3 | `motor-S2-puerta2`, `correccion-motor-S2-r1` | **aprobado** (r2: 0 altos, 0 medios) | Ventanas con hueco = `invalid_argument`; tipo `Instant` (exige «Z»); metas con 3 estados (Q2-18 A); `provisional` null en flujo y proyección. 32 fixtures, 3185 de 3185 hojas, 159 de 159 de `model.json` |
| Motor S5 y oráculo | contrato 1.5 | `motor-S5-puerta2`, `correccion-motor-S5-r1`, `correccion-motor-S5-r2` | **no aprobado** (r3: 0 altos, 1 medio, 6 bajos) | Una sola «mejor tienda» (Q2-05 A, Q2-08 A); granel con un solo redondeo (en uso, `S5-carritos-granel`); lectura PROFECO con la lista cerrada de §11; `created_at` con zona. 36 fixtures, 1359 hojas esperadas, 958 de 958 hojas oráculo contra motor |
| U3 Inicio · U4 Ajustes | U3 1.2 → **1.5**; U4 1.4 → **1.6** (parte a 1.4 → 1.6, parte b 1.3 → 1.5) | `U4-puerta2`, `correccion-U4-r1`, `correccion-U4-r2` | **aprobado** (r3: 5 bajos) | Regla general de vista previa (Q2-09 A); P-26.08 y P-26.09 con claves de `app_setting` (Q2-11 B) y valores de K50; Q2-66 B, Q2-67 opción 1, Q2-68 C, Q2-69 (a), Q2-70 (a), Q2-76 |
| DESIGN.md y tokens C2 | 1.0 (sin cambio de versión) | `U4-puerta2` | dentro de la revisión de U4 | Aprobados tal cual por Axel (Q2-96 A). Solo 4 líneas: se quitan las marcas de Q2-81, Q2-82 y Q2-84 (ratificaciones, K51). Tokens sin cambio |
| Criterios | 1.2 → **1.3** (INDICE 1.3; S1 1.2, S2 1.3, S4 1.3, S5 1.3, S7 1.3, C1 1.3, U3 1.3, U4 1.3) | `criterios-puerta2`, `correccion-criterios-r1`, `correccion-criterios-r2` | **aprobado** (r3: 4 bajos) | K31/K44 de PROVISIONAL a FIRME con clave y valor (Q2-72 a Q2-77); S7, C1-22 y U4-12 leen los avisos en `app_setting`. 261 criterios (201 alta y 60 media) en 18 archivos. **S5-07 sigue NO CUMPLE** |
| Bloqueos y decisiones | — | `bloqueos-respuestas`, `correccion-bloqueos-r1`, `correccion-bloqueos-r2` | **no aprobado** (r3: 1 alto, hoy cerrado de hecho; 4 bajos) | 27 bloqueos con su sección de la puerta 2 y la línea «Estado del bloqueo (puerta 2)»; `decisiones.md` §K46 a §K52; tabla de las 98 respuestas |
| Términos de las tiendas | `docs/estado/S5-terminos.md` | `S5-terminos` | sin revisión (0 rondas) | Q2-64 = A intentada: 3 de 12 solicitudes, las 3 portadas con desafío anti-robot; 3 veredictos «no verificado»; no se leyó ninguna página legal |
| Integración | Makefile, `INDICE.md`, README | `integrador-p2c` | — | Makefile: 4 rótulos de `verify-s1` pasan a «esquema 1.4» y «00001 a 00004» (el comodín ya tomaba la 00004). `INDICE.md`: conteo final y alta de los 9 bloqueos de P2c que faltaban. README: estado de la puerta 2 |

## 2. Verificaciones del integrador

### 2.1 Verificación global

`make verify-all GOOSE=/home/deepsec/finanzas-hogar-orq/gobin/goose` (bitácora completa en `/tmp/p2c-verify-all.log`, fuera del repositorio): **exit 0** en 21.7 s. Líneas clave:

- `Up aplicado con sqlite3 3.53.4: backend/db/migrations/00004_puerta2.sql` · `tablas=55 (STRICT=55) · índices=152 · disparadores=79` · `Down en orden inverso (00004,00003,00002,00001) deja 0 objetos` · `goose v3.28.0: up, down-to 0 y up otra vez sin errores`.
- C1: `RESULTADO: los 21 archivos coinciden byte a byte con el generador` · `comprobaciones OK: 33` · `RESULTADO: TODO OK`.
- S1: `KUESKI OK … esquema 1.4`, `K44 OK`, `PROG14 OK`, `K47 OK`, `DOWN14 OK`, `GOOSE OK`, `RESUMEN TODO OK` · `s1_verifica.py aplicó las 4 migraciones del comodín`.
- U1: `comprobaciones OK: 63 · INFO: 5 · PENDIENTE: 0`.
- Criterios: `TOTAL 261 (alta 201 / media 60 / baja 0)` · `TOTAL comprobadas: 178 | con diferencias: 0` · `INDICE contra conteo propio: 0 diferencias`.
- K25: `motor_cred` 93.7 %, `motor_fin` 89.2 %, `precios` 96.0 %, `recetas` 93.2 %.
- Oráculos: S2 32 fixtures (3185 de 3185 hojas, 159 de 159 de `model.json`); S3 906 de 906 (437 de 437); S4 26 fixtures (8451 de 8451); S5 36 fixtures (958 de 958); `pre_verify.py: 15 comprobaciones OK · TODO COINCIDE`; `model.json` y `verify3.json` intactos.
- Recetario, tokens C2, contraste (0 bajo umbral), U4 (14 de 14 secciones) y paquete (`sha256sum -c: 21 archivos sin diferencias`).
- `== Verificación ampliada OK (WF-0, C1, S1, U1, criterios y WF-1)`.

Después de mis cambios (Makefile, `INDICE.md`, README y estos reportes) la volví a correr: el resultado está en `integrador-p2c.md`.

### 2.2 Paquete para Claude Design

- `bash docs/design/tools/paquete_entrada.sh` → **exit 1**. Salida:
  ```
  1/5 PDF y DOCX
    DESIGN: camino html+chromium
    U3-dashboard: camino html+chromium
    U4-ajustes: camino html+chromium
    C1-resumen: sin cambios en las fuentes
    DESIGN.docx generado
  2/5 Copias tal cual
  3/5 Páginas del reporte previo
    reporte-p01..p11: sin cambios en las fuentes
  4/5 Revisión de contenido (sin código)
  ERROR: U4-ajustes.pdf trae 2 líneas con /api/v1
  ```
- **Causa (HECHO).** U4 1.6 cita la ruta sin método HTTP en 2 sitios: `docs/design/U4-ajustes.md` l. 490 («(`/api/v1/settings/reports`, N-296 y N-297)», que viene de `U4-ajustes-b.md` l. 35) y l. 848 («(`/api/v1/settings/reports`)», de `U4-ajustes-b.md` l. 393). El filtro `filtrar_api` del script (decisión 42, M-6) solo quita la forma «(GET `/api/v1/…`)», así que `pdftotext` las encuentra. U3 1.5 trae 13 rutas, todas con método, y el filtro las quita.
- **Qué hice.** El script se detuvo en el paso 4 con `DESIGN.pdf`, `DESIGN.docx`, `U3-dashboard.pdf` y `U4-ajustes.pdf` nuevos y `MANIFIESTO.md` viejo. Así `verify-wf1-paquete` habría fallado. Restauré `docs/design/entrada/` y `docs/design/tools/paquete/.cache/` con `git checkout` (deshice mi corrida; no edité nada a mano). Después: `sha256sum -c` contra `MANIFIESTO.md` da 21 OK.
- **No lo arreglé:** U4 es un artefacto aprobado y el script no está en mi ALCANCE. Hay dos arreglos posibles (§5, I-1).

### 2.3 Otras comprobaciones

- Estado de los 84 bloqueos con la regla de conteo de la puerta 2 (última línea «Estado del bloqueo (puerta 2)» o «- **Estado:**»; script propio en `/tmp`): `total 84 · RESUELTO 63 · ABIERTO 21`; 13 archivos sin versionar, los 13 de P2c. Tres de ellos (`C1-hora-reportes-job`, `C1-reparto-alcance` y `S5-precio-vigente-captura`) no traen línea de estado. Por su texto («Sin decisión de Axel no se aplica»; «PROVISIONAL, en uso: el motor») están abiertos.
- `grep` del enum de `SavingsGoal.status` en C1 1.3: `on_track`, `behind` y `achieved` en `extractos/inicio.yaml`, `ajustes-b.yaml` y `presupuesto-b.yaml` (esquema `SavingsGoal`). `C1-datos.md` no cita Q2-18 (§5, I-2).
- Hoja §8 de `PUERTA-2.md`: 98 filas; 16 con «sin recomendación» (Q2-11 y Q2-64, que Axel contestó de forma explícita; Q2-48, Q2-50 a Q2-52 y Q2-56 a Q2-58, «en uso A»; Q2-85, Q2-87 a Q2-91 y Q2-94), y 3 con una parte sin recomendación (Q2-01 caso 2, la variante de Q2-07 y Q2-12 «C o A»).

## 3. Corte y continuación

El run se cortó a las 22:35 por un reinicio de la sesión. `bloqueos-respuestas` se reanudó sobre sus parciales; S1 y contratos arrancaron de cero; `S5-terminos` quedó «no verificado» por el desafío anti-robot. Las rondas que siguieron están en los reportes `correccion-*` y `REVISION-*` de esta carpeta.

## 4. Artefactos no aprobados tras 3 rondas (para que el orquestador los escale)

El alto común a 4.1, 4.2 y 4.3 era la verificación global rota en `verify-c1` («gen2: tipo de parámetro desconocido: hora HH:MM»; luego FALLA K44 y C1-15). **Hoy pasa** (§2.1): lo cerraron `correccion-C1-r1` y `r2` (00:19 y 00:36), después de esas revisiones r3 (23:10 a 23:46). No hubo ronda 4 que lo confirme. Esto es lo que queda abierto según cada revisión:

### 4.1 Bloqueos (`REVISION-bloqueos-r3.md`)
- H-r3-01 (alto): cerrado de hecho (ver arriba).
- H-r3-02 (baja): la pregunta 4 de `S2-ventanas` se describe distinto según el lugar.
- H-r3-03 (baja): cifras sin signo de pesos.
- H-r3-04 (baja): `RESPUESTAS.md` no lista Q2-12 como abierta (le toca al orquestador).
- H-r3-05 (baja): la tabla de la puerta 2 del INDICE no sigue el formato de las puertas 1 y 1b.

### 4.2 S1 (`REVISION-S1-r3.md`)
- H-01 (alto): cerrado de hecho.
- B-01 a B-07 (bajos): celda de `credit.vat_on_interest_bp` con la frase de 1.3; sección «propuesta S1» del OCR; marcas de la 00003 que citan un bloqueo resuelto; herramientas editadas fuera del ALCANCE (B-04, ver §5); `S1-puerta2.md` cita el CHECK de la ronda 0; INDICE sin 2 bloqueos (B-06, **cerrado por este integrador**); huecos de `s1_verifica.py` en S1-16, K47 y S1-13 (B-07).

### 4.3 Contratos de motores (`REVISION-contratos-r3.md`)
- A-1 (alto): cerrado de hecho.
- M-1 (media): el registro de bloqueos no tenía los bloqueos de las correcciones. **En `INDICE.md` queda cerrado** (§6). Falta en `decisiones.md` §K46, fuera de mi ALCANCE (§9).
- B-1 a B-5 (bajos): `correccion-contratos-r1.md` no remite a r2; `S2-ventanas.md` l. 85 conserva el texto de «la pregunta se redacta»; bajos de S2 de la ronda 1 sin atender (p. ej. S2 l. 359 «27 fixtures» contra 30); «Cómo seguí» de `S3-iva-lectura` pregunta 2; bitácora en presente con bloqueos ya resueltos.

### 4.4 Motor S5 (`REVISION-motor-S5-r3.md`)
- **M-1 (media, abierto):** 6 textos de `created_at` en que un lado acepta y el otro da `invalid_date`. El motor acepta la hora de un dígito (`T7:00:00-06:00`) y la coma decimal; el oráculo acepta un salto de línea final y dígitos no ASCII; el año 0000 lo acepta el motor y no el oráculo. El bloqueo `S5-precio-vigente-captura` afirma que no hay divergencia. No mueve cifras canon. El revisor pide al integrador decidir entre una corrección r3 acotada o registrarlo como pendiente. **No abrí la corrección: no está en mi ALCANCE.** Lo registro como pendiente de WF-3 (§8) y lo escalo.
- B-1 (independencia, §15.3): **registro la excepción** que piden `correccion-motor-S5-r1` y `r2`. Las dos correcciones tocaron el motor y el oráculo. El agente declara que escribió primero el oráculo, desde el contrato y RFC 3339, y que leyó `eval.go` después. Es una excepción al reparto de §15.3, sin aprobación de Axel.
- B-2: faltan marcas PROVISIONAL en `compare.go` (`S5-carritos-granel`) y en `profeco/links.go` y `filter.go` (`S5-profeco-filtro`).
- B-3: la tabla de §15 no lista 7 fixtures nuevas.
- B-4: el INDICE no listaba `S5-motor-oraculo-reglas` ni `S5-precio-vigente-captura`. **Cerrado por este integrador.**
- B-5: la pregunta 5 de `S5-motor-oraculo-reglas` no trae la cifra medida.
- B-6: `backend/internal/precios/types.go:4` cita el contrato 1.4 (el vigente es 1.5).

### 4.5 S5-terminos (no verificado; se escala a Axel)
- **HECHO.** Las 3 portadas (walmart.com.mx, bodegaaurrera.com.mx y sams.com.mx) devolvieron un desafío anti-robot. No se saltó y no se leyó ninguna página legal. `docs/estado/S5-terminos.md` tiene 3 veredictos «no verificado». El bloqueo `S5-terminos` sigue **ABIERTO** en su archivo y S5-07 sigue **NO CUMPLE**.
- **Respuesta de Axel ya dada** (puerta 2 bis, `RESPUESTAS.md`): «Ignóralo». La aplicación no lee sitios de tiendas. Un agente externo (una sesión de Claude Code de Axel con Claude in Chrome) trae precios, productos y recetas por la API. Se registra como **K53 en P2d**, con el cierre de `S5-terminos`. Hasta entonces, las consecuencias (importación en C1 con credencial de servicio y el peso de esos precios en la comparativa) son PROVISIONAL o PREGUNTA de la puerta 4.

## 5. Hallazgos del integrador

- **I-1 (media; bloquea el paquete).** U4 1.6 cita `/api/v1/settings/reports` sin método HTTP en 2 sitios, y el filtro de `paquete_entrada.sh` no los quita (§2.2). Dos arreglos, cada uno fuera de mi ALCANCE: (a) el dueño de U4 escribe la ruta con su método, p. ej. «(GET `/api/v1/settings/reports`)», o la quita, como en U3; o (b) el dueño del script agrega a `filtrar_api` las formas «(`/api/v1/…`)» y «(`/api/v1/…`, …)». Después hay que correr el script y `make verify-all`.
- **I-2 (media; inconsistencia de contrato).** C1 1.3 conserva `SavingsGoal.status` = `on_track`, `behind` y `achieved`. Q2-18 A (tres estados: en curso, cumplida y vencida), S2 1.3, U3 y U4 usan `in_progress`, `achieved` y `overdue`. Lo anotó `REVISION-U4-r2` («fuera del alcance de U4») y la revisión de C1 no lo levantó. El arreglo va en el generador de C1.
- **I-3 (baja).** `docs/design/DESIGN.md` l. 933 («Avisos y frecuencias») sigue PROVISIONAL y cita `S1-parametros-avisos`, que ya está RESUELTO (Q2-11 B). Lo pidieron `U4-puerta2` y `REVISION-U4-r1` B-2. Le toca al dueño de DESIGN.
- **I-4 (baja).** `C1-hora-reportes-job.md` y `C1-reparto-alcance.md` no traen línea de estado («Estado: ABIERTO»). Lo anotó `REVISION-U4-r2`. En el INDICE ya cuentan como abiertos.
- **I-5 (observación).** `REVISION-motor-S5-r3` y `REVISION-criterios-r3` (notas al integrador) listan marcas PROVISIONAL viejas que citan bloqueos ya resueltos: `docs/contratos/motores/productos-formato.md` l. 67 (`S4-asignacion-mercado`), `S2-interfaz.md` l. 354 (`S2-flujo-quincenal` sin la pregunta), `backend/internal/precios/ocr/types.go` l. 13, `oraculo/pre_genera_fixtures.py` l. 254, 10 líneas en `backend/internal/recetas`, y `gen2_schemas.py` l. 364 y 374 (`S1-parametros-avisos`). No las verifiqué una por una; quedan para WF-3 (§8).
- **I-6 (observación, sin verificar por mí).** `revisor-U4-r3` deja una nota sobre C1 1.3 `getCashFlowFortnights`: `planned_savings_cents` es `required` y `balance_cents` se describe restando el ahorro. S2 1.3 §3.11 lo deja en `null` y sin restar, con `S2-ahorro-quincenal` abierto.
- **I-7 (PREGUNTA al orquestador).** `RESPUESTAS.md` (puerta 2 bis) habla de «70 recomendadas» y «28 sin recomendación». La hoja §8 tiene 16 filas «sin recomendación» y 3 parciales (§2.3), y K46 cuenta 15 sin respuesta. No cambia nada de lo aplicado: Axel no pidió excepciones y las «sin recomendación» quedan PROVISIONAL. Pero 70 y 28 no cuadran con la hoja. ¿De dónde sale 28? No ajusté ninguna cifra.
- **S1 B-04 (ratificación).** `S1-puerta2` editó `backend/db/tools/gen_dict.py` (+44) y `s1_verifica.py` (+110), fuera de su ALCANCE. **Las integro en el commit:** `verify-s1` las usa y da `RESUMEN TODO OK`, y el diccionario sale idéntico de `gen_dict.py`. Queda registrado como excepción de alcance.

## 6. Bloqueos

**HECHO. Al cierre: 84 archivos, 63 resueltos y 21 abiertos.**
- **Los 6 abiertos de WF-1** tienen respuestas parciales o «sin recomendación»: `S2-flujo-quincenal` (Q2-12 «C o A»), `S2-ventanas` (pregunta 4, casos 1 y 2), `S5-carrito-comparativa` (Q2-48), `S5-carritos-detalle` (Q2-50 a Q2-52), `S5-profeco-filtro` (Q2-56 a Q2-58) y `S5-terminos` (K53 en P2d).
- **Los 13 nuevos de P2c**, todos abiertos y ya listados en `INDICE.md` («Bloqueos nuevos de P2c»):
  - `S2-ahorro-quincenal`, `S2-metas-vencida`, `S4-presentacion-respaldo` y `S5-carritos-granel` (`contratos-puerta2`);
  - `C1-hora-reportes-job` y `C1-reparto-alcance` (`C1-puerta2`);
  - `S1-reparto-edicion-gasto` (`correccion-S1-r1`);
  - `S3-formula-tarjeta-pdf`, `S3-iva-lectura` y `S4-empates-sin-confirmar` (`correccion-contratos-r1`);
  - `S5-motor-oraculo-reglas` (`correccion-motor-S5-r1`);
  - `S5-precio-vigente-captura` (`correccion-motor-S5-r2`);
  - `U4-reparto-preseleccion` (`correccion-U4-r1`).
- **Los 2 de la Pista 3** (`P3-contraste-chip-hover` y `P3-huecos-componentes`) son de su dueño.

Los que pide el encargo:
- **`S2-ventanas`, pregunta 4:** sigue abierta. El caso 1 lleva la regla 1A como PROVISIONAL (motor y oráculo: lista con hueco = `invalid_argument`). El caso 2 (ventanas que no duran un mes) va sin recomendación.
- **`S4-presentacion-respaldo`:** abierto. Rige la regla del motor, PROVISIONAL: la de la primera receta, luego la típica y luego la primera del catálogo.
- **`S5-carritos-granel`:** abierto. Rige la regla del motor, un solo redondeo, PROVISIONAL. El oráculo ya está alineado (`motor-S5-puerta2`).

Preguntas que siguen abiertas «sin recomendación» y pasan a WF-3 o a la puerta 4, con la opción en uso:
- Q2-01 caso 2; la variante de Q2-07; Q2-12;
- Q2-48, Q2-50 a Q2-52 y Q2-56 a Q2-58 (en uso A);
- Q2-85, Q2-87 a Q2-91 y Q2-94 (ratificaciones: siguen como están).

Ninguna se decidió en P2c.

## 7. Propuestas para Axel (de los reportes; ninguna está aplicada)

1. `C1-hora-reportes-job`: A, mandan las claves `reports.*` y la API reescribe el cron del job (`C1-puerta2`).
2. `C1-reparto-alcance`: B, «también los ya capturados» abarca todos los gastos que usan la regla (`C1-puerta2`).
3. `U4-reparto-preseleccion`: A, «Solo a los gastos nuevos» preseleccionada (`U4-puerta2`; la corrección r1 no recomienda).
4. Confirmar que `ocr.timeout_seconds` va en Ajustes › «Datos y respaldos» (`S1-puerta2` y `U4-puerta2`).
5. Confirmar que los créditos que ya existían sin porcentaje de IVA pasan como «incluye IVA» y no como «falta capturar» (supuesto 3 de `S1-puerta2`).
6. Rangos de seguridad para el OCR, la anticipación de precios viejos y la ventana de fechas. Hoy solo se pide 0 o más (`S1-puerta2`).
7. Entregar los 2 PDF de estado de cuenta, de una tarjeta y de un crédito de tienda (`S3-formula-tarjeta-pdf`).
8. Confirmar o corregir los empates de nutrición (Q2-90, `S4-empates-sin-confirmar`).
9. Al registrar K53 en P2d: reformular S5-07 («la aplicación no lee sitios de tiendas») o retirarlo con firma de Axel. Agregar un criterio de C1 para la importación del agente externo cuando la API exista (`criterios-puerta2`).
10. Criterios propios para Q2-67, Q2-69 (a) y Q2-82, si Axel los quiere (`criterios-puerta2`).
11. Contestar por id las preguntas de §6.

## 8. Pendientes

### 8.1 Para WF-3
1. Los 21 bloqueos abiertos, con la opción en uso PROVISIONAL (§6).
2. Cambios que deben implementar los motores: las listas de `contratos-puerta2.md`, `correccion-contratos-r1.md` y `correccion-contratos-r2.md`. La revisión r3 comprobó que cubren los 101 bloques del diff y los 106 renglones de código.
3. Motor S5: M-1, con gramáticas estrictas en ASCII en los dos lados y fixtures `s5-ingesta-*`; corregir el texto del bloqueo y sumar el año 0000 a la pregunta 3. También B-2, B-3, B-5 y B-6 (§4.4).
4. C1: `SavingsGoal.status` según Q2-18 A (I-2) y la nota de `getCashFlowFortnights` (I-6). También el «hora HH:MM» viejo en `gen2_schemas.py` l. 364 y 374.
5. Marcas PROVISIONAL viejas (I-5) y los bajos de §4.1 a §4.3.
6. Heredados de WF-1 (`docs/estado/wf1/INTEGRACION.md` §8) que siguen vigentes: medir el OCR con las 5 fotos (S5-06, K24); `backend/internal/dominio` con 85 % (K25); I-1 e I-2 de WF-1; dónde vive la base de PROFECO (`S5-profeco-almacen`); cola sin conexión (K9).
7. K53 (P2d): la importación en C1 para el agente externo, con credencial de servicio; S5-07 deja de aplicar.

### 8.2 Para la Pista 3
1. El paquete de `docs/design/entrada/` sigue en la versión de WF-1 (U3 1.2, U4 1.4) hasta que se arregle I-1. Para las maquetas, la fuente vigente es `docs/design/U3-dashboard.md` 1.5 y `docs/design/U4-ajustes.md` 1.6, no los PDF del paquete.
2. Las metas de ahorro usan los estados de Q2-18 A (`in_progress`, `achieved` y `overdue`), como U3 y U4, aunque C1 1.3 diga otra cosa (I-2).
3. P-26.12 con la preselección «Solo a los gastos nuevos», PROVISIONAL (`U4-reparto-preseleccion`).
4. Q2-98 A: nombres reales y cifras sintéticas. Q2-95 A: WF-3 no espera a la Pista 3.
5. Sus 2 bloqueos propios (`P3-contraste-chip-hover` y `P3-huecos-componentes`).

## 9. Solicitudes fuera del alcance del integrador

1. **Dueño de U4 o de `paquete_entrada.sh`:** arreglar I-1 y regenerar el paquete con el script.
2. **Agente de C1 (generador):** I-2 y las líneas 364 y 374 de `gen2_schemas.py`.
3. **Orquestador o P2d:** registrar en `decisiones.md` §K46 los 13 bloqueos nuevos de P2c (M-1 de contratos r3). Registrar K53 y cerrar `S5-terminos`. Agregar Q2-12 a `RESPUESTAS.md` como abierta (H-r3-04). Contestar I-7.
4. **Dueño de DESIGN.md:** I-3.
5. **Dueño de C1 (bloqueos):** I-4, la línea de estado de los 2 bloqueos de C1.
6. **Orquestador:** decidir si abre una corrección r3 acotada del motor S5 (M-1) o lo deja para WF-3. Los no aprobados de §4 se escalan.

## 10. Commit

- **Comando:** `git add -A -- . ':!docs/design/sistema' ':!docs/design/maquetas' ':!docs/estado/wfp3' ':!docs/design/salida' ':!docs/design/tools/ds'`. Antes comprobé que el índice estaba vacío, para no arrastrar nada que otro agente hubiera preparado.
- **Por qué excluyo también esas rutas.** Son de la Pista 3, que sigue en curso:
  - `docs/design/salida/iter-N/` son las copias de cada vuelta (P3-06; README, «Cómo sigue» 5): 210 archivos sin versionar, de 01:25.
  - `docs/design/tools/ds/verificar_sinteticos.py` (modificado a las 01:24) verifica `docs/design/maquetas/datos-sinteticos.json` (inventario y quincena en curso).
  - Ningún reporte de P2c los menciona.
  - Lo confirmó la propia Pista 3: su commit 0271f3b [P3c] (01:40) versionó `salida/iter-1` a `iter-3` y `tools/ds`, y no tocó ninguna ruta fuera de las suyas (`git show --name-only`). Después siguió escribiendo: `docs/design/salida/iter-4/` sin versionar y `tools/ds/generar_sinteticos.py` y `verificar_sinteticos.py` modificados. Por eso excluyo `docs/design/salida` y `docs/design/tools/ds` completos (P2c no escribió en ellas).
  - Meterlos en el commit «[Puerta 2]» atribuiría a P2c trabajo a medias de la Pista 3. Siguen en el árbol de trabajo, sin tocar.
- `docs/estado/S5-terminos.md` sí entra: es el resultado de Q2-64.
- El hash y `git show --stat` van en `integrador-p2c.md` y en el reporte al orquestador.

## 11. Acciones y KB por agente (de la sección «Consumo aproximado» de cada reporte)

| Agente | Acciones | KB leídos | KB escritos |
|---|---|---|---|
| bloqueos-respuestas | 19 | ~190 | ~7 |
| revisor-bloqueos-r1 | 33 | ~150 | ~14 |
| correccion-bloqueos-r1 | 14 | ~45 | ~9 |
| revisor-bloqueos-r2 | ~31 | ~120 | ~19 |
| correccion-bloqueos-r2 | 27 | ~130 | ~17 |
| revisor-bloqueos-r3 | 32 | ~120 | ~21 |
| S5-terminos | 8 | ~12 | ~7 |
| S1-puerta2 | 46 | ~260 | ~75 (más 190 del diccionario que reescribe el generador) |
| revisor-S1-r1 | 31 | ~165 | ~22 |
| correccion-S1-r1 | ~42 | ~120 | ~27 |
| revisor-S1-r2 | ~30 | ~150 | ~32 |
| correccion-S1-r2 | ~36 | ~115 | ~31 |
| revisor-S1-r3 | ~40 | ~150 | ~37 |
| contratos-puerta2 | 59 | ~300 | ~95 |
| revisor-contratos-r1 | 53 | ~190 | ~16 |
| correccion-contratos-r1 | 40 | ~190 | ~50 |
| revisor-contratos-r2 | 55 | ~170 | ~16 |
| correccion-contratos-r2 | 34 | ~125 | ~35 |
| revisor-contratos-r3 | 44 | ~150 | ~25 |
| C1-puerta2 | ~44 | ~210 | ~32 |
| revisor-C1-r1 | ~40 | ~200 | ~22 |
| correccion-C1-r1 | 38 | ~95 | ~22 |
| revisor-C1-r2 | ~50 | ~150 | ~21 |
| correccion-C1-r2 | 22 | ~75 | no consta en el extracto |
| revisor-C1-r3 | ~38 | ~115 | ~22 |
| U4-puerta2 | ~50 | ~250 | ~12 netos (más un script de 30 en /tmp) |
| revisor-U4-r1 | ~47 | ~210 | ~18 |
| correccion-U4-r1 | 49 | ~150 | ~30 |
| revisor-U4-r2 | ~38 | ~210 | ~19 |
| correccion-U4-r2 | 27 | ~75 | ~9 |
| revisor-U4-r3 | 39 | ~95 | no consta en el extracto |
| motor-S2-puerta2 | 46 | ~135 | ~45 (más ~70 generados) |
| revisor-motor-S2-r1 | 42 | ~175 | ~22 |
| correccion-motor-S2-r1 | ~16 | ~30 | ~5 |
| revisor-motor-S2-r2 | ~31 | ~150 | ~22 |
| motor-S5-puerta2 | ~50 | ~230 | ~150 (de ellos ~95 generados) |
| revisor-motor-S5-r1 | ~47 | ~180 | ~45 |
| correccion-motor-S5-r1 | ~33 | ~140 | ~30 |
| revisor-motor-S5-r2 | ~30 | ~120 | ~30 |
| correccion-motor-S5-r2 | 34 | ~150 | ~26 |
| revisor-motor-S5-r3 | 33 | ~175 | ~21 |
| criterios-puerta2 | ~58 | ~190 | ~75 |
| revisor-criterios-r1 | ~50 | ~200 | ~20 |
| correccion-criterios-r1 | 12 | ~75 | ~5 |
| revisor-criterios-r2 | ~46 | ~150 | ~19 |
| correccion-criterios-r2 | 8 | ~45 | ~5 |
| revisor-criterios-r3 | ~38 | ~120 | ~20 |
| integrador-p2c | ver `integrador-p2c.md` | | |

Ningún agente pasó el tope de 60 acciones. Los que pasaron de 150 KB escritos lo hicieron con archivos generados por scripts (fixtures, diccionario), no con texto a mano.
