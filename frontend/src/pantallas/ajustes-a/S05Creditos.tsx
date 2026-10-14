// P-26.05 Créditos y condiciones (U4 parte a): líneas y condiciones con vistas previas 2 y 4–10, interruptor A/B con confirmación,
// tasas por plazo, calendario de promociones y comparador (restricciones K16, pesos, parámetros K44).
import { useEffect, useState } from 'react'
import { Plus } from 'lucide-react'
import { nuevaClave, pedir } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { usarConsulta } from '../../api/consulta'
import { ErrorApi } from '../../api/errores'
import { aCentavos, dinero, fecha, hoyISO, porcentaje } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { Enlace } from '../../app/enrutador'
import { Boton } from '../../componentes/Boton'
import { Campo } from '../../componentes/Campo'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Hoja } from '../../componentes/Hoja'
import { textos } from '../../textos/es-MX'
import type { Cambio } from '../../vistas-bd/tipos'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { aFraccion, aPuntosBase, corteYPago, multDiv, normalizar, pctDeFraccion, pctDePuntosBase } from './calculos'
import { Aviso, Bloque, DialogoCritico, Editable, ErrorGuardado, Fila, Parametro, textoDinero, textoFallo, vDia, vDinero, vEntero, vFecha, vNombre, vOpcion, vPct, vPlazo, type Previa } from './comunes'
import { comoLista } from './comunes'
import { guardarVarias, nota as notaGrupo, resumen, siNo } from './grupo'
import { t } from './textos'
const k = t.creditos
type Cr = Esquemas['Credit']; type Tasa = Esquemas['CreditTermRate']; type Promo = Esquemas['Promotion']; type Cal = Esquemas['CreditPurchasePreview']; type Compra = Esquemas['CreditPurchase']
const CAPITAL_EJEMPLO = 309450
const opciones = (m: Record<string, string>) => Object.entries(m).map(([valor, texto]) => ({ valor, texto }))
const tasaDe = (num?: number | null, den?: number | null, dec?: number | null, pct?: string) => (pct ? `${pct} %`.replace(/ % %$/, ' %') : num != null && den ? pctDeFraccion(num, den, dec ?? 2) : t.comun.sinValor)
const textoTasa = (num?: number | null, den?: number | null, dec?: number | null) => (num != null && den ? pctDeFraccion(num, den, dec ?? 2).replace(' %', '') : '')

/** Faltantes de K18 / RN-26 calculados en el navegador (vista previa 4 (a)). */
function faltantes(c: Partial<Cr>): string[] {
  if (c.kind !== 'bank_card' && c.kind !== 'store_credit') return []
  const f: string[] = []
  if (c.annual_rate_num == null) f.push('annual_rate'); if (c.rate_base == null) f.push('rate_base'); if (c.vat_included == null) f.push('vat')
  if (c.min_payment_pct_bp == null) f.push('min_payment_pct_bp'); if (c.min_payment_fixed_cents == null) f.push('min_payment_fixed_cents')
  return f
}
const textoFaltan = (f: string[]) => (f.length ? f.map((x) => k.faltan[x]).join(', ') : k.sinFaltantes)

/** Intereses de la compra de ejemplo (vistas previas 6 y 7): línea revolvente, con la tabla por plazo (RN-08, compra planeada, igual que la vista previa 4 de P-26.02);
 *  otras líneas, los de su última compra (PROVISIONAL, bloqueo WF3-ajustes-compra-ejemplo). */
function interesEjemplo(revolvente: boolean, capital: number, plazo: number, u?: Compra): Partial<Esquemas['CreditPurchaseCreate']> {
  if (revolvente) return { interest_source: 'term_table' }
  if (!u) return {}
  return { interest_source: u.interest_source, rate_basis: u.rate_basis, rate_num: u.rate_num, rate_den: u.rate_den, rate_decimals: u.rate_decimals, interest_includes_vat: u.interest_includes_vat,
    ...(capital === u.principal_cents && plazo === u.term_count ? { interest_cents: u.interest_cents } : {}) }
}
const fechasDe = (p: Cal) => p.schedule.map((x) => fecha(x.due_date)).join(' · ')
const pagosDe = (p: Cal) => { const m = p.schedule.map((x) => x.amount_cents); return m.length && m.every((x) => x === m[0]) ? k.cadaPago(dinero(m[0]), m.length) : m.map((x) => dinero(x)).join(' · ') }
/** Calendario de la compra de ejemplo «Antes → Después»; sin `B`, el «Después» queda «Vista previa pendiente». */
const resumenCal = (A: Cal, B: Cal | null): Cambio[] => [
  { etiqueta: k.fechasPago, antes: fechasDe(A), despues: B ? fechasDe(B) : t.comun.previaPendiente },
  { etiqueta: k.pagoCada, antes: pagosDe(A), despues: B ? pagosDe(B) : t.comun.previaPendiente },
  { etiqueta: k.totalPagos, antes: dinero(A.total_cents), despues: B ? dinero(B.total_cents) : t.comun.previaPendiente },
]

/** Tabla francesa con fracciones exactas (S3-interfaz §5.4, Q2-28 (a)): pago fijo y total de intereses. */
export function amortizar(capital: number, num: number, den: number, periodosAnio: number, n: number, ivaBp?: number) {
  let a = BigInt(num), b = BigInt(den) * BigInt(periodosAnio)
  if (ivaBp) { a *= BigInt(10000 + ivaBp); b *= 10000n }
  const C = BigInt(capital), u = (a + b) ** BigInt(n), bn = b ** BigInt(n)
  const N = C * a * u, D = b * (u - bn)
  const pago = D > 0n ? Number((2n * N + D) / (2n * D)) : Math.round(capital / n)
  let saldo = capital, intereses = 0
  for (let i = 0; i < n; i++) { const int = multDiv(saldo, Number(a), Number(b)); intereses += int; saldo = i === n - 1 ? 0 : saldo - (pago - int) }
  return { pago, intereses }
}

