// WF-3a bis · área Gastos y Presupuesto: una prueba por necesidad implementada (C3-matriz puntos 2 y 3), con la API
// simulada (AreaFalsa + las rutas nuevas de abajo, con la forma de C1 v1.4). Las de P-03 se saltan solas mientras
// src/shell/HojaRegistro.tsx no monte las piezas del área (solicitud a frontend-base) y corren en cuanto las monte.
import { expect, test, type Locator, type Page } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { AreaFalsa } from './datos'

const t0 = '2026-10-01T09:00:00-06:00'
const PNG = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+ip1sAAAAASUVORK5CYII=', 'base64')
const pagina = (items: unknown[]) => ({ items, next_cursor: null })
let area: AreaFalsa
let vistas: URLSearchParams[], borrados: string[], parches: { ruta: string; cuerpo: Record<string, unknown> }[]
test.beforeEach(async ({ context }) => {
  const api = new ApiFalsa(); await api.instalar(context)
  area = new AreaFalsa(api); await area.instalar(context)
  vistas = []; borrados = []; parches = []
  await context.route(/\/api\/v1\//, async (route) => {
    const pet = route.request(); const u = new URL(pet.url()); const ruta = `${pet.method()} ${u.pathname}`
    const json = (cuerpo: unknown, estado = 200) => route.fulfill({ status: estado, contentType: 'application/json', body: JSON.stringify(cuerpo) })
    if (ruta === 'GET /api/v1/views/expense') {
      vistas.push(u.searchParams); const mode = u.searchParams.get('mode')
      if (mode === 'board') return json({ ...pagina([]), mode, groups: [{ key: '1', label: 'Despensa', count: 1, total_cents: 150000, share_bp: 5556 }, { key: '6', label: 'Predial', count: 1, total_cents: 120000, share_bp: 4444 }] })
      if (mode === 'calendar') return json({ ...pagina([]), mode, groups: [{ key: '2026-10-05', label: '5 oct', count: 1, total_cents: 150000 }, { key: '2026-10-20', label: '20 oct', count: 1, total_cents: 120000 }] })
      return json({ ...pagina([{ id: 11, title: 'Súper semanal', date: '2026-10-05', amount_cents: 150000, image_attachment_id: 701, record: {} }, { id: 12, title: 'Predial', record: {} }]), mode, groups: [] })
    }
    if (ruta === 'GET /api/v1/attachments/701/file') return route.fulfill({ status: 200, contentType: 'image/png', body: PNG })
    if (ruta === 'DELETE /api/v1/attachments/701') { borrados.push('701'); return route.fulfill({ status: 204 }) }
    if (/^PATCH \/api\/v1\/(expenses|incomes)\/\d+$/.test(ruta)) { parches.push({ ruta, cuerpo: pet.postDataJSON() }); return json({ id: Number(u.pathname.split('/').pop()), ...pet.postDataJSON(), created_at: t0, updated_at: t0 }) }
    if (ruta === 'GET /api/v1/products') return json(pagina((u.searchParams.get('q') ?? '').toLowerCase().startsWith('le') ? [{ id: 5, name: 'Leche entera', base_unit: 'ml', is_active: true, created_at: t0, updated_at: t0 }] : []))
    if (ruta === 'GET /api/v1/products/5') return json({ id: 5, name: 'Leche entera', base_unit: 'ml', is_active: true, created_at: t0, updated_at: t0,
      presentations: [{ id: 51, product_id: 5, name: 'Caja de 1 L', size_milli: 1000000, is_bulk: false, is_active: true, created_at: t0, updated_at: t0 }] })
    if (ruta === 'GET /api/v1/shopping-list') return json({ start_date: '2026-10-12', end_date: '2026-10-18', total_cents: 5800, groups: [{ purchase_channel: 'chain', store_id: 2, total_cents: 5800,
      items: [{ id: 81, shopping_list_id: 3, origin: 'calculated', product_id: 5, presentation_id: 51, packages_to_buy: 2, is_bought: true, created_at: t0, updated_at: t0 }] }] })
    if (ruta === 'GET /api/v1/incomes') return json(pagina([{ id: 31, person_id: 1, kind: 'fixed', description: 'Quincena', amount_cents: 1500000, income_date: '2026-10-15', notes: 'Depósito', created_at: t0, updated_at: t0 }]))
    if (ruta === 'GET /api/v1/credits') return json(pagina([
      { id: 2, name: 'Línea de la tienda', kind: 'store_credit', holder_person_id: 1, payment_frequency: 'fortnightly', availability_reading: 'A', availability_confirmed: true, is_active: true, created_at: t0, updated_at: t0 },
      { id: 1, name: 'Tarjeta azul', kind: 'bank_card', holder_person_id: 1, payment_frequency: 'monthly', availability_reading: 'A', availability_confirmed: true, is_active: true, created_at: t0, updated_at: t0 },
      { id: 3, name: 'Línea cerrada', kind: 'revolving_line', holder_person_id: 1, payment_frequency: 'monthly', availability_reading: 'A', availability_confirmed: true, is_active: false, created_at: t0, updated_at: t0 }]))
    return route.fallback()
  })
})
const movil = () => test.info().project.name === 'movil'
const tocar = (l: Locator) => (movil() ? l.tap() : l.click())
const cargada = (img: Locator) => expect.poll(() => img.evaluate((e) => (e as HTMLImageElement).complete && (e as HTMLImageElement).naturalWidth > 0)).toBe(true)
async function abrirHoja(page: Page, pestana: 'Gasto' | 'Ingreso') {
  await page.goto('/gastos?periodo=2026-10')
  await tocar(page.getByRole('main').getByRole('button', { name: 'Registrar gasto' }))
  const hoja = page.getByRole('dialog', { name: 'Registrar' })
  await tocar(hoja.getByRole('tab', { name: pestana }))
  return hoja
}

