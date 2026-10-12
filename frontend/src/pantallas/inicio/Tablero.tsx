// Resumen de P-01 = tablero de control (U3 2.0 §2, K54): K-1 gasto del mes (indicador principal), B-1 barómetro, K-4 pagos
// próximos, K-2 ingreso, K-3 créditos, G-3, G-1, G-6, G-2, G-5 (= T-3), G-4 (= T-2) y «Más del hogar», en ese orden en el DOM
// (orden del foco de 360 en todos los anchos, §2.4). Cada bloque hace su propia petición y tiene sus estados (§1.3, §7); G-1 y
// G-6 comparten una. La interfaz no calcula porcentajes, diferencias ni estados: llegan de C1 v1.5 (motores S2 y S3).
import type { ReactNode } from 'react'
import { useState } from 'react'
import { ArrowDownLeft, ArrowUpRight, Banknote, CalendarClock, CircleCheck, CreditCard, Info, TriangleAlert, Wallet } from 'lucide-react'
import { Enlace } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { usarConsulta, type EstadoConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import type { Respuesta } from '../../api/tipos'
import { dinero, fecha, fechaConDia, mes, porcentaje } from '../../api/formatos'
import { Barometro } from '../../componentes/Barometro'
import { Boton } from '../../componentes/Boton'
import { patron } from '../../componentes/Grafica'
import { cn } from '../../componentes/utils'
import { textos } from '../../textos/es-MX'
import { CAJA, MasDelHogar, Tarjeta, Titulo } from './Resumen'
import { T2, T3 } from './Tendencias'
import { compacto, ContenedorTabla, ejes, FiguraGrafica, token } from './GraficaInicio'
import { EnlaceAccion, EsqueletoFilas, Estados, FallaBloque, IconoAviso, IconoInfo, Insignia, InsigniaPago, TABLET, diasAviso, estadoPago, karen, listaDe, nombreCredito, usarMedia, usarRecargaAlCambiar, useInicio } from './comun'
import { mesDe, sumarDias, sumarMeses } from './fechas'
import { MESES_LARGOS, textosInicio as t } from './textos'

const T = t.tablero, G = t.graficas, B = textos.barometro
const nombreMes = (m: string) => MESES_LARGOS[Number(m.slice(5, 7)) - 1]
const mesLargo = (m: string) => { const n = nombreMes(m); return `${n.charAt(0).toUpperCase()}${n.slice(1)} de ${m.slice(0, 4)}` }
const nada = t.comun.sinDato
const val = (v?: number) => (v === undefined ? nada : dinero(v))
const ESQUELETO = <div aria-hidden className="h-60 w-full rounded-md bg-skeleton-on-card tablet:h-72" />
const CAJA_GRAFICA = 'flex min-w-0 flex-col gap-3 rounded-md border border-border-divider bg-surface-card p-4 tablet:p-6'
const LIGA_FILA = 'flex min-h-11 flex-col justify-center gap-1 rounded-sm py-1 text-text-default no-underline hover:bg-surface-hover-subtle focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus'
type Mensual = Respuesta<'/api/v1/dashboard/monthly-overview', 'get'>

/** Rejilla de §2.4: 1 columna; 768 con 8 columnas (K-1 + B-1, K-4 entera, K-2 + K-3, una gráfica por fila); desde 1024 con 12
 *  (K-1 + B-1; K-2, K-3 y K-4 de 4; gráficas 6 + 6). K-4 va antes que K-2 en el DOM y se coloca a la derecha en 1024. */
const POS = {
  k1: 'tablet:col-span-4 desktop:col-span-6',
  b1: 'tablet:col-span-4 desktop:col-span-6',
  k4: 'tablet:col-span-8 desktop:col-span-4 desktop:col-start-9 desktop:row-start-2',
  k2: 'tablet:col-span-4 desktop:col-span-4 desktop:col-start-1 desktop:row-start-2',
  k3: 'tablet:col-span-4 desktop:col-span-4 desktop:col-start-5 desktop:row-start-2',
  g: 'tablet:col-span-8 desktop:col-span-6',
  mas: 'tablet:col-span-8 desktop:col-span-12',
}
const Bloque = ({ id, className, children }: { id: string; className: string; children: ReactNode }) => <div data-bloque={id} className={cn('min-w-0', className)}>{children}</div>

export default function Tablero() {
  const c = useInicio()
  const actual = mesDe(c.hoy)
  // G-1 y G-6: UNA petición de 17 meses (U3 §2.1; horizonte PROVISIONAL, bloqueo U3-tablero-horizonte).
  const mensual = usarConsulta('get', '/api/v1/dashboard/monthly-overview', { consulta: { from: sumarMeses(actual, -11), to: sumarMeses(actual, 5) } })
  const capacidad = usarConsulta('get', '/api/v1/dashboard/debt-capacity', { consulta: { period: actual } })
  usarRecargaAlCambiar(mensual.recargar)
  return (
    <section data-prueba="tablero" aria-label={T.region} className="grid grid-cols-1 gap-4 tablet:grid-cols-8 desktop:grid-cols-12 desktop:gap-6">
      <Bloque id="K-1" className={POS.k1}><K1 /></Bloque>
      <Bloque id="B-1" className={POS.b1}><B1 q={capacidad} /></Bloque>
      <Bloque id="K-4" className={POS.k4}><K4 /></Bloque>
      <Bloque id="K-2" className={POS.k2}><K2 /></Bloque>
      <Bloque id="K-3" className={POS.k3}><K3 /></Bloque>
      <Bloque id="G-3" className={POS.g}><G3 /></Bloque>
      <Bloque id="G-1" className={POS.g}><G1 q={mensual} /></Bloque>
      <Bloque id="G-6" className={POS.g}><G6 q={mensual} escala={capacidad.datos?.scale_max_pct} /></Bloque>
      <Bloque id="G-2" className={POS.g}><G2 /></Bloque>
      <Bloque id="G-5" className={POS.g}><T3 id="g5" sinFiltros /></Bloque>
      <Bloque id="G-4" className={POS.g}><T2 id="g4" /></Bloque>
      <Bloque id="mas" className={POS.mas}><MasDelHogar /></Bloque>
    </section>
  )
}

/** «[diferencia] más que [mes]», «menos que» o «Igual que», con icono de rol (DESIGN §3.3 y §7.6: la palabra va siempre). */
function Comparacion({ etiqueta, dif, mesAnterior }: { etiqueta?: string; dif: number; mesAnterior: string }) {
  const icono = dif > 0 ? <ArrowUpRight aria-hidden className="size-5 shrink-0 text-role-expense" /> : dif < 0 ? <ArrowDownLeft aria-hidden className="size-5 shrink-0 text-text-secondary" /> : null
  return <span className="inline-flex items-center gap-1">{icono}{etiqueta ? `${etiqueta} · ` : ''}{T.comparar(dinero(Math.abs(dif)), dif, mesAnterior)}</span>
}

// ---- K-1 · Gasto del mes (M-01 + 8 campos de C1 v1.5) ----
function K1() {
  const c = useInicio()
  return <Tarjeta n="K-1" titulo={T.k1.titulo} icono={<Wallet aria-hidden />} consulta={c.plan} falla={T.k1.falla} accionBase={{ a: rutaDe('P-06'), texto: T.k1.accion }} vista={(d) => {
    const conPlan = d.planned_cents !== undefined
    const usado = d.used_pct, resta = d.remaining_cents
    const rebasado = resta !== undefined && resta < 0
    const ant = d.previous_month ? nombreMes(d.previous_month) : undefined
    const alerta = d.status === 'warning' || d.status === 'exceeded'
    const sinPrevio = d.previous_actual_cents === undefined && d.previous_cash_cents === undefined
    const estado = d.status === 'warning' ? T.k1.vasEn(usado ?? 0) : d.status === 'exceeded' ? T.k1.rebasaste : undefined
    return {
      valor: dinero(d.actual_cents), etiquetaValor: T.k1.gastadoEn(nombreMes(d.month)),
      cuerpo: conPlan && usado !== undefined ? (
        <div role="meter" aria-label={T.k1.barra} aria-valuemin={0} aria-valuemax={Math.max(100, usado)} aria-valuenow={usado} aria-valuetext={T.k1.delPlan(porcentaje(usado))} className="h-[var(--size-meter-height)] w-full overflow-hidden rounded-full bg-viz-track">
          <div className={cn('h-full rounded-full', alerta ? 'bg-role-warning' : 'bg-viz-meter')} style={{ width: `${Math.min(100, Math.max(0, usado))}%` }} />
        </div>) : undefined,
      secundario: conPlan ? <>
        <span>{T.k1.dePlan(usado !== undefined ? porcentaje(usado) : nada, dinero(d.planned_cents!), rebasado ? nada : val(resta))}</span>
        {rebasado && <span className="inline-flex items-center gap-1 text-role-warning"><IconoAviso />{T.k1.rebasado(dinero(resta!))}</span>}
      </> : <><span>{`${nada} · ${T.k1.sinPlan(nombreMes(d.month))}`}</span><span>{T.k1.llevas(dinero(d.actual_cents))}</span></>,
      extra: (d.cash_cents !== undefined || ant) ? <>
        {d.cash_cents !== undefined && <span>{T.k1.efectivo(dinero(d.cash_cents))}</span>}
        {ant && sinPrevio && <span>{T.sinDatosDe(ant)}</span>}
        {ant && !sinPrevio && d.actual_difference_cents !== undefined && <Comparacion etiqueta={T.k1.gastado} dif={d.actual_difference_cents} mesAnterior={ant} />}
        {ant && !sinPrevio && d.cash_difference_cents !== undefined && <Comparacion etiqueta={T.k1.enEfectivo} dif={d.cash_difference_cents} mesAnterior={ant} />}
      </> : undefined,
      estado,
      insignias: (estado || (d.categories_exceeded_count ?? 0) > 0) ? <>
        {estado && <Insignia rol="advertencia" icono={<IconoAviso />}>{estado}</Insignia>}
        {(d.categories_exceeded_count ?? 0) > 0 && <Insignia rol="advertencia" icono={<IconoAviso />}>{T.k1.rebasadas(d.categories_exceeded_count!)}</Insignia>}
      </> : undefined,
      accion: conPlan ? { a: rutaDe('P-06'), texto: T.k1.accion } : { a: rutaDe('P-26', { seccion: '02' }), texto: T.k1.accionVacio },
    }
  }} />
}

// ---- K-2 · Ingreso del mes (N-100, M-18) ----
function K2() {
  const c = useInicio()
  const q = usarConsulta('get', '/api/v1/incomes/month', { consulta: { period: mesDe(c.hoy) } })
  return <Tarjeta n="K-2" titulo={T.k2.titulo} icono={<Banknote aria-hidden />} consulta={q} falla={T.k2.falla} accionBase={{ a: rutaDe('P-08'), texto: T.k2.accion }} vista={(d) => {
    const qs = d.fortnights ?? []
    if (!qs.length && d.extras_cents === 0) return { valor: nada, secundario: <span>{T.k2.faltaIngreso}</span>, accion: { a: rutaDe('P-26', { seccion: '03' }), texto: T.k2.accionVacio } }
    const esperado = qs.some((x) => x.status === 'expected')
    const ant = d.previous_month ? nombreMes(d.previous_month) : undefined
    return {
      // Quincenas del mes = capturadas + esperadas sin capturar (U3 K-2: suma de lo que está a la vista, del motor S2).
      valor: dinero(d.captured_cents + d.expected_pending_cents, 'ingreso'), rolValor: 'ingreso', iconoValor: <ArrowDownLeft aria-hidden />, etiquetaValor: T.k2.deQuincenas,
      secundario: <span>{d.extras_cents > 0 ? T.k2.extra(dinero(d.extras_cents, 'ingreso')) : T.k2.sinExtra}</span>,
      extra: <>
        <span>{T.k2.capturadas(qs.filter((x) => x.status === 'captured').length, qs.length)}</span>
        {ant && (d.previous_fortnights_cents === undefined ? <span>{T.sinDatosDe(ant)}</span> : d.fortnights_difference_cents !== undefined && <Comparacion dif={d.fortnights_difference_cents} mesAnterior={ant} />)}
      </>,
      estado: esperado ? T.k2.esperado : undefined,
      insignias: esperado ? <Insignia rol="info" icono={<IconoInfo />}>{T.k2.esperado}</Insignia> : undefined,
      accion: { a: rutaDe('P-08'), texto: T.k2.accion },
    }
  }} />
}

// ---- K-3 · Créditos (M-12 a hoy, M-19 y disponible M-04 con la lectura activa) ----
function K3() {
  const c = useInicio()
  const q = usarConsulta('get', '/api/v1/credits/totals')
  return <Tarjeta n="K-3" titulo={T.k3.titulo} icono={<CreditCard aria-hidden />} consulta={q} falla={T.k3.falla} accionBase={{ a: rutaDe('P-17'), texto: T.k3.accion }} vista={(d) => {
    if (c.creditos.datos && !listaDe(c.creditos).length) return { valor: T.k3.sinCreditos, accion: { a: rutaDe('P-26', { seccion: '05' }), texto: T.k3.accionVacio } }
    const m = nombreMes(mesDe(c.hoy))
    const sinPagos = d.month_total_cents === 0 && !d.month_payment_count
    const av = c.proximo.datos?.availability ?? [], linea = av[0]
    const disponible = c.proximo.error ? <span className="border-t border-border-divider pt-2">{T.k3.disponibleSinDato}</span> : linea && (
      <span className="flex flex-col gap-1 border-t border-border-divider pt-2">
        <span className="tabular-nums text-text-default">{t.t3.disponibleHoy(dinero(linea.available_cents), linea.availability_reading)}</span>
        {!linea.availability_confirmed && <Insignia rol="info" icono={<IconoInfo />} etiqueta={t.t3.reposicion(nombreCredito(c, linea.credit_id))}>{t.t3.pendienteConfirmar}</Insignia>}
        {av.length > 1 && <span>{t.t3.masLineas(av.length - 1)}</span>}
      </span>
    )
    return {
      valor: dinero(d.live_debt_cents), rolValor: 'deuda', iconoValor: <CreditCard aria-hidden />, etiquetaValor: T.k3.deudaViva,
      secundario: <>{sinPagos ? <span>{T.k3.sinPagos(m)}</span> : <span>{T.k3.pagosDe(m, dinero(d.month_total_cents), d.month_payment_count)}</span>}{disponible}</>,
      extra: sinPagos ? undefined : <>
        <span>{t.comun.hogarYKaren(dinero(d.month_household_cents), d.month_external_cents > 0 ? dinero(d.month_external_cents) : undefined)}</span>
        {d.month_pending_household_cents !== undefined && <span>{T.k3.faltaPagar(dinero(d.month_pending_household_cents))}</span>}
      </>,
      accion: { a: rutaDe('P-17'), texto: T.k3.accion },
    }
  }} />
}

// ---- K-4 · Pagos próximos (M-13): vencidos (si los hay) y de hoy a hoy + 30, N de getDashboardSettings (PROVISIONAL) ----
function K4() {
  const c = useInicio()
  const n = c.ajustesTablero.datos?.parameters?.['dashboard.payment_reminder_count']?.value
  const N = typeof n === 'number' && n > 0 ? n : undefined
  const np = c.proximo.datos?.next_payment
  const conVencidos = N !== undefined && np?.is_overdue === true
  const vencidos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { status: 'overdue', sort: 'due_date', limit: N } }, conVencidos)
  const proximos = usarConsulta('get', '/api/v1/scheduled-payments', { consulta: { status: 'pending', sort: 'due_date', limit: N, period: `${c.hoy}..${sumarDias(c.hoy, 30)}` } }, N !== undefined)
  const idK = karen(c)?.id, aviso = diasAviso(c)
  const sinN = c.ajustesTablero.error || c.ajustesTablero.nuncaCargada || (!c.ajustesTablero.cargando && N === undefined)
  return (
    <section data-prueba="tarjeta-K-4" aria-label={T.k4.titulo} className={CAJA}>
      <Titulo icono={<CalendarClock aria-hidden />} titulo={T.k4.titulo} />
      {sinN ? <FallaBloque texto={T.k4.falla} alReintentar={c.ajustesTablero.recargar} />
        : N === undefined ? <EsqueletoFilas n={3} />
        : <Estados consulta={proximos} falla={T.k4.falla} esqueleto={<EsqueletoFilas n={3} />}>{(d) => {
          const filas = [...(conVencidos ? vencidos.datos?.items ?? [] : []), ...d.items].slice(0, N)
          if (!filas.length) {
            if (c.creditos.datos && !listaDe(c.creditos).length) return <><p className="text-small text-text-secondary">{T.k4.sinCreditos}</p><EnlaceAccion a={rutaDe('P-26', { seccion: '05' })} texto={T.k4.accionVacio} /></>
            return <p className="text-small text-text-secondary">{T.k4.sinPagos}</p>
          }
          return <>
            <ul className="flex flex-col divide-y divide-border-divider">
              {filas.map((p) => {
                const partes = p.shares ?? [], suma = (xs: typeof partes) => xs.reduce((a, s) => a + s.amount_cents, 0)
                const k = suma(partes.filter((s) => s.person_id === idK)), h = suma(partes.filter((s) => s.person_id !== idK))
                return <li key={p.id}><Enlace a={p.credit_id !== undefined ? rutaDe('P-18', { id: p.credit_id }) : rutaDe('P-17')} className={LIGA_FILA}>
                  <span className="flex flex-wrap items-center gap-x-1 gap-y-1 text-small tabular-nums">{T.k4.fila(fechaConDia(p.due_date), nombreCredito(c, p.credit_id), dinero(p.amount_cents))}{' · '}<InsigniaPago e={estadoPago(p, c.hoy, aviso, idK)} /></span>
                  {partes.length > 0 && <span className="hidden text-small text-text-secondary tabular-nums phone:block">{t.comun.hogarYKaren(dinero(h), k > 0 ? dinero(k) : undefined)}</span>}
                </Enlace></li>
              })}
            </ul>
            {d.next_cursor && <p className="text-small text-text-secondary">{T.k4.hayMas}</p>}
          </>
        }}</Estados>}
      <EnlaceAccion a={`${rutaDe('P-01')}?sub=calendario`} texto={T.k4.accion} />
    </section>
  )
}

