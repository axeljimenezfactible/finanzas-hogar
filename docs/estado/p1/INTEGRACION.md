# Integración del flujo P1 · Aplicar respuestas de la puerta 1 (11-oct-2026)

Integrador: `integrador-p1`. Fuente de las decisiones: `docs/plan/decisiones.md` §K (K1 a K34, respuestas de Axel del 11-oct-2026). Reglas: `docs/plan/protocolo-agentes.md` (§4, reparto fino: ≤ 400 KB leídos, ≤ 150 KB escritos y ≤ 60 acciones por agente).

## 1. Resultado de las revisiones

| Artefacto | Versión | Revisión | Veredicto | Informe final |
|---|---|---|---|---|
| Modelo de datos (S1) | esquema 1.2 (migración 00002) | 3 rondas | **Aprobado** · 0 altos, 0 medios, 9 bajos | `docs/estado/p1/REVISION-S1-r3.md` |
| Arquitectura de información (U1) | 1.1 | 3 rondas | **Aprobado** · 0 altos, 0 medios, 7 bajos | `docs/estado/p1/REVISION-U1-r3.md` |
| Colorimetría (U2a) | 1.1 | 2 rondas | **Aprobado** · 0 altos, 0 medios, 5 bajos | `docs/estado/p1/REVISION-U2a-r2.md` |
| Criterios de aceptación | 1.1 (18 archivos) | 3 rondas | **Aprobado** · 0 altos, 0 medios, 4 bajos | `docs/estado/p1/REVISION-criterios-r3.md` |
| Contrato C1 | v1.1 (`info.version` 1.1.0) con extractos | 3 rondas (tope) | **No aprobado** · 2 altos, 2 medios, 5 bajos | `docs/estado/p1/REVISION-C1-r3.md` |

C1 agotó sus 3 rondas sin aprobación. El integrador **no** corrige artefactos revisados: los 4 hallazgos altos y medios quedan en §5 y §7 para que el orquestador decida (ronda extra de corrección de C1 antes de WF-1, o pasarlos a WF-1 con su dueño).

## 2. Hecho por artefacto

**S1 · esquema 1.2** (`docs/estado/p1/S1-esquema-1.2.md`, correcciones r1 y r2)
- Migración `backend/db/migrations/00002_puerta1.sql` (Up y Down): borrador de OCR en `attachment` (K3), tabla `telegram_chat_candidate` (K4), `idempotency_key` y `captured_at` en `expense` y `weight_entry` (K9), decimales capturados de las tasas (K12), `rate_base` y pago mínimo en `credit` para tarjeta y crédito de tienda (K18), `health_profile.goal` (K23) y ajustes sembrados de K6, K11 y K21. El Down deja `sqlite_master` igual al esquema 1.1, byte por byte.
- Diccionario 1.2 y reglas de negocio (28 reglas) con K1–K23; comentarios de 00001 sin cambiar el esquema 1.1.
- `backend/db/tools/s1_verifica.py` (antes en el Anexo B del reporte de WF-0) y `gen_dict.py`.
- Resultado: 55 tablas STRICT, 152 índices, 71 disparadores; regresión Kueski igual al oráculo.

**C1 · v1.1** (`docs/estado/p1/C1-v1.1.md`, correcciones r1 y r2)
- Generador `backend/internal/api/c1gen/` actualizado: aplica todas las migraciones, aplica §K y genera **extractos por módulo** (`gen7_extractos.py` → `docs/contratos/extractos/`: 17 YAML más `INDICE.md`).
- `openapi.yaml` 1.1.0: 201 operaciones en 162 rutas; `routes.go` con 201 rutas; `C1-datos.md` con matriz de 227 filas y 203 necesidades. Incluye `scope` de presupuesto (K2), borrador de OCR (K3), candidatos de chat con `chat_id_masked` (K4), `family` de 3 valores (K5), valores técnicos de K6, `Idempotency-Key` y `captured_at` (K9), decimales de tasas (K12), tope por persona (K13) y campos de K18.

