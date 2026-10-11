# Generador C1 · parte 6: llena C1-datos.md desde openapi.json, U1 y el esquema de S1
import json, re, os
from c1paths import ROOT, OUT, WORK, HERE
d = json.load(open(os.path.join(WORK, 'openapi.json'), encoding='utf-8'))
schema = json.load(open(os.path.join(WORK, 'schema.json'), encoding='utf-8'))
u1 = open(os.path.join(ROOT, 'docs/design/U1-arquitectura.md'), encoding='utf-8').read().split('\n')
# --- necesidades de U1 ---
NEEDS = {}
for l in u1:
    m = re.match(r'^\| (N-\d{3}) \| ([LE]) \| (.+)$', l)
    if m:
        cols = [c.strip() for c in m.group(3).split(' | ')]
        name = cols[1] if re.fullmatch(r'\d{2}', cols[0]) else cols[0]
        NEEDS[m.group(1)] = (m.group(2), name)
nums = sorted(int(n[2:]) for n in NEEDS)
def expand(cell):
    out = []
    for a, b in re.findall(r'N-(\d{3})(?: a N-(\d{3}))?', cell):
        if b: out += [f'N-{n:03d}' for n in nums if int(a) <= n <= int(b)]
        else: out.append(f'N-{a}')
    return out
SCREENS = []
i0 = next(i for i, l in enumerate(u1) if l.startswith('### 3.29'))
for l in u1[i0:]:
    if l.startswith('### 3.30'): break
    m = re.match(r'^\| (Global \(G-01 a G-07\)|P-[\d.]+) \| (.*?) \| (.*?) \| (.*?) \|$', l)
    if m:
        scr = 'Global' if m.group(1).startswith('Global') else m.group(1)
        SCREENS.append((scr, expand(m.group(2)), expand(m.group(3))))
# --- operaciones por necesidad ---
OPS = []
for path, ops in d['paths'].items():
    for meth, o in ops.items():
        OPS.append(dict(m=meth.upper(), path=path, oid=o['operationId'], tag=o['tags'][0], needs=o.get('x-necesidades', []), tables=o.get('x-tablas', []), params=o.get('parameters', []), public=o.get('security') == [], idem=any(p.get('$ref', '').endswith('IdempotencyKey') for p in o.get('parameters', [])), desc=o.get('description', '')))
BYN = {}
for o in OPS:
    for n in o['needs']: BYN.setdefault(n, []).append(o)
def opstr(o): return f"{o['m']} `{o['path']}`"
LOCAL = {'N-003': 'Recientes en el dispositivo (U1: Local)', 'N-301': 'Tema por dispositivo (U1: Local)', 'N-317': 'Cola sin conexión en el teléfono (U1: Local; envía con N-054 y N-240, K9)'}
# Corrección r2: con el esquema 1.1 de S1 ninguna necesidad queda sin endpoint por bloqueo. La lista de compras sigue la opción A
# C1 v1.1 (puerta 1): K1 a K23 cerraron la lista de compras y los huecos U1 ↔ S1; quedan con cobertura parcial solo las necesidades que
# dependen de un bloqueo abierto (valores operativos K31, redondeo del promedio semanal y ajuste por objetivo de la meta de kcal).
# Corrección r4 (H-r3-02): N-120 por la semana por omisión de getMealPlanWeek (regla del domingo, bloqueo U1-lista-vigente-domingo).
BLOCKED = {}
# C1 v1.2 (puerta 1b): K36, K38, K39 y K43 cerraron esas marcas; quedan FIRMES citando su decisión. N-297 y N-300 son FIRMES
# (claves de K43), pero dónde viven send_time y advance_days sigue PROVISIONAL por el bloqueo abierto S1-parametros-avisos.
PARTIAL = {}
FIRM = {'N-023': 'K43', 'N-064': 'K43', 'N-226': 'K43', 'N-297': 'K43', 'N-300': 'K43', 'N-042': 'K39', 'N-230': 'K36', 'N-311': 'K36', 'N-313': 'K36', 'N-120': 'K38'}
PLACE = {'N-297': 'S1-parametros-avisos', 'N-300': 'S1-parametros-avisos'}
def status(n):
    if n in LOCAL: return 'Local'
    if n in BLOCKED: return f'BLOQUEO `{BLOCKED[n]}`'
    if n in PARTIAL: return f'PROVISIONAL `{PARTIAL[n]}`'
    if n in PLACE: return f'FIRME ({FIRM[n]}) · dónde viven `send_time` y `advance_days`: PROVISIONAL `{PLACE[n]}`'
    if n in FIRM: return f'FIRME ({FIRM[n]})'
    return 'OK' if BYN.get(n) else 'SIN ENDPOINT'
