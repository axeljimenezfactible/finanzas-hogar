// Humo real (C3-09 parcial, S6-12): entra con la cuenta de desarrollo y abre Inicio contra el servidor de casa
// en marcha (ops/dev-run.sh) con la semilla. Sin context.route: todas las cifras vienen de la API real.
import { readFileSync } from 'node:fs'
import os from 'node:os'
import path from 'node:path'
import { expect, test } from '@playwright/test'
import { dinero } from '../src/api/formatos'

// La contraseña se lee del archivo de semilla de desarrollo (modo 600, fuera del repositorio); nunca se imprime.
function contrasenaDesarrollo(usuario: string): string {
  const archivo = process.env.FH_DEV_SEMILLA ?? path.join(os.homedir(), '.config', 'finanzas-hogar', 'dev-semilla.txt')
  for (const renglon of readFileSync(archivo, 'utf8').split('\n')) {
    const i = renglon.indexOf('=')
    if (!renglon.startsWith('#') && i > 0 && renglon.slice(0, i).trim() === usuario) return renglon.slice(i + 1).trim()
  }
  throw new Error(`No hay contraseña de desarrollo para «${usuario}» en el archivo de semilla`)
}

// Las 6 tarjetas de U3 en orden fijo y la operación de C1 de cada una.
const TARJETAS = [
  { n: 1, ruta: '/api/v1/dashboard/month-vs-plan', op: 'getDashboardMonthVsPlan' },
  { n: 2, ruta: '/api/v1/dashboard/fortnight-flow', op: 'getDashboardFortnightFlow' },
  { n: 3, ruta: '/api/v1/dashboard/next-payment', op: 'getDashboardNextPayment' },
  { n: 4, ruta: '/api/v1/dashboard/week-meals', op: 'getDashboardWeekMeals' },
  { n: 5, ruta: '/api/v1/dashboard/best-basket', op: 'getDashboardBestBasket' },
  { n: 6, ruta: '/api/v1/dashboard/savings', op: 'getDashboardSavings' },
]

// Todas las cifras en centavos (*_cents) de una respuesta, ya con el formato de pantalla (con y sin signo de gasto).
function cifras(o: unknown, salida = new Set<string>(), clave = ''): Set<string> {
  if (Array.isArray(o)) o.forEach((v) => cifras(v, salida, clave))
  else if (o && typeof o === 'object') for (const [k, v] of Object.entries(o)) cifras(v, salida, k)
  else if (typeof o === 'number' && Number.isSafeInteger(o) && clave.endsWith('_cents') && o !== 0) {
    salida.add(dinero(o)); salida.add(dinero(Math.abs(o)))
  }
  return salida
}

test('humo real · entrar como axel, Inicio con las 6 tarjetas desde la API y el disponible A $10,000.00', async ({ page }) => {
  const base = new URL(test.info().project.use.baseURL ?? '')
  const respuestas = new Map<string, unknown>()
  const ajenas = new Set<string>()
  page.on('request', (r) => { const u = new URL(r.url()); if (/^https?:$/.test(u.protocol) && u.host !== base.host) ajenas.add(u.host) })
  page.on('response', async (r) => {
    const u = new URL(r.url())
    const t = TARJETAS.find((x) => x.ruta === u.pathname)
    if (t && r.status() === 200 && u.host === base.host) respuestas.set(t.ruta, await r.json().catch(() => undefined))
  })

  await page.goto('/inicio')
  await expect(page.getByRole('heading', { name: 'Iniciar sesión' })).toBeVisible()
  await page.getByLabel('Usuario').fill('axel')
  await page.getByLabel('Contraseña').fill(contrasenaDesarrollo('axel'))
  await page.getByRole('button', { name: 'Entrar' }).click()
  await expect(page).toHaveURL(/\/inicio/)
  await expect(page.getByRole('heading', { level: 1, name: 'Inicio' })).toBeVisible()

  for (const t of TARJETAS) {
    const tarjeta = page.getByTestId(`tarjeta-${t.n}`)
    await expect(tarjeta, `tarjeta ${t.n} visible`).toBeVisible()
    await expect.poll(() => respuestas.has(t.ruta), { message: `${t.op} respondió 200 desde ${base.host}` }).toBe(true)
    // Cargada = la tarjeta es el enlace al detalle (no la caja de carga, vacío o error).
    await expect(tarjeta, `tarjeta ${t.n} cargada`).toHaveJSProperty('tagName', 'A')
    const texto = (await tarjeta.innerText()).replace(/\s+/g, ' ')
    const esperadas = [...cifras(respuestas.get(t.ruta))]
    const halladas = esperadas.filter((c) => texto.includes(c))
    test.info().annotations.push({ type: `tarjeta-${t.n} ${t.op}`, description: `cifras de la API: ${esperadas.length}; en pantalla: ${halladas.join(' | ') || '—'} · texto: ${texto.slice(0, 160)}` })
    if (esperadas.length) expect(halladas.length, `tarjeta ${t.n}: alguna cifra de ${t.op} en pantalla`).toBeGreaterThan(0)
  }

  // C3-09 (parcial), cifras de oraculo/model.json: próximo pago 973.75 y disponible de la lectura activa (B 6,905.50).
  await expect(page.getByTestId('tarjeta-3')).toContainText('$973.75')
  await expect(page.getByTestId('tarjeta-3')).toContainText('$6,905.50')

  // Disponible inicial A 10,000.00 (credito.eventos[0].disponible_A): Inicio → Detalle → «Disponible por línea».
  const vista = page.waitForResponse((r) => /\/api\/v1\/credits\/\d+\/availability-preview$/.test(new URL(r.url()).pathname) && r.status() === 200)
  await page.getByRole('tab', { name: 'Detalle' }).click()
  const dispo = (await (await vista).json()) as { reading_a: { available_today_cents: number }; reading_b: { available_today_cents: number } }
  expect(dispo.reading_a.available_today_cents, 'API: disponible lectura A').toBe(1_000_000)
  expect(dispo.reading_b.available_today_cents, 'API: disponible lectura B').toBe(690_550)
  await expect(page.getByText('$10,000.00').first()).toBeVisible()
  await expect(page.getByText(dinero(dispo.reading_b.available_today_cents)).first()).toBeVisible()
  expect([...ajenas], 'peticiones a hosts externos').toEqual([])
})
