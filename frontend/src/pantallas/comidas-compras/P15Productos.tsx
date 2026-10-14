// P-15 Productos y precios (U1 §2.1–§2.2, §3.17 N-170 a N-176). Los precios de PROFECO no se editan (N-176).
// WF-3a bis: N-173 con listNutritionFacts y N-175 con edición en el lugar (updateProduct, updatePresentation y `Ctrl+/`).
import { useDeferredValue, useState } from 'react'
import { Package } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import type { PropsPantalla } from '../../app/registro'
import type { Esquemas } from '../../api/tipos'
import { nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { useConexion } from '../../api/conexion'
import { aCentavos, dinero, fecha, hoyISO } from '../../api/formatos'
import { conConsulta, navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Franja, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { Grafica } from '../../componentes/Grafica'
import { useAvisos } from '../../componentes/Avisos'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { CeldaEditable } from '../../vistas-bd/CeldaEditable'
import type { Propiedad } from '../../vistas-bd/tipos'
import { BarraAcciones, MenuMas, Seccion, Selector, Tiendita, aMilli, cantidad, filasDe, invalidar, mensajeError, numero, porId, usarCatalogo, nombreTienda } from './comun'
import { textosCC as t } from './textos'

type Producto = Esquemas['Product']
type Precio = Esquemas['Price']
type Presentacion = Esquemas['Presentation']
/** C1: base_unit_price_cents es el precio por 1 kg, 1 l o 1 pieza (K7); el contenido viene en milésimas de g, ml o pieza. */
const MILLI_POR_UNIDAD_PRECIO: Record<string, number> = { g: 1_000_000, ml: 1_000_000, piece: 1_000 }
/** a × b ÷ c con HALF_UP en enteros (§6: centavos enteros, nunca flotantes en el resultado). */
const mulDivHalfUp = (a: number, b: number, c: number) => Math.floor((2 * a * b + c) / (2 * c))
const EDITABLES = ['manual', 'ticket_manual', 'ticket_ocr']
const opc = (r: Record<string, string>) => Object.entries(r).map(([valor, texto]) => ({ valor, texto }))

export default function P15Productos({ sub, consulta }: PropsPantalla) {
  const conexion = useConexion()
  const tiendas = porId(usarCatalogo('tiendas'))
  const gruposL = usarCatalogo('grupos')
  const grupos = porId(gruposL)
  const [f, setF] = useState({ q: '', tienda: '', grupo: '', frescura: consulta.get('frescura') ?? '', fuente: '' })
  const q = useDeferredValue(f.q.trim())
  const activa = sub === 'detalle'
  const productos = usarConsulta('get', '/api/v1/products', { consulta: { q: q || undefined, store_id: f.tienda ? [Number(f.tienda)] : undefined, product_group_id: f.grupo ? Number(f.grupo) : undefined, freshness: (f.frescura || undefined) as Producto['freshness'], source: (f.fuente || undefined) as Precio['source'] | undefined, limit: 200 } }, activa)
  const resumen = usarConsulta('get', '/api/v1/prices/summary', {}, sub === 'resumen')
  const sinRed = conexion ? undefined : textos.conexion.requiere
  const hoja = consulta.get('hoja')
  const cerrar = () => navegar(conConsulta({ hoja: null, id: null }), { reemplazar: true })
  const mejor = (p: Producto) => p.best_price
  const opcionesGrupo = [{ valor: '', texto: t.sinCategoria }, ...(gruposL ?? []).map((g) => ({ valor: String(g.id), texto: g.name }))]
  const editarProducto = async (p: Producto, cuerpo: Esquemas['ProductUpdate']) => {
    await pedir('patch', '/api/v1/products/{product_id}', { ruta: { product_id: p.id }, cuerpo })
    invalidar('productos'); productos.recargar()
  }
  const grupoDe = (p: Producto) => (p.product_group_id ? String(p.product_group_id) : '')
  return (
    <div className="flex flex-col gap-4">
      <BarraAcciones acciones={<>
        <Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => navegar(conConsulta({ hoja: 'precio' }))}>{t.capturarPrecio}</Boton>
        <MenuMas acciones={[{ texto: t.fotoTicket, alHacer: () => navegar(rutaDe('P-05')) }, { texto: t.nuevoProducto, alHacer: () => navegar(conConsulta({ hoja: 'nuevo' })), deshabilitadoPor: sinRed }]} />
      </>} />
      {sub === 'resumen' && <ConDatos consulta={resumen}>{(r) => (
        <div className="flex flex-col gap-4">
          <div className="grid gap-3 tablet:grid-cols-3">
            {([['fresh', t.conPrecioFresco, r.fresh_count], ['stale', t.conPrecioViejo, r.stale_count], ['no_price', t.sinPrecioCuenta, r.no_price_count]] as const).map(([k, titulo, n]) =>
              <TarjetaMetrica key={k} titulo={titulo} valor={numero(n)} nota={t.frescuras[k]} rol={k === 'fresh' ? undefined : 'warning'} alAbrir={() => navegar(conConsulta({ sub: null, frescura: k }))} />)}
          </div>
          <Seccion titulo={t.ultimaFecha}>
            <ul className="flex flex-col gap-1">{r.stores.map((s) => <li key={s.store_id} className="flex items-center justify-between text-small"><Tiendita id={s.store_id} mapa={tiendas} /><span>{s.latest_observed_date ? fecha(s.latest_observed_date) : t.sinPrecio}</span></li>)}</ul>
          </Seccion>
        </div>
      )}</ConDatos>}
      {sub === 'tendencias' && <Historial />}
      {activa && <>
        <div className="grid gap-2 tablet:grid-cols-3 desktop:grid-cols-5">
          <Campo etiqueta={t.buscar} type="search" value={f.q} onChange={(e) => setF({ ...f, q: e.target.value })} />
          <Selector id="p15-tienda" etiqueta={t.tienda} valor={f.tienda} vacio={t.todas} alCambiar={(v) => setF({ ...f, tienda: v })} opciones={[...tiendas.values()].map((s) => ({ valor: String(s.id), texto: s.name }))} />
          <Selector id="p15-grupo" etiqueta={t.categoria} valor={f.grupo} vacio={t.todas} alCambiar={(v) => setF({ ...f, grupo: v })} opciones={(gruposL ?? []).map((g) => ({ valor: String(g.id), texto: g.name }))} />
          <Selector id="p15-frescura" etiqueta={t.frescura} valor={f.frescura} vacio={t.todas} alCambiar={(v) => setF({ ...f, frescura: v })} opciones={opc(t.frescuras)} />
          <Selector id="p15-fuente" etiqueta={t.fuente} valor={f.fuente} vacio={t.todas} alCambiar={(v) => setF({ ...f, fuente: v })} opciones={opc(t.fuentes)} />
        </div>
        <ConDatos consulta={productos}>{(d) => {
          const filas = filasDe<Producto>(d)
          if (!filas.length) return <Vacio icono={<Package aria-hidden className="size-6" size={ICONO_MD} />} texto={t.sinProductos} />
          return <VistaBD def={{
            clave: 'P-15-productos', filas, idFila: (p) => p.id, vistas: ['tabla', 'tablero', 'galeria'], titulo: (p) => p.name, alAbrir: (p) => navegar(conConsulta({ hoja: 'producto', id: String(p.id) })),
            columnasTablero: { propiedad: 'frescura', columnas: Object.values(t.frescuras) },
            propiedades: [
              { id: 'nombre', nombre: t.nombre, valor: (p) => p.name }, // abre la ficha; el nombre se edita en el lugar dentro de ella
              { id: 'categoria', nombre: t.categoria, valor: (p) => grupos.get(p.product_group_id ?? 0)?.name ?? '', pintar: (p) => grupos.get(p.product_group_id ?? 0)?.name ?? t.sinCategoria,
                editar: { tipo: 'opcion', opciones: opcionesGrupo, textoInicial: grupoDe, guardar: (p, v) => editarProducto(p, { product_group_id: v ? Number(v) : null }) } },
              { id: 'unidad', nombre: t.unidadBase, valor: (p) => t.unidadesNombre[p.base_unit] },
              { id: 'mejor', nombre: t.mejorPrecio, valor: (p) => mejor(p)?.base_unit_price_cents ?? null, alinear: 'derecha',
                pintar: (p) => { const b = mejor(p); return b ? <span className="flex flex-col items-end"><span className="tabular-nums">{t.porUnidad(dinero(b.base_unit_price_cents), t.unidadPrecio[p.base_unit])}</span><Tiendita id={b.store_id} mapa={tiendas} /></span> : '—' } },
              { id: 'ultimo', nombre: t.ultimoPrecio, valor: (p) => p.last_observed_date ?? null, pintar: (p) => (p.last_observed_date ? fecha(p.last_observed_date) : '—') },
              { id: 'frescura', nombre: t.frescura, valor: (p) => t.frescuras[p.freshness ?? 'no_price'], grupo: (p) => t.frescuras[p.freshness ?? 'no_price'] },
            ],
            // Edición en grupo (U1 §5.3 `Ctrl+/`): categoría de varios productos, con vista previa sumada antes de aplicar.
            editarVarias: {
              propiedades: [{ id: 'product_group_id', nombre: t.categoria, tipo: 'opcion', opciones: opcionesGrupo }],
              vistaPrevia: async (sel, _prop, v) => [{ etiqueta: t.cambianCategoria, antes: numero(0), despues: numero(sel.filter((p) => grupoDe(p) !== v).length) }],
              guardar: async (sel, _prop, v) => {
                for (const p of sel.filter((x) => grupoDe(x) !== v)) await pedir('patch', '/api/v1/products/{product_id}', { ruta: { product_id: p.id }, cuerpo: { product_group_id: v ? Number(v) : null } })
                invalidar('productos'); productos.recargar()
              },
            },
          }} />
        }}</ConDatos>
      </>}
      {hoja === 'producto' && consulta.get('id') && <HojaProducto id={Number(consulta.get('id'))} cerrar={cerrar} alCambiar={productos.recargar} />}
      {hoja === 'precio' && <HojaPrecio cerrar={cerrar} alGuardar={productos.recargar} />}
      {hoja === 'nuevo' && <HojaNuevo cerrar={cerrar} alGuardar={productos.recargar} />}
    </div>
  )
}

function Historial() {
  const productos = usarCatalogo('productos') ?? []
  const tiendas = porId(usarCatalogo('tiendas'))
  const [id, setId] = useState('')
  const precios = usarConsulta('get', '/api/v1/products/{product_id}/prices', { ruta: { product_id: Number(id) || 0 }, consulta: { sort: 'observed_date', limit: 200 } }, !!id)
  const p = productos.find((x) => String(x.id) === id)
  return (
    <div className="flex flex-col gap-3">
      <Selector id="p15-historial" etiqueta={t.producto} valor={id} alCambiar={setId} vacio={t.elegirProducto} opciones={productos.map((x) => ({ valor: String(x.id), texto: x.name }))} />
      {!id ? <Vacio texto={t.elegirParaHistorial} /> : <ConDatos consulta={precios}>{(d) => {
        const l = filasDe<Precio>(d).filter((x) => x.base_unit_price_cents != null)
        if (!l.length) return <Vacio texto={textos.vacio.grafica} />
        const fechas = [...new Set(l.map((x) => x.observed_date))].sort()
        const ids = [...new Set(l.map((x) => x.store_id))]
        // Precio vigente por tienda: el último capturado hasta cada fecha (se mantiene hasta que cambia)
        const series = ids.map((s) => { let ult = 0; return { nombre: nombreTienda(tiendas.get(s), tiendas), valores: fechas.map((fe) => { const x = l.filter((y) => y.store_id === s && y.observed_date === fe).pop(); if (x) ult = x.base_unit_price_cents ?? ult; return ult }) } })
        return <Grafica titulo={t.historial} resumen={t.resumenHistorial(p?.name ?? '')} tipo="lineas" categorias={fechas.map(fecha)} series={series} />
      }}</ConDatos>}
    </div>
  )
}

/** Esc dentro de una edición en el lugar descarta solo esa edición (U4 §0.2), no cierra la hoja: Radix escucha Esc antes
 *  que la celda, así que la hoja ignora el cierre mientras el foco está en un control de `[data-edicion-en-lugar]`. */
const editandoEnLugar = () => { const el = document.activeElement; return el instanceof HTMLElement && /^(INPUT|SELECT)$/.test(el.tagName) && !!el.closest('[data-edicion-en-lugar]') }

function HojaProducto({ id, cerrar, alCambiar }: { id: number; cerrar: () => void; alCambiar: () => void }) {
  const avisar = useAvisos()
  const conexion = useConexion()
  const tiendas = porId(usarCatalogo('tiendas'))
  const det = usarConsulta('get', '/api/v1/products/{product_id}', { ruta: { product_id: id } })
  const precios = usarConsulta('get', '/api/v1/products/{product_id}/prices', { ruta: { product_id: id }, consulta: { sort: '-observed_date', limit: 100 } })
  const nutricion = usarConsulta('get', '/api/v1/products/{product_id}/nutrition-facts', { ruta: { product_id: id }, consulta: { limit: 50 } })
  const [corrige, setCorrige] = useState<{ id: number; precio: string; confirmarBorrar?: boolean } | null>(null)
  const [nut, setNut] = useState({ kcal: '', cita: '' })
  const [pres, setPres] = useState({ nombre: '', tam: '' })
  const sinRed = conexion ? undefined : textos.conexion.requiere
  const hacer = async (f: () => Promise<unknown>) => { try { await f(); avisar({ texto: t.guardado, rol: 'exito' }); setCorrige(null); det.recargar(); precios.recargar(); nutricion.recargar() } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) } }
  // N-175: edición en el lugar de la presentación (CeldaEditable guarda con Enter o al salir; Esc descarta)
  const editarPres = async (x: Presentacion, cuerpo: Esquemas['PresentationUpdate']) => {
    await pedir('patch', '/api/v1/presentations/{presentation_id}', { ruta: { presentation_id: x.id }, cuerpo })
    det.recargar(); precios.recargar(); invalidar('productos'); alCambiar()
  }
  const propNombre: Propiedad<Esquemas['ProductDetail']> = { id: 'nombre', nombre: t.nombre, valor: (x) => x.name,
    editar: { tipo: 'texto', textoInicial: (x) => x.name, guardar: async (x, v) => {
      await pedir('patch', '/api/v1/products/{product_id}', { ruta: { product_id: x.id }, cuerpo: { name: v.trim() } })
      det.recargar(); invalidar('productos'); alCambiar()
    } } }
  return (
    <Hoja abierta alCambiar={(a) => !a && !editandoEnLugar() && cerrar()} titulo={det.datos?.name ?? t.producto}>
      <ConDatos consulta={det}>{(p) => {
        const nombreP = (pid: number) => p.presentations.find((x) => x.id === pid)?.name ?? '—'
        return (
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-1" data-edicion-en-lugar><span className="text-small font-medium text-text-secondary">{t.nombre}</span>
              <CeldaEditable fila={p} prop={propNombre}><span className="font-medium">{p.name}</span></CeldaEditable></div>
            <Seccion titulo={t.presentaciones}>
              {p.presentations.length > 0 && <div data-edicion-en-lugar><VistaBD def={{
                clave: 'P-15-presentaciones', filas: p.presentations, idFila: (x) => x.id, vistas: ['tabla'], titulo: (x) => x.name,
                propiedades: [
                  { id: 'nombre', nombre: t.presentacion, valor: (x) => x.name, editar: { tipo: 'texto', textoInicial: (x) => x.name, guardar: (x, v) => editarPres(x, { name: v.trim() }) } },
                  { id: 'marca', nombre: t.marca, valor: (x) => x.brand ?? '', pintar: (x) => x.brand || '—', editar: { tipo: 'texto', textoInicial: (x) => x.brand ?? '', guardar: (x, v) => editarPres(x, { brand: v.trim() || null }) } },
                  { id: 'tamano', nombre: `${t.tamano} (${t.unidades[p.base_unit]})`, valor: (x) => x.size_milli, alinear: 'derecha',
                    pintar: (x) => <span className="tabular-nums">{cantidad(x.size_milli, p.base_unit)}{x.is_bulk ? ` · ${t.granel}` : ''}</span>,
                    editar: { tipo: 'numero', textoInicial: (x) => String(x.size_milli / 1000),
                      // El tamaño cambia el precio por kg, l o pieza: «Antes → Después» con el último precio de la presentación
                      vistaPrevia: async (x, v) => {
                        const nuevo = aMilli(v)
                        if (nuevo === null) throw new Error(t.cantidadInvalida)
                        const suyos = filasDe<Precio>(precios.datos ?? { items: [] }).filter((y) => y.presentation_id === x.id)
                        const cambios = [{ etiqueta: t.tamano, antes: cantidad(x.size_milli, p.base_unit), despues: cantidad(nuevo, p.base_unit) }]
                        const u = suyos[0]
                        if (u && u.base_unit_price_cents != null) cambios.push({ etiqueta: t.ultimoPrecioUnidad(t.unidadPrecio[p.base_unit], nombreTienda(tiendas.get(u.store_id), tiendas), fecha(u.observed_date)), antes: dinero(u.base_unit_price_cents), despues: dinero(mulDivHalfUp(u.price_cents, MILLI_POR_UNIDAD_PRECIO[p.base_unit], nuevo)) })
                        if (suyos.length) cambios.push({ etiqueta: t.preciosSeRecalculan, antes: numero(0), despues: numero(suyos.length) })
                        return cambios
                      },
                      guardar: (x, v) => editarPres(x, { size_milli: aMilli(v)! }) } },
                ],
              }} /></div>}
              <div className="grid grid-cols-2 items-end gap-2">
                <Campo etiqueta={t.nuevaPresentacion} value={pres.nombre} onChange={(e) => setPres({ ...pres, nombre: e.target.value })} />
                <Campo etiqueta={`${t.tamano} (${t.unidades[p.base_unit]})`} inputMode="decimal" value={pres.tam} onChange={(e) => setPres({ ...pres, tam: e.target.value })} />
              </div>
              <Boton deshabilitadoPor={sinRed} disabled={!pres.nombre.trim() || aMilli(pres.tam) === null} onClick={() => hacer(() => pedir('post', '/api/v1/products/{product_id}/presentations', { ruta: { product_id: id }, cuerpo: { name: pres.nombre.trim(), size_milli: aMilli(pres.tam)! }, idempotencia: nuevaClave() }).then(() => setPres({ nombre: '', tam: '' })))}>{t.nuevaPresentacion}</Boton>
            </Seccion>
            <Seccion titulo={t.precios}>
              <ConDatos consulta={precios}>{(d) => (
                <ul className="flex flex-col divide-y divide-border-divider">{filasDe<Precio>(d).map((x) => (
                  <li key={x.id} className="flex flex-col gap-1 py-2">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <Tiendita id={x.store_id} mapa={tiendas} />
                      <span className="tabular-nums font-medium">{dinero(x.price_cents)}{x.is_stale && <span className="ml-1 rounded-sm bg-role-warning-surface px-1 text-caption text-role-warning">{t.viejo}</span>}</span>
                    </div>
                    <span className="text-caption tabular-nums text-text-secondary">{[nombreP(x.presentation_id), x.base_unit_price_cents != null ? t.porUnidad(dinero(x.base_unit_price_cents), t.unidadPrecio[p.base_unit]) : '', fecha(x.observed_date), t.fuentes[x.source], x.age_days != null ? t.diasN(x.age_days) : ''].filter(Boolean).join(' · ')}</span>
                    {!EDITABLES.includes(x.source) ? <span className="text-caption text-text-secondary">{t.profecoNoSeEdita}</span> : corrige?.id === x.id ? (
                      corrige.confirmarBorrar
                        ? <div className="flex flex-wrap items-center gap-2"><span className="text-small">{t.confirmarEliminarPrecio}</span><Boton onClick={() => setCorrige(null)}>{textos.app.cancelar}</Boton><Boton variante="destructivo" deshabilitadoPor={sinRed} onClick={() => hacer(() => pedir('delete', '/api/v1/prices/{price_id}', { ruta: { price_id: x.id } }))}>{t.eliminar}</Boton></div>
                        : <div className="flex flex-wrap items-end gap-2"><Campo etiqueta={t.precio} inputMode="decimal" value={corrige.precio} onChange={(e) => setCorrige({ ...corrige, precio: e.target.value })} ayuda={aCentavos(corrige.precio) ? `${dinero(x.price_cents)} → ${dinero(aCentavos(corrige.precio)!)}` : undefined} />
                          <Boton onClick={() => setCorrige(null)}>{textos.app.cancelar}</Boton>
                          <Boton variante="primario" deshabilitadoPor={sinRed} disabled={!aCentavos(corrige.precio)} onClick={() => hacer(() => pedir('patch', '/api/v1/prices/{price_id}', { ruta: { price_id: x.id }, cuerpo: { price_cents: aCentavos(corrige.precio)! } }))}>{textos.app.guardar}</Boton></div>
                    ) : <div className="flex gap-2"><Boton variante="enlace" onClick={() => setCorrige({ id: x.id, precio: String(x.price_cents / 100) })}>{t.corregir}</Boton><Boton variante="enlace" onClick={() => setCorrige({ id: x.id, precio: '', confirmarBorrar: true })}>{t.eliminar}</Boton></div>}
                  </li>
                ))}</ul>
              )}</ConDatos>
            </Seccion>
            <Seccion titulo={t.nutricion}>
              <ConDatos consulta={nutricion}>{(dn) => {
                const datos = filasDe<Esquemas['NutritionFact']>(dn)
                return datos.length ? <ul className="flex flex-col gap-1 text-small" data-prueba="nutricion">{datos.map((n) => <li key={n.id}><span className="tabular-nums">{t.kcal100(numero(n.kcal_milli_per_100g / 1000), p.base_unit === 'ml' ? 'ml' : 'g')}</span> · {t.fuentesNutricion[n.source]}{n.source_ref ? ` · ${n.source_ref}` : ''}{n.is_preferred ? ` · ${t.preferida}` : ''}</li>)}</ul> : <p className="text-small text-text-secondary">{t.sinNutricion}</p>
              }}</ConDatos>
              <div className="grid grid-cols-2 items-end gap-2">
                <Campo etiqueta={t.kcalPor100} inputMode="decimal" value={nut.kcal} onChange={(e) => setNut({ ...nut, kcal: e.target.value })} />
                <Campo etiqueta={t.cita} value={nut.cita} onChange={(e) => setNut({ ...nut, cita: e.target.value })} />
              </div>
              <Boton deshabilitadoPor={sinRed} disabled={aMilli(nut.kcal) === null || !nut.cita.trim()} onClick={() => hacer(() => pedir('post', '/api/v1/products/{product_id}/nutrition-facts', { ruta: { product_id: id }, cuerpo: { kcal_milli_per_100g: aMilli(nut.kcal)!, source: 'manual', source_ref: nut.cita.trim() }, idempotencia: nuevaClave() }).then(() => setNut({ kcal: '', cita: '' })))}>{t.agregarNutricion}</Boton>
            </Seccion>
          </div>
        )
      }}</ConDatos>
    </Hoja>
  )
}

