// Tema: por defecto según el sistema, con interruptor (decisión G, U5-10). Se guarda por dispositivo.
// Claro y oscuro los resuelven las variables --ui-* de tokens.css con [data-theme]; nunca inversión de color.
import { useSyncExternalStore } from 'react'
export type Tema = 'sistema' | 'claro' | 'oscuro'
const CLAVE = 'fh-tema'
const oyentes = new Set<() => void>()
export function leerTema(): Tema { try { const t = localStorage.getItem(CLAVE); return t === 'claro' || t === 'oscuro' ? t : 'sistema' } catch { return 'sistema' } }
export function ponerTema(t: Tema) {
  try { if (t === 'sistema') localStorage.removeItem(CLAVE); else localStorage.setItem(CLAVE, t) } catch { /* sin almacenamiento */ }
  if (t === 'sistema') document.documentElement.removeAttribute('data-theme')
  else document.documentElement.setAttribute('data-theme', t === 'claro' ? 'light' : 'dark')
  oyentes.forEach((f) => f())
}
/** /tema alterna sistema → claro → oscuro → sistema (U1 §5.2). */
export const siguienteTema = (t: Tema): Tema => (t === 'sistema' ? 'claro' : t === 'claro' ? 'oscuro' : 'sistema')
export const useTema = () => useSyncExternalStore((f) => { oyentes.add(f); return () => oyentes.delete(f) }, leerTema)