export default function S05Creditos() {
  const creditos = comoLista<Esquemas['Credit']>(usarConsulta('get', '/api/v1/credits', { consulta: { limit: 200 } }))
  const promos = comoLista<Esquemas['Promotion']>(usarConsulta('get', '/api/v1/promotions', { consulta: { limit: 200, sort: 'valid_from' } }))
  const optimizador = usarConsulta('get', '/api/v1/settings/optimizer')
  const [abierto, setAbierto] = useState<number | null>(null)
  const [nuevo, setNuevo] = useState(false)
  const [nuevaPromo, setNuevaPromo] = useState(false)
  // Alt+N: nuevo registro en el bloque con foco (U4 P-26.05 «Atajos»): promoción dentro del calendario; línea en lo demás.
  useEffect(() => { const f = (e: KeyboardEvent) => { if (e.altKey && e.key.toLowerCase() === 'n') { e.preventDefault(); if ((document.activeElement as HTMLElement | null)?.closest('section[aria-labelledby="b-promos"]')) setNuevaPromo(true); else setNuevo(true) } }; window.addEventListener('keydown', f); return () => window.removeEventListener('keydown', f) }, [])
  return (
    <div className="flex flex-col gap-4">
      <Bloque titulo={k.lineas} id="b-lineas" accion={<Boton variante="primario" onClick={() => setNuevo(true)}><Plus aria-hidden className="size-5" />{k.nuevo}</Boton>}>
        <ConDatos consulta={creditos}>{(d) => d.items.length === 0 ? <Vacio texto={k.sin} accion={<Boton variante="primario" onClick={() => setNuevo(true)}>{k.nuevo}</Boton>} /> : (
          <VistaBD def={{ clave: 'P-26.05-lineas', filas: d.items, idFila: (c) => c.id, vistas: ['tabla', 'galeria', 'tablero'], titulo: (c) => c.name, alAbrir: (c) => setAbierto(c.id),
            columnasTablero: { propiedad: 'kind', columnas: Object.values(k.tipos) },
            propiedades: [
              { id: 'name', nombre: k.nombre, valor: (c) => c.name, pintar: (c) => c.name },
              { id: 'kind', nombre: k.tipo, valor: (c) => c.kind, grupo: (c) => k.tipos[c.kind], pintar: (c) => k.tipos[c.kind] },
              { id: 'live_debt_cents', nombre: k.deuda, alinear: 'derecha', valor: (c) => c.live_debt_cents ?? 0, pintar: (c) => <span className="tabular-nums">{dinero(c.live_debt_cents ?? 0)}</span> },
              { id: 'available_cents', nombre: k.disponible, alinear: 'derecha', valor: (c) => c.available_cents ?? 0, pintar: (c) => <span className="tabular-nums">{c.available_cents != null ? dinero(c.available_cents) : '—'}</span> },
              { id: 'next_payment', nombre: k.proximoPago, valor: (c) => c.next_payment_date, pintar: (c) => (c.next_payment_date ? <span className="tabular-nums">{fecha(c.next_payment_date)} · {dinero(c.next_payment_cents ?? 0)}</span> : '—') },
              { id: 'missing', nombre: k.faltaCapturar, valor: (c) => (c.missing_terms ?? []).length, pintar: (c) => ((c.missing_terms ?? []).length ? <span className="text-role-warning">{k.faltaCapturar}: {textoFaltan(c.missing_terms ?? [])}</span> : '—') },
            ] }} />
        )}</ConDatos>
        {abierto != null && creditos.datos && <Linea c={creditos.datos.items.find((x) => x.id === abierto)!} alCambiar={creditos.recargar} alCerrar={() => setAbierto(null)} />}
      </Bloque>
      <Bloque titulo={k.promociones} id="b-promos" accion={<Boton variante="primario" onClick={() => setNuevaPromo(true)}><Plus aria-hidden className="size-5" />{k.nuevaPromo}</Boton>}>
        <ConDatos consulta={promos}>{(d) => <Promociones filas={d.items} creditos={creditos.datos?.items ?? []} alCambiar={promos.recargar} alNueva={() => setNuevaPromo(true)} />}</ConDatos>
      </Bloque>
      <Bloque titulo={k.comparador} id="b-comparador">
        <ConDatos consulta={optimizador}>{(o) => <Comparador o={o} alCambiar={optimizador.recargar} />}</ConDatos>
      </Bloque>
      <NuevoCredito abierta={nuevo} alCambiar={setNuevo} alCrear={creditos.recargar} />
      <NuevaPromocion abierta={nuevaPromo} alCambiar={setNuevaPromo} creditos={creditos.datos?.items ?? []} alCrear={promos.recargar} />
    </div>
  )
}

