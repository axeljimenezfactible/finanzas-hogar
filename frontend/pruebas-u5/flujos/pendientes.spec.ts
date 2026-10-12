// WF-3b · u5-flujos · CONTINUACIÓN 1: lo que quedó pendiente de la primera pasada. Corre contra la COPIA de prueba
// (servidor-prueba.sh → https://127.0.0.1:8444; escribe datos): F-5 pasos 3-5 (generar y descargar el PDF), F-4 paso 4
// (Marcar pago), F-2 paso 2 (Copiar semana anterior) y F-3 ruta A repetida con una lista que sí tiene precios (H-3).
//   FH_BASE=https://127.0.0.1:8444 npx playwright test -c pruebas-u5/flujos/playwright.config.ts pendientes
import { readFileSync, statSync } from 'node:fs'
import { expect, test } from '@playwright/test'
import { api, entrar, Flujo } from './apoyo'

test.describe.configure({ mode: 'default' })
const norm = (s: string) => s.replace(/\s+/g, ' ').trim()
const cuerpo = async (page: import('@playwright/test').Page) => norm(await page.locator('body').innerText())

test('F-5 · Reporte: generar, abrir la vista web y descargar el PDF (U1 pasos 1-5)', async ({ page }) => {
  await entrar(page)
  const f = new Flujo('F-5b', 'Ver o generar un reporte (completo)', page)
  await f.paso('Toque 1: ☰', 'P-01 → G-01', async () => { await f.t.tocar(page.getByRole('button', { name: 'Abrir menú' }), '☰'); await expect(page.getByRole('link', { name: 'Reportes', exact: true })).toBeVisible() })
  await f.paso('Toque 2: «Reportes»', 'G-01 → P-23', async () => { await f.t.tocar(page.getByRole('link', { name: 'Reportes', exact: true }), 'Reportes'); await expect(page.getByRole('heading', { level: 1, name: 'Centro de reportes' })).toBeVisible() })
  await f.paso('Toques 3-4: «Generar reporte mensual» → «Generar» (mes oct-2026 parcial, sin Telegram)', 'P-23', async () => {
    await f.t.tocar(page.getByRole('button', { name: 'Generar reporte mensual' }), 'Generar reporte mensual')
    await expect(page.getByLabel(/Enviarlo por Telegram/)).not.toBeChecked()
    await f.t.tocar(page.getByRole('button', { name: 'Generar', exact: true }), 'Generar')
    await expect(page.getByText('Mensual PDF').first()).toBeVisible({ timeout: 20000 })
    await f.captura('F5b-fila-generada')
    return 'aparece la fila «Mensual PDF · oct-2026» con «Sin envío», «Descargar PDF» y «Compartir»'
  })
  await f.paso('Toque 5: tocar la fila «Mensual · oct-2026»', 'P-23 → P-24', async () => {
    await f.t.tocar(page.getByRole('button', { name: /^Mensual/ }).first(), 'fila Mensual oct-2026')
    await expect(page.getByRole('heading', { level: 1, name: 'Reporte mensual' })).toBeVisible()
    await expect(page).toHaveURL(/\/reportes\/mensual\?mes=2026-10/)
    await page.waitForTimeout(800)
    await f.captura('F5b-vista-web-del-mes')
    return 'vista web: ' + norm(await page.locator('main').innerText()).slice(0, 260)
  })
  let descargado = ''
  await f.paso('Toque 6: «Descargar PDF» en la vista del mes', 'P-24', async () => {
    const espera = page.waitForEvent('download', { timeout: 20000 })
    await f.t.tocar(page.getByRole('button', { name: 'Descargar PDF' }).first(), 'Descargar PDF')
    const d = await espera
    const ruta = await d.path()
    const bytes = readFileSync(ruta)
    expect(bytes.subarray(0, 5).toString('latin1'), 'el archivo es un PDF').toBe('%PDF-')
    descargado = `${d.suggestedFilename()} · ${statSync(ruta).size} bytes · cabecera ${bytes.subarray(0, 8).toString('latin1').trim()}`
    return descargado
  })
  await f.paso('Contrato: GET del PDF con la sesión (tipo y tamaño)', 'API', async () => {
    const r = await page.evaluate(async () => {
      const lista = await (await fetch('/api/v1/reports?limit=5', { credentials: 'same-origin' })).json()
      const it = (lista.items ?? lista)[0]
      return { it: it ? { id: it.id, type: it.type ?? it.report_type, period: it.period ?? it.period_label, claves: Object.keys(it) } : null }
    })
    return JSON.stringify(r)
  }, false, true)
  const reg = f.guardar()
  expect(reg.ok, JSON.stringify(reg.pasos.filter((p) => !p.ok))).toBe(true)
})

