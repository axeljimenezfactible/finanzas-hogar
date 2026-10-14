// Flujos de U1 §4 que pasan por las pantallas del área (412×915 táctil): F-4 (crédito y calendario, marcar pago),
// F-5 (ver, descargar y generar el reporte; estado de cuenta de Karen) y la variante de F-6 desde Salud.
import { expect, test } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { instalarArea } from './datos'

test('F-4 · Crédito → Calendario → pago → marcar la parte del hogar (y deshacer)', async ({ page }) => {
  await new ApiFalsa().instalar(page.context())
  const esc = await instalarArea(page.context())
  await page.goto('/creditos/1')
  await page.getByRole('tablist', { name: 'Crédito' }).getByRole('tab', { name: 'Calendario' }).click()
  const fila = page.getByTestId('pago-agenda').first()
  await expect(fila).toContainText('Vie 30 oct 2026')
  await expect(fila).toContainText('$973.75')
  await expect(fila).toContainText('Hogar $598.82 · Karen $374.93 · pendiente')
  await fila.click()
  const hoja = page.getByRole('dialog', { name: 'Pago programado' })
  await expect(hoja.getByText('Capital que se repone')).toBeVisible()
  await hoja.getByRole('button', { name: 'Marcar pagada la parte de Hogar' }).click()
  await expect(page.getByText('Parte marcada como pagada')).toBeVisible()
  const pago = esc.find((e) => e.ruta === '/api/v1/scheduled-payments/101/payments')!
  expect(pago.cuerpo).toMatchObject({ payer_person_id: 1, amount_cents: 59882 })
  expect(pago.clave!.length).toBeGreaterThanOrEqual(16)
  await page.getByRole('button', { name: 'Deshacer' }).click()
  await expect.poll(() => esc.some((e) => e.metodo === 'DELETE' && e.ruta === '/api/v1/credit-payments/900')).toBe(true)
})

test('F-5 · Reportes → mensual sep-2026 → Descargar PDF; mes sin PDF → Generar PDF', async ({ page }) => {
  await new ApiFalsa().instalar(page.context())
  const esc = await instalarArea(page.context())
  await page.goto('/reportes')
  await page.getByRole('button', { name: 'Mensual PDF' }).first().click()
  await expect(page).toHaveURL(/\/reportes\/mensual\?mes=2026-09/)
  const descarga = page.waitForEvent('download')
  await page.getByRole('button', { name: 'Descargar PDF' }).click()
  expect((await descarga).suggestedFilename()).toBe('reporte-mensual-2026-09.pdf')
  await page.getByLabel('Mes del reporte').selectOption('2026-10')
  await page.getByRole('button', { name: 'Generar PDF' }).click()
  await expect(page.getByRole('button', { name: 'Descargar PDF' })).toBeVisible()
  const gen = esc.find((e) => e.ruta === '/api/v1/reports/monthly')!
  expect(gen.cuerpo).toEqual({ month: '2026-10', send_telegram: false })
  expect(gen.clave!.length).toBeGreaterThanOrEqual(16)
})

test('F-5 · Estado de cuenta de Karen → Generar PDF → Compartir', async ({ page }) => {
  await new ApiFalsa().instalar(page.context())
  const esc = await instalarArea(page.context())
  await page.goto('/reportes')
  await page.getByRole('button', { name: 'Estado de cuenta de Karen', exact: true }).first().click()
  await expect(page).toHaveURL(/\/reportes\/karen$/)
  await page.getByRole('button', { name: 'Generar PDF' }).click()
  await expect(page.getByRole('button', { name: 'Compartir' }).first()).toBeVisible()
  expect(esc.find((e) => e.ruta === '/api/v1/reports/external-statements')!.cuerpo).toEqual({ person_id: 3, period_start_date: '2026-10-01', period_end_date: '2026-10-31' })
})

test('F-6 variante · Salud → Registrar peso → Guardar recarga el historial', async ({ page }) => {
  const api = new ApiFalsa()
  await api.instalar(page.context())
  await instalarArea(page.context())
  let lecturas = 0
  page.on('request', (r) => { if (r.method() === 'GET' && new URL(r.url()).pathname === '/api/v1/weight-entries') lecturas++ })
  await page.goto('/salud')
  await expect(page.getByTestId('ficha-salud').first()).toBeVisible()
  const antes = lecturas
  await page.getByRole('button', { name: 'Registrar peso' }).click()
  await page.getByRole('dialog').getByLabel('Peso (kg)').fill('84.2')
  await page.getByRole('dialog').getByRole('button', { name: 'Guardar' }).click()
  await expect.poll(() => api.pesos.length).toBe(1)
  await expect.poll(() => lecturas).toBeGreaterThan(antes)
})
