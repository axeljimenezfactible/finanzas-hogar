// Botón (DESIGN §7.2) con el patrón de shadcn/ui (cva + cn) y solo tokens C2.
import type { ComponentProps } from 'react'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from './utils'
export const estilosBoton = cva(
  'inline-flex select-none items-center justify-center gap-2 rounded-md text-small font-medium transition-colors duration-150 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus aria-disabled:cursor-not-allowed aria-disabled:border aria-disabled:border-dashed aria-disabled:border-border-control aria-disabled:bg-surface-disabled aria-disabled:text-text-disabled',
  {
    variants: {
      variante: {
        primario: 'bg-action-primary-bg text-text-on-primary hover:bg-action-primary-bg-hover active:bg-action-primary-bg-pressed',
        secundario: 'border border-border-control bg-surface-card text-text-default hover:border-border-control-hover hover:bg-surface-hover active:bg-surface-pressed',
        fantasma: 'text-text-default hover:bg-surface-hover active:bg-surface-pressed',
        destructivo: 'bg-action-destructive-bg text-text-on-destructive hover:bg-action-destructive-bg-hover active:bg-action-destructive-bg-pressed',
        enlace: 'px-0 text-text-link underline-offset-2 hover:text-text-link-hover hover:underline',
      },
      tam: { normal: 'h-9 min-w-11 px-3 coarse:h-11', hoja: 'h-12 w-full px-4 text-body', icono: 'size-9 coarse:size-11' },
    },
    defaultVariants: { variante: 'secundario', tam: 'normal' },
  },
)
type Props = ComponentProps<'button'> & VariantProps<typeof estilosBoton> & { /** Motivo en texto: deja el botón visible pero inactivo (§7.0). */ deshabilitadoPor?: string }
export function Boton({ variante, tam, className, deshabilitadoPor, onClick, children, type = 'button', ...resto }: Props) {
  return (
    <button
      type={type}
      className={cn(estilosBoton({ variante, tam }), className)}
      aria-disabled={deshabilitadoPor ? true : undefined}
      title={deshabilitadoPor}
      onClick={deshabilitadoPor ? (e) => e.preventDefault() : onClick}
      {...resto}
    >
      {children}
      {deshabilitadoPor && <span className="sr-only"> · {deshabilitadoPor}</span>}
    </button>
  )
}