**U1 · 1.1** (`docs/estado/p1/U1-1.1.md`, correcciones r1 y r2)
- `docs/design/U1-arquitectura.md` 1.1 con «Cambios de la versión 1.1» (16 cambios): K2, K3, K4, K5, K8 (firme), K9 (N-317, cola sin conexión), K10, K13, K16 («Faltan restricciones»), K18 («Falta capturar»), K21 y K22; necesidades nuevas N-316 a N-319.
- Verificador movido a `docs/design/tools/verificar_u1.py` (antes anexo del reporte de WF-0).

**U2a · colorimetría 1.1** (`docs/estado/p1/U2a-1.1.md`, corrección r1)
- Gasto pasa a 45° (K26) con barrido nuevo de 40° a 45°; Karen y Bodega Aurrera reubicadas; la marca entra al barrido de pares cercanos (R3-04); K27 con ΔL ≥ 0.10 entre series consecutivas.
- `colorimetria.md`, `contraste.csv`, `C2-colores.json`, `paleta.html` y `paleta.png` regenerados por `docs/design/tools/colorimetria.py`: 190 tokens y 202 pares, 0 bajo umbral.

**Criterios · 1.1** (`docs/estado/p1/criterios-S-1.1.md`, `criterios-UC-1.1.md`, correcciones r1 y r2)
- Los 18 archivos pasan a la versión 1.1 con su lista de cambios; U1-15, U1-16, C1-16 a C1-21 y U5-13 son nuevos. Total: 257 criterios (197 alta y 60 media).
- `docs/criterios/tools/`: `verif_final.sh` y 8 scripts (conteo, traza, referencias, citas, cifras canon, S1-09 en SQL, equivalencia con el oráculo y matriz C1/U1).

**Integración** (este agente)
- `Makefile`: `verify-s1` llama a `backend/db/tools/s1_verifica.py` (ya no extrae el anexo); objetivos nuevos `verify-u1` y `verify-criterios`; `verify-all` los incluye; requisito Go 1.27.
- `README.md`: estado al 11-oct-2026, árbol con las carpetas nuevas y tabla de versiones vigentes.
- `docs/criterios/INDICE.md`: conteo propio y reconciliado (§10), cifra de §1 y fecha.
- `docs/estado/bloqueos/INDICE.md`: sección «Bloqueos nuevos de P1» con los 8 abiertos.

## 3. Verificaciones del integrador

`make verify-all GOOSE=<scratchpad>/gobin/goose` → código 0, unos 10 s. Salida resumida:
- **SQL:** Up de `00001_esquema.sql` y luego de `00002_puerta1.sql` (orden del wildcard ordenado); 55 tablas (55 STRICT), 152 índices, 71 disparadores, `foreign_key_check` 0 filas, `integrity_check` ok; **Down en orden inverso (00002 y luego 00001) deja 0 objetos**; Up otra vez: 55 tablas; goose v3.28.0 up, down-to 0 y up sin errores.
- **Go:** go1.27.2; `go build ./...`, `go vet ./...` y `gofmt -l .` sin errores ni archivos por formatear.
- **Colorimetría:** «VERIFICACIÓN OK», 202 filas de contraste con 0 `cumple=false`; `colorimetria.md`, `contraste.csv`, `C2-colores.json` y `paleta.html` idénticos byte a byte al regenerado.
- **OpenAPI:** `openapi-spec-validator` 0.9.0 → `openapi.yaml: OK`.
- **C1:** `build.py --check` → «los 21 archivos coinciden byte a byte con el generador»: `openapi.yaml`, `routes.go`, `C1-datos.md` y los **18 de `docs/contratos/extractos/`**. Así, `verify-c1` ya cubre los extractos. `verify_c1.py --db` → 24 comprobaciones OK, «TODO OK».
- **S1:** `s1_verifica.py` → 23 líneas OK (S1-01 a S1-17, KUESKI, K12, NEG, NEG12, DOWN y GOOSE) y «RESUMEN TODO OK».
- **U1:** `verificar_u1.py` → 53 comprobaciones OK, 9 INFO, 0 PENDIENTE, «TODO OK».
- **Criterios:** `verif_final.sh` → 18 archivos y 257 criterios; 0 sin línea en §5; 257 IDs citados y 0 inexistentes; 178 cifras canon sin diferencias; S1-09 en SQL bien; 12 de 12 cifras del calendario iguales al oráculo; 0 necesidades sin operación; 0 residuos «N-xx»; 0 rutas `/tmp`. Las 5 citas «SIN FUENTE» son informativas (hallazgo bajo B-04 de criterios r3).
- **Conteo propio de criterios** (script sobre `^\*\*<ID>-nn … (severidad)` en los 18 archivos): 257 (197 / 60 / 0), sin IDs repetidos y todos en la versión 1.1. Coincide con la tabla de `docs/criterios/INDICE.md` §2 y con `traza.py`. Ver §10 del índice.
- Prueba del objetivo nuevo: `make -n verify-all` sin errores y luego la corrida real. Las 13 cadenas esperadas de `verify-criterios` y los dos «0» (comprobaciones 5 y 10) se encontraron en la salida.
- Segunda corrida después de editar README y los dos índices: ver §3.1.

