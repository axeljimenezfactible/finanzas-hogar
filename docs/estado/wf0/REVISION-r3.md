# Revisión WF-0 «Cimientos» · ronda 3 de 3

Revisor independiente `revisor-wf0-r3` · 11-oct-2026 · Puerta 1 (C1 + U1 + criterios; S1 y U2a como entregables de WF-0).
Evalúo **solo** contra `docs/criterios/` (S1, U1, C1 y los de U2 que aplican a la colorimetría de WF-0). No escribí nada de lo revisado y no corregí nada. Todo se reejecutó con scripts propios (scratchpad `rev-r3/`: `u2a_r3.py`, `u2a_pairs.py`, `chart_r3.py`, `matriz_r3.py`, `c103_r3.py`, `c1conv_r3.py`, `decisiones_r3.py`, `prov_r3.py`, `criterios_r3.py`, `regresion_r3.py`), además de los verificadores de los agentes. No se aprobó nada por la palabra de un reporte.

Artefactos revisados (sha256, primeros 16):
- S1: migración `747fc8fd88e576f4` · diccionario `edda526135849067` · reglas `b93aaf9bc094ccd9` (sin cambios desde r2).
- C1: `openapi.yaml` `95c8866589d173bc` · `C1-datos.md` `e3ec3c882ee4e274` · `routes.go` `b4b4d33296eb4940` (regenerados por `correccion-C1-r2`).
- U1: `45fcfb61da0ca09f` (sin cambios desde r1).
- U2a: `colorimetria.md` `beadf8708b445a60` · `contraste.csv` `9f995320c6fed1f6` · `C2-colores.json` `1446c402b3226d8e` · `colorimetria.py` `c70e643bb4300282` (sin cambios desde r2).
- Criterios: 18 archivos + INDICE sin cambios desde r2 (10-oct 23:56). Oráculo `model.json` `1a57066689ddeed3` (igual al del INDICE).

## Veredicto

**APROBADO en ronda 3**: 0 hallazgos de severidad alta, 0 de media y 10 de baja (6 heredados de r2, cuyos archivos no se tocaron en la ronda, y 4 nuevos).
- **R2-01 y R2-02 (alta, C1) quedan cerrados.** C1 se regeneró contra el esquema 1.1 de S1. `build.py --check`: IDÉNTICO en los 3 archivos. `verify_c1.py --db`: 22 OK, 0 FALLA, «TODO OK». Mi matriz propia da 0 necesidades sin operación ni Local, 0 pantallas sin operación y una sola operación sin necesidad, la excepción «estáticos» (`getOpenApiDocument`). Mi C1-03 propio da 54 de 54 tablas con operación y 0 tablas inexistentes.
- Los bajos se registran aquí (INDICE §1: «baja se registra en el informe del revisor»). Ninguno cambia cifras, esquema ni rutas, y ninguno bloquea la puerta 1.

## 1. Verificaciones ejecutadas por el revisor