function Linea({ c, alCambiar, alCerrar }: { c: Cr; alCambiar: () => void; alCerrar: () => void }) {
  const hoy = hoyISO()
  const guardar = async (cambio: Esquemas['CreditUpdate']) => { await pedir('patch', '/api/v1/credits/{credit_id}', { ruta: { credit_id: c.id }, cuerpo: cambio }); alCambiar() }
  const revolvente = c.kind === 'revolving_line', tarjeta = c.kind === 'bank_card' || c.kind === 'store_credit'
  const notaComision: Previa = { notas: [k.notaComision] }
  // Compra de ejemplo de las vistas previas 6 y 7 (U4 P-26.05): capital y plazo de la última compra de la línea, editables; sin compras, vacíos.
  const compras = comoLista<Compra>(usarConsulta('get', '/api/v1/credit-purchases', { consulta: { credit_id: [c.id], limit: 200 } }))
  const ultima = [...(compras.datos?.items ?? [])].sort((a, b) => b.purchase_date.localeCompare(a.purchase_date) || b.id - a.id)[0]
  const [ej, setEj] = useState<{ capital: string; plazo: string } | null>(null)
  const ejemplo = ej ?? { capital: ultima ? textoDinero(ultima.principal_cents) : '', plazo: ultima ? String(ultima.term_count) : '' }
  const fechaEj = ultima?.purchase_date ?? hoy
  const faltaEjemplo = !!(vDinero(true)(ejemplo.capital) ?? vPlazo(ejemplo.plazo))
  const calendario = (frecuencia: Cr['payment_frequency'], compra: string) => {
    const capital = aCentavos(ejemplo.capital) ?? 0, plazo = Number(ejemplo.plazo)
    return pedir('post', '/api/v1/credit-purchases/preview', { cuerpo: { credit_id: c.id, kind: ultima?.kind ?? 'installment_plan', purchase_date: compra, principal_cents: capital, term_count: plazo, payment_frequency: frecuencia, ...interesEjemplo(c.kind === 'revolving_line', capital, plazo, ultima) } })
  }
  const notasEjemplo = () => [k.ejemploTexto(dinero(aCentavos(ejemplo.capital) ?? 0), Number(ejemplo.plazo), fecha(fechaEj)), ...(c.kind === 'revolving_line' ? [] : [k.ejemploProvisional])]
  const sinEjemplo: Previa = { notas: [k.escribeEjemplo], bloquea: true }
  // Vista previa 6: la misma compra con la frecuencia guardada («Antes») y con la tecleada («Después»).
  const previaFrecuencia = async (s: string): Promise<Previa> => {
    if (faltaEjemplo) return sinEjemplo
    const [A, B] = await Promise.all([calendario(c.payment_frequency, fechaEj), calendario(s as Cr['payment_frequency'], fechaEj)])
    return { cambios: resumenCal(A, B), notas: notasEjemplo() }
  }
  // Vista previa 7 (PROVISIONAL, SC-U4-1): el motor toma los días mínimos de la línea guardada; se pinta con el valor guardado y el «Después» queda pendiente.
  const previaDias = async (): Promise<Previa> => {
    if (faltaEjemplo) return sinEjemplo
    const [H, E] = await Promise.all([calendario(c.payment_frequency, hoy), calendario(c.payment_frequency, fechaEj)])
    return { cambios: [{ etiqueta: k.primerPagoHoy, antes: H.schedule[0] ? fecha(H.schedule[0].due_date) : '—', despues: t.comun.previaPendiente }, ...resumenCal(E, null)], notas: notasEjemplo(), pendiente: `${t.comun.provisional} (SC-U4-1)` }
  }
  const editorEjemplo = (
    <fieldset className="flex flex-wrap items-end gap-2 rounded-md border border-border-divider p-2" data-prueba="compra-ejemplo">
      <legend className="px-1 text-small text-text-secondary">{k.ejemplo}</legend>
      <Campo etiqueta={k.ejCapital} inputMode="decimal" className="w-40" value={ejemplo.capital} onChange={(e) => setEj({ ...ejemplo, capital: e.target.value })} />
      <Campo etiqueta={k.ejPlazo} inputMode="numeric" className="w-32" value={ejemplo.plazo} onChange={(e) => setEj({ ...ejemplo, plazo: e.target.value })} />
    </fieldset>
  )
  const depEjemplo = `${ejemplo.capital}|${ejemplo.plazo}|${fechaEj}`
  const previaCondicion = (cambio: Partial<Cr>): Previa => {
    const d = { ...c, ...cambio }
    return { cambios: [{ etiqueta: k.faltaCapturar, antes: textoFaltan(faltantes(c)), despues: textoFaltan(faltantes(d)) }, { etiqueta: k.tasaAnual, antes: tasaDe(c.annual_rate_num, c.annual_rate_den, c.annual_rate_decimals), despues: tasaDe(d.annual_rate_num, d.annual_rate_den, d.annual_rate_decimals) }],
      pendiente: `${k.pagoMinPend} · ${t.comun.provisional} (S3-formula-tarjeta-pdf)` }
  }
  const previaLimite = async (s: string): Promise<Previa> => {
    const delta = (aCentavos(s) ?? 0) - (c.credit_limit_cents ?? 0)
    if (!revolvente) return { cambios: [{ etiqueta: k.disponible, antes: dinero(c.available_cents ?? 0), despues: dinero((c.available_cents ?? 0) + delta) }] }
    const a = await pedir('get', '/api/v1/credits/{credit_id}/availability-preview', { ruta: { credit_id: c.id } })
    const orden = c.availability_reading === 'A' ? (['A', 'B'] as const) : (['B', 'A'] as const)
    const cambios: Cambio[] = orden.flatMap((l) => { const r = l === 'A' ? a.reading_a : a.reading_b; return [{ etiqueta: `${k.abLimite(l)} · ${k.hoy}`, antes: dinero(r.available_today_cents), despues: dinero(r.available_today_cents + delta) }, { etiqueta: `${k.abLimite(l)} · ${k.minimo}`, antes: dinero(r.horizon_min_cents), despues: dinero(r.horizon_min_cents + delta) }] })
    return { cambios }
  }
  const previaTasaPrestamo = async (s: string): Promise<Previa> => {
    const f = aFraccion(s)
    // R2-05: la compra vigente más reciente del préstamo (filtro status=active, orden -purchase_date) con SU frecuencia de pagos.
    const compras = await pedir('get', '/api/v1/credit-purchases', { consulta: { credit_id: [c.id], status: 'active', sort: '-purchase_date,-id', limit: 1 } })
    const x = (compras.items as Esquemas['CreditPurchase'][])[0]
    const periodos = { weekly: 52, fortnightly: 24, monthly: 12 }[x?.payment_frequency ?? c.payment_frequency]
    const iva = c.vat_included === false ? c.vat_on_interest_bp : undefined
    const cambios: Cambio[] = [{ etiqueta: k.tasaAnual, antes: tasaDe(c.annual_rate_num, c.annual_rate_den, c.annual_rate_decimals), despues: f ? pctDeFraccion(f.num, f.den, f.decimales) : '—' }]
    if (x && f && c.annual_rate_num != null && c.annual_rate_den) {
      const A = amortizar(x.principal_cents, c.annual_rate_num, c.annual_rate_den, periodos, x.term_count, iva), B = amortizar(x.principal_cents, f.num, f.den, periodos, x.term_count, iva)
      cambios.push({ etiqueta: k.proximoPago, antes: dinero(A.pago), despues: dinero(B.pago) }, { etiqueta: k.intereses.split(' de una')[0], antes: dinero(A.intereses), despues: dinero(B.intereses) })
      return { cambios, notas: [k.ejemploTexto(dinero(x.principal_cents), x.term_count, fecha(x.first_payment_date ?? x.purchase_date))] }
    }
    return { cambios }
  }
  const tasaCampo = (etiqueta: string, previa: (s: string) => Previa | Promise<Previa>) => (
    <Editable etiqueta={etiqueta} tipo="numero" valor={textoTasa(c.annual_rate_num, c.annual_rate_den, c.annual_rate_decimals)} mostrar={<span className="tabular-nums">{tasaDe(c.annual_rate_num, c.annual_rate_den, c.annual_rate_decimals, c.annual_rate_pct)}</span>} validar={vPct(4)}
      vistaPrevia={previa} guardar={(s) => { const f = aFraccion(s)!; return guardar({ annual_rate_num: f.num, annual_rate_den: f.den, annual_rate_decimals: f.decimales }) }} />
  )
  const dia = (campo: 'statement_day' | 'payment_due_day', etiqueta: string) => (
    <Editable etiqueta={etiqueta} tipo="numero" valor={String(c[campo] ?? '')} mostrar={c[campo] ? String(c[campo]) : t.comun.sinValor} validar={vDia}
      vistaPrevia={(s) => { const nuevo = { ...c, [campo]: Number(s) }; const A = c.statement_day && c.payment_due_day ? corteYPago(hoy, c.statement_day, c.payment_due_day) : null; const B = nuevo.statement_day && nuevo.payment_due_day ? corteYPago(hoy, nuevo.statement_day, nuevo.payment_due_day) : null
        return { cambios: [{ etiqueta: k.proximoCorte, antes: A ? fecha(A.corte) : '—', despues: B ? fecha(B.corte) : '—' }, { etiqueta: k.pagoCorte, antes: A ? fecha(A.pago) : '—', despues: B ? fecha(B.pago) : '—' }], notas: [k.recalcula] } }}
      guardar={(s) => guardar({ [campo]: Number(s) })} />
  )
  const dineroCampo = (campo: 'annual_fee_cents' | 'late_fee_cents' | 'opening_fee_cents' | 'min_payment_fixed_cents', etiqueta: string, previa: (s: string) => Previa) => (
    <Editable etiqueta={etiqueta} tipo="dinero" valor={textoDinero(c[campo])} mostrar={<span className="tabular-nums">{c[campo] != null ? dinero(c[campo]!) : t.comun.sinValor}</span>} validar={vDinero()} vistaPrevia={previa} guardar={(s) => guardar({ [campo]: aCentavos(s) ?? 0 })} />
  )
  return (
    <Hoja abierta alCambiar={(a) => !a && alCerrar()} titulo={c.name}>
      <div className="flex flex-col gap-2" data-prueba="linea-credito">
        {faltantes(c).length > 0 && <Aviso rol="advertencia">{k.faltaCapturar}: {textoFaltan(faltantes(c))}</Aviso>}
        <Fila etiqueta={k.nombre}><Editable etiqueta={k.nombre} valor={c.name} mostrar={c.name} validar={vNombre} guardar={(s) => guardar({ name: s.trim() })} /></Fila>
        <Fila etiqueta={k.tipo}><span className="px-2">{k.tipos[c.kind]}</span></Fila>
        <Fila etiqueta={k.institucion}><Editable etiqueta={k.institucion} valor={c.lender ?? ''} mostrar={c.lender ?? t.comun.sinValor} validar={(s) => (s.length > 120 ? t.v.v09 : null)} guardar={(s) => guardar({ lender: s.trim() || null })} /></Fila>
        <Fila etiqueta={k.titular}><span className="px-2">Axel</span></Fila>
        <Fila etiqueta={k.limite}><Editable etiqueta={k.limite} tipo="dinero" valor={textoDinero(c.credit_limit_cents)} mostrar={<span className="tabular-nums">{c.credit_limit_cents != null ? dinero(c.credit_limit_cents) : t.comun.sinValor}</span>} validar={vDinero()} vistaPrevia={previaLimite} guardar={(s) => guardar({ credit_limit_cents: aCentavos(s) ?? 0 })} /></Fila>
        <Fila etiqueta={k.frecuencia}><Editable etiqueta={k.frecuencia} tipo="opcion" valor={c.payment_frequency} mostrar={k.frecuencias[c.payment_frequency]} opciones={opciones(k.frecuencias)} vistaPrevia={previaFrecuencia} extra={editorEjemplo} dependencias={depEjemplo} guardar={(s) => guardar({ payment_frequency: s as Cr['payment_frequency'] })} /></Fila>
        <Fila etiqueta={k.diasMin}><Editable etiqueta={k.diasMin} tipo="numero" valor={String(c.first_payment_min_days ?? '')} mostrar={c.first_payment_min_days != null ? String(c.first_payment_min_days) : t.comun.sinValor} validar={vEntero(0)} vistaPrevia={() => previaDias()} extra={editorEjemplo} dependencias={depEjemplo} guardar={(s) => guardar({ first_payment_min_days: Number(s) })} /></Fila>
        <Fila etiqueta={k.comAnual}>{dineroCampo('annual_fee_cents', k.comAnual, () => notaComision)}</Fila>
        <Fila etiqueta={k.comTardio}>{dineroCampo('late_fee_cents', k.comTardio, () => notaComision)}</Fila>
        <Fila etiqueta={k.comApertura}>{dineroCampo('opening_fee_cents', k.comApertura, () => notaComision)}</Fila>
        {revolvente && <LecturaAB c={c} guardar={guardar} />}
        {revolvente && <Fila etiqueta={k.anclaA}><Editable etiqueta={k.anclaA} tipo="fecha" valor={c.reading_a_anchor_date ?? ''} mostrar={c.reading_a_anchor_date ? fecha(c.reading_a_anchor_date) : t.comun.sinValor} validar={vFecha()} guardar={(s) => guardar({ reading_a_anchor_date: s })} /></Fila>}
        {tarjeta && <>
          <Fila etiqueta={k.tasaAnual}>{tasaCampo(k.tasaAnual, (s) => { const f = aFraccion(s); return previaCondicion(f ? { annual_rate_num: f.num, annual_rate_den: f.den, annual_rate_decimals: f.decimales } : {}) })}</Fila>
          <Fila etiqueta={k.base}><Editable etiqueta={k.base} tipo="opcion" valor={c.rate_base ?? ''} mostrar={c.rate_base ? (c.rate_base === 'daily' ? k.diaria : k.mensual) : t.comun.sinValor} opciones={[{ valor: '', texto: t.v.v22 }, { valor: 'daily', texto: k.diaria }, { valor: 'monthly', texto: k.mensual }]} validar={vOpcion()} vistaPrevia={(s) => previaCondicion({ rate_base: s as Cr['rate_base'] })} guardar={(s) => guardar({ rate_base: s as Cr['rate_base'] })} /></Fila>
          <Fila etiqueta={k.iva}><Editable etiqueta={k.iva} tipo="opcion" valor={c.vat_included == null ? '' : c.vat_included ? 'si' : 'no'} mostrar={c.vat_included == null ? t.comun.sinValor : c.vat_included ? k.ivaSi : k.ivaNo} opciones={[{ valor: '', texto: t.v.v22 }, { valor: 'si', texto: k.ivaSi }, { valor: 'no', texto: k.ivaNo }]} validar={vOpcion()}
            vistaPrevia={(s) => previaCondicion({ vat_included: s === 'si' })} guardar={(s) => guardar(s === 'si' ? { vat_included: true, vat_on_interest_bp: null } : { vat_included: false })} /></Fila>
          {c.vat_included === false && <Fila etiqueta={k.ivaPct}><Editable etiqueta={k.ivaPct} tipo="numero" valor={c.vat_on_interest_bp != null ? String(c.vat_on_interest_bp / 100) : ''} mostrar={c.vat_on_interest_bp != null ? pctDePuntosBase(c.vat_on_interest_bp) : `${t.comun.sinValor} · ${t.comun.provisional}`} validar={vPct(2)} guardar={(s) => guardar({ vat_on_interest_bp: aPuntosBase(s) })} /></Fila>}
          <Fila etiqueta={k.pagoMinPct}><Editable etiqueta={k.pagoMinPct} tipo="numero" valor={c.min_payment_pct_bp != null ? String(c.min_payment_pct_bp / 100) : ''} mostrar={c.min_payment_pct_bp != null ? pctDePuntosBase(c.min_payment_pct_bp) : t.comun.sinValor} validar={vPct(2)} vistaPrevia={(s) => previaCondicion({ min_payment_pct_bp: aPuntosBase(s) ?? undefined })} guardar={(s) => guardar({ min_payment_pct_bp: aPuntosBase(s) })} /></Fila>
          <Fila etiqueta={k.pagoMinFijo}>{dineroCampo('min_payment_fixed_cents', k.pagoMinFijo, (s) => previaCondicion({ min_payment_fixed_cents: aCentavos(s) ?? undefined }))}</Fila>
        </>}
        {c.kind === 'bank_card' && <><Fila etiqueta={k.corte}>{dia('statement_day', k.corte)}</Fila><Fila etiqueta={k.pago}>{dia('payment_due_day', k.pago)}</Fila></>}
        {c.kind === 'term_loan' && <Fila etiqueta={k.tasaAnual}>{tasaCampo(k.tasaAnual, previaTasaPrestamo)}</Fila>}
        <Fila etiqueta={t.comun.activa}><Editable etiqueta={t.comun.activa} tipo="interruptor" valor={c.is_active ? 'si' : 'no'} mostrar={null} vistaPrevia={() => ({})} guardar={(s) => guardar({ is_active: s === 'si' })} /></Fila>
        <Fila etiqueta={k.deuda}><span className="px-2 tabular-nums">{dinero(c.live_debt_cents ?? 0)} · {k.comprasVivas}: {c.live_purchase_count ?? 0}</span></Fila>
        {revolvente && <TasasPlazo credito={c} />}
      </div>
    </Hoja>
  )
}

