// WF-3a · front-variantes: renderiza CADA historia de Storybook en claro y oscuro, en escritorio (puntero fino, con hover)
// y en móvil (360 px, táctil), provoca hover, presionado y foco con el puntero y el teclado reales, mide texto, iconos,
// bordes y marcas contra su fondo efectivo y calcula el contraste con la fórmula de docs/design/tools/contraste_componentes.py
// (scripts/contraste-variantes.py la importa: no se reescribe). Además comprueba: 0 errores de consola, 0 hosts externos,
// 0 desbordes horizontales a 360 px y áreas táctiles de 44 px con puntero grueso.
//
// Uso:  npm run build-storybook && npm run verificar:variantes
//   FH_CHROMIUM=<ruta a chrome-headless-shell>   (si el Chromium de Playwright no está descargado)
//   --salida <dir>                                carpeta del CSV/JSON (por omisión ../datos/tmp/variantes)
//   --solo <texto>                                solo las historias cuyo id contiene el texto
// Código de salida: 0 si no hay fallas fuera de las historias marcadas «hallazgo»; 1 si las hay; 2 si el entorno falla.
import { createServer } from 'node:http'
import { spawnSync } from 'node:child_process'
import { existsSync, mkdirSync, readdirSync, readFileSync, writeFileSync, statSync } from 'node:fs'
import path from 'node:path'
import { chromium } from '@playwright/test'
import { desdeFront, desdeRepo } from './comun.mjs'

const args = process.argv.slice(2)
const opcion = (n, d) => { const i = args.indexOf(n); return i >= 0 ? args[i + 1] : d }
const SALIDA = path.resolve(opcion('--salida', desdeRepo('datos/tmp/variantes')))
const SOLO = opcion('--solo', '')
const ESTATICO = desdeFront('storybook-static')
mkdirSync(SALIDA, { recursive: true })
if (!existsSync(path.join(ESTATICO, 'index.json'))) { console.error('Falta storybook-static: corre «npm run build-storybook».'); process.exit(2) }

// ---------------------------------------------------------------- servidor estático local (sin red externa)
const TIPOS = { '.html': 'text/html; charset=utf-8', '.js': 'text/javascript', '.mjs': 'text/javascript', '.css': 'text/css', '.json': 'application/json', '.woff2': 'font/woff2', '.svg': 'image/svg+xml', '.png': 'image/png', '.ico': 'image/x-icon', '.txt': 'text/plain', '.map': 'application/json' }
const servidor = createServer((req, res) => {
  const ruta = decodeURIComponent(new URL(req.url, 'http://x').pathname)
  let archivo = path.join(ESTATICO, ruta === '/' ? 'index.html' : ruta)
  if (!archivo.startsWith(ESTATICO) || !existsSync(archivo) || statSync(archivo).isDirectory()) { res.writeHead(404); res.end('no existe'); return }
  res.writeHead(200, { 'Content-Type': TIPOS[path.extname(archivo)] ?? 'application/octet-stream' })
  res.end(readFileSync(archivo))
})
await new Promise((r) => servidor.listen(0, '127.0.0.1', r))
const BASE = `http://127.0.0.1:${servidor.address().port}`

// ---------------------------------------------------------------- tokens --ui-* (para nombrar los colores medidos)
const tokensCss = readFileSync(desdeFront('src/estilos/tokens.css'), 'utf8')
const VARS = [...new Set([...tokensCss.matchAll(/^\s*(--ui-[a-z0-9-]+)\s*:/gm)].map((m) => m[1]))]