rows = ['| Pantalla | Necesidad | Tipo | Qué necesita (U1) | Operación(es) C1 | Estado |', '|---|---|---|---|---|---|']
for scr, rd, wr in SCREENS:
    for n in rd + wr:
        t, name = NEEDS.get(n, ('?', '?'))
        ops = BYN.get(n, [])
        cell = ' · '.join(opstr(o) for o in ops) if ops else (LOCAL.get(n) or 'sin endpoint (ver estado)')
        rows.append(f'| {scr} | {n} | {t} | {name} | {cell} | {status(n)} |')
MATRIZ = '\n'.join(rows)
# --- resumen y huecos ---
nset = set(NEEDS); mapped = {n for n in nset if BYN.get(n)}
sin_ep = sorted(nset - mapped)
sin_n = [o for o in OPS if not o['needs']]
unknown_n = sorted({n for o in OPS for n in o['needs']} - nset)
RESUMEN = '\n'.join([
 f'- Necesidades de U1: **{len(nset)}** (N-001 a {max(nset)}). Con operación: **{len(mapped)}**. Locales sin API (U1): **{len(set(LOCAL) & nset)}** ({", ".join(sorted(set(LOCAL) & nset))}). Sin endpoint por bloqueo: **{len(set(BLOCKED) & nset)}**{(" (" + ", ".join(sorted(BLOCKED)) + ")") if BLOCKED else ""}. Con cobertura parcial por un bloqueo abierto de §13: **{len(set(PARTIAL) & mapped)}**. FIRMES desde la puerta 1b (antes con marca): **{len(set(FIRM) & mapped)}**; en {len(PLACE)} de ellas solo queda pendiente dónde viven `send_time` y `advance_days` (bloqueo abierto `S1-parametros-avisos`, §13).',
 f'- Operaciones del contrato: **{len(OPS)}** en **{len(d["paths"])}** rutas. Con al menos una N-xxx: **{len(OPS) - len(sin_n)}**. Excepciones técnicas sin N-xxx: **{len(sin_n)}** (§11.2). N-xxx citadas por el contrato que no existen en U1: **{len(unknown_n)}**.',
 f'- Pantallas de la matriz de U1 §3.29: **{len(SCREENS)}** (Global, P-00 a P-25 y P-26.01 a P-26.14); filas de la matriz: **{len(rows) - 2}**.',
 f'- Tablas de S1: **{len(schema)}**; con al menos una operación: **{sum(1 for t in schema if any(t in o["tables"] for o in OPS))}** (§11.4). `idempotency_key` es infraestructura de S6: la declaran en `x-tablas` las **{sum(o["idem"] for o in OPS)}** operaciones con `Idempotency-Key` (§1 punto 8).'])
