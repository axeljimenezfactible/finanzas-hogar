// Una prueba por pantalla (P-17..P-21, P-23..P-25) en 412×915 táctil (movil) y 1440×900 (escritorio), con la API falsa.
// Comprueba sus elementos principales, que no haya desplazamiento horizontal y, en móvil, objetivos táctiles ≥ 44 px.
import { expect, test, type Locator, type Page } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { controlesChicos, dineroSinTabulares, textosEnIngles } from '../apoyo/medidas'
import { KUESKI, instalarArea } from './datos'

async function abrir(page: Page, ruta: string) {
  await new ApiFalsa().instalar(page.context())
  const escrituras = await instalarArea(page.context())
  await page.goto(ruta)
  return escrituras
}
async function revisar(page: Page, accion: Locator) {
  await expect(accion).toBeVisible()
  expect(await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth)).toBeLessThanOrEqual(0)
  if (test.info().project.name === 'movil') expect((await accion.boundingBox())!.height).toBeGreaterThanOrEqual(44)
}

test('P-17 Todas las líneas: totales, aviso A/B, tarjetas de línea y «Falta capturar»', async ({ page }) => {
  await abrir(page, '/creditos')
  await expect(page.getByText('Deuda viva total')).toBeVisible()
  await expect(page.getByText('$3,895.00').first()).toBeVisible()
  await expect(page.getByText(/Lectura B activa · Kueski Pay · Reposición de capital pendiente de confirmar con Kueski/)).toBeVisible()
  await expect(page.getByLabel('Líneas', { exact: true }).getByRole('link', { name: 'Kueski Pay' })).toBeVisible()
  await expect(page.getByTestId('tarjeta-linea')).toHaveCount(2)
  await expect(page.getByText('Falta capturar').first()).toBeVisible()
  await revisar(page, page.getByRole('button', { name: 'Registrar compra a crédito' }))
})

test('P-18 Crédito: revolvente con lecturas A/B y tasas por plazo; tarjeta con su estado', async ({ page }) => {
  await abrir(page, '/creditos/1')
  await expect(page.getByText('Disponible · lectura B (activa)')).toBeVisible()
  await expect(page.getByText('Mínimo A $8,094.50 · mínimo B $5,000.00')).toBeVisible()
  await expect(page.getByText('17.27 %')).toBeVisible()
  await expect(page.getByText('25.8685 %')).toBeVisible()
  await revisar(page, page.getByRole('button', { name: 'Marcar pago' }))
  await page.goto('/creditos/2')
  await expect(page.getByRole('heading', { name: 'Estado de la tarjeta' })).toBeVisible()
  await expect(page.getByText('Falta capturar').first()).toBeVisible()
})

test('P-17 N-193: el aviso A/B sale de getCredit (lectura activa y reposición sin confirmar)', async ({ page }) => {
  const ctx = page.context()
  await new ApiFalsa().instalar(ctx)
  await instalarArea(ctx)
  // listCredits trae lectura B; getCredit trae lectura A: el aviso debe mostrar lo de getCredit (C1 x-necesidades N-193).
  let confirmada = false
  await ctx.route(/\/api\/v1\/credits\/1$/, (route) => route.fulfill({ status: 200, contentType: 'application/json', headers: { 'X-Request-Id': 'prueba' }, body: JSON.stringify({ ...KUESKI, availability_reading: 'A', availability_confirmed: confirmada }) }))
  const pedido = page.waitForRequest((r) => r.method() === 'GET' && new URL(r.url()).pathname === '/api/v1/credits/1')
  await page.goto('/creditos')
  await pedido
  await expect(page.getByText(/Lectura A activa · Kueski Pay · Reposición de capital pendiente de confirmar con Kueski/)).toBeVisible()
  await expect(page.getByText(/Lectura B activa/)).toHaveCount(0)
  expect(await textosEnIngles(page)).toEqual([])
  await revisar(page, page.getByRole('link', { name: 'Confirmar en Ajustes › Créditos y condiciones' }))
  confirmada = true
  await page.reload()
  await expect(page.getByTestId('tarjeta-linea')).toHaveCount(2)
  await expect(page.getByTestId('aviso-ab')).toHaveCount(0)
})

