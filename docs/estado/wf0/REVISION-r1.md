# Revisión WF-0 «Cimientos» · ronda 1 de 3

Revisor independiente `revisor-wf0-r1` · 10-oct-2026 · Puerta 1 (C1 + U1 + criterios; S1 y U2a como entregables de WF-0).
Evalúo **solo** contra `docs/criterios/` (S1, U1, C1 y los de U2 que aplican a la colorimetría de WF-0). Todo se reejecutó; no se aprobó nada por la palabra de un reporte.

## Veredicto

**NO APROBADO en ronda 1**: 1 hallazgo de severidad alta (H-01), 7 de severidad media y 4 de severidad baja.
- Lo que sí está sólido: el esquema aplica limpio (sqlite3 y goose v3.28.0, arriba/abajo/arriba), el backend compila y pasa `go vet`, el OpenAPI 3.1.1 valida, las convenciones únicas de C1 (errores, paginación, filtros, fechas y dinero) se cumplen al 100 %, la colorimetría se regenera byte a byte con 0 pares bajo umbral, el flujo F-1 demuestra 3 toques y la regresión Kueski completa pasa con 0 diferencias en 109 valores.
- Lo que bloquea: la matriz U1 ↔ C1 no está completa. Cuatro necesidades con datos del servidor no tienen endpoint (N-151, N-152, N-153 y N-165, lista de compras), porque el esquema de S1 no tiene dónde guardarlas, y el paso 3 de la ruta A del flujo clave F-3 depende de una de ellas. C1-02 y U1-03 (alta) no admiten esa excepción.

## 1. Verificaciones ejecutadas por el revisor

| # | Qué pide la tarea | Comando o método | Resultado |
|---|---|---|---|
| 1 | Aplicar el SQL en base temporal y `foreign_key_check` | Sección Up extraída con awk → `sqlite3 -bail datos/tmp/rev.db`; además `goose -dir backend/db/migrations sqlite3 datos/tmp/rev-goose.db up/down/up` (goose v3.28.0) | `UP OK` · 49 tablas STRICT, 100 índices con nombre, 62 disparadores · `PRAGMA foreign_key_check` = 0 filas · `integrity_check` = ok · goose: «OK 00001_esquema.sql», down deja 0 tablas, up otra vez ok. Carpeta `datos/tmp` borrada al final |
| 2 | `go build ./...` y `go vet ./...` | En `backend/` con Go 1.26.5; además `gofmt -l .` | BUILD OK · VET OK · gofmt sin salida. `go.mod` declara `go 1.26` (bloqueo `S6-version-go`) |
| 3 | Ejecutar `colorimetria.py`; `contraste.csv` sin `cumple=false`; hex de C2 = colorimetria.md | El mismo script copiado a un árbol espejo del scratchpad (`--sin-png`, para no reescribir el repositorio) y comparación con `cmp`; luego un script propio (WCAG y OKLab implementados aparte, sin reutilizar código) | Regenera md, csv, json y html **idénticos byte a byte** a los del repositorio · 202 filas, 0 `cumple=false` · recálculo propio de los **202** pares: 0 difieren más de 0.01, 0 bajo umbral · 190 tokens del JSON = 190 filas del Anexo A y 171 valores de las tablas del cuerpo: 0 discrepancias |
| 4 | Matriz C1: toda N-xxx con endpoint y todo endpoint con pantalla o justificación | Script propio: U1 §3 (199 N-xxx) y §3.29 (220 pares) ↔ C1 §10 (220 filas) ↔ `x-necesidades` del OpenAPI (193 operaciones) | Pares U1 = pares C1 (0 diferencias); `x-necesidades` = C1 §10 (0 diferencias); 0 rutas inexistentes. **N-xxx sin endpoint: N-003 y N-301 (locales, válidas) y N-151, N-152, N-153, N-165 (bloqueo `C1-lista-compras`)**. Operación sin N-xxx: `GET /openapi.yaml` (excepción estática justificada en C1 §11.2). El script del propio U1 (reejecutado desde su anexo) da la misma falla en U1-03 |
| 5 | Cada decisión A–I cubierta por S1, U1 o C1 | Lectura dirigida y consultas al esquema y al OpenAPI | Tabla del §3: 40 decisiones aplicables, 39 cubiertas (varias con PROVISIONAL o bloqueo) y 1 parcial porque depende de WF-1 (G7); 3 de proceso no aplican. Una con incoherencia menor (tema, H-04) y una con hallazgo de color (100°, H-03) |
| 6 | PROVISIONAL, supuestos y bloqueos: ¿alguno esconde un valor inventado? | grep en los 7 artefactos y 5 reportes; lectura de los 30 bloqueos; revisión de semillas, `default`, límites y ejemplos del OpenAPI | Toda marca PROVISIONAL remite a un bloqueo existente (U1 10, C1 44 + 63 en YAML, diccionario 6, reglas 22, migración 4). Semillas de la migración: solo valores decididos. Valores técnicos de C1 (24 h, 90 días, 12 caracteres, 200, 10 MB, 1 GB, SameSite=Strict) marcados y con bloqueo. Dos puntos menores en §5 (tope de Karen y descripción por defecto del gasto) |
| 7 | Criterios medibles y trazables al plan | Script propio sobre los 18 archivos | 248 criterios (189 alta, 59 media) = INDICE; 248 con «Verifica» y «Evidencia»; 0 enunciados vagos sin cifra; **35 sin línea explícita en su §5 Trazabilidad**; 4 inconsistencias con los entregables (H-06, H-07, H-10) |
| 8 | Registrar gasto en ≤ 3 toques demostrado en U1 | Lectura de U1 §4 (regla de conteo) y F-1 | Demostrado: «+» (1), chip de categoría (2), «Guardar» (3); regla de conteo escrita; fecha, persona y forma de pago por defecto, reparto desde la categoría. Depende de que la descripción no sea obligatoria (ver H-02) y de la verificación en Android de U5 |
| 9 | Regresión de S1 (Kueski 389500 y 239528) | Script propio: inserta las 4 despensas, sus partes, los 16 pagos con partes y el papel de contado; reconstruye eventos y series desde la base | Suma de pagos de la despensa 1 = **389500**; parte de Axel = **239528** (por pagos y en la compra). 17 eventos, series A y B de 14 valores, mínimos 8,094.50 / 5,000.00 y efectivo y devengado de P1–P6 por mes calendario: **109 valores, 0 diferencias** |

