# Convenciones del backend · WF-3a «Construcción»

Versión 1.1 · 12-oct-2026 · agente **backend-base** (Opus); 1.1 de `correccion-backend-core-r1`: ruta real del paquete de gastos (§2, §14) y `areas.go` completo (hallazgos M2 y H4 de REVISION-backend-core-r1). Obligatorias para todos los agentes de backend de WF-3 (api-*, auth-sistema, s7-*). Lo que no esté aquí se pregunta (regla suprema); no se inventa.

## 1. Mapa del código (qué ya existe y no se reescribe)

| Ruta | Qué es | Quién lo cambia |
|---|---|---|
| `backend/internal/api/routes.go`, `types_gen.go` | Generados por `c1gen/build.py` (gen5 y **gen8_types.py**): 216 rutas, 343 tipos de `components/schemas` y `Operations` (sesión, Idempotency-Key, cuerpo, 2xx de cada operación) | Nadie a mano. `python3 -I backend/internal/api/c1gen/build.py --check` = 0 |
| `backend/internal/api/server.go`, `http.go`, `nullable.go` | Registry, middleware (sesión, agente externo, idempotencia, errores, 404/405, cabeceras, registro de accesos), helpers HTTP | backend-base |
| `backend/internal/dominio/` | SQLite (modernc), migraciones embebidas (goose biblioteca), repositorios comunes, dinero, fechas, periodos, quincenas, respaldos | backend-base (cobertura ≥ 85 %) |
| `backend/internal/jobs/` | Planificador `time.Ticker` + tabla `job`, cron de 5 campos en America/Mexico_City | backend-base |
| `backend/internal/ia/` | Proveedores de IA apagados (`Off`) y `Registry` | backend-base |
| `backend/internal/testutil/` | Base temporal migrada + fixtures, reloj simulable, servidor TLS con sesión, validación contra `openapi.yaml` | backend-base |
| `backend/cmd/finanzas/` | `serve`, `migrate`, `semilla`, `respaldo`; `areas.go` junta los `Register` | backend-base (cada área pide su línea en `areas.go`) |
| `backend/db/embed.go` | `go:embed migrations/*.sql` (paquete `finanzas-hogar/db`) | backend-base |
| Motores `motor_fin`, `motor_cred`, `recetas`, `precios` | Aprobados con oráculos | Nadie: se llaman, no se reescriben ni cambian cifras |

## 2. Un paquete por agente y cómo registra sus operaciones

| Agente | Paquete | Agente | Paquete |
|---|---|---|---|
| api-gastos-inicio | `internal/api/gastos-inicio` (`package gastos`; import `gastos "finanzas-hogar/internal/api/gastos-inicio"`) | api-creditos | `internal/api/creditos` |
| api-presupuesto | `internal/api/presupuesto` | api-ajustes | `internal/api/ajustes` |
| api-comidas-salud | `internal/api/comidas` | auth-sistema | `internal/auth` |
| api-compras-a | `internal/api/comprasa` | s7-avisos | `internal/telegram` |
| api-compras-b | `internal/api/comprasb` | s7-pdf | `internal/reportes` |

(`internal/auth`, `internal/telegram` e `internal/reportes` son los nombres que mide S6-02.) El directorio de api-gastos-inicio es `gastos-inicio` (así lo pidió su ALCANCE) y el paquete se llama `gastos`: se importa con el alias `gastos`. Cada paquete exporta **una** función:

```go
func Register(reg *api.Registry) { reg.Handle("listUsers", listUsers(reg)); /* … */ }
```

- `reg.Handle(operationID, h)` entra en pánico si la operación no está en C1 o ya tiene dueño: dos agentes no pueden reclamar la misma.
- `backend/cmd/finanzas/areas.go` las junta (ya trae las 10 áreas; `TestAreas_SinOperacionesEn501` exige las 216 operaciones): una línea `gastos.Register,` y su import por agente (fuera del alcance de cada área: se pide en `solicitudes_fuera_de_alcance` con esa línea exacta, o la agrega el integrador). Lo no registrado responde **501** `not_implemented`.
- `api.RegisterBase` pone `getHealth` y `getOpenApiDocument` solo si nadie los registró (son de auth-sistema; al registrarlos, los suyos ganan).
- Dependencias: `reg.Deps()` → `DB`, `Clock` (usar siempre `Clock.Now()`, nunca `time.Now()`), `Log` (slog), `IA`, `Jobs`, `Service`, `Config` (`DataDir`, `BackupDir`, `TmpDir`, `OpenAPIPath`, `FrontendDir`).

