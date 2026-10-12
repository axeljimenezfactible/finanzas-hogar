// P-26.13 Salud y metas de ahorro (U4 parte b; decisión D, K7, K23, K36, RN-17, RN-18, Q2-16 A, Q2-17 C, Q2-18 A).
// Vistas previas en el navegador; al guardar manda el servidor. Aportación de una meta vencida sin cumplir:
// «Vista previa pendiente» (PROVISIONAL, bloqueo S2-metas-vencida).
import { useEffect, useRef, useState, type ReactNode } from 'react'
import { pedir, nuevaClave } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { aCentavos, dinero, fecha, hoyISO, kilos, mes } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { Campo, estilosControl } from '../../componentes/Campo'
import { Hoja } from '../../componentes/Hoja'
import { Esqueleto, Vacio } from '../../componentes/Estados'
import { cn } from '../../componentes/utils'
import { abrirRegistro } from '../../shell/estado'
import { VistaBD } from '../../vistas-bd/VistaBD'
import type { Cambio } from '../../vistas-bd/tipos'
import { AntesDespues, Bloque, Cargar, FilaParametro, IndicadorGuardado, Marco, unir, useGuardado, useSinRed } from './comun'
import { invalido } from './P2612Personas'
import { aportacionNecesaria, avanceDecimas, claseImc, decimas, estadoMeta, imcDecimas, metaSugerida, numero, type InsumosSalud } from './calculos'
import { textosAjustesB as T } from './textos'

const t = T.s
type Pe = Esquemas['Person']; type HP = Esquemas['HealthProfile']; type G = Esquemas['SavingsGoal']
type HS = Esquemas['HealthSettings']['parameters']
const guardarCasa = (c: Esquemas['HealthSettingsUpdate']) => pedir('patch', '/api/v1/settings/health', { cuerpo: c })

export default function P2613Salud() {
  const datos = unir<[{ items: Pe[] }, { items: HP[] }, Esquemas['HealthSettings'], { items: G[] }]>([
    usarConsulta('get', '/api/v1/persons', { consulta: { limit: 50 } }),
    usarConsulta('get', '/api/v1/health-profiles', { consulta: { limit: 50 } }),
    usarConsulta('get', '/api/v1/settings/health'),
    usarConsulta('get', '/api/v1/savings-goals', { consulta: { limit: 200 } }),
  ])
  useEffect(() => { const f = () => datos.recargar(); window.addEventListener('fh:datos-cambiaron', f); return () => window.removeEventListener('fh:datos-cambiaron', f) }, [datos.recargar]) // eslint-disable-line react-hooks/exhaustive-deps
  return (
    <Marco id="13">
      <p className="text-body text-text-secondary">{t.proposito}</p>
      <Cargar c={datos} que={t.que} esqueleto={<div className="grid gap-3 desktop:grid-cols-2"><Esqueleto className="h-64" /><Esqueleto className="h-64" /></div>}>{([ps, hp, hs, gs]) => {
        const titulares = ps.items.filter((p) => p.participation !== 'external')
        return (
          <>
            <div className="grid grid-cols-1 gap-4 desktop:grid-cols-2">
              {titulares.map((p) => <Perfil key={p.id} p={p} perfil={hp.items.find((x) => x.person_id === p.id)} casa={hs.parameters} recargar={datos.recargar} />)}
            </div>
            <Casa casa={hs.parameters} perfiles={hp.items} personas={titulares} recargar={datos.recargar} />
            <Metas metas={gs.items} recargar={datos.recargar} />
          </>
        )
      }}</Cargar>
    </Marco>
  )
}

const insumos = (p: HP, casa: HS): InsumosSalud => ({ sex: p.sex, birth_date: p.birth_date, height_mm: p.height_mm, weight_g: p.current_weight_g, activity_level: p.activity_level, goal: p.goal,
  ajusteBajar: p.goal_adjust_lose_kcal ?? p.household_goal_adjust_lose_kcal ?? casa['health.goal_adjust_lose_kcal'].value ?? 0,
  ajusteSubir: p.goal_adjust_gain_kcal ?? p.household_goal_adjust_gain_kcal ?? casa['health.goal_adjust_gain_kcal'].value ?? 0 })
