// API falsa del área Créditos, Salud y Reportes: respuestas de C1 v1.4 con la semilla canónica de U1 §3.21
// (despensa 1: capital $3,094.50, intereses $800.50, 4 quincenas, pagos de $973.75; hogar $598.82 y Karen $374.93).
// Se registra después de ApiFalsa (base): lo que no atiende va a route.fallback(). Sin red externa.
import type { BrowserContext, Route } from '@playwright/test'
import type { Cuerpo, Esquemas, Metodo, Respuesta, paths } from '../../src/api/tipos'
/** Ruta del contrato con cada {parámetro} cambiado por un número, p. ej. /api/v1/credits/1 (para tipar las respuestas falsas contra C1, H-11). */
type Concreta<S extends string> = S extends `${infer A}{${string}}${infer B}` ? `${A}${number}${Concreta<B>}` : S
/** Cuerpo de la respuesta 200/201/202 de «MÉTODO /ruta concreta» según C1 (esquema.gen.ts). */
type CuerpoDe<K extends string> = { [P in keyof paths & string]: { [M in Metodo]: K extends `${Uppercase<M>} ${Concreta<P>}` ? Respuesta<P, M> : never }[Metodo] }[keyof paths & string]
const t = '2026-10-01T09:00:00-06:00'
const base = { created_at: t, updated_at: t }
const pagina = <T,>(items: T[]) => ({ items, next_cursor: null })
const FECHAS = ['2026-10-30', '2026-11-15', '2026-11-30', '2026-12-15']
export const PAGOS: Esquemas['ScheduledPayment'][] = FECHAS.map((due_date, i) => ({ id: 101 + i, credit_purchase_id: 1, credit_id: 1, seq: i + 1, due_date, amount_cents: 97375, principal_cents: 77363 + (i === 3 ? 1 : 0), interest_cents: 20012 - (i === 3 ? 1 : 0), status: 'pending', balance_after_cents: 97375 * (3 - i), ...base,
  shares: [{ person_id: 1, amount_cents: 59882, paid_cents: 0, status: 'pending' }, { person_id: 3, amount_cents: 37493, paid_cents: 0, status: 'pending' }] }))
export const KUESKI: Respuesta<'/api/v1/credits/{credit_id}', 'get'> = { id: 1, name: 'Kueski Pay', kind: 'revolving_line', lender: 'Kueski', holder_person_id: 1, credit_limit_cents: 1200000, payment_frequency: 'fortnightly', availability_reading: 'B', availability_confirmed: false, is_active: true, live_debt_cents: 389500, available_cents: 810500, next_payment_date: '2026-10-30', next_payment_cents: 97375, live_purchase_count: 1, ...base }
const TARJETA: Respuesta<'/api/v1/credits/{credit_id}', 'get'> = { id: 2, name: 'Tarjeta Oro', kind: 'bank_card', lender: 'Banco', holder_person_id: 1, credit_limit_cents: 3000000, payment_frequency: 'monthly', availability_reading: 'A', availability_confirmed: true, is_active: true, statement_day: 5, payment_due_day: 25, annual_rate_pct: '45.5', live_debt_cents: 120000, available_cents: 2880000, live_purchase_count: 0, missing_terms: ['min_payment_pct_bp'], ...base }
export const COMPRA: Respuesta<'/api/v1/credit-purchases/{credit_purchase_id}', 'get'> = { id: 1, credit_id: 1, kind: 'installment_plan', description: 'Despensa Walmart', purchase_date: '2026-10-15', principal_cents: 309450, interest_cents: 80050, total_cents: 389500, term_count: 4, payment_frequency: 'fortnightly', first_payment_date: '2026-10-30', interest_source: 'quoted', interest_includes_vat: true, status: 'active', payments_made: 0, payments_total: 4, balance_cents: 389500, ...base,
  shares: [{ id: 1, credit_purchase_id: 1, person_id: 1, principal_cents: 190298, interest_cents: 49230, total_cents: 239528, fortnight_payment_cents: 59882, ...base }, { id: 2, credit_purchase_id: 1, person_id: 3, principal_cents: 119152, interest_cents: 30820, total_cents: 149972, fortnight_payment_cents: 37493, ...base }] }
