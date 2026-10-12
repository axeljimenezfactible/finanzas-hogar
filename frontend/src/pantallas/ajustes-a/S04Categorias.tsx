// P-26.04 Categorías (U4 parte a): nombre, jerarquía, tipo, reparto por defecto con vista previa ($199.00 y gastos fijos), orden y archivo.
import { useEffect, useRef, useState } from 'react'
import { Plus } from 'lucide-react'
import { nuevaClave, pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { usarConsulta } from '../../api/consulta'
import { dinero } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { Campo, estilosControl } from '../../componentes/Campo'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { useAvisos } from '../../componentes/Avisos'
import { ErrorApi } from '../../api/errores'
import { textos } from '../../textos/es-MX'
import { ListaOrdenable } from '../../vistas-bd/ListaOrdenable'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { repartir } from './calculos'
import { Bloque, BloquePrevia, Editable, vNombre, type Previa } from './comunes'
import { comoLista } from './comunes'
import { t } from './textos'
const k = t.categorias
type Cat = Esquemas['Category']; type Regla = Esquemas['CostShareRule']; type Per = Esquemas['Person']; type PI = Esquemas['PlannedItem']
const EJEMPLO = 19900

export function partesTexto(monto: number, regla: Regla | undefined, personas: Per[]) {
  const r = repartir(monto, regla)
  if (!r) return '—'
  return r.map((x) => `${x.persona === regla?.remainder_person_id ? t.presupuesto.resto.split(' ')[0] : (personas.find((p) => p.id === x.persona)?.display_name ?? '')} ${dinero(x.cents)}`).join(' · ')
}
export default function S04Categorias() {
  const cats = comoLista<Esquemas['Category']>(usarConsulta('get', '/api/v1/categories', { consulta: { limit: 200 } }))
  const reglas = comoLista<Esquemas['CostShareRule']>(usarConsulta('get', '/api/v1/cost-share-rules', { consulta: { limit: 200 } }))
  const personas = comoLista<Esquemas['Person']>(usarConsulta('get', '/api/v1/persons'))
  const fijos = comoLista<Esquemas['PlannedItem']>(usarConsulta('get', '/api/v1/planned-items', { consulta: { kind: 'fixed_expense', limit: 200 } }))
  const [nueva, setNueva] = useState(false)
  const [orden, setOrden] = useState(false)
  // U4 §0.2 punto 2 y l. 235: al soltar una tarjeta del tablero en otro reparto queda ahí con la vista previa;
  // «Aplicar» (o Enter, con el foco en el botón) guarda con updateCategory y Esc la regresa a su columna sin guardar.
  const [mov, setMov] = useState<{ c: Cat; regla: number | null; columna: string } | null>(null)
  const avisarMov = useAvisos()
  const panelMov = useRef<HTMLDivElement>(null)
  useEffect(() => {
    if (!mov) return
    // El foco va a «Aplicar» (Enter aplica) después de que dnd-kit devuelve el foco a la tarjeta al soltar.
    let r2 = 0
    const r1 = requestAnimationFrame(() => { r2 = requestAnimationFrame(() => panelMov.current?.querySelector<HTMLButtonElement>('button')?.focus()) })
    const f = (e: KeyboardEvent) => { if (e.key === 'Escape') { e.preventDefault(); setMov(null) } }
    window.addEventListener('keydown', f)
    return () => { cancelAnimationFrame(r1); cancelAnimationFrame(r2); window.removeEventListener('keydown', f) }
  }, [mov])
  useEffect(() => { const f = (e: KeyboardEvent) => { if (e.altKey && e.key.toLowerCase() === 'n') { e.preventDefault(); setNueva(true) } }; window.addEventListener('keydown', f); return () => window.removeEventListener('keydown', f) }, [])
  const rs = reglas.datos?.items ?? [], ps = personas.datos?.items ?? []
  const reglaDe = (id?: number | null) => rs.find((r) => (id != null ? r.id === id : r.is_default))
  return (
    <Bloque titulo={k.titulo} accion={<div className="flex gap-2"><Boton onClick={() => setOrden(!orden)}>{k.reordenar}</Boton><Boton variante="primario" onClick={() => setNueva(true)}><Plus aria-hidden className="size-4" />{k.nueva}</Boton></div>}>
      <ConDatos consulta={cats}>{(d) => {
        const filas = [...d.items].sort((a, b) => a.sort_order - b.sort_order)
        if (!filas.length) return <Vacio texto={k.sin} accion={<Boton variante="primario" onClick={() => setNueva(true)}>{k.nueva}</Boton>} />
        const guardar = (c: Cat) => async (cambio: Esquemas['CategoryUpdate']) => { await pedir('patch', '/api/v1/categories/{category_id}', { ruta: { category_id: c.id }, cuerpo: cambio }); cats.recargar() }
        const previaReparto = (c: Cat, nuevo: string): Previa => {
          const antes = reglaDe(c.default_cost_share_rule_id), despues = reglaDe(nuevo ? Number(nuevo) : null)
          const sin = (fijos.datos?.items ?? []).filter((x: PI) => x.category_id === c.id && x.cost_share_rule_id == null && x.is_active)
          const cambios: { etiqueta: string; antes: string; despues: string }[] = [{ etiqueta: k.ejemplo, antes: partesTexto(EJEMPLO, antes, ps), despues: partesTexto(EJEMPLO, despues, ps) }]
          sin.forEach((x) => cambios.push({ etiqueta: `${k.fijosCambian(sin.length)} · ${x.description}`, antes: partesTexto(x.amount_cents, antes, ps), despues: partesTexto(x.amount_cents, despues, ps) }))
          return { cambios, notas: [k.aplicaDesde] }
        }
        const opRegla = [{ valor: '', texto: k.porDefecto }, ...rs.map((r) => ({ valor: String(r.id), texto: r.name }))]
        if (orden) return <ListaOrdenable items={filas} idDe={(c) => String(c.id)} pintar={(c) => <span>{c.name}</span>} alReordenar={async (n) => { await pedir('put', '/api/v1/categories/order', { cuerpo: { ids: n.map((c) => c.id) } as Esquemas['IdList'] }); cats.recargar() }} />
        const aplicarMov = async () => {
          if (!mov) return
          try { await guardar(mov.c)({ default_cost_share_rule_id: mov.regla }); setMov(null) } catch (e) { avisarMov({ texto: e instanceof ErrorApi ? e.message : textos.errores.generico, rol: 'error' }) }
        }
        return <>{mov && (
          <div key={mov.c.id} ref={panelMov} className="mb-3 flex flex-col gap-2 rounded-md border border-border-divider p-3" data-prueba="previa-tablero">
            <p className="text-small font-medium">{k.mover(mov.c.name, mov.columna)}</p>
            <BloquePrevia previa={previaReparto(mov.c, mov.regla == null ? '' : String(mov.regla))} />
            <div className="flex flex-wrap gap-2"><Boton variante="primario" onClick={aplicarMov}>{k.aplicar}</Boton><Boton onClick={() => setMov(null)}>{t.comun.cancelar}</Boton></div>
          </div>)}
        <VistaBD def={{ clave: 'P-26.04-categorias', filas, idFila: (c) => c.id, vistas: ['tabla', 'tablero'], titulo: (c) => c.name,
          columnasTablero: { propiedad: 'reparto', columnas: [k.porDefecto, ...rs.map((r) => r.name)], moverA: async (c, columna) => setMov({ c, columna, regla: columna === k.porDefecto ? null : rs.find((r) => r.name === columna)?.id ?? null }) },
          propiedades: [
            { id: 'name', nombre: k.nombre, valor: (c) => c.name, pintar: (c) => <Editable etiqueta={`${k.nombre} · ${c.name}`} valor={c.name} mostrar={c.name} validar={vNombre} guardar={(s) => guardar(c)({ name: s.trim() })} /> },
            { id: 'code', nombre: k.clave, valor: (c) => c.code ?? '', pintar: (c) => <span className="px-2 text-small text-text-secondary">{c.code ?? '—'}</span> },
            { id: 'parent_id', nombre: k.padre, valor: (c) => c.parent_id, pintar: (c) => <Editable etiqueta={`${k.padre} · ${c.name}`} tipo="opcion" valor={String(c.parent_id ?? '')} mostrar={filas.find((x) => x.id === c.parent_id)?.name ?? k.raiz} opciones={[{ valor: '', texto: k.raiz }, ...filas.filter((x) => x.id !== c.id).map((x) => ({ valor: String(x.id), texto: x.name }))]} guardar={(s) => guardar(c)({ parent_id: s ? Number(s) : null })} /> },
            { id: 'kind', nombre: k.tipo, valor: (c) => c.kind, grupo: (c) => (c.kind === 'fixed' ? k.fijo : c.kind === 'variable' ? k.variable : k.sinClasificar), pintar: (c) => <Editable etiqueta={`${k.tipo} · ${c.name}`} tipo="opcion" valor={c.kind ?? ''} mostrar={c.kind === 'fixed' ? k.fijo : c.kind === 'variable' ? k.variable : k.sinClasificar} opciones={[{ valor: '', texto: k.sinClasificar }, { valor: 'fixed', texto: k.fijo }, { valor: 'variable', texto: k.variable }]} guardar={(s) => guardar(c)({ kind: (s || null) as Cat['kind'] | null })} /> },
            { id: 'reparto', nombre: k.reparto, valor: (c) => c.default_cost_share_rule_id, grupo: (c) => (mov?.c.id === c.id ? mov.columna : reglaDe(c.default_cost_share_rule_id)?.name ?? k.porDefecto), pintar: (c) => <Editable etiqueta={`${k.reparto} · ${c.name}`} tipo="opcion" valor={String(c.default_cost_share_rule_id ?? '')} mostrar={rs.find((r) => r.id === c.default_cost_share_rule_id)?.name ?? k.porDefecto} opciones={opRegla} vistaPrevia={(s) => previaReparto(c, s)} guardar={(s) => guardar(c)({ default_cost_share_rule_id: s ? Number(s) : null })} /> },
            { id: 'expense_count', nombre: k.gastos, alinear: 'derecha', valor: (c) => c.expense_count ?? 0, pintar: (c) => <span className="px-2 tabular-nums">{c.expense_count ?? 0}</span> },
            { id: 'is_active', nombre: k.activa, valor: (c) => c.is_active, pintar: (c) => <Editable etiqueta={`${k.activa} · ${c.name}`} tipo="interruptor" valor={c.is_active ? 'si' : 'no'} mostrar={null}
              vistaPrevia={(s) => ({ notas: s === 'no' ? [...(c.code === 'groceries' || c.code === 'eating_out' ? [k.usanCalculos] : []), k.archivar(c.expense_count ?? 0)] : [] })} guardar={(s) => guardar(c)({ is_active: s === 'si' })} /> },
          ] }} /></>
      }}</ConDatos>
      <NuevaCategoria abierta={nueva} alCambiar={setNueva} reglas={rs} alCrear={cats.recargar} />
    </Bloque>
  )
}
function NuevaCategoria({ abierta, alCambiar, reglas, alCrear }: { abierta: boolean; alCambiar: (a: boolean) => void; reglas: Regla[]; alCrear: () => void }) {
  const conexion = useConexion()
  const avisar = useAvisos()
  const [v, setV] = useState({ name: '', kind: '', regla: '' })
  const [err, setErr] = useState<string | null>(null)
  const [clave] = useState(nuevaClave)
  const guardar = async () => {
    const e = vNombre(v.name); setErr(e)
    if (e) return
    try { await pedir('post', '/api/v1/categories', { idempotencia: clave, cuerpo: { name: v.name.trim(), ...(v.kind ? { kind: v.kind as 'fixed' } : {}), ...(v.regla ? { default_cost_share_rule_id: Number(v.regla) } : {}), is_active: true } } as never); setV({ name: '', kind: '', regla: '' }); alCambiar(false); alCrear() }
    catch (x) { if (x instanceof Error && 'estado' in x && (x as { estado: number }).estado === 409) setErr(t.v.v12); else avisar({ texto: t.comun.noGuardo, rol: 'error' }) }
  }
  return (
    <Hoja abierta={abierta} alCambiar={alCambiar} titulo={k.nueva} pie={<div className="flex justify-end gap-2"><Boton onClick={() => alCambiar(false)}>{textos.app.cancelar}</Boton><Boton variante="primario" deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => void guardar()}>{textos.app.guardar}</Boton></div>}>
      <form className="flex flex-col gap-3" onSubmit={(e) => { e.preventDefault(); void guardar() }} onKeyDown={(e) => { if (e.key === 'Enter' && e.ctrlKey) { e.preventDefault(); void guardar() } }}>
        <Campo etiqueta={k.nombre} value={v.name} error={err ?? undefined} onChange={(e) => setV({ ...v, name: e.target.value })} />
        <label className="flex flex-col gap-1 text-small font-medium text-text-secondary">{k.tipo}<select className={estilosControl} value={v.kind} onChange={(e) => setV({ ...v, kind: e.target.value })}><option value="">{k.sinClasificar}</option><option value="fixed">{k.fijo}</option><option value="variable">{k.variable}</option></select></label>
        <label className="flex flex-col gap-1 text-small font-medium text-text-secondary">{k.reparto}<select className={estilosControl} value={v.regla} onChange={(e) => setV({ ...v, regla: e.target.value })}><option value="">{k.porDefecto}</option>{reglas.map((r) => <option key={r.id} value={r.id}>{r.name}</option>)}</select></label>
      </form>
    </Hoja>
  )
}
