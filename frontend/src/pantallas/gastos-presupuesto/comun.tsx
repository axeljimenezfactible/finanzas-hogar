// Piezas compartidas del área Gastos y Presupuesto: catálogos, selector de mes, barra de avance, selector nativo,
// confirmación y utilidades de periodo. Solo tokens C2; textos en ./textos.
import { useCallback, useEffect, useState, type ReactNode } from 'react'
import { ChevronLeft, ChevronRight, TriangleAlert } from 'lucide-react'
import { ICONO_MD } from '../../componentes/Icono'
import { todas } from '../../api/cliente'
import { useConexion } from '../../api/conexion'
import { ErrorApi } from '../../api/errores'
import { aCentavos as aCentavosLocal, hoyISO, mes, porcentaje } from '../../api/formatos'
import type { Esquemas } from '../../api/tipos'
import { conConsulta, navegar } from '../../app/enrutador'
import { Boton } from '../../componentes/Boton'
import { Campo as CampoTexto, estilosControl } from '../../componentes/Campo'
import { Hoja } from '../../componentes/Hoja'
import { cn } from '../../componentes/utils'
import { textos } from '../../textos/es-MX'
import { t } from './textos'

// ── Periodo (filtro global `period` de C1: aquí solo el mes calendario AAAA-MM) ──
export const mesActual = () => hoyISO().slice(0, 7)
export function periodoDe(consulta: URLSearchParams): string {
  const p = consulta.get('periodo')
  return p && /^\d{4}-(0[1-9]|1[0-2])$/.test(p) ? p : mesActual()
}
export function sumarMes(aaaamm: string, n: number): string {
  const total = Number(aaaamm.slice(0, 4)) * 12 + Number(aaaamm.slice(5, 7)) - 1 + n
  return `${Math.floor(total / 12)}-${String((total % 12) + 1).padStart(2, '0')}`
}
export const primerDia = (aaaamm: string) => `${aaaamm}-01`
export function ultimoDia(aaaamm: string): string {
  const d = new Date(Date.UTC(Number(aaaamm.slice(0, 4)), Number(aaaamm.slice(5, 7)), 0))
  return `${aaaamm}-${String(d.getUTCDate()).padStart(2, '0')}`
}
export function SelectorMes({ valor }: { valor: string }) {
  const ir = (n: number) => navegar(conConsulta({ periodo: sumarMes(valor, n) }), { reemplazar: true })
  return (
    <div className="flex items-center gap-1" role="group" aria-label={t.comun.periodo}>
      <Boton tam="icono" variante="fantasma" aria-label={t.comun.mesAnterior} onClick={() => ir(-1)}><ChevronLeft aria-hidden className="size-6" size={ICONO_MD} /></Boton>
      <span className="min-w-24 text-center text-body font-medium" aria-live="polite">{mes(valor)}</span>
      <Boton tam="icono" variante="fantasma" aria-label={t.comun.mesSiguiente} onClick={() => ir(1)}><ChevronRight aria-hidden className="size-6" size={ICONO_MD} /></Boton>
    </div>
  )
}

// ── Catálogos (categorías, personas, tiendas, reglas de reparto): una carga por sesión de la página ──
type Catalogos = { categorias: Esquemas['Category'][]; personas: Esquemas['Person'][]; tiendas: Esquemas['Store'][]; reglas: Esquemas['CostShareRule'][] }
let cache: Promise<Catalogos> | null = null
const cargar = () => (cache ??= Promise.all([
  todas('get', '/api/v1/categories').catch(() => []),
  todas('get', '/api/v1/persons').catch(() => []),
  todas('get', '/api/v1/stores').catch(() => []),
  todas('get', '/api/v1/cost-share-rules').catch(() => []),
]).then(([categorias, personas, tiendas, reglas]) => {
  const c = { categorias, personas, tiendas, reglas } as Catalogos
  if (!categorias.length) cache = null // reintenta en la siguiente pantalla
  return c
}))
export function useCatalogos() {
  const [c, setC] = useState<Catalogos>({ categorias: [], personas: [], tiendas: [], reglas: [] })
  useEffect(() => { let vivo = true; cargar().then((x) => vivo && setC(x)); return () => { vivo = false } }, [])
  const nombre = <T extends { id: number }>(lista: T[], id: number | null | undefined, campo: (x: T) => string) => (id == null ? '—' : (lista.find((x) => x.id === id) ? campo(lista.find((x) => x.id === id)!) : '—'))
  return {
    ...c,
    categoria: (id?: number | null) => nombre(c.categorias, id, (x) => x.name),
    persona: (id?: number | null) => nombre(c.personas, id, (x) => x.display_name),
    tienda: (id?: number | null) => nombre(c.tiendas, id, (x) => x.name),
    regla: (id?: number | null) => nombre(c.reglas, id, (x) => x.name),
  }
}

