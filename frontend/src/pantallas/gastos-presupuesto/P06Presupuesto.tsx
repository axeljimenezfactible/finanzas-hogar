// P-06 Presupuesto del mes (U1 §2.1, §2.2, §3.8: N-080 a N-087; K2 alcance del ajuste). Resumen* · Detalle · Tendencias · Calendario.
import { useState } from 'react'
import { SlidersHorizontal } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { conConsulta, navegar } from '../../app/enrutador'
import { nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { aCentavos, dinero, fecha, mes, porcentaje } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { useAvisos } from '../../componentes/Avisos'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Franja, Vacio } from '../../componentes/Estados'
import { Grafica } from '../../componentes/Grafica'
import { Hoja } from '../../componentes/Hoja'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { textos } from '../../textos/es-MX'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { BarraAvance, Bloque, Fila, SelectorMes, Selector, centavosATexto, mensajeDe, periodoDe, primerDia, rejillaTarjetas, ultimoDia, useCatalogos, useSinRed } from './comun'
import { t } from './textos'

type Crudo = Esquemas['BudgetCategoryRows']['rows'][number]
type FilaCat = { category_id: number; planned_cents: number; actual_cents: number; difference_cents: number; used_pct: number; status: 'normal' | 'warning' | 'exceeded'; expense_count: number }
const normal = (f: Crudo): FilaCat => ({ category_id: f.category_id ?? 0, planned_cents: f.planned_cents ?? 0, actual_cents: f.actual_cents ?? 0, difference_cents: f.difference_cents ?? 0, used_pct: f.used_pct ?? 0, status: f.status ?? 'normal', expense_count: f.expense_count ?? 0 })

export default function P06Presupuesto({ sub, consulta }: PropsPantalla) {
  const periodo = periodoDe(consulta)
  const cats = useCatalogos()
  const sinRed = useSinRed()
  const porCat = usarConsulta('get', '/api/v1/budget/categories', { consulta: { period: periodo } })
  const resumen = usarConsulta('get', '/api/v1/budget/month', { consulta: { period: periodo } })
  const [ajuste, setAjuste] = useState<number | null>(null)
  const recargar = () => { porCat.recargar(); resumen.recargar() }
  const ajustar = (categoria?: number) => setAjuste(categoria ?? porCat.datos?.rows[0]?.category_id ?? cats.categorias[0]?.id ?? 0)
  const vacio = <Vacio texto={t.p06.vacio} accion={<Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => ajustar()}>{t.p06.ajustar}</Boton>} />
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-2">
        <SelectorMes valor={periodo} />
        <Boton variante="primario" className="ml-auto" deshabilitadoPor={sinRed} onClick={() => ajustar()}><SlidersHorizontal aria-hidden className="size-5" />{t.p06.ajustar}</Boton>
      </div>
      {sub === 'resumen' && <Resumen periodo={periodo} c={resumen} vacio={vacio} />}
      {sub === 'detalle' && <Detalle periodo={periodo} c={porCat} vacio={vacio} alAjustar={sinRed ? undefined : ajustar} />}
      {sub === 'tendencias' && <Tendencias consulta={consulta} />}
      {sub === 'calendario' && <Calendario periodo={periodo} />}
      {ajuste !== null && <Ajustar periodo={periodo} categoria={ajuste} filas={(porCat.datos?.rows ?? []).map(normal)} alCerrar={(ok) => { setAjuste(null); if (ok) recargar() }} />}
    </div>
  )
}

