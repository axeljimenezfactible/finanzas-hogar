// Edición en grupo (U4 §0.2 punto 2 y §0.7; U1 §5.3 `Ctrl+/`): se elige la propiedad y el valor nuevo, se muestra la
// vista previa SUMADA «Antes → Después» sobre las filas seleccionadas y solo entonces se puede aplicar.
// `Enter` aplica a todas; `Esc` cancela sin cambiar ninguna. Sin vista previa calculada no aplica.
import { useEffect, useRef, useState, type KeyboardEvent } from 'react'
import { ArrowRight } from 'lucide-react'
import { textos } from '../textos/es-MX'
import { estilosControl } from '../componentes/Campo'
import { Boton } from '../componentes/Boton'
import { useAvisos } from '../componentes/Avisos'
import { ErrorApi } from '../api/errores'
import { useConexion } from '../api/conexion'
import type { Cambio, EditarVarias } from './tipos'
type Previa = Cambio[] | 'calculando' | 'error' | null
export function BloqueCambios({ cambios }: { cambios: Cambio[] }) {
  return (
    <dl aria-live="polite" data-prueba="vista-previa" className="rounded-md border border-border-divider bg-surface-raised p-2 text-small tabular-nums">
      <div className="flex gap-2 text-caption text-text-secondary"><span className="flex-1" /><span>{textos.vistas.antes}</span><ArrowRight aria-hidden className="size-5" /><span>{textos.vistas.despues}</span></div>
      {cambios.map((c) => <div key={c.etiqueta} className="flex items-center gap-2"><dt className="flex-1 text-text-secondary">{c.etiqueta}</dt><dd>{c.antes}</dd><ArrowRight aria-label={textos.vistas.cambiaA} className="size-5" /><dd className="font-medium">{c.despues}</dd></div>)}
    </dl>
  )
}
export function EdicionVarias<T>({ ev, filas, alTerminar }: { ev: EditarVarias<T>; filas: T[]; alTerminar: (aplicado: boolean) => void }) {
  const [propId, setPropId] = useState(ev.propiedades[0]?.id ?? '')
  const [valor, setValor] = useState('')
  const [previa, setPrevia] = useState<Previa>(null)
  const [guardando, setGuardando] = useState(false)
  const conexion = useConexion()
  const avisar = useAvisos()
  const turno = useRef(0)
  const primero = useRef<HTMLSelectElement & HTMLInputElement>(null)
  const prop = ev.propiedades.find((p) => p.id === propId)
  useEffect(() => { primero.current?.focus() }, [])
  const calcular = (p: string, v: string) => {
    const t = ++turno.current
    if (v === '') return setPrevia(null)
    setPrevia('calculando')
    ev.vistaPrevia(filas, p, v).then((c) => { if (t === turno.current) setPrevia(c) }, () => { if (t === turno.current) setPrevia('error') })
  }
  const listo = Array.isArray(previa)
  const motivo = !conexion ? textos.conexion.requiere : guardando ? textos.vistas.guardando : previa === 'calculando' ? textos.vistas.calculando : !listo ? textos.vistas.eligeValor : undefined
  const aplicar = async () => {
    if (motivo) return
    setGuardando(true)
    try {
      await ev.guardar(filas, propId, valor)
      avisar({ texto: textos.vistas.guardadoVarias(filas.length), rol: 'exito' })
      alTerminar(true)
    } catch (e) {
      avisar({ texto: e instanceof ErrorApi ? e.message : textos.errores.generico, rol: 'error' })
      setGuardando(false)
    }
  }
  const teclas = (e: KeyboardEvent) => {
    if (e.key === 'Escape') { e.preventDefault(); e.stopPropagation(); alTerminar(false) }
    else if (e.key === 'Enter' && !(e.target instanceof HTMLButtonElement)) { e.preventDefault(); void aplicar() }
  }
  const controlValor = prop?.tipo === 'opcion'
    ? <select ref={ev.propiedades.length === 1 ? primero : undefined} aria-label={textos.vistas.valorNuevo} value={valor} className={estilosControl} onChange={(e) => { setValor(e.target.value); calcular(propId, e.target.value) }}>
        <option value="">{textos.vistas.eligeOpcion}</option>
        {prop.opciones?.map((o) => <option key={o.valor} value={o.valor}>{o.texto}</option>)}
      </select>
    : <input ref={ev.propiedades.length === 1 ? primero : undefined} aria-label={textos.vistas.valorNuevo} value={valor} className={estilosControl} type={prop?.tipo === 'fecha' ? 'date' : 'text'} inputMode={prop?.tipo === 'dinero' || prop?.tipo === 'numero' ? 'decimal' : undefined} onChange={(e) => { setValor(e.target.value); calcular(propId, e.target.value) }} />
  return (
    <section role="region" aria-label={textos.vistas.edicionGrupo(filas.length)} data-prueba="edicion-varias" onKeyDown={teclas} className="flex flex-col gap-2 rounded-md border border-border-divider bg-surface-card p-3">
      <h3 className="text-small font-medium text-text-default">{textos.vistas.edicionGrupo(filas.length)}</h3>
      <div className="flex flex-wrap items-center gap-2">
        {ev.propiedades.length > 1 && (
          <select ref={primero} aria-label={textos.vistas.propiedad} value={propId} className={estilosControl + ' w-auto'} onChange={(e) => { setPropId(e.target.value); setValor(''); calcular(e.target.value, '') }}>
            {ev.propiedades.map((p) => <option key={p.id} value={p.id}>{p.nombre}</option>)}
          </select>
        )}
        {ev.propiedades.length === 1 && <span className="text-small text-text-secondary">{prop?.nombre}</span>}
        <div className="min-w-40 flex-1">{controlValor}</div>
      </div>
      {previa === 'calculando' && <p role="status" className="text-small text-text-secondary">{textos.vistas.calculando}</p>}
      {previa === 'error' && <p role="status" className="text-small text-role-error">{textos.vistas.noSeCalculo} · <Boton variante="enlace" onClick={() => calcular(propId, valor)}>{textos.app.reintentar}</Boton></p>}
      {Array.isArray(previa) && (previa.length ? <BloqueCambios cambios={previa} /> : <p role="status" className="text-small text-text-secondary">{textos.vistas.sinEfecto}</p>)}
      <p className="text-caption text-text-secondary">{textos.vistas.ayudaGrupo}</p>
      <div className="flex flex-wrap gap-2">
        <Boton variante="primario" deshabilitadoPor={motivo} onClick={() => void aplicar()}>{textos.vistas.aplicarA(filas.length)}</Boton>
        <Boton variante="secundario" onClick={() => alTerminar(false)}>{textos.app.cancelar}</Boton>
      </div>
    </section>
  )
}