test('F-4 · Marcar pago (U1 paso 4) y recálculo del disponible', async ({ page }) => {
  await entrar(page)
  const f = new Flujo('F-4b', 'Marcar pago de la parte del hogar', page)
  let antes = ''
  await f.paso('Toques 1-3: tarjeta Próximo pago → Calendario → pago del 30-oct', 'P-01 → P-18 → G-06', async () => {
    await f.t.tocar(page.getByTestId('tarjeta-3'), 'tarjeta Próximo pago')
    await expect(page.getByText('Disponible · lectura A')).toBeVisible()
    antes = norm(await page.locator('main').innerText())
    await f.t.tocar(page.getByRole('tab', { name: 'Calendario' }), 'pestaña Calendario')
    await f.t.tocar(page.getByRole('button', { name: /^Vie 30 oct 2026/ }).first(), 'pago del 30-oct')
    await expect(page.getByText('Pago programado').first()).toBeVisible()
  })
  await f.paso('Toque 4: «Marcar pagada la parte de Hogar»', 'G-06', async () => {
    await f.t.tocar(page.getByRole('button', { name: 'Marcar pagada la parte de Hogar' }), 'Marcar pagada la parte de Hogar')
    await expect(page.getByText('Parte marcada como pagada')).toBeVisible()
    await f.captura('F4b-parte-pagada')
    const fila = norm(await page.getByRole('button', { name: /^Vie 30 oct 2026/ }).first().innerText())
    expect(fila).toMatch(/pagada \/ pendiente/)
    return `la fila queda «${fila}» y aparece «Deshacer» (el marcado es inmediato: no hay «Guardar»; U1 dice 6 toques, aquí son 4)`
  })
  await f.paso('Recálculo: Resumen del crédito después de marcar', 'P-18', async () => {
    await page.getByRole('tab', { name: 'Resumen' }).tap()
    await page.waitForTimeout(1200)
    const despues = norm(await page.locator('main').innerText())
    await f.captura('F4b-resumen-despues')
    const disp = (t: string) => (t.match(/Disponible[^$]{0,30}\$[0-9,]+\.[0-9]{2}/g) ?? []).join(' | ') + ' · Próximo pago ' + (t.match(/Próximo pago \$[0-9,.]+ [^$]{0,60}\$[0-9,.]+ · Karen \$[0-9,.]+/)?.[0] ?? '?')
    await page.reload(); await expect(page.getByText('Disponible · lectura A')).toBeVisible(); await page.waitForTimeout(800)
    const recargado = norm(await page.locator('main').innerText())
    return `antes → ${disp(antes)} || después → ${disp(despues)} || tras recargar → ${disp(recargado)} · deuda viva antes/después/recargado: ${antes.match(/Deuda viva \$[0-9,.]+/)?.[0]} / ${despues.match(/Deuda viva \$[0-9,.]+/)?.[0]} / ${recargado.match(/Deuda viva \$[0-9,.]+/)?.[0]}`
  })
  await f.paso('Deshacer (la copia vuelve a la semilla)', 'P-18', async () => {
    await page.getByRole('tab', { name: 'Calendario' }).tap()
    const fila = norm(await page.getByRole('button', { name: /^Vie 30 oct 2026/ }).first().innerText())
    const d = page.getByRole('button', { name: 'Deshacer' })
    if (await d.count()) { await d.tap(); await page.waitForTimeout(1000) }
    return `fila antes de deshacer: ${fila}; después: ${norm(await page.getByRole('button', { name: /^Vie 30 oct 2026/ }).first().innerText())}`
  }, false, true)
  const reg = f.guardar()
  expect(reg.ok, JSON.stringify(reg.pasos.filter((p) => !p.ok))).toBe(true)
})

test('F-2 · Copiar semana anterior (U1 paso 2: ⋯ → Copiar semana anterior = +2 toques)', async ({ page }) => {
  await entrar(page)
  const f = new Flujo('F-2b', 'Planear la semana · copiar semana anterior', page)
  await f.paso('Toque 1: Comidas', 'P-01 → P-10', async () => { await f.t.tocar(page.getByRole('link', { name: 'Comidas', exact: true }), 'Comidas'); await expect(page.getByRole('heading', { level: 1, name: 'Plan de la semana' })).toBeVisible() })
  await f.paso('Toques 2-3: «⋯» → «Copiar semana anterior»', 'P-10', async () => {
    await f.t.tocar(page.getByRole('button', { name: 'Más opciones' }), '⋯ Más opciones')
    const resp = page.waitForResponse((r) => r.url().includes('/meal-plan/copy-previous-week'))
    await f.t.tocar(page.getByRole('menuitem', { name: 'Copiar semana anterior' }), 'Copiar semana anterior')
    const r = await resp
    const cuerpoResp = await r.text()
    await page.waitForTimeout(500)
    await f.captura('F2b-copiar-semana')
    const t = await cuerpo(page)
    return `POST copy-previous-week → ${r.status()} ${cuerpoResp.slice(0, 200)} · aviso en pantalla: ${t.match(/(Se copiaron|copiad|No hay|nada que copiar|Ya hay)[^.]{0,120}/i)?.[0] ?? '(ninguno visible)'}`
  })
  const reg = f.guardar()
  expect(reg.ok, JSON.stringify(reg.pasos.filter((p) => !p.ok))).toBe(true)
})

