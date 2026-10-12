// P-02 Todos los gastos (U1 §2.1, §2.2, §3.4: N-040 a N-047). Subsecciones Resumen · Detalle* · Tendencias · Calendario.
import { useState } from 'react'
import { Camera, Filter, Plus } from 'lucide-react'
import type { PropsPantalla } from '../../app/registro'
import { rutaDe } from '../../app/rutas'
import { conConsulta, navegar } from '../../app/enrutador'
import { pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { aCentavos, dinero, fecha, mes } from '../../api/formatos'
import type { Consulta, Esquemas } from '../../api/tipos'
import { Boton } from '../../componentes/Boton'
import { ConDatos, Vacio } from '../../componentes/Estados'
import { Grafica } from '../../componentes/Grafica'
import { TarjetaMetrica } from '../../componentes/TarjetaMetrica'
import { useAvisos } from '../../componentes/Avisos'
import { abrirRegistro } from '../../shell/estado'
import { textos } from '../../textos/es-MX'
import { VistaBD } from '../../vistas-bd/VistaBD'
import type { Propiedad } from '../../vistas-bd/tipos'
import { SelectorMes, Selector, centavosATexto, mensajeDe, periodoDe, rejillaTarjetas, useAlCambiarDatos, useCatalogos } from './comun'
import { t } from './textos'

type Gasto = Esquemas['Expense']
type Filtro = Consulta<'/api/v1/expenses', 'get'>
const FILTROS = ['familia', 'tipo', 'pago', 'ticket', 'categoria'] as const

export default function P02Gastos({ sub, consulta }: PropsPantalla) {
  const periodo = periodoDe(consulta)
  const filtro: Filtro = { period: periodo }
  const familia = consulta.get('familia'); if (familia === 'household' || familia === 'external' || familia === 'joint') filtro.family = familia
  const tipo = consulta.get('tipo'); if (tipo === 'fixed' || tipo === 'variable' || tipo === 'upcoming') filtro.kind = [tipo]
  const pago = consulta.get('pago'); if (pago === 'cash' || pago === 'debit' || pago === 'transfer' || pago === 'credit') filtro.payment_method = [pago]
  const ticket = consulta.get('ticket'); if (ticket === 'si' || ticket === 'no') filtro.has_attachment = ticket === 'si'
  const cat = Number(consulta.get('categoria')); if (cat > 0) filtro.category_id = [cat]
  const conFiltros = FILTROS.some((f) => consulta.get(f))
  const [verFiltros, setVerFiltros] = useState(conFiltros)
  const cats = useCatalogos()
  const quitar = () => navegar(conConsulta(Object.fromEntries(FILTROS.map((f) => [f, null]))), { reemplazar: true })
  const poner = (k: string) => (v: string) => navegar(conConsulta({ [k]: v || null }), { reemplazar: true })
  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap items-center gap-2">
        <SelectorMes valor={periodo} />
        <Boton variante="fantasma" aria-expanded={verFiltros} onClick={() => setVerFiltros(!verFiltros)}><Filter aria-hidden className="size-5" />{t.p02.filtros}{conFiltros && ` (${FILTROS.filter((f) => consulta.get(f)).length})`}</Boton>
        <div className="ml-auto flex gap-2">
          <Boton onClick={() => navegar(rutaDe('P-05'))}><Camera aria-hidden className="size-5" />{t.p02.fotoTicket}</Boton>
          <Boton variante="primario" onClick={() => abrirRegistro('gasto')}><Plus aria-hidden className="size-5" />{t.p02.registrar}</Boton>
        </div>
      </div>
      {verFiltros && (
        <div className="grid grid-cols-2 gap-3 rounded-lg border border-border-divider bg-surface-card p-3 desktop:grid-cols-6">
          <Selector etiqueta={t.p02.familia} valor={familia ?? ''} vacio={t.comun.todos} alCambiar={poner('familia')} opciones={[{ valor: 'household', texto: t.comun.hogar }, { valor: 'external', texto: t.comun.karen }, { valor: 'joint', texto: t.comun.conjunto }]} />
          <Selector etiqueta={t.comun.tipo} valor={tipo ?? ''} vacio={t.comun.todos} alCambiar={poner('tipo')} opciones={Object.entries(t.tipos).map(([valor, texto]) => ({ valor, texto }))} />
          <Selector etiqueta={t.comun.formaPago} valor={pago ?? ''} vacio={t.comun.todos} alCambiar={poner('pago')} opciones={(['cash', 'debit', 'transfer', 'credit'] as const).map((f) => ({ valor: f, texto: textos.registro.formas[f] }))} />
          <Selector etiqueta={t.p02.conTicket} valor={ticket ?? ''} vacio={t.comun.todos} alCambiar={poner('ticket')} opciones={[{ valor: 'si', texto: t.comun.si }, { valor: 'no', texto: t.comun.no }]} />
          <Selector etiqueta={t.comun.categoria} valor={cat > 0 ? String(cat) : ''} vacio={t.comun.todos} alCambiar={poner('categoria')} opciones={cats.categorias.map((c) => ({ valor: String(c.id), texto: c.name }))} />
          {conFiltros && <Boton variante="enlace" className="self-end" onClick={quitar}>{t.p02.quitarFiltros}</Boton>}
        </div>
      )}
      {sub === 'resumen' ? <Resumen filtro={filtro} /> : sub === 'tendencias' ? <Tendencias filtro={filtro} /> : <Lista key={JSON.stringify(filtro) + sub} filtro={filtro} soloCalendario={sub === 'calendario'} conFiltros={conFiltros} quitar={quitar} periodo={periodo} />}
    </div>
  )
}

