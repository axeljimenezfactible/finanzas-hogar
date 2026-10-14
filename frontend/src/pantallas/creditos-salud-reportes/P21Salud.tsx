// P-21 Salud (U1 §2.1, §2.2, §3.23: N-230..N-233, N-242, N-327; K23, K36). Dos columnas (Axel y Ximena), kcal de hoy,
// historial de peso editable y tendencias. Las cifras (IMC, meta sugerida) las calcula el servidor (motor S4).
import { useEffect, useState } from 'react'
import { Scale, Trash2, UserPen } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { navegar } from '../../app/enrutador'
import { usarConsulta } from '../../api/consulta'
import { pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { aGramos, fecha, kilos } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ChipPersona, ORDEN_PERSONAS } from '../../componentes/Chip'
import { useAvisos } from '../../componentes/Avisos'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Grafica } from '../../componentes/Grafica'
import { Hoja } from '../../componentes/Hoja'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { abrirRegistro } from '../../shell/estado'
import { Barra, Datos, Seccion, alCambiarDatos, itemsDe, mensajeError, usarPersonas, type Personas } from './comun'
import { ACTIVIDAD, IMC, OBJETIVO, SEXO, t } from './textos'

const s = t.salud
type Perfil = Esquemas['HealthProfile']
type Peso = Esquemas['WeightEntry']
type Hoy = Esquemas['TodayKcal']
const imc = (d?: number | null) => (d == null ? t.comun.sinDato : (d / 10).toFixed(1))

export default function P21Salud({ sub }: PropsPantalla) {
  const per = usarPersonas()
  const perfiles = usarConsulta('get', '/api/v1/health-profiles')
  const hoy = usarConsulta('get', '/api/v1/health/today-kcal', undefined, sub === 'resumen')
  const pesos = usarConsulta('get', '/api/v1/weight-entries', { consulta: { sort: '-measured_date', limit: 200 } })
  const kcal = usarConsulta('get', '/api/v1/meal-plan/trends', undefined, sub === 'tendencias')
  useEffect(() => alCambiarDatos(() => { pesos.recargar(); perfiles.recargar() }), [pesos.recargar, perfiles.recargar])
  const orden = (id: number) => { const i = ORDEN_PERSONAS.indexOf(per.lista.find((p) => p.id === id)?.code ?? ''); return i < 0 ? 99 : i }
  const codigo = (id: number) => per.lista.find((p) => p.id === id)?.code ?? 'shared'
  return (
    <div className="flex flex-col gap-6 tabular-nums">
      <Barra>
        <Boton variante="secundario" onClick={() => navegar(rutaDe('P-26', { seccion: '13' }))}><UserPen aria-hidden className="size-5" />{s.editarPerfil}</Boton>
        <Boton variante="primario" onClick={() => abrirRegistro('peso')}><Scale aria-hidden className="size-5" />{s.registrarPeso}</Boton>
      </Barra>
      <ConDatos consulta={perfiles}>{(d) => {
        const ps = itemsDe<Perfil>(d).sort((a, b) => orden(a.person_id) - orden(b.person_id))
        if (sub === 'detalle') return <ConDatos consulta={pesos}>{(w) => <Historial pesos={itemsDe<Peso>(w)} per={per} alCambiar={pesos.recargar} />}</ConDatos>
        if (sub === 'tendencias') return <div className="grid grid-cols-1 gap-6 desktop:grid-cols-2">
          <ConDatos consulta={pesos}>{(w) => <>{ps.map((p) => { const serie = itemsDe<Peso>(w).filter((x) => x.person_id === p.person_id).reverse()
            return <Grafica key={p.id} titulo={`${s.graficaPeso} · ${per.nombre(p.person_id)}`} resumen={s.resumenPeso} tipo="lineas" categorias={serie.map((x) => fecha(x.measured_date))} series={[{ nombre: per.nombre(p.person_id), valores: serie.map((x) => x.weight_g) }]} formato={kilos} /> })}</>}</ConDatos>
          <ConDatos consulta={kcal}>{(k) => <>{ps.map((p) => <Grafica key={p.id} titulo={`${s.kcalSemana} · ${per.nombre(p.person_id)}`} resumen={s.resumenKcal} tipo="lineas" categorias={k.weeks.map((x) => fecha(x.week_start))} formato={t.comun.kcal}
            series={[{ nombre: per.nombre(p.person_id), valores: k.weeks.map((x) => x.persons.find((y) => y.person_id === p.person_id)?.avg_kcal_per_day ?? 0) }, { nombre: s.metaDe(per.nombre(p.person_id)), valores: k.weeks.map((x) => x.persons.find((y) => y.person_id === p.person_id)?.daily_kcal_target ?? 0) }]} />)}</>}</ConDatos>
        </div>
        if (!ps.length) return <Vacio texto={s.completa} accion={<Boton variante="primario" onClick={() => navegar(rutaDe('P-26', { seccion: '13' }))}>{s.editarPerfil}</Boton>} />
        return <div className="grid grid-cols-1 gap-6 desktop:grid-cols-2">{ps.map((p) => <Ficha key={p.id} p={p} nombre={per.nombre(p.person_id)} codigo={codigo(p.person_id)} hoy={hoy.datos} alCambiar={perfiles.recargar} />)}</div>
      }}</ConDatos>
    </div>
  )
}

