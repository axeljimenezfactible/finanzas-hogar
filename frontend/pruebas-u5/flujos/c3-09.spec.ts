// C3-09 · Cifras de pantalla = API = oráculo (WF-3b u5-flujos). SOLO LECTURA: se corre contra el servidor de desarrollo
// con la semilla real (https://finanzas.local:8443 por omisión). Lee en pantalla (innerText de tablas y tarjetas) y compara
// con oraculo/model.json, al centavo. Nada se corrige: cada diferencia queda en el informe con pantalla y cifra.
import { mkdirSync, readFileSync, writeFileSync } from 'node:fs'
import path from 'node:path'
import { expect, test, type Page } from '@playwright/test'
import { entrar, RAIZ } from './apoyo'

const modelo = JSON.parse(readFileSync(path.join(RAIZ, 'oraculo/model.json'), 'utf8'))
const centavos = (x: number) => Math.round(x * 100)
const MES = ['ene', 'feb', 'mar', 'abr', 'may', 'jun', 'jul', 'ago', 'sep', 'oct', 'nov', 'dic']
const diaCorto = (iso: string) => { const [a, m, d] = iso.split('-').map(Number); return `${d}-${MES[m - 1]}-${a}` }
const mesCorto = (iso: string) => { const [a, m] = iso.split('-').map(Number); return `${MES[m - 1]}-${a}` }
/** «$1,234.56» → 123456 centavos (NaN si no es dinero). */
const aCentavos = (s: string) => { const m = /^-?\$?(-?)([0-9][0-9,]*\.[0-9]{2})$/.exec(s.trim()); return m ? (s.trim().startsWith('-') || m[1] === '-' ? -1 : 1) * Math.round(parseFloat(m[2].replace(/,/g, '')) * 100) : NaN }

interface Fila { n: number; pantalla: string; cifra: string; mostrado: string; esperado: string; dif: number }
const filas: Fila[] = []
function comparar(pantalla: string, cifra: string, mostradoC: number, esperadoC: number) {
  const f = { n: filas.length + 1, pantalla, cifra, mostrado: Number.isNaN(mostradoC) ? 'NO ENCONTRADA' : (mostradoC / 100).toFixed(2), esperado: (esperadoC / 100).toFixed(2), dif: Number.isNaN(mostradoC) ? NaN : mostradoC - esperadoC }
  filas.push(f)
}
async function celdas(page: Page, tabla: ReturnType<Page['locator']>) {
  return tabla.locator('tr').evaluateAll((trs) => trs.map((tr) => [...tr.querySelectorAll('th,td')].map((c) => (c as HTMLElement).innerText.replace(/\s+/g, ' ').trim())))
}
const primerDinero = (t: string, despuesDe: RegExp) => { const m = despuesDe.exec(t); return m ? aCentavos(m[1]) : NaN }

