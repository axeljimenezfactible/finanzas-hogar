// Pruebas de P-26 Ajustes, secciones 01 a 07 (front-ajustes-a). Corren en 412×915 táctil (movil) y 1440×900 (escritorio).
// API simulada: ApiFalsa de la base + rutas del área con datos sintéticos (cifras del oráculo citadas en U4).
import { expect, test, type BrowserContext, type Page } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { controlesChicos, dineroSinTabulares, textosEnIngles } from '../apoyo/medidas'
import { moverTarjetaA } from './tablero'
import type { Cuerpo, Esquemas, Respuesta, RutasCon } from '../../src/api/tipos'
import { erroresC1 } from './validar-c1'
/** Ruta del contrato con cada {parámetro} cambiado por un número, p. ej. /api/v1/credits/1 (para tipar las respuestas falsas contra C1, H-11). */
type Concreta<S extends string> = S extends `${infer A}{${string}}${infer B}` ? `${A}${number}${Concreta<B>}` : S
type DatosC1 = { [P in RutasCon<'get'> & string as Concreta<P>]?: Respuesta<P, 'get'> }
const t0 = '2026-10-01T09:00:00-06:00'
const sello = { created_at: t0, updated_at: t0 }
const param = (key: string, value: number, unit: string) => ({ key, value_type: 'integer' as const, unit, value, initial_value: value, is_initial: true })
const pagina = <T,>(items: T[]) => ({ items, next_cursor: null })
/** Fecha de hace n días en la zona de la casa (para la frescura de precios, relativa a hoy). */
const haceDias = (n: number) => new Date(Date.now() - n * 86_400_000).toLocaleDateString('en-CA', { timeZone: 'America/Mexico_City' })
const SEGURIDAD = { parameters: { 'auth.max_failed_attempts': param('auth.max_failed_attempts', 5, 'intentos'), 'auth.lockout_window_minutes': param('auth.lockout_window_minutes', 15, 'minutos'), 'session.max_age_days': param('session.max_age_days', 90, 'días'), 'auth.password_min_length': param('auth.password_min_length', 12, 'caracteres') } } satisfies Respuesta<'/api/v1/settings/security', 'get'>
const KUESKI = { id: 1, name: 'Kueski Pay', kind: 'revolving_line', holder_person_id: 1, credit_limit_cents: 1000000, payment_frequency: 'fortnightly', availability_reading: 'B', availability_confirmed: false, reading_a_anchor_date: '2026-10-10', is_active: true, live_debt_cents: 309450, available_cents: 690550, next_payment_date: '2026-10-30', next_payment_cents: 97375, live_purchase_count: 1, missing_terms: [], ...sello } satisfies Respuesta<'/api/v1/credits/{credit_id}', 'get'>
const DATOS: Record<string, unknown> = {
  '/api/v1/users': pagina([{ id: 1, person_id: 1, username: 'axel', display_name: 'Axel', participation: 'holder', is_active: true, last_login_at: '2026-10-12T08:00:00-06:00', ...sello }, { id: 2, person_id: 2, username: 'ximena', display_name: 'Ximena', participation: 'co_holder', is_active: true, ...sello }]),
  '/api/v1/sessions': pagina([{ id: 7, user_account_id: 1, device_label: 'Teléfono de Axel', last_seen_at: '2026-10-12T09:00:00-06:00', expires_at: '2027-01-10T09:00:00-06:00', is_current: true, username: 'axel', ...sello }, { id: 8, user_account_id: 2, device_label: 'Laptop de Ximena', last_seen_at: '2026-10-11T21:00:00-06:00', expires_at: '2027-01-09T21:00:00-06:00', username: 'ximena', ...sello }]),
  '/api/v1/settings/security': SEGURIDAD,
  '/api/v1/budgets': pagina([{ id: 1, category_id: 1, amount_cents: 400000, valid_from_date: '2026-10-01', ...sello }]),
  '/api/v1/cost-share-rules': pagina([{ id: 1, name: 'Hogar', remainder_person_id: 1, is_template: false, is_default: true, is_active: true, parts: [], ...sello }, { id: 2, name: 'Conjunto 50/50', remainder_person_id: 1, is_template: false, is_default: false, is_active: true, parts: [{ id: 1, rule_id: 2, person_id: 3, pct_bp: 5000, ...sello }], ...sello }]),
  '/api/v1/planned-items': pagina([
    { id: 1, kind: 'fixed_expense', description: 'Papel', amount_cents: 19900, frequency: 'monthly', interval_count: 1, start_date: '2026-10-30', category_id: 6, cost_share_rule_id: 2, accrual_months: 1, is_active: true, ...sello },
    { id: 2, kind: 'income', description: 'Ingreso quincenal', amount_cents: 1500000, frequency: 'fortnightly', interval_count: 1, start_date: '2026-10-15', person_id: 1, accrual_months: 1, is_active: true, ...sello },
    { id: 3, kind: 'fixed_expense', description: 'Despensa bimestral', amount_cents: 309450, frequency: 'monthly', interval_count: 2, start_date: '2026-10-10', category_id: 1, payment_method: 'credit', credit_id: 1, term_count: 4, accrual_months: 2, is_active: true, ...sello }]),
  '/api/v1/settings/projection': { repeat_cadence_enabled: false },
  '/api/v1/credits': pagina([
    KUESKI,
    { id: 2, name: 'Tarjeta de la casa', kind: 'bank_card', holder_person_id: 1, payment_frequency: 'monthly', availability_reading: 'B', availability_confirmed: true, is_active: true, missing_terms: ['annual_rate', 'rate_base', 'vat', 'min_payment_pct_bp', 'min_payment_fixed_cents'], ...sello }]),
  '/api/v1/credits/1': KUESKI,
  '/api/v1/credits/1/availability-preview': { credit_id: 1, availability_reading: 'B', availability_confirmed: false, reading_a: { available_today_cents: 1000000, horizon_min_cents: 809450 }, reading_b: { available_today_cents: 690550, horizon_min_cents: 500000 } },
  '/api/v1/credits/1/availability-series': { credit_id: 1, events: [], min_a_cents: 809450, min_b_cents: 500000 },
  '/api/v1/credits/1/term-rates': pagina([{ id: 1, credit_id: 1, term_count: 2, rate_num: 1727, rate_den: 10000, rate_decimals: 2, rate_pct: '17.27', source: 'quote', ...sello }, { id: 2, credit_id: 1, term_count: 4, rate_num: 80050, rate_den: 309450, rate_decimals: 4, rate_pct: '25.8685', source: 'quote', quote_date: '2026-10-10', ...sello }]),
  '/api/v1/credit-purchases': pagina([{ id: 1, credit_id: 1, kind: 'installment_plan', purchase_date: '2026-10-10', principal_cents: 309450, interest_cents: 80050, total_cents: 389500, first_payment_date: '2026-10-30', term_count: 4, payment_frequency: 'fortnightly', interest_source: 'quoted', rate_num: 80050, rate_den: 309450, rate_basis: 'term_total', interest_includes_vat: true, status: 'active', ...sello }]),
  '/api/v1/promotions': pagina([{ id: 1, merchant_name: 'Walmart', store_id: 1, term_count: 3, rate_num: 0, rate_den: 1, min_purchase_cents: 100000, valid_from: '2026-10-01', valid_to: '2026-10-31', is_active: true, is_current: true, ...sello }]),
  '/api/v1/settings/optimizer': { missing_constraints: ['min_available_buffer_cents', 'max_live_purchases', 'zero_debt_by_date'], parameters: { 'comparator.timeout_seconds': param('comparator.timeout_seconds', 5, 'segundos'), 'optimizer.date_window_days': param('optimizer.date_window_days', 7, 'días') }, weights: { lowest_interest: 1, even_low_payment: 1, max_available: 1, most_value_per_peso: 1 } },
  '/api/v1/stores': pagina([{ id: 1, kind: 'chain', name: 'Walmart', assisted_reading_enabled: false, assisted_reading_verdict: 'missing', requires_membership: false, has_membership: false, is_active: true, ...sello }, { id: 3, kind: 'chain', name: "Sam's Club", assisted_reading_enabled: false, assisted_reading_verdict: 'missing', requires_membership: true, has_membership: false, is_active: true, ...sello }, { id: 4, kind: 'market', name: 'Tianguis del martes', visit_weekday: 2, assisted_reading_enabled: false, requires_membership: false, has_membership: false, is_active: true, product_group_ids: [1], ...sello }]),
  '/api/v1/products': pagina([{ id: 1, name: 'Aceite de cocina', base_unit: 'ml', is_active: true, last_observed_date: haceDias(25), ...sello }, { id: 2, name: 'Arroz', base_unit: 'g', is_active: true, last_observed_date: haceDias(10), ...sello }, { id: 3, name: 'Frijol', base_unit: 'g', is_active: true, last_observed_date: haceDias(40), ...sello }, { id: 4, name: 'Jitomate', base_unit: 'g', is_active: true, ...sello }]),
  '/api/v1/product-groups': pagina([{ id: 1, name: 'Verduras', purchase_channel: 'market', sort_order: 1, is_active: true, ...sello }]),
  '/api/v1/price-sources': { manual_enabled: true, ocr_engine: 'Tesseract 5.5.3', assisted_reading_enabled: false, sites: [], profeco_price_count: 0, profeco_coverage: [] },
  '/api/v1/settings/prices': { assisted_reading_enabled: false, parameters: { 'price.stale_after_days': param('price.stale_after_days', 30, 'días'), 'price.zero_weight_after_days': param('price.zero_weight_after_days', 90, 'días'), 'price.profeco_import_max_gb': param('price.profeco_import_max_gb', 1, 'GB') } },
  '/api/v1/ai/providers': { providers: [{ provider: 'local_model', is_enabled: false, allowed_uses: [], is_available: false, is_configured: false }] },
  '/api/v1/settings/recipes': { default_base_servings: 2, meal_slots: [{ id: 1, code: 'breakfast', name: 'Desayuno', sort_order: 1, is_optional: false, ...sello }, { id: 2, code: 'lunch', name: 'Comida', sort_order: 2, is_optional: false, ...sello }], recipe_counts: [{ origin: 'claude_seed', count: 12 }, { origin: 'household', count: 3 }], nutrition_sources: [{ source: 'usda_fdc', food_count: 0 }] },
  '/api/v1/recipes': pagina([{ id: 1, name: 'Tinga de pollo', base_servings: 2, origin: 'claude_seed', is_active: true, meal_slot_ids: [2], ...sello }]),
  '/api/v1/health-profiles': pagina([{ id: 1, person_id: 1, ...sello }, { id: 2, person_id: 2, ...sello }]),
} satisfies DatosC1
/** previewCreditPurchase falsa: quincenal da la serie del oráculo que cita U4 P-26.05 punto 6; mensual, fechas sintéticas solo de prueba. */
const vistaPreviaCompra = (frecuencia?: string): Respuesta<'/api/v1/credit-purchases/preview', 'post'> => {
  const fechas = frecuencia === 'monthly' ? ['2026-11-10', '2026-12-10', '2027-01-10', '2027-02-10'] : ['2026-10-30', '2026-11-15', '2026-11-30', '2026-12-15']
  return { principal_cents: 309450, interest_cents: 80050, total_cents: 389500, shares: [], schedule: fechas.map((due_date, i) => ({ seq: i + 1, due_date, amount_cents: 97375, principal_cents: i === 3 ? 77361 : 77363, interest_cents: i === 3 ? 20014 : 20012, shares: [] })) }
}
/** Respuesta falsa de cada escritura que ejercen las pruebas, con el código y el esquema de C1 (R2-01). undefined = sin respuesta falsa. */
function respuestaEscritura(metodo: string, ruta: string, cuerpo: unknown): [number, unknown?] | undefined {
  const id = Number(ruta.split('/')[4])
  switch (`${metodo} ${ruta.replace(/\/\d+(?=\/|$)/g, '/{id}')}`) {
    case 'PATCH /api/v1/settings/security': {
      const cambios = (cuerpo as Cuerpo<'/api/v1/settings/security', 'patch'>).parameters ?? {}
      return [200, { parameters: Object.fromEntries(Object.entries(SEGURIDAD.parameters).map(([k, p]) => [k, k in cambios ? { ...p, value: Number(cambios[k as keyof typeof cambios]), is_initial: Number(cambios[k as keyof typeof cambios]) === p.initial_value } : p])) as typeof SEGURIDAD.parameters } satisfies Respuesta<'/api/v1/settings/security', 'patch'>]
    }
    case 'PATCH /api/v1/budgets/{id}': return [200, { id, category_id: 1, amount_cents: (cuerpo as Cuerpo<'/api/v1/budgets/{budget_id}', 'patch'>).amount_cents ?? 400000, valid_from_date: '2026-10-01', ...sello } satisfies Respuesta<'/api/v1/budgets/{budget_id}', 'patch'>]
    case 'PATCH /api/v1/credits/{id}': {
      const c = cuerpo as Cuerpo<'/api/v1/credits/{credit_id}', 'patch'>
      return [200, { ...KUESKI, id, ...(c.availability_reading ? { availability_reading: c.availability_reading } : {}), ...(c.payment_frequency ? { payment_frequency: c.payment_frequency } : {}) } satisfies Respuesta<'/api/v1/credits/{credit_id}', 'patch'>]
    }
    case 'POST /api/v1/promotions': {
      const c = cuerpo as Cuerpo<'/api/v1/promotions', 'post'>
      return [201, { id: 2, merchant_name: c.merchant_name, credit_id: c.credit_id ?? undefined, term_count: c.term_count, rate_num: c.rate_num ?? 0, rate_den: c.rate_den ?? 1, min_purchase_cents: c.min_purchase_cents ?? 0, valid_from: c.valid_from, valid_to: c.valid_to, is_active: c.is_active ?? true, is_current: true, ...sello } satisfies Respuesta<'/api/v1/promotions', 'post'>]
    }
    case 'POST /api/v1/recipes': return [201, { id: 3, name: 'Caldo de verduras', base_servings: 2, origin: 'household', is_active: true, meal_slot_ids: [2], ...sello } satisfies Respuesta<'/api/v1/recipes', 'post'>]
    case 'POST /api/v1/recipes/{id}/duplicate': return [201, { id: 2, name: 'Tinga de pollo', base_servings: 2, origin: 'household', is_active: true, meal_slot_ids: [2], ...sello } satisfies Respuesta<'/api/v1/recipes/{recipe_id}/duplicate', 'post'>]
    case 'PATCH /api/v1/categories/{id}': return [200, { id, name: 'Despensa', sort_order: id, is_active: true, default_cost_share_rule_id: (cuerpo as Cuerpo<'/api/v1/categories/{category_id}', 'patch'>).default_cost_share_rule_id ?? undefined, ...sello } satisfies Respuesta<'/api/v1/categories/{category_id}', 'patch'>]
    case 'PATCH /api/v1/planned-items/{id}': {
      const c = cuerpo as Cuerpo<'/api/v1/planned-items/{planned_item_id}', 'patch'>
      return [200, { id, kind: 'fixed_expense', description: 'Partida', amount_cents: 19900, frequency: 'monthly', interval_count: 1, start_date: '2026-10-30', category_id: c.category_id ?? 6, accrual_months: 1, is_active: c.is_active ?? true, ...sello } satisfies Respuesta<'/api/v1/planned-items/{planned_item_id}', 'patch'>]
    }
    case 'PATCH /api/v1/promotions/{id}': {
      const c = cuerpo as Cuerpo<'/api/v1/promotions/{promotion_id}', 'patch'>
      return [200, { id, merchant_name: 'Promoción', term_count: c.term_count ?? 3, rate_num: 0, rate_den: 1, min_purchase_cents: c.min_purchase_cents ?? 100000, valid_from: '2026-10-01', valid_to: c.valid_to ?? '2026-10-31', is_active: c.is_active ?? true, is_current: c.is_active ?? true, ...sello } satisfies Respuesta<'/api/v1/promotions/{promotion_id}', 'patch'>]
    }
    case 'PATCH /api/v1/recipes/{id}': return [200, { id, name: 'Receta', base_servings: 2, origin: 'claude_seed', is_active: true, meal_slot_ids: (cuerpo as Cuerpo<'/api/v1/recipes/{recipe_id}', 'patch'>).meal_slot_ids ?? [2], ...sello } satisfies Respuesta<'/api/v1/recipes/{recipe_id}', 'patch'>]
    case 'POST /api/v1/sessions/revoke': return [200, { revoked_count: 1 } satisfies Respuesta<'/api/v1/sessions/revoke', 'post'>]
  }
  return metodo === 'DELETE' ? [204] : undefined
}
/** Errores de C1 del tráfico falso de la prueba en curso: cuerpos de petición y respuestas con su código (R2-01). Cada prueba termina con cero. */
let erroresTrafico: string[] = []
test.afterEach(() => expect(erroresTrafico).toEqual([]))
async function preparar(context: BrowserContext) {
  await new ApiFalsa().instalar(context)
  erroresTrafico = []
  const escrituras: { metodo: string; url: string; cuerpo: unknown }[] = []
  await context.route(/\/api\/v1\//, async (route) => {
    const pet = route.request(), u = new URL(pet.url())
    const cuerpo: unknown = pet.postData() ? pet.postDataJSON() : null
    const responder = (estado: number, json?: unknown) => {
      erroresTrafico.push(...erroresC1(pet.method(), u.pathname, estado, cuerpo, json ?? null))
      return json === undefined ? route.fulfill({ status: estado }) : route.fulfill({ status: estado, json, headers: { 'X-Request-Id': 'prueba' } })
    }
    if (pet.method() === 'GET') return DATOS[u.pathname] !== undefined ? responder(200, DATOS[u.pathname]) : route.fallback()
    if (['/api/v1/expenses', '/api/v1/weight-entries', '/api/v1/auth/login', '/api/v1/auth/logout'].includes(u.pathname)) return route.fallback()
    if (u.pathname === '/api/v1/credit-purchases/preview') return responder(200, vistaPreviaCompra((cuerpo as { payment_frequency?: string }).payment_frequency))
    escrituras.push({ metodo: pet.method(), url: u.pathname + u.search, cuerpo })
    const r = respuestaEscritura(pet.method(), u.pathname, cuerpo)
    if (!r) { erroresTrafico.push(`${pet.method()} ${u.pathname}: la prueba no tiene respuesta falsa de C1 para esta escritura`); return route.fulfill({ status: 501 }) }
    return responder(r[0], r[1])
  })
  return escrituras
}
const sinDesborde = async (page: Page) => expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth + 1)).toBe(true)
const SECCIONES: [string, string, RegExp[]][] = [
  ['01', 'Cuenta y roles', [/Participante externa, sin cuenta/, /Laptop de Ximena/, /Parámetros de acceso/]],
  ['02', 'Presupuesto', [/Mes calendario; quincenas/, /Gastos fijos/, /Despensa programada/]],
  ['03', 'Ingresos', [/Monto de la quincena del 15/, /\$15,000\.00/, /No aporta ingresos al modelo/]],
  ['04', 'Categorías', [/Despensa/, /El de por defecto/, /Comida fuera/]],
  ['05', 'Créditos y condiciones', [/Kueski Pay/, /Falta capturar/, /Faltan restricciones/]],
  ['06', 'Tiendas y precios', [/Tianguis del martes/, /Fuentes de precio/, /queda viejo el 9-nov-2026/]],
  ['07', 'Recetas', [/Tinga de pollo/, /Bases nutricionales/, /Restricciones alimentarias/]],
]