// ---------------------------------------------------------------- medición dentro de la página (autocontenida)
function medirEnPagina({ raizSel, vars, estado }) {
  const cv = document.createElement('canvas'); cv.width = cv.height = 1
  const cx = cv.getContext('2d', { willReadFrequently: true })
  const aRGBA = (c) => {
    if (!c || c === 'transparent' || c === 'none') return [0, 0, 0, 0]
    const m = c.match(/^rgba?\(([^)]+)\)$/)
    if (m && !m[1].includes('%') && !m[1].includes('none')) {
      const p = m[1].split(/[ ,/]+/).filter(Boolean).map(Number)
      return [p[0], p[1], p[2], p.length > 3 ? p[3] : 1]
    }
    cx.clearRect(0, 0, 1, 1); cx.fillStyle = '#010203'; cx.fillStyle = c
    if (cx.fillStyle === '#010203' && !/#?010203/.test(c)) return null
    cx.fillRect(0, 0, 1, 1)
    const d = cx.getImageData(0, 0, 1, 1).data
    return [d[0], d[1], d[2], d[3] / 255]
  }
  const sobre = (f, b) => [0, 1, 2].map((i) => Math.round(f[i] * f[3] + b[i] * (1 - f[3])))
  const hex = (rgb) => '#' + rgb.map((v) => v.toString(16).padStart(2, '0')).join('').toUpperCase()
  // Mapa hex → nombres de token (tema vigente); solo colores opacos.
  const nombres = {}
  const cs0 = getComputedStyle(document.documentElement)
  for (const v of vars) {
    const c = aRGBA(cs0.getPropertyValue(v).trim())
    if (!c || c[3] < 1) continue
    ;(nombres[hex(c.slice(0, 3))] ||= []).push(v)
  }
  const decorativos = new Set(['--ui-border-divider', '--ui-viz-grid', '--ui-viz-track', '--ui-skeleton-on-page', '--ui-skeleton-on-card'].map((v) => { const c = aRGBA(cs0.getPropertyValue(v).trim()); return c ? hex(c.slice(0, 3)) : '' }))
  const fondoEfectivo = (el) => {
    const cadena = []
    for (let e = el; e; e = e.parentElement) cadena.push(e)
    let base = [255, 255, 255]
    for (const e of cadena.reverse()) {
      const c = aRGBA(getComputedStyle(e).backgroundColor)
      if (c && c[3] > 0) base = sobre(c, base)
    }
    return base
  }
  const visible = (el) => {
    const cs = getComputedStyle(el)
    if (cs.display === 'none' || cs.visibility === 'hidden' || Number(cs.opacity) === 0) return false
    const r = el.getBoundingClientRect()
    return r.width > 1 && r.height > 1
  }
  const excluido = (el) => !!el.closest('.sb-wrapper, #storybook-docs, [data-aria-hidden="true"], script, style, noscript')
  const dentro = (el, raiz) => !raiz || raiz === el || raiz.contains(el)
  const raiz = raizSel ? document.querySelector(raizSel) : null
  const sal = []
  const nom = (h) => (nombres[h] ?? []).map((v) => v.slice(2).replace(/-/g, '.'))
  const descr = (el) => {
    const t = (el.getAttribute('aria-label') ?? el.textContent ?? '').trim().replace(/\s+/g, ' ').slice(0, 40)
    const c = typeof el.className === 'string' ? el.className.split(/\s+/).filter(Boolean)[0] : ''
    return `${el.localName}${t ? ` «${t}»` : ''}${c ? ` .${c}` : ''}`
  }
  const anota = (tipo, el, fg, bg, umbral, extra = {}) => {
    const f = hex(fg), b = hex(bg)
    sal.push({ tipo, estado, elemento: descr(el), fg: f, bg: b, tokFg: nom(f), tokBg: nom(b), umbral, ...extra })
  }
  const todos = [...document.body.querySelectorAll('*')].filter((el) => !excluido(el) && dentro(el, raiz) && visible(el))
  for (const el of todos) {
    const cs = getComputedStyle(el)
    const ns = el.namespaceURI
    const esSvg = ns === 'http://www.w3.org/2000/svg'
    // Texto (nodos de texto propios).
    const tieneTexto = [...el.childNodes].some((n) => n.nodeType === 3 && n.textContent.trim())
    if (tieneTexto && (!esSvg || el.localName === 'text' || el.localName === 'tspan')) {
      const fgRaw = esSvg ? cs.fill : cs.color
      const c = aRGBA(fgRaw)
      if (c && c[3] > 0) {
        const bg = fondoEfectivo(el)
        const px = parseFloat(cs.fontSize), peso = Number(cs.fontWeight)
        const grande = px >= 24 || (px >= 18.66 && peso >= 700)
        anota('texto', el, sobre(c, bg), bg, grande ? 3 : 4.5, { px, peso })
      }
    }
    // Valor y marcador de campos.
    if (el.localName === 'input' || el.localName === 'textarea') {
      const bg = fondoEfectivo(el)
      if (el.value) { const c = aRGBA(cs.color); if (c) anota('texto', el, sobre(c, bg), bg, 4.5, { valor: true }) }
      else if (el.placeholder) { const c = aRGBA(getComputedStyle(el, '::placeholder').color); if (c) anota('texto', el, sobre(c, bg), bg, 4.5, { marcador: true }) }
    }
    // Iconos (lucide) y marcas de forma (chips): no texto, 3:1.
    if (esSvg && el.localName === 'svg' && (el.classList.contains('lucide') || el.getAttribute('class')?.includes('lucide'))) {
      const c = aRGBA(cs.stroke === 'none' ? cs.color : cs.stroke)
      if (c && c[3] > 0) { const bg = fondoEfectivo(el.parentElement ?? el); anota('icono', el, sobre(c, bg), bg, 3) }
    }
    if (!esSvg && el.getAttribute('aria-hidden') === 'true' && el.children.length === 0 && !el.textContent.trim()) {
      // Marca de forma: relleno o borde grueso de color de persona/tienda.
      const rel = aRGBA(cs.backgroundColor)
      const bw = parseFloat(cs.borderTopWidth)
      const col = rel && rel[3] > 0 ? rel : bw >= 2 ? aRGBA(cs.borderTopColor) : null
      if (col && col[3] > 0 && !decorativos.has(hex(col.slice(0, 3)))) { const bg = fondoEfectivo(el.parentElement); anota('marca', el, sobre(col, bg), bg, 3) }
    }
    // Marcas de series de gráficas (ECharts, SVG): solo las que usan un color --ui-chart-N.
    if (esSvg && ['path', 'rect', 'circle', 'polygon', 'polyline', 'ellipse'].includes(el.localName) && !el.closest('svg.lucide')) {
      for (const prop of ['fill', 'stroke']) {
        const c = aRGBA(cs[prop])
        if (!c || c[3] === 0) continue
        const h = hex(c.slice(0, 3))
        if ((nombres[h] ?? []).some((v) => /^--ui-chart-\d$/.test(v) || v === '--ui-viz-other')) {
          const bg = fondoEfectivo(el.closest('svg')?.parentElement ?? el)
          anota('marca-serie', el, sobre(c, bg), bg, 3, { propiedad: prop })
        } else if (prop === 'stroke' && parseFloat(cs.strokeWidth) >= 1 && !decorativos.has(h) && !(nombres[h] ?? []).length) {
          const bg = fondoEfectivo(el.closest('svg')?.parentElement ?? el)
          anota('fuera-de-tokens', el, sobre(c, bg), bg, 0, { propiedad: prop })
        }
      }
    }
    // Bordes de controles: contra el fondo de afuera y contra el interior (los divisores son decorativos).
    if (!esSvg) {
      const lados = ['Top', 'Right', 'Bottom', 'Left'].filter((l) => parseFloat(cs[`border${l}Width`]) > 0 && cs[`border${l}Style`] !== 'none')
      const vistos = new Set()
      for (const l of lados) {
        const c = aRGBA(cs[`border${l}Color`])
        if (!c || c[3] === 0) continue
        const h = hex(c.slice(0, 3))
        if (vistos.has(h)) continue
        vistos.add(h)
        if (decorativos.has(h)) continue
        const fuera = fondoEfectivo(el.parentElement)
        anota('borde-fuera', el, sobre(c, fuera), fuera, 3, { lado: l })
        const propio = aRGBA(cs.backgroundColor)
        if (propio && propio[3] === 1) anota('borde-dentro', el, sobre(c, propio.slice(0, 3)), propio.slice(0, 3), 3, { lado: l })
      }
    }
  }
  // Controles deshabilitados: DESIGN.md §7.0 pide relleno ui.surface.disabled, texto ui.text.disabled y borde punteado en ui.border.control.
  const tok = (v) => { const c = aRGBA(cs0.getPropertyValue(v).trim()); return c ? hex(c.slice(0, 3)) : '' }
  for (const el of todos.filter((e) => e.matches('button, input, select, textarea, [aria-disabled="true"]') && (e.disabled || e.getAttribute('aria-disabled') === 'true'))) {
    const cs = getComputedStyle(el)
    const dif = []
    const bg = aRGBA(cs.backgroundColor), col = aRGBA(cs.color), bc = aRGBA(cs.borderTopColor)
    if (!bg || bg[3] < 1 || hex(bg.slice(0, 3)) !== tok('--ui-surface-disabled')) dif.push(`relleno ${bg ? hex(bg.slice(0, 3)) : '?'} en lugar de ${tok('--ui-surface-disabled')}`)
    if (!col || hex(col.slice(0, 3)) !== tok('--ui-text-disabled')) dif.push(`texto ${col ? hex(col.slice(0, 3)) : '?'} en lugar de ${tok('--ui-text-disabled')}`)
    if (cs.borderTopStyle !== 'dashed' || parseFloat(cs.borderTopWidth) < 1 || !bc || hex(bc.slice(0, 3)) !== tok('--ui-border-control')) dif.push(`borde ${cs.borderTopStyle} ${cs.borderTopWidth} en lugar de punteado de 1px en ui.border.control`)
    sal.push({ tipo: 'deshabilitado-estilo', estado, elemento: descr(el), fg: '', bg: '', tokFg: [], tokBg: [], umbral: 0, diferencias: dif })
  }
  // Anillo de foco del elemento activo (fuera del elemento, sobre el fondo de su padre).
  const act = document.activeElement
  if (estado === 'foco' && act && act !== document.body && dentro(act, raiz)) {
    const cs = getComputedStyle(act)
    const c = aRGBA(cs.outlineColor)
    const grosor = parseFloat(cs.outlineWidth)
    const ok = cs.outlineStyle !== 'none' && grosor >= 2
    if (c && ok) { const bg = fondoEfectivo(act.parentElement ?? act); anota('foco', act, sobre(c, bg), bg, 3, { grosor, separacion: parseFloat(cs.outlineOffset) }) }
    else sal.push({ tipo: 'foco-ausente', estado, elemento: descr(act), fg: '', bg: '', tokFg: [], tokBg: [], umbral: 0, grosor, estilo: cs.outlineStyle })
  }
  return sal
}

