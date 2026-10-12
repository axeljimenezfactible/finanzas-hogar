// Avisos breves (DESIGN §7.15, G-07): uno visible a la vez; éxito 5 s (pausa con hover o foco); error sin cierre automático.
import { createContext, useCallback, useContext, useEffect, useRef, useState, type ReactNode } from 'react'
import { CircleAlert, CircleCheck, Info, TriangleAlert, Undo2, X } from 'lucide-react'
import { ICONO_MD } from './Icono'
import { textos } from '../textos/es-MX'
export type RolAviso = 'exito' | 'error' | 'advertencia' | 'info'
export interface Aviso { texto: string; rol?: RolAviso; accion?: { texto: string; alHacer: () => void | Promise<void> } }
const Ctx = createContext<(a: Aviso) => void>(() => undefined)
export const useAvisos = () => useContext(Ctx)
const ICONO = { exito: CircleCheck, error: CircleAlert, advertencia: TriangleAlert, info: Info }
const COLOR = { exito: 'text-role-success', error: 'text-role-error', advertencia: 'text-role-warning', info: 'text-role-info' }
export function ProveedorAvisos({ children }: { children: ReactNode }) {
  const [cola, setCola] = useState<Aviso[]>([])
  const [pausa, setPausa] = useState(false)
  const actual = cola[0]
  const cerrar = useCallback(() => setCola((c) => c.slice(1)), [])
  const mostrar = useCallback((a: Aviso) => setCola((c) => [...c, a]), [])
  const t = useRef<ReturnType<typeof setTimeout>>(undefined)
  useEffect(() => {
    if (!actual || pausa || actual.rol === 'error') return
    t.current = setTimeout(cerrar, 5000)
    return () => clearTimeout(t.current)
  }, [actual, pausa, cerrar])
  const rol = actual?.rol ?? 'exito'
  const Icono = ICONO[rol]
  return (
    <Ctx.Provider value={mostrar}>
      {children}
      <div className="pointer-events-none fixed inset-x-2 bottom-20 z-50 flex justify-center tablet:inset-x-auto tablet:right-6 tablet:bottom-6">
        {actual && (
          <div
            role={rol === 'error' ? 'alert' : 'status'}
            aria-live={rol === 'error' ? 'assertive' : 'polite'}
            onMouseEnter={() => setPausa(true)} onMouseLeave={() => setPausa(false)} onFocus={() => setPausa(true)} onBlur={() => setPausa(false)}
            className="pointer-events-auto flex w-full max-w-120 items-center gap-2 rounded-md border border-border-divider bg-surface-raised px-3 py-2 text-small text-text-default shadow-2"
          >
            <Icono aria-hidden className={`size-5 shrink-0 ${COLOR[rol]}`} />
            <span className="flex-1">{actual.texto}</span>
            {actual.accion && (
              <button type="button" className="inline-flex h-9 items-center gap-1 px-2 font-medium text-text-link coarse:h-11" onClick={async () => { const a = actual.accion!; cerrar(); await a.alHacer() }}>
                {actual.accion.texto === textos.app.deshacer && <Undo2 aria-hidden className="size-5" />}{actual.accion.texto}
              </button>
            )}
            <button type="button" aria-label={textos.app.cerrar} className="inline-flex size-9 items-center justify-center coarse:size-11" onClick={cerrar}><X aria-hidden className="size-6" size={ICONO_MD} /></button>
          </div>
        )}
      </div>
    </Ctx.Provider>
  )
}
