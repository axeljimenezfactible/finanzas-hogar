// P-10 Plan de la semana (U1 §2.1–§2.2, §3.12 N-120 a N-129; flujo F-2). Semana por defecto de §0.2 (K38).
import { useEffect, useMemo, useState } from 'react'
import { ChevronLeft, ChevronRight, Utensils } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import type { PropsPantalla } from '../../app/registro'
import type { Esquemas } from '../../api/tipos'
import { nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { useConexion } from '../../api/conexion'
import { dinero, fechaConDia, hoyISO, mes as textoMes } from '../../api/formatos'
import { conConsulta, navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ChipSeleccionable, MarcaPersona } from '../../componentes/Chip'
import { ConDatos, Franja, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { Grafica } from '../../componentes/Grafica'
import { useAvisos } from '../../componentes/Avisos'
import { cn } from '../../componentes/utils'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { BarraAcciones, MenuMas, SelectorSemana, conSemana, diaCorto, diasDe, kcal, mensajeError, personasPlan, porId, porciones, semanaDe, useSemana, usarCatalogo } from './comun'
import { textosCC as t } from './textos'

type Plan = Esquemas['MealPlanWeek']
type Celda = Plan['cells'][number]
interface Abierta { fecha: string; tiempo: number }

export default function P10Plan({ sub, consulta }: PropsPantalla) {
  const sem = useSemana(consulta)
  const conexion = useConexion()
  const avisar = useAvisos()
  const plan = usarConsulta('get', '/api/v1/meal-plan/week', { consulta: { period: sem.semana } })
  const tiemposCat = usarCatalogo('tiempos')
  const tiempos = useMemo(() => [...(tiemposCat ?? [])].sort((a, b) => a.sort_order - b.sort_order), [tiemposCat])
  const personas = personasPlan(usarCatalogo('personas'))
  const recetas = usarCatalogo('recetas')
  const mapaRecetas = porId(recetas)
  const [abierta, setAbierta] = useState<Abierta | null>(null)
  const sinRed = conexion ? undefined : textos.conexion.requiere
  // «/comida» y «Agregar al plan» (P-12) abren el selector con ?hoja=asignar
  const hojaQ = consulta.get('hoja')
  useEffect(() => {
    if (hojaQ === 'asignar' && tiempos.length) {
      const hoy = hoyISO(); const dias = diasDe(sem.semana)
      setAbierta({ fecha: dias.includes(hoy) ? hoy : dias[0], tiempo: tiempos[0].id })
    }
  }, [hojaQ, tiempos, sem.semana])
  const cerrar = () => { setAbierta(null); if (consulta.get('hoja')) navegar(conConsulta({ hoja: null, receta: null }), { reemplazar: true }) }

  async function generar() {
    try {
      await pedir('post', '/api/v1/shopping-list/generate', { consulta: { week: sem.semana }, idempotencia: nuevaClave() })
      avisar({ texto: t.listaGenerada, rol: 'exito' })
      navegar(conSemana(rutaDe('P-13'), sem.semana, sem.porDefecto))
    } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
  }
  async function copiar() {
    try {
      const r = await pedir('post', '/api/v1/meal-plan/copy-previous-week', { consulta: { period: sem.semana }, idempotencia: nuevaClave() })
      avisar({ texto: t.copiadas(r?.created_count ?? 0), rol: 'exito' }); plan.recargar()
    } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
  }
  const nombreTiempo = (id: number) => tiempos.find((x) => x.id === id)?.name ?? ''
  const nombrePersona = (id: number) => personas.find((p) => p.id === id)?.display_name ?? ''
  const obligatorias = 7 * tiempos.filter((x) => !x.is_optional).length

  return (
    <div className="flex flex-col gap-4">
      <BarraAcciones acciones={<>
        <Boton variante="primario" deshabilitadoPor={sinRed} onClick={generar}>{t.generarLista}</Boton>
        <MenuMas acciones={[{ texto: t.copiarSemana, alHacer: copiar, deshabilitadoPor: sinRed }, { texto: t.asignarReceta, alHacer: () => navegar(conConsulta({ hoja: 'asignar' })), deshabilitadoPor: sinRed }]} />
      </>}><SelectorSemana {...sem} /></BarraAcciones>
      {sub === 'tendencias' ? <Tendencias personas={personas} /> : sub === 'calendario' ? <Calendario semana={sem.semana} /> : (
        <ConDatos consulta={plan}>{(p) => {
          const llenas = p.cells.filter((c) => c.recipe_id || c.is_eating_out)
          if (sub === 'resumen') return (
            <div className="grid gap-3 tablet:grid-cols-2 desktop:grid-cols-3">
              <TarjetaMetrica titulo={t.costoSemana} valor={dinero(p.week_cost_cents)} rol="expense" alAbrir={() => navegar(conConsulta({ sub: null }))} />
              <TarjetaMetrica titulo={t.sinArmar} valor={t.deTotal(p.unplanned_slots, obligatorias)} rol={p.unplanned_slots ? 'warning' : undefined} alAbrir={() => navegar(conConsulta({ sub: null }))} />
              {personas.map((per) => {
                const ks = p.days.map((d) => d.persons.find((x) => x.person_id === per.id)).filter((x) => x !== undefined)
                const prom = ks.length ? Math.floor(ks.reduce((s, x) => s + x.kcal, 0) / ks.length + 0.5) : null
                const meta = ks.find((x) => x.daily_kcal_target)?.daily_kcal_target
                return <TarjetaMetrica key={per.id} titulo={t.kcalDia(per.display_name)} valor={kcal(prom)} nota={meta ? t.meta(kcal(meta), kcal((prom ?? 0) - meta)) : t.sinMeta} />
              })}
              <TarjetaMetrica titulo={t.sobrantes} valor={porciones(p.leftovers.reduce((s, l) => s + l.servings_milli, 0))} nota={t.sobrantesNota} />
            </div>
          )
          return (
            <div className="flex flex-col gap-4">
              {llenas.length === 0 && <Vacio icono={<Utensils aria-hidden className="size-6" size={ICONO_MD} />} texto={t.semanaSinArmar} accion={<Boton deshabilitadoPor={sinRed} onClick={copiar}>{t.copiarSemana}</Boton>} />}
              <div className="grid gap-3 desktop:grid-cols-7" data-prueba="cuadricula">
                {diasDe(sem.semana).map((f) => {
                  const dia = p.days.find((d) => d.plan_date === f)
                  return (
                    <section key={f} aria-label={fechaConDia(f)} className={cn('flex flex-col gap-2 rounded-md border border-border-divider bg-surface-card p-3', f === hoyISO() && 'border-border-selected')}>
                      <div className="flex items-center justify-between gap-2">
                        <h2 className="text-small font-semibold">{diaCorto(f)}</h2>
                        {f === hoyISO() && <span className="rounded-full bg-surface-selected px-2 text-caption">{t.hoy}</span>}
                        <span className="text-small tabular-nums text-text-secondary">{dinero(dia?.cost_cents ?? 0)}</span>
                      </div>
                      {tiempos.map((ti) => (
                        <Casilla key={ti.id} celdas={p.cells.filter((c) => c.plan_date === f && c.meal_slot_id === ti.id)} tiempo={ti.name} personas={personas}
                          receta={(id) => mapaRecetas.get(id)?.name ?? '—'} alAbrir={() => setAbierta({ fecha: f, tiempo: ti.id })} />
                      ))}
                      {dia && <p className="text-caption text-text-secondary">{dia.persons.map((x) => `${nombrePersona(x.person_id)} ${kcal(x.kcal)}${x.difference_kcal != null ? ` (${x.difference_kcal > 0 ? '+' : ''}${x.difference_kcal.toLocaleString('es-MX')})` : ''}`).join(' · ')}</p>}
                    </section>
                  )
                })}
              </div>
              <h2 className="text-title-3 font-semibold">{t.asignaciones}</h2>
              <VistaBD def={{
                clave: 'P-10-asignaciones', filas: llenas, idFila: (c) => `${c.plan_date}-${c.meal_slot_id}-${c.person_id}`, vistas: ['tabla', 'tablero', 'galeria'],
                titulo: (c) => (c.is_eating_out ? t.fuera : mapaRecetas.get(c.recipe_id ?? 0)?.name ?? '—'),
                columnasTablero: { propiedad: 'tiempo', columnas: tiempos.map((x) => x.name) },
                alAbrir: (c) => setAbierta({ fecha: c.plan_date, tiempo: c.meal_slot_id }),
                propiedades: [
                  { id: 'dia', nombre: t.dia, valor: (c) => c.plan_date, pintar: (c) => fechaConDia(c.plan_date) },
                  { id: 'tiempo', nombre: t.tiempo, valor: (c) => nombreTiempo(c.meal_slot_id), grupo: (c) => nombreTiempo(c.meal_slot_id) },
                  { id: 'persona', nombre: t.persona, valor: (c) => nombrePersona(c.person_id) },
                  { id: 'receta', nombre: t.receta, valor: (c) => (c.is_eating_out ? t.fuera : mapaRecetas.get(c.recipe_id ?? 0)?.name ?? '—') },
                  { id: 'porciones', nombre: t.porcionesElegidas, valor: (c) => c.servings_milli ?? null, pintar: (c) => (c.servings_milli ? porciones(c.servings_milli) : '—'), alinear: 'derecha' },
                  { id: 'costo', nombre: t.costo, valor: (c) => c.cost_cents ?? null, pintar: (c) => (c.cost_cents != null ? dinero(c.cost_cents) : '—'), alinear: 'derecha' },
                  { id: 'kcal', nombre: t.kcal, valor: (c) => c.kcal ?? null, pintar: (c) => kcal(c.kcal), alinear: 'derecha' },
                ],
              }} />
              {abierta && <HojaAsignar abierta={abierta} cerrar={cerrar} semana={sem.semana} celdas={p.cells} tiempos={tiempos} personas={personas} recetas={recetas ?? []} preferida={Number(consulta.get('receta')) || undefined} recargar={plan.recargar} />}
            </div>
          )
        }}</ConDatos>
      )}
    </div>
  )
}

function Casilla({ celdas, tiempo, personas, receta, alAbrir }: { celdas: Celda[]; tiempo: string; personas: Esquemas['Person'][]; receta: (id: number) => string; alAbrir: () => void }) {
  const llenas = celdas.filter((c) => c.recipe_id || c.is_eating_out)
  const iguales = llenas.length === personas.length && llenas.length > 1 && llenas.every((c) => c.recipe_id === llenas[0].recipe_id && c.is_eating_out === llenas[0].is_eating_out && c.servings_milli === llenas[0].servings_milli)
  const nombre = (id: number) => personas.find((p) => p.id === id)?.display_name ?? ''
  const lineas = iguales
    ? [{ quien: t.ambosNombres(nombre(personas[0].id), nombre(personas[1].id)), c: llenas[0], c2: true, costo: llenas.reduce((s, x) => s + (x.cost_cents ?? 0), 0) }]
    : llenas.map((c) => ({ quien: nombre(c.person_id), c, c2: false, costo: c.cost_cents ?? 0 }))
  return (
    <button type="button" onClick={alAbrir} aria-label={`${tiempo}: ${llenas.length ? lineas.map((l) => `${l.quien} ${l.c.is_eating_out ? t.fuera : receta(l.c.recipe_id ?? 0)}`).join(', ') : t.vacia}`}
      className="flex min-h-11 w-full flex-col items-start gap-0.5 rounded-sm border border-border-divider px-2 py-1 text-left hover:bg-surface-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus">
      <span className="text-caption font-medium text-text-secondary">{tiempo}</span>
      {llenas.length === 0 ? <span className="text-small text-text-secondary">{t.vacia}</span> : lineas.map((l, k) => (
        <span key={`${l.c.person_id}-${k}`} className="flex flex-col text-small">
          <span className="text-caption text-text-secondary">{l.quien}{l.c.servings_milli ? ` · ${l.c2 ? t.porcionC(porciones(l.c.servings_milli)) : t.porciones(porciones(l.c.servings_milli))}` : ''}</span>
          <span className="font-medium">{l.c.is_eating_out ? t.fuera : receta(l.c.recipe_id ?? 0)}</span>
          {!l.c.is_eating_out && <span className="text-caption tabular-nums text-text-secondary">{kcal(l.c.kcal)}{l.c2 ? ' c/u' : ''} · {dinero(l.costo)}</span>}
        </span>
      ))}
    </button>
  )
}

function HojaAsignar({ abierta, cerrar, semana, celdas, tiempos, personas, recetas, preferida, recargar }: {
  abierta: Abierta; cerrar: () => void; semana: string; celdas: Celda[]; tiempos: Esquemas['MealSlot'][]; personas: Esquemas['Person'][]; recetas: Esquemas['Recipe'][]; preferida?: number; recargar: () => void
}) {
  const avisar = useAvisos()
  const conexion = useConexion()
  const [tiempo, setTiempo] = useState(abierta.tiempo)
  const [dias, setDias] = useState<string[]>([abierta.fecha])
  const [quien, setQuien] = useState<number | 'ambos'>('ambos')
  const [cuantas, setCuantas] = useState('1')
  const [texto, setTexto] = useState('')
  const [ocupado, setOcupado] = useState(false)
  const ids = quien === 'ambos' ? personas.map((p) => p.id) : [quien]
  const enCasilla = (f: string) => celdas.filter((c) => c.plan_date === f && c.meal_slot_id === tiempo && ids.includes(c.person_id) && c.meal_plan_entry_id)
  const llena = enCasilla(abierta.fecha).some((c) => c.recipe_id || c.is_eating_out)
  const sinRed = conexion ? undefined : textos.conexion.requiere
  const nTiempo = tiempos.find((x) => x.id === tiempo)?.name ?? ''
  const q = texto.trim().toLowerCase()
  const lista = recetas.filter((r) => r.is_active && (!r.meal_slot_ids || r.meal_slot_ids.includes(tiempo)) && (!q || r.name.toLowerCase().includes(q)))
    .sort((a, b) => Number(b.id === preferida) - Number(a.id === preferida) || (b.times_used_8w ?? 0) - (a.times_used_8w ?? 0) || a.name.localeCompare(b.name, 'es-MX'))
  async function guardar(cambio: { recipe_id: number | null; servings_milli: number | null; is_eating_out: boolean } | null) {
    setOcupado(true)
    try {
      for (const f of dias) {
        const existentes = enCasilla(f)
        for (const c of existentes) {
          if (cambio === null) await pedir('delete', '/api/v1/meal-plan-entries/{meal_plan_entry_id}', { ruta: { meal_plan_entry_id: c.meal_plan_entry_id! } })
          else await pedir('patch', '/api/v1/meal-plan-entries/{meal_plan_entry_id}', { ruta: { meal_plan_entry_id: c.meal_plan_entry_id! }, cuerpo: cambio })
        }
        const faltan = ids.filter((id) => !existentes.some((c) => c.person_id === id))
        if (cambio && faltan.length) await pedir('post', '/api/v1/meal-plan-entries', { idempotencia: nuevaClave(), cuerpo: { plan_dates: [f], meal_slot_id: tiempo, person_ids: faltan, is_eating_out: cambio.is_eating_out, ...(cambio.recipe_id ? { recipe_id: cambio.recipe_id, servings_milli: cambio.servings_milli ?? 1000 } : {}) } })
      }
      avisar({ texto: cambio === null ? t.casillaVaciada : t.recetaAsignada, rol: 'exito' })
      cerrar()
    } catch (e) {
      avisar({ texto: `${t.casillaNoGuardada}. ${mensajeError(e)}`, rol: 'error' })
    } finally { setOcupado(false); recargar() }
  }
  const milli = Math.round(Number(cuantas.replace(',', '.')) * 1000)
  const porcionValida = Number.isFinite(milli) && milli > 0
  return (
    <Hoja abierta alCambiar={(a) => !a && cerrar()} titulo={t.elegirReceta(diaCorto(abierta.fecha), nTiempo)}
      pie={<div className="flex flex-col gap-2 tablet:flex-row">
        <Boton tam="hoja" deshabilitadoPor={sinRed} disabled={ocupado} onClick={() => guardar({ recipe_id: null, servings_milli: null, is_eating_out: true })}>{t.fuera}</Boton>
        {llena && <Boton tam="hoja" variante="destructivo" deshabilitadoPor={sinRed} disabled={ocupado} onClick={() => guardar(null)}>{t.vaciarCasilla}</Boton>}
      </div>}>
      <div className="flex flex-col gap-3">
        <fieldset className="flex flex-col gap-1"><legend className="text-small font-medium text-text-secondary">{t.tiempo}</legend>
          <div className="flex flex-wrap gap-2">{tiempos.map((x) => <ChipSeleccionable key={x.id} seleccionado={x.id === tiempo} onClick={() => setTiempo(x.id)}>{x.name}</ChipSeleccionable>)}</div>
        </fieldset>
        <fieldset className="flex flex-col gap-1"><legend className="text-small font-medium text-text-secondary">{t.dias}</legend>
          <div className="flex flex-wrap gap-2">{diasDe(semana).map((f) => <ChipSeleccionable key={f} seleccionado={dias.includes(f)} onClick={() => setDias((d) => (d.includes(f) ? (d.length > 1 ? d.filter((x) => x !== f) : d) : [...d, f]))}>{diaCorto(f)}</ChipSeleccionable>)}</div>
        </fieldset>
        <fieldset className="flex flex-col gap-1"><legend className="text-small font-medium text-text-secondary">{t.persona}</legend>
          <div className="flex flex-wrap gap-2">
            <ChipSeleccionable seleccionado={quien === 'ambos'} marca={<MarcaPersona codigo="shared" />} onClick={() => setQuien('ambos')}>{t.ambos}</ChipSeleccionable>
            {personas.map((p) => <ChipSeleccionable key={p.id} seleccionado={quien === p.id} marca={<MarcaPersona codigo={p.code} />} onClick={() => setQuien(p.id)}>{p.display_name}</ChipSeleccionable>)}
          </div>
        </fieldset>
        <Campo etiqueta={t.porcionesPersona} inputMode="decimal" value={cuantas} onChange={(e) => setCuantas(e.target.value)} error={porcionValida ? undefined : t.cantidadInvalida} />
        <Campo etiqueta={t.buscarReceta} type="search" value={texto} onChange={(e) => setTexto(e.target.value)} ayuda={t.ordenUsadas} />
        {lista.length === 0 ? <Vacio texto={t.sinRecetasTiempo} /> : (
          <ul className="flex flex-col divide-y divide-border-divider rounded-md border border-border-divider" aria-label={t.receta}>
            {lista.map((r) => (
              <li key={r.id}>
                <button type="button" disabled={ocupado || !porcionValida} aria-disabled={sinRed ? true : undefined} onClick={() => !sinRed && guardar({ recipe_id: r.id, servings_milli: milli, is_eating_out: false })}
                  className={cn('flex min-h-11 w-full items-center justify-between gap-2 px-3 py-2 text-left hover:bg-surface-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus', r.id === preferida && 'bg-surface-selected')}>
                  <span className="font-medium">{r.name}</span>
                  <span className="text-caption tabular-nums text-text-secondary">{kcal(r.kcal_per_serving)} · {r.cost_per_serving_cents != null ? dinero(r.cost_per_serving_cents) : '—'}</span>
                </button>
              </li>
            ))}
          </ul>
        )}
        {sinRed && <Franja rol="sinConexion">{sinRed}</Franja>}
      </div>
    </Hoja>
  )
}

function Tendencias({ personas }: { personas: Esquemas['Person'][] }) {
  const c = usarConsulta('get', '/api/v1/meal-plan/trends', {})
  return (
    <ConDatos consulta={c}>{(d) => {
      if (!d.weeks.length) return <Vacio texto={textos.vacio.grafica} />
      const cats = d.weeks.map((w) => diaCorto(w.week_start) + ' ' + textoMes(w.week_start.slice(0, 7)))
      const serie = (id: number, campo: 'avg_kcal_per_day' | 'daily_kcal_target') => d.weeks.map((w) => w.persons.find((x) => x.person_id === id)?.[campo] ?? 0)
      return (
        <div className="grid gap-4 desktop:grid-cols-2">
          <Grafica titulo={t.costoSemanal} resumen={t.resumenCosto(d.weeks.length)} tipo="barras" categorias={cats} series={[{ nombre: t.costo, valores: d.weeks.map((w) => w.cost_cents) }]} />
          <Grafica titulo={t.kcalPromedio} resumen={t.resumenKcal} tipo="lineas" categorias={cats} formato={kcal}
            series={personas.flatMap((p) => [{ nombre: p.display_name, valores: serie(p.id, 'avg_kcal_per_day') }, { nombre: t.metaDe(p.display_name), valores: serie(p.id, 'daily_kcal_target') }])} />
        </div>
      )
    }}</ConDatos>
  )
}

function Calendario({ semana }: { semana: string }) {
  const [mes, setMes] = useState(() => diasDe(semana)[3].slice(0, 7))
  const c = usarConsulta('get', '/api/v1/meal-plan/month', { consulta: { period: mes } })
  const mover = (n: number) => { const [a, m] = mes.split('-').map(Number); const d = new Date(Date.UTC(a, m - 1 + n, 1)); setMes(d.toISOString().slice(0, 7)) }
  const huecos = (new Date(`${mes}-01T12:00:00Z`).getUTCDay() + 6) % 7
  const ROL: Record<string, string> = { complete: 'bg-role-success-surface', incomplete: 'bg-role-warning-surface', empty: 'bg-surface-card' }
  return (
    <section className="flex flex-col gap-3" aria-label={t.calendarioMes}>
      <div className="flex items-center gap-2">
        <Boton tam="icono" variante="fantasma" aria-label={t.mesAnterior} onClick={() => mover(-1)}><ChevronLeft aria-hidden className="size-6" size={ICONO_MD} /></Boton>
        <h2 className="text-title-3 font-semibold">{textoMes(mes)}</h2>
        <Boton tam="icono" variante="fantasma" aria-label={t.mesSiguiente} onClick={() => mover(1)}><ChevronRight aria-hidden className="size-6" size={ICONO_MD} /></Boton>
      </div>
      <ConDatos consulta={c}>{(d) => (
        <ol className="grid grid-cols-7 gap-1">
          {Array.from({ length: huecos }, (_, i) => <li key={`h${i}`} aria-hidden />)}
          {d.days.map((x) => (
            <li key={x.plan_date}>
              <button type="button" onClick={() => navegar(conConsulta({ sub: null, semana: semanaDe(x.plan_date) }))}
                className={cn('flex min-h-11 w-full flex-col items-center rounded-sm border border-border-divider p-1 text-small focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-border-focus', ROL[x.status])}
                aria-label={`${fechaConDia(x.plan_date)}: ${t.estadosDia[x.status]}`}>
                <span className="tabular-nums">{Number(x.plan_date.slice(8))}</span>
                <span className="hidden text-caption text-text-secondary tablet:inline">{t.estadosDia[x.status]}</span>
              </button>
            </li>
          ))}
        </ol>
      )}</ConDatos>
      <ul className="flex flex-wrap gap-3 text-caption text-text-secondary">
        {Object.entries(t.estadosDia).map(([k, v]) => <li key={k} className="flex items-center gap-1"><span aria-hidden className={cn('inline-block size-3 rounded-sm border border-border-divider', ROL[k])} />{v}</li>)}
      </ul>
    </section>
  )
}
