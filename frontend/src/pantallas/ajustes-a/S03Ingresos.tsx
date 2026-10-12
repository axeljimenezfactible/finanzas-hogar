// P-26.03 Ingresos (U4 parte a): ingreso quincenal esperado de Axel (días 15 y 30 o último), desde cuándo rige; Ximena no aporta (decisión B).
import { useState } from 'react'
import { nuevaClave, pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { usarConsulta } from '../../api/consulta'
import { aCentavos, dinero, fecha, hoyISO } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { abrirRegistro } from '../../shell/estado'
import { textos } from '../../textos/es-MX'
import { Bloque, Editable, Fila, textoDinero, vDinero, vFecha, type Previa } from './comunes'
import { comoLista } from './comunes'
import { t } from './textos'
const g = t.ingresos
type PI = Esquemas['PlannedItem']
const mesDeIngreso = (x: PI) => x.amount_cents + (x.second_half_amount_cents ?? x.amount_cents)
const texto = (x: PI) => g.previa(fecha(x.start_date), dinero(x.amount_cents), dinero(x.second_half_amount_cents ?? x.amount_cents), dinero(mesDeIngreso(x)))

export default function S03Ingresos() {
  const partidas = comoLista<Esquemas['PlannedItem']>(usarConsulta('get', '/api/v1/planned-items', { consulta: { kind: 'income', limit: 50 } }))
  const personas = comoLista<Esquemas['Person']>(usarConsulta('get', '/api/v1/persons'))
  const [nuevo, setNuevo] = useState(false)
  const titular = personas.datos?.items.find((x) => x.participation === 'holder')
  const ximena = personas.datos?.items.find((x) => x.participation === 'co_holder')
  return (
    <div className="flex flex-col gap-4">
      <ConDatos consulta={partidas}>{(d) => d.items.length === 0
        ? <Vacio texto={g.sinIngreso} accion={<Boton variante="primario" onClick={() => setNuevo(true)}>{g.capturar}</Boton>} />
        : <>{d.items.map((x) => <Partida key={x.id} x={x} titular={titular?.display_name ?? ''} alCambiar={partidas.recargar} />)}</>}</ConDatos>
      {ximena && <Bloque titulo={ximena.display_name}><p className="text-body text-text-secondary">{g.ximena}</p></Bloque>}
      <Boton variante="enlace" className="w-fit min-h-11" onClick={() => abrirRegistro('ingreso')}>{g.extra}</Boton>
      <NuevoIngreso abierta={nuevo} alCambiar={setNuevo} titular={titular?.id} alCrear={partidas.recargar} />
    </div>
  )
}
function Partida({ x, titular, alCambiar }: { x: PI; titular: string; alCambiar: () => void }) {
  const guardar = async (c: Esquemas['PlannedItemUpdate']) => { await pedir('patch', '/api/v1/planned-items/{planned_item_id}', { ruta: { planned_item_id: x.id }, cuerpo: c }); alCambiar() }
  const previa = (c: Partial<PI>): Previa => ({ cambios: [{ etiqueta: g.esperado, antes: texto(x), despues: texto({ ...x, ...c }) }], notas: [g.aviso] })
  return (
    <Bloque titulo={x.description || g.concepto}>
      <Fila etiqueta={g.persona}><span className="px-2">{titular}</span></Fila>
      <Fila etiqueta={g.frecuencia}><span className="px-2">{g.quincenal}</span></Fila>
      <div className="flex flex-col desktop:grid desktop:grid-cols-2 desktop:gap-x-6">
        <Fila etiqueta={g.q15}><Editable etiqueta={g.q15} tipo="dinero" valor={textoDinero(x.amount_cents)} mostrar={<span className="text-title-3 tabular-nums">{dinero(x.amount_cents)}</span>} validar={vDinero(true)} vistaPrevia={(s) => previa({ amount_cents: aCentavos(s) ?? 0 })} guardar={(s) => guardar({ amount_cents: aCentavos(s) ?? 0 })} /></Fila>
        <Fila etiqueta={g.q30}><Editable etiqueta={g.q30} tipo="dinero" valor={textoDinero(x.second_half_amount_cents)} mostrar={<span className="text-title-3 tabular-nums">{x.second_half_amount_cents != null ? dinero(x.second_half_amount_cents) : g.igual}</span>} validar={vDinero(true, true)}
          vistaPrevia={(s) => previa({ second_half_amount_cents: s ? (aCentavos(s) ?? 0) : undefined })} guardar={(s) => guardar({ second_half_amount_cents: s ? (aCentavos(s) ?? 0) : null })} /></Fila>
      </div>
      <Fila etiqueta={g.rige}><Editable etiqueta={g.rige} tipo="fecha" valor={x.start_date} mostrar={fecha(x.start_date)} validar={vFecha()} vistaPrevia={(s) => previa({ start_date: s })} guardar={(s) => guardar({ start_date: s })} /></Fila>
    </Bloque>
  )
}
function NuevoIngreso({ abierta, alCambiar, titular, alCrear }: { abierta: boolean; alCambiar: (a: boolean) => void; titular?: number; alCrear: () => void }) {
  const conexion = useConexion()
  const [v, setV] = useState({ a: '', b: '', desde: hoyISO() })
  const [err, setErr] = useState<Record<string, string | null>>({})
  const [clave] = useState(nuevaClave)
  const guardar = async () => {
    const e = { a: vDinero(true)(v.a), b: vDinero(true, true)(v.b), desde: vFecha()(v.desde) }
    setErr(e)
    if (Object.values(e).some(Boolean) || !titular) return
    await pedir('post', '/api/v1/planned-items', { idempotencia: clave, cuerpo: { kind: 'income', description: g.concepto, person_id: titular, amount_cents: aCentavos(v.a) ?? 0, ...(v.b ? { second_half_amount_cents: aCentavos(v.b) ?? 0 } : {}), frequency: 'fortnightly', interval_count: 1, start_date: v.desde } })
    alCambiar(false); alCrear()
  }
  return (
    <Hoja abierta={abierta} alCambiar={alCambiar} titulo={g.capturar} pie={<div className="flex justify-end gap-2"><Boton onClick={() => alCambiar(false)}>{textos.app.cancelar}</Boton><Boton variante="primario" deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => void guardar()}>{textos.app.guardar}</Boton></div>}>
      <form className="flex flex-col gap-3" onKeyDown={(e) => { if (e.key === 'Enter' && e.ctrlKey) { e.preventDefault(); void guardar() } }} onSubmit={(e) => { e.preventDefault(); void guardar() }}>
        <Campo etiqueta={g.q15} inputMode="decimal" value={v.a} error={err.a ?? undefined} onChange={(e) => setV({ ...v, a: e.target.value })} />
        <Campo etiqueta={g.q30} inputMode="decimal" value={v.b} error={err.b ?? undefined} ayuda={g.igual} onChange={(e) => setV({ ...v, b: e.target.value })} />
        <Campo etiqueta={g.rige} type="date" value={v.desde} error={err.desde ?? undefined} onChange={(e) => setV({ ...v, desde: e.target.value })} />
      </form>
    </Hoja>
  )
}
