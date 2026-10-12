// P-26.09 Notificaciones (U4 parte b; decisión F, K4, K44, K47, K50). Claves alerts.* con getAlertSettings/updateAlertSettings;
// chats de Telegram detectados por el bot (K4): nunca se muestran el chat_id completo ni el token (K30).
import { useRef, useState, type ReactNode } from 'react'
import { useSesion } from '../../app/sesion'
import { pedir, nuevaClave } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import { dinero, fechaConDia, fechaHora, hoyISO } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { estilosControl } from '../../componentes/Campo'
import { Hoja } from '../../componentes/Hoja'
import { Esqueleto, Vacio } from '../../componentes/Estados'
import { useAvisos } from '../../componentes/Avisos'
import { cn } from '../../componentes/utils'
import { Bloque, Cargar, DialogoCritico, FilaParametro, IndicadorGuardado, Interruptor, Marco, unir, useEsEscritorio, useGuardado, useSinRed } from './comun'
import { Destinatarios } from './P2608Reportes'
import { divHalfUp, sumarDias } from './calculos'
import { DIAS_SEMANA, textosAjustesB as T } from './textos'
import { GrupoVarias } from './varias'
import type { EditarVarias } from '../../vistas-bd/tipos'

const t = T.n
type N = Esquemas['Notification']; type R = Esquemas['NotificationRecipient']; type Pe = Esquemas['Person']; type C = Esquemas['TelegramChatCandidate']
type Params = Esquemas['AlertSettings']['parameters']
type Clave = keyof Params
const ALERTAS = ['payment_due', 'budget_threshold', 'stale_prices', 'meal_plan_missing'] as const
const guardarAlerta = (parameters?: Esquemas['AlertSettingsUpdate']['parameters'], reset?: Clave[]) =>
  pedir('patch', '/api/v1/settings/alerts', { cuerpo: reset ? { reset_parameters: reset } : { parameters } })
const sinAnio = (iso: string) => fechaConDia(iso).replace(/ \d{4}$/, '')

export default function P2609Notificaciones() {
  const mesActual = hoyISO().slice(0, 7)
  const datos = unir<[{ items: N[] }, Esquemas['AlertSettings'], { items: R[] }, { items: C[] }, { items: Pe[] }]>([
    usarConsulta('get', '/api/v1/notifications', { consulta: { limit: 50 } }),
    usarConsulta('get', '/api/v1/settings/alerts'),
    usarConsulta('get', '/api/v1/notification-recipients', { consulta: { limit: 50 } }),
    usarConsulta('get', '/api/v1/telegram-chat-candidates', { consulta: { limit: 50 } }),
    usarConsulta('get', '/api/v1/persons', { consulta: { limit: 50 } }),
  ])
  const presupuesto = usarConsulta('get', '/api/v1/budget/month', { consulta: { period: mesActual } })
  const pago = usarConsulta('get', '/api/v1/dashboard/next-payment')
  const esqueleto = <div className="flex flex-col gap-3">{[0, 1, 2].map((i) => <Esqueleto key={i} className="h-24" />)}</div>
  return (
    <Marco id="09">
      <p className="text-body text-text-secondary">{t.proposito}</p>
      <Cargar c={datos} que={t.que} esqueleto={esqueleto}>{([ns, aj, rs, cs, ps]) => {
        const ctx = { params: aj.parameters, recips: rs.items, personas: ps.items, recargar: datos.recargar, total: presupuesto.datos?.planned_spending_cents, vence: pago.datos?.next_payment?.due_date }
        const pendientes = cs.items.filter((c) => !c.confirmed_at)
        const alertas = ALERTAS.flatMap((code) => ns.items.filter((x) => x.code === code))
        return (
          <>
            {pendientes.length > 0 && <Detectados cs={pendientes} personas={ps.items} recargar={datos.recargar} />}
            <Bloque titulo={t.alertas} accion={<PruebaAviso alertas={ns.items.filter((n) => (ALERTAS as readonly string[]).includes(n.code))} personas={ps.items} recips={rs.items} />}>
              <GrupoVarias filas={alertas} idFila={(n) => n.id} nombre={(n) => n.name} ev={grupoAlertas(alertas, datos.recargar)}>{(casilla) => (
                <div className="flex flex-col gap-3">
                  {ALERTAS.map((code) => { const n = ns.items.find((x) => x.code === code); return n ? <Alerta key={code} n={n} casilla={casilla(n)} {...ctx} /> : <p key={code} role="alert" className="text-small text-role-error">{T.r.faltaReporte}</p> })}
                </div>
              )}</GrupoVarias>
              <FilaParametro etiqueta={t.horaAvisos} tipo="hora" param={aj.parameters['alerts.send_time']} mensaje={T.validar.hora}
                guardar={async (v: string) => { await guardarAlerta({ 'alerts.send_time': v }); datos.recargar() }} restablecer={async () => { await guardarAlerta(undefined, ['alerts.send_time']); datos.recargar() }} />
            </Bloque>
            <Confirmados recips={rs.items} personas={ps.items} recargar={datos.recargar} />
            {pendientes.length === 0 && <Detectados cs={pendientes} personas={ps.items} recargar={datos.recargar} />}
          </>
        )
      }}</Cargar>
    </Marco>
  )
}