| # | Qué pide la tarea | Comando o método | Resultado |
|---|---|---|---|
| 1 | Aplicar el SQL en `datos/tmp/rev.db` y correr `foreign_key_check` | Sección Up extraída con awk → `sqlite3 -bail datos/tmp/rev.db`. Down sobre una copia. Además `goose v3.28.0 -dir backend/db/migrations sqlite3 datos/tmp/rev-goose.db up/down/up` | `UP_OK`. 54 tablas, las 54 STRICT. 147 índices (110 explícitos y 37 automáticos). 68 disparadores. `PRAGMA foreign_key_check` = **0 filas** e `integrity_check` = ok. Down deja 0 objetos. goose: «OK 00001_esquema.sql»; el down deja solo `goose_db_version` y `sqlite_sequence`, y el segundo up da 56 (54 + 2).<br>Semillas: 7 `app_setting` (`price.stale_after_days` 30, `auth.*` 5 y 15, pesos iguales del optimizador, `ai.*` false, lectura asistida false), 4 `meal_slot` (colación opcional) y 7 `notification` (semanal y plan sin armar en domingo, mensual el día 1, pago [1, 0], presupuesto 80/100).<br>Negativas rechazadas: cuenta de Karen (disparador), hash que no es argon2id, segunda co-titular (UNIQUE), `12.5` en centavos (STRICT), fecha `10-oct-2026`, ventana del 10 al 9 como mes, `pay_date` 2027-02-30 y pasar a externa a una persona con cuenta. La quincena 2 de feb-2027 paga el 2027-02-28. `datos/tmp` borrada al final |
| 2 | `go build ./...` y `go vet ./...` | En `backend/` con go1.26.5; además `gofmt -l .`, `go list ./...` y un cruce propio de `routes.go` con `openapi.yaml` | BUILD_OK y VET_OK; `gofmt` sin salida; 2 paquetes. `routes.go` = `openapi.yaml`: 197 = 197 en método, ruta, `operationId` y pública (S6-03). `go.mod` declara `go 1.26` (bloqueo `S6-version-go`; S6-01 es de WF-3) |
| 3 | `colorimetria.py` regenera `contraste.csv` sin `cumple=false`; hex de C2 = `colorimetria.md` | El script del repositorio, copiado a un árbol espejo del scratchpad, corrido con `--sin-png` (no reescribe el repositorio) y comparado con `cmp`. Además `u2a_r3.py`, con WCAG y OKLab escritos aparte, sin importar el script | «VERIFICACIÓN OK». md, csv, json y html **idénticos byte a byte** al repositorio.<br>Recálculo propio: 202 filas, 0 `cumple≠true`, 0 bajo umbral (mínimo de texto 4.73, de UI 3.05) y 0 razones que difieran en más de 0.005. Hex del CSV = hex del token.<br>190 tokens del JSON = 190 filas del Anexo A, 0 hex distintos. Las tablas de valores finales §5.1, §5.2, §5.3, §6.1, §7.1, §7.2, §8.1 y §10 usan solo hex de token. Plan §4.3 (marca y neutro 50–950): 22 de 22 |
| 4 | Matriz C1: toda N-xxx con endpoint y todo endpoint con pantalla o justificación | `matriz_r3.py` propio: U1 §3 y §3.29 ↔ `x-necesidades` ↔ C1 §10. Además `verify_c1.py --db` y `build.py --check` de C1, `u1_verify.py` (extraído del anexo del reporte de U1) y `openapi-spec-validator` 0.9.0 | **U1**: 199 N-xxx (136 L, 63 E), 0 duplicadas, Local N-003 y N-301; 41 pantallas y 220 pares en §3.29.<br>**OpenAPI**: 3.1.1, 158 rutas, 197 operaciones.<br>**Cobertura**: N-xxx sin operación ni Local = **0**; pantallas sin operación ni Local = **0**; operaciones sin N-xxx = solo `getOpenApiDocument` (C1 §11.2); N-xxx citadas que no existen en U1 = 0.<br>**C1 §10**: 220 filas = 220 pares de U1, 0 rutas inexistentes; cada operación citada lista su N-xxx en `x-necesidades` (0 incoherencias). Estados: OK 198, Local 2, PROVISIONAL 20, BLOQUEO 0.<br>**Lista de compras**: N-151 y N-153 → `updateShoppingListItem`; N-152 → `createShoppingListItem`; N-165 → `assignShoppingListStore`; N-127 → `generateShoppingList`; N-150 y N-073 → `getShoppingList`.<br>**Herramientas de los agentes**: `build.py --check` IDÉNTICO ×3; `verify_c1.py` «TODO OK»; validador «OK», salida 0. `u1_verify.py` da «OK U1-03 N-xxx sin endpoint en C1: []» y una FALLA por `GET /openapi.yaml`, que es la excepción de C1-02 (R3-03) |
| 5 | Cada decisión A–I cubierta por S1, U1 o C1 | `decisiones_r3.py`: 35 renglones comprobados contra la base, la OpenAPI y U1, más lectura dirigida | 35 de 35 con respaldo (tabla de §3) |
| 6 | PROVISIONAL, supuestos y bloqueos: ¿esconden un valor inventado? | `prov_r3.py` sobre migración, diccionario, reglas, C1, U1, colorimetría y `main.go`, y sobre los nodos `x-provisional` de la OpenAPI. Lectura de los 32 bloqueos y de las secciones Supuestos de los 10 reportes | **Marcas en texto**: 108; las de contenido citan un bloqueo existente a ±3 líneas. Las 20 que no lo citan cerca son leyendas, renglones de índice cuyo detalle sí lo cita, o la retención de idempotencia (`C1-parametros-tecnicos`).<br>**OpenAPI**: 54 nodos `x-provisional: true`, todos con un bloqueo existente en su texto; 78 apariciones de «PROVISIONAL», 0 sin bloqueo cerca.<br>**Cifras de C1 §9**: 18 distintas; 16 están en `model.json` y 2 en `data.json` (3,233.00 = `totales.walmart`; 138.50 = 3,233.00 − 3,094.50).<br>Ningún número inventado entra en cifras, esquema ni contrato. Quedan 4 límites técnicos de validación sin marca (R3-02) |
| 7 | Criterios medibles y trazables | `criterios_r3.py` sobre los 18 archivos | 248 criterios (189 alta y 59 media), igual que el INDICE. Los 248 tienen «Verifica» y «Evidencia» y aparecen en la §5 de su archivo (contando rangos); los 18 archivos tienen las 5 secciones. Sin cambios desde r2: siguen las incoherencias menores de R2-07 |
| 8 | Registrar gasto en ≤ 3 toques demostrado en U1 | Lectura de U1 §4 (regla de conteo) y F-1; `ExpenseCreate` en la OpenAPI; `u1_verify.py` | **Demostrado.** «+» (1), chip de categoría (2) y «Guardar» (3). Teclear no cuenta. Fecha, persona, forma de pago y tipo tienen valor por defecto (N-052) y el reparto sale de la categoría.<br>`ExpenseCreate` exige solo `category_id`, `amount_cents`, `expense_date` y `payment_method`, todos cubiertos en esos 3 toques. La descripción es opcional en S1 (NULL) y en C1 (sin relleno).<br>`u1_verify.py`: «U1-04 … = 3». La prueba en Android sigue a cargo de U5 |
| 9 | Regresión de S1 (Kueski 389500 y 239528) | `regresion_r3.py` propio sobre una copia de la base: 3 personas; Kueski (revolvente, límite 1000000, lectura B); 2 tasas por plazo; 4 despensas como gasto (devengo de 2 meses) y compra a crédito (intereses cotizados con IVA) con sus partes; 16 pagos con 32 partes; 2 papeles de contado de Karen. Lectura de vuelta desde la base | Suma de pagos de la despensa 1 = **389500**. Parte de Axel = **239528**, en la compra y como Σ de sus pagos. 98 valores comparados, **0 diferencias**:<br>- 4 despensas con partes.<br>- 16 pagos con partes y Σ partes = monto.<br>- 17 eventos con disponible A y B.<br>- Series A y B de 14 valores; mínimos 809450 y 500000.<br>- Efectivo y devengado de oct-2026 a mar-2027 = `periodos[0..5]`.<br>- 0 eventos entre los días 1 y 9.<br>`fk_check` 0, `integrity_check` ok |

## 2. Tabla de criterios

Estados: **Cumple** · **Cumple con observación** · **No cumple** · **No aplica en WF-0**.

