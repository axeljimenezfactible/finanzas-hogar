// P-26.02 · «Parámetros del tablero» (agente frontend-cierre, WF-3a cierre). Decisión 4 del orquestador = opción A del bloqueo
// docs/estado/bloqueos/C3-ajustes-tablero.md (PROVISIONAL hasta la puerta 4): las cinco claves de K54 en Ajustes › Presupuesto con
// getDashboardSettings / updateDashboardSettings (N-275, N-278), vista previa «Antes → Después» (U4 §0.2) y «Restablecer» (U4 §0.5).
// Corre en 412×915 táctil (movil) y 1440×900 (escritorio). Las respuestas falsas se validan contra C1 (erroresC1, H-11).
import { expect, test, type BrowserContext, type Page } from '@playwright/test'
import { ApiFalsa } from '../apoyo/api-falsa'
import { controlesChicos, textosEnIngles } from '../apoyo/medidas'
import { erroresC1 } from './validar-c1'
import { unidadesDelServidor } from './unidades-servidor'

// `unit` = lo que manda el servidor real (tercera columna de @param de 00005_tablero.sql, una descripción larga), no una unidad
// inventada (M-1 de REVISION-frontend-cierre-r1). P-26.02 pinta la unidad corta de textos.ts: PROVISIONAL, bloqueo
// docs/estado/bloqueos/C1-unidad-parametros-tablero.md.
const UNIT = unidadesDelServidor()
type Valores = Record<string, [valor: number, inicial: number, unidad: string]>
const INICIALES: Valores = {
  'debt.capacity_warning_pct': [30, 30, UNIT['debt.capacity_warning_pct']], 'debt.capacity_limit_pct': [40, 40, UNIT['debt.capacity_limit_pct']],
  'dashboard.payment_reminder_count': [5, 5, UNIT['dashboard.payment_reminder_count']],
  'budget.signal_months': [4, 3, UNIT['budget.signal_months']], 'budget.signal_margin_pct': [10, 10, UNIT['budget.signal_margin_pct']],
}
const lectura = (v: Valores) => ({ parameters: Object.fromEntries(Object.entries(v).map(([key, [value, initial_value, unit]]) =>
  [key, { key, value_type: 'integer', unit, value, initial_value, is_initial: value === initial_value }])) })

async function preparar(context: BrowserContext) {
  await new ApiFalsa().instalar(context)
  const valores: Valores = JSON.parse(JSON.stringify(INICIALES))
  const escrituras: { metodo: string; cuerpo: any }[] = [], errores: string[] = []
  await context.route(/\/api\/v1\/settings\/(dashboard|projection)$/, async (route) => {
    const pet = route.request(), ruta = new URL(pet.url()).pathname
    const cuerpo = pet.postData() ? pet.postDataJSON() : null
    const responder = (estado: number, json: unknown) => { errores.push(...erroresC1(pet.method(), ruta, estado, cuerpo, json)); return route.fulfill({ status: estado, json }) }
    if (ruta.endsWith('/projection')) return pet.method() === 'GET' ? responder(200, { repeat_cadence_enabled: false }) : route.fallback()
    if (pet.method() === 'GET') return responder(200, lectura(valores))
    escrituras.push({ metodo: pet.method(), cuerpo })
    for (const [k, n] of Object.entries(cuerpo?.parameters ?? {})) valores[k][0] = n as number
    for (const k of cuerpo?.reset_parameters ?? []) valores[k][0] = valores[k][1]
    return responder(200, lectura(valores))
  })
  return { escrituras, errores }
}
const bloque = (page: Page) => page.getByRole('region', { name: 'Parámetros del tablero' })
const editar = async (page: Page, etiqueta: string) => {
  await bloque(page).getByRole('button', { name: `Editar ${etiqueta}` }).click()
  return bloque(page).getByRole('textbox', { name: etiqueta })
}

