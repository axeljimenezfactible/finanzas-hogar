// Cola sin conexión (K9, U1-15, C1-19, U5-13): SOLO gastos y peso. Vive en IndexedDB del teléfono,
// sobrevive a recargar, se envía sola al volver la red, en orden de captura y sin duplicar: cada captura
// conserva su Idempotency-Key y su captured_at (momento de captura, con desplazamiento) en cada reintento.
import { useSyncExternalStore } from 'react'
import { nuevaClave, pedir } from '../api/cliente'
import { ErrorApi, esSinConexion } from '../api/errores'
import { ahoraConDesplazamiento } from '../api/formatos'
import { hayConexion } from '../api/conexion'
import type { Esquemas } from '../api/tipos'
import { poner, quitar, todos } from './idb'
import { textos } from '../textos/es-MX'

export type CuerpoGasto = Esquemas['ExpenseCreate']
export type CuerpoPeso = Esquemas['WeightEntryCreate']
interface Base { clave: string; creadoEn: number; estado: 'pendiente' | 'fallido'; intentos: number; error: string | null }
export type Captura = Base & ({ tipo: 'gasto'; cuerpo: CuerpoGasto } | { tipo: 'peso'; cuerpo: CuerpoPeso })

let cache: Captura[] = []
const oyentes = new Set<() => void>()
const canal = typeof BroadcastChannel !== 'undefined' ? new BroadcastChannel('fh-cola') : null
async function refrescar(difundir = true) {
  try { cache = await todos<Captura>('cola') } catch { cache = [] }
  oyentes.forEach((f) => f())
  if (difundir) canal?.postMessage('cambio')
}
canal?.addEventListener('message', () => void refrescar(false))

/** Crea la captura con su clave y su captured_at, y la guarda en el teléfono. */
async function guardar(c: Captura) {
  await poner('cola', c)
  await refrescar()
  return c
}
export function encolarGasto(cuerpo: CuerpoGasto, clave = nuevaClave()) {
  return guardar({ clave, tipo: 'gasto', cuerpo: { ...cuerpo, captured_at: cuerpo.captured_at ?? ahoraConDesplazamiento() }, creadoEn: Date.now(), estado: 'pendiente', intentos: 0, error: null })
}
export function encolarPeso(cuerpo: CuerpoPeso, clave = nuevaClave()) {
  return guardar({ clave, tipo: 'peso', cuerpo: { ...cuerpo, captured_at: cuerpo.captured_at ?? ahoraConDesplazamiento() }, creadoEn: Date.now(), estado: 'pendiente', intentos: 0, error: null })
}
/** Quitar de la cola (Deshacer o «Descartar» confirmado por la persona; nunca sin confirmación). */
export async function quitarDeCola(clave: string) {
  await quitar('cola', clave)
  await refrescar()
}
export async function reintentar(clave: string) {
  const c = cache.find((x) => x.clave === clave)
  if (c) await poner('cola', { ...c, estado: 'pendiente', error: null })
  await refrescar()
  return enviarCola()
}

/** enviadas: llegaron al servidor · fallidas: rechazadas (4xx, quedan «fallido») · retenidas: falla del servidor o
 *  sesión vencida (5xx, 401) por primera vez: siguen «pendiente», se reintentan solas y se avisa «No se pudo enviar · Revisar». */
export interface ResultadoEnvio { enviadas: number; fallidas: number; retenidas: number }
let enviando: Promise<ResultadoEnvio> | null = null
async function ejecutar(): Promise<ResultadoEnvio> {
  const r: ResultadoEnvio = { enviadas: 0, fallidas: 0, retenidas: 0 }
  for (const c of await todos<Captura>('cola')) {
    if (c.estado === 'fallido') continue
    try {
      if (c.tipo === 'gasto') await pedir('post', '/api/v1/expenses', { cuerpo: c.cuerpo, idempotencia: c.clave })
      else await pedir('post', '/api/v1/weight-entries', { cuerpo: c.cuerpo, idempotencia: c.clave })
      await quitar('cola', c.clave)
      r.enviadas++
    } catch (e) {
      // Sin red: la captura sigue pendiente y se reintenta al volver la conexión (la franja ya dice «Sin conexión»).
      if (esSinConexion(e)) break
      // Sesión vencida o servidor con falla (U5-13 paso 7, U1 §2.4): sigue pendiente y se reintenta sola, pero queda
      // marcada con el motivo para que la franja la cuente («No se pudo enviar · n en cola») y se avise una sola vez
      // por captura (los reintentos automáticos de 30 s no repiten el aviso; «Reintentar envío» sí).
      if (e instanceof ErrorApi && (e.estado === 401 || e.estado >= 500)) {
        if (c.error === null) r.retenidas++
        await poner('cola', { ...c, estado: 'pendiente', intentos: c.intentos + 1, error: e.message || textos.errores.generico })
        break
      }
      await poner('cola', { ...c, estado: 'fallido', intentos: c.intentos + 1, error: e instanceof ErrorApi ? e.message : textos.errores.generico })
      r.fallidas++
    }
  }
  await refrescar()
  if (r.enviadas || r.fallidas || r.retenidas) window.dispatchEvent(new CustomEvent<ResultadoEnvio>('fh:cola-enviada', { detail: r }))
  return r
}
/** Envía lo pendiente. Un solo envío a la vez (también entre pestañas con Web Locks). */
export function enviarCola(): Promise<ResultadoEnvio> {
  if (!enviando) {
    const p = navigator.locks ? (navigator.locks.request('fh-cola-envio', ejecutar) as unknown as Promise<ResultadoEnvio>) : ejecutar()
    enviando = p.finally(() => { enviando = null })
  }
  return enviando
}

let iniciada = false
/** Arranque (una vez, en App): carga la cola y programa los envíos automáticos. */
export function iniciarCola() {
  if (iniciada) return
  iniciada = true
  const intentar = () => { if (hayConexion() && cache.some((c) => c.estado === 'pendiente')) void enviarCola() }
  void refrescar(false).then(intentar)
  window.addEventListener('online', intentar)
  window.addEventListener('fh:conexion-recuperada', intentar)
  // PROPUESTA (K31): reintento cada 30 s mientras haya pendientes y conexión.
  setInterval(intentar, 30_000)
}

const suscribir = (f: () => void) => { oyentes.add(f); return () => oyentes.delete(f) }
const instantanea = () => cache
export function useCola() {
  const capturas = useSyncExternalStore(suscribir, instantanea)
  return {
    capturas,
    total: capturas.length,
    fallidas: capturas.filter((c) => c.estado === 'fallido').length,
    /** Pendientes cuyo último envío falló por el servidor o la sesión (5xx, 401); se reintentan solas. */
    retenidas: capturas.filter((c) => c.estado === 'pendiente' && c.error !== null).length,
  }
}