### S1 · Dominio y datos (esquema 1.1, sin cambios desde r2)
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| S1-01 Entidades completas | alta | Cumple | 31 de 31 entidades del plan en el diccionario. 54 tablas en la base, las 54 con sección; columnas idénticas en nombre, orden y tipo |
| S1-02 Migraciones en SQLite limpio | alta | Cumple | sqlite3: Up ok, `integrity_check` ok, `foreign_key_check` vacío. goose v3.28.0 up/down/up ok |
| S1-03 Cada regla con ejemplo del reporte | alta | Cumple con observación | Sin cambios desde r2 (22 de 22 reglas con ejemplo). Las 10 reglas sin datos en el oráculo llevan ejemplo «ilustrativo», caso que el criterio no prevé (R2-07) |
| S1-04 Ninguna tabla huérfana | alta | Cumple | 54 de 54 con «Usado por». Muestra del revisor contra U1: `notification_subscription` → P-26.09; `inventory_item` → P-16 y P-26.14; `weight_entry` → P-21 y P-22; `income` → P-08 y P-26.03; `backup` → P-26.11; `expense_share` → S2, S7 y filtros |
| S1-05 Dinero exacto | alta | Cumple | 36 columnas `_cents`, todas INTEGER; 0 REAL, FLOAT, DOUBLE ni NUMERIC; `12.5` rechazado |
| S1-06 Fechas y zona | media | Cumple | `CHECK (date(x) IS x)`; `10-oct-2026` rechazado; zona declarada |
| S1-07 Personas y acceso | alta | Cumple | Cuenta de Karen, segunda co-titular y paso a externa con cuenta rechazados. Sin columna de rol. Ingreso con disparador de solo titular |
| S1-08 Reparto y participantes | alta | Cumple | Σ partes = total en las 4 despensas y en los 16 pagos (regresión). Partes por persona en `credit_purchase_share` |
| S1-09 Periodos e ingresos | alta | Cumple | `CHECK` de mes calendario y quincenas; Q2 de feb-2027 paga el 28; ventana del 10 al 9 rechazada; equivalencia de 12 de 12 cifras por mes calendario; 0 eventos entre los días 1 y 9 |
| S1-10 Créditos | alta | Cumple | 4 tipos con CHECK; tasas 53450/309450 y 80050/309450 cargadas; `availability_reading DEFAULT 'B'` y `availability_confirmed`; promoción con comercio, plazo, mínimo y vigencia |
| S1-11 Comida, salud, inventario, precios | media | Cumple | Sin cambios desde r2: 12 de 12 |
| S1-12 Ida y vuelta de la semilla | alta | Cumple | Regresión propia: 98 valores, 0 diferencias; 389500 y 239528 |
| S1-13 Ajustes tiene dónde guardarse | media | Cumple | Tabla de 14 secciones en el diccionario (sin cambios desde r2); Apariencia local, igual que U1 N-301 |
| S1-14 Sesiones y credenciales | alta | Cumple | Columnas sensibles: solo `password_hash` (CHECK argon2id) y `token_hash`. 0 tokens de Telegram (`grep "[0-9]{8,}:AA"` = 0) |
| S1-15 Integridad y consultas de C1 | media | Cumple | FK (muestra del revisor): `shopping_list_item` 4, `scheduled_payment_share` 2, `meal_plan_entry` 3, `price` 4, `savings_movement` 1. Índices con `person_id`, `category_id`, `store_id` y `credit_id`, y con las fechas de periodo (`expense_date`, `due_date`) |
| S1-16 Sin valores inventados | alta | Cumple con observación | Las semillas salen de las decisiones y toda marca PROVISIONAL cita un bloqueo. Siguen la redacción del tope de Karen (R2-03) y dos límites técnicos sin marca (R3-02) |
| S1-17 Idioma y nombres | media | Cumple | Todos los identificadores en `[a-z0-9_]`; textos en es-MX |

### U1 · Arquitectura de información (documento sin cambios desde r1)
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| U1-01 Barra lateral de 9 entradas | alta | Cumple | Sin cambios: 9 entradas en el orden exacto |
| U1-02 Inventario con identificador | alta | Cumple | 27 P-xx + 14 P-26.xx = T 41; 199 N-xxx únicas, cada una con tipo L o E |
| U1-03 Matriz pantalla ↔ necesidad ↔ endpoint | alta | **Cumple** (antes No cumple) | 0 N-xxx sin operación ni Local y 0 pantallas sin operación. La única operación sin N-xxx es la excepción «estáticos». El reporte de U1 sigue diciendo «PENDIENTE» y su script marca esa excepción como falla (R3-03) |
| U1-04 Gasto en ≤ 3 toques | alta | Cumple | Regla escrita y tabla de F-1: 3 toques; campos obligatorios cubiertos por valores por defecto |
| U1-05 Seis flujos clave | alta | Cumple con observación | 8 flujos (F-1 a F-6 y los complementarios F-7 y F-8), todos los pasos con pantalla existente. Los 4 obligatorios están en F-1, F-2, F-7 y F-8; el criterio dice «6» (R2-07) |
| U1-06 Vistas por entidad | alta | Cumple | 31 entidades, coherentes con C1 §7 (23 con vistas y 8 con motivo) |
| U1-07 Menú «/», atajos, edición | media | Cumple | 27 comandos con destino; 15 atajos, 0 duplicados |
| U1-08 Ajustes con 14 secciones | alta | Cumple | 14 de 14; C1 §11.5 con 14 filas |
| U1-09 Cobertura de las decisiones | alta | Cumple | Sin cambios (26 capacidades); cruce de §3 sin faltantes |
| U1-10 Dos titulares, Karen sin acceso | alta | Cumple | Sin pantalla de acceso de Karen; 56 menciones, todas en las áreas permitidas |
| U1-11 Móvil primero, una acción primaria | media | Cumple con observación | §6.4 completa; §6.1 punto 4 sigue en 24 × 24 px frente a ≥ 44 px (R2-05) |
| U1-12 Consistencia con S1 y C1 | media | Cumple con observación | C1 ya coincide con S1 1.1 (R2-02 cerrado). Diferencia menor en la semana por omisión de la lista de compras (R3-01) |
| U1-13 Idioma | media | Cumple | 77 nombres visibles, 0 en inglés |
| U1-14 Supuestos y bloqueos | media | Cumple | Sección Supuestos; 2 bloqueos U1; PROVISIONAL marcado |

