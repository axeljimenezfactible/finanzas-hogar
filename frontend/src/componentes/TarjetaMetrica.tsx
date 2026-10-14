// Tarjeta de métrica (DESIGN §7.6): etiqueta, cifra exacta con cifras tabulares, nota (también tabular: suele llevar dinero,
// DESIGN §3.4); tocar abre el detalle.
import type { ReactNode } from 'react'
import { cn } from './utils'
type Rol = 'income' | 'expense' | 'savings' | 'debt' | 'warning'
const COLOR: Record<Rol, string> = { income: 'text-role-income', expense: 'text-role-expense', savings: 'text-role-savings', debt: 'text-role-debt', warning: 'text-role-warning' }
interface Props { titulo: string; valor: string; nota?: ReactNode; rol?: Rol; icono?: ReactNode; alAbrir?: () => void; children?: ReactNode }
export function TarjetaMetrica({ titulo, valor, nota, rol, icono, alAbrir, children }: Props) {
  const Contenedor = alAbrir ? 'button' : 'div'
  return (
    <Contenedor
      {...(alAbrir ? { type: 'button' as const, onClick: alAbrir } : {})}
      className={cn('flex w-full flex-col gap-1 rounded-lg border border-border-divider bg-surface-card p-4 text-left', alAbrir && 'border-border-control hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus')}
    >
      <span className="flex items-center gap-2 text-small font-medium text-text-secondary">{icono}{titulo}</span>
      <span className={cn('whitespace-nowrap text-title-2 font-semibold tabular-nums', rol ? COLOR[rol] : 'text-text-default')}>{valor}</span>
      {nota && <span className="text-small text-text-secondary tabular-nums">{nota}</span>}
      {children}
    </Contenedor>
  )
}