test('P-26 muestra las 14 secciones de Ajustes', async ({ page, context }) => {
  await preparar(context)
  await page.goto('/ajustes')
  const nav = page.getByRole('navigation', { name: 'Secciones de Ajustes' })
  await expect(nav.getByRole('link')).toHaveCount(14)
  await expect(nav.getByRole('link', { name: /Cuenta y roles/ })).toBeVisible()
  await sinDesborde(page)
})
for (const [id, nombre, textos] of SECCIONES) {
  test(`P-26.${id} ${nombre} monta sus elementos principales`, async ({ page, context }) => {
    await preparar(context)
    await page.goto(`/ajustes/${id}`)
    await expect(page.getByRole('heading', { level: 2, name: nombre })).toBeVisible()
    for (const re of textos) await expect(page.getByText(re).first()).toBeVisible()
    await sinDesborde(page)
    if (test.info().project.name === 'movil') {
      const alturas = await page.locator('article button:visible').evaluateAll((bs) => bs.slice(0, 12).map((b) => b.getBoundingClientRect().height))
      for (const h of alturas) expect(h).toBeGreaterThanOrEqual(44)
    }
  })
}
test('P-26.01 parámetro K44 con vista previa «Antes → Después» y guardado', async ({ page, context }) => {
  const escrituras = await preparar(context)
  await page.goto('/ajustes/01')
  await page.getByRole('button', { name: /Editar Intentos fallidos permitidos/ }).click()
  const campo = page.getByRole('textbox', { name: 'Intentos fallidos permitidos en la ventana' })
  await campo.fill('6')
  await expect(page.getByTestId('vista-previa')).toContainText('el intento 7 fallido')
  await campo.press('Enter')
  await expect(page.getByText('Guardado')).toBeVisible()
  expect(escrituras.at(-1)).toMatchObject({ metodo: 'PATCH', url: '/api/v1/settings/security', cuerpo: { parameters: { 'auth.max_failed_attempts': 6 } } })
})
test('P-26.01 cerrar sesiones es acción crítica con «Cancelar» enfocado', async ({ page, context }) => {
  await preparar(context)
  await page.goto('/ajustes/01')
  await page.getByLabel('Qué sesiones cerrar').selectOption('all_except_current')
  await page.getByRole('button', { name: 'Cerrar sesiones' }).click()
  await expect(page.getByRole('alertdialog')).toContainText('Cerrar 1 sesión')
  await expect(page.getByRole('button', { name: 'Cancelar' })).toBeFocused()
})
test('P-26.02 cambiar un monto pide el alcance con «De este mes en adelante» preseleccionado', async ({ page, context }) => {
  const escrituras = await preparar(context)
  await page.goto('/ajustes/02')
  await page.getByRole('button', { name: /Editar Monto Despensa/ }).click()
  await page.getByRole('textbox', { name: 'Monto Despensa' }).fill('4500')
  await page.getByRole('textbox', { name: 'Monto Despensa' }).press('Enter')
  await expect(page.getByRole('radio', { name: 'De este mes en adelante' })).toBeChecked()
  await page.getByRole('button', { name: 'Confirmar' }).click()
  await expect.poll(() => escrituras.at(-1)?.url ?? '').toContain('scope=onward')
  expect(escrituras.at(-1)?.cuerpo).toEqual({ amount_cents: 450000 })
})
test('R4-01 · P-26.02 un 409 budget_overlap al confirmar el monto muestra V-17 con la hoja abierta; otro fallo muestra «No se pudo guardar…» y Reintentar guarda', async ({ page, context }) => {
  const escrituras = await preparar(context)
  const fallosPagina: string[] = []
  page.on('pageerror', (e) => fallosPagina.push(e.message))
  let modo: 'traslape' | 'falla' | 'bien' = 'traslape'
  const intentos: unknown[] = []
  await context.route(/\/api\/v1\/budgets\/\d+/, async (route) => {
    const pet = route.request(), u = new URL(pet.url())
    if (pet.method() !== 'PATCH' || modo === 'bien') return route.fallback()
    const cuerpo: unknown = pet.postDataJSON()
    intentos.push(cuerpo)
    if (modo === 'traslape') {
      const json = { code: 'budget_overlap', message: 'Ya hay un monto de Despensa en esos meses.', request_id: 'prueba' } satisfies Esquemas['Error']
      erroresTrafico.push(...erroresC1('PATCH', u.pathname, 409, cuerpo, json))
      return route.fulfill({ status: 409, json, headers: { 'X-Request-Id': 'prueba' } })
    }
    // 500 internal_error: C1 lo declara en `default` con el esquema Error; validar-c1 solo lee códigos explícitos, así que este no pasa por erroresC1.
    return route.fulfill({ status: 500, json: { code: 'internal_error', message: 'Error interno del servidor.', request_id: 'prueba' } satisfies Esquemas['Error'], headers: { 'X-Request-Id': 'prueba' } })
  })
  await page.goto('/ajustes/02')
  await page.getByRole('button', { name: /Editar Monto Despensa/ }).click()
  await page.getByRole('textbox', { name: 'Monto Despensa' }).fill('4500')
  await page.getByRole('textbox', { name: 'Monto Despensa' }).press('Enter')
  const hoja = page.getByRole('dialog')
  await hoja.getByRole('radio', { name: 'Solo este mes' }).check()
  await hoja.getByRole('button', { name: 'Confirmar' }).click()
  // 409 → V-17 debajo del monto; la hoja sigue abierta con el monto y el alcance capturados.
  await expect(hoja.getByRole('alert').filter({ hasText: 'Esta categoría ya tiene un monto en esos meses. Ajusta el que existe.' })).toBeVisible()
  await expect(hoja.getByRole('radio', { name: 'Solo este mes' })).toBeChecked()
  await expect(hoja).toContainText('$4,500.00')
  // Cualquier otro fallo → texto de U4 §0.3 con Reintentar, sin perder lo capturado.
  modo = 'falla'
  await hoja.getByRole('button', { name: 'Confirmar' }).click()
  const fallo = hoja.getByRole('alert').filter({ hasText: 'No se pudo guardar. Revisa la conexión con el servidor de casa' })
  await expect(fallo).toBeVisible()
  await expect(hoja.getByText('Esta categoría ya tiene un monto en esos meses.')).toHaveCount(0)
  await expect(hoja.getByRole('radio', { name: 'Solo este mes' })).toBeChecked()
  modo = 'bien'
  await fallo.getByRole('button', { name: 'Reintentar' }).click()
  await expect(hoja).toBeHidden()
  await expect.poll(() => escrituras.at(-1)?.url ?? '').toContain('scope=this_month')
  expect(escrituras.at(-1)?.cuerpo).toEqual({ amount_cents: 450000 })
  expect(intentos).toEqual([{ amount_cents: 450000 }, { amount_cents: 450000 }])
  expect(fallosPagina).toEqual([])
})
test('P-26.05 interruptor A/B muestra las dos lecturas y confirma', async ({ page, context }) => {
  const escrituras = await preparar(context)
  await page.goto('/ajustes/05')
  await page.getByText('Kueski Pay').first().click()
  await expect(page.getByText('Pendiente confirmar con el acreedor').first()).toBeVisible()
  await page.getByRole('button', { name: 'Usar la lectura A' }).click()
  const lecturas = page.getByTestId('lecturas-ab')
  await expect(lecturas).toContainText('$10,000.00')
  await expect(lecturas).toContainText('$6,905.50')
  await expect(lecturas).toContainText('$5,000.00')
  await page.getByRole('dialog', { name: 'Lectura del disponible' }).getByRole('button', { name: 'Usar la lectura A' }).click()
  await expect.poll(() => escrituras.at(-1)?.cuerpo).toEqual({ availability_reading: 'A' })
})
test('P-26.06 días para peso cero menores que precio viejo da V-25', async ({ page, context }) => {
  await preparar(context)
  await page.goto('/ajustes/06')
  await page.getByRole('button', { name: /Editar Días en que el peso/ }).click()
  await page.getByRole('textbox', { name: /Días en que el peso/ }).fill('20')
  await page.getByRole('textbox', { name: /Días en que el peso/ }).press('Enter')
  await expect(page.getByText('Escribe 30 días o más')).toBeVisible()
})

