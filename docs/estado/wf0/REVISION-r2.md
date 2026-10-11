# Revisión WF-0 «Cimientos» · ronda 2 de 3

Revisor independiente `revisor-wf0-r2` · 11-oct-2026 · Puerta 1 (C1 + U1 + criterios; S1 y U2a como entregables de WF-0).
Evalúo **solo** contra `docs/criterios/` (S1, U1, C1 y los de U2 que aplican a la colorimetría de WF-0). No escribí nada de lo revisado. Todo se reejecutó con scripts propios (scratchpad `rev-r2/`) además de los verificadores de los agentes; no se aprobó nada por la palabra de un reporte.

Artefactos revisados (sha256, primeros 16): migración `747fc8fd88e576f4` · diccionario `edda526135849067` · reglas `b93aaf9bc094ccd9` · `openapi.yaml` `316ce18d6d1e4300` · `C1-datos.md` `12681dbf9f918e33` · `routes.go` `e0de67b9ebd257c6` · U1 `45fcfb61da0ca09f` (sin cambios desde r1) · `colorimetria.md` `beadf8708b445a60` · `contraste.csv` `9f995320c6fed1f6` · `C2-colores.json` `1446c402b3226d8e` · `colorimetria.py` `c70e643bb4300282` · `model.json` `1a57066689ddeed3` (igual al del INDICE).

## Veredicto

**NO APROBADO en ronda 2**: 2 hallazgos de severidad alta (R2-01, R2-02), 0 de media y 6 de baja.
- Lo que quedó resuelto: S1 corrigió todo lo suyo (H-01 en el esquema, H-02, H-04 y H-11): 54 tablas, aplica limpio, ida y vuelta del oráculo con 0 diferencias y evidencia por criterio. U2a resolvió H-03 (serie 6 en 225°, zona 80°–120° libre y verificada por script). Los criterios resolvieron H-06, H-07 y H-08. El generador de C1 ya vive en el repositorio y se reejecuta (H-05 para C1).
- Lo que bloquea: **C1 no se regeneró contra el esquema 1.1 de S1**. `openapi.yaml`, `routes.go` y `C1-datos.md` son del 10-oct 23:50:01 y la migración del 23:52:33. Por eso (a) las 4 necesidades de la lista de compras siguen sin operación (H-01 sigue abierto, ahora del lado de C1) y (b) el contrato describe columnas que ya no existen y omite 5 tablas nuevas (C1-03 falla). El propio verificador del agente (`verify_c1.py`) da `RESULTADO: 2 falla(s)` y `build.py --check` da `DIFIERE` en 2 de 3 archivos. Ninguna corrección de la ronda 1 tenía asignada esa regeneración (`correccion-C1-r1` solo tenía H-05 y lo declaró en sus Pendientes 1 y Supuestos 2).

## 1. Verificaciones ejecutadas por el revisor