function Resumen({ filtro }: { filtro: Filtro }) {
  const { period, family, category_id } = filtro
  const r = usarConsulta('get', '/api/v1/expenses/summary', { consulta: { period, family, category_id } })
  const tot = usarConsulta('get', '/api/v1/expenses/totals', { consulta: { period, family, category_id } })
  useAlCambiarDatos(r.recargar, tot.recargar)
  const cats = useCatalogos()
  const detalle = () => navegar(conConsulta({ sub: 'detalle' }))
  return (
    <ConDatos consulta={r}>{(s) => (
      <div className={rejillaTarjetas}>
        <TarjetaMetrica titulo={t.p02.gastoPeriodo} rol="expense" valor={dinero(s.household_share_cents)} alAbrir={detalle}
          nota={s.difference_cents === 0 ? t.p02.igualAnterior : t.p02.contraAnterior(dinero(Math.abs(s.difference_cents)), s.difference_cents > 0)} />
        <TarjetaMetrica titulo={t.p02.promDiario} valor={s.daily_average_cents == null ? '—' : dinero(s.daily_average_cents)} alAbrir={detalle} />
        <TarjetaMetrica titulo={t.p02.promSemanal} valor={s.weekly_average_cents == null ? '—' : dinero(s.weekly_average_cents)} alAbrir={detalle} />
        <TarjetaMetrica titulo={t.p02.promMensual} valor={s.monthly_average_cents == null ? '—' : dinero(s.monthly_average_cents)} alAbrir={detalle} />
        <TarjetaMetrica titulo={t.p02.mayorCategoria} valor={s.top_category_household_share_cents == null ? '—' : dinero(s.top_category_household_share_cents)} nota={cats.categoria(s.top_category_id)}
          alAbrir={s.top_category_id ? () => navegar(conConsulta({ sub: 'detalle', categoria: String(s.top_category_id) })) : detalle} />
        <TarjetaMetrica titulo={t.p02.parteKaren} valor={tot.datos ? dinero(tot.datos.external_share_cents) : '—'} nota={tot.datos && `${t.p02.numero}: ${tot.datos.expense_count}`} alAbrir={detalle} />
      </div>
    )}</ConDatos>
  )
}

