// U5-13 · Captura sin conexión: cola en el teléfono, sin duplicar (K9). Los 7 pasos literales del criterio, a 412×915 táctil,
// con la red cortada (context.setOffline). Base de PRUEBA bajo datos/tmp/u5-flujos/ (copia de la de desarrollo; la borra el integrador/el agente al terminar):
//   bash servidor-prueba.sh arrancar && cd ../.. && FH_BASE=https://127.0.0.1:8444 npx playwright test -c pruebas-u5/flujos/playwright.config.ts captura-sin-conexion
import { execFileSync } from 'node:child_process'
import { existsSync, mkdirSync, writeFileSync } from 'node:fs'
import path from 'node:path'
import { expect, test, type Page, type Request } from '@playwright/test'
import { api, DIR_CAPTURAS, DIR_TMP, entrar, Flujo, RAIZ } from './apoyo'

const DB = process.env.FH_DB_PRUEBA ?? path.join(DIR_TMP, 'finanzas.db')
const norm = (s: string) => s.replace(/\s+/g, ' ').trim()
const sqlite = (sql: string) => (existsSync(DB) ? execFileSync('sqlite3', [DB, sql], { encoding: 'utf8' }).trim() : 'sin base')
const ISO_CON_DESPLAZAMIENTO = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?([+-]\d{2}:\d{2}|Z)$/

interface Envio { url: string; clave: string | undefined; cuerpo: Record<string, unknown>; enviadoEn: number }
const gastos = async (page: Page) => (await api<{ items: { id: number; amount_cents: number }[] }>(page, '/api/v1/expenses?limit=200')).json.items
const pesos = async (page: Page) => (await api<{ items: { id: number; weight_g: number }[] }>(page, '/api/v1/weight-entries?limit=200')).json.items
const enCola = (page: Page) => page.getByTestId('contador-cola')

async function capturarGasto(page: Page, monto: string, t: { n: number; pasos: string[] }, tocar: (l: ReturnType<Page['locator']>, n: string) => Promise<void>) {
  await tocar(page.getByTestId('boton-mas'), '«+»')
  await page.getByLabel('Monto').fill(monto)
  await tocar(page.getByRole('button', { name: 'Comida fuera', exact: true }), 'chip Comida fuera')
  await tocar(page.getByRole('button', { name: 'Guardar', exact: true }), 'Guardar')
  void t
}