## 3. Manejadores, tipos y errores

- Firma: `api.HandlerFunc func(w http.ResponseWriter, r *http.Request) error`. Se responde con `api.WriteJSON(w, status, v)` o `api.NoContent(w)` y se devuelve `nil`; o se devuelve un error y la API escribe el **esquema único de error** (C1 §3) con `request_id`.
- Errores de negocio (`internal/dominio`): `dominio.Invalid(fields…)` 422 `validation_failed`; `dominio.NotFound(msg)`/`ErrNotFound`/`sql.ErrNoRows` 404; `dominio.Conflict(code, msg)` 409 (`in_use`, `already_exists`, `has_credit_purchase`, `ai_disabled`…); `dominio.Forbidden(code, msg)` 403; `dominio.Unavailable(msg)` 503; `dominio.BadRequest(msg)` 400. Para lo demás: `&api.Problem{Status, Code, Message, Fields, RetryAt}` (p. ej. 413, 415, 429 con `RetryAt`). Cualquier otro error = 500 `internal_error`; su texto va solo al registro.
- `code` en inglés estable de la tabla de C1 §3; `message` en **español de México**, lenguaje simple (qué pasó y qué hacer). Nunca en un mensaje: contraseñas, tokens, cookies, `chat_id` completo.
- Tipos de entrada y salida: los de `types_gen.go` (`api.ExpenseCreate`, `api.Expense`, `api.ExpensePage`…). Reglas del generador: requerido → valor; opcional → puntero (`api.Ptr(x)`) u `omitzero`; tipo con `null` → `api.Nullable[T]` (tres estados: ausente, `null`, valor; `api.Some(v)`, `api.Null[T]()`), que también sirve para JSON Merge Patch (`PATCH`: `Set=false` no se toca, `Set && !Valid` borra). Columna NULL en la base → campo omitido en la respuesta (C1 §1 punto 6).
- Dinero: `int64` en centavos siempre; redondeo `dominio.RoundHalfUp`, `MulDivRound`, `ApplyBasisPoints`, `PercentOf` (ROUND_HALF_UP). Las series de pagos («el último ajusta centavos») son de los motores. `dominio.FormatPesos` solo para mensajes.
- Fechas: `dominio.Date` (`ParseDate` estricto, `String()` AAAA-MM-DD), `dominio.Today(clock)`, `FortnightOf`, `Fortnight.PayDay()` (15 y 30, febrero el último día), `ParsePeriod` (filtro `period`). Instantes: la base guarda UTC con `dominio.FormatDBTime(now)`; la respuesta usa `api.APITime(columna)` (−06:00). Entradas con «Z» se aceptan (`dominio.ParseDateTime`).

## 4. Entrada: decodificar y validar

- `reg.DecodeJSON(r, &in)`: JSON inválido, tipo incorrecto (`973.75` o `"973.75"` en un `*_cents`), campo desconocido o datos de sobra → **400** `bad_request` (con `fields` cuando se sabe el campo); más de 8 MiB → 413. Es la regla de C1 §3: el formato es 400, la regla de negocio es 422.
- `var v api.Validator`; `v.Required(ok, "campo")`, `v.PositiveCents(c, "amount_cents")`, `v.Check(cond, campo, code, msg)`, `v.Date(s, campo)` (400 si no existe); al final `return v.Err()` (nil o 422 con todos los campos).
- Parámetros: `api.PathID(r, "expense_id")`, `api.QueryIDs(r, "person_id")`, `api.QueryPeriod(r)`, `api.QueryDate(r, "from")`.

## 5. Colecciones (C1 §5)

```go
p, err := api.ParsePage(r, []string{"expense_date", "amount_cents"}, "-expense_date") // orden por defecto: el de la descripción en C1
q := `SELECT … FROM expense e WHERE … ` + p.OrderBy(map[string]string{"expense_date": "e.expense_date", "id": "e.id"}) + " " + p.LimitOffset()
// … leer filas en items (LimitOffset pide una de más) …
return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p)) // {"items": [...], "next_cursor": "…" | null}
```

