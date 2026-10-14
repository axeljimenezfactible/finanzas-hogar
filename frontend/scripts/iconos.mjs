// «npm run iconos»: íconos de la PWA (192, 512 y enmascarable) con el ícono «wallet» de Lucide
// autoalojado (docs/design/fuentes/lucide/sprite.svg) y los colores de tokens.css (tema claro).
// También escribe src/estilos/colores-manifest.json para el manifest. Usa rsvg-convert local.
import { readFileSync, writeFileSync, mkdirSync } from 'node:fs'
import { execFileSync } from 'node:child_process'
import { desdeRepo, desdeFront } from './comun.mjs'
const css = readFileSync(desdeRepo('docs/contratos/tokens.css'), 'utf8')
const raiz = css.slice(css.indexOf(':root {'), css.indexOf('@media (prefers-color-scheme: dark)'))
const vars = Object.fromEntries([...raiz.matchAll(/(--[a-z0-9-]+):\s*([^;]+);/g)].map((m) => [m[1], m[2].trim()]))
const resolver = (v, n = 0) => { const m = /^var\((--[a-z0-9-]+)\)$/.exec(v); return m && n < 10 ? resolver(vars[m[1]], n + 1) : v }
const marca = resolver('var(--ui-action-primary-bg)'), sobreMarca = resolver('var(--ui-text-on-primary)'), fondo = resolver('var(--ui-surface-page)')
for (const c of [marca, sobreMarca, fondo]) if (!/^#[0-9A-Fa-f]{6}$/.test(c)) throw new Error(`color no resuelto: ${c}`)
const sprite = readFileSync(desdeRepo('docs/design/fuentes/lucide/sprite.svg'), 'utf8')
const simbolo = /<symbol[^>]*id="wallet"[^>]*>([\s\S]*?)<\/symbol>/.exec(sprite)?.[1]
if (!simbolo) throw new Error('falta el ícono wallet en sprite.svg')
const svg = (escala) => `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="512" height="512"><rect width="24" height="24" rx="${escala === 1 ? 5 : 0}" fill="${marca}"/><g transform="translate(12 12) scale(${0.62 * escala}) translate(-12 -12)" fill="none" stroke="${sobreMarca}" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">${simbolo}</g></svg>`
mkdirSync(desdeFront('public/icons'), { recursive: true })
writeFileSync(desdeFront('public/favicon.svg'), svg(1))
const png = (contenido, tam, salida) => {
  const tmp = desdeFront(`public/icons/.tmp-${tam}.svg`); writeFileSync(tmp, contenido)
  execFileSync('rsvg-convert', ['-w', String(tam), '-h', String(tam), '-o', desdeFront(salida), tmp])
  execFileSync('rm', ['-f', tmp])
}
png(svg(1), 192, 'public/icons/icon-192.png')
png(svg(1), 512, 'public/icons/icon-512.png')
png(svg(0.8), 512, 'public/icons/icon-maskable-512.png')
writeFileSync(desdeFront('src/estilos/colores-manifest.json'), JSON.stringify({ tema: marca, fondo }, null, 2) + '\n')
console.log(`iconos: marca ${marca}, fondo ${fondo}`)