Pruebas negativas sobre el esquema (revisor): rechaza cuenta para Karen, segundo titular, hash que no es argon2id, ingreso de Ximena, plan de comidas de Karen, 12.5 en `_cents`, fechas `10-oct-2026` y `2027-02-30`, tarjeta sin corte, mercado sin día de visita, quincena 2 de febrero pagando el 30 y borrar a Karen con historia. Acepta (documentado como validación de la API): partes 60 % + 50 % en una regla (invariante 8 del diccionario), una tercera cuenta para otra co-titular y cambiar la participación de Ximena a externa con cuenta viva (invariante 12; la API no tiene operaciones para ninguna de las dos cosas).

## 2. Tabla de criterios

Estados: **Cumple** · **Cumple con observación** · **No cumple** · **No aplica en WF-0**.

### S1 · Dominio y datos
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| S1-01 Entidades completas | alta | Cumple | Mapa del diccionario: 31 de 31 entidades → `nombre_codigo`; las 49 tablas existen en la base migrada |
| S1-02 Migraciones en SQLite limpio | alta | Cumple | sqlite3 y goose v3.28.0: up ok, `integrity_check` ok, `foreign_key_check` vacío; down/up ok |
| S1-03 Cada regla con ejemplo del reporte | alta | Cumple con observación | Las reglas mínimas (redondeo, último pago, devengado, quincena y febrero, reparto, A/B, frescura) citan cifras de `model.json`/`data.json` que existen. 10 reglas sin datos en el oráculo (salud, mercados, Sam's, ahorro…) llevan ejemplo marcado «ilustrativo» o «hipotético»; las cifras fuera del oráculo en RN-07 y RN-08 son contrafactuales explicados. Ver H-10 |
| S1-04 Ninguna tabla huérfana | alta | Cumple | 49 de 49 con «Usado por». Comprobadas a mano: `store_product_group`, `credit_statement`, `login_attempt`, `notification_delivery`, `nutrition_fact` |
| S1-05 Dinero exacto | alta | Cumple | 0 columnas REAL/FLOAT/DOUBLE; 35 columnas `_cents`, todas INTEGER; 12.5 rechazado |
| S1-06 Fechas y zona | media | Cumple | grep de formato de pantalla en migraciones = 0; CHECK `date(x) IS x` rechaza `10-oct-2026` |
| S1-07 Personas y acceso | alta | Cumple con observación | Cuenta de Karen e ingreso de Ximena rechazados; sin columna de rol. La base admite una tercera cuenta (H-11) |
| S1-08 Reparto y participantes | alta | Cumple | Suma ≤ 100 % documentada como invariante 8 (API, prueba para S6); despensa 1 cabe sin pérdida (regresión) |
| S1-09 Periodos e ingresos | alta | Cumple con observación | Meses y quincenas con CHECK; febrero 2027 Q2 paga el 28. La prueba literal (insertar P1–P6 del 10 al 9) es incompatible con la decisión B y el CHECK la rechaza; la equivalencia P1–P6 ≡ oct-2026…mar-2027 se verificó (H-07) |
| S1-10 Créditos | alta | Cumple | 4 tipos con CHECK; tasas como fracción exacta 53450/309450 (17.27 %) y 80050/309450 (25.8685 %); promoción con comercio, plazo, mínimo y vigencia; A/B con `DEFAULT 'B'` y `availability_confirmed` |
| S1-11 Comida, salud, inventario, precios | media | Cumple | 12 de 12 puntos con tabla y columna |
| S1-12 Ida y vuelta de la semilla del oráculo | alta | Cumple con observación | Verificado por el revisor (109 valores, 0 diferencias). La prueba del agente solo cubrió la despensa 1 y no quedó en el repositorio (H-04) |
| S1-13 Ajustes tiene dónde guardarse | media | No cumple | Falta la tabla sección → campos → tabla.columna; «Apariencia» sin destino coherente (H-04) |
| S1-14 Sesiones y credenciales | alta | Cumple | Solo hash argon2id (CHECK) y `token_hash` SHA-256; grep sin valores secretos; 0 tokens de Telegram |
| S1-15 Integridad y consultas de C1 | media | Cumple | FK en `expense`, `credit_purchase`, `scheduled_payment_share`, `price`, `meal_plan_entry`; índices para persona, categoría, tienda, periodo (fechas) y crédito; «familia» sale de las partes por persona (índices en `*_share.person_id`) |
| S1-16 Sin valores inventados | alta | Cumple con observación | Semillas decididas; PROVISIONAL con bloqueo. Interpretación no declarada del tope de Karen para la semilla de S8 (H-11) |
| S1-17 Idioma y nombres | media | Cumple | Identificadores en inglés, glosario `nombre_codigo`, textos en es-MX |

Además: el reporte de S1 no tiene ninguna línea de evidencia que empiece con su ID (convención de `INDICE.md` §1) → H-04.

### U1 · Arquitectura de información
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| U1-01 Barra lateral de 9 entradas | alta | Cumple | 9 entradas en el orden exacto |
| U1-02 Inventario con identificador | alta | Cumple | 27 pantallas P-00…P-26 únicas + 14 secciones P-26.01…P-26.14; U1 fija P-xx para pantallas y N-xxx para necesidades (ver H-06) |
| U1-03 Matriz pantalla ↔ endpoint completa | alta | **No cumple** | N-151, N-152, N-153 y N-165 sin endpoint (H-01) |
| U1-04 Gasto en ≤ 3 toques | alta | Cumple | 3 toques con regla escrita (U1 §4) y valores por defecto |
| U1-05 Seis flujos clave | alta | Cumple con observación | 8 flujos (6 de la tarea + F-7 comparador y F-8 ticket), todos los pasos con pantalla existente; el criterio dice 6 (H-10). F-3 ruta A paso 3 sin endpoint (H-01) |
| U1-06 Vistas por entidad | alta | Cumple | 31 entidades: 23 con vistas, filtro y orden; 8 con motivo |
| U1-07 Menú «/», atajos, edición | media | Cumple | 27 comandos con destino; 15 atajos sin duplicados; §2.6 por pantalla |
| U1-08 Ajustes con 14 secciones | alta | Cumple | 14 de 14 en orden y con nombre exacto, todas en la matriz |
| U1-09 Cobertura de las decisiones | alta | Cumple | 26 capacidades con P-xx y N-xxx |
| U1-10 Dos titulares, Karen sin acceso | alta | Cumple | Sin pantalla de acceso de Karen; solo reparto, partes, participación, estado de cuenta y Ajustes |
| U1-11 Móvil primero, una acción primaria | media | Cumple con observación | §6.4 con 27 filas; objetivo táctil mínimo de 24 px frente a 44 px del plan (H-09) |
| U1-12 Consistencia con S1 y C1 | media | Cumple con observación | Toda tabla de S1 tiene pantalla o motor; en sentido contrario hay necesidades de U1 sin respaldo en S1 (H-01, H-02) |
| U1-13 Idioma | media | Cumple | 77 nombres visibles, 0 en inglés |
| U1-14 Supuestos y bloqueos | media | Cumple | Sección Supuestos; 2 bloqueos U1; PROVISIONAL marcado |

### C1 · Contrato de datos y API
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| C1-01 OpenAPI 3.1 válido | alta | Cumple | `uvx --from openapi-spec-validator==0.9.0 openapi-spec-validator` → OK, salida 0; 3.1.1, 154 rutas, 193 operaciones |
| C1-02 Toda N-xxx con endpoint y viceversa | alta | **No cumple** | 4 N-xxx sin endpoint por bloqueo; no son excepciones técnicas (H-01) |
| C1-03 Endpoints respaldados por tablas | alta | Cumple | 49 de 49 tablas citadas; 906 propiedades: 824 columnas, 71 calculadas, 11 anidadas, 0 sin respaldo (script propio) |
| C1-04 Una convención de errores | alta | Cumple | 874 de 874 respuestas 4xx/5xx/default → `Error` (script propio) |
| C1-05 Una convención de paginación | alta | Cumple | 39 de 39 colecciones con `cursor`/`limit`/`sort` y `Page` |
| C1-06 Una convención de filtros | alta | Cumple | 6 de 6 en `components/parameters`; los homónimos en línea son parámetros de ruta; `family` PROVISIONAL con bloqueo |
| C1-07 Una convención de fechas | alta | Cumple con observación | Solo `date` y `date-time`; quincenas y febrero declarados. El patrón de `DateTime` acepta «Z» (H-12) |
| C1-08 Una convención de dinero | alta | Cumple | 253 propiedades `_cents` → `MoneyCents` (integer); 0 `number` |
| C1-09 Métricas con fórmula y periodo | alta | Cumple | 17 métricas; 6 tarjetas; 10,000.00 / 6,905.50, 30-oct-2026 973.75 y P1 1,947.50 / 973.75 presentes y del oráculo |
| C1-10 Recursos y vistas de las 31 entidades | alta | Cumple | §7 con 31 filas |
| C1-11 Cobertura funcional | alta | Cumple | 35 capacidades y 14 secciones (Apariencia local) |
| C1-12 Seguridad | alta | Cumple | Públicas solo login y salud; 0 campos sensibles en respuestas; cookie HttpOnly/Secure/SameSite; 429 `too_many_attempts` |
| C1-13 Versionado y estructura | media | Cumple | 0 `operationId` repetidos; 9 etiquetas de la barra lateral; C1 v1.0 |
| C1-14 Idioma y ejemplos | media | Cumple | 13 ejemplos válidos con jsonschema; 2,607 descripciones sin inglés frecuente |
| C1-15 PROVISIONAL marcado | alta | Cumple | 63 marcas en YAML y 44 líneas en C1, todas con bloqueo existente |

### U2 · aplicados a la colorimetría de WF-0 (U2a)
| Criterio | Sev. | Estado | Evidencia del revisor |
|---|---|---|---|
| U2-01 DESIGN.md con 10 secciones | alta | No aplica en WF-0 | DESIGN.md es de WF-1 |
| U2-02 Matiz C aplicado | alta | Cumple con observación | Matices medidos a ±0.7° de lo declarado; marca y neutro 22 de 22 del plan; ingreso 135° es la opción que prevé decisiones I. La serie 6 de gráficas usa el 100° reservado (H-03) |
| U2-03 Método y saturación | alta | Cumple | Desviación de L por paso ≤ 0.0016; croma máx. marca 0.110, neutro 0.012, roles 0.141 (redondeo a hex) |
| U2-04 0 pares bajo umbral | alta | Cumple | 202 pares, 0 bajo umbral; recálculo propio de los 202 coincide ±0.01 |
| U2-05 Pares estructurales | alta | Cumple | 11.53, 5.40, 3.47, 3.98, 14.07, 7.97, 10.38, 1.20, 1.90; borde de campo 3.71 (el plan decía ≥ 4.2, documentado) |
| U2-06 Roles como texto | alta | Cumple | Paso 700 en claro (5.07–5.78) y 400 en oscuro (7.60–8.39) |
| U2-07 Botón primario | alta | Cumple | Blanco/marca 700 = 5.77; marca 700/fondo = 5.53; oscuro 10.47 y 10.38; blanco sobre marca 600 no aparece como texto |
| U2-08 Colisiones resueltas | alta | Cumple | Pares < 20° y ΔL < 0.15: 3 en claro y 4 en oscuro, los mismos que lista §7.7 con su refuerzo; gasto y error con bloqueo `U2-gasto-error` |
| U2-09 Paleta de 8 series | alta | Cumple | ΔL consecutivo mín. 0.179 (claro) y 0.223 (oscuro); contraste mín. 3.15 y 3.11; simulación documentada; 28 pares informativos (bloqueo `criterios-deltal`) |
| U2-10 a U2-13 | alta/media | No aplica en WF-0 | Tipografía, iconos, espaciado, componentes y patrones: WF-1 |
| U2-14 Idioma y supuestos | media | Cumple | es-MX, tokens en inglés, Supuestos y bloqueo |

### Criterios (archivo `docs/criterios/`)
| Aspecto | Estado | Evidencia |
|---|---|---|
| Completos y con estructura | Cumple | 19 archivos; 248 criterios = INDICE (189 / 59); 5 secciones en cada archivo |
| Medibles | Cumple | 248 de 248 con «Verifica» (comando o cálculo) y «Evidencia»; 0 enunciados vagos sin cifra |
| Trazables al plan | Cumple con observación | 213 de 248 con línea en su §5; 35 sin ella (H-08) |
| Coherentes con los entregables y decisiones | Cumple con observación | N-xx (H-06), S1-09 (H-07), U2-02, U1-05, C1-01 y C1-02 (H-10) |

## 3. Cobertura de las decisiones A a I

| Id | Decisión (resumen) | S1 (esquema o regla) | U1 (pantalla) | C1 (operación) | Estado |
|---|---|---|---|---|---|
| A1 | Co-titulares con los mismos permisos; usuario y contraseña; sesión larga por dispositivo; cierre desde Ajustes; argon2id; 5 en 15 min | `user_account` sin rol, CHECK argon2id; `device_session`; `login_attempt` + `app_setting auth.*` | P-00, P-26.01 (N-010, N-011, N-271, N-272) | `login` (429), `logout`, `listSessions`, `revokeSession(s)` | Cubierta (duración PROVISIONAL, `criterios-login-bloqueo`) |
| A2 | Karen externa sin cuenta; solo despensa y compras Kueski; paga directo; informativa; sin finanzas de su familia; estado de cuenta PDF | `participation = external` + disparador; `pays_share_directly`; partes; `generated_report` | P-25, N-026, N-210 | `getPersonParticipation`, `getExternalStatement`, `generateExternalStatement`; sin alta de cuentas | Cubierta (reposición por parte PROVISIONAL) |
| A3 | Comidas y salud solo Axel y Ximena | Disparadores en `meal_plan_entry`, `health_profile`, `weight_entry` | P-10, P-21 | 422 `external_person` | Cubierta |
| A4 | Telegram a Axel y Ximena; destinatarios por aviso | `notification_recipient.chat_id`, `notification_subscription` | P-26.09 (N-298, N-299) | `setNotificationRecipients`, `createNotificationRecipient` | Cubierta (registro del chat PROVISIONAL) |
| B1 | MXN, centavos, HALF_UP, último ajusta; zona; fechas; es-MX | `_cents` INTEGER STRICT; CHECK ISO; RN-01, RN-07 | §0.4 | `MoneyCents`, `Date`, `DateTime` | Cubierta |
| B2 | Mes calendario; quincena 15/30; despensa devengada a la mitad | `period` (CHECK); `accrual_months` | §0.2, P-06 (N-083) | `period`, `getAccruedGroceries` | Cubierta (rangos PROVISIONAL) |
| B3 | Ingresos solo de Axel, 15/30, extras | `income`, `planned_item` + disparadores | P-08 | `incomes`, `planned-items` | Cubierta |
| B4 | Reparto % por gasto o categoría, 50/50; qué se comparte con Karen por gasto | `cost_share_rule(_part)`, reglas en categoría, gasto y renglón; `credit_purchase_share` | N-051, N-281, N-309 | `cost-share-rules`, filtro `family` | Cubierta (`family` PROVISIONAL) |
| B5 | Fondo y metas con fecha, aportación, avance | `savings_goal`, `savings_movement` | P-09, N-025 | `savings-goals`, M-07, M-16 | Cubierta |
| B6 | Comida fuera como gasto con persona; plan marca «fuera» | `eating_out` (semilla S8), `expense.person_id`, `is_eating_out` | N-124 | `createMealPlanEntries` | Cubierta |
| C1 | 4 tipos de crédito | `credit.kind` | P-17, P-18, P-26.05 | `credits` | Cubierta |
| C2 | Tasa por compra con IVA; tabla por plazo 17.27 / 25.8685; promociones | `interest_cents` (quoted), `credit_term_rate` (fracción), `promotion` | N-218, N-286, N-287 | `term-rates`, `promotions` | Cubierta (ver H-12) |
| C3 | A/B por defecto B, aviso «pendiente confirmar» | `availability_reading DEFAULT 'B'`, `availability_confirmed` | N-022, N-284, N-285 | `previewAvailability` | Cubierta |
| C4 | Comparador: escenarios, 3 restricciones, 4 objetivos, recomendación, sin solver | `planned_item`, `app_setting optimizer.*` | P-20, F-7 | `compareScenarios` (restricciones, ventana, pesos; 4 óptimos, recomendación, descartados) | Cubierta (valores PROVISIONAL) |
| C5 | Regresión Kueski | Regresión del revisor: 0 diferencias | — | Ejemplos M-03, M-04, M-11, M-12 | Cubierta en WF-0 (motores en WF-1) |
| D1 | Recetario de Claude + propias; restricciones vacías | `recipe.origin`, `dietary_restrictions` | P-11, P-12 | `recipes` | Cubierta |
| D2 | 4 tiempos; porciones base 2 | `meal_slot` (4), `base_servings` 2 | P-10 | `meal-slots`, plan | Cubierta |
| D3 | Perfil, IMC, meta kcal, kcal por receta y día; macros después | `health_profile`, `weight_entry`, macros en `nutrition_fact` | P-21, P-22 | `health-profiles`, `getTodayKcal` | Cubierta (meta kcal PROVISIONAL) |
| D4 | Fuentes nutricionales autorizadas | `nutrition_fact.source`, licencia | N-294, N-295 | `reloadNutritionBase` | Cubierta |
| D5 | Inventario híbrido; sobrantes; mercados con día de visita | `inventory_item`, `store` market + `visit_weekday`, `store_product_group` | P-16, P-13, P-26.06 | `inventory-items`, `stores` | Cubierta (movimientos PROVISIONAL) |
| E1 | Manual, OCR con confirmación, lectura asistida apagada, PROFECO | `price.source`, `ocr_status`, `assisted_reading_enabled = 0` | P-05, F-8, N-291 a N-293 | `confirmTicketScan`, `importProfeco`, 409 | Cubierta |
| E2 | Viejo a los 30 días, pesa menos | `price.stale_after_days = 30` | N-170, N-171, N-161 | frescura en precios | Cubierta (peso PROVISIONAL) |
| E3 | Sam's: punto de equilibrio | `requires_membership`, `membership_annual_fee_cents` | N-164 | `getBasketComparison` | Cubierta (cálculo PROVISIONAL) |
| E4 | Unidad base, sustitutos, mejor tienda | `base_unit`, `size_milli`, RN-20 | N-161 a N-163 | comparativa | Cubierta |
| F1 | Diario, semanal, mensual PDF día 1, vista web; estado de Karen | `notification` (7 avisos), `job`, `generated_report` | P-23, P-24, P-25, P-26.08 | reportes | Cubierta (hora del diario sin valor, bloqueo) |
| F2 | 4 alertas | Semillas `payment_due` [1, 0], `budget_threshold` 80/100, `stale_prices`, `meal_plan_missing` | P-26.09 | `notifications`, `getBudgetAlerts` | Cubierta |
| F3 | MultiX texto plano, token en tiempo de ejecución | Sin token en la base | N-299 | Sin campos sensibles | Cubierta (prefijo: orquestador) |
| G1 | Fiel a Notion | — | §1, §2.5, §2.6, §5 | `queryView` (4 modos) | Cubierta |
| G2 | Android y escritorio; PWA mkcert; ≤ 3 toques | — | F-1, §6, N-274 | `getCaCertificate` | Cubierta (PWA en WF-3) |
| G3 | Tema del sistema con interruptor | `appearance.*` sin nombrar | P-26.10 local | sin API | Cubierta con incoherencia (H-04) |
| G4 | Paleta nueva; matiz I | — | — | — (U2a: 190 tokens) | Cubierta |
| G5 | Resumen ≤ 6 tarjetas, Detalle, Tendencias, Calendario | — | P-01 (6 tarjetas) | M-01 a M-07 | Cubierta |
| G6 | Ajustes, 14 secciones | «Usado por» | P-26.01 a P-26.14 | C1 §11.5 | Cubierta (S1-13, H-04) |
| G7 | Autoalojado; cifras tabulares; WCAG 2.2 AA | — | §0.4 | — (U2a: contraste) | Parcial: tipografía e iconos en WF-1; táctil (H-09) |
| H1 | Solo LAN, finanzas.local:8443 | — | — | Base del contrato; esqueleto en 127.0.0.1 (`criterios-bind`) | Cubierta |
| H2 | Go + SQLite + systemd; respaldos locales | `backup` (`integrity_check`, `restore_tested_at`) | P-26.11 | `backups` | Cubierta (systemd en WF-4) |
| H3 | IA opcional apagada | `ai.*` = false | N-306, N-307 | `updateAiProvider`, 409 `ai_disabled` | Cubierta |
| H5 | Semilla e importador | `planned_item` | N-304, N-305 | `importPlannedItems`, `getSeedStatus` | Cubierta |
| I1 | Pila A | goose (validado), Go `net/http` | — | `routes.go` | Cubierta (Go 1.27: bloqueo `S6-version-go`) |
| I2 | Matiz C y reglas | — | — | — (U2a) | Cubierta con hallazgo (H-03) |
| H4, H6, J | Proceso (paquetes, repositorio, puertas) | — | — | — | No aplica al producto |

## 4. Bloqueos, supuestos y PROVISIONAL

- 30 bloqueos en `docs/estado/bloqueos/` (16 de criterios, 5 de S1, 2 de U1, 1 de U2, 5 de C1 y 1 de S6), todos con el formato del protocolo.
- Los bloqueos de criterios y de S1 son preguntas de negocio legítimas.
- `C1-lista-compras` y 5 de los 7 puntos de `C1-huecos-u1-s1` (descripción del gasto, tiempos de receta, movimientos de inventario, presupuesto vigente, monto por quincena) no esconden valores inventados. Son incoherencias internas de WF-0 que S1 puede corregir sin Axel (H-01, H-02). Solo la elección «marcas compartidas, locales o ninguna» de la lista y el registro del chat son preguntas de producto.
- Ningún valor inventado entró sin marca. Hay dos puntos menores:
  - C1 llena la descripción del gasto con el nombre de la categoría (PROVISIONAL, para no romper F-1).
  - S1 propone sembrar el tope por quincena de Karen = 374.93 (`karen_fija.pago`) en S8, interpretado como «tope» del objetivo 2 del optimizador y sin declararlo (H-11).

## 5. Hallazgos

| Id | Sev. | Artefacto | Criterio | Descripción | Ruta | Corrección sugerida |
|---|---|---|---|---|---|---|
| H-01 | **alta** | S1 | C1-02, U1-03 (U1-12) | La lista de compras de U1 (marcar comprado N-151, agregar artículo N-152, cambiar tienda o cantidad N-153, «usar esta tienda» N-165) no tiene tabla en S1. Por eso no tiene endpoint en C1. El paso 3 de la ruta A del flujo clave F-3 («Usar Bodega Aurrera») queda sin operación. U1 lo pidió a S1 (solicitud 3 de su reporte); el esquema no lo incluyó ni lo bloqueó, y el bloqueo lo abrió C1 después. Las necesidades bloqueadas no son excepciones técnicas de C1-02/U1-03 | `backend/db/migrations/00001_esquema.sql`; `docs/contratos/C1-datos.md` §11.3; `docs/estado/bloqueos/C1-lista-compras.md` | En la ronda 2, cerrar la matriz con una opción del bloqueo marcada PROVISIONAL hasta la puerta 1. Opción A: S1 agrega `shopping_list` y `shopping_list_item` (semana, producto o texto libre, cantidad, tienda, comprado, origen) y regenera el diccionario; C1 regenera OpenAPI, `routes.go` y la matriz con las 4 operaciones. Opción B: U1 declara N-151 a N-153 y N-165 como Local. Reejecutar la matriz: 0 N-xxx sin endpoint ni declaración Local |
| H-02 | media | S1 | U1-12, C1-15 | Siete huecos S1 ↔ U1 (`C1-huecos-u1-s1`) más la clave de idempotencia obligan a C1 a cubrir 13 necesidades con lecturas PROVISIONALES que recortan U1 o fijan un comportamiento. Ejemplo: `expense.description` NOT NULL contra N-054 opcional del flujo F-1; C1 guarda el nombre de la categoría. Otros: un solo tiempo por receta, movimientos de inventario sin motivo, presupuesto sin vigencia, un solo monto quincenal | `backend/db/migrations/00001_esquema.sql`; `docs/contratos/C1-datos.md` §14; `docs/estado/bloqueos/C1-huecos-u1-s1.md` | S1 corrige lo que no requiere a Axel: descripción NULL, tabla receta-tiempo, `inventory_movement` con motivo, vigencia del presupuesto, monto por quincena, `idempotency_key`. Deja como bloqueo solo lo de negocio (registro del chat, si el presupuesto rige a meses siguientes, si se guardan correcciones a medias del OCR). Regenerar diccionario y C1 |
| H-03 | media | U2a | U2-02 | La serie 6 de gráficas («mostaza») usa el matiz 100°, que decisiones I y U2-02 reservan **solo** para «acción crítica» de confirmación. El documento no declara la excepción | `docs/design/colorimetria.md` §8.1; `docs/contratos/C2-colores.json` `color.chart.6.*`; `docs/design/tools/colorimetria.py` | Mover la serie 6 a un matiz no reservado y regenerar con el script (ΔL, CVD y contraste se reverifican solos). Otra opción: justificar por escrito la excepción para gráficas y llevarla a la puerta 2 |
| H-04 | media | S1 | S1-12, S1-13; INDICE §1 | Evidencia de S1 incompleta. (a) El reporte no tiene líneas `S1-xx`. (b) La ida y vuelta pedida (4 despensas, 16 pagos, 17 eventos, series A/B) no la corrió el agente (solo la despensa 1) y su verificador quedó en el scratchpad efímero. (c) Falta la tabla de las 14 secciones → campos → `tabla.columna`. (d) «Apariencia» se mapea a claves `appearance.*` «por cuenta» sin nombre, mientras U1 y C1 la dejan local por dispositivo | `docs/estado/wf0/S1-modelo-datos.md`; `docs/contratos/diccionario-datos.md` («Claves de app_setting previstas sin valor») | Agregar las 17 líneas de evidencia. Guardar en el repositorio una prueba de ida y vuelta reejecutable (el revisor obtuvo 109 valores y 0 diferencias). Agregar la tabla de S1-13 con 14 de 14 secciones. Alinear Apariencia con U1 (local por dispositivo) o nombrar sus claves |
| H-05 | media | C1 | S6-03, C1-02 (reproducibilidad) | `openapi.yaml` (557 KB), `routes.go` («DO NOT EDIT») y la matriz de C1 §10 salen de generadores (`gen1…gen6`, `verify_c1.py`) que solo existen en el scratchpad efímero. Lo mismo pasa con el verificador de S1, `u1_verify.py` (solo en el anexo de su reporte) y `gen_canon.py` de criterios. La ronda 2 (H-01, H-02) los necesita | scratchpad `…/c1/`; `backend/internal/api/routes.go` | El integrador copia los generadores y verificadores al repositorio (p. ej. `docs/contratos/tools/c1/`, `docs/design/verificar_u1.py`, `ops/tools/`) antes de la ronda 2. Los reportes y criterios citan esas rutas |
| H-06 | media | criterios | U1-02, U1-04, U1-05, U1-08, C1-02, C1-10 (y C3/U5) | Los criterios usan «N-xx» como identificador de pantalla. U1 y C1 usan P-xx para pantallas y N-xxx para necesidades de datos (la lectura que también da esta tarea). El bloqueo `criterios-nxx` sigue abierto y U1 y S6a pidieron alinearlo | `docs/criterios/U1.md`, `C1.md`, `C3.md`, `U5.md`; `docs/estado/bloqueos/criterios-nxx.md` | Reescribir: pantalla = P-xx (P-26.xx); necesidad = N-xxx. C1-02 cruza N-xxx ↔ `operationId` y P-xx ↔ al menos una operación o «Local». Cerrar `criterios-nxx` |
| H-07 | media | criterios | S1-09 (y S2-09, R-45…R-50) | S1-09 pide insertar los periodos P1–P6 del oráculo (del día 10 al 9). Eso contradice la decisión B (mes calendario) y el CHECK de `period` lo rechaza (probado). La equivalencia P1–P6 = oct-2026…mar-2027 se cumple (12 de 12 cifras), pero solo porque en la semilla ningún evento cae entre los días 1 y 9 | `docs/criterios/S1.md` S1-09; `docs/criterios/S2.md` | Cambiar la verificación: meses calendario + prueba de equivalencia con el oráculo. Agregar una nota para S2: la equivalencia depende de la semilla; con eventos entre los días 1 y 9 hay que comparar contra las ventanas del oráculo en las fixtures |
| H-08 | media | criterios | Trazabilidad (todas las fases) | 35 de 248 criterios no aparecen en la sección «Trazabilidad» de su archivo: C1-04, S1-15, U1-14, C2-09/10/11, C3-03/10/11, P3-13, S2-08/09/14/16, S3-04/11/20, S4-06/15, S5-13, S6-08/10/15/16, S7-14, S8-08, U3-08/10/12, U4-07/09/10/11, U5-06/12 | `docs/criterios/*.md` §5 | Agregar la línea que los traza al plan, a decisiones o al protocolo, o marcarlos como convención propia del directorio |
| H-09 | baja | U1 | U1-11 (U2-11) | U1 §6.1 fija objetivos táctiles «de al menos 24 × 24 px». El plan §4.1 punto 9 y U2-11 fijan ≥ 44 px | `docs/design/U1-arquitectura.md` §6.1 | Citar ≥ 44 px (plan §4.1) como mínimo táctil y conservar 24 px solo como piso WCAG |
| H-10 | baja | criterios | U2-02, U1-05, C1-01, C1-02, S1-03 | U2-02 da ingreso 150° sin la opción «o U2 propone 135°» de decisiones I. U1-05 espera exactamente 6 flujos, pero exige 4 propios que no coinciden con los 6 de la tarea (se entregaron 6 + 2). El comando de C1-01 usa un módulo que no está en el Python del sistema. C1-02 no prevé necesidades bloqueadas. S1-03 no dice cómo tratar reglas sin datos en el oráculo | `docs/criterios/U2.md`, `U1.md`, `C1.md`, `S1.md` | Admitir 135° como propuesta para la puerta 2; «≥ 6 flujos con los 4 obligatorios»; equivalente `uvx --from openapi-spec-validator==0.9.0 …`; aclarar que una necesidad bloqueada no cuenta como cubierta; permitir ejemplos «ilustrativos» marcados solo para reglas sin datos en el oráculo |
| H-11 | baja | S1 | S1-07, S1-16 | La base acepta una tercera cuenta (otra co-titular) y deja una cuenta viva si la participación cambia a externa (probado). La API lo impide: no hay operaciones de alta ni de cambio de participación. La semilla planeada para S8 fija el tope de Karen = 374.93 (`karen_fija.pago`) como tope del objetivo 2 sin declararlo como supuesto | `backend/db/migrations/00001_esquema.sql`; `docs/estado/wf0/S1-modelo-datos.md` (Pendientes 3) | Disparadores en `person` (participación) y límite de dos cuentas como defensa en profundidad. Declarar el tope de Karen como supuesto PROVISIONAL o agregarlo a `criterios-optimizador` |
| H-12 | baja | C1 | C1-07; decisión C | El patrón de `DateTime` acepta «Z», pero C1 §4 dice que los instantes llevan el desplazamiento de America/Mexico_City. `rate_pct` con 4 decimales mostraría 17.2725 % a 2 quincenas, contra el 17.27 % de la decisión C | `docs/contratos/openapi.yaml` (`DateTime`, `rate_pct`); `docs/contratos/C1-datos.md` §4 | Restringir el patrón de salida o documentar que «Z» solo se acepta en entradas. Definir la regla de visualización de tasas (como se capturó, o 2 decimales cuando la fracción viene de una cotización) |

## 6. Lo que la ronda 2 debe traer para aprobar
1. H-01 resuelto: matriz con 0 N-xxx sin endpoint ni declaración Local. El revisor reejecuta `matriz_rev.py` (o el verificador que el integrador deje en el repositorio).
2. H-02 a H-08 corregidos o justificados por escrito (media).
3. Los generadores y verificadores en rutas del repositorio (H-05), para reejecutar sin copiar anexos.
