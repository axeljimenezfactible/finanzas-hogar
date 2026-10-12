// P-07 Flujo y proyección (U1 §2.1, §2.2, §3.9: N-090 a N-093; interruptor A/B de los meses 7 a 12, N-318, K21). Resumen* · Detalle · Tendencias.
import { useState } from 'react'
import { Calculator } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { conConsulta, navegar } from '../../app/enrutador'
import { pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { dinero, fecha, hoyISO, mes } from '../../api/formatos'
import { useAvisos } from '../../componentes/Avisos'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Grafica } from '../../componentes/Grafica'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { cn } from '../../componentes/utils'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { Bloque, Fila, mensajeDe, rejillaTarjetas, useCatalogos, useSinRed } from './comun'
import { t } from './textos'

export default function P07Flujo({ sub, consulta }: PropsPantalla) {
  const ajustes = usarConsulta('get', '/api/v1/settings/projection')
  const flujo = usarConsulta('get', '/api/v1/cash-flow/fortnights')
  const meses = usarConsulta('get', '/api/v1/cash-flow/accrual-vs-cash')
  const avisar = useAvisos()
  const sinRed = useSinRed()
  const [cambiando, setCambiando] = useState(false)
  const encendido = ajustes.datos?.repeat_cadence_enabled ?? false
  const alternar = async () => {
    setCambiando(true)
    try { await pedir('patch', '/api/v1/settings/projection', { cuerpo: { repeat_cadence_enabled: !encendido } }); ajustes.recargar(); flujo.recargar(); meses.recargar() }
    catch (e) { avisar({ texto: mensajeDe(e), rol: 'error' }) } finally { setCambiando(false) }
  }
  const modo = flujo.datos?.projection_mode ?? meses.datos?.projection_mode
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-3">
        <button type="button" role="switch" aria-checked={encendido} aria-disabled={!!sinRed || cambiando || !ajustes.datos} title={sinRed} onClick={() => !sinRed && !cambiando && ajustes.datos && alternar()}
          className="flex min-h-11 items-center gap-3 rounded-md px-2 text-left text-small text-text-default hover:bg-surface-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus aria-disabled:text-text-disabled">
          <span aria-hidden className={cn('flex h-6 w-11 shrink-0 items-center rounded-full p-0.5 transition-colors', encendido ? 'bg-action-primary-bg' : 'bg-viz-track')}>
            <span className={cn('size-5 rounded-full bg-surface-card shadow-1 transition-transform', encendido && 'translate-x-5')} />
          </span>{t.p07.interruptor}
        </button>
        {modo && <span className="text-caption text-text-secondary">{t.p07.modo(modo === 'with_repeat_cadence' ? t.p07.modoB : t.p07.modoA)}</span>}
        <Boton variante="primario" className="ml-auto" onClick={() => navegar(rutaDe('P-20'))}><Calculator aria-hidden className="size-5" />{t.p07.simular}</Boton>
      </div>
      <ConDatos consulta={flujo}>{(f) => {
        const filas = f.rows as Required<(typeof f.rows)[number]>[]
        if (!filas.length) return <Vacio texto={t.p07.vacio} accion={<Boton onClick={() => navegar(rutaDe('P-26', { seccion: '03' }))}>{t.p07.irAjustes}</Boton>} />
        if (sub === 'detalle') return <Detalle filas={filas} meses={meses} />
        if (sub === 'tendencias') return <Tendencias meses={meses} />
        return <Resumen filas={filas} meses={meses} consulta={consulta} />
      }}</ConDatos>
    </div>
  )
}

type Quincena = Required<{ start_date: string; end_date: string; pay_date: string; income_cents: number; income_is_expected: boolean; credit_payments_cents: number; fixed_and_upcoming_cents: number; variable_estimate_cents: number; planned_savings_cents?: number; balance_cents: number; cumulative_balance_cents: number; is_actual: boolean }>
type Meses = ReturnType<typeof usarConsulta<'/api/v1/cash-flow/accrual-vs-cash'>>

