// P-19 Compra a crédito (U1 §2.1, §2.2, §3.21: N-215..N-219). Alta con vista previa del cálculo (motor S3, servidor)
// y detalle con partes por persona y pagos con estado por parte. Sin conexión: «Requiere conexión» (U1 §2.4).
import { useEffect, useMemo, useState } from 'react'
import { CircleCheck, Pencil, Receipt, Trash2 } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { conConsulta, navegar } from '../../app/enrutador'
import { usarConsulta } from '../../api/consulta'
import { nuevaClave, pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { ErrorApi } from '../../api/errores'
import { aCentavos, dinero, fecha, fechaConDia, hoyISO } from '../../api/formatos'
import type { Cuerpo, Esquemas } from '../../api/tipos'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { useAvisos } from '../../componentes/Avisos'
import { ConDatos, Franja } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { itemsDe, Barra, Datos, HojaPago, Seccion, Selector, agrupar, mensajeError, textoPartes, usarPersonas, type Personas } from './comun'
import { ESTADO_COMPRA, ESTADO_PARTE, TIPO_COMPRA, plazo, t } from './textos'

const k = t.compra
type Alta = Cuerpo<'/api/v1/credit-purchases', 'post'>
type Previa = Esquemas['CreditPurchasePreview']

export default function P19Compra({ params, consulta }: PropsPantalla) {
  const enLinea = useConexion()
  if (!enLinea) return <Franja rol="sinConexion">{textos.conexion.requiere}</Franja>
  return params.compra === 'nueva' ? <AltaCompra creditoInicial={Number(params.id)} consulta={consulta} /> : <DetalleCompra credito={Number(params.id)} id={Number(params.compra)} consulta={consulta} />
}

function AltaCompra({ creditoInicial, consulta }: { creditoInicial: number; consulta: URLSearchParams }) {
  const avisar = useAvisos()
  const per = usarPersonas()
  const lineas = usarConsulta('get', '/api/v1/credits', { consulta: { limit: 200 } })
  const [f, setF] = useState({ credito: String(creditoInicial || ''), tipo: 'installment_plan', comercio: consulta.get('comercio') ?? '', fecha: consulta.get('fecha') ?? hoyISO(), capital: consulta.get('monto') ?? '', intereses: '', plazo: consulta.get('plazo') ?? '', promocion: consulta.get('promocion') ?? '', nota: '' })
  const [partes, setPartes] = useState<Record<number, string>>({})
  const [errores, setErrores] = useState<Record<string, string>>({})
  const [clave, setClave] = useState(nuevaClave)
  const promos = usarConsulta('get', '/api/v1/promotions', { consulta: { credit_id: [Number(f.credito)], valid_on: f.fecha } }, !!f.credito && /^\d{4}-\d{2}-\d{2}$/.test(f.fecha))
  const cambiar = (c: Partial<typeof f>) => setF({ ...f, ...c })
  const cuerpo = useMemo((): Alta | null => {
    const capital = aCentavos(f.capital), n = Number(f.plazo), cred = Number(f.credito)
    if (!capital || capital <= 0 || !Number.isInteger(n) || n < 1 || !cred || !f.fecha) return null
    const interes = f.intereses.trim() ? aCentavos(f.intereses) : null
    const lista = Object.entries(partes).filter(([, v]) => v.trim()).map(([p, v]) => ({ person_id: Number(p), principal_cents: aCentavos(v) ?? 0 }))
    return { credit_id: cred, kind: f.tipo as Alta['kind'], purchase_date: f.fecha, principal_cents: capital, term_count: n,
      ...(f.comercio.trim() ? { description: f.comercio.trim() } : {}), ...(f.nota.trim() ? { notes: f.nota.trim() } : {}),
      ...(interes != null ? { interest_cents: interes, interest_source: 'quoted' as const } : f.promocion ? { promotion_id: Number(f.promocion), interest_source: 'promotion' as const } : { interest_source: 'term_table' as const }),
      ...(f.promocion ? { promotion_id: Number(f.promocion) } : {}), ...(lista.length ? { shares: lista } : {}) }
  }, [f, partes])
  const sumaMal = !!cuerpo?.shares && cuerpo.shares.reduce((a, s) => a + s.principal_cents, 0) !== cuerpo.principal_cents
  const [previa, setPrevia] = useState<{ estado: 'vacia' | 'calculando' | 'error' | 'lista'; datos?: Previa }>({ estado: 'vacia' })
  useEffect(() => {
    if (!cuerpo || sumaMal) return setPrevia({ estado: 'vacia' })
    setPrevia((p) => ({ ...p, estado: 'calculando' }))
    const ctl = new AbortController()
    const h = setTimeout(() => pedir('post', '/api/v1/credit-purchases/preview', { cuerpo, senal: ctl.signal }).then((datos) => setPrevia({ estado: 'lista', datos }), (e) => { if ((e as DOMException).name !== 'AbortError') setPrevia({ estado: 'error' }) }), 350)
    return () => { clearTimeout(h); ctl.abort() }
  }, [JSON.stringify(cuerpo), sumaMal]) // eslint-disable-line react-hooks/exhaustive-deps
  const guardar = async () => {
    const e: Record<string, string> = {}
    if (!aCentavos(f.capital)) e.principal_cents = k.capitalInvalido
    if (!(Number(f.plazo) >= 1)) e.term_count = k.plazoInvalido
    if (sumaMal) e.shares = k.sumaPartes
    setErrores(e)
    if (Object.keys(e).length || !cuerpo || previa.estado !== 'lista') return
    try {
      const r = await pedir('post', '/api/v1/credit-purchases', { cuerpo, idempotencia: clave })
      setClave(nuevaClave())
      avisar({ texto: k.guardada, rol: 'exito' })
      navegar(rutaDe('P-19', { id: r.credit_id, compra: r.id }), { reemplazar: true })
    } catch (err) {
      if (err instanceof ErrorApi && err.campos.length) setErrores(Object.fromEntries(err.campos.map((c) => [c.field, c.message])))
      else avisar({ texto: mensajeError(err), rol: 'error' })
    }
  }
  const motivoNo = previa.estado === 'error' ? k.sinCalculo : previa.estado !== 'lista' ? k.completa : undefined
  return (
    <div className="grid grid-cols-1 gap-6 tabular-nums desktop:grid-cols-2">
      <Seccion titulo={k.datos}>
        <div className="flex flex-col gap-3">
          <Selector etiqueta={t.creditos.linea} valor={f.credito} alCambiar={(v) => cambiar({ credito: v, promocion: '' })} opciones={itemsDe<Esquemas['Credit']>(lineas.datos).map((l) => ({ valor: String(l.id), texto: l.name }))} />
          <Selector etiqueta={k.tipo} valor={f.tipo} alCambiar={(v) => cambiar({ tipo: v })} opciones={Object.entries(TIPO_COMPRA).map(([valor, texto]) => ({ valor, texto }))} />
          <Campo etiqueta={k.comercio} value={f.comercio} error={errores.description} onChange={(e) => cambiar({ comercio: e.target.value })} />
          <Campo etiqueta={k.fecha} type="date" value={f.fecha} error={errores.purchase_date} onChange={(e) => cambiar({ fecha: e.target.value, promocion: '' })} />
          <Campo etiqueta={k.capital} inputMode="decimal" value={f.capital} error={errores.principal_cents} onChange={(e) => cambiar({ capital: e.target.value })} />
          <Campo etiqueta={k.intereses} inputMode="decimal" ayuda={k.interesesAyuda} value={f.intereses} error={errores.interest_cents} onChange={(e) => cambiar({ intereses: e.target.value })} />
          <Campo etiqueta={k.plazo} inputMode="numeric" value={f.plazo} error={errores.term_count} onChange={(e) => cambiar({ plazo: e.target.value })} />
          <Selector etiqueta={k.promocion} valor={f.promocion} alCambiar={(v) => cambiar({ promocion: v })} opciones={[{ valor: '', texto: k.sinPromocion }, ...itemsDe<Esquemas['Promotion']>(promos.datos).map((p) => ({ valor: String(p.id), texto: `${p.merchant_name} · ${plazo(p.term_count)}` }))]} />
          <fieldset className="flex flex-col gap-2 rounded-lg border border-border-divider p-3"><legend className="px-1 text-small font-medium">{k.reparto}</legend>
            <p className="text-caption text-text-secondary">{k.repartoAyuda}</p>
            {per.lista.filter((p) => p.is_active).map((p) => <Campo key={p.id} etiqueta={p.display_name} inputMode="decimal" value={partes[p.id] ?? ''} onChange={(e) => setPartes({ ...partes, [p.id]: e.target.value })} />)}
            {(sumaMal || errores.shares) && <p role="alert" className="text-small text-role-error">{k.sumaPartes}</p>}
          </fieldset>
          <Campo etiqueta={k.nota} value={f.nota} onChange={(e) => cambiar({ nota: e.target.value })} />
          <Boton variante="primario" tam="hoja" deshabilitadoPor={motivoNo} onClick={guardar}>{k.guardar}</Boton>
        </div>
      </Seccion>
      <Seccion titulo={k.vistaPrevia}>
        <div aria-live="polite">{previa.estado === 'calculando' ? <p className="text-small text-text-secondary">{k.calculando}</p> : previa.estado === 'error' ? <Franja rol="error" alerta>{k.sinCalculo}</Franja> : previa.estado === 'vacia' || !previa.datos ? <p className="text-small text-text-secondary">{k.completa}</p> : <Previa p={previa.datos} per={per} />}</div>
      </Seccion>
    </div>
  )
}
function Previa({ p, per }: { p: Previa; per: Personas }) {
  return <div className="flex flex-col gap-4">
    <Datos filas={[[t.creditos.capital, dinero(p.principal_cents)], [t.creditos.intereses, dinero(p.interest_cents)], [t.comun.total, dinero(p.total_cents)]]} />
    {p.promotion_applies === true && <Franja rol="info">{k.promoAplica}</Franja>}
    {p.promotion_applies === false && <Franja rol="advertencia">{k.promoNoAplica(p.promotion_issue ?? '')}</Franja>}
    <Seccion titulo={k.partes}><Datos filas={p.shares.map((s) => [per.nombre(s.person_id), `${dinero(s.total_cents)} · ${k.pagoQuincena.toLowerCase()} ${dinero(s.fortnight_payment_cents)}`])} /></Seccion>
    <Seccion titulo={k.calendario}><Datos filas={p.schedule.map((r) => [fechaConDia(r.due_date), `${dinero(r.amount_cents)} · ${textoPartes(agrupar(r.shares, (x) => x.amount_cents, per))}`])} /></Seccion>
  </div>
}

function DetalleCompra({ credito, id, consulta }: { credito: number; id: number; consulta: URLSearchParams }) {
  const avisar = useAvisos()
  const per = usarPersonas()
  const compra = usarConsulta('get', '/api/v1/credit-purchases/{credit_purchase_id}', { ruta: { credit_purchase_id: id } })
  const linea = usarConsulta('get', '/api/v1/credits/{credit_id}', { ruta: { credit_id: credito } })
  const pagos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { credit_purchase_id: id, limit: 200 } })
  const [hoja, setHoja] = useState<'editar' | 'eliminar' | null>(null)
  const [ed, setEd] = useState({ comercio: '', nota: '' })
  const recargar = () => { pagos.recargar(); compra.recargar() }
  const guardar = async () => {
    try { await pedir('patch', '/api/v1/credit-purchases/{credit_purchase_id}', { ruta: { credit_purchase_id: id }, cuerpo: { description: ed.comercio.trim() || null, notes: ed.nota.trim() || null } }); setHoja(null); compra.recargar(); avisar({ texto: k.guardado, rol: 'exito' }) } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
  }
  const eliminar = async () => {
    try { await pedir('delete', '/api/v1/credit-purchases/{credit_purchase_id}', { ruta: { credit_purchase_id: id } }); avisar({ texto: k.eliminada, rol: 'exito' }); navegar(rutaDe('P-18', { id: credito }), { reemplazar: true }) } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
  }
  const pendiente = itemsDe<Esquemas['ScheduledPayment']>(pagos.datos).find((p) => p.status !== 'paid')
  return (
    <div className="flex flex-col gap-6 tabular-nums">
      <ConDatos consulta={compra}>{(x) => <>
        <Barra>
          {x.expense_id && <Boton variante="fantasma" onClick={() => navegar(rutaDe('P-04', { id: x.expense_id! }))}><Receipt aria-hidden className="size-5" />{k.verGasto}</Boton>}
          <Boton variante="secundario" onClick={() => { setEd({ comercio: x.description ?? '', nota: x.notes ?? '' }); setHoja('editar') }}><Pencil aria-hidden className="size-5" />{t.comun.editar}</Boton>
          <Boton variante="secundario" onClick={() => setHoja('eliminar')}><Trash2 aria-hidden className="size-5" />{t.comun.eliminar}</Boton>
          <Boton variante="primario" deshabilitadoPor={pendiente ? undefined : t.creditos.sinPagosPendientes} onClick={() => pendiente && navegar(conConsulta({ hoja: 'pago', pago: String(pendiente.id) }))}><CircleCheck aria-hidden className="size-5" />{t.creditos.marcarPago}</Boton>
        </Barra>
        <div className="grid grid-cols-1 gap-6 desktop:grid-cols-2">
          <Seccion titulo={k.datos}><Datos filas={[[t.creditos.linea, linea.datos?.name ?? ''], [t.creditos.comercio, x.description ?? t.comun.sinDato], [t.comun.fecha, fecha(x.purchase_date)], [t.creditos.capital, dinero(x.principal_cents)], [t.creditos.intereses, dinero(x.interest_cents)], [t.comun.total, dinero(x.total_cents)], [t.creditos.plazo, plazo(x.term_count, x.payment_frequency)], [t.creditos.promocion, x.promotion_id ? (x.rate_pct ? t.creditos.tasaCapturada(x.rate_pct) : TIPO_COMPRA[x.kind]) : k.sinPromocion], [t.comun.estado, ESTADO_COMPRA[x.status]], ...(x.notes ? [[t.comun.nota, x.notes] as [string, string]] : [])]} /></Seccion>
          <Seccion titulo={k.partes}><Datos filas={(x.shares ?? []).map((s) => [per.nombre(s.person_id), `${dinero(s.principal_cents)} + ${dinero(s.interest_cents)} = ${dinero(s.total_cents)}${s.fortnight_payment_cents != null ? ` · ${k.pagoQuincena.toLowerCase()} ${dinero(s.fortnight_payment_cents)}` : ''}`])} /></Seccion>
        </div>
      </>}</ConDatos>
      <Seccion titulo={k.pagos}><ConDatos consulta={pagos}>{(d) => <VistaBD def={{ clave: 'P-19-pagos', filas: d.items, idFila: (p) => p.id, vistas: ['tabla'], titulo: (p) => fecha(p.due_date), alAbrir: (p) => navegar(conConsulta({ hoja: 'pago', pago: String(p.id) })),
        propiedades: [{ id: 'due', nombre: t.comun.fecha, valor: (p) => p.due_date, pintar: (p) => fechaConDia(p.due_date) }, { id: 'amount', nombre: t.comun.monto, valor: (p) => p.amount_cents, pintar: (p) => dinero(p.amount_cents), alinear: 'derecha' }, { id: 'partes', nombre: t.creditos.partes, valor: (p) => (p.shares ?? []).map((s) => `${per.parte(s.person_id)} ${dinero(s.amount_cents)} · ${ESTADO_PARTE[s.status].toLowerCase()}`).join(' · ') }, { id: 'principal', nombre: t.creditos.capitalRepuesto, valor: (p) => p.principal_cents, pintar: (p) => dinero(p.principal_cents), alinear: 'derecha' }, { id: 'status', nombre: t.comun.estado, valor: (p) => ESTADO_PARTE[p.status ?? 'pending'] }] }} />}</ConDatos></Seccion>
      <Hoja abierta={hoja === 'editar'} alCambiar={(a) => !a && setHoja(null)} titulo={k.editar} pie={<Boton variante="primario" tam="hoja" onClick={guardar}>{textos.app.guardar}</Boton>}>
        <div className="flex flex-col gap-3"><Campo etiqueta={k.comercio} value={ed.comercio} onChange={(e) => setEd({ ...ed, comercio: e.target.value })} /><Campo etiqueta={k.nota} value={ed.nota} onChange={(e) => setEd({ ...ed, nota: e.target.value })} /></div>
      </Hoja>
      <Hoja abierta={hoja === 'eliminar'} alCambiar={(a) => !a && setHoja(null)} titulo={k.eliminar} pie={<div className="flex gap-2"><Boton className="flex-1" onClick={() => setHoja(null)}>{textos.app.cancelar}</Boton><Boton className="flex-1" variante="destructivo" onClick={eliminar}>{t.comun.eliminar}</Boton></div>}>
        <p className="text-body">{k.confirmarEliminar}</p>
      </Hoja>
      <HojaPago pagos={pagos.datos?.items ?? []} per={per} alCambiar={recargar} consulta={consulta} />
    </div>
  )
}
