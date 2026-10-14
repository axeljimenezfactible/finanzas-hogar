// Enrutador mínimo sobre History API (sin dependencias). navegar() y <Enlace> mantienen la SPA.
import { useSyncExternalStore, type ComponentProps } from 'react'
const EVENTO = 'fh:navegacion'
const suscribir = (f: () => void) => { window.addEventListener('popstate', f); window.addEventListener(EVENTO, f); return () => { window.removeEventListener('popstate', f); window.removeEventListener(EVENTO, f) } }
const instantanea = () => location.pathname + location.search
export function useUbicacion() {
  const actual = useSyncExternalStore(suscribir, instantanea)
  const u = new URL(actual, location.origin)
  return { ruta: u.pathname, consulta: u.searchParams, completa: actual }
}
export function navegar(a: string, opciones: { reemplazar?: boolean } = {}) {
  if (a === instantanea()) return
  history[opciones.reemplazar ? 'replaceState' : 'pushState'](null, '', a)
  window.dispatchEvent(new Event(EVENTO))
  window.scrollTo(0, 0)
}
/** Cambia un parámetro de la consulta (sub=, vista=, hoja=) sin perder los demás. */
export function conConsulta(cambios: Record<string, string | null>) {
  const u = new URL(location.href)
  for (const [k, v] of Object.entries(cambios)) v === null ? u.searchParams.delete(k) : u.searchParams.set(k, v)
  return u.pathname + (u.search ? u.search : '')
}
export function Enlace({ a, onClick, ...resto }: ComponentProps<'a'> & { a: string }) {
  return <a href={a} onClick={(e) => { onClick?.(e); if (e.defaultPrevented || e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey) return; e.preventDefault(); navegar(a) }} {...resto} />
}
