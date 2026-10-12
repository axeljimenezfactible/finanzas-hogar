// Cálculos en el navegador para las vistas previas de Ajustes (U4 §0.2): centavos enteros, ROUND_HALF_UP, sin flotantes en dinero.
import type { Esquemas } from '../../api/tipos'
type Partida = Pick<Esquemas['PlannedItem'], 'frequency' | 'interval_count' | 'start_date' | 'end_date' | 'amount_cents' | 'second_half_amount_cents'>
type Regla = Esquemas['CostShareRule']

/** ROUND_HALF_UP(a × b ÷ c) con enteros exactos (a, b ≥ 0; c > 0). */
export function multDiv(a: number, b: number, c: number): number {
  const n = BigInt(a) * BigInt(b), d = BigInt(c)
  return Number((2n * n + d) / (2n * d))
}

/** «25.8685» → 258685/1000000 con 4 decimales (tasa como fracción exacta, diccionario §sufijos). */
export function aFraccion(texto: string, maxDec = 4): { num: number; den: number; decimales: number } | null {
  const s = texto.trim().replace(/\s*%$/, '')
  const m = /^(\d+)(?:\.(\d+))?$/.exec(s)
  if (!m || (m[2]?.length ?? 0) > maxDec) return null
  const dec = m[2]?.length ?? 0
  return { num: Number(m[1] + (m[2] ?? '')), den: 10 ** (dec + 2), decimales: dec }
}
/** Porcentaje con hasta 2 decimales → puntos base («12.5» → 1250). */
export const aPuntosBase = (texto: string): number | null => { const f = aFraccion(texto, 2); return f ? multDiv(f.num, 10000, f.den) : null }
export const pctDePuntosBase = (bp: number) => `${(bp / 100).toFixed(2).replace(/\.?0+$/, '')} %`
/** Texto «25.8685 %» con los decimales capturados (K12). */
export function pctDeFraccion(num: number, den: number, decimales = 2): string {
  const escala = 10 ** decimales
  const v = multDiv(num, 100 * escala, den)
  const ent = Math.floor(v / escala), frac = String(v % escala).padStart(decimales, '0')
  return decimales ? `${ent}.${frac} %` : `${ent} %`
}

// ——— Fechas ISO (AAAA-MM-DD) sin zona: aritmética de calendario.
const partes = (iso: string) => iso.split('-').map(Number) as [number, number, number]
const iso = (a: number, m: number, d: number) => `${a}-${String(m).padStart(2, '0')}-${String(d).padStart(2, '0')}`
export const ultimoDia = (a: number, m: number) => new Date(Date.UTC(a, m, 0)).getUTCDate()
export function sumarDias(f: string, n: number): string { const [a, m, d] = partes(f); const x = new Date(Date.UTC(a, m - 1, d + n)); return iso(x.getUTCFullYear(), x.getUTCMonth() + 1, x.getUTCDate()) }
function sumarMeses(f: string, n: number, dia: number): string { const [a, m] = partes(f); const t = a * 12 + (m - 1) + n; const aa = Math.floor(t / 12), mm = (t % 12) + 1; return iso(aa, mm, Math.min(dia, ultimoDia(aa, mm))) }
export const mesDe = (f: string) => f.slice(0, 7)
export function mesSiguiente(aaaamm: string): string { return sumarMeses(`${aaaamm}-01`, 1, 1).slice(0, 7) }

/** Días de pago 15 y 30 (o último día) desde una fecha (RN-06). */
function quincenasDesde(f: string): string[] {
  const [a, m] = partes(f); const r: string[] = []
  for (let i = 0; i < 400; i++) { const t = a * 12 + m - 1 + Math.floor(i / 2); const aa = Math.floor(t / 12), mm = (t % 12) + 1; r.push(iso(aa, mm, i % 2 ? Math.min(30, ultimoDia(aa, mm)) : 15)) }
  return r.filter((x) => x >= f)
}