function Resumen({ filas, meses, consulta }: { filas: Quincena[]; meses: Meses; consulta: URLSearchParams }) {
  const hoy = hoyISO()
  const fechaDeuda = consulta.get('deuda') ?? filas[filas.length - 1].end_date
  const deuda = usarConsulta('get', '/api/v1/debt/live', { consulta: { date: fechaDeuda } })
  const supuestos = usarConsulta('get', '/api/v1/cash-flow/assumptions')
  const cats = useCatalogos()
  const [verSupuestos, setVerSupuestos] = useState(false)
  const actual = filas.find((q) => q.start_date <= hoy && hoy <= q.end_date) ?? filas[0]
  const ultimoMes = meses.datos?.months[meses.datos.months.length - 1]
  const detalle = () => navegar(conConsulta({ sub: 'detalle' }))
  return (
    <div className="flex flex-col gap-6">
      <div className={rejillaTarjetas}>
        <TarjetaMetrica titulo={t.p07.saldoQuincena} valor={dinero(actual.balance_cents)} nota={`${t.p07.pago}: ${fecha(actual.pay_date)}`} alAbrir={detalle} />
        <TarjetaMetrica titulo={t.p07.saldoAcumulado} valor={dinero(filas[filas.length - 1].cumulative_balance_cents)} nota={fecha(filas[filas.length - 1].end_date)} alAbrir={detalle} />
        <TarjetaMetrica titulo={t.p07.diferenciaAcumulada} valor={ultimoMes ? dinero(ultimoMes.cumulative_difference_cents ?? 0) : '—'} nota={ultimoMes && mes(ultimoMes.month)} alAbrir={() => navegar(conConsulta({ sub: 'tendencias' }))} />
        <TarjetaMetrica titulo={t.p07.deudaViva} rol="debt" valor={deuda.datos ? dinero(deuda.datos.total_cents) : '—'} nota={t.p07.deudaNota(fecha(fechaDeuda))}>
          <Campo etiqueta={t.p07.fechaDeuda} type="date" value={fechaDeuda} onChange={(e) => e.target.value && navegar(conConsulta({ deuda: e.target.value }), { reemplazar: true })} />
        </TarjetaMetrica>
      </div>
      <Bloque titulo={t.p07.supuestos} accion={<Boton variante="fantasma" aria-expanded={verSupuestos} onClick={() => setVerSupuestos(!verSupuestos)}>{t.p07.verSupuestos}</Boton>}>
        {verSupuestos && <ConDatos consulta={supuestos}>{(s) => (
          <div>
            <Fila etiqueta={t.p07.metodo}>{s.variable_estimate_method}</Fila>
            <Fila etiqueta={t.p07.variable}>{s.variable_estimate_cents == null ? '—' : dinero(s.variable_estimate_cents)}</Fila>
            <Fila etiqueta={t.p07.interruptor}>{s.projection_mode === 'with_repeat_cadence' ? t.p07.modoB : t.p07.modoA}</Fila>
            <h3 className="pt-3 text-small font-semibold text-text-default">{t.p07.partidas}</h3>
            {s.planned_items.map((p) => <Fila key={p.id} etiqueta={`${p.description}${p.category_id ? ` · ${cats.categoria(p.category_id)}` : ''}`}>{dinero(p.amount_cents)} · {t.p07.frecuencias[p.frequency]} · {fecha(p.start_date)}</Fila>)}
          </div>)}
        </ConDatos>}
      </Bloque>
    </div>
  )
}