function LecturaAB({ c, guardar }: { c: Cr; guardar: (x: Esquemas['CreditUpdate']) => Promise<void> }) {
  const conexion = useConexion()
  const [abierta, setAbierta] = useState(false)
  const [elegida, setElegida] = useState<'A' | 'B'>(c.availability_reading)
  const previa = usarConsulta('get', '/api/v1/credits/{credit_id}/availability-preview', { ruta: { credit_id: c.id } }, abierta)
  const serie = usarConsulta('get', '/api/v1/credits/{credit_id}/availability-series', { ruta: { credit_id: c.id } }, abierta)
  // N-284 (U1 §3.29; C3-matriz punto 2): el valor del interruptor A/B y el aviso «pendiente confirmar» salen de getCredit.
  const detalle = usarConsulta('get', '/api/v1/credits/{credit_id}', { ruta: { credit_id: c.id } })
  const linea = detalle.datos ?? c
  const otra = linea.availability_reading === 'A' ? 'B' : 'A'
  const guardarYLeer = async (x: Esquemas['CreditUpdate']) => { await guardar(x); detalle.recargar() }
  return (
    <>
      <Fila etiqueta={k.lectura}>
        <div className="flex flex-col gap-2">
          <span className="px-2" data-prueba="lectura-activa">{linea.availability_reading === 'A' ? k.lecturaA : k.lecturaB}</span>
          <Boton className="w-fit" deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => { setElegida(otra); setAbierta(true) }}>{k.usar(otra)}</Boton>
          {!linea.availability_confirmed && <Aviso rol="advertencia">{k.pendienteConfirmar}</Aviso>}
          <label className="flex min-h-11 items-center gap-2 px-2"><input type="checkbox" className="size-5" checked={linea.availability_confirmed} disabled={!conexion} onChange={(e) => void guardarYLeer({ availability_confirmed: e.target.checked })} />{k.yaConfirme}</label>
        </div>
      </Fila>
      <Hoja abierta={abierta} alCambiar={setAbierta} titulo={k.lectura} pie={<div className="flex justify-end gap-2"><Boton onClick={() => setAbierta(false)}>{textos.app.cancelar}</Boton><Boton variante="primario" deshabilitadoPor={!previa.datos ? t.comun.noCalculo : undefined} onClick={async () => { await guardarYLeer({ availability_reading: elegida }); setAbierta(false) }}>{k.usar(elegida)}</Boton></div>}>
        <ConDatos consulta={previa}>{(p) => (
          <div className="flex flex-col gap-3 desktop:grid desktop:grid-cols-2" data-prueba="lecturas-ab">
            {(['A', 'B'] as const).map((l) => { const r = l === 'A' ? p.reading_a : p.reading_b; return (
              <div key={l} className="flex flex-col gap-1 rounded-md border border-border-divider p-3">
                <h4 className="text-body font-medium">{l === 'A' ? k.lecturaA : k.lecturaB}</h4>
                <p className="tabular-nums">{k.hoy}: {dinero(r.available_today_cents)}</p><p className="tabular-nums">{k.minimo}: {dinero(r.horizon_min_cents)}</p>
                {serie.datos && <p className="text-small tabular-nums text-text-secondary">{k.serie}: {serie.datos.events.map((e) => dinero(l === 'A' ? e.available_a_cents : e.available_b_cents)).join(' · ')}</p>}
              </div>) })}
            {!linea.availability_confirmed && <Aviso rol="advertencia">{k.pendienteConfirmar}</Aviso>}
          </div>
        )}</ConDatos>
      </Hoja>
    </>
  )
}

