// Hojas y diálogos (DESIGN §7.12, U1 §1.3): hoja a pantalla completa en móvil, ventana modal en escritorio. Radix Dialog.
import type { ReactNode } from 'react'
import { Dialog } from 'radix-ui'
import { X } from 'lucide-react'
import { ICONO_MD } from './Icono'
import { textos } from '../textos/es-MX'
interface Props { abierta: boolean; alCambiar: (a: boolean) => void; titulo: string; descripcion?: string; children: ReactNode; pie?: ReactNode; alAbrirFoco?: (e: Event) => void }
export function Hoja({ abierta, alCambiar, titulo, descripcion, children, pie, alAbrirFoco }: Props) {
  return (
    <Dialog.Root open={abierta} onOpenChange={alCambiar}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-40 bg-overlay-scrim" />
        <Dialog.Content
          onOpenAutoFocus={alAbrirFoco}
          aria-describedby={descripcion ? undefined : undefined}
          className="fixed inset-0 z-50 flex flex-col bg-surface-raised text-text-default tablet:inset-auto tablet:top-1/2 tablet:left-1/2 tablet:max-h-[85dvh] tablet:w-full tablet:max-w-lg tablet:-translate-x-1/2 tablet:-translate-y-1/2 tablet:rounded-lg tablet:border tablet:border-border-divider tablet:shadow-2"
        >
          <div className="flex items-center gap-2 border-b border-border-divider px-4 py-2">
            <Dialog.Title className="flex-1 text-title-3 font-semibold">{titulo}</Dialog.Title>
            <Dialog.Close aria-label={textos.app.cerrar} className="inline-flex size-9 items-center justify-center rounded-md hover:bg-surface-hover coarse:size-11"><X aria-hidden className="size-6" size={ICONO_MD} /></Dialog.Close>
          </div>
          {descripcion ? <Dialog.Description className="px-4 pt-2 text-small text-text-secondary">{descripcion}</Dialog.Description> : <Dialog.Description className="sr-only">{titulo}</Dialog.Description>}
          <div className="flex-1 overflow-y-auto px-4 py-3">{children}</div>
          {pie && <div className="border-t border-border-divider px-4 py-3 pb-[max(env(safe-area-inset-bottom),12px)]">{pie}</div>}
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  )
}