test('P-26.05 vista previa 6: frecuencia de pagos con la compra de ejemplo', async ({ page, context }) => {
  const escrituras = await preparar(context)
  await page.goto('/ajustes/05')
  await page.getByText('Kueski Pay').first().click()
  await page.getByRole('button', { name: /Editar Frecuencia de pagos/ }).click()
  await expect(page.getByLabel('Capital de ejemplo')).toHaveValue('3094.50')
  await expect(page.getByLabel('Plazo de ejemplo (pagos)')).toHaveValue('4')
  await page.getByRole('combobox', { name: 'Frecuencia de pagos' }).selectOption('monthly')
  const previa = page.getByTestId('vista-previa')
  await expect(previa).toContainText('30-oct-2026 · 15-nov-2026 · 30-nov-2026 · 15-dic-2026')
  await expect(previa).toContainText('$973.75 cada uno (4 pagos)')
  await expect(previa).toContainText('$3,895.00')
  await expect(previa).toContainText('10-nov-2026')
  await page.getByRole('combobox', { name: 'Frecuencia de pagos' }).press('Enter')
  await expect.poll(() => escrituras.at(-1)?.cuerpo).toEqual({ payment_frequency: 'monthly' })
})
test('P-26.05 sin compras pide la compra de ejemplo y no guarda la frecuencia', async ({ page, context }) => {
  const escrituras = await preparar(context)
  await context.route('**/api/v1/credit-purchases?*', (r) => (r.request().method() === 'GET' ? r.fulfill({ json: pagina([]) }) : r.fallback()))
  await page.goto('/ajustes/05')
  await page.getByText('Kueski Pay').first().click()
  await page.getByRole('button', { name: /Editar Frecuencia de pagos/ }).click()
  await page.getByRole('combobox', { name: 'Frecuencia de pagos' }).selectOption('weekly')
  await expect(page.getByTestId('vista-previa')).toContainText('Escribe un monto y un plazo de ejemplo')
  await page.getByRole('combobox', { name: 'Frecuencia de pagos' }).press('Enter')
  expect(escrituras.filter((e) => e.metodo === 'PATCH')).toHaveLength(0)
  await page.getByLabel('Capital de ejemplo').fill('3094.50')
  await page.getByLabel('Plazo de ejemplo (pagos)').fill('4')
  await expect(page.getByTestId('vista-previa')).toContainText('$3,895.00')
})
test('P-26.05 vista previa 7: primer pago con el valor guardado y «Después» pendiente', async ({ page, context }) => {
  await preparar(context)
  await page.goto('/ajustes/05')
  await page.getByText('Kueski Pay').first().click()
  await page.getByRole('button', { name: /Editar Días mínimos entre compra y primer pago/ }).click()
  await page.getByRole('textbox', { name: 'Días mínimos entre compra y primer pago' }).fill('10')
  const previa = page.getByTestId('vista-previa')
  await expect(previa).toContainText('Primer pago de una compra hecha hoy')
  await expect(previa).toContainText('30-oct-2026')
  await expect(previa).toContainText('Vista previa pendiente')
  await expect(previa).toContainText('PROVISIONAL (SC-U4-1)')
})
test('P-26.05 alta de promoción con Ctrl+Enter', async ({ page, context }) => {
  const escrituras = await preparar(context)
  await page.goto('/ajustes/05')
  await page.getByRole('button', { name: 'Nueva promoción' }).first().click()
  const hoja = page.getByRole('dialog', { name: 'Nueva promoción' })
  await hoja.getByLabel('Comercio').fill('Chedraui')
  await hoja.getByLabel('Línea').selectOption({ label: 'Kueski Pay' })
  await hoja.getByLabel('Plazo (pagos)').fill('4')
  await hoja.getByLabel('Tasa promocional por todo el plazo').fill('0')
  await hoja.getByLabel('Vigencia: inicio').fill('2026-10-30')
  await hoja.getByLabel('Vigencia: fin').fill('2026-10-30')
  await hoja.getByLabel('Vigencia: fin').press('Control+Enter')
  await expect.poll(() => escrituras.at(-1)?.url).toBe('/api/v1/promotions')
  expect(escrituras.at(-1)?.cuerpo).toMatchObject({ merchant_name: 'Chedraui', credit_id: 1, term_count: 4, rate_num: 0, valid_from: '2026-10-30', valid_to: '2026-10-30', is_active: true })
})
test('P-26.07 duplicar receta abre P-12 de la copia', async ({ page, context }) => {
  const escrituras = await preparar(context)
  await page.goto('/ajustes/07')
  await page.getByRole('button', { name: 'Duplicar receta', exact: true }).last().click()
  await expect.poll(() => escrituras.at(-1)?.url).toBe('/api/v1/recipes/1/duplicate')
  await expect(page).toHaveURL(/\/comidas\/recetas\/2$/)
})