// ---- B-1 · Barómetro de capacidad de endeudamiento (M-20, getDashboardDebtCapacity) ----
function B1({ q }: { q: EstadoConsulta<Respuesta<'/api/v1/dashboard/debt-capacity', 'get'>> }) {
  if (esServidorViejo(q.error)) return (
    <section data-prueba="tarjeta-B-1" aria-label={T.b1.titulo} className={CAJA}>
      <Titulo icono={<CreditCard aria-hidden />} titulo={T.b1.titulo} />
      <ServidorViejo />
      <EnlaceAccion a={rutaDe('P-17')} texto={T.b1.accion} />
    </section>
  )
  return <Tarjeta n="B-1" titulo={T.b1.titulo} icono={<CreditCard aria-hidden />} consulta={q} falla={T.b1.falla} accionBase={{ a: rutaDe('P-17'), texto: T.b1.accion }} vista={(d) => {
    const sinIngreso = d.income_cents === undefined
    const conValor = d.ratio_pct !== undefined && d.status !== undefined
    const texto = conValor ? T.b1.valorTexto(`${d.ratio_pct} %`, d.limit_pct, B.colores[d.status!]) : T.b1.sinValorTexto(d.limit_pct)
    return {
      cuerpo: <Barometro titulo={T.b1.titulo} valor={conValor ? d.ratio_pct : undefined} estado={conValor ? d.status : undefined} aviso={d.warning_pct} tope={d.limit_pct} escala={d.scale_max_pct} textoAccesible={texto} />,
      secundario: sinIngreso ? <span>{T.b1.faltaIngreso}</span> : <>
        <span>{T.b1.pagosDe(dinero(d.payments_cents), dinero(d.income_cents!))}</span>
        {!d.has_credits && <span>{T.b1.sinCreditos}</span>}
      </>,
      extra: sinIngreso ? undefined : <>
        {d.headroom_cents !== undefined && (d.headroom_cents >= 0 ? <span>{T.b1.teCaben(dinero(d.headroom_cents))}</span>
          : <span className="inline-flex items-center gap-1 text-role-warning"><IconoAviso />{T.b1.arriba(dinero(d.headroom_cents))}</span>)}
        {d.ratio_total_pct !== undefined && <span className="hidden tablet:inline">{T.b1.conKaren(d.ratio_total_pct)}</span>}
      </>,
      insignias: d.income_is_expected && !sinIngreso ? <Insignia rol="info" icono={<IconoInfo />}>{T.b1.esperado}</Insignia> : undefined,
      estado: conValor ? `${d.ratio_pct} % · ${B.estados[d.status!]}` : nada,
      accion: sinIngreso ? { a: rutaDe('P-26', { seccion: '03' }), texto: T.b1.accionVacio } : { a: rutaDe('P-17'), texto: T.b1.accion },
    }
  }} />
}