function TasasPlazo({ credito }: { credito: Cr }) {
  const tasas = comoLista<Esquemas['CreditTermRate']>(usarConsulta('get', '/api/v1/credits/{credit_id}/term-rates', { ruta: { credit_id: credito.id } }))
  const [borrar, setBorrar] = useState<Tasa | null>(null)
  const [nueva, setNueva] = useState({ plazo: '', tasa: '' })
  const [err, setErr] = useState<string | null>(null)
  const [fallo, setFallo] = useState<string | null>(null)
  const conexion = useConexion()
  const guardar = (x: Tasa) => async (c: Esquemas['CreditTermRateUpdate']) => { await pedir('patch', '/api/v1/credit-term-rates/{credit_term_rate_id}', { ruta: { credit_term_rate_id: x.id }, cuerpo: c }); tasas.recargar() }
  const interes = (num: number, den: number) => dinero(multDiv(CAPITAL_EJEMPLO, num, den))
  const crear = async () => {
    const f = aFraccion(nueva.tasa); const e = vPlazo(nueva.plazo) ?? vPct(4)(nueva.tasa); setErr(e); setFallo(null)
    if (e || !f) return
    try { await pedir('post', '/api/v1/credits/{credit_id}/term-rates', { ruta: { credit_id: credito.id }, idempotencia: nuevaClave(), cuerpo: { term_count: Number(nueva.plazo), rate_num: f.num, rate_den: f.den, rate_decimals: f.decimales, source: 'manual' } } as never) }
    catch (x) { setFallo(textoFallo(x, { already_exists: t.v.v12 })); return }
    setNueva({ plazo: '', tasa: '' }); tasas.recargar()
  }
  return (
    <div className="flex flex-col gap-2 pt-2">
      <h4 className="text-body font-medium">{k.tasas}</h4>
      <ConDatos consulta={tasas}>{(d) => d.items.length === 0 ? <Vacio texto={k.sinTasas} /> : (
        <VistaBD def={{ clave: 'P-26.05-tasas', filas: [...d.items].sort((a, b) => a.term_count - b.term_count), idFila: (x) => x.id, vistas: ['tabla'], titulo: (x) => String(x.term_count),
          propiedades: [
            { id: 'term_count', nombre: k.plazo, alinear: 'derecha', valor: (x) => x.term_count, pintar: (x) => <Editable etiqueta={k.plazo} tipo="numero" valor={String(x.term_count)} mostrar={String(x.term_count)} validar={(s) => vPlazo(s) ?? (d.items.some((o) => o.id !== x.id && o.term_count === Number(s)) ? t.v.v12 : null)} guardar={(s) => guardar(x)({ term_count: Number(s) })} /> },
            { id: 'rate', nombre: k.tasaPlazo, alinear: 'derecha', valor: (x) => x.rate_num / x.rate_den, pintar: (x) => <Editable etiqueta={`${k.tasaPlazo} · ${x.term_count}`} tipo="numero" valor={x.rate_pct ?? textoTasa(x.rate_num, x.rate_den, x.rate_decimals)} mostrar={<span className="tabular-nums">{tasaDe(x.rate_num, x.rate_den, x.rate_decimals, x.rate_pct)}</span>} validar={vPct(4)}
              vistaPrevia={(s) => { const f = aFraccion(s)!; return { cambios: [{ etiqueta: k.intereses, antes: interes(x.rate_num, x.rate_den), despues: interes(f.num, f.den) }], notas: [k.interesesEjemplo(pctDeFraccion(f.num, f.den, f.decimales), x.term_count, interes(f.num, f.den))] } }}
              guardar={(s) => { const f = aFraccion(s)!; return guardar(x)({ rate_num: f.num, rate_den: f.den, rate_decimals: f.decimales }) }} /> },
            { id: 'source', nombre: k.origen, valor: (x) => x.source, pintar: (x) => <Editable etiqueta={k.origen} tipo="opcion" valor={x.source} mostrar={x.source === 'quote' ? k.cotizacion : k.tecleada} opciones={[{ valor: 'quote', texto: k.cotizacion }, { valor: 'manual', texto: k.tecleada }]} guardar={(s) => guardar(x)({ source: s as Tasa['source'] })} /> },
            { id: 'quote_date', nombre: k.fechaCot, valor: (x) => x.quote_date, pintar: (x) => <Editable etiqueta={k.fechaCot} tipo="fecha" valor={x.quote_date ?? ''} mostrar={x.quote_date ? fecha(x.quote_date) : '—'} validar={vFecha(true)} guardar={(s) => guardar(x)({ quote_date: s || null })} /> },
            { id: 'notes', nombre: k.nota, valor: (x) => x.notes ?? '', pintar: (x) => <Editable etiqueta={k.nota} valor={x.notes ?? ''} mostrar={x.notes || '—'} guardar={(s) => guardar(x)({ notes: s || null })} /> },
            { id: 'borrar', nombre: t.comun.eliminar, valor: () => '', pintar: (x) => <Boton variante="fantasma" deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => setBorrar(x)}>{t.comun.eliminar}</Boton> },
          ] }} />
      )}</ConDatos>
      <div className="flex flex-wrap items-end gap-2">
        <Campo etiqueta={k.plazo} inputMode="numeric" className="w-32" value={nueva.plazo} onChange={(e) => setNueva({ ...nueva, plazo: e.target.value })} />
        <Campo etiqueta={k.tasaPlazo} inputMode="decimal" className="w-40" value={nueva.tasa} onChange={(e) => setNueva({ ...nueva, tasa: e.target.value })} />
        <Boton deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => void crear()}><Plus aria-hidden className="size-5" />{k.nuevaTasa}</Boton>
      </div>
      {err && <p className="text-small text-role-error">{err}</p>}
      {fallo && <ErrorGuardado texto={fallo} alReintentar={() => void crear()} />}
      <DialogoCritico abierto={!!borrar} alCambiar={(a) => !a && setBorrar(null)} titulo={borrar ? k.eliminarTasa(borrar.term_count) : ''} verbo={borrar ? k.eliminarTasa(borrar.term_count) : ''}
        alConfirmar={async () => { if (borrar) await pedir('delete', '/api/v1/credit-term-rates/{credit_term_rate_id}', { ruta: { credit_term_rate_id: borrar.id } }); tasas.recargar() }}>{k.eliminarTasaTexto}</DialogoCritico>
    </div>
  )
}