test('P-02 N-043 · tablero por categoría con el total de queryView en la cabecera', async ({ page }) => {
  await page.goto('/gastos?periodo=2026-10')
  await tocar(page.getByRole('tab', { name: 'Tablero' }))
  await expect(page.getByRole('heading', { level: 3, name: /^Despensa · \$1,500\.00 \(56 %\)/ })).toBeVisible()
  const q = vistas.find((v) => v.get('mode') === 'board')!
  expect([q.get('group_by'), q.get('period')]).toEqual(['category_id', '2026-10'])
})

test('P-02 N-045 · calendario con la suma del día de queryView', async ({ page }) => {
  await page.goto('/gastos?periodo=2026-10&sub=calendario')
  await expect(page.getByText('Total del día · $1,500.00 · 1 gasto')).toBeVisible()
  await expect(page.getByText('Total del día · $1,200.00 · 1 gasto')).toBeVisible()
  expect(vistas.find((v) => v.get('mode') === 'calendar')?.get('date_field')).toBe('expense_date')
})

test('P-02 N-046 · galería de tickets con la foto vía downloadAttachment; la tabla no descarga fotos', async ({ page }) => {
  await page.goto('/gastos?periodo=2026-10')
  await expect(page.getByRole('cell', { name: 'Sí', exact: true })).toBeVisible()
  await expect(page.locator('table img')).toHaveCount(0)
  await tocar(page.getByRole('tab', { name: 'Galería' }))
  const foto = page.getByRole('img', { name: 'Foto del ticket: Súper semanal' })
  await expect(foto).toHaveAttribute('src', '/api/v1/attachments/701/file')
  await cargada(foto)
  await expect(page.getByText('Sin ticket')).toBeVisible()
})