// ── Recarga cuando otra hoja guardó un gasto (fh:datos-cambiaron) ──
export function useAlCambiarDatos(...recargas: (() => void)[]) {
  const f = useCallback(() => recargas.forEach((r) => r()), recargas) // eslint-disable-line react-hooks/exhaustive-deps
  useEffect(() => { window.addEventListener('fh:datos-cambiaron', f); return () => window.removeEventListener('fh:datos-cambiaron', f) }, [f])
}

/** Motivo para desactivar una escritura sin conexión (U1 §2.4: solo lectura). */
export const useSinRed = () => (useConexion() ? undefined : textos.conexion.requiere)
export const mensajeDe = (e: unknown) => (e instanceof ErrorApi ? e.message : textos.errores.generico)

// ── Barra de avance (DESIGN §7.6): pista 8 px, llenado neutro; aviso y rebasado según el estado que da el API ──
export function BarraAvance({ pct, estado, excedente }: { pct: number; estado?: 'normal' | 'warning' | 'exceeded' | null; excedente?: string }) {
  const alerta = estado === 'warning' || estado === 'exceeded'
  return (
    <div className="flex flex-col gap-1">
      <div role="meter" aria-valuemin={0} aria-valuemax={100} aria-valuenow={Math.min(pct, 100)} aria-valuetext={t.comun.delPlan(porcentaje(pct))} className="h-2 w-full overflow-hidden rounded-full bg-viz-track">
        <div className={cn('h-full rounded-full', alerta ? 'bg-role-warning' : 'bg-viz-meter')} style={{ width: `${Math.max(0, Math.min(pct, 100))}%` }} />
      </div>
      {alerta && (
        <span className="flex items-center gap-1 text-caption text-role-warning">
          <TriangleAlert aria-hidden className="size-5 shrink-0" />{t.comun.usado(porcentaje(pct))}{estado === 'exceeded' && <> · {t.comun.rebasado}{excedente && <span className="tabular-nums"> {excedente}</span>}</>}
        </span>
      )}
    </div>
  )
}

// ── Selector nativo con etiqueta (DESIGN §7.3) ──
export function Selector({ etiqueta, valor, opciones, alCambiar, className, vacio }: { etiqueta: string; valor: string; opciones: { valor: string; texto: string }[]; alCambiar: (v: string) => void; className?: string; vacio?: string }) {
  return (
    <label className={cn('flex flex-col gap-1 text-small text-text-secondary', className)}>
      {etiqueta}
      <select value={valor} onChange={(e) => alCambiar(e.target.value)} className={cn(estilosControl, 'w-full')}>
        {vacio !== undefined && <option value="">{vacio}</option>}
        {opciones.map((o) => <option key={o.valor} value={o.valor}>{o.texto}</option>)}
      </select>
    </label>
  )
}

/** Título de bloque dentro de la pantalla (h2), con acción opcional a la derecha. */
export function Bloque({ titulo, accion, children }: { titulo: string; accion?: ReactNode; children: ReactNode }) {
  return (
    <section className="flex flex-col gap-3">
      <div className="flex flex-wrap items-center justify-between gap-2"><h2 className="text-title-3 font-semibold text-text-default">{titulo}</h2>{accion}</div>
      {children}
    </section>
  )
}

