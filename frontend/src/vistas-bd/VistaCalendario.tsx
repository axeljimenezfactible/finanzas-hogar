// Vista calendario (DESIGN §7.5.4): semana de lunes a domingo (K8). Desde 768 px, mes en cuadrícula; por debajo de 768,
// agenda por día (lista agrupada por fecha con encabezado fijo «Vie 30 oct 2026»), con renglones de al menos 44 px
// táctiles (U5-06, §4.6): en 360 y 412 una cuadrícula de 7 columnas no deja 44 px de ancho por renglón.
import { useState, useSyncExternalStore } from 'react'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { ICONO_MD } from '../componentes/Icono'
import { textos } from '../textos/es-MX'
import { fechaConDia, hoyISO, mes } from '../api/formatos'
import { Boton } from '../componentes/Boton'
import type { DefVista } from './tipos'
const DIAS = ['Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb', 'Dom']
const mover = (aaaamm: string, d: number) => { const [a, m] = aaaamm.split('-').map(Number); const x = new Date(Date.UTC(a, m - 1 + d, 1)); return x.toISOString().slice(0, 7) }
/** ¿La pantalla mide al menos `breakpoint.tablet` (768 px)? Se lee en vivo (girar el teléfono o cambiar la ventana). */
const CONSULTA_MES = '(min-width: 48rem)'
const suscribirAncho = (f: () => void) => { const mq = matchMedia(CONSULTA_MES); mq.addEventListener('change', f); return () => mq.removeEventListener('change', f) }
export function usarVistaMes() { return useSyncExternalStore(suscribirAncho, () => matchMedia(CONSULTA_MES).matches, () => true) }
export function VistaCalendario<T>({ def, filas }: { def: DefVista<T>; filas: T[] }) {
  const [actual, setActual] = useState(hoyISO().slice(0, 7))
  const enMes = usarVistaMes()
  const [a, m] = actual.split('-').map(Number)
  const primero = (new Date(Date.UTC(a, m - 1, 1)).getUTCDay() + 6) % 7
  const dias = new Date(Date.UTC(a, m, 0)).getUTCDate()
  const fechas = Array.from({ length: dias }, (_, i) => `${actual}-${String(i + 1).padStart(2, '0')}`)
  const celdas = [...Array(primero).fill(null), ...fechas]
  const porDia = Object.groupBy(filas, (f) => def.fecha?.(f)?.slice(0, 10) ?? '')
  const conRegistros = fechas.filter((c) => porDia[c]?.length)
  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-2">
        <Boton tam="icono" variante="fantasma" aria-label={textos.vistas.mesAnterior} onClick={() => setActual(mover(actual, -1))}><ChevronLeft aria-hidden className="size-6" size={ICONO_MD} /></Boton>
        <h3 className="text-body font-medium">{mes(actual)}</h3>
        <Boton tam="icono" variante="fantasma" aria-label={textos.vistas.mesSiguiente} onClick={() => setActual(mover(actual, 1))}><ChevronRight aria-hidden className="size-6" size={ICONO_MD} /></Boton>
      </div>
      {enMes ? (
        <div role="grid" className="grid grid-cols-7 gap-px rounded-md border border-border-divider bg-border-divider text-small">
          <div role="row" className="contents">{DIAS.map((d) => <div role="columnheader" key={d} className="bg-surface-card p-1 text-center text-caption text-text-secondary">{d}</div>)}</div>
          {celdas.map((c, i) => (
            <div role="gridcell" key={c ?? `v${i}`} aria-label={c ? fechaConDia(c) : undefined} className="min-h-20 min-w-0 bg-surface-card p-1">
              {c && <span className="text-caption tabular-nums text-text-secondary">{Number(c.slice(8))}</span>}
              {c && porDia[c]?.map((f) => <button type="button" key={def.idFila(f)} onClick={() => def.alAbrir?.(f)} className="mt-1 block min-h-6 w-full truncate rounded-sm bg-surface-hover px-1 text-left text-caption tabular-nums coarse:min-h-11 coarse:py-3">{def.titulo(f)}</button>)}
            </div>
          ))}
        </div>
      ) : conRegistros.length === 0 ? (
        <p className="text-small text-text-secondary">{textos.vistas.sinRegistrosMes(mes(actual))}</p>
      ) : (
        <ol data-prueba="agenda" aria-label={mes(actual)} className="flex flex-col rounded-md border border-border-divider bg-surface-card text-small">
          {conRegistros.map((c) => (
            <li key={c} className="border-b border-border-divider last:border-b-0">
              <h4 className="sticky top-0 bg-surface-card px-2 pt-2 text-caption font-medium text-text-secondary">{fechaConDia(c)}</h4>
              <ul className="flex flex-col">
                {porDia[c]!.map((f) => (
                  <li key={def.idFila(f)}>
                    <button type="button" onClick={() => def.alAbrir?.(f)} className="flex min-h-11 w-full items-center rounded-sm px-2 text-left tabular-nums hover:bg-surface-hover focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-border-focus">{def.titulo(f)}</button>
                  </li>
                ))}
              </ul>
            </li>
          ))}
        </ol>
      )}
    </div>
  )
}
