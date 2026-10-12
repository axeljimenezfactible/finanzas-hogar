// P-26.06 Tiendas y precios (U4 parte a): tiendas y mercados, membresía, lectura asistida con veredicto (acción crítica),
// fuentes de precio, corte PROFECO y parámetros de precios (K44) con la vista previa de frescura (K22, K40).
import { useEffect, useState, type KeyboardEvent as EventoTecla, type ReactNode } from 'react'
import { ChevronLeft, ChevronRight, Plus, Upload } from 'lucide-react'
import { nuevaClave, pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { usarConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import { aCentavos, dinero, fecha, fechaHora, hoyISO } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Enlace } from '../../app/enrutador'
import { Boton } from '../../componentes/Boton'
import { Campo, estilosControl } from '../../componentes/Campo'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { textos } from '../../textos/es-MX'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { pesoPrecio, sumarDias } from './calculos'
import { Aviso, Bloque, Editable, Fila, Interruptor, Parametro, textoDinero, vDinero, vEntero, vNombre, vOpcion, type Previa } from './comunes'
import { comoLista } from './comunes'
import { t } from './textos'
const s = t.tiendas
type St = Esquemas['Store']
const FECHA_EJEMPLO = '2026-10-10' // precio de la semilla (U4 P-26.06, vista previa 1)
const dias = s.dias.map((d, i) => ({ valor: String(i + 1), texto: d }))

export function frescura(viejo: number, cero: number) {
  const medio = Math.round((viejo + cero) / 2)
  return s.frescura(fecha(FECHA_EJEMPLO), fecha(sumarDias(FECHA_EJEMPLO, viejo)), viejo, pesoPrecio(viejo, viejo, cero), medio, pesoPrecio(medio, viejo, cero), fecha(sumarDias(FECHA_EJEMPLO, cero)), cero)
}
/** Días entre dos fechas ISO (aaaa-mm-dd), sin hora. */
const diasEntre = (desde: string, hasta: string) => { const u = (f: string) => { const [a, m, d] = f.split('-').map(Number); return Date.UTC(a, m - 1, d) }; return Math.round((u(hasta) - u(desde)) / 86_400_000) }
/** Lunes de la semana de una fecha (K8: la semana va de lunes a domingo). */
const lunesDe = (f: string) => { const [a, m, d] = f.split('-').map(Number); return sumarDias(f, -((new Date(Date.UTC(a, m - 1, d)).getUTCDay() + 6) % 7)) }
export default function S06Tiendas() {
  const tiendas = comoLista<Esquemas['Store']>(usarConsulta('get', '/api/v1/stores', { consulta: { limit: 200 } }))
  const grupos = comoLista<Esquemas['ProductGroup']>(usarConsulta('get', '/api/v1/product-groups', { consulta: { limit: 200 } }))
  const fuentes = usarConsulta('get', '/api/v1/price-sources')
  const ajustes = usarConsulta('get', '/api/v1/settings/prices')
  const ia = usarConsulta('get', '/api/v1/ai/providers')
  // Fecha del último precio de cada producto (listProducts.last_observed_date) para contar los que pasan de fresco a viejo (U4 P-26.06, vista previa 1).
  const productos = comoLista<Esquemas['Product'] & { last_observed_date?: string }>(usarConsulta('get', '/api/v1/products', { consulta: { limit: 200 } }))
  const [abierta, setAbierta] = useState<number | null>(null)
  const [nueva, setNueva] = useState(false)
  useEffect(() => { const f = (e: KeyboardEvent) => { if (e.altKey && e.key.toLowerCase() === 'n') { e.preventDefault(); setNueva(true) } }; window.addEventListener('keydown', f); return () => window.removeEventListener('keydown', f) }, [])
  const gs = grupos.datos?.items ?? []
  return (
    <div className="flex flex-col gap-4">
      <Bloque titulo={s.tiendas} id="b-tiendas" accion={<Boton variante="primario" onClick={() => setNueva(true)}><Plus aria-hidden className="size-4" />{s.nueva}</Boton>}>
        <ConDatos consulta={tiendas}>{(d) => {
          if (!d.items.length) return <Vacio texto={s.sin} accion={<Boton variante="primario" onClick={() => setNueva(true)}>{s.nueva}</Boton>} />
          const guardar = (x: St) => async (c: Esquemas['StoreUpdate']) => { await pedir('patch', '/api/v1/stores/{store_id}', { ruta: { store_id: x.id }, cuerpo: c }); tiendas.recargar() }
          return <>
            <VistaBD def={{ clave: 'P-26.06-tiendas', filas: d.items, idFila: (x) => x.id, vistas: ['tabla', 'galeria', 'tablero'], titulo: (x) => x.name, alAbrir: (x) => setAbierta(x.id),
              columnasTablero: { propiedad: 'kind', columnas: Object.values(s.tipos) },
              propiedades: [
                { id: 'name', nombre: s.nombre, valor: (x) => x.name, pintar: (x) => <button type="button" className="min-h-11 text-left text-text-link desktop:min-h-0" onClick={() => setAbierta(x.id)}>{x.name}</button> },
                { id: 'kind', nombre: s.tipo, valor: (x) => x.kind, grupo: (x) => s.tipos[x.kind], pintar: (x) => s.tipos[x.kind] },
                { id: 'visit_weekday', nombre: s.diaVisita, valor: (x) => x.visit_weekday, pintar: (x) => (x.kind === 'market' ? <Editable etiqueta={`${s.diaVisita} · ${x.name}`} tipo="opcion" valor={String(x.visit_weekday ?? '')} mostrar={x.visit_weekday ? s.dias[x.visit_weekday - 1] : t.comun.sinValor} opciones={[{ valor: '', texto: t.v.v11 }, ...dias]} validar={(v) => (v ? null : t.v.v11)} guardar={(v) => guardar(x)({ visit_weekday: Number(v) })} /> : '—') },
                { id: 'groups', nombre: s.queCompra, valor: (x) => (x.product_group_ids ?? []).length, pintar: (x) => (x.product_group_ids ?? []).map((id) => gs.find((g) => g.id === id)?.name).filter(Boolean).join(', ') || s.gruposVacios },
                { id: 'membership', nombre: s.costo, alinear: 'derecha', valor: (x) => x.membership_annual_fee_cents, pintar: (x) => (x.requires_membership ? (x.membership_annual_fee_cents != null ? <span className="tabular-nums">{dinero(x.membership_annual_fee_cents)}</span> : s.capturaCosto) : '—') },
                { id: 'verdict', nombre: s.veredicto, valor: (x) => x.assisted_reading_verdict, pintar: (x) => (x.kind === 'chain' ? s.veredictos[x.assisted_reading_verdict ?? 'missing'] : '—') },
                { id: 'is_active', nombre: t.comun.activa, valor: (x) => x.is_active, pintar: (x) => (x.is_active ? t.comun.si : t.comun.no) },
              ] }} />
            <CalendarioVisitas tiendas={d.items} alAbrir={setAbierta} />
            {abierta != null && d.items.find((x) => x.id === abierta) && <Detalle x={d.items.find((x) => x.id === abierta)!} cadenas={d.items.filter((x) => x.kind === 'chain')} grupos={gs} guardar={guardar(d.items.find((x) => x.id === abierta)!)} alGrupos={tiendas.recargar} alCerrar={() => setAbierta(null)} />}
          </>
        }}</ConDatos>
      </Bloque>
      <div className="flex flex-col gap-4 desktop:grid desktop:grid-cols-2">
        <Bloque titulo={s.fuentes} id="b-fuentes">
          <ConDatos consulta={fuentes}>{(f) => <Fuentes f={f} ia={ia.datos?.providers ?? []} maxGb={ajustes.datos?.parameters['price.profeco_import_max_gb'].value ?? null} alCambiar={() => { fuentes.recargar(); ajustes.recargar() }} />}</ConDatos>
        </Bloque>
        <Bloque titulo={s.parametros} id="b-parametros-precios">
          <ConDatos consulta={ajustes}>{(a) => {
            const p = a.parameters, viejo = p['price.stale_after_days'].value ?? 30, cero = p['price.zero_weight_after_days'].value ?? 90
            type K = keyof typeof p
            const patch = (c: Esquemas['PriceSettingsUpdate']) => pedir('patch', '/api/v1/settings/prices', { cuerpo: c }).then(() => ajustes.recargar())
            const prop = (cl: K) => ({ param: p[cl], guardar: (n: number) => patch({ parameters: { [cl]: n } }), restablecer: () => patch({ reset_parameters: [cl] }) })
            // Cuántos pasan hoy de fresco a viejo con el valor nuevo: productos cuyo último precio tiene entre «nuevo» y «vigente − 1» días.
            // PROVISIONAL (bloqueo WF3-ajustes-precios-conteo): C1 no tiene una vista previa del conteo; se cuenta con la misma unidad que getPriceSummary (producto por su último precio).
            const conteo = (v: number): ReactNode => {
              const pr = productos.datos
              if (!pr || pr.next_cursor) return <span data-prueba="conteo-viejos">{`${s.conteoTitulo}: ${t.comun.previaPendiente} · ${t.comun.provisional} (${s.bloqueoConteo})`}</span>
              const hoy = hoyISO()
              const n = pr.items.filter((x) => { if (!x.last_observed_date) return false; const e = diasEntre(x.last_observed_date, hoy); return e >= v && e < viejo }).length
              return <span data-prueba="conteo-viejos" className="tabular-nums">{`${s.conteo(n)} · ${t.comun.provisional} (${s.bloqueoConteo})`}</span>
            }
            const previa = (v: number, c: number, conConteo: boolean): Previa => ({ cambios: [{ etiqueta: s.pViejo.split(' para')[0], antes: frescura(viejo, cero), despues: frescura(v, c) }], ...(conConteo ? { notas: [conteo(v)] } : {}) })
            return <>
              <p className="text-small text-text-secondary">{frescura(viejo, cero)}</p>
              <Parametro etiqueta={s.pViejo} {...prop('price.stale_after_days')} validar={(v) => vEntero(1)(v) ?? (Number(v) > cero ? t.v.v25(Number(v)) : null)} vistaPrevia={(n) => previa(n, cero, true)} />
              <Parametro etiqueta={s.pCero} {...prop('price.zero_weight_after_days')} validar={(v) => vEntero(0)(v) ?? (Number(v) < viejo ? t.v.v25(viejo) : null)} vistaPrevia={(n) => previa(viejo, n, false)} />
              <Parametro etiqueta={s.pMax} {...prop('price.profeco_import_max_gb')} />
            </>
          }}</ConDatos>
        </Bloque>
      </div>
      <NuevaTienda abierta={nueva} alCambiar={setNueva} alCrear={tiendas.recargar} />
    </div>
  )
}

function Detalle({ x, cadenas, grupos, guardar, alGrupos, alCerrar }: { x: St; cadenas: St[]; grupos: Esquemas['ProductGroup'][]; guardar: (c: Esquemas['StoreUpdate']) => Promise<void>; alGrupos: () => void; alCerrar: () => void }) {
  const texto = (campo: 'profeco_name' | 'website_url' | 'address' | 'municipality' | 'state' | 'notes', etiqueta: string) => (
    <Fila etiqueta={etiqueta}><Editable etiqueta={etiqueta} tipo={campo === 'notes' ? 'area' : 'texto'} valor={x[campo] ?? ''} mostrar={x[campo] || t.comun.sinValor} validar={campo === 'profeco_name' ? (v) => (v.length > 120 ? t.v.v09 : null) : undefined} guardar={(v) => guardar({ [campo]: v.trim() || null })} /></Fila>
  )
  const veredicto = x.assisted_reading_verdict ?? 'missing'
  const [sel, setSel] = useState<number[]>(x.product_group_ids ?? [])
  return (
    <Hoja abierta alCambiar={(a) => !a && alCerrar()} titulo={x.name}>
      <div className="flex flex-col gap-1">
        <Fila etiqueta={s.nombre}><Editable etiqueta={s.nombre} valor={x.name} mostrar={x.name} validar={vNombre} guardar={(v) => guardar({ name: v.trim() })} /></Fila>
        <Fila etiqueta={s.tipo}><span className="px-2">{s.tipos[x.kind]}</span></Fila>
        {x.kind === 'branch' && <Fila etiqueta={s.cadena}><Editable etiqueta={s.cadena} tipo="opcion" valor={String(x.parent_id ?? '')} mostrar={cadenas.find((c) => c.id === x.parent_id)?.name ?? t.comun.sinValor} opciones={[{ valor: '', texto: t.v.v22 }, ...cadenas.map((c) => ({ valor: String(c.id), texto: c.name }))]} validar={vOpcion()} guardar={(v) => guardar({ parent_id: Number(v) })} /></Fila>}
        {texto('profeco_name', s.profeco)}
        {x.kind === 'chain' && texto('website_url', s.sitio)}
        {texto('address', s.direccion)}{texto('municipality', s.municipio)}{texto('state', s.entidad)}
        <Fila etiqueta={s.queCompra}>
          <fieldset className="flex flex-wrap gap-2"><legend className="sr-only">{s.queCompra}</legend>
            {grupos.map((g) => <label key={g.id} className="flex min-h-11 items-center gap-2 px-2"><input type="checkbox" className="size-5" checked={sel.includes(g.id)} onChange={async (e) => { const n = e.target.checked ? [...sel, g.id] : sel.filter((i) => i !== g.id); setSel(n); await pedir('put', '/api/v1/stores/{store_id}/product-groups', { ruta: { store_id: x.id }, cuerpo: { ids: n } }); alGrupos() }} />{g.name}</label>)}
          </fieldset>
        </Fila>
        <Fila etiqueta={s.exige}><Interruptor etiqueta={s.exige} encendido={x.requires_membership} guardar={(v) => guardar({ requires_membership: v })} /></Fila>
        {x.requires_membership && <>
          <Fila etiqueta={s.tiene}><Interruptor etiqueta={s.tiene} encendido={x.has_membership} guardar={(v) => guardar({ has_membership: v })} /></Fila>
          <Fila etiqueta={s.costo}><Editable etiqueta={s.costo} tipo="dinero" valor={textoDinero(x.membership_annual_fee_cents)} mostrar={x.membership_annual_fee_cents != null ? <span className="tabular-nums">{dinero(x.membership_annual_fee_cents)}</span> : t.comun.sinValor} validar={vDinero()} vistaPrevia={() => ({ pendiente: `${s.equilibrio} · ${t.comun.provisional}` })} guardar={(v) => guardar({ membership_annual_fee_cents: aCentavos(v) ?? 0 })} /></Fila>
          <Fila etiqueta={s.equilibrio}><span className="px-2 text-text-secondary">{x.membership_annual_fee_cents == null ? s.capturaCosto : `${t.comun.previaPendiente} · ${t.comun.provisional}`}</span></Fila>
        </>}
        {x.kind === 'chain' && <Fila etiqueta={s.lectura} nota={`${s.veredicto}: ${s.veredictos[veredicto]}`}>
          <Interruptor etiqueta={s.lectura} encendido={x.assisted_reading_enabled} bloqueadoPor={veredicto !== 'allowed' ? t.v.v19 : undefined} guardar={(v) => guardar({ assisted_reading_enabled: v })} critico={{ titulo: s.activarVerbo, verbo: s.activarVerbo, texto: s.activar(x.name) }} />
        </Fila>}
        <Fila etiqueta={t.comun.activa}><Interruptor etiqueta={t.comun.activa} encendido={x.is_active} guardar={(v) => guardar({ is_active: v })} /></Fila>
        {texto('notes', s.notas)}
      </div>
    </Hoja>
  )
}

function Fuentes({ f, ia, maxGb, alCambiar }: { f: Esquemas['PriceSources']; ia: Esquemas['AiProvider'][]; maxGb: number | null; alCambiar: () => void }) {
  const conexion = useConexion()
  const [archivo, setArchivo] = useState<File | null>(null)
  const [msg, setMsg] = useState<{ rol: 'exito' | 'error'; texto: string } | null>(null)
  const [importando, setImportando] = useState(false)
  const iaOcr = ia.some((p) => p.is_enabled && p.allowed_uses.includes('ticket_ocr'))
  const max = maxGb == null ? null : maxGb * 1024 ** 3
  const tam = (b: number) => (b >= 1024 ** 3 ? `${(b / 1024 ** 3).toFixed(2)} GB` : `${(b / 1024 ** 2).toFixed(1)} MB`)
  const importar = async () => {
    if (!archivo) return
    if (max != null && archivo.size > max) return setMsg({ rol: 'error', texto: t.v.v20(`${maxGb} GB`) })
    const datos = new FormData(); datos.append('file', archivo)
    setImportando(true); setMsg(null)
    try { const r = await pedir('post', '/api/v1/profeco-imports', { idempotencia: nuevaClave(), cuerpo: datos as never }); setMsg({ rol: 'exito', texto: s.importado(r.loaded_count, r.skipped_count) }); alCambiar() }
    catch (e) { setMsg({ rol: 'error', texto: e instanceof ErrorApi ? e.message : t.comun.noGuardo }) }
    finally { setImportando(false) }
  }
  return (
    <div className="flex flex-col">
      <Fila etiqueta={s.manual}><span className="px-2">{s.siempre}</span></Fila>
      <Fila etiqueta={s.ocr}><span className="px-2">{f.ocr_engine} · {s.conConfirmacion}</span></Fila>
      <Fila etiqueta={s.ia}><span className="flex flex-wrap items-center gap-2 px-2">{iaOcr ? s.iaEncendida : s.iaApagada}<Enlace a="/ajustes/11" className="inline-flex min-h-11 items-center text-text-link">{s.iaEnlace}</Enlace></span></Fila>
      <Fila etiqueta={s.lecturaGeneral} nota={s.lecturaGeneralNota}>
        <Interruptor etiqueta={s.lecturaGeneral} encendido={f.assisted_reading_enabled} guardar={async (v) => { await pedir('patch', '/api/v1/settings/prices', { cuerpo: { assisted_reading_enabled: v } }); alCambiar() }} critico={{ titulo: s.activarVerbo, verbo: s.activarVerbo, texto: s.activar(s.sitios) }} />
      </Fila>
      <Fila etiqueta={s.profecoTitulo} nota={s.profecoNota}>
        {f.profeco_last_import_at ? <dl className="grid grid-cols-2 gap-x-3 px-2 text-small tabular-nums"><dt className="text-text-secondary">{s.ultimoCorte}</dt><dd>{fechaHora(f.profeco_last_import_at)}</dd><dt className="text-text-secondary">{s.registros}</dt><dd>{f.profeco_price_count.toLocaleString('es-MX')}</dd>
          <dt className="text-text-secondary">{s.cobertura}</dt><dd>{f.profeco_coverage.map((c) => `${c.municipality ?? s.municipioSin}: ${c.price_count.toLocaleString('es-MX')}`).join(' · ')}</dd></dl>
          : <span className="px-2 text-text-secondary">{s.sinCortes}</span>}
      </Fila>
      <Fila etiqueta={s.importar}>
        <div className="flex flex-col gap-2">
          <label className="flex flex-col gap-1 text-small text-text-secondary">{s.archivo}<input type="file" accept=".rar,.csv" className="min-h-11 text-small" onChange={(e) => { setArchivo(e.target.files?.[0] ?? null); setMsg(null) }} /></label>
          {archivo && <p className="text-small tabular-nums">{s.archivoInfo(archivo.name, tam(archivo.size), maxGb == null ? t.comun.sinValor : `${maxGb} GB`)}</p>}
          <Boton className="w-fit" deshabilitadoPor={!conexion ? textos.conexion.requiere : !archivo ? s.archivo : importando ? s.importando : undefined} onClick={() => void importar()}><Upload aria-hidden className="size-4" />{s.importar}</Boton>
          {/* Avance de la importación (U4 P-26.06, «Estados»). La tarea en listJobs queda PROVISIONAL: no hay job de importación PROFECO (bloqueo S6-compras-b-fuentes, punto 2). */}
          <p role="status" aria-live="polite" className="text-small text-text-secondary" data-prueba="avance-profeco">{importando ? `${s.importando} · ${s.avanceTarea}` : ''}</p>
          {msg && <Aviso rol={msg.rol}>{msg.texto}</Aviso>}
        </div>
      </Fila>
    </div>
  )
}

// Calendario por día de visita, recurrente (U4 P-26.06, «Vista de base de datos»): mercados y tianguis activos en su día de la semana,
// semana de lunes a domingo con ←/→ para la semana anterior o siguiente («Atajos»). VistaCalendario de la base solo pinta fechas fijas por mes.
function CalendarioVisitas({ tiendas, alAbrir }: { tiendas: St[]; alAbrir: (id: number) => void }) {
  const [lunes, setLunes] = useState(() => lunesDe(hoyISO()))
  const mercados = tiendas.filter((x) => x.kind === 'market' && x.is_active && x.visit_weekday)
  const mover = (n: number) => setLunes((l) => sumarDias(l, 7 * n))
  const flechas = (e: EventoTecla<HTMLElement>) => {
    if ((e.key !== 'ArrowLeft' && e.key !== 'ArrowRight') || e.altKey || e.ctrlKey || e.metaKey || (e.target as HTMLElement).closest('input, select, textarea, [contenteditable="true"]')) return
    e.preventDefault(); mover(e.key === 'ArrowLeft' ? -1 : 1)
  }
  return (
    <section aria-labelledby="t-calendario-visitas" className="mt-4 flex flex-col gap-2" onKeyDown={flechas} data-prueba="calendario-visitas">
      <div className="flex flex-wrap items-center gap-2">
        <h4 id="t-calendario-visitas" className="text-body font-semibold">{s.calendario}</h4>
        <span className="flex items-center gap-1">
          <Boton tam="icono" variante="fantasma" className="shrink-0" aria-label={s.semanaAnterior} onClick={() => mover(-1)}><ChevronLeft aria-hidden className="size-6" /></Boton>
          <span className="text-small tabular-nums" aria-live="polite">{s.semana(fecha(lunes), fecha(sumarDias(lunes, 6)))}</span>
          <Boton tam="icono" variante="fantasma" className="shrink-0" aria-label={s.semanaSiguiente} onClick={() => mover(1)}><ChevronRight aria-hidden className="size-6" /></Boton>
        </span>
      </div>
      {!mercados.length ? <p className="text-small text-text-secondary">{s.sinMercados}</p> : (
        <ol className="grid grid-cols-1 gap-px rounded-md border border-border-divider bg-border-divider text-small desktop:grid-cols-7">
          {s.dias.map((d, i) => { const f = sumarDias(lunes, i); const del = mercados.filter((x) => x.visit_weekday === i + 1); return (
            <li key={f} aria-label={`${d} ${fecha(f)}`} className="flex flex-col gap-1 bg-surface-card p-2 desktop:min-h-20">
              <span className="text-caption tabular-nums text-text-secondary">{d} {fecha(f)}</span>
              {del.map((x) => <button type="button" key={x.id} onClick={() => alAbrir(x.id)} className="min-h-11 rounded-sm bg-surface-hover px-2 text-left text-small text-text-default desktop:min-h-9">{x.name}</button>)}
            </li>) })}
        </ol>)}
    </section>
  )
}

function NuevaTienda({ abierta, alCambiar, alCrear }: { abierta: boolean; alCambiar: (a: boolean) => void; alCrear: () => void }) {
  const conexion = useConexion()
  const [v, setV] = useState({ kind: '', name: '', dia: '' })
  const [err, setErr] = useState<Record<string, string | null>>({})
  const [clave] = useState(nuevaClave)
  const guardar = async () => {
    const e = { kind: vOpcion()(v.kind), name: vNombre(v.name), dia: v.kind === 'market' && !v.dia ? t.v.v11 : null }; setErr(e)
    if (Object.values(e).some(Boolean)) return
    try { await pedir('post', '/api/v1/stores', { idempotencia: clave, cuerpo: { kind: v.kind as St['kind'], name: v.name.trim(), ...(v.dia ? { visit_weekday: Number(v.dia) } : {}), is_active: true } }); setV({ kind: '', name: '', dia: '' }); alCambiar(false); alCrear() }
    catch (x) { setErr({ name: x instanceof ErrorApi && x.estado === 409 ? t.v.v12 : t.comun.noGuardo }) }
  }
  return (
    <Hoja abierta={abierta} alCambiar={alCambiar} titulo={s.nueva} pie={<div className="flex justify-end gap-2"><Boton onClick={() => alCambiar(false)}>{textos.app.cancelar}</Boton><Boton variante="primario" deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => void guardar()}>{textos.app.guardar}</Boton></div>}>
      <form className="flex flex-col gap-3" onSubmit={(e) => { e.preventDefault(); void guardar() }} onKeyDown={(e) => { if (e.key === 'Enter' && e.ctrlKey) { e.preventDefault(); void guardar() } }}>
        <label className="flex flex-col gap-1 text-small font-medium text-text-secondary">{s.tipo}<select className={estilosControl} value={v.kind} onChange={(e) => setV({ ...v, kind: e.target.value })}><option value="">{t.v.v22}</option>{Object.entries(s.tipos).map(([valor, texto]) => <option key={valor} value={valor}>{texto}</option>)}</select>{err.kind && <span className="text-role-error">{err.kind}</span>}</label>
        <Campo etiqueta={s.nombre} value={v.name} error={err.name ?? undefined} onChange={(e) => setV({ ...v, name: e.target.value })} />
        {v.kind === 'market' && <label className="flex flex-col gap-1 text-small font-medium text-text-secondary">{s.diaVisita}<select className={estilosControl} value={v.dia} onChange={(e) => setV({ ...v, dia: e.target.value })}><option value="">{t.v.v11}</option>{dias.map((d) => <option key={d.valor} value={d.valor}>{d.texto}</option>)}</select>{err.dia && <span className="text-role-error">{err.dia}</span>}</label>}
      </form>
    </Hoja>
  )
}
