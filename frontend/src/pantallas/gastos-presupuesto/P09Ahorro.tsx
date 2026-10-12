// P-09 Ahorro y metas (U1 §2.1, §2.2, §3.11: N-110 a N-116). Resumen* · Detalle · Tendencias · Calendario. Las metas se definen en P-26.13.
import { useState } from 'react'
import { ArrowDownToLine, ArrowUpFromLine, Plus, Trash2 } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { conConsulta, navegar } from '../../app/enrutador'
import { nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { aCentavos, dinero, fecha, hoyISO, mes } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { useAvisos } from '../../componentes/Avisos'
import { Boton } from '../../componentes/Boton'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Grafica } from '../../componentes/Grafica'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { textos } from '../../textos/es-MX'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { Calendario } from './P06Presupuesto'
import { BarraAvance, Bloque, Confirmar, HojaCaptura, SelectorMes, Selector, centavosATexto, mensajeDe, periodoDe, rejillaTarjetas, useSinRed } from './comun'
import { t } from './textos'

type Meta = Esquemas['SavingsGoal']
type Mov = Esquemas['SavingsMovement']

export default function P09Ahorro({ sub, consulta }: PropsPantalla) {
  const periodo = periodoDe(consulta)
  const metasC = usarConsulta('get', '/api/v1/savings-goals', { consulta: { status: 'active' } })
  const sinRed = useSinRed()
  const avisar = useAvisos()
  const hoja = consulta.get('hoja')
  const [tipo, setTipo] = useState<'deposit' | 'withdrawal'>('deposit')
  const [destino, setDestino] = useState('')
  const [vuelta, setVuelta] = useState(0)
  const metas = (metasC.datos?.items ?? []) as Meta[]
  const nombreMeta = (id: number) => metas.find((m) => m.id === id)?.name ?? '—'
  const abrir = (h: string | null, k: 'deposit' | 'withdrawal' = 'deposit') => { setTipo(k); navegar(conConsulta({ hoja: h })) }
  const nueva = () => navegar(rutaDe('P-26', { seccion: '13' }))
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-2">
        {(sub === 'detalle' || sub === 'calendario') && <SelectorMes valor={periodo} />}
        <div className="ml-auto flex flex-wrap gap-2">
          <Boton onClick={nueva}><Plus aria-hidden className="size-5" />{t.p09.nuevaMeta}</Boton>
          <Boton deshabilitadoPor={sinRed} onClick={() => abrir('movimiento', 'withdrawal')}><ArrowUpFromLine aria-hidden className="size-5" />{t.p09.retirar}</Boton>
          <Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => abrir('movimiento')}><ArrowDownToLine aria-hidden className="size-5" />{t.p09.aportar}</Boton>
        </div>
      </div>
      {sub === 'calendario' ? <Calendario periodo={periodo} tipos={['savings_goal', 'savings_movement']} clave="P-09-calendario" /> : sub === 'tendencias' ? <Tendencias /> : (
        <ConDatos key={vuelta} consulta={metasC}>{() => metas.length === 0 ? <Vacio texto={t.p09.vacio} accion={<Boton variante="primario" onClick={nueva}>{t.p09.nuevaMeta}</Boton>} />
          : sub === 'detalle' ? <Detalle metas={metas} periodo={periodo} nombreMeta={nombreMeta} /> : <Resumen metas={metas} />}
        </ConDatos>)}
      <HojaCaptura titulo={tipo === 'deposit' ? t.p09.aportar : t.p09.retirar} abierta={hoja === 'movimiento'} alCerrar={() => navegar(conConsulta({ hoja: null }))} fechaInicial={hoyISO()}
        extra={<>
          <Selector etiqueta={t.comun.tipo} valor={tipo} alCambiar={(v) => setTipo(v as 'deposit' | 'withdrawal')} opciones={Object.entries(t.p09.tipos).map(([valor, texto]) => ({ valor, texto }))} />
          <Selector etiqueta={t.p09.destino} valor={destino || String(metas[0]?.id ?? '')} alCambiar={setDestino} opciones={metas.map((m) => ({ valor: String(m.id), texto: m.name }))} />
        </>}
        alGuardar={async (c) => {
          const meta = Number(destino || metas[0]?.id)
          if (!meta) throw new Error(t.p09.vacio)
          await pedir('post', '/api/v1/savings-movements', { idempotencia: nuevaClave(), cuerpo: { savings_goal_id: meta, movement_date: c.fecha, kind: tipo, amount_cents: c.monto, notes: c.nota || undefined } })
          avisar({ texto: t.p09.registrado, rol: 'exito' }); metasC.recargar(); setVuelta((n) => n + 1)
        }} />
    </div>
  )
}

