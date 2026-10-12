// U5-13 / U5-09 (base): cola sin conexión de gastos y peso (K9), 3 toques (K10), sin duplicar.
// U5-13 paso 7: si el envío falla (422, 500 o 401) la captura sigue en la cola, la franja la cuenta y se avisa.
import { expect, test, type Locator } from '@playwright/test'
import { ApiFalsa } from './apoyo/api-falsa'
const ISO_CON_DESPLAZAMIENTO = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d{1,6})?[+-]\d{2}:\d{2}$/
test('captura sin conexión: «Sin conexión», contador, recarga, envío al volver la red, sin duplicar', async ({ page, context }) => {
  const api = new ApiFalsa()
  await api.instalar(context)
  let toques = 0
  const tocar = async (l: Locator) => { await l.tap(); toques++ }
  await page.goto('/inicio')
  await page.evaluate(() => navigator.serviceWorker.ready.then(() => undefined))
  await page.reload() // el service worker ya controla la página: la app abre sin red
  await expect(page.getByRole('heading', { level: 1, name: 'Inicio' })).toBeVisible()
  await expect.poll(() => api.peticiones.some((p) => p.url.startsWith('/api/v1/categories'))).toBe(true)
  await context.setOffline(true)
  await expect(page.getByText('Sin conexión con el servidor de casa')).toBeVisible()
  // F-1 sin conexión: «+», categoría, Guardar = 3 toques (teclear no cuenta).
  await tocar(page.getByRole('button', { name: 'Registrar gasto' }).last())
  await expect(page.getByLabel('Monto')).toBeFocused()
  await page.keyboard.type('185.50')
  await tocar(page.getByRole('button', { name: 'Despensa' }))
  await tocar(page.getByRole('button', { name: 'Guardar' }))
  console.log(`U5-09 toques para registrar un gasto (sin conexión): ${toques}`)
  expect(toques).toBeLessThanOrEqual(3)
  await expect(page.getByRole('status').filter({ hasText: 'Gasto guardado en el teléfono' })).toBeVisible()
  await expect(page.getByText('1 en cola')).toBeVisible()
  // F-6 sin conexión: «+», Peso, Guardar.
  toques = 0
  await tocar(page.getByRole('button', { name: 'Registrar gasto' }).last())
  await tocar(page.getByRole('tab', { name: 'Peso' }))
  await page.getByLabel('Peso (kg)').fill('72.5')
  await tocar(page.getByRole('button', { name: 'Guardar' }))
  expect(toques).toBeLessThanOrEqual(3)
  await expect(page.getByText('2 en cola')).toBeVisible()
  await page.screenshot({ path: '../datos/tmp/playwright-frontend/sin-conexion.png' })
  // La cola sobrevive a recargar la página sin red.
  await page.reload()
  await expect(page.getByText('2 en cola')).toBeVisible()
  expect(api.gastos.length + api.pesos.length).toBe(0)
  // Vuelve la red: la cola se envía sola, el contador baja a 0, un registro por captura.
  await context.setOffline(false)
  await expect(page.getByTestId('contador-cola')).toHaveCount(0)
  await expect.poll(() => [api.gastos.length, api.pesos.length]).toEqual([1, 1])
  const [g] = api.gastos, [w] = api.pesos
  for (const c of [g, w]) { expect(c.clave.length).toBeGreaterThanOrEqual(16); expect(c.clave.length).toBeLessThanOrEqual(128); expect(c.cuerpo.captured_at).toMatch(ISO_CON_DESPLAZAMIENTO) }
  expect(g.cuerpo).toMatchObject({ amount_cents: 18550, category_id: 1, payment_method: 'cash', kind: 'variable' })
  expect(w.cuerpo).toMatchObject({ weight_g: 72500, person_id: 1 })
  console.log(`U5-13 clave gasto=${g.clave} captured_at=${g.cuerpo.captured_at}; clave peso=${w.clave} captured_at=${w.cuerpo.captured_at}`)
  // Reintento forzado de la misma captura (misma clave, mismo cuerpo): no duplica.
  const estado = await page.evaluate(async ({ clave, cuerpo }) => (await fetch('/api/v1/expenses', { method: 'POST', headers: { 'Content-Type': 'application/json', 'Idempotency-Key': clave }, body: JSON.stringify(cuerpo) })).status, g)
  expect(estado).toBe(200)
  expect(api.gastos.length).toBe(1)
  console.log(`U5-13 registros: gastos ${api.gastos.length}, pesos ${api.pesos.length} (antes 0 y 0); reintento forzado → ${estado}`)
})
test('si un envío falla, la captura sigue en la cola y se avisa', async ({ page, context }) => {
  const api = new ApiFalsa()
  await api.instalar(context)
  await page.goto('/inicio')
  await expect.poll(() => api.peticiones.some((p) => p.url.startsWith('/api/v1/categories'))).toBe(true)
  await context.setOffline(true)
  await page.getByRole('button', { name: 'Registrar gasto' }).last().tap()
  await page.keyboard.type('99')
  await page.getByRole('button', { name: 'Comida fuera' }).tap()
  await page.getByRole('button', { name: 'Guardar' }).tap()
  await expect(page.getByText('1 en cola')).toBeVisible()
  api.fallarSiguienteGasto = true
  await context.setOffline(false)
  await expect(page.getByTestId('franja-conexion')).toContainText('No se pudo enviar')
  await expect(page.getByText('1 en cola')).toBeVisible()
  await expect(page.getByRole('alert').filter({ hasText: 'No se pudo enviar · Revisar' })).toBeVisible({ timeout: 8000 }) // va después del aviso de éxito (uno a la vez, §7.15)
  expect(api.gastos.length).toBe(0)
})

