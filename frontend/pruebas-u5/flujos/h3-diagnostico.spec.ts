// WF-3b · continuación 1 · Diagnóstico de H-3 (cobertura 0 de 31 tras F-2). Corre contra la COPIA (8444). Escribe en la copia.
import { writeFileSync } from 'node:fs'
import path from 'node:path'
import { test } from '@playwright/test'
import { api, DIR_TMP, entrar } from './apoyo'

test('H-3 diagnóstico', async ({ page }) => {
  await entrar(page)
  await page.waitForURL(/\/inicio/); await page.waitForTimeout(800)
  const sal: Record<string, unknown> = {}
  const gen = await api<any>(page, '/api/v1/shopping-list/generate?week=2026-W42', { method: 'POST', headers: { 'Idempotency-Key': 'h3-diag-0000000000000001' }, body: {} })
  sal.generar = { estado: gen.estado, json: gen.json }
  const lista = await api<any>(page, '/api/v1/shopping-list?week=2026-W42')
  sal.lista_estado = lista.estado
  const items = (lista.json.groups ?? []).flatMap((g: any) => g.items)
  sal.n_items = items.length
  sal.items_muestra = items.slice(0, 4)
  sal.con_producto = items.filter((i: any) => i.product_id != null).length
  sal.con_cantidad = items.filter((i: any) => i.quantity_milli > 0 || i.needed_milli > 0).length
  sal.con_precio = items.filter((i: any) => i.best_price).length
  sal.total_estimado = lista.json.estimated_total_cents ?? lista.json.total_cents
  sal.claves_lista = Object.keys(lista.json)
  const mc = await api<any>(page, '/api/v1/dashboard/best-basket?week=2026-W42')
  sal.mejor_canasta = mc.json
  const cmp = await api<any>(page, '/api/v1/basket-comparison?week=2026-W42&basket=shopping_list')
  sal.comparacion_estado = cmp.estado
  sal.comparacion = cmp.estado === 200 ? { cheapest: cmp.json.cheapest_store, stores: (cmp.json.stores ?? []).map((s: any) => ({ s: s.store_id, tot: s.total_cents, id: s.identical_count, su: s.substitute_count, mi: s.missing_count })) } : cmp.json
  writeFileSync(path.join(DIR_TMP, 'h3.json'), JSON.stringify(sal, null, 1))
})

test('H-3 detalle', async ({ page }) => {
  await entrar(page); await page.waitForURL(/\/inicio/); await page.waitForTimeout(800)
  const lista = await api<any>(page, '/api/v1/shopping-list?week=2026-W42')
  const items = (lista.json.groups ?? []).flatMap((g: any) => g.items)
  const cmp = await api<any>(page, '/api/v1/basket-comparison?week=2026-W42&basket=shopping_list')
  const tiendas = await api<any>(page, '/api/v1/stores')
  const prods = await api<any>(page, '/api/v1/products?limit=300')
  const nombre = new Map<number, string>((prods.json.items ?? []).map((p: any) => [p.id, `${p.name}|best=${p.best_price ? p.best_price.price_cents + '@' + p.best_price.store_id : '-'}`]))
  writeFileSync(path.join(DIR_TMP, 'h3b.json'), JSON.stringify({
    tiendas: (tiendas.json.items ?? tiendas.json).map?.((t: any) => ({ id: t.id, n: t.name })),
    claves_item: Object.keys(items[0]), claves_cmp: Object.keys(cmp.json), cheapest: cmp.json.cheapest_store ?? null,
    store2: (cmp.json.stores ?? []).find((s: any) => s.store_id === 2), 
    items: items.map((i: any) => ({ p: i.product_id, nom: nombre.get(i.product_id), need: i.needed_milli, pk: i.packages_to_buy, pres: i.recommended_presentation_id, est: i.estimated_cost_cents ?? i.price_cents ?? null })),
    lineas_cmp: (cmp.json.lines ?? cmp.json.items ?? []).slice(0, 30),
  }, null, 1))
})