### 3.1 Segunda corrida y prueba negativa
- `make verify-all GOOSE=…/goose` después de editar `README.md`, `docs/criterios/INDICE.md` y `docs/estado/bloqueos/INDICE.md`: código 0, con las mismas cifras de §3 (21 archivos de C1 idénticos, 24 OK en `verify_c1.py`, «RESUMEN TODO OK» en S1, 53 OK en U1, 257 criterios y 0 diferencias). Al terminar no queda `datos/tmp`.
- **Prueba negativa de `verify-criterios`:** el Makefile acepta `CRITERIOS_SH=<script>` (por omisión, `docs/criterios/tools/verif_final.sh`). Con un script falso que devuelve la salida real alterada («con diferencias: 2» y la comprobación 10 en 3), `make verify-criterios` termina con error y muestra «FALLA, no aparece en la salida: | con diferencias: 0» y «FALLA, la comprobación 10 de verif_final.sh da «3» (se espera 0)». Con el script real termina en 0. `verif_final.sh` siempre sale con 0, así que el Makefile es el que decide si pasa.
- `make clean`: «Limpio: sin datos/tmp, sin backend/bin y sin __pycache__».

## 4. Marcas PROVISIONAL y bloqueos que citan (`grep -c PROVISIONAL`, 11-oct-2026)

| Archivo | Líneas | Bloqueos citados en la línea o en sus vecinas (estado) |
|---|---|---|
| `backend/db/migrations/00001_esquema.sql` | 0 | — |
| `backend/db/migrations/00002_puerta1.sql` | 0 | — |
| `docs/contratos/diccionario-datos.md` | 1 | Solo la leyenda (línea 19) |
| `docs/contratos/reglas-negocio.md` | 11 | `S1-fuentes-precio` (abierto), `criterios-frescura-extremos` (abierto), `S1-ajuste-objetivo-kcal` (abierto), `criterios-promedios-redondeo` (abierto); historial de `S1-frescura-peso` y `S1-meta-kcal` (resueltos); 1 línea es la leyenda (línea 9) |
| `docs/design/U1-arquitectura.md` | 7 | `U1-lista-vigente-domingo` (abierto); K31 con `criterios-umbrales-ausentes` (resuelto, definición de «casillas sin armar»); historial de `criterios-frescura-peso` y `criterios-invariante-deuda` (resueltos) |
| `docs/contratos/openapi.yaml` | 6 | `C1-valores-k31` (abierto), `S1-ajuste-objetivo-kcal` (abierto), `criterios-promedios-redondeo` (abierto) |
| `docs/contratos/C1-datos.md` | 18 | `C1-valores-k31` (abierto), `S1-ajuste-objetivo-kcal` (abierto), `criterios-promedios-redondeo` (abierto); historial de `C1-lista-compras` (resuelto); 5 líneas son leyenda, historial o el título de §13 |