type Ctx = { n: N; params: Params; recips: R[]; personas: Pe[]; recargar: () => void; total?: number; vence?: string }
function Parametros({ n, params, recargar, total, vence }: Ctx) {
  // Umbral por categoría al tocar una categoría de ejemplo (U4 P-26.09, vista previa de umbrales; H-08): partida de getBudgetByCategory del mes en curso.
  const esUmbral = n.code === 'budget_threshold'
  const porCategoria = usarConsulta('get', '/api/v1/budget/categories', { consulta: { period: hoyISO().slice(0, 7) } }, esUmbral)
  const categorias = usarConsulta('get', '/api/v1/categories', { consulta: { limit: 200 } }, esUmbral)
  const [cat, setCat] = useState('')
  const nombreCat = (id: number) => (categorias.datos?.items as Esquemas['Category'][] | undefined)?.find((c) => c.id === id)?.name ?? ''
  const ejemplo = porCategoria.datos?.rows?.find((r) => String(r.category_id) === cat)
  const planeado = ejemplo?.planned_cents ?? 0
  const lineaCategoria = (w: number, l: number) => (ejemplo ? t.umbralesCategoria(nombreCat(ejemplo.category_id), dinero(planeado), dinero(divHalfUp(planeado * w, 100)), dinero(divHalfUp(planeado * l, 100))) : null)
  const p = <V extends number | string | number[]>(clave: Clave, etiqueta: string, tipo: 'entero' | 'lista' | 'diaSemana', mensaje: string, unidad?: string, previa?: (v: V) => { cambios: { etiqueta: string; antes: string; despues: string }[]; nota?: ReactNode } | null, ayuda?: string) => (
    <FilaParametro<V> key={clave} etiqueta={etiqueta} tipo={tipo} param={params[clave] as never} mensaje={mensaje} unidad={unidad} ayuda={ayuda} vistaPrevia={previa}
      guardar={async (v) => { await guardarAlerta({ [clave]: v }); recargar() }} restablecer={async () => { await guardarAlerta(undefined, [clave]); recargar() }} />
  )
  const avisoA = (pct: number) => (total === undefined ? '—' : dinero(divHalfUp(total * pct, 100)))
  const umbral = (clave: 'alerts.budget_warning_pct' | 'alerts.budget_limit_pct', etiqueta: string) => p<number>(clave, etiqueta, 'entero', T.validar.porcentaje, '%', (v) => {
    const w = clave === 'alerts.budget_warning_pct' ? v : (params['alerts.budget_warning_pct'].value ?? 0)
    const l = clave === 'alerts.budget_limit_pct' ? v : (params['alerts.budget_limit_pct'].value ?? 0)
    return { cambios: [{ etiqueta, antes: avisoA(params[clave].value ?? 0), despues: avisoA(v) }], nota: <>{total === undefined ? T.comun.pendiente : t.umbrales(dinero(total), avisoA(w), avisoA(l))}{ejemplo && <span className="block">{lineaCategoria(w, l)}</span>}</> }
  })
  if (esUmbral) return <>{umbral('alerts.budget_warning_pct', t.primerAviso)}{umbral('alerts.budget_limit_pct', t.tope)}<p className="text-caption text-text-secondary">{t.reglaPresupuesto}</p>
    <label className="flex flex-col gap-1 text-small text-text-secondary">{t.categoriaEjemplo}
      <select className={estilosControl} value={cat} onChange={(e) => setCat(e.target.value)}><option value="">{t.elegirCategoria}</option>{(porCategoria.datos?.rows ?? []).filter((r) => (r.planned_cents ?? 0) > 0).map((r) => <option key={r.category_id} value={String(r.category_id)}>{nombreCat(r.category_id)}</option>)}</select>
    </label>
    {ejemplo && <p className="text-small tabular-nums text-text-secondary" data-prueba="umbral-categoria">{lineaCategoria(params['alerts.budget_warning_pct'].value ?? 0, params['alerts.budget_limit_pct'].value ?? 0)}</p>}</>
  if (n.code === 'payment_due') return p<number[]>('alerts.payment_due_days_before', t.diasPago, 'lista', T.validar.diasPago, undefined, (v) => {
    const fechas = (l: number[]) => (vence ? l.map((d) => `el ${sinAnio(sumarDias(vence, -d))}`).join(' y ') : '—')
    return { cambios: [{ etiqueta: t.diasPago, antes: fechas(params['alerts.payment_due_days_before'].value ?? []), despues: fechas(v) }], nota: vence ? t.proximoPago(fechaConDia(vence), fechas(v)) : t.sinPago }
  }, t.diasPagoAyuda)
  if (n.code === 'stale_prices') return p<number>('alerts.stale_prices_advance_days', t.diasDespensa, 'entero', T.validar.diasDespensa, 'días', (v) => ({ cambios: [], nota: t.llegara(v) }))
  return p<number>('alerts.meal_plan_missing_weekday', t.diaRevision, 'diaSemana', T.validar.diaSemana, undefined, (v) => ({ cambios: [], nota: t.seRevisa(DIAS_SEMANA[v - 1]) }))
}
function resumen({ n, params }: Ctx): string {
  if (n.code === 'payment_due') return t.resumenDias(params['alerts.payment_due_days_before'].value ?? [])
  if (n.code === 'budget_threshold') return `${params['alerts.budget_warning_pct'].value} % y ${params['alerts.budget_limit_pct'].value} %`
  if (n.code === 'stale_prices') return t.llegara(params['alerts.stale_prices_advance_days'].value ?? 0)
  return t.seRevisa(DIAS_SEMANA[(params['alerts.meal_plan_missing_weekday'].value ?? 7) - 1])
}
/** `Ctrl+/` (U4 P-26.09, Atajos; H-08): activa o desactiva a la vez las alertas seleccionadas. Antes de aplicar, la vista previa
 *  sumada dice cuántas alertas quedan activas y cuáles cambian; `Enter` aplica a todas (un `updateNotification` por alerta) y
 *  `Esc` cancela sin cambiar ninguna. */