### C1 · Contrato de datos y API (regenerado en la corrección r2)
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| C1-01 OpenAPI 3.1 válido | alta | Cumple | `openapi-spec-validator` 0.9.0 → OK, salida 0; 3.1.1, 158 rutas, 197 operaciones, `info.version` 1.0.0 |
| C1-02 Toda N-xxx con operación y viceversa | alta | **Cumple** (antes No cumple) | Matriz propia: 0 N-xxx sin operación ni Local, 0 pantallas sin operación, 1 excepción justificada (§11.2), §10 = 220 pares, 0 filas BLOQUEO. `verify_c1.py`: OK |
| C1-03 Endpoints respaldados por tablas | alta | **Cumple** (antes No cumple) | Propio: 54 de 54 tablas citadas en `x-tablas`, 0 inexistentes. 95 esquemas con `x-tabla`: toda propiedad es columna, campo calculado o de entrada, o lista u objeto anidado de un esquema hijo cuyas propiedades son columnas (`CostShareRulePartInput`, `CreditPurchaseShareInput` y los demás). `verify_c1.py`: 967 propiedades, 0 sin columna; 8 `x-columna`, todas existen |
| C1-04 Una convención de errores | alta | Cumple | 895 de 895 respuestas 4xx, 5xx y default → `components/responses` → `Error` |
| C1-05 Una convención de paginación | alta | Cumple | 39 de 39 listas con `Cursor`, `Limit` y `Sort` y la envoltura `Page` |
| C1-06 Una convención de filtros | alta | Cumple | 6 de 6 en `components/parameters` (`person_id`, `family`, `category_id`, `store_id`, `period`, `credit_id`), 0 sinónimos; `family` PROVISIONAL con bloqueo |
| C1-07 Una convención de fechas | alta | Cumple con observación | Solo `date` y `date-time`. El patrón de `DateTime` admite «Z» (R2-06) |
| C1-08 Una convención de dinero | alta | Cumple | 256 propiedades `_cents` → `MoneyCents` (integer); 0 `type: number` |
| C1-09 Métricas con fórmula y periodo | alta | Cumple | M-01 a M-17 con fórmula, periodo, filtros, tablas y endpoint. Ejemplos 10,000.00 / 6,905.50, 973.75 y 1,947.50 presentes y del oráculo |
| C1-10 Recursos y vistas de las 31 entidades | alta | Cumple | §7 con 31 filas, coherente con U1 §2.5 (lista de compras y movimientos de inventario dentro de sus entidades) |
| C1-11 Cobertura funcional | alta | Cumple | §11.6 con 35 capacidades, todas con `operationId` existente; las escrituras de la lista de compras ya existen; §11.5 con las 14 secciones |
| C1-12 Seguridad | alta | Cumple | Públicas solo `getHealth` y `login`. Cookie `fh_session` HttpOnly, Secure y SameSite (PROVISIONAL). Propiedades sensibles solo de entrada (`LoginRequest.password`, `PasswordChangeRequest`). 429 en el login |
| C1-13 Versionado y estructura | media | Cumple | 197 `operationId` únicas; 9 etiquetas de la barra; «C1 v1.0» con las filas «corrección r1» y «corrección r2» |
| C1-14 Idioma y ejemplos | media | Cumple | `verify_c1.py`: ejemplos válidos con jsonschema y descripciones en español |
| C1-15 PROVISIONAL marcado | alta | Cumple con observación | 54 nodos `x-provisional` y 42 líneas de C1, todos con un bloqueo existente. `C1-lista-compras` y `C1-huecos-u1-s1` plantean ya la misma pregunta que sus pares de S1. Quedan 2 límites técnicos sin marca (R3-02) y un detalle de la DateTime (R2-06) |