Bloqueos abiertos citados: `C1-valores-k31`, `S1-ajuste-objetivo-kcal`, `S1-fuentes-precio`, `U1-lista-vigente-domingo`, `criterios-frescura-extremos` y `criterios-promedios-redondeo`. Los otros dos bloqueos nuevos (`U2a-rosas-ximena-bodega` y `U2a-walmart-oscuro`) marcan tokens en `C2-colores.json` y `colorimetria.md`, fuera de esta lista de archivos. **Hueco conocido:** `openapi.yaml` aplica la regla del domingo en `getMealPlanWeek` sin marca ni cita de `U1-lista-vigente-domingo` (hallazgo alto H-r3-02 de C1). Ninguna cifra del oráculo depende de una marca PROVISIONAL.

## 5. Hallazgos que quedaron registrados

**C1 (no aprobado; abiertos para decisión del orquestador)**
- H-r3-01 · alta · `getBasketComparison` lee la lista de compras sin la llave `week` (usa `period` opcional). Dueño: generador `backend/internal/api/c1gen/` (`gen3_ops.py`).
- H-r3-02 · alta · `getMealPlanWeek` aplica la regla del domingo sin marca PROVISIONAL ni cita de `U1-lista-vigente-domingo` (abierto). Dueño: generador de C1.
- H-r3-03 · media · `docs/design/U1-arquitectura.md` línea 66 (§0.2) sigue llamando `period` a la semana de la lista; C1 usa `week`. Dueño: U1. También `docs/criterios/U1.md` U1-16 sigue con `period` (4 apariciones; solicitud 3 de `correccion-C1-r2.md` sin aplicar).
- H-r3-04 · media · `docs/contratos/diccionario-datos.md` no registra `this_month`, `onward` ni `chat_id_masked`, y la línea 530 (`budget.valid_to_date`) no dice que se borran los renglones posteriores. Dueño: S1 (solicitudes 1 y 2 de `correccion-C1-r2.md`).
- Bajos: H-r3-05 (M-06 «lista vigente» en §9), H-r3-06 (`vat_on_interest_bp` NULL no distingue «falta capturar», K18), H6 (descripciones que repiten el identificador), H7 (PROVISIONAL de `DashboardWeekMeals` y `NotificationUpdate` sin `x-provisional`), H8 (`LoginRequest.username` sin 3–64, K6).

**S1 (aprobado; 9 bajos)**: R3-01 (`s1_verifica.py` no compara con el oráculo las cifras de frases «ilustrativo» con ruta); R2-02 (1,814.50 pasa como cuenta citada); R2-03 y R2-04 (el reporte `S1-esquema-1.2.md` conserva supuestos y cifras de antes de las correcciones); R2-05 (K18 sin regla «todo o nada» en la base; IVA NULL = «ya incluye»); R2-06 (RN-12 conserva «la que S1 recomienda»); R2-07 (`credit.rate_base` frente a `credit_purchase.rate_basis`, nombres que se confunden); R2-08 (cabecera de 00002 con fecha en formato de pantalla); R2-09 (10 ejemplos sin la leyenda «(sin dato en el oráculo)»).

**U1 (aprobado; 7 bajos)**: B1 (N-316 asigna el chat por defecto al usuario de la sesión; K4 dice «ella confirma»); B2 (el verificador no revisa `scope`, `family`, `Idempotency-Key` ni la fecha de captura); B3 («Para C1» pide `period` para N-151 y N-153); B4 (los domingos, dos tarjetas de Inicio muestran semanas distintas); B5 (`U1-1.1.md` con cifras viejas); B6 (la adenda K18 de «Para C1» está en `correccion-U1-r2.md` y no en `U1-1.1.md`); B7 («Para C1» no pide quitar `updatePerson` de N-289).

**Colorimetría (aprobada; 5 bajos)**: B-01 (gasto declarado en 45°; 10 de 15 tokens miden 45.03° a 45.17°); B-02 (el bloqueo `U2a-rosas-ximena-bodega` dice que la marca va en la introducción de §7 y está en §7.5); H-02 (§6.3 sin nota del cambio de matiz de Gasto); H-03 (marca 300 frente a Conjunto y Walmart sin tabla); H-04 (Gasto cerca de las series 3 y 4; las series no entran en el barrido de U2-08).

