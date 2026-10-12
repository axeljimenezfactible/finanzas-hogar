// P-14 Comparar canasta (U1 §2.1–§2.2, §3.16 N-160 a N-165; flujo F-3). Cifras del motor S5 tal como llegan de C1.
import { useState } from 'react'
import { Scale } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import type { PropsPantalla } from '../../app/registro'
import type { Esquemas } from '../../api/tipos'
import { nuevaClave, pedir } from '../../api/cliente'
import { usarConsulta } from '../../api/consulta'
import { useConexion } from '../../api/conexion'
import { dinero, fecha, porcentaje } from '../../api/formatos'
import { navegar } from '../../app/enrutador'
import { rutaDe } from '../../app/rutas'
import { textos } from '../../textos/es-MX'
import { Boton } from '../../componentes/Boton'
import { ConDatos, Franja, Vacio } from '../../componentes/Estados'
import { useAvisos } from '../../componentes/Avisos'
import { cn } from '../../componentes/utils'
import { VistaBD } from '../../vistas-bd/VistaBD'
import { BarraAcciones, Seccion, SelectorSemana, Selector, Tiendita, cantidad, conSemana, filasDe, mensajeError, porId, textoSemana, usarCatalogo, useSemana } from './comun'
import { textosCC as t } from './textos'

type Comparacion = Esquemas['BasketComparison']
const mas = (c: number) => (c > 0 ? dinero(c, 'ingreso') : dinero(c))

