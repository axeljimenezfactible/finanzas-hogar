// P-20 Comparador de escenarios (U1 §2.1, §2.2, §3.22: N-225, N-226; K16). El servidor evalúa (motor S3); aquí solo
// se capturan los datos y se muestran: recomendación, 4 óptimos lado a lado, válidos y descartados con su motivo.
import { useState } from 'react'
import { Enlace, navegar } from '../../app/enrutador'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { usarConsulta } from '../../api/consulta'
import { aCentavos, dinero, fecha, hoyISO, porcentaje } from '../../api/formatos'
import type { Cuerpo, Esquemas } from '../../api/tipos'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ChipSeleccionable } from '../../componentes/Chip'
import { Cargando, ErrorCarga, Franja } from '../../componentes/Estados'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { Datos, Seccion, itemsDe, rutaCondiciones, usarPersonas } from './comun'
import { OBJETIVOS, RESTRICCION, plazo, t } from './textos'

const k = t.comparador
type Escenario = Esquemas['Scenario']
type Resultado = Esquemas['ScenarioResult']
type Pedido = Cuerpo<'/api/v1/scenario-comparisons', 'post'>
const PESOS = Object.keys(OBJETIVOS) as (keyof typeof OBJETIVOS)[]
const entero = (v: string) => (v.trim() && /^\d+$/.test(v.trim()) ? Number(v.trim()) : undefined)
/** Centavos → texto editable («1500.00»), sin flotantes. */
const textoCentavos = (c?: number | null) => (c == null ? '' : `${Math.floor(c / 100)}.${String(c % 100).padStart(2, '0')}`)
type Restriccion = 'colchon' | 'vivas' | 'cero'