| # | Qué pide la tarea | Comando o método | Resultado |
|---|---|---|---|
| 1 | SQL en base temporal y `foreign_key_check` | Sección Up extraída con awk → `sqlite3 -bail datos/tmp/rev.db`; Down sobre copia; además `goose -dir backend/db/migrations sqlite3 datos/tmp/rev-goose.db up/down/up` (goose v3.28.0, `goose --version`) | `UP_OK` · 54 tablas, las 54 STRICT · 147 índices · 68 disparadores · `PRAGMA foreign_key_check` = 0 filas · `integrity_check` = ok · Down deja 0 objetos · goose: «OK 00001_esquema.sql», down deja solo `goose_db_version` y `sqlite_sequence`, up otra vez 54 + 2. Semillas: 7 `app_setting`, 4 `meal_slot`, 7 `notification`, todas con valores de decisiones A, C, D, E, F y H. `datos/tmp` borrada al final |
| 2 | `go build ./...` y `go vet ./...` | En `backend/` con go1.26.5; además `gofmt -l .` y `go list ./...` | BUILD_OK · VET_OK · gofmt sin salida · 2 paquetes. `go.mod` declara `go 1.26` (bloqueo `S6-version-go`, decisión I pide 1.27) |
| 3 | `colorimetria.py`: `contraste.csv` sin `cumple=false`; hex de C2 = `colorimetria.md` | Script del repositorio copiado a un árbol espejo del scratchpad, `--sin-png` (no reescribe el repositorio), `cmp` contra el repositorio; luego `color_r2.py` y `u2_r2.py` propios (WCAG y OKLab escritos aparte, sin importar el script) | «VERIFICACIÓN OK»; md, csv, json y html **idénticos byte a byte** · 202 filas, 0 `cumple=false` · recálculo propio de los 202 pares: 0 difieren (truncado a 2 decimales), 0 bajo umbral · 190 tokens del JSON = 190 filas del Anexo A, 0 hex distintos; hex del CSV = JSON · las tablas de valores finales del cuerpo (§5.1–§6.2, §7.1, §7.2, §8.1, §8.3, §10) usan solo hex de tokens; los hex que no son token están solo en comparaciones, simulaciones y candidatos (§6.3, §7.5, §8.4, §8.7, §8.8) |
| 4 | Matriz C1: toda N-xxx con endpoint y todo endpoint con pantalla o justificación | `matriz_r2.py` propio (U1 §3 y §3.29 ↔ `x-necesidades` ↔ C1 §10); además `verify_c1.py --db datos/tmp/rev-goose.db` y `u1_verify.py` de los agentes | U1: 199 N-xxx (136 L, 63 E), Local N-003 y N-301 · 193 operaciones · **N-xxx sin operación ni Local: 4 (N-151, N-152, N-153, N-165)** · operación sin N-xxx: solo `getOpenApiDocument` (excepción «estáticos», C1 §11.2) · 41 pantallas, 0 sin operación ni Local · §10: 220 filas = 220 pares de U1, 4 en estado BLOQUEO. `verify_c1.py`: `FALLA C1-02`, `FALLA C1-03`, `RESULTADO: 2 falla(s)`. `u1_verify.py`: `FALLA U1-03 N-xxx sin endpoint en C1: ['N-151', 'N-152', 'N-153', 'N-165']` |
| 5 | Cada decisión A–I cubierta por S1, U1 o C1 | `decisiones_r2.py` (37 comprobaciones contra esquema, OpenAPI y U1) + lectura dirigida | §3: todas las decisiones de producto tienen respaldo en S1 y U1; C1 tiene operación para todas, pero en 5 renglones (B1, B2, B3, D1, D5) describe el esquema anterior (R2-02) y en D5 le faltan las escrituras de la lista (R2-01) |
| 6 | PROVISIONAL, supuestos y bloqueos: ¿esconden un valor inventado? | `provisional_r2.py` sobre migración, diccionario, reglas, OpenAPI, C1, U1, colorimetría y `main.go`; lectura de los 32 bloqueos y de las secciones Supuestos de 5 reportes | Toda marca PROVISIONAL de contenido remite a un bloqueo existente (las que no lo citan a ±3 líneas son leyendas o renglones de índice cuyo detalle sí lo cita). Ningún número inventado entra en cifras, esquema o contrato. Dos incoherencias de texto: C1 afirma premisas que S1 1.1 ya no cumple (R2-02) y S1 afirma que el oráculo no trae tope de Karen cuando sus notas dicen «cap de 374.93» (R2-03) |
| 7 | Criterios medibles y trazables | `criterios_r2.py` y `citas_r2.py` propios sobre los 18 archivos + INDICE | 248 criterios (189 alta, 59 media) = INDICE; 248 con «Verifica» y «Evidencia»; 248 de 248 en su §5 (con rangos); toda línea de §5 cita plan, decisiones, protocolo o «convención propia». 5 «Verifica» de inspección (P3-04, P3-11, S1-10, S4-09, S4-13) con resultado observable. 27 de 288 fragmentos citados en §5 son paráfrasis, no citas literales. Siguen las 4 incoherencias de H-10 y aparece una nueva (rutas `/tmp`) → R2-07 |
| 8 | Gasto en ≤ 3 toques demostrado en U1 | Lectura de U1 §4 (regla de conteo) y F-1; `u1_verify.py` | Demostrado: «+» (1), chip de categoría (2), «Guardar» (3); teclear no cuenta; fecha, persona y forma de pago por defecto; reparto desde la categoría; descripción no obligatoria (S1 1.1 ya la admite NULL). La prueba en Android sigue siendo de U5 |
| 9 | Regresión de S1 (Kueski 389500 y 239528) | `regresion_r2.py` propio: inserta 3 personas, Kueski, 4 despensas como gasto + compra a crédito con partes, 16 pagos con 32 partes y el papel de contado; reconstruye eventos, series y meses desde la base. Además `s1_verifica.py` extraído del Anexo B del reporte de S1 | Suma de pagos de la despensa 1 = **389500**; parte de Axel = **239528** (en la compra y como Σ de sus pagos). 17 eventos con A y B, series A/B de 14 valores, mínimos 809450 / 500000, P1–P6 ≡ oct-2026…mar-2027 en efectivo y devengado (total y por persona), 0 eventos entre los días 1 y 9: **142 valores, 0 diferencias**, `fk_check` 0. `s1_verifica.py`: «RESUMEN TODO OK», 164 valores, 46 de 46 negativas, goose [0, 0, 0] |

Otras comprobaciones: `build.py --check` del generador de C1 → `DIFIERE docs/contratos/openapi.yaml (repositorio 570515 bytes, generador 573729 bytes)`, `IDÉNTICO routes.go`, `DIFIERE C1-datos.md`; el generador, aun regenerando, sigue informando `N sin endpoint: [... 'N-151', 'N-152', 'N-153', 'N-165' ...]` y `tablas sin op: [idempotency_key, recipe_meal_slot, inventory_movement, shopping_list, shopping_list_item]`. `uvx --offline --from openapi-spec-validator==0.9.0 openapi-spec-validator docs/contratos/openapi.yaml` → `OK`, salida 0. Diccionario contra base (`dicc_r2.py`): 54 de 54 tablas con sección, columnas idénticas en nombre, orden y tipo, 54 con «Usado por», 36 columnas `_cents` INTEGER, 0 REAL/FLOAT/DOUBLE/NUMERIC/BLOB, 31 de 31 entidades del plan.

## 2. Tabla de criterios

Estados: **Cumple** · **Cumple con observación** · **No cumple** · **No aplica en WF-0**.