test('P-02 Ctrl+/ · editar el tipo de varias filas con vista previa sumada; Esc no guarda y Enter aplica', async ({ page }) => {
  test.skip(movil(), 'Atajo de teclado: solo con teclado físico (proyecto escritorio)')
  await page.goto('/gastos?periodo=2026-10')
  await page.getByRole('checkbox', { name: 'Seleccionar todas' }).check()
  await page.keyboard.press('Control+/')
  const bloque = page.getByTestId('edicion-varias')
  await bloque.locator('select').nth(0).selectOption('kind')
  await bloque.locator('select').nth(1).selectOption('fixed')
  await expect(page.getByTestId('vista-previa')).toContainText('$2,700.00')
  await page.keyboard.press('Escape')
  await expect(bloque).toHaveCount(0)
  expect(parches).toHaveLength(0)
  await page.getByRole('checkbox', { name: 'Seleccionar todas' }).focus()
  await page.keyboard.press('Control+/')
  await bloque.locator('select').nth(0).selectOption('kind')
  await bloque.locator('select').nth(1).selectOption('fixed')
  await expect(page.getByTestId('vista-previa')).toContainText('$2,700.00')
  await page.keyboard.press('Enter')
  await expect.poll(() => parches.length).toBe(2)
  expect(parches.map((p) => p.cuerpo)).toEqual([{ kind: 'fixed' }, { kind: 'fixed' }])
})

test('P-04 N-064 · foto del ticket y «Quitar adjunto» con deleteAttachment', async ({ page }) => {
  await page.goto('/gastos/11')
  await expect(page.getByRole('heading', { name: 'Ticket adjunto' })).toBeVisible()
  await cargada(page.getByRole('img', { name: 'Foto del ticket' }))
  await tocar(page.getByRole('main').getByRole('button', { name: 'Quitar adjunto' }))
  const dialogo = page.getByRole('dialog', { name: 'Quitar adjunto' })
  await tocar(dialogo.getByRole('button', { name: 'Quitar adjunto' }))
  await expect.poll(() => borrados).toEqual(['701'])
  await expect(page.getByText('Ticket quitado del gasto')).toBeVisible()
})

test('P-05 N-072 · emparejar un renglón con el catálogo (listProducts y getProduct)', async ({ page }) => {
  await page.goto('/gastos/ticket?lectura=501')
  await tocar(page.getByRole('button', { name: 'Buscar en catálogo' }).first())
  await page.getByLabel('Buscar en catálogo · Renglón 1').fill('leche')
  await tocar(page.getByRole('button', { name: 'Leche entera' }))
  await tocar(page.getByRole('button', { name: 'Caja de 1 L' }))
  await expect(page.getByText('Leche entera · Caja de 1 L')).toBeVisible()
  await page.getByLabel('Categoría').selectOption({ label: 'Despensa' })
  await tocar(page.getByRole('button', { name: 'Confirmar ticket' }))
  await expect(page).toHaveURL(/\/gastos\/11$/)
  const cuerpo = area.confirmados[0] as { expense: { items: { presentation_id?: number }[] } }
  expect(cuerpo.expense.items[0].presentation_id).toBe(51)
})

test('P-05 N-073 · «Terminar compra» propone lo comprado de la lista de la semana (getShoppingList)', async ({ page }) => {
  await page.goto('/gastos/ticket?semana=2026-W42')
  await page.getByLabel('Elegir de la galería').setInputFiles({ name: 'ticket.png', mimeType: 'image/png', buffer: PNG })
  await expect(page).toHaveURL(/semana=2026-W42.*lectura=501|lectura=501.*semana=2026-W42/)
  await expect(page.getByText('1 artículo marcado como comprado en la lista de la semana 42')).toBeVisible()
  await tocar(page.getByRole('button', { name: 'Agregar 1 renglón' }))
  await expect(page.getByText('Renglones agregados desde la lista de compras')).toBeVisible()
  await expect(page.getByLabel('Renglón 3 · Texto leído')).toHaveValue('Leche entera')
  await page.getByLabel('Categoría').selectOption({ label: 'Despensa' })
  await tocar(page.getByRole('button', { name: 'Confirmar ticket' }))
  await expect(page).toHaveURL(/\/gastos\/11$/)
  const cuerpo = area.confirmados[0] as { expense: { items: { presentation_id?: number; quantity_milli: number }[] } }
  expect(cuerpo.expense.items[2]).toMatchObject({ presentation_id: 51, quantity_milli: 2000 })
})