export default function P20Comparador({ consulta }: PropsPantalla) {
  const enLinea = useConexion()
  const per = usarPersonas()
  const lineas = usarConsulta('get', '/api/v1/credits', { consulta: { limit: 200 } })
  const [f, setF] = useState({ monto: consulta.get('monto') ?? '', fecha: hoyISO(), holgura: '', comercio: '', quien: 'hogar', pct: '50', colchon: '', vivas: '', cero: '' })
  const [pesos, setPesos] = useState<Record<string, string>>({})
  // N-225: las restricciones duras y los pesos globales (P-26.05, getOptimizerSettings) llegan precargados con «De Ajustes»;
  // lo que se cambie aquí vale solo para esta comparación (las globales no cambian).
  const ajustes = usarConsulta('get', '/api/v1/settings/optimizer')
  const [tocado, setTocado] = useState<Record<string, boolean>>({})
  const g = ajustes.datos
  const deAjustes: Record<Restriccion, string> = { colchon: textoCentavos(g?.min_available_buffer_cents), vivas: g?.max_live_purchases != null ? String(g.max_live_purchases) : '', cero: g?.zero_debt_by_date ?? '' }
  const valor = (c: Restriccion) => (tocado[c] ? f[c] : deAjustes[c])
  const valorPeso = (p: (typeof PESOS)[number]) => (tocado[`peso.${p}`] ? (pesos[p] ?? '') : g ? String(g.weights[p]) : '')
  const ayuda = (clave: string, actual: string, global: string) => (!global ? undefined : !tocado[clave] || actual === global ? k.deAjustes : k.soloEsta)
  const sumaPesos = PESOS.reduce((a, p) => a + (entero(valorPeso(p)) ?? 0), 0)
  /** Peso normalizado en % entero, ROUND_HALF_UP en enteros. */
  const normalizado = (p: (typeof PESOS)[number]) => Math.floor(((entero(valorPeso(p)) ?? 0) * 200 + sumaPesos) / (2 * sumaPesos))
  const [error, setError] = useState('')
  const [res, setRes] = useState<{ estado: 'inicial' | 'evaluando' | 'error' | 'listo'; datos?: Resultado; error?: unknown }>({ estado: 'inicial' })
  if (!enLinea) return <Franja rol="sinConexion">{textos.conexion.requiere}</Franja>
  const nombre = (id: number) => itemsDe<Esquemas['Credit']>(lineas.datos).find((l) => l.id === id)?.name ?? ''
  const cambiar = (c: Partial<typeof f>) => setF({ ...f, ...c })
  const cambiarRestriccion = (c: Restriccion, v: string) => { setF({ ...f, [c]: v }); setTocado({ ...tocado, [c]: true }) }
  const comparar = async () => {
    const monto = aCentavos(f.monto)
    if (!monto || monto <= 0) return setError(textos.registro.montoInvalido)
    setError('')
    const hogar = per.titular?.id, karen = per.externaPrincipal?.id
    const pct = Math.min(100, entero(f.pct) ?? 0), parteHogar = Math.floor((monto * pct + 50) / 100)
    const shares = f.quien === 'hogar' && hogar ? [{ person_id: hogar, principal_cents: monto }] : f.quien === 'karen' && karen ? [{ person_id: karen, principal_cents: monto }] : hogar && karen ? [{ person_id: hogar, principal_cents: parteHogar }, { person_id: karen, principal_cents: monto - parteHogar }].filter((s) => s.principal_cents > 0) : undefined
    const colchon = valor('colchon').trim() ? aCentavos(valor('colchon')) : null
    const cuerpo: Pedido = { amount_cents: monto, desired_date: f.fecha, ...(f.comercio.trim() ? { merchant_name: f.comercio.trim() } : {}), ...(shares ? { shares } : {}),
      ...(entero(f.holgura) != null ? { date_window_days: entero(f.holgura) } : {}), ...(colchon != null ? { min_available_buffer_cents: colchon } : {}), ...(entero(valor('vivas')) != null ? { max_live_purchases: entero(valor('vivas')) } : {}), ...(valor('cero') ? { zero_debt_by_date: valor('cero') } : {}),
      ...(PESOS.some((p) => valorPeso(p).trim()) ? { weights: Object.fromEntries(PESOS.map((p) => [p, entero(valorPeso(p)) ?? 0])) as Pedido['weights'] } : {}) }
    setRes({ estado: 'evaluando' })
    try { setRes({ estado: 'listo', datos: await pedir('post', '/api/v1/scenario-comparisons', { cuerpo }) }) } catch (e) { setRes({ estado: 'error', error: e }) }
  }
  const usar = (s: Escenario) => {
    const q = new URLSearchParams({ monto: f.monto, fecha: s.purchase_date, plazo: String(s.term_count), ...(s.promotion_id ? { promocion: String(s.promotion_id) } : {}), ...(f.comercio.trim() ? { comercio: f.comercio.trim() } : {}) })
    navegar(`${rutaDe('P-19', { id: s.credit_id, compra: 'nueva' })}?${q}`)
  }
  return (
    <div className="grid grid-cols-1 gap-6 tabular-nums desktop:grid-cols-3">
      <div className="flex flex-col gap-4">
        <Seccion titulo={k.datos}>
          <Campo etiqueta={k.monto} inputMode="decimal" value={f.monto} error={error} onChange={(e) => cambiar({ monto: e.target.value })} />
          <Campo etiqueta={k.fecha} type="date" value={f.fecha} onChange={(e) => cambiar({ fecha: e.target.value })} />
          <Campo etiqueta={k.holgura} inputMode="numeric" value={f.holgura} onChange={(e) => cambiar({ holgura: e.target.value })} />
          <Campo etiqueta={k.comercio} value={f.comercio} onChange={(e) => cambiar({ comercio: e.target.value })} />
          <fieldset className="flex flex-col gap-2"><legend className="text-small font-medium">{k.quien}</legend>
            <div className="flex flex-wrap gap-2">{[['hogar', t.comun.hogar], ['karen', per.externaPrincipal?.display_name ?? 'Karen'], ['conjunto', k.conjunto]].map(([v, n]) => <ChipSeleccionable key={v} seleccionado={f.quien === v} onClick={() => cambiar({ quien: v })}>{n}</ChipSeleccionable>)}</div>
            {f.quien === 'conjunto' && <Campo etiqueta={k.pctHogar} inputMode="numeric" value={f.pct} onChange={(e) => cambiar({ pct: e.target.value })} />}
          </fieldset>
        </Seccion>
        <details className="rounded-lg border border-border-divider p-3" open={res.datos?.status === 'missing_constraints' || undefined}>
          <summary className="flex min-h-11 cursor-pointer items-center text-small font-medium">{k.restricciones}</summary>
          <div className="flex flex-col gap-3 pt-2"><p className="text-caption text-text-secondary">{k.restriccionesAyuda}</p>
            <Campo etiqueta={k.colchon} inputMode="decimal" value={valor('colchon')} ayuda={ayuda('colchon', valor('colchon'), deAjustes.colchon)} onChange={(e) => cambiarRestriccion('colchon', e.target.value)} />
            <Campo etiqueta={k.maxVivas} inputMode="numeric" value={valor('vivas')} ayuda={ayuda('vivas', valor('vivas'), deAjustes.vivas)} onChange={(e) => cambiarRestriccion('vivas', e.target.value)} />
            <Campo etiqueta={k.deudaCero} type="date" value={valor('cero')} ayuda={ayuda('cero', valor('cero'), deAjustes.cero)} onChange={(e) => cambiarRestriccion('cero', e.target.value)} />
          </div>
        </details>
        <details className="rounded-lg border border-border-divider p-3">
          <summary className="flex min-h-11 cursor-pointer items-center text-small font-medium">{k.pesos}</summary>
          <div className="flex flex-col gap-3 pt-2"><p className="text-caption text-text-secondary">{k.pesosAyuda}</p>
            {PESOS.map((p) => <Campo key={p} etiqueta={OBJETIVOS[p]} inputMode="numeric" value={valorPeso(p)} ayuda={ayuda(`peso.${p}`, valorPeso(p), g ? String(g.weights[p]) : '')} onChange={(e) => { setPesos({ ...pesos, [p]: e.target.value }); setTocado({ ...tocado, [`peso.${p}`]: true }) }} />)}
            {sumaPesos > 0 && <p className="text-small text-text-secondary" data-prueba="pesos-normalizados">{k.normalizados}: {PESOS.map((p) => `${OBJETIVOS[p]} ${porcentaje(normalizado(p))}`).join(' · ')}</p>}
          </div>
        </details>
        <section aria-labelledby="t-topes" className="rounded-lg border border-border-divider p-3" data-prueba="topes-persona">
          <h3 id="t-topes" className="text-small font-medium">{k.topes}</h3>
          <ul className="flex flex-col gap-1 pt-2 text-small">
            {per.lista.map((p) => <li key={p.id} className="flex justify-between gap-2"><span>{p.display_name}</span><span className="tabular-nums">{p.fortnight_payment_cap_cents != null ? dinero(p.fortnight_payment_cap_cents) : k.sinTope}</span></li>)}
          </ul>
          <p className="pt-2 text-caption text-text-secondary">{k.topesAyuda}</p>
        </section>
        <Boton variante="primario" tam="hoja" onClick={comparar}>{k.comparar}</Boton>
      </div>
      <div className="flex flex-col gap-6 desktop:col-span-2" aria-live="polite">
        {res.estado === 'inicial' && <p className="text-small text-text-secondary">{k.inicial}</p>}
        {res.estado === 'evaluando' && <><p className="text-small text-text-secondary">{k.evaluando}</p><Cargando /></>}
        {res.estado === 'error' && <ErrorCarga error={res.error} alReintentar={comparar} />}
        {res.estado === 'listo' && res.datos && <ResultadoVista r={res.datos} nombre={nombre} usar={usar} />}
      </div>
    </div>
  )
}