// Elementos interactivos: se marcan con data-fv para recorrerlos.
function marcarInteractivos() {
  const sel = 'button, a[href], input:not([type=hidden]), select, textarea, summary, [role=button], [role=link], [tabindex="0"]'
  const lista = [...document.body.querySelectorAll(sel)].filter((el) => {
    if (el.closest('.sb-wrapper, #storybook-docs, [data-aria-hidden="true"]')) return false
    const cs = getComputedStyle(el), r = el.getBoundingClientRect()
    return cs.display !== 'none' && cs.visibility !== 'hidden' && r.width > 1 && r.height > 1
  })
  return lista.map((el, i) => {
    el.setAttribute('data-fv', String(i))
    const r = el.getBoundingClientRect()
    const control = el.matches('button, a[href], input, select, textarea, summary, [role=button]')
    const t = (el.getAttribute('aria-label') ?? el.textContent ?? '').trim().replace(/\s+/g, ' ').slice(0, 40)
    // DESIGN.md §7.0: el anillo de foco aparece sin transición.
    const cs = getComputedStyle(el)
    const props = cs.transitionProperty.split(',').map((x) => x.trim()), durs = cs.transitionDuration.split(',').map((x) => parseFloat(x))
    const transicionFoco = props.some((p, k) => ['all', 'outline', 'outline-color', 'outline-width', 'outline-offset'].includes(p) && (durs[k % durs.length] || 0) > 0)
    return { i, control, transicionFoco, ancho: Math.round(r.width * 10) / 10, alto: Math.round(r.height * 10) / 10, desc: `${el.localName}${t ? ` «${t}»` : ''}` }
  })
}

