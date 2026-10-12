import { useEffect } from 'react'
import { ProveedorAvisos, useAvisos } from '../componentes/Avisos'
import { ProveedorIconos } from '../componentes/Icono'
import { Cargando } from '../componentes/Estados'
import { iniciarCola, type ResultadoEnvio } from '../cola/cola'
import { textos } from '../textos/es-MX'
import { ProveedorSesion, useSesion } from './sesion'
import { IniciarSesion } from './IniciarSesion'
import { navegar, useUbicacion } from './enrutador'
import { Layout } from '../shell/Layout'
function Raiz() {
  const { estado } = useSesion()
  const { ruta } = useUbicacion()
  const avisar = useAvisos()
  useEffect(() => {
    iniciarCola()
    const alEnviar = (e: Event) => {
      const r = (e as CustomEvent<ResultadoEnvio>).detail
      if (r.enviadas) avisar({ texto: textos.conexion.enviadas(r.enviadas), rol: 'exito' })
      // U5-13 paso 7: tanto el rechazo (4xx) como la falla del servidor o la sesión vencida (5xx, 401) avisan; la captura sigue en la cola.
      if (r.fallidas || r.retenidas) avisar({ texto: `${textos.conexion.noSeEnvio} · ${textos.conexion.revisar}`, rol: 'error' })
    }
    window.addEventListener('fh:cola-enviada', alEnviar)
    return () => window.removeEventListener('fh:cola-enviada', alEnviar)
  }, [avisar])
  useEffect(() => { if (estado === 'activa' && (ruta === '/' || ruta === '/entrar')) navegar('/inicio', { reemplazar: true }) }, [estado, ruta])
  if (estado === 'cargando') return <div className="p-6"><Cargando /></div>
  if (estado === 'sin-sesion') return <IniciarSesion />
  return <Layout />
}
export function App() {
  return <ProveedorIconos><ProveedorAvisos><ProveedorSesion><Raiz /></ProveedorSesion></ProveedorAvisos></ProveedorIconos>
}
