// Variantes de estados vacío, carga, error y franjas (DESIGN.md §7.14, G-07).
import type { Meta, StoryObj } from '@storybook/react-vite'
import { Receipt } from 'lucide-react'
import { Boton } from './Boton'
import { Cargando, ConDatos, Esqueleto, ErrorCarga, Franja, Vacio } from './Estados'
import { ErrorApi, ErrorSinConexion } from '../api/errores'
import type { EstadoConsulta } from '../api/consulta'
import { textos } from '../textos/es-MX'
import { Grupo } from '../../.storybook/apoyo'

const meta = { title: 'Componentes/Estados', component: Franja, args: { rol: 'info', children: 'Aviso' } } satisfies Meta<typeof Franja>
export default meta
type Historia = StoryObj<typeof meta>

const consulta = (p: Partial<EstadoConsulta<string[]>>): EstadoConsulta<string[]> => ({ datos: undefined, cargando: false, error: null, desdeCopia: null, nuncaCargada: false, recargar: () => undefined, ...p })

export const Franjas: Historia = {
  render: () => (
    <div className="flex max-w-2xl flex-col gap-3">
      <Franja rol="info">{textos.registro.ingresoIr}</Franja>
      <Franja rol="advertencia" accion={<Boton variante="enlace">Capturar</Boton>}>Falta capturar el precio de la leche</Franja>
      <Franja rol="error" alerta accion={<Boton variante="enlace">{textos.app.reintentar}</Boton>}>{textos.errores.generico}</Franja>
      <Franja rol="sinConexion" alerta accion={<Boton variante="enlace">{textos.conexion.revisar}</Boton>}>{textos.conexion.sinConexion} · {textos.conexion.enCola(2)}</Franja>
      <Franja rol="info" accion={<Boton variante="enlace">{textos.app.reintentar}</Boton>}>{textos.app.tardando}</Franja>
    </div>
  ),
}

export const Carga: Historia = {
  name: 'Carga (esqueleto)',
  render: () => (
    <div className="flex max-w-md flex-col gap-6">
      <Grupo titulo="Esqueleto por omisión" nota="A los 10 s agrega «Esto está tardando más de lo normal · Reintentar»."><Cargando alReintentar={() => undefined} /></Grupo>
      <Grupo titulo="Esqueleto con la forma final"><Cargando><Esqueleto className="h-24" /><Esqueleto className="w-1/3" /></Cargando></Grupo>
    </div>
  ),
}

export const Vacios: Historia = {
  name: 'Vacío',
  render: () => (
    <Vacio icono={<Receipt />} texto="Aún no hay gastos este mes" accion={<Boton variante="primario">{textos.barra.registrarGasto}</Boton>} />
  ),
}

export const Errores: Historia = {
  render: () => (
    <div className="flex max-w-2xl flex-col gap-3">
      <ErrorCarga error={new ErrorSinConexion()} alReintentar={() => undefined} />
      <ErrorCarga error={new ErrorApi(404, 'not_found', 'No existe')} alReintentar={() => undefined} />
      <ErrorCarga error={new ErrorApi(409, 'conflict', 'Este registro cambió. Recarga para ver la versión nueva.')} alReintentar={() => undefined} />
      <ErrorCarga error={new Error('x')} />
    </div>
  ),
}

export const ConDatosEstados: Historia = {
  name: 'ConDatos (todos los estados)',
  render: () => (
    <div className="flex max-w-xl flex-col gap-6">
      <Grupo titulo="Cargando"><ConDatos consulta={consulta({ cargando: true })}>{() => null}</ConDatos></Grupo>
      <Grupo titulo="Sin conexión y nunca cargada"><ConDatos consulta={consulta({ nuncaCargada: true })}>{() => null}</ConDatos></Grupo>
      <Grupo titulo="Con error"><ConDatos consulta={consulta({ error: new ErrorSinConexion() })}>{() => null}</ConDatos></Grupo>
      <Grupo titulo="Con datos de la copia del teléfono"><ConDatos consulta={consulta({ datos: ['Tortillas', 'Leche'], desdeCopia: Date.UTC(2026, 9, 12, 15, 30) })}>{(d) => <p className="text-body text-text-default">{d.join(', ')}</p>}</ConDatos></Grupo>
      <Grupo titulo="Con datos"><ConDatos consulta={consulta({ datos: ['Tortillas', 'Leche'] })}>{(d) => <p className="text-body text-text-default">{d.join(', ')}</p>}</ConDatos></Grupo>
    </div>
  ),
}
