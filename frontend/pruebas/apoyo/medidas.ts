// Mediciones de DESIGN en el DOM real (U5-06): área de toque ≥ 44 × 44 px (§4.6 size.touch.min), iconos Lucide de
// 20 o 24 px con trazo de 1.5 px dibujados, 24 px en los botones de solo icono (§5.2) y cifras tabulares en todo
// dinero visible (§3.4).
import type { Page } from '@playwright/test'
export interface Control { control: string; ancho: number; alto: number }
export interface IconoMal { clase: string; ancho: number; alto: number; trazo: number; control: string; regla: string }
/** Controles visibles cuya área de toque es menor que 44 × 44 px. La etiqueta que envuelve una casilla cuenta como su
 *  área de toque (tocarla marca la casilla). Se omiten los ocultos (sr-only, fuera de pantalla, aria-hidden, inert). */
export function controlesChicos(page: Page): Promise<{ total: number; chicos: Control[] }> {
  return page.evaluate(() => {
    const SEL = 'a[href], button, input:not([type=hidden]), select, textarea, summary, [role=button], [role=tab], [role=checkbox], [role=switch], [role=radio], [role=menuitem], [role=option], [role=slider]'
    const chicos: { control: string; ancho: number; alto: number }[] = []
    let total = 0
    for (const el of Array.from(document.querySelectorAll<HTMLElement>(SEL))) {
      const r = el.getBoundingClientRect(), cs = getComputedStyle(el)
      if (r.width <= 1 || r.height <= 1 || cs.visibility === 'hidden' || el.closest('[aria-hidden="true"], [inert]')) continue
      if (r.right <= 0 || r.left >= window.innerWidth) continue
      if (cs.clipPath === 'inset(50%)' || cs.clip === 'rect(0px, 0px, 0px, 0px)') continue
      const etiqueta = (el as HTMLInputElement).labels?.[0]
      const caja = etiqueta && etiqueta.contains(el) ? etiqueta.getBoundingClientRect() : r
      total++
      if (caja.width + 0.5 < 44 || caja.height + 0.5 < 44) {
        const nombre = el.getAttribute('aria-label') || el.textContent?.trim().slice(0, 40) || el.getAttribute('name') || ''
        chicos.push({ control: `${el.tagName.toLowerCase()}${el.getAttribute('role') ? `[${el.getAttribute('role')}]` : ''} «${nombre}»`, ancho: Math.round(caja.width), alto: Math.round(caja.height) })
      }
    }
    return { total, chicos }
  })
}
/** Iconos Lucide visibles fuera de §5.2:
 *  - tamaño distinto de 20 o 24 px, o trazo dibujado distinto de 1.5 px (stroke-width en unidades del viewBox de 24 ×
 *    ancho / 24);
 *  - regla de uso: el icono de un botón o enlace sin texto visible (solo icono, con su nombre en aria-label o en texto
 *    sr-only) mide 24 px (`size.icon.md`: «solos en barras… botones de solo icono»). Queda fuera la barra lateral
 *    («Navegación principal»), a la que §5.2 asigna `size.icon.sm` en todas sus filas, incluidas las flechas
 *    Desplegar/Plegar de sección. */
export function iconosFueraDeNorma(page: Page): Promise<{ total: number; mal: IconoMal[] }> {
  return page.evaluate(() => {
    const CONTROL = 'button, a[href], [role=button], [role=link], [role=tab], [role=menuitem], [role=option]'
    const oculto = (el: Element) => {
      const cs = getComputedStyle(el), r = el.getBoundingClientRect()
      return cs.display === 'none' || cs.visibility === 'hidden' || cs.clipPath === 'inset(50%)' || cs.clip === 'rect(0px, 0px, 0px, 0px)' || r.width <= 1 || r.height <= 1
    }
    /** ¿El control muestra texto? Cuenta solo el texto que se ve (no el sr-only, ni el oculto, ni el de los svg). */
    const conTexto = (c: Element) => {
      const w = document.createTreeWalker(c, NodeFilter.SHOW_TEXT)
      for (let n = w.nextNode(); n; n = w.nextNode()) {
        if (!n.textContent?.trim()) continue
        let visible = true
        for (let p = n.parentElement; p && visible; p = p === c ? null : p.parentElement) if (p.closest('svg') || oculto(p)) visible = false
        if (visible) return true
      }
      return false
    }
    const nombre = (c: Element) => c.getAttribute('aria-label') || c.getAttribute('title') || c.textContent?.trim().slice(0, 40) || ''
    const mal: { clase: string; ancho: number; alto: number; trazo: number; control: string; regla: string }[] = []
    let total = 0
    for (const svg of Array.from(document.querySelectorAll<SVGSVGElement>('svg.lucide'))) {
      const r = svg.getBoundingClientRect()
      if (r.width <= 1 || r.height <= 1 || svg.closest('[inert]')) continue
      total++
      const vb = svg.viewBox.baseVal?.width || 24
      const dibujado = (Number(svg.getAttribute('stroke-width')) * r.width) / vb
      const ancho = Math.round(r.width * 10) / 10, alto = Math.round(r.height * 10) / 10
      const control = svg.closest(CONTROL)
      const fila = { clase: svg.getAttribute('class') ?? '', ancho, alto, trazo: Math.round(dibujado * 100) / 100, control: control ? `${control.tagName.toLowerCase()} «${nombre(control)}»` : '' }
      if (!((ancho === 20 || ancho === 24) && alto === ancho && Math.abs(dibujado - 1.5) < 0.01)) mal.push({ ...fila, regla: 'tamaño o trazo' })
      else if (control && ancho !== 24 && !control.closest('nav[aria-label="Navegación principal"]') && !conTexto(control)) mal.push({ ...fila, regla: 'solo icono → 24 px' })
    }
    return { total, mal }
  })
}
/** Textos visibles con dinero ($1,234.56) cuyo elemento no tiene font-variant-numeric: tabular-nums. Se omiten las
 *  gráficas (ECharts dibuja en canvas o SVG y sus ejes compactos los permite DESIGN §7.7). */
export function dineroSinTabulares(page: Page): Promise<string[]> {
  return page.evaluate(() => {
    const fuera: string[] = []
    const w = document.createTreeWalker(document.body, NodeFilter.SHOW_TEXT)
    for (let n = w.nextNode(); n; n = w.nextNode()) {
      const txt = n.textContent ?? ''
      if (!/\$\d/.test(txt)) continue
      const el = n.parentElement
      if (!el || el.closest('svg, canvas, option, script, style, [aria-hidden="true"]')) continue
      const r = el.getBoundingClientRect()
      if (r.width <= 1 || r.height <= 1) continue
      if (!getComputedStyle(el).fontVariantNumeric.includes('tabular-nums')) fuera.push(txt.trim().slice(0, 80))
    }
    return fuera
  })
}