SIN_N = '\n'.join(['| Operación | Motivo |', '|---|---|'] + [f"| {opstr(o)} (`{o['oid']}`) | Estático y técnico: sirve este mismo contrato para desarrollo y para la prueba de contrato de S6 (S6-03). Requiere sesión. Excepción «estáticos» de C1-02. |" for o in sin_n]) + '\n\nNota: `POST /api/v1/auth/login` (N-010, N-011) y `GET /health` (N-004, franja «Sin conexión») sí tienen necesidad en U1, aunque C1-02 los admite como excepciones.'
N_SIN_EP = '\n'.join(['| Necesidad | Qué es | Motivo |', '|---|---|---|'] + [f'| {n} | {NEEDS[n][1]} | ' + (LOCAL[n] if n in LOCAL else f'Sin endpoint por un bloqueo abierto: `{BLOCKED.get(n, "?")}`') + ' |' for n in sin_ep])
# --- filtros ---
FN = {'person_id': 'persona', 'family': 'familia', 'category_id': 'categoría', 'store_id': 'tienda', 'period': 'periodo', 'credit_id': 'crédito'}
REFN = {'PersonId': 'person_id', 'Family': 'family', 'CategoryId': 'category_id', 'StoreId': 'store_id', 'Period': 'period', 'CreditId': 'credit_id'}
# H-r2-01 (revisión C1-r2): la última columna marca la llave `week` de la lista de compras (ListWeek), que no es filtro.
fr = ['| Operación | Ruta | persona | familia | categoría | tienda | periodo | crédito | semana de lista (`week`, llave) |', '|---|---|---|---|---|---|---|---|---|']
for o in OPS:
    refs = {p['$ref'].split('/')[-1] for p in o['params'] if '$ref' in p}
    have = {REFN[r] for r in refs if r in REFN}
    if have or 'ListWeek' in refs: fr.append(f"| `{o['oid']}` | {o['m']} `{o['path']}` | " + ' | '.join('✓' if k in have else '' for k in FN) + ' | ' + ('✓' if 'ListWeek' in refs else '') + ' |')
