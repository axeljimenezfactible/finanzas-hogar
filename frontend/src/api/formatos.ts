// Formatos de pantalla (U1 §0.4, DESIGN §3.3). Dinero en centavos enteros; nunca punto flotante.
const MESES = ['ene', 'feb', 'mar', 'abr', 'may', 'jun', 'jul', 'ago', 'sep', 'oct', 'nov', 'dic']
const DIAS = ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb']
export const ZONA = 'America/Mexico_City'
export const MENOS = '−'

/** 123456 → «$1,234.56»; con signo: 'gasto' → «−$1,234.56», 'ingreso' → «+$1,234.56»; negativos siempre con U+2212. */
export function dinero(centavos: number, signo: 'ninguno' | 'gasto' | 'ingreso' = 'ninguno'): string {
  if (!Number.isSafeInteger(centavos)) throw new Error('centavos no enteros')
  const neg = centavos < 0
  const abs = Math.abs(centavos)
  const enteros = String(Math.floor(abs / 100)).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
  const texto = `$${enteros}.${String(abs % 100).padStart(2, '0')}`
  if (neg || (signo === 'gasto' && abs > 0)) return MENOS + texto
  return signo === 'ingreso' && abs > 0 ? '+' + texto : texto
}

/** Texto tecleado → centavos con ROUND_HALF_UP («185.5» → 18550; «1,000.005» → 100001). null si no es válido. */
export function aCentavos(texto: string): number | null {
  const limpio = texto.replace(/[$,\s]/g, '')
  const m = /^(\d{1,12})(?:\.(\d*))?$/.exec(limpio)
  if (!m) return null
  const dec = (m[2] ?? '').padEnd(3, '0')
  let c = Number(m[1]) * 100 + Number(dec.slice(0, 2))
  if (Number(dec[2]) >= 5) c += 1
  return c
}

/** Kilos tecleados → gramos enteros («72.5» → 72500; ROUND_HALF_UP al gramo). */
export function aGramos(texto: string): number | null {
  const m = /^(\d{1,3})(?:\.(\d*))?$/.exec(texto.replace(/[\s,]/g, '.').replace(/\.(?=.*\.)/g, ''))
  if (!m) return null
  const dec = (m[2] ?? '').padEnd(4, '0')
  let g = Number(m[1]) * 1000 + Number(dec.slice(0, 3))
  if (Number(dec[3]) >= 5) g += 1
  return g > 0 ? g : null
}
/** 72500 → «72.5 kg» (sin ceros de sobra; precisión exacta del registro). */
export const kilos = (g: number) => `${(g / 1000).toFixed(3).replace(/\.?0+$/, '')} kg`

const partes = (fecha: string) => {
  const [a, m, d] = fecha.slice(0, 10).split('-').map(Number)
  return { a, m, d, dia: new Date(Date.UTC(a, m - 1, d)).getUTCDay() }
}
/** «2026-10-10» → «10-oct-2026». */
export function fecha(iso: string): string {
  const { a, m, d } = partes(iso)
  return `${d}-${MESES[m - 1]}-${a}`
}
/** «2026-10-30» → «Vie 30 oct 2026». */
export function fechaConDia(iso: string): string {
  const { a, m, d, dia } = partes(iso)
  return `${DIAS[dia]} ${d} ${MESES[m - 1]} ${a}`
}
/** «2026-10» → «oct-2026». */
export const mes = (aaaamm: string) => `${MESES[Number(aaaamm.slice(5, 7)) - 1]}-${aaaamm.slice(0, 4)}`
/** Porcentaje de avance entero con ROUND_HALF_UP: 79.5 → «80 %». */
export const porcentaje = (n: number) => `${Math.floor(n + 0.5)} %`

function enZona(instante = new Date()) {
  const f = new Intl.DateTimeFormat('en-CA', { timeZone: ZONA, year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit', hourCycle: 'h23' })
  const p = Object.fromEntries(f.formatToParts(instante).map((x) => [x.type, x.value]))
  return { fecha: `${p.year}-${p.month}-${p.day}`, hora: `${p.hour}:${p.minute}:${p.second}`, p }
}
/** Hoy en America/Mexico_City (AAAA-MM-DD), para valores por defecto. */
export const hoyISO = (instante = new Date()) => enZona(instante).fecha
/** Instante ISO 8601 con desplazamiento de America/Mexico_City (C1 DateTime; `captured_at`). */
export function ahoraConDesplazamiento(instante = new Date()): string {
  const { fecha: f, hora, p } = enZona(instante)
  const local = Date.UTC(+p.year, +p.month - 1, +p.day, +p.hour, +p.minute, +p.second)
  const min = Math.round((local - Math.floor(instante.getTime() / 1000) * 1000) / 60000)
  const s = min < 0 ? '-' : '+'
  const a = Math.abs(min)
  return `${f}T${hora}${s}${String(Math.floor(a / 60)).padStart(2, '0')}:${String(a % 60).padStart(2, '0')}`
}
/** Hora local «21:30» de un instante (para «Cargado a las 21:30»). */
export const horaCorta = (instante: Date | number) => enZona(new Date(instante)).hora.slice(0, 5)
/** Fecha-hora de API («2026-10-10T21:30:00-06:00») → «10-oct-2026 21:30». */
export const fechaHora = (iso: string) => `${fecha(enZona(new Date(iso)).fecha)} ${horaCorta(new Date(iso))}`