function Tarjeta({ s, titulo, nombre, usar, principal }: { s: Escenario; titulo: string; nombre: (id: number) => string; usar: (s: Escenario) => void; principal?: boolean }) {
  return (
    <article className="flex flex-col gap-2 rounded-lg border border-border-divider bg-surface-card p-4" data-prueba="escenario">
      <h3 className="text-small font-semibold text-text-default">{titulo}</h3>
      <p className="text-body text-text-default">{nombre(s.credit_id)} · {fecha(s.purchase_date)} · {plazo(s.term_count)}</p>
      <Datos className="border-0 p-0" filas={[[k.intereses, dinero(s.interest_cents)], [k.pagoMax, dinero(s.max_fortnight_payment_cents)], [k.diferencia, dinero(s.fortnight_payment_spread_cents)], [k.dispMin, dinero(s.min_available_cents)], [k.vivasMax, String(s.max_live_purchases)], [k.deudaCeroEl, fecha(s.zero_debt_date)], [k.rendimiento, s.value_per_peso_milli == null ? t.comun.sinDato : (s.value_per_peso_milli / 1000).toFixed(3)]]} />
      <Boton variante={principal ? 'primario' : 'secundario'} onClick={() => usar(s)}>{principal ? k.usarRecomendado : k.usarEste}</Boton>
    </article>
  )
}
function ResultadoVista({ r, nombre, usar }: { r: Resultado; nombre: (id: number) => string; usar: (s: Escenario) => void }) {
  const cuenta = new Map<string, number>()
  for (const x of r.rejected) cuenta.set(x.violated_constraint, (cuenta.get(x.violated_constraint) ?? 0) + 1)
  const peor = [...cuenta.entries()].sort((a, b) => b[1] - a[1])[0]
  return <>
    <p className="text-small text-text-secondary">{k.resultado(r.evaluated_count, r.valid.length, r.rejected.length)}</p>
    {r.status === 'missing_constraints' && <Franja rol="advertencia" accion={<Enlace a={rutaCondiciones()} className="inline-flex min-h-11 items-center text-small text-text-link underline">{k.capturar}</Enlace>}>{k.faltanRestricciones} · {k.faltan(r.missing_constraints.map((m) => RESTRICCION[m] ?? m).join(', '))}</Franja>}
    {r.status === 'no_valid_scenario' && <Franja rol="error" alerta>{k.ninguno}{peor ? ` · ${k.masDescarta(RESTRICCION[peor[0]] ?? peor[0], peor[1])}` : ''}</Franja>}
    {r.status === 'ok' && r.recommendation && <Tarjeta s={r.recommendation} titulo={k.recomendado((r.recommendation.score_bp ?? 0) / 100 + '')} nombre={nombre} usar={usar} principal />}
    {r.status === 'ok' && <Seccion titulo={k.optimos}><div className="grid grid-cols-1 gap-3 tablet:grid-cols-2 wide:grid-cols-4">{(Object.keys(OBJETIVOS) as (keyof typeof OBJETIVOS)[]).map((o) => r.optima[o] && <Tarjeta key={o} s={r.optima[o]} titulo={OBJETIVOS[o]} nombre={nombre} usar={usar} />)}</div></Seccion>}
    {r.valid.length > 0 && <Seccion titulo={k.validos}><VistaBD def={{ clave: 'P-20-validos', filas: r.valid.map((s, i) => ({ ...s, i })), idFila: (s) => s.i, vistas: ['tabla'], titulo: (s) => nombre(s.credit_id), alAbrir: (s) => usar(s),
      propiedades: [{ id: 'l', nombre: t.creditos.linea, valor: (s) => nombre(s.credit_id) }, { id: 'f', nombre: t.comun.fecha, valor: (s) => s.purchase_date, pintar: (s) => fecha(s.purchase_date) }, { id: 'p', nombre: t.creditos.plazo, valor: (s) => s.term_count, pintar: (s) => plazo(s.term_count) }, { id: 'i', nombre: k.intereses, valor: (s) => s.interest_cents, pintar: (s) => dinero(s.interest_cents), alinear: 'derecha' }, { id: 'q', nombre: k.pagoMax, valor: (s) => s.max_fortnight_payment_cents, pintar: (s) => dinero(s.max_fortnight_payment_cents), alinear: 'derecha' }, { id: 'd', nombre: k.dispMin, valor: (s) => s.min_available_cents, pintar: (s) => dinero(s.min_available_cents), alinear: 'derecha' }] }} /></Seccion>}
    {r.rejected.length > 0 && <Seccion titulo={k.descartados}><Datos filas={r.rejected.map((x, i) => [`${i + 1}. ${nombre(x.scenario.credit_id)} · ${fecha(x.scenario.purchase_date)} · ${plazo(x.scenario.term_count)}`, RESTRICCION[x.violated_constraint] ?? x.violated_constraint])} /></Seccion>}
  </>
}