async function capturarSinConexion(page: import('@playwright/test').Page, api: ApiFalsa, context: import('@playwright/test').BrowserContext) {
  await page.goto('/inicio')
  await expect.poll(() => api.peticiones.some((p) => p.url.startsWith('/api/v1/categories'))).toBe(true)
  await context.setOffline(true)
  await page.getByRole('button', { name: 'Registrar gasto' }).last().tap()
  await page.keyboard.type('99')
  await page.getByRole('button', { name: 'Comida fuera' }).tap()
  await page.getByRole('button', { name: 'Guardar' }).tap()
  await expect(page.getByText('1 en cola')).toBeVisible()
}
test('falla del servidor (500) al enviar la cola: sigue en la cola, la franja la cuenta y se avisa', async ({ page, context }) => {
  const api = new ApiFalsa()
  await api.instalar(context)
  await capturarSinConexion(page, api, context)
  api.fallaServidorGastos = true
  await context.setOffline(false)
  const franja = page.getByTestId('franja-conexion')
  await expect.poll(() => api.peticiones.filter((p) => p.metodo === 'POST' && p.url === '/api/v1/expenses').length).toBeGreaterThanOrEqual(1)
  await expect(franja).toContainText('No se pudo enviar')
  await expect(franja.getByTestId('contador-cola')).toHaveText('1 en cola')
  await expect(page.getByRole('alert').filter({ hasText: 'No se pudo enviar · Revisar' })).toBeVisible({ timeout: 8000 })
  expect(api.gastos.length).toBe(0)
  console.log(`U5-13 paso 7 (500): franja visible, «1 en cola», aviso «No se pudo enviar · Revisar»; registros ${api.gastos.length}`)
  // Sigue guardada en el teléfono: tras recargar (otro intento con 500) la franja la sigue contando.
  await page.reload()
  await expect(page.getByTestId('contador-cola')).toHaveText('1 en cola')
  await expect(page.getByTestId('franja-conexion')).toContainText('No se pudo enviar')
  // El servidor se recupera: «Revisar» → «Reintentar envío» la manda una sola vez y la franja desaparece.
  api.fallaServidorGastos = false
  await page.getByTestId('franja-conexion').getByRole('button', { name: 'Revisar' }).tap()
  await page.getByRole('button', { name: 'Reintentar envío' }).tap()
  await expect.poll(() => api.gastos.length).toBe(1)
  await expect(page.getByTestId('franja-conexion')).toHaveCount(0)
  expect(api.gastos[0].cuerpo).toMatchObject({ amount_cents: 9900, category_id: 2 })
})
test('sesión vencida (401) al enviar la cola: la captura no se pierde, se avisa y se envía al volver a entrar', async ({ page, context }) => {
  const api = new ApiFalsa()
  await api.instalar(context)
  await capturarSinConexion(page, api, context)
  api.conSesion = false // la sesión venció mientras el teléfono estaba sin red
  await context.setOffline(false)
  await expect(page.getByRole('alert').filter({ hasText: 'No se pudo enviar · Revisar' })).toBeVisible({ timeout: 8000 })
  await expect(page.getByLabel('Usuario')).toBeVisible() // 401 lleva a P-00
  expect(api.gastos.length).toBe(0)
  await page.getByLabel('Usuario').fill('axel')
  await page.getByLabel('Contraseña').fill('contrasena-de-prueba') // credencial falsa de la API falsa
  await page.getByRole('button', { name: 'Entrar' }).tap()
  const franja = page.getByTestId('franja-conexion')
  await expect(franja).toContainText('No se pudo enviar')
  await expect(franja.getByTestId('contador-cola')).toHaveText('1 en cola')
  console.log('U5-13 paso 7 (401): aviso «No se pudo enviar · Revisar», P-00, y al volver a entrar la franja cuenta «1 en cola»')
  await franja.getByRole('button', { name: 'Revisar' }).tap()
  await page.getByRole('button', { name: 'Reintentar envío' }).tap()
  await expect.poll(() => api.gastos.length).toBe(1)
  await expect(page.getByTestId('franja-conexion')).toHaveCount(0)
})