test('P-26.06 calendario por día de visita: el mercado en su día y ←/→ cambia de semana', async ({ page, context }) => {
  await preparar(context)
  await page.goto('/ajustes/06')
  const cal = page.getByTestId('calendario-visitas')
  const martes = cal.getByRole('listitem').nth(1)
  await expect(martes).toContainText('martes')
  await expect(martes.getByRole('button', { name: 'Tianguis del martes' })).toBeVisible()
  await expect(cal.getByRole('listitem').nth(0).getByRole('button')).toHaveCount(0)
  const semana = cal.getByText(/^Semana del /)
  const antes = await semana.textContent()
  await cal.getByRole('button', { name: 'Semana siguiente' }).focus()
  await page.keyboard.press('ArrowRight')
  await expect(semana).not.toHaveText(antes ?? '')
  await page.keyboard.press('ArrowLeft')
  await expect(semana).toHaveText(antes ?? '')
  await martes.getByRole('button', { name: 'Tianguis del martes' }).click()
  await expect(page.getByRole('dialog')).toContainText('Tianguis del martes')
})

test('P-26.06 vista previa de frescura: cuántos pasan hoy de fresco a viejo (PROVISIONAL)', async ({ page, context }) => {
  await preparar(context)
  await page.goto('/ajustes/06')
  await page.getByRole('button', { name: /Editar Días para que un precio sea viejo/ }).click()
  await page.getByRole('textbox', { name: /Días para que un precio sea viejo/ }).fill('20')
  await expect(page.getByTestId('vista-previa')).toContainText('queda viejo')
  await expect(page.getByTestId('conteo-viejos')).toHaveText('Con este valor, 1 producto pasa hoy de precio fresco a viejo · PROVISIONAL (bloqueo WF3-ajustes-precios-conteo)')
})

