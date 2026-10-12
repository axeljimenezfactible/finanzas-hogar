// Vista tabla (DESIGN §7.5.2) con TanStack Table 9: orden por columna y filtro de texto; agrupar por propiedad.
import { useMemo, useState } from 'react'
import { ArrowDown, ArrowUp } from 'lucide-react'
import { columnFilteringFeature, createFilteredRowModel, createSortedRowModel, filterFns, globalFilteringFeature, rowSortingFeature, sortFns, tableFeatures, useTable, type RowData, type SortingState } from '@tanstack/react-table'
import { textos } from '../textos/es-MX'
import { cn } from '../componentes/utils'
import { CeldaEditable } from './CeldaEditable'
import type { DefVista, Propiedad } from './tipos'
const caracteristicas = tableFeatures({ rowSortingFeature, columnFilteringFeature, globalFilteringFeature, sortedRowModel: createSortedRowModel(), filteredRowModel: createFilteredRowModel(), sortFns, filterFns })
export function Celda<T>({ fila, prop }: { fila: T; prop: Propiedad<T> }) {
  const contenido = prop.pintar ? prop.pintar(fila) : String(prop.valor(fila) ?? '')
  return prop.editar ? <CeldaEditable fila={fila} prop={prop}>{contenido}</CeldaEditable> : <>{contenido}</>
}
export function VistaTabla<T>({ def, texto, agruparPor }: { def: DefVista<T>; texto: string; agruparPor?: string }) {
  const [orden, setOrden] = useState<SortingState>([])
  const columnas = useMemo(() => def.propiedades.map((p) => ({ id: p.id, header: p.nombre, accessorFn: (f: RowData) => p.valor(f as T) ?? '' })), [def.propiedades])
  const tabla = useTable({ features: caracteristicas, columns: columnas, data: def.filas as RowData[], state: { sorting: orden, globalFilter: texto }, onSortingChange: setOrden, globalFilterFn: 'includesString', getRowId: (f: RowData) => String(def.idFila(f as T)) })
  const filas = tabla.getRowModel().rows
  const prop = def.propiedades.find((p) => p.id === agruparPor)
  const grupos = prop ? Object.entries(Object.groupBy(filas, (r) => (prop.grupo ?? ((f: T) => String(prop.valor(f) ?? textos.vistas.sinValor)))(r.original as T))) : [['', filas] as const]
  const props = def.propiedades
  return (
    <div className="overflow-x-auto">
      <table className="w-full border-collapse text-small">
        <thead>
          <tr className="border-b border-border-divider text-text-secondary">
            {tabla.getHeaderGroups()[0].headers.map((h) => {
              const p = props.find((x) => x.id === h.column.id)!
              const s = h.column.getIsSorted()
              return (
                <th key={h.id} scope="col" aria-sort={s === 'asc' ? 'ascending' : s === 'desc' ? 'descending' : 'none'} className={cn('px-2 py-1 font-medium', p.alinear === 'derecha' ? 'text-right' : 'text-left')}>
                  <button type="button" className="inline-flex h-9 items-center gap-1 coarse:h-11 coarse:min-w-11" onClick={h.column.getToggleSortingHandler()}>
                    {p.nombre}{s === 'asc' && <ArrowUp aria-label={textos.vistas.ascendente} className="size-5" />}{s === 'desc' && <ArrowDown aria-label={textos.vistas.descendente} className="size-5" />}
                  </button>
                </th>
              )
            })}
          </tr>
        </thead>
        {grupos.map(([g, rs]) => (
          <tbody key={g}>
            {prop && <tr><th colSpan={props.length} scope="rowgroup" className="bg-surface-hover-subtle px-2 py-1 text-left font-medium">{g} · {textos.vistas.registros(rs?.length ?? 0)}</th></tr>}
            {rs?.map((r) => (
              <tr key={r.id} className="border-b border-border-divider hover:bg-surface-hover-subtle">
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
