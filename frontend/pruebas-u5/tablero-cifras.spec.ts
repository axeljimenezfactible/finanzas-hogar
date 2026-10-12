// U5-02 sobre P-01 Resumen (tablero de control): guarda el texto visible del tablero junto con las respuestas crudas de /api/v1/dashboard/*
// que lo alimentaron (misma carga), para compararlos entre sí y contra el oráculo (python3 -I frontend/pruebas-u5/tablero_cifras.py).
// Salida: docs/design/capturas-u5/evidencia/tablero-cifras/<combinación>.json (subcarpeta: no cuenta en U5-01). Solo lectura.
import { test as base } from '@playwright/test'
import fs from 'node:fs'
import path from 'node:path'
import { EVID, abrirSesion, listo, type Sesion } from './apoyo'

const test = base.extend<object, { sesion: Sesion }>({
  sesion: [async ({ browser }, use, wi) => { const s = await abrirSesion(browser, wi.project.use as never); await use(s); await s.context.close() }, { scope: 'worker' }],
})
test('P-01 Resumen · cifras en pantalla y respuestas del API', async ({ sesion }, info) => {
  const dir = path.join(EVID, 'tablero-cifras'); fs.mkdirSync(dir, { recursive: true })
  const page = await sesion.context.newPage()
  const api: Record<string, unknown> = {}
  page.on('response', async (r) => {
    const u = new URL(r.url())
    if (u.pathname.startsWith('/api/v1/dashboard/') || u.pathname === '/api/v1/budget/signals') { try { api[u.pathname + u.search] = await r.json() } catch { /* sin cuerpo */ } }
  })
  try {
    await page.goto('/inicio?sub=resumen')
    await listo(page)
    await page.waitForTimeout(1500)
    const texto = await page.evaluate(() => document.body.innerText)
    const aria = await page.evaluate(() => [...document.querySelectorAll('[role=meter],[role=progressbar],[aria-valuetext],[aria-label*="barómetro" i]')].map((e) => ({ role: e.getAttribute('role'), label: e.getAttribute('aria-label'), now: e.getAttribute('aria-valuenow'), text: e.getAttribute('aria-valuetext') })))
    fs.writeFileSync(path.join(dir, `${info.project.name}.json`), JSON.stringify({ combinacion: info.project.name, url: '/inicio?sub=resumen', texto_visible: texto, medidores: aria, api }, null, 1))
  } finally { await page.close() }
})