function Detalle({ filas, meses }: { filas: Quincena[]; meses: Meses }) {
  const d = (n: number | undefined) => dinero(n ?? 0)
  return (
    <div className="flex flex-col gap-6">
      <Bloque titulo={t.p07.quincenasTitulo}>
        <VistaBD def={{ clave: 'P-07-quincenas', filas, idFila: (q) => q.start_date, vistas: ['tabla'], titulo: (q) => fecha(q.pay_date),
          propiedades: [
            { id: 'pay_date', nombre: t.p07.pago, valor: (q) => q.pay_date, pintar: (q) => fecha(q.pay_date) },
            { id: 'income_cents', nombre: t.p07.ingreso, valor: (q) => q.income_cents, pintar: (q) => `${d(q.income_cents)}${q.income_is_expected ? ` (${t.p07.esperado})` : ''}`, alinear: 'derecha' },
            { id: 'credit_payments_cents', nombre: t.p07.creditos, valor: (q) => q.credit_payments_cents, pintar: (q) => d(q.credit_payments_cents), alinear: 'derecha' },
            { id: 'fixed_and_upcoming_cents', nombre: t.p07.fijos, valor: (q) => q.fixed_and_upcoming_cents, pintar: (q) => d(q.fixed_and_upcoming_cents), alinear: 'derecha' },
            { id: 'variable_estimate_cents', nombre: t.p07.variableCol, valor: (q) => q.variable_estimate_cents, pintar: (q) => d(q.variable_estimate_cents), alinear: 'derecha' },
            { id: 'planned_savings_cents', nombre: t.p07.ahorroCol, valor: (q) => q.planned_savings_cents ?? 0, pintar: (q) => d(q.planned_savings_cents), alinear: 'derecha' },
            { id: 'balance_cents', nombre: t.p07.saldo, valor: (q) => q.balance_cents, pintar: (q) => d(q.balance_cents), alinear: 'derecha' },
            { id: 'cumulative_balance_cents', nombre: t.p07.acumulado, valor: (q) => q.cumulative_balance_cents, pintar: (q) => d(q.cumulative_balance_cents), alinear: 'derecha' },
            { id: 'is_actual', nombre: t.comun.estado, valor: (q) => (q.is_actual ? t.p07.real : t.p07.proyectado), grupo: (q) => (q.is_actual ? t.p07.real : t.p07.proyectado) },
          ] }} />
      </Bloque>
      <Bloque titulo={t.p07.mesesTitulo}>
        <ConDatos consulta={meses}>{(m) => (
          <VistaBD def={{ clave: 'P-07-meses', filas: m.months, idFila: (x) => x.month ?? '', vistas: ['tabla'], titulo: (x) => mes(x.month ?? ''),
            propiedades: [
              { id: 'month', nombre: t.p07.mes, valor: (x) => x.month ?? '', pintar: (x) => mes(x.month ?? '') },
              { id: 'accrued_cents', nombre: t.p07.devengado, valor: (x) => x.accrued_cents ?? 0, pintar: (x) => d(x.accrued_cents), alinear: 'derecha' },
              { id: 'cash_cents', nombre: t.p07.efectivoCol, valor: (x) => x.cash_cents ?? 0, pintar: (x) => d(x.cash_cents), alinear: 'derecha' },
              { id: 'accrued_household_cents', nombre: `${t.p07.devengado} · ${t.comun.hogar}`, valor: (x) => x.accrued_household_cents ?? 0, pintar: (x) => d(x.accrued_household_cents), alinear: 'derecha' },
              { id: 'accrued_external_cents', nombre: `${t.p07.devengado} · ${t.comun.karen}`, valor: (x) => x.accrued_external_cents ?? 0, pintar: (x) => d(x.accrued_external_cents), alinear: 'derecha' },
              { id: 'cumulative_difference_cents', nombre: t.p07.diferenciaCol, valor: (x) => x.cumulative_difference_cents ?? 0, pintar: (x) => d(x.cumulative_difference_cents), alinear: 'derecha' },
            ] }} />)}
        </ConDatos>
      </Bloque>
    </div>
  )
}

function Tendencias({ meses }: { meses: Meses }) {
  return (
    <ConDatos consulta={meses}>{(m) => m.months.length === 0 ? <Vacio texto={t.p07.vacio} /> : (
      <Grafica titulo={t.p07.mesesTitulo} resumen={t.p07.resumenGrafica} tipo="lineas" categorias={m.months.map((x) => mes(x.month ?? ''))}
        series={[{ nombre: t.p07.devengado, valores: m.months.map((x) => x.accrued_cents ?? 0) }, { nombre: t.p07.efectivoCol, valores: m.months.map((x) => x.cash_cents ?? 0) }, { nombre: t.p07.diferenciaCol, valores: m.months.map((x) => x.cumulative_difference_cents ?? 0) }]} />)}
    </ConDatos>
  )
}