### S1 · Dominio y datos
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| S1-01 Entidades completas | alta | Cumple | 31 de 31 entidades en el mapa del diccionario con su `nombre_codigo`; las 54 tablas tienen sección y existen en la base |
| S1-02 Migraciones en SQLite limpio | alta | Cumple | sqlite3: Up ok, `integrity_check` ok, `foreign_key_check` vacío; goose v3.28.0 up/down/up ok |
| S1-03 Cada regla con ejemplo del reporte | alta | Cumple con observación | 22 de 22 reglas con ejemplo; toda ruta citada de `model.json`/`data.json` existe y su valor aparece en el texto (`reglas2_r2.py`). 10 reglas sin datos en el oráculo llevan ejemplo «ilustrativo» (RN-10, RN-11, RN-13 a RN-19, RN-22); el criterio literal no prevé ese caso (R2-07) |
| S1-04 Ninguna tabla huérfana | alta | Cumple | 54 de 54 con «Usado por». Comprobadas contra U1: `shopping_list` y `shopping_list_item` (P-13 N-150 a N-153, P-14 N-165), `inventory_movement` (P-16 N-181, N-183), `recipe_meal_slot` (P-10 N-122, P-11 N-135), `idempotency_key` (justificada por escrito: API y job de purga de S6) |
| S1-05 Dinero exacto | alta | Cumple | 0 columnas REAL/FLOAT/DOUBLE; 36 `_cents`, todas INTEGER; `12.5` rechazado (negativas de `s1_verifica.py`) |
| S1-06 Fechas y zona | media | Cumple | grep de formato de pantalla en la migración = 0; `CHECK (date(x) IS x)`; zona declarada |
| S1-07 Personas y acceso | alta | Cumple | Cuenta de Karen, tercera cuenta e ingreso de Ximena rechazados; paso a externa con cuenta rechazado (`trg_person_external_without_access`); sin columna de rol. H-11 cerrado |
| S1-08 Reparto y participantes | alta | Cumple | Suma ≤ 100 % como invariante 8 (API, prueba de S6); despensa 1 cabe sin pérdida: Σ partes = total en las 4 despensas (regresión) |
| S1-09 Periodos e ingresos | alta | Cumple | 18 filas, quincena 2 de feb-2027 paga el 28, ventana del 10 al 9 rechazada (`s1_verifica.py`); equivalencia 12 de 12 por mes calendario y 0 eventos entre los días 1 y 9 (regresión propia) |
| S1-10 Créditos | alta | Cumple | 4 tipos con CHECK; tasas 53450/309450 y 80050/309450; `availability_reading DEFAULT 'B'` y `availability_confirmed`; promoción con comercio, plazo, mínimo y vigencia |
| S1-11 Comida, salud, inventario, precios | media | Cumple | 12 de 12 puntos con tabla y columna |
| S1-12 Ida y vuelta de la semilla | alta | Cumple | Propia: 142 valores, 0 diferencias; del agente: 164 valores, 0 diferencias; 389500 y 239528 |
| S1-13 Ajustes tiene dónde guardarse | media | Cumple | Tabla «14 secciones → campos → tabla.columna» en el diccionario, 14 de 14; Apariencia local por dispositivo, igual que U1 (N-301). H-04 cerrado |
| S1-14 Sesiones y credenciales | alta | Cumple | Solo `password_hash` argon2id (CHECK) y `token_hash`; 0 tokens de Telegram |
| S1-15 Integridad y consultas de C1 | media | Cumple | `foreign_key_list` (muestra del revisor): `shopping_list_item` 4, `inventory_movement` 3, `budget` 1, `credit_purchase_share` 2, `expense_share` 2; los 6 índices de filtros existen |
| S1-16 Sin valores inventados | alta | Cumple con observación | Semillas decididas; PROVISIONAL con bloqueo; Supuestos sin cambios de cifras. La descripción del tope de Karen afirma algo que el oráculo contradice (R2-03) |
| S1-17 Idioma y nombres | media | Cumple | 633 identificadores en `[a-z0-9_]`, glosario sin repetidos, textos en es-MX |

### U1 · Arquitectura de información (documento sin cambios desde r1)
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| U1-01 Barra lateral de 9 entradas | alta | Cumple | 9 en el orden exacto; Buscar y Registrar gasto son acciones fijas (§1.2) |
| U1-02 Inventario con identificador | alta | Cumple | 27 P-xx + 14 P-26.xx = T 41; 199 N-xxx únicas con tipo L o E |
| U1-03 Matriz pantalla ↔ necesidad ↔ endpoint | alta | **No cumple** | N-151, N-152, N-153 y N-165 sin operación (R2-01); el reporte de U1 aún dice «U1-03: PENDIENTE» |
| U1-04 Gasto en ≤ 3 toques | alta | Cumple | 3 toques con regla escrita (§4) y valores por defecto |
| U1-05 Seis flujos clave | alta | Cumple con observación | 8 flujos (6 + F-7 y F-8), todos los pasos con pantalla existente; el criterio dice 6 (R2-07). F-3 ruta A paso 3 sin operación (R2-01) |
| U1-06 Vistas por entidad | alta | Cumple | 31 entidades: 23 con vistas, filtro y orden; 8 con motivo (sin cambios) |
| U1-07 Menú «/», atajos, edición | media | Cumple | 27 comandos con destino; 15 atajos, 0 duplicados (`u1_verify.py`) |
| U1-08 Ajustes con 14 secciones | alta | Cumple | 14 de 14 en orden y con nombre exacto |
| U1-09 Cobertura de las decisiones | alta | Cumple | 26 capacidades con P-xx y N-xxx |
| U1-10 Dos titulares, Karen sin acceso | alta | Cumple | Sin pantalla de acceso de Karen; 56 menciones en las áreas permitidas |
| U1-11 Móvil primero, una acción primaria | media | Cumple con observación | §6.4 con 27 filas; §6.1 punto 4 sigue con 24 × 24 px frente a ≥ 44 px (R2-05) |
| U1-12 Consistencia con S1 y C1 | media | Cumple con observación | Toda tabla de S1 tiene pantalla o motor; C1 no coincide con S1 1.1 en campos de presupuesto, receta, gasto e inventario (R2-02) |
| U1-13 Idioma | media | Cumple | 77 nombres visibles, 0 en inglés |
| U1-14 Supuestos y bloqueos | media | Cumple | Sección Supuestos; 2 bloqueos U1; PROVISIONAL marcado |