// ---------------------------------------------------------------- recorrido
const indice = JSON.parse(readFileSync(path.join(ESTATICO, 'index.json'), 'utf8'))
const historias = Object.values(indice.entries).filter((e) => e.type === 'story' && e.id.includes(SOLO))
console.log(`Historias: ${historias.length}`)
if (!historias.length) { console.error('Ninguna historia coincide: nada que medir.'); process.exit(2) }

// Chromium: FH_CHROMIUM, el de Playwright si está descargado o el chrome-headless-shell más nuevo de ~/.cache/ms-playwright.
function buscarChromium() {
  if (process.env.FH_CHROMIUM) return process.env.FH_CHROMIUM
  if (existsSync(chromium.executablePath())) return undefined
  const cache = path.join(process.env.HOME ?? '', '.cache/ms-playwright')
  const dirs = existsSync(cache) ? readdirSync(cache).filter((d) => d.startsWith('chromium_headless_shell-')).sort().reverse() : []
  for (const d of dirs) { const ex = path.join(cache, d, 'chrome-headless-shell-linux64/chrome-headless-shell'); if (existsSync(ex)) return ex }
  return undefined
}
const navegador = await chromium.launch({ executablePath: buscarChromium() })
const CONTEXTOS = {
  escritorio: { viewport: { width: 1280, height: 800 }, hasTouch: false, isMobile: false },
  movil: { viewport: { width: 360, height: 800 }, hasTouch: true, isMobile: true, deviceScaleFactor: 2 },
}
const medidas = [], tamanos = [], problemas = []
let paginas = 0
for (const [nombreCtx, cfg] of Object.entries(CONTEXTOS)) {
  const ctx = await navegador.newContext({ ...cfg, locale: 'es-MX', timezoneId: 'America/Mexico_City', reducedMotion: 'no-preference' })
  const page = await ctx.newPage()
  const hosts = new Set()
  let consola = []
  page.on('request', (r) => { try { const u = new URL(r.url()); if (!['127.0.0.1', 'localhost'].includes(u.hostname) && !['data:', 'blob:'].includes(u.protocol)) hosts.add(u.host) } catch { /* sin URL */ } })
  page.on('console', (m) => { if (m.type() === 'error') consola.push(m.text().slice(0, 200)) })
  page.on('pageerror', (e) => consola.push(`pageerror: ${String(e.message).slice(0, 200)}`))
  for (const h of historias) {
    for (const tema of ['claro', 'oscuro']) {
      consola = []
      const base = { historia: h.id, titulo: `${h.title} › ${h.name}`, tema, contexto: nombreCtx, hallazgo: (h.tags ?? []).includes('hallazgo') }
      await page.goto(`${BASE}/iframe.html?id=${h.id}&viewMode=story&globals=tema:${tema}`, { waitUntil: 'load' })
      try {
        await page.waitForFunction(() => document.body.classList.contains('sb-show-main') || document.body.classList.contains('sb-show-errordisplay'), null, { timeout: 15000 })
      } catch { problemas.push({ ...base, problema: 'La historia no terminó de renderizar en 15 s' }); continue }
      if (await page.evaluate(() => document.body.classList.contains('sb-show-errordisplay'))) { problemas.push({ ...base, problema: 'Storybook muestra error en la historia' }); continue }
      await page.evaluate(() => document.fonts.ready)
      await page.waitForTimeout(250)
      paginas++
      const temaAplicado = await page.evaluate(() => document.documentElement.getAttribute('data-theme'))
      if (temaAplicado !== (tema === 'claro' ? 'light' : 'dark')) problemas.push({ ...base, problema: `data-theme=${temaAplicado}, se esperaba ${tema}` })
      const reg = (lista) => lista.forEach((m) => medidas.push({ ...base, ...m }))
      // 1) reposo: toda la página
      reg(await page.evaluate(medirEnPagina, { raizSel: null, vars: VARS, estado: 'reposo' }))
      const interactivos = await page.evaluate(marcarInteractivos)
      for (const it of interactivos.filter((x) => x.control)) tamanos.push({ ...base, elemento: it.desc, ancho: it.ancho, alto: it.alto })
      for (const it of interactivos.filter((x) => x.transicionFoco)) medidas.push({ ...base, tipo: 'foco-transicion', estado: 'foco', elemento: it.desc, fg: '', bg: '', tokFg: [], tokBg: [], umbral: 0 })
      // A partir de aquí los estados se miden ya asentados (sin transiciones ni animaciones en curso).
      await page.addStyleTag({ content: '*, *::before, *::after { transition: none !important; animation: none !important; }' })
      // 2) desborde horizontal de la página
      const desborde = await page.evaluate(() => document.documentElement.scrollWidth - document.documentElement.clientWidth)
      if (desborde > 0) problemas.push({ ...base, problema: `Desborde horizontal de ${desborde} px` })
      // 3) hover y presionado (solo escritorio: en Android no existe hover)
      if (nombreCtx === 'escritorio') {
        for (const it of interactivos) {
          const sel = `[data-fv="${it.i}"]`
          try {
            await page.hover(sel, { timeout: 1500 })
            reg(await page.evaluate(medirEnPagina, { raizSel: sel, vars: VARS, estado: 'hover' }))
            await page.mouse.down()
            reg(await page.evaluate(medirEnPagina, { raizSel: sel, vars: VARS, estado: 'presionado' }))
            await page.mouse.move(2, 2)
            await page.mouse.up()
          } catch { /* oculto por una capa (hoja abierta): se mide con el teclado */ await page.mouse.up().catch(() => {}) }
        }
        await page.mouse.move(2, 2)
      }
      // 4) foco con el teclado (Tab): el anillo y el contenido del elemento enfocado
      await page.evaluate(() => { document.activeElement?.blur?.(); window.scrollTo(0, 0) })
      await page.keyboard.press('Shift')
      const vistos = new Set()
      for (let n = 0; n < interactivos.length * 2 + 4; n++) {
        await page.keyboard.press('Tab')
        const marca = await page.evaluate(() => document.activeElement?.getAttribute?.('data-fv') ?? null)
        if (marca === null || vistos.has(marca)) { if (marca !== null && vistos.size >= interactivos.length) break; continue }
        vistos.add(marca)
        reg(await page.evaluate(medirEnPagina, { raizSel: `[data-fv="${marca}"]`, vars: VARS, estado: 'foco' }))
      }
      if (consola.length) problemas.push({ ...base, problema: `Errores de consola: ${[...new Set(consola)].slice(0, 3).join(' | ')}` })
    }
  }
  if (hosts.size) problemas.push({ contexto: nombreCtx, problema: `Hosts externos: ${[...hosts].join(', ')}` })
  console.log(`${nombreCtx}: hosts externos ${hosts.size}`)
  await ctx.close()
}
await navegador.close()
servidor.close()

const entrada = path.join(SALIDA, 'medidas.json')
writeFileSync(entrada, JSON.stringify({ historias: historias.length, paginas, medidas, tamanos, problemas }))
// Cálculo del contraste con la fórmula de docs/design/tools/contraste_componentes.py
const r = spawnSync('python3', ['-I', desdeFront('scripts/contraste-variantes.py'), entrada, SALIDA, ...(SOLO ? ['--parcial'] : [])], { stdio: 'inherit' })
process.exit(r.status ?? 2)