// ---- Piezas de las gráficas del tablero ----
/** Servidor anterior a C1 v1.5 (U3 2.0 §2, reglas del tablero): la operación nueva responde 501 `not_implemented` o 404
 *  «La ruta no existe» (C1 no declara 404 en getDashboardDebtCapacity ni en getDashboardMonthlyOverview, así que un 404 ahí
 *  solo es ruta inexistente). Solo se consulta en los bloques que dependen solo de una operación nueva: B-1, G-1 y G-6. */
const esServidorViejo = (e: unknown) => e instanceof ErrorApi && (e.estado === 501 || e.estado === 404)
/** «Disponible cuando el servidor se actualice» en `font.size.caption` y rol información (par C-073 de contraste-componentes):
 *  sin cifra y sin «Reintentar», porque reintentar no lo arregla. */
const ServidorViejo = () => (
  <p role="status" className="flex w-fit max-w-full items-center gap-1 rounded-sm bg-role-info-surface px-2 py-1 text-caption text-role-info [&_svg]:size-5 [&_svg]:shrink-0">
    <Info aria-hidden />{T.servidorViejo}
  </p>
)
/** Estados de una gráfica con su título visible mientras carga o falla (§7: «ídem T-1»). `soloNueva`: el bloque depende solo
 *  de una operación de C1 v1.5 y, con un servidor anterior, muestra ServidorViejo en lugar de la falla. */
