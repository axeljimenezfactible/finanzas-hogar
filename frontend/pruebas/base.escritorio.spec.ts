// Piezas de la base montadas en el banco de pruebas (pruebas/banco/, Vite de desarrollo, sin dist/). Corre en
// «movil» (412×915 táctil) y en «escritorio» (1440×900). Cubre los hallazgos de frontend-base de WF-3a bis:
// N-04 (descripción accesible de la gráfica), H-03 (calendario, «Ver como tabla», descripción emergente en oscuro),
// H-07/H-08 (selección de filas y Ctrl+/ con vista previa sumada), H-10 (foco sin transición, series sin recolorear,
// deshabilitado de §7.0) y N-06 (barrera de errores).
import { expect, test, type Page } from '@playwright/test'
import { iniciarBanco, type Banco } from './apoyo/banco'
import { controlesChicos, textosEnIngles } from './apoyo/medidas'
let banco: Banco
test.beforeAll(async () => { test.setTimeout(120_000); banco = await iniciarBanco() })
test.afterAll(async () => { await banco?.cerrar() })
/** Color de un token --ui-* resuelto a rgb() en el tema activo. */
const token = (page: Page, v: string) => page.evaluate((nombre) => {
  const d = document.createElement('div'); d.style.color = getComputedStyle(document.documentElement).getPropertyValue(nombre).trim()
  document.body.append(d); const c = getComputedStyle(d).color; d.remove(); return c
}, v)
const guardados = (page: Page) => page.evaluate(() => window.__guardados)

test('H-07/H-08 · Ctrl+/: selección de filas, vista previa sumada, Esc cancela sin guardar y Enter aplica a todas', async ({ page, isMobile }) => {
  test.skip(isMobile, 'Ctrl+/ es atajo de escritorio (U4 §0.7: en Android no hay atajos); en táctil lo cubre la prueba siguiente')
  await page.goto(banco.url('vistabd'))
  await page.getByRole('checkbox', { name: 'Seleccionar Internet' }).check()
  await page.getByRole('checkbox', { name: 'Seleccionar Gimnasio' }).check()
  await expect(page.getByText('2 seleccionadas')).toBeVisible()
  await page.keyboard.press('Control+/')
  const bloque = page.getByRole('region', { name: 'Editar 2 registros a la vez' })
  await expect(bloque).toBeVisible()
  await expect(bloque.getByRole('combobox', { name: 'Propiedad' })).toBeFocused()
  await page.keyboard.press('Enter') // sin valor ni vista previa: no aplica
  expect(await guardados(page)).toEqual([])
  const valor = bloque.getByRole('combobox', { name: 'Valor nuevo' })
  await valor.selectOption('no')
  const previa = bloque.getByTestId('vista-previa')
  await expect(previa).toContainText('Gasto fijo mensual')
  await expect(previa).toContainText('$1,049.00') // Internet 599.00 + Gimnasio 450.00 activos
  await expect(previa).toContainText('$0.00') //   los dos pasan a inactivos
  console.log(`H-07 vista previa sumada: «${(await previa.innerText()).replace(/\s+/g, ' ')}»`)
  await valor.press('Escape')
  await expect(bloque).toHaveCount(0)
  expect(await guardados(page)).toEqual([]) // Esc cancela sin cambiar ninguna
  await expect(page.getByRole('checkbox', { name: 'Seleccionar Internet' })).toBeChecked() // la selección sigue
  await expect(page.getByRole('checkbox', { name: 'Seleccionar Gimnasio' })).toBeFocused() // el foco regresa a donde estaba (la última casilla)
  await page.keyboard.press('Control+/')
  await expect(bloque).toBeVisible()
  await valor.selectOption('no')
  await expect(previa).toContainText('$0.00')
  await valor.press('Enter')
  await expect(bloque).toHaveCount(0)
  expect(await guardados(page)).toEqual([{ ids: [1, 2], propiedad: 'activo', valor: 'no' }])
  await expect(page.getByText('Guardado en 2 registros')).toBeVisible()
  await expect(page.getByRole('checkbox', { name: 'Seleccionar Internet' })).not.toBeChecked() // al aplicar se limpia
  // «Seleccionar todas» y la segunda propiedad (categoría); Ctrl+/ sin selección solo avisa.
  await page.getByRole('checkbox', { name: 'Seleccionar todas' }).check()
  await expect(page.getByText('3 seleccionadas')).toBeVisible()
  await page.keyboard.press('Control+/')
  const bloque3 = page.getByRole('region', { name: 'Editar 3 registros a la vez' })
  await bloque3.getByRole('combobox', { name: 'Propiedad' }).selectOption('categoria')
  await bloque3.getByRole('combobox', { name: 'Valor nuevo' }).selectOption('ocio')
  await expect(page.getByTestId('vista-previa')).toContainText('Registros que cambian de categoría')
  await page.getByRole('button', { name: 'Aplicar a 3' }).click()
  expect((await guardados(page)).at(-1)).toEqual({ ids: [1, 2, 3], propiedad: 'categoria', valor: 'ocio' })
  await page.getByRole('checkbox', { name: 'Seleccionar Internet' }).focus()
  await page.keyboard.press('Control+/')
  await expect(page.getByRole('region', { name: /^Editar/ })).toHaveCount(0)
  await expect(page.getByText('Selecciona filas para editarlas a la vez')).toBeVisible({ timeout: 12_000 }) // tras el aviso anterior
})

