// Vista calendario (DESIGN §7.5.4): mes en cuadrícula, semana de lunes a domingo (K8).
import { useState } from 'react'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { ICONO_MD } from '../componentes/Icono'
import { textos } from '../textos/es-MX'
import { fechaConDia, hoyISO, mes } from '../api/formatos'
import { Boton } from '../componentes/Boton'
import type { DefVista } from './tipos'
const DIAS = ['Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb', 'Dom']
const mover = (aaaamm: string, d: number) => { const [a, m] = aaaamm.split('-').map(Number); const x = new Date(Date.UTC(a, m - 1 + d, 1)); return x.toISOString().slice(0, 7) }
export function VistaCalendario<T>({ def, filas }: { def: DefVista<T>; filas: T[] }) {
  const [actual, setActual] = useState(hoyISO().slice(0, 7))
  const [a, m] = actual.split('-').map(Number)
  const primero = (new Date(Date.UTC(a, m - 1, 1)).getUTCDay() + 6) % 7
  const dias = new Date(Date.UTC(a, m, 0)).getUTCDate()
  const celdas = [...Array(primero).fill(null), ...Array.from({ length: dias }, (_, i) => `${actual}-${String(i + 1).padStart(2, '0')}`)]
  const porDia = Object.groupBy(filas, (f) => def.fecha?.(f)?.slice(0, 10) ?? '')
  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-2">
        <Boton tam="icono" variante="fantasma" aria-label={textos.vistas.mesAnterior} onClick={() => setActual(mover(actual, -1))}><ChevronLeft aria-hidden className="size-6" size={ICONO_MD} /></Boton>
        <h3 className="text-body font-medium">{mes(actual)}</h3>
        <Boton tam="icono" variante="fantasma" aria-label={textos.vistas.mesSiguiente} onClick={() => setActual(mover(actual, 1))}><ChevronRight aria-hidden className="size-6" size={ICONO_MD} /></Boton>
      </div>
      <div role="grid" className="grid grid-cols-7 gap-px rounded-md border border-border-divider bg-border-divider text-small">
        <div role="row" className="contents">{DIAS.map((d) => <div role="columnheader" key={d} className="bg-surface-card p-1 text-center text-caption text-text-secondary">{d}</div>)}</div>
        {celdas.map((c, i) => (
          <div role="gridcell" key={c ?? `v${i}`} aria-label={c ? fechaConDia(c) : undefined} className="min-h-20 bg-surface-card p-1">
            {c && <span className="text-caption tabular-nums text-text-secondary">{Number(c.slice(8))}</span>}
            {c && porDia[c]?.map((f) => <button type="button" key={def.idFila(f)} onClick={() => def.alAbrir?.(f)} className="mt-1 block w-full truncate rounded-sm bg-surface-hover px-1 text-left text-caption tabular-nums coarse:min-h-11 coarse:py-3">{def.titulo(f)}</button>)}
          </div>
        ))}
      </div>
    </div>
  )
}
