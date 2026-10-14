// Ajustes 8 a 14 (U4 parte b): una prueba por pantalla en 412×915 (móvil) y 1440×900 (escritorio), con API falsa.
import { expect, test, type Page } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { erroresTrafico, instalarAjustesB, RESPUESTAS, type Escritura } from './datos'
import { erroresC1 } from '../ajustes-a/validar-c1'
import type { Esquemas } from '../../src/api/tipos'
import { moverTarjetaA } from '../ajustes-a/tablero'
import { controlesChicos, textosEnIngles } from '../apoyo/medidas'

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

test('R4-02 · P-26.13 «Usar el de la casa» y «Usar la sugerida» con el servidor en 500 muestran «No se pudo guardar…» debajo del campo', async ({ page, context }) => {
  const fallosPagina: string[] = []
  page.on('pageerror', (e) => fallosPagina.push(e.message))
  const lectura = RESPUESTAS['GET /api/v1/health-profiles'] as { items: Esquemas['HealthProfile'][]; next_cursor: null }
  const conAjuste = { ...lectura, items: [{ ...lectura.items[0], goal_adjust_lose_kcal: 300 }] }
  const intentos: unknown[] = []
  await context.route(/\/api\/v1\/health-profiles/, async (route) => {
    const req = route.request(), u = new URL(req.url())
    if (req.method() === 'GET') { erroresTrafico.push(...erroresC1('GET', u.pathname, 200, null, conAjuste)); return route.fulfill({ status: 200, json: conAjuste }) }
    if (req.method() !== 'PATCH') return route.fallback()
    intentos.push(req.postDataJSON())
    // 500 internal_error: C1 lo declara en `default` con el esquema Error; validar-c1 solo lee códigos explícitos, así que este no pasa por erroresC1.
    return route.fulfill({ status: 500, json: { code: 'internal_error', message: 'Error interno del servidor.', request_id: 'prueba' } satisfies Esquemas['Error'] })
  })
  await page.goto('/ajustes/13')
  await titulo(page, 'Salud y metas de ahorro')
  const mensaje = 'No se pudo guardar. Revisa la conexión con el servidor de casa'
  const casa = page.locator('div', { has: page.getByLabel('Mi ajuste para bajar (kcal)') }).last()
  await casa.getByRole('button', { name: 'Usar el de la casa' }).click()
  await expect(casa.getByText(mensaje)).toBeVisible()
  await expect(casa.getByRole('button', { name: 'Reintentar' })).toBeVisible()
  const sugerida = page.locator('div', { has: page.getByLabel('Meta confirmada (kcal por día)') }).last()
  await sugerida.getByRole('button', { name: 'Usar la sugerida' }).click()
  await expect(sugerida.getByText(mensaje)).toBeVisible()
  await expect(page.getByText(mensaje)).toHaveCount(2)
  expect(intentos).toEqual([{ goal_adjust_lose_kcal: null }, { daily_kcal_target: 2556 }])
  expect(fallosPagina).toEqual([])
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
  await page.getByRole('button', { name: 'Tablero', exact: true }).click()
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

// R2-06: una prueba por operación que edita una clave `reports.*` o `alerts.*` y comprueba el cuerpo (respuesta falsa válida en C1).
test('R2-06 · P-26.08: editar la hora del resumen diario manda updateReportSettings solo con esa clave', async ({ page }) => {
  await page.goto('/ajustes/08')
  await titulo(page, 'Reportes y frecuencias')
  const fila = page.getByTestId('reporte-daily_summary')
  await expect(fila).toBeVisible() // ya cargó: en escritorio es la fila de la tabla y en móvil la tarjeta
  const editar = fila.getByRole('button', { name: 'Día, hora y destinatarios' })
  if (await editar.count()) await editar.click() // en móvil, el editor está en la hoja de la tarjeta
  const hora = page.locator('[id="param-reports.daily_send_time"]')
  await hora.fill('22:15')
  await hora.press('Enter')
  await expect.poll(() => escrituras.find((e) => e.ruta === '/api/v1/settings/reports')).toEqual({ metodo: 'PATCH', ruta: '/api/v1/settings/reports', cuerpo: { parameters: { 'reports.daily_send_time': '22:15' } } })
})

test('R2-06 · P-26.09: editar la hora de los avisos manda updateAlertSettings solo con esa clave', async ({ page }) => {
  await page.goto('/ajustes/09')
  await titulo(page, 'Notificaciones')
  const hora = page.getByLabel('Hora de los avisos con hora')
  await hora.fill('08:30')
  await hora.press('Enter')
  await expect.poll(() => escrituras.find((e) => e.ruta === '/api/v1/settings/alerts')).toEqual({ metodo: 'PATCH', ruta: '/api/v1/settings/alerts', cuerpo: { parameters: { 'alerts.send_time': '08:30' } } })
})

// H-08: `Ctrl+/` en 09, 13 y 14 (U4 §0.2 punto 2 y §0.7). En escritorio se abre con el atajo y se cierra con Esc o Enter;
// en táctil (U4 §0.7: en Android no hay atajos) el mismo bloque se abre con «Editar seleccionadas» y se usan sus botones.
const abrirGrupo = (page: Page, movil: boolean) => (movil ? page.getByRole('button', { name: /Editar seleccionadas/ }).click() : page.keyboard.press('Control+/'))
const cancelarGrupo = (page: Page, movil: boolean) => (movil ? page.getByRole('region', { name: /^Editar/ }).getByRole('button', { name: 'Cancelar' }).click() : page.keyboard.press('Escape'))
const aplicarGrupo = (page: Page, movil: boolean, n: number) => (movil ? page.getByRole('button', { name: `Aplicar a ${n}` }).click() : page.keyboard.press('Enter'))
/** El foco regresa a donde estaba: la última casilla marcada (atajo) o «Editar seleccionadas» (táctil). */
const focoRegresa = (page: Page, movil: boolean, casilla: string) => expect(movil ? page.getByRole('button', { name: /Editar seleccionadas/ }) : page.getByRole('checkbox', { name: casilla })).toBeFocused()
const parches = () => escrituras.filter((e) => e.metodo === 'PATCH').map((e) => [e.ruta, e.cuerpo])

test('H-08 · P-26.09 Ctrl+/: activar o desactivar varias alertas con vista previa sumada; Esc no guarda y Enter aplica a todas', async ({ page, isMobile }) => {
  await page.goto('/ajustes/09')
  await titulo(page, 'Notificaciones')
  await page.getByRole('checkbox', { name: 'Seleccionar Pago próximo' }).check()
  await page.getByRole('checkbox', { name: 'Seleccionar Precios viejos antes del día de despensa' }).check()
  await expect(page.getByText('2 seleccionadas')).toBeVisible()
  if (isMobile) { // casillas, «Seleccionar todas» y la barra de selección con área de toque ≥ 44 (U5-06)
    const c = await controlesChicos(page)
    expect(c.chicos.filter((x) => /[Ss]elecci/.test(x.control))).toEqual([])
  }
  await abrirGrupo(page, isMobile)
  const bloque = page.getByRole('region', { name: 'Editar 2 registros a la vez' })
  const valor = bloque.getByRole('combobox', { name: 'Valor nuevo' })
  await expect(valor).toBeFocused()
  await valor.selectOption('no')
  const previa = bloque.getByTestId('vista-previa')
  await expect(previa).toContainText('Alertas activas')
  await expect(previa).toContainText('4 de 4')
  await expect(previa).toContainText('2 de 4')
  await expect(previa).toContainText('Pago próximo')
  await expect(previa).toContainText('Inactiva')
  console.log(`H-08 09 vista previa: «${(await previa.innerText()).replace(/\s+/g, ' ')}»`)
  expect(await textosEnIngles(page)).toEqual([]) // U5-07 con el bloque abierto
  await cancelarGrupo(page, isMobile)
  await expect(bloque).toHaveCount(0)
  expect(parches()).toEqual([]) // Esc cancela sin cambiar ninguna
  await expect(page.getByRole('checkbox', { name: 'Seleccionar Pago próximo' })).toBeChecked() // la selección sigue
  await focoRegresa(page, isMobile, 'Seleccionar Precios viejos antes del día de despensa')
  await abrirGrupo(page, isMobile)
  await valor.selectOption('no')
  await expect(previa).toContainText('2 de 4')
  await aplicarGrupo(page, isMobile, 2)
  await expect.poll(parches).toEqual([['/api/v1/notifications/4', { is_enabled: false }], ['/api/v1/notifications/6', { is_enabled: false }]])
  await expect(page.getByText('Guardado en 2 registros')).toBeVisible()
  await expect(page.getByRole('checkbox', { name: 'Seleccionar Pago próximo' })).not.toBeChecked()
})

test('H-08 · P-26.13 Ctrl+/: «Activa» de varias metas con el total reservado sumado; Esc no guarda y Enter aplica a todas', async ({ page, isMobile }) => {
  await page.goto('/ajustes/13')
  await titulo(page, 'Salud y metas de ahorro')
  await page.getByRole('checkbox', { name: 'Seleccionar Vacaciones' }).check()
  await page.getByRole('checkbox', { name: 'Seleccionar Fondo de ahorro' }).check()
  await abrirGrupo(page, isMobile)
  const bloque = page.getByRole('region', { name: 'Editar 2 registros a la vez' })
  const valor = bloque.getByRole('combobox', { name: 'Valor nuevo' })
  await valor.selectOption('no')
  const previa = bloque.getByTestId('vista-previa')
  await expect(previa).toContainText('Total reservado para metas en')
  await expect(previa).toContainText('$1,500.00') // 1,000.00 de Vacaciones + 500.00 del Fondo, activas
  await expect(previa).toContainText('$0.00') //     las dos pasan a inactivas
  console.log(`H-08 13 vista previa: «${(await previa.innerText()).replace(/\s+/g, ' ')}»`)
  expect(await textosEnIngles(page)).toEqual([]) // U5-07 con el bloque abierto
  await cancelarGrupo(page, isMobile)
  await expect(bloque).toHaveCount(0)
  expect(parches()).toEqual([])
  await focoRegresa(page, isMobile, 'Seleccionar Fondo de ahorro')
  await abrirGrupo(page, isMobile)
  await valor.selectOption('no')
  await expect(previa).toContainText('$0.00')
  await aplicarGrupo(page, isMobile, 2)
  await expect.poll(parches).toEqual([['/api/v1/savings-goals/1', { is_active: false }], ['/api/v1/savings-goals/2', { is_active: false }]])
  await expect(page.getByText('Guardado en 2 registros')).toBeVisible()
})

test('H-08 · P-26.14 Ctrl+/: nivel de varias filas con «{k} pasan a Reponer…» y las filas que cambian; Esc no guarda y Enter aplica a todas', async ({ page, isMobile }) => {
  await page.goto('/ajustes/14')
  await titulo(page, 'Inventario')
  await page.getByRole('checkbox', { name: 'Seleccionar Arroz' }).check()
  await page.getByRole('checkbox', { name: 'Seleccionar Frijol' }).check()
  await abrirGrupo(page, isMobile)
  const bloque = page.getByRole('region', { name: 'Editar 2 registros a la vez' })
  await bloque.getByRole('combobox', { name: 'Propiedad' }).selectOption({ label: 'Nivel de reposición' })
  const valor = bloque.getByRole('textbox', { name: 'Valor nuevo' })
  await valor.fill('2000') // 2,000 g: Frijol (1,500 g) pasa a Reponer; Arroz sigue Agotado
  const previa = bloque.getByTestId('vista-previa')
  await expect(previa).toContainText('1 pasa a Reponer, 0 a Agotado y 0 a Suficiente')
  await expect(previa).toContainText('Suficiente 1 · Reponer 0 · Agotado 1')
  await expect(previa).toContainText('Suficiente 0 · Reponer 1 · Agotado 1')
  await expect(previa).toContainText('Frijol')
  await expect(previa).not.toContainText('Arroz') // solo las filas que cambian
  console.log(`H-08 14 vista previa: «${(await previa.innerText()).replace(/\s+/g, ' ')}»`)
  expect(await textosEnIngles(page)).toEqual([]) // U5-07 con el bloque abierto
  await cancelarGrupo(page, isMobile)
  await expect(bloque).toHaveCount(0)
  expect(parches()).toEqual([])
  await focoRegresa(page, isMobile, 'Seleccionar Frijol')
  await abrirGrupo(page, isMobile)
  await bloque.getByRole('combobox', { name: 'Propiedad' }).selectOption({ label: 'Modo' })
  await bloque.getByRole('combobox', { name: 'Valor nuevo' }).selectOption('consumption_rate')
  await expect(previa).toContainText('Falta el ritmo') // sin ritmo capturado (U4 P-26.14)
  await bloque.getByRole('combobox', { name: 'Propiedad' }).selectOption({ label: 'Nivel de reposición' })
  await valor.fill('2000')
  await expect(previa).toContainText('1 pasa a Reponer')
  if (!isMobile) await valor.press('Enter'); else await aplicarGrupo(page, true, 2)
  await expect.poll(parches).toEqual([['/api/v1/inventory-items/2', { reorder_level_milli: 2000000 }], ['/api/v1/inventory-items/3', { reorder_level_milli: 2000000 }]])
  await expect(page.getByText('Guardado en 2 registros')).toBeVisible()
})
