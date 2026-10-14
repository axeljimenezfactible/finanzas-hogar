// WF-3a bis (front-comidas-compras-fix): N-123 en P-12 y N-173 / N-175 en P-15 con la API simulada (C1 v1.4), en 412×915
// táctil y en 1440×900. Datos falsos, sin credenciales. Cada operación nueva se cuenta con las peticiones que llegan al mock.
import { expect, test, type BrowserContext, type Page } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { dineroSinTabulares, textosEnIngles } from '../apoyo/medidas'
import { PRODUCTOS, instalarArea, type Pet } from './datos'

const HOY = new Date('2026-10-12T10:00:00-06:00') // lunes: la semana por defecto es 2026-W42 (12 al 18-oct)
const T = '2026-10-01T09:00:00-06:00'
const pag = (items: unknown[]) => ({ items, next_cursor: null })
const presentacion = (id: number, name: string, size_milli: number, extra = {}) => ({ id, product_id: 1, name, size_milli, is_bulk: false, is_active: true, created_at: T, updated_at: T, ...extra })
const DETALLE = { ...PRODUCTOS[0], presentations: [presentacion(1, 'Granel', 1_000_000, { is_bulk: true }), presentacion(2, 'Bolsa', 500_000, { brand: 'Del campo' })], nutrition_facts: [] }
// Bodega Aurrera (tienda 2), 10-oct-2026: $28.90 por 1 kg → $28.90 por kg (C1: precio por 1 kg, 1 l o 1 pieza)
const PRECIOS = pag([{ id: 11, product_id: 1, presentation_id: 1, store_id: 2, price_cents: 2890, base_unit_price_cents: 2890, observed_date: '2026-10-10', source: 'ticket_manual', is_stale: false, age_days: 2, created_at: T, updated_at: T }])
const NUTRICION = pag([{ id: 5, product_id: 1, kcal_milli_per_100g: 18_000, source: 'usda_fdc', source_ref: 'FDC 170457', is_preferred: true, created_at: T, updated_at: T }])

/** Rutas de P-15 que no están en datos.ts: ficha, precios, nutrición (N-173) y las dos escrituras de N-175. */
async function instalarFicha(ctx: BrowserContext, pets: Pet[]) {
  await ctx.route(/\/api\/v1\/(products\/\d+(\/prices|\/nutrition-facts)?|presentations\/\d+)(\?.*)?$/, async (route) => {
    const r = route.request(); const u = new URL(r.url()); const m = `${r.method()} ${u.pathname}`
    const cuerpo = r.postData() ? r.postDataJSON() : null
    let resp: unknown
    if (m === 'GET /api/v1/products/1') resp = DETALLE
    else if (m === 'GET /api/v1/products/1/prices') resp = PRECIOS
    else if (m === 'GET /api/v1/products/1/nutrition-facts') resp = NUTRICION
    else if (r.method() === 'PATCH' && /^\/api\/v1\/products\/\d+$/.test(u.pathname)) resp = { ...DETALLE, id: Number(u.pathname.split('/').pop()), ...cuerpo }
    else if (r.method() === 'PATCH' && /^\/api\/v1\/presentations\/\d+$/.test(u.pathname)) resp = { ...DETALLE.presentations[0], ...cuerpo }
    if (resp === undefined) return route.fallback()
    pets.push({ metodo: r.method(), ruta: u.pathname, consulta: u.searchParams, cuerpo, clave: r.headers()['idempotency-key'] })
    return route.fulfill({ json: resp })
  })
}
async function abrir(page: Page, ruta: string, pets: Pet[] = []) {
  await new ApiFalsa().instalar(page.context())
  await instalarArea(page.context(), pets)
  await instalarFicha(page.context(), pets)
  await page.clock.setFixedTime(HOY)
  await page.goto(ruta)
}
const escrituras = (pets: Pet[], metodo: string, ruta: RegExp) => pets.filter((p) => p.metodo === metodo && ruta.test(p.ruta))

