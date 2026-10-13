// P-01 Inicio (U3, U5-06/07/09, C3-03): corre a 412×915 táctil y a 1440×900. API simulada con context.route.
import { expect, test, type Locator, type Page } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { instalarInicio } from './datos'
const reloj = (page: Page, iso = '2026-10-10T12:00:00-06:00') => page.clock.setFixedTime(new Date(iso))
const tarjeta = (page: Page, n: number | string) => page.locator(`[data-prueba="tarjeta-${n}"]`)
// Desde U3 2.0 (K54) el Resumen es el tablero (pruebas en tablero.escritorio.spec.ts); las tarjetas 2, 4, 5 y 6 siguen en «Más del hogar».

test('P-01 Más del hogar: tarjetas 2, 4, 5 y 6 en orden, cifras de la semilla y una petición por tarjeta', async ({ page, context }) => {
  const reg: string[] = []
  await new ApiFalsa().instalar(context); await instalarInicio(context, {}, reg); await reloj(page)
  await page.goto('/inicio')
  await expect(page.getByRole('heading', { level: 1, name: 'Inicio' })).toBeVisible()
  const tarjetas = page.getByTestId('mas-del-hogar').locator('[data-prueba^="tarjeta-"]')
  await expect(tarjetas).toHaveCount(4)
  expect(await tarjetas.evaluateAll((t) => t.map((x) => x.getAttribute('data-prueba')))).toEqual([2, 4, 5, 6].map((n) => `tarjeta-${n}`))
  await expect(tarjeta(page, 2)).toContainText('Del 1 al 15 oct')
  await expect(tarjeta(page, 2)).toContainText('Ingreso esperado')
  await expect(tarjeta(page, 4)).toContainText('Aún no hay plan para esta semana')
  for (const txt of ['Bodega Aurrera $3,094.50', '$138.50 menos que Walmart', 'Tu lista está vacía; se compara la compra del 10-oct-2026']) await expect(tarjeta(page, 5)).toContainText(txt)
  await expect(tarjeta(page, 6)).toContainText('Aún no tienes fondo ni metas de ahorro')
  // A-02, A-03, A-04: semana en curso para comidas y lista vigente (K38) para la canasta; getAlertSettings una sola vez.
  expect(reg).toContain('/api/v1/dashboard/week-meals?period=2026-W41')
  expect(reg).toContain('/api/v1/dashboard/best-basket?week=2026-W41')
  expect(reg.filter((r) => r.startsWith('/api/v1/settings/alerts'))).toHaveLength(1)
  // A-08 y A-09: cada tarjeta es un solo enlace ≥ 44 px; sin desplazamiento horizontal.
  for (const n of [2, 4, 5, 6]) {
    const t = tarjeta(page, n)
    expect(await t.evaluate((x) => x.tagName)).toBe('A')
    await expect(t.locator('a, button')).toHaveCount(0)
    expect((await t.boundingBox())!.height).toBeGreaterThanOrEqual(44)
  }
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true)
  const texto = await page.locator('main').innerText()
  expect(texto).not.toMatch(/\b(Save|Cancel|Settings|Dashboard|Loading|Error|Home|Delete|Search|Add)\b/)
  expect(texto).not.toMatch(/\d{4}-\d{2}-\d{2}/)
})

test('P-01 estados: una tarjeta con falla no tumba las demás y «Reintentar» la recupera', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await instalarInicio(context, { '/api/v1/dashboard/month-vs-plan': 500 }); await reloj(page)
  await page.goto('/inicio')
  await expect(tarjeta(page, 'K-1')).toContainText('No se pudo cargar el gasto del mes. Revisa que estés en la red de casa.')
  await expect(tarjeta(page, 'K-1').getByRole('link', { name: 'Ver presupuesto' })).toBeVisible()
  await expect(tarjeta(page, 'K-3')).toContainText('Pagos de octubre: $973.75 · 1 pago')
  await context.unrouteAll({ behavior: 'ignoreErrors' })
  await new ApiFalsa().instalar(context); await instalarInicio(context)
  await tarjeta(page, 'K-1').getByRole('button', { name: 'Reintentar' }).click()
  await expect(tarjeta(page, 'K-1')).toContainText('Llevas $1,197.64 gastado')
})

test('P-01 A-17: el Jue 29 oct 2026 el pago de K-4 dice «Vence mañana»', async ({ page, context }) => {
  await new ApiFalsa().instalar(context)
  await instalarInicio(context, { '/api/v1/dashboard/next-payment': { next_payment: { due_date: '2026-10-30', credit_id: 1, amount_cents: 97375, household_share_cents: 59882, external_share_cents: 37493, days_left: 1, payment_count: 1, is_overdue: false }, availability: [] } })
  await reloj(page, '2026-10-29T12:00:00-06:00')
  await page.goto('/inicio')
  await expect(tarjeta(page, 'K-4')).toContainText('Vie 30 oct 2026 · Kueski · $973.75 · Vence mañana')
})