**Criterios (aprobados; 4 bajos)**: B-01 (`INDICE.md`: la cifra 248 de §1 ya quedó en 257 en esta integración; siguen las notas 8 y 10 de §5 y §6 sin los bloqueos nuevos); B-02 (C1-06 y C1-07 conservan «quincena 15/30» junto a K8); B-03 (U2 §5 marca 135° como convención propia y decisiones I lo prevé); B-04 (5 citas de §5 con fragmentos en otro orden: «SIN FUENTE» en `citas.py`).

## 6. Bloqueos nuevos de P1 (8, todos abiertos)
`docs/estado/bloqueos/C1-valores-k31.md` (C1), `S1-ajuste-objetivo-kcal.md` (S1), `S1-fuentes-precio.md` (S1), `U1-lista-vigente-domingo.md` (U1), `U2a-rosas-ximena-bodega.md` (colorimetría), `U2a-walmart-oscuro.md` (colorimetría), `criterios-frescura-extremos.md` (criterios) y `criterios-promedios-redondeo.md` (criterios). Pregunta, marca y fase que lo usa: `docs/estado/bloqueos/INDICE.md`, sección «Bloqueos nuevos de P1». Ninguno cambia cifras del oráculo ni el esquema.

## 7. Pendientes para WF-1 (con ruta y dueño)
1. **Decisión sobre C1** (orquestador): H-r3-01 y H-r3-02 en `backend/internal/api/c1gen/` y regenerar con `build.py`; H-r3-03 en `docs/design/U1-arquitectura.md` línea 66 y `docs/criterios/U1.md` U1-16 (`period` → `week`); H-r3-04 en `docs/contratos/diccionario-datos.md` línea 530 y la tabla de `telegram_chat_candidate`, y en el comentario de `budget.valid_to_date` (`00001_esquema.sql` línea 299).
2. **K18 e IVA «falta capturar»** (S1, luego C1 y U1): `vat_on_interest_bp` NULL hoy significa «ya incluye IVA». Hace falta un tercer estado o una columna `interest_includes_vat` que admita NULL, y sumarla a `missing_terms` (S1 R2-05, C1 H-r3-06 y U1 §8.4). Lo decide S1 con una migración nueva; si hay dos lecturas, va como bloqueo.
3. **«Para C1» de U1** (`docs/estado/p1/U1-1.1.md`): fusionar la adenda K18 de `correccion-U1-r2.md` y agregar «N-289: quitar `updatePerson`» (B6 y B7). C1 v1.1 ya expone los campos de K18 (`verificar_u1.py` da OK).
4. **Criterios** (agente de criterios de WF-1): `docs/criterios/S4.md` S4-07, citar RN-17 1.2 y `S1-ajuste-objetivo-kcal`; `docs/criterios/S5.md` S5-04, sumar `S1-fuentes-precio` (hoy 0 citas en los dos); `docs/criterios/C1.md`, criterio para los extractos y «bloqueo abierto = sin Respuesta de Axel» en C1-15 (solicitud 3 de `C1-v1.1.md`); notas 8 y 10 y §6 de `docs/criterios/INDICE.md` (B-01); B-02 a B-04 de criterios r3.
5. **Una sola copia de los verificadores**: `docs/estado/wf0/S1-modelo-datos.md` (Anexos A y B) y `docs/estado/wf0/U1-arquitectura.md` (anexo `u1_verify.py`) todavía traen el código. Las copias vigentes son `backend/db/tools/` y `docs/design/tools/verificar_u1.py`. El Makefile ya no usa los anexos. Hay que cambiarlos por un enlace.
6. **U4 y U2 (WF-1)**: en `docs/design/U4-ajustes.md`, los campos de `ProjectionSettings` (interruptor de K21, que hoy queda en «Presupuesto» a reserva de U4), de las restricciones con «faltan restricciones» (K16), del objetivo de peso (K23) y del tope por persona (K13). En U2 `DESIGN.md`, lo que decidan los dos bloqueos de U2a, la regla de refuerzo entre la marca y las categóricas, los pares cercanos de §7.7 y H-04 (series).
7. **Otros**: `docs/estado/wf0/REVISION-r3.md` §6, marcar R3-04 como atendido por la colorimetría 1.1. U1 decide si existe «Descartar chat detectado» (N-xxx nueva; solicitud 1 de `C1-v1.1.md`). En el esquema, renombrar `credit.rate_base` o `credit_purchase.rate_basis` (R2-07). Actualizar los reportes `S1-esquema-1.2.md` y `U1-1.1.md` con las cifras finales (R2-03, R2-04 y B5).
8. **Axel**: los 8 bloqueos de §6. Para cerrar la marca de C1 de H-r3-02 sin rodeos sirve responder `U1-lista-vigente-domingo`.

