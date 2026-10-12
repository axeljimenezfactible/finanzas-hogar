// P-26.14 Inventario (U4 parte b; decisión D, K7, RN-15, Q2-70 a). Modo, nivel y ritmo con vista previa en el navegador
// (estado, existencia estimada, días restantes y consumo por semana); la existencia se ajusta en P-16.
import { useEffect, useRef, useState } from 'react'
import { CircleAlert, CircleCheck, TriangleAlert } from 'lucide-react'
import { pedir, nuevaClave, todas } from '../../api/cliente'
import { usarConsulta, type EstadoConsulta } from '../../api/consulta'
import { fecha, hoyISO } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { Campo, estilosControl } from '../../componentes/Campo'
import { Hoja } from '../../componentes/Hoja'
import { Esqueleto, Vacio } from '../../componentes/Estados'
import { cn } from '../../componentes/utils'
import { VistaBD } from '../../vistas-bd/VistaBD'
import type { Cambio } from '../../vistas-bd/tipos'
import { AntesDespues, Bloque, Cargar, IndicadorGuardado, Marco, unir, useGuardado, useSinRed } from './comun'
import { ErrorApi } from '../../api/errores'
import { useAvisos } from '../../componentes/Avisos'
import { textos } from '../../textos/es-MX'
import { invalido } from './P2612Personas'
import { aMilesimas, cantidad, consumoSemanal, diasRestantes, estadoInventario, existenciaHoy, type InsumosInventario } from './calculos'
import { textosAjustesB as T } from './textos'

const t = T.i
type I = Esquemas['InventoryItem']; type P = Esquemas['Product']
const ORDEN = { out: 0, reorder: 1, sufficient: 2 } as Record<string, number>
const actualizar = (id: number, cuerpo: Esquemas['InventoryItemUpdate']) => pedir('patch', '/api/v1/inventory-items/{inventory_item_id}', { ruta: { inventory_item_id: id }, cuerpo })
const ICONO = { sufficient: CircleCheck, reorder: TriangleAlert, out: CircleAlert }
const COLOR = { sufficient: 'text-role-success', reorder: 'text-role-warning', out: 'text-role-error' }
function Estado({ e }: { e: 'sufficient' | 'reorder' | 'out' | null | undefined }) {
  if (!e) return <span className="text-text-secondary">{t.faltaRitmo}</span>
  const Icono = ICONO[e]
  return <span className="inline-flex items-center gap-1"><Icono aria-hidden className={cn('size-5', COLOR[e])} />{t.estados[e]}</span>
}

export default function P2614Inventario() {
  const productos = usarTodos()
  const datos = unir<[{ items: I[] }, P[]]>([usarConsulta('get', '/api/v1/inventory-items', { consulta: { limit: 200 } }), productos])
  const [nuevo, setNuevo] = useState(false)
  const sinRed = useSinRed()
  useEffect(() => { const f = (e: KeyboardEvent) => { if (e.altKey && e.key.toLowerCase() === 'n') { e.preventDefault(); setNuevo(true) } }; window.addEventListener('keydown', f); return () => window.removeEventListener('keydown', f) }, [])
  const boton = <Boton variante="secundario" deshabilitadoPor={sinRed} onClick={() => setNuevo(true)}>{t.agregar}</Boton>
  return (
    <Marco id="14">
      <p className="text-body text-text-secondary">{t.proposito}</p>
      <Cargar c={datos} que={t.que} esqueleto={<div className="flex flex-col gap-2">{[0, 1, 2, 3].map((i) => <Esqueleto key={i} className="h-10" />)}</div>}>{([inv, ps]) => (
        <Bloque titulo={t.titulo} accion={boton} ayuda={t.regla}>
          <Tabla items={inv.items} productos={ps} recargar={datos.recargar} vacio={<Vacio texto={t.vacio} accion={boton} />} />
          {nuevo && <Nuevo productos={ps.filter((p) => !inv.items.some((i) => i.product_id === p.id))} cerrar={() => setNuevo(false)} recargar={datos.recargar} />}
        </Bloque>
      )}</Cargar>
    </Marco>
  )
}
/** Catálogo completo de productos (lista corta, C1 §5: todas las páginas). */
function usarTodos(): EstadoConsulta<P[]> {
  const [e, setE] = useState<Omit<EstadoConsulta<P[]>, 'recargar'>>({ datos: undefined, cargando: true, error: null, desdeCopia: null, nuncaCargada: false })
  const [vuelta, setVuelta] = useState(0)
  useEffect(() => { let vivo = true; todas('get', '/api/v1/products', { consulta: { limit: 200 } }).then((d) => vivo && setE({ datos: d as P[], cargando: false, error: null, desdeCopia: null, nuncaCargada: false }), (error) => vivo && setE((x) => ({ ...x, cargando: false, error }))); return () => { vivo = false } }, [vuelta])
  return { ...e, recargar: () => setVuelta((v) => v + 1) }
}