function EstadosGrafica<D>({ titulo, consulta, soloNueva = false, children }: { titulo: string; consulta: EstadoConsulta<D>; soloNueva?: boolean; children: (d: D) => ReactNode }) {
  if (soloNueva && esServidorViejo(consulta.error)) return <section aria-label={titulo} className={CAJA_GRAFICA}><span className="text-small font-medium text-text-secondary">{titulo}</span><ServidorViejo /></section>
  if (consulta.datos !== undefined && !consulta.error && !consulta.nuncaCargada) return <Estados consulta={consulta} falla={G.falla} esqueleto={ESQUELETO}>{children}</Estados>
  return <section aria-label={titulo} className={CAJA_GRAFICA}><span className="text-small font-medium text-text-secondary">{titulo}</span><Estados consulta={consulta} falla={G.falla} esqueleto={ESQUELETO}>{children}</Estados></section>
}
function VacioBloque({ id, titulo, texto, accion }: { id: string; titulo: string; texto: ReactNode; accion?: ReactNode }) {
  return <section data-prueba={`grafica-${id}`} aria-label={titulo} className={CAJA_GRAFICA}><span className="text-small font-medium text-text-secondary">{titulo}</span><p className="flex items-center gap-2 text-body text-text-secondary">{texto}</p>{accion}</section>
}
/** Figura con barras horizontales en HTML (G-2, G-3): mismas reglas que FiguraGrafica (figcaption, resumen, «Ver como tabla»
 *  con las mismas cifras, tabla siempre disponible para lectores de pantalla) y filas que pueden ser enlaces de 44 px. */