`limit` 50 por defecto, de 1 a 200 (fuera → 400); `sort` solo con las propiedades que el área permita (el `id` cierra el orden); el cursor es opaco (desplazamiento + huella de orden y filtros: un cursor de otra consulta → 400).

## 6. Idempotency-Key (C1 §1 punto 8; K6, K9)

- Automática en las **44** operaciones con `Idempotent: true` en `api.Operations`: falta o formato malo → 400; misma clave + mismo cuerpo → la respuesta guardada con `Idempotent-Replayed: true` (201 se repite como **200** si la operación declara 200); otro cuerpo → 409 `idempotency_key_reused`; en curso → 409 `conflict`. Solo se guardan respuestas 2xx: un 4xx libera la clave. Retención: `idempotency.retention_hours` (24) y el job `idempotency_purge`.
- **K9 (gastos y pesos)**: `createExpense` y `createWeightEntry` guardan además `api.IdempotencyKey(r)` en `expense.idempotency_key` / `weight_entry.idempotency_key`; si el renglón ya existe (pasadas las 24 h), responden `api.MarkReplayed(w)` + **200** con el registro guardado, sin crear otro.
- El agente externo (bearer, K53) no tiene cuenta: su clave no se guarda en `idempotency_key` (PROVISIONAL, bloqueo `S5-agente-externo`).

## 7. Sesión y autenticación

- Middleware automático según `Operations[i].Auth`: `session` (cookie `fh_session`, renueva 90 días con cada uso y reemite la cookie), `public` (`login`, `getHealth`), `service` (`Authorization: Bearer`; `Deps.Service == nil` → 401, PROVISIONAL).
- `api.Session(r)` (`SessionID`, `UserID`, `PersonID`, `Username`) y `api.Actor(r)` (para `audit_log`). Axel y Ximena tienen los mismos permisos: **nada filtra por quién inició sesión** (solo valores por defecto).
- Para auth-sistema: `dominio.NewSessionToken`, `CreateSession`, `RevokeSession`, `RevokeOtherSessions`, `ListSessions`, `api.SetSessionCookie` / `ClearSessionCookie` / `SessionMaxAge`; `dominio.HashPassword(pw, dominio.DefaultArgon2)` / `VerifyPassword` (argon2id PHC); `dominio.CheckLockout(ctx, q, usuario, ip, now, max, ventana)` + `RecordLoginAttempt`. **Un intento rechazado por bloqueo (429) no se anota como fallo**, si no la ventana deslizante nunca se libera (S6-06 b). `retry_at` = `Lockout.RetryAt`.

## 8. Datos, transacciones y bitácora

- Repositorios: funciones que reciben `dominio.Querier` (`*sql.DB` o `*sql.Tx`). Escrituras de varios renglones: `dominio.WithTx(ctx, db, func(tx *sql.Tx) error {…})`.
- **Cada alta, cambio o baja escribe `audit_log`** (no hay disparadores): `before, _ := dominio.RowJSON(ctx, tx, "expense", id)` … `dominio.Audit(ctx, tx, api.Actor(r), "update", "expense", id, before, after, api.RequestID(r))`.
- `created_at`/`updated_at` los pone la base; si una regla compara con el reloj de la API, escribe el instante con `dominio.FormatDBTime(clock.Now())`. Violaciones de CHECK o disparadores: `dominio.IsConstraint(err)` (→ 422 o 409 según la regla), `dominio.IsUnique(err)` (→ 409 `already_exists`).
- Esquema: 1.4 (00001–00004). **Nunca se edita una migración aplicada**: cambio = migración nueva `0000N_*.sql` y un bloqueo `docs/estado/bloqueos/S6-<tema>.md`.

## 9. Ajustes (`app_setting`, K44)

- Leer: `dominio.SettingInt/IntList/Time/Bool(ctx, q, "clave")`. GET de settings: `dominio.ReadParameter(ctx, q, clave)` → `value`, `initial_value`, `is_initial`, `value_type`, `unit` (sale de las anotaciones `@param` de las migraciones, igual que C1 §11.7).
- PATCH de settings: `dominio.ApplyParams(ctx, tx, dominio.ParamChange{Set, Reset}, clavesDeLaSección, api.Actor(r), api.RequestID(r))`: valida tipo, clave de la sección, «cambiar y restablecer a la vez» (422) y el rango de los disparadores (422 `out_of_range`); escribe `audit_log`. Claves sin `@param` (`ai.*`, `optimizer.weights`): `dominio.SetSetting`.

