// Cálculos de vista previa en el navegador (U4 parte b). Enteros y ROUND_HALF_UP (protocolo §6);
// al guardar manda el servidor. Fórmulas: RN-15, RN-17, RN-18, K7, Q2-16 A, Q2-18 A, Q2-66 B, Q2-70 (a).
import { fecha, fechaConDia, hoyISO, horaCorta } from '../../api/formatos'

/** ROUND_HALF_UP de a ÷ b para enteros (b > 0); simétrico para negativos. */
export function divHalfUp(a: number, b: number): number {
  const s = a < 0 ? -1 : 1
  return s * Math.floor((2 * Math.abs(a) + b) / (2 * b))
}
const partes = (iso: string) => iso.split('-').map(Number) as [number, number, number]
export const diasDelMes = (a: number, m: number) => new Date(Date.UTC(a, m, 0)).getUTCDate()
export function sumarDias(iso: string, n: number): string {
  const [a, m, d] = partes(iso)
  return new Date(Date.UTC(a, m - 1, d + n)).toISOString().slice(0, 10)
}
/** Día ISO 8601: 1 = lunes … 7 = domingo. */
export function diaIso(iso: string): number {
  const [a, m, d] = partes(iso)
  const w = new Date(Date.UTC(a, m - 1, d)).getUTCDay()
  return w === 0 ? 7 : w
}
export const esHora = (t: string) => /^([01]\d|2[0-3]):[0-5]\d$/.test(t)

export type Programa = { tipo: 'diario' } | { tipo: 'semanal'; dia: number } | { tipo: 'mensual'; dia: number }
/** Próximo envío con el día y la hora tecleados (hora del centro de México). Q2-66 B: en un mes corto, el último día. */
export function proximoEnvio(p: Programa, hora: string, ahora = new Date()): { fecha: string; hora: string } | null {
  if (!esHora(hora)) return null
  const hoy = hoyISO(ahora)
  const ya = horaCorta(ahora)
  for (let i = 0; i < 400; i++) {
    const f = sumarDias(hoy, i)
    if (i === 0 && hora <= ya) continue
    const [a, m, d] = partes(f)
    const ok = p.tipo === 'diario' || (p.tipo === 'semanal' ? diaIso(f) === p.dia : d === Math.min(p.dia, diasDelMes(a, m)))
    if (ok) return { fecha: f, hora }
  }
  return null
}
export const textoEnvio = (r: { fecha: string; hora: string } | null) => (r ? `${fechaConDia(r.fecha)}, ${r.hora}` : '—')
/** Ejemplo de mes corto para el día 29, 30 o 31: el próximo febrero que no tiene ese día. */
export function ejemploMesCorto(dia: number, hoy = hoyISO()): string {
  let [a, m] = partes(hoy)
  for (let i = 0; i < 48; i++) {
    if (diasDelMes(a, m) < dia) return fecha(`${a}-${String(m).padStart(2, '0')}-${String(diasDelMes(a, m)).padStart(2, '0')}`)
    m += 1
    if (m > 12) { m = 1; a += 1 }
  }
  return ''
}

// ── Salud (RN-17, K7, K23, K36) ──
export const FACTOR_MILESIMAS: Record<string, number> = { sedentary: 1200, light: 1375, moderate: 1550, active: 1725, very_active: 1900 }
export interface InsumosSalud { sex?: string | null; birth_date?: string | null; height_mm?: number | null; weight_g?: number | null; activity_level?: string | null; goal?: string | null; ajusteBajar: number; ajusteSubir: number }
export function edad(nacimiento: string, hoy = hoyISO()): number {
  const [a, m, d] = partes(nacimiento)
  const [ha, hm, hd] = partes(hoy)
  return ha - a - (hm < m || (hm === m && hd < d) ? 1 : 0)
}
/** IMC en décimas: weight_g × 1000 ÷ height_mm², una vez con HALF_UP. */
export const imcDecimas = (g: number, mm: number) => divHalfUp(g * 10000, mm * mm)
export const claseImc = (dec: number) => (dec < 185 ? 'underweight' : dec < 250 ? 'normal' : dec < 300 ? 'overweight' : 'obesity')
const num = new Intl.NumberFormat('es-MX', { maximumFractionDigits: 6 })
export const numero = (n: number) => num.format(n)
/** Meta sugerida (RN-17): Mifflin-St Jeor × factor ∓ ajuste, un solo HALF_UP al final. Devuelve lo que falta si no se puede. */
export function metaSugerida(x: InsumosSalud, hoy = hoyISO()): { kcal: number; desglose: string } | { falta: string } {
  const orden: [keyof InsumosSalud, string][] = [['sex', 'sexo'], ['birth_date', 'nacimiento'], ['height_mm', 'estatura'], ['weight_g', 'peso'], ['activity_level', 'actividad'], ['goal', 'objetivo']]
  for (const [k, n] of orden) if (x[k] === undefined || x[k] === null || x[k] === '') return { falta: n }
  const basal800 = 8 * x.weight_g! + 500 * x.height_mm! - 4000 * edad(x.birth_date!, hoy) + 800 * (x.sex === 'male' ? 5 : -161)
  const f = FACTOR_MILESIMAS[x.activity_level!]
  const ajuste = x.goal === 'lose' ? -x.ajusteBajar : x.goal === 'gain' ? x.ajusteSubir : 0
  const exacto = basal800 * f + ajuste * 800000
  const kcal = divHalfUp(exacto, 800000)
  const tramo = `${numero(basal800 / 800)} × ${numero(f / 1000)} = ${numero((basal800 * f) / 800000)}`
  const conAjuste = ajuste ? ` ${ajuste < 0 ? '−' : '+'} ${numero(Math.abs(ajuste))} = ${numero(exacto / 800000)}` : ''
  return { kcal, desglose: `Basal ${tramo}${conAjuste} → ${numero(kcal)} kcal` }
}