test('C3-09 · cifras de Inicio, créditos y flujo = oráculo (model.json)', async ({ page }) => {
  const ev = modelo.credito.eventos as { fecha: string; tipo: string; monto: number; disponible_A: number; disponible_B: number }[]
  await entrar(page)
  await expect(page.getByTestId('tarjeta-3')).toHaveJSProperty('tagName', 'A')

  // ── Inicio · Resumen · tarjeta «Próximo pago» (tarjeta 3)
  const t3 = (await page.getByTestId('tarjeta-3').innerText()).replace(/\s+/g, ' ')
  comparar('Inicio · tarjeta Próximo pago', 'Próximo pago (30-oct-2026)', primerDinero(t3, /Próximo pago (\$[0-9,.]+)/), centavos(ev[1].monto))
  expect(t3).toContain(diaCorto(modelo.despensas[0].pagos[0].fecha).replace('-', ' ').replace('-', ' ').replace('oct', 'oct')) // «30 oct 2026»
  comparar('Inicio · tarjeta Próximo pago', 'Disponible hoy, lectura B', primerDinero(t3, /Disponible hoy (\$[0-9,.]+)/), centavos(ev[0].disponible_B))
  comparar('Inicio · tarjeta Próximo pago', 'Parte del hogar del pago 1', primerDinero(t3, /Hogar (\$[0-9,.]+)/), centavos(modelo.despensas[0].pagos[0].axel))
  comparar('Inicio · tarjeta Próximo pago', 'Parte de Karen del pago 1', primerDinero(t3, /Karen (\$[0-9,.]+)/), centavos(modelo.despensas[0].pagos[0].karen))

  // ── Inicio · Detalle · «Disponible por línea» y mínimos
  await page.getByRole('tab', { name: 'Detalle' }).click()
  const disp = page.getByText('Disponible por línea').locator('xpath=ancestor::section[1]')
  await expect(page.getByText('Lectura A').first()).toBeVisible()
  const det = (await page.locator('main').innerText()).replace(/\s+/g, ' ')
  comparar('Inicio · Detalle · Disponible por línea', 'Disponible inicial, lectura A', primerDinero(det, /Lectura A (\$[0-9,.]+)/), centavos(ev[0].disponible_A))
  comparar('Inicio · Detalle · Disponible por línea', 'Disponible inicial, lectura B', primerDinero(det, /Lectura B (\$[0-9,.]+)/), centavos(ev[0].disponible_B))
  comparar('Inicio · Detalle · Disponible por línea', 'Mínimo del horizonte, lectura A', primerDinero(det, /Lectura A \$[0-9,.]+ Mínimo del horizonte (\$[0-9,.]+)/), centavos(modelo.credito.minimo_A))
  comparar('Inicio · Detalle · Disponible por línea', 'Mínimo del horizonte, lectura B', primerDinero(det, /Lectura B \$[0-9,.]+ Activa Mínimo del horizonte (\$[0-9,.]+)/), centavos(modelo.credito.minimo_B))
  void disp

  // ── Inicio · Tendencias · «Devengado contra efectivo» P1–P6 (oct-2026 … mar-2027)
  await page.getByRole('tab', { name: 'Tendencias' }).click()
  const tabla = page.locator('table').filter({ has: page.locator('caption', { hasText: 'Devengado contra efectivo' }) })
  await expect(tabla).toBeVisible()
  const rows = await celdas(page, tabla)
  const cab = rows[0]
  const col = (nombre: string) => cab.findIndex((c) => c.startsWith(nombre))
  const cDev = col('Devengado'), cEfe = col('Efectivo'), cDif = col('Diferencia acumulada')
  const cDevH = cab.findIndex((c) => c === 'Devengado del hogar'), cDevK = cab.findIndex((c) => c === 'Devengado de Karen')
  const cEfeH = cab.findIndex((c) => c === 'Efectivo del hogar'), cEfeK = cab.findIndex((c) => c === 'Efectivo de Karen')
  const periodos = modelo.periodos as { periodo: string; inicio: string; devengado: { total: number; axel: number; karen: number }; efectivo: { total: number; axel: number; karen: number }; diferencia_acumulada: number }[]
  for (const p of periodos) {
    const mes = mesCorto(p.inicio)
    const r = rows.find((x) => x[0] === mes)
    const g = (i: number) => (r && i >= 0 ? aCentavos(r[i]) : NaN)
    const pant = `Inicio · Tendencias · Devengado contra efectivo (${mes})`
    comparar(pant, `${p.periodo} devengado total`, g(cDev), centavos(p.devengado.total))
    comparar(pant, `${p.periodo} efectivo total`, g(cEfe), centavos(p.efectivo.total))
    comparar(pant, `${p.periodo} diferencia acumulada`, g(cDif), centavos(p.diferencia_acumulada))
    comparar(pant, `${p.periodo} devengado hogar`, g(cDevH), centavos(p.devengado.axel))
    comparar(pant, `${p.periodo} devengado Karen`, g(cDevK), centavos(p.devengado.karen))
    comparar(pant, `${p.periodo} efectivo hogar`, g(cEfeH), centavos(p.efectivo.axel))
    comparar(pant, `${p.periodo} efectivo Karen`, g(cEfeK), centavos(p.efectivo.karen))
  }
  // Acumulados al cierre de P6 (horizonte): total_devengado y total_efectivo de model.json
  const r6 = rows.find((x) => x[0] === mesCorto(periodos[5].inicio))
  const cDevAc = cab.findIndex((c) => c === 'Devengado acumulado'), cEfeAc = cab.findIndex((c) => c === 'Efectivo acumulado')
  comparar('Inicio · Tendencias · Devengado contra efectivo (mar-2027)', 'Devengado acumulado a P6', r6 ? aCentavos(r6[cDevAc]) : NaN, centavos(modelo.horizonte.total_devengado))
  comparar('Inicio · Tendencias · Devengado contra efectivo (mar-2027)', 'Efectivo acumulado a P6', r6 ? aCentavos(r6[cEfeAc]) : NaN, centavos(modelo.horizonte.total_efectivo))

  // ── Créditos · línea Kueski · Resumen (tasa, compras) y Detalle (4 despensas)
  await page.goto('/creditos/1')
  await expect(page.getByText('Tasas por plazo')).toBeVisible()
  const res = (await page.locator('main').innerText()).replace(/\s+/g, ' ')
  const m4 = /4 quincenas ([0-9.]+) %/.exec(res)
  comparar('Créditos · línea · Resumen', 'Tasa a 4 quincenas (%, 4 decimales)', m4 ? Math.round(parseFloat(m4[1]) * 10000) : NaN, Math.round(modelo.parametros.tasa_pct * 10000))
  comparar('Créditos · línea · Resumen', 'Próximo pago', primerDinero(res, /Próximo pago (\$[0-9,.]+)/), centavos(ev[1].monto))
  comparar('Créditos · línea · Resumen', 'Disponible lectura A', primerDinero(res, /Disponible · lectura A (\$[0-9,.]+)/), centavos(ev[0].disponible_A))
  comparar('Créditos · línea · Resumen', 'Disponible lectura B', primerDinero(res, /Disponible · lectura B \(activa\) (\$[0-9,.]+)/), centavos(ev[0].disponible_B))
  await page.getByRole('tab', { name: 'Detalle' }).click()
  const tc = page.locator('table').filter({ hasText: 'Intereses (IVA incl.)' }).first()
  await expect(tc).toBeVisible()
  const rc = await celdas(page, tc)
  const ccap = rc[0].findIndex((c) => c.startsWith('Capital')), cint = rc[0].findIndex((c) => c.startsWith('Intereses')), ctot = rc[0].findIndex((c) => c === 'Total')
  for (const d of modelo.despensas as { n: number; capital: number; intereses: number; total: number }[]) {
    const r = rc.find((x) => x[0].startsWith(`Despensa ${d.n}`))
    const pant = `Créditos · línea · Detalle · Despensa ${d.n}`
    comparar(pant, `Despensa ${d.n} capital`, r ? aCentavos(r[ccap]) : NaN, centavos(d.capital))
    comparar(pant, `Despensa ${d.n} intereses`, r ? aCentavos(r[cint]) : NaN, centavos(d.intereses))
    comparar(pant, `Despensa ${d.n} total`, r ? aCentavos(r[ctot]) : NaN, centavos(d.total))
  }

  // ── Créditos · línea · Tendencias · serie de disponible A y B por evento (sin «papel», que no mueve la línea)
  await page.getByRole('tab', { name: 'Tendencias' }).click()
  const ts = page.locator('table').filter({ hasText: 'Lectura A' }).filter({ hasText: 'Capital repuesto' }).first()
  await expect(ts).toBeVisible()
  const rs = await celdas(page, ts)
  const cA = rs[0].findIndex((c) => c === 'Lectura A'), cB = rs[0].findIndex((c) => c === 'Lectura B'), cMo = rs[0].findIndex((c) => c === 'Monto')
  for (const e of ev.filter((x) => x.tipo !== 'papel' && x.tipo !== 'inicio')) {
    const r = rs.find((x) => x[0] === diaCorto(e.fecha))
    const pant = `Créditos · línea · Tendencias · serie (${diaCorto(e.fecha)})`
    comparar(pant, `${e.tipo} ${diaCorto(e.fecha)}: monto`, r ? aCentavos(r[cMo]) : NaN, centavos(e.monto))
    comparar(pant, `${e.tipo} ${diaCorto(e.fecha)}: disponible A`, r ? aCentavos(r[cA]) : NaN, centavos(e.disponible_A))
    comparar(pant, `${e.tipo} ${diaCorto(e.fecha)}: disponible B`, r ? aCentavos(r[cB]) : NaN, centavos(e.disponible_B))
  }

  // ── Flujo y proyección · Deuda viva al 10-abr-2027
  await page.goto('/presupuesto/flujo')
  await page.getByLabel('Fecha de deuda viva').fill('2027-04-10')
  const tarj = page.locator('main').getByText('Lo que queda por pagar después del 10-abr-2027')
  await expect(tarj).toBeVisible()
  const fl = (await page.locator('main').innerText()).replace(/\s+/g, ' ')
  comparar('Flujo y proyección · Resumen', 'Deuda viva al 10-abr-2027', primerDinero(fl, /Deuda viva (\$[0-9,.]+) Lo que queda por pagar después del 10-abr-2027/), centavos(modelo.horizonte.deuda_viva_10_abr))

  // ── Salida
  const dir = path.join(RAIZ, 'datos/tmp/u5-flujos'); mkdirSync(dir, { recursive: true })
  writeFileSync(path.join(dir, 'c3-09.json'), JSON.stringify(filas, null, 1))
  const malas = filas.filter((f) => f.dif !== 0)
  console.log(`C3-09: ${filas.length} cifras leídas en pantalla · diferencias != 0: ${malas.length}`)
  for (const f of filas) console.log(`${String(f.n).padStart(3)} | ${f.pantalla} | ${f.cifra} | pantalla ${f.mostrado} | oráculo ${f.esperado} | dif ${f.dif}`)
  expect(filas.length, 'al menos 10 cifras').toBeGreaterThanOrEqual(10)
  expect(malas, `diferencias: ${JSON.stringify(malas)}`).toEqual([])
})
