// Estados comunes (U1 §2.4, DESIGN §7.14): carga con esqueleto, vacío, error, sin conexión y franjas.
import { useEffect, useState, type ReactNode } from 'react'
import { CircleAlert, Info, RefreshCw, TriangleAlert, WifiOff } from 'lucide-react'
import { textos } from '../textos/es-MX'
import { ErrorApi, esSinConexion } from '../api/errores'
import { horaCorta } from '../api/formatos'
import type { EstadoConsulta } from '../api/consulta'
import { Boton } from './Boton'
import { cn } from './utils'
export const Esqueleto = ({ className }: { className?: string }) => <div aria-hidden className={cn('h-6 rounded-md bg-skeleton-on-page', className)} />
export function Cargando({ alReintentar, children }: { alReintentar?: () => void; children?: ReactNode }) {
  const [tarde, setTarde] = useState(false)
  useEffect(() => { const t = setTimeout(() => setTarde(true), 10_000); return () => clearTimeout(t) }, [])
  return (
    <div aria-busy="true" aria-label={textos.app.cargando} className="flex flex-col gap-3">
      {children ?? <><Esqueleto className="w-2/3" /><Esqueleto /><Esqueleto className="w-1/2" /></>}
      {tarde && <Franja rol="info" accion={alReintentar && <Boton variante="enlace" onClick={alReintentar}>{textos.app.reintentar}</Boton>}>{textos.app.tardando}</Franja>}
    </div>
  )
}
const ROL = {
  info: { icono: Info, clase: 'bg-role-info-surface', color: 'text-role-info' },
  advertencia: { icono: TriangleAlert, clase: 'bg-role-warning-surface', color: 'text-role-warning' },
  error: { icono: CircleAlert, clase: 'bg-role-error-surface', color: 'text-role-error' },
  sinConexion: { icono: WifiOff, clase: 'bg-role-warning-surface', color: 'text-role-warning' },
}
/** Franja tipo «callout» de Notion. */
export function Franja({ rol, children, accion, alerta }: { rol: keyof typeof ROL; children: ReactNode; accion?: ReactNode; alerta?: boolean }) {
  const { icono: Icono, clase, color } = ROL[rol]
  return (
    <div role={alerta ? 'alert' : 'status'} className={cn('flex items-center gap-2 rounded-md p-3 text-small text-text-default', clase)}>
      <Icono aria-hidden className={cn('size-5 shrink-0', color)} />
      <span className="flex-1">{children}</span>
      {accion}
    </div>
  )
}
export function Vacio({ icono, texto, accion }: { icono?: ReactNode; texto: string; accion?: ReactNode }) {
  return (
    <div className="flex flex-col items-center gap-3 py-10 text-center text-text-secondary">
      {icono && <span aria-hidden className="text-icon-default [&_svg]:size-6">{icono}</span>}
      <p className="text-body">{texto}</p>
      {accion}
    </div>
  )
}
export function ErrorCarga({ error, alReintentar }: { error: unknown; alReintentar?: () => void }) {
  const texto = esSinConexion(error) ? textos.errores.red : error instanceof ErrorApi ? (error.estado === 404 ? textos.errores.noExiste : error.message) : textos.errores.generico
  return <Franja rol="error" alerta accion={alReintentar && <Boton variante="enlace" onClick={alReintentar}><RefreshCw aria-hidden className="size-5" />{textos.app.reintentar}</Boton>}>{texto}</Franja>
}
/** Pinta una consulta con todos sus estados; `children` recibe los datos ya cargados. */
export function ConDatos<T>({ consulta, children, esqueleto }: { consulta: EstadoConsulta<T>; children: (d: T) => ReactNode; esqueleto?: ReactNode }) {
  if (consulta.cargando && consulta.datos === undefined) return <Cargando alReintentar={consulta.recargar}>{esqueleto}</Cargando>
  if (consulta.nuncaCargada) return <Franja rol="sinConexion">{textos.conexion.paginaNoCargada}</Franja>
  if (consulta.error || consulta.datos === undefined) return <ErrorCarga error={consulta.error} alReintentar={consulta.recargar} />
  return (
    <>
      {consulta.desdeCopia && <p className="text-caption text-text-secondary">{textos.conexion.cargadoA(horaCorta(consulta.desdeCopia))}</p>}
      {children(consulta.datos)}
    </>
  )
}
