// Chips de persona y tienda (DESIGN §7.4, K41): forma + nombre + posición, nunca solo color.
import type { ComponentProps, ReactNode } from 'react'
import { Check } from 'lucide-react'
import { cn } from './utils'
const PERSONA: Record<string, string> = { axel: 'bg-person-axel', ximena: 'bg-person-ximena', karen: 'bg-person-karen' }
const TIENDA: Record<string, string> = { walmart: 'bg-store-walmart', 'bodega-aurrera': 'bg-store-bodega-aurrera', 'sams-club': 'bg-store-sams-club', market: 'bg-store-market' }
/** Orden fijo de personas y tiendas (DESIGN §7.4). */
export const ORDEN_PERSONAS = ['axel', 'ximena', 'karen', 'shared']
export const ORDEN_TIENDAS = ['walmart', 'bodega-aurrera', 'sams-club', 'market']
export function MarcaPersona({ codigo }: { codigo: string }) {
  return codigo in PERSONA
    ? <span aria-hidden className={cn('inline-block size-2.5 shrink-0 rounded-full', PERSONA[codigo])} />
    : <span aria-hidden className="inline-block size-2.5 shrink-0 rounded-full border-2 border-person-shared" />
}
export function MarcaTienda({ codigo }: { codigo: string }) {
  return <span aria-hidden className={cn('inline-block size-2.5 shrink-0', TIENDA[codigo] ?? 'bg-store-market')} />
}
const base = 'inline-flex h-6 items-center gap-1 border border-border-divider bg-surface-raised px-2 text-small text-text-default'
export const ChipPersona = ({ codigo, nombre }: { codigo: string; nombre: string }) => (
  <span className={cn(base, 'rounded-full')}><MarcaPersona codigo={codigo} />{nombre}</span>
)
export const ChipTienda = ({ codigo, nombre }: { codigo: string; nombre: string }) => (
  <span className={cn(base, 'rounded-sm')}><MarcaTienda codigo={codigo} />{nombre}</span>
)
type Sel = ComponentProps<'button'> & { seleccionado: boolean; marca?: ReactNode; forma?: 'pildora' | 'rectangulo' }
/** Chip seleccionable (filtros y categorías de F-1): <button aria-pressed>, 44 px con puntero grueso. */
export function ChipSeleccionable({ seleccionado, marca, forma = 'pildora', className, children, ...resto }: Sel) {
  return (
    <button
      type="button"
      aria-pressed={seleccionado}
      className={cn(base, 'h-8 cursor-pointer px-3 enabled:hover:bg-surface-hover disabled:cursor-not-allowed disabled:border-dashed disabled:border-border-control disabled:bg-surface-disabled disabled:text-text-disabled focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus coarse:h-11', forma === 'pildora' ? 'rounded-full' : 'rounded-sm', seleccionado && 'border-border-selected bg-surface-selected', className)}
      {...resto}
    >
      {seleccionado && <Check aria-hidden className="size-5 text-icon-selected" />}
      {marca}
      {children}
    </button>
  )
}
