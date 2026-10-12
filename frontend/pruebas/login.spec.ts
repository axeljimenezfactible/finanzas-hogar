// Prueba base: P-00 Iniciar sesión (K11), bloqueo 429 y U5-08 (sin peticiones a hosts externos).
import { expect, test } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'
import { ApiFalsa } from './apoyo/api-falsa'
test('P-00: entra con usuario y contraseña y vuelve a Inicio; sin hosts externos', async ({ page, context }) => {
  const api = new ApiFalsa(); api.conSesion = false
  await api.instalar(context)
  const hosts = new Set<string>()
  page.on('request', (r) => hosts.add(new URL(r.url()).hostname))
  await page.goto('/inicio')
  await expect(page.getByRole('heading', { name: 'Iniciar sesión' })).toBeVisible()
  const axe = await new AxeBuilder({ page: page as never }).withTags(['wcag2a', 'wcag2aa', 'wcag21aa', 'wcag22aa']).analyze()
  expect(axe.violations.map((v) => `${v.id}: ${v.help}`)).toEqual([])
  await page.getByLabel('Usuario').fill('axel')
  await page.getByLabel('Contraseña').fill('equivocada')
  await page.getByRole('button', { name: 'Entrar' }).click()
  await expect(page.getByRole('alert')).toContainText('Usuario o contraseña incorrectos.')
  await page.getByLabel('Contraseña').fill('contrasena-de-prueba')
  await page.getByRole('button', { name: 'Entrar' }).click()
  await expect(page).toHaveURL(/\/inicio$/)
  await expect(page.getByRole('heading', { level: 1, name: 'Inicio' })).toBeVisible()
  const externos = [...hosts].filter((h) => !['localhost', '127.0.0.1', 'finanzas.local'].includes(h))
  console.log(`U5-08 hosts pedidos: ${[...hosts].join(', ')}; externos: ${externos.length}`)
  expect(externos).toEqual([])
})
test('P-00: el bloqueo por intentos dice la hora en español', async ({ page, context }) => {
  const api = new ApiFalsa(); api.conSesion = false; api.limitarLogin = true
  await api.instalar(context)
  await page.goto('/entrar')
  await page.getByLabel('Usuario').fill('axel')
  await page.getByLabel('Contraseña').fill('cualquiera-de-prueba')
  await page.getByRole('button', { name: 'Entrar' }).click()
  await expect(page.getByRole('alert')).toContainText('Demasiados intentos. Puedes volver a intentar a las 21:45.')
})
