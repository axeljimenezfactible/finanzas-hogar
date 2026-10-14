// Cascarón: barra lateral, barra superior, franja sin conexión, pantalla registrada y hojas globales.
import { Suspense, useEffect } from 'react'
import { Dialog } from 'radix-ui'
import { textos } from '../textos/es-MX'
import { navegar, useUbicacion } from '../app/enrutador'
import { resolver, type Subseccion } from '../app/rutas'
import { PANTALLAS_REGISTRADAS } from '../app/registro'
import { Cargando, Vacio } from '../componentes/Estados'
import { BarreraErrores } from '../componentes/BarreraErrores'
import { BarraLateral, ContenidoBarra } from './BarraLateral'
import { BarraInferior, BarraSuperior } from './Barras'
import { FranjaSinConexion, PanelCola } from './FranjaSinConexion'
import { HojaRegistro, precargarRegistro } from './HojaRegistro'
import { MenuComandos } from './MenuComandos'
import { abrirMenu, abrirRegistro, cambiar, useShell } from './estado'
const enCampo = (e: KeyboardEvent) => { const t = e.target as HTMLElement; return t.isContentEditable || ['INPUT', 'TEXTAREA', 'SELECT'].includes(t.tagName) }
export function Layout() {
  const { ruta, consulta } = useUbicacion()
  const { panelMovil, barraPlegada } = useShell()
  const r = resolver(ruta)
  useEffect(() => { if (!r) navegar('/inicio', { reemplazar: true }) }, [r])
  useEffect(() => { precargarRegistro() }, [])
  useEffect(() => {
    const h = consulta.get('hoja')
    if (h === 'gasto' || h === 'peso') { abrirRegistro(h); navegar(ruta, { reemplazar: true }) }
  }, [consulta, ruta])
  useEffect(() => {
    const teclas = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && (e.key === 'k' || e.key === 'p')) { e.preventDefault(); abrirMenu() }
      else if (e.altKey && e.key.toLowerCase() === 'g') { e.preventDefault(); abrirRegistro('gasto') }
      else if (e.ctrlKey && e.key === '\\') { e.preventDefault(); cambiar({ barraPlegada: !barraPlegada }) }
      else if (e.key === '/' && !enCampo(e) && !e.ctrlKey) { e.preventDefault(); abrirMenu('/') }
    }
    window.addEventListener('keydown', teclas)
    return () => window.removeEventListener('keydown', teclas)
  }, [barraPlegada])
  if (!r) return null
  const { pantalla: p, params } = r
  const sub = (p.subs?.includes(consulta.get('sub') as Subseccion) ? consulta.get('sub') : p.sub ?? null) as Subseccion | null
  const Pantalla = PANTALLAS_REGISTRADAS[p.id]
  return (
    <div className="flex min-h-dvh">
      <a href="#contenido" className="sr-only focus:not-sr-only focus:absolute focus:z-50 focus:bg-surface-raised focus:p-2">{textos.app.saltarContenido}</a>
      <BarraLateral />
      <Dialog.Root open={panelMovil} onOpenChange={(a) => cambiar({ panelMovil: a })}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 z-40 bg-overlay-scrim tablet:hidden" />
          <Dialog.Content className="fixed inset-y-0 left-0 z-50 w-72 max-w-[85vw] bg-surface-card tablet:hidden" onClick={(e) => { if ((e.target as HTMLElement).closest('a')) cambiar({ panelMovil: false }) }}>
            <Dialog.Title className="sr-only">{textos.barra.navegacion}</Dialog.Title>
            <Dialog.Description className="sr-only">{textos.barra.navegacion}</Dialog.Description>
            <ContenidoBarra />
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>
      <div className="flex min-w-0 flex-1 flex-col pb-[calc(var(--spacing)*14_+_env(safe-area-inset-bottom))] tablet:pb-0">
        <BarraSuperior p={p} sub={sub} />
        <FranjaSinConexion />
        <main id="contenido" className="mx-auto w-full max-w-6xl flex-1 px-4 py-4">
          <BarreraErrores key={p.id + JSON.stringify(params)} clave={`${sub ?? ''}?${consulta.toString()}`}>
            <Suspense fallback={<Cargando />}>
              {Pantalla ? <Pantalla key={p.id + JSON.stringify(params)} params={params} sub={sub} consulta={consulta} /> : <Vacio texto={`${textos.pendiente.titulo}. ${textos.pendiente.cuerpo(p.id)}`} />}
            </Suspense>
          </BarreraErrores>
        </main>
      </div>
      <BarraInferior />
      <HojaRegistro />
      <MenuComandos />
      <PanelCola />
    </div>
  )
}
