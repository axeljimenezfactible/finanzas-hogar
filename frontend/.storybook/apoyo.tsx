// Apoyo de las historias: rejillas de variantes con solo clases de los tokens de C2.
import type { ReactNode } from 'react'

/** Bloque con título (h2) que agrupa una fila o una cuadrícula de variantes. */
export function Grupo({ titulo, nota, children }: { titulo: string; nota?: string; children: ReactNode }) {
  return (
    <section className="flex flex-col gap-3">
      <h2 className="text-title-3 font-semibold text-text-default">{titulo}</h2>
      {nota && <p className="text-small text-text-secondary">{nota}</p>}
      {children}
    </section>
  )
}

/** Fila horizontal que se parte en líneas (variantes una junto a otra). */
export const Fila = ({ children, className = '' }: { children: ReactNode; className?: string }) => (
  <div className={`flex flex-wrap items-center gap-2 ${className}`}>{children}</div>
)

/** Variante con su etiqueta, para que el nombre del estado sea texto y no solo posición. */
export function Etiquetada({ etiqueta, children }: { etiqueta: string; children: ReactNode }) {
  return (
    <div className="flex flex-col items-start gap-1">
      <span className="text-caption text-text-secondary">{etiqueta}</span>
      {children}
    </div>
  )
}

/** Superficies donde se mide cada componente (DESIGN.md §7.1): página, tarjeta y lo que flota. */
export function Superficies({ children }: { children: (nombre: string) => ReactNode }) {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col gap-2 rounded-md bg-surface-page p-3 border border-border-divider">
        <span className="text-caption text-text-secondary">Sobre la página</span>{children('página')}
      </div>
      <div className="flex flex-col gap-2 rounded-md bg-surface-card p-3">
        <span className="text-caption text-text-secondary">Sobre una tarjeta</span>{children('tarjeta')}
      </div>
      <div className="flex flex-col gap-2 rounded-md bg-surface-raised p-3 border border-border-divider shadow-2">
        <span className="text-caption text-text-secondary">Sobre una superficie flotante</span>{children('flotante')}
      </div>
    </div>
  )
}