test('P-01 Detalle: D-3 con tres pagos y sin pedir vencidos (A-16), lecturas A y B, y participación de Karen', async ({ page, context }) => {
  const reg: string[] = []
  await new ApiFalsa().instalar(context); await instalarInicio(context, {}, reg); await reloj(page)
  await page.goto('/inicio?sub=detalle')
  const d3 = page.getByRole('region', { name: 'Pagos y disponible' })
  for (const f of ['Vie 30 oct 2026 · Kueski', 'Dom 15 nov 2026 · Kueski', 'Lun 30 nov 2026 · Kueski', 'Lectura A $10,000.00', 'Lectura B $6,905.50', 'Mínimo del horizonte $8,094.50', 'Mínimo del horizonte $5,000.00']) await expect(d3).toContainText(f)
  await expect(d3.getByText('Pendiente', { exact: true })).toHaveCount(3)
  expect(reg.some((r) => r.includes('status=overdue'))).toBe(false)
  expect(reg.some((r) => r.startsWith('/api/v1/scheduled-payments') && decodeURIComponent(r).includes('period=2026-10-10..2026-12-09'))).toBe(true)
  const d7 = page.getByRole('region', { name: 'Participación de Karen' })
  await expect(d7).toContainText('Su parte de la despensa este mes: $749.86')
  await expect(d7).toContainText('Karen paga directo a Kueski su parte; su parte es informativa y repone la línea cuando paga')
  await expect(d7.getByRole('link', { name: 'Estado de cuenta de Karen' })).toBeVisible()
})

const d3 = (page: Page) => page.getByRole('region', { name: 'Pagos y disponible' })
const PARTES = 'Hogar $598.82 · Karen $374.93'

test('P-01 Detalle: D-3 no pinta ninguna fila mientras listPersons no responde; nunca «Hogar $973.75» sin la parte de Karen (M-1 de REVISION-frontend-cierre-r2, C3-09, R-13)', async ({ page, context }) => {
  const reg: string[] = []
  // Registra si en algún instante la página dijo «Hogar $973.75» (el pago entero cargado al hogar porque aún no se sabe quién es Karen).
  await page.addInitScript(() => {
    const w = window as unknown as { __hogarTodo: boolean }; w.__hogarTodo = false
    new MutationObserver(() => { if (document.body?.textContent?.includes('Hogar $973.75')) w.__hogarTodo = true }).observe(document, { subtree: true, childList: true, characterData: true })
  })
  await new ApiFalsa().instalar(context); await instalarInicio(context, {}, reg); await reloj(page)
  let soltar = () => {}
  const personas = new Promise<void>((r) => { soltar = r })
  await context.route(/\/api\/v1\/persons(\?|$)/, async (route) => { await personas; await route.fallback() })
  await page.goto('/inicio?sub=detalle')
  // Los pagos próximos y el disponible ya respondieron y listPersons sigue retenida: esqueleto, sin filas ni cifras repartidas.
  await expect.poll(() => reg.some((x) => x.startsWith('/api/v1/scheduled-payments') && x.includes('status=pending'))).toBe(true)
  await expect(d3(page)).toContainText('Lectura B $6,905.50')
  await page.waitForTimeout(400)
  await expect(d3(page).getByRole('listitem')).toHaveCount(0)
  await expect(d3(page)).not.toContainText('$973.75')
  await expect(d3(page)).not.toContainText('Hogar')
  soltar()
  // Al soltarla, las tres filas con las dos partes separadas (oráculo R-13).
  for (const f of ['Vie 30 oct 2026 · Kueski', 'Dom 15 nov 2026 · Kueski', 'Lun 30 nov 2026 · Kueski']) await expect(d3(page)).toContainText(f)
  await expect(d3(page).getByText(PARTES, { exact: true })).toHaveCount(3)
  expect(await page.evaluate(() => (window as unknown as { __hogarTodo: boolean }).__hogarTodo)).toBe(false)
})

test('P-01 Detalle: D-3 con listPersons en falla muestra su estado de error con «Reintentar» y sin cifras repartidas; «Reintentar» lo recupera (PROVISIONAL, bloqueo U3-tablero-k4-personas)', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await instalarInicio(context, { '/api/v1/persons': 500 }); await reloj(page)
  await page.goto('/inicio?sub=detalle')
  await expect(d3(page)).toContainText('No se pudieron cargar los pagos. Revisa que estés en la red de casa.')
  await expect(d3(page)).toContainText('Lectura B $6,905.50') // el disponible no depende de listPersons
  await expect(d3(page)).not.toContainText('$973.75')
  await expect(d3(page)).not.toContainText('Hogar')
  await context.unrouteAll({ behavior: 'ignoreErrors' })
  await new ApiFalsa().instalar(context); await instalarInicio(context)
  await d3(page).getByRole('button', { name: 'Reintentar' }).click()
  await expect(d3(page)).toContainText('Vie 30 oct 2026 · Kueski')
  await expect(d3(page).getByText(PARTES, { exact: true })).toHaveCount(3)
})

