// U5-08 (estático): hosts que aparecen en dist/. Distingue literales que nunca se piden por red
// (espacios de nombres XML y ligas de documentación dentro de mensajes de error de librerías).
import { readFileSync, readdirSync, statSync } from 'node:fs'
import path from 'node:path'
const dir = process.argv[2] || 'dist'
const PERMITIDOS = new Set(['localhost', '127.0.0.1', 'finanzas.local'])
// Literales que nunca se piden: xmlns de SVG (w3.org), ligas en mensajes de error de React (react.dev) y de
// workbox (bit.ly), comentario de licencia del CSS de Tailwind (tailwindcss.com) y texto de la licencia OFL
// de Inter que se distribuye junto a la fuente (github.com, scripts.sil.org, solo en LICENSE-Inter.txt).
const NO_PETICION = new Set(['www.w3.org', 'react.dev', 'reactjs.org', 'mswjs.io', 'bit.ly', 'tailwindcss.com', 'github.com', 'scripts.sil.org'])
const hallados = new Map()
const recorrer = (d) => { for (const n of readdirSync(d)) { const f = path.join(d, n); if (statSync(f).isDirectory()) recorrer(f); else if (/\.(js|css|html|webmanifest|json|svg|txt)$/.test(n)) for (const m of readFileSync(f, 'utf8').matchAll(/https?:\/\/([a-z0-9.-]+\.[a-z]{2,}|localhost)/gi)) { const h = m[1].toLowerCase(); hallados.set(h, (hallados.get(h) || new Set()).add(path.relative(dir, f))) } } }
recorrer(dir)
const externos = [...hallados.keys()].filter((h) => !PERMITIDOS.has(h) && !NO_PETICION.has(h))
console.log(`U5-08 dist: hosts externos que se pedirían = ${externos.length}`)
for (const [h, fs] of hallados) console.log(`  ${PERMITIDOS.has(h) ? 'permitido' : NO_PETICION.has(h) ? 'literal-no-peticion' : 'EXTERNO'} ${h} en ${[...fs].join(', ')}`)
process.exit(externos.length ? 1 : 0)
