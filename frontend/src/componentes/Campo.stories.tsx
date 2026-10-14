// Variantes del campo de formulario (DESIGN.md §7.3): reposo, ayuda, error, deshabilitado y superficies.
import type { Meta, StoryObj } from '@storybook/react-vite'
import { Campo } from './Campo'
import { textos } from '../textos/es-MX'
import { Grupo, Superficies } from '../../.storybook/apoyo'

const meta = {
  title: 'Componentes/Campo',
  component: Campo,
  args: { etiqueta: textos.registro.monto, placeholder: '0.00', inputMode: 'decimal' },
  argTypes: { error: { control: 'text' }, ayuda: { control: 'text' }, disabled: { control: 'boolean' } },
  decorators: [(Historia) => <div className="max-w-sm"><Historia /></div>],
} satisfies Meta<typeof Campo>
export default meta
type Historia = StoryObj<typeof meta>

/** Editable con los controles de Storybook. */
export const Editable: Historia = {}

export const Estados: Historia = {
  render: () => (
    <div className="flex max-w-sm flex-col gap-4">
      <Grupo titulo="Reposo"><Campo etiqueta={textos.registro.descripcion} placeholder="Tortillas" /></Grupo>
      <Grupo titulo="Con valor"><Campo etiqueta={textos.registro.monto} defaultValue="1,234.56" inputMode="decimal" /></Grupo>
      <Grupo titulo="Con ayuda"><Campo etiqueta={textos.registro.pesoKg} placeholder="72.5" inputMode="decimal" ayuda="Escribe el peso en kilos." /></Grupo>
      <Grupo titulo="Con error"><Campo etiqueta={textos.registro.monto} defaultValue="0" inputMode="decimal" error={textos.registro.montoInvalido} /></Grupo>
      <Grupo titulo="Etiqueta oculta"><Campo etiqueta={textos.menu.titulo} ocultarEtiqueta placeholder={textos.menu.placeholder} /></Grupo>
    </div>
  ),
}

/** Hallazgo de la revisión: Campo no trae estilos de deshabilitado (DESIGN.md §7.3 pide relleno ui.surface.disabled, texto ui.text.disabled y borde punteado). */
export const Deshabilitado: Historia = {
  name: 'Deshabilitado',
  tags: ['hallazgo'],
  args: { disabled: true, defaultValue: 'No editable' },
}

export const Superficie: Historia = {
  name: 'Sobre página, tarjeta y superficie flotante',
  render: () => (
    <Superficies>
      {(nombre) => (
        <div className="flex max-w-sm flex-col gap-3">
          <Campo etiqueta={textos.registro.monto} placeholder="0.00" inputMode="decimal" />
          <Campo etiqueta={textos.registro.descripcion} defaultValue={`Mercado (${nombre})`} error={textos.errores.revisaCampos} />
        </div>
      )}
    </Superficies>
  ),
}
