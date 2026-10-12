// C3-08 · Flujos complementarios de U1 §4: F-7 (compra a crédito con el comparador, obligatorio (b) de U1-05) y F-8
// (ticket con OCR, obligatorio (d)) de extremo a extremo contra el servidor REAL, a 412×915 táctil, sin API simulada.
// Escriben datos: corren contra la COPIA de prueba (servidor-prueba.sh → https://127.0.0.1:8444, base en datos/tmp/u5-flujos/).
//   bash pruebas-u5/flujos/servidor-prueba.sh arrancar && FH_BASE=https://127.0.0.1:8444 npx playwright test -c pruebas-u5/flujos/playwright.config.ts f7-f8
// Toques (K10): SOLO pulsaciones en controles (abrir, elegir, guardar); teclear y los valores por defecto no cuentan.
// F-8 usa la imagen SINTÉTICA datos/semilla/tickets-sinteticos/sintetico-01.png (no es foto real ni de una tienda) y el
// OCR local del servidor (Tesseract, spa; job «ticket_ocr», revisión cada 30 s).
import path from 'node:path'
import { expect, test } from '@playwright/test'
import { api, entrar, Flujo, RAIZ } from './apoyo'

test.describe.configure({ mode: 'default' })
const norm = (s: string) => s.replace(/\s+/g, ' ').trim()
const cuerpo = async (page: import('@playwright/test').Page) => norm(await page.locator('main').innerText())