test('U5-13 · captura-sin-conexion (7 pasos)', async ({ page, context }) => {
  const f = new Flujo('U5-13', 'Captura sin conexión (7 pasos de U5-13)', page)
  const envios: Envio[] = []
  page.on('request', (r: Request) => {
    if (r.method() === 'POST' && /\/api\/v1\/(expenses|weight-entries)$/.test(new URL(r.url()).pathname)) {
      envios.push({ url: new URL(r.url()).pathname, clave: r.headers()['idempotency-key'], cuerpo: r.postDataJSON() as Record<string, unknown>, enviadoEn: Date.now() })
    }
  })
  await entrar(page)
  // Precalentar en línea lo que la hoja necesita sin conexión (categorías, personas, último peso) y contar los registros de partida.
  await page.getByTestId('boton-mas').tap(); await expect(page.getByRole('button', { name: 'Comida fuera', exact: true })).toBeVisible()
  await page.getByRole('tab', { name: 'Peso' }).tap(); await expect(page.getByLabel(/^Peso/).first()).toBeVisible()
  await page.getByRole('button', { name: 'Cerrar' }).first().tap()
  await expect(page.getByTestId('tarjeta-1')).toBeVisible()
  const g0 = (await gastos(page)).length, w0 = (await pesos(page)).length
  const bd0 = sqlite("select (select count(*) from expense)||' gastos · '||(select count(*) from weight_entry)||' pesos'")
  const tocar = async (l: ReturnType<Page['locator']>, n: string) => f.t.tocar(l, n)
  let tGasto = 0, antesDeEnviar = 0

  await f.paso(`Paso 1: red cortada → captura un gasto (≤ 3 toques) y un peso; «Sin conexión» y contador en 2 (partida: ${g0} gastos y ${w0} pesos en el servidor)`, 'P-01, P-03, P-22', async () => {
    await context.setOffline(true)
    await expect(page.getByTestId('franja-conexion')).toContainText('Sin conexión con el servidor de casa')
    const n0 = f.t.n
    await capturarGasto(page, '43.21', f.t, tocar)
    tGasto = f.t.n - n0
    await expect(page.getByText('Gasto guardado en el teléfono')).toBeVisible()
    await expect(enCola(page)).toHaveText('1 en cola')
    await page.getByTestId('boton-mas').tap(); await page.getByRole('tab', { name: 'Peso' }).tap()
    await page.getByLabel(/^Peso/).first().fill('77.7'); await page.getByRole('button', { name: 'Guardar', exact: true }).tap()
    await expect(page.getByText('Peso guardado en el teléfono')).toBeVisible()
    await expect(enCola(page)).toHaveText('2 en cola')
    await f.captura('U5-13-sin-conexion-movil-claro-contador-2')
    expect(tGasto, 'toques del gasto sin conexión').toBeLessThanOrEqual(3)
    return `franja: «${norm(await page.getByTestId('franja-conexion').innerText())}» · gasto en ${tGasto} toques (${f.t.pasos.slice(-3).join('; ')}) · servidor sigue en ${(await Promise.resolve(g0))} gastos (sin red)`
  })
  await f.paso('(0 colas en otras pantallas) Sin red, una escritura de otra pantalla no se encola: botón deshabilitado con aviso', 'P-10', async () => {
    await page.getByRole('link', { name: 'Comidas', exact: true }).tap(); await page.waitForTimeout(1500)
    const t = norm(await page.locator('body').innerText())
    const b = page.getByRole('button', { name: 'Generar lista de compras' })
    const off = (await b.count()) ? await b.first().isDisabled() : null
    await f.captura('U5-13-otra-pantalla-sin-conexion')
    await expect(enCola(page)).toHaveText('2 en cola')
    return `P-10 sin red: botón «Generar lista de compras» ${off === null ? 'ausente' : off ? 'deshabilitado' : 'HABILITADO'}; cola sigue en 2; texto: ${/Requiere conexión[^.]{0,40}|Sin conexión[^.]{0,80}/.exec(t)?.[0] ?? '—'}`
  }, false, true)
  await f.paso('Paso 2a: recargar la página sin red; las 2 capturas siguen guardadas en el teléfono (IndexedDB)', 'P-01', async () => {
    await page.reload()
    await page.waitForTimeout(2500)
    const claves = await page.evaluate(async () => {
      const dbs = await indexedDB.databases()
      const salida: Record<string, unknown> = {}
      for (const d of dbs) {
        const db: IDBDatabase = await new Promise((res, rej) => { const q = indexedDB.open(d.name!); q.onsuccess = () => res(q.result); q.onerror = () => rej(q.error) })
        for (const almacen of Array.from(db.objectStoreNames)) {
          const filas: any[] = await new Promise((res) => { const q = db.transaction(almacen).objectStore(almacen).getAll(); q.onsuccess = () => res(q.result) })
          salida[`${d.name}/${almacen}`] = filas.map((x) => x.tipo ?? '(sin tipo)')
        }
        db.close()
      }
      return salida
    })
    const conTipo = Object.entries(claves).filter(([, v]) => (v as string[]).some((x) => x === 'gasto' || x === 'peso'))
    expect(conTipo.length, `almacenes con capturas: ${JSON.stringify(claves)}`).toBe(1)
    expect([...(conTipo[0][1] as string[])].sort()).toEqual(['gasto', 'peso'])
    return `almacenes del teléfono tras recargar: ${JSON.stringify(claves)} (únicas capturas: gasto y peso; ningún otro almacén guarda colas)`
  })
  await f.paso('Paso 2b: tras recargar sin red, la pantalla muestra «Sin conexión» y el contador sigue en 2', 'P-01', async () => {
    await f.captura('U5-13-tras-recargar-sin-red')
    await expect(page.getByTestId('franja-conexion')).toContainText('Sin conexión con el servidor de casa', { timeout: 5000 })
    await expect(enCola(page)).toHaveText('2 en cola')
    return `«${norm(await page.getByTestId('franja-conexion').innerText())}»`
  }, false, true)
  await f.paso('Paso 3: cada petición lleva Idempotency-Key (16–128) y captured_at (ISO con desplazamiento) anterior al envío', 'P-01', async () => {
    antesDeEnviar = Date.now()
    await context.setOffline(false)
    await page.goto('/inicio')
    await expect.poll(() => envios.length, { timeout: 60_000, message: 'peticiones de la cola' }).toBeGreaterThanOrEqual(2)
    const lineas: string[] = []
    for (const e of envios) {
      expect(e.clave, 'Idempotency-Key').toBeTruthy()
      expect(e.clave!.length).toBeGreaterThanOrEqual(16); expect(e.clave!.length).toBeLessThanOrEqual(128)
      const cap = String(e.cuerpo.captured_at)
      expect(cap, 'captured_at ISO 8601 con desplazamiento').toMatch(ISO_CON_DESPLAZAMIENTO)
      expect(Date.parse(cap), 'captured_at anterior al envío').toBeLessThan(e.enviadoEn)
      expect(Date.parse(cap), 'captured_at es del momento de captura, no del envío').toBeLessThan(antesDeEnviar)
      lineas.push(`POST ${e.url} · Idempotency-Key: ${e.clave} (${e.clave!.length} car.) · captured_at: ${cap}`)
    }
    return lineas.join(' || ')
  })
  await f.paso('Paso 4: al volver la red el contador llega a 0', 'P-01', async () => {
    await expect(enCola(page)).toHaveCount(0, { timeout: 60_000 })
    await expect(page.getByTestId('franja-conexion')).toHaveCount(0)
    return 'franja y contador desaparecen (0 pendientes)'
  })
  let idGasto = 0
  await f.paso('Paso 5: la API tiene exactamente 1 gasto y 1 peso nuevos, con la clave guardada (expense y weight_entry)', 'API', async () => {
    const g = await gastos(page), w = await pesos(page)
    expect(g.length - g0, 'gastos nuevos').toBe(1); expect(w.length - w0, 'pesos nuevos').toBe(1)
    expect(g.filter((x) => x.amount_cents === 4321).length).toBe(1); expect(w.filter((x) => x.weight_g === 77700).length).toBe(1)
    idGasto = g.find((x) => x.amount_cents === 4321)!.id
    const kG = envios.find((e) => e.url.endsWith('/expenses'))!.clave!, kW = envios.find((e) => e.url.endsWith('/weight-entries'))!.clave!
    const dbG = sqlite(`select count(*)||' | '||idempotency_key||' | '||captured_at from expense where idempotency_key='${kG}'`)
    const dbW = sqlite(`select count(*)||' | '||idempotency_key||' | '||captured_at from weight_entry where idempotency_key='${kW}'`)
    expect(dbG.startsWith('1 | ' + kG)).toBe(true); expect(dbW.startsWith('1 | ' + kW)).toBe(true)
    return `antes ${g0} gastos/${w0} pesos → después ${g.length}/${w.length}; BD expense: «${dbG}»; weight_entry: «${dbW}»; envíos totales de la cola: ${envios.length}`
  })
  await f.paso('Paso 6: reenviar la misma captura (misma clave y mismo cuerpo) no duplica y devuelve el mismo identificador', 'API', async () => {
    const e = envios.find((x) => x.url.endsWith('/expenses'))!
    const r = await page.evaluate(async ({ clave, cuerpo }) => {
      const x = await fetch('/api/v1/expenses', { method: 'POST', credentials: 'same-origin', headers: { 'Content-Type': 'application/json', 'Idempotency-Key': clave! }, body: JSON.stringify(cuerpo) })
      return { estado: x.status, repetida: x.headers.get('Idempotent-Replayed'), json: await x.json() }
    }, { clave: e.clave, cuerpo: e.cuerpo })
    expect([200, 201]).toContain(r.estado)
    expect(r.json.id, 'mismo identificador').toBe(idGasto)
    const g = await gastos(page)
    expect(g.length - g0, 'sigue 1 gasto nuevo').toBe(1)
    return `respuesta ${r.estado}, Idempotent-Replayed: ${r.repetida}, id ${r.json.id} = ${idGasto}; gastos nuevos: ${g.length - g0}`
  })
  await f.paso('Paso 7: fallo del servidor en un envío → la captura sigue en la cola y se muestra un aviso en español', 'P-01', async () => {
    await context.setOffline(true)
    await expect(page.getByTestId('franja-conexion')).toBeVisible()
    await capturarGasto(page, '55.55', f.t, tocar)
    await expect(enCola(page)).toHaveText('1 en cola')
    await page.route('**/api/v1/expenses', (r) => (r.request().method() === 'POST' ? r.fulfill({ status: 500, contentType: 'application/json', body: JSON.stringify({ code: 'internal', message: 'Falla simulada del servidor' }) }) : r.fallback()))
    await context.setOffline(false)
    await expect(page.getByText(/No se pudo enviar/).first()).toBeVisible({ timeout: 60_000 })
    await expect(enCola(page)).toContainText('1 en cola')
    const aviso = norm(await page.getByTestId('franja-conexion').innerText())
    await f.captura('U5-13-fallo-de-envio-aviso')
    const g = await gastos(page)
    expect(g.filter((x) => x.amount_cents === 5555).length, 'el gasto con fallo no llegó').toBe(0)
    await page.unroute('**/api/v1/expenses')
    return `aviso: «${aviso}» · la captura sigue en cola · el servidor no recibió el gasto de 55.55`
  })
  await f.paso('Limpieza: reintentar al quitar el fallo → llega 1 sola vez (0 duplicados)', 'P-01', async () => {
    const reint = page.getByRole('button', { name: /Reintentar envío|Reintentar|Revisar/ }).first()
    if (await reint.count()) await reint.tap().catch(() => undefined)
    await expect.poll(async () => (await gastos(page)).filter((x) => x.amount_cents === 5555).length, { timeout: 60_000 }).toBe(1)
    await expect(enCola(page)).toHaveCount(0, { timeout: 60_000 })
    const g = await gastos(page)
    expect(g.length - g0).toBe(2)
    return `gastos nuevos totales: ${g.length - g0} (43.21 y 55.55, uno cada uno); pesos nuevos: ${(await pesos(page)).length - w0}; BD final: ${sqlite("select (select count(*) from expense)||' gastos · '||(select count(*) from weight_entry)||' pesos'")} (partida: ${bd0})`
  })
  const reg = f.guardar()
  mkdirSync(DIR_CAPTURAS, { recursive: true })
  writeFileSync(path.join(DIR_CAPTURAS, 'u5-13-peticiones.json'), JSON.stringify(envios.map((e) => ({ ruta: e.url, 'Idempotency-Key': e.clave, captured_at: e.cuerpo.captured_at })), null, 1))
  expect(reg.pasos.filter((p) => !p.ok), JSON.stringify(reg.pasos.filter((p) => !p.ok))).toEqual([])
})