function Promociones({ filas, creditos, alCambiar, alNueva }: { filas: Promo[]; creditos: Cr[]; alCambiar: () => void; alNueva: () => void }) {
  const [borrar, setBorrar] = useState<Promo | null>(null)
  const guardar = (x: Promo) => async (c: Esquemas['PromotionUpdate']) => { await pedir('patch', '/api/v1/promotions/{promotion_id}', { ruta: { promotion_id: x.id }, cuerpo: c }); alCambiar() }
  const nota = (): Previa => ({ notas: [k.notaPromo] })
  const tasa = (x: Promo) => (x.rate_num === 0 ? k.sinIntereses : tasaDe(x.rate_num, x.rate_den, x.rate_decimals, x.rate_pct))
  // Ctrl+/ (U4 P-26.05 «Atajos»: «editar a la vez varias promociones»). Campos PROVISIONALES (bloqueo WF3a-bis-ajustes-a §2):
  // activa, línea, plazo, compra mínima y fin de vigencia; cada uno con el aviso de R5-3 (punto 11) y, si cambia la vigencia, el conteo de vigentes hoy.
  const lineaDe = (id?: number | null) => creditos.find((c) => c.id === id)?.name ?? k.cualquiera
  const errorGrupo = (sel: Promo[], prop: string, v: string): string | null =>
    prop === 'term_count' ? vPlazo(v) : prop === 'min_purchase_cents' ? vDinero()(v) : prop === 'valid_to' ? (vFecha()(v) ?? (sel.some((x) => v < x.valid_from) ? t.v.v08 : null)) : null
  const conValor = (x: Promo, prop: string, v: string): Promo => ({ ...x, ...cuerpoGrupo(prop, v) } as Promo)
  const cuerpoGrupo = (prop: string, v: string): Esquemas['PromotionUpdate'] =>
    prop === 'is_active' ? { is_active: v === 'si' } : prop === 'credit_id' ? { credit_id: v === 'cualquiera' ? null : Number(v) } : prop === 'term_count' ? { term_count: Number(v) } : prop === 'min_purchase_cents' ? { min_purchase_cents: aCentavos(v) ?? 0 } : { valid_to: v }
  const mostrarGrupo = (x: Promo, prop: string): string =>
    prop === 'is_active' ? (x.is_active ? t.comun.si : t.comun.no) : prop === 'credit_id' ? lineaDe(x.credit_id) : prop === 'term_count' ? String(x.term_count) : prop === 'min_purchase_cents' ? dinero(x.min_purchase_cents) : fecha(x.valid_to)
  const nombreGrupo: Record<string, string> = { is_active: t.comun.activa, credit_id: k.linea, term_count: k.plazoPagos, min_purchase_cents: k.compraMin, valid_to: k.hasta }
  const previaGrupo = (sel: Promo[], prop: string, v: string): Cambio[] => {
    const e = errorGrupo(sel, prop, v)
    if (e) return [notaGrupo(e)]
    const nuevas = sel.map((x) => conValor(x, prop, v)), hoy = hoyISO()
    const cambios: Cambio[] = [{ etiqueta: nombreGrupo[prop], antes: resumen(sel.map((x) => mostrarGrupo(x, prop))), despues: mostrarGrupo(nuevas[0], prop) }]
    if (prop === 'is_active' || prop === 'valid_to') cambios.push({ etiqueta: t.grupo.vigentesHoy, antes: t.grupo.deN(sel.filter((x) => x.is_current).length, sel.length), despues: t.grupo.deN(nuevas.filter((x) => x.is_active && x.valid_from <= hoy && hoy <= x.valid_to).length, sel.length) })
    cambios.push(notaGrupo(k.notaPromo))
    return cambios
  }
  if (!filas.length) return <Vacio texto={k.sinPromos} accion={<Boton variante="primario" onClick={alNueva}>{k.nuevaPromo}</Boton>} />
  // ←/→: mes anterior o siguiente del calendario (U4 P-26.05 «Atajos»), fuera de los campos de texto.
  const flechas = (e: React.KeyboardEvent<HTMLDivElement>) => {
    if ((e.key !== 'ArrowLeft' && e.key !== 'ArrowRight') || e.altKey || e.ctrlKey || e.metaKey || (e.target as HTMLElement).closest('input, select, textarea, [contenteditable="true"]')) return
    const b = e.currentTarget.querySelector<HTMLButtonElement>(`button[aria-label="${e.key === 'ArrowLeft' ? textos.vistas.mesAnterior : textos.vistas.mesSiguiente}"]`)
    if (b) { e.preventDefault(); b.click() }
  }
  return (
    <div onKeyDown={flechas} data-prueba="promociones">
      <VistaBD def={{ clave: 'P-26.05-promociones', filas, idFila: (x) => x.id, vistas: ['calendario', 'tabla', 'tablero'], fecha: (x) => x.valid_from, titulo: (x) => `${x.merchant_name} · ${tasa(x)}`,
        columnasTablero: { propiedad: 'credit_id', columnas: [k.cualquiera, ...creditos.map((c) => c.name)] },
        editarVarias: {
          propiedades: [
            { id: 'is_active', nombre: t.comun.activa, tipo: 'opcion', opciones: siNo },
            { id: 'credit_id', nombre: k.linea, tipo: 'opcion', opciones: [{ valor: 'cualquiera', texto: k.cualquiera }, ...creditos.map((c) => ({ valor: String(c.id), texto: c.name }))] },
            { id: 'term_count', nombre: k.plazoPagos, tipo: 'numero' },
            { id: 'min_purchase_cents', nombre: k.compraMin, tipo: 'dinero' },
            { id: 'valid_to', nombre: k.hasta, tipo: 'fecha' },
          ],
          vistaPrevia: async (sel, prop, v) => previaGrupo(sel, prop, v),
          guardar: async (sel, prop, v) => {
            const e = errorGrupo(sel, prop, v)
            if (e) throw new ErrorApi(422, 'validation_error', e)
            await guardarVarias(sel, (x) => pedir('patch', '/api/v1/promotions/{promotion_id}', { ruta: { promotion_id: x.id }, cuerpo: cuerpoGrupo(prop, v) }), alCambiar)
          },
        },
        propiedades: [
          { id: 'merchant_name', nombre: k.comercio, valor: (x) => x.merchant_name, pintar: (x) => <Editable etiqueta={k.comercio} valor={x.merchant_name} mostrar={x.merchant_name} validar={vNombre} guardar={(s) => guardar(x)({ merchant_name: s.trim() })} /> },
          { id: 'credit_id', nombre: k.linea, valor: (x) => x.credit_id, grupo: (x) => creditos.find((c) => c.id === x.credit_id)?.name ?? k.cualquiera, pintar: (x) => creditos.find((c) => c.id === x.credit_id)?.name ?? k.cualquiera },
          { id: 'term_count', nombre: k.plazoPagos, alinear: 'derecha', valor: (x) => x.term_count, pintar: (x) => <Editable etiqueta={k.plazoPagos} tipo="numero" valor={String(x.term_count)} mostrar={String(x.term_count)} validar={vPlazo} vistaPrevia={nota} guardar={(s) => guardar(x)({ term_count: Number(s) })} /> },
          { id: 'rate', nombre: k.tasaPromo, alinear: 'derecha', valor: (x) => x.rate_num / x.rate_den, pintar: (x) => <Editable etiqueta={k.tasaPromo} tipo="numero" valor={textoTasa(x.rate_num, x.rate_den, x.rate_decimals)} mostrar={<span className="tabular-nums">{tasa(x)}</span>} validar={vPct(4)} vistaPrevia={nota} guardar={(s) => { const f = aFraccion(s)!; return guardar(x)({ rate_num: f.num, rate_den: f.den, rate_decimals: f.decimales }) }} /> },
          { id: 'min_purchase_cents', nombre: k.compraMin, alinear: 'derecha', valor: (x) => x.min_purchase_cents, pintar: (x) => <Editable etiqueta={k.compraMin} tipo="dinero" valor={textoDinero(x.min_purchase_cents)} mostrar={<span className="tabular-nums">{dinero(x.min_purchase_cents)}</span>} validar={vDinero()} vistaPrevia={nota} guardar={(s) => guardar(x)({ min_purchase_cents: aCentavos(s) ?? 0 })} /> },
          { id: 'valid_from', nombre: k.desde, valor: (x) => x.valid_from, pintar: (x) => <Editable etiqueta={k.desde} tipo="fecha" valor={x.valid_from} mostrar={fecha(x.valid_from)} validar={(s) => vFecha()(s) ?? (s > x.valid_to ? t.v.v08 : null)} vistaPrevia={nota} guardar={(s) => guardar(x)({ valid_from: s })} /> },
          { id: 'valid_to', nombre: k.hasta, valor: (x) => x.valid_to, pintar: (x) => <Editable etiqueta={k.hasta} tipo="fecha" valor={x.valid_to} mostrar={fecha(x.valid_to)} validar={(s) => vFecha()(s) ?? (s < x.valid_from ? t.v.v08 : null)} vistaPrevia={nota} guardar={(s) => guardar(x)({ valid_to: s })} /> },
          { id: 'is_current', nombre: k.vigente, valor: (x) => x.is_current, pintar: (x) => (x.is_current ? t.comun.si : t.comun.no) },
          { id: 'borrar', nombre: t.comun.eliminar, valor: () => '', pintar: (x) => <Boton variante="fantasma" onClick={() => setBorrar(x)}>{t.comun.eliminar}</Boton> },
        ] }} />
      <DialogoCritico abierto={!!borrar} alCambiar={(a) => !a && setBorrar(null)} titulo={borrar ? k.eliminarPromo(borrar.merchant_name) : ''} verbo={borrar ? k.eliminarPromo(borrar.merchant_name) : ''}
        alConfirmar={async () => { if (borrar) await pedir('delete', '/api/v1/promotions/{promotion_id}', { ruta: { promotion_id: borrar.id } }); alCambiar() }}
        porCodigo={{ in_use: { texto: t.v.v18, accion: { texto: t.comun.desactivar, alHacer: async () => { if (borrar) await pedir('patch', '/api/v1/promotions/{promotion_id}', { ruta: { promotion_id: borrar.id }, cuerpo: { is_active: false } }); alCambiar() } } } }}>{k.eliminarPromoTexto}</DialogoCritico>
    </div>
  )
}

