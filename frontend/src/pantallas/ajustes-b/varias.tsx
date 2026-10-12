// Edición en grupo con `Ctrl+/` (U4 §0.2 punto 2 y §0.7; U1 §5.3) para listas de esta área que no son una VistaBD
// (P-26.09: las alertas son tarjetas). Mismo comportamiento que `VistaBD` con `def.editarVarias` (CONVENCIONES-frontend
// §7): casilla por fila y «Seleccionar todas», «n seleccionadas · Editar seleccionadas (Ctrl+/) · Quitar selección», y el
// bloque `EdicionVarias` de la base con la vista previa sumada («Antes → Después»): `Enter` aplica a todas y `Esc`
// cancela sin cambiar ninguna; el foco regresa a donde estaba.
import { useEffect, useRef, useState, type KeyboardEvent, type ReactNode } from 'react'
import { Boton } from '../../componentes/Boton'
import { useAvisos } from '../../componentes/Avisos'
import { textos } from '../../textos/es-MX'
import { EdicionVarias } from '../../vistas-bd/EdicionVarias'
import type { EditarVarias } from '../../vistas-bd/tipos'

/** Casilla de 20 px dentro de un área de toque de 36 px (fino) o 44 px (grueso), como la de la tabla de la base. */
export function Casilla({ etiqueta, marcada, mixta = false, alCambiar, texto }: { etiqueta: string; marcada: boolean; mixta?: boolean; alCambiar: (v: boolean) => void; texto?: string }) {
  const ref = useRef<HTMLInputElement>(null)
  useEffect(() => { if (ref.current) ref.current.indeterminate = mixta }, [mixta])
  return (
    <label className="inline-flex min-h-9 min-w-9 cursor-pointer items-center justify-center gap-2 text-small text-text-default coarse:min-h-11 coarse:min-w-11">
      <input ref={ref} type="checkbox" aria-label={etiqueta} checked={marcada} onChange={(e) => alCambiar(e.target.checked)} className="size-5 cursor-pointer accent-action-primary-bg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus" />
      {texto && <span aria-hidden>{texto}</span>}
    </label>
  )
}

export function GrupoVarias<T>({ filas, idFila, nombre, ev, children }: {
  filas: T[]; idFila: (f: T) => string | number; nombre: (f: T) => string; ev: EditarVarias<T>
  /** Pinta la lista; `casilla(fila)` da la casilla de selección de esa fila. */
  children: (casilla: (f: T) => ReactNode) => ReactNode
}) {
  const [sel, setSel] = useState<Set<string>>(() => new Set())
  const [enGrupo, setEnGrupo] = useState(false)
  const avisar = useAvisos()
  const raiz = useRef<HTMLDivElement>(null)
  const regreso = useRef<HTMLElement | null>(null)
  const ids = filas.map((f) => String(idFila(f)))
  const seleccionadas = filas.filter((f) => sel.has(String(idFila(f))))
  const abrir = () => {
    if (!seleccionadas.length) return avisar({ texto: textos.vistas.eligeFilas, rol: 'info' })
    regreso.current = document.activeElement as HTMLElement | null
    setEnGrupo(true)
  }
  const cerrar = (aplicado: boolean) => {
    setEnGrupo(false)
    if (aplicado) setSel(new Set())
    const r = regreso.current
    requestAnimationFrame(() => (r?.isConnected ? r : raiz.current?.querySelector<HTMLElement>('input[type=checkbox]'))?.focus())
  }
  const teclas = (e: KeyboardEvent) => { if ((e.ctrlKey || e.metaKey) && e.key === '/') { e.preventDefault(); abrir() } }
  const alternar = (id: string, v: boolean) => { const s = new Set(sel); if (v) s.add(id); else s.delete(id); setSel(s) }
  const casilla = (f: T) => { const id = String(idFila(f)); return <Casilla etiqueta={textos.vistas.seleccionar(nombre(f))} marcada={sel.has(id)} alCambiar={(v) => alternar(id, v)} /> }
  const n = seleccionadas.length
  return (
    <div ref={raiz} className="flex flex-col gap-3" onKeyDown={teclas}>
      <div className="flex flex-wrap items-center gap-2">
        <Casilla etiqueta={textos.vistas.seleccionarTodas} texto={textos.vistas.seleccionarTodas} marcada={n > 0 && n === ids.length} mixta={n > 0 && n < ids.length} alCambiar={(v) => setSel(v ? new Set(ids) : new Set())} />
        {n > 0 && (<>
          <span aria-live="polite" className="text-small text-text-secondary">{textos.vistas.seleccionadas(n)}</span>
          <Boton variante="secundario" aria-keyshortcuts="Control+/" onClick={abrir}>{textos.vistas.editarSeleccionadas}<kbd className="hidden rounded-sm border border-border-divider px-1 text-caption text-text-secondary desktop:inline">{textos.vistas.atajoEditar}</kbd></Boton>
          <Boton variante="fantasma" onClick={() => { setSel(new Set()); setEnGrupo(false) }}>{textos.vistas.quitarSeleccion}</Boton>
        </>)}
      </div>
      {enGrupo && n > 0 && <EdicionVarias key={[...sel].join(',')} ev={ev} filas={seleccionadas} alTerminar={cerrar} />}
      {children(casilla)}
    </div>
  )
}
