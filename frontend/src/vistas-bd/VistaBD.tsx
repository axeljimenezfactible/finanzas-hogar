// Vista de base de datos completa: barra de la vista (cambiar vista, buscar, agrupar) + la vista elegida.
// La vista elegida se recuerda por dispositivo (Local), como pide U1 §5.2. Con `def.editarVarias`, la tabla trae
// selección de filas y `Ctrl+/` (o «Editar seleccionadas») abre la edición en grupo con vista previa sumada (U4 §0.2).
import { useRef, useState, type KeyboardEvent } from 'react'
import { CalendarDays, Images, Kanban, Table2 } from 'lucide-react'
import { textos } from '../textos/es-MX'
import { Vacio } from '../componentes/Estados'
import { Boton } from '../componentes/Boton'
import { useAvisos } from '../componentes/Avisos'
import { EdicionVarias } from './EdicionVarias'
import { cn } from '../componentes/utils'
import { estilosControl } from '../componentes/Campo'
import type { DefVista, TipoVista } from './tipos'
import { VistaTabla } from './VistaTabla'
import { VistaTablero } from './VistaTablero'
import { VistaCalendario } from './VistaCalendario'
import { VistaGaleria } from './VistaGaleria'
const ICONO = { tabla: Table2, tablero: Kanban, calendario: CalendarDays, galeria: Images }
const leer = (k: string) => { try { return localStorage.getItem(k) } catch { return null } }
export function VistaBD<T>({ def }: { def: DefVista<T> }) {
  const guardada = leer(`fh-vista:${def.clave}`) as TipoVista | null
  const [vista, setVista] = useState<TipoVista>(guardada && def.vistas.includes(guardada) ? guardada : def.vistas[0])
  const [texto, setTexto] = useState('')
  const [grupo, setGrupo] = useState('')
  const [sel, setSel] = useState<Set<string>>(() => new Set())
  const [enGrupo, setEnGrupo] = useState(false)
  const avisar = useAvisos()
  const raiz = useRef<HTMLDivElement>(null)
  const regreso = useRef<HTMLElement | null>(null)
  const ev = def.editarVarias
  const seleccionadas = ev ? def.filas.filter((f) => sel.has(String(def.idFila(f)))) : []
  const abrirGrupo = () => {
    if (!seleccionadas.length) return avisar({ texto: textos.vistas.eligeFilas, rol: 'info' })
    regreso.current = document.activeElement as HTMLElement | null
    setEnGrupo(true)
  }
  const cerrarGrupo = (aplicado: boolean) => {
    setEnGrupo(false)
    if (aplicado) setSel(new Set())
    const r = regreso.current
    requestAnimationFrame(() => (r?.isConnected ? r : raiz.current?.querySelector<HTMLElement>('input[type=checkbox]'))?.focus())
  }
  const teclas = (e: KeyboardEvent) => {
    if (ev && vista === 'tabla' && (e.ctrlKey || e.metaKey) && e.key === '/') { e.preventDefault(); abrirGrupo() }
  }
  const elegir = (v: TipoVista) => { setVista(v); try { localStorage.setItem(`fh-vista:${def.clave}`, v) } catch { /* sin almacenamiento */ } }
  const filtradas = texto ? def.filas.filter((f) => def.propiedades.some((p) => String(p.valor(f) ?? '').toLocaleLowerCase('es-MX').includes(texto.toLocaleLowerCase('es-MX')))) : def.filas
  return (
    <div ref={raiz} className="flex flex-col gap-3" onKeyDown={teclas}>
      <div className="flex flex-wrap items-center gap-2">
        <div role="tablist" aria-label="Vistas" className="flex gap-1">
          {def.vistas.map((v) => { const I = ICONO[v]; return (
            <button key={v} role="tab" type="button" aria-selected={vista === v} onClick={() => elegir(v)} className={cn('inline-flex h-9 items-center gap-1 rounded-md px-2 text-small coarse:h-11', vista === v ? 'bg-surface-selected font-medium' : 'hover:bg-surface-hover')}>
              <I aria-hidden className="size-5" />{textos.vistas[v]}
            </button>) })}
        </div>
        <input type="search" value={texto} onChange={(e) => setTexto(e.target.value)} aria-label={textos.vistas.filtrar} placeholder={textos.vistas.filtrar} className={cn(estilosControl, 'w-40 flex-1')} />
        {vista === 'tabla' && (
          <select aria-label={textos.vistas.agrupar} value={grupo} onChange={(e) => setGrupo(e.target.value)} className={cn(estilosControl, 'w-auto')}>
            <option value="">{textos.vistas.sinGrupo}</option>
            {def.propiedades.filter((p) => p.grupo).map((p) => <option key={p.id} value={p.id}>{textos.vistas.agrupar}: {p.nombre}</option>)}
          </select>
        )}
      </div>
      {ev && vista === 'tabla' && seleccionadas.length > 0 && (
        <div className="flex flex-wrap items-center gap-2">
          <span aria-live="polite" className="text-small text-text-secondary">{textos.vistas.seleccionadas(seleccionadas.length)}</span>
          <Boton variante="secundario" aria-keyshortcuts="Control+/" onClick={abrirGrupo}>{textos.vistas.editarSeleccionadas}<kbd className="hidden rounded-sm border border-border-divider px-1 text-caption text-text-secondary desktop:inline">{textos.vistas.atajoEditar}</kbd></Boton>
          <Boton variante="fantasma" onClick={() => { setSel(new Set()); setEnGrupo(false) }}>{textos.vistas.quitarSeleccion}</Boton>
        </div>
      )}
      {ev && enGrupo && seleccionadas.length > 0 && <EdicionVarias key={[...sel].join(',')} ev={ev} filas={seleccionadas} alTerminar={cerrarGrupo} />}
      {def.filas.length === 0 ? (def.vacio ?? <Vacio texto={textos.vacio.generico} />)
        : vista === 'tabla' ? <VistaTabla def={def} texto={texto} agruparPor={grupo || undefined} seleccion={ev ? { ids: sel, cambiar: setSel } : undefined} />
        : vista === 'tablero' ? <VistaTablero def={def} filas={filtradas} />
        : vista === 'calendario' ? <VistaCalendario def={def} filas={filtradas} />
        : <VistaGaleria def={def} filas={filtradas} />}
    </div>
  )
}
