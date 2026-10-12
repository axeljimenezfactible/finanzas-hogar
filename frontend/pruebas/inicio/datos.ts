// Datos de prueba de P-01 con las cifras del oráculo y la semilla (U3: R-01, R-11, R-13 a R-15, R-45, R-46, R-61, R-62, R-64, M-01, M-08).
// Formas de C1 v1.4; ingresos y plan son marcadores (U3 los deja como marcador). Sin credenciales.
import type { BrowserContext } from '@playwright/test'
const t = '2026-10-01T09:00:00-06:00'
const pago = (id: number, due_date: string) => ({ id, credit_purchase_id: 1, seq: id, due_date, amount_cents: 97375, principal_cents: 90000, interest_cents: 7375, created_at: t, updated_at: t, credit_id: 1, status: 'pending', shares: [{ person_id: 1, amount_cents: 59882, paid_cents: 0, status: 'pending' }, { person_id: 3, amount_cents: 37493, paid_cents: 0, status: 'pending' }] })
const meses = Array.from({ length: 12 }, (_, i) => { const n = 9 + i; return `${2026 + Math.floor(n / 12)}-${String((n % 12) + 1).padStart(2, '0')}` })
export function respuestas(): Record<string, unknown> {
  return {
    '/api/v1/settings/alerts': { parameters: { 'alerts.payment_due_days_before': { key: 'alerts.payment_due_days_before', value_type: 'integer_list', unit: 'días', value: [1, 0] } } },
    '/api/v1/credits': { items: [{ id: 1, name: 'Kueski', kind: 'revolving_line', holder_person_id: 1, payment_frequency: 'fortnightly', availability_reading: 'B', availability_confirmed: false, is_active: true, created_at: t, updated_at: t }], next_cursor: null },
    '/api/v1/stores': { items: [['Walmart', 'chain'], ['Bodega Aurrera', 'chain'], ["Sam's Club", 'chain']].map(([name, kind], i) => ({ id: i + 1, name, kind, assisted_reading_enabled: false, requires_membership: i === 2, has_membership: false, is_active: true, created_at: t, updated_at: t })), next_cursor: null },
    '/api/v1/dashboard/month-vs-plan': { month: '2026-10', actual_cents: 119764, days_left: 22 },
    '/api/v1/dashboard/fortnight-flow': { start_date: '2026-10-01', end_date: '2026-10-15', income_cents: 1500000, income_is_expected: true, spent_cents: 119764, committed_cents: 0, closing_balance_cents: 1380236 },
    '/api/v1/dashboard/next-payment': { next_payment: { due_date: '2026-10-30', credit_id: 1, amount_cents: 97375, household_share_cents: 59882, external_share_cents: 37493, days_left: 20, payment_count: 1, is_overdue: false }, availability: [{ credit_id: 1, availability_reading: 'B', available_cents: 690550, availability_confirmed: false }] },
    '/api/v1/dashboard/week-meals': { week_start: '2026-10-05', week_end: '2026-10-11', cost_cents: 0, persons: [], unplanned_slots: 21 },
    '/api/v1/dashboard/best-basket': { basket: 'expense', expense_id: 10, store_id: 2, total_cents: 309450, difference_vs_second_cents: 13850, second_store_id: 1, covered_count: 33, item_count: 36, stale_price_count: 0, assumes_membership: false },
    '/api/v1/expenses/10': { id: 10, expense_date: '2026-10-10', amount_cents: 309450, category_id: 1, person_id: 1, created_at: t, updated_at: t },
    '/api/v1/dashboard/savings': { saved_total_cents: 0, contributed_month_cents: 0, planned_month_cents: 0 },
    '/api/v1/persons/3/participation': { person_id: 3, month: '2026-10', accrued_groceries_cents: 74986, purchases: [{ credit_purchase_id: 1, description: 'Despensa 1', purchase_date: '2026-10-10', total_cents: 389500, payments_made: 0, payments_total: 4, balance_cents: 149972 }], payment_shares: [{ scheduled_payment_id: 1, due_date: '2026-10-30', amount_cents: 37493, status: 'pending' }], live_debt_cents: 149972 },
    '/api/v1/scheduled-payments': { items: [pago(1, '2026-10-30'), pago(2, '2026-11-15'), pago(3, '2026-11-30')], next_cursor: null },
    '/api/v1/credits/1/availability-preview': { credit_id: 1, availability_reading: 'B', availability_confirmed: false, reading_a: { available_today_cents: 1000000, horizon_min_cents: 809450 }, reading_b: { available_today_cents: 690550, horizon_min_cents: 500000 } },
    '/api/v1/cash-flow/fortnights': { projection_mode: 'contracted_only', rows: [{ start_date: '2026-10-01', end_date: '2026-10-15', pay_date: '2026-10-15', income_cents: 1500000, income_is_expected: true, credit_payments_cents: 0, fixed_and_upcoming_cents: 0, variable_estimate_cents: 0, balance_cents: 1380236, cumulative_balance_cents: 1380236, is_actual: true }] },
    '/api/v1/cash-flow/accrual-vs-cash': { projection_mode: 'contracted_only', months: meses.map((month, i) => {
      const acc = i < 2 ? 194750 : 0, cash = i === 0 ? 97375 : i === 1 ? 194750 : 0
      const cumAcc = i === 0 ? 194750 : 389500, cumCash = i === 0 ? 97375 : 292125
      return { month, accrued_cents: acc, cash_cents: cash, accrued_household_cents: acc ? 119764 : 0, accrued_external_cents: acc ? 74986 : 0, cash_household_cents: cash ? 59882 * (i + 1) : 0, cash_external_cents: cash ? 37493 * (i + 1) : 0, cumulative_accrued_cents: cumAcc, cumulative_cash_cents: cumCash, cumulative_difference_cents: cumAcc - cumCash, is_canonical: i < 6 }
    }) },
    '/api/v1/budget/monthly-trend': { months: [] },
    '/api/v1/credits/debt-trend': { months: meses.map((month, i) => ({ month, live_debt_cents: Math.max(0, 389500 - i * 97375), credits: [{ credit_id: 1, live_debt_cents: Math.max(0, 389500 - i * 97375) }], interest_paid_cents: 0, interest_pending_cents: 0 })) },
    '/api/v1/meal-plan/trends': { weeks: [] },
    '/api/v1/calendar/events': { items: [
      { event_type: 'income', date: '2026-10-15', title: 'Quincena de Axel', amount_cents: 1500000, status: 'expected', ref_table: 'income', ref_id: 1 },
      { event_type: 'credit_payment', date: '2026-10-30', title: 'Pago', amount_cents: 97375, household_share_cents: 59882, external_share_cents: 37493, status: 'pending', credit_id: 1, ref_table: 'scheduled_payment', ref_id: 1 },
    ], next_cursor: null },
  }
}
/** Instala las respuestas de Inicio (después de ApiFalsa); `cambios` reemplaza o falla rutas: número = estado HTTP de error. */
export async function instalarInicio(ctx: BrowserContext, cambios: Record<string, unknown> = {}, registro: string[] = []) {
  const r = { ...respuestas(), ...cambios }
  await ctx.route(/\/api\/v1\//, async (route) => {
    const u = new URL(route.request().url())
    if (route.request().method() !== 'GET' || !(u.pathname in r)) return route.fallback()
    registro.push(u.pathname + u.search)
    const v = r[u.pathname]
    if (typeof v === 'number') return route.fulfill({ status: v, contentType: 'application/json', body: JSON.stringify({ code: 'internal', message: 'Falla de prueba', request_id: 'prueba' }) })
    return route.fulfill({ status: 200, contentType: 'application/json', headers: { 'X-Request-Id': 'prueba' }, body: JSON.stringify(v) })
  })
}
