// P-10 a P-16 montadas en 412×915 (móvil) y 1440×900 (escritorio) con la API simulada: elementos principales y sin desborde horizontal.
import { expect, test, type Page } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { instalarArea, type Pet } from './datos'
const HOY = new Date('2026-10-12T10:00:00-06:00')
async function abrir(page: Page, ruta: string, peticiones: Pet[] = []) {
  await new ApiFalsa().instalar(page.context())
  await instalarArea(page.context(), peticiones)
  await page.clock.setFixedTime(HOY)
  await page.goto(ruta)
}
const sinDesborde = async (page: Page) => expect(await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth)).toBeLessThanOrEqual(0)

test('P-10 Plan de la semana: semana por defecto, cuadrícula, casillas y totales', async ({ page }) => {
  await abrir(page, '/comidas')
  await expect(page.getByRole('heading', { level: 1, name: 'Plan de la semana' })).toBeVisible()
  await expect(page.getByTestId('semana')).toHaveText('Semana 42 · 12-oct-2026 al 18-oct-2026')
  await expect(page.getByTestId('cuadricula').getByText('Picadillo de res con papa').first()).toBeVisible()
  await expect(page.getByRole('button', { name: 'Generar lista de compras' })).toBeVisible()
  await expect(page.getByRole('button', { name: /^Cena: Vacía/ }).first()).toBeVisible()
  await sinDesborde(page)
  await page.getByRole('tab', { name: 'Resumen' }).click()
  await expect(page.getByText('$68.76').first()).toBeVisible()
  await expect(page.getByText('20 de 21')).toBeVisible()
})

test('P-10 Tendencias (core R5-01): las semanas se rotulan con fechaConDia, «Lun 5 oct 2026», en el eje y en «Ver como tabla»', async ({ page }) => {
  await abrir(page, '/comidas')
  await page.getByRole('tab', { name: 'Tendencias' }).click()
  for (const titulo of ['Costo semanal de comidas', 'kcal promedio por día contra la meta']) {
    const figura = page.locator('figure').filter({ hasText: titulo })
    await expect(figura.locator('svg text', { hasText: 'Lun 5 oct 2026' }).first()).toBeAttached()
    await figura.getByText('Ver como tabla').click()
    await expect(figura.getByRole('rowheader')).toHaveText(['Lun 5 oct 2026', 'Lun 12 oct 2026'])
  }
  // Ningún formato híbrido «Lun 5 oct-2026» ni ISO en la vista, tampoco dentro del SVG.
  const texto = await page.locator('main').evaluate((m) => m.textContent ?? '')
  expect(texto).not.toMatch(/\b(Lun|Mar|Mié|Jue|Vie|Sáb|Dom) \d{1,2} [a-z]{3}-\d{4}/)
  expect(texto).not.toMatch(/\b\d{4}-\d{2}-\d{2}\b/)
  await sinDesborde(page)
})
test('P-11 Recetario: galería con costo y kcal por porción', async ({ page }) => {
  await abrir(page, '/comidas/recetas')
  await expect(page.getByText('Yogur con manzana y avena').first()).toBeVisible()
  await expect(page.getByRole('button', { name: 'Nueva receta' })).toBeVisible()
  await sinDesborde(page)
})
test('P-12 Receta: escalado, totales, ingredientes sin precio y sin fuente', async ({ page }) => {
  const pets: Pet[] = []
  await abrir(page, '/comidas/recetas/1', pets)
  await expect(page.getByRole('heading', { level: 2, name: 'Picadillo de res con papa' })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Agregar al plan' })).toBeVisible()
  await expect(page.getByText('$68.76').first()).toBeVisible()
  await expect(page.getByText('sin fuente').first()).toBeVisible()
  await page.getByRole('button', { name: 'Una porción más' }).click()
  await expect(page.getByTestId('porciones')).toHaveText('3')
  expect(pets.some((p) => p.ruta === '/api/v1/recipes/1' && p.consulta.get('servings_milli') === '3000')).toBe(true)
  await sinDesborde(page)
})
test('P-13 Lista de compras: total, grupos por tienda y marcar comprado con la semana', async ({ page }) => {
  const pets: Pet[] = []
  await abrir(page, '/compras', pets)
  await expect(page.getByTestId('total-lista')).toHaveText('$77.80')
  await expect(page.getByText('Día de visita: 17-oct-2026')).toHaveCount(1)
  await expect(page.getByText('Mercado Portales').first()).toBeVisible()
  await page.getByRole('checkbox', { name: 'Marcar Leche entera como comprado' }).first().click()
  await expect.poll(() => pets.find((p) => p.metodo === 'PATCH')?.consulta.get('week')).toBe('2026-W42')
  await sinDesborde(page)
})
test('P-14 Comparar canasta: cifras canon de la semilla', async ({ page }) => {
  await abrir(page, '/compras/comparar')
  await expect(page.getByTestId('frase-barata')).toHaveText('Bodega Aurrera es la más barata: $3,094.50.')
  const walmart = page.getByTestId('tienda-1')
  await expect(walmart).toContainText('$3,233.00')
  await expect(walmart).toContainText('+$138.50 contra la más barata · +$100.50 en idénticos')
  await expect(walmart).toContainText('33 idénticos · 3 sustitutos')
  await expect(page.getByText('Empates: 8 · 36 artículos en total')).toBeVisible()
  await expect(page.getByRole('button', { name: 'Usar Bodega Aurrera' })).toBeVisible()
  await sinDesborde(page)
})
test('P-15 Productos y precios: catálogo y hoja de captura', async ({ page }) => {
  await abrir(page, '/compras/productos')
  await expect(page.getByText('Jitomate saladet').first()).toBeVisible()
  await sinDesborde(page)
  await page.getByRole('button', { name: 'Capturar precio' }).click()
  await expect(page.getByRole('dialog', { name: 'Capturar precio' })).toBeVisible()
})
test('P-16 Inventario: agotado primero y resumen por estado', async ({ page }) => {
  await abrir(page, '/compras/inventario')
  await expect(page.getByText('Jitomate saladet').first()).toBeVisible()
  await expect(page.getByText('Agotado', { exact: true }).filter({ visible: true }).first()).toBeVisible()
  await sinDesborde(page)
  await page.getByRole('tab', { name: 'Resumen' }).click()
  await expect(page.getByRole('button', { name: /Agotado/ }).first()).toBeVisible()
})