function HojaPrecio({ cerrar, alGuardar }: { cerrar: () => void; alGuardar: () => void }) {
  const avisar = useAvisos()
  const conexion = useConexion()
  const productos = usarCatalogo('productos') ?? []
  const tiendas = usarCatalogo('tiendas') ?? []
  const [v, setV] = useState({ producto: '', presentacion: '', tienda: '', precio: '', fecha: hoyISO() })
  const det = usarConsulta('get', '/api/v1/products/{product_id}', { ruta: { product_id: Number(v.producto) || 0 } }, !!v.producto)
  const [error, setError] = useState<string | null>(null)
  async function guardar() {
    const c = aCentavos(v.precio)
    if (!v.presentacion || !v.tienda) return setError(t.elegir)
    if (!c) return setError(t.montoInvalido)
    try { await pedir('post', '/api/v1/prices', { cuerpo: { presentation_id: Number(v.presentacion), store_id: Number(v.tienda), price_cents: c, observed_date: v.fecha }, idempotencia: nuevaClave() }); avisar({ texto: t.precioGuardado, rol: 'exito' }); invalidar('productos'); alGuardar(); cerrar() } catch (e) { setError(mensajeError(e)) }
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={t.capturarPrecio} pie={<Boton tam="hoja" variante="primario" deshabilitadoPor={conexion ? undefined : textos.conexion.requiere} onClick={guardar}>{textos.app.guardar}</Boton>}>
      <div className="flex flex-col gap-3">
        <Selector id="p15-p-producto" etiqueta={t.producto} valor={v.producto} vacio={t.elegirProducto} alCambiar={(x) => setV({ ...v, producto: x, presentacion: '' })} opciones={productos.map((p) => ({ valor: String(p.id), texto: p.name }))} />
        <Selector id="p15-p-presentacion" etiqueta={t.presentacion} valor={v.presentacion} vacio={t.elegir} alCambiar={(x) => setV({ ...v, presentacion: x })} opciones={(det.datos?.presentations ?? []).map((p) => ({ valor: String(p.id), texto: `${p.name} · ${cantidad(p.size_milli, det.datos?.base_unit)}` }))} />
        <Selector id="p15-p-tienda" etiqueta={t.tienda} valor={v.tienda} vacio={t.elegir} alCambiar={(x) => setV({ ...v, tienda: x })} opciones={tiendas.filter((s) => s.is_active).map((s) => ({ valor: String(s.id), texto: s.name }))} />
        <Campo etiqueta={t.precio} inputMode="decimal" value={v.precio} onChange={(e) => setV({ ...v, precio: e.target.value })} />
        <Campo etiqueta={t.fecha} type="date" value={v.fecha} onChange={(e) => setV({ ...v, fecha: e.target.value })} ayuda={fecha(v.fecha)} />
        {error && <Franja rol="error" alerta>{error}</Franja>}
      </div>
    </Hoja>
  )
}

