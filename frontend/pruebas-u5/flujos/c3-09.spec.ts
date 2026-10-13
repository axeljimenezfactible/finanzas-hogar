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
/** Compara una cifra de pantalla con su esperado: dinero en centavos (por omisión) o un porcentaje entero (`'%'`). */
function comparar(pantalla: string, cifra: string, mostradoC: number, esperadoC: number, unidad: '$' | '%' = '$') {
  const ver = (x: number) => (unidad === '%' ? `${x} %` : (x / 100).toFixed(2))
  const f = { n: filas.length + 1, pantalla, cifra, mostrado: Number.isNaN(mostradoC) ? 'NO ENCONTRADA' : ver(mostradoC), esperado: ver(esperadoC), dif: Number.isNaN(mostradoC) ? NaN : mostradoC - esperadoC }
  filas.push(f)
}
const norm = (s: string) => s.replace(/\s+/g, ' ').trim()
/** Hoy en la zona de la casa (AAAA-MM-DD), igual que hoyISO() de la app. */
const hoyCasa = () => new Date().toLocaleDateString('en-CA', { timeZone: 'America/Mexico_City' })
/** ROUND_HALF_UP de a ÷ b con enteros ≥ 0 (protocolo §6). */
const divRedondeo = (a: number, b: number) => Math.floor((2 * a + b) / (2 * b))
async function celdas(page: Page, tabla: ReturnType<Page['locator']>) {
  return tabla.locator('tr').evaluateAll((trs) => trs.map((tr) => [...tr.querySelectorAll('th,td')].map((c) => (c as HTMLElement).innerText.replace(/\s+/g, ' ').trim())))
}
const primerDinero = (t: string, despuesDe: RegExp) => { const m = despuesDe.exec(t); return m ? aCentavos(m[1]) : NaN }