## 10. Jobs

En `Register`: `reg.Deps().Jobs.Register(jobs.Def{Code: "daily_summary", Name: "Resumen diario", Schedule: "30 21 * * *", Run: f})` (antes de `Start`, que corre `main`). `Schedule` es solo el valor inicial del renglón `job` (Ajustes lo cambia con `updateJob`); `f(ctx, run)` recibe `run.Scheduled` (America/Mexico_City) y `run.RequestID` para `audit_log`. Un disparo corre una sola vez; un job vencido tras un reinicio corre una vez. «Correr ahora»: `Jobs.RunNow(ctx, code)`. Jobs base: `backup_daily` (03:15, conserva 30; PROVISIONAL) e `idempotency_purge` (cada hora al minuto 5).

## 11. IA

`p, err := reg.Deps().IA.Active(ctx)`; `errors.Is(err, ia.ErrDisabled)` → `dominio.Conflict("ai_disabled", "La IA está apagada. Puedes encenderla en Ajustes.")`. Apagada = cero llamadas de salida. Interruptores: `ai.local_model_enabled`, `ai.claude_api_enabled` (`ia.SettingKey(code)`).

## 12. Pruebas (S6-02, S6-03, S6-04)

- `s := testutil.NewServer(t, gastos.Register)` (TLS, base en `datos/tmp/` migrada, fixtures Axel, Ximena y Karen sin cuenta, reloj en 12-oct-2026 12:00 `s.Clock`); `c := s.SessionFor(t, testutil.AxelUser)`; `res := c.Do("POST", "/api/v1/…", cuerpo, "Idempotency-Key", "…")` **valida la respuesta contra `openapi.yaml`** (código declarado + esquema; error = prueba fallida). `c.Raw` no valida (para cuerpos malos a propósito). `s.Anonymous(t)` sin sesión. `testutil.NewDB(t)` para pruebas de repositorio.
- **Una prueba por operación** cuyo nombre lleve el operationId: `func TestOp_listUsers(t *testing.T)` (así sale la tabla operación → prueba de S6-03 con `go test -v`). Cubrir 2xx y los 4xx propios; nada de comprobar solo el 200.
- Datos de ejemplo de interfaz: `docs/design/maquetas/datos-sinteticos.json`; nunca credenciales reales (contraseña falsa `testutil.FakeCredential`). Temporales solo con `testutil.TempDir(t)`.

## 13. Ejemplo completo (una operación con su prueba)

```go
// internal/auth/usuarios.go
package auth

func Register(reg *api.Registry) { reg.Handle("listUsers", listUsers(reg)) }

func listUsers(reg *api.Registry) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		p, err := api.ParsePage(r, []string{"username"}, "username")
		if err != nil {
			return err
		}
		us, err := dominio.ListUsers(r.Context(), reg.Deps().DB)
		if err != nil {
			return err
		}
		items := make([]api.UserAccount, 0, len(us))
		for _, u := range us {
			it := api.UserAccount{ID: u.ID, PersonID: u.PersonID, Username: u.Username, IsActive: u.IsActive,
				CreatedAt: api.APITime(u.CreatedAt), UpdatedAt: api.APITime(u.UpdatedAt), DisplayName: api.Ptr(u.DisplayName)}
			if u.LastLoginAt.Valid {
				it.LastLoginAt = api.Ptr(api.APITime(u.LastLoginAt.String))
			}
			items = append(items, it) // nunca password_hash
		}
		return api.WriteJSON(w, http.StatusOK, api.PageOf(items, p))
	}
}

// internal/auth/usuarios_test.go
func TestOp_listUsers(t *testing.T) {
	s := testutil.NewServer(t, auth.Register)
	res := s.SessionFor(t, testutil.XimenaUser).Do("GET", "/api/v1/users", nil) // valida contra openapi.yaml
	var pg api.UserAccountPage
	res.JSON(t, &pg)
	if res.StatusCode != 200 || len(pg.Items) != 2 || pg.NextCursor.Valid {
		t.Fatalf("%d %+v", res.StatusCode, pg)
	}
	if s.Anonymous(t).Do("GET", "/api/v1/users", nil).StatusCode != 401 {
		t.Error("sin sesión debe ser 401")
	}
}
```