function Resumen({ periodo, c, vacio }: { periodo: string; c: ReturnType<typeof usarConsulta<'/api/v1/budget/month'>>; vacio: React.ReactNode }) {
  const alertas = usarConsulta('get', '/api/v1/budget/alerts')
  const quincenas = usarConsulta('get', '/api/v1/budget/fortnights', { consulta: { period: periodo } })
  const cats = useCatalogos()
  const detalle = () => navegar(conConsulta({ sub: 'detalle' }))
  return (
    <ConDatos consulta={c}>{(m) => m.planned_spending_cents === 0 && m.actual_cents === 0 ? vacio : (
      <div className="flex flex-col gap-6">
        {(alertas.datos?.rows ?? []).length > 0 && (
          <div className="flex flex-col gap-2" aria-label={t.p06.alertas}>
            {alertas.datos!.rows.map((a) => <Franja key={`${a.category_id}-${a.threshold_pct}`} rol="advertencia">{t.p06.alerta(a.category_id == null ? t.p06.total : cats.categoria(a.category_id), porcentaje(a.used_pct), fecha(a.crossed_date))}</Franja>)}
          </div>)}
        <div className={rejillaTarjetas}>
          <TarjetaMetrica titulo={t.p06.ingreso} rol="income" valor={dinero(m.income_cents)} alAbrir={() => navegar(rutaDe('P-08'))} />
          <TarjetaMetrica titulo={t.p06.plan} valor={dinero(m.planned_spending_cents)} alAbrir={detalle} />
          <TarjetaMetrica titulo={t.p06.ahorro} rol="savings" valor={dinero(m.planned_savings_cents)} alAbrir={() => navegar(rutaDe('P-09'))} />
          <TarjetaMetrica titulo={t.p06.sinAsignar} valor={dinero(m.unassigned_cents)} alAbrir={detalle} />
          <TarjetaMetrica titulo={t.p06.real} rol="expense" valor={dinero(m.actual_cents)} alAbrir={detalle}>
            <BarraAvance nombre={t.p06.real} pct={m.used_pct ?? 0} estado={m.status} excedente={m.remaining_cents < 0 ? dinero(-m.remaining_cents, 'gasto') : undefined} />
          </TarjetaMetrica>
          <TarjetaMetrica titulo={t.p06.restante} valor={dinero(m.remaining_cents)} nota={t.p06.diasRestantes(m.days_left)} alAbrir={detalle} />
        </div>
        <Bloque titulo={t.p06.quincenas}>
          <ConDatos consulta={quincenas}>{(q) => <div>{q.fortnights.map((f) => (
            <Fila key={f.half} etiqueta={`${t.p06.quincena(f.half)} · ${fecha(f.start_date)} – ${fecha(f.end_date)}`}>
              <span className="flex w-full flex-col gap-1"><span>{t.p06.realCol} {dinero(f.actual_cents)} · {t.p06.planCol} {dinero(f.planned_cents)}</span>
                <BarraAvance nombre={t.p06.quincena(f.half)} pct={f.planned_cents > 0 ? Math.floor((f.actual_cents * 100) / f.planned_cents) : 0} /></span>
            </Fila>))}</div>}
          </ConDatos>
        </Bloque>
      </div>)}
    </ConDatos>
  )
}

function Detalle({ periodo, c, vacio, alAjustar }: { periodo: string; c: ReturnType<typeof usarConsulta<'/api/v1/budget/categories'>>; vacio: React.ReactNode; alAjustar?: (categoria: number) => void }) {
  const cats = useCatalogos()
  const despensa = usarConsulta('get', '/api/v1/budget/accrued-groceries', { consulta: { period: periodo } })
  const verGastos = (f: FilaCat) => navegar(`${rutaDe('P-02')}?periodo=${periodo}&categoria=${f.category_id}`)
  return (
    <div className="flex flex-col gap-6">
      <ConDatos consulta={c}>{(d) => (
        <VistaBD def={{
          clave: 'P-06-categorias', filas: d.rows.map(normal), idFila: (f) => f.category_id, vistas: ['tabla', 'tablero'], vacio, titulo: (f) => cats.categoria(f.category_id),
          alAbrir: alAjustar ? (f) => alAjustar(f.category_id) : verGastos,
          columnasTablero: { propiedad: 'status', columnas: Object.values(t.estadosPresupuesto) },
          propiedades: [
            { id: 'category_id', nombre: t.comun.categoria, valor: (f) => cats.categoria(f.category_id) },
            { id: 'planned_cents', nombre: t.p06.planCol, valor: (f) => f.planned_cents, pintar: (f) => dinero(f.planned_cents), alinear: 'derecha' },
            { id: 'actual_cents', nombre: t.p06.realCol, valor: (f) => f.actual_cents, pintar: (f) => dinero(f.actual_cents), alinear: 'derecha' },
            { id: 'difference_cents', nombre: t.p06.diferencia, valor: (f) => f.difference_cents, pintar: (f) => dinero(f.difference_cents), alinear: 'derecha' },
            { id: 'used_pct', nombre: t.p06.pctUsado, valor: (f) => f.used_pct, pintar: (f) => <span className="flex min-w-24 flex-col gap-1"><span className="tabular-nums">{porcentaje(f.used_pct)}</span><BarraAvance nombre={cats.categoria(f.category_id)} pct={f.used_pct} estado={f.status} /></span> },
            { id: 'status', nombre: t.comun.estado, valor: (f) => t.estadosPresupuesto[f.status], grupo: (f) => t.estadosPresupuesto[f.status] },
            { id: 'expense_count', nombre: t.p06.numGastos, valor: (f) => f.expense_count, pintar: (f) => <Boton variante="enlace" onClick={(e) => { e.stopPropagation(); verGastos(f) }}>{t.p06.verGastos} ({f.expense_count})</Boton>, alinear: 'derecha' },
          ],
        }} />)}
      </ConDatos>
      <Bloque titulo={t.p06.despensa}>
        <ConDatos consulta={despensa}>{(g) => g.rows.length === 0 ? <p className="text-small text-text-secondary">{t.p06.sinDespensa}</p> : (
          <VistaBD def={{ clave: 'P-06-despensa', filas: g.rows, idFila: (r) => r.expense_id, vistas: ['tabla'], titulo: (r) => fecha(r.expense_date), alAbrir: (r) => navegar(rutaDe('P-04', { id: r.expense_id })),
            propiedades: [
              { id: 'expense_date', nombre: t.p06.fechaCompra, valor: (r) => r.expense_date, pintar: (r) => fecha(r.expense_date) },
              { id: 'total_cents', nombre: t.p06.totalIntereses, valor: (r) => r.total_cents, pintar: (r) => dinero(r.total_cents), alinear: 'derecha' },
              { id: 'accrued_cents', nombre: t.p06.devengadoMes, valor: (r) => r.accrued_cents, pintar: (r) => dinero(r.accrued_cents), alinear: 'derecha' },
              { id: 'household_share_cents', nombre: t.p02.parteHogar, valor: (r) => r.household_share_cents, pintar: (r) => dinero(r.household_share_cents), alinear: 'derecha' },
              { id: 'external_share_cents', nombre: t.p02.parteKaren, valor: (r) => r.external_share_cents, pintar: (r) => dinero(r.external_share_cents), alinear: 'derecha' },
            ] }} />)}
        </ConDatos>
      </Bloque>
    </div>
  )
}