### C1 · Contrato de datos y API
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| C1-01 OpenAPI 3.1 válido | alta | Cumple | openapi-spec-validator 0.9.0 → OK, salida 0; 3.1.1, 154 rutas, 193 operaciones |
| C1-02 Toda N-xxx con operación y viceversa | alta | **No cumple** | 4 N-xxx sin operación ni Local; una necesidad bloqueada no cuenta como cubierta (R2-01) |
| C1-03 Endpoints respaldados por tablas | alta | **No cumple** | 5 propiedades sin columna en el esquema vigente y 5 tablas de S1 sin operación ni justificación en C1 (R2-02) |
| C1-04 Una convención de errores | alta | Cumple | 874 de 874 respuestas de error → `Error` |
| C1-05 Una convención de paginación | alta | Cumple | 39 de 39 colecciones con `cursor`/`limit`/`sort` y `Page` |
| C1-06 Una convención de filtros | alta | Cumple | 6 de 6 por `$ref`, 0 sinónimos; `family` PROVISIONAL con bloqueo |
| C1-07 Una convención de fechas | alta | Cumple con observación | Solo `date` y `date-time`; el patrón de `DateTime` acepta «Z» y `rate_pct` lleva 4 decimales (R2-06, antes H-12) |
| C1-08 Una convención de dinero | alta | Cumple | 253 propiedades `_cents` → `MoneyCents`; 0 `number` |
| C1-09 Métricas con fórmula y periodo | alta | Cumple | Sin cambios desde r1; ejemplos 10,000.00 / 6,905.50, 973.75 y 1,947.50 presentes y del oráculo |
| C1-10 Recursos y vistas de las 31 entidades | alta | Cumple | §7 con 31 filas (las tablas nuevas de S1 son complementarias, no entidades del plan) |
| C1-11 Cobertura funcional | alta | Cumple con observación | 35 capacidades; la lista de compras solo tiene lectura (`getShoppingList`) y sus escrituras faltan (R2-01) |
| C1-12 Seguridad | alta | Cumple | Públicas solo login y salud; 0 campos sensibles; cookie; 429 `too_many_attempts` |
| C1-13 Versionado y estructura | media | Cumple | 0 `operationId` repetidos; 9 etiquetas; C1 v1.0 con fila «corrección r1» |
| C1-14 Idioma y ejemplos | media | Cumple | 13 ejemplos válidos con jsonschema; 2,607 descripciones en español |
| C1-15 PROVISIONAL marcado | alta | Cumple con observación | 63 marcas en YAML y 44 líneas en C1, todas con bloqueo existente; pero `C1-huecos-u1-s1`, `C1-lista-compras` y varias descripciones parten de premisas que S1 1.1 ya no cumple (R2-02) |

### U2 · aplicados a la colorimetría de WF-0 (U2a)
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| U2-01 DESIGN.md con 10 secciones | alta | No aplica en WF-0 | DESIGN.md es de WF-1 |
| U2-02 Matiz C aplicado | alta | Cumple | Matiz medido (OKLab propio, croma > 0.02): marca 279.0–280.6°, éxito 149.1–150.4°, error 19.5–20.3°, advertencia 73.8–75.7°, información 238.7–240.1°, ingreso 134.5–135.5° (propuesta 135° que permite decisiones I), gasto 23.8–25.1°, ahorro 199.8–200.9°, deuda 309.4–310.2°, acción crítica 99.5–100.8°; marca y neutro 50/950 = plan. **0 tokens fuera de `color.critical` en 80°–120°**; serie 6 en 225.3°. H-03 cerrado |
| U2-03 Método y saturación | alta | Cumple | Desviación de L por paso ≤ 0.0010; croma máx. marca 0.1103, neutro 0.0121, roles 0.1408 (redondeo a hex) |
| U2-04 0 pares bajo umbral | alta | Cumple | 202 pares, 0 bajo umbral; recálculo propio de los 202 sin diferencias |
| U2-05 Pares estructurales | alta | Cumple | 11.53, 5.40, 3.47, 3.98, 14.07, 7.97, 10.38, 1.20, 1.90; borde de campo 3.71 documentado |
| U2-06 Roles como texto | alta | Cumple | Paso 700 (claro) y 400 (oscuro) ≥ 4.5 sobre fondo y tarjeta en los 8 roles |
| U2-07 Botón primario | alta | Cumple | Blanco/marca 700 = 5.77; marca 700/fondo = 5.53; marca 950/marca 300 = 10.47; marca 300/fondo 950 = 10.38 |
| U2-08 Colisiones resueltas | alta | Cumple | Pares < 20° y ΔL < 0.15: 3 en claro y 4 en oscuro, los mismos de §7.7 con su refuerzo |
| U2-09 Paleta de 8 series | alta | Cumple | ΔL consecutivo mín. (L OKLab propio) 0.174 claro y 0.218 oscuro; script: gris 0.179 / 0.223; contraste mín. 3.14 / 3.10; simulación documentada |
| U2-10 a U2-13 | alta/media | No aplica en WF-0 | WF-1 |
| U2-14 Idioma y supuestos | media | Cumple | es-MX, tokens en inglés, Supuestos (ancho de la zona) y bloqueo `U2-gasto-error` |

### Criterios (archivo `docs/criterios/`)
| Aspecto | Estado | Evidencia |
|---|---|---|
| Completos y con estructura | Cumple | 18 archivos + INDICE; 248 criterios = INDICE (189 / 59); secciones 1 a 5 en todos |
| Medibles | Cumple | 248 de 248 con «Verifica» y «Evidencia»; los 5 de inspección tienen resultado observable |
| Trazables al plan | Cumple con observación | 248 de 248 en su §5; 27 de 288 fragmentos citados son paráfrasis (R2-07) |
| Coherentes con entregables y decisiones | Cumple con observación | H-06, H-07 y H-08 cerrados; quedan U2-02 (135°), U1-05 (6 flujos), C1-01 (módulo ausente), S1-03 (ejemplos ilustrativos), rutas `/tmp` y `verify_c1.py` sin citar (R2-07) |

## 3. Cobertura de las decisiones A a I