### U2 · aplicados a la colorimetría de WF-0 (U2a, sin cambios desde r2)
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| U2-01 DESIGN.md con 10 secciones | alta | No aplica en WF-0 | DESIGN.md es de WF-1 |
| U2-02 Matiz C aplicado | alta | Cumple | Matiz medido con OKLab propio, en pasos con croma > 0.02: marca 279.0–280.6°, éxito 149.1–150.4°, error 19.5–20.3°, advertencia 73.8–75.7°, información 238.7–240.1°, ingreso 134.5–135.5° (la propuesta de 135° que admiten las decisiones I), gasto 23.8–25.1°, ahorro 199.8–200.9°, deuda 309.4–310.2° y acción crítica 99.5–100.8°. Marca y neutro 50/950 = plan. 0 tokens fuera de `color.critical` entre 80° y 120° |
| U2-03 Método y saturación | alta | Cumple | Desviación máxima de L por paso: 0.0016 (redondeo a hex). Croma máximo: marca 0.1103, neutro 0.0121, roles ≤ 0.1408 |
| U2-04 0 pares bajo umbral | alta | Cumple | 202 pares, 0 bajo umbral; recálculo propio de los 202 sin diferencias |
| U2-05 Pares estructurales | alta | Cumple | Recalculados: 11.53 (difiere 0.13 del plan; documentado con ambos valores), 5.39, 3.46, 3.98, 14.06, 7.97, 10.37, 1.19 y 1.89; borde de campo 3.71 documentado |
| U2-06 Roles como texto | alta | Cumple | Paso 700 en claro y 400 en oscuro, ≥ 4.5 sobre fondo y tarjeta en los 8 roles (mínimo 4.73) |
| U2-07 Botón primario | alta | Cumple | Blanco / marca 700 = 5.77; marca 700 / fondo = 5.53; marca 950 / marca 300 = 10.47; marca 300 / fondo 950 = 10.37 |
| U2-08 Colisiones resueltas | alta | Cumple | Pares con menos de 20° y ΔL < 0.15 entre roles y categóricas: 3 en claro y 4 en oscuro, los mismos de §7.7 con su refuerzo. Karen pasa a 50° y Bodega Aurrera a 50° (lejos de advertencia 75°). Fuera del alcance del criterio, la marca queda cerca de Walmart y Conjunto (R3-04) |
| U2-09 Paleta de 8 series | alta | Cumple | ΔL consecutivo mínimo (OKLab propio): 0.174 en claro y 0.218 en oscuro. Contraste mínimo contra fondo y tarjeta: 3.15 y 3.11. Los 28 pares dan 0.0 por diseño de dos niveles (bloqueo `criterios-deltal`). Simulación documentada |
| U2-10 a U2-13 | alta/media | No aplica en WF-0 | WF-1 |
| U2-14 Idioma y supuestos | media | Cumple | es-MX, tokens en inglés, Supuestos y bloqueo `U2-gasto-error` |

### Criterios (`docs/criterios/`, sin cambios desde r2)
| Aspecto | Estado | Evidencia |
|---|---|---|
| Completos y con estructura | Cumple | 18 archivos + INDICE; 248 criterios = INDICE (189 / 59); secciones 1 a 5 en todos |
| Medibles | Cumple | 248 de 248 con «Verifica» (comando, script o cálculo) y «Evidencia» |
| Trazables al plan | Cumple con observación | 248 de 248 en la §5 de su archivo, contando rangos. Siguen las paráfrasis entre comillas (R2-07) |
| Coherentes con entregables y decisiones | Cumple con observación | Siguen las incoherencias menores de R2-07: U2-02 con 135°, U1-05 con «6», orden de comandos de C1-01, S1-03 con ejemplos ilustrativos, rutas `/tmp` en los «Verifica» y `verify_c1.py` sin citar |

## 3. Cobertura de las decisiones A a I