function FiguraBarras({ id, titulo, subtitulo, resumen, tabla, notas, children }: { id: string; titulo: string; subtitulo: string; resumen: string; tabla: { columnas: string[]; filas: string[][] }; notas?: ReactNode; children: ReactNode }) {
  const [comoTabla, setComoTabla] = useState(false)
  return (
    <figure data-prueba={`grafica-${id}`} className={CAJA_GRAFICA}>
      <figcaption className="flex flex-col gap-1">
        <span className="text-small font-medium text-text-secondary">{titulo}</span>
        <span className="text-caption text-text-secondary">{subtitulo}</span>
        <span className="text-small text-text-default tabular-nums">{resumen}</span>
      </figcaption>
      {!comoTabla && children}
      <ContenedorTabla comoTabla={comoTabla} titulo={titulo}>
        <table className="w-full text-small tabular-nums">
          <caption className="sr-only">{titulo}</caption>
          <thead><tr>{tabla.columnas.map((c, i) => <th key={c} scope="col" className={cn('whitespace-nowrap px-2 py-1 font-medium text-text-secondary', i ? 'text-right' : 'text-left')}>{c}</th>)}</tr></thead>
          <tbody>{tabla.filas.map((f) => <tr key={f[0]} className="border-t border-border-divider">{f.map((v, i) => (i ? <td key={i} className="whitespace-nowrap px-2 py-1 text-right">{v}</td> : <th key={i} scope="row" className="whitespace-nowrap px-2 py-1 text-left font-normal">{v}</th>))}</tr>)}</tbody>
        </table>
      </ContenedorTabla>
      {notas}
      <Boton variante="enlace" className="w-max" aria-pressed={comoTabla} onClick={() => setComoTabla((x) => !x)}>{comoTabla ? t.acciones.verComoGrafica : t.acciones.verComoTabla}</Boton>
    </figure>
  )
}
const tooltipDinero = { trigger: 'axis', valueFormatter: (v: unknown) => (typeof v === 'number' ? dinero(v) : nada) }