| Id | Decisión (resumen) | S1 (esquema o regla) | U1 (pantalla) | C1 (operación) | Estado |
|---|---|---|---|---|---|
| A1 | Co-titulares, mismos permisos; usuario y contraseña; sesión larga; cierre desde Ajustes; argon2id; 5 en 15 min | `user_account` sin rol, CHECK argon2id; `device_session.expires_at`; `login_attempt`; semilla `auth.*` 5 y 15 | P-00, P-26.01 | `login` (429), `logout`, `listSessions`, `revokeSession(s)` | Cubierta (duración PROVISIONAL, `criterios-login-bloqueo`) |
| A2 | Karen externa sin cuenta; despensa y compras Kueski; paga directo; informativa; estado de cuenta PDF | `participation = external` + 9 disparadores; `pays_share_directly`; partes; `generated_report` | P-25, P-26.12 | `getPersonParticipation`, `getExternalStatement`, `generateExternalStatement`; sin alta de cuentas ni personas | Cubierta (reposición por parte PROVISIONAL) |
| A3 | Comidas y salud solo Axel y Ximena | Disparadores en `meal_plan_entry`, `health_profile`, `weight_entry` | P-10, P-21 | 422 `external_person` | Cubierta |
| A4 | Telegram a Axel y Ximena, destinatarios por aviso | `notification_recipient.chat_id`, `notification_subscription` | P-26.09 | `setNotificationRecipients`, `createNotificationRecipient` | Cubierta (registro del chat: `S1-huecos-negocio`) |
| B1 | MXN, centavos, HALF_UP, último ajusta; zona; fechas ISO; es-MX | 36 `_cents` INTEGER STRICT; CHECK ISO; RN-01, RN-07; C-16 descripción NULL | §0.4 | `MoneyCents`, `Date`, `DateTime` | Cubierta; C1 aún rellena la descripción con la categoría (R2-02) |
| B2 | Mes calendario; quincena 15/30; despensa devengada a la mitad | `period` con CHECK; `accrual_months`; `budget` con vigencia | §0.2, P-06 | `period`, `getAccruedGroceries`, presupuestos | Cubierta en S1 y U1; C1 describe `budget.period_id`, que ya no existe (R2-02) |
| B3 | Ingresos solo de Axel, 15/30, extras | `income`, `planned_item` + disparadores; `second_half_amount_cents` | P-08, P-26.03 | `incomes`, `planned-items` | Cubierta; C1 sin el monto de la quincena del 30 (R2-02) |
| B4 | Reparto % por gasto o categoría, 50/50; qué se comparte con Karen | `cost_share_rule(_part)`, regla en categoría, gasto y renglón; `credit_purchase_share` | N-051, N-281, N-309 | `cost-share-rules`, filtro `family` | Cubierta (`family` PROVISIONAL) |
| B5 | Fondo y metas con fecha | `savings_goal` (`target_date`, `planned_monthly_cents`), `savings_movement` | P-09 | `savings-goals`, métricas de ahorro | Cubierta |
| B6 | Comida fuera como gasto con persona; plan «fuera» | `expense.person_id`, `meal_plan_entry.is_eating_out` | N-124 | `createMealPlanEntries` | Cubierta |
| C1 | 4 tipos de crédito | `credit.kind` con CHECK | P-17, P-18, P-26.05 | `credits` | Cubierta |
| C2 | Tasa por compra con IVA; tabla por plazo 17.27 / 25.8685; promociones | `interest_cents` cotizado, `interest_includes_vat`, `credit_term_rate` (fracción), `promotion` | N-218, N-286, N-287 | `term-rates`, `promotions` | Cubierta (visualización de tasa: R2-06) |
| C3 | A/B por defecto B, aviso «pendiente confirmar» | `availability_reading DEFAULT 'B'`, `availability_confirmed` | N-284, N-285 | `previewAvailability` | Cubierta |
| C4 | Comparador: 3 restricciones, 4 objetivos, recomendación, sin solver | `app_setting optimizer.*` (pesos iguales sembrados; restricciones sin valor) | P-20, F-7 | `compareScenarios` | Cubierta (valores PROVISIONAL, `S1-restricciones-optimizador`) |
| C5 | Regresión Kueski | Regresión del revisor: 142 valores, 0 diferencias | — | Ejemplos del oráculo | Cubierta en WF-0 (motores en WF-1) |
| D1 | Recetario de Claude + propias; restricciones vacías | `recipe.origin`, `health_profile.dietary_restrictions`, `recipe_meal_slot` | P-11, P-12 | `recipes` | Cubierta en S1 y U1; C1 conserva un solo tiempo por receta (R2-02) |
| D2 | 4 tiempos; porciones base 2 | `meal_slot` (4, colación opcional), `base_servings DEFAULT 2` | P-10 | `listMealSlots`, plan | Cubierta |
| D3 | Perfil, IMC, meta kcal, kcal por receta y día; macros después | `health_profile`, `weight_entry`, macros en `nutrition_fact` | P-21, P-22 | `health-profiles`, `getTodayKcal` | Cubierta (meta kcal PROVISIONAL) |
| D4 | Fuentes nutricionales autorizadas | `nutrition_fact.source`, `source_license` | P-26.07 | `reloadNutritionBase` | Cubierta |
| D5 | Inventario híbrido; sobrantes; mercados con día de visita; lista de compras | `inventory_item`, `inventory_movement`, `store.visit_weekday`, `store_product_group`, `shopping_list(_item)` | P-13, P-16, P-26.06 | `inventory-items`, `stores`, `getShoppingList` | S1 y U1 cubren; **C1 sin escrituras de la lista (R2-01) y con movimientos derivados de la bitácora (R2-02)** |
| E1 | Manual, OCR con confirmación, lectura asistida apagada, PROFECO | `price.source`, `store.assisted_reading_enabled = 0`, semilla `price.assisted_reading_enabled = false` | P-05, F-8 | `confirmTicketScan`, `importProfeco` | Cubierta |
| E2 | Viejo a los 30 días, pesa menos | semilla `price.stale_after_days = 30` | N-170, N-171 | frescura en precios | Cubierta (peso PROVISIONAL) |
| E3 | Sam's: punto de equilibrio | `requires_membership`, `has_membership`, `membership_annual_fee_cents` | N-164 | `getBasketComparison` | Cubierta (cálculo PROVISIONAL) |
| E4 | Unidad base, sustitutos, mejor tienda | `product.base_unit`, `presentation`, RN-20 | N-161 a N-163 | comparativa | Cubierta |
| F1 | Diario, semanal, mensual PDF día 1, vista web; estado de Karen | `notification` (7 avisos), `job`, `generated_report` | P-23, P-24, P-25, P-26.08 | reportes | Cubierta (hora del diario sin valor, bloqueo) |
| F2 | 4 alertas | Semillas `payment_due` [1, 0], `budget_threshold` 80/100, `stale_prices`, `meal_plan_missing` | P-26.09 | `notifications`, `getBudgetAlerts` | Cubierta |
| F3 | MultiX texto plano, token en tiempo de ejecución | Sin token en la base (C-14) | N-299 | 0 campos sensibles | Cubierta |
| G1 | Fiel a Notion | — | §1, §2.5, §2.6, §5 | `queryView` (4 modos) | Cubierta |
| G2 | Android y escritorio; PWA mkcert; ≤ 3 toques | `idempotency_key` para reintentos | F-1, §6 | `getCaCertificate` | Cubierta (PWA en WF-3) |
| G3 | Tema del sistema con interruptor | Sin dato en la base (local) | P-26.10 (N-301, Local) | sin API | Cubierta (coherente en S1, U1 y C1; H-04 cerrado) |
| G4 | Paleta nueva; matiz de I | — | — | — (U2a: 190 tokens) | Cubierta |
| G5 | Resumen ≤ 6 tarjetas, Detalle, Tendencias, Calendario | — | P-01 | 6 métricas de tablero | Cubierta |
| G6 | Ajustes, 14 secciones | Tabla 14 secciones → tabla.columna | P-26.01 a P-26.14 | C1 §11 | Cubierta |
| G7 | Autoalojado; cifras tabulares; WCAG 2.2 AA | — | §0.4 | — (U2a: contraste) | Parcial: tipografía e iconos son de WF-1; táctil (R2-05) |
| H1 | Solo LAN, `finanzas.local:8443` | — | — | `servers` = `https://finanzas.local:8443`; esqueleto en 127.0.0.1 (`criterios-bind`) | Cubierta |
| H2 | Go + SQLite + systemd; respaldos locales | `backup` (`integrity_check`, `restore_tested_at`) | P-26.11 | `listBackups`, `createBackup`, `restoreBackup` | Cubierta (systemd en WF-4) |
| H3 | IA opcional apagada | Semillas `ai.* = false` | P-26.11 | `updateAiProvider`, `suggestRecipeWithAi` | Cubierta |
| H5 | Semilla e importador | `planned_item` | P-26.11 | `importPlannedItems`, `getSeedStatus` | Cubierta |
| I1 | Pila A | goose v3.28.0 (validado), Go `net/http` | — | `routes.go` | Cubierta (Go 1.27: bloqueo `S6-version-go`) |
| I2 | Matiz C y reglas | — | — | — (U2a) | Cubierta (H-03 cerrado) |
| H4, H6, J | Proceso (paquetes, repositorio, puertas) | — | — | — | No aplica al producto |