| Id | Decisión (resumen) | S1 (esquema o regla) | U1 (pantalla o necesidad) | C1 (operación) | Estado |
|---|---|---|---|---|---|
| A1 | Co-titulares con los mismos permisos; usuario y contraseña; sesión larga; cierre desde Ajustes; argon2id; 5 en 15 min | `user_account` sin rol, CHECK argon2id, `ux_person_single_co_holder`; `device_session.expires_at`; `login_attempt`; semilla `auth.*` 5 y 15 | P-00, P-26.01, N-010, N-272 | `login` (429), `logout`, `listSessions`, `revokeSession(s)`, `changeOwnPassword` | Cubierta (duración PROVISIONAL, `criterios-login-bloqueo`) |
| A2 | Karen externa sin cuenta; despensa y compras con Kueski; paga directo; parte informativa; estado de cuenta PDF | `participation = external`, `pays_share_directly`, `trg_user_account_not_external_*`, partes por persona, `generated_report` | P-25, N-261, N-309 | `getPersonParticipation`, `getExternalStatement`, `generateExternalStatement`; sin alta de cuentas | Cubierta (reposición por parte PROVISIONAL) |
| A3 | Comidas y salud solo de Axel y Ximena | Disparadores en `meal_plan_entry`, `health_profile` y `weight_entry` | P-10, P-21 | `getMealPlanWeek`, `listHealthProfiles` | Cubierta |
| A4 | Telegram a Axel y Ximena, destinatarios por aviso | `notification_recipient.chat_id`, `notification_subscription` | P-26.09, N-299, N-300 | `listNotificationRecipients`, `setNotificationRecipients` | Cubierta (registro del chat: `S1-huecos-negocio`) |
| B1 | MXN, centavos, HALF_UP, el último pago ajusta; ISO; zona; es-MX | 36 `_cents` INTEGER STRICT; CHECK ISO; RN-01 y RN-07 | §0.2 y §0.4 | `MoneyCents`, `Date`, `DateTime` | Cubierta |
| B2 | Mes calendario; quincena 15/30; despensa devengada a la mitad | `period` con CHECK; `expense.accrual_months`; `budget` con vigencia | P-06, N-085 | `getBudgetMonth`, `getBudgetFortnights`, `getAccruedGroceries` | Cubierta (rangos de quincena PROVISIONAL) |
| B3 | Ingresos solo de Axel, 15/30, más extras | `income` con disparador de solo titular; `planned_item.second_half_amount_cents` | P-08, N-101, N-279 | `listIncomes`, `createIncome`, `planned-items` | Cubierta |
| B4 | Reparto por porcentaje, por gasto o categoría, 50/50; qué se comparte con Karen | `cost_share_rule(_part)`, regla de la categoría y del gasto, `credit_purchase_share` | N-051, N-281, N-309 | `listCostShareRules`, `updateCostShareRule`, filtro `family` | Cubierta (`family` PROVISIONAL) |
| B5 | Fondo y metas con fecha | `savings_goal.target_date`, `planned_monthly_cents`, `savings_movement` | P-09, N-113 | `listSavingsGoals`, `createSavingsGoal` | Cubierta |
| B6 | Comida fuera como gasto con persona; «fuera» en el plan | `meal_plan_entry.is_eating_out`, `expense.person_id` | N-124 | `createMealPlanEntries` | Cubierta |
| C1 | 4 tipos de crédito | `credit.kind` con CHECK de 4 valores | P-17, P-18, P-26.05 | `listCredits`, `createCredit` | Cubierta |
| C2 | Tasa por compra con IVA; tabla por plazo 17.27 / 25.8685; promociones | `interest_includes_vat`, `credit_term_rate` como fracción, `promotion.min_purchase_cents` | N-218, N-286, N-287 | `listCreditTermRates`, `listPromotions` | Cubierta (decimales de la tasa en pantalla: R2-06) |
| C3 | A/B por defecto B, aviso «pendiente confirmar» | `availability_reading DEFAULT 'B'`, `availability_confirmed` | N-284, N-285 | `previewAvailability` | Cubierta |
| C4 | Comparador: 3 restricciones, 4 objetivos, recomendación, sin solver | Semilla `optimizer.weights` con pesos iguales; restricciones sin valor | P-20, F-7 | `compareScenarios` | Cubierta (valores PROVISIONAL, `S1-restricciones-optimizador`) |
| C5 | Regresión Kueski | Regresión propia: 98 valores, 0 diferencias | — | Ejemplos del oráculo en §9 | Cubierta en WF-0 (los motores son de WF-1) |
| D1 | Recetario de Claude más las propias; restricciones vacías | `recipe.origin`, `health_profile.dietary_restrictions`, `recipe_meal_slot` | P-11, P-12 | `listRecipes`, `createRecipe` (`meal_slot_ids`) | Cubierta |
| D2 | 4 tiempos; porciones base 2 | `meal_slot` (4, colación opcional), `base_servings DEFAULT 2` | P-10 | `listMealSlots` | Cubierta |
| D3 | Perfil, IMC, meta de kcal, kcal por receta y día; macros después | `health_profile`, `weight_entry`, `nutrition_fact` | P-21, P-22 | `getTodayKcal` | Cubierta (meta de kcal PROVISIONAL) |
| D4 | Fuentes nutricionales autorizadas | `nutrition_fact.source` | P-26.07 | `reloadNutritionBase` | Cubierta |
| D5 | Inventario híbrido; sobrantes; mercados con día de visita; lista de compras | `inventory_item`, `inventory_movement.reason`, `store.visit_weekday`, `shopping_list_item.is_bought` | P-13, P-16, P-26.06 | `getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `updateShoppingListItem`, `assignShoppingListStore`, `adjustInventoryItem`, `listInventoryMovements` | Cubierta (lista PROVISIONAL, opción A de `C1-lista-compras`; R2-01 cerrado) |
| E1 | Manual, OCR con confirmación, lectura asistida apagada, PROFECO | `price.source`, `store.assisted_reading_enabled`, semilla `price.assisted_reading_enabled = false` | P-05, F-8 | `createTicketScan`, `confirmTicketScan`, `importProfeco` | Cubierta |
| E2 | Viejo a los 30 días y pesa menos | Semilla `price.stale_after_days = 30` | N-170 | Frescura en precios | Cubierta (peso PROVISIONAL) |
| E3 | Sam's: punto de equilibrio | `store.requires_membership`, `membership_annual_fee_cents` | N-164 | `getBasketComparison` | Cubierta (cálculo PROVISIONAL) |
| E4 | Unidad base, sustitutos, mejor tienda | `product.base_unit`, `presentation` | N-161, N-163 | `listBaskets`, `getBasketComparison` | Cubierta |
| F1 | Diario, semanal, mensual PDF el día 1, vista web; estado de cuenta de Karen | Semillas de `notification` (diario, semanal en domingo, mensual el día 1), `job`, `generated_report` | P-23 a P-25, P-26.08 | Reportes, `generateExternalStatement` | Cubierta (hora del diario: bloqueo) |
| F2 | 4 alertas | Semillas `payment_due` [1, 0], `budget_threshold` 80/100, `stale_prices`, `meal_plan_missing` | P-26.09 | `getBudgetAlerts`, `notifications` | Cubierta |
| F3 | MultiX en texto plano, token en tiempo de ejecución | Sin token en la base | N-299 | 0 campos sensibles en respuestas | Cubierta |
| G1 | Fiel a Notion: barra, vistas, «/», atajos, edición en el lugar | — | §1, §2.5, §2.6, §5 | `queryView` (4 modos) | Cubierta |
| G2 | Android y escritorio; PWA con mkcert; ≤ 3 toques | `idempotency_key` para reintentos | F-1, §6 | `getCaCertificate` | Cubierta (PWA en WF-3) |
| G3 | Tema del sistema con interruptor | Local por dispositivo | N-301 (Local) | Sin API | Cubierta |
| G4 · I2 | Paleta nueva; matiz C y sus reglas | — | — | U2a: 190 tokens | Cubierta |
| G5 | Resumen ≤ 6 tarjetas; Detalle, Tendencias, Calendario | — | P-01 | M-01 a M-07 | Cubierta |
| G6 | Ajustes con 14 secciones | Tabla de 14 secciones → tabla.columna | P-26.01 a P-26.14 | C1 §11.5 | Cubierta |
| G7 | Autoalojado; cifras tabulares; WCAG 2.2 AA | — | §0.4 | U2a: contraste | Parcial en WF-0: tipografía e iconos son de WF-1; táctil (R2-05) |
| H1 | Solo LAN, `finanzas.local:8443` | — | — | `servers` = `https://finanzas.local:8443`; el esqueleto escucha en 127.0.0.1 (`criterios-bind`) | Cubierta |
| H2 | Go + SQLite + systemd; respaldos locales | `backup.integrity_check`, `restore_tested_at` | P-26.11 | `listBackups`, `createBackup`, `restoreBackup` | Cubierta (systemd en WF-4) |
| H3 | IA opcional y apagada | Semillas `ai.* = false` | N-306 | `updateAiProvider` | Cubierta |
| H5 | Semilla e importador | `planned_item` | P-26.11 | `importPlannedItems` | Cubierta |
| I1 | Pila A | goose v3.28.0 validado; Go `net/http` sin dependencias | — | `routes.go` | Cubierta (Go 1.27: bloqueo `S6-version-go`) |
| H4, H6, J | Proceso (paquetes, repositorio, puertas) | — | — | — | No aplica al producto |

