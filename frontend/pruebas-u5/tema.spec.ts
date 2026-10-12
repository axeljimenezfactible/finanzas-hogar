// U5-10 · El tema sigue al sistema (colorScheme emulado, sin tocar el interruptor) y el interruptor lo cambia (P-26.10 y la barra lateral).
import { test, expect } from '@playwright/test'
import fs from 'node:fs'
import path from 'node:path'
import { EVID, abrirSesion, fondoC2, listo } from './apoyo'

test('U5-10 · interruptor del tema: sistema claro/oscuro, forzar claro y oscuro, volver a sistema', async ({ browser }, info) => {
  const disp = info.project.name.split('-')[0]
  const s = await abrirSesion(browser, info.project.use as never)
  const { page } = s
  const lee = () => page.evaluate(() => ({ bg: getComputedStyle(document.body).backgroundColor, dt: document.documentElement.getAttribute('data-theme'), guardado: localStorage.getItem('fh-tema') }))
  const claro = fondoC2('claro').rgb, oscuro = fondoC2('oscuro').rgb
  const pasos: Record<string, unknown>[] = []
  const paso = async (que: string, esperaBg: string, esperaDt: string | null) => {
    const r = await lee(); pasos.push({ paso: que, ...r, esperado_bg: esperaBg, esperado_data_theme: esperaDt })
    expect(r.bg, que).toBe(esperaBg); expect(r.dt, que).toBe(esperaDt)
  }
  try {
    await page.emulateMedia({ colorScheme: 'light' }); await page.goto('/inicio'); await listo(page)
    await paso('sistema=claro, sin tocar el interruptor', claro, null)
    await page.emulateMedia({ colorScheme: 'dark' }); await page.waitForTimeout(300)
    await paso('sistema=oscuro, sin tocar el interruptor', oscuro, null)
    await page.goto('/ajustes/10'); await listo(page)
    const radio = (n: string) => page.getByRole('radio', { name: n })
    await radio('Claro').click(); await page.waitForTimeout(300)
    await paso('interruptor Claro con sistema oscuro', claro, 'light')
    await page.reload(); await listo(page)
    await paso('Claro persiste al recargar', claro, 'light')
    await page.emulateMedia({ colorScheme: 'light' })
    await radio('Oscuro').click(); await page.waitForTimeout(300)
    await paso('interruptor Oscuro con sistema claro', oscuro, 'dark')
    await radio('Sistema').click(); await page.waitForTimeout(300)
    await paso('interruptor Sistema (sistema claro)', claro, null)
    await page.emulateMedia({ colorScheme: 'dark' }); await page.waitForTimeout(300)
    await paso('Sistema sigue al sistema oscuro', oscuro, null)
    if (disp === 'escritorio') {
      const b = page.getByRole('button', { name: /Cambiar tema/ }).first()
      await b.click(); await page.waitForTimeout(300); await paso('barra lateral: sistema → claro', claro, 'light')
      await b.click(); await page.waitForTimeout(300); await paso('barra lateral: claro → oscuro', oscuro, 'dark')
      await b.click(); await page.waitForTimeout(300); await paso('barra lateral: oscuro → sistema (sistema oscuro)', oscuro, null)
    }
  } finally {
    fs.mkdirSync(EVID, { recursive: true })
    fs.writeFileSync(path.join(EVID, `tema-${info.project.name}.json`), JSON.stringify({ criterio: 'U5-10', combinacion: info.project.name, token_claro: fondoC2('claro'), token_oscuro: fondoC2('oscuro'), pasos }, null, 1))
    await s.context.close()
  }
})