test('H-07/H-08 · en táctil «Editar seleccionadas» abre el mismo bloque; casillas con área de toque ≥ 44', async ({ page, isMobile }) => {
  await page.goto(banco.url('vistabd'))
  const casilla = page.getByRole('checkbox', { name: 'Seleccionar Streaming' })
  if (isMobile) await casilla.tap(); else await casilla.check()
  await expect(casilla).toBeChecked()
  await page.getByRole('button', { name: /Editar seleccionadas/ }).click()
  const bloque = page.getByRole('region', { name: 'Editar 1 registro' })
  await bloque.getByRole('combobox', { name: 'Propiedad' }).selectOption('activo')
  const aplicar = bloque.getByRole('button', { name: /Aplicar a 1/ })
  await expect(aplicar).toHaveAttribute('aria-disabled', 'true') // sin vista previa no aplica
  await bloque.getByRole('combobox', { name: 'Valor nuevo' }).selectOption('si')
  await expect(bloque.getByTestId('vista-previa')).toContainText('$1,348.00') // 599 + 450 + 299 activos
  await expect(aplicar).not.toHaveAttribute('aria-disabled', 'true')
  if (isMobile) {
    const c = await controlesChicos(page)
    console.log(`H-07 táctil: ${c.total} controles, menores de 44 px: ${c.chicos.map((x) => `${x.control} ${x.ancho}×${x.alto}`).join('; ') || 'ninguno'}`)
    expect(c.chicos).toEqual([])
  }
  await aplicar.click()
  expect(await guardados(page)).toEqual([{ ids: [3], propiedad: 'activo', valor: 'si' }])
})