function Tendencias({ consulta }: { consulta: URLSearchParams }) {
  const cat = Number(consulta.get('categoria'))
  const cats = useCatalogos()
  const r = usarConsulta('get', '/api/v1/budget/monthly-trend', { consulta: cat > 0 ? { category_id: [cat] } : {} })
  return (
    <div className="flex flex-col gap-3">
      <Selector className="tablet:w-64" etiqueta={t.comun.categoria} valor={cat > 0 ? String(cat) : ''} vacio={t.p06.total} alCambiar={(v) => navegar(conConsulta({ categoria: v || null }), { reemplazar: true })} opciones={cats.categorias.map((c) => ({ valor: String(c.id), texto: c.name }))} />
      <ConDatos consulta={r}>{(d) => d.months.length === 0 ? <Vacio texto={textos.vacio.grafica} /> : (
        <Grafica titulo={t.p06.tendencia} resumen={t.p06.resumenTendencia} tipo="barras" categorias={d.months.map((m) => mes(m.month))}
          series={[{ nombre: t.p06.planCol, valores: d.months.map((m) => m.planned_cents ?? 0) }, { nombre: t.p06.realCol, valores: d.months.map((m) => m.actual_cents) }]} />)}
      </ConDatos>
    </div>
  )
}

/** N-086: calendario del mes limitado a eventos de dinero (listCalendarEvents con los tipos de dinero). */
export function Calendario({ periodo, tipos = ['income', 'credit_payment', 'fixed_expense', 'upcoming_expense'], clave = 'P-06-calendario' }: { periodo: string; tipos?: Esquemas['CalendarEvent']['event_type'][]; clave?: string }) {
  const r = usarConsulta('get', '/api/v1/calendar/events', { consulta: { period: periodo, event_type: tipos, limit: 200 } })
  return (
    <ConDatos consulta={r}>{(d) => (
      <VistaBD def={{ clave, filas: d.items as Esquemas['CalendarEvent'][], idFila: (e) => `${e.ref_table}-${e.ref_id}-${e.date}`, vistas: ['calendario'], fecha: (e) => e.date, vacio: <Vacio texto={t.p06.sinEventos} />,
        titulo: (e) => <span>{e.title}{e.amount_cents != null && <span className="tabular-nums"> · {dinero(e.amount_cents)}</span>}</span>,
        propiedades: [{ id: 'title', nombre: t.comun.concepto, valor: (e) => e.title }, { id: 'amount_cents', nombre: t.comun.monto, valor: (e) => e.amount_cents ?? null, pintar: (e) => (e.amount_cents == null ? '—' : dinero(e.amount_cents)), alinear: 'derecha' }] }} />)}
    </ConDatos>
  )
}