const textoMeta = (x: InsumosSalud) => { const m = metaSugerida(x); return 'falta' in m ? t.falta(t.faltas[m.falta]) : m.desglose }
const textoImc = (g?: number | null, mm?: number | null) => (!mm ? t.falta(t.faltas.estatura) : !g ? t.falta(t.faltas.peso) : `${decimas(imcDecimas(g, mm))} (${t.clases[claseImc(imcDecimas(g, mm))]})`)
const aMm = (s: string) => { const x = s.trim(); if (!/^\d+(\.\d{1,3})?$/.test(x)) return null; const [e, d = ''] = x.split('.'); const mm = Number(e) * 1000 + Number(d.padEnd(3, '0')); return mm > 0 ? mm : null }

type CampoP = { etiqueta: string; clave: keyof Esquemas['HealthProfileUpdate']; inicial: string; tipo?: 'texto' | 'fecha' | 'opcion'; opciones?: [string, string][]; interpretar: (s: string) => unknown; mensaje?: string; previa?: (s: string) => Cambio[]; extra?: ReactNode; atenuado?: boolean; resaltado?: boolean }
function CampoPerfil({ c, perfil, recargar }: { c: CampoP; perfil: HP; recargar: () => void }) {
  const sinRed = useSinRed()
  const g = useGuardado()
  const [v, setV] = useState(c.inicial)
  const [error, setError] = useState('')
  useEffect(() => setV(c.inicial), [c.inicial])
  const id = `perfil-${perfil.person_id}-${c.clave}`
  const guardar = async () => {
    if (v === c.inicial) return
    const x = c.interpretar(v)
    if (x === undefined) { setError(c.mensaje ?? ''); setV(c.inicial); return }
    setError('')
    if (!(await g.ejecutar(async () => { await pedir('patch', '/api/v1/health-profiles/{health_profile_id}', { ruta: { health_profile_id: perfil.id }, cuerpo: { [c.clave]: x } }); recargar() }))) setV(c.inicial)
  }
  const comun = { id, disabled: !!sinRed, 'aria-invalid': error ? true : undefined, onBlur: () => void guardar(), onKeyDown: (e: React.KeyboardEvent) => { if (e.key === 'Enter') { e.preventDefault(); void guardar() } else if (e.key === 'Escape') { setV(c.inicial); setError('') } } }
  return (
    <div className={cn('flex flex-col gap-1', c.atenuado && 'opacity-60', c.resaltado && 'order-first rounded-md bg-surface-selected p-2')}>
      <label htmlFor={id} className="text-small font-medium">{c.etiqueta}</label>
      {c.tipo === 'opcion'
        ? <select {...comun} className={estilosControl} value={v} onChange={(e) => setV(e.target.value)}>{c.opciones!.map(([k, n]) => <option key={k} value={k}>{n}</option>)}</select>
        : <input {...comun} type={c.tipo === 'fecha' ? 'date' : 'text'} className={cn(estilosControl, 'tabular-nums')} value={v} onChange={(e) => setV(e.target.value)} />}
      {c.extra}
      {error ? <p className="text-caption text-role-error">{error}</p> : <IndicadorGuardado g={g} />}
      {v !== c.inicial && c.previa && c.interpretar(v) !== undefined && <AntesDespues cambios={c.previa(v)} />}
    </div>
  )
}
function Perfil({ p, perfil, casa, recargar }: { p: Pe; perfil?: HP; casa: HS; recargar: () => void }) {
  const sinRed = useSinRed()
  const g = useGuardado()
  const clave = useRef(nuevaClave())
  if (!perfil) return (
    <Bloque titulo={t.perfil(p.display_name)}>
      <Vacio texto={t.sinPerfil} accion={<Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => void g.ejecutar(async () => { await pedir('post', '/api/v1/health-profiles', { cuerpo: { person_id: p.id }, idempotencia: clave.current }); recargar() })}>{t.crearPerfil}</Boton>} />
      <IndicadorGuardado g={g} />
    </Bloque>
  )
  const base = insumos(perfil, casa)
  const previaCon = (cambio: Partial<InsumosSalud>, imc = false): Cambio[] => {
    const nuevo = { ...base, ...cambio }
    const l = [{ etiqueta: t.sugerida, antes: textoMeta(base), despues: textoMeta(nuevo) }]
    return imc ? [{ etiqueta: t.imc, antes: textoImc(base.weight_g, base.height_mm), despues: textoImc(nuevo.weight_g, nuevo.height_mm) }, ...l] : l
  }
  const entero = (s: string) => (s.trim() === '' ? null : /^\d+$/.test(s.trim()) ? Number(s) : undefined)
  const meta = metaSugerida(base)
  const mantener = perfil.goal === 'maintain'
  const ajuste = (k: 'goal_adjust_lose_kcal' | 'goal_adjust_gain_kcal', etiqueta: string, casaN: number, objetivo: string): CampoP => ({
    etiqueta, clave: k, inicial: perfil[k] === undefined || perfil[k] === null ? '' : String(perfil[k]), interpretar: entero, mensaje: T.validar.kcal, atenuado: mantener, resaltado: perfil.goal === objetivo,
    previa: (s) => previaCon(k === 'goal_adjust_lose_kcal' ? { ajusteBajar: entero(s) ?? casaN } : { ajusteSubir: entero(s) ?? casaN }),
    extra: <span className="flex flex-wrap items-center gap-2 text-caption text-text-secondary">{mantener ? t.sinAjuste : t.elDeLaCasa(casaN)}
      {perfil[k] !== undefined && perfil[k] !== null && <Boton variante="enlace" deshabilitadoPor={sinRed} onClick={() => void pedir('patch', '/api/v1/health-profiles/{health_profile_id}', { ruta: { health_profile_id: perfil.id }, cuerpo: { [k]: null } }).then(recargar)}>{t.usarCasa}</Boton>}</span>,
  })
  const campos: CampoP[] = [
    { etiqueta: t.estatura, clave: 'height_mm', inicial: perfil.height_mm ? (perfil.height_mm / 1000).toFixed(2) : '', interpretar: (s) => aMm(s) ?? undefined, mensaje: t.estaturaMal, previa: (s) => previaCon({ height_mm: aMm(s) }, true) },
    { etiqueta: t.nacimiento, clave: 'birth_date', tipo: 'fecha', inicial: perfil.birth_date ?? '', interpretar: (s) => (s && s <= hoyISO() ? s : undefined), mensaje: t.futura, previa: (s) => previaCon({ birth_date: s }),
      extra: perfil.birth_date && <span className="text-caption text-text-secondary">{fecha(perfil.birth_date)}{perfil.age_years !== undefined && ` · ${perfil.age_years} años`}</span> },
    { etiqueta: t.sexo, clave: 'sex', tipo: 'opcion', opciones: [['', t.sinElegir], ['female', t.sexos.female], ['male', t.sexos.male]], inicial: perfil.sex ?? '', interpretar: (s) => s || null, previa: (s) => previaCon({ sex: s || null }) },
    { etiqueta: t.actividad, clave: 'activity_level', tipo: 'opcion', opciones: [['', t.sinElegir], ...Object.entries(t.actividades)], inicial: perfil.activity_level ?? '', interpretar: (s) => s || null, previa: (s) => previaCon({ activity_level: s || null }) },
    { etiqueta: t.objetivo, clave: 'goal', tipo: 'opcion', opciones: [['', t.sinElegir], ...Object.entries(t.objetivos)], inicial: perfil.goal ?? '', interpretar: (s) => s || null, previa: (s) => previaCon({ goal: s || null }) },
    ajuste('goal_adjust_lose_kcal', t.ajusteBajar, perfil.household_goal_adjust_lose_kcal ?? casa['health.goal_adjust_lose_kcal'].value ?? 0, 'lose'),
    ajuste('goal_adjust_gain_kcal', t.ajusteSubir, perfil.household_goal_adjust_gain_kcal ?? casa['health.goal_adjust_gain_kcal'].value ?? 0, 'gain'),
    { etiqueta: t.confirmada, clave: 'daily_kcal_target', inicial: perfil.daily_kcal_target ? String(perfil.daily_kcal_target) : '', interpretar: (s) => (s.trim() === '' ? null : /^\d+$/.test(s.trim()) && Number(s) > 0 ? Number(s) : undefined), mensaje: t.kcalMal,
      extra: perfil.suggested_kcal_target !== undefined && <Boton variante="enlace" deshabilitadoPor={sinRed} onClick={() => void pedir('patch', '/api/v1/health-profiles/{health_profile_id}', { ruta: { health_profile_id: perfil.id }, cuerpo: { daily_kcal_target: perfil.suggested_kcal_target } }).then(recargar)}>{t.usarSugerida}</Boton> },
    { etiqueta: t.restricciones, clave: 'dietary_restrictions', inicial: perfil.dietary_restrictions ?? '', interpretar: (s) => s.trim() || null },
    { etiqueta: t.notas, clave: 'notes', inicial: perfil.notes ?? '', interpretar: (s) => s.trim() || null },
  ]
  return (
    <Bloque titulo={t.perfil(p.display_name)}>
      <dl className="grid grid-cols-2 gap-2 rounded-md bg-surface-raised p-3 text-small" data-prueba={`perfil-${p.code}`}>
        <dt className="text-text-secondary">{t.peso}</dt>
        <dd className="tabular-nums">{perfil.current_weight_g ? `${kilos(perfil.current_weight_g)} · ${fecha(perfil.current_weight_date!)}` : <Boton variante="enlace" onClick={() => abrirRegistro('peso')}>{`${t.sinPeso} · ${t.registrarPeso}`}</Boton>}</dd>
        <dt className="text-text-secondary">{t.imc}</dt>
        <dd className="tabular-nums">{perfil.bmi_tenths !== undefined && perfil.bmi_class ? `${decimas(perfil.bmi_tenths)} (${t.clases[perfil.bmi_class]})` : textoImc(perfil.current_weight_g, perfil.height_mm)}</dd>
        <dt className="text-text-secondary">{t.sugerida}</dt>
        <dd className="tabular-nums">{perfil.suggested_kcal_target !== undefined ? t.kcalDia(numero(perfil.suggested_kcal_target)) : '—'}<span className="block text-caption text-text-secondary">{'falta' in meta ? t.falta(t.faltas[meta.falta]) : meta.desglose}</span></dd>
      </dl>
      <div className="flex flex-col gap-3">{campos.map((c) => <CampoPerfil key={c.clave} c={c} perfil={perfil} recargar={recargar} />)}</div>
    </Bloque>
  )
}
function Casa({ casa, perfiles, personas, recargar }: { casa: HS; perfiles: HP[]; personas: Pe[]; recargar: () => void }) {
  const usan = (k: 'goal_adjust_lose_kcal' | 'goal_adjust_gain_kcal') => perfiles.filter((p) => p[k] === undefined || p[k] === null)
  const nombres = personas.filter((p) => usan('goal_adjust_lose_kcal').some((x) => x.person_id === p.id)).map((p) => p.display_name).join(' y ')
  const fila = (clave: 'health.goal_adjust_lose_kcal' | 'health.goal_adjust_gain_kcal', etiqueta: string) => (
    <FilaParametro etiqueta={etiqueta} tipo="entero" unidad="kcal" param={casa[clave]} mensaje={T.validar.kcal}
      vistaPrevia={(v: number) => ({ cambios: usan(clave === 'health.goal_adjust_lose_kcal' ? 'goal_adjust_lose_kcal' : 'goal_adjust_gain_kcal').map((p) => {
        const b = insumos({ ...p, household_goal_adjust_lose_kcal: undefined, household_goal_adjust_gain_kcal: undefined }, casa)
        const n = clave === 'health.goal_adjust_lose_kcal' ? { ...b, ajusteBajar: v } : { ...b, ajusteSubir: v }
        const quien = personas.find((x) => x.id === p.person_id)?.display_name ?? ''
        return { etiqueta: `${t.sugerida} · ${quien}`, antes: textoMeta(b), despues: textoMeta(n) }
      }) })}
      guardar={async (v: number) => { await guardarCasa({ parameters: { [clave]: v } }); recargar() }} restablecer={async () => { await guardarCasa({ reset_parameters: [clave] }); recargar() }} />
  )
  return (
    <Bloque titulo={t.ajusteCasa} ayuda={t.usanCasa(nombres)}>
      {fila('health.goal_adjust_lose_kcal', t.casaBajar)}
      {fila('health.goal_adjust_gain_kcal', t.casaSubir)}
    </Bloque>
  )
}

