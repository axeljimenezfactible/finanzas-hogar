// P-01 Resumen = tablero de control (U3 2.0 §2, K54): corre a 412×915 táctil y a 1440×900, en claro y en oscuro. API simulada.
// Cifras de la semilla y del oráculo (R-13, R-64, M-01; A-05, A-06); el ingreso de las pruebas de barómetro es SINTÉTICO.
import { expect, test, type Page } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'
import { ApiFalsa } from '../apoyo/api-falsa'
import { controlesChicos, dineroSinTabulares, iconosFueraDeNorma, textosEnIngles } from '../apoyo/medidas'
import { instalarInicio, rangoMeses } from './datos'
const reloj = (page: Page, iso = '2026-10-10T12:00:00-06:00') => page.clock.setFixedTime(new Date(iso))
const bloque = (page: Page, id: string) => page.locator(`[data-bloque="${id}"]`)
const ORDEN = ['K-1', 'B-1', 'K-4', 'K-2', 'K-3', 'G-3', 'G-1', 'G-6', 'G-2', 'G-5', 'G-4', 'mas']
/** Capacidad SINTÉTICA (ingreso de prueba, no del oráculo) con el porcentaje y el estado que mandaría el motor. */
const capacidad = (ratio: number, status: string, extra: Record<string, unknown> = {}) => ({ month: '2026-10', income_cents: 2000000, income_is_expected: false, extras_cents: 0, payments_cents: 59882, payments_total_cents: 97375, payment_count: 1, ratio_pct: ratio, ratio_total_pct: ratio + 2, warning_pct: 30, limit_pct: 40, scale_max_pct: 60, status, headroom_cents: 740118, has_credits: true, ...extra })