function Ficha({ p, nombre, codigo, hoy, alCambiar }: { p: Perfil; nombre: string; codigo: string; hoy?: Hoy; alCambiar: () => void }) {
  const avisar = useAvisos()
  const enLinea = useConexion()
  const sin = enLinea ? undefined : textos.conexion.requiere
  const confirmar = async () => {
    try { await pedir('patch', '/api/v1/health-profiles/{health_profile_id}', { ruta: { health_profile_id: p.id }, cuerpo: { daily_kcal_target: p.suggested_kcal_target } }); alCambiar(); avisar({ texto: s.metaGuardada, rol: 'exito' }) } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
  }
  const h = hoy?.persons.find((x) => x.person_id === p.person_id)
  const meta = p.daily_kcal_target
  return (
    <section aria-label={nombre} className="flex flex-col gap-4" data-prueba="ficha-salud">
      <h2 className="flex items-center gap-2 text-title-3 font-semibold"><ChipPersona codigo={codigo} nombre={nombre} /></h2>
      <Datos filas={[
        [s.altura, p.height_mm ? `${p.height_mm / 10} cm` : t.comun.sinDato], [s.edad, p.age_years != null ? s.anios(p.age_years) : t.comun.sinDato], [s.sexo, p.sex ? SEXO[p.sex] : t.comun.sinDato],
        [s.actividad, p.activity_level ? ACTIVIDAD[p.activity_level] : t.comun.sinDato], [s.objetivo, p.goal ? OBJETIVO[p.goal] : t.comun.sinDato], [s.restricciones, p.dietary_restrictions || t.comun.sinDato],
        [s.pesoActual, p.current_weight_g ? `${kilos(p.current_weight_g)}${p.current_weight_date ? ` · ${fecha(p.current_weight_date)}` : ''}` : t.comun.sinDato],
        [s.imc, p.bmi_tenths != null ? `${imc(p.bmi_tenths)} · ${p.bmi_class ? IMC[p.bmi_class] : ''}` : t.comun.sinDato],
      ]} />
      <Seccion titulo={s.meta}>
        <Datos filas={[[s.metaConfirmada, meta ? t.comun.kcal(meta) : <span className="text-role-warning">— · {s.sinConfirmar}</span>], [s.metaSugerida, p.suggested_kcal_target ? t.comun.kcal(p.suggested_kcal_target) : t.comun.sinDato]]} />
        {p.suggested_kcal_target != null && p.suggested_kcal_target !== meta && <Boton variante="primario" deshabilitadoPor={sin} onClick={confirmar}>{s.confirmar(t.comun.kcal(p.suggested_kcal_target))}</Boton>}
      </Seccion>
      <Seccion titulo={s.ajuste}>
        <div className="grid grid-cols-1 gap-3 tablet:grid-cols-2">
          <Ajuste p={p} campo="goal_adjust_lose_kcal" casa={p.household_goal_adjust_lose_kcal ?? 0} etiqueta={s.ajusteBajar} alCambiar={alCambiar} sin={sin} />
          <Ajuste p={p} campo="goal_adjust_gain_kcal" casa={p.household_goal_adjust_gain_kcal ?? 0} etiqueta={s.ajusteSubir} alCambiar={alCambiar} sin={sin} />
        </div>
      </Seccion>
      <Seccion titulo={s.hoy}>
        {h ? <Datos filas={[[s.planeadas, t.comun.kcal(h.planned_kcal)], [s.meta, meta && h.daily_kcal_target != null ? t.comun.kcal(h.daily_kcal_target) : `— · ${s.sinConfirmar}`], ...(meta && h.daily_kcal_target && h.difference_kcal != null ? [[s.diferencia, h.difference_kcal >= 0 ? s.porEncima(t.comun.kcal(h.difference_kcal)) : s.porDebajo(t.comun.kcal(-h.difference_kcal))] as [string, string]] : []), ...(h.eating_out_slots ? [['', s.fuera(h.eating_out_slots)] as [string, string]] : [])]} /> : <p className="text-small text-text-secondary">{t.comun.sinDato}</p>}
      </Seccion>
    </section>
  )
}

