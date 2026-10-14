// Lista que se reordena arrastrando (p. ej. P-26.04 «orden» de categorías) con @dnd-kit/sortable; también con teclado.
import { DndContext, KeyboardSensor, PointerSensor, TouchSensor, closestCenter, useSensor, useSensors, type DragEndEvent } from '@dnd-kit/core'
import { SortableContext, arrayMove, sortableKeyboardCoordinates, useSortable, verticalListSortingStrategy } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { GripVertical } from 'lucide-react'
import { ICONO_MD } from '../componentes/Icono'
import type { ReactNode } from 'react'
function Renglon({ id, children }: { id: string; children: ReactNode }) {
  const { attributes, listeners, setNodeRef, transform, transition } = useSortable({ id })
  return (
    <li ref={setNodeRef} style={{ transform: CSS.Transform.toString(transform), transition }} className="flex items-center gap-2 rounded-md border border-border-divider bg-surface-card p-2">
      <button type="button" {...attributes} {...listeners} className="inline-flex size-9 cursor-grab items-center justify-center coarse:size-11"><GripVertical aria-hidden className="size-6" size={ICONO_MD} /></button>
      {children}
    </li>
  )
}
export function ListaOrdenable<T>({ items, idDe, pintar, alReordenar }: { items: T[]; idDe: (x: T) => string; pintar: (x: T) => ReactNode; alReordenar: (nuevo: T[]) => void }) {
  const sensores = useSensors(useSensor(PointerSensor), useSensor(TouchSensor, { activationConstraint: { delay: 250, tolerance: 6 } }), useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }))
  const alSoltar = (e: DragEndEvent) => {
    if (!e.over || e.active.id === e.over.id) return
    const de = items.findIndex((x) => idDe(x) === e.active.id), a = items.findIndex((x) => idDe(x) === e.over!.id)
    alReordenar(arrayMove(items, de, a))
  }
  return (
    <DndContext sensors={sensores} collisionDetection={closestCenter} onDragEnd={alSoltar}>
      <SortableContext items={items.map(idDe)} strategy={verticalListSortingStrategy}>
        <ul className="flex flex-col gap-2">{items.map((x) => <Renglon key={idDe(x)} id={idDe(x)}>{pintar(x)}</Renglon>)}</ul>
      </SortableContext>
    </DndContext>
  )
}
