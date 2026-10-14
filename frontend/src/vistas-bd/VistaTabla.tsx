// Vista tabla (DESIGN §7.5.2) con TanStack Table 9: orden por columna y filtro de texto; agrupar por propiedad.
import { isValidElement, useEffect, useMemo, useRef, useState } from 'react'
import { ArrowDown, ArrowUp } from 'lucide-react'
import { columnFilteringFeature, createFilteredRowModel, createSortedRowModel, filterFns, globalFilteringFeature, rowSortingFeature, sortFns, tableFeatures, useTable, type RowData, type SortingState } from '@tanstack/react-table'
import { textos } from '../textos/es-MX'
import { cn } from '../componentes/utils'
import { CeldaEditable } from './CeldaEditable'
import type { DefVista, Propiedad } from './tipos'
const caracteristicas = tableFeatures({ rowSortingFeature, columnFilteringFeature, globalFilteringFeature, sortedRowModel: createSortedRowModel(), filteredRowModel: createFilteredRowModel(), sortFns, filterFns })
export function Celda<T>({ fila, prop }: { fila: T; prop: Propiedad<T> }) {
  const crudo = prop.pintar ? prop.pintar(fila) : String(prop.valor(fila) ?? '')
  const contenido = crudo === '' || crudo === null || crudo === undefined ? textos.vistas.sinValor : crudo // R3-02: nunca un botón sin nombre
  return prop.editar ? <CeldaEditable fila={fila} prop={prop}>{contenido}</CeldaEditable> : <>{contenido}</>
}
/** Selección de filas (solo con `def.editarVarias`): ids como texto. */
export interface Seleccion { ids: Set<string>; cambiar: (s: Set<string>) => void }
/** Casilla de 20 px dentro de un área de toque de 36 px (fino) o 44 px (grueso): la etiqueta la marca. */
function Casilla({ etiqueta, marcada, mixta = false, alCambiar }: { etiqueta: string; marcada: boolean; mixta?: boolean; alCambiar: (v: boolean) => void }) {
  const ref = useRef<HTMLInputElement>(null)
  useEffect(() => { if (ref.current) ref.current.indeterminate = mixta }, [mixta])
  return (
    <label className="inline-flex size-9 cursor-pointer items-center justify-center coarse:size-11">
      <input ref={ref} type="checkbox" aria-label={etiqueta} checked={marcada} onChange={(e) => alCambiar(e.target.checked)} className="size-5 cursor-pointer accent-action-primary-bg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus" />
    </label>
  )
}
/** Nombre de la fila para su casilla: el título si es texto; si no, la primera propiedad pintada o su valor. */
function nombreFila<T>(def: DefVista<T>, f: T): string {
  const t = def.titulo(f)
  if (typeof t === 'string' || typeof t === 'number') return String(t)
  const p = def.propiedades[0]
  const pintado = p?.pintar?.(f)
  if (typeof pintado === 'string' || typeof pintado === 'number') return String(pintado)
  if (isValidElement(pintado)) return String(p.valor(f) ?? def.idFila(f))
  return String(p?.valor(f) ?? def.idFila(f))
}
export function VistaTabla<T>({ def, texto, agruparPor, seleccion }: { def: DefVista<T>; texto: string; agruparPor?: string; seleccion?: Seleccion }) {
  const [orden, setOrden] = useState<SortingState>([])
  const columnas = useMemo(() => def.propiedades.map((p) => ({ id: p.id, header: p.nombre, accessorFn: (f: RowData) => p.valor(f as T) ?? '' })), [def.propiedades])
  const tabla = useTable({ features: caracteristicas, columns: columnas, data: def.filas as RowData[], state: { sorting: orden, globalFilter: texto }, onSortingChange: setOrden, globalFilterFn: 'includesString', getRowId: (f: RowData) => String(def.idFila(f as T)) })
  const filas = tabla.getRowModel().rows
  const prop = def.propiedades.find((p) => p.id === agruparPor)
  const grupos = prop ? Object.entries(Object.groupBy(filas, (r) => (prop.grupo ?? ((f: T) => String(prop.valor(f) ?? textos.vistas.sinValor)))(r.original as T))) : [['', filas] as const]
  const props = def.propiedades
  const visibles = filas.map((r) => r.id)
  const marcadas = seleccion ? visibles.filter((id) => seleccion.ids.has(id)).length : 0
  const alternar = (id: string, v: boolean) => { if (!seleccion) return; const s = new Set(seleccion.ids); if (v) s.add(id); else s.delete(id); seleccion.cambiar(s) }
  const todas = (v: boolean) => { if (!seleccion) return; const s = new Set(seleccion.ids); for (const id of visibles) { if (v) s.add(id); else s.delete(id) } seleccion.cambiar(s) }
  return (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse text-small">
        <thead>
          <tr className="border-b border-border-divider text-text-secondary">
            {seleccion && <th scope="col" className="w-11 px-1"><Casilla etiqueta={textos.vistas.seleccionarTodas} marcada={visibles.length > 0 && marcadas === visibles.length} mixta={marcadas > 0 && marcadas < visibles.length} alCambiar={todas} /></th>}
            {tabla.getHeaderGroups()[0].headers.map((h) => {
              const p = props.find((x) => x.id === h.column.id)!
              const s = h.column.getIsSorted()
              return (
                <th key={h.id} scope="col" aria-sort={s === 'asc' ? 'ascending' : s === 'desc' ? 'descending' : 'none'} className={cn('px-2 py-1 font-medium', p.alinear === 'derecha' ? 'text-right' : 'text-left')}>
                  {p.nombre ? <button type="button" className="inline-flex h-9 items-center gap-1 coarse:h-11 coarse:min-w-11" onClick={h.column.getToggleSortingHandler()}>
                    {p.nombre}{s === 'asc' && <ArrowUp aria-label={textos.vistas.ascendente} className="size-5" />}{s === 'desc' && <ArrowDown aria-label={textos.vistas.descendente} className="size-5" />}
                  </button> : <span className="sr-only">{textos.vistas.columnaAcciones}</span>}
                </th>
              )
            })}
          </tr>
        </thead>
        {grupos.map(([g, rs]) => (
          <tbody key={g}>
            {prop && <tr><th colSpan={props.length + (seleccion ? 1 : 0)} scope="rowgroup" className="bg-surface-hover-subtle px-2 py-1 text-left font-medium">{g} · {textos.vistas.registros(rs?.length ?? 0)}</th></tr>}
            {rs?.map((r) => (
              <tr key={r.id} data-seleccionada={seleccion?.ids.has(r.id) || undefined} className="border-b border-border-divider hover:bg-surface-hover-subtle data-seleccionada:bg-surface-selected">
                {seleccion && <td className="px-1"><Casilla etiqueta={textos.vistas.seleccionar(nombreFila(def, r.original as T))} marcada={seleccion.ids.has(r.id)} alCambiar={(v) => alternar(r.id, v)} /></td>}
                {props.map((p, i) => (
                  <td key={p.id} className={cn('px-2 py-1', p.alinear === 'derecha' && 'text-right tabular-nums')}>
                    {i === 0 && def.alAbrir && !p.editar ? <button type="button" className="inline-flex items-center text-left font-medium hover:underline coarse:min-h-11 coarse:min-w-11" onClick={() => def.alAbrir!(r.original as T)}><Celda fila={r.original as T} prop={p} /></button> : <Celda fila={r.original as T} prop={p} />}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        ))}
      </table>
    </div>
  )
}