function Comparador({ o, alCambiar }: { o: Esquemas['OptimizerSettings']; alCambiar: () => void }) {
  const patch = async (c: Esquemas['OptimizerSettingsUpdate']) => { await pedir('patch', '/api/v1/settings/optimizer', { cuerpo: c }); alCambiar() }
  const claves = ['lowest_interest', 'even_low_payment', 'max_available', 'most_value_per_peso'] as const
  const norm = normalizar(claves.map((x) => o.weights[x]))
  const [falloPesos, setFalloPesos] = useState<string | null>(null)
  const pesosIguales = async () => {
    setFalloPesos(null)
    try { await patch({ weights: { lowest_interest: 1, even_low_payment: 1, max_available: 1, most_value_per_peso: 1 } }) } catch (x) { setFalloPesos(textoFallo(x)) }
  }
  const previaPeso = (cl: (typeof claves)[number], n: number): Previa => { const nuevo = normalizar(claves.map((x) => (x === cl ? n : o.weights[x]))); return { cambios: claves.map((x, i) => ({ etiqueta: k.objetivos[x], antes: porcentaje(norm[i]), despues: porcentaje(nuevo[i]) })) } }
  type P = keyof Esquemas['OptimizerSettings']['parameters']
  const param = (cl: P, etiqueta: string, validar?: (s: string) => string | null) => <Parametro etiqueta={etiqueta} param={o.parameters[cl]} validar={validar} guardar={(n) => patch({ parameters: { [cl]: n } })} restablecer={() => patch({ reset_parameters: [cl] })} />
  return (
    <div className="flex flex-col gap-2">
      {o.missing_constraints.length > 0 && <Aviso rol="advertencia">{k.faltanRestr(o.missing_constraints.map((x) => k.restr[x]).join(', '))}</Aviso>}
      <Fila etiqueta={k.colchon}><Editable etiqueta={k.colchon} tipo="dinero" valor={textoDinero(o.min_available_buffer_cents)} mostrar={o.min_available_buffer_cents != null ? <span className="tabular-nums">{dinero(o.min_available_buffer_cents)}</span> : t.comun.sinValor} validar={vDinero()} guardar={(s) => patch({ min_available_buffer_cents: aCentavos(s) ?? 0 })} /></Fila>
      <Fila etiqueta={k.maxCompras}><Editable etiqueta={k.maxCompras} tipo="numero" valor={String(o.max_live_purchases ?? '')} mostrar={o.max_live_purchases != null ? String(o.max_live_purchases) : t.comun.sinValor} validar={vEntero(0)} guardar={(s) => patch({ max_live_purchases: Number(s) })} /></Fila>
      <Fila etiqueta={k.deudaCero}><Editable etiqueta={k.deudaCero} tipo="fecha" valor={o.zero_debt_by_date ?? ''} mostrar={o.zero_debt_by_date ? fecha(o.zero_debt_by_date) : t.comun.sinValor} validar={vFecha()} guardar={(s) => patch({ zero_debt_by_date: s })} /></Fila>
      {param('optimizer.date_window_days', k.ventana)}
      {param('comparator.timeout_seconds', k.tiempoMax, vEntero(1, 30))}
      <h4 className="pt-2 text-body font-medium">{k.pesos}</h4>
      {claves.map((cl, i) => (
        <Fila key={cl} etiqueta={k.objetivos[cl]} nota={k.pesoNorm(porcentaje(norm[i]))}>
          <Editable etiqueta={k.objetivos[cl]} tipo="numero" valor={String(o.weights[cl])} mostrar={<span className="tabular-nums">{o.weights[cl]}</span>}
            validar={(s) => vEntero(0)(s) ?? (claves.every((x) => (x === cl ? Number(s) : o.weights[x]) === 0) ? t.v.pesos : null)} vistaPrevia={(s) => previaPeso(cl, Number(s))} guardar={(s) => patch({ weights: { ...o.weights, [cl]: Number(s) } })} />
        </Fila>
      ))}
      <Boton className="w-fit" onClick={() => void pesosIguales()}>{k.pesosIguales}</Boton>
      {falloPesos && <ErrorGuardado texto={falloPesos} alReintentar={() => void pesosIguales()} />}
      <Fila etiqueta={k.tope}><Enlace a="/ajustes/12" className="flex min-h-11 items-center px-2 text-text-link">{k.topeEnlace}</Enlace></Fila>
      
    </div>
  )
}

