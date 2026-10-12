// Piezas compartidas de Comidas y Compras (P-10 a P-16): semanas ISO (K8, K38), catálogos del contrato,
// cantidades en milésimas (Milli de C1), marca de tienda y el menú «⋯». Solo tokens de C2.
import { useEffect, useState, type ReactNode } from 'react'
import { DropdownMenu } from 'radix-ui'
import { ChevronLeft, ChevronRight, Ellipsis } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import { todas } from '../../api/cliente'
import type { Esquemas } from '../../api/tipos'
import { fecha, hoyISO } from '../../api/formatos'
import { ErrorApi } from '../../api/errores'
import { conConsulta, navegar } from '../../app/enrutador'
import { Boton } from '../../componentes/Boton'
import { estilosControl } from '../../componentes/Campo'
import { ChipTienda } from '../../componentes/Chip'
import { textos } from '../../textos/es-MX'
import { textosCC as t } from './textos'

export type Tienda = Esquemas['Store']
export type Producto = Esquemas['Product']

// ── Semanas: lunes a domingo (K8); los domingos la semana por defecto es la siguiente (K38, U1 §0.2) ──
const DIA = 86_400_000
const utc = (iso: string) => Date.UTC(Number(iso.slice(0, 4)), Number(iso.slice(5, 7)) - 1, Number(iso.slice(8, 10)))
const aISO = (ms: number) => new Date(ms).toISOString().slice(0, 10)
export const sumarDias = (iso: string, n: number) => aISO(utc(iso) + n * DIA)
export function semanaDe(iso: string): string {
  const d = new Date(utc(iso))
  d.setUTCDate(d.getUTCDate() + 4 - (d.getUTCDay() || 7))
  const anio = d.getUTCFullYear()
  const n = Math.ceil(((+d - Date.UTC(anio, 0, 1)) / DIA + 1) / 7)
  return `${anio}-W${String(n).padStart(2, '0')}`
}
export function lunesDe(semana: string): string {
  const anio = Number(semana.slice(0, 4))
  const ene4 = Date.UTC(anio, 0, 4)
  return aISO(ene4 - ((new Date(ene4).getUTCDay() || 7) - 1) * DIA + (Number(semana.slice(6)) - 1) * 7 * DIA)
}
export const diasDe = (semana: string) => Array.from({ length: 7 }, (_, i) => sumarDias(lunesDe(semana), i))
export const moverSemana = (semana: string, n: number) => semanaDe(sumarDias(lunesDe(semana), 7 * n))
export const semanaPorDefecto = (hoy = hoyISO()) => semanaDe(new Date(utc(hoy)).getUTCDay() === 0 ? sumarDias(hoy, 1) : hoy)
export const textoSemana = (s: string) => t.semana(Number(s.slice(6)), fecha(lunesDe(s)), fecha(sumarDias(lunesDe(s), 6)))
const DIAS_CORTOS = ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb']
export const diaCorto = (iso: string) => `${DIAS_CORTOS[new Date(utc(iso)).getUTCDay()]} ${Number(iso.slice(8, 10))}`

/** Semana que muestra la pantalla: la de ?semana= o la por defecto, fijada al abrir (U1 §0.2). */
export function useSemana(consulta: URLSearchParams) {
  const [porDefecto] = useState(() => semanaPorDefecto())
  const q = consulta.get('semana')
  const semana = q && /^\d{4}-W\d{2}$/.test(q) ? q : porDefecto
  const cambiar = (s: string) => navegar(conConsulta({ semana: s === porDefecto ? null : s }), { reemplazar: true })
  return { semana, porDefecto, cambiar }
}
export const conSemana = (ruta: string, semana: string, porDefecto: string) => (semana === porDefecto ? ruta : `${ruta}?semana=${semana}`)