// ── Metas de ahorro (RN-18, M-16, Q2-16 A, Q2-17 C, Q2-18 A) ──
export const avanceDecimas = (ahorrado: number, objetivo: number) => (objetivo > 0 ? divHalfUp(ahorrado * 1000, objetivo) : 0)
export const decimas = (d: number) => `${Math.trunc(d / 10)}.${Math.abs(d % 10)}`
export function mesesQueFaltan(fechaMeta: string, hoy = hoyISO()): number {
  const [a, m] = partes(fechaMeta)
  const [ha, hm] = partes(hoy)
  return Math.max(0, a * 12 + m - (ha * 12 + hm) + 1)
}
export function estadoMeta(ahorrado: number, objetivo: number | undefined | null, fechaMeta: string | undefined | null, hoy = hoyISO()): 'achieved' | 'overdue' | 'in_progress' {
  if (objetivo && ahorrado >= objetivo) return 'achieved'
  if (fechaMeta && fechaMeta < hoy) return 'overdue'
  return 'in_progress'
}
/** Aportación necesaria; null = vencida sin cumplir (PROVISIONAL, bloqueo S2-metas-vencida). */
export function aportacionNecesaria(ahorrado: number, objetivo: number, fechaMeta: string, hoy = hoyISO()): number | null {
  const resto = objetivo - ahorrado
  if (resto <= 0) return 0
  const meses = mesesQueFaltan(fechaMeta, hoy)
  return meses === 0 ? null : divHalfUp(resto, meses)
}

// ── Inventario (RN-15, K7, Q2-70 a) en milésimas de la unidad base ──
export interface InsumosInventario { tracking_mode: string; quantity_milli: number; counted_date?: string | null; rate_quantity_milli?: number | null; rate_period_days?: number | null; reorder_level_milli?: number | null }
const diasEntre = (a: string, b: string) => Math.round((Date.parse(`${b}T00:00:00Z`) - Date.parse(`${a}T00:00:00Z`)) / 86400000)
export function existenciaHoy(x: InsumosInventario, hoy = hoyISO()): number | null {
  if (x.tracking_mode !== 'consumption_rate') return x.quantity_milli
  if (!x.rate_quantity_milli || !x.rate_period_days) return null
  const dias = x.counted_date ? Math.max(0, diasEntre(x.counted_date, hoy)) : 0
  return Math.max(0, x.quantity_milli - divHalfUp(x.rate_quantity_milli * dias, x.rate_period_days))
}
export function estadoInventario(existencia: number | null, nivel?: number | null): 'sufficient' | 'reorder' | 'out' | null {
  if (existencia === null) return null
  if (existencia === 0) return 'out'
  if (nivel !== undefined && nivel !== null && existencia <= nivel) return 'reorder'
  return 'sufficient'
}
export function diasRestantes(x: InsumosInventario, hoy = hoyISO()): number | null {
  const e = existenciaHoy(x, hoy)
  if (x.tracking_mode !== 'consumption_rate' || e === null) return null
  return Math.floor((e * x.rate_period_days!) / x.rate_quantity_milli!)
}
export const consumoSemanal = (x: InsumosInventario) => (x.tracking_mode === 'consumption_rate' && x.rate_quantity_milli && x.rate_period_days ? divHalfUp(x.rate_quantity_milli * 7, x.rate_period_days) : null)
const cant = new Intl.NumberFormat('es-MX', { maximumFractionDigits: 3 })
export const cantidad = (milli: number) => cant.format(milli / 1000)
/** «1.5» → 1500 milésimas (hasta 3 decimales, sin flotantes). */
export function aMilesimas(t: string): number | null {
  const s = t.trim().replace(/,/g, '')
  if (!/^\d+(\.\d{1,3})?$/.test(s)) return null
  const [e, d = ''] = s.split('.')
  return Number(e) * 1000 + Number(d.padEnd(3, '0'))
}
/** Porcentaje con hasta 2 decimales → puntos base («50» → 5000). */
export function aPuntosBase(t: string): number | null {
  const s = t.trim().replace('%', '').trim()
  if (!/^\d+(\.\d{1,2})?$/.test(s)) return null
  const [e, d = ''] = s.split('.')
  return Number(e) * 100 + Number(d.padEnd(2, '0'))
}
export const pb = (bp: number) => `${Math.trunc(bp / 100)}.${String(bp % 100).padStart(2, '0')} %`
