// P-12 Receta (U1 §2.1–§2.2, §3.14 N-140 a N-144): escalado, costo y kcal del motor S4; «nueva» crea (N-136).
import { useState } from 'react'
import { Minus, Plus, Trash2 } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import type { PropsPantalla } from '../../app/registro'
import type { Esquemas } from '../../api/tipos'
import { nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { useConexion } from '../../api/conexion'
import { dinero, fecha, hoyISO } from '../../api/formatos'
import { navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { Campo, estilosControl } from '../../componentes/Campo'
import { ChipSeleccionable } from '../../componentes/Chip'
import { ConDatos, Franja } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { useAvisos } from '../../componentes/Avisos'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { BarraAcciones, MenuMas, Seccion, aMilli, cantidad, invalidar, kcal, mensajeError, porId, porciones, usarCatalogo } from './comun'
import { textosCC as t } from './textos'

type Cuerpo = Esquemas['RecipeCreate']
type Detalle = Esquemas['RecipeDetail']
let borrador: Cuerpo | null = null
/** Borrador de «Sugerir con IA» (P-11): se revisa aquí antes de guardar (N-137). */
export const guardarBorrador = (b: Cuerpo) => { borrador = b }

export default function P12Receta({ params }: PropsPantalla) {
  return params.id === 'nueva' ? <NuevaReceta /> : <DetalleReceta id={Number(params.id)} />
}

function NuevaReceta() {
  const avisar = useAvisos()
  const [inicial] = useState<Cuerpo>(() => { const b = borrador ?? { name: '', base_servings: 2, ingredients: [] }; borrador = null; return b })
  async function guardar(c: Cuerpo) {
    const r = await pedir('post', '/api/v1/recipes', { cuerpo: c, idempotencia: nuevaClave() })
    invalidar('recetas'); avisar({ texto: t.recetaCreada, rol: 'exito' })
    navegar(rutaDe('P-12', { id: r.id }), { reemplazar: true })
  }
  return <FormReceta inicial={inicial} textoGuardar={t.crearReceta} alGuardar={guardar} />
}

function DetalleReceta({ id }: { id: number }) {
  const conexion = useConexion()
  const avisar = useAvisos()
  const [porc, setPorc] = useState<number | null>(null)
  const [hoja, setHoja] = useState<'editar' | 'eliminar' | null>(null)
  const det = usarConsulta('get', '/api/v1/recipes/{recipe_id}', { ruta: { recipe_id: id }, consulta: porc ? { servings_milli: porc } : {} })
  const uso = usarConsulta('get', '/api/v1/recipes/{recipe_id}/usage', { ruta: { recipe_id: id } })
  const productos = porId(usarCatalogo('productos'))
  const tiendas = porId(usarCatalogo('tiendas'))
  const tiempos = usarCatalogo('tiempos') ?? []
  const sinRed = conexion ? undefined : textos.conexion.requiere
  async function duplicar() {
    try {
      const r = await pedir('post', '/api/v1/recipes/{recipe_id}/duplicate', { ruta: { recipe_id: id }, idempotencia: nuevaClave() })
      invalidar('recetas'); avisar({ texto: t.duplicada, rol: 'exito' }); navegar(rutaDe('P-12', { id: r.id }))
    } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
  }
  return (
    <ConDatos consulta={det}>{(d) => {
      const r = d.recipe
      const planeada = uso.datos?.next_plan_date && uso.datos.next_plan_date >= hoyISO() ? uso.datos.next_plan_date : undefined
      return (
        <div className="flex flex-col gap-4">
          <BarraAcciones acciones={<>
            <Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => navegar(`${rutaDe('P-10')}?hoja=asignar&receta=${id}`)}>{t.agregarAlPlan}</Boton>
            <MenuMas acciones={[{ texto: t.editar, alHacer: () => setHoja('editar'), deshabilitadoPor: sinRed }, { texto: t.duplicar, alHacer: duplicar, deshabilitadoPor: sinRed }, { texto: t.eliminar, alHacer: () => setHoja('eliminar'), deshabilitadoPor: sinRed }]} />
          </>}>
            <h2 className="text-title-2 font-semibold">{r.name}</h2>
            <p className="text-small text-text-secondary">
              {[(r.meal_slot_ids ?? []).map((x) => tiempos.find((y) => y.id === x)?.name).filter(Boolean).join(', '), `${t.porcionesBase}: ${r.base_servings}`, `${t.tiempoTotal}: ${t.minutos((r.prep_minutes ?? 0) + (r.cook_minutes ?? 0))}`, `${t.origen}: ${t.origenes[r.origin]}`].filter(Boolean).join(' · ')}
            </p>
            {uso.datos && <p className="text-small text-text-secondary">{t.uso(uso.datos.times_used_8w, uso.datos.next_plan_date ? fecha(uso.datos.next_plan_date) : undefined)}</p>}
          </BarraAcciones>
          <div role="group" aria-label={t.porcionesElegidas} className="flex items-center gap-2">
            <span className="text-small font-medium">{t.porcionesElegidas}</span>
            <Boton tam="icono" aria-label={t.menosPorcion} onClick={() => setPorc(Math.max(1000, d.servings_milli - 1000))}><Minus aria-hidden className="size-6" size={ICONO_MD} /></Boton>
            <output className="min-w-11 text-center text-title-3 font-semibold tabular-nums" data-prueba="porciones">{porciones(d.servings_milli)}</output>
            <Boton tam="icono" aria-label={t.masPorcion} onClick={() => setPorc(d.servings_milli + 1000)}><Plus aria-hidden className="size-6" size={ICONO_MD} /></Boton>
          </div>
          <div className="grid grid-cols-2 gap-3 desktop:grid-cols-4">
            <TarjetaMetrica titulo={t.costoTotal} valor={d.total_cost_cents != null ? dinero(d.total_cost_cents) : '—'} rol="expense" />
            <TarjetaMetrica titulo={t.costoPorcion} valor={d.cost_per_serving_cents != null ? dinero(d.cost_per_serving_cents) : '—'} />
            <TarjetaMetrica titulo={t.kcalTotal} valor={kcal(d.total_kcal)} />
            <TarjetaMetrica titulo={t.kcalPorcion} valor={kcal(d.kcal_per_serving)} />
          </div>
          <Seccion titulo={t.ingredientes}>
            <VistaBD def={{
              clave: 'P-12-ingredientes', filas: d.ingredients, idFila: (i) => i.recipe_ingredient_id, vistas: ['tabla'], titulo: (i) => productos.get(i.product_id)?.name ?? '—',
              propiedades: [
                { id: 'producto', nombre: t.producto, valor: (i) => i.display_text ?? productos.get(i.product_id)?.name ?? '—', pintar: (i) => <span>{i.display_text ?? productos.get(i.product_id)?.name ?? '—'}{i.is_optional && <span className="text-text-secondary"> · {t.opcional}</span>}</span> },
                { id: 'cantidad', nombre: t.cantidad, valor: (i) => i.quantity_milli, pintar: (i) => cantidad(i.quantity_milli, productos.get(i.product_id)?.base_unit), alinear: 'derecha' },
                { id: 'comprar', nombre: t.comprar, valor: (i) => i.packages_to_buy ?? null, pintar: (i) => (i.packages_to_buy ? t.piezas(i.packages_to_buy) : '—'), alinear: 'derecha' },
                { id: 'costo', nombre: t.costo, valor: (i) => i.cost_cents ?? null, alinear: 'derecha',
                  pintar: (i) => (i.cost_cents != null ? <span className="flex flex-col items-end"><span className="tabular-nums">{dinero(i.cost_cents)}</span>{i.price_store_id && i.price_observed_date && <span className="text-caption text-text-secondary">{t.precioDe(tiendas.get(i.price_store_id)?.name ?? '', fecha(i.price_observed_date))}</span>}</span> : <span title={t.sinPrecioMotivo}>— <span className="text-caption text-text-secondary">{t.sinPrecio}</span></span>) },
                { id: 'kcal', nombre: t.kcal, valor: (i) => i.kcal ?? null, alinear: 'derecha', pintar: (i) => (i.kcal != null ? `${kcal(i.kcal)} · ${i.nutrition_source ? t.fuentesNutricion[i.nutrition_source] : t.sinFuente}` : t.sinFuente) },
              ],
            }} />
          </Seccion>
          <Seccion titulo={t.pasos}>
            {r.instructions?.trim() ? (
              <ol className="flex list-decimal flex-col gap-2 pl-6">{r.instructions.split(/\n\s*\n/).map((p, i) => <li key={i} className="whitespace-pre-line text-body">{p.trim()}</li>)}</ol>
            ) : <p className="text-small text-text-secondary">{t.sinPasos}</p>}
            {r.notes && <Franja rol="info">{r.notes}</Franja>}
          </Seccion>
          {hoja === 'editar' && (
            <Hoja abierta alCambiar={(a) => !a && setHoja(null)} titulo={t.editarReceta}>
              <FormReceta inicial={aCuerpo(d)} textoGuardar={textos.app.guardar} alGuardar={async (c) => {
                await pedir('patch', '/api/v1/recipes/{recipe_id}', { ruta: { recipe_id: id }, cuerpo: c })
                invalidar('recetas'); avisar({ texto: t.guardado, rol: 'exito' }); setHoja(null); det.recargar()
              }} />
            </Hoja>
          )}
          {hoja === 'eliminar' && (
            <Hoja abierta alCambiar={(a) => !a && setHoja(null)} titulo={t.eliminarReceta}
              pie={<div className="flex gap-2"><Boton tam="hoja" onClick={() => setHoja(null)}>{textos.app.cancelar}</Boton>
                <Boton tam="hoja" variante="destructivo" deshabilitadoPor={sinRed} onClick={async () => {
                  try { await pedir('delete', '/api/v1/recipes/{recipe_id}', { ruta: { recipe_id: id }, consulta: planeada ? { confirm_planned: true } : {} }); invalidar('recetas'); avisar({ texto: t.eliminado, rol: 'exito' }); navegar(rutaDe('P-11')) } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
                }}>{t.eliminar}</Boton></div>}>
              <p>{t.confirmarEliminar(r.name)}</p>
              {planeada && <Franja rol="advertencia">{t.enPlanes(fecha(planeada))}</Franja>}
            </Hoja>
          )}
        </div>
      )
    }}</ConDatos>
  )
}

const aCuerpo = (d: Detalle): Cuerpo => ({
  name: d.recipe.name, base_servings: d.recipe.base_servings, prep_minutes: d.recipe.prep_minutes, cook_minutes: d.recipe.cook_minutes, instructions: d.recipe.instructions, notes: d.recipe.notes, meal_slot_ids: d.recipe.meal_slot_ids,
  ingredients: d.ingredients.map((i, n) => ({ product_id: i.product_id, presentation_id: i.presentation_id, quantity_milli: Math.round((i.quantity_milli * d.recipe.base_servings * 1000) / d.servings_milli), display_text: i.display_text, is_optional: i.is_optional, sort_order: n + 1 })),
})

interface Renglon { product_id: string; cantidad: string; display_text: string; presentation_id?: number; is_optional?: boolean }
function FormReceta({ inicial, textoGuardar, alGuardar }: { inicial: Cuerpo; textoGuardar: string; alGuardar: (c: Cuerpo) => Promise<void> }) {
  const conexion = useConexion()
  const productos = usarCatalogo('productos') ?? []
  const tiempos = usarCatalogo('tiempos') ?? []
  const [nombre, setNombre] = useState(inicial.name)
  const [slots, setSlots] = useState<number[]>(inicial.meal_slot_ids ?? [])
  const [base, setBase] = useState(String(inicial.base_servings ?? 2))
  const [prep, setPrep] = useState(inicial.prep_minutes != null ? String(inicial.prep_minutes) : '')
  const [coc, setCoc] = useState(inicial.cook_minutes != null ? String(inicial.cook_minutes) : '')
  const [pasos, setPasos] = useState(inicial.instructions ?? '')
  const [notas, setNotas] = useState(inicial.notes ?? '')
  const [ings, setIngs] = useState<Renglon[]>((inicial.ingredients ?? []).map((i) => ({ product_id: String(i.product_id), cantidad: String(i.quantity_milli / 1000), display_text: i.display_text ?? '', presentation_id: i.presentation_id, is_optional: i.is_optional })))
  const [error, setError] = useState<string | null>(null)
  const [ocupado, setOcupado] = useState(false)
  const cambiar = (n: number, c: Partial<Renglon>) => setIngs((l) => l.map((x, i) => (i === n ? { ...x, ...c } : x)))
  async function enviar() {
    setError(null)
    if (!nombre.trim()) return setError(t.faltaNombre)
    const lista = ings.filter((i) => i.product_id).map((i, n) => ({ product_id: Number(i.product_id), quantity_milli: aMilli(i.cantidad), display_text: i.display_text.trim() || undefined, presentation_id: i.presentation_id, is_optional: i.is_optional, sort_order: n + 1 }))
    if (lista.some((i) => i.quantity_milli === null)) return setError(t.cantidadInvalida)
    setOcupado(true)
    try {
      await alGuardar({ name: nombre.trim(), base_servings: Number(base) || 2, prep_minutes: prep ? Number(prep) : undefined, cook_minutes: coc ? Number(coc) : undefined, instructions: pasos.trim() || undefined, notes: notas.trim() || undefined, meal_slot_ids: slots,
        ingredients: lista.map((i) => ({ ...i, quantity_milli: i.quantity_milli as number })) })
    } catch (e) { setError(mensajeError(e)) } finally { setOcupado(false) }
  }
  return (
    <form className="flex flex-col gap-3" onSubmit={(e) => { e.preventDefault(); void enviar() }}>
      <Campo etiqueta={t.nombre} value={nombre} onChange={(e) => setNombre(e.target.value)} required />
      <fieldset className="flex flex-col gap-1"><legend className="text-small font-medium text-text-secondary">{t.tiempos}</legend>
        <div className="flex flex-wrap gap-2">{tiempos.map((x) => <ChipSeleccionable key={x.id} seleccionado={slots.includes(x.id)} onClick={() => setSlots((s) => (s.includes(x.id) ? s.filter((y) => y !== x.id) : [...s, x.id]))}>{x.name}</ChipSeleccionable>)}</div>
      </fieldset>
      <div className="grid grid-cols-3 gap-2">
        <Campo etiqueta={t.porcionesBase} inputMode="numeric" value={base} onChange={(e) => setBase(e.target.value)} />
        <Campo etiqueta={t.prepMin} inputMode="numeric" value={prep} onChange={(e) => setPrep(e.target.value)} />
        <Campo etiqueta={t.coccionMin} inputMode="numeric" value={coc} onChange={(e) => setCoc(e.target.value)} />
      </div>
      <fieldset className="flex flex-col gap-2"><legend className="text-small font-medium text-text-secondary">{t.ingredientes}</legend>
        {ings.map((i, n) => {
          const p = productos.find((x) => String(x.id) === i.product_id)
          return (
            <div key={n} className="grid grid-cols-2 items-end gap-2 rounded-md border border-border-divider p-2 tablet:grid-cols-4">
              <label className="col-span-2 flex flex-col gap-1 text-small font-medium text-text-secondary">{t.producto}
                <select className={`${estilosControl} coarse:h-11`} value={i.product_id} onChange={(e) => cambiar(n, { product_id: e.target.value })}>
                  <option value="">{t.elegirProducto}</option>
                  {productos.map((x) => <option key={x.id} value={x.id}>{x.name}</option>)}
                </select>
              </label>
              <Campo etiqueta={`${t.cantidad}${p ? ` (${t.unidades[p.base_unit]})` : ''}`} inputMode="decimal" value={i.cantidad} onChange={(e) => cambiar(n, { cantidad: e.target.value })} />
              <div className="flex items-end gap-1">
                <Campo className="flex-1" etiqueta={t.textoVisible} value={i.display_text} onChange={(e) => cambiar(n, { display_text: e.target.value })} />
                <Boton tam="icono" variante="fantasma" aria-label={t.quitarIngrediente} onClick={() => setIngs((l) => l.filter((_, k) => k !== n))}><Trash2 aria-hidden className="size-6" size={ICONO_MD} /></Boton>
              </div>
            </div>
          )
        })}
        <Boton onClick={() => setIngs((l) => [...l, { product_id: '', cantidad: '', display_text: '' }])}><Plus aria-hidden className="size-5" />{t.agregarIngrediente}</Boton>
      </fieldset>
      <label className="flex flex-col gap-1 text-small font-medium text-text-secondary">{t.instrucciones}
        <textarea rows={6} value={pasos} onChange={(e) => setPasos(e.target.value)} className="rounded-md border border-border-control bg-surface-field p-3 text-body text-text-default focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus" />
      </label>
      <Campo etiqueta={t.notas} value={notas} onChange={(e) => setNotas(e.target.value)} />
      {error && <Franja rol="error" alerta>{error}</Franja>}
      <Boton type="submit" tam="hoja" variante="primario" disabled={ocupado} deshabilitadoPor={conexion ? undefined : textos.conexion.requiere}>{textoGuardar}</Boton>
    </form>
  )
}