export function SelectorSemana({ semana, porDefecto, cambiar }: ReturnType<typeof useSemana>) {
  return (
    <div role="group" aria-label={t.cambiarSemana} className="flex flex-wrap items-center gap-1">
      <Boton tam="icono" variante="fantasma" aria-label={t.semanaAnterior} onClick={() => cambiar(moverSemana(semana, -1))}><ChevronLeft aria-hidden className="size-6" size={ICONO_MD} /></Boton>
      <span className="text-body font-medium" data-prueba="semana">{textoSemana(semana)}</span>
      <Boton tam="icono" variante="fantasma" aria-label={t.semanaSiguiente} onClick={() => cambiar(moverSemana(semana, 1))}><ChevronRight aria-hidden className="size-6" size={ICONO_MD} /></Boton>
      {semana === porDefecto
        ? <span className="rounded-full bg-surface-selected px-2 text-caption text-text-secondary">{t.semanaVigente}</span>
        : <Boton variante="enlace" onClick={() => cambiar(porDefecto)}>{t.irSemanaVigente}</Boton>}
    </div>
  )
}

// ── Catálogos (una carga por sesión de la página; se invalidan al crear) ──
type Catalogos = { tiendas: Tienda; productos: Producto; personas: Esquemas['Person']; tiempos: Esquemas['MealSlot']; recetas: Esquemas['Recipe']; grupos: Esquemas['ProductGroup'] }
type Clave = keyof Catalogos
const cargadores: { [K in Clave]: () => Promise<unknown> } = {
  tiendas: () => todas('get', '/api/v1/stores', { consulta: { limit: 200 } }),
  productos: () => todas('get', '/api/v1/products', { consulta: { limit: 200 } }),
  personas: () => todas('get', '/api/v1/persons'),
  tiempos: () => todas('get', '/api/v1/meal-slots', { consulta: { limit: 200 } }),
  recetas: () => todas('get', '/api/v1/recipes', { consulta: { limit: 200 } }),
  grupos: () => todas('get', '/api/v1/product-groups', { consulta: { limit: 200 } }),
}
const cache: Partial<Record<Clave, Promise<unknown>>> = {}
export const invalidar = (k: Clave) => { delete cache[k] }
export function usarCatalogo<K extends Clave>(k: K): Catalogos[K][] | undefined {
  const [datos, setDatos] = useState<Catalogos[K][]>()
  useEffect(() => {
    let vivo = true
    const p = (cache[k] ??= cargadores[k]().catch((e) => { delete cache[k]; throw e })) as Promise<Catalogos[K][]>
    p.then((v) => vivo && setDatos(v), () => vivo && setDatos([]))
    return () => { vivo = false }
  }, [k])
  return datos
}
export const porId = <T extends { id: number }>(l?: T[]) => new Map((l ?? []).map((x) => [x.id, x]))
/** Personas con plan de comidas: titular y co-titular (Karen es externa, decisión A). */
export const personasPlan = (l?: Esquemas['Person'][]) => (l ?? []).filter((p) => p.participation !== 'external' && p.is_active)

