// P-13 Lista de compras (U1 §2.1–§2.2, §3.15 N-150 a N-153). Toda lectura y escritura lleva la semana que se ve (week, R3-01).
import { useState } from 'react'
import { ShoppingCart } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import type { PropsPantalla } from '../../app/registro'
import type { Esquemas } from '../../api/tipos'
import { nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { useConexion } from '../../api/conexion'
import { dinero, fecha } from '../../api/formatos'
import { conConsulta, navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Franja, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { useAvisos } from '../../componentes/Avisos'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { BarraAcciones, MenuMas, Seccion, SelectorSemana, Selector, Tiendita, aMilli, cantidad, conSemana, mensajeError, porId, usarCatalogo, useSemana, type Tienda, nombreTienda } from './comun'
import { textosCC as t } from './textos'

type Lista = Esquemas['ShoppingList']
type Grupo = Lista['groups'][number]
type Item = Esquemas['ShoppingListItem'] & { grupo: string }

export default function P13Lista({ consulta }: PropsPantalla) {
  const sem = useSemana(consulta)
  const conexion = useConexion()
  const avisar = useAvisos()
  const lista = usarConsulta('get', '/api/v1/shopping-list', { consulta: { week: sem.semana } })
  const productos = porId(usarCatalogo('productos'))
  const tiendasL = usarCatalogo('tiendas')
  const tiendas = porId(tiendasL)
  const [editar, setEditar] = useState<Item | null>(null)
  const sinRed = conexion ? undefined : textos.conexion.requiere
  const nombre = (i: Esquemas['ShoppingListItem']) => i.free_text ?? productos.get(i.product_id ?? 0)?.name ?? '—'
  const unidad = (i: Esquemas['ShoppingListItem']) => i.free_text_unit ?? productos.get(i.product_id ?? 0)?.base_unit
  const etiqueta = (g: Grupo) => (g.store_id ? nombreTienda(tiendas.get(g.store_id), tiendas) : g.purchase_channel === 'market' ? t.mercados : t.cadenas)
  async function cambiar(i: Esquemas['ShoppingListItem'], cuerpo: Esquemas['ShoppingListItemUpdate']) {
    try { await pedir('patch', '/api/v1/shopping-list/items/{shopping_list_item_id}', { ruta: { shopping_list_item_id: i.id }, consulta: { week: sem.semana }, cuerpo }) }
    catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
    lista.recargar()
  }
  const cerrarHoja = () => navegar(conConsulta({ hoja: null }), { reemplazar: true })
  return (
    <div className="flex flex-col gap-4">
      <BarraAcciones acciones={<>
        <Boton variante="primario" onClick={() => navegar(`${rutaDe('P-05')}?semana=${sem.semana}`)}>{t.terminarCompra}</Boton>
        <MenuMas acciones={[{ texto: t.compararTiendas, alHacer: () => navegar(conSemana(rutaDe('P-14'), sem.semana, sem.porDefecto)) }, { texto: t.agregarArticulo, alHacer: () => navegar(conConsulta({ hoja: 'articulo' })), deshabilitadoPor: sinRed }]} />
      </>}><SelectorSemana {...sem} /></BarraAcciones>
      <ConDatos consulta={lista}>{(l) => {
        const items: Item[] = l.groups.flatMap((g) => g.items.map((i) => ({ ...i, grupo: etiqueta(g) })))
        if (!items.length) return l.generated_at
          ? <Vacio icono={<ShoppingCart aria-hidden className="size-6" size={ICONO_MD} />} texto={t.nadaQueComprar} />
          : <Vacio icono={<ShoppingCart aria-hidden className="size-6" size={ICONO_MD} />} texto={t.sinPlan} accion={<Boton onClick={() => navegar(conSemana(rutaDe('P-10'), sem.semana, sem.porDefecto))}>{t.planear}</Boton>} />
        const comprados = items.filter((i) => i.is_bought).length
        const grupos = l.groups.map((g) => ({ g, nombre: etiqueta(g) }))
        return (
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-1 rounded-md border border-border-divider bg-surface-card p-4">
              <span className="text-small text-text-secondary">{t.totalEstimado}</span>
              <span className="text-title-1 font-semibold tabular-nums" data-prueba="total-lista">{dinero(l.total_cents)}</span>
              <span className="text-small text-text-secondary">{t.resumenLista(items.length, comprados)}</span>
            </div>
            <Seccion titulo={t.totalesTienda}>
              <ul className="grid gap-2 tablet:grid-cols-2 desktop:grid-cols-3">
                {grupos.map(({ g, nombre: n }, k) => (
                  <li key={`${k}-${g.store_id ?? g.purchase_channel}`} className="flex items-center justify-between gap-2 rounded-md border border-border-divider bg-surface-card px-3 py-2">
                    <span className="flex flex-col">
                      {g.store_id ? <Tiendita id={g.store_id} mapa={tiendas} /> : <span className="font-medium">{n}</span>}
                      {g.purchase_channel === 'market' && g.visit_date && <span className="text-caption text-text-secondary">{t.diaVisita(fecha(g.visit_date))}</span>}
                    </span>
                    <span className="tabular-nums font-medium">{dinero(g.total_cents)}</span>
                  </li>
                ))}
              </ul>
            </Seccion>
            <VistaBD def={{
              clave: 'P-13-lista', filas: items, idFila: (i) => i.id, vistas: ['tabla', 'tablero'], titulo: nombre, alAbrir: (i) => setEditar(i),
              columnasTablero: { propiedad: 'tienda', columnas: grupos.map((x) => x.nombre), moverA: async (i, col) => { const g = grupos.find((x) => x.nombre === col)?.g; if (g?.store_id) await cambiar(i, { store_id: g.store_id }) } },
              propiedades: [
                { id: 'comprado', nombre: t.comprado, valor: (i) => i.is_bought, pintar: (i) => (
                  <label className="inline-flex size-11 items-center justify-center" onClick={(e) => e.stopPropagation()}>
                    <input type="checkbox" className="size-5 accent-action-primary-bg" checked={i.is_bought} aria-label={t.marcarComprado(nombre(i))} aria-disabled={sinRed ? true : undefined} onChange={() => !sinRed && cambiar(i, { is_bought: !i.is_bought })} />
                  </label>) },
                { id: 'articulo', nombre: t.articulo, valor: nombre, pintar: (i) => <span className={i.is_bought ? 'text-text-secondary line-through' : 'font-medium'}>{nombre(i)}</span> },
                { id: 'cantidad', nombre: t.cantidad, valor: (i) => i.quantity_milli ?? i.needed_milli ?? null, alinear: 'derecha',
                  pintar: (i) => <span className="flex flex-col items-end tabular-nums"><span>{cantidad(i.quantity_milli ?? i.needed_milli, unidad(i))}</span>{i.packages_to_buy ? <span className="text-caption text-text-secondary">{t.piezas(i.packages_to_buy)}</span> : null}</span> },
                { id: 'tienda', nombre: t.tienda, valor: (i) => i.grupo, grupo: (i) => i.grupo, pintar: (i) => <Tiendita id={i.assigned_store_id ?? i.store_id ?? i.recommended_store_id} mapa={tiendas} /> },
                { id: 'precio', nombre: t.precioEstimado, valor: (i) => i.estimated_price_cents ?? null, alinear: 'derecha',
                  pintar: (i) => (i.estimated_price_cents != null ? <span className="tabular-nums">{dinero(i.estimated_price_cents)}{i.is_stale_price && <span className="ml-1 rounded-sm bg-role-warning-surface px-1 text-caption text-role-warning">{t.viejo}</span>}</span> : '—') },
                { id: 'motivo', nombre: t.motivo, valor: (i) => (i.reasons ?? []).map((r) => t.motivos[r]).join(', ') },
              ],
            }} />
            <Franja rol="info">{t.noCreaGasto}</Franja>
            {editar && <HojaArticulo titulo={t.editarArticulo} item={editar} tiendas={tiendasL ?? []} unidad={unidad(editar)} cerrar={() => setEditar(null)}
              alGuardar={async (c) => { await pedir('patch', '/api/v1/shopping-list/items/{shopping_list_item_id}', { ruta: { shopping_list_item_id: editar.id }, consulta: { week: sem.semana }, cuerpo: { quantity_milli: c.quantity_milli, store_id: c.store_id ?? null } }); lista.recargar() }} />}
          </div>
        )
      }}</ConDatos>
      {consulta.get('hoja') === 'articulo' && <HojaArticulo titulo={t.agregarArticulo} tiendas={tiendasL ?? []} cerrar={cerrarHoja}
        alGuardar={async (c) => { await pedir('post', '/api/v1/shopping-list/items', { consulta: { week: sem.semana }, cuerpo: c, idempotencia: nuevaClave() }); avisar({ texto: t.articuloAgregado, rol: 'exito' }); lista.recargar() }} />}
    </div>
  )
}

function HojaArticulo({ titulo, item, tiendas, unidad, cerrar, alGuardar }: { titulo: string; item?: Item; tiendas: Tienda[]; unidad?: string | null; cerrar: () => void; alGuardar: (c: Esquemas['ShoppingListItemCreate']) => Promise<void> }) {
  const conexion = useConexion()
  const productos = usarCatalogo('productos') ?? []
  const [producto, setProducto] = useState('')
  const [libre, setLibre] = useState('')
  const [uLibre, setULibre] = useState<'g' | 'ml' | 'piece'>('piece')
  const [cant, setCant] = useState(item?.quantity_milli != null ? String(item.quantity_milli / 1000) : item?.needed_milli != null ? String(item.needed_milli / 1000) : '')
  const [tienda, setTienda] = useState(String(item?.assigned_store_id ?? item?.store_id ?? ''))
  const [error, setError] = useState<string | null>(null)
  const u = item ? unidad : producto ? productos.find((p) => String(p.id) === producto)?.base_unit : uLibre
  async function guardar() {
    setError(null)
    const q = aMilli(cant)
    if (q === null) return setError(t.cantidadInvalida)
    if (!item && !producto && !libre.trim()) return setError(t.faltaProducto)
    try {
      await alGuardar({ quantity_milli: q, store_id: tienda ? Number(tienda) : undefined, ...(item ? {} : producto ? { product_id: Number(producto) } : { free_text: libre.trim(), free_text_unit: uLibre }) })
      cerrar()
    } catch (e) { setError(mensajeError(e)) }
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={titulo} pie={<Boton tam="hoja" variante="primario" deshabilitadoPor={conexion ? undefined : textos.conexion.requiere} onClick={guardar}>{textos.app.guardar}</Boton>}>
      <div className="flex flex-col gap-3">
        {item ? <p className="font-medium">{item.grupo}</p> : <>
          <Selector id="p13-producto" etiqueta={t.producto} valor={producto} alCambiar={setProducto} vacio={t.elegirProducto} opciones={productos.map((p) => ({ valor: String(p.id), texto: p.name }))} />
          {!producto && <div className="grid grid-cols-3 gap-2">
            <Campo className="col-span-2" etiqueta={t.textoLibre} value={libre} onChange={(e) => setLibre(e.target.value)} />
            <Selector id="p13-unidad" etiqueta={t.unidad} valor={uLibre} alCambiar={(v) => setULibre(v as 'g' | 'ml' | 'piece')} opciones={Object.entries(t.unidadesNombre).map(([valor, texto]) => ({ valor, texto }))} />
          </div>}
        </>}
        <Campo etiqueta={`${t.cantidad}${u ? ` (${t.unidades[u]})` : ''}`} inputMode="decimal" value={cant} onChange={(e) => setCant(e.target.value)} />
        <Selector id="p13-tienda" etiqueta={t.tienda} valor={tienda} alCambiar={setTienda} vacio={t.sinTienda} opciones={tiendas.filter((s) => s.is_active).map((s) => ({ valor: String(s.id), texto: s.name }))} />
        {error && <Franja rol="error" alerta>{error}</Franja>}
      </div>
    </Hoja>
  )
}
