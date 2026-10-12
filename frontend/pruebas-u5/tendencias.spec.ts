// U5-06 #14 (DESIGN §7.7, §8.3 y §8.5) · P-01 ▸ Tendencias en las 4 combinaciones (proyecto = combinación): captura de la
// pestaña y comprobación de que ningún texto de las gráficas (etiquetas de datos, ejes, leyenda, «Hoy», «Proyección») se
// encima con otro ni sale del lienzo. Cada caja se reduce 20 % arriba y abajo (el cuadro de la fuente lleva el interlineado:
// el «$0» del eje Y y el primer mes del eje X se tocan 3 px sin que las letras se toquen). Evidencia en docs/design/capturas-u5/evidencia/tendencias/ (subcarpeta: no cuenta en U5-01).
import { test as base, expect } from '@playwright/test'
import fs from 'node:fs'
import path from 'node:path'
import { EVID, abrirSesion, listo, type Sesion } from './apoyo'

const DIR = path.join(EVID, 'tendencias')
const test = base.extend<object, { sesion: Sesion }>({
  sesion: [async ({ browser }, use, wi) => { const s = await abrirSesion(browser, wi.project.use as never); await use(s); await s.context.close() }, { scope: 'worker' }],
})
test('P-01 Tendencias · etiquetas de gráficas sin encimarse', async ({ sesion }, info) => {
  fs.mkdirSync(DIR, { recursive: true })
  const page = await sesion.context.newPage()
  try {
    await page.goto('/inicio?sub=tendencias')
    await listo(page)
    await page.locator('figure[data-prueba^="grafica-"] svg').first().waitFor({ state: 'visible', timeout: 30_000 })
    await page.waitForTimeout(800)
    await page.mouse.move(0, 0)
    const punteroGrueso = await page.evaluate(() => matchMedia('(any-pointer: coarse)').matches)
    await page.screenshot({ path: path.join(DIR, `P-01-tendencias-${info.project.name}.png`), fullPage: true, animations: 'disabled' })
    const figuras = await page.evaluate(() => [...document.querySelectorAll('figure[data-prueba^="grafica-"]')].map((f) => {
      const lienzo = [...f.querySelectorAll('svg')].find((x) => x.querySelector('text')) // el SVG de ECharts (no los iconos de los chips)
      const svg = lienzo?.getBoundingClientRect()
      const cajas = [...(lienzo?.querySelectorAll('text') ?? [])].map((t) => { const r = t.getBoundingClientRect(); const m = r.height * 0.2; return { t: (t.textContent ?? '').trim(), x0: r.left, y0: r.top + m, x1: r.right, y1: r.bottom - m } })
        .filter((b) => b.t && b.x1 - b.x0 > 0 && b.y1 - b.y0 > 0)
      const encimadas: string[] = []
      for (let i = 0; i < cajas.length; i++) for (let j = i + 1; j < cajas.length; j++) {
        const a = cajas[i], b = cajas[j]
        const dx = Math.min(a.x1, b.x1) - Math.max(a.x0, b.x0), dy = Math.min(a.y1, b.y1) - Math.max(a.y0, b.y0)
        if (dx > 1 && dy > 1) encimadas.push(`«${a.t}» × «${b.t}» (${dx.toFixed(1)} × ${dy.toFixed(1)} px)`)
      }
      const fuera = svg ? cajas.filter((b) => b.x0 < svg.left - 0.5 || b.x1 > svg.right + 0.5 || b.y0 < svg.top - 0.5 || b.y1 > svg.bottom + 0.5).map((b) => b.t) : []
      return { figura: f.getAttribute('data-prueba'), titulo: f.querySelector('figcaption span')?.textContent ?? '', textos: cajas.map((b) => b.t), encimadas, fuera_del_lienzo: fuera }
    }))
    fs.writeFileSync(path.join(DIR, `P-01-tendencias-${info.project.name}.json`), JSON.stringify({ combinacion: info.project.name, ruta: '/inicio?sub=tendencias', puntero_grueso: punteroGrueso, figuras }, null, 1))
    expect(figuras.length, 'gráficas dibujadas en Tendencias').toBeGreaterThan(0)
    if (info.project.name.startsWith('android')) expect(punteroGrueso).toBe(true)
    expect(figuras.flatMap((f) => f.encimadas.map((e) => `${f.titulo}: ${e}`)), 'textos encimados').toEqual([])
    expect(figuras.flatMap((f) => f.fuera_del_lienzo.map((e) => `${f.titulo}: ${e}`)), 'textos fuera del lienzo').toEqual([])
  } finally { await page.close() }
})
