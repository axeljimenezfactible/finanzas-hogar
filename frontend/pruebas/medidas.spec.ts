// U5-06 · DESIGN §4.6, §5.2 y §3.4 medidos en el DOM a 412 px táctil (proyecto «movil»), pantalla por pantalla de
// Inicio, Gastos y presupuesto, Comidas y compras, y en los estados abiertos del cascarón (hoja F-1 con «Más opciones»,
// panel de la cola, menú «/» y vista previa «Antes → Después»): todo control visible con área de toque ≥ 44 × 44 px,
// todo icono Lucide de 20 o 24 px con trazo de 1.5 px dibujados (24 px en los botones de solo icono) y todo dinero
// visible con cifras tabulares.
import { expect, test, type BrowserContext, type Page } from '@playwright/test'
import { ApiFalsa } from './apoyo/api-falsa'
import { controlesChicos, dineroSinTabulares, iconosFueraDeNorma } from './apoyo/medidas'
import { instalarInicio } from './inicio/datos'
import { AreaFalsa } from './gastos-presupuesto/datos'
import { instalarArea, type Pet } from './comidas-compras/datos'
type Instalar = (ctx: BrowserContext, page: Page) => Promise<void>
const GRUPOS: { nombre: string; instalar: Instalar; rutas: string[] }[] = [
  {
    nombre: 'Inicio',
    instalar: async (ctx, page) => { await new ApiFalsa().instalar(ctx); await instalarInicio(ctx); await page.clock.setFixedTime(new Date('2026-10-10T12:00:00-06:00')) },
    rutas: ['/inicio', '/inicio?sub=detalle', '/inicio?sub=tendencias', '/inicio?sub=calendario'],
  },
  {
    nombre: 'Gastos y presupuesto',
    instalar: async (ctx) => { const api = new ApiFalsa(); await api.instalar(ctx); await new AreaFalsa(api).instalar(ctx) },
    rutas: ['/gastos?periodo=2026-10', '/gastos?periodo=2026-10&sub=resumen', '/gastos?periodo=2026-10&sub=tendencias', '/gastos?periodo=2026-10&sub=calendario', '/gastos/11', '/gastos/ticket',
      '/presupuesto?periodo=2026-10', '/presupuesto?periodo=2026-10&sub=detalle', '/presupuesto?periodo=2026-10&sub=tendencias', '/presupuesto/flujo', '/presupuesto/flujo?sub=detalle',
      '/presupuesto/ingresos?periodo=2026-10', '/presupuesto/ingresos?periodo=2026-10&sub=resumen', '/presupuesto/ahorro', '/presupuesto/ahorro?sub=detalle'],
  },
  {
    nombre: 'Comidas y compras',
    instalar: async (ctx, page) => { await new ApiFalsa().instalar(ctx); await instalarArea(ctx, []); await page.clock.setFixedTime(new Date('2026-10-12T10:00:00-06:00')) },
    rutas: ['/comidas', '/comidas?sub=resumen', '/comidas/recetas', '/comidas/recetas/1', '/compras', '/compras/comparar', '/compras/comparar?sub=detalle', '/compras/productos', '/compras/inventario', '/compras/inventario?sub=resumen'],
  },
]
for (const g of GRUPOS) {
  test(`U5-06 a 412 px táctil · ${g.nombre}: controles ≥ 44 × 44, iconos 20/24 px con trazo 1.5 (24 en solo icono) y dinero tabular`, async ({ page, context }) => {
    test.setTimeout(150_000)
    await g.instalar(context, page)
    const chicos: string[] = [], iconos: string[] = [], dinero: string[] = []
    let controles = 0, totalIconos = 0
    for (const ruta of g.rutas) {
      await page.goto(ruta)
      await expect(page.getByRole('heading', { level: 1 })).toBeVisible()
      await page.waitForLoadState('networkidle')
      await page.waitForTimeout(300)
      const c = await controlesChicos(page), i = await iconosFueraDeNorma(page), d = await dineroSinTabulares(page)
      controles += c.total; totalIconos += i.total
      chicos.push(...c.chicos.map((x) => `${ruta} ${x.control} ${x.ancho}×${x.alto}`))
      iconos.push(...i.mal.map((x) => `${ruta} ${x.regla}: ${x.control} ${x.clase} ${x.ancho}×${x.alto} trazo ${x.trazo}`))
      dinero.push(...d.map((x) => `${ruta} «${x}»`))
    }
    console.log(`U5-06 ${g.nombre} (${g.rutas.length} vistas): controles ${controles}, menores de 44 px ${chicos.length}; iconos ${totalIconos}, fuera de 20/24 px, de trazo 1.5 o de 24 px en solo icono: ${iconos.length}; dinero sin tabulares ${dinero.length}`)
    for (const x of [...chicos, ...iconos, ...dinero]) console.log(`   ${x}`)
    expect(controles).toBeGreaterThan(0)
    expect(totalIconos).toBeGreaterThan(0)
    expect(chicos).toEqual([])
    expect(iconos).toEqual([])
    expect(dinero).toEqual([])
  })
}
test('U5-06 a 412 px táctil · estados abiertos: hoja F-1 con «Más opciones», panel de la cola, menú «/» y vista previa', async ({ page, context }) => {
  test.setTimeout(150_000)
  const api = new ApiFalsa()
  await api.instalar(context); await new AreaFalsa(api).instalar(context)
  const chicos: string[] = [], iconos: string[] = [], dinero: string[] = [], estados: string[] = []
  let controles = 0, totalIconos = 0
  const medir = async (estado: string) => {
    await page.waitForTimeout(300)
    const c = await controlesChicos(page), i = await iconosFueraDeNorma(page), d = await dineroSinTabulares(page)
    controles += c.total; totalIconos += i.total
    estados.push(`${estado} (${c.total} controles, ${i.total} iconos)`)
    chicos.push(...c.chicos.map((x) => `${estado}: ${x.control} ${x.ancho}×${x.alto}`))
    iconos.push(...i.mal.map((x) => `${estado}: ${x.regla}: ${x.control} ${x.clase} ${x.ancho}×${x.alto} trazo ${x.trazo}`))
    dinero.push(...d.map((x) => `${estado}: «${x}»`))
  }
  // Hoja F-1 (P-03): en reposo, con «Más opciones» (todas las categorías y pago a crédito) y en sus otras pestañas.
  await page.goto('/gastos?periodo=2026-10')
  await expect(page.getByRole('heading', { level: 1 })).toBeVisible()
  await expect.poll(() => api.peticiones.some((p) => p.url.startsWith('/api/v1/categories'))).toBe(true)
  await page.getByRole('button', { name: 'Registrar gasto' }).last().tap()
  const hoja = page.getByRole('dialog', { name: 'Registrar' })
  await expect(hoja.getByLabel('Monto')).toBeFocused()
  await medir('F-1 Gasto')
  await hoja.getByRole('button', { name: 'Más opciones' }).tap()
  await expect(hoja.getByRole('button', { name: 'Menos opciones' })).toBeVisible()
  await hoja.getByRole('button', { name: 'Más', exact: true }).tap()
  await hoja.getByLabel('Forma de pago').selectOption('credit')
  await expect(hoja.getByRole('button', { name: 'Ir a compra a crédito' })).toBeVisible()
  await medir('F-1 Gasto · Más opciones')
  for (const pestana of ['Ingreso', 'Peso', 'Ticket']) {
    await hoja.getByRole('tab', { name: pestana }).tap()
    await expect(hoja.getByRole('tab', { name: pestana })).toHaveAttribute('aria-selected', 'true')
    await medir(`F-1 ${pestana}`)
  }
  await hoja.getByRole('button', { name: 'Cerrar' }).tap()
  await expect(hoja).toHaveCount(0)
  // Franja «Sin conexión · 1 en cola», aviso con «Deshacer» y panel de la cola (también al confirmar un descarte).
  await context.setOffline(true)
  await expect(page.getByText('Sin conexión con el servidor de casa')).toBeVisible()
  await page.getByRole('button', { name: 'Registrar gasto' }).last().tap()
  await page.keyboard.type('99')
  await page.getByRole('button', { name: 'Comida fuera' }).tap()
  await page.getByRole('button', { name: 'Guardar' }).tap()
  await expect(page.getByText('1 en cola')).toBeVisible()
  await medir('Franja sin conexión y aviso')
  await page.getByTestId('franja-conexion').getByRole('button', { name: 'Revisar' }).tap()
  const panel = page.getByRole('dialog', { name: 'Capturas en el teléfono' })
  await expect(panel.getByRole('button', { name: 'Descartar' })).toBeVisible()
  await medir('Panel de la cola')
  await panel.getByRole('button', { name: 'Descartar' }).tap()
  await expect(panel.getByRole('button', { name: 'Cancelar' })).toBeVisible()
  await medir('Panel de la cola · confirmar descarte')
  await panel.getByRole('button', { name: 'Cancelar' }).tap()
  await panel.getByRole('button', { name: 'Cerrar' }).tap()
  await expect(panel).toHaveCount(0)
  await context.setOffline(false)
  await expect(page.getByTestId('contador-cola')).toHaveCount(0)
  // Menú «/» (G-03/G-04) con la lista de comandos.
  await page.locator('body').press('/')
  const menu = page.getByRole('dialog', { name: 'Buscar o escribir «/» para comandos' })
  await expect(menu.getByRole('textbox')).toHaveValue('/')
  await expect(menu.getByRole('option').first()).toBeVisible()
  await medir('Menú «/»')
  await page.keyboard.press('Escape')
  await expect(menu).toHaveCount(0)
  // Vista previa «Antes → Después» de la edición en el lugar (P-04, monto) y su estado «No se pudo calcular · Reintentar».
  await page.goto('/gastos/11')
  await expect(page.getByRole('heading', { level: 1, name: 'Gasto' })).toBeVisible()
  await page.getByRole('button', { name: '$1,850.50' }).first().tap()
  const monto = page.getByRole('textbox', { name: 'Monto' })
  await monto.fill('1900')
  await expect(page.getByText('Después')).toBeVisible()
  await medir('Vista previa «Antes → Después»')
  await monto.fill('abc')
  await expect(page.getByText('No se pudo calcular el efecto')).toBeVisible()
  await medir('Vista previa · No se pudo calcular')
  await monto.press('Escape')
  console.log(`U5-06 estados abiertos (${estados.length}): controles ${controles}, menores de 44 px ${chicos.length}; iconos ${totalIconos}, fuera de 20/24 px, de trazo 1.5 o de 24 px en solo icono: ${iconos.length}; dinero sin tabulares ${dinero.length}`)
  for (const x of [...estados, ...chicos, ...iconos, ...dinero]) console.log(`   ${x}`)
  expect(estados.length).toBe(11)
  expect(chicos).toEqual([])
  expect(iconos).toEqual([])
  expect(dinero).toEqual([])
})
test('U5-06 P-13: la casilla «Marcar … como comprado» dibuja 20 px dentro de un área de toque de 44 × 44 que la marca', async ({ page, context }) => {
  const pets: Pet[] = []
  await new ApiFalsa().instalar(context); await instalarArea(context, pets); await page.clock.setFixedTime(new Date('2026-10-12T10:00:00-06:00'))
  await page.goto('/compras')
  const casilla = page.getByRole('checkbox', { name: 'Marcar Leche entera como comprado' }).first()
  await expect(casilla).toBeVisible()
  const dibujo = (await casilla.boundingBox())!
  const area = await casilla.evaluate((el) => { const r = (el as HTMLInputElement).labels![0].getBoundingClientRect(); return { x: r.x, y: r.y, ancho: r.width, alto: r.height } })
  console.log(`U5-06 P-13 casilla: dibujo ${dibujo.width}×${dibujo.height}, área de toque (etiqueta) ${area.ancho}×${area.alto}`)
  expect(area.ancho).toBeGreaterThanOrEqual(44)
  expect(area.alto).toBeGreaterThanOrEqual(44)
  // Tocar la esquina del área, fuera del dibujo de 20 px, marca el artículo (PATCH con la semana).
  expect(area.x + 3 < dibujo.x && area.y + 3 < dibujo.y).toBe(true)
  await page.touchscreen.tap(area.x + 3, area.y + 3)
  await expect.poll(() => pets.find((p) => p.metodo === 'PATCH')?.consulta.get('week')).toBe('2026-W42')
})