function Ajuste({ p, campo, casa, etiqueta, alCambiar, sin }: { p: Perfil; campo: 'goal_adjust_lose_kcal' | 'goal_adjust_gain_kcal'; casa: number; etiqueta: string; alCambiar: () => void; sin?: string }) {
  const avisar = useAvisos()
  const actual = p[campo] == null ? '' : String(p[campo])
  const [v, setV] = useState(actual)
  const [error, setError] = useState('')
  useEffect(() => setV(actual), [actual])
  const guardar = async (valor: number | null) => {
    try { await pedir('patch', '/api/v1/health-profiles/{health_profile_id}', { ruta: { health_profile_id: p.id }, cuerpo: { [campo]: valor } }); setError(''); alCambiar(); avisar({ texto: s.ajusteGuardado, rol: 'exito' }) } catch (e) { setV(actual); avisar({ texto: mensajeError(e), rol: 'error' }) }
  }
  const alSalir = () => {
    if (v === actual) return
    if (!v.trim()) return void guardar(null)
    if (!/^\d+$/.test(v.trim())) { setError(s.ajusteInvalido); setV(actual); return }
    void guardar(Number(v.trim()))
  }
  return (
    <div className="flex flex-col gap-2">
      <Campo etiqueta={etiqueta} inputMode="numeric" value={v} placeholder={String(casa)} ayuda={actual ? undefined : s.deLaCasa(casa)} error={error} disabled={!!sin} onChange={(e) => setV(e.target.value)} onBlur={alSalir} onKeyDown={(e) => { if (e.key === 'Enter') (e.target as HTMLInputElement).blur(); if (e.key === 'Escape') setV(actual) }} />
      {actual && <Boton variante="enlace" className="self-start" deshabilitadoPor={sin} onClick={() => { setV(''); void guardar(null) }}>{s.restablecer}</Boton>}
    </div>
  )
}

function Historial({ pesos, per, alCambiar }: { pesos: Peso[]; per: Personas; alCambiar: () => void }) {
  const avisar = useAvisos()
  const enLinea = useConexion()
  const [borrar, setBorrar] = useState<Peso | null>(null)
  const eliminar = async () => {
    if (!borrar) return
    try { await pedir('delete', '/api/v1/weight-entries/{weight_entry_id}', { ruta: { weight_entry_id: borrar.id } }); setBorrar(null); alCambiar(); avisar({ texto: s.eliminado, rol: 'exito' }) } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
  }
  return <Seccion titulo={s.historial}>
    <VistaBD def={{ clave: 'P-21-pesos', filas: pesos, idFila: (x) => x.id, vistas: ['tabla'], titulo: (x) => kilos(x.weight_g), vacio: <Vacio texto={s.completa} />,
      propiedades: [
        { id: 'persona', nombre: s.persona, valor: (x) => per.nombre(x.person_id), grupo: (x) => per.nombre(x.person_id) },
        { id: 'measured_date', nombre: s.fecha, valor: (x) => x.measured_date, pintar: (x) => fecha(x.measured_date) },
        { id: 'weight_g', nombre: s.peso, valor: (x) => x.weight_g, pintar: (x) => kilos(x.weight_g), alinear: 'derecha',
          ...(enLinea ? { editar: { tipo: 'numero' as const, textoInicial: (x: Peso) => String(x.weight_g / 1000),
            vistaPrevia: async (x: Peso, nuevo: string) => { const g = aGramos(nuevo); if (!g || g <= 0) throw new RangeError(s.pesoInvalido); return [{ etiqueta: s.peso, antes: kilos(x.weight_g), despues: kilos(g) }] },
            guardar: async (x: Peso, nuevo: string) => { const g = aGramos(nuevo); if (!g || g <= 0) throw new RangeError(s.pesoInvalido); await pedir('patch', '/api/v1/weight-entries/{weight_entry_id}', { ruta: { weight_entry_id: x.id }, cuerpo: { weight_g: g } }); alCambiar(); avisar({ texto: s.corregido, rol: 'exito' }) } } } : {}) },
        { id: 'bmi', nombre: s.imc, valor: (x) => x.bmi_tenths ?? null, pintar: (x) => imc(x.bmi_tenths), alinear: 'derecha' },
        { id: 'acciones', nombre: s.acciones, valor: () => '', pintar: (x) => <Boton variante="fantasma" tam="icono" aria-label={s.eliminarPeso(fecha(x.measured_date))} deshabilitadoPor={enLinea ? undefined : textos.conexion.requiere} onClick={(e) => { e.stopPropagation(); setBorrar(x) }}><Trash2 aria-hidden className="size-5" /></Boton> },
      ] }} />
    <Hoja abierta={!!borrar} alCambiar={(a) => !a && setBorrar(null)} titulo={s.eliminarPeso(borrar ? fecha(borrar.measured_date) : '')} pie={<div className="flex gap-2"><Boton className="flex-1" onClick={() => setBorrar(null)}>{textos.app.cancelar}</Boton><Boton className="flex-1" variante="destructivo" onClick={eliminar}>{t.comun.eliminar}</Boton></div>}>
      <p className="text-body">{s.confirmarEliminar}</p>
    </Hoja>
  </Seccion>
}
