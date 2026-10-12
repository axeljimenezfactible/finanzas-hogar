// Pruebas del área Gastos y Presupuesto (WF-3a): una por pantalla en 412×915 (movil) y 1440×900 (escritorio),
// más F-1 (U5-09: ≤ 3 toques y el gasto aparece en la tabla de P-02) y el flujo complementario del ticket (P-05 → P-04).
import { expect, test, type Locator } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { AreaFalsa, sinDesborde } from './datos'

let api: ApiFalsa, area: AreaFalsa
test.beforeEach(async ({ context }) => {
  api = new ApiFalsa(); await api.instalar(context)
  area = new AreaFalsa(api); await area.instalar(context)
})
const movil = () => test.info().project.name === 'movil'
const tocar = (l: Locator) => (movil() ? l.tap() : l.click())
async function alto44(l: Locator) { if (movil()) expect((await l.boundingBox())!.height).toBeGreaterThanOrEqual(44) }

test('P-02 Todos los gastos: tabla, totales, filtros y resumen', async ({ page }) => {
  await page.goto('/gastos?periodo=2026-10')
  await expect(page.getByRole('heading', { level: 1, name: 'Todos los gastos' })).toBeVisible()
  await expect(page.getByText('Súper semanal')).toBeVisible()
  await expect(page.getByText('$1,850.50').first()).toBeVisible()
  await expect(page.getByLabel('Totales del filtro')).toContainText('$3,050.50')
  await expect(page.getByText('oct-2026', { exact: true })).toBeVisible()
  const registrar = page.getByRole('main').getByRole('button', { name: 'Registrar gasto' })
  await alto44(registrar)
  await page.getByRole('button', { name: 'Filtros' }).click()
  await expect(page.getByLabel('Familia')).toBeVisible()
  await sinDesborde(page)
  await page.goto('/gastos?periodo=2026-10&sub=resumen')
  await expect(page.getByText('Gasto del hogar en el periodo')).toBeVisible()
  await expect(page.getByText('$200.00 más que el periodo anterior')).toBeVisible()
})