export default function P14Canasta({ sub, consulta }: PropsPantalla) {
  const sem = useSemana(consulta)
  const conexion = useConexion()
  const avisar = useAvisos()
  const [canasta, setCanasta] = useState('')
  const canastas = usarConsulta('get', '/api/v1/baskets', { consulta: { week: sem.semana, limit: 50 } })
  const gasto = canasta ? Number(canasta) : undefined
  const comp = usarConsulta('get', '/api/v1/basket-comparison', { consulta: gasto ? { basket: 'expense', expense_id: gasto, week: sem.semana } : { basket: 'shopping_list', week: sem.semana } })
  const tiendas = porId(usarCatalogo('tiendas'))
  const productos = porId(usarCatalogo('productos'))
  const nombreT = (id?: number) => (id ? tiendas.get(id)?.name ?? '—' : '—')
  const sinRed = conexion ? undefined : textos.conexion.requiere
  async function usar(storeId: number) {
    try {
      const r = await pedir('post', '/api/v1/shopping-list/assign-store', { consulta: { week: sem.semana }, cuerpo: { store_id: storeId }, idempotencia: nuevaClave() })
      avisar({ texto: t.tiendaAsignada(r.assigned_count, r.without_price_count), rol: 'exito' })
      navegar(conSemana(rutaDe('P-13'), sem.semana, sem.porDefecto))
    } catch (e) { avisar({ texto: mensajeError(e), rol: 'error' }) }
  }
  const opciones = filasDe<Esquemas['Basket']>(canastas.datos ?? {}).filter((b) => b.basket === 'expense' && b.expense_id)
    .map((b) => ({ valor: String(b.expense_id), texto: t.compraDel(b.expense_date ? fecha(b.expense_date) : '—', nombreT(b.store_id), b.total_cents != null ? dinero(b.total_cents) : '—') }))
  return (
    <div className="flex flex-col gap-4">
      <BarraAcciones acciones={null}><SelectorSemana {...sem} /></BarraAcciones>
      <Selector id="p14-canasta" etiqueta={t.canasta} valor={canasta} alCambiar={setCanasta} opciones={[{ valor: '', texto: t.listaSemana(textoSemana(sem.semana)) }, ...opciones]} />
      <ConDatos consulta={comp}>{(c) => {
        if (!c.items.length) return <Vacio icono={<Scale aria-hidden className="size-6" size={ICONO_MD} />} texto={t.canastaVacia} />
        const tiendasC = [...c.stores].sort((a, b) => a.total_cents - b.total_cents)
        const barata = tiendasC.find((s) => s.store_id === c.cheapest_store_id)
        if (!barata || tiendasC.every((s) => s.missing_count >= c.items.length)) return <Vacio icono={<Scale aria-hidden className="size-6" size={ICONO_MD} />} texto={t.sinPreciosComparar} accion={<Boton onClick={() => navegar(`${rutaDe('P-15')}?hoja=precio`)}>{t.capturarPrecios}</Boton>} />
        const acciones = (
          <div className="flex flex-col gap-2 tablet:flex-row tablet:justify-end">
            <Boton onClick={() => navegar(`${rutaDe('P-20')}?monto=${barata.total_cents}`)}>{t.simularCredito}</Boton>
            <Boton variante="primario" deshabilitadoPor={gasto ? t.soloLista : sinRed} onClick={() => usar(barata.store_id)}>{t.usarTienda(nombreT(barata.store_id))}</Boton>
          </div>
        )
        if (sub === 'detalle') return <><Detalle c={c} nombreT={nombreT} nombreP={(id) => productos.get(id)?.name ?? '—'} unidad={(id) => productos.get(id)?.base_unit} tiendas={tiendasC.map((s) => s.store_id)} />{acciones}</>
        return (
          <div className="flex flex-col gap-4">
            <p className="text-body tabular-nums" data-prueba="frase-barata">{t.masBarata(nombreT(barata.store_id), dinero(barata.total_cents))}</p>
            <ul className="grid gap-3 tablet:grid-cols-2">
              {tiendasC.map((s) => {
                const conPrecio = c.items.length - s.missing_count
                return (
                  <li key={s.store_id} data-prueba={`tienda-${s.store_id}`} className={cn('flex flex-col gap-1 rounded-md border bg-surface-card p-4', s.store_id === barata.store_id ? 'border-border-selected' : 'border-border-divider')}>
                    <div className="flex flex-wrap items-center gap-2"><Tiendita id={s.store_id} mapa={tiendas} />{s.assumes_membership && <span className="text-caption text-text-secondary">({t.membresiaSupuesta})</span>}{s.store_id === barata.store_id && <span className="rounded-full bg-surface-selected px-2 text-caption">{t.laMasBarata}</span>}</div>
                    <span className="text-title-2 font-semibold tabular-nums">{dinero(s.total_cents)}</span>
                    {s.store_id !== barata.store_id && <span className="text-small tabular-nums">{t.contraMasBarata(mas(s.difference_vs_cheapest_cents ?? s.total_cents - barata.total_cents))} · {t.enIdenticos(mas(s.identical_total_cents - barata.identical_total_cents))}</span>}
                    <span className="text-small text-text-secondary">{t.cobertura(conPrecio, c.items.length, porcentaje((conPrecio * 100) / c.items.length))}</span>
                    <span className="text-small text-text-secondary">{[t.identicos(s.identical_count), t.sustitutos(s.substitute_count), t.sinPrecioN(s.missing_count), t.viejos(s.stale_count)].join(' · ')}</span>
                  </li>
                )
              })}
            </ul>
            <Seccion titulo={t.ganadores}>
              <ul className="flex flex-col gap-1">
                {[...c.winners].sort((a, b) => b.wins - a.wins).map((w) => <li key={w.store_id} className="flex items-center justify-between gap-2 text-small"><Tiendita id={w.store_id} mapa={tiendas} /><span className="tabular-nums font-medium">{w.wins.toLocaleString('es-MX')}</span></li>)}
              </ul>
              <p className="text-small text-text-secondary">{t.empates(c.ties, c.items.length)}</p>
            </Seccion>
            {c.membership && <Membresia m={c.membership} nombre={nombreT(c.membership.store_id)} />}
            {acciones}
          </div>
        )
      }}</ConDatos>
    </div>
  )
}

