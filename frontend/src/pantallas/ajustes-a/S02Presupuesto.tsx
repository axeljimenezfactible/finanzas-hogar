// P-26.02 Presupuesto (U4 parte a): montos por categoría con alcance (K2), gastos fijos con vista previa 3,
// despensa programada con vista previa 4 (previewCreditPurchase) e interruptor de la proyección de los meses 7 a 12 (K21).
import { useEffect, useState } from 'react'
import { CircleAlert, Plus } from 'lucide-react'
import { nuevaClave, pedir } from '../../api/cliente'
import { ErrorApi } from '../../api/errores'
import { useConexion } from '../../api/conexion'
import { usarConsulta } from '../../api/consulta'
import { aCentavos, dinero, fecha, hoyISO, mes } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { Campo, estilosControl } from '../../componentes/Campo'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { textos } from '../../textos/es-MX'
import type { Cambio } from '../../vistas-bd/tipos'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { fechasEnMes, mesDe, mesSiguiente, montoEn, proximas, repartir, totalEnMes, ultimoDia } from './calculos'
import { Aviso, Bloque, DialogoCritico, Editable, ErrorGuardado, Fila, Interruptor, textoDinero, textoFallo, vDinero, vEntero, vFecha, vNombre, vOpcion, vPlazo, type Previa } from './comunes'
import { comoLista } from './comunes'
import { guardarVarias, nota, siNo } from './grupo'
import { t } from './textos'
const p = t.presupuesto
type PI = Esquemas['PlannedItem']; type Cat = Esquemas['Category']; type Regla = Esquemas['CostShareRule']; type Per = Esquemas['Person']
const esDespensa = (x: PI) => x.kind !== 'income' && x.accrual_months === 2

export default function S02Presupuesto() {
  const hoy = hoyISO(), mesActual = mesDe(hoy)
  const presupuestos = comoLista<Esquemas['Budget']>(usarConsulta('get', '/api/v1/budgets', { consulta: { period: mesActual, limit: 200 } }))
  const categorias = comoLista<Esquemas['Category']>(usarConsulta('get', '/api/v1/categories', { consulta: { limit: 200 } }))
  const reglas = comoLista<Esquemas['CostShareRule']>(usarConsulta('get', '/api/v1/cost-share-rules', { consulta: { limit: 200 } }))
  const personas = comoLista<Esquemas['Person']>(usarConsulta('get', '/api/v1/persons'))
  const partidas = comoLista<Esquemas['PlannedItem']>(usarConsulta('get', '/api/v1/planned-items', { consulta: { limit: 200 } }))
  const creditos = comoLista<Esquemas['Credit']>(usarConsulta('get', '/api/v1/credits', { consulta: { status: 'active', limit: 200 } }))
  const proyeccion = usarConsulta('get', '/api/v1/settings/projection')
  const [nuevo, setNuevo] = useState(false)
  useEffect(() => { const f = (e: KeyboardEvent) => { if (e.altKey && e.key.toLowerCase() === 'n') { e.preventDefault(); setNuevo(true) } }; window.addEventListener('keydown', f); return () => window.removeEventListener('keydown', f) }, [])
  const ctx = { cats: categorias.datos?.items ?? [], reglas: reglas.datos?.items ?? [], personas: personas.datos?.items ?? [], creditos: creditos.datos?.items ?? [] }
  return (
    <div className="flex flex-col gap-4">
      <Bloque titulo={p.proyeccion} id="b-proyeccion">
        <ConDatos consulta={proyeccion}>{(d) => <Interruptor etiqueta={p.proyeccion} encendido={d.repeat_cadence_enabled} guardar={async (v) => { await pedir('patch', '/api/v1/settings/projection', { cuerpo: { repeat_cadence_enabled: v } }); proyeccion.recargar() }} />}</ConDatos>
        <p className="text-small tabular-nums text-text-secondary">{p.proyeccionTexto}</p>
      </Bloque>
      <div className="flex flex-col gap-4 wide:grid wide:grid-cols-2">
        <Bloque titulo={p.montos} id="b-montos">
          <Fila etiqueta={p.periodo}><span className="px-2 text-body">{p.periodoTexto}</span></Fila>
          <ConDatos consulta={presupuestos}>{(b) => <Montos presupuestos={b.items} cats={ctx.cats} mesActual={mesActual} alCambiar={presupuestos.recargar} />}</ConDatos>
        </Bloque>
        <Bloque titulo={p.fijos} id="b-fijos" accion={<Boton variante="primario" onClick={() => setNuevo(true)}><Plus aria-hidden className="size-5" />{p.nuevoFijo}</Boton>}>
          <ConDatos consulta={partidas}>{(d) => <Fijos todas={d.items} ctx={ctx} mesActual={mesActual} hoy={hoy} alCambiar={partidas.recargar} alNuevo={() => setNuevo(true)} />}</ConDatos>
        </Bloque>
      </div>
      <Bloque titulo={p.despensa} id="b-despensa">
        <ConDatos consulta={partidas}>{(d) => <Despensa item={d.items.find(esDespensa)} creditos={ctx.creditos} hoy={hoy} alCambiar={partidas.recargar} />}</ConDatos>
      </Bloque>
      <NuevoFijo abierta={nuevo} alCambiar={setNuevo} cats={ctx.cats} alCrear={partidas.recargar} />
    </div>
  )
}