/** N-087 Ajustar presupuesto con alcance K2 («De este mes en adelante» preseleccionado) y vista previa «Antes → Después». */
function Ajustar({ periodo, categoria, filas, alCerrar }: { periodo: string; categoria: number; filas: FilaCat[]; alCerrar: (ok: boolean) => void }) {
  const cats = useCatalogos()
  const avisar = useAvisos()
  const [cat, setCat] = useState(categoria)
  const actual = filas.find((f) => f.category_id === cat)
  const [monto, setMonto] = useState(actual ? centavosATexto(actual.planned_cents) : '')
  const [alcance, setAlcance] = useState<'onward' | 'this_month'>('onward')
  const [error, setError] = useState('')
  const [enCurso, setEnCurso] = useState(false)
  const nuevo = aCentavos(monto)
  const valido = nuevo != null && nuevo >= 0
  const guardar = async () => {
    if (!valido) return setError(textos.registro.montoInvalido)
    setEnCurso(true)
    try {
      const vigentes = await pedir('get', '/api/v1/budgets', { consulta: { period: periodo, category_id: [cat] } })
      const renglon = (vigentes.items as Esquemas['Budget'][]).find((b) => b.category_id === cat)
      if (renglon) await pedir('patch', '/api/v1/budgets/{budget_id}', { ruta: { budget_id: renglon.id }, consulta: { scope: alcance, month: periodo }, cuerpo: { amount_cents: nuevo } })
      else await pedir('post', '/api/v1/budgets', { idempotencia: nuevaClave(), cuerpo: { category_id: cat, amount_cents: nuevo, valid_from_date: primerDia(periodo), valid_to_date: alcance === 'this_month' ? ultimoDia(periodo) : undefined } })
      avisar({ texto: t.p06.ajustado, rol: 'exito' }); alCerrar(true)
    } catch (e) { setError(mensajeDe(e)) } finally { setEnCurso(false) }
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && alCerrar(false)} titulo={t.p06.ajustar} descripcion={mes(periodo)}
      pie={<div className="flex justify-end gap-2"><Boton onClick={() => alCerrar(false)}>{textos.app.cancelar}</Boton><Boton variante="primario" aria-busy={enCurso} deshabilitadoPor={valido ? undefined : textos.registro.montoInvalido} onClick={guardar}>{textos.app.guardar}</Boton></div>}>
      <div className="flex flex-col gap-4">
        <Selector etiqueta={t.comun.categoria} valor={String(cat)} alCambiar={(v) => { const n = Number(v); setCat(n); const f = filas.find((x) => x.category_id === n); setMonto(f ? centavosATexto(f.planned_cents) : '') }} opciones={cats.categorias.map((c) => ({ valor: String(c.id), texto: c.name }))} />
        <Campo etiqueta={t.comun.monto} inputMode="decimal" value={monto} error={error || undefined} onChange={(e) => { setMonto(e.target.value); setError('') }} />
        <fieldset className="flex flex-col gap-2"><legend className="text-small text-text-secondary">{t.p06.alcance}</legend>
          {([['onward', t.p06.enAdelante], ['this_month', t.p06.soloEste]] as const).map(([v, n]) => (
            <label key={v} className="flex min-h-11 items-center gap-2 text-body text-text-default"><input type="radio" name="alcance" className="size-5" checked={alcance === v} onChange={() => setAlcance(v)} />{n}</label>))}
        </fieldset>
        {valido && (
          <dl className="grid grid-cols-3 gap-2 rounded-md bg-surface-hover-subtle p-3 text-small tabular-nums" aria-label={`${t.p06.antes} → ${t.p06.despues}`}>
            <dt className="text-text-secondary">{t.p06.planCol}</dt><dd>{dinero(actual?.planned_cents ?? 0)}</dd><dd>→ {dinero(nuevo)}</dd>
            <dt className="text-text-secondary">{t.p06.restante}</dt><dd>{dinero((actual?.planned_cents ?? 0) - (actual?.actual_cents ?? 0))}</dd><dd>→ {dinero(nuevo - (actual?.actual_cents ?? 0))}</dd>
          </dl>)}
      </div>
    </Hoja>
  )
}