for (const tema of ['light', 'dark'] as const) {
  test(`N-04 y H-03 · gráfica en ${tema === 'light' ? 'claro' : 'oscuro'}: aria-label en español con dinero formateado, descripción emergente con tokens y series sin recolorear`, async ({ page, isMobile }) => {
    await page.emulateMedia({ colorScheme: tema })
    await page.goto(banco.url('grafica'))
    const g = page.getByRole('img', { name: /^Gasto por categoría\./ })
    await expect(g.locator('svg')).toBeVisible()
    await expect(g).toHaveAttribute('aria-label', 'Gasto por categoría. Despensa sumó $4,000.00 en oct-2026')
    const ingles = await textosEnIngles(page)
    console.log(`N-04 (${tema}): aria-label «${await g.getAttribute('aria-label')}»; inglés: ${ingles.length}`)
    expect(ingles).toEqual([])
    const rellenos = () => g.locator('svg path').evaluateAll((ps) => ps.map((p) => p.getAttribute('fill')).filter((f) => f && f !== 'none').join('|')) // la guía vertical (sin relleno) aparece al pasar el puntero
    const antes = await rellenos()
    const caja = (await g.boundingBox())!
    await page.mouse.move(caja.x + caja.width * 0.72, caja.y + caja.height * 0.55)
    const tip = page.locator('.fh-descripcion-grafica')
    await expect(tip).toBeVisible()
    await expect(tip).toContainText('$4,000.00')
    const m = await tip.evaluate((el) => {
      const cs = getComputedStyle(el)
      const textos = Array.from(el.querySelectorAll('*')).filter((x) => Array.from(x.childNodes).some((n) => n.nodeType === 3 && n.textContent?.trim())).map((x) => getComputedStyle(x).color)
      return { fondo: cs.backgroundColor, borde: cs.borderTopColor, texto: cs.color, textos: [...new Set(textos)], radio: cs.borderTopLeftRadius }
    })
    const esperado = { fondo: await token(page, '--ui-surface-raised'), borde: await token(page, '--ui-border-divider'), texto: await token(page, '--ui-text-default') }
    console.log(`H-03 descripción emergente (${tema}): fondo ${m.fondo}, borde ${m.borde}, texto ${m.texto} ${m.textos.join(' ')}, radio ${m.radio}; tokens ${JSON.stringify(esperado)}`)
    expect(m.fondo).toBe(esperado.fondo)
    expect(m.borde).toBe(esperado.borde)
    expect(m.texto).toBe(esperado.texto)
    expect(m.textos).toEqual([esperado.texto])
    expect(m.radio).toBe('10px')
    expect(await rellenos()).toBe(antes) // H-10 FV-04: al pasar el puntero las series no cambian de color
    const resumen = page.getByText('Ver como tabla')
    const alto = (await resumen.boundingBox())!.height
    console.log(`H-03 «Ver como tabla» (${isMobile ? '412 táctil' : '1440 fino'}): ${Math.round(alto)} px`)
    expect(alto).toBeGreaterThanOrEqual(isMobile ? 44 : 36)
  })
}

test('H-03 · calendario: en 412 táctil agenda por día con renglones ≥ 44; desde 768 la cuadrícula del mes', async ({ page, isMobile }) => {
  for (const tema of ['light', 'dark'] as const) {
    await page.emulateMedia({ colorScheme: tema })
    await page.goto(banco.url('calendario'))
    const evento = page.getByRole('button', { name: 'Walmart · Sin intereses' })
    await expect(evento).toBeVisible()
    const c = (await evento.boundingBox())!
    console.log(`H-03 evento «Walmart · Sin intereses» (${isMobile ? '412 táctil' : '1440'}, ${tema}): ${Math.round(c.width)}×${Math.round(c.height)}`)
    if (isMobile) {
      await expect(page.getByTestId('agenda')).toBeVisible()
      await expect(page.getByRole('grid')).toHaveCount(0)
      expect(c.width).toBeGreaterThanOrEqual(44)
      expect(c.height).toBeGreaterThanOrEqual(44)
      expect((await controlesChicos(page)).chicos).toEqual([])
    } else {
      await expect(page.getByRole('grid')).toBeVisible()
    }
  }
})

test('H-10 · foco sin transición y deshabilitado de §7.0 en Campo y ChipSeleccionable', async ({ page }) => {
  await page.goto(banco.url('controles'))
  const transicion = await page.getByRole('button', { name: 'Guardar' }).evaluate((b) => getComputedStyle(b).transitionProperty)
  console.log(`H-10 Botón transition-property: ${transicion}`)
  expect(transicion).not.toMatch(/outline|all/)
  const sup = await token(page, '--ui-surface-disabled'), txt = await token(page, '--ui-text-disabled')
  for (const el of [page.getByLabel('Monto deshabilitado'), page.getByRole('button', { name: 'Ocio' })]) {
    const m = await el.evaluate((x) => { const cs = getComputedStyle(x); return { borde: cs.borderTopStyle, fondo: cs.backgroundColor, texto: cs.color, cursor: cs.cursor } })
    console.log(`H-10 deshabilitado: ${JSON.stringify(m)} (tokens ${sup} / ${txt})`)
    expect(m).toEqual({ borde: 'dashed', fondo: sup, texto: txt, cursor: 'not-allowed' })
  }
})

test('N-06 · una pantalla que falla al pintarse muestra «No se pudo mostrar esta página · Reintentar»', async ({ page }) => {
  await page.goto(banco.url('barrera'))
  await expect(page.getByRole('alert')).toContainText('No se pudo mostrar esta página')
  await expect(page.getByRole('alert').getByRole('button', { name: 'Reintentar' })).toBeVisible()
})