// ---- G-1 · Ingreso, gasto y pagos por mes (M-21): columnas agrupadas, 3 series en orden fijo con patrón ----
function G1({ q }: { q: EstadoConsulta<Mensual> }) {
  const c = useInicio(), ancho = usarMedia(TABLET), actual = mesDe(c.hoy)
  return <EstadosGrafica titulo={T.g1.titulo} consulta={q} soloNueva>{(d) => {
    const m = d.months.filter((x) => !x.is_projection && x.month <= actual).slice(-12)
    if (!m.length || m.every((x) => x.income_total_cents === undefined && x.accrued_cents === undefined))
      return <VacioBloque id="g1" titulo={T.g1.titulo} texto={T.g1.vacio} accion={<EnlaceAccion a={rutaDe('P-26', { seccion: '03' })} texto={T.g1.accionVacio} />} />
    const iu = Math.max(0, m.findIndex((x) => x.month === actual)), u = m[iu] ?? m[m.length - 1]
    const cs = m.map((x) => mes(x.month)), sup = token('--ui-surface-card')
    const serie = (nombre: string, i: number, v: (number | undefined)[]) => ({
      name: nombre, type: 'bar', barGap: '8%', itemStyle: { color: token(`--ui-chart-${i + 1}`), decal: patron(i, sup) }, emphasis: { disabled: true },
      data: v.map((x, j) => (x === undefined ? null : j === iu ? { value: x, label: { show: true, position: 'top', color: token('--ui-text-default'), fontSize: 12, formatter: () => dinero(x) } } : x)),
    })
    return <FiguraGrafica id="g1" titulo={T.g1.titulo} subtitulo={T.g1.subtitulo(mes(m[0].month), mes(m[m.length - 1].month))} dependencias={[d, ancho]}
      resumen={T.g1.resumen(mes(u.month), val(u.income_total_cents), val(u.accrued_cents), dinero(u.credit_payments_cents))} notas={[T.g1.nota]}
      tabla={{ columnas: [G.mes, T.g1.ingreso, T.g1.quincenas, T.g1.extras, T.g1.gasto, T.g1.pagos], filas: m.map((x) => [mes(x.month), val(x.income_total_cents), val(x.income_fortnights_cents), val(x.income_extras_cents), val(x.accrued_cents), dinero(x.credit_payments_cents)]) }}
      opcion={() => ({ ...ejes(cs, compacto, ancho), tooltip: tooltipDinero, series: [serie(T.g1.ingreso, 0, m.map((x) => x.income_total_cents)), serie(T.g1.gasto, 1, m.map((x) => x.accrued_cents)), serie(T.g1.pagos, 2, m.map((x) => x.credit_payments_cents))] })} />
  }}</EstadosGrafica>
}