function NuevoCredito({ abierta, alCambiar, alCrear }: { abierta: boolean; alCambiar: (a: boolean) => void; alCrear: () => void }) {
  const conexion = useConexion()
  const personas = comoLista<Esquemas['Person']>(usarConsulta('get', '/api/v1/persons', undefined, abierta))
  const [v, setV] = useState({ name: '', kind: '', frecuencia: '' })
  const [err, setErr] = useState<Record<string, string | null>>({})
  const [fallo, setFallo] = useState<string | null>(null)
  const [clave, setClave] = useState(nuevaClave)
  const guardar = async () => {
    const e = { name: vNombre(v.name), kind: vOpcion()(v.kind), frecuencia: vOpcion()(v.frecuencia) }; setErr(e); setFallo(null)
    const titular = personas.datos?.items.find((x) => x.participation === 'holder')
    if (Object.values(e).some(Boolean) || !titular) return
    try { await pedir('post', '/api/v1/credits', { idempotencia: clave, cuerpo: { name: v.name.trim(), kind: v.kind as Cr['kind'], holder_person_id: titular.id, payment_frequency: v.frecuencia as Cr['payment_frequency'], ...(v.kind === 'revolving_line' ? { availability_reading: 'B' as const } : {}) } }) }
    catch (x) { setFallo(textoFallo(x, { already_exists: t.v.v12 })); return }
    setV({ name: '', kind: '', frecuencia: '' }); setClave(nuevaClave()); alCambiar(false); alCrear()
  }
  const sel = (id: 'kind' | 'frecuencia', etiqueta: string, m: Record<string, string>) => (
    <label className="flex flex-col gap-1 text-small font-medium text-text-secondary">{etiqueta}
      <select className="h-11 rounded-md border border-border-control bg-surface-field px-3 text-body text-text-default focus-visible:outline-2 focus-visible:outline-border-focus" value={v[id]} onChange={(e) => setV({ ...v, [id]: e.target.value })}><option value="">{t.v.v22}</option>{opciones(m).map((o) => <option key={o.valor} value={o.valor}>{o.texto}</option>)}</select>
      {err[id] && <span className="text-role-error">{err[id]}</span>}
    </label>
  )
  return (
    <Hoja abierta={abierta} alCambiar={alCambiar} titulo={k.nuevo} pie={<div className="flex flex-col gap-2">{fallo && <ErrorGuardado texto={fallo} alReintentar={() => void guardar()} />}<div className="flex justify-end gap-2"><Boton onClick={() => alCambiar(false)}>{textos.app.cancelar}</Boton><Boton variante="primario" deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => void guardar()}>{textos.app.guardar}</Boton></div></div>}>
      <form className="flex flex-col gap-3" onSubmit={(e) => { e.preventDefault(); void guardar() }} onKeyDown={(e) => { if (e.key === 'Enter' && e.ctrlKey) { e.preventDefault(); void guardar() } }}>
        <Campo etiqueta={k.nombre} value={v.name} error={err.name ?? undefined} onChange={(e) => setV({ ...v, name: e.target.value })} />
        {sel('kind', k.tipo, k.tipos)}{sel('frecuencia', k.frecuencia, k.frecuencias)}
      </form>
    </Hoja>
  )
}

/** Alta de promoción (U4 P-26.05 bloque 3, createPromotion): se guarda con «Guardar» o Ctrl+Enter (§0.2, formularios de alta). */
function NuevaPromocion({ abierta, alCambiar, creditos, alCrear }: { abierta: boolean; alCambiar: (a: boolean) => void; creditos: Cr[]; alCrear: () => void }) {
  const conexion = useConexion()
  const tiendas = comoLista<Esquemas['Store']>(usarConsulta('get', '/api/v1/stores', { consulta: { limit: 200 } }, abierta))
  const vacio = { comercio: '', tienda: '', linea: '', plazo: '', tasa: '', minimo: '', desde: '', hasta: '', fuente: '' }
  const [v, setV] = useState(vacio)
  const [err, setErr] = useState<Partial<Record<keyof typeof vacio, string | null>>>({})
  const [fallo, setFallo] = useState<string | null>(null)
  const [clave, setClave] = useState(nuevaClave)
  const guardar = async () => {
    const e = { comercio: vNombre(v.comercio), plazo: vPlazo(v.plazo), tasa: vPct(4, true)(v.tasa), minimo: vDinero(false, true)(v.minimo), desde: vFecha()(v.desde), hasta: vFecha()(v.hasta) ?? (v.desde && v.hasta < v.desde ? t.v.v08 : null), fuente: v.fuente.length > 500 ? t.v.v09 : null }
    setErr(e); setFallo(null)
    if (Object.values(e).some(Boolean)) return
    const f = v.tasa.trim() ? aFraccion(v.tasa) : null
    try { await pedir('post', '/api/v1/promotions', { idempotencia: clave, cuerpo: { merchant_name: v.comercio.trim(), term_count: Number(v.plazo), valid_from: v.desde, valid_to: v.hasta, is_active: true,
      ...(v.linea ? { credit_id: Number(v.linea) } : {}), ...(v.tienda ? { store_id: Number(v.tienda) } : {}), ...(f ? { rate_num: f.num, rate_den: f.den, rate_decimals: f.decimales } : {}),
      ...(v.minimo.trim() ? { min_purchase_cents: aCentavos(v.minimo) ?? 0 } : {}), ...(v.fuente.trim() ? { source_note: v.fuente.trim() } : {}) } }) }
    catch (x) { setFallo(textoFallo(x, { already_exists: t.v.v12 })); return }
    setV(vacio); setClave(nuevaClave()); alCambiar(false); alCrear()
  }
  const lista = (id: 'linea' | 'tienda', etiqueta: string, primera: string, ops: { valor: string; texto: string }[]) => (
    <label className="flex flex-col gap-1 text-small font-medium text-text-secondary">{etiqueta}
      <select className="h-11 rounded-md border border-border-control bg-surface-field px-3 text-body text-text-default focus-visible:outline-2 focus-visible:outline-border-focus" value={v[id]} onChange={(e) => setV({ ...v, [id]: e.target.value })}><option value="">{primera}</option>{ops.map((o) => <option key={o.valor} value={o.valor}>{o.texto}</option>)}</select>
    </label>
  )
  const campo = (id: 'comercio' | 'plazo' | 'tasa' | 'minimo' | 'desde' | 'hasta' | 'fuente', etiqueta: string, extra: Partial<React.InputHTMLAttributes<HTMLInputElement>> = {}) => (
    <Campo etiqueta={etiqueta} value={v[id]} error={err[id] ?? undefined} onChange={(e) => setV({ ...v, [id]: e.target.value })} {...extra} />
  )
  return (
    <Hoja abierta={abierta} alCambiar={alCambiar} titulo={k.nuevaPromo} pie={<div className="flex flex-col gap-2">{fallo && <ErrorGuardado texto={fallo} alReintentar={() => void guardar()} />}<div className="flex justify-end gap-2"><Boton onClick={() => alCambiar(false)}>{textos.app.cancelar}</Boton><Boton variante="primario" deshabilitadoPor={!conexion ? textos.conexion.requiere : undefined} onClick={() => void guardar()}>{textos.app.guardar}</Boton></div></div>}>
      <form className="flex flex-col gap-3" data-prueba="nueva-promocion" onSubmit={(e) => { e.preventDefault(); void guardar() }} onKeyDown={(e) => { if (e.key === 'Enter' && e.ctrlKey) { e.preventDefault(); void guardar() } }}>
        {campo('comercio', k.comercio)}
        {lista('tienda', k.tienda, k.tiendaNinguna, (tiendas.datos?.items ?? []).map((x) => ({ valor: String(x.id), texto: x.name })))}
        {lista('linea', k.linea, k.cualquiera, creditos.map((x) => ({ valor: String(x.id), texto: x.name })))}
        {campo('plazo', k.plazoPagos, { inputMode: 'numeric' })}
        {campo('tasa', k.tasaPromo, { inputMode: 'decimal', ayuda: k.sinIntereses } as never)}
        {campo('minimo', k.compraMin, { inputMode: 'decimal' })}
        {campo('desde', k.desde, { type: 'date' })}
        {campo('hasta', k.hasta, { type: 'date' })}
        {campo('fuente', k.fuente)}
      </form>
    </Hoja>
  )
}