test('P-18 N-208: el calendario de la línea lee listCalendarEvents filtrado por la línea', async ({ page }) => {
  const ctx = page.context()
  await new ApiFalsa().instalar(ctx)
  await instalarArea(ctx)
  const pedido = page.waitForRequest((r) => new URL(r.url()).pathname === '/api/v1/calendar/events')
  await page.goto('/creditos/1?sub=calendario')
  const q = new URL((await pedido).url()).searchParams
  expect(q.get('credit_id')).toBe('1')
  expect(q.get('event_type')).toBe('credit_payment,promotion,card_statement')
  expect(q.get('period')).toMatch(/^\d{4}-\d{2}-01\.\.\d{4}-\d{2}-\d{2}$/)
  const primera = page.getByTestId('agenda-linea').getByTestId('pago-agenda').first()
  await expect(primera).toContainText('Vie 30 oct 2026 · Pago de Kueski Pay')
  await expect(primera).toContainText('$973.75')
  await expect(primera).toContainText('Hogar $598.82 · Karen $374.93 · pendiente')
  await expect(page.getByText(/Promoción: Walmart a 4 pagos · Vigente del/)).toBeVisible()
  await expect(page.getByText('Corte de Tarjeta Oro')).toHaveCount(0)
  await expect(page.getByRole('heading', { name: 'Calendario de la línea' })).toBeVisible({ visible: test.info().project.name !== 'movil' })
  await revisar(page, primera)
  // U5-06 (44 px en táctil), U5-07 (sin inglés) y cifras tabulares en la pestaña nueva.
  if (test.info().project.name === 'movil') expect((await controlesChicos(page)).chicos).toEqual([])
  expect(await textosEnIngles(page)).toEqual([])
  expect(await dineroSinTabulares(page)).toEqual([])
  await primera.click()
  await expect(page.getByText('Pago programado').first()).toBeVisible()
  await page.goto('/creditos/2?sub=calendario')
  await expect(page.getByTestId('agenda-linea')).toContainText('Corte de Tarjeta Oro')
  await expect(page.getByTestId('agenda-linea')).toContainText('Saldo al corte')
  await expect(page.getByText('Pago de Kueski Pay')).toHaveCount(0)
})

test('P-17 y P-18 tendencias: «Ver como tabla» ≥ 44 px táctil y ≥ 36 px con puntero fino (H-03)', async ({ page }) => {
  const minimo = test.info().project.name === 'movil' ? 44 : 36
  for (const ruta of ['/creditos?sub=tendencias', '/creditos/1?sub=tendencias']) {
    await abrir(page, ruta)
    const ver = page.getByText('Ver como tabla', { exact: true })
    await expect(ver.first()).toBeVisible()
    for (const caja of await ver.all()) expect((await caja.boundingBox())!.height).toBeGreaterThanOrEqual(minimo)
    await page.context().unrouteAll({ behavior: 'ignoreErrors' })
  }
})

test('P-19 Compra a crédito: detalle con partes y pagos; alta con vista previa', async ({ page }) => {
  await abrir(page, '/creditos/1/compras/1')
  await expect(page.getByText('Despensa Walmart').first()).toBeVisible()
  await expect(page.getByText('$3,895.00').first()).toBeVisible()
  await expect(page.getByText('Vie 30 oct 2026').first()).toBeVisible()
  await revisar(page, page.getByRole('button', { name: 'Marcar pago' }))
  await page.goto('/creditos/1/compras/nueva')
  await page.getByLabel('Capital', { exact: true }).fill('3094.50')
  await page.getByLabel('Plazo (número de pagos)').fill('4')
  await expect(page.getByText('$800.50')).toBeVisible()
  await expect(page.getByRole('button', { name: 'Guardar compra' })).not.toHaveAttribute('aria-disabled', 'true')
})