test('P-26.07 nueva receta con Alt+N: pide nombre y tiempo apto, crea y abre P-12', async ({ page, context }) => {
  const escrituras = await preparar(context)
  await page.goto('/ajustes/07')
  await expect(page.getByText('Tinga de pollo').first()).toBeVisible()
  await page.keyboard.press('Alt+n')
  const hoja = page.getByRole('dialog')
  await hoja.getByRole('textbox', { name: 'Nombre' }).fill('Caldo de verduras')
  await hoja.getByRole('textbox', { name: 'Nombre' }).press('Control+Enter')
  await expect(hoja.getByText('Elige al menos un tiempo de comida.')).toBeVisible()
  expect(escrituras.filter((e) => e.metodo === 'POST')).toHaveLength(0)
  await hoja.getByRole('checkbox', { name: 'Comida' }).check()
  await hoja.getByRole('button', { name: 'Guardar' }).click()
  await expect.poll(() => escrituras.at(-1)?.url).toBe('/api/v1/recipes')
  expect(escrituras.at(-1)?.cuerpo).toEqual({ name: 'Caldo de verduras', meal_slot_ids: [2] })
  await expect(page).toHaveURL(/\/comidas\/recetas\/3$/)
})

test('P-26.04 tablero: soltar en otro reparto muestra la vista previa; Esc la regresa sin guardar y Enter aplica', async ({ page, context }) => {
  const escrituras = await preparar(context)
  await page.goto('/ajustes/04')
  await page.getByRole('button', { name: 'Tablero', exact: true }).click()
  await moverTarjetaA(page, 'Despensa', 'Conjunto 50/50')
  const previa = page.getByTestId('previa-tablero')
  await expect(previa).toContainText('«Despensa» pasa a Conjunto 50/50')
  await expect(previa).toContainText('$99.50')
  await expect(page.getByRole('region', { name: 'Conjunto 50/50' }).getByText('Despensa')).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(previa).toHaveCount(0)
  await expect(page.getByRole('region', { name: 'Hogar' }).getByText('Despensa')).toBeVisible()
  expect(escrituras.filter((e) => e.metodo === 'PATCH')).toHaveLength(0)
  await moverTarjetaA(page, 'Despensa', 'Conjunto 50/50')
  await expect(page.getByRole('button', { name: 'Aplicar' })).toBeFocused()
  await page.keyboard.press('Enter')
  await expect.poll(() => escrituras.at(-1)).toMatchObject({ metodo: 'PATCH', url: '/api/v1/categories/1', cuerpo: { default_cost_share_rule_id: 2 } })
  await expect(previa).toHaveCount(0)
})

test('U5-06 en 412 táctil: flechas de semana de 44 px en 06 (R2-02) y aviso del alcance de 02 con cifras tabulares (H-09)', async ({ page, context }) => {
  test.skip(test.info().project.name !== 'movil', 'medida táctil')
  await preparar(context)
  await page.goto('/ajustes/06')
  await expect(page.getByTestId('calendario-visitas')).toBeVisible()
  const { chicos } = await controlesChicos(page)
  expect(chicos.filter((c) => /Semana/.test(c.control))).toEqual([])
  for (const n of ['Semana anterior', 'Semana siguiente']) {
    const caja = await page.getByRole('button', { name: n }).boundingBox()
    expect(Math.min(caja!.width, caja!.height)).toBeGreaterThanOrEqual(44)
  }
  // H-03 (frontend-base): en < 768 el calendario de promociones de 05 es agenda por día; el evento mide ≥ 44 × 44.
  await page.goto('/ajustes/05')
  const evento = page.getByRole('button', { name: /Walmart · Sin intereses/ }).first()
  await expect(evento).toBeVisible()
  const ce = (await evento.boundingBox())!
  console.log(`H-03 /ajustes/05 evento «Walmart · Sin intereses» a 412 táctil: ${Math.round(ce.width)}×${Math.round(ce.height)}`)
  expect(Math.min(ce.width, ce.height)).toBeGreaterThanOrEqual(44)
  expect((await controlesChicos(page)).chicos).toEqual([])
  await page.goto('/ajustes/02')
  await page.getByRole('button', { name: /Editar Monto Despensa/ }).click()
  await page.getByRole('textbox', { name: 'Monto Despensa' }).fill('4500')
  await page.getByRole('textbox', { name: 'Monto Despensa' }).press('Enter')
  await expect(page.getByText(/De este mes en adelante: oct-2026/).first()).toBeVisible()
  expect(await dineroSinTabulares(page)).toEqual([])
})