test('F-3 ruta A repetida (H-3): lista de la semana con y sin precios en todas las tiendas', async ({ page }) => {
  await entrar(page)
  const f = new Flujo('F-3A2', 'Comparar canasta · ruta A con la lista del plan', page)
  await f.paso('Preparación (no cuenta): generar la lista de la semana 42 desde el plan de la semilla', 'API', async () => {
    const g = await api<any>(page, '/api/v1/shopping-list/generate?week=2026-W42', { method: 'POST', headers: { 'Idempotency-Key': 'u5-f3a2-generar-0000000001' }, body: {} })
    expect(g.estado).toBe(200)
    await page.goto('/inicio'); await expect(page.getByTestId('tarjeta-5')).toHaveJSProperty('tagName', 'A'); await page.waitForTimeout(1500)
    return `lista con ${g.json.groups.flatMap((x: any) => x.items).length} artículos; tarjeta: «${norm(await page.getByTestId('tarjeta-5').innerText())}»`
  })
  await f.paso('Observación A1: la tarjeta y P-14 con la lista completa (¿hay tienda con todos los precios?)', 'P-01 → P-14', async () => {
    await f.t.tocar(page.getByTestId('tarjeta-5'), 'tarjeta Mejor canasta (A1)')
    await page.waitForTimeout(1500)
    await f.captura('F3A2-destino-de-la-tarjeta')
    const h1 = norm(await page.getByRole('heading', { level: 1 }).first().innerText())
    return `la tarjeta lleva a «${h1}» (${new URL(page.url()).pathname}), NO a Comparar canasta (P-14): con ${'la lista completa'} ninguna tienda tiene todos los precios`
  })
  await f.paso('Preparación (no cuenta; datos de prueba por API en la COPIA): registrar un precio de $50.00 en Vallejo para cada artículo de la lista sin precio (equivale a «Registrar precios»)', 'API', async () => {
    const cmp = await api<any>(page, '/api/v1/basket-comparison?week=2026-W42&basket=shopping_list')
    const sinPrecio = new Set<number>(cmp.json.items.filter((l: any) => l.prices.length === 0).map((l: any) => l.product_id))
    const lista = await api<any>(page, '/api/v1/shopping-list?week=2026-W42')
    let n = 0
    for (const it of lista.json.groups.flatMap((g: any) => g.items)) {
      if (!sinPrecio.has(it.product_id)) continue
      const r = await api(page, '/api/v1/prices', { method: 'POST', headers: { 'Idempotency-Key': `u5-f3a2-precio-${it.recommended_presentation_id}-0000` }, body: { presentation_id: it.recommended_presentation_id, store_id: 2, price_cents: 5000, observed_date: '2026-10-10' } })
      expect(r.estado, JSON.stringify(r.json)).toBe(201); n++
    }
    const c2 = await api<any>(page, '/api/v1/basket-comparison?week=2026-W42&basket=shopping_list')
    return `${n} precios registrados en Vallejo; comparación: ${JSON.stringify(c2.json.stores.map((x: any) => ({ tienda: x.store_id, total: x.total_cents, faltan: x.missing_count })))}; cheapest_store: ${c2.json.cheapest_store ?? 'ausente'}`
  })
  const g = new Flujo('F-3A3', 'Comparar canasta · ruta A (lista con precios)', page)
  await g.paso('Partida: Inicio con la lista que sí tiene precios', 'P-01', async () => {
    await page.goto('/inicio'); await expect(page.getByTestId('tarjeta-5')).toHaveJSProperty('tagName', 'A'); await page.waitForTimeout(1500)
    return `tarjeta: «${norm(await page.getByTestId('tarjeta-5').innerText())}»`
  })
  await g.paso('Toque 1: tarjeta «Mejor canasta»', 'P-01 → P-14', async () => {
    await g.t.tocar(page.getByTestId('tarjeta-5'), 'tarjeta Mejor canasta')
    await expect(page.getByRole('heading', { level: 1, name: 'Comparar canasta' })).toBeVisible(); await page.waitForTimeout(1500)
    await g.captura('F3A3-p14-resumen')
    return norm(await page.locator('main').innerText()).slice(0, 400)
  })
  await g.paso('Toque 2: pestaña «Detalle»', 'P-14', async () => { await g.t.tocar(page.getByRole('tab', { name: 'Detalle' }), 'pestaña Detalle'); await page.waitForTimeout(1000) })
  await g.paso('Toque 3: «Usar <tienda ganadora>»', 'P-14 → P-13', async () => {
    const b = page.getByRole('button', { name: /^Usar / }).first()
    await expect(b).toBeVisible({ timeout: 4000 })
    const nombre = norm(await b.innerText())
    await g.t.tocar(b, nombre)
    await expect(page).toHaveURL(/\/compras(\?|$)/); await page.waitForTimeout(1200)
    await g.captura('F3A3-lista-con-tienda')
    return `asignada: ${nombre} · P-13: ${norm(await page.locator('main').innerText()).slice(0, 300)}`
  })
  const r1 = f.guardar(); const r2 = g.guardar()
  expect(r1.pasos[0].ok && r2.ok, JSON.stringify([...r1.pasos, ...r2.pasos].filter((p) => !p.ok))).toBe(true)
})