## 4. Bloqueos, supuestos y PROVISIONAL

- Hay 32 bloqueos en `docs/estado/bloqueos/`, todos con el formato del protocolo: 16 de criterios (uno ya resuelto), 7 de S1, 2 de U1, 1 de U2, 5 de C1 y 1 de S6.
- Los pares de C1 y S1 ya plantean una sola pregunta: `C1-lista-compras` = `S1-lista-compras`; `C1-huecos-u1-s1` = `S1-huecos-negocio` (3 puntos). Detalle menor: `S1-huecos-negocio` punto 1 lista también N-275 (lectura), que C1 ya no marca PROVISIONAL; la pregunta es la misma.
- **Ningún valor inventado entra en cifras, esquema ni contrato.**
  - Las semillas de la migración salen de las decisiones A, C, D, E, F y H.
  - Los valores técnicos de C1 llevan marca y bloqueo: 200, 24 h, `SameSite=Strict`, 12 caracteres, 1 GB y 90 días de sesión (`C1-parametros-tecnicos`, `criterios-login-bloqueo`).
  - Las definiciones de los campos calculados llevan `C1-campos-calculados`.
  - Las 18 cifras de C1 §9 salen de `model.json` o de `data.json`.
- Quedan 4 límites de validación sin marca ni bloqueo, todos técnicos y sin efecto en cifras (R3-02).
- **Supuestos de los reportes**: ninguno cambia cifras, esquema, endpoints ni tokens.
  - Los de `correccion-C1-r2` describen la forma de la API de la lista de compras (cubierta por `C1-lista-compras`), el código 409 `budget_overlap` (tabla única de códigos) y el movimiento inicial de inventario.
  - El supuesto 4 de U1 («los domingos se abre la semana siguiente») no quedó reflejado en el valor por omisión de C1 (R3-01).

## 5. Estado de los hallazgos de la ronda 2

| Id r2 | Sev. | Estado en r3 | Evidencia |
|---|---|---|---|
| R2-01 | alta | **Cerrado** | 5 operaciones de la lista de compras (N-127, N-150 a N-153, N-165, N-073), PROVISIONALES con bloqueo; matriz propia sin faltantes; `verify_c1.py` y `u1_verify.py` con U1-03 en OK |
| R2-02 | alta | **Cerrado** | `build.py --check` IDÉNTICO ×3. `Budget` con vigencia, `meal_slot_ids` sobre `recipe_meal_slot`, descripción opcional sin relleno, `InventoryMovement` = `inventory_movement`, `second_half_amount_cents` y `idempotency_key` declarada. En `openapi.yaml`: 0 apariciones de `period_id`, `default_meal_slot_id` y `37493`; C1-03 propio sin faltantes |
| R2-03 | baja | Abierto (S1 sin cambios) | La descripción de `person.fortnight_payment_cap_cents` sigue afirmando que el oráculo no trae un tope |
| R2-04 | baja | Abierto | `reglas-negocio.md` sigue en «Versión: 1.0» |
| R2-05 | baja | Abierto | U1 §6.1 punto 4: 24 × 24 px |
| R2-06 | baja | Abierto | `DateTime` admite «Z»; `rate_pct` con 4 decimales. `correccion-C1-r2` lo dejó en Pendientes 1 y sugiere preguntar a Axel cuántos decimales se muestran |
| R2-07 | baja | Abierto | Criterios sin cambios |
| R2-08 | baja | Abierto | `gen_dict.py`, `s1_verifica.py`, `u1_verify.py` y los scripts de criterios siguen solo en anexos de reportes |

## 6. Hallazgos de la ronda 3

Todos son de severidad **baja**: se registran y no bloquean la puerta 1.

