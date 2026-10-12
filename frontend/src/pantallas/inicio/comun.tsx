// Piezas comunes de P-01 (U3 §1): tarjeta-enlace, insignias, estados por bloque, estado de un pago (§1.4) y datos compartidos.
import { Component, createContext, useContext, useEffect, useState, type ReactNode } from 'react'
import { ListFilter, X } from 'lucide-react'
import { Hoja } from '../../componentes/Hoja'
import { ChipSeleccionable } from '../../componentes/Chip'
import { textos } from '../../textos/es-MX'
import { ChevronRight, CircleAlert, CircleCheck, CircleDashed, CircleDot, CalendarClock, CalendarX, Info, RefreshCw, TriangleAlert } from 'lucide-react'
import { Enlace } from '../../app/enrutador'
import { usarConsulta, type EstadoConsulta } from '../../api/consulta'
import { horaCorta, hoyISO } from '../../api/formatos'
import { Boton } from '../../componentes/Boton'
import { Cargando, Esqueleto } from '../../componentes/Estados'
import { cn } from '../../componentes/utils'
import type { Esquemas, Respuesta } from '../../api/tipos'
import { difDias, quincenaDe, mesDe } from './fechas'
import { textosInicio as t } from './textos'

export const entero = (n: number) => String(Math.round(n)).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
/** Décimas de punto porcentual → «47.5 %» (C1 `progress_tenths_pct`). */
export const decimas = (n: number) => `${Math.trunc(n / 10)}.${Math.abs(n % 10)} %`

type RolInsignia = 'advertencia' | 'info' | 'exito' | 'neutro'
const INSIGNIA: Record<RolInsignia, string> = { advertencia: 'bg-role-warning-surface text-role-warning', info: 'bg-role-info-surface text-role-info', exito: 'bg-role-success-surface text-role-success', neutro: 'bg-surface-raised text-text-secondary border border-border-divider' }
export function Insignia({ rol, icono, children, etiqueta }: { rol: RolInsignia; icono: ReactNode; children: ReactNode; etiqueta?: string }) {
  return <span aria-label={etiqueta} className={cn('inline-flex h-6 w-max items-center gap-1 rounded-sm px-2 text-small font-medium [&_svg]:size-5', INSIGNIA[rol])}>{icono}{children}</span>
}
export const IconoAviso = () => <TriangleAlert aria-hidden />
export const IconoInfo = () => <Info aria-hidden />
export const IconoOk = () => <CircleCheck aria-hidden />

/** Pie de tarjeta o panel: texto de enlace con chevron (la zona entera es el enlace). */
export const PieEnlace = ({ texto }: { texto: string }) => (
  <span className="mt-auto inline-flex min-h-11 items-center gap-1 text-small font-medium text-text-link underline underline-offset-2">{texto}<ChevronRight aria-hidden className="size-5" /></span>
)
export function EnlaceAccion({ a, texto }: { a: string; texto: string }) {
  return <Enlace a={a} className="inline-flex min-h-11 items-center gap-1 rounded-sm text-small font-medium text-text-link underline underline-offset-2 hover:text-text-link-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus">{texto}<ChevronRight aria-hidden className="size-5" /></Enlace>
}

export function FallaBloque({ texto, alReintentar }: { texto: string; alReintentar?: () => void }) {
  return (
    <div role="alert" className="flex flex-col gap-2 rounded-md bg-role-error-surface p-3 text-small text-text-default">
      <span className="flex items-start gap-2"><CircleAlert aria-hidden className="size-5 shrink-0 text-role-error" />{texto}</span>
      {alReintentar && <Boton variante="secundario" className="w-max" onClick={alReintentar}><RefreshCw aria-hidden className="size-5" />{t.acciones.reintentar}</Boton>}
    </div>
  )
}
export const CargadoA = ({ desde }: { desde: number | null }) => (desde ? <p className="text-caption text-text-secondary">{t.comun.cargadoA(horaCorta(desde))}</p> : null)

