// Piezas de P-03 «Registrar» que aporta el área Gastos y Presupuesto (WF-3a bis, U1 §3.5):
// - N-053: con la forma de pago «Crédito», elegir la línea activa (nombre y tipo, orden alfabético; listCredits) y seguir
//   en P-19 prellenada con monto, fecha y comercio.
// - N-102: pestaña «Ingreso»: registrar el ingreso extra en la misma hoja (createIncome, con clave de idempotencia).
// Las monta src/shell/HojaRegistro.tsx (cambio pedido a frontend-base); P-08 usa FormIngreso en su hoja «Registrar ingreso
// extra», así que N-102 tiene una sola implementación. No importa ./comun para no sumar sus lecturas a P-03 en la matriz C3.
import { useRef, useState, type FormEvent, type RefObject } from 'react'
import { nuevaClave, pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { usarConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import { aCentavos, fecha, hoyISO } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { useAvisos } from '../../componentes/Avisos'
import { Boton } from '../../componentes/Boton'
import { Campo, estilosControl } from '../../componentes/Campo'
import { ConDatos, Franja } from '../../componentes/Estados'
import { cn } from '../../componentes/utils'
import { textos } from '../../textos/es-MX'
import { t } from './textos'

const r = t.registro
/** Centavos → «1234.56» (P-19 lee `monto` como texto en pesos; nunca flotantes, protocolo §6). */
const aTexto = (c: number) => `${Math.floor(c / 100)}.${String(c % 100).padStart(2, '0')}`

/** N-053. `comercio` es la descripción o comercio que se escribió en la hoja. */
export function LineaCredito({ monto, dia, comercio, alSalir }: { monto: string; dia: string; comercio: string; alSalir: () => void }) {
  const enLinea = useConexion()
  const lineas = usarConsulta('get', '/api/v1/credits', { consulta: { limit: 200 } })
  const [linea, setLinea] = useState('')
  if (!enLinea) return <Franja rol="sinConexion">{textos.conexion.requiere}</Franja>
  return (
    <ConDatos consulta={lineas}>{(d) => {
      const activas = (d.items as Esquemas['Credit'][]).filter((c) => c.is_active).sort((a, b) => a.name.localeCompare(b.name, 'es-MX'))
      if (!activas.length) return (
        <Franja rol="info" accion={<Boton variante="enlace" onClick={() => { alSalir(); navegar(rutaDe('P-26', { seccion: '05' })) }}>{r.irAjustesCreditos}</Boton>}>{r.sinLineas}</Franja>)
      const elegida = activas.find((c) => String(c.id) === linea)
      const seguir = () => {
        if (!elegida) return
        const q = new URLSearchParams()
        const c = aCentavos(monto); if (c != null && c > 0) q.set('monto', aTexto(c))
        if (dia) q.set('fecha', dia)
        if (comercio.trim()) q.set('comercio', comercio.trim())
        alSalir()
        navegar(`${rutaDe('P-19', { id: elegida.id, compra: 'nueva' })}?${q}`)
      }
      return (
        <div className="flex flex-col gap-2">
          <label className="flex flex-col gap-1 text-small font-medium text-text-secondary">{r.linea}
            <select value={linea} onChange={(e) => setLinea(e.target.value)} className={cn(estilosControl, 'w-full')}>
              <option value="">{r.elegirLinea}</option>
              {activas.map((c) => <option key={c.id} value={c.id}>{c.name} · {r.tiposLinea[c.kind] ?? c.kind}</option>)}
            </select>
          </label>
          <p className="text-caption text-text-secondary">{r.avisoCredito}</p>
          <Boton variante="secundario" deshabilitadoPor={elegida ? undefined : r.elegirLinea} onClick={seguir}>{r.continuarCompra}</Boton>
        </div>
      )
    }}</ConDatos>
  )
}

/** N-102: fecha (hoy por defecto), concepto, monto y nota. Sin conexión: «Requiere conexión» en lugar del contenido (U1 §2.4). */
export function FormIngreso({ refMonto, alTerminar }: { refMonto?: RefObject<HTMLInputElement | null>; alTerminar: () => void }) {
  const enLinea = useConexion()
  const avisar = useAvisos()
  const clave = useRef(nuevaClave()) // un reintento de la misma captura reutiliza la clave (C1 §4)
  const [v, setV] = useState({ monto: '', concepto: '', dia: hoyISO(), nota: '' })
  const [err, setErr] = useState<Record<string, string | undefined>>({})
  const [enCurso, setEnCurso] = useState(false)
  if (!enLinea) return <Franja rol="sinConexion">{textos.conexion.requiere}</Franja>
  const poner = (k: keyof typeof v) => (e: { target: { value: string } }) => { setV({ ...v, [k]: e.target.value }); setErr({}) }
  const guardar = async (e: FormEvent) => {
    e.preventDefault()
    const c = aCentavos(v.monto)
    const errores: Record<string, string> = {}
    if (c == null || c <= 0) errores.amount_cents = textos.registro.montoInvalido
    if (!v.concepto.trim()) errores.description = r.conceptoFalta
    if (!v.dia) errores.income_date = textos.errores.revisaCampos
    if (Object.keys(errores).length) return setErr(errores)
    setEnCurso(true)
    try {
      await pedir('post', '/api/v1/incomes', { idempotencia: clave.current, cuerpo: { kind: 'extra', description: v.concepto.trim(), amount_cents: c!, income_date: v.dia, notes: v.nota.trim() || undefined } })
      clave.current = nuevaClave()
      avisar({ texto: r.ingresoRegistrado, rol: 'exito' })
      window.dispatchEvent(new CustomEvent('fh:datos-cambiaron', { detail: { tipo: 'ingreso' } }))
      alTerminar()
    } catch (x) {
      setErr(x instanceof ErrorApi
        ? { amount_cents: x.campo('amount_cents'), description: x.campo('description'), income_date: x.campo('income_date'), general: x.message }
        : { general: textos.errores.generico })
    } finally { setEnCurso(false) }
  }
  return (
    <form noValidate onSubmit={guardar} className="flex flex-col gap-3">
      <Campo ref={refMonto} etiqueta={t.comun.monto} inputMode="decimal" autoComplete="off" enterKeyHint="next" placeholder="$0.00" value={v.monto} onChange={poner('monto')} error={err.amount_cents}
        className="[&_input]:h-12 [&_input]:text-title-2 [&_input]:tabular-nums" />
      <Campo etiqueta={t.comun.concepto} value={v.concepto} onChange={poner('concepto')} error={err.description} />
      <Campo etiqueta={t.comun.fecha} type="date" value={v.dia} onChange={poner('dia')} error={err.income_date} ayuda={v.dia ? fecha(v.dia) : undefined} />
      <Campo etiqueta={t.comun.nota} value={v.nota} onChange={poner('nota')} />
      {err.general && <p role="alert" className="text-small text-role-error">{err.general}</p>}
      <Boton type="submit" variante="primario" tam="hoja" aria-busy={enCurso}>{textos.app.guardar}</Boton>
    </form>
  )
}