test('P-12 N-123: «Agregar al plan» asigna en la misma pantalla con createMealPlanEntries (días, tiempo apto, Ambos, porciones)', async ({ page }) => {
  const pets: Pet[] = []
  await abrir(page, '/comidas/recetas/1', pets)
  await page.getByRole('button', { name: 'Agregar al plan' }).click()
  const hoja = page.getByRole('dialog', { name: 'Agregar al plan · Picadillo de res con papa' })
  await expect(hoja).toBeVisible()
  await expect(hoja.getByTestId('semana')).toContainText('12-oct-2026') // semana por defecto (U1 §0.2)
  await expect(hoja.getByRole('button', { name: 'Desayuno', exact: true })).toHaveCount(0) // la receta es apta para comida y cena
  await expect(hoja.getByRole('button', { name: 'Comida', exact: true })).toHaveAttribute('aria-pressed', 'true')
  await expect(hoja.getByRole('button', { name: 'Ambos' })).toHaveAttribute('aria-pressed', 'true')
  await expect(hoja.getByRole('button', { name: 'Lun 12' })).toHaveAttribute('aria-pressed', 'true') // hoy
  await expect(hoja.getByText('Porciones por día en total: 2')).toBeVisible() // decisión D: con Ambos son 2
  await hoja.getByRole('button', { name: 'Cena', exact: true }).click()
  await hoja.getByRole('button', { name: 'Mar 13' }).click()
  await hoja.getByRole('button', { name: 'Agregar al plan' }).click()
  await expect(page.getByText('Receta agregada al plan en 4 casillas')).toBeVisible()
  await expect(hoja).toHaveCount(0)
  const post = escrituras(pets, 'POST', /^\/api\/v1\/meal-plan-entries$/)
  expect(post).toHaveLength(1)
  expect(post[0].cuerpo).toEqual({ plan_dates: ['2026-10-12', '2026-10-13'], meal_slot_id: 3, person_ids: [1, 2], recipe_id: 1, servings_milli: 1000, is_eating_out: false })
  expect(post[0].clave).toBeTruthy()
  await expect(page).toHaveURL(/\/comidas\/recetas\/1$/) // ya no manda a P-10 ?hoja=asignar
  expect(await textosEnIngles(page)).toEqual([])
})

test('P-12 N-123: un 409 del servidor se muestra en la hoja y no se cierra', async ({ page, context }) => {
  await abrir(page, '/comidas/recetas/1')
  await context.route('**/api/v1/meal-plan-entries', (r) => r.request().method() === 'POST'
    ? r.fulfill({ status: 409, json: { code: 'already_exists', message: 'Esa receta ya está en esa casilla.', request_id: 'prueba-409' } }) : r.fallback())
  await page.getByRole('button', { name: 'Agregar al plan' }).click()
  const hoja = page.getByRole('dialog', { name: 'Agregar al plan · Picadillo de res con papa' })
  await hoja.getByLabel('Porciones por persona').fill('0')
  await expect(hoja.getByRole('button', { name: 'Agregar al plan' })).toBeDisabled() // porción inválida: no envía
  await hoja.getByLabel('Porciones por persona').fill('1.5')
  await expect(hoja.getByText('Porciones por día en total: 3')).toBeVisible()
  await hoja.getByRole('button', { name: 'Agregar al plan' }).click()
  await expect(hoja.getByText('Esa receta ya está en esa casilla.')).toBeVisible()
  await expect(hoja).toBeVisible()
})