## 4. Bloqueos, supuestos y PROVISIONAL

- 32 bloqueos en `docs/estado/bloqueos/` (16 de criterios, uno ya resuelto; 7 de S1; 2 de U1; 1 de U2; 5 de C1; 1 de S6), todos con el formato del protocolo. Los nuevos de la ronda 1 (`S1-lista-compras`, `S1-huecos-negocio`) son preguntas de producto legítimas y se responden junto con `C1-lista-compras` y `C1-huecos-u1-s1`.
- Ningún valor inventado entra en cifras, esquema o contrato sin marca: las semillas de la migración son valores de decisiones; los valores técnicos de C1 (200, 24 h, 12 caracteres, 1 GB, `SameSite=Strict`) siguen marcados con `C1-parametros-tecnicos`; las definiciones de campos calculados con `C1-campos-calculados`.
- Incoherencias que sí deben corregirse antes de la puerta 1:
  - `C1-huecos-u1-s1` y `C1-lista-compras` dicen en «Cómo seguí» que S1 no tiene dónde guardar los datos; `S1-huecos-negocio` y `S1-lista-compras` dicen que 4 de los 7 huecos y la lista ya están en el esquema. Axel recibiría la misma pregunta con dos planteamientos opuestos (R2-02).
  - El tope de Karen: S1 dice «el oráculo no trae un tope»; `model.json` `papel[2].nota` y `papel[4].nota` dicen «fuera del cap de 374.93»; C1 dice «Karen 37493 en el oráculo». No hay número inventado (la columna queda NULL), pero las afirmaciones se contradicen (R2-02, R2-03).
- Supuestos de los reportes: ninguno cambia cifras, esquema, endpoints ni tokens. El ancho de la zona reservada de U2a (20° a cada lado) no mueve ningún token: la serie 6 queda a 125° del opuesto y advertencia (75°, decidida) impide una zona más ancha.

## 5. Estado de los hallazgos de la ronda 1

