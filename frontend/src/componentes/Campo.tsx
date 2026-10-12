// Campo de formulario (DESIGN §7.3): etiqueta, control, error con aria-invalid y aria-describedby.
import { useId, type ComponentProps } from 'react'
import { CircleAlert } from 'lucide-react'
import { cn } from './utils'
export const estilosControl = 'h-9 w-full rounded-md border border-border-control bg-surface-field px-3 text-body text-text-default placeholder:text-text-secondary hover:border-border-control-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus aria-invalid:border-border-error coarse:h-11'
type Props = ComponentProps<'input'> & { etiqueta: string; error?: string; ayuda?: string; ocultarEtiqueta?: boolean }
export function Campo({ etiqueta, error, ayuda, ocultarEtiqueta, className, id, ...resto }: Props) {
  const auto = useId()
  const idc = id ?? auto
  const desc = error ? `${idc}-error` : ayuda ? `${idc}-ayuda` : undefined
  return (
    <div className={cn('flex flex-col gap-1', className)}>
      <label htmlFor={idc} className={cn('text-small font-medium text-text-secondary', ocultarEtiqueta && 'sr-only')}>{etiqueta}</label>
      <input id={idc} aria-invalid={error ? true : undefined} aria-describedby={desc} className={estilosControl} {...resto} />
      {error ? (
        <p id={desc} className="flex items-center gap-1 text-small text-role-error"><CircleAlert aria-hidden className="size-5" />{error}</p>
      ) : ayuda ? <p id={desc} className="text-small text-text-secondary">{ayuda}</p> : null}
    </div>
  )
}
