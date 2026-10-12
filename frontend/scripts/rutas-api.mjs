// C3-03: rutas de API que invoca el frontend contra openapi.yaml (0 fuera del contrato) y contra
// backend/internal/api/routes.go (rutas registradas). Solo cuenta llamadas por el cliente (pedir/pedirPagina/todas/encolar).
import { readFileSync, readdirSync, statSync } from 'node:fs'
import path from 'node:path'
import { parse } from 'yaml'
import { desdeRepo, desdeFront } from './comun.mjs'
const doc = parse(readFileSync(desdeRepo('docs/contratos/openapi.yaml'), 'utf8'))
const contrato = new Set(Object.entries(doc.paths).flatMap(([r, it]) => Object.keys(it).filter((m) => it[m]?.operationId).map((m) => `${m.toUpperCase()} ${r}`)))
let registradas = new Set()
try { registradas = new Set([...readFileSync(desdeRepo('backend/internal/api/routes.go'), 'utf8').matchAll(/Method:\s*"(\w+)",\s*Pattern:\s*"([^"]+)"/g)].map((m) => `${m[1]} ${m[2]}`)) } catch { /* sin routes.go */ }
const archivos = []
const recorrer = (d) => { for (const n of readdirSync(d)) { const f = path.join(d, n); if (statSync(f).isDirectory()) { if (n !== 'mocks') recorrer(f) } else if (/\.(ts|tsx)$/.test(n) && !n.endsWith('.gen.ts')) archivos.push(f) } }
recorrer(desdeFront('src'))
const usadas = new Map()
const re = /\b(?:pedir|pedirPagina|todas|encolar|usarConsulta)\(\s*'(get|post|put|patch|delete)'\s*,\s*'([^']+)'/g
for (const f of archivos) for (const m of readFileSync(f, 'utf8').matchAll(re)) usadas.set(`${m[1].toUpperCase()} ${m[2]}`, path.relative(desdeFront(''), f))
const fuera = [...usadas.keys()].filter((k) => !contrato.has(k))
const faltan = registradas.size ? [...usadas.keys()].filter((k) => !registradas.has(k)) : []
console.log(`C3-03 frontend: ${usadas.size} rutas invocadas; fuera de openapi.yaml: ${fuera.length}; sin registrar en routes.go: ${faltan.length} (registradas: ${registradas.size})`)
for (const k of fuera) console.log(`  FUERA ${k} (${usadas.get(k)})`)
for (const k of faltan) console.log(`  SIN-REGISTRAR ${k} (${usadas.get(k)})`)
process.exit(fuera.length || faltan.length ? 1 : 0)