/** Confirmación de una acción destructiva (U1 §2.2: «Eliminar (con confirmación)»). */
export function Confirmar({ abierta, alCambiar, titulo, texto, accion, alConfirmar }: { abierta: boolean; alCambiar: (a: boolean) => void; titulo: string; texto: string; accion: string; alConfirmar: () => Promise<void> }) {
  const [enCurso, setEnCurso] = useState(false)
  return (
    <Hoja abierta={abierta} alCambiar={alCambiar} titulo={titulo} descripcion={texto}
      pie={<div className="flex justify-end gap-2"><Boton onClick={() => alCambiar(false)}>{textos.app.cancelar}</Boton><Boton variante="destructivo" aria-busy={enCurso} onClick={async () => { setEnCurso(true); try { await alConfirmar() } finally { setEnCurso(false) } }}>{accion}</Boton></div>}>
      <p className="text-body text-text-default">{texto}</p>
    </Hoja>
  )
}

/** Fila de propiedad: etiqueta a la izquierda y valor a la derecha (móvil en dos renglones). */
export function Fila({ etiqueta, children }: { etiqueta: string; children: ReactNode }) {
  return (
    <div className="flex flex-col gap-1 border-b border-border-divider py-2 tablet:flex-row tablet:items-center tablet:gap-4">
      <span className="text-small text-text-secondary tablet:w-48 tablet:shrink-0">{etiqueta}</span>
      <div className="min-h-11 flex-1 text-body text-text-default tabular-nums tablet:min-h-9 coarse:min-h-11 flex items-center">{children}</div>
    </div>
  )
}
export const rejillaTarjetas = 'grid grid-cols-1 gap-3 tablet:grid-cols-2 desktop:grid-cols-3'
/** Centavos → texto editable «1234.56» sin flotantes (protocolo §6). */
export const centavosATexto = (c: number) => `${c < 0 ? '-' : ''}${Math.floor(Math.abs(c) / 100)}.${String(Math.abs(c) % 100).padStart(2, '0')}`

/** Hoja de captura con monto, fecha y campos extra (P-08 quincena/extra, P-09 movimiento). Valida monto > 0 (C1 MoneyCents). */
export function HojaCaptura({ titulo, abierta, alCerrar, montoInicial = '', fechaInicial, extra, alGuardar }: { titulo: string; abierta: boolean; alCerrar: () => void; montoInicial?: string; fechaInicial: string; extra?: ReactNode; alGuardar: (c: { monto: number; fecha: string; nota: string }) => Promise<void> }) {
  const [v, setV] = useState({ monto: montoInicial, fecha: fechaInicial, nota: '', error: '', general: '' })
  const [enCurso, setEnCurso] = useState(false)
  useEffect(() => { if (abierta) setV({ monto: montoInicial, fecha: fechaInicial, nota: '', error: '', general: '' }) }, [abierta, montoInicial, fechaInicial])
  const sinRed = useSinRed()
  const guardar = async () => {
    const m = aCentavosLocal(v.monto)
    if (m == null || m <= 0) return setV({ ...v, error: textos.registro.montoInvalido })
    setEnCurso(true)
    try { await alGuardar({ monto: m, fecha: v.fecha, nota: v.nota }); alCerrar() } catch (e) { setV({ ...v, general: mensajeDe(e) }) } finally { setEnCurso(false) }
  }
  return (
    <Hoja abierta={abierta} alCambiar={(a) => !a && alCerrar()} titulo={titulo}
      pie={<div className="flex justify-end gap-2"><Boton onClick={alCerrar}>{textos.app.cancelar}</Boton><Boton variante="primario" aria-busy={enCurso} deshabilitadoPor={sinRed} onClick={guardar}>{textos.app.guardar}</Boton></div>}>
      <div className="flex flex-col gap-3">
        {extra}
        <CampoTexto etiqueta={t.comun.monto} inputMode="decimal" value={v.monto} error={v.error || undefined} onChange={(e) => setV({ ...v, monto: e.target.value, error: '' })} />
        <CampoTexto etiqueta={t.comun.fecha} type="date" value={v.fecha} onChange={(e) => setV({ ...v, fecha: e.target.value })} />
        <CampoTexto etiqueta={t.comun.nota} value={v.nota} onChange={(e) => setV({ ...v, nota: e.target.value })} />
        {v.general && <p role="alert" className="text-small text-role-error">{v.general}</p>}
      </div>
    </Hoja>
  )
}
