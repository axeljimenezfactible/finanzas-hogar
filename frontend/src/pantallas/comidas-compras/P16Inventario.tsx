// P-16 Inventario (U1 §2.1–§2.2, §3.18 N-180 a N-184). Estado y días restantes vienen del motor S4 por C1.
import { useDeferredValue, useState } from 'react'
import { Archive } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import type { PropsPantalla } from '../../app/registro'
import type { Esquemas } from '../../api/tipos'
import { nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { useConexion } from '../../api/conexion'
import { fecha, hoyISO } from '../../api/formatos'
import { conConsulta, navegar } from '../../app/enrutador'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Franja, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { useAvisos } from '../../componentes/Avisos'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { BarraAcciones, MenuMas, Selector, aMilli, cantidad, filasDe, mensajeError, numero, porId, usarCatalogo } from './comun'
import { textosCC as t } from './textos'

type Item = Esquemas['InventoryItem']
const ORDEN: Record<string, number> = { out: 0, reorder: 1, sufficient: 2 }
const opc = (r: Record<string, string>) => Object.entries(r).map(([valor, texto]) => ({ valor, texto }))

export default function P16Inventario({ sub, consulta }: PropsPantalla) {
  const conexion = useConexion()
  const avisar = useAvisos()
  const productos = porId(usarCatalogo('productos'))
  const recetas = porId(usarCatalogo('recetas'))
  const [f, setF] = useState({ q: '', modo: '', estado: '' })
  const q = useDeferredValue(f.q.trim())
  const items = usarConsulta('get', '/api/v1/inventory-items', { consulta: { q: q || undefined, tracking_mode: (f.modo || undefined) as Item['tracking_mode'] | undefined, status: (f.estado || undefined) as Item['status'], limit: 200 } }, sub !== 'resumen')
  const resumen = usarConsulta('get', '/api/v1/inventory/summary', {}, sub === 'resumen')
  const sinRed = conexion ? undefined : textos.conexion.requiere
  const hoja = consulta.get('hoja')
  const cerrar = () => navegar(conConsulta({ hoja: null, id: null }), { reemplazar: true })
  const nombre = (i: Item) => (i.recipe_id ? t.sobranteDe(recetas.get(i.recipe_id)?.name ?? '—') : productos.get(i.product_id ?? 0)?.name ?? '—')
  const unidad = (i: Item) => (i.recipe_id ? 'porcion' : productos.get(i.product_id ?? 0)?.base_unit)
  const cant = (m: number | undefined, i: Item) => (unidad(i) === 'porcion' ? (m == null ? '—' : t.porciones(numero(m / 1000))) : cantidad(m, unidad(i)))
  async function ajustar(i: Item, milli: number, razon: 'count' | 'consumption' | 'waste' = 'count') {
    await pedir('post', '/api/v1/inventory-items/{inventory_item_id}/adjustments', { ruta: { inventory_item_id: i.id }, cuerpo: { reason: razon, quantity_milli: milli, counted_date: hoyISO() }, idempotencia: nuevaClave() })
    avisar({ texto: t.existenciaAjustada, rol: 'exito' }); items.recargar()
  }
  return (
    <div className="flex flex-col gap-4">
      <BarraAcciones acciones={<>
        <Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => navegar(conConsulta({ hoja: 'existencia' }))}>{t.ajustarExistencia}</Boton>
        <MenuMas acciones={[{ texto: t.registrarSobrante, alHacer: () => navegar(conConsulta({ hoja: 'sobrante' })), deshabilitadoPor: sinRed }]} />
      </>} />
      {sub === 'resumen' ? <ConDatos consulta={resumen}>{(r) => (
        <div className="grid gap-3 tablet:grid-cols-3">
          {([['out', r.out_count], ['reorder', r.reorder_count], ['sufficient', r.sufficient_count]] as const).map(([k, n]) =>
            <TarjetaMetrica key={k} titulo={t.estados[k]} valor={numero(n)} rol={k === 'sufficient' ? undefined : 'warning'} alAbrir={() => { setF({ ...f, estado: k }); navegar(conConsulta({ sub: null })) }} />)}
        </div>
      )}</ConDatos> : <>
        <div className="grid gap-2 tablet:grid-cols-3">
          <Campo etiqueta={t.buscar} type="search" value={f.q} onChange={(e) => setF({ ...f, q: e.target.value })} />
          <Selector id="p16-modo" etiqueta={t.modo} valor={f.modo} vacio={t.todos} alCambiar={(v) => setF({ ...f, modo: v })} opciones={opc(t.modos)} />
          <Selector id="p16-estado" etiqueta={t.estado} valor={f.estado} vacio={t.todos} alCambiar={(v) => setF({ ...f, estado: v })} opciones={opc(t.estados)} />
        </div>
        <ConDatos consulta={items}>{(d) => {
          const filas = filasDe<Item>(d).sort((a, b) => (ORDEN[a.status ?? 'sufficient'] - ORDEN[b.status ?? 'sufficient']) || nombre(a).localeCompare(nombre(b), 'es-MX'))
          if (!filas.length) return <Vacio icono={<Archive aria-hidden className="size-6" size={ICONO_MD} />} texto={t.inventarioVacio} />
          return <VistaBD def={{
            clave: 'P-16-inventario', filas, idFila: (i) => i.id, vistas: ['tabla', 'tablero', 'galeria'], titulo: nombre, alAbrir: (i) => navegar(conConsulta({ hoja: 'movimientos', id: String(i.id) })),
            columnasTablero: { propiedad: 'estado', columnas: [t.estados.out, t.estados.reorder, t.estados.sufficient] },
            propiedades: [
              { id: 'articulo', nombre: t.articulo, valor: nombre },
              { id: 'modo', nombre: t.modo, valor: (i) => t.modos[i.tracking_mode] },
              { id: 'existencia', nombre: t.existencia, valor: (i) => i.estimated_quantity_milli ?? i.quantity_milli, pintar: (i) => cant(i.estimated_quantity_milli ?? i.quantity_milli, i), alinear: 'derecha',
                editar: { tipo: 'numero', textoInicial: (i) => String(i.quantity_milli / 1000),
                  vistaPrevia: async (i, v) => { const m = aMilli(v); if (m === null) throw new RangeError(t.cantidadInvalida); return [{ etiqueta: t.existencia, antes: cant(i.quantity_milli, i), despues: cant(m, i) }] },
                  guardar: async (i, v) => { const m = aMilli(v); if (m === null) throw new RangeError(t.cantidadInvalida); await ajustar(i, m) } } },
              { id: 'consumo', nombre: t.consumoSemana, valor: (i) => i.weekly_consumption_milli ?? null, pintar: (i) => cant(i.weekly_consumption_milli, i), alinear: 'derecha' },
              { id: 'dias', nombre: t.diasRestantes, valor: (i) => i.days_left ?? null, pintar: (i) => (i.days_left != null ? t.diasN(i.days_left) : '—'), alinear: 'derecha' },
              { id: 'entrada', nombre: t.ultimaEntrada, valor: (i) => i.counted_date ?? null, pintar: (i) => (i.counted_date ? fecha(i.counted_date) : '—') },
              { id: 'estado', nombre: t.estado, valor: (i) => t.estados[i.status ?? 'sufficient'], grupo: (i) => t.estados[i.status ?? 'sufficient'],
                pintar: (i) => <span className={i.status === 'sufficient' ? '' : 'rounded-sm bg-role-warning-surface px-1 text-role-warning'}>{t.estados[i.status ?? 'sufficient']}</span> },
            ],
          }} />
        }}</ConDatos>
      </>}
      {hoja === 'movimientos' && consulta.get('id') && <HojaMovimientos id={Number(consulta.get('id'))} titulo={(() => { const i = filasDe<Item>(items.datos ?? {}).find((x) => x.id === Number(consulta.get('id'))); return i ? nombre(i) : t.movimientos })()} cerrar={cerrar} />}
      {hoja === 'existencia' && <HojaExistencia filas={filasDe<Item>(items.datos ?? {})} nombre={nombre} cerrar={cerrar} alGuardar={ajustar} />}
      {hoja === 'sobrante' && <HojaSobrante cerrar={cerrar} alGuardar={items.recargar} />}
    </div>
  )
}