## 14. Reparto de operaciones (la verdad para WF-3a; 216 de 216)

| Agente | Operaciones |
|---|---|
| api-gastos-inicio (26), paquete `internal/api/gastos-inicio` | `gastos.yaml` sin `uploadExpenseAttachment`, `downloadAttachment`, `deleteAttachment`; `inicio.yaml` sin `getHealth`, `getCurrentSession` |
| api-presupuesto (28) | `presupuesto-a.yaml`, `presupuesto-b.yaml` |
| api-comidas-salud (26) | `comidas.yaml`, `salud.yaml`, `createHealthProfile`, `updateHealthProfile`, `getHealthSettings`, `updateHealthSettings` |
| api-compras-a (18) | `compras-tiendas-a.yaml`, `createStore`, `updateStore`, `setStoreProductGroups` |
| api-compras-b (23) | `compras-tiendas-b.yaml` (incluye las 4 del agente externo, K53: `getExternalAgentReference`, `importExternalCatalog`, `importExternalPrices`, `importExternalRecipes`), `getPriceSources`, `getPriceSettings`, `updatePriceSettings`, `importProfeco`, `getRecipeSettings`, `reloadNutritionBase`, `updateInventoryItem` |
| api-creditos (29) | `creditos-a.yaml`, `creditos-b.yaml`, `createCredit`, `updateCredit`, `createCreditTermRate`, `updateCreditTermRate`, `deleteCreditTermRate`, `createPromotion`, `updatePromotion`, `deletePromotion`, `getOptimizerSettings`, `updateOptimizerSettings` |
| api-ajustes (19) | `listCategories`, `createCategory`, `updateCategory`, `reorderCategories`, `listPersons`, `updatePerson`, `listCostShareRules`, `updateCostShareRule`, `listPlannedItems`, `createPlannedItem`, `updatePlannedItem`, `previewPlannedItemsImport`, `importPlannedItems`, `getProjectionSettings`, `updateProjectionSettings`, `createSavingsGoal`, `updateSavingsGoal`, `getAlertSettings`, `updateAlertSettings` |
| auth-sistema (27) | `login`, `logout`, `changeOwnPassword`, `listUsers`, `listSessions`, `revokeSession`, `revokeSessions`, `getCurrentSession`, `getHealth`, `getOpenApiDocument`, `listBackups`, `createBackup`, `restoreBackup`, `getSeedStatus`, `getSecuritySettings`, `updateSecuritySettings`, `getDataSettings`, `updateDataSettings`, `listAiProviders`, `updateAiProvider`, `getSystemStatus`, `getCaCertificate`, `uploadExpenseAttachment`, `downloadAttachment`, `deleteAttachment`, `listJobs`, `updateJob` |
| s7-avisos (10) | `listNotifications`, `updateNotification`, `setNotificationRecipients`, `listNotificationRecipients`, `updateNotificationRecipient`, `sendTestNotification`, `listTelegramChatCandidates`, `confirmTelegramChatCandidate`, `dismissTelegramChatCandidate`, `listNotificationDeliveries` |
| s7-pdf (10) | `listReports`, `generateMonthlyReport`, `generateExternalStatement`, `downloadReportFile`, `sendReport`, `listMonthlyReportMonths`, `getMonthlyReport`, `getExternalStatement` y **`getReportSettings`, `updateReportSettings`** (sin dueño en el reparto del orquestador; asignadas aquí por backend-base: sección «Reportes y frecuencias», claves `reports.*` que usan los reportes programados de S7. PROPUESTA: el orquestador puede moverlas a s7-avisos) |

## 15. Reglas que no cambian

Sin red en ejecución (ni CDN ni telemetría); identificadores en inglés y textos al usuario en español de México; nada de credenciales reales en código, pruebas, configuración o reportes; el token del bot nunca se imprime; temporales en `datos/tmp/`; valores operativos sin decisión = PROVISIONAL con justificación y bloqueo; verificación: `make verify-wf3` (y `make verify-all GOOSE=…` no se rompe).