// ——— Montos por categoría y total, con el diálogo de alcance (vista previa 1)
function Montos({ presupuestos, cats, mesActual, alCambiar }: { presupuestos: Esquemas['Budget'][]; cats: Cat[]; mesActual: string; alCambiar: () => void }) {
  const [pendiente, setPendiente] = useState<{ cat: Cat | null; b?: Esquemas['Budget']; monto: number } | null>(null)
  const [alcance, setAlcance] = useState<'onward' | 'this_month'>('onward')
  /** Fallo al confirmar el monto (R4-01): `traslape` = 409 budget_overlap (V-17); `guardar` = cualquier otro (U4 §0.3). */
  const [fallo, setFallo] = useState<'traslape' | 'guardar' | null>(null)
  const [quitar, setQuitar] = useState<{ cat: Cat; b: Esquemas['Budget'] } | null>(null)
  const filas = [...cats.filter((c) => c.is_active).map((c) => ({ cat: c as Cat | null, b: presupuestos.find((x) => x.category_id === c.id) })), { cat: null, b: presupuestos.find((x) => x.category_id == null) }]
  const total = presupuestos.filter((x) => x.category_id != null).reduce((s, x) => s + x.amount_cents, 0)
  const nombre = (c: Cat | null) => c?.name ?? p.total
  const sig = mesSiguiente(mesActual)
  const confirmar = async () => {
    if (!pendiente) return
    const [a, m] = mesActual.split('-').map(Number)
    setFallo(null)
    try {
      if (pendiente.b) await pedir('patch', '/api/v1/budgets/{budget_id}', { ruta: { budget_id: pendiente.b.id }, consulta: { scope: alcance, month: mesActual }, cuerpo: { amount_cents: pendiente.monto } })
      else await pedir('post', '/api/v1/budgets', { idempotencia: nuevaClave(), cuerpo: { category_id: pendiente.cat?.id, amount_cents: pendiente.monto, valid_from_date: `${mesActual}-01`, ...(alcance === 'this_month' ? { valid_to_date: `${mesActual}-${ultimoDia(a, m)}` } : {}) } })
    } catch (e) {
      // U4 P-26.02 «Si se traslapa con otro monto» → V-17; cualquier otro fallo, el texto de §0.3. La hoja sigue abierta con el monto y el alcance capturados.
      setFallo(e instanceof ErrorApi && e.codigo === 'budget_overlap' ? 'traslape' : 'guardar')
      return
    }
    setPendiente(null); alCambiar()
  }
  const cerrar = () => { setPendiente(null); setFallo(null) }
  if (!presupuestos.length && !cats.length) return <Vacio texto={p.sinMontos} />
  return (
    <>
      {!presupuestos.length && <Aviso rol="info">{p.sinMontos} · {p.fijarMonto}</Aviso>}
      <ul className="flex flex-col">
        {filas.map(({ cat, b }) => (
          <li key={cat?.id ?? 'total'} className="flex flex-col gap-1 border-b border-border-divider py-1 last:border-b-0 desktop:grid desktop:grid-cols-4 desktop:items-center desktop:gap-2">
            <span className={cat ? "text-body desktop:col-span-2" : "text-body font-medium desktop:col-span-2"}>{nombre(cat)}{b?.is_single_month && <span className="ml-2 rounded-full bg-surface-selected px-2 text-caption">{p.soloEsteMes}</span>}
              {b?.is_single_month && b.following_month_amount_cents != null && <span className="block text-caption text-text-secondary">{p.siguienteMes(dinero(b.following_month_amount_cents))}</span>}</span>
            <Editable etiqueta={`${p.monto} ${nombre(cat)}`} tipo="dinero" valor={textoDinero(b?.amount_cents)} mostrar={<span className="w-full text-right tabular-nums">{b ? dinero(b.amount_cents) : p.fijarMonto}</span>}
              validar={vDinero()} guardar={async (s) => { setAlcance('onward'); setFallo(null); setPendiente({ cat, b, monto: aCentavos(s) ?? 0 }) }} />
            {cat && b ? <Boton variante="fantasma" onClick={() => setQuitar({ cat, b })}>{t.comun.eliminar}</Boton> : <span />}
          </li>
        ))}
      </ul>
      <Hoja abierta={!!pendiente} alCambiar={(a) => !a && cerrar()} titulo={p.alcance}
        pie={<div className="flex justify-end gap-2"><Boton onClick={cerrar}>{textos.app.cancelar}</Boton><Boton variante="primario" onClick={() => void confirmar()}>{textos.app.confirmar}</Boton></div>}>
        {pendiente && (
          <fieldset className="flex flex-col gap-3">
            <legend className="sr-only">{p.alcance}</legend>
            {(['onward', 'this_month'] as const).map((a) => (
              <label key={a} className="flex min-h-11 items-center gap-2"><input type="radio" name="alcance" className="size-5" checked={alcance === a} onChange={() => { setAlcance(a); setFallo(null) }} />{a === 'onward' ? p.enAdelante : p.soloMes}</label>
            ))}
            <Aviso rol="info"><span className="tabular-nums">{alcance === 'onward' ? p.previaAdelante(mes(mesActual), dinero(pendiente.monto)) : p.previaSolo(mes(mesActual), dinero(pendiente.monto), mes(sig), dinero(pendiente.b?.amount_cents ?? 0))}</span></Aviso>
            {pendiente.cat && <p className="text-small tabular-nums">{p.totalMes(mes(mesActual))}: {dinero(total)} → {dinero(total - (pendiente.b?.amount_cents ?? 0) + pendiente.monto)}</p>}
            {fallo === 'traslape' && <p role="alert" className="flex items-start gap-1 text-small text-role-error"><CircleAlert aria-hidden className="size-5 shrink-0" />{t.v.v17}</p>}
            {fallo === 'guardar' && <p role="alert" className="flex flex-wrap items-center gap-1 text-small text-role-error"><CircleAlert aria-hidden className="size-5" />{t.comun.noGuardo} · <button type="button" className="min-h-11 text-text-link underline-offset-2 hover:underline desktop:min-h-0" onClick={() => void confirmar()}>{textos.app.reintentar}</button></p>}
          </fieldset>
        )}
      </Hoja>
      <DialogoCritico abierto={!!quitar} alCambiar={(a) => !a && setQuitar(null)} titulo={quitar ? p.quitar(quitar.cat.name) : ''} verbo={quitar ? p.quitar(quitar.cat.name) : ''}
        alConfirmar={async () => { if (quitar) await pedir('delete', '/api/v1/budgets/{budget_id}', { ruta: { budget_id: quitar.b.id } }); alCambiar() }}>{quitar ? p.quitarTexto(quitar.cat.name) : ''}</DialogoCritico>
    </>
  )
}

