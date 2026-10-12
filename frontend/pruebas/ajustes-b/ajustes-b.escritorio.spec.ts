// Ajustes 8 a 14 (U4 parte b): una prueba por pantalla en 412×915 (móvil) y 1440×900 (escritorio), con API falsa.
import { expect, test, type Page } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { erroresTrafico, instalarAjustesB, type Escritura } from './datos'
import { moverTarjetaA } from '../ajustes-a/tablero'

let escrituras: Escritura[] = []
test.beforeEach(async ({ context }) => {
  escrituras = []
  await new ApiFalsa().instalar(context)
  await instalarAjustesB(context, escrituras)
})
test.afterEach(() => expect(erroresTrafico).toEqual([]))
const titulo = (page: Page, texto: string) => expect(page.getByRole('heading', { level: 2, name: texto })).toBeVisible()

test('P-26.08 Reportes y frecuencias: tres reportes, próximo envío y bajo demanda', async ({ page }) => {
  await page.goto('/ajustes/08')
  await titulo(page, 'Reportes y frecuencias')
  for (const r of ['Resumen diario', 'Resumen semanal', 'Reporte mensual en PDF']) await expect(page.getByText(r, { exact: true }).first()).toBeVisible()
  await expect(page.getByText('Lun 12 oct 2026, 21:30').first()).toBeVisible()
  await expect(page.getByText('No se enviará mientras esté apagado').first()).toBeVisible()
  await expect(page.getByText('Sin destinatarios: no llegará a nadie').first()).toBeVisible()
  await page.getByRole('button', { name: 'Generar reporte mensual ahora' }).click()
  await expect(page.getByRole('dialog')).toContainText('Enviarlo también por Telegram')
  await expect(page.getByRole('checkbox')).not.toBeChecked()
})

test('P-26.09 Notificaciones: alertas, chats confirmados y chat detectado con descarte crítico', async ({ page }) => {
  await page.goto('/ajustes/09')
  await titulo(page, 'Notificaciones')
  await expect(page.getByText('Pago próximo', { exact: true })).toBeVisible()
  await expect(page.getByText('Ximena aún no confirma su chat')).toBeVisible()
  await expect(page.getByText('Xime', { exact: true })).toBeVisible()
  await expect(page.getByText(/•••4821/)).toBeVisible()
  await page.getByRole('button', { name: 'Descartar', exact: true }).click()
  const dialogo = page.getByRole('alertdialog')
  await expect(dialogo).toContainText('Descartar el chat de “Xime” (•••4821)')
  await expect(dialogo.getByRole('button', { name: 'Cancelar' })).toBeFocused()
  await dialogo.getByRole('button', { name: 'Descartar chat' }).click()
  await expect.poll(() => escrituras.some((e) => e.metodo === 'DELETE' && e.ruta === '/api/v1/telegram-chat-candidates/5')).toBe(true)
})

test('P-26.09 umbral por categoría: al elegir una categoría de ejemplo muestra sus montos (H-08)', async ({ page }) => {
  await page.goto('/ajustes/09')
  await titulo(page, 'Notificaciones')
  const alerta = page.getByTestId('alerta-budget_threshold')
  await expect(alerta).toBeVisible()
  const editar = alerta.getByRole('button', { name: 'Editar', exact: true })
  if (await editar.count()) await editar.click()
  const selector = page.getByLabel('Categoría de ejemplo')
  await expect(selector.locator('option')).toHaveText(['Elige una categoría', 'Despensa'])
  await selector.selectOption({ label: 'Despensa' })
  await expect(page.getByTestId('umbral-categoria')).toHaveText('Con el presupuesto de Despensa de este mes ($8,000.00), el primer aviso llega al gastar $6,400.00 y el de tope al gastar $8,000.00')
})

test('P-26.10 Apariencia: el tema cambia al instante y se guarda en el dispositivo', async ({ page }) => {
  await page.goto('/ajustes/10')
  await titulo(page, 'Apariencia')
  await page.getByRole('radio', { name: 'Oscuro' }).click()
  await expect(page.locator('html')).toHaveAttribute('data-theme', 'dark')
  await expect(page.getByText('Guardado')).toBeVisible()
  await expect(page.getByRole('switch', { name: 'Atajos de una tecla' })).toBeChecked()
})

