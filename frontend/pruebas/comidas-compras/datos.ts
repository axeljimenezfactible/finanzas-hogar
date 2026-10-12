// Datos de prueba de Comidas y Compras con la forma de C1 v1.4 (sin credenciales reales). P-14 usa las cifras canon de S5.
import type { BrowserContext } from '@playwright/test'
const T = '2026-10-01T09:00:00-06:00'
const pag = (items: unknown[]) => ({ items, next_cursor: null })
const tienda = (id: number, name: string, kind: string, extra = {}) => ({ id, kind, name, assisted_reading_enabled: false, requires_membership: false, has_membership: false, is_active: true, created_at: T, updated_at: T, ...extra })
export const TIENDAS = [tienda(1, 'Walmart', 'chain'), tienda(2, 'Bodega Aurrera', 'chain'), tienda(3, "Sam's Club", 'chain', { requires_membership: true }), tienda(4, 'Mercado Portales', 'market', { visit_weekday: 6 })]
const precio = (store_id: number, c: number, fresco = true) => ({ store_id, presentation_id: 1, price_cents: c, base_unit_price_cents: Math.round(c / 10), observed_date: '2026-10-10', source: 'ticket_manual', is_stale: !fresco })
export const PRODUCTOS = [
  { id: 1, name: 'Jitomate saladet', product_group_id: 1, base_unit: 'g', is_active: true, created_at: T, updated_at: T, best_price: precio(2, 2890), last_observed_date: '2026-10-10', freshness: 'fresh' },
  { id: 2, name: 'Leche entera', product_group_id: 2, base_unit: 'ml', is_active: true, created_at: T, updated_at: T, last_observed_date: '2026-08-01', freshness: 'stale' },
  { id: 3, name: 'Huevo blanco', product_group_id: 2, base_unit: 'piece', is_active: true, created_at: T, updated_at: T, freshness: 'no_price' },
]
const GRUPOS = [{ id: 1, name: 'Frutas y verduras', purchase_channel: 'market', sort_order: 1, is_active: true, created_at: T, updated_at: T }, { id: 2, name: 'Lácteos y huevo', purchase_channel: 'chain', sort_order: 2, is_active: true, created_at: T, updated_at: T }]
const TIEMPOS = [['breakfast', 'Desayuno', false], ['lunch', 'Comida', false], ['dinner', 'Cena', false], ['snack', 'Colación', true]].map(([code, name, is_optional], i) => ({ id: i + 1, code, name, sort_order: i + 1, is_optional, created_at: T, updated_at: T }))
const receta = (id: number, name: string, slots: number[], c: number, k: number, usada: number) => ({ id, name, base_servings: 2, prep_minutes: 15, cook_minutes: 25, origin: 'claude_seed', is_active: true, created_at: T, updated_at: T, cost_per_serving_cents: c, kcal_per_serving: k, times_used_8w: usada, meal_slot_ids: slots, instructions: 'Picar la papa.\n\nGuisar con la carne.' })
export const RECETAS = [receta(1, 'Picadillo de res con papa', [2, 3], 3438, 743, 3), receta(2, 'Yogur con manzana y avena', [1, 4], 3704, 347, 5)]
const dias = ['2026-10-12', '2026-10-13', '2026-10-14', '2026-10-15', '2026-10-16', '2026-10-17', '2026-10-18']
export const SEMANA = {
  week_start: dias[0], week_end: dias[6], week_cost_cents: 6876, unplanned_slots: 20, leftovers: [{ recipe_id: 1, servings_milli: 1000 }],
  cells: [1, 2].map((p) => ({ plan_date: dias[0], meal_slot_id: 2, person_id: p, meal_plan_entry_id: 9 + p, recipe_id: 1, servings_milli: 1000, is_eating_out: false, cost_cents: 3438, kcal: 743, price_store_id: 2, price_observed_date: '2026-10-10' })),
  days: dias.map((d, i) => ({ plan_date: d, cost_cents: i ? 0 : 6876, persons: [{ person_id: 1, kcal: i ? 0 : 743, daily_kcal_target: 2364, difference_kcal: i ? -2364 : -1621 }, { person_id: 2, kcal: i ? 0 : 743, daily_kcal_target: 1900, difference_kcal: i ? -1900 : -1157 }] })),
}
const item = (id: number, product_id: number, q: number, store: number, est: number, extra = {}) => ({ id, shopping_list_id: 5, origin: 'calculated', product_id, needed_milli: q, quantity_milli: q, assigned_store_id: store, is_bought: false, created_at: T, updated_at: T, packages_to_buy: 1, estimated_price_cents: est, reasons: ['recipe'], recipe_ids: [1], ...extra })
export const LISTA = {
  shopping_list_id: 5, start_date: dias[0], end_date: dias[6], generated_at: T, total_cents: 7780,
  groups: [{ purchase_channel: 'market', store_id: 4, visit_weekday: 6, visit_date: '2026-10-17', items: [item(1, 1, 1000000, 4, 2890)], total_cents: 2890 },
    { purchase_channel: 'chain', store_id: 2, items: [item(2, 2, 2000000, 2, 4890, { is_stale_price: true, reasons: ['replenishment'] })], total_cents: 4890 }],
}
// Canon S5 (U1 §3.16): Bodega $3,094.50 contra Walmart $3,233.00 (+$138.50); idénticos +$100.50; 33 idénticos y 3 sustitutos; Bodega gana 22, empates 8, Walmart 3.
export const COMPARACION = {
  basket: 'shopping_list', cheapest_store_id: 2, ties: 8, winners: [{ store_id: 2, wins: 22 }, { store_id: 1, wins: 3 }],
  stores: [{ store_id: 1, total_cents: 323300, identical_count: 33, identical_total_cents: 287400, substitute_count: 3, missing_count: 0, stale_count: 5, difference_vs_cheapest_cents: 13850, assumes_membership: false },
    { store_id: 2, total_cents: 309450, identical_count: 33, identical_total_cents: 277350, substitute_count: 3, missing_count: 0, stale_count: 2, difference_vs_cheapest_cents: 0, assumes_membership: false }],
  items: Array.from({ length: 36 }, (_, i) => ({ product_id: (i % 3) + 1, quantity_milli: 1000, is_substitute: i >= 33, substitute_note: i >= 33 ? 'Otra presentación' : undefined, winner_store_id: i < 22 ? 2 : i < 25 ? 1 : undefined, difference_cents: 36 - i, prices: [precio(1, 9000 + i), precio(2, 8600 + i, i % 5 !== 0)] })),
  membership: { store_id: 3, basket_savings_cents: -5000, annual_savings_cents: -21240, annual_fee_cents: 50000, verdict: 'not_worth_it' },
}
const CANASTAS = [{ basket: 'shopping_list', item_count: 36 }, { basket: 'expense', expense_id: 77, expense_date: '2026-10-10', store_id: 1, total_cents: 323300, item_count: 36 }]
const INVENTARIO = [
  { id: 1, product_id: 2, tracking_mode: 'consumption_rate', quantity_milli: 3000000, counted_date: '2026-10-05', created_at: T, updated_at: T, status: 'sufficient', estimated_quantity_milli: 2000000, weekly_consumption_milli: 1500000, days_left: 9 },
  { id: 2, product_id: 1, tracking_mode: 'stock', quantity_milli: 0, counted_date: '2026-10-10', created_at: T, updated_at: T, status: 'out' },
]
const DETALLE_RECETA = (sm: number) => ({
  recipe: RECETAS[0], servings_milli: sm, total_cost_cents: Math.round((6876 * sm) / 2000), cost_per_serving_cents: 3438, total_kcal: Math.round((1486 * sm) / 2000), kcal_per_serving: 743,
  ingredients: [{ recipe_ingredient_id: 1, product_id: 1, quantity_milli: 250000 * (sm / 2000), is_optional: false, packages_to_buy: 1, cost_cents: 723, price_store_id: 2, price_observed_date: '2026-10-10', kcal: 45, nutrition_source: 'usda_fdc' },
    { recipe_ingredient_id: 2, product_id: 3, quantity_milli: 2000 * (sm / 2000), is_optional: true }],
})
export interface Pet { metodo: string; ruta: string; consulta: URLSearchParams; cuerpo: unknown; clave?: string }
/** Registra las rutas del área encima de ApiFalsa (Playwright prueba primero la última ruta registrada). */
export async function instalarArea(ctx: BrowserContext, peticiones: Pet[] = []) {
  await ctx.route(/\/api\/v1\//, async (route) => {
    const r = route.request(); const u = new URL(r.url()); const clave = `${r.method()} ${u.pathname}`
    const pet: Pet = { metodo: r.method(), ruta: u.pathname, consulta: u.searchParams, cuerpo: r.postData() ? r.postDataJSON() : null, clave: r.headers()['idempotency-key'] }
    const json = (cuerpo: unknown, status = 200) => { peticiones.push(pet); return route.fulfill({ status, contentType: 'application/json', body: JSON.stringify(cuerpo) }) }
    const tabla: Record<string, () => unknown> = {
      'GET /api/v1/stores': () => pag(TIENDAS), 'GET /api/v1/products': () => pag(PRODUCTOS), 'GET /api/v1/product-groups': () => pag(GRUPOS),
      'GET /api/v1/meal-slots': () => pag(TIEMPOS), 'GET /api/v1/recipes': () => pag(RECETAS), 'GET /api/v1/meal-plan/week': () => SEMANA,
      'GET /api/v1/meal-plan/trends': () => ({ weeks: [{ week_start: '2026-10-05', cost_cents: 98000, persons: [{ person_id: 1, avg_kcal_per_day: 2200, daily_kcal_target: 2364 }] }, { week_start: '2026-10-12', cost_cents: 6876, persons: [{ person_id: 1, avg_kcal_per_day: 743, daily_kcal_target: 2364 }] }] }),
      'GET /api/v1/meal-plan/month': () => ({ month: '2026-10', days: Array.from({ length: 31 }, (_, i) => ({ plan_date: `2026-10-${String(i + 1).padStart(2, '0')}`, status: i < 11 ? 'complete' : i === 11 ? 'incomplete' : 'empty' })) }),
      'GET /api/v1/recipes/1': () => DETALLE_RECETA(Number(u.searchParams.get('servings_milli') ?? 2000)), 'GET /api/v1/recipes/1/usage': () => ({ times_used_8w: 3, next_plan_date: '2026-10-12' }),
      'GET /api/v1/shopping-list': () => LISTA, 'GET /api/v1/baskets': () => pag(CANASTAS), 'GET /api/v1/basket-comparison': () => COMPARACION,
      'GET /api/v1/prices/summary': () => ({ fresh_count: 1, stale_count: 1, no_price_count: 1, stores: [{ store_id: 1, latest_observed_date: '2026-10-10' }, { store_id: 2, latest_observed_date: '2026-10-10' }] }),
      'GET /api/v1/inventory-items': () => pag(INVENTARIO), 'GET /api/v1/inventory/summary': () => ({ sufficient_count: 1, reorder_count: 0, out_count: 1 }),
      'POST /api/v1/meal-plan-entries': () => ({ entries: [] }), 'POST /api/v1/shopping-list/generate': () => LISTA,
      'POST /api/v1/shopping-list/assign-store': () => ({ assigned_count: 36, without_price_count: 0 }), 'PATCH /api/v1/shopping-list/items/2': () => ({ ...LISTA.groups[1].items[0], is_bought: true }),
    }
    const f = tabla[clave]
    if (!f) return route.fallback()
    return json(f(), r.method() === 'POST' && clave.endsWith('meal-plan-entries') ? 201 : 200)
  })
}