// ——— Gastos fijos con la vista previa 3 (a)–(g)
type Ctx = { cats: Cat[]; reglas: Regla[]; personas: Per[]; creditos: Esquemas['Credit'][] }
function reglaDe(x: PI, ctx: Ctx): Regla | undefined {
  const id = x.cost_share_rule_id ?? ctx.cats.find((c) => c.id === x.category_id)?.default_cost_share_rule_id
  return ctx.reglas.find((r) => (id != null ? r.id === id : r.is_default))
}
function partesMes(x: PI, ctx: Ctx, aaaamm: string): Map<number, number> | null {
  const r = new Map<number, number>()
  if (!x.is_active) return r
  for (const f of fechasEnMes(x, aaaamm)) { const ps = repartir(montoEn(x, f), reglaDe(x, ctx)); if (!ps) return null; ps.forEach((s) => r.set(s.persona, (r.get(s.persona) ?? 0) + s.cents)) }
  return r
}
export function previaFijo(antes: PI, despues: PI, todas: PI[], ctx: Ctx, mesActual: string, hoy: string): Previa {
  if (!despues.frequency || !despues.start_date) return { pendiente: p.completa }
  const fechas = (x: PI) => proximas(x, hoy).map(fecha).join(' · ') || '—'
  const conFecha = (x: PI) => fechasEnMes(x, mesActual).length > 0
  let m = mesActual
  if (!conFecha(antes) && !conFecha(despues)) { const f = proximas(despues, hoy, 1)[0] ?? proximas(antes, hoy, 1)[0]; if (!f) return { notas: [p.sinFechas] }; m = mesDe(f) }
  const linea = (x: PI) => (!x.is_active ? p.inactiva : p.vecesPor(fechasEnMes(x, m).length, dinero(x.amount_cents), dinero(totalEnMes(x, m))))
  const suma = (lista: PI[]) => lista.filter((x) => x.is_active && x.kind === 'fixed_expense' && !esDespensa(x)).reduce((s, x) => s + totalEnMes(x, m), 0)
  const resto = todas.filter((x) => x.id !== antes.id)
  const cambios: Cambio[] = [
    { etiqueta: p.proximas, antes: fechas(antes), despues: fechas(despues) },
    { etiqueta: p.enMes(mes(m)), antes: linea(antes), despues: linea(despues) },
    { etiqueta: p.fijosDe(mes(m)), antes: dinero(suma([...resto, antes])), despues: dinero(suma([...resto, despues])) },
  ]
  const pa = partesMes(antes, ctx, m), pd = partesMes(despues, ctx, m)
  const notas: string[] = []
  if (!pd) { const r = reglaDe(despues, ctx); const fija = r?.parts?.find((x) => x.fixed_cents != null && x.fixed_cents > despues.amount_cents); notas.push(p.parteFijaPasa(ctx.personas.find((x) => x.id === fija?.person_id)?.display_name ?? '', dinero(fija?.fixed_cents ?? 0), dinero(despues.amount_cents))) }
  else {
    const ids = [...new Set([...(pa?.keys() ?? []), ...pd.keys()])]
    ids.forEach((id) => { const r = reglaDe(despues, ctx) ?? reglaDe(antes, ctx); const nom = id === r?.remainder_person_id ? p.resto : (ctx.personas.find((x) => x.id === id)?.display_name ?? ''); cambios.push({ etiqueta: `${nom} · ${mes(m)}`, antes: dinero(pa?.get(id) ?? 0), despues: dinero(pd.get(id) ?? 0) }) })
  }
  if (antes.person_id !== despues.person_id) { const n = (id?: number) => ctx.personas.find((x) => x.id === id)?.display_name ?? p.hogar; cambios.push({ etiqueta: p.atribuye, antes: n(antes.person_id), despues: n(despues.person_id) }) }
  notas.push(p.avisoP07)
  return { cambios, notas }
}
/** Vista previa 3 (h): edición múltiple con Ctrl+/ de la categoría o de «activo», sumada sobre las partidas elegidas en el mes en curso. */
export function previaFijosGrupo(sel: PI[], prop: string, valor: string, todas: PI[], ctx: Ctx, m: string): Cambio[] {
  const cambiar = (x: PI): PI => (prop === 'category_id' ? { ...x, category_id: Number(valor) } : { ...x, is_active: valor === 'si' })
  const ids = new Set(sel.map((x) => x.id))
  const suma = (lista: PI[]) => lista.filter((x) => x.is_active && x.kind === 'fixed_expense' && !esDespensa(x)).reduce((s, x) => s + totalEnMes(x, m), 0)
  const cambios: Cambio[] = [{ etiqueta: p.fijosDe(mes(m)), antes: dinero(suma(todas)), despues: dinero(suma(todas.map((x) => (ids.has(x.id) ? cambiar(x) : x)))) }]
  const pa = new Map<number, number>(), pd = new Map<number, number>(), restos = new Set<number>()
  let pasa: PI | null = null
  for (const x of sel) {
    const y = cambiar(x), a = partesMes(x, ctx, m), d = partesMes(y, ctx, m)
    if (!d) { pasa ??= y; continue }
    a?.forEach((c, id) => pa.set(id, (pa.get(id) ?? 0) + c))
    d.forEach((c, id) => pd.set(id, (pd.get(id) ?? 0) + c))
    for (const r of [reglaDe(x, ctx), reglaDe(y, ctx)]) if (r?.remainder_person_id != null) restos.add(r.remainder_person_id)
  }
  ctx.personas.filter((o) => pa.has(o.id) || pd.has(o.id)).forEach((o) => cambios.push({ etiqueta: `${restos.has(o.id) ? p.resto : o.display_name} · ${mes(m)}`, antes: dinero(pa.get(o.id) ?? 0), despues: dinero(pd.get(o.id) ?? 0) }))
  if (pasa) { const fija = reglaDe(pasa, ctx)?.parts?.find((x) => x.fixed_cents != null && x.fixed_cents > pasa.amount_cents); cambios.push(nota(p.parteFijaPasa(ctx.personas.find((x) => x.id === fija?.person_id)?.display_name ?? '', dinero(fija?.fixed_cents ?? 0), dinero(pasa.amount_cents)))) }
  if (sel.some((x) => !x.frequency || !x.start_date)) cambios.push(nota(p.completa))
  if (prop === 'category_id') cambios.push(nota(t.grupo.tomanCategoria(sel.filter((x) => x.cost_share_rule_id == null).length, sel.length)))
  cambios.push(nota(p.avisoP07))
  return cambios
}
function Fijos({ todas, ctx, mesActual, hoy, alCambiar, alNuevo }: { todas: PI[]; ctx: Ctx; mesActual: string; hoy: string; alCambiar: () => void; alNuevo: () => void }) {
  const filas = todas.filter((x) => x.kind === 'fixed_expense' && !esDespensa(x)).sort((a, b) => Number(a.start_date.slice(8)) - Number(b.start_date.slice(8)))
  const guardar = (x: PI) => async (cambio: Partial<PI>) => { await pedir('patch', '/api/v1/planned-items/{planned_item_id}', { ruta: { planned_item_id: x.id }, cuerpo: cambio as Esquemas['PlannedItemUpdate'] }); alCambiar() }
  const previa = (x: PI) => (cambio: Partial<PI>) => previaFijo(x, { ...x, ...cambio }, todas, ctx, mesActual, hoy)
  const opc = (lista: { id: number; nombre: string }[], vacio: string) => [{ valor: '', texto: vacio }, ...lista.map((o) => ({ valor: String(o.id), texto: o.nombre }))]
  const num = (s: string) => (s ? Number(s) : null)
  const campo = (x: PI, id: keyof PI, nombre: string, tipo: 'texto' | 'dinero' | 'numero' | 'fecha' | 'opcion' | 'interruptor', a: { valor: string; mostrar: React.ReactNode; aCambio: (s: string) => Partial<PI>; validar?: (s: string) => string | null; opciones?: { valor: string; texto: string }[]; conPrevia?: boolean }) =>
    <Editable etiqueta={`${nombre} · ${x.description}`} tipo={tipo} valor={a.valor} mostrar={a.mostrar} opciones={a.opciones} validar={a.validar} vistaPrevia={a.conPrevia ? (s) => previa(x)(a.aCambio(s)) : undefined} guardar={(s) => guardar(x)(a.aCambio(s))} key={String(id)} />
  if (!filas.length) return <Vacio texto={p.sinFijos} accion={<div className="flex flex-col items-center gap-1"><Boton variante="primario" onClick={alNuevo}>{p.nuevoFijo}</Boton><span className="text-small text-text-secondary">{p.importarCsv}</span></div>} />
  const cats = opc(ctx.cats.map((c) => ({ id: c.id, nombre: c.name })), t.v.v22)
  return (
    <VistaBD def={{ clave: 'P-26.02-fijos', filas, idFila: (x) => x.id, vistas: ['tabla'], titulo: (x) => x.description,
      // Ctrl+/ (U4 P-26.02 «Atajos» y punto h): categoría o activo de varias partidas, un updatePlannedItem por partida.
      editarVarias: {
        propiedades: [
          { id: 'category_id', nombre: p.categoria, tipo: 'opcion', opciones: ctx.cats.filter((c) => c.is_active).map((c) => ({ valor: String(c.id), texto: c.name })) },
          { id: 'is_active', nombre: t.comun.activo, tipo: 'opcion', opciones: siNo },
        ],
        vistaPrevia: async (sel, prop, valor) => previaFijosGrupo(sel, prop, valor, todas, ctx, mesActual),
        guardar: (sel, prop, valor) => guardarVarias(sel, (x) => pedir('patch', '/api/v1/planned-items/{planned_item_id}', { ruta: { planned_item_id: x.id }, cuerpo: prop === 'category_id' ? { category_id: Number(valor) } : { is_active: valor === 'si' } }), alCambiar),
      },
      propiedades: [
        { id: 'description', nombre: p.concepto, valor: (x) => x.description, pintar: (x) => campo(x, 'description', p.concepto, 'texto', { valor: x.description, mostrar: x.description, aCambio: (s) => ({ description: s.trim() }), validar: vNombre }) },
        { id: 'category_id', nombre: p.categoria, valor: (x) => x.category_id, grupo: (x) => ctx.cats.find((c) => c.id === x.category_id)?.name ?? '', pintar: (x) => campo(x, 'category_id', p.categoria, 'opcion', { valor: String(x.category_id ?? ''), mostrar: ctx.cats.find((c) => c.id === x.category_id)?.name ?? '—', opciones: cats, aCambio: (s) => ({ category_id: num(s) ?? undefined }), validar: vOpcion(), conPrevia: x.cost_share_rule_id == null }) },
        { id: 'amount_cents', nombre: p.montoVez, alinear: 'derecha', valor: (x) => x.amount_cents, pintar: (x) => campo(x, 'amount_cents', p.montoVez, 'dinero', { valor: textoDinero(x.amount_cents), mostrar: <span className="tabular-nums">{dinero(x.amount_cents)}</span>, aCambio: (s) => ({ amount_cents: aCentavos(s) ?? 0 }), validar: vDinero(true), conPrevia: true }) },
        { id: 'frequency', nombre: p.periodicidad, valor: (x) => x.frequency, pintar: (x) => campo(x, 'frequency', p.periodicidad, 'opcion', { valor: x.frequency, mostrar: p.frecuencias[x.frequency], opciones: Object.entries(p.frecuencias).map(([valor, texto]) => ({ valor, texto })), aCambio: (s) => ({ frequency: s as PI['frequency'] }), validar: vOpcion(), conPrevia: true }) },
        { id: 'interval_count', nombre: p.cadaCuanto, alinear: 'derecha', valor: (x) => x.interval_count, pintar: (x) => campo(x, 'interval_count', p.cadaCuanto, 'numero', { valor: String(x.interval_count), mostrar: String(x.interval_count), aCambio: (s) => ({ interval_count: Number(s) }), validar: vEntero(1, 24), conPrevia: true }) },
        { id: 'start_date', nombre: p.primeraFecha, valor: (x) => x.start_date, pintar: (x) => campo(x, 'start_date', p.primeraFecha, 'fecha', { valor: x.start_date, mostrar: fecha(x.start_date), aCambio: (s) => ({ start_date: s }), validar: vFecha(), conPrevia: true }) },
        { id: 'end_date', nombre: p.fechaFinal, valor: (x) => x.end_date, pintar: (x) => campo(x, 'end_date', p.fechaFinal, 'fecha', { valor: x.end_date ?? '', mostrar: x.end_date ? fecha(x.end_date) : p.sinFin, aCambio: (s) => ({ end_date: s || undefined }), validar: (s) => vFecha(true)(s) ?? (s && s < x.start_date ? t.v.v08 : null), conPrevia: true }) },
        { id: 'person_id', nombre: p.persona, valor: (x) => x.person_id, pintar: (x) => campo(x, 'person_id', p.persona, 'opcion', { valor: String(x.person_id ?? ''), mostrar: ctx.personas.find((o) => o.id === x.person_id)?.display_name ?? p.hogar, opciones: opc(ctx.personas.map((o) => ({ id: o.id, nombre: o.display_name })), p.hogar), aCambio: (s) => ({ person_id: num(s) ?? undefined }), conPrevia: true }) },
        { id: 'cost_share_rule_id', nombre: p.reparto, valor: (x) => x.cost_share_rule_id, pintar: (x) => campo(x, 'cost_share_rule_id', p.reparto, 'opcion', { valor: String(x.cost_share_rule_id ?? ''), mostrar: ctx.reglas.find((r) => r.id === x.cost_share_rule_id)?.name ?? t.categorias.porDefecto, opciones: opc(ctx.reglas.map((r) => ({ id: r.id, nombre: r.name })), t.categorias.porDefecto), aCambio: (s) => ({ cost_share_rule_id: num(s) ?? undefined }), conPrevia: true }) },
        { id: 'payment_method', nombre: p.formaPago, valor: (x) => x.payment_method, pintar: (x) => campo(x, 'payment_method', p.formaPago, 'opcion', { valor: x.payment_method ?? '', mostrar: x.payment_method ? p.formas[x.payment_method] : '—', opciones: [{ valor: '', texto: '—' }, ...Object.entries(p.formas).map(([valor, texto]) => ({ valor, texto }))], aCambio: (s) => ({ payment_method: (s || undefined) as PI['payment_method'] }) }) },
        { id: 'is_active', nombre: t.comun.activo, valor: (x) => x.is_active, grupo: (x) => (x.is_active ? t.comun.activo : t.comun.no), pintar: (x) => campo(x, 'is_active', t.comun.activo, 'interruptor', { valor: x.is_active ? 'si' : 'no', mostrar: null, aCambio: (s) => ({ is_active: s === 'si' }), conPrevia: true }) },
      ] }} />
  )
}

