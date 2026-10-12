// P-26.12 Personas y reparto (U4 parte b; decisión A, B, K13, K14, RN-03, Q2-68 C). Las partes se guardan juntas
// tras el diálogo de alcance (scope onward|all de C1 v1.3). Preselección PROVISIONAL (bloqueo U4-reparto-preseleccion);
// el conteo y la fecha de «también los ya capturados», «Vista previa pendiente» (bloqueo C1-reparto-alcance).
import { useState } from 'react'
import { CircleAlert } from 'lucide-react'
import { pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import { aCentavos, dinero } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { estilosControl } from '../../componentes/Campo'
import { ChipPersona } from '../../componentes/Chip'
import { Hoja } from '../../componentes/Hoja'
import { Esqueleto } from '../../componentes/Estados'
import { cn } from '../../componentes/utils'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { AntesDespues, Bloque, Cargar, IndicadorGuardado, Marco, unir, useGuardado, useSinRed } from './comun'
import { aPuntosBase, divHalfUp, pb } from './calculos'
import { textosAjustesB as T } from './textos'

const t = T.p
type Pe = Esquemas['Person']; type Regla = Esquemas['CostShareRule']
export const invalido = (m: string): never => { throw new ErrorApi(422, 'validation_failed', m) }
const centavosATexto = (c?: number | null) => (c === undefined || c === null ? '' : (c / 100).toFixed(2))
const actualizarPersona = (id: number, cuerpo: Esquemas['PersonUpdate']) => pedir('patch', '/api/v1/persons/{person_id}', { ruta: { person_id: id }, cuerpo })

export default function P2612Personas() {
  const datos = unir<[{ items: Pe[] }, { items: Regla[] }]>([
    usarConsulta('get', '/api/v1/persons', { consulta: { limit: 50 } }),
    usarConsulta('get', '/api/v1/cost-share-rules', { consulta: { limit: 50 } }),
  ])
  return (
    <Marco id="12">
      <p className="text-body text-text-secondary">{t.proposito}</p>
      <Cargar c={datos} que={t.que} esqueleto={<div className="flex flex-col gap-3"><Esqueleto className="h-32" /><Esqueleto className="h-32" /></div>}>{([ps, rs]) => (
        <>
          <Bloque titulo={t.personas} ayuda={t.efectoTope('$X').replace('$X', '…')}>
            <VistaBD def={{
              clave: 'P-26.12-personas', filas: ps.items, idFila: (p) => p.id, vistas: ['tabla'], titulo: (p) => p.display_name,
              propiedades: [
                { id: 'display_name', nombre: t.nombre, valor: (p) => p.display_name, pintar: (p) => <ChipPersona codigo={p.code} nombre={p.display_name} />,
                  editar: { tipo: 'texto', textoInicial: (p) => p.display_name, guardar: async (p, v) => { if (!v.trim()) invalido(T.validar.nombre); await actualizarPersona(p.id, { display_name: v.trim() }); datos.recargar() } } },
                { id: 'participation', nombre: t.papel, valor: (p) => t.papeles[p.participation] },
                { id: 'has_account', nombre: t.cuenta, valor: (p) => ((p.has_account ?? p.participation !== 'external') ? T.comun.si : T.comun.no) },
                { id: 'pays', nombre: t.pagaDirecto, valor: (p) => (p.pays_share_directly ? T.comun.si : T.comun.no) },
                { id: 'hace', nombre: t.queHace, valor: (p) => t.hace[p.participation] },
                { id: 'cap', nombre: t.tope, valor: (p) => p.fortnight_payment_cap_cents ?? null, alinear: 'derecha', pintar: (p) => <span className="tabular-nums">{p.fortnight_payment_cap_cents === undefined || p.fortnight_payment_cap_cents === null ? t.sinTope : dinero(p.fortnight_payment_cap_cents)}</span>,
                  editar: {
                    tipo: 'dinero', textoInicial: (p) => centavosATexto(p.fortnight_payment_cap_cents),
                    vistaPrevia: async (p, v) => {
                      const c = v.trim() === '' ? null : aCentavos(v)
                      const antes = p.fortnight_payment_cap_cents ?? null
                      return [{ etiqueta: t.tope, antes: antes === null ? t.sinTope : dinero(antes), despues: c === null ? (v.trim() ? T.validar.monto : t.sinTope) : `${dinero(c)} · ${t.efectoTope(dinero(c))}` }]
                    },
                    guardar: async (p, v) => { const c = v.trim() === '' ? null : aCentavos(v); if (v.trim() && (c === null || c < 0)) invalido(T.validar.monto); await actualizarPersona(p.id, { fortnight_payment_cap_cents: c }); datos.recargar() },
                  } },
                { id: 'notes', nombre: t.notas, valor: (p) => p.notes ?? '', pintar: (p) => p.notes || T.comun.sinValor, editar: { tipo: 'texto', textoInicial: (p) => p.notes ?? '', guardar: async (p, v) => { await actualizarPersona(p.id, { notes: v.trim() || null }); datos.recargar() } } },
              ],
            }} />
          </Bloque>
          <Bloque titulo={t.reglas}>
            <div className="grid grid-cols-1 gap-3 desktop:grid-cols-3">{rs.items.filter((r) => r.is_template).map((r) => <TarjetaRegla key={r.id} r={r} personas={ps.items} recargar={datos.recargar} />)}</div>
          </Bloque>
        </>
      )}</Cargar>
    </Marco>
  )
}

const COLOR: Record<string, string> = { axel: 'bg-person-axel', ximena: 'bg-person-ximena', karen: 'bg-person-karen', shared: 'bg-person-shared' }
const nombre = (ps: Pe[], id: number) => ps.find((p) => p.id === id)?.display_name ?? ''
function resumenRegla(r: { parts: { person_id: number; pct_bp?: number; fixed_cents?: number }[]; remainder_person_id: number }, ps: Pe[]) {
  const pct = r.parts.reduce((s, p) => s + (p.pct_bp ?? 0), 0)
  const partes = r.parts.filter((p) => p.person_id !== r.remainder_person_id).map((p) => `${nombre(ps, p.person_id)} ${p.pct_bp !== undefined ? pb(p.pct_bp) : dinero(p.fixed_cents ?? 0)}`)
  return [...partes, t.elResto(pb(Math.max(0, 10000 - pct)), `${nombre(ps, r.remainder_person_id)}, ${t.hogar}`)].join(' · ')
}
function TarjetaRegla({ r, personas, recargar }: { r: Regla; personas: Pe[]; recargar: () => void }) {
  const sinRed = useSinRed()
  const g = useGuardado()
  const [nombreRegla, setNombre] = useState(r.name)
  const [abierta, setAbierta] = useState(false)
  const guardar = (cuerpo: Esquemas['CostShareRuleUpdate']) => g.ejecutar(async () => { await pedir('patch', '/api/v1/cost-share-rules/{cost_share_rule_id}', { ruta: { cost_share_rule_id: r.id }, cuerpo }); recargar() })
  const partes = (r.parts ?? []).filter((p) => p.person_id !== r.remainder_person_id)
  const pct = partes.reduce((s, p) => s + (p.pct_bp ?? 0), 0)
  return (
    <article className="flex flex-col gap-2 rounded-md border border-border-divider p-3" data-prueba={`regla-${r.id}`}>
      <label className="flex flex-col gap-1 text-caption text-text-secondary">{T.s.nombre}
        <input className={estilosControl} value={nombreRegla} disabled={!!sinRed} onChange={(e) => setNombre(e.target.value)}
          onBlur={() => { if (nombreRegla === r.name) return; if (!nombreRegla.trim()) { setNombre(r.name); return } void guardar({ name: nombreRegla.trim() }) }}
          onKeyDown={(e) => { if (e.key === 'Enter') e.currentTarget.blur(); else if (e.key === 'Escape') setNombre(r.name) }} />
      </label>
      {r.is_default && <span className="self-start rounded-sm bg-role-info-surface px-2 py-1 text-caption">{t.porDefecto}</span>}
      <div aria-hidden className="flex h-2 overflow-hidden rounded-full bg-viz-track">
        {partes.map((p) => <span key={p.person_id} className={cn('h-full', COLOR[personas.find((x) => x.id === p.person_id)?.code ?? ''] ?? 'bg-person-shared')} style={{ width: `${(p.pct_bp ?? 0) / 100}%` }} />)}
        <span className={cn('h-full', COLOR[personas.find((x) => x.id === r.remainder_person_id)?.code ?? ''] ?? 'bg-person-shared')} style={{ width: `${Math.max(0, 10000 - pct) / 100}%` }} />
      </div>
      <p className="text-small tabular-nums">{`${r.name}: ${resumenRegla({ parts: r.parts ?? [], remainder_person_id: r.remainder_person_id }, personas)}`}</p>
      <p className="text-caption text-text-secondary">{`${t.resto}: ${nombre(personas, r.remainder_person_id)}`}</p>
      <Boton deshabilitadoPor={sinRed} onClick={() => setAbierta(true)}>{t.editarPartes}</Boton>
      <IndicadorGuardado g={g} />
      {abierta && <EditorPartes r={r} personas={personas} cerrar={() => setAbierta(false)} guardar={guardar} />}
    </article>
  )
}
type Borrador = { person_id: number; modo: 'pct' | 'fijo'; texto: string }
function EditorPartes({ r, personas, cerrar, guardar }: { r: Regla; personas: Pe[]; cerrar: () => void; guardar: (c: Esquemas['CostShareRuleUpdate']) => Promise<boolean> }) {
  const otras = personas.filter((p) => p.id !== r.remainder_person_id)
  const [b, setB] = useState<Borrador[]>(otras.map((p) => { const x = r.parts?.find((q) => q.person_id === p.id); return { person_id: p.id, modo: x?.fixed_cents !== undefined ? 'fijo' : 'pct', texto: x ? (x.fixed_cents !== undefined ? (x.fixed_cents / 100).toFixed(2) : (x.pct_bp! / 100).toFixed(2)) : '' } }))
  const [monto, setMonto] = useState('100.00')
  const [alcance, setAlcance] = useState<'onward' | 'all' | null>(null)
  const partes = b.filter((x) => x.texto.trim() !== '').map((x) => (x.modo === 'pct' ? { person_id: x.person_id, pct_bp: aPuntosBase(x.texto) } : { person_id: x.person_id, fixed_cents: aCentavos(x.texto) }))
  const malas = partes.some((p) => ('pct_bp' in p ? p.pct_bp == null || p.pct_bp > 10000 : p.fixed_cents == null || p.fixed_cents < 0))
  const suma = partes.reduce((s, p) => s + ('pct_bp' in p ? p.pct_bp ?? 0 : 0), 0)
  const quien = nombre(personas, r.remainder_person_id)
  const error = malas ? t.porcentajeMal : suma > 10000 ? t.v16(pb(suma).replace(' %', ''), quien) : ''
  const M = aCentavos(monto) ?? 0
  const reparte = (ps: { person_id: number; pct_bp?: number | null; fixed_cents?: number | null }[]) => {
    const l = ps.map((p) => ({ n: nombre(personas, p.person_id), c: p.pct_bp != null ? divHalfUp(M * p.pct_bp, 10000) : p.fixed_cents ?? 0 }))
    const resto = M - l.reduce((s, x) => s + x.c, 0)
    return [...l.map((x) => `${x.n} ${dinero(x.c)}`), `${quien}, ${t.hogar} ${dinero(resto)}`].join(' y ')
  }
  const confirmar = async () => {
    if (!alcance) return
    const ok = await guardar({ parts: partes.map((p) => ('pct_bp' in p ? { person_id: p.person_id, pct_bp: p.pct_bp! } : { person_id: p.person_id, fixed_cents: p.fixed_cents! })), scope: alcance })
    if (ok) cerrar()
  }
  return (
    <>
      <Hoja abierta={alcance === null} alCambiar={(a) => !a && cerrar()} titulo={`${t.editarPartes} · ${r.name}`}
        pie={<Boton variante="primario" tam="hoja" deshabilitadoPor={error || undefined} onClick={() => setAlcance('onward')}>{T.comun.guardar}</Boton>}>
        <div className="flex flex-col gap-3">
          {b.map((x, i) => (
            <fieldset key={x.person_id} className="flex flex-wrap items-end gap-2">
              <legend className="text-small font-medium">{nombre(personas, x.person_id)}</legend>
              <select aria-label={`${t.partes} · ${nombre(personas, x.person_id)}`} className={cn(estilosControl, 'w-36')} value={x.modo} onChange={(e) => setB(b.map((y, j) => (j === i ? { ...y, modo: e.target.value as 'pct' | 'fijo' } : y)))}>
                <option value="pct">%</option><option value="fijo">{t.montoFijo}</option>
              </select>
              <input aria-label={`${nombre(personas, x.person_id)} ${x.modo === 'pct' ? '%' : t.montoFijo}`} inputMode="decimal" className={cn(estilosControl, 'w-32 tabular-nums')} value={x.texto} onChange={(e) => setB(b.map((y, j) => (j === i ? { ...y, texto: e.target.value } : y)))} />
            </fieldset>
          ))}
          <p className="text-small tabular-nums">{t.elResto(pb(Math.max(0, 10000 - suma)), `${quien}, ${t.hogar}`)}</p>
          {error && <p role="alert" className="inline-flex items-center gap-1 text-small text-role-error"><CircleAlert aria-hidden className="size-5" />{error}</p>}
          <label className="flex items-center gap-2 text-small">{t.conGasto}<input inputMode="decimal" className={cn(estilosControl, 'w-32 tabular-nums')} value={monto} onChange={(e) => setMonto(e.target.value)} /></label>
          {!error && <AntesDespues cambios={[{ etiqueta: t.efectoParte(dinero(M), '').replace(/: $/, ''), antes: reparte(r.parts ?? []), despues: reparte(partes) }]} />}
        </div>
      </Hoja>
      <Hoja abierta={alcance !== null} alCambiar={(a) => !a && setAlcance(null)} titulo={t.alcance}
        pie={<div className="flex flex-col gap-2"><Boton variante="primario" tam="hoja" onClick={() => void confirmar()}>{T.n.confirmar}</Boton><Boton tam="hoja" onClick={() => setAlcance(null)}>{T.comun.cancelar}</Boton></div>}>
        <div role="radiogroup" aria-label={t.alcance} className="flex flex-col gap-3">
          {(['onward', 'all'] as const).map((v) => (
            <label key={v} className="flex min-h-11 items-start gap-2 rounded-md border border-border-control p-3 text-small has-[:checked]:border-border-selected has-[:checked]:bg-surface-selected">
              <input type="radio" name="alcance" className="mt-1 size-5 accent-action-primary-bg" checked={alcance === v} onChange={() => setAlcance(v)} />
              <span><span className="block font-medium">{v === 'onward' ? t.soloNuevos : t.tambien}</span><span className="text-text-secondary">{v === 'onward' ? t.soloNuevosAyuda : t.tambienAyuda}</span></span>
            </label>
          ))}
          <p className="text-caption text-text-secondary">{`${t.preseleccion}: ${t.soloNuevos}`}</p>
        </div>
      </Hoja>
    </>
  )
}
