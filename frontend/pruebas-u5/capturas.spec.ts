// U5-01..U5-04, U5-08, U5-10 · Una prueba por pantalla (P-xx y P-26.xx de U1) en cada combinación (proyecto):
// captura de página completa + axe-core (5 etiquetas WCAG) + contraste medido en el DOM + hosts pedidos + comprobación del tema.
// El título lleva «axe» y «contraste» para que --grep axe / --grep contraste las seleccionen (axe y contraste se corren sobre CADA captura).
import { test as base, expect } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'
import fs from 'node:fs'
import path from 'node:path'
import { AQUI, EVID, ETIQUETAS_AXE, PARTES, REPO, SALIDA, abrirAnonima, abrirSesion, fondoC2, listo, nombresRutas, pantallasU1, resolverIds, ruta, type Sesion } from './apoyo'

const CONTRASTE_JS = fs.readFileSync(path.join(AQUI, 'contraste-dom.js'), 'utf8')
const test = base.extend<object, { sesion: Sesion }>({
  sesion: [async ({ browser }, use, wi) => { const s = await abrirSesion(browser, wi.project.use as never); await use(s); await s.context.close() }, { scope: 'worker' }],
})
for (const d of [SALIDA, path.join(SALIDA, 'axe'), path.join(SALIDA, 'contraste'), EVID, PARTES]) fs.mkdirSync(d, { recursive: true })

// Cifras de la semilla que deben verse en pantalla (U5-02): salen de oraculo/model.json (se comprueba que existan ahí).
const modelo = fs.readFileSync(path.join(REPO, 'oraculo', 'model.json'), 'utf8')
const CIFRAS: { etiqueta: string; texto: string; enModelo: string }[] = [
  { etiqueta: 'disponible inicial A', texto: '$10,000.00', enModelo: '10000.0' },
  { etiqueta: 'disponible inicial B', texto: '$6,905.50', enModelo: '6905.5' },
  { etiqueta: 'pago del 30-oct-2026', texto: '$973.75', enModelo: '973.75' },
  { etiqueta: 'deuda viva al 10-abr-2027', texto: '$6,293.44', enModelo: '6293.44' },
  { etiqueta: 'compra del 10-oct-2026', texto: '$3,094.50', enModelo: '3094.5' },
  { etiqueta: 'parte de Karen (capital)', texto: '$1,191.50', enModelo: '1191.5' },
].map((c) => ({ ...c, enModelo: String(modelo.includes(c.enModelo)) }))

// U5-07 sobre el texto visible (innerText) de cada captura: (a) formato de cada importe, (b) 0 fechas ISO, (c) 0 palabras de
// interfaz en inglés de la lista del criterio (más las de la descripción automática de ECharts y «Fixture») y (d) 0 claves
// internas en snake_case (N-1 de REVISION-u5-r2: la nota de la semilla mostraba «Fixture s2-metas.json, savings_goals[2]»).
const INGLES_U507 = /\b(Save|Cancel|Settings|Dashboard|Loading|Error|Home|Delete|Search|Add|Submit|Close|Fixture|This is a chart|series count|The data is as follows)\b/g
const CLAVE_INTERNA = /\b[a-z]+_[a-z_]+\b/g
const FECHA_ISO = /[0-9]{4}-[0-9]{2}-[0-9]{2}/g
const IMPORTE = /-?\$\s?[0-9](?:[0-9,]*[0-9])?(?:\.[0-9]+)?(?:\s?mil)?/g
const IMPORTE_OK = /^-?\$[0-9]{1,3}(?:,[0-9]{3})*\.[0-9]{2}$/
const IMPORTE_EJE = /^\$[0-9]+(?:\.[0-9])?\s?mil$/ // compacto: solo en ejes de gráficas (DESIGN §7.7)
function revisarTexto(t: string) {
  const unicos = (re: RegExp) => [...new Set([...t.matchAll(re)].map((m) => m[0]))]
  const importes = [...t.matchAll(IMPORTE)].map((m) => m[0])
  return { importes: importes.length, importes_fuera_de_formato: [...new Set(importes.filter((x) => !IMPORTE_OK.test(x) && !IMPORTE_EJE.test(x)))], fechas_iso: unicos(FECHA_ISO), ingles: unicos(INGLES_U507), claves_internas: unicos(CLAVE_INTERNA) }
}