| Id r1 | Sev. | Estado en r2 | Evidencia |
|---|---|---|---|
| H-01 | alta | **Abierto (lado C1)** | S1 agregó las tablas; C1 no tiene operaciones para N-151, N-152, N-153 ni N-165 → R2-01 |
| H-02 | media | S1 cerrado; C1 abierto | Esquema corregido; C1 sigue con las lecturas PROVISIONALES anteriores y con propiedades sin columna → R2-02 (sube a alta por C1-03) |
| H-03 | media | Cerrado | Serie 6 en 225°; 0 tokens en 80°–120° fuera de `color.critical` (cálculo propio) |
| H-04 | media | Cerrado | 17 líneas S1-xx; `s1_verifica.py` reejecutable desde el Anexo B (TODO OK); tabla de 14 secciones; Apariencia local |
| H-05 | media | Cerrado para C1; parcial en el resto | `backend/internal/api/c1gen/` con `build.py --check` y `verify_c1.py`; los scripts de S1, U1 y criterios siguen solo en anexos de reportes o en el scratchpad → R2-08 |
| H-06 | media | Cerrado | Criterios con P-xx y N-xxx; `criterios-nxx` resuelto |
| H-07 | media | Cerrado | S1-09 con meses calendario y prueba de equivalencia; nota en S2-09 |
| H-08 | media | Cerrado | 248 de 248 criterios en su §5 (cálculo propio) |
| H-09 | baja | Abierto | U1 §6.1 punto 4 sigue en 24 × 24 px → R2-05 |
| H-10 | baja | Abierto | U2-02, U1-05, C1-01 y S1-03 sin cambio → R2-07 |
| H-11 | baja | Cerrado en el esquema | Dos cuentas como máximo y disparador; queda la redacción del tope → R2-03 |
| H-12 | baja | Abierto | `DateTime` acepta «Z»; `rate_pct` con 4 decimales → R2-06 |

## 6. Hallazgos de la ronda 2