FILTROS = '\n'.join(fr)
# --- entidades (31) ---
ENT = [('persona', 'person', 'listPersons, updatePerson, getPersonParticipation', 'Tabla* (P-26.12)'),
 ('cuenta_usuario', 'user_account', 'listUsers, login, getCurrentSession, changeOwnPassword (sin alta: Karen no tiene cuenta)', 'Tabla* (P-26.01)'),
 ('sesion_dispositivo', 'device_session', 'listSessions, revokeSession, revokeSessions, logout', 'Tabla* (P-26.01)'),
 ('periodo', 'period', 'listPeriods (selector); motor S2 (meses y quincenas de los cálculos; el presupuesto usa vigencia por fechas)', 'Sin lista: dato derivado (U1)'),
 ('ingreso', 'income + planned_item', 'listIncomes, getIncomeMonth, createIncome, updateIncome, deleteIncome, getIncomeMonthlyTrend; listPlannedItems, createPlannedItem, updatePlannedItem', 'Tabla*, calendario (income_date)'),
 ('categoria', 'category', 'listCategories, createCategory, updateCategory, reorderCategories', 'Tabla*, tablero (default_cost_share_rule_id)'),
 ('gasto', 'expense + expense_item + expense_share + attachment', 'listExpenses, getExpense, createExpense, updateExpense, deleteExpense, duplicateExpense, markExpensePaid, listExpenseItems, uploadExpenseAttachment, downloadAttachment, deleteAttachment, ticket-scans', 'Tabla*, tablero (category_id), calendario (expense_date), galería (tickets)'),
 ('regla_reparto', 'cost_share_rule + cost_share_rule_part', 'listCostShareRules, updateCostShareRule; reglas de un solo uso al crear gasto o renglón', 'Sin lista: propiedad de la categoría y del gasto (U1)'),
 ('meta_ahorro', 'savings_goal', 'listSavingsGoals, createSavingsGoal, updateSavingsGoal, getSavingsSummary', 'Tabla (P-26.13), tablero (status), calendario (target_date), galería* (P-09)'),
 ('movimiento_ahorro', 'savings_movement', 'listSavingsMovements, createSavingsMovement, updateSavingsMovement, deleteSavingsMovement', 'Tabla*, calendario (movement_date)'),
 ('credito', 'credit + credit_statement', 'listCredits, getCredit, createCredit, updateCredit, getCreditTotals, previewAvailability, getAvailabilitySeries, getCardStatus, getDebtTrend', 'Tabla (P-26.05), tablero (kind), galería* (P-17)'),
 ('tasa_por_plazo', 'credit_term_rate', 'listCreditTermRates, createCreditTermRate, updateCreditTermRate, deleteCreditTermRate', 'Tabla*'),
 ('compra_credito', 'credit_purchase + credit_purchase_share', 'listCreditPurchases, getCreditPurchase, createCreditPurchase, previewCreditPurchase, updateCreditPurchase, deleteCreditPurchase', 'Tabla*, tablero (status), calendario (purchase_date)'),
 ('pago_programado', 'scheduled_payment + scheduled_payment_share + credit_payment', 'listScheduledPayments, markScheduledPaymentPaid, deleteCreditPayment, listCalendarEvents', 'Tabla, tablero (status), calendario* (due_date)'),
 ('promocion', 'promotion', 'listPromotions, createPromotion, updatePromotion, deletePromotion', 'Tabla, tablero (credit_id), calendario* (valid_from a valid_to)'),
 ('tienda', 'store + store_product_group + product_group', 'listStores, createStore, updateStore, setStoreProductGroups, listProductGroups', 'Tabla*, tablero (kind), calendario (visit_weekday, recurrente), galería'),
 ('producto', 'product', 'listProducts, getProduct, createProduct, updateProduct', 'Tabla*, tablero (freshness), galería'),
 ('presentacion', 'presentation', 'createPresentation, updatePresentation; dentro de getProduct', 'Sin lista propia: tabla dentro del producto (U1)'),
 ('precio', 'price', 'listProductPrices, createPrice, updatePrice, deletePrice, getPriceSummary, importProfeco', 'Tabla*, tablero (freshness)'),
 ('receta', 'recipe + recipe_meal_slot', 'listRecipes, getRecipe, createRecipe, updateRecipe, duplicateRecipe, deleteRecipe, getRecipeUsage, suggestRecipeWithAi', 'Tabla, tablero (tiempos aptos, recipe_meal_slot), galería*'),
 ('ingrediente_receta', 'recipe_ingredient', 'dentro de getRecipe, createRecipe y updateRecipe', 'Sin lista propia: tabla dentro de la receta (U1)'),
 ('tiempo_comida', 'meal_slot', 'listMealSlots (solo lectura)', 'Sin lista: catálogo fijo de 4 (U1)'),
 ('plan_comida', 'meal_plan_entry', 'getMealPlanWeek, getMealPlanMonth, getMealPlanTrends, createMealPlanEntries, updateMealPlanEntry, deleteMealPlanEntry, copyPreviousWeek', 'Tabla, tablero (meal_slot_id), calendario* (plan_date), galería (recipe_id)'),
 ('inventario_item', 'inventory_item + inventory_movement', 'listInventoryItems, createInventoryItem, updateInventoryItem, adjustInventoryItem, listInventoryMovements, getInventorySummary', 'Tabla*, tablero (status), galería'),
 ('perfil_salud', 'health_profile', 'listHealthProfiles, createHealthProfile, updateHealthProfile, getTodayKcal', 'Sin lista: una ficha por persona (U1)'),
 ('peso_historico', 'weight_entry', 'listWeightEntries, createWeightEntry, updateWeightEntry, deleteWeightEntry', 'Tabla*'),
 ('nutricion_ingrediente', 'nutrition_fact', 'listNutritionFacts, createNutritionFact, reloadNutritionBase, getRecipeSettings', 'Sin lista: propiedad del producto (U1)'),
 ('notificacion', 'notification + notification_delivery', 'listNotifications, updateNotification, sendTestNotification, listNotificationDeliveries', 'Tabla*, tablero (notification_id), calendario (sent_at) en avisos enviados'),
 ('destinatario', 'notification_recipient + notification_subscription + telegram_chat_candidate', 'listNotificationRecipients, listTelegramChatCandidates, confirmTelegramChatCandidate, updateNotificationRecipient, setNotificationRecipients', 'Sin lista: se elige dentro de cada aviso (U1)'),
 ('job', 'job', 'listJobs, updateJob', 'Tabla*, tablero (code), calendario (next_run_at)'),
 ('respaldo', 'backup', 'listBackups, createBackup, restoreBackup', 'Tabla*')]
