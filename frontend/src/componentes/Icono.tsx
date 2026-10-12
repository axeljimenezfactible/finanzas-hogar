// Iconos Lucide (DESIGN §5.2): solo dos tamaños, 20 px (`size.icon.sm`) y 24 px (`size.icon.md`), y trazo constante
// de 1.5 px dibujados (`icon.stroke`) con el ancho de trazo absoluto del componente: stroke-width 1.8 a 20 px y 1.5 a
// 24 px en unidades del viewBox de 24 (`icon.stroke.user-units.sm` · `.md`).
// ProveedorIconos fija esos valores por omisión para toda la app (20 px). Un icono de 24 px (solo icono, barras,
// estados vacíos) lleva className="size-6" y size={ICONO_MD}, para que el trazo se calcule sobre 24 px.
import type { ReactNode } from 'react'
import { LucideProvider } from 'lucide-react'
export const ICONO_SM = 20
export const ICONO_MD = 24
export const TRAZO_ICONO = 1.5
export function ProveedorIconos({ children }: { children: ReactNode }) {
  return <LucideProvider size={ICONO_SM} strokeWidth={TRAZO_ICONO} absoluteStrokeWidth>{children}</LucideProvider>
}