| Id | Sev. | Artefacto | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|---|
| R2-01 | **alta** | C1 | C1-02, U1-03 (C1-11, U1-05) | H-01 sigue abierto. S1 1.1 ya tiene `shopping_list` y `shopping_list_item` (PROVISIONAL, opción A), pero C1 no se regeneró: N-151 (marcar comprado), N-152 (agregar artículo), N-153 (cambiar tienda o cantidad) y N-165 («Usar esta tienda», paso 3 de la ruta A de F-3) siguen sin operación, con estado BLOQUEO y el motivo «S1 no tiene dónde guardarla», que ya es falso. `verify_c1.py` → `FALLA C1-02 … 4 ['N-151', 'N-152', 'N-153', 'N-165']`; `u1_verify.py` → `FALLA U1-03`; `matriz_r2.py` → lo mismo. N-127 y N-073 siguen PROVISIONALES sin persistencia | `docs/contratos/openapi.yaml`; `docs/contratos/C1-datos.md` §10, §11.3, §13, §14; `backend/internal/api/routes.go`; `backend/internal/api/c1gen/gen3_ops.py`; `docs/estado/bloqueos/C1-lista-compras.md` | Asignar a C1 la regeneración. En `gen3_ops.py`, agregar operaciones sobre `shopping_list_item` para N-151, N-152, N-153 y N-165, y la persistencia de N-127 y N-073 (todas PROVISIONALES con `S1-lista-compras`/`C1-lista-compras` hasta la puerta 1). Correr `build.py`, `build.py --check` (IDÉNTICO ×3), `verify_c1.py --db` (0 fallas), el validador y `go build`/`go vet`. Actualizar «Cómo seguí» de `C1-lista-compras` para que ambos bloqueos planteen una sola pregunta |
| R2-02 | **alta** | C1 | C1-03 (C1-15, U1-12) | C1 describe el esquema anterior. `build.py --check` → `DIFIERE openapi.yaml (570515 → 573729 bytes)` y `DIFIERE C1-datos.md`. `verify_c1.py --db` → `FALLA C1-03`: 5 propiedades sin columna (`Budget.period_id`, `BudgetCreate.period_id`, `Recipe.default_meal_slot_id`, `RecipeCreate.default_meal_slot_id`, `RecipeUpdate.default_meal_slot_id`) y 5 tablas sin operación ni justificación (`idempotency_key`, `inventory_movement`, `recipe_meal_slot`, `shopping_list`, `shopping_list_item`). Textos contrarios a S1 1.1: `ExpenseCreate.description` «si falta, la API usa el nombre de la categoría (PROVISIONAL: S1 lo exige NOT NULL)» y `Expense` la exige, contra C-16 («admite NULL; la API no inventa un texto de relleno»); `fortnight_payment_cap_cents` «(Karen 37493 en el oráculo)» contra «NULL = sin tope, lo da Axel»; movimientos de inventario derivados de la bitácora y motivo no guardado, contra `inventory_movement.reason`; presupuesto por mes contra vigencia; un tiempo por receta; un solo monto quincenal. Quien implemente WF-1 o WF-3 con C1 programaría contra columnas que no existen | `docs/contratos/openapi.yaml`; `docs/contratos/C1-datos.md` §10, §13, §14; `backend/internal/api/c1gen/gen2_schemas.py`, `gen3_ops.py`, `c1_plantilla.md`; `docs/estado/bloqueos/C1-huecos-u1-s1.md` | En la misma regeneración de R2-01: ajustar `gen2_schemas.py` y `gen3_ops.py` al esquema 1.1 (vigencia de `budget`, tiempos de receta como lista sobre `recipe_meal_slot`, descripción opcional sin relleno, `inventory_movement` como origen de N-181 y N-183, `second_half_amount_cents` para N-279 y N-280, `idempotency_key` declarada en C1 §11 como infraestructura de S6 con su justificación, tope de Karen sin cifra). Quitar las marcas PROVISIONAL ya resueltas y dejar en `C1-huecos-u1-s1` solo los 3 puntos de `S1-huecos-negocio`. Verificar: `build.py --check` IDÉNTICO y `verify_c1.py --db` sin fallas en C1-03 |
| R2-03 | baja | S1 | S1-16 | La descripción de `person.fortnight_payment_cap_cents` afirma «el oráculo no trae un tope (374.93 de `karen_fija.pago` es la parte de Karen por pago, no un tope)». Pero `model.json` `papel[2].nota` («credito sigue activo; fuera del cap de 374.93») y `papel[4].nota` («fuera del cap de 374.93») hablan de un tope de 374.93 para Karen. No entra ninguna cifra (NULL), pero se presenta una interpretación como hecho, contraria a C1 | `backend/db/migrations/00001_esquema.sql` (columna `fortnight_payment_cap_cents`); `docs/contratos/diccionario-datos.md` (`person`) | Redacción neutral y la pregunta como bloqueo, por ejemplo en `S1-restricciones-optimizador`: «¿El tope por quincena de Karen del objetivo 2 es 374.93 (nota «cap» del oráculo) o lo fija Axel?». Regenerar el diccionario |
| R2-04 | baja | S1 | Trazabilidad de cambios (INDICE §1) | `reglas-negocio.md` sigue diciendo «Versión: 1.0», aunque la ronda 1 cambió C-13, C-15, C-16, RN-05, RN-14, RN-15 y RN-22; el diccionario sí pasó a 1.1 con su lista de cambios | `docs/contratos/reglas-negocio.md` (encabezado) | Subir a 1.1 con la lista de cambios, como en el diccionario |
| R2-05 | baja | U1 | U1-11 (U2-11) | H-09 sin atender: §6.1 punto 4 fija objetivos táctiles «de al menos 24 × 24 px»; el plan §4.1 punto 9 y U2-11 piden ≥ 44 px | `docs/design/U1-arquitectura.md` §6.1 | Citar ≥ 44 px como mínimo táctil y dejar 24 px solo como piso WCAG 2.5.8 |
| R2-06 | baja | C1 | C1-07; decisión C | H-12 sin atender: el patrón de `DateTime` acepta «Z», pero C1 §4 dice que los instantes llevan el desplazamiento de America/Mexico_City; `rate_pct` con 4 decimales mostraría 17.2725 % a 2 quincenas frente al 17.27 % de la decisión C | `docs/contratos/openapi.yaml` (`DateTime`, `rate_pct`); `docs/contratos/C1-datos.md` §4 | Incluirlo en la regeneración de R2-01/R2-02: restringir el patrón de salida o documentar que «Z» solo se acepta en entradas; definir la regla de visualización de tasas |
| R2-07 | baja | criterios | U2-02, U1-05, C1-01, C1-02, C1-03, S1-03, S6-03; protocolo §7 | Siguen las incoherencias de H-10: U2-02 da ingreso 150° sin «o U2 propone 135°» (U2a entregó 135° como propuesta); U1-05 pide exactamente 6 flujos con 4 obligatorios (U1 entrega 6 + F-7 y F-8); el primer comando de C1-01 usa un módulo que no está en el Python del sistema (salida 1; la alternativa `uvx --from openapi-spec-validator==0.9.0` sí funciona); S1-03 exige cifras del oráculo en reglas sin datos en el oráculo. Nuevas: los «Verifica» escriben en `/tmp` (S1-01, S1-02, S1-05, S1-15, S2 a S7 y C2), contra protocolo §7 y el uso de `datos/tmp` del flujo; C1-02, C1-03 y S6-03 no citan el verificador entregado `backend/internal/api/c1gen/verify_c1.py` ni `build.py --check`; 27 de 288 fragmentos citados en §5 son paráfrasis | `docs/criterios/U2.md`, `U1.md`, `C1.md`, `S1.md`, `S2.md` a `S7.md`, `C2.md`, `S6.md` | Admitir 135° como propuesta para la puerta 2; «≥ 6 flujos con los 4 obligatorios»; poner primero el comando con `uvx`; permitir ejemplos «ilustrativos» marcados solo en reglas sin datos en el oráculo; cambiar `/tmp/...` por `datos/tmp/...`; citar `verify_c1.py --db` y `build.py --check`; volver literales las citas o quitarles las comillas |
| R2-08 | baja | S1 | INDICE §1 (reejecución); H-05 | El diccionario dice que se genera de la migración y «no se edita a mano», pero su generador (`gen_dict.py`) y la prueba `s1_verifica.py` solo existen en los Anexos A y B del reporte de S1. Igual `u1_verify.py` (anexo del reporte de U1) y los scripts de criterios (`traza.py`, `citas.py`, `canon.py`, `s109_*.py`, `matriz_c102.py`, en el anexo de `correccion-criterios-r1.md`). Se pueden extraer y correr (el revisor extrajo `s1_verifica.py`: TODO OK), pero no son archivos del repositorio | `docs/estado/wf0/S1-modelo-datos.md` (Anexos A y B); `docs/estado/wf0/U1-arquitectura.md` (Anexo); `docs/estado/wf0/correccion-criterios-r1.md` (Anexo) | El integrador copia los scripts a rutas del repositorio (p. ej. `backend/db/tools/gen_dict.py` y `s1_verifica.py`, `docs/design/tools/verificar_u1.py`, `docs/criterios/tools/`) y los reportes citan esas rutas |

## 7. Lo que la ronda 3 debe traer para aprobar
1. R2-01 y R2-02 resueltos con **una sola regeneración de C1** contra la migración `747fc8fd…` (o la que quede tras responder los bloqueos): `python3 -I backend/internal/api/c1gen/build.py --check` IDÉNTICO ×3 y `uv run --no-project --with pyyaml --with jsonschema python3 -I backend/internal/api/c1gen/verify_c1.py --db <base migrada con goose>` con `RESULTADO` sin fallas; `go build ./...` y `go vet ./...` limpios. El revisor reejecutará `matriz_r2.py` (0 N-xxx sin operación ni Local) y `u1_verify.py`.
2. La asignación de la ronda 3 debe nombrar a C1 como dueño de la regeneración: en la ronda 1 nadie la tenía (`correccion-C1-r1` solo tenía H-05).
3. Los bajos (R2-03 a R2-08) se registran; conviene cerrarlos en la misma ronda porque casi todos caen en los mismos archivos.
