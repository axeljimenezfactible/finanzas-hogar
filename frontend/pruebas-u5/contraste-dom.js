// U5-04 · Contraste medido en el DOM (WCAG 2.x). Se evalúa dentro de la página con page.evaluate("(" + este archivo + ")(alcance)").
// Texto ≥ 4.5:1; texto grande (≥ 24 px, o ≥ 18.66 px en negrita) y componentes de UI / íconos ≥ 3:1.
// Color de texto y fondo efectivo salen de los estilos computados; el fondo se compone por ascendencia (capas con alfa hasta una opaca).
// Los canales se redondean a 8 bits (como pinta Chrome) ANTES de calcular la razón, así cualquiera la reproduce desde los hex.
(function (alcanceSel) {
  const raiz = (alcanceSel && document.querySelector(alcanceSel)) || document.body
  const ctx = document.createElement('canvas').getContext('2d', { willReadFrequently: true })
  const memo = new Map()
  function parse(s) {
    if (memo.has(s)) return memo.get(s)
    let m, o
    if ((m = /^rgba?\(([^)]+)\)$/.exec(s))) { const p = m[1].split(/[\s,\/]+/).filter(Boolean).map(parseFloat); o = [p[0], p[1], p[2], p.length > 3 ? p[3] : 1] }
    else if ((m = /^color\(srgb\s+([^)]+)\)$/.exec(s))) { const p = m[1].split(/[\s\/]+/).filter(Boolean).map(parseFloat); o = [p[0] * 255, p[1] * 255, p[2] * 255, p.length > 3 ? p[3] : 1] }
    else { ctx.clearRect(0, 0, 1, 1); ctx.fillStyle = '#000'; ctx.fillStyle = s; ctx.fillRect(0, 0, 1, 1); const d = ctx.getImageData(0, 0, 1, 1).data; o = [d[0], d[1], d[2], d[3] / 255] }
    memo.set(s, o); return o
  }
  const sobre = (t, b) => { const a = t[3] + b[3] * (1 - t[3]); return a === 0 ? [0, 0, 0, 0] : [0, 1, 2].map((i) => (t[i] * t[3] + b[i] * b[3] * (1 - t[3])) / a).concat([a]) }
  const redondea = (c) => [Math.round(c[0]), Math.round(c[1]), Math.round(c[2])]
  const lin = (v) => { v /= 255; return v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4) }
  const lum = (c) => 0.2126 * lin(c[0]) + 0.7152 * lin(c[1]) + 0.0722 * lin(c[2])
  const razon = (a, b) => { const x = lum(a), y = lum(b); return (Math.max(x, y) + 0.05) / (Math.min(x, y) + 0.05) }
  const hex = (c) => '#' + c.map((v) => v.toString(16).padStart(2, '0')).join('').toUpperCase()
  function opacidad(el) { let o = 1; for (let n = el; n && n.nodeType === 1; n = n.parentElement) o *= parseFloat(getComputedStyle(n).opacity); return o }
  // Fondo efectivo de `el` (incluido su propio fondo si propio=true). indeterminado = hay imagen o degradado antes de una capa opaca.
  function fondo(el, propio) {
    const capas = []; let indet = false; let acum = 0
    for (let n = propio ? el : el.parentElement; n && n.nodeType === 1; n = n.parentElement) {
      const cs = getComputedStyle(n)
      if (cs.backgroundImage && cs.backgroundImage !== 'none') indet = true
      const c = parse(cs.backgroundColor)
      if (c[3] > 0) { capas.push(c); acum = acum + c[3] * (1 - acum); if (c[3] >= 0.999) break }
    }
    let f = [255, 255, 255, 1]
    if (acum < 0.999) { const cv = getComputedStyle(document.documentElement).backgroundColor; const base = parse(cv); f = base[3] > 0 ? base : (matchMedia('(prefers-color-scheme: dark)').matches && getComputedStyle(document.documentElement).colorScheme.includes('dark') ? [18, 18, 18, 1] : [255, 255, 255, 1]) }
    for (let i = capas.length - 1; i >= 0; i--) f = sobre(capas[i], f)
    return { c: f, indet }
  }
  function sel(el) {
    const p = []; let n = el
    for (let i = 0; n && n.nodeType === 1 && i < 4; i++) {
      let s = n.tagName.toLowerCase()
      if (n.id) { p.unshift(s + '#' + n.id); break }
      const dp = n.getAttribute('data-prueba')
      if (dp) s += '[data-prueba="' + dp + '"]'; else { const c = [...n.classList].slice(0, 2).join('.'); if (c) s += '.' + c }
      p.unshift(s); n = n.parentElement
    }
    return p.join(' > ')
  }
  const visible = (el, cs) => { if (cs.display === 'none' || cs.visibility === 'hidden' || cs.visibility === 'collapse') return false; const r = el.getBoundingClientRect(); return r.width >= 2 && r.height >= 2 && opacidad(el) > 0.01 }
  const deshab = (el) => !!el.closest(':disabled,[aria-disabled="true"],fieldset:disabled')
  const grupos = new Map(); const bajos = []; const indeterminados = []
  const excl = { deshabilitados: 0, ui_sin_limite_visible: 0, graficos_marcas_no_medibles_en_dom: 0 }
  let evaluados = 0
  function registrar(tipo, el, fg, bg, umbral, grande, texto, indet, sub) {
    const fgc = redondea(fg), bgc = redondea(bg.c)
    const r = razon(fgc, bgc); const clave = [tipo, sub || '', hex(fgc), hex(bgc), umbral].join('|')
    let g = grupos.get(clave)
    if (!g) { g = { tipo, subtipo: sub || null, fg: hex(fgc), bg: hex(bgc), umbral, grande, razon: Math.round(r * 1000) / 1000, elementos: 0, ejemplos: [] }; grupos.set(clave, g) }
    g.elementos++; evaluados++
    const e = { nodo: sel(el), texto: (texto || '').trim().slice(0, 40) }
    if (g.ejemplos.length < 3) g.ejemplos.push(e)
    if (indet) indeterminados.push({ tipo, subtipo: sub || null, ...e, fg: hex(fgc), bg: hex(bgc), razon: g.razon })
    if (r < umbral) bajos.push({ tipo, subtipo: sub || null, ...e, fg: hex(fgc), bg: hex(bgc), razon: g.razon, umbral })
  }
  const enGrafica = (el) => !!el.closest('[_echarts_instance_]')
  const todos = [raiz, ...raiz.querySelectorAll('*')]
  for (const el of todos) {
    const tag = el.tagName.toLowerCase()
    if (['script', 'style', 'noscript', 'canvas', 'option', 'path', 'g', 'rect', 'circle', 'line', 'polyline', 'polygon', 'ellipse', 'defs', 'clippath', 'tspan', 'title', 'desc'].includes(tag)) continue
    const cs = getComputedStyle(el)
    if (!visible(el, cs)) continue
    const alfaOp = opacidad(el)
    // Texto propio (nodos de texto directos) o valor de un campo.
    let texto = ''
    for (const nodo of el.childNodes) if (nodo.nodeType === 3) texto += nodo.textContent
    if ((tag === 'input' && !['checkbox', 'radio', 'button', 'submit', 'hidden', 'file', 'range', 'color', 'image'].includes(el.type) || tag === 'textarea') && el.value) texto = el.value
    if (tag === 'select') texto = el.selectedOptions[0] ? el.selectedOptions[0].textContent : ''
    if (tag === 'text' ) texto = el.textContent
    if (texto.trim()) {
      if (deshab(el)) { excl.deshabilitados++ }
      else {
        const bg = fondo(el, tag !== 'text' && tag !== 'svg' && true)
        const c = parse(cs.color && tag !== 'text' ? cs.color : (cs.fill !== 'none' ? cs.fill : cs.color))
        const fgC = sobre([c[0], c[1], c[2], c[3] * alfaOp], bg.c)
        const fs = parseFloat(cs.fontSize), fw = parseInt(cs.fontWeight, 10) || 400
        const grande = fs >= 24 || (fs >= 18.66 && fw >= 700)
        registrar('texto', el, fgC, bg, grande ? 3 : 4.5, grande, texto, bg.indet)
      }
    }
    // Íconos (svg sueltos, no gráficas).
    if (tag === 'svg' && !enGrafica(el) && !el.closest('svg svg') && el.getBoundingClientRect().width >= 8) {
      const trazo = cs.stroke !== 'none' && cs.stroke ? cs.stroke : cs.fill
      if (trazo && trazo !== 'none') {
        if (deshab(el)) { excl.deshabilitados++ } else {
          const bg = fondo(el, false); const c = parse(trazo)
          registrar('icono', el, sobre([c[0], c[1], c[2], c[3] * alfaOp], bg.c), bg, 3, false, el.getAttribute('aria-label') || el.getAttribute('class') || '', bg.indet)
        }
      }
    }
    // Componentes de UI con límite visible (borde o relleno distinto del fondo que los rodea).
    if (el.matches('button,[role=button],[role=tab],[role=checkbox],[role=switch],[role=radio],[role=combobox],input:not([type=hidden]),select,textarea')) {
      if (deshab(el)) { excl.deshabilitados++; continue }
      const padre = fondo(el, false); const propio = fondo(el, true)
      let mejor = 0; let pareja = null
      const lados = ['Top', 'Right', 'Bottom', 'Left']
      for (const l of lados) {
        if (parseFloat(cs['border' + l + 'Width']) >= 1 && cs['border' + l + 'Style'] !== 'none') {
          const b = parse(cs['border' + l + 'Color']); if (b[3] <= 0) continue
          const bc = sobre(b, padre.c); const r = razon(redondea(bc), redondea(padre.c)); if (r > mejor) { mejor = r; pareja = [bc] }
        }
      }
      const rf = razon(redondea(propio.c), redondea(padre.c)); if (rf > mejor && rf > 1.001) { mejor = rf; pareja = [propio.c] }
      const sombra = cs.boxShadow !== 'none' && /0px 0px 0px [1-9]/.test(cs.boxShadow)
      if (!pareja && !sombra) { excl.ui_sin_limite_visible++; continue }
      // subtipo: campo = su límite es lo que lo identifica (WCAG 1.4.11 lo exige); seleccion = botón con texto cuyo estado se marca con relleno o borde; boton = botón con texto.
      const rol = el.getAttribute('role') || ''
      const sub = (tag === 'input' || tag === 'select' || tag === 'textarea' || ['checkbox', 'radio', 'switch', 'combobox'].includes(rol)) ? 'campo'
        : (rol === 'tab' || el.hasAttribute('aria-pressed') || el.hasAttribute('aria-selected') || el.hasAttribute('aria-checked') || el.hasAttribute('aria-current') || el.hasAttribute('data-state')) ? 'seleccion' : 'boton'
      if (pareja) registrar('ui', el, pareja[0], padre, 3, false, el.getAttribute('aria-label') || el.textContent || el.getAttribute('name') || '', padre.indet || propio.indet, sub)
    }
  }
  if (document.querySelector('[_echarts_instance_]')) excl.graficos_marcas_no_medibles_en_dom = raiz.querySelectorAll('[_echarts_instance_] path, [_echarts_instance_] rect').length
  const pares = [...grupos.values()].sort((a, b) => a.razon - b.razon)
  return { alcance: alcanceSel || 'body', evaluados, pares_unicos: pares.length, bajo_umbral: bajos.length, indeterminados: indeterminados.length, excluidos: excl, bajos: bajos.slice(0, 200), indeterminados_detalle: indeterminados.slice(0, 50), pares }
})
