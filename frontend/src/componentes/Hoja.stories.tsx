// Variantes de la hoja y el diálogo (DESIGN.md §7.12): pantalla completa en móvil, ventana modal desde 768 px.
import { useState } from 'react'
import type { Meta, StoryObj } from '@storybook/react-vite'
import { Hoja } from './Hoja'
import { Boton } from './Boton'
import { Campo } from './Campo'
import { textos } from '../textos/es-MX'

const meta = { title: 'Componentes/Hoja', component: Hoja, parameters: { layout: 'fullscreen' } } satisfies Meta<typeof Hoja>
export default meta
type Historia = StoryObj<typeof meta>

const pie = <Boton variante="primario" tam="hoja">{textos.app.guardar}</Boton>

/** Abierta desde el inicio, para verla y medirla. */
export const Abierta: Historia = {
  args: { abierta: true, alCambiar: () => undefined, titulo: textos.registro.titulo, children: null },
  render: (args) => (
    <Hoja {...args} pie={pie}>
      <div className="flex flex-col gap-4">
        <Campo etiqueta={textos.registro.monto} placeholder="0.00" inputMode="decimal" />
        <Campo etiqueta={textos.registro.descripcion} placeholder="Tortillas" />
      </div>
    </Hoja>
  ),
}

export const ConDescripcion: Historia = {
  name: 'Con descripción y error en el campo',
  args: { abierta: true, alCambiar: () => undefined, titulo: textos.registro.pesoKg, descripcion: 'Escribe tu peso en kilos.', children: null },
  render: (args) => (
    <Hoja {...args} pie={pie}>
      <Campo etiqueta={textos.registro.pesoKg} defaultValue="0" inputMode="decimal" error={textos.registro.pesoInvalido} />
    </Hoja>
  ),
}

function Disparador() {
  const [abierta, setAbierta] = useState(false)
  return (
    <div className="p-4">
      <Boton variante="primario" onClick={() => setAbierta(true)}>{textos.barra.registrarGasto}</Boton>
      <Hoja abierta={abierta} alCambiar={setAbierta} titulo={textos.registro.titulo} pie={pie}>
        <Campo etiqueta={textos.registro.monto} placeholder="0.00" inputMode="decimal" />
      </Hoja>
    </div>
  )
}
/** Se abre con el botón; Esc la cierra y devuelve el foco al disparador. */
export const Interactiva: Historia = { args: { abierta: false, alCambiar: () => undefined, titulo: '', children: null }, render: () => <Disparador /> }