test('P-01 Detalle: en D-3 una parte de una persona que no está en listPersons no se suma al hogar (M-1 de REVISION-frontend-cierre-r2)', async ({ page, context }) => {
  const parte = (person_id: number, amount_cents: number) => ({ person_id, amount_cents, paid_cents: 0, status: 'pending' })
  const pago = { id: 1, credit_purchase_id: 1, seq: 1, due_date: '2026-10-30', amount_cents: 97375, principal_cents: 90000, interest_cents: 7375, created_at: '2026-10-01T09:00:00-06:00', updated_at: '2026-10-01T09:00:00-06:00', credit_id: 1, status: 'pending', shares: [parte(1, 59882), parte(9, 37493)] }
  await new ApiFalsa().instalar(context); await instalarInicio(context, { '/api/v1/scheduled-payments': { items: [pago], next_cursor: null } }); await reloj(page)
  await page.goto('/inicio?sub=detalle')
  await expect(d3(page)).toContainText('Vie 30 oct 2026 · Kueski')
  await expect(d3(page).getByRole('listitem')).toHaveCount(1)
  await expect(d3(page)).toContainText('Pendiente')
  await expect(d3(page)).not.toContainText('Hogar')
})

test('P-01 Tendencias: T-2 con R-45 y R-46 en la tabla, figuras con figcaption y la pestaña recordada', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await instalarInicio(context); await reloj(page)
  await page.goto('/inicio?sub=tendencias')
  const t2 = page.getByTestId('grafica-t2')
  await expect(t2.locator('figcaption')).toContainText('Devengado contra efectivo')
  await expect(t2.locator('figcaption')).toContainText('devengado acumulado $3,895.00 y efectivo acumulado $2,921.25')
  await t2.getByRole('button', { name: 'Ver como tabla' }).click()
  await expect(t2.getByRole('row', { name: /oct-2026/ })).toContainText('$1,947.50$973.75$973.75')
  await expect(t2.getByRole('row', { name: /nov-2026/ })).toContainText('$1,947.50$1,947.50$973.75')
  expect((await t2.innerText()).toLowerCase()).not.toContain('deuda')
  await expect(page.getByTestId('grafica-t3').locator('figcaption')).toContainText('Deuda viva al cierre de oct-2026: $3,895.00')
  await expect(page.getByText('Aún no hay presupuesto capturado')).toBeVisible()
  await expect(page.getByText('Aún no hay semanas con plan de comidas')).toBeVisible()
  await page.goto('/inicio')
  await expect(page).toHaveURL(/sub=tendencias/)
})

test('P-01 Calendario: pagos con su parte, quincenas K8 y vista según el ancho', async ({ page, context }) => {
  await new ApiFalsa().instalar(context); await instalarInicio(context); await reloj(page)
  await page.goto('/inicio?sub=calendario')
  if (page.viewportSize()!.width < 768) {
    const ag = page.getByTestId('agenda')
    await expect(ag).toContainText('Vie 30 oct 2026')
    await expect(ag).toContainText('Kueski · $973.75')
    await expect(ag).toContainText('Hogar $598.82 · Karen $374.93')
    await expect(ag).toContainText('2.ª quincena · del 16 al 31')
  } else {
    const mes = page.getByTestId('mes')
    await expect(mes).toContainText('octubre 2026')
    await expect(mes).toContainText('1.ª quincena · del 1 al 15')
    await mes.getByRole('button', { name: /Vie 30 oct 2026/ }).click()
    await expect(page.getByRole('region', { name: 'Día seleccionado' })).toContainText('Hogar $598.82 · Karen $374.93')
  }
})

test('F-1 desde Inicio: 3 toques y Inicio vuelve a pedir el mes y la quincena; F-3 y F-4 paso 1', async ({ page, context }) => {
  const reg: string[] = []
  await new ApiFalsa().instalar(context); await instalarInicio(context, {}, reg); await reloj(page)
  await page.goto('/inicio')
  await expect(tarjeta(page, 'K-4')).toContainText('$973.75')
  const antes = reg.filter((r) => r.startsWith('/api/v1/dashboard/month-vs-plan')).length
  let toques = 0
  const tocar = async (l: Locator) => { if (page.viewportSize()!.width < 768) await l.tap(); else await l.click(); toques++ }
  await tocar(page.getByRole('button', { name: 'Registrar gasto' }).last())
  await page.getByLabel('Monto').fill('185.50')
  await tocar(page.getByRole('button', { name: 'Despensa' }))
  await tocar(page.getByRole('button', { name: 'Guardar' }))
  console.log(`U5-09 toques desde Inicio (${page.viewportSize()!.width} px): ${toques}`)
  expect(toques).toBeLessThanOrEqual(3)
  await expect.poll(() => reg.filter((r) => r.startsWith('/api/v1/dashboard/month-vs-plan')).length).toBeGreaterThan(antes)
  await expect.poll(() => reg.filter((r) => r.startsWith('/api/v1/dashboard/fortnight-flow')).length).toBeGreaterThan(1)
  await tarjeta(page, 5).click()
  await expect(page).toHaveURL(/\/compras\/comparar$/)
  await page.goBack()
  await tarjeta(page, 'K-4').getByRole('link', { name: /Vie 30 oct 2026/ }).click()
  await expect(page).toHaveURL(/\/creditos\/1$/)
})