test('P-04 Gasto: propiedades, reparto, renglones y gasto borrado', async ({ page }) => {
  await page.goto('/gastos/11')
  await expect(page.getByRole('heading', { level: 1, name: 'Gasto' })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Propiedades' })).toBeVisible()
  await expect(page.getByText('Parte de Karen').first()).toBeVisible()
  await expect(page.getByText('$350.50').first()).toBeVisible()
  await expect(page.getByText('Leche 1 L')).toBeVisible()
  await expect(page.getByRole('button', { name: 'Duplicar' })).toBeVisible()
  await sinDesborde(page)
  await page.goto('/gastos/999')
  await expect(page.getByText('Este gasto ya no existe')).toBeVisible()
  await expect(page.getByRole('button', { name: 'Volver' })).toBeVisible()
})

test('P-05 Capturar ticket: pasos, por confirmar y flujo foto → confirmar → P-04', async ({ page }) => {
  await page.goto('/gastos/ticket')
  await expect(page.getByRole('heading', { level: 1, name: 'Capturar ticket' })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Tomar foto' })).toBeVisible()
  await expect(page.getByText('Con borrador', { exact: false })).toBeVisible()
  await alto44(page.getByRole('button', { name: 'Tomar foto' }))
  await page.getByLabel('Elegir de la galería').setInputFiles({ name: 'ticket.png', mimeType: 'image/png', buffer: Buffer.from('falso') })
  await expect(page).toHaveURL(/lectura=501/)
  await expect(page.getByRole('heading', { name: 'Renglones' })).toBeVisible()
  await expect(page.getByText('Confianza: Baja')).toBeVisible()
  await expect(page.getByText('La suma de los renglones cuadra con el total.')).toBeVisible()
  await sinDesborde(page)
  await page.getByLabel('Categoría').selectOption({ label: 'Despensa' })
  await tocar(page.getByRole('button', { name: 'Confirmar ticket' }))
  await expect(page).toHaveURL(/\/gastos\/11$/)
  const cuerpo = area.confirmados[0] as { expense: { amount_cents: number; items: unknown[]; category_id: number } }
  expect(cuerpo.expense).toMatchObject({ amount_cents: 25000, category_id: 1 })
  expect(cuerpo.expense.items).toHaveLength(2)
})

test('P-06 Presupuesto del mes: resumen, alertas, detalle y ajuste con alcance', async ({ page }) => {
  await page.goto('/presupuesto?periodo=2026-10')
  await expect(page.getByRole('heading', { level: 1, name: 'Presupuesto del mes' })).toBeVisible()
  await expect(page.getByText('$30,000.00')).toBeVisible()
  await expect(page.getByRole('meter').first()).toHaveAttribute('aria-valuetext', '84 % del plan')
  await expect(page.getByText('Despensa: 103 % usado desde el 9-oct-2026')).toBeVisible()
  await sinDesborde(page)
  await page.goto('/presupuesto?periodo=2026-10&sub=detalle')
  await expect(page.getByText('Despensa devengada')).toBeVisible()
  await page.getByRole('main').getByRole('button', { name: 'Ajustar presupuesto' }).click()
  const hoja = page.getByRole('dialog', { name: 'Ajustar presupuesto' })
  await expect(hoja.getByLabel('De este mes en adelante')).toBeChecked()
  await hoja.getByLabel('Monto').fill('9000')
  await expect(hoja.getByText('→ $9,000.00')).toBeVisible()
  await hoja.getByRole('button', { name: 'Guardar' }).click()
  await expect.poll(() => area.presupuestos.length).toBe(1)
  expect(area.presupuestos[0].url).toContain('scope=onward')
  expect(area.presupuestos[0].cuerpo).toEqual({ amount_cents: 900000 })
})

test('P-07 Flujo y proyección: interruptor A/B, deuda viva y supuestos', async ({ page }) => {
  await page.goto('/presupuesto/flujo')
  await expect(page.getByRole('heading', { level: 1, name: 'Flujo y proyección' })).toBeVisible()
  const interruptor = page.getByRole('switch')
  await expect(interruptor).toHaveAttribute('aria-checked', 'false')
  await expect(page.getByText('Proyección: Solo lo contratado')).toBeVisible()
  await expect(page.getByText('$6,293.44')).toBeVisible()
  await alto44(interruptor)
  await interruptor.click()
  await expect(interruptor).toHaveAttribute('aria-checked', 'true')
  await expect(page.getByText('Proyección: Con la despensa y el papel repetidos')).toBeVisible()
  await page.getByRole('button', { name: 'Ver supuestos' }).click()
  await expect(page.getByText('Promedio de 3 meses')).toBeVisible()
  await sinDesborde(page)
})

test('P-08 Ingresos: quincenas, extras y capturar quincena', async ({ page }) => {
  await page.goto('/presupuesto/ingresos?periodo=2026-10')
  await expect(page.getByRole('heading', { level: 1, name: 'Ingresos' })).toBeVisible()
  await expect(page.getByText('Venta de bicicleta')).toBeVisible()
  await expect(page.getByText('Esperada')).toBeVisible()
  await sinDesborde(page)
  await page.getByRole('main').getByRole('button', { name: 'Capturar quincena' }).first().click()
  const hoja = page.getByRole('dialog', { name: 'Capturar quincena' })
  await expect(hoja.getByLabel('Monto')).toHaveValue('15000.00')
  await hoja.getByRole('button', { name: 'Guardar' }).click()
  await expect.poll(() => area.ingresos.length).toBe(1)
  expect(area.ingresos[0].cuerpo).toMatchObject({ kind: 'fixed', amount_cents: 1500000, income_date: '2026-10-30' })
  expect(area.ingresos[0].clave.length).toBeGreaterThanOrEqual(16)
})

test('P-09 Ahorro y metas: totales, metas con avance y aportación', async ({ page }) => {
  await page.goto('/presupuesto/ahorro')
  await expect(page.getByRole('heading', { level: 1, name: 'Ahorro y metas' })).toBeVisible()
  await expect(page.getByText('$45,000.00')).toBeVisible()
  await expect(page.getByText('Vacaciones')).toBeVisible()
  await expect(page.getByRole('meter', { name: '' }).first()).toBeVisible()
  await sinDesborde(page)
  await page.getByRole('main').getByRole('button', { name: 'Registrar aportación' }).click()
  const hoja = page.getByRole('dialog', { name: 'Registrar aportación' })
  await hoja.getByLabel('Monto').fill('500')
  await hoja.getByRole('button', { name: 'Guardar' }).click()
  await expect.poll(() => area.movimientos.length).toBe(1)
  expect(area.movimientos[0].cuerpo).toMatchObject({ kind: 'deposit', amount_cents: 50000, savings_goal_id: 1 })
})

test('F-1 (U5-09): registrar gasto en ≤ 3 toques y verlo en la tabla de P-02', async ({ page }) => {
  await page.goto('/gastos?periodo=2026-10')
  await expect(page.getByText('Súper semanal')).toBeVisible()
  let toques = 0
  const toque = async (l: Locator) => { await tocar(l); toques++ }
  await toque(page.getByRole('button', { name: 'Registrar gasto' }).last())
  await page.getByLabel('Monto').fill('185.5')
  await toque(page.getByRole('dialog').getByRole('button', { name: 'Comida fuera' }))
  await toque(page.getByRole('dialog').getByRole('button', { name: 'Guardar' }))
  console.log(`U5-09 toques para registrar un gasto (P-02, ${test.info().project.name}): ${toques}`)
  expect(toques).toBeLessThanOrEqual(3)
  await expect(page.getByText('$185.50').first()).toBeVisible()
  expect(api.gastos).toHaveLength(1)
})
