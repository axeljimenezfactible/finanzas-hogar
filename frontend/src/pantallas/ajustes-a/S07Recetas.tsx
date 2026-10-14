// P-26.07 Recetas (U4 parte a): porciones base, tiempos de comida, recetario en línea, bases nutricionales con recarga,
// restricciones alimentarias por persona (vacías, decisión D) y estado de la IA (decisión H).
import { useEffect, useState } from 'react'
import { Plus } from 'lucide-react'
import { pedir, nuevaClave } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { usarConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import { fecha } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Enlace, navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { textos } from '../../textos/es-MX'
import type { Cambio } from '../../vistas-bd/tipos'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { Aviso, Bloque, BloquePrevia, DialogoCritico, Editable, Fila, vNombre } from './comunes'
import { comoLista } from './comunes'
import { guardarVarias, nota } from './grupo'
import { t } from './textos'
const r = t.recetas
type Rec = Esquemas['Recipe']; type Fuente = Esquemas['RecipeSettings']['nutrition_sources'][number]
type Tiempo = Esquemas['RecipeSettings']['meal_slots'][number]
// Ctrl+/ (U4 P-26.07 «Atajos»: «cambiar a la vez los tiempos aptos de varias»). La selección múltiple se edita en grupo con
// «Agregar {tiempo}» o «Quitar {tiempo}» (valor `+id` o `-id`): PROVISIONAL, bloqueo WF3a-bis-ajustes-a §3.
const nuevosTiempos = (x: Rec, v: string) => { const id = Number(v.slice(1)), a = x.meal_slot_ids ?? []; return v.startsWith('+') ? [...new Set([...a, id])].sort((p, q) => p - q) : a.filter((s) => s !== id) }
const sinTiempo = (sel: Rec[], v: string) => sel.filter((x) => nuevosTiempos(x, v).length === 0).length
export function previaTiempos(sel: Rec[], v: string, slots: Tiempo[]): Cambio[] {
  const id = Number(v.slice(1)), n = sel.length, nombre = slots.find((m) => m.id === id)?.name ?? ''
  const tiene = (ids: number[]) => ids.includes(id)
  const cambios: Cambio[] = [{ etiqueta: t.grupo.aptasPara(nombre), antes: t.grupo.deN(sel.filter((x) => tiene(x.meal_slot_ids ?? [])).length, n), despues: t.grupo.deN(sel.filter((x) => tiene(nuevosTiempos(x, v))).length, n) }]
  const k = sinTiempo(sel, v)
  if (k) cambios.push(nota(t.grupo.sinTiempo(k)))
  return cambios
}

export default function S07Recetas() {
  const ajustes = usarConsulta('get', '/api/v1/settings/recipes')
  const recetas = comoLista<Esquemas['Recipe']>(usarConsulta('get', '/api/v1/recipes', { consulta: { limit: 200 } }))
  const perfiles = comoLista<Esquemas['HealthProfile']>(usarConsulta('get', '/api/v1/health-profiles'))
  const personas = comoLista<Esquemas['Person']>(usarConsulta('get', '/api/v1/persons'))
  const ia = usarConsulta('get', '/api/v1/ai/providers')
  const [borrar, setBorrar] = useState<{ x: Rec; n: number } | null>(null)
  const [nueva, setNueva] = useState(false)
  // Alt+N: nueva receta (U4 P-26.07, «Atajos»).
  useEffect(() => { const f = (e: KeyboardEvent) => { if (e.altKey && e.key.toLowerCase() === 'n') { e.preventDefault(); setNueva(true) } }; window.addEventListener('keydown', f); return () => window.removeEventListener('keydown', f) }, [])
  const slots = ajustes.datos?.meal_slots ?? []
  const iaRecetas = (ia.datos?.providers ?? []).some((p) => p.is_enabled && p.allowed_uses.includes('recipe_suggestions'))
  return (
    <div className="flex flex-col gap-4">
      <Bloque titulo={r.porciones} id="b-recetas-base">
        <ConDatos consulta={ajustes}>{(a) => <>
          <Fila etiqueta={r.porciones}><span className="px-2">{r.porcionesTexto(a.default_base_servings)}</span></Fila>
          <Fila etiqueta={r.tiempos}><span className="px-2">{a.meal_slots.map((m) => m.name).join(' · ')}</span></Fila>
          <Fila etiqueta={r.porOrigen}><span className="px-2 tabular-nums">{a.recipe_counts.map((c) => `${r.origenes[c.origin]}: ${c.count}`).join(' · ')}</span></Fila>
        </>}</ConDatos>
      </Bloque>
      <Bloque titulo={r.recetario} id="b-recetario" accion={<Boton variante="primario" onClick={() => setNueva(true)}><Plus aria-hidden className="size-5" />{r.nueva}</Boton>}>
        <ConDatos consulta={recetas}>{(d) => {
          if (!d.items.length) return <Vacio texto={r.sin} accion={<Boton variante="primario" onClick={() => setNueva(true)}>{r.nueva}</Boton>} />
          const guardar = (x: Rec) => async (c: Esquemas['RecipeUpdate']) => { await pedir('patch', '/api/v1/recipes/{recipe_id}', { ruta: { recipe_id: x.id }, cuerpo: c }); recetas.recargar() }
          return <VistaBD def={{ clave: 'P-26.07-recetario', filas: d.items, idFila: (x) => x.id, vistas: ['tabla', 'tablero', 'galeria'], titulo: (x) => x.name,
            columnasTablero: { propiedad: 'origin', columnas: Object.values(r.origenes) },
            editarVarias: {
              propiedades: [{ id: 'meal_slot_ids', nombre: r.aptos, tipo: 'opcion', opciones: slots.flatMap((m) => [{ valor: `+${m.id}`, texto: t.grupo.agregar(m.name) }, { valor: `-${m.id}`, texto: t.grupo.quitar(m.name) }]) }],
              vistaPrevia: async (sel, _p, v) => previaTiempos(sel, v, slots),
              guardar: async (sel, _p, v) => {
                if (sinTiempo(sel, v)) throw new ErrorApi(422, 'validation_error', t.v.tiempos)
                await guardarVarias(sel, (x) => pedir('patch', '/api/v1/recipes/{recipe_id}', { ruta: { recipe_id: x.id }, cuerpo: { meal_slot_ids: nuevosTiempos(x, v) } }), recetas.recargar)
              },
            },
            propiedades: [
              { id: 'name', nombre: r.nombre, valor: (x) => x.name, pintar: (x) => <Editable etiqueta={`${r.nombre} · ${x.name}`} valor={x.name} mostrar={x.name} validar={vNombre} guardar={(s) => guardar(x)({ name: s.trim() })} /> },
              { id: 'slots', nombre: r.aptos, valor: (x) => (x.meal_slot_ids ?? []).length, pintar: (x) => <Tiempos x={x} slots={slots} guardar={guardar(x)} /> },
              { id: 'origin', nombre: r.origen, valor: (x) => x.origin, grupo: (x) => r.origenes[x.origin], pintar: (x) => r.origenes[x.origin] },
              { id: 'abrir', nombre: r.duplicar, valor: () => '', pintar: (x) => <span className="flex flex-wrap gap-2"><Enlace a={rutaDe('P-12', { id: x.id })} className="inline-flex min-h-11 items-center text-text-link desktop:min-h-0">{x.name}</Enlace><Boton variante="fantasma" onClick={async () => { const n = await pedir('post', '/api/v1/recipes/{recipe_id}/duplicate', { ruta: { recipe_id: x.id }, idempotencia: nuevaClave() }); navegar(rutaDe('P-12', { id: n.id })) }}>{r.duplicar}</Boton>
                <Boton variante="fantasma" onClick={async () => { const u = await pedir('get', '/api/v1/recipes/{recipe_id}/usage', { ruta: { recipe_id: x.id } }); setBorrar({ x, n: u.times_used_8w }) }}>{t.comun.eliminar}</Boton></span> },
            ] }} />
        }}</ConDatos>
        <DialogoCritico abierto={!!borrar} alCambiar={(a) => !a && setBorrar(null)} titulo={borrar ? r.eliminar(borrar.x.name) : ''} verbo={borrar ? r.eliminar(borrar.x.name) : ''}
          alConfirmar={async () => { if (borrar) await pedir('delete', '/api/v1/recipes/{recipe_id}', { ruta: { recipe_id: borrar.x.id } }); recetas.recargar() }}>{borrar ? r.uso(borrar.n) : ''}</DialogoCritico>
        <NuevaReceta abierta={nueva} alCambiar={setNueva} slots={slots} />
      </Bloque>
      <div className="flex flex-col gap-4 desktop:grid desktop:grid-cols-2">
        <Bloque titulo={r.bases} id="b-bases">
          <ConDatos consulta={ajustes}>{(a) => <ul className="flex flex-col">{a.nutrition_sources.map((f) => <BaseNutricional key={f.source} f={f} alCambiar={ajustes.recargar} />)}</ul>}</ConDatos>
        </Bloque>
        <Bloque titulo={r.restricciones} id="b-restricciones">
          <ConDatos consulta={perfiles}>{(d) => <>{d.items.map((p) => { const nombre = personas.datos?.items.find((x) => x.id === p.person_id)?.display_name ?? ''; return (
            <Fila key={p.id} etiqueta={r.de(nombre)}><Editable etiqueta={r.de(nombre)} tipo="area" valor={p.dietary_restrictions ?? ''} mostrar={p.dietary_restrictions || t.comun.sinValor}
              guardar={async (s) => { await pedir('patch', '/api/v1/health-profiles/{health_profile_id}', { ruta: { health_profile_id: p.id }, cuerpo: { dietary_restrictions: s.trim() || null } }); perfiles.recargar() }} /></Fila>) })}</>}</ConDatos>
          <Fila etiqueta={r.iaRecetas}><span className="flex flex-wrap items-center gap-2 px-2">{iaRecetas ? t.tiendas.iaEncendida : t.tiendas.iaApagada}<Enlace a="/ajustes/11" className="inline-flex min-h-11 items-center text-text-link">{t.tiendas.iaEnlace}</Enlace></span></Fila>
        </Bloque>
      </div>
    </div>
  )
}
function Tiempos({ x, slots, guardar }: { x: Rec; slots: Esquemas['MealSlot'][]; guardar: (c: Esquemas['RecipeUpdate']) => Promise<void> }) {
  const conexion = useConexion()
  const [error, setError] = useState<string | null>(null)
  const ids = x.meal_slot_ids ?? []
  return (
    <fieldset className="flex flex-wrap gap-x-2"><legend className="sr-only">{r.aptos} · {x.name}</legend>
      {slots.map((m) => <label key={m.id} className="flex min-h-11 items-center gap-1 text-small desktop:min-h-9"><input type="checkbox" className="size-5" disabled={!conexion} checked={ids.includes(m.id)}
        onChange={async (e) => { const n = e.target.checked ? [...ids, m.id] : ids.filter((i) => i !== m.id); if (!n.length) return setError(t.v.tiempos); setError(null); await guardar({ meal_slot_ids: n }) }} />{m.name}</label>)}
      {error && <span className="w-full text-small text-role-error">{error}</span>}
    </fieldset>
  )
}
function BaseNutricional({ f, alCambiar }: { f: Fuente; alCambiar: () => void }) {
  const conexion = useConexion()
  const [abierta, setAbierta] = useState(false)
  const [msg, setMsg] = useState<{ rol: 'exito' | 'error'; texto: string } | null>(null)
  const recargable = f.source === 'usda_fdc' || f.source === 'open_food_facts'
  const nombre = r.fuentes[f.source]
  return (
    <li className="flex flex-col gap-1 border-b border-border-divider py-2 last:border-b-0">
      <div className="flex flex-wrap items-center gap-x-4 gap-y-1"><span className="min-w-36 font-medium">{nombre}</span>
        <span className="text-small tabular-nums text-text-secondary">{f.loaded_at ? `${r.cargada}: ${fecha(f.loaded_at.slice(0, 10))} · ${r.alimentos}: ${f.food_count.toLocaleString('es-MX')}` : r.sinDatos(nombre)}</span>
        {recargable && <Boton className="ml-auto" deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => setAbierta(true)}>{r.recargar}</Boton>}</div>
      {msg && <Aviso rol={msg.rol}>{msg.texto}</Aviso>}
      <Hoja abierta={abierta} alCambiar={setAbierta} titulo={r.recargar} pie={<div className="flex justify-end gap-2"><Boton onClick={() => setAbierta(false)}>{textos.app.cancelar}</Boton><Boton variante="primario" onClick={async () => {
        try { const x = await pedir('post', '/api/v1/nutrition-facts/reload', { idempotencia: nuevaClave(), cuerpo: { source: f.source as 'usda_fdc' } }); setMsg({ rol: 'exito', texto: r.recargada(x.loaded_count) }); alCambiar() } catch (e) { setMsg({ rol: 'error', texto: e instanceof ErrorApi ? e.message : t.comun.noGuardo }) } setAbierta(false) }}>{textos.app.confirmar}</Boton></div>}>
        <BloquePrevia previa={{ notas: [r.recargarTexto(nombre, f.loaded_at ? fecha(f.loaded_at.slice(0, 10)) : r.sinCarga, f.food_count)] }} />
      </Hoja>
    </li>
  )
}

