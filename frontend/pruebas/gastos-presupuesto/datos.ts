// API del área Gastos y Presupuesto para Playwright (context.route sobre ApiFalsa, convenciones §8). Datos de ejemplo
// con la forma de C1 v1.4 (centavos enteros, fechas ISO); no son la semilla real. Lo que no atiende → route.fallback().
import type { BrowserContext, Page } from '@playwright/test'
import { expect } from '@playwright/test'
import type { ApiFalsa } from '../apoyo/api-falsa'
const t = '2026-10-01T09:00:00-06:00'
const pagina = (items: unknown[]) => ({ items, next_cursor: null })
export const GASTOS = [
  { id: 11, kind: 'variable', description: 'Súper semanal', category_id: 1, amount_cents: 185050, expense_date: '2026-10-05', accrual_months: 1, payment_method: 'cash', person_id: 1, cost_share_rule_id: null, household_share_cents: 150000, external_share_cents: 35050, has_attachment: true, shares: [{ person_id: 1, amount_cents: 150000 }, { person_id: 3, amount_cents: 35050 }], accruals: [{ month: '2026-10', amount_cents: 185050 }], cash_flows: [{ date: '2026-10-05', amount_cents: 185050 }], created_at: t, updated_at: t },
  { id: 12, kind: 'upcoming', description: 'Predial', category_id: 6, amount_cents: 120000, expense_date: '2026-10-20', accrual_months: 1, payment_method: 'transfer', person_id: 1, household_share_cents: 120000, external_share_cents: 0, has_attachment: false, created_at: t, updated_at: t },
]
const LECTURA = { id: 501, kind: 'ticket_photo', mime_type: 'image/jpeg', size_bytes: 1000, sha256: 'a'.repeat(64), ocr_status: 'done', created_at: t, updated_at: t, has_draft: false, draft: null, display_store_id: null, display_total_cents: 25000, proposed_store_id: null, proposed_expense_date: '2026-10-11', proposed_total_cents: 25000,
  lines: [{ line_no: 1, text: 'LECHE 1L', product_id: null, presentation_id: null, quantity_milli: 2000, unit_price_cents: 2900, line_total_cents: 5800, cost_share_rule_id: null, confidence: 'high' }, { line_no: 2, text: 'HUEVO 30', product_id: null, presentation_id: null, quantity_milli: 1000, unit_price_cents: 19200, line_total_cents: 19200, cost_share_rule_id: null, confidence: 'low' }] }