const ESC = (credit_id: number, term_count: number, interest_cents: number, score_bp: number): Respuesta<'/api/v1/scenario-comparisons', 'post'>['valid'][number] => ({ credit_id, term_count, purchase_date: '2026-10-15', interest_cents, max_fortnight_payment_cents: 97375, fortnight_payment_spread_cents: 0, min_available_cents: 500000, max_live_purchases: 2, zero_debt_date: '2026-12-15', score_bp })
const REPORTES: Esquemas['GeneratedReport'][] = [
  { id: 11, kind: 'monthly_pdf', period_start_date: '2026-09-01', period_end_date: '2026-09-30', sha256: 'a'.repeat(64), size_bytes: 52000, generated_at: '2026-10-01T08:00:00-06:00', sent_at: '2026-10-01T08:01:00-06:00', delivery_status: 'sent', ...base },
  { id: 12, kind: 'external_statement_pdf', person_id: 3, period_start_date: '2026-09-01', period_end_date: '2026-09-30', sha256: 'b'.repeat(64), size_bytes: 21000, generated_at: '2026-10-02T12:20:00-06:00', delivery_status: 'not_sent', ...base },
  { id: 10, kind: 'monthly_pdf', period_start_date: '2026-08-01', period_end_date: '2026-08-31', sha256: 'c'.repeat(64), size_bytes: 50000, generated_at: '2026-09-01T08:00:00-06:00', delivery_status: 'failed', ...base },
]
const mensual = (month: string, conPdf: boolean): Respuesta<'/api/v1/monthly-report', 'get'> => ({ month,
  budget: { month, income_cents: 4200000, planned_spending_cents: 3100000, planned_savings_cents: 600000, unassigned_cents: 500000, actual_cents: 2950000, remaining_cents: 150000, used_pct: 95, status: 'warning', days_left: 0 },
  categories: { month, rows: [{ category_id: 1, planned_cents: 900000, actual_cents: 875050, difference_cents: 24950, used_pct: 97, status: 'warning', expense_count: 14 }] },
  fortnights: { month, fortnights: [{ half: 1, start_date: `${month}-01`, end_date: `${month}-15`, pay_date: `${month}-15`, actual_cents: 1500000, planned_cents: 1550000 }, { half: 2, start_date: `${month}-16`, end_date: `${month}-30`, pay_date: `${month}-30`, actual_cents: 1450000, planned_cents: 1550000 }] },
  credits: { live_debt_cents: 389500, fortnight_total_cents: 97375, fortnight_household_cents: 59882, fortnight_external_cents: 37493, month_total_cents: 194750, month_household_cents: 119764, month_external_cents: 74986, interest_paid_month_cents: 40024, interest_paid_year_cents: 40024 },
  live_debt_end_cents: 389500, meals: { weeks: [{ week_start: '2026-09-07', cost_cents: 185000, persons: [{ person_id: 1, avg_kcal_per_day: 2400, daily_kcal_target: 2364 }] }] },
  external_participation: [{ person_id: 3, month, accrued_groceries_cents: 149972, purchases: [], payment_shares: [{ scheduled_payment_id: 101, due_date: '2026-10-30', amount_cents: 37493, status: 'pending' }], live_debt_cents: 149972 }],
  ...(conPdf ? { pdf: REPORTES[0] } : {}) })

/** Eventos del calendario (C1 listCalendarEvents) de las dos líneas: pagos de Kueski con sus partes, la promoción Walmart y un corte de la tarjeta. */
const EVENTOS: Esquemas['CalendarEvent'][] = [
  { event_type: 'promotion', date: '2026-10-01', end_date: '2026-10-31', title: 'Promoción: Walmart a 4 pagos', credit_id: 1, ref_table: 'promotion', ref_id: 5 },
  ...PAGOS.map((p): Esquemas['CalendarEvent'] => ({ event_type: 'credit_payment', date: p.due_date, title: 'Pago de Kueski Pay', amount_cents: p.amount_cents, household_share_cents: 59882, external_share_cents: 37493, status: 'pending', credit_id: 1, ref_table: 'scheduled_payment', ref_id: p.id })),
  { event_type: 'card_statement', date: '2026-11-05', title: 'Corte de Tarjeta Oro', amount_cents: 120000, credit_id: 2, ref_table: 'credit_statement', ref_id: 7 },
]
/** Filtra EVENTOS como el servidor: period `desde..hasta`, credit_id y event_type separados por coma; orden por fecha. */
function eventosDe(q: URLSearchParams) {
  const [desde, hasta] = (q.get('period') ?? '').split('..')
  const tipos = (q.get('event_type') ?? '').split(',').filter(Boolean), lineas = (q.get('credit_id') ?? '').split(',').filter(Boolean).map(Number)
  return EVENTOS.filter((e) => (!tipos.length || tipos.includes(e.event_type)) && (!lineas.length || lineas.includes(e.credit_id ?? 0)) && (!desde || ((e.end_date ?? e.date) >= desde && e.date <= hasta)))
    .sort((a, b) => a.date.localeCompare(b.date))
}