function grupoAlertas(alertas: N[], recargar: () => void): EditarVarias<N> {
  const estado = (v: boolean) => (v ? t.activa : t.inactiva)
  return {
    propiedades: [{ id: 'is_enabled', nombre: t.activa, tipo: 'opcion', opciones: [{ valor: 'si', texto: T.comun.si }, { valor: 'no', texto: T.comun.no }] }],
    vistaPrevia: async (sel, _p, v) => {
      const activar = v === 'si'
      const cambian = sel.filter((n) => n.is_enabled !== activar)
      const antes = alertas.filter((n) => n.is_enabled).length
      return [
        { etiqueta: t.alertasActivas, antes: t.deTotal(antes, alertas.length), despues: t.deTotal(antes + (activar ? cambian.length : -cambian.length), alertas.length) },
        ...cambian.map((n) => ({ etiqueta: n.name, antes: estado(n.is_enabled), despues: estado(activar) })),
      ]
    },
    guardar: async (sel, _p, v) => {
      try { for (const n of sel) await pedir('patch', '/api/v1/notifications/{notification_id}', { ruta: { notification_id: n.id }, cuerpo: { is_enabled: v === 'si' } }) } finally { recargar() }
    },
  }
}
function Alerta({ casilla, ...c }: Ctx & { casilla?: ReactNode }) {
  const g = useGuardado()
  const sinRed = useSinRed()
  const escritorio = useEsEscritorio()
  const [abierta, setAbierta] = useState(false)
  const activa = (v: boolean) => g.ejecutar(async () => { await pedir('patch', '/api/v1/notifications/{notification_id}', { ruta: { notification_id: c.n.id }, cuerpo: { is_enabled: v } }); c.recargar() })
  const cuerpo = <><Destinatarios {...c} /><Parametros {...c} /></>
  return (
    <article className="flex flex-col gap-2 rounded-md border border-border-divider p-3" data-prueba={`alerta-${c.n.code}`}>
      <div className="flex flex-wrap items-center justify-between gap-2">
        <span className="flex items-center gap-1">{casilla}<h4 className="text-body font-semibold">{c.n.name}</h4></span>
        <span className="flex items-center gap-2"><Interruptor ocultarEtiqueta etiqueta={`${t.activa} · ${c.n.name}`} activo={c.n.is_enabled} deshabilitadoPor={sinRed} alCambiar={(v) => void activa(v)} /><IndicadorGuardado g={g} /></span>
      </div>
      <p className="text-small text-text-secondary">{resumen(c)}</p>
      {escritorio ? cuerpo : (
        <>
          <Boton tam="hoja" onClick={() => setAbierta(true)}>{T.comun.editar}</Boton>
          <Hoja abierta={abierta} alCambiar={setAbierta} titulo={c.n.name}><div className="flex flex-col gap-2">{cuerpo}</div></Hoja>
        </>
      )}
    </article>
  )
}
function Confirmados({ recips, personas, recargar }: { recips: R[]; personas: Pe[]; recargar: () => void }) {
  return (
    <Bloque titulo={t.confirmados}>
      <div className="flex flex-col gap-3">
        {personas.filter((p) => p.participation !== 'external').map((p) => {
          const r = recips.find((x) => x.person_id === p.id && x.verified_at)
          return (
            <div key={p.id} className="flex flex-col gap-2 border-b border-border-divider pb-3 last:border-b-0 tablet:flex-row tablet:items-center" data-prueba={`chat-${p.code}`}>
              <span className="text-body font-medium tablet:w-40">{p.display_name}</span>
              {r ? <ChatConfirmado r={r} nombre={p.display_name} recargar={recargar} /> : <span className="text-small text-text-secondary">{t.sinChat(p.display_name)}</span>}
            </div>
          )
        })}
      </div>
    </Bloque>
  )
}
function ChatConfirmado({ r, nombre, recargar }: { r: R; nombre: string; recargar: () => void }) {
  const g = useGuardado()
  const sinRed = useSinRed()
  const [etiqueta, setEtiqueta] = useState(r.label ?? '')
  const guardar = (cuerpo: Esquemas['NotificationRecipientUpdate']) => g.ejecutar(async () => { await pedir('patch', '/api/v1/notification-recipients/{recipient_id}', { ruta: { recipient_id: r.id }, cuerpo }); recargar() })
  return (
    <div className="flex flex-1 flex-col gap-2 tablet:flex-row tablet:flex-wrap tablet:items-center">
      <span className="text-small">{`${t.registrado}: ${t.confirmadoEl(fechaHora(r.verified_at!))}`}</span>
      <label className="flex items-center gap-2 text-small">{t.etiqueta}
        <input className={cn(estilosControl, 'w-40')} value={etiqueta} disabled={!!sinRed} onChange={(e) => setEtiqueta(e.target.value)}
          onBlur={() => etiqueta !== (r.label ?? '') && void guardar({ label: etiqueta.trim() || null })}
          onKeyDown={(e) => { if (e.key === 'Enter') e.currentTarget.blur(); else if (e.key === 'Escape') setEtiqueta(r.label ?? '') }} />
      </label>
      <Interruptor etiqueta={`${t.recibe} · ${nombre}`} activo={r.is_active} deshabilitadoPor={sinRed} alCambiar={(v) => void guardar({ is_active: v })} />
      <IndicadorGuardado g={g} />
    </div>
  )
}
function Detectados({ cs, personas, recargar }: { cs: C[]; personas: Pe[]; recargar: () => void }) {
  return (
    <Bloque titulo={t.detectados} ayuda={t.ayudaBot} accion={<Boton variante="fantasma" onClick={recargar}>{t.actualizar}</Boton>}>
      {cs.length === 0 ? <Vacio texto={t.sinDetectados} /> : <div className="flex flex-col gap-3">{cs.map((c) => <Candidato key={c.id} c={c} personas={personas} recargar={recargar} />)}</div>}
    </Bloque>
  )
}
function Candidato({ c, personas, recargar }: { c: C; personas: Pe[]; recargar: () => void }) {
  const { sesion } = useSesion()
  const avisar = useAvisos()
  const sinRed = useSinRed()
  const g = useGuardado()
  const conCuenta = personas.filter((p) => p.participation !== 'external')
  const [persona, setPersona] = useState(String(sesion?.person_id ?? conCuenta[0]?.id ?? ''))
  const [dialogo, setDialogo] = useState(false)
  const clave = useRef(nuevaClave())
  const fin = `•••${(c.chat_id_masked ?? '').replace(/^[•*.\s]+/, '').slice(-4)}`
  const confirmar = () => g.ejecutar(async () => {
    try {
      await pedir('post', '/api/v1/telegram-chat-candidates/{candidate_id}/confirm', { ruta: { candidate_id: c.id }, cuerpo: { person_id: Number(persona) }, idempotencia: clave.current })
      avisar({ texto: t.chatConfirmado, rol: 'exito' })
    } catch (e) {
      if (e instanceof ErrorApi && e.estado === 409) { avisar({ texto: t.yaConfirmado, rol: 'advertencia' }); return }
      throw e
    } finally { recargar() }
  })
  const descartar = () => g.ejecutar(async () => { setDialogo(false); await pedir('delete', '/api/v1/telegram-chat-candidates/{candidate_id}', { ruta: { candidate_id: c.id } }); avisar({ texto: t.descartado }); recargar() })
  return (
    <div className="flex flex-col gap-2 rounded-md border border-border-divider p-3 tablet:flex-row tablet:flex-wrap tablet:items-center" data-prueba={`candidato-${c.id}`}>
      <div className="flex flex-col tablet:w-56"><span className="text-body font-medium">{c.telegram_name}</span><span className="text-caption text-text-secondary">{`${t.detectado}: ${fechaHora(c.detected_at)} · ${t.terminacion}: ${fin}`}</span></div>
      <label className="flex items-center gap-2 text-small">{t.confirmarComo}
        <select className={cn(estilosControl, 'w-36')} value={persona} onChange={(e) => setPersona(e.target.value)}>{conCuenta.map((p) => <option key={p.id} value={p.id}>{p.display_name}</option>)}</select>
      </label>
      <Boton variante="primario" tam="hoja" className="tablet:w-auto" deshabilitadoPor={sinRed} onClick={() => void confirmar()}>{t.confirmar}</Boton>
      <Boton variante="fantasma" deshabilitadoPor={sinRed} onClick={() => setDialogo(true)}>{t.descartar}</Boton>
      <IndicadorGuardado g={g} />
      <DialogoCritico abierto={dialogo} alCambiar={setDialogo} titulo={t.descartarTitulo} texto={t.descartarTexto(c.telegram_name, fin)} verbo={t.descartarTitulo} alConfirmar={() => void descartar()} />
    </div>
  )
}
function PruebaAviso({ alertas, personas, recips }: { alertas: N[]; personas: Pe[]; recips: R[] }) {
  const avisar = useAvisos()
  const sinRed = useSinRed()
  const [abierta, setAbierta] = useState(false)
  const confirmadas = personas.filter((p) => recips.some((r) => r.person_id === p.id && r.verified_at))
  const [alerta, setAlerta] = useState(String(alertas[0]?.id ?? ''))
  const [persona, setPersona] = useState(String(confirmadas[0]?.id ?? ''))
  const [error, setError] = useState('')
  const enviar = async () => {
    setError('')
    try {
      await pedir('post', '/api/v1/notifications/{notification_id}/test', { ruta: { notification_id: Number(alerta) }, cuerpo: { person_id: Number(persona) } })
      setAbierta(false)
      avisar({ texto: t.pruebaEnviada(personas.find((p) => p.id === Number(persona))?.display_name ?? ''), rol: 'exito' })
    } catch { setError(t.pruebaFallo) }
  }
  return (
    <>
      <Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => setAbierta(true)}>{t.prueba}</Boton>
      <Hoja abierta={abierta} alCambiar={setAbierta} titulo={t.prueba} pie={<Boton variante="primario" tam="hoja" deshabilitadoPor={sinRed ?? (persona ? undefined : T.r.soloConfirmados)} onClick={() => void enviar()}>{t.enviar}</Boton>}>
        <div className="flex flex-col gap-3">
          <label className="flex flex-col gap-1 text-small font-medium">{t.alerta}<select className={estilosControl} value={alerta} onChange={(e) => setAlerta(e.target.value)}>{alertas.map((a) => <option key={a.id} value={a.id}>{a.name}</option>)}</select></label>
          <label className="flex flex-col gap-1 text-small font-medium">{t.persona}<select className={estilosControl} value={persona} onChange={(e) => setPersona(e.target.value)}>{confirmadas.map((p) => <option key={p.id} value={p.id}>{p.display_name}</option>)}</select></label>
          {error && <p role="alert" className="text-small text-role-error">{error}</p>}
        </div>
      </Hoja>
    </>
  )
}
