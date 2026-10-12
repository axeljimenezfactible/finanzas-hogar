// P-03 Registrar gasto (hoja global) con sus pestañas Gasto · Ingreso · Peso · Ticket, y P-22 Registrar peso.
// F-1 en 3 toques (K10): «+» abre con el monto enfocado y los valores por defecto (N-052); chip de categoría; Guardar.
// Sin conexión (K9): Guardar deja el gasto o el peso en la cola del teléfono con su clave y su captured_at.
// U1 §3.5: la pestaña «Ingreso» registra el ingreso extra en la misma hoja (N-102) y la forma de pago «A crédito» lista
// las líneas activas (N-053); las dos piezas son del área Gastos y Presupuesto (pantallas/gastos-presupuesto/hoja-registro).
import { useMemo, useRef, useState, type FormEvent } from 'react'
import { textos } from '../textos/es-MX'
import { nuevaClave, pedir } from '../api/cliente'
import { usarConsulta } from '../api/consulta'
import { ErrorApi, esSinConexion } from '../api/errores'
import { aCentavos, aGramos, fecha, hoyISO, kilos } from '../api/formatos'
import { hayConexion, useConexion } from '../api/conexion'
import { encolarGasto, encolarPeso, quitarDeCola, type CuerpoGasto } from '../cola/cola'
import { Hoja } from '../componentes/Hoja'
import { Boton } from '../componentes/Boton'
import { Campo, estilosControl } from '../componentes/Campo'
import { ChipSeleccionable, MarcaPersona } from '../componentes/Chip'
import { Franja } from '../componentes/Estados'
import { useAvisos } from '../componentes/Avisos'
import { cn } from '../componentes/utils'
import { navegar } from '../app/enrutador'
import { rutaDe } from '../app/rutas'
import { useSesion } from '../app/sesion'
import { FormIngreso, LineaCredito } from '../pantallas/gastos-presupuesto/hoja-registro'
import { avisarCambio, cambiar, useShell, type TipoRegistro } from './estado'
const TIPOS: TipoRegistro[] = ['gasto', 'ingreso', 'peso', 'ticket']
type Persona = { id: number; code: string; display_name: string; participation: string }
export function HojaRegistro() {
  const { registro } = useShell()
  const ref = useRef<HTMLInputElement>(null)
  const cerrar = () => cambiar({ registro: null })
  return (
    <Hoja abierta={registro !== null} alCambiar={(a) => !a && cerrar()} titulo={textos.registro.titulo} alAbrirFoco={(e) => { e.preventDefault(); ref.current?.focus() }}>
      <div role="tablist" aria-label={textos.registro.titulo} className="mb-3 grid grid-cols-4 gap-1 rounded-md bg-surface-hover-subtle p-1">
        {TIPOS.map((t) => (
          <button key={t} role="tab" type="button" aria-selected={registro === t} onClick={() => cambiar({ registro: t })}
            className={cn('h-9 rounded-sm text-small coarse:h-11', registro === t ? 'bg-surface-card font-medium shadow-1' : 'text-text-secondary hover:bg-surface-hover')}>{textos.registro.tipos[t]}</button>
        ))}
      </div>
      {registro === 'gasto' && <FormGasto refMonto={ref} alTerminar={cerrar} />}
      {registro === 'peso' && <FormPeso refPeso={ref} alTerminar={cerrar} />}
      {registro === 'ingreso' && <FormIngreso refMonto={ref} alTerminar={cerrar} />}
      {registro === 'ticket' && <IrA texto={textos.registro.ticketIr} boton={textos.registro.irTicket} a={rutaDe('P-05')} />}
    </Hoja>
  )
}
function IrA({ texto, boton, a }: { texto: string; boton: string; a: string }) {
  const conectado = useConexion()
  if (!conectado) return <Franja rol="sinConexion">{textos.conexion.requiere}</Franja>
  return <div className="flex flex-col gap-3"><p className="text-body">{texto}</p><Boton variante="primario" onClick={() => { cambiar({ registro: null }); navegar(a) }}>{boton}</Boton></div>
}
/** Deja en el teléfono la última carga de lo que P-03 y P-22 necesitan sin conexión (categorías y personas). */
export function precargarRegistro() {
  void pedir('get', '/api/v1/categories', { consulta: { status: 'active', limit: 200 } }).catch(() => undefined)
  void pedir('get', '/api/v1/persons').catch(() => undefined)
}
function usarPersonas() {
  const c = usarConsulta('get', '/api/v1/persons')
  return (c.datos?.items ?? []) as Persona[]
}
function FormGasto({ refMonto, alTerminar }: { refMonto: React.RefObject<HTMLInputElement | null>; alTerminar: () => void }) {
  const { sesion } = useSesion()
  const avisar = useAvisos()
  const cats = usarConsulta('get', '/api/v1/categories', { consulta: { status: 'active', limit: 200 } })
  const personas = usarPersonas()
  const [monto, setMonto] = useState('')
  const [categoria, setCategoria] = useState<number | null>(null)
  const [todas, setTodas] = useState(false)
  const [mas, setMas] = useState(false)
  const [dia, setDia] = useState(hoyISO())
  const [persona, setPersona] = useState<number | null>(sesion?.person_id ?? null)
  const [forma, setForma] = useState<CuerpoGasto['payment_method']>('cash')
  const [descripcion, setDescripcion] = useState('')
  const [nota, setNota] = useState('')
  const [errores, setErrores] = useState<Record<string, string>>({})
  const [guardando, setGuardando] = useState(false)
  const activas = useMemo(() => [...(cats.datos?.items ?? [])].filter((c) => c.is_active).sort((a, b) => a.sort_order - b.sort_order), [cats.datos])
  // N-050: las 6 con más gastos en 90 días; empates y relleno por el orden de P-26.04.
  const sugeridas = useMemo(() => [...activas].sort((a, b) => (b.recent_usage_count ?? 0) - (a.recent_usage_count ?? 0) || a.sort_order - b.sort_order).slice(0, 6), [activas])
  const visibles = todas ? activas : sugeridas
  const alGuardar = async (e: FormEvent) => {
    e.preventDefault()
    const centavos = aCentavos(monto)
    const errs: Record<string, string> = {}
    if (!centavos || centavos <= 0) errs.amount_cents = textos.registro.montoInvalido
    if (!categoria) errs.category_id = textos.registro.elegirCategoria
    setErrores(errs)
    if (Object.keys(errs).length || guardando) return
    if (forma === 'credit') return
    const cuerpo: CuerpoGasto = { kind: dia > hoyISO() ? 'upcoming' : 'variable', category_id: categoria!, amount_cents: centavos!, expense_date: dia, payment_method: forma, ...(persona ? { person_id: persona } : {}), ...(descripcion.trim() ? { description: descripcion.trim() } : {}), ...(nota.trim() ? { notes: nota.trim() } : {}) }
    const clave = nuevaClave()
    const enCola = async () => {
      await encolarGasto(cuerpo, clave)
      avisar({ texto: textos.registro.gastoEnCola, rol: 'exito', accion: { texto: textos.app.deshacer, alHacer: () => quitarDeCola(clave) } })
    }
    setGuardando(true)
    try {
      if (!hayConexion()) await enCola()
      else {
        const g = await pedir('post', '/api/v1/expenses', { cuerpo, idempotencia: clave })
        avisar({ texto: textos.registro.gastoGuardado, rol: 'exito', accion: { texto: textos.app.deshacer, alHacer: async () => { await pedir('delete', '/api/v1/expenses/{expense_id}', { ruta: { expense_id: g.id } }); avisarCambio('gasto') } } })
        avisarCambio('gasto')
      }
      alTerminar()
    } catch (err) {
      if (esSinConexion(err)) { await enCola(); alTerminar() }
      else if (err instanceof ErrorApi) setErrores(Object.fromEntries([['_', err.message], ...err.campos.map((c) => [c.field, c.message])]))
      else setErrores({ _: textos.errores.generico })
    } finally { setGuardando(false) }
  }
  return (
    <form onSubmit={alGuardar} className="flex flex-col gap-4" noValidate onKeyDown={(e) => { if (e.key === 'Enter' && e.ctrlKey) void alGuardar(e) }}>
      <Campo ref={refMonto} etiqueta={textos.registro.monto} inputMode="decimal" autoComplete="off" enterKeyHint="done" placeholder="$0.00" value={monto} onChange={(e) => setMonto(e.target.value)} error={errores.amount_cents} className="[&_input]:h-12 [&_input]:text-title-2 [&_input]:tabular-nums" />
      <fieldset className="flex flex-col gap-2">
        <legend className="mb-1 text-small font-medium text-text-secondary">{textos.registro.categoria}</legend>
        {!cats.cargando && activas.length === 0 && <p className="text-small text-text-secondary">{textos.registro.sinCategorias}</p>}
        <div className="flex flex-wrap gap-2">
          {visibles.map((c) => <ChipSeleccionable key={c.id} seleccionado={categoria === c.id} onClick={() => setCategoria(c.id)}>{c.name}</ChipSeleccionable>)}
          {!todas && activas.length > 6 && <ChipSeleccionable seleccionado={false} onClick={() => setTodas(true)}>{textos.app.mas}</ChipSeleccionable>}
        </div>
        {errores.category_id && <p className="text-small text-role-error">{errores.category_id}</p>}
      </fieldset>
      <Boton variante="enlace" className="self-start" aria-expanded={mas} onClick={() => setMas(!mas)}>{mas ? textos.app.menosOpciones : textos.app.masOpciones}</Boton>
      {mas && (
        <div className="grid grid-cols-1 gap-3 phone:grid-cols-2">
          <Campo etiqueta={textos.registro.fecha} type="date" value={dia} onChange={(e) => setDia(e.target.value)} ayuda={fecha(dia)} />
          <label className="flex flex-col gap-1 text-small font-medium text-text-secondary">{textos.registro.persona}
            <select className={estilosControl} value={persona ?? ''} onChange={(e) => setPersona(e.target.value ? Number(e.target.value) : null)}>
              {personas.map((p) => <option key={p.id} value={p.id}>{p.display_name}</option>)}
            </select></label>
          <label className="flex flex-col gap-1 text-small font-medium text-text-secondary">{textos.registro.formaPago}
            <select className={estilosControl} value={forma} onChange={(e) => setForma(e.target.value as CuerpoGasto['payment_method'])}>
              {(['cash', 'debit', 'transfer', 'credit'] as const).map((f) => <option key={f} value={f}>{textos.registro.formas[f]}</option>)}
            </select></label>
          <Campo etiqueta={textos.registro.descripcion} value={descripcion} onChange={(e) => setDescripcion(e.target.value)} />
          <Campo etiqueta={textos.registro.nota} value={nota} onChange={(e) => setNota(e.target.value)} className="phone:col-span-2" />
          {forma === 'credit' && <div className="phone:col-span-2"><LineaCredito monto={monto} dia={dia} comercio={descripcion} alSalir={() => cambiar({ registro: null })} /></div>}
        </div>
      )}
      {errores._ && <Franja rol="error" alerta>{errores._}</Franja>}
      <Boton type="submit" variante="primario" tam="hoja" disabled={guardando}>{textos.app.guardar}</Boton>
    </form>
  )
}
function FormPeso({ refPeso, alTerminar }: { refPeso: React.RefObject<HTMLInputElement | null>; alTerminar: () => void }) {
  const { sesion } = useSesion()
  const avisar = useAvisos()
  const personas = usarPersonas().filter((p) => p.participation !== 'external')
  const [persona, setPersona] = useState<number | null>(sesion?.person_id ?? null)
  const [peso, setPeso] = useState('')
  const [error, setError] = useState<string | null>(null)
  const ultimo = usarConsulta('get', '/api/v1/weight-entries', { consulta: { person_id: persona ? [persona] : undefined, sort: '-measured_date', limit: 1 } }, persona !== null)
  const previo = ultimo.datos?.items?.[0]
  const alGuardar = async (e: FormEvent) => {
    e.preventDefault()
    const g = aGramos(peso)
    if (!g || !persona) return setError(textos.registro.pesoInvalido)
    const cuerpo = { person_id: persona, measured_date: hoyISO(), weight_g: g }
    const clave = nuevaClave()
    const enCola = async () => { await encolarPeso(cuerpo, clave); avisar({ texto: textos.registro.pesoEnCola, rol: 'exito', accion: { texto: textos.app.deshacer, alHacer: () => quitarDeCola(clave) } }) }
    try {
      if (!hayConexion()) await enCola()
      else {
        const w = await pedir('post', '/api/v1/weight-entries', { cuerpo, idempotencia: clave })
        avisar({ texto: textos.registro.pesoGuardado, rol: 'exito', accion: { texto: textos.app.deshacer, alHacer: async () => { await pedir('delete', '/api/v1/weight-entries/{weight_entry_id}', { ruta: { weight_entry_id: w.id } }); avisarCambio('peso') } } })
        avisarCambio('peso')
      }
      alTerminar()
    } catch (err) {
      if (esSinConexion(err)) { await enCola(); alTerminar() } else setError(err instanceof ErrorApi ? err.campo('weight_g') ?? err.message : textos.errores.generico)
    }
  }
  return (
    <form onSubmit={alGuardar} className="flex flex-col gap-4" noValidate>
      <div className="flex flex-wrap gap-2" role="group" aria-label={textos.registro.persona}>
        {personas.map((p) => <ChipSeleccionable key={p.id} seleccionado={persona === p.id} marca={<MarcaPersona codigo={p.code} />} onClick={() => setPersona(p.id)}>{p.display_name}</ChipSeleccionable>)}
      </div>
      <Campo ref={refPeso} etiqueta={textos.registro.pesoKg} inputMode="decimal" autoComplete="off" value={peso} onChange={(e) => setPeso(e.target.value)} error={error ?? undefined}
        ayuda={previo ? textos.registro.ultimoPeso(kilos(previo.weight_g), fecha(previo.measured_date)) : textos.registro.sinPesoPrevio} className="[&_input]:h-12 [&_input]:text-title-2 [&_input]:tabular-nums" />
      <Boton type="submit" variante="primario" tam="hoja">{textos.app.guardar}</Boton>
    </form>
  )
}