function Membresia({ m, nombre }: { m: NonNullable<Comparacion['membership']>; nombre: string }) {
  const filas: [string, string][] = [[t.ahorroCanasta, dinero(m.basket_savings_cents)], [t.ahorroAnual, m.annual_savings_cents != null ? dinero(m.annual_savings_cents) : '—'], [t.costoAnual, m.annual_fee_cents != null ? dinero(m.annual_fee_cents) : '—'], [t.equilibrio, m.break_even_purchases != null ? t.compras(m.break_even_purchases) : '—']]
  return (
    <Seccion titulo={t.membresia(nombre)}>
      <Franja rol={m.verdict === 'worth_it' ? 'info' : 'advertencia'}>{t.veredictos[m.verdict]}</Franja>
      <dl className="grid grid-cols-2 gap-2 desktop:grid-cols-4">
        {filas.map(([k, v]) => <div key={k} className="flex flex-col rounded-md border border-border-divider bg-surface-card p-3"><dt className="text-small text-text-secondary">{k}</dt><dd className="font-semibold tabular-nums">{v}</dd></div>)}
      </dl>
    </Seccion>
  )
}

function Detalle({ c, nombreT, nombreP, unidad, tiendas }: { c: Comparacion; nombreT: (id?: number) => string; nombreP: (id: number) => string; unidad: (id: number) => string | undefined; tiendas: number[] }) {
  type Fila = Comparacion['items'][number]
  const filas = [...c.items].sort((a, b) => (b.difference_cents ?? 0) - (a.difference_cents ?? 0))
  return (
    <VistaBD def={{
      clave: 'P-14-detalle', filas, idFila: (i: Fila) => i.product_id, vistas: ['tabla'], titulo: (i: Fila) => nombreP(i.product_id),
      propiedades: [
        { id: 'producto', nombre: t.producto, valor: (i: Fila) => nombreP(i.product_id), pintar: (i: Fila) => <span className="flex flex-col"><span className="font-medium">{nombreP(i.product_id)}</span>{i.is_substitute && <span className="text-caption text-text-secondary">{t.sustituto}{i.substitute_note ? `: ${i.substitute_note}` : ''}</span>}</span> },
        { id: 'cantidad', nombre: t.cantidad, valor: (i: Fila) => i.quantity_milli, pintar: (i: Fila) => cantidad(i.quantity_milli, unidad(i.product_id)), alinear: 'derecha' },
        ...tiendas.map((s) => ({
          id: `tienda-${s}`, nombre: nombreT(s), alinear: 'derecha' as const, valor: (i: Fila) => i.prices.find((p) => p.store_id === s)?.price_cents ?? null,
          pintar: (i: Fila) => {
            const p = i.prices.find((x) => x.store_id === s)
            if (!p) return <span className="text-text-secondary">{t.sinPrecio}</span>
            return (
              <span className={cn('flex flex-col items-end', i.winner_store_id === s && 'font-semibold')}>
                <span className="tabular-nums">{dinero(p.price_cents)}{p.is_stale && <span className="ml-1 rounded-sm bg-role-warning-surface px-1 text-caption text-role-warning">{t.viejo}</span>}</span>
                <span className="text-caption font-normal text-text-secondary">{t.porUnidad(dinero(p.base_unit_price_cents), t.unidades[unidad(i.product_id) ?? 'piece'])} · {fecha(p.observed_date)} · {t.fuentes[p.source]}</span>
              </span>
            )
          },
        })),
        { id: 'ganadora', nombre: t.ganadora, valor: (i: Fila) => nombreT(i.winner_store_id) },
        { id: 'diferencia', nombre: t.diferencia, valor: (i: Fila) => i.difference_cents ?? null, pintar: (i: Fila) => (i.difference_cents != null ? mas(i.difference_cents) : '—'), alinear: 'derecha' },
      ],
    }} />
  )
}
