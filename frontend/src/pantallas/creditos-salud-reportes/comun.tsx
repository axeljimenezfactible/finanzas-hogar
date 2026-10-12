// Piezas compartidas del área: personas y partes (U1 §0.5), secciones, pares dato-valor, «Falta capturar» (K18),
// PDF (descargar y compartir), agenda de pagos (N-192/N-208) y la hoja «Marcar pago» (G-06, N-210).
import { Fragment, useState, type ReactNode } from 'react'
import { usarConsulta } from '../../api/consulta'
import { pedir, nuevaClave } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { ErrorApi } from '../../api/errores'
import { aCentavos, dinero, fecha, fechaConDia, hoyISO, mes } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { Campo, estilosControl } from '../../componentes/Campo'
import { useAvisos } from '../../componentes/Avisos'
import { Hoja } from '../../componentes/Hoja'
import { cn } from '../../componentes/utils'
import { Enlace, conConsulta, navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { textos } from '../../textos/es-MX'
import { ESTADO_PARTE, t } from './textos'

export type Pago = Esquemas['ScheduledPayment']
/** `Page & { items?: T[] }` (openapi-typescript) deja `items` como `unknown[] & T[]`; esto devuelve `T[]`. */
export const itemsDe = <T,>(d: { items: unknown[] } | undefined) => (d?.items ?? []) as T[]
type Parte = NonNullable<Pago['shares']>[number]

export function usarPersonas() {
  const c = usarConsulta('get', '/api/v1/persons')
  const lista = itemsDe<Esquemas['Person']>(c.datos)
  const nombre = (id: number) => lista.find((p) => p.id === id)?.display_name ?? t.comun.persona(id)
  const externa = (id: number) => lista.find((p) => p.id === id)?.participation === 'external'
  /** U1 §0.5: la parte de la persona externa lleva su nombre (Karen); las de Axel y Ximena son «Hogar». */
  const parte = (id: number) => (externa(id) ? nombre(id) : t.comun.hogar)
  return { lista, nombre, externa, parte, titular: lista.find((p) => p.participation === 'holder'), externaPrincipal: lista.find((p) => p.participation === 'external') }
}
export type Personas = ReturnType<typeof usarPersonas>

/** Suma las partes por etiqueta (Hogar, Karen) en orden fijo: Hogar primero. */
export function agrupar<T extends { person_id: number }>(partes: T[] | undefined, monto: (p: T) => number, per: Personas) {
  const m = new Map<string, number>()
  for (const p of partes ?? []) m.set(per.parte(p.person_id), (m.get(per.parte(p.person_id)) ?? 0) + monto(p))
  return [...m.entries()].sort(([a], [b]) => (a === t.comun.hogar ? -1 : b === t.comun.hogar ? 1 : a.localeCompare(b, 'es')))
}
export const textoPartes = (partes: [string, number][]) => partes.map(([k, v]) => `${k} ${dinero(v)}`).join(' · ')

export function Seccion({ titulo, children, accion }: { titulo: string; children: ReactNode; accion?: ReactNode }) {
  return (
    <section aria-label={titulo} className="flex flex-col gap-3">
      <div className="flex flex-wrap items-center justify-between gap-2"><h2 className="text-title-3 font-semibold text-text-default">{titulo}</h2>{accion}</div>
      {children}
    </section>
  )
}
export const Rejilla = ({ children, anchas }: { children: ReactNode; anchas?: boolean }) => (
  <div className={cn('grid grid-cols-1 gap-3 tablet:grid-cols-2', anchas ? 'desktop:grid-cols-2' : 'desktop:grid-cols-3')}>{children}</div>
)
export function Datos({ filas, className }: { filas: [string, ReactNode][]; className?: string }) {
  return (
    <dl className={cn('grid grid-cols-2 gap-x-4 gap-y-2 rounded-lg border border-border-divider bg-surface-card p-4 text-small', className)}>
      {filas.map(([k, v]) => <Fragment key={k}><dt className="text-text-secondary">{k}</dt><dd className="text-right tabular-nums text-text-default">{v}</dd></Fragment>)}
    </dl>
  )
}
export const Barra = ({ children }: { children: ReactNode }) => <div className="flex flex-wrap items-center justify-end gap-2">{children}</div>
export function Selector({ etiqueta, valor, alCambiar, opciones, className }: { etiqueta: string; valor: string; alCambiar: (v: string) => void; opciones: { valor: string; texto: string }[]; className?: string }) {
  return (
    <label className={cn('flex flex-col gap-1 text-small font-medium text-text-default', className)}>{etiqueta}
      <select value={valor} onChange={(e) => alCambiar(e.target.value)} className={estilosControl}>{opciones.map((o) => <option key={o.valor} value={o.valor}>{o.texto}</option>)}</select>
    </label>
  )
}
/** Ajustes › Créditos y condiciones (P-26.05) en esa línea. `?credito=` es PROVISIONAL (bloqueo WF3-creditos-salud-reportes). */
export const rutaCondiciones = (credito?: number) => rutaDe('P-26', { seccion: '05' }) + (credito ? `?credito=${credito}` : '')
export function FaltaCapturar({ credito, que }: { credito: number; que?: string[] }) {
  return (
    <span className="inline-flex flex-wrap items-center justify-end gap-1 text-role-warning">
      <span title={que?.map((q) => t.creditos.faltan[q] ?? q).join(', ')}>{t.creditos.faltaCapturar}</span>·
      <Enlace a={rutaCondiciones(credito)} className="inline-flex min-h-11 items-center text-text-link underline">{t.creditos.capturar}</Enlace>
    </span>
  )
}
export const mensajeError = (e: unknown) => (e instanceof ErrorApi ? e.message : textos.errores.generico)
/** Escucha `fh:datos-cambiaron` (peso o gasto guardado en P-03/P-22). */
export function alCambiarDatos(f: () => void) { window.addEventListener('fh:datos-cambiaron', f); return () => window.removeEventListener('fh:datos-cambiaron', f) }

async function traerPdf(id: number) { return (await pedir('get', '/api/v1/reports/{report_id}/file', { ruta: { report_id: id } })) as unknown as Blob }
function descargar(blob: Blob, nombre: string) {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url; a.download = nombre; a.click()
  setTimeout(() => URL.revokeObjectURL(url), 1000)
}
/** Descargar PDF (N-249) y Compartir con la hoja de compartir de Android (F-5); sin Web Share, descarga. */
export function AccionesPdf({ id, nombre, principal }: { id: number; nombre: string; principal?: boolean }) {
  const avisar = useAvisos()
  const enLinea = useConexion()
  const sin = enLinea ? undefined : textos.conexion.requiere
  const bajar = async (compartir: boolean) => {
    try {
      const blob = await traerPdf(id)
      const archivo = new File([blob], nombre, { type: 'application/pdf' })
      if (compartir && navigator.canShare?.({ files: [archivo] })) await navigator.share({ files: [archivo], title: nombre })
      else descargar(blob, nombre)
    } catch (e) { if ((e as DOMException).name !== 'AbortError') avisar({ texto: t.comun.noSePudoDescargar, rol: 'error' }) }
  }
  return (
    <div className="flex flex-wrap gap-2">
      <Boton variante={principal ? 'primario' : 'secundario'} deshabilitadoPor={sin} onClick={() => bajar(false)}>{t.comun.descargar}</Boton>
      <Boton variante="secundario" deshabilitadoPor={sin} onClick={() => bajar(true)}>{t.comun.compartir}</Boton>
    </div>
  )
}

/** Agenda de pagos por fecha (N-192, N-208). Tocar un pago abre su vista previa con «Marcar pago». */
export function Agenda({ pagos, per, nombreLinea }: { pagos: Pago[]; per: Personas; nombreLinea: (id?: number) => string }) {
  if (!pagos.length) return <p className="py-4 text-small text-text-secondary">{t.creditos.sinPagos}</p>
  return (
    <ul className="flex flex-col divide-y divide-border-divider rounded-lg border border-border-divider bg-surface-card">
      {pagos.map((p) => (
        <li key={p.id}>
          <button type="button" data-prueba="pago-agenda" onClick={() => navegar(conConsulta({ hoja: 'pago', pago: String(p.id) }))}
            className="flex min-h-11 w-full flex-col gap-1 px-4 py-3 text-left hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus">
            <span className="flex flex-wrap justify-between gap-2 text-body text-text-default"><span>{fechaConDia(p.due_date)} · {nombreLinea(p.credit_id)}</span><span className="tabular-nums font-medium">{dinero(p.amount_cents)}</span></span>
            <span className="text-small tabular-nums text-text-secondary">{(p.shares ?? []).length ? agruparEstado(p.shares!, per) : ESTADO_PARTE[p.status ?? 'pending']}</span>
          </button>
        </li>
      ))}
    </ul>
  )
}
const agruparEstado = (partes: Parte[], per: Personas) => agrupar(partes, (s) => s.amount_cents, per).map(([k, v]) => `${k} ${dinero(v)}`).join(' · ') + ' · ' + [...new Set(partes.map((s) => ESTADO_PARTE[s.status]))].join(' / ').toLowerCase()

/** Hoja «Marcar pago» (G-06, N-210): una parte por persona, fecha real (hoy) y monto programado por defecto; Deshacer desmarca. */
export function HojaPago({ pagos, per, alCambiar, consulta }: { pagos: Pago[]; per: Personas; alCambiar: () => void; consulta: URLSearchParams }) {
  const abierta = consulta.get('hoja') === 'pago'
  const pago = pagos.find((p) => String(p.id) === consulta.get('pago'))
  const cerrar = () => navegar(conConsulta({ hoja: null, pago: null }))
  return (
    <Hoja abierta={abierta} alCambiar={(a) => !a && cerrar()} titulo={t.pago.titulo}>
      {pago ? <FormPago key={pago.id} pago={pago} per={per} alCambiar={alCambiar} alCerrar={cerrar} /> : (
        <div className="flex flex-col gap-2">
          <p className="text-small text-text-secondary">{pagos.some((p) => p.status !== 'paid') ? t.pago.elegir : t.pago.sinPendientes}</p>
          {pagos.filter((p) => p.status !== 'paid').slice(0, 10).map((p) => <Boton key={p.id} className="justify-between" onClick={() => navegar(conConsulta({ pago: String(p.id) }))}><span>{fechaConDia(p.due_date)}</span><span className="tabular-nums">{dinero(p.amount_cents)}</span></Boton>)}
        </div>
      )}
    </Hoja>
  )
}
function FormPago({ pago, per, alCambiar, alCerrar }: { pago: Pago; per: Personas; alCambiar: () => void; alCerrar: () => void }) {
  const avisar = useAvisos()
  const enLinea = useConexion()
  const [fecha, setFecha] = useState(hoyISO())
  const [montos, setMontos] = useState<Record<number, string>>(() => Object.fromEntries((pago.shares ?? []).map((s) => [s.person_id, ((s.amount_cents - s.paid_cents) / 100).toFixed(2)])))
  const [error, setError] = useState<Record<number, string>>({})
  const [claves] = useState<Record<number, string>>(() => Object.fromEntries((pago.shares ?? []).map((s) => [s.person_id, nuevaClave()])))
  const marcar = async (s: Parte) => {
    const c = aCentavos(montos[s.person_id] ?? '')
    if (!c || c <= 0) return setError({ ...error, [s.person_id]: t.pago.montoInvalido })
    try {
      const r = await pedir('post', '/api/v1/scheduled-payments/{scheduled_payment_id}/payments', { ruta: { scheduled_payment_id: pago.id }, cuerpo: { payer_person_id: s.person_id, paid_date: fecha, amount_cents: c }, idempotencia: claves[s.person_id] })
      alCambiar()
      alCerrar() // cierra la hoja: el aviso con «Deshacer» queda al alcance (fuera del diálogo modal)
      avisar({ texto: t.pago.marcado, rol: 'exito', accion: { texto: textos.app.deshacer, alHacer: async () => { await pedir('delete', '/api/v1/credit-payments/{credit_payment_id}', { ruta: { credit_payment_id: r.id } }); alCambiar(); avisar({ texto: t.pago.desmarcado }) } } })
    } catch (e) { setError({ ...error, [s.person_id]: mensajeError(e) }) }
  }
  return (
    <div className="flex flex-col gap-4">
      <Datos filas={[[t.comun.fecha, fechaConDia(pago.due_date)], [t.comun.monto, dinero(pago.amount_cents)], [t.pago.capitalRepuesto, dinero(pago.principal_cents)], [t.creditos.interesesCorto, dinero(pago.interest_cents)]]} />
      <Campo etiqueta={t.pago.fechaReal} type="date" value={fecha} onChange={(e) => setFecha(e.target.value)} />
      {(pago.shares ?? []).map((s) => (
        <div key={s.person_id} className="flex flex-col gap-2 rounded-lg border border-border-divider p-3">
          <p className="flex justify-between text-small"><span className="font-medium">{t.pago.parteDe(per.parte(s.person_id))}</span><span className="tabular-nums">{dinero(s.amount_cents)} · {ESTADO_PARTE[s.status]}</span></p>
          {s.status !== 'paid' && <>
            <Campo etiqueta={t.pago.montoParte(per.parte(s.person_id))} inputMode="decimal" value={montos[s.person_id] ?? ''} error={error[s.person_id]} onChange={(e) => setMontos({ ...montos, [s.person_id]: e.target.value })} />
            <Boton variante="primario" deshabilitadoPor={enLinea ? undefined : textos.conexion.requiere} onClick={() => marcar(s)}>{t.pago.marcarParte(per.parte(s.person_id))}</Boton>
          </>}
        </div>
      ))}
    </div>
  )
}

/** Suma días a una fecha ISO (AAAA-MM-DD), sin zonas: aritmética de calendario. */
export function sumarDias(iso: string, n: number) { const d = new Date(`${iso}T12:00:00Z`); d.setUTCDate(d.getUTCDate() + n); return d.toISOString().slice(0, 10) }
/** Agenda: pagos sin pagar con fecha pasada y los de hoy a 60 días (U1 N-192), por fecha. */
export const enAgenda = (pagos: Pago[], hoy = hoyISO()) => pagos.filter((p) => (p.due_date < hoy ? p.status !== 'paid' : p.due_date <= sumarDias(hoy, 60))).sort((a, b) => a.due_date.localeCompare(b.due_date) || a.id - b.id)
/** Periodo de un reporte: un mes calendario como «sep-2026»; si no, «1-sep-2026 al 15-sep-2026». */
export function periodoTexto(a: string, b: string) {
  const fin = new Date(`${b}T12:00:00Z`); fin.setUTCDate(fin.getUTCDate() + 1)
  return a.endsWith('-01') && fin.getUTCDate() === 1 && a.slice(0, 7) === b.slice(0, 7) ? mes(a.slice(0, 7)) : `${fecha(a)} al ${fecha(b)}`
}