/** Fechas de una partida planificada (regla del esquema 00001, U4 P-26.02 vista previa 3 (a)), hasta un tope. */
export function fechasDe(p: Partida, hasta: string): string[] {
  if (!p.start_date) return []
  const k = Math.max(1, p.interval_count || 1), fin = p.end_date && p.end_date < hasta ? p.end_date : hasta
  const r: string[] = []
  if (p.frequency === 'once') return p.start_date <= fin ? [p.start_date] : []
  if (p.frequency === 'fortnightly') { const q = quincenasDesde(p.start_date); for (let i = 0; i < q.length && q[i] <= fin; i += k) r.push(q[i]); return r }
  const dia = partes(p.start_date)[2]
  for (let i = 0; i < 2000; i++) {
    const f = p.frequency === 'weekly' ? sumarDias(p.start_date, 7 * k * i) : sumarMeses(p.start_date, (p.frequency === 'yearly' ? 12 : 1) * k * i, dia)
    if (f > fin) break
    r.push(f)
  }
  return r
}
export const proximas = (p: Partida, desde: string, n = 3) => fechasDe(p, sumarDias(desde, 3660)).filter((f) => f >= desde).slice(0, n)
export const montoEn = (p: Partida, f: string) => (p.frequency === 'fortnightly' && p.second_half_amount_cents != null && partes(f)[2] !== 15 ? p.second_half_amount_cents : p.amount_cents)
export const fechasEnMes = (p: Partida, aaaamm: string) => fechasDe(p, `${aaaamm}-31`).filter((f) => mesDe(f) === aaaamm)
export const totalEnMes = (p: Partida, aaaamm: string) => fechasEnMes(p, aaaamm).reduce((s, f) => s + montoEn(p, f), 0)

/** Reparto de un monto con una regla (RN-03, K14): % con ROUND_HALF_UP, fijo tal cual, el resto al titular. null si una parte fija pasa del monto. */
export function repartir(monto: number, regla: Regla | undefined): { persona: number; cents: number }[] | null {
  if (!regla) return []
  const partesR = (regla.parts ?? []).filter((x) => x.person_id !== regla.remainder_person_id)
  const r = partesR.map((x) => ({ persona: x.person_id, cents: x.fixed_cents ?? multDiv(monto, x.pct_bp ?? 0, 10000) }))
  const usado = r.reduce((s, x) => s + x.cents, 0)
  if (usado > monto) return null
  return [...r, { persona: regla.remainder_person_id, cents: monto - usado }]
}

/** Primera fecha ≥ desde cuyo día es `dia` (recortado al último del mes). */
export function proximoDia(desde: string, dia: number, estricto = false): string {
  for (let i = 0; i < 3; i++) { const f = sumarMeses(`${mesDe(desde)}-01`, i, dia); if (estricto ? f > desde : f >= desde) return f }
  return sumarMeses(`${mesDe(desde)}-01`, 3, dia)
}
/** Calendario de la tarjeta (S3-interfaz §5.3, Q2-32 A): corte y pago de ese corte. */
export function corteYPago(hoy: string, corte: number, pago: number) { const c = proximoDia(hoy, corte); return { corte: c, pago: proximoDia(c, pago, true) } }

/** Peso de un precio por antigüedad (K22, K40): 1 antes de viejo; (cero − días) ÷ (cero − viejo + 1) después; 0 desde «cero». */
export function pesoPrecio(dias: number, viejo: number, cero: number): string {
  if (dias < viejo) return '1'
  if (dias >= cero) return '0'
  return ((cero - dias) / (cero - viejo + 1)).toFixed(4)
}
/** Pesos normalizados en % (decisión C). */
export function normalizar(pesos: number[]): number[] { const s = pesos.reduce((a, b) => a + b, 0); return s ? pesos.map((p) => (p * 100) / s) : pesos.map(() => 0) }