function HojaNuevo({ cerrar, alGuardar }: { cerrar: () => void; alGuardar: () => void }) {
  const avisar = useAvisos()
  const conexion = useConexion()
  const grupos = usarCatalogo('grupos') ?? []
  const [v, setV] = useState({ nombre: '', grupo: '', unidad: 'g', pres: '', tam: '' })
  const [error, setError] = useState<string | null>(null)
  async function guardar() {
    if (!v.nombre.trim()) return setError(t.faltaNombre)
    const tam = aMilli(v.tam)
    try {
      await pedir('post', '/api/v1/products', { cuerpo: { name: v.nombre.trim(), base_unit: v.unidad as Producto['base_unit'], product_group_id: v.grupo ? Number(v.grupo) : undefined, presentations: v.pres.trim() && tam ? [{ name: v.pres.trim(), size_milli: tam }] : undefined }, idempotencia: nuevaClave() })
      invalidar('productos'); avisar({ texto: t.productoCreado, rol: 'exito' }); alGuardar(); cerrar()
    } catch (e) { setError(mensajeError(e)) }
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={t.nuevoProducto} pie={<Boton tam="hoja" variante="primario" deshabilitadoPor={conexion ? undefined : textos.conexion.requiere} onClick={guardar}>{textos.app.guardar}</Boton>}>
      <div className="flex flex-col gap-3">
        <Campo etiqueta={t.nombre} value={v.nombre} onChange={(e) => setV({ ...v, nombre: e.target.value })} />
        <Selector id="p15-n-grupo" etiqueta={t.categoria} valor={v.grupo} vacio={t.elegir} alCambiar={(x) => setV({ ...v, grupo: x })} opciones={grupos.map((g) => ({ valor: String(g.id), texto: g.name }))} />
        <Selector id="p15-n-unidad" etiqueta={t.unidadBase} valor={v.unidad} alCambiar={(x) => setV({ ...v, unidad: x })} opciones={opc(t.unidadesNombre)} />
        <div className="grid grid-cols-2 gap-2">
          <Campo etiqueta={t.presentacionInicial} value={v.pres} onChange={(e) => setV({ ...v, pres: e.target.value })} />
          <Campo etiqueta={`${t.tamano} (${t.unidades[v.unidad]})`} inputMode="decimal" value={v.tam} onChange={(e) => setV({ ...v, tam: e.target.value })} />
        </div>
        {error && <Franja rol="error" alerta>{error}</Franja>}
      </div>
    </Hoja>
  )
}
