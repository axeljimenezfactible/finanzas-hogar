// «npm run tokens»: copia TAL CUAL los tokens de C2 y las fuentes Inter autoalojadas.
// Nunca se editan a mano en frontend/: la fuente es docs/contratos y docs/design/fuentes.
import { copyFileSync, mkdirSync, readFileSync } from 'node:fs'
import { createHash } from 'node:crypto'
import { desdeRepo, desdeFront } from './comun.mjs'
const sha = (f) => createHash('sha256').update(readFileSync(f)).digest('hex').slice(0, 12)
const copias = [
  [desdeRepo('docs/contratos/tokens.css'), desdeFront('src/estilos/tokens.css')],
  [desdeRepo('docs/contratos/tailwind.theme.css'), desdeFront('src/estilos/tailwind.theme.css')],
  [desdeRepo('docs/design/fuentes/inter/web/InterVariable.woff2'), desdeFront('public/fonts/InterVariable.woff2')],
  [desdeRepo('docs/design/fuentes/inter/web/InterVariable-Italic.woff2'), desdeFront('public/fonts/InterVariable-Italic.woff2')],
  [desdeRepo('docs/design/fuentes/inter/LICENSE.txt'), desdeFront('public/fonts/LICENSE-Inter.txt')],
]
mkdirSync(desdeFront('public/fonts'), { recursive: true })
for (const [de, a] of copias) {
  copyFileSync(de, a)
  if (sha(de) !== sha(a)) throw new Error(`copia distinta: ${a}`)
  console.log(`tokens: ${a.replace(desdeFront(''), '')} = ${sha(a)}`)
}
