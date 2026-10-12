// Flujos de U1 que pasan por Comidas y Compras, en 412×915 táctil: F-2 (planear la semana) y F-3 rutas A y B (comparar canasta).
import { expect, test } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { instalarArea, type Pet } from './datos'
test.beforeEach(async ({ page }) => { await page.clock.setFixedTime(new Date('2026-10-12T10:00:00-06:00')) })

test('F-2: casilla vacía → receta → generar lista (toques en controles: 3 desde P-10)', async ({ page, context }) => {
  const pets: Pet[] = []
  await new ApiFalsa().instalar(context); await instalarArea(context, pets)
  await page.goto('/comidas')
  await page.getByRole('button', { name: /^Cena: Vacía/ }).first().tap()
  await expect(page.getByRole('dialog', { name: 'Elegir receta · Lun 12 · Cena' })).toBeVisible()
  await page.getByRole('button', { name: /Picadillo de res con papa/ }).tap()
  await expect.poll(() => pets.find((p) => p.metodo === 'POST' && p.ruta === '/api/v1/meal-plan-entries')?.cuerpo).toEqual({ plan_dates: ['2026-10-12'], meal_slot_id: 3, person_ids: [1, 2], is_eating_out: false, recipe_id: 1, servings_milli: 1000 })
  expect(pets.find((p) => p.ruta === '/api/v1/meal-plan-entries')?.clave?.length).toBeGreaterThanOrEqual(16)
  await page.getByRole('button', { name: 'Generar lista de compras' }).tap()
  await expect(page).toHaveURL(/\/compras$/)
  expect(pets.find((p) => p.ruta === '/api/v1/shopping-list/generate')?.consulta.get('week')).toBe('2026-W42')
  await expect(page.getByTestId('total-lista')).toHaveText('$77.80')
})
test('F-3 ruta A: Resumen → Detalle → «Usar Bodega Aurrera» abre la lista', async ({ page, context }) => {
  const pets: Pet[] = []
  await new ApiFalsa().instalar(context); await instalarArea(context, pets)
  await page.goto('/compras/comparar')
  await expect(page.getByTestId('frase-barata')).toContainText('$3,094.50')
  await page.getByRole('tab', { name: 'Detalle' }).tap()
  await page.getByRole('button', { name: 'Usar Bodega Aurrera' }).tap()
  await expect(page).toHaveURL(/\/compras$/)
  const asignar = pets.find((p) => p.ruta === '/api/v1/shopping-list/assign-store')
  expect(asignar?.cuerpo).toEqual({ store_id: 2 })
  expect(asignar?.consulta.get('week')).toBe('2026-W42')
})
test('F-3 ruta B: Lista → «⋯» → Comparar tiendas → compra del 10-oct-2026', async ({ page, context }) => {
  const pets: Pet[] = []
  await new ApiFalsa().instalar(context); await instalarArea(context, pets)
  await page.goto('/compras')
  await page.getByRole('button', { name: 'Más opciones', exact: true }).last().tap()
  await page.getByRole('menuitem', { name: 'Comparar tiendas' }).tap()
  await expect(page).toHaveURL(/\/compras\/comparar$/)
  await page.locator('select#p14-canasta').selectOption({ label: 'Compra del 10-oct-2026 · Walmart · $3,233.00' })
  await expect.poll(() => pets.some((p) => p.ruta === '/api/v1/basket-comparison' && p.consulta.get('basket') === 'expense' && p.consulta.get('expense_id') === '77')).toBe(true)
})