// ── Tiendas: la marca de color sale del nombre de la cadena (C1 no trae código de tienda: PROVISIONAL) ──
export function codigoTienda(s: Tienda | undefined, mapa?: Map<number, Tienda>): string {
  if (!s) return 'market'
  if (s.kind === 'market') return 'market'
  if (s.kind === 'branch' && s.parent_id && mapa?.get(s.parent_id)) return codigoTienda(mapa.get(s.parent_id), mapa)
  return s.name.normalize('NFD').replace(/[̀-ͯ']/g, '').toLowerCase().trim().replace(/\s+/g, '-')
}
export function Tiendita({ id, mapa }: { id?: number | null; mapa: Map<number, Tienda> }) {
  const s = id ? mapa.get(id) : undefined
  return s ? <ChipTienda codigo={codigoTienda(s, mapa)} nombre={s.name} /> : <span className="text-text-secondary">{t.sinTienda}</span>
}

// ── Cantidades (Milli: milésimas de la unidad base) ──
const num = new Intl.NumberFormat('es-MX', { maximumFractionDigits: 2 })
export const numero = (n: number) => num.format(n)
export function cantidad(milli: number | null | undefined, unidad?: string | null): string {
  if (milli == null) return '—'
  const v = milli / 1000
  if (unidad === 'g') return v >= 1000 ? `${num.format(v / 1000)} kg` : `${num.format(v)} g`
  if (unidad === 'ml') return v >= 1000 ? `${num.format(v / 1000)} l` : `${num.format(v)} ml`
  return `${num.format(v)} ${v === 1 ? 'pza' : 'pzas'}`
}
export const kcal = (n: number | null | undefined) => (n == null ? '—' : `${num.format(n)} kcal`)
export const porciones = (milli: number) => num.format(milli / 1000)
/** «1.5» → 1500 (enteros, sin flotantes); null si no es mayor que cero. */
export function aMilli(texto: string): number | null {
  const s = texto.trim().replace(',', '.')
  if (!/^\d+(\.\d{1,3})?$/.test(s)) return null
  const [e, d = ''] = s.split('.')
  const v = Number(e) * 1000 + Number(d.padEnd(3, '0'))
  return v > 0 ? v : null
}
export const mensajeError = (e: unknown) => (e instanceof ErrorApi ? e.message : t.noGuardado)

// ── Controles pequeños ──
export function Selector({ etiqueta, valor, alCambiar, opciones, vacio, id }: { etiqueta: string; valor: string; alCambiar: (v: string) => void; opciones: { valor: string; texto: string }[]; vacio?: string; id: string }) {
  return (
    <label htmlFor={id} className="flex flex-col gap-1 text-small font-medium">
      {etiqueta}
      <select id={id} value={valor} onChange={(e) => alCambiar(e.target.value)} className={`${estilosControl} coarse:h-11`}>
        {vacio !== undefined && <option value="">{vacio}</option>}
        {opciones.map((o) => <option key={o.valor} value={o.valor}>{o.texto}</option>)}
      </select>
    </label>
  )
}
export interface AccionMas { texto: string; alHacer: () => void; deshabilitadoPor?: string }
export function MenuMas({ acciones }: { acciones: AccionMas[] }) {
  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger asChild>
        <Boton tam="icono" variante="secundario" aria-label={textos.app.masOpciones}><Ellipsis aria-hidden className="size-6" size={ICONO_MD} /></Boton>
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content align="end" sideOffset={4} className="z-50 min-w-56 rounded-md border border-border-divider bg-surface-raised p-1 text-text-default shadow-1">
          {acciones.map((a) => (
            <DropdownMenu.Item key={a.texto} disabled={!!a.deshabilitadoPor} onSelect={a.alHacer} className="flex h-9 cursor-pointer items-center rounded-sm px-3 text-small outline-none data-[highlighted]:bg-surface-hover data-[disabled]:cursor-not-allowed data-[disabled]:text-text-disabled coarse:h-11">
              {a.texto}{a.deshabilitadoPor && <span className="sr-only"> · {a.deshabilitadoPor}</span>}
            </DropdownMenu.Item>
          ))}
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  )
}
/** Barra de acciones de la pantalla: a la izquierda el contexto, a la derecha la primaria y «⋯». */
export const BarraAcciones = ({ children, acciones }: { children?: ReactNode; acciones: ReactNode }) => (
  <div className="flex flex-col gap-2 tablet:flex-row tablet:items-center tablet:justify-between">
    <div className="min-w-0">{children}</div>
    <div className="flex flex-wrap items-center gap-2">{acciones}</div>
  </div>
)
export const Seccion = ({ titulo, children, accion }: { titulo: string; children: ReactNode; accion?: ReactNode }) => (
  <section className="flex flex-col gap-2">
    <div className="flex items-center justify-between gap-2"><h2 className="text-title-3 font-semibold">{titulo}</h2>{accion}</div>
    {children}
  </section>
)
/** Filas de una página de C1 (Page.items llega como unknown[] en el tipo generado por la composición allOf). */
export const filasDe = <T,>(d: { items?: unknown[] }) => (d.items ?? []) as T[]
