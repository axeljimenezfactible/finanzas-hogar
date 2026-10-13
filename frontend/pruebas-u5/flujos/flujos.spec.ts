// C3-08 · U5-09 · Los 6 flujos clave de U1 §4 (F-1…F-6) de extremo a extremo contra el servidor REAL, a 412×915 táctil.
// Corre contra la COPIA de prueba (servidor-prueba.sh → https://127.0.0.1:8444, base en datos/tmp/u5-flujos/) porque
// F-1, F-2, F-3, F-4, F-5 y F-6 escriben datos; la base de desarrollo no se toca.
//   bash servidor-prueba.sh arrancar && cd ../.. && FH_BASE=https://127.0.0.1:8444 npx playwright test -c pruebas-u5/flujos/playwright.config.ts flujos
// Toques (K10): SOLO pulsaciones en controles (abrir, elegir, guardar); teclear y los valores por defecto no cuentan.
import { expect, test } from '@playwright/test'
import { api, entrar, Flujo } from './apoyo'

test.describe.configure({ mode: 'default' }) // sigue tras un fallo: cada flujo informa el suyo
const hoyLargo = /12-oct-2026/
const norm = (s: string) => s.replace(/\s+/g, ' ').trim()

test('F-1 · Registrar gasto desde Inicio (U5-09: ≤ 3 toques, el gasto aparece en Inicio y en la tabla)', async ({ page }) => {
  await entrar(page)
  const f = new Flujo('F-1', 'Registrar gasto desde Inicio', page)
  let previo = ''
  // U3 2.0 (K54): el gasto del mes lo muestra K-1 «Gasto del mes» del tablero (antes, la tarjeta 1 «Mes contra plan»).
  await f.paso('Partida: app abierta en Inicio (Resumen) con sesión (no cuenta toques)', 'P-01', async () => {
    await expect(page.getByTestId('tarjeta-K-1')).toHaveJSProperty('tagName', 'A')
    previo = norm(await page.getByTestId('tarjeta-K-1').innerText())
    return `K-1 «Gasto del mes» antes: ${previo.slice(0, 120)}`
  })
  await f.paso('Toque 1: «+» de la barra inferior', 'P-01 → P-03', async () => {
    await f.t.tocar(page.getByTestId('boton-mas'), '«+» barra inferior')
    await expect(page.getByLabel('Monto')).toBeVisible()
    await expect(page.getByLabel('Monto')).toBeFocused()
    const chips = await page.getByRole('button', { name: /^(Despensa|Comida fuera|Papel)$/ }).count()
    return `hoja abierta, monto enfocado, ${chips} categorías sugeridas visibles; tipo Gasto por defecto`
  })
  await f.paso('Teclear 185.50 (no cuenta)', 'P-03', async () => { await page.getByLabel('Monto').fill('185.50') })
  await f.paso('Toque 2: chip de categoría «Despensa»', 'P-03', async () => {
    await f.t.tocar(page.getByRole('button', { name: 'Despensa', exact: true }), 'chip Despensa')
  })
  await f.paso('Toque 3: «Guardar»', 'P-03 → P-01', async () => {
    await f.captura('F1-hoja-antes-de-guardar')
    await f.t.tocar(page.getByRole('button', { name: 'Guardar', exact: true }), 'Guardar')
    await expect(page.getByText('Gasto guardado')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Deshacer' })).toBeVisible()
    return 'aviso «Gasto guardado · Deshacer»'
  })
  await f.paso('El gasto aparece en Inicio (K-1 «Gasto del mes» y Detalle → Presupuesto por categoría)', 'P-01', async () => {
    await expect.poll(async () => norm(await page.getByTestId('tarjeta-K-1').innerText()), { message: 'K-1 «Gasto del mes» cambió' }).not.toBe(previo)
    const despues = norm(await page.getByTestId('tarjeta-K-1').innerText())
    await page.getByRole('tab', { name: 'Detalle' }).tap()
    await expect(page.getByText('Presupuesto por categoría')).toBeVisible()
    await expect(page.getByText('$1,290.39 de $2,000.00')).toBeVisible() // 1,197.64 + 50 % de 185.50: Despensa se reparte Conjunto 50/50 y el presupuesto cuenta la parte del hogar
    await f.captura('F1-inicio-con-el-gasto')
    return `K-1 después: ${despues.slice(0, 120)} · Detalle: Despensa $1,290.39 de $2,000.00 (antes $1,197.64; +$92.75 = parte del hogar del reparto 50/50 de $185.50)`
  })
  await f.paso('El gasto aparece en la vista de tabla de Todos los gastos', 'P-02', async () => {
    await page.goto('/gastos?sub=detalle')
    const fila = page.getByRole('row').filter({ hasText: '$185.50' }).filter({ hasText: 'Despensa' })
    await expect(fila.first()).toBeVisible()
    await expect(fila.first()).toContainText(hoyLargo)
    await f.captura('F1-tabla-todos-los-gastos')
    return `fila: ${norm(await fila.first().innerText())}`
  })
  await f.paso('El servidor tiene el gasto (API)', 'API getExpenses', async () => {
    const r = await api<{ items: { amount_cents: number }[] }>(page, '/api/v1/expenses?limit=200')
    expect(r.estado).toBe(200)
    const lista = r.json.items
    const n = lista.filter((g) => g.amount_cents === 18550).length
    expect(n).toBe(1)
    return '1 gasto de 18550 centavos'
  })
  const reg = f.guardar()
  expect(f.t.n, `toques: ${f.t.pasos.join(' · ')}`).toBeLessThanOrEqual(3)
  expect(reg.ok, JSON.stringify(reg.pasos.filter((p) => !p.ok))).toBe(true)
})

test('F-3A0 · «Mejor canasta» con la lista vacía (estado de la semilla, antes de F-2)', async ({ page }) => {
  await entrar(page)
  const f = new Flujo('F-3A0', 'Comparar canasta · ruta A con la lista vacía (semilla tal cual)', page)
  await f.paso('Toque 1: tarjeta «Mejor canasta» (dice que compara la compra del 10-oct-2026)', 'P-01 → P-14', async () => {
    await expect(page.getByTestId('tarjeta-5')).toHaveJSProperty('tagName', 'A')
    const antes = norm(await page.getByTestId('tarjeta-5').innerText())
    await f.t.tocar(page.getByTestId('tarjeta-5'), 'tarjeta Mejor canasta')
    await expect(page.getByRole('heading', { level: 1, name: 'Comparar canasta' })).toBeVisible()
    await page.waitForTimeout(1500)
    const t = norm(await page.locator('main').innerText())
    await f.captura('F3A0-comparar-con-lista-vacia')
    expect(t, `la tarjeta decía «${antes}»; P-14 debía mostrar la comparación de esa canasta`).toContain('más barata')
    return `P-14: ${t.slice(0, 300)}`
  })
  f.guardar()
})

test('F-2 · Planear la semana de comidas', async ({ page }) => {
  await entrar(page)
  const f = new Flujo('F-2', 'Planear la semana de comidas', page)
  await f.paso('Toque 1: «Comidas» de la barra inferior', 'P-01 → P-10', async () => {
    await f.t.tocar(page.getByRole('link', { name: 'Comidas', exact: true }), 'Comidas (barra inferior)')
    await expect(page.getByRole('heading', { level: 1, name: 'Plan de la semana' })).toBeVisible()
    await expect(page.getByText(/Semana 42/)).toBeVisible()
    return 'plan en Detalle, semana 42 (12 al 18-oct-2026), semana vigente'
  })
  await f.paso('Toque 2: casilla vacía Lun · Desayuno', 'P-10', async () => {
    await f.t.tocar(page.getByRole('button', { name: 'Desayuno: Vacía' }).first(), 'casilla Lun · Desayuno')
    await expect(page.getByRole('button', { name: 'Fuera de casa' })).toBeVisible()
    return 'selector de recetas abierto (tiempo Desayuno, día Lun 12 marcado)'
  })
  await f.paso('Toque 3: receta «Avena con manzana y canela»', 'P-10', async () => {
    await f.t.tocar(page.getByRole('button', { name: /^Avena con manzana y canela/ }), 'receta Avena con manzana y canela')
    await expect(page.getByRole('button', { name: 'Fuera de casa' })).toBeHidden()
    await expect(page.getByRole('button', { name: /^Desayuno: Axel y Ximena.*Avena con manzana/ }).first()).toBeVisible()
    return 'la casilla muestra «Axel y Ximena · Avena con manzana y canela»'
  })
  await f.paso('Toque 4: «Generar lista de compras»', 'P-10 → P-13', async () => {
    await f.captura('F2-plan-con-receta')
    await f.t.tocar(page.getByRole('button', { name: 'Generar lista de compras' }), 'Generar lista de compras')
    await expect(page).toHaveURL(/\/compras(\?|$)/)
    await expect(page.getByRole('heading', { level: 1, name: 'Lista de compras' })).toBeVisible()
    await page.waitForTimeout(1500)
    const t = norm(await page.locator('main').innerText())
    expect(t).not.toMatch(/lista está vacía|Aún no hay/i)
    await f.captura('F2-lista-generada')
    return `lista abierta: ${t.slice(0, 260)}`
  })
  const reg = f.guardar()
  expect(reg.ok, JSON.stringify(reg.pasos.filter((p) => !p.ok))).toBe(true)
})

test('F-3 · Comparar tienda para una canasta (ruta A: condición de U1 §4 con la lista de F-2; ruta B: otra canasta)', async ({ page }) => {
  await entrar(page)
  // U1 §4 F-3 (aclarado por u3-v2, opción (c) de C3-flujos punto 2, PROVISIONAL hasta la puerta 4): la ruta A exige una lista
  // vigente con artículos y con TODOS sus precios en al menos una tienda. La lista que deja F-2 (artículos de recetas sin precio en
  // la semilla) no la cumple: U1 manda que «Mejor canasta» muestre su vacío y lleve a P-15 «Registrar precios», no a P-14. Eso se
  // comprueba aquí; la ruta A completa (3 toques, con una lista que sí cumple) la corre F-3A3 en pendientes.spec.ts.
  const f = new Flujo('F-3A', 'Comparar canasta · ruta A: condición de U1 §4 con la lista que generó F-2', page)
  await f.paso('Condición de la ruta A (no cuenta toques): la lista de F-2 no tiene todos sus precios en ninguna tienda → «Mejor canasta» en vacío con «Registrar precios» → P-15', 'P-01', async () => {
    const c = page.getByTestId('tarjeta-5')
    await expect(c).toHaveJSProperty('tagName', 'A')
    await expect(c).toContainText('Aún no hay precios para comparar')
    await expect(c).toContainText('Registrar precios')
    const href = await c.getAttribute('href')
    expect(href, 'la acción del vacío lleva a P-15 Productos y precios, no a P-14').toMatch(/^\/compras\/productos(\?|$)/)
    const lista = await api<{ groups: { items: unknown[] }[] }>(page, '/api/v1/shopping-list?week=2026-W42')
    expect(lista.estado).toBe(200)
    const n = lista.json.groups.flatMap((g) => g.items).length
    expect(n, 'F-2 dejó una lista con artículos').toBeGreaterThan(0)
    await f.captura('F3A-tarjeta-sin-precios')
    return `lista de la semana con ${n} artículos; tarjeta: «${norm(await c.innerText())}» → ${href}. Ruta A completa: F-3A3 (pendientes.spec.ts)`
  })
  const regA = f.guardar()

  // Ruta B: Compras → ⋯ → Comparar tiendas → selector de canasta → «Compra del 10-oct-2026»
  await page.goto('/inicio')
  await expect(page.getByTestId('tarjeta-5')).toBeVisible()
  const g = new Flujo('F-3B', 'Comparar canasta · ruta B (compra del 10-oct-2026)', page)
  await g.paso('Toque 1: «Compras» de la barra inferior', 'P-01 → P-13', async () => {
    await g.t.tocar(page.getByRole('link', { name: 'Compras', exact: true }), 'Compras (barra inferior)')
    await expect(page.getByRole('heading', { level: 1, name: 'Lista de compras' })).toBeVisible()
  })
  await g.paso('Toques 2-3: «⋯» → «Comparar tiendas»', 'P-13 → P-14', async () => {
    await g.t.tocar(page.getByRole('button', { name: 'Más opciones' }), '⋯ Más opciones')
    await g.t.tocar(page.getByText('Comparar tiendas', { exact: true }), 'Comparar tiendas')
    await expect(page.getByRole('heading', { level: 1, name: 'Comparar canasta' })).toBeVisible()
  })
  await g.paso('Toques 4-5: selector de canasta → «Compra del 10-oct-2026»', 'P-14', async () => {
    // Selector nativo: abrirlo (1 toque) y elegir la opción (1 toque).
    await g.t.tocar(page.locator('#p14-canasta'), 'selector de canasta')
    await page.locator('#p14-canasta').selectOption('1')
    g.t.n++; g.t.pasos.push(`${g.t.n}. opción «Compra del 10-oct-2026»`)
    await page.waitForTimeout(2000)
    const t = norm(await page.locator('main').innerText())
    await g.captura('F3B-comparar-compra-10-oct')
    return `Resumen: ${t.slice(0, 700)}`
  })
  // Cifra de la ruta B (correccion-c3-r2): con `basket = expense` la API cotiza cada otra tienda con los productos y las
  // cantidades de la compra base (K46, Q2-47 = A, FIRME; S5-interfaz §7.2 «API»: «la API no reproduce +138.50 con la
  // semilla»). El ejemplo de U1 1.2 §4 F-3 (Walmart $3,233.00, +$138.50) es el canon carrito contra carrito, que
  // reproducen CompareCarts y la fixture s5-comparativa-canon, no la pantalla. Se comprueba pantalla = API, y en idénticos +$100.50 (= U1).
  await g.paso('Cifra FIRME (K46, Q2-47 = A): la pantalla da los totales de getBasketComparison con la compra del 10-oct-2026 como canasta base', 'P-14 = API', async () => {
    const t = norm(await page.locator('main').innerText())
    const r = await api<{ cheapest_store_id?: number; stores: { store_id: number; total_cents: number; difference_vs_cheapest_cents?: number | null; identical_total_cents?: number; missing_count: number }[] }>(page, '/api/v1/basket-comparison?basket=expense&expense_id=1&week=2026-W42')
    expect(r.estado, JSON.stringify(r.json)).toBe(200)
    const mx = (c: number) => '$' + (c / 100).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
    const barata = r.json.stores.find((s) => s.store_id === r.json.cheapest_store_id)
    expect(barata, 'hay tienda más barata').toBeTruthy()
    expect(barata!.total_cents, 'U1 §4 F-3: Bodega Aurrera $3,094.50').toBe(309450)
    const otras = r.json.stores.filter((s) => s.store_id !== barata!.store_id && s.missing_count === 0)
    expect(otras.length, 'al menos otra tienda completa').toBeGreaterThan(0)
    expect(t).toContain(mx(barata!.total_cents))
    const partes: string[] = []
    for (const s of otras) {
      expect(t, `total de la tienda ${s.store_id}`).toContain(mx(s.total_cents))
      expect(t, `diferencia de la tienda ${s.store_id}`).toContain('+' + mx(s.difference_vs_cheapest_cents ?? 0))
      const idem = (s.identical_total_cents ?? 0) - (barata!.identical_total_cents ?? 0)
      expect(t, `idénticos de la tienda ${s.store_id}`).toContain(`+${mx(idem)} en idénticos`)
      partes.push(`tienda ${s.store_id}: ${mx(s.total_cents)}, +${mx(s.difference_vs_cheapest_cents ?? 0)}, +${mx(idem)} en idénticos`)
    }
    expect(t, 'U1 §4 F-3: en idénticos +$100.50').toContain('+$100.50 en idénticos')
    return `pantalla = API: más barata ${mx(barata!.total_cents)} (tienda ${barata!.store_id}); ${partes.join('; ')}. U1 1.2 dice Walmart $3,233.00 (+$138.50): canon de CompareCarts, que la API no reproduce por Q2-47 = A`
  })
  const regB = g.guardar()
  expect(regA.ok, JSON.stringify(regA.pasos.filter((p) => !p.ok))).toBe(true)
  expect(regB.ok, JSON.stringify(regB.pasos.filter((p) => !p.ok))).toBe(true)
})

test('F-4 · Ver un crédito y su calendario de pagos', async ({ page }) => {
  await entrar(page)
  const f = new Flujo('F-4', 'Ver un crédito y su calendario de pagos', page)
  // U1 §4 F-4 paso 1 (U3 2.0, K54): se toca el pago en K-4 «Pagos próximos»; cada fila es un enlace a P-18 de su línea.
  await f.paso('Toque 1: el pago del 30-oct-2026 en K-4 «Pagos próximos»', 'P-01 → P-18', async () => {
    const fila = page.getByTestId('tarjeta-K-4').getByRole('link', { name: /^Vie 30 oct 2026/ })
    await expect(fila).toHaveCount(1)
    await f.t.tocar(fila, 'pago del 30-oct-2026 en K-4')
    await expect(page.getByRole('heading', { level: 1, name: 'Crédito' })).toBeVisible()
    await expect(page.getByText('Disponible · lectura A')).toBeVisible()
    const t = norm(await page.locator('main').innerText())
    for (const c of ['Deuda viva', '$10,000.00', '$6,905.50', 'Pendiente de confirmar|pendiente de confirmar', '$973.75']) expect(t).toMatch(new RegExp(c.replace(/\$/g, '\\$')))
    return 'Resumen: deuda viva, disponible A $10,000.00 y B $6,905.50 con aviso «pendiente de confirmar», próximo pago $973.75'
  })
  await f.paso('Toque 2: pestaña «Calendario»', 'P-18', async () => {
    await f.t.tocar(page.getByRole('tab', { name: 'Calendario' }), 'pestaña Calendario')
    const fila = page.getByRole('button', { name: /^Vie 30 oct 2026/ }).first()
    await expect(fila).toBeVisible()
    const t = norm(await fila.innerText())
    expect(t).toMatch(/\$973\.75/); expect(t).toMatch(/Hogar \$598\.82/); expect(t).toMatch(/Karen \$374\.93/); expect(t).toMatch(/pendiente/)
    await f.captura('F4-calendario')
    return `primera fila: ${t}`
  })
  await f.paso('Toque 3: tocar el pago (vista previa)', 'P-18 → G-06', async () => {
    await f.t.tocar(page.getByRole('button', { name: /^Vie 30 oct 2026/ }).first(), 'pago del 30-oct-2026')
    await page.waitForTimeout(1200)
    const t = norm(await page.locator('body').innerText())
    expect(t).toMatch(/Pago programado/)
    await f.captura('F4-vista-previa-pago')
    return 'vista previa: ' + norm(await page.locator('body').innerText()).split('Pago programado').slice(1).join('Pago programado').slice(0, 700)
  })
  const reg = f.guardar()
  expect(reg.ok, JSON.stringify(reg.pasos.filter((p) => !p.ok))).toBe(true)
})

test('F-5 · Ver o generar un reporte', async ({ page }) => {
  await entrar(page)
  const f = new Flujo('F-5', 'Ver o generar un reporte', page)
  await f.paso('Toque 1: ☰ (Abrir menú)', 'P-01 → G-01', async () => {
    await f.t.tocar(page.getByRole('button', { name: 'Abrir menú' }), '☰')
    await expect(page.getByRole('link', { name: 'Reportes', exact: true })).toBeVisible()
  })
  await f.paso('Toque 2: «Reportes»', 'G-01 → P-23', async () => {
    await f.t.tocar(page.getByRole('link', { name: 'Reportes', exact: true }), 'Reportes')
    await expect(page.getByRole('heading', { level: 1, name: 'Centro de reportes' })).toBeVisible()
    return norm(await page.locator('main').innerText()).slice(0, 400)
  })
  await f.paso('Toque 3: «Generar reporte mensual» (no hay fila «Mensual · sep-2026» que abrir)', 'P-23', async () => {
    await f.t.tocar(page.getByRole('button', { name: 'Generar reporte mensual' }), 'Generar reporte mensual')
    await page.waitForTimeout(2500)
    await f.captura('F5-generar')
    return norm(await page.locator('body').innerText()).slice(0, 500)
  })
  const reg = f.guardar()
  expect(reg.ok, JSON.stringify(reg.pasos.filter((p) => !p.ok))).toBe(true)
})

test('F-6 · Capturar peso', async ({ page }) => {
  await entrar(page)
  const f = new Flujo('F-6', 'Capturar peso', page)
  await f.paso('Toque 1: «+» de la barra inferior', 'P-01 → P-03', async () => {
    await f.t.tocar(page.getByTestId('boton-mas'), '«+» barra inferior')
    await expect(page.getByRole('tab', { name: 'Peso' })).toBeVisible()
  })
  await f.paso('Toque 2: «Peso» en el selector de tipo', 'P-03 → P-22', async () => {
    await f.t.tocar(page.getByRole('tab', { name: 'Peso' }), 'tipo Peso')
    await expect(page.getByLabel(/Peso/).first()).toBeVisible()
    const t = norm(await page.locator('body').innerText())
    return `campo de peso; ${/Antes[^.]*|Sin peso previo[^.]*/.exec(t)?.[0] ?? 'sin texto «Antes»'}`
  })
  await f.paso('Teclear 78.4 (no cuenta)', 'P-22', async () => { await page.getByLabel(/^Peso/).first().fill('78.4') })
  await f.paso('Toque 3: «Guardar»', 'P-22 → P-01', async () => {
    await f.t.tocar(page.getByRole('button', { name: 'Guardar', exact: true }), 'Guardar')
    await expect(page.getByText('Peso guardado')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Deshacer' })).toBeVisible()
    return 'aviso «Peso guardado · Deshacer»'
  })
  await f.paso('Salud muestra el peso (IMC e historial)', 'P-21', async () => {
    await page.goto('/salud?sub=detalle')
    await page.waitForTimeout(2000)
    const t = norm(await page.locator('main').innerText())
    expect(t).toContain('78.4')
    await f.captura('F6-salud')
    return t.slice(0, 400)
  })
  const reg = f.guardar()
  expect(f.t.n, `toques: ${f.t.pasos.join(' · ')}`).toBe(3)
  expect(reg.ok, JSON.stringify(reg.pasos.filter((p) => !p.ok))).toBe(true)
})