test('P-26.11 Datos y respaldos: respaldos, parámetro con vista previa y activación crítica de la IA', async ({ page }) => {
  await page.goto('/ajustes/11')
  await titulo(page, 'Datos y respaldos')
  await expect(page.getByRole('button', { name: 'Respaldar ahora' }).first()).toBeVisible()
  await expect(page.getByText(/Último respaldo correcto: 12-oct-2026 03:00/)).toBeVisible()
  await expect(page.getByText('Falló').first()).toBeVisible()
  await expect(page.getByText('CSV de gastos fijos e ingresos')).toBeVisible()
  const campo = page.getByLabel('Tamaño máximo de un adjunto de gasto')
  await campo.fill('12')
  await expect(page.getByText('Con 12 MB, un ticket o adjunto más pesado se rechaza')).toBeVisible()
  await campo.press('Enter')
  await expect.poll(() => escrituras.find((e) => e.ruta === '/api/v1/settings/data')?.cuerpo).toEqual({ parameters: { 'attachment.max_size_mb': 12 } })
  await expect(page.getByRole('button', { name: 'Restablecer al valor inicial' })).toBeVisible()
  await page.getByRole('switch', { name: 'Activa · Modelo local (Ollama)' }).click()
  const dialogo = page.getByRole('alertdialog')
  await expect(dialogo).toContainText('Siempre pedirá confirmación antes de guardar un ticket')
  await dialogo.getByRole('button', { name: 'Activar modelo local' }).click()
  await expect.poll(() => escrituras.find((e) => e.ruta === '/api/v1/ai/providers/local_model')?.cuerpo).toEqual({ is_enabled: true })
})

test('P-26.12 Personas y reparto: papeles, reglas y regla V-16', async ({ page }) => {
  await page.goto('/ajustes/12')
  await titulo(page, 'Personas y reparto')
  await expect(page.getByText('Participante externa').first()).toBeVisible()
  await expect(page.getByText('Conjunto 50/50: Karen 50.00 % · el resto 50.00 % (Axel, hogar)')).toBeVisible()
  await page.locator('[data-prueba="regla-3"]').getByRole('button', { name: 'Editar partes' }).click()
  await page.getByLabel('Ximena %').fill('60')
  await expect(page.getByRole('alert')).toContainText('Las partes con porcentaje suman 110.00 %. No pueden pasar de 100 %: lo que falta lo recibe Axel.')
})

test('P-26.13 Salud y metas: perfil con meta sugerida, perfil faltante y metas', async ({ page }) => {
  await page.goto('/ajustes/13')
  await titulo(page, 'Salud y metas de ahorro')
  await expect(page.getByText('Basal 1,648.75 × 1.55 = 2,555.5625 → 2,556 kcal')).toBeVisible()
  await expect(page.getByText('Sin perfil de salud')).toBeVisible()
  await expect(page.getByText('Vacaciones').first()).toBeVisible()
  await page.getByLabel('Objetivo').first().selectOption('lose')
  await expect(page.getByRole('status', { name: 'Vista previa del efecto' }).first()).toBeVisible()
})

test('P-26.14 Inventario: modos, estado con texto y días restantes', async ({ page }) => {
  await page.goto('/ajustes/14')
  await titulo(page, 'Inventario')
  await expect(page.getByText('Aceite de cocina').first()).toBeVisible()
  await expect(page.getByText('Agotado').first()).toBeVisible()
  await expect(page.getByRole('button', { name: 'Agregar producto al inventario' })).toBeVisible()
})

test('P-26.14 tablero: soltar en el otro modo muestra el estado antes → después; Esc la regresa sin guardar y Enter aplica', async ({ page }) => {
  await page.goto('/ajustes/14')
  await page.getByRole('tab', { name: 'Tablero' }).click()
  const mover = () => moverTarjetaA(page, 'Arroz', 'Ritmo de consumo')
  await mover()
  const previa = page.getByTestId('previa-tablero')
  await expect(previa).toContainText('«Arroz» pasa a')
  await expect(previa).toContainText('Estado: Agotado')
  await expect(previa).toContainText('Falta el ritmo')
  await page.keyboard.press('Escape')
  await expect(previa).toHaveCount(0)
  expect(escrituras.filter((e) => e.metodo === 'PATCH')).toHaveLength(0)
  await mover()
  await expect(page.getByRole('button', { name: 'Aplicar' })).toBeFocused()
  await page.keyboard.press('Enter')
  await expect.poll(() => escrituras.find((e) => e.metodo === 'PATCH')).toMatchObject({ ruta: '/api/v1/inventory-items/2', cuerpo: { tracking_mode: 'consumption_rate' } })
})

test('R2-03: en 12, 13 y 14 toda celda editable tiene nombre visible, también vacía', async ({ page }) => {
  for (const [ruta, titulo] of [['/ajustes/12', 'Personas y reparto'], ['/ajustes/13', 'Salud y metas'], ['/ajustes/14', 'Inventario']]) {
    await page.goto(ruta)
    await expect(page.getByRole('heading', { level: 2, name: titulo })).toBeVisible()
    await page.waitForLoadState('networkidle')
    await expect(page.getByText('Sin valor').first()).toBeVisible()
    const sinNombre = await page.locator('button:visible').evaluateAll((bs) => bs.filter((b) => !(b.getAttribute('aria-label') || b.getAttribute('aria-labelledby') || b.textContent?.trim() || b.getAttribute('title'))).map((b) => b.outerHTML.slice(0, 120)))
    expect(sinNombre, ruta).toEqual([])
  }
})