// ——— WF-3a bis (front-ajustes-a-fix): Ctrl+/ con vista previa sumada (H-07), N-284, R2-05 y R2-07.
/** Respuesta GET propia de una prueba (Playwright prueba primero la última ruta registrada); también se valida contra C1. */
async function conGet<P extends keyof DatosC1>(context: BrowserContext, ruta: P, json: NonNullable<DatosC1[P]>, alPedir?: (u: URL) => void) {
  await context.route((u) => u.pathname === ruta, (route) => {
    const pet = route.request()
    if (pet.method() !== 'GET') return route.fallback()
    const u = new URL(pet.url()); alPedir?.(u)
    erroresTrafico.push(...erroresC1('GET', u.pathname, 200, null, json))
    return route.fulfill({ status: 200, json, headers: { 'X-Request-Id': 'prueba' } })
  })
}
/** Marca las filas, abre «Editar n registros a la vez» con Ctrl+/ y elige propiedad y valor; devuelve el bloque. */
async function editarEnGrupo(page: Page, filas: string[], propiedad: string | null, valor: string, dentro = page.locator('body')) {
  const bloque = page.getByTestId('edicion-varias')
  if (!(await bloque.isVisible())) {
    for (const f of filas) await dentro.getByRole('checkbox', { name: `Seleccionar ${f}`, exact: true }).check()
    await dentro.getByRole('checkbox', { name: `Seleccionar ${filas[0]}`, exact: true }).focus()
    await page.keyboard.press('Control+/')
  }
  await expect(bloque).toBeVisible()
  if (propiedad) await bloque.getByRole('combobox', { name: 'Propiedad' }).selectOption(propiedad)
  await bloque.getByRole('combobox', { name: 'Valor nuevo' }).selectOption(valor)
  return bloque
}
const parches = (escrituras: { metodo: string; url: string; cuerpo: unknown }[], prefijo: string) => escrituras.filter((e) => e.metodo === 'PATCH' && e.url.startsWith(prefijo))

test('P-26.02 Ctrl+/ en gastos fijos: «activo» de 2 partidas con la vista previa 3 (h) sumada; Esc no guarda y Enter aplica a las 2', async ({ page, context }) => {
  const escrituras = await preparar(context)
  const base = (DATOS['/api/v1/planned-items'] as { items: Record<string, unknown>[] }).items
  await conGet(context, '/api/v1/planned-items', pagina([...base, { id: 4, kind: 'fixed_expense', description: 'Internet', amount_cents: 50000, frequency: 'monthly', interval_count: 1, start_date: '2026-10-20', category_id: 4, accrual_months: 1, is_active: true, ...sello }]) as NonNullable<DatosC1['/api/v1/planned-items']>)
  await page.goto('/ajustes/02')
  const bloque = await editarEnGrupo(page, ['Papel', 'Internet'], 'is_active', 'no')
  const previa = bloque.getByTestId('vista-previa')
  // (c) Papel $199.00 + Internet $500.00 en el mes; (d) partes: Karen 50 % de Papel; el resto, Hogar; (f) aviso de P-07.
  await expect(previa).toContainText(/Gastos fijos de \w{3}-\d{4}\s*\$699\.00\s*\$0\.00/)
  await expect(previa).toContainText(/Karen · \w{3}-\d{4}\s*\$99\.50\s*\$0\.00/)
  await expect(previa).toContainText(/Hogar \(el resto\) · \w{3}-\d{4}\s*\$599\.50\s*\$0\.00/)
  await expect(previa).toContainText('Flujo y proyección se recalcula')
  expect(await textosEnIngles(page)).toEqual([])
  expect(await dineroSinTabulares(page)).toEqual([])
  await page.keyboard.press('Escape')
  await expect(bloque).toBeHidden()
  expect(parches(escrituras, '/api/v1/planned-items/')).toEqual([])
  await editarEnGrupo(page, ['Papel', 'Internet'], 'category_id', '5')
  await expect(previa).toContainText('1 de 2 partidas toman el reparto de su categoría')
  await bloque.getByRole('combobox', { name: 'Valor nuevo' }).press('Enter')
  await expect(bloque).toBeHidden()
  // Un updatePlannedItem por partida, en el orden de la tabla (por día de la primera fecha: Internet el 20, Papel el 30).
  expect(parches(escrituras, '/api/v1/planned-items/').map((e) => [e.url, e.cuerpo])).toEqual([['/api/v1/planned-items/4', { category_id: 5 }], ['/api/v1/planned-items/1', { category_id: 5 }]])
})

test('P-26.04 Ctrl+/ en categorías: reparto de 2 con el ejemplo de $199.00 y los gastos fijos que cambian; Esc no guarda y Enter aplica', async ({ page, context }) => {
  const escrituras = await preparar(context)
  await page.goto('/ajustes/04')
  const bloque = await editarEnGrupo(page, ['Despensa', 'Comida fuera'], 'reparto', '2')
  const previa = bloque.getByTestId('vista-previa')
  await expect(previa).toContainText(/2 categorías pasan a Conjunto 50\/50\s*0 de 2\s*2 de 2/)
  await expect(previa).toContainText(/Ejemplo con un gasto de \$199\.00 · El de por defecto\s*Hogar \$199\.00\s*Karen \$99\.50 · Hogar \$99\.50/)
  // La despensa bimestral ($3,094.50) no tiene reparto propio: toma el de su categoría (R5-2).
  await expect(previa).toContainText(/1 gasto fijo sin reparto propio cambia sus partes · Karen\s*\$0\.00\s*\$1,547\.25/)
  await expect(previa).toContainText('Aplica a los gastos que registres desde ahora')
  expect(await textosEnIngles(page)).toEqual([])
  await page.keyboard.press('Escape')
  await expect(bloque).toBeHidden()
  expect(parches(escrituras, '/api/v1/categories/')).toEqual([])
  await editarEnGrupo(page, ['Despensa', 'Comida fuera'], 'reparto', '2')
  await bloque.getByRole('combobox', { name: 'Valor nuevo' }).press('Enter')
  await expect(bloque).toBeHidden()
  expect(parches(escrituras, '/api/v1/categories/').map((e) => [e.url, e.cuerpo])).toEqual([['/api/v1/categories/1', { default_cost_share_rule_id: 2 }], ['/api/v1/categories/2', { default_cost_share_rule_id: 2 }]])
})

