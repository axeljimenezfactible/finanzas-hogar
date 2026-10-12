// Vista tablero (DESIGN §7.5.3) con dnd-kit: arrastrar una tarjeta a otra columna cambia su propiedad.
// Con teclado: Espacio toma la tarjeta, flechas la mueven, Espacio la suelta (KeyboardSensor).
import { DndContext, KeyboardSensor, PointerSensor, TouchSensor, useDraggable, useDroppable, useSensor, useSensors, type DragEndEvent } from '@dnd-kit/core'
import { CSS } from '@dnd-kit/utilities'
import type { ReactNode } from 'react'
import { textos } from '../textos/es-MX'
import { useAvisos } from '../componentes/Avisos'
import { ErrorApi } from '../api/errores'
import type { DefVista } from './tipos'
function Columna({ id, titulo, n, children }: { id: string; titulo: string; n: number; children: ReactNode }) {
  const { setNodeRef, isOver } = useDroppable({ id })
  return (
    <section ref={setNodeRef} aria-label={titulo} className={`flex w-72 shrink-0 flex-col gap-2 rounded-lg p-2 ${isOver ? 'bg-surface-selected' : 'bg-surface-hover-subtle'}`}>
      <h3 className="px-1 text-small font-medium text-text-secondary">{titulo} · {n}</h3>
      {children}
    </section>
  )
}
function Tarjeta({ id, children, alAbrir }: { id: string; children: ReactNode; alAbrir?: () => void }) {
  const { attributes, listeners, setNodeRef, transform } = useDraggable({ id })
  return (
    <div ref={setNodeRef} style={{ transform: CSS.Translate.toString(transform) }} {...attributes} {...listeners} onDoubleClick={alAbrir}
      className="cursor-grab rounded-md border border-border-divider bg-surface-card p-3 text-small shadow-1 focus-visible:outline-2 focus-visible:outline-border-focus">
      {children}
    </div>
  )
}
export function VistaTablero<T>({ def, filas }: { def: DefVista<T>; filas: T[] }) {
  const avisar = useAvisos()
  const sensores = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 6 } }), useSensor(TouchSensor, { activationConstraint: { delay: 250, tolerance: 6 } }), useSensor(KeyboardSensor))
  const cfg = def.columnasTablero!
  const prop = def.propiedades.find((p) => p.id === cfg.propiedad)!
  const grupoDe = (f: T) => (prop.grupo ?? ((x: T) => String(prop.valor(x) ?? textos.vistas.sinValor)))(f)
  const alSoltar = async (e: DragEndEvent) => {
    const fila = filas.find((f) => String(def.idFila(f)) === String(e.active.id))
    const destino = e.over?.id as string | undefined
    if (!fila || !destino || destino === grupoDe(fila) || !cfg.moverA) return
    try { await cfg.moverA(fila, destino) } catch (err) { avisar({ texto: err instanceof ErrorApi ? err.message : textos.errores.generico, rol: 'error' }) }
  }
  return (
    <DndContext sensors={sensores} onDragEnd={alSoltar}>
      <div className="flex gap-3 overflow-x-auto pb-2">
        {cfg.columnas.map((c) => {
          const enCol = filas.filter((f) => grupoDe(f) === c)
          return (
            <Columna key={c} id={c} titulo={c} n={enCol.length}>
              {enCol.map((f) => <Tarjeta key={def.idFila(f)} id={String(def.idFila(f))} alAbrir={def.alAbrir && (() => def.alAbrir!(f))}>{def.titulo(f)}</Tarjeta>)}
            </Columna>
          )
        })}
      </div>
    </DndContext>
  )
}
