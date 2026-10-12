// Apoyo de las pruebas U5: lista T de pantallas desde U1 (por script), rutas, sesión y colores de C2.
import fs from 'node:fs'
import { fileURLToPath } from 'node:url'
import os from 'node:os'
import path from 'node:path'
import type { Browser, BrowserContext, Page } from '@playwright/test'

export const AQUI = path.dirname(fileURLToPath(import.meta.url))
export const REPO = path.resolve(AQUI, '..', '..')
export const SALIDA = path.join(REPO, 'docs', 'design', 'capturas-u5')
export const EVID = path.join(SALIDA, 'evidencia')
export const PARTES = path.join(REPO, 'datos', 'tmp', 'u5-partes')
export const ETIQUETAS_AXE = ['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'wcag22aa']

/** T: todos los P-xx y P-26.xx de docs/design/U1-arquitectura.md (el inventario que declara U1-02), con su nombre de la tabla §2.1/§2.3. */
export function pantallasU1(): { id: string; nombre: string }[] {
  const md = fs.readFileSync(path.join(REPO, 'docs', 'design', 'U1-arquitectura.md'), 'utf8')
  const ids = [...new Set([...md.matchAll(/\bP-\d{2}(?:\.\d{2})?\b/g)].map((m) => m[0]))].sort()
  return ids.map((id) => {
    const fila = md.split('\n').find((l) => l.startsWith(`| ${id} | `))
    return { id, nombre: fila ? fila.split('|')[2].trim() : '' }
  })
}
/** Ruta de cada pantalla. P-03 y P-22 son hojas: se abren desde su pantalla de origen. */
export function ruta(id: string, ids: { gasto: string; receta: string; credito: string; compra: string }): string {
  if (/^P-26\.\d\d$/.test(id)) return `/ajustes/${id.slice(5)}`
  const m: Record<string, string> = {
    'P-00': '/entrar', 'P-01': '/inicio', 'P-02': '/gastos', 'P-03': '/inicio', 'P-04': `/gastos/${ids.gasto}`, 'P-05': '/gastos/ticket',
    'P-06': '/presupuesto', 'P-07': '/presupuesto/flujo', 'P-08': '/presupuesto/ingresos', 'P-09': '/presupuesto/ahorro',
    'P-10': '/comidas', 'P-11': '/comidas/recetas', 'P-12': `/comidas/recetas/${ids.receta}`,
    'P-13': '/compras', 'P-14': '/compras/comparar', 'P-15': '/compras/productos', 'P-16': '/compras/inventario',
    'P-17': '/creditos', 'P-18': `/creditos/${ids.credito}`, 'P-19': `/creditos/${ids.credito}/compras/${ids.compra}`, 'P-20': '/creditos/comparador',
    'P-21': '/salud', 'P-22': '/salud', 'P-23': '/reportes', 'P-24': '/reportes/mensual', 'P-25': '/reportes/karen', 'P-26': '/ajustes',
  }
  if (!m[id]) throw new Error(`Sin ruta para ${id}`)
  return m[id]
}
/** Nombres de rutas.ts (frontend) para comprobar el h1. */
export function nombresRutas(): Record<string, string> {
  const t = fs.readFileSync(path.join(REPO, 'frontend', 'src', 'app', 'rutas.ts'), 'utf8')
  const o: Record<string, string> = {}
  for (const m of t.matchAll(/id: '(P-\d\d)', nombre: '([^']+)'/g)) o[m[1]] = m[2]
  return o
}
function credenciales(): { usuario: string; clave: string } {
  // La contraseña de desarrollo vive FUERA del repositorio (modo 600). Nunca se imprime ni se guarda en evidencias.
  const f = process.env.FINANZAS_SEED_PASSWORD_FILE ?? path.join(os.homedir(), '.config', 'finanzas-hogar', 'dev-semilla.txt')
  const linea = fs.readFileSync(f, 'utf8').split('\n').find((l) => l.startsWith('axel='))
  if (!linea) throw new Error('No hay renglón axel= en el archivo de contraseñas de desarrollo')
  return { usuario: 'axel', clave: linea.slice('axel='.length).trim() }
}
type Uso = Record<string, unknown>
const OPC = ['baseURL', 'viewport', 'hasTouch', 'isMobile', 'colorScheme', 'locale', 'timezoneId', 'deviceScaleFactor', 'ignoreHTTPSErrors']
export const opcionesContexto = (uso: Uso) => Object.fromEntries(OPC.filter((k) => uso[k] !== undefined).map((k) => [k, uso[k]]))
export interface Sesion { context: BrowserContext; page: Page; ids?: { gasto: string; receta: string; credito: string; compra: string } }
export async function abrirAnonima(browser: Browser, uso: Uso): Promise<Sesion> {
  const context = await browser.newContext(opcionesContexto(uso) as never)
  return { context, page: await context.newPage() }
}
/** Inicia sesión UNA vez por contexto (formulario real de P-00). */
export async function abrirSesion(browser: Browser, uso: Uso): Promise<Sesion> {
  const s = await abrirAnonima(browser, uso)
  const { usuario, clave } = credenciales()
  await s.page.goto('/entrar')
  await s.page.getByLabel('Usuario').fill(usuario)
  await s.page.getByLabel('Contraseña').fill(clave)
  await s.page.getByRole('button', { name: 'Entrar' }).click()
  await s.page.waitForURL(/\/inicio$/, { timeout: 30_000 })
  return s
}
/** Ids de datos sembrados para las pantallas con parámetro (P-04, P-12, P-18, P-19), leídos del API con la sesión abierta. */
export async function resolverIds(s: Sesion) {
  if (s.ids) return s.ids
  const get = async (p: string) => (await s.page.evaluate(async (u) => (await fetch(u, { credentials: 'include' })).json(), p)) as { items?: Record<string, unknown>[] }
  const gastos = (await get('/api/v1/expenses?limit=200')).items ?? []
  const compra = gastos.find((g) => g.amount_cents === 309450) ?? gastos[0]
  const recetas = (await get('/api/v1/recipes?limit=5')).items ?? []
  const creditos = (await get('/api/v1/credits?limit=5')).items ?? []
  const cid = String((creditos[0] as { id?: number } | undefined)?.id ?? 1)
  const compras = (await get(`/api/v1/credit-purchases?credit_id=${cid}&limit=5`)).items ?? []
  s.ids = { gasto: String(compra?.id ?? 1), receta: String((recetas[0] as { id?: number } | undefined)?.id ?? 1), credito: cid, compra: String((compras[0] as { id?: number } | undefined)?.id ?? 1) }
  return s.ids
}
/** Colores de C2 (docs/contratos/C2-tokens.json): fondo de página por tema, resuelto a rgb. */
export function fondoC2(tema: 'claro' | 'oscuro'): { rgb: string; hex: string; ref: string } {
  const t = JSON.parse(fs.readFileSync(path.join(REPO, 'docs', 'contratos', 'C2-tokens.json'), 'utf8'))
  const ref: string = t.ui.surface.page[tema === 'claro' ? 'light' : 'dark']['$value']
  const hex: string = ref.replace(/[{}]/g, '').split('.').reduce((o: Record<string, never>, k) => (o as never)[k], t as never)['$value' as never]
  const n = parseInt(hex.slice(1), 16)
  return { rgb: `rgb(${(n >> 16) & 255}, ${(n >> 8) & 255}, ${n & 255})`, hex, ref: `ui.surface.page.${tema === 'claro' ? 'light' : 'dark'} = ${ref}` }
}
/** Espera a que la pantalla termine de cargar (h1 visible, red en reposo, sin esqueletos) y a que acaben las animaciones de gráficas. */
export async function listo(page: Page): Promise<{ cargando: boolean }> {
  await page.locator('h1').first().waitFor({ state: 'visible', timeout: 30_000 })
  await page.waitForLoadState('networkidle', { timeout: 10_000 }).catch(() => undefined)
  let cargando = true
  for (let i = 0; i < 40 && cargando; i++) {
    cargando = (await page.locator('[aria-busy="true"]').count()) > 0
    if (cargando) await page.waitForTimeout(250)
  }
  await page.waitForTimeout(900)
  return { cargando }
}
