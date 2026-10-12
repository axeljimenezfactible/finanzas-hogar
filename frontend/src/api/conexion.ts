// Estado de conexión con el servidor de casa (G-07, N-004): navigator.onLine + último resultado de red.
import { useSyncExternalStore } from 'react'
let conServidor = typeof navigator === 'undefined' ? true : navigator.onLine
const oyentes = new Set<() => void>()
let sondeo: ReturnType<typeof setInterval> | null = null
const avisar = () => oyentes.forEach((f) => f())
export const hayConexion = () => conServidor && navigator.onLine
export function marcarConServidor() {
  if (conServidor) return
  conServidor = true
  if (sondeo) clearInterval(sondeo)
  sondeo = null
  avisar()
  window.dispatchEvent(new Event('fh:conexion-recuperada'))
}
export function marcarSinServidor() {
  if (!conServidor) return
  conServidor = false
  avisar()
  // Sondeo de /health cada 15 s mientras no haya servidor (PROPUESTA K31: valor configurable).
  sondeo ??= setInterval(() => {
    if (!navigator.onLine) return
    fetch('/health', { cache: 'no-store' }).then((r) => r.ok && marcarConServidor(), () => undefined)
  }, 15_000)
}
if (typeof window !== 'undefined') {
  window.addEventListener('online', () => { conServidor = false; marcarConServidor() })
  window.addEventListener('offline', () => { conServidor = true; marcarSinServidor() })
}
const suscribir = (f: () => void) => { oyentes.add(f); return () => oyentes.delete(f) }
export const useConexion = () => useSyncExternalStore(suscribir, hayConexion)