test('P-15 N-173: la ficha del producto lee la nutrición con listNutritionFacts (kcal, fuente y cita)', async ({ page }) => {
  const pets: Pet[] = []
  await abrir(page, '/compras/productos?hoja=producto&id=1', pets)
  const hoja = page.getByRole('dialog', { name: 'Jitomate saladet' })
  await expect(hoja.getByTestId('nutricion')).toHaveText('18 kcal por 100 g · USDA · FDC 170457 · preferida')
  expect(pets.some((p) => p.metodo === 'GET' && p.ruta === '/api/v1/products/1/nutrition-facts')).toBe(true)
  await expect(hoja.getByText('$28.90 por kg')).toBeVisible() // precio por 1 kg (C1, K7), no «por g»
  expect(await dineroSinTabulares(page)).toEqual([])
  expect(await textosEnIngles(page)).toEqual([])
})

test('P-15 N-175: categoría en el lugar del catálogo con updateProduct y precio por kg', async ({ page }) => {
  const pets: Pet[] = []
  await abrir(page, '/compras/productos', pets)
  await expect(page.getByText('$2.89 por kg')).toBeVisible()
  await page.getByRole('button', { name: 'Lácteos y huevo' }).first().click() // Leche entera
  await page.getByRole('combobox', { name: 'Categoría' }).last().selectOption('1')
  await page.keyboard.press('Enter')
  await expect.poll(() => escrituras(pets, 'PATCH', /^\/api\/v1\/products\/\d+$/).length).toBe(1)
  const [w] = escrituras(pets, 'PATCH', /^\/api\/v1\/products\/\d+$/)
  expect(w.ruta).toBe('/api/v1/products/2')
  expect(w.cuerpo).toEqual({ product_group_id: 1 })
})