ENTIDADES = '\n'.join(['| # | Entidad (§3.4.1) | Tabla(s) S1 | Operaciones | Vistas (U1 §2.5.a; * = por defecto) |', '|---|---|---|---|---|'] + [f'| {i} | {e} | `{t}` | {ops} | {v} |' for i, (e, t, ops, v) in enumerate(ENT, 1)])
VIS = [('expense', 'table, board, calendar, gallery', 'category_id, kind, payment_method, person_id', 'expense_date', 'parte del hogar', 'N-043, N-045, N-046'),
 ('income', 'table, calendar', '—', 'income_date', 'amount_cents', 'N-104 (también listCalendarEvents)'),
 ('category', 'table, board', 'default_cost_share_rule_id', '—', '—', 'N-281'),
 ('savings_goal', 'table, board, calendar, gallery', 'status', 'target_date', 'saved_cents', 'N-111'),
 ('savings_movement', 'table, calendar', 'savings_goal_id, kind', 'movement_date', 'amount_cents (aportación +, retiro −)', 'N-115'),
 ('credit', 'table, board, gallery', 'kind', '—', 'live_debt_cents', 'N-190'),
 ('credit_term_rate', 'table', '—', '—', '—', 'N-286'),
 ('credit_purchase', 'table, board, calendar', 'status, credit_id', 'purchase_date', 'total_cents', 'N-202'),
 ('scheduled_payment', 'table, board, calendar', 'status, credit_id', 'due_date', 'amount_cents', 'N-192'),
 ('promotion', 'table, board, calendar', 'credit_id', 'valid_from (franja hasta valid_to)', '—', 'N-287'),
 ('store', 'table, board, calendar, gallery', 'kind', 'visit_weekday (recurrente)', '—', 'N-290'),
 ('product', 'table, board, gallery', 'freshness, product_group_id', '—', '—', 'N-170'),
 ('price', 'table, board', 'is_stale, source, store_id', 'observed_date', '—', 'N-171'),
 ('recipe', 'table, board, gallery', 'meal_slot_id (tiempos aptos de recipe_meal_slot: una receta con varios tiempos aparece en cada columna), origin', '—', '—', 'N-135'),
 ('meal_plan_entry', 'table, board, calendar, gallery', 'meal_slot_id, person_id', 'plan_date', '—', 'N-128'),
 ('inventory_item', 'table, board, gallery', 'status, tracking_mode', '—', '—', 'N-180'),
 ('notification_delivery', 'table, board, calendar', 'notification_id, status', 'sent_at (día local)', '—', 'N-247'),
 ('job', 'table, board, calendar', 'code, last_status', 'next_run_at (día local)', '—', 'N-246, N-302'),
 ('person', 'table', '—', '—', '—', 'N-309'), ('user_account', 'table', '—', '—', '—', 'N-270'), ('device_session', 'table', 'user_account_id', '—', '—', 'N-271'),
 ('weight_entry', 'table', 'person_id', 'measured_date', '—', 'N-232'), ('backup', 'table', 'status, kind', 'started_at (día local)', '—', 'N-302')]
VISTAS = '\n'.join(['| `entity` | Modos | `group_by` permitidos | `date_field` | Importe principal (`total_cents`) | Necesidades |', '|---|---|---|---|---|---|'] + [f'| `{a}` | {b} | {c} | {dd} | {e} | {f} |' for a, b, c, dd, e, f in VIS]) + '\n\nOtro valor de `group_by` o `date_field` responde 422 `validation_failed`. `queryView` declara en `x-necesidades` las necesidades cuya vista principal es de tablero, calendario o galería; las vistas de tabla de cada entidad también se pueden pedir aquí, aunque cada pantalla usa la colección tipada de su recurso.'
# --- tablas ---
tabla_ops = {t: [o['oid'] for o in OPS if t in o['tables']] for t in schema}
TABLAS = '\n'.join(['| Tabla S1 | Operaciones | Ejemplos |', '|---|---|---|'] + [f'| `{t}` | {len(v)} | {", ".join("`%s`" % x for x in v[:5])}{" …" if len(v) > 5 else ""} |' for t, v in tabla_ops.items()])
sin_tabla = [t for t, v in tabla_ops.items() if not v]
bad_tab = sorted({t for o in OPS for t in o['tables']} - set(schema))
# --- ajustes ---
SEC = {'P-26.01': 'Cuenta y roles', 'P-26.02': 'Presupuesto', 'P-26.03': 'Ingresos', 'P-26.04': 'Categorías', 'P-26.05': 'Créditos y condiciones (incluye promociones)', 'P-26.06': "Tiendas y precios (incluye mercados y membresía Sam's)", 'P-26.07': 'Recetas', 'P-26.08': 'Reportes y frecuencias', 'P-26.09': 'Notificaciones', 'P-26.10': 'Apariencia', 'P-26.11': 'Datos y respaldos', 'P-26.12': 'Personas y reparto', 'P-26.13': 'Salud y metas de ahorro', 'P-26.14': 'Inventario'}
aj = ['| Sección | Nombre | Necesidades | Operaciones |', '|---|---|---|---|']
for scr, rd, wr in SCREENS:
    if scr in SEC:
        ops = []
        for n in rd + wr:
            for o in BYN.get(n, []):
                if o['oid'] not in ops: ops.append(o['oid'])
        aj.append(f"| {scr} | {SEC[scr]} | {', '.join(rd + wr)} | {', '.join('`%s`' % x for x in ops) if ops else 'Sin API: el tema vive en el dispositivo (N-301, Local)'} |")
