// Manejadores de msw del tablero (solo «npm run dev:mocks»; nunca en el build): las operaciones nuevas de C1 v1.5 (K54) con
// datos SINTÉTICOS coherentes entre sí (no son del oráculo): ingreso de $20,000.00 al mes, pagos del hogar de la semilla
// (R-13: $598.82 en octubre; $1,197.64 y $1,482.57 en noviembre y diciembre) y los umbrales iniciales 30 % y 40 % (K54).
// El estado y el porcentaje van escritos como los mandaría el motor: la interfaz no los calcula.
import { http, HttpResponse, type HttpHandler } from 'msw'
const INGRESO = 2000000
const PAGOS: Record<string, [number, number, number, 'normal' | 'warning' | 'exceeded']> = {
  '2026-10': [59882, 97375, 3, 'normal'], '2026-11': [119764, 194750, 6, 'normal'], '2026-12': [148257, 223243, 7, 'normal'],
}
const parametro = (key: string, value: number, unit: string) => ({ key, value_type: 'integer', unit, value })
const ajustes = { parameters: {
  'debt.capacity_warning_pct': parametro('debt.capacity_warning_pct', 30, '%'), 'debt.capacity_limit_pct': parametro('debt.capacity_limit_pct', 40, '%'),
  'dashboard.payment_reminder_count': parametro('dashboard.payment_reminder_count', 5, 'pagos'), 'budget.signal_months': parametro('budget.signal_months', 3, 'meses'),
  'budget.signal_margin_pct': parametro('budget.signal_margin_pct', 10, '%'),
} }
const mesMas = (m: string, n: number) => { const t = Number(m.slice(0, 4)) * 12 + Number(m.slice(5, 7)) - 1 + n; return `${Math.floor(t / 12)}-${String((t % 12) + 1).padStart(2, '0')}` }
export const manejadores: HttpHandler[] = [
  http.get('/api/v1/settings/dashboard', () => HttpResponse.json(ajustes)),
  http.get('/api/v1/dashboard/debt-capacity', ({ request }) => {
    const month = new URL(request.url).searchParams.get('period') ?? '2026-10'
    const [h, tot, n, status] = PAGOS[month] ?? [0, 0, 0, 'normal']
    return HttpResponse.json({ month, income_cents: INGRESO, income_is_expected: false, extras_cents: 0, payments_cents: h, payments_total_cents: tot, payment_count: h ? 1 : 0, ratio_pct: n, ratio_total_pct: n + 2, warning_pct: 30, limit_pct: 40, scale_max_pct: 60, status, headroom_cents: 800000 - h, has_credits: true })
  }),
  http.get('/api/v1/dashboard/monthly-overview', ({ request }) => {
    const u = new URL(request.url), desde = u.searchParams.get('from') ?? '2025-11', hasta = u.searchParams.get('to') ?? '2027-03'
    const months = []
    for (let m = desde; m <= hasta && months.length < 24; m = mesMas(m, 1)) {
      const [h, tot, n, status] = PAGOS[m] ?? [0, 0, 0, 'normal'], proy = m > '2026-10'
      months.push({ month: m, is_projection: proy, income_total_cents: INGRESO, income_fortnights_cents: INGRESO, income_extras_cents: 0, income_is_expected: proy, ...(m === '2026-10' ? { accrued_cents: 119764, cash_cents: 59882 } : {}), credit_payments_cents: h, credit_payments_total_cents: tot, debt_ratio_pct: n, debt_status: status })
    }
    return HttpResponse.json({ warning_pct: 30, limit_pct: 40, months })
  }),
  http.get('/api/v1/budget/signals', ({ request }) => HttpResponse.json({ month: new URL(request.url).searchParams.get('period') ?? '2026-10', months_considered: 3, margin_pct: 10, rows: [] })),
]
