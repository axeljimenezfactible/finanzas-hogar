// Estado global del cascarón: hoja de registro (P-03 / P-22), menú «/» y barra lateral.
import { useSyncExternalStore } from 'react'
export type TipoRegistro = 'gasto' | 'ingreso' | 'peso' | 'ticket'
interface Estado { registro: TipoRegistro | null; menu: string | null; barraPlegada: boolean; panelMovil: boolean; panelCola: boolean }
const leerPlegada = () => { try { return localStorage.getItem('fh-barra-plegada') === '1' } catch { return false } }
let estado: Estado = { registro: null, menu: null, barraPlegada: leerPlegada(), panelMovil: false, panelCola: false }
const oyentes = new Set<() => void>()
export function cambiar(p: Partial<Estado>) {
  estado = { ...estado, ...p }
  if ('barraPlegada' in p) try { localStorage.setItem('fh-barra-plegada', p.barraPlegada ? '1' : '0') } catch { /* nada */ }
  oyentes.forEach((f) => f())
}
export const abrirRegistro = (tipo: TipoRegistro = 'gasto') => cambiar({ registro: tipo, menu: null, panelMovil: false })
export const abrirMenu = (texto = '') => cambiar({ menu: texto, panelMovil: false })
export const useShell = () => useSyncExternalStore((f) => { oyentes.add(f); return () => oyentes.delete(f) }, () => estado)
/** Avisa a las pantallas que deben volver a pedir sus datos (p. ej. Inicio tras guardar un gasto, F-1 paso 4). */
export const avisarCambio = (tipo: 'gasto' | 'peso') => window.dispatchEvent(new CustomEvent('fh:datos-cambiaron', { detail: { tipo } }))