AJUSTES = '\n'.join(aj)
CAP = [('Las 14 secciones de Ajustes', '§11.5'), ('Sesiones por dispositivo y su cierre', 'listSessions, revokeSession, revokeSessions, logout'), ('Interruptor A/B y su vista previa', 'getCredit, updateCredit, previewAvailability'),
 ('Comparador de escenarios: 4 objetivos, restricciones, pesos y recomendación', 'compareScenarios, getOptimizerSettings, updateOptimizerSettings'), ('Promociones', 'listPromotions, createPromotion, updatePromotion, deletePromotion'), ('Tasa por plazo', 'listCreditTermRates, createCreditTermRate, updateCreditTermRate, deleteCreditTermRate'),
 ('Calendario unificado de pagos', 'listScheduledPayments, listCalendarEvents'), ('Créditos de los 4 tipos y sus condiciones', 'listCredits, createCredit, updateCredit, getCardStatus, listScheduledPayments (amortización y abonos)'), ('Compras a crédito con partes por persona', 'createCreditPurchase, previewCreditPurchase, getCreditPurchase, markScheduledPaymentPaid, deleteCreditPayment'),
 ('Plan de comidas', 'getMealPlanWeek, createMealPlanEntries, updateMealPlanEntry, deleteMealPlanEntry, copyPreviousWeek, getMealPlanMonth, getMealPlanTrends'), ('Lista de compras', 'getShoppingList, generateShoppingList, createShoppingListItem, updateShoppingListItem, assignShoppingListStore (compartida en el servidor, K1)'), ('Recetas e ingredientes', 'listRecipes, getRecipe, createRecipe, updateRecipe, duplicateRecipe, deleteRecipe, getRecipeUsage, suggestRecipeWithAi'),
 ('Perfil de salud y peso', 'listHealthProfiles, createHealthProfile, updateHealthProfile, getTodayKcal, listWeightEntries, createWeightEntry, updateWeightEntry, deleteWeightEntry'), ('Precios con fuente y frescura', 'listProductPrices, createPrice, updatePrice, deletePrice, getPriceSummary, getPriceSettings, getPriceSources'),
 ('OCR con confirmación humana', 'createTicketScan, listTicketScans, getTicketScan, confirmTicketScan'), ('Mercados con día de visita', 'listStores (kind=market), createStore, updateStore, setStoreProductGroups'), ("Membresía de Sam's (punto de equilibrio)", 'getBasketComparison (membership), updateStore (membership_annual_fee_cents)'),
 ('Gastos con reparto y adjuntos', 'createExpense, updateExpense, listCostShareRules, uploadExpenseAttachment, downloadAttachment, deleteAttachment'), ('Ingresos fijos y extra', 'createIncome, listIncomes, getIncomeMonth, updateIncome, deleteIncome, listPlannedItems'), ('Categorías', 'listCategories, createCategory, updateCategory, reorderCategories'),
 ('Metas y fondo de ahorro', 'listSavingsGoals, createSavingsGoal, updateSavingsGoal, getSavingsSummary, createSavingsMovement'), ('Tiendas, productos, presentaciones y equivalencias', 'listStores, listProducts, createProduct, getProduct, createPresentation, updatePresentation'), ('Comparador de canasta', 'listBaskets, getBasketComparison, assignShoppingListStore'),
 ('Inventario', 'listInventoryItems, createInventoryItem, updateInventoryItem, adjustInventoryItem, listInventoryMovements, getInventorySummary'), ('Nutrición por ingrediente', 'listNutritionFacts, createNutritionFact, reloadNutritionBase'), ('Notificaciones y destinatarios', 'listNotifications, updateNotification, setNotificationRecipients, listNotificationRecipients, listTelegramChatCandidates, confirmTelegramChatCandidate, listNotificationDeliveries'),
 ('Disparo manual del resumen diario y de una alerta', 'sendTestNotification'), ('Reporte mensual PDF (generar, listar, descargar, reenviar)', 'generateMonthlyReport, listReports, downloadReportFile, sendReport'), ('Estado de cuenta de Karen bajo demanda', 'getExternalStatement, generateExternalStatement'), ('Vista web mensual', 'getMonthlyReport, listMonthlyReportMonths'),
 ('Jobs', 'listJobs, updateJob'), ('Respaldos (listar, ejecutar, restaurar)', 'listBackups, createBackup, restoreBackup'), ('IA opcional apagada por defecto', 'listAiProviders, updateAiProvider'), ('Vistas de base de datos genéricas', 'queryView'), ('Métricas del dashboard', 'getDashboardMonthVsPlan, getDashboardFortnightFlow, getDashboardNextPayment, getDashboardWeekMeals, getDashboardBestBasket, getDashboardSavings')]