// Nueva receta (U4 P-26.07, «Nueva receta · duplicar receta»: acciones que abren P-12). Aquí se capturan el nombre y los tiempos aptos
// (lo que se edita en la tabla, V-01, V-09 y «al menos un tiempo»); createRecipe la crea y se abre P-12 para ingredientes y pasos.
function NuevaReceta({ abierta, alCambiar, slots }: { abierta: boolean; alCambiar: (a: boolean) => void; slots: Esquemas['MealSlot'][] }) {
  const conexion = useConexion()
  const [nombre, setNombre] = useState('')
  const [ids, setIds] = useState<number[]>([])
  const [err, setErr] = useState<{ nombre?: string | null; tiempos?: string | null }>({})
  const [clave] = useState(nuevaClave)
  const guardar = async () => {
    const e = { nombre: vNombre(nombre), tiempos: ids.length ? null : t.v.tiempos }; setErr(e)
    if (e.nombre || e.tiempos) return
    try { const n = await pedir('post', '/api/v1/recipes', { idempotencia: clave, cuerpo: { name: nombre.trim(), meal_slot_ids: ids } }); setNombre(''); setIds([]); alCambiar(false); navegar(rutaDe('P-12', { id: n.id })) }
    catch (x) { setErr({ nombre: x instanceof ErrorApi && x.estado === 409 ? t.v.v12 : t.comun.noGuardo }) }
  }
  return (
    <Hoja abierta={abierta} alCambiar={alCambiar} titulo={r.nueva} pie={<div className="flex justify-end gap-2"><Boton onClick={() => alCambiar(false)}>{textos.app.cancelar}</Boton><Boton variante="primario" deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => void guardar()}>{textos.app.guardar}</Boton></div>}>
      <form className="flex flex-col gap-3" onSubmit={(e) => { e.preventDefault(); void guardar() }} onKeyDown={(e) => { if (e.key === 'Enter' && e.ctrlKey) { e.preventDefault(); void guardar() } }}>
        <Campo etiqueta={r.nombre} value={nombre} error={err.nombre ?? undefined} onChange={(e) => setNombre(e.target.value)} />
        <fieldset className="flex flex-col gap-1"><legend className="text-small font-medium text-text-secondary">{r.aptos}</legend>
          <span className="flex flex-wrap gap-x-3">{slots.map((m) => <label key={m.id} className="flex min-h-11 items-center gap-2 text-small"><input type="checkbox" className="size-5" checked={ids.includes(m.id)} onChange={(e) => setIds(e.target.checked ? [...ids, m.id] : ids.filter((i) => i !== m.id))} />{m.name}</label>)}</span>
          {err.tiempos && <span className="text-small text-role-error">{err.tiempos}</span>}
        </fieldset>
      </form>
    </Hoja>
  )
}
