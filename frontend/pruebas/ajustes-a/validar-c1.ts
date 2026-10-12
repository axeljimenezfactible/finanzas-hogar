// Valida contra C1 (docs/contratos/openapi.yaml 3.1, JSON Schema 2020-12) el tráfico falso de las pruebas de Ajustes (R2-01):
// el cuerpo JSON de cada petición y cada respuesta falsa con su código. Un código que C1 no declara para la operación es error.
// Lo usan ajustes-a y ajustes-b. Sin red: lee el contrato del repositorio.
import { readFileSync } from 'node:fs'
import { fileURLToPath } from 'node:url'
import { parse } from 'yaml'
import Ajv2020, { type ValidateFunction } from 'ajv/dist/2020.js'
type Obj = Record<string, unknown>
interface C1 { doc: Obj; ajv: Ajv2020; rutas: { re: RegExp; plantilla: string; params: number }[]; compiladas: Map<string, ValidateFunction> }
let c1: C1 | undefined
function cargar(): C1 {
  if (c1) return c1
  const doc = parse(readFileSync(fileURLToPath(new URL('../../../docs/contratos/openapi.yaml', import.meta.url)), 'utf8')) as Obj
  const ajv = new Ajv2020({ strict: false, validateFormats: false, validateSchema: false, allErrors: true })
  ajv.addSchema(doc, 'c1')
  const rutas = Object.keys(doc.paths as Obj).map((plantilla) => ({ plantilla, params: (plantilla.match(/\{/g) ?? []).length, re: new RegExp('^' + plantilla.replace(/\{[^}]+\}/g, '[^/]+') + '$') }))
  rutas.sort((a, b) => a.params - b.params)
  c1 = { doc, ajv, rutas, compiladas: new Map() }
  return c1
}
/** Sigue los $ref locales (#/components/…) de un objeto de OpenAPI (respuesta o cuerpo). */
function resolver(doc: Obj, o: unknown): Obj | undefined {
  let x = o as Obj | undefined
  while (x && typeof x.$ref === 'string') x = (x.$ref as string).split('/').slice(1).reduce<Obj>((a, k) => a[k.replace(/~1/g, '/').replace(/~0/g, '~')] as Obj, doc)
  return x
}
function validador(c: C1, clave: string, esquema: unknown): ValidateFunction {
  let v = c.compiladas.get(clave)
  if (!v) { v = c.ajv.compile(JSON.parse(JSON.stringify(esquema).replaceAll('"#/components/', '"c1#/components/'))); c.compiladas.set(clave, v) }
  return v
}
const json = (o: Obj | undefined) => ((o?.content as Obj | undefined)?.['application/json'] as Obj | undefined)?.schema
/** Errores de C1 de una petición falsa: `ruta` es la ruta concreta (sin consulta). Lista vacía = válida. */
export function erroresC1(metodo: string, ruta: string, estado: number, peticion: unknown, respuesta: unknown): string[] {
  const c = cargar(), m = metodo.toLowerCase()
  const r = c.rutas.find((x) => x.re.test(ruta) && (c.doc.paths as Record<string, Obj>)[x.plantilla][m])
  if (!r) return [`${metodo} ${ruta}: la operación no está en C1`]
  const op = (c.doc.paths as Record<string, Obj>)[r.plantilla][m] as Obj
  const nombre = `${metodo} ${r.plantilla} (${op.operationId as string})`
  const errores: string[] = []
  const cuerpo = json(resolver(c.doc, op.requestBody))
  if (cuerpo && peticion !== null && peticion !== undefined) {
    const v = validador(c, `${m} ${r.plantilla} cuerpo`, cuerpo)
    if (!v(peticion)) errores.push(`${nombre} cuerpo: ${c.ajv.errorsText(v.errors)}`)
  }
  const resp = resolver(c.doc, (op.responses as Obj)[String(estado)])
  if (!resp) return [...errores, `${nombre}: C1 no declara la respuesta ${estado} (declara ${Object.keys(op.responses as Obj).join(', ')})`]
  const esquema = json(resp)
  if (!esquema) { if (respuesta !== null && respuesta !== undefined) errores.push(`${nombre} ${estado}: C1 no declara cuerpo y la prueba manda uno`); return errores }
  const v = validador(c, `${m} ${r.plantilla} ${estado}`, esquema)
  if (!v(respuesta)) errores.push(`${nombre} ${estado}: ${c.ajv.errorsText(v.errors)}`)
  return errores
}