test('P-26.02 «Parámetros del tablero»: las 5 claves de K54 con su valor inicial y «Restablecer» solo en la que cambió', async ({ page, context }) => {
  const { errores } = await preparar(context)
  await page.goto('/ajustes/02')
  const b = bloque(page)
  await expect(b).toBeVisible()
  for (const k of Object.keys(INICIALES)) await expect(b.getByTestId(`parametro-${k}`)).toBeVisible()
  // El dato falso trae la descripción larga del servidor (sin unidades inventadas) y la fila no la repite: «valor unidad-corta».
  for (const [k, [, , unit]] of Object.entries(INICIALES)) {
    expect(unit.length, `unit de ${k} sale de @param de 00005`).toBeGreaterThan(20)
    await expect(b).not.toContainText(unit)
  }
  // (La descripción larga trae «de 1 a 100», que contradice V-23; el único «de 1 a 100» que queda es el de la nota de rangos.)
  await expect(b.getByText(/de 1 a 100/)).toHaveCount(1)
  await expect(b.getByText(/de 1 a 100/)).toContainText('Rangos:')
  await expect(b).toContainText('Valor inicial: 30 %')
  await expect(b).toContainText('Valor inicial: 40 %')
  await expect(b).toContainText('Valor inicial: 5 pagos')
  await expect(b).toContainText('Valor inicial: 3 meses')
  await expect(b).toContainText('Valor inicial: 10 puntos')
  await expect(b.getByTestId('parametro-budget.signal_months')).toContainText('4 meses')
  await expect(b).toContainText('PROVISIONAL (bloqueo C1-tablero)')
  await expect(b.getByRole('button', { name: 'Restablecer al valor inicial' })).toHaveCount(1) // solo budget.signal_months (4, inicial 3)
  expect(await textosEnIngles(page)).toEqual([])
  if (test.info().project.name === 'movil') expect((await controlesChicos(page)).chicos).toEqual([]) // 44 × 44 solo con puntero grueso (U5-06)
  expect(errores).toEqual([])
})

test('P-26.02 tope del barómetro: vista previa de tramos y escala; Esc no guarda y Enter manda updateDashboardSettings', async ({ page, context }) => {
  const { escrituras, errores } = await preparar(context)
  await page.goto('/ajustes/02')
  const campo = await editar(page, 'Barómetro: tope antes del rojo (% del ingreso del mes)')
  await campo.fill('45')
  const previa = bloque(page).getByTestId('vista-previa')
  await expect(previa).toContainText('Verde hasta 30 % · amarillo hasta 40 % · rojo arriba de 40 %')
  await expect(previa).toContainText('Verde hasta 30 % · amarillo hasta 45 % · rojo arriba de 45 %')
  await expect(previa).toContainText('De 0 a 60 %')
  await expect(previa).toContainText('De 0 a 70 %') // múltiplo de 10 ≥ 1.5 × 45 = 67.5 (U3 §B-1)
  await campo.press('Escape')
  await expect(previa).toHaveCount(0)
  expect(escrituras).toEqual([])
  const otra = await editar(page, 'Barómetro: tope antes del rojo (% del ingreso del mes)')
  await otra.fill('45')
  await otra.press('Enter')
  await expect(page.getByText('Guardado').first()).toBeVisible()
  expect(escrituras).toEqual([{ metodo: 'PATCH', cuerpo: { parameters: { 'debt.capacity_limit_pct': 45 } } }])
  await expect(bloque(page).getByTestId('parametro-debt.capacity_limit_pct')).toContainText('45 %')
  expect(errores).toEqual([])
})

test('P-26.02 fin del verde igual al tope no se guarda (V-23 con el rango que deja el tope); los demás muestran su vista previa', async ({ page, context }) => {
  const { escrituras } = await preparar(context)
  await page.goto('/ajustes/02')
  const verde = await editar(page, 'Barómetro: fin del tramo verde (% del ingreso del mes)')
  await verde.fill('40')
  await verde.press('Enter')
  await expect(bloque(page)).toContainText('Escribe un número entero del 1 al 39.')
  expect(escrituras).toEqual([])
  await verde.press('Escape')
  const pagos = await editar(page, 'Pagos próximos que muestra Inicio')
  await pagos.fill('8')
  await expect(bloque(page).getByTestId('vista-previa')).toContainText('Hasta 5 pagos')
  await expect(bloque(page).getByTestId('vista-previa')).toContainText('Hasta 8 pagos')
  await pagos.press('Escape')
  const margen = await editar(page, 'Sugerencia de presupuesto: margen sobre el plan (puntos de %)')
  await margen.fill('15')
  await expect(bloque(page).getByTestId('vista-previa')).toContainText('Subir si pasa de 110 % del plan · recortar si queda abajo de 90 %')
  await expect(bloque(page).getByTestId('vista-previa')).toContainText('Subir si pasa de 115 % del plan · recortar si queda abajo de 85 %')
  await margen.press('Escape')
  expect(escrituras).toEqual([])
})

test('P-26.02 «Restablecer al valor inicial» manda reset_parameters con la clave', async ({ page, context }) => {
  const { escrituras, errores } = await preparar(context)
  await page.goto('/ajustes/02')
  await bloque(page).getByTestId('parametro-budget.signal_months').getByRole('button', { name: 'Restablecer al valor inicial' }).click()
  await expect(page.getByText('Valor restablecido')).toBeVisible()
  expect(escrituras).toEqual([{ metodo: 'PATCH', cuerpo: { reset_parameters: ['budget.signal_months'] } }])
  await expect(bloque(page).getByRole('button', { name: 'Restablecer al valor inicial' })).toHaveCount(0)
  expect(errores).toEqual([])
})
