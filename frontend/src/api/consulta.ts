// Hook de lectura con estados (carga, error, sin conexión con copia o sin copia) y recarga.
import { useCallback, useEffect, useRef, useState } from 'react'
import { armarUrl, pedir, type Opciones } from './cliente'
import { esSinConexion } from './errores'
import { leerLectura } from './lecturas'
import type { Respuesta, RutasCon } from './tipos'

export interface EstadoConsulta<T> {
  datos: T | undefined
  cargando: boolean
  error: unknown
  /** Si los datos vienen de la última carga del teléfono (sin conexión): cuándo se cargaron. */
  desdeCopia: number | null
  /** true cuando no hay conexión y esta lectura nunca se cargó en el teléfono. */
  nuncaCargada: boolean
  recargar: () => void
}

export function usarConsulta<P extends RutasCon<'get'>>(metodo: 'get', ruta: P, opciones?: Opciones<P, 'get'>, activa = true): EstadoConsulta<Respuesta<P, 'get'>> {
  const clave = JSON.stringify([ruta, opciones ?? {}])
  const [estado, setEstado] = useState<Omit<EstadoConsulta<Respuesta<P, 'get'>>, 'recargar'>>({ datos: undefined, cargando: activa, error: null, desdeCopia: null, nuncaCargada: false })
  const [vuelta, setVuelta] = useState(0)
  const opRef = useRef(opciones)
  opRef.current = opciones
  useEffect(() => {
    if (!activa) return
    const ctl = new AbortController()
    setEstado((e) => ({ ...e, cargando: true, error: null }))
    pedir(metodo, ruta, { ...(opRef.current as object), senal: ctl.signal } as Opciones<P, 'get'>).then(
      (datos) => setEstado({ datos, cargando: false, error: null, desdeCopia: null, nuncaCargada: false }),
      async (error) => {
        if (ctl.signal.aborted) return
        if (esSinConexion(error)) {
          const o = opRef.current as { ruta?: object; consulta?: object } | undefined
          const copia = await leerLectura(armarUrl(ruta as string, o?.ruta, o?.consulta))
          if (copia) return setEstado({ datos: copia.datos as Respuesta<P, 'get'>, cargando: false, error: null, desdeCopia: copia.cargadoEn, nuncaCargada: false })
          return setEstado({ datos: undefined, cargando: false, error, desdeCopia: null, nuncaCargada: true })
        }
        setEstado({ datos: undefined, cargando: false, error, desdeCopia: null, nuncaCargada: false })
      },
    )
    const recuperada = () => setVuelta((v) => v + 1)
    window.addEventListener('fh:conexion-recuperada', recuperada)
    return () => { ctl.abort(); window.removeEventListener('fh:conexion-recuperada', recuperada) }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [clave, vuelta, activa])
  const recargar = useCallback(() => setVuelta((v) => v + 1), [])
  return { ...estado, recargar }
}