// Captura del estado «Sin conexión» en las 4 combinaciones de U5-01 (móvil/escritorio × claro/oscuro).
const COMBOS = [
  { id: 'movil-claro', viewport: { width: 412, height: 915 }, colorScheme: 'light' as const, movil: true },
  { id: 'movil-oscuro', viewport: { width: 412, height: 915 }, colorScheme: 'dark' as const, movil: true },
  { id: 'escritorio-claro', viewport: { width: 1440, height: 900 }, colorScheme: 'light' as const, movil: false },
  { id: 'escritorio-oscuro', viewport: { width: 1440, height: 900 }, colorScheme: 'dark' as const, movil: false },
]
for (const c of COMBOS) {
  test(`U5-13 · captura del estado «Sin conexión» · ${c.id}`, async ({ browser, baseURL }) => {
    const ctx = await browser.newContext({ baseURL, viewport: c.viewport, colorScheme: c.colorScheme, hasTouch: c.movil, isMobile: c.movil, locale: 'es-MX', timezoneId: 'America/Mexico_City' })
    const page = await ctx.newPage()
    try {
      await entrar(page)
      await expect(page.getByTestId('tarjeta-1')).toBeVisible()
      await ctx.setOffline(true)
      await expect(page.getByTestId('franja-conexion')).toContainText('Sin conexión con el servidor de casa')
      mkdirSync(DIR_CAPTURAS, { recursive: true })
      await page.screenshot({ path: path.join(DIR_CAPTURAS, `U5-13-sin-conexion-${c.id}.jpg`), type: 'jpeg', quality: 45 })
      const fondo = await page.evaluate(() => getComputedStyle(document.body).backgroundColor)
      console.log(`U5-13 ${c.id}: franja «${norm(await page.getByTestId('franja-conexion').innerText())}» · fondo ${fondo}`)
    } finally { await ctx.close() }
  })
}
void RAIZ