test('P-15 N-175: nombre y presentación en el lugar (updateProduct, updatePresentation) con «Antes → Después» del tamaño', async ({ page }) => {
  const pets: Pet[] = []
  await abrir(page, '/compras/productos?hoja=producto&id=1', pets)
  const hoja = page.getByRole('dialog', { name: 'Jitomate saladet' })
  // Nombre del producto: Enter guarda
  await hoja.getByRole('button', { name: 'Jitomate saladet' }).click()
  await hoja.getByRole('textbox', { name: 'Nombre' }).fill('Jitomate saladet rojo')
  await hoja.getByRole('textbox', { name: 'Nombre' }).press('Enter')
  await expect.poll(() => escrituras(pets, 'PATCH', /^\/api\/v1\/products\/1$/).map((p) => p.cuerpo)).toEqual([{ name: 'Jitomate saladet rojo' }])
  // Tamaño: la vista previa sale antes de guardar; Esc descarta sin escribir
  await hoja.getByRole('button', { name: '1 kg · a granel' }).click()
  const tam = hoja.getByRole('table').getByRole('textbox', { name: 'Tamaño (g)' })
  await tam.fill('900')
  const previa = hoja.getByTestId('vista-previa')
  await expect(previa).toContainText('Tamaño')
  await expect(previa).toContainText('1 kg')
  await expect(previa).toContainText('900 g')
  await expect(previa).toContainText('Último precio por kg (Bodega Aurrera, 10-oct-2026)')
  await expect(previa).toContainText('$28.90')
  await expect(previa).toContainText('$32.11') // 2,890 × 1,000,000 ÷ 900,000 = 3,211.1 → 3,211 (HALF_UP)
  console.log(`N-175 vista previa del tamaño: «${(await previa.innerText()).replace(/\s+/g, ' ')}»`)
  await tam.press('Escape')
  await expect(previa).toHaveCount(0)
  await expect(hoja).toBeVisible() // Esc descarta solo la edición; la hoja sigue abierta
  await expect(hoja.getByRole('button', { name: '1 kg · a granel' })).toBeVisible()
  expect(escrituras(pets, 'PATCH', /^\/api\/v1\/presentations\//)).toEqual([])
  // Enter guarda con la vista previa ya calculada
  await hoja.getByRole('button', { name: '1 kg · a granel' }).click()
  await tam.fill('900')
  await expect(previa).toContainText('$32.11')
  await tam.press('Enter')
  await expect.poll(() => escrituras(pets, 'PATCH', /^\/api\/v1\/presentations\/1$/).map((p) => p.cuerpo)).toEqual([{ size_milli: 900_000 }])
  // Marca de la segunda presentación
  await hoja.getByRole('button', { name: 'Del campo' }).click()
  await hoja.getByRole('textbox', { name: 'Marca' }).fill('Campo verde')
  await hoja.getByRole('textbox', { name: 'Marca' }).press('Enter')
  await expect.poll(() => escrituras(pets, 'PATCH', /^\/api\/v1\/presentations\/2$/).map((p) => p.cuerpo)).toEqual([{ brand: 'Campo verde' }])
  expect(await textosEnIngles(page)).toEqual([])
})

test('P-15 N-175 · Ctrl+/: categoría de varios productos con vista previa sumada; Esc no escribe y Enter aplica a todas', async ({ page, isMobile }) => {
  test.skip(isMobile, 'Ctrl+/ es atajo de escritorio (U4 §0.7); en táctil queda la edición por celda de la prueba anterior')
  const pets: Pet[] = []
  await abrir(page, '/compras/productos', pets)
  await page.getByRole('checkbox', { name: 'Seleccionar Leche entera' }).check()
  await page.getByRole('checkbox', { name: 'Seleccionar Huevo blanco' }).check()
  await page.keyboard.press('Control+/')
  const bloque = page.getByRole('region', { name: 'Editar 2 registros a la vez' })
  await expect(bloque).toBeVisible()
  const valor = bloque.getByRole('combobox', { name: 'Valor nuevo' })
  await valor.selectOption('1')
  const previa = bloque.getByTestId('vista-previa')
  await expect(previa).toContainText('Productos que cambian de categoría')
  await expect(previa).toContainText('2')
  await valor.press('Escape')
  await expect(bloque).toHaveCount(0)
  expect(escrituras(pets, 'PATCH', /^\/api\/v1\/products\//)).toEqual([])
  console.log(`Foco tras Esc: ${await page.evaluate(() => { const e = document.activeElement; return e ? `${e.tagName} ${e.getAttribute('aria-label') ?? ''}` : 'ninguno' })}`)
  await page.getByRole('checkbox', { name: 'Seleccionar Huevo blanco' }).focus() // la vista vuelve a tener el foco
  await page.keyboard.press('Control+/')
  await valor.selectOption('1')
  await expect(previa).toContainText('2')
  await valor.press('Enter')
  await expect(bloque).toHaveCount(0)
  await expect.poll(() => escrituras(pets, 'PATCH', /^\/api\/v1\/products\//).map((p) => `${p.ruta} ${JSON.stringify(p.cuerpo)}`).sort())
    .toEqual(['/api/v1/products/2 {"product_group_id":1}', '/api/v1/products/3 {"product_group_id":1}'])
})

test('P-14 con sucursales (como la semilla de desarrollo): «Cadena · Sucursal» y las cifras que manda la API', async ({ page, context }) => {
  await abrir(page, '/compras/comparar')
  const T0 = '2026-10-01T09:00:00-06:00'
  const s = (id: number, name: string, kind: string, extra = {}) => ({ id, kind, name, assisted_reading_enabled: false, requires_membership: false, has_membership: false, is_active: true, created_at: T0, updated_at: T0, ...extra })
  await context.route('**/api/v1/stores?*', (r) => r.fulfill({ json: pag([s(1, 'Bodega Aurrera', 'chain'), s(2, 'Vallejo', 'branch', { parent_id: 1 }), s(3, 'Walmart', 'chain'), s(4, 'Acueducto de Guadalupe', 'branch', { parent_id: 3 })]) }))
  await page.reload()
  await expect(page.getByTestId('frase-barata')).toHaveText('Bodega Aurrera · Vallejo es la más barata: $3,094.50.')
  await expect(page.getByRole('button', { name: 'Usar Bodega Aurrera · Vallejo' })).toBeVisible()
})
