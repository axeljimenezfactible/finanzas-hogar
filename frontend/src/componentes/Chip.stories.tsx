// Variantes de chips de persona y tienda (DESIGN.md §7.4, K41): forma + nombre + posición, nunca solo color.
import { useState } from 'react'
import type { Meta, StoryObj } from '@storybook/react-vite'
import { ChipPersona, ChipSeleccionable, ChipTienda, MarcaPersona, MarcaTienda, ORDEN_PERSONAS, ORDEN_TIENDAS } from './Chip'
import { Etiquetada, Fila, Grupo, Superficies } from '../../.storybook/apoyo'

const meta = {
  title: 'Componentes/Chip',
  component: ChipPersona,
  args: { codigo: 'axel', nombre: 'Axel' },
  argTypes: { codigo: { control: 'select', options: ORDEN_PERSONAS } },
} satisfies Meta<typeof ChipPersona>
export default meta
type Historia = StoryObj<typeof meta>

const PERSONAS: Record<string, string> = { axel: 'Axel', ximena: 'Ximena', karen: 'Karen', shared: 'Conjunto' }
const TIENDAS: Record<string, string> = { walmart: 'Walmart', 'bodega-aurrera': 'Bodega Aurrera', 'sams-club': "Sam's Club", market: 'Mercado' }

/** Editable con los controles de Storybook. */
export const Editable: Historia = {}

export const PersonasYTiendas: Historia = {
  name: 'Personas y tiendas (orden fijo)',
  render: () => (
    <div className="flex flex-col gap-6">
      <Grupo titulo="Personas" nota="Círculo relleno; Conjunto, anillo hueco. Orden: Axel, Ximena, Karen, Conjunto.">
        <Fila>{ORDEN_PERSONAS.map((c) => <ChipPersona key={c} codigo={c} nombre={PERSONAS[c]} />)}</Fila>
      </Grupo>
      <Grupo titulo="Tiendas" nota="Cuadrado relleno. Orden: Walmart, Bodega Aurrera, Sam's Club, Mercado.">
        <Fila>{ORDEN_TIENDAS.map((c) => <ChipTienda key={c} codigo={c} nombre={TIENDAS[c]} />)}</Fila>
      </Grupo>
      <Grupo titulo="Solo la marca de forma" nota="Decorativa (aria-hidden): el nombre en texto es el portador.">
        <Fila>
          {ORDEN_PERSONAS.map((c) => <MarcaPersona key={c} codigo={c} />)}
          {ORDEN_TIENDAS.map((c) => <MarcaTienda key={c} codigo={c} />)}
        </Fila>
      </Grupo>
    </div>
  ),
}

function Alternable({ forma, conMarca }: { forma: 'pildora' | 'rectangulo'; conMarca?: boolean }) {
  const [activo, setActivo] = useState(false)
  return (
    <ChipSeleccionable seleccionado={activo} forma={forma} onClick={() => setActivo(!activo)} marca={conMarca ? <MarcaTienda codigo="walmart" /> : undefined}>
      {conMarca ? 'Walmart' : 'Comida'}
    </ChipSeleccionable>
  )
}

export const Seleccionables: Historia = {
  render: () => (
    <div className="flex flex-col gap-6">
      <Grupo titulo="Píldora" nota="Alto 32 px con puntero fino y 44 px con puntero grueso. Tocar alterna aria-pressed.">
        <Fila>
          <Etiquetada etiqueta="Sin elegir"><ChipSeleccionable seleccionado={false}>Comida</ChipSeleccionable></Etiquetada>
          <Etiquetada etiqueta="Elegido"><ChipSeleccionable seleccionado>Comida</ChipSeleccionable></Etiquetada>
          <Etiquetada etiqueta="Deshabilitado"><ChipSeleccionable seleccionado={false} disabled>Comida</ChipSeleccionable></Etiquetada>
          <Etiquetada etiqueta="Alternable"><Alternable forma="pildora" /></Etiquetada>
        </Fila>
      </Grupo>
      <Grupo titulo="Rectángulo con marca (tienda)">
        <Fila>
          <Etiquetada etiqueta="Sin elegir"><ChipSeleccionable seleccionado={false} forma="rectangulo" marca={<MarcaTienda codigo="walmart" />}>Walmart</ChipSeleccionable></Etiquetada>
          <Etiquetada etiqueta="Elegido"><ChipSeleccionable seleccionado forma="rectangulo" marca={<MarcaTienda codigo="walmart" />}>Walmart</ChipSeleccionable></Etiquetada>
          <Etiquetada etiqueta="Alternable"><Alternable forma="rectangulo" conMarca /></Etiquetada>
        </Fila>
      </Grupo>
    </div>
  ),
}

export const Superficie: Historia = {
  name: 'Sobre página, tarjeta, superficie flotante y fila elegida',
  render: () => (
    <div className="flex flex-col gap-4">
      <Superficies>
        {() => (
          <Fila>
            {ORDEN_PERSONAS.map((c) => <ChipPersona key={c} codigo={c} nombre={PERSONAS[c]} />)}
            {ORDEN_TIENDAS.map((c) => <ChipTienda key={c} codigo={c} nombre={TIENDAS[c]} />)}
          </Fila>
        )}
      </Superficies>
      <div className="flex flex-col gap-2 rounded-md bg-surface-selected p-3">
        <span className="text-caption text-text-secondary">Sobre una fila elegida</span>
        <Fila>{ORDEN_PERSONAS.map((c) => <ChipPersona key={c} codigo={c} nombre={PERSONAS[c]} />)}</Fila>
      </div>
    </div>
  ),
}
