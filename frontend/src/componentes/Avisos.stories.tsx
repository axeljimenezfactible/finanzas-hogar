// Variantes de los avisos breves (DESIGN.md §7.15): éxito, información, advertencia y error, con y sin acción.
import { useEffect } from 'react'
import type { Meta, StoryObj } from '@storybook/react-vite'
import { ProveedorAvisos, useAvisos, type Aviso } from './Avisos'
import { Boton } from './Boton'
import { textos } from '../textos/es-MX'

/** Dispara el aviso al montar, para verlo (y medirlo) sin tocar nada. */
function Disparar({ aviso }: { aviso: Aviso }) {
  const mostrar = useAvisos()
  useEffect(() => { mostrar(aviso) }, [mostrar, aviso])
  return <p className="text-small text-text-secondary">El aviso breve aparece abajo (a la derecha desde 768 px).</p>
}
const deshacer = { texto: textos.app.deshacer, alHacer: () => undefined }
const historia = (aviso: Aviso) => ({ render: () => <ProveedorAvisos><Disparar aviso={aviso} /></ProveedorAvisos> })

const meta = { title: 'Componentes/Avisos', component: ProveedorAvisos, args: { children: null } } satisfies Meta<typeof ProveedorAvisos>
export default meta
type Historia = StoryObj<typeof meta>

export const Exito: Historia = { name: 'Éxito con Deshacer', ...historia({ texto: 'Guardado', rol: 'exito', accion: deshacer }) }
export const ExitoEnCola: Historia = { name: 'Éxito sin conexión (texto largo)', ...historia({ texto: textos.registro.gastoEnCola, rol: 'exito', accion: deshacer }) }
export const Informacion: Historia = { name: 'Información', ...historia({ texto: textos.registro.ingresoIr, rol: 'info' }) }
export const Advertencia: Historia = { name: 'Advertencia', ...historia({ texto: textos.conexion.sinConexion, rol: 'advertencia' }) }
export const Error: Historia = { name: 'Error (sin cierre automático)', ...historia({ texto: textos.errores.generico, rol: 'error' }) }

function Probador() {
  const mostrar = useAvisos()
  return (
    <div className="flex flex-wrap gap-2">
      <Boton variante="primario" onClick={() => mostrar({ texto: 'Guardado', rol: 'exito', accion: deshacer })}>Mostrar éxito</Boton>
      <Boton onClick={() => mostrar({ texto: textos.errores.generico, rol: 'error' })}>Mostrar error</Boton>
    </div>
  )
}
/** Uno visible a la vez y los demás en cola; el de éxito dura 5 s y se pausa con hover o foco. */
export const Interactivo: Historia = { render: () => <ProveedorAvisos><Probador /></ProveedorAvisos> }