// ---- G-6 · Tendencia del barómetro (M-21): línea con proyección punteada, referencias en los umbrales y «Hoy» ----
function G6({ q, escala }: { q: EstadoConsulta<Mensual>; escala?: number }) {
  const c = useInicio(), ancho = usarMedia(TABLET), actual = mesDe(c.hoy)
  return <EstadosGrafica titulo={T.g6.titulo} consulta={q} soloNueva>{(d) => {
    const m = d.months.filter((x) => x.month >= sumarMeses(actual, -6) && x.month <= sumarMeses(actual, 5))
    if (!m.length || m.every((x) => x.debt_ratio_pct === undefined))
      return <VacioBloque id="g6" titulo={T.g6.titulo} texto={T.g6.vacio} accion={<EnlaceAccion a={rutaDe('P-26', { seccion: '03' })} texto={T.g6.accionVacio} />} />
    const cs = m.map((x) => mes(x.month)), hoyI = m.findIndex((x) => x.month === actual)
    const color = (e?: 'normal' | 'warning' | 'exceeded') => (e ? B.colores[e] : nada)
    const pct = (v?: number) => (v === undefined ? nada : `${v} %`)
    const real = m.map((x) => (!x.is_projection ? x.debt_ratio_pct ?? null : null))
    const proy = m.map((x, i) => (x.is_projection || i === hoyI ? x.debt_ratio_pct ?? null : null))
    const proyectados = m.filter((x) => x.is_projection && x.debt_ratio_pct !== undefined)
    const mayor = proyectados.reduce<typeof proyectados[number] | undefined>((a, b) => (!a || b.debt_ratio_pct! > a.debt_ratio_pct! ? b : a), undefined)
    const u = m[hoyI] ?? m[m.length - 1]
    const tope = Math.max(escala ?? d.limit_pct, ...m.map((x) => x.debt_ratio_pct ?? 0))
    const c1 = token('--ui-chart-1'), base = token('--ui-viz-baseline'), sec = token('--ui-text-secondary')
    const referencia = (v: number) => ({ yAxis: v, lineStyle: { color: base, width: 1, type: 'solid' }, label: { formatter: `${v} %`, position: 'end', color: sec } })
    return <FiguraGrafica id="g6" titulo={T.g6.titulo} subtitulo={T.g1.subtitulo(cs[0], cs[cs.length - 1])} dependencias={[d, ancho, escala]}
      resumen={mayor ? T.g6.resumen(mes(u.month), pct(u.debt_ratio_pct), color(u.debt_status), mes(mayor.month), pct(mayor.debt_ratio_pct)) : T.g6.resumenSinProyeccion(mes(u.month), pct(u.debt_ratio_pct), color(u.debt_status))}
      notas={[T.g6.nota]}
      tabla={{ columnas: [G.mes, T.g6.capacidad, T.g6.estado, T.g1.pagos, T.g1.ingreso, T.g6.proyeccion], filas: m.map((x) => [mes(x.month), x.debt_ratio_pct === undefined ? T.g6.faltaIngreso : pct(x.debt_ratio_pct), color(x.debt_status), dinero(x.credit_payments_cents), val(x.income_total_cents), x.is_projection ? T.g6.si : T.g6.no]) }}
      opcion={() => {
        const e = ejes(cs, (v: number) => `${v} %`, ancho)
        return { ...e, yAxis: { ...e.yAxis, min: 0, max: tope },
          tooltip: { trigger: 'axis', formatter: (ps: { dataIndex: number }[]) => { const x = m[ps[0]?.dataIndex ?? 0]; return x.debt_ratio_pct === undefined ? `${mesLargo(x.month)}: ${T.g6.faltaIngreso}` : T.g6.punto(mesLargo(x.month), x.debt_ratio_pct, color(x.debt_status)) } },
          series: [
            { name: T.g6.serie, type: 'line', symbol: 'circle', symbolSize: 8, connectNulls: false, emphasis: { disabled: true }, lineStyle: { color: c1, width: 2 }, itemStyle: { color: c1 },
              data: real.map((v, i) => (v !== null && i === hoyI ? { value: v, label: { show: true, position: 'top', color: token('--ui-text-default'), fontSize: 12, formatter: `${v} %` } } : v)),
              markLine: { symbol: 'none', silent: true, data: [referencia(d.warning_pct), referencia(d.limit_pct), ...(hoyI >= 0 ? [{ xAxis: cs[hoyI], lineStyle: { color: base, width: 1, type: 'solid' }, label: { formatter: G.hoy, color: sec } }] : [])] } },
            { name: T.g6.serie, type: 'line', symbol: 'circle', symbolSize: 8, connectNulls: false, emphasis: { disabled: true }, lineStyle: { color: c1, width: 2, type: [8, 4] }, itemStyle: { color: c1 }, data: proy },
          ] }
      }} />
  }}</EstadosGrafica>
}

// ---- G-2 · Deuda viva por línea (M-12): barras horizontales en orden descendente con el valor al final ----
function G2() {
  const c = useInicio()
  const q = usarConsulta('get', '/api/v1/debt/live', { consulta: { date: c.hoy } })
  return <EstadosGrafica titulo={T.g2.titulo} consulta={q}>{(d) => {
    const filas = d.credits.filter((x) => x.total_cents > 0).sort((a, b) => b.total_cents - a.total_cents)
    if (!filas.length || d.total_cents === 0) return <VacioBloque id="g2" titulo={T.g2.titulo} texto={<><CircleCheck aria-hidden className="size-5 text-role-success" />{T.g2.vacio}</>} />
    const max = filas[0].total_cents
    return <FiguraBarras id="g2" titulo={T.g2.titulo} subtitulo={T.g2.subtitulo(fecha(d.date))} resumen={T.g2.resumen(dinero(d.total_cents), nombreCredito(c, filas[0].credit_id), dinero(filas[0].total_cents))}
      tabla={{ columnas: [T.g2.linea, T.g2.monto], filas: filas.map((x) => [nombreCredito(c, x.credit_id), dinero(x.total_cents)]) }}>
      <ul aria-hidden className="flex flex-col gap-3">
        {filas.map((x) => <li key={x.credit_id} className="flex flex-col gap-1">
          <span className="text-small text-text-default">{nombreCredito(c, x.credit_id)}</span>
          <span className="flex items-center gap-2"><span className="block h-[var(--size-viz-bar-max)] shrink-0 rounded-r-[var(--size-viz-bar-radius)] bg-chart-1" style={{ width: `calc((100% - 7rem) * ${x.total_cents / max})` }} /><span className="whitespace-nowrap text-small text-text-default tabular-nums">{dinero(x.total_cents)}</span></span>
        </li>)}
      </ul>
    </FiguraBarras>
  }}</EstadosGrafica>
}

