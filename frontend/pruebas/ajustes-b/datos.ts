// API falsa de Ajustes 8 a 14 (datos de ejemplo de prueba, sin credenciales). Se registra después de ApiFalsa.
import type { BrowserContext } from '@playwright/test'
import type { Cuerpo, Esquemas, Respuesta, RutasCon } from '../../src/api/tipos'
import { erroresC1 } from '../ajustes-a/validar-c1'
/** Ruta del contrato con cada {parámetro} cambiado por un número, p. ej. /api/v1/credits/1 (para tipar las respuestas falsas contra C1, H-11). */
type Concreta<S extends string> = S extends `${infer A}{${string}}${infer B}` ? `${A}${number}${Concreta<B>}` : S
type RespuestasGet = { [P in RutasCon<'get'> & string as `GET ${Concreta<P>}`]?: Respuesta<P, 'get'> }
const t = '2026-10-01T09:00:00-06:00'
const sello = { created_at: t, updated_at: t }
const p = <T extends string, V>(key: string, value_type: T, unit: string, value: V, initial_value: V = value) => ({ key, value_type, unit, value, initial_value, is_initial: JSON.stringify(value) === JSON.stringify(initial_value) })
const nombres = [['daily_summary', 'Resumen diario'], ['weekly_summary', 'Resumen semanal'], ['monthly_report', 'Reporte mensual en PDF'], ['payment_due', 'Pago próximo'], ['budget_threshold', 'Presupuesto al 80 % y al 100 %'], ['stale_prices', 'Precios viejos antes del día de despensa'], ['meal_plan_missing', 'Plan de comidas sin armar']] as const
export const PERSONAS = ([['axel', 'Axel', 'holder', true], ['ximena', 'Ximena', 'co_holder', true], ['karen', 'Karen', 'external', false]] as const).map(([code, display_name, participation, has_account], i) => ({ id: i + 1, code, display_name, participation, has_account, pays_share_directly: code === 'karen', is_active: true, ...sello }))
export const RESPUESTAS: Record<string, unknown> = {
  'GET /api/v1/persons': { items: PERSONAS, next_cursor: null },
  'GET /api/v1/notifications': { items: nombres.map(([code, name], i) => ({ id: i + 1, code, name, is_enabled: i !== 1, params_json: {}, recipient_ids: i === 0 ? [10] : [], job_id: i + 1, ...sello })), next_cursor: null },
  'GET /api/v1/jobs': { items: [
    { id: 1, code: 'daily_summary', name: 'Resumen diario', schedule: '30 21 * * *', is_enabled: true, notification_id: 1, last_status: 'ok', last_run_at: '2026-10-11T21:30:00-06:00', next_run_at: '2026-10-12T21:30:00-06:00', ...sello },
    { id: 8, code: 'backup_daily', name: 'Respaldo diario', schedule: '0 3 * * *', is_enabled: true, next_run_at: '2026-10-13T03:00:00-06:00', ...sello },
  ], next_cursor: null },
  'GET /api/v1/settings/reports': { parameters: { 'reports.daily_send_time': p('reports.daily_send_time', 'time', '', '21:30'), 'reports.weekly_weekday': p('reports.weekly_weekday', 'integer', '', 7), 'reports.weekly_send_time': p('reports.weekly_send_time', 'time', '', '19:00'), 'reports.monthly_day_of_month': p('reports.monthly_day_of_month', 'integer', '', 31, 1), 'reports.monthly_send_time': p('reports.monthly_send_time', 'time', '', '08:00') } },
  'GET /api/v1/notification-recipients': { items: [{ id: 10, person_id: 1, channel: 'telegram', label: 'Teléfono de Axel', verified_at: '2026-10-02T10:00:00-06:00', is_active: true, is_registered: true, ...sello }], next_cursor: null },
  'GET /api/v1/telegram-chat-candidates': { items: [{ id: 5, telegram_name: 'Xime', detected_at: '2026-10-11T20:15:00-06:00', chat_id_masked: '•••4821', ...sello }], next_cursor: null },
  'GET /api/v1/monthly-report/months': { months: [{ month: '2026-09', is_partial: false }, { month: '2026-10', is_partial: true }] },
  'GET /api/v1/settings/alerts': { parameters: { 'alerts.budget_warning_pct': p('alerts.budget_warning_pct', 'integer', '%', 80), 'alerts.budget_limit_pct': p('alerts.budget_limit_pct', 'integer', '%', 100), 'alerts.payment_due_days_before': p('alerts.payment_due_days_before', 'integer_list', 'días', [1, 0]), 'alerts.send_time': p('alerts.send_time', 'time', '', '09:00'), 'alerts.stale_prices_advance_days': p('alerts.stale_prices_advance_days', 'integer', 'días', 1), 'alerts.meal_plan_missing_weekday': p('alerts.meal_plan_missing_weekday', 'integer', '', 7) } },
  'GET /api/v1/budget/categories': { month: '2026-10', rows: [{ category_id: 1, planned_cents: 800000, actual_cents: 600000, difference_cents: 200000, used_pct: 75, status: 'normal', expense_count: 5 }, { category_id: 2, planned_cents: 0, actual_cents: 0, expense_count: 0 }] },
  'GET /api/v1/budget/month': { month: '2026-10', income_cents: 4000000, planned_spending_cents: 2500000, planned_savings_cents: 300000, unassigned_cents: 0, actual_cents: 1000000, remaining_cents: 1500000, days_left: 19 },
  'GET /api/v1/dashboard/next-payment': { next_payment: { due_date: '2026-10-30', credit_id: 1, amount_cents: 120000, household_share_cents: 60000, external_share_cents: 60000, days_left: 18, payment_count: 1, is_overdue: false }, availability: [{ credit_id: 1, availability_reading: 'B', available_cents: 610500, availability_confirmed: false }] },
  'GET /api/v1/backups': { items: [
    { id: 2, kind: 'scheduled', status: 'ok', started_at: '2026-10-12T03:00:00-06:00', finished_at: '2026-10-12T03:00:09-06:00', file_path: '/respaldos/fh-20261012.db', size_bytes: 5242880, sha256: 'a1b2c3d4e5f6a7b8c9d0', integrity_check: 'ok', restore_tested_at: '2026-10-12T03:01:00-06:00', restore_test_result: 'ok', ...sello },
    { id: 1, kind: 'manual', status: 'failed', started_at: '2026-10-11T12:00:00-06:00', error: 'Sin espacio en el disco de respaldos', ...sello },
  ], next_cursor: null },
  'GET /api/v1/system/seed-status': { items: [{ code: 'purchase_2026_10_10', is_loaded: true, row_count: 31 }, { code: 'kueski_credit', is_loaded: true, row_count: 12 }, { code: 'fixed_items_csv', is_loaded: false, row_count: 0 }] },
  'GET /api/v1/ai/providers': { providers: [{ provider: 'local_model', is_enabled: false, allowed_uses: ['ticket_ocr', 'recipe_suggestions'], is_available: false, is_configured: true }, { provider: 'claude_api', is_enabled: false, allowed_uses: ['ticket_ocr', 'recipe_suggestions'], is_available: false, is_configured: false }] },
  'GET /api/v1/system/status': { version: '0.1.0', started_at: '2026-10-12T08:00:00-06:00', database_size_bytes: 3145728 },
  'GET /api/v1/settings/data': { parameters: { 'attachment.max_size_mb': p('attachment.max_size_mb', 'integer', 'MB', 10), 'idempotency.retention_hours': p('idempotency.retention_hours', 'integer', 'h', 24), 'ocr.timeout_seconds': p('ocr.timeout_seconds', 'integer', 's', 45, 30) } },
  'GET /api/v1/cost-share-rules': { items: [
    { id: 1, name: 'Hogar', remainder_person_id: 1, is_template: true, is_default: true, is_active: true, parts: [], ...sello },
    { id: 2, name: 'Karen', remainder_person_id: 1, is_template: true, is_default: false, is_active: true, parts: [{ id: 1, rule_id: 2, person_id: 3, pct_bp: 10000, ...sello }], ...sello },
    { id: 3, name: 'Conjunto 50/50', remainder_person_id: 1, is_template: true, is_default: false, is_active: true, parts: [{ id: 2, rule_id: 3, person_id: 3, pct_bp: 5000, ...sello }], ...sello },
  ], next_cursor: null },
  'GET /api/v1/health-profiles': { items: [{ id: 1, person_id: 1, height_mm: 1750, birth_date: '1991-05-01', sex: 'male', activity_level: 'moderate', goal: 'maintain', current_weight_g: 72500, current_weight_date: '2026-10-10', age_years: 35, bmi_tenths: 237, bmi_class: 'normal', suggested_kcal_target: 2556, household_goal_adjust_lose_kcal: 0, household_goal_adjust_gain_kcal: 0, ...sello }], next_cursor: null },
  'GET /api/v1/settings/health': { parameters: { 'health.goal_adjust_lose_kcal': p('health.goal_adjust_lose_kcal', 'integer', 'kcal', 0), 'health.goal_adjust_gain_kcal': p('health.goal_adjust_gain_kcal', 'integer', 'kcal', 0) } },
  'GET /api/v1/savings-goals': { items: [
    { id: 1, name: 'Vacaciones', kind: 'goal', target_cents: 629344, target_date: '2026-12-31', planned_monthly_cents: 100000, start_date: '2026-08-01', is_active: true, saved_cents: 314672, progress_tenths_pct: 500, required_monthly_cents: 104891, status: 'in_progress', ...sello },
    { id: 2, name: 'Fondo de ahorro', kind: 'fund', planned_monthly_cents: 50000, start_date: '2026-01-01', is_active: true, saved_cents: 800000, status: 'in_progress', ...sello },
  ], next_cursor: null },
  'GET /api/v1/inventory-items': { items: [
    { id: 1, product_id: 1, tracking_mode: 'consumption_rate', quantity_milli: 2838000, counted_date: '2026-09-12', rate_quantity_milli: 946000, rate_period_days: 21, reorder_level_milli: 946000, status: 'sufficient', estimated_quantity_milli: 1486571, weekly_consumption_milli: 315333, days_left: 32, ...sello },
    { id: 2, product_id: 2, tracking_mode: 'stock', quantity_milli: 0, status: 'out', ...sello },
    { id: 3, product_id: 3, tracking_mode: 'stock', quantity_milli: 1500000, reorder_level_milli: 500000, status: 'sufficient', ...sello },
  ], next_cursor: null },
  'GET /api/v1/products': { items: [{ id: 1, name: 'Aceite de cocina', base_unit: 'ml', is_active: true, ...sello }, { id: 2, name: 'Arroz', base_unit: 'g', is_active: true, ...sello }, { id: 3, name: 'Frijol', base_unit: 'g', is_active: true, ...sello }], next_cursor: null },
} satisfies RespuestasGet
export interface Escritura { metodo: string; ruta: string; cuerpo: unknown }
/** Errores de C1 del tráfico falso de la prueba en curso: cuerpos de petición y respuestas con su código (R2-01). Cada prueba termina con cero. */
export const erroresTrafico: string[] = []
type ConParametros = { parameters: Record<string, { value: unknown; initial_value: unknown; is_initial: boolean }> }
/** La respuesta de un PATCH de parámetros: los de la lectura con los valores nuevos del cuerpo. */
function conCambios<R extends ConParametros>(lectura: R, cuerpo: unknown): R {
  const cambios = (cuerpo as { parameters?: Record<string, unknown> }).parameters ?? {}
  return { ...lectura, parameters: Object.fromEntries(Object.entries(lectura.parameters).map(([k, v]) => [k, k in cambios ? { ...v, value: cambios[k], is_initial: JSON.stringify(cambios[k]) === JSON.stringify(v.initial_value) } : v])) }
}
/** Respuesta falsa de cada escritura que ejercen las pruebas, con el código y el esquema de C1 (R2-01). undefined = sin respuesta falsa. */
function respuestaEscritura(metodo: string, ruta: string, cuerpo: unknown): [number, unknown?] | undefined {
  const leer = <P extends RutasCon<'get'>>(r: P) => RESPUESTAS[`GET ${r}`] as Respuesta<P, 'get'>
  switch (`${metodo} ${ruta.replace(/\/\d+(?=\/|$)/g, '/{id}')}`) {
    case 'PATCH /api/v1/settings/data': return [200, conCambios(leer('/api/v1/settings/data'), cuerpo) satisfies Respuesta<'/api/v1/settings/data', 'patch'>]
    case 'PATCH /api/v1/settings/reports': return [200, conCambios(leer('/api/v1/settings/reports'), cuerpo) satisfies Respuesta<'/api/v1/settings/reports', 'patch'>]
    case 'PATCH /api/v1/settings/alerts': return [200, conCambios(leer('/api/v1/settings/alerts'), cuerpo) satisfies Respuesta<'/api/v1/settings/alerts', 'patch'>]
    case 'PATCH /api/v1/inventory-items/{id}': {
      const actual = (leer('/api/v1/inventory-items').items as Esquemas['InventoryItem'][]).find((x) => x.id === Number(ruta.split('/')[4]))!
      const c = cuerpo as Cuerpo<'/api/v1/inventory-items/{inventory_item_id}', 'patch'>
      return [200, { ...actual, tracking_mode: c.tracking_mode ?? actual.tracking_mode, ...(typeof c.reorder_level_milli === 'number' ? { reorder_level_milli: c.reorder_level_milli } : {}) } satisfies Respuesta<'/api/v1/inventory-items/{inventory_item_id}', 'patch'>]
    }
    case 'PATCH /api/v1/notifications/{id}': {
      const actual = (leer('/api/v1/notifications').items as Esquemas['Notification'][]).find((x) => x.id === Number(ruta.split('/')[4]))!
      const c = cuerpo as Cuerpo<'/api/v1/notifications/{notification_id}', 'patch'>
      return [200, { ...actual, is_enabled: c.is_enabled ?? actual.is_enabled } satisfies Respuesta<'/api/v1/notifications/{notification_id}', 'patch'>]
    }
    case 'PATCH /api/v1/savings-goals/{id}': {
      const actual = (leer('/api/v1/savings-goals').items as Esquemas['SavingsGoal'][]).find((x) => x.id === Number(ruta.split('/')[4]))!
      const c = cuerpo as Cuerpo<'/api/v1/savings-goals/{savings_goal_id}', 'patch'>
      return [200, { ...actual, is_active: c.is_active ?? actual.is_active } satisfies Respuesta<'/api/v1/savings-goals/{savings_goal_id}', 'patch'>]
    }
    case 'PATCH /api/v1/ai/providers/local_model': {
      const actual = leer('/api/v1/ai/providers').providers.find((x) => x.provider === 'local_model')!
      return [200, { ...actual, ...(cuerpo as Cuerpo<'/api/v1/ai/providers/{provider}', 'patch'>) } satisfies Respuesta<'/api/v1/ai/providers/{provider}', 'patch'>]
    }
  }
  return metodo === 'DELETE' ? [204] : undefined
}
/** Instala las respuestas del área. Las escrituras se guardan en `escrituras` y responden con el código y el esquema de C1; todo el tráfico falso se valida contra C1 en `erroresTrafico`. */
export async function instalarAjustesB(ctx: BrowserContext, escrituras: Escritura[] = []) {
  erroresTrafico.length = 0
  await ctx.route(/\/api\/v1\//, async (route) => {
    const req = route.request()
    const u = new URL(req.url())
    const clave = `${req.method()} ${u.pathname}`
    let cuerpo: unknown = null
    try { cuerpo = req.postDataJSON() } catch { cuerpo = req.postData() }
    const responder = (estado: number, json?: unknown) => {
      erroresTrafico.push(...erroresC1(req.method(), u.pathname, estado, cuerpo, json ?? null))
      return json === undefined ? route.fulfill({ status: estado }) : route.fulfill({ status: estado, json })
    }
    if (req.method() === 'GET') return RESPUESTAS[clave] !== undefined ? responder(200, RESPUESTAS[clave]) : route.fallback()
    escrituras.push({ metodo: req.method(), ruta: u.pathname, cuerpo })
    const r = respuestaEscritura(req.method(), u.pathname, cuerpo)
    if (!r) { erroresTrafico.push(`${clave}: la prueba no tiene respuesta falsa de C1 para esta escritura`); return route.fulfill({ status: 501 }) }
    return responder(r[0], r[1])
  })
}