test('P-20 Comparador: recomendación, 4 óptimos y «Usar escenario» abre P-19 prellenada', async ({ page }) => {
  await abrir(page, '/creditos/comparador')
  await page.getByLabel('Monto').fill('3094.50')
  await revisar(page, page.getByRole('button', { name: 'Comparar escenarios' }))
  // N-225: restricciones duras y pesos de Ajustes precargados con «De Ajustes»; tope por persona de P-26.12.
  await page.getByText('Restricciones duras (solo para esta comparación)').click()
  await expect(page.getByLabel('Colchón mínimo de disponible')).toHaveValue('500.00')
  await expect(page.getByLabel('Máximo de compras vivas a la vez')).toHaveValue('3')
  await expect(page.getByText('De Ajustes').first()).toBeVisible()
  await page.getByLabel('Máximo de compras vivas a la vez').fill('2')
  await expect(page.getByText('Solo para esta comparación; el de Ajustes no cambia.')).toBeVisible()
  await page.getByText('Pesos de los 4 objetivos').click()
  await expect(page.getByTestId('pesos-normalizados')).toContainText('25 %')
  await expect(page.getByTestId('topes-persona')).toContainText('Sin tope')
  const peticion = page.waitForRequest((q) => q.url().endsWith('/api/v1/scenario-comparisons') && q.method() === 'POST')
  await page.getByRole('button', { name: 'Comparar escenarios' }).click()
  expect((await peticion).postDataJSON()).toMatchObject({ min_available_buffer_cents: 50000, max_live_purchases: 2, weights: { lowest_interest: 1, even_low_payment: 1, max_available: 1, most_value_per_peso: 1 } })
  await expect(page.getByText('Escenario recomendado · puntaje 87.34 de 100')).toBeVisible()
  await expect(page.getByText('24 escenarios evaluados · 2 válidos · 1 descartados')).toBeVisible()
  await expect(page.getByTestId('escenario')).toHaveCount(5)
  await page.getByRole('button', { name: 'Usar escenario recomendado' }).click()
  await expect(page).toHaveURL(/\/creditos\/1\/compras\/nueva\?monto=3094\.50&fecha=2026-10-15&plazo=4/)
})

test('P-21 Salud: dos fichas, IMC, meta por confirmar, ajuste por objetivo e historial', async ({ page }) => {
  await abrir(page, '/salud')
  await expect(page.getByTestId('ficha-salud')).toHaveCount(2)
  await expect(page.getByText('26.7 · Sobrepeso')).toBeVisible()
  await expect(page.getByRole('button', { name: 'Confirmar 1,816 kcal' })).toBeVisible()
  await expect(page.getByText('El de la casa: 0 kcal').first()).toBeVisible()
  await revisar(page, page.getByRole('button', { name: 'Registrar peso' }))
  await page.goto('/salud?sub=detalle')
  await expect(page.getByText('84.6 kg').first()).toBeVisible()
})

test('P-23 Centro de reportes: próximos envíos, reportes generados y avisos enviados', async ({ page }) => {
  await abrir(page, '/reportes')
  await expect(page.getByRole('heading', { name: 'Próximos envíos' })).toBeVisible()
  await expect(page.getByText('Resumen diario').first()).toBeVisible()
  await expect(page.getByRole('button', { name: 'Mensual PDF' }).first()).toBeVisible()
  await expect(page.getByRole('button', { name: 'Reintentar envío' })).toBeVisible()
  await expect(page.getByText(/Resumen del 11-oct-2026/)).toBeVisible()
  await revisar(page, page.getByRole('button', { name: 'Generar reporte mensual' }))
})

test('P-24 Reporte mensual: último mes cerrado, secciones y PDF', async ({ page }) => {
  await abrir(page, '/reportes/mensual')
  await expect(page.getByLabel('Mes del reporte')).toHaveValue('2026-09')
  await expect(page.getByRole('heading', { name: 'Presupuesto por categoría' })).toBeVisible()
  await expect(page.getByText('$42,000.00')).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Participación de Karen' })).toBeVisible()
  await revisar(page, page.getByRole('button', { name: 'Descargar PDF' }))
})

test('P-25 Estado de cuenta de Karen: saldo, partes y generar PDF', async ({ page }) => {
  await abrir(page, '/reportes/karen')
  await expect(page.getByText('Saldo pendiente total')).toBeVisible()
  await expect(page.getByText('$1,499.72').first()).toBeVisible()
  await expect(page.getByText(/Karen paga su parte directo a Kueski/)).toBeVisible()
  await revisar(page, page.getByRole('button', { name: 'Generar PDF' }))
})