test('P-26.05 Ctrl+/ en promociones: desactivar 2 con «vigentes hoy» sumado y el aviso de R5-3; Esc no guarda y Enter aplica', async ({ page, context }) => {
  const escrituras = await preparar(context)
  const walmart = (DATOS['/api/v1/promotions'] as { items: Record<string, unknown>[] }).items[0]
  await conGet(context, '/api/v1/promotions', pagina([walmart, { id: 2, merchant_name: 'Liverpool', term_count: 6, rate_num: 0, rate_den: 1, min_purchase_cents: 300000, valid_from: '2026-10-01', valid_to: '2026-12-31', is_active: true, is_current: true, ...sello }]) as NonNullable<DatosC1['/api/v1/promotions']>)
  await page.goto('/ajustes/05')
  const promos = page.getByTestId('promociones')
  await promos.getByRole('button', { name: 'Tabla', exact: true }).click()
  const bloque = await editarEnGrupo(page, ['Walmart · Sin intereses', 'Liverpool · Sin intereses'], 'is_active', 'no', promos)
  const previa = bloque.getByTestId('vista-previa')
  await expect(previa).toContainText(/Activa\s*Sí\s*No/)
  await expect(previa).toContainText(/Promociones vigentes hoy\s*2 de 2\s*0 de 2/)
  await expect(previa).toContainText('cambia los escenarios del comparador')
  expect(await textosEnIngles(page)).toEqual([])
  // U5-06: en 412 táctil, casillas, listas y botones del bloque miden ≥ 44 px (con puntero fino la norma es otra).
  if (test.info().project.name === 'movil') expect((await controlesChicos(page)).chicos).toEqual([])
  await page.keyboard.press('Escape')
  await expect(bloque).toBeHidden()
  expect(parches(escrituras, '/api/v1/promotions/')).toEqual([])
  await editarEnGrupo(page, ['Walmart · Sin intereses', 'Liverpool · Sin intereses'], 'is_active', 'no', promos)
  await bloque.getByRole('combobox', { name: 'Valor nuevo' }).press('Enter')
  await expect(bloque).toBeHidden()
  expect(parches(escrituras, '/api/v1/promotions/').map((e) => [e.url, e.cuerpo])).toEqual([['/api/v1/promotions/1', { is_active: false }], ['/api/v1/promotions/2', { is_active: false }]])
})

test('P-26.07 Ctrl+/ en el recetario: tiempos aptos de 2 recetas; quitar el único tiempo no guarda; Esc no guarda y Enter aplica', async ({ page, context }) => {
  const escrituras = await preparar(context)
  const tinga = (DATOS['/api/v1/recipes'] as { items: Record<string, unknown>[] }).items[0]
  await conGet(context, '/api/v1/recipes', pagina([tinga, { id: 3, name: 'Caldo de verduras', base_servings: 2, origin: 'household', is_active: true, meal_slot_ids: [2], ...sello }]) as NonNullable<DatosC1['/api/v1/recipes']>)
  await page.goto('/ajustes/07')
  const bloque = await editarEnGrupo(page, ['Tinga de pollo', 'Caldo de verduras'], null, '-2')
  const previa = bloque.getByTestId('vista-previa')
  await expect(previa).toContainText(/Recetas aptas para Comida\s*2 de 2\s*0 de 2/)
  await expect(previa).toContainText('2 recetas se quedarían sin tiempo de comida. Elige al menos un tiempo de comida.')
  expect(await textosEnIngles(page)).toEqual([])
  await bloque.getByRole('combobox', { name: 'Valor nuevo' }).press('Enter')
  await expect(page.getByText('Elige al menos un tiempo de comida.').last()).toBeVisible()
  await expect(bloque).toBeVisible()
  expect(parches(escrituras, '/api/v1/recipes/')).toEqual([])
  await page.keyboard.press('Escape')
  await expect(bloque).toBeHidden()
  await editarEnGrupo(page, ['Tinga de pollo', 'Caldo de verduras'], null, '+1')
  await expect(previa).toContainText(/Recetas aptas para Desayuno\s*0 de 2\s*2 de 2/)
  await bloque.getByRole('combobox', { name: 'Valor nuevo' }).press('Enter')
  await expect(bloque).toBeHidden()
  expect(parches(escrituras, '/api/v1/recipes/').map((e) => [e.url, e.cuerpo])).toEqual([['/api/v1/recipes/1', { meal_slot_ids: [1, 2] }], ['/api/v1/recipes/3', { meal_slot_ids: [1, 2] }]])
})

test('P-26.05 N-284: la lectura A/B y el aviso «pendiente confirmar» salen de getCredit', async ({ page, context }) => {
  await preparar(context)
  let leida = 0
  await conGet(context, '/api/v1/credits/1', { ...(DATOS['/api/v1/credits/1'] as NonNullable<DatosC1['/api/v1/credits/1']>), availability_reading: 'A', availability_confirmed: true }, () => leida++)
  await page.goto('/ajustes/05')
  await page.getByText('Kueski Pay').first().click()
  // La lista dice B y sin confirmar; getCredit dice A y confirmada: manda getCredit.
  await expect(page.getByTestId('lectura-activa')).toHaveText('A: el límite está libre hoy y los pagos suman encima')
  await expect(page.getByText('Pendiente confirmar con el acreedor')).toHaveCount(0)
  expect(leida).toBeGreaterThan(0)
})

test('P-26.05 R2-05 vista previa 10: la compra vigente más reciente del préstamo, con su frecuencia y su primera fecha', async ({ page, context }) => {
  const escrituras = await preparar(context)
  const lineas = (DATOS['/api/v1/credits'] as { items: Record<string, unknown>[] }).items
  await conGet(context, '/api/v1/credits', pagina([...lineas, { id: 3, name: 'Préstamo de nómina', kind: 'term_loan', holder_person_id: 1, payment_frequency: 'fortnightly', annual_rate_num: 24, annual_rate_den: 100, annual_rate_decimals: 2, annual_rate_pct: '24.00', availability_reading: 'B', availability_confirmed: true, is_active: true, ...sello }]) as NonNullable<DatosC1['/api/v1/credits']>)
  const consultas: URLSearchParams[] = []
  await conGet(context, '/api/v1/credit-purchases', pagina([{ id: 9, credit_id: 3, kind: 'installment_plan', purchase_date: '2026-10-05', principal_cents: 1200000, interest_cents: 0, total_cents: 1200000, first_payment_date: '2026-11-05', term_count: 12, payment_frequency: 'monthly', interest_source: 'quoted', rate_num: 24, rate_den: 100, rate_basis: 'annual', interest_includes_vat: true, status: 'active', ...sello }]) as NonNullable<DatosC1['/api/v1/credit-purchases']>, (u) => consultas.push(u.searchParams))
  await page.goto('/ajustes/05')
  await page.getByText('Préstamo de nómina').first().click()
  await page.getByRole('button', { name: /Editar Tasa anual, como se capturó/ }).click()
  await page.getByRole('textbox', { name: 'Tasa anual, como se capturó' }).fill('30')
  const previa = page.getByTestId('vista-previa')
  await expect(previa).toContainText('Compra de ejemplo: $12,000.00 a 12 pagos, del 5-nov-2026')
  await expect(previa).toContainText('Próximo pago')
  const q = consultas.find((c) => c.get('status'))
  expect(q?.get('status')).toBe('active')
  expect(q?.get('sort')).toBe('-purchase_date,-id')
  expect(q?.get('limit')).toBe('1')
  await page.keyboard.press('Escape')
  expect(parches(escrituras, '/api/v1/credits/')).toEqual([])
})

test('P-26.06 R2-07: el punto de equilibrio de la membresía dice PROVISIONAL y cita su bloqueo', async ({ page, context }) => {
  await preparar(context)
  const tiendas = (DATOS['/api/v1/stores'] as { items: Record<string, unknown>[] }).items.map((x) => (x.name === "Sam's Club" ? { ...x, membership_annual_fee_cents: 50000 } : x))
  await conGet(context, '/api/v1/stores', pagina(tiendas) as NonNullable<DatosC1['/api/v1/stores']>)
  await page.goto('/ajustes/06')
  await page.getByRole('button', { name: "Sam's Club" }).first().click()
  await expect(page.getByTestId('equilibrio')).toHaveText('Vista previa pendiente · PROVISIONAL (bloqueo WF3-ajustes-a §2)')
})

