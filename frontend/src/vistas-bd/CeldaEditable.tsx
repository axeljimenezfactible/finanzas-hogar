// Edición en el lugar con vista previa «Antes → Después» (U4 §0.2): Enter o salir guarda; Esc descarta.
// Si la vista previa no se puede calcular, dice «No se pudo calcular el efecto · Reintentar» y no guarda.
import { useEffect, useRef, useState, type ReactNode } from 'react'
import { ArrowRight } from 'lucide-react'
import { textos } from '../textos/es-MX'
import { estilosControl } from '../componentes/Campo'
import { useAvisos } from '../componentes/Avisos'
import { ErrorApi } from '../api/errores'
import type { Cambio, Propiedad } from './tipos'
export function CeldaEditable<T>({ fila, prop, children }: { fila: T; prop: Propiedad<T>; children: ReactNode }) {
  const ed = prop.editar!
  const [editando, setEditando] = useState(false)
  const [valor, setValor] = useState('')
  const [previa, setPrevia] = useState<Cambio[] | 'error' | null>(null)
  const [guardando, setGuardando] = useState(false)
  const avisar = useAvisos()
  const ref = useRef<HTMLInputElement & HTMLSelectElement>(null)
  useEffect(() => { if (editando) ref.current?.focus() }, [editando])
  const inicial = ed.textoInicial(fila)
  const calcular = (v: string) => {
    if (!ed.vistaPrevia || v === inicial) return setPrevia(null)
    ed.vistaPrevia(fila, v).then(setPrevia, () => setPrevia('error'))
  }
  const cancelar = () => { setEditando(false); setPrevia(null) }
  const guardar = async () => {
    if (guardando) return
    if (valor === inicial) return cancelar()
    if (ed.vistaPrevia && (previa === null || previa === 'error')) return // la vista previa debe verse antes de guardar
    setGuardando(true)
    try {
      await ed.guardar(fila, valor)
      avisar({ texto: textos.vistas.guardadoEnLugar, rol: 'exito' })
      cancelar()
    } catch (e) {
      avisar({ texto: e instanceof ErrorApi ? e.message : textos.errores.generico, rol: 'error' })
      cancelar() // regresa al valor anterior (U1 §2.4)
    } finally { setGuardando(false) }
  }
  if (!editando) return <button type="button" className="w-full min-h-9 text-inherit rounded-sm px-1 text-left coarse:min-h-11 coarse:min-w-11 hover:bg-surface-hover focus-visible:outline-2 focus-visible:outline-border-focus" onClick={() => { setValor(inicial); setEditando(true) }}>{children}</button>
  const comunes = {
    ref, value: valor, 'aria-label': prop.nombre, className: estilosControl,
    onKeyDown: (e: React.KeyboardEvent) => { if (e.key === 'Enter') { e.preventDefault(); void guardar() } if (e.key === 'Escape') { e.preventDefault(); cancelar() } },
    onBlur: () => void guardar(),
  }
  return (
    <div className="flex flex-col gap-1">
      {ed.tipo === 'opcion'
        ? <select {...comunes} onChange={(e) => { setValor(e.target.value); calcular(e.target.value) }}>{ed.opciones?.map((o) => <option key={o.valor} value={o.valor}>{o.texto}</option>)}</select>
        : <input {...comunes} type={ed.tipo === 'fecha' ? 'date' : 'text'} inputMode={ed.tipo === 'dinero' || ed.tipo === 'numero' ? 'decimal' : undefined} onChange={(e) => { setValor(e.target.value); calcular(e.target.value) }} />}
      {previa === 'error' && <p role="status" className="text-small text-role-error">No se pudo calcular el efecto · <button type="button" className="text-text-link coarse:inline-flex coarse:min-h-11 coarse:items-center" onMouseDown={(e) => e.preventDefault()} onClick={() => calcular(valor)}>{textos.app.reintentar}</button></p>}
      {Array.isArray(previa) && previa.length > 0 && (
        <dl aria-live="polite" className="rounded-md border border-border-divider bg-surface-raised p-2 text-small tabular-nums">
          <div className="flex gap-2 text-caption text-text-secondary"><span className="flex-1" /><span>{textos.vistas.antes}</span><ArrowRight aria-hidden className="size-5" /><span>{textos.vistas.despues}</span></div>
          {previa.map((c) => <div key={c.etiqueta} className="flex items-center gap-2"><dt className="flex-1 text-text-secondary">{c.etiqueta}</dt><dd>{c.antes}</dd><ArrowRight aria-label="cambia a" className="size-5" /><dd className="font-medium">{c.despues}</dd></div>)}
        </dl>
      )}
    </div>
  )
}
