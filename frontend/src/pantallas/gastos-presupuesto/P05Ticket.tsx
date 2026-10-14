// P-05 Capturar ticket (U1 §2.1, §2.2, §3.7: N-070 a N-076; K3 «Guardar para después»). Pasos: Foto · Lectura · Revisar · Confirmar.
// Sin conexión: aviso «Requiere conexión» en lugar del contenido (U1 §2.4). Nada se guarda sin «Confirmar ticket» (decisión E).
// WF-3a bis: cada renglón se empareja con el catálogo (listProducts y getProduct, N-072) y, si se llega desde P-13
// «Terminar compra» (`?semana=2026-W42`), lo marcado como comprado se propone como renglones (getShoppingList, N-073).
import { useEffect, useRef, useState } from 'react'
import { Camera, ImageUp, Plus, X } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { conConsulta, navegar } from '../../app/enrutador'
import { armarUrl, nuevaClave, pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { usarConsulta } from '../../api/consulta'
import { aCentavos, aGramos, dinero, fechaHora } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { useAvisos } from '../../componentes/Avisos'
import { Boton } from '../../componentes/Boton'
import { Campo, estilosControl } from '../../componentes/Campo'
import { ConDatos, Esqueleto, Franja, Vacio } from '../../componentes/Estados'
import { cn } from '../../componentes/utils'
import { abrirRegistro } from '../../shell/estado'
import { textos } from '../../textos/es-MX'
import { Bloque, Selector, centavosATexto, mensajeDe, useCatalogos } from './comun'
import { t } from './textos'

type Lectura = Esquemas['TicketScan']
interface Renglon { clave: string; text: string; cantidad: string; precio: string; importe: string; cost_share_rule_id: string; product_id?: number; presentation_id?: number; confidence?: string; nombre?: string }
interface Borrador { store_id: string; expense_date: string; total: string; category_id: string; payment_method: 'cash' | 'debit' | 'transfer'; cost_share_rule_id: string; renglones: Renglon[] }
const milli = (m?: number) => (m == null ? '' : `${Math.floor(m / 1000)}${m % 1000 ? `.${String(m % 1000).padStart(3, '0').replace(/0+$/, '')}` : ''}`)
const txt = (c?: number | null) => (c == null ? '' : centavosATexto(c))
const id = (v: string) => (v ? Number(v) : undefined)

export default function P05Ticket({ consulta }: PropsPantalla) {
  const conexion = useConexion()
  if (!conexion) return <Franja rol="sinConexion">{textos.conexion.requiere}</Franja>
  const lectura = Number(consulta.get('lectura'))
  const semana = /^\d{4}-W\d{2}$/.test(consulta.get('semana') ?? '') ? consulta.get('semana')! : undefined
  return lectura > 0 ? <VerLectura key={lectura} id={lectura} semana={semana} /> : <Inicio />
}

function Pasos({ actual }: { actual: number }) {
  return (
    <ol aria-label={t.p05.paso(actual, t.p05.pasos[actual - 1])} className="grid grid-cols-4 gap-1">
      {t.p05.pasos.map((n, i) => (
        <li key={n} aria-current={i + 1 === actual ? 'step' : undefined} className={cn('flex flex-col gap-1 text-caption', i + 1 <= actual ? 'text-text-default' : 'text-text-secondary')}>
          <span className={cn('h-1 rounded-full', i + 1 <= actual ? 'bg-action-primary-bg' : 'bg-viz-track')} />{i + 1}. {n}
        </li>))}
    </ol>
  )
}

function Inicio() {
  const camara = useRef<HTMLInputElement>(null)
  const galeria = useRef<HTMLInputElement>(null)
  const [foto, setFoto] = useState<string | null>(null)
  const avisar = useAvisos()
  const cats = useCatalogos()
  const lista = usarConsulta('get', '/api/v1/ticket-scans', { consulta: { limit: 50 } })
  const subir = async (f: File | undefined) => {
    if (!f) return
    setFoto(URL.createObjectURL(f))
    const datos = new FormData(); datos.append('file', f)
    try { const r = await pedir('post', '/api/v1/ticket-scans', { cuerpo: datos, idempotencia: nuevaClave() }); navegar(conConsulta({ lectura: String(r.id) })) }
    catch (e) { setFoto(null); avisar({ texto: mensajeDe(e), rol: 'error' }) }
  }
  return (
    <div className="flex flex-col gap-6">
      <Pasos actual={foto ? 2 : 1} />
      {foto ? (
        <div className="flex flex-col gap-3" aria-busy>
          <img src={foto} alt={t.p05.tomarFoto} className="w-full rounded-lg border border-border-divider object-contain tablet:w-96" />
          <p role="status" className="text-body text-text-default">{t.p05.leyendo}</p>
          {Array.from({ length: 4 }, (_, i) => <Esqueleto key={i} />)}
        </div>
      ) : (
        <div className="flex flex-col gap-2 tablet:flex-row">
          <Boton variante="primario" tam="hoja" className="tablet:w-auto" onClick={() => camara.current?.click()}><Camera aria-hidden className="size-5" />{t.p05.tomarFoto}</Boton>
          <Boton tam="hoja" className="tablet:w-auto" onClick={() => galeria.current?.click()}><ImageUp aria-hidden className="size-5" />{t.p05.galeria}</Boton>
          <input ref={camara} type="file" hidden accept="image/*" capture="environment" aria-label={t.p05.tomarFoto} onChange={(e) => subir(e.target.files?.[0])} />
          <input ref={galeria} type="file" hidden accept="image/jpeg,image/png,image/webp" aria-label={t.p05.galeria} onChange={(e) => subir(e.target.files?.[0])} />
        </div>
      )}
      <Bloque titulo={t.p05.porConfirmar}>
        <ConDatos consulta={lista}>{(d) => d.items.length === 0 ? <Vacio texto={t.p05.ningunoPorConfirmar} /> : (
          <ul className="flex flex-col divide-y divide-border-divider rounded-lg border border-border-divider bg-surface-card">
            {(d.items as Lectura[]).map((s) => (
              <li key={s.id}>
                <button type="button" onClick={() => navegar(conConsulta({ lectura: String(s.id) }))} className="flex min-h-11 w-full items-center justify-between gap-3 px-3 py-2 text-left hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus">
                  <span className="flex flex-col"><span className="text-body text-text-default">{cats.tienda(s.display_store_id)}</span><span className="text-caption text-text-secondary">{fechaHora(s.created_at)}{s.has_draft && ` · ${t.p05.conBorrador}`}</span></span>
                  <span className="tabular-nums text-body font-semibold">{s.display_total_cents == null ? '—' : dinero(s.display_total_cents)}</span>
                </button>
              </li>))}
          </ul>)}
        </ConDatos>
      </Bloque>
    </div>
  )
}

function VerLectura({ id: lecturaId, semana }: { id: number; semana?: string }) {
  const r = usarConsulta('get', '/api/v1/ticket-scans/{attachment_id}', { ruta: { attachment_id: lecturaId } })
  useEffect(() => { if (r.datos?.ocr_status === 'pending') { const h = setTimeout(r.recargar, 2000); return () => clearTimeout(h) } }, [r.datos, r.recargar])
  const foto = <img src={armarUrl('/api/v1/attachments/{attachment_id}/file', { attachment_id: lecturaId })} alt={t.p05.tomarFoto} className="w-full rounded-lg border border-border-divider object-contain tablet:w-96" />
  return (
    <ConDatos consulta={r}>{(s) => {
      if (s.ocr_status === 'pending') return <div className="flex flex-col gap-4"><Pasos actual={2} />{foto}<p role="status">{t.p05.leyendo}</p>{Array.from({ length: 4 }, (_, i) => <Esqueleto key={i} />)}</div>
      if (s.ocr_status === 'failed') return (
        <div className="flex flex-col gap-4"><Pasos actual={2} />
          <Franja rol="error" alerta>{t.p05.noSeLeyo}</Franja>{foto}
          <div className="flex flex-wrap gap-2"><Boton variante="primario" onClick={() => navegar(conConsulta({ lectura: null }))}>{t.p05.repetir}</Boton><Boton onClick={() => abrirRegistro('gasto')}>{t.p05.aMano}</Boton></div>
        </div>)
      return <Editor s={s} foto={foto} semana={semana} />
    }}</ConDatos>
  )
}

function inicial(s: Lectura): Borrador {
  const d = s.has_draft ? s.draft : null
  const lineas = (d?.lines ?? s.lines ?? []) as { line_no: number; text?: string; product_id?: number; presentation_id?: number; quantity_milli?: number; unit_price_cents?: number; line_total_cents?: number; cost_share_rule_id?: number; confidence?: string }[]
  const forma = d?.payment_method && d.payment_method !== 'credit' ? d.payment_method : 'cash'
  return {
    store_id: String(d?.store_id ?? s.proposed_store_id ?? ''), expense_date: d?.expense_date ?? s.proposed_expense_date ?? '', total: txt(d?.total_cents ?? s.proposed_total_cents),
    category_id: String(d?.category_id ?? ''), payment_method: forma, cost_share_rule_id: String(d?.cost_share_rule_id ?? ''),
    renglones: lineas.map((l) => ({ clave: nuevaClave(), text: l.text ?? '', cantidad: milli(l.quantity_milli), precio: txt(l.unit_price_cents), importe: txt(l.line_total_cents), cost_share_rule_id: String(l.cost_share_rule_id ?? ''), product_id: l.product_id, presentation_id: l.presentation_id, confidence: (s.lines ?? []).find((p) => p.line_no === l.line_no)?.confidence })),
  }
}

function Editor({ s, foto, semana }: { s: Lectura; foto: React.ReactNode; semana?: string }) {
  const [b, setB] = useState<Borrador>(() => inicial(s))
  const [enCurso, setEnCurso] = useState<'borrador' | 'confirmar' | null>(null)
  const [error, setError] = useState('')
  const cats = useCatalogos()
  const avisar = useAvisos()
  const poner = (k: keyof Borrador) => (v: string) => setB({ ...b, [k]: v })
  const cambiarR = (i: number, c: Partial<Renglon>) => setB({ ...b, renglones: b.renglones.map((r, j) => (j === i ? { ...r, ...c } : r)) })
  const suma = b.renglones.reduce((a, r) => a + (aCentavos(r.importe) ?? 0), 0)
  const total = aCentavos(b.total) ?? 0
  const lineas = () => b.renglones.map((r, i) => ({ line_no: i + 1, text: r.text, description: r.text || t.p05.renglon(i + 1), product_id: r.product_id, presentation_id: r.presentation_id, quantity_milli: aGramos(r.cantidad) ?? 1000, unit_price_cents: aCentavos(r.precio) ?? undefined, line_total_cents: aCentavos(r.importe) ?? 0, cost_share_rule_id: id(r.cost_share_rule_id) }))
  const guardarDespues = async () => {
    setEnCurso('borrador')
    try {
      await pedir('put', '/api/v1/ticket-scans/{attachment_id}/draft', { ruta: { attachment_id: s.id }, cuerpo: { store_id: id(b.store_id), expense_date: b.expense_date || undefined, total_cents: aCentavos(b.total) ?? undefined, category_id: id(b.category_id), payment_method: b.payment_method, cost_share_rule_id: id(b.cost_share_rule_id), lines: lineas().map(({ description: _d, ...l }) => l) } })
      avisar({ texto: t.p05.guardadoDespues, rol: 'exito' }); navegar(conConsulta({ lectura: null }))
    } catch (e) { avisar({ texto: mensajeDe(e), rol: 'error' }) } finally { setEnCurso(null) }
  }
  const confirmar = async () => {
    if (!b.category_id || total <= 0 || !b.expense_date) return setError(t.p05.falta)
    setError(''); setEnCurso('confirmar')
    try {
      const items = lineas().map(({ text: _t, product_id: _p, ...l }) => l)
      const r = await pedir('post', '/api/v1/ticket-scans/{attachment_id}/confirm', { ruta: { attachment_id: s.id }, idempotencia: nuevaClave(), cuerpo: { expense: { category_id: Number(b.category_id), amount_cents: total, expense_date: b.expense_date, payment_method: b.payment_method, store_id: id(b.store_id), cost_share_rule_id: id(b.cost_share_rule_id), items } } })
      avisar({ texto: textos.registro.tipos.ticket + ' · ' + t.comun.guardado, rol: 'exito' }); navegar(rutaDe('P-04', { id: r.expense_id }))
    } catch (e) { setError(mensajeDe(e)) } finally { setEnCurso(null) }
  }
  const reglas = cats.reglas.map((x) => ({ valor: String(x.id), texto: x.name }))
  return (
    <div className="flex flex-col gap-6">
      <Pasos actual={enCurso === 'confirmar' ? 4 : 3} />
      <div className="grid gap-6 desktop:grid-cols-2">
        <div className="flex flex-col gap-2">{foto}<Boton variante="enlace" onClick={() => navegar(conConsulta({ lectura: null }))}>{t.p05.repetir}</Boton></div>
        <Bloque titulo={t.p05.encabezado}>
          <div className="grid grid-cols-1 gap-3 tablet:grid-cols-2">
            <Selector etiqueta={t.comun.tienda} valor={b.store_id} vacio={t.comun.ninguna} alCambiar={poner('store_id')} opciones={cats.tiendas.map((x) => ({ valor: String(x.id), texto: x.name }))} />
            <Campo etiqueta={t.comun.fecha} type="date" value={b.expense_date} onChange={(e) => poner('expense_date')(e.target.value)} />
            <Campo etiqueta={t.p05.total} inputMode="decimal" value={b.total} onChange={(e) => poner('total')(e.target.value)} />
            <Selector etiqueta={t.comun.categoria} valor={b.category_id} vacio={t.comun.ninguna} alCambiar={poner('category_id')} opciones={cats.categorias.map((x) => ({ valor: String(x.id), texto: x.name }))} />
            <Selector etiqueta={t.comun.formaPago} valor={b.payment_method} alCambiar={poner('payment_method')} opciones={(['cash', 'debit', 'transfer'] as const).map((f) => ({ valor: f, texto: textos.registro.formas[f] }))} />
            <Selector etiqueta={t.p05.cuenta} valor={b.cost_share_rule_id} vacio={t.comun.ninguna} alCambiar={poner('cost_share_rule_id')} opciones={reglas} />
          </div>
          <p className="text-caption text-text-secondary">{t.p05.creditoPendiente}</p>
        </Bloque>
      </div>
      {semana && <DeLaLista semana={semana} alAgregar={(rs, tienda) => setB((x) => ({ ...x, store_id: x.store_id || (tienda ? String(tienda) : ''), renglones: [...x.renglones, ...rs] }))} />}
      <Bloque titulo={t.p05.renglones} accion={<Boton onClick={() => setB({ ...b, renglones: [...b.renglones, { clave: nuevaClave(), text: '', cantidad: '1', precio: '', importe: '', cost_share_rule_id: '' }] })}><Plus aria-hidden className="size-5" />{t.p05.agregarRenglon}</Boton>}>
        <ol className="flex flex-col gap-3">
          {b.renglones.map((r, i) => (
            <li key={r.clave} className={cn('grid grid-cols-2 gap-2 rounded-lg border bg-surface-card p-3 desktop:grid-cols-7 desktop:items-end', r.confidence === 'low' ? 'border-role-warning' : 'border-border-divider')}>
              <Campo className="col-span-2" etiqueta={`${t.p05.renglon(i + 1)} · ${t.p05.texto}`} value={r.text} onChange={(e) => cambiarR(i, { text: e.target.value })} />
              <Campo etiqueta={t.p05.cantidad} inputMode="decimal" value={r.cantidad} onChange={(e) => cambiarR(i, { cantidad: e.target.value })} />
              <Campo etiqueta={t.p05.precio} inputMode="decimal" value={r.precio} onChange={(e) => cambiarR(i, { precio: e.target.value })} />
              <Campo etiqueta={t.p05.importe} inputMode="decimal" value={r.importe} onChange={(e) => cambiarR(i, { importe: e.target.value })} />
              <Selector etiqueta={t.p05.cuenta} valor={r.cost_share_rule_id} vacio={t.comun.ninguna} alCambiar={(v) => cambiarR(i, { cost_share_rule_id: v })} opciones={reglas} />
              <Catalogo r={r} n={i + 1} alElegir={(c) => { cambiarR(i, c); avisar({ texto: t.p05.emparejado(i + 1), rol: 'exito' }) }} />
              <div className="flex items-center justify-between gap-2">
                {r.confidence && <span className={cn('rounded-full px-2 text-caption', r.confidence === 'low' ? 'bg-role-warning-surface text-role-warning' : 'bg-surface-hover text-text-secondary')}>{t.p05.confianza}: {t.p05.niveles[r.confidence]}</span>}
                <Boton tam="icono" variante="fantasma" aria-label={`${t.p05.quitarRenglon} ${i + 1}`} onClick={() => setB({ ...b, renglones: b.renglones.filter((_, j) => j !== i) })}><X aria-hidden className="size-6" size={ICONO_MD} /></Boton>
              </div>
            </li>))}
        </ol>
        <p className="text-small text-text-secondary">{suma === total ? t.p05.cuadra : t.p05.noCuadra(dinero(Math.abs(total - suma)))}</p>
      </Bloque>
      {error && <Franja rol="error" alerta>{error}</Franja>}
      <div className="flex flex-col-reverse gap-2 tablet:flex-row tablet:justify-end">
        <Boton onClick={() => navegar(conConsulta({ lectura: null }))}>{t.p05.cancelar}</Boton>
        <Boton aria-busy={enCurso === 'borrador'} onClick={guardarDespues}>{t.p05.guardarDespues}</Boton>
        <Boton variante="primario" aria-busy={enCurso === 'confirmar'} onClick={confirmar}>{t.p05.confirmar}</Boton>
      </div>
    </div>
  )
}

type Producto = Esquemas['Product']
/** N-072: «Producto del catálogo» del renglón: busca por texto (listProducts `q`), elige el producto y su presentación
 *  (getProduct). La confirmación manda `presentation_id` (ExpenseItemInput); el borrador guarda también `product_id`. */
function Catalogo({ r, n, alElegir }: { r: Renglon; n: number; alElegir: (c: Partial<Renglon>) => void }) {
  const [abierto, setAbierto] = useState(false)
  const [texto, setTexto] = useState('')
  const [q, setQ] = useState('')
  const [elegido, setElegido] = useState<Producto | null>(null)
  useEffect(() => { const h = setTimeout(() => setQ(texto.trim()), 300); return () => clearTimeout(h) }, [texto])
  const res = usarConsulta('get', '/api/v1/products', { consulta: { q, limit: 10 } }, abierto && q.length >= 2)
  const det = usarConsulta('get', '/api/v1/products/{product_id}', { ruta: { product_id: elegido?.id ?? 0 } }, !!elegido)
  const actual = usarConsulta('get', '/api/v1/products/{product_id}', { ruta: { product_id: r.product_id ?? 0 } }, !!r.product_id && !r.nombre)
  const nombreActual = r.nombre ?? (actual.datos ? [actual.datos.name, actual.datos.presentations.find((p) => p.id === r.presentation_id)?.name].filter(Boolean).join(' · ') : r.product_id ? '…' : t.p05.sinProducto)
  const elegir = (c: Partial<Renglon>) => { alElegir(c); setAbierto(false); setElegido(null); setTexto('') }
  const opcion = 'flex min-h-11 w-full items-center px-3 text-left text-body hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-border-focus'
  return (
    <div className="col-span-2 flex flex-col gap-2 desktop:col-span-7">
      <div className="flex flex-wrap items-center gap-2 text-small">
        <span className="text-text-secondary">{t.p05.catalogo}:</span><span className="text-text-default">{nombreActual}</span>
        <Boton variante="enlace" aria-expanded={abierto} onClick={() => { setAbierto(!abierto); setElegido(null) }}>{abierto ? t.p05.cerrarBusqueda : t.p05.buscarCatalogo}</Boton>
      </div>
      {abierto && !elegido && (
        <div className="flex flex-col gap-2">
          <input type="search" autoFocus value={texto} onChange={(e) => setTexto(e.target.value)} aria-label={`${t.p05.buscarCatalogo} · ${t.p05.renglon(n)}`} placeholder={t.p05.buscarCatalogo} className={cn(estilosControl, 'w-full')} />
          {q.length < 2 ? <p className="text-caption text-text-secondary">{t.p05.escribeDos}</p> : (
            <ConDatos consulta={res}>{(d) => (d.items as Producto[]).length === 0 ? <p className="text-small text-text-secondary">{t.p05.sinResultados}</p> : (
              <ul className="flex flex-col divide-y divide-border-divider rounded-md border border-border-divider bg-surface-card">
                {(d.items as Producto[]).map((p) => <li key={p.id}><button type="button" className={opcion} onClick={() => setElegido(p)}>{p.name}</button></li>)}
              </ul>)}
            </ConDatos>)}
        </div>
      )}
      {abierto && elegido && (
        <ConDatos consulta={det}>{(p) => (
          <div className="flex flex-col gap-2" role="group" aria-label={t.p05.elegirPresentacion(p.name)}>
            <p className="text-small text-text-secondary">{t.p05.elegirPresentacion(p.name)}</p>
            <ul className="flex flex-col divide-y divide-border-divider rounded-md border border-border-divider bg-surface-card">
              {p.presentations.filter((x) => x.is_active).map((x) => (
                <li key={x.id}><button type="button" className={opcion} onClick={() => elegir({ product_id: p.id, presentation_id: x.id, nombre: `${p.name} · ${x.name}` })}>{x.name}{x.brand && ` · ${x.brand}`}</button></li>))}
              <li><button type="button" className={opcion} onClick={() => elegir({ product_id: p.id, presentation_id: undefined, nombre: p.name })}>{t.p05.soloProducto}</button></li>
            </ul>
          </div>)}
        </ConDatos>
      )}
    </div>
  )
}

/** N-073: artículos marcados como comprados en la lista de la semana que mostraba P-13 (semana explícita, R3-01), como
 *  renglones propuestos que la persona agrega con un toque. PROVISIONAL (bloqueo WF3a-bis-gastos-presupuesto): no se
 *  combinan solos con los renglones leídos; cantidad = piezas de la presentación (`packages_to_buy`) y tienda del encabezado
 *  solo si todos los artículos son de la misma tienda. */
function DeLaLista({ semana, alAgregar }: { semana: string; alAgregar: (r: Renglon[], tienda?: number) => void }) {
  const l = usarConsulta('get', '/api/v1/shopping-list', { consulta: { week: semana, is_bought: true } })
  const [estado, setEstado] = useState<'listo' | 'agregando' | 'agregado'>('listo')
  const avisar = useAvisos()
  if (estado === 'agregado') return <Franja rol="info">{t.p05.agregadosDeLista}</Franja>
  return (
    <ConDatos consulta={l}>{(d) => {
      const items = d.groups.flatMap((g) => g.items.filter((it) => it.is_bought).map((it) => ({ it, tienda: it.store_id ?? it.assigned_store_id ?? g.store_id })))
      if (!items.length) return <Franja rol="info">{t.p05.listaVacia}</Franja>
      const agregar = async () => {
        setEstado('agregando')
        try {
          const nombres = await Promise.all(items.map(({ it }) => (it.product_id ? pedir('get', '/api/v1/products/{product_id}', { ruta: { product_id: it.product_id } }) : Promise.resolve(null))))
          const rs = items.map(({ it }, k): Renglon => {
            const p = nombres[k], pres = it.presentation_id ?? it.recommended_presentation_id
            const np = p?.presentations.find((x) => x.id === pres)?.name
            return { clave: nuevaClave(), text: p?.name ?? it.free_text ?? '', cantidad: it.packages_to_buy ? String(it.packages_to_buy) : '', precio: '', importe: '', cost_share_rule_id: '',
              product_id: it.product_id, presentation_id: pres, nombre: p ? [p.name, np].filter(Boolean).join(' · ') : undefined }
          })
          const tiendas = [...new Set(items.map((x) => x.tienda).filter((x) => x != null))]
          alAgregar(rs, tiendas.length === 1 ? tiendas[0] : undefined); setEstado('agregado')
        } catch (e) { setEstado('listo'); avisar({ texto: mensajeDe(e), rol: 'error' }) }
      }
      return <Franja rol="info" accion={<Boton aria-busy={estado === 'agregando'} onClick={agregar}><Plus aria-hidden className="size-5" />{t.p05.agregarDeLista(items.length)}</Boton>}>{t.p05.deLaLista(items.length, semana.slice(6))}</Franja>
    }}</ConDatos>
  )
}