/** Aísla un bloque: si sus datos no tienen la forma de C1, muestra su falla y las demás tarjetas siguen (U3 §1.2). */
export class Aislado extends Component<{ texto: string; alReintentar?: () => void; children: ReactNode }, { fallo: boolean }> {
  state = { fallo: false }
  componentDidCatch() { this.setState({ fallo: true }) }
  render() { return this.state.fallo ? <FallaBloque texto={this.props.texto} alReintentar={() => { this.setState({ fallo: false }); this.props.alReintentar?.() }} /> : this.props.children }
}

/** Estados de una lectura (U3 §1.3): esqueleto con la forma final, aviso a los 10 s, falla con su texto, sin conexión. */
export function Estados<T>({ consulta, falla, esqueleto, children }: { consulta: EstadoConsulta<T>; falla: string; esqueleto: ReactNode; children: (d: T) => ReactNode }) {
  if (consulta.cargando && consulta.datos === undefined) return <EsperaLarga alReintentar={consulta.recargar}>{esqueleto}</EsperaLarga>
  if (consulta.nuncaCargada) return <p role="status" className="rounded-md bg-role-warning-surface p-3 text-small text-text-default">{t.comun.nuncaCargada}</p>
  if (consulta.error || consulta.datos === undefined) return <FallaBloque texto={falla} alReintentar={consulta.recargar} />
  return <Aislado texto={falla} alReintentar={consulta.recargar}><CargadoA desde={consulta.desdeCopia} /><Pintar f={children} d={consulta.datos} /></Aislado>
}
const EsperaLarga = Cargando
/** Llama a la función de pintado dentro del límite de Aislado (un dato con otra forma no tumba la página). */
function Pintar<T>({ f, d }: { f: (d: T) => ReactNode; d: T }) { return <>{f(d)}</> }
export const EsqueletoTarjeta = () => <div className="flex flex-col gap-2"><Esqueleto className="w-1/2 bg-skeleton-on-card" /><Esqueleto className="h-8 w-2/3 bg-skeleton-on-card" /><Esqueleto className="w-5/6 bg-skeleton-on-card" /></div>
export const EsqueletoFilas = ({ n }: { n: number }) => <div className="flex flex-col gap-2">{Array.from({ length: n }, (_, i) => <Esqueleto key={i} className="h-10 bg-skeleton-on-card" />)}</div>

