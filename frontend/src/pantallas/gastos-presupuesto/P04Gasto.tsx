// P-04 Gasto (U1 §2.1, §2.2, §3.6: N-060 a N-066). Propiedades editables, reparto, devengo, efectivo, renglones y adjunto.
// WF-3a bis: foto del ticket (downloadAttachment) y «Quitar adjunto» (deleteAttachment, N-064). Expense solo trae
// has_attachment: el id del adjunto sale de queryView en galería (image_attachment_id) o de la respuesta al adjuntar
// (PROVISIONAL, bloqueo WF3a-bis-gastos-presupuesto: un PDF adjuntado antes no se puede ubicar).
import { useRef, useState } from 'react'
import { Copy, CreditCard, Paperclip, Trash2, CircleCheck, Unlink } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { navegar } from '../../app/enrutador'
import { armarUrl, nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import { aCentavos, dinero, fecha, hoyISO, mes } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { useAvisos } from '../../componentes/Avisos'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Esqueleto, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { textos } from '../../textos/es-MX'
import { CeldaEditable } from '../../vistas-bd/CeldaEditable'
import { VistaBD } from '../../vistas-bd/VistaBD'
import type { Propiedad } from '../../vistas-bd/tipos'
import { Bloque, Confirmar, Fila, centavosATexto, mensajeDe, useCatalogos, useSinRed } from './comun'
import { t } from './textos'

type Gasto = Esquemas['Expense']
const milesimas = (m: number) => `${Math.floor(m / 1000)}${m % 1000 ? `.${String(m % 1000).padStart(3, '0').replace(/0+$/, '')}` : ''}`

export default function P04Gasto({ params }: PropsPantalla) {
  const id = Number(params.id)
  const g = usarConsulta('get', '/api/v1/expenses/{expense_id}', { ruta: { expense_id: id } })
  if (g.error instanceof ErrorApi && g.error.estado === 404)
    return <Vacio texto={t.p04.noExiste} accion={<Boton onClick={() => navegar(rutaDe('P-02'))}>{textos.app.volver}</Boton>} />
  return (
    <ConDatos consulta={g} esqueleto={<div className="flex flex-col gap-3">{Array.from({ length: 8 }, (_, i) => <Esqueleto key={i} />)}</div>}>
      {(gasto) => <Detalle gasto={gasto} recargar={g.recargar} />}
    </ConDatos>
  )
}

function Detalle({ gasto, recargar }: { gasto: Gasto; recargar: () => void }) {
  const cats = useCatalogos()
  const avisar = useAvisos()
  const sinRed = useSinRed()
  const archivo = useRef<HTMLInputElement>(null)
  const [hoja, setHoja] = useState<'pagado' | 'eliminar' | 'quitar' | null>(null)
  const [subido, setSubido] = useState<{ id: number; mime: string } | null>(null)
  const galeria = usarConsulta('get', '/api/v1/views/{entity}', { ruta: { entity: 'expense' }, consulta: { mode: 'gallery', period: gasto.expense_date.slice(0, 7), limit: 200 } }, gasto.has_attachment)
  const fotoId = (galeria.datos?.items as Esquemas['ViewItem'][] | undefined)?.find((i) => i.id === gasto.id)?.image_attachment_id
  const adjunto = fotoId ? { id: fotoId, mime: 'image' } : subido
  const urlAdjunto = adjunto ? armarUrl('/api/v1/attachments/{attachment_id}/file', { attachment_id: adjunto.id }) : ''
  const [pago, setPago] = useState({ monto: centavosATexto(gasto.amount_cents), fecha: hoyISO(), error: '' })
  const ruta = { expense_id: gasto.id }
  const parchar = async (_: Gasto, cuerpo: Esquemas['ExpenseUpdate']) => { await pedir('patch', '/api/v1/expenses/{expense_id}', { ruta, cuerpo }); recargar() }
  const opc = <T extends { id: number }>(l: T[], n: (x: T) => string) => l.map((x) => ({ valor: String(x.id), texto: n(x) }))
  const props: Propiedad<Gasto>[] = [
    { id: 'expense_date', nombre: t.comun.fecha, valor: (x) => x.expense_date, pintar: (x) => fecha(x.expense_date), editar: { tipo: 'fecha', textoInicial: (x) => x.expense_date, guardar: (x, v) => parchar(x, { expense_date: v }) } },
    { id: 'amount_cents', nombre: t.comun.monto, valor: (x) => x.amount_cents, pintar: (x) => dinero(x.amount_cents), editar: { tipo: 'dinero', textoInicial: (x) => centavosATexto(x.amount_cents),
      vistaPrevia: async (x, v) => { const c = aCentavos(v); if (c == null || c <= 0) throw new Error(textos.registro.montoInvalido); return [{ etiqueta: t.comun.monto, antes: dinero(x.amount_cents), despues: dinero(c) }] },
      guardar: (x, v) => parchar(x, { amount_cents: aCentavos(v)! }) } },
    { id: 'category_id', nombre: t.comun.categoria, valor: (x) => cats.categoria(x.category_id), editar: { tipo: 'opcion', opciones: opc(cats.categorias, (c) => c.name), textoInicial: (x) => String(x.category_id), guardar: (x, v) => parchar(x, { category_id: Number(v) }) } },
    { id: 'person_id', nombre: t.comun.persona, valor: (x) => cats.persona(x.person_id), editar: { tipo: 'opcion', opciones: opc(cats.personas, (p) => p.display_name), textoInicial: (x) => String(x.person_id ?? ''), guardar: (x, v) => parchar(x, { person_id: v ? Number(v) : null }) } },
    { id: 'description', nombre: t.comun.descripcion, valor: (x) => x.description ?? '—', editar: { tipo: 'texto', textoInicial: (x) => x.description ?? '', guardar: (x, v) => parchar(x, { description: v || null }) } },
    { id: 'store_id', nombre: t.comun.tienda, valor: (x) => cats.tienda(x.store_id), editar: { tipo: 'opcion', opciones: opc(cats.tiendas, (s) => s.name), textoInicial: (x) => String(x.store_id ?? ''), guardar: (x, v) => parchar(x, { store_id: v ? Number(v) : null }) } },
    { id: 'payment_method', nombre: t.comun.formaPago, valor: (x) => textos.registro.formas[x.payment_method] },
    { id: 'kind', nombre: t.comun.tipo, valor: (x) => t.tipos[x.kind], editar: { tipo: 'opcion', opciones: Object.entries(t.tipos).map(([valor, texto]) => ({ valor, texto })), textoInicial: (x) => x.kind, guardar: (x, v) => parchar(x, { kind: v as Gasto['kind'] }) } },
    { id: 'cost_share_rule_id', nombre: t.comun.reparto, valor: (x) => cats.regla(x.cost_share_rule_id) },
    { id: 'accrual_months', nombre: t.p04.mesesDevengo, valor: (x) => x.accrual_months, editar: { tipo: 'numero', textoInicial: (x) => String(x.accrual_months),
      vistaPrevia: async (x, v) => { const n = Number(v); if (!Number.isInteger(n) || n < 1) throw new Error(textos.errores.revisaCampos); return [{ etiqueta: t.p04.mesesDevengo, antes: String(x.accrual_months), despues: String(n) }] },
      guardar: (x, v) => parchar(x, { accrual_months: Number(v) }) } },
    { id: 'notes', nombre: t.comun.nota, valor: (x) => x.notes ?? '—', editar: { tipo: 'texto', textoInicial: (x) => x.notes ?? '', guardar: (x, v) => parchar(x, { notes: v || null }) } },
  ]
  const accion = async (f: () => Promise<unknown>, ok: string, despues?: (r: unknown) => void) => {
    try { const r = await f(); avisar({ texto: ok, rol: 'exito' }); despues ? despues(r) : recargar() } catch (e) { avisar({ texto: mensajeDe(e), rol: 'error' }) }
  }
  const adjuntar = (f: File | undefined) => {
    if (!f) return
    const datos = new FormData(); datos.append('file', f); datos.append('kind', 'ticket_photo')
    accion(() => pedir('post', '/api/v1/expenses/{expense_id}/attachments', { ruta, cuerpo: datos, idempotencia: nuevaClave() }), t.p04.adjuntado,
      (r) => { const a = r as Esquemas['Attachment']; setSubido({ id: a.id, mime: a.mime_type }); recargar(); galeria.recargar() })
  }
  const irCompra = async () => {
    try { const c = await pedir('get', '/api/v1/credit-purchases/{credit_purchase_id}', { ruta: { credit_purchase_id: gasto.credit_purchase_id! } }); navegar(rutaDe('P-19', { id: c.credit_id, compra: c.id })) }
    catch (e) { avisar({ texto: mensajeDe(e), rol: 'error' }) }
  }
  const marcar = () => {
    const c = aCentavos(pago.monto)
    if (c == null || c <= 0) return setPago({ ...pago, error: textos.registro.montoInvalido })
    accion(() => pedir('post', '/api/v1/expenses/{expense_id}/mark-paid', { ruta, cuerpo: { amount_cents: c, expense_date: pago.fecha }, idempotencia: nuevaClave() }), t.p04.marcado, () => { setHoja(null); recargar() })
  }
  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-wrap gap-2">
        {!gasto.has_attachment && <Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => archivo.current?.click()}><Paperclip aria-hidden className="size-5" />{t.p04.adjuntar}</Boton>}
        <input ref={archivo} type="file" hidden accept="image/jpeg,image/png,image/webp,application/pdf" onChange={(e) => adjuntar(e.target.files?.[0])} />
        <Boton deshabilitadoPor={sinRed} onClick={() => accion(() => pedir('post', '/api/v1/expenses/{expense_id}/duplicate', { ruta, idempotencia: nuevaClave() }), t.p04.duplicado, (r) => navegar(rutaDe('P-04', { id: (r as Gasto).id })))}><Copy aria-hidden className="size-5" />{t.p04.duplicar}</Boton>
        {gasto.kind === 'upcoming' && <Boton deshabilitadoPor={sinRed} onClick={() => setHoja('pagado')}><CircleCheck aria-hidden className="size-5" />{t.p04.marcarPagado}</Boton>}
        {gasto.credit_purchase_id != null && <Boton onClick={irCompra}><CreditCard aria-hidden className="size-5" />{t.p04.compraCredito}</Boton>}
        <Boton variante="destructivo" deshabilitadoPor={sinRed} onClick={() => setHoja('eliminar')}><Trash2 aria-hidden className="size-5" />{t.comun.eliminar}</Boton>
      </div>
      {gasto.has_attachment && (
        <Bloque titulo={t.p04.ticketTitulo} accion={<Boton deshabilitadoPor={sinRed ?? (adjunto ? undefined : t.p04.sinUbicar)} onClick={() => setHoja('quitar')}><Unlink aria-hidden className="size-5" />{t.p04.quitarAdjunto}</Boton>}>
          {adjunto ? (
            adjunto.mime.startsWith('image')
              ? <a href={urlAdjunto} target="_blank" rel="noopener" className="self-start rounded-lg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus" aria-label={t.p04.abrirArchivo}>
                  <img src={urlAdjunto} alt={t.p04.fotoTicket} className="max-h-96 w-full rounded-lg border border-border-divider object-contain tablet:w-96" />
                </a>
              : <a href={urlAdjunto} target="_blank" rel="noopener" className="self-start text-body text-text-link underline">{t.p04.abrirArchivo}</a>
          ) : galeria.cargando ? <Esqueleto /> : <p className="text-small text-text-secondary">{t.p04.sinUbicar}</p>}
        </Bloque>
      )}
      <Bloque titulo={t.p04.propiedades}>
        <div>{props.map((p) => <Fila key={p.id} etiqueta={p.nombre}>{p.editar && !sinRed ? <CeldaEditable fila={gasto} prop={p}>{p.pintar ? p.pintar(gasto) : String(p.valor(gasto))}</CeldaEditable> : (p.pintar ? p.pintar(gasto) : String(p.valor(gasto)))}</Fila>)}</div>
      </Bloque>
      <Bloque titulo={t.p04.repartoTitulo}>
        <div>
          {(gasto.shares ?? []).map((s) => <Fila key={s.person_id} etiqueta={cats.persona(s.person_id)}>{dinero(s.amount_cents)}</Fila>)}
          <Fila etiqueta={t.p04.parteHogar}><span className="font-semibold">{dinero(gasto.household_share_cents ?? 0)}</span></Fila>
          <Fila etiqueta={t.p04.parteKaren}><span className="font-semibold">{dinero(gasto.external_share_cents ?? 0)}</span></Fila>
        </div>
      </Bloque>
      <div className="grid gap-6 desktop:grid-cols-2">
        <Bloque titulo={t.p04.devengo}><div>{(gasto.accruals ?? []).map((a) => <Fila key={a.month} etiqueta={mes(a.month)}>{dinero(a.amount_cents)}</Fila>)}</div></Bloque>
        <Bloque titulo={t.p04.efectivo}><div>{(gasto.cash_flows ?? []).map((c) => <Fila key={c.date} etiqueta={fecha(c.date)}>{dinero(c.amount_cents)}</Fila>)}</div></Bloque>
      </div>
      <Renglones id={gasto.id} />
      <Hoja abierta={hoja === 'pagado'} alCambiar={(a) => setHoja(a ? 'pagado' : null)} titulo={t.p04.marcarPagado}
        pie={<div className="flex justify-end gap-2"><Boton onClick={() => setHoja(null)}>{textos.app.cancelar}</Boton><Boton variante="primario" onClick={marcar}>{textos.app.guardar}</Boton></div>}>
        <div className="flex flex-col gap-3">
          <Campo etiqueta={t.p04.montoReal} inputMode="decimal" value={pago.monto} error={pago.error || undefined} onChange={(e) => setPago({ ...pago, monto: e.target.value, error: '' })} />
          <Campo etiqueta={t.p04.fechaReal} type="date" value={pago.fecha} onChange={(e) => setPago({ ...pago, fecha: e.target.value })} />
        </div>
      </Hoja>
      <Confirmar abierta={hoja === 'quitar'} alCambiar={(a) => setHoja(a ? 'quitar' : null)} titulo={t.p04.quitarAdjunto} texto={t.p04.confirmarQuitar} accion={t.p04.quitarAdjunto}
        alConfirmar={() => accion(() => pedir('delete', '/api/v1/attachments/{attachment_id}', { ruta: { attachment_id: adjunto!.id } }), t.p04.quitado, () => { setHoja(null); setSubido(null); recargar() })} />
      <Confirmar abierta={hoja === 'eliminar'} alCambiar={(a) => setHoja(a ? 'eliminar' : null)} titulo={t.comun.eliminar} texto={t.p04.confirmarEliminar} accion={t.comun.eliminar}
        alConfirmar={() => accion(() => pedir('delete', '/api/v1/expenses/{expense_id}', { ruta }), t.p04.eliminado, () => navegar(rutaDe('P-02')))} />
    </div>
  )
}