allids = {o['oid'] for o in OPS}
cap_bad = [x for _, ops in CAP for x in re.findall(r'\b([a-z]+[A-Z][A-Za-z]*|login|logout|search)\b', ops) if re.match(r'^(list|get|create|update|delete|preview|mark|compare|upload|download|duplicate|confirm|revoke|reorder|set|send|generate|import|reload|restore|adjust|copy|suggest|logout|login|change|query|search)', x) and x not in allids]
CAPACIDADES = '\n'.join(['| # | Capacidad | Operaciones |', '|---|---|---|'] + [f'| {i} | {c} | {ops} |' for i, (c, ops) in enumerate(CAP, 1)])
OPERACIONES = '\n'.join(['| Método y ruta | operationId | Etiqueta | N-xxx | Tablas S1 (`x-tablas`) |', '|---|---|---|---|---|'] + [f"| {opstr(o)} | `{o['oid']}` | {o['tag']} | {', '.join(o['needs']) or '—'} | {', '.join(o['tables']) or '— (archivo, no base)'} |" for o in sorted(OPS, key=lambda x: (x['path'], x['m']))])
# C1 v1.2 (puerta 1b, K44): §11.7, tabla de parámetros operativos. Sale de las anotaciones @param y @param_sin_valor de las
# migraciones (params.json de parse_schema.py, que ya comparó cada valor inicial con la siembra) y de los esquemas de settings
# de openapi.json: cada clave tiene que tener una sola lectura (GET con parameters) y un solo cambio con «restablecer».
PRM = json.load(open(os.path.join(WORK, 'params.json'), encoding='utf-8'))
SCH = d['components']['schemas']
def _pk(ref_, campo):
    s = SCH.get((ref_ or '').split('/')[-1], {}).get('properties', {}).get(campo, {})
    return set(s.get('properties', {})) if campo == 'parameters' else set(s.get('items', {}).get('enum', []))
prd, pwr = {}, {}
for _p, _item in d['paths'].items():
    for _m, _o in _item.items():
        if not isinstance(_o, dict) or 'operationId' not in _o: continue
        _rs = _o.get('responses', {}).get('200', {}).get('content', {}).get('application/json', {}).get('schema', {}).get('$ref')
        _rq = _o.get('requestBody', {}).get('content', {}).get('application/json', {}).get('schema', {}).get('$ref')
        if _m == 'get':
            for k in _pk(_rs, 'parameters'): prd.setdefault(k, []).append(_o['operationId'])
        if _m == 'patch':
            for k in _pk(_rq, 'parameters') & _pk(_rq, 'reset_parameters'): pwr.setdefault(k, []).append(_o['operationId'])