function NuevoFijo({ abierta, alCambiar, cats, alCrear }: { abierta: boolean; alCambiar: (a: boolean) => void; cats: Cat[]; alCrear: () => void }) {
  const conexion = useConexion()
  const vacio = { description: '', category_id: '', amount: '', frequency: '', interval_count: '1', start_date: '' }
  const [v, setV] = useState(vacio)
  const [err, setErr] = useState<Record<string, string | null>>({})
  const [fallo, setFallo] = useState<string | null>(null)
  const [clave, setClave] = useState(nuevaClave)
  const guardar = async () => {
    const e = { description: vNombre(v.description), category_id: vOpcion()(v.category_id), amount: vDinero(true)(v.amount), frequency: vOpcion()(v.frequency), interval_count: vEntero(1, 24)(v.interval_count), start_date: vFecha()(v.start_date) }
    setErr(e); setFallo(null)
    if (Object.values(e).some(Boolean)) return
    try { await pedir('post', '/api/v1/planned-items', { idempotencia: clave, cuerpo: { kind: 'fixed_expense', description: v.description.trim(), category_id: Number(v.category_id), amount_cents: aCentavos(v.amount) ?? 0, frequency: v.frequency as PI['frequency'], interval_count: Number(v.interval_count), start_date: v.start_date, is_active: true } }) }
    catch (x) { setFallo(textoFallo(x)); return }
    setV(vacio); setClave(nuevaClave()); alCambiar(false); alCrear()
  }
  const sel = (id: keyof typeof v, etiqueta: string, opciones: { valor: string; texto: string }[]) => (
    <label className="flex flex-col gap-1 text-small font-medium text-text-secondary">{etiqueta}
      <select className={estilosControl} value={v[id]} aria-invalid={err[id] ? true : undefined} onChange={(e) => setV({ ...v, [id]: e.target.value })}><option value="">{t.v.v22}</option>{opciones.map((o) => <option key={o.valor} value={o.valor}>{o.texto}</option>)}</select>
      {err[id] && <span className="text-role-error">{err[id]}</span>}
    </label>
  )
  return (
    <Hoja abierta={abierta} alCambiar={alCambiar} titulo={p.nuevoFijo} pie={<div className="flex flex-col gap-2">{fallo && <ErrorGuardado texto={fallo} alReintentar={() => void guardar()} />}<div className="flex justify-end gap-2"><Boton onClick={() => alCambiar(false)}>{textos.app.cancelar}</Boton><Boton variante="primario" deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => void guardar()}>{textos.app.guardar}</Boton></div></div>}>
      <form className="flex flex-col gap-3" onKeyDown={(e) => { if (e.key === 'Enter' && e.ctrlKey) { e.preventDefault(); void guardar() } }} onSubmit={(e) => { e.preventDefault(); void guardar() }}>
        <Campo etiqueta={p.concepto} value={v.description} error={err.description ?? undefined} onChange={(e) => setV({ ...v, description: e.target.value })} />
        {sel('category_id', p.categoria, cats.filter((c) => c.is_active).map((c) => ({ valor: String(c.id), texto: c.name })))}
        <Campo etiqueta={p.montoVez} inputMode="decimal" value={v.amount} error={err.amount ?? undefined} onChange={(e) => setV({ ...v, amount: e.target.value })} />
        {sel('frequency', p.periodicidad, Object.entries(p.frecuencias).map(([valor, texto]) => ({ valor, texto })))}
        <Campo etiqueta={p.cadaCuanto} inputMode="numeric" value={v.interval_count} error={err.interval_count ?? undefined} onChange={(e) => setV({ ...v, interval_count: e.target.value })} />
        <Campo etiqueta={p.primeraFecha} type="date" value={v.start_date} error={err.start_date ?? undefined} onChange={(e) => setV({ ...v, start_date: e.target.value })} />
      </form>
    </Hoja>
  )
}