function Tendencias({ filtro }: { filtro: Filtro }) {
  const r = usarConsulta('get', '/api/v1/expenses/monthly-trend', { consulta: { family: filtro.family, person_id: filtro.person_id } })
  const cats = useCatalogos()
  return (
    <ConDatos consulta={r}>{(d) => {
      if (!d.months.length) return <Vacio texto={textos.vacio.grafica} />
      const suma = new Map<number, number>()
      d.months.forEach((m) => m.categories.forEach((c) => suma.set(c.category_id ?? 0, (suma.get(c.category_id ?? 0) ?? 0) + c.household_share_cents)))
      const top = [...suma.entries()].sort((a, b) => b[1] - a[1]).slice(0, 6).map(([id]) => id)
      const series = top.map((id) => ({ nombre: cats.categoria(id), valores: d.months.map((m) => m.categories.find((c) => c.category_id === id)?.household_share_cents ?? 0) }))
      if (suma.size > 6) series.push({ nombre: t.p02.otras, valores: d.months.map((m) => m.categories.filter((c) => !top.includes(c.category_id ?? 0)).reduce((a, c) => a + c.household_share_cents, 0)) })
      return <Grafica titulo={t.p02.tendencia} resumen={t.p02.resumenTendencia} tipo="barras" categorias={d.months.map((m) => mes(m.month))} series={series} />
    }}</ConDatos>
  )
}