const pantallas = pantallasU1()
const nombres = nombresRutas()
test.describe('capturas', () => {
  test('T = 41 pantallas de U1 (14 secciones P-26.xx incluidas)', () => {
    expect(pantallas.filter((p) => !p.id.includes('.')).length).toBe(27)
    expect(pantallas.filter((p) => p.id.startsWith('P-26.')).length).toBe(14)
    expect(pantallas.length).toBe(41)
  })
  for (const { id, nombre: nombreU1 } of pantallas) {
    test(`${id} ${nombreU1} · captura, axe y contraste`, async ({ sesion, browser }, info) => {
      const [disp, tema] = info.project.name.split('-') as ['android' | 'escritorio', 'claro' | 'oscuro']
      const nombre = `${id}-${disp}-${tema}`
      // Página nueva por pantalla (N-2 de REVISION-u5-r2): al reutilizar la página de la sesión tras cerrar la anónima de P-00,
      // Chromium dejaba de emular el puntero táctil y las capturas Android salían con controles de 36 px. La cookie es del contexto.
      const ses = id === 'P-00' ? await abrirAnonima(browser, info.project.use as never) : { context: sesion.context, page: await sesion.context.newPage() }
      const { page, context } = ses
      const hosts = new Set<string>()
      const alPedir = (r: { url(): string }) => { try { const u = new URL(r.url()); if (/^(https?|wss?):$/.test(u.protocol)) hosts.add(u.host) } catch { /* data: / blob: no son red */ } }
      context.on('request', alPedir)
      try {
        const ids = id === 'P-00' ? { gasto: '1', receta: '1', credito: '1', compra: '1' } : await resolverIds(sesion)
        await page.goto(ruta(id, ids))
        let { cargando } = await listo(page)
        let alcance: string | undefined
        if (id === 'P-03') { await page.getByRole('button', { name: 'Registrar gasto' }).first().click(); await page.locator('[role=dialog]').first().waitFor({ state: 'visible' }); alcance = '[role=dialog]'; await page.waitForTimeout(600) }
        if (id === 'P-22') { await page.getByRole('button', { name: /Registrar peso/ }).first().click(); await page.locator('[role=dialog]').first().waitFor({ state: 'visible' }); alcance = '[role=dialog]'; await page.waitForTimeout(600) }
        await page.mouse.move(0, 0)
        // Android: el puntero grueso debe estar activo al capturar (DESIGN §4.6, controles de 44 px con `coarse:`).
        const punteroGrueso = await page.evaluate(() => matchMedia('(any-pointer: coarse)').matches)
        // 1 · Captura de página completa
        const png = path.join(SALIDA, `${nombre}.png`)
        await page.screenshot({ path: png, fullPage: true, animations: 'disabled' })
        // 2 · axe-core
        const axe = await new AxeBuilder({ page: page as never }).withTags(ETIQUETAS_AXE).analyze()
        const serias = axe.violations.filter((v) => v.impact === 'serious' || v.impact === 'critical')
        fs.writeFileSync(path.join(SALIDA, 'axe', `${nombre}.json`), JSON.stringify({
          pantalla: id, nombre: nombreU1, combinacion: info.project.name, url: page.url().replace(/^https?:\/\/[^/]+/, ''), etiquetas: ETIQUETAS_AXE, version_axe: axe.testEngine.version,
          resumen: { violaciones: axe.violations.length, serias_o_criticas: serias.length, incompletos: axe.incomplete.length, pasadas: axe.passes.length, no_aplican: axe.inapplicable.length },
          violations: axe.violations.map((v) => ({ id: v.id, impact: v.impact, help: v.help, helpUrl: v.helpUrl, tags: v.tags, nodes: v.nodes.map((n) => ({ target: n.target, html: n.html.slice(0, 300), failureSummary: n.failureSummary })) })),
          incomplete: axe.incomplete.map((v) => ({ id: v.id, impact: v.impact, help: v.help, nodes: v.nodes.map((n) => ({ target: n.target, html: n.html.slice(0, 300), message: n.any[0]?.message ?? n.all[0]?.message })) })),
        }, null, 1))
        // 3 · Contraste en el DOM
        const cont = await page.evaluate(`(${CONTRASTE_JS})(${JSON.stringify(alcance ?? null)})`) as { evaluados: number; bajo_umbral: number; indeterminados: number }
        fs.writeFileSync(path.join(SALIDA, 'contraste', `${nombre}.json`), JSON.stringify({ pantalla: id, combinacion: info.project.name, ...cont }, null, 1))
        // 4 · Tema real y datos de la página
        const esperado = fondoC2(tema)
        const dom = await page.evaluate(() => ({
          bg: getComputedStyle(document.body).backgroundColor, dataTheme: document.documentElement.getAttribute('data-theme'),
          oscuroSistema: matchMedia('(prefers-color-scheme: dark)').matches, h1: document.querySelector('h1')?.textContent?.trim() ?? '',
          texto: document.body.innerText, errorPantalla: /No se pudo mostrar esta página/.test(document.body.innerText),
          ancho: document.documentElement.scrollWidth, alto: document.documentElement.scrollHeight,
        }))
        const u507 = revisarTexto(dom.texto)
        const cifras = CIFRAS.filter((c) => dom.texto.includes(c.texto)).map((c) => c.texto)
        const manifiesto = {
          pantalla: id, nombre: nombreU1, combinacion: info.project.name, archivo: `${nombre}.png`, ruta: page.url().replace(/^https?:\/\/[^/]+/, ''), h1: dom.h1,
          h1_coincide: id === 'P-03' || id === 'P-22' || id.includes('.') || id === 'P-26' ? null : dom.h1 === (nombres[id] ?? dom.h1),
          cargando_al_capturar: cargando, puntero_grueso: punteroGrueso, error_de_pantalla: dom.errorPantalla, sesion: id === 'P-00' ? 'sin sesión (formulario)' : 'iniciada una vez por contexto',
          tema: { esperado: tema, fondo_body: dom.bg, fondo_token_c2: esperado.rgb, token: esperado.ref, coincide: dom.bg === esperado.rgb, data_theme: dom.dataTheme, prefers_dark: dom.oscuroSistema },
          dimensiones_pagina: { ancho: dom.ancho, alto: dom.alto }, cifras_visibles: cifras, axe: { serias_o_criticas: serias.length, violaciones: axe.violations.length },
          contraste: { evaluados: cont.evaluados, bajo_umbral: cont.bajo_umbral, indeterminados: cont.indeterminados }, u5_07: u507,
        }
        fs.writeFileSync(path.join(PARTES, `${nombre}.json`), JSON.stringify({ nombre, manifiesto, hosts: [...hosts].sort() }))
        // El tema oscuro debe ser real (fondo = token de C2), no inversión; en claro, el token claro. No se oculta nada: fallan las aserciones y queda el JSON.
        expect(dom.bg, 'fondo del body = token de C2 del tema').toBe(esperado.rgb)
        expect(dom.dataTheme, 'sin tocar el interruptor').toBeNull()
        expect(dom.oscuroSistema).toBe(tema === 'oscuro')
        if (disp === 'android') expect(punteroGrueso, 'Android con (any-pointer: coarse)').toBe(true)
        expect(u507.fechas_iso, 'U5-07 fechas ISO visibles').toEqual([])
        expect(u507.ingles, 'U5-07 palabras en inglés visibles').toEqual([])
        expect(u507.claves_internas, 'U5-07 claves internas visibles').toEqual([])
        expect([...hosts].filter((h) => !/^(finanzas\.local|localhost|127\.0\.0\.1)(:\d+)?$/.test(h)), 'U5-08 hosts externos').toEqual([])
      } finally {
        context.off('request', alPedir)
        if (id === 'P-00') await ses.context.close()
        else await page.close()
      }
    })
  }
})
// Las cifras de la semilla deben existir en oraculo/model.json (U5-02); queda dicho en evidencia/cifras.json al cerrar.
test.afterAll(() => { fs.writeFileSync(path.join(EVID, 'cifras-oraculo.json'), JSON.stringify(CIFRAS, null, 1)) })