export class AreaFalsa {
  ingresos: { clave: string; cuerpo: Record<string, unknown> }[] = []
  movimientos: { clave: string; cuerpo: Record<string, unknown> }[] = []
  presupuestos: { url: string; cuerpo: unknown }[] = []
  proyeccion = false
  confirmados: Record<string, unknown>[] = []
  constructor(private api: ApiFalsa) {}
  async instalar(ctx: BrowserContext) {
    await ctx.route(/\/api\/v1\//, async (route) => {
      const pet = route.request(); const u = new URL(pet.url()); const ruta = `${pet.method()} ${u.pathname}`
      const json = (cuerpo: unknown, estado = 200) => route.fulfill({ status: estado, contentType: 'application/json', body: JSON.stringify(cuerpo) })
      const clave = pet.headers()['idempotency-key'] ?? ''
      const capturados = this.api.gastos.map((g, i) => ({ id: 100 + i, kind: 'variable', accrual_months: 1, household_share_cents: g.cuerpo.amount_cents, external_share_cents: 0, created_at: t, updated_at: t, ...g.cuerpo }))
      switch (ruta) {
        case 'GET /api/v1/expenses': return json(pagina([...capturados, ...GASTOS]))
        case 'GET /api/v1/expenses/totals': return json({ expense_count: 2 + capturados.length, amount_cents: 305050, household_share_cents: 270000, external_share_cents: 35050 })
        case 'GET /api/v1/expenses/summary': return json({ period_start: '2026-10-01', period_end: '2026-10-31', household_share_cents: 270000, previous_household_share_cents: 250000, difference_cents: 20000, daily_average_cents: 22500, weekly_average_cents: 157500, monthly_average_cents: 675000, top_category_id: 1, top_category_household_share_cents: 150000 })
        case 'GET /api/v1/expenses/monthly-trend': return json({ months: [{ month: '2026-09', categories: [{ category_id: 1, household_share_cents: 400000 }] }, { month: '2026-10', categories: [{ category_id: 1, household_share_cents: 150000 }, { category_id: 6, household_share_cents: 120000 }] }] })
        case 'GET /api/v1/expenses/11': return json(GASTOS[0])
        case 'GET /api/v1/expenses/999': return json({ code: 'not_found', message: 'No existe.', request_id: 'prueba' }, 404)
        case 'GET /api/v1/expenses/11/items': return json(pagina([{ id: 1, expense_id: 11, line_no: 1, description: 'Leche 1 L', quantity_milli: 2000, unit_price_cents: 2900, line_total_cents: 5800, created_at: t, updated_at: t, household_share_cents: 5800 }]))
        case 'GET /api/v1/ticket-scans': return json(pagina([{ ...LECTURA, id: 502, has_draft: true }]))
        case 'POST /api/v1/ticket-scans': return clave.length >= 16 ? json({ ...LECTURA, ocr_status: 'pending' }, 202) : json({ code: 'bad_request', message: 'Falta la clave.', request_id: 'prueba' }, 400)
        case 'GET /api/v1/ticket-scans/501': return json(LECTURA)
        case 'GET /api/v1/attachments/501/file': return route.fulfill({ status: 200, contentType: 'image/png', body: Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+ip1sAAAAASUVORK5CYII=', 'base64') })
        case 'POST /api/v1/ticket-scans/501/confirm': this.confirmados.push(pet.postDataJSON()); return json({ expense_id: 11, price_ids: [1, 2], inventory_item_ids: [] }, 201)
        case 'GET /api/v1/budget/month': return json({ month: '2026-10', income_cents: 3000000, planned_spending_cents: 2200000, planned_savings_cents: 300000, unassigned_cents: 500000, actual_cents: 1850000, remaining_cents: 350000, used_pct: 84, status: 'warning', days_left: 19 })
        case 'GET /api/v1/budget/categories': return json({ month: '2026-10', rows: [{ category_id: 1, planned_cents: 800000, actual_cents: 820000, difference_cents: -20000, used_pct: 103, status: 'exceeded', expense_count: 6 }, { category_id: 3, planned_cents: 200000, actual_cents: 50000, difference_cents: 150000, used_pct: 25, status: 'normal', expense_count: 2 }] })
        case 'GET /api/v1/budget/fortnights': return json({ month: '2026-10', fortnights: [{ half: 1, start_date: '2026-10-01', end_date: '2026-10-15', pay_date: '2026-10-15', actual_cents: 900000, planned_cents: 1100000 }, { half: 2, start_date: '2026-10-16', end_date: '2026-10-31', pay_date: '2026-10-30', actual_cents: 950000, planned_cents: 1100000 }] })
        case 'GET /api/v1/budget/alerts': return json({ rows: [{ category_id: 1, used_pct: 103, threshold_pct: 100, crossed_date: '2026-10-09' }] })
        case 'GET /api/v1/budget/accrued-groceries': return json({ month: '2026-10', rows: [{ expense_id: 11, expense_date: '2026-09-12', total_cents: 600000, accrued_cents: 300000, household_share_cents: 250000, external_share_cents: 50000 }] })
        case 'GET /api/v1/budget/monthly-trend': return json({ months: [{ month: '2026-09', planned_cents: 2200000, actual_cents: 2100000 }, { month: '2026-10', planned_cents: 2200000, actual_cents: 1850000 }] })
        case 'GET /api/v1/budgets': return json(pagina([{ id: 7, category_id: 1, amount_cents: 800000, valid_from_date: '2026-01-01', created_at: t, updated_at: t }]))
        case 'PATCH /api/v1/budgets/7': this.presupuestos.push({ url: u.pathname + u.search, cuerpo: pet.postDataJSON() }); return json({ id: 7, category_id: 1, amount_cents: 900000, valid_from_date: '2026-10-01', created_at: t, updated_at: t })
        case 'GET /api/v1/settings/projection': return json({ repeat_cadence_enabled: this.proyeccion })
        case 'PATCH /api/v1/settings/projection': this.proyeccion = pet.postDataJSON().repeat_cadence_enabled; return json({ repeat_cadence_enabled: this.proyeccion })
        case 'GET /api/v1/cash-flow/fortnights': return json({ projection_mode: this.proyeccion ? 'with_repeat_cadence' : 'contracted_only', rows: [
          { start_date: '2026-10-01', end_date: '2026-10-15', pay_date: '2026-10-15', income_cents: 1500000, income_is_expected: false, credit_payments_cents: 300000, fixed_and_upcoming_cents: 400000, variable_estimate_cents: 500000, planned_savings_cents: 150000, balance_cents: 150000, cumulative_balance_cents: 150000, is_actual: true },
          { start_date: '2026-10-16', end_date: '2026-10-31', pay_date: '2026-10-30', income_cents: 1500000, income_is_expected: true, credit_payments_cents: 300000, fixed_and_upcoming_cents: 400000, variable_estimate_cents: 500000, planned_savings_cents: 150000, balance_cents: 150000, cumulative_balance_cents: 300000, is_actual: false }] })
        case 'GET /api/v1/cash-flow/accrual-vs-cash': return json({ projection_mode: 'contracted_only', months: [{ month: '2026-10', accrued_cents: 1800000, cash_cents: 1700000, accrued_household_cents: 1500000, accrued_external_cents: 300000, cash_household_cents: 1400000, cash_external_cents: 300000, cumulative_accrued_cents: 1800000, cumulative_cash_cents: 1700000, cumulative_difference_cents: 100000, is_canonical: true }] })
        case 'GET /api/v1/cash-flow/assumptions': return json({ projection_mode: 'contracted_only', planned_items: [], variable_estimate_method: 'Promedio de 3 meses', variable_estimate_cents: 1000000 })
        case 'GET /api/v1/debt/live': return json({ date: u.searchParams.get('date'), total_cents: 629344, credits: [], purchases: [] })
        case 'GET /api/v1/incomes/month': return json({ month: '2026-10', fortnights: [{ pay_date: '2026-10-15', status: 'captured', expected_cents: 1500000, income_id: 31, amount_cents: 1500000 }, { pay_date: '2026-10-30', status: 'expected', expected_cents: 1500000, income_id: 0, amount_cents: 0 }], extras: [{ id: 40, person_id: 1, kind: 'extra', description: 'Venta de bicicleta', amount_cents: 250000, income_date: '2026-10-08', created_at: t, updated_at: t }], captured_cents: 1500000, expected_pending_cents: 1500000, extras_cents: 250000, total_cents: 3250000 })
        case 'GET /api/v1/incomes/monthly-trend': return json({ months: [{ month: '2026-10', fixed_cents: 3000000, extra_cents: 250000, total_cents: 3250000 }] })
        case 'POST /api/v1/incomes': this.ingresos.push({ clave, cuerpo: pet.postDataJSON() }); return json({ id: 50, person_id: 1, ...pet.postDataJSON(), created_at: t, updated_at: t }, 201)
        case 'GET /api/v1/savings/summary': return json({ saved_total_cents: 4500000, goals: [], contributed_month_cents: 150000, planned_month_cents: 300000 })
        case 'GET /api/v1/savings-goals': return json(pagina([{ id: 1, name: 'Fondo de emergencia', kind: 'fund', start_date: '2026-01-01', is_active: true, created_at: t, updated_at: t, saved_cents: 3000000 }, { id: 2, name: 'Vacaciones', kind: 'goal', target_cents: 3000000, target_date: '2027-04-01', start_date: '2026-01-01', is_active: true, created_at: t, updated_at: t, saved_cents: 1500000, progress_pct: 50, required_monthly_cents: 250000, status: 'in_progress' }]))
        case 'GET /api/v1/savings-movements': return json(pagina([{ id: 9, savings_goal_id: 2, movement_date: '2026-10-02', kind: 'deposit', amount_cents: 150000, created_at: t, updated_at: t }]))
        case 'POST /api/v1/savings-movements': this.movimientos.push({ clave, cuerpo: pet.postDataJSON() }); return json({ id: 10, ...pet.postDataJSON(), created_at: t, updated_at: t }, 201)
        case 'GET /api/v1/savings/monthly-trend': return json({ months: [{ month: '2026-10', saved_end_cents: 4500000, contributed_cents: 150000, planned_cents: 300000 }] })
        case 'GET /api/v1/calendar/events': return json(pagina([{ event_type: 'income', date: '2026-10-15', title: 'Quincena', amount_cents: 1500000, ref_table: 'income', ref_id: 31 }]))
      }
      return route.fallback()
    })
  }
}
/** Móvil primero: sin desplazamiento horizontal de la página (la tabla se desplaza dentro de su caja). */
export async function sinDesborde(page: Page) {
  const [ancho, vista] = await page.evaluate(() => [document.documentElement.scrollWidth, window.innerWidth])
  expect(ancho).toBeLessThanOrEqual(vista)
}