## 8. Paquetes instalados
Ningún agente de P1 instaló paquetes del sistema. Usaron goose v3.28.0 del scratchpad del orquestador, sqlite3 3.53.4, Go 1.27.2 (instalado antes del flujo, K28), Python del sistema, chromium ya instalado (para `paleta.png`) y la caché de `uv`/`uvx` (pyyaml, jsonschema y openapi-spec-validator 0.9.0).

## 9. Consumo por agente (de sus reportes; cifras aproximadas)

| Agente | Acciones | Leído (KB) | Escrito (KB) | Tope |
|---|---|---|---|---|
| S1-esquema-1.2 | 33 | 230 | 100 | no |
| revisor-S1-r1 | 32 | 195 | 37 | no |
| correccion-S1-r1 | 33 | 100 | 23 | no |
| revisor-S1-r2 | 37 | 210 | 31 | no |
| correccion-S1-r2 | 22 | 80 | 18 | no |
| revisor-S1-r3 | 35 | 210 | 47 | no |
| C1-v1.1 | 38 | 205 | 85 (más las salidas del generador) | no |
| revisor-C1-r1 | 45 | 155 | 35 | no |
| correccion-C1-r1 | 26 | 105 | 18 (a mano) | no |
| revisor-C1-r2 | 41 | 135 | 33 | no |
| correccion-C1-r2 | 34 | 115 | 24 (a mano) | no |
| revisor-C1-r3 | 38 | 140 | 33 | no |
| U1-1.1 | 36 | 195 | 65 | no |
| revisor-U1-r1 | 29 | 150 | 34 | no |
| correccion-U1-r1 | 33 | 125 | 16 | no |
| revisor-U1-r2 | 53 | 187 | 33 | no |
| correccion-U1-r2 | 26 | 135 | 28 | no |
| revisor-U1-r3 | 27 | 195 | 98 (72 temporales) | no |
| U2a-1.1 | 36 | 165 | 62 (más las salidas del script) | no |
| revisor-U2a-r1 | 20 | 140 | 37 | no |
| correccion-U2a-r1 | 32 | 125 | 50 | no |
| revisor-U2a-r2 | 29 | 130 | 40 | no |
| criterios-S-1.1 | ~40 | 220 | 95 | no |
| criterios-UC-1.1 | 34 | 165 | 85 | no |
| revisor-criterios-r1 | 34 | 200 | 45 | no |
| correccion-criterios-r1 | 27 | 85 | 35 | no |
| revisor-criterios-r2 | 29 | 190 | 37 | no |
| correccion-criterios-r2 | 8 | 60 | 6 | no |
| revisor-criterios-r3 | ~48 | 185 | 20 | no |
| integrador-p1 | ~36 | ~185 | ~55 | no |
| **Total (30 agentes)** | **~991** | **~4,717** | **~1,325** | ninguno |

Ningún agente pasó los límites del protocolo §4: el máximo fue 53 acciones (revisor-U1-r2), 230 KB leídos (S1-esquema-1.2) y 100 KB escritos (S1-esquema-1.2). El plan contra real en tokens (K33) lo reporta el orquestador, porque los agentes no ven su conteo de tokens.
