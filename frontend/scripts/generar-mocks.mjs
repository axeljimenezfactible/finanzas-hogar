// «npm run mocks:gen»: una entrada por operación de C1 (openapi.yaml) con su respuesta de éxito
// armada con los ejemplos del contrato. La usa src/mocks/manejadores.ts (solo «npm run dev:mocks»).
import { readFileSync, writeFileSync } from 'node:fs'
import { parse } from 'yaml'
import { desdeRepo, desdeFront } from './comun.mjs'
const doc = parse(readFileSync(desdeRepo('docs/contratos/openapi.yaml'), 'utf8'))
const ref = (r) => r.split('/').slice(1).reduce((o, k) => o[k], doc)
function muestra(s, prof = 0) {
  if (!s || prof > 7) return null
  if (s.$ref) return muestra(ref(s.$ref), prof + 1)
  if (s.examples?.length) return s.examples[0]
  if (s.example !== undefined) return s.example
  if (s.default !== undefined) return s.default
  if (s.const !== undefined) return s.const
  if (s.enum) return s.enum[0]
  if (s.allOf) return Object.assign({}, ...s.allOf.map((x) => muestra(x, prof + 1)))
  if (s.oneOf || s.anyOf) return muestra((s.oneOf || s.anyOf)[0], prof + 1)
  const tipo = Array.isArray(s.type) ? s.type.find((t) => t !== 'null') : s.type
  if (tipo === 'object' || s.properties) return Object.fromEntries(Object.entries(s.properties || {}).map(([k, v]) => [k, muestra(v, prof + 1)]))
  if (tipo === 'array') return s.items ? [muestra(s.items, prof + 1)] : []
  if (tipo === 'integer' || tipo === 'number') return s.minimum ?? 1
  if (tipo === 'boolean') return false
  if (tipo === 'string') return s.format === 'date' ? '2026-10-21' : s.format === 'date-time' ? '2026-10-21T09:00:00-06:00' : 'texto'
  return null
}
const ops = []
for (const [ruta, item] of Object.entries(doc.paths)) for (const [metodo, op] of Object.entries(item)) {
  if (!op?.operationId) continue
  const estado = ['200', '201', '202', '204'].find((c) => op.responses?.[c])
  const r = estado && op.responses[estado]
  const esquema = r?.content?.['application/json']?.schema
  ops.push({ id: op.operationId, metodo, ruta, estado: Number(estado || 200), ejemplo: esquema ? muestra(esquema) : null, binario: !!r?.content && !esquema })
}
// Datos de ejemplo para las listas que usa la base (personas y categorías) desde datos-sinteticos.json.
const sint = JSON.parse(readFileSync(desdeRepo('docs/design/maquetas/datos-sinteticos.json'), 'utf8'))
const ahora = '2026-10-21T09:00:00-06:00'
const PART = { Titular: 'holder', 'Co-titular': 'co_holder' }
const personas = sint.personas.map((p, i) => ({ id: i + 1, code: p.id, display_name: p.nombre, participation: PART[p.rol] ?? 'external', pays_share_directly: false, is_active: true, has_account: !!p.tiene_cuenta, created_at: ahora, updated_at: ahora }))
const categorias = sint.categorias.map((c, i) => ({ id: i + 1, code: c.id, name: c.nombre, kind: c.tipo === 'fijo' ? 'fixed' : 'variable', sort_order: c.orden, is_active: c.activa, expense_count: c.num_gastos ?? 0, recent_usage_count: c.num_gastos ?? 0, created_at: ahora, updated_at: ahora }))
writeFileSync(desdeFront('src/mocks/sinteticos.gen.json'), JSON.stringify({ personas, categorias }) + '\n')
writeFileSync(desdeFront('src/mocks/operaciones.gen.json'), JSON.stringify(ops) + '\n')
console.log(`mocks: ${ops.length} operaciones de C1 ${doc.info.version}`)