function Renglones({ id }: { id: number }) {
  const r = usarConsulta('get', '/api/v1/expenses/{expense_id}/items', { ruta: { expense_id: id }, consulta: { limit: 200 } })
  type Renglon = Esquemas['ExpenseItem']
  return (
    <Bloque titulo={t.p04.renglones}>
      <ConDatos consulta={r}>{(d) => (
        <VistaBD def={{
          clave: 'P-04-renglones', filas: d.items, idFila: (x: Renglon) => x.id, vistas: ['tabla'], titulo: (x: Renglon) => x.description, vacio: <Vacio texto={t.p04.sinRenglones} />,
          propiedades: [
            { id: 'line_no', nombre: '#', valor: (x: Renglon) => x.line_no },
            { id: 'description', nombre: t.p04.producto, valor: (x: Renglon) => x.description },
            { id: 'quantity_milli', nombre: t.p04.cantidad, valor: (x: Renglon) => x.quantity_milli, pintar: (x: Renglon) => milesimas(x.quantity_milli), alinear: 'derecha' },
            { id: 'unit_price_cents', nombre: t.p04.precioUnitario, valor: (x: Renglon) => x.unit_price_cents ?? null, pintar: (x: Renglon) => (x.unit_price_cents == null ? '—' : dinero(x.unit_price_cents)), alinear: 'derecha' },
            { id: 'line_total_cents', nombre: t.p04.importe, valor: (x: Renglon) => x.line_total_cents, pintar: (x: Renglon) => dinero(x.line_total_cents), alinear: 'derecha' },
            { id: 'household_share_cents', nombre: t.p04.parteHogar, valor: (x: Renglon) => x.household_share_cents ?? 0, pintar: (x: Renglon) => dinero(x.household_share_cents ?? 0), alinear: 'derecha' },
          ],
        }} />
      )}</ConDatos>
    </Bloque>
  )
}