_pb = [p['key'] for p in PRM['params'] if len(prd.get(p['key'], [])) != 1 or len(pwr.get(p['key'], [])) != 1]
assert not _pb, f'parámetros K44 sin una lectura y un cambio: {_pb}'
PARAMS_T = '\n'.join(['| Clave de `app_setting` | Tipo | Unidad | Valor inicial | Decisión | Sección de Ajustes (propuesta de K44) | Sembrada en | Lectura | Cambio y «restablecer» |', '|---|---|---|---|---|---|---|---|---|'] +
    [f"| `{p['key']}` | {p['type']} | {p['unit']} | `{p['initial_json']}` | {p['decision']} | {p['section']} | {p['migration']} | `{prd[p['key']][0]}` | `{pwr[p['key']][0]}` |" for p in PRM['params']])
PARAMS_SV = '\n'.join(['| Parámetro | Para qué | Por qué no tiene valor sembrado aquí |', '|---|---|---|'] + [f"| {p['name']} | {p['purpose']} | {p['reason']} |" for p in PRM['params_sin_valor']])
# C1 v1.2 · corrección r1 (A-01 de docs/estado/p1c/REVISION-C1-r1.md, C1-23): validación de send_time y advance_days (K43) tal como
# la declaran las propiedades de params_json en openapi; Notification y NotificationUpdate tienen que traer las mismas.
_pjk = {n: SCH[n]['properties']['params_json'].get('properties', {}) for n in ('Notification', 'NotificationUpdate')}
assert _pjk['Notification'] == _pjk['NotificationUpdate'] and set(_pjk['Notification']) == {'send_time', 'advance_days'}, _pjk
def _rng(v): return ', '.join(f'`{k}: {v[k]}`' for k in ('type', 'pattern', 'minimum') if k in v) + (f", valor por omisión `{v['default']}`" if 'default' in v else ', sin valor por omisión')
PARAMS_K43 = '\n'.join(['Validación de las claves de K43 en `params_json` de `Notification` y `NotificationUpdate` (`updateNotification`). El tipo y el rango son FIRMES (K43); solo su lugar es PROVISIONAL, bloqueo `S1-parametros-avisos` (§13):', ''] +
    [f"- `{k}`: {_rng(v)}." for k, v in _pjk['NotificationUpdate'].items()] +
    ['- `weekday` y `day_of_month` (decisión F) no se declaran como propiedades: entran por `additionalProperties: true` con el valor que sembró la 00001.'])
INCONS = open(os.path.join(HERE, 'c1_inconsistencias.md'), encoding='utf-8').read().strip()
t = open(os.path.join(HERE, 'c1_plantilla.md'), encoding='utf-8').read()
repl = {'N_TABLAS': str(len(schema)), 'N_TABLAS_OP': str(sum(1 for v in tabla_ops.values() if v)), 'N_RUTAS': str(len(d['paths'])), 'N_OPS': str(len(OPS)), 'N_ESQ': str(len(d['components']['schemas'])), 'N_IDEM': str(sum(o['idem'] for o in OPS)),
        'N_LISTAS': str(sum(1 for o in OPS if any(p.get('$ref', '').endswith('/Cursor') for p in o['params']))), 'FILTROS': FILTROS, 'ENTIDADES': ENTIDADES, 'VISTAS': VISTAS, 'MATRIZ': MATRIZ,
        'RESUMEN': RESUMEN, 'SIN_N': SIN_N, 'N_SIN_EP': N_SIN_EP, 'TABLAS': TABLAS, 'AJUSTES': AJUSTES, 'CAPACIDADES': CAPACIDADES, 'OPERACIONES': OPERACIONES, 'INCONSISTENCIAS': INCONS, 'PARAMS': PARAMS_T, 'PARAMS_SV': PARAMS_SV, 'PARAMS_K43': PARAMS_K43}
for k, v in repl.items(): t = t.replace('{{%s}}' % k, v)
assert '{{' not in t, re.findall(r'{{\w+}}', t)
open(os.path.join(OUT, 'docs/contratos/C1-datos.md'), 'w', encoding='utf-8').write(t)
print('C1-datos.md', len(t), 'bytes; N sin endpoint:', sin_ep, '; ops sin N:', [o['oid'] for o in sin_n], '; N desconocidas:', unknown_n, '; tablas sin op:', sin_tabla, '; tablas inexistentes:', bad_tab, '; capacidades con id malo:', cap_bad, '; pantallas:', len(SCREENS))