const actualizarMeta = (id: number, cuerpo: Esquemas['SavingsGoalUpdate']) => pedir('patch', '/api/v1/savings-goals/{savings_goal_id}', { ruta: { savings_goal_id: id }, cuerpo })
const necesaria = (m: G, objetivo = m.target_cents, fechaMeta = m.target_date) => {
  if (m.kind !== 'goal' || !objetivo || !fechaMeta) return '—'
  const n = aportacionNecesaria(m.saved_cents ?? 0, objetivo, fechaMeta)
  return n === null ? T.comun.pendiente : dinero(n)
}
function previaMeta(m: G, objetivo = m.target_cents, fechaMeta = m.target_date): Cambio[] {
  const ah = m.saved_cents ?? 0
  const l: Cambio[] = []
  if (objetivo !== m.target_cents && objetivo) l.push({ etiqueta: t.avance, antes: m.target_cents ? t.efectoAvance(dinero(ah), dinero(m.target_cents), decimas(avanceDecimas(ah, m.target_cents))) : '—', despues: t.efectoAvance(dinero(ah), dinero(objetivo), decimas(avanceDecimas(ah, objetivo))) })
  const nueva = necesaria(m, objetivo, fechaMeta)
  l.push({ etiqueta: t.efectoNecesaria, antes: necesaria(m), despues: nueva === T.comun.pendiente ? t.vencidaPendiente : nueva })
  l.push({ etiqueta: t.estado, antes: t.estados[estadoMeta(ah, m.target_cents, m.target_date)], despues: t.estados[estadoMeta(ah, objetivo, fechaMeta)] })
  return l
}
function Metas({ metas, recargar }: { metas: G[]; recargar: () => void }) {
  const sinRed = useSinRed()
  const [nueva, setNueva] = useState(false)
  useEffect(() => { const f = (e: KeyboardEvent) => { if (e.altKey && e.key.toLowerCase() === 'n') { e.preventDefault(); setNueva(true) } }; window.addEventListener('keydown', f); return () => window.removeEventListener('keydown', f) }, [])
  const reservado = (l: G[]) => l.filter((m) => m.is_active).reduce((s, m) => s + (m.planned_monthly_cents ?? 0), 0)
  const mesHoy = mes(hoyISO().slice(0, 7))
  const filas = [...metas].sort((a, b) => (a.target_date ?? '9999').localeCompare(b.target_date ?? '9999'))
  const boton = <Boton variante="primario" deshabilitadoPor={sinRed} onClick={() => setNueva(true)}>{t.nuevaMeta}</Boton>
  return (
    <Bloque titulo={t.metas} accion={boton}>
      <VistaBD def={{
        clave: 'P-26.13-metas', filas, idFila: (m) => m.id, vistas: ['tabla', 'tablero', 'calendario', 'galeria'], fecha: (m) => m.target_date, titulo: (m) => m.name,
        columnasTablero: { propiedad: 'status', columnas: [t.estados.in_progress, t.estados.achieved, t.estados.overdue] },
        vacio: <Vacio texto={t.sinMetas} accion={boton} />,
        propiedades: [
          { id: 'name', nombre: t.nombre, valor: (m) => m.name, editar: { tipo: 'texto', textoInicial: (m) => m.name, guardar: async (m, v) => { if (!v.trim()) invalido(t.nombreMal); await actualizarMeta(m.id, { name: v.trim() }); recargar() } } },
          { id: 'kind', nombre: t.tipo, valor: (m) => t.tipos[m.kind], editar: { tipo: 'opcion', opciones: [{ valor: 'fund', texto: t.tipos.fund }, { valor: 'goal', texto: t.tipos.goal }], textoInicial: (m) => m.kind, guardar: async (m, v) => { await actualizarMeta(m.id, { kind: v as G['kind'] }); recargar() } } },
          { id: 'target_cents', nombre: t.monto, alinear: 'derecha', valor: (m) => m.target_cents ?? null, pintar: (m) => <span className="tabular-nums">{m.target_cents ? dinero(m.target_cents) : '—'}</span>,
            editar: { tipo: 'dinero', textoInicial: (m) => (m.target_cents ? (m.target_cents / 100).toFixed(2) : ''), vistaPrevia: async (m, v) => { const c = aCentavos(v); if (!c || c <= 0) throw new Error(t.montoMal); return previaMeta(m, c) },
              guardar: async (m, v) => { const c = aCentavos(v); if (!c || c <= 0) invalido(t.montoMal); await actualizarMeta(m.id, { target_cents: c }); recargar() } } },
          { id: 'target_date', nombre: t.fechaMeta, valor: (m) => m.target_date ?? null, pintar: (m) => (m.target_date ? fecha(m.target_date) : '—'),
            editar: { tipo: 'fecha', textoInicial: (m) => m.target_date ?? '', vistaPrevia: async (m, v) => previaMeta(m, m.target_cents, v), guardar: async (m, v) => { if (v && v < m.start_date) invalido(t.fechaMal); await actualizarMeta(m.id, { target_date: v || null }); recargar() } } },
          { id: 'planned_monthly_cents', nombre: t.aportacion, alinear: 'derecha', valor: (m) => m.planned_monthly_cents ?? null, pintar: (m) => <span className="tabular-nums">{m.planned_monthly_cents !== undefined ? dinero(m.planned_monthly_cents) : '—'}</span>,
            editar: { tipo: 'dinero', textoInicial: (m) => (m.planned_monthly_cents !== undefined ? (m.planned_monthly_cents / 100).toFixed(2) : ''),
              vistaPrevia: async (m, v) => { const c = v.trim() === '' ? 0 : aCentavos(v) ?? 0; const otros = reservado(metas.filter((x) => x.id !== m.id)); return [{ etiqueta: t.efectoReserva(mesHoy), antes: dinero(reservado(metas)), despues: dinero(otros + (m.is_active ? c : 0)) }] },
              guardar: async (m, v) => { const c = v.trim() === '' ? null : aCentavos(v); if (v.trim() && c === null) invalido(t.aportMal); await actualizarMeta(m.id, { planned_monthly_cents: c }); recargar() } } },
          { id: 'start_date', nombre: t.inicio, valor: (m) => m.start_date, pintar: (m) => fecha(m.start_date), editar: { tipo: 'fecha', textoInicial: (m) => m.start_date, guardar: async (m, v) => { if (!v) invalido(t.inicioMal); await actualizarMeta(m.id, { start_date: v }); recargar() } } },
          { id: 'is_active', nombre: t.activa, valor: (m) => (m.is_active ? T.comun.si : T.comun.no),
            editar: { tipo: 'opcion', opciones: [{ valor: 'si', texto: T.comun.si }, { valor: 'no', texto: T.comun.no }], textoInicial: (m) => (m.is_active ? 'si' : 'no'),
              vistaPrevia: async (m, v) => { const otros = metas.filter((x) => x.id !== m.id); const d = reservado([...otros, { ...m, is_active: v === 'si' }]); return [{ etiqueta: t.efectoReserva(mesHoy), antes: dinero(reservado(metas)), despues: m.planned_monthly_cents ? dinero(d) : `${dinero(d)} · ${t.sinAportacion}` }] },
              guardar: async (m, v) => { await actualizarMeta(m.id, { is_active: v === 'si' }); recargar() } } },
          { id: 'saved_cents', nombre: t.ahorrado, alinear: 'derecha', valor: (m) => m.saved_cents ?? 0, pintar: (m) => <span className="tabular-nums">{dinero(m.saved_cents ?? 0)}</span> },
          { id: 'progress', nombre: t.avance, alinear: 'derecha', valor: (m) => m.progress_tenths_pct ?? null, pintar: (m) => (m.progress_tenths_pct !== undefined ? <Avance d={m.progress_tenths_pct} /> : '—') },
          { id: 'required', nombre: t.necesaria, alinear: 'derecha', valor: (m) => m.required_monthly_cents ?? null, pintar: (m) => <span className="tabular-nums">{m.required_monthly_cents !== undefined ? dinero(m.required_monthly_cents) : '—'}</span> },
          { id: 'status', nombre: t.estado, valor: (m) => (m.status ? t.estados[m.status] : null), grupo: (m) => (m.status ? t.estados[m.status] : t.estados.in_progress) },
          { id: 'notes', nombre: t.notas, valor: (m) => m.notes ?? '', pintar: (m) => m.notes || T.comun.sinValor, editar: { tipo: 'texto', textoInicial: (m) => m.notes ?? '', guardar: async (m, v) => { await actualizarMeta(m.id, { notes: v.trim() || null }); recargar() } } },
        ],
      }} />
      {nueva && <NuevaMeta cerrar={() => setNueva(false)} recargar={recargar} />}
    </Bloque>
  )
}
const Avance = ({ d }: { d: number }) => (
  <span className="inline-flex items-center gap-2 tabular-nums">
    <span aria-hidden className="h-2 w-16 overflow-hidden rounded-full bg-viz-track"><span className="block h-full bg-role-savings" style={{ width: `${Math.min(100, d / 10)}%` }} /></span>{`${decimas(d)} %`}
  </span>
)
function NuevaMeta({ cerrar, recargar }: { cerrar: () => void; recargar: () => void }) {
  const sinRed = useSinRed()
  const g = useGuardado()
  const clave = useRef(nuevaClave())
  const [f, setF] = useState({ name: '', kind: '', target: '', target_date: '', planned: '', start_date: '', notes: '' })
  const [errores, setErrores] = useState<Record<string, string>>({})
  const poner = (k: keyof typeof f) => (e: { target: { value: string } }) => setF({ ...f, [k]: e.target.value })
  const guardar = async () => {
    const e: Record<string, string> = {}
    if (!f.name.trim()) e.name = t.nombreMal
    if (f.kind !== 'fund' && f.kind !== 'goal') e.kind = t.tipoMal
    const target = f.target.trim() ? aCentavos(f.target) : null
    if (f.kind === 'goal' && (!target || target <= 0)) e.target = t.montoMal
    if (!f.start_date) e.start_date = t.inicioMal
    if (f.kind === 'goal' && (!f.target_date || f.target_date < f.start_date)) e.target_date = t.fechaMal
    const planned = f.planned.trim() ? aCentavos(f.planned) : null
    if (f.planned.trim() && planned === null) e.planned = t.aportMal
    setErrores(e)
    if (Object.keys(e).length) return
    const ok = await g.ejecutar(async () => {
      await pedir('post', '/api/v1/savings-goals', { idempotencia: clave.current, cuerpo: { name: f.name.trim(), kind: f.kind as G['kind'], start_date: f.start_date, ...(target ? { target_cents: target } : {}), ...(f.target_date ? { target_date: f.target_date } : {}), ...(planned !== null ? { planned_monthly_cents: planned } : {}), ...(f.notes.trim() ? { notes: f.notes.trim() } : {}) } })
      recargar()
    })
    if (ok) cerrar()
  }
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={t.nuevaMeta} pie={<><Boton variante="primario" tam="hoja" deshabilitadoPor={sinRed} onClick={() => void guardar()}>{T.comun.guardar}</Boton><IndicadorGuardado g={g} /></>}>
      <form className="flex flex-col gap-3" onKeyDown={(e) => { if (e.ctrlKey && e.key === 'Enter') { e.preventDefault(); void guardar() } }} onSubmit={(e) => { e.preventDefault(); void guardar() }}>
        <Campo etiqueta={t.nombre} value={f.name} onChange={poner('name')} error={errores.name} />
        <label className="flex flex-col gap-1 text-small font-medium">{t.tipo}
          <select className={estilosControl} value={f.kind} onChange={poner('kind')} aria-invalid={errores.kind ? true : undefined}><option value="">{t.tipoMal}</option><option value="fund">{t.tipos.fund}</option><option value="goal">{t.tipos.goal}</option></select>
          {errores.kind && <span className="text-caption text-role-error">{errores.kind}</span>}
        </label>
        <Campo etiqueta={t.monto} inputMode="decimal" value={f.target} onChange={poner('target')} error={errores.target} />
        <Campo etiqueta={t.inicio} type="date" value={f.start_date} onChange={poner('start_date')} error={errores.start_date} />
        <Campo etiqueta={t.fechaMeta} type="date" value={f.target_date} onChange={poner('target_date')} error={errores.target_date} />
        <Campo etiqueta={t.aportacion} inputMode="decimal" value={f.planned} onChange={poner('planned')} error={errores.planned} />
        <Campo etiqueta={t.notas} value={f.notes} onChange={poner('notes')} />
      </form>
    </Hoja>
  )
}