for (const tema of ['light', 'dark'] as const) {
  test(`P-01 tablero (${tema}): bloques en orden, cifras de la semilla y una petición por bloque (A-01, A-02, A-05, A-06, A-09, A-14)`, async ({ page, context }) => {
    const reg: string[] = []
    await page.emulateMedia({ colorScheme: tema })
    await new ApiFalsa().instalar(context); await instalarInicio(context, {}, reg); await reloj(page)
    await page.goto('/inicio')
    await expect(page.getByRole('heading', { level: 1, name: 'Inicio' })).toBeVisible()
    expect(await page.evaluate(() => matchMedia('(prefers-color-scheme: dark)').matches)).toBe(tema === 'dark')
    // A-01: orden del DOM = orden de apilado de 360/412 en todos los anchos.
    await expect(page.locator('[data-bloque]')).toHaveCount(ORDEN.length)
    expect(await page.locator('[data-bloque]').evaluateAll((xs) => xs.map((x) => x.getAttribute('data-bloque')))).toEqual(ORDEN)
    // A-06: K-1 en octubre de 2026 (P1): gastado $1,197.64, efectivo $598.82, «Sin datos de septiembre»; sin plan, «Capturar presupuesto».
    const k1 = bloque(page, 'K-1')
    for (const x of ['Gasto del mes', '$1,197.64', 'gastado en octubre', '— · Aún no hay presupuesto para octubre', 'Llevas $1,197.64 gastado', 'Pagado en efectivo $598.82', 'Sin datos de septiembre', 'Capturar presupuesto']) await expect(k1).toContainText(x)
    // A-05: K-3 y K-4 con la semilla el 10-oct-2026.
    const k3 = bloque(page, 'K-3')
    for (const x of ['$3,895.00', 'deuda viva total', 'Pagos de octubre: $973.75 · 1 pago', 'Hogar $598.82 · Karen $374.93', 'Falta pagar $598.82', 'Disponible hoy $6,905.50 · Lectura B', 'Pendiente confirmar']) await expect(k3).toContainText(x)
    const k4 = bloque(page, 'K-4')
    await expect(k4.getByRole('listitem')).toHaveCount(1)
    await expect(k4).toContainText('Vie 30 oct 2026 · Kueski · $973.75 · Pendiente')
    await expect(k4).toContainText('Hogar $598.82 · Karen $374.93')
    await expect(k4.getByRole('link', { name: 'Ver calendario' })).toHaveAttribute('href', /\/inicio\?sub=calendario$/)
    // Sin ingreso en la semilla: K-2, B-1 y G-6 con su vacío; G-3 sin plan.
    await expect(bloque(page, 'K-2')).toContainText('Falta tu ingreso quincenal')
    await expect(bloque(page, 'K-2')).toContainText('Capturar ingreso')
    const b1 = bloque(page, 'B-1')
    await expect(b1.getByTestId('barometro-digito')).toHaveText('—')
    await expect(b1).toContainText('Falta tu ingreso del mes para calcular la capacidad')
    // Sin dato no hay `meter` (sin `aria-valuenow` sería inválido, axe `aria-required-attr`): imagen con el texto accesible (DESIGN §8.9).
    await expect(b1.getByRole('meter')).toHaveCount(0)
    const sinDato = b1.getByRole('img', { name: 'Capacidad de endeudamiento sin dato, tope 40 %', exact: true })
    await expect(sinDato).toHaveAttribute('data-prueba', 'barometro')
    for (const a of ['aria-valuenow', 'aria-valuemin', 'aria-valuemax', 'aria-valuetext']) expect(await sinDato.getAttribute(a)).toBeNull()
    await expect(b1.getByTestId('barometro-aguja')).toHaveCount(0)
    await expect(bloque(page, 'G-6')).toContainText('Falta tu ingreso para calcular la capacidad')
    await expect(bloque(page, 'G-3')).toContainText('Aún no hay presupuesto por categoría')
    // G-1 con el gasto de octubre (periodos[0].devengado.axel) y la nota fija; G-2 con la deuda viva por línea.
    await expect(page.getByTestId('grafica-g1').locator('figcaption')).toContainText('En oct-2026: ingreso —, gasto $1,197.64 y pagos de crédito $598.82')
    await expect(page.getByTestId('grafica-g1')).toContainText('Los pagos de crédito pagan compras que ya cuentan en el gasto de su mes; no se suman al gasto')
    await expect(page.getByTestId('grafica-g2').locator('figcaption')).toContainText('Deuda viva hoy: $3,895.00; la mayor, Kueski con $3,895.00')
    await expect(page.getByTestId('grafica-g4').locator('figcaption')).toContainText('Devengado contra efectivo')
    await expect(page.getByTestId('grafica-g5').locator('figcaption')).toContainText('Deuda viva proyectada')
    // «Más del hogar» conserva las tarjetas 2, 4, 5 y 6 (§2.5).
    const mas = bloque(page, 'mas')
    await expect(mas.getByRole('heading', { level: 2, name: 'Más del hogar' })).toBeVisible()
    expect(await mas.locator('[data-prueba^="tarjeta-"]').evaluateAll((xs) => xs.map((x) => x.getAttribute('data-prueba')))).toEqual(['tarjeta-2', 'tarjeta-4', 'tarjeta-5', 'tarjeta-6'])
    await expect(mas).toContainText('Bodega Aurrera $3,094.50')
    // A-02: una petición por operación; G-1 y G-6 de una sola (from = mes − 11, to = mes + 5); ajustes leídos una vez.
    const de = (p: string) => reg.filter((r) => r.startsWith(p))
    expect(de('/api/v1/dashboard/monthly-overview')).toEqual(['/api/v1/dashboard/monthly-overview?from=2025-11&to=2027-03'])
    expect(de('/api/v1/settings/alerts')).toHaveLength(1)
    expect(de('/api/v1/settings/dashboard')).toHaveLength(1)
    expect(de('/api/v1/dashboard/debt-capacity')).toEqual(['/api/v1/dashboard/debt-capacity?period=2026-10'])
    expect(de('/api/v1/incomes/month')).toEqual(['/api/v1/incomes/month?period=2026-10'])
    expect(de('/api/v1/budget/signals')).toEqual(['/api/v1/budget/signals?period=2026-10'])
    expect(de('/api/v1/scheduled-payments').map(decodeURIComponent)).toEqual(['/api/v1/scheduled-payments?status=pending&sort=due_date&limit=5&period=2026-10-10..2026-11-09'])
    // A-08: K-1, K-2, K-3 y B-1 son un solo enlace cada uno (la tarjeta entera).
    for (const id of ['K-1', 'K-2', 'K-3', 'B-1']) {
      const t = bloque(page, id).locator(`[data-prueba="tarjeta-${id}"]`)
      expect(await t.evaluate((x) => x.tagName)).toBe('A')
      await expect(t.locator('a, button')).toHaveCount(0)
    }
    // A-09, A-10, A-14, U5-06/07: ≥ 44 px, sin desplazamiento horizontal, dígito sin recortar, tabulares, sin inglés ni fechas ISO.
    // 44 × 44 con puntero grueso (412 táctil); con puntero fino la barra lateral usa 32 px (DESIGN §7.9).
    if (await page.evaluate(() => matchMedia('(any-pointer: coarse)').matches)) expect((await controlesChicos(page)).chicos).toEqual([])
    expect((await iconosFueraDeNorma(page)).mal).toEqual([])
    expect(await dineroSinTabulares(page)).toEqual([])
    expect(await textosEnIngles(page)).toEqual([])
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true)
    const texto = await page.locator('main').innerText()
    expect(texto).not.toMatch(/\d{4}-\d{2}-\d{2}/)
    // U5-03/U5-06: axe WCAG 2.0–2.2 A/AA sobre el vacío de la semilla (B-1, K-2 y G-6 sin ingreso), en cada ancho y tema.
    const graves = (await new AxeBuilder({ page: page as never }).withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'wcag22aa']).analyze()).violations
      .filter((v) => v.impact === 'serious' || v.impact === 'critical').map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(' ')).join(' | ')}`)
    expect(graves).toEqual([])
  })
}

test('P-01 B-1: el barómetro pinta lo que manda C1 y ningún umbral es constante (A-18)', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await reloj(page)
  const casos: [number, string, string, string][] = [[30, 'normal', 'Verde · dentro del rango', 'verde'], [31, 'warning', 'Amarillo · cerca del tope', 'amarillo'], [40, 'warning', 'Amarillo · cerca del tope', 'amarillo'], [41, 'exceeded', 'Rojo · arriba del tope', 'rojo']]
  for (const [n, st, insignia, color] of casos) {
    await context.unrouteAll({ behavior: 'ignoreErrors' }); await new ApiFalsa().instalar(context)
    await instalarInicio(context, { '/api/v1/dashboard/debt-capacity': capacidad(n, st) })
    await page.goto('/inicio')
    const b1 = bloque(page, 'B-1')
    await expect(b1.getByTestId('barometro-digito')).toHaveText(`${n} %`)
    await expect(b1).toContainText(insignia)
    const medidor = b1.getByRole('meter', { name: 'Capacidad de endeudamiento' })
    await expect(medidor).toHaveAttribute('aria-valuetext', `Capacidad de endeudamiento ${n} % de 40 %: ${color}`)
    await expect(medidor).toHaveAttribute('aria-valuenow', String(n))
    await expect(b1.getByTestId('barometro-aguja')).toHaveCount(1)
    await expect(b1).toContainText('Pagos $598.82 de un ingreso de $20,000.00')
    await expect(b1).toContainText('Te caben $7,401.18 más al mes antes del rojo')
  }
  // Con warning_pct = 25 en la respuesta, la marca del tramo se mueve (los umbrales llegan de C1, K44).
  await context.unrouteAll({ behavior: 'ignoreErrors' }); await new ApiFalsa().instalar(context)
  await instalarInicio(context, { '/api/v1/dashboard/debt-capacity': capacidad(27, 'warning', { warning_pct: 25, headroom_cents: -1200 }) })
  await page.goto('/inicio')
  await expect(bloque(page, 'B-1').getByTestId('barometro-marca')).toHaveText(['25 %', '40 %'])
  await expect(bloque(page, 'B-1')).toContainText('−$12.00 arriba del tope')
  // Los tramos usan los tokens de rol (sin colores sueltos): éxito, advertencia y error distintos entre sí.
  const trazos = await bloque(page, 'B-1').locator('path[data-tramo]').evaluateAll((ps) => ps.map((p) => getComputedStyle(p).stroke))
  expect(new Set(trazos).size).toBe(3)
})

/** Respuestas SINTÉTICAS (ingreso de prueba) con las 6 gráficas del tablero pintadas: G-1, G-6, G-2, G-3, G-4 y G-5. */
const sinteticas = () => {
  const ratio: Record<string, [number, string]> = { '2026-10': [27, 'normal'], '2026-11': [35, 'warning'], '2026-12': [44, 'exceeded'] }
  return {
    '/api/v1/dashboard/debt-capacity': capacidad(27, 'normal'),
    '/api/v1/dashboard/monthly-overview': (u: URL) => ({ warning_pct: 30, limit_pct: 40, months: rangoMeses(u.searchParams.get('from')!, u.searchParams.get('to')!).map((month) => {
      const r = ratio[month], proy = month > '2026-10'
      return { month, is_projection: proy, income_is_expected: proy, ...(month >= '2026-09' ? { income_total_cents: 2000000, income_fortnights_cents: 2000000, income_extras_cents: 0 } : {}), ...(month === '2026-10' ? { accrued_cents: 119764, cash_cents: 59882 } : {}), credit_payments_cents: month === '2026-10' ? 59882 : 0, credit_payments_total_cents: 0, ...(r ? { debt_ratio_pct: r[0], debt_status: r[1] } : {}) }
    }) }),
    '/api/v1/settings/alerts': { parameters: { 'alerts.payment_due_days_before': { value: [1, 0] }, 'alerts.budget_warning_pct': { value: 80 }, 'alerts.budget_limit_pct': { value: 100 } } },
    '/api/v1/budget/categories': { month: '2026-10', rows: [
      { category_id: 1, planned_cents: 300000, actual_cents: 330000, difference_cents: -30000, used_pct: 110, status: 'exceeded', expense_count: 4 },
      { category_id: 2, planned_cents: 100000, actual_cents: 50000, difference_cents: 50000, used_pct: 50, status: 'normal', expense_count: 1 },
      { category_id: 3, planned_cents: 100000, actual_cents: 85000, difference_cents: 15000, used_pct: 85, status: 'warning', expense_count: 2 },
      { category_id: 4, actual_cents: 12000, expense_count: 1 },
    ] },
    '/api/v1/budget/signals': { month: '2026-10', months_considered: 3, margin_pct: 10, rows: [{ category_id: 1, signal: 'raise', months: [] }, { category_id: 2, signal: 'cut', months: [] }, { category_id: 3, signal: 'none', months: [] }] },
  }
}

test('P-01 tablero con ingreso SINTÉTICO: G-1, G-6 y G-3 con sugerencias y desglose a P-02 (A-19, A-20, A-21 paso 1)', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await reloj(page)
  await instalarInicio(context, sinteticas())
  await page.goto('/inicio')
  const g3 = page.getByTestId('grafica-g3')
  await expect(g3.locator('figcaption')).toContainText('En octubre: 2 categorías arriba de 80 % y 1 rebasada')
  const filas = g3.getByRole('listitem')
  await expect(filas).toHaveCount(3)
  await expect(filas.nth(0)).toContainText('Rebasada')
  await expect(filas.nth(0)).toContainText('Sugerencia: subir presupuesto')
  await expect(filas.nth(1)).toContainText('Vas en 85 %')
  await expect(filas.nth(2)).toContainText('Sugerencia: recortar presupuesto')
  await expect(filas.nth(0).locator('[aria-label$="3 meses seguidos arriba del plan por más de 10 %"]')).toHaveCount(1)
  await expect(g3).toContainText('1 categoría con gasto y sin plan')
  await expect(g3.locator('[data-marca="aviso"]')).toHaveCount(3)
  const g1 = page.getByTestId('grafica-g1')
  await expect(g1.locator('figcaption')).toContainText('En oct-2026: ingreso $20,000.00, gasto $1,197.64 y pagos de crédito $598.82')
  await g1.getByRole('button', { name: 'Ver como tabla' }).click()
  await expect(g1.getByRole('row', { name: /oct-2026/ })).toContainText('$20,000.00')
  await expect(g1.getByRole('row', { name: /ago-2026/ })).toContainText('—')
  const g6 = page.getByTestId('grafica-g6')
  await expect(g6.locator('figcaption')).toContainText('oct-2026: 27 % (verde); el mayor de la proyección es dic-2026 con 44 %')
  await expect(g6).toContainText('Proyección con pagos programados e ingreso esperado')
  for (const x of ['30 %', '40 %', 'Hoy']) await expect(g6.locator('svg').getByText(x, { exact: true }).first()).toBeAttached()
  await g6.getByRole('button', { name: 'Ver como tabla' }).click()
  await expect(g6.getByRole('row', { name: /nov-2026/ })).toContainText('35 %amarillo')
  await expect(g6.getByRole('row', { name: /jul-2026/ })).toContainText('— · Falta ingreso')
  expect(await textosEnIngles(page)).toEqual([])
  // Desglose §2.2 paso 1: la barra de la categoría abre P-02 con la categoría y el mes.
  await filas.nth(0).getByRole('link').click()
  await expect(page).toHaveURL(/\/gastos\?categoria=1&periodo=2026-10$/)
})

test('P-01 estados del tablero: un bloque con falla no tumba los demás y «Reintentar» lo recupera (A-07)', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await instalarInicio(context, { '/api/v1/dashboard/debt-capacity': 500, '/api/v1/settings/dashboard': 500, '/api/v1/dashboard/monthly-overview': 500 }); await reloj(page)
  await page.goto('/inicio')
  await expect(bloque(page, 'B-1')).toContainText('No se pudo calcular la capacidad de endeudamiento. Revisa que estés en la red de casa.')
  await expect(bloque(page, 'B-1').getByRole('link', { name: 'Ver créditos' })).toBeVisible()
  await expect(bloque(page, 'K-4')).toContainText('No se pudieron cargar los pagos. Revisa que estés en la red de casa.')
  await expect(bloque(page, 'G-1')).toContainText('No se pudo cargar la gráfica.')
  await expect(bloque(page, 'G-6')).toContainText('No se pudo cargar la gráfica.')
  await expect(bloque(page, 'K-1')).toContainText('$1,197.64')
  await expect(bloque(page, 'K-3')).toContainText('Pagos de octubre: $973.75 · 1 pago')
  await context.unrouteAll({ behavior: 'ignoreErrors' })
  await new ApiFalsa().instalar(context); await instalarInicio(context)
  await bloque(page, 'B-1').getByRole('button', { name: 'Reintentar' }).click()
  await expect(bloque(page, 'B-1')).toContainText('Falta tu ingreso del mes para calcular la capacidad')
  await bloque(page, 'K-4').getByRole('button', { name: 'Reintentar' }).click()
  await expect(bloque(page, 'K-4')).toContainText('Vie 30 oct 2026 · Kueski · $973.75 · Pendiente')
})

test('P-01 servidor anterior a C1 v1.5 (501 o 404 en las operaciones nuevas): B-1, G-1 y G-6 dicen «Disponible cuando el servidor se actualice», sin cifra ni «Reintentar» (U3 2.0 §2, U3-03)', async ({ page, context }) => {
  await reloj(page)
  for (const estado of [501, 404]) {
    await context.unrouteAll({ behavior: 'ignoreErrors' })
    await new ApiFalsa().instalar(context); await instalarInicio(context, { '/api/v1/dashboard/debt-capacity': estado, '/api/v1/dashboard/monthly-overview': estado })
    await page.goto('/inicio')
    for (const id of ['B-1', 'G-1', 'G-6']) {
      const b = bloque(page, id)
      await expect(b.getByRole('status')).toHaveText('Disponible cuando el servidor se actualice')
      await expect(b.getByRole('alert')).toHaveCount(0)
      await expect(b).not.toContainText('Revisa que estés en la red de casa')
      await expect(b.getByRole('button', { name: 'Reintentar' })).toHaveCount(0)
      expect(await b.innerText()).not.toMatch(/\$|\d\s?%/)
    }
    await expect(bloque(page, 'B-1').getByRole('link', { name: 'Ver créditos' })).toBeVisible()
    // Los bloques de operaciones que ya existen siguen con sus cifras de la semilla.
    await expect(bloque(page, 'K-1')).toContainText('$1,197.64')
    await expect(bloque(page, 'K-3')).toContainText('Pagos de octubre: $973.75 · 1 pago')
  }
})

for (const tema of ['light', 'dark'] as const) {
  test(`P-01 tablero (${tema}): axe WCAG 2.0–2.2 A/AA sin infracciones serias ni críticas con las 6 gráficas, ocultas y como tabla (U5-06, U5-03)`, async ({ page, context }) => {
    await page.emulateMedia({ colorScheme: tema })
    await new ApiFalsa().instalar(context); await reloj(page)
    await instalarInicio(context, sinteticas())
    await page.goto('/inicio')
    for (const g of ['g3', 'g1', 'g6', 'g2', 'g5', 'g4']) await expect(page.getByTestId(`grafica-${g}`)).toBeVisible()
    const graves = async () => (await new AxeBuilder({ page: page as never }).withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'wcag22aa']).analyze()).violations
      .filter((v) => v.impact === 'serious' || v.impact === 'critical').map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(' ')).join(' | ')}`)
    expect(await graves()).toEqual([])
    // La tabla oculta no es región desplazable; ya visible, es región con foco y nombre = título (FiguraBarras y FiguraGrafica).
    for (const [g, titulo] of [['g2', 'Deuda viva por línea'], ['g1', 'Ingreso, gasto y pagos por mes']] as const) {
      const f = page.getByTestId(`grafica-${g}`)
      await expect(f.getByRole('region')).toHaveCount(0)
      await f.getByRole('button', { name: 'Ver como tabla' }).click()
      const region = f.getByRole('region', { name: titulo })
      await expect(region).toHaveAttribute('tabindex', '0')
      await region.focus(); await expect(region).toBeFocused()
    }
    expect(await graves()).toEqual([])
  })
}

test('P-01 K-1 en noviembre de 2026 (P2): «Igual que octubre» en lo gastado y «$598.82 más que octubre» en efectivo (A-06)', async ({ page, context }) => {
  await new ApiFalsa().instalar(context)
  await instalarInicio(context, { '/api/v1/dashboard/month-vs-plan': { month: '2026-11', actual_cents: 119764, days_left: 20, cash_cents: 119764, previous_month: '2026-10', previous_actual_cents: 119764, previous_cash_cents: 59882, actual_difference_cents: 0, cash_difference_cents: 59882, categories_warning_count: 0, categories_exceeded_count: 0 } })
  await reloj(page, '2026-11-10T12:00:00-06:00')
  await page.goto('/inicio')
  const k1 = bloque(page, 'K-1')
  await expect(k1).toContainText('gastado en noviembre')
  await expect(k1).toContainText('Gastado · Igual que octubre')
  await expect(k1).toContainText('Efectivo · $598.82 más que octubre')
})