export interface Registro { metodo: string; ruta: string; cuerpo: unknown; clave?: string }
/** Instala las rutas del área; devuelve las escrituras recibidas para revisarlas en la prueba. */
export async function instalarArea(ctx: BrowserContext) {
  const escrituras: Registro[] = []
  await ctx.route(/\/api\/v1\//, async (route: Route) => {
    const pet = route.request(), u = new URL(pet.url()), r = `${pet.method()} ${u.pathname}`, q = u.searchParams
    const json = <K extends string>(_clave: K, cuerpo: CuerpoDe<K>, status = 200) => route.fulfill({ status, contentType: 'application/json', headers: { 'X-Request-Id': 'prueba' }, body: JSON.stringify(cuerpo) })
    if (pet.method() !== 'GET') escrituras.push({ metodo: pet.method(), ruta: u.pathname, cuerpo: pet.postDataJSON?.() ?? null, clave: pet.headers()['idempotency-key'] })
    switch (r) {
      case 'GET /api/v1/credits': return json(r, pagina([KUESKI, TARJETA]))
      case 'GET /api/v1/credits/1': return json(r, KUESKI)
      case 'GET /api/v1/calendar/events': return json(r, pagina(eventosDe(q)))
      case 'GET /api/v1/credits/2': return json(r, TARJETA)
      case 'GET /api/v1/credits/totals': return json(r, mensual('2026-10', false).credits)
      case 'GET /api/v1/credits/debt-trend': return json(r, { months: ['2026-10', '2026-11', '2026-12'].map((month, i) => ({ month, live_debt_cents: 389500 - i * 194750, credits: [{ credit_id: 1, live_debt_cents: 389500 - i * 194750 }], interest_paid_cents: 40024, interest_pending_cents: 40026 - i * 20013 })) })
      case 'GET /api/v1/credits/1/availability-preview': return json(r, { credit_id: 1, availability_reading: 'B', availability_confirmed: false, reading_a: { available_today_cents: 1119455, horizon_min_cents: 809450 }, reading_b: { available_today_cents: 810500, horizon_min_cents: 500000 } })
      case 'GET /api/v1/credits/1/availability-series': return json(r, { credit_id: 1, events: [{ event_date: '2026-10-12', event_type: 'start', description: 'Inicio', amount_cents: 0, principal_restored_cents: 0, interest_not_restored_cents: 0, available_a_cents: 1200000, available_b_cents: 1200000 }, { event_date: '2026-10-15', event_type: 'purchase', description: 'Despensa Walmart', amount_cents: 389500, principal_restored_cents: 0, interest_not_restored_cents: 0, available_a_cents: 809450, available_b_cents: 810500 }], min_a_cents: 809450, min_b_cents: 500000 })
      case 'GET /api/v1/credits/1/term-rates': return json(r, pagina([{ id: 1, credit_id: 1, term_count: 2, rate_num: 1727, rate_den: 10000, source: 'quote', rate_pct: '17.27', ...base }, { id: 2, credit_id: 1, term_count: 4, rate_num: 258685, rate_den: 1000000, source: 'quote', rate_pct: '25.8685', ...base }]))
      case 'GET /api/v1/credits/2/card-status': return json(r, { credit_id: 2, statement_date: '2026-10-05', due_date: '2026-10-25', statement_balance_cents: 120000, no_interest_payment_cents: 120000, current_balance_cents: 120000, installment_plans: [] })
      case 'GET /api/v1/scheduled-payments': return json(r, pagina(PAGOS.filter(() => !q.get('credit_id') || q.get('credit_id') === '1')))
      case 'POST /api/v1/scheduled-payments/101/payments': return json(r, { id: 900, credit_id: 1, scheduled_payment_id: 101, ...(pet.postDataJSON() as Cuerpo<'/api/v1/scheduled-payments/{scheduled_payment_id}/payments', 'post'>), ...base }, 201)
      case 'DELETE /api/v1/credit-payments/900': return route.fulfill({ status: 204 })
      case 'GET /api/v1/credit-purchases': return json(r, pagina([COMPRA]))
      case 'GET /api/v1/credit-purchases/1': return json(r, COMPRA)
      case 'POST /api/v1/credit-purchases/preview': return json(r, { principal_cents: 309450, interest_cents: 80050, total_cents: 389500, shares: [{ person_id: 1, principal_cents: 190298, interest_cents: 49230, total_cents: 239528, fortnight_payment_cents: 59882 }, { person_id: 3, principal_cents: 119152, interest_cents: 30820, total_cents: 149972, fortnight_payment_cents: 37493 }], schedule: PAGOS.map((p) => ({ seq: p.seq, due_date: p.due_date, amount_cents: p.amount_cents, principal_cents: p.principal_cents, interest_cents: p.interest_cents, shares: (p.shares ?? []).map((s) => ({ person_id: s.person_id, amount_cents: s.amount_cents })) })) })
      case 'GET /api/v1/promotions': return json(r, pagina([{ id: 5, credit_id: 1, merchant_name: 'Walmart', term_count: 4, rate_num: 0, rate_den: 1, min_purchase_cents: 150000, valid_from: '2026-10-01', valid_to: '2026-10-31', is_active: true, ...base }]))
      case 'GET /api/v1/settings/optimizer': return json(r, { missing_constraints: [], parameters: { 'comparator.timeout_seconds': { key: 'comparator.timeout_seconds', value_type: 'integer', unit: 'segundos', value: 5, initial_value: 5, is_initial: true }, 'optimizer.date_window_days': { key: 'optimizer.date_window_days', value_type: 'integer', unit: 'días', value: 7, initial_value: 7, is_initial: true } }, min_available_buffer_cents: 50000, max_live_purchases: 3, weights: { lowest_interest: 1, even_low_payment: 1, max_available: 1, most_value_per_peso: 1 } } satisfies Respuesta<'/api/v1/settings/optimizer', 'get'>)
      case 'POST /api/v1/scenario-comparisons': return json(r, { status: 'ok', missing_constraints: [], evaluated_count: 24, valid: [ESC(1, 4, 80050, 8734), ESC(1, 2, 53440, 7100)], optima: { lowest_interest: ESC(1, 2, 53440, 7100), even_low_payment: ESC(1, 4, 80050, 8734), max_available: ESC(1, 4, 80050, 8734), most_value_per_peso: ESC(1, 4, 80050, 8734) }, recommendation: ESC(1, 4, 80050, 8734), rejected: [{ scenario: ESC(2, 6, 120000, 0), violated_constraint: 'max_live_purchases' }] })
      case 'GET /api/v1/health-profiles': return json(r, pagina([
        { id: 1, person_id: 1, height_mm: 1780, sex: 'male', activity_level: 'moderate', daily_kcal_target: 2364, goal: 'lose', current_weight_g: 84600, current_weight_date: '2026-10-05', age_years: 36, bmi_tenths: 267, bmi_class: 'overweight', suggested_kcal_target: 2364, household_goal_adjust_lose_kcal: 0, household_goal_adjust_gain_kcal: 0, ...base },
        { id: 2, person_id: 2, height_mm: 1630, sex: 'female', activity_level: 'light', daily_kcal_target: 1900, goal: 'maintain', current_weight_g: 62800, current_weight_date: '2026-10-05', age_years: 33, bmi_tenths: 236, bmi_class: 'normal', suggested_kcal_target: 1816, household_goal_adjust_lose_kcal: 0, household_goal_adjust_gain_kcal: 0, ...base }]))
      case 'GET /api/v1/health/today-kcal': return json(r, { date: '2026-10-12', persons: [{ person_id: 1, planned_kcal: 2413, daily_kcal_target: 2364, difference_kcal: 49, eating_out_slots: 0 }, { person_id: 2, planned_kcal: 2325, daily_kcal_target: 1900, difference_kcal: 425, eating_out_slots: 1 }] })
      case 'GET /api/v1/weight-entries': return json(r, pagina([{ id: 1, person_id: 1, measured_date: '2026-10-05', weight_g: 84600, bmi_tenths: 267, ...base }, { id: 2, person_id: 1, measured_date: '2026-09-28', weight_g: 84900, bmi_tenths: 268, ...base }, { id: 3, person_id: 2, measured_date: '2026-10-05', weight_g: 62800, bmi_tenths: 236, ...base }]))
      case 'GET /api/v1/meal-plan/trends': return json(r, { weeks: [{ week_start: '2026-09-28', cost_cents: 185000, persons: [{ person_id: 1, avg_kcal_per_day: 2400, daily_kcal_target: 2364 }, { person_id: 2, avg_kcal_per_day: 1950, daily_kcal_target: 1900 }] }] })
      case 'GET /api/v1/reports': return json(r, pagina(REPORTES.filter((x) => !q.get('kind') || x.kind === q.get('kind'))))
      case 'GET /api/v1/jobs': return json(r, pagina([['daily_summary', 'Resumen diario', '2026-10-12T21:30:00-06:00'], ['weekly_summary', 'Resumen semanal', '2026-10-18T19:00:00-06:00'], ['monthly_report', 'Reporte mensual', '2026-11-01T08:00:00-06:00'], ['backup_daily', 'Respaldo', '2026-10-13T03:00:00-06:00']].map(([code, name, next_run_at], i) => ({ id: i + 1, code, name, schedule: '0 0 * * *', is_enabled: true, next_run_at, last_status: 'ok', ...base }))))
      case 'GET /api/v1/notification-deliveries': return json(r, pagina([{ id: 1, notification_id: 1, recipient_id: 1, dedupe_key: 'd-1', status: 'sent', notification_code: 'daily_summary', message_text: 'Resumen del 11-oct-2026: gastaste $951.30 (parte del hogar) en 2 gastos.', sent_at: '2026-10-11T21:30:00-06:00', ...base }]))
      case 'GET /api/v1/monthly-report/months': return json(r, { months: [{ month: '2026-10', is_partial: true, generated_report_id: 0 }, { month: '2026-09', is_partial: false, generated_report_id: 11 }] })
      case 'GET /api/v1/monthly-report': return json(r, mensual(q.get('period') ?? '2026-09', (q.get('period') ?? '2026-09') === '2026-09'))
      case 'GET /api/v1/categories': return route.fallback()
      case 'GET /api/v1/stores': return json(r, pagina([]))
      case 'POST /api/v1/reports/monthly': { const b = pet.postDataJSON() as { month: string }; return json(r, { ...REPORTES[0], id: 20, period_start_date: `${b.month}-01`, period_end_date: `${b.month}-31`, delivery_status: 'not_sent', generated_at: '2026-10-12T10:00:00-06:00' }, 202) }
      case 'POST /api/v1/reports/external-statements': { const b = pet.postDataJSON() as { period_start_date: string; period_end_date: string }; return json(r, { ...REPORTES[1], id: 21, ...b, generated_at: '2026-10-12T10:00:00-06:00' }, 202) }
      case 'POST /api/v1/reports/10/send': return json(r, { ...REPORTES[2], delivery_status: 'sent' }, 202)
      case 'GET /api/v1/persons/3/statement': return json(r, { person_id: 3, period_start_date: '2026-10-01', period_end_date: '2026-10-31', purchases: [], grocery_shares: [{ credit_purchase_id: 1, purchase_date: '2026-10-15', principal_cents: 119152, interest_cents: 30820, total_cents: 149972 }], payment_shares: [{ person_id: 3, amount_cents: 37493, paid_cents: 0, status: 'pending' }], paid_in_period_cents: 0, pending_total_cents: 149972 })
    }
    if (/^GET \/api\/v1\/reports\/\d+\/file$/.test(r)) return route.fulfill({ status: 200, contentType: 'application/pdf', body: '%PDF-1.4\n% prueba\n%%EOF\n' })
    return route.fallback()
  })
  return escrituras
}
