// G-02 barra superior (migas, título, pestañas de subsección) y G-05 barra inferior (solo < 768 px).
import { ChevronLeft, Home, Menu, Plus, Search, ShoppingCart, Utensils } from 'lucide-react'
import { ICONO_MD } from '../componentes/Icono'
import { textos } from '../textos/es-MX'
import { Enlace, conConsulta, useUbicacion } from '../app/enrutador'
import { PRIMERA, pantalla, rutaDe, type Pantalla, type Subseccion } from '../app/rutas'
import { cn } from '../componentes/utils'
import { abrirMenu, abrirRegistro, cambiar, useShell } from './estado'
export function BarraSuperior({ p, sub }: { p: Pantalla; sub: Subseccion | null }) {
  const { barraPlegada } = useShell()
  const padre = p.padre ? pantalla(p.padre) : null
  const seccion = p.seccion ? textos.secciones[p.seccion] : ''
  return (
    <header className="flex flex-col gap-2 border-b border-border-divider bg-surface-page px-4 pt-2">
      <div className="flex items-center gap-2">
        <button type="button" aria-label={textos.barra.abrirMenu} onClick={() => cambiar({ panelMovil: true })} className="inline-flex size-11 items-center justify-center rounded-sm hover:bg-surface-hover tablet:hidden"><Menu aria-hidden className="size-6" size={ICONO_MD} /></button>
        {barraPlegada && <button type="button" aria-label={textos.barra.desplegar} onClick={() => cambiar({ barraPlegada: false })} className="hidden size-9 items-center justify-center rounded-sm hover:bg-surface-hover tablet:inline-flex coarse:size-11"><Menu aria-hidden className="size-6" size={ICONO_MD} /></button>}
        {padre && <Enlace a={rutaDe(padre.id)} className="inline-flex min-h-11 items-center gap-1 text-small text-text-secondary tablet:hidden"><ChevronLeft aria-hidden className="size-5" />{padre.nombre}</Enlace>}
        <nav aria-label="Migas" className="hidden text-small text-text-secondary tablet:block">
          <ol className="flex gap-1">
            {p.seccion && <li><Enlace a={rutaDe(PRIMERA[p.seccion], { seccion: '01' })} className="hover:underline">{seccion}</Enlace></li>}
            {padre && <li>/ <Enlace a={rutaDe(padre.id)} className="hover:underline">{padre.nombre}</Enlace></li>}
            {p.nombre !== seccion && <li aria-current="page">/ {p.nombre}</li>}
          </ol>
        </nav>
        <button type="button" onClick={() => abrirMenu()} aria-label={textos.barra.buscar} className="ml-auto hidden size-9 items-center justify-center rounded-sm hover:bg-surface-hover tablet:inline-flex coarse:size-11"><Search aria-hidden className="size-6" size={ICONO_MD} /></button>
      </div>
      <h1 className="text-title-1 font-bold text-text-default">{p.nombre}</h1>
      {p.subs && (
        <div role="tablist" aria-label={p.nombre} className="-mb-px flex gap-1 overflow-x-auto">
          {p.subs.map((s) => (
            <Enlace key={s} role="tab" aria-selected={sub === s} a={conConsulta({ sub: s === p.sub ? null : s })}
              className={cn('inline-flex h-9 shrink-0 items-center border-b-2 px-2 text-small coarse:h-11', sub === s ? 'border-border-selected font-medium text-text-default' : 'border-transparent text-text-secondary hover:text-text-default')}>{textos.subsecciones[s]}</Enlace>
          ))}
        </div>
      )}
    </header>
  )
}
// DESIGN §7.10: fondo ui.surface.raised; destinos de 56 px de alto y ≥ 64 px de ancho; icono size.icon.md (24 px) en
// ui.text.secondary; el destino actual lleva la píldora de 56 × 32 (size.bottombar.indicator-*) en ui.surface.selected
// detrás del icono ui.icon.selected, y el nombre en ui.text.default con peso medio.
const destino = 'group flex min-w-16 flex-1 flex-col items-center justify-center gap-1 text-caption text-text-secondary focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-border-focus aria-[current=page]:font-medium aria-[current=page]:text-text-default'
const pildora = 'inline-flex h-8 w-14 items-center justify-center rounded-full group-aria-[current=page]:bg-surface-selected group-aria-[current=page]:text-icon-selected'
export function BarraInferior() {
  const { ruta } = useUbicacion()
  const enlace = (id: string, nombre: string, Icono: typeof Home) => {
    const a = rutaDe(id)
    return <Enlace a={a} className={destino} aria-current={ruta.startsWith(a) ? 'page' : undefined}><span data-prueba="pildora" className={pildora}><Icono aria-hidden className="size-6" size={ICONO_MD} /></span>{nombre}</Enlace>
  }
  return (
    <nav aria-label={textos.barra.inferior} className="fixed inset-x-0 bottom-0 z-30 flex h-14 border-t border-border-divider bg-surface-raised pb-[env(safe-area-inset-bottom)] tablet:hidden">
      {enlace('P-01', textos.secciones.inicio, Home)}
      {enlace('P-10', textos.secciones.comidas, Utensils)}
      <button type="button" aria-label={textos.barra.registrarGasto} data-prueba="boton-mas" onClick={() => abrirRegistro('gasto')} className={destino}>
        <span className="inline-flex size-10 items-center justify-center rounded-full bg-action-primary-bg text-text-on-primary"><Plus aria-hidden className="size-6" size={ICONO_MD} /></span>
      </button>
      {enlace('P-13', textos.barra.compras, ShoppingCart)}
      <button type="button" onClick={() => abrirMenu()} className={destino}><span className={pildora}><Search aria-hidden className="size-6" size={ICONO_MD} /></span>{textos.barra.buscar}</button>
    </nav>
  )
}