| Id | Sev. | Artefacto | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|---|
| R3-01 | baja | C1 | U1-12, C1-02 | Las operaciones de la lista de compras toman por omisión «la semana en curso». U1 §0.2 define la «lista de compras vigente» como la semana en curso y, los domingos, la siguiente (supuesto 4 del reporte de U1); N-150 y N-165 usan «lista vigente». Si el cliente no manda `period` un domingo, «Usar esta tienda» y «Agregar artículo» actúan sobre otra semana que la que muestra P-13 | `docs/contratos/openapi.yaml` (`getShoppingList`, `generateShoppingList`, `createShoppingListItem`, `assignShoppingListStore`); `backend/internal/api/c1gen/gen3_ops.py` | En la próxima regeneración (tras la puerta 1): igualar el valor por omisión a la regla de U1, o volver obligatorio `period` en las escrituras y documentarlo |
| R3-02 | baja | C1 | C1-15, S1-16 | Límites técnicos de validación sin marca PROVISIONAL ni bloqueo: `Cursor` `maxLength` 512 y texto de búsqueda `maxLength` 200 en la OpenAPI; usuario de 3 a 64 caracteres y clave de idempotencia de 16 a 128 en CHECK de la migración (C1 §1 punto 8 la declara convención). No afectan cifras ni rutas, pero C1-15 pide marcar todo valor que no esté en el plan | `docs/contratos/openapi.yaml` (`components/parameters`); `backend/db/migrations/00001_esquema.sql` (`user_account.username`, `idempotency_key.idempotency_key`); `docs/estado/bloqueos/C1-parametros-tecnicos.md` | Agregarlos a `C1-parametros-tecnicos`, o declararlos una sola vez como convención técnica en C1 §1 y en `reglas-negocio.md` §Convenciones |
| R3-03 | baja | U1 | U1-03, U1-14 (INDICE §1, reejecución) | El reporte de U1 sigue con «U1-03: PENDIENTE» (líneas 35 y 198) y su Pendiente 1 ya está resuelto. Su `u1_verify.py` termina con «1 falla(s)» por `GET /openapi.yaml`, que es la excepción «estáticos» que admiten C1-02 y U1-03 (C1 §11.2). Un lector del reporte vería un criterio alto pendiente que ya cumple | `docs/estado/wf0/U1-arquitectura.md` (§Verificación, §Pendientes y Anexo) | El integrador anota en el reporte el cierre de U1-03 (salida de esta ronda) y ajusta el script para aceptar las excepciones declaradas en C1 §11.2 (o lo incorpora al repositorio, R2-08) |
| R3-04 | baja | U2a | U2-08 (adyacente); decisiones G | La revisión de pares cercanos de §7.7 cubre roles, acción crítica y categóricas, pero no la marca. En claro, Walmart (262°) queda a 18.0° de la marca 700 con ΔL 0.000, y Conjunto (279°, croma bajo) a 0.9° con ΔL 0.094; en oscuro, Walmart queda a 17.8° de la marca 400. El criterio no lo exige (la marca no es un rol), pero un chip de Walmart junto a un botón o enlace de marca puede confundirse | `docs/design/colorimetria.md` §7.7; `docs/design/tools/colorimetria.py` | En WF-1 (DESIGN.md): agregar la marca al barrido de §7.7 y declarar el refuerzo (etiqueta o forma), o separar Walmart de la marca |
| R2-03 | baja | S1 | S1-16 | Heredado sin cambios (ver §5) | `backend/db/migrations/00001_esquema.sql` (`person.fortnight_payment_cap_cents`); `docs/contratos/diccionario-datos.md` | Redacción neutral y la pregunta como bloqueo (p. ej. en `S1-restricciones-optimizador`): «¿el tope por quincena de Karen es 374.93 (nota «cap» del oráculo) o lo fija Axel?» |
| R2-04 | baja | S1 | Trazabilidad de cambios | Heredado sin cambios | `docs/contratos/reglas-negocio.md` (encabezado) | Subir a 1.1 con su lista de cambios |
| R2-05 | baja | U1 | U1-11 (U2-11) | Heredado sin cambios | `docs/design/U1-arquitectura.md` §6.1 punto 4 | ≥ 44 px como mínimo táctil; 24 px solo como piso WCAG 2.5.8 |
| R2-06 | baja | C1 | C1-07; decisión C | Heredado sin cambios | `docs/contratos/openapi.yaml` (`DateTime`, `CreditTermRate.rate_pct`); `docs/contratos/C1-datos.md` §4 | Documentar que «Z» solo se acepta en entradas; preguntar a Axel en la puerta 1 cuántos decimales de tasa se muestran (17.27 frente a 25.8685) |
| R2-07 | baja | criterios | U2-02, U1-05, C1-01, C1-02, C1-03, S1-03, S6-03; protocolo §7 | Heredado sin cambios | `docs/criterios/U2.md`, `U1.md`, `C1.md`, `S1.md`, `S2.md` a `S7.md`, `C2.md` | Admitir 135° como propuesta; «≥ 6 flujos con los 4 obligatorios»; poner primero el comando `uvx`; permitir ejemplos «ilustrativos» marcados; usar `datos/tmp` en lugar de `/tmp`; citar `verify_c1.py --db` y `build.py --check` |
| R2-08 | baja | S1 | INDICE §1 (reejecución) | Heredado sin cambios | `docs/estado/wf0/S1-modelo-datos.md` (Anexos A y B); `docs/estado/wf0/U1-arquitectura.md` (Anexo); `docs/estado/wf0/correccion-criterios-r1.md` (Anexo) | El integrador copia los scripts a rutas del repositorio y los reportes las citan |

## 7. Para la puerta 1

- **Recibe Axel:**
  - C1 v1.0 (`openapi.yaml`, `C1-datos.md`).
  - U1.
  - Los criterios de las 18 fases.
  - Informativos: S1 1.1 y la colorimetría de U2a.
- **Preguntas abiertas que bloquean comportamiento, no la puerta.** Todas están marcadas PROVISIONAL:
  - `C1-lista-compras` = `S1-lista-compras`.
  - `C1-huecos-u1-s1` = `S1-huecos-negocio`.
  - `C1-filtro-familia`.
  - `C1-parametros-tecnicos`; conviene sumar R3-02.
  - `C1-campos-calculados`.
  - `U1-periodos-semana-quincena` y `U1-captura-fuera-de-casa`.
  - `U2-gasto-error`.
  - `S6-version-go`.
  - Los 15 abiertos de criterios.
  - Los 5 restantes de S1: `S1-frescura-peso`, `S1-meta-kcal`, `S1-redondeo-series`, `S1-reposicion-parte-externa` y `S1-restricciones-optimizador`.
- **Pregunta sugerida, nacida de R2-06:** ¿cuántos decimales de tasa se muestran en pantalla?
- **Para el integrador**: cerrar en el commit de WF-0, o dejar registrados, R3-01 a R3-04 y R2-03 a R2-08. Ninguno requiere otra ronda de revisión.
