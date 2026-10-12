// Vista galería (DESIGN §7.5.5): tarjetas con el título y las primeras propiedades.
import type { DefVista } from './tipos'
export function VistaGaleria<T>({ def, filas }: { def: DefVista<T>; filas: T[] }) {
  return (
    <ul className="grid grid-cols-1 gap-3 phone:grid-cols-2 desktop:grid-cols-3">
      {filas.map((f) => (
        <li key={def.idFila(f)}>
          <button type="button" onClick={() => def.alAbrir?.(f)} className="flex w-full flex-col gap-1 rounded-lg border border-border-divider bg-surface-card p-3 text-left hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-border-focus">
            <span className="text-body font-medium">{def.titulo(f)}</span>
            {def.propiedades.slice(1, 4).map((p) => <span key={p.id} className="text-small text-text-secondary tabular-nums"><span className="sr-only">{p.nombre}: </span>{p.pintar ? p.pintar(f) : String(p.valor(f) ?? '')}</span>)}
          </button>
        </li>
      ))}
    </ul>
  )
}