// ——— áreas R5-01 (WF-3b, frontend-min-fix): U4 §0.3 «Error al guardar» en los manejadores sin `catch` y V-18 al eliminar una promoción en uso.
// Cada prueba simula el fallo con page.route (gana sobre las rutas del contexto de preparar()) y comprueba el texto y 0 pageerror.
const NO_GUARDO = 'No se pudo guardar. Revisa la conexión con el servidor de casa'
const ERROR_500 = { code: 'internal_error', message: 'Error interno del servidor.', request_id: 'prueba' } satisfies Esquemas['Error']
const fallosDe = (page: Page) => { const f: string[] = []; page.on('pageerror', (e) => f.push(e.message)); return f }

test('áreas R5-01 · P-26.05 «Pesos iguales» con 500: avisa «No se pudo guardar…» y Reintentar guarda', async ({ page, context }) => {
  await preparar(context)
  const fallosPagina = fallosDe(page)
  const parches: unknown[] = []
  let falla = true
  // El fake del área no responde PATCH /settings/optimizer: el éxito del reintento se responde aquí, validado contra C1.
  await page.route('**/api/v1/settings/optimizer', (route) => {
    const pet = route.request(), u = new URL(pet.url())
    if (pet.method() !== 'PATCH') return route.fallback()
    const cuerpo = pet.postDataJSON() as Cuerpo<'/api/v1/settings/optimizer', 'patch'>
    parches.push(cuerpo)
    if (falla) return route.fulfill({ status: 500, json: ERROR_500, headers: { 'X-Request-Id': 'prueba' } })
    const json = { ...(DATOS['/api/v1/settings/optimizer'] as Esquemas['OptimizerSettings']), weights: cuerpo.weights! }
    erroresTrafico.push(...erroresC1('PATCH', u.pathname, 200, cuerpo, json))
    return route.fulfill({ status: 200, json, headers: { 'X-Request-Id': 'prueba' } })
  })
  await page.goto('/ajustes/05')
  await page.getByRole('button', { name: 'Pesos iguales' }).click()
  const aviso = page.getByRole('alert').filter({ hasText: NO_GUARDO })
  await expect(aviso).toBeVisible()
  falla = false
  await aviso.getByRole('button', { name: 'Reintentar' }).click()
  await expect(aviso).toBeHidden()
  const iguales = { weights: { lowest_interest: 1, even_low_payment: 1, max_available: 1, most_value_per_peso: 1 } }
  await expect.poll(() => parches).toEqual([iguales, iguales])
  expect(fallosPagina).toEqual([])
})

test('áreas R5-01 · P-26.05 alta de promoción con 500: la hoja queda abierta con lo capturado y Reintentar usa la misma clave', async ({ page, context }) => {
  const escrituras = await preparar(context)
  const fallosPagina = fallosDe(page)
  const claves: (string | undefined)[] = []
  let falla = true
  await page.route('**/api/v1/promotions', (route) => {
    if (route.request().method() !== 'POST') return route.fallback()
    claves.push(route.request().headers()['idempotency-key'])
    return falla ? route.fulfill({ status: 500, json: ERROR_500, headers: { 'X-Request-Id': 'prueba' } }) : route.fallback()
  })
  await page.goto('/ajustes/05')
  await page.getByRole('button', { name: 'Nueva promoción' }).first().click()
  const hoja = page.getByRole('dialog', { name: 'Nueva promoción' })
  await hoja.getByLabel('Comercio').fill('Chedraui')
  await hoja.getByLabel('Plazo (pagos)').fill('4')
  await hoja.getByLabel('Vigencia: inicio').fill('2026-10-30')
  await hoja.getByLabel('Vigencia: fin').fill('2026-10-30')
  await hoja.getByLabel('Vigencia: fin').press('Control+Enter')
  const aviso = hoja.getByRole('alert').filter({ hasText: NO_GUARDO })
  await expect(aviso).toBeVisible()
  await expect(hoja.getByLabel('Comercio')).toHaveValue('Chedraui')
  await expect(hoja.getByLabel('Plazo (pagos)')).toHaveValue('4')
  falla = false
  await aviso.getByRole('button', { name: 'Reintentar' }).click()
  await expect(hoja).toBeHidden()
  await expect.poll(() => escrituras.at(-1)?.url).toBe('/api/v1/promotions')
  expect(escrituras.at(-1)?.cuerpo).toMatchObject({ merchant_name: 'Chedraui', term_count: 4, valid_from: '2026-10-30', valid_to: '2026-10-30', is_active: true })
  expect(claves).toHaveLength(2)
  expect(claves[0]).toBeTruthy()
  expect(claves[1]).toBe(claves[0])
  expect(fallosPagina).toEqual([])
})

test('áreas R5-01 · P-26.05 eliminar una promoción con 500: el diálogo crítico queda abierto con «No se pudo guardar…» y Reintentar borra', async ({ page, context }) => {
  const escrituras = await preparar(context)
  const fallosPagina = fallosDe(page)
  let falla = true
  await page.route(/\/api\/v1\/promotions\/\d+$/, (route) => (route.request().method() === 'DELETE' && falla ? route.fulfill({ status: 500, json: ERROR_500, headers: { 'X-Request-Id': 'prueba' } }) : route.fallback()))
  await page.goto('/ajustes/05')
  const promos = page.getByTestId('promociones')
  await promos.getByRole('button', { name: 'Tabla', exact: true }).click()
  await promos.getByRole('row').filter({ hasText: 'Walmart' }).getByRole('button', { name: 'Eliminar' }).click()
  const dialogo = page.getByRole('alertdialog')
  await dialogo.getByRole('button', { name: 'Eliminar la promoción de Walmart' }).click()
  const aviso = dialogo.getByRole('alert').filter({ hasText: NO_GUARDO })
  await expect(aviso).toBeVisible()
  await expect(dialogo).toBeVisible()
  falla = false
  await aviso.getByRole('button', { name: 'Reintentar' }).click()
  await expect(dialogo).toBeHidden()
  await expect.poll(() => escrituras.filter((e) => e.metodo === 'DELETE').map((e) => e.url)).toEqual(['/api/v1/promotions/1'])
  expect(fallosPagina).toEqual([])
})

test('áreas R5-01 · P-26.05 eliminar una promoción en uso (409 in_use): V-18 y «Desactivar» la apaga con is_active = false', async ({ page, context }) => {
  const escrituras = await preparar(context)
  const fallosPagina = fallosDe(page)
  await page.route(/\/api\/v1\/promotions\/\d+$/, (route) => {
    const pet = route.request(), u = new URL(pet.url())
    if (pet.method() !== 'DELETE') return route.fallback()
    const json = { code: 'in_use', message: 'La promoción tiene compras que la usan.', request_id: 'prueba' } satisfies Esquemas['Error']
    erroresTrafico.push(...erroresC1('DELETE', u.pathname, 409, null, json))
    return route.fulfill({ status: 409, json, headers: { 'X-Request-Id': 'prueba' } })
  })
  await page.goto('/ajustes/05')
  const promos = page.getByTestId('promociones')
  await promos.getByRole('button', { name: 'Tabla', exact: true }).click()
  await promos.getByRole('row').filter({ hasText: 'Walmart' }).getByRole('button', { name: 'Eliminar' }).click()
  const dialogo = page.getByRole('alertdialog')
  await dialogo.getByRole('button', { name: 'Eliminar la promoción de Walmart' }).click()
  await expect(dialogo.getByRole('alert')).toHaveText('No se puede borrar porque hay registros que la usan. Puedes desactivarla.')
  await expect(dialogo.getByRole('button', { name: 'Eliminar la promoción de Walmart' })).toHaveCount(0)
  await dialogo.getByRole('button', { name: 'Desactivar' }).click()
  await expect(dialogo).toBeHidden()
  await expect.poll(() => escrituras.filter((e) => e.metodo === 'PATCH').map((e) => [e.url, e.cuerpo])).toEqual([['/api/v1/promotions/1', { is_active: false }]])
  expect(fallosPagina).toEqual([])
})