test('C3-09 · cifras de Inicio, créditos y flujo = oráculo (model.json)', async ({ page }) => {
  const ev = modelo.credito.eventos as { fecha: string; tipo: string; monto: number; disponible_A: number; disponible_B: number }[]
  await entrar(page)
  const hoy = hoyCasa(), mesHoy = hoy.slice(0, 7)
  type Despensa = { n: number; fecha: string; total: number; pagos: { fecha: string; monto: number; axel: number; karen: number }[] }
  const despensas = modelo.despensas as Despensa[]
  // Oráculo de la deuda viva a una fecha D (igual que horizonte.deuda_viva_10_abr): compras con fecha ≤ D, pagos con fecha > D.
  // Decisión 1 del orquestador (WF-3a cierre, PROVISIONAL): la deuda viva «hoy» de K-3 y P-17 = getLiveDebt?date=hoy.
  const deudaVivaAl = (d: string) => despensas.filter((x) => x.fecha <= d).flatMap((x) => x.pagos).filter((q) => q.fecha > d).reduce((a, q) => a + centavos(q.monto), 0)
  // Pagos del mes en curso de todas las líneas (RN-30, M-19; decisión 2): suma de los pagos con fecha en el mes.
  const pagosMes = despensas.flatMap((x) => x.pagos).filter((q) => q.fecha.startsWith(mesHoy))
  const suma = (k: 'monto' | 'axel' | 'karen') => pagosMes.reduce((a, q) => a + centavos(q[k]), 0)
  const MESES_LARGOS = ['enero', 'febrero', 'marzo', 'abril', 'mayo', 'junio', 'julio', 'agosto', 'septiembre', 'octubre', 'noviembre', 'diciembre']
  const mesLargo = MESES_LARGOS[Number(mesHoy.slice(5)) - 1]

  // ── Inicio · Resumen · K-4 «Pagos próximos» (U3 2.0): primera fila = pago del 30-oct-2026 con sus partes (R-13)
  const p1 = despensas[0].pagos[0]
  const k4 = page.getByTestId('tarjeta-K-4')
  const fila1 = k4.getByRole('link', { name: /^Vie 30 oct 2026/ })
  await expect(fila1).toBeVisible()
  // La fila se lee ya asentada: con las personas cargadas aparece «Karen» (su parte > 0 en el oráculo, U3 K-4). Mientras no cargan,
  // la fila dice un instante «Hogar $973.75» sin Karen (todo al hogar): queda como solicitud fuera de alcance para Tablero.tsx K-4.
  await expect(fila1).toContainText('Karen')
  const tk4 = norm(await fila1.innerText())
  console.log(`K-4 · primera fila: «${tk4}»`)
  expect(tk4).toContain(diaCorto(p1.fecha).replace(/-/g, ' ')) // «30 oct 2026»
  comparar('Inicio · K-4 Pagos próximos', 'Próximo pago (30-oct-2026)', primerDinero(tk4, /^Vie 30 oct 2026 · [^$]*?(\$[0-9,.]+)/), centavos(p1.monto))
  comparar('Inicio · K-4 Pagos próximos', 'Parte del hogar del pago 1', primerDinero(tk4, /Hogar (\$[0-9,.]+)/), centavos(p1.axel))
  comparar('Inicio · K-4 Pagos próximos', 'Parte de Karen del pago 1', primerDinero(tk4, /Karen (\$[0-9,.]+)/), centavos(p1.karen))

  // ── Inicio · Resumen · K-3 «Créditos»: deuda viva de hoy, pagos del mes y partes, disponible hoy (lectura B)
  const k3 = page.getByTestId('tarjeta-K-3')
  await expect(k3).toContainText('deuda viva total')
  await expect(k3).toContainText('Disponible hoy') // getDashboardNextPayment ya respondió (bloque «Disponible hoy»)
  await expect(k3).toContainText('Karen')
  const tk3 = norm(await k3.innerText())
  console.log(`K-3: «${tk3}»`)
  comparar('Inicio · K-3 Créditos', `Deuda viva de hoy (${hoy})`, primerDinero(tk3, /(\$[0-9,.]+) deuda viva total/), deudaVivaAl(hoy))
  comparar('Inicio · K-3 Créditos', `Pagos de ${mesLargo} (todas las líneas)`, primerDinero(tk3, new RegExp(`Pagos de ${mesLargo}: (\\$[0-9,.]+)`)), suma('monto'))
  comparar('Inicio · K-3 Créditos', `Parte del hogar de los pagos de ${mesLargo}`, primerDinero(tk3, /Hogar (\$[0-9,.]+)/), suma('axel'))
  comparar('Inicio · K-3 Créditos', `Parte de Karen de los pagos de ${mesLargo}`, primerDinero(tk3, /Karen (\$[0-9,.]+)/), suma('karen'))
  comparar('Inicio · K-3 Créditos', 'Disponible hoy, lectura B', primerDinero(tk3, /Disponible hoy (\$[0-9,.]+)/), centavos(ev[0].disponible_B))

  // ── Inicio · Resumen · B-1 barómetro: pagos del hogar del mes = oráculo; % y margen = fórmula de U3 §B-1 sobre las cifras
  // de la pantalla (el ingreso no está en el oráculo: model.json no trae ingresos).
  const b1 = page.getByTestId('tarjeta-B-1')
  await expect(b1.getByTestId('barometro-digito')).toBeVisible()
  const tb1 = norm(await b1.innerText())
  const pagosB1 = primerDinero(tb1, /Pagos (\$[0-9,.]+) de un ingreso de/), ingreso = primerDinero(tb1, /de un ingreso de (\$[0-9,.]+)/)
  const tope = Number(/de ([0-9]+) % de tope/.exec(tb1)?.[1] ?? NaN)
  const digito = Number(/^([0-9]+) %$/.exec(norm(await b1.getByTestId('barometro-digito').innerText()))?.[1] ?? NaN)
  comparar('Inicio · B-1 Capacidad de endeudamiento', `Pagos del hogar de ${mesLargo}`, pagosB1, suma('axel'))
  comparar('Inicio · B-1 Capacidad de endeudamiento', 'Capacidad usada (pagos × 100 ÷ ingreso, ROUND_HALF_UP)', digito, divRedondeo(suma('axel') * 100, ingreso), '%')
  comparar('Inicio · B-1 Capacidad de endeudamiento', 'Te caben (tope × ingreso ÷ 100 − pagos)', primerDinero(tb1, /Te caben (\$[0-9,.]+) más al mes/), divRedondeo(tope * ingreso, 100) - suma('axel'))

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

  // ── Créditos · P-17 Todas las líneas · «Deuda viva total» y la de la línea, a hoy (decisión 1: = getLiveDebt?date=hoy)
  await page.goto('/creditos')
  await expect(page.getByText('Deuda viva total').first()).toBeVisible()
  await page.waitForTimeout(800)
  const p17 = norm(await page.locator('main').innerText())
  comparar('Créditos · P-17 Todas las líneas', `Deuda viva total de hoy (${hoy})`, primerDinero(p17, /Deuda viva total (\$[0-9,.]+)/), deudaVivaAl(hoy))
  comparar('Créditos · P-17 Todas las líneas', `Deuda viva de la línea de hoy (${hoy})`, primerDinero(p17, /Deuda viva (?!total)(\$[0-9,.]+)/), deudaVivaAl(hoy))

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
  expect(filas.filter((f) => /^(Inicio · K-|Inicio · B-1|Créditos · P-17)/.test(f.pantalla)).length, 'C3-09 de U3 2.0: al menos 10 cifras del tablero y de P-17').toBeGreaterThanOrEqual(10)
  expect(malas, `diferencias: ${JSON.stringify(malas)}`).toEqual([])
})