// ---- G-3 · Presupuesto por categoría (M-09) con la sugerencia subir o recortar (M-22, PROVISIONAL) ----
function G3() {
  const c = useInicio(), amplio = usarMedia(TABLET), periodo = mesDe(c.hoy)
  const q = usarConsulta('get', '/api/v1/budget/categories', { consulta: { period: periodo } })
  const senales = usarConsulta('get', '/api/v1/budget/signals', { consulta: { period: periodo } })
  const umbral = (k: 'alerts.budget_warning_pct' | 'alerts.budget_limit_pct') => { const v = c.alertas.datos?.parameters?.[k]?.value; return typeof v === 'number' ? v : undefined }
  const w = umbral('alerts.budget_warning_pct'), l = umbral('alerts.budget_limit_pct')
  return <EstadosGrafica titulo={T.g3.titulo} consulta={q}>{(d) => {
    const cats = listaDe(c.categorias), nombre = (id: number) => cats.find((x) => x.id === id)?.name ?? nada
    const conPlan = d.rows.filter((r) => r.planned_cents !== undefined && r.used_pct !== undefined).sort((a, b) => b.used_pct! - a.used_pct!)
    const sinPlan = d.rows.filter((r) => r.planned_cents === undefined && r.actual_cents > 0).length
    const capturar = <EnlaceAccion a={rutaDe('P-26', { seccion: '02' })} texto={T.g3.accionVacio} />
    if (!conPlan.length) return <VacioBloque id="g3" titulo={T.g3.titulo} texto={T.g3.vacio} accion={capturar} />
    const visibles = amplio ? conPlan : conPlan.slice(0, 8)
    const eje = Math.max(100, ...visibles.map((r) => r.used_pct!))
    const pos = (n: number) => `${(Math.min(n, eje) / eje) * 100}%`
    const s = senales.error ? undefined : senales.datos
    const senal = (id: number) => s?.rows.find((r) => r.category_id === id)?.signal
    const alerta = (r: { status?: string }) => r.status === 'warning' || r.status === 'exceeded'
    const textoSenal = (x?: string) => (x === 'raise' ? T.g3.subir : x === 'cut' ? T.g3.recortar : T.g3.ninguna)
    return <FiguraBarras id="g3" titulo={T.g3.titulo} subtitulo={T.g3.subtitulo(mes(d.month ?? periodo))}
      resumen={T.g3.resumen(nombreMes(periodo), visibles.filter(alerta).length, w ?? 0, visibles.filter((r) => r.status === 'exceeded').length)}
      tabla={{ columnas: [T.g3.categoria, T.g3.real, T.g3.plan, T.g3.usado, T.g3.sugerencia], filas: visibles.map((r) => [nombre(r.category_id), dinero(r.actual_cents), dinero(r.planned_cents!), porcentaje(r.used_pct!), textoSenal(senal(r.category_id))]) }}
      notas={<>
        {!amplio && conPlan.length > visibles.length && <EnlaceAccion a={rutaDe('P-06')} texto={T.g3.verTodas} />}
        {sinPlan > 0 && <p className="flex flex-wrap items-center gap-2 text-small text-text-secondary">{T.g3.sinPlan(sinPlan)}{capturar}</p>}
      </>}>
      <div aria-hidden className="relative h-5 text-caption text-text-secondary tabular-nums">
        {w !== undefined && <span className="absolute -translate-x-1/2" style={{ left: pos(w) }}>{T.g3.marca(w)}</span>}
        {l !== undefined && <span className="absolute -translate-x-1/2" style={{ left: pos(l) }}>{T.g3.marca(l)}</span>}
      </div>
      <ul className="flex flex-col gap-1">
        {visibles.map((r) => {
          const x = senal(r.category_id), n = nombre(r.category_id)
          return <li key={r.category_id}><Enlace a={`${rutaDe('P-02')}?categoria=${r.category_id}&periodo=${periodo}`} aria-label={T.g3.abrir(n)} className={LIGA_FILA}>
            <span className="flex flex-wrap items-center gap-x-2 text-small tabular-nums">{T.g3.etiqueta(n, dinero(r.actual_cents), dinero(r.planned_cents!), porcentaje(r.used_pct!))}
              {alerta(r) && <span className="inline-flex items-center gap-1 text-role-warning"><TriangleAlert aria-hidden className="size-5" />{r.status === 'exceeded' ? T.g3.rebasada : T.g3.vasEn(r.used_pct!)}</span>}
            </span>
            <span className="relative block h-[var(--size-viz-bar-max)] w-full">
              <span className={cn('absolute inset-y-0 left-0 rounded-r-[var(--size-viz-bar-radius)]', alerta(r) ? 'bg-role-warning' : 'bg-chart-1')} style={{ width: pos(r.used_pct!) }} />
              {w !== undefined && <span aria-hidden data-marca="aviso" className="absolute -inset-y-1 border-l border-dashed border-viz-baseline" style={{ left: pos(w) }} />}
              {l !== undefined && <span aria-hidden data-marca="tope" className="absolute -inset-y-1 border-l-2 border-viz-baseline" style={{ left: pos(l) }} />}
            </span>
            {(x === 'raise' || x === 'cut') && s && <Insignia rol="info" icono={<IconoInfo />} etiqueta={T.g3.motivo(n, s.months_considered, x === 'raise', s.margin_pct)}>{textoSenal(x)}</Insignia>}
          </Enlace></li>
        })}
      </ul>
    </FiguraBarras>
  }}</EstadosGrafica>
}