function HojaMovimientos({ id, titulo, cerrar }: { id: number; titulo: string; cerrar: () => void }) {
  const m = usarConsulta('get', '/api/v1/inventory-items/{inventory_item_id}/movements', { ruta: { inventory_item_id: id }, consulta: { sort: '-movement_date', limit: 100 } })
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={titulo} descripcion={t.movimientos}>
      <ConDatos consulta={m}>{(d) => {
        const l = filasDe<Esquemas['InventoryMovement']>(d)
        if (!l.length) return <Vacio texto={t.sinMovimientos} />
        return <ul className="flex flex-col divide-y divide-border-divider">{l.map((x) => (
          <li key={x.id} className="flex items-center justify-between gap-2 py-2 text-small">
            <span className="flex flex-col"><span className="font-medium">{t.tiposMov[x.kind]}{x.reason ? ` · ${t.razones[x.reason]}` : ''}</span><span className="text-caption text-text-secondary">{fecha(x.movement_date)}{x.notes ? ` · ${x.notes}` : ''}</span></span>
            <span className="tabular-nums">{x.quantity_delta_milli > 0 ? '+' : '−'}{numero(Math.abs(x.quantity_delta_milli) / 1000)}</span>
          </li>
        ))}</ul>
      }}</ConDatos>
    </Hoja>
  )
}

