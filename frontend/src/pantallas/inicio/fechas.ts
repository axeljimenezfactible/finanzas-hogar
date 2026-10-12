// Fechas de Inicio (K8: quincenas 1–15 y 16–último día; semana ISO de lunes a domingo; K38: los domingos la lista vigente es la de la semana siguiente).
// Todo en fechas de calendario AAAA-MM-DD (America/Mexico_City vía hoyISO), sin horas ni flotantes.
const MESES = ['ene', 'feb', 'mar', 'abr', 'may', 'jun', 'jul', 'ago', 'sep', 'oct', 'nov', 'dic']
const utc = (iso: string) => { const [a, m, d] = iso.slice(0, 10).split('-').map(Number); return Date.UTC(a, m - 1, d) }
const aIso = (ms: number) => new Date(ms).toISOString().slice(0, 10)
export const sumarDias = (iso: string, n: number) => aIso(utc(iso) + n * 86_400_000)
/** Días de `a` a `b` (b − a). */
export const difDias = (a: string, b: string) => Math.round((utc(b) - utc(a)) / 86_400_000)
/** 0 = lunes … 6 = domingo. */
export const diaSemanaLunes = (iso: string) => (new Date(utc(iso)).getUTCDay() + 6) % 7
export const lunesDe = (iso: string) => sumarDias(iso, -diaSemanaLunes(iso))
export const ultimoDia = (aaaamm: string) => new Date(Date.UTC(Number(aaaamm.slice(0, 4)), Number(aaaamm.slice(5, 7)), 0)).getUTCDate()
/** Semana ISO 8601 «2026-W41» de una fecha. */
export function semanaISO(iso: string): string {
  const jueves = sumarDias(lunesDe(iso), 3)
  const anio = Number(jueves.slice(0, 4))
  const n = Math.floor(difDias(`${anio}-01-01`, jueves) / 7) + 1
  return `${anio}-W${String(n).padStart(2, '0')}`
}
/** Lista vigente (K38, U1 §0.2): los domingos, la semana que empieza el lunes siguiente. */
export const semanaVigente = (iso: string) => semanaISO(diaSemanaLunes(iso) === 6 ? sumarDias(iso, 1) : iso)
export const mesDe = (iso: string) => iso.slice(0, 7)
export const quincenaDe = (iso: string) => `${mesDe(iso)}-${Number(iso.slice(8, 10)) <= 15 ? 'q1' : 'q2'}`
export const diaDe = (iso: string) => Number(iso.slice(8, 10))
export const mesAbrev = (iso: string) => MESES[Number(iso.slice(5, 7)) - 1]
/** «2026-10» ± n meses. */
export function sumarMeses(aaaamm: string, n: number) {
  const t = Number(aaaamm.slice(0, 4)) * 12 + Number(aaaamm.slice(5, 7)) - 1 + n
  return `${Math.floor(t / 12)}-${String((t % 12) + 1).padStart(2, '0')}`
}
/** «Lun 5 al Dom 11 oct» (o «Lun 28 sep al Dom 4 oct» si cambia el mes). */
export function rangoSemana(inicio: string, fin: string) {
  const mismo = inicio.slice(5, 7) === fin.slice(5, 7)
  return `Lun ${diaDe(inicio)}${mismo ? '' : ` ${mesAbrev(inicio)}`} al Dom ${diaDe(fin)} ${mesAbrev(fin)}`
}
/** Día de pago de quincena (15 y 30, o el último día del mes si es menor; K8). */
export const esDiaDePago = (iso: string) => { const d = diaDe(iso); return d === 15 || d === Math.min(30, ultimoDia(mesDe(iso))) }
/** Rejilla de un mes de lunes a domingo (K8): semanas con 7 fechas; las de otros meses en null. */
export function rejillaMes(aaaamm: string): (string | null)[][] {
  const primero = `${aaaamm}-01`
  const u = ultimoDia(aaaamm)
  const celdas: (string | null)[] = Array.from({ length: diaSemanaLunes(primero) }, () => null)
  for (let d = 1; d <= u; d++) celdas.push(`${aaaamm}-${String(d).padStart(2, '0')}`)
  while (celdas.length % 7) celdas.push(null)
  return Array.from({ length: celdas.length / 7 }, (_, i) => celdas.slice(i * 7, i * 7 + 7))
}