test('P-08 N-105 · corregir en el lugar la quincena capturada (listIncomes) con «Antes → Después»', async ({ page }) => {
  await page.goto('/presupuesto/ingresos?periodo=2026-10')
  await tocar(page.getByRole('button', { name: /Depósito/ }))
  await page.keyboard.press('ControlOrMeta+a'); await page.keyboard.type('Depósito de nómina'); await page.keyboard.press('Enter')
  await expect.poll(() => parches).toEqual([{ ruta: 'PATCH /api/v1/incomes/31', cuerpo: { notes: 'Depósito de nómina' } }])
  await tocar(page.getByRole('button', { name: /\$15,000\.00/ }).first())
  await page.keyboard.press('ControlOrMeta+a'); await page.keyboard.type('15100')
  await expect(page.getByText('$15,100.00').first()).toBeVisible()
  await page.keyboard.press('Enter')
  await expect.poll(() => parches.length).toBe(2)
  expect(parches[1]).toEqual({ ruta: 'PATCH /api/v1/incomes/31', cuerpo: { amount_cents: 1510000 } })
})

test('P-08 N-102 · ingreso extra con la pieza compartida (createIncome con clave)', async ({ page }) => {
  await page.goto('/presupuesto/ingresos?periodo=2026-10')
  await tocar(page.getByRole('button', { name: 'Registrar ingreso extra' }))
  const hoja = page.getByRole('dialog', { name: 'Registrar ingreso extra' })
  await hoja.getByLabel('Monto').fill('850.5')
  await hoja.getByLabel('Concepto').fill('Venta de libros')
  await tocar(hoja.getByRole('button', { name: 'Guardar' }))
  await expect.poll(() => area.ingresos.length).toBe(1)
  expect(area.ingresos[0].cuerpo).toMatchObject({ kind: 'extra', description: 'Venta de libros', amount_cents: 85050 })
  expect(area.ingresos[0].clave.length).toBeGreaterThanOrEqual(16)
})

test('P-03 N-102 · la pestaña «Ingreso» registra el ingreso en la misma hoja', async ({ page }) => {
  const hoja = await abrirHoja(page, 'Ingreso')
  // R4-01: ya no hay desvío a P-08; la hoja tiene el formulario de N-102 (si vuelve el desvío, la prueba falla, no se salta).
  await expect(hoja.getByRole('button', { name: 'Ir a Ingresos' })).toHaveCount(0)
  await expect(hoja.getByLabel('Monto')).toBeVisible()
  await hoja.getByLabel('Monto').fill('1200')
  await hoja.getByLabel('Concepto').fill('Reembolso')
  await tocar(hoja.getByRole('button', { name: 'Guardar' }))
  await expect.poll(() => area.ingresos.length).toBe(1)
  expect(area.ingresos[0].cuerpo).toMatchObject({ kind: 'extra', description: 'Reembolso', amount_cents: 120000 })
  await expect(page.getByRole('dialog', { name: 'Registrar' })).toHaveCount(0)
})

test('P-03 N-053 · forma de pago «Crédito» elige la línea activa (listCredits) y sigue en P-19 prellenada', async ({ page }) => {
  const hoja = await abrirHoja(page, 'Gasto')
  await hoja.getByLabel('Monto').fill('185.5')
  await tocar(hoja.getByRole('button', { name: 'Más opciones' }))
  await hoja.getByLabel('Forma de pago').selectOption('credit')
  // R4-01: la forma de pago «A crédito» lista las líneas en la hoja (N-053), sin el enlace a P-17.
  await expect(hoja.getByRole('button', { name: 'Ir a compra a crédito' })).toHaveCount(0)
  const linea = hoja.getByLabel('Línea de crédito')
  await expect(linea.locator('option')).toHaveText(['Elige la línea', 'Línea de la tienda · Crédito de tienda', 'Tarjeta azul · Tarjeta bancaria'])
  await linea.selectOption({ label: 'Línea de la tienda · Crédito de tienda' })
  await tocar(hoja.getByRole('button', { name: 'Continuar en Compra a crédito' }))
  await expect(page).toHaveURL(/\/creditos\/2\/compras\/nueva\?monto=185\.50&fecha=\d{4}-\d{2}-\d{2}/)
})