function TarjetaMeta({ m }: { m: Meta }) {
  return (
    <div className="flex flex-col gap-2">
      <span className="flex justify-between text-small text-text-secondary"><span>{m.kind === 'fund' ? t.p09.fondo : t.p09.estados[m.status ?? 'in_progress']}</span>{m.target_date && <span>{fecha(m.target_date)}</span>}</span>
      <span className="tabular-nums text-body text-text-default">{dinero(m.saved_cents ?? 0)}{m.target_cents != null && <span className="text-text-secondary"> / {dinero(m.target_cents)}</span>}</span>
      {m.target_cents != null && <BarraAvance pct={m.progress_pct ?? 0} />}
      {m.required_monthly_cents != null && <span className="text-caption text-text-secondary">{t.p09.mensualNecesaria}: <span className="tabular-nums">{dinero(m.required_monthly_cents)}</span></span>}
    </div>
  )
}

function Resumen({ metas }: { metas: Meta[] }) {
  const r = usarConsulta('get', '/api/v1/savings/summary')
  const detalle = () => navegar(conConsulta({ sub: 'detalle' }))
  return (
    <div className="flex flex-col gap-6">
      <ConDatos consulta={r}>{(s) => (
        <div className={rejillaTarjetas}>
          <TarjetaMetrica titulo={t.p09.ahorrado} rol="savings" valor={dinero(s.saved_total_cents)} alAbrir={detalle} />
          <TarjetaMetrica titulo={t.p09.aportadoMes} rol="savings" valor={dinero(s.contributed_month_cents)} nota={t.p09.contraPlan(dinero(s.planned_month_cents))} alAbrir={detalle}>
            <BarraAvance pct={s.planned_month_cents > 0 ? Math.floor((s.contributed_month_cents * 100) / s.planned_month_cents) : 0} />
          </TarjetaMetrica>
        </div>)}
      </ConDatos>
      <Bloque titulo={t.p09.metas}>
        <ul className={rejillaTarjetas}>
          {metas.map((m) => <li key={m.id} className="rounded-lg border border-border-divider bg-surface-card p-4"><h3 className="pb-2 text-body font-semibold text-text-default">{m.name}</h3><TarjetaMeta m={m} /></li>)}
        </ul>
      </Bloque>
    </div>
  )
}

