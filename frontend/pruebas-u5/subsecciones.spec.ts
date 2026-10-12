// U5-01, U5-03, U5-04 y U5-10 para P-01 Inicio en sus 4 subsecciones (Resumen, Detalle, Tendencias, Calendario) × 4 combinaciones.
// capturas.spec.ts toma P-01 en su subsección por omisión (Resumen); aquí van las otras tres y también el Resumen, con axe y contraste
// propios de cada subsección. Salida en docs/design/capturas-u5/evidencia/inicio-subsecciones/ (subcarpeta: no cuenta en U5-01).
// Uso (desde frontend/): npx playwright test -c <config con testMatch /subsecciones\.spec\.ts/> (ver docs/estado/wf3/u5-tablero.md).
import { test as base, expect } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'
import fs from 'node:fs'
import path from 'node:path'
import { AQUI, EVID, ETIQUETAS_AXE, abrirSesion, fondoC2, listo, type Sesion } from './apoyo'

const CONTRASTE_JS = fs.readFileSync(path.join(AQUI, 'contraste-dom.js'), 'utf8')
const DIR = path.join(EVID, 'inicio-subsecciones')
const SUBS = ['resumen', 'detalle', 'tendencias', 'calendario'] as const
const test = base.extend<object, { sesion: Sesion }>({
  sesion: [async ({ browser }, use, wi) => { const s = await abrirSesion(browser, wi.project.use as never); await use(s); await s.context.close() }, { scope: 'worker' }],
})
for (const sub of SUBS) {
  test(`P-01 ${sub} · captura, axe y contraste`, async ({ sesion }, info) => {
    fs.mkdirSync(DIR, { recursive: true })
    const [disp, tema] = info.project.name.split('-') as ['android' | 'escritorio', 'claro' | 'oscuro']
    const nombre = `P-01-${sub}-${disp}-${tema}`
    const page = await sesion.context.newPage()
    const hosts = new Set<string>()
    page.on('request', (r) => { try { const u = new URL(r.url()); if (/^(https?|wss?):$/.test(u.protocol)) hosts.add(u.host) } catch { /* data: / blob: */ } })
    try {
      await page.goto(`/inicio?sub=${sub}`)
      await listo(page)
      if (sub === 'tendencias') await page.locator('figure[data-prueba^="grafica-"] svg').first().waitFor({ state: 'visible', timeout: 30_000 })
      await page.waitForTimeout(900)
      await page.mouse.move(0, 0)
      const punteroGrueso = await page.evaluate(() => matchMedia('(any-pointer: coarse)').matches)
      await page.screenshot({ path: path.join(DIR, `${nombre}.png`), fullPage: true, animations: 'disabled' })
      const axe = await new AxeBuilder({ page: page as never }).withTags(ETIQUETAS_AXE).analyze()
      const serias = axe.violations.filter((v) => v.impact === 'serious' || v.impact === 'critical')
      const cont = await page.evaluate(`(${CONTRASTE_JS})(null)`) as { evaluados: number; bajo_umbral: number; indeterminados: number; [k: string]: unknown }
      const esperado = fondoC2(tema)
      const dom = await page.evaluate(() => ({ bg: getComputedStyle(document.body).backgroundColor, dataTheme: document.documentElement.getAttribute('data-theme'), texto: document.body.innerText, ancho: document.documentElement.scrollWidth, alto: document.documentElement.scrollHeight, url: location.pathname + location.search }))
      fs.writeFileSync(path.join(DIR, `${nombre}.json`), JSON.stringify({
        pantalla: 'P-01', subseccion: sub, combinacion: info.project.name, url: dom.url, dimensiones_pagina: { ancho: dom.ancho, alto: dom.alto }, puntero_grueso: punteroGrueso,
        tema: { esperado: tema, fondo_body: dom.bg, fondo_token_c2: esperado.rgb, coincide: dom.bg === esperado.rgb, data_theme: dom.dataTheme },
        axe: { violaciones: axe.violations.length, serias_o_criticas: serias.length, incompletos: axe.incomplete.length, reglas: axe.violations.map((v) => ({ id: v.id, impact: v.impact, nodos: v.nodes.map((n) => n.target.join(' ')) })) },
        contraste: cont, hosts: [...hosts].sort(), texto_visible: dom.texto,
      }, null, 1))
      expect(serias.map((v) => `${v.id}: ${v.nodes.map((n) => n.target.join(' ')).join(' | ')}`), 'axe serias/críticas').toEqual([])
      expect(cont.bajo_umbral, 'contraste: pares bajo umbral').toBe(0)
      expect(dom.bg, 'fondo del body = token de C2').toBe(esperado.rgb)
      expect(dom.dataTheme).toBeNull()
      expect(dom.ancho, 'sin desborde horizontal').toBeLessThanOrEqual(disp === 'android' ? 412 : 1440)
      expect([...hosts].filter((h) => !/^(finanzas\.local|localhost|127\.0\.0\.1)(:\d+)?$/.test(h)), 'hosts externos').toEqual([])
    } finally { await page.close() }
  })
}
