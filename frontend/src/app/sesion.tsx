// Carga de sesión (K11: sesión larga por dispositivo con cookie HttpOnly fh_session; el frontend nunca ve el token).
// Sin conexión, usa la última sesión cargada en el teléfono para poder capturar gastos y peso (K9).
import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from 'react'
import { pedir } from '../api/cliente'
import { esSinConexion } from '../api/errores'
import { borrarLecturas, leerLectura } from '../api/lecturas'
import type { Esquemas } from '../api/tipos'
import { navegar } from './enrutador'
export type Sesion = Esquemas['CurrentSession']
type Estado = 'cargando' | 'sin-sesion' | 'activa'
interface Ctx { sesion: Sesion | null; estado: Estado; vencida: boolean; entrar: (usuario: string, contrasena: string) => Promise<void>; salir: () => Promise<void> }
const Contexto = createContext<Ctx | null>(null)
export const useSesion = () => useContext(Contexto)!
const CLAVE_VOLVER = 'fh-volver-a'
export function ProveedorSesion({ children }: { children: ReactNode }) {
  const [sesion, setSesion] = useState<Sesion | null>(null)
  const [estado, setEstado] = useState<Estado>('cargando')
  const [vencida, setVencida] = useState(false)
  useEffect(() => {
    pedir('get', '/api/v1/auth/me').then(
      (s) => { setSesion(s); setEstado('activa') },
      async (e) => {
        const copia = esSinConexion(e) ? await leerLectura('/api/v1/auth/me') : undefined
        if (copia) { setSesion(copia.datos as Sesion); setEstado('activa') } else setEstado('sin-sesion')
      },
    )
    const alVencer = () => {
      try { sessionStorage.setItem(CLAVE_VOLVER, location.pathname + location.search) } catch { /* nada */ }
      setVencida(true); setSesion(null); setEstado('sin-sesion')
    }
    window.addEventListener('fh:sesion-vencida', alVencer)
    return () => window.removeEventListener('fh:sesion-vencida', alVencer)
  }, [])
  const entrar = useCallback(async (username: string, password: string) => {
    const s = await pedir('post', '/api/v1/auth/login', { cuerpo: { username, password } })
    setSesion(s); setEstado('activa'); setVencida(false)
    let volver = '/inicio'
    try { volver = sessionStorage.getItem(CLAVE_VOLVER) ?? volver; sessionStorage.removeItem(CLAVE_VOLVER) } catch { /* nada */ }
    navegar(volver.startsWith('/entrar') ? '/inicio' : volver, { reemplazar: true })
  }, [])
  const salir = useCallback(async () => {
    try { await pedir('post', '/api/v1/auth/logout') } finally {
      await borrarLecturas() // la cola NO se borra: las capturas pendientes nunca se pierden (U1 §2.4)
      setSesion(null); setEstado('sin-sesion'); navegar('/entrar', { reemplazar: true })
    }
  }, [])
  return <Contexto.Provider value={{ sesion, estado, vencida, entrar, salir }}>{children}</Contexto.Provider>
}
