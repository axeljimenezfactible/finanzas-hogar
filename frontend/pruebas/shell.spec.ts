// Prueba base del cascarón: barra lateral (9 secciones en orden), menú «/», tema según el sistema e interruptor.
import { expect, test } from '@playwright/test'
import { ApiFalsa } from './apoyo/api-falsa'
const SECCIONES = ['Inicio', 'Gastos', 'Presupuesto', 'Comidas', 'Compras y tiendas', 'Créditos', 'Salud', 'Reportes', 'Ajustes']
test('barra lateral, navegación y menú «/»', async ({ page, context, isMobile }) => {
  await new ApiFalsa().instalar(context)
  await page.goto('/inicio')
  if (isMobile) {
    await expect(page.getByRole('navigation', { name: 'Navegación inferior' })).toBeVisible()
    await page.getByRole('button', { name: 'Abrir menú' }).click()
  }
  const nav = page.getByRole('navigation', { name: 'Navegación principal' }).last()
  for (const s of SECCIONES) await expect(nav.getByRole('link', { name: s, exact: true })).toBeVisible()
  const orden = await nav.getByRole('link').evaluateAll((ls) => ls.map((l) => l.textContent))
  expect(orden.filter((t) => SECCIONES.includes(t ?? ''))).toEqual(SECCIONES)
  await nav.getByRole('link', { name: 'Gastos', exact: true }).click()
  await expect(page).toHaveURL(/\/gastos$/)
  await expect(page.getByRole('heading', { level: 1, name: 'Todos los gastos' })).toBeVisible()
  await page.keyboard.press('Control+k')
  const menu = page.getByRole('dialog')
  await menu.getByRole('textbox').fill('/peso')
  await menu.getByRole('textbox').press('Enter')
  await expect(page.getByRole('dialog', { name: 'Registrar' })).toBeVisible()
  await expect(page.getByRole('tab', { name: 'Peso' })).toHaveAttribute('aria-selected', 'true')
})
test('U5-10: tema según el sistema y el interruptor lo cambia', async ({ page, context }) => {
  await new ApiFalsa().instalar(context)
  await page.emulateMedia({ colorScheme: 'dark' })
  await page.goto('/inicio')
  const fondo = () => page.evaluate(() => getComputedStyle(document.body).backgroundColor)
  const oscuro = await fondo()
  await page.emulateMedia({ colorScheme: 'light' })
  const claro = await fondo()
  expect(oscuro).not.toEqual(claro)
  await page.evaluate(() => localStorage.setItem('fh-tema', 'oscuro'))
  await page.reload()
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark')
  expect(await fondo()).toEqual(oscuro)
})