/** Recarga al guardar un gasto o un peso (fh:datos-cambiaron, CONVENCIONES §4). */
export function usarRecargaAlCambiar(...recargas: (() => void)[]) {
  useEffect(() => {
    const f = () => recargas.forEach((r) => r())
    window.addEventListener('fh:datos-cambiaron', f)
    return () => window.removeEventListener('fh:datos-cambiaron', f)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, recargas)
}

// ---- Datos compartidos de Inicio: se leen una sola vez por visita (A-02: las alertas se leen una vez) ----
export type Persona = Esquemas['Person']
type Lista<T> = { items: T[]; next_cursor?: string | null }
export interface Compartido {
  hoy: string
  alertas: EstadoConsulta<{ parameters: { 'alerts.payment_due_days_before': { value: number[] | null } } }>
  personas: EstadoConsulta<Lista<Persona>>
  creditos: EstadoConsulta<Lista<Esquemas['Credit']>>
  tiendas: EstadoConsulta<Lista<Esquemas['Store']>>
  categorias: EstadoConsulta<Lista<Esquemas['Category']>>
  proximo: EstadoConsulta<Respuesta<'/api/v1/dashboard/next-payment', 'get'>>
  flujo: EstadoConsulta<Respuesta<'/api/v1/dashboard/fortnight-flow', 'get'>>
  plan: EstadoConsulta<Respuesta<'/api/v1/dashboard/month-vs-plan', 'get'>>
}
const Ctx = createContext<Compartido | null>(null)
export const useInicio = () => useContext(Ctx)!
export function ProveedorInicio({ children }: { children: ReactNode }) {
  const hoy = hoyISO()
  const alertas = usarConsulta('get', '/api/v1/settings/alerts')
  const personas = usarConsulta('get', '/api/v1/persons', { consulta: { limit: 200 } })
  const creditos = usarConsulta('get', '/api/v1/credits', { consulta: { limit: 200 } })
  const tiendas = usarConsulta('get', '/api/v1/stores', { consulta: { limit: 200 } })
  const categorias = usarConsulta('get', '/api/v1/categories', { consulta: { limit: 200 } })
  const proximo = usarConsulta('get', '/api/v1/dashboard/next-payment')
  const flujo = usarConsulta('get', '/api/v1/dashboard/fortnight-flow', { consulta: { period: quincenaDe(hoy) } })
  const plan = usarConsulta('get', '/api/v1/dashboard/month-vs-plan', { consulta: { period: mesDe(hoy) } })
  usarRecargaAlCambiar(plan.recargar, flujo.recargar)
  const valor = { hoy, alertas, personas, creditos, tiendas, categorias, proximo, flujo, plan } as Compartido
  return <Ctx.Provider value={valor}>{children}</Ctx.Provider>
}
export const listaDe = <T,>(c: EstadoConsulta<Lista<T>>): T[] => (Array.isArray(c.datos?.items) ? c.datos.items : [])
export const nombreCredito = (c: Compartido, id?: number) => listaDe(c.creditos).find((x) => x.id === id)?.name ?? t.comun.sinDato
export const karen = (c: Compartido) => listaDe(c.personas).find((p) => p.participation === 'external')
export const titulares = (c: Compartido) => listaDe(c.personas).filter((p) => p.participation !== 'external').sort((a, b) => ORDEN.indexOf(a.code) - ORDEN.indexOf(b.code))
const ORDEN = ['axel', 'ximena', 'karen', 'shared']
/** Días de aviso «Por vencer» (K44): solo de la lectura de alertas (C1-22); si falla, ninguno (U3 §1.2). */
export const diasAviso = (c: Compartido): number[] => { const v = c.alertas.datos?.parameters?.['alerts.payment_due_days_before']?.value; return Array.isArray(v) ? v : [] }
/** Color de tienda de C2 por nombre de la cadena. PROVISIONAL: C1 no trae un código de tienda (bloqueo WF3-inicio-tiendas). */
export function codigoTienda(s?: { name: string; kind: string }) {
  if (!s) return 'market'
  if (s.kind === 'market') return 'market'
  const n = s.name.normalize('NFD').replace(/[̀-ͯ']/g, '').toLowerCase().replace(/\s+/g, '-')
  return ['walmart', 'bodega-aurrera', 'sams-club'].find((c) => n.startsWith(c)) ?? 'market'
}

// ---- Estado de un pago (U3 §1.4, DESIGN §7.8): el primero que se cumpla ----
export interface EstadoPago { texto: string; icono: ReactNode; rol: RolInsignia; linea?: string }
type Parte = { person_id: number; amount_cents: number; paid_cents: number }
export function estadoPago(p: { status?: string; due_date: string; shares?: Parte[] }, hoy: string, aviso: number[], idKaren?: number): EstadoPago {
  const partes = p.shares ?? []
  const pagada = (s: Parte) => s.paid_cents >= s.amount_cents
  const k = partes.filter((s) => s.person_id === idKaren), h = partes.filter((s) => s.person_id !== idKaren)
  const falta = partes.some(pagada) && partes.some((s) => !pagada(s)) ? (k.length && k.every(pagada) ? t.pago.faltaHogar : h.length && h.every(pagada) ? t.pago.faltaKaren : undefined) : undefined
  if (p.status === 'paid') return { texto: t.pago.pagado, icono: <CircleCheck aria-hidden />, rol: 'exito' }
  if (p.status === 'overdue') return { texto: t.pago.vencido, icono: <CalendarX aria-hidden />, rol: 'advertencia', linea: [t.pago.sinMarcar, falta].filter(Boolean).join(' · ') }
  const d = difDias(hoy, p.due_date)
  if (aviso.includes(d)) return { texto: d === 0 ? t.pago.venceHoy : d === 1 ? t.pago.venceManana : t.pago.venceEn(d), icono: <CalendarClock aria-hidden />, rol: 'advertencia', linea: falta }
  if (falta) return { texto: falta, icono: <CircleDot aria-hidden />, rol: 'neutro' }
  return { texto: t.pago.pendiente, icono: <CircleDashed aria-hidden />, rol: 'neutro' }
}
export const InsigniaPago = ({ e }: { e: EstadoPago }) => <span className="flex flex-col gap-1"><Insignia rol={e.rol} icono={e.icono}>{e.texto}</Insignia>{e.linea && <span className="text-caption text-text-secondary">{e.linea}</span>}</span>

/** Título de panel de Detalle o bloque (font.size.title-3). */
export const TituloPanel = ({ children, id }: { children: ReactNode; id?: string }) => <h2 id={id} className="flex items-center gap-2 text-title-3 font-semibold text-text-default">{children}</h2>
export const Panel = ({ titulo, children, className, etiqueta }: { titulo: ReactNode; children: ReactNode; className?: string; etiqueta: string }) => (
  <section aria-label={etiqueta} className={cn('flex flex-col gap-3 rounded-md border border-border-divider bg-surface-card p-4 tablet:p-6', className)}><TituloPanel>{titulo}</TituloPanel>{children}</section>
)

// ---- Filtros (U3 §6): botón «Filtrar» que abre una hoja con chips; chips activos con «Quitar filtro» ----
export interface GrupoFiltro { id: string; nombre: string; multiple?: boolean; opciones: { valor: string; texto: string; marca?: ReactNode }[] }
export function Filtros({ grupos, valores, alCambiar }: { grupos: GrupoFiltro[]; valores: Record<string, string[]>; alCambiar: (id: string, v: string[]) => void }) {
  const [abierta, setAbierta] = useState(false)
  const activos = grupos.flatMap((g) => (valores[g.id] ?? []).map((v) => ({ g, v, texto: g.opciones.find((o) => o.valor === v)?.texto ?? v })))
  const alternar = (g: GrupoFiltro, v: string) => { const a = valores[g.id] ?? []; alCambiar(g.id, a.includes(v) ? a.filter((x) => x !== v) : g.multiple ? [...a, v] : [v]) }
  return (
    <div className="flex flex-wrap items-center gap-2">
      <Boton variante="secundario" onClick={() => setAbierta(true)}><ListFilter aria-hidden className="size-5" />{t.detalle.filtrar}</Boton>
      {activos.map(({ g, v, texto }) => (
        <ChipSeleccionable key={g.id + v} seleccionado aria-label={textos.vistas.quitarFiltro(`${g.nombre} ${texto}`)} onClick={() => alternar(g, v)}>{g.nombre}: {texto}<X aria-hidden className="size-5" /></ChipSeleccionable>
      ))}
      <Hoja abierta={abierta} alCambiar={setAbierta} titulo={t.detalle.filtros} pie={<Boton variante="secundario" tam="hoja" onClick={() => setAbierta(false)}>{t.detalle.listo}</Boton>}>
        <div className="flex flex-col gap-4">
          {grupos.map((g) => (
            <fieldset key={g.id} className="flex flex-col gap-2">
              <legend className="mb-2 text-small font-medium text-text-secondary">{g.nombre}</legend>
              <div className="flex flex-wrap gap-2">{g.opciones.map((o) => <ChipSeleccionable key={o.valor} seleccionado={(valores[g.id] ?? []).includes(o.valor)} marca={o.marca} onClick={() => alternar(g, o.valor)}>{o.texto}</ChipSeleccionable>)}</div>
            </fieldset>
          ))}
        </div>
      </Hoja>
    </div>
  )
}
/** Ancho de la ventana por puntos de quiebre de C2 (tablet 768, escritorio 1024, amplio 1440). */
export function usarMedia(consulta: string) {
  const [si, setSi] = useState(() => typeof matchMedia === 'function' && matchMedia(consulta).matches)
  useEffect(() => { const m = matchMedia(consulta); const f = () => setSi(m.matches); f(); m.addEventListener('change', f); return () => m.removeEventListener('change', f) }, [consulta])
  return si
}
export const TABLET = '(min-width: 48rem)'
export const AMPLIO = '(min-width: 90rem)'