function Lista({ filtro, soloCalendario, conFiltros, quitar, periodo }: { filtro: Filtro; soloCalendario: boolean; conFiltros: boolean; quitar: () => void; periodo: string }) {
  const lista = usarConsulta('get', '/api/v1/expenses', { consulta: { ...filtro, limit: 50 } })
  const tot = usarConsulta('get', '/api/v1/expenses/totals', { consulta: { period: filtro.period, family: filtro.family, category_id: filtro.category_id } })
  const [mas, setMas] = useState<{ filas: Gasto[]; cursor: string | null | undefined }>({ filas: [], cursor: undefined })
  useAlCambiarDatos(lista.recargar, tot.recargar, () => setMas({ filas: [], cursor: undefined }))
  const cats = useCatalogos()
  const avisar = useAvisos()
  const parchar = async (g: Gasto, cuerpo: Esquemas['ExpenseUpdate']) => {
    // Si falla, CeldaEditable regresa al valor anterior y avisa (U1 §2.2).
    await pedir('patch', '/api/v1/expenses/{expense_id}', { ruta: { expense_id: g.id }, cuerpo }); lista.recargar(); tot.recargar(); setMas({ filas: [], cursor: undefined })
  }
  const opcionesCat = cats.categorias.map((c) => ({ valor: String(c.id), texto: c.name }))
  const propiedades: Propiedad<Gasto>[] = [
    { id: 'expense_date', nombre: t.comun.fecha, valor: (g) => g.expense_date, pintar: (g) => fecha(g.expense_date), editar: { tipo: 'fecha', textoInicial: (g) => g.expense_date, guardar: (g, v) => parchar(g, { expense_date: v }) } },
    { id: 'description', nombre: t.comun.descripcion, valor: (g) => g.description ?? '', editar: { tipo: 'texto', textoInicial: (g) => g.description ?? '', guardar: (g, v) => parchar(g, { description: v || null }) } },
    { id: 'category_id', nombre: t.comun.categoria, valor: (g) => cats.categoria(g.category_id), grupo: (g) => cats.categoria(g.category_id), editar: { tipo: 'opcion', opciones: opcionesCat, textoInicial: (g) => String(g.category_id), guardar: (g, v) => parchar(g, { category_id: Number(v) }) } },
    { id: 'amount_cents', nombre: t.comun.monto, valor: (g) => g.amount_cents, pintar: (g) => dinero(g.amount_cents), alinear: 'derecha',
      editar: { tipo: 'dinero', textoInicial: (g) => centavosATexto(g.amount_cents),
        vistaPrevia: async (g, v) => { const c = aCentavos(v); if (c == null || c <= 0) throw new Error(textos.registro.montoInvalido); return [{ etiqueta: t.comun.monto, antes: dinero(g.amount_cents), despues: dinero(c) }] },
        guardar: (g, v) => parchar(g, { amount_cents: aCentavos(v)! }) } },
    { id: 'person_id', nombre: t.comun.persona, valor: (g) => cats.persona(g.person_id), grupo: (g) => cats.persona(g.person_id) },
    { id: 'reparto', nombre: t.comun.reparto, valor: (g) => g.household_share_cents ?? 0, alinear: 'derecha', pintar: (g) => `${t.comun.hogar} ${dinero(g.household_share_cents ?? 0)} · ${t.comun.karen} ${dinero(g.external_share_cents ?? 0)}` },
    { id: 'kind', nombre: t.comun.tipo, valor: (g) => t.tipos[g.kind], grupo: (g) => t.tipos[g.kind], editar: { tipo: 'opcion', opciones: Object.entries(t.tipos).map(([valor, texto]) => ({ valor, texto })), textoInicial: (g) => g.kind, guardar: (g, v) => parchar(g, { kind: v as Gasto['kind'] }) } },
    { id: 'payment_method', nombre: t.comun.formaPago, valor: (g) => textos.registro.formas[g.payment_method], grupo: (g) => textos.registro.formas[g.payment_method] },
    { id: 'store_id', nombre: t.comun.tienda, valor: (g) => cats.tienda(g.store_id), grupo: (g) => cats.tienda(g.store_id) },
  ]
  const masPaginas = async (cursor: string) => {
    try { const p = await pedir('get', '/api/v1/expenses', { consulta: { ...filtro, limit: 50, cursor } }); setMas((m) => ({ filas: [...m.filas, ...p.items], cursor: p.next_cursor })) }
    catch (e) { avisar({ texto: mensajeDe(e), rol: 'error' }) }
  }
  return (
    <ConDatos consulta={lista}>{(d) => {
      const filas = [...d.items, ...mas.filas]
      const cursor = mas.cursor === undefined ? d.next_cursor : mas.cursor
      const vacio = conFiltros
        ? <Vacio texto={t.p02.sinCoincidencias} accion={<Boton onClick={quitar}>{t.p02.quitarFiltros}</Boton>} />
        : <Vacio texto={t.p02.vacio(mes(periodo))} accion={<Boton variante="primario" onClick={() => abrirRegistro('gasto')}>{t.p02.registrar}</Boton>} />
      return (
        <div className="flex flex-col gap-3">
          {tot.datos && filas.length > 0 && (
            <dl aria-label={t.p02.totales} className="grid grid-cols-2 gap-2 rounded-lg bg-surface-card p-3 text-small desktop:grid-cols-4">
              {[[t.p02.numero, String(tot.datos.expense_count)], [t.p02.sumaMontos, dinero(tot.datos.amount_cents)], [t.p02.parteHogar, dinero(tot.datos.household_share_cents)], [t.p02.parteKaren, dinero(tot.datos.external_share_cents)]].map(([k, v]) => (
                <div key={k} className="flex flex-col"><dt className="text-text-secondary">{k}</dt><dd className="font-semibold tabular-nums text-text-default">{v}</dd></div>))}
            </dl>
          )}
          <VistaBD def={{
            clave: soloCalendario ? 'P-02-calendario' : 'P-02-gastos', filas, idFila: (g) => g.id, propiedades,
            vistas: soloCalendario ? ['calendario'] : ['tabla', 'tablero', 'calendario', 'galeria'], fecha: (g) => g.expense_date,
            columnasTablero: { propiedad: 'category_id', columnas: cats.categorias.map((c) => c.name), moverA: (g, col) => { const c = cats.categorias.find((x) => x.name === col); return c ? parchar(g, { category_id: c.id }) : Promise.resolve() } },
            titulo: (g) => <span className="flex justify-between gap-2"><span>{g.description || cats.categoria(g.category_id)}{g.kind === 'upcoming' && <span className="text-caption text-role-warning"> · {t.p02.venidero}</span>}</span><span className="tabular-nums">{dinero(g.amount_cents, 'gasto')}</span></span>,
            alAbrir: (g) => navegar(rutaDe('P-04', { id: g.id })), vacio,
          }} />
          {cursor && <Boton className="self-center" onClick={() => masPaginas(cursor)}>{t.p02.mostrarMas}</Boton>}
        </div>
      )
    }}</ConDatos>
  )
}
