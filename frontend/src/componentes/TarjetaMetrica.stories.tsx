// Variantes de la tarjeta de métrica (DESIGN.md §7.6): cifra neutra y de rol, nota, interactiva y ancho de 360.
import type { Meta, StoryObj } from '@storybook/react-vite'
import { ArrowDownRight, ArrowUpRight, PiggyBank, Receipt, TriangleAlert, Wallet } from 'lucide-react'
import { TarjetaMetrica } from './TarjetaMetrica'
import { dinero } from '../api/formatos'
import { Grupo } from '../../.storybook/apoyo'

const meta = {
  id: 'tarjeta-metrica',
  title: 'Componentes/Tarjeta de métrica',
  component: TarjetaMetrica,
  args: { titulo: 'Gasto del mes', valor: dinero(824010, 'gasto'), nota: '12 % más que septiembre', rol: 'expense' },
  argTypes: { rol: { control: 'select', options: [undefined, 'income', 'expense', 'savings', 'debt', 'warning'] } },
  decorators: [(Historia) => <div className="max-w-xs"><Historia /></div>],
} satisfies Meta<typeof TarjetaMetrica>
export default meta
type Historia = StoryObj<typeof meta>

/** Editable con los controles de Storybook. */
export const Editable: Historia = {}

export const Roles: Historia = {
  name: 'Cifra neutra y de rol',
  render: () => (
    <div className="grid max-w-3xl grid-cols-1 gap-4 tablet:grid-cols-2 desktop:grid-cols-3">
      <TarjetaMetrica titulo="Saldo disponible" valor={dinero(1000000)} nota="Cuenta A" icono={<Wallet aria-hidden className="size-5" />} />
      <TarjetaMetrica titulo="Ingresos del mes" valor={dinero(3250000, 'ingreso')} nota={<span className="inline-flex items-center gap-1"><ArrowUpRight aria-hidden className="size-5" />8 % más que septiembre</span>} rol="income" />
      <TarjetaMetrica titulo="Gasto del mes" valor={dinero(824010, 'gasto')} nota={<span className="inline-flex items-center gap-1"><ArrowDownRight aria-hidden className="size-5" />12 % menos que septiembre</span>} rol="expense" icono={<Receipt aria-hidden className="size-5" />} />
      <TarjetaMetrica titulo="Ahorro" valor={dinero(500000)} nota="Meta del año" rol="savings" icono={<PiggyBank aria-hidden className="size-5" />} />
      <TarjetaMetrica titulo="Deuda total" valor={dinero(9500050)} nota="Pago el 30-oct-2026" rol="debt" />
      <TarjetaMetrica titulo="Presupuesto de comida" valor="82 % usado" nota="Falta poco para el límite" rol="warning" icono={<TriangleAlert aria-hidden className="size-5" />} />
    </div>
  ),
}

export const Interactiva: Historia = {
  name: 'Interactiva (abre el detalle)',
  render: () => (
    <div className="flex max-w-xs flex-col gap-4">
      <Grupo titulo="Tarjeta enlace" nota="Hover, foco y presionado se ven con el puntero y el teclado.">
        <TarjetaMetrica titulo="Gasto del mes" valor={dinero(824010, 'gasto')} nota="12 % más que septiembre" rol="expense" alAbrir={() => undefined} />
      </Grupo>
    </div>
  ),
}

/** `$99,999.99` debe caber completo a 360 px (P5): 360 menos 16 px de margen por lado. */
export const ImporteLargo: Historia = {
  name: 'Importe largo a 360 px',
  parameters: { viewport: { defaultViewport: 'm360' } },
  render: () => (
    <div className="w-[328px]">
      <TarjetaMetrica titulo="Gasto del mes" valor={dinero(9999999, 'gasto')} nota="12 % más que septiembre" rol="expense" />
    </div>
  ),
}
