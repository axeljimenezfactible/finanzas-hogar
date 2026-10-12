// G-01 Barra lateral (U1 §1.2, DESIGN §7.9): 9 secciones en orden fijo, ▸ despliega, nombre abre la primera
// subpágina; Créditos lista sus líneas; colapsable (‹‹ o Ctrl+\) y panel flotante al acercar el puntero al borde.
import { useState } from 'react'
import { ChevronRight, ChevronsLeft, LogOut, Monitor, Moon, Plus, Search, Sun, UserRound } from 'lucide-react'
import { ICONO_MD } from '../componentes/Icono'
import { textos } from '../textos/es-MX'
import { usarConsulta } from '../api/consulta'
import type { Esquemas } from '../api/tipos'
import { Enlace, useUbicacion } from '../app/enrutador'
import { PRIMERA, SECCIONES, SECCIONES_AJUSTES, pantalla, resolver, rutaDe, type IdSeccion } from '../app/rutas'
import { useSesion } from '../app/sesion'
import { ponerTema, siguienteTema, useTema } from '../app/tema'
import { cn } from '../componentes/utils'
import { abrirMenu, abrirRegistro, cambiar, useShell } from './estado'
const leerAbiertas = (): Record<string, boolean> => { try { return JSON.parse(localStorage.getItem('fh-barra-abiertas') ?? '{}') } catch { return {} } }
const fila = 'flex h-8 w-full items-center gap-2 rounded-sm px-2 text-small text-text-default hover:bg-surface-hover coarse:h-11'
export function InterruptorTema() {
  const tema = useTema()
  const Icono = tema === 'claro' ? Sun : tema === 'oscuro' ? Moon : Monitor
  return (
    <button type="button" className={fila} onClick={() => ponerTema(siguienteTema(tema))} aria-label={`${textos.tema.cambiar}: ${textos.tema[tema]}`} data-prueba="interruptor-tema">
      <Icono aria-hidden className="size-5 text-icon-default" />{textos.tema.titulo}: {textos.tema[tema]}
    </button>
  )
}
export function ContenidoBarra() {
  const { ruta } = useUbicacion()
  const { sesion, salir } = useSesion()
  const actual = resolver(ruta)?.pantalla
  const [abiertas, setAbiertas] = useState(leerAbiertas)
  const [menuUsuario, setMenuUsuario] = useState(false)
  const lineas = usarConsulta('get', '/api/v1/credits')
  const alternar = (id: string, abierta: boolean) => { const n = { ...abiertas, [id]: !abierta }; setAbiertas(n); try { localStorage.setItem('fh-barra-abiertas', JSON.stringify(n)) } catch { /* nada */ } }
  const hijosDe = (id: IdSeccion) => {
    const base = SECCIONES.find((s) => s.id === id)!.hijos.map((h) => ({ clave: h, nombre: pantalla(h).nombre, a: rutaDe(h) }))
    if (id === 'creditos') base.splice(1, 0, ...((lineas.datos?.items ?? []) as Esquemas['Credit'][]).map((l) => ({ clave: `c${l.id}`, nombre: l.name, a: rutaDe('P-18', { id: l.id }) })))
    if (id === 'ajustes') return SECCIONES_AJUSTES.map((s) => ({ clave: s.id, nombre: `${s.id} ${s.nombre}`, a: `/ajustes/${s.id}` }))
    return base
  }
  return (
    <nav aria-label={textos.barra.navegacion} className="flex h-full flex-col gap-1 p-2">
      <div className="relative">
        <button type="button" className={cn(fila, 'h-10 font-semibold')} aria-expanded={menuUsuario} onClick={() => setMenuUsuario(!menuUsuario)}>
          <UserRound aria-hidden className="size-5" /><span className="flex-1 truncate text-left">{textos.app.nombre}</span><span className="text-caption font-normal text-text-secondary">{sesion?.display_name}</span>
        </button>
        {menuUsuario && (
          <div className="absolute inset-x-0 top-full z-10 flex flex-col rounded-md border border-border-divider bg-surface-raised p-1 shadow-1">
            <Enlace a="/ajustes/01" className={fila}>{textos.barra.cuenta}</Enlace>
            <Enlace a="/ajustes/10" className={fila}>{textos.barra.apariencia}</Enlace>
            <button type="button" className={fila} onClick={() => void salir()}><LogOut aria-hidden className="size-5" />{textos.barra.cerrarSesion}</button>
          </div>
        )}
      </div>
      <button type="button" className={fila} onClick={() => abrirMenu()}><Search aria-hidden className="size-5 text-icon-default" />{textos.barra.buscar}</button>
      <button type="button" className={fila} onClick={() => abrirRegistro('gasto')}><Plus aria-hidden className="size-5 text-icon-default" />{textos.barra.registrarGasto}</button>
      <ul className="mt-2 flex flex-1 flex-col gap-0.5 overflow-y-auto">
        {SECCIONES.map((s) => {
          const nombre = textos.secciones[s.id]
          const hijos = hijosDe(s.id)
          const abierta = abiertas[s.id] ?? (s.id !== 'ajustes' && actual?.seccion === s.id)
          return (
            <li key={s.id}>
              <div className={cn('flex items-center rounded-sm', actual?.seccion === s.id && !hijos.length && 'bg-surface-selected')}>
                {hijos.length > 0 ? (
                  <button type="button" aria-expanded={abierta} aria-label={abierta ? textos.barra.plegarSeccion(nombre) : textos.barra.desplegarSeccion(nombre)} onClick={() => alternar(s.id, abierta)} className="inline-flex size-8 shrink-0 items-center justify-center rounded-sm hover:bg-surface-hover coarse:size-11">
                    <ChevronRight aria-hidden className={cn('size-5 transition-transform duration-150', abierta && 'rotate-90')} />
                  </button>
                ) : <span className="size-8 shrink-0 coarse:size-11" />}
                <Enlace a={rutaDe(PRIMERA[s.id], { seccion: '01' })} className={cn(fila, 'font-medium')} aria-current={actual?.seccion === s.id && !hijos.length ? 'page' : undefined}>{nombre}</Enlace>
              </div>
              {abierta && hijos.length > 0 && (
                <ul className="ml-8 flex flex-col gap-0.5">
                  {hijos.map((h) => <li key={h.clave}><Enlace a={h.a} className={cn(fila, ruta === h.a && 'bg-surface-selected font-medium')} aria-current={ruta === h.a ? 'page' : undefined}>{h.nombre}</Enlace></li>)}
                </ul>
              )}
            </li>
          )
        })}
      </ul>
      <InterruptorTema />
    </nav>
  )
}
export function BarraLateral() {
  const { barraPlegada } = useShell()
  const [flotante, setFlotante] = useState(false)
  if (barraPlegada) return (
    <>
      <div className="fixed inset-y-0 left-0 z-30 hidden w-2 tablet:block" onMouseEnter={() => setFlotante(true)} />
      {flotante && <aside onMouseLeave={() => setFlotante(false)} className="fixed inset-y-0 left-0 z-40 hidden w-64 border-r border-border-divider bg-surface-raised shadow-2 tablet:block"><ContenidoBarra /></aside>}
    </>
  )
  return (
    <aside className="sticky top-0 hidden h-dvh w-64 shrink-0 flex-col border-r border-border-divider bg-surface-card tablet:flex">
      <button type="button" aria-label={textos.barra.plegar} title="Ctrl+\" onClick={() => cambiar({ barraPlegada: true })} className="absolute top-2 right-2 inline-flex size-8 items-center justify-center rounded-sm hover:bg-surface-hover"><ChevronsLeft aria-hidden className="size-6" size={ICONO_MD} /></button>
      <ContenidoBarra />
    </aside>
  )
}