test('F-7 · Planear una compra a crédito con el comparador (P-01 → G-01 → P-17 → P-20 → P-19)', async ({ page }) => {
  test.setTimeout(180_000)
  await entrar(page)
  // ── Parte A: con la semilla tal cual (sin restricciones globales en P-26.05): K16 «Faltan restricciones · Capturar»
  const f = new Flujo('F-7A', 'Compra a crédito con el comparador · con la semilla (sin restricciones, K16)', page)
  await f.paso('Partida (sin tocar): las restricciones globales de P-26.05 en la semilla (API getOptimizerSettings)', 'API', async () => {
    const r = await api<Record<string, unknown>>(page, '/api/v1/settings/optimizer')
    expect(r.estado).toBe(200)
    const { min_available_buffer_cents: colchon, max_live_purchases: vivas, zero_debt_by_date: cero } = r.json
    return `colchón mínimo ${JSON.stringify(colchon)} · máximo de compras vivas ${JSON.stringify(vivas)} · deuda en cero ${JSON.stringify(cero)}`
  })
  await f.paso('Toque 1: ☰ (Abrir menú)', 'P-01 → G-01', async () => {
    await f.t.tocar(page.getByRole('button', { name: 'Abrir menú' }), '☰')
    await expect(page.getByRole('link', { name: 'Créditos', exact: true })).toBeVisible()
  })
  await f.paso('Toque 2: «Créditos»', 'G-01 → P-17', async () => {
    await f.t.tocar(page.getByRole('link', { name: 'Créditos', exact: true }), 'Créditos')
    await expect(page).toHaveURL(/\/creditos(\?|$)/)
    await expect(page.getByRole('button', { name: 'Simular compra' })).toBeVisible()
    return (await cuerpo(page)).slice(0, 200)
  })
  await f.paso('Toque 3: «Simular compra» (U1: «⋯» → «Simular compra», 2 toques; aquí el botón está a la vista en la barra de acciones)', 'P-17 → P-20', async () => {
    await f.t.tocar(page.getByRole('button', { name: 'Simular compra' }), 'Simular compra')
    await expect(page).toHaveURL(/\/creditos\/comparador/)
    await expect(page.getByRole('heading', { level: 1, name: 'Comparador de escenarios' })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Comparar escenarios' })).toBeVisible()
    return (await cuerpo(page)).slice(0, 300)
  })
  await f.paso('Teclear el monto 3094.50 (no cuenta; total de la compra del 10-oct-2026 de la semilla)', 'P-20', async () => {
    await page.getByLabel('Monto', { exact: true }).fill('3094.50')
  })
  let estadoA = ''
  await f.paso('Toque 4: «Comparar escenarios» → «Faltan restricciones · Capturar», sin recomendación (K16)', 'P-20', async () => {
    const resp = page.waitForResponse((r) => r.url().includes('/api/v1/scenario-comparisons') && r.request().method() === 'POST')
    await f.t.tocar(page.getByRole('button', { name: 'Comparar escenarios' }), 'Comparar escenarios')
    const j = await (await resp).json()
    estadoA = j.status
    expect(j.status, 'status de la API').toBe('missing_constraints')
    expect(j.recommendation ?? null, 'sin recomendación').toBeNull()
    await expect(page.getByText(/Faltan restricciones/)).toBeVisible()
    await expect(page.getByRole('link', { name: 'Capturar', exact: true })).toBeVisible()
    await expect(page.getByRole('button', { name: 'Usar escenario recomendado' })).toHaveCount(0)
    await page.getByText(/Faltan restricciones/).scrollIntoViewIfNeeded()
    await f.captura('F7-faltan-restricciones')
    return `API: status ${j.status}, faltan ${JSON.stringify(j.missing_constraints)}, ${j.evaluated_count} evaluados · ${j.valid?.length} válidos · ${j.rejected?.length} descartados, recomendación null · pantalla: ${(await cuerpo(page)).match(/[0-9]+ escenarios evaluados[^.]*|Faltan restricciones[^.]*\./g)?.join(' | ')}`
  })
  await f.paso('Toque 5: «Capturar» abre Ajustes › Créditos y condiciones', 'P-20 → P-26.05', async () => {
    await f.t.tocar(page.getByRole('link', { name: 'Capturar', exact: true }), 'Capturar')
    await expect(page).toHaveURL(/\/ajustes\/05/)
    await expect(page.getByRole('heading', { name: 'Créditos y condiciones' }).first()).toBeVisible()
    return (await cuerpo(page)).slice(0, 200)
  })
  const regA = f.guardar()

  // ── Parte B: el resto de la cadena de U1 (pasos 4 a 7) con restricciones tecleadas SOLO para esta comparación
  // (valores de prueba; las globales de P-26.05 no cambian: se comprueba al final).
  const g = new Flujo('F-7B', 'Compra a crédito con el comparador · hasta P-19 con restricciones de esta comparación', page)
  await g.paso('Volver al comparador (atrás del navegador; no cuenta: el recorrido corto de U1 es ☰ → Créditos → Simular compra)', 'P-26.05 → P-20', async () => {
    await page.goBack()
    await expect(page.getByRole('heading', { level: 1, name: 'Comparador de escenarios' })).toBeVisible()
  })
  await g.paso('Teclear el monto 3094.50 (no cuenta)', 'P-20', async () => { await page.getByLabel('Monto', { exact: true }).fill('3094.50') })
  await g.paso('Toque 1: abrir «Restricciones duras (solo para esta comparación)» y teclear valores de prueba: colchón 500.00, máximo 3 compras vivas y deuda en cero el 2027-12-31 (el servidor exige las 3, C1; la pantalla dice «opcional» para la fecha: ver C3-flujos punto 6)', 'P-20', async () => {
    await g.t.tocar(page.getByText('Restricciones duras (solo para esta comparación)'), 'Restricciones duras')
    await page.getByLabel('Colchón mínimo de disponible').fill('500.00')
    await page.getByLabel('Máximo de compras vivas a la vez').fill('3')
    await page.getByLabel(/^Fecha de deuda en cero/).fill('2027-12-31')
  })
  let rec: { credit_id: number; term_count: number; purchase_date: string; interest_cents: number } | null = null
  await g.paso('Toque 2: «Comparar escenarios» → recomendación con su puntaje y los 4 óptimos apilados', 'P-20', async () => {
    const resp = page.waitForResponse((r) => r.url().includes('/api/v1/scenario-comparisons') && r.request().method() === 'POST')
    await g.t.tocar(page.getByRole('button', { name: 'Comparar escenarios' }), 'Comparar escenarios')
    const j = await (await resp).json()
    expect(j.status, `status de la API (faltan ${JSON.stringify(j.missing_constraints)})`).toBe('ok')
    rec = j.recommendation
    expect(rec).not.toBeNull()
    await expect(page.getByText(/Escenario recomendado · puntaje/)).toBeVisible()
    await expect(page.getByText('Óptimo de cada objetivo')).toBeVisible()
    await page.getByText(/Escenario recomendado · puntaje/).scrollIntoViewIfNeeded()
    await g.captura('F7-recomendacion')
    return `API: ${j.evaluated_count} evaluados · ${j.valid.length} válidos · ${j.rejected.length} descartados; recomendado: línea ${rec!.credit_id}, ${rec!.term_count} pagos, compra ${rec!.purchase_date}, intereses ${rec!.interest_cents} centavos · pantalla: ${(await cuerpo(page)).match(/Escenario recomendado · puntaje [0-9.]+ de 100/)?.[0]}`
  })
  await g.paso('Toque 3: «Usar escenario recomendado» → P-19 prellenada, sin guardar nada', 'P-20 → P-19', async () => {
    const antes = await api<{ items: unknown[] }>(page, '/api/v1/credit-purchases?limit=200')
    await g.t.tocar(page.getByRole('button', { name: 'Usar escenario recomendado' }), 'Usar escenario recomendado')
    await expect(page).toHaveURL(/\/creditos\/\d+\/compras\/nueva\?/)
    await expect(page.getByLabel('Capital', { exact: true })).toHaveValue('3094.50')
    await expect(page.getByLabel('Plazo (número de pagos)')).toHaveValue(String(rec!.term_count))
    await expect(page.getByLabel('Fecha de compra')).toHaveValue(rec!.purchase_date)
    const despues = await api<{ items: unknown[] }>(page, '/api/v1/credit-purchases?limit=200')
    expect(despues.json.items.length, 'P-19 prellenada no guarda nada').toBe(antes.json.items.length)
    await expect.poll(async () => cuerpo(page), { message: 'vista previa del cálculo (N-216)' }).toMatch(/Vista previa del cálculo.*Capital \$3,094\.50.*Intereses \(IVA incl\.\) \$[0-9,]+\.[0-9]{2}.*Total \$[0-9,]+\.[0-9]{2}/)
    await g.captura('F7-compra-prellenada')
    const previa = (await cuerpo(page)).match(/Vista previa del cálculo.{0,160}/)?.[0]
    return `URL ${new URL(page.url()).pathname}${new URL(page.url()).search}; ${previa}; compras antes ${antes.json.items.length} = después ${despues.json.items.length}; intereses vacíos → «${await page.getByText(/Vacío: se usa la tasa/).first().innerText()}»`
  })
  let compraId = 0
  await g.paso('Toque 4: «Guardar compra» (intereses vacíos: tasa de la tabla por plazo de la línea, de la semilla) → crea la compra con sus pagos', 'P-19', async () => {
    const resp = page.waitForResponse((r) => /\/api\/v1\/credit-purchases$/.test(new URL(r.url()).pathname) && r.request().method() === 'POST')
    await g.t.tocar(page.getByRole('button', { name: 'Guardar compra' }), 'Guardar compra')
    const r = await resp
    expect(r.status()).toBe(201)
    compraId = (await r.json()).id
    await expect(page.getByText('Compra guardada')).toBeVisible()
    await expect(page).toHaveURL(new RegExp(`/creditos/\\d+/compras/${compraId}$`))
    const c = await api<{ principal_cents: number; interest_cents: number; term_count: number; purchase_date: string }>(page, `/api/v1/credit-purchases/${compraId}`)
    expect(c.estado).toBe(200)
    expect(c.json.principal_cents).toBe(309450)
    expect(c.json.term_count).toBe(rec!.term_count)
    const pagos = await api<{ items: { amount_cents: number; due_date: string }[] }>(page, `/api/v1/scheduled-payments?limit=200&credit_purchase_id=${compraId}`)
    await g.captura('F7-compra-guardada')
    return `compra ${compraId}: capital ${c.json.principal_cents}, intereses ${c.json.interest_cents}, ${c.json.term_count} pagos, fecha ${c.json.purchase_date}; pagos programados (API) ${pagos.estado} · ${pagos.json?.items?.length ?? '?'} renglones`
  })
  await g.paso('Las restricciones globales de P-26.05 siguen sin capturar (las de esta comparación no se guardan)', 'API', async () => {
    const r = await api<Record<string, unknown>>(page, '/api/v1/settings/optimizer')
    expect(r.json.min_available_buffer_cents ?? null).toBeNull()
    expect(r.json.max_live_purchases ?? null).toBeNull()
    return 'colchón y máximo de compras vivas globales: null'
  })
  const regB = g.guardar()
  expect(estadoA).toBe('missing_constraints')
  expect(regA.ok, JSON.stringify(regA.pasos.filter((p) => !p.ok))).toBe(true)
  expect(regB.ok, JSON.stringify(regB.pasos.filter((p) => !p.ok))).toBe(true)
})

test('F-8 · Foto de ticket → OCR → confirmación → precios guardados (P-01 → P-03 → P-05 → P-04)', async ({ page }) => {
  test.setTimeout(240_000)
  await entrar(page)
  const f = new Flujo('F-8', 'Ticket con OCR', page)
  const foto = path.join(RAIZ, 'datos/semilla/tickets-sinteticos/sintetico-01.png')
  await f.paso('Toque 1: «+» de la barra inferior', 'P-01 → P-03', async () => {
    await f.t.tocar(page.getByTestId('boton-mas'), '«+» barra inferior')
    await expect(page.getByRole('tab', { name: 'Ticket' })).toBeVisible()
  })
  await f.paso('Toque 2: «Ticket» en el selector de tipo', 'P-03', async () => {
    await f.t.tocar(page.getByRole('tab', { name: 'Ticket' }), 'tipo Ticket')
    await expect(page.getByRole('button', { name: 'Capturar ticket' })).toBeVisible()
    return norm(await page.getByRole('dialog').innerText().catch(() => '')).slice(0, 200)
  })
  await f.paso('Toque 3: «Capturar ticket» (U1: «Ticket» abre directo la cámara; aquí la hoja lleva a P-05)', 'P-03 → P-05', async () => {
    await f.t.tocar(page.getByRole('button', { name: 'Capturar ticket' }), 'Capturar ticket')
    await expect(page).toHaveURL(/\/gastos\/ticket/)
    await expect(page.getByRole('button', { name: 'Tomar foto' })).toBeVisible()
  })
  await f.paso('Toques 4-5: «Tomar foto» y aceptar la foto en la cámara (imagen sintética de prueba)', 'P-05', async () => {
    const elegir = page.waitForEvent('filechooser')
    await f.t.tocar(page.getByRole('button', { name: 'Tomar foto' }), 'Tomar foto')
    await (await elegir).setFiles(foto)
    f.t.n++; f.t.pasos.push(`${f.t.n}. aceptar la foto (cámara)`)
    await expect(page.getByText('Leyendo el ticket…').first()).toBeVisible()
    await expect(page).toHaveURL(/lectura=\d+/)
    return `«Leyendo el ticket…» visible; ${new URL(page.url()).search}`
  })
  const lecturaId = Number(new URL(page.url()).searchParams.get('lectura'))
  let propuesta = ''
  await f.paso('El OCR local lee el ticket y propone los renglones (espera al job, ≤ 120 s)', 'P-05', async () => {
    await expect(page.getByText('Encabezado', { exact: true })).toBeVisible({ timeout: 120_000 })
    await expect(page.getByLabel(/^Renglón 1 · Texto leído/)).toBeVisible()
    const s = await api<{ ocr_status: string; lines: { line_no: number; text: string; line_total_cents: number; confidence: string; presentation_id?: number }[]; proposed_total_cents: number | null; proposed_store_id: number | null; proposed_expense_date: string | null }>(page, `/api/v1/ticket-scans/${lecturaId}`)
    expect(s.json.ocr_status).toBe('done')
    expect(s.json.lines.length).toBeGreaterThan(0)
    propuesta = s.json.lines.map((l) => `${l.text} ${l.line_total_cents} (${l.confidence}${l.presentation_id ? `, pres. ${l.presentation_id}` : ''})`).join(' | ')
    await f.captura('F8-renglones-propuestos')
    return `OCR ${s.json.ocr_status}: ${s.json.lines.length} renglones: ${propuesta}; total propuesto ${s.json.proposed_total_cents}, tienda ${s.json.proposed_store_id}, fecha ${s.json.proposed_expense_date}`
  })
  const renglones = page.locator('li').filter({ has: page.getByLabel(/^Renglón \d+ · Texto leído/) })
  await f.paso('Corrección de 1 renglón (U1 paso 5, opcional): el de «HUEVO BLANCO 12 PZA» → Huevo blanco, cartón 12 piezas (busca en el catálogo)', 'P-05', async () => {
    let i = -1
    const n = await renglones.count()
    for (let k = 0; k < n && i < 0; k++) if (/HUEVO/i.test(await renglones.nth(k).getByLabel(/^Renglón \d+ · Texto leído/).inputValue())) i = k
    expect(i, `renglón de HUEVO entre ${n}`).toBeGreaterThanOrEqual(0)
    const fila = renglones.nth(i)
    if (/cartón 12 piezas/.test(await fila.innerText())) return `renglón ${i + 1}: ya venía emparejado por el OCR, sin corrección`
    await f.t.tocar(fila.getByRole('button', { name: 'Buscar en catálogo' }), 'Buscar en catálogo')
    await fila.getByRole('searchbox').fill('huevo blanco')
    await f.t.tocar(fila.getByRole('button', { name: 'Huevo blanco', exact: true }), 'producto Huevo blanco')
    await f.t.tocar(fila.getByRole('button', { name: /^Huevo blanco, cartón 12 piezas/ }), 'presentación cartón 12 piezas')
    await expect(fila).toContainText('Huevo blanco · Huevo blanco, cartón 12 piezas')
    return `renglón ${i + 1} («${await fila.getByLabel(/^Renglón \d+ · Texto leído/).inputValue()}») emparejado con la presentación 116`
  })
  await f.paso('Corrección: quitar los renglones sin importe que el OCR leyó del encabezado del ticket (1 toque cada uno)', 'P-05', async () => {
    const quitados: string[] = []
    for (let vuelta = 0; vuelta < 10; vuelta++) {
      const n = await renglones.count()
      let i = -1
      for (let k = 0; k < n && i < 0; k++) if (!(await renglones.nth(k).getByLabel('Importe', { exact: true }).inputValue())) i = k
      if (i < 0) break
      quitados.push(await renglones.nth(i).getByLabel(/^Renglón \d+ · Texto leído/).inputValue())
      await f.t.tocar(page.getByRole('button', { name: `Quitar renglón ${i + 1}`, exact: true }), `Quitar renglón ${i + 1}`)
    }
    return `${quitados.length} quitados: ${quitados.map((x) => `«${x}»`).join(', ')}; quedan ${await renglones.count()}`
  })
  await f.paso('Encabezado: tienda y categoría (selectores nativos: abrir + elegir = 2 toques cada uno; Bodega Aurrera y Despensa); fecha y total se teclean si el OCR no los propuso (no cuentan)', 'P-05', async () => {
    const tienda = page.getByRole('combobox', { name: /^Tienda o comercio/ })
    if (!(await tienda.inputValue())) { await f.t.tocar(tienda, 'selector de tienda'); await tienda.selectOption({ label: 'Bodega Aurrera' }); f.t.n++; f.t.pasos.push(`${f.t.n}. opción «Bodega Aurrera»`) }
    const cat = page.getByRole('combobox', { name: /^Categoría/ })
    if (!(await cat.inputValue())) { await f.t.tocar(cat, 'selector de categoría'); await cat.selectOption({ label: 'Despensa' }); f.t.n++; f.t.pasos.push(`${f.t.n}. opción «Despensa»`) }
    const campoFecha = page.getByLabel('Fecha', { exact: true }).first()
    const fecha = await campoFecha.inputValue()
    if (!fecha) await campoFecha.fill('2026-10-12')
    const campoTotal = page.getByLabel('Total del ticket', { exact: true })
    const total = await campoTotal.inputValue()
    if (!total) await campoTotal.fill('251.53') // suma de los 6 importes impresos del ticket sintético (esperado/sintetico-01.csv)
    await expect(page.getByText('La suma de los renglones cuadra con el total.')).toBeVisible()
    return `tienda ${await tienda.inputValue()} · categoría ${await cat.inputValue()} · fecha propuesta ${fecha || '(vacía → tecleada 2026-10-12)'} · total propuesto ${total || '(vacío → tecleado 251.53)'} · la suma cuadra`
  })
  let gastoId = 0
  let precios: number[] = []
  await f.paso('Toque final: «Confirmar ticket» → crea gasto, renglones, precios con fuente ticket e inventario; abre el gasto', 'P-05 → P-04', async () => {
    const resp = page.waitForResponse((r) => /\/ticket-scans\/\d+\/confirm$/.test(new URL(r.url()).pathname) && r.request().method() === 'POST')
    await f.t.tocar(page.getByRole('button', { name: 'Confirmar ticket' }), 'Confirmar ticket')
    const r = await resp
    const j = await r.json()
    expect(r.status(), JSON.stringify(j)).toBe(201)
    gastoId = j.expense_id; precios = j.price_ids
    await expect(page).toHaveURL(new RegExp(`/gastos/${gastoId}$`))
    await expect(page.getByText('Renglones del ticket').or(page.getByText('Producto')).first()).toBeVisible()
    await f.captura('F8-gasto-creado')
    return `gasto ${gastoId}; precios ${JSON.stringify(j.price_ids)}; inventario ${JSON.stringify(j.inventory_item_ids)} · P-04: ${(await cuerpo(page)).slice(0, 300)}`
  })
  await f.paso('El precio quedó guardado (API listProductPrices del producto Huevo blanco, fuente ticket_ocr, $45.90 en la tienda elegida)', 'API', async () => {
    expect(precios.length).toBeGreaterThan(0)
    const r = await api<{ items: { id: number; presentation_id: number; price_cents: number; source: string; store_id: number; observed_date: string }[] }>(page, '/api/v1/products/109/prices?limit=200')
    expect(r.estado).toBe(200)
    const p = r.json.items.find((x) => precios.includes(x.id))
    expect(p, JSON.stringify(r.json.items.slice(0, 3))).toBeTruthy()
    expect(p!.presentation_id).toBe(116)
    expect(p!.price_cents).toBe(4590)
    expect(p!.source).toBe('ticket_ocr')
    return `precio ${p!.id}: presentación ${p!.presentation_id}, ${p!.price_cents} centavos, fuente ${p!.source}, tienda ${p!.store_id}, fecha ${p!.observed_date}`
  })
  await f.paso('El gasto está en la tabla de Todos los gastos', 'P-02', async () => {
    const g = await api<{ amount_cents: number; category_id: number; store_id: number | null; items?: unknown[] }>(page, `/api/v1/expenses/${gastoId}`)
    expect(g.estado).toBe(200)
    return `gasto ${gastoId}: ${g.json.amount_cents} centavos, categoría ${g.json.category_id}, tienda ${g.json.store_id}`
  })
  const reg = f.guardar()
  expect(reg.ok, JSON.stringify(reg.pasos.filter((p) => !p.ok))).toBe(true)
})