function Tabla({ items, productos, recargar, vacio }: { items: I[]; productos: P[]; recargar: () => void; vacio: React.ReactNode }) {
  const prod = (i: I) => productos.find((p) => p.id === i.product_id)
  const nombre = (i: I) => (i.product_id ? prod(i)?.name ?? '' : t.sobrante(String(i.recipe_id ?? '')))
  const unidad = (i: I) => (i.product_id ? t.unidades[prod(i)?.base_unit ?? ''] ?? '' : t.porciones)
  const ins = (i: I): InsumosInventario => ({ tracking_mode: i.tracking_mode, quantity_milli: i.quantity_milli, counted_date: i.counted_date, rate_quantity_milli: i.rate_quantity_milli, rate_period_days: i.rate_period_days, reorder_level_milli: i.reorder_level_milli })
  const resumen = (x: InsumosInventario, u: string) => {
    const e = existenciaHoy(x)
    const l: [string, string][] = [[t.efectoEstado, e === null ? t.faltaRitmo : t.estados[estadoInventario(e, x.reorder_level_milli)!]]]
    if (x.tracking_mode === 'consumption_rate') {
      l.push([t.efectoEstimada, e === null ? t.faltaRitmo : `${cantidad(e)} ${u}`])
      const d = diasRestantes(x); const s = consumoSemanal(x)
      l.push([t.efectoDias, d === null ? t.faltaRitmo : String(d)], [t.efectoSemana, s === null ? t.faltaRitmo : `${cantidad(s)} ${u}`])
    }
    return l
  }
  const previa = (i: I, cambio: Partial<InsumosInventario>): Cambio[] => {
    const a = resumen(ins(i), unidad(i)); const b = resumen({ ...ins(i), ...cambio }, unidad(i))
    return b.map(([k, v]) => ({ etiqueta: k, antes: a.find((x) => x[0] === k)?.[1] ?? '—', despues: v }))
  }
  // `Ctrl+/` sobre {n} filas (U4 P-26.14, vista previa; H-08): estado de cada fila antes → después con la regla K7 y RN-15 (en el
  // navegador, como la vista previa de una sola fila), «{k} pasan a Reponer, {j} a Agotado y {s} a Suficiente» con el conteo de
  // estados de las seleccionadas antes → después, y la lista de las filas que cambian.
  const estadoDe = (x: InsumosInventario) => estadoInventario(existenciaHoy(x), x.reorder_level_milli)
  const textoEstado = (e: string | null) => (e ? t.estados[e] : t.faltaRitmo)
  const previaVarias = (sel: I[], cambio: Partial<InsumosInventario>): Cambio[] => {
    const l = sel.map((i) => ({ i, a: estadoDe(ins(i)), b: estadoDe({ ...ins(i), ...cambio }) }))
    const cambian = l.filter((x) => x.a !== x.b)
    const pasan = (e: string) => cambian.filter((x) => x.b === e).length
    const conteo = (k: 'a' | 'b') => {
      const c: [string, number][] = (['sufficient', 'reorder', 'out'] as const).map((e) => [t.estados[e], l.filter((x) => x[k] === e).length])
      const sinRitmo = l.filter((x) => x[k] === null).length
      return t.conteo(sinRitmo ? [...c, [t.faltaRitmo, sinRitmo]] : c)
    }
    return [
      { etiqueta: t.pasan(pasan('reorder'), pasan('out'), pasan('sufficient')), antes: conteo('a'), despues: conteo('b') },
      ...cambian.map((x) => ({ etiqueta: nombre(x.i), antes: textoEstado(x.a), despues: textoEstado(x.b) })),
    ]
  }
  const milli = (v: string, msg: string, positivo = false) => { if (v.trim() === '') return null; const m = aMilesimas(v); if (m === null || (positivo && m <= 0)) invalido(msg); return m }
  const cambioVarias = (p: string, v: string): Esquemas['InventoryItemUpdate'] & Partial<InsumosInventario> =>
    (p === 'tracking_mode' ? { tracking_mode: v as I['tracking_mode'] } : { reorder_level_milli: milli(v, t.nivelMal) })
  const filas = [...items].sort((a, b) => (ORDEN[a.status ?? ''] ?? 3) - (ORDEN[b.status ?? ''] ?? 3) || nombre(a).localeCompare(nombre(b), 'es-MX'))
  const [mov, setMov] = useState<{ i: I; modo: I['tracking_mode'] } | null>(null)
  const avisar = useAvisos()
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
  const MODOS = [t.modos.stock, t.modos.consumption_rate]
  const aplicarMov = async () => {
    if (!mov) return
    try { await actualizar(mov.i.id, { tracking_mode: mov.modo }); setMov(null); recargar() } catch (e) { avisar({ texto: e instanceof ErrorApi ? e.message : textos.errores.generico, rol: 'error' }) }
  }
  return (<>
    {mov && (
      <div key={mov.i.id} ref={panelMov} className="flex flex-col gap-2 rounded-md border border-border-divider p-3" data-prueba="previa-tablero">
        <p className="text-small font-medium">{t.mover(nombre(mov.i), t.modos[mov.modo])}</p>
        <AntesDespues cambios={previa(mov.i, { tracking_mode: mov.modo })} />
        <div className="flex flex-wrap gap-2"><Boton variante="primario" onClick={aplicarMov}>{T.comun.aplicar}</Boton><Boton onClick={() => setMov(null)}>{T.comun.cancelar}</Boton></div>
      </div>)}
    <VistaBD def={{
      clave: 'P-26.14-inventario', filas, idFila: (i) => i.id, vistas: ['tabla', 'tablero'], titulo: nombre, vacio,
      // U4 §0.2 punto 2 y P-26.14 (H-08): al soltar la tarjeta en el otro modo queda ahí con la vista previa del estado antes → después;
      // «Aplicar» (o Enter, con el foco en el botón) guarda y Esc la regresa a su columna sin guardar. Nada se guarda al soltar.
      columnasTablero: { propiedad: 'tracking_mode', columnas: MODOS, moverA: async (i, columna) => setMov({ i, modo: columna === t.modos.stock ? 'stock' : 'consumption_rate' }) },
      // `Ctrl+/` (U4 P-26.14, Atajos): modo o nivel de las filas seleccionadas con la vista previa sumada; `Enter` aplica a todas
      // (un `updateInventoryItem` por fila) y `Esc` cancela sin cambiar ninguna. Un nivel que no es cantidad no calcula vista previa.
      editarVarias: {
        propiedades: [
          { id: 'tracking_mode', nombre: t.modo, tipo: 'opcion', opciones: [{ valor: 'stock', texto: t.modos.stock }, { valor: 'consumption_rate', texto: t.modos.consumption_rate }] },
          { id: 'reorder', nombre: t.nivel, tipo: 'numero' },
        ],
        vistaPrevia: async (sel, p, v) => previaVarias(sel, cambioVarias(p, v)),
        guardar: async (sel, p, v) => {
          const c = cambioVarias(p, v)
          try { for (const i of sel) await actualizar(i.id, c) } finally { recargar() }
        },
      },
      propiedades: [
        { id: 'producto', nombre: t.producto, valor: nombre },
        { id: 'tracking_mode', nombre: t.modo, valor: (i) => t.modos[i.tracking_mode], grupo: (i) => t.modos[mov?.i.id === i.id ? mov.modo : i.tracking_mode],
          editar: { tipo: 'opcion', opciones: [{ valor: 'stock', texto: t.modos.stock }, { valor: 'consumption_rate', texto: t.modos.consumption_rate }], textoInicial: (i) => i.tracking_mode,
            vistaPrevia: async (i, v) => previa(i, { tracking_mode: v }), guardar: async (i, v) => { await actualizar(i.id, { tracking_mode: v as I['tracking_mode'] }); recargar() } } },
        { id: 'unidad', nombre: t.unidad, valor: unidad },
        { id: 'reorder', nombre: t.nivel, alinear: 'derecha', valor: (i) => i.reorder_level_milli ?? null, pintar: (i) => (i.reorder_level_milli !== undefined && i.reorder_level_milli !== null ? `${cantidad(i.reorder_level_milli)} ${unidad(i)}` : '—'),
          editar: { tipo: 'numero', textoInicial: (i) => (i.reorder_level_milli != null ? String(i.reorder_level_milli / 1000) : ''), vistaPrevia: async (i, v) => previa(i, { reorder_level_milli: v.trim() === '' ? null : aMilesimas(v) }),
            guardar: async (i, v) => { await actualizar(i.id, { reorder_level_milli: milli(v, t.nivelMal) }); recargar() } } },
        { id: 'rate_q', nombre: t.ritmoCant, alinear: 'derecha', valor: (i) => i.rate_quantity_milli ?? null, pintar: (i) => (i.rate_quantity_milli ? `${cantidad(i.rate_quantity_milli)} ${unidad(i)}` : '—'),
          editar: { tipo: 'numero', textoInicial: (i) => (i.rate_quantity_milli ? String(i.rate_quantity_milli / 1000) : ''), vistaPrevia: async (i, v) => previa(i, { rate_quantity_milli: aMilesimas(v) }),
            guardar: async (i, v) => { await actualizar(i.id, { rate_quantity_milli: milli(v, t.cantMal, true) }); recargar() } } },
        { id: 'rate_d', nombre: t.ritmoDias, alinear: 'derecha', valor: (i) => i.rate_period_days ?? null, pintar: (i) => (i.rate_period_days != null ? String(i.rate_period_days) : T.comun.sinValor),
          editar: { tipo: 'numero', textoInicial: (i) => (i.rate_period_days ? String(i.rate_period_days) : ''), vistaPrevia: async (i, v) => previa(i, { rate_period_days: /^\d+$/.test(v.trim()) ? Number(v) : null }),
            guardar: async (i, v) => { if (v.trim() && !(/^\d+$/.test(v.trim()) && Number(v) > 0)) invalido(t.diasMal); await actualizar(i.id, { rate_period_days: v.trim() ? Number(v) : null }); recargar() } } },
        { id: 'expires', nombre: t.caducidad, valor: (i) => i.expires_date ?? null, pintar: (i) => (i.expires_date ? fecha(i.expires_date) : '—'),
          editar: { tipo: 'fecha', textoInicial: (i) => i.expires_date ?? '', guardar: async (i, v) => { await actualizar(i.id, { expires_date: v || null }); recargar() } } },
        { id: 'existencia', nombre: t.existencia, alinear: 'derecha', valor: (i) => i.estimated_quantity_milli ?? i.quantity_milli,
          pintar: (i) => (i.tracking_mode === 'consumption_rate' && !(i.rate_quantity_milli && i.rate_period_days) ? t.faltaRitmo : `${cantidad(i.estimated_quantity_milli ?? i.quantity_milli)} ${unidad(i)}`) },
        { id: 'status', nombre: t.estado, valor: (i) => (i.status ? t.estados[i.status] : null), pintar: (i) => <Estado e={i.status ?? estadoInventario(existenciaHoy(ins(i)), i.reorder_level_milli)} /> },
        { id: 'days_left', nombre: t.diasRestantes, alinear: 'derecha', valor: (i) => i.days_left ?? null, pintar: (i) => (i.days_left !== undefined ? String(i.days_left) : '—') },
        { id: 'notes', nombre: t.notas, valor: (i) => i.notes ?? '', pintar: (i) => i.notes || T.comun.sinValor, editar: { tipo: 'texto', textoInicial: (i) => i.notes ?? '', guardar: async (i, v) => { await actualizar(i.id, { notes: v.trim() || null }); recargar() } } },
      ],
    }} />
  </>)
}
function Nuevo({ productos, cerrar, recargar }: { productos: P[]; cerrar: () => void; recargar: () => void }) {
  const sinRed = useSinRed()
  const g = useGuardado()
  const clave = useRef(nuevaClave())
  const [producto, setProducto] = useState('')
  const [modo, setModo] = useState<I['tracking_mode']>('stock')
  const [existencia, setExistencia] = useState('')
  const [error, setError] = useState('')
  const guardar = async () => {
    const q = existencia.trim() === '' ? 0 : aMilesimas(existencia)
    if (!producto) { setError(t.elegir); return }
    if (q === null) { setError(t.nivelMal); return }
    setError('')
    if (await g.ejecutar(async () => { await pedir('post', '/api/v1/inventory-items', { idempotencia: clave.current, cuerpo: { product_id: Number(producto), tracking_mode: modo, quantity_milli: q, counted_date: hoyISO() } }); recargar() })) cerrar()
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={t.agregar} pie={<><Boton variante="primario" tam="hoja" deshabilitadoPor={sinRed} onClick={() => void guardar()}>{T.comun.guardar}</Boton><IndicadorGuardado g={g} /></>}>
      <div className="flex flex-col gap-3" onKeyDown={(e) => { if (e.ctrlKey && e.key === 'Enter') void guardar() }}>
        <label className="flex flex-col gap-1 text-small font-medium">{t.producto}
          <select className={estilosControl} value={producto} onChange={(e) => setProducto(e.target.value)}><option value="">{t.elegir}</option>{productos.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}</select>
        </label>
        <label className="flex flex-col gap-1 text-small font-medium">{t.modo}
          <select className={estilosControl} value={modo} onChange={(e) => setModo(e.target.value as I['tracking_mode'])}><option value="stock">{t.modos.stock}</option><option value="consumption_rate">{t.modos.consumption_rate}</option></select>
        </label>
        <Campo etiqueta={t.existenciaInicial} inputMode="decimal" value={existencia} onChange={(e) => setExistencia(e.target.value)} />
        {error && <p role="alert" className="text-small text-role-error">{error}</p>}
      </div>
    </Hoja>
  )
}