function Detalle({ metas, periodo, nombreMeta }: { metas: Meta[]; periodo: string; nombreMeta: (id: number) => string }) {
  const movs = usarConsulta('get', '/api/v1/savings-movements', { consulta: { period: periodo, limit: 200 } })
  const sinRed = useSinRed()
  const avisar = useAvisos()
  const [borrar, setBorrar] = useState<number | null>(null)
  const parchar = async (id: number, cuerpo: Esquemas['SavingsMovementUpdate']) => { await pedir('patch', '/api/v1/savings-movements/{savings_movement_id}', { ruta: { savings_movement_id: id }, cuerpo }); movs.recargar() }
  return (
    <div className="flex flex-col gap-6">
      <Bloque titulo={t.p09.metas}>
        <VistaBD def={{ clave: 'P-09-metas', filas: metas, idFila: (m) => m.id, vistas: ['galeria', 'tablero', 'calendario'], fecha: (m) => m.target_date, titulo: (m) => <span className="flex flex-col gap-2"><span>{m.name}</span><TarjetaMeta m={m} /></span>,
          columnasTablero: { propiedad: 'status', columnas: Object.values(t.p09.estados) },
          propiedades: [
            { id: 'name', nombre: t.comun.descripcion, valor: (m) => m.name },
            { id: 'status', nombre: t.comun.estado, valor: (m) => t.p09.estados[m.status ?? 'in_progress'], grupo: (m) => t.p09.estados[m.status ?? 'in_progress'] },
            { id: 'saved_cents', nombre: t.p09.ahorrado, valor: (m) => m.saved_cents ?? 0, pintar: (m) => dinero(m.saved_cents ?? 0), alinear: 'derecha' },
            { id: 'target_cents', nombre: t.p09.meta, valor: (m) => m.target_cents ?? null, pintar: (m) => (m.target_cents == null ? '—' : dinero(m.target_cents)), alinear: 'derecha' },
            { id: 'planned_monthly_cents', nombre: t.p09.mensualPlan, valor: (m) => m.planned_monthly_cents ?? null, pintar: (m) => (m.planned_monthly_cents == null ? '—' : dinero(m.planned_monthly_cents)), alinear: 'derecha' },
          ] }} />
      </Bloque>
      <Bloque titulo={t.p09.movimientos}>
        <ConDatos consulta={movs}>{(d) => (
          <VistaBD def={{ clave: 'P-09-movimientos', filas: d.items as Mov[], idFila: (x) => x.id, vistas: ['tabla'], titulo: (x) => nombreMeta(x.savings_goal_id), vacio: <Vacio texto={t.p09.sinMovimientos} />,
            propiedades: [
              { id: 'movement_date', nombre: t.comun.fecha, valor: (x) => x.movement_date, pintar: (x) => fecha(x.movement_date), editar: { tipo: 'fecha', textoInicial: (x) => x.movement_date, guardar: (x, v) => parchar(x.id, { movement_date: v }) } },
              { id: 'kind', nombre: t.comun.tipo, valor: (x) => t.p09.tipos[x.kind], grupo: (x) => t.p09.tipos[x.kind] },
              { id: 'amount_cents', nombre: t.comun.monto, valor: (x) => x.amount_cents, alinear: 'derecha', pintar: (x) => dinero(x.amount_cents, x.kind === 'withdrawal' ? 'gasto' : 'ingreso'),
                editar: { tipo: 'dinero', textoInicial: (x) => centavosATexto(x.amount_cents), guardar: (x, v) => parchar(x.id, { amount_cents: aCentavos(v)! }),
                  vistaPrevia: async (x, v) => { const c = aCentavos(v); if (c == null || c <= 0) throw new Error(textos.registro.montoInvalido); return [{ etiqueta: t.comun.monto, antes: dinero(x.amount_cents), despues: dinero(c) }] } } },
              { id: 'savings_goal_id', nombre: t.p09.destino, valor: (x) => nombreMeta(x.savings_goal_id), grupo: (x) => nombreMeta(x.savings_goal_id),
                editar: { tipo: 'opcion', opciones: metas.map((m) => ({ valor: String(m.id), texto: m.name })), textoInicial: (x) => String(x.savings_goal_id), guardar: (x, v) => parchar(x.id, { savings_goal_id: Number(v) }),
                  vistaPrevia: async (x, v) => [{ etiqueta: t.p09.destino, antes: nombreMeta(x.savings_goal_id), despues: nombreMeta(Number(v)) }] } },
              { id: 'notes', nombre: t.comun.nota, valor: (x) => x.notes ?? '', editar: { tipo: 'texto', textoInicial: (x) => x.notes ?? '', guardar: (x, v) => parchar(x.id, { notes: v || null }) } },
              { id: 'accion', nombre: '', valor: () => '', pintar: (x) => <Boton tam="icono" variante="fantasma" aria-label={t.comun.eliminar} deshabilitadoPor={sinRed} onClick={(e) => { e.stopPropagation(); setBorrar(x.id) }}><Trash2 aria-hidden className="size-6" size={ICONO_MD} /></Boton> },
            ] }} />)}
        </ConDatos>
      </Bloque>
      <Confirmar abierta={borrar !== null} alCambiar={(a) => !a && setBorrar(null)} titulo={t.comun.eliminar} texto={t.comun.confirmarEliminar} accion={t.comun.eliminar}
        alConfirmar={async () => { try { await pedir('delete', '/api/v1/savings-movements/{savings_movement_id}', { ruta: { savings_movement_id: borrar! } }); avisar({ texto: t.comun.eliminado, rol: 'exito' }); setBorrar(null); movs.recargar() } catch (e) { avisar({ texto: mensajeDe(e), rol: 'error' }) } }} />
    </div>
  )
}

function Tendencias() {
  const r = usarConsulta('get', '/api/v1/savings/monthly-trend')
  return (
    <ConDatos consulta={r}>{(d) => d.months.length === 0 ? <Vacio texto={textos.vacio.grafica} /> : (
      <Grafica titulo={t.p09.tendencia} resumen={t.p09.resumenTendencia} tipo="lineas" categorias={d.months.map((m) => mes(m.month ?? ''))}
        series={[{ nombre: t.p09.alCierre, valores: d.months.map((m) => m.saved_end_cents ?? 0) }, { nombre: t.p09.aportado, valores: d.months.map((m) => m.contributed_cents ?? 0) }, { nombre: t.p09.planificado, valores: d.months.map((m) => m.planned_cents ?? 0) }]} />)}
    </ConDatos>
  )
}