function HojaExistencia({ filas, nombre, cerrar, alGuardar }: { filas: Item[]; nombre: (i: Item) => string; cerrar: () => void; alGuardar: (i: Item, m: number, r: 'count' | 'consumption' | 'waste') => Promise<void> }) {
  const conexion = useConexion()
  const [v, setV] = useState({ item: '', cant: '', razon: 'count' })
  const [error, setError] = useState<string | null>(null)
  const i = filas.find((x) => String(x.id) === v.item)
  async function guardar() {
    const m = aMilli(v.cant)
    if (!i) return setError(t.elegir)
    if (m === null) return setError(t.cantidadInvalida)
    try { await alGuardar(i, m, v.razon as 'count'); cerrar() } catch (e) { setError(mensajeError(e)) }
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={t.ajustarExistencia} pie={<Boton tam="hoja" variante="primario" deshabilitadoPor={conexion ? undefined : textos.conexion.requiere} onClick={guardar}>{textos.app.guardar}</Boton>}>
      <div className="flex flex-col gap-3">
        <Selector id="p16-e-item" etiqueta={t.articulo} valor={v.item} vacio={t.elegir} alCambiar={(x) => setV({ ...v, item: x })} opciones={filas.map((x) => ({ valor: String(x.id), texto: nombre(x) }))} />
        <Campo etiqueta={t.cantidadNueva} inputMode="decimal" value={v.cant} onChange={(e) => setV({ ...v, cant: e.target.value })} />
        <Selector id="p16-e-razon" etiqueta={t.motivo} valor={v.razon} alCambiar={(x) => setV({ ...v, razon: x })} opciones={opc(t.razones)} />
        {error && <Franja rol="error" alerta>{error}</Franja>}
      </div>
    </Hoja>
  )
}

function HojaSobrante({ cerrar, alGuardar }: { cerrar: () => void; alGuardar: () => void }) {
  const avisar = useAvisos()
  const conexion = useConexion()
  const recetas = usarCatalogo('recetas') ?? []
  const [v, setV] = useState({ receta: '', porc: '' })
  const [error, setError] = useState<string | null>(null)
  async function guardar() {
    const m = aMilli(v.porc)
    if (!v.receta) return setError(t.elegir)
    if (m === null) return setError(t.cantidadInvalida)
    try {
      await pedir('post', '/api/v1/inventory-items', { cuerpo: { recipe_id: Number(v.receta), tracking_mode: 'stock', quantity_milli: m, counted_date: hoyISO() }, idempotencia: nuevaClave() })
      avisar({ texto: t.sobranteRegistrado, rol: 'exito' }); alGuardar(); cerrar()
    } catch (e) { setError(mensajeError(e)) }
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={t.registrarSobrante} pie={<Boton tam="hoja" variante="primario" deshabilitadoPor={conexion ? undefined : textos.conexion.requiere} onClick={guardar}>{textos.app.guardar}</Boton>}>
      <div className="flex flex-col gap-3">
        <Selector id="p16-s-receta" etiqueta={t.recetaOrigen} valor={v.receta} vacio={t.elegir} alCambiar={(x) => setV({ ...v, receta: x })} opciones={recetas.map((r) => ({ valor: String(r.id), texto: r.name }))} />
        <Campo etiqueta={t.porcionesSobrantes} inputMode="decimal" value={v.porc} onChange={(e) => setV({ ...v, porc: e.target.value })} />
        {error && <Franja rol="error" alerta>{error}</Franja>}
      </div>
    </Hoja>
  )
}