// ——— Despensa programada (vista previa 4). Alta PROVISIONAL: bloqueo WF3-ajustes-a-despensa (tipo de partida).
function Despensa({ item, creditos, hoy, alCambiar }: { item?: PI; creditos: Esquemas['Credit'][]; hoy: string; alCambiar: () => void }) {
  if (!item) return <Vacio texto={p.sinDespensa} />
  const guardar = async (c: Partial<PI>) => { await pedir('patch', '/api/v1/planned-items/{planned_item_id}', { ruta: { planned_item_id: item.id }, cuerpo: c as Esquemas['PlannedItemUpdate'] }); alCambiar() }
  const pagos = async (x: PI): Promise<string> => {
    if (!x.credit_id || !x.term_count) return p.eligeLinea
    const f = proximas(x, hoy, 1)[0]
    if (!f) return '—'
    const r = await pedir('post', '/api/v1/credit-purchases/preview', { cuerpo: { credit_id: x.credit_id, kind: 'installment_plan', purchase_date: f, principal_cents: x.amount_cents, term_count: x.term_count, interest_source: 'term_table' } as Esquemas['CreditPurchaseCreate'] })
    return `${r.schedule.map((s) => `${fecha(s.due_date)} ${dinero(s.amount_cents)}`).join(' · ')} · ${p.totalIntereses}: ${dinero(r.total_cents)}`
  }
  const previa = async (c: Partial<PI>): Promise<Previa> => {
    const d = { ...item, ...c }
    const [a, b] = await Promise.all([pagos(item), pagos(d)])
    const linea = creditos.find((x) => x.id === d.credit_id)
    return { cambios: [{ etiqueta: p.proximas, antes: proximas(item, hoy).map(fecha).join(' · '), despues: proximas(d, hoy).map(fecha).join(' · ') }, { etiqueta: p.pagos, antes: a, despues: b }], notas: linea ? [p.disponibleBaja(linea.name, dinero(d.amount_cents))] : [p.eligeLinea] }
  }
  const lineas = [{ valor: '', texto: t.v.v22 }, ...creditos.map((c) => ({ valor: String(c.id), texto: c.name }))]
  return (
    <div className="flex flex-col desktop:grid desktop:grid-cols-2 desktop:gap-x-6">
      <Fila etiqueta={p.montoBase}><Editable etiqueta={p.montoBase} tipo="dinero" valor={textoDinero(item.amount_cents)} mostrar={<span className="tabular-nums">{dinero(item.amount_cents)}</span>} validar={vDinero(true)} vistaPrevia={(s) => previa({ amount_cents: aCentavos(s) ?? 0 })} guardar={(s) => guardar({ amount_cents: aCentavos(s) ?? 0 })} /></Fila>
      <Fila etiqueta={p.frecuencia}><span className="px-2">{p.bimestral}</span></Fila>
      <Fila etiqueta={p.devengo}><span className="px-2">{p.mitad}</span></Fila>
      <Fila etiqueta={p.dia}><Editable etiqueta={p.dia} tipo="fecha" valor={item.start_date} mostrar={fecha(item.start_date)} validar={vFecha()} vistaPrevia={(s) => previa({ start_date: s })} guardar={(s) => guardar({ start_date: s })} /></Fila>
      <Fila etiqueta={p.proxima}><span className="px-2">{item.next_date ? fecha(item.next_date) : (proximas(item, hoy, 1)[0] ? fecha(proximas(item, hoy, 1)[0]) : '—')}</span></Fila>
      <Fila etiqueta={p.linea}><Editable etiqueta={p.linea} tipo="opcion" opciones={lineas} valor={String(item.credit_id ?? '')} mostrar={creditos.find((c) => c.id === item.credit_id)?.name ?? '—'} validar={vOpcion()} vistaPrevia={(s) => previa({ credit_id: Number(s), payment_method: 'credit' })} guardar={(s) => guardar({ credit_id: Number(s), payment_method: 'credit' })} /></Fila>
      <Fila etiqueta={p.plazo}><Editable etiqueta={p.plazo} tipo="numero" valor={String(item.term_count ?? '')} mostrar={item.term_count ? String(item.term_count) : '—'} validar={vPlazo} vistaPrevia={(s) => previa({ term_count: Number(s) })} guardar={(s) => guardar({ term_count: Number(s) })} /></Fila>
    </div>
  )
}
